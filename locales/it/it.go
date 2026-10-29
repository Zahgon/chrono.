// Ported from src/locales/it/index.ts.

// Package it is chrono's Italian locale: Casual and Strict parsers and the configuration
// constructors (ITDefaultConfiguration's createCasualConfiguration and createConfiguration).
package it

import (
	chrono "github.com/wanasit/chrono-go"
	"github.com/wanasit/chrono-go/internal/it"
)

// Casual is the Chrono object configured for parsing casual Italian.
var Casual *chrono.Chrono = it.Casual

// Strict is the Chrono object configured for parsing strict Italian.
var Strict *chrono.Chrono = it.Strict

// Parse is Casual.Parse.
func Parse(text string, ref chrono.ParsingReference, option *chrono.ParsingOption) []*chrono.ParsingResult {
	return Casual.Parse(text, ref, option)
}

// ParseDate is Casual.ParseDate.
func ParseDate(text string, ref chrono.ParsingReference, option *chrono.ParsingOption) (float64, bool) {
	return Casual.ParseDate(text, ref, option)
}

// CreateCasualConfiguration creates a casual Italian configuration (fresh parsers and refiners).
func CreateCasualConfiguration() *chrono.Configuration {
	return it.CreateCasualConfiguration()
}

// CreateConfiguration creates an Italian configuration (fresh parsers and refiners); strictMode
// makes the time unit mentions strict (the JS default is true).
func CreateConfiguration(strictMode bool) *chrono.Configuration {
	return it.CreateConfiguration(strictMode)
}
