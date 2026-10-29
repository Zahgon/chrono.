// Ported from src/locales/de/parsers/DECasualDateParser.ts.

package de

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

var casualDatePattern = core.MustRegExp(
	"(jetzt|heute|morgen|übermorgen|uebermorgen|gestern|vorgestern|letzte\\s*nacht)"+
		"(?:\\s*(morgen|vormittag|mittags?|nachmittag|abend|nacht|mitternacht))?"+
		"(?=\\W|$)",
	"i",
)

var reLetzteNacht = core.MustRegExp("letzte\\s*nacht", "")

const (
	casualDateDateGroup = 1
	casualDateTimeGroup = 2
)

// CasualDateParser is DECasualDateParser ("jetzt", "heute", "morgen", "gestern abend", ...).
type CasualDateParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewCasualDateParser returns a new DECasualDateParser.
func NewCasualDateParser() *CasualDateParser {
	p := &CasualDateParser{}
	p.Self = p
	return p
}

// InnerPattern returns the casual date pattern.
func (p *CasualDateParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return casualDatePattern
}

func dayDuration(days float64) *core.Duration {
	return core.NewDuration(core.DurationValue{Unit: "day", Value: days})
}

// InnerExtract returns the components of the casual date (and time).
func (p *CasualDateParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	targetDate := context.Reference.DateWithAdjustedTimezone()
	dateKeyword := js.ToLowerCase(match.Get(casualDateDateGroup))
	timeKeyword := js.ToLowerCase(match.Get(casualDateTimeGroup))

	component := context.CreateParsingComponents(nil)
	switch dateKeyword {
	case "jetzt":
		component = core.Now(context.Reference)

	case "heute":
		component = core.Today(context.Reference)

	case "morgen":
		targetDate = core.AddDuration(targetDate, dayDuration(1))
		core.AssignSimilarDate(component, targetDate)
		core.ImplySimilarTime(component, targetDate)

	case "übermorgen", "uebermorgen":
		targetDate = core.AddDuration(targetDate, dayDuration(2))
		core.AssignSimilarDate(component, targetDate)
		core.ImplySimilarTime(component, targetDate)

	case "gestern":
		targetDate = core.AddDuration(targetDate, dayDuration(-1))
		core.AssignSimilarDate(component, targetDate)
		core.ImplySimilarTime(component, targetDate)

	case "vorgestern":
		targetDate = core.AddDuration(targetDate, dayDuration(-2))
		core.AssignSimilarDate(component, targetDate)
		core.ImplySimilarTime(component, targetDate)

	default:
		if reLetzteNacht.Exec(dateKeyword) != nil {
			if targetDate.GetHours() > 6 {
				targetDate = core.AddDuration(targetDate, dayDuration(-1))
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
