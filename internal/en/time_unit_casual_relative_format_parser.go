// Ported from src/locales/en/parsers/ENTimeUnitCasualRelativeFormatParser.ts.

package en

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

var (
	timeUnitCasualRelativePattern       = core.MustRegExp("(this|last|past|next|after|\\+|-)\\s*("+TimeUnitsPattern+")(?=\\W|$)", "i")
	timeUnitCasualRelativePatternNoAbbr = core.MustRegExp("(this|last|past|next|after|\\+|-)\\s*("+TimeUnitsNoAbbrPattern+")(?=\\W|$)", "i")
)

// TimeUnitCasualRelativeFormatParser is ENTimeUnitCasualRelativeFormatParser ("next 2 weeks",
// "+3 days").
type TimeUnitCasualRelativeFormatParser struct {
	core.AbstractParserWithWordBoundaryChecking
	allowAbbreviations bool
}

// NewTimeUnitCasualRelativeFormatParser returns a new ENTimeUnitCasualRelativeFormatParser (the
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
	case "last", "past", "-":
		duration = core.ReverseDuration(duration)
	}
	return core.CreateRelativeFromReference(context.Reference, duration)
}
