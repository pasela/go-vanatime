# go-vanatime

Go library for dealing with Vana'diel time from Final Fantasy XI.
Converting between realtime and Vana'diel time, and so on.

## Examples

See `_example` directory.

```go
// Current Vana'diel time
vt := vanatime.Now()
fmt.Println(vt)
//=> 1313-04-13 21:20:27 Lightsday Waxing Crescent (33%)

// From the earth time
et := time.Now()
vt = vanatime.FromEarth(et)
fmt.Println(et)
fmt.Println(vt)
//=> 2018-11-05 21:58:25.119486 +0900 JST m=+0.000910642
//=> 1313-04-13 21:20:27 Lightsday Waxing Crescent (33%)

// Specified Vana'diel time
vt = vanatime.Date(1300, 2, 3, 0, 0, 0, 0)
fmt.Println(vt)
//=> 1300-02-03 00:00:00 Firesday Waning Gibbous (76%)

// Weekday and MoonPhase in Japanese
fmt.Println(
    vt.Weekday().StringLocale("ja"),
    vt.Moon().Phase().StringLocale("ja"),
)
//=> 火曜日 居待月

// To the earth time
et = vt.Earth()
fmt.Println(et)
//=> 2018-04-29 21:07:12 +0900 JST

// Add 3 days
vt = vt.Add(3 * vanatime.Day)
fmt.Println(vt)
//=> 1300-02-06 00:00:00 Windsday Waning Gibbous (69%)

// Formatting
fmt.Println(vt.Strftime("%Y/%m/%d %H:%M:%S"))
//=> 1300/02/06 00:00:00
```

## License

MIT

## Author

pasela
