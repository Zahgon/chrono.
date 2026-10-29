// Ported from src/locales/nl/parsers/NLCasualDateParser.ts.

package nl

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

var casualDatePattern = core.MustRegExp("(nu|vandaag|morgen|morgend|gisteren)(?=\\W|$)", "i")

// CasualDateParser is NLCasualDateParser ("nu", "vandaag", "morgen", "gisteren").
type CasualDateParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewCasualDateParser returns a new NLCasualDateParser.
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
	lowerText := js.ToLowerCase(match.Get(0))
	component := context.CreateParsingComponents(nil)

	switch lowerText {
	case "nu":
		return core.Now(context.Reference)

	case "vandaag":
		return core.Today(context.Reference)

	case "morgen", "morgend":
		return core.Tomorrow(context.Reference)

	case "gisteren":
		return core.Yesterday(context.Reference)
	}

	return component
}
