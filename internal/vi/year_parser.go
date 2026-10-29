// Ported from src/locales/vi/parsers/VIYearParser.ts.

package vi

import "github.com/wanasit/chrono-go/internal/core"

// năm 1975  |  năm 43 TCN  |  179 TCN (bare BC year without năm prefix)
var yearPattern = core.MustRegExp("(?:\\bnăm\\s*("+YearPattern+")|\\b([0-9]{1,4})\\s*(TCN))(?=\\W|$)", "i")

const (
	yearWithNamGroup  = 1 // năm YYYY or năm YYYY TCN
	bareBCYearGroup   = 2 // bare YYYY in "YYYY TCN"
	bareBCSuffixGroup = 3 // "TCN" in bare form
)

// YearParser is VIYearParser ("năm 1975", "năm 43 TCN", "179 TCN").
type YearParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewYearParser returns a new VIYearParser.
func NewYearParser() *YearParser {
	p := &YearParser{}
	p.Self = p
	return p
}

// InnerPattern returns the year pattern.
func (p *YearParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return yearPattern
}

// jsString is String(v) for a regex group: "undefined" when the group did not participate.
func jsString(match *core.RegExpMatch, i int) string {
	if !match.Has(i) {
		return "undefined"
	}
	return match.Get(i)
}

// InnerExtract returns the result of the year.
func (p *YearParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	var yearText string
	if match.Truthy(yearWithNamGroup) {
		yearText = match.Get(yearWithNamGroup)
	} else {
		yearText = jsString(match, bareBCYearGroup) + " " + jsString(match, bareBCSuffixGroup)
	}
	result := context.CreateParsingResult(match.Index, match.Get(0), nil, nil)
	result.Start.Assign(core.ComponentYear, ParseYear(yearText))
	result.Start.Imply(core.ComponentMonth, 1)
	result.Start.Imply(core.ComponentDay, 1)
	return result
}
