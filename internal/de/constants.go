// Ported from src/locales/de/constants.ts.

package de

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

// WeekdayDictionary is WEEKDAY_DICTIONARY.
var WeekdayDictionary = core.NewDictionary([]core.Entry[float64]{
	{Key: "sonntag", Value: 0},
	{Key: "so", Value: 0},
	{Key: "montag", Value: 1},
	{Key: "mo", Value: 1},
	{Key: "dienstag", Value: 2},
	{Key: "di", Value: 2},
	{Key: "mittwoch", Value: 3},
	{Key: "mi", Value: 3},
	{Key: "donnerstag", Value: 4},
	{Key: "do", Value: 4},
	{Key: "freitag", Value: 5},
	{Key: "fr", Value: 5},
	{Key: "samstag", Value: 6},
	{Key: "sa", Value: 6},
}...)

// MonthDictionary is MONTH_DICTIONARY.
var MonthDictionary = core.NewDictionary([]core.Entry[float64]{
	{Key: "januar", Value: 1},
	{Key: "jänner", Value: 1},
	{Key: "janner", Value: 1},
	{Key: "jan", Value: 1},
	{Key: "jan.", Value: 1},
	{Key: "februar", Value: 2},
	{Key: "feber", Value: 2},
	{Key: "feb", Value: 2},
	{Key: "feb.", Value: 2},
	{Key: "märz", Value: 3},
	{Key: "maerz", Value: 3},
	{Key: "mär", Value: 3},
	{Key: "mär.", Value: 3},
	{Key: "mrz", Value: 3},
	{Key: "mrz.", Value: 3},
	{Key: "april", Value: 4},
	{Key: "apr", Value: 4},
	{Key: "apr.", Value: 4},
	{Key: "mai", Value: 5},
	{Key: "juni", Value: 6},
	{Key: "jun", Value: 6},
	{Key: "jun.", Value: 6},
	{Key: "juli", Value: 7},
	{Key: "jul", Value: 7},
	{Key: "jul.", Value: 7},
	{Key: "august", Value: 8},
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
	{Key: "dezember", Value: 12},
	{Key: "dez", Value: 12},
	{Key: "dez.", Value: 12},
}...)

// IntegerWordDictionary is INTEGER_WORD_DICTIONARY.
var IntegerWordDictionary = core.NewDictionary([]core.Entry[float64]{
	{Key: "eins", Value: 1},
	{Key: "eine", Value: 1},
	{Key: "einem", Value: 1},
	{Key: "einen", Value: 1},
	{Key: "einer", Value: 1},
	{Key: "zwei", Value: 2},
	{Key: "drei", Value: 3},
	{Key: "vier", Value: 4},
	{Key: "fünf", Value: 5},
	{Key: "fuenf", Value: 5},
	{Key: "sechs", Value: 6},
	{Key: "sieben", Value: 7},
	{Key: "acht", Value: 8},
	{Key: "neun", Value: 9},
	{Key: "zehn", Value: 10},
	{Key: "elf", Value: 11},
	{Key: "zwölf", Value: 12},
	{Key: "zwoelf", Value: 12},
}...)

// TimeUnitDictionary is TIME_UNIT_DICTIONARY.
var TimeUnitDictionary = core.NewDictionary([]core.Entry[string]{
	{Key: "sek", Value: "second"},
	{Key: "sekunde", Value: "second"},
	{Key: "sekunden", Value: "second"},
	{Key: "min", Value: "minute"},
	{Key: "minute", Value: "minute"},
	{Key: "minuten", Value: "minute"},
	{Key: "h", Value: "hour"},
	{Key: "std", Value: "hour"},
	{Key: "stunde", Value: "hour"},
	{Key: "stunden", Value: "hour"},
	{Key: "tag", Value: "day"},
	{Key: "tage", Value: "day"},
	{Key: "tagen", Value: "day"},
	{Key: "woche", Value: "week"},
	{Key: "wochen", Value: "week"},
	{Key: "monat", Value: "month"},
	{Key: "monate", Value: "month"},
	{Key: "monaten", Value: "month"},
	{Key: "monats", Value: "month"},
	{Key: "quartal", Value: "quarter"},
	{Key: "quartals", Value: "quarter"},
	{Key: "quartale", Value: "quarter"},
	{Key: "quartalen", Value: "quarter"},
	{Key: "a", Value: "year"},
	{Key: "j", Value: "year"},
	{Key: "jr", Value: "year"},
	{Key: "jahr", Value: "year"},
	{Key: "jahre", Value: "year"},
	{Key: "jahren", Value: "year"},
	{Key: "jahres", Value: "year"},
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
	"|[0-9]+|[0-9]+\\.[0-9]+|halb?|halbe?|einigen?|wenigen?|mehreren?)"

var (
	reWenigen  = core.MustRegExp("wenigen", "")
	reHalb     = core.MustRegExp("halb", "")
	reHalben   = core.MustRegExp("halben", "")
	reEinigen  = core.MustRegExp("einigen", "")
	reMehreren = core.MustRegExp("mehreren", "")
)

// ParseNumberPattern parses a NumberPattern match.
func ParseNumberPattern(match string) float64 {
	num := js.ToLowerCase(match)
	if v, ok := IntegerWordDictionary.Get(num); ok {
		return v
	} else if num == "ein" || num == "einer" || num == "einem" || num == "einen" || num == "eine" {
		return 1
	} else if reWenigen.Exec(num) != nil {
		return 2
	} else if reHalb.Exec(num) != nil || reHalben.Exec(num) != nil {
		return 0.5
	} else if reEinigen.Exec(num) != nil {
		return 3
	} else if reMehreren.Exec(num) != nil {
		return 7
	}

	return js.ParseFloat(num)
}

//-----------------------------

// YearPattern is YEAR_PATTERN.
const YearPattern = "(?:[0-9]{1,4}(?:\\s*[vn]\\.?\\s*(?:C(?:hr)?|(?:u\\.?|d\\.?(?:\\s*g\\.?)?)?\\s*Z)\\.?|\\s*(?:u\\.?|d\\.?(?:\\s*g\\.)?)\\s*Z\\.?)?)"

var (
	reYearV     = core.MustRegExp("v", "i")
	reYearN     = core.MustRegExp("n", "i")
	reYearZ     = core.MustRegExp("z", "i")
	reNonDigits = core.MustRegExp("[^0-9]+", "gi") // replace() resets its lastIndex
)

// ParseYear parses a YearPattern match.
func ParseYear(match string) float64 {
	if reYearV.Test(match) {
		// v.Chr.
		return -js.ParseInt(reNonDigits.Replace(match, ""), 0)
	}

	if reYearN.Test(match) {
		// n.Chr.
		return js.ParseInt(reNonDigits.Replace(match, ""), 0)
	}

	if reYearZ.Test(match) {
		// n.Chr. as "uZ" or "dgZ"
		return js.ParseInt(reNonDigits.Replace(match, ""), 0)
	}

	rawYearNumber := js.ParseInt(match, 0)
	return core.FindMostLikelyADYear(rawYearNumber)
}

//-----------------------------

var singleTimeUnitPattern = "(" + NumberPattern + ")\\s{0,5}(" + core.MatchAnyPattern(TimeUnitDictionary.Keys()) + ")\\s{0,5}"
var singleTimeUnitRegex = core.MustRegExp(singleTimeUnitPattern, "i")

// TimeUnitsPattern is TIME_UNITS_PATTERN.
var TimeUnitsPattern = core.RepeatedTimeunitPattern("", singleTimeUnitPattern, core.DefaultTimeunitConnectorPattern)

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
