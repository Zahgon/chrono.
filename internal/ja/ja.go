// Ported from src/locales/ja/index.ts.

// Package ja is chrono's Japanese locale (parsers, refiners and configuration). The public
// package github.com/wanasit/chrono-go/locales/ja exposes it.
package ja

import "github.com/wanasit/chrono-go/internal/core"

// Casual is the Chrono object configured for parsing casual Japanese.
var Casual = core.NewChrono(CreateCasualConfiguration())

// Strict is the Chrono object configured for parsing strict Japanese.
var Strict = core.NewChrono(CreateConfiguration(true))

// Parse is Casual.Parse.
func Parse(text string, ref core.ParsingReference, option *core.ParsingOption) []*core.ParsingResult {
	return Casual.Parse(text, ref, option)
}

// ParseDate is Casual.ParseDate.
func ParseDate(text string, ref core.ParsingReference, option *core.ParsingOption) (float64, bool) {
	return Casual.ParseDate(text, ref, option)
}

// CreateCasualConfiguration creates a casual configuration for Japanese: the configuration of
// CreateConfiguration(false) with JPCasualDateParser first.
func CreateCasualConfiguration() *core.Configuration {
	option := CreateConfiguration(false)
	option.Parsers = append([]core.Parser{NewCasualDateParser()}, option.Parsers...)
	return option
}

// CreateConfiguration creates a configuration for Japanese (the JS default of strictMode is
// true).
func CreateConfiguration(strictMode bool) *core.Configuration {
	configuration := core.IncludeCommonConfiguration(
		&core.Configuration{
			Parsers: []core.Parser{
				NewStandardParser(),
				NewWeekdayParser(),
				NewWeekdayWithParenthesesParser(),
				NewSlashDateFormatParser(),
				NewTimeExpressionParser(),
			},
			Refiners: []core.Refiner{
				NewMergeWeekdayComponentRefiner(), // Note: should be before JPMergeDateTimeRefiner and JPMergeDateRangeRefiner
				NewMergeDateTimeRefiner(),
				NewMergeDateRangeRefiner(),
			},
		},
		strictMode,
	)

	// Note: Remove because it is not used in Japanese grammar
	refiners := []core.Refiner{}
	for _, refiner := range configuration.Refiners {
		if _, ok := refiner.(*core.MergeWeekdayComponentRefiner); !ok {
			refiners = append(refiners, refiner)
		}
	}
	configuration.Refiners = refiners

	return configuration
}
