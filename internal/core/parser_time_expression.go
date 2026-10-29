// Ported from src/common/parsers/AbstractTimeExpressionParser.ts.

package core

import (
	"math"

	"github.com/wanasit/chrono-go/internal/js"
)

func primaryTimePattern(leftBoundary, primaryPrefix, primarySuffix, flags string) *RegExp {
	return MustRegExp(
		leftBoundary+
			primaryPrefix+
			"(\\d{1,4})"+
			"(?:"+
			"(?:\\.|:|：)"+
			"(\\d{1,2})"+
			"(?:"+
			"(?::|：)"+
			"(\\d{2})"+
			"(?:\\.(\\d{1,6}))?"+
			")?"+
			")?"+
			"(?:\\s*(a\\.m\\.|p\\.m\\.|am?|pm?))?"+
			primarySuffix,
		flags,
	)
}

func followingTimePatten(followingPhase, followingSuffix string) *RegExp {
	return MustRegExp(
		"^("+followingPhase+")"+
			"(\\d{1,4})"+
			"(?:"+
			"(?:\\.|\\:|\\：)"+
			"(\\d{1,2})"+
			"(?:"+
			"(?:\\.|\\:|\\：)"+
			"(\\d{1,2})(?:\\.(\\d{1,6}))?"+
			")?"+
			")?"+
			"(?:\\s*(a\\.m\\.|p\\.m\\.|am?|pm?))?"+
			followingSuffix,
		"i",
	)
}

const (
	hourGroup        = 2
	minuteGroup      = 3
	secondGroup      = 4
	milliSecondGroup = 5
	amPmHourGroup    = 6
)

// TimeExpressionParserImpl holds the abstract and overridable methods of
// AbstractTimeExpressionParser. PrimaryPrefix and FollowingPhase are abstract; the others have
// defaults on the abstract struct (promoted through embedding) and can be overridden by defining
// them on the embedding type (call p.AbstractTimeExpressionParser.X(...) for `super.x(...)`).
type TimeExpressionParserImpl interface {
	PrimaryPrefix() string
	FollowingPhase() string
	PatternFlags() string
	PrimaryPatternLeftBoundary() string
	PrimarySuffix() string
	FollowingSuffix() string
	ExtractPrimaryTimeComponents(context *ParsingContext, match *RegExpMatch) *ParsingComponents
	ExtractFollowingTimeComponents(context *ParsingContext, match *RegExpMatch, result *ParsingResult) *ParsingComponents
}

// AbstractTimeExpressionParser parses time expressions ("10:30", "5pm", "10:00 - 11:00").
// Embed it, implement PrimaryPrefix and FollowingPhase, and set Self (and StrictMode) in the
// constructor.
type AbstractTimeExpressionParser struct {
	Self       TimeExpressionParserImpl
	StrictMode bool

	cachedPrimaryPrefix      *string
	cachedPrimarySuffix      *string
	cachedPrimaryTimePattern *RegExp

	cachedFollowingPhase      *string
	cachedFollowingSuffix     *string
	cachedFollowingTimePatten *RegExp
}

// PatternFlags is the default "i".
func (p *AbstractTimeExpressionParser) PatternFlags() string { return "i" }

// PrimaryPatternLeftBoundary is the default `(^|\s|T|\b)`.
func (p *AbstractTimeExpressionParser) PrimaryPatternLeftBoundary() string {
	return "(^|\\s|T|\\b)"
}

// PrimarySuffix is the default `(?!/)(?=\W|$)`.
func (p *AbstractTimeExpressionParser) PrimarySuffix() string { return "(?!/)(?=\\W|$)" }

// FollowingSuffix is the default `(?!/)(?=\W|$)`.
func (p *AbstractTimeExpressionParser) FollowingSuffix() string { return "(?!/)(?=\\W|$)" }

// Pattern returns the primary time pattern (cached by prefix and suffix).
func (p *AbstractTimeExpressionParser) Pattern(context *ParsingContext) *RegExp {
	return p.GetPrimaryTimePatternThroughCache()
}

var (
	reFourDigitsStart       = MustRegExp("^\\d{4}", "")
	reThreeFourDigitsStart  = MustRegExp("^\\d{3,4}", "")
	reSignedTwoToFourDigits = MustRegExp("^\\s*([+-])\\s*\\d{2,4}$", "")
	reSignedHHMM            = MustRegExp("^\\s*([+-])\\s*\\d{2}\\W\\d{2}", "")
	reSignedThreeFourDigits = MustRegExp("^\\s*([+-])\\s*\\d{3,4}$", "")
	reSingleDigit           = MustRegExp("^\\d$", "")
	reThreeOrMoreDigits     = MustRegExp("^\\d\\d\\d+$", "")
	reDigitThenAP           = MustRegExp("\\d[apAP]$", "")
	reEndingWithNumbers     = MustRegExp("[^\\d:.](\\d[\\d.]+)$", "")
	reDotTwoDigits          = MustRegExp("\\d(\\.\\d{2})+$", "")
	reNumberDashNumber      = MustRegExp("^\\d+-\\d+$", "")
	reEndingWithRange       = MustRegExp("[^\\d:.](\\d[\\d.]+)\\s*-\\s*(\\d[\\d.]+)$", "")
)

// Extract parses the primary time and an optional following ("- 11:00") time.
func (p *AbstractTimeExpressionParser) Extract(context *ParsingContext, match *RegExpMatch) any {
	startComponents := p.Self.ExtractPrimaryTimeComponents(context, match)
	if startComponents == nil {
		// If the match seem like a year e.g. "2013.12:...",
		// then skips the year part and try matching again.
		if reFourDigitsStart.Exec(match.Get(0)) != nil {
			match.Index += 4 // Skip over potential overlapping pattern
			return nil
		}

		match.Index += js.Len(match.Get(0)) // Skip over potential overlapping pattern
		return nil
	}

	index := match.Index + js.Len(match.Get(1))
	text := js.SubstringFrom(match.Get(0), js.Len(match.Get(1)))
	result := context.CreateParsingResult(index, text, startComponents, nil)
	match.Index += js.Len(match.Get(0)) // Skip over potential overlapping pattern

	remainingText := context.SubstringFrom(match.Index)
	followingPattern := p.GetFollowingTimePatternThroughCache()
	followingMatch := followingPattern.Exec(remainingText)

	// Pattern "456-12", "2022-12" should not be time without proper context
	if reThreeFourDigitsStart.Exec(text) != nil && followingMatch != nil {
		// e.g. "2022-12"
		if reSignedTwoToFourDigits.Exec(followingMatch.Get(0)) != nil {
			return nil
		}
		// e.g. "2022-12:01..."
		if reSignedHHMM.Exec(followingMatch.Get(0)) != nil {
			return nil
		}
	}

	if followingMatch == nil ||
		// Pattern "YY.YY -XXXX" is more like timezone offset
		reSignedThreeFourDigits.Exec(followingMatch.Get(0)) != nil {
		return nilIfNoResult(p.checkAndReturnWithoutFollowingPattern(result))
	}

	result.End = p.Self.ExtractFollowingTimeComponents(context, followingMatch, result)
	if result.End != nil {
		result.Text += followingMatch.Get(0)
	}

	return nilIfNoResult(p.checkAndReturnWithFollowingPattern(result))
}

func nilIfNoResult(r *ParsingResult) any {
	if r == nil {
		return nil
	}
	return r
}

// ExtractPrimaryTimeComponents extracts the primary time from the match (nil when it is not a
// time).
func (p *AbstractTimeExpressionParser) ExtractPrimaryTimeComponents(context *ParsingContext, match *RegExpMatch) *ParsingComponents {
	components := context.CreateParsingComponents(nil)
	minute := 0.0
	meridiem := -1.0 // null

	// ----- Hours
	hour := js.ParseInt(match.Get(hourGroup), 0)
	if !match.Has(hourGroup) {
		hour = math.NaN()
	}
	if hour > 100 {
		// When time is like '2019', it is more likely a year.
		if js.Len(match.Get(hourGroup)) == 4 && !match.Has(minuteGroup) && !match.Truthy(amPmHourGroup) {
			return nil
		}

		if p.StrictMode || match.Has(minuteGroup) {
			return nil
		}

		minute = math.Mod(hour, 100)
		hour = math.Floor(hour / 100)
	}

	if hour > 24 {
		return nil
	}

	// ----- Minutes
	if match.Has(minuteGroup) {
		if js.Len(match.Get(minuteGroup)) == 1 && !match.Truthy(amPmHourGroup) {
			// Skip single digit minute e.g. "at 1.1 xx"
			return nil
		}

		minute = js.ParseInt(match.Get(minuteGroup), 0)
	}

	if minute >= 60 {
		return nil
	}

	if hour > 12 {
		meridiem = PM
	}

	// ----- AM & PM
	if match.Has(amPmHourGroup) {
		if hour > 12 {
			return nil
		}
		ampm := js.ToLowerCase(FirstChar(match.Get(amPmHourGroup)))
		if ampm == "a" {
			meridiem = AM
			if hour == 12 {
				hour = 0
			}
		}

		if ampm == "p" {
			meridiem = PM
			if hour != 12 {
				hour += 12
			}
		}
	}

	components.Assign(ComponentHour, hour)
	components.Assign(ComponentMinute, minute)

	if meridiem != -1 {
		components.Assign(ComponentMeridiem, meridiem)
	} else {
		if hour < 12 {
			components.Imply(ComponentMeridiem, AM)
		} else {
			components.Imply(ComponentMeridiem, PM)
		}
	}

	// ----- Millisecond
	if match.Has(milliSecondGroup) {
		millisecond := js.ParseInt(js.Substring(match.Get(milliSecondGroup), 0, 3), 0)
		if millisecond >= 1000 {
			return nil
		}
		components.Assign(ComponentMillisecond, millisecond)
	}

	// ----- Second
	if match.Has(secondGroup) {
		second := js.ParseInt(match.Get(secondGroup), 0)
		if second >= 60 {
			return nil
		}
		components.Assign(ComponentSecond, second)
	}

	return components
}

// ExtractFollowingTimeComponents extracts the end time of a range from the following match.
func (p *AbstractTimeExpressionParser) ExtractFollowingTimeComponents(context *ParsingContext, match *RegExpMatch, result *ParsingResult) *ParsingComponents {
	components := context.CreateParsingComponents(nil)

	// ----- Millisecond
	if match.Has(milliSecondGroup) {
		millisecond := js.ParseInt(js.Substring(match.Get(milliSecondGroup), 0, 3), 0)
		if millisecond >= 1000 {
			return nil
		}
		components.Assign(ComponentMillisecond, millisecond)
	}

	// ----- Second
	if match.Has(secondGroup) {
		second := js.ParseInt(match.Get(secondGroup), 0)
		if second >= 60 {
			return nil
		}
		components.Assign(ComponentSecond, second)
	}

	hour := js.ParseInt(match.Get(hourGroup), 0)
	if !match.Has(hourGroup) {
		hour = math.NaN()
	}
	minute := 0.0
	meridiem := -1.0

	// ----- Minute
	if match.Has(minuteGroup) {
		minute = js.ParseInt(match.Get(minuteGroup), 0)
	} else if hour > 100 {
		minute = math.Mod(hour, 100)
		hour = math.Floor(hour / 100)
	}

	if minute >= 60 || hour > 24 {
		return nil
	}

	if hour >= 12 {
		meridiem = PM
	}

	// ----- AM & PM
	if match.Has(amPmHourGroup) {
		if hour > 12 {
			return nil
		}

		ampm := js.ToLowerCase(FirstChar(match.Get(amPmHourGroup)))
		if ampm == "a" {
			meridiem = AM
			if hour == 12 {
				hour = 0
				if !components.IsCertain(ComponentDay) {
					components.Imply(ComponentDay, components.V(ComponentDay)+1)
				}
			}
		}

		if ampm == "p" {
			meridiem = PM
			if hour != 12 {
				hour += 12
			}
		}

		if !result.Start.IsCertain(ComponentMeridiem) {
			if meridiem == AM {
				result.Start.Imply(ComponentMeridiem, AM)

				if startHour, ok := result.Start.Get(ComponentHour); ok && startHour == 12 {
					result.Start.Assign(ComponentHour, 0)
				}
			} else {
				result.Start.Imply(ComponentMeridiem, PM)

				if startHour, ok := result.Start.Get(ComponentHour); !ok || startHour != 12 {
					result.Start.Assign(ComponentHour, result.Start.V(ComponentHour)+12)
				}
			}
		}
	}

	components.Assign(ComponentHour, hour)
	components.Assign(ComponentMinute, minute)

	if meridiem >= 0 {
		components.Assign(ComponentMeridiem, meridiem)
	} else {
		startAtPM := result.Start.IsCertain(ComponentMeridiem) && result.Start.V(ComponentHour) > 12
		if startAtPM {
			if result.Start.V(ComponentHour)-12 > hour {
				// 10pm - 1 (am)
				components.Imply(ComponentMeridiem, AM)
			} else if hour <= 12 {
				components.Assign(ComponentHour, hour+12)
				components.Assign(ComponentMeridiem, PM)
			}
		} else if hour > 12 {
			components.Imply(ComponentMeridiem, PM)
		} else if hour <= 12 {
			components.Imply(ComponentMeridiem, AM)
		}
	}

	if components.Date() < result.Start.Date() {
		components.Imply(ComponentDay, components.V(ComponentDay)+1)
	}

	return components
}

func (p *AbstractTimeExpressionParser) checkAndReturnWithoutFollowingPattern(result *ParsingResult) *ParsingResult {
	// Single digit (e.g "1") should not be counted as time expression (without proper context)
	if reSingleDigit.Exec(result.Text) != nil {
		return nil
	}

	// Three or more digit (e.g. "203", "2014") should not be counted as time expression
	if reThreeOrMoreDigits.Exec(result.Text) != nil {
		return nil
	}

	// Instead of "am/pm", it ends with "a" or "p" (e.g "1a", "123p"), this seems unlikely
	if reDigitThenAP.Exec(result.Text) != nil {
		return nil
	}

	// If it ends only with numbers or dots
	if endingWithNumbers := reEndingWithNumbers.Exec(result.Text); endingWithNumbers != nil {
		endingNumbers := endingWithNumbers.Get(1)

		// In strict mode (e.g. "at 1" or "at 1.2"), this should not be accepted
		if p.StrictMode {
			return nil
		}

		// If it ends only with dot single digit, e.g. "at 1.2"
		if js.Includes(endingNumbers, ".") && reDotTwoDigits.Exec(endingNumbers) == nil {
			return nil
		}

		// If it ends only with numbers above 24, e.g. "at 25"
		endingNumberVal := js.ParseInt(endingNumbers, 0)
		if endingNumberVal > 24 {
			return nil
		}
	}

	return result
}

func (p *AbstractTimeExpressionParser) checkAndReturnWithFollowingPattern(result *ParsingResult) *ParsingResult {
	if reNumberDashNumber.Exec(result.Text) != nil {
		return nil
	}

	// If it ends only with numbers or dots
	if endingWithNumbers := reEndingWithRange.Exec(result.Text); endingWithNumbers != nil {
		// In strict mode (e.g. "at 1-3" or "at 1.2 - 2.3"), this should not be accepted
		if p.StrictMode {
			return nil
		}

		startingNumbers := endingWithNumbers.Get(1)
		endingNumbers := endingWithNumbers.Get(2)
		// If it ends only with dot single digit, e.g. "at 1.2"
		if js.Includes(endingNumbers, ".") && reDotTwoDigits.Exec(endingNumbers) == nil {
			return nil
		}

		// If it ends only with numbers above 24, e.g. "at 25"
		endingNumberVal := js.ParseInt(endingNumbers, 0)
		startingNumberVal := js.ParseInt(startingNumbers, 0)
		if endingNumberVal > 24 || startingNumberVal > 24 {
			return nil
		}
	}

	return result
}

// GetPrimaryTimePatternThroughCache returns the primary pattern, rebuilt when the prefix or the
// suffix changed.
func (p *AbstractTimeExpressionParser) GetPrimaryTimePatternThroughCache() *RegExp {
	primaryPrefix := p.Self.PrimaryPrefix()
	primarySuffix := p.Self.PrimarySuffix()

	if p.cachedPrimaryPrefix != nil && *p.cachedPrimaryPrefix == primaryPrefix &&
		p.cachedPrimarySuffix != nil && *p.cachedPrimarySuffix == primarySuffix {
		return p.cachedPrimaryTimePattern
	}

	p.cachedPrimaryTimePattern = primaryTimePattern(
		p.Self.PrimaryPatternLeftBoundary(),
		primaryPrefix,
		primarySuffix,
		p.Self.PatternFlags(),
	)
	p.cachedPrimaryPrefix = &primaryPrefix
	p.cachedPrimarySuffix = &primarySuffix
	return p.cachedPrimaryTimePattern
}

// GetFollowingTimePatternThroughCache returns the following pattern, rebuilt when the phase or
// the suffix changed.
func (p *AbstractTimeExpressionParser) GetFollowingTimePatternThroughCache() *RegExp {
	followingPhase := p.Self.FollowingPhase()
	followingSuffix := p.Self.FollowingSuffix()

	if p.cachedFollowingPhase != nil && *p.cachedFollowingPhase == followingPhase &&
		p.cachedFollowingSuffix != nil && *p.cachedFollowingSuffix == followingSuffix {
		return p.cachedFollowingTimePatten
	}

	p.cachedFollowingTimePatten = followingTimePatten(followingPhase, followingSuffix)
	p.cachedFollowingPhase = &followingPhase
	p.cachedFollowingSuffix = &followingSuffix
	return p.cachedFollowingTimePatten
}
