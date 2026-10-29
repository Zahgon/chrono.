// Ported from src/common/refiners/OverlapRemovalRefiner.ts.

package core

import "github.com/wanasit/chrono-go/internal/js"

// OverlapRemovalRefiner removes the shorter of two overlapping results.
type OverlapRemovalRefiner struct{}

// NewOverlapRemovalRefiner returns a new OverlapRemovalRefiner.
func NewOverlapRemovalRefiner() *OverlapRemovalRefiner { return &OverlapRemovalRefiner{} }

// Refine keeps the longer result of each overlapping pair.
func (r *OverlapRemovalRefiner) Refine(context *ParsingContext, results []*ParsingResult) []*ParsingResult {
	if len(results) < 2 {
		return results
	}

	filteredResults := []*ParsingResult{}
	prevResult := results[0]
	for i := 1; i < len(results); i++ {
		result := results[i]
		if result.Index >= prevResult.Index+js.Len(prevResult.Text) {
			filteredResults = append(filteredResults, prevResult)
			prevResult = result
			continue
		}

		// If overlap, compare the length and discard the shorter one
		if js.Len(result.Text) > js.Len(prevResult.Text) {
			prevResult = result
		}
	}

	// The last one
	if prevResult != nil {
		filteredResults = append(filteredResults, prevResult)
	}
	return filteredResults
}
