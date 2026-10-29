// Ported from src/locales/fi/constants.ts.

package fi

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

// WeekdayDictionary is WEEKDAY_DICTIONARY.
var WeekdayDictionary = core.NewDictionary([]core.Entry[float64]{
	{Key: "sunnuntai", Value: 0},
	{Key: "sunnuntaina", Value: 0},
	{Key: "su", Value: 0},
	{Key: "maanantai", Value: 1},
	{Key: "maanantaina", Value: 1},
	{Key: "ma", Value: 1},
	{Key: "tiistai", Value: 2},
	{Key: "tiistaina", Value: 2},
	{Key: "ti", Value: 2},
	{Key: "keskiviikko", Value: 3},
	{Key: "keskiviikkona", Value: 3},
	{Key: "ke", Value: 3},
	{Key: "torstai", Value: 4},
	{Key: "torstaina", Value: 4},
	{Key: "to", Value: 4},
	{Key: "perjantai", Value: 5},
	{Key: "perjantaina", Value: 5},
	{Key: "pe", Value: 5},
	{Key: "lauantai", Value: 6},
	{Key: "lauantaina", Value: 6},
	{Key: "la", Value: 6},
}...)

// MonthDictionary is MONTH_DICTIONARY.
var MonthDictionary = core.NewDictionary([]core.Entry[float64]{
	{Key: "tammikuu", Value: 1},
	{Key: "tammikuuta", Value: 1},
	{Key: "tammikuun", Value: 1},
	{Key: "tammi", Value: 1},
	{Key: "helmikuu", Value: 2},
	{Key: "helmikuuta", Value: 2},
	{Key: "helmikuun", Value: 2},
	{Key: "helmi", Value: 2},
	{Key: "maaliskuu", Value: 3},
	{Key: "maaliskuuta", Value: 3},
	{Key: "maaliskuun", Value: 3},
	{Key: "maalis", Value: 3},
	{Key: "huhtikuu", Value: 4},
	{Key: "huhtikuuta", Value: 4},
	{Key: "huhtikuun", Value: 4},
	{Key: "huhti", Value: 4},
	{Key: "toukokuu", Value: 5},
	{Key: "toukokuuta", Value: 5},
	{Key: "toukokuun", Value: 5},
	{Key: "touko", Value: 5},
	{Key: "kesäkuu", Value: 6},
	{Key: "kesäkuuta", Value: 6},
	{Key: "kesäkuun", Value: 6},
	{Key: "kesä", Value: 6},
	{Key: "heinäkuu", Value: 7},
	{Key: "heinäkuuta", Value: 7},
	{Key: "heinäkuun", Value: 7},
	{Key: "heinä", Value: 7},
	{Key: "elokuu", Value: 8},
	{Key: "elokuuta", Value: 8},
	{Key: "elokuun", Value: 8},
	{Key: "elo", Value: 8},
	{Key: "syyskuu", Value: 9},
	{Key: "syyskuuta", Value: 9},
	{Key: "syyskuun", Value: 9},
	{Key: "syys", Value: 9},
	{Key: "lokakuu", Value: 10},
	{Key: "lokakuuta", Value: 10},
	{Key: "lokakuun", Value: 10},
	{Key: "loka", Value: 10},
	{Key: "marraskuu", Value: 11},
	{Key: "marraskuuta", Value: 11},
	{Key: "marraskuun", Value: 11},
	{Key: "marras", Value: 11},
	{Key: "joulukuu", Value: 12},
	{Key: "joulukuuta", Value: 12},
	{Key: "joulukuun", Value: 12},
	{Key: "joulu", Value: 12},
}...)

// IntegerWordDictionary is INTEGER_WORD_DICTIONARY.
var IntegerWordDictionary = core.NewDictionary([]core.Entry[float64]{
	{Key: "yksi", Value: 1},
	{Key: "yhden", Value: 1},
	{Key: "kaksi", Value: 2},
	{Key: "kahden", Value: 2},
	{Key: "kolme", Value: 3},
	{Key: "kolmen", Value: 3},
	{Key: "neljä", Value: 4},
	{Key: "neljän", Value: 4},
	{Key: "viisi", Value: 5},
	{Key: "viiden", Value: 5},
	{Key: "kuusi", Value: 6},
	{Key: "kuuden", Value: 6},
	{Key: "seitsemän", Value: 7},
	{Key: "kahdeksan", Value: 8},
	{Key: "yhdeksän", Value: 9},
	{Key: "kymmenen", Value: 10},
}...)

// TimeUnitDictionary is TIME_UNIT_DICTIONARY.
var TimeUnitDictionary = core.NewDictionary([]core.Entry[string]{
	{Key: "s", Value: "second"},
	{Key: "sek", Value: "second"},
	{Key: "sekunti", Value: "second"},
	{Key: "sekuntia", Value: "second"},
	{Key: "sekunnin", Value: "second"},
	{Key: "min", Value: "minute"},
	{Key: "minuutti", Value: "minute"},
	{Key: "minuuttia", Value: "minute"},
	{Key: "minuutin", Value: "minute"},
	{Key: "t", Value: "hour"},
	{Key: "tunti", Value: "hour"},
	{Key: "tuntia", Value: "hour"},
	{Key: "tunnin", Value: "hour"},
	{Key: "pv", Value: "day"},
	{Key: "päivä", Value: "day"},
	{Key: "päivää", Value: "day"},
	{Key: "päivän", Value: "day"},
	{Key: "vk", Value: "week"},
	{Key: "viikko", Value: "week"},
	{Key: "viikkoa", Value: "week"},
	{Key: "viikon", Value: "week"},
	{Key: "kk", Value: "month"},
	{Key: "kuukausi", Value: "month"},
	{Key: "kuukautta", Value: "month"},
	{Key: "kuukauden", Value: "month"},
	{Key: "vuosi", Value: "year"},
	{Key: "vuotta", Value: "year"},
	{Key: "vuoden", Value: "year"},
}...)

// TimeUnitNoAbbrDictionary is TIME_UNIT_NO_ABBR_DICTIONARY.
var TimeUnitNoAbbrDictionary = core.NewDictionary([]core.Entry[string]{
	{Key: "sekunti", Value: "second"},
	{Key: "sekuntia", Value: "second"},
	{Key: "sekunnin", Value: "second"},
	{Key: "minuutti", Value: "minute"},
	{Key: "minuuttia", Value: "minute"},
	{Key: "minuutin", Value: "minute"},
	{Key: "tunti", Value: "hour"},
	{Key: "tuntia", Value: "hour"},
	{Key: "tunnin", Value: "hour"},
	{Key: "päivä", Value: "day"},
	{Key: "päivää", Value: "day"},
	{Key: "päivän", Value: "day"},
	{Key: "viikko", Value: "week"},
	{Key: "viikkoa", Value: "week"},
	{Key: "viikon", Value: "week"},
	{Key: "kuukausi", Value: "month"},
	{Key: "kuukautta", Value: "month"},
	{Key: "kuukauden", Value: "month"},
	{Key: "vuosi", Value: "year"},
	{Key: "vuotta", Value: "year"},
	{Key: "vuoden", Value: "year"},
}...)

// ParseDuration is parseDuration: the duration of a TimeUnitsPattern match (never nil).
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
	unit, ok := TimeUnitDictionary.Get(js.ToLowerCase(match.Get(2)))
	if !ok {
		// fragments[undefined] = num
		unit = "undefined"
	}
	fragments.Set(unit, num)
}

// NumberPattern is NUMBER_PATTERN.
var NumberPattern = "(?:" + core.MatchAnyPattern(IntegerWordDictionary.Keys()) + "|\\d+)"

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

// ParseNumberPattern is parseNumberPattern: a NumberPattern match as a number.
func ParseNumberPattern(match string) float64 {
	num := js.ToLowerCase(match)
	if v, ok := IntegerWordDictionary.Get(num); ok {
		return v
	}
	return js.ParseInt(num, 0)
}

var reDigits = core.MustRegExp("\\d+", "")

// ParseYear is parseYear: a year number or word.
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
