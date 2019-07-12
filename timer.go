package vanatime

import (
	"sync"
	"time"
)

// The Timer type represents a single event.
// When the Timer expires, the current time will be sent on C,
// unless the Timer was created by AfterFunc.
// A Timer must be created with NewTimer or AfterFunc.
type Timer struct {
	C <-chan Time

	c          chan Time
	f          timerFunc
	earthTimer *time.Timer
	stop       chan struct{}
	wg         sync.WaitGroup
}

type timerFunc func(c chan<- Time, t Time)

// NewTimer creates a new Timer that will send
// the current time on its channel after at least duration d.
func NewTimer(d Duration) *Timer {
	return newTimer(d, func(c chan<- Time, t Time) {
		c <- t
	})
}

func newTimer(d Duration, f timerFunc) *Timer {
	earthTimer := time.NewTimer(vd2ed(d))

	c := make(chan Time, 1)
	t := &Timer{
		C:          c,
		c:          c,
		f:          f,
		earthTimer: earthTimer,
		stop:       make(chan struct{}),
	}

	t.start()

	return t
}

func (t *Timer) start() {
	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		for {
			select {
			case value, ok := <-t.earthTimer.C:
				if ok {
					t.f(t.c, FromEarth(value))
				} else {
					close(t.c)
				}
				return

			case <-t.stop:
				return
			}
		}
	}()
}

// Stop prevents the Timer from firing.
// It returns true if the call stops the timer, false if the timer has already
// expired or been stopped.
// Stop does not close the channel, to prevent a read from the channel succeeding
// incorrectly.
//
// To prevent a timer created with NewTimer from firing after a call to Stop,
// check the return value and drain the channel.
// For example, assuming the program has not received from t.C already:
//
// 	if !t.Stop() {
// 		<-t.C
// 	}
//
// This cannot be done concurrent to other receives from the Timer's
// channel.
//
// For a timer created with AfterFunc(d, f), if t.Stop returns false, then the timer
// has already expired and the function f has been started in its own goroutine;
// Stop does not wait for f to complete before returning.
// If the caller needs to know whether f is completed, it must coordinate
// with f explicitly.
func (t *Timer) Stop() {
	if t.earthTimer != nil {
		t.earthTimer.Stop()
		if t.stop != nil {
			close(t.stop)
			t.wg.Wait()
		}
	}
}

// Reset changes the timer to expire after duration d.
// It returns true if the timer had been active, false if the timer had
// expired or been stopped.
//
// Resetting a timer must take care not to race with the send into t.C
// that happens when the current timer expires.
// If a program has already received a value from t.C, the timer is known
// to have expired, and t.Reset can be used directly.
// If a program has not yet received a value from t.C, however,
// the timer must be stopped and—if Stop reports that the timer expired
// before being stopped—the channel explicitly drained:
//
// 	if !t.Stop() {
// 		<-t.C
// 	}
// 	t.Reset(d)
//
// This should not be done concurrent to other receives from the Timer's
// channel.
//
// Note that it is not possible to use Reset's return value correctly, as there
// is a race condition between draining the channel and the new timer expiring.
// Reset should always be invoked on stopped or expired channels, as described above.
// The return value exists to preserve compatibility with existing programs.
func (t *Timer) Reset(d Duration) bool {
	return t.earthTimer.Reset(vd2ed(d))
}

// AfterFunc waits for the duration to elapse and then calls f in its own
// goroutine. It returns a Timer that can be used to cancel the call using
// its Stop method.
func AfterFunc(d Duration, f func()) *Timer {
	return newTimer(d, func(c chan<- Time, t Time) {
		go f()
	})
}

// After waits for the duration to elapse and then sends the current time
// on the returned channel.
// It is equivalent to NewTimer(d).C.
// The underlying Timer is not recovered by the garbage collector
// until the timer fires. If efficiency is a concern, use NewTimer
// instead and call Timer.Stop if the timer is no longer needed.
func After(d Duration) <-chan Time {
	return NewTimer(d).C
}

// Sleep pauses the current goroutine for at least the duration d.
// A negative or zero duration causes Sleep to return immediately.
func Sleep(d Duration) {
	time.Sleep(vd2ed(d))
}
