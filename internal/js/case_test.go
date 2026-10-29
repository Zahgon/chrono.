package js

import (
	"testing"
)

func TestCaseVectors(t *testing.T) {
	recs := readVectors(t, "case.jsonl.gz")
	cpm := &mismatches{t: t, what: "code point mapping"}
	sig := &mismatches{t: t, what: "final sigma probes"}
	ctx := &mismatches{t: t, what: "context strings"}
	listed := map[int]bool{}
	for _, r := range recs {
		switch str(r, "k") {
		case "meta":
			t.Logf("vectors from %s", str(r, "runtime"))
		case "map":
			cp := int(num(r, "cp"))
			listed[cp] = true
			s := FromCodePoint(cp)
			l, u := ToLowerCase(s), ToUpperCase(s)
			cpm.check(l == str(r, "l") && u == str(r, "u"), "U+%04X: got %q %q want %q %q", cp, l, u, str(r, "l"), str(r, "u"))
		case "sigma":
			for _, run := range arr(r, "runs") {
				rr := run.(Array)
				lo, hi, code := int(rr[0].(float64)), int(rr[1].(float64)), int(rr[2].(float64))
				for cp := lo; cp <= hi; cp++ {
					x := FromCodePoint(cp)
					probes := [4]string{"A\u03a3" + x, "A\u03a3" + x + "B", x + "\u03a3", "A" + x + "\u03a3"}
					got := 0
					for i, p := range probes {
						u := CodeUnits(ToLowerCase(Concat(p)))
						var c uint16
						if i < 2 {
							c = u[1]
						} else {
							c = u[len(u)-1]
						}
						if c == 0x3c2 {
							got |= 1 << i
						}
					}
					sig.check(got == code, "U+%04X: probe bits %04b want %04b", cp, got, code)
				}
			}
		case "str":
			s := str(r, "s")
			l, u := ToLowerCase(s), ToUpperCase(s)
			ctx.check(l == str(r, "l") && u == str(r, "u"), "%+q: got %+q %+q want %+q %+q", s, l, u, str(r, "l"), str(r, "u"))
		}
	}
	for cp := 0; cp < 0x110000; cp++ {
		if listed[cp] {
			continue
		}
		s := FromCodePoint(cp)
		cpm.check(ToLowerCase(s) == s && ToUpperCase(s) == s, "U+%04X should map to itself", cp)
	}
	cpm.done()
	sig.done()
	ctx.done()
}
