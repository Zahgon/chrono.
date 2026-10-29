// Ported from src/locales/en/parsers/ENCasualDateParser.ts.

package en

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
	"github.com/wanasit/chrono-go/internal/jsdate"
)

var casualDatePattern = core.MustRegExp("(now|today|tonight|tomorrow|overmorrow|tmr|tmrw|yesterday|last\\s*night)(?=\\W|$)", "i")

var reLastNight = core.MustRegExp("last\\s*night", "")

// CasualDateParser is ENCasualDateParser ("now", "today", "tomorrow", ...).
type CasualDateParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewCasualDateParser returns a new ENCasualDateParser.
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
	case "now":
		component = core.Now(context.Reference)
	case "today":
		component = core.Today(context.Reference)
	case "yesterday":
		component = core.Yesterday(context.Reference)
	case "tomorrow", "tmr", "tmrw":
		component = core.Tomorrow(context.Reference)
	case "tonight":
		component = core.Tonight(context.Reference, 22)
	case "overmorrow":
		component = core.TheDayAfter(context.Reference, 2)
	default:
		if reLastNight.Exec(lowerText) != nil {
			if targetDate.GetHours() > 6 {
				previousDay := jsdate.FromTime(targetDate.Time())
				previousDay.SetDate(previousDay.GetDate() - 1)
				targetDate = previousDay
			}

			core.AssignSimilarDate(component, targetDate)
			component.Imply(core.ComponentHour, 0)
		}
	}
	component.AddTag("parser/ENCasualDateParser")
	return component
}
