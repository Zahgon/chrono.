package de_test

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	chrono "github.com/wanasit/chrono-go"
	"github.com/wanasit/chrono-go/internal/js"
	"github.com/wanasit/chrono-go/internal/jsdate"
	"github.com/wanasit/chrono-go/locales/de"
)

// testdata/calls.jsonl.gz holds German parse calls with the output of the original chrono (node,
// rendered by dev/print.js): {"call":{"id","tz","chrono","text","ref","opt","fn"},"want"}.
// chrono is "de.casual", "de.strict", "de.ccL" (createCasualConfiguration(L)) or "de.cSL"
// (createConfiguration(S, L)), L and S being 0 or 1. A call that throws records
// "ERROR <name>: <message>". Ambiguous time zone overrides are encoded as
// {"dst":m,"nonDst":m,"start":["nth",month,weekday,n,hour]|["last",month,weekday,hour],"end":...}.

var chronoCache = map[string]*chrono.Chrono{}

func resolveChrono(t *testing.T, name string) *chrono.Chrono {
	if c, ok := chronoCache[name]; ok {
		return c
	}
	_, k, _ := strings.Cut(name, ".")
	var c *chrono.Chrono
	switch {
	case k == "casual":
		c = de.Casual
	case k == "strict":
		c = de.Strict
	case len(k) == 3 && k[:2] == "cc":
		c = chrono.NewChrono(de.CreateCasualConfiguration(k[2] == '1'))
	case len(k) == 3 && k[0] == 'c':
		c = chrono.NewChrono(de.CreateConfiguration(k[1] == '1', k[2] == '1'))
	default:
		t.Fatalf("unknown chrono %q", name)
	}
	chronoCache[name] = c
	return c
}

func decodeRule(a js.Array) func(float64) float64 {
	n := func(i int) int { return int(a[i].(float64)) }
	if a[0].(string) == "nth" {
		return func(y float64) float64 { return chrono.GetNthWeekdayOfMonth(y, n(1), n(2), n(3), n(4)) }
	}
	return func(y float64) float64 { return chrono.GetLastWeekdayOfMonth(y, n(1), n(2), n(3)) }
}

func decodeCall(c *js.Object) (text string, ref chrono.ParsingReference, opt *chrono.ParsingOption) {
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
	opt = &chrono.ParsingOption{}
	opt.ForwardDate, _ = o.Value("forwardDate").(bool)
	if tzs, ok := o.Value("timezones").(*js.Object); ok {
		opt.Timezones = map[string]chrono.TimezoneAbbr{}
		for _, k := range tzs.Keys() {
			switch v := tzs.Value(k).(type) {
			case float64:
				opt.Timezones[k] = chrono.TimezoneAbbr{Offset: v}
			case *js.Object:
				opt.Timezones[k] = chrono.TimezoneAbbr{Ambiguous: &chrono.AmbiguousTimezoneMap{
					TimezoneOffsetDuringDst: v.Value("dst").(float64),
					TimezoneOffsetNonDst:    v.Value("nonDst").(float64),
					DstStart:                decodeRule(v.Value("start").(js.Array)),
					DstEnd:                  decodeRule(v.Value("end").(js.Array)),
				}}
			}
		}
	}
	return text, ref, opt
}

// The canonical rendering of dev/print.js.

var renderComponents = []chrono.Component{"year", "month", "day", "weekday", "hour", "minute", "second", "millisecond", "meridiem", "timezoneOffset"}

func num(f float64) string { return js.NumberToString(f) }

func strList(xs []string) string {
	parts := make([]string, len(xs))
	for i, x := range xs {
		parts[i] = js.Quote(x)
	}
	return "[" + strings.Join(parts, ",") + "]"
}

func renderComponentsOf(c *chrono.ParsingComponents) string {
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

func renderParse(results []*chrono.ParsingResult) string {
	parts := make([]string, len(results))
	for i, r := range results {
		parts[i] = "[" + num(float64(r.Index)) + "," + js.Quote(r.Text) + "," + renderComponentsOf(r.Start) + "," +
			renderComponentsOf(r.End) + "," + strList(r.Tags()) + "]"
	}
	return "[" + strings.Join(parts, ",") + "]"
}

// inZone runs fn with the JavaScript local time zone set to name.
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

// catchJSError runs fn and returns the *chrono.Error it panics with (nil when it returns).
func catchJSError(fn func()) (err *chrono.Error) {
	defer func() {
		if r := recover(); r != nil {
			e, ok := r.(*chrono.Error)
			if !ok {
				panic(r)
			}
			err = e
		}
	}()
	fn()
	return nil
}

func TestRecordedCalls(t *testing.T) {
	f, err := os.Open("testdata/calls.jsonl.gz")
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
				c := resolveChrono(t, name)
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
					got = "ERROR " + err.Error()
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
	if total < 2000 {
		t.Errorf("only %d recorded calls", total)
	}
	if failed > 0 {
		t.Errorf("%d of %d calls differ", failed, total)
	}
}

// jsNames are the constructor names of the original configurations' parsers and refiners.
var jsNames = map[string]string{
	"casual": "DETimeUnitAgoFormatParser DECasualDateParser DECasualTimeParser ISOFormatParser ISOFormatParser SlashDateFormatParser DETimeExpressionParser DESpecificTimeExpressionParser DEMonthNameLittleEndianParser DEWeekdayParser DETimeUnitWithinFormatParser | OverlapRemovalRefiner ExtractTimezoneOffsetRefiner MergeWeekdayComponentRefiner DEMergeDateRangeRefiner DEMergeDateTimeRefiner ExtractTimezoneAbbrRefiner OverlapRemovalRefiner ForwardDateRefiner UnlikelyFormatFilter",
	"strict": "ISOFormatParser ISOFormatParser SlashDateFormatParser DETimeExpressionParser DESpecificTimeExpressionParser DEMonthNameLittleEndianParser DEWeekdayParser DETimeUnitWithinFormatParser | OverlapRemovalRefiner ExtractTimezoneOffsetRefiner MergeWeekdayComponentRefiner DEMergeDateRangeRefiner DEMergeDateTimeRefiner ExtractTimezoneAbbrRefiner OverlapRemovalRefiner ForwardDateRefiner UnlikelyFormatFilter",
}

// goToJSName maps a Go parser/refiner type to the original class name.
func goToJSName(v any) string {
	name := fmt.Sprintf("%T", v)
	pkg, typ, _ := strings.Cut(strings.TrimPrefix(name, "*"), ".")
	if pkg == "core" {
		return typ
	}
	if typ == "TimeUnitRelativeFormatParser" {
		return "DETimeUnitAgoFormatParser" // the class name in DETimeUnitRelativeFormatParser.ts
	}
	return strings.ToUpper(pkg) + typ
}

func configNames(c *chrono.Configuration) string {
	var p, r []string
	for _, x := range c.Parsers {
		p = append(p, goToJSName(x))
	}
	for _, x := range c.Refiners {
		r = append(r, goToJSName(x))
	}
	return strings.Join(p, " ") + " | " + strings.Join(r, " ")
}

// TestConfigurationClassOrder compares the parser/refiner classes and order with the original.
func TestConfigurationClassOrder(t *testing.T) {
	for _, le := range []bool{true, false} {
		if got := configNames(de.CreateCasualConfiguration(le)); got != jsNames["casual"] {
			t.Errorf("casual(%v):\n got %s\nwant %s", le, got, jsNames["casual"])
		}
		for _, strict := range []bool{true, false} {
			if got := configNames(de.CreateConfiguration(strict, le)); got != jsNames["strict"] {
				t.Errorf("config(%v, %v):\n got %s\nwant %s", strict, le, got, jsNames["strict"])
			}
		}
	}
	if len(de.Casual.Parsers) != len(de.CreateCasualConfiguration(true).Parsers) || len(de.Strict.Refiners) != len(de.CreateConfiguration(true, true).Refiners) {
		t.Error("Casual/Strict do not hold the configurations")
	}
}
