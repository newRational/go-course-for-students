package dd

import (
	"errors"
	"flag"
	"lecture03_homework/dd/lib"
	"os"
	"strings"
)

func ParseFlags() (*Options, error) {
	opts := &Options{}

	DefineFlags(opts)
	flag.Parse()
	invalidFlags := ValidateFlags(opts)
	adjustFlags(opts, invalidFlags)

	return opts, invalidFlags
}

func DefineFlags(opts *Options) {
	flag.StringVar(&opts.From, "from", stdin, "Sets file to read. By default: stdin")

	flag.StringVar(&opts.To, "to", stdout, "Sets file to write. By default: stdout")

	flag.Int64Var(&opts.Offset, "offset", defaultOffset, "Sets the number of bytes to skip from the beginning of the "+
		"input file for writing. By default: 0")

	flag.Int64Var(&opts.Limit, "limit", defaultLimit, "Sets the maximum number of bytes to read. "+
		"By default all content is copied starting with -offset")

	flag.Int64Var(&opts.BlockSize, "block-size", defaultBlockSize, "Sets the size of one block in bytes "+
		"for reading and writing. By default: 1")

	opts.Conv = flag.String("conv", defaultConvType, "Sets text conversion options. "+
		"By default original text is copied without changes")
}

func ValidateFlags(opts *Options) error {
	return errors.Join(
		validateInput(opts.From),
		validateOutput(opts.To),
		validateOffset(opts.From, opts.Offset),
		validateLimit(opts.Limit),
		validateBlockSize(opts.BlockSize),
		validateConv(opts.Conv),
	)
}

func validateFile(path string) error {
	_, err := os.Stat(path)
	return err
}

func validateInput(from string) error {
	if from == stdin {
		return nil
	}
	return validateFile(from)
}

func validateOutput(to string) error {
	if to == stdout {
		return nil
	}
	if validateFile(to) == nil {
		return os.ErrExist
	}
	return nil
}

func validateOffset(from string, offset int64) error {
	if offset < 0 {
		return errors.New("negative offset")
	}
	if from == stdin {
		return nil
	}
	if validateInput(from) == nil && fileSize(from) < offset {
		return errors.New("offset is greater than input size")
	}

	return nil
}

func validateLimit(limit int64) error {
	if limit < 0 && limit != NoLimit {
		return errors.New("negative limit")
	}

	return nil
}

func validateBlockSize(blockSize int64) error {
	if blockSize < 0 {
		return errors.New("negative block-size")
	}

	return nil
}

func validateConv(conv *string) error {
	readConvTypes := strings.Split(*conv, ",")

	err := errors.Join(
		validateConvExistence(readConvTypes),
		validateNonContradictory(readConvTypes),
	)

	return err
}

// validateConvExistence проверяет считанные значения флага conv
// на существование (на корректность ввода)
func validateConvExistence(readConvTypes []string) error {
	var errs []error
	for _, v := range readConvTypes {
		errs = append(errs, validateConvType(v))
	}
	return errors.Join(errs...)
}

func validateConvType(readConvType string) error {
	switch readConvType {
	case changeNothing, upperCase, lowerCase, trimSpaces:
		return nil
	}
	return errors.New(readConvType + ": unexpected conv type")
}

// validateNonContradictory определяет валидность
// набора, состоящего из типов конвертации, в соответствии
// с его содержимым
func validateNonContradictory(readConvTypes []string) error {
	switch len(readConvTypes) {
	case 1, 0:
		return nil
	case 2:
		if readConvTypes[0] == trimSpaces || readConvTypes[1] == trimSpaces {
			return nil
		}
	}
	return errors.New("invalid set of conv types")
}

// adjustFlags корректирует флаги для их более
// удобного использования
func adjustFlags(opts *Options, invalidFlags error) {
	if invalidFlags != nil {
		return
	}

	if opts.From != stdin && validateInput(opts.From) == nil {
		configureLimit(opts)
	}
}

// configureLimit уточняет значение limit
func configureLimit(opts *Options) {
	if opts.Limit == NoLimit {
		opts.Limit = fileSize(opts.From) - opts.Offset
	} else {
		opts.Limit = lib.MinInt64(opts.Limit, fileSize(opts.From))
	}
}

func fileSize(path string) int64 {
	fileInfo, _ := os.Stat(path)
	return fileInfo.Size()
}
