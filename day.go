package vanatime

import "golang.org/x/text/language"

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

var dayNames = map[language.Tag][8]string{
	language.English: [8]string{
		"Firesday",
		"Earthsday",
		"Watersday",
		"Windsday",
		"Iceday",
		"Lightningday",
		"Lightsday",
		"Darksday",
	},
	language.Japanese: [8]string{
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

var dayLangs language.Matcher

func init() {
	var keys []language.Tag
	for k, _ := range dayNames {
		keys = append(keys, k)
	}
	dayLangs = language.NewMatcher(keys)
}

// String returns the English name of the day ("Firesday", "Earthsday", ...).
func (w Weekday) String() string {
	return defaultDayNames[w]
}

// String returns the name of the day by specified locale.
func (w Weekday) StringLocale(locale string) string {
	userTag := language.Make(locale)
	tag, _, _ := dayLangs.Match(userTag)
	if names, ok := dayNames[tag]; ok {
		return names[w]
	}
	return w.String()
}
