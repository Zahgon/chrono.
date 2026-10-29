// Ported from src/locales/fi/refiners/FIMergeDateTimeRefiner.ts.

package fi

import "github.com/wanasit/chrono-go/internal/core"

// MergeDateTimeRefiner is FIMergeDateTimeRefiner ("huomenna [klo] 15").
type MergeDateTimeRefiner struct {
	core.AbstractMergeDateTimeRefiner
}

// NewMergeDateTimeRefiner returns a new FIMergeDateTimeRefiner.
func NewMergeDateTimeRefiner() *MergeDateTimeRefiner {
	r := &MergeDateTimeRefiner{}
	r.Init(r)
	return r
}

// PatternBetween returns a new date/time connector pattern (a new RegExp per call, as in the
// original).
func (r *MergeDateTimeRefiner) PatternBetween() *core.RegExp {
	return core.MustRegExp("^\\s*(T|klo|kello|,|-)?\\s*$", "")
}
