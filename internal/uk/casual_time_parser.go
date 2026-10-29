// Ported from src/locales/uk/parsers/UKCasualTimeParser.ts.

package uk

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
	"github.com/wanasit/chrono-go/internal/jsdate"
)

var (
	reLastNight   = core.MustRegExp("минулої\\s*ночі", "")
	reLastEvening = core.MustRegExp("минулого\\s*вечора", "")
	reNextNight   = core.MustRegExp("наступної\\s*ночі", "")
	reThisNight   = core.MustRegExp("цієї\\s*ночі", "")
)

// CasualTimeParser is UKCasualTimeParser ("зараз", "вранці", "ввечері", "опівночі", ...).
type CasualTimeParser struct {
	AbstractParserWithLeftRightBoundaryChecking
}

// NewCasualTimeParser returns a new UKCasualTimeParser.
func NewCasualTimeParser() *CasualTimeParser {
	p := &CasualTimeParser{}
	p.Self = p
	return p
}

// InnerPatternString returns the casual time pattern.
func (p *CasualTimeParser) InnerPatternString(context *core.ParsingContext) string {
	return "(зараз|минулого\\s*вечора|минулої\\s*ночі|наступної\\s*ночі|сьогодні\\s*вночі|цієї\\s*ночі|цього ранку|вранці|ранку|зранку|опівдні|ввечері|вечора|опівночі|вночі)"
}

// InnerExtract returns the components of the casual time.
func (p *CasualTimeParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	targetDate := jsdate.FromTime(context.RefDate)
	lowerText := js.ToLowerCase(match.Get(0))
	component := context.CreateParsingComponents(nil)

	if lowerText == "зараз" {
		return core.Now(context.Reference)
	}
	if lowerText == "ввечері" || lowerText == "вечора" {
		return core.Evening(context.Reference, 20)
	}
	if js.EndsWith(lowerText, "вранці") || js.EndsWith(lowerText, "ранку") || js.EndsWith(lowerText, "зранку") {
		return core.Morning(context.Reference, 6)
	}
	if js.EndsWith(lowerText, "опівдні") {
		return core.Noon(context.Reference)
	}
	if reLastNight.Exec(lowerText) != nil {
		return core.LastNight(context.Reference, 0)
	}
	if reLastEvening.Exec(lowerText) != nil {
		return core.YesterdayEvening(context.Reference, 20)
	}
	if reNextNight.Exec(lowerText) != nil {
		daysToAdd := 2.0
		if targetDate.GetHours() < 22 {
			daysToAdd = 1
		}
		nextDay := jsdate.FromTime(targetDate.Time())
		nextDay.SetDate(nextDay.GetDate() + daysToAdd)
		core.AssignSimilarDate(component, nextDay)
		component.Imply(core.ComponentHour, 1)
	}
	if reThisNight.Exec(lowerText) != nil {
		return core.Midnight(context.Reference)
	}
	if js.EndsWith(lowerText, "опівночі") || js.EndsWith(lowerText, "вночі") {
		return core.Midnight(context.Reference)
	}
	return component
}
