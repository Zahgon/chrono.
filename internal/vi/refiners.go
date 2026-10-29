// Ported from src/locales/vi/refiners/VIMergeDateRangeRefiner.ts,
// src/locales/vi/refiners/VIMergeDateTimeRefiner.ts and
// src/locales/vi/refiners/VIMergeWeekdayComponentRefiner.ts.

package vi

import "github.com/wanasit/chrono-go/internal/core"

var mergeDateRangePattern = core.MustRegExp("^\\s*(?:–|-|đến|tới|và)\\s*$", "")

// MergeDateRangeRefiner is VIMergeDateRangeRefiner ("15 tháng 3 [đến] 20 tháng 3").
type MergeDateRangeRefiner struct {
	core.AbstractMergeDateRangeRefiner
}

// NewMergeDateRangeRefiner returns a new VIMergeDateRangeRefiner.
func NewMergeDateRangeRefiner() *MergeDateRangeRefiner {
	r := &MergeDateRangeRefiner{}
	r.Init(r)
	return r
}

// PatternBetween returns the range connector pattern.
func (r *MergeDateRangeRefiner) PatternBetween() *core.RegExp { return mergeDateRangePattern }

var mergeDateTimePattern = core.MustRegExp("^\\s*(?:lúc|vào|,|T|-)?\\s*$", "")

// MergeDateTimeRefiner is VIMergeDateTimeRefiner ("ngày mai [lúc] 7 giờ").
type MergeDateTimeRefiner struct {
	core.AbstractMergeDateTimeRefiner
}

// NewMergeDateTimeRefiner returns a new VIMergeDateTimeRefiner.
func NewMergeDateTimeRefiner() *MergeDateTimeRefiner {
	r := &MergeDateTimeRefiner{}
	r.Init(r)
	return r
}

// PatternBetween returns the date/time connector pattern.
func (r *MergeDateTimeRefiner) PatternBetween() *core.RegExp { return mergeDateTimePattern }

// MergeWeekdayComponentRefiner is VIMergeWeekdayComponentRefiner, the common
// MergeWeekdayComponentRefiner unchanged: [thứ Hai] [lúc 7 giờ] => [thứ Hai lúc 7 giờ].
type MergeWeekdayComponentRefiner struct {
	*core.MergeWeekdayComponentRefiner
}

// NewMergeWeekdayComponentRefiner returns a new VIMergeWeekdayComponentRefiner.
func NewMergeWeekdayComponentRefiner() *MergeWeekdayComponentRefiner {
	return &MergeWeekdayComponentRefiner{core.NewMergeWeekdayComponentRefiner()}
}
