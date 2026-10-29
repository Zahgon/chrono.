package jsdate

import (
	"math"
	"sync/atomic"
	"time"
)

// zone is the local time zone. It is read through SetLocation/location so tests can swap it;
// by default it is time.Local (the zone named by the TZ environment variable, with
// time.LoadLocation semantics).
var zone atomic.Pointer[time.Location]

// SetLocation makes loc the local time zone of every Date (nil restores time.Local) and returns
// the previous one. It exists for tests and harnesses; programs normally rely on TZ.
func SetLocation(loc *time.Location) *time.Location {
	prev := location()
	zone.Store(loc)
	return prev
}

// location returns the current local time zone.
func location() *time.Location {
	if loc := zone.Load(); loc != nil {
		return loc
	}
	return time.Local
}

// offsetAt returns the local zone's UTC offset in seconds in effect at the instant sec (Unix
// seconds).
func offsetAt(loc *time.Location, sec int64) int64 {
	_, off := time.Unix(sec, 0).In(loc).Zone()
	return int64(off)
}

// localOffset returns the offset in milliseconds that LocalTime adds to the time value t (an
// integral, finite time value): the zone's exact offset at t, seconds included.
func localOffset(t float64) float64 {
	sec := floorDiv(int64(t), 1000)
	return float64(offsetAt(location(), sec) * 1000)
}

// toLocal is LocalTime(t).
func toLocal(t float64) float64 { return t + localOffset(t) }

// offsetFromLocal returns the offset in milliseconds to subtract from the local time value l to
// get UTC, as V8 computes it with ICU's getOffsetFromLocal(UCAL_TZ_LOCAL_FORMER,
// UCAL_TZ_LOCAL_FORMER): for every transition at instant T from offset b to offset a, local
// times at or after T+max(a, b) use a, so a wall time in a gap is read with the offset before
// the transition (it moves forward) and a wall time that occurs twice is read as the earlier
// instant. The transitions come from the zone itself (time.Time.ZoneBounds); time.Date is not
// used because its choice differs in gaps and overlaps.
func offsetFromLocal(l float64) float64 {
	loc := location()
	lsec := floorDiv(int64(l), 1000) // transitions and offsets are whole seconds
	const window = 2 * 86400         // larger than any UTC offset
	t := lsec - window
	off := offsetAt(loc, t)
	for {
		_, end := time.Unix(t, 0).In(loc).ZoneBounds()
		if end.IsZero() {
			break
		}
		tr := end.Unix()
		if tr <= t { // defensive: ZoneBounds must advance
			tr = t + 1
		}
		if tr > lsec+window {
			break
		}
		before := offsetAt(loc, tr-1)
		after := offsetAt(loc, tr)
		if lsec >= tr+max(before, after) {
			off = after
		}
		t = tr
	}
	return float64(off * 1000)
}

// localToUTC converts a local time value to a time value as V8 does for the constructor and
// the local setters: values beyond V8's pre-conversion bound are NaN.
func localToUTC(l float64) float64 {
	if math.IsNaN(l) || l < -maxTimeBeforeUTC || l > maxTimeBeforeUTC {
		return math.NaN()
	}
	l = math.Trunc(l)
	return l - offsetFromLocal(l)
}

func floorDiv(a, b int64) int64 {
	q := a / b
	if (a%b != 0) && ((a < 0) != (b < 0)) {
		q--
	}
	return q
}
