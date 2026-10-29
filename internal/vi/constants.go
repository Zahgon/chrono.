// Ported from src/locales/vi/constants.ts.

package vi

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

// WeekdayDictionary is WEEKDAY_DICTIONARY.
var WeekdayDictionary = core.NewDictionary([]core.Entry[float64]{
	{Key: "chủ nhật", Value: 0},
	{Key: "cn", Value: 0},
	{Key: "thứ hai", Value: 1},
	{Key: "t2", Value: 1},
	{Key: "thứ ba", Value: 2},
	{Key: "t3", Value: 2},
	{Key: "thứ tư", Value: 3},
	{Key: "t4", Value: 3},
	{Key: "thứ năm", Value: 4},
	{Key: "t5", Value: 4},
	{Key: "thứ sáu", Value: 5},
	{Key: "t6", Value: 5},
	{Key: "thứ bảy", Value: 6},
	{Key: "t7", Value: 6},
}...)

// MonthDictionary is MONTH_DICTIONARY.
var MonthDictionary = core.NewDictionary([]core.Entry[float64]{
	{Key: "tháng 1", Value: 1},
	{Key: "tháng một", Value: 1},
	{Key: "tháng giêng", Value: 1},
	{Key: "tháng 2", Value: 2},
	{Key: "tháng hai", Value: 2},
	{Key: "tháng 3", Value: 3},
	{Key: "tháng ba", Value: 3},
	{Key: "tháng 4", Value: 4},
	{Key: "tháng tư", Value: 4},
	{Key: "tháng 5", Value: 5},
	{Key: "tháng năm", Value: 5},
	{Key: "tháng 6", Value: 6},
	{Key: "tháng sáu", Value: 6},
	{Key: "tháng 7", Value: 7},
	{Key: "tháng bảy", Value: 7},
	{Key: "tháng 8", Value: 8},
	{Key: "tháng tám", Value: 8},
	{Key: "tháng 9", Value: 9},
	{Key: "tháng chín", Value: 9},
	{Key: "tháng 10", Value: 10},
	{Key: "tháng mười", Value: 10},
	{Key: "tháng 11", Value: 11},
	{Key: "tháng mười một", Value: 11},
	{Key: "tháng 12", Value: 12},
	{Key: "tháng mười hai", Value: 12},
	{Key: "tháng chạp", Value: 12},
}...)

// IntegerWordDictionary is INTEGER_WORD_DICTIONARY.
var IntegerWordDictionary = core.NewDictionary([]core.Entry[float64]{
	{Key: "một", Value: 1},
	{Key: "hai", Value: 2},
	{Key: "ba", Value: 3},
	{Key: "bốn", Value: 4},
	{Key: "năm", Value: 5},
	{Key: "sáu", Value: 6},
	{Key: "bảy", Value: 7},
	{Key: "tám", Value: 8},
	{Key: "chín", Value: 9},
	{Key: "mười", Value: 10},
	{Key: "mười một", Value: 11},
	{Key: "mười hai", Value: 12},
}...)

// TimeUnitDictionary is TIME_UNIT_DICTIONARY.
var TimeUnitDictionary = core.NewDictionary([]core.Entry[string]{
	{Key: "giây", Value: "second"},
	{Key: "phút", Value: "minute"},
	{Key: "giờ", Value: "hour"},
	{Key: "ngày", Value: "day"},
	{Key: "tuần", Value: "week"},
	{Key: "tháng", Value: "month"},
	{Key: "năm", Value: "year"},
}...)

// NumberPattern is NUMBER_PATTERN.
var NumberPattern = "(?:" + core.MatchAnyPattern(IntegerWordDictionary.Keys()) + "|[0-9]+|[0-9]+\\.[0-9]+)"

// ParseNumberPattern is parseNumberPattern: a NumberPattern match as a number.
func ParseNumberPattern(match string) float64 {
	num := js.ToLowerCase(match)
	if v, ok := IntegerWordDictionary.Get(num); ok {
		return v
	}
	return js.ParseFloat(num)
}

// YearPattern is YEAR_PATTERN: YYYY, YYYY TCN (Trước Công nguyên = BC).
const YearPattern = "(?:[0-9]{1,4}(?:\\s*TCN)?)"

var (
	reNonDigits = core.MustRegExp("[^0-9]+", "g")
	reTCN       = core.MustRegExp("TCN", "")
)

// ParseYear is parseYear: a YearPattern match as a year (negative for TCN).
func ParseYear(match string) float64 {
	upper := js.ToUpperCase(match)
	num := js.ParseInt(reNonDigits.Replace(match, ""), 0)
	if reTCN.Test(upper) {
		return -num
	}
	return core.FindMostLikelyADYear(num)
}

var singleTimeUnitPattern = "(" + NumberPattern + ")\\s{0,5}(" + core.MatchAnyPattern(TimeUnitDictionary.Keys()) + ")\\s{0,5}"
var singleTimeUnitRegex = core.MustRegExp(singleTimeUnitPattern, "i")

// TimeUnitsPattern is TIME_UNITS_PATTERN.
var TimeUnitsPattern = core.RepeatedTimeunitPattern("", singleTimeUnitPattern, core.DefaultTimeunitConnectorPattern)

// ParseDuration is parseDuration: the duration of a TimeUnitsPattern match (never nil, possibly
// empty).
func ParseDuration(timeunitText string) *core.Duration {
	fragments := core.NewDuration()
	remainingText := timeunitText
	match := singleTimeUnitRegex.Exec(remainingText)
	for match != nil {
		num := ParseNumberPattern(match.Get(1))
		unit, ok := TimeUnitDictionary.Get(js.ToLowerCase(match.Get(2)))
		if !ok {
			// fragments[undefined] = num
			unit = "undefined"
		}
		fragments.Set(unit, num)
		remainingText = js.SubstringFrom(remainingText, js.Len(match.Get(0)))
		match = singleTimeUnitRegex.Exec(remainingText)
	}
	return fragments
}
