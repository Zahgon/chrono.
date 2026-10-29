// Ported from src/locales/ru/parsers/RURelativeDateFormatParser.ts.

package ru

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
	"github.com/wanasit/chrono-go/internal/jsdate"
)

const (
	relativeModifierWordGroup = 1
	relativeRelativeWordGroup = 2
)

var (
	reWeek  = core.MustRegExp("week", "i")
	reMonth = core.MustRegExp("month", "i")
	reYear  = core.MustRegExp("year", "i")
)

// RelativeDateFormatParser is RURelativeDateFormatParser ("в прошлом месяце", "на следующей
// неделе", "в этом году", ...).
type RelativeDateFormatParser struct {
	AbstractParserWithLeftRightBoundaryChecking
}

// NewRelativeDateFormatParser returns a new RURelativeDateFormatParser.
func NewRelativeDateFormatParser() *RelativeDateFormatParser {
	p := &RelativeDateFormatParser{}
	p.Self = p
	return p
}

// InnerPatternString returns the relative date pattern.
func (p *RelativeDateFormatParser) InnerPatternString(context *core.ParsingContext) string {
	return "(в прошлом|на прошлой|на следующей|в следующем|на этой|в этом)\\s*(" + core.MatchAnyPattern(
		TimeUnitDictionary.Keys(),
	) + ")"
}

// InnerExtract returns the components of the relative date.
func (p *RelativeDateFormatParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	modifier := js.ToLowerCase(match.Get(relativeModifierWordGroup))
	unitWord := js.ToLowerCase(match.Get(relativeRelativeWordGroup))
	timeunit, timeunitDefined := TimeUnitDictionary.Get(unitWord)
	timeunitKey := timeunit
	if !timeunitDefined {
		timeunitKey = "undefined" // timeUnits[undefined]
	}

	if modifier == "на следующей" || modifier == "в следующем" {
		timeUnits := core.NewDuration()
		timeUnits.Set(timeunitKey, 1)
		return core.CreateRelativeFromReference(context.Reference, timeUnits)
	}

	if modifier == "в прошлом" || modifier == "на прошлой" {
		timeUnits := core.NewDuration()
		timeUnits.Set(timeunitKey, -1)
		return core.CreateRelativeFromReference(context.Reference, timeUnits)
	}

	components := context.CreateParsingComponents(nil)
	date := jsdate.FromTime(context.Reference.Instant)

	if !timeunitDefined {
		panic(&core.Error{Name: "TypeError", Message: "Cannot read properties of undefined (reading 'match')"})
	}

	// This week
	if reWeek.Exec(timeunit) != nil {
		date.SetDate(date.GetDate() - date.GetDay())
		components.Imply(core.ComponentDay, date.GetDate())
		components.Imply(core.ComponentMonth, date.GetMonth()+1)
		components.Imply(core.ComponentYear, date.GetFullYear())
	} else if reMonth.Exec(timeunit) != nil {
		// This month
		date.SetDate(1)
		components.Imply(core.ComponentDay, date.GetDate())
		components.Assign(core.ComponentYear, date.GetFullYear())
		components.Assign(core.ComponentMonth, date.GetMonth()+1)
	} else if reYear.Exec(timeunit) != nil {
		// This year
		date.SetDate(1)
		date.SetMonth(0)
		components.Imply(core.ComponentDay, date.GetDate())
		components.Imply(core.ComponentMonth, date.GetMonth()+1)
		components.Assign(core.ComponentYear, date.GetFullYear())
	}

	return components
}
