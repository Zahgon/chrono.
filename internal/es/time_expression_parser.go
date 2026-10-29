// Ported from src/locales/es/parsers/ESTimeExpressionParser.ts.

package es

import "github.com/wanasit/chrono-go/internal/core"

// TimeExpressionParser is ESTimeExpressionParser ("a las 10:30", "de 10 a 11").
type TimeExpressionParser struct {
	core.AbstractTimeExpressionParser
}

// NewTimeExpressionParser returns a new ESTimeExpressionParser (JS strictMode defaults to false).
func NewTimeExpressionParser() *TimeExpressionParser {
	p := &TimeExpressionParser{}
	p.Self = p
	p.StrictMode = false
	return p
}

// PrimaryPrefix returns the "a las"/"de" prefix pattern.
func (p *TimeExpressionParser) PrimaryPrefix() string {
	return "(?:(?:aslas|deslas|las?|al?|de|del)\\s*)?"
}

// FollowingPhase returns the range connector pattern.
func (p *TimeExpressionParser) FollowingPhase() string {
	return "\\s*(?:\\-|\\–|\\~|\\〜|a(?:l)?|\\?)\\s*"
}
