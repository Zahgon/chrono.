// RegExp and RegExpMatch: the JavaScript RegExp objects used by the ported code (no direct TS
// counterpart; every pattern is run by internal/jsregex with its original source and flags).

package core

import (
	"sync"

	"github.com/wanasit/chrono-go/internal/js"
	"github.com/wanasit/chrono-go/internal/jsregex"
)

// RegExp is an ECMAScript regular expression. Like a JavaScript RegExp object it carries a
// lastIndex, used (and updated) by Exec and Test only when the g or y flag is set.
type RegExp struct {
	re *jsregex.Regexp
	// LastIndex is the regular expression's lastIndex property.
	LastIndex int
}

// NewRegExp is `new RegExp(source, flags)`; the error's text is V8's SyntaxError message.
func NewRegExp(source, flags string) (*RegExp, error) {
	re, err := jsregex.Compile(source, flags)
	if err != nil {
		return nil, err
	}
	return &RegExp{re: re}, nil
}

// MustRegExp is NewRegExp for patterns that are known to be valid (a module-level `new RegExp`
// or regex literal); an invalid pattern panics with a SyntaxError *Error, as JavaScript throws.
func MustRegExp(source, flags string) *RegExp {
	re, err := NewRegExp(source, flags)
	if err != nil {
		panic(&Error{Name: "SyntaxError", Message: err.Error()})
	}
	return re
}

var stringPatternCache sync.Map // source -> *jsregex.Regexp

// StringPattern returns `new RegExp(source)` (no flags) for String.prototype.match(string),
// compiling each distinct source once.
func StringPattern(source string) *RegExp {
	if re, ok := stringPatternCache.Load(source); ok {
		return &RegExp{re: re.(*jsregex.Regexp)}
	}
	r := MustRegExp(source, "")
	stringPatternCache.Store(source, r.re)
	return r
}

// Source is RegExp.prototype.source.
func (r *RegExp) Source() string { return r.re.Source() }

// Flags is RegExp.prototype.flags.
func (r *RegExp) Flags() string { return r.re.Flags() }

// String is RegExp.prototype.toString.
func (r *RegExp) String() string { return r.re.String() }

// Exec is RegExp.prototype.exec: nil when there is no match. With the g or y flag it starts at
// LastIndex and updates it.
func (r *RegExp) Exec(s string) *RegExpMatch {
	stateful := r.re.Global() || r.re.Sticky()
	start := 0
	if stateful {
		start = r.LastIndex
		if start < 0 {
			start = 0
		}
		if start > js.Len(s) {
			r.LastIndex = 0
			return nil
		}
	}
	m := r.re.ExecAt(s, start)
	if m == nil {
		if stateful {
			r.LastIndex = 0
		}
		return nil
	}
	if stateful {
		r.LastIndex = m.Ends[0]
	}
	return newRegExpMatch(m)
}

// Test is RegExp.prototype.test.
func (r *RegExp) Test(s string) bool { return r.Exec(s) != nil }

// MatchString is String.prototype.match(regexp) for a non-global regular expression (the same as
// Exec); it returns nil when there is no match.
func (r *RegExp) MatchString(s string) *RegExpMatch { return r.Exec(s) }

// Replace is String.prototype.replace(regexp, replacement) (all matches with the g flag, else the
// first one; `$` patterns are substituted).
func (r *RegExp) Replace(s, replacement string) string {
	if r.re.Global() {
		r.LastIndex = 0
	}
	return r.re.Replace(s, replacement)
}

// Regexp returns the underlying compiled regular expression.
func (r *RegExp) Regexp() *jsregex.Regexp { return r.re }

// RegExpMatch is what RegExp.prototype.exec returns.
type RegExpMatch struct {
	// Index is the UTF-16 index of the match (Chrono adjusts it to the whole text before
	// Extract).
	Index int
	// Values holds match[0], match[1], ...; "" for an undefined group.
	Values []string
	// Defined is false where the group is undefined.
	Defined []bool
	// Input is the string the match was run on.
	Input string

	groups map[string]namedGroup
}

type namedGroup struct {
	value   string
	defined bool
}

func newRegExpMatch(m *jsregex.Match) *RegExpMatch {
	out := &RegExpMatch{
		Index:   m.Index,
		Values:  append([]string(nil), m.Groups...),
		Defined: append([]bool(nil), m.Defined...),
		Input:   m.Input,
	}
	if len(m.Names) > 0 {
		out.groups = make(map[string]namedGroup, len(m.Names))
		for name, i := range m.Names {
			out.groups[name] = namedGroup{m.Groups[i], m.Defined[i]}
		}
	}
	return out
}

// Len is match.length.
func (m *RegExpMatch) Len() int { return len(m.Values) }

// Get returns match[i] ("" when undefined or out of range).
func (m *RegExpMatch) Get(i int) string {
	if i < 0 || i >= len(m.Values) {
		return ""
	}
	return m.Values[i]
}

// Has reports whether match[i] is defined (`match[i] != null`).
func (m *RegExpMatch) Has(i int) bool {
	return i >= 0 && i < len(m.Defined) && m.Defined[i]
}

// Truthy reports whether match[i] is truthy (defined and not empty).
func (m *RegExpMatch) Truthy(i int) bool { return m.Has(i) && m.Values[i] != "" }

// Group returns match.groups[name] and whether it is defined.
func (m *RegExpMatch) Group(name string) (string, bool) {
	g, ok := m.groups[name]
	if !ok {
		return "", false
	}
	return g.value, g.defined
}

// HasGroups reports whether the regular expression had named groups (match.groups !== undefined).
func (m *RegExpMatch) HasGroups() bool { return m.groups != nil }

// Set is `match[i] = value` (defined).
func (m *RegExpMatch) Set(i int, value string, defined bool) {
	for len(m.Values) <= i {
		m.Values = append(m.Values, "")
		m.Defined = append(m.Defined, false)
	}
	m.Values[i] = value
	m.Defined[i] = defined
}
