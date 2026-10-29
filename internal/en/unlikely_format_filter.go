// Ported from src/locales/en/refiners/ENUnlikelyFormatFilter.ts.

package en

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

var reInBefore = core.MustRegExp("\\b(in)$", "i")

// UnlikelyFormatFilter is ENUnlikelyFormatFilter: drops "may" (the modal verb) and a dangling
// "the second".
type UnlikelyFormatFilter struct {
	core.Filter
}

// NewUnlikelyFormatFilter returns a new ENUnlikelyFormatFilter.
func NewUnlikelyFormatFilter() *UnlikelyFormatFilter {
	f := &UnlikelyFormatFilter{}
	f.Self = f
	return f
}

// IsValid reports whether the result is kept.
func (f *UnlikelyFormatFilter) IsValid(context *core.ParsingContext, result *core.ParsingResult) bool {
	text := js.Trim(result.Text)

	// If the result is consists of the whole text (e.g. "2024", "May", etc), then it is unlikely
	// to be a date.
	if text == js.Trim(context.Text) {
		return true
	}

	// In English, the word "may" is a month name, but it is also a modal verb.
	if js.ToLowerCase(text) == "may" {
		textBefore := js.Trim(context.Substring(0, result.Index))
		if reInBefore.Exec(textBefore) == nil {
			return false
		}
	}

	// In English, "the second" could refer to the ordinal number or timeunit.
	if js.EndsWith(js.ToLowerCase(text), "the second") {
		return false
	}
	return true
}
