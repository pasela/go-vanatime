package vanatime_test

import (
	"runtime"
	"sync/atomic"
	"testing"

	. "github.com/pasela/go-vanatime"
)

// Test the basic function calling behavior. Correct queueing
// behavior is tested elsewhere, since After and AfterFunc share
// the same code.
func TestAfterFunc(t *testing.T) {
	i := 10
	c := make(chan bool)
	var f func()
	f = func() {
		i--
		if i >= 0 {
			AfterFunc(0, f)
			Sleep(1 * Second)
		} else {
			c <- true
		}
	}

	AfterFunc(0, f)
	<-c
}

func TestAfterStress(t *testing.T) {
	stop := uint32(0)
	go func() {
		for atomic.LoadUint32(&stop) == 0 {
			runtime.GC()
			// Yield so that the OS can wake up the timer thread,
			// so that it can generate channel sends for the main goroutine,
			// which will eventually set stop = 1 for us.
			Sleep(Millisecond)
		}
	}()
	ticker := NewTicker(1)
	for i := 0; i < 100; i++ {
		<-ticker.C
	}
	ticker.Stop()
	atomic.StoreUint32(&stop, 1)
}
