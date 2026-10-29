// Ported from src/locales/vi/parsers/VIWeekdayParser.ts.

package vi

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

var weekdayPattern = core.MustRegExp(
	"("+
		core.MatchAnyPattern(WeekdayDictionary.Keys())+
		")"+
		// 'sau khi' is a conjunction ("after when") — exclude via negative lookahead
		"(?:\\s*(này|tới|sau(?!\\s*khi)|qua))?"+
		"(?=\\W|$)",
	"i",
)

const (
	weekdayWeekdayGroup  = 1
	weekdayModifierGroup = 2
)

// WeekdayParser is VIWeekdayParser ("thứ hai", "chủ nhật tới", "t6 qua").
type WeekdayParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewWeekdayParser returns a new VIWeekdayParser.
func NewWeekdayParser() *WeekdayParser {
	p := &WeekdayParser{}
	p.Self = p
	return p
}

// InnerPattern returns the weekday pattern.
func (p *WeekdayParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return weekdayPattern
}

// InnerExtract returns the components of the weekday.
func (p *WeekdayParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	dowText := js.ToLowerCase(match.Get(weekdayWeekdayGroup))
	dow, ok := WeekdayDictionary.Get(dowText)
	if !ok {
		return nil
	}

	modifierType := ""
	if match.Truthy(weekdayModifierGroup) {
		m := js.ToLowerCase(match.Get(weekdayModifierGroup))
		if js.Includes(m, "tới") || js.Includes(m, "sau") {
			modifierType = "next"
		} else if js.Includes(m, "qua") {
			modifierType = "last"
		}
	}

	return core.CreateParsingComponentsAtWeekday(context.Reference, dow, modifierType)
}
