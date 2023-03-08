package dd

// ConvTypes
const (
	Default = iota
	UpperCase
	LowerCase
	TrimSpaces
)

type Options struct {
	From      string
	To        string
	Offset    int
	Limit     int
	BlockSize int
	ConvType  int
}
