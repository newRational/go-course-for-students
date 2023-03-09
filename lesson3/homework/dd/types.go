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
	DefaultFrom      = "stdin"
	DefaultTo        = "stdout"
	DefaultOffset    = 0
	DefaultLimit     = -1
	DefaultBlockSize = 2
	DefaultConvType  = ChangeNothing
)

// ConvTypes
const (
	ChangeNothing = ""
	UpperCase     = "upper_case"
	LowerCase     = "lower_case"
	TrimSpaces    = "trim_spaces"
)
