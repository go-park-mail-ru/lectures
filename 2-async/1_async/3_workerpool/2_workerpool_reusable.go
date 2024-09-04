package main

type workerPool struct {
	locker chan struct{}
	queue  chan func()
}

func NewWorkerPool(maxWorkers int) workerPool {
	wp := workerPool{
		locker: make(chan struct{}, maxWorkers),
		queue:  make(chan func(), maxWorkers),
	}
	go wp.work()

	return wp
}

func (wp workerPool) work() {
	for {
		task := <-wp.queue
		wp.locker <- struct{}{}

		go func() {
			defer func() { <-wp.locker }()
			task()
		}()
	}
}

func (wp workerPool) PutTask(task func()) {
	wp.queue <- task
}
