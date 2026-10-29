// Ported from src/locales/nl/index.ts.

// Package nl is chrono's Dutch locale (parsers, refiners and configuration). The public package
// github.com/wanasit/chrono-go/locales/nl exposes it.
package nl

import "github.com/wanasit/chrono-go/internal/core"

// Casual is the Chrono object configured for parsing casual Dutch.
var Casual = core.NewChrono(CreateCasualConfiguration(true))

// Strict is the Chrono object configured for parsing strict Dutch.
var Strict = core.NewChrono(CreateConfiguration(true, true))

// Parse is Casual.Parse.
func Parse(text string, ref core.ParsingReference, option *core.ParsingOption) []*core.ParsingResult {
	return Casual.Parse(text, ref, option)
}

// ParseDate is Casual.ParseDate.
func ParseDate(text string, ref core.ParsingReference, option *core.ParsingOption) (float64, bool) {
	return Casual.ParseDate(text, ref, option)
}

// CreateCasualConfiguration creates a casual Dutch configuration: the configuration of
// CreateConfiguration(false, littleEndian) with the casual parsers in front.
func CreateCasualConfiguration(littleEndian bool) *core.Configuration {
	option := CreateConfiguration(false, littleEndian)
	option.Parsers = append([]core.Parser{NewCasualDateParser()}, option.Parsers...)
	option.Parsers = append([]core.Parser{NewCasualTimeParser()}, option.Parsers...)
	option.Parsers = append([]core.Parser{NewCasualDateTimeParser()}, option.Parsers...)
	option.Parsers = append([]core.Parser{NewMonthNameParser()}, option.Parsers...)
	option.Parsers = append([]core.Parser{NewRelativeDateFormatParser()}, option.Parsers...)
	option.Parsers = append([]core.Parser{NewTimeUnitCasualRelativeFormatParser()}, option.Parsers...)
	return option
}

// CreateConfiguration creates a Dutch configuration. strictMode makes the time expressions and
// time unit mentions strict; littleEndian reads day-first slash dates.
func CreateConfiguration(strictMode, littleEndian bool) *core.Configuration {
	return core.IncludeCommonConfiguration(
		&core.Configuration{
			Parsers: []core.Parser{
				core.NewSlashDateFormatParser(littleEndian),
				NewTimeUnitWithinFormatParser(),
				NewMonthNameMiddleEndianParser(),
				NewMonthNameParser(),
				NewWeekdayParser(),
				NewCasualYearMonthDayParser(),
				NewSlashMonthFormatParser(),
				NewTimeExpressionParser(strictMode),
				NewTimeUnitAgoFormatParser(strictMode),
				NewTimeUnitLaterFormatParser(strictMode),
			},
			Refiners: []core.Refiner{NewMergeDateTimeRefiner(), NewMergeDateRangeRefiner()},
		},
		strictMode,
	)
}
