package vanatime

import "time"

// vanatime is an abstraction of Vana'diel dates and times from Final Fantasy XI.
// Time is stored internally as the number of microseconds since C.E. 0001-01-01 00:00:00.
//
// Vana'diel time spec:
//
//     One year   = 12 months = 360 days
//     One month  = 30 days
//     One day    = 24 hours
//     One hour   = 60 minutes
//     One minute = 60 seconds
//     One second = 0.04 seconds of the earth's (1/25th of a second)
//
//     Vana'diel second         = 0.04 earth seconds (1/25th of a second)
//     Vana'diel minute         = 2.4 earth seconds
//     Vana'diel hour           = 2 minutes 24 earth seconds
//     Vana'diel day            = 57 minutes 36 earth seconds
//     Vana'diel week           = 7 hours 40 minutes 48 earth seconds
//     Vana'diel calendar month = 1 day 4 hours 48 earth minutes
//     Vana'diel lunar month    = 3 days 14 hours 24 earth minutes
//     Vana'diel year           = 14 days 9 hours 36 earth minutes
//
//     Each full lunar cycle lasts for 84 Vana'diel days.
//     Vana'diel has 12 distinct moon phases.
//     Japanese client expresses moon phases by 12 kinds of texts. (percentage is not displayed in Japanese client)
//     Non-Japanese client expresses moon phases by 7 kinds of texts and percentage.
//
// C.E. = Crystal Era
//
//     A.D. -91270800 => 1967/02/10 00:00:00 +0900
//     C.E. 0         => 0001/01/01 00:00:00
//
//     A.D. 2002/01/01(Tue) 00:00:00 JST
//     C.E. 0886/01/01(Fir) 00:00:00
//
//     A.D. 2047/10/22(Tue) 01:00:00 JST
//     C.E. 2047/10/22(Wat) 01:00:00
//
//     A.D. 2047/10/21(Mon) 15:37:30 UTC
//     C.E. 2047/10/21(Win) 15:37:30

const (
	TimeScale         int   = 25 // Vana'diel time goes 25 times faster than the Earth
	BaseYear          int   = 886
	BaseTime          int64 = (int64(BaseYear) * int64(Year)) / int64(TimeScale)
	EarthBaseTime     int64 = 1009810800 * int64(Second) // 2002-01-01 00:00:00.000 JST
	VanaEarthDiffTime int64 = BaseTime - EarthBaseTime
	MoonCycleDays     int   = 84 // Vana'diel moon cycle lasts 84 days
)

// A Time represents an instant in Vana'diel time with microsecond precision.
type Time struct {
	// the time as microseconds since C.E. 0001-01-01 00:00:00
	time int64
}

// Now returns the current Vana'diel time.
func Now() Time {
	return earth2vana(time.Now())
}

// Date returns the Time corresponding to given arguments.
func Date(year, mon, day, hour, min, sec, usec int) Time {
	sec, usec = norm(sec, usec, 1e6)
	min, sec = norm(min, sec, 60)
	hour, min = norm(hour, min, 60)
	day, hour = norm(day, hour, 24)
	mon, day = norm(mon, day, 30)
	year, mon = norm(year, mon, 12)

	return Time{
		(int64((year - 1)) * int64(Year)) +
			(int64((mon - 1)) * int64(Month)) +
			(int64((day - 1)) * int64(Day)) +
			(int64(hour) * int64(Hour)) +
			(int64(min) * int64(Minute)) +
			(int64(sec) * int64(Second)) +
			int64(usec),
	}
}

// FromEarth returns the Time corresponding to the given Earth time.
func FromEarth(earth time.Time) Time {
	return earth2vana(earth)
}

// FromInt64 returns the Time corresponding to the given Vana'diel time (since C.E. 0001-01-01 00:00:00).
func FromInt64(time int64) Time {
	return Time{
		time: time,
	}
}

// IsZero reports whether t represents the zero time instant.
func (t Time) IsZero() bool {
	return t.time == 0
}

// Before reports whether the time instant t is before u.
func (t Time) Before(u Time) bool {
	return t.time < u.time
}

// After reports whether the time instant t is after u.
func (t Time) After(u Time) bool {
	return t.time > u.time
}

// Equal reports whether t and u represent the same time instant.
func (t Time) Equal(u Time) bool {
	return t.time == u.time
}

// Add returns the time t+d.
func (t Time) Add(d Duration) Time {
	return Time{t.time + int64(d)}
}

// AddDate returns the time corresponding to adding the given number of years, months and days to t.
func (t Time) AddDate(years int, months int, days int) Time {
	year, month, day, _ := t.Date()
	hour, min, sec := t.Clock()
	return Date(year+years, month+months, day+days, hour, min, sec, t.Microsecond())
}

func lessThanHalf(x, y Duration) bool {
	return uint64(x)+uint64(x) < uint64(y)
}

// Truncate returns the result of rounding t down to a multiple of d.
//
// Truncate operates on the time as an absolute duration since the zero
// time; it does not operate on the presentation form of the time. Thus,
// Truncate(Hour) may return a time with a non-zero minute.
func (t Time) Truncate(d Duration) Time {
	if d <= 0 {
		return t
	}
	r := Duration(t.time % int64(d))
	return t.Add(-r)
}

// Round returns the result of rounding t to the nearest multiple of d.
// The rounding behavior for halfway values is to round up.
//
// Round operates on the time as an absolute duration since the zero
// time; it does not operate on the presentation form of the time. Thus,
// Round(Hour) may return a time with a non-zero minute.
func (t Time) Round(d Duration) Time {
	if d <= 0 {
		return t
	}
	r := Duration(t.time % int64(d))
	if lessThanHalf(r, d) {
		return t.Add(-r)
	}
	return t.Add(d - r)
}

// Sub returns the duration t-u. If the result exceeds the maximum (or minimum)
// value that can be stored in a Duration, the maximum (or minimum) duration
// will be returned. To compute t-d for a duration d, use t.Add(-d).
func (t Time) Sub(u Time) Duration {
	d := Duration(t.time - u.time)
	switch {
	case u.Add(d).Equal(t):
		return d
	case t.Before(u):
		return minDuration
	default:
		return maxDuration
	}
}

// Earth returns the time of Earth.
func (t Time) Earth() time.Time {
	return vana2earth(t)
}

// Date returns the year, month, day and day of the year in which t occurs.
func (t Time) Date() (year, mon, day, yday int) {
	year = int(t.time/int64(Year)) + 1
	mon = int(t.time%int64(Year)/int64(Month)) + 1
	day = int(t.time%int64(Month)/int64(Day)) + 1
	yday = (mon-1)*30 + day
	return
}

// Year returns the year in which t occurs.
func (t Time) Year() int {
	year, _, _, _ := t.Date()
	return year
}

// Month returns the month of the year specified by t.
func (t Time) Month() int {
	_, mon, _, _ := t.Date()
	return mon
}

// Day returns the day of the month specified by t.
func (t Time) Day() int {
	_, _, day, _ := t.Date()
	return day
}

// YearDay returns the day of the year specified by t, in the range [1,360].
func (t Time) YearDay() int {
	_, _, _, yday := t.Date()
	return yday
}

// Weekday returns the day of the week specified by t.
func (t Time) Weekday() Weekday {
	wday := int(t.time % int64(Week) / int64(Day))
	return Weekday(wday)
}

// Clock returns the hour, minute, and second within the day specified by t.
func (t Time) Clock() (hour, min, sec int) {
	hour = int(t.time % int64(Day) / int64(Hour))
	min = int(t.time % int64(Hour) / int64(Minute))
	sec = int(t.time % int64(Minute) / int64(Second))
	return
}

// Hour returns the hour within the day specified by t, in the range [0, 23].
func (t Time) Hour() int {
	hour, _, _ := t.Clock()
	return hour
}

// Minute returns the minute offset within the hour specified by t, in the range [0, 59].
func (t Time) Minute() int {
	_, min, _ := t.Clock()
	return min
}

// Second returns the second offset within the minute specified by t, in the range [0, 59].
func (t Time) Second() int {
	_, _, sec := t.Clock()
	return sec
}

// Microsecond returns the microsecond offset within the second specified by t, in the range [0, 999999].
func (t Time) Microsecond() int {
	return int(t.time % int64(Second))
}

// Int64 returns t as a int64 since C.E. 0001-01-01 00:00:00.
func (t Time) Int64() int64 {
	return t.time
}

// Moon returns the moon specified by t.
func (t Time) Moon() Moon {
	var days int = int(t.time / int64(Day))
	timeOfMoon := (((int64(days) + 12) % 7) * int64(Day)) + (t.time % int64(Day))

	return Moon{
		days:       days,
		timeOfMoon: timeOfMoon,
	}
}

// String returns the time formatted using the format string.
//	"%Y-%m-%d %H:%M:%S"
func (t Time) String() string {
	m := t.Moon()
	return t.Strftime("%Y-%m-%d %H:%M:%S") + " " + t.Weekday().String() + " " + m.String()
}

func earth2vana(etime time.Time) Time {
	return Time{e2v(etime.UnixNano() / 1000)}
}

func e2v(etime int64) int64 {
	return (etime+VanaEarthDiffTime)*int64(TimeScale) - int64(Year)
}

func vana2earth(vtime Time) time.Time {
	usec := v2e(vtime.time)
	return time.Unix(usec/int64(Second), usec%int64(Second))
}

func v2e(vtime int64) int64 {
	return ((vtime + int64(Year)) / int64(TimeScale)) - VanaEarthDiffTime
}

// from https://golang.org/src/time/time.go
func norm(hi, lo, base int) (nhi, nlo int) {
	if lo < 0 {
		n := (-lo-1)/base + 1
		hi -= n
		lo += n * base
	}
	if lo >= base {
		n := lo / base
		hi += n
		lo -= n * base
	}
	return hi, lo
}
