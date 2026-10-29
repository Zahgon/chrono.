# chrono-go

A Go port of [chrono](https://github.com/wanasit/chrono) v2.10.1, the natural language date parser
by Wanasit Tanakitrungruang (MIT license, see `LICENSE.txt`). It reproduces the original's results
exactly — the same results, texts, indices, components, certainty, tags and dates — for every
input: casual and strict English, UK-style English, and the other locales under `locales/`.

## Usage

```go
import (
	"fmt"
	"time"

	chrono "github.com/wanasit/chrono-go"
	"github.com/wanasit/chrono-go/locales/en"
)

ref := chrono.ParsingReference{Instant: float64(time.Now().UnixMilli())}

results := chrono.Parse("An appointment on Sep 12-13 at 5pm", ref, nil)
for _, r := range results {
	fmt.Println(r.Index, r.Text, r.Start.Date(), r.End != nil)
}

// The start date of the first result, as a time value (milliseconds since the epoch).
if ms, ok := chrono.ParseDate("next Friday", ref, &chrono.ParsingOption{ForwardDate: true}); ok {
	fmt.Println(time.UnixMilli(int64(ms)))
}

// Strict mode, UK-style (day-first) English, a reference time zone (an abbreviation or minutes).
chrono.Strict.Parse("2 weeks ago", ref, nil)
en.GB.Parse("12/7/2014", chrono.ParsingReference{Instant: ref.Instant, Timezone: "BST"}, nil)
```

Like the JavaScript original, this package works on JavaScript values: dates are time values
(`float64` milliseconds, `NaN` for an invalid date) computed in the process's local time zone (the
`TZ` environment variable), result indices and text lengths count UTF-16 code units, and a missing
component (`null`) is reported by `Get` as `ok == false`.

A `Chrono` is its list of parsers and refiners (`Parsers`, `Refiners`), applied in order; clone one
with `Clone`, reorder or remove entries, or add your own `Parser` / `Refiner` implementations:

```go
type sleepsParser struct{}

func (sleepsParser) Pattern(*chrono.ParsingContext) *chrono.RegExp {
	re, _ := chrono.NewRegExp(`(\d+)\s*sleeps`, "i")
	return re
}

func (sleepsParser) Extract(ctx *chrono.ParsingContext, m *chrono.RegExpMatch) any {
	return chrono.Components{{Component: "year", Value: 2030}, {Component: "month", Value: 1}}
}

custom := chrono.Casual.Clone()
custom.Parsers = append(custom.Parsers, sleepsParser{})
```

A `Chrono` (like its parsers, which cache their patterns) is not safe for concurrent use; use one
per goroutine (`Clone` shares the parser objects, so create fresh ones with the locale packages'
`CreateCasualConfiguration` / `CreateConfiguration` and `NewChrono`).

Where the original throws (for example a reference time zone named like an `Object.prototype`
member, such as `"constructor"`), the port panics with a `*chrono.Error` carrying the JavaScript
error's name and message.

## Attribution

chrono is Copyright (c) 2014, Wanasit Tanakitrungruang, and released under the MIT license; this
port keeps that notice in `LICENSE.txt`. Every Go file names the TypeScript file it was ported from.
