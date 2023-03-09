package dd

import (
	"errors"
	"flag"
	"os"
)

func ParseFlags() (*Options, []error) {
	var opts Options

	DefineFlags(&opts)
	flag.Parse()
	invalidFlags := ValidateFlags(&opts)

	return &opts, invalidFlags
}

func DefineFlags(opts *Options) {
	flag.StringVar(&opts.From, "from", DefaultFrom, "Sets file to read. By default - stdin")

	flag.StringVar(&opts.To, "to", DefaultTo, "Sets file to write. By default - stdout")

	flag.IntVar(&opts.Offset, "offset", DefaultOffset, "Sets the number of bytes to skip from the beginning of the "+
		"input file for writing. By default - 0")

	flag.IntVar(&opts.Limit, "limit", DefaultLimit, "Sets the maximum number of bytes to read. "+
		"By default all content is copied starting with -offset")

	flag.IntVar(&opts.BlockSize, "block-size", DefaultBlockSize, "Sets the size of one block in bytes "+
		"for reading and writing. By default - 1")

	flag.StringVar(&opts.Conv, "conv", DefaultConvType, "Sets text conversion options. "+
		"By default original text is copied without changes")
}

func ValidateFlags(opts *Options) (report []error) {
	report = appendIfErr(report, validateInputFile(opts.From))
	report = appendIfErr(report, validateOutputFile(opts.To))
	report = appendIfErr(report, validateLimit(opts.Limit))
	report = appendIfErr(report, validateOffset(opts.From, opts.Offset))
	report = appendIfErr(report, validateBlockSize(opts.BlockSize))
	report = appendIfErr(report, validateConv(opts.Conv))
	return
}

func validateFile(path string) error {
	_, err := os.Stat(path)
	return err
}

func validateInputFile(path string) error {
	if path == "stdin" {
		return nil
	}
	return validateFile(path)
}

func validateOutputFile(path string) error {
	if path == "stdout" {
		return nil
	}
	return validateFile(path)
}

func validateOffset(path string, offset int) error {
	file, err := os.Stat(path)

	if err != nil {
		return err
	} else if file.Size() < int64(offset) {
		return errors.New("offset is greater than input file size")
	} else if offset < 0 {
		return errors.New("negative offset")
	}

	return nil
}

func validateLimit(limit int) error {
	if limit < 0 {
		return errors.New("negative limit")
	}

	return nil
}

func validateBlockSize(blockSize int) error {
	if blockSize < 0 {
		return errors.New("negative block-size")
	}

	return nil
}

func validateConv(conv string) error {
	switch conv {
	case ChangeNothing:
		return nil
	case UpperCase:
		return nil
	case LowerCase:
		return nil
	case TrimSpaces:
		return nil
	default:
		return errors.New("unknown conv value: " + conv)
	}
}

func appendIfErr(errors []error, err error) []error {
	if err != nil {
		return append(errors, err)
	}

	return errors
}
