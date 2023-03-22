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
	last := make(chan any)

	o := stages[0](in)
	for i := 1; i < len(stages); i++ {
		o = stages[i](o)
	}

	go func() {
		defer close(last)
		ok := true
		var data any
		for ok && ctx.Err() == nil {
			select {
			case data, ok = <-o:
				if !ok {
					break
				}
				last <- data
			default:
			}
		}
	}()

	return last
}
