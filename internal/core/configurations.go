// Ported from src/configurations.ts.

package core

// IncludeCommonConfiguration adds the common parsers and refiners to the configuration (in
// place) and returns it.
func IncludeCommonConfiguration(configuration *Configuration, strictMode bool) *Configuration {
	configuration.Parsers = append([]Parser{NewISOFormatParser()}, configuration.Parsers...)

	configuration.Refiners = append([]Refiner{NewMergeWeekdayComponentRefiner()}, configuration.Refiners...)
	configuration.Refiners = append([]Refiner{NewExtractTimezoneOffsetRefiner()}, configuration.Refiners...)
	configuration.Refiners = append([]Refiner{NewOverlapRemovalRefiner()}, configuration.Refiners...)

	// Unlike ExtractTimezoneOffsetRefiner, this refiner relies on knowing both date and time in
	// cases where the tz is ambiguous (in terms of DST/non-DST). It therefore needs to be applied
	// as late as possible in the parsing.
	configuration.Refiners = append(configuration.Refiners, NewExtractTimezoneAbbrRefiner(nil))
	configuration.Refiners = append(configuration.Refiners, NewOverlapRemovalRefiner())
	configuration.Refiners = append(configuration.Refiners, NewForwardDateRefiner())
	configuration.Refiners = append(configuration.Refiners, NewUnlikelyFormatFilter(strictMode))
	return configuration
}
