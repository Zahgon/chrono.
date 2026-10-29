// Ported from src/locales/fr/parsers/FRTimeExpressionParser.ts.

package fr

import "github.com/wanasit/chrono-go/internal/core"

var reYearLike = core.MustRegExp("^\\s*\\d{4}\\s*$", "")

// TimeExpressionParser is FRTimeExpressionParser ("à 10:30", "10:00 - 11:00").
type TimeExpressionParser struct {
	core.AbstractTimeExpressionParser
}

// NewTimeExpressionParser returns a new FRTimeExpressionParser (JS strictMode defaults to false).
func NewTimeExpressionParser() *TimeExpressionParser {
	p := &TimeExpressionParser{}
	p.Self = p
	p.StrictMode = false
	return p
}

// PrimaryPrefix returns the "à" prefix pattern.
func (p *TimeExpressionParser) PrimaryPrefix() string {
	return "(?:(?:[àa])\\s*)?"
}

// FollowingPhase returns the range connector pattern.
func (p *TimeExpressionParser) FollowingPhase() string {
	return "\\s*(?:\\-|\\–|\\~|\\〜|[àa]|\\?)\\s*"
}

// ExtractPrimaryTimeComponents skips matches that look like a year (e.g. 2020).
func (p *TimeExpressionParser) ExtractPrimaryTimeComponents(context *core.ParsingContext, match *core.RegExpMatch) *core.ParsingComponents {
	// This looks more like a year e.g. 2020
	if reYearLike.Exec(match.Get(0)) != nil {
		return nil
	}

	return p.AbstractTimeExpressionParser.ExtractPrimaryTimeComponents(context, match)
}
