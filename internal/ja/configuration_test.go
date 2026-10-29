package ja

import (
	"fmt"
	"strings"
	"testing"
)

// The parser and refiner order of node's chrono.ja.casual / chrono.ja.strict (constructor names).
func TestConfigurationOrder(t *testing.T) {
	refiners := "*core.OverlapRemovalRefiner,*core.ExtractTimezoneOffsetRefiner,*ja.MergeWeekdayComponentRefiner,*ja.MergeDateTimeRefiner,*ja.MergeDateRangeRefiner,*core.ExtractTimezoneAbbrRefiner,*core.OverlapRemovalRefiner,*core.ForwardDateRefiner,*core.UnlikelyFormatFilter"
	strict := "*core.ISOFormatParser,*ja.StandardParser,*ja.WeekdayParser,*ja.WeekdayWithParenthesesParser,*ja.SlashDateFormatParser,*ja.TimeExpressionParser"
	for _, c := range []struct {
		name, parsers string
		got           func() (string, string)
	}{
		{"casual", "*ja.CasualDateParser," + strict, func() (string, string) {
			cfg := CreateCasualConfiguration()
			return typeNames(cfg.Parsers), typeNames(cfg.Refiners)
		}},
		{"strict", strict, func() (string, string) {
			cfg := CreateConfiguration(true)
			return typeNames(cfg.Parsers), typeNames(cfg.Refiners)
		}},
		{"non-strict", strict, func() (string, string) {
			cfg := CreateConfiguration(false)
			return typeNames(cfg.Parsers), typeNames(cfg.Refiners)
		}},
	} {
		p, r := c.got()
		if p != c.parsers || r != refiners {
			t.Errorf("%s:\n got %s | %s\nwant %s | %s", c.name, p, r, c.parsers, refiners)
		}
	}
	if CreateConfiguration(true).Parsers[1] == CreateConfiguration(true).Parsers[1] {
		t.Error("CreateConfiguration must build fresh parser objects")
	}
}

func typeNames[T any](xs []T) string {
	names := make([]string, len(xs))
	for i, x := range xs {
		names[i] = fmt.Sprintf("%T", x)
	}
	return strings.Join(names, ",")
}
