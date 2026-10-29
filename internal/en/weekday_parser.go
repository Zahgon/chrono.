// Ported from src/locales/en/parsers/ENWeekdayParser.ts.

package en

import (
	"math"

	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

var weekdayPattern = core.MustRegExp(
	"(?:(?:\\,|\\(|\\（)\\s*)?"+
		"(?:on\\s*?)?"+
		"(?:(this|last|past|next)\\s*)?"+
		"("+core.MatchAnyPattern(WeekdayDictionary.Keys())+"|weekend|weekday)"+
		"(?:\\s*(?:\\,|\\)|\\）))?"+
		"(?:\\s*(?:of\\s*)?(this|last|past|next)\\s*week)?"+
		"(?=\\W|$)",
	"i",
)

const (
	weekdayPrefixGroup  = 1
	weekdayWeekdayGroup = 2
	weekdayPostfixGroup = 3
)

// WeekdayParser is ENWeekdayParser ("Monday", "next Friday", "this weekend").
type WeekdayParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewWeekdayParser returns a new ENWeekdayParser.
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
	modifierWord := match.Get(weekdayPrefixGroup)
	if !match.Truthy(weekdayPrefixGroup) {
		modifierWord = match.Get(weekdayPostfixGroup)
	}
	modifierWord = js.ToLowerCase(modifierWord)

	modifier := ""
	if modifierWord == "last" || modifierWord == "past" {
		modifier = "last"
	} else if modifierWord == "next" {
		modifier = "next"
	} else if modifierWord == "this" {
		modifier = "this"
	}

	weekdayWord := js.ToLowerCase(match.Get(weekdayWeekdayGroup))
	var weekday float64
	if v, ok := WeekdayDictionary.Get(weekdayWord); ok {
		weekday = v
	} else if weekdayWord == "weekend" {
		// 'This/next weekend' means the coming Saturday, 'last weekend' means last Sunday.
		if modifier == "last" {
			weekday = core.Sunday
		} else {
			weekday = core.Saturday
		}
	} else if weekdayWord == "weekday" {
		// On weekend ref, this means the coming Monday or last Friday; on weekday ref, the
		// next/last working day.
		refWeekday := context.Reference.DateWithAdjustedTimezone().GetDay()
		if refWeekday == core.Sunday || refWeekday == core.Saturday {
			if modifier == "last" {
				weekday = core.Friday
			} else {
				weekday = core.Monday
			}
		} else {
			weekday = refWeekday - 1
			if modifier == "last" {
				weekday = weekday - 1
			} else {
				weekday = weekday + 1
			}
			weekday = math.Mod(weekday, 5) + 1
		}
	} else {
		return nil
	}

	return core.CreateParsingComponentsAtWeekday(context.Reference, weekday, modifier)
}
