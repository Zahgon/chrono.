// Ported from src/locales/it/parsers/ITRelativeDateFormatParser.ts.

package it

import (
	"github.com/wanasit/chrono-go/internal/core"
	"github.com/wanasit/chrono-go/internal/js"
	"github.com/wanasit/chrono-go/internal/jsdate"
)

// Italian uses "il mese prossimo" (article + noun + modifier) more commonly than "il prossimo
// mese".
var relativeDateFormatPattern = core.MustRegExp(
	"(?:"+
		// Pattern 1: "il prossimo mese" / "la prossima settimana" (modifier before noun)
		"(?:il\\s*|la\\s*|l'\\s*)?(questo|questa|quest'|scorso|scorsa|prossimo|prossima|dopo\\s*questo|dopo\\s*questa)\\s*("+core.MatchAnyPattern(TimeUnitDictionary.Keys())+")"+
		"|"+
		// Pattern 2: "il mese prossimo" / "la settimana scorsa" (modifier after noun)
		"(?:il\\s*|la\\s*|l'\\s*)("+core.MatchAnyPattern(TimeUnitDictionary.Keys())+")\\s*(scorso|scorsa|prossimo|prossima)"+
		")(?=\\W|$)",
	"i",
)

const (
	// Pattern 1 groups
	relativeModifierWordGroup1 = 1
	relativeRelativeWordGroup1 = 2
	// Pattern 2 groups
	relativeRelativeWordGroup2 = 3
	relativeModifierWordGroup2 = 4
)

var (
	reSettimanaI = core.MustRegExp("settimana", "i")
	reMeseI      = core.MustRegExp("mese", "i")
	reAnnoI      = core.MustRegExp("anno", "i")
)

// RelativeDateFormatParser is ITRelativeDateFormatParser ("il mese prossimo", "la settimana
// scorsa", "quest'anno").
type RelativeDateFormatParser struct {
	core.AbstractParserWithWordBoundaryChecking
}

// NewRelativeDateFormatParser returns a new ITRelativeDateFormatParser.
func NewRelativeDateFormatParser() *RelativeDateFormatParser {
	p := &RelativeDateFormatParser{}
	p.Self = p
	return p
}

// InnerPattern returns the relative date pattern.
func (p *RelativeDateFormatParser) InnerPattern(context *core.ParsingContext) *core.RegExp {
	return relativeDateFormatPattern
}

// InnerExtract returns the components of the relative date.
func (p *RelativeDateFormatParser) InnerExtract(context *core.ParsingContext, match *core.RegExpMatch) any {
	// Determine which pattern matched
	var modifier, unitWord string

	if match.Truthy(relativeModifierWordGroup1) {
		// Pattern 1: "il prossimo mese"
		modifier = js.ToLowerCase(match.Get(relativeModifierWordGroup1))
		unitWord = js.ToLowerCase(match.Get(relativeRelativeWordGroup1))
	} else {
		// Pattern 2: "il mese prossimo"
		modifier = js.ToLowerCase(match.Get(relativeModifierWordGroup2))
		unitWord = js.ToLowerCase(match.Get(relativeRelativeWordGroup2))
	}

	timeunit := timeunitLookup(unitWord)

	if modifier == "prossimo" || modifier == "prossima" || js.StartsWith(modifier, "dopo", 0) {
		timeUnits := core.NewDuration()
		timeUnits.Set(timeunit, 1)
		return core.CreateRelativeFromReference(context.Reference, timeUnits)
	}

	if modifier == "scorso" || modifier == "scorsa" {
		timeUnits := core.NewDuration()
		timeUnits.Set(timeunit, -1)
		return core.CreateRelativeFromReference(context.Reference, timeUnits)
	}

	components := context.CreateParsingComponents(nil)
	date := jsdate.FromTime(context.Reference.Instant)

	if reSettimanaI.Exec(unitWord) != nil {
		// Questa settimana
		date.SetDate(date.GetDate() - date.GetDay())
		components.Imply(core.ComponentDay, date.GetDate())
		components.Imply(core.ComponentMonth, date.GetMonth()+1)
		components.Imply(core.ComponentYear, date.GetFullYear())
	} else if reMeseI.Exec(unitWord) != nil {
		// Questo mese
		date.SetDate(1)
		components.Imply(core.ComponentDay, date.GetDate())
		components.Assign(core.ComponentYear, date.GetFullYear())
		components.Assign(core.ComponentMonth, date.GetMonth()+1)
	} else if reAnnoI.Exec(unitWord) != nil {
		// Quest'anno
		date.SetDate(1)
		date.SetMonth(0)
		components.Imply(core.ComponentDay, date.GetDate())
		components.Imply(core.ComponentMonth, date.GetMonth()+1)
		components.Assign(core.ComponentYear, date.GetFullYear())
	}

	return components
}
