// Ported from src/locales/uk/refiners/UKMergeDateRangeRefiner.ts and
// src/locales/uk/refiners/UKMergeDateTimeRefiner.ts.

package uk

import "github.com/wanasit/chrono-go/internal/core"

var mergeDateRangePattern = core.MustRegExp("^\\s*(і до|і по|до|по|-)\\s*$", "i")

// MergeDateRangeRefiner is UKMergeDateRangeRefiner: "[з|із] 06.09.1989 [до|по] 11.12.1996",
// "[з|із] п'ятниці і до середи".
type MergeDateRangeRefiner struct {
	core.AbstractMergeDateRangeRefiner
}

// NewMergeDateRangeRefiner returns a new UKMergeDateRangeRefiner.
func NewMergeDateRangeRefiner() *MergeDateRangeRefiner {
	r := &MergeDateRangeRefiner{}
	r.Init(r)
	return r
}

// PatternBetween returns the range connector pattern.
func (r *MergeDateRangeRefiner) PatternBetween() *core.RegExp { return mergeDateRangePattern }

// MergeDateTimeRefiner is UKMergeDateTimeRefiner: "2020-02-13 [в/у/о] 6:00", "Завтра [,] 7:00".
type MergeDateTimeRefiner struct {
	core.AbstractMergeDateTimeRefiner
}

// NewMergeDateTimeRefiner returns a new UKMergeDateTimeRefiner.
func NewMergeDateTimeRefiner() *MergeDateTimeRefiner {
	r := &MergeDateTimeRefiner{}
	r.Init(r)
	return r
}

// PatternBetween returns a new date/time connector pattern (a new RegExp per call, as in the
// original).
func (r *MergeDateTimeRefiner) PatternBetween() *core.RegExp {
	return core.MustRegExp("^\\s*(T|в|у|о|,|-)?\\s*$", "")
}
