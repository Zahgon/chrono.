// Ported from src/locales/fi/parsers/FICasualTimeParser.ts.

package fi

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
	"github.com/wanasit/chrono-go/internal/jsdate"
)

var casualTimePattern = core.MustRegExp("(tänä\\s*)?(aamulla|aamuna|aamupäivällä|päivällä|iltapäivällä|illalla|yöllä|keskiyöllä)(?=\\W|$)", "i")

// CasualTimeParser is FICasualTimeParser ("aamulla", "tänä iltana", "keskiyöllä", ...).
type CasualTimeParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewCasualTimeParser returns a new FICasualTimeParser.
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

// ExtractTimeComponents is FICasualTimeParser.extractTimeComponents: implies the time of a
// casual time keyword on component.
func ExtractTimeComponents(component *core.ParsingComponents, timeKeywordPattern string) *core.ParsingComponents {
	switch timeKeywordPattern {
	case "aamulla", "aamuna":
		component.Imply(core.ComponentHour, 6)
		component.Imply(core.ComponentMinute, 0)
		component.Imply(core.ComponentSecond, 0)
		component.Imply(core.ComponentMeridiem, core.AM)
		component.AddTag("casualReference/morning")

	case "aamupäivällä":
		component.Imply(core.ComponentHour, 9)
		component.Imply(core.ComponentMinute, 0)
		component.Imply(core.ComponentSecond, 0)
		component.Imply(core.ComponentMeridiem, core.AM)
		component.AddTag("casualReference/morning")

	case "päivällä":
		component.Imply(core.ComponentHour, 12)
		component.Imply(core.ComponentMinute, 0)
		component.Imply(core.ComponentSecond, 0)
		component.Imply(core.ComponentMeridiem, core.AM)
		component.AddTag("casualReference/noon")

	case "iltapäivällä":
		component.Imply(core.ComponentHour, 15)
		component.Imply(core.ComponentMinute, 0)
		component.Imply(core.ComponentSecond, 0)
		component.Imply(core.ComponentMeridiem, core.PM)
		component.AddTag("casualReference/afternoon")

	case "illalla":
		component.Imply(core.ComponentHour, 18)
		component.Imply(core.ComponentMinute, 0)
		component.Imply(core.ComponentSecond, 0)
		component.Imply(core.ComponentMeridiem, core.PM)
		component.AddTag("casualReference/evening")

	case "yöllä":
		component.Imply(core.ComponentHour, 22)
		component.Imply(core.ComponentMinute, 0)
		component.Imply(core.ComponentSecond, 0)
		component.Imply(core.ComponentMeridiem, core.PM)
		component.AddTag("casualReference/evening")

	case "keskiyöllä":
		if component.V(core.ComponentHour) > 1 {
			component.AddDurationAsImplied(core.NewDuration(core.DurationValue{Unit: "day", Value: 1}))
		}

		component.Imply(core.ComponentHour, 0)
		component.Imply(core.ComponentMinute, 0)
		component.Imply(core.ComponentSecond, 0)
		component.Imply(core.ComponentMeridiem, core.AM)
		component.AddTag("casualReference/midnight")
	}
	return component
}
