// Ported from src/locales/nl/index.ts.

// Package nl is chrono's Dutch locale: the Casual and Strict parsers and the configuration
// constructors.
package nl

import (
	chrono "github.com/wanasit/chrono-go"
	"github.com/wanasit/chrono-go/internal/nl"
)

// Casual is the Chrono object configured for parsing casual Dutch
// (createCasualConfiguration(true)).
var Casual *chrono.Chrono = nl.Casual

// Strict is the Chrono object configured for parsing strict Dutch
// (createConfiguration(true, true)).
var Strict *chrono.Chrono = nl.Strict

// Parse is Casual.Parse.
func Parse(text string, ref chrono.ParsingReference, option *chrono.ParsingOption) []*chrono.ParsingResult {
	return Casual.Parse(text, ref, option)
}

// ParseDate is Casual.ParseDate.
func ParseDate(text string, ref chrono.ParsingReference, option *chrono.ParsingOption) (float64, bool) {
	return Casual.ParseDate(text, ref, option)
}

// CreateCasualConfiguration creates a casual Dutch configuration (fresh parsers and refiners).
func CreateCasualConfiguration(littleEndian bool) *chrono.Configuration {
	return nl.CreateCasualConfiguration(littleEndian)
}

// CreateConfiguration creates a Dutch configuration (fresh parsers and refiners).
func CreateConfiguration(strictMode, littleEndian bool) *chrono.Configuration {
	return nl.CreateConfiguration(strictMode, littleEndian)
}
