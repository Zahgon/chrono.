// Ported from src/locales/zh/hant/parsers/ZHHantDeadlineFormatParser.ts.

package hant

import (
	"strings"

	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
	"github.com/wanasit/chrono-go/internal/jsdate"
)

var deadlineFormatPattern = core.MustRegExp(
	"(\\d+|["+
		strings.Join(Number.Keys(), "")+
		"]+|半|幾)(?:\\s*)"+
		"(?:個)?"+
		"(秒(?:鐘)?|分鐘|小時|鐘|日|天|星期|禮拜|月|年)"+
		"(?:(?:之|過)?後|(?:之)?內)",
	"i",
)

var deadlineDateUnitPattern = core.MustRegExp("[日天星禮月年]", "")

const (
	deadlineNumberGroup = 1
	deadlineUnitGroup   = 2
)

// DeadlineFormatParser is ZHHantDeadlineFormatParser ("5天內", "三個月後", "半小時之後", ...).
type DeadlineFormatParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewDeadlineFormatParser returns a new ZHHantDeadlineFormatParser.
func NewDeadlineFormatParser() *DeadlineFormatParser {
	p := &DeadlineFormatParser{}
	p.Self = p
	return p
}

// InnerPattern returns the deadline pattern.
func (p *DeadlineFormatParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return deadlineFormatPattern
}

// InnerExtract returns the result at the end of the duration.
func (p *DeadlineFormatParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	result := context.CreateParsingResult(match.Index, match.Get(0), nil, nil)

	number := js.ParseInt(match.Get(deadlineNumberGroup), 0)
	if js.IsNaN(number) {
		number = ZhStringToNumber(match.Get(deadlineNumberGroup))
	}

	if js.IsNaN(number) {
		str := match.Get(deadlineNumberGroup)
		if str == "幾" {
			number = 3
		} else if str == "半" {
			number = 0.5
		} else {
			//just in case
			return nil
		}
	}

	duration := core.NewDuration()
	unit := match.Get(deadlineUnitGroup)
	unitAbbr := core.FirstChar(unit)

	if deadlineDateUnitPattern.Exec(unitAbbr) != nil {
		if unitAbbr == "日" || unitAbbr == "天" {
			duration.Set("day", number)
		} else if unitAbbr == "星" || unitAbbr == "禮" {
			duration.Set("week", number)
		} else if unitAbbr == "月" {
			duration.Set("month", number)
		} else if unitAbbr == "年" {
			duration.Set("year", number)
		}

		date := core.AddDuration(jsdate.FromTime(context.RefDate), duration)
		result.Start.Assign(core.ComponentYear, date.GetFullYear())
		result.Start.Assign(core.ComponentMonth, date.GetMonth()+1)
		result.Start.Assign(core.ComponentDay, date.GetDate())
		return result
	}

	if unitAbbr == "秒" {
		duration.Set("second", number)
	} else if unitAbbr == "分" {
		duration.Set("minute", number)
	} else if unitAbbr == "小" || unitAbbr == "鐘" {
		duration.Set("hour", number)
	}

	date := core.AddDuration(jsdate.FromTime(context.RefDate), duration)
	result.Start.Imply(core.ComponentYear, date.GetFullYear())
	result.Start.Imply(core.ComponentMonth, date.GetMonth()+1)
	result.Start.Imply(core.ComponentDay, date.GetDate())
	result.Start.Assign(core.ComponentHour, date.GetHours())
	result.Start.Assign(core.ComponentMinute, date.GetMinutes())
	result.Start.Assign(core.ComponentSecond, date.GetSeconds())
	return result
}
