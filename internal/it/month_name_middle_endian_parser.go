// Ported from src/locales/it/parsers/ITMonthNameMiddleEndianParser.ts.

package it

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

var monthNameMiddleEndianPattern = core.MustRegExp(
	"("+core.MatchAnyPattern(MonthDictionary.Keys())+")"+
		"(?:-|/|\\s*,?\\s*)"+
		"("+OrdinalNumberPattern+")(?!\\s*(?:am|pm))\\s*"+
		"(?:"+
		"(?:al|\\-)\\s*"+
		"("+OrdinalNumberPattern+")\\s*"+
		")?"+
		"(?:"+
		"(?:-|/|\\s*,\\s*|\\s+)"+
		"("+YearPattern+")"+
		")?"+
		"(?=\\W|$)(?!\\:\\d)",
	"i",
)

const (
	middleEndianMonthNameGroup = 1
	middleEndianDateGroup      = 2
	middleEndianDateToGroup    = 3
	middleEndianYearGroup      = 4
)

var reYearLikeDay = core.MustRegExp("^2[0-5]$", "")

// MonthNameMiddleEndianParser is ITMonthNameMiddleEndianParser: dates that begin with the month
// name ("Gennaio 13", "Gennaio 13, 2012", "Gennaio 13 - 15, 2012"). With shouldSkipYearLikeDate,
// "Gennaio 21" (a day that looks like a year) is skipped.
type MonthNameMiddleEndianParser struct {
	core.AbstractParserWithWordBoundaryChecking
	ShouldSkipYearLikeDate bool
}

// NewMonthNameMiddleEndianParser returns a new ITMonthNameMiddleEndianParser.
func NewMonthNameMiddleEndianParser(shouldSkipYearLikeDate bool) *MonthNameMiddleEndianParser {
	p := &MonthNameMiddleEndianParser{ShouldSkipYearLikeDate: shouldSkipYearLikeDate}
	p.Self = p
	return p
}

// InnerPattern returns the middle-endian month name pattern.
func (p *MonthNameMiddleEndianParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return monthNameMiddleEndianPattern
}

// InnerExtract returns the components (or range result) of the date.
func (p *MonthNameMiddleEndianParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	month := monthLookup(js.ToLowerCase(match.Get(middleEndianMonthNameGroup)))
	day := ParseOrdinalNumberPattern(match.Get(middleEndianDateGroup))
	if day > 31 {
		return nil
	}

	// Skip the case where the day looks like a year (ex: Gennaio 21)
	if p.ShouldSkipYearLikeDate {
		if !match.Truthy(middleEndianDateToGroup) && !match.Truthy(middleEndianYearGroup) && reYearLikeDay.Exec(match.Get(middleEndianDateGroup)) != nil {
			return nil
		}
	}
	components := context.CreateParsingComponents(core.Components{
		{Component: core.ComponentDay, Value: day},
		{Component: core.ComponentMonth, Value: month},
	}).AddTag("parser/ITMonthNameMiddleEndianParser")

	if match.Truthy(middleEndianYearGroup) {
		year := ParseYear(match.Get(middleEndianYearGroup))
		components.Assign(core.ComponentYear, year)
	} else {
		year := core.FindYearClosestToRef(context.RefDate, day, core.Num(month))
		components.Imply(core.ComponentYear, year)
	}

	if !match.Truthy(middleEndianDateToGroup) {
		return components
	}

	// Text can be 'range' value. Such as 'Gennaio 12 - 13, 2012'
	endDate := ParseOrdinalNumberPattern(match.Get(middleEndianDateToGroup))
	result := context.CreateParsingResult(match.Index, match.Get(0), nil, nil)
	result.Start = components
	result.End = components.Clone()
	result.End.Assign(core.ComponentDay, endDate)

	return result
}
