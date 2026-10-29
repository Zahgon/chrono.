// Ported from src/locales/sv/index.ts.

// Package sv is chrono's Swedish locale (parsers and configuration). The public package
// github.com/wanasit/chrono-go/locales/sv exposes it.
package sv

import "github.com/wanasit/chrono-go/internal/core"

// Casual is the Chrono object configured for parsing casual Swedish.
var Casual = core.NewChrono(CreateCasualConfiguration(true))

// Strict is the Chrono object configured for parsing strict Swedish.
var Strict = core.NewChrono(CreateConfiguration(true, true))

// Parse is Casual.Parse.
func Parse(text string, ref core.ParsingReference, option *core.ParsingOption) []*core.ParsingResult {
	return Casual.Parse(text, ref, option)
}

// ParseDate is Casual.ParseDate.
func ParseDate(text string, ref core.ParsingReference, option *core.ParsingOption) (float64, bool) {
	return Casual.ParseDate(text, ref, option)
}

// CreateCasualConfiguration creates a casual Swedish configuration: the configuration of
// CreateConfiguration(false, littleEndian) with the casual date parser in front.
func CreateCasualConfiguration(littleEndian bool) *core.Configuration {
	option := CreateConfiguration(false, littleEndian)
	option.Parsers = append([]core.Parser{NewCasualDateParser()}, option.Parsers...)
	return option
}

// CreateConfiguration creates a Swedish configuration. strictMode is passed to the common
// configuration; littleEndian reads day-first slash dates.
func CreateConfiguration(strictMode, littleEndian bool) *core.Configuration {
	return core.IncludeCommonConfiguration(
		&core.Configuration{
			Parsers: []core.Parser{
				core.NewISOFormatParser(),
				core.NewSlashDateFormatParser(littleEndian),
				NewMonthNameLittleEndianParser(),
				NewWeekdayParser(),
				NewTimeUnitCasualRelativeFormatParser(true),
			},
			Refiners: []core.Refiner{},
		},
		strictMode,
	)
}
