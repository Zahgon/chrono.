// Ported from src/locales/ru/parsers/RUTimeExpressionParser.ts.

package ru

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

// TimeExpressionParser is RUTimeExpressionParser ("в 10:30", "6 вечера", "с 8 до 10 утра").
type TimeExpressionParser struct {
	core.AbstractTimeExpressionParser
}

// NewTimeExpressionParser returns a new RUTimeExpressionParser.
func NewTimeExpressionParser(strictMode bool) *TimeExpressionParser {
	p := &TimeExpressionParser{}
	p.Self = p
	p.StrictMode = strictMode
	return p
}

// PatternFlags returns REGEX_PARTS.flags.
func (p *TimeExpressionParser) PatternFlags() string {
	return RegexParts.Flags
}

// PrimaryPatternLeftBoundary returns the Unicode-aware left boundary.
func (p *TimeExpressionParser) PrimaryPatternLeftBoundary() string {
	return "(^|\\s|T|(?:[^\\p{L}\\p{N}_]))"
}

// FollowingPhase returns the range connector pattern.
func (p *TimeExpressionParser) FollowingPhase() string {
	return "\\s*(?:\\-|\\–|\\~|\\〜|до|и|по|\\?)\\s*"
}

// PrimaryPrefix returns the "в"/"с" prefix pattern.
func (p *TimeExpressionParser) PrimaryPrefix() string {
	return "(?:(?:в|с)\\s*)??"
}

// PrimarySuffix returns the "утра"/"вечера"/"после полудня" suffix pattern.
func (p *TimeExpressionParser) PrimarySuffix() string {
	return "(?:\\s*(?:утра|вечера|после полудня))?(?!\\/)" + RegexParts.RightBoundary
}

// ExtractPrimaryTimeComponents adjusts the meridiem for "вечера", "после полудня" and "утра".
func (p *TimeExpressionParser) ExtractPrimaryTimeComponents(context *core.ParsingContext, match *core.RegExpMatch) *core.ParsingComponents {
	components := p.AbstractTimeExpressionParser.ExtractPrimaryTimeComponents(context, match)
	if components != nil {
		if js.EndsWith(match.Get(0), "вечера") {
			hour := components.V(core.ComponentHour)
			if hour >= 6 && hour < 12 {
				components.Assign(core.ComponentHour, components.V(core.ComponentHour)+12)
				components.Assign(core.ComponentMeridiem, core.PM)
			} else if hour < 6 {
				components.Assign(core.ComponentMeridiem, core.AM)
			}
		}

		if js.EndsWith(match.Get(0), "после полудня") {
			components.Assign(core.ComponentMeridiem, core.PM)
			hour := components.V(core.ComponentHour)
			if hour >= 0 && hour <= 6 {
				components.Assign(core.ComponentHour, components.V(core.ComponentHour)+12)
			}
		}

		if js.EndsWith(match.Get(0), "утра") {
			components.Assign(core.ComponentMeridiem, core.AM)
			hour := components.V(core.ComponentHour)
			if hour < 12 {
				components.Assign(core.ComponentHour, components.Raw(core.ComponentHour))
			}
		}
	}

	return components
}
