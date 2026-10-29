// Ported from src/locales/uk/parsers/UKWeekdayParser.ts.

package uk

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

const (
	weekdayPrefixGroup  = 1
	weekdayWeekdayGroup = 2
	weekdayPostfixGroup = 3
)

// WeekdayParser is UKWeekdayParser ("у п'ятницю", "наступний понеділок", "в середу на минулому
// тижні").
type WeekdayParser struct {
	AbstractParserWithLeftRightBoundaryChecking
}

// NewWeekdayParser returns a new UKWeekdayParser.
func NewWeekdayParser() *WeekdayParser {
	p := &WeekdayParser{}
	p.Self = p
	return p
}

// InnerPatternString returns the weekday pattern.
func (p *WeekdayParser) InnerPatternString(context *core.ParsingContext) string {
	return "(?:(?:,|\\(|（)\\s*)?" +
		"(?:в\\s*?)?" +
		"(?:у\\s*?)?" +
		"(?:(цей|минулого|минулий|попередній|попереднього|наступного|наступний|наступному)\\s*)?" +
		"(" + core.MatchAnyPattern(WeekdayDictionary.Keys()) + ")" +
		"(?:\\s*(?:,|\\)|）))?" +
		"(?:\\s*(на|у|в)\\s*(цьому|минулому|наступному)\\s*тижні)?"
}

// InnerExtract returns the components at the weekday. The original lowercases with
// toLocaleLowerCase(), which equals toLowerCase() in the default (non-Turkic, non-Lithuanian)
// locale.
func (p *WeekdayParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	dayOfWeek := js.ToLowerCase(match.Get(weekdayWeekdayGroup))
	weekday := core.Undefined
	if v, ok := WeekdayDictionary.Get(dayOfWeek); ok {
		weekday = v
	}
	prefix := match.Get(weekdayPrefixGroup)
	// POSTFIX_GROUP is the (на|у|в) group, as in the original.
	postfix := match.Get(weekdayPostfixGroup)
	modifierWord := prefix
	if !match.Truthy(weekdayPrefixGroup) {
		modifierWord = postfix
	}
	modifierWord = js.ToLowerCase(modifierWord)

	modifier := ""
	if modifierWord == "минулого" ||
		modifierWord == "минулий" ||
		modifierWord == "попередній" ||
		modifierWord == "попереднього" {
		modifier = "last"
	} else if modifierWord == "наступного" || modifierWord == "наступний" {
		modifier = "next"
	} else if modifierWord == "цей" || modifierWord == "цього" || modifierWord == "цьому" {
		modifier = "this"
	}

	return core.CreateParsingComponentsAtWeekday(context.Reference, weekday, modifier)
}
