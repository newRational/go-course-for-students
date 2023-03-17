package storage

import (
	"context"
	"sync"
)

// Result represents the Size function result
type Result struct {
	// Total Size of File objects
	Size int64
	// Count is a count of File objects processed
	Count int64
}

type DirSizer interface {
	// Size calculate a size of given Dir, receive a ctx and the root Dir instance
	// will return Result or error if happened
	Size(ctx context.Context, d Dir) (Result, error)
}

// sizer implement the DirSizer interface
type sizer struct {
	// maxWorkersCount number of workers for asynchronous run
	maxWorkersCount int

	// TODO: add other fields as you wish
}

// NewSizer returns new DirSizer instance
func NewSizer() DirSizer {
	return &sizer{}
}

func (a *sizer) Size(ctx context.Context, d Dir) (Result, error) {
	var res Result

	if err := ctx.Err(); err != nil {
		return res, err
	}

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

	if len(dirs) == 0 {
		return res, nil
	}

	chDirsErr := make(chan error, len(dirs))
	chDirsRes := make(chan Result, len(dirs))

	var wg sync.WaitGroup
	wg.Add(len(dirs))

	for _, dir := range dirs {
		go func(dir Dir) {
			defer wg.Done()
			a.processDir(ctx, dir, chDirsRes, chDirsErr)
		}(dir)
	}

	wg.Wait()
	close(chDirsErr)
	close(chDirsRes)

	for e := range chDirsErr {
		if e != nil {
			return res, e
		}
	}

	for r := range chDirsRes {
		res.Size += r.Size
		res.Count += r.Count
	}

	return res, nil
}

func (a *sizer) processDir(ctx context.Context, dir Dir, chRes chan<- Result, chErr chan<- error) {
	var res Result

	if err := ctx.Err(); err != nil {
		chErr <- err
		chRes <- res
		return
	}

	dirRes, err := a.Size(ctx, dir)
	if err != nil {
		chErr <- err
		chRes <- res
		return
	}

	res.Size += dirRes.Size
	res.Count += dirRes.Count

	chErr <- nil
	chRes <- res
}
