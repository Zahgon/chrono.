// Ported from src/locales/pt/parsers/PTCasualTimeParser.ts.

package pt

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
	"github.com/wanasit/chrono-go/internal/jsdate"
)

var casualTimePattern = core.MustRegExp("(?:esta\\s*)?(manha|manhã|tarde|meia-noite|meio-dia|noite)(?=\\W|$)", "i")

// CasualTimeParser is PTCasualTimeParser ("esta manhã", "tarde", "meia-noite", ...).
type CasualTimeParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewCasualTimeParser returns a new PTCasualTimeParser.
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
	component := context.CreateParsingComponents(nil)
	switch js.ToLowerCase(match.Get(1)) {
	case "tarde":
		component.Imply(core.ComponentMeridiem, core.PM)
		component.Imply(core.ComponentHour, 15)
		component.AddTag("casualReference/afternoon")

	case "noite":
		component.Imply(core.ComponentMeridiem, core.PM)
		component.Imply(core.ComponentHour, 22)
		component.AddTag("casualReference/evening")

	case "manha", "manhã":
		component.Imply(core.ComponentMeridiem, core.AM)
		component.Imply(core.ComponentHour, 6)
		component.AddTag("casualReference/morning")

	case "meia-noite":
		nextDay := jsdate.FromTime(targetDate.Time())
		nextDay.SetDate(nextDay.GetDate() + 1)
		core.AssignSimilarDate(component, nextDay)
		core.ImplySimilarTime(component, nextDay)
		component.Imply(core.ComponentHour, 0)
		component.Imply(core.ComponentMinute, 0)
		component.Imply(core.ComponentSecond, 0)
		component.AddTag("casualReference/midnight")

	case "meio-dia":
		component.Imply(core.ComponentMeridiem, core.AM)
		component.Imply(core.ComponentHour, 12)
		component.AddTag("casualReference/noon")
	}

	return component
}
