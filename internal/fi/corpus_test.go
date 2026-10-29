package fi_test

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/fi"
	"github.com/wanasit/chrono-go/internal/js"
	"github.com/wanasit/chrono-go/internal/jsdate"
)

// testdata/fi_calls.jsonl.gz holds Finnish parse calls with the output of the original chrono
// (node, rendered by dev/print.js): {"call":{"id","tz","chrono","text","ref","opt","fn"},"want"}.
// chrono is "fi.casual", "fi.strict", "fi.ccX" (createCasualConfiguration(X)) or "fi.cXY"
// (createConfiguration(X, Y)), with X and Y "1" (true) or "0" (false).

func chronos() map[string]*core.Chrono {
	m := map[string]*core.Chrono{"fi.casual": fi.Casual, "fi.strict": fi.Strict}
	k := func(b bool) string {
		if b {
			return "1"
		}
		return "0"
	}
	for _, a := range []bool{true, false} {
		m["fi.cc"+k(a)] = core.NewChrono(fi.CreateCasualConfiguration(a))
		for _, b := range []bool{true, false} {
			m["fi.c"+k(a)+k(b)] = core.NewChrono(fi.CreateConfiguration(a, b))
		}
	}
	return m
}

func TestRecordedCalls(t *testing.T) {
	replay(t, "testdata/fi_calls.jsonl.gz", chronos())
}

func TestConfigurationOrder(t *testing.T) {
	names := func(xs any) string {
		var out []string
		switch v := xs.(type) {
		case []core.Parser:
			for _, x := range v {
				out = append(out, fmt.Sprintf("%T", x))
			}
		case []core.Refiner:
			for _, x := range v {
				out = append(out, fmt.Sprintf("%T", x))
			}
		}
		return strings.Join(out, " ")
	}
	c := fi.CreateCasualConfiguration(true)
	wantParsers := "*fi.TimeUnitCasualRelativeFormatParser *fi.CasualDateParser *fi.CasualTimeParser *core.ISOFormatParser " +
		"*core.ISOFormatParser *core.SlashDateFormatParser *fi.TimeExpressionParser *fi.MonthNameLittleEndianParser " +
		"*fi.WeekdayParser *fi.TimeUnitWithinFormatParser *fi.TimeUnitAgoFormatParser"
	if got := names(c.Parsers); got != wantParsers {
		t.Errorf("casual parsers:\n got %s\nwant %s", got, wantParsers)
	}
	wantRefiners := "*core.OverlapRemovalRefiner *core.ExtractTimezoneOffsetRefiner *core.MergeWeekdayComponentRefiner " +
		"*fi.MergeDateRangeRefiner *fi.MergeDateTimeRefiner *core.ExtractTimezoneAbbrRefiner *core.OverlapRemovalRefiner " +
		"*core.ForwardDateRefiner *core.UnlikelyFormatFilter"
	if got := names(c.Refiners); got != wantRefiners {
		t.Errorf("casual refiners:\n got %s\nwant %s", got, wantRefiners)
	}
	if s := fi.CreateConfiguration(true, false); len(s.Parsers) != 8 || len(s.Refiners) != 9 {
		t.Errorf("strict: %d parsers, %d refiners", len(s.Parsers), len(s.Refiners))
	}
	if a, b := fi.CreateCasualConfiguration(true), fi.CreateCasualConfiguration(true); a.Parsers[0] == b.Parsers[0] {
		t.Error("configurations share parser objects")
	}
}

// ---------------------------------------------------------------------------------- replay

var renderComponents = []core.Component{"year", "month", "day", "weekday", "hour", "minute", "second", "millisecond", "meridiem", "timezoneOffset"}

func num(f float64) string { return js.NumberToString(f) }

func strList(xs []string) string {
	parts := make([]string, len(xs))
	for i, x := range xs {
		parts[i] = js.Quote(x)
	}
	return "[" + strings.Join(parts, ",") + "]"
}

func renderComponentsOf(c *core.ParsingComponents) string {
	if c == nil {
		return "null"
	}
	vals := make([]string, len(renderComponents))
	for i, k := range renderComponents {
		vs := "null"
		if v, ok := c.Get(k); ok {
			vs = num(v)
		}
		cert := 0
		if c.IsCertain(k) {
			cert = 1
		}
		vals[i] = fmt.Sprintf("[%s,%d]", vs, cert)
	}
	var certain []string
	for _, x := range c.GetCertainComponents() {
		certain = append(certain, string(x))
	}
	return "[[" + strings.Join(vals, ",") + "]," + strList(certain) + "," + num(c.Date()) + "," + strList(c.Tags()) + "]"
}

func renderParse(results []*core.ParsingResult) string {
	parts := make([]string, len(results))
	for i, r := range results {
		parts[i] = "[" + num(float64(r.Index)) + "," + js.Quote(r.Text) + "," + renderComponentsOf(r.Start) + "," +
			renderComponentsOf(r.End) + "," + strList(r.Tags()) + "]"
	}
	return "[" + strings.Join(parts, ",") + "]"
}

func decodeCall(c *js.Object) (text string, ref core.ParsingReference, opt *core.ParsingOption) {
	text, _ = c.Value("text").(string)
	switch r := c.Value("ref").(type) {
	case float64:
		ref.Instant = r
	case *js.Object:
		ref.Instant, _ = r.Value("instant").(float64)
		if tz, ok := r.Get("timezone"); ok {
			ref.Timezone = tz
		}
	}
	o, ok := c.Value("opt").(*js.Object)
	if !ok {
		return text, ref, nil
	}
	opt = &core.ParsingOption{}
	opt.ForwardDate, _ = o.Value("forwardDate").(bool)
	if tzs, ok := o.Value("timezones").(*js.Object); ok {
		opt.Timezones = map[string]core.TimezoneAbbr{}
		for _, k := range tzs.Keys() {
			if v, ok := tzs.Value(k).(float64); ok {
				opt.Timezones[k] = core.TimezoneAbbr{Offset: v}
			}
		}
	}
	return text, ref, opt
}

func inZone(t *testing.T, name string, fn func()) {
	t.Helper()
	loc, err := time.LoadLocation(name)
	if err != nil {
		t.Fatalf("LoadLocation(%q): %v", name, err)
	}
	prev := jsdate.SetLocation(loc)
	defer jsdate.SetLocation(prev)
	fn()
}

func catchJSError(fn func()) (err *core.Error) {
	defer func() {
		if r := recover(); r != nil {
			e, ok := r.(*core.Error)
			if !ok {
				panic(r)
			}
			err = e
		}
	}()
	fn()
	return nil
}

func replay(t *testing.T, path string, all map[string]*core.Chrono) {
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	zr, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	type rec struct {
		call *js.Object
		want string
	}
	byZone := map[string][]rec{}
	var zones []string
	sc := bufio.NewScanner(zr)
	sc.Buffer(make([]byte, 1<<20), 1<<26)
	for sc.Scan() {
		v, err := js.Parse(sc.Text())
		if err != nil {
			t.Fatal(err)
		}
		o := v.(*js.Object)
		call := o.Value("call").(*js.Object)
		tz, _ := call.Value("tz").(string)
		if _, seen := byZone[tz]; !seen {
			zones = append(zones, tz)
		}
		byZone[tz] = append(byZone[tz], rec{call, o.Value("want").(string)})
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}
	total, failed := 0, 0
	for _, tz := range zones {
		inZone(t, tz, func() {
			for _, r := range byZone[tz] {
				total++
				name, _ := r.call.Value("chrono").(string)
				c, ok := all[name]
				if !ok {
					t.Fatalf("unknown chrono %q", name)
				}
				text, ref, opt := decodeCall(r.call)
				var got string
				if err := catchJSError(func() {
					if fn, _ := r.call.Value("fn").(string); fn == "parseDate" {
						if d, ok := c.ParseDate(text, ref, opt); ok {
							got = num(d)
						} else {
							got = "null"
						}
					} else {
						got = renderParse(c.Parse(text, ref, opt))
					}
				}); err != nil {
					got = "ERROR"
				}
				if got != r.want {
					failed++
					if failed <= 10 {
						t.Errorf("%v (TZ=%s) %q:\n got %s\nwant %s", r.call.Value("id"), tz, text, got, r.want)
					}
				}
			}
		})
	}
	if total < 3000 {
		t.Errorf("only %d recorded calls", total)
	}
	if failed > 0 {
		t.Errorf("%d of %d calls differ", failed, total)
	}
}
