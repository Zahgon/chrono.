// Ported from src/locales/zh/hans/index.ts.

// Package hans is chrono's Simplified Chinese locale (parsers, refiners and configuration). The
// public package github.com/wanasit/chrono-go/locales/zh/hans exposes it.
package hans

import "github.com/wanasit/chrono-go/internal/core"

// Hans is the `hans` Chrono object of the original module (casual configuration).
var Hans = core.NewChrono(CreateCasualConfiguration())

// Casual is the Chrono object configured for parsing casual Simplified Chinese.
var Casual = core.NewChrono(CreateCasualConfiguration())

// Strict is the Chrono object configured for parsing strict Simplified Chinese.
var Strict = core.NewChrono(CreateConfiguration())

// Parse is Casual.Parse.
func Parse(text string, ref core.ParsingReference, option *core.ParsingOption) []*core.ParsingResult {
	return Casual.Parse(text, ref, option)
}

// ParseDate is Casual.ParseDate.
func ParseDate(text string, ref core.ParsingReference, option *core.ParsingOption) (float64, bool) {
	return Casual.ParseDate(text, ref, option)
}

// CreateCasualConfiguration creates a casual configuration for Simplified Chinese: the
// configuration of CreateConfiguration with ZHHansCasualDateParser first.
func CreateCasualConfiguration() *core.Configuration {
	option := CreateConfiguration()
	option.Parsers = append([]core.Parser{NewCasualDateParser()}, option.Parsers...)
	return option
}

// CreateConfiguration creates a configuration for Simplified Chinese.
func CreateConfiguration() *core.Configuration {
	configuration := core.IncludeCommonConfiguration(
		&core.Configuration{
			Parsers: []core.Parser{
				NewDateParser(),
				NewRelationWeekdayParser(),
				NewWeekdayParser(),
				NewTimeExpressionParser(),
				NewDeadlineFormatParser(),
				NewAgoFormatParser(),
			},
			Refiners: []core.Refiner{NewMergeDateRangeRefiner(), NewMergeDateTimeRefiner()},
		},
		false,
	)

	// REMOVE ExtractTimezoneOffsetRefiner
	configuration.Refiners = withoutTimezoneOffsetRefiner(configuration.Refiners)

	return configuration
}

// withoutTimezoneOffsetRefiner is
// `refiners.filter((refiner) => !(refiner instanceof ExtractTimezoneOffsetRefiner))`.
func withoutTimezoneOffsetRefiner(refiners []core.Refiner) []core.Refiner {
	out := []core.Refiner{}
	for _, refiner := range refiners {
		if _, ok := refiner.(*core.ExtractTimezoneOffsetRefiner); !ok {
			out = append(out, refiner)
		}
	}
	return out
}
