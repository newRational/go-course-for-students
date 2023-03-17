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
	//maxWorkersCount int

	// TODO: add other fields as you wish
}

// NewSizer returns new DirSizer instance
func NewSizer() DirSizer {
	return &sizer{}
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

	filesRes, err := a.processFiles(ctx, files)
	if err != nil {
		return res, err
	}
	res.add(filesRes)

	if len(dirs) != 0 {
		dirsRes, err := a.processDirsAsync(ctx, dirs)
		if err != nil {
			return res, err
		}
		res.add(dirsRes)
	}

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

	chErr := make(chan error, len(dirs))
	chRes := make(chan Result, len(dirs))

	var wg sync.WaitGroup
	wg.Add(len(dirs))

	for _, dir := range dirs {
		go func(dir Dir) {
			defer wg.Done()
			a.processDir(ctx, dir, chRes, chErr)
		}(dir)
	}

	wg.Wait()
	close(chErr)
	close(chRes)

	for e := range chErr {
		if e != nil {
			return res, e
		}
	}

	for r := range chRes {
		res.Size += r.Size
		res.Count += r.Count
	}

	return res, nil
}

func (a *sizer) processDir(ctx context.Context, dir Dir, chRes chan<- Result, chErr chan<- error) {
	defer func() {
		if ctxErr := ctx.Err(); ctxErr != nil {
			chErr <- ctxErr
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

	chErr <- nil
	chRes <- res
}
