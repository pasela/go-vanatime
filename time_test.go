package vanatime_test

import (
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
