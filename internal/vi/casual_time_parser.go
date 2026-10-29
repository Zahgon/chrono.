// Ported from src/locales/vi/parsers/VICasualTimeParser.ts.

package vi

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
	"github.com/wanasit/chrono-go/internal/jsdate"
)

// buổi sáng  |  buổi trưa  |  buổi chiều  |  buổi tối  |  nửa đêm
var casualTimePattern = core.MustRegExp("(buổi\\s*)?(sáng sớm|sáng|trưa|chiều|tối|đêm|nửa đêm|bình minh)(?=\\W|$)", "i")

// CasualTimeParser is VICasualTimeParser ("buổi sáng", "chiều", "nửa đêm", ...).
type CasualTimeParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewCasualTimeParser returns a new VICasualTimeParser.
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
	component := context.CreateParsingComponents(nil)
	core.ImplySimilarTime(component, jsdate.FromTime(context.RefDate))
	return ExtractTimeComponents(component, js.ToLowerCase(match.Get(2)))
}

// ExtractTimeComponents is VICasualTimeParser.extractTimeComponents: implies the time of a
// casual time keyword on component.
func ExtractTimeComponents(component *core.ParsingComponents, keyword string) *core.ParsingComponents {
	switch keyword {
	case "b\u00ecnh minh", "s\u00e1ng s\u1edbm":
		component.Imply(core.ComponentHour, 6)
		component.Imply(core.ComponentMinute, 0)
		component.Imply(core.ComponentMeridiem, core.AM)
	case "s\u00e1ng":
		component.Imply(core.ComponentHour, 9)
		component.Imply(core.ComponentMinute, 0)
		component.Imply(core.ComponentMeridiem, core.AM)
	case "tr\u01b0a":
		component.Imply(core.ComponentHour, 12)
		component.Imply(core.ComponentMinute, 0)
		component.Imply(core.ComponentMeridiem, core.PM) // noon = 12:00 PM in chrono's 12-hour convention
	case "chi\u1ec1u":
		component.Imply(core.ComponentHour, 15)
		component.Imply(core.ComponentMinute, 0)
		component.Imply(core.ComponentMeridiem, core.PM)
	case "t\u1ed1i":
		component.Imply(core.ComponentHour, 19)
		component.Imply(core.ComponentMinute, 0)
		component.Imply(core.ComponentMeridiem, core.PM)
	case "\u0111\u00eam":
		component.Imply(core.ComponentHour, 22)
		component.Imply(core.ComponentMinute, 0)
		component.Imply(core.ComponentMeridiem, core.PM)
	case "n\u1eeda \u0111\u00eam":
		component.Imply(core.ComponentHour, 0)
		component.Imply(core.ComponentMinute, 0)
		component.Imply(core.ComponentMeridiem, core.AM)
	}
	return component
}
