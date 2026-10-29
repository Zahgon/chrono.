package js

import (
	"fmt"
	"strconv"
	"testing"
)

// miscOutputs mirrors the output list of dev/js/gen_misc.mjs.
func miscOutputs(s, t string, a, b int, fill string) []any {
	sa := func(x []string) Array {
		out := Array{}
		for _, v := range x {
			out = append(out, v)
		}
		return out
	}
	return []any{
		Len(s), Slice(s, a, b), SliceFrom(s, a), Substring(s, a, b), SubstringFrom(s, a), Substr(s, a, b), SubstrFrom(s, a),
		IndexOf(s, t, a), LastIndexOf(s, t), LastIndexOfFrom(s, t, a), Includes(s, t), StartsWith(s, t, a), EndsWith(s, t), EndsWithAt(s, t, a),
		sa(Split(s, t, -1)), sa(Split(s, t, max(a, 0))), Trim(s), TrimStart(s), TrimEnd(s), PadStart(s, b+5, fill), PadEnd(s, b+5, fill),
		Repeat(s, max(b, 0)%4), CharAt(s, a), CharCodeAt(s, a), CodePointAt(s, a),
		Compare(s, t), Concat(s, t), IsWellFormed(s), ToWellFormed(s),
	}
}

func sameValue(got, want any) bool {
	switch g := got.(type) {
	case int:
		w, ok := want.(float64)
		return ok && float64(g) == w
	case Array:
		w, ok := want.(Array)
		if !ok || len(w) != len(g) {
			return false
		}
		for i := range g {
			if !sameValue(g[i], w[i]) {
				return false
			}
		}
		return true
	}
	return got == want
}

func TestMiscVectors(t *testing.T) {
	recs := readVectors(t, "misc.jsonl.gz")
	sm := &mismatches{t: t, what: "string methods"}
	nm := &mismatches{t: t, what: "Math.round / ToInt32 / ToUint32"}
	for _, r := range recs {
		switch str(r, "k") {
		case "meta":
			t.Logf("vectors from %s", str(r, "runtime"))
		case "str":
			s, tt, fill := str(r, "s"), str(r, "t"), str(r, "fill")
			a, b := int(num(r, "a")), int(num(r, "b"))
			got := miscOutputs(s, tt, a, b, fill)
			for j, w := range arr(r, "o") {
				sm.check(sameValue(got[j], w), "#%d s=%+q t=%+q a=%d b=%d fill=%+q: got %+q want %+q", j, s, tt, a, b, fill, fmt.Sprint(got[j]), fmt.Sprint(w))
			}
		case "num":
			x, _ := strconv.ParseFloat(str(r, "x"), 64)
			o := arr(r, "o")
			nm.check(sameNumber(MathRound(x), o[0].(string)), "Math.round(%v) = %v want %v", x, MathRound(x), o[0])
			nm.check(float64(ToInt32(x)) == o[1].(float64), "ToInt32(%v) = %v want %v", x, ToInt32(x), o[1])
			nm.check(float64(ToUint32(x)) == o[2].(float64), "ToUint32(%v) = %v want %v", x, ToUint32(x), o[2])
		}
	}
	sm.done()
	nm.done()
}
