// Ported from src/locales/es/index.ts.

// Package es is chrono's Spanish locale: Casual and Strict parsers and the configuration
// constructors.
package es

import (
	chrono "github.com/wanasit/chrono-go"
	"github.com/wanasit/chrono-go/internal/es"
)

// Casual is the Chrono object configured for parsing casual Spanish.
var Casual *chrono.Chrono = es.Casual

// Strict is the Chrono object configured for parsing strict Spanish.
var Strict *chrono.Chrono = es.Strict

// Parse is Casual.Parse.
func Parse(text string, ref chrono.ParsingReference, option *chrono.ParsingOption) []*chrono.ParsingResult {
	return Casual.Parse(text, ref, option)
}

// ParseDate is Casual.ParseDate.
func ParseDate(text string, ref chrono.ParsingReference, option *chrono.ParsingOption) (float64, bool) {
	return Casual.ParseDate(text, ref, option)
}

// CreateCasualConfiguration creates a casual Spanish configuration (fresh parsers and refiners);
// littleEndian reads day-first slash dates (the JS default is true).
func CreateCasualConfiguration(littleEndian bool) *chrono.Configuration {
	return es.CreateCasualConfiguration(littleEndian)
}

// CreateConfiguration creates a Spanish configuration (fresh parsers and refiners); the JS
// defaults are strictMode = true and littleEndian = true.
func CreateConfiguration(strictMode, littleEndian bool) *chrono.Configuration {
	return es.CreateConfiguration(strictMode, littleEndian)
}
