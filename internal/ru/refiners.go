// Ported from src/locales/ru/refiners/RUMergeDateRangeRefiner.ts and
// src/locales/ru/refiners/RUMergeDateTimeRefiner.ts.

package ru

import "github.com/wanasit/chrono-go/internal/core"

var mergeDateRangePattern = core.MustRegExp("^\\s*(и до|и по|до|по|-)\\s*$", "i")

// MergeDateRangeRefiner is RUMergeDateRangeRefiner: "c 06.09.1989 [до|по] 11.12.1996",
// "c пятницы и до среды".
type MergeDateRangeRefiner struct {
	core.AbstractMergeDateRangeRefiner
}

// NewMergeDateRangeRefiner returns a new RUMergeDateRangeRefiner.
func NewMergeDateRangeRefiner() *MergeDateRangeRefiner {
	r := &MergeDateRangeRefiner{}
	r.Init(r)
	return r
}

// PatternBetween returns the range connector pattern.
func (r *MergeDateRangeRefiner) PatternBetween() *core.RegExp { return mergeDateRangePattern }

// MergeDateTimeRefiner is RUMergeDateTimeRefiner: "2020-02-13 [в] 6:00", "Завтра [,] 7:00".
type MergeDateTimeRefiner struct {
	core.AbstractMergeDateTimeRefiner
}

// NewMergeDateTimeRefiner returns a new RUMergeDateTimeRefiner.
func NewMergeDateTimeRefiner() *MergeDateTimeRefiner {
	r := &MergeDateTimeRefiner{}
	r.Init(r)
	return r
}

// PatternBetween returns a new date/time connector pattern (a new RegExp per call, as in the
// original).
func (r *MergeDateTimeRefiner) PatternBetween() *core.RegExp {
	return core.MustRegExp("^\\s*(T|в|,|-)?\\s*$", "")
}
