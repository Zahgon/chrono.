// Ported from src/locales/fr/parsers/FRWeekdayParser.ts.

package fr

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

var weekdayPattern = core.MustRegExp(
	"(?:(?:\\,|\\(|\\（)\\s*)?"+
		"(?:(?:ce)\\s*)?"+
		"("+core.MatchAnyPattern(WeekdayDictionary.Keys())+")"+
		"(?:\\s*(?:\\,|\\)|\\）))?"+
		"(?:\\s*(dernier|prochain)\\s*)?"+
		"(?=\\W|\\d|$)",
	"i",
)

const (
	weekdayWeekdayGroup = 1
	weekdayPostfixGroup = 2
)

// WeekdayParser is FRWeekdayParser ("lundi", "mardi prochain", "ce vendredi").
type WeekdayParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewWeekdayParser returns a new FRWeekdayParser.
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
	weekday, ok := WeekdayDictionary.Get(dayOfWeek)
	if !ok {
		return nil
	}

	suffix := match.Get(weekdayPostfixGroup)
	suffix = js.ToLowerCase(suffix)

	modifier := ""
	if suffix == "dernier" {
		modifier = "last"
	} else if suffix == "prochain" {
		modifier = "next"
	}

	return core.CreateParsingComponentsAtWeekday(context.Reference, weekday, modifier)
}
