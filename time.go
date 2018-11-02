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

func Now() Time {
	return earth2vana(time.Now())
}

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

func FromEarth(earth time.Time) Time {
	return earth2vana(earth)
}

func FromInt64(time int64) Time {
	return Time{
		time: time,
	}
}

func (t Time) IsZero() bool {
	return t.time == 0
}

func (t Time) Before(u Time) bool {
	return t.time < u.time
}

func (t Time) After(u Time) bool {
	return t.time > u.time
}

func (t Time) Equal(u Time) bool {
	return t.time == u.time
}

func (t Time) Add(d Duration) Time {
	return Time{t.time + int64(d)}
}

func (t Time) AddDate(years int, months int, days int) Time {
	year, month, day, _ := t.Date()
	hour, min, sec := t.Clock()
	return Date(year+years, month+months, day+days, hour, min, sec, t.Microsecond())
}

func (t Time) Earth() time.Time {
	return vana2earth(t)
}

func (t Time) Date() (year, mon, day, yday int) {
	year = int(t.time/int64(Year)) + 1
	mon = int(t.time%int64(Year)/int64(Month)) + 1
	day = int(t.time%int64(Month)/int64(Day)) + 1
	yday = (mon-1)*30 + day
	return
}

func (t Time) Year() int {
	year, _, _, _ := t.Date()
	return year
}

func (t Time) Month() int {
	_, mon, _, _ := t.Date()
	return mon
}

func (t Time) Day() int {
	_, _, day, _ := t.Date()
	return day
}

func (t Time) YearDay() int {
	_, _, _, yday := t.Date()
	return yday
}

func (t Time) Weekday() Weekday {
	wday := int(t.time % int64(Week) / int64(Day))
	return Weekday(wday)
}

func (t Time) Clock() (hour, min, sec int) {
	hour = int(t.time % int64(Day) / int64(Hour))
	min = int(t.time % int64(Hour) / int64(Minute))
	sec = int(t.time % int64(Minute) / int64(Second))
	return
}

func (t Time) Hour() int {
	hour, _, _ := t.Clock()
	return hour
}

func (t Time) Minute() int {
	_, min, _ := t.Clock()
	return min
}

func (t Time) Second() int {
	_, _, sec := t.Clock()
	return sec
}

func (t Time) Microsecond() int {
	return int(t.time % int64(Second))
}

func (t Time) Int64() int64 {
	return t.time
}

func (t Time) Moon() Moon {
	var days int = int(t.time / int64(Day))
	timeOfMoon := (((int64(days) + 12) % 7) * int64(Day)) + (t.time % int64(Day))

	return Moon{
		days:       days,
		timeOfMoon: timeOfMoon,
	}
}

func (t Time) String() string {
	return t.Strftime("%Y-%m-%d %H:%M:%S")
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
