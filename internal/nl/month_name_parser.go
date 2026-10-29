// Ported from src/locales/nl/parsers/NLMonthNameParser.ts.

package nl

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

var monthNamePattern = core.MustRegExp(
	"("+core.MatchAnyPattern(MonthDictionary.Keys())+")"+
		"\\s*"+
		"(?:"+
		"[,-]?\\s*("+YearPattern+")?"+
		")?"+
		"(?=[^\\s\\w]|\\s+[^0-9]|\\s+$|$)",
	"i",
)

const (
	monthNameMonthNameGroup = 1
	monthNameYearGroup      = 2
)

// MonthNameParser is NLMonthNameParser, for a month name and a year: "januari, 2012",
// "januari 2012", "januari".
type MonthNameParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewMonthNameParser returns a new NLMonthNameParser.
func NewMonthNameParser() *MonthNameParser {
	p := &MonthNameParser{}
	p.Self = p
	return p
}

// InnerPattern returns the month name pattern.
func (p *MonthNameParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return monthNamePattern
}

// InnerExtract returns the components of the month (and year).
func (p *MonthNameParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	components := context.CreateParsingComponents(nil)
	components.Imply(core.ComponentDay, 1)

	monthName := match.Get(monthNameMonthNameGroup)
	month := monthLookup(js.ToLowerCase(monthName))
	components.Assign(core.ComponentMonth, month)

	if match.Truthy(monthNameYearGroup) {
		year := ParseYear(match.Get(monthNameYearGroup))
		components.Assign(core.ComponentYear, year)
	} else {
		year := core.FindYearClosestToRef(context.RefDate, 1, core.Num(month))
		components.Imply(core.ComponentYear, year)
	}

	return components
}
