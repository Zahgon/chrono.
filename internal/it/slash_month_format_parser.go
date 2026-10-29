// Ported from src/locales/it/parsers/ITSlashMonthFormatParser.ts.

package it

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

var slashMonthFormatPattern = core.MustRegExp("([0-9]|0[1-9]|1[012])/([0-9]{4})"+"", "i")

const (
	slashMonthMonthGroup = 1
	slashMonthYearGroup  = 2
)

// SlashMonthFormatParser is ITSlashMonthFormatParser: month/year ("06/2005").
type SlashMonthFormatParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewSlashMonthFormatParser returns a new ITSlashMonthFormatParser.
func NewSlashMonthFormatParser() *SlashMonthFormatParser {
	p := &SlashMonthFormatParser{}
	p.Self = p
	return p
}

// InnerPattern returns the month/year pattern.
func (p *SlashMonthFormatParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return slashMonthFormatPattern
}

// InnerExtract returns the components of the month.
func (p *SlashMonthFormatParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	year := js.ParseInt(match.Get(slashMonthYearGroup), 0)
	month := js.ParseInt(match.Get(slashMonthMonthGroup), 0)
	return context.CreateParsingComponents(nil).Imply(core.ComponentDay, 1).Assign(core.ComponentMonth, month).Assign(core.ComponentYear, year)
}
