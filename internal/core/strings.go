// JavaScript string helpers used by the ported files (no direct TS counterpart).

package core

import "github.com/wanasit/chrono-go/internal/js"

func jsLen(s string) int { return js.Len(s) }

// FirstChar is `s[0]` (the first UTF-16 code unit as a string; "" for an empty string, where JS
// gives undefined).
func FirstChar(s string) string { return js.CharAt(s, 0) }
