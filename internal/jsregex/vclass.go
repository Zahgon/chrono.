package jsregex

import "sort"

// v-mode (unicodeSets) character classes: ClassUnion, ClassIntersection,
// ClassSubtraction, nested classes, \q{...} and properties of strings,
// following V8's ParseClassUnion / ParseClassSetOperand. With the i flag V8
// closes each operand under simple case folding when it is parsed (except
// \p{ASCII}, \p{Any}, \p{Assigned} and single characters of \q{...}) and
// does not close the resulting class again; strings still match
// case-insensitively. V8 decides the folding of operands by the leaked i
// flag (see leakIcase) and that of strings by the group's own flags.

type operandType uint8

const (
	opCharacter operandType = iota
	opRange
	opClassEscape
	opStringDisjunction
	opNested
)

// classValue is a set of code points plus strings (never of length 1).
// may is V8's static "may contain strings" property, which decides whether
// negation is allowed regardless of the actual contents.
type classValue struct {
	set  rangeSet
	strs map[string][]rune
	may  bool
}

func (v *classValue) addString(s []rune) {
	if v.strs == nil {
		v.strs = map[string][]rune{}
	}
	v.strs[string(s)] = s
}

func (v *classValue) union(o *classValue) {
	v.may = v.may || o.may
	v.set = append(v.set, o.set...).norm()
	for _, s := range o.strs {
		v.addString(s)
	}
}

func isClassSetSyntaxChar(c rune) bool {
	switch c {
	case '(', ')', '[', ']', '{', '}', '/', '-', '\\', '|':
		return true
	}
	return false
}

func (p *parser) isClassSetReservedDoublePunctuator(c rune) bool {
	switch c {
	case '&', '!', '#', '$', '%', '*', '+', ',', '.', ':', ';', '<', '=', '>', '?', '@', '^', '`', '~':
		return p.next() == c
	}
	return false
}

// addFolded adds [lo, hi] (closed under case folding with i).
func (p *parser) addFolded(v *classValue, lo, hi rune) {
	r := rangeSet{lo, hi}
	if p.leakIcase {
		r = getFold(true).closure(r)
	}
	v.set = append(v.set, r...).norm()
}

func (p *parser) canonString(s []rune) []rune {
	if !p.flags.icase {
		return s
	}
	f := getFold(true)
	out := make([]rune, len(s))
	for i, c := range s {
		out[i] = f.canonical(c)
	}
	return out
}

// foldValue applies the i flag to a property operand: \p is closed, \P is
// the complement of the closed set (V8 folds before complementing in v mode).
func (p *parser) foldValue(v *classValue, negated bool, n1, n2 string) *classValue {
	if !p.leakIcase {
		return v
	}
	f := getFold(true)
	if n2 == "" && (n1 == "ASCII" || n1 == "Any" || n1 == "Assigned") {
		return v
	}
	if negated {
		base, _, _ := p.propertySet(n1, n2, false)
		v.set = f.closure(base).complement(0x10FFFF)
	} else {
		v.set = f.closure(v.set)
	}
	if v.strs != nil {
		old := v.strs
		v.strs = nil
		for _, s := range old {
			v.addString(p.canonString(s))
		}
	}
	return v
}

func (p *parser) classTermV(v *classValue) *term {
	t := &term{typ: tClass, set: v.set.norm(), strIcase: p.flags.icase}
	for _, s := range v.strs {
		t.strs = append(t.strs, s)
	}
	sort.Slice(t.strs, func(a, b int) bool {
		if len(t.strs[a]) != len(t.strs[b]) {
			return len(t.strs[a]) > len(t.strs[b])
		}
		return string(t.strs[a]) < string(t.strs[b])
	})
	return t
}

func (p *parser) parseClassSetCharacter() rune {
	c := p.cur
	if c == '\\' {
		switch p.next() {
		case 'b':
			p.advanceN(2)
			return '\b'
		case endMarker:
			p.fail(errEscapeAtEnd)
		}
		return p.parseCharacterEscape(true)
	}
	if c == endMarker {
		p.fail(errUnterminatedClass)
	}
	if isClassSetSyntaxChar(c) {
		p.fail(errInvalidCharInClass)
	}
	if p.isClassSetReservedDoublePunctuator(c) {
		p.fail(errInvalidClassSetOp)
	}
	p.advance()
	return c
}

func (p *parser) parseClassStringDisjunction(v *classValue) {
	p.advanceN(2)
	if p.cur != '{' {
		p.fail(errInvalidEscape)
	}
	p.advance()
	var s []rune
	flush := func() {
		if len(s) == 1 {
			v.set = append(v.set, s[0], s[0]).norm()
		} else {
			v.addString(p.canonString(s))
			v.may = true
		}
		s = nil
	}
	for p.cur != '}' {
		if p.cur == '|' {
			flush()
			p.advance()
			continue
		}
		s = append(s, p.parseClassSetCharacter())
	}
	flush()
	p.advance()
}

// parseClassSetOperand parses one operand. Characters are returned (not
// added); escapes and \q are merged into v; nested classes are returned.
func (p *parser) parseClassSetOperand(v *classValue) (operandType, *classValue, rune) {
	if p.cur == '\\' {
		switch nx := p.next(); nx {
		case 'q':
			p.parseClassStringDisjunction(v)
			return opStringDisjunction, nil, 0
		case 'p', 'P':
			p.advanceN(2)
			if n1, n2, ok := p.parsePropertyClassName(); ok {
				if set, strs, ok := p.propertySet(n1, n2, nx == 'P'); ok {
					pv := &classValue{set: set, may: strs != nil}
					for _, s := range strs {
						pv.addString(s)
					}
					v.union(p.foldValue(pv, nx == 'P', n1, n2))
					return opClassEscape, nil, 0
				}
			}
			p.fail(errInvalidClassProperty)
		case 'd', 'D', 's', 'S', 'w', 'W':
			s := p.classEscapeSet(nx, p.leakIcase)
			p.advanceN(2)
			if p.leakIcase {
				s = getFold(true).closure(s.norm())
			}
			v.set = append(v.set, s...).norm()
			return opClassEscape, nil, 0
		}
	}
	if p.cur == '[' {
		p.advance()
		neg := false
		if p.cur == '^' {
			neg = true
			p.advance()
		}
		return opNested, p.parseClassBodyV(neg), 0
	}
	return opCharacter, nil, p.parseClassSetCharacter()
}

func (p *parser) negateV(v *classValue, neg bool) *classValue {
	if !neg {
		return v
	}
	if v.may {
		p.fail(errNegatedClassWithStrings)
	}
	return &classValue{set: v.set.norm().complement(0x10FFFF)}
}

// parseClassBodyV parses ClassContents after '[' / '[^' in v mode.
func (p *parser) parseClassBodyV(neg bool) *classValue {
	if p.cur == ']' {
		p.advance()
		return p.negateV(&classValue{}, neg)
	}
	acc := &classValue{}
	typ, nested, ch := p.parseClassSetOperand(acc)
	first := acc
	if typ == opNested {
		first = nested
	}
	if (p.cur == '-' || p.cur == '&') && p.next() == p.cur {
		if typ == opCharacter {
			p.addFolded(acc, ch, ch)
		}
		if p.cur == '-' {
			return p.parseClassSetOp(neg, first, true)
		}
		return p.parseClassSetOp(neg, first, false)
	}
	return p.parseClassUnion(neg, acc, nested, typ, ch)
}

func (p *parser) parseClassUnion(neg bool, acc, firstNested *classValue, lastType operandType, ch rune) *classValue {
	result := &classValue{}
	if firstNested != nil {
		result.union(firstNested)
	}
	for p.hasMore && p.cur != ']' {
		if p.cur == '-' && p.next() == '-' {
			p.fail(errInvalidClassSetOp)
		}
		if p.cur == '-' {
			if lastType != opCharacter {
				p.advance()
				if p.cur == endMarker {
					p.fail(errUnterminatedClass)
				}
				p.fail(errInvalidCharacterClass)
			}
			from := ch
			p.advance()
			var typ operandType
			typ, _, ch = p.parseClassSetOperand(acc)
			if typ != opCharacter {
				p.fail(errInvalidCharacterClass)
			}
			if from > ch {
				p.fail(errOutOfOrderClass)
			}
			p.addFolded(acc, from, ch)
			lastType = opRange
			continue
		}
		if lastType == opCharacter {
			p.addFolded(acc, ch, ch)
		}
		var nested *classValue
		lastType, nested, ch = p.parseClassSetOperand(acc)
		if nested != nil {
			result.union(nested)
		}
	}
	if !p.hasMore {
		p.fail(errUnterminatedClass)
	}
	p.advance()
	if lastType == opCharacter {
		p.addFolded(acc, ch, ch)
	}
	result.union(acc)
	return p.negateV(result, neg)
}

// parseClassSetOp parses ClassSubtraction (sub) or ClassIntersection.
func (p *parser) parseClassSetOp(neg bool, first *classValue, sub bool) *classValue {
	cur := first
	may := first.may
	op := byte('&')
	if sub {
		op = '-'
		// The result of a subtraction may contain strings iff its first
		// operand may; V8 rejects the negation before parsing further.
		if neg && may {
			p.fail(errNegatedClassWithStrings)
		}
	}
	for p.hasMore && p.cur != ']' {
		if p.cur != rune(op) || p.next() != rune(op) {
			p.fail(errInvalidClassSetOp)
		}
		p.advanceN(2)
		if !sub && p.cur == '&' {
			p.fail(errInvalidCharInClass)
		}
		v := &classValue{}
		typ, nested, ch := p.parseClassSetOperand(v)
		switch typ {
		case opNested:
			v = nested
		case opCharacter:
			p.addFolded(v, ch, ch)
		}
		if !sub {
			may = may && v.may
		}
		nv := &classValue{}
		if sub {
			nv.set = cur.set.norm().subtract(v.set.norm())
		} else {
			nv.set = cur.set.norm().intersect(v.set.norm())
		}
		for k, s := range cur.strs {
			if _, in := v.strs[k]; in != sub {
				nv.addString(s)
			}
		}
		cur = nv
	}
	if !p.hasMore {
		p.fail(errUnterminatedClass)
	}
	p.advance()
	cur.may = may
	return p.negateV(cur, neg)
}
