// Ported from src/locales/it/parsers/ITTimeUnitCasualRelativeFormatParser.ts.

package it

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

var (
	timeUnitCasualRelativePattern = core.MustRegExp(
		"(?:il\\s*|la\\s*|l'\\s*)?(questo|questa|quest'|scorso|scorsa|prossimo|prossima|dopo|\\+|-)\\s*("+TimeUnitsPattern+")(?=\\W|$)",
		"i",
	)
	timeUnitCasualRelativePatternNoAbbr = core.MustRegExp(
		"(?:il\\s*|la\\s*|l'\\s*)?(questo|questa|quest'|scorso|scorsa|prossimo|prossima|dopo|\\+|-)\\s*("+TimeUnitsNoAbbrPattern+")(?=\\W|$)",
		"i",
	)
)

// TimeUnitCasualRelativeFormatParser is ITTimeUnitCasualRelativeFormatParser ("prossimi 2
// giorni", "+3 giorni").
type TimeUnitCasualRelativeFormatParser struct {
	core.AbstractParserWithWordBoundaryChecking
	allowAbbreviations bool
}

// NewTimeUnitCasualRelativeFormatParser returns a new ITTimeUnitCasualRelativeFormatParser (the
// JS default for allowAbbreviations is true).
func NewTimeUnitCasualRelativeFormatParser(allowAbbreviations bool) *TimeUnitCasualRelativeFormatParser {
	p := &TimeUnitCasualRelativeFormatParser{allowAbbreviations: allowAbbreviations}
	p.Self = p
	return p
}

// InnerPattern returns the pattern with or without abbreviations.
func (p *TimeUnitCasualRelativeFormatParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	if p.allowAbbreviations {
		return timeUnitCasualRelativePattern
	}
	return timeUnitCasualRelativePatternNoAbbr
}

// InnerExtract returns the relative components.
func (p *TimeUnitCasualRelativeFormatParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	prefix := js.ToLowerCase(match.Get(1))
	duration := ParseDuration(match.Get(2))
	if duration == nil {
		return nil
	}
	switch prefix {
	case "scorso", "scorsa", "-":
		duration = core.ReverseDuration(duration)
	}
	return core.CreateRelativeFromReference(context.Reference, duration)
}
