package racer

import (
	"fmt"
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

func TestFailures(t *testing.T) {
	failureFn := func(chan struct{}) (interface{}, error) {
		return nil, fmt.Errorf("I failed")
	}
	slowlyButSurely := func(chan struct{}) (interface{}, error) {
		time.Sleep(time.Millisecond)
		return true, nil
	}

	var racers []Racer
	for i := 0; i < 10; i++ {
		racers = append(racers, failureFn)
	}
	racers = append(racers, slowlyButSurely)
	res, err := Race(nil, racers...)
	assert.NoError(t, err)
	assert.Equal(t, true, res)
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
