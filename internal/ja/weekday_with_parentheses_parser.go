// Ported from src/locales/ja/parsers/JPWeekdayWithParenthesesParser.ts.

package ja

import (
	"strings"

	"github.com/wanasit/chrono-go/internal/core"
)

var weekdayWithParenthesesPattern = core.MustRegExp(
	"(?:\\(|\\（)(?<weekday>"+strings.Join(WeekdayOffset.Keys(), "|")+")(?:\\)|\\）)",
	"i",
)

// WeekdayWithParenthesesParser is JPWeekdayWithParenthesesParser: a weekday in parentheses,
// e.g. "(水)", "（土）".
type WeekdayWithParenthesesParser struct {
	_ byte // non-zero size: each constructor call returns a distinct object, as `new` does
}

// NewWeekdayWithParenthesesParser returns a new JPWeekdayWithParenthesesParser.
func NewWeekdayWithParenthesesParser() *WeekdayWithParenthesesParser {
	return &WeekdayWithParenthesesParser{}
}

// Pattern returns the weekday-in-parentheses pattern.
func (p *WeekdayWithParenthesesParser) Pattern(context *core.ParsingContext) *core.RegExp {
	return weekdayWithParenthesesPattern
}

// Extract returns the components at the weekday.
func (p *WeekdayWithParenthesesParser) Extract(context *core.ParsingContext, match *core.RegExpMatch) any {
	dayOfWeek, _ := match.Group("weekday")
	offset, ok := WeekdayOffset.Get(dayOfWeek)
	if !ok {
		return nil
	}
	return core.CreateParsingComponentsAtWeekday(context.Reference, offset, "")
}
