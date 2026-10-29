package nl_test

import (
	"encoding/json"
	"os"
	"testing"

	chrono "github.com/wanasit/chrono-go"
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/locales/nl"
)

// testdata/patterns.json holds [source, flags] of parser.pattern(context) for every parser of the
// original configurations, with the options {} and {forwardDate: true}.
func TestParserPatterns(t *testing.T) {
	data, err := os.ReadFile("testdata/patterns.json")
	if err != nil {
		t.Fatal(err)
	}
	var want map[string][][][2]string
	if err := json.Unmarshal(data, &want); err != nil {
		t.Fatal(err)
	}
	configs := map[string]*chrono.Configuration{
		"casual":   nl.CreateCasualConfiguration(true),
		"strict":   nl.CreateConfiguration(true, true),
		"casual0":  nl.CreateCasualConfiguration(false),
		"config00": nl.CreateConfiguration(false, false),
	}
	options := []*core.ParsingOption{{}, {ForwardDate: true}}
	for name, cfg := range configs {
		for oi, opt := range options {
			ctx := core.NewParsingContext("x", core.ParsingReference{Instant: 1344600000000}, opt)
			w := want[name][oi]
			if len(w) != len(cfg.Parsers) {
				t.Fatalf("%s: %d parsers, want %d", name, len(cfg.Parsers), len(w))
			}
			for i, p := range cfg.Parsers {
				re := p.Pattern(ctx)
				if re.Source() != w[i][0] || re.Flags() != w[i][1] {
					t.Errorf("%s option %d parser %d (%T):\n got /%s/%s\nwant /%s/%s", name, oi, i, p, re.Source(), re.Flags(), w[i][0], w[i][1])
				}
			}
		}
	}
}
