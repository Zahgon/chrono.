// Ported from src/locales/it/parsers/ITCasualDateParser.ts.

package it

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
	"github.com/wanasit/chrono-go/internal/jsdate"
)

var casualDatePattern = core.MustRegExp("(adesso|ora|oggi|stasera|stanotte|domani|dopodomani|ieri|ieri\\s*sera|ieri\\s*notte)(?=\\W|$)", "i")

var (
	reIeriSera  = core.MustRegExp("ieri\\s*sera", "")
	reIeriNotte = core.MustRegExp("ieri\\s*notte", "")
)

// CasualDateParser is ITCasualDateParser ("adesso", "oggi", "domani", "ieri sera", ...).
type CasualDateParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewCasualDateParser returns a new ITCasualDateParser.
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
	case "adesso", "ora":
		component = core.Now(context.Reference)

	case "oggi":
		component = core.Today(context.Reference)

	case "ieri":
		component = core.Yesterday(context.Reference)

	case "domani":
		component = core.Tomorrow(context.Reference)

	case "stasera", "stanotte":
		component = core.Tonight(context.Reference, 22)

	case "dopodomani":
		component = core.TheDayAfter(context.Reference, 2)

	default:
		if reIeriSera.Exec(lowerText) != nil || reIeriNotte.Exec(lowerText) != nil {
			if targetDate.GetHours() > 6 {
				previousDay := jsdate.FromTime(targetDate.Time())
				previousDay.SetDate(previousDay.GetDate() - 1)
				targetDate = previousDay
			}

			core.AssignSimilarDate(component, targetDate)
			component.Imply(core.ComponentHour, 0)
		}
	}
	component.AddTag("parser/ITCasualDateParser")
	return component
}
