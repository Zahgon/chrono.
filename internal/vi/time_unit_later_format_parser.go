// Ported from src/locales/vi/parsers/VITimeUnitLaterFormatParser.ts.

package vi

import "github.com/wanasit/chrono-go/internal/core"

// 2 ngày sau  |  3 tuần nữa  |  1 tháng tới
var timeUnitLaterPattern = core.MustRegExp(
	"("+TimeUnitsPattern+")"+"\\s{0,5}(?:sau|n\u1eefa|t\u1edbi|ti\u1ebfp)(?=\\W|$)",
	"i",
)

// TimeUnitLaterFormatParser is VITimeUnitLaterFormatParser ("2 ngày sau", "3 tuần nữa").
type TimeUnitLaterFormatParser struct {
	core.AbstractParserWithWordBoundaryChecking
	strictMode bool
}

// NewTimeUnitLaterFormatParser returns a new VITimeUnitLaterFormatParser (strictMode does not
// change the pattern; the original's default is false).
func NewTimeUnitLaterFormatParser(strictMode bool) *TimeUnitLaterFormatParser {
	p := &TimeUnitLaterFormatParser{strictMode: strictMode}
	p.Self = p
	return p
}

// InnerPattern returns the "... sau/nữa/tới/tiếp" pattern (VI has no unit abbreviations, so
// strict and casual patterns are identical).
func (p *TimeUnitLaterFormatParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return timeUnitLaterPattern
}

// InnerExtract returns the components after the reference.
func (p *TimeUnitLaterFormatParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	duration := ParseDuration(match.Get(1))
	if duration == nil {
		return nil
	}
	return core.CreateRelativeFromReference(context.Reference, duration)
}
