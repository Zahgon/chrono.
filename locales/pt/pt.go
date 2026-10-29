// Ported from src/locales/pt/index.ts.

// Package pt is chrono's Portuguese locale: Casual and Strict parsers and the configuration
// constructors.
package pt

import (
	chrono "github.com/wanasit/chrono-go"
	"github.com/wanasit/chrono-go/internal/pt"
)

// Casual is the Chrono object configured for parsing casual Portuguese.
var Casual *chrono.Chrono = pt.Casual

// Strict is the Chrono object configured for parsing strict Portuguese.
var Strict *chrono.Chrono = pt.Strict

// Parse is Casual.Parse.
func Parse(text string, ref chrono.ParsingReference, option *chrono.ParsingOption) []*chrono.ParsingResult {
	return Casual.Parse(text, ref, option)
}

// ParseDate is Casual.ParseDate.
func ParseDate(text string, ref chrono.ParsingReference, option *chrono.ParsingOption) (float64, bool) {
	return Casual.ParseDate(text, ref, option)
}

// CreateCasualConfiguration creates a casual Portuguese configuration (fresh parsers and refiners);
// littleEndian reads day-first slash dates (the JS default is true).
func CreateCasualConfiguration(littleEndian bool) *chrono.Configuration {
	return pt.CreateCasualConfiguration(littleEndian)
}

// CreateConfiguration creates a Portuguese configuration (fresh parsers and refiners); the JS
// defaults are strictMode = true and littleEndian = true.
func CreateConfiguration(strictMode, littleEndian bool) *chrono.Configuration {
	return pt.CreateConfiguration(strictMode, littleEndian)
}
