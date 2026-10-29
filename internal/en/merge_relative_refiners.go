// Ported from src/locales/en/refiners/ENMergeRelativeAfterDateRefiner.ts and
// src/locales/en/refiners/ENMergeRelativeFollowByDateRefiner.ts.

package en

import "github.com/wanasit/chrono-go/internal/core"

var (
	rePositiveFollowing = core.MustRegExp("^[+-]", "i")
	reNegativeFollowing = core.MustRegExp("^-", "i")
	reBlankI            = core.MustRegExp("^\\s*$", "i")
	reEarlierReference  = core.MustRegExp("\\s+(before|from)$", "i")
	reLaterReference    = core.MustRegExp("\\s+(after|since)$", "i")
)

func isPositiveFollowingReference(result *core.ParsingResult) bool {
	return rePositiveFollowing.Exec(result.Text) != nil
}

func isNegativeFollowingReference(result *core.ParsingResult) bool {
	return reNegativeFollowing.Exec(result.Text) != nil
}

// MergeRelativeAfterDateRefiner is ENMergeRelativeAfterDateRefiner: a relative date/time after an
// absolute date ("[2020-02-13] [+2 weeks]").
type MergeRelativeAfterDateRefiner struct {
	core.MergingRefiner
}

// NewMergeRelativeAfterDateRefiner returns a new ENMergeRelativeAfterDateRefiner.
func NewMergeRelativeAfterDateRefiner() *MergeRelativeAfterDateRefiner {
	r := &MergeRelativeAfterDateRefiner{}
	r.Self = r
	return r
}

// ShouldMergeResults reports adjacent results where the second starts with + or -.
func (r *MergeRelativeAfterDateRefiner) ShouldMergeResults(textBetween string, currentResult, nextResult *core.ParsingResult, context *core.ParsingContext) bool {
	if reBlankI.Exec(textBetween) == nil {
		return false
	}
	return isPositiveFollowingReference(nextResult) || isNegativeFollowingReference(nextResult)
}

// MergeResults applies the relative duration to the absolute date.
func (r *MergeRelativeAfterDateRefiner) MergeResults(textBetween string, currentResult, nextResult *core.ParsingResult, context *core.ParsingContext) *core.ParsingResult {
	timeUnits := ParseDuration(nextResult.Text)
	if isNegativeFollowingReference(nextResult) {
		timeUnits = core.ReverseDuration(timeUnits)
	}

	components := core.CreateRelativeFromReference(core.ReferenceFromDate(currentResult.Start.Date()), timeUnits)

	return core.NewParsingResult(
		currentResult.Reference,
		currentResult.Index,
		currentResult.Text+textBetween+nextResult.Text,
		components,
		nil,
	)
}

func hasImpliedEarlierReferenceDate(result *core.ParsingResult) bool {
	return reEarlierReference.Exec(result.Text) != nil
}

func hasImpliedLaterReferenceDate(result *core.ParsingResult) bool {
	return reLaterReference.Exec(result.Text) != nil
}

// MergeRelativeFollowByDateRefiner is ENMergeRelativeFollowByDateRefiner: a relative date/time
// followed by an absolute date ("[2 weeks before] [2020-02-13]").
type MergeRelativeFollowByDateRefiner struct {
	core.MergingRefiner
}

// NewMergeRelativeFollowByDateRefiner returns a new ENMergeRelativeFollowByDateRefiner.
func NewMergeRelativeFollowByDateRefiner() *MergeRelativeFollowByDateRefiner {
	r := &MergeRelativeFollowByDateRefiner{}
	r.Self = r
	return r
}

// PatternBetween returns the blank pattern.
func (r *MergeRelativeFollowByDateRefiner) PatternBetween() *core.RegExp { return reBlankI }

// ShouldMergeResults reports "<relative> before|from|after|since" followed by an absolute date.
func (r *MergeRelativeFollowByDateRefiner) ShouldMergeResults(textBetween string, currentResult, nextResult *core.ParsingResult, context *core.ParsingContext) bool {
	// Dates need to be next to each other to get merged
	if r.PatternBetween().Exec(textBetween) == nil {
		return false
	}

	// Check if any relative tokens were swallowed by the first date.
	if !hasImpliedEarlierReferenceDate(currentResult) && !hasImpliedLaterReferenceDate(currentResult) {
		return false
	}

	// make sure that <date2> implies an absolute date
	return core.Truthy(nextResult.Start.V(core.ComponentDay)) && core.Truthy(nextResult.Start.V(core.ComponentMonth)) && core.Truthy(nextResult.Start.V(core.ComponentYear))
}

// MergeResults applies the relative duration to the following date.
func (r *MergeRelativeFollowByDateRefiner) MergeResults(textBetween string, currentResult, nextResult *core.ParsingResult, context *core.ParsingContext) *core.ParsingResult {
	duration := ParseDuration(currentResult.Text)
	if hasImpliedEarlierReferenceDate(currentResult) {
		duration = core.ReverseDuration(duration)
	}

	components := core.CreateRelativeFromReference(core.ReferenceFromDate(nextResult.Start.Date()), duration)

	return core.NewParsingResult(
		nextResult.Reference,
		currentResult.Index,
		currentResult.Text+textBetween+nextResult.Text,
		components,
		nil,
	)
}
