package main

import (
	"flag"
	"fmt"
	"lecture03_homework/dd"
	"os"
)

//type Options struct {
//	From      string
//	To        string
//	Offset    int
//	Limit     int
//	BlockSize int
//	ConvType  int
//	// todo: add required flags
//}

func ParseFlags() (*dd.Options, error) {
	var opts dd.Options

	flag.StringVar(&opts.From, "from", "stdin", "File to read. By default - stdin")

	flag.StringVar(&opts.To, "to", "stdout", "File to write. By default - stdout")

	flag.IntVar(&opts.Offset, "offset", 0, "Sets the number of bytes to skip from the beginning of the "+
		"source file for writing. By default value - 0")

	flag.IntVar(&opts.Limit, "limit", -1, "Sets the maximum number of bytes to read. "+
		"Copy all content starting with -offset by default")

	flag.IntVar(&opts.BlockSize, "block-size", 1, "Sets the size of one block in bytes for reading and writing.")

	flag.IntVar(&opts.ConvType, "conv", 0, "Sets text conversion options. "+
		"By default original text is copied without changes")

	// todo: parse and validate all flags

	flag.Parse()

	return &opts, nil
}

func main() {
	opts, err := ParseFlags()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "can not parse flags:", err)
		os.Exit(1)
	}

	fmt.Println(opts)

	// todo: implement the functional requirements described in read.me
}
