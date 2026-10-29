// Ported from src/locales/fi/index.ts.

// Package fi is chrono's Finnish locale: the Casual and Strict parsers and the configuration
// constructors.
package fi

import (
	chrono "github.com/wanasit/chrono-go"
	"github.com/wanasit/chrono-go/internal/fi"
)

// Casual is the Chrono object configured for parsing casual Finnish.
var Casual *chrono.Chrono = fi.Casual

// Strict is the Chrono object configured for parsing strict Finnish.
var Strict *chrono.Chrono = fi.Strict

// Parse is Casual.Parse.
func Parse(text string, ref chrono.ParsingReference, option *chrono.ParsingOption) []*chrono.ParsingResult {
	return Casual.Parse(text, ref, option)
}

// ParseDate is Casual.ParseDate.
func ParseDate(text string, ref chrono.ParsingReference, option *chrono.ParsingOption) (float64, bool) {
	return Casual.ParseDate(text, ref, option)
}

// CreateCasualConfiguration creates a casual Finnish configuration (fresh parsers and refiners).
// littleEndian reads day-first slash dates (the original's default is true).
func CreateCasualConfiguration(littleEndian bool) *chrono.Configuration {
	return fi.CreateCasualConfiguration(littleEndian)
}

// CreateConfiguration creates a Finnish configuration (fresh parsers and refiners). The original's
// defaults are strictMode = true and littleEndian = true.
func CreateConfiguration(strictMode, littleEndian bool) *chrono.Configuration {
	return fi.CreateConfiguration(strictMode, littleEndian)
}
