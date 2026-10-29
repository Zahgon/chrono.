// Ported from src/locales/ru/parsers/RUTimeUnitWithinFormatParser.ts.

package ru

import (
	"sync"

	"github.com/wanasit/chrono-go/internal/core"
)

var withinPattern = "(?:(?:около|примерно)\\s*(?:~\\s*)?)?(" + TimeUnitsPattern + ")" + RegexParts.RightBoundary

// The original builds a new RegExp on every innerPattern call, so the default
// innerPatternHasChange always rebuilds the boundary pattern. Both RegExps are equivalent and
// non-global (no lastIndex state), so the port compiles each variant once: the parse output is
// the same and a parse does not recompile the large time unit pattern (milliseconds per compile).
var (
	withinRegexOptionalPrefix = sync.OnceValue(func() *core.RegExp {
		return core.MustRegExp(withinPattern, RegexParts.Flags)
	})
	withinRegexWithPrefix = sync.OnceValue(func() *core.RegExp {
		return core.MustRegExp("(?:в течение|в течении)\\s*"+withinPattern, RegexParts.Flags)
	})
)

// TimeUnitWithinFormatParser is RUTimeUnitWithinFormatParser ("в течение 2 часов"; with
// forwardDate also a bare "2 часа").
type TimeUnitWithinFormatParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewTimeUnitWithinFormatParser returns a new RUTimeUnitWithinFormatParser.
func NewTimeUnitWithinFormatParser() *TimeUnitWithinFormatParser {
	p := &TimeUnitWithinFormatParser{}
	p.Self = p
	return p
}

// PatternLeftBoundary returns REGEX_PARTS.leftBoundary.
func (p *TimeUnitWithinFormatParser) PatternLeftBoundary() string {
	return RegexParts.LeftBoundary
}

// InnerPattern returns the pattern for the forwardDate option (the prefix is optional with
// forwardDate).
func (p *TimeUnitWithinFormatParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	if context.Option.ForwardDate {
		return withinRegexOptionalPrefix()
	}
	return withinRegexWithPrefix()
}

// InnerExtract returns the components after the reference.
func (p *TimeUnitWithinFormatParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	timeUnits := ParseDuration(match.Get(1))
	return core.CreateRelativeFromReference(context.Reference, timeUnits)
}
