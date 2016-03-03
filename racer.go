package racer

import (
	"errors"
	"reflect"
	"time"
)

var (
	ErrTimeout = errors.New("all operations timed out")
	ErrKilled  = errors.New("all operations were killed")
)

var (
	errUnknown  = errors.New("unknown")
	errNoRacers = errors.New("no racers")
)

type Racer func(chan struct{}) (interface{}, error)

type Options struct {
	Timeout time.Duration
	Delay   time.Duration
	Kill    chan struct{}
}

func Race(opts *Options, racers ...Racer) (interface{}, error) {
	if len(racers) == 0 {
		return nil, errNoRacers
	}

	done := make(chan struct{})
	defer close(done)

	cases := make([]reflect.SelectCase, 0, len(racers)+10)

	for _, r := range racers {
		cases = append(cases, newCase(r, done))
	}

	if opts != nil && opts.Timeout != 0 {
		cases = append(cases, newCase(newTimeoutRacer(opts.Timeout), done))
	}

	if opts != nil && opts.Kill != nil {
		cases = append(cases, newCase(newKillRacer(opts.Kill), done))
	}

	_, val, _ := reflect.Select(cases)
	res := val.Interface().(result)
	return res.res, res.err
}

func newCase(r Racer, done chan struct{}) reflect.SelectCase {
	return reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(run(r, done)),
	}
}

func newKillRacer(kill chan struct{}) func(chan struct{}) (interface{}, error) {
	return func(done chan struct{}) (interface{}, error) {
		select {
		case <-kill:
			return nil, ErrKilled
		case <-done:
			return nil, nil
		}
	}
}

func newTimeoutRacer(d time.Duration) func(chan struct{}) (interface{}, error) {
	return func(done chan struct{}) (interface{}, error) {
		select {
		case <-time.After(d):
			return nil, ErrTimeout
		case <-done:
			return nil, nil
		}
	}
}

func run(r Racer, kill chan struct{}) <-chan result {
	ch := make(chan result)
	go func() {
		res, err := r(kill)
		ch <- result{
			res: res,
			err: err,
		}
	}()
	return ch
}

type result struct {
	res interface{}
	err error
}
