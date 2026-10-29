// Ported from src/locales/de/parsers/DESpecificTimeExpressionParser.ts.

package de

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

// specificTimeFirstPattern matches
//
//	8h10m00s
//	8h10m00
//	8h10
//	8 Uhr
//	8h10m00s Uhr
//	8:10 Uhr
var specificTimeFirstPattern = core.MustRegExp(
	"(^|\\s|T)"+
		"(?:(?:um|von)\\s*)?"+
		"(\\d{1,2})(?:h|:)?"+
		"(?:(\\d{1,2})(?:m|:)?)?"+
		"(?:(\\d{1,2})(?:s)?)?"+
		"(?:\\s*Uhr)?"+
		"(?:\\s*(morgens|vormittags|nachmittags|abends|nachts|am\\s+(?:Morgen|Vormittag|Nachmittag|Abend)|in\\s+der\\s+Nacht))?"+
		"(?=\\W|$)",
	"i",
)

var specificTimeSecondPattern = core.MustRegExp(
	"^\\s*(\\-|\\–|\\~|\\〜|bis(?:\\s+um)?|\\?)\\s*"+
		"(\\d{1,2})(?:h|:)?"+
		"(?:(\\d{1,2})(?:m|:)?)?"+
		"(?:(\\d{1,2})(?:s)?)?"+
		"(?:\\s*Uhr)?"+
		"(?:\\s*(morgens|vormittags|nachmittags|abends|nachts|am\\s+(?:Morgen|Vormittag|Nachmittag|Abend)|in\\s+der\\s+Nacht))?"+
		"(?=\\W|$)",
	"i",
)

const (
	specificHourGroup     = 2
	specificMinuteGroup   = 3
	specificSecondGroup   = 4
	specificAmPmHourGroup = 5
)

var (
	reFourDigitsOnly  = core.MustRegExp("^\\d{4}$", "")
	reMorgenVormittag = core.MustRegExp("morgen|vormittag", "")
	reNachmittagAbend = core.MustRegExp("nachmittag|abend", "")
	reNacht           = core.MustRegExp("nacht", "")
)

// SpecificTimeExpressionParser is DESpecificTimeExpressionParser ("8 Uhr", "8h10", "8:10 Uhr
// abends", "8 - 10 Uhr").
type SpecificTimeExpressionParser struct{}

// NewSpecificTimeExpressionParser returns a new DESpecificTimeExpressionParser.
func NewSpecificTimeExpressionParser() *SpecificTimeExpressionParser {
	return &SpecificTimeExpressionParser{}
}

// Pattern returns the first time pattern.
func (p *SpecificTimeExpressionParser) Pattern(context *core.ParsingContext) *core.RegExp {
	return specificTimeFirstPattern
}

// Extract returns the time (range) result, or nil (moving match.Index past the match).
func (p *SpecificTimeExpressionParser) Extract(context *core.ParsingContext, match *core.RegExpMatch) any {
	result := context.CreateParsingResult(match.Index+js.Len(match.Get(1)), js.SubstringFrom(match.Get(0), js.Len(match.Get(1))), nil, nil)

	// This looks more like a year e.g. 2020
	if reFourDigitsOnly.Exec(result.Text) != nil {
		match.Index += js.Len(match.Get(0))
		return nil
	}

	result.Start = extractSpecificTimeComponent(result.Start.Clone(), match)
	if result.Start == nil {
		match.Index += js.Len(match.Get(0))
		return nil
	}

	remainingText := context.SubstringFrom(match.Index + js.Len(match.Get(0)))
	secondMatch := specificTimeSecondPattern.Exec(remainingText)
	if secondMatch != nil {
		result.End = extractSpecificTimeComponent(result.Start.Clone(), secondMatch)
		if result.End != nil {
			result.Text += secondMatch.Get(0)
		}
	}

	return result
}

// extractSpecificTimeComponent is the private static extractTimeComponent (nil = null).
func extractSpecificTimeComponent(extractingComponents *core.ParsingComponents, match *core.RegExpMatch) *core.ParsingComponents {
	hour := 0.0
	minute := 0.0
	meridiem, hasMeridiem := 0.0, false

	// ----- Hours
	hour = js.ParseInt(match.Get(specificHourGroup), 0)

	// ----- Minutes
	if match.Has(specificMinuteGroup) {
		minute = js.ParseInt(match.Get(specificMinuteGroup), 0)
	}

	if minute >= 60 || hour > 24 {
		return nil
	}

	if hour >= 12 {
		meridiem, hasMeridiem = core.PM, true
	}

	// ----- AM & PM
	if match.Has(specificAmPmHourGroup) {
		if hour > 12 {
			return nil
		}
		ampm := js.ToLowerCase(match.Get(specificAmPmHourGroup))
		if reMorgenVormittag.Exec(ampm) != nil {
			meridiem, hasMeridiem = core.AM, true
			if hour == 12 {
				hour = 0
			}
		}

		if reNachmittagAbend.Exec(ampm) != nil {
			meridiem, hasMeridiem = core.PM, true
			if hour != 12 {
				hour += 12
			}
		}

		if reNacht.Exec(ampm) != nil {
			if hour == 12 {
				meridiem, hasMeridiem = core.AM, true
				hour = 0
			} else if hour < 6 {
				meridiem, hasMeridiem = core.AM, true
			} else {
				meridiem, hasMeridiem = core.PM, true
				hour += 12
			}
		}
	}

	extractingComponents.Assign(core.ComponentHour, hour)
	extractingComponents.Assign(core.ComponentMinute, minute)
	if hasMeridiem {
		extractingComponents.Assign(core.ComponentMeridiem, meridiem)
	} else {
		if hour < 12 {
			extractingComponents.Imply(core.ComponentMeridiem, core.AM)
		} else {
			extractingComponents.Imply(core.ComponentMeridiem, core.PM)
		}
	}

	// ----- Second
	if match.Has(specificSecondGroup) {
		second := js.ParseInt(match.Get(specificSecondGroup), 0)
		if second >= 60 {
			return nil
		}

		extractingComponents.Assign(core.ComponentSecond, second)
	}

	return extractingComponents
}
