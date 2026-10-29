// Package core is the language-independent engine of chrono-go: the result and component types,
// the Chrono parse loop, time zones, date calculations, pattern helpers, the common parsers and
// refiners shared by every locale, and the abstract parser/refiner building blocks the locale
// packages extend.
//
// Every file is a line-by-line port of the corresponding TypeScript file of chrono v2.10.1
// (https://github.com/wanasit/chrono). JavaScript semantics are reproduced exactly: strings are
// UTF-16 (Go strings hold WTF-8, every index is a UTF-16 code unit index), regular expressions are
// ECMAScript regular expressions (internal/jsregex), dates are JavaScript Dates (internal/jsdate),
// numbers are float64 and object key order follows JavaScript's own-property order.
//
// The public package github.com/wanasit/chrono-go re-exports the API types of this package as
// aliases; see PORTING.md at the module root for how locale packages use it.
package core
