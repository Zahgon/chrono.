// Ported from src/locales/en/parsers/ENTimeUnitAgoFormatParser.ts.

package en

import "github.com/wanasit/chrono-go/internal/core"

var (
	timeUnitAgoPattern       = core.MustRegExp("("+TimeUnitsPattern+")\\s{0,5}(?:ago|before|earlier)(?=\\W|$)", "i")
	timeUnitAgoStrictPattern = core.MustRegExp("("+TimeUnitsNoAbbrPattern+")\\s{0,5}(?:ago|before|earlier)(?=\\W|$)", "i")
)

// TimeUnitAgoFormatParser is ENTimeUnitAgoFormatParser ("3 days ago").
type TimeUnitAgoFormatParser struct {
	core.AbstractParserWithWordBoundaryChecking
	strictMode bool
}

// NewTimeUnitAgoFormatParser returns a new ENTimeUnitAgoFormatParser.
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
	duration := ParseDuration(match.Get(1))
	if duration == nil {
		return nil
	}
	return core.CreateRelativeFromReference(context.Reference, core.ReverseDuration(duration))
}
