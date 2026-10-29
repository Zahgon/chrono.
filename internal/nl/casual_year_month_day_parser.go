// Ported from src/locales/nl/parsers/NLCasualYearMonthDayParser.ts.

package nl

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

// casualYearMonthDayPattern is a date format with "/" (or "." or a space) between numbers like
// SlashDateFormatParser, but with the year before the month and the date: YYYY/MM/DD,
// YYYY-MM-DD, YYYY.MM.DD.
var casualYearMonthDayPattern = core.MustRegExp(
	"([0-9]{4})[\\.\\/\\s]"+
		"(?:("+core.MatchAnyPattern(MonthDictionary.Keys())+")|([0-9]{1,2}))[\\.\\/\\s]"+
		"([0-9]{1,2})"+
		"(?=\\W|$)",
	"i",
)

const (
	ymdYearNumberGroup  = 1
	ymdMonthNameGroup   = 2
	ymdMonthNumberGroup = 3
	ymdDateNumberGroup  = 4
)

// CasualYearMonthDayParser is NLCasualYearMonthDayParser ("2012/08/10", "2012 aug 10").
type CasualYearMonthDayParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewCasualYearMonthDayParser returns a new NLCasualYearMonthDayParser.
func NewCasualYearMonthDayParser() *CasualYearMonthDayParser {
	p := &CasualYearMonthDayParser{}
	p.Self = p
	return p
}

// InnerPattern returns the year-month-day pattern.
func (p *CasualYearMonthDayParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return casualYearMonthDayPattern
}

// InnerExtract returns the date components (nil for an invalid month).
func (p *CasualYearMonthDayParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	var month float64
	if match.Truthy(ymdMonthNumberGroup) {
		month = js.ParseInt(match.Get(ymdMonthNumberGroup), 0)
	} else {
		month = monthLookup(js.ToLowerCase(match.Get(ymdMonthNameGroup)))
	}

	if core.Num(month) < 1 || core.Num(month) > 12 {
		return nil
	}

	year := js.ParseInt(match.Get(ymdYearNumberGroup), 0)
	day := js.ParseInt(match.Get(ymdDateNumberGroup), 0)

	return core.Components{
		{Component: core.ComponentDay, Value: day},
		{Component: core.ComponentMonth, Value: month},
		{Component: core.ComponentYear, Value: year},
	}
}
