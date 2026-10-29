package chrono_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	chrono "github.com/wanasit/chrono-go"
	"github.com/wanasit/chrono-go/internal/js"
	"github.com/wanasit/chrono-go/internal/jsdate"
)

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
