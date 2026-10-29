// Ported from src/locales/it/constants.ts.

package it

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

func concat[V any](parts ...[]core.Entry[V]) []core.Entry[V] {
	var out []core.Entry[V]
	for _, p := range parts {
		out = append(out, p...)
	}
	return out
}

// WeekdayDictionary is WEEKDAY_DICTIONARY.
var WeekdayDictionary = core.NewDictionary([]core.Entry[float64]{
	{Key: "domenica", Value: 0},
	{Key: "dom", Value: 0},
	{Key: "dom.", Value: 0},
	{Key: "lunedì", Value: 1},
	{Key: "lunedi", Value: 1},
	{Key: "lun", Value: 1},
	{Key: "lun.", Value: 1},
	{Key: "martedì", Value: 2},
	{Key: "martedi", Value: 2},
	{Key: "mar", Value: 2},
	{Key: "mar.", Value: 2},
	{Key: "mercoledì", Value: 3},
	{Key: "mercoledi", Value: 3},
	{Key: "mer", Value: 3},
	{Key: "mer.", Value: 3},
	{Key: "giovedì", Value: 4},
	{Key: "giovedi", Value: 4},
	{Key: "gio", Value: 4},
	{Key: "gio.", Value: 4},
	{Key: "venerdì", Value: 5},
	{Key: "venerdi", Value: 5},
	{Key: "ven", Value: 5},
	{Key: "ven.", Value: 5},
	{Key: "sabato", Value: 6},
	{Key: "sab", Value: 6},
	{Key: "sab.", Value: 6},
}...)

// FullMonthNameDictionary is FULL_MONTH_NAME_DICTIONARY.
var FullMonthNameDictionary = core.NewDictionary([]core.Entry[float64]{
	{Key: "gennaio", Value: 1},
	{Key: "febbraio", Value: 2},
	{Key: "marzo", Value: 3},
	{Key: "aprile", Value: 4},
	{Key: "maggio", Value: 5},
	{Key: "giugno", Value: 6},
	{Key: "luglio", Value: 7},
	{Key: "agosto", Value: 8},
	{Key: "settembre", Value: 9},
	{Key: "ottobre", Value: 10},
	{Key: "novembre", Value: 11},
	{Key: "dicembre", Value: 12},
}...)

// MonthDictionary is MONTH_DICTIONARY.
var MonthDictionary = core.NewDictionary(concat(FullMonthNameDictionary.Entries(), []core.Entry[float64]{
	{Key: "gen", Value: 1},
	{Key: "gen.", Value: 1},
	{Key: "feb", Value: 2},
	{Key: "feb.", Value: 2},
	{Key: "mar", Value: 3},
	{Key: "mar.", Value: 3},
	{Key: "apr", Value: 4},
	{Key: "apr.", Value: 4},
	{Key: "mag", Value: 5},
	{Key: "mag.", Value: 5},
	{Key: "giu", Value: 6},
	{Key: "giu.", Value: 6},
	{Key: "lug", Value: 7},
	{Key: "lug.", Value: 7},
	{Key: "ago", Value: 8},
	{Key: "ago.", Value: 8},
	{Key: "set", Value: 9},
	{Key: "set.", Value: 9},
	{Key: "sett", Value: 9},
	{Key: "sett.", Value: 9},
	{Key: "ott", Value: 10},
	{Key: "ott.", Value: 10},
	{Key: "nov", Value: 11},
	{Key: "nov.", Value: 11},
	{Key: "dic", Value: 12},
	{Key: "dic.", Value: 12},
})...)

// IntegerWordDictionary is INTEGER_WORD_DICTIONARY.
var IntegerWordDictionary = core.NewDictionary([]core.Entry[float64]{
	{Key: "uno", Value: 1},
	{Key: "una", Value: 1},
	{Key: "un", Value: 1},
	{Key: "due", Value: 2},
	{Key: "tre", Value: 3},
	{Key: "quattro", Value: 4},
	{Key: "cinque", Value: 5},
	{Key: "sei", Value: 6},
	{Key: "sette", Value: 7},
	{Key: "otto", Value: 8},
	{Key: "nove", Value: 9},
	{Key: "dieci", Value: 10},
	{Key: "undici", Value: 11},
	{Key: "dodici", Value: 12},
}...)

// OrdinalWordDictionary is ORDINAL_WORD_DICTIONARY.
var OrdinalWordDictionary = core.NewDictionary([]core.Entry[float64]{
	{Key: "primo", Value: 1},
	{Key: "prima", Value: 1},
	{Key: "1°", Value: 1},
	{Key: "1ª", Value: 1},
	{Key: "secondo", Value: 2},
	{Key: "seconda", Value: 2},
	{Key: "2°", Value: 2},
	{Key: "2ª", Value: 2},
	{Key: "terzo", Value: 3},
	{Key: "terza", Value: 3},
	{Key: "3°", Value: 3},
	{Key: "3ª", Value: 3},
	{Key: "quarto", Value: 4},
	{Key: "quarta", Value: 4},
	{Key: "4°", Value: 4},
	{Key: "4ª", Value: 4},
	{Key: "quinto", Value: 5},
	{Key: "quinta", Value: 5},
	{Key: "5°", Value: 5},
	{Key: "5ª", Value: 5},
	{Key: "sesto", Value: 6},
	{Key: "sesta", Value: 6},
	{Key: "6°", Value: 6},
	{Key: "6ª", Value: 6},
	{Key: "settimo", Value: 7},
	{Key: "settima", Value: 7},
	{Key: "7°", Value: 7},
	{Key: "7ª", Value: 7},
	{Key: "ottavo", Value: 8},
	{Key: "ottava", Value: 8},
	{Key: "8°", Value: 8},
	{Key: "8ª", Value: 8},
	{Key: "nono", Value: 9},
	{Key: "nona", Value: 9},
	{Key: "9°", Value: 9},
	{Key: "9ª", Value: 9},
	{Key: "decimo", Value: 10},
	{Key: "decima", Value: 10},
	{Key: "10°", Value: 10},
	{Key: "10ª", Value: 10},
	{Key: "undicesimo", Value: 11},
	{Key: "undicesima", Value: 11},
	{Key: "11°", Value: 11},
	{Key: "11ª", Value: 11},
	{Key: "dodicesimo", Value: 12},
	{Key: "dodicesima", Value: 12},
	{Key: "12°", Value: 12},
	{Key: "12ª", Value: 12},
	{Key: "tredicesimo", Value: 13},
	{Key: "tredicesima", Value: 13},
	{Key: "13°", Value: 13},
	{Key: "13ª", Value: 13},
	{Key: "quattordicesimo", Value: 14},
	{Key: "quattordicesima", Value: 14},
	{Key: "14°", Value: 14},
	{Key: "14ª", Value: 14},
	{Key: "quindicesimo", Value: 15},
	{Key: "quindicesima", Value: 15},
	{Key: "15°", Value: 15},
	{Key: "15ª", Value: 15},
	{Key: "sedicesimo", Value: 16},
	{Key: "sedicesima", Value: 16},
	{Key: "16°", Value: 16},
	{Key: "16ª", Value: 16},
	{Key: "diciassettesimo", Value: 17},
	{Key: "diciassettesima", Value: 17},
	{Key: "17°", Value: 17},
	{Key: "17ª", Value: 17},
	{Key: "diciottesimo", Value: 18},
	{Key: "diciottesima", Value: 18},
	{Key: "18°", Value: 18},
	{Key: "18ª", Value: 18},
	{Key: "diciannovesimo", Value: 19},
	{Key: "diciannovesima", Value: 19},
	{Key: "19°", Value: 19},
	{Key: "19ª", Value: 19},
	{Key: "ventesimo", Value: 20},
	{Key: "ventesima", Value: 20},
	{Key: "20°", Value: 20},
	{Key: "20ª", Value: 20},
	{Key: "ventunesimo", Value: 21},
	{Key: "ventunesima", Value: 21},
	{Key: "21°", Value: 21},
	{Key: "21ª", Value: 21},
	{Key: "ventiduesimo", Value: 22},
	{Key: "ventiduesima", Value: 22},
	{Key: "22°", Value: 22},
	{Key: "22ª", Value: 22},
	{Key: "ventitreesimo", Value: 23},
	{Key: "ventitreesima", Value: 23},
	{Key: "23°", Value: 23},
	{Key: "23ª", Value: 23},
	{Key: "ventiquattresimo", Value: 24},
	{Key: "ventiquattresima", Value: 24},
	{Key: "24°", Value: 24},
	{Key: "24ª", Value: 24},
	{Key: "venticinquesimo", Value: 25},
	{Key: "venticinquesima", Value: 25},
	{Key: "25°", Value: 25},
	{Key: "25ª", Value: 25},
	{Key: "ventiseiesimo", Value: 26},
	{Key: "ventiseiesima", Value: 26},
	{Key: "26°", Value: 26},
	{Key: "26ª", Value: 26},
	{Key: "ventisettesimo", Value: 27},
	{Key: "ventisettesima", Value: 27},
	{Key: "27°", Value: 27},
	{Key: "27ª", Value: 27},
	{Key: "ventottesimo", Value: 28},
	{Key: "ventottesima", Value: 28},
	{Key: "28°", Value: 28},
	{Key: "28ª", Value: 28},
	{Key: "ventinovesimo", Value: 29},
	{Key: "ventinovesima", Value: 29},
	{Key: "29°", Value: 29},
	{Key: "29ª", Value: 29},
	{Key: "trentesimo", Value: 30},
	{Key: "trentesima", Value: 30},
	{Key: "30°", Value: 30},
	{Key: "30ª", Value: 30},
	{Key: "trentunesimo", Value: 31},
	{Key: "trentunesima", Value: 31},
	{Key: "31°", Value: 31},
	{Key: "31ª", Value: 31},
}...)

// TimeUnitDictionaryNoAbbr is TIME_UNIT_DICTIONARY_NO_ABBR.
var TimeUnitDictionaryNoAbbr = core.NewDictionary([]core.Entry[string]{
	{Key: "secondo", Value: "second"},
	{Key: "secondi", Value: "second"},
	{Key: "minuto", Value: "minute"},
	{Key: "minuti", Value: "minute"},
	{Key: "ora", Value: "hour"},
	{Key: "ore", Value: "hour"},
	{Key: "giorno", Value: "day"},
	{Key: "giorni", Value: "day"},
	{Key: "settimana", Value: "week"},
	{Key: "settimane", Value: "week"},
	{Key: "mese", Value: "month"},
	{Key: "mesi", Value: "month"},
	{Key: "trimestre", Value: "quarter"},
	{Key: "trimestri", Value: "quarter"},
	{Key: "anno", Value: "year"},
	{Key: "anni", Value: "year"},
}...)

// TimeUnitDictionary is TIME_UNIT_DICTIONARY.
var TimeUnitDictionary = core.NewDictionary(concat([]core.Entry[string]{
	{Key: "s", Value: "second"},
	{Key: "sec", Value: "second"},
	{Key: "secondo", Value: "second"},
	{Key: "secondi", Value: "second"},
	{Key: "m", Value: "minute"},
	{Key: "min", Value: "minute"},
	{Key: "minuto", Value: "minute"},
	{Key: "minuti", Value: "minute"},
	{Key: "h", Value: "hour"},
	{Key: "ora", Value: "hour"},
	{Key: "ore", Value: "hour"},
	{Key: "g", Value: "day"},
	{Key: "gg", Value: "day"},
	{Key: "giorno", Value: "day"},
	{Key: "giorni", Value: "day"},
	{Key: "sett", Value: "week"},
	{Key: "settimana", Value: "week"},
	{Key: "settimane", Value: "week"},
	{Key: "mese", Value: "month"},
	{Key: "mesi", Value: "month"},
	{Key: "trim", Value: "quarter"},
	{Key: "trimestre", Value: "quarter"},
	{Key: "trimestri", Value: "quarter"},
	{Key: "anno", Value: "year"},
	{Key: "anni", Value: "year"},
}, TimeUnitDictionaryNoAbbr.Entries())...)

// monthLookup is MONTH_DICTIONARY[name] (undefined when missing).
func monthLookup(name string) float64 {
	if v, ok := MonthDictionary.Get(name); ok {
		return v
	}
	return core.Undefined
}

// timeunitLookup is TIME_UNIT_DICTIONARY[word] used as an object key ("undefined" when missing).
func timeunitLookup(word string) string {
	if v, ok := TimeUnitDictionary.Get(word); ok {
		return v
	}
	return "undefined"
}

//-----------------------------

// NumberPattern is NUMBER_PATTERN.
var NumberPattern = "(?:" + core.MatchAnyPattern(IntegerWordDictionary.Keys()) +
	"|[0-9]+|[0-9]+\\.[0-9]+|mezz[oa]?(?:\\s{0,2})?|un(?:'|\\s{0,2})?|qualche|alcuni|paio\\s{0,2}(?:di)?)"

var (
	reQualche = core.MustRegExp("qualche", "")
	reMezz    = core.MustRegExp("mezz", "")
	rePaio    = core.MustRegExp("paio", "")
	reAlcuni  = core.MustRegExp("alcuni", "")
)

// ParseNumberPattern parses a NumberPattern match.
func ParseNumberPattern(match string) float64 {
	num := js.ToLowerCase(match)
	if v, ok := IntegerWordDictionary.Get(num); ok {
		return v
	} else if num == "un" || num == "un'" || num == "una" || num == "uno" {
		return 1
	} else if reQualche.Exec(num) != nil {
		return 3
	} else if reMezz.Exec(num) != nil {
		return 0.5
	} else if rePaio.Exec(num) != nil {
		return 2
	} else if reAlcuni.Exec(num) != nil {
		return 7
	}
	return js.ParseFloat(num)
}

//-----------------------------

// OrdinalNumberPattern is ORDINAL_NUMBER_PATTERN.
var OrdinalNumberPattern = "(?:" + core.MatchAnyPattern(OrdinalWordDictionary.Keys()) + "|[0-9]{1,2}(?:°|ª|º)?)"

var reOrdinalSuffix = core.MustRegExp("(?:°|ª|º)$", "i")

// ParseOrdinalNumberPattern parses an OrdinalNumberPattern match.
func ParseOrdinalNumberPattern(match string) float64 {
	num := js.ToLowerCase(match)
	if v, ok := OrdinalWordDictionary.Get(num); ok {
		return v
	}

	num = reOrdinalSuffix.Replace(num, "")
	return js.ParseInt(num, 0)
}

//-----------------------------

// YearPattern is YEAR_PATTERN.
const YearPattern = "(?:[1-9][0-9]{0,3}\\s{0,2}(?:BE|AD|BC|BCE|CE)|[1-2][0-9]{3}|[5-9][0-9]|2[0-5])"

var (
	reBE   = core.MustRegExp("BE", "i")
	reBCE  = core.MustRegExp("BCE?", "i")
	reADCE = core.MustRegExp("(AD|CE)", "i")
)

// ParseYear parses a YearPattern match.
func ParseYear(match string) float64 {
	if reBE.Test(match) {
		// Buddhist Era
		match = reBE.Replace(match, "")
		return js.ParseInt(match, 0) - 543
	}

	if reBCE.Test(match) {
		// Before Christ, Before Common Era
		match = reBCE.Replace(match, "")
		return -js.ParseInt(match, 0)
	}

	if reADCE.Test(match) {
		// Anno Domini, Common Era
		match = reADCE.Replace(match, "")
		return js.ParseInt(match, 0)
	}

	rawYearNumber := js.ParseInt(match, 0)
	return core.FindMostLikelyADYear(rawYearNumber)
}

//-----------------------------

var singleTimeUnitPattern = "(" + NumberPattern + ")\\s{0,3}(" + core.MatchAnyPattern(TimeUnitDictionary.Keys()) + ")"
var singleTimeUnitRegex = core.MustRegExp(singleTimeUnitPattern, "i")

var singleTimeUnitNoAbbrPattern = "(" + NumberPattern + ")\\s{0,3}(" + core.MatchAnyPattern(TimeUnitDictionaryNoAbbr.Keys()) + ")"

const timeUnitConnectorPattern = "\\s{0,5},?(?:\\s*e)?\\s{0,5}"

// TimeUnitsPattern is TIME_UNITS_PATTERN.
var TimeUnitsPattern = core.RepeatedTimeunitPattern(
	"(?:(?:circa|approssimativamente)\\s{0,3})?",
	singleTimeUnitPattern,
	timeUnitConnectorPattern,
)

// TimeUnitsNoAbbrPattern is TIME_UNITS_NO_ABBR_PATTERN.
var TimeUnitsNoAbbrPattern = core.RepeatedTimeunitPattern(
	"(?:(?:circa|approssimativamente)\\s{0,3})?",
	singleTimeUnitNoAbbrPattern,
	timeUnitConnectorPattern,
)

// ParseDuration parses a TimeUnitsPattern match (nil when no time unit is found).
func ParseDuration(timeunitText string) *core.Duration {
	fragments := core.NewDuration()
	remainingText := timeunitText
	match := singleTimeUnitRegex.Exec(remainingText)
	for match != nil {
		collectDateTimeFragment(fragments, match)
		remainingText = js.Trim(js.SubstringFrom(remainingText, js.Len(match.Get(0))))
		match = singleTimeUnitRegex.Exec(remainingText)
	}
	if fragments.Len() == 0 {
		return nil
	}
	return fragments
}

var reLettersOnly = core.MustRegExp("^[a-zA-Z]+$", "")

func collectDateTimeFragment(fragments *core.Duration, match *core.RegExpMatch) {
	if reLettersOnly.Exec(match.Get(0)) != nil {
		return
	}
	num := ParseNumberPattern(match.Get(1))
	unit := timeunitLookup(js.ToLowerCase(match.Get(2)))
	fragments.Set(unit, num)
}
