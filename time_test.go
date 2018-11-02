package vanatime_test

import (
	"fmt"
	"testing"
	"time"

	vanatime "github.com/pasela/go-vanatime"
)

var locJA *time.Location

func init() {
	ja, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	locJA = ja
}

func TestEpoch(t *testing.T) {
	vt := vanatime.Date(1, 1, 1, 0, 0, 0, 0)
	got := vt.Earth()
	want := time.Date(1967, 2, 10, 0, 0, 0, 0, locJA)

	if !got.Equal(want) {
		t.Fatalf("want %v, but %v:", want, got)
	}
}

func TestFromEarth(t *testing.T) {
	et := time.Date(1967, 2, 10, 0, 0, 0, 0, locJA)
	got := vanatime.FromEarth(et)
	want := vanatime.Date(1, 1, 1, 0, 0, 0, 0)

	if !got.Equal(want) {
		t.Fatalf("want %v, but %v:", want, got)
	}
}

func formatVanatime(vt vanatime.Time) string {
	moon := vt.Moon()
	return fmt.Sprintf(
		"%s %s %s (%d%%)",
		vt.Strftime("%Y-%m-%d %H:%M:%S %A"),
		moon.Phase().StringLocale("ja"),
		moon.Phase().String(),
		moon.Percent(),
	)
}

func TestEarthVariation(t *testing.T) {
	patterns := [][2]string{
		{"2018-11-01 00:00:00", "1312-12-11 00:00:00 Iceday 下弦の月 Last Quarter (57%)"},
		{"2018-11-02 00:00:00", "1313-01-06 00:00:00 Lightningday 新月 New Moon (2%)"},
		{"2018-11-03 00:00:00", "1313-02-01 00:00:00 Lightsday 十日夜 Waxing Gibbous (62%)"},
	}

	layout := "2006-01-02 15:04:05"
	for i, pattern := range patterns {
		et, err := time.ParseInLocation(layout, pattern[0], locJA)
		if err != nil {
			t.Fatalf("[%d]: %s", i, err)
		}
		vt := vanatime.FromEarth(et)
		vs := formatVanatime(vt)

		if vs != pattern[1] {
			t.Errorf(`[%d]: want "%s", but "%s"`, i, pattern[1], vs)
		}
	}
}

func TestVanaVariation(t *testing.T) {
	patterns := []struct {
		V vanatime.Time
		E string
	}{
		{vanatime.Date(1000, 3, 1, 0, 0, 0, 0), "2006-07-03 00:00:00.000000 Monday"},
		{vanatime.Date(1000, 3, 2, 0, 0, 0, 0), "2006-07-03 00:57:36.000000 Monday"},
		{vanatime.Date(1000, 3, 3, 0, 0, 0, 0), "2006-07-03 01:55:12.000000 Monday"},
	}

	layout := "2006-01-02 15:04:05.000000 Monday"
	for i, pattern := range patterns {
		et := pattern.V.Earth()
		want, err := time.ParseInLocation(layout, pattern.E, locJA)
		if err != nil {
			t.Fatalf("[%d]: %s", i, err)
		}

		if !et.Equal(want) {
			t.Errorf(`[%d]: want "%s", but "%s"`, i, want, et)
		}
	}
}

func TestAdd(t *testing.T) {
	vt := vanatime.Date(1000, 3, 1, 0, 0, 0, 0)
	want := vanatime.Date(1000, 3, 8, 12, 34, 56, 0)
	got := vt.Add(7*vanatime.Day +
		12*vanatime.Hour +
		34*vanatime.Minute +
		56*vanatime.Second)

	if !got.Equal(want) {
		t.Errorf(`want "%v", but "%v"`, want, got)
	}
}

func TestAddDate(t *testing.T) {
	vt := vanatime.Date(1000, 3, 1, 0, 0, 0, 0)
	want := vanatime.Date(1001, 5, 4, 0, 0, 0, 0)
	got := vt.AddDate(1, 2, 3)

	if !got.Equal(want) {
		t.Errorf(`want "%v", but "%v"`, want, got)
	}
}

func TestTruncate(t *testing.T) {
	vt := vanatime.Date(650, 3, 11, 12, 34, 56, 0)
	patterns := []struct {
		D    vanatime.Duration
		Want vanatime.Time
	}{
		{vanatime.Hour, vanatime.Date(650, 3, 11, 12, 0, 0, 0)},
		{2 * vanatime.Hour, vanatime.Date(650, 3, 11, 12, 0, 0, 0)},
		{5 * vanatime.Hour, vanatime.Date(650, 3, 11, 10, 0, 0, 0)},
		{30 * vanatime.Minute, vanatime.Date(650, 3, 11, 12, 30, 0, 0)},
	}

	for i, pattern := range patterns {
		got := vt.Truncate(pattern.D)

		if !got.Equal(pattern.Want) {
			t.Errorf(`[%d]: want "%s", but "%s"`, i, pattern.Want, got)
		}
	}
}

func TestRound(t *testing.T) {
	vt := vanatime.Date(650, 3, 11, 12, 34, 56, 0)
	patterns := []struct {
		D    vanatime.Duration
		Want vanatime.Time
	}{
		{vanatime.Hour, vanatime.Date(650, 3, 11, 13, 0, 0, 0)},
		{2 * vanatime.Hour, vanatime.Date(650, 3, 11, 12, 0, 0, 0)},
		{5 * vanatime.Hour, vanatime.Date(650, 3, 11, 15, 0, 0, 0)},
		{30 * vanatime.Minute, vanatime.Date(650, 3, 11, 12, 30, 0, 0)},
		{20 * vanatime.Minute, vanatime.Date(650, 3, 11, 12, 40, 0, 0)},
		{15 * vanatime.Minute, vanatime.Date(650, 3, 11, 12, 30, 0, 0)},
		{10 * vanatime.Minute, vanatime.Date(650, 3, 11, 12, 30, 0, 0)},
		{10 * vanatime.Second, vanatime.Date(650, 3, 11, 12, 35, 0, 0)},
	}

	for i, pattern := range patterns {
		got := vt.Round(pattern.D)

		if !got.Equal(pattern.Want) {
			t.Errorf(`[%d]: want "%s", but "%s"`, i, pattern.Want, got)
		}
	}
}
