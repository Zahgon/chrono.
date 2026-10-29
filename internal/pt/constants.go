// Ported from src/locales/pt/constants.ts.

package pt

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

// WeekdayDictionary is WEEKDAY_DICTIONARY.
var WeekdayDictionary = core.NewDictionary([]core.Entry[float64]{
	{Key: "domingo", Value: 0},
	{Key: "dom", Value: 0},
	{Key: "segunda", Value: 1},
	{Key: "segunda-feira", Value: 1},
	{Key: "seg", Value: 1},
	{Key: "terça", Value: 2},
	{Key: "terça-feira", Value: 2},
	{Key: "ter", Value: 2},
	{Key: "quarta", Value: 3},
	{Key: "quarta-feira", Value: 3},
	{Key: "qua", Value: 3},
	{Key: "quinta", Value: 4},
	{Key: "quinta-feira", Value: 4},
	{Key: "qui", Value: 4},
	{Key: "sexta", Value: 5},
	{Key: "sexta-feira", Value: 5},
	{Key: "sex", Value: 5},
	{Key: "sábado", Value: 6},
	{Key: "sabado", Value: 6},
	{Key: "sab", Value: 6},
}...)

// MonthDictionary is MONTH_DICTIONARY.
var MonthDictionary = core.NewDictionary([]core.Entry[float64]{
	{Key: "janeiro", Value: 1},
	{Key: "jan", Value: 1},
	{Key: "jan.", Value: 1},
	{Key: "fevereiro", Value: 2},
	{Key: "fev", Value: 2},
	{Key: "fev.", Value: 2},
	{Key: "março", Value: 3},
	{Key: "mar", Value: 3},
	{Key: "mar.", Value: 3},
	{Key: "abril", Value: 4},
	{Key: "abr", Value: 4},
	{Key: "abr.", Value: 4},
	{Key: "maio", Value: 5},
	{Key: "mai", Value: 5},
	{Key: "mai.", Value: 5},
	{Key: "junho", Value: 6},
	{Key: "jun", Value: 6},
	{Key: "jun.", Value: 6},
	{Key: "julho", Value: 7},
	{Key: "jul", Value: 7},
	{Key: "jul.", Value: 7},
	{Key: "agosto", Value: 8},
	{Key: "ago", Value: 8},
	{Key: "ago.", Value: 8},
	{Key: "setembro", Value: 9},
	{Key: "set", Value: 9},
	{Key: "set.", Value: 9},
	{Key: "outubro", Value: 10},
	{Key: "out", Value: 10},
	{Key: "out.", Value: 10},
	{Key: "novembro", Value: 11},
	{Key: "nov", Value: 11},
	{Key: "nov.", Value: 11},
	{Key: "dezembro", Value: 12},
	{Key: "dez", Value: 12},
	{Key: "dez.", Value: 12},
}...)

// monthLookup is MONTH_DICTIONARY[name] (undefined when missing).
func monthLookup(name string) float64 {
	if v, ok := MonthDictionary.Get(name); ok {
		return v
	}
	return core.Undefined
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
