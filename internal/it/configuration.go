// Ported from src/locales/it/configuration.ts.

package it

import "github.com/wanasit/chrono-go/internal/core"

// CreateCasualConfiguration creates a default casual configuration for Italian
// (ITDefaultConfiguration.createCasualConfiguration): the configuration of
// CreateConfiguration(false) plus the casual parsers and the unlikely format filter.
func CreateCasualConfiguration() *core.Configuration {
	option := CreateConfiguration(false)
	option.Parsers = append(option.Parsers, NewCasualDateParser())
	option.Parsers = append(option.Parsers, NewCasualTimeParser())
	option.Parsers = append(option.Parsers, NewMonthNameParser())
	option.Parsers = append(option.Parsers, NewRelativeDateFormatParser())
	option.Parsers = append(option.Parsers, NewTimeUnitCasualRelativeFormatParser(true))
	option.Refiners = append(option.Refiners, NewUnlikelyFormatFilter())
	return option
}

// CreateConfiguration creates a default configuration for Italian
// (ITDefaultConfiguration.createConfiguration). strictMode makes the time unit mentions strict.
func CreateConfiguration(strictMode bool) *core.Configuration {
	options := core.IncludeCommonConfiguration(
		&core.Configuration{
			Parsers: []core.Parser{
				core.NewSlashDateFormatParser(true), // Italian uses day/month/year format (littleEndian)
				NewTimeUnitWithinFormatParser(strictMode),
				NewMonthNameLittleEndianParser(),
				NewMonthNameMiddleEndianParser( /*shouldSkipYearLikeDate=*/ true),
				NewWeekdayParser(),
				NewSlashMonthFormatParser(),
				NewTimeExpressionParser(strictMode),
				NewTimeUnitAgoFormatParser(strictMode),
				NewTimeUnitLaterFormatParser(strictMode),
			},
			Refiners: []core.Refiner{NewMergeDateTimeRefiner()},
		},
		strictMode,
	)
	options.Parsers = append([]core.Parser{NewYearMonthDayParser( /*strictMonthDateOrder=*/ strictMode)}, options.Parsers...)

	// These relative-dates consideration should be done before other common refiners.
	options.Refiners = append([]core.Refiner{NewMergeRelativeFollowByDateRefiner()}, options.Refiners...)
	options.Refiners = append([]core.Refiner{NewMergeRelativeAfterDateRefiner()}, options.Refiners...)
	options.Refiners = append([]core.Refiner{core.NewOverlapRemovalRefiner()}, options.Refiners...)

	// Re-apply the date time refiner again after the timezone refinement and exclusion in common
	// refiners.
	options.Refiners = append(options.Refiners, NewMergeDateTimeRefiner())

	// Extract year after merging date and time
	options.Refiners = append(options.Refiners, NewExtractYearSuffixRefiner())

	// Keep the date range refiner at the end (after all other refinements).
	options.Refiners = append(options.Refiners, NewMergeDateRangeRefiner())
	return options
}
