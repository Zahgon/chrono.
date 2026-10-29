// Ported from src/locales/sv/index.ts.

// Package sv is chrono's Swedish locale: the Casual and Strict parsers and the configuration
// constructors.
package sv

import (
	chrono "github.com/wanasit/chrono-go"
	"github.com/wanasit/chrono-go/internal/sv"
)

// Casual is the Chrono object configured for parsing casual Swedish
// (createCasualConfiguration(true)).
var Casual *chrono.Chrono = sv.Casual

// Strict is the Chrono object configured for parsing strict Swedish
// (createConfiguration(true, true)).
var Strict *chrono.Chrono = sv.Strict

// Parse is Casual.Parse.
func Parse(text string, ref chrono.ParsingReference, option *chrono.ParsingOption) []*chrono.ParsingResult {
	return Casual.Parse(text, ref, option)
}

// ParseDate is Casual.ParseDate.
func ParseDate(text string, ref chrono.ParsingReference, option *chrono.ParsingOption) (float64, bool) {
	return Casual.ParseDate(text, ref, option)
}

// CreateCasualConfiguration creates a casual Swedish configuration (fresh parsers and refiners).
func CreateCasualConfiguration(littleEndian bool) *chrono.Configuration {
	return sv.CreateCasualConfiguration(littleEndian)
}

// CreateConfiguration creates a Swedish configuration (fresh parsers and refiners).
func CreateConfiguration(strictMode, littleEndian bool) *chrono.Configuration {
	return sv.CreateConfiguration(strictMode, littleEndian)
}
