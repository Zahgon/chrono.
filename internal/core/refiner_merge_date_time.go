// Ported from src/common/refiners/AbstractMergeDateTimeRefiner.ts.

package core

// MergeDateTimeRefinerImpl holds the abstract method of AbstractMergeDateTimeRefiner.
type MergeDateTimeRefinerImpl interface {
	MergingRefinerImpl
	PatternBetween() *RegExp
}

// AbstractMergeDateTimeRefiner merges a date-only result and a time-only result when the text
// between them matches PatternBetween. Embed it, implement PatternBetween, and call Init with the
// embedding value in the constructor.
type AbstractMergeDateTimeRefiner struct {
	MergingRefiner
	DateTimeSelf MergeDateTimeRefinerImpl
}

// Init sets the virtual-dispatch receivers to self (the embedding value).
func (r *AbstractMergeDateTimeRefiner) Init(self MergeDateTimeRefinerImpl) {
	r.DateTimeSelf = self
	r.MergingRefiner.Self = self
}

// ShouldMergeResults reports a date-only / time-only pair with a matching text between.
func (r *AbstractMergeDateTimeRefiner) ShouldMergeResults(textBetween string, currentResult, nextResult *ParsingResult, context *ParsingContext) bool {
	return ((currentResult.Start.IsOnlyDate() && nextResult.Start.IsOnlyTime()) ||
		(nextResult.Start.IsOnlyDate() && currentResult.Start.IsOnlyTime())) &&
		r.DateTimeSelf.PatternBetween().Exec(textBetween) != nil
}

// MergeResults merges the date and the time.
func (r *AbstractMergeDateTimeRefiner) MergeResults(textBetween string, currentResult, nextResult *ParsingResult, context *ParsingContext) *ParsingResult {
	var result *ParsingResult
	if currentResult.Start.IsOnlyDate() {
		result = MergeDateTimeResult(currentResult, nextResult)
	} else {
		result = MergeDateTimeResult(nextResult, currentResult)
	}

	result.Index = currentResult.Index
	result.Text = currentResult.Text + textBetween + nextResult.Text
	return result
}
