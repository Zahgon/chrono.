// Ported from src/locales/zh/hant/parsers/ZHHantWeekdayParser.ts.

package hant

import (
	"math"
	"strings"

	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/jsdate"
)

var weekdayPattern = core.MustRegExp("(?:星期|禮拜|週)(?<weekday>"+strings.Join(WeekdayOffset.Keys(), "|")+")", "")

// WeekdayParser is ZHHantWeekdayParser ("星期一", "禮拜天", "週五", ...).
type WeekdayParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewWeekdayParser returns a new ZHHantWeekdayParser.
func NewWeekdayParser() *WeekdayParser {
	p := &WeekdayParser{}
	p.Self = p
	return p
}

// InnerPattern returns the weekday pattern.
func (p *WeekdayParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return weekdayPattern
}

// InnerExtract returns the result at the closest such weekday.
func (p *WeekdayParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	result := context.CreateParsingResult(match.Index, match.Get(0), nil, nil)

	dayOfWeek, _ := match.Group("weekday")
	offset, ok := WeekdayOffset.Get(dayOfWeek)
	if !ok {
		return nil
	}

	date := jsdate.FromTime(jsdate.FromTime(context.RefDate).Time())
	const startMomentFixed = false
	refOffset := date.GetDay()

	diff := offset - refOffset
	if math.Abs(diff-7) < math.Abs(diff) {
		diff -= 7
	}
	if math.Abs(diff+7) < math.Abs(diff) {
		diff += 7
	}

	date.SetDate(date.GetDate() + diff)

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
