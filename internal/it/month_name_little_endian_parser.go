// Ported from src/locales/it/parsers/ITMonthNameLittleEndianParser.ts.

package it

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

var monthNameLittleEndianPattern = core.MustRegExp(
	"(?:il\\s{0,3})?"+
		"("+OrdinalNumberPattern+")"+
		"(?:"+
		"\\s{0,3}(?:al|\\-|\\–|fino\\s*al?)?\\s{0,3}"+
		"("+OrdinalNumberPattern+")"+
		")?"+
		"(?:-|/|\\s{0,3}(?:di)?\\s{0,3})"+
		"("+core.MatchAnyPattern(MonthDictionary.Keys())+")"+
		"(?:"+
		"(?:-|/|,?\\s{0,3})"+
		"("+YearPattern+"(?!\\w))"+
		")?"+
		"(?=\\W|$)",
	"i",
)

const (
	littleEndianDateGroup      = 1
	littleEndianDateToGroup    = 2
	littleEndianMonthNameGroup = 3
	littleEndianYearGroup      = 4
)

// MonthNameLittleEndianParser is ITMonthNameLittleEndianParser ("12 gennaio 2012", "il 1° - 3
// marzo").
type MonthNameLittleEndianParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewMonthNameLittleEndianParser returns a new ITMonthNameLittleEndianParser.
func NewMonthNameLittleEndianParser() *MonthNameLittleEndianParser {
	p := &MonthNameLittleEndianParser{}
	p.Self = p
	return p
}

// InnerPattern returns the little-endian month name pattern.
func (p *MonthNameLittleEndianParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return monthNameLittleEndianPattern
}

// InnerExtract returns the result of the date (range).
func (p *MonthNameLittleEndianParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	result := context.CreateParsingResult(match.Index, match.Get(0), nil, nil)

	month := monthLookup(js.ToLowerCase(match.Get(littleEndianMonthNameGroup)))
	day := ParseOrdinalNumberPattern(match.Get(littleEndianDateGroup))
	if day > 31 {
		// e.g. "[96 Ago]" => "9[6 Ago]", we need to shift away from the next number
		match.Index = match.Index + js.Len(match.Get(littleEndianDateGroup))
		return nil
	}

	result.Start.Assign(core.ComponentMonth, month)
	result.Start.Assign(core.ComponentDay, day)

	if match.Truthy(littleEndianYearGroup) {
		yearNumber := ParseYear(match.Get(littleEndianYearGroup))
		result.Start.Assign(core.ComponentYear, yearNumber)
	} else {
		year := core.FindYearClosestToRef(context.RefDate, day, core.Num(month))
		result.Start.Imply(core.ComponentYear, year)
	}

	if match.Truthy(littleEndianDateToGroup) {
		endDate := ParseOrdinalNumberPattern(match.Get(littleEndianDateToGroup))

		result.End = result.Start.Clone()
		result.End.Assign(core.ComponentDay, endDate)
	}

	return result
}
