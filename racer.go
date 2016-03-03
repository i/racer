package racer

import (
	"errors"
	"reflect"
	"time"
)

// Errors returned by racer
var (
	ErrTimeout = errors.New("all operations timed out")
	ErrKilled  = errors.New("all operations were killed")
)

var (
	errUnknown  = errors.New("unknown")
	errNoRacers = errors.New("no racers")
)

// Racer is defined as a function that respects a done channel and returns an
// interface and an error. The done channel is closed when the race is
// considered to be over. This happens when another function successfully
// finishes without an error or some other condition is met (such as a timeout).
type Racer func(done chan struct{}) (result interface{}, err error)

// Options is a list of optional params to a race.
type Options struct {
	// Timeout the maximum time to wait before ending the race.
	Timeout time.Duration
	// Kill is a channel to signal all racers to stop. Close it to end the race.
	Kill chan struct{}
}

// Race starts a race between multiple racers. It works by calling each of the
// racers and waiting until one of them successfully finishes execution or a
// stopping condition is met. Racers that return a non-nil zero are
// disqualified and the race keeps going.
func Race(opts *Options, racers ...Racer) (res interface{}, err error) {
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

	for {
		i, val, _ := reflect.Select(cases)
		res := val.Interface().(result)

		if res.err != nil {
			if res.err == ErrKilled || res.err == ErrTimeout {
				return nil, res.err
			}
			cases = append(cases[:i], cases[i+1:]...)
			continue
		}
		return res.res, nil
	}
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
