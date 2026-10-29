// Ported from src/locales/ja/parsers/JPWeekdayParser.ts.

package ja

import (
	"strings"

	"github.com/wanasit/chrono-go/internal/core"
)

var weekdayPattern = core.MustRegExp(
	"((?<prefix>前の|次の|今週))?(?<weekday>"+strings.Join(WeekdayOffset.Keys(), "|")+")(?:曜日|曜)",
	"i",
)

var (
	weekdayLastPattern = core.MustRegExp("前の", "")
	weekdayNextPattern = core.MustRegExp("次の", "")
	weekdayThisPattern = core.MustRegExp("今週", "")
)

// WeekdayParser is JPWeekdayParser ("月曜日", "前の金曜", "今週水曜日", ...).
type WeekdayParser struct {
	_ byte // non-zero size: each constructor call returns a distinct object, as `new` does
}

// NewWeekdayParser returns a new JPWeekdayParser.
func NewWeekdayParser() *WeekdayParser { return &WeekdayParser{} }

// Pattern returns the weekday pattern.
func (p *WeekdayParser) Pattern(context *core.ParsingContext) *core.RegExp {
	return weekdayPattern
}

// Extract returns the components at the weekday.
func (p *WeekdayParser) Extract(context *core.ParsingContext, match *core.RegExpMatch) any {
	dayOfWeek, _ := match.Group("weekday")
	offset, ok := WeekdayOffset.Get(dayOfWeek)
	if !ok {
		return nil
	}

	prefix, _ := match.Group("prefix") // `|| ""`: an undefined group is ""

	modifier := ""
	if weekdayLastPattern.Exec(prefix) != nil {
		modifier = "last"
	} else if weekdayNextPattern.Exec(prefix) != nil {
		modifier = "next"
	} else if weekdayThisPattern.Exec(prefix) != nil {
		modifier = "this"
	}
	// TODO: handle 先週, 来週. They are different from last and next.

	return core.CreateParsingComponentsAtWeekday(context.Reference, offset, modifier)
}
