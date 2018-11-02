package vanatime_test

import (
	"testing"

	vanatime "github.com/pasela/go-vanatime"
)

func TestMicroseconds(t *testing.T) {
	d := 100 * vanatime.Microsecond
	got := d.Microseconds()
	want := int64(100)
	if got != want {
		t.Fatalf("want %v, but %v:", want, got)
	}
}

func TestSeconds(t *testing.T) {
	d := 100*vanatime.Second + 50*vanatime.Microsecond
	got := d.Seconds()
	want := float64(100.000050)
	if got != want {
		t.Fatalf("want %v, but %v:", want, got)
	}
}

func TestMinutes(t *testing.T) {
	d := 3*vanatime.Minute + 30*vanatime.Second + 300*vanatime.Microsecond
	got := d.Minutes()
	want := float64(3.500005)
	if got != want {
		t.Fatalf("want %v, but %v:", want, got)
	}
}

func TestHours(t *testing.T) {
	d := 3*vanatime.Hour + 30*vanatime.Minute + 45*vanatime.Second + 900*vanatime.Microsecond
	got := d.Hours()
	want := float64(3.51250025)
	if got != want {
		t.Fatalf("want %v, but %v:", want, got)
	}
}

func TestTruncate(t *testing.T) {
	d := 3*vanatime.Hour + 30*vanatime.Minute + 45*vanatime.Second + 900*vanatime.Microsecond

	cases := [][]vanatime.Duration{
		{vanatime.Hour, 3 * vanatime.Hour},
		{vanatime.Minute, 3*vanatime.Hour + 30*vanatime.Minute},
		{20 * vanatime.Minute, 3*vanatime.Hour + 20*vanatime.Minute},
		{vanatime.Second, 3*vanatime.Hour + 30*vanatime.Minute + 45*vanatime.Second},
		{30 * vanatime.Second, 3*vanatime.Hour + 30*vanatime.Minute + 30*vanatime.Second},
	}
	for i, c := range cases {
		got := d.Truncate(c[0])
		want := c[1]
		if got != want {
			t.Errorf("[%d]: want %v, but %v:", i, want, got)
		}
	}
}

func TestRound(t *testing.T) {
	d := 3*vanatime.Hour + 30*vanatime.Minute + 45*vanatime.Second + 900*vanatime.Microsecond

	cases := [][]vanatime.Duration{
		{vanatime.Hour, 4 * vanatime.Hour},
		{3 * vanatime.Hour, 3 * vanatime.Hour},
		{30 * vanatime.Minute, 3*vanatime.Hour + 30*vanatime.Minute},
		{15 * vanatime.Minute, 3*vanatime.Hour + 30*vanatime.Minute},
		{20 * vanatime.Minute, 3*vanatime.Hour + 40*vanatime.Minute},
	}
	for i, c := range cases {
		got := d.Round(c[0])
		want := c[1]
		if got != want {
			t.Errorf("[%d]: want %v, but %v:", i, want, got)
		}
	}
}

func TestString(t *testing.T) {
	cases := []struct {
		D    vanatime.Duration
		Want string
	}{
		{vanatime.Hour, "1h0m0s"},
		{3 * vanatime.Hour, "3h0m0s"},
		{300 * vanatime.Hour, "300h0m0s"},
		{3*vanatime.Hour + 12*vanatime.Minute, "3h12m0s"},
		{3*vanatime.Hour + 72*vanatime.Minute, "4h12m0s"},
		{12*vanatime.Minute + 34*vanatime.Second, "12m34s"},
		{34*vanatime.Second + 56*vanatime.Millisecond, "34.056s"},
		{34*vanatime.Second + 56*vanatime.Microsecond, "34.000056s"},
		{3 * vanatime.Millisecond, "3ms"},
		{3 * vanatime.Microsecond, "3µs"},
	}
	for i, c := range cases {
		got := c.D.String()
		if got != c.Want {
			t.Errorf("[%d]: want %v, but %v:", i, c.Want, got)
		}
	}
}

func TestParseDuration(t *testing.T) {
	cases := []struct {
		S    string
		Want vanatime.Duration
	}{
		{"1h", vanatime.Hour},
		{"1h0m0s", vanatime.Hour},
		{"3h0m0s", 3 * vanatime.Hour},
		{"300h0m0s", 300 * vanatime.Hour},
		{"3h12m0s", 3*vanatime.Hour + 12*vanatime.Minute},
		{"3h72m0s", 4*vanatime.Hour + 12*vanatime.Minute},
		{"12m34s", 12*vanatime.Minute + 34*vanatime.Second},
		{"34.056s", 34*vanatime.Second + 56*vanatime.Millisecond},
		{"34.000056s", 34*vanatime.Second + 56*vanatime.Microsecond},
		{"3ms", 3 * vanatime.Millisecond},
		{"3µs", 3 * vanatime.Microsecond},
	}
	for i, c := range cases {
		got, err := vanatime.ParseDuration(c.S)
		if err != nil {
			t.Fatalf("[%d]: error %s", i, err)
		} else {
			if got != c.Want {
				t.Errorf("[%d]: want %v, but %v:", i, c.Want, got)
			}
		}
	}
}

func TestParseDurationError(t *testing.T) {
	cases := []string{
		"1ns",
	}
	for i, c := range cases {
		_, err := vanatime.ParseDuration(c)
		if err == nil {
			t.Fatalf("[%d]: want error, but nil", i)
		}
	}
}
