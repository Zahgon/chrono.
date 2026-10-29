// Ported from src/locales/sv/parsers/SVTimeUnitCasualRelativeFormatParser.ts.

package sv

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

var (
	timeUnitCasualRelativePattern = core.MustRegExp(
		"(denna|den här|förra|passerade|nästa|kommande|efter|\\+|-)\\s*("+TimeUnitsPattern+")(?=\\W|$)",
		"i",
	)
	timeUnitCasualRelativePatternNoAbbr = core.MustRegExp(
		"(denna|den här|förra|passerade|nästa|kommande|efter|\\+|-)\\s*("+TimeUnitsNoAbbrPattern+")(?=\\W|$)",
		"i",
	)
)

// TimeUnitCasualRelativeFormatParser is SVTimeUnitCasualRelativeFormatParser ("nästa 2 veckor",
// "förra 3 dagar", "+2 timmar").
type TimeUnitCasualRelativeFormatParser struct {
	core.AbstractParserWithWordBoundaryChecking
	allowAbbreviations bool
}

// NewTimeUnitCasualRelativeFormatParser returns a new SVTimeUnitCasualRelativeFormatParser
// (the original's default is allowAbbreviations = true).
func NewTimeUnitCasualRelativeFormatParser(allowAbbreviations bool) *TimeUnitCasualRelativeFormatParser {
	p := &TimeUnitCasualRelativeFormatParser{allowAbbreviations: allowAbbreviations}
	p.Self = p
	return p
}

// InnerPattern returns the pattern with or without abbreviated time units.
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
	case "förra", "passerade", "-":
		duration = core.ReverseDuration(duration)
	}

	return core.CreateRelativeFromReference(context.Reference, duration)
}
