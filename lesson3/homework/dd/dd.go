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

func Start(opts *Options) (err error) {
	r, err := readCloser(opts.From)
	if err != nil {
		return err
	}

	w, err := writeCloser(opts.To)
	if err != nil {
		return err
	}

	err = process(r, w, opts)

	defer func() {
		err = errors.Join(r.Close(), w.Close(), err)
	}()

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
		readBytes = correctBlock(r, b)
		trimFirst = true
	}

	for remainingBytesCount > 0 {
		if !trimFirst {
			readBytes = readBlock(r, block)
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

		if block != nil {
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

func readBlock(r io.Reader, block []byte) []byte {
	readBytesCount, _ := r.Read(block)

	if readBytesCount < len(block) {
		return block[:readBytesCount]
	}

	return correctBlock(r, block)
}

func correctBlock(r io.Reader, block []byte) []byte {
	runeStart, count := findStartByteFromBack(block)
	diff := runeLen(runeStart) - count

	tmp := make([]byte, diff)
	_, _ = r.Read(tmp)
	return append(block, tmp...)
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

	for remainingBytesCount >= opts.BlockSize {
		readBytesCount, _ := r.Read(block)
		if readBytesCount == 0 {
			fmt.Fprintln(os.Stderr, "offset is greater than input size")
			return errors.New("offset is greater than input size")
		}
		remainingBytesCount -= int64(readBytesCount)
	}

	_, _ = r.Read(block[:remainingBytesCount])

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
