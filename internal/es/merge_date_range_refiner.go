// Ported from src/locales/es/refiners/ESMergeDateRangeRefiner.ts.

package es

import "github.com/wanasit/chrono-go/internal/core"

var mergeDateRangePattern = core.MustRegExp("^\\s*(?:-)\\s*$", "i")

// MergeDateRangeRefiner is ESMergeDateRangeRefiner: merges the before and after results of a
// range ("lunes [-] viernes").
type MergeDateRangeRefiner struct {
	core.AbstractMergeDateRangeRefiner
}

// NewMergeDateRangeRefiner returns a new ESMergeDateRangeRefiner.
func NewMergeDateRangeRefiner() *MergeDateRangeRefiner {
	r := &MergeDateRangeRefiner{}
	r.Init(r)
	return r
}

// PatternBetween returns the range connector pattern.
func (r *MergeDateRangeRefiner) PatternBetween() *core.RegExp { return mergeDateRangePattern }
