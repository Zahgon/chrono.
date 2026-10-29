// Ported from src/locales/en/parsers/ENRelativeDateFormatParser.ts.

package en

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
	"github.com/wanasit/chrono-go/internal/jsdate"
)

var relativeDateFormatPattern = core.MustRegExp(
	"(this|last|past|next|after\\s*this)\\s*("+core.MatchAnyPattern(TimeUnitDictionary.Keys())+")(?=\\s*)"+"(?=\\W|$)",
	"i",
)

const (
	relativeModifierWordGroup = 1
	relativeRelativeWordGroup = 2
)

var (
	reWeekI  = core.MustRegExp("week", "i")
	reMonthI = core.MustRegExp("month", "i")
	reYearI  = core.MustRegExp("year", "i")
)

// RelativeDateFormatParser is ENRelativeDateFormatParser ("next week", "last month", "this
// year").
type RelativeDateFormatParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewRelativeDateFormatParser returns a new ENRelativeDateFormatParser.
func NewRelativeDateFormatParser() *RelativeDateFormatParser {
	p := &RelativeDateFormatParser{}
	p.Self = p
	return p
}

// InnerPattern returns the relative date pattern.
func (p *RelativeDateFormatParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return relativeDateFormatPattern
}

func timeunitLookup(word string) string {
	if v, ok := TimeUnitDictionary.Get(word); ok {
		return v
	}
	return "undefined"
}

// InnerExtract returns the components of the relative date.
func (p *RelativeDateFormatParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	modifier := js.ToLowerCase(match.Get(relativeModifierWordGroup))
	unitWord := js.ToLowerCase(match.Get(relativeRelativeWordGroup))
	timeunit := timeunitLookup(unitWord)

	if modifier == "next" || js.StartsWith(modifier, "after", 0) {
		timeUnits := core.NewDuration()
		timeUnits.Set(timeunit, 1)
		return core.CreateRelativeFromReference(context.Reference, timeUnits)
	}

	if modifier == "last" || modifier == "past" {
		timeUnits := core.NewDuration()
		timeUnits.Set(timeunit, -1)
		return core.CreateRelativeFromReference(context.Reference, timeUnits)
	}

	components := context.CreateParsingComponents(nil)
	date := jsdate.FromTime(context.Reference.Instant)

	if reWeekI.Exec(unitWord) != nil {
		// This week
		date.SetDate(date.GetDate() - date.GetDay())
		components.Imply(core.ComponentDay, date.GetDate())
		components.Imply(core.ComponentMonth, date.GetMonth()+1)
		components.Imply(core.ComponentYear, date.GetFullYear())
	} else if reMonthI.Exec(unitWord) != nil {
		// This month
		date.SetDate(1)
		components.Imply(core.ComponentDay, date.GetDate())
		components.Assign(core.ComponentYear, date.GetFullYear())
		components.Assign(core.ComponentMonth, date.GetMonth()+1)
	} else if reYearI.Exec(unitWord) != nil {
		// This year
		date.SetDate(1)
		date.SetMonth(0)
		components.Imply(core.ComponentDay, date.GetDate())
		components.Imply(core.ComponentMonth, date.GetMonth()+1)
		components.Assign(core.ComponentYear, date.GetFullYear())
	}
	return components
}
