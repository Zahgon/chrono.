// Ported from src/locales/it/parsers/ITCasualTimeParser.ts.

package it

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

var casualTimePattern = core.MustRegExp(
	"(stamattina|stasera|stanotte|stamani)|(questa\\s*)?(mattina|pomeriggio|sera|notte|mezzanotte|mezzogiorno)(?=\\W|$)",
	"i",
)

// CasualTimeParser is ITCasualTimeParser ("stamattina", "questa sera", "mezzanotte", ...).
type CasualTimeParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewCasualTimeParser returns a new ITCasualTimeParser.
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

	// Handle compound words like "stamattina", "stasera", etc.
	if match.Truthy(1) {
		switch js.ToLowerCase(match.Get(1)) {
		case "stamattina", "stamani":
			component = core.Morning(context.Reference, 6)
		case "stasera", "stanotte":
			component = core.Evening(context.Reference, 20)
		}
	} else if match.Truthy(3) {
		switch js.ToLowerCase(match.Get(3)) {
		case "pomeriggio":
			component = core.Afternoon(context.Reference, 15)
		case "sera", "notte":
			component = core.Evening(context.Reference, 20)
		case "mezzanotte":
			component = core.Midnight(context.Reference)
		case "mattina":
			component = core.Morning(context.Reference, 6)
		case "mezzogiorno":
			component = core.Noon(context.Reference)
		}
	}
	if component != nil {
		component.AddTag("parser/ITCasualTimeParser")
		return component
	}
	return nil
}
