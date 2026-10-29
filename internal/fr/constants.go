// Ported from src/locales/fr/constants.ts.

package fr

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

// WeekdayDictionary is WEEKDAY_DICTIONARY.
var WeekdayDictionary = core.NewDictionary([]core.Entry[float64]{
	{Key: "dimanche", Value: 0},
	{Key: "dim", Value: 0},
	{Key: "lundi", Value: 1},
	{Key: "lun", Value: 1},
	{Key: "mardi", Value: 2},
	{Key: "mar", Value: 2},
	{Key: "mercredi", Value: 3},
	{Key: "mer", Value: 3},
	{Key: "jeudi", Value: 4},
	{Key: "jeu", Value: 4},
	{Key: "vendredi", Value: 5},
	{Key: "ven", Value: 5},
	{Key: "samedi", Value: 6},
	{Key: "sam", Value: 6},
}...)

// MonthDictionary is MONTH_DICTIONARY.
var MonthDictionary = core.NewDictionary([]core.Entry[float64]{
	{Key: "janvier", Value: 1},
	{Key: "jan", Value: 1},
	{Key: "jan.", Value: 1},
	{Key: "janv", Value: 1},
	{Key: "janv.", Value: 1},
	{Key: "février", Value: 2},
	{Key: "fév", Value: 2},
	{Key: "fév.", Value: 2},
	{Key: "févr", Value: 2},
	{Key: "févr.", Value: 2},
	{Key: "fevrier", Value: 2},
	{Key: "fev", Value: 2},
	{Key: "fev.", Value: 2},
	{Key: "fevr", Value: 2},
	{Key: "fevr.", Value: 2},
	{Key: "mars", Value: 3},
	{Key: "mar", Value: 3},
	{Key: "mar.", Value: 3},
	{Key: "avril", Value: 4},
	{Key: "avr", Value: 4},
	{Key: "avr.", Value: 4},
	{Key: "mai", Value: 5},
	{Key: "juin", Value: 6},
	{Key: "juin.", Value: 6},
	{Key: "jun", Value: 6},
	{Key: "jun.", Value: 6},
	{Key: "juillet", Value: 7},
	{Key: "juil", Value: 7},
	{Key: "juil.", Value: 7},
	{Key: "jul", Value: 7},
	{Key: "jul.", Value: 7},
	{Key: "août", Value: 8},
	{Key: "août.", Value: 8},
	{Key: "aout", Value: 8},
	{Key: "aout.", Value: 8},
	{Key: "septembre", Value: 9},
	{Key: "sep", Value: 9},
	{Key: "sep.", Value: 9},
	{Key: "sept", Value: 9},
	{Key: "sept.", Value: 9},
	{Key: "octobre", Value: 10},
	{Key: "oct", Value: 10},
	{Key: "oct.", Value: 10},
	{Key: "novembre", Value: 11},
	{Key: "nov", Value: 11},
	{Key: "nov.", Value: 11},
	{Key: "décembre", Value: 12},
	{Key: "decembre", Value: 12},
	{Key: "dec", Value: 12},
	{Key: "dec.", Value: 12},
	{Key: "déc", Value: 12},
	{Key: "déc.", Value: 12},
}...)

// IntegerWordDictionary is INTEGER_WORD_DICTIONARY.
var IntegerWordDictionary = core.NewDictionary([]core.Entry[float64]{
	{Key: "un", Value: 1},
	{Key: "deux", Value: 2},
	{Key: "trois", Value: 3},
	{Key: "quatre", Value: 4},
	{Key: "cinq", Value: 5},
	{Key: "six", Value: 6},
	{Key: "sept", Value: 7},
	{Key: "huit", Value: 8},
	{Key: "neuf", Value: 9},
	{Key: "dix", Value: 10},
	{Key: "onze", Value: 11},
	{Key: "douze", Value: 12},
	{Key: "treize", Value: 13},
}...)

// TimeUnitDictionary is TIME_UNIT_DICTIONARY.
var TimeUnitDictionary = core.NewDictionary([]core.Entry[string]{
	{Key: "sec", Value: "second"},
	{Key: "seconde", Value: "second"},
	{Key: "secondes", Value: "second"},
	{Key: "min", Value: "minute"},
	{Key: "mins", Value: "minute"},
	{Key: "minute", Value: "minute"},
	{Key: "minutes", Value: "minute"},
	{Key: "h", Value: "hour"},
	{Key: "hr", Value: "hour"},
	{Key: "hrs", Value: "hour"},
	{Key: "heure", Value: "hour"},
	{Key: "heures", Value: "hour"},
	{Key: "jour", Value: "day"},
	{Key: "jours", Value: "day"},
	{Key: "semaine", Value: "week"},
	{Key: "semaines", Value: "week"},
	{Key: "mois", Value: "month"},
	{Key: "trimestre", Value: "quarter"},
	{Key: "trimestres", Value: "quarter"},
	{Key: "ans", Value: "year"},
	{Key: "année", Value: "year"},
	{Key: "années", Value: "year"},
}...)

// monthLookup is MONTH_DICTIONARY[name] (undefined when missing).
func monthLookup(name string) float64 {
	if v, ok := MonthDictionary.Get(name); ok {
		return v
	}
	return core.Undefined
}

// timeunitLookup is TIME_UNIT_DICTIONARY[word] used as an object key ("undefined" when missing).
func timeunitLookup(word string) string {
	if v, ok := TimeUnitDictionary.Get(word); ok {
		return v
	}
	return "undefined"
}

//-----------------------------

// NumberPattern is NUMBER_PATTERN.
var NumberPattern = "(?:" + core.MatchAnyPattern(IntegerWordDictionary.Keys()) +
	"|[0-9]+|[0-9]+\\.[0-9]+|une?\\b|quelques?|demi-?)"

var (
	reQuelques = core.MustRegExp("quelques?", "")
	reDemi     = core.MustRegExp("demi-?", "")
)

// ParseNumberPattern parses a NumberPattern match.
func ParseNumberPattern(match string) float64 {
	num := js.ToLowerCase(match)
	if v, ok := IntegerWordDictionary.Get(num); ok {
		return v
	} else if num == "une" || num == "un" {
		return 1
	} else if reQuelques.Exec(num) != nil {
		return 3
	} else if reDemi.Exec(num) != nil {
		return 0.5
	}
	return js.ParseFloat(num)
}

//-----------------------------

// OrdinalNumberPattern is ORDINAL_NUMBER_PATTERN.
const OrdinalNumberPattern = "(?:[0-9]{1,2}(?:er)?)"

var reOrdinalSuffix = core.MustRegExp("(?:er)$", "i")

// ParseOrdinalNumberPattern parses an OrdinalNumberPattern match.
func ParseOrdinalNumberPattern(match string) float64 {
	num := js.ToLowerCase(match)
	num = reOrdinalSuffix.Replace(num, "")
	return js.ParseInt(num, 0)
}

//-----------------------------
// 88 p. Chr. n.
// 234 AC

// YearPattern is YEAR_PATTERN.
const YearPattern = "(?:[1-9][0-9]{0,3}\\s*(?:AC|AD|p\\.\\s*C(?:hr?)?\\.\\s*n\\.)|[1-2][0-9]{3}|[5-9][0-9])"

var (
	reAC       = core.MustRegExp("AC", "i")
	reBC       = core.MustRegExp("BC", "i")
	reAD       = core.MustRegExp("AD", "i")
	reC        = core.MustRegExp("C", "i")
	reNonDigit = core.MustRegExp("[^\\d]+", "i")
)

// ParseYear parses a YearPattern match.
func ParseYear(match string) float64 {
	if reAC.Test(match) {
		match = reBC.Replace(match, "")
		return -js.ParseInt(match, 0)
	}

	if reAD.Test(match) || reC.Test(match) {
		match = reNonDigit.Replace(match, "")
		return js.ParseInt(match, 0)
	}

	yearNumber := js.ParseInt(match, 0)
	if yearNumber < 100 {
		if yearNumber > 50 {
			yearNumber = yearNumber + 1900
		} else {
			yearNumber = yearNumber + 2000
		}
	}
	return yearNumber
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
	unit := timeunitLookup(js.ToLowerCase(match.Get(2)))
	fragments.Set(unit, num)
}
