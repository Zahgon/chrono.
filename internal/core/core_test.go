package core

import (
	"math"
	"slices"
	"testing"
	"time"

	"github.com/wanasit/chrono-go/internal/jsdate"
)

// Expected values in this file were recorded with the original chrono (node).

func TestMatchAnyPatternKeyOrder(t *testing.T) {
	d := NewDictionary(
		Entry[int]{"b", 1}, Entry[int]{"10", 2}, Entry[int]{"a", 3},
		Entry[int]{"2", 4}, Entry[int]{"ab.", 5}, Entry[int]{"cd", 6},
	)
	if got, want := d.Keys(), []string{"2", "10", "b", "a", "ab.", "cd"}; !slices.Equal(got, want) {
		t.Errorf("Keys() = %q, want %q", got, want)
	}
	if got, want := MatchAnyPattern(d.Keys()), `(?:ab\.|10|cd|2|b|a)`; got != want {
		t.Errorf("MatchAnyPattern = %s, want %s", got, want)
	}
}

func TestAddDurationRewritesDuration(t *testing.T) {
	prev := jsdate.SetLocation(time.UTC)
	defer jsdate.SetLocation(prev)

	d := NewDuration(DurationValue{"y", 1.5}, DurationValue{"d", 0.25})
	r := AddDuration(jsdate.FromTime(0), d)
	if r.Time() != 47196000000 {
		t.Errorf("AddDuration = %v, want 47196000000", r.Time())
	}
	want := []DurationValue{{"year", 1.5}, {"day", 0.25}, {"month", 6}, {"hour", 6}}
	if got := d.Values(); !slices.Equal(got, want) {
		t.Errorf("duration after AddDuration = %v, want %v", got, want)
	}
}

func TestAddDurationNilThrows(t *testing.T) {
	defer func() {
		e, ok := recover().(*Error)
		if !ok || e.Name != "TypeError" || e.Message != "Cannot read properties of null (reading 'y')" {
			t.Errorf("recovered %v", e)
		}
	}()
	AddDuration(jsdate.FromTime(0), nil)
}

func TestNullAndUndefinedComponentValues(t *testing.T) {
	ref := ReferenceFromDate(0)
	c := NewParsingComponents(ref, Components{{ComponentDay, Undefined}, {ComponentHour, Null}})
	if v, ok := c.Get(ComponentDay); ok || !math.IsNaN(v) {
		t.Errorf("Get(day) = %v, %v; want NaN (undefined), false", v, ok)
	}
	if v, ok := c.Get(ComponentHour); ok || v != 0 {
		t.Errorf("Get(hour) = %v, %v; want 0 (null), false", v, ok)
	}
	if !c.IsCertain(ComponentDay) || !c.IsCertain(ComponentHour) || c.IsNull(ComponentDay) || !c.IsNull(ComponentHour) {
		t.Error("certainty / null checks of stored null and undefined")
	}
	if got := c.GetCertainComponents(); !slices.Equal(got, []Component{ComponentDay, ComponentHour}) {
		t.Errorf("GetCertainComponents = %v", got)
	}
}

func TestAssignKeepsKeyOrderAndDeleteMovesIt(t *testing.T) {
	c := NewParsingComponents(ReferenceFromDate(0), nil)
	c.Assign(ComponentMonth, 1).Assign(ComponentDay, 2).Assign(ComponentMonth, 3)
	c.Delete(ComponentMonth)
	c.Assign(ComponentMonth, 4)
	if got := c.GetCertainComponents(); !slices.Equal(got, []Component{ComponentDay, ComponentMonth}) {
		t.Errorf("GetCertainComponents = %v", got)
	}
}

func TestRegExpLastIndex(t *testing.T) {
	re := MustRegExp("a", "g")
	if m := re.Exec("xaa"); m == nil || m.Index != 1 || re.LastIndex != 2 {
		t.Fatalf("first exec: %+v lastIndex %d", m, re.LastIndex)
	}
	if m := re.Exec("xaa"); m == nil || m.Index != 2 || re.LastIndex != 3 {
		t.Fatalf("second exec: %+v lastIndex %d", m, re.LastIndex)
	}
	if m := re.Exec("xaa"); m != nil || re.LastIndex != 0 {
		t.Fatalf("third exec: %+v lastIndex %d", m, re.LastIndex)
	}
}

func TestTimezoneOffsets(t *testing.T) {
	prev := jsdate.SetLocation(time.UTC)
	defer jsdate.SetLocation(prev)
	// CET at 2021-07-01 (DST) and 2021-01-01; ET around 2021-03-14 02:00 local (here UTC).
	cases := []struct {
		abbr    string
		instant float64
		want    float64
	}{
		{"CET", 1625097600000, 120}, {"CET", 1609459200000, 60},
		{"ET", 1615687140000, -300}, {"ET", 1615687260000, -240}, {"JST", 0, 540},
	}
	for _, c := range cases {
		got, ok := ToTimezoneOffset(c.abbr, jsdate.FromTime(c.instant), nil)
		if !ok || got != c.want {
			t.Errorf("ToTimezoneOffset(%s, %v) = %v, %v; want %v", c.abbr, c.instant, got, ok, c.want)
		}
	}
	if _, ok := ToTimezoneOffset("XYZ", jsdate.FromTime(0), nil); ok {
		t.Error("unknown abbreviation should be null")
	}
	if _, ok := ToTimezoneOffset("ET", nil, nil); ok {
		t.Error("ambiguous abbreviation without a date should be null")
	}
}
