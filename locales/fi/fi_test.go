package fi_test

import (
	"testing"
	"time"

	chrono "github.com/wanasit/chrono-go"
	"github.com/wanasit/chrono-go/locales/fi"
)

func TestPublicAPI(t *testing.T) {
	ref := chrono.ParsingReference{Instant: float64(time.Date(2012, 8, 10, 12, 0, 0, 0, time.UTC).UnixMilli()), Timezone: float64(0)}
	results := fi.Parse("Tapaaminen huomenna klo 15.30", ref, nil)
	if len(results) != 1 || results[0].Text != "huomenna klo 15.30" || results[0].Index != 11 {
		t.Fatalf("Parse: %+v", results)
	}
	if d, ok := results[0].Start.Get("day"); !ok || d != 11 {
		t.Errorf("day = %v, %v", d, ok)
	}
	if h, ok := results[0].Start.Get("hour"); !ok || h != 15 {
		t.Errorf("hour = %v, %v", h, ok)
	}
	if _, ok := fi.ParseDate("ei päivämäärää", ref, nil); ok {
		t.Error("ParseDate found a date")
	}
	if d, ok := fi.ParseDate("3 päivää sitten", ref, nil); !ok || d != ref.Instant-3*86400000 {
		t.Errorf("ParseDate = %v, %v", d, ok)
	}
	if c := fi.CreateCasualConfiguration(false); len(c.Parsers) != 11 || len(c.Refiners) != 9 {
		t.Errorf("casual: %d parsers, %d refiners", len(c.Parsers), len(c.Refiners))
	}
	if c := fi.CreateConfiguration(true, true); len(c.Parsers) != 8 || len(c.Refiners) != 9 {
		t.Errorf("strict: %d parsers, %d refiners", len(c.Parsers), len(c.Refiners))
	}
	if len(fi.Strict.Parsers) != 8 || len(fi.Casual.Parsers) != 11 {
		t.Errorf("Casual/Strict: %d/%d parsers", len(fi.Casual.Parsers), len(fi.Strict.Parsers))
	}
}
