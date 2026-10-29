// Ported from src/locales/ru/parsers/RUTimeUnitCasualRelativeFormatParser.ts.

package ru

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

// TimeUnitCasualRelativeFormatParser is RUTimeUnitCasualRelativeFormatParser ("через 2 часа",
// "последние 3 дня", "+5 минут").
type TimeUnitCasualRelativeFormatParser struct {
	AbstractParserWithLeftRightBoundaryChecking
}

// NewTimeUnitCasualRelativeFormatParser returns a new RUTimeUnitCasualRelativeFormatParser.
func NewTimeUnitCasualRelativeFormatParser() *TimeUnitCasualRelativeFormatParser {
	p := &TimeUnitCasualRelativeFormatParser{}
	p.Self = p
	return p
}

// InnerPatternString returns the casual relative pattern.
func (p *TimeUnitCasualRelativeFormatParser) InnerPatternString(context *core.ParsingContext) string {
	return "(эти|последние|прошлые|следующие|после|спустя|через|\\+|-)\\s*(" + TimeUnitsPattern + ")"
}

// InnerExtract returns the components relative to the reference.
func (p *TimeUnitCasualRelativeFormatParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	prefix := js.ToLowerCase(match.Get(1))
	timeUnits := ParseDuration(match.Get(2))
	switch prefix {
	case "последние", "прошлые", "-":
		timeUnits = core.ReverseDuration(timeUnits)
	}

	return core.CreateRelativeFromReference(context.Reference, timeUnits)
}
