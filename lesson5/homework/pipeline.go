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
	o := stages[0](in)
	for i := 1; i < len(stages); i++ {
		o = stages[i](o)
	}

	return o
}
