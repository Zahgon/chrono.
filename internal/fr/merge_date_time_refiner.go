// Ported from src/locales/fr/refiners/FRMergeDateTimeRefiner.ts.

package fr

import "github.com/wanasit/chrono-go/internal/core"

// The original builds a new (stateless, non-global) RegExp on every patternBetween call; one
// package-level value is equivalent.
var mergeDateTimePattern = core.MustRegExp("^\\s*(T|à|a|au|vers|de|,|-)?\\s*$", "")

// MergeDateTimeRefiner is FRMergeDateTimeRefiner: merges a date-only result and a time-only
// result ("demain [à] 10h").
type MergeDateTimeRefiner struct {
	core.AbstractMergeDateTimeRefiner
}

// NewMergeDateTimeRefiner returns a new FRMergeDateTimeRefiner.
func NewMergeDateTimeRefiner() *MergeDateTimeRefiner {
	r := &MergeDateTimeRefiner{}
	r.Init(r)
	return r
}

// PatternBetween returns the date/time connector pattern.
func (r *MergeDateTimeRefiner) PatternBetween() *core.RegExp { return mergeDateTimePattern }
