// Ported from src/locales/ru/parsers/RUMonthNameParser.ts.

package ru

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

const (
	monthNameMonthNameGroup = 2
	monthNameYearGroup      = 3
)

// MonthNameParser is RUMonthNameParser: a month name with an optional year ("Январь, 2012",
// "Январь 2012", "Январь").
type MonthNameParser struct {
	AbstractParserWithLeftBoundaryChecking
}

// NewMonthNameParser returns a new RUMonthNameParser.
func NewMonthNameParser() *MonthNameParser {
	p := &MonthNameParser{}
	p.Self = p
	return p
}

// InnerPatternString returns the month name pattern.
func (p *MonthNameParser) InnerPatternString(context *core.ParsingContext) string {
	return "((?:в)\\s*)?" +
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

	// skip some unlikely words "янв", "фер", ..
	if full, ok := FullMonthNameDictionary.Get(monthName); js.Len(match.Get(0)) <= 3 && !(ok && core.Truthy(full)) {
		return nil
	}

	result := context.CreateParsingResultToIndex(match.Index, match.Index+js.Len(match.Get(0)), nil, nil)
	result.Start.Imply(core.ComponentDay, 1)

	month := monthLookup(monthName)
	result.Start.Assign(core.ComponentMonth, month)

	if match.Truthy(monthNameYearGroup) {
		year := ParseYear(match.Get(monthNameYearGroup))
		result.Start.Assign(core.ComponentYear, year)
	} else {
		year := core.FindYearClosestToRef(context.RefDate, 1, core.Num(month))
		result.Start.Imply(core.ComponentYear, year)
	}

	return result
}
