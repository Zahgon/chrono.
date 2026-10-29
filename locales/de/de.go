// Ported from src/locales/de/index.ts.

// Package de is chrono's German locale: the Casual and Strict parsers and the configuration
// constructors.
package de

import (
	chrono "github.com/wanasit/chrono-go"
	"github.com/wanasit/chrono-go/internal/de"
)

// Casual is the Chrono object configured for parsing casual German
// (createCasualConfiguration(true)).
var Casual *chrono.Chrono = de.Casual

// Strict is the Chrono object configured for parsing strict German
// (createConfiguration(true, true)).
var Strict *chrono.Chrono = de.Strict

// Parse is Casual.Parse.
func Parse(text string, ref chrono.ParsingReference, option *chrono.ParsingOption) []*chrono.ParsingResult {
	return Casual.Parse(text, ref, option)
}

// ParseDate is Casual.ParseDate.
func ParseDate(text string, ref chrono.ParsingReference, option *chrono.ParsingOption) (float64, bool) {
	return Casual.ParseDate(text, ref, option)
}

// CreateCasualConfiguration creates a casual German configuration (fresh parsers and refiners).
func CreateCasualConfiguration(littleEndian bool) *chrono.Configuration {
	return de.CreateCasualConfiguration(littleEndian)
}

// CreateConfiguration creates a German configuration (fresh parsers and refiners).
func CreateConfiguration(strictMode, littleEndian bool) *chrono.Configuration {
	return de.CreateConfiguration(strictMode, littleEndian)
}
