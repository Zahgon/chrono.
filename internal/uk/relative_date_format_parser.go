// Ported from src/locales/uk/parsers/UKRelativeDateFormatParser.ts.

package uk

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

// RelativeDateFormatParser is UKRelativeDateFormatParser ("в минулому місяці", "наступного
// тижня", "у цьому році", ...).
type RelativeDateFormatParser struct {
	AbstractParserWithLeftRightBoundaryChecking
}

// NewRelativeDateFormatParser returns a new UKRelativeDateFormatParser.
func NewRelativeDateFormatParser() *RelativeDateFormatParser {
	p := &RelativeDateFormatParser{}
	p.Self = p
	return p
}

// InnerPatternString returns the relative date pattern.
func (p *RelativeDateFormatParser) InnerPatternString(context *core.ParsingContext) string {
	return "(в минулому|у минулому|на минулому|минулого|на наступному|в наступному|у наступному|наступного|на цьому|в цьому|у цьому|цього)\\s*" +
		"(" + core.MatchAnyPattern(TimeUnitDictionary.Keys()) + ")(?=\\s*)"
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

	if modifier == "на наступному" ||
		modifier == "в наступному" ||
		modifier == "у наступному" ||
		modifier == "наступного" {
		timeUnits := core.NewDuration()
		timeUnits.Set(timeunitKey, 1)
		return core.CreateRelativeFromReference(context.Reference, timeUnits)
	}

	if modifier == "на минулому" ||
		modifier == "в минулому" ||
		modifier == "у минулому" ||
		modifier == "минулого" {
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
