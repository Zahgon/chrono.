package jsregex

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestAPI(t *testing.T) {
	re := MustCompile(`(?<y>\d{4})-(?<m>\d\d)|/`, "dgi")
	if re.Source() != `(?<y>\d{4})-(?<m>\d\d)|\/` || re.Flags() != "dgi" || re.String() != `/(?<y>\d{4})-(?<m>\d\d)|\//dgi` {
		t.Errorf("source/flags/string: %q %q %q", re.Source(), re.Flags(), re.String())
	}
	if !re.Global() || !re.IgnoreCase() || re.Multiline() || re.Sticky() || re.Unicode() || re.NumGroups() != 2 {
		t.Error("flag accessors")
	}
	if g := re.GroupNames(); len(g) != 2 || g[0] != "y" || g[1] != "m" {
		t.Errorf("GroupNames = %q", g)
	}
	m := re.Exec("on 2024-05!")
	if m == nil || m.Index != 3 || m.Groups[1] != "2024" || m.Names["m"] != 2 || m.Starts[2] != 8 || m.Ends[2] != 10 {
		t.Fatalf("Exec = %+v", m)
	}
	if got := re.ReplaceFunc("2024-05 and 1999-12", func(m *Match) string { return m.Groups[2] + "/" + m.Groups[1] }); got != "05/2024 and 12/1999" {
		t.Errorf("ReplaceFunc = %q", got)
	}
	if got := MustCompile("", "").Source(); got != "(?:)" {
		t.Errorf("empty source = %q", got)
	}
	if got := MustCompile("\n/[/]", "").Source(); got != `\n\/[/]` {
		t.Errorf("escaped source = %q", got)
	}
	_, err := Compile("(", "")
	if se, ok := err.(*SyntaxError); !ok || se.Message != "Invalid regular expression: /(/: Unterminated group" {
		t.Errorf("error = %v", err)
	}
	if _, err := Compile("a", "gg"); err == nil || err.Error() != "Invalid flags supplied to RegExp constructor 'gg'" {
		t.Errorf("flags error = %v", err)
	}
	defer func() {
		if recover() == nil {
			t.Error("MustCompile did not panic")
		}
	}()
	MustCompile("[", "")
}

// TestUTF16Indices checks that indices count UTF-16 code units and that lone
// surrogates (WTF-8) round-trip.
func TestUTF16Indices(t *testing.T) {
	m := MustCompile("b", "").Exec("😀b")
	if m.Index != 2 {
		t.Errorf("index = %d", m.Index)
	}
	lone := "\xed\xa0\xbd" // U+D83D
	if m := MustCompile(`\ud83d`, "").Exec("x" + lone); m == nil || m.Groups[0] != lone || m.Index != 1 {
		t.Errorf("lone surrogate match = %+v", m)
	}
	if got := MustCompile(".", "g").Match("😀"); len(got) != 2 || got[0] != "\xed\xa0\xbd" {
		t.Errorf("non-unicode . on astral = %q", got)
	}
}

// TestLongSubjects checks matching on 20k code-unit subjects (no Go stack
// overflow, reasonable time).
func TestLongSubjects(t *testing.T) {
	ab := strings.Repeat("ab", 10000)
	cases := []struct {
		pat, flags, s string
		want          bool
	}{
		{`^(a|b)*c$`, "", ab + "c", true},
		{`^(?:a|b)*$`, "", ab, true},
		{`^(?:(a)|b)+$`, "", ab, true},
		{`^(ab)+\1$`, "", ab, true},
		{`(?<=^(?:ab)*)$`, "", ab, true},
		{`^.*?x`, "s", ab, false},
		{`^[\s\S]*$`, "u", strings.Repeat("😀", 10000), true},
		{`^(?:\p{L}|\s)+$`, "u", strings.Repeat("ab ", 7000), true},
		{`(\w+)\s(\w+)$`, "", strings.Repeat("x", 20000) + " y", true},
	}
	for _, c := range cases {
		t0 := time.Now()
		got := MustCompile(c.pat, c.flags).Test(c.s)
		if got != c.want {
			t.Errorf("/%s/%s on %d units = %v", c.pat, c.flags, len(c.s), got)
		}
		if d := time.Since(t0); d > 20*time.Second {
			t.Errorf("/%s/%s took %v", c.pat, c.flags, d)
		}
	}
}

func libraryRegexes(tb testing.TB, n int) [][2]string {
	lines := gzLines(tb.(*testing.T), "testdata/chrono_library.jsonl.gz")
	var out [][2]string
	for _, line := range lines[1:] {
		var v vector
		if err := json.Unmarshal(line, &v); err != nil || v.Err.valid {
			continue
		}
		out = append(out, [2]string{v.P.s, v.F})
		if len(out) == n {
			break
		}
	}
	return out
}

// TestCompileSpeed compiles 200 of the regexes chrono builds (the chrono library corpus).
func TestCompileSpeed(t *testing.T) {
	res := libraryRegexes(t, 200)
	best := time.Hour
	for k := 0; k < 5; k++ {
		t0 := time.Now()
		for _, r := range res {
			MustCompile(r[0], r[1])
		}
		if d := time.Since(t0); d < best {
			best = d
		}
	}
	t.Logf("compiled %d library regexes in %v", len(res), best)
	if best > 5*time.Second {
		t.Errorf("compiling %d regexes took %v", len(res), best)
	}
}

func BenchmarkCompile200(b *testing.B) {
	t := &testing.T{}
	res := libraryRegexes(t, 200)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, r := range res {
			MustCompile(r[0], r[1])
		}
	}
}

func BenchmarkTestShortWord(b *testing.B) {
	re := MustCompile(`[aeiou](pp|ll|ss|ff|gg|tt|rr|bb|nn|mm)ed$`, "")
	words := []string{"walked", "stopped", "happy", "running", "tagged"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		re.Test(words[i%len(words)])
	}
}

func BenchmarkTestAnchored(b *testing.B) {
	re := MustCompile(`^[0-9]+(st|nd|rd|th)$`, "i")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		re.Test("3rd")
	}
}
