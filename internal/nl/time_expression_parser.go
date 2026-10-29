// Ported from src/locales/nl/parsers/NLTimeExpressionParser.ts.

package nl

import "github.com/wanasit/chrono-go/internal/core"

var reYearLikeTime = core.MustRegExp("^\\s*\\d{4}\\s*$", "")

// TimeExpressionParser is NLTimeExpressionParser ("om 10:30", "10 uur", "8:00 - 10:00").
type TimeExpressionParser struct {
	core.AbstractTimeExpressionParser
}

// NewTimeExpressionParser returns a new NLTimeExpressionParser.
func NewTimeExpressionParser(strictMode bool) *TimeExpressionParser {
	p := &TimeExpressionParser{}
	p.Self = p
	p.StrictMode = strictMode
	return p
}

// PrimaryPrefix returns the "om" prefix pattern.
func (p *TimeExpressionParser) PrimaryPrefix() string {
	return "(?:(?:om)\\s*)?"
}

// FollowingPhase returns the range connector pattern.
func (p *TimeExpressionParser) FollowingPhase() string {
	return "\\s*(?:\\-|\\–|\\~|\\〜|om|\\?)\\s*"
}

// PrimarySuffix returns the "uur" suffix pattern.
func (p *TimeExpressionParser) PrimarySuffix() string {
	return "(?:\\s*(?:uur))?(?!/)(?=\\W|$)"
}

// ExtractPrimaryTimeComponents rejects year-like matches ("2020").
func (p *TimeExpressionParser) ExtractPrimaryTimeComponents(context *core.ParsingContext, match *core.RegExpMatch) *core.ParsingComponents {
	// This looks more like a year e.g. 2020
	if reYearLikeTime.Exec(match.Get(0)) != nil {
		return nil
	}

	return p.AbstractTimeExpressionParser.ExtractPrimaryTimeComponents(context, match)
}
