package dd

import "io"

type CloserReaderAt interface {
	io.ReaderAt
	io.Closer
}

type Options struct {
	From      string
	To        string
	Offset    int
	Limit     int
	BlockSize int
	Conv      *string
}

// Default Options values
const (
	Stdin            = "stdin"
	Stdout           = "stdout"
	DefaultOffset    = 0
	DefaultLimit     = NoLimit
	DefaultBlockSize = 4
	DefaultConvType  = ""
)

const NoLimit = -1

func convTypes() []string {
	return []string{"", "upper_case", "lower_case", "trim_spaces"}
}
