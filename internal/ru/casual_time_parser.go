// Ported from src/locales/ru/parsers/RUCasualTimeParser.ts.

package ru

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
	"github.com/wanasit/chrono-go/internal/jsdate"
)

var (
	reNoon          = core.MustRegExp("в\\s*полдень", "")
	reLastNight     = core.MustRegExp("прошлой\\s*ночью", "")
	reLastEvening   = core.MustRegExp("прошлым\\s*вечером", "")
	reNextNight     = core.MustRegExp("следующей\\s*ночью", "")
	reMidnightMatch = core.MustRegExp("в\\s*полночь", "")
)

// CasualTimeParser is RUCasualTimeParser ("сейчас", "утром", "вечером", "в полночь", ...).
type CasualTimeParser struct {
	AbstractParserWithLeftRightBoundaryChecking
}

// NewCasualTimeParser returns a new RUCasualTimeParser.
func NewCasualTimeParser() *CasualTimeParser {
	p := &CasualTimeParser{}
	p.Self = p
	return p
}

// InnerPatternString returns the casual time pattern.
func (p *CasualTimeParser) InnerPatternString(context *core.ParsingContext) string {
	return "(сейчас|прошлым\\s*вечером|прошлой\\s*ночью|следующей\\s*ночью|сегодня\\s*ночью|этой\\s*ночью|ночью|этим утром|утром|утра|в\\s*полдень|вечером|вечера|в\\s*полночь)"
}

// InnerExtract returns the components of the casual time.
func (p *CasualTimeParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	targetDate := jsdate.FromTime(context.RefDate)
	lowerText := js.ToLowerCase(match.Get(0))
	component := context.CreateParsingComponents(nil)

	if lowerText == "сейчас" {
		return core.Now(context.Reference)
	}
	if lowerText == "вечером" || lowerText == "вечера" {
		return core.Evening(context.Reference, 20)
	}
	if js.EndsWith(lowerText, "утром") || js.EndsWith(lowerText, "утра") {
		return core.Morning(context.Reference, 6)
	}
	if reNoon.Exec(lowerText) != nil {
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
		component.Imply(core.ComponentHour, 0)
	}
	if reMidnightMatch.Exec(lowerText) != nil || js.EndsWith(lowerText, "ночью") {
		return core.Midnight(context.Reference)
	}
	return component
}
