package js

import (
	"math"
	"testing"
)

// TestChronoStringVectors checks the String.prototype and global operations chrono uses against
// node (dev/js/gen_chrono_strings.mjs): toLowerCase, toUpperCase, trim, Number, parseFloat,
// parseInt, isNaN and replace/replaceAll with a string search value.
func TestChronoStringVectors(t *testing.T) {
	recs := readVectors(t, "chrono_strings.jsonl.gz")
	m := &mismatches{t: t, what: "chrono string operations"}
	for _, r := range recs {
		if str(r, "k") == "meta" {
			t.Logf("vectors from %s", str(r, "runtime"))
			continue
		}
		s := str(r, "s")
		m.check(ToLowerCase(s) == str(r, "lo"), "toLowerCase(%+q) = %+q want %+q", s, ToLowerCase(s), str(r, "lo"))
		m.check(ToUpperCase(s) == str(r, "up"), "toUpperCase(%+q) = %+q want %+q", s, ToUpperCase(s), str(r, "up"))
		m.check(Trim(s) == str(r, "tr"), "trim(%+q) = %+q want %+q", s, Trim(s), str(r, "tr"))
		m.check(sameNumber(ToNumber(s), str(r, "n")), "Number(%+q) = %v want %s", s, ToNumber(s), str(r, "n"))
		m.check(sameNumber(ParseFloat(s), str(r, "f")), "parseFloat(%+q) = %v want %s", s, ParseFloat(s), str(r, "f"))
		m.check(sameNumber(ParseInt(s, 0), str(r, "i")), "parseInt(%+q) = %v want %s", s, ParseInt(s, 0), str(r, "i"))
		m.check(sameNumber(ParseInt(s, 10), str(r, "i10")), "parseInt(%+q, 10) = %v want %s", s, ParseInt(s, 10), str(r, "i10"))
		nan, _ := r.Value("nan").(bool)
		m.check(IsNaN(ToNumber(s)) == nan, "isNaN(%+q) = %v want %v", s, IsNaN(ToNumber(s)), nan)
		for _, x := range arr(r, "rep") {
			a := x.(Array)
			search, repl := a[0].(string), a[1].(string)
			m.check(Replace(s, search, repl) == a[2].(string), "%+q.replace(%+q, %+q) = %+q want %+q", s, search, repl, Replace(s, search, repl), a[2])
			m.check(ReplaceAll(s, search, repl) == a[3].(string), "%+q.replaceAll(%+q, %+q) = %+q want %+q", s, search, repl, ReplaceAll(s, search, repl), a[3])
		}
	}
	m.done()
}

func TestIsFinite(t *testing.T) {
	for _, c := range []struct {
		x    float64
		want bool
	}{{0, true}, {-1.5, true}, {math.NaN(), false}, {math.Inf(1), false}, {math.Inf(-1), false}} {
		if IsFinite(c.x) != c.want {
			t.Errorf("IsFinite(%v) = %v", c.x, !c.want)
		}
	}
}
