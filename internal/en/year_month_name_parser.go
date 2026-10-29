// Ported from src/locales/en/parsers/ENYearMonthNameParser.ts.

package en

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

const yearMonthNameYearPattern = "(?:[1-9][0-9]{0,3}\\s{0,2}(?:BE|AD|BC|BCE|CE)|[1-9][0-9]{3})"

var yearMonthNamePattern = core.MustRegExp(
	"("+yearMonthNameYearPattern+")"+
		"(?:\\s*[-.\\/,]?\\s*|\\s+of\\s+)"+
		"("+core.MatchAnyPattern(MonthDictionary.Keys())+")"+
		"(?=[^\\s\\w]|\\s+[^0-9]|\\s+$|$)",
	"i",
)

const (
	yearMonthNameYearGroup      = 1
	yearMonthNameMonthNameGroup = 2
)

// YearMonthNameParser is ENYearMonthNameParser ("2012 January").
type YearMonthNameParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewYearMonthNameParser returns a new ENYearMonthNameParser.
func NewYearMonthNameParser() *YearMonthNameParser {
	p := &YearMonthNameParser{}
	p.Self = p
	return p
}

// InnerPattern returns the year-month-name pattern.
func (p *YearMonthNameParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return yearMonthNamePattern
}

// InnerExtract returns the result of the year and month.
func (p *YearMonthNameParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	year := ParseYear(match.Get(yearMonthNameYearGroup))
	monthName := js.ToLowerCase(match.Get(yearMonthNameMonthNameGroup))
	month := monthLookup(monthName)

	result := context.CreateParsingResult(match.Index, match.Get(0), nil, nil)
	result.Start.Imply(core.ComponentDay, 1)
	result.Start.Assign(core.ComponentMonth, month)
	result.Start.Assign(core.ComponentYear, year)
	result.Start.AddTag("parser/ENYearMonthNameParser")
	return result
}
