package dd

import (
	"errors"
	"flag"
	"os"
	"strings"
)

func ParseFlags() (*Options, []error) {
	opts := &Options{}

	DefineFlags(opts)
	flag.Parse()
	invalidFlags := ValidateFlags(opts)
	adjustFlags(opts, invalidFlags)

	return opts, invalidFlags
}

func DefineFlags(opts *Options) {
	flag.StringVar(&opts.From, "from", stdin, "Sets file to read. By default - stdin")

	flag.StringVar(&opts.To, "to", stdout, "Sets file to write. By default - stdout")

	flag.Int64Var(&opts.Offset, "offset", defaultOffset, "Sets the number of bytes to skip from the beginning of the "+
		"input file for writing. By default - 0")

	flag.Int64Var(&opts.Limit, "limit", defaultLimit, "Sets the maximum number of bytes to read. "+
		"By default all content is copied starting with -offset")

	flag.Int64Var(&opts.BlockSize, "block-size", defaultBlockSize, "Sets the size of one block in bytes "+
		"for reading and writing. By default - 1")

	opts.Conv = flag.String("conv", defaultConvType, "Sets text conversion options. "+
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
	if path == stdin {
		return nil
	}
	return validateFile(path)
}

func validateOutputFile(path string) error {
	if path == stdout {
		return nil
	}
	return validateFile(path)
}

func validateOffset(path string, offset int64) error {
	if path == stdin {
		return nil
	}

	file, err := os.Stat(path)

	if offset < 0 {
		return errors.New("negative offset")
	} else if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	} else if file.Size() < offset {
		return errors.New("offset is greater than input file size")
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

func adjustFlags(opts *Options, invalidFlags []error) {
	if invalidFlags != nil {
		return
	}

	if isNotStdin(opts.From) {
		configureLimit(opts)
	}
}

func isStdin(from string) bool {
	return from == stdin
}
func isNotStdin(from string) bool {
	return from != stdin
}

func configureLimit(opts *Options) {
	if opts.Limit == NoLimit {
		opts.Limit = fileSize(opts.From)
	}
}

func fileSize(path string) int64 {
	file, _ := os.Open(path)
	fileInfo, _ := file.Stat()
	return fileInfo.Size()
}
