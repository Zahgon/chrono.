// Ported from src/locales/nl/parsers/NLCasualDateTimeParser.ts.

package nl

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
	"github.com/wanasit/chrono-go/internal/jsdate"
)

var casualDateTimePattern = core.MustRegExp("(gisteren|morgen|van)(ochtend|middag|namiddag|avond|nacht)(?=\\W|$)", "i")

const (
	casualDateTimeDateGroup      = 1
	casualDateTimeTimeOfDayGroup = 2
)

// CasualDateTimeParser is NLCasualDateTimeParser: the combined words "morgenochtend",
// "morgenmiddag", "morgennamiddag", "morgenavond", "morgennacht", "vanochtend", "vanmiddag",
// "vannamiddag", "vanavond", "vannacht", "gisterenochtend", "gisterenmiddag",
// "gisterennamiddag", "gisterenavond", "gisterennacht".
type CasualDateTimeParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewCasualDateTimeParser returns a new NLCasualDateTimeParser.
func NewCasualDateTimeParser() *CasualDateTimeParser {
	p := &CasualDateTimeParser{}
	p.Self = p
	return p
}

// InnerPattern returns the combined casual date-time pattern.
func (p *CasualDateTimeParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return casualDateTimePattern
}

// InnerExtract returns the components of the casual date and time of day.
func (p *CasualDateTimeParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	dateText := js.ToLowerCase(match.Get(casualDateTimeDateGroup))
	timeText := js.ToLowerCase(match.Get(casualDateTimeTimeOfDayGroup))
	component := context.CreateParsingComponents(nil)
	targetDate := jsdate.FromTime(context.RefDate)

	switch dateText {
	case "gisteren":
		previousDay := jsdate.FromTime(targetDate.Time())
		previousDay.SetDate(previousDay.GetDate() - 1)
		core.AssignSimilarDate(component, previousDay)
	case "van":
		core.AssignSimilarDate(component, targetDate)
	case "morgen":
		nextDay := jsdate.FromTime(targetDate.Time())
		nextDay.SetDate(nextDay.GetDate() + 1)
		core.AssignSimilarDate(component, nextDay)
		core.ImplySimilarTime(component, nextDay)
	}

	switch timeText {
	case "ochtend":
		component.Imply(core.ComponentMeridiem, core.AM)
		component.Imply(core.ComponentHour, 6)
	case "middag":
		component.Imply(core.ComponentMeridiem, core.AM)
		component.Imply(core.ComponentHour, 12)
	case "namiddag":
		component.Imply(core.ComponentMeridiem, core.PM)
		component.Imply(core.ComponentHour, 15)

	case "avond":
		component.Imply(core.ComponentMeridiem, core.PM)
		component.Imply(core.ComponentHour, 20)
	}

	return component
}
