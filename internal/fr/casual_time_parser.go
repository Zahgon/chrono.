// Ported from src/locales/fr/parsers/FRCasualTimeParser.ts.

package fr

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

var casualTimePattern = core.MustRegExp("(cet?)?\\s*(matin|soir|après-midi|aprem|a midi|à minuit)(?=\\W|$)", "i")

// CasualTimeParser is FRCasualTimeParser ("ce matin", "ce soir", "à minuit", ...).
type CasualTimeParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewCasualTimeParser returns a new FRCasualTimeParser.
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
	suffixLower := js.ToLowerCase(match.Get(2))
	component := context.CreateParsingComponents(nil)

	switch suffixLower {
	case "après-midi", "aprem":
		component.Imply(core.ComponentHour, 14)
		component.Imply(core.ComponentMinute, 0)
		component.Imply(core.ComponentMeridiem, core.PM)
		component.AddTag("casualReference/afternoon")

	case "soir":
		component.Imply(core.ComponentHour, 18)
		component.Imply(core.ComponentMinute, 0)
		component.Imply(core.ComponentMeridiem, core.PM)
		component.AddTag("casualReference/evening")

	case "matin":
		component.Imply(core.ComponentHour, 8)
		component.Imply(core.ComponentMinute, 0)
		component.Imply(core.ComponentMeridiem, core.AM)
		component.AddTag("casualReference/morning")

	case "a midi":
		component.Imply(core.ComponentHour, 12)
		component.Imply(core.ComponentMinute, 0)
		component.Imply(core.ComponentMeridiem, core.AM)
		component.AddTag("casualReference/noon")

	case "à minuit":
		component.Imply(core.ComponentHour, 0)
		component.Imply(core.ComponentMeridiem, core.AM)
		component.AddTag("casualReference/midnight")
	}

	return component
}
