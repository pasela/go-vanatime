package vanatime

import (
	"errors"
	"strings"
	"time"
)

// A Duration represents the elapsed vana'diel time between two instants as an int64 microsecond count.
type Duration int64

const (
	minDuration Duration = -1 << 63
	maxDuration Duration = 1<<63 - 1
)

const (
	Microsecond Duration = 1
	Millisecond          = 1000 * Microsecond
	Second               = 1000 * Millisecond
	Minute               = 60 * Second
	Hour                 = 60 * Minute
	Day                  = 24 * Hour
	Week                 = 8 * Day
	Month                = 30 * Day
	Year                 = 360 * Day
)

func Since(t Time) Duration {
	return Now().Sub(t)
}

func Until(t Time) Duration {
	return t.Sub(Now())
}

func (d Duration) Microseconds() int64 {
	return int64(d)
}

func (d Duration) Seconds() float64 {
	sec := d / Second
	usec := d % Second
	return float64(sec) + float64(usec)/1e6
}

func (d Duration) Minutes() float64 {
	min := d / Minute
	usec := d % Minute
	return float64(min) + float64(usec)/(60*1e6)
}

func (d Duration) Hours() float64 {
	hour := d / Hour
	usec := d % Hour
	return float64(hour) + float64(usec)/(60*60*1e6)
}

func (d Duration) Truncate(m Duration) Duration {
	if m <= 0 {
		return d
	}
	return d - d%m
}

func (d Duration) Round(m Duration) Duration {
	t := time.Duration(d * 1000).Round(time.Duration(m * 1000))
	return Duration(t / 1000)
}

func (d Duration) String() string {
	return time.Duration(d * 1000).String()
}

func ParseDuration(s string) (Duration, error) {
	if strings.Contains(s, "ns") {
		return 0, errors.New("vanatime: unknown unit ns in duration " + s)
	}

	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, err
	}
	return Duration(d / 1000), nil
}
