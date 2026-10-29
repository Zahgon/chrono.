// Ported from src/common/parsers/ISOFormatParser.ts.

package core

import "github.com/wanasit/chrono-go/internal/js"

// ISO 8601: YYYY-MM-DD, YYYY-MM-DDThh:mmTZD, YYYY-MM-DDThh:mm:ssTZD, YYYY-MM-DDThh:mm:ss.sTZD
// (TZD = Z or +hh:mm or -hh:mm).
var isoPattern = MustRegExp(
	"([0-9]{4})\\-([0-9]{1,2})\\-([0-9]{1,2})"+
		"(?:T"+
		"([0-9]{1,2}):([0-9]{1,2})"+
		"(?:"+
		":([0-9]{1,2})(?:\\.(\\d{1,4}))?"+
		")?"+
		"("+
		"Z|([+-]\\d{2}):?(\\d{2})?"+
		")?"+
		")?"+
		"(?=\\W|$)",
	"i",
)

const (
	isoYearNumberGroup        = 1
	isoMonthNumberGroup       = 2
	isoDateNumberGroup        = 3
	isoHourNumberGroup        = 4
	isoMinuteNumberGroup      = 5
	isoSecondNumberGroup      = 6
	isoMillisecondNumberGroup = 7
	isoTzdGroup               = 8
	isoTzdHourOffsetGroup     = 9
	isoTzdMinuteOffsetGroup   = 10
)

// ISOFormatParser parses ISO 8601 dates and date-times.
type ISOFormatParser struct {
	AbstractParserWithWordBoundaryChecking
}

// NewISOFormatParser returns a new ISOFormatParser.
func NewISOFormatParser() *ISOFormatParser {
	p := &ISOFormatParser{}
	p.Self = p
	return p
}

// InnerPattern returns the ISO pattern.
func (p *ISOFormatParser) InnerPattern(context *ParsingContext) *RegExp { return isoPattern }

func parseIntGroup(match *RegExpMatch, i int) float64 {
	return js.ParseInt(match.Get(i), 0)
}

// InnerExtract builds the components of the ISO date.
func (p *ISOFormatParser) InnerExtract(context *ParsingContext, match *RegExpMatch) any {
	components := context.CreateParsingComponents(Components{
		{ComponentYear, parseIntGroup(match, isoYearNumberGroup)},
		{ComponentMonth, parseIntGroup(match, isoMonthNumberGroup)},
		{ComponentDay, parseIntGroup(match, isoDateNumberGroup)},
	})
	if match.Has(isoHourNumberGroup) {
		components.Assign(ComponentHour, parseIntGroup(match, isoHourNumberGroup))
		components.Assign(ComponentMinute, parseIntGroup(match, isoMinuteNumberGroup))

		if match.Has(isoSecondNumberGroup) {
			components.Assign(ComponentSecond, parseIntGroup(match, isoSecondNumberGroup))
		}

		if match.Has(isoMillisecondNumberGroup) {
			components.Assign(ComponentMillisecond, parseIntGroup(match, isoMillisecondNumberGroup))
		}
		if match.Has(isoTzdGroup) {
			// The Zulu time zone (Z) is equivalent to UTC
			offset := 0.0
			if match.Truthy(isoTzdHourOffsetGroup) {
				hourOffset := parseIntGroup(match, isoTzdHourOffsetGroup)
				minuteOffset := 0.0
				if match.Has(isoTzdMinuteOffsetGroup) {
					minuteOffset = parseIntGroup(match, isoTzdMinuteOffsetGroup)
				}
				offset = hourOffset * 60
				if offset < 0 {
					offset -= minuteOffset
				} else {
					offset += minuteOffset
				}
			}
			components.Assign(ComponentTimezoneOffset, offset)
		}
	}
	return components.AddTag("parser/ISOFormatParser")
}
