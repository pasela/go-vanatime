package vanatime

import (
	"errors"
	"sync"
	"time"
)

// A Ticker holds a channel that delivers `ticks' of a clock
// at intervals.
type Ticker struct {
	C <-chan Time

	c           chan Time
	earthTicker *time.Ticker
	stop        chan struct{}
	wg          sync.WaitGroup
}

// NewTicker returns a new Ticker containing a channel that will send the
// time with a period specified by the duration argument.
// It adjusts the intervals or drops ticks to make up for slow receivers.
// The duration d must be greater than zero; if not, NewTicker will panic.
// Stop the ticker to release associated resources.
func NewTicker(d Duration) *Ticker {
	if d <= 0 {
		panic(errors.New("non-positive interval for NewTicker"))
	}

	earthTicker := time.NewTicker(vd2ed(d))

	c := make(chan Time, 1)
	t := &Ticker{
		C:           c,
		c:           c,
		earthTicker: earthTicker,
		stop:        make(chan struct{}),
	}

	t.start()

	return t
}

func (t *Ticker) start() {
	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		for {
			select {
			case value, ok := <-t.earthTicker.C:
				if ok {
					t.c <- FromEarth(value)
				} else {
					close(t.c)
					return
				}

			case <-t.stop:
				return
			}
		}
	}()
}

// Stop turns off a ticker. After Stop, no more ticks will be sent.
// Stop does not close the channel, to prevent a concurrent goroutine
// reading from the channel from seeing an erroneous "tick".
func (t *Ticker) Stop() {
	t.earthTicker.Stop()
	close(t.stop)
	t.wg.Wait()
}

// Tick is a convenience wrapper for NewTicker providing access to the ticking
// channel only. While Tick is useful for clients that have no need to shut down
// the Ticker, be aware that without a way to shut it down the underlying
// Ticker cannot be recovered by the garbage collector; it "leaks".
// Unlike NewTicker, Tick will return nil if d <= 0.
func Tick(d Duration) <-chan Time {
	if d <= 0 {
		return nil
	}
	return NewTicker(d).C
}
