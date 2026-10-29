package zh

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
	"github.com/wanasit/chrono-go/internal/jsdate"

	"github.com/wanasit/chrono-go/internal/zh/hans"
	"github.com/wanasit/chrono-go/internal/zh/hant"
)

// testdata/zh_calls.jsonl.gz holds parse calls with the output of the original chrono (node, orig/dist/cjs,
// rendered by dev/print.js): {"call":{"id","tz","chrono","text","ref","opt","fn"},"want"}.
// Chrono names "<locale>.new.casual" / "<locale>.new.cfg[.<strictMode>]" are new Chrono objects
// built from the configuration constructors.

var newChronos = map[string]*core.Chrono{}

func resolveChrono(t *testing.T, name string) *core.Chrono {
	switch name {
	case "zh.casual":
		return Casual
	case "zh.strict":
		return Strict
	case "zh.hans.casual":
		return hans.Casual
	case "zh.hans.strict":
		return hans.Strict
	case "zh.hant.casual":
		return hant.Casual
	case "zh.hant.strict":
		return hant.Strict
	}
	if c, ok := newChronos[name]; ok {
		return c
	}
	var c *core.Chrono
	switch name {
	case "zh.new.casual":
		c = core.NewChrono(CreateCasualConfiguration())
	case "zh.new.cfg":
		c = core.NewChrono(CreateConfiguration())
	case "zh.hans.new.casual":
		c = core.NewChrono(hans.CreateCasualConfiguration())
	case "zh.hans.new.cfg":
		c = core.NewChrono(hans.CreateConfiguration())
	case "zh.hant.new.casual":
		c = core.NewChrono(hant.CreateCasualConfiguration())
	case "zh.hant.new.cfg":
		c = core.NewChrono(hant.CreateConfiguration())
	default:
		t.Fatalf("unknown chrono %q", name)
	}
	newChronos[name] = c
	return c
}

// The canonical rendering of dev/print.js.

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

// catchJSError runs fn and returns the *core.Error it panics with (nil when it returns).
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

func TestRecordedCalls(t *testing.T) {
	f, err := os.Open("testdata/zh_calls.jsonl.gz")
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
		loc, err := time.LoadLocation(tz)
		if err != nil {
			t.Fatalf("LoadLocation(%q): %v", tz, err)
		}
		prev := jsdate.SetLocation(loc)
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
				got = "ERROR"
			}
			if got != r.want {
				failed++
				if failed <= 10 {
					t.Errorf("%v (TZ=%s) %q:\n got %s\nwant %s", r.call.Value("id"), tz, text, got, r.want)
				}
			}
		}
		jsdate.SetLocation(prev)
	}
	if total < 7000 {
		t.Errorf("only %d recorded calls", total)
	}
	if failed > 0 {
		t.Errorf("%d of %d calls differ", failed, total)
	}
}
