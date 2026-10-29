// Ported from src/locales/es/parsers/ESTimeUnitWithinFormatParser.ts.

package es

import "github.com/wanasit/chrono-go/internal/core"

// The original builds a new (stateless, non-global) RegExp on every innerPattern call; one
// package-level value is equivalent.
var timeUnitWithinPattern = core.MustRegExp("(?:en|por|durante|de|dentro de)\\s*("+TimeUnitsPattern+")(?=\\W|$)", "i")

// TimeUnitWithinFormatParser is ESTimeUnitWithinFormatParser ("en 3 días").
type TimeUnitWithinFormatParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewTimeUnitWithinFormatParser returns a new ESTimeUnitWithinFormatParser.
func NewTimeUnitWithinFormatParser() *TimeUnitWithinFormatParser {
	p := &TimeUnitWithinFormatParser{}
	p.Self = p
	return p
}

// InnerPattern returns the "en/por/durante/de/dentro de" pattern.
func (p *TimeUnitWithinFormatParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return timeUnitWithinPattern
}

// InnerExtract returns the components after the reference.
func (p *TimeUnitWithinFormatParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	timeUnits := ParseDuration(match.Get(1))
	return core.CreateRelativeFromReference(context.Reference, timeUnits)
}
