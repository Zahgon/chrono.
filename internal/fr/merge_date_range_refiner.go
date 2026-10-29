// Ported from src/locales/fr/refiners/FRMergeDateRangeRefiner.ts.

package fr

import "github.com/wanasit/chrono-go/internal/core"

var mergeDateRangePattern = core.MustRegExp("^\\s*(à|a|au|-)\\s*$", "i")

// MergeDateRangeRefiner is FRMergeDateRangeRefiner: merges the before and after results of a
// range ("lundi [au] vendredi").
type MergeDateRangeRefiner struct {
	core.AbstractMergeDateRangeRefiner
}

// NewMergeDateRangeRefiner returns a new FRMergeDateRangeRefiner.
func NewMergeDateRangeRefiner() *MergeDateRangeRefiner {
	r := &MergeDateRangeRefiner{}
	r.Init(r)
	return r
}

// PatternBetween returns the range connector pattern.
func (r *MergeDateRangeRefiner) PatternBetween() *core.RegExp { return mergeDateRangePattern }
