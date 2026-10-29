package it_test

import (
	"fmt"
	"testing"

	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/it"
)

// The parser and refiner order of every configuration constructor, as the original builds it
// (constructor names recorded from orig/dist/cjs, mapped to the Go types).
var configurationOrders = []struct {
	cfg      *core.Configuration
	parsers  []string
	refiners []string
}{
	{it.CreateCasualConfiguration(),
		[]string{"*it.YearMonthDayParser", "*core.ISOFormatParser", "*core.SlashDateFormatParser", "*it.TimeUnitWithinFormatParser", "*it.MonthNameLittleEndianParser", "*it.MonthNameMiddleEndianParser", "*it.WeekdayParser", "*it.SlashMonthFormatParser", "*it.TimeExpressionParser", "*it.TimeUnitAgoFormatParser", "*it.TimeUnitLaterFormatParser", "*it.CasualDateParser", "*it.CasualTimeParser", "*it.MonthNameParser", "*it.RelativeDateFormatParser", "*it.TimeUnitCasualRelativeFormatParser"},
		[]string{"*core.OverlapRemovalRefiner", "*it.MergeRelativeAfterDateRefiner", "*it.MergeRelativeFollowByDateRefiner", "*core.OverlapRemovalRefiner", "*core.ExtractTimezoneOffsetRefiner", "*core.MergeWeekdayComponentRefiner", "*it.MergeDateTimeRefiner", "*core.ExtractTimezoneAbbrRefiner", "*core.OverlapRemovalRefiner", "*core.ForwardDateRefiner", "*core.UnlikelyFormatFilter", "*it.MergeDateTimeRefiner", "*it.ExtractYearSuffixRefiner", "*it.MergeDateRangeRefiner", "*it.UnlikelyFormatFilter"}},
	{it.CreateConfiguration(true),
		[]string{"*it.YearMonthDayParser", "*core.ISOFormatParser", "*core.SlashDateFormatParser", "*it.TimeUnitWithinFormatParser", "*it.MonthNameLittleEndianParser", "*it.MonthNameMiddleEndianParser", "*it.WeekdayParser", "*it.SlashMonthFormatParser", "*it.TimeExpressionParser", "*it.TimeUnitAgoFormatParser", "*it.TimeUnitLaterFormatParser"},
		[]string{"*core.OverlapRemovalRefiner", "*it.MergeRelativeAfterDateRefiner", "*it.MergeRelativeFollowByDateRefiner", "*core.OverlapRemovalRefiner", "*core.ExtractTimezoneOffsetRefiner", "*core.MergeWeekdayComponentRefiner", "*it.MergeDateTimeRefiner", "*core.ExtractTimezoneAbbrRefiner", "*core.OverlapRemovalRefiner", "*core.ForwardDateRefiner", "*core.UnlikelyFormatFilter", "*it.MergeDateTimeRefiner", "*it.ExtractYearSuffixRefiner", "*it.MergeDateRangeRefiner"}},
	{it.CreateConfiguration(false),
		[]string{"*it.YearMonthDayParser", "*core.ISOFormatParser", "*core.SlashDateFormatParser", "*it.TimeUnitWithinFormatParser", "*it.MonthNameLittleEndianParser", "*it.MonthNameMiddleEndianParser", "*it.WeekdayParser", "*it.SlashMonthFormatParser", "*it.TimeExpressionParser", "*it.TimeUnitAgoFormatParser", "*it.TimeUnitLaterFormatParser"},
		[]string{"*core.OverlapRemovalRefiner", "*it.MergeRelativeAfterDateRefiner", "*it.MergeRelativeFollowByDateRefiner", "*core.OverlapRemovalRefiner", "*core.ExtractTimezoneOffsetRefiner", "*core.MergeWeekdayComponentRefiner", "*it.MergeDateTimeRefiner", "*core.ExtractTimezoneAbbrRefiner", "*core.OverlapRemovalRefiner", "*core.ForwardDateRefiner", "*core.UnlikelyFormatFilter", "*it.MergeDateTimeRefiner", "*it.ExtractYearSuffixRefiner", "*it.MergeDateRangeRefiner"}},
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
