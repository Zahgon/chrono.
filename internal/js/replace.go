package js

import "math"

// Replace is String.prototype.replace(search, replacement) with a string search value: the first
// occurrence of search (UTF-16 code unit comparison; an empty search matches at index 0) is
// replaced by GetSubstitution(replacement). With a string search there are no captures, so only
// "$$", "$&", "$`" and "$'" are special; "$1", "$<name>" and any other "$" stay literal.
// For a regular-expression search use jsregex.Regexp.Replace.
func Replace(s, search, replacement string) string {
	u, w := CodeUnits(s), CodeUnits(search)
	p := indexUnits(u, w, 0)
	if p < 0 {
		return s
	}
	return FromCodeUnits(appendSubstitution(append([]uint16(nil), u[:p]...), u, w, p, CodeUnits(replacement), u[p+len(w):]))
}

// ReplaceAll is String.prototype.replaceAll(search, replacement) with a string search value.
// An empty search matches before every code unit and at the end.
func ReplaceAll(s, search, replacement string) string {
	u, w, r := CodeUnits(s), CodeUnits(search), CodeUnits(replacement)
	adv := max(len(w), 1)
	var positions []int
	for p := indexUnits(u, w, 0); p >= 0; {
		positions = append(positions, p)
		if p+adv > len(u) {
			break
		}
		p = indexUnits(u, w, p+adv)
	}
	if len(positions) == 0 {
		return s
	}
	var out []uint16
	end := 0
	for _, p := range positions {
		out = append(out, u[end:p]...)
		out = appendSubstitution(out, u, w, p, r, nil)
		end = p + len(w)
	}
	out = append(out, u[end:]...)
	return FromCodeUnits(out)
}

// appendSubstitution appends GetSubstitution(matched=w at position p of u, no captures, no named
// captures, replacement r) and then tail.
func appendSubstitution(out, u, w []uint16, p int, r, tail []uint16) []uint16 {
	for i := 0; i < len(r); i++ {
		c := r[i]
		if c != '$' || i+1 >= len(r) {
			out = append(out, c)
			continue
		}
		switch r[i+1] {
		case '$':
			out = append(out, '$')
			i++
		case '&':
			out = append(out, w...)
			i++
		case '`':
			out = append(out, u[:p]...)
			i++
		case '\'':
			out = append(out, u[min(p+len(w), len(u)):]...)
			i++
		default:
			out = append(out, c)
		}
	}
	return append(out, tail...)
}

// IsNaN is the global isNaN applied to a number (isNaN(x) for a string is IsNaN(ToNumber(s))).
func IsNaN(x float64) bool { return math.IsNaN(x) }

// IsFinite is the global isFinite applied to a number.
func IsFinite(x float64) bool { return !math.IsNaN(x) && !math.IsInf(x, 0) }
