package dd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"
)

func Start(opts *Options) {
	r := readSeekCloser(opts.From)
	w := writerCloser(opts.To)

	process(r, w, opts)

	closeStream(r)
	closeStream(w)
}

func readSeekCloser(from string) io.ReadSeekCloser {
	if from == Stdin {
		return os.Stdin
	}
	file, _ := os.Open(from)
	return file
}

func writerCloser(to string) io.WriteCloser {
	if to == Stdout {
		return os.Stdout
	}
	file, _ := os.Create(to)
	return file
}

func process(rs io.ReadSeeker, w io.Writer, opts *Options) {
	if isStdin(rs) {
		processFromStdin(rs, w, opts)
	} else {
		processFromFile(rs, w, opts)
	}
}

func isStdin(rs io.ReadSeeker) bool {
	file, ok := rs.(*os.File)
	if !ok || file.Fd() != uintptr(syscall.Stdin) {
		return false
	}
	return true
}

func processFromStdin(rs io.ReadSeeker, w io.Writer, opts *Options) {
	block := make([]byte, opts.BlockSize)

	skipToOffset(rs, opts)

	for {
		readBytesCount, _ := rs.Read(block)
		if readBytesCount == 0 {
			break
		}
		if readBytesCount < opts.BlockSize {
			block = block[:readBytesCount]
		}

		block = convertBlock(block, opts.Conv)
		writeBlock(w, block)
	}
}

func processFromFile(rs io.ReadSeeker, w io.Writer, opts *Options) {
	block := make([]byte, opts.BlockSize)

	seekFromStart(rs, opts.Offset)

	for {
		readBytesCount, _ := rs.Read(block)
		if readBytesCount == 0 {
			break
		}
		if readBytesCount < opts.BlockSize {
			block = block[:readBytesCount]
		}

		block = convertBlock(block, opts.Conv)
		writeBlock(w, block)
	}
}

func skipToOffset(rs io.ReadSeeker, opts *Options) {
	block := make([]byte, opts.BlockSize)

	diff := opts.Offset

	for diff > opts.BlockSize {
		readBytesCount, _ := rs.Read(block)

		if readBytesCount < opts.BlockSize {
			fmt.Fprintln(os.Stderr, "offset is greater than input data size")
			return
		}

		diff -= readBytesCount
	}

	rs.Read(make([]byte, 0, diff))
}

func seekFromStart(rs io.Seeker, offset int) {
	_, err := rs.Seek(int64(offset), io.SeekStart)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func writeBlock(w io.Writer, block []byte) {
	_, err := w.Write(block)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func closeStream(c io.Closer) {
	if err := c.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func convertBlock(block []byte, conv string) []byte {
	str := string(block)

	switch conv {
	case UpperCase:
		str = strings.ToUpper(str)
	case LowerCase:
		str = strings.ToLower(str)
	}

	return []byte(str)
}
