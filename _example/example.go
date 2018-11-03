package main

import (
	"fmt"
	"time"

	vanatime "github.com/pasela/go-vanatime"
)

func main() {
	// Current Vana'diel time
	vt := vanatime.Now()
	fmt.Println(vt)

	// From the earth time
	et := time.Now()
	vt = vanatime.FromEarth(et)
	fmt.Println(et)
	fmt.Println(vt)

	// Specified Vana'diel time
	vt = vanatime.Date(1300, 2, 3, 0, 0, 0, 0)
	fmt.Println(vt)

	// Weekday and MoonPhase in Japanese
	fmt.Println(
		vt.Weekday().StringLocale("ja"),
		vt.Moon().Phase().StringLocale("ja"),
	)

	// to the earth time
	et = vt.Earth()
	fmt.Println(et)
}
