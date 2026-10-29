// Ported from src/locales/zh/index.ts.

// Package zh is chrono's Chinese locale: the Casual and Strict parsers and the configuration
// constructors.
package zh

import (
	chrono "github.com/wanasit/chrono-go"
	impl "github.com/wanasit/chrono-go/internal/zh"
)

// Casual is the Chrono object configured for parsing casual Chinese.
var Casual *chrono.Chrono = impl.Casual

// Strict is the Chrono object configured for parsing strict Chinese.
var Strict *chrono.Chrono = impl.Strict

// Parse is Casual.Parse.
func Parse(text string, ref chrono.ParsingReference, option *chrono.ParsingOption) []*chrono.ParsingResult {
	return Casual.Parse(text, ref, option)
}

// ParseDate is Casual.ParseDate.
func ParseDate(text string, ref chrono.ParsingReference, option *chrono.ParsingOption) (float64, bool) {
	return Casual.ParseDate(text, ref, option)
}

// CreateCasualConfiguration creates a casual Chinese configuration (fresh parsers and refiners).
func CreateCasualConfiguration() *chrono.Configuration {
	return impl.CreateCasualConfiguration()
}

// CreateConfiguration creates a Chinese configuration (fresh parsers and refiners).
func CreateConfiguration() *chrono.Configuration {
	return impl.CreateConfiguration()
}
