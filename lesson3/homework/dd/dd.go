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

		readBytes = append(startByte, readBytes...)

		readBytes, opts.Limit, err = correctBlock(r, readBytes, opts.Limit)
		if err != nil {
			return err
		}

		if readBytesCount == 0 || opts.Limit < 0 {
			return nil
		}

		if opts.Limit < opts.BlockSize {
			readBytes = readBytes[:opts.Limit]
		}

		readBytes, count, err := ca.startAll(r, readBytes)
		if err != nil {
			return err
		}

		_, err = w.Write(readBytes)
		if err != nil {
			return err
		}

		opts.Limit -= int64(readBytesCount + trimmedCount + count)

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

	return block, nil
}

func correctBlock(r io.Reader, block []byte, remainingBytesCount int64) ([]byte, int64, error) {
	if len(block) == 0 {
		return block, 0, nil
	}
	_, startByteIndex, count := findStartByteFromBack(block)

	tmp := make([]byte, count)
	copy(tmp, block[startByteIndex:])
	b := make([]byte, 1)
	var err error
	diff := 0
	for !utf8.Valid(tmp) && err == nil {
		_, err = r.Read(b)
		tmp = append(tmp, b...)
		diff++
	}

	if remainingBytesCount < int64(count+diff) {
		return block, remainingBytesCount, nil
	}

	if err != nil && err != io.EOF {
		return nil, remainingBytesCount - int64(diff), err
	}
	return append(block, tmp[count:]...), remainingBytesCount - int64(diff), nil
}

func findStartByteFromBack(block []byte) (byte, int, int) {
	bytesFromBackCount := 1
	l := len(block)
	for !utf8.RuneStart(block[l-bytesFromBackCount]) {
		bytesFromBackCount++
	}
	return block[l-bytesFromBackCount], l - bytesFromBackCount, bytesFromBackCount
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
