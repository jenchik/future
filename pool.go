package future

import (
	"github.com/jenchik/thread"
)

type (
	poolFutures struct {
		allowExceed bool
		jobsChan    chan func()
		g           *thread.GroupWorkers
	}
)

// NewPoolFutures pool with ready handlers that can serve futures.
// maxGoroutines - number of handlers(goroutines) in the pool
// taskQueueLen - task queue for pool. If queue is full, then waiting or create new goroutine
// allowExceed - true, If pool is busy and queue is full, then a new goroutine will be created
//               false, waiting until the queue is free
func NewPoolFutures(maxGoroutines, taskQueueLen int, allowExceed bool) *poolFutures {
	p := &poolFutures{
		allowExceed: allowExceed,
		g:           thread.NewGroupWorkers(),
		jobsChan:    make(chan func(), taskQueueLen),
	}

	for i := 0; i < maxGoroutines; i++ {
		p.g.AddAsWorker(func(w *thread.Worker) {
			stop := w.StopC()
			for {
				select {
				case <-stop:
					return
				case job := <-p.jobsChan:
					job()
				}
			}
		}, nil)
	}

	return p
}

func (p *poolFutures) Close() error {
	return p.g.Close()
}

func (p *poolFutures) Wait() {
	p.g.Wait()
}

func (p *poolFutures) StopC() <-chan struct{} {
	return p.g.StopC()
}

func (p *poolFutures) Count() int {
	return p.g.Count()
}

func (p *poolFutures) Queue() int {
	return len(p.jobsChan)
}

func (p *poolFutures) process(f *futureResult) {
}

func (p *poolFutures) AddTask(fn func() (interface{}, error)) Future {
	f := &futureResult{
		done: make(chan struct{}),
	}

	do := func() {
		f.run(fn, func() {
			close(f.done)
		})
	}

	if p.allowExceed {
		select {
		case p.jobsChan <- do:
		default:
			go do()
		}
	} else {
		p.jobsChan <- do
	}

	return f
}
