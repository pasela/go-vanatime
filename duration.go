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

// ParseDuration parses a duration string. A duration string is a possibly
// signed sequence of decimal numbers, each with optional fraction and a unit
// suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are
// "us" (or "µs"), "ms", "s", "m", "h".
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

// Microseconds returns the duration as an integer microsecond count.
func (d Duration) Microseconds() int64 {
	return int64(d)
}

// Seconds returns the duration as a floating point number of seconds.
func (d Duration) Seconds() float64 {
	sec := d / Second
	usec := d % Second
	return float64(sec) + float64(usec)/1e6
}

// Minutes returns the duration as a floating point number of minutes.
func (d Duration) Minutes() float64 {
	min := d / Minute
	usec := d % Minute
	return float64(min) + float64(usec)/(60*1e6)
}

// Hours returns the duration as a floating point number of hours.
func (d Duration) Hours() float64 {
	hour := d / Hour
	usec := d % Hour
	return float64(hour) + float64(usec)/(60*60*1e6)
}

// Truncate returns the result of rounding d toward zero to a multiple of m.
// If m <= 0, Truncate returns d unchanged.
func (d Duration) Truncate(m Duration) Duration {
	if m <= 0 {
		return d
	}
	return d - d%m
}

// Round returns the result of rounding d to the nearest multiple of m.
// The rounding behavior for halfway values is to round away from zero.
// If the result exceeds the maximum (or minimum) value that can be stored
// in a Duration, Round returns the maximum (or minimum) duration.
// If m <= 0, Round returns d unchanged.
func (d Duration) Round(m Duration) Duration {
	t := time.Duration(d * 1000).Round(time.Duration(m * 1000))
	return Duration(t / 1000)
}

// String returns a string representing the duration in the form "72h3m0.5s".
// Leading zero units are omitted. As a special case, durations less than one
// second format use a smaller unit (milli-, microseconds) to ensure that the
// leading digit is non-zero. The zero duration formats as 0s.
func (d Duration) String() string {
	return time.Duration(d * 1000).String()
}

// Since returns the time elapsed since t. It is shorthand for time.Now().Sub(t).
func Since(t Time) Duration {
	return Now().Sub(t)
}

// Until returns the duration until t. It is shorthand for t.Sub(time.Now()).
func Until(t Time) Duration {
	return t.Sub(Now())
}

func vd2ed(d Duration) time.Duration {
	return time.Duration(int64(d*1000) / int64(TimeScale))
}
