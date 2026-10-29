// Ported from src/locales/ru/parsers/RUWeekdayParser.ts.

package ru

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

const (
	weekdayPrefixGroup  = 1
	weekdayWeekdayGroup = 2
	weekdayPostfixGroup = 3
)

// WeekdayParser is RUWeekdayParser ("в пятницу", "следующий понедельник", "в среду на прошлой
// неделе").
type WeekdayParser struct {
	AbstractParserWithLeftRightBoundaryChecking
}

// NewWeekdayParser returns a new RUWeekdayParser.
func NewWeekdayParser() *WeekdayParser {
	p := &WeekdayParser{}
	p.Self = p
	return p
}

// InnerPatternString returns the weekday pattern.
func (p *WeekdayParser) InnerPatternString(context *core.ParsingContext) string {
	return "(?:(?:,|\\(|（)\\s*)?" +
		"(?:в\\s*?)?" +
		"(?:(эту|этот|прошлый|прошлую|следующий|следующую|следующего)\\s*)?" +
		"(" + core.MatchAnyPattern(WeekdayDictionary.Keys()) + ")" +
		"(?:\\s*(?:,|\\)|）))?" +
		"(?:\\s*на\\s*(этой|прошлой|следующей)\\s*неделе)?"
}

// InnerExtract returns the components at the weekday.
func (p *WeekdayParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	dayOfWeek := js.ToLowerCase(match.Get(weekdayWeekdayGroup))
	weekday := core.Undefined
	if v, ok := WeekdayDictionary.Get(dayOfWeek); ok {
		weekday = v
	}
	prefix := match.Get(weekdayPrefixGroup)
	postfix := match.Get(weekdayPostfixGroup)
	modifierWord := prefix
	if !match.Truthy(weekdayPrefixGroup) {
		modifierWord = postfix
	}
	modifierWord = js.ToLowerCase(modifierWord)

	modifier := ""
	if modifierWord == "прошлый" || modifierWord == "прошлую" || modifierWord == "прошлой" {
		modifier = "last"
	} else if modifierWord == "следующий" ||
		modifierWord == "следующую" ||
		modifierWord == "следующей" ||
		modifierWord == "следующего" {
		modifier = "next"
	} else if modifierWord == "этот" || modifierWord == "эту" || modifierWord == "этой" {
		modifier = "this"
	}

	return core.CreateParsingComponentsAtWeekday(context.Reference, weekday, modifier)
}
