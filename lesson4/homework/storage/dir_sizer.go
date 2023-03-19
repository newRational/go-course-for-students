package storage

import (
	"context"
	"sync"
)

type Result struct {
	Size  int64
	Count int64
}

type DirSizer interface {
	Size(ctx context.Context, d Dir) (Result, error)
}

type sizer struct {
	maxWorkersCount int
}

func NewSizer() DirSizer {
	return &sizer{maxWorkersCount: 1}
}

func (a *sizer) Size(ctx context.Context, d Dir) (res Result, err error) {
	defer func() {
		if ctxErr := ctx.Err(); ctxErr != nil {
			err = ctx.Err()
		}
	}()

	dirs, files, err := d.Ls(ctx)
	if err != nil {
		return res, err
	}

	for _, file := range files {
		size, err := file.Stat(ctx)
		if err != nil {
			return res, err
		}
		res.Size += size
		res.Count++
	}

	chRes := make(chan Result, len(dirs))
	chErr := make(chan error, len(dirs))

	var wg sync.WaitGroup
	wg.Add(len(dirs))

	for _, dir := range dirs {
		go func(chRes chan<- Result, chErr chan<- error, dir Dir) {
			defer wg.Done()
			res, err := a.Size(ctx, dir)
			chRes <- res
			chErr <- err
		}(chRes, chErr, dir)
	}

	wg.Wait()

	close(chRes)
	close(chErr)

	for err := range chErr {
		if err != nil {
			return res, err
		}
	}

	for r := range chRes {
		res.Size += r.Size
		res.Count += r.Count
	}

	return
}
