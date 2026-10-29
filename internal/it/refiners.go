// Ported from src/locales/it/refiners/ITMergeDateRangeRefiner.ts and
// src/locales/it/refiners/ITMergeDateTimeRefiner.ts.

package it

import "github.com/wanasit/chrono-go/internal/core"

var mergeDateRangePattern = core.MustRegExp("^\\s*(a|al|-|–|fino\\s*a|fino\\s*al)\\s*$", "i")

// MergeDateRangeRefiner is ITMergeDateRangeRefiner: "2020-02-13 [a] 2020-02-13",
// "Mercoledì [-] Venerdì".
type MergeDateRangeRefiner struct {
	core.AbstractMergeDateRangeRefiner
}

// NewMergeDateRangeRefiner returns a new ITMergeDateRangeRefiner.
func NewMergeDateRangeRefiner() *MergeDateRangeRefiner {
	r := &MergeDateRangeRefiner{}
	r.Init(r)
	return r
}

// PatternBetween returns the range connector pattern.
func (r *MergeDateRangeRefiner) PatternBetween() *core.RegExp { return mergeDateRangePattern }

// The original builds a new (stateless, non-global) RegExp on every patternBetween call; one
// package-level value is equivalent.
var mergeDateTimePattern = core.MustRegExp("^\\s*(T|alle?|dopo\\s*le|prima\\s*delle?|,|-|\\.|∙|:)?\\s*$", "")

// MergeDateTimeRefiner is ITMergeDateTimeRefiner: "2020-02-13 [alle] 18", "Domani [dopo le] 7".
type MergeDateTimeRefiner struct {
	core.AbstractMergeDateTimeRefiner
}

// NewMergeDateTimeRefiner returns a new ITMergeDateTimeRefiner.
func NewMergeDateTimeRefiner() *MergeDateTimeRefiner {
	r := &MergeDateTimeRefiner{}
	r.Init(r)
	return r
}

// PatternBetween returns the date/time connector pattern.
func (r *MergeDateTimeRefiner) PatternBetween() *core.RegExp { return mergeDateTimePattern }
