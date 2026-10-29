// Ported from src/locales/pt/parsers/PTTimeExpressionParser.ts.

package pt

import "github.com/wanasit/chrono-go/internal/core"

// TimeExpressionParser is PTTimeExpressionParser ("às 10:30", "das 10 às 11").
type TimeExpressionParser struct {
	core.AbstractTimeExpressionParser
}

// NewTimeExpressionParser returns a new PTTimeExpressionParser (JS strictMode defaults to false).
func NewTimeExpressionParser() *TimeExpressionParser {
	p := &TimeExpressionParser{}
	p.Self = p
	p.StrictMode = false
	return p
}

// PrimaryPrefix returns the "às"/"das" prefix pattern.
func (p *TimeExpressionParser) PrimaryPrefix() string {
	return "(?:(?:ao?|às?|das|da|de|do)\\s*)?"
}

// FollowingPhase returns the range connector pattern.
func (p *TimeExpressionParser) FollowingPhase() string {
	return "\\s*(?:\\-|\\–|\\~|\\〜|a(?:o)?|\\?)\\s*"
}
