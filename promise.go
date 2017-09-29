package future

type Promise interface {
	Wait(<-chan struct{}) (interface{}, error)
	Then(func(value interface{}) (interface{}, error)) Promise
}

func NewPromise(f func() (interface{}, error)) Promise {
	return newFutureResult(f)
}

func (p *futureResult) Then(f func(value interface{}) (interface{}, error)) Promise {
	next := &futureResult{
		done: make(chan struct{}),
	}

	go func() {
		<-p.done
		if p.err != nil {
			next.value, next.err = p.value, p.err
			close(next.done)
			return
		}

		defer close(next.done)
		next.value, next.err = f(p.value)
	}()

	return next
}
