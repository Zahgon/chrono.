// Ported from src/locales/it/parsers/ITTimeUnitWithinFormatParser.ts.

package it

import "github.com/wanasit/chrono-go/internal/core"

var (
	withinPatternWithOptionalPrefix = core.MustRegExp(
		"(?:(?:entro|tra|fra|in|per)\\s*)?"+
			"(?:(?:circa|approssimativamente)\\s*(?:~\\s*)?)?("+TimeUnitsPattern+")(?=\\W|$)",
		"i",
	)
	withinPatternWithPrefix = core.MustRegExp(
		"(?:entro|tra|fra|in|per)\\s*"+"(?:(?:circa|approssimativamente)\\s*(?:~\\s*)?)?("+TimeUnitsPattern+")(?=\\W|$)",
		"i",
	)
	withinPatternWithPrefixStrict = core.MustRegExp(
		"(?:entro|tra|fra|in|per)\\s*"+
			"(?:(?:circa|approssimativamente)\\s*(?:~\\s*)?)?("+TimeUnitsNoAbbrPattern+")(?=\\W|$)",
		"i",
	)
	rePerTheUnit = core.MustRegExp("^per\\s*(il|la|l')\\s*\\w+", "")
)

// TimeUnitWithinFormatParser is ITTimeUnitWithinFormatParser ("tra 3 giorni", "entro 2 ore").
type TimeUnitWithinFormatParser struct {
	core.AbstractParserWithWordBoundaryChecking
	strictMode bool
}

// NewTimeUnitWithinFormatParser returns a new ITTimeUnitWithinFormatParser.
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
	// Exclude "per il/la unit" phases, e.g. "per l'anno"
	if rePerTheUnit.Exec(match.Get(0)) != nil {
		return nil
	}
	timeUnits := ParseDuration(match.Get(1))
	if timeUnits == nil {
		return nil
	}
	return core.CreateRelativeFromReference(context.Reference, timeUnits)
}
