// Ported from src/locales/en/parsers/ENCasualTimeParser.ts.

package en

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

var casualTimePattern = core.MustRegExp("(?:this)?\\s{0,3}(morning|afternoon|evening|night|midnight|midday|noon)(?=\\W|$)", "i")

// CasualTimeParser is ENCasualTimeParser ("morning", "noon", "midnight", ...).
type CasualTimeParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewCasualTimeParser returns a new ENCasualTimeParser.
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
	var component *core.ParsingComponents
	switch js.ToLowerCase(match.Get(1)) {
	case "afternoon":
		component = core.Afternoon(context.Reference, 15)
	case "evening", "night":
		component = core.Evening(context.Reference, 20)
	case "midnight":
		component = core.Midnight(context.Reference)
	case "morning":
		component = core.Morning(context.Reference, 6)
	case "noon", "midday":
		component = core.Noon(context.Reference)
	}
	if component != nil {
		component.AddTag("parser/ENCasualTimeParser")
		return component
	}
	return nil
}
