package dd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

func Start(opts *Options) error {
	r, err := readCloser(opts.From)
	if err != nil {
		return err
	}

	w, err := writeCloser(opts.To)
	if err != nil {
		return err
	}

	defer func(r io.ReadCloser, w io.WriteCloser) error {
		err = r.Close()
		if err != nil {
			return err
		}
		err = w.Close()
		if err != nil {
			return err
		}
		return nil
	}(r, w)

	err = process(r, w, opts)

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

func process(r io.Reader, w io.Writer, opts *Options) error {
	block := make([]byte, opts.BlockSize)
	var trimmedRightBytes []byte
	var trimFirst bool
	var readBytes []byte

	doTrim := containsString(opts.Conv, TrimSpaces)

	err := skipOffset(r, opts, block)
	if err != nil {
		return err
	}

	remainingBytesCount := opts.Limit

	if doTrim {
		b, trimmedCount, err := trimLeft(r)
		if err != nil {
			return err
		}
		remainingBytesCount -= int64(trimmedCount)
		if readBytes, err = correctBlock(r, b); err != nil {
			return err
		}
		trimFirst = true
	}

	for remainingBytesCount > 0 {
		if !trimFirst {
			if readBytes, err = readBlock(r, block); err != nil {
				return err
			}
		}
		readBytesCount := len(readBytes)

		if readBytesCount == 0 {
			return nil
		}

		if remainingBytesCount < opts.BlockSize {
			readBytes = readBytes[:remainingBytesCount]
		}

		if doTrim {
			readBytes, trimmedRightBytes = trimRight(readBytes, trimmedRightBytes)
		}

		if readBytes != nil {
			if _, err = w.Write(convert(readBytes, opts.Conv)); err != nil {
				return err
			}
		}

		remainingBytesCount -= int64(readBytesCount)
		trimFirst = false
	}

	return nil
}

func trimRight(block, previousTrimmedBytes []byte) (bytesToWrite, allTrimmedBytes []byte) {
	leftBytes, rightSpaceBytes := splitBlock(block)

	if len(leftBytes) == 0 {
		allTrimmedBytes = append(previousTrimmedBytes, rightSpaceBytes...)
	} else {
		bytesToWrite = append(previousTrimmedBytes, leftBytes...)
		allTrimmedBytes = make([]byte, len(rightSpaceBytes))
		copy(allTrimmedBytes, rightSpaceBytes)
	}

	return
}

func splitBlock(block []byte) ([]byte, []byte) {
	leftBytes := []byte(strings.TrimRightFunc(string(block), unicode.IsSpace))
	rightSpaceBytes := block[len(leftBytes):]

	return leftBytes, rightSpaceBytes
}

func readBlock(r io.Reader, block []byte) ([]byte, error) {
	readBytesCount, err := r.Read(block)
	if err != nil {
		return nil, err
	}

	if readBytesCount < len(block) {
		return block[:readBytesCount], nil
	}

	return correctBlock(r, block)
}

func correctBlock(r io.Reader, block []byte) ([]byte, error) {
	runeStart, count := findStartByteFromBack(block)
	diff := runeLen(runeStart) - count

	tmp := make([]byte, diff)
	_, err := r.Read(tmp)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return append(block, tmp...), nil
}

func findStartByteFromBack(block []byte) (byte, int) {
	i := 1
	l := len(block)
	for !utf8.RuneStart(block[l-i]) {
		i++
	}
	return block[l-i], i
}

func runeLen(b byte) int {
	if b&0b11110000 == 0b11110000 {
		return 4
	} else if b&0b11100000 == 0b11100000 {
		return 3
	} else if b&0b11000000 == 0b11000000 {
		return 2
	} else if b&0b10000000 == 0b10000000 {
		return 0
	} else {
		return 1
	}
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
	}
	return str
}

func containsString(bunch *string, target string) bool {
	content := strings.Split(*bunch, ",")
	for _, v := range content {
		if target == v {
			return true
		}
	}
	return false
}

func skipOffset(r io.Reader, opts *Options, block []byte) error {
	remainingBytesCount := opts.Offset

	for remainingBytesCount > 0 {
		readBytesCount, err := r.Read(block)
		if err != nil && err != io.EOF {
			return err
		}
		if int64(readBytesCount) < opts.BlockSize {
			fmt.Fprintln(os.Stderr, "offset is greater than input size")
			return errors.New("offset is greater than input size")
		}
		remainingBytesCount -= int64(readBytesCount)
	}

	_, err := r.Read(block[:remainingBytesCount])
	if err != nil {
		return err
	}

	return nil
}

func trimLeft(r io.Reader) ([]byte, int, error) {
	b := make([]byte, 1)
	_, err := r.Read(b)
	if err != nil {
		return nil, 0, err
	}
	n := 0
	for unicode.IsSpace(rune(b[0])) {
		_, err = r.Read(b)
		if err != nil {
			return nil, 0, err
		}
		n++
	}
	return b, n, nil
}
