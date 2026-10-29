// Ported from src/locales/nl/parsers/NLTimeUnitLaterFormatParser.ts.

package nl

import "github.com/wanasit/chrono-go/internal/core"

var (
	timeUnitLaterPattern = core.MustRegExp(
		""+"("+TimeUnitsPattern+")"+"(later|na|vanaf nu|voortaan|vooruit|uit)"+"(?=(?:\\W|$))",
		"i",
	)
	timeUnitLaterStrictPattern = core.MustRegExp(""+"("+TimeUnitsPattern+")"+"(later|vanaf nu)"+"(?=(?:\\W|$))", "i")
)

const laterGroupNumTimeunits = 1

// TimeUnitLaterFormatParser is NLTimeUnitLaterFormatParser ("2 dagen later").
type TimeUnitLaterFormatParser struct {
	core.AbstractParserWithWordBoundaryChecking
	strictMode bool
}

// NewTimeUnitLaterFormatParser returns a new NLTimeUnitLaterFormatParser.
func NewTimeUnitLaterFormatParser(strictMode bool) *TimeUnitLaterFormatParser {
	p := &TimeUnitLaterFormatParser{strictMode: strictMode}
	p.Self = p
	return p
}

// InnerPattern returns the strict or casual pattern.
func (p *TimeUnitLaterFormatParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	if p.strictMode {
		return timeUnitLaterStrictPattern
	}
	return timeUnitLaterPattern
}

// InnerExtract returns the components after the reference.
func (p *TimeUnitLaterFormatParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	fragments := ParseDuration(match.Get(laterGroupNumTimeunits))
	return core.CreateRelativeFromReference(context.Reference, fragments)
}
