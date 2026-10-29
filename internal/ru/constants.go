// Ported from src/locales/ru/constants.ts.

package ru

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

// RegexParts is REGEX_PARTS: the Unicode-aware word boundaries and the flags of the Russian
// patterns.
var RegexParts = struct {
	LeftBoundary  string
	RightBoundary string
	Flags         string
}{
	LeftBoundary:  "([^\\p{L}\\p{N}_]|^)",
	RightBoundary: "(?=[^\\p{L}\\p{N}_]|$)",
	Flags:         "iu",
}

// WeekdayDictionary is WEEKDAY_DICTIONARY.
var WeekdayDictionary = core.NewDictionary([]core.Entry[float64]{
	{Key: "воскресенье", Value: 0},
	{Key: "воскресенья", Value: 0},
	{Key: "вск", Value: 0},
	{Key: "вск.", Value: 0},
	{Key: "понедельник", Value: 1},
	{Key: "понедельника", Value: 1},
	{Key: "пн", Value: 1},
	{Key: "пн.", Value: 1},
	{Key: "вторник", Value: 2},
	{Key: "вторника", Value: 2},
	{Key: "вт", Value: 2},
	{Key: "вт.", Value: 2},
	{Key: "среда", Value: 3},
	{Key: "среды", Value: 3},
	{Key: "среду", Value: 3},
	{Key: "ср", Value: 3},
	{Key: "ср.", Value: 3},
	{Key: "четверг", Value: 4},
	{Key: "четверга", Value: 4},
	{Key: "чт", Value: 4},
	{Key: "чт.", Value: 4},
	{Key: "пятница", Value: 5},
	{Key: "пятницу", Value: 5},
	{Key: "пятницы", Value: 5},
	{Key: "пт", Value: 5},
	{Key: "пт.", Value: 5},
	{Key: "суббота", Value: 6},
	{Key: "субботу", Value: 6},
	{Key: "субботы", Value: 6},
	{Key: "сб", Value: 6},
	{Key: "сб.", Value: 6},
}...)

// FullMonthNameDictionary is FULL_MONTH_NAME_DICTIONARY.
var FullMonthNameDictionary = core.NewDictionary([]core.Entry[float64]{
	{Key: "январь", Value: 1},
	{Key: "января", Value: 1},
	{Key: "январе", Value: 1},
	{Key: "февраль", Value: 2},
	{Key: "февраля", Value: 2},
	{Key: "феврале", Value: 2},
	{Key: "март", Value: 3},
	{Key: "марта", Value: 3},
	{Key: "марте", Value: 3},
	{Key: "апрель", Value: 4},
	{Key: "апреля", Value: 4},
	{Key: "апреле", Value: 4},
	{Key: "май", Value: 5},
	{Key: "мая", Value: 5},
	{Key: "мае", Value: 5},
	{Key: "июнь", Value: 6},
	{Key: "июня", Value: 6},
	{Key: "июне", Value: 6},
	{Key: "июль", Value: 7},
	{Key: "июля", Value: 7},
	{Key: "июле", Value: 7},
	{Key: "август", Value: 8},
	{Key: "августа", Value: 8},
	{Key: "августе", Value: 8},
	{Key: "сентябрь", Value: 9},
	{Key: "сентября", Value: 9},
	{Key: "сентябре", Value: 9},
	{Key: "октябрь", Value: 10},
	{Key: "октября", Value: 10},
	{Key: "октябре", Value: 10},
	{Key: "ноябрь", Value: 11},
	{Key: "ноября", Value: 11},
	{Key: "ноябре", Value: 11},
	{Key: "декабрь", Value: 12},
	{Key: "декабря", Value: 12},
	{Key: "декабре", Value: 12},
}...)

// MonthDictionary is MONTH_DICTIONARY.
var MonthDictionary = core.NewDictionary(append(FullMonthNameDictionary.Entries(), []core.Entry[float64]{
	{Key: "янв", Value: 1},
	{Key: "янв.", Value: 1},
	{Key: "фев", Value: 2},
	{Key: "фев.", Value: 2},
	{Key: "мар", Value: 3},
	{Key: "мар.", Value: 3},
	{Key: "апр", Value: 4},
	{Key: "апр.", Value: 4},
	{Key: "авг", Value: 8},
	{Key: "авг.", Value: 8},
	{Key: "сен", Value: 9},
	{Key: "сен.", Value: 9},
	{Key: "окт", Value: 10},
	{Key: "окт.", Value: 10},
	{Key: "ноя", Value: 11},
	{Key: "ноя.", Value: 11},
	{Key: "дек", Value: 12},
	{Key: "дек.", Value: 12},
}...)...)

// IntegerWordDictionary is INTEGER_WORD_DICTIONARY.
var IntegerWordDictionary = core.NewDictionary([]core.Entry[float64]{
	{Key: "один", Value: 1},
	{Key: "одна", Value: 1},
	{Key: "одной", Value: 1},
	{Key: "одну", Value: 1},
	{Key: "две", Value: 2},
	{Key: "два", Value: 2},
	{Key: "двух", Value: 2},
	{Key: "три", Value: 3},
	{Key: "трех", Value: 3},
	{Key: "трёх", Value: 3},
	{Key: "четыре", Value: 4},
	{Key: "четырех", Value: 4},
	{Key: "четырёх", Value: 4},
	{Key: "пять", Value: 5},
	{Key: "пяти", Value: 5},
	{Key: "шесть", Value: 6},
	{Key: "шести", Value: 6},
	{Key: "семь", Value: 7},
	{Key: "семи", Value: 7},
	{Key: "восемь", Value: 8},
	{Key: "восьми", Value: 8},
	{Key: "девять", Value: 9},
	{Key: "девяти", Value: 9},
	{Key: "десять", Value: 10},
	{Key: "десяти", Value: 10},
	{Key: "одиннадцать", Value: 11},
	{Key: "одиннадцати", Value: 11},
	{Key: "двенадцать", Value: 12},
	{Key: "двенадцати", Value: 12},
}...)

// OrdinalWordDictionary is ORDINAL_WORD_DICTIONARY.
var OrdinalWordDictionary = core.NewDictionary([]core.Entry[float64]{
	{Key: "первое", Value: 1},
	{Key: "первого", Value: 1},
	{Key: "второе", Value: 2},
	{Key: "второго", Value: 2},
	{Key: "третье", Value: 3},
	{Key: "третьего", Value: 3},
	{Key: "четвертое", Value: 4},
	{Key: "четвертого", Value: 4},
	{Key: "пятое", Value: 5},
	{Key: "пятого", Value: 5},
	{Key: "шестое", Value: 6},
	{Key: "шестого", Value: 6},
	{Key: "седьмое", Value: 7},
	{Key: "седьмого", Value: 7},
	{Key: "восьмое", Value: 8},
	{Key: "восьмого", Value: 8},
	{Key: "девятое", Value: 9},
	{Key: "девятого", Value: 9},
	{Key: "десятое", Value: 10},
	{Key: "десятого", Value: 10},
	{Key: "одиннадцатое", Value: 11},
	{Key: "одиннадцатого", Value: 11},
	{Key: "двенадцатое", Value: 12},
	{Key: "двенадцатого", Value: 12},
	{Key: "тринадцатое", Value: 13},
	{Key: "тринадцатого", Value: 13},
	{Key: "четырнадцатое", Value: 14},
	{Key: "четырнадцатого", Value: 14},
	{Key: "пятнадцатое", Value: 15},
	{Key: "пятнадцатого", Value: 15},
	{Key: "шестнадцатое", Value: 16},
	{Key: "шестнадцатого", Value: 16},
	{Key: "семнадцатое", Value: 17},
	{Key: "семнадцатого", Value: 17},
	{Key: "восемнадцатое", Value: 18},
	{Key: "восемнадцатого", Value: 18},
	{Key: "девятнадцатое", Value: 19},
	{Key: "девятнадцатого", Value: 19},
	{Key: "двадцатое", Value: 20},
	{Key: "двадцатого", Value: 20},
	{Key: "двадцать первое", Value: 21},
	{Key: "двадцать первого", Value: 21},
	{Key: "двадцать второе", Value: 22},
	{Key: "двадцать второго", Value: 22},
	{Key: "двадцать третье", Value: 23},
	{Key: "двадцать третьего", Value: 23},
	{Key: "двадцать четвертое", Value: 24},
	{Key: "двадцать четвертого", Value: 24},
	{Key: "двадцать пятое", Value: 25},
	{Key: "двадцать пятого", Value: 25},
	{Key: "двадцать шестое", Value: 26},
	{Key: "двадцать шестого", Value: 26},
	{Key: "двадцать седьмое", Value: 27},
	{Key: "двадцать седьмого", Value: 27},
	{Key: "двадцать восьмое", Value: 28},
	{Key: "двадцать восьмого", Value: 28},
	{Key: "двадцать девятое", Value: 29},
	{Key: "двадцать девятого", Value: 29},
	{Key: "тридцатое", Value: 30},
	{Key: "тридцатого", Value: 30},
	{Key: "тридцать первое", Value: 31},
	{Key: "тридцать первого", Value: 31},
}...)

// TimeUnitDictionary is TIME_UNIT_DICTIONARY.
var TimeUnitDictionary = core.NewDictionary([]core.Entry[string]{
	{Key: "сек", Value: "second"},
	{Key: "секунда", Value: "second"},
	{Key: "секунд", Value: "second"},
	{Key: "секунды", Value: "second"},
	{Key: "секунду", Value: "second"},
	{Key: "секундочка", Value: "second"},
	{Key: "секундочки", Value: "second"},
	{Key: "секундочек", Value: "second"},
	{Key: "секундочку", Value: "second"},
	{Key: "мин", Value: "minute"},
	{Key: "минута", Value: "minute"},
	{Key: "минут", Value: "minute"},
	{Key: "минуты", Value: "minute"},
	{Key: "минуту", Value: "minute"},
	{Key: "минуток", Value: "minute"},
	{Key: "минутки", Value: "minute"},
	{Key: "минутку", Value: "minute"},
	{Key: "минуточек", Value: "minute"},
	{Key: "минуточки", Value: "minute"},
	{Key: "минуточку", Value: "minute"},
	{Key: "час", Value: "hour"},
	{Key: "часов", Value: "hour"},
	{Key: "часа", Value: "hour"},
	{Key: "часу", Value: "hour"},
	{Key: "часиков", Value: "hour"},
	{Key: "часика", Value: "hour"},
	{Key: "часике", Value: "hour"},
	{Key: "часик", Value: "hour"},
	{Key: "день", Value: "day"},
	{Key: "дня", Value: "day"},
	{Key: "дней", Value: "day"},
	{Key: "суток", Value: "day"},
	{Key: "сутки", Value: "day"},
	{Key: "неделя", Value: "week"},
	{Key: "неделе", Value: "week"},
	{Key: "недели", Value: "week"},
	{Key: "неделю", Value: "week"},
	{Key: "недель", Value: "week"},
	{Key: "недельке", Value: "week"},
	{Key: "недельки", Value: "week"},
	{Key: "неделек", Value: "week"},
	{Key: "месяц", Value: "month"},
	{Key: "месяце", Value: "month"},
	{Key: "месяцев", Value: "month"},
	{Key: "месяца", Value: "month"},
	{Key: "квартал", Value: "quarter"},
	{Key: "квартале", Value: "quarter"},
	{Key: "кварталов", Value: "quarter"},
	{Key: "год", Value: "year"},
	{Key: "года", Value: "year"},
	{Key: "году", Value: "year"},
	{Key: "годов", Value: "year"},
	{Key: "лет", Value: "year"},
	{Key: "годик", Value: "year"},
	{Key: "годика", Value: "year"},
	{Key: "годиков", Value: "year"},
}...)

//-----------------------------

// NumberPattern is NUMBER_PATTERN.
var NumberPattern = "(?:" + core.MatchAnyPattern(IntegerWordDictionary.Keys()) +
	"|[0-9]+|[0-9]+\\.[0-9]+|пол|несколько|пар(?:ы|у)|\\s{0,3})"

var (
	reNeskolko = core.MustRegExp("несколько", "")
	rePol      = core.MustRegExp("пол", "")
	rePar      = core.MustRegExp("пар", "")
)

// ParseNumberPattern parses a NumberPattern match.
func ParseNumberPattern(match string) float64 {
	num := js.ToLowerCase(match)
	if v, ok := IntegerWordDictionary.Get(num); ok {
		return v
	}
	if reNeskolko.Exec(num) != nil {
		return 3
	} else if rePol.Exec(num) != nil {
		return 0.5
	} else if rePar.Exec(num) != nil {
		return 2
	} else if num == "" {
		return 1
	}
	return js.ParseFloat(num)
}

//-----------------------------

// OrdinalNumberPattern is ORDINAL_NUMBER_PATTERN.
var OrdinalNumberPattern = "(?:" + core.MatchAnyPattern(OrdinalWordDictionary.Keys()) + "|[0-9]{1,2}(?:го|ого|е|ое)?)"

// ParseOrdinalNumberPattern parses an OrdinalNumberPattern match.
func ParseOrdinalNumberPattern(match string) float64 {
	num := js.ToLowerCase(match)
	if v, ok := OrdinalWordDictionary.Get(num); ok {
		return v
	}
	return js.ParseInt(num, 0)
}

//-----------------------------

const yearSuffix = "(?:\\s+(?:году|года|год|г|г.))?"

// YearPattern is YEAR_PATTERN.
const YearPattern = "(?:[1-9][0-9]{0,3}" + yearSuffix + "\\s*(?:н.э.|до н.э.|н. э.|до н. э.)|[1-2][0-9]{3}" + yearSuffix + "|[5-9][0-9]" + yearSuffix + ")"

var (
	reYearWord = core.MustRegExp("(год|года|г|г.)", "i")
	reBCE      = core.MustRegExp("(до н.э.|до н. э.)", "i")
	reCE       = core.MustRegExp("(н. э.|н.э.)", "i")
)

// ParseYear parses a YearPattern match.
func ParseYear(match string) float64 {
	if reYearWord.Test(match) {
		match = reYearWord.Replace(match, "")
	}

	if reBCE.Test(match) {
		//Before Common Era
		match = reBCE.Replace(match, "")
		return -js.ParseInt(match, 0)
	}

	if reCE.Test(match) {
		//Common Era
		match = reCE.Replace(match, "")
		return js.ParseInt(match, 0)
	}

	rawYearNumber := js.ParseInt(match, 0)
	return core.FindMostLikelyADYear(rawYearNumber)
}

//-----------------------------

var singleTimeUnitPattern = "(" + NumberPattern + ")\\s{0,3}(" + core.MatchAnyPattern(TimeUnitDictionary.Keys()) + ")"
var singleTimeUnitRegex = core.MustRegExp(singleTimeUnitPattern, "i")

// TimeUnitsPattern is TIME_UNITS_PATTERN.
var TimeUnitsPattern = core.RepeatedTimeunitPattern("(?:(?:около|примерно)\\s{0,3})?", singleTimeUnitPattern, core.DefaultTimeunitConnectorPattern)

// ParseDuration parses a TimeUnitsPattern match (an empty duration when no time unit is found).
func ParseDuration(timeunitText string) *core.Duration {
	fragments := core.NewDuration()
	remainingText := timeunitText
	match := singleTimeUnitRegex.Exec(remainingText)
	for match != nil {
		collectDateTimeFragment(fragments, match)
		remainingText = js.Trim(js.SubstringFrom(remainingText, js.Len(match.Get(0))))
		match = singleTimeUnitRegex.Exec(remainingText)
	}
	return fragments
}

func collectDateTimeFragment(fragments *core.Duration, match *core.RegExpMatch) {
	num := ParseNumberPattern(match.Get(1))
	unit, ok := TimeUnitDictionary.Get(js.ToLowerCase(match.Get(2)))
	if !ok {
		unit = "undefined" // fragments[undefined] = num
	}
	fragments.Set(unit, num)
}
