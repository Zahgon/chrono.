// Ported from src/common/refiners/ForwardDateRefiner.ts.

package core

import "github.com/wanasit/chrono-go/internal/jsdate"

// ForwardDateRefiner enforces the forwardDate option: results with missing components are moved
// into the future instead of the past.
type ForwardDateRefiner struct{}

// NewForwardDateRefiner returns a new ForwardDateRefiner.
func NewForwardDateRefiner() *ForwardDateRefiner { return &ForwardDateRefiner{} }

// Refine moves the results after the reference date when the forwardDate option is set.
func (r *ForwardDateRefiner) Refine(context *ParsingContext, results []*ParsingResult) []*ParsingResult {
	if !context.Option.ForwardDate {
		return results
	}

	for _, result := range results {
		refDate := context.Reference.DateWithAdjustedTimezone()

		if result.Start.IsOnlyTime() && context.Reference.Instant > result.Start.Date() {
			refDate := context.Reference.DateWithAdjustedTimezone()
			refFollowingDay := refDate.Clone()
			refFollowingDay.SetDate(refFollowingDay.GetDate() + 1)

			ImplySimilarDate(result.Start, refFollowingDay)
			if result.End != nil && result.End.IsOnlyTime() {
				ImplySimilarDate(result.End, refFollowingDay)
				if result.Start.Date() > result.End.Date() {
					refFollowingDay.SetDate(refFollowingDay.GetDate() + 1)
					ImplySimilarDate(result.End, refFollowingDay)
				}
			}
		}

		if result.Start.IsOnlyWeekdayComponent() && refDate.Time() > result.Start.Date() {
			daysToAdd := forwardDays(refDate, result.Start.V(ComponentWeekday))
			forwardedWeekday := AddDuration(refDate, NewDuration(DurationValue{"day", daysToAdd}))
			ImplySimilarDate(result.Start, forwardedWeekday)

			if result.End != nil && result.Start.Date() > result.End.Date() {
				daysToAdd := forwardDays(refDate, result.Start.V(ComponentWeekday))
				forwardedWeekday := AddDuration(refDate, NewDuration(DurationValue{"day", daysToAdd}))
				ImplySimilarDate(result.End, forwardedWeekday)
			}
		}

		// In case where we know the month, but not which year (e.g. "in December", "25th
		// December"), try move to another year (up-to 3 times)
		if result.Start.IsDateWithUnknownYear() && refDate.Time() > result.Start.Date() {
			for i := 0; i < 3 && refDate.Time() > result.Start.Date(); i++ {
				result.Start.Imply(ComponentYear, result.Start.V(ComponentYear)+1)

				if result.End != nil && !result.End.IsCertain(ComponentYear) {
					result.End.Imply(ComponentYear, result.End.V(ComponentYear)+1)
				}
			}
		}
	}
	return results
}

// forwardDays is `getDaysForwardToWeekday(refDate, weekday) || 7`.
func forwardDays(refDate *jsdate.Date, weekday float64) float64 {
	d := GetDaysForwardToWeekday(refDate, weekday)
	if !Truthy(d) {
		return 7
	}
	return d
}
