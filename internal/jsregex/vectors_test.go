package jsregex

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

// jsStr is a JS string as recorded by dev/regex/lib.mjs: a JSON string, or
// {"u":[code units]} when it contains lone surrogates.
type jsStr struct {
	s     string
	valid bool // false for JSON null (undefined)
}

func (j *jsStr) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		*j = jsStr{}
		return nil
	}
	if len(b) > 0 && b[0] == '{' {
		var o struct{ U []uint16 }
		if err := json.Unmarshal(b, &o); err != nil {
			return err
		}
		*j = jsStr{fromUTF16(o.U), true}
		return nil
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	*j = jsStr{s, true}
	return nil
}

// enc mirrors dev/regex/lib.mjs enc for comparisons.
func enc(s string) any {
	u := toUTF16(nil, s)
	for i := 0; i < len(u); i++ {
		if isLead(u[i]) && i+1 < len(u) && isTrail(u[i+1]) {
			i++
			continue
		}
		if isLead(u[i]) || isTrail(u[i]) {
			return map[string]any{"u": u}
		}
	}
	return s
}

// vecCase is one subject. exec, test, search and match are derived from the
// matchAll records A (for a fresh regexp they are determined by A[0]).
type vecCase struct {
	S  jsStr             `json:"s"`
	Li *int              `json:"li"`
	Xl json.RawMessage   `json:"xl"`
	A  []json.RawMessage `json:"a"`
	R  []jsStr           `json:"r"`
	Sp []json.RawMessage
}

type vector struct {
	P   jsStr     `json:"p"`
	F   string    `json:"f"`
	Err jsStr     `json:"err"`
	Src jsStr     `json:"src"`
	Fl  string    `json:"fl"`
	Str jsStr     `json:"str"`
	Ng  int       `json:"ng"`
	C   []vecCase `json:"c"`
}

func (c *vecCase) UnmarshalJSON(b []byte) error {
	type plain vecCase
	var aux struct {
		plain
		Sp []json.RawMessage `json:"sp"`
	}
	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}
	*c = vecCase(aux.plain)
	c.Sp = aux.Sp
	return nil
}

// canon re-encodes a JSON value so that equal values compare equal as strings.
func canon(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return "marshal error: " + err.Error()
	}
	var x any
	if err := json.Unmarshal(b, &x); err != nil {
		return "unmarshal error: " + err.Error()
	}
	b, _ = json.Marshal(x)
	return string(b)
}

func canonRaw(r json.RawMessage) string {
	if len(r) == 0 {
		return "null"
	}
	var x any
	if err := json.Unmarshal(r, &x); err != nil {
		return "unmarshal error: " + err.Error()
	}
	b, _ := json.Marshal(x)
	return string(b)
}

// execRec mirrors dev/regex/lib.mjs execRec.
func execRec(re *Regexp, m *Match) any {
	if m == nil {
		return nil
	}
	idx := make([]any, len(m.Starts))
	for i := range m.Starts {
		if m.Defined[i] {
			idx[i] = []int{m.Starts[i], m.Ends[i]}
		}
	}
	r := []any{m.Index, idx}
	if names := re.GroupNames(); len(names) > 0 {
		g := make([]any, len(names))
		for i, n := range names {
			if v, ok := m.Named(n); ok {
				g[i] = []any{n, enc(v)}
			} else {
				g[i] = []any{n, nil}
			}
		}
		r = append(r, g)
	}
	return r
}

type mismatch struct {
	what       string
	got, want  string
	pat, flags string
	subject    string
}

func (m mismatch) String() string {
	return fmt.Sprintf("/%s/%s on %q: %s\n   got:  %s\n   want: %s", m.pat, m.flags, m.subject, m.what, m.got, m.want)
}

// gzLines returns the non-empty lines of a gzipped JSONL file.
func gzLines(t *testing.T, path string) [][]byte {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	zr, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	sc := bufio.NewScanner(zr)
	sc.Buffer(make([]byte, 1<<20), 64<<20)
	var out [][]byte
	for sc.Scan() {
		if line := bytes.TrimSpace(sc.Bytes()); len(line) > 0 {
			out = append(out, append([]byte(nil), line...))
		}
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}
	if len(out) == 0 {
		t.Fatalf("%s: empty", path)
	}
	return out
}

func loadVectors(t *testing.T, path string) ([]string, []vector) {
	t.Helper()
	lines := gzLines(t, path)
	var h struct{ Templates []string }
	if err := json.Unmarshal(lines[0], &h); err != nil {
		t.Fatalf("header: %v", err)
	}
	vecs := make([]vector, len(lines)-1)
	for i, line := range lines[1:] {
		if err := json.Unmarshal(line, &vecs[i]); err != nil {
			t.Fatalf("bad vector %s: %v", line, err)
		}
	}
	return h.Templates, vecs
}

// checkVector compares one recorded vector with this package's behaviour.
func checkVector(templates []string, v *vector) (out []mismatch) {
	pat := v.P.s
	bad := func(what, subj string, got, want any) {
		out = append(out, mismatch{what, fmt.Sprint(got), fmt.Sprint(want), pat, v.F, subj})
	}
	defer func() {
		if r := recover(); r != nil {
			bad("panic", "", r, "no panic")
		}
	}()
	re, err := Compile(pat, v.F)
	if v.Err.valid {
		if err == nil {
			bad("compile", "", "ok", v.Err.s)
		} else if err.Error() != v.Err.s {
			bad("compile", "", err.Error(), v.Err.s)
		}
		return
	}
	if err != nil {
		bad("compile", "", err.Error(), "ok")
		return
	}
	if re.Source() != v.Src.s {
		bad("source", "", re.Source(), v.Src.s)
	}
	if re.Flags() != v.Fl {
		bad("flags", "", re.Flags(), v.Fl)
	}
	if re.String() != v.Str.s {
		bad("string", "", re.String(), v.Str.s)
	}
	if re.NumGroups() != v.Ng {
		bad("ngroups", "", re.NumGroups(), v.Ng)
	}
	for _, c := range v.C {
		s := c.S.s
		u := toUTF16(nil, s)
		var first []json.RawMessage // [index, indices, groups?] of the first match
		wantExec := json.RawMessage("null")
		if len(c.A) > 0 {
			wantExec = c.A[0]
			_ = json.Unmarshal(c.A[0], &first)
		}
		if g, w := canon(execRec(re, re.Exec(s))), canonRaw(wantExec); g != w {
			bad("exec", s, g, w)
		}
		if c.Li != nil {
			if g, w := canon(execRec(re, re.ExecAt(s, *c.Li))), canonRaw(c.Xl); g != w {
				bad(fmt.Sprintf("exec@%d", *c.Li), s, g, w)
			}
		}
		all := re.MatchAll(s)
		recs := make([]any, len(all))
		for i, m := range all {
			recs[i] = execRec(re, m)
		}
		if g, w := canon(recs), canon(c.A); g != w {
			bad("matchAll", s, g, w)
		}
		if g := re.Test(s); g != (len(c.A) > 0) {
			bad("test", s, g, len(c.A) > 0)
		}
		wantSearch := -1
		if first != nil {
			_ = json.Unmarshal(first[0], &wantSearch)
		}
		if g := re.Search(s); g != wantSearch {
			bad("search", s, g, wantSearch)
		}
		for i, tm := range templates {
			want := s
			if c.R != nil {
				want = c.R[i].s
			}
			if g := re.Replace(s, tm); g != want {
				bad("replace "+tm, s, canon(enc(g)), canon(enc(want)))
			}
		}
		// match: whole matches when global, else the groups of the first match.
		slices := func(rec json.RawMessage, onlyWhole bool) []any {
			var r []json.RawMessage
			_ = json.Unmarshal(rec, &r)
			var idx [][]int
			_ = json.Unmarshal(r[1], &idx)
			if onlyWhole {
				idx = idx[:1]
			}
			out := make([]any, len(idx))
			for i, p := range idx {
				if p != nil {
					out[i] = enc(fromUTF16(u[p[0]:p[1]]))
				}
			}
			return out
		}
		var wm any
		if re.Global() && len(c.A) > 0 {
			var l []any
			for _, rec := range c.A {
				l = append(l, slices(rec, true)[0])
			}
			wm = l
		} else if len(c.A) > 0 {
			wm = slices(c.A[0], false)
		}
		var gm any
		if re.Global() {
			if ms := re.Match(s); ms != nil {
				l := make([]any, len(ms))
				for i, x := range ms {
					l[i] = enc(x)
				}
				gm = l
			}
		} else if m := re.Exec(s); m != nil {
			l := make([]any, len(m.Groups))
			for i, x := range m.Groups {
				if m.Defined[i] {
					l[i] = enc(x)
				}
			}
			gm = l
		}
		if g, w := canon(gm), canon(wm); g != w {
			bad("match", s, g, w)
		}
		if len(c.Sp) == 2 {
			var lim int
			_ = json.Unmarshal(c.Sp[0], &lim)
			parts, defined := re.Split(s, lim)
			gp := make([]any, len(parts))
			for i, x := range parts {
				if defined[i] {
					gp[i] = enc(x)
				}
			}
			if g, w := canon(gp), canonRaw(c.Sp[1]); g != w {
				bad(fmt.Sprintf("split(%d)", lim), s, g, w)
			}
		}
	}
	return out
}

func runVectorFile(t *testing.T, path string) {
	templates, vecs := loadVectors(t, path)
	if len(vecs) == 0 {
		t.Fatal("no vectors")
	}
	var fails []mismatch
	failedVecs := 0
	var slow []string
	for i := range vecs {
		t0 := time.Now()
		ms := checkVector(templates, &vecs[i])
		if d := time.Since(t0); d > 20*time.Second {
			slow = append(slow, fmt.Sprintf("/%s/%s took %v", vecs[i].P.s, vecs[i].F, d))
		}
		if len(ms) > 0 {
			failedVecs++
			fails = append(fails, ms...)
		}
	}
	for _, s := range slow {
		t.Logf("slow: %s", s)
	}
	t.Logf("%s: %d vectors, %d mismatching vectors, %d mismatches", path, len(vecs), failedVecs, len(fails))
	if len(fails) > 0 {
		var b strings.Builder
		for i, f := range fails {
			if i == 60 {
				fmt.Fprintf(&b, "... and %d more\n", len(fails)-60)
				break
			}
			b.WriteString(f.String())
			b.WriteByte('\n')
		}
		t.Errorf("%d mismatches:\n%s", len(fails), b.String())
	}
}

// TestVectors checks the differential fuzz vectors recorded from node by
// dev/regex/gen_vectors.mjs (JSREGEX_VECTORS overrides the file).
func TestVectors(t *testing.T) {
	path := "testdata/vectors.jsonl.gz"
	if p := os.Getenv("JSREGEX_VECTORS"); p != "" {
		path = p
	}
	runVectorFile(t, path)
}
