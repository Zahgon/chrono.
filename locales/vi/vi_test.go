package vi_test

import (
	"testing"
	"time"

	chrono "github.com/wanasit/chrono-go"
	"github.com/wanasit/chrono-go/locales/vi"
)

func TestPublicAPI(t *testing.T) {
	ref := chrono.ParsingReference{Instant: float64(time.Date(2012, 8, 10, 12, 0, 0, 0, time.UTC).UnixMilli()), Timezone: float64(0)}
	results := vi.Parse("Cuộc họp ngày mai lúc 7 giờ 30 phút tối", ref, nil)
	if len(results) != 1 || results[0].Text != "ngày mai lúc 7 giờ 30 phút tối" || results[0].Index != 9 {
		t.Fatalf("Parse: %+v", results)
	}
	if d, ok := results[0].Start.Get("day"); !ok || d != 11 {
		t.Errorf("day = %v, %v", d, ok)
	}
	if h, ok := results[0].Start.Get("hour"); !ok || h != 19 {
		t.Errorf("hour = %v, %v", h, ok)
	}
	if _, ok := vi.ParseDate("không có ngày", ref, nil); ok {
		t.Error("ParseDate found a date")
	}
	if d, ok := vi.ParseDate("2 ngày trước", ref, nil); !ok || d != ref.Instant-2*86400000 {
		t.Errorf("ParseDate = %v, %v", d, ok)
	}
	if c := vi.CreateCasualConfiguration(false); len(c.Parsers) != 14 || len(c.Refiners) != 10 {
		t.Errorf("casual: %d parsers, %d refiners", len(c.Parsers), len(c.Refiners))
	}
	if c := vi.CreateConfiguration(true, true); len(c.Parsers) != 11 || len(c.Refiners) != 10 {
		t.Errorf("strict: %d parsers, %d refiners", len(c.Parsers), len(c.Refiners))
	}
	if len(vi.Strict.Parsers) != 11 || len(vi.Casual.Parsers) != 14 {
		t.Errorf("Casual/Strict: %d/%d parsers", len(vi.Casual.Parsers), len(vi.Strict.Parsers))
	}
}
