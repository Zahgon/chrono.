package ru_test

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"os"
	"slices"
	"strings"
	"testing"
	"time"

	chrono "github.com/wanasit/chrono-go"
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
	"github.com/wanasit/chrono-go/internal/jsdate"
	"github.com/wanasit/chrono-go/internal/ru"
	locale "github.com/wanasit/chrono-go/locales/ru"
)

// testdata/ru_calls.jsonl.gz holds Russian parse calls with the output of the original chrono
// (node, rendered by dev/print.js): {"call":{"id","tz","chrono","text","ref","opt","fn"},"want"}.
// Chrono names: ru.casual, ru.strict, ru.cfgCasual (new Chrono(createCasualConfiguration())),
// ru.cfgTrue / ru.cfgFalse (new Chrono(createConfiguration(b))), ru.casualRev / ru.strictRev
// (a clone with the parsers and refiners reversed). want is ERROR when the original throws.

func reversed(c *chrono.Chrono) *chrono.Chrono {
	c = c.Clone()
	slices.Reverse(c.Parsers)
	slices.Reverse(c.Refiners)
	return c
}

func chronos() map[string]*chrono.Chrono {
	return map[string]*chrono.Chrono{
		"ru.casual":    locale.Casual,
		"ru.strict":    locale.Strict,
		"ru.cfgCasual": chrono.NewChrono(locale.CreateCasualConfiguration()),
		"ru.cfgTrue":   chrono.NewChrono(locale.CreateConfiguration(true)),
		"ru.cfgFalse":  chrono.NewChrono(locale.CreateConfiguration(false)),
		"ru.casualRev": reversed(locale.Casual),
		"ru.strictRev": reversed(locale.Strict),
	}
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
			opt.Timezones[k] = chrono.TimezoneAbbr{Offset: tzs.Value(k).(float64)}
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
	f, err := os.Open("testdata/ru_calls.jsonl.gz")
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
	all := chronos()
	total, failed, errors := 0, 0, 0
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
				if e := catchJSError(func() {
					if fn, _ := r.call.Value("fn").(string); fn == "parseDate" {
						if d, ok := c.ParseDate(text, ref, opt); ok {
							got = num(d)
						} else {
							got = "null"
						}
					} else {
						got = renderParse(c.Parse(text, ref, opt))
					}
				}); e != nil {
					if e.Error() != "TypeError: Cannot read properties of undefined (reading 'match')" {
						t.Errorf("%v: unexpected error %v", r.call.Value("id"), e)
					}
					errors++
					got = "ERROR"
				}
				if got != r.want {
					failed++
					if failed <= 10 {
						t.Errorf("%v (TZ=%s, %s) %q:\n got %s\nwant %s", r.call.Value("id"), tz, name, text, got, r.want)
					}
				}
			}
		})
	}
	if total < 3000 || errors == 0 {
		t.Errorf("only %d recorded calls (%d errors)", total, errors)
	}
	if failed > 0 {
		t.Errorf("%d of %d calls differ", failed, total)
	}
}

// jsName maps a Go parser/refiner type to the constructor name of the original object.
func jsName(v any) string {
	name := fmt.Sprintf("%T", v)
	if after, ok := strings.CutPrefix(name, "*ru."); ok {
		return "RU" + after
	}
	return strings.TrimPrefix(name, "*core.")
}

func TestConfigurations(t *testing.T) {
	casualParsers := "RUTimeUnitCasualRelativeFormatParser,RURelativeDateFormatParser,RUMonthNameParser,RUCasualTimeParser,RUCasualDateParser,ISOFormatParser,SlashDateFormatParser,RUTimeUnitWithinFormatParser,RUMonthNameLittleEndianParser,RUWeekdayParser,RUTimeExpressionParser,RUTimeUnitAgoFormatParser"
	strictParsers := "ISOFormatParser,SlashDateFormatParser,RUTimeUnitWithinFormatParser,RUMonthNameLittleEndianParser,RUWeekdayParser,RUTimeExpressionParser,RUTimeUnitAgoFormatParser"
	refiners := "OverlapRemovalRefiner,ExtractTimezoneOffsetRefiner,MergeWeekdayComponentRefiner,RUMergeDateTimeRefiner,RUMergeDateRangeRefiner,ExtractTimezoneAbbrRefiner,OverlapRemovalRefiner,ForwardDateRefiner,UnlikelyFormatFilter"
	cases := []struct {
		name     string
		cfg      *chrono.Configuration
		parsers  string
		refiners string
	}{
		{"casual", &chrono.Configuration{Parsers: locale.Casual.Parsers, Refiners: locale.Casual.Refiners}, casualParsers, refiners},
		{"strict", &chrono.Configuration{Parsers: locale.Strict.Parsers, Refiners: locale.Strict.Refiners}, strictParsers, refiners},
		{"createCasualConfiguration()", locale.CreateCasualConfiguration(), casualParsers, refiners},
		{"createConfiguration(true)", locale.CreateConfiguration(true), strictParsers, refiners},
		{"createConfiguration(false)", locale.CreateConfiguration(false), strictParsers, refiners},
	}
	for _, c := range cases {
		var ps, rs []string
		for _, p := range c.cfg.Parsers {
			ps = append(ps, jsName(p))
		}
		for _, r := range c.cfg.Refiners {
			rs = append(rs, jsName(r))
		}
		if got := strings.Join(ps, ","); got != c.parsers {
			t.Errorf("%s parsers:\n got %s\nwant %s", c.name, got, c.parsers)
		}
		if got := strings.Join(rs, ","); got != c.refiners {
			t.Errorf("%s refiners:\n got %s\nwant %s", c.name, got, c.refiners)
		}
	}
	a, b := locale.CreateCasualConfiguration(), locale.CreateCasualConfiguration()
	for i := range a.Parsers {
		if a.Parsers[i] == b.Parsers[i] {
			t.Errorf("parser %d is shared between configurations", i)
		}
	}
}

// entries renders Object.entries(dict) as "key=value|...".
func entries[V any](d *core.Dictionary[V]) string {
	var parts []string
	for _, e := range d.Entries() {
		parts = append(parts, fmt.Sprintf("%s=%v", e.Key, e.Value))
	}
	return strings.Join(parts, "|")
}

func TestPatterns(t *testing.T) {
	// The exported patterns of src/locales/ru/constants.ts, as the original builds them.
	f, err := os.ReadFile("testdata/ru_patterns.json")
	if err != nil {
		t.Fatal(err)
	}
	v, err := js.Parse(string(f))
	if err != nil {
		t.Fatal(err)
	}
	want := v.(*js.Object)
	got := map[string]string{
		"NUMBER_PATTERN":         ru.NumberPattern,
		"ORDINAL_NUMBER_PATTERN": ru.OrdinalNumberPattern,
		"YEAR_PATTERN":           ru.YearPattern,
		"TIME_UNITS_PATTERN":     ru.TimeUnitsPattern,
		"leftBoundary":           ru.RegexParts.LeftBoundary,
		"rightBoundary":          ru.RegexParts.RightBoundary,
		"flags":                  ru.RegexParts.Flags,
	}
	for k, g := range got {
		if w, _ := want.Value(k).(string); g != w {
			t.Errorf("%s:\n got %s\nwant %s", k, g, w)
		}
	}
	dicts := map[string]string{
		"WEEKDAY_DICTIONARY":         entries(ru.WeekdayDictionary),
		"FULL_MONTH_NAME_DICTIONARY": entries(ru.FullMonthNameDictionary),
		"MONTH_DICTIONARY":           entries(ru.MonthDictionary),
		"INTEGER_WORD_DICTIONARY":    entries(ru.IntegerWordDictionary),
		"ORDINAL_WORD_DICTIONARY":    entries(ru.OrdinalWordDictionary),
		"TIME_UNIT_DICTIONARY":       entries(ru.TimeUnitDictionary),
	}
	for k, g := range dicts {
		if w, _ := want.Value(k).(string); g != w {
			t.Errorf("%s differs:\n got %s\nwant %s", k, g, w)
		}
	}
}
