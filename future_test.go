package future

import (
	"errors"
	"testing"
	"time"

	"context"

	"github.com/stretchr/testify/assert"
)

func TestFutureError(t *testing.T) {
	f1 := NewFuture(func() (interface{}, error) {
		time.Sleep(100 * time.Millisecond)
		return nil, errors.New("test error")
	})

	value, err := f1.Wait(nil)
	assert.Error(t, err)
	assert.Nil(t, value)
}

func TestFutureAsync(t *testing.T) {
	start := time.Now()
	f1 := NewFuture(func() (interface{}, error) {
		time.Sleep(100 * time.Millisecond)
		return 42, nil
	})

	f2 := NewFuture(func() (interface{}, error) {
		time.Sleep(100 * time.Millisecond)
		return 43, nil
	})
	value, err := f1.Wait(nil)
	assert.Equal(t, 42, value)
	assert.NoError(t, err)

	value, err = f2.Wait(nil)
	assert.Equal(t, 43, value)
	assert.NoError(t, err)

	assert.InDelta(t, 0.1, time.Since(start).Seconds(), 0.01)
}

func TestFutureGroupRead(t *testing.T) {
	cancel := make(chan struct{})
	start := time.Now()
	f1 := NewFuture(func() (interface{}, error) {
		time.Sleep(100 * time.Millisecond)
		return 42, nil
	})

	f2 := NewFuture(func() (interface{}, error) {
		time.Sleep(200 * time.Millisecond)
		return 43, nil
	})

	time.AfterFunc(100*time.Millisecond, func() {
		close(cancel)
	})

	value, err := f1.Wait(nil)
	assert.Equal(t, 42, value)
	assert.NoError(t, err)

	value, err = f2.Wait(cancel)
	assert.Equal(t, nil, value)
	assert.Error(t, err)

	value, err = f1.Wait(nil)
	assert.Equal(t, 42, value)
	assert.NoError(t, err)

	assert.InDelta(t, 0.1, time.Since(start).Seconds(), 0.01)
}

func TestFutureWithCancel(t *testing.T) {
	cancel := make(chan struct{})
	start := time.Now()
	f1 := NewFuture(func() (interface{}, error) {
		time.Sleep(50 * time.Millisecond)
		return 42, nil
	})

	f2 := NewFuture(func() (interface{}, error) {
		time.Sleep(100 * time.Millisecond)
		return 43, nil
	})

	f3 := NewFuture(func() (interface{}, error) {
		time.Sleep(500 * time.Millisecond)
		return 44, nil
	})

	time.AfterFunc(100*time.Millisecond, func() {
		close(cancel)
	})

	value, err := f1.Wait(cancel)
	assert.Equal(t, 42, value)
	assert.NoError(t, err)

	value, err = f2.Wait(nil)
	assert.Equal(t, 43, value)
	assert.NoError(t, err)

	value, err = f3.Wait(cancel)
	assert.Equal(t, nil, value)
	assert.Error(t, err)

	assert.InDelta(t, 0.1, time.Since(start).Seconds(), 0.01)
}

func TestFutureWithTimeout(t *testing.T) {
	start := time.Now()
	f := NewFuture(func() (interface{}, error) {
		time.Sleep(1 * time.Second)
		return 42, nil
	})

	value, err := f.WaitWithTimeout(100 * time.Millisecond)
	assert.Error(t, err)
	assert.Equal(t, ErrTimeout, err)
	assert.Nil(t, value)

	assert.InDelta(t, 0.1, time.Since(start).Seconds(), 0.01)
}

func TestFutureWithTimeoutComplete(t *testing.T) {
	start := time.Now()
	f := NewFuture(func() (interface{}, error) {
		time.Sleep(100 * time.Millisecond)
		return 42, nil
	})

	value, err := f.WaitWithTimeout(1 * time.Second)
	assert.Equal(t, 42, value)
	assert.NoError(t, err)

	assert.InDelta(t, 0.1, time.Since(start).Seconds(), 0.01)
}

func TestFutureWithContext(t *testing.T) {
	f := NewFuture(func() (interface{}, error) {
		time.Sleep(100 * time.Millisecond)
		return 42, nil
	})

	ctx := context.Background()
	value, err := f.WaitWithContext(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 42, value)
}

func TestFutureWithContextCancel(t *testing.T) {
	f := NewFuture(func() (interface{}, error) {
		time.Sleep(100 * time.Millisecond)
		return 42, nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	value, err := f.WaitWithContext(ctx)
	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
	assert.Nil(t, value)
}

func TestFutureWithContextTimeout(t *testing.T) {
	f := NewFuture(func() (interface{}, error) {
		time.Sleep(100 * time.Millisecond)
		return 42, nil
	})

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Millisecond)
	value, err := f.WaitWithContext(ctx)
	assert.Error(t, err)
	assert.Equal(t, context.DeadlineExceeded, err)
	assert.Nil(t, value)
}
