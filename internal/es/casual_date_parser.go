// Ported from src/locales/es/parsers/ESCasualDateParser.ts.

package es

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

var casualDatePattern = core.MustRegExp("(ahora|hoy|mañana|ayer)(?=\\W|$)", "i")

// CasualDateParser is ESCasualDateParser ("ahora", "hoy", "mañana", "ayer").
type CasualDateParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewCasualDateParser returns a new ESCasualDateParser.
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
	case "ahora":
		return core.Now(context.Reference)

	case "hoy":
		return core.Today(context.Reference)

	case "mañana":
		return core.Tomorrow(context.Reference)

	case "ayer":
		return core.Yesterday(context.Reference)
	}

	return component
}
