// Ported from src/locales/nl/parsers/NLMonthNameMiddleEndianParser.ts.

package nl

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

var monthNameMiddleEndianPattern = core.MustRegExp(
	"(?:on\\s*?)?"+
		"("+OrdinalNumberPattern+")"+
		"(?:\\s*"+
		"(?:tot|\\-|\\–|until|through|till|\\s)\\s*"+
		"("+OrdinalNumberPattern+")"+
		")?"+
		"(?:-|/|\\s*(?:of)?\\s*)"+
		"("+
		core.MatchAnyPattern(MonthDictionary.Keys())+
		")"+
		"(?:"+
		"(?:-|/|,?\\s*)"+
		"("+YearPattern+"(?![^\\s]\\d))"+
		")?"+
		"(?=\\W|$)",
	"i",
)

const (
	middleEndianMonthNameGroup = 3
	middleEndianDateGroup      = 1
	middleEndianDateToGroup    = 2
	middleEndianYearGroup      = 4
)

// MonthNameMiddleEndianParser is NLMonthNameMiddleEndianParser, for the BE/NL date format with
// the month's name in full writing: "1 januari 2019", "01 januari 2019", "10 januari 2019",
// "13 januari", "10 - 25 maart", "10 - 25 maart 2019", "1 aug 2019",
// "1 september 200 voor Christus", "1 september 2002 na Christus", "19 januari 87",
// "12de juli 2013", "1ste november 2013".
type MonthNameMiddleEndianParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewMonthNameMiddleEndianParser returns a new NLMonthNameMiddleEndianParser.
func NewMonthNameMiddleEndianParser() *MonthNameMiddleEndianParser {
	p := &MonthNameMiddleEndianParser{}
	p.Self = p
	return p
}

// InnerPattern returns the day-month-year pattern.
func (p *MonthNameMiddleEndianParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return monthNameMiddleEndianPattern
}

// InnerExtract returns the components of the date, or the result of the date range.
func (p *MonthNameMiddleEndianParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	month := monthLookup(js.ToLowerCase(match.Get(middleEndianMonthNameGroup)))
	day := ParseOrdinalNumberPattern(match.Get(middleEndianDateGroup))
	if day > 31 {
		// e.g. "[96 Aug]" => "9[6 Aug]", we need to shift away from the next number
		match.Index = match.Index + js.Len(match.Get(middleEndianDateGroup))
		return nil
	}

	components := context.CreateParsingComponents(core.Components{
		{Component: core.ComponentDay, Value: day},
		{Component: core.ComponentMonth, Value: month},
	})

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
