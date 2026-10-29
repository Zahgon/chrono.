// Ported from src/locales/uk/constants.ts.

package uk

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

// RegexParts is REGEX_PARTS: the Unicode-aware word boundaries and the flags of the Ukrainian
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
	{Key: "неділя", Value: 0},
	{Key: "неділі", Value: 0},
	{Key: "неділю", Value: 0},
	{Key: "нд", Value: 0},
	{Key: "нд.", Value: 0},
	{Key: "понеділок", Value: 1},
	{Key: "понеділка", Value: 1},
	{Key: "пн", Value: 1},
	{Key: "пн.", Value: 1},
	{Key: "вівторок", Value: 2},
	{Key: "вівторка", Value: 2},
	{Key: "вт", Value: 2},
	{Key: "вт.", Value: 2},
	{Key: "середа", Value: 3},
	{Key: "середи", Value: 3},
	{Key: "середу", Value: 3},
	{Key: "ср", Value: 3},
	{Key: "ср.", Value: 3},
	{Key: "четвер", Value: 4},
	{Key: "четверга", Value: 4},
	{Key: "четвергу", Value: 4},
	{Key: "чт", Value: 4},
	{Key: "чт.", Value: 4},
	{Key: "п'ятниця", Value: 5},
	{Key: "п'ятниці", Value: 5},
	{Key: "п'ятницю", Value: 5},
	{Key: "пт", Value: 5},
	{Key: "пт.", Value: 5},
	{Key: "субота", Value: 6},
	{Key: "суботи", Value: 6},
	{Key: "суботу", Value: 6},
	{Key: "сб", Value: 6},
	{Key: "сб.", Value: 6},
}...)

// FullMonthNameDictionary is FULL_MONTH_NAME_DICTIONARY.
var FullMonthNameDictionary = core.NewDictionary([]core.Entry[float64]{
	{Key: "січень", Value: 1},
	{Key: "січня", Value: 1},
	{Key: "січні", Value: 1},
	{Key: "лютий", Value: 2},
	{Key: "лютого", Value: 2},
	{Key: "лютому", Value: 2},
	{Key: "березень", Value: 3},
	{Key: "березня", Value: 3},
	{Key: "березні", Value: 3},
	{Key: "квітень", Value: 4},
	{Key: "квітня", Value: 4},
	{Key: "квітні", Value: 4},
	{Key: "травень", Value: 5},
	{Key: "травня", Value: 5},
	{Key: "травні", Value: 5},
	{Key: "червень", Value: 6},
	{Key: "червня", Value: 6},
	{Key: "червні", Value: 6},
	{Key: "липень", Value: 7},
	{Key: "липня", Value: 7},
	{Key: "липні", Value: 7},
	{Key: "серпень", Value: 8},
	{Key: "серпня", Value: 8},
	{Key: "серпні", Value: 8},
	{Key: "вересень", Value: 9},
	{Key: "вересня", Value: 9},
	{Key: "вересні", Value: 9},
	{Key: "жовтень", Value: 10},
	{Key: "жовтня", Value: 10},
	{Key: "жовтні", Value: 10},
	{Key: "листопад", Value: 11},
	{Key: "листопада", Value: 11},
	{Key: "листопаду", Value: 11},
	{Key: "грудень", Value: 12},
	{Key: "грудня", Value: 12},
	{Key: "грудні", Value: 12},
}...)

// MonthDictionary is MONTH_DICTIONARY.
var MonthDictionary = core.NewDictionary(append(FullMonthNameDictionary.Entries(), []core.Entry[float64]{
	{Key: "січ", Value: 1},
	{Key: "січ.", Value: 1},
	{Key: "лют", Value: 2},
	{Key: "лют.", Value: 2},
	{Key: "бер", Value: 3},
	{Key: "бер.", Value: 3},
	{Key: "квіт", Value: 4},
	{Key: "квіт.", Value: 4},
	{Key: "трав", Value: 5},
	{Key: "трав.", Value: 5},
	{Key: "черв", Value: 6},
	{Key: "черв.", Value: 6},
	{Key: "лип", Value: 7},
	{Key: "лип.", Value: 7},
	{Key: "серп", Value: 8},
	{Key: "серп.", Value: 8},
	{Key: "сер", Value: 8},
	{Key: "cер.", Value: 8},
	{Key: "вер", Value: 9},
	{Key: "вер.", Value: 9},
	{Key: "верес", Value: 9},
	{Key: "верес.", Value: 9},
	{Key: "жовт", Value: 10},
	{Key: "жовт.", Value: 10},
	{Key: "листоп", Value: 11},
	{Key: "листоп.", Value: 11},
	{Key: "груд", Value: 12},
	{Key: "груд.", Value: 12},
}...)...)

// IntegerWordDictionary is INTEGER_WORD_DICTIONARY.
var IntegerWordDictionary = core.NewDictionary([]core.Entry[float64]{
	{Key: "один", Value: 1},
	{Key: "одна", Value: 1},
	{Key: "одної", Value: 1},
	{Key: "одну", Value: 1},
	{Key: "дві", Value: 2},
	{Key: "два", Value: 2},
	{Key: "двох", Value: 2},
	{Key: "три", Value: 3},
	{Key: "трьох", Value: 3},
	{Key: "чотири", Value: 4},
	{Key: "чотирьох", Value: 4},
	{Key: "п'ять", Value: 5},
	{Key: "п'яти", Value: 5},
	{Key: "шість", Value: 6},
	{Key: "шести", Value: 6},
	{Key: "сім", Value: 7},
	{Key: "семи", Value: 7},
	{Key: "вісім", Value: 8},
	{Key: "восьми", Value: 8},
	{Key: "дев'ять", Value: 9},
	{Key: "дев'яти", Value: 9},
	{Key: "десять", Value: 10},
	{Key: "десяти", Value: 10},
	{Key: "одинадцять", Value: 11},
	{Key: "одинадцяти", Value: 11},
	{Key: "дванадцять", Value: 12},
	{Key: "дванадцяти", Value: 12},
}...)

// OrdinalWordDictionary is ORDINAL_WORD_DICTIONARY.
var OrdinalWordDictionary = core.NewDictionary([]core.Entry[float64]{
	{Key: "перше", Value: 1},
	{Key: "першого", Value: 1},
	{Key: "друге", Value: 2},
	{Key: "другого", Value: 2},
	{Key: "третє", Value: 3},
	{Key: "третього", Value: 3},
	{Key: "четверте", Value: 4},
	{Key: "четвертого", Value: 4},
	{Key: "п'яте", Value: 5},
	{Key: "п'ятого", Value: 5},
	{Key: "шосте", Value: 6},
	{Key: "шостого", Value: 6},
	{Key: "сьоме", Value: 7},
	{Key: "сьомого", Value: 7},
	{Key: "восьме", Value: 8},
	{Key: "восьмого", Value: 8},
	{Key: "дев'яте", Value: 9},
	{Key: "дев'ятого", Value: 9},
	{Key: "десяте", Value: 10},
	{Key: "десятого", Value: 10},
	{Key: "одинадцяте", Value: 11},
	{Key: "одинадцятого", Value: 11},
	{Key: "дванадцяте", Value: 12},
	{Key: "дванадцятого", Value: 12},
	{Key: "тринадцяте", Value: 13},
	{Key: "тринадцятого", Value: 13},
	{Key: "чотирнадцяте", Value: 14},
	{Key: "чотинрнадцятого", Value: 14},
	{Key: "п'ятнадцяте", Value: 15},
	{Key: "п'ятнадцятого", Value: 15},
	{Key: "шістнадцяте", Value: 16},
	{Key: "шістнадцятого", Value: 16},
	{Key: "сімнадцяте", Value: 17},
	{Key: "сімнадцятого", Value: 17},
	{Key: "вісімнадцяте", Value: 18},
	{Key: "вісімнадцятого", Value: 18},
	{Key: "дев'ятнадцяте", Value: 19},
	{Key: "дев'ятнадцятого", Value: 19},
	{Key: "двадцяте", Value: 20},
	{Key: "двадцятого", Value: 20},
	{Key: "двадцять перше", Value: 21},
	{Key: "двадцять першого", Value: 21},
	{Key: "двадцять друге", Value: 22},
	{Key: "двадцять другого", Value: 22},
	{Key: "двадцять третє", Value: 23},
	{Key: "двадцять третього", Value: 23},
	{Key: "двадцять четверте", Value: 24},
	{Key: "двадцять четвертого", Value: 24},
	{Key: "двадцять п'яте", Value: 25},
	{Key: "двадцять п'ятого", Value: 25},
	{Key: "двадцять шосте", Value: 26},
	{Key: "двадцять шостого", Value: 26},
	{Key: "двадцять сьоме", Value: 27},
	{Key: "двадцять сьомого", Value: 27},
	{Key: "двадцять восьме", Value: 28},
	{Key: "двадцять восьмого", Value: 28},
	{Key: "двадцять дев'яте", Value: 29},
	{Key: "двадцять дев'ятого", Value: 29},
	{Key: "тридцяте", Value: 30},
	{Key: "тридцятого", Value: 30},
	{Key: "тридцять перше", Value: 31},
	{Key: "тридцять першого", Value: 31},
}...)

// TimeUnitDictionary is TIME_UNIT_DICTIONARY.
var TimeUnitDictionary = core.NewDictionary([]core.Entry[string]{
	{Key: "сек", Value: "second"},
	{Key: "секунда", Value: "second"},
	{Key: "секунд", Value: "second"},
	{Key: "секунди", Value: "second"},
	{Key: "секунду", Value: "second"},
	{Key: "секундочок", Value: "second"},
	{Key: "секундочки", Value: "second"},
	{Key: "секундочку", Value: "second"},
	{Key: "хв", Value: "minute"},
	{Key: "хвилина", Value: "minute"},
	{Key: "хвилин", Value: "minute"},
	{Key: "хвилини", Value: "minute"},
	{Key: "хвилину", Value: "minute"},
	{Key: "хвилинок", Value: "minute"},
	{Key: "хвилинки", Value: "minute"},
	{Key: "хвилинку", Value: "minute"},
	{Key: "хвилиночок", Value: "minute"},
	{Key: "хвилиночки", Value: "minute"},
	{Key: "хвилиночку", Value: "minute"},
	{Key: "год", Value: "hour"},
	{Key: "година", Value: "hour"},
	{Key: "годин", Value: "hour"},
	{Key: "години", Value: "hour"},
	{Key: "годину", Value: "hour"},
	{Key: "годинка", Value: "hour"},
	{Key: "годинок", Value: "hour"},
	{Key: "годинки", Value: "hour"},
	{Key: "годинку", Value: "hour"},
	{Key: "день", Value: "day"},
	{Key: "дня", Value: "day"},
	{Key: "днів", Value: "day"},
	{Key: "дні", Value: "day"},
	{Key: "доба", Value: "day"},
	{Key: "добу", Value: "day"},
	{Key: "тиждень", Value: "week"},
	{Key: "тижню", Value: "week"},
	{Key: "тижня", Value: "week"},
	{Key: "тижні", Value: "week"},
	{Key: "тижнів", Value: "week"},
	{Key: "місяць", Value: "month"},
	{Key: "місяців", Value: "month"},
	{Key: "місяці", Value: "month"},
	{Key: "місяця", Value: "month"},
	{Key: "квартал", Value: "quarter"},
	{Key: "кварталу", Value: "quarter"},
	{Key: "квартала", Value: "quarter"},
	{Key: "кварталів", Value: "quarter"},
	{Key: "кварталі", Value: "quarter"},
	{Key: "рік", Value: "year"},
	{Key: "року", Value: "year"},
	{Key: "році", Value: "year"},
	{Key: "років", Value: "year"},
	{Key: "роки", Value: "year"},
}...)

//--------------------------------

// NumberPattern is NUMBER_PATTERN.
var NumberPattern = "(?:" + core.MatchAnyPattern(IntegerWordDictionary.Keys()) +
	"|[0-9]+|[0-9]+\\.[0-9]+|пів|декілька|пар(?:у)|\\s{0,3})"

var (
	reDekilka = core.MustRegExp("декілька", "")
	rePiv     = core.MustRegExp("пів", "")
	rePar     = core.MustRegExp("пар", "")
)

// ParseNumberPattern parses a NumberPattern match.
func ParseNumberPattern(match string) float64 {
	num := js.ToLowerCase(match)
	if v, ok := IntegerWordDictionary.Get(num); ok {
		return v
	}
	if reDekilka.Exec(num) != nil {
		return 2
	} else if rePiv.Exec(num) != nil {
		return 0.5
	} else if rePar.Exec(num) != nil {
		return 2
	} else if num == "" {
		return 1
	}
	return js.ParseFloat(num)
}

// OrdinalNumberPattern is ORDINAL_NUMBER_PATTERN.
var OrdinalNumberPattern = "(?:" + core.MatchAnyPattern(OrdinalWordDictionary.Keys()) + "|[0-9]{1,2}(?:го|ого|е)?)"

// ParseOrdinalNumberPattern parses an OrdinalNumberPattern match.
func ParseOrdinalNumberPattern(match string) float64 {
	num := js.ToLowerCase(match)
	if v, ok := OrdinalWordDictionary.Get(num); ok {
		return v
	}
	return js.ParseInt(num, 0)
}

const yearSuffix = "(?:\\s+(?:року|рік|р|р.))?"

// YearPattern is YEAR_PATTERN.
const YearPattern = "(?:[1-9][0-9]{0,3}" + yearSuffix + "\\s*(?:н.е.|до н.е.|н. е.|до н. е.)|[1-2][0-9]{3}" + yearSuffix + "|[5-9][0-9]" + yearSuffix + ")"

var (
	reYearWord = core.MustRegExp("(рік|року|р|р.)", "i")
	reBCE      = core.MustRegExp("(до н.е.|до н. е.)", "i")
	reCE       = core.MustRegExp("(н. е.|н.е.)", "i")
)

// ParseYearPattern parses a YearPattern match.
func ParseYearPattern(match string) float64 {
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

var singleTimeUnitPattern = "(" + NumberPattern + ")\\s{0,3}(" + core.MatchAnyPattern(TimeUnitDictionary.Keys()) + ")"
var singleTimeUnitRegex = core.MustRegExp(singleTimeUnitPattern, "i")

// TimeUnitsPattern is TIME_UNITS_PATTERN.
var TimeUnitsPattern = core.RepeatedTimeunitPattern(
	"(?:(?:близько|приблизно)\\s{0,3})?",
	singleTimeUnitPattern,
	core.DefaultTimeunitConnectorPattern,
)

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
