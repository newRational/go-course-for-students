package dd

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
	defaultConvType  = ChangeNothing
)

const NoLimit = -1

// ConvTypes
const (
	ChangeNothing = ""
	UpperCase     = "upper_case"
	LowerCase     = "lower_case"
	TrimSpaces    = "trim_spaces"
)
