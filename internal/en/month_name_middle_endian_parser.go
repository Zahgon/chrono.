// Ported from src/locales/en/parsers/ENMonthNameMiddleEndianParser.ts.

package en

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

var monthNameMiddleEndianPattern = core.MustRegExp(
	"("+core.MatchAnyPattern(MonthDictionary.Keys())+")"+
		"(?:-|/|\\s*,?\\s*)"+
		"("+OrdinalNumberPattern+")(?!\\s*(?:am|pm))\\s*"+
		"(?:"+
		"(?:to|\\-)\\s*"+
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

var reTwoDigits = core.MustRegExp("^\\d{2}$", "")

// MonthNameMiddleEndianParser is ENMonthNameMiddleEndianParser ("January 13, 2012",
// "January 13 - 15, 2012").
type MonthNameMiddleEndianParser struct {
	core.AbstractParserWithWordBoundaryChecking
	ShouldSkipYearLikeDate bool
}

// NewMonthNameMiddleEndianParser returns a new ENMonthNameMiddleEndianParser.
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

	// Skip the case where the day looks like a year (ex: January 21)
	if p.ShouldSkipYearLikeDate {
		if !match.Truthy(middleEndianDateToGroup) && !match.Truthy(middleEndianYearGroup) && reTwoDigits.Exec(match.Get(middleEndianDateGroup)) != nil {
			return nil
		}
	}
	components := context.CreateParsingComponents(core.Components{
		{Component: core.ComponentDay, Value: day},
		{Component: core.ComponentMonth, Value: month},
	}).AddTag("parser/ENMonthNameMiddleEndianParser")

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

	// Text can be 'range' value. Such as 'January 12 - 13, 2012'
	endDate := ParseOrdinalNumberPattern(match.Get(middleEndianDateToGroup))
	result := context.CreateParsingResult(match.Index, match.Get(0), nil, nil)
	result.Start = components
	result.End = components.Clone()
	result.End.Assign(core.ComponentDay, endDate)
	return result
}
