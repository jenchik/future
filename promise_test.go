package future

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPromise(t *testing.T) {
	f := func() (interface{}, error) {
		time.Sleep(100 * time.Millisecond)
		return 42, nil
	}
	p := NewPromise(f)
	value, err := p.Wait(nil)

	assert.Equal(t, 42, value)
	assert.NoError(t, err)
}

func TestPromiseThen(t *testing.T) {
	p := NewPromise(func() (interface{}, error) {
		time.Sleep(100 * time.Millisecond)
		return 42, nil
	})
	value, err := p.Then(func(value interface{}) (interface{}, error) {
		time.Sleep(1 * time.Second)
		return value.(int) + 3, nil
	}).Wait(nil)

	assert.Equal(t, 45, value)
	assert.NoError(t, err)
}

func TestPromiseThenWithCancel(t *testing.T) {
	cancel := make(chan struct{})
	p := NewPromise(func() (interface{}, error) {
		time.Sleep(500 * time.Millisecond)
		return 42, nil
	})
	time.AfterFunc(100*time.Millisecond, func() {
		close(cancel)
	})
	value, err := p.Then(func(value interface{}) (interface{}, error) {
		time.Sleep(1 * time.Second)
		return value.(int) + 3, nil
	}).Wait(cancel)

	assert.Equal(t, nil, value)
	assert.Error(t, err)
}

func TestPromiseErrorThen(t *testing.T) {
	p := NewPromise(func() (interface{}, error) {
		time.Sleep(100 * time.Millisecond)
		return nil, errors.New("error!")
	})
	value, err := p.Then(func(value interface{}) (interface{}, error) {
		time.Sleep(100 * time.Millisecond)
		return value.(int) + 3, nil
	}).Wait(nil)

	assert.Nil(t, value)
	assert.Error(t, err)
}

func TestPromiseChain(t *testing.T) {
	value, err := NewPromise(func() (interface{}, error) {
		time.Sleep(10 * time.Millisecond)
		return 20, nil
	}).Then(func(value interface{}) (interface{}, error) {
		time.Sleep(10 * time.Millisecond)
		return value.(int) - 10, nil
	}).Then(func(value interface{}) (interface{}, error) {
		time.Sleep(10 * time.Millisecond)
		return value.(int) * 3, nil
	}).Then(func(value interface{}) (interface{}, error) {
		time.Sleep(10 * time.Millisecond)
		return value.(int) / 5, nil
	}).Wait(nil)

	assert.Equal(t, 6, value)
	assert.NoError(t, err)
}

func TestPromiseChainDelay(t *testing.T) {
	start := time.Now()

	p1 := NewPromise(func() (interface{}, error) {
		time.Sleep(100 * time.Millisecond)
		return 42, nil
	})

	p2 := p1.Then(func(value interface{}) (interface{}, error) {
		time.Sleep(100 * time.Millisecond)
		return value.(int) + 3, nil
	})

	assert.InDelta(t, 0.0, time.Since(start).Seconds(), 0.01)

	value, err := p2.Wait(nil)
	assert.Equal(t, 45, value)
	assert.NoError(t, err)

	assert.InDelta(t, 0.2, time.Since(start).Seconds(), 0.05)
}
