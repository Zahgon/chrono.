// Helpers for JavaScript object semantics used by the ported files (no direct TS counterpart).

package core

import (
	"math"

	"github.com/wanasit/chrono-go/internal/js"
)

const (
	undefinedBits = 0x7ff8_0000_0000_0d1f
	nullBits      = 0x7ff8_0000_0000_0d1e
)

// Undefined is a NaN with a private payload that stands for the JavaScript value undefined where
// the original stores it in a component or duration object (for example the result of a failed
// dictionary lookup). Store it only as it is: arithmetic on it yields an ordinary NaN in JS, so
// compute with Num(v) when a looked-up value may be Undefined.
var Undefined = math.Float64frombits(undefinedBits)

// Null is a NaN with a private payload that stands for the JavaScript value null stored as a
// component value (for example `assign("hour", other.get("hour"))` when the other hour is null).
// ParsingComponents.Raw returns it; pass it on unchanged.
var Null = math.Float64frombits(nullBits)

// IsUndefined reports whether v is the Undefined marker.
func IsUndefined(v float64) bool { return math.Float64bits(v) == undefinedBits }

// IsNullMarker reports whether v is the Null marker.
func IsNullMarker(v float64) bool { return math.Float64bits(v) == nullBits }

// IsNullish reports whether v is the Null or the Undefined marker (JS `v == null`).
func IsNullish(v float64) bool { return IsUndefined(v) || IsNullMarker(v) }

// Num is JS ToNumber for a value that may be a marker: Undefined becomes NaN, Null becomes 0, and
// any other value is returned unchanged.
func Num(v float64) float64 {
	if IsUndefined(v) {
		return math.NaN()
	}
	if IsNullMarker(v) {
		return 0
	}
	return v
}

// Truthy is JavaScript ToBoolean for a number.
func Truthy(v float64) bool { return v != 0 && !math.IsNaN(v) }

// objectPrototypeKeys are the properties every plain JavaScript object inherits from
// Object.prototype; `obj[key]` on a dictionary yields a function (or an object) for them.
var objectPrototypeKeys = map[string]bool{
	"constructor": true, "__defineGetter__": true, "__defineSetter__": true, "hasOwnProperty": true,
	"__lookupGetter__": true, "__lookupSetter__": true, "isPrototypeOf": true,
	"propertyIsEnumerable": true, "toString": true, "valueOf": true, "__proto__": true,
	"toLocaleString": true,
}

// IsObjectPrototypeKey reports whether a plain JavaScript object `{}` has an inherited property
// named key (so `dict[key]` is not undefined even though dict has no such own property).
func IsObjectPrototypeKey(key string) bool { return objectPrototypeKeys[key] }

// numEntry is one property of a numeric JS object.
type numEntry struct {
	key string
	v   float64
}

// numObject is a JavaScript object whose values are numbers, iterating in own-property order:
// canonical array-index keys first (ascending), then the other keys in insertion order. Setting
// an existing key keeps its position; deleting a key and setting it again moves it to the end.
type numObject struct {
	entries []numEntry
}

func (o *numObject) find(key string) int {
	for i := range o.entries {
		if o.entries[i].key == key {
			return i
		}
	}
	return -1
}

func (o *numObject) has(key string) bool { return o.find(key) >= 0 }

func (o *numObject) get(key string) (float64, bool) {
	if i := o.find(key); i >= 0 {
		return o.entries[i].v, true
	}
	return 0, false
}

func (o *numObject) set(key string, v float64) {
	if i := o.find(key); i >= 0 {
		o.entries[i].v = v
		return
	}
	if n, ok := js.ArrayIndex(key); ok {
		pos := 0
		for pos < len(o.entries) {
			m, isIndex := js.ArrayIndex(o.entries[pos].key)
			if !isIndex || m > n {
				break
			}
			pos++
		}
		o.entries = append(o.entries, numEntry{})
		copy(o.entries[pos+1:], o.entries[pos:])
		o.entries[pos] = numEntry{key, v}
		return
	}
	o.entries = append(o.entries, numEntry{key, v})
}

func (o *numObject) del(key string) {
	if i := o.find(key); i >= 0 {
		o.entries = append(o.entries[:i:i], o.entries[i+1:]...)
	}
}

func (o *numObject) clone() numObject {
	return numObject{entries: append([]numEntry(nil), o.entries...)}
}

// Entry is one key/value pair of a Dictionary, as written in a JavaScript object literal.
type Entry[V any] struct {
	Key   string
	Value V
}

// Dictionary is a JavaScript object literal used as a dictionary (for example a locale's
// MONTH_DICTIONARY). Keys iterate in JavaScript own-property order, which matters for
// MatchAnyPattern: canonical array-index keys first in ascending order, then the other keys in
// insertion order; a repeated key keeps its first position and takes the last value (as in an
// object literal or a spread).
type Dictionary[V any] struct {
	keys []string
	vals map[string]V
	obj  numObject // key order bookkeeping
}

// NewDictionary builds a dictionary from object-literal entries in source order. Spread
// (`{...A, k: v}`) is NewDictionary(append(A.Entries(), Entry{k, v})...).
func NewDictionary[V any](entries ...Entry[V]) *Dictionary[V] {
	d := &Dictionary[V]{vals: map[string]V{}}
	for _, e := range entries {
		d.Set(e.Key, e.Value)
	}
	return d
}

// Set is `dict[key] = value`.
func (d *Dictionary[V]) Set(key string, value V) {
	if d.vals == nil {
		d.vals = map[string]V{}
	}
	if _, ok := d.vals[key]; !ok {
		d.obj.set(key, 0)
		d.keys = nil
	}
	d.vals[key] = value
}

// Get returns the own property key (`dict[key]` when it is not undefined).
func (d *Dictionary[V]) Get(key string) (V, bool) {
	v, ok := d.vals[key]
	return v, ok
}

// Has reports whether key is an own property.
func (d *Dictionary[V]) Has(key string) bool {
	_, ok := d.vals[key]
	return ok
}

// Keys returns Object.keys(dict).
func (d *Dictionary[V]) Keys() []string {
	if d.keys == nil {
		d.keys = make([]string, len(d.obj.entries))
		for i, e := range d.obj.entries {
			d.keys[i] = e.key
		}
	}
	return append([]string(nil), d.keys...)
}

// Entries returns Object.entries(dict) as Entry values, in key order.
func (d *Dictionary[V]) Entries() []Entry[V] {
	out := make([]Entry[V], 0, len(d.obj.entries))
	for _, e := range d.obj.entries {
		out = append(out, Entry[V]{e.key, d.vals[e.key]})
	}
	return out
}

// Len returns the number of keys.
func (d *Dictionary[V]) Len() int { return len(d.obj.entries) }
