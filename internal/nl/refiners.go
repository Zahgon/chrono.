// Ported from src/locales/nl/refiners/NLMergeDateRangeRefiner.ts and
// src/locales/nl/refiners/NLMergeDateTimeRefiner.ts.

package nl

import "github.com/wanasit/chrono-go/internal/core"

var mergeDateRangePattern = core.MustRegExp("^\\s*(tot|-)\\s*$", "i")

// MergeDateRangeRefiner is NLMergeDateRangeRefiner: merges "2020-02-13 [tot] 2020-02-13" and
// "Wednesday [-] Friday" into a range.
type MergeDateRangeRefiner struct {
	core.AbstractMergeDateRangeRefiner
}

// NewMergeDateRangeRefiner returns a new NLMergeDateRangeRefiner.
func NewMergeDateRangeRefiner() *MergeDateRangeRefiner {
	r := &MergeDateRangeRefiner{}
	r.Init(r)
	return r
}

// PatternBetween returns the range connector pattern.
func (r *MergeDateRangeRefiner) PatternBetween() *core.RegExp { return mergeDateRangePattern }

// The original builds a new RegExp on every call; without a g/y flag that is unobservable, so it
// is compiled once.
var mergeDateTimePattern = core.MustRegExp("^\\s*(om|na|voor|in de|,|-)?\\s*$", "")

// MergeDateTimeRefiner is NLMergeDateTimeRefiner: merges a date-only result and a time-only
// result ("2020-02-13 [om] 18:00", "morgen [na] 7 uur").
type MergeDateTimeRefiner struct {
	core.AbstractMergeDateTimeRefiner
}

// NewMergeDateTimeRefiner returns a new NLMergeDateTimeRefiner.
func NewMergeDateTimeRefiner() *MergeDateTimeRefiner {
	r := &MergeDateTimeRefiner{}
	r.Init(r)
	return r
}

// PatternBetween returns the date/time connector pattern.
func (r *MergeDateTimeRefiner) PatternBetween() *core.RegExp { return mergeDateTimePattern }
