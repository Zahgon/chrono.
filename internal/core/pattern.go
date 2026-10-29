// Ported from src/utils/pattern.ts.

package core

import (
	"strings"

	"github.com/wanasit/chrono-go/internal/js"
)

// DefaultTimeunitConnectorPattern is the default connectorPattern of repeatedTimeunitPattern.
const DefaultTimeunitConnectorPattern = "\\s{0,5},?\\s{0,5}"

var captureGroupOpen = MustRegExp("\\((?!\\?)", "g")
var literalDot = MustRegExp("\\.", "g")

// RepeatedTimeunitPattern is repeatedTimeunitPattern(prefix, singleTimeunitPattern,
// connectorPattern); pass DefaultTimeunitConnectorPattern for the default connector.
func RepeatedTimeunitPattern(prefix, singleTimeunitPattern, connectorPattern string) string {
	singleTimeunitPatternNoCapture := captureGroupOpen.Replace(singleTimeunitPattern, "(?:")
	return prefix + singleTimeunitPatternNoCapture + "(?:" + connectorPattern + singleTimeunitPatternNoCapture + "){0,10}"
}

// ExtractTerms returns the terms of a dictionary: the keys of an object or Map in their
// iteration order, or a copy of an array. Pass dict.Keys() for a Dictionary.
func ExtractTerms(terms []string) []string {
	return append([]string(nil), terms...)
}

// MatchAnyPattern returns `(?:term1|term2|...)` with the terms sorted by UTF-16 length, longest
// first (a stable sort: equal lengths keep the dictionary order) and every "." escaped.
func MatchAnyPattern(terms []string) string {
	sorted := ExtractTerms(terms)
	js.SortSlice(sorted, func(a, b string) float64 { return float64(js.Len(b) - js.Len(a)) })
	joinedTerms := literalDot.Replace(strings.Join(sorted, "|"), "\\.")
	return "(?:" + joinedTerms + ")"
}
