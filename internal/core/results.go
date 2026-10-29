// Ported from src/results.ts.

package core

import (
	"math"
	"strings"

	"github.com/wanasit/chrono-go/internal/jsdate"
)

// ReferenceWithTimezone is the reference instant plus the reference time zone offset (minutes east
// of UTC, nil for JS null) of one parse call.
type ReferenceWithTimezone struct {
	// Instant is the reference Date's time value.
	Instant float64
	// TimezoneOffset is the reference offset in minutes, nil when null.
	TimezoneOffset *float64
}

// NewReferenceWithTimezone is `new ReferenceWithTimezone(instant, timezoneOffset)`.
func NewReferenceWithTimezone(instant float64, timezoneOffset *float64) *ReferenceWithTimezone {
	return &ReferenceWithTimezone{Instant: instant, TimezoneOffset: timezoneOffset}
}

// ReferenceFromDate is `ReferenceWithTimezone.fromDate(date)` (date is a time value).
func ReferenceFromDate(date float64) *ReferenceWithTimezone {
	return NewReferenceWithTimezone(date, nil)
}

// ReferenceFromInput is `ReferenceWithTimezone.fromInput(input, timezoneOverrides)`.
func ReferenceFromInput(input ParsingReference, timezoneOverrides map[string]TimezoneAbbr) *ReferenceWithTimezone {
	instant := input.Instant
	offset, ok := ToTimezoneOffset(input.Timezone, jsdate.FromTime(instant), timezoneOverrides)
	if !ok {
		return NewReferenceWithTimezone(instant, nil)
	}
	return NewReferenceWithTimezone(instant, &offset)
}

func (r *ReferenceWithTimezone) instantDate() *jsdate.Date { return jsdate.FromTime(r.Instant) }

// DateWithAdjustedTimezone returns a JS date (system timezone) with the { year, month, day, hour,
// minute, second } equal to the reference (getDateWithAdjustedTimezone). The result is a fresh
// Date the caller may modify.
func (r *ReferenceWithTimezone) DateWithAdjustedTimezone() *jsdate.Date {
	date := jsdate.FromTime(r.Instant)
	if r.TimezoneOffset != nil {
		date.SetMinutes(date.GetMinutes() - r.GetSystemTimezoneAdjustmentMinute(r.instantDate(), nil))
	}
	return date
}

// GetDateWithAdjustedTimezone is getDateWithAdjustedTimezone() as a time value.
func (r *ReferenceWithTimezone) GetDateWithAdjustedTimezone() float64 {
	return r.DateWithAdjustedTimezone().Time()
}

// GetSystemTimezoneAdjustmentMinute returns the number of minutes between the JS date's time zone
// and the reference time zone. date nil is `new Date()`; overrideTimezoneOffset nil is undefined.
func (r *ReferenceWithTimezone) GetSystemTimezoneAdjustmentMinute(date *jsdate.Date, overrideTimezoneOffset *float64) float64 {
	if date == nil {
		date = jsdate.Now()
	}
	currentTimezoneOffset := -date.GetTimezoneOffset()
	targetTimezoneOffset := currentTimezoneOffset
	if overrideTimezoneOffset != nil {
		targetTimezoneOffset = *overrideTimezoneOffset
	} else if r.TimezoneOffset != nil {
		targetTimezoneOffset = *r.TimezoneOffset
	}
	return currentTimezoneOffset - targetTimezoneOffset
}

// GetTimezoneOffset returns the reference offset, or the system offset at the instant when null.
func (r *ReferenceWithTimezone) GetTimezoneOffset() float64 {
	if r.TimezoneOffset != nil {
		return *r.TimezoneOffset
	}
	return -r.instantDate().GetTimezoneOffset()
}

// ParsingComponents is a collection of parsed date/time components with their certainty: a
// component is certain (known), implied, or unknown (null).
type ParsingComponents struct {
	knownValues   numObject
	impliedValues numObject
	reference     *ReferenceWithTimezone
	tags          []string
}

// NewParsingComponents is `new ParsingComponents(reference, knownComponents)`; known nil is an
// absent argument.
func NewParsingComponents(reference *ReferenceWithTimezone, known Components) *ParsingComponents {
	c := &ParsingComponents{reference: reference}
	for _, kv := range known {
		c.knownValues.set(string(kv.Component), kv.Value)
	}
	date := reference.DateWithAdjustedTimezone()
	c.Imply(ComponentDay, date.GetDate())
	c.Imply(ComponentMonth, date.GetMonth()+1)
	c.Imply(ComponentYear, date.GetFullYear())
	c.Imply(ComponentHour, 12)
	c.Imply(ComponentMinute, 0)
	c.Imply(ComponentSecond, 0)
	c.Imply(ComponentMillisecond, 0)
	return c
}

// CreateRelativeFromReference is `ParsingComponents.createRelativeFromReference(reference,
// duration)`. Pass EmptyDuration for the default argument; a nil duration is JS null (and throws,
// as in the original). The duration is modified as addDuration modifies it.
func CreateRelativeFromReference(reference *ReferenceWithTimezone, duration *Duration) *ParsingComponents {
	date := AddDuration(reference.DateWithAdjustedTimezone(), duration)

	components := NewParsingComponents(reference, nil)
	components.AddTag("result/relativeDate")
	if duration.Has("hour") || duration.Has("minute") || duration.Has("second") || duration.Has("millisecond") {
		components.AddTag("result/relativeDateAndTime")
		AssignSimilarTime(components, date)
		AssignSimilarDate(components, date)
		components.Assign(ComponentTimezoneOffset, reference.GetTimezoneOffset())
	} else {
		ImplySimilarTime(components, date)
		components.Imply(ComponentTimezoneOffset, reference.GetTimezoneOffset())

		if duration.Has("day") {
			components.Assign(ComponentDay, date.GetDate())
			components.Assign(ComponentMonth, date.GetMonth()+1)
			components.Assign(ComponentYear, date.GetFullYear())
			components.Assign(ComponentWeekday, date.GetDay())
		} else if duration.Has("week") {
			components.Assign(ComponentDay, date.GetDate())
			components.Assign(ComponentMonth, date.GetMonth()+1)
			components.Assign(ComponentYear, date.GetFullYear())
			components.Imply(ComponentWeekday, date.GetDay())
		} else {
			components.Imply(ComponentDay, date.GetDate())
			if duration.Has("month") {
				components.Assign(ComponentMonth, date.GetMonth()+1)
				components.Assign(ComponentYear, date.GetFullYear())
			} else {
				components.Imply(ComponentMonth, date.GetMonth()+1)
				if duration.Has("year") {
					components.Assign(ComponentYear, date.GetFullYear())
				} else {
					components.Imply(ComponentYear, date.GetFullYear())
				}
			}
		}
	}
	return components
}

// Reference returns the components' reference.
func (c *ParsingComponents) Reference() *ReferenceWithTimezone { return c.reference }

// Get returns the certain or implied value of the component. It returns (0, false) when the
// value is null and (NaN, false) when it is undefined, so arithmetic on the returned value
// matches JavaScript's (null -> 0, undefined -> NaN) and `ok` is `get(c) != null`.
func (c *ParsingComponents) Get(component Component) (float64, bool) {
	v := c.Raw(component)
	if IsUndefined(v) {
		return math.NaN(), false
	}
	if IsNullMarker(v) {
		return 0, false
	}
	return v, true
}

// Raw returns get(component) as a value that can be passed on to Assign/Imply unchanged: the
// number, or the Null / Undefined marker.
func (c *ParsingComponents) Raw(component Component) float64 {
	if v, ok := c.knownValues.get(string(component)); ok {
		return v
	}
	if v, ok := c.impliedValues.get(string(component)); ok {
		return v
	}
	return Null
}

// V returns the value of Get without the ok flag (null -> 0, undefined -> NaN), for arithmetic.
func (c *ParsingComponents) V(component Component) float64 {
	v, _ := c.Get(component)
	return v
}

// IsNull reports whether `get(component) === null` (missing, or stored as null; undefined is not
// === null).
func (c *ParsingComponents) IsNull(component Component) bool {
	return IsNullMarker(c.Raw(component))
}

// IsCertain reports whether the component is known.
func (c *ParsingComponents) IsCertain(component Component) bool {
	return c.knownValues.has(string(component))
}

// GetCertainComponents returns Object.keys(knownValues).
func (c *ParsingComponents) GetCertainComponents() []Component {
	out := make([]Component, len(c.knownValues.entries))
	for i, e := range c.knownValues.entries {
		out[i] = Component(e.key)
	}
	return out
}

// Imply sets the implied value unless the component is known.
func (c *ParsingComponents) Imply(component Component, value float64) *ParsingComponents {
	if c.knownValues.has(string(component)) {
		return c
	}
	c.impliedValues.set(string(component), value)
	return c
}

// Assign sets the known value and removes the implied one.
func (c *ParsingComponents) Assign(component Component, value float64) *ParsingComponents {
	c.knownValues.set(string(component), value)
	c.impliedValues.del(string(component))
	return c
}

// AddDurationAsImplied adds the duration to these components, marking the modified date and/or
// time as implied.
func (c *ParsingComponents) AddDurationAsImplied(duration *Duration) *ParsingComponents {
	currentDate := c.dateWithoutTimezoneAdjustment()
	date := AddDuration(currentDate, duration)
	if duration.Has("day") || duration.Has("week") || duration.Has("month") || duration.Has("year") {
		c.Delete(ComponentDay, ComponentWeekday, ComponentMonth, ComponentYear)
		c.Imply(ComponentDay, date.GetDate())
		c.Imply(ComponentWeekday, date.GetDay())
		c.Imply(ComponentMonth, date.GetMonth()+1)
		c.Imply(ComponentYear, date.GetFullYear())
	}
	if duration.Has("second") || duration.Has("minute") || duration.Has("hour") {
		c.Delete(ComponentSecond, ComponentMinute, ComponentHour)
		c.Imply(ComponentSecond, date.GetSeconds())
		c.Imply(ComponentMinute, date.GetMinutes())
		c.Imply(ComponentHour, date.GetHours())
	}
	return c
}

// Delete removes the components (known and implied).
func (c *ParsingComponents) Delete(components ...Component) {
	for _, component := range components {
		c.knownValues.del(string(component))
		c.impliedValues.del(string(component))
	}
}

// Clone copies the known and implied values (not the tags, as in the original).
func (c *ParsingComponents) Clone() *ParsingComponents {
	component := NewParsingComponents(c.reference, nil)
	component.knownValues = c.knownValues.clone()
	component.impliedValues = c.impliedValues.clone()
	return component
}

// IsOnlyDate reports that no hour, minute or second is certain.
func (c *ParsingComponents) IsOnlyDate() bool {
	return !c.IsCertain(ComponentHour) && !c.IsCertain(ComponentMinute) && !c.IsCertain(ComponentSecond)
}

// IsOnlyTime reports that no weekday, day, month or year is certain.
func (c *ParsingComponents) IsOnlyTime() bool {
	return !c.IsCertain(ComponentWeekday) && !c.IsCertain(ComponentDay) && !c.IsCertain(ComponentMonth) && !c.IsCertain(ComponentYear)
}

// IsOnlyWeekdayComponent reports that the weekday is certain but the day and month are not.
func (c *ParsingComponents) IsOnlyWeekdayComponent() bool {
	return c.IsCertain(ComponentWeekday) && !c.IsCertain(ComponentDay) && !c.IsCertain(ComponentMonth)
}

// IsDateWithUnknownYear reports that the month is certain but the year is not.
func (c *ParsingComponents) IsDateWithUnknownYear() bool {
	return c.IsCertain(ComponentMonth) && !c.IsCertain(ComponentYear)
}

// IsValidDate checks the calendar validity of the components in UTC.
func (c *ParsingComponents) IsValidDate() bool {
	year, yearOK := c.Get(ComponentYear)
	month := c.V(ComponentMonth)
	day, dayOK := c.Get(ComponentDay)
	hour, hourOK := c.Get(ComponentHour)
	minute, minuteOK := c.Get(ComponentMinute)
	date := jsdate.FromTime(jsdate.UTC(year, month-1, day, hour, minute, c.V(ComponentSecond), c.V(ComponentMillisecond)))
	date.SetUTCFullYear(year)

	// `x !== null` is true for undefined (NaN) as well; NaN never compares equal.
	if !yearOK || date.GetUTCFullYear() != year {
		return false
	}
	if _, ok := c.Get(ComponentMonth); !ok || date.GetUTCMonth() != month-1 {
		return false
	}
	if !dayOK || date.GetUTCDate() != day {
		return false
	}
	if hourOK && date.GetUTCHours() != hour {
		return false
	}
	if minuteOK && date.GetUTCMinutes() != minute {
		return false
	}
	return true
}

// Date returns the time value of the components (NaN = Invalid Date).
func (c *ParsingComponents) Date() float64 { return c.DateObject().Time() }

// DateObject is date(): a new JS Date for the components.
func (c *ParsingComponents) DateObject() *jsdate.Date {
	timezoneOffset, ok := c.Get(ComponentTimezoneOffset)
	if !ok {
		if c.reference.TimezoneOffset == nil {
			// No timezone known: interpret the components as system-local time.
			return c.dateWithoutTimezoneAdjustment()
		}
		timezoneOffset = *c.reference.TimezoneOffset
	}
	date := jsdate.FromTime(jsdate.UTC(c.V(ComponentYear), c.V(ComponentMonth)-1, c.V(ComponentDay),
		c.V(ComponentHour), c.V(ComponentMinute), c.V(ComponentSecond), c.V(ComponentMillisecond)))
	date.SetUTCFullYear(c.V(ComponentYear))
	return jsdate.FromTime(date.Time() - timezoneOffset*60000)
}

// AddTag adds a debugging tag.
func (c *ParsingComponents) AddTag(tag string) *ParsingComponents {
	for _, t := range c.tags {
		if t == tag {
			return c
		}
	}
	c.tags = append(c.tags, tag)
	return c
}

// AddTags adds the tags in order.
func (c *ParsingComponents) AddTags(tags []string) *ParsingComponents {
	for _, tag := range tags {
		c.AddTag(tag)
	}
	return c
}

// Tags returns the tags in Set insertion order.
func (c *ParsingComponents) Tags() []string { return append([]string(nil), c.tags...) }

// HasTagPrefix reports whether some tag starts with prefix.
func (c *ParsingComponents) HasTagPrefix(prefix string) bool {
	for _, t := range c.tags {
		if strings.HasPrefix(t, prefix) {
			return true
		}
	}
	return false
}

func (c *ParsingComponents) dateWithoutTimezoneAdjustment() *jsdate.Date {
	date := jsdate.NewLocal(c.V(ComponentYear), c.V(ComponentMonth)-1, c.V(ComponentDay),
		c.V(ComponentHour), c.V(ComponentMinute), c.V(ComponentSecond), c.V(ComponentMillisecond))
	date.SetFullYear(c.V(ComponentYear))
	return date
}

// ParsingResult is one parsed result: a date/time (or range) mentioned in the input.
type ParsingResult struct {
	// RefDate is the reference instant's time value.
	RefDate float64
	// Index is the UTF-16 index of Text in the input.
	Index int
	// Text is the matched text.
	Text string
	// Reference is the parse call's reference.
	Reference *ReferenceWithTimezone
	// Start is the start (or only) date.
	Start *ParsingComponents
	// End is the end of a range, nil when absent.
	End *ParsingComponents
}

// NewParsingResult is `new ParsingResult(reference, index, text, start, end)`; a nil start is
// replaced by new components.
func NewParsingResult(reference *ReferenceWithTimezone, index int, text string, start, end *ParsingComponents) *ParsingResult {
	if start == nil {
		start = NewParsingComponents(reference, nil)
	}
	return &ParsingResult{RefDate: reference.Instant, Index: index, Text: text, Reference: reference, Start: start, End: end}
}

// Clone copies the result and its components.
func (r *ParsingResult) Clone() *ParsingResult {
	result := NewParsingResult(r.Reference, r.Index, r.Text, nil, nil)
	if r.Start != nil {
		result.Start = r.Start.Clone()
	} else {
		result.Start = nil
	}
	if r.End != nil {
		result.End = r.End.Clone()
	}
	return result
}

// Date returns the start date's time value.
func (r *ParsingResult) Date() float64 { return r.Start.Date() }

// AddTag adds the tag to the start and end components.
func (r *ParsingResult) AddTag(tag string) *ParsingResult {
	r.Start.AddTag(tag)
	if r.End != nil {
		r.End.AddTag(tag)
	}
	return r
}

// AddTags adds the tags to the start and end components.
func (r *ParsingResult) AddTags(tags []string) *ParsingResult {
	r.Start.AddTags(tags)
	if r.End != nil {
		r.End.AddTags(tags)
	}
	return r
}

// Tags returns the start tags followed by the new end tags (Set order).
func (r *ParsingResult) Tags() []string {
	combined := &ParsingComponents{tags: r.Start.Tags()}
	if r.End != nil {
		combined.AddTags(r.End.tags)
	}
	return combined.tags
}
