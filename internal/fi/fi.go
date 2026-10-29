// Ported from src/locales/fi/index.ts.

// Package fi is chrono's Finnish locale (parsers, refiners and configuration). The public
// package github.com/wanasit/chrono-go/locales/fi exposes it.
package fi

import "github.com/wanasit/chrono-go/internal/core"

// Casual is the Chrono object configured for parsing casual Finnish.
var Casual = core.NewChrono(CreateCasualConfiguration(true))

// Strict is the Chrono object configured for parsing strict Finnish.
var Strict = core.NewChrono(CreateConfiguration(true, true))

// Parse is Casual.Parse.
func Parse(text string, ref core.ParsingReference, option *core.ParsingOption) []*core.ParsingResult {
	return Casual.Parse(text, ref, option)
}

// ParseDate is Casual.ParseDate.
func ParseDate(text string, ref core.ParsingReference, option *core.ParsingOption) (float64, bool) {
	return Casual.ParseDate(text, ref, option)
}

// CreateCasualConfiguration creates a casual Finnish configuration: the configuration of
// CreateConfiguration(false, littleEndian) with the casual parsers in front. The original's
// default littleEndian is true.
func CreateCasualConfiguration(littleEndian bool) *core.Configuration {
	option := CreateConfiguration(false, littleEndian)
	option.Parsers = append([]core.Parser{NewCasualTimeParser()}, option.Parsers...)
	option.Parsers = append([]core.Parser{NewCasualDateParser()}, option.Parsers...)
	option.Parsers = append([]core.Parser{NewTimeUnitCasualRelativeFormatParser(true)}, option.Parsers...)
	return option
}

// CreateConfiguration creates a Finnish configuration. The original's defaults are
// strictMode = true and littleEndian = true.
func CreateConfiguration(strictMode, littleEndian bool) *core.Configuration {
	return core.IncludeCommonConfiguration(
		&core.Configuration{
			Parsers: []core.Parser{
				core.NewISOFormatParser(),
				core.NewSlashDateFormatParser(littleEndian),
				NewTimeExpressionParser(),
				NewMonthNameLittleEndianParser(),
				NewWeekdayParser(),
				NewTimeUnitWithinFormatParser(),
				NewTimeUnitAgoFormatParser(),
			},
			Refiners: []core.Refiner{NewMergeDateRangeRefiner(), NewMergeDateTimeRefiner()},
		},
		strictMode,
	)
}
