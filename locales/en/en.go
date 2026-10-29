// Ported from src/locales/en/index.ts.

// Package en is chrono's English locale: Casual, Strict and GB (UK-style, day-first) parsers and
// the configuration constructors.
package en

import (
	chrono "github.com/wanasit/chrono-go"
	"github.com/wanasit/chrono-go/internal/en"
)

// Casual is the Chrono object configured for parsing casual English (the same object as
// chrono.Casual).
var Casual *chrono.Chrono = en.Casual

// Strict is the Chrono object configured for parsing strict English (the same object as
// chrono.Strict).
var Strict *chrono.Chrono = en.Strict

// GB is the Chrono object configured for parsing UK-style English.
var GB *chrono.Chrono = en.GB

// Parse is Casual.Parse.
func Parse(text string, ref chrono.ParsingReference, option *chrono.ParsingOption) []*chrono.ParsingResult {
	return Casual.Parse(text, ref, option)
}

// ParseDate is Casual.ParseDate.
func ParseDate(text string, ref chrono.ParsingReference, option *chrono.ParsingOption) (float64, bool) {
	return Casual.ParseDate(text, ref, option)
}

// CreateCasualConfiguration creates a casual English configuration (fresh parsers and refiners).
func CreateCasualConfiguration(littleEndian bool) *chrono.Configuration {
	return en.CreateCasualConfiguration(littleEndian)
}

// CreateConfiguration creates an English configuration (fresh parsers and refiners).
func CreateConfiguration(strictMode, littleEndian bool) *chrono.Configuration {
	return en.CreateConfiguration(strictMode, littleEndian)
}
