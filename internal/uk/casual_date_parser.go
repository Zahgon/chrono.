// Ported from src/locales/uk/parsers/UKCasualDateParser.ts.

package uk

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

// CasualDateParser is UKCasualDateParser ("сьогодні", "завтра", "позавчора", ...).
type CasualDateParser struct {
	AbstractParserWithLeftRightBoundaryChecking
}

// NewCasualDateParser returns a new UKCasualDateParser.
func NewCasualDateParser() *CasualDateParser {
	p := &CasualDateParser{}
	p.Self = p
	return p
}

// InnerPatternString returns the casual date pattern.
func (p *CasualDateParser) InnerPatternString(context *core.ParsingContext) string {
	return "(?:з|із|від)?\\s*(сьогодні|вчора|завтра|післязавтра|післяпіслязавтра|позапозавчора|позавчора)"
}

// InnerExtract returns the components of the casual date.
func (p *CasualDateParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	lowerText := js.ToLowerCase(match.Get(1))
	component := context.CreateParsingComponents(nil)

	switch lowerText {
	case "сьогодні":
		return core.Today(context.Reference)

	case "вчора":
		return core.Yesterday(context.Reference)

	case "завтра":
		return core.Tomorrow(context.Reference)

	case "післязавтра":
		return core.TheDayAfter(context.Reference, 2)

	case "післяпіслязавтра":
		return core.TheDayAfter(context.Reference, 3)

	case "позавчора":
		return core.TheDayBefore(context.Reference, 2)

	case "позапозавчора":
		return core.TheDayBefore(context.Reference, 3)
	}

	return component
}
