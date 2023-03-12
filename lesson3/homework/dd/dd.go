package dd

import (
	"errors"
	"io"
	"os"
	"strings"
	"unicode"
)

func Start(opts *Options) (err error) {
	r, err := readCloser(opts.From)
	if err != nil {
		return err
	}

	w, err := writeCloser(opts.To)
	if err != nil {
		return err
	}

	defer func() {
		err = errors.Join(r.Close(), w.Close())
	}()

	return process(r, w, opts)
}

func readCloser(from string) (io.ReadCloser, error) {
	if from == stdin {
		return os.Stdin, nil
	}

	file, err := os.Open(from)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func writeCloser(to string) (io.WriteCloser, error) {
	if to == stdout {
		return os.Stdout, nil
	}

	file, err := os.Create(to)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func process(r io.Reader, w io.Writer, opts *Options) error {
	if opts.Limit != NoLimit {
		r = io.LimitReader(r, opts.Limit+opts.Offset)
	}
	bytes, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	convertedBytes := convert(bytes[opts.Offset:], opts.Conv)

	if _, err = w.Write(convertedBytes); err != nil {
		return err
	}

	return nil
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
	case UpperCase:
		str = strings.ToUpper(str)
	case LowerCase:
		str = strings.ToLower(str)
	case TrimSpaces:
		str = strings.TrimFunc(str, unicode.IsSpace)
	}
	return str
}
