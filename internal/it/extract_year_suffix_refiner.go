// Ported from src/locales/it/refiners/ITExtractYearSuffixRefiner.ts.

package it

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

var yearSuffixPattern = core.MustRegExp("^\\s*("+YearPattern+")", "i")

const yearSuffixYearGroup = 1

// ExtractYearSuffixRefiner is ITExtractYearSuffixRefiner: a year following a date without one.
type ExtractYearSuffixRefiner struct{}

// NewExtractYearSuffixRefiner returns a new ITExtractYearSuffixRefiner.
func NewExtractYearSuffixRefiner() *ExtractYearSuffixRefiner { return &ExtractYearSuffixRefiner{} }

// Refine assigns a following year to dates with an unknown year.
func (r *ExtractYearSuffixRefiner) Refine(context *core.ParsingContext, results []*core.ParsingResult) []*core.ParsingResult {
	for _, result := range results {
		if !result.Start.IsDateWithUnknownYear() {
			continue
		}
		suffix := context.SubstringFrom(result.Index + js.Len(result.Text))
		match := yearSuffixPattern.Exec(suffix)
		if match == nil {
			continue
		}
		// If the suffix match is just a short number, e.g. "14/4 90", don't assume it year.
		if js.Len(js.Trim(match.Get(0))) <= 3 {
			continue
		}
		year := ParseYear(match.Get(yearSuffixYearGroup))
		if result.End != nil {
			result.End.Assign(core.ComponentYear, year)
		}
		result.Start.Assign(core.ComponentYear, year)
		result.Text += match.Get(0)
	}
	return results
}
