// Ported from src/locales/es/refiners/ESMergeDateTimeRefiner.ts.

package es

import "github.com/wanasit/chrono-go/internal/core"

// The original builds a new (stateless, non-global) RegExp on every patternBetween call; one
// package-level value is equivalent.
var mergeDateTimePattern = core.MustRegExp("^\\s*(?:,|de|aslas|a)?\\s*$", "")

// MergeDateTimeRefiner is ESMergeDateTimeRefiner: merges a date-only result and a time-only
// result ("mañana [a] las 10").
type MergeDateTimeRefiner struct {
	core.AbstractMergeDateTimeRefiner
}

// NewMergeDateTimeRefiner returns a new ESMergeDateTimeRefiner.
func NewMergeDateTimeRefiner() *MergeDateTimeRefiner {
	r := &MergeDateTimeRefiner{}
	r.Init(r)
	return r
}

// PatternBetween returns the date/time connector pattern.
func (r *MergeDateTimeRefiner) PatternBetween() *core.RegExp { return mergeDateTimePattern }
