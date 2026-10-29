// Ported from src/locales/uk/parsers/UKTimeUnitAgoFormatParser.ts.

package uk

import "github.com/wanasit/chrono-go/internal/core"

// TimeUnitAgoFormatParser is UKTimeUnitAgoFormatParser ("3 дні тому").
type TimeUnitAgoFormatParser struct {
	AbstractParserWithLeftBoundaryChecking
}

// NewTimeUnitAgoFormatParser returns a new UKTimeUnitAgoFormatParser.
func NewTimeUnitAgoFormatParser() *TimeUnitAgoFormatParser {
	p := &TimeUnitAgoFormatParser{}
	p.Self = p
	return p
}

// InnerPatternString returns the "... тому" pattern.
func (p *TimeUnitAgoFormatParser) InnerPatternString(context *core.ParsingContext) string {
	return "(" + TimeUnitsPattern + ")\\s{0,5}тому(?=(?:\\W|$))"
}

// InnerExtract returns the components before the reference.
func (p *TimeUnitAgoFormatParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	timeUnits := ParseDuration(match.Get(1))
	outputTimeUnits := core.ReverseDuration(timeUnits)

	return core.CreateRelativeFromReference(context.Reference, outputTimeUnits)
}
