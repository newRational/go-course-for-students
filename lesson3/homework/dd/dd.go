package dd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

func Start(opts *Options) {
	r := closerReaderAt(opts.From)
	w := writerCloser(opts.To)

	process(r, w, opts)

	closeStream(r)
	closeStream(w)
}

func closerReaderAt(from string) CloserReaderAt {
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

func process(r CloserReaderAt, w io.Writer, opts *Options) {
	bytes := make([]byte, opts.Limit)

	readBytesAt(r, bytes, opts.Offset)
	convertedBytes := convert(bytes, opts.Conv)
	writeBytes(w, convertedBytes)
}

func readBytesAt(r io.ReaderAt, bytes []byte, offset int) {
	_, err := r.ReadAt(bytes, int64(offset))
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot read bytes at %d\n", offset)
	}
}

func writeBytes(w io.Writer, bytes []byte) {
	_, err := w.Write(bytes)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func closeStream(c io.Closer) {
	if err := c.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func convert(bytes []byte, conv *string) []byte {
	str := string(bytes)
	readConvTypes := strings.Split(*conv, ",")

	for _, v := range readConvTypes {
		applyConv(str, v)
	}

	return []byte(str)
}

func applyConv(str, conv string) string {
	switch conv {
	case "upper_case":
		str = strings.ToUpper(str)
	case "lower_case":
		str = strings.ToLower(str)
	case "trim_spaces":
		str = strings.TrimFunc(str, unicode.IsSpace)
	}
	return str
}
