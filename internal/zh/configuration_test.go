package zh

import (
	"fmt"
	"strings"
	"testing"

	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/zh/hans"
	"github.com/wanasit/chrono-go/internal/zh/hant"
)

// The parser and refiner order of node's chrono.zh, chrono.zh.hans and chrono.zh.hant
// configurations (constructor names).
func TestConfigurationOrder(t *testing.T) {
	refs := func(s string) string {
		return "*core.OverlapRemovalRefiner,*core.MergeWeekdayComponentRefiner,*" + s + ".MergeDateRangeRefiner,*" + s +
			".MergeDateTimeRefiner,*core.ExtractTimezoneAbbrRefiner,*core.OverlapRemovalRefiner,*core.ForwardDateRefiner,*core.UnlikelyFormatFilter"
	}
	script := func(s string) string {
		return "*core.ISOFormatParser,*" + s + ".DateParser,*" + s + ".RelationWeekdayParser,*" + s + ".WeekdayParser,*" + s +
			".TimeExpressionParser,*" + s + ".DeadlineFormatParser,*" + s + ".AgoFormatParser"
	}
	zhStrict := "*core.ISOFormatParser,*hant.DateParser,*hans.DateParser,*hant.RelationWeekdayParser,*hans.RelationWeekdayParser," +
		"*hant.WeekdayParser,*hans.WeekdayParser,*hant.TimeExpressionParser,*hans.TimeExpressionParser," +
		"*hant.DeadlineFormatParser,*hans.DeadlineFormatParser"
	for _, c := range []struct {
		name, parsers, refiners string
		cfg                     *core.Configuration
	}{
		{"zh casual", "*hant.CasualDateParser," + zhStrict, refs("hant"), CreateCasualConfiguration()},
		{"zh", zhStrict, refs("hant"), CreateConfiguration()},
		{"hans casual", "*hans.CasualDateParser," + script("hans"), refs("hans"), hans.CreateCasualConfiguration()},
		{"hans", script("hans"), refs("hans"), hans.CreateConfiguration()},
		{"hant casual", "*hant.CasualDateParser," + script("hant"), refs("hant"), hant.CreateCasualConfiguration()},
		{"hant", script("hant"), refs("hant"), hant.CreateConfiguration()},
	} {
		p, r := typeNames(c.cfg.Parsers), typeNames(c.cfg.Refiners)
		if p != c.parsers || r != c.refiners {
			t.Errorf("%s:\n got %s | %s\nwant %s | %s", c.name, p, r, c.parsers, c.refiners)
		}
	}
}

func typeNames[T any](xs []T) string {
	names := make([]string, len(xs))
	for i, x := range xs {
		names[i] = fmt.Sprintf("%T", x)
	}
	return strings.Join(names, ",")
}
