// Ported from src/locales/ja/parsers/JPTimeExpressionParser.ts.

package ja

import (
	"strings"

	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

var timeFirstRegPattern = core.MustRegExp(
	"(?:"+
		"(午前|午後|A.M.|P.M.|AM|PM)"+
		")?"+
		"(?:[\\s,，、]*)"+
		"(?:([0-9０-９]+|["+
		strings.Join(Number.Keys(), "")+
		"]+)(?:\\s*)(?:時(?!間)|:|：)"+
		"(?:\\s*)"+
		"([0-9０-９]+|半|["+
		strings.Join(Number.Keys(), "")+
		"]+)?(?:\\s*)(?:分|:|：)?"+
		"(?:\\s*)"+
		"([0-9０-９]+|["+
		strings.Join(Number.Keys(), "")+
		"]+)?(?:\\s*)(?:秒)?)"+
		"(?:\\s*(A.M.|P.M.|AM?|PM?))?",
	"i",
)

var timeSecondRegPattern = core.MustRegExp(
	"(?:^\\s*(?:から|\\-|\\–|\\－|\\~|\\〜)\\s*)"+
		"(?:"+
		"(午前|午後|A.M.|P.M.|AM|PM)"+
		")?"+
		"(?:[\\s,，、]*)"+
		"(?:([0-9０-９]+|["+
		strings.Join(Number.Keys(), "")+
		"]+)(?:\\s*)(?:時|:|：)"+
		"(?:\\s*)"+
		"([0-9０-９]+|半|["+
		strings.Join(Number.Keys(), "")+
		"]+)?(?:\\s*)(?:分|:|：)?"+
		"(?:\\s*)"+
		"([0-9０-９]+|["+
		strings.Join(Number.Keys(), "")+
		"]+)?(?:\\s*)(?:秒)?)"+
		"(?:\\s*(A.M.|P.M.|AM?|PM?))?",
	"i",
)

var timeWordCharPattern = core.MustRegExp("\\w", "")

const (
	timeAmPmHourGroup1 = 1
	timeHourGroup      = 2
	timeMinuteGroup    = 3
	timeSecondGroup    = 4
	timeAmPmHourGroup2 = 5
)

// TimeExpressionParser is JPTimeExpressionParser ("午後3時", "10時30分", "9:00〜18:00", ...).
type TimeExpressionParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewTimeExpressionParser returns a new JPTimeExpressionParser.
func NewTimeExpressionParser() *TimeExpressionParser {
	p := &TimeExpressionParser{}
	p.Self = p
	return p
}

// InnerPattern returns the time expression pattern.
func (p *TimeExpressionParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return timeFirstRegPattern
}

// optionalGroup is a regex group value that may be undefined.
type optionalGroup struct {
	value   string
	defined bool
}

func groupOf(match *core.RegExpMatch, i int) optionalGroup {
	return optionalGroup{match.Get(i), match.Has(i)}
}

// amPmOf is `match[AM_PM_HOUR_GROUP_1] ?? match[AM_PM_HOUR_GROUP_2]`.
func amPmOf(match *core.RegExpMatch) optionalGroup {
	if match.Has(timeAmPmHourGroup1) {
		return groupOf(match, timeAmPmHourGroup1)
	}
	return groupOf(match, timeAmPmHourGroup2)
}

// InnerExtract returns the time (range) result.
func (p *TimeExpressionParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	// This pattern can be overlapped Ex. [12] AM, 1[2] AM
	if match.Index > 0 && timeWordCharPattern.Exec(js.CharAt(context.Text, match.Index-1)) != nil {
		return nil
	}

	result := context.CreateParsingResult(match.Index, match.Get(0), nil, nil)
	result.Start = createTimeComponents(
		context,
		groupOf(match, timeHourGroup),
		groupOf(match, timeMinuteGroup),
		groupOf(match, timeSecondGroup),
		amPmOf(match),
	)
	if result.Start == nil {
		match.Index += js.Len(match.Get(0)) // Skip over potential overlapping pattern
		return nil
	}

	// =============================================================================================
	//                  Extracting the 'to' chunk
	// =============================================================================================

	match = timeSecondRegPattern.Exec(context.SubstringFrom(result.Index + js.Len(result.Text)))
	if match == nil {
		return result
	}

	result.Text = result.Text + match.Get(0)
	result.End = createTimeComponents(
		context,
		groupOf(match, timeHourGroup),
		groupOf(match, timeMinuteGroup),
		groupOf(match, timeSecondGroup),
		amPmOf(match),
	)
	if result.End == nil {
		return nil
	}
	if !result.End.IsCertain(core.ComponentMeridiem) && result.Start.IsCertain(core.ComponentMeridiem) {
		result.End.Imply(core.ComponentMeridiem, result.Start.Raw(core.ComponentMeridiem))
		if result.Start.V(core.ComponentMeridiem) == core.PM {
			if result.Start.V(core.ComponentHour)-12 > result.End.V(core.ComponentHour) {
				// 10pm - 1 (am)
				result.End.Imply(core.ComponentMeridiem, core.AM)
			} else if result.End.V(core.ComponentHour) < 12 {
				result.End.Assign(core.ComponentHour, result.End.V(core.ComponentHour)+12)
			}
		}
	}
	if result.End.Date() < result.Start.Date() {
		result.End.Imply(core.ComponentDay, result.End.V(core.ComponentDay)+1)
	}
	return result
}

// parseNumber is `parseInt(toHankaku(text))`, falling back to jaStringToNumber(text) on NaN.
func parseNumber(text string) float64 {
	n := js.ParseInt(ToHankaku(text), 0)
	if js.IsNaN(n) {
		n = JaStringToNumber(text)
	}
	return n
}

func createTimeComponents(
	context *core.ParsingContext,
	matchHour optionalGroup,
	matchMinute optionalGroup,
	matchSecond optionalGroup,
	matchAmPm optionalGroup,
) *core.ParsingComponents {
	hour := 0.0
	meridiem := -1.0
	targetComponents := context.CreateParsingComponents(nil)
	hour = parseNumber(matchHour.value)
	if hour > 24 {
		return nil
	}

	if matchMinute.defined && matchMinute.value != "" {
		var minute float64
		if matchMinute.value == "半" {
			minute = 30
		} else {
			minute = parseNumber(matchMinute.value)
		}
		if minute >= 60 {
			return nil
		}
		targetComponents.Assign(core.ComponentMinute, minute)
	}
	if matchSecond.defined && matchSecond.value != "" {
		second := parseNumber(matchSecond.value)
		if second >= 60 {
			return nil
		}
		targetComponents.Assign(core.ComponentSecond, second)
	}

	if matchAmPm.defined && matchAmPm.value != "" {
		if hour > 12 {
			return nil
		}
		amPmString := matchAmPm.value
		if amPmString == "午前" || js.ToLowerCase(core.FirstChar(amPmString)) == "a" {
			meridiem = core.AM
			if hour == 12 {
				hour = 0
			}
		} else if amPmString == "午後" || js.ToLowerCase(core.FirstChar(amPmString)) == "p" {
			meridiem = core.PM
			if hour != 12 {
				hour += 12
			}
		}
	}

	targetComponents.Assign(core.ComponentHour, hour)

	if meridiem >= 0 {
		targetComponents.Assign(core.ComponentMeridiem, meridiem)
	} else {
		if hour < 12 {
			targetComponents.Imply(core.ComponentMeridiem, core.AM)
		} else {
			targetComponents.Imply(core.ComponentMeridiem, core.PM)
		}
	}
	return targetComponents
}
