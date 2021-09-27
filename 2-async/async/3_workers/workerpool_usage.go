package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Job interface {
	Do(ctx context.Context)
}

type WorkerPool interface {
	Do(j Job)
	Stop()
}

func NewWorkerPool(workersNum int) WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	input := make(chan Job, workersNum)

	wg := &sync.WaitGroup{}
	wg.Add(workersNum)
	for i := 0; i < workersNum; i++ {
		go func(ctx context.Context, wg *sync.WaitGroup) {
			for {
				select {
				case <-ctx.Done():
					wg.Done()
					return
				case job := <-input:
					job.Do(ctx)
				}
			}
		}(ctx, wg)
	}
	return &Pool{
		wg:     wg,
		cancel: cancel,
		input:  input,
	}
}

type Pool struct {
	wg     *sync.WaitGroup
	cancel context.CancelFunc
	input  chan Job
}

func (p *Pool) Stop() {
	p.cancel()
	p.wg.Wait()
}

func (p *Pool) Do(j Job) {
	p.input <- j
}

type LongJob struct {
	result *int64
	input  int64
}

func (j *LongJob) Do(ctx context.Context) {
	longOperation := func() int64 {
		time.Sleep(100 * time.Millisecond)
		return j.input * 2
	}

	atomic.StoreInt64(j.result, longOperation())
}

func main() {
	const (
		workers = 10
		tasks   = 1000
	)

	pool := NewWorkerPool(workers)

	res := make([]int64, tasks)

	for i := 0; i < tasks; i++ {
		pool.Do(&LongJob{input: int64(i), result: &res[i]})
	}

	pool.Stop()

	for _, val := range res {
		fmt.Println(val)
	}
}
