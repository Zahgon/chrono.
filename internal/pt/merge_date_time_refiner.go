// Ported from src/locales/pt/refiners/PTMergeDateTimeRefiner.ts.

package pt

import "github.com/wanasit/chrono-go/internal/core"

// The original builds a new (stateless, non-global) RegExp on every patternBetween call; one
// package-level value is equivalent.
var mergeDateTimePattern = core.MustRegExp("^\\s*(?:,|à)?\\s*$", "")

// MergeDateTimeRefiner is PTMergeDateTimeRefiner: merges a date-only result and a time-only
// result ("amanhã [à]s 10").
type MergeDateTimeRefiner struct {
	core.AbstractMergeDateTimeRefiner
}

// NewMergeDateTimeRefiner returns a new PTMergeDateTimeRefiner.
func NewMergeDateTimeRefiner() *MergeDateTimeRefiner {
	r := &MergeDateTimeRefiner{}
	r.Init(r)
	return r
}

// PatternBetween returns the date/time connector pattern.
func (r *MergeDateTimeRefiner) PatternBetween() *core.RegExp { return mergeDateTimePattern }
