// Port of com/google/javascript/jscomp/regex/RegExpTree.java @ closure-compiler v20261021

/*
 * Copyright 2011 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package regex

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// RegExpTree_DIGITS is [0-9].
var RegExpTree_DIGITS = CharRanges_Inclusive('0', '9')

// RegExpTree_UCASE_LETTERS is [A-Z].
var RegExpTree_UCASE_LETTERS = CharRanges_Inclusive('A', 'Z')

// RegExpTree_LCASE_LETTERS is [a-z].
var RegExpTree_LCASE_LETTERS = CharRanges_Inclusive('a', 'z')

// RegExpTree_LETTERS is [A-Za-z].
var RegExpTree_LETTERS = RegExpTree_UCASE_LETTERS.Union(RegExpTree_LCASE_LETTERS)

// RegExpTree_WORD_CHARS is \w.
var RegExpTree_WORD_CHARS = RegExpTree_DIGITS.Union(RegExpTree_LETTERS).Union(CharRanges_WithMembers('_'))

// RegExpTree_INVERSE_WORD_CHARS is \W.
var RegExpTree_INVERSE_WORD_CHARS = CharRanges_ALL_CODE_UNITS.Difference(RegExpTree_WORD_CHARS)

// RegExpTree_SPACE_CHARS is \s.
var RegExpTree_SPACE_CHARS = CharRanges_WithMembers(
	'\t',
	'\n',
	'\u000b',
	'\u000c',
	'\r',
	' ',
	'\u00a0',
	// Unicode 3.0 Zs
	'\u1680',
	'\u180e',
	'\u2000',
	'\u2001',
	'\u2002',
	'\u2003',
	'\u2004',
	'\u2005',
	'\u2006',
	'\u2007',
	'\u2008',
	'\u2009',
	'\u200a',
	// Line terminator chars
	'\u2028',
	'\u2029',
	// Unicode 3.0 Zs
	'\u202f',
	'\u205f',
	'\u3000',
	// Byte order marker is a space character in ES5 but not ES3.
	'\ufeff')

// RegExpTree_IE_SPACE_CHARS: IE is broken around \s. IE (6, 7, 8 at least),
// only recognize these.
var RegExpTree_IE_SPACE_CHARS = CharRanges_WithMembers('\t', '\n', '\u000b', '\u000c', '\r', ' ')

// RegExpTree_IE_SPEC_ERRORS: IE is broken around \s. IE (6, 7, 8 at least),
// only recognize these.
var RegExpTree_IE_SPEC_ERRORS = RegExpTree_SPACE_CHARS.Difference(RegExpTree_IE_SPACE_CHARS)

// RegExpTree_NAMED_CHAR_GROUPS maps the class escapes \d \D \s \S \w \W to
// their ranges (in this order).
var RegExpTree_NAMED_CHAR_GROUPS = []struct {
	letter uint16
	ranges *CharRanges
}{
	{'d', RegExpTree_DIGITS},
	{'D', CharRanges_ALL_CODE_UNITS.Difference(RegExpTree_DIGITS)},
	{'s', RegExpTree_SPACE_CHARS},
	{'S', CharRanges_ALL_CODE_UNITS.Difference(RegExpTree_SPACE_CHARS)},
	{'w', RegExpTree_WORD_CHARS},
	{'W', RegExpTree_INVERSE_WORD_CHARS},
}

func regExpTree_namedCharGroup(letter uint16) *CharRanges {
	for _, g := range RegExpTree_NAMED_CHAR_GROUPS {
		if g.letter == letter {
			return g.ranges
		}
	}
	return nil
}

// RegExpTree_DOT_CHARSET is the charset matched by '.'.
var RegExpTree_DOT_CHARSET = NewRegExpTree_Charset(
	CharRanges_ALL_CODE_UNITS.Difference(
		CharRanges_WithMembers('\n', '\r', '\u2028', '\u2029')),
	CharRanges_EMPTY)

// RegExpTree_Charset represents a set of possible characters structured as
// [a-zA-Z] or [^a-zA-Z].
//
// Java: private static final class RegExpTree.Charset
type RegExpTree_Charset struct {
	ranges *CharRanges

	// Code units that were mentioned explicitly and that might be matched by
	// a group according to ECMAScript 5 but would not because of
	// specification violations in IE.
	ieExplicits *CharRanges
}

// NewRegExpTree_Charset constructs a Charset.
//
// Java: Charset(CharRanges, CharRanges)
func NewRegExpTree_Charset(ranges *CharRanges, ieExplicits *CharRanges) *RegExpTree_Charset {
	return &RegExpTree_Charset{
		ranges:      ranges,
		ieExplicits: ieExplicits,
	}
}

func (c *RegExpTree_Charset) ContainsAnchor() bool    { return false }
func (c *RegExpTree_Charset) HasCapturingGroup() bool { return false }
func (c *RegExpTree_Charset) NumCapturingGroups() int { return 0 }
func (c *RegExpTree_Charset) Children() []RegExpTree  { return emptyRegExpChildren }

// RegExpTree_Charset_ComplexityWordFolded is the complexity of ranges or of
// its inverse, whichever is smaller.
//
// Java: Charset.complexityWordFolded(CharRanges)
func RegExpTree_Charset_ComplexityWordFolded(ranges *CharRanges) int {
	return min(
		RegExpTree_Charset_ComplexityWordFoldedHelper(ranges),
		1+RegExpTree_Charset_ComplexityWordFoldedHelper(CharRanges_ALL_CODE_UNITS.Difference(ranges)))
}

// RegExpTree_Charset_ComplexityWordFoldedHelper is the complexity of ranges,
// using \w or \W where that is shorter.
//
// Java: Charset.complexityWordFoldedHelper(CharRanges)
func RegExpTree_Charset_ComplexityWordFoldedHelper(ranges *CharRanges) int {
	complexity := RegExpTree_DecomposedCharset_Complexity(ranges)
	if ranges.ContainsAll(RegExpTree_WORD_CHARS) {
		complexity = min(complexity, 1+RegExpTree_DecomposedCharset_Complexity(ranges.Difference(RegExpTree_WORD_CHARS)))
	}
	if ranges.ContainsAll(RegExpTree_INVERSE_WORD_CHARS) {
		complexity = min(
			complexity,
			1+RegExpTree_DecomposedCharset_Complexity(ranges.Difference(RegExpTree_INVERSE_WORD_CHARS)))
	}
	return complexity
}

// Simplify picks the simplest equivalent charset.
//
// Java: Charset.simplify(String)
func (c *RegExpTree_Charset) Simplify(flags string) RegExpTree {
	if c.ranges.IsEmpty() {
		return RegExpTree_NEVER_MATCHES
	}
	best := c.ranges
	if strings.IndexByte(flags, 'i') >= 0 {
		// Java: Set<CharRanges> options = new LinkedHashSet<>();
		// Go-only: language-adapter — a LinkedHashSet with CharRanges value
		// equality, as an insertion-ordered slice.
		var optionsBuf [7]*CharRanges
		options := optionsBuf[:0]
		add := func(option *CharRanges) {
			for _, o := range options {
				if o.Equals(option) {
					return
				}
			}
			options = append(options, option)
		}
		add(CaseCanonicalize_ExpandToAllMatched(c.ranges))
		add(CaseCanonicalize_ReduceToMinimum(c.ranges))

		lcaseLetters := c.ranges.Intersection(RegExpTree_LCASE_LETTERS)
		ucaseLetters := c.ranges.Intersection(RegExpTree_UCASE_LETTERS)

		lcaseLettersToUpper := lcaseLetters.Shift(-32)
		ucaseLettersToLower := ucaseLetters.Shift(32)

		add(c.ranges.Union(ucaseLettersToLower))
		add(c.ranges.Union(lcaseLettersToUpper))
		add(c.ranges.Union(lcaseLettersToUpper).Union(ucaseLettersToLower))

		add(c.ranges.Union(ucaseLettersToLower).Difference(ucaseLetters))
		add(c.ranges.Union(lcaseLettersToUpper).Difference(lcaseLetters))

		bestComplexity := RegExpTree_Charset_ComplexityWordFolded(c.ranges)

		for _, option := range options {
			complexity := RegExpTree_Charset_ComplexityWordFolded(option)
			if complexity < bestComplexity {
				bestComplexity = complexity
				best = option
			}
		}
	}

	if best.GetNumRanges() == 1 && best.End(0)-best.Start(0) == 1 {
		return NewRegExpTree_Text(stringFromUTF16Char(uint16(best.Start(0))))
	}

	if !best.Equals(c.ranges) {
		return NewRegExpTree_Charset(best, c.ieExplicits)
	}

	return c
}

// IsCaseSensitive reports whether the ranges left after removing the named
// groups are case sensitive.
//
// Java: Charset.isCaseSensitive()
func (c *RegExpTree_Charset) IsCaseSensitive() bool {
	// We could test
	//     !ranges.equals(CaseCanonicalize.expandToAllMatched(ranges))
	// but we get better optimizations by leaving the 'i' flag on in most
	// cases.

	// Check whether skipping all the character groups that are known
	// case-insensitive leaves us with something that matches the above
	// definition.
	withoutNamedGroups := c.Decompose().ranges
	return !withoutNamedGroups.Equals(CaseCanonicalize_ExpandToAllMatched(withoutNamedGroups))
}

// DecomposeRanges extracts named groups (\d, \s, ...) from ranges while that
// makes it simpler.
//
// Java: Charset.decompose(CharRanges, boolean)
func (c *RegExpTree_Charset) DecomposeRanges(ranges *CharRanges, inverted bool) *RegExpTree_DecomposedCharset {
	var namedGroups strings.Builder
	rangesInterIeExplicits := ranges.Intersection(c.ieExplicits)
	for {
		groupName := uint16(0)
		var simplest *CharRanges
		minComplexity := RegExpTree_DecomposedCharset_Complexity(ranges)
		for _, g := range RegExpTree_NAMED_CHAR_GROUPS {
			if ranges.ContainsAll(g.ranges) {
				withoutGroup := ranges.Difference(g.ranges).Union(rangesInterIeExplicits)
				complexity := RegExpTree_DecomposedCharset_Complexity(withoutGroup)
				if complexity < minComplexity {
					simplest = withoutGroup
					groupName = g.letter
					minComplexity = complexity
				}
			}
		}
		if simplest != nil {
			namedGroups.WriteByte('\\')
			namedGroups.WriteByte(byte(groupName))
			ranges = simplest
		} else {
			break
		}
	}
	return NewRegExpTree_DecomposedCharset(inverted, ranges, namedGroups.String())
}

// Decompose picks the simpler of the positive and the inverted
// decomposition.
//
// Java: Charset.decompose()
func (c *RegExpTree_Charset) Decompose() *RegExpTree_DecomposedCharset {
	negRanges := CharRanges_ALL_CODE_UNITS.Difference(c.ranges)
	if !c.ieExplicits.IsEmpty() {
		if negRanges.Intersection(c.ieExplicits).IsEmpty() {
			return c.DecomposeRanges(c.ranges, false)
		} else if c.ranges.Intersection(c.ieExplicits).IsEmpty() {
			return c.DecomposeRanges(negRanges, true)
		}
	}
	positive := c.DecomposeRanges(c.ranges, false)
	negative := c.DecomposeRanges(negRanges, true)
	if positive.Complexity() <= negative.Complexity() {
		return positive
	}
	return negative
}

// AppendSourceCode emits '.' or the decomposed charset.
//
// Java: Charset.appendSourceCode(StringBuilder)
func (c *RegExpTree_Charset) AppendSourceCode(sb *strings.Builder) {
	if RegExpTree_DOT_CHARSET.ranges.Equals(c.ranges) {
		sb.WriteByte('.')
		return
	}
	c.Decompose().AppendSourceCode(sb)
}

// AppendDebugInfo emits the ranges.
//
// Java: Charset.appendDebugInfo(StringBuilder)
func (c *RegExpTree_Charset) AppendDebugInfo(sb *strings.Builder) {
	sb.WriteString(c.ranges.String())
}

// Equals compares the ranges.
//
// Java: Charset.equals(Object)
func (c *RegExpTree_Charset) Equals(o any) bool {
	charset, ok := o.(*RegExpTree_Charset)
	return ok && c.ranges.Equals(charset.ranges)
}

// HashCode hashes the ranges.
//
// Java: Charset.hashCode()
func (c *RegExpTree_Charset) HashCode() int32 {
	return int32(uint32(c.ranges.HashCode()) ^ 0xdede2246)
}

// Java: RegExpTree.toString()
func (c *RegExpTree_Charset) String() string { return regExpTreeString(c) }

// RegExpTree_DecomposedCharset is the internal representation for []
// charsets.
//
// Java: record RegExpTree.DecomposedCharset(boolean, CharRanges, String)
type RegExpTree_DecomposedCharset struct {
	inverted    bool
	ranges      *CharRanges
	namedGroups string
}

// NewRegExpTree_DecomposedCharset constructs a DecomposedCharset.
//
// Java: DecomposedCharset(boolean, CharRanges, String)
func NewRegExpTree_DecomposedCharset(inverted bool, ranges *CharRanges, namedGroups string) *RegExpTree_DecomposedCharset {
	d := &RegExpTree_DecomposedCharset{}
	d.inverted = inverted
	d.ranges = ranges
	d.namedGroups = namedGroups
	return d
}

// Complexity estimates the length of the rendered charset.
//
// Java: DecomposedCharset.complexity()
func (d *RegExpTree_DecomposedCharset) Complexity() int {
	inverted := 0
	if d.inverted {
		inverted = 1
	}
	return inverted + jutil.StringLength(d.namedGroups) + RegExpTree_DecomposedCharset_Complexity(d.ranges)
}

// AppendSourceCode emits the charset.
//
// Java: DecomposedCharset.appendSourceCode(StringBuilder)
func (d *RegExpTree_DecomposedCharset) AppendSourceCode(sb *strings.Builder) {
	if d.ranges.IsEmpty() {
		if !d.inverted && jutil.StringLength(d.namedGroups) == 2 {
			sb.WriteString(d.namedGroups)
			return
		} else if d.ranges.IsEmpty() && d.namedGroups == "" {
			if d.inverted {
				sb.WriteString("[\\S\\s]")
			} else {
				sb.WriteString("(?!)")
			}
			return
		}
	}
	sb.WriteByte('[')
	if d.inverted {
		sb.WriteByte('^')
	}
	sb.WriteString(d.namedGroups)
	rangesStartCharset := !d.inverted && d.namedGroups == ""
	emitDashAtEnd := false
	for i, n := 0, d.ranges.GetNumRanges(); i < n; i++ {
		start := uint16(d.ranges.Start(i))
		end := uint16(d.ranges.End(i) - 1)
		switch int(end) - int(start) {
		case 0:
			if start == '-' {
				// Put it at the end where it doesn't need escaping.
				emitDashAtEnd = true
			} else {
				RegExpTree_DecomposedCharset_EscapeRangeCharOnto(start, rangesStartCharset, i == 0, i+1 == n, sb)
			}
		case 1:
			RegExpTree_DecomposedCharset_EscapeRangeCharOnto(start, rangesStartCharset, i == 0, false, sb)
			RegExpTree_DecomposedCharset_EscapeRangeCharOnto(end, rangesStartCharset, false, i+1 == n, sb)
		default:
			RegExpTree_DecomposedCharset_EscapeRangeCharOnto(start, rangesStartCharset, i == 0, false, sb)
			sb.WriteByte('-')
			RegExpTree_DecomposedCharset_EscapeRangeCharOnto(end, rangesStartCharset, false, true, sb)
		}
	}
	if emitDashAtEnd {
		sb.WriteByte('-')
	}
	sb.WriteByte(']')
}

// RegExpTree_DecomposedCharset_EscapeRangeCharOnto escapes a character that
// appears in a charset.
//
// Java: DecomposedCharset.escapeRangeCharOnto(char, boolean, boolean, boolean, StringBuilder)
func RegExpTree_DecomposedCharset_EscapeRangeCharOnto(
	ch uint16, startIsFlush bool, atStart bool, atEnd bool, sb *strings.Builder) {
	switch ch {
	case '\b':
		sb.WriteString("\\b")
	case '^':
		if atStart && startIsFlush {
			sb.WriteString("\\^")
		} else {
			sb.WriteByte('^')
		}
	case '-':
		if atStart || atEnd {
			sb.WriteByte('-')
		} else {
			sb.WriteString("\\-")
		}
	case '\\', ']':
		sb.WriteByte('\\')
		sb.WriteByte(byte(ch))
	default:
		RegExpTree_EscapeCharOnto(ch, sb)
	}
}

// RegExpTree_DecomposedCharset_Complexity estimates the rendered length of
// ranges.
//
// Java: DecomposedCharset.complexity(CharRanges)
func RegExpTree_DecomposedCharset_Complexity(ranges *CharRanges) int {
	complexity := 0
	for i, n := 0, ranges.GetNumRanges(); i < n; i++ {
		start := ranges.Start(i)
		end := ranges.End(i) - 1
		if start < 0x20 || start >= 0x7f {
			if start >= 0x100 {
				complexity += 6
			} else {
				complexity += 4
			}
		} else {
			complexity++
		}
		switch end - start {
		case 0:
			continue
		case 1:
		default:
			complexity += 1
		}
		if end < 0x20 || end >= 0x7f {
			if end >= 0x100 {
				complexity += 6
			} else {
				complexity += 4
			}
		} else {
			complexity++
		}
	}
	return complexity
}

// Equals is the record's equals: the same class and equal components
// (Objects.equals for the references).
//
// Java: DecomposedCharset.equals(Object)
func (d *RegExpTree_DecomposedCharset) Equals(o any) bool {
	that, ok := o.(*RegExpTree_DecomposedCharset)
	if !ok || that == nil {
		return false
	}
	rangesEqual := d.ranges == that.ranges ||
		(d.ranges != nil && that.ranges != nil && d.ranges.Equals(that.ranges))
	return d.inverted == that.inverted && rangesEqual && d.namedGroups == that.namedGroups
}

// HashCode is the record's hashCode: a combination of the components'
// hash codes.
//
// Go-only: language-adapter — java.lang.Record leaves the combination
// unspecified, and nothing hashes a DecomposedCharset, so this one only
// keeps equal records at equal hash codes: 31 * (31 * h(inverted) +
// h(ranges)) + h(namedGroups), with Boolean.hashCode's 1231/1237 and 0 for
// null.
//
// Java: DecomposedCharset.hashCode()
func (d *RegExpTree_DecomposedCharset) HashCode() int32 {
	h := int32(1237)
	if d.inverted {
		h = 1231
	}
	rangesHash := int32(0)
	if d.ranges != nil {
		rangesHash = d.ranges.HashCode()
	}
	h = 31*h + rangesHash
	return 31*h + jutil.StringHashCode(d.namedGroups)
}
