// Ported from src/locales/fi/parsers/FITimeExpressionParser.ts.

package fi

import "github.com/wanasit/chrono-go/internal/core"

// TimeExpressionParser is FITimeExpressionParser ("klo 14.30", "kello 8-10").
type TimeExpressionParser struct {
	core.AbstractTimeExpressionParser
}

// NewTimeExpressionParser returns a new FITimeExpressionParser (constructed without arguments
// in the original, i.e. strictMode false).
func NewTimeExpressionParser() *TimeExpressionParser {
	p := &TimeExpressionParser{}
	p.Self = p
	p.StrictMode = false
	return p
}

// PrimaryPrefix returns the "klo"/"kello" prefix pattern.
func (p *TimeExpressionParser) PrimaryPrefix() string {
	return "(?:(?:klo|kello)\\s*)?"
}

// FollowingPhase returns the range connector pattern.
func (p *TimeExpressionParser) FollowingPhase() string {
	return "\\s*(?:\\-|\\–|\\~|\\〜)\\s*"
}

var reYearLike = core.MustRegExp("^\\s*\\d{4}\\s*$", "")

// ExtractPrimaryTimeComponents rejects a four-digit number (more likely a year, e.g. 2020).
func (p *TimeExpressionParser) ExtractPrimaryTimeComponents(context *core.ParsingContext, match *core.RegExpMatch) *core.ParsingComponents {
	// This looks more like a year e.g. 2020
	if reYearLike.Exec(match.Get(0)) != nil {
		return nil
	}

	return p.AbstractTimeExpressionParser.ExtractPrimaryTimeComponents(context, match)
}
