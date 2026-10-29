// Ported from src/locales/ja/parsers/JPCasualDateParser.ts.

package ja

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/jsdate"
)

var casualDatePattern = core.MustRegExp("今日|きょう|本日|ほんじつ|昨日|きのう|明日|あした|今夜|こんや|今夕|こんゆう|今晩|こんばん|今朝|けさ", "i")

func normalizeTextToKanji(text string) string {
	switch text {
	case "きょう":
		return "今日"
	case "ほんじつ":
		return "本日"
	case "きのう":
		return "昨日"
	case "あした":
		return "明日"
	case "こんや":
		return "今夜"
	case "こんゆう":
		return "今夕"
	case "こんばん":
		return "今晩"
	case "けさ":
		return "今朝"
	default:
		return text
	}
}

// CasualDateParser is JPCasualDateParser ("今日", "昨日", "明日", "今夜", "今朝", ...).
type CasualDateParser struct {
	_ byte // non-zero size: each constructor call returns a distinct object, as `new` does
}

// NewCasualDateParser returns a new JPCasualDateParser.
func NewCasualDateParser() *CasualDateParser { return &CasualDateParser{} }

// Pattern returns the casual date pattern.
func (p *CasualDateParser) Pattern(context *core.ParsingContext) *core.RegExp {
	return casualDatePattern
}

// Extract returns the components of the casual date.
func (p *CasualDateParser) Extract(context *core.ParsingContext, match *core.RegExpMatch) any {
	text := normalizeTextToKanji(match.Get(0))

	components := context.CreateParsingComponents(nil)

	switch text {
	case "昨日":
		return core.Yesterday(context.Reference)

	case "明日":
		return core.Tomorrow(context.Reference)

	case "本日", "今日":
		return core.Today(context.Reference)
	}

	if text == "今夜" || text == "今夕" || text == "今晩" {
		components.Imply(core.ComponentHour, 22)
		components.Assign(core.ComponentMeridiem, core.PM)
	} else if core.StringPattern("今朝").Exec(text) != nil {
		components.Imply(core.ComponentHour, 6)
		components.Assign(core.ComponentMeridiem, core.AM)
	}

	date := jsdate.FromTime(context.RefDate)
	components.Assign(core.ComponentDay, date.GetDate())
	components.Assign(core.ComponentMonth, date.GetMonth()+1)
	components.Assign(core.ComponentYear, date.GetFullYear())
	return components
}
