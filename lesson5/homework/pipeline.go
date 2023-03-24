package executor

import (
	"context"
)

type (
	In  <-chan any
	Out = In
)

type Stage func(in In) (out Out)

func ExecutePipeline(ctx context.Context, in In, stages ...Stage) Out {
	out := make(chan any)

	for _, s := range stages {
		in = s(in)
	}

	go func() {
		defer close(out)
		for {
			select {
			case data, ok := <-in:
				if !ok {
					return
				}
				out <- data
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}
