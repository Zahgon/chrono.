// Ported from src/locales/vi/parsers/VITimeExpressionParser.ts.

package vi

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

// lúc 7 giờ 30 phút  |  7 giờ sáng  |  7 giờ  |  15:30
var timeExpressionPattern = core.MustRegExp(
	"(?:l\u00fac\\s*|v\u00e0o\\s*)?"+
		"([0-9]{1,2})"+
		"(?:\\s*gi\u1edd\\s*([0-9]{1,2})?\\s*(?:ph\u00fat\\s*)?"+
		"(s\u00e1ng|tr\u01b0a|chi\u1ec1u|t\u1ed1i|\u0111\u00eam)?"+
		"|:([0-9]{2}))"+
		"(?=\\W|$)",
	"i",
)

const (
	timeHourGroup        = 1
	timeMinuteGioGroup   = 2
	timeMeridiemGroup    = 3
	timeMinuteColonGroup = 4
)

// TimeExpressionParser is VITimeExpressionParser ("lúc 7 giờ 30 phút", "7 giờ sáng", "15:30").
type TimeExpressionParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewTimeExpressionParser returns a new VITimeExpressionParser.
func NewTimeExpressionParser() *TimeExpressionParser {
	p := &TimeExpressionParser{}
	p.Self = p
	return p
}

// InnerPattern returns the time pattern.
func (p *TimeExpressionParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return timeExpressionPattern
}

// InnerExtract returns the result of the time.
func (p *TimeExpressionParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	hour := js.ParseInt(match.Get(timeHourGroup), 0)
	if hour > 23 {
		return nil
	}

	result := context.CreateParsingResult(match.Index, match.Get(0), nil, nil)
	result.Start.Assign(core.ComponentHour, hour)

	var minute float64
	if match.Truthy(timeMinuteColonGroup) {
		minute = js.ParseInt(match.Get(timeMinuteColonGroup), 0)
	} else if match.Truthy(timeMinuteGioGroup) {
		minute = js.ParseInt(match.Get(timeMinuteGioGroup), 0)
	} else {
		minute = 0
	}
	if minute >= 60 {
		return nil
	}
	result.Start.Assign(core.ComponentMinute, minute)

	// match[MERIDIEM_GROUP]?.toLowerCase(): undefined never equals the strings below.
	meridiem := ""
	hasMeridiem := match.Has(timeMeridiemGroup)
	if hasMeridiem {
		meridiem = js.ToLowerCase(match.Get(timeMeridiemGroup))
	}
	if hasMeridiem && meridiem == "sáng" {
		// "12 giờ sáng" = midnight (00:00), matching EN convention
		result.Start.Assign(core.ComponentMeridiem, core.AM)
		if hour == 12 {
			result.Start.Assign(core.ComponentHour, 0)
		}
	} else if hasMeridiem && meridiem == "trưa" {
		// "trưa" = noon/midday (~11 AM - 1 PM)
		// "1 giờ trưa" = 13:00, but "11 giờ trưa" = 11:00 (approaching noon)
		if hour < 10 {
			result.Start.Assign(core.ComponentMeridiem, core.PM)
			result.Start.Assign(core.ComponentHour, hour+12)
		} else {
			// 10-12: keep as-is (AM for 10-11, PM for 12)
			if hour >= 12 {
				result.Start.Assign(core.ComponentMeridiem, core.PM)
			} else {
				result.Start.Assign(core.ComponentMeridiem, core.AM)
			}
		}
	} else if hasMeridiem && (meridiem == "chiều" || meridiem == "tối" || meridiem == "đêm") {
		result.Start.Assign(core.ComponentMeridiem, core.PM)
		if hour < 12 {
			result.Start.Assign(core.ComponentHour, hour+12)
		}
	}

	result.Start.Imply(core.ComponentSecond, 0)
	result.Start.Imply(core.ComponentMillisecond, 0)
	return result
}
