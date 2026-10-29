// Ported from src/locales/uk/index.ts.

// Package uk is chrono's Ukrainian locale: the Casual and Strict parsers and the configuration
// constructors.
package uk

import (
	chrono "github.com/wanasit/chrono-go"
	"github.com/wanasit/chrono-go/internal/uk"
)

// Casual is the Chrono object configured for parsing casual Ukrainian.
var Casual *chrono.Chrono = uk.Casual

// Strict is the Chrono object configured for parsing strict Ukrainian.
var Strict *chrono.Chrono = uk.Strict

// Parse is Casual.Parse.
func Parse(text string, ref chrono.ParsingReference, option *chrono.ParsingOption) []*chrono.ParsingResult {
	return Casual.Parse(text, ref, option)
}

// ParseDate is Casual.ParseDate.
func ParseDate(text string, ref chrono.ParsingReference, option *chrono.ParsingOption) (float64, bool) {
	return Casual.ParseDate(text, ref, option)
}

// CreateCasualConfiguration creates a casual Ukrainian configuration (fresh parsers and refiners).
func CreateCasualConfiguration() *chrono.Configuration {
	return uk.CreateCasualConfiguration()
}

// CreateConfiguration creates a Ukrainian configuration (fresh parsers and refiners); strictMode
// makes the time unit mentions strict, not casual.
func CreateConfiguration(strictMode bool) *chrono.Configuration {
	return uk.CreateConfiguration(strictMode)
}
