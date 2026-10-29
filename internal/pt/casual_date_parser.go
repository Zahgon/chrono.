// Ported from src/locales/pt/parsers/PTCasualDateParser.ts.

package pt

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

var casualDatePattern = core.MustRegExp("(agora|hoje|amanha|amanhã|ontem)(?=\\W|$)", "i")

// CasualDateParser is PTCasualDateParser ("agora", "hoje", "amanhã", "ontem").
type CasualDateParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewCasualDateParser returns a new PTCasualDateParser.
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
	case "agora":
		return core.Now(context.Reference)

	case "hoje":
		return core.Today(context.Reference)

	case "amanha", "amanhã":
		return core.Tomorrow(context.Reference)

	case "ontem":
		return core.Yesterday(context.Reference)
	}

	return component
}
