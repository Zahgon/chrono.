// Ported from src/locales/zh/hant/parsers/ZHHantAgoFormatParser.ts.

package hant

import (
	"strings"

	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
	"github.com/wanasit/chrono-go/internal/jsdate"
)

var agoFormatPattern = core.MustRegExp(
	"(\\d+|["+
		strings.Join(Number.Keys(), "")+
		"]+|半|幾)(?:\\s*)"+
		"(?:個)?"+
		"(秒(?:鐘)?|分鐘|小時|鐘|日|天|星期|禮拜|月|年)"+
		"(?:之)?前",
	"i",
)

var agoDateUnitPattern = core.MustRegExp("[日天星禮月年]", "")

const (
	agoNumberGroup = 1
	agoUnitGroup   = 2
)

// AgoFormatParser is ZHHantAgoFormatParser ("5天前", "三個月之前", "半小時前", ...).
type AgoFormatParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewAgoFormatParser returns a new ZHHantAgoFormatParser.
func NewAgoFormatParser() *AgoFormatParser {
	p := &AgoFormatParser{}
	p.Self = p
	return p
}

// InnerPattern returns the "ago" pattern.
func (p *AgoFormatParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return agoFormatPattern
}

// InnerExtract returns the result at the duration before the reference.
func (p *AgoFormatParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	result := context.CreateParsingResult(match.Index, match.Get(0), nil, nil)

	number := js.ParseInt(match.Get(agoNumberGroup), 0)
	if js.IsNaN(number) {
		number = ZhStringToNumber(match.Get(agoNumberGroup))
	}

	if js.IsNaN(number) {
		str := match.Get(agoNumberGroup)
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
	unit := match.Get(agoUnitGroup)
	unitAbbr := core.FirstChar(unit)

	if agoDateUnitPattern.Exec(unitAbbr) != nil {
		if unitAbbr == "日" || unitAbbr == "天" {
			duration.Set("day", number)
		} else if unitAbbr == "星" || unitAbbr == "禮" {
			duration.Set("week", number)
		} else if unitAbbr == "月" {
			duration.Set("month", number)
		} else if unitAbbr == "年" {
			duration.Set("year", number)
		}

		duration = core.ReverseDuration(duration)

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

	duration = core.ReverseDuration(duration)

	date := core.AddDuration(jsdate.FromTime(context.RefDate), duration)
	result.Start.Imply(core.ComponentYear, date.GetFullYear())
	result.Start.Imply(core.ComponentMonth, date.GetMonth()+1)
	result.Start.Imply(core.ComponentDay, date.GetDate())
	result.Start.Assign(core.ComponentHour, date.GetHours())
	result.Start.Assign(core.ComponentMinute, date.GetMinutes())
	result.Start.Assign(core.ComponentSecond, date.GetSeconds())
	return result
}
