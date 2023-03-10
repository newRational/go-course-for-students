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
	if from == stdin {
		return os.Stdin
	}
	file, err := os.Open(from)
	reportIfErr(err)
	return file
}

func writerCloser(to string) io.WriteCloser {
	if to == stdout {
		return os.Stdout
	}
	file, err := os.Create(to)
	reportIfErr(err)
	return file
}

func process(r io.ReaderAt, w io.Writer, opts *Options) {
	if isStdin(opts.From) {
		processFromStdin(w, opts)
	} else {
		processFromFile(r, w, opts)
	}
}

func processFromStdin(w io.Writer, opts *Options) {
	tmpFilePath := ".tmp"
	r := readerOnNewFile(tmpFilePath, opts)

	copyAndConvert(r, w, opts)

	closeStream(r)
	removeFile(tmpFilePath)
}

func readerOnNewFile(newFilePath string, opts *Options) CloserReaderAt {
	createAndFillFile(newFilePath, os.Stdin)
	opts.From = newFilePath
	configureLimit(opts)
	return closerReaderAt(newFilePath)
}

func createAndFillFile(outputPath string, r io.Reader) {
	w := writerCloser(outputPath)
	copyFile(w, r)
	closeStream(w)
}

func processFromFile(r io.ReaderAt, w io.Writer, opts *Options) {
	copyAndConvert(r, w, opts)
}

func copyAndConvert(r io.ReaderAt, w io.Writer, opts *Options) {
	bytes := make([]byte, opts.Limit)

	readBytesAt(r, bytes, opts.Offset)
	convertedBytes := convert(bytes, opts.Conv)
	writeBytes(w, convertedBytes)
}

func readBytesAt(r io.ReaderAt, bytes []byte, offset int64) {
	_, err := r.ReadAt(bytes, offset)
	reportIfErr(err, io.EOF)
}

func convert(bytes []byte, conv *string) []byte {
	str := string(bytes)
	readConvTypes := strings.Split(*conv, ",")

	for _, v := range readConvTypes {
		str = applyConv(str, v)
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

func writeBytes(w io.Writer, bytes []byte) {
	_, err := w.Write(bytes)
	reportIfErr(err)
}

func closeStream(c io.Closer) {
	err := c.Close()
	reportIfErr(err)
}

func removeFile(path string) {
	err := os.Remove(path)
	reportIfErr(err)
}

func copyFile(w io.Writer, r io.Reader) {
	_, err := io.Copy(w, r)
	reportIfErr(err)
}

func reportIfErr(err error, except ...error) {
	if isExceptOrNil(err, except...) {
		return
	}
	fmt.Fprintln(os.Stderr, err)
}

func isExceptOrNil(err error, except ...error) bool {
	for _, e := range except {
		if err == e {
			return true
		}
	}
	return err == nil
}
