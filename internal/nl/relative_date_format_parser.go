// Ported from src/locales/nl/parsers/NLRelativeDateFormatParser.ts.

package nl

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
	"github.com/wanasit/chrono-go/internal/jsdate"
)

var relativeDateFormatPattern = core.MustRegExp(
	"(dit|deze|(?:aan)?komend|volgend|afgelopen|vorig)e?\\s*("+core.MatchAnyPattern(TimeUnitDictionary.Keys())+")(?=\\s*)"+
		"(?=\\W|$)",
	"i",
)

const (
	relativeModifierWordGroup = 1
	relativeRelativeWordGroup = 2
)

var (
	reWeek  = core.MustRegExp("week", "i")
	reMaand = core.MustRegExp("maand", "i")
	reJaar  = core.MustRegExp("jaar", "i")
)

// RelativeDateFormatParser is NLRelativeDateFormatParser ("deze week", "volgende maand",
// "afgelopen jaar").
type RelativeDateFormatParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewRelativeDateFormatParser returns a new NLRelativeDateFormatParser.
func NewRelativeDateFormatParser() *RelativeDateFormatParser {
	p := &RelativeDateFormatParser{}
	p.Self = p
	return p
}

// InnerPattern returns the relative date pattern.
func (p *RelativeDateFormatParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return relativeDateFormatPattern
}

// InnerExtract returns the components of the relative date.
func (p *RelativeDateFormatParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	modifier := js.ToLowerCase(match.Get(relativeModifierWordGroup))
	unitWord := js.ToLowerCase(match.Get(relativeRelativeWordGroup))
	timeunit := timeUnitLookup(unitWord)

	if modifier == "volgend" || modifier == "komend" || modifier == "aankomend" {
		timeUnits := core.NewDuration()
		timeUnits.Set(timeunit, 1)
		return core.CreateRelativeFromReference(context.Reference, timeUnits)
	}

	if modifier == "afgelopen" || modifier == "vorig" {
		timeUnits := core.NewDuration()
		timeUnits.Set(timeunit, -1)
		return core.CreateRelativeFromReference(context.Reference, timeUnits)
	}

	components := context.CreateParsingComponents(nil)
	date := jsdate.FromTime(context.Reference.Instant)

	if reWeek.Exec(unitWord) != nil {
		// This week
		date.SetDate(date.GetDate() - date.GetDay())
		components.Imply(core.ComponentDay, date.GetDate())
		components.Imply(core.ComponentMonth, date.GetMonth()+1)
		components.Imply(core.ComponentYear, date.GetFullYear())
	} else if reMaand.Exec(unitWord) != nil {
		// This month
		date.SetDate(1)
		components.Imply(core.ComponentDay, date.GetDate())
		components.Assign(core.ComponentYear, date.GetFullYear())
		components.Assign(core.ComponentMonth, date.GetMonth()+1)
	} else if reJaar.Exec(unitWord) != nil {
		// This year
		date.SetDate(1)
		date.SetMonth(0)
		components.Imply(core.ComponentDay, date.GetDate())
		components.Imply(core.ComponentMonth, date.GetMonth()+1)
		components.Assign(core.ComponentYear, date.GetFullYear())
	}

	return components
}
