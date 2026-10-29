// Ported from src/locales/fr/index.ts.

// Package fr is chrono's French locale (parsers, refiners and configuration). The public package
// github.com/wanasit/chrono-go/locales/fr exposes it.
package fr

import "github.com/wanasit/chrono-go/internal/core"

// Casual is the Chrono object configured for parsing casual French
// (createCasualConfiguration() with the JS default littleEndian = true).
var Casual = core.NewChrono(CreateCasualConfiguration(true))

// Strict is the Chrono object configured for parsing strict French
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

// CreateCasualConfiguration creates a casual French configuration: the configuration of
// CreateConfiguration(false, littleEndian) with the casual parsers in front.
func CreateCasualConfiguration(littleEndian bool) *core.Configuration {
	option := CreateConfiguration(false, littleEndian)
	option.Parsers = append([]core.Parser{NewCasualDateParser()}, option.Parsers...)
	option.Parsers = append([]core.Parser{NewCasualTimeParser()}, option.Parsers...)
	option.Parsers = append([]core.Parser{NewTimeUnitRelativeFormatParser()}, option.Parsers...)
	return option
}

// CreateConfiguration creates a French configuration. strictMode is passed to the common
// configuration; littleEndian reads day-first slash dates.
func CreateConfiguration(strictMode, littleEndian bool) *core.Configuration {
	return core.IncludeCommonConfiguration(
		&core.Configuration{
			Parsers: []core.Parser{
				core.NewSlashDateFormatParser(littleEndian),
				NewMonthNameLittleEndianParser(),
				NewTimeExpressionParser(),
				NewSpecificTimeExpressionParser(),
				NewTimeUnitAgoFormatParser(),
				NewTimeUnitWithinFormatParser(),
				NewWeekdayParser(),
			},
			Refiners: []core.Refiner{NewMergeDateTimeRefiner(), NewMergeDateRangeRefiner()},
		},
		strictMode,
	)
}
