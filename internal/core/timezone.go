// Ported from src/timezone.ts.

package core

import "github.com/wanasit/chrono-go/internal/jsdate"

// TimezoneAbbrMapDefault is TIMEZONE_ABBR_MAP: time zone abbreviations and their offsets in
// minutes. CET, CT, ET, MT and PT are ambiguous (DST-dependent).
var TimezoneAbbrMapDefault = map[string]TimezoneAbbr{
	"ACDT":  {Offset: 630},
	"ACST":  {Offset: 570},
	"ADT":   {Offset: -180},
	"AEDT":  {Offset: 660},
	"AEST":  {Offset: 600},
	"AFT":   {Offset: 270},
	"AKDT":  {Offset: -480},
	"AKST":  {Offset: -540},
	"ALMT":  {Offset: 360},
	"AMST":  {Offset: -180},
	"AMT":   {Offset: -240},
	"ANAST": {Offset: 720},
	"ANAT":  {Offset: 720},
	"AQTT":  {Offset: 300},
	"ART":   {Offset: -180},
	"AST":   {Offset: -240},
	"AWDT":  {Offset: 540},
	"AWST":  {Offset: 480},
	"AZOST": {Offset: 0},
	"AZOT":  {Offset: -60},
	"AZST":  {Offset: 300},
	"AZT":   {Offset: 240},
	"BNT":   {Offset: 480},
	"BOT":   {Offset: -240},
	"BRST":  {Offset: -120},
	"BRT":   {Offset: -180},
	"BST":   {Offset: 60},
	"BTT":   {Offset: 360},
	"CAST":  {Offset: 480},
	"CAT":   {Offset: 120},
	"CCT":   {Offset: 390},
	"CDT":   {Offset: -300},
	"CEST":  {Offset: 120},
	"CET": {Ambiguous: &AmbiguousTimezoneMap{
		TimezoneOffsetDuringDst: 2 * 60,
		TimezoneOffsetNonDst:    60,
		DstStart:                func(year float64) float64 { return GetLastWeekdayOfMonth(year, March, Sunday, 2) },
		DstEnd:                  func(year float64) float64 { return GetLastWeekdayOfMonth(year, October, Sunday, 3) },
	}},
	"CHADT": {Offset: 825},
	"CHAST": {Offset: 765},
	"CKT":   {Offset: -600},
	"CLST":  {Offset: -180},
	"CLT":   {Offset: -240},
	"COT":   {Offset: -300},
	"CST":   {Offset: -360},
	"CT": {Ambiguous: &AmbiguousTimezoneMap{
		TimezoneOffsetDuringDst: -5 * 60,
		TimezoneOffsetNonDst:    -6 * 60,
		DstStart:                func(year float64) float64 { return GetNthWeekdayOfMonth(year, March, Sunday, 2, 2) },
		DstEnd:                  func(year float64) float64 { return GetNthWeekdayOfMonth(year, November, Sunday, 1, 2) },
	}},
	"CVT":   {Offset: -60},
	"CXT":   {Offset: 420},
	"ChST":  {Offset: 600},
	"DAVT":  {Offset: 420},
	"EASST": {Offset: -300},
	"EAST":  {Offset: -360},
	"EAT":   {Offset: 180},
	"ECT":   {Offset: -300},
	"EDT":   {Offset: -240},
	"EEST":  {Offset: 180},
	"EET":   {Offset: 120},
	"EGST":  {Offset: 0},
	"EGT":   {Offset: -60},
	"EST":   {Offset: -300},
	"ET": {Ambiguous: &AmbiguousTimezoneMap{
		TimezoneOffsetDuringDst: -4 * 60,
		TimezoneOffsetNonDst:    -5 * 60,
		DstStart:                func(year float64) float64 { return GetNthWeekdayOfMonth(year, March, Sunday, 2, 2) },
		DstEnd:                  func(year float64) float64 { return GetNthWeekdayOfMonth(year, November, Sunday, 1, 2) },
	}},
	"FJST":  {Offset: 780},
	"FJT":   {Offset: 720},
	"FKST":  {Offset: -180},
	"FKT":   {Offset: -240},
	"FNT":   {Offset: -120},
	"GALT":  {Offset: -360},
	"GAMT":  {Offset: -540},
	"GET":   {Offset: 240},
	"GFT":   {Offset: -180},
	"GILT":  {Offset: 720},
	"GMT":   {Offset: 0},
	"GST":   {Offset: 240},
	"GYT":   {Offset: -240},
	"HAA":   {Offset: -180},
	"HAC":   {Offset: -300},
	"HADT":  {Offset: -540},
	"HAE":   {Offset: -240},
	"HAP":   {Offset: -420},
	"HAR":   {Offset: -360},
	"HAST":  {Offset: -600},
	"HAT":   {Offset: -90},
	"HAY":   {Offset: -480},
	"HKT":   {Offset: 480},
	"HLV":   {Offset: -210},
	"HNA":   {Offset: -240},
	"HNC":   {Offset: -360},
	"HNE":   {Offset: -300},
	"HNP":   {Offset: -480},
	"HNR":   {Offset: -420},
	"HNT":   {Offset: -150},
	"HNY":   {Offset: -540},
	"HOVT":  {Offset: 420},
	"ICT":   {Offset: 420},
	"IDT":   {Offset: 180},
	"IOT":   {Offset: 360},
	"IRDT":  {Offset: 270},
	"IRKST": {Offset: 540},
	"IRKT":  {Offset: 540},
	"IRST":  {Offset: 210},
	"IST":   {Offset: 330},
	"JST":   {Offset: 540},
	"KGT":   {Offset: 360},
	"KRAST": {Offset: 480},
	"KRAT":  {Offset: 480},
	"KST":   {Offset: 540},
	"KUYT":  {Offset: 240},
	"LHDT":  {Offset: 660},
	"LHST":  {Offset: 630},
	"LINT":  {Offset: 840},
	"MAGST": {Offset: 720},
	"MAGT":  {Offset: 720},
	"MART":  {Offset: -510},
	"MAWT":  {Offset: 300},
	"MDT":   {Offset: -360},
	"MESZ":  {Offset: 120},
	"MEZ":   {Offset: 60},
	"MHT":   {Offset: 720},
	"MMT":   {Offset: 390},
	"MSD":   {Offset: 240},
	"MSK":   {Offset: 180},
	"MST":   {Offset: -420},
	"MT": {Ambiguous: &AmbiguousTimezoneMap{
		TimezoneOffsetDuringDst: -6 * 60,
		TimezoneOffsetNonDst:    -7 * 60,
		DstStart:                func(year float64) float64 { return GetNthWeekdayOfMonth(year, March, Sunday, 2, 2) },
		DstEnd:                  func(year float64) float64 { return GetNthWeekdayOfMonth(year, November, Sunday, 1, 2) },
	}},
	"MUT":   {Offset: 240},
	"MVT":   {Offset: 300},
	"MYT":   {Offset: 480},
	"NCT":   {Offset: 660},
	"NDT":   {Offset: -90},
	"NFT":   {Offset: 690},
	"NOVST": {Offset: 420},
	"NOVT":  {Offset: 360},
	"NPT":   {Offset: 345},
	"NST":   {Offset: -150},
	"NUT":   {Offset: -660},
	"NZDT":  {Offset: 780},
	"NZST":  {Offset: 720},
	"OMSST": {Offset: 420},
	"OMST":  {Offset: 420},
	"PDT":   {Offset: -420},
	"PET":   {Offset: -300},
	"PETST": {Offset: 720},
	"PETT":  {Offset: 720},
	"PGT":   {Offset: 600},
	"PHOT":  {Offset: 780},
	"PHT":   {Offset: 480},
	"PKT":   {Offset: 300},
	"PMDT":  {Offset: -120},
	"PMST":  {Offset: -180},
	"PONT":  {Offset: 660},
	"PST":   {Offset: -480},
	"PT": {Ambiguous: &AmbiguousTimezoneMap{
		TimezoneOffsetDuringDst: -7 * 60,
		TimezoneOffsetNonDst:    -8 * 60,
		DstStart:                func(year float64) float64 { return GetNthWeekdayOfMonth(year, March, Sunday, 2, 2) },
		DstEnd:                  func(year float64) float64 { return GetNthWeekdayOfMonth(year, November, Sunday, 1, 2) },
	}},
	"PWT":   {Offset: 540},
	"PYST":  {Offset: -180},
	"PYT":   {Offset: -240},
	"RET":   {Offset: 240},
	"SAMT":  {Offset: 240},
	"SAST":  {Offset: 120},
	"SBT":   {Offset: 660},
	"SCT":   {Offset: 240},
	"SGT":   {Offset: 480},
	"SRT":   {Offset: -180},
	"SST":   {Offset: -660},
	"TAHT":  {Offset: -600},
	"TFT":   {Offset: 300},
	"TJT":   {Offset: 300},
	"TKT":   {Offset: 780},
	"TLT":   {Offset: 540},
	"TMT":   {Offset: 300},
	"TVT":   {Offset: 720},
	"ULAT":  {Offset: 480},
	"UTC":   {Offset: 0},
	"UYST":  {Offset: -120},
	"UYT":   {Offset: -180},
	"UZT":   {Offset: 300},
	"VET":   {Offset: -210},
	"VLAST": {Offset: 660},
	"VLAT":  {Offset: 660},
	"VUT":   {Offset: 660},
	"WAST":  {Offset: 120},
	"WAT":   {Offset: 60},
	"WEST":  {Offset: 60},
	"WESZ":  {Offset: 60},
	"WET":   {Offset: 0},
	"WEZ":   {Offset: 0},
	"WFT":   {Offset: 720},
	"WGST":  {Offset: -120},
	"WGT":   {Offset: -180},
	"WIB":   {Offset: 420},
	"WIT":   {Offset: 540},
	"WITA":  {Offset: 480},
	"WST":   {Offset: 780},
	"WT":    {Offset: 0},
	"YAKST": {Offset: 600},
	"YAKT":  {Offset: 600},
	"YAPT":  {Offset: 600},
	"YEKST": {Offset: 360},
	"YEKT":  {Offset: 360},
}

// GetNthWeekdayOfMonth returns the Date (time value) of the nth occurrence of weekday in the
// month (1-12) of year, at the given hour (local time).
func GetNthWeekdayOfMonth(year float64, month, weekday, n, hour int) float64 {
	dayOfMonth := 0
	i := 0
	for i < n {
		dayOfMonth++
		date := jsdate.NewLocal(year, float64(month-1), float64(dayOfMonth))
		if date.GetDay() == float64(weekday) {
			i++
		}
	}
	return jsdate.NewLocal(year, float64(month-1), float64(dayOfMonth), float64(hour)).Time()
}

// GetLastWeekdayOfMonth returns the Date (time value) of the last occurrence of weekday in the
// month (1-12) of year, at the given hour (local time).
func GetLastWeekdayOfMonth(year float64, month, weekday, hour int) float64 {
	// Find the first weekday of the next month, compare with the given weekday, and use the
	// difference to determine how many days to subtract from the first of the next month.
	oneIndexedWeekday := float64(weekday)
	if weekday == 0 {
		oneIndexedWeekday = 7
	}
	date := jsdate.NewLocal(year, float64(month-1+1), 1, 12)
	firstWeekdayNextMonth := date.GetDay()
	if date.GetDay() == 0 {
		firstWeekdayNextMonth = 7
	}
	var dayDiff float64
	if firstWeekdayNextMonth == oneIndexedWeekday {
		dayDiff = 7
	} else if firstWeekdayNextMonth < oneIndexedWeekday {
		dayDiff = 7 + firstWeekdayNextMonth - oneIndexedWeekday
	} else {
		dayDiff = firstWeekdayNextMonth - oneIndexedWeekday
	}
	date.SetDate(date.GetDate() - dayDiff)
	return jsdate.NewLocal(year, float64(month-1), date.GetDate(), float64(hour)).Time()
}

// ToTimezoneOffset finds the offset (minutes) of timezoneInput: nil (JS null/undefined), a number
// (float64 or int, returned as is) or a string looked up in timezoneOverrides, then in
// TimezoneAbbrMapDefault. Ambiguous abbreviations use date (nil for JS null) to pick the DST or
// non-DST offset. ok is false when the result is null.
func ToTimezoneOffset(timezoneInput any, date *jsdate.Date, timezoneOverrides map[string]TimezoneAbbr) (offset float64, ok bool) {
	var name string
	switch v := timezoneInput.(type) {
	case nil:
		return 0, false
	case float64:
		return v, true
	case int:
		return float64(v), true
	case string:
		name = v
	default:
		return 0, false
	}

	matchedTimezone, found := timezoneOverrides[name]
	if !found {
		matchedTimezone, found = TimezoneAbbrMapDefault[name]
	}
	if !found {
		if !IsObjectPrototypeKey(name) {
			return 0, false
		}
		// The lookup yields an Object.prototype member: not a number, and it has no dstStart.
		if date == nil {
			return 0, false
		}
		throwTypeError("matchedTimezone.dstStart is not a function")
	}
	// This means that we have matched an unambiguous timezone
	if matchedTimezone.Ambiguous == nil {
		return matchedTimezone.Offset, true
	}

	// Without refDate as context, there's no way to know if DST or non-DST offset should be used.
	if date == nil {
		return 0, false
	}
	amb := matchedTimezone.Ambiguous
	// Return DST offset if the refDate is during daylight savings
	if date.Time() > amb.DstStart(date.GetFullYear()) && !(date.Time() > amb.DstEnd(date.GetFullYear())) {
		return amb.TimezoneOffsetDuringDst, true
	}
	// refDate is not during DST => return non-DST offset
	return amb.TimezoneOffsetNonDst, true
}
