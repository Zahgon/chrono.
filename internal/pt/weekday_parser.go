// Ported from src/locales/pt/parsers/PTWeekdayParser.ts.

package pt

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

var weekdayPattern = core.MustRegExp(
	"(?:(?:\\,|\\(|\\（)\\s*)?"+
		"(?:(este|esta|passado|pr[oó]ximo)\\s*)?"+
		"("+core.MatchAnyPattern(WeekdayDictionary.Keys())+")"+
		"(?:\\s*(?:\\,|\\)|\\）))?"+
		"(?:\\s*(este|esta|passado|pr[óo]ximo)\\s*semana)?"+
		"(?=\\W|\\d|$)",
	"i",
)

const (
	weekdayPrefixGroup  = 1
	weekdayWeekdayGroup = 2
	weekdayPostfixGroup = 3
)

// WeekdayParser is PTWeekdayParser ("segunda-feira", "próximo sábado").
type WeekdayParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewWeekdayParser returns a new PTWeekdayParser.
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
	if norm == "passado" {
		modifier = "this"
	} else if norm == "próximo" || norm == "proximo" {
		modifier = "next"
	} else if norm == "este" {
		modifier = "this"
	}

	return core.CreateParsingComponentsAtWeekday(context.Reference, weekday, modifier)
}
