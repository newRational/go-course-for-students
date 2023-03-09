package dd

type Options struct {
	From      string
	To        string
	Offset    int
	Limit     int
	BlockSize int
	Conv      string
}

// Default Options values
const (
	Stdin            = "stdin"
	Stdout           = "stdout"
	DefaultOffset    = 0
	DefaultLimit     = NoLimit
	DefaultBlockSize = 100
	DefaultConvType  = ChangeNothing
)

const NoLimit = -1

// ConvTypes
const (
	ChangeNothing = ""
	UpperCase     = "upper_case"
	LowerCase     = "lower_case"
	TrimSpaces    = "trim_spaces"
)
