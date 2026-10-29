// Ported from src/locales/nl/parsers/NLTimeUnitAgoFormatParser.ts.

package nl

import "github.com/wanasit/chrono-go/internal/core"

var (
	timeUnitAgoPattern       = core.MustRegExp(""+"("+TimeUnitsPattern+")"+"(?:geleden|voor|eerder)(?=(?:\\W|$))", "i")
	timeUnitAgoStrictPattern = core.MustRegExp(""+"("+TimeUnitsPattern+")"+"geleden(?=(?:\\W|$))", "i")
)

// TimeUnitAgoFormatParser is NLTimeUnitAgoFormatParser ("3 dagen geleden").
type TimeUnitAgoFormatParser struct {
	core.AbstractParserWithWordBoundaryChecking
	strictMode bool
}

// NewTimeUnitAgoFormatParser returns a new NLTimeUnitAgoFormatParser.
func NewTimeUnitAgoFormatParser(strictMode bool) *TimeUnitAgoFormatParser {
	p := &TimeUnitAgoFormatParser{strictMode: strictMode}
	p.Self = p
	return p
}

// InnerPattern returns the strict or casual pattern.
func (p *TimeUnitAgoFormatParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	if p.strictMode {
		return timeUnitAgoStrictPattern
	}
	return timeUnitAgoPattern
}

// InnerExtract returns the components before the reference.
func (p *TimeUnitAgoFormatParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	timeUnits := ParseDuration(match.Get(1))
	outputTimeUnits := core.ReverseDuration(timeUnits)

	return core.CreateRelativeFromReference(context.Reference, outputTimeUnits)
}
