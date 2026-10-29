// Ported from src/locales/en/parsers/ENTimeUnitWithinFormatParser.ts.

package en

import "github.com/wanasit/chrono-go/internal/core"

var (
	withinPatternWithOptionalPrefix = core.MustRegExp(
		"(?:(?:within|in|for)\\s*)?"+
			"(?:(?:about|around|roughly|approximately|just)\\s*(?:~\\s*)?)?("+TimeUnitsPattern+")(?=\\W|$)",
		"i",
	)
	withinPatternWithPrefix = core.MustRegExp(
		"(?:within|in|for)\\s*"+
			"(?:(?:about|around|roughly|approximately|just)\\s*(?:~\\s*)?)?("+TimeUnitsPattern+")(?=\\W|$)",
		"i",
	)
	withinPatternWithPrefixStrict = core.MustRegExp(
		"(?:within|in|for)\\s*"+
			"(?:(?:about|around|roughly|approximately|just)\\s*(?:~\\s*)?)?("+TimeUnitsNoAbbrPattern+")(?=\\W|$)",
		"i",
	)
	reForTheUnit = core.MustRegExp("^for\\s*the\\s*\\w+", "")
)

// TimeUnitWithinFormatParser is ENTimeUnitWithinFormatParser ("in 3 days", "within 2 hours").
type TimeUnitWithinFormatParser struct {
	core.AbstractParserWithWordBoundaryChecking
	strictMode bool
}

// NewTimeUnitWithinFormatParser returns a new ENTimeUnitWithinFormatParser.
func NewTimeUnitWithinFormatParser(strictMode bool) *TimeUnitWithinFormatParser {
	p := &TimeUnitWithinFormatParser{strictMode: strictMode}
	p.Self = p
	return p
}

// InnerPattern returns the pattern for the mode and the forwardDate option.
func (p *TimeUnitWithinFormatParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	if p.strictMode {
		return withinPatternWithPrefixStrict
	}
	if context.Option.ForwardDate {
		return withinPatternWithOptionalPrefix
	}
	return withinPatternWithPrefix
}

// InnerExtract returns the components after the reference.
func (p *TimeUnitWithinFormatParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	// Exclude "for the unit" phases, e.g. "for the year"
	if reForTheUnit.Exec(match.Get(0)) != nil {
		return nil
	}
	timeUnits := ParseDuration(match.Get(1))
	if timeUnits == nil {
		return nil
	}
	return core.CreateRelativeFromReference(context.Reference, timeUnits)
}
