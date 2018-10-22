package vanatime

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var defaultFormatPadding = map[rune]string{
	'e': " ", 'k': " ", 'A': " ", 'n': " ", 't': " ", '%': " ",
}

func formatPadding(r rune) string {
	if pad, ok := defaultFormatPadding[r]; ok {
		return pad
	}
	return "0"
}

var defaultFormatWidth = map[rune]int{
	'y': 2, 'm': 2, 'd': 2, 'e': 2, 'H': 2, 'k': 2, 'M': 2, 'S': 2,
	'j': 3, 'L': 3,
	'N': 6,
}

func formatWidth(r rune) int {
	if width, ok := defaultFormatWidth[r]; ok {
		return width
	}
	return 0
}

// Strftime formats Vana'diel time according to the directives in the format string.
// The directives begins with a percent (%) character. Any text not listed
// as a directive will be passed through to the output string.
//
// The directive consists of a percent (%) character, zero or more flags,
// optional minimum field width and a conversion specifier as follows.
//
//     %<flags><width><conversion>
//
// Flags:
//
//     -  don't pad a numerical output.
//     _  use spaces for padding.
//     0  use zeros for padding.
//     ^  upcase the result string.
//     #  change case.
//
// The minimum field width specifies the minimum width.
//
// Format directives:
//
//     Date (Year, Month, Day):
//       %Y - Year with century (can be negative)
//               -0001, 0000, 1995, 2009, 14292, etc.
//       %C - year / 100 (round down.  20 in 2009)
//       %y - year % 100 (00..99)
//
//       %m - Month of the year, zero-padded (01..12)
//               %_m  blank-padded ( 1..12)
//               %-m  no-padded (1..12)
//
//       %d - Day of the month, zero-padded (01..30)
//               %-d  no-padded (1..30)
//       %e - Day of the month, blank-padded ( 1..30)
//
//       %j - Day of the year (001..360)
//
//     Time (Hour, Minute, Second, Subsecond):
//       %H - Hour of the day, 24-hour clock, zero-padded (00..23)
//       %k - Hour of the day, 24-hour clock, blank-padded ( 0..23)
//
//       %M - Minute of the hour (00..59)
//
//       %S - Second of the minute (00..59)
//
//       %L - Millisecond of the second (000..999)
//       %N - Fractional seconds digits, default is 6 digits (microsecond)
//               %3N  millisecond (3 digits)
//               %6N  microsecond (6 digits)
//
//     Weekday:
//       %A - The full weekday name (``Firesday'')
//               %^A  uppercased (``FIRESDAY'')
//       %w - Day of the week (Firesday is 0, 0..7)
//
//     Seconds since the Epoch:
//       %s - Number of seconds since 0001-01-01 00:00:00
//
//     Literal string:
//       %n - Newline character (\n)
//       %t - Tab character (\t)
//       %% - Literal ``%'' character
//
//     Combination:
//       %F - The ISO 8601 date format (%Y-%m-%d)
//       %X - Same as %T
//       %R - 24-hour time (%H:%M)
//       %T - 24-hour time (%H:%M:%S)
func (t Time) Strftime(format string) string {
	year, mon, day, yday := t.Date()
	hour, min, sec := t.Clock()
	usec := t.Microsecond()
	wday := t.Weekday()

	source := map[rune]int64{
		'Y': int64(year), 'C': int64(year / 100), 'y': int64(year % 100),
		'm': int64(mon), 'd': int64(day), 'e': int64(day), 'j': int64(yday),
		'H': int64(hour), 'k': int64(hour), 'M': int64(min), 'S': int64(sec),
		'L': int64(usec), 'N': int64(usec),
		'A': int64(wday), 'w': int64(wday), 's': t.time,
	}

	format = normalizeFormat(format)
	re := regexp.MustCompile(`%([-_0^#]+)?(\d+)?([YCymdejHkMSLNAawsnt%])`)
	result := re.ReplaceAllStringFunc(format, func(substr string) string {
		parts := re.FindStringSubmatch(substr)
		flags := parts[1]
		var width int
		if parts[2] != "" {
			width, _ = strconv.Atoi(parts[2])
		}
		conversion := []rune(parts[3])[0]
		upcase := false
		padding := formatPadding(conversion)
		if width == 0 {
			width = formatWidth(conversion)
		}

		if flags != "" {
			for _, c := range flags {
				switch c {
				case '-':
					padding = ""
				case '_':
					padding = " "
				case '0':
					padding = "0"
				case '^', '#':
					upcase = true
				}
			}
		}

		var value string
		switch conversion {
		case 'L', 'N':
			var v int
			if width <= 6 {
				v = usec / (100000 / int(math.Pow10(width-1)))
			} else {
				v = usec * int(math.Pow10(width-6))
			}
			value = strconv.Itoa(v)

		case 'A':
			value = wday.String()

		case 'n':
			value = "\n"
		case 't':
			value = "\t"
		case '%':
			value = "%"

		default:
			if v, ok := source[conversion]; ok {
				value = strconv.FormatInt(v, 10)
			} else {
				value = ""
			}
		}

		length := utf8.RuneCountInString(value)
		if width > 0 && padding != "" && length < width {
			value = strings.Repeat(padding, width-length) + value
		}

		if upcase {
			value = strings.ToUpper(value)
		}

		return value
	})

	return result
}

func normalizeFormat(format string) string {
	re := regexp.MustCompile(`%([-_0^#]+)?(\d+)?([FXRT])`)
	return re.ReplaceAllStringFunc(format, func(substr string) string {
		parts := re.FindStringSubmatch(substr)
		switch parts[3] {
		case "F":
			return "%Y-%m-%d"
		case "T", "X":
			return "%H:%M:%S"
		case "R":
			return "%H:%M"
		default:
			return substr
		}
	})
}
