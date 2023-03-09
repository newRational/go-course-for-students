package main

import (
	"fmt"
	"lecture03_homework/dd"
	"os"
)

func main() {
	opts, err := dd.ParseFlags()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "cannot parse flags:", err)
		os.Exit(1)
	}

	fmt.Println(opts)

	dd.Start(opts)
}
