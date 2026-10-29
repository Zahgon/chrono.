// Ported from src/locales/ja/index.ts.

// Package ja is chrono's Japanese locale: the Casual and Strict parsers and the configuration
// constructors.
package ja

import (
	chrono "github.com/wanasit/chrono-go"
	"github.com/wanasit/chrono-go/internal/ja"
)

// Casual is the Chrono object configured for parsing casual Japanese.
var Casual *chrono.Chrono = ja.Casual

// Strict is the Chrono object configured for parsing strict Japanese.
var Strict *chrono.Chrono = ja.Strict

// Parse is Casual.Parse.
func Parse(text string, ref chrono.ParsingReference, option *chrono.ParsingOption) []*chrono.ParsingResult {
	return Casual.Parse(text, ref, option)
}

// ParseDate is Casual.ParseDate.
func ParseDate(text string, ref chrono.ParsingReference, option *chrono.ParsingOption) (float64, bool) {
	return Casual.ParseDate(text, ref, option)
}

// CreateCasualConfiguration creates a casual Japanese configuration (fresh parsers and
// refiners).
func CreateCasualConfiguration() *chrono.Configuration {
	return ja.CreateCasualConfiguration()
}

// CreateConfiguration creates a Japanese configuration (fresh parsers and refiners).
func CreateConfiguration(strictMode bool) *chrono.Configuration {
	return ja.CreateConfiguration(strictMode)
}
