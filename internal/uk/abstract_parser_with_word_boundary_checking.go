// Ported from src/locales/uk/parsers/AbstractParserWithWordBoundaryChecking.ts.

package uk

import "github.com/wanasit/chrono-go/internal/core"

// LeftBoundaryParserImpl holds the abstract methods of AbstractParserWithLeftBoundaryChecking.
type LeftBoundaryParserImpl interface {
	core.WordBoundaryParserImpl
	InnerPatternString(context *core.ParsingContext) string
}

// AbstractParserWithLeftBoundaryChecking is a word boundary parser with the Ukrainian (Unicode-aware)
// left boundary whose inner pattern is built from InnerPatternString with the "iu" flags, once per
// parser object (the inner pattern is never considered changed). Embed it, implement
// InnerPatternString and InnerExtract, and set Self to the embedding value.
type AbstractParserWithLeftBoundaryChecking struct {
	core.AbstractParserWithWordBoundaryChecking
}

func (p *AbstractParserWithLeftBoundaryChecking) innerPatternString(context *core.ParsingContext) string {
	return p.Self.(LeftBoundaryParserImpl).InnerPatternString(context)
}

// PatternLeftBoundary returns REGEX_PARTS.leftBoundary.
func (p *AbstractParserWithLeftBoundaryChecking) PatternLeftBoundary() string {
	return RegexParts.LeftBoundary
}

// InnerPattern returns a new RegExp of InnerPatternString with REGEX_PARTS.flags.
func (p *AbstractParserWithLeftBoundaryChecking) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return core.MustRegExp(p.innerPatternString(context), RegexParts.Flags)
}

// InnerPatternHasChange is always false.
func (p *AbstractParserWithLeftBoundaryChecking) InnerPatternHasChange(context *core.ParsingContext, currentInnerPattern *core.RegExp) bool {
	return false
}

// AbstractParserWithLeftRightBoundaryChecking is AbstractParserWithLeftBoundaryChecking with the
// Ukrainian right boundary appended to the inner pattern.
type AbstractParserWithLeftRightBoundaryChecking struct {
	AbstractParserWithLeftBoundaryChecking
}

// InnerPattern returns a new RegExp of InnerPatternString + REGEX_PARTS.rightBoundary with
// REGEX_PARTS.flags.
func (p *AbstractParserWithLeftRightBoundaryChecking) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return core.MustRegExp(p.innerPatternString(context)+RegexParts.RightBoundary, RegexParts.Flags)
}
