package dd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

func Start(opts *Options) {
	r := readAtCloser(opts.From)
	w := writeCloser(opts.To)

	process(r, w, opts)

	closeStream(r)
	closeStream(w)
}

func readAtCloser(from string) ReadAtCloser {
	if isStdin(from) {
		return os.Stdin
	}
	file, err := os.Open(from)
	reportIfErr(err)
	return file
}

func writeCloser(to string) io.WriteCloser {
	if isStdout(to) {
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

// processFromStdin сначала создает новый readAtCloser на
// основе временного файла, затем обрабатывает данные,
// после чего временный файл удаляется
func processFromStdin(w io.Writer, opts *Options) {
	tmpFilePath := ".tmp"
	r := readAtCloserOnNewFile(tmpFilePath, opts)

	copyAndConvert(r, w, opts)

	closeStream(r)
	removeFile(tmpFilePath)
}

func processFromFile(r io.ReaderAt, w io.Writer, opts *Options) {
	copyAndConvert(r, w, opts)
}

// Данные из readerAt копируются в слайс, затем конвертируются,
// после чего выводятся
func copyAndConvert(r io.ReaderAt, w io.Writer, opts *Options) {
	bytes := make([]byte, opts.Limit)

	readBytesAt(r, bytes, opts.Offset)
	convertedBytes := convert(bytes, opts.Conv)
	writeBytes(w, convertedBytes)
}

// convert применяет каждое указанное правило
// конвертации к считанному слайсу байтов
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
	case UpperCase:
		str = strings.ToUpper(str)
	case LowerCase:
		str = strings.ToLower(str)
	case TrimSpaces:
		str = strings.Trim(str, "\x00")
		str = strings.TrimFunc(str, unicode.IsSpace)
	}
	return str
}

// readAtCloserOnNewFile сначала создается новый файл и заполняется
// содежимым из stdin, затем проверятся offset на валидность, после чего
// изменяются значения полей экземпляра Options и возвращается новый
// readAtCloser
func readAtCloserOnNewFile(newFilePath string, opts *Options) ReadAtCloser {
	createAndFillFile(newFilePath, os.Stdin, opts)
	reportIfErr(validateOffset(newFilePath, opts.Offset))
	opts.From = newFilePath
	configureLimit(opts)
	return readAtCloser(newFilePath)
}

func createAndFillFile(outputPath string, r io.Reader, opts *Options) {
	w := writeCloser(outputPath)
	copyAccordingToLimit(w, r, opts)
	closeStream(w)
}

func copyAccordingToLimit(w io.Writer, r io.Reader, opts *Options) {
	var err error
	if opts.Limit == NoLimit {
		_, err = io.Copy(w, r)
	} else {
		_, err = io.CopyN(w, r, opts.Limit+opts.Offset)
	}
	reportIfErr(err, io.EOF)
}

func readBytesAt(r io.ReaderAt, bytes []byte, offset int64) {
	_, err := r.ReadAt(bytes, offset)
	reportIfErr(err, io.EOF)
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

func isStdout(to string) bool {
	return to == stdout
}

func isStdin(from string) bool {
	return from == stdin
}

func reportIfErr(err error, except ...error) {
	if isExceptOrNil(err, except...) {
		return
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func isExceptOrNil(err error, except ...error) bool {
	for _, e := range except {
		if err == e {
			return true
		}
	}
	return err == nil
}
