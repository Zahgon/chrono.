// Ported from src/locales/zh/hant/refiners/ZHHantMergeDateRangeRefiner.ts and
// src/locales/zh/hant/refiners/ZHHantMergeDateTimeRefiner.ts.

package hant

import "github.com/wanasit/chrono-go/internal/core"

var mergeDateRangePattern = core.MustRegExp("^\\s*(至|到|\\-|\\~|～|－|ー)\\s*$", "i")

// MergeDateRangeRefiner is ZHHantMergeDateRangeRefiner: "9月3日[至]9月5日".
type MergeDateRangeRefiner struct {
	core.AbstractMergeDateRangeRefiner
}

// NewMergeDateRangeRefiner returns a new ZHHantMergeDateRangeRefiner.
func NewMergeDateRangeRefiner() *MergeDateRangeRefiner {
	r := &MergeDateRangeRefiner{}
	r.Init(r)
	return r
}

// PatternBetween returns the range connector pattern.
func (r *MergeDateRangeRefiner) PatternBetween() *core.RegExp { return mergeDateRangePattern }

var mergeDateTimePattern = core.MustRegExp("^\\s*$", "i")

// MergeDateTimeRefiner is ZHHantMergeDateTimeRefiner: a date directly followed by a time.
type MergeDateTimeRefiner struct {
	core.AbstractMergeDateTimeRefiner
}

// NewMergeDateTimeRefiner returns a new ZHHantMergeDateTimeRefiner.
func NewMergeDateTimeRefiner() *MergeDateTimeRefiner {
	r := &MergeDateTimeRefiner{}
	r.Init(r)
	return r
}

// PatternBetween returns the date/time connector pattern.
func (r *MergeDateTimeRefiner) PatternBetween() *core.RegExp { return mergeDateTimePattern }
