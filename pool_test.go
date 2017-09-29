package future

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPoolFutureEmpty(t *testing.T) {
	p := NewPoolFutures(5, 10, false)

	assert.Equal(t, 5, p.Count())
	assert.Zero(t, p.Queue())
}

func TestPoolFutureAddTasks(t *testing.T) {
	start := time.Now()
	p := NewPoolFutures(2, 1, false)

	assert.Equal(t, 2, p.Count())
	assert.Equal(t, 0, p.Queue())

	lock1 := make(chan struct{})
	f1 := p.AddTask(func() (interface{}, error) {
		<-lock1
		return 100, nil
	})

	f2 := p.AddTask(func() (interface{}, error) {
		time.Sleep(time.Millisecond * 100)
		return 200, nil
	})

	time.Sleep(time.Microsecond)
	assert.Equal(t, 0, p.Queue())

	f3 := p.AddTask(func() (interface{}, error) {
		v, e := f1.Wait(nil)
		return v, e
	})
	assert.Equal(t, 1, p.Queue())

	var f4 Future
	ready := make(chan struct{})
	go func() {
		f4 = p.AddTask(func() (interface{}, error) {
			time.Sleep(time.Millisecond * 100)
			return 400, nil
		})
		close(ready)
	}()
	assert.Equal(t, 1, p.Queue())

	close(lock1)
	assert.Equal(t, 1, p.Queue())

	<-ready
	assert.Equal(t, 0, p.Queue())

	val, err := f3.Wait(nil)
	assert.Equal(t, 0, p.Queue())
	assert.Equal(t, 100, val)
	assert.NoError(t, err)

	val, err = f2.Wait(nil)
	assert.Equal(t, 200, val)
	assert.NoError(t, err)

	val, err = f4.Wait(nil)
	assert.Equal(t, 400, val)
	assert.NoError(t, err)

	assert.InDelta(t, 0.1, time.Since(start).Seconds(), 0.01)
}

func TestPoolFutureAddTasksWithoutLimit(t *testing.T) {
	start := time.Now()
	p := NewPoolFutures(2, 1, true)

	assert.Equal(t, 2, p.Count())
	assert.Equal(t, 0, p.Queue())

	var wg sync.WaitGroup
	ready := make(chan struct{})
	wg.Add(10)
	for i := 0; i < 10; i++ {
		p.AddTask(func() (interface{}, error) {
			defer wg.Done()
			<-ready
			time.Sleep(time.Millisecond * 100)
			return i * 10, nil
		})
	}
	assert.Equal(t, 2, p.Count())
	assert.Equal(t, 1, p.Queue())

	close(ready)
	assert.Equal(t, 2, p.Count())
	assert.Equal(t, 1, p.Queue())

	wg.Wait()
	assert.Equal(t, 2, p.Count())
	assert.Equal(t, 0, p.Queue())

	// 2 worker (100ms) + 1 in queue (100ms)
	assert.InDelta(t, 0.2, time.Since(start).Seconds(), 0.01)
}
