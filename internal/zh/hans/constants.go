// Ported from src/locales/zh/hans/constants.ts.

package hans

import (
	"math"

	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

// Number is NUMBER: the Chinese (simplified) numerals.
var Number = core.NewDictionary([]core.Entry[float64]{
	{Key: "零", Value: 0},
	{Key: "〇", Value: 0},
	{Key: "一", Value: 1},
	{Key: "二", Value: 2},
	{Key: "两", Value: 2},
	{Key: "三", Value: 3},
	{Key: "四", Value: 4},
	{Key: "五", Value: 5},
	{Key: "六", Value: 6},
	{Key: "七", Value: 7},
	{Key: "八", Value: 8},
	{Key: "九", Value: 9},
	{Key: "十", Value: 10},
}...)

// WeekdayOffset is WEEKDAY_OFFSET: the weekday characters (after 星期/礼拜/周).
var WeekdayOffset = core.NewDictionary([]core.Entry[float64]{
	{Key: "天", Value: 0},
	{Key: "日", Value: 0},
	{Key: "一", Value: 1},
	{Key: "二", Value: 2},
	{Key: "三", Value: 3},
	{Key: "四", Value: 4},
	{Key: "五", Value: 5},
	{Key: "六", Value: 6},
}...)

// numberValue is NUMBER[char] in arithmetic (undefined is NaN).
func numberValue(char string) float64 {
	if v, ok := Number.Get(char); ok {
		return v
	}
	return math.NaN()
}

// ZhStringToNumber is zhStringToNumber: the value of a Chinese numeral ("二十五" = 25).
func ZhStringToNumber(text string) float64 {
	number := 0.0

	for i := 0; i < js.Len(text); i++ {
		char := js.CharAt(text, i)
		if char == "十" {
			if number == 0 {
				number = numberValue(char)
			} else {
				number = number * numberValue(char)
			}
		} else {
			number += numberValue(char)
		}
	}

	return number
}

// ZhStringToYear is zhStringToYear: a year written digit by digit ("二〇一二" = 2012).
func ZhStringToYear(text string) float64 {
	str := ""

	for i := 0; i < js.Len(text); i++ {
		char := js.CharAt(text, i)
		if v, ok := Number.Get(char); ok {
			str = str + js.NumberToString(v)
		} else {
			str = str + "undefined"
		}
	}

	return js.ParseInt(str, 0)
}
