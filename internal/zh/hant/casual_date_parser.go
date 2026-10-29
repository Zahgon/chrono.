// Ported from src/locales/zh/hant/parsers/ZHHantCasualDateParser.ts.

package hant

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/jsdate"
)

const (
	casualNowGroup   = 1
	casualDayGroup1  = 2
	casualTimeGroup1 = 3
	casualTimeGroup2 = 4
	casualDayGroup3  = 5
	casualTimeGroup3 = 6
)

// CasualDateParser is ZHHantCasualDateParser ("而家", "今晚", "聽日朝早", "下晝", ...).
type CasualDateParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewCasualDateParser returns a new ZHHantCasualDateParser.
func NewCasualDateParser() *CasualDateParser {
	p := &CasualDateParser{}
	p.Self = p
	return p
}

// InnerPattern returns the casual date pattern (a new RegExp per call, as in the original).
func (p *CasualDateParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return core.MustRegExp(
		"(而家|立(?:刻|即)|即刻)|"+
			"(今|明|前|大前|後|大後|聽|昨|尋|琴)(早|朝|晚)|"+
			"(上(?:午|晝)|朝(?:早)|早(?:上)|下(?:午|晝)|晏(?:晝)|晚(?:上)|夜(?:晚)?|中(?:午)|凌(?:晨))|"+
			"(今|明|前|大前|後|大後|聽|昨|尋|琴)(?:日|天)"+
			"(?:[\\s|,|，]*)"+
			"(?:(上(?:午|晝)|朝(?:早)|早(?:上)|下(?:午|晝)|晏(?:晝)|晚(?:上)|夜(?:晚)?|中(?:午)|凌(?:晨)))?",
		"i",
	)
}

// InnerExtract returns the result of the casual date.
func (p *CasualDateParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	index := match.Index
	result := context.CreateParsingResult(index, match.Get(0), nil, nil)

	refDate := jsdate.FromTime(context.RefDate)
	date := jsdate.FromTime(refDate.Time())

	if match.Truthy(casualNowGroup) {
		result.Start.Imply(core.ComponentHour, refDate.GetHours())
		result.Start.Imply(core.ComponentMinute, refDate.GetMinutes())
		result.Start.Imply(core.ComponentSecond, refDate.GetSeconds())
		result.Start.Imply(core.ComponentMillisecond, refDate.GetMilliseconds())
	} else if match.Truthy(casualDayGroup1) {
		day1 := match.Get(casualDayGroup1)
		time1 := match.Get(casualTimeGroup1)

		if day1 == "明" || day1 == "聽" {
			// Check not "Tomorrow" on late night
			if refDate.GetHours() > 1 {
				date.SetDate(date.GetDate() + 1)
			}
		} else if day1 == "昨" || day1 == "尋" || day1 == "琴" {
			date.SetDate(date.GetDate() - 1)
		} else if day1 == "前" {
			date.SetDate(date.GetDate() - 2)
		} else if day1 == "大前" {
			date.SetDate(date.GetDate() - 3)
		} else if day1 == "後" {
			date.SetDate(date.GetDate() + 2)
		} else if day1 == "大後" {
			date.SetDate(date.GetDate() + 3)
		}

		if time1 == "早" || time1 == "朝" {
			result.Start.Imply(core.ComponentHour, 6)
		} else if time1 == "晚" {
			result.Start.Imply(core.ComponentHour, 22)
			result.Start.Imply(core.ComponentMeridiem, 1)
		}
	} else if match.Truthy(casualTimeGroup2) {
		timeString2 := match.Get(casualTimeGroup2)
		time2 := core.FirstChar(timeString2)
		if time2 == "早" || time2 == "朝" || time2 == "上" {
			result.Start.Imply(core.ComponentHour, 6)
		} else if time2 == "下" || time2 == "晏" {
			result.Start.Imply(core.ComponentHour, 15)
			result.Start.Imply(core.ComponentMeridiem, 1)
		} else if time2 == "中" {
			result.Start.Imply(core.ComponentHour, 12)
			result.Start.Imply(core.ComponentMeridiem, 1)
		} else if time2 == "夜" || time2 == "晚" {
			result.Start.Imply(core.ComponentHour, 22)
			result.Start.Imply(core.ComponentMeridiem, 1)
		} else if time2 == "凌" {
			result.Start.Imply(core.ComponentHour, 0)
		}
	} else if match.Truthy(casualDayGroup3) {
		day3 := match.Get(casualDayGroup3)

		if day3 == "明" || day3 == "聽" {
			// Check not "Tomorrow" on late night
			if refDate.GetHours() > 1 {
				date.SetDate(date.GetDate() + 1)
			}
		} else if day3 == "昨" || day3 == "尋" || day3 == "琴" {
			date.SetDate(date.GetDate() - 1)
		} else if day3 == "前" {
			date.SetDate(date.GetDate() - 2)
		} else if day3 == "大前" {
			date.SetDate(date.GetDate() - 3)
		} else if day3 == "後" {
			date.SetDate(date.GetDate() + 2)
		} else if day3 == "大後" {
			date.SetDate(date.GetDate() + 3)
		}

		if match.Truthy(casualTimeGroup3) {
			timeString3 := match.Get(casualTimeGroup3)
			time3 := core.FirstChar(timeString3)
			if time3 == "早" || time3 == "朝" || time3 == "上" {
				result.Start.Imply(core.ComponentHour, 6)
			} else if time3 == "下" || time3 == "晏" {
				result.Start.Imply(core.ComponentHour, 15)
				result.Start.Imply(core.ComponentMeridiem, 1)
			} else if time3 == "中" {
				result.Start.Imply(core.ComponentHour, 12)
				result.Start.Imply(core.ComponentMeridiem, 1)
			} else if time3 == "夜" || time3 == "晚" {
				result.Start.Imply(core.ComponentHour, 22)
				result.Start.Imply(core.ComponentMeridiem, 1)
			} else if time3 == "凌" {
				result.Start.Imply(core.ComponentHour, 0)
			}
		}
	}

	result.Start.Assign(core.ComponentDay, date.GetDate())
	result.Start.Assign(core.ComponentMonth, date.GetMonth()+1)
	result.Start.Assign(core.ComponentYear, date.GetFullYear())

	return result
}
