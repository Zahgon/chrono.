// Ported from src/locales/de/parsers/DETimeExpressionParser.ts.

package de

import "github.com/wanasit/chrono-go/internal/core"

var reYearLikeTime = core.MustRegExp("^\\s*\\d{4}\\s*$", "")

// TimeExpressionParser is DETimeExpressionParser ("um 10:30", "von 8 bis 10").
type TimeExpressionParser struct {
	core.AbstractTimeExpressionParser
}

// NewTimeExpressionParser returns a new DETimeExpressionParser (constructed without arguments
// in the original, so not strict).
func NewTimeExpressionParser() *TimeExpressionParser {
	p := &TimeExpressionParser{}
	p.Self = p
	p.StrictMode = false
	return p
}

// PrimaryPrefix returns the "um"/"von" prefix pattern.
func (p *TimeExpressionParser) PrimaryPrefix() string {
	return "(?:(?:um|von)\\s*)?"
}

// FollowingPhase returns the range connector pattern.
func (p *TimeExpressionParser) FollowingPhase() string {
	return "\\s*(?:\\-|\\–|\\~|\\〜|bis)\\s*"
}

// ExtractPrimaryTimeComponents rejects year-like matches ("2020").
func (p *TimeExpressionParser) ExtractPrimaryTimeComponents(context *core.ParsingContext, match *core.RegExpMatch) *core.ParsingComponents {
	// This looks more like a year e.g. 2020
	if reYearLikeTime.Exec(match.Get(0)) != nil {
		return nil
	}

	return p.AbstractTimeExpressionParser.ExtractPrimaryTimeComponents(context, match)
}
