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

	if opts.Limit == NoLimit {
		err = processNoLimit(r, w, opts)
	} else {
		err = processLimit(r, w, opts)
	}

	return err
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

func processLimit(r io.Reader, w io.Writer, opts *Options) error {
	block := make([]byte, opts.BlockSize)
	remainingBytesCount := opts.Limit

	err := skipOffset(r, opts, block)
	if err != nil {
		return err
	}

	for remainingBytesCount > 0 {
		readBytesCount, _ := r.Read(block)
		if readBytesCount == 0 {
			break
		}

		if remainingBytesCount < opts.BlockSize {
			block = block[:remainingBytesCount]
		}
		if int64(readBytesCount) < opts.BlockSize {
			block = block[:readBytesCount]
		}
		convertedBytes := convert(block, opts.Conv)

		if _, err = w.Write(convertedBytes); err != nil {
			return err
		}
		remainingBytesCount -= int64(readBytesCount)
	}

	return nil
}

func processNoLimit(r io.Reader, w io.Writer, opts *Options) error {
	block := make([]byte, opts.BlockSize)

	err := skipOffset(r, opts, block)
	if err != nil {
		return err
	}

	for {
		readBytesCount, _ := r.Read(block)
		if readBytesCount == 0 {
			break
		}

		if int64(readBytesCount) < opts.BlockSize {
			block = block[:readBytesCount]
		}

		convertedBytes := convert(block, opts.Conv)

		if _, err = w.Write(convertedBytes); err != nil {
			return err
		}
	}

	return nil
}

func skipOffset(r io.Reader, opts *Options, block []byte) error {
	remainingBytesCount := opts.Offset

	for remainingBytesCount > opts.BlockSize {
		readBytesCount, _ := r.Read(block)
		if readBytesCount == 0 {
			return errors.New("offset is greater than input size")
		}
		remainingBytesCount -= int64(readBytesCount)
	}

	_, _ = r.Read(block[:remainingBytesCount])
	return nil
}

/* ---------------- Считывание по блока ------------------
1. 	en, block-size = 4, offset = 0, limit = no, conv = no				+
2.	en, block-size = 4, offset = 0, limit = no, conv = upper_case		+
3.	en, block-size = 4, offset = 0, limit = no, conv = lower_case		+

4. 	en, block-size = 4, offset = 0, limit = set, conv = no				+
5. 	en, block-size = 4, offset = 0, limit = set, conv = upper_case		+
6. 	en, block-size = 4, offset = 0, limit = set, conv = lower_case		+

7. 	en, block-size = 4, offset = set, limit = no, conv = no				+
8. 	en, block-size = 4, offset = set, limit = no, conv = upper_case		+
9. 	en, block-size = 4, offset = set, limit = no, conv = lower_case		+

10.	en, block-size = 4, offset = set, limit = set, conv = no			+
11.	en, block-size = 4, offset = set, limit = set, conv = upper_case	+
12.	en, block-size = 4, offset = set, limit = set, conv = lower_case	+
---------------------------------------------------------------- */

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
