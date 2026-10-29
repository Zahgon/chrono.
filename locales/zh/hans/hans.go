// Ported from src/locales/zh/hans/index.ts.

// Package hans is chrono's Simplified Chinese locale: the Casual and Strict parsers and the configuration
// constructors.
package hans

import (
	chrono "github.com/wanasit/chrono-go"
	impl "github.com/wanasit/chrono-go/internal/zh/hans"
)

// Casual is the Chrono object configured for parsing casual Simplified Chinese.
var Casual *chrono.Chrono = impl.Casual

// Strict is the Chrono object configured for parsing strict Simplified Chinese.
var Strict *chrono.Chrono = impl.Strict

// Parse is Casual.Parse.
func Parse(text string, ref chrono.ParsingReference, option *chrono.ParsingOption) []*chrono.ParsingResult {
	return Casual.Parse(text, ref, option)
}

// ParseDate is Casual.ParseDate.
func ParseDate(text string, ref chrono.ParsingReference, option *chrono.ParsingOption) (float64, bool) {
	return Casual.ParseDate(text, ref, option)
}

// CreateCasualConfiguration creates a casual Simplified Chinese configuration (fresh parsers and refiners).
func CreateCasualConfiguration() *chrono.Configuration {
	return impl.CreateCasualConfiguration()
}

// CreateConfiguration creates a Simplified Chinese configuration (fresh parsers and refiners).
func CreateConfiguration() *chrono.Configuration {
	return impl.CreateConfiguration()
}
