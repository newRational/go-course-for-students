package dd

import (
	"errors"
	"io"
	"math"
	"strings"
	"unicode"
)

type ConvApplier struct {
	conversions map[string]Converter
}

func newConvApplier() *ConvApplier {
	return &ConvApplier{
		conversions: map[string]Converter{},
	}
}

func newConvApplierConversions(conversions *string) *ConvApplier {
	readConvTypes := strings.Split(*conversions, ",")
	ca := newConvApplier()

	for _, v := range readConvTypes {
		switch v {
		case upperCase:
			ca.add(upperCase, &UpperCase{})
		case lowerCase:
			ca.add(lowerCase, &LowerCase{})
		case trimSpaces:
			ca.add("trim_left", &TrimLeft{
				mustApply: true,
			})
			ca.add("trim_right", &TrimRight{})
		}
	}

	return ca
}

func (c *ConvApplier) startAll(r io.Reader, block []byte) ([]byte, int, error) {
	var errs error
	var err error
	var bytesCount int
	var n int

	for _, v := range c.conversions {
		if v.convName() == "trim_left" {
			continue
		}
		//fmt.Fprintln(os.Stderr, v.convName())
		block, n, err = v.convert(r, block)
		bytesCount += n
		errs = errors.Join(err, errs)
	}

	return block, bytesCount, errs
}

func (c *ConvApplier) startSingle(r io.Reader, name string, block []byte) ([]byte, int, error) {
	if conv := c.conversions[name]; conv != nil {
		return conv.convert(r, block)
	}
	return nil, 0, nil
}

func (c *ConvApplier) add(name string, conv Converter) {
	c.conversions[name] = conv
}

type Converter interface {
	convert(r io.Reader, block []byte) ([]byte, int, error)
	convName() string
}

type UpperCase struct {
}

type LowerCase struct {
}

type TrimLeft struct {
	mustApply bool
}

type TrimRight struct {
	trimmedBytes []byte
}

func (uc *UpperCase) convert(_ io.Reader, block []byte) ([]byte, int, error) {
	str := string(block)
	str = strings.ToUpper(str)
	return []byte(str), 0, nil
}

func (uc *UpperCase) convName() string {
	return upperCase
}

func (lc *LowerCase) convert(_ io.Reader, block []byte) ([]byte, int, error) {
	str := string(block)
	str = strings.ToLower(str)
	return []byte(str), 0, nil
}

func (lc *LowerCase) convName() string {
	return lowerCase
}

func (tl *TrimLeft) convert(r io.Reader, block []byte) ([]byte, int, error) {
	//fmt.Println("\tcalled\t")
	if !tl.mustApply {
		return nil, 0, nil
	}
	//fmt.Println("\tcalled must\t")
	tl.mustApply = false

	b := make([]byte, 1)
	if _, err := r.Read(b); err != nil && err != io.EOF {
		return nil, 0, err
	}

	trimmedLeftBytesCount := 0
	for unicode.IsSpace(rune(b[0])) {
		if _, err := r.Read(b); err != nil {
			return nil, 0, err
		}
		trimmedLeftBytesCount++
	}

	//fmt.Println("returned:", string(b))

	return b, trimmedLeftBytesCount, nil
}

func (tl *TrimLeft) convName() string {
	return "trim_left"
}

func (tr *TrimRight) convert(_ io.Reader, block []byte) ([]byte, int, error) {
	if len(block) == 0 {
		return block, 0, nil
	}

	var bytesToWrite []byte

	leftBytes, rightSpaceBytes := tr.splitBlock(block)
	rightSpaceBytesCount := len(rightSpaceBytes)

	if len(leftBytes) == 0 {
		tr.trimmedBytes = append(tr.trimmedBytes, rightSpaceBytes...)
	} else {
		bytesToWrite = append(tr.trimmedBytes, leftBytes...)
		tr.trimmedBytes = make([]byte, len(rightSpaceBytes))
		copy(tr.trimmedBytes, rightSpaceBytes)
	}

	return bytesToWrite, rightSpaceBytesCount, nil
}

func (tr *TrimRight) convName() string {
	return "trim_right"
}

func (tr *TrimRight) splitBlock(block []byte) ([]byte, []byte) {
	leftBytes := []byte(strings.TrimRightFunc(string(block), unicode.IsSpace))
	rightSpaceBytes := block[len(leftBytes):]

	return leftBytes, rightSpaceBytes
}

type Options struct {
	From      string
	To        string
	Offset    int64
	Limit     int64
	BlockSize int64
	Conv      *string
}

// Default Options values
const (
	stdin            = "stdin"
	stdout           = "stdout"
	defaultOffset    = 0
	defaultLimit     = NoLimit
	defaultBlockSize = 1
	defaultConvType  = changeNothing
)

const NoLimit = math.MaxInt64

// ConvTypes
const (
	changeNothing = ""
	upperCase     = "upper_case"
	lowerCase     = "lower_case"
	trimSpaces    = "trim_spaces"
)
