// Ported from src/locales/nl/parsers/NLTimeUnitWithinFormatParser.ts.

package nl

import "github.com/wanasit/chrono-go/internal/core"

// The original's innerPattern() builds a new, identical RegExp on every call (so the word
// boundary cache rebuilds its pattern on every parse). The pattern has no g/y flag, so this is
// unobservable and the pattern is compiled once.
var timeUnitWithinPattern = core.MustRegExp("(?:binnen|in|binnen de|voor)\\s*"+"("+TimeUnitsPattern+")"+"(?=\\W|$)", "i")

// TimeUnitWithinFormatParser is NLTimeUnitWithinFormatParser ("binnen 3 dagen", "in 2 uur").
type TimeUnitWithinFormatParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewTimeUnitWithinFormatParser returns a new NLTimeUnitWithinFormatParser.
func NewTimeUnitWithinFormatParser() *TimeUnitWithinFormatParser {
	p := &TimeUnitWithinFormatParser{}
	p.Self = p
	return p
}

// InnerPattern returns the "within" pattern.
func (p *TimeUnitWithinFormatParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return timeUnitWithinPattern
}

// InnerExtract returns the components after the reference.
func (p *TimeUnitWithinFormatParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	timeUnits := ParseDuration(match.Get(1))
	return core.CreateRelativeFromReference(context.Reference, timeUnits)
}
