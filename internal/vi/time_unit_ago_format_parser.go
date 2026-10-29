// Ported from src/locales/vi/parsers/VITimeUnitAgoFormatParser.ts.

package vi

import "github.com/wanasit/chrono-go/internal/core"

// 2 ngày trước  |  3 tháng qua  |  1 năm trước
var timeUnitAgoPattern = core.MustRegExp("("+TimeUnitsPattern+")"+"\\s{0,5}(?:tr\u01b0\u1edbc|qua)(?=\\W|$)", "i")

// TimeUnitAgoFormatParser is VITimeUnitAgoFormatParser ("2 ngày trước").
type TimeUnitAgoFormatParser struct {
	core.AbstractParserWithWordBoundaryChecking
	strictMode bool
}

// NewTimeUnitAgoFormatParser returns a new VITimeUnitAgoFormatParser (strictMode does not change
// the pattern; the original's default is false).
func NewTimeUnitAgoFormatParser(strictMode bool) *TimeUnitAgoFormatParser {
	p := &TimeUnitAgoFormatParser{strictMode: strictMode}
	p.Self = p
	return p
}

// InnerPattern returns the "... trước/qua" pattern (VI has no unit abbreviations, so strict and
// casual patterns are identical).
func (p *TimeUnitAgoFormatParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
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
