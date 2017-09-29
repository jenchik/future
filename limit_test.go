package future

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLimitFutureEmpty(t *testing.T) {
	l := NewLimitFutures(10)

	assert.Equal(t, 10, l.Count())
	assert.Zero(t, l.Queue())
}

func TestLimitFutureAddTasks(t *testing.T) {
	l := NewLimitFutures(10)

	f1 := l.AddTask(func() (interface{}, error) {
		time.Sleep(time.Millisecond)
		return 100, nil
	})

	assert.Equal(t, 10, l.Count())
	assert.Equal(t, 1, l.Queue())

	f2 := l.AddTask(func() (interface{}, error) {
		assert.Equal(t, 2, l.Queue())
		v, e := f1.Wait(nil)
		assert.Equal(t, 1, l.Queue())
		return v, e
	})

	assert.Equal(t, 10, l.Count())
	assert.Equal(t, 2, l.Queue())

	val, err := f2.Wait(nil)

	assert.Equal(t, 10, l.Count())
	assert.Zero(t, l.Queue())

	assert.Equal(t, 100, val)
	assert.NoError(t, err)
}

func TestLimitFutureVerifyLimits(t *testing.T) {
	start := time.Now()
	l := NewLimitFutures(10)

	var f Future
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		k := i
		f = l.AddTask(func() (interface{}, error) {
			time.Sleep(time.Millisecond * 100)
			return 100 + k, nil
		})
		go func(ff Future) {
			ff.Wait(nil)
			wg.Done()
		}(f)
	}

	assert.Equal(t, 10, l.Count())
	assert.Equal(t, 10, l.Queue())

	wg.Wait()

	assert.Equal(t, 10, l.Count())
	assert.Zero(t, l.Queue())

	assert.InDelta(t, 0.1, time.Since(start).Seconds(), 0.01)
}

func TestLimitFutureVerifyLimits2(t *testing.T) {
	start := time.Now()
	l := NewLimitFutures(10)

	var wg, wr sync.WaitGroup
	wg.Add(15)
	wr.Add(10)
	go func() {
		for i := 0; i < 15; i++ {
			k := i
			l.AddTask(func() (interface{}, error) {
				defer wg.Done()
				time.Sleep(time.Millisecond * 100)
				return 100 + k, nil
			})
			if i < 10 {
				wr.Done()
			}
		}
	}()

	wr.Wait()

	assert.Equal(t, 10, l.Count())
	assert.Equal(t, 10, l.Queue())

	wg.Wait()

	assert.Equal(t, 10, l.Count())
	assert.Zero(t, l.Queue())

	assert.InDelta(t, 0.2, time.Since(start).Seconds(), 0.01, time.Since(start).String())
}

func TestLimitFutureVerifyNoLimits(t *testing.T) {
	start := time.Now()
	l := NewLimitFutures(-1)

	var f Future
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		k := i
		f = l.AddTask(func() (interface{}, error) {
			time.Sleep(time.Millisecond * 100)
			return 100 + k, nil
		})
		go func(ff Future) {
			ff.Wait(nil)
			wg.Done()
		}(f)
	}

	assert.Zero(t, l.Count())
	assert.Zero(t, l.Queue())

	wg.Wait()

	assert.Zero(t, l.Count())
	assert.Zero(t, l.Queue())

	assert.InDelta(t, 0.1, time.Since(start).Seconds(), 0.01)
}
