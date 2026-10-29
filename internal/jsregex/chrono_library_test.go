package jsregex

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
)

// TestChronoLibrary checks every regex chrono builds (composed parser patterns of every locale
// configuration, every new RegExp executed while loading and parsing, every regex literal and
// String.prototype.match string argument), recorded by dev/regex/gen_chrono_library.mjs: exec
// on every subject (index, every group, named groups), matchAll for global regexes, and replace
// with a few templates on matching subjects. CHRONO_REGEX_VECTORS overrides the file.
func TestChronoLibrary(t *testing.T) {
	path := "testdata/chrono_library.jsonl.gz"
	if p := os.Getenv("CHRONO_REGEX_VECTORS"); p != "" {
		path = p
	}
	lines := gzLines(t, path)
	var h struct {
		Templates []string
		Subjects  []jsStr
	}
	if err := json.Unmarshal(lines[0], &h); err != nil {
		t.Fatalf("header: %v", err)
	}
	subjects := make([]string, len(h.Subjects))
	for i, s := range h.Subjects {
		subjects[i] = s.s
	}
	type indexed struct {
		I int
		V json.RawMessage
	}
	decodePairs := func(raw []json.RawMessage) ([]indexed, error) {
		out := make([]indexed, len(raw))
		for k, r := range raw {
			var pair []json.RawMessage
			if err := json.Unmarshal(r, &pair); err != nil || len(pair) != 2 {
				return nil, fmt.Errorf("bad pair %s", r)
			}
			if err := json.Unmarshal(pair[0], &out[k].I); err != nil {
				return nil, err
			}
			out[k].V = pair[1]
		}
		return out, nil
	}
	var fails []string
	bad := func(p, f, subj, what string, got, want any) {
		fails = append(fails, fmt.Sprintf("/%s/%s on %q: %s\n   got:  %v\n   want: %v", p, f, subj, what, got, want))
	}
	execs, matched := 0, 0
	for _, line := range lines[1:] {
		var v struct {
			P   jsStr
			F   string
			O   []string
			Src jsStr
			Fl  string
			Str jsStr
			Ng  int
			M   []json.RawMessage
			A   []json.RawMessage
			R   []json.RawMessage
		}
		if err := json.Unmarshal(line, &v); err != nil {
			t.Fatalf("bad vector: %v", err)
		}
		p, f := v.P.s, v.F
		re, err := Compile(p, f)
		if err != nil {
			bad(p, f, "", "compile", err, "ok")
			continue
		}
		if re.Source() != v.Src.s || re.Flags() != v.Fl || re.String() != v.Str.s || re.NumGroups() != v.Ng {
			bad(p, f, "", "source/flags/string/ngroups", []any{re.Source(), re.Flags(), re.String(), re.NumGroups()}, []any{v.Src.s, v.Fl, v.Str.s, v.Ng})
		}
		m, err := decodePairs(v.M)
		if err != nil {
			t.Fatal(err)
		}
		want := make(map[int]json.RawMessage, len(m))
		for _, x := range m {
			want[x.I] = x.V
		}
		for i, s := range subjects {
			execs++
			got := re.Exec(s)
			w, ok := want[i]
			if !ok {
				if got != nil {
					bad(p, f, s, "exec", canon(execRec(re, got)), "null")
				}
				continue
			}
			matched++
			if g, wr := canon(execRec(re, got)), canonRaw(w); g != wr {
				bad(p, f, s, "exec", g, wr)
				continue
			}
			// Group values must be the slices the indices name.
			u := toUTF16(nil, s)
			for k := range got.Groups {
				if got.Defined[k] && got.Groups[k] != fromUTF16(u[got.Starts[k]:got.Ends[k]]) {
					bad(p, f, s, fmt.Sprintf("group %d value", k), got.Groups[k], "slice of indices")
				}
			}
			if !re.Test(s) {
				bad(p, f, s, "test", false, true)
			}
		}
		all, err := decodePairs(v.A)
		if err != nil {
			t.Fatal(err)
		}
		for _, x := range all {
			ms := re.MatchAll(subjects[x.I])
			recs := make([]any, len(ms))
			for k, mm := range ms {
				recs[k] = execRec(re, mm)
			}
			if g, wr := canon(recs), canonRaw(x.V); g != wr {
				bad(p, f, subjects[x.I], "matchAll", g, wr)
			}
		}
		reps, err := decodePairs(v.R)
		if err != nil {
			t.Fatal(err)
		}
		for _, x := range reps {
			var wr []jsStr
			if err := json.Unmarshal(x.V, &wr); err != nil {
				t.Fatal(err)
			}
			for k, tm := range h.Templates {
				if g := re.Replace(subjects[x.I], tm); g != wr[k].s {
					bad(p, f, subjects[x.I], "replace "+tm, canon(enc(g)), canon(enc(wr[k].s)))
				}
			}
		}
	}
	t.Logf("%d regexes, %d subjects, %d exec calls, %d matches, %d mismatches", len(lines)-1, len(subjects), execs, matched, len(fails))
	if len(fails) > 0 {
		if len(fails) > 40 {
			fails = append(fails[:40], fmt.Sprintf("... and %d more", len(fails)-40))
		}
		t.Errorf("mismatches:\n%s", strings.Join(fails, "\n"))
	}
}
