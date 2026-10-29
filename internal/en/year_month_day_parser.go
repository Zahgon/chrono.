// Ported from src/locales/en/parsers/ENYearMonthDayParser.ts.

package en

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

// Date format with slash "/" between numbers like ENSlashDateFormatParser, but this parser
// expects the year before the month and date: YYYY/MM/DD, YYYY-MM-DD, YYYY.MM.DD.
var yearMonthDayPattern = core.MustRegExp(
	"([0-9]{4})[-\\.\\/\\s]"+
		"(?:("+core.MatchAnyPattern(MonthDictionary.Keys())+")|([0-9]{1,2}))[-\\.\\/\\s]"+
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

// YearMonthDayParser is ENYearMonthDayParser.
type YearMonthDayParser struct {
	core.AbstractParserWithWordBoundaryChecking
	strictMonthDateOrder bool
}

// NewYearMonthDayParser returns a new ENYearMonthDayParser.
func NewYearMonthDayParser(strictMonthDateOrder bool) *YearMonthDayParser {
	p := &YearMonthDayParser{strictMonthDateOrder: strictMonthDateOrder}
	p.Self = p
	return p
}

// InnerPattern returns the year-month-day pattern.
func (p *YearMonthDayParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return yearMonthDayPattern
}

// InnerExtract returns the { day, month, year } object.
func (p *YearMonthDayParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	year := js.ParseInt(match.Get(ymdYearNumberGroup), 0)
	day := js.ParseInt(match.Get(ymdDateNumberGroup), 0)
	var month float64
	if match.Truthy(ymdMonthNumberGroup) {
		month = js.ParseInt(match.Get(ymdMonthNumberGroup), 0)
	} else {
		month = monthLookup(js.ToLowerCase(match.Get(ymdMonthNameGroup)))
	}

	if core.Num(month) < 1 || core.Num(month) > 12 {
		if p.strictMonthDateOrder {
			return nil
		}
		if day >= 1 && day <= 12 {
			month, day = day, month
		}
	}
	if core.Num(day) < 1 || core.Num(day) > 31 {
		return nil
	}

	return core.Components{
		{Component: core.ComponentDay, Value: day},
		{Component: core.ComponentMonth, Value: month},
		{Component: core.ComponentYear, Value: year},
	}
}
