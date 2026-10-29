// Ported from src/common/abstractRefiners.ts.

package core

// FilterImpl is the abstract method of a Filter.
type FilterImpl interface {
	IsValid(context *ParsingContext, result *ParsingResult) bool
}

// Filter is a Refiner that keeps the results for which Self.IsValid is true. Embed it and set
// Self to the embedding value.
type Filter struct {
	Self FilterImpl
}

// Refine is `results.filter((r) => this.isValid(context, r))`.
func (f *Filter) Refine(context *ParsingContext, results []*ParsingResult) []*ParsingResult {
	out := []*ParsingResult{}
	for _, r := range results {
		if f.Self.IsValid(context, r) {
			out = append(out, r)
		}
	}
	return out
}

// MergingRefinerImpl holds the abstract methods of a MergingRefiner.
type MergingRefinerImpl interface {
	ShouldMergeResults(textBetween string, currentResult, nextResult *ParsingResult, context *ParsingContext) bool
	MergeResults(textBetween string, currentResult, nextResult *ParsingResult, context *ParsingContext) *ParsingResult
}

// MergingRefiner is a Refiner that merges consecutive results. Embed it and set Self to the
// embedding value.
type MergingRefiner struct {
	Self MergingRefinerImpl
}

// Refine merges each pair of consecutive results for which Self.ShouldMergeResults is true.
func (m *MergingRefiner) Refine(context *ParsingContext, results []*ParsingResult) []*ParsingResult {
	if len(results) < 2 {
		return results
	}

	mergedResults := []*ParsingResult{}
	curResult := results[0]
	var nextResult *ParsingResult

	for i := 1; i < len(results); i++ {
		nextResult = results[i]

		textBetween := context.Substring(curResult.Index+jsLen(curResult.Text), nextResult.Index)
		if !m.Self.ShouldMergeResults(textBetween, curResult, nextResult, context) {
			mergedResults = append(mergedResults, curResult)
			curResult = nextResult
		} else {
			left := curResult
			right := nextResult
			mergedResult := m.Self.MergeResults(textBetween, left, right, context)
			curResult = mergedResult
		}
	}

	if curResult != nil {
		mergedResults = append(mergedResults, curResult)
	}
	return mergedResults
}
