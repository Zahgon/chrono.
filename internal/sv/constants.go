// Ported from src/locales/sv/constants.ts.

package sv

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

// WeekdayDictionary is WEEKDAY_DICTIONARY.
var WeekdayDictionary = core.NewDictionary([]core.Entry[float64]{
	{Key: "söndag", Value: 0},
	{Key: "sön", Value: 0},
	{Key: "so", Value: 0},
	{Key: "måndag", Value: 1},
	{Key: "mån", Value: 1},
	{Key: "må", Value: 1},
	{Key: "tisdag", Value: 2},
	{Key: "tis", Value: 2},
	{Key: "ti", Value: 2},
	{Key: "onsdag", Value: 3},
	{Key: "ons", Value: 3},
	{Key: "on", Value: 3},
	{Key: "torsdag", Value: 4},
	{Key: "tors", Value: 4},
	{Key: "to", Value: 4},
	{Key: "fredag", Value: 5},
	{Key: "fre", Value: 5},
	{Key: "fr", Value: 5},
	{Key: "lördag", Value: 6},
	{Key: "lör", Value: 6},
	{Key: "lö", Value: 6},
}...)

// MonthDictionary is MONTH_DICTIONARY.
var MonthDictionary = core.NewDictionary([]core.Entry[float64]{
	{Key: "januari", Value: 1},
	{Key: "jan", Value: 1},
	{Key: "jan.", Value: 1},
	{Key: "februari", Value: 2},
	{Key: "feb", Value: 2},
	{Key: "feb.", Value: 2},
	{Key: "mars", Value: 3},
	{Key: "mar", Value: 3},
	{Key: "mar.", Value: 3},
	{Key: "april", Value: 4},
	{Key: "apr", Value: 4},
	{Key: "apr.", Value: 4},
	{Key: "maj", Value: 5},
	{Key: "juni", Value: 6},
	{Key: "jun", Value: 6},
	{Key: "jun.", Value: 6},
	{Key: "juli", Value: 7},
	{Key: "jul", Value: 7},
	{Key: "jul.", Value: 7},
	{Key: "augusti", Value: 8},
	{Key: "aug", Value: 8},
	{Key: "aug.", Value: 8},
	{Key: "september", Value: 9},
	{Key: "sep", Value: 9},
	{Key: "sep.", Value: 9},
	{Key: "sept", Value: 9},
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

// OrdinalNumberDictionary is ORDINAL_NUMBER_DICTIONARY.
var OrdinalNumberDictionary = core.NewDictionary([]core.Entry[float64]{
	{Key: "första", Value: 1},
	{Key: "andra", Value: 2},
	{Key: "tredje", Value: 3},
	{Key: "fjärde", Value: 4},
	{Key: "femte", Value: 5},
	{Key: "sjätte", Value: 6},
	{Key: "sjunde", Value: 7},
	{Key: "åttonde", Value: 8},
	{Key: "nionde", Value: 9},
	{Key: "tionde", Value: 10},
	{Key: "elfte", Value: 11},
	{Key: "tolfte", Value: 12},
	{Key: "trettonde", Value: 13},
	{Key: "fjortonde", Value: 14},
	{Key: "femtonde", Value: 15},
	{Key: "sextonde", Value: 16},
	{Key: "sjuttonde", Value: 17},
	{Key: "artonde", Value: 18},
	{Key: "nittonde", Value: 19},
	{Key: "tjugonde", Value: 20},
	{Key: "tjugoförsta", Value: 21},
	{Key: "tjugoandra", Value: 22},
	{Key: "tjugotredje", Value: 23},
	{Key: "tjugofjärde", Value: 24},
	{Key: "tjugofemte", Value: 25},
	{Key: "tjugosjätte", Value: 26},
	{Key: "tjugosjunde", Value: 27},
	{Key: "tjugoåttonde", Value: 28},
	{Key: "tjugonionde", Value: 29},
	{Key: "trettionde", Value: 30},
	{Key: "trettioförsta", Value: 31},
}...)

// IntegerWordDictionary is INTEGER_WORD_DICTIONARY. As in the original, the key for 30
// ("trettiо") ends with a Cyrillic "о" (U+043E).
var IntegerWordDictionary = core.NewDictionary([]core.Entry[float64]{
	{Key: "en", Value: 1},
	{Key: "ett", Value: 1},
	{Key: "två", Value: 2},
	{Key: "tre", Value: 3},
	{Key: "fyra", Value: 4},
	{Key: "fem", Value: 5},
	{Key: "sex", Value: 6},
	{Key: "sju", Value: 7},
	{Key: "åtta", Value: 8},
	{Key: "nio", Value: 9},
	{Key: "tio", Value: 10},
	{Key: "elva", Value: 11},
	{Key: "tolv", Value: 12},
	{Key: "tretton", Value: 13},
	{Key: "fjorton", Value: 14},
	{Key: "femton", Value: 15},
	{Key: "sexton", Value: 16},
	{Key: "sjutton", Value: 17},
	{Key: "arton", Value: 18},
	{Key: "nitton", Value: 19},
	{Key: "tjugo", Value: 20},
	{Key: "trettiо", Value: 30},
	{Key: "fyrtio", Value: 40},
	{Key: "femtio", Value: 50},
	{Key: "sextio", Value: 60},
	{Key: "sjuttio", Value: 70},
	{Key: "åttio", Value: 80},
	{Key: "nittio", Value: 90},
	{Key: "hundra", Value: 100},
	{Key: "tusen", Value: 1000},
}...)

// TimeUnitDictionary is TIME_UNIT_DICTIONARY. As in the original, the first "kvartаl" key has a
// Cyrillic "а" (U+0430).
var TimeUnitDictionary = core.NewDictionary([]core.Entry[string]{
	{Key: "sek", Value: "second"},
	{Key: "sekund", Value: "second"},
	{Key: "sekunder", Value: "second"},
	{Key: "min", Value: "minute"},
	{Key: "minut", Value: "minute"},
	{Key: "minuter", Value: "minute"},
	{Key: "tim", Value: "hour"},
	{Key: "timme", Value: "hour"},
	{Key: "timmar", Value: "hour"},
	{Key: "dag", Value: "day"},
	{Key: "dagar", Value: "day"},
	{Key: "vecka", Value: "week"},
	{Key: "veckor", Value: "week"},
	{Key: "mån", Value: "month"},
	{Key: "månad", Value: "month"},
	{Key: "månader", Value: "month"},
	{Key: "år", Value: "year"},
	{Key: "kvartаl", Value: "quarter"},
	{Key: "kvartal", Value: "quarter"},
}...)

// TimeUnitNoAbbrDictionary is TIME_UNIT_NO_ABBR_DICTIONARY.
var TimeUnitNoAbbrDictionary = core.NewDictionary([]core.Entry[string]{
	{Key: "sekund", Value: "second"},
	{Key: "sekunder", Value: "second"},
	{Key: "minut", Value: "minute"},
	{Key: "minuter", Value: "minute"},
	{Key: "timme", Value: "hour"},
	{Key: "timmar", Value: "hour"},
	{Key: "dag", Value: "day"},
	{Key: "dagar", Value: "day"},
	{Key: "vecka", Value: "week"},
	{Key: "veckor", Value: "week"},
	{Key: "månad", Value: "month"},
	{Key: "månader", Value: "month"},
	{Key: "år", Value: "year"},
	{Key: "kvartal", Value: "quarter"},
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

// NumberPattern is NUMBER_PATTERN.
var NumberPattern = "(?:" + core.MatchAnyPattern(IntegerWordDictionary.Keys()) + "|\\d+)"

// OrdinalNumberPattern is ORDINAL_NUMBER_PATTERN.
var OrdinalNumberPattern = "(?:" + core.MatchAnyPattern(OrdinalNumberDictionary.Keys()) + "|\\d{1,2}(?:e|:e))"

// TimeUnitPattern is TIME_UNIT_PATTERN.
var TimeUnitPattern = "(?:" + core.MatchAnyPattern(TimeUnitDictionary.Keys()) + ")"

var singleTimeUnitPattern = "(" + NumberPattern + ")\\s{0,5}(" + core.MatchAnyPattern(TimeUnitDictionary.Keys()) + ")\\s{0,5}"
var singleTimeUnitRegex = core.MustRegExp(singleTimeUnitPattern, "i")

var singleTimeUnitNoAbbrPattern = "(" + NumberPattern + ")\\s{0,5}(" + core.MatchAnyPattern(
	TimeUnitNoAbbrDictionary.Keys(),
) + ")\\s{0,5}"

// TimeUnitsPattern is TIME_UNITS_PATTERN.
var TimeUnitsPattern = core.RepeatedTimeunitPattern("", singleTimeUnitPattern, core.DefaultTimeunitConnectorPattern)

// TimeUnitsNoAbbrPattern is TIME_UNITS_NO_ABBR_PATTERN.
var TimeUnitsNoAbbrPattern = core.RepeatedTimeunitPattern("", singleTimeUnitNoAbbrPattern, core.DefaultTimeunitConnectorPattern)

// ParseNumberPattern parses a NumberPattern match.
func ParseNumberPattern(match string) float64 {
	num := js.ToLowerCase(match)
	if v, ok := IntegerWordDictionary.Get(num); ok {
		return v
	}
	return js.ParseInt(num, 0)
}

// ParseOrdinalNumberPattern parses an OrdinalNumberPattern match.
func ParseOrdinalNumberPattern(match string) float64 {
	num := js.ToLowerCase(match)
	if v, ok := OrdinalNumberDictionary.Get(num); ok {
		return v
	}
	return js.ParseInt(num, 0)
}

var reDigits = core.MustRegExp("\\d+", "")

// ParseYear parses a year: digits (two-digit years mapped to the most likely AD year) or an
// integer word.
func ParseYear(match string) float64 {
	if reDigits.Test(match) {
		yearNumber := js.ParseInt(match, 0)
		if yearNumber < 100 {
			yearNumber = core.FindMostLikelyADYear(yearNumber)
		}
		return yearNumber
	}

	num := js.ToLowerCase(match)
	if v, ok := IntegerWordDictionary.Get(num); ok {
		return v
	}

	return js.ParseInt(match, 0)
}
