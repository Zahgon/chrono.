package chrono_test

import (
	"bufio"
	"compress/gzip"
	"os"
	"testing"

	chrono "github.com/wanasit/chrono-go"
	"github.com/wanasit/chrono-go/internal/js"
	"github.com/wanasit/chrono-go/locales/en"
)

// testdata/en_calls.jsonl.gz holds English parse calls with the output of the original chrono
// (node, rendered by dev/print.js): {"call":{"id","tz","chrono","text","ref","opt","fn"},"want"}.
// Ambiguous time zone overrides are encoded as
// {"dst":m,"nonDst":m,"start":["nth",month,weekday,n,hour]|["last",month,weekday,hour],"end":...}.

func resolveChrono(t *testing.T, name string) *chrono.Chrono {
	switch name {
	case "", "casual":
		return chrono.Casual
	case "strict":
		return chrono.Strict
	case "en.GB":
		return en.GB
	}
	t.Fatalf("unknown chrono %q", name)
	return nil
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

func TestRecordedEnglishCalls(t *testing.T) {
	f, err := os.Open("testdata/en_calls.jsonl.gz")
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
	if total < 5000 {
		t.Errorf("only %d recorded calls", total)
	}
	if failed > 0 {
		t.Errorf("%d of %d calls differ", failed, total)
	}
}
