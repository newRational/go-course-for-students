package storage

import (
	"context"
	"github.com/gammazero/workerpool"
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
	wp *workerpool.WorkerPool
}

// NewSizer returns new DirSizer instance
func NewSizer() DirSizer {
	return &sizer{
		wp: workerpool.New(4),
	}
}

func (r *Result) add(other Result) {
	r.Size += other.Size
	r.Count += other.Count
}

func (a *sizer) Size(ctx context.Context, d Dir) (res Result, err error) {
	defer func() {
		if ctxErr := ctx.Err(); ctxErr != nil {
			err = ctxErr
		}
	}()

	dirs, files, err := d.Ls(ctx)
	if err != nil {
		return res, err
	}

	if len(dirs) != 0 {
		dirsRes, err := a.processDirsAsync(ctx, dirs)
		if err != nil {
			return res, err
		}
		res.add(dirsRes)
	}

	filesRes, err := a.processFiles(ctx, files)
	if err != nil {
		return res, err
	}
	res.add(filesRes)

	return res, nil
}

func (a *sizer) processFiles(ctx context.Context, files []File) (res Result, err error) {
	defer func() {
		if ctxErr := ctx.Err(); ctxErr != nil {
			err = ctxErr
		}
	}()

	for _, file := range files {
		size, err := file.Stat(ctx)
		if err != nil {
			return res, err
		}
		res.Size += size
		res.Count++
	}

	return res, nil
}

func (a *sizer) processDirsAsync(ctx context.Context, dirs []Dir) (res Result, err error) {
	defer func() {
		if ctxErr := ctx.Err(); ctxErr != nil {
			err = ctxErr
		}
	}()

	chCtxErr := make(chan error, 1)
	chErr := make(chan error, len(dirs))
	chRes := make(chan Result, len(dirs))

	var wg sync.WaitGroup
	wg.Add(len(dirs))

	for _, dir := range dirs {
		f := closure(ctx, a, dir, &wg, chRes, chErr, chCtxErr)
		a.wp.Submit(f)
	}

	wg.Wait()
	close(chCtxErr)
	close(chErr)
	close(chRes)

	if err := <-chCtxErr; err != nil {
		return res, err
	}

	for err := range chErr {
		if err != nil {
			return res, err
		}
	}

	for r := range chRes {
		res.Size += r.Size
		res.Count += r.Count
	}

	return res, nil
}

func (a *sizer) processDir(ctx context.Context, dir Dir, chRes chan<- Result, chErr, chCtxErr chan<- error) {
	defer func() {
		if ctxErr := ctx.Err(); ctxErr != nil && len(chCtxErr) == 0 {
			chCtxErr <- ctxErr
		}
	}()

	var res Result

	dirRes, err := a.Size(ctx, dir)

	if err != nil {
		chErr <- err
		chRes <- res
		return
	}

	res.Size += dirRes.Size
	res.Count += dirRes.Count

	chRes <- res
}

// closure returns closure for WorkerPool's Submit method
func closure(ctx context.Context, a *sizer, dir Dir, wg *sync.WaitGroup, chRes chan<- Result, chErr, chCtxErr chan<- error) func() {
	return func() {
		defer wg.Done()
		a.processDir(ctx, dir, chRes, chErr, chCtxErr)
	}
}
