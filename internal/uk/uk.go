// Ported from src/locales/uk/index.ts.

// Package uk is chrono's Ukrainian locale (parsers, refiners and configuration). The public
// package github.com/wanasit/chrono-go/locales/uk exposes it.
package uk

import "github.com/wanasit/chrono-go/internal/core"

// Casual is the Chrono object configured for parsing casual Ukrainian.
var Casual = core.NewChrono(CreateCasualConfiguration())

// Strict is the Chrono object configured for parsing strict Ukrainian.
var Strict = core.NewChrono(CreateConfiguration(true))

// CreateCasualConfiguration creates a default casual configuration for Ukrainian: the
// configuration of CreateConfiguration(false) plus the casual parsers.
func CreateCasualConfiguration() *core.Configuration {
	option := CreateConfiguration(false)
	option.Parsers = append([]core.Parser{NewCasualDateParser()}, option.Parsers...)
	option.Parsers = append([]core.Parser{NewCasualTimeParser()}, option.Parsers...)
	option.Parsers = append([]core.Parser{NewMonthNameParser()}, option.Parsers...)
	option.Parsers = append([]core.Parser{NewRelativeDateFormatParser()}, option.Parsers...)
	option.Parsers = append([]core.Parser{NewTimeUnitCasualRelativeFormatParser()}, option.Parsers...)
	return option
}

// CreateConfiguration creates a default configuration for Ukrainian; strictMode makes the time
// unit mentions strict, not casual.
func CreateConfiguration(strictMode bool) *core.Configuration {
	return core.IncludeCommonConfiguration(
		&core.Configuration{
			Parsers: []core.Parser{
				core.NewISOFormatParser(),
				core.NewSlashDateFormatParser(true),
				NewTimeUnitWithinFormatParser(),
				NewMonthNameLittleEndianParser(),
				NewWeekdayParser(),
				NewTimeExpressionParser(strictMode),
				NewTimeUnitAgoFormatParser(),
			},
			Refiners: []core.Refiner{NewMergeDateTimeRefiner(), NewMergeDateRangeRefiner()},
		},
		strictMode,
	)
}

// Parse is Casual.Parse.
func Parse(text string, ref core.ParsingReference, option *core.ParsingOption) []*core.ParsingResult {
	return Casual.Parse(text, ref, option)
}

// ParseDate is Casual.ParseDate.
func ParseDate(text string, ref core.ParsingReference, option *core.ParsingOption) (float64, bool) {
	return Casual.ParseDate(text, ref, option)
}
