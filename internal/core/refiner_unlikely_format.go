// Ported from src/common/refiners/UnlikelyFormatFilter.ts.

package core

import "strings"

var reUnlikelyNumber = MustRegExp("^\\d*(\\.\\d*)?$", "")

// UnlikelyFormatFilter removes number-only, invalid and (in strict mode) weekday-only results.
type UnlikelyFormatFilter struct {
	Filter
	strictMode bool
}

// NewUnlikelyFormatFilter returns a new UnlikelyFormatFilter.
func NewUnlikelyFormatFilter(strictMode bool) *UnlikelyFormatFilter {
	f := &UnlikelyFormatFilter{strictMode: strictMode}
	f.Self = f
	return f
}

// IsValid reports whether the result is kept.
func (f *UnlikelyFormatFilter) IsValid(context *ParsingContext, result *ParsingResult) bool {
	if reUnlikelyNumber.Exec(strings.Replace(result.Text, " ", "", 1)) != nil {
		return false
	}

	if !result.Start.IsValidDate() {
		return false
	}

	if result.End != nil && !result.End.IsValidDate() {
		return false
	}

	if f.strictMode {
		return f.isStrictModeValid(context, result)
	}
	return true
}

func (f *UnlikelyFormatFilter) isStrictModeValid(context *ParsingContext, result *ParsingResult) bool {
	return !result.Start.IsOnlyWeekdayComponent()
}
