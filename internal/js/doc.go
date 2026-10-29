// Package js reproduces the JavaScript semantics of V8 (as shipped in Node.js 24, ICU 78.3,
// Unicode 17.0) for strings, numbers, JSON and Array.prototype.sort, for the chrono port: the
// String.prototype operations chrono uses (toLowerCase, toUpperCase, trim, substring, indexOf,
// split, replace/replaceAll with a string search value, ...), Number, parseInt, parseFloat,
// isNaN and JavaScript number-to-string conversion. String.prototype.match, and replace or split
// with a regular expression, are methods of jsregex.Regexp (s.match("x") is
// jsregex.MustCompile("x", "").Exec(s)).
//
// # String model
//
// A JavaScript string is represented by a Go string holding its WTF-8 encoding: UTF-8 in which a
// lone (unpaired) surrogate code unit is stored as its three-byte generalised UTF-8 form
// (0xED 0xA0..0xBF 0x80..0xBF). Every index, length and position used by this package counts
// UTF-16 code units, as JavaScript does. Functions that build strings always produce canonical
// WTF-8, where a high surrogate followed by a low surrogate is encoded as the four-byte form of
// the supplementary code point. Inputs that join a lone high and a lone low surrogate byte form
// (for example with Go's + operator) are still read as the code unit sequence they spell, and
// bytes that are not WTF-8 at all are read as U+FFFD.
//
// # Errors
//
// Operations that throw in JavaScript on invalid arguments (for example a RangeError from
// Number.prototype.toFixed(101)) panic with an *Error carrying the JavaScript error name and
// message. JSON parsing returns its SyntaxError as an *Error value instead.
package js
