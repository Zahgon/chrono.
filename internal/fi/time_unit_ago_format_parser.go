// Ported from src/locales/fi/parsers/FITimeUnitAgoFormatParser.ts.

package fi

import "github.com/wanasit/chrono-go/internal/core"

// The original builds this RegExp anew on every innerPattern() call; it has no g/y flag, so a
// shared object behaves the same.
var timeUnitAgoPattern = core.MustRegExp("("+TimeUnitsPattern+")\\s*sitten(?=\\W|$)", "i")

// TimeUnitAgoFormatParser is FITimeUnitAgoFormatParser ("3 päivää sitten").
type TimeUnitAgoFormatParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewTimeUnitAgoFormatParser returns a new FITimeUnitAgoFormatParser.
func NewTimeUnitAgoFormatParser() *TimeUnitAgoFormatParser {
	p := &TimeUnitAgoFormatParser{}
	p.Self = p
	return p
}

// InnerPattern returns the "... sitten" pattern.
func (p *TimeUnitAgoFormatParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return timeUnitAgoPattern
}

// InnerExtract returns the components before the reference.
func (p *TimeUnitAgoFormatParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	timeUnits := ParseDuration(match.Get(1))
	outputTimeUnits := core.ReverseDuration(timeUnits)

	return core.CreateRelativeFromReference(context.Reference, outputTimeUnits)
}
