// Ported from src/locales/ja/constants.ts.

package ja

import (
	"math"

	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
	"github.com/wanasit/chrono-go/internal/jsregex"
)

// Number is NUMBER: the Japanese numerals.
var Number = core.NewDictionary([]core.Entry[float64]{
	{Key: "零", Value: 0},
	{Key: "〇", Value: 0},
	{Key: "一", Value: 1},
	{Key: "二", Value: 2},
	{Key: "三", Value: 3},
	{Key: "四", Value: 4},
	{Key: "五", Value: 5},
	{Key: "六", Value: 6},
	{Key: "七", Value: 7},
	{Key: "八", Value: 8},
	{Key: "九", Value: 9},
	{Key: "十", Value: 10},
}...)

// WeekdayOffset is WEEKDAY_OFFSET: the weekday kanji.
var WeekdayOffset = core.NewDictionary([]core.Entry[float64]{
	{Key: "日", Value: 0},
	{Key: "月", Value: 1},
	{Key: "火", Value: 2},
	{Key: "水", Value: 3},
	{Key: "木", Value: 4},
	{Key: "金", Value: 5},
	{Key: "土", Value: 6},
}...)

var (
	hankakuRightSingleQuote = core.MustRegExp("\\u2019", "g")
	hankakuRightDoubleQuote = core.MustRegExp("\\u201D", "g")
	hankakuIdeographicSpace = core.MustRegExp("\\u3000", "g")
	hankakuFullwidthYen     = core.MustRegExp("\\uFFE5", "g")
	hankakuAlphaNum         = core.MustRegExp("[\\uFF01\\uFF03-\\uFF06\\uFF08\\uFF09\\uFF0C-\\uFF19\\uFF1C-\\uFF1F\\uFF21-\\uFF3B\\uFF3D\\uFF3F\\uFF41-\\uFF5B\\uFF5D\\uFF5E]", "g")
)

// ToHankaku is toHankaku (to-hankaku.js 1.0.1 by think49, MIT): converts full-width ASCII
// characters (and a few quotes and spaces) to their ASCII forms.
func ToHankaku(text string) string {
	s := hankakuRightSingleQuote.Replace(text, "'")
	s = hankakuRightDoubleQuote.Replace(s, `"`)
	s = hankakuIdeographicSpace.Replace(s, " ")
	s = hankakuFullwidthYen.Replace(s, "¥")
	return hankakuAlphaNum.Regexp().ReplaceFunc(s, func(m *jsregex.Match) string {
		return alphaNum(m.Groups[0])
	})
}

func alphaNum(token string) string {
	return js.FromCharCode(js.CharCodeAt(token, 0) - 65248)
}

// numberValue is NUMBER[char] in arithmetic (undefined is NaN).
func numberValue(char string) float64 {
	if v, ok := Number.Get(char); ok {
		return v
	}
	return math.NaN()
}

// JaStringToNumber is jaStringToNumber: the value of a numeral written in kanji ("二十五" = 25).
func JaStringToNumber(text string) float64 {
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
