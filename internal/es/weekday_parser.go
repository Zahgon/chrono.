// Ported from src/locales/es/parsers/ESWeekdayParser.ts.

package es

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

var weekdayPattern = core.MustRegExp(
	"(?:(?:\\,|\\(|\\（)\\s*)?"+
		"(?:(este|esta|pasado|pr[oó]ximo)\\s*)?"+
		"("+core.MatchAnyPattern(WeekdayDictionary.Keys())+")"+
		"(?:\\s*(?:\\,|\\)|\\）))?"+
		"(?:\\s*(este|esta|pasado|pr[óo]ximo)\\s*semana)?"+
		"(?=\\W|\\d|$)",
	"i",
)

const (
	weekdayPrefixGroup  = 1
	weekdayWeekdayGroup = 2
	weekdayPostfixGroup = 3
)

// WeekdayParser is ESWeekdayParser ("lunes", "próximo viernes").
type WeekdayParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewWeekdayParser returns a new ESWeekdayParser.
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

	norm := ""
	if match.Truthy(weekdayPrefixGroup) {
		norm = match.Get(weekdayPrefixGroup)
	} else if match.Truthy(weekdayPostfixGroup) {
		norm = match.Get(weekdayPostfixGroup)
	}
	norm = js.ToLowerCase(norm)

	modifier := ""
	if norm == "pasado" {
		modifier = "this"
	} else if norm == "próximo" || norm == "proximo" {
		modifier = "next"
	} else if norm == "este" {
		modifier = "this"
	}

	return core.CreateParsingComponentsAtWeekday(context.Reference, weekday, modifier)
}
