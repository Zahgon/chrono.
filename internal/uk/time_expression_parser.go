// Ported from src/locales/uk/parsers/UKTimeExpressionParser.ts.

package uk

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

// TimeExpressionParser is UKTimeExpressionParser ("о 10:30", "6 вечора", "з 8 до 10 ранку").
type TimeExpressionParser struct {
	core.AbstractTimeExpressionParser
}

// NewTimeExpressionParser returns a new UKTimeExpressionParser.
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
	return "\\s*(?:\\-|\\–|\\~|\\〜|до|і|по|\\?)\\s*"
}

// PrimaryPrefix returns the "в"/"у"/"о"/"об"/"з"/"із"/"від" prefix pattern.
func (p *TimeExpressionParser) PrimaryPrefix() string {
	return "(?:(?:в|у|о|об|з|із|від)\\s*)??"
}

// PrimarySuffix returns the "ранку"/"вечора"/"по обіді"/"після обіду" suffix pattern.
func (p *TimeExpressionParser) PrimarySuffix() string {
	return "(?:\\s*(?:ранку|вечора|по обіді|після обіду))?(?!\\/)" + RegexParts.RightBoundary
}

// ExtractPrimaryTimeComponents adjusts the meridiem for "вечора", "по обіді"/"після обіду" and "ранку".
func (p *TimeExpressionParser) ExtractPrimaryTimeComponents(context *core.ParsingContext, match *core.RegExpMatch) *core.ParsingComponents {
	components := p.AbstractTimeExpressionParser.ExtractPrimaryTimeComponents(context, match)
	if components != nil {
		if js.EndsWith(match.Get(0), "вечора") {
			hour := components.V(core.ComponentHour)
			if hour >= 6 && hour < 12 {
				components.Assign(core.ComponentHour, components.V(core.ComponentHour)+12)
				components.Assign(core.ComponentMeridiem, core.PM)
			} else if hour < 6 {
				components.Assign(core.ComponentMeridiem, core.AM)
			}
		}

		if js.EndsWith(match.Get(0), "по обіді") || js.EndsWith(match.Get(0), "після обіду") {
			components.Assign(core.ComponentMeridiem, core.PM)
			hour := components.V(core.ComponentHour)
			if hour >= 0 && hour <= 6 {
				components.Assign(core.ComponentHour, components.V(core.ComponentHour)+12)
			}
		}

		if js.EndsWith(match.Get(0), "ранку") {
			components.Assign(core.ComponentMeridiem, core.AM)
			hour := components.V(core.ComponentHour)
			if hour < 12 {
				components.Assign(core.ComponentHour, components.Raw(core.ComponentHour))
			}
		}
	}

	return components
}
