package dd

import (
	"errors"
	"flag"
	"lecture03_homework/dd/lib"
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
	if isStdin(path) {
		return nil
	}
	return validateFile(path)
}

func validateOutputFile(path string) error {
	if isStdout(path) {
		return nil
	}
	if validateFile(path) == nil {
		return errors.New("the file already exists")
	}
	return nil
}

func validateOffset(path string, offset int64) error {
	if offset < 0 {
		return errors.New("negative offset")
	}
	if isStdin(path) {
		return nil
	}

	if fileSize(path) < offset {
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
	var convErrors []error
	readConvTypes := strings.Split(*conv, ",")

	convErrors = appendIfErr(convErrors, validateConvExistence(readConvTypes)...)
	convErrors = appendIfErr(convErrors, validateNonContradictory(readConvTypes))

	return convErrors
}

// validateConvExistence проверяет считанные значения флага conv
// на существование (на корректность ввода)
func validateConvExistence(readConvTypes []string) []error {
	var typeErrors []error
	var res error

	for _, v := range readConvTypes {
		res = validateConvType(v)
		typeErrors = appendIfErr(typeErrors, res)
	}

	return typeErrors
}

func validateConvType(readConvType string) error {
	switch readConvType {
	case ChangeNothing, UpperCase, LowerCase, TrimSpaces:
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
		if readConvTypes[0] == TrimSpaces || readConvTypes[1] == TrimSpaces {
			return nil
		}
	}
	return errors.New("invalid set of conv types")
}

func appendIfErr(errors []error, possibleErrors ...error) []error {
	for _, e := range possibleErrors {
		if e != nil {
			errors = append(errors, e)
		}
	}
	return errors
}

// adjustFlags корректирует флаги для их более
// удобного использования
func adjustFlags(opts *Options, invalidFlags []error) {
	if invalidFlags != nil {
		return
	}

	if isNotStdin(opts.From) {
		configureLimit(opts)
	}
}

// configureLimit уточняет значение limit
func configureLimit(opts *Options) {
	if opts.Limit == NoLimit {
		opts.Limit = fileSize(opts.From)
	} else {
		opts.Limit = lib.MinInt64(opts.Limit, fileSize(opts.From))
	}
}

func isNotStdin(from string) bool {
	return from != stdin
}

func fileSize(path string) int64 {
	fileInfo, _ := os.Stat(path)
	return fileInfo.Size()
}
