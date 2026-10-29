// Ported from src/locales/es/constants.ts.

package es

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

// WeekdayDictionary is WEEKDAY_DICTIONARY.
var WeekdayDictionary = core.NewDictionary([]core.Entry[float64]{
	{Key: "domingo", Value: 0},
	{Key: "dom", Value: 0},
	{Key: "lunes", Value: 1},
	{Key: "lun", Value: 1},
	{Key: "martes", Value: 2},
	{Key: "mar", Value: 2},
	{Key: "miércoles", Value: 3},
	{Key: "miercoles", Value: 3},
	{Key: "mié", Value: 3},
	{Key: "mie", Value: 3},
	{Key: "jueves", Value: 4},
	{Key: "jue", Value: 4},
	{Key: "viernes", Value: 5},
	{Key: "vie", Value: 5},
	{Key: "sábado", Value: 6},
	{Key: "sabado", Value: 6},
	{Key: "sáb", Value: 6},
	{Key: "sab", Value: 6},
}...)

// MonthDictionary is MONTH_DICTIONARY.
var MonthDictionary = core.NewDictionary([]core.Entry[float64]{
	{Key: "enero", Value: 1},
	{Key: "ene", Value: 1},
	{Key: "ene.", Value: 1},
	{Key: "febrero", Value: 2},
	{Key: "feb", Value: 2},
	{Key: "feb.", Value: 2},
	{Key: "marzo", Value: 3},
	{Key: "mar", Value: 3},
	{Key: "mar.", Value: 3},
	{Key: "abril", Value: 4},
	{Key: "abr", Value: 4},
	{Key: "abr.", Value: 4},
	{Key: "mayo", Value: 5},
	{Key: "may", Value: 5},
	{Key: "may.", Value: 5},
	{Key: "junio", Value: 6},
	{Key: "jun", Value: 6},
	{Key: "jun.", Value: 6},
	{Key: "julio", Value: 7},
	{Key: "jul", Value: 7},
	{Key: "jul.", Value: 7},
	{Key: "agosto", Value: 8},
	{Key: "ago", Value: 8},
	{Key: "ago.", Value: 8},
	{Key: "septiembre", Value: 9},
	{Key: "setiembre", Value: 9},
	{Key: "sep", Value: 9},
	{Key: "sep.", Value: 9},
	{Key: "octubre", Value: 10},
	{Key: "oct", Value: 10},
	{Key: "oct.", Value: 10},
	{Key: "noviembre", Value: 11},
	{Key: "nov", Value: 11},
	{Key: "nov.", Value: 11},
	{Key: "diciembre", Value: 12},
	{Key: "dic", Value: 12},
	{Key: "dic.", Value: 12},
}...)

// IntegerWordDictionary is INTEGER_WORD_DICTIONARY.
var IntegerWordDictionary = core.NewDictionary([]core.Entry[float64]{
	{Key: "uno", Value: 1},
	{Key: "dos", Value: 2},
	{Key: "tres", Value: 3},
	{Key: "cuatro", Value: 4},
	{Key: "cinco", Value: 5},
	{Key: "seis", Value: 6},
	{Key: "siete", Value: 7},
	{Key: "ocho", Value: 8},
	{Key: "nueve", Value: 9},
	{Key: "diez", Value: 10},
	{Key: "once", Value: 11},
	{Key: "doce", Value: 12},
	{Key: "trece", Value: 13},
}...)

// TimeUnitDictionary is TIME_UNIT_DICTIONARY.
var TimeUnitDictionary = core.NewDictionary([]core.Entry[string]{
	{Key: "sec", Value: "second"},
	{Key: "segundo", Value: "second"},
	{Key: "segundos", Value: "second"},
	{Key: "min", Value: "minute"},
	{Key: "mins", Value: "minute"},
	{Key: "minuto", Value: "minute"},
	{Key: "minutos", Value: "minute"},
	{Key: "h", Value: "hour"},
	{Key: "hr", Value: "hour"},
	{Key: "hrs", Value: "hour"},
	{Key: "hora", Value: "hour"},
	{Key: "horas", Value: "hour"},
	{Key: "día", Value: "day"},
	{Key: "días", Value: "day"},
	{Key: "semana", Value: "week"},
	{Key: "semanas", Value: "week"},
	{Key: "mes", Value: "month"},
	{Key: "meses", Value: "month"},
	{Key: "cuarto", Value: "quarter"},
	{Key: "cuartos", Value: "quarter"},
	{Key: "año", Value: "year"},
	{Key: "años", Value: "year"},
}...)

// monthLookup is MONTH_DICTIONARY[name] (undefined when missing).
func monthLookup(name string) float64 {
	if v, ok := MonthDictionary.Get(name); ok {
		return v
	}
	return core.Undefined
}

//-----------------------------

// NumberPattern is NUMBER_PATTERN.
var NumberPattern = "(?:" + core.MatchAnyPattern(IntegerWordDictionary.Keys()) +
	"|[0-9]+|[0-9]+\\.[0-9]+|un?|uno?|una?|algunos?|unos?|demi-?)"

var (
	reAlgunos = core.MustRegExp("algunos?", "")
	reUnos    = core.MustRegExp("unos?", "")
	reMedia   = core.MustRegExp("media?", "")
)

// ParseNumberPattern parses a NumberPattern match.
func ParseNumberPattern(match string) float64 {
	num := js.ToLowerCase(match)
	if v, ok := IntegerWordDictionary.Get(num); ok {
		return v
	} else if num == "un" || num == "una" || num == "uno" {
		return 1
	} else if reAlgunos.Exec(num) != nil {
		return 3
	} else if reUnos.Exec(num) != nil {
		return 3
	} else if reMedia.Exec(num) != nil {
		return 0.5
	}
	return js.ParseFloat(num)
}

//-----------------------------
// 88 p. Chr. n.
// 234 AC

// YearPattern is YEAR_PATTERN.
const YearPattern = "[0-9]{1,4}(?![^\\s]\\d)(?:\\s*[a|d]\\.?\\s*c\\.?|\\s*a\\.?\\s*d\\.?)?"

var (
	reYearDigits = core.MustRegExp("^[0-9]{1,4}$", "")
	reAC         = core.MustRegExp("a\\.?\\s*c\\.?", "i")
)

// ParseYear parses a YearPattern match.
func ParseYear(match string) float64 {
	if reYearDigits.Exec(match) != nil {
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

	if reAC.Exec(match) != nil {
		match = reAC.Replace(match, "")
		return -js.ParseInt(match, 0)
	}

	return js.ParseInt(match, 0)
}

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
	unit, ok := TimeUnitDictionary.Get(js.ToLowerCase(match.Get(2)))
	if !ok {
		unit = "undefined"
	}
	fragments.Set(unit, num)
}
