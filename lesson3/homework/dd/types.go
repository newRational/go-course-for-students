package dd

import "io"

type CloserReaderAt interface {
	io.ReaderAt
	io.Closer
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
	defaultBlockSize = 4
	defaultConvType  = ""
)

const NoLimit = -1

func convTypes() []string {
	return []string{"", "upper_case", "lower_case", "trim_spaces"}
}
