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
		skipOffset(rs, opts)
	} else {
		seekFromStart(rs, opts.Offset)
	}
	startCopy(rs, w, opts)
}

func isStdin(rs io.ReadSeeker) bool {
	file, ok := rs.(*os.File)
	if !ok || file.Fd() != uintptr(syscall.Stdin) {
		return false
	}
	return true
}

func startCopy(rs io.ReadSeeker, w io.Writer, opts *Options) {
	block := make([]byte, opts.BlockSize)

	for {
		readBytesCount := readBlock(rs, block)
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

func skipOffset(rs io.ReadSeeker, opts *Options) {
	block := make([]byte, opts.BlockSize)

	remainingBytesCount := opts.Offset

	for remainingBytesCount > opts.BlockSize {
		readBytesCount := readBlock(rs, block)
		if !noMoreBytes(readBytesCount, opts.BlockSize) {
			return
		}
		remainingBytesCount -= readBytesCount
	}

	readBlock(rs, make([]byte, remainingBytesCount))
}

func noMoreBytes(readBytesCount, blockSize int) bool {
	if readBytesCount < blockSize {
		fmt.Fprintln(os.Stderr, "offset is greater than input data size")
		return false
	}
	return true
}

func seekFromStart(rs io.Seeker, offset int) {
	_, err := rs.Seek(int64(offset), io.SeekStart)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func readBlock(r io.Reader, block []byte) int {
	readBytesCount, err := r.Read(block)
	if err != nil && err != io.EOF {
		fmt.Fprintln(os.Stderr, err)
	}
	return readBytesCount
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
