// Ported from src/locales/sv/parsers/SVCasualDateParser.ts.

package sv

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
	"github.com/wanasit/chrono-go/internal/jsdate"
)

var casualDatePattern = core.MustRegExp(
	"(nu|idag|imorgon|övermorgon|igår|förrgår|i\\s*förrgår)"+
		"(?:\\s*(?:på\\s*)?(morgonen?|förmiddagen?|middagen?|eftermiddagen?|kvällen?|natten?|midnatt))?"+
		"(?=\\W|$)",
	"i",
)

const (
	casualDateDateGroup = 1
	casualDateTimeGroup = 2
)

// CasualDateParser is SVCasualDateParser ("nu", "idag", "imorgon", "igår på kvällen", ...).
type CasualDateParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewCasualDateParser returns a new SVCasualDateParser.
func NewCasualDateParser() *CasualDateParser {
	p := &CasualDateParser{}
	p.Self = p
	return p
}

// InnerPattern returns the casual date pattern.
func (p *CasualDateParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return casualDatePattern
}

// InnerExtract returns the components of the casual date (and time).
func (p *CasualDateParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	targetDate := jsdate.FromTime(context.RefDate)
	dateKeyword := js.ToLowerCase(match.Get(casualDateDateGroup))
	timeKeyword := js.ToLowerCase(match.Get(casualDateTimeGroup))

	component := context.CreateParsingComponents(nil)
	switch dateKeyword {
	case "nu":
		component = core.Now(context.Reference)

	case "idag":
		component = core.Today(context.Reference)

	case "imorgon", "imorn":
		nextDay := jsdate.FromTime(targetDate.Time())
		nextDay.SetDate(nextDay.GetDate() + 1)
		core.AssignSimilarDate(component, nextDay)
		core.ImplySimilarTime(component, nextDay)

	case "igår":
		previousDay := jsdate.FromTime(targetDate.Time())
		previousDay.SetDate(previousDay.GetDate() - 1)
		core.AssignSimilarDate(component, previousDay)
		core.ImplySimilarTime(component, previousDay)

	case "förrgår", "i förrgår":
		twoDaysAgo := jsdate.FromTime(targetDate.Time())
		twoDaysAgo.SetDate(twoDaysAgo.GetDate() - 2)
		core.AssignSimilarDate(component, twoDaysAgo)
		core.ImplySimilarTime(component, twoDaysAgo)
	}

	switch timeKeyword {
	case "morgon", "morgonen":
		component.Imply(core.ComponentHour, 6)
		component.Imply(core.ComponentMinute, 0)
		component.Imply(core.ComponentSecond, 0)
		component.Imply(core.ComponentMillisecond, 0)

	case "förmiddag", "förmiddagen":
		component.Imply(core.ComponentHour, 9)
		component.Imply(core.ComponentMinute, 0)
		component.Imply(core.ComponentSecond, 0)
		component.Imply(core.ComponentMillisecond, 0)

	case "middag", "middagen":
		component.Imply(core.ComponentHour, 12)
		component.Imply(core.ComponentMinute, 0)
		component.Imply(core.ComponentSecond, 0)
		component.Imply(core.ComponentMillisecond, 0)

	case "eftermiddag", "eftermiddagen":
		component.Imply(core.ComponentHour, 15)
		component.Imply(core.ComponentMinute, 0)
		component.Imply(core.ComponentSecond, 0)
		component.Imply(core.ComponentMillisecond, 0)

	case "kväll", "kvällen":
		component.Imply(core.ComponentHour, 20)
		component.Imply(core.ComponentMinute, 0)
		component.Imply(core.ComponentSecond, 0)
		component.Imply(core.ComponentMillisecond, 0)

	case "natt", "natten", "midnatt":
		if timeKeyword == "midnatt" {
			component.Imply(core.ComponentHour, 0)
		} else {
			component.Imply(core.ComponentHour, 2)
		}
		component.Imply(core.ComponentMinute, 0)
		component.Imply(core.ComponentSecond, 0)
		component.Imply(core.ComponentMillisecond, 0)
	}

	return component
}
