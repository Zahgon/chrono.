package es_test

import (
	"fmt"
	"testing"

	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/es"
)

// The parser and refiner order of every configuration constructor, as the original builds it
// (constructor names recorded from orig/dist/cjs, mapped to the Go types).
var configurationOrders = []struct {
	cfg      *core.Configuration
	parsers  []string
	refiners []string
}{
	{es.CreateCasualConfiguration(true),
		[]string{"*core.ISOFormatParser", "*core.SlashDateFormatParser", "*es.WeekdayParser", "*es.TimeExpressionParser", "*es.MonthNameLittleEndianParser", "*es.TimeUnitWithinFormatParser", "*es.CasualDateParser", "*es.CasualTimeParser"},
		[]string{"*core.OverlapRemovalRefiner", "*core.ExtractTimezoneOffsetRefiner", "*core.MergeWeekdayComponentRefiner", "*es.MergeDateTimeRefiner", "*es.MergeDateRangeRefiner", "*core.ExtractTimezoneAbbrRefiner", "*core.OverlapRemovalRefiner", "*core.ForwardDateRefiner", "*core.UnlikelyFormatFilter"}},
	{es.CreateConfiguration(true, true),
		[]string{"*core.ISOFormatParser", "*core.SlashDateFormatParser", "*es.WeekdayParser", "*es.TimeExpressionParser", "*es.MonthNameLittleEndianParser", "*es.TimeUnitWithinFormatParser"},
		[]string{"*core.OverlapRemovalRefiner", "*core.ExtractTimezoneOffsetRefiner", "*core.MergeWeekdayComponentRefiner", "*es.MergeDateTimeRefiner", "*es.MergeDateRangeRefiner", "*core.ExtractTimezoneAbbrRefiner", "*core.OverlapRemovalRefiner", "*core.ForwardDateRefiner", "*core.UnlikelyFormatFilter"}},
	{es.CreateConfiguration(true, false),
		[]string{"*core.ISOFormatParser", "*core.SlashDateFormatParser", "*es.WeekdayParser", "*es.TimeExpressionParser", "*es.MonthNameLittleEndianParser", "*es.TimeUnitWithinFormatParser"},
		[]string{"*core.OverlapRemovalRefiner", "*core.ExtractTimezoneOffsetRefiner", "*core.MergeWeekdayComponentRefiner", "*es.MergeDateTimeRefiner", "*es.MergeDateRangeRefiner", "*core.ExtractTimezoneAbbrRefiner", "*core.OverlapRemovalRefiner", "*core.ForwardDateRefiner", "*core.UnlikelyFormatFilter"}},
	{es.CreateCasualConfiguration(false),
		[]string{"*core.ISOFormatParser", "*core.SlashDateFormatParser", "*es.WeekdayParser", "*es.TimeExpressionParser", "*es.MonthNameLittleEndianParser", "*es.TimeUnitWithinFormatParser", "*es.CasualDateParser", "*es.CasualTimeParser"},
		[]string{"*core.OverlapRemovalRefiner", "*core.ExtractTimezoneOffsetRefiner", "*core.MergeWeekdayComponentRefiner", "*es.MergeDateTimeRefiner", "*es.MergeDateRangeRefiner", "*core.ExtractTimezoneAbbrRefiner", "*core.OverlapRemovalRefiner", "*core.ForwardDateRefiner", "*core.UnlikelyFormatFilter"}},
	{es.CreateConfiguration(false, true),
		[]string{"*core.ISOFormatParser", "*core.SlashDateFormatParser", "*es.WeekdayParser", "*es.TimeExpressionParser", "*es.MonthNameLittleEndianParser", "*es.TimeUnitWithinFormatParser"},
		[]string{"*core.OverlapRemovalRefiner", "*core.ExtractTimezoneOffsetRefiner", "*core.MergeWeekdayComponentRefiner", "*es.MergeDateTimeRefiner", "*es.MergeDateRangeRefiner", "*core.ExtractTimezoneAbbrRefiner", "*core.OverlapRemovalRefiner", "*core.ForwardDateRefiner", "*core.UnlikelyFormatFilter"}},
	{es.CreateConfiguration(false, false),
		[]string{"*core.ISOFormatParser", "*core.SlashDateFormatParser", "*es.WeekdayParser", "*es.TimeExpressionParser", "*es.MonthNameLittleEndianParser", "*es.TimeUnitWithinFormatParser"},
		[]string{"*core.OverlapRemovalRefiner", "*core.ExtractTimezoneOffsetRefiner", "*core.MergeWeekdayComponentRefiner", "*es.MergeDateTimeRefiner", "*es.MergeDateRangeRefiner", "*core.ExtractTimezoneAbbrRefiner", "*core.OverlapRemovalRefiner", "*core.ForwardDateRefiner", "*core.UnlikelyFormatFilter"}},
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
