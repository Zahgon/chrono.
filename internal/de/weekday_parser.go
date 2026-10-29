// Ported from src/locales/de/parsers/DEWeekdayParser.ts.

package de

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

var weekdayPattern = core.MustRegExp(
	"(?:(?:\\,|\\(|\\（)\\s*)?"+
		"(?:a[mn]\\s*?)?"+
		"(?:(diese[mn]|letzte[mn]|n(?:ä|ae)chste[mn])\\s*)?"+
		"("+core.MatchAnyPattern(WeekdayDictionary.Keys())+")"+
		"(?:\\s*(?:\\,|\\)|\\）))?"+
		"(?:\\s*(diese|letzte|n(?:ä|ae)chste)\\s*woche)?"+
		"(?=\\W|$)",
	"i",
)

const (
	weekdayPrefixGroup  = 1
	weekdaySuffixGroup  = 3
	weekdayWeekdayGroup = 2
)

var (
	reModLetzte = core.MustRegExp("letzte", "")
	reModChste  = core.MustRegExp("chste", "")
	reModDiese  = core.MustRegExp("diese", "")
)

// WeekdayParser is DEWeekdayParser ("Montag", "nächsten Freitag", "Di letzte Woche").
type WeekdayParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewWeekdayParser returns a new DEWeekdayParser.
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
	if reModLetzte.Exec(modifierWord) != nil {
		modifier = "last"
	} else if reModChste.Exec(modifierWord) != nil {
		modifier = "next"
	} else if reModDiese.Exec(modifierWord) != nil {
		modifier = "this"
	}

	return core.CreateParsingComponentsAtWeekday(context.Reference, offset, modifier)
}
