// Ported from src/locales/uk/parsers/UKMonthNameLittleEndianParser.ts.

package uk

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

const (
	littleEndianDateGroup      = 1
	littleEndianDateToGroup    = 2
	littleEndianMonthNameGroup = 3
	littleEndianYearGroup      = 4
)

// MonthNameLittleEndianParser is UKMonthNameLittleEndianParser: a day (or day range) before a
// month name, with an optional year ("10 серпня 2012", "з 10 по 12 серпня").
type MonthNameLittleEndianParser struct {
	AbstractParserWithLeftRightBoundaryChecking
}

// NewMonthNameLittleEndianParser returns a new UKMonthNameLittleEndianParser.
func NewMonthNameLittleEndianParser() *MonthNameLittleEndianParser {
	p := &MonthNameLittleEndianParser{}
	p.Self = p
	return p
}

// InnerPatternString returns the little-endian month name pattern.
func (p *MonthNameLittleEndianParser) InnerPatternString(context *core.ParsingContext) string {
	return "(?:з|із)?\\s*(" + OrdinalNumberPattern + ")" +
		"(?:" +
		"\\s{0,3}(?:по|-|–|до)?\\s{0,3}" +
		"(" + OrdinalNumberPattern + ")" +
		")?" +
		"(?:-|\\/|\\s{0,3}(?:of)?\\s{0,3})" +
		"(" + core.MatchAnyPattern(MonthDictionary.Keys()) + ")" +
		"(?:" +
		"(?:-|\\/|,?\\s{0,3})" +
		"(" + YearPattern + "(?![^\\s]\\d))" +
		")?"
}

// InnerExtract returns the result for the day (range), month and year.
func (p *MonthNameLittleEndianParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	result := context.CreateParsingResult(match.Index, match.Get(0), nil, nil)

	month := monthLookup(js.ToLowerCase(match.Get(littleEndianMonthNameGroup)))
	day := ParseOrdinalNumberPattern(match.Get(littleEndianDateGroup))
	if day > 31 {
		// e.g. "[96 Aug]" => "9[6 Aug]", we need to shift away from the next number
		match.Index = match.Index + js.Len(match.Get(littleEndianDateGroup))
		return nil
	}

	result.Start.Assign(core.ComponentMonth, month)
	result.Start.Assign(core.ComponentDay, day)

	if match.Truthy(littleEndianYearGroup) {
		yearNumber := ParseYearPattern(match.Get(littleEndianYearGroup))
		result.Start.Assign(core.ComponentYear, yearNumber)
	} else {
		year := core.FindYearClosestToRef(context.Reference.Instant, day, core.Num(month))
		result.Start.Imply(core.ComponentYear, year)
	}

	if match.Truthy(littleEndianDateToGroup) {
		endDate := ParseOrdinalNumberPattern(match.Get(littleEndianDateToGroup))

		result.End = result.Start.Clone()
		result.End.Assign(core.ComponentDay, endDate)
	}

	return result
}

// monthLookup is MONTH_DICTIONARY[name] (undefined when missing).
func monthLookup(name string) float64 {
	if v, ok := MonthDictionary.Get(name); ok {
		return v
	}
	return core.Undefined
}
