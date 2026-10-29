// Ported from src/locales/de/parsers/DECasualTimeParser.ts.

package de

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
	"github.com/wanasit/chrono-go/internal/jsdate"
)

var casualTimePattern = core.MustRegExp("(diesen)?\\s*(morgen|vormittag|mittags?|nachmittag|abend|nacht|mitternacht)(?=\\W|$)", "i")

// CasualTimeParser is DECasualTimeParser ("morgen", "abend", "diesen nachmittag", ...).
type CasualTimeParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewCasualTimeParser returns a new DECasualTimeParser.
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
	timeKeywordPattern := js.ToLowerCase(match.Get(2))
	component := context.CreateParsingComponents(nil)
	core.ImplySimilarTime(component, targetDate)
	return ExtractTimeComponents(component, timeKeywordPattern)
}

// ExtractTimeComponents is the static DECasualTimeParser.extractTimeComponents: it implies the
// time of the keyword on component and returns component.
func ExtractTimeComponents(component *core.ParsingComponents, timeKeywordPattern string) *core.ParsingComponents {
	switch timeKeywordPattern {
	case "morgen":
		component.Imply(core.ComponentHour, 6)
		component.Imply(core.ComponentMinute, 0)
		component.Imply(core.ComponentSecond, 0)
		component.Imply(core.ComponentMeridiem, core.AM)
		component.AddTag("casualReference/morning")

	case "vormittag":
		component.Imply(core.ComponentHour, 9)
		component.Imply(core.ComponentMinute, 0)
		component.Imply(core.ComponentSecond, 0)
		component.Imply(core.ComponentMeridiem, core.AM)
		component.AddTag("casualReference/morning")

	case "mittag", "mittags":
		component.Imply(core.ComponentHour, 12)
		component.Imply(core.ComponentMinute, 0)
		component.Imply(core.ComponentSecond, 0)
		component.Imply(core.ComponentMeridiem, core.AM)
		component.AddTag("casualReference/noon")

	case "nachmittag":
		component.Imply(core.ComponentHour, 15)
		component.Imply(core.ComponentMinute, 0)
		component.Imply(core.ComponentSecond, 0)
		component.Imply(core.ComponentMeridiem, core.PM)
		component.AddTag("casualReference/afternoon")

	case "abend":
		component.Imply(core.ComponentHour, 18)
		component.Imply(core.ComponentMinute, 0)
		component.Imply(core.ComponentSecond, 0)
		component.Imply(core.ComponentMeridiem, core.PM)
		component.AddTag("casualReference/evening")

	case "nacht":
		component.Imply(core.ComponentHour, 22)
		component.Imply(core.ComponentMinute, 0)
		component.Imply(core.ComponentSecond, 0)
		component.Imply(core.ComponentMeridiem, core.PM)
		component.AddTag("casualReference/evening")

	case "mitternacht":
		if component.V(core.ComponentHour) > 1 {
			component.AddDurationAsImplied(dayDuration(1))
		}

		component.Imply(core.ComponentHour, 0)
		component.Imply(core.ComponentMinute, 0)
		component.Imply(core.ComponentSecond, 0)
		component.Imply(core.ComponentMeridiem, core.AM)
		component.AddTag("casualReference/midnight")
	}
	return component
}
