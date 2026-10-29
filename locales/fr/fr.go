// Ported from src/locales/fr/index.ts.

// Package fr is chrono's French locale: Casual and Strict parsers and the configuration
// constructors.
package fr

import (
	chrono "github.com/wanasit/chrono-go"
	"github.com/wanasit/chrono-go/internal/fr"
)

// Casual is the Chrono object configured for parsing casual French.
var Casual *chrono.Chrono = fr.Casual

// Strict is the Chrono object configured for parsing strict French.
var Strict *chrono.Chrono = fr.Strict

// Parse is Casual.Parse.
func Parse(text string, ref chrono.ParsingReference, option *chrono.ParsingOption) []*chrono.ParsingResult {
	return Casual.Parse(text, ref, option)
}

// ParseDate is Casual.ParseDate.
func ParseDate(text string, ref chrono.ParsingReference, option *chrono.ParsingOption) (float64, bool) {
	return Casual.ParseDate(text, ref, option)
}

// CreateCasualConfiguration creates a casual French configuration (fresh parsers and refiners);
// littleEndian reads day-first slash dates (the JS default is true).
func CreateCasualConfiguration(littleEndian bool) *chrono.Configuration {
	return fr.CreateCasualConfiguration(littleEndian)
}

// CreateConfiguration creates a French configuration (fresh parsers and refiners); the JS
// defaults are strictMode = true and littleEndian = true.
func CreateConfiguration(strictMode, littleEndian bool) *chrono.Configuration {
	return fr.CreateConfiguration(strictMode, littleEndian)
}
