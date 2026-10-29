// Ported from src/common/parsers/AbstractParserWithWordBoundary.ts.

package core

import "github.com/wanasit/chrono-go/internal/js"

// WordBoundaryParserImpl holds the abstract and overridable methods of
// AbstractParserWithWordBoundaryChecking. InnerPatternHasChange and PatternLeftBoundary have
// default implementations on the abstract struct (promoted through embedding); define them on the
// embedding type to override.
type WordBoundaryParserImpl interface {
	InnerPattern(context *ParsingContext) *RegExp
	InnerExtract(context *ParsingContext, match *RegExpMatch) any
	InnerPatternHasChange(context *ParsingContext, currentInnerPattern *RegExp) bool
	PatternLeftBoundary() string
}

// AbstractParserWithWordBoundaryChecking is a parser that checks for a word boundary before
// applying the inner pattern and extraction. Embed it in a parser type, implement InnerPattern
// and InnerExtract, and set Self to the embedding value in the constructor.
type AbstractParserWithWordBoundaryChecking struct {
	Self WordBoundaryParserImpl

	cachedInnerPattern *RegExp
	cachedPattern      *RegExp
}

// InnerPatternHasChange is the default: `this.innerPattern(context) !== currentInnerPattern`
// (object identity).
func (p *AbstractParserWithWordBoundaryChecking) InnerPatternHasChange(context *ParsingContext, currentInnerPattern *RegExp) bool {
	return p.Self.InnerPattern(context) != currentInnerPattern
}

// PatternLeftBoundary is the default left boundary `(\W|^)`.
func (p *AbstractParserWithWordBoundaryChecking) PatternLeftBoundary() string {
	return "(\\W|^)"
}

// Pattern returns the cached boundary-prefixed pattern, rebuilding it when the inner pattern
// changed.
func (p *AbstractParserWithWordBoundaryChecking) Pattern(context *ParsingContext) *RegExp {
	if p.cachedInnerPattern != nil {
		if !p.Self.InnerPatternHasChange(context, p.cachedInnerPattern) {
			return p.cachedPattern
		}
	}
	p.cachedInnerPattern = p.Self.InnerPattern(context)
	p.cachedPattern = MustRegExp(p.Self.PatternLeftBoundary()+p.cachedInnerPattern.Source(), p.cachedInnerPattern.Flags())
	return p.cachedPattern
}

// Extract strips the boundary group from the match (shifting the groups left by one, as the
// original does in place) and calls InnerExtract.
func (p *AbstractParserWithWordBoundaryChecking) Extract(context *ParsingContext, match *RegExpMatch) any {
	header := ""
	if match.Has(1) {
		header = match.Values[1]
	}
	headerLength := js.Len(header)
	match.Index = match.Index + headerLength
	match.Values[0] = js.SubstringFrom(match.Values[0], headerLength)
	for i := 2; i < len(match.Values); i++ {
		match.Values[i-1] = match.Values[i]
		match.Defined[i-1] = match.Defined[i]
	}
	return p.Self.InnerExtract(context, match)
}
