// Ported from src/locales/zh/index.ts.

// Package zh is chrono's combined Chinese locale (Traditional and Simplified parsers side by
// side). The public package github.com/wanasit/chrono-go/locales/zh exposes it; the scripts
// alone are in internal/zh/hans and internal/zh/hant.
package zh

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/zh/hans"
	"github.com/wanasit/chrono-go/internal/zh/hant"
)

// Casual is the Chrono object configured for parsing casual Chinese.
var Casual = core.NewChrono(CreateCasualConfiguration())

// Strict is the Chrono object configured for parsing strict Chinese.
var Strict = core.NewChrono(CreateConfiguration())

// Parse is Casual.Parse.
func Parse(text string, ref core.ParsingReference, option *core.ParsingOption) []*core.ParsingResult {
	return Casual.Parse(text, ref, option)
}

// ParseDate is Casual.ParseDate.
func ParseDate(text string, ref core.ParsingReference, option *core.ParsingOption) (float64, bool) {
	return Casual.ParseDate(text, ref, option)
}

// CreateCasualConfiguration creates a casual configuration for Chinese: the configuration of
// CreateConfiguration with ZHHantCasualDateParser first.
func CreateCasualConfiguration() *core.Configuration {
	option := CreateConfiguration()
	option.Parsers = append([]core.Parser{hant.NewCasualDateParser()}, option.Parsers...)
	return option
}

// CreateConfiguration creates a configuration for Chinese.
func CreateConfiguration() *core.Configuration {
	configuration := core.IncludeCommonConfiguration(
		&core.Configuration{
			Parsers: []core.Parser{
				hant.NewDateParser(),
				hans.NewDateParser(),
				hant.NewRelationWeekdayParser(),
				hans.NewRelationWeekdayParser(),
				hant.NewWeekdayParser(),
				hans.NewWeekdayParser(),
				hant.NewTimeExpressionParser(),
				hans.NewTimeExpressionParser(),
				hant.NewDeadlineFormatParser(),
				hans.NewDeadlineFormatParser(),
			},
			Refiners: []core.Refiner{hant.NewMergeDateRangeRefiner(), hant.NewMergeDateTimeRefiner()},
		},
		false,
	)

	// REMOVE ExtractTimezoneOffsetRefiner
	refiners := []core.Refiner{}
	for _, refiner := range configuration.Refiners {
		if _, ok := refiner.(*core.ExtractTimezoneOffsetRefiner); !ok {
			refiners = append(refiners, refiner)
		}
	}
	configuration.Refiners = refiners

	return configuration
}
