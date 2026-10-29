// Ported from src/locales/vi/parsers/VIMonthYearParser.ts.

package vi

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

// tháng 3 năm 1975  |  tháng ba năm 1975  |  tháng chạp  |  tháng giêng/1975
var monthYearPattern = core.MustRegExp(
	"("+core.MatchAnyPattern(MonthDictionary.Keys())+")"+"(?:\\s*(?:năm|/)\\s*("+YearPattern+"))?"+"(?=\\W|$)",
	"i",
)

const (
	monthYearMonthGroup = 1
	monthYearYearGroup  = 2
)

// MonthYearParser is VIMonthYearParser ("tháng 3 năm 1975", "tháng chạp").
type MonthYearParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewMonthYearParser returns a new VIMonthYearParser.
func NewMonthYearParser() *MonthYearParser {
	p := &MonthYearParser{}
	p.Self = p
	return p
}

// InnerPattern returns the month (and year) pattern.
func (p *MonthYearParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return monthYearPattern
}

// InnerExtract returns the result of the month.
func (p *MonthYearParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	month, ok := MonthDictionary.Get(js.ToLowerCase(match.Get(monthYearMonthGroup)))
	if !ok || !core.Truthy(month) {
		return nil
	}

	result := context.CreateParsingResult(match.Index, match.Get(0), nil, nil)
	result.Start.Assign(core.ComponentMonth, month)
	result.Start.Imply(core.ComponentDay, 1)
	if match.Truthy(monthYearYearGroup) {
		result.Start.Assign(core.ComponentYear, ParseYear(match.Get(monthYearYearGroup)))
	} else {
		result.Start.Imply(core.ComponentYear, context.Reference.DateWithAdjustedTimezone().GetFullYear())
	}
	return result
}
