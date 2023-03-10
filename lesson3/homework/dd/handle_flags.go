package dd

import (
	"errors"
	"flag"
	"os"
	"strings"
)

func ParseFlags() (*Options, []error) {
	var opts Options

	DefineFlags(&opts)
	flag.Parse()
	invalidFlags := ValidateFlags(&opts)

	return &opts, invalidFlags
}

func DefineFlags(opts *Options) {
	flag.StringVar(&opts.From, "from", Stdin, "Sets file to read. By default - stdin")

	flag.StringVar(&opts.To, "to", Stdout, "Sets file to write. By default - stdout")

	flag.IntVar(&opts.Offset, "offset", DefaultOffset, "Sets the number of bytes to skip from the beginning of the "+
		"input file for writing. By default - 0")

	flag.IntVar(&opts.Limit, "limit", DefaultLimit, "Sets the maximum number of bytes to read. "+
		"By default all content is copied starting with -offset")

	flag.IntVar(&opts.BlockSize, "block-size", DefaultBlockSize, "Sets the size of one block in bytes "+
		"for reading and writing. By default - 1")

	opts.Conv = flag.String("conv", DefaultConvType, "Sets text conversion options. "+
		"By default original text is copied without changes")
}

func ValidateFlags(opts *Options) (report []error) {
	report = appendIfErr(report, validateInputFile(opts.From))
	report = appendIfErr(report, validateOutputFile(opts.To))
	report = appendIfErr(report, validateOffset(opts.From, opts.Offset))
	report = appendIfErr(report, validateLimit(opts.Limit))
	report = appendIfErr(report, validateBlockSize(opts.BlockSize))
	report = appendIfErr(report, validateConv(opts.Conv)...)
	return
}

func validateFile(path string) error {
	_, err := os.Stat(path)
	return err
}

func validateInputFile(path string) error {
	if path == Stdin {
		return nil
	}
	return validateFile(path)
}

func validateOutputFile(path string) error {
	if path == Stdout {
		return nil
	}
	return validateFile(path)
}

func validateOffset(path string, offset int) error {
	if path == Stdin {
		return nil
	}

	file, err := os.Stat(path)

	if offset < 0 {
		return errors.New("negative offset")
	} else if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	} else if file.Size() < int64(offset) {
		return errors.New("offset is greater than input file size")
	}

	return nil
}

func validateLimit(limit int) error {
	if limit < 0 && limit != NoLimit {
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

func validateConv(conv *string) []error {
	readConvTypes := strings.Split(*conv, ",")
	var typeErrors []error
	var res error

	for _, v := range readConvTypes {
		res = validateConvType(v)
		appendIfErr(typeErrors, res)
	}

	return typeErrors
}

func validateConvType(cType string) error {
	convTypes := convTypes()
	for _, v := range convTypes {
		if cType == v {
			return nil
		}
	}
	return errors.New(cType + ": unexpected conv type")
}

func appendIfErr(errors []error, possibleErrors ...error) []error {
	for _, e := range possibleErrors {
		if e != nil {
			errors = append(errors, e)
		}
	}
	return errors
}
