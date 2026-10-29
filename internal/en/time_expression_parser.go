// Ported from src/locales/en/parsers/ENTimeExpressionParser.ts.

package en

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

// TimeExpressionParser is ENTimeExpressionParser.
type TimeExpressionParser struct {
	core.AbstractTimeExpressionParser
}

// NewTimeExpressionParser returns a new ENTimeExpressionParser.
func NewTimeExpressionParser(strictMode bool) *TimeExpressionParser {
	p := &TimeExpressionParser{}
	p.Self = p
	p.StrictMode = strictMode
	return p
}

// FollowingPhase returns the range connector pattern.
func (p *TimeExpressionParser) FollowingPhase() string {
	return "\\s*(?:\\-|\\–|\\~|\\〜|to|until|through|till|\\?)\\s*"
}

// PrimaryPrefix returns the "at"/"from" prefix pattern.
func (p *TimeExpressionParser) PrimaryPrefix() string {
	return "(?:(?:at|from)\\s*)??"
}

// PrimarySuffix returns the "o'clock"/"at night"/"in the morning" suffix pattern.
func (p *TimeExpressionParser) PrimarySuffix() string {
	return "(?:\\s*(?:o\\W*clock|at\\s*night|in\\s*the\\s*(?:morning|afternoon)))?(?!/)(?=\\W|$)"
}

// ExtractPrimaryTimeComponents adjusts the meridiem for "night", "afternoon" and "morning".
func (p *TimeExpressionParser) ExtractPrimaryTimeComponents(context *core.ParsingContext, match *core.RegExpMatch) *core.ParsingComponents {
	components := p.AbstractTimeExpressionParser.ExtractPrimaryTimeComponents(context, match)
	if components == nil {
		return components
	}

	if js.EndsWith(match.Get(0), "night") {
		hour := components.V(core.ComponentHour)
		if hour >= 6 && hour < 12 {
			components.Assign(core.ComponentHour, components.V(core.ComponentHour)+12)
			components.Assign(core.ComponentMeridiem, core.PM)
		} else if hour < 6 {
			components.Assign(core.ComponentMeridiem, core.AM)
		}
	}

	if js.EndsWith(match.Get(0), "afternoon") {
		components.Assign(core.ComponentMeridiem, core.PM)
		hour := components.V(core.ComponentHour)
		if hour >= 0 && hour <= 6 {
			components.Assign(core.ComponentHour, components.V(core.ComponentHour)+12)
		}
	}

	if js.EndsWith(match.Get(0), "morning") {
		components.Assign(core.ComponentMeridiem, core.AM)
		hour := components.V(core.ComponentHour)
		if hour < 12 {
			components.Assign(core.ComponentHour, components.Raw(core.ComponentHour))
		}
	}

	return components.AddTag("parser/ENTimeExpressionParser")
}

// ExtractFollowingTimeComponents tags the following components.
func (p *TimeExpressionParser) ExtractFollowingTimeComponents(context *core.ParsingContext, match *core.RegExpMatch, result *core.ParsingResult) *core.ParsingComponents {
	followingComponents := p.AbstractTimeExpressionParser.ExtractFollowingTimeComponents(context, match, result)
	if followingComponents != nil {
		followingComponents.AddTag("parser/ENTimeExpressionParser")
	}
	return followingComponents
}
