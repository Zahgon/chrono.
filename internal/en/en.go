// Ported from src/locales/en/index.ts.

// Package en is chrono's English locale (parsers, refiners and configuration). The public
// packages github.com/wanasit/chrono-go and github.com/wanasit/chrono-go/locales/en expose it.
package en

import "github.com/wanasit/chrono-go/internal/core"

func init() {
	core.DefaultConfiguration = func() *core.Configuration { return CreateCasualConfiguration(false) }
}

// Casual is the Chrono object configured for parsing casual English.
var Casual = core.NewChrono(CreateCasualConfiguration(false))

// Strict is the Chrono object configured for parsing strict English.
var Strict = core.NewChrono(CreateConfiguration(true, false))

// GB is the Chrono object configured for parsing UK-style English.
var GB = core.NewChrono(CreateCasualConfiguration(true))

// Parse is Casual.Parse.
func Parse(text string, ref core.ParsingReference, option *core.ParsingOption) []*core.ParsingResult {
	return Casual.Parse(text, ref, option)
}

// ParseDate is Casual.ParseDate.
func ParseDate(text string, ref core.ParsingReference, option *core.ParsingOption) (float64, bool) {
	return Casual.ParseDate(text, ref, option)
}
