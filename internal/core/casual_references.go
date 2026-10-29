// Ported from src/common/casualReferences.ts.

package core

import "github.com/wanasit/chrono-go/internal/jsdate"

// Now returns the components of "now" at the reference.
func Now(reference *ReferenceWithTimezone) *ParsingComponents {
	targetDate := reference.DateWithAdjustedTimezone()
	component := NewParsingComponents(reference, Components{})
	AssignSimilarDate(component, targetDate)
	AssignSimilarTime(component, targetDate)
	component.Assign(ComponentTimezoneOffset, reference.GetTimezoneOffset())
	component.AddTag("casualReference/now")
	return component
}

// Today returns the components of "today".
func Today(reference *ReferenceWithTimezone) *ParsingComponents {
	targetDate := reference.DateWithAdjustedTimezone()
	component := NewParsingComponents(reference, Components{})
	AssignSimilarDate(component, targetDate)
	ImplySimilarTime(component, targetDate)
	component.Delete(ComponentMeridiem)
	component.AddTag("casualReference/today")
	return component
}

// Yesterday returns the components of "yesterday".
func Yesterday(reference *ReferenceWithTimezone) *ParsingComponents {
	return TheDayBefore(reference, 1).AddTag("casualReference/yesterday")
}

// Tomorrow returns the components of "tomorrow".
func Tomorrow(reference *ReferenceWithTimezone) *ParsingComponents {
	return TheDayAfter(reference, 1).AddTag("casualReference/tomorrow")
}

// TheDayBefore returns the components numDay days before the reference.
func TheDayBefore(reference *ReferenceWithTimezone, numDay float64) *ParsingComponents {
	return TheDayAfter(reference, -numDay)
}

// TheDayAfter returns the components nDays days after the reference.
func TheDayAfter(reference *ReferenceWithTimezone, nDays float64) *ParsingComponents {
	targetDate := reference.DateWithAdjustedTimezone()
	component := NewParsingComponents(reference, Components{})
	newDate := jsdate.FromTime(targetDate.Time())
	newDate.SetDate(newDate.GetDate() + nDays)

	AssignSimilarDate(component, newDate)
	ImplySimilarTime(component, newDate)
	component.Delete(ComponentMeridiem)
	return component
}

// Tonight returns the components of "tonight" (the JS default implyHour is 22).
func Tonight(reference *ReferenceWithTimezone, implyHour float64) *ParsingComponents {
	targetDate := reference.DateWithAdjustedTimezone()
	component := NewParsingComponents(reference, Components{})
	AssignSimilarDate(component, targetDate)
	component.Imply(ComponentHour, implyHour)
	component.Imply(ComponentMeridiem, PM)
	component.AddTag("casualReference/tonight")
	return component
}

// LastNight returns the components of "last night" (the JS default implyHour is 0).
func LastNight(reference *ReferenceWithTimezone, implyHour float64) *ParsingComponents {
	targetDate := reference.DateWithAdjustedTimezone()
	component := NewParsingComponents(reference, Components{})
	if targetDate.GetHours() < 6 {
		targetDate = jsdate.FromTime(targetDate.Time() - 24*60*60*1000)
	}
	AssignSimilarDate(component, targetDate)
	component.Imply(ComponentHour, implyHour)
	return component
}

// Evening returns the components of "evening" (the JS default implyHour is 20).
func Evening(reference *ReferenceWithTimezone, implyHour float64) *ParsingComponents {
	component := NewParsingComponents(reference, Components{})
	component.Imply(ComponentMeridiem, PM)
	component.Imply(ComponentHour, implyHour)
	component.AddTag("casualReference/evening")
	return component
}

// YesterdayEvening returns the components of "yesterday evening" (the JS default implyHour is
// 20).
func YesterdayEvening(reference *ReferenceWithTimezone, implyHour float64) *ParsingComponents {
	targetDate := reference.DateWithAdjustedTimezone()
	component := NewParsingComponents(reference, Components{})
	targetDate = jsdate.FromTime(targetDate.Time() - 24*60*60*1000)
	AssignSimilarDate(component, targetDate)
	component.Imply(ComponentHour, implyHour)
	component.Imply(ComponentMeridiem, PM)
	component.AddTag("casualReference/yesterday")
	component.AddTag("casualReference/evening")
	return component
}

// Midnight returns the components of "midnight".
func Midnight(reference *ReferenceWithTimezone) *ParsingComponents {
	component := NewParsingComponents(reference, Components{})
	if reference.DateWithAdjustedTimezone().GetHours() > 2 {
		// Unless it's very early morning (0~2AM), we assume the midnight is the coming midnight.
		component.AddDurationAsImplied(NewDuration(DurationValue{"day", 1}))
	}
	component.Assign(ComponentHour, 0)
	component.Imply(ComponentMinute, 0)
	component.Imply(ComponentSecond, 0)
	component.Imply(ComponentMillisecond, 0)
	component.AddTag("casualReference/midnight")
	return component
}

// Morning returns the components of "morning" (the JS default implyHour is 6).
func Morning(reference *ReferenceWithTimezone, implyHour float64) *ParsingComponents {
	component := NewParsingComponents(reference, Components{})
	component.Imply(ComponentMeridiem, AM)
	component.Imply(ComponentHour, implyHour)
	component.Imply(ComponentMinute, 0)
	component.Imply(ComponentSecond, 0)
	component.Imply(ComponentMillisecond, 0)
	component.AddTag("casualReference/morning")
	return component
}

// Afternoon returns the components of "afternoon" (the JS default implyHour is 15).
func Afternoon(reference *ReferenceWithTimezone, implyHour float64) *ParsingComponents {
	component := NewParsingComponents(reference, Components{})
	component.Imply(ComponentMeridiem, PM)
	component.Imply(ComponentHour, implyHour)
	component.Imply(ComponentMinute, 0)
	component.Imply(ComponentSecond, 0)
	component.Imply(ComponentMillisecond, 0)
	component.AddTag("casualReference/afternoon")
	return component
}

// Noon returns the components of "noon".
func Noon(reference *ReferenceWithTimezone) *ParsingComponents {
	component := NewParsingComponents(reference, Components{})
	component.Imply(ComponentMeridiem, AM)
	component.Assign(ComponentHour, 12)
	component.Imply(ComponentMinute, 0)
	component.Imply(ComponentSecond, 0)
	component.Imply(ComponentMillisecond, 0)
	component.AddTag("casualReference/noon")
	return component
}
