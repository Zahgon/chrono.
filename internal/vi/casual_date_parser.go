// Ported from src/locales/vi/parsers/VICasualDateParser.ts.

package vi

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

var casualDatePattern = core.MustRegExp("\\b(hôm nay|hôm qua|hôm kia|ngày mai|ngày kia|bây giờ|lúc này)(?=\\W|$)", "i")

// CasualDateParser is VICasualDateParser ("hôm nay", "ngày mai", "bây giờ", ...).
type CasualDateParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewCasualDateParser returns a new VICasualDateParser.
func NewCasualDateParser() *CasualDateParser {
	p := &CasualDateParser{}
	p.Self = p
	return p
}

// InnerPattern returns the casual date pattern.
func (p *CasualDateParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return casualDatePattern
}

// InnerExtract returns the components of the casual date.
func (p *CasualDateParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	switch js.ToLowerCase(match.Get(1)) {
	case "bây giờ", "lúc này":
		return core.Now(context.Reference)
	case "hôm nay":
		return core.Today(context.Reference)
	case "hôm qua":
		return core.Yesterday(context.Reference)
	case "hôm kia":
		return core.TheDayBefore(context.Reference, 2)
	case "ngày mai":
		return core.Tomorrow(context.Reference)
	case "ngày kia":
		return core.TheDayAfter(context.Reference, 2)
	}
	return context.CreateParsingComponents(nil)
}
