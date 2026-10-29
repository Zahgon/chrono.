package js

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

// JSONAppender is implemented by values that write their own JSON text. StringifyIndent inserts
// that text unchanged.
type JSONAppender interface {
	AppendJSON(b []byte) []byte
}

const hexDigits = "0123456789abcdef"

// Quote is JSON.stringify(s) for a string: it escapes '"', '\\', \b \f \n \r \t, the other C0
// controls as \u00xx and lone surrogates as \udxxx (lowercase hex), and nothing else.
func Quote(s string) string { return string(AppendQuote(nil, s)) }

// AppendQuote appends Quote(s) to b.
func AppendQuote(b []byte, s string) []byte {
	b = append(b, '"')
	start := 0
	for i := 0; i < len(s); {
		c := s[i]
		if c >= 0x20 && c < 0x80 && c != '"' && c != '\\' {
			i++
			continue
		}
		if c >= 0x80 {
			r, size := decodeUnit(s, i)
			switch {
			case r >= 0xd800 && r < 0xdc00 && i+size < len(s):
				if r2, size2 := decodeUnit(s, i+size); r2 >= 0xdc00 && r2 < 0xe000 {
					// A pair spelled as two lone-surrogate forms.
					b = append(b, s[start:i]...)
					p := (r-0xd800)<<10 + (r2 - 0xdc00) + 0x10000
					b = appendCodePoint(b, p)
					i += size + size2
					start = i
					continue
				}
				fallthrough
			case r >= 0xd800 && r < 0xe000:
				b = append(b, s[start:i]...)
				b = append(b, '\\', 'u', hexDigits[r>>12], hexDigits[r>>8&0xf], hexDigits[r>>4&0xf], hexDigits[r&0xf])
				i += size
				start = i
				continue
			case r == 0xfffd && size == 1:
				b = append(b, s[start:i]...)
				b = append(b, "\ufffd"...)
				i++
				start = i
				continue
			}
			i += size
			continue
		}
		b = append(b, s[start:i]...)
		switch c {
		case '"':
			b = append(b, '\\', '"')
		case '\\':
			b = append(b, '\\', '\\')
		case '\b':
			b = append(b, '\\', 'b')
		case '\f':
			b = append(b, '\\', 'f')
		case '\n':
			b = append(b, '\\', 'n')
		case '\r':
			b = append(b, '\\', 'r')
		case '\t':
			b = append(b, '\\', 't')
		default:
			b = append(b, '\\', 'u', '0', '0', hexDigits[c>>4], hexDigits[c&0xf])
		}
		i++
		start = i
	}
	b = append(b, s[start:]...)
	return append(b, '"')
}

// Stringify is JSON.stringify(v). Supported values: nil (null), Undefined, bool, the Go integer
// and float types (as numbers; NaN and the infinities become null, -0 becomes 0), string,
// *Object, Array, *Array, any other slice or array, and JSONAppender. For a top-level Undefined
// JavaScript returns undefined; Stringify returns "". Other types panic with a TypeError.
func Stringify(v any) string { return StringifyIndent(v, "") }

// StringifyIndent is JSON.stringify(v, null, indent): indent is cut to its first 10 code units
// and an empty indent gives the compact form.
func StringifyIndent(v any, indent string) string {
	if Len(indent) > 10 {
		indent = Slice(indent, 0, 10)
	}
	if IsUndefined(v) {
		return ""
	}
	w := jsonWriter{gap: indent}
	w.value(v, "")
	return string(w.b)
}

// StringifyIndentN is JSON.stringify(v, null, n): n (clamped to 10) spaces of indentation.
func StringifyIndentN(v any, n int) string {
	return StringifyIndent(v, strings.Repeat(" ", min(max(n, 0), 10)))
}

type jsonWriter struct {
	b   []byte
	gap string
}

// skip reports whether v is omitted as an object member.
func skip(v any) bool { return IsUndefined(v) }

func (w *jsonWriter) open(c byte, indent string, empty bool) string {
	w.b = append(w.b, c)
	if w.gap == "" || empty {
		return indent
	}
	return indent + w.gap
}

func (w *jsonWriter) sep(first bool, inner string) {
	if !first {
		w.b = append(w.b, ',')
	}
	if w.gap != "" {
		w.b = append(w.b, '\n')
		w.b = append(w.b, inner...)
	}
}

func (w *jsonWriter) close(c byte, indent string, any bool) {
	if w.gap != "" && any {
		w.b = append(w.b, '\n')
		w.b = append(w.b, indent...)
	}
	w.b = append(w.b, c)
}

func (w *jsonWriter) number(f float64) {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		w.b = append(w.b, "null"...)
		return
	}
	w.b = append(w.b, NumberToString(f)...)
}

func (w *jsonWriter) value(v any, indent string) {
	switch x := v.(type) {
	case nil:
		w.b = append(w.b, "null"...)
	case undefinedType:
		w.b = append(w.b, "null"...)
	case bool:
		w.b = strconv.AppendBool(w.b, x)
	case float64:
		w.number(x)
	case float32:
		w.number(float64(x))
	case int:
		w.number(float64(x))
	case int8:
		w.number(float64(x))
	case int16:
		w.number(float64(x))
	case int32:
		w.number(float64(x))
	case int64:
		w.number(float64(x))
	case uint:
		w.number(float64(x))
	case uint8:
		w.number(float64(x))
	case uint16:
		w.number(float64(x))
	case uint32:
		w.number(float64(x))
	case uint64:
		w.number(float64(x))
	case string:
		w.b = AppendQuote(w.b, x)
	case *Object:
		if x == nil {
			w.b = append(w.b, "null"...)
			return
		}
		w.object(x, indent)
	case Array:
		w.array(len(x), func(i int) any { return x[i] }, indent)
	case *Array:
		if x == nil {
			w.b = append(w.b, "null"...)
			return
		}
		w.array(len(*x), func(i int) any { return (*x)[i] }, indent)
	case []any:
		w.array(len(x), func(i int) any { return x[i] }, indent)
	case []string:
		w.array(len(x), func(i int) any { return x[i] }, indent)
	case []float64:
		w.array(len(x), func(i int) any { return x[i] }, indent)
	case []int:
		w.array(len(x), func(i int) any { return x[i] }, indent)
	case JSONAppender:
		w.b = x.AppendJSON(w.b)
	default:
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Slice, reflect.Array:
			if rv.Kind() == reflect.Slice && rv.IsNil() {
				w.array(0, nil, indent)
				return
			}
			w.array(rv.Len(), func(i int) any { return rv.Index(i).Interface() }, indent)
			return
		case reflect.Pointer:
			if rv.IsNil() {
				w.b = append(w.b, "null"...)
				return
			}
		}
		panic(&Error{Name: "TypeError", Message: fmt.Sprintf("js.Stringify: unsupported type %T", v)})
	}
}

func (w *jsonWriter) object(o *Object, indent string) {
	keys := o.Keys()
	count := 0
	for _, k := range keys {
		if !skip(o.vals[k]) {
			count++
		}
	}
	inner := w.open('{', indent, count == 0)
	first := true
	for _, k := range keys {
		v := o.vals[k]
		if skip(v) {
			continue
		}
		w.sep(first, inner)
		first = false
		w.b = AppendQuote(w.b, k)
		w.b = append(w.b, ':')
		if w.gap != "" {
			w.b = append(w.b, ' ')
		}
		w.value(v, inner)
	}
	w.close('}', indent, count > 0)
}

func (w *jsonWriter) array(n int, at func(int) any, indent string) {
	inner := w.open('[', indent, n == 0)
	for i := range n {
		w.sep(i == 0, inner)
		w.value(at(i), inner)
	}
	w.close(']', indent, n > 0)
}
