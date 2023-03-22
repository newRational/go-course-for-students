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
	o := make(chan any)

	for _, s := range stages {
		in = s(in)
	}

	go func() {
		defer close(o)
		ok := true
		var data any
		for ok && ctx.Err() == nil {
			select {
			case data, ok = <-in:
				if !ok {
					break
				}
				o <- data
			default:
			}
		}
	}()

	return o
}
