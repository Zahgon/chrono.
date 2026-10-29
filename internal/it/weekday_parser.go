// Ported from src/locales/it/parsers/ITWeekdayParser.ts.

package it

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

var weekdayPattern = core.MustRegExp(
	"(?:(?:\\,|\\(|\\（)\\s*)?"+
		"(?:(questo|questa|quest'|scorso|scorsa|prossimo|prossima)\\s*)?"+
		"("+core.MatchAnyPattern(WeekdayDictionary.Keys())+"|weekend|fine\\s*settimana)"+
		"(?:\\s*(?:\\,|\\)|\\）))?"+
		"(?:\\s*(scorso|scorsa|prossimo|prossima)(?:\\s*settimana)?)?"+
		"(?=\\W|$)",
	"i",
)

const (
	weekdayPrefixGroup  = 1
	weekdayWeekdayGroup = 2
	weekdayPostfixGroup = 3
)

var reFineSettimana = core.MustRegExp("fine\\s*settimana", "")

// WeekdayParser is ITWeekdayParser ("lunedì", "venerdì prossimo", "questo weekend").
type WeekdayParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewWeekdayParser returns a new ITWeekdayParser.
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
	modifierWord := ""
	if match.Truthy(weekdayPrefixGroup) {
		modifierWord = match.Get(weekdayPrefixGroup)
	} else if match.Truthy(weekdayPostfixGroup) {
		modifierWord = match.Get(weekdayPostfixGroup)
	}
	modifierWord = js.ToLowerCase(modifierWord)

	modifier := ""
	if modifierWord == "scorso" || modifierWord == "scorsa" {
		modifier = "last"
	} else if modifierWord == "prossimo" || modifierWord == "prossima" {
		modifier = "next"
	} else if modifierWord == "questo" || modifierWord == "questa" || modifierWord == "quest'" {
		modifier = "this"
	}

	weekdayWord := js.ToLowerCase(match.Get(weekdayWeekdayGroup))
	var weekday float64
	if v, ok := WeekdayDictionary.Get(weekdayWord); ok {
		weekday = v
	} else if weekdayWord == "weekend" || reFineSettimana.Exec(weekdayWord) != nil {
		// 'Questo/prossimo weekend' means the coming Saturday, 'scorso weekend' means last Sunday.
		if modifier == "last" {
			weekday = core.Sunday
		} else {
			weekday = core.Saturday
		}
	} else {
		return nil
	}

	return core.CreateParsingComponentsAtWeekday(context.Reference, weekday, modifier)
}
