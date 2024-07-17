package storage

import (
	"context"
	"sync"
	"sync/atomic"
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
	jobs            chan job
	// TODO: add other fields as you wish
}

func (a *sizer) enqueue(j job) {
	j.wg.Add(1)
	select {
	case a.jobs <- j:
	default:
		j.do(a.enqueue)
		j.wg.Done()
	}
}

// NewSizer returns new DirSizer instance
func NewSizer() DirSizer {
	jobs := make(chan job)
	maxWorkersCount := 5

	sizer := &sizer{maxWorkersCount: maxWorkersCount, jobs: jobs}

	for i := 0; i < maxWorkersCount; i++ {
		go func() {
			for j := range jobs {
				j.do(sizer.enqueue)
				j.wg.Done()
			}
		}()
	}

	return sizer
}

type job struct {
	directory Dir
	file      File
	context   context.Context
	res       *Result
	err       *isError
	wg        *sync.WaitGroup
}

type isError struct {
	err error
}

func (j *job) do(enqueue func(job)) {
	if j.directory != nil {
		dirs, files, err := j.directory.Ls(j.context)
		if err != nil {
			*j.err = isError{err: err}
			return
		}
		for _, d := range dirs {
			enqueue(job{d, nil, j.context, j.res, j.err, j.wg})
		}
		for _, f := range files {
			enqueue(job{nil, f, j.context, j.res, j.err, j.wg})
		}
		return
	}
	val, err := j.file.Stat(j.context)
	if err != nil {
		*j.err = isError{err: err}
		return
	}
	atomic.AddInt64(&j.res.Size, val)
	atomic.AddInt64(&j.res.Count, 1)
}

func (a *sizer) Size(ctx context.Context, d Dir) (Result, error) {
	wg := new(sync.WaitGroup)

	res := &Result{0, 0}
	err := &isError{nil}

	a.enqueue(job{d, nil, ctx, res, err, wg})

	wg.Wait()
	close(a.jobs)
	return *res, err.err

}
