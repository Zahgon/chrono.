// Ported from src/locales/uk/parsers/UKTimeUnitCasualRelativeFormatParser.ts.

package uk

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

// TimeUnitCasualRelativeFormatParser is UKTimeUnitCasualRelativeFormatParser ("через 2 години",
// "останні 3 дні", "+5 хвилин").
type TimeUnitCasualRelativeFormatParser struct {
	AbstractParserWithLeftRightBoundaryChecking
}

// NewTimeUnitCasualRelativeFormatParser returns a new UKTimeUnitCasualRelativeFormatParser.
func NewTimeUnitCasualRelativeFormatParser() *TimeUnitCasualRelativeFormatParser {
	p := &TimeUnitCasualRelativeFormatParser{}
	p.Self = p
	return p
}

// InnerPatternString returns the casual relative pattern.
func (p *TimeUnitCasualRelativeFormatParser) InnerPatternString(context *core.ParsingContext) string {
	return "(ці|останні|минулі|майбутні|наступні|після|через|\\+|-)\\s*(" + TimeUnitsPattern + ")"
}

// InnerExtract returns the components relative to the reference.
func (p *TimeUnitCasualRelativeFormatParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	prefix := js.ToLowerCase(match.Get(1))
	// match[3] is left over from the boundary group shift: the same text as match[2].
	timeUnits := ParseDuration(match.Get(3))
	switch prefix {
	case "останні", "минулі", "-":
		timeUnits = core.ReverseDuration(timeUnits)
	}

	return core.CreateRelativeFromReference(context.Reference, timeUnits)
}
