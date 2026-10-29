// Ported from src/chrono.ts.

package core

import (
	"reflect"

	"github.com/wanasit/chrono-go/internal/js"
)

// Parser is chrono's Parser: Pattern returns the regular expression run over the text and Extract
// turns each match into a result. Extract returns nil (null), *ParsingComponents,
// *ParsingResult or Components (a nil pointer or a nil Components is null).
type Parser interface {
	Pattern(context *ParsingContext) *RegExp
	Extract(context *ParsingContext, match *RegExpMatch) any
}

// Refiner is chrono's Refiner: it takes the list of results and returns another list.
type Refiner interface {
	Refine(context *ParsingContext, results []*ParsingResult) []*ParsingResult
}

// Configuration is an ordered list of parsers and refiners.
type Configuration struct {
	Parsers  []Parser
	Refiners []Refiner
}

// Chrono is the Chrono object: its parsers and refiners are applied in order.
type Chrono struct {
	Parsers  []Parser
	Refiners []Refiner
}

// DefaultConfiguration creates the configuration used by NewChrono(nil) (English casual); it is
// set by the English locale package.
var DefaultConfiguration func() *Configuration

// NewChrono is `new Chrono(configuration)`; nil uses DefaultConfiguration.
func NewChrono(configuration *Configuration) *Chrono {
	if configuration == nil {
		configuration = DefaultConfiguration()
	}
	return &Chrono{
		Parsers:  append([]Parser{}, configuration.Parsers...),
		Refiners: append([]Refiner{}, configuration.Refiners...),
	}
}

// Clone creates a shallow copy of the Chrono object with the same parsers and refiners.
func (c *Chrono) Clone() *Chrono {
	return NewChrono(&Configuration{
		Parsers:  append([]Parser{}, c.Parsers...),
		Refiners: append([]Refiner{}, c.Refiners...),
	})
}

// ParseDate returns the start date (time value) of the first result; ok is false when there is
// no result (null).
func (c *Chrono) ParseDate(text string, ref ParsingReference, option *ParsingOption) (float64, bool) {
	results := c.Parse(text, ref, option)
	if len(results) > 0 {
		return results[0].Start.Date(), true
	}
	return 0, false
}

// Parse parses the text with the reference and option (nil: absent).
func (c *Chrono) Parse(text string, ref ParsingReference, option *ParsingOption) []*ParsingResult {
	context := NewParsingContext(text, ref, option)

	results := []*ParsingResult{}
	for _, parser := range c.Parsers {
		parsedResults := executeParser(context, parser)
		results = append(results, parsedResults...)
	}

	js.SortSlice(results, func(a, b *ParsingResult) float64 { return float64(a.Index - b.Index) })

	for _, refiner := range c.Refiners {
		results = refiner.Refine(context, results)
	}
	return results
}

// isNullResult reports whether an Extract result is falsy (null, undefined or a nil pointer).
func isNullResult(result any) bool {
	if result == nil {
		return true
	}
	switch v := result.(type) {
	case *ParsingComponents:
		return v == nil
	case *ParsingResult:
		return v == nil
	case Components:
		return v == nil
	}
	rv := reflect.ValueOf(result)
	switch rv.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Interface:
		return rv.IsNil()
	}
	return false
}

func executeParser(context *ParsingContext, parser Parser) []*ParsingResult {
	results := []*ParsingResult{}
	pattern := parser.Pattern(context)

	// The original's string slicing, on the UTF-16 code units of the text.
	units := context.textUnits()
	remainingStart := 0
	remainingText := context.Text
	substringFrom := func(start int) {
		start = max(0, min(start, len(units)))
		remainingStart = start
		remainingText = js.FromCodeUnits(units[start:])
	}
	match := pattern.Exec(remainingText)

	for match != nil {
		// Calculate match index on the full text;
		index := match.Index + len(units) - (len(units) - remainingStart)
		match.Index = index

		result := parser.Extract(context, match)
		if isNullResult(result) {
			// If fails, move on by 1
			substringFrom(match.Index + 1)
			match = pattern.Exec(remainingText)
			continue
		}

		var parsedResult *ParsingResult
		switch r := result.(type) {
		case *ParsingResult:
			parsedResult = r
		case *ParsingComponents:
			parsedResult = context.CreateParsingResult(match.Index, match.Get(0), nil, nil)
			parsedResult.Start = r
		case Components:
			parsedResult = context.CreateParsingResult(match.Index, match.Get(0), context.CreateParsingComponents(r), nil)
		default:
			// Any other truthy value is used as an (empty) components object.
			parsedResult = context.CreateParsingResult(match.Index, match.Get(0), context.CreateParsingComponents(Components{}), nil)
		}

		parsedIndex := parsedResult.Index
		parsedText := parsedResult.Text

		results = append(results, parsedResult)
		substringFrom(parsedIndex + js.Len(parsedText))
		match = pattern.Exec(remainingText)
	}
	return results
}

// ParsingContext is the context of one parse call.
type ParsingContext struct {
	// Text is the input text.
	Text string
	// Option is `option ?? {}`.
	Option ParsingOption
	// Reference is the reference instant and time zone.
	Reference *ReferenceWithTimezone
	// RefDate is the reference instant's time value (deprecated in the original; use Reference).
	RefDate float64

	unitsText string   // the Text that units were computed from
	units     []uint16 // UTF-16 code units of unitsText
}

// textUnits returns the UTF-16 code units of Text (cached).
func (c *ParsingContext) textUnits() []uint16 {
	if c.units == nil || c.unitsText != c.Text {
		c.unitsText = c.Text
		c.units = js.CodeUnits(c.Text)
	}
	return c.units
}

// TextLength is `context.text.length`.
func (c *ParsingContext) TextLength() int { return len(c.textUnits()) }

// Substring is `context.text.substring(a, b)`.
func (c *ParsingContext) Substring(a, b int) string {
	u := c.textUnits()
	a, b = max(0, min(a, len(u))), max(0, min(b, len(u)))
	if a > b {
		a, b = b, a
	}
	return js.FromCodeUnits(u[a:b])
}

// SubstringFrom is `context.text.substring(a)`.
func (c *ParsingContext) SubstringFrom(a int) string { return c.Substring(a, len(c.textUnits())) }

// NewParsingContext is `new ParsingContext(text, refDate, option)`.
func NewParsingContext(text string, ref ParsingReference, option *ParsingOption) *ParsingContext {
	c := &ParsingContext{Text: text}
	if option != nil {
		c.Option = *option
	}
	c.Reference = ReferenceFromInput(ref, c.Option.Timezones)
	c.RefDate = c.Reference.Instant
	return c
}

// CreateParsingComponents is context.createParsingComponents(components); nil components is an
// absent argument.
func (c *ParsingContext) CreateParsingComponents(components Components) *ParsingComponents {
	return NewParsingComponents(c.Reference, components)
}

// CreateParsingResult is context.createParsingResult(index, text, start, end); nil start or end
// is absent (start then defaults to new components).
func (c *ParsingContext) CreateParsingResult(index int, text string, start, end *ParsingComponents) *ParsingResult {
	return NewParsingResult(c.Reference, index, text, start, end)
}

// CreateParsingResultToIndex is context.createParsingResult(index, endIndex, start, end): the
// text is `context.text.substring(index, endIndex)`.
func (c *ParsingContext) CreateParsingResultToIndex(index, endIndex int, start, end *ParsingComponents) *ParsingResult {
	return c.CreateParsingResult(index, c.Substring(index, endIndex), start, end)
}
