package dd

import (
	"fmt"
	"log"
	"os"
)

func ReadFileWithChunk(opts *Options) error {
	chunk := make([]byte, opts.BlockSize)

	file, err := os.Open(opts.From)
	if err != nil {
		log.Fatal(err)
	}

	file.Seek(int64(opts.Offset), 0)

	for {
		bytesRead, _ := file.Read(chunk)
		if bytesRead == 0 {
			break
		}
		process(chunk)
	}

	if err := file.Close(); err != nil {
		return err
	}

	return nil
}

func process(chunk []byte) {
	fmt.Print(string(chunk))
}
