// Ported from src/locales/fr/parsers/FRSpecificTimeExpressionParser.ts.

package fr

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

// 8h10m00s
// 8h10m00
// 8h10
var specificTimeFirstPattern = core.MustRegExp(
	"(^|\\s|T)"+
		"(?:(?:[àa])\\s*)?"+
		"(\\d{1,2})(?:h|:)?"+
		"(?:(\\d{1,2})(?:m|:)?)?"+
		"(?:(\\d{1,2})(?:s|:)?)?"+
		"(?:\\s*(A\\.M\\.|P\\.M\\.|AM?|PM?))?"+
		"(?=\\W|$)",
	"i",
)

var specificTimeSecondPattern = core.MustRegExp(
	"^\\s*(\\-|\\–|\\~|\\〜|[àa]|\\?)\\s*"+
		"(\\d{1,2})(?:h|:)?"+
		"(?:(\\d{1,2})(?:m|:)?)?"+
		"(?:(\\d{1,2})(?:s|:)?)?"+
		"(?:\\s*(A\\.M\\.|P\\.M\\.|AM?|PM?))?"+
		"(?=\\W|$)",
	"i",
)

const (
	specificHourGroup     = 2
	specificMinuteGroup   = 3
	specificSecondGroup   = 4
	specificAmPmHourGroup = 5
)

var reFourDigits = core.MustRegExp("^\\d{4}$", "")

// SpecificTimeExpressionParser is FRSpecificTimeExpressionParser ("8h10m00s", "à 8h10 - 9h").
type SpecificTimeExpressionParser struct{}

// NewSpecificTimeExpressionParser returns a new FRSpecificTimeExpressionParser.
func NewSpecificTimeExpressionParser() *SpecificTimeExpressionParser {
	return &SpecificTimeExpressionParser{}
}

// Pattern returns the first time pattern.
func (p *SpecificTimeExpressionParser) Pattern(context *core.ParsingContext) *core.RegExp {
	return specificTimeFirstPattern
}

// Extract returns the result of the time (range).
func (p *SpecificTimeExpressionParser) Extract(context *core.ParsingContext, match *core.RegExpMatch) any {
	result := context.CreateParsingResult(match.Index+js.Len(match.Get(1)), js.SubstringFrom(match.Get(0), js.Len(match.Get(1))), nil, nil)

	// This looks more like a year e.g. 2020
	if reFourDigits.Exec(result.Text) != nil {
		match.Index += js.Len(match.Get(0))
		return nil
	}

	result.Start = extractTimeComponent(result.Start.Clone(), match)
	if result.Start == nil {
		match.Index += js.Len(match.Get(0))
		return nil
	}

	remainingText := context.SubstringFrom(match.Index + js.Len(match.Get(0)))
	secondMatch := specificTimeSecondPattern.Exec(remainingText)
	if secondMatch != nil {
		result.End = extractTimeComponent(result.Start.Clone(), secondMatch)
		if result.End != nil {
			result.Text += secondMatch.Get(0)
		}
	}

	return result
}

func extractTimeComponent(extractingComponents *core.ParsingComponents, match *core.RegExpMatch) *core.ParsingComponents {
	var hour float64
	var minute float64
	meridiem := -1.0 // null

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
		meridiem = core.PM
	}

	// ----- AM & PM
	if match.Has(specificAmPmHourGroup) {
		if hour > 12 {
			return nil
		}
		ampm := js.ToLowerCase(core.FirstChar(match.Get(specificAmPmHourGroup)))
		if ampm == "a" {
			meridiem = core.AM
			if hour == 12 {
				hour = 0
			}
		}

		if ampm == "p" {
			meridiem = core.PM
			if hour != 12 {
				hour += 12
			}
		}
	}

	extractingComponents.Assign(core.ComponentHour, hour)
	extractingComponents.Assign(core.ComponentMinute, minute)
	if meridiem != -1 {
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
