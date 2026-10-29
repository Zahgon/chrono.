// Ported from src/locales/sv/parsers/SVWeekdayParser.ts.

package sv

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

var weekdayPattern = core.MustRegExp(
	"(?:(?:\\,|\\(|\\（)\\s*)?"+
		"(?:på\\s*?)?"+
		"(?:(förra|senaste|nästa|kommande)\\s*)?"+
		"("+core.MatchAnyPattern(WeekdayDictionary.Keys())+")"+
		"(?:\\s*(?:\\,|\\)|\\）))?"+
		"(?:\\s*(förra|senaste|nästa|kommande)\\s*vecka)?"+
		"(?=\\W|$)",
	"i",
)

const (
	weekdayPrefixGroup  = 1
	weekdaySuffixGroup  = 3
	weekdayWeekdayGroup = 2
)

var (
	reModLast = core.MustRegExp("förra|senaste", "")
	reModNext = core.MustRegExp("nästa|kommande", "")
)

// WeekdayParser is SVWeekdayParser ("måndag", "nästa fredag", "tisdag förra vecka").
type WeekdayParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewWeekdayParser returns a new SVWeekdayParser.
func NewWeekdayParser() *WeekdayParser {
	p := &WeekdayParser{}
	p.Self = p
	return p
}

// InnerPattern returns the weekday pattern.
func (p *WeekdayParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return weekdayPattern
}

// InnerExtract returns the components at the weekday.
func (p *WeekdayParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	dayOfWeek := js.ToLowerCase(match.Get(weekdayWeekdayGroup))
	offset := weekdayLookup(dayOfWeek)

	modifierWord := match.Get(weekdayPrefixGroup)
	if !match.Truthy(weekdayPrefixGroup) {
		modifierWord = match.Get(weekdaySuffixGroup)
	}
	modifierWord = js.ToLowerCase(modifierWord)

	modifier := ""
	if reModLast.Exec(modifierWord) != nil {
		modifier = "last"
	} else if reModNext.Exec(modifierWord) != nil {
		modifier = "next"
	}

	return core.CreateParsingComponentsAtWeekday(context.Reference, offset, modifier)
}
