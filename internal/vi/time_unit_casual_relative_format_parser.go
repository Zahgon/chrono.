// Ported from src/locales/vi/parsers/VITimeUnitCasualRelativeFormatParser.ts.

package vi

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

// tuần này  |  tháng trước  |  năm sau  |  tuần tới  |  2 tuần trước
// NUMBER is optional so bare unit words ("tuần này") are matched
var casualUnitPattern = "(?:" + NumberPattern + "\\s{0,5})?(?:" + core.MatchAnyPattern(TimeUnitDictionary.Keys()) + ")"

var timeUnitCasualRelativePattern = core.MustRegExp(
	"(này|trước|qua|sau|tới|tiếp)\\s*("+
		casualUnitPattern+
		")"+
		"|("+
		casualUnitPattern+
		")\\s*(này|trước|qua|sau|tới|tiếp)"+
		"(?=\\W|$)",
	"i",
)

// TimeUnitCasualRelativeFormatParser is VITimeUnitCasualRelativeFormatParser ("tuần này",
// "tháng trước", "2 tuần trước").
type TimeUnitCasualRelativeFormatParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewTimeUnitCasualRelativeFormatParser returns a new VITimeUnitCasualRelativeFormatParser.
func NewTimeUnitCasualRelativeFormatParser() *TimeUnitCasualRelativeFormatParser {
	p := &TimeUnitCasualRelativeFormatParser{}
	p.Self = p
	return p
}

// InnerPattern returns the casual relative pattern.
func (p *TimeUnitCasualRelativeFormatParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return timeUnitCasualRelativePattern
}

// firstTruthy is `a || b || ""` on regex groups.
func firstTruthy(match *core.RegExpMatch, a, b int) string {
	if match.Truthy(a) {
		return match.Get(a)
	}
	if match.Truthy(b) {
		return match.Get(b)
	}
	return ""
}

// InnerExtract returns the components relative to the reference.
func (p *TimeUnitCasualRelativeFormatParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	// prefix form (trước tuần) or suffix form (tuần trước)
	modifier := js.ToLowerCase(firstTruthy(match, 1, 4))
	unitText := js.ToLowerCase(firstTruthy(match, 2, 3))

	duration := ParseDuration(unitText)
	if duration.Len() == 0 {
		// bare unit word with no number — implies 1
		unit, ok := TimeUnitDictionary.Get(unitText)
		if !ok || unit == "" {
			return nil
		}
		duration = core.NewDuration(core.DurationValue{Unit: unit, Value: 1})
	}

	if modifier == "trước" || modifier == "qua" {
		duration = core.ReverseDuration(duration)
	}
	return core.CreateRelativeFromReference(context.Reference, duration)
}
