// Ported from src/types.ts.

package core

// Component names a date/time component: "year", "month", "day", "weekday", "hour", "minute",
// "second", "millisecond", "meridiem" or "timezoneOffset".
type Component string

// The date/time components.
const (
	ComponentYear           Component = "year"
	ComponentMonth          Component = "month"
	ComponentDay            Component = "day"
	ComponentWeekday        Component = "weekday"
	ComponentHour           Component = "hour"
	ComponentMinute         Component = "minute"
	ComponentSecond         Component = "second"
	ComponentMillisecond    Component = "millisecond"
	ComponentMeridiem       Component = "meridiem"
	ComponentTimezoneOffset Component = "timezoneOffset"
)

// AllComponents lists the ten components in their canonical order.
var AllComponents = []Component{
	ComponentYear, ComponentMonth, ComponentDay, ComponentWeekday, ComponentHour,
	ComponentMinute, ComponentSecond, ComponentMillisecond, ComponentMeridiem, ComponentTimezoneOffset,
}

// ComponentValue is one entry of a `{ [c in Component]?: number }` object; Components keeps the
// object's key order.
type ComponentValue struct {
	Component Component
	Value     float64
}

// Components is a `{ [c in Component]?: number }` object literal, in key order. A Value equal to
// Undefined (see IsUndefined) is a property whose value is undefined.
type Components []ComponentValue

// Meridiem values (JS enum Meridiem).
const (
	AM = 0
	PM = 1
)

// Weekday values (JS enum Weekday).
const (
	Sunday    = 0
	Monday    = 1
	Tuesday   = 2
	Wednesday = 3
	Thursday  = 4
	Friday    = 5
	Saturday  = 6
)

// Month values (JS enum Month).
const (
	January   = 1
	February  = 2
	March     = 3
	April     = 4
	May       = 5
	June      = 6
	July      = 7
	August    = 8
	September = 9
	October   = 10
	November  = 11
	December  = 12
)

// AmbiguousTimezoneMap is a time zone abbreviation whose offset depends on whether daylight saving
// time is in effect. DstStart and DstEnd return a Date time value for the given year.
type AmbiguousTimezoneMap struct {
	TimezoneOffsetDuringDst float64
	TimezoneOffsetNonDst    float64
	DstStart                func(year float64) float64
	DstEnd                  func(year float64) float64
}

// TimezoneAbbr is one value of a TimezoneAbbrMap: a fixed offset in minutes, or, when Ambiguous
// is not nil, an AmbiguousTimezoneMap.
type TimezoneAbbr struct {
	Offset    float64
	Ambiguous *AmbiguousTimezoneMap
}

// TimezoneAbbrMap maps time zone abbreviations to offsets (JS TimezoneAbbrMap).
type TimezoneAbbrMap = map[string]TimezoneAbbr

// ParsingOption is chrono's ParsingOption; a nil *ParsingOption is an absent option.
type ParsingOption struct {
	// ForwardDate makes the results refer to dates after the reference date.
	ForwardDate bool
	// Timezones holds additional time zone abbreviations (nil: none).
	Timezones map[string]TimezoneAbbr
}

// ParsingReference is chrono's ParsingReference. Instant is a JavaScript Date's time value
// (milliseconds since the epoch). Timezone is nil, a string (an abbreviation) or a float64
// (minutes east of UTC).
type ParsingReference struct {
	Instant  float64
	Timezone any
}

// Error is a JavaScript exception thrown by the original code (for example a TypeError). The Go
// port panics with an *Error exactly where chrono throws.
type Error struct {
	Name    string
	Message string
}

// Error returns "Name: Message".
func (e *Error) Error() string { return e.Name + ": " + e.Message }

func throwTypeError(msg string) {
	panic(&Error{Name: "TypeError", Message: msg})
}
