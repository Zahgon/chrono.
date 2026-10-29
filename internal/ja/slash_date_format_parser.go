// Ported from src/locales/ja/parsers/JPSlashDateFormatParser.ts.

package ja

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

// Date format with slash "/" between numbers. YYYY/MM/DD or MM/DD. Big-endian slash format is
// used in Japan. For examples: 7/31, 2020/7/12.
var slashDatePattern = core.MustRegExp(
	"([0-9０-９]{4}[\\/|\\／])?"+"([0-1０-１]{0,1}[0-9０-９]{1})(?:[\\/|\\／]([0-3０-３]{0,1}[0-9０-９]{1}))",
	"i",
)

const (
	slashYearGroup  = 1
	slashMonthGroup = 2
	slashDayGroup   = 3
)

// SlashDateFormatParser is JPSlashDateFormatParser ("7/31", "2020/7/12", full-width digits too).
type SlashDateFormatParser struct {
	_ byte // non-zero size: each constructor call returns a distinct object, as `new` does
}

// NewSlashDateFormatParser returns a new JPSlashDateFormatParser.
func NewSlashDateFormatParser() *SlashDateFormatParser { return &SlashDateFormatParser{} }

// Pattern returns the slash date pattern.
func (p *SlashDateFormatParser) Pattern(context *core.ParsingContext) *core.RegExp {
	return slashDatePattern
}

// Extract returns the components of the slash date, or nil for an invalid month or day.
func (p *SlashDateFormatParser) Extract(context *core.ParsingContext, match *core.RegExpMatch) any {
	result := context.CreateParsingComponents(nil)
	month := js.ParseInt(ToHankaku(match.Get(slashMonthGroup)), 0)
	day := js.ParseInt(ToHankaku(match.Get(slashDayGroup)), 0)
	if month < 1 || month > 12 {
		return nil
	}

	if day < 1 || day > 31 {
		return nil
	}

	result.Assign(core.ComponentDay, day)
	result.Assign(core.ComponentMonth, month)

	if match.Truthy(slashYearGroup) {
		rawYearNumber := js.ParseInt(ToHankaku(match.Get(slashYearGroup)), 0)
		year := core.FindMostLikelyADYear(rawYearNumber)
		result.Assign(core.ComponentYear, year)
	} else {
		year := core.FindYearClosestToRef(context.Reference.Instant, day, month)
		result.Imply(core.ComponentYear, year)
	}

	return result
}
