package dd

import (
	"errors"
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

	defer func() {
		err = errors.Join(r.Close(), w.Close())
	}()

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
	var trimmedLeftCount int
	var trimmedRightBytes []byte
	doTrim := containsString(opts.Conv, TrimSpaces)

	err := skipOffset(r, opts, block)
	if err != nil {
		return err
	}

	remainingBytesCount := opts.Limit

	if opts.Limit == NoLimit {
		remainingBytesCount = 1
	}

	if doTrim {
		if trimmedLeftCount, err = trimLeft(r, w, block); err != nil {
			return err
		}
	}

	if opts.Limit != NoLimit {
		remainingBytesCount -= int64(trimmedLeftCount)
	}

	for remainingBytesCount > 0 {
		readBytes, _ := readBlock(r, block)
		readBytesCount := len(readBytes)

		if readBytesCount == 0 {
			break
		}

		if opts.Limit != NoLimit && remainingBytesCount < opts.BlockSize {
			readBytes = readBytes[:remainingBytesCount]
		}

		if doTrim {
			readBytes, trimmedRightBytes = trimRight(readBytes, trimmedRightBytes)
		}

		convertedBytes := convert(readBytes, opts.Conv)

		if readBytes != nil {
			if _, err = w.Write(convertedBytes); err != nil {
				return err
			}
		}

		if opts.Limit != NoLimit {
			remainingBytesCount -= int64(readBytesCount)
		}
	}

	return nil
}

func skipOffset(r io.Reader, opts *Options, block []byte) error {
	remainingBytesCount := opts.Offset

	for remainingBytesCount >= opts.BlockSize {
		readBytesCount, _ := r.Read(block)
		if readBytesCount == 0 {
			return errors.New("offset is greater than input size")
		}
		remainingBytesCount -= int64(readBytesCount)
	}

	_, _ = r.Read(block[:remainingBytesCount])
	return nil
}

func trimLeft(r io.Reader, w io.Writer, block []byte) (int, error) {
	var str string
	var trimmedCount int

	for len(str) == 0 {
		readBytes, _ := readBlock(r, block)
		readBytesCount := len(readBytes)

		str = strings.TrimLeftFunc(string(readBytes), unicode.IsSpace)
		trimmedCount += readBytesCount - len(str)

		if readBytesCount == 0 {
			break
		}
	}

	if _, err := w.Write([]byte(str)); err != nil {
		return trimmedCount, err
	}

	return trimmedCount, nil
}

func trimRight(block, previousTrimmedBytes []byte) (bytesToWrite []byte, allTrimmedBytes []byte) {
	//fmt.Println("block:\t\t\t", "|"+string(block)+"|")
	//fmt.Println("previousTrimmedBytes:\t", "|"+string(previousTrimmedBytes)+"|")

	leftBytes, rightSpaceBytes := splitBlock(block)

	//fmt.Println("leftBytes:\t\t", "|"+string(leftBytes)+"|")
	//fmt.Println("rightBytes:\t\t", "|"+string(rightSpaceBytes)+"|")

	if len(leftBytes) == 0 {
		allTrimmedBytes = append(previousTrimmedBytes, rightSpaceBytes...)

		//fmt.Println("1) allTrimmedBytes:\t", "|"+string(allTrimmedBytes)+"|")
		//fmt.Println()
	} else {
		bytesToWrite = append(previousTrimmedBytes, leftBytes...)
		copy(allTrimmedBytes, rightSpaceBytes)

		//fmt.Println("2) bytesToWrite:\t", "|"+string(bytesToWrite)+"|")
		//fmt.Println("2) rightSpaceBytes:\t", "|"+string(rightSpaceBytes)+"|")
		//fmt.Println("2) allTrimmedBytes:\t", "|"+string(allTrimmedBytes)+"|")
		//fmt.Println()
	}
	return
}

func splitBlock(block []byte) ([]byte, []byte) {
	leftBytes := []byte(strings.TrimRightFunc(string(block), unicode.IsSpace))
	rightSpaceBytes := block[len(leftBytes):]

	return leftBytes, rightSpaceBytes
}

func readBlock(r io.Reader, block []byte) ([]byte, error) {
	readBytesCount, _ := r.Read(block)

	if readBytesCount < len(block) {
		return block[:readBytesCount], nil
	}

	runeStart, count := findStartByteFromBack(block)
	diff := runeLen(runeStart) - count

	tmp := make([]byte, diff)
	_, _ = r.Read(tmp)
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
