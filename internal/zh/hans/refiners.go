// Ported from src/locales/zh/hans/refiners/ZHHansMergeDateRangeRefiner.ts and
// src/locales/zh/hans/refiners/ZHHansMergeDateTimeRefiner.ts.

package hans

import "github.com/wanasit/chrono-go/internal/core"

var mergeDateRangePattern = core.MustRegExp("^\\s*(至|到|-|~|～|－|ー)\\s*$", "i")

// MergeDateRangeRefiner is ZHHansMergeDateRangeRefiner: "9月3日[至]9月5日".
type MergeDateRangeRefiner struct {
	core.AbstractMergeDateRangeRefiner
}

// NewMergeDateRangeRefiner returns a new ZHHansMergeDateRangeRefiner.
func NewMergeDateRangeRefiner() *MergeDateRangeRefiner {
	r := &MergeDateRangeRefiner{}
	r.Init(r)
	return r
}

// PatternBetween returns the range connector pattern.
func (r *MergeDateRangeRefiner) PatternBetween() *core.RegExp { return mergeDateRangePattern }

var mergeDateTimePattern = core.MustRegExp("^\\s*$", "i")

// MergeDateTimeRefiner is ZHHansMergeDateTimeRefiner: a date directly followed by a time.
type MergeDateTimeRefiner struct {
	core.AbstractMergeDateTimeRefiner
}

// NewMergeDateTimeRefiner returns a new ZHHansMergeDateTimeRefiner.
func NewMergeDateTimeRefiner() *MergeDateTimeRefiner {
	r := &MergeDateTimeRefiner{}
	r.Init(r)
	return r
}

// PatternBetween returns the date/time connector pattern.
func (r *MergeDateTimeRefiner) PatternBetween() *core.RegExp { return mergeDateTimePattern }
