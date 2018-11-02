package vanatime

import "math"

// MOON_BASE_TIME  = 0 - (ONE_DAY * 12) # Start of New moon (10%)
//
// moon_time     = @time - MOON_BASE_TIME
// @moon_age     = (moon_time / ONE_DAY / 7 % (MAX_MOON_AGE + 1)).floor
// @time_of_moon = ((moon_time / ONE_DAY % 7) * ONE_DAY).floor
//               + (@hour * ONE_HOUR)
//               + (@min * ONE_MINUTE)
//               + (@sec * ONE_SECOND)
//               + (@usec)

// C.E. 0001/01/01 00:00:00 => WXC 19%
// C.E. 0886/01/01 00:00:00 => NM  10%
//
// 0% NM   7% WXC  40% FQM  57% WXG   90% FM  93% WNG  60% LQM  43% WNC  10% NM
// 2% NM  10% WXC  43% FQM  60% WXG   93% FM  90% WNG  57% LQM  40% WNC   7% NM
// 5% NM  12% WXC  45% FQM  62% WXG   95% FM  88% WNG  55% LQM  38% WNC   5% NM
//        14% WXC  48% FQM  64% WXG   98% FM  86% WNG  52% LQM  36% WNC   2% NM
//        17% WXC  50% FQM  67% WXG  100% FM  83% WNG  50% LQM  33% WNC
//        19% WXC  52% FQM  69% WXG   98% FM  81% WNG  48% LQM  31% WNC
//        21% WXC  55% FQM  71% WXG   95% FM  79% WNG  45% LQM  29% WNC
//        24% WXC           74% WXG           76% WNG           26% WNC
//        26% WXC           76% WXG           74% WNG           24% WNC
//        29% WXC           79% WXG           71% WNG           21% WNC
//        31% WXC           81% WXG           69% WNG           19% WNC
//        33% WXC           83% WXG           67% WNG           17% WNC
//        36% WXC           86% WXG           64% WNG           14% WNC
//        38% WXC           88% WXG           62% WNG           12% WNC

type MoonPhase int

const (
	// 新月
	NewMoon MoonPhase = iota

	// 三日月
	WaxingCrescent1

	// 七日月
	WaxingCrescent2

	// 上弦の月
	FirstQuarter

	// 十日夜
	WaxingGibbous1

	// 十三夜
	WaxingGibbous2

	// 満月
	FullMoon

	// 十六夜
	WaningGibbous1

	// 居待月
	WaningGibbous2

	// 下弦の月
	LastQuarter

	// 二十日余月
	WaningCrescent1

	// 二十六夜
	WaningCrescent2
)

var defaultMoonNames = [...]string{
	"New Moon",
	"Waxing Crescent",
	"Waxing Crescent",
	"First Quarter",
	"Waxing Gibbous",
	"Waxing Gibbous",
	"Full Moon",
	"Waning Gibbous",
	"Waning Gibbous",
	"Last Quarter",
	"Waning Crescent",
	"Waning Crescent",
}

var moonNames = map[string][12]string{
	"en": [12]string{
		"New Moon",
		"Waxing Crescent",
		"Waxing Crescent",
		"First Quarter",
		"Waxing Gibbous",
		"Waxing Gibbous",
		"Full Moon",
		"Waning Gibbous",
		"Waning Gibbous",
		"Last Quarter",
		"Waning Crescent",
		"Waning Crescent",
	},
	"ja": [12]string{
		"新月",
		"三日月",
		"七日月",
		"上弦の月",
		"十日夜",
		"十三夜",
		"満月",
		"十六夜",
		"居待月",
		"下弦の月",
		"二十日余月",
		"二十六夜",
	},
}

func (m MoonPhase) String() string {
	return defaultMoonNames[m]
}

func (m MoonPhase) StringLocale(locale string) string {
	if names, ok := moonNames[locale]; ok {
		return names[m]
	}
	return m.String()
}

type Moon struct {
	days       int
	timeOfMoon int64
}

func (m Moon) Percent() int {
	percent := math.Round(float64((m.days+8)%MoonCycleDays) * (200.0 / float64(MoonCycleDays)))
	if percent > 100.0 {
		percent = 200.0 - percent
	}
	return int(percent)
}

func (m Moon) Phase() MoonPhase {
	return MoonPhase(((m.days + 12) / 7) % 12)
}

func (m Moon) TimeOfMoon() int64 {
	return m.timeOfMoon
}
