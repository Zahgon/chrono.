// Ported from src/locales/it/parsers/ITTimeExpressionParser.ts.

package it

import "github.com/wanasit/chrono-go/internal/core"

var (
	reDiSeraNotte  = core.MustRegExp("di\\s*(?:sera|notte)", "i")
	reDiPomeriggio = core.MustRegExp("di\\s*pomeriggio", "i")
	reDiMattina    = core.MustRegExp("di\\s*mattina", "i")
)

// TimeExpressionParser is ITTimeExpressionParser ("alle 10:30", "dalle 8 alle 10", "8 di sera").
type TimeExpressionParser struct {
	core.AbstractTimeExpressionParser
}

// NewTimeExpressionParser returns a new ITTimeExpressionParser.
func NewTimeExpressionParser(strictMode bool) *TimeExpressionParser {
	p := &TimeExpressionParser{}
	p.Self = p
	p.StrictMode = strictMode
	return p
}

// FollowingPhase returns the range connector pattern.
func (p *TimeExpressionParser) FollowingPhase() string {
	return "\\s*(?:\\-|\\–|\\~|\\〜|a|fino\\s*a|alle?|\\?)\\s*"
}

// PrimaryPrefix returns the "alle"/"dalle" prefix pattern.
func (p *TimeExpressionParser) PrimaryPrefix() string {
	return "(?:(?:alle?|dalle?)\\s*)??"
}

// PrimarySuffix returns the "di sera/notte/mattina/pomeriggio" suffix pattern.
func (p *TimeExpressionParser) PrimarySuffix() string {
	return "(?:\\s*(?:di\\s*(?:sera|notte|mattina|pomeriggio)))?(?!/)(?=\\W|$)"
}

// ExtractPrimaryTimeComponents adjusts the meridiem for "di sera", "di notte", "di pomeriggio"
// and "di mattina".
func (p *TimeExpressionParser) ExtractPrimaryTimeComponents(context *core.ParsingContext, match *core.RegExpMatch) *core.ParsingComponents {
	components := p.AbstractTimeExpressionParser.ExtractPrimaryTimeComponents(context, match)
	if components == nil {
		return components
	}

	if reDiSeraNotte.Exec(match.Get(0)) != nil {
		hour := components.V(core.ComponentHour)
		if hour >= 6 && hour < 12 {
			components.Assign(core.ComponentHour, components.V(core.ComponentHour)+12)
			components.Assign(core.ComponentMeridiem, core.PM)
		} else if hour < 6 {
			components.Assign(core.ComponentMeridiem, core.AM)
		}
	}

	if reDiPomeriggio.Exec(match.Get(0)) != nil {
		components.Assign(core.ComponentMeridiem, core.PM)
		hour := components.V(core.ComponentHour)
		if hour >= 0 && hour <= 6 {
			components.Assign(core.ComponentHour, components.V(core.ComponentHour)+12)
		}
	}

	if reDiMattina.Exec(match.Get(0)) != nil {
		components.Assign(core.ComponentMeridiem, core.AM)
		hour := components.V(core.ComponentHour)
		if hour < 12 {
			components.Assign(core.ComponentHour, components.Raw(core.ComponentHour))
		}
	}

	return components.AddTag("parser/ITTimeExpressionParser")
}

// ExtractFollowingTimeComponents tags the following components.
func (p *TimeExpressionParser) ExtractFollowingTimeComponents(context *core.ParsingContext, match *core.RegExpMatch, result *core.ParsingResult) *core.ParsingComponents {
	followingComponents := p.AbstractTimeExpressionParser.ExtractFollowingTimeComponents(context, match, result)
	if followingComponents != nil {
		followingComponents.AddTag("parser/ITTimeExpressionParser")
	}
	return followingComponents
}
