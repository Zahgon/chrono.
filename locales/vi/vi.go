// Ported from src/locales/vi/index.ts.

// Package vi is chrono's Vietnamese locale: the Casual and Strict parsers and the configuration
// constructors.
package vi

import (
	chrono "github.com/wanasit/chrono-go"
	"github.com/wanasit/chrono-go/internal/vi"
)

// Casual is the Chrono object configured for parsing casual Vietnamese.
var Casual *chrono.Chrono = vi.Casual

// Strict is the Chrono object configured for parsing strict Vietnamese.
var Strict *chrono.Chrono = vi.Strict

// Parse is Casual.Parse.
func Parse(text string, ref chrono.ParsingReference, option *chrono.ParsingOption) []*chrono.ParsingResult {
	return Casual.Parse(text, ref, option)
}

// ParseDate is Casual.ParseDate.
func ParseDate(text string, ref chrono.ParsingReference, option *chrono.ParsingOption) (float64, bool) {
	return Casual.ParseDate(text, ref, option)
}

// CreateCasualConfiguration creates a casual Vietnamese configuration (fresh parsers and refiners).
// littleEndian reads day-first slash dates (the original's default is true).
func CreateCasualConfiguration(littleEndian bool) *chrono.Configuration {
	return vi.CreateCasualConfiguration(littleEndian)
}

// CreateConfiguration creates a Vietnamese configuration (fresh parsers and refiners). The original's
// defaults are strictMode = true and littleEndian = true.
func CreateConfiguration(strictMode, littleEndian bool) *chrono.Configuration {
	return vi.CreateConfiguration(strictMode, littleEndian)
}
