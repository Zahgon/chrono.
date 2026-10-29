// Ported from src/index.ts, src/chrono.ts, src/results.ts, src/types.ts and src/timezone.ts (the
// public API; the implementation is in internal/core and internal/en).

// Package chrono is a Go port of chrono (https://github.com/wanasit/chrono) v2.10.1, a natural
// language date parser. It reproduces the original's results exactly, including its JavaScript
// semantics: indices and lengths are UTF-16 code units, regular expressions are ECMAScript regular
// expressions and dates are JavaScript Dates in the process's local time zone (TZ).
//
// The package's Casual and Strict objects parse English (as chrono.casual and chrono.strict do);
// the other locales are in the packages under locales/.
package chrono

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/en"
)

// Component names a date/time component: "year", "month", "day", "weekday", "hour", "minute",
// "second", "millisecond", "meridiem" or "timezoneOffset".
type Component = core.Component

// The date/time components.
const (
	ComponentYear           = core.ComponentYear
	ComponentMonth          = core.ComponentMonth
	ComponentDay            = core.ComponentDay
	ComponentWeekday        = core.ComponentWeekday
	ComponentHour           = core.ComponentHour
	ComponentMinute         = core.ComponentMinute
	ComponentSecond         = core.ComponentSecond
	ComponentMillisecond    = core.ComponentMillisecond
	ComponentMeridiem       = core.ComponentMeridiem
	ComponentTimezoneOffset = core.ComponentTimezoneOffset
)

// Meridiem values.
const (
	AM = core.AM
	PM = core.PM
)

// Weekday values.
const (
	Sunday    = core.Sunday
	Monday    = core.Monday
	Tuesday   = core.Tuesday
	Wednesday = core.Wednesday
	Thursday  = core.Thursday
	Friday    = core.Friday
	Saturday  = core.Saturday
)

// ComponentValue is one entry of a `{ [c in Component]?: number }` object; Components keeps the
// object's key order.
type ComponentValue = core.ComponentValue

// Components is a `{ [c in Component]?: number }` object in key order.
type Components = core.Components

// ParsingReference is chrono's ParsingReference. Instant is a JavaScript Date's time value
// (milliseconds since the epoch). Timezone is nil, a string (an abbreviation) or a float64
// (minutes east of UTC).
type ParsingReference = core.ParsingReference

// TimezoneAbbr is one value of a TimezoneAbbrMap: a fixed offset in minutes, or, when Ambiguous
// is not nil, an AmbiguousTimezoneMap.
type TimezoneAbbr = core.TimezoneAbbr

// AmbiguousTimezoneMap is a DST-dependent time zone abbreviation.
type AmbiguousTimezoneMap = core.AmbiguousTimezoneMap

// ParsingOption is chrono's ParsingOption; a nil *ParsingOption is an absent option.
type ParsingOption = core.ParsingOption

// GetNthWeekdayOfMonth returns the Date (time value) of the nth weekday (0 = Sunday) of the month
// (1-12) of year at the given hour, in local time.
func GetNthWeekdayOfMonth(year float64, month, weekday, n, hour int) float64 {
	return core.GetNthWeekdayOfMonth(year, month, weekday, n, hour)
}

// GetLastWeekdayOfMonth returns the Date (time value) of the last weekday (0 = Sunday) of the
// month (1-12) of year at the given hour, in local time.
func GetLastWeekdayOfMonth(year float64, month, weekday, hour int) float64 {
	return core.GetLastWeekdayOfMonth(year, month, weekday, hour)
}

// ReferenceWithTimezone is the reference instant and time zone offset of a parse call.
type ReferenceWithTimezone = core.ReferenceWithTimezone

// ParsingComponents is a collection of parsed date/time components (certain, implied or
// unknown).
type ParsingComponents = core.ParsingComponents

// NewParsingComponents is `new ParsingComponents(reference, knownComponents)`.
func NewParsingComponents(reference *ReferenceWithTimezone, known Components) *ParsingComponents {
	return core.NewParsingComponents(reference, known)
}

// ParsingResult is one parsed result.
type ParsingResult = core.ParsingResult

// NewParsingResult is `new ParsingResult(reference, index, text, start, end)`.
func NewParsingResult(reference *ReferenceWithTimezone, index int, text string, start, end *ParsingComponents) *ParsingResult {
	return core.NewParsingResult(reference, index, text, start, end)
}

// RegExp is an ECMAScript regular expression.
type RegExp = core.RegExp

// NewRegExp is `new RegExp(source, flags)`; the error's text is V8's SyntaxError message.
func NewRegExp(source, flags string) (*RegExp, error) { return core.NewRegExp(source, flags) }

// RegExpMatch is what RegExp.prototype.exec returns.
type RegExpMatch = core.RegExpMatch

// ParsingContext is the context of one parse call.
type ParsingContext = core.ParsingContext

// Parser recognizes one date format. Extract returns nil (null), *ParsingComponents,
// *ParsingResult or Components.
type Parser = core.Parser

// Refiner refines the list of results.
type Refiner = core.Refiner

// Configuration is an ordered list of parsers and refiners.
type Configuration = core.Configuration

// Chrono applies its parsers and refiners in order.
type Chrono = core.Chrono

// Error is a JavaScript exception of the original (e.g. a TypeError); chrono-go panics with an
// *Error exactly where chrono throws.
type Error = core.Error

// NewChrono is `new Chrono(configuration)`; nil is the English casual configuration.
func NewChrono(configuration *Configuration) *Chrono { return core.NewChrono(configuration) }

// Casual is the English casual Chrono (chrono.casual).
var Casual = en.Casual

// Strict is the English strict Chrono (chrono.strict).
var Strict = en.Strict

// Parse is Casual.Parse (chrono.parse).
func Parse(text string, ref ParsingReference, option *ParsingOption) []*ParsingResult {
	return Casual.Parse(text, ref, option)
}

// ParseDate is Casual.ParseDate (chrono.parseDate); ok is false for null.
func ParseDate(text string, ref ParsingReference, option *ParsingOption) (float64, bool) {
	return Casual.ParseDate(text, ref, option)
}
