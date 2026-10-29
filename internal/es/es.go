// Ported from src/locales/es/index.ts.

// Package es is chrono's Spanish locale (parsers, refiners and configuration). The public package
// github.com/wanasit/chrono-go/locales/es exposes it.
package es

import "github.com/wanasit/chrono-go/internal/core"

// Casual is the Chrono object configured for parsing casual Spanish
// (createCasualConfiguration() with the JS default littleEndian = true).
var Casual = core.NewChrono(CreateCasualConfiguration(true))

// Strict is the Chrono object configured for parsing strict Spanish
// (createConfiguration(true) with the JS default littleEndian = true).
var Strict = core.NewChrono(CreateConfiguration(true, true))

// Parse is Casual.Parse.
func Parse(text string, ref core.ParsingReference, option *core.ParsingOption) []*core.ParsingResult {
	return Casual.Parse(text, ref, option)
}

// ParseDate is Casual.ParseDate.
func ParseDate(text string, ref core.ParsingReference, option *core.ParsingOption) (float64, bool) {
	return Casual.ParseDate(text, ref, option)
}

// CreateCasualConfiguration creates a casual Spanish configuration: the configuration of
// CreateConfiguration(false, littleEndian) plus the casual parsers.
func CreateCasualConfiguration(littleEndian bool) *core.Configuration {
	option := CreateConfiguration(false, littleEndian)
	option.Parsers = append(option.Parsers, NewCasualDateParser())
	option.Parsers = append(option.Parsers, NewCasualTimeParser())
	return option
}

// CreateConfiguration creates a Spanish configuration. strictMode is passed to the common
// configuration; littleEndian reads day-first slash dates.
func CreateConfiguration(strictMode, littleEndian bool) *core.Configuration {
	return core.IncludeCommonConfiguration(
		&core.Configuration{
			Parsers: []core.Parser{
				core.NewSlashDateFormatParser(littleEndian),
				NewWeekdayParser(),
				NewTimeExpressionParser(),
				NewMonthNameLittleEndianParser(),
				NewTimeUnitWithinFormatParser(),
			},
			Refiners: []core.Refiner{NewMergeDateTimeRefiner(), NewMergeDateRangeRefiner()},
		},
		strictMode,
	)
}
