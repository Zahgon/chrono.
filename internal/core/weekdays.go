// Ported from src/calculation/weekdays.ts.

package core

import "github.com/wanasit/chrono-go/internal/jsdate"

// CreateParsingComponentsAtWeekday returns the parsing components at the weekday considering the
// modifier ("this", "next", "last", or "" for none: the weekday closest to the reference). The
// time and time zone are assumed to be similar to the reference.
func CreateParsingComponentsAtWeekday(reference *ReferenceWithTimezone, weekday float64, modifier string) *ParsingComponents {
	refDate := reference.DateWithAdjustedTimezone()
	daysToWeekday := GetDaysToWeekday(refDate, weekday, modifier)

	components := NewParsingComponents(reference, nil)
	components = components.AddDurationAsImplied(NewDuration(DurationValue{"day", daysToWeekday}))
	components.Assign(ComponentWeekday, weekday)
	return components
}

// GetDaysToWeekday returns the number of days from refDate to the weekday ("this", "next", "last"
// or "" for the closest).
func GetDaysToWeekday(refDate *jsdate.Date, weekday float64, modifier string) float64 {
	refWeekday := refDate.GetDay()
	switch modifier {
	case "this":
		return GetDaysForwardToWeekday(refDate, weekday)
	case "last":
		return GetBackwardDaysToWeekday(refDate, weekday)
	case "next":
		// From Sunday, the next Sunday is 7 days later.
		if refWeekday == Sunday {
			if weekday == Sunday {
				return 7
			}
			return weekday
		}
		// From Saturday, the next Saturday is 7 days later, the next Sunday is 8-days later.
		if refWeekday == Saturday {
			if weekday == Saturday {
				return 7
			}
			if weekday == Sunday {
				return 8
			}
			return 1 + weekday
		}
		// From weekdays, next Mon is the following week's Mon, and so on.
		if weekday < refWeekday && weekday != Sunday {
			return GetDaysForwardToWeekday(refDate, weekday)
		}
		return GetDaysForwardToWeekday(refDate, weekday) + 7
	}
	return GetDaysToWeekdayClosest(refDate, weekday)
}

// GetDaysToWeekdayClosest returns the days to the closest such weekday (forward or backward).
func GetDaysToWeekdayClosest(refDate *jsdate.Date, weekday float64) float64 {
	backward := GetBackwardDaysToWeekday(refDate, weekday)
	forward := GetDaysForwardToWeekday(refDate, weekday)
	if forward < -backward {
		return forward
	}
	return backward
}

// GetDaysForwardToWeekday returns the days forward (0..6) to the weekday.
func GetDaysForwardToWeekday(refDate *jsdate.Date, weekday float64) float64 {
	refWeekday := refDate.GetDay()
	forwardCount := weekday - refWeekday
	if forwardCount < 0 {
		forwardCount += 7
	}
	return forwardCount
}

// GetBackwardDaysToWeekday returns the days backward (-7..-1) to the weekday.
func GetBackwardDaysToWeekday(refDate *jsdate.Date, weekday float64) float64 {
	refWeekday := refDate.GetDay()
	backwardCount := weekday - refWeekday
	if backwardCount >= 0 {
		backwardCount -= 7
	}
	return backwardCount
}
