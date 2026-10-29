// Ported from src/locales/en/refiners/ENMergeDateRangeRefiner.ts and
// src/locales/en/refiners/ENMergeDateTimeRefiner.ts.

package en

import "github.com/wanasit/chrono-go/internal/core"

var mergeDateRangePattern = core.MustRegExp("^\\s*(to|-|–|until|through|till)\\s*$", "i")

// MergeDateRangeRefiner is ENMergeDateRangeRefiner: "2020-02-13 [to] 2020-02-13",
// "Wednesday [-] Friday".
type MergeDateRangeRefiner struct {
	core.AbstractMergeDateRangeRefiner
}

// NewMergeDateRangeRefiner returns a new ENMergeDateRangeRefiner.
func NewMergeDateRangeRefiner() *MergeDateRangeRefiner {
	r := &MergeDateRangeRefiner{}
	r.Init(r)
	return r
}

// PatternBetween returns the range connector pattern.
func (r *MergeDateRangeRefiner) PatternBetween() *core.RegExp { return mergeDateRangePattern }

// MergeDateTimeRefiner is ENMergeDateTimeRefiner: "2020-02-13 [at] 6pm", "Tomorrow [after] 7am".
type MergeDateTimeRefiner struct {
	core.AbstractMergeDateTimeRefiner
}

// NewMergeDateTimeRefiner returns a new ENMergeDateTimeRefiner.
func NewMergeDateTimeRefiner() *MergeDateTimeRefiner {
	r := &MergeDateTimeRefiner{}
	r.Init(r)
	return r
}

// PatternBetween returns a new date/time connector pattern (a new RegExp per call, as in the
// original).
func (r *MergeDateTimeRefiner) PatternBetween() *core.RegExp {
	return core.MustRegExp("^\\s*(T|at|after|before|on|of|,|-|\\.|∙|:)?\\s*$", "")
}
