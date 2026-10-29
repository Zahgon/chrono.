// Ported from src/locales/ja/parsers/JPStandardParser.ts.

package ja

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

var standardPattern = core.MustRegExp(
	"(?:(?:([同今本])|((昭和|平成|令和)?([0-9０-９]{1,4}|元)))年\\s*)?([0-9０-９]{1,2})月\\s*([0-9０-９]{1,2})日",
	"i",
)

const (
	standardSpecialYearGroup = 1
	standardTypicalYearGroup = 2
	standardEraGroup         = 3
	standardYearNumberGroup  = 4
	standardMonthGroup       = 5
	standardDayGroup         = 6
)

// StandardParser is JPStandardParser ("2012年3月31日", "令和元年5月1日", "同年4月2日", ...).
type StandardParser struct {
	_ byte // non-zero size: each constructor call returns a distinct object, as `new` does
}

// NewStandardParser returns a new JPStandardParser.
func NewStandardParser() *StandardParser { return &StandardParser{} }

// Pattern returns the standard date pattern.
func (p *StandardParser) Pattern(context *core.ParsingContext) *core.RegExp {
	return standardPattern
}

// Extract returns the components of the date.
func (p *StandardParser) Extract(context *core.ParsingContext, match *core.RegExpMatch) any {
	month := js.ParseInt(ToHankaku(match.Get(standardMonthGroup)), 0)
	day := js.ParseInt(ToHankaku(match.Get(standardDayGroup)), 0)
	components := context.CreateParsingComponents(core.Components{
		{Component: core.ComponentDay, Value: day},
		{Component: core.ComponentMonth, Value: month},
	})

	if match.Truthy(standardSpecialYearGroup) && core.StringPattern("同|今|本").Exec(match.Get(standardSpecialYearGroup)) != nil {
		components.Assign(core.ComponentYear, context.Reference.DateWithAdjustedTimezone().GetFullYear())
	}

	if match.Truthy(standardTypicalYearGroup) {
		yearNumText := match.Get(standardYearNumberGroup)

		var year float64
		if match.Has(standardYearNumberGroup) && yearNumText == "元" {
			year = 1
		} else {
			year = js.ParseInt(ToHankaku(yearNumText), 0)
		}
		era, eraDefined := match.Get(standardEraGroup), match.Has(standardEraGroup)
		if eraDefined && era == "令和" {
			year += 2018
		} else if eraDefined && era == "平成" {
			year += 1988
		} else if eraDefined && era == "昭和" {
			year += 1925
		}

		components.Assign(core.ComponentYear, year)
	} else {
		year := core.FindYearClosestToRef(context.RefDate, day, month)
		components.Imply(core.ComponentYear, year)
	}

	return components
}
