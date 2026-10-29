// Ported from src/locales/fr/parsers/FRTimeUnitAgoFormatParser.ts.

package fr

import "github.com/wanasit/chrono-go/internal/core"

// The original builds a new (stateless, non-global) RegExp on every innerPattern call; one
// package-level value is equivalent.
var timeUnitAgoPattern = core.MustRegExp("il y a\\s*("+TimeUnitsPattern+")(?=(?:\\W|$))", "i")

// TimeUnitAgoFormatParser is FRTimeUnitAgoFormatParser ("il y a 3 jours").
type TimeUnitAgoFormatParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewTimeUnitAgoFormatParser returns a new FRTimeUnitAgoFormatParser.
func NewTimeUnitAgoFormatParser() *TimeUnitAgoFormatParser {
	p := &TimeUnitAgoFormatParser{}
	p.Self = p
	return p
}

// InnerPattern returns the "il y a" pattern.
func (p *TimeUnitAgoFormatParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return timeUnitAgoPattern
}

// InnerExtract returns the components before the reference.
func (p *TimeUnitAgoFormatParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	timeUnits := ParseDuration(match.Get(1))
	outputTimeUnits := core.ReverseDuration(timeUnits)

	return core.CreateRelativeFromReference(context.Reference, outputTimeUnits)
}
