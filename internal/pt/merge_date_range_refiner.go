// Ported from src/locales/pt/refiners/PTMergeDateRangeRefiner.ts.

package pt

import "github.com/wanasit/chrono-go/internal/core"

var mergeDateRangePattern = core.MustRegExp("^\\s*(?:-)\\s*$", "i")

// MergeDateRangeRefiner is PTMergeDateRangeRefiner: merges the before and after results of a
// range ("segunda [-] sexta").
type MergeDateRangeRefiner struct {
	core.AbstractMergeDateRangeRefiner
}

// NewMergeDateRangeRefiner returns a new PTMergeDateRangeRefiner.
func NewMergeDateRangeRefiner() *MergeDateRangeRefiner {
	r := &MergeDateRangeRefiner{}
	r.Init(r)
	return r
}

// PatternBetween returns the range connector pattern.
func (r *MergeDateRangeRefiner) PatternBetween() *core.RegExp { return mergeDateRangePattern }
