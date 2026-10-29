// Ported from src/locales/en/constants.ts.

package en

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

func concat[V any](parts ...[]core.Entry[V]) []core.Entry[V] {
	var out []core.Entry[V]
	for _, p := range parts {
		out = append(out, p...)
	}
	return out
}

// WeekdayDictionary is WEEKDAY_DICTIONARY.
var WeekdayDictionary = core.NewDictionary([]core.Entry[float64]{
	{Key: "sunday", Value: 0},
	{Key: "sun", Value: 0},
	{Key: "sun.", Value: 0},
	{Key: "monday", Value: 1},
	{Key: "mon", Value: 1},
	{Key: "mon.", Value: 1},
	{Key: "tuesday", Value: 2},
	{Key: "tue", Value: 2},
	{Key: "tue.", Value: 2},
	{Key: "wednesday", Value: 3},
	{Key: "wed", Value: 3},
	{Key: "wed.", Value: 3},
	{Key: "thursday", Value: 4},
	{Key: "thurs", Value: 4},
	{Key: "thurs.", Value: 4},
	{Key: "thur", Value: 4},
	{Key: "thur.", Value: 4},
	{Key: "thu", Value: 4},
	{Key: "thu.", Value: 4},
	{Key: "friday", Value: 5},
	{Key: "fri", Value: 5},
	{Key: "fri.", Value: 5},
	{Key: "saturday", Value: 6},
	{Key: "sat", Value: 6},
	{Key: "sat.", Value: 6},
}...)

// FullMonthNameDictionary is FULL_MONTH_NAME_DICTIONARY.
var FullMonthNameDictionary = core.NewDictionary([]core.Entry[float64]{
	{Key: "january", Value: 1},
	{Key: "february", Value: 2},
	{Key: "march", Value: 3},
	{Key: "april", Value: 4},
	{Key: "may", Value: 5},
	{Key: "june", Value: 6},
	{Key: "july", Value: 7},
	{Key: "august", Value: 8},
	{Key: "september", Value: 9},
	{Key: "october", Value: 10},
	{Key: "november", Value: 11},
	{Key: "december", Value: 12},
}...)

// MonthDictionary is MONTH_DICTIONARY.
var MonthDictionary = core.NewDictionary(concat(FullMonthNameDictionary.Entries(), []core.Entry[float64]{
	{Key: "jan", Value: 1},
	{Key: "jan.", Value: 1},
	{Key: "feb", Value: 2},
	{Key: "feb.", Value: 2},
	{Key: "mar", Value: 3},
	{Key: "mar.", Value: 3},
	{Key: "apr", Value: 4},
	{Key: "apr.", Value: 4},
	{Key: "jun", Value: 6},
	{Key: "jun.", Value: 6},
	{Key: "jul", Value: 7},
	{Key: "jul.", Value: 7},
	{Key: "aug", Value: 8},
	{Key: "aug.", Value: 8},
	{Key: "sep", Value: 9},
	{Key: "sep.", Value: 9},
	{Key: "sept", Value: 9},
	{Key: "sept.", Value: 9},
	{Key: "oct", Value: 10},
	{Key: "oct.", Value: 10},
	{Key: "nov", Value: 11},
	{Key: "nov.", Value: 11},
	{Key: "dec", Value: 12},
	{Key: "dec.", Value: 12},
})...)

// IntegerWordDictionary is INTEGER_WORD_DICTIONARY.
var IntegerWordDictionary = core.NewDictionary([]core.Entry[float64]{
	{Key: "one", Value: 1},
	{Key: "two", Value: 2},
	{Key: "three", Value: 3},
	{Key: "four", Value: 4},
	{Key: "five", Value: 5},
	{Key: "six", Value: 6},
	{Key: "seven", Value: 7},
	{Key: "eight", Value: 8},
	{Key: "nine", Value: 9},
	{Key: "ten", Value: 10},
	{Key: "eleven", Value: 11},
	{Key: "twelve", Value: 12},
}...)

// OrdinalWordDictionary is ORDINAL_WORD_DICTIONARY.
var OrdinalWordDictionary = core.NewDictionary([]core.Entry[float64]{
	{Key: "first", Value: 1},
	{Key: "second", Value: 2},
	{Key: "third", Value: 3},
	{Key: "fourth", Value: 4},
	{Key: "fifth", Value: 5},
	{Key: "sixth", Value: 6},
	{Key: "seventh", Value: 7},
	{Key: "eighth", Value: 8},
	{Key: "ninth", Value: 9},
	{Key: "tenth", Value: 10},
	{Key: "eleventh", Value: 11},
	{Key: "twelfth", Value: 12},
	{Key: "thirteenth", Value: 13},
	{Key: "fourteenth", Value: 14},
	{Key: "fifteenth", Value: 15},
	{Key: "sixteenth", Value: 16},
	{Key: "seventeenth", Value: 17},
	{Key: "eighteenth", Value: 18},
	{Key: "nineteenth", Value: 19},
	{Key: "twentieth", Value: 20},
	{Key: "twenty first", Value: 21},
	{Key: "twenty-first", Value: 21},
	{Key: "twenty second", Value: 22},
	{Key: "twenty-second", Value: 22},
	{Key: "twenty third", Value: 23},
	{Key: "twenty-third", Value: 23},
	{Key: "twenty fourth", Value: 24},
	{Key: "twenty-fourth", Value: 24},
	{Key: "twenty fifth", Value: 25},
	{Key: "twenty-fifth", Value: 25},
	{Key: "twenty sixth", Value: 26},
	{Key: "twenty-sixth", Value: 26},
	{Key: "twenty seventh", Value: 27},
	{Key: "twenty-seventh", Value: 27},
	{Key: "twenty eighth", Value: 28},
	{Key: "twenty-eighth", Value: 28},
	{Key: "twenty ninth", Value: 29},
	{Key: "twenty-ninth", Value: 29},
	{Key: "thirtieth", Value: 30},
	{Key: "thirty first", Value: 31},
	{Key: "thirty-first", Value: 31},
}...)

// TimeUnitDictionaryNoAbbr is TIME_UNIT_DICTIONARY_NO_ABBR.
var TimeUnitDictionaryNoAbbr = core.NewDictionary([]core.Entry[string]{
	{Key: "second", Value: "second"},
	{Key: "seconds", Value: "second"},
	{Key: "minute", Value: "minute"},
	{Key: "minutes", Value: "minute"},
	{Key: "hour", Value: "hour"},
	{Key: "hours", Value: "hour"},
	{Key: "day", Value: "day"},
	{Key: "days", Value: "day"},
	{Key: "week", Value: "week"},
	{Key: "weeks", Value: "week"},
	{Key: "month", Value: "month"},
	{Key: "months", Value: "month"},
	{Key: "quarter", Value: "quarter"},
	{Key: "quarters", Value: "quarter"},
	{Key: "year", Value: "year"},
	{Key: "years", Value: "year"},
}...)

// TimeUnitDictionary is TIME_UNIT_DICTIONARY.
var TimeUnitDictionary = core.NewDictionary(concat([]core.Entry[string]{
	{Key: "s", Value: "second"},
	{Key: "sec", Value: "second"},
	{Key: "second", Value: "second"},
	{Key: "seconds", Value: "second"},
	{Key: "m", Value: "minute"},
	{Key: "min", Value: "minute"},
	{Key: "mins", Value: "minute"},
	{Key: "minute", Value: "minute"},
	{Key: "minutes", Value: "minute"},
	{Key: "h", Value: "hour"},
	{Key: "hr", Value: "hour"},
	{Key: "hrs", Value: "hour"},
	{Key: "hour", Value: "hour"},
	{Key: "hours", Value: "hour"},
	{Key: "d", Value: "day"},
	{Key: "day", Value: "day"},
	{Key: "days", Value: "day"},
	{Key: "w", Value: "week"},
	{Key: "week", Value: "week"},
	{Key: "weeks", Value: "week"},
	{Key: "mo", Value: "month"},
	{Key: "mon", Value: "month"},
	{Key: "mos", Value: "month"},
	{Key: "month", Value: "month"},
	{Key: "months", Value: "month"},
	{Key: "qtr", Value: "quarter"},
	{Key: "quarter", Value: "quarter"},
	{Key: "quarters", Value: "quarter"},
	{Key: "y", Value: "year"},
	{Key: "yr", Value: "year"},
	{Key: "year", Value: "year"},
	{Key: "years", Value: "year"},
}, TimeUnitDictionaryNoAbbr.Entries())...)

//-----------------------------

// NumberPattern is NUMBER_PATTERN.
var NumberPattern = "(?:" + core.MatchAnyPattern(IntegerWordDictionary.Keys()) +
	"|[0-9]+|[0-9]+\\.[0-9]+|half(?:\\s{0,2}an?)?|an?\\b(?:\\s{0,2}few)?|few|several|the|a?\\s{0,2}couple\\s{0,2}(?:of)?)"

var (
	reFew     = core.MustRegExp("few", "")
	reHalf    = core.MustRegExp("half", "")
	reCouple  = core.MustRegExp("couple", "")
	reSeveral = core.MustRegExp("several", "")
)

// ParseNumberPattern parses a NumberPattern match.
func ParseNumberPattern(match string) float64 {
	num := js.ToLowerCase(match)
	if v, ok := IntegerWordDictionary.Get(num); ok {
		return v
	} else if num == "a" || num == "an" || num == "the" {
		return 1
	} else if reFew.Exec(num) != nil {
		return 3
	} else if reHalf.Exec(num) != nil {
		return 0.5
	} else if reCouple.Exec(num) != nil {
		return 2
	} else if reSeveral.Exec(num) != nil {
		return 7
	}
	return js.ParseFloat(num)
}

//-----------------------------

// OrdinalNumberPattern is ORDINAL_NUMBER_PATTERN.
var OrdinalNumberPattern = "(?:" + core.MatchAnyPattern(OrdinalWordDictionary.Keys()) + "|[0-9]{1,2}(?:st|nd|rd|th)?)"

var reOrdinalSuffix = core.MustRegExp("(?:st|nd|rd|th)$", "i")

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
const YearPattern = "(?:[1-9][0-9]{0,3}\\s{0,2}(?:BE|AD|BC|BCE|CE)|[1-9][0-9]{3}|[0-9]{2}(?!\\w|:\\d|\\s+(?:am|pm|o\\s*clock|jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec)))"

var (
	reBE   = core.MustRegExp("BE", "i")
	reBCE  = core.MustRegExp("BCE?", "i")
	reADCE = core.MustRegExp("(AD|CE)", "i")
)

// ParseYear parses a YearPattern match.
func ParseYear(match string) float64 {
	if reBE.Test(match) {
		// Buddhist Era
		match = reBE.Replace(match, "")
		return js.ParseInt(match, 0) - 543
	}

	if reBCE.Test(match) {
		// Before Christ, Before Common Era
		match = reBCE.Replace(match, "")
		return -js.ParseInt(match, 0)
	}

	if reADCE.Test(match) {
		// Anno Domini, Common Era
		match = reADCE.Replace(match, "")
		return js.ParseInt(match, 0)
	}

	rawYearNumber := js.ParseInt(match, 0)
	return core.FindMostLikelyADYear(rawYearNumber)
}

//-----------------------------

var singleTimeUnitPattern = "(" + NumberPattern + ")\\s{0,3}(" + core.MatchAnyPattern(TimeUnitDictionary.Keys()) + ")"
var singleTimeUnitRegex = core.MustRegExp(singleTimeUnitPattern, "i")

var singleTimeUnitNoAbbrPattern = "(" + NumberPattern + ")\\s{0,3}(" + core.MatchAnyPattern(TimeUnitDictionaryNoAbbr.Keys()) + ")"

const timeUnitConnectorPattern = "\\s{0,5},?(?:\\s*and)?\\s{0,5}"

// TimeUnitsPattern is TIME_UNITS_PATTERN.
var TimeUnitsPattern = core.RepeatedTimeunitPattern(
	"(?:(?:about|around)\\s{0,3})?",
	singleTimeUnitPattern,
	timeUnitConnectorPattern,
)

// TimeUnitsNoAbbrPattern is TIME_UNITS_NO_ABBR_PATTERN.
var TimeUnitsNoAbbrPattern = core.RepeatedTimeunitPattern(
	"(?:(?:about|around)\\s{0,3})?",
	singleTimeUnitNoAbbrPattern,
	timeUnitConnectorPattern,
)

// ParseDuration parses a TimeUnitsPattern match (nil when no time unit is found).
func ParseDuration(timeunitText string) *core.Duration {
	fragments := core.NewDuration()
	remainingText := timeunitText
	match := singleTimeUnitRegex.Exec(remainingText)
	for match != nil {
		collectDateTimeFragment(fragments, match)
		remainingText = js.Trim(js.SubstringFrom(remainingText, js.Len(match.Get(0))))
		match = singleTimeUnitRegex.Exec(remainingText)
	}
	if fragments.Len() == 0 {
		return nil
	}
	return fragments
}

var reLettersOnly = core.MustRegExp("^[a-zA-Z]+$", "")

func collectDateTimeFragment(fragments *core.Duration, match *core.RegExpMatch) {
	if reLettersOnly.Exec(match.Get(0)) != nil {
		return
	}
	num := ParseNumberPattern(match.Get(1))
	unit, ok := TimeUnitDictionary.Get(js.ToLowerCase(match.Get(2)))
	if !ok {
		unit = "undefined"
	}
	fragments.Set(unit, num)
}
