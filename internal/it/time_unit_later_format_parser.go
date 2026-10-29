// Ported from src/locales/it/parsers/ITTimeUnitLaterFormatParser.ts.

package it

import "github.com/wanasit/chrono-go/internal/core"

var (
	timeUnitLaterPattern = core.MustRegExp(
		"("+TimeUnitsPattern+")\\s{0,5}(?:dopo|più\\s*tardi|da\\s*adesso|da\\s*ora)"+"(?=(?:\\W|$))",
		"i",
	)
	timeUnitLaterStrictPattern = core.MustRegExp(
		"("+TimeUnitsNoAbbrPattern+")\\s{0,5}(dopo|più\\s*tardi|da\\s*adesso|da\\s*ora)(?=\\W|$)",
		"i",
	)
)

const laterGroupNumTimeunits = 1

// TimeUnitLaterFormatParser is ITTimeUnitLaterFormatParser ("3 giorni dopo", "2 ore da adesso").
type TimeUnitLaterFormatParser struct {
	core.AbstractParserWithWordBoundaryChecking
	strictMode bool
}

// NewTimeUnitLaterFormatParser returns a new ITTimeUnitLaterFormatParser.
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
	timeUnits := ParseDuration(match.Get(laterGroupNumTimeunits))
	if timeUnits == nil {
		return nil
	}
	return core.CreateRelativeFromReference(context.Reference, timeUnits)
}
