// Ported from src/locales/fi/parsers/FICasualDateParser.ts.

package fi

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

var casualDatePattern = core.MustRegExp(
	"(nyt|tänään|huomenna|ylihuomenna|eilen|toissapäivänä|viime\\s*yönä)"+
		"(?:\\s*(aamulla|aamuna|aamupäivällä|päivällä|iltapäivällä|illalla|yöllä|keskiyöllä))?"+
		"(?=\\W|$)",
	"i",
)

const (
	casualDateDateGroup = 1
	casualDateTimeGroup = 2
)

var reViimeYona = core.MustRegExp("viime\\s*yönä", "")

// CasualDateParser is FICasualDateParser ("tänään", "huomenna illalla", "viime yönä", ...).
type CasualDateParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewCasualDateParser returns a new FICasualDateParser.
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
	targetDate := context.Reference.DateWithAdjustedTimezone()
	dateKeyword := js.ToLowerCase(match.Get(casualDateDateGroup))
	timeKeyword := js.ToLowerCase(match.Get(casualDateTimeGroup))

	component := context.CreateParsingComponents(nil)
	switch dateKeyword {
	case "nyt":
		component = core.Now(context.Reference)

	case "tänään":
		component = core.Today(context.Reference)

	case "huomenna":
		targetDate = core.AddDuration(targetDate, core.NewDuration(core.DurationValue{Unit: "day", Value: 1}))
		core.AssignSimilarDate(component, targetDate)
		core.ImplySimilarTime(component, targetDate)

	case "ylihuomenna":
		targetDate = core.AddDuration(targetDate, core.NewDuration(core.DurationValue{Unit: "day", Value: 2}))
		core.AssignSimilarDate(component, targetDate)
		core.ImplySimilarTime(component, targetDate)

	case "eilen":
		targetDate = core.AddDuration(targetDate, core.NewDuration(core.DurationValue{Unit: "day", Value: -1}))
		core.AssignSimilarDate(component, targetDate)
		core.ImplySimilarTime(component, targetDate)

	case "toissapäivänä":
		targetDate = core.AddDuration(targetDate, core.NewDuration(core.DurationValue{Unit: "day", Value: -2}))
		core.AssignSimilarDate(component, targetDate)
		core.ImplySimilarTime(component, targetDate)

	default:
		if reViimeYona.Exec(dateKeyword) != nil {
			if targetDate.GetHours() > 6 {
				targetDate = core.AddDuration(targetDate, core.NewDuration(core.DurationValue{Unit: "day", Value: -1}))
			}

			core.AssignSimilarDate(component, targetDate)
			component.Imply(core.ComponentHour, 0)
		}
	}

	if timeKeyword != "" {
		component = ExtractTimeComponents(component, timeKeyword)
	}

	return component
}
