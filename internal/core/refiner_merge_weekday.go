// Ported from src/common/refiners/MergeWeekdayComponentRefiner.ts.

package core

// MergeWeekdayComponentRefiner merges a weekday into a following, more complete date:
// [Sunday] [12/7/2014] => [Sunday 12/7/2014].
type MergeWeekdayComponentRefiner struct {
	MergingRefiner
}

// NewMergeWeekdayComponentRefiner returns a new MergeWeekdayComponentRefiner.
func NewMergeWeekdayComponentRefiner() *MergeWeekdayComponentRefiner {
	r := &MergeWeekdayComponentRefiner{}
	r.Self = r
	return r
}

var reCommaSpaces = MustRegExp("^,?\\s*$", "")

// MergeResults copies the weekday into the next result.
func (r *MergeWeekdayComponentRefiner) MergeResults(textBetween string, currentResult, nextResult *ParsingResult, context *ParsingContext) *ParsingResult {
	newResult := nextResult.Clone()
	newResult.Index = currentResult.Index
	newResult.Text = currentResult.Text + textBetween + newResult.Text

	newResult.Start.Assign(ComponentWeekday, currentResult.Start.Raw(ComponentWeekday))
	if newResult.End != nil {
		newResult.End.Assign(ComponentWeekday, currentResult.Start.Raw(ComponentWeekday))
	}
	return newResult
}

// ShouldMergeResults reports a weekday-only result followed by a date.
func (r *MergeWeekdayComponentRefiner) ShouldMergeResults(textBetween string, currentResult, nextResult *ParsingResult, context *ParsingContext) bool {
	weekdayThenNormalDate := currentResult.Start.IsOnlyWeekdayComponent() &&
		!currentResult.Start.IsCertain(ComponentHour) &&
		nextResult.Start.IsCertain(ComponentDay)
	return weekdayThenNormalDate && reCommaSpaces.Exec(textBetween) != nil
}
