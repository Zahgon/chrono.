// Ported from src/common/refiners/AbstractMergeDateRangeRefiner.ts.

package core

// MergeDateRangeRefinerImpl holds the abstract method of AbstractMergeDateRangeRefiner.
type MergeDateRangeRefinerImpl interface {
	MergingRefinerImpl
	PatternBetween() *RegExp
}

// AbstractMergeDateRangeRefiner merges two consecutive results into a range when the text between
// them matches PatternBetween. Embed it, implement PatternBetween, and call Init with the
// embedding value in the constructor.
type AbstractMergeDateRangeRefiner struct {
	MergingRefiner
	RangeSelf MergeDateRangeRefinerImpl
}

// Init sets the virtual-dispatch receivers to self (the embedding value).
func (r *AbstractMergeDateRangeRefiner) Init(self MergeDateRangeRefinerImpl) {
	r.RangeSelf = self
	r.MergingRefiner.Self = self
}

// ShouldMergeResults reports that neither result is a range and the text between matches.
func (r *AbstractMergeDateRangeRefiner) ShouldMergeResults(textBetween string, currentResult, nextResult *ParsingResult, context *ParsingContext) bool {
	return currentResult.End == nil && nextResult.End == nil && r.RangeSelf.PatternBetween().Exec(textBetween) != nil
}

// MergeResults merges the two results into a range.
func (r *AbstractMergeDateRangeRefiner) MergeResults(textBetween string, fromResult, toResult *ParsingResult, context *ParsingContext) *ParsingResult {
	if !fromResult.Start.IsOnlyWeekdayComponent() && !toResult.Start.IsOnlyWeekdayComponent() {
		for _, key := range toResult.Start.GetCertainComponents() {
			if !fromResult.Start.IsCertain(key) {
				fromResult.Start.Imply(key, toResult.Start.Raw(key))
			}
		}

		for _, key := range fromResult.Start.GetCertainComponents() {
			if !toResult.Start.IsCertain(key) {
				toResult.Start.Imply(key, fromResult.Start.Raw(key))
			}
		}
	}
	if fromResult.Start.Date() > toResult.Start.Date() {
		fromDate := fromResult.Start.DateObject()
		toDate := toResult.Start.DateObject()

		if toResult.Start.IsOnlyWeekdayComponent() && AddDuration(toDate, NewDuration(DurationValue{"day", 7})).Time() > fromDate.Time() {
			toDate = AddDuration(toDate, NewDuration(DurationValue{"day", 7}))
			toResult.Start.Imply(ComponentDay, toDate.GetDate())
			toResult.Start.Imply(ComponentMonth, toDate.GetMonth()+1)
			toResult.Start.Imply(ComponentYear, toDate.GetFullYear())
		} else if fromResult.Start.IsOnlyWeekdayComponent() && AddDuration(fromDate, NewDuration(DurationValue{"day", -7})).Time() < toDate.Time() {
			fromDate = AddDuration(fromDate, NewDuration(DurationValue{"day", -7}))
			fromResult.Start.Imply(ComponentDay, fromDate.GetDate())
			fromResult.Start.Imply(ComponentMonth, fromDate.GetMonth()+1)
			fromResult.Start.Imply(ComponentYear, fromDate.GetFullYear())
		} else if toResult.Start.IsDateWithUnknownYear() && AddDuration(toDate, NewDuration(DurationValue{"year", 1})).Time() > fromDate.Time() {
			toDate = AddDuration(toDate, NewDuration(DurationValue{"year", 1}))
			toResult.Start.Imply(ComponentYear, toDate.GetFullYear())
		} else if fromResult.Start.IsDateWithUnknownYear() && AddDuration(fromDate, NewDuration(DurationValue{"year", -1})).Time() < toDate.Time() {
			fromDate = AddDuration(fromDate, NewDuration(DurationValue{"year", -1}))
			fromResult.Start.Imply(ComponentYear, fromDate.GetFullYear())
		} else {
			toResult, fromResult = fromResult, toResult
		}
	}
	result := fromResult.Clone()
	result.Start = fromResult.Start
	result.End = toResult.Start
	result.Index = min(fromResult.Index, toResult.Index)
	if fromResult.Index < toResult.Index {
		result.Text = fromResult.Text + textBetween + toResult.Text
	} else {
		result.Text = toResult.Text + textBetween + fromResult.Text
	}
	return result
}
