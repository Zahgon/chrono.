// Ported from src/locales/nl/parsers/NLCasualTimeParser.ts.

package nl

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
	"github.com/wanasit/chrono-go/internal/jsdate"
)

var casualTimePattern = core.MustRegExp("(deze)?\\s*(namiddag|avond|middernacht|ochtend|middag|'s middags|'s avonds|'s ochtends)(?=\\W|$)", "i")

const (
	casualTimeDayGroup    = 1
	casualTimeMomentGroup = 2
)

// CasualTimeParser is NLCasualTimeParser ("avond", "deze middag", "'s ochtends", ...).
type CasualTimeParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewCasualTimeParser returns a new NLCasualTimeParser.
func NewCasualTimeParser() *CasualTimeParser {
	p := &CasualTimeParser{}
	p.Self = p
	return p
}

// InnerPattern returns the casual time pattern.
func (p *CasualTimeParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return casualTimePattern
}

// InnerExtract returns the components of the casual time.
func (p *CasualTimeParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	targetDate := jsdate.FromTime(context.RefDate)
	component := context.CreateParsingComponents(nil)

	if match.Has(casualTimeDayGroup) && match.Get(casualTimeDayGroup) == "deze" {
		refDate := jsdate.FromTime(context.RefDate)
		component.Assign(core.ComponentDay, refDate.GetDate())
		component.Assign(core.ComponentMonth, refDate.GetMonth()+1)
		component.Assign(core.ComponentYear, refDate.GetFullYear())
	}

	switch js.ToLowerCase(match.Get(casualTimeMomentGroup)) {
	case "namiddag", "'s namiddags":
		component.Imply(core.ComponentMeridiem, core.PM)
		component.Imply(core.ComponentHour, 15)
		component.AddTag("casualReference/afternoon")

	case "avond", "'s avonds'":
		component.Imply(core.ComponentMeridiem, core.PM)
		component.Imply(core.ComponentHour, 20)
		component.AddTag("casualReference/evening")

	case "middernacht":
		nextDay := jsdate.FromTime(targetDate.Time())
		nextDay.SetDate(nextDay.GetDate() + 1)
		core.AssignSimilarDate(component, nextDay)
		core.ImplySimilarTime(component, nextDay)
		component.Imply(core.ComponentHour, 0)
		component.Imply(core.ComponentMinute, 0)
		component.Imply(core.ComponentSecond, 0)
		component.AddTag("casualReference/midnight")

	case "ochtend", "'s ochtends":
		component.Imply(core.ComponentMeridiem, core.AM)
		component.Imply(core.ComponentHour, 6)
		component.AddTag("casualReference/morning")

	case "middag", "'s middags":
		component.Imply(core.ComponentMeridiem, core.AM)
		component.Imply(core.ComponentHour, 12)
		component.AddTag("casualReference/noon")
	}

	return component
}
