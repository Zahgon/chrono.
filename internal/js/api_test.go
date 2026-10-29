package js

import (
	"errors"
	"math"
	"strings"
	"testing"
)

func TestToLocaleStringExamples(t *testing.T) {
	for _, c := range []struct {
		x    float64
		want string
	}{
		{1.0005, "1.001"}, {0.0125, "0.013"}, {2.5e-7, "0"}, {math.Copysign(0, -1), "-0"},
		{1e21, "1,000,000,000,000,000,000,000"}, {math.NaN(), "NaN"}, {math.Inf(1), "∞"}, {math.Inf(-1), "-∞"},
		{-1234.5678, "-1,234.568"}, {-0.0004, "-0"}, {-0.0005, "-0.001"}, {999.9995, "1,000"}, {1234567, "1,234,567"},
	} {
		if got := ToLocaleString(c.x); got != c.want {
			t.Errorf("ToLocaleString(%v) = %q want %q", c.x, got, c.want)
		}
	}
}

type rawJSON string

func (r rawJSON) AppendJSON(b []byte) []byte { return append(b, r...) }

func TestStringifyGoTypes(t *testing.T) {
	o := NewObject()
	o.Set("b", 1)
	o.Set("10", []string{"x", "\u2028"})
	o.Set("2", []float64{math.NaN(), math.Copysign(0, -1), 1e21})
	o.Set("a", Undefined)
	o.Set("01", []int{1, 2})
	o.Set("raw", rawJSON(`{"k":true}`))
	o.Set("u", Array{Undefined, nil, false})
	o.Set("b", "again")
	want := `{"2":[null,0,1e+21],"10":["x","` + "\u2028" + `"],"b":"again","01":[1,2],"raw":{"k":true},"u":[null,null,false]}`
	if got := Stringify(o); got != want {
		t.Errorf("Stringify = %s\nwant      %s", got, want)
	}
	if got := StringifyIndentN(Array{1, NewObject(), Array{}}, 2); got != "[\n  1,\n  {},\n  []\n]" {
		t.Errorf("StringifyIndentN = %q", got)
	}
	if got := Stringify(Undefined); got != "" {
		t.Errorf("Stringify(Undefined) = %q", got)
	}
	o.Delete("b")
	o.Set("b", 3)
	if got := strings.Join(o.Keys(), ","); got != "2,10,a,01,raw,u,b" {
		t.Errorf("Keys after delete/set = %s", got)
	}
	for _, k := range []string{"0", "4294967294"} {
		if _, ok := ArrayIndex(k); !ok {
			t.Errorf("%q should be an array index", k)
		}
	}
	for _, k := range []string{"", "01", "-1", "+1", "1.0", "4294967295", "99999999999"} {
		if _, ok := ArrayIndex(k); ok {
			t.Errorf("%q should not be an array index", k)
		}
	}
}

func expectThrow(t *testing.T, name, msg string, f func()) {
	t.Helper()
	defer func() {
		t.Helper()
		var e *Error
		r := recover()
		if err, ok := r.(error); !ok || !errors.As(err, &e) || e.Name != name || e.Message != msg {
			t.Errorf("got panic %v, want %s: %s", r, name, msg)
		}
	}()
	f()
}

func TestThrows(t *testing.T) {
	expectThrow(t, "RangeError", "toFixed() digits argument must be between 0 and 100", func() { ToFixed(1, 101) })
	expectThrow(t, "RangeError", "toPrecision() argument must be between 1 and 100", func() { ToPrecision(1, 0) })
	expectThrow(t, "RangeError", "toString() radix must be between 2 and 36", func() { NumberToStringRadix(1, 37) })
	expectThrow(t, "RangeError", "Invalid count value: -1", func() { Repeat("a", -1) })
	expectThrow(t, "RangeError", "Invalid code point 1114112", func() { FromCodePoint(0x110000) })
}

func TestNoPanics(t *testing.T) {
	rng := &mulberry32{a: 1}
	alphabet := []string{"{", "}", "[", "]", ",", ":", "\"", "\\", "u", "0", "1", "-", ".", "e", "t", "n", " ", "\x01", "\xed\xa0\x80", "\xff", "\u00e9", "\U0001F600", "\u03a3"}
	for i := 0; i < 20000; i++ {
		var b strings.Builder
		for n := rng.next() % 20; n > 0; n-- {
			b.WriteString(alphabet[rng.next()%uint32(len(alphabet))])
		}
		s := b.String()
		if v, err := Parse(s); err == nil {
			_ = Stringify(v)
		}
		_ = ToLowerCase(s) + ToUpperCase(s) + Quote(s) + Trim(s)
		_, _, _ = ToNumber(s), ParseFloat(s), ParseInt(s, int(rng.next()%40))
		k := int(rng.next()%30) - 10
		_ = Slice(s, k, -k) + Substring(s, k, 3) + Substr(s, k, k) + CharAt(s, k)
		_ = Split(s, CharAt(s, 0), -1)
		_ = IndexOf(s, "\xed\xa0\x80", k) + CodePointAt(s, k) + Compare(s, "a")
	}
}
