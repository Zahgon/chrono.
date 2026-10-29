// Ported from src/locales/vi/parsers/VIStandardParser.ts.

package vi

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

// ngày 15 tháng 3 năm 1975  |  15 tháng 3 năm 1975  |  15 tháng 3
var standardPattern = core.MustRegExp(
	"(?:ng\u00e0y\\s*)?"+
		"([0-9]{1,2})"+
		"\\s*th\u00e1ng\\s*"+
		"([0-9]{1,2})"+
		"(?:\\s*n\u0103m\\s*("+
		YearPattern+
		"))?"+
		"(?=\\W|$)",
	"i",
)

const (
	standardDayGroup   = 1
	standardMonthGroup = 2
	standardYearGroup  = 3
)

// StandardParser is VIStandardParser ("ngày 15 tháng 3 năm 1975", "15 tháng 3").
type StandardParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewStandardParser returns a new VIStandardParser.
func NewStandardParser() *StandardParser {
	p := &StandardParser{}
	p.Self = p
	return p
}

// InnerPattern returns the day-month(-year) pattern.
func (p *StandardParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return standardPattern
}

// InnerExtract returns the result of the date.
func (p *StandardParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	day := js.ParseInt(match.Get(standardDayGroup), 0)
	month := js.ParseInt(match.Get(standardMonthGroup), 0)
	if day > 31 || month > 12 {
		return nil
	}

	result := context.CreateParsingResult(match.Index, match.Get(0), nil, nil)
	result.Start.Assign(core.ComponentDay, day)
	result.Start.Assign(core.ComponentMonth, month)

	if match.Truthy(standardYearGroup) {
		result.Start.Assign(core.ComponentYear, ParseYear(match.Get(standardYearGroup)))
	} else {
		result.Start.Imply(core.ComponentYear, core.FindYearClosestToRef(context.RefDate, day, month))
	}
	return result
}
