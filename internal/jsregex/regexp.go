// Package jsregex implements ECMAScript regular expressions with the exact
// syntax, semantics and error messages of V8's RegExp (Node.js 24, ICU 78,
// Unicode 17). The chrono port compiles every parser pattern and every regex
// of the original with it; the regexes chrono builds are checked against node
// in chrono_library_test.go (vectors from dev/regex/gen_chrono_library.mjs).
//
// Strings are Go strings in WTF-8 (lone surrogates encoded as 3-byte
// sequences); every index in the API is a UTF-16 code-unit index.
// Non-unicode patterns match UTF-16 code units, u/v patterns code points.
// A *Regexp is immutable and safe for concurrent use; it never carries a
// lastIndex between calls.
package jsregex

import (
	"fmt"
	"strings"
	"sync"
)

// SyntaxError is the error returned for invalid patterns or flags. Message
// is exactly the message of the SyntaxError V8 throws.
type SyntaxError struct{ Message string }

// Error returns the message.
func (e *SyntaxError) Error() string { return e.Message }

// Regexp is a compiled ECMAScript regular expression.
type Regexp struct {
	source                               string
	flags                                string
	hasIndices, global, icase, multiline bool
	dotAll, unicode, vmode, sticky       bool
	prog                                 mnode
	ncap                                 int
	nloops                               int
	anchorStart, anchorEnd               bool     // every alternative starts with ^ / ends with $ (without m)
	maxLen                               int      // maximum match length in code units (infinity if unbounded)
	capNames                             []string // by group index; "" when unnamed
	groupNames                           []string // distinct names in order of first appearance
	pool                                 sync.Pool
}

// Match is the result of a successful match.
type Match struct {
	// Index is the UTF-16 index of the match start.
	Index int
	// Groups[0] is the whole match, Groups[i] capture i ("" when undefined).
	Groups []string
	// Defined reports whether each group participated in the match.
	Defined []bool
	// Starts and Ends are UTF-16 indices of each group (-1 when undefined).
	Starts, Ends []int
	// Names maps each group name to its group index: the participating group
	// for duplicate names, else the first group with that name.
	Names map[string]int
	// Input is the subject string.
	Input string
}

// Named returns the value of the named group and whether it is defined.
func (m *Match) Named(name string) (string, bool) {
	i, ok := m.Names[name]
	if !ok {
		return "", false
	}
	return m.Groups[i], m.Defined[i]
}

func flagsError(flags string) error {
	return &SyntaxError{Message: "Invalid flags supplied to RegExp constructor '" + flags + "'"}
}

// Compile is new RegExp(source, flags).
func Compile(source, flags string) (re *Regexp, err error) {
	re = &Regexp{source: source}
	seen := map[rune]bool{}
	for _, f := range flags {
		if seen[f] {
			return nil, flagsError(flags)
		}
		seen[f] = true
		switch f {
		case 'd':
			re.hasIndices = true
		case 'g':
			re.global = true
		case 'i':
			re.icase = true
		case 'm':
			re.multiline = true
		case 's':
			re.dotAll = true
		case 'u':
			re.unicode = true
		case 'v':
			re.vmode = true
		case 'y':
			re.sticky = true
		default:
			return nil, flagsError(flags)
		}
	}
	if re.unicode && re.vmode {
		return nil, flagsError(flags)
	}
	var b strings.Builder
	for _, f := range "dgimsuvy" {
		if seen[f] {
			b.WriteRune(f)
		}
	}
	re.flags = b.String()

	canonFlags := re.flags
	defer func() {
		if r := recover(); r != nil {
			re = nil
			if pe, ok := r.(parseError); ok {
				err = &SyntaxError{Message: "Invalid regular expression: /" + source + "/" + canonFlags + ": " + pe.msg}
				return
			}
			err = fmt.Errorf("jsregex: internal error compiling /%s/%s: %v", source, flags, r)
		}
	}()
	p := &parser{in: toUTF16(nil, source), unicode: re.unicode, vmode: re.vmode,
		flags: modFlags{icase: re.icase, multiline: re.multiline, dotAll: re.dotAll}}
	ast := p.parse()
	re.ncap = p.captures
	re.capNames = make([]string, re.ncap+1)
	seenName := map[string]bool{}
	for _, g := range p.named {
		re.capNames[g.index] = fromUTF16(g.name)
	}
	for i := 1; i <= re.ncap; i++ {
		if n := re.capNames[i]; n != "" && !seenName[n] {
			seenName[n] = true
			re.groupNames = append(re.groupNames, n)
		}
	}
	c := &compiler{unicode: re.unicode || re.vmode, vmode: re.vmode, ncap: re.ncap}
	re.prog = c.comp(ast, acceptNode{}, false)
	re.anchorStart, re.anchorEnd = anchoredAt(ast, false), anchoredAt(ast, true)
	re.maxLen = ast.maxMatch(re.unicode || re.vmode)
	re.nloops = c.nloops
	return re, nil
}

// MustCompile is like Compile but panics on error.
func MustCompile(source, flags string) *Regexp {
	re, err := Compile(source, flags)
	if err != nil {
		panic(err)
	}
	return re
}

// Source returns re.source as V8 prints it.
func (re *Regexp) Source() string { return escapeSource(re.source) }

// Flags returns the canonical flags string ("dgimsuvy" order).
func (re *Regexp) Flags() string { return re.flags }

// String returns "/source/flags".
func (re *Regexp) String() string { return "/" + re.Source() + "/" + re.flags }

// Global reports the g flag.
func (re *Regexp) Global() bool { return re.global }

// IgnoreCase reports the i flag.
func (re *Regexp) IgnoreCase() bool { return re.icase }

// Multiline reports the m flag.
func (re *Regexp) Multiline() bool { return re.multiline }

// Sticky reports the y flag.
func (re *Regexp) Sticky() bool { return re.sticky }

// Unicode reports the u flag.
func (re *Regexp) Unicode() bool { return re.unicode }

// UnicodeSets reports the v flag.
func (re *Regexp) UnicodeSets() bool { return re.vmode }

// DotAll reports the s flag.
func (re *Regexp) DotAll() bool { return re.dotAll }

// HasIndices reports the d flag.
func (re *Regexp) HasIndices() bool { return re.hasIndices }

// NumGroups returns the number of capturing groups.
func (re *Regexp) NumGroups() int { return re.ncap }

// GroupNames returns the distinct group names in order of first appearance.
func (re *Regexp) GroupNames() []string { return append([]string(nil), re.groupNames...) }

func (re *Regexp) fullUnicode() bool { return re.unicode || re.vmode }

func (re *Regexp) getMachine(in []uint16) *machine {
	m, _ := re.pool.Get().(*machine)
	if m == nil {
		m = &machine{caps: make([]int, 2*(re.ncap+1)), open: make([]int, re.ncap+1), loops: make([]loopState, re.nloops)}
	}
	m.in = in
	m.stack = m.stack[:0]
	return m
}

func (re *Regexp) putMachine(m *machine) {
	m.in = nil
	if cap(m.ubuf) > 1<<12 {
		m.ubuf = nil
	}
	re.pool.Put(m)
}

// run searches for a match starting at lastIndex (only at lastIndex when
// sticky). stepBack is V8's OptionallyStepBackToLeadSurrogate, compiled
// into global and sticky u/v regexps: when lastIndex is inside a surrogate
// pair, matching is first attempted from the lead surrogate and then from
// lastIndex itself.
func (re *Regexp) run(m *machine, lastIndex int, sticky, stepBack bool) bool {
	n := len(m.in)
	if lastIndex > n {
		return false
	}
	start := lastIndex
	if !sticky {
		// No match can start elsewhere (V8 makes the same deductions).
		if re.anchorStart {
			sticky = true
		}
		if re.anchorEnd && re.maxLen < infinity && n-re.maxLen > start {
			start = n - re.maxLen
		}
	}
	if stepBack && re.fullUnicode() && m.midPair(start) {
		if !sticky {
			start--
		} else if re.attempt(m, start-1) {
			return true
		}
	}
	for {
		if re.attempt(m, start) {
			return true
		}
		if sticky || start >= n {
			return false
		}
		// V8 retries at every code unit, even inside surrogate pairs.
		start++
	}
}

// attempt tries to match at exactly start.
func (re *Regexp) attempt(m *machine, start int) bool {
	for k := range m.caps {
		m.caps[k] = -1
	}
	if re.prog.match(m, start) {
		m.caps[0], m.caps[1] = start, m.matchEnd
		return true
	}
	return false
}

type subject struct {
	s     string
	u     []uint16
	ascii bool
}

func newSubject(s string) *subject {
	sub := &subject{s: s, ascii: isASCII(s)}
	if sub.ascii {
		sub.u = make([]uint16, len(s))
		for i := 0; i < len(s); i++ {
			sub.u[i] = uint16(s[i])
		}
	} else {
		sub.u = toUTF16(nil, s)
	}
	return sub
}

func (sub *subject) slice(a, b int) string {
	if sub.ascii {
		return sub.s[a:b]
	}
	return fromUTF16(sub.u[a:b])
}

func (re *Regexp) buildMatch(sub *subject, m *machine) *Match {
	n := re.ncap + 1
	res := &Match{Index: m.caps[0], Groups: make([]string, n), Defined: make([]bool, n),
		Starts: make([]int, n), Ends: make([]int, n), Input: sub.s}
	for i := 0; i < n; i++ {
		s, e := m.caps[2*i], m.caps[2*i+1]
		if s >= 0 && e >= 0 {
			res.Groups[i] = sub.slice(s, e)
			res.Defined[i] = true
			res.Starts[i], res.Ends[i] = s, e
		} else {
			res.Starts[i], res.Ends[i] = -1, -1
		}
	}
	if len(re.groupNames) > 0 {
		res.Names = make(map[string]int, len(re.groupNames))
		for i := 1; i < n; i++ {
			name := re.capNames[i]
			if name == "" {
				continue
			}
			if j, ok := res.Names[name]; !ok || (!res.Defined[j] && res.Defined[i]) {
				res.Names[name] = i
			}
		}
	}
	return res
}

func (re *Regexp) execSub(sub *subject, lastIndex int, sticky bool) *Match {
	m := re.getMachine(sub.u)
	defer re.putMachine(m)
	if !re.run(m, lastIndex, sticky, re.global || re.sticky) {
		return nil
	}
	return re.buildMatch(sub, m)
}

// ExecAt runs RegExpBuiltinExec with the given lastIndex. As in JS,
// lastIndex is only used when the regexp is global or sticky (otherwise
// matching starts at 0); sticky regexps only match at lastIndex.
func (re *Regexp) ExecAt(s string, lastIndex int) *Match {
	if !re.global && !re.sticky {
		lastIndex = 0
	}
	if lastIndex < 0 {
		lastIndex = 0
	}
	return re.execSub(newSubject(s), lastIndex, re.sticky)
}

// Exec is re.exec(s) with lastIndex 0.
func (re *Regexp) Exec(s string) *Match { return re.ExecAt(s, 0) }

// Test is re.test(s) with lastIndex 0.
func (re *Regexp) Test(s string) bool {
	var u []uint16
	m := re.getMachine(nil)
	if u == nil {
		m.ubuf = toUTF16(m.ubuf[:0], s)
		u = m.ubuf
	}
	m.in = u
	ok := re.run(m, 0, re.sticky, false)
	re.putMachine(m)
	return ok
}

// advanceIndex is AdvanceStringIndex.
func advanceIndex(u []uint16, i int, unicode bool) int {
	if !unicode || i+1 >= len(u) {
		return i + 1
	}
	if isLead(u[i]) && isTrail(u[i+1]) {
		return i + 2
	}
	return i + 1
}

func (re *Regexp) allMatches(sub *subject) []*Match {
	var out []*Match
	m := re.getMachine(sub.u)
	defer re.putMachine(m)
	last := 0
	for {
		if !re.run(m, last, re.sticky, true) {
			return out
		}
		mt := re.buildMatch(sub, m)
		out = append(out, mt)
		end := m.caps[1]
		if end == m.caps[0] {
			last = advanceIndex(sub.u, end, re.fullUnicode())
		} else {
			last = end
		}
	}
}

// MatchAll is [...s.matchAll(re)] for a global copy of re.
func (re *Regexp) MatchAll(s string) []*Match { return re.allMatches(newSubject(s)) }

// Match is s.match(re): the groups of the first match when re is not
// global, all whole matches when it is; nil when there is no match.
func (re *Regexp) Match(s string) []string {
	sub := newSubject(s)
	if !re.global {
		mt := re.execSub(sub, 0, re.sticky)
		if mt == nil {
			return nil
		}
		return mt.Groups
	}
	all := re.allMatches(sub)
	if len(all) == 0 {
		return nil
	}
	out := make([]string, len(all))
	for i, mt := range all {
		out[i] = mt.Groups[0]
	}
	return out
}

// Search is s.search(re).
func (re *Regexp) Search(s string) int {
	mt := re.execSub(newSubject(s), 0, re.sticky)
	if mt == nil {
		return -1
	}
	return mt.Index
}

func (re *Regexp) replaceMatches(sub *subject) []*Match {
	if re.global {
		return re.allMatches(sub)
	}
	if mt := re.execSub(sub, 0, re.sticky); mt != nil {
		return []*Match{mt}
	}
	return nil
}

// Replace is s.replace(re, replacement) with GetSubstitution patterns
// ($$, $&, $`, $', $n, $nn, $<name>); all matches are replaced when re is global.
func (re *Regexp) Replace(s, replacement string) string {
	sub := newSubject(s)
	ms := re.replaceMatches(sub)
	if len(ms) == 0 {
		return s
	}
	tmpl := toUTF16(nil, replacement)
	var out []uint16
	next := 0
	for _, mt := range ms {
		pos := mt.Starts[0]
		if pos < next {
			continue
		}
		out = append(out, sub.u[next:pos]...)
		out = re.substitute(out, sub, mt, tmpl)
		next = mt.Ends[0]
	}
	out = append(out, sub.u[next:]...)
	return fromUTF16(out)
}

// ReplaceFunc is s.replace(re, fn).
func (re *Regexp) ReplaceFunc(s string, fn func(m *Match) string) string {
	sub := newSubject(s)
	ms := re.replaceMatches(sub)
	if len(ms) == 0 {
		return s
	}
	var out []uint16
	next := 0
	for _, mt := range ms {
		pos := mt.Starts[0]
		if pos < next {
			continue
		}
		out = append(out, sub.u[next:pos]...)
		out = toUTF16(out, fn(mt))
		next = mt.Ends[0]
	}
	out = append(out, sub.u[next:]...)
	return fromUTF16(out)
}

func (re *Regexp) substitute(out []uint16, sub *subject, mt *Match, t []uint16) []uint16 {
	m := re.ncap
	pos, end := mt.Starts[0], mt.Ends[0]
	for i := 0; i < len(t); i++ {
		c := t[i]
		if c != '$' || i+1 >= len(t) {
			out = append(out, c)
			continue
		}
		d := t[i+1]
		switch {
		case d == '$':
			out = append(out, '$')
			i++
		case d == '&':
			out = append(out, sub.u[pos:end]...)
			i++
		case d == '`':
			out = append(out, sub.u[:pos]...)
			i++
		case d == '\'':
			out = append(out, sub.u[end:]...)
			i++
		case d >= '0' && d <= '9':
			digits := 1
			idx := int(d - '0')
			if i+2 < len(t) && t[i+2] >= '0' && t[i+2] <= '9' {
				digits = 2
				idx = idx*10 + int(t[i+2]-'0')
			}
			if idx > m && digits == 2 {
				digits = 1
				idx = int(d - '0')
			}
			if idx >= 1 && idx <= m {
				if mt.Defined[idx] {
					out = append(out, sub.u[mt.Starts[idx]:mt.Ends[idx]]...)
				}
			} else {
				out = append(out, t[i:i+1+digits]...)
			}
			i += digits
		case d == '<':
			gt := -1
			for k := i + 2; k < len(t); k++ {
				if t[k] == '>' {
					gt = k
					break
				}
			}
			if gt < 0 || len(re.groupNames) == 0 {
				out = append(out, '$', '<')
				i++
				continue
			}
			name := fromUTF16(t[i+2 : gt])
			if v, ok := mt.Named(name); ok {
				out = toUTF16(out, v)
			}
			i = gt
		default:
			out = append(out, '$')
		}
	}
	return out
}

// Split is s.split(re, limit); limit < 0 means no limit. defined[i] is false
// for parts that come from undefined captures. Like V8, non-sticky regexps
// use the fast path (a plain search from each position, which may stop
// inside a surrogate pair) and sticky ones the specification's algorithm.
func (re *Regexp) Split(s string, limit int) (parts []string, defined []bool) {
	lim := uint32(0xFFFFFFFF)
	if limit >= 0 {
		lim = uint32(limit)
	}
	parts, defined = []string{}, []bool{}
	if lim == 0 {
		return
	}
	sub := newSubject(s)
	size := len(sub.u)
	m := re.getMachine(sub.u)
	defer re.putMachine(m)
	if size == 0 {
		if re.run(m, 0, re.sticky, true) {
			return
		}
		return []string{s}, []bool{true}
	}
	uni := re.fullUnicode()
	push := func(a, b int, ok bool) bool {
		if ok {
			parts = append(parts, sub.slice(a, b))
		} else {
			parts = append(parts, "")
		}
		defined = append(defined, ok)
		return uint32(len(parts)) == lim
	}
	pushCaptures := func() bool {
		for i := 1; i <= re.ncap; i++ {
			cs, ce := m.caps[2*i], m.caps[2*i+1]
			if push(cs, ce, cs >= 0 && ce >= 0) {
				return true
			}
		}
		return false
	}
	if !re.sticky {
		next, last := 0, 0
		for next != size {
			if !re.run(m, next, false, re.global) {
				break
			}
			from, to := m.caps[0], m.caps[1]
			if from == size {
				break
			}
			if to == next && to == last {
				next = advanceIndex(sub.u, next, uni)
				continue
			}
			if push(last, from, true) || pushCaptures() {
				return
			}
			last, next = to, to
		}
		push(last, size, true)
		return
	}
	// V8's spec path advances by code points only with the u flag (not v).
	uni = re.unicode
	p, q := 0, 0
	for q < size {
		if !re.run(m, q, true, true) {
			q = advanceIndex(sub.u, q, uni)
			continue
		}
		e := m.caps[1]
		if e > size {
			e = size
		}
		if e == p {
			q = advanceIndex(sub.u, q, uni)
			continue
		}
		if push(p, q, true) {
			return
		}
		p = e
		if pushCaptures() {
			return
		}
		q = p
	}
	push(p, size, true)
	return
}

// escapeSource implements V8's EscapeRegExpSource.
func escapeSource(src string) string {
	if src == "" {
		return "(?:)"
	}
	u := toUTF16(nil, src)
	var out []uint16
	inClass := false
	for i := 0; i < len(u); i++ {
		c := u[i]
		switch {
		case c == '\\':
			if i+1 < len(u) && isLineTerm(u[i+1]) {
				continue
			}
			out = append(out, c)
			if i+1 < len(u) {
				i++
				out = append(out, u[i])
			}
		case c == '/' && !inClass:
			out = append(out, '\\', '/')
		case c == '[':
			inClass = true
			out = append(out, c)
		case c == ']':
			inClass = false
			out = append(out, c)
		case c == '\n':
			out = append(out, '\\', 'n')
		case c == '\r':
			out = append(out, '\\', 'r')
		case c == 0x2028:
			out = append(out, toUTF16(nil, `\u2028`)...)
		case c == 0x2029:
			out = append(out, toUTF16(nil, `\u2029`)...)
		default:
			out = append(out, c)
		}
	}
	return fromUTF16(out)
}

// anchoredAt reports whether every alternative of t begins with ^ (or, with
// end, ends with $), outside lookarounds and without the m flag.
func anchoredAt(t *term, end bool) bool {
	switch t.typ {
	case tAssert:
		return (!end && t.assert == aStartInput) || (end && t.assert == aEndInput)
	case tCapture, tGroup:
		return anchoredAt(t.sub, end)
	case tSeq:
		if end {
			return anchoredAt(t.subs[len(t.subs)-1], end)
		}
		return anchoredAt(t.subs[0], end)
	case tAlt:
		for _, a := range t.subs {
			if !anchoredAt(a, end) {
				return false
			}
		}
		return true
	}
	return false
}
