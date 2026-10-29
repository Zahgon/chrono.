// Ported from src/locales/vi/parsers/VITimeUnitWithinFormatParser.ts.

package vi

import "github.com/wanasit/chrono-go/internal/core"

// trong 2 ngày  |  trong 3 tuần  |  trong vòng 1 tháng
var timeUnitWithinPattern = core.MustRegExp("(?:trong\\s*(?:v\u00f2ng\\s*)?)"+"("+TimeUnitsPattern+")(?=\\W|$)", "i")

// TimeUnitWithinFormatParser is VITimeUnitWithinFormatParser ("trong vòng 1 tháng").
type TimeUnitWithinFormatParser struct {
	core.AbstractParserWithWordBoundaryChecking
	strictMode bool
}

// NewTimeUnitWithinFormatParser returns a new VITimeUnitWithinFormatParser (strictMode does not
// change the pattern; the original's default is false).
func NewTimeUnitWithinFormatParser(strictMode bool) *TimeUnitWithinFormatParser {
	p := &TimeUnitWithinFormatParser{strictMode: strictMode}
	p.Self = p
	return p
}

// InnerPattern returns the "trong (vòng) ..." pattern (VI has no unit abbreviations, so strict
// and casual patterns are identical).
func (p *TimeUnitWithinFormatParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return timeUnitWithinPattern
}

// InnerExtract returns the components after the reference.
func (p *TimeUnitWithinFormatParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	timeUnits := ParseDuration(match.Get(1))
	if timeUnits == nil {
		return nil
	}
	return core.CreateRelativeFromReference(context.Reference, timeUnits)
}
