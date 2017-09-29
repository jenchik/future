package future

import (
	"context"
	"time"
)

type (
	Task func() (interface{}, error)

	Async func(Task) Future

	Future interface {
		Wait(<-chan struct{}) (interface{}, error)
		WaitWithTimeout(time.Duration) (interface{}, error)
		WaitWithContext(context.Context) (interface{}, error)
	}

	futureResult struct {
		done  chan struct{}
		value interface{}
		err   error
	}
)

var (
	ErrTimeout  = context.DeadlineExceeded
	ErrCanceled = context.Canceled
)

func NewFuture(fn func() (interface{}, error)) Future {
	return newFutureResult(fn)
}

func newFutureResult(t Task) *futureResult {
	f := &futureResult{
		done: make(chan struct{}),
	}

	go f.run(t, func() {
		close(f.done)
	})

	return f
}

func (f *futureResult) run(t Task, done func()) {
	if done != nil {
		defer done()
	}

	f.value, f.err = t()
}

func (f *futureResult) Wait(c <-chan struct{}) (interface{}, error) {
	select {
	case <-c:
		return nil, ErrCanceled
	case <-f.done:
		return f.value, f.err
	}
}

func (f *futureResult) WaitWithTimeout(timeout time.Duration) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return f.WaitWithContext(ctx)
}

func (f *futureResult) WaitWithContext(ctx context.Context) (interface{}, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-f.done:
		return f.value, f.err
	}
}
