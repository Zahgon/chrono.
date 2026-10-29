// Ported from src/locales/zh/hant/parsers/ZHHantTimeExpressionParser.ts.

package hant

import (
	"math"
	"strings"

	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
	"github.com/wanasit/chrono-go/internal/jsdate"
)

var timeFirstRegPattern = core.MustRegExp(
	"(?:由|從|自)?"+
		"(?:"+
		"(今|明|前|大前|後|大後|聽|昨|尋|琴)(早|朝|晚)|"+
		"(上(?:午|晝)|朝(?:早)|早(?:上)|下(?:午|晝)|晏(?:晝)|晚(?:上)|夜(?:晚)?|中(?:午)|凌(?:晨))|"+
		"(今|明|前|大前|後|大後|聽|昨|尋|琴)(?:日|天)"+
		"(?:[\\s,，]*)"+
		"(?:(上(?:午|晝)|朝(?:早)|早(?:上)|下(?:午|晝)|晏(?:晝)|晚(?:上)|夜(?:晚)?|中(?:午)|凌(?:晨)))?"+
		")?"+
		"(?:[\\s,，]*)"+
		"(?:(\\d+|["+
		strings.Join(Number.Keys(), "")+
		"]+)(?:\\s*)(?:點|時|:|：)"+
		"(?:\\s*)"+
		"(\\d+|半|正|整|["+
		strings.Join(Number.Keys(), "")+
		"]+)?(?:\\s*)(?:分|:|：)?"+
		"(?:\\s*)"+
		"(\\d+|["+
		strings.Join(Number.Keys(), "")+
		"]+)?(?:\\s*)(?:秒)?)"+
		"(?:\\s*(A.M.|P.M.|AM?|PM?))?",
	"i",
)

var timeSecondRegPattern = core.MustRegExp(
	"(?:^\\s*(?:到|至|\\-|\\–|\\~|\\〜)\\s*)"+
		"(?:"+
		"(今|明|前|大前|後|大後|聽|昨|尋|琴)(早|朝|晚)|"+
		"(上(?:午|晝)|朝(?:早)|早(?:上)|下(?:午|晝)|晏(?:晝)|晚(?:上)|夜(?:晚)?|中(?:午)|凌(?:晨))|"+
		"(今|明|前|大前|後|大後|聽|昨|尋|琴)(?:日|天)"+
		"(?:[\\s,，]*)"+
		"(?:(上(?:午|晝)|朝(?:早)|早(?:上)|下(?:午|晝)|晏(?:晝)|晚(?:上)|夜(?:晚)?|中(?:午)|凌(?:晨)))?"+
		")?"+
		"(?:[\\s,，]*)"+
		"(?:(\\d+|["+
		strings.Join(Number.Keys(), "")+
		"]+)(?:\\s*)(?:點|時|:|：)"+
		"(?:\\s*)"+
		"(\\d+|半|正|整|["+
		strings.Join(Number.Keys(), "")+
		"]+)?(?:\\s*)(?:分|:|：)?"+
		"(?:\\s*)"+
		"(\\d+|["+
		strings.Join(Number.Keys(), "")+
		"]+)?(?:\\s*)(?:秒)?)"+
		"(?:\\s*(A.M.|P.M.|AM?|PM?))?",
	"i",
)

var (
	timeWordCharPattern   = core.MustRegExp("\\w", "")
	timeDigitsOnlyPattern = core.MustRegExp("^\\d+$", "")
)

const (
	timeDayGroup1        = 1
	timeZhAmPmHourGroup1 = 2
	timeZhAmPmHourGroup2 = 3
	timeDayGroup3        = 4
	timeZhAmPmHourGroup3 = 5
	timeHourGroup        = 6
	timeMinuteGroup      = 7
	timeSecondGroup      = 8
	timeAmPmHourGroup    = 9
)

// TimeExpressionParser is ZHHantTimeExpressionParser ("下午3點", "聽日朝早8點半", "6點到9點", ...).
type TimeExpressionParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewTimeExpressionParser returns a new ZHHantTimeExpressionParser.
func NewTimeExpressionParser() *TimeExpressionParser {
	p := &TimeExpressionParser{}
	p.Self = p
	return p
}

// PatternLeftBoundary returns a capturing group to ensure that the match index is correct in the
// base class while avoiding matching CJK characters as word boundaries.
func (p *TimeExpressionParser) PatternLeftBoundary() string {
	return "()"
}

// InnerPattern returns the time expression pattern.
func (p *TimeExpressionParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return timeFirstRegPattern
}

// InnerExtract returns the time (range) result.
func (p *TimeExpressionParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	// This pattern can be overlaped Ex. [12] AM, 1[2] AM
	if match.Index > 0 && timeWordCharPattern.Exec(js.CharAt(context.Text, match.Index-1)) != nil {
		return nil
	}

	result := context.CreateParsingResult(match.Index, match.Get(0), nil, nil)
	startMoment := jsdate.FromTime(jsdate.FromTime(context.Reference.Instant).Time())

	// ----- Day
	if match.Truthy(timeDayGroup1) {
		day1 := match.Get(timeDayGroup1)
		if day1 == "明" || day1 == "聽" {
			// Check not "Tomorrow" on late night
			if jsdate.FromTime(context.Reference.Instant).GetHours() > 1 {
				startMoment.SetDate(startMoment.GetDate() + 1)
			}
		} else if day1 == "昨" || day1 == "尋" || day1 == "琴" {
			startMoment.SetDate(startMoment.GetDate() - 1)
		} else if day1 == "前" {
			startMoment.SetDate(startMoment.GetDate() - 2)
		} else if day1 == "大前" {
			startMoment.SetDate(startMoment.GetDate() - 3)
		} else if day1 == "後" {
			startMoment.SetDate(startMoment.GetDate() + 2)
		} else if day1 == "大後" {
			startMoment.SetDate(startMoment.GetDate() + 3)
		}
		result.Start.Assign(core.ComponentDay, startMoment.GetDate())
		result.Start.Assign(core.ComponentMonth, startMoment.GetMonth()+1)
		result.Start.Assign(core.ComponentYear, startMoment.GetFullYear())
	} else if match.Truthy(timeDayGroup3) {
		day3 := match.Get(timeDayGroup3)
		if day3 == "明" || day3 == "聽" {
			startMoment.SetDate(startMoment.GetDate() + 1)
		} else if day3 == "昨" || day3 == "尋" || day3 == "琴" {
			startMoment.SetDate(startMoment.GetDate() - 1)
		} else if day3 == "前" {
			startMoment.SetDate(startMoment.GetDate() - 2)
		} else if day3 == "大前" {
			startMoment.SetDate(startMoment.GetDate() - 3)
		} else if day3 == "後" {
			startMoment.SetDate(startMoment.GetDate() + 2)
		} else if day3 == "大後" {
			startMoment.SetDate(startMoment.GetDate() + 3)
		}
		result.Start.Assign(core.ComponentDay, startMoment.GetDate())
		result.Start.Assign(core.ComponentMonth, startMoment.GetMonth()+1)
		result.Start.Assign(core.ComponentYear, startMoment.GetFullYear())
	} else {
		result.Start.Imply(core.ComponentDay, startMoment.GetDate())
		result.Start.Imply(core.ComponentMonth, startMoment.GetMonth()+1)
		result.Start.Imply(core.ComponentYear, startMoment.GetFullYear())
	}

	hour := 0.0
	minute := 0.0
	meridiem := -1.0

	// ----- Second
	if match.Truthy(timeSecondGroup) {
		second := js.ParseInt(match.Get(timeSecondGroup), 0)
		if js.IsNaN(second) {
			second = ZhStringToNumber(match.Get(timeSecondGroup))
		}
		if second >= 60 {
			return nil
		}
		result.Start.Assign(core.ComponentSecond, second)
	}

	hour = js.ParseInt(match.Get(timeHourGroup), 0)
	if js.IsNaN(hour) {
		hour = ZhStringToNumber(match.Get(timeHourGroup))
	}

	// ----- Minutes
	if match.Truthy(timeMinuteGroup) {
		if match.Get(timeMinuteGroup) == "半" {
			minute = 30
		} else if match.Get(timeMinuteGroup) == "正" || match.Get(timeMinuteGroup) == "整" {
			minute = 0
		} else {
			minute = js.ParseInt(match.Get(timeMinuteGroup), 0)
			if js.IsNaN(minute) {
				minute = ZhStringToNumber(match.Get(timeMinuteGroup))
			}
		}
	} else if hour > 100 {
		minute = math.Mod(hour, 100)
		hour = math.Floor(hour / 100)
	}

	if minute >= 60 {
		return nil
	}

	if hour > 24 {
		return nil
	}
	if hour >= 12 {
		meridiem = 1
	}

	// ----- AM & PM
	if match.Truthy(timeAmPmHourGroup) {
		if hour > 12 {
			return nil
		}
		ampm := js.ToLowerCase(core.FirstChar(match.Get(timeAmPmHourGroup)))
		if ampm == "a" {
			meridiem = 0
			if hour == 12 {
				hour = 0
			}
		}

		if ampm == "p" {
			meridiem = 1
			if hour != 12 {
				hour += 12
			}
		}
	} else if match.Truthy(timeZhAmPmHourGroup1) {
		zhAMPMString1 := match.Get(timeZhAmPmHourGroup1)
		zhAMPM1 := core.FirstChar(zhAMPMString1)
		if zhAMPM1 == "朝" || zhAMPM1 == "早" {
			meridiem = 0
			if hour == 12 {
				hour = 0
			}
		} else if zhAMPM1 == "晚" {
			meridiem = 1
			if hour != 12 {
				hour += 12
			}
		}
	} else if match.Truthy(timeZhAmPmHourGroup2) {
		zhAMPMString2 := match.Get(timeZhAmPmHourGroup2)
		zhAMPM2 := core.FirstChar(zhAMPMString2)
		if zhAMPM2 == "上" || zhAMPM2 == "朝" || zhAMPM2 == "早" || zhAMPM2 == "凌" {
			meridiem = 0
			if hour == 12 {
				hour = 0
			}
		} else if zhAMPM2 == "下" || zhAMPM2 == "晏" || zhAMPM2 == "晚" {
			meridiem = 1
			if hour != 12 {
				hour += 12
			}
		}
	} else if match.Truthy(timeZhAmPmHourGroup3) {
		zhAMPMString3 := match.Get(timeZhAmPmHourGroup3)
		zhAMPM3 := core.FirstChar(zhAMPMString3)
		if zhAMPM3 == "上" || zhAMPM3 == "朝" || zhAMPM3 == "早" || zhAMPM3 == "凌" {
			meridiem = 0
			if hour == 12 {
				hour = 0
			}
		} else if zhAMPM3 == "下" || zhAMPM3 == "晏" || zhAMPM3 == "晚" {
			meridiem = 1
			if hour != 12 {
				hour += 12
			}
		}
	}

	result.Start.Assign(core.ComponentHour, hour)
	result.Start.Assign(core.ComponentMinute, minute)

	if meridiem >= 0 {
		result.Start.Assign(core.ComponentMeridiem, meridiem)
	} else {
		if hour < 12 {
			result.Start.Imply(core.ComponentMeridiem, 0)
		} else {
			result.Start.Imply(core.ComponentMeridiem, 1)
		}
	}

	// ==============================================================
	//                  Extracting the 'to' chunk
	// ==============================================================

	secondMatch := timeSecondRegPattern.Exec(context.SubstringFrom(result.Index + js.Len(result.Text)))
	if secondMatch == nil {
		// Not accept number only result
		if timeDigitsOnlyPattern.Exec(result.Text) != nil {
			return nil
		}
		return result
	}

	endMoment := jsdate.FromTime(startMoment.Time())
	if secondMatch.Truthy(timeDayGroup1) || secondMatch.Truthy(timeDayGroup3) {
		// If the end time expression has a relative day specified (e.g. "Tomorrow"),
		// we should reset the end date to the reference date to calculate the relative date correctly.
		endMoment = jsdate.FromTime(jsdate.FromTime(context.Reference.Instant).Time())
	}
	result.End = context.CreateParsingComponents(nil)

	// ----- Day
	if secondMatch.Truthy(timeDayGroup1) {
		day1 := secondMatch.Get(timeDayGroup1)
		if day1 == "明" || day1 == "聽" {
			// Check not "Tomorrow" on late night
			if jsdate.FromTime(context.Reference.Instant).GetHours() > 1 {
				endMoment.SetDate(endMoment.GetDate() + 1)
			}
		} else if day1 == "昨" || day1 == "尋" || day1 == "琴" {
			endMoment.SetDate(endMoment.GetDate() - 1)
		} else if day1 == "前" {
			endMoment.SetDate(endMoment.GetDate() - 2)
		} else if day1 == "大前" {
			endMoment.SetDate(endMoment.GetDate() - 3)
		} else if day1 == "後" {
			endMoment.SetDate(endMoment.GetDate() + 2)
		} else if day1 == "大後" {
			endMoment.SetDate(endMoment.GetDate() + 3)
		}
		result.End.Assign(core.ComponentDay, endMoment.GetDate())
		result.End.Assign(core.ComponentMonth, endMoment.GetMonth()+1)
		result.End.Assign(core.ComponentYear, endMoment.GetFullYear())
	} else if secondMatch.Truthy(timeDayGroup3) {
		day3 := secondMatch.Get(timeDayGroup3)
		if day3 == "明" || day3 == "聽" {
			endMoment.SetDate(endMoment.GetDate() + 1)
		} else if day3 == "昨" || day3 == "尋" || day3 == "琴" {
			endMoment.SetDate(endMoment.GetDate() - 1)
		} else if day3 == "前" {
			endMoment.SetDate(endMoment.GetDate() - 2)
		} else if day3 == "大前" {
			endMoment.SetDate(endMoment.GetDate() - 3)
		} else if day3 == "後" {
			endMoment.SetDate(endMoment.GetDate() + 2)
		} else if day3 == "大後" {
			endMoment.SetDate(endMoment.GetDate() + 3)
		}
		result.End.Assign(core.ComponentDay, endMoment.GetDate())
		result.End.Assign(core.ComponentMonth, endMoment.GetMonth()+1)
		result.End.Assign(core.ComponentYear, endMoment.GetFullYear())
	} else {
		result.End.Imply(core.ComponentDay, endMoment.GetDate())
		result.End.Imply(core.ComponentMonth, endMoment.GetMonth()+1)
		result.End.Imply(core.ComponentYear, endMoment.GetFullYear())
	}

	hour = 0
	minute = 0
	meridiem = -1

	// ----- Second
	if secondMatch.Truthy(timeSecondGroup) {
		second := js.ParseInt(secondMatch.Get(timeSecondGroup), 0)
		if js.IsNaN(second) {
			second = ZhStringToNumber(secondMatch.Get(timeSecondGroup))
		}

		if second >= 60 {
			return nil
		}
		result.End.Assign(core.ComponentSecond, second)
	}

	hour = js.ParseInt(secondMatch.Get(timeHourGroup), 0)
	if js.IsNaN(hour) {
		hour = ZhStringToNumber(secondMatch.Get(timeHourGroup))
	}

	// ----- Minutes
	if secondMatch.Truthy(timeMinuteGroup) {
		if secondMatch.Get(timeMinuteGroup) == "半" {
			minute = 30
		} else if secondMatch.Get(timeMinuteGroup) == "正" || secondMatch.Get(timeMinuteGroup) == "整" {
			minute = 0
		} else {
			minute = js.ParseInt(secondMatch.Get(timeMinuteGroup), 0)
			if js.IsNaN(minute) {
				minute = ZhStringToNumber(secondMatch.Get(timeMinuteGroup))
			}
		}
	} else if hour > 100 {
		minute = math.Mod(hour, 100)
		hour = math.Floor(hour / 100)
	}

	if minute >= 60 {
		return nil
	}

	if hour > 24 {
		return nil
	}
	if hour >= 12 {
		meridiem = 1
	}

	// ----- AM & PM
	if secondMatch.Truthy(timeAmPmHourGroup) {
		if hour > 12 {
			return nil
		}
		ampm := js.ToLowerCase(core.FirstChar(secondMatch.Get(timeAmPmHourGroup)))
		if ampm == "a" {
			meridiem = 0
			if hour == 12 {
				hour = 0
			}
		}

		if ampm == "p" {
			meridiem = 1
			if hour != 12 {
				hour += 12
			}
		}

		if !result.Start.IsCertain(core.ComponentMeridiem) {
			if meridiem == 0 {
				result.Start.Imply(core.ComponentMeridiem, 0)

				if result.Start.V(core.ComponentHour) == 12 {
					result.Start.Assign(core.ComponentHour, 0)
				}
			} else {
				result.Start.Imply(core.ComponentMeridiem, 1)

				if result.Start.V(core.ComponentHour) != 12 {
					result.Start.Assign(core.ComponentHour, result.Start.V(core.ComponentHour)+12)
				}
			}
		}
	} else if secondMatch.Truthy(timeZhAmPmHourGroup1) {
		zhAMPMString1 := secondMatch.Get(timeZhAmPmHourGroup1)
		zhAMPM1 := core.FirstChar(zhAMPMString1)
		if zhAMPM1 == "朝" || zhAMPM1 == "早" {
			meridiem = 0
			if hour == 12 {
				hour = 0
			}
		} else if zhAMPM1 == "晚" {
			meridiem = 1
			if hour != 12 {
				hour += 12
			}
		}
	} else if secondMatch.Truthy(timeZhAmPmHourGroup2) {
		zhAMPMString2 := secondMatch.Get(timeZhAmPmHourGroup2)
		zhAMPM2 := core.FirstChar(zhAMPMString2)
		if zhAMPM2 == "上" || zhAMPM2 == "朝" || zhAMPM2 == "早" || zhAMPM2 == "凌" {
			meridiem = 0
			if hour == 12 {
				hour = 0
			}
		} else if zhAMPM2 == "下" || zhAMPM2 == "晏" || zhAMPM2 == "晚" {
			meridiem = 1
			if hour != 12 {
				hour += 12
			}
		}
	} else if secondMatch.Truthy(timeZhAmPmHourGroup3) {
		zhAMPMString3 := secondMatch.Get(timeZhAmPmHourGroup3)
		zhAMPM3 := core.FirstChar(zhAMPMString3)
		if zhAMPM3 == "上" || zhAMPM3 == "朝" || zhAMPM3 == "早" || zhAMPM3 == "凌" {
			meridiem = 0
			if hour == 12 {
				hour = 0
			}
		} else if zhAMPM3 == "下" || zhAMPM3 == "晏" || zhAMPM3 == "晚" {
			meridiem = 1
			if hour != 12 {
				hour += 12
			}
		}
	}

	result.Text = result.Text + secondMatch.Get(0)
	result.End.Assign(core.ComponentHour, hour)
	result.End.Assign(core.ComponentMinute, minute)
	if meridiem >= 0 {
		result.End.Assign(core.ComponentMeridiem, meridiem)
	} else {
		startAtPM := result.Start.IsCertain(core.ComponentMeridiem) && result.Start.V(core.ComponentMeridiem) == 1
		if startAtPM && result.Start.V(core.ComponentHour) > hour {
			// 10pm - 1 (am)
			result.End.Imply(core.ComponentMeridiem, 0)
		} else if hour > 12 {
			result.End.Imply(core.ComponentMeridiem, 1)
		}
	}

	if result.End.Date() < result.Start.Date() {
		result.End.Imply(core.ComponentDay, result.End.V(core.ComponentDay)+1)
	}

	return result
}
