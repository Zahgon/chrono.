// Ported from src/locales/de/refiners/DEMergeDateRangeRefiner.ts and
// src/locales/de/refiners/DEMergeDateTimeRefiner.ts.

package de

import "github.com/wanasit/chrono-go/internal/core"

var mergeDateRangePattern = core.MustRegExp("^\\s*(bis(?:\\s*(?:am|zum))?|-)\\s*$", "i")

// MergeDateRangeRefiner is DEMergeDateRangeRefiner: merges "<date> [bis] <date>" and
// "<date> [-] <date>" into a range.
type MergeDateRangeRefiner struct {
	core.AbstractMergeDateRangeRefiner
}

// NewMergeDateRangeRefiner returns a new DEMergeDateRangeRefiner.
func NewMergeDateRangeRefiner() *MergeDateRangeRefiner {
	r := &MergeDateRangeRefiner{}
	r.Init(r)
	return r
}

// PatternBetween returns the range connector pattern.
func (r *MergeDateRangeRefiner) PatternBetween() *core.RegExp { return mergeDateRangePattern }

// The original builds a new RegExp on every call; without a g/y flag that is unobservable, so it
// is compiled once.
var mergeDateTimePattern = core.MustRegExp("^\\s*(T|um|am|,|-)?\\s*$", "")

// MergeDateTimeRefiner is DEMergeDateTimeRefiner: merges a date-only result and a time-only
// result ("Morgen [um] 10 Uhr").
type MergeDateTimeRefiner struct {
	core.AbstractMergeDateTimeRefiner
}

// NewMergeDateTimeRefiner returns a new DEMergeDateTimeRefiner.
func NewMergeDateTimeRefiner() *MergeDateTimeRefiner {
	r := &MergeDateTimeRefiner{}
	r.Init(r)
	return r
}

// PatternBetween returns the date/time connector pattern.
func (r *MergeDateTimeRefiner) PatternBetween() *core.RegExp { return mergeDateTimePattern }
