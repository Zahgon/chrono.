// Ported from src/locales/fr/parsers/FRTimeUnitRelativeFormatParser.ts.

package fr

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

// The original builds a new (stateless, non-global) RegExp on every innerPattern call; one
// package-level value is equivalent.
var timeUnitRelativePattern = core.MustRegExp(
	"(?:les?|la|l'|du|des?)\\s*"+
		"("+NumberPattern+")?"+
		"(?:\\s*(prochaine?s?|derni[eè]re?s?|pass[ée]e?s?|pr[ée]c[ée]dents?|suivante?s?))?"+
		"\\s*("+core.MatchAnyPattern(TimeUnitDictionary.Keys())+")"+
		"(?:\\s*(prochaine?s?|derni[eè]re?s?|pass[ée]e?s?|pr[ée]c[ée]dents?|suivante?s?))?"+
		"(?=\\W|$)",
	"i",
)

var (
	reDerniere  = core.MustRegExp("derni[eè]re?s?", "")
	rePassee    = core.MustRegExp("pass[ée]e?s?", "")
	rePrecedent = core.MustRegExp("pr[ée]c[ée]dents?", "")
)

// TimeUnitRelativeFormatParser is FRTimeUnitRelativeFormatParser ("la semaine prochaine", "les 2
// derniers jours"). (The TypeScript class is also named FRTimeUnitAgoFormatParser.)
type TimeUnitRelativeFormatParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewTimeUnitRelativeFormatParser returns a new FRTimeUnitRelativeFormatParser.
func NewTimeUnitRelativeFormatParser() *TimeUnitRelativeFormatParser {
	p := &TimeUnitRelativeFormatParser{}
	p.Self = p
	return p
}

// InnerPattern returns the relative time unit pattern.
func (p *TimeUnitRelativeFormatParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return timeUnitRelativePattern
}

// InnerExtract returns the relative components (nil without a modifier).
func (p *TimeUnitRelativeFormatParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	num := 1.0
	if match.Truthy(1) {
		num = ParseNumberPattern(match.Get(1))
	}
	unit := timeunitLookup(js.ToLowerCase(match.Get(3)))
	timeUnits := core.NewDuration()
	timeUnits.Set(unit, num)

	// Modifier
	modifier := ""
	if match.Truthy(2) {
		modifier = match.Get(2)
	} else if match.Truthy(4) {
		modifier = match.Get(4)
	}
	modifier = js.ToLowerCase(modifier)
	if modifier == "" {
		return nil
	}

	if reDerniere.Test(modifier) || rePassee.Test(modifier) || rePrecedent.Test(modifier) {
		timeUnits = core.ReverseDuration(timeUnits)
	}

	return core.CreateRelativeFromReference(context.Reference, timeUnits)
}
