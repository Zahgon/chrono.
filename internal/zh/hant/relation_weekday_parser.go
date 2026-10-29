// Ported from src/locales/zh/hant/parsers/ZHHantRelationWeekdayParser.ts.

package hant

import (
	"math"
	"strings"

	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/jsdate"
)

var relationWeekdayPattern = core.MustRegExp(
	"(?<prefix>上|今|下|這|呢)(?:個)?(?:星期|禮拜|週)(?<weekday>"+strings.Join(WeekdayOffset.Keys(), "|")+")",
	"",
)

// RelationWeekdayParser is ZHHantRelationWeekdayParser ("上星期一", "下個禮拜五", "呢個星期日", ...).
type RelationWeekdayParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewRelationWeekdayParser returns a new ZHHantRelationWeekdayParser.
func NewRelationWeekdayParser() *RelationWeekdayParser {
	p := &RelationWeekdayParser{}
	p.Self = p
	return p
}

// InnerPattern returns the relative weekday pattern.
func (p *RelationWeekdayParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return relationWeekdayPattern
}

// InnerExtract returns the result at the weekday.
func (p *RelationWeekdayParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	result := context.CreateParsingResult(match.Index, match.Get(0), nil, nil)

	dayOfWeek, _ := match.Group("weekday")
	offset, ok := WeekdayOffset.Get(dayOfWeek)
	if !ok {
		return nil
	}

	modifier := ""
	prefix, _ := match.Group("prefix")

	if prefix == "上" {
		modifier = "last"
	} else if prefix == "下" {
		modifier = "next"
	} else if prefix == "今" || prefix == "這" || prefix == "呢" {
		modifier = "this"
	}

	date := jsdate.FromTime(jsdate.FromTime(context.RefDate).Time())
	startMomentFixed := false
	refOffset := date.GetDay()

	if modifier == "last" || modifier == "past" {
		date.SetDate(date.GetDate() + (offset - 7 - refOffset))
		startMomentFixed = true
	} else if modifier == "next" {
		date.SetDate(date.GetDate() + (offset + 7 - refOffset))
		startMomentFixed = true
	} else if modifier == "this" {
		date.SetDate(date.GetDate() + (offset - refOffset))
	} else {
		diff := offset - refOffset
		if math.Abs(diff-7) < math.Abs(diff) {
			diff -= 7
		}
		if math.Abs(diff+7) < math.Abs(diff) {
			diff += 7
		}
		date.SetDate(date.GetDate() + diff)
	}

	result.Start.Assign(core.ComponentWeekday, offset)
	if startMomentFixed {
		result.Start.Assign(core.ComponentDay, date.GetDate())
		result.Start.Assign(core.ComponentMonth, date.GetMonth()+1)
		result.Start.Assign(core.ComponentYear, date.GetFullYear())
	} else {
		result.Start.Imply(core.ComponentDay, date.GetDate())
		result.Start.Imply(core.ComponentMonth, date.GetMonth()+1)
		result.Start.Imply(core.ComponentYear, date.GetFullYear())
	}

	return result
}
