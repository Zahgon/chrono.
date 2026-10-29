// Ported from src/locales/ja/refiners/JPMergeDateRangeRefiner.ts,
// src/locales/ja/refiners/JPMergeDateTimeRefiner.ts and
// src/locales/ja/refiners/JPMergeWeekdayComponentRefiner.ts.

package ja

import "github.com/wanasit/chrono-go/internal/core"

var mergeDateRangePattern = core.MustRegExp("^\\s*(から|－|ー|-|～|~)\\s*$", "i")

// MergeDateRangeRefiner is JPMergeDateRangeRefiner: "2月11日[ー]2月13日", "水曜日[から]日曜日".
type MergeDateRangeRefiner struct {
	core.AbstractMergeDateRangeRefiner
}

// NewMergeDateRangeRefiner returns a new JPMergeDateRangeRefiner.
func NewMergeDateRangeRefiner() *MergeDateRangeRefiner {
	r := &MergeDateRangeRefiner{}
	r.Init(r)
	return r
}

// PatternBetween returns the range connector pattern.
func (r *MergeDateRangeRefiner) PatternBetween() *core.RegExp { return mergeDateRangePattern }

var mergeDateTimePattern = core.MustRegExp("^\\s*(の)?\\s*$", "i")

// MergeDateTimeRefiner is JPMergeDateTimeRefiner: "7月12日[の]午後3時".
type MergeDateTimeRefiner struct {
	core.AbstractMergeDateTimeRefiner
}

// NewMergeDateTimeRefiner returns a new JPMergeDateTimeRefiner.
func NewMergeDateTimeRefiner() *MergeDateTimeRefiner {
	r := &MergeDateTimeRefiner{}
	r.Init(r)
	return r
}

// PatternBetween returns the date/time connector pattern.
func (r *MergeDateTimeRefiner) PatternBetween() *core.RegExp { return mergeDateTimePattern }

var mergeWeekdayBetweenPattern = core.MustRegExp("^[,、の]?\\s*$", "")

// MergeWeekdayComponentRefiner is JPMergeWeekdayComponentRefiner: it merges a weekday into the
// preceding, more complete date ([2014/7/12], [(土)] => [2014/7/12 (土)]; [7月12日], [土曜日] =>
// [7月12日 土曜日]). Unlike the common MergeWeekdayComponentRefiner, the date comes first.
type MergeWeekdayComponentRefiner struct {
	core.MergingRefiner
}

// NewMergeWeekdayComponentRefiner returns a new JPMergeWeekdayComponentRefiner.
func NewMergeWeekdayComponentRefiner() *MergeWeekdayComponentRefiner {
	r := &MergeWeekdayComponentRefiner{}
	r.Self = r
	return r
}

// MergeResults copies the weekday of the next result into the current one.
func (r *MergeWeekdayComponentRefiner) MergeResults(textBetween string, currentResult, nextResult *core.ParsingResult, context *core.ParsingContext) *core.ParsingResult {
	newResult := currentResult.Clone()
	newResult.Text = currentResult.Text + textBetween + nextResult.Text

	newResult.Start.Assign(core.ComponentWeekday, nextResult.Start.Raw(core.ComponentWeekday))
	if newResult.End != nil {
		newResult.End.Assign(core.ComponentWeekday, nextResult.Start.Raw(core.ComponentWeekday))
	}

	return newResult
}

// ShouldMergeResults reports a date followed by a weekday-only result.
func (r *MergeWeekdayComponentRefiner) ShouldMergeResults(textBetween string, currentResult, nextResult *core.ParsingResult, context *core.ParsingContext) bool {
	normalDateThenWeekday := currentResult.Start.IsCertain(core.ComponentDay) &&
		nextResult.Start.IsOnlyWeekdayComponent() &&
		!nextResult.Start.IsCertain(core.ComponentHour)
	return normalDateThenWeekday && mergeWeekdayBetweenPattern.Exec(textBetween) != nil
}
