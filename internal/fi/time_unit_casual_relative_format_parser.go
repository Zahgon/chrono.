// Ported from src/locales/fi/parsers/FITimeUnitCasualRelativeFormatParser.ts.

package fi

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

var timeUnitCasualRelativePattern = core.MustRegExp(
	"(seuraava|seuraavat|seuraavien|edellinen|edelliset|edellisten|viimeiset|viimeisten|kuluneet|kuluneiden|\\+|-)\\s*("+TimeUnitsPattern+")(?=\\W|$)",
	"i",
)

var timeUnitCasualRelativePatternNoAbbr = core.MustRegExp(
	"(seuraava|seuraavat|seuraavien|edellinen|edelliset|edellisten|viimeiset|viimeisten|kuluneet|kuluneiden|\\+|-)\\s*("+TimeUnitsNoAbbrPattern+")(?=\\W|$)",
	"i",
)

// TimeUnitCasualRelativeFormatParser is FITimeUnitCasualRelativeFormatParser ("seuraavat 2
// viikkoa", "edelliset 3 päivää", "+2 tuntia").
type TimeUnitCasualRelativeFormatParser struct {
	core.AbstractParserWithWordBoundaryChecking
	allowAbbreviations bool
}

// NewTimeUnitCasualRelativeFormatParser returns a new FITimeUnitCasualRelativeFormatParser
// (the original's default allowAbbreviations is true).
func NewTimeUnitCasualRelativeFormatParser(allowAbbreviations bool) *TimeUnitCasualRelativeFormatParser {
	p := &TimeUnitCasualRelativeFormatParser{allowAbbreviations: allowAbbreviations}
	p.Self = p
	return p
}

// InnerPattern returns the pattern with or without unit abbreviations.
func (p *TimeUnitCasualRelativeFormatParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	if p.allowAbbreviations {
		return timeUnitCasualRelativePattern
	}
	return timeUnitCasualRelativePatternNoAbbr
}

// InnerExtract returns the components relative to the reference.
func (p *TimeUnitCasualRelativeFormatParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	prefix := js.ToLowerCase(match.Get(1))
	duration := ParseDuration(match.Get(2))
	if duration == nil {
		return nil
	}

	switch prefix {
	case "edellinen", "edelliset", "edellisten", "viimeiset", "viimeisten", "kuluneet", "kuluneiden", "-":
		duration = core.ReverseDuration(duration)
	}

	return core.CreateRelativeFromReference(context.Reference, duration)
}
