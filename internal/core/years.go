// Ported from src/calculation/years.ts.

package core

import (
	"math"

	"github.com/wanasit/chrono-go/internal/jsdate"
)

// FindMostLikelyADYear finds the most likely year from a raw number: 1997 => 1997, 97 => 1997,
// 12 => 2012.
func FindMostLikelyADYear(yearNumber float64) float64 {
	if yearNumber < 100 {
		if yearNumber > 50 {
			yearNumber = yearNumber + 1900
		} else {
			yearNumber = yearNumber + 2000
		}
	}
	return yearNumber
}

// FindYearClosestToRef returns the year in which day/month is closest to refDate (a time value).
func FindYearClosestToRef(refDate float64, day, month float64) float64 {
	ref := jsdate.FromTime(refDate)
	date := ref.Clone()
	date.SetMonth(month - 1)
	date.SetDate(day)
	nextYear := AddDuration(date, NewDuration(DurationValue{"year", 1}))
	lastYear := AddDuration(date, NewDuration(DurationValue{"year", -1}))
	if math.Abs(nextYear.Time()-ref.Time()) < math.Abs(date.Time()-ref.Time()) {
		date = nextYear
	} else if math.Abs(lastYear.Time()-ref.Time()) < math.Abs(date.Time()-ref.Time()) {
		date = lastYear
	}
	return date.GetFullYear()
}
