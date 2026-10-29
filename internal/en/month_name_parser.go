// Ported from src/locales/en/parsers/ENMonthNameParser.ts.

package en

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

var monthNamePattern = core.MustRegExp(
	"((?:in)\\s*)?"+
		"("+core.MatchAnyPattern(MonthDictionary.Keys())+")"+
		"\\s*"+
		"(?:"+
		"(?:,|-|of)?\\s*("+YearPattern+")?"+
		")?"+
		"(?=[^\\s\\w]|\\s+[^0-9]|\\s+$|$)",
	"i",
)

const (
	monthNamePrefixGroup    = 1
	monthNameMonthNameGroup = 2
	monthNameYearGroup      = 3
)

// MonthNameParser is ENMonthNameParser: a month name with an optional year ("January 2012",
// "(in) Jan").
type MonthNameParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewMonthNameParser returns a new ENMonthNameParser.
func NewMonthNameParser() *MonthNameParser {
	p := &MonthNameParser{}
	p.Self = p
	return p
}

// InnerPattern returns the month name pattern.
func (p *MonthNameParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return monthNamePattern
}

// InnerExtract returns the result for the month (and year).
func (p *MonthNameParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	monthName := js.ToLowerCase(match.Get(monthNameMonthNameGroup))

	// skip some unlikely words "jan", "mar", ..
	if full, ok := FullMonthNameDictionary.Get(monthName); js.Len(match.Get(0)) <= 3 && !(ok && core.Truthy(full)) {
		return nil
	}

	result := context.CreateParsingResultToIndex(
		match.Index+js.Len(match.Get(monthNamePrefixGroup)),
		match.Index+js.Len(match.Get(0)),
		nil, nil,
	)
	result.Start.Imply(core.ComponentDay, 1)
	result.Start.AddTag("parser/ENMonthNameParser")

	month := monthLookup(monthName)
	result.Start.Assign(core.ComponentMonth, month)

	if match.Truthy(monthNameYearGroup) {
		year := ParseYear(match.Get(monthNameYearGroup))
		result.Start.Assign(core.ComponentYear, year)
	} else {
		year := core.FindYearClosestToRef(context.RefDate, 1, month)
		result.Start.Imply(core.ComponentYear, year)
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
