package dd

import (
	"errors"
	"fmt"
	"io"
	"os"
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
	ca := newConvApplierConversions(opts.Conv)
	var readBytes []byte

	err := skipOffset(r, opts, block)
	if err != nil {
		return err
	}

	for opts.Limit > 0 {
		startByte, trimmedCount, err := ca.startSingle(r, "trim_left", block)
		if err != nil {
			return err
		}
		readBytes, err = readBlock(r, block)
		readBytesCount := len(readBytes)
		if err != nil {
			return err
		}

		//fmt.Println("startByte:", startByte, ":", string(startByte))
		//fmt.Println("trimmedCount:", trimmedCount)
		//fmt.Println("readBytes:", readBytes, ":", string(readBytes))

		readBytes = append(startByte, readBytes...)
		//fmt.Println("readBytes 2:", readBytes, ":", string(readBytes))

		readBytes, _ = correctBlock(r, readBytes)
		//fmt.Println("readBytes 3:", readBytes, ":", string(readBytes))

		if readBytesCount == 0 {
			return nil
		}

		if opts.Limit < opts.BlockSize {
			readBytes = readBytes[:opts.Limit]
		}

		_, count, _ := ca.startAll(r, readBytes)
		//fmt.Println("readBytes:", readBytes, ":", string(readBytes))

		_, err = w.Write(readBytes)
		if err != nil {
			return err
		}

		opts.Limit -= int64(readBytesCount + trimmedCount + count)
		//fmt.Println("\tchecked\n")
	}

	return nil
}

func readBlock(r io.Reader, block []byte) ([]byte, error) {
	readBytesCount, err := r.Read(block)

	if err != nil && err != io.EOF {
		return nil, err
	}

	if readBytesCount < len(block) {
		return block[:readBytesCount], nil
	}

	return correctBlock(r, block)
}

func correctBlock(r io.Reader, block []byte) ([]byte, error) {
	if len(block) == 0 {
		return block, nil
	}
	startByte, count := findStartByteFromBack(block)
	diff := runeLen(startByte) - count

	tmp := make([]byte, diff)
	_, err := r.Read(tmp)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return append(block, tmp...), nil
}

func findStartByteFromBack(block []byte) (byte, int) {
	bytesFromBackCount := 1
	l := len(block)
	for !utf8.RuneStart(block[l-bytesFromBackCount]) {
		bytesFromBackCount++
	}
	return block[l-bytesFromBackCount], bytesFromBackCount
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

func skipOffset(r io.Reader, opts *Options, block []byte) error {
	remainingBytesCount := opts.Offset

	for remainingBytesCount >= opts.BlockSize {
		readBytesCount, err := r.Read(block)
		if err != nil && err != io.EOF {
			return err
		}
		if readBytesCount == 0 {
			fmt.Fprintln(os.Stderr, "offset is greater than input size")
			return errors.New("offset is greater than input size")
		}
		remainingBytesCount -= int64(readBytesCount)
	}

	_, err := r.Read(block[:remainingBytesCount])
	if err != nil && err != io.EOF {
		return err
	}

	return nil
}
