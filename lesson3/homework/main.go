package main

import (
	"flag"
	"fmt"
	"lecture03_homework/dd"
	"os"
)

func ParseFlags() (*dd.Options, []error) {
	var opts dd.Options

	dd.DefineFlags(&opts)
	flag.Parse()
	invalidFlags := dd.ValidateFlags(&opts)

	return &opts, invalidFlags
}

func main() {
	opts, err := ParseFlags()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "cannot parse flags:", err)
		os.Exit(1)
	}

	fmt.Println(opts)

	// todo: implement the functional requirements described in read.me
}
