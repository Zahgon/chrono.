// Ported from src/locales/ru/index.ts.

// Package ru is chrono's Russian locale: the Casual and Strict parsers and the configuration
// constructors.
package ru

import (
	chrono "github.com/wanasit/chrono-go"
	"github.com/wanasit/chrono-go/internal/ru"
)

// Casual is the Chrono object configured for parsing casual Russian.
var Casual *chrono.Chrono = ru.Casual

// Strict is the Chrono object configured for parsing strict Russian.
var Strict *chrono.Chrono = ru.Strict

// Parse is Casual.Parse.
func Parse(text string, ref chrono.ParsingReference, option *chrono.ParsingOption) []*chrono.ParsingResult {
	return Casual.Parse(text, ref, option)
}

// ParseDate is Casual.ParseDate.
func ParseDate(text string, ref chrono.ParsingReference, option *chrono.ParsingOption) (float64, bool) {
	return Casual.ParseDate(text, ref, option)
}

// CreateCasualConfiguration creates a casual Russian configuration (fresh parsers and refiners).
func CreateCasualConfiguration() *chrono.Configuration {
	return ru.CreateCasualConfiguration()
}

// CreateConfiguration creates a Russian configuration (fresh parsers and refiners); strictMode
// makes the time unit mentions strict, not casual.
func CreateConfiguration(strictMode bool) *chrono.Configuration {
	return ru.CreateConfiguration(strictMode)
}
