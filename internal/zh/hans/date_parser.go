// Ported from src/locales/zh/hans/parsers/ZHHansDateParser.ts.

package hans

import (
	"strings"

	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
	"github.com/wanasit/chrono-go/internal/jsdate"
)

const (
	dateYearGroup  = 1
	dateMonthGroup = 2
	dateDayGroup   = 3
)

// DateParser is ZHHansDateParser ("2016年9月3号", "二〇一六年九月三日", "9月3日", ...).
type DateParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewDateParser returns a new ZHHansDateParser.
func NewDateParser() *DateParser {
	p := &DateParser{}
	p.Self = p
	return p
}

// InnerPattern returns the date pattern (a new RegExp per call, as in the original).
func (p *DateParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return core.MustRegExp(
		"("+
			"\\d{2,4}|"+
			"["+
			strings.Join(Number.Keys(), "")+
			"]{4}|"+
			"["+
			strings.Join(Number.Keys(), "")+
			"]{2}"+
			")?"+
			"(?:\\s*)"+
			"(?:年)?"+
			"(?:[\\s|,|，]*)"+
			"("+
			"\\d{1,2}|"+
			"["+
			strings.Join(Number.Keys(), "")+
			"]{1,3}"+
			")"+
			"(?:\\s*)"+
			"(?:月)"+
			"(?:\\s*)"+
			"("+
			"\\d{1,2}|"+
			"["+
			strings.Join(Number.Keys(), "")+
			"]{1,3}"+
			")?"+
			"(?:\\s*)"+
			"(?:日|号)?",
		"",
	)
}

// InnerExtract returns the result of the date.
func (p *DateParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	result := context.CreateParsingResult(match.Index, match.Get(0), nil, nil)

	//Month
	month := js.ParseInt(match.Get(dateMonthGroup), 0)
	if js.IsNaN(month) {
		month = ZhStringToNumber(match.Get(dateMonthGroup))
	}
	result.Start.Assign(core.ComponentMonth, month)

	//Day
	if match.Truthy(dateDayGroup) {
		day := js.ParseInt(match.Get(dateDayGroup), 0)
		if js.IsNaN(day) {
			day = ZhStringToNumber(match.Get(dateDayGroup))
		}
		result.Start.Assign(core.ComponentDay, day)
	} else {
		result.Start.Imply(core.ComponentDay, jsdate.FromTime(context.RefDate).GetDate())
	}

	//Year
	if match.Truthy(dateYearGroup) {
		year := js.ParseInt(match.Get(dateYearGroup), 0)
		if js.IsNaN(year) {
			year = ZhStringToYear(match.Get(dateYearGroup))
		}
		result.Start.Assign(core.ComponentYear, year)
	} else {
		result.Start.Imply(core.ComponentYear, jsdate.FromTime(context.RefDate).GetFullYear())
	}

	return result
}
