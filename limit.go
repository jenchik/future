package future

type (
	limitFutures struct {
		jobsChan   chan struct{}
		maxThreads int
	}
)

// NewLimitFutures max limitation goroutines, which can serve the futures
// If maxGoroutines < 1, then not limit
func NewLimitFutures(maxGoroutines int) *limitFutures {
	l := &limitFutures{
		maxThreads: maxGoroutines,
	}
	if l.maxThreads > 0 {
		l.jobsChan = make(chan struct{}, l.maxThreads)
	}

	return l
}

func (l *limitFutures) Count() int {
	return cap(l.jobsChan)
}

func (l *limitFutures) Queue() int {
	return len(l.jobsChan)
}

func (l *limitFutures) process(t Task) Future {
	f := &futureResult{
		done: make(chan struct{}),
	}

	l.jobsChan <- struct{}{}

	go f.run(t, func() {
		close(f.done)
		<-l.jobsChan
	})

	return f
}

func (l *limitFutures) AddTask(f func() (interface{}, error)) Future {
	if l.maxThreads < 1 {
		return NewFuture(f)
	}

	return l.process(f)
}
