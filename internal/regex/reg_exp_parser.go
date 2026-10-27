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

// RegExpTree_ParentheticalType is the kind of a parenthesized group.
//
// Java: private enum RegExpTree.ParentheticalType
type RegExpTree_ParentheticalType int32

const (
	RegExpTree_ParentheticalType_CAPTURING RegExpTree_ParentheticalType = iota
	RegExpTree_ParentheticalType_NONCAPTURING
	RegExpTree_ParentheticalType_POSITIVE_LOOKAHEAD
	RegExpTree_ParentheticalType_NEGATIVE_LOOKAHEAD
	RegExpTree_ParentheticalType_POSITIVE_LOOKBEHIND
	RegExpTree_ParentheticalType_NEGATIVE_LOOKBEHIND
	RegExpTree_ParentheticalType_NAMED_GROUPS
)

var regExpTree_ParentheticalType_names = [...]string{
	"CAPTURING",
	"NONCAPTURING",
	"POSITIVE_LOOKAHEAD",
	"NEGATIVE_LOOKAHEAD",
	"POSITIVE_LOOKBEHIND",
	"NEGATIVE_LOOKBEHIND",
	"NAMED_GROUPS",
}

// String is the constant's name.
//
// Java: ParentheticalType.name()
func (p RegExpTree_ParentheticalType) String() string {
	return regExpTree_ParentheticalType_names[p]
}

// RegExpTree_ParseRegExp parses a regular expression to an AST.
//
// pattern is the "foo" from /foo/i; flags is the "i" from /foo/i.
//
// Java: RegExpTree.parseRegExp(String, String)
func RegExpTree_ParseRegExp(pattern string, flags string) RegExpTree {
	return NewRegExpTree_Parser(pattern, flags).ParseTopLevel()
}

// RegExpTree_Parser is a recursive descent parser that closes over pattern
// and flags.
//
// Java: local class Parser in RegExpTree.parseRegExp(String, String); the
// captured pattern and flags are constructor parameters.
type RegExpTree_Parser struct {
	pattern string
	flags   string

	// Go-only: language-adapter — the UTF-16 code units of pattern, so that
	// pattern.charAt(i) is O(1) as in Java.
	patternChars []uint16

	// The number of characters in pattern consumed.
	pos int

	// The number of capturing groups seen so far.
	numCapturingGroups int

	// The names of capturing groups in the regex expression
	capturingGroupNames map[string]bool

	// The length of pattern.
	limit int

	// Boolean indicating whether we should look for named capture group
	// backreferences
	lookForNamedCaptureBackreferences bool
}

// NewRegExpTree_Parser constructs the parser for pattern and flags.
//
// Java: new Parser() (inside parseRegExp)
func NewRegExpTree_Parser(pattern string, flags string) *RegExpTree_Parser {
	p := &RegExpTree_Parser{}
	p.pattern = pattern
	p.flags = flags
	p.patternChars = jutil.StringToUTF16(pattern)
	p.capturingGroupNames = make(map[string]bool)
	p.limit = len(p.patternChars)
	return p
}

// charAt is pattern.charAt(i).
//
// Go-only: language-adapter — String.charAt on the UTF-16 view, with Java's
// StringIndexOutOfBoundsException.
func (p *RegExpTree_Parser) charAt(i int) uint16 {
	if i < 0 || i >= len(p.patternChars) {
		panic(jutil.NewStringIndexOutOfBoundsExceptionForIndex(i, len(p.patternChars)))
	}
	return p.patternChars[i]
}

// substring is pattern.substring(begin, end).
//
// Go-only: language-adapter — String.substring on the UTF-16 view, with
// Java's StringIndexOutOfBoundsException.
func (p *RegExpTree_Parser) substring(begin, end int) string {
	if begin < 0 || end > len(p.patternChars) || begin > end {
		panic(jutil.NewStringIndexOutOfBoundsExceptionForRange(begin, end, len(p.patternChars)))
	}
	return jutil.StringFromUTF16(p.patternChars[begin:end])
}

// substringFrom is pattern.substring(begin).
//
// Go-only: language-adapter — see substring.
func (p *RegExpTree_Parser) substringFrom(begin int) string {
	return p.substring(begin, len(p.patternChars))
}

// ParseTopLevel parses the whole pattern.
//
// Java: Parser.parseTopLevel()
func (p *RegExpTree_Parser) ParseTopLevel() RegExpTree {
	// First assume there are no named capture backreferences,
	// because the spec says they should only be recognized
	// if the pattern contains at least one named capture group.
	p.pos = 0
	p.numCapturingGroups = 0
	p.lookForNamedCaptureBackreferences = false
	out := p.Parse()

	// If there were named capture groups, we must parse the pattern string again
	// so we can check for any backreferences to them.
	if len(p.capturingGroupNames) != 0 {
		p.pos = 0
		p.numCapturingGroups = 0
		p.lookForNamedCaptureBackreferences = true
		out = p.Parse()
	}

	if p.pos < p.limit { // Unmatched closed paren maybe.
		panic(jutil.NewIllegalArgumentExceptionMsg(p.substringFrom(p.pos)))
	}
	return out
}

// Parse parses a disjunction up to the end of the pattern or an unmatched
// close paren.
//
// Java: Parser.parse()
func (p *RegExpTree_Parser) Parse() RegExpTree {
	// Collects ["foo", "bar", "baz"] for /foo|bar|baz/.
	var alternatives []RegExpTree
	// The last item parsed within an alternation.
	var preceder RegExpTree

topLoop:
	for p.pos < p.limit {
		ch := p.charAt(p.pos)
		var atom RegExpTree
		switch ch {
		case '[':
			atom = p.ParseCharset()
		case '(':
			atom = p.ParseParenthetical()
		case ')':
			break topLoop
		case '\\':
			atom = p.ParseEscape()
		case '^', '$':
			atom = NewRegExpTree_Anchor(ch)
			p.pos++
		case '.':
			// We represent . as a character set to make it easy to simplify
			// things like /.|[\r\n]/.
			atom = RegExpTree_DOT_CHARSET
			p.pos++
		case '|':
			// An alternative may be empty as in /foo||bar/.
			// The '|' is consumed below.
			atom = RegExpTree_Empty_INSTANCE
		default:
			// Find a run of concatenated characters to avoid building a
			// tree node per literal character.
			start := p.pos
			end := p.pos + 1
		charsLoop:
			for end < p.limit {
				switch p.charAt(end) {
				case '[', '(', ')', '\\', '^', '$', '|', '.', '*', '+', '?', '{':
					break charsLoop
				default:
					// Repetition binds more tightly than concatenation.
					// Only consume up to "foo" in /foob*/ so that the suffix
					// operator parser below has the right precedence.
					if end+1 >= p.limit || !p.IsRepetitionStart(p.charAt(end+1)) {
						end++
					} else {
						break charsLoop
					}
				}
			}
			atom = NewRegExpTree_Text(p.substring(start, end))
			p.pos = end
		}
		if p.pos < p.limit && p.IsRepetitionStart(p.charAt(p.pos)) {
			atom = p.ParseRepetition(atom)
		}
		if preceder == nil {
			preceder = atom
		} else {
			preceder = NewRegExpTree_ConcatenationOfPair(preceder, atom)
		}
		// If this is an alternative in a alternation, then add it to the
		// list of complete alternatives, and reset the parser state for the
		// next alternative.
		if p.pos < p.limit && p.charAt(p.pos) == '|' {
			if alternatives == nil {
				alternatives = []RegExpTree{}
			}
			alternatives = append(alternatives, preceder)
			preceder = nil
			p.pos++
		}
	}
	// An alternative may have no parsed content blank as in /foo|/.
	if preceder == nil {
		preceder = RegExpTree_Empty_INSTANCE
	}
	if alternatives != nil {
		alternatives = append(alternatives, preceder)
		return NewRegExpTree_Alternation(alternatives)
	} else {
		return preceder
	}
}

// ParseParenthetical handles capturing groups (...), non-capturing groups
// (?:...), and lookahead assertions (?=...).
//
// Java: Parser.parseParenthetical()
func (p *RegExpTree_Parser) ParseParenthetical() RegExpTree {
	jutil.CheckState(p.charAt(p.pos) == '(')
	start := p.pos
	p.pos++
	var type_ RegExpTree_ParentheticalType
	var captureName *string
	if p.pos < p.limit && p.charAt(p.pos) == '?' {
		if p.pos+1 < p.limit {
			ch := p.charAt(p.pos + 1)
			switch ch {
			case ':':
				// (?:...) Non-capturing groups.
				p.pos += 2
				type_ = RegExpTree_ParentheticalType_NONCAPTURING
			// (?=...) and (?!...) Lookahead Assertions
			case '=':
				p.pos += 2
				type_ = RegExpTree_ParentheticalType_POSITIVE_LOOKAHEAD
			case '!':
				p.pos += 2
				type_ = RegExpTree_ParentheticalType_NEGATIVE_LOOKAHEAD
			// (?<=...) and (?<!...) Lookbehind Assertions, (?<name>) named groups
			case '<':
				if p.pos+2 < p.limit && p.charAt(p.pos+2) == '=' {
					p.pos += 3
					type_ = RegExpTree_ParentheticalType_POSITIVE_LOOKBEHIND
				} else if p.pos+2 < p.limit && p.charAt(p.pos+2) == '!' {
					p.pos += 3
					type_ = RegExpTree_ParentheticalType_NEGATIVE_LOOKBEHIND
				} else {
					p.pos += 2
					name := p.ScanNamedGroupName()
					captureName = &name
					p.capturingGroupNames[*captureName] = true
					type_ = RegExpTree_ParentheticalType_NAMED_GROUPS
				}
			default:
				panic(jutil.NewIllegalArgumentExceptionMsg(
					"Malformed parenthetical: " + p.substringFrom(start)))
			}
		} else {
			panic(jutil.NewIllegalArgumentExceptionMsg(
				"Malformed parenthetical: " + p.substringFrom(start)))
		}
	} else {
		type_ = RegExpTree_ParentheticalType_CAPTURING
	}
	body := p.Parse()
	if p.pos < p.limit && p.charAt(p.pos) == ')' {
		p.pos++
	} else {
		panic(jutil.NewIllegalArgumentExceptionMsg(
			"Unclosed parenthetical group: " + p.substringFrom(start)))
	}
	switch type_ {
	case RegExpTree_ParentheticalType_CAPTURING:
		p.numCapturingGroups++
		return NewRegExpTree_CapturingGroup(body)
	case RegExpTree_ParentheticalType_NONCAPTURING:
		return body
	case RegExpTree_ParentheticalType_POSITIVE_LOOKAHEAD:
		return NewRegExpTree_LookaheadAssertion(body, true)
	case RegExpTree_ParentheticalType_NEGATIVE_LOOKAHEAD:
		return NewRegExpTree_LookaheadAssertion(body, false)
	case RegExpTree_ParentheticalType_POSITIVE_LOOKBEHIND:
		return NewRegExpTree_LookbehindAssertion(body, true)
	case RegExpTree_ParentheticalType_NEGATIVE_LOOKBEHIND:
		return NewRegExpTree_LookbehindAssertion(body, false)
	case RegExpTree_ParentheticalType_NAMED_GROUPS:
		if captureName != nil {
			p.numCapturingGroups++
			return NewRegExpTree_NamedCaptureGroup(body, *captureName)
		} else {
			panic(jutil.NewIllegalArgumentExceptionMsg(
				"Malformed named capture group: " + p.substringFrom(start)))
		}
	}
	panic(jutil.NewAssertionError("Unrecognized ParentheticalType " + type_.String()))
}

// ScanNamedGroupName scans the pattern for a named group name. Assumes that
// pos points to the character after '<'. Returns the group name.
//
// Java: Parser.scanNamedGroupName()
func (p *RegExpTree_Parser) ScanNamedGroupName() string {
	start := p.pos
	end := p.pos
	if !RegExpTree_IsIdentifierStart(p.charAt(start)) {
		panic(jutil.NewIllegalArgumentExceptionMsg(
			"Invalid capture group name: <" + p.substringFrom(start)))
	}
	end++
	for end < p.limit {
		if p.charAt(end) == '>' {
			p.pos = end + 1
			return p.substring(start, end)
		} else if RegExpTree_IsIdentifierPart(p.charAt(end)) {
			end++
		} else {
			panic(jutil.NewIllegalArgumentExceptionMsg(
				"Invalid capture group name: <" + p.substringFrom(start)))
		}
	}
	panic(jutil.NewIllegalArgumentExceptionMsg(
		"Malformed named capture group: <" + p.substringFrom(start)))
}

// ParseCharset parses a square bracketed character set. Standalone character
// groups /\d/ are handled by ParseEscape.
//
// Java: Parser.parseCharset()
func (p *RegExpTree_Parser) ParseCharset() RegExpTree {
	jutil.CheckState(p.charAt(p.pos) == '[')
	p.pos++

	isCaseInsensitive := strings.IndexByte(p.flags, 'i') >= 0
	inverse := p.pos < p.limit && p.charAt(p.pos) == '^'
	if inverse {
		p.pos++
	}
	ranges := CharRanges_EMPTY
	ieExplicits := CharRanges_EMPTY
	for p.pos < p.limit && p.charAt(p.pos) != ']' {
		ch := p.charAt(p.pos)
		var start int
		if ch == '\\' {
			p.pos++
			possibleGroupName := p.charAt(p.pos)
			group := regExpTree_namedCharGroup(possibleGroupName)
			if group != nil {
				p.pos++
				ranges = ranges.Union(group)
				continue
			}
			start = p.ParseEscapeChar()
		} else {
			start = int(ch)
			p.pos++
		}
		end := start
		if p.pos+1 < p.limit && p.charAt(p.pos) == '-' && p.charAt(p.pos+1) != ']' {
			p.pos++
			ch = p.charAt(p.pos)
			if ch == '\\' {
				p.pos++
				end = p.ParseEscapeChar()
			} else {
				end = int(ch)
				p.pos++
			}
		}
		range_ := CharRanges_Inclusive(start, end)
		ranges = ranges.Union(range_)
		if RegExpTree_IE_SPEC_ERRORS.Contains(start) && RegExpTree_IE_SPEC_ERRORS.Contains(end) {
			ieExplicits = ieExplicits.Union(range_.Intersection(RegExpTree_IE_SPEC_ERRORS))
		}
		if isCaseInsensitive {
			// If the flags contain the 'i' flag, then it is not correct to
			// say that [^a-z] contains the letter 'A', or that [a-z] does not
			// contain the letter 'A'.
			// We expand out letter groups here so that parse returns something
			// that is valid independent of flags.
			// Calls to simplify(flags) may later reintroduce flag assumptions.
			// but without this step, later steps might conflate
			//     /[a-z]/i
			// and
			//     /[^\0-`{-\uffff]/i
			// which matches nothing because the information about whether the
			// ^ is present has been lost during optimizations and charset
			// unionizing as in /[...]|[^...]/.
			ranges = CaseCanonicalize_ExpandToAllMatched(ranges)
		}
	}
	p.pos++ // Consume ']'

	if inverse {
		ranges = CharRanges_ALL_CODE_UNITS.Difference(ranges)
	}

	return NewRegExpTree_Charset(ranges, ieExplicits)
}

// ParseEscapeChar parses an escape to a code point. Some of the characters
// parsed here have special meanings in various contexts, so contexts must
// filter those instead. E.g. '\b' means a different thing inside a charset
// than without.
//
// Java: Parser.parseEscapeChar()
func (p *RegExpTree_Parser) ParseEscapeChar() int {
	p.pos++
	ch := p.charAt(p.pos - 1)
	switch ch {
	case 'b':
		return '\b'
	case 'f':
		return '\f'
	case 'n':
		return '\n'
	case 'r':
		return '\r'
	case 't':
		return '\t'
	case 'u':
		if strings.Contains(p.flags, "u") && p.pos < p.limit && p.charAt(p.pos) == '{' {
			return p.ParseBracedUnicodeEscape()
		}
		return p.ParseHex(4)
	case 'v':
		return '\u000b'
	case 'x':
		return p.ParseHex(2)
	default:
		if '0' <= ch && ch <= '7' {
			codeUnit := ch - '0'
			// Allow octal literals in the range \0-\377.
			// \41 might be a group, but \041 is not a group.
			// We read, but do not emit octal literals since they
			// are deprecated in ES5.
			extra := 1
			if ch <= '3' {
				extra = 2
			}
			if ch == '0' {
				extra++
			}
			octLimit := min(p.limit, p.pos+extra)
			for p.pos < octLimit {
				ch = p.charAt(p.pos)
				if '0' <= ch && ch <= '7' {
					codeUnit = (codeUnit << 3) + (ch - '0')
					p.pos++
				} else {
					break
				}
			}
			return int(codeUnit)
		}
		return int(ch)
	}
}

// ParseEscape parses an escape that appears outside a charset.
//
// Java: Parser.parseEscape()
func (p *RegExpTree_Parser) ParseEscape() RegExpTree {
	jutil.CheckState(p.charAt(p.pos) == '\\')
	start := p.pos
	p.pos++
	ch := p.charAt(p.pos)
	if ch == 'b' || ch == 'B' {
		p.pos++
		return NewRegExpTree_WordBoundary(ch)
	} else if (ch == 'p' || ch == 'P') && strings.Contains(p.flags, "u") {
		// handle ES2018 unicode property tests, e.g.
		// /\p{ASCII_Hex_Digit=true}/ only hex digits
		// /\P{Script=Greek}/  no greek letters
		negated := ch == 'P'
		p.pos++
		if p.pos < p.limit && p.charAt(p.pos) == '{' {
			lhs := jutil.NewStringBuilder()
			for {
				p.pos++
				if !(p.pos < p.limit) {
					break
				}
				ch = p.charAt(p.pos)
				if !(ch == '_' ||
					('a' <= ch && ch <= 'z') ||
					('A' <= ch && ch <= 'Z') ||
					('0' <= ch && ch <= '9')) {
					break
				}
				lhs.AppendChar(ch)
			}
			if p.pos < p.limit && ch == '}' {
				// Case of shorthand like /\p{ASCII_Hex_Digit}/u
				p.pos++
				return NewRegExpTree_UnicodePropertyEscape(nil, lhs.String(), negated)
			} else if p.pos < p.limit && ch == '=' {
				// Case of having '=' like /\p{Script=Greek}/u
				rhs := jutil.NewStringBuilder()
				for {
					p.pos++
					if !(p.pos < p.limit) {
						break
					}
					ch = p.charAt(p.pos)
					if !(ch == '_' ||
						('a' <= ch && ch <= 'z') ||
						('A' <= ch && ch <= 'Z') ||
						('0' <= ch && ch <= '9')) {
						break
					}
					rhs.AppendChar(ch)
				}
				if p.pos < p.limit && ch == '}' {
					p.pos++
					lhsString := lhs.String()
					return NewRegExpTree_UnicodePropertyEscape(&lhsString, rhs.String(), negated)
				} else {
					panic(jutil.NewIllegalArgumentExceptionMsg(
						"Malformed Unicode Property Escape: expected '}' after " +
							p.substring(start, p.pos)))
				}
			} else {
				panic(jutil.NewIllegalArgumentExceptionMsg(
					"Malformed Unicode Property Escape: expected '=' or '}' after " +
						p.substring(start, p.pos)))
			}
		} else {
			panic(jutil.NewIllegalArgumentExceptionMsg(
				"Malformed Unicode Property Escape: expected '{' after " +
					p.substring(start, p.pos)))
		}
	} else if '1' <= ch && ch <= '9' {
		p.pos++
		possibleGroupIndex := int(ch - '0')
		if p.numCapturingGroups >= possibleGroupIndex {
			if p.pos < p.limit {
				next := p.charAt(p.pos)
				if '0' <= next && next <= '9' {
					twoDigitGroupIndex := possibleGroupIndex*10 + int(next-'0')
					if p.numCapturingGroups >= twoDigitGroupIndex {
						p.pos++
						possibleGroupIndex = twoDigitGroupIndex
					}
				}
			}
			return NewRegExpTree_BackReference(possibleGroupIndex)
		} else {
			// \1 - \7 are octal escapes if there is no such group.
			// \8 and \9 are the literal characters '8' and '9' if there
			// is no such group.
			c := ch
			if possibleGroupIndex <= 7 {
				c = uint16(possibleGroupIndex)
			}
			return NewRegExpTree_Text(jutil.StringFromChar(c))
		}
	} else if p.lookForNamedCaptureBackreferences &&
		ch == 'k' &&
		p.pos+1 < p.limit &&
		p.charAt(p.pos+1) == '<' &&
		// According to the spec
		// https://github.com/tc39/proposal-regexp-named-groups#backwards-compatibility-of-new-syntax
		// we want to treat \k as a normal string if there are no named capturing groups present
		len(p.capturingGroupNames) != 0 {
		p.pos += 2
		potentialName := p.ScanNamedGroupName()
		if !p.capturingGroupNames[potentialName] {
			panic(jutil.NewIllegalArgumentExceptionMsg(
				"Invalid named capture referenced: " + p.substringFrom(start)))
		}
		return NewRegExpTree_NamedBackReference(potentialName)
	} else {
		charGroup := regExpTree_namedCharGroup(ch)
		if charGroup != nil { // Handle \d, etc.
			p.pos++
			return NewRegExpTree_Charset(charGroup, CharRanges_EMPTY)
		}
		return NewRegExpTree_Text(jutil.StringFromUTF16(jutil.CharacterToChars(p.ParseEscapeChar())))
	}
}

// ParseHex parses n hex digits to a code-unit.
//
// Java: Parser.parseHex(int)
func (p *RegExpTree_Parser) ParseHex(n int) int {
	if p.pos+n > p.limit {
		panic(jutil.NewIllegalArgumentExceptionMsg("Abbreviated hex escape " + p.substringFrom(p.pos)))
	}
	if n > 7 {
		// We need to guard the MSB to prevent overflow.
		panic(jutil.NewIllegalArgumentExceptionMsg(
			"Cannot parse hexadecimal encoding wider than 28 bits: " +
				p.substring(p.pos, p.pos+n)))
	}

	result := 0
	for n--; n >= 0; n-- {
		ch := p.charAt(p.pos)
		var digit int
		if '0' <= ch && ch <= '9' {
			digit = int(ch - '0')
		} else if 'a' <= ch && ch <= 'f' {
			digit = int(ch) + (10 - 'a')
		} else if 'A' <= ch && ch <= 'F' {
			digit = int(ch) + (10 - 'A')
		} else {
			panic(jutil.NewIllegalArgumentExceptionMsg(p.substringFrom(p.pos)))
		}
		p.pos++
		result = (result << 4) | digit
	}
	return result
}

// ParseBracedUnicodeEscape parses \u{...} (with the 'u' flag).
//
// Java: Parser.parseBracedUnicodeEscape()
func (p *RegExpTree_Parser) ParseBracedUnicodeEscape() int {
	openBrace := p.pos
	p.pos++
	jutil.CheckState(p.charAt(p.pos-1) == '{')

	closeBrace := p.pos
	for closeBrace < p.limit && p.charAt(closeBrace) != '}' {
		closeBrace++
	}
	if closeBrace == p.limit {
		panic(jutil.NewIllegalArgumentExceptionMsg(
			"Malformed unicode escape: expected '}' after " + p.substringFrom(openBrace)))
	} else if closeBrace == p.pos {
		panic(jutil.NewIllegalArgumentExceptionMsg("Empty unicode escape"))
	}

	result := p.ParseHex(closeBrace - p.pos)
	if result > 0x10FFFF {
		panic(jutil.NewIllegalArgumentExceptionMsg(
			"Unicode must be at most 0x10FFFF: " + p.substring(openBrace+1, p.pos)))
	}
	p.pos++ // Consume the close brace.
	return result
}

// IsRepetitionStart reports whether ch starts a quantifier.
//
// Java: Parser.isRepetitionStart(char)
func (p *RegExpTree_Parser) IsRepetitionStart(ch uint16) bool {
	switch ch {
	case '?', '*', '+', '{':
		return true
	default:
		return false
	}
}

// ParseRepetition parses a repetition. x? is treated as a repetition -- an
// optional production can be matched 0 or 1 time.
//
// Java: Parser.parseRepetition(RegExpTree)
func (p *RegExpTree_Parser) ParseRepetition(body RegExpTree) RegExpTree {
	if p.pos == p.limit {
		return body
	}
	var min_ int
	var max_ int
	switch p.charAt(p.pos) {
	case '+':
		p.pos++
		min_ = 1
		max_ = javaIntegerMaxValue
	case '*':
		p.pos++
		min_ = 0
		max_ = javaIntegerMaxValue
	case '?':
		p.pos++
		min_ = 0
		max_ = 1
	case '{':
		p.pos++
		start := p.pos
		end := -1 // Java: pattern.indexOf('}', start)
		for i := start; i < p.limit; i++ {
			if p.patternChars[i] == '}' {
				end = i
				break
			}
		}
		if end < 0 {
			p.pos = start - 1
			return body
		}
		counts := p.substring(start, end)
		p.pos = end + 1
		comma := jutil.StringIndexOf(counts, ",")
		jutil.Try(func() {
			if comma >= 0 {
				min_ = int(jutil.IntegerParseInt(jutil.StringSubstring(counts, 0, comma), 10))
			} else {
				min_ = int(jutil.IntegerParseInt(counts, 10))
			}
			if comma >= 0 {
				if comma+1 != jutil.StringLength(counts) {
					max_ = int(jutil.IntegerParseInt(jutil.StringSubstringFrom(counts, comma+1), 10))
				} else {
					max_ = javaIntegerMaxValue
				}
			} else {
				max_ = min_
			}
		}, func(ex jutil.INumberFormatException) {
			min_ = -1
			max_ = -1
		})
		if min_ < 0 || min_ > max_ {
			// Treat the open curly bracket literally.
			p.pos = start - 1
			return body
		}
	default:
		return body
	}
	greedy := true
	if p.pos < p.limit && p.charAt(p.pos) == '?' {
		greedy = false
		p.pos++
	}
	return NewRegExpTree_Repetition(body, min_, max_, greedy)
}

// RegExpTree_IsIdentifierStart reports whether ch is a valid Javascript
// identifier start character.
//
// Java: RegExpTree.isIdentifierStart(char)
func RegExpTree_IsIdentifierStart(ch uint16) bool {
	// TODO(yitingwang) This and the one in Scanner.java should share the same implementation
	// Most code is written in pure ASCII create a fast path here.
	if ch <= 127 {
		return (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || (ch == '_' || ch == '$')
	}

	return jutil.CharacterIsLetter(ch)
}

// RegExpTree_IsIdentifierPart reports whether ch is allowed in Javascript
// identifiers.
//
// Java: RegExpTree.isIdentifierPart(char)
func RegExpTree_IsIdentifierPart(ch uint16) bool {
	// TODO(yitingwang) This and the one in Scanner.java should share the same implementation
	// Most code is written in pure ASCII create a fast path here.
	if ch <= 127 {
		return (ch >= 'A' && ch <= 'Z') ||
			(ch >= 'a' && ch <= 'z') ||
			(ch >= '0' && ch <= '9') ||
			(ch == '_' || ch == '$') // _ or $
	}
	// TODO: identifier part character classes
	// CombiningMark
	//   Non-Spacing mark (Mn)
	//   Combining spacing mark(Mc)
	// Connector punctuation (Pc)
	// Zero Width Non-Joiner
	// Zero Width Joiner
	return RegExpTree_IsIdentifierStart(ch) || jutil.CharacterIsDigit(ch)
}
