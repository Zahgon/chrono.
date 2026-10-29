package jsdate

import (
	"math"
	"testing"
	"time"
)

func inZone(t *testing.T, name string) {
	t.Helper()
	loc, err := time.LoadLocation(name)
	if err != nil {
		t.Fatal(err)
	}
	prev := SetLocation(loc)
	t.Cleanup(func() { SetLocation(prev) })
}

func iso(t *testing.T, d *Date) string {
	t.Helper()
	s, ok := d.ToISOString()
	if !ok {
		return "Invalid"
	}
	return s
}

// Values recorded with node 24 (V8, ICU 78.3, tz 2026e).
func TestNewYorkExamples(t *testing.T) {
	inZone(t, "America/New_York")
	// Local mean time -4:56:02: exact in both directions, truncated by getTimezoneOffset.
	d := NewLocal(1800, 0, 1)
	if got := iso(t, d); got != "1800-01-01T04:56:02.000Z" {
		t.Errorf("new Date(1800,0,1) = %s", got)
	}
	if got := d.GetTimezoneOffset(); got != 296 {
		t.Errorf("getTimezoneOffset = %v", got)
	}
	lmt := FromTime(-5364662400000)
	if h, m, s := lmt.GetHours(), lmt.GetMinutes(), lmt.GetSeconds(); h != 19 || m != 3 || s != 58 {
		t.Errorf("LMT wall time %v:%v:%v", h, m, s)
	}
	// Gap: 02:30 does not exist, read with the offset before (EST), so it moves forward.
	if got := iso(t, NewLocal(2021, 2, 14, 2, 30)); got != "2021-03-14T07:30:00.000Z" {
		t.Errorf("gap = %s", got)
	}
	// Overlap: 01:30 occurs twice, the earlier instant (EDT) wins.
	if got := iso(t, NewLocal(2021, 10, 7, 1, 30)); got != "2021-11-07T05:30:00.000Z" {
		t.Errorf("overlap = %s", got)
	}
	inv := FromTime(math.NaN())
	if got := inv.SetHours(1); !math.IsNaN(got) {
		t.Errorf("setHours on invalid = %v", got)
	}
	if got := inv.SetFullYear(2020); got != 1577854800000 {
		t.Errorf("setFullYear on invalid = %v", got)
	}
	if got := iso(t, NewLocal(99, 0)); got != "1999-01-01T05:00:00.000Z" {
		t.Errorf("new Date(99, 0) = %s", got)
	}
	if got := UTC(5); got != -2051222400000 {
		t.Errorf("Date.UTC(5) = %v", got)
	}
	if got := UTC(-1, 0); got != -62198755200000 {
		t.Errorf("Date.UTC(-1, 0) = %v", got)
	}
	if got := iso(t, NewLocal(2020.9, 1.7, 3.5, -1.5)); got != "2020-02-03T04:00:00.000Z" {
		t.Errorf("fractional = %s", got)
	}
}

func TestLordHoweExamples(t *testing.T) {
	inZone(t, "Australia/Lord_Howe")
	if got := iso(t, NewLocal(1800, 0, 1)); got != "1799-12-31T13:23:40.000Z" {
		t.Errorf("new Date(1800,0,1) = %s", got)
	}
	if got := NewLocal(1800, 0, 1).GetTimezoneOffset(); got != -636 {
		t.Errorf("getTimezoneOffset = %v", got)
	}
}

func TestToISOStringRange(t *testing.T) {
	for _, c := range []struct {
		t    float64
		want string
	}{
		{0, "1970-01-01T00:00:00.000Z"},
		{8.64e15, "+275760-09-13T00:00:00.000Z"},
		{-8.64e15, "-271821-04-20T00:00:00.000Z"},
		{-62198755200000, "-000001-01-01T00:00:00.000Z"},
	} {
		if got := iso(t, FromTime(c.t)); got != c.want {
			t.Errorf("toISOString(%v) = %s, want %s", c.t, got, c.want)
		}
	}
	if _, ok := FromTime(8.64e15 + 1).ToISOString(); ok {
		t.Error("toISOString beyond TimeClip should fail")
	}
}
