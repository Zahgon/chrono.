package fr_test

import (
	"fmt"
	"testing"

	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/fr"
)

// The parser and refiner order of every configuration constructor, as the original builds it
// (constructor names recorded from orig/dist/cjs, mapped to the Go types).
var configurationOrders = []struct {
	cfg      *core.Configuration
	parsers  []string
	refiners []string
}{
	{fr.CreateCasualConfiguration(true),
		[]string{"*fr.TimeUnitRelativeFormatParser", "*fr.CasualTimeParser", "*fr.CasualDateParser", "*core.ISOFormatParser", "*core.SlashDateFormatParser", "*fr.MonthNameLittleEndianParser", "*fr.TimeExpressionParser", "*fr.SpecificTimeExpressionParser", "*fr.TimeUnitAgoFormatParser", "*fr.TimeUnitWithinFormatParser", "*fr.WeekdayParser"},
		[]string{"*core.OverlapRemovalRefiner", "*core.ExtractTimezoneOffsetRefiner", "*core.MergeWeekdayComponentRefiner", "*fr.MergeDateTimeRefiner", "*fr.MergeDateRangeRefiner", "*core.ExtractTimezoneAbbrRefiner", "*core.OverlapRemovalRefiner", "*core.ForwardDateRefiner", "*core.UnlikelyFormatFilter"}},
	{fr.CreateConfiguration(true, true),
		[]string{"*core.ISOFormatParser", "*core.SlashDateFormatParser", "*fr.MonthNameLittleEndianParser", "*fr.TimeExpressionParser", "*fr.SpecificTimeExpressionParser", "*fr.TimeUnitAgoFormatParser", "*fr.TimeUnitWithinFormatParser", "*fr.WeekdayParser"},
		[]string{"*core.OverlapRemovalRefiner", "*core.ExtractTimezoneOffsetRefiner", "*core.MergeWeekdayComponentRefiner", "*fr.MergeDateTimeRefiner", "*fr.MergeDateRangeRefiner", "*core.ExtractTimezoneAbbrRefiner", "*core.OverlapRemovalRefiner", "*core.ForwardDateRefiner", "*core.UnlikelyFormatFilter"}},
	{fr.CreateConfiguration(true, false),
		[]string{"*core.ISOFormatParser", "*core.SlashDateFormatParser", "*fr.MonthNameLittleEndianParser", "*fr.TimeExpressionParser", "*fr.SpecificTimeExpressionParser", "*fr.TimeUnitAgoFormatParser", "*fr.TimeUnitWithinFormatParser", "*fr.WeekdayParser"},
		[]string{"*core.OverlapRemovalRefiner", "*core.ExtractTimezoneOffsetRefiner", "*core.MergeWeekdayComponentRefiner", "*fr.MergeDateTimeRefiner", "*fr.MergeDateRangeRefiner", "*core.ExtractTimezoneAbbrRefiner", "*core.OverlapRemovalRefiner", "*core.ForwardDateRefiner", "*core.UnlikelyFormatFilter"}},
	{fr.CreateCasualConfiguration(false),
		[]string{"*fr.TimeUnitRelativeFormatParser", "*fr.CasualTimeParser", "*fr.CasualDateParser", "*core.ISOFormatParser", "*core.SlashDateFormatParser", "*fr.MonthNameLittleEndianParser", "*fr.TimeExpressionParser", "*fr.SpecificTimeExpressionParser", "*fr.TimeUnitAgoFormatParser", "*fr.TimeUnitWithinFormatParser", "*fr.WeekdayParser"},
		[]string{"*core.OverlapRemovalRefiner", "*core.ExtractTimezoneOffsetRefiner", "*core.MergeWeekdayComponentRefiner", "*fr.MergeDateTimeRefiner", "*fr.MergeDateRangeRefiner", "*core.ExtractTimezoneAbbrRefiner", "*core.OverlapRemovalRefiner", "*core.ForwardDateRefiner", "*core.UnlikelyFormatFilter"}},
	{fr.CreateConfiguration(false, true),
		[]string{"*core.ISOFormatParser", "*core.SlashDateFormatParser", "*fr.MonthNameLittleEndianParser", "*fr.TimeExpressionParser", "*fr.SpecificTimeExpressionParser", "*fr.TimeUnitAgoFormatParser", "*fr.TimeUnitWithinFormatParser", "*fr.WeekdayParser"},
		[]string{"*core.OverlapRemovalRefiner", "*core.ExtractTimezoneOffsetRefiner", "*core.MergeWeekdayComponentRefiner", "*fr.MergeDateTimeRefiner", "*fr.MergeDateRangeRefiner", "*core.ExtractTimezoneAbbrRefiner", "*core.OverlapRemovalRefiner", "*core.ForwardDateRefiner", "*core.UnlikelyFormatFilter"}},
	{fr.CreateConfiguration(false, false),
		[]string{"*core.ISOFormatParser", "*core.SlashDateFormatParser", "*fr.MonthNameLittleEndianParser", "*fr.TimeExpressionParser", "*fr.SpecificTimeExpressionParser", "*fr.TimeUnitAgoFormatParser", "*fr.TimeUnitWithinFormatParser", "*fr.WeekdayParser"},
		[]string{"*core.OverlapRemovalRefiner", "*core.ExtractTimezoneOffsetRefiner", "*core.MergeWeekdayComponentRefiner", "*fr.MergeDateTimeRefiner", "*fr.MergeDateRangeRefiner", "*core.ExtractTimezoneAbbrRefiner", "*core.OverlapRemovalRefiner", "*core.ForwardDateRefiner", "*core.UnlikelyFormatFilter"}},
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
