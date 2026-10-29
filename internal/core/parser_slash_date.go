// Ported from src/common/parsers/SlashDateFormatParser.ts.

package core

import (
	"strings"

	"github.com/wanasit/chrono-go/internal/js"
)

// Date format with slash "/" (or dot ".") between numbers: 7/10, 7/12/2020, 7.12.2020.
var slashDatePattern = MustRegExp(
	"([^\\d]|^)"+
		"([0-3]{0,1}[0-9]{1})[\\/\\.\\-]([0-3]{0,1}[0-9]{1})"+
		"(?:[\\/\\.\\-]([0-9]{4}|[0-9]{2}))?"+
		"(\\W|$)",
	"i",
)

const (
	slashOpeningGroup       = 1
	slashEndingGroup        = 5
	slashFirstNumbersGroup  = 2
	slashSecondNumbersGroup = 3
	slashYearGroup          = 4
)

var (
	reSlashTextBefore   = StringPattern("\\d/?$")
	reSlashTextAfter    = StringPattern("^/?\\d")
	reVersionTwoParts   = MustRegExp("^\\d\\.\\d$", "")
	reVersionThreeParts = MustRegExp("^\\d\\.\\d{1,2}\\.\\d{1,2}\\s*$", "")
)

// SlashDateFormatParser parses numeric dates separated by "/", "." or "-".
type SlashDateFormatParser struct {
	GroupNumberMonth int
	GroupNumberDay   int
}

// NewSlashDateFormatParser returns a parser reading day-first dates when littleEndian.
func NewSlashDateFormatParser(littleEndian bool) *SlashDateFormatParser {
	p := &SlashDateFormatParser{GroupNumberMonth: slashFirstNumbersGroup, GroupNumberDay: slashSecondNumbersGroup}
	if littleEndian {
		p.GroupNumberMonth = slashSecondNumbersGroup
		p.GroupNumberDay = slashFirstNumbersGroup
	}
	return p
}

// Pattern returns the slash date pattern.
func (p *SlashDateFormatParser) Pattern(context *ParsingContext) *RegExp { return slashDatePattern }

// Extract builds the result of a slash date.
func (p *SlashDateFormatParser) Extract(context *ParsingContext, match *RegExpMatch) any {
	// Because of how pattern is executed on remaining text in `chrono.ts`, the character before
	// the match could still be a number (e.g. X[X/YY/ZZ] or XX[/YY/ZZ] or [XX/YY/]ZZ).
	index := match.Index + js.Len(match.Get(slashOpeningGroup))
	indexEnd := match.Index + js.Len(match.Get(0)) - js.Len(match.Get(slashEndingGroup))
	if index > 0 {
		textBefore := context.Substring(0, index)
		if reSlashTextBefore.Exec(textBefore) != nil {
			return nil
		}
	}
	if indexEnd < context.TextLength() {
		textAfter := context.SubstringFrom(indexEnd)
		if reSlashTextAfter.Exec(textAfter) != nil {
			return nil
		}
	}

	text := context.Substring(index, indexEnd)

	// '1.12', '1.12.12' is more like a version numbers
	if reVersionTwoParts.Exec(text) != nil || reVersionThreeParts.Exec(text) != nil {
		return nil
	}

	// MM/dd -> OK
	// MM.dd -> NG
	if !match.Truthy(slashYearGroup) && !strings.Contains(text, "/") {
		return nil
	}

	result := context.CreateParsingResult(index, text, nil, nil)
	month := js.ParseInt(match.Get(p.GroupNumberMonth), 0)
	day := js.ParseInt(match.Get(p.GroupNumberDay), 0)
	if month < 1 || month > 12 {
		if month > 12 {
			if day >= 1 && day <= 12 && month <= 31 {
				day, month = month, day
			} else {
				return nil
			}
		}
	}

	if day < 1 || day > 31 {
		return nil
	}

	result.Start.Assign(ComponentDay, day)
	result.Start.Assign(ComponentMonth, month)

	if match.Truthy(slashYearGroup) {
		rawYearNumber := js.ParseInt(match.Get(slashYearGroup), 0)
		year := FindMostLikelyADYear(rawYearNumber)
		result.Start.Assign(ComponentYear, year)
	} else {
		year := FindYearClosestToRef(context.RefDate, day, month)
		result.Start.Imply(ComponentYear, year)
	}

	return result.AddTag("parser/SlashDateFormatParser")
}
