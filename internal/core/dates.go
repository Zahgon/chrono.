// Ported from src/utils/dates.ts.

package core

import "github.com/wanasit/chrono-go/internal/jsdate"

func meridiemOf(hours float64) float64 {
	if hours < 12 {
		return AM
	}
	return PM
}

// AssignSimilarDate assigns (force updates) the component to the same day as target.
func AssignSimilarDate(component *ParsingComponents, target *jsdate.Date) {
	component.Assign(ComponentDay, target.GetDate())
	component.Assign(ComponentMonth, target.GetMonth()+1)
	component.Assign(ComponentYear, target.GetFullYear())
}

// AssignSimilarTime assigns (force updates) the component to the same time as target.
func AssignSimilarTime(component *ParsingComponents, target *jsdate.Date) {
	component.Assign(ComponentHour, target.GetHours())
	component.Assign(ComponentMinute, target.GetMinutes())
	component.Assign(ComponentSecond, target.GetSeconds())
	component.Assign(ComponentMillisecond, target.GetMilliseconds())
	component.Assign(ComponentMeridiem, meridiemOf(target.GetHours()))
}

// ImplySimilarDate implies (weakly updates) the component to the same day as target.
func ImplySimilarDate(component *ParsingComponents, target *jsdate.Date) {
	component.Imply(ComponentDay, target.GetDate())
	component.Imply(ComponentMonth, target.GetMonth()+1)
	component.Imply(ComponentYear, target.GetFullYear())
}

// ImplySimilarTime implies (weakly updates) the component to the same time as target.
func ImplySimilarTime(component *ParsingComponents, target *jsdate.Date) {
	component.Imply(ComponentHour, target.GetHours())
	component.Imply(ComponentMinute, target.GetMinutes())
	component.Imply(ComponentSecond, target.GetSeconds())
	component.Imply(ComponentMillisecond, target.GetMilliseconds())
	component.Imply(ComponentMeridiem, meridiemOf(target.GetHours()))
}
