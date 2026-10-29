// Ported from src/locales/fi/parsers/FIMonthNameLittleEndianParser.ts.

package fi

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

var monthNameLittleEndianPattern = core.MustRegExp(
	"([0-9]{1,2})\\.?"+
		"(?:\\s*(?:\\-|\\–|\\s)\\s*([0-9]{1,2})\\.?)?\\s*"+
		"("+core.MatchAnyPattern(MonthDictionary.Keys())+")"+
		"(?:(?:-|/|,?\\s*)([0-9]{4}(?![^\\s]\\d)))?"+
		"(?=\\W|$)",
	"i",
)

const (
	littleEndianDateGroup      = 1
	littleEndianDateToGroup    = 2
	littleEndianMonthNameGroup = 3
	littleEndianYearGroup      = 4
)

// MonthNameLittleEndianParser is FIMonthNameLittleEndianParser ("15. tammikuuta 2024",
// "1.-3. maaliskuuta").
type MonthNameLittleEndianParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewMonthNameLittleEndianParser returns a new FIMonthNameLittleEndianParser.
func NewMonthNameLittleEndianParser() *MonthNameLittleEndianParser {
	p := &MonthNameLittleEndianParser{}
	p.Self = p
	return p
}

// InnerPattern returns the little-endian month name pattern.
func (p *MonthNameLittleEndianParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return monthNameLittleEndianPattern
}

func monthLookup(name string) float64 {
	if v, ok := MonthDictionary.Get(name); ok {
		return v
	}
	return core.Undefined
}

// InnerExtract returns the result of the date (range).
func (p *MonthNameLittleEndianParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	result := context.CreateParsingResult(match.Index, match.Get(0), nil, nil)

	month := monthLookup(js.ToLowerCase(match.Get(littleEndianMonthNameGroup)))
	day := js.ParseInt(match.Get(littleEndianDateGroup), 0)
	if day > 31 {
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
		endDate := js.ParseInt(match.Get(littleEndianDateToGroup), 0)
		result.End = result.Start.Clone()
		result.End.Assign(core.ComponentDay, endDate)
	}

	return result
}
