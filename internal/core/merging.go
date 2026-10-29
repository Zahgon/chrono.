// Ported from src/calculation/mergingCalculation.ts.

package core

import "github.com/wanasit/chrono-go/internal/jsdate"

// MergeDateTimeResult merges a date-only result with a time-only result.
func MergeDateTimeResult(dateResult, timeResult *ParsingResult) *ParsingResult {
	result := dateResult.Clone()
	beginDate := dateResult.Start
	beginTime := timeResult.Start

	result.Start = MergeDateTimeComponent(beginDate, beginTime)
	if dateResult.End != nil || timeResult.End != nil {
		endDate := dateResult.End
		if endDate == nil {
			endDate = dateResult.Start
		}
		endTime := timeResult.End
		if endTime == nil {
			endTime = timeResult.Start
		}
		endDateTime := MergeDateTimeComponent(endDate, endTime)

		if dateResult.End == nil && endDateTime.Date() < result.Start.Date() {
			// For example, "Tuesday 9pm - 1am" the ending should actually be 1am on the next day.
			nextDay := jsdate.FromTime(endDateTime.Date())
			nextDay.SetDate(nextDay.GetDate() + 1)
			if endDateTime.IsCertain(ComponentDay) {
				AssignSimilarDate(endDateTime, nextDay)
			} else {
				ImplySimilarDate(endDateTime, nextDay)
			}
		}

		result.End = endDateTime
	}
	return result
}

// MergeDateTimeComponent merges the time of timeComponent into a clone of dateComponent.
func MergeDateTimeComponent(dateComponent, timeComponent *ParsingComponents) *ParsingComponents {
	dateTimeComponent := dateComponent.Clone()

	if timeComponent.IsCertain(ComponentHour) {
		dateTimeComponent.Assign(ComponentHour, timeComponent.Raw(ComponentHour))
		dateTimeComponent.Assign(ComponentMinute, timeComponent.Raw(ComponentMinute))

		if timeComponent.IsCertain(ComponentSecond) {
			dateTimeComponent.Assign(ComponentSecond, timeComponent.Raw(ComponentSecond))

			if timeComponent.IsCertain(ComponentMillisecond) {
				dateTimeComponent.Assign(ComponentMillisecond, timeComponent.Raw(ComponentMillisecond))
			} else {
				dateTimeComponent.Imply(ComponentMillisecond, timeComponent.Raw(ComponentMillisecond))
			}
		} else {
			dateTimeComponent.Imply(ComponentSecond, timeComponent.Raw(ComponentSecond))
			dateTimeComponent.Imply(ComponentMillisecond, timeComponent.Raw(ComponentMillisecond))
		}
	} else {
		dateTimeComponent.Imply(ComponentHour, timeComponent.Raw(ComponentHour))
		dateTimeComponent.Imply(ComponentMinute, timeComponent.Raw(ComponentMinute))
		dateTimeComponent.Imply(ComponentSecond, timeComponent.Raw(ComponentSecond))
		dateTimeComponent.Imply(ComponentMillisecond, timeComponent.Raw(ComponentMillisecond))
	}

	if timeComponent.IsCertain(ComponentTimezoneOffset) {
		dateTimeComponent.Assign(ComponentTimezoneOffset, timeComponent.Raw(ComponentTimezoneOffset))
	}

	_, dateMeridiemOK := dateComponent.Get(ComponentMeridiem)
	dateHasMeaningfulMeridiem := dateMeridiemOK &&
		(dateComponent.IsCertain(ComponentMeridiem) || dateComponent.HasTagPrefix("casualReference/"))

	if timeComponent.IsCertain(ComponentMeridiem) {
		dateTimeComponent.Assign(ComponentMeridiem, timeComponent.Raw(ComponentMeridiem))
	} else if _, ok := timeComponent.Get(ComponentMeridiem); ok && !dateHasMeaningfulMeridiem {
		// Let the time component guide meridiem only when the date part doesn't already carry a
		// meaningful (certain or casual) hint.
		dateTimeComponent.Imply(ComponentMeridiem, timeComponent.Raw(ComponentMeridiem))
	}

	// `get("meridiem") == Meridiem.PM && get("hour") < 12` (null -> 0 < 12; undefined -> false)
	if meridiem, ok := dateTimeComponent.Get(ComponentMeridiem); ok && meridiem == PM && dateTimeComponent.V(ComponentHour) < 12 {
		if timeComponent.IsCertain(ComponentHour) {
			dateTimeComponent.Assign(ComponentHour, dateTimeComponent.V(ComponentHour)+12)
		} else {
			dateTimeComponent.Imply(ComponentHour, dateTimeComponent.V(ComponentHour)+12)
		}
	}

	dateTimeComponent.AddTags(dateComponent.Tags())
	dateTimeComponent.AddTags(timeComponent.Tags())
	return dateTimeComponent
}
