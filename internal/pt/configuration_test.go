package pt_test

import (
	"fmt"
	"testing"

	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/pt"
)

// The parser and refiner order of every configuration constructor, as the original builds it
// (constructor names recorded from orig/dist/cjs, mapped to the Go types).
var configurationOrders = []struct {
	cfg      *core.Configuration
	parsers  []string
	refiners []string
}{
	{pt.CreateCasualConfiguration(true),
		[]string{"*core.ISOFormatParser", "*core.SlashDateFormatParser", "*pt.WeekdayParser", "*pt.TimeExpressionParser", "*pt.MonthNameLittleEndianParser", "*pt.CasualDateParser", "*pt.CasualTimeParser"},
		[]string{"*core.OverlapRemovalRefiner", "*core.ExtractTimezoneOffsetRefiner", "*core.MergeWeekdayComponentRefiner", "*pt.MergeDateTimeRefiner", "*pt.MergeDateRangeRefiner", "*core.ExtractTimezoneAbbrRefiner", "*core.OverlapRemovalRefiner", "*core.ForwardDateRefiner", "*core.UnlikelyFormatFilter"}},
	{pt.CreateConfiguration(true, true),
		[]string{"*core.ISOFormatParser", "*core.SlashDateFormatParser", "*pt.WeekdayParser", "*pt.TimeExpressionParser", "*pt.MonthNameLittleEndianParser"},
		[]string{"*core.OverlapRemovalRefiner", "*core.ExtractTimezoneOffsetRefiner", "*core.MergeWeekdayComponentRefiner", "*pt.MergeDateTimeRefiner", "*pt.MergeDateRangeRefiner", "*core.ExtractTimezoneAbbrRefiner", "*core.OverlapRemovalRefiner", "*core.ForwardDateRefiner", "*core.UnlikelyFormatFilter"}},
	{pt.CreateConfiguration(true, false),
		[]string{"*core.ISOFormatParser", "*core.SlashDateFormatParser", "*pt.WeekdayParser", "*pt.TimeExpressionParser", "*pt.MonthNameLittleEndianParser"},
		[]string{"*core.OverlapRemovalRefiner", "*core.ExtractTimezoneOffsetRefiner", "*core.MergeWeekdayComponentRefiner", "*pt.MergeDateTimeRefiner", "*pt.MergeDateRangeRefiner", "*core.ExtractTimezoneAbbrRefiner", "*core.OverlapRemovalRefiner", "*core.ForwardDateRefiner", "*core.UnlikelyFormatFilter"}},
	{pt.CreateCasualConfiguration(false),
		[]string{"*core.ISOFormatParser", "*core.SlashDateFormatParser", "*pt.WeekdayParser", "*pt.TimeExpressionParser", "*pt.MonthNameLittleEndianParser", "*pt.CasualDateParser", "*pt.CasualTimeParser"},
		[]string{"*core.OverlapRemovalRefiner", "*core.ExtractTimezoneOffsetRefiner", "*core.MergeWeekdayComponentRefiner", "*pt.MergeDateTimeRefiner", "*pt.MergeDateRangeRefiner", "*core.ExtractTimezoneAbbrRefiner", "*core.OverlapRemovalRefiner", "*core.ForwardDateRefiner", "*core.UnlikelyFormatFilter"}},
	{pt.CreateConfiguration(false, true),
		[]string{"*core.ISOFormatParser", "*core.SlashDateFormatParser", "*pt.WeekdayParser", "*pt.TimeExpressionParser", "*pt.MonthNameLittleEndianParser"},
		[]string{"*core.OverlapRemovalRefiner", "*core.ExtractTimezoneOffsetRefiner", "*core.MergeWeekdayComponentRefiner", "*pt.MergeDateTimeRefiner", "*pt.MergeDateRangeRefiner", "*core.ExtractTimezoneAbbrRefiner", "*core.OverlapRemovalRefiner", "*core.ForwardDateRefiner", "*core.UnlikelyFormatFilter"}},
	{pt.CreateConfiguration(false, false),
		[]string{"*core.ISOFormatParser", "*core.SlashDateFormatParser", "*pt.WeekdayParser", "*pt.TimeExpressionParser", "*pt.MonthNameLittleEndianParser"},
		[]string{"*core.OverlapRemovalRefiner", "*core.ExtractTimezoneOffsetRefiner", "*core.MergeWeekdayComponentRefiner", "*pt.MergeDateTimeRefiner", "*pt.MergeDateRangeRefiner", "*core.ExtractTimezoneAbbrRefiner", "*core.OverlapRemovalRefiner", "*core.ForwardDateRefiner", "*core.UnlikelyFormatFilter"}},
}

func TestConfigurationOrder(t *testing.T) {
	for i, c := range configurationOrders {
		var ps, rs []string
		for _, p := range c.cfg.Parsers {
			ps = append(ps, fmt.Sprintf("%T", p))
		}
		for _, r := range c.cfg.Refiners {
			rs = append(rs, fmt.Sprintf("%T", r))
		}
		if fmt.Sprint(ps) != fmt.Sprint(c.parsers) {
			t.Errorf("configuration %d parsers:\n got %v\nwant %v", i, ps, c.parsers)
		}
		if fmt.Sprint(rs) != fmt.Sprint(c.refiners) {
			t.Errorf("configuration %d refiners:\n got %v\nwant %v", i, rs, c.refiners)
		}
	}
}
