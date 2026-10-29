// Ported from src/locales/fi/parsers/FITimeUnitWithinFormatParser.ts.

package fi

import "github.com/wanasit/chrono-go/internal/core"

// The original builds this RegExp anew on every innerPattern() call; it has no g/y flag, so a
// shared object behaves the same.
var timeUnitWithinPattern = core.MustRegExp("("+TimeUnitsPattern+")\\s*(?:sisällä|kuluessa|päästä)(?=\\W|$)", "i")

// TimeUnitWithinFormatParser is FITimeUnitWithinFormatParser ("2 tunnin päästä",
// "viikon sisällä").
type TimeUnitWithinFormatParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewTimeUnitWithinFormatParser returns a new FITimeUnitWithinFormatParser.
func NewTimeUnitWithinFormatParser() *TimeUnitWithinFormatParser {
	p := &TimeUnitWithinFormatParser{}
	p.Self = p
	return p
}

// InnerPattern returns the "... sisällä/kuluessa/päästä" pattern.
func (p *TimeUnitWithinFormatParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return timeUnitWithinPattern
}

// InnerExtract returns the components after the reference.
func (p *TimeUnitWithinFormatParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	timeUnits := ParseDuration(match.Get(1))
	return core.CreateRelativeFromReference(context.Reference, timeUnits)
}
