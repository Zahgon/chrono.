// Package jsdate is the JavaScript Date as V8 (Node.js 24) implements it, for the chrono port.
//
// A Date holds a time value: float64 milliseconds since the Unix epoch, NaN when the date is
// invalid, always an integer within ±8.64e15 otherwise (TimeClip). The constructor, Date.UTC and
// every getter and setter follow ECMAScript's MakeDay/MakeTime/MakeDate/TimeClip with V8's
// argument handling (fractional arguments are truncated toward zero, NaN and ±Infinity make the
// result invalid, years 0 to 99 in the constructor and Date.UTC mean 1900 to 1999).
//
// Local time is the zone named by the TZ environment variable as Go loads it (time.Local; no tz
// data is embedded), with V8's rules on top: LocalTime(t) adds the zone's exact offset at t
// (seconds included, e.g. local mean time), a local wall time that does not exist (DST gap) is
// read with the offset in effect before the transition so it moves forward, a wall time that
// occurs twice is read as the earlier instant, and GetTimezoneOffset truncates the offset to
// whole minutes toward zero. SetLocation swaps the zone (for tests).
package jsdate

import (
	"fmt"
	"math"
	"time"
)

const (
	msPerDay    = 86400000.0
	msPerHour   = 3600000.0
	msPerMinute = 60000.0
	msPerSecond = 1000.0
	maxTime     = 8.64e15
	// maxTimeBeforeUTC is V8's bound on a local time value before it is converted to UTC
	// (DateCache::kMaxTimeBeforeUTCInMs: kMaxTimeInMs + 30 days).
	maxTimeBeforeUTC = maxTime + 30*msPerDay
)

// Date is a JavaScript Date. The zero value is the epoch; use FromTime, NewLocal or Now.
type Date struct {
	t float64 // time value; NaN when invalid
}

// FromTime is new Date(t): the time value TimeClip(t).
func FromTime(t float64) *Date { return &Date{timeClip(t)} }

// Now is new Date(): the current time to the millisecond.
func Now() *Date { return &Date{float64(time.Now().UnixMilli())} }

// NewLocal is new Date(year, month[, date, hours, minutes, seconds, ms]) in local time. Missing
// arguments default to month 0, date 1 and zero time; with a single argument it is
// new Date(year, 0) (use FromTime for new Date(t)); with none it is Now. Arguments after the
// seventh are ignored, as in JavaScript.
func NewLocal(args ...float64) *Date {
	if len(args) == 0 {
		return Now()
	}
	return &Date{timeClip(localToUTC(dateArgs(args)))}
}

// UTC is Date.UTC(year[, month, date, hours, minutes, seconds, ms]); Date.UTC() is NaN.
func UTC(args ...float64) float64 {
	if len(args) == 0 {
		return math.NaN()
	}
	return timeClip(dateArgs(args))
}

// Clone is new Date(d).
func (d *Date) Clone() *Date { return &Date{d.t} }

// Time is getTime/valueOf: the time value, NaN when invalid.
func (d *Date) Time() float64 { return d.t }

// IsValid reports whether the time value is not NaN.
func (d *Date) IsValid() bool { return !math.IsNaN(d.t) }

// String is Date.prototype.toISOString, or "Invalid Date"; it exists for debugging output.
func (d *Date) String() string {
	if s, ok := d.ToISOString(); ok {
		return s
	}
	return "Invalid Date"
}

// ---------------------------------------------------------------------------------------------
// abstract operations

// toIntegerOrInfinity is ToIntegerOrInfinity (NaN and -0 become +0).
func toIntegerOrInfinity(x float64) float64 {
	if math.IsNaN(x) {
		return 0
	}
	if math.IsInf(x, 0) {
		return x
	}
	if t := math.Trunc(x); t != 0 {
		return t
	}
	return 0
}

var dayFromMonth = [2][12]int64{
	{0, 31, 59, 90, 120, 151, 181, 212, 243, 273, 304, 334},
	{0, 31, 60, 91, 121, 152, 182, 213, 244, 274, 305, 335},
}

// makeDay is V8's MakeDay: year and month are truncated to integers and must lie within
// ±1e6 and ±1e7, date must be finite.
func makeDay(year, month, date float64) float64 {
	if !(-1000000 <= year && year <= 1000000 && -10000000 <= month && month <= 10000000) ||
		math.IsNaN(date) || math.IsInf(date, 0) {
		return math.NaN()
	}
	y, m := int64(year), int64(month)
	y += m / 12
	m %= 12
	if m < 0 {
		m += 12
		y--
	}
	const yearDelta = 399999
	const base = 365*(1970+yearDelta) + (1970+yearDelta)/4 - (1970+yearDelta)/100 + (1970+yearDelta)/400
	day := 365*(y+yearDelta) + (y+yearDelta)/4 - (y+yearDelta)/100 + (y+yearDelta)/400 - base
	leap := 0
	if y%4 == 0 && (y%100 != 0 || y%400 == 0) {
		leap = 1
	}
	day += dayFromMonth[leap][m]
	return float64(day-1) + toIntegerOrInfinity(date)
}

// makeTime is MakeTime.
func makeTime(h, m, s, ms float64) float64 {
	if !finite(h) || !finite(m) || !finite(s) || !finite(ms) {
		return math.NaN()
	}
	return toIntegerOrInfinity(h)*msPerHour + toIntegerOrInfinity(m)*msPerMinute +
		toIntegerOrInfinity(s)*msPerSecond + toIntegerOrInfinity(ms)
}

// makeDate is MakeDate.
func makeDate(day, tm float64) float64 {
	if !finite(day) || !finite(tm) {
		return math.NaN()
	}
	return day*msPerDay + tm
}

// timeClip is TimeClip.
func timeClip(t float64) float64 {
	if !finite(t) || math.Abs(t) > maxTime {
		return math.NaN()
	}
	return toIntegerOrInfinity(t)
}

func finite(x float64) bool { return !math.IsNaN(x) && !math.IsInf(x, 0) }

// dateArgs reads the (year[, month, date, hours, minutes, seconds, ms]) arguments of the
// constructor and Date.UTC and returns the time value they describe before UTC conversion.
func dateArgs(args []float64) float64 {
	n := [7]float64{math.NaN(), 0, 1, 0, 0, 0, 0}
	copy(n[:], args)
	if !math.IsNaN(n[0]) {
		if y := toIntegerOrInfinity(n[0]); 0 <= y && y <= 99 {
			n[0] = 1900 + y
		}
	}
	return makeDate(makeDay(n[0], n[1], n[2]), makeTime(n[3], n[4], n[5], n[6]))
}

// ---------------------------------------------------------------------------------------------
// fields

type fields struct {
	year, month, date, hours, minutes, seconds, ms, weekday float64
}

func nanFields() fields {
	n := math.NaN()
	return fields{n, n, n, n, n, n, n, n}
}

func posMod(a, b float64) float64 {
	r := math.Mod(a, b)
	if r < 0 {
		r += b
	}
	if r == 0 {
		return 0
	}
	return r
}

// fieldsOf splits a (finite, integral) time value into calendar fields.
func fieldsOf(t float64) fields {
	day := math.Floor(t / msPerDay)
	within := t - day*msPerDay
	y, m, dd := civilFromDays(int64(day))
	return fields{
		year:    float64(y),
		month:   float64(m - 1),
		date:    float64(dd),
		hours:   math.Floor(within / msPerHour),
		minutes: posMod(math.Floor(within/msPerMinute), 60),
		seconds: posMod(math.Floor(within/msPerSecond), 60),
		ms:      posMod(within, 1000),
		weekday: posMod(day+4, 7),
	}
}

// civilFromDays converts days since 1970-01-01 to a proleptic Gregorian year, month (1-12) and
// day of month.
func civilFromDays(z int64) (y, m, d int64) {
	z += 719468
	era := z / 146097
	if z < 0 && z%146097 != 0 {
		era--
	}
	doe := z - era*146097
	yoe := (doe - doe/1460 + doe/36524 - doe/146096) / 365
	y = yoe + era*400
	doy := doe - (365*yoe + yoe/4 - yoe/100)
	mp := (5*doy + 2) / 153
	d = doy - (153*mp+2)/5 + 1
	if mp < 10 {
		m = mp + 3
	} else {
		m = mp - 9
	}
	if m <= 2 {
		y++
	}
	return
}

func (d *Date) local() fields {
	if math.IsNaN(d.t) {
		return nanFields()
	}
	return fieldsOf(toLocal(d.t))
}

func (d *Date) utc() fields {
	if math.IsNaN(d.t) {
		return nanFields()
	}
	return fieldsOf(d.t)
}

// GetFullYear is Date.prototype.getFullYear.
func (d *Date) GetFullYear() float64 { return d.local().year }

// GetMonth is Date.prototype.getMonth (0-11).
func (d *Date) GetMonth() float64 { return d.local().month }

// GetDate is Date.prototype.getDate.
func (d *Date) GetDate() float64 { return d.local().date }

// GetDay is Date.prototype.getDay (0 = Sunday).
func (d *Date) GetDay() float64 { return d.local().weekday }

// GetHours is Date.prototype.getHours.
func (d *Date) GetHours() float64 { return d.local().hours }

// GetMinutes is Date.prototype.getMinutes.
func (d *Date) GetMinutes() float64 { return d.local().minutes }

// GetSeconds is Date.prototype.getSeconds.
func (d *Date) GetSeconds() float64 { return d.local().seconds }

// GetMilliseconds is Date.prototype.getMilliseconds.
func (d *Date) GetMilliseconds() float64 { return d.local().ms }

// GetUTCFullYear is Date.prototype.getUTCFullYear.
func (d *Date) GetUTCFullYear() float64 { return d.utc().year }

// GetUTCMonth is Date.prototype.getUTCMonth.
func (d *Date) GetUTCMonth() float64 { return d.utc().month }

// GetUTCDate is Date.prototype.getUTCDate.
func (d *Date) GetUTCDate() float64 { return d.utc().date }

// GetUTCDay is Date.prototype.getUTCDay.
func (d *Date) GetUTCDay() float64 { return d.utc().weekday }

// GetUTCHours is Date.prototype.getUTCHours.
func (d *Date) GetUTCHours() float64 { return d.utc().hours }

// GetUTCMinutes is Date.prototype.getUTCMinutes.
func (d *Date) GetUTCMinutes() float64 { return d.utc().minutes }

// GetUTCSeconds is Date.prototype.getUTCSeconds.
func (d *Date) GetUTCSeconds() float64 { return d.utc().seconds }

// GetUTCMilliseconds is Date.prototype.getUTCMilliseconds.
func (d *Date) GetUTCMilliseconds() float64 { return d.utc().ms }

// GetTimezoneOffset is Date.prototype.getTimezoneOffset as V8 computes it:
// (t - LocalTime(t)) / 60000 truncated toward zero, so an offset with seconds (local mean time)
// loses them here while the local getters keep them.
func (d *Date) GetTimezoneOffset() float64 {
	if math.IsNaN(d.t) {
		return math.NaN()
	}
	diff := int64(d.t) - int64(toLocal(d.t))
	return float64(diff / 60000)
}

// ---------------------------------------------------------------------------------------------
// setters

type field int

const (
	fYear field = iota
	fMonth
	fDate
	fHours
	fMinutes
	fSeconds
	fMs
)

// set implements every setter: f is the first field set, args its value and the optional
// following ones, utc selects the setUTC variant. It stores and returns the new time value.
func (d *Date) set(f field, utc bool, args []float64) float64 {
	arg := func(i int, def float64) float64 {
		if i < len(args) {
			return args[i]
		}
		return def
	}
	t := d.t
	if math.IsNaN(t) {
		if f != fYear {
			return d.t
		}
		t = 0 // setFullYear/setUTCFullYear on an invalid date start from +0
	} else if !utc {
		t = toLocal(t)
	}
	c := fieldsOf(t)
	day := math.Floor(t / msPerDay)
	within := t - day*msPerDay
	v := args[0]
	var nt float64
	switch f {
	case fYear:
		nt = makeDate(makeDay(v, arg(1, c.month), arg(2, c.date)), within)
	case fMonth:
		nt = makeDate(makeDay(c.year, v, arg(1, c.date)), within)
	case fDate:
		nt = makeDate(makeDay(c.year, c.month, v), within)
	case fHours:
		nt = makeDate(day, makeTime(v, arg(1, c.minutes), arg(2, c.seconds), arg(3, c.ms)))
	case fMinutes:
		nt = makeDate(day, makeTime(c.hours, v, arg(1, c.seconds), arg(2, c.ms)))
	case fSeconds:
		nt = makeDate(day, makeTime(c.hours, c.minutes, v, arg(1, c.ms)))
	case fMs:
		nt = makeDate(day, makeTime(c.hours, c.minutes, c.seconds, v))
	}
	if utc {
		d.t = timeClip(nt)
	} else {
		d.t = timeClip(localToUTC(nt))
	}
	return d.t
}

// SetTime is Date.prototype.setTime.
func (d *Date) SetTime(t float64) float64 {
	d.t = timeClip(t)
	return d.t
}

// SetFullYear is Date.prototype.setFullYear(year[, month[, date]]).
func (d *Date) SetFullYear(y float64, rest ...float64) float64 {
	return d.set(fYear, false, append([]float64{y}, rest...))
}

// SetMonth is Date.prototype.setMonth(month[, date]).
func (d *Date) SetMonth(m float64, rest ...float64) float64 {
	return d.set(fMonth, false, append([]float64{m}, rest...))
}

// SetDate is Date.prototype.setDate.
func (d *Date) SetDate(date float64) float64 { return d.set(fDate, false, []float64{date}) }

// SetHours is Date.prototype.setHours(hours[, minutes[, seconds[, ms]]]).
func (d *Date) SetHours(h float64, rest ...float64) float64 {
	return d.set(fHours, false, append([]float64{h}, rest...))
}

// SetMinutes is Date.prototype.setMinutes(minutes[, seconds[, ms]]).
func (d *Date) SetMinutes(m float64, rest ...float64) float64 {
	return d.set(fMinutes, false, append([]float64{m}, rest...))
}

// SetSeconds is Date.prototype.setSeconds(seconds[, ms]).
func (d *Date) SetSeconds(s float64, rest ...float64) float64 {
	return d.set(fSeconds, false, append([]float64{s}, rest...))
}

// SetMilliseconds is Date.prototype.setMilliseconds.
func (d *Date) SetMilliseconds(ms float64) float64 { return d.set(fMs, false, []float64{ms}) }

// SetUTCFullYear is Date.prototype.setUTCFullYear(year[, month[, date]]).
func (d *Date) SetUTCFullYear(y float64, rest ...float64) float64 {
	return d.set(fYear, true, append([]float64{y}, rest...))
}

// SetUTCMonth is Date.prototype.setUTCMonth(month[, date]).
func (d *Date) SetUTCMonth(m float64, rest ...float64) float64 {
	return d.set(fMonth, true, append([]float64{m}, rest...))
}

// SetUTCDate is Date.prototype.setUTCDate.
func (d *Date) SetUTCDate(date float64) float64 { return d.set(fDate, true, []float64{date}) }

// SetUTCHours is Date.prototype.setUTCHours(hours[, minutes[, seconds[, ms]]]).
func (d *Date) SetUTCHours(h float64, rest ...float64) float64 {
	return d.set(fHours, true, append([]float64{h}, rest...))
}

// SetUTCMinutes is Date.prototype.setUTCMinutes(minutes[, seconds[, ms]]).
func (d *Date) SetUTCMinutes(m float64, rest ...float64) float64 {
	return d.set(fMinutes, true, append([]float64{m}, rest...))
}

// SetUTCSeconds is Date.prototype.setUTCSeconds(seconds[, ms]).
func (d *Date) SetUTCSeconds(s float64, rest ...float64) float64 {
	return d.set(fSeconds, true, append([]float64{s}, rest...))
}

// SetUTCMilliseconds is Date.prototype.setUTCMilliseconds.
func (d *Date) SetUTCMilliseconds(ms float64) float64 { return d.set(fMs, true, []float64{ms}) }

// ToISOString is Date.prototype.toISOString; ok is false where V8 throws RangeError (an
// invalid date).
func (d *Date) ToISOString() (s string, ok bool) {
	if math.IsNaN(d.t) {
		return "", false
	}
	f := d.utc()
	y := int64(f.year)
	var ys string
	switch {
	case y < 0:
		ys = fmt.Sprintf("-%06d", -y)
	case y > 9999:
		ys = fmt.Sprintf("+%06d", y)
	default:
		ys = fmt.Sprintf("%04d", y)
	}
	return fmt.Sprintf("%s-%02d-%02dT%02d:%02d:%02d.%03dZ", ys, int(f.month)+1, int(f.date),
		int(f.hours), int(f.minutes), int(f.seconds), int(f.ms)), true
}
