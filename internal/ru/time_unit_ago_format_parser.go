// Ported from src/locales/ru/parsers/RUTimeUnitAgoFormatParser.ts.

package ru

import "github.com/wanasit/chrono-go/internal/core"

// TimeUnitAgoFormatParser is RUTimeUnitAgoFormatParser ("3 дня назад").
type TimeUnitAgoFormatParser struct {
	AbstractParserWithLeftBoundaryChecking
}

// NewTimeUnitAgoFormatParser returns a new RUTimeUnitAgoFormatParser.
func NewTimeUnitAgoFormatParser() *TimeUnitAgoFormatParser {
	p := &TimeUnitAgoFormatParser{}
	p.Self = p
	return p
}

// InnerPatternString returns the "... назад" pattern.
func (p *TimeUnitAgoFormatParser) InnerPatternString(context *core.ParsingContext) string {
	return "(" + TimeUnitsPattern + ")\\s{0,5}назад(?=(?:\\W|$))"
}

// InnerExtract returns the components before the reference.
func (p *TimeUnitAgoFormatParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	timeUnits := ParseDuration(match.Get(1))
	outputTimeUnits := core.ReverseDuration(timeUnits)

	return core.CreateRelativeFromReference(context.Reference, outputTimeUnits)
}
