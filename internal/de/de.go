// Ported from src/locales/de/index.ts.

// Package de is chrono's German locale (parsers, refiners and configuration). The public package
// github.com/wanasit/chrono-go/locales/de exposes it.
package de

import "github.com/wanasit/chrono-go/internal/core"

// Casual is the Chrono object configured for parsing casual German.
var Casual = core.NewChrono(CreateCasualConfiguration(true))

// Strict is the Chrono object configured for parsing strict German.
var Strict = core.NewChrono(CreateConfiguration(true, true))

// Parse is Casual.Parse.
func Parse(text string, ref core.ParsingReference, option *core.ParsingOption) []*core.ParsingResult {
	return Casual.Parse(text, ref, option)
}

// ParseDate is Casual.ParseDate.
func ParseDate(text string, ref core.ParsingReference, option *core.ParsingOption) (float64, bool) {
	return Casual.ParseDate(text, ref, option)
}

// CreateCasualConfiguration creates a casual German configuration: the configuration of
// CreateConfiguration(false, littleEndian) with the casual parsers in front.
func CreateCasualConfiguration(littleEndian bool) *core.Configuration {
	option := CreateConfiguration(false, littleEndian)
	option.Parsers = append([]core.Parser{NewCasualTimeParser()}, option.Parsers...)
	option.Parsers = append([]core.Parser{NewCasualDateParser()}, option.Parsers...)
	option.Parsers = append([]core.Parser{NewTimeUnitRelativeFormatParser()}, option.Parsers...)
	return option
}

// CreateConfiguration creates a German configuration. strictMode is passed to the common
// configuration; littleEndian reads day-first slash dates.
func CreateConfiguration(strictMode, littleEndian bool) *core.Configuration {
	return core.IncludeCommonConfiguration(
		&core.Configuration{
			Parsers: []core.Parser{
				core.NewISOFormatParser(),
				core.NewSlashDateFormatParser(littleEndian),
				NewTimeExpressionParser(),
				NewSpecificTimeExpressionParser(),
				NewMonthNameLittleEndianParser(),
				NewWeekdayParser(),
				NewTimeUnitWithinFormatParser(),
			},
			Refiners: []core.Refiner{NewMergeDateRangeRefiner(), NewMergeDateTimeRefiner()},
		},
		strictMode,
	)
}
