// Ported from src/locales/nl/parsers/NLTimeUnitCasualRelativeFormatParser.ts.

package nl

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

var timeUnitCasualRelativePattern = core.MustRegExp(
	"(dit|deze|vorig|afgelopen|(?:aan)?komend|over|\\+|-)e?\\s*("+TimeUnitsPattern+")(?=\\W|$)",
	"i",
)

const (
	casualRelativePrefixWordGroup   = 1
	casualRelativeTimeUnitWordGroup = 2
)

// TimeUnitCasualRelativeFormatParser is NLTimeUnitCasualRelativeFormatParser
// ("volgende 2 weken", "afgelopen 3 dagen", "+2 uur").
type TimeUnitCasualRelativeFormatParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewTimeUnitCasualRelativeFormatParser returns a new NLTimeUnitCasualRelativeFormatParser.
func NewTimeUnitCasualRelativeFormatParser() *TimeUnitCasualRelativeFormatParser {
	p := &TimeUnitCasualRelativeFormatParser{}
	p.Self = p
	return p
}

// InnerPattern returns the casual relative pattern.
func (p *TimeUnitCasualRelativeFormatParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return timeUnitCasualRelativePattern
}

// InnerExtract returns the components relative to the reference.
func (p *TimeUnitCasualRelativeFormatParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	prefix := js.ToLowerCase(match.Get(casualRelativePrefixWordGroup))
	timeUnits := ParseDuration(match.Get(casualRelativeTimeUnitWordGroup))
	switch prefix {
	case "vorig", "afgelopen", "-":
		timeUnits = core.ReverseDuration(timeUnits)
	}

	return core.CreateRelativeFromReference(context.Reference, timeUnits)
}
