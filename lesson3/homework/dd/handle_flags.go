package dd

import (
	"flag"
	"lecture03_homework/lib/e"
	"os"
)

func DefineFlags(opts *Options) {
	flag.StringVar(&opts.From, "from", DefaultFrom, "Sets file to read. By default - stdin")

	flag.StringVar(&opts.To, "to", DefaultTo, "Sets file to write. By default - stdout")

	flag.IntVar(&opts.Offset, "offset", DefaultOffset, "Sets the number of bytes to skip from the beginning of the "+
		"input file for writing. By default - 0")

	flag.IntVar(&opts.Limit, "limit", DefaultLimit, "Sets the maximum number of bytes to read. "+
		"By default all content is copied starting with -offset")

	flag.IntVar(&opts.BlockSize, "block-size", DefaultBlockSize, "Sets the size of one block in bytes "+
		"for reading and writing. By default - 1")

	flag.IntVar(&opts.ConvType, "conv", DefaultConvType, "Sets text conversion options. "+
		"By default original text is copied without changes")
}

func ValidateFlags(opts *Options) (report []error) {
	report = appendIfNotNil(report, validateInputFile(opts.From))
	report = appendIfNotNil(report, validateOutputFile(opts.To))

	return
}

func validateInputFile(path string) error {
	_, err := os.Stat(path)

	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		return e.Wrap("input file doesn't exist", err)
	}

	return e.Wrap("invalid input file", err)
}

func validateOutputFile(path string) error {
	_, err := os.Stat(path)

	if err == nil {
		return e.Wrap("output file is already exist", err)
	}
	if os.IsNotExist(err) {
		return nil
	}

	return e.Wrap("invalid output file", err)
}

func appendIfNotNil(errors []error, err error) []error {
	if err != nil {
		return append(errors, err)
	}
	return errors
}
