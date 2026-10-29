// Ported from src/calculation/duration.ts.

package core

import (
	"math"

	"github.com/wanasit/chrono-go/internal/js"
	"github.com/wanasit/chrono-go/internal/jsdate"
)

// DurationValue is one `unit: value` property of a Duration literal.
type DurationValue struct {
	Unit  string
	Value float64
}

// Duration is a directed time duration as a JavaScript object of values by time unit ("year",
// "month", "week", "day", "hour", "minute", "second", "millisecond", "quarter", or the short
// keys "y", "mo", "M", "w", "d", "h", "m", "s", "ms"). Positive values point into the future.
// Keys keep JavaScript object order; `"unit" in duration` is Has. A nil *Duration is JS null.
type Duration struct {
	obj numObject
}

// NewDuration builds a duration literal, e.g. NewDuration(DurationValue{"day", 1}).
func NewDuration(values ...DurationValue) *Duration {
	d := &Duration{}
	for _, v := range values {
		d.Set(v.Unit, v.Value)
	}
	return d
}

// EmptyDuration is the explicit empty duration { day: 0, second: 0, millisecond: 0 } (one shared
// object, as in the original).
var EmptyDuration = NewDuration(DurationValue{"day", 0}, DurationValue{"second", 0}, DurationValue{"millisecond", 0})

// Set is `duration[unit] = value`.
func (d *Duration) Set(unit string, value float64) { d.obj.set(unit, value) }

// Get returns `duration[unit]` and whether the property exists.
func (d *Duration) Get(unit string) (float64, bool) { return d.obj.get(unit) }

// Has is `unit in duration`.
func (d *Duration) Has(unit string) bool { return d.obj.has(unit) }

// Delete is `delete duration[unit]`.
func (d *Duration) Delete(unit string) { d.obj.del(unit) }

// Keys is Object.keys(duration).
func (d *Duration) Keys() []string {
	out := make([]string, len(d.obj.entries))
	for i, e := range d.obj.entries {
		out[i] = e.key
	}
	return out
}

// Len is Object.keys(duration).length.
func (d *Duration) Len() int { return len(d.obj.entries) }

// Values returns the properties in key order.
func (d *Duration) Values() []DurationValue {
	out := make([]DurationValue, len(d.obj.entries))
	for i, e := range d.obj.entries {
		out[i] = DurationValue{e.key, e.v}
	}
	return out
}

// truthy is `if (duration[key])`.
func (d *Duration) truthy(key string) bool {
	v, ok := d.obj.get(key)
	return ok && !IsUndefined(v) && Truthy(v)
}

// num is `duration[key]` used as a number (undefined -> NaN).
func (d *Duration) num(key string) float64 {
	v, ok := d.obj.get(key)
	if !ok {
		return math.NaN()
	}
	return Num(v)
}

// orZero is `duration?.key ?? 0`.
func (d *Duration) orZero(key string) float64 {
	v, ok := d.obj.get(key)
	if !ok || IsUndefined(v) {
		return 0
	}
	return v
}

// AddDuration returns the date after adding the duration to ref (a new Date). As in the original,
// it rewrites short unit keys of duration to full keys and carries fractional parts into smaller
// units of duration. A nil duration throws the original's TypeError.
func AddDuration(ref *jsdate.Date, duration *Duration) *jsdate.Date {
	date := ref.Clone()
	if duration == nil {
		throwTypeError("Cannot read properties of null (reading 'y')")
	}

	shortKeys := [][2]string{
		{"y", "year"}, {"mo", "month"}, {"M", "month"}, {"w", "week"}, {"d", "day"},
		{"h", "hour"}, {"m", "minute"}, {"s", "second"}, {"ms", "millisecond"},
	}
	for _, sk := range shortKeys {
		if duration.truthy(sk[0]) {
			v, _ := duration.Get(sk[0])
			duration.Set(sk[1], v)
			duration.Delete(sk[0])
		}
	}

	if duration.Has("year") {
		floor := math.Floor(duration.num("year"))
		date.SetFullYear(date.GetFullYear() + floor)
		remainingFraction := duration.num("year") - floor
		if remainingFraction > 0 {
			duration.Set("month", duration.orZero("month"))
			duration.Set("month", duration.num("month")+remainingFraction*12)
		}
	}
	if duration.Has("quarter") {
		floor := math.Floor(duration.num("quarter"))
		date.SetMonth(date.GetMonth() + floor*3)
	}
	if duration.Has("month") {
		floor := math.Floor(duration.num("month"))
		date.SetMonth(date.GetMonth() + floor)
		remainingFraction := duration.num("month") - floor
		if remainingFraction > 0 {
			duration.Set("week", duration.orZero("week"))
			duration.Set("week", duration.num("week")+remainingFraction*4)
		}
	}
	if duration.Has("week") {
		floor := math.Floor(duration.num("week"))
		date.SetDate(date.GetDate() + floor*7)
		remainingFraction := duration.num("week") - floor
		if remainingFraction > 0 {
			duration.Set("day", duration.orZero("day"))
			duration.Set("day", duration.num("day")+js.MathRound(remainingFraction*7))
		}
	}
	if duration.Has("day") {
		floor := math.Floor(duration.num("day"))
		date.SetDate(date.GetDate() + floor)
		remainingFraction := duration.num("day") - floor
		if remainingFraction > 0 {
			duration.Set("hour", duration.orZero("hour"))
			duration.Set("hour", duration.num("hour")+js.MathRound(remainingFraction*24))
		}
	}
	if duration.Has("hour") {
		floor := math.Floor(duration.num("hour"))
		date.SetHours(date.GetHours() + floor)
		remainingFraction := duration.num("hour") - floor
		if remainingFraction > 0 {
			duration.Set("minute", duration.orZero("minute"))
			duration.Set("minute", duration.num("minute")+js.MathRound(remainingFraction*60))
		}
	}
	if duration.Has("minute") {
		floor := math.Floor(duration.num("minute"))
		date.SetMinutes(date.GetMinutes() + floor)
		remainingFraction := duration.num("minute") - floor
		if remainingFraction > 0 {
			duration.Set("second", duration.orZero("second"))
			duration.Set("second", duration.num("second")+js.MathRound(remainingFraction*60))
		}
	}
	if duration.Has("second") {
		floor := math.Floor(duration.num("second"))
		date.SetSeconds(date.GetSeconds() + floor)
		remainingFraction := duration.num("second") - floor
		if remainingFraction > 0 {
			duration.Set("millisecond", duration.orZero("millisecond"))
			duration.Set("millisecond", duration.num("millisecond")+js.MathRound(remainingFraction*1000))
		}
	}
	if duration.Has("millisecond") {
		floor := math.Floor(duration.num("millisecond"))
		date.SetMilliseconds(date.GetMilliseconds() + floor)
	}
	return date
}

// ReverseDuration returns the reversed duration (into the past instead of the future). A nil
// duration (JS null) gives an empty duration, as `for..in` over null does nothing.
func ReverseDuration(duration *Duration) *Duration {
	reversed := &Duration{}
	if duration == nil {
		return reversed
	}
	for _, e := range duration.obj.entries {
		reversed.Set(e.key, -Num(e.v))
	}
	return reversed
}
