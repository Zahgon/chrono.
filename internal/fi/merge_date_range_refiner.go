// Ported from src/locales/fi/refiners/FIMergeDateRangeRefiner.ts.

package fi

import "github.com/wanasit/chrono-go/internal/core"

var mergeDateRangePattern = core.MustRegExp("^\\s*(-|–)\\s*$", "i")

// MergeDateRangeRefiner is FIMergeDateRangeRefiner ("1.3.2024 - 5.3.2024").
type MergeDateRangeRefiner struct {
	core.AbstractMergeDateRangeRefiner
}

// NewMergeDateRangeRefiner returns a new FIMergeDateRangeRefiner.
func NewMergeDateRangeRefiner() *MergeDateRangeRefiner {
	r := &MergeDateRangeRefiner{}
	r.Init(r)
	return r
}

// PatternBetween returns the range connector pattern.
func (r *MergeDateRangeRefiner) PatternBetween() *core.RegExp { return mergeDateRangePattern }
