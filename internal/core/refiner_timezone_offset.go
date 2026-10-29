// Ported from src/common/refiners/ExtractTimezoneOffsetRefiner.ts.

package core

import "github.com/wanasit/chrono-go/internal/js"

var timezoneOffsetPattern = MustRegExp("^\\s*(?:\\(?(?:GMT|UTC)\\s?)?([+-])(\\d{1,2})(?::?(\\d{2}))?\\)?", "i")

const (
	timezoneOffsetSignGroup         = 1
	timezoneOffsetHourOffsetGroup   = 2
	timezoneOffsetMinuteOffsetGroup = 3
)

// ExtractTimezoneOffsetRefiner extracts a "+hh:mm" / "GMT+hh" offset following a result.
type ExtractTimezoneOffsetRefiner struct{}

// NewExtractTimezoneOffsetRefiner returns a new ExtractTimezoneOffsetRefiner.
func NewExtractTimezoneOffsetRefiner() *ExtractTimezoneOffsetRefiner {
	return &ExtractTimezoneOffsetRefiner{}
}

// Refine assigns the offset following each result without a certain offset.
func (r *ExtractTimezoneOffsetRefiner) Refine(context *ParsingContext, results []*ParsingResult) []*ParsingResult {
	for _, result := range results {
		if result.Start.IsCertain(ComponentTimezoneOffset) {
			continue
		}

		suffix := context.SubstringFrom(result.Index + js.Len(result.Text))
		match := timezoneOffsetPattern.Exec(suffix)
		if match == nil {
			continue
		}

		hourOffset := js.ParseInt(match.Get(timezoneOffsetHourOffsetGroup), 0)
		minuteText := match.Get(timezoneOffsetMinuteOffsetGroup)
		if !match.Truthy(timezoneOffsetMinuteOffsetGroup) {
			minuteText = "0"
		}
		minuteOffset := js.ParseInt(minuteText, 0)
		timezoneOffset := hourOffset*60 + minuteOffset
		// No timezones have offsets greater than 14 hours, so disregard this match
		if timezoneOffset > 14*60 {
			continue
		}
		if match.Has(timezoneOffsetSignGroup) && match.Get(timezoneOffsetSignGroup) == "-" {
			timezoneOffset = -timezoneOffset
		}

		if result.End != nil {
			result.End.Assign(ComponentTimezoneOffset, timezoneOffset)
		}

		result.Start.Assign(ComponentTimezoneOffset, timezoneOffset)
		result.Text += match.Get(0)
	}
	return results
}
