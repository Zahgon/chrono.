// Ported from src/locales/de/parsers/DETimeUnitRelativeFormatParser.ts.

package de

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

// The original's innerPattern() builds a new, identical RegExp on every call (so the word
// boundary cache rebuilds its pattern on every parse). The pattern has no g/y flag, so this is
// unobservable and the pattern is compiled once.
var timeUnitRelativePattern = core.MustRegExp(
	"(?:\\s*((?:nächste|kommende|folgende|letzte|vergangene|vorige|vor(?:her|an)gegangene)(?:s|n|m|r)?|vor|in)\\s*)?"+
		"("+NumberPattern+")?"+
		"(?:\\s*(nächste|kommende|folgende|letzte|vergangene|vorige|vor(?:her|an)gegangene)(?:s|n|m|r)?)?"+
		"\\s*("+core.MatchAnyPattern(TimeUnitDictionary.Keys())+")(?=\\W|$)",
	"i",
)

var (
	reVor       = core.MustRegExp("vor", "")
	reLetzte    = core.MustRegExp("letzte", "")
	reVergangen = core.MustRegExp("vergangen", "")
)

// TimeUnitRelativeFormatParser is the class of DETimeUnitRelativeFormatParser.ts (named
// DETimeUnitAgoFormatParser there): "vor 2 Tagen", "in 3 Stunden", "nächste Woche".
type TimeUnitRelativeFormatParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewTimeUnitRelativeFormatParser returns a new DETimeUnitRelativeFormatParser.
func NewTimeUnitRelativeFormatParser() *TimeUnitRelativeFormatParser {
	p := &TimeUnitRelativeFormatParser{}
	p.Self = p
	return p
}

// InnerPattern returns the relative time unit pattern.
func (p *TimeUnitRelativeFormatParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return timeUnitRelativePattern
}

// InnerExtract returns the components relative to the reference (nil without a modifier).
func (p *TimeUnitRelativeFormatParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	num := 1.0
	if match.Truthy(2) {
		num = ParseNumberPattern(match.Get(2))
	}
	unit := timeUnitLookup(js.ToLowerCase(match.Get(4)))
	timeUnits := core.NewDuration()
	timeUnits.Set(unit, num)

	// Modifier
	modifier := match.Get(1)
	if !match.Truthy(1) {
		modifier = match.Get(3)
	}
	modifier = js.ToLowerCase(modifier)
	if modifier == "" {
		return nil
	}

	if reVor.Test(modifier) || reLetzte.Test(modifier) || reVergangen.Test(modifier) {
		timeUnits = core.ReverseDuration(timeUnits)
	}

	return core.CreateRelativeFromReference(context.Reference, timeUnits)
}
