// Ported from src/locales/nl/constants.ts.

package nl

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

// WeekdayDictionary is WEEKDAY_DICTIONARY.
var WeekdayDictionary = core.NewDictionary([]core.Entry[float64]{
	{Key: "zondag", Value: 0},
	{Key: "zon", Value: 0},
	{Key: "zon.", Value: 0},
	{Key: "zo", Value: 0},
	{Key: "zo.", Value: 0},
	{Key: "maandag", Value: 1},
	{Key: "ma", Value: 1},
	{Key: "ma.", Value: 1},
	{Key: "dinsdag", Value: 2},
	{Key: "din", Value: 2},
	{Key: "din.", Value: 2},
	{Key: "di", Value: 2},
	{Key: "di.", Value: 2},
	{Key: "woensdag", Value: 3},
	{Key: "woe", Value: 3},
	{Key: "woe.", Value: 3},
	{Key: "wo", Value: 3},
	{Key: "wo.", Value: 3},
	{Key: "donderdag", Value: 4},
	{Key: "dond", Value: 4},
	{Key: "dond.", Value: 4},
	{Key: "do", Value: 4},
	{Key: "do.", Value: 4},
	{Key: "vrijdag", Value: 5},
	{Key: "vrij", Value: 5},
	{Key: "vrij.", Value: 5},
	{Key: "vr", Value: 5},
	{Key: "vr.", Value: 5},
	{Key: "zaterdag", Value: 6},
	{Key: "zat", Value: 6},
	{Key: "zat.", Value: 6},
	{Key: "za", Value: 6},
	{Key: "za.", Value: 6},
}...)

// MonthDictionary is MONTH_DICTIONARY.
var MonthDictionary = core.NewDictionary([]core.Entry[float64]{
	{Key: "januari", Value: 1},
	{Key: "jan", Value: 1},
	{Key: "jan.", Value: 1},
	{Key: "februari", Value: 2},
	{Key: "feb", Value: 2},
	{Key: "feb.", Value: 2},
	{Key: "maart", Value: 3},
	{Key: "mar", Value: 3},
	{Key: "mar.", Value: 3},
	{Key: "mrt", Value: 3},
	{Key: "mrt.", Value: 3},
	{Key: "april", Value: 4},
	{Key: "apr", Value: 4},
	{Key: "apr.", Value: 4},
	{Key: "mei", Value: 5},
	{Key: "juni", Value: 6},
	{Key: "jun", Value: 6},
	{Key: "jun.", Value: 6},
	{Key: "juli", Value: 7},
	{Key: "jul", Value: 7},
	{Key: "jul.", Value: 7},
	{Key: "augustus", Value: 8},
	{Key: "aug", Value: 8},
	{Key: "aug.", Value: 8},
	{Key: "september", Value: 9},
	{Key: "sep", Value: 9},
	{Key: "sep.", Value: 9},
	{Key: "sept", Value: 9},
	{Key: "sept.", Value: 9},
	{Key: "oktober", Value: 10},
	{Key: "okt", Value: 10},
	{Key: "okt.", Value: 10},
	{Key: "november", Value: 11},
	{Key: "nov", Value: 11},
	{Key: "nov.", Value: 11},
	{Key: "december", Value: 12},
	{Key: "dec", Value: 12},
	{Key: "dec.", Value: 12},
}...)

// IntegerWordDictionary is INTEGER_WORD_DICTIONARY.
var IntegerWordDictionary = core.NewDictionary([]core.Entry[float64]{
	{Key: "een", Value: 1},
	{Key: "twee", Value: 2},
	{Key: "drie", Value: 3},
	{Key: "vier", Value: 4},
	{Key: "vijf", Value: 5},
	{Key: "zes", Value: 6},
	{Key: "zeven", Value: 7},
	{Key: "acht", Value: 8},
	{Key: "negen", Value: 9},
	{Key: "tien", Value: 10},
	{Key: "elf", Value: 11},
	{Key: "twaalf", Value: 12},
}...)

// OrdinalWordDictionary is ORDINAL_WORD_DICTIONARY.
var OrdinalWordDictionary = core.NewDictionary([]core.Entry[float64]{
	{Key: "eerste", Value: 1},
	{Key: "tweede", Value: 2},
	{Key: "derde", Value: 3},
	{Key: "vierde", Value: 4},
	{Key: "vijfde", Value: 5},
	{Key: "zesde", Value: 6},
	{Key: "zevende", Value: 7},
	{Key: "achtste", Value: 8},
	{Key: "negende", Value: 9},
	{Key: "tiende", Value: 10},
	{Key: "elfde", Value: 11},
	{Key: "twaalfde", Value: 12},
	{Key: "dertiende", Value: 13},
	{Key: "veertiende", Value: 14},
	{Key: "vijftiende", Value: 15},
	{Key: "zestiende", Value: 16},
	{Key: "zeventiende", Value: 17},
	{Key: "achttiende", Value: 18},
	{Key: "negentiende", Value: 19},
	{Key: "twintigste", Value: 20},
	{Key: "eenentwintigste", Value: 21},
	{Key: "tweeëntwintigste", Value: 22},
	{Key: "drieentwintigste", Value: 23},
	{Key: "vierentwintigste", Value: 24},
	{Key: "vijfentwintigste", Value: 25},
	{Key: "zesentwintigste", Value: 26},
	{Key: "zevenentwintigste", Value: 27},
	{Key: "achtentwintig", Value: 28},
	{Key: "negenentwintig", Value: 29},
	{Key: "dertigste", Value: 30},
	{Key: "eenendertigste", Value: 31},
}...)

// TimeUnitDictionary is TIME_UNIT_DICTIONARY.
var TimeUnitDictionary = core.NewDictionary([]core.Entry[string]{
	{Key: "sec", Value: "second"},
	{Key: "second", Value: "second"},
	{Key: "seconden", Value: "second"},
	{Key: "min", Value: "minute"},
	{Key: "mins", Value: "minute"},
	{Key: "minute", Value: "minute"},
	{Key: "minuut", Value: "minute"},
	{Key: "minuten", Value: "minute"},
	{Key: "minuutje", Value: "minute"},
	{Key: "h", Value: "hour"},
	{Key: "hr", Value: "hour"},
	{Key: "hrs", Value: "hour"},
	{Key: "uur", Value: "hour"},
	{Key: "u", Value: "hour"},
	{Key: "uren", Value: "hour"},
	{Key: "dag", Value: "day"},
	{Key: "dagen", Value: "day"},
	{Key: "week", Value: "week"},
	{Key: "weken", Value: "week"},
	{Key: "maand", Value: "month"},
	{Key: "maanden", Value: "month"},
	{Key: "jaar", Value: "year"},
	{Key: "jr", Value: "year"},
	{Key: "jaren", Value: "year"},
}...)

// monthLookup is `MONTH_DICTIONARY[name]` (undefined when name is not a key).
func monthLookup(name string) float64 {
	if v, ok := MonthDictionary.Get(name); ok {
		return v
	}
	return core.Undefined
}

// weekdayLookup is `WEEKDAY_DICTIONARY[name]` (undefined when name is not a key).
func weekdayLookup(name string) float64 {
	if v, ok := WeekdayDictionary.Get(name); ok {
		return v
	}
	return core.Undefined
}

// timeUnitLookup is `TIME_UNIT_DICTIONARY[name]` used as an object key ("undefined" when name is
// not a key).
func timeUnitLookup(name string) string {
	if v, ok := TimeUnitDictionary.Get(name); ok {
		return v
	}
	return "undefined"
}

//-----------------------------

// NumberPattern is NUMBER_PATTERN.
var NumberPattern = "(?:" + core.MatchAnyPattern(IntegerWordDictionary.Keys()) +
	"|[0-9]+|[0-9]+[\\.,][0-9]+|halve?|half|paar)"

var reHalve = core.MustRegExp("halve?", "")

// ParseNumberPattern parses a NumberPattern match.
func ParseNumberPattern(match string) float64 {
	num := js.ToLowerCase(match)
	if v, ok := IntegerWordDictionary.Get(num); ok {
		return v
	} else if num == "paar" {
		return 2
	} else if num == "half" || reHalve.Exec(num) != nil {
		return 0.5
	}
	// Replace "," with "." to support some European languages
	return js.ParseFloat(js.Replace(num, ",", "."))
}

//-----------------------------

// OrdinalNumberPattern is ORDINAL_NUMBER_PATTERN.
var OrdinalNumberPattern = "(?:" + core.MatchAnyPattern(OrdinalWordDictionary.Keys()) + "|[0-9]{1,2}(?:ste|de)?)"

var reOrdinalSuffix = core.MustRegExp("(?:ste|de)$", "i")

// ParseOrdinalNumberPattern parses an OrdinalNumberPattern match.
func ParseOrdinalNumberPattern(match string) float64 {
	num := js.ToLowerCase(match)
	if v, ok := OrdinalWordDictionary.Get(num); ok {
		return v
	}
	num = reOrdinalSuffix.Replace(num, "")
	return js.ParseInt(num, 0)
}

//-----------------------------

// YearPattern is YEAR_PATTERN.
const YearPattern = "(?:[1-9][0-9]{0,3}\\s*(?:voor Christus|na Christus)|[1-2][0-9]{3}|[5-9][0-9])"

var (
	reVoorChristus = core.MustRegExp("voor Christus", "i")
	reNaChristus   = core.MustRegExp("na Christus", "i")
)

// ParseYear parses a YearPattern match.
func ParseYear(match string) float64 {
	if reVoorChristus.Test(match) {
		// Before Christ
		match = reVoorChristus.Replace(match, "")
		return -js.ParseInt(match, 0)
	}

	if reNaChristus.Test(match) {
		match = reNaChristus.Replace(match, "")
		return js.ParseInt(match, 0)
	}

	rawYearNumber := js.ParseInt(match, 0)
	return core.FindMostLikelyADYear(rawYearNumber)
}

//-----------------------------

var singleTimeUnitPattern = "(" + NumberPattern + ")\\s{0,5}(" + core.MatchAnyPattern(TimeUnitDictionary.Keys()) + ")\\s{0,5}"
var singleTimeUnitRegex = core.MustRegExp(singleTimeUnitPattern, "i")

// TimeUnitsPattern is TIME_UNITS_PATTERN.
var TimeUnitsPattern = core.RepeatedTimeunitPattern("(?:(?:binnen|in)\\s*)?", singleTimeUnitPattern, core.DefaultTimeunitConnectorPattern)

// ParseDuration parses a TimeUnitsPattern match (an empty duration when no time unit is found).
func ParseDuration(timeunitText string) *core.Duration {
	fragments := core.NewDuration()
	remainingText := timeunitText
	match := singleTimeUnitRegex.Exec(remainingText)
	for match != nil {
		collectDateTimeFragment(fragments, match)
		remainingText = js.SubstringFrom(remainingText, js.Len(match.Get(0)))
		match = singleTimeUnitRegex.Exec(remainingText)
	}
	return fragments
}

func collectDateTimeFragment(fragments *core.Duration, match *core.RegExpMatch) {
	num := ParseNumberPattern(match.Get(1))
	unit := timeUnitLookup(js.ToLowerCase(match.Get(2)))
	fragments.Set(unit, num)
}
