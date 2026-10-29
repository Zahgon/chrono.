// Ported from src/locales/nl/parsers/NLWeekdayParser.ts.

package nl

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

var weekdayPattern = core.MustRegExp(
	"(?:(?:\\,|\\(|\\（)\\s*)?"+
		"(?:op\\s*?)?"+
		"(?:(deze|vorige|volgende)\\s*(?:week\\s*)?)?"+
		"("+core.MatchAnyPattern(WeekdayDictionary.Keys())+")"+
		"(?=\\W|$)",
	"i",
)

const (
	weekdayPrefixGroup  = 1
	weekdayWeekdayGroup = 2
	// weekdayPostfixGroup is POSTFIX_GROUP: the pattern has no third group, so after the word
	// boundary shift match[3] repeats the weekday group (as in the original).
	weekdayPostfixGroup = 3
)

// WeekdayParser is NLWeekdayParser ("maandag", "volgende week vrijdag", "vorige dinsdag").
type WeekdayParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewWeekdayParser returns a new NLWeekdayParser.
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
	weekday := weekdayLookup(dayOfWeek)

	modifierWord := match.Get(weekdayPrefixGroup)
	if !match.Truthy(weekdayPrefixGroup) {
		modifierWord = match.Get(weekdayPostfixGroup)
	}
	modifierWord = js.ToLowerCase(modifierWord)

	modifier := ""
	if modifierWord == "vorige" {
		modifier = "last"
	} else if modifierWord == "volgende" {
		modifier = "next"
	} else if modifierWord == "deze" {
		modifier = "this"
	}

	return core.CreateParsingComponentsAtWeekday(context.Reference, weekday, modifier)
}
