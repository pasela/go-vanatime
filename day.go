package vanatime

// A Weekday specifies a day of the week in Vana'diel (Firesday = 0, ...).
type Weekday int

const (
	Firesday Weekday = iota
	Earthsday
	Watersday
	Windsday
	Iceday
	Lightningday
	Lightsday
	Darksday
)

var defaultDayNames = [...]string{
	"Firesday",
	"Earthsday",
	"Watersday",
	"Windsday",
	"Iceday",
	"Lightningday",
	"Lightsday",
	"Darksday",
}

var dayNames = map[string][8]string{
	"en": [8]string{
		"Firesday",
		"Earthsday",
		"Watersday",
		"Windsday",
		"Iceday",
		"Lightningday",
		"Lightsday",
		"Darksday",
	},
	"ja": [8]string{
		"火曜日",
		"土曜日",
		"水曜日",
		"風曜日",
		"氷曜日",
		"雷曜日",
		"光曜日",
		"闇曜日",
	},
}

func (w Weekday) String() string {
	return defaultDayNames[w]
}

func (w Weekday) StringLocale(locale string) string {
	if names, ok := dayNames[locale]; ok {
		return names[w]
	}
	return w.String()
}
