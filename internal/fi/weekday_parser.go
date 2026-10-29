// Ported from src/locales/fi/parsers/FIWeekdayParser.ts.

package fi

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

var weekdayPattern = core.MustRegExp(
	"(?:(?:\\,|\\(|\\（)\\s*)?"+
		"(?:(viime|edellinen|edellisenä|ensi|seuraava|seuraavana|tämä|tänä)\\s*)?"+
		"("+core.MatchAnyPattern(WeekdayDictionary.Keys())+")"+
		"(?:\\s*(?:\\,|\\)|\\）))?"+
		"(?:\\s*(viime|ensi|seuraava)\\s*viikolla)?"+
		"(?=\\W|$)",
	"i",
)

const (
	weekdayPrefixGroup  = 1
	weekdayWeekdayGroup = 2
	weekdaySuffixGroup  = 3
)

var (
	reModifierLast = core.MustRegExp("viime|edellinen|edellisenä", "")
	reModifierNext = core.MustRegExp("ensi|seuraava|seuraavana", "")
)

// WeekdayParser is FIWeekdayParser ("maanantaina", "ensi perjantai", "tiistai ensi viikolla").
type WeekdayParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewWeekdayParser returns a new FIWeekdayParser.
func NewWeekdayParser() *WeekdayParser {
	p := &WeekdayParser{}
	p.Self = p
	return p
}

// InnerPattern returns the weekday pattern.
func (p *WeekdayParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return weekdayPattern
}

// InnerExtract returns the components of the weekday.
func (p *WeekdayParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	dayOfWeek := js.ToLowerCase(match.Get(weekdayWeekdayGroup))
	offset, ok := WeekdayDictionary.Get(dayOfWeek)
	if !ok {
		offset = core.Undefined
	}

	// modifierWord = (prefix || postfix) || ""
	modifierWord := match.Get(weekdayPrefixGroup)
	if modifierWord == "" {
		modifierWord = match.Get(weekdaySuffixGroup)
	}
	modifierWord = js.ToLowerCase(modifierWord)

	modifier := ""
	if reModifierLast.Exec(modifierWord) != nil {
		modifier = "last"
	} else if reModifierNext.Exec(modifierWord) != nil {
		modifier = "next"
	}

	return core.CreateParsingComponentsAtWeekday(context.Reference, offset, modifier)
}
