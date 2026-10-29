// Ported from src/locales/ru/parsers/RUCasualDateParser.ts.

package ru

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
)

// CasualDateParser is RUCasualDateParser ("сегодня", "завтра", "позавчера", ...).
type CasualDateParser struct {
	AbstractParserWithLeftRightBoundaryChecking
}

// NewCasualDateParser returns a new RUCasualDateParser.
func NewCasualDateParser() *CasualDateParser {
	p := &CasualDateParser{}
	p.Self = p
	return p
}

// InnerPatternString returns the casual date pattern.
func (p *CasualDateParser) InnerPatternString(context *core.ParsingContext) string {
	return "(?:с|со)?\\s*(сегодня|вчера|завтра|послезавтра|послепослезавтра|позапозавчера|позавчера)"
}

// InnerExtract returns the components of the casual date.
func (p *CasualDateParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	lowerText := js.ToLowerCase(match.Get(1))
	component := context.CreateParsingComponents(nil)

	switch lowerText {
	case "сегодня":
		return core.Today(context.Reference)

	case "вчера":
		return core.Yesterday(context.Reference)

	case "завтра":
		return core.Tomorrow(context.Reference)

	case "послезавтра":
		return core.TheDayAfter(context.Reference, 2)

	case "послепослезавтра":
		return core.TheDayAfter(context.Reference, 3)

	case "позавчера":
		return core.TheDayBefore(context.Reference, 2)

	case "позапозавчера":
		return core.TheDayBefore(context.Reference, 3)
	}

	return component
}
