package js

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

// relative resolves a relative index the way String.prototype.slice does: negative values count
// from the end, and the result is clamped to [0, n].
func relative(i, n int) int {
	if i < 0 {
		i += n
		if i < 0 {
			return 0
		}
	}
	if i > n {
		return n
	}
	return i
}

func clamp(i, n int) int {
	if i < 0 {
		return 0
	}
	if i > n {
		return n
	}
	return i
}

// CharAt is String.prototype.charAt: the code unit at UTF-16 index i, or "" when out of range.
func CharAt(s string, i int) string {
	if isASCII(s) {
		if i < 0 || i >= len(s) {
			return ""
		}
		return s[i : i+1]
	}
	u := CodeUnits(s)
	if i < 0 || i >= len(u) {
		return ""
	}
	return FromCodeUnits(u[i : i+1])
}

// CharCodeAt is String.prototype.charCodeAt; it returns -1 where JavaScript returns NaN.
func CharCodeAt(s string, i int) int {
	if isASCII(s) {
		if i < 0 || i >= len(s) {
			return -1
		}
		return int(s[i])
	}
	u := CodeUnits(s)
	if i < 0 || i >= len(u) {
		return -1
	}
	return int(u[i])
}

// CodePointAt is String.prototype.codePointAt; it returns -1 where JavaScript returns undefined.
func CodePointAt(s string, i int) int {
	u := CodeUnits(s)
	if i < 0 || i >= len(u) {
		return -1
	}
	c := int(u[i])
	if c >= 0xd800 && c < 0xdc00 && i+1 < len(u) && u[i+1] >= 0xdc00 && u[i+1] < 0xe000 {
		return (c-0xd800)<<10 + (int(u[i+1]) - 0xdc00) + 0x10000
	}
	return c
}

// Slice is String.prototype.slice(start, end); negative indices count from the end.
func Slice(s string, start, end int) string {
	if isASCII(s) {
		from, to := relative(start, len(s)), relative(end, len(s))
		if from >= to {
			return ""
		}
		return s[from:to]
	}
	u := CodeUnits(s)
	return unitsString(u, relative(start, len(u)), relative(end, len(u)))
}

// SliceFrom is String.prototype.slice(start).
func SliceFrom(s string, start int) string {
	if isASCII(s) {
		return s[relative(start, len(s)):]
	}
	u := CodeUnits(s)
	return unitsString(u, relative(start, len(u)), len(u))
}

// Substring is String.prototype.substring(a, b): negative indices become 0 and the bounds are
// swapped when a > b.
func Substring(s string, a, b int) string {
	n := Len(s)
	a, b = clamp(a, n), clamp(b, n)
	if a > b {
		a, b = b, a
	}
	if isASCII(s) {
		return s[a:b]
	}
	return unitsString(CodeUnits(s), a, b)
}

// SubstringFrom is String.prototype.substring(a).
func SubstringFrom(s string, a int) string { return Substring(s, a, Len(s)) }

// Substr is String.prototype.substr(start, length).
func Substr(s string, start, length int) string {
	n := Len(s)
	from := relative(start, n)
	if length > n-from {
		length = n - from
	}
	if length <= 0 {
		return ""
	}
	if isASCII(s) {
		return s[from : from+length]
	}
	return unitsString(CodeUnits(s), from, from+length)
}

// SubstrFrom is String.prototype.substr(start).
func SubstrFrom(s string, start int) string { return SliceFrom(s, start) }

// plainUTF8 reports whether byte-level searching in s gives the same matches as code unit level
// searching, which holds when s contains no lone-surrogate byte forms.
func plainUTF8(s string) bool { return !hasSurrogateBytes(s) }

func indexUnits(u, w []uint16, from int) int {
	for i := from; i+len(w) <= len(u); i++ {
		match := true
		for j := range w {
			if u[i+j] != w[j] {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}

func lastIndexUnits(u, w []uint16, from int) int {
	if from > len(u)-len(w) {
		from = len(u) - len(w)
	}
	for i := from; i >= 0; i-- {
		match := true
		for j := range w {
			if u[i+j] != w[j] {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}

// IndexOf is String.prototype.indexOf(search, from), in UTF-16 code units.
func IndexOf(s, search string, from int) int {
	if isASCII(s) && isASCII(search) {
		from = clamp(from, len(s))
		if i := strings.Index(s[from:], search); i >= 0 {
			return from + i
		}
		return -1
	}
	u := CodeUnits(s)
	return indexUnits(u, CodeUnits(search), clamp(from, len(u)))
}

// LastIndexOf is String.prototype.lastIndexOf(search).
func LastIndexOf(s, search string) int {
	if isASCII(s) && isASCII(search) {
		return strings.LastIndex(s, search)
	}
	u := CodeUnits(s)
	return lastIndexUnits(u, CodeUnits(search), len(u))
}

// LastIndexOfFrom is String.prototype.lastIndexOf(search, from).
func LastIndexOfFrom(s, search string, from int) int {
	u := CodeUnits(s)
	return lastIndexUnits(u, CodeUnits(search), clamp(from, len(u)))
}

// Includes is String.prototype.includes(search).
func Includes(s, search string) bool {
	if plainUTF8(s) && plainUTF8(search) {
		return strings.Contains(s, search)
	}
	return IndexOf(s, search, 0) >= 0
}

// IncludesFrom is String.prototype.includes(search, from).
func IncludesFrom(s, search string, from int) bool { return IndexOf(s, search, from) >= 0 }

// StartsWith is String.prototype.startsWith(search, pos).
func StartsWith(s, search string, pos int) bool {
	if pos <= 0 && plainUTF8(s) && plainUTF8(search) {
		return strings.HasPrefix(s, search)
	}
	u, w := CodeUnits(s), CodeUnits(search)
	pos = clamp(pos, len(u))
	if pos+len(w) > len(u) {
		return false
	}
	for j := range w {
		if u[pos+j] != w[j] {
			return false
		}
	}
	return true
}

// EndsWith is String.prototype.endsWith(search).
func EndsWith(s, search string) bool {
	if plainUTF8(s) && plainUTF8(search) {
		return strings.HasSuffix(s, search)
	}
	return EndsWithAt(s, search, Len(s))
}

// EndsWithAt is String.prototype.endsWith(search, endPosition).
func EndsWithAt(s, search string, end int) bool {
	u, w := CodeUnits(s), CodeUnits(search)
	end = clamp(end, len(u))
	start := end - len(w)
	if start < 0 {
		return false
	}
	for j := range w {
		if u[start+j] != w[j] {
			return false
		}
	}
	return true
}

// Split is String.prototype.split(sep, limit) with a string separator. An empty separator splits
// into single code units (so a surrogate pair becomes two lone surrogates). A negative limit means
// no limit.
func Split(s, sep string, limit int) []string {
	if limit == 0 {
		return []string{}
	}
	var parts []string
	switch {
	case sep == "":
		if isASCII(s) {
			parts = make([]string, len(s))
			for i := range parts {
				parts[i] = s[i : i+1]
			}
		} else {
			u := CodeUnits(s)
			parts = make([]string, len(u))
			for i := range u {
				parts[i] = FromCodeUnits(u[i : i+1])
			}
		}
	case plainUTF8(s) && plainUTF8(sep):
		parts = strings.Split(s, sep)
	default:
		u, w := CodeUnits(s), CodeUnits(sep)
		p := 0
		for {
			q := indexUnits(u, w, p)
			if q < 0 {
				break
			}
			parts = append(parts, FromCodeUnits(u[p:q]))
			p = q + len(w)
		}
		parts = append(parts, FromCodeUnits(u[p:]))
	}
	if limit > 0 && len(parts) > limit {
		parts = parts[:limit]
	}
	return parts
}

// IsWhiteSpace reports whether r is JavaScript WhiteSpace or LineTerminator, the set that
// String.prototype.trim removes.
func IsWhiteSpace(r rune) bool {
	switch r {
	case '\t', '\n', '\v', '\f', '\r', ' ', 0xa0, 0x1680, 0x2028, 0x2029, 0x202f, 0x205f, 0x3000, 0xfeff:
		return true
	}
	return r >= 0x2000 && r <= 0x200a
}

// Trim is String.prototype.trim.
func Trim(s string) string { return strings.TrimFunc(s, IsWhiteSpace) }

// TrimStart is String.prototype.trimStart.
func TrimStart(s string) string { return strings.TrimLeftFunc(s, IsWhiteSpace) }

// TrimEnd is String.prototype.trimEnd.
func TrimEnd(s string) string { return strings.TrimRightFunc(s, IsWhiteSpace) }

// Repeat is String.prototype.repeat; a negative count panics with a RangeError.
func Repeat(s string, count int) string {
	if count < 0 {
		throwRange("Invalid count value: " + strconv.Itoa(count))
	}
	if count == 0 || s == "" {
		return ""
	}
	if hasSurrogateBytes(s) {
		u := CodeUnits(s)
		out := make([]uint16, 0, len(u)*count)
		for range count {
			out = append(out, u...)
		}
		return FromCodeUnits(out)
	}
	return strings.Repeat(s, count)
}

// pad builds the filler of String.prototype.padStart/padEnd.
func pad(s string, maxLength int, fill string) (string, bool) {
	n := Len(s)
	if maxLength <= n || fill == "" {
		return "", false
	}
	need := maxLength - n
	w := CodeUnits(fill)
	out := make([]uint16, need)
	for i := range out {
		out[i] = w[i%len(w)]
	}
	return FromCodeUnits(out), true
}

// PadStart is String.prototype.padStart(maxLength, fill); pass " " for the default filler.
func PadStart(s string, maxLength int, fill string) string {
	p, ok := pad(s, maxLength, fill)
	if !ok {
		return s
	}
	return concat(p, s)
}

// PadEnd is String.prototype.padEnd(maxLength, fill); pass " " for the default filler.
func PadEnd(s string, maxLength int, fill string) string {
	p, ok := pad(s, maxLength, fill)
	if !ok {
		return s
	}
	return concat(s, p)
}

// Concat joins strings as JavaScript's + does, merging a lone high surrogate at the end of one
// part with a lone low surrogate at the start of the next into one pair.
func Concat(parts ...string) string {
	var b []byte
	for _, p := range parts {
		b = appendJoin(b, p)
	}
	return string(b)
}

func concat(a, b string) string {
	return string(appendJoin([]byte(a), b))
}

// appendJoin appends s to b, joining a trailing lone high surrogate with a leading low one.
func appendJoin(b []byte, s string) []byte {
	if n := len(b); n >= 3 && len(s) >= 3 && b[n-3] == 0xed && b[n-2] >= 0xa0 && b[n-2] <= 0xaf &&
		s[0] == 0xed && s[1] >= 0xb0 && s[1] <= 0xbf {
		hi, _ := decodeUnit(string(b[n-3:]), 0)
		lo, _ := decodeUnit(s, 0)
		b = utf8.AppendRune(b[:n-3], (hi-0xd800)<<10+(lo-0xdc00)+0x10000)
		return append(b, s[3:]...)
	}
	return append(b, s...)
}

// FromCharCode is String.fromCharCode: each argument is truncated to 16 bits.
func FromCharCode(codes ...int) string {
	u := make([]uint16, len(codes))
	for i, c := range codes {
		u[i] = uint16(c)
	}
	return FromCodeUnits(u)
}

// FromCodePoint is String.fromCodePoint; a value outside 0..0x10FFFF panics with a RangeError.
func FromCodePoint(cps ...int) string {
	var b []byte
	var u []uint16
	for _, c := range cps {
		if c < 0 || c > 0x10ffff {
			throwRange("Invalid code point " + strconv.Itoa(c))
		}
		if c >= 0x10000 {
			c -= 0x10000
			u = append(u, uint16(0xd800+(c>>10)), uint16(0xdc00+(c&0x3ff)))
		} else {
			u = append(u, uint16(c))
		}
	}
	return string(appendWTF8(b, u))
}

// Compare compares strings as JavaScript's relational operators do (UTF-16 code unit order). It
// returns -1, 0 or 1.
func Compare(a, b string) int {
	i := 0
	for i < len(a) && i < len(b) && a[i] == b[i] {
		i++
	}
	if i == len(a) && i == len(b) {
		return 0
	}
	if (i == len(a) || a[i] < 0x80) && (i == len(b) || b[i] < 0x80) {
		switch {
		case i == len(a):
			return -1
		case i == len(b):
			return 1
		case a[i] < b[i]:
			return -1
		}
		return 1
	}
	ua, ub := CodeUnits(a), CodeUnits(b)
	for j := 0; j < len(ua) && j < len(ub); j++ {
		if ua[j] != ub[j] {
			if ua[j] < ub[j] {
				return -1
			}
			return 1
		}
	}
	switch {
	case len(ua) < len(ub):
		return -1
	case len(ua) > len(ub):
		return 1
	}
	return 0
}

// IsWellFormed is String.prototype.isWellFormed.
func IsWellFormed(s string) bool {
	return utf8.ValidString(s) || !hasSurrogateBytes(FromCodeUnits(CodeUnits(s)))
}

// ToWellFormed is String.prototype.toWellFormed: lone surrogates become U+FFFD.
func ToWellFormed(s string) string {
	if utf8.ValidString(s) {
		return s
	}
	u := CodeUnits(s)
	for i := 0; i < len(u); i++ {
		c := u[i]
		if c >= 0xd800 && c < 0xdc00 && i+1 < len(u) && u[i+1] >= 0xdc00 && u[i+1] < 0xe000 {
			i++
			continue
		}
		if c >= 0xd800 && c < 0xe000 {
			u[i] = 0xfffd
		}
	}
	return FromCodeUnits(u)
}
