package js

import (
	"math"
	"math/big"
	"strconv"
	"strings"
)

// shortestDigits returns the shortest round-trip decimal digits of a finite f > 0 and the
// exponent n with f = 0.digits × 10^n.
func shortestDigits(f float64) (string, int) {
	s := strconv.FormatFloat(f, 'e', -1, 64)
	mant, exp, _ := strings.Cut(s, "e")
	e, _ := strconv.Atoi(exp)
	return strings.Replace(mant, ".", "", 1), e + 1
}

// NumberToString is String(f) for a number (Number::toString, radix 10).
func NumberToString(f float64) string {
	switch {
	case math.IsNaN(f):
		return "NaN"
	case f == 0:
		return "0"
	case math.IsInf(f, 1):
		return "Infinity"
	case math.IsInf(f, -1):
		return "-Infinity"
	}
	if f == math.Trunc(f) && math.Abs(f) < 1e15 {
		return strconv.FormatInt(int64(f), 10)
	}
	var b strings.Builder
	if f < 0 {
		b.WriteByte('-')
		f = -f
	}
	digits, n := shortestDigits(f)
	k := len(digits)
	switch {
	case k <= n && n <= 21:
		b.WriteString(digits)
		b.WriteString(strings.Repeat("0", n-k))
	case 0 < n && n <= 21:
		b.WriteString(digits[:n])
		b.WriteByte('.')
		b.WriteString(digits[n:])
	case -6 < n && n <= 0:
		b.WriteString("0.")
		b.WriteString(strings.Repeat("0", -n))
		b.WriteString(digits)
	default:
		b.WriteString(digits[:1])
		if k > 1 {
			b.WriteByte('.')
			b.WriteString(digits[1:])
		}
		b.WriteByte('e')
		if n-1 >= 0 {
			b.WriteByte('+')
		}
		b.WriteString(strconv.Itoa(n - 1))
	}
	return b.String()
}

// dblExponent is double-conversion's Double::Exponent: f = significand × 2^exponent with an
// integer significand of at most 53 bits.
func dblExponent(f float64) int {
	biased := int(math.Float64bits(f)>>52) & 0x7ff
	if biased == 0 {
		return -1074
	}
	return biased - 1075
}

const radixChars = "0123456789abcdefghijklmnopqrstuvwxyz"

// NumberToStringRadix is Number.prototype.toString(radix) for radix 2..36, using V8's
// DoubleToRadixCString for radices other than 10. An invalid radix panics with a RangeError.
func NumberToStringRadix(f float64, radix int) string {
	if radix < 2 || radix > 36 {
		throwRange("toString() radix must be between 2 and 36")
	}
	if radix == 10 || math.IsNaN(f) || math.IsInf(f, 0) {
		return NumberToString(f)
	}
	if f == 0 {
		return "0"
	}
	negative := f < 0
	if negative {
		f = -f
	}
	integer := math.Floor(f)
	fraction := f - integer
	delta := float64(0.5 * (math.Nextafter(f, math.Inf(1)) - f))
	delta = math.Max(math.Nextafter(0, 1), delta)
	var frac []byte
	if fraction >= delta {
		frac = append(frac, '.')
		for {
			// The conversions keep the compiler from fusing these products into later adds.
			fraction = float64(fraction * float64(radix))
			delta = float64(delta * float64(radix))
			digit := int(fraction)
			frac = append(frac, radixChars[digit])
			fraction -= float64(digit)
			if fraction > 0.5 || (fraction == 0.5 && digit&1 != 0) {
				if fraction+delta > 1 {
					for {
						last := len(frac) - 1
						if last == 0 {
							integer++
							frac = frac[:0]
							break
						}
						c := frac[last]
						frac = frac[:last]
						d := strings.IndexByte(radixChars, c)
						if d+1 < radix {
							frac = append(frac, radixChars[d+1])
							break
						}
					}
					break
				}
			}
			if !(fraction >= delta) {
				break
			}
		}
		if len(frac) == 0 {
			frac = nil
		}
	}
	var rev []byte
	for dblExponent(integer/float64(radix)) > 0 {
		integer /= float64(radix)
		rev = append(rev, '0')
	}
	for {
		remainder := math.Mod(integer, float64(radix))
		rev = append(rev, radixChars[int(remainder)])
		integer = (integer - remainder) / float64(radix)
		if !(integer > 0) {
			break
		}
	}
	out := make([]byte, 0, len(rev)+len(frac)+1)
	if negative {
		out = append(out, '-')
	}
	for i := len(rev) - 1; i >= 0; i-- {
		out = append(out, rev[i])
	}
	return string(append(out, frac...))
}

// exactDecimal returns the exact decimal digits of a finite x > 0 (no trailing zeros) and the
// position of the decimal point: x = 0.digits × 10^point.
func exactDecimal(x float64) ([]byte, int) {
	mant, exp := math.Frexp(x)
	m := int64(mant * (1 << 53))
	exp -= 53
	n := big.NewInt(m)
	var digits string
	var point int
	if exp >= 0 {
		n.Lsh(n, uint(exp))
		digits = n.String()
		point = len(digits)
	} else {
		five := new(big.Int).Exp(big.NewInt(5), big.NewInt(int64(-exp)), nil)
		n.Mul(n, five)
		digits = n.String()
		point = len(digits) + exp
	}
	d := []byte(strings.TrimRight(digits, "0"))
	if len(d) == 0 {
		d = []byte("0")
	}
	return d, point
}

// roundHalfUp rounds 0.digits × 10^point to keep leading digits (keep may be <= 0), rounding
// ties away from zero. It returns the kept digits (exactly keep of them when keep > 0, or "1"
// after a carry when keep == 0) and the possibly incremented point.
func roundHalfUp(digits []byte, point, keep int) ([]byte, int) {
	if keep < 0 {
		return nil, point
	}
	out := make([]byte, keep)
	for i := range keep {
		if i < len(digits) {
			out[i] = digits[i]
		} else {
			out[i] = '0'
		}
	}
	if keep < len(digits) && digits[keep] >= '5' {
		i := keep - 1
		for ; i >= 0; i-- {
			if out[i] == '9' {
				out[i] = '0'
				continue
			}
			out[i]++
			break
		}
		if i < 0 {
			out = append([]byte{'1'}, out...)
			point++
			if keep > 0 {
				out = out[:keep]
			}
		}
	}
	return out, point
}

// ToFixed is Number.prototype.toFixed(digits). The result is computed from the exact binary
// value with ties rounded up; |x| >= 1e21 falls back to NumberToString. digits outside 0..100
// panics with a RangeError.
func ToFixed(x float64, digits int) string {
	if digits < 0 || digits > 100 {
		throwRange("toFixed() digits argument must be between 0 and 100")
	}
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return NumberToString(x)
	}
	sign := ""
	if x < 0 {
		sign = "-"
		x = -x
	}
	if x >= 1e21 {
		return sign + NumberToString(x)
	}
	m := "0"
	if x != 0 {
		d, point := exactDecimal(x)
		kept, np := roundHalfUp(d, point, point+digits)
		if point+digits > 0 && np > point {
			kept = append(kept, '0')
		}
		if m = strings.TrimLeft(string(kept), "0"); m == "" {
			m = "0"
		}
	}
	if digits != 0 {
		k := len(m)
		if k <= digits {
			m = strings.Repeat("0", digits+1-k) + m
			k = digits + 1
		}
		m = m[:k-digits] + "." + m[k-digits:]
	}
	return sign + m
}

// ToExponential is Number.prototype.toExponential(digits); digits < 0 selects the
// undefined-argument form (as many digits as needed). digits > 100 panics with a RangeError.
func ToExponential(x float64, digits int) string {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return NumberToString(x)
	}
	if digits > 100 {
		throwRange("toExponential() argument must be between 0 and 100")
	}
	sign := ""
	if x < 0 {
		sign = "-"
		x = -x
	}
	var m string
	e := 0
	switch {
	case x == 0:
		m = strings.Repeat("0", max(digits, 0)+1)
	case digits < 0:
		var n int
		m, n = shortestDigits(x)
		e = n - 1
	default:
		d, point := exactDecimal(x)
		kept, np := roundHalfUp(d, point, digits+1)
		m = string(kept)
		e = np - 1
	}
	if len(m) > 1 {
		m = m[:1] + "." + m[1:]
	}
	c := "+"
	if e < 0 {
		c = "-"
		e = -e
	}
	return sign + m + "e" + c + strconv.Itoa(e)
}

// ToPrecision is Number.prototype.toPrecision(precision) with ties rounded up; precision outside
// 1..100 panics with a RangeError.
func ToPrecision(x float64, precision int) string {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return NumberToString(x)
	}
	if precision < 1 || precision > 100 {
		throwRange("toPrecision() argument must be between 1 and 100")
	}
	sign := ""
	if x < 0 {
		sign = "-"
		x = -x
	}
	m := strings.Repeat("0", precision)
	e := 0
	if x != 0 {
		d, point := exactDecimal(x)
		kept, np := roundHalfUp(d, point, precision)
		m = string(kept)
		e = np - 1
		if e < -6 || e >= precision {
			if precision != 1 {
				m = m[:1] + "." + m[1:]
			}
			c := "+"
			if e < 0 {
				c = "-"
				e = -e
			}
			return sign + m + "e" + c + strconv.Itoa(e)
		}
	}
	switch {
	case e == precision-1:
		return sign + m
	case e >= 0:
		return sign + m[:e+1] + "." + m[e+1:]
	}
	return sign + "0." + strings.Repeat("0", -(e+1)) + m
}

// ToLocaleString is Number.prototype.toLocaleString() with no arguments in the en-US locale, as
// V8 with ICU 78 formats it: the shortest round-trip decimal of x is rounded half away from zero
// to at most three fraction digits, trailing fraction zeros are dropped, the integer part is
// grouped by thousands with ",", negative values (also -0 and values rounding to zero) keep
// their "-", and the infinities print as "∞".
func ToLocaleString(x float64) string {
	switch {
	case math.IsNaN(x):
		return "NaN"
	case math.IsInf(x, 1):
		return "∞"
	case math.IsInf(x, -1):
		return "-∞"
	}
	sign := ""
	if math.Signbit(x) {
		sign = "-"
		x = -x
	}
	var intPart, fracPart string
	if x == 0 {
		intPart = "0"
	} else {
		d, n := shortestDigits(x)
		kept, point := roundHalfUp([]byte(d), n, n+3)
		k := string(kept)
		switch {
		case k == "" || strings.Trim(k, "0") == "":
			intPart = "0"
		case point <= 0:
			intPart = "0"
			fracPart = strings.Repeat("0", -point) + k
		case point >= len(k):
			intPart = k + strings.Repeat("0", point-len(k))
		default:
			intPart, fracPart = k[:point], k[point:]
		}
		fracPart = strings.TrimRight(fracPart, "0")
	}
	var b strings.Builder
	b.WriteString(sign)
	for i := 0; i < len(intPart); i++ {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteByte(intPart[i])
	}
	if fracPart != "" {
		b.WriteByte('.')
		b.WriteString(fracPart)
	}
	return b.String()
}
