// Ported from src/locales/uk/parsers/UKMonthNameParser.ts.

package uk

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

const (
	monthNameMonthNameGroup = 2
	monthNameYearGroup      = 3
)

// MonthNameParser is UkMonthNameParser: a month name with an optional year ("Cічень, 2012",
// "Січень 2012", "Січень").
type MonthNameParser struct {
	AbstractParserWithLeftBoundaryChecking
}

// NewMonthNameParser returns a new UkMonthNameParser.
func NewMonthNameParser() *MonthNameParser {
	p := &MonthNameParser{}
	p.Self = p
	return p
}

// InnerPatternString returns the month name pattern.
func (p *MonthNameParser) InnerPatternString(context *core.ParsingContext) string {
	return "((?:в|у)\\s*)?" +
		"(" + core.MatchAnyPattern(MonthDictionary.Keys()) + ")" +
		"\\s*" +
		"(?:" +
		"[,-]?\\s*(" + YearPattern + ")?" +
		")?" +
		"(?=[^\\s\\w]|\\s+[^0-9]|\\s+$|$)"
}

// InnerExtract returns the result for the month (and year).
func (p *MonthNameParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	monthName := js.ToLowerCase(match.Get(monthNameMonthNameGroup))

	// skip some unlikely words "січ", "лют", ..
	if full, ok := FullMonthNameDictionary.Get(monthName); js.Len(match.Get(0)) <= 3 && !(ok && core.Truthy(full)) {
		return nil
	}

	result := context.CreateParsingResultToIndex(match.Index, match.Index+js.Len(match.Get(0)), nil, nil)
	result.Start.Imply(core.ComponentDay, 1)

	month := monthLookup(monthName)
	result.Start.Assign(core.ComponentMonth, month)

	if match.Truthy(monthNameYearGroup) {
		year := ParseYearPattern(match.Get(monthNameYearGroup))
		result.Start.Assign(core.ComponentYear, year)
	} else {
		year := core.FindYearClosestToRef(context.Reference.Instant, 1, core.Num(month))
		result.Start.Imply(core.ComponentYear, year)
	}

	return result
}
