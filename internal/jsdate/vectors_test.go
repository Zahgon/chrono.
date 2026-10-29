package jsdate

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strings"
	"testing"
	"time"
)

// jsNum decodes a number recorded by dev/jsdate/gen_jsdate.mjs: a JSON number, or the string
// "NaN", "Infinity", "-Infinity" or "-0".
func jsNum(v any) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case string:
		switch x {
		case "NaN":
			return math.NaN()
		case "Infinity":
			return math.Inf(1)
		case "-Infinity":
			return math.Inf(-1)
		case "-0":
			return math.Copysign(0, -1)
		}
	}
	panic(fmt.Sprintf("bad number %v", v))
}

// same compares numbers as Object.is does, except that -0 and +0 are equal when want is +0
// recorded as a plain 0 (JSON loses nothing: -0 is recorded as "-0").
func same(got float64, want any) bool {
	w := jsNum(want)
	if math.IsNaN(w) {
		return math.IsNaN(got)
	}
	return got == w && math.Signbit(got) == math.Signbit(w)
}

type vecCase struct {
	Z   string
	Ops [][]any
	R   [][]any
	U   []any
	ISO *string
}

func setter(d *Date, name string, a []float64) float64 {
	switch name {
	case "setFullYear":
		return d.SetFullYear(a[0], a[1:]...)
	case "setMonth":
		return d.SetMonth(a[0], a[1:]...)
	case "setDate":
		return d.SetDate(a[0])
	case "setHours":
		return d.SetHours(a[0], a[1:]...)
	case "setMinutes":
		return d.SetMinutes(a[0], a[1:]...)
	case "setSeconds":
		return d.SetSeconds(a[0], a[1:]...)
	case "setMilliseconds":
		return d.SetMilliseconds(a[0])
	case "setTime":
		return d.SetTime(a[0])
	case "setUTCFullYear":
		return d.SetUTCFullYear(a[0], a[1:]...)
	case "setUTCMonth":
		return d.SetUTCMonth(a[0], a[1:]...)
	case "setUTCDate":
		return d.SetUTCDate(a[0])
	case "setUTCHours":
		return d.SetUTCHours(a[0], a[1:]...)
	case "setUTCMinutes":
		return d.SetUTCMinutes(a[0], a[1:]...)
	case "setUTCSeconds":
		return d.SetUTCSeconds(a[0], a[1:]...)
	case "setUTCMilliseconds":
		return d.SetUTCMilliseconds(a[0])
	}
	panic("unknown setter " + name)
}

// TestVectors replays the operation chains recorded from node (V8) in every grading time zone
// (dev/jsdate/gen_jsdate.mjs; JSDATE_VECTORS overrides the file). Each case sets the local zone
// from its zone name with time.LoadLocation.
func TestVectors(t *testing.T) {
	path := "testdata/vectors.jsonl.gz"
	if p := os.Getenv("JSDATE_VECTORS"); p != "" {
		path = p
	}
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
	sc.Buffer(make([]byte, 1<<20), 1<<26)
	if !sc.Scan() {
		t.Fatal("empty vector file")
	}
	var meta struct {
		Meta  string
		Zones []string
	}
	if err := json.Unmarshal(sc.Bytes(), &meta); err != nil {
		t.Fatal(err)
	}
	t.Logf("vectors from %s", meta.Meta)
	defer SetLocation(SetLocation(nil))
	locs := map[string]*time.Location{}
	type stat struct{ chains, ops, bad int }
	stats := map[string]*stat{}
	var fails []string
	for sc.Scan() {
		var c vecCase
		if err := json.Unmarshal(sc.Bytes(), &c); err != nil {
			t.Fatal(err)
		}
		loc, ok := locs[c.Z]
		if !ok {
			if loc, err = time.LoadLocation(c.Z); err != nil {
				t.Fatal(err)
			}
			locs[c.Z] = loc
			stats[c.Z] = &stat{}
		}
		SetLocation(loc)
		st := stats[c.Z]
		st.chains++
		var errs []string
		check := func(what string, got float64, want any) {
			st.ops++
			if !same(got, want) {
				errs = append(errs, fmt.Sprintf("%s = %v, want %v", what, got, want))
			}
		}
		var d *Date
		for i, op := range c.Ops {
			name := op[0].(string)
			args := make([]float64, len(op)-1)
			for k := range args {
				args[k] = jsNum(op[k+1])
			}
			var ret float64
			switch name {
			case "L":
				d = NewLocal(args...)
				ret = d.Time()
			case "T":
				d = FromTime(args[0])
				ret = d.Time()
			case "U":
				ret = UTC(args...)
				d = FromTime(ret)
			default:
				ret = setter(d, name, args)
			}
			r := c.R[i]
			label := fmt.Sprintf("op %d %v", i, op)
			check(label+" result", ret, r[0])
			if len(r) == 1 {
				if d.IsValid() {
					errs = append(errs, label+": date valid, want invalid")
				}
				for _, g := range []float64{d.GetFullYear(), d.GetMonth(), d.GetDate(), d.GetHours(), d.GetMinutes(), d.GetSeconds(), d.GetMilliseconds(), d.GetDay(), d.GetTimezoneOffset()} {
					check(label+" getter on invalid date", g, "NaN")
				}
				continue
			}
			got := []float64{d.GetFullYear(), d.GetMonth(), d.GetDate(), d.GetHours(), d.GetMinutes(), d.GetSeconds(), d.GetMilliseconds(), d.GetDay(), d.GetTimezoneOffset()}
			names := []string{"getFullYear", "getMonth", "getDate", "getHours", "getMinutes", "getSeconds", "getMilliseconds", "getDay", "getTimezoneOffset"}
			for k, g := range got {
				check(label+" "+names[k], g, r[k+1])
			}
		}
		ug := []float64{d.GetUTCFullYear(), d.GetUTCMonth(), d.GetUTCDate(), d.GetUTCHours(), d.GetUTCMinutes(), d.GetUTCSeconds(), d.GetUTCMilliseconds(), d.GetUTCDay()}
		for k, g := range ug {
			check(fmt.Sprintf("UTC getter %d", k), g, c.U[k])
		}
		iso, ok := d.ToISOString()
		st.ops++
		if (c.ISO == nil) != !ok || (ok && iso != *c.ISO) {
			errs = append(errs, fmt.Sprintf("toISOString = %q,%v want %v", iso, ok, c.ISO))
		}
		if len(errs) > 0 {
			st.bad++
			if len(fails) < 30 {
				fails = append(fails, fmt.Sprintf("[%s] %v:\n    %s", c.Z, c.Ops, strings.Join(errs, "\n    ")))
			}
		}
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}
	total, bad := 0, 0
	for _, z := range meta.Zones {
		st := stats[z]
		if st == nil {
			t.Errorf("no vectors for %s", z)
			continue
		}
		t.Logf("%-20s %6d chains %8d operations %d mismatching chains", z, st.chains, st.ops, st.bad)
		total += st.ops
		bad += st.bad
	}
	t.Logf("total %d operations, %d mismatching chains", total, bad)
	if bad > 0 {
		t.Errorf("mismatches:\n%s", strings.Join(fails, "\n"))
	}
}
