package racer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRace(t *testing.T) {
	slowFn := func(chan struct{}) (interface{}, error) {
		time.Sleep(2 * time.Millisecond)
		return "I'm slow", nil
	}
	fastFn := func(chan struct{}) (interface{}, error) {
		time.Sleep(time.Millisecond)
		return "I'm fast", nil
	}

	res, err := Race(nil, slowFn, fastFn)
	assert.NoError(t, err)
	assert.Equal(t, "I'm fast", res)
}

func TestRaceNoRacers(t *testing.T) {
	res, err := Race(nil)
	assert.Equal(t, errNoRacers, err)
	assert.Nil(t, res)
}

func TestTimeoutFirst(t *testing.T) {
	slowFn1 := func(chan struct{}) (interface{}, error) {
		time.Sleep(time.Millisecond)
		return nil, nil
	}
	slowFn2 := func(chan struct{}) (interface{}, error) {
		time.Sleep(time.Millisecond)
		return nil, nil
	}

	opts := &Options{Timeout: time.Nanosecond}
	_, err := Race(opts, slowFn1, slowFn2)
	assert.Equal(t, ErrTimeout, err)
}

func TestTimeoutLater(t *testing.T) {
	fn1 := func(chan struct{}) (interface{}, error) {
		time.Sleep(time.Millisecond)
		return nil, nil
	}
	fn2 := func(chan struct{}) (interface{}, error) {
		time.Sleep(time.Millisecond)
		return nil, nil
	}

	opts := &Options{Timeout: time.Millisecond * 5}
	_, err := Race(opts, fn1, fn2)
	assert.NoError(t, err)
}

func TestKillFirst(t *testing.T) {
	fn1 := func(chan struct{}) (interface{}, error) {
		time.Sleep(time.Second)
		return nil, nil
	}
	fn2 := func(chan struct{}) (interface{}, error) {
		time.Sleep(time.Second)
		return nil, nil
	}

	killCh := func() chan struct{} {
		ch := make(chan struct{})
		go func() {
			time.Sleep(time.Millisecond)
			close(ch)
		}()
		return ch
	}()

	opts := &Options{Kill: killCh}
	_, err := Race(opts, fn1, fn2)
	assert.Equal(t, ErrKilled, err)
}

func TestKillLater(t *testing.T) {
	fn1 := func(chan struct{}) (interface{}, error) {
		time.Sleep(time.Millisecond)
		return nil, nil
	}
	fn2 := func(chan struct{}) (interface{}, error) {
		time.Sleep(time.Millisecond)
		return nil, nil
	}

	killCh := func() chan struct{} {
		ch := make(chan struct{})
		go func() {
			time.Sleep(time.Second)
			close(ch)
		}()
		return ch
	}()

	opts := &Options{Kill: killCh}
	_, err := Race(opts, fn1, fn2)
	assert.NoError(t, err)
}
