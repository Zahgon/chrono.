package chrono_test

import (
	"bufio"
	"os"
	"slices"
	"strings"
	"testing"

	chrono "github.com/wanasit/chrono-go"
	"github.com/wanasit/chrono-go/internal/js"
)

// funcParser is a custom Parser built from two functions (a JS object literal parser).
type funcParser struct {
	pattern func(context *chrono.ParsingContext) *chrono.RegExp
	extract func(context *chrono.ParsingContext, match *chrono.RegExpMatch) any
}

func (p *funcParser) Pattern(context *chrono.ParsingContext) *chrono.RegExp {
	return p.pattern(context)
}
func (p *funcParser) Extract(context *chrono.ParsingContext, match *chrono.RegExpMatch) any {
	return p.extract(context, match)
}

// funcRefiner is a custom Refiner built from a function.
type funcRefiner func(context *chrono.ParsingContext, results []*chrono.ParsingResult) []*chrono.ParsingResult

func (f funcRefiner) Refine(context *chrono.ParsingContext, results []*chrono.ParsingResult) []*chrono.ParsingResult {
	return f(context, results)
}

func mustRegExp(t *testing.T, source, flags string) *chrono.RegExp {
	re, err := chrono.NewRegExp(source, flags)
	if err != nil {
		t.Fatal(err)
	}
	return re
}

// probeScenarios mirror the JavaScript scenarios recorded in testdata/probes.txt.
func probeScenarios(t *testing.T) map[string]func() []*chrono.ParsingResult {
	ref := chrono.ParsingReference{Instant: 1700000000000}
	return map[string]func() []*chrono.ParsingResult{
		"remove-iso": func() []*chrono.ParsingResult {
			c := chrono.Casual.Clone()
			c.Parsers = slices.Delete(c.Parsers, 1, 2)
			return c.Parse("2021-03-04T05:06:07Z and tomorrow 2021-03-04", ref, nil)
		},
		"reverse-refiners": func() []*chrono.ParsingResult {
			c := chrono.Casual.Clone()
			slices.Reverse(c.Refiners)
			return c.Parse("Monday 5pm to Friday 3pm EST, Sat 12/7/2014", ref, nil)
		},
		"no-refiners": func() []*chrono.ParsingResult {
			c := chrono.Strict.Clone()
			c.Refiners = nil
			return c.Parse("March 5 at 10:30 - 11:45 PST", ref, nil)
		},
		"custom-components": func() []*chrono.ParsingResult {
			c := chrono.Casual.Clone()
			c.Parsers = append(c.Parsers, &funcParser{
				pattern: func(*chrono.ParsingContext) *chrono.RegExp { return mustRegExp(t, "(\\d+)\\s*sleeps?", "i") },
				extract: func(_ *chrono.ParsingContext, m *chrono.RegExpMatch) any {
					return chrono.Components{{Component: "year", Value: 2030}, {Component: "month", Value: 1}, {Component: "day", Value: js.ParseInt(m.Values[1], 0)}}
				},
			})
			return c.Parse("5 sleeps and 7 sleep then 40 sleeps", ref, nil)
		},
		"custom-result-null": func() []*chrono.ParsingResult {
			c := chrono.Casual.Clone()
			c.Parsers = append([]chrono.Parser{&funcParser{
				pattern: func(*chrono.ParsingContext) *chrono.RegExp { return mustRegExp(t, "at (\\d)", "i") },
				extract: func(ctx *chrono.ParsingContext, m *chrono.RegExpMatch) any {
					if m.Values[1] == "0" {
						return nil
					}
					return ctx.CreateParsingResult(m.Index, m.Values[0], ctx.CreateParsingComponents(chrono.Components{
						{Component: "hour", Value: js.ParseInt(m.Values[1], 0)}, {Component: "minute", Value: 0},
					}), nil)
				},
			}}, c.Parsers...)
			return c.Parse("at 0 then at 5 and at 7pm", ref, nil)
		},
		"custom-global": func() []*chrono.ParsingResult {
			c := chrono.NewChrono(&chrono.Configuration{Parsers: []chrono.Parser{&funcParser{
				pattern: func(*chrono.ParsingContext) *chrono.RegExp { return mustRegExp(t, "x(\\d)", "g") },
				extract: func(ctx *chrono.ParsingContext, m *chrono.RegExpMatch) any {
					return ctx.CreateParsingComponents(chrono.Components{{Component: "month", Value: js.ParseInt(m.Values[1], 0)}})
				},
			}}, Refiners: []chrono.Refiner{}})
			return c.Parse("x1 x2 yy x3 x4", ref, nil)
		},
		"custom-refiners": func() []*chrono.ParsingResult {
			c := chrono.Casual.Clone()
			c.Refiners = append([]chrono.Refiner{funcRefiner(func(_ *chrono.ParsingContext, rs []*chrono.ParsingResult) []*chrono.ParsingResult {
				out := []*chrono.ParsingResult{}
				for _, r := range rs {
					if r.Index > 3 {
						out = append(out, r)
					}
				}
				return out
			})}, c.Refiners...)
			c.Refiners = append(c.Refiners, funcRefiner(func(_ *chrono.ParsingContext, rs []*chrono.ParsingResult) []*chrono.ParsingResult {
				for _, r := range rs {
					r.AddTag("custom")
				}
				return rs
			}))
			return c.Parse("Mon tomorrow at 5 and next week", ref, nil)
		},
		"default-chrono": func() []*chrono.ParsingResult {
			return chrono.NewChrono(nil).Parse("next tuesday 10am", ref, nil)
		},
		"components-assign": func() []*chrono.ParsingResult {
			c := chrono.Casual.Clone()
			c.Parsers = []chrono.Parser{&funcParser{
				pattern: func(*chrono.ParsingContext) *chrono.RegExp { return mustRegExp(t, "zz", "") },
				extract: func(ctx *chrono.ParsingContext, _ *chrono.RegExpMatch) any {
					return ctx.CreateParsingComponents(nil).Assign("hour", 3).Imply("minute", 7)
				},
			}}
			return c.Parse("a zz b zz", ref, nil)
		},
		"cache-forward": func() []*chrono.ParsingResult {
			r1 := chrono.Casual.Parse("2 days", ref, &chrono.ParsingOption{ForwardDate: true})
			r2 := chrono.Casual.Parse("2 days", ref, nil)
			return append(r1, r2...)
		},
		"word-boundary-subclass": func() []*chrono.ParsingResult {
			c := chrono.Casual.Clone()
			last := c.Parsers[len(c.Parsers)-1]
			c.Parsers = append([]chrono.Parser{last}, c.Parsers[:len(c.Parsers)-1]...)
			return c.Parse("+2 weeks after next friday", ref, nil)
		},
	}
}

func TestProbeScenarios(t *testing.T) {
	f, err := os.Open("testdata/probes.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	scenarios := probeScenarios(t)
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1<<20), 1<<24)
	n := 0
	for sc.Scan() {
		parts := strings.SplitN(sc.Text(), "\t", 3)
		zone, name, want := parts[0], parts[1], parts[2]
		fn, ok := scenarios[name]
		if !ok {
			t.Fatalf("no Go scenario %q", name)
		}
		inZone(t, zone, func() {
			if got := renderParse(fn()); got != want {
				t.Errorf("%s (TZ=%s):\n got %s\nwant %s", name, zone, got, want)
			}
		})
		n++
	}
	if n != 2*len(scenarios) {
		t.Errorf("ran %d recorded probes for %d scenarios", n, len(scenarios))
	}
}

func TestPrototypeTimezoneThrows(t *testing.T) {
	err := catchJSError(func() {
		chrono.Parse("tomorrow", chrono.ParsingReference{Instant: 0, Timezone: "constructor"}, nil)
	})
	if err == nil || err.Name != "TypeError" || err.Message != "matchedTimezone.dstStart is not a function" {
		t.Fatalf("got %v", err)
	}
}

func TestParserOrderAndIdentity(t *testing.T) {
	if chrono.Casual.Parsers[0] == chrono.Strict.Parsers[0] {
		t.Error("casual and strict share parser objects")
	}
	if got := len(chrono.Casual.Parsers); got != 17 {
		t.Errorf("casual has %d parsers, want 17", got)
	}
	if got := len(chrono.Casual.Refiners); got != 15 {
		t.Errorf("casual has %d refiners, want 15", got)
	}
	if got := len(chrono.Strict.Parsers); got != 12 {
		t.Errorf("strict has %d parsers, want 12", got)
	}
}
