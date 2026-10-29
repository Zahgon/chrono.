// Ported from src/locales/it/refiners/ITUnlikelyFormatFilter.ts.

package it

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

// UnlikelyFormatFilter is ITUnlikelyFormatFilter: drops a "secondo" (ordinal or time unit) that
// is not the whole text.
type UnlikelyFormatFilter struct {
	core.Filter
}

// NewUnlikelyFormatFilter returns a new ITUnlikelyFormatFilter.
func NewUnlikelyFormatFilter() *UnlikelyFormatFilter {
	f := &UnlikelyFormatFilter{}
	f.Self = f
	return f
}

// IsValid reports whether the result is kept.
func (f *UnlikelyFormatFilter) IsValid(context *core.ParsingContext, result *core.ParsingResult) bool {
	text := js.Trim(result.Text)

	// If the result is consists of the whole text (e.g. "2024", "Maggio", etc),
	// then it is unlikely to be a date.
	if text == js.Trim(context.Text) {
		return true
	}

	// In Italian, "secondo" could refer to the ordinal number or timeunit.
	if js.EndsWith(js.ToLowerCase(text), "il secondo") || js.ToLowerCase(text) == "secondo" {
		// (The original computes the trimmed text after the result only for a debug message.)
		return false
	}

	return true
}
