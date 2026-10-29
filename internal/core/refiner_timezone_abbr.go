// Ported from src/common/refiners/ExtractTimezoneAbbrRefiner.ts.

package core

import "github.com/wanasit/chrono-go/internal/js"

var timezoneNamePattern = MustRegExp("^\\s*,?\\s*\\(?([A-Z]{2,4})\\)?(?=\\W|$)", "i")

// ExtractTimezoneAbbrRefiner extracts a time zone abbreviation following a result.
type ExtractTimezoneAbbrRefiner struct {
	timezoneOverrides map[string]TimezoneAbbr
}

// NewExtractTimezoneAbbrRefiner returns the refiner with optional constructor overrides (nil:
// undefined).
func NewExtractTimezoneAbbrRefiner(timezoneOverrides map[string]TimezoneAbbr) *ExtractTimezoneAbbrRefiner {
	return &ExtractTimezoneAbbrRefiner{timezoneOverrides: timezoneOverrides}
}

// Refine assigns the offset of an abbreviation following each result.
func (r *ExtractTimezoneAbbrRefiner) Refine(context *ParsingContext, results []*ParsingResult) []*ParsingResult {
	timezoneOverrides := context.Option.Timezones

	for _, result := range results {
		suffix := context.SubstringFrom(result.Index + js.Len(result.Text))
		match := timezoneNamePattern.Exec(suffix)
		if match == nil {
			continue
		}

		timezoneAbbr := js.ToUpperCase(match.Get(1))
		refDate := result.Start.DateObject()
		tzOverrides := map[string]TimezoneAbbr{}
		for k, v := range r.timezoneOverrides {
			tzOverrides[k] = v
		}
		for k, v := range timezoneOverrides {
			tzOverrides[k] = v
		}
		extractedTimezoneOffset, ok := ToTimezoneOffset(timezoneAbbr, refDate, tzOverrides)
		if !ok {
			continue
		}

		currentTimezoneOffset := result.Start.Raw(ComponentTimezoneOffset)
		// `currentTimezoneOffset !== null && extractedTimezoneOffset != currentTimezoneOffset`
		if !IsNullMarker(currentTimezoneOffset) && !(extractedTimezoneOffset == Num(currentTimezoneOffset) && !IsUndefined(currentTimezoneOffset)) {
			// We may already have extracted the timezone offset e.g. "11 am GMT+0900 (JST)"
			// - if they are equal, we also want to take the abbreviation text into result
			// - if they are not equal, we trust the offset more
			if result.Start.IsCertain(ComponentTimezoneOffset) {
				continue
			}

			// This is often because it's relative time with inferred timezone (e.g. in 1 hour,
			// tomorrow). Then, we want to double-check the abbr case (e.g. "GET" not "get")
			if timezoneAbbr != match.Get(1) {
				continue
			}
		}

		if result.Start.IsOnlyDate() {
			// If the time is not explicitly mentioned, we also want to double-check the abbr case
			if timezoneAbbr != match.Get(1) {
				continue
			}
		}

		result.Text += match.Get(0)

		if !result.Start.IsCertain(ComponentTimezoneOffset) {
			result.Start.Assign(ComponentTimezoneOffset, extractedTimezoneOffset)
		}

		if result.End != nil && !result.End.IsCertain(ComponentTimezoneOffset) {
			result.End.Assign(ComponentTimezoneOffset, extractedTimezoneOffset)
		}
	}
	return results
}
