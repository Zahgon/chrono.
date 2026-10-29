// Ported from src/locales/fr/parsers/FRCasualDateParser.ts.

package fr

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
	"github.com/wanasit/chrono-go/internal/jsdate"
)

var casualDatePattern = core.MustRegExp("(maintenant|aujourd'hui|demain|hier|cette\\s*nuit|la\\s*veille)(?=\\W|$)", "i")

var (
	reCetteNuit = core.MustRegExp("cette\\s*nuit", "")
	reLaVeille  = core.MustRegExp("la\\s*veille", "")
)

// CasualDateParser is FRCasualDateParser ("maintenant", "aujourd'hui", "demain", ...).
type CasualDateParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewCasualDateParser returns a new FRCasualDateParser.
func NewCasualDateParser() *CasualDateParser {
	p := &CasualDateParser{}
	p.Self = p
	return p
}

// InnerPattern returns the casual date pattern.
func (p *CasualDateParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return casualDatePattern
}

// InnerExtract returns the components of the casual date.
func (p *CasualDateParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	targetDate := jsdate.FromTime(context.RefDate)
	lowerText := js.ToLowerCase(match.Get(0))
	component := context.CreateParsingComponents(nil)

	switch lowerText {
	case "maintenant":
		return core.Now(context.Reference)

	case "aujourd'hui":
		return core.Today(context.Reference)

	case "hier":
		return core.Yesterday(context.Reference)

	case "demain":
		return core.Tomorrow(context.Reference)

	default:
		if reCetteNuit.Exec(lowerText) != nil {
			core.AssignSimilarDate(component, targetDate)
			component.Imply(core.ComponentHour, 22)
			component.Imply(core.ComponentMeridiem, core.PM)
		} else if reLaVeille.Exec(lowerText) != nil {
			previousDay := jsdate.FromTime(targetDate.Time())
			previousDay.SetDate(previousDay.GetDate() - 1)
			core.AssignSimilarDate(component, previousDay)
			component.Imply(core.ComponentHour, 0)
		}
	}

	return component
}
