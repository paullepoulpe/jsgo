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
	"math"
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// javaIntegerMaxValue is Integer.MAX_VALUE.
const javaIntegerMaxValue = math.MaxInt32

// RegExpTree is an AST for JavaScript regular expressions.
//
// Java: abstract class RegExpTree. Its subclasses embed RegExpTreeBase
// (PORTING.md §4).
type RegExpTree interface {
	// Simplify returns a simpler regular expression that is semantically the
	// same assuming the given flags (e.g. "igm").
	Simplify(flags string) RegExpTree
	// IsCaseSensitive reports whether the presence or absence of an "i" flag
	// would change the meaning of this regular expression.
	IsCaseSensitive() bool
	// ContainsAnchor reports whether the regular expression contains an
	// anchor: ^ or $.
	ContainsAnchor() bool
	// HasCapturingGroup reports whether the regular expression contains
	// capturing groups.
	HasCapturingGroup() bool
	// NumCapturingGroups is the number of capturing groups.
	NumCapturingGroups() int
	// Children returns the children of this node.
	Children() []RegExpTree
	// AppendSourceCode appends this regular expression source to the given
	// buffer.
	AppendSourceCode(sb *jutil.StringBuilder)
	AppendDebugInfo(sb *jutil.StringBuilder)
	String() string
	Equals(o any) bool
	HashCode() int32
	GetRegExpTreeBase() *RegExpTreeBase
}

// RegExpTreeBase holds the base implementations of RegExpTree.
type RegExpTreeBase struct {
	self RegExpTree
}

// InitRegExpTreeBase is the Go form of the RegExpTree constructor for
// subclasses.
//
// Java: RegExpTree()
func (t *RegExpTreeBase) InitRegExpTreeBase(self RegExpTree) {
	t.self = self
}

// GetRegExpTreeBase returns the base struct (field access through the
// interface).
//
// Go-only: language-adapter — field access through the interface
// (PORTING.md §4).
func (t *RegExpTreeBase) GetRegExpTreeBase() *RegExpTreeBase {
	return t
}

// HasCapturingGroup reports whether the regular expression contains
// capturing groups.
//
// Java: RegExpTree.hasCapturingGroup()
func (t *RegExpTreeBase) HasCapturingGroup() bool {
	return t.self.NumCapturingGroups() != 0
}

// String renders the regular expression as a literal, e.g. "/foo/".
//
// Java: RegExpTree.toString()
func (t *RegExpTreeBase) String() string {
	sb := jutil.NewStringBuilder()
	sb.AppendChar('/')
	t.self.AppendSourceCode(sb)
	// Don't emit a regular expression that looks like a line comment start.
	if sb.Length() == 1 {
		sb.Append("(?:)")
	}
	sb.AppendChar('/')
	return sb.String()
}

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
	capturingGroupNames *jutil.LinkedHashSet[string]

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
	p.capturingGroupNames = jutil.NewLinkedHashSet[string]()
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
	if !p.capturingGroupNames.IsEmpty() {
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
					p.capturingGroupNames.Add(*captureName)
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
			group := RegExpTree_NAMED_CHAR_GROUPS.Get(possibleGroupName)
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
		!p.capturingGroupNames.IsEmpty() {
		p.pos += 2
		potentialName := p.ScanNamedGroupName()
		if !p.capturingGroupNames.Contains(potentialName) {
			panic(jutil.NewIllegalArgumentExceptionMsg(
				"Invalid named capture referenced: " + p.substringFrom(start)))
		}
		return NewRegExpTree_NamedBackReference(potentialName)
	} else {
		charGroup := RegExpTree_NAMED_CHAR_GROUPS.Get(ch)
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

// RegExpTree_MatchesWholeInput reports true if, but not necessarily always
// when, the given regular expression must match the whole input or none of
// it.
//
// Java: RegExpTree.matchesWholeInput(RegExpTree, String)
func RegExpTree_MatchesWholeInput(t RegExpTree, flags string) bool {
	if strings.IndexByte(flags, 'm') >= 0 {
		return false
	}

	c, ok := t.(*RegExpTree_Concatenation)
	if !ok {
		return false
	}

	if len(c.elements) == 0 {
		return false
	}
	first := c.elements[0]
	last := c.elements[len(c.elements)-1]
	firstAnchor, ok1 := first.(*RegExpTree_Anchor)
	lastAnchor, ok2 := last.(*RegExpTree_Anchor)
	if !(ok1 && ok2) {
		return false
	}
	return firstAnchor.type_ == '^' && lastAnchor.type_ == '$'
}

// RegExpTree_RegExpTreeAtom represents a node that never has children such
// as an anchor or charset.
//
// Java: abstract static class RegExpTree.RegExpTreeAtom
type RegExpTree_RegExpTreeAtom interface {
	RegExpTree
	GetRegExpTree_RegExpTreeAtomBase() *RegExpTree_RegExpTreeAtomBase
}

// RegExpTree_RegExpTreeAtomBase holds RegExpTreeAtom's base implementations.
type RegExpTree_RegExpTreeAtomBase struct {
	RegExpTreeBase
}

// InitRegExpTree_RegExpTreeAtomBase is the Go form of the RegExpTreeAtom
// constructor for subclasses.
//
// Java: RegExpTreeAtom()
func (a *RegExpTree_RegExpTreeAtomBase) InitRegExpTree_RegExpTreeAtomBase(self RegExpTree) {
	a.InitRegExpTreeBase(self)
}

// GetRegExpTree_RegExpTreeAtomBase returns the base struct.
//
// Go-only: language-adapter — field access through the interface
// (PORTING.md §4).
func (a *RegExpTree_RegExpTreeAtomBase) GetRegExpTree_RegExpTreeAtomBase() *RegExpTree_RegExpTreeAtomBase {
	return a
}

// IsCaseSensitive is false for atoms unless overridden.
//
// Java: RegExpTreeAtom.isCaseSensitive()
func (a *RegExpTree_RegExpTreeAtomBase) IsCaseSensitive() bool {
	return false
}

// ContainsAnchor is false for atoms unless overridden.
//
// Java: RegExpTreeAtom.containsAnchor()
func (a *RegExpTree_RegExpTreeAtomBase) ContainsAnchor() bool {
	return false
}

// NumCapturingGroups is 0 for atoms.
//
// Java: RegExpTreeAtom.numCapturingGroups()
func (a *RegExpTree_RegExpTreeAtomBase) NumCapturingGroups() int {
	return 0
}

// Children is empty for atoms.
//
// Java: RegExpTreeAtom.children()
func (a *RegExpTree_RegExpTreeAtomBase) Children() []RegExpTree {
	return []RegExpTree{}
}

// RegExpTree_Empty represents an empty portion of a RegExp such as the middle
// of "||".
//
// Java: private static final class RegExpTree.Empty
type RegExpTree_Empty struct {
	RegExpTree_RegExpTreeAtomBase
}

// NewRegExpTree_Empty constructs an Empty.
//
// Java: Empty()
func NewRegExpTree_Empty() *RegExpTree_Empty {
	e := &RegExpTree_Empty{}
	e.InitRegExpTree_RegExpTreeAtomBase(e)
	return e
}

// RegExpTree_Empty_INSTANCE is the shared Empty.
var RegExpTree_Empty_INSTANCE = NewRegExpTree_Empty()

// Simplify returns e.
//
// Java: Empty.simplify(String)
func (e *RegExpTree_Empty) Simplify(flags string) RegExpTree {
	return e
}

// AppendSourceCode emits nothing.
//
// Java: Empty.appendSourceCode(StringBuilder)
func (e *RegExpTree_Empty) AppendSourceCode(sb *jutil.StringBuilder) {
	// No output
}

// AppendDebugInfo emits nothing.
//
// Java: Empty.appendDebugInfo(StringBuilder)
func (e *RegExpTree_Empty) AppendDebugInfo(sb *jutil.StringBuilder) {
	// No output
}

// Equals reports whether o is an Empty.
//
// Java: Empty.equals(Object)
func (e *RegExpTree_Empty) Equals(o any) bool {
	_, ok := o.(*RegExpTree_Empty)
	return ok
}

// HashCode is a constant.
//
// Java: Empty.hashCode()
func (e *RegExpTree_Empty) HashCode() int32 {
	return 0x7ee06141
}

// RegExpTree_Anchor represents an anchor, namely ^ or $.
//
// Java: public static final class RegExpTree.Anchor
type RegExpTree_Anchor struct {
	RegExpTree_RegExpTreeAtomBase
	type_ uint16
}

// NewRegExpTree_Anchor constructs an Anchor.
//
// Java: Anchor(char)
func NewRegExpTree_Anchor(type_ uint16) *RegExpTree_Anchor {
	a := &RegExpTree_Anchor{}
	a.InitRegExpTree_RegExpTreeAtomBase(a)
	a.type_ = type_
	return a
}

// Simplify returns a.
//
// Java: Anchor.simplify(String)
func (a *RegExpTree_Anchor) Simplify(flags string) RegExpTree {
	return a
}

// ContainsAnchor is true.
//
// Java: Anchor.containsAnchor()
func (a *RegExpTree_Anchor) ContainsAnchor() bool {
	return true
}

// AppendSourceCode emits ^ or $.
//
// Java: Anchor.appendSourceCode(StringBuilder)
func (a *RegExpTree_Anchor) AppendSourceCode(sb *jutil.StringBuilder) {
	sb.AppendChar(a.type_)
}

// AppendDebugInfo emits the anchor character.
//
// Java: Anchor.appendDebugInfo(StringBuilder)
func (a *RegExpTree_Anchor) AppendDebugInfo(sb *jutil.StringBuilder) {
	sb.AppendChar(a.type_)
}

// Equals compares anchor types.
//
// Java: Anchor.equals(Object)
func (a *RegExpTree_Anchor) Equals(o any) bool {
	anchor, ok := o.(*RegExpTree_Anchor)
	return ok && a.type_ == anchor.type_
}

// HashCode hashes the anchor type.
//
// Java: Anchor.hashCode()
func (a *RegExpTree_Anchor) HashCode() int32 {
	return int32(uint32(a.type_) ^ 0xe85317ff)
}

// RegExpTree_WordBoundary represents \b or \B.
//
// Java: public static final class RegExpTree.WordBoundary
type RegExpTree_WordBoundary struct {
	RegExpTree_RegExpTreeAtomBase
	type_ uint16
}

// NewRegExpTree_WordBoundary constructs a WordBoundary.
//
// Java: WordBoundary(char)
func NewRegExpTree_WordBoundary(type_ uint16) *RegExpTree_WordBoundary {
	w := &RegExpTree_WordBoundary{}
	w.InitRegExpTree_RegExpTreeAtomBase(w)
	w.type_ = type_
	return w
}

// Simplify returns w.
//
// Java: WordBoundary.simplify(String)
func (w *RegExpTree_WordBoundary) Simplify(flags string) RegExpTree {
	return w
}

// AppendSourceCode emits \b or \B.
//
// Java: WordBoundary.appendSourceCode(StringBuilder)
func (w *RegExpTree_WordBoundary) AppendSourceCode(sb *jutil.StringBuilder) {
	sb.AppendChar('\\').AppendChar(w.type_)
}

// AppendDebugInfo emits the boundary type.
//
// Java: WordBoundary.appendDebugInfo(StringBuilder)
func (w *RegExpTree_WordBoundary) AppendDebugInfo(sb *jutil.StringBuilder) {
	sb.AppendChar(w.type_)
}

// Equals compares boundary types.
//
// Java: WordBoundary.equals(Object)
func (w *RegExpTree_WordBoundary) Equals(o any) bool {
	wordBoundary, ok := o.(*RegExpTree_WordBoundary)
	return ok && w.type_ == wordBoundary.type_
}

// HashCode hashes the boundary type.
//
// Java: WordBoundary.hashCode()
func (w *RegExpTree_WordBoundary) HashCode() int32 {
	return 0x5673aa29 ^ int32(w.type_)
}

// RegExpTree_BackReference represents a reference to a previous group such
// as \1 or \2.
//
// Java: public static final class RegExpTree.BackReference
type RegExpTree_BackReference struct {
	RegExpTree_RegExpTreeAtomBase
	groupIndex int
}

// NewRegExpTree_BackReference constructs a BackReference.
//
// Java: BackReference(int)
func NewRegExpTree_BackReference(groupIndex int) *RegExpTree_BackReference {
	b := &RegExpTree_BackReference{}
	b.InitRegExpTree_RegExpTreeAtomBase(b)
	jutil.CheckArgument(groupIndex >= 0 && groupIndex <= 99)
	b.groupIndex = groupIndex
	return b
}

// Simplify returns b.
//
// Java: BackReference.simplify(String)
func (b *RegExpTree_BackReference) Simplify(flags string) RegExpTree {
	return b
}

// AppendSourceCode emits \N.
//
// Java: BackReference.appendSourceCode(StringBuilder)
func (b *RegExpTree_BackReference) AppendSourceCode(sb *jutil.StringBuilder) {
	sb.AppendChar('\\').AppendInt(int64(b.groupIndex))
}

// AppendDebugInfo emits the group index.
//
// Java: BackReference.appendDebugInfo(StringBuilder)
func (b *RegExpTree_BackReference) AppendDebugInfo(sb *jutil.StringBuilder) {
	sb.AppendInt(int64(b.groupIndex))
}

// Equals compares group indices.
//
// Java: BackReference.equals(Object)
func (b *RegExpTree_BackReference) Equals(o any) bool {
	backReference, ok := o.(*RegExpTree_BackReference)
	return ok && b.groupIndex == backReference.groupIndex
}

// HashCode hashes the group index.
//
// Java: BackReference.hashCode()
func (b *RegExpTree_BackReference) HashCode() int32 {
	return int32(uint32(0xff072663) ^ uint32(int32(b.groupIndex)))
}

// RegExpTree_NamedBackReference represents a reference to a previous named
// group.
//
// Java: public static final class RegExpTree.NamedBackReference
type RegExpTree_NamedBackReference struct {
	RegExpTree_RegExpTreeAtomBase
	groupName string
}

// NewRegExpTree_NamedBackReference constructs a NamedBackReference.
//
// Java: NamedBackReference(String)
func NewRegExpTree_NamedBackReference(groupName string) *RegExpTree_NamedBackReference {
	n := &RegExpTree_NamedBackReference{}
	n.InitRegExpTree_RegExpTreeAtomBase(n)
	n.groupName = groupName
	return n
}

// Simplify returns n.
//
// Java: NamedBackReference.simplify(String)
func (n *RegExpTree_NamedBackReference) Simplify(flags string) RegExpTree {
	return n
}

// AppendSourceCode emits \k<name>.
//
// Java: NamedBackReference.appendSourceCode(StringBuilder)
func (n *RegExpTree_NamedBackReference) AppendSourceCode(sb *jutil.StringBuilder) {
	sb.Append("\\k<").Append(n.groupName).AppendChar('>')
}

// AppendDebugInfo emits the group name.
//
// Java: NamedBackReference.appendDebugInfo(StringBuilder)
func (n *RegExpTree_NamedBackReference) AppendDebugInfo(sb *jutil.StringBuilder) {
	sb.Append(n.groupName)
}

// Equals compares group names.
//
// Java: NamedBackReference.equals(Object)
func (n *RegExpTree_NamedBackReference) Equals(o any) bool {
	namedBackReference, ok := o.(*RegExpTree_NamedBackReference)
	return ok && n.groupName == namedBackReference.groupName
}

// HashCode hashes the group name.
//
// Java: NamedBackReference.hashCode()
func (n *RegExpTree_NamedBackReference) HashCode() int32 {
	return jutil.StringHashCode(n.groupName) // Java: Objects.hashCode(groupName)
}

// RegExpTree_Text represents a run of non-special characters such as
// "foobar".
//
// Java: public static final class RegExpTree.Text
type RegExpTree_Text struct {
	RegExpTree_RegExpTreeAtomBase
	text string
}

// NewRegExpTree_Text constructs a Text.
//
// Java: Text(String)
func NewRegExpTree_Text(text string) *RegExpTree_Text {
	t := &RegExpTree_Text{}
	t.InitRegExpTree_RegExpTreeAtomBase(t)
	t.text = text
	return t
}

// RegExpTree_Text_EscapeRegularCharOnto escapes ch onto sb. ch is the
// code-unit to escape; next is the next code-unit or -1 if indeterminable.
//
// Java: Text.escapeRegularCharOnto(char, int, StringBuilder)
func RegExpTree_Text_EscapeRegularCharOnto(ch uint16, next int, sb *jutil.StringBuilder) {
	switch ch {
	case '$', '^', '*', '(', ')', '+', '[', '|', '.', '/', '?':
		sb.AppendChar('\\').AppendChar(ch)
	case '{':
		// If possibly part of a repetition, then escape.
		// Concatenation is handled by the digitsMightBleed check.
		if '0' <= next && next <= '9' {
			sb.AppendChar('\\')
		}
		sb.AppendChar(ch)
	default:
		RegExpTree_EscapeCharOnto(ch, sb)
	}
}

// Simplify drops empty text and case-canonicalizes under the 'i' flag.
//
// Java: Text.simplify(String)
func (t *RegExpTree_Text) Simplify(flags string) RegExpTree {
	n := jutil.StringLength(t.text)
	if n == 0 {
		return RegExpTree_Empty_INSTANCE
	}
	if strings.IndexByte(flags, 'i') >= 0 {
		canonicalized := CaseCanonicalize_CaseCanonicalize(t.text)
		if t.text != canonicalized {
			return NewRegExpTree_Text(canonicalized)
		}
	}
	return t
}

// IsCaseSensitive reports whether any character is case sensitive.
//
// Java: Text.isCaseSensitive()
func (t *RegExpTree_Text) IsCaseSensitive() bool {
	chars := jutil.StringToUTF16(t.text) // Java: text.charAt(i)
	for i, n := 0, len(chars); i < n; i++ {
		if CaseCanonicalize_CASE_SENSITIVE.Contains(int(chars[i])) {
			return true
		}
	}
	return false
}

// AppendSourceCode emits the escaped text.
//
// Java: Text.appendSourceCode(StringBuilder)
func (t *RegExpTree_Text) AppendSourceCode(sb *jutil.StringBuilder) {
	chars := jutil.StringToUTF16(t.text) // Java: text.charAt(i)
	for i, n := 0, len(chars); i < n; i++ {
		next := -1
		if i+1 < n {
			next = int(chars[i+1])
		}
		RegExpTree_Text_EscapeRegularCharOnto(chars[i], next, sb)
	}
}

// AppendDebugInfo emits `text`.
//
// Java: Text.appendDebugInfo(StringBuilder)
func (t *RegExpTree_Text) AppendDebugInfo(sb *jutil.StringBuilder) {
	sb.AppendChar('`').Append(t.text).AppendChar('`')
}

// Equals compares the text.
//
// Java: Text.equals(Object)
func (t *RegExpTree_Text) Equals(o any) bool {
	text, ok := o.(*RegExpTree_Text)
	return ok && t.text == text.text
}

// HashCode hashes the text.
//
// Java: Text.hashCode()
func (t *RegExpTree_Text) HashCode() int32 {
	return jutil.StringHashCode(t.text) ^ 0x617e310
}

// RegExpTree_Repetition represents a repeating item such as ...+, ...*, or
// ...{0,1}.
//
// Java: public static final class RegExpTree.Repetition
type RegExpTree_Repetition struct {
	RegExpTreeBase
	body   RegExpTree
	min    int
	max    int
	greedy bool
}

// NewRegExpTree_Repetition constructs a Repetition.
//
// Java: Repetition(RegExpTree, int, int, boolean)
func NewRegExpTree_Repetition(body RegExpTree, min_ int, max_ int, greedy bool) *RegExpTree_Repetition {
	r := &RegExpTree_Repetition{}
	r.InitRegExpTreeBase(r)
	r.body = body
	r.min = min_
	r.max = max_
	r.greedy = greedy
	return r
}

// Simplify folds nested and trivial repetitions.
//
// Java: Repetition.simplify(String)
func (r *RegExpTree_Repetition) Simplify(flags string) RegExpTree {
	body := r.body.Simplify(flags)
	if r.max == 0 && !body.HasCapturingGroup() {
		return RegExpTree_Empty_INSTANCE
	}
	if _, ok := body.(*RegExpTree_Empty); ok || RegExpTree_NEVER_MATCHES.Equals(body) {
		return body
	}
	min_ := r.min
	max_ := r.max
	if rbody, ok := body.(*RegExpTree_Repetition); ok {
		if rbody.greedy == r.greedy {
			lmin := int64(min_) * int64(rbody.min)
			lmax := int64(max_) * int64(rbody.max)
			if lmin < javaIntegerMaxValue {
				body = rbody.body
				min_ = int(lmin)
				if lmax >= javaIntegerMaxValue {
					max_ = javaIntegerMaxValue
				} else {
					max_ = int(lmax)
				}
			}
		}
	}
	if min_ == 1 && max_ == 1 {
		return body
	}
	greedy := r.greedy || min_ == max_
	if body.Equals(r.body) && min_ == r.min && max_ == r.max && greedy == r.greedy {
		return r
	}
	return NewRegExpTree_Repetition(body, min_, max_, greedy).Simplify(flags)
}

// IsCaseSensitive delegates to the body.
//
// Java: Repetition.isCaseSensitive()
func (r *RegExpTree_Repetition) IsCaseSensitive() bool {
	return r.body.IsCaseSensitive()
}

// ContainsAnchor delegates to the body.
//
// Java: Repetition.containsAnchor()
func (r *RegExpTree_Repetition) ContainsAnchor() bool {
	return r.body.ContainsAnchor()
}

// NumCapturingGroups delegates to the body.
//
// Java: Repetition.numCapturingGroups()
func (r *RegExpTree_Repetition) NumCapturingGroups() int {
	return r.body.NumCapturingGroups()
}

// Children is [body].
//
// Java: Repetition.children()
func (r *RegExpTree_Repetition) Children() []RegExpTree {
	return []RegExpTree{r.body}
}

// AppendBodySourceCode emits the body, parenthesized when needed.
//
// Java: Repetition.appendBodySourceCode(StringBuilder)
func (r *RegExpTree_Repetition) AppendBodySourceCode(sb *jutil.StringBuilder) {
	_, isAlternation := r.body.(*RegExpTree_Alternation)
	_, isConcatenation := r.body.(*RegExpTree_Concatenation)
	_, isRepetition := r.body.(*RegExpTree_Repetition)
	text, isText := r.body.(*RegExpTree_Text)
	if isAlternation ||
		isConcatenation ||
		isRepetition ||
		(isText && jutil.StringLength(text.text) > 1) {
		sb.Append("(?:")
		r.body.AppendSourceCode(sb)
		sb.AppendChar(')')
	} else {
		r.body.AppendSourceCode(sb)
	}
}

// RegExpTree_Repetition_SuffixLen is the length of the suffix that
// AppendSourceCode renders for min and max.
//
// Java: Repetition.suffixLen(int, int)
func RegExpTree_Repetition_SuffixLen(min_ int, max_ int) int {
	// This mirrors the branches that renders a suffix in appendSourceCode below.
	if max_ == javaIntegerMaxValue {
		switch min_ {
		case 0, // *
			1: // +
			return 1
		default: // {3,}
			return 3 + RegExpTree_Repetition_NumDecimalDigits(min_)
		}
	}
	if min_ == 0 && max_ == 1 {
		return 1 // ?
	}
	if min_ == max_ {
		if min_ == 1 {
			return 0 // No suffix needed for {1}.
		}
		return 2 + RegExpTree_Repetition_NumDecimalDigits(min_) // {4}
	}
	return 3 + RegExpTree_Repetition_NumDecimalDigits(min_) + RegExpTree_Repetition_NumDecimalDigits(max_) // {2,7}
}

// RegExpTree_Repetition_NumDecimalDigits is the number of decimal digits of
// a non-negative n.
//
// Java: Repetition.numDecimalDigits(int)
func RegExpTree_Repetition_NumDecimalDigits(n int) int {
	if n < 0 {
		// Negative values should not be passed in.
		panic(jutil.NewAssertionErrorNoMessage())
		// If changing this code to support negative values,
		// Integer.MIN_VALUE is a corner-case..
	}
	nDigits := 1
	for n >= 10 {
		nDigits++
		n /= 10
	}
	return nDigits
}

// AppendSourceCode emits the body and the shortest suffix.
//
// Java: Repetition.appendSourceCode(StringBuilder)
func (r *RegExpTree_Repetition) AppendSourceCode(sb *jutil.StringBuilder) {
	bodyStart := sb.Length()
	r.AppendBodySourceCode(sb)
	bodyEnd := sb.Length()
	bodyLen := bodyEnd - bodyStart
	min_ := r.min
	max_ := r.max
	if (min_ >= 2 && max_ == javaIntegerMaxValue) || max_-min_ <= 1 {
		var expanded int
		// If min == max then we want to try expanding to the limit and
		// attach the empty suffix, which is equivalent to min = max = 1,
		// i.e. /a/ vs /a{1}/.
		if min_ == max_ ||
			// Give aa+ preference over aaa*.
			max_ == javaIntegerMaxValue {
			expanded = min_ - 1
		} else {
			expanded = min_
		}
		expandedMin := min_ - expanded
		expandedMax := max_ - expanded
		if max_ == javaIntegerMaxValue {
			expandedMax = max_
		}
		suffixLen := RegExpTree_Repetition_SuffixLen(min_, max_)
		expandedSuffixLen := RegExpTree_Repetition_SuffixLen(expandedMin, expandedMax)
		// Go-only: language-adapter — Java int arithmetic (wraps at 32 bits).
		if int32(bodyLen)*int32(expanded)+int32(expandedSuffixLen) < int32(suffixLen) && !r.body.HasCapturingGroup() {
			// a{2} -> aa
			// a{2,} -> aa+
			// a{2,3} -> aaa?
			for expanded--; expanded >= 0; expanded-- {
				sb.AppendCharSequenceRange(sb, bodyStart, bodyEnd)
			}
			min_ = expandedMin
			max_ = expandedMax
		}
	}

	if max_ == javaIntegerMaxValue {
		switch min_ {
		case 0:
			sb.AppendChar('*')
		case 1:
			sb.AppendChar('+')
		default:
			sb.AppendChar('{').AppendInt(int64(min_)).Append(",}")
		}
	} else if min_ == 0 && max_ == 1 {
		sb.AppendChar('?')
	} else if min_ == max_ {
		if min_ != 1 {
			sb.AppendChar('{').AppendInt(int64(min_)).AppendChar('}')
		}
	} else {
		sb.AppendChar('{').AppendInt(int64(min_)).AppendChar(',').AppendInt(int64(max_)).AppendChar('}')
	}
	if !r.greedy {
		sb.AppendChar('?')
	}
}

// AppendDebugInfo emits the bounds and greediness.
//
// Java: Repetition.appendDebugInfo(StringBuilder)
func (r *RegExpTree_Repetition) AppendDebugInfo(sb *jutil.StringBuilder) {
	sb.Append(" min=").AppendInt(int64(r.min)).Append(", max=").AppendInt(int64(r.max))
	if !r.greedy {
		sb.Append("  not_greedy")
	}
}

// Equals compares body, bounds and greediness.
//
// Java: Repetition.equals(Object)
func (r *RegExpTree_Repetition) Equals(o any) bool {
	that, ok := o.(*RegExpTree_Repetition)
	if !ok {
		return false
	}
	return r.body.Equals(that.body) &&
		r.min == that.min &&
		r.max == that.max &&
		r.greedy == that.greedy
}

// HashCode hashes body, bounds and greediness.
//
// Java: Repetition.hashCode()
func (r *RegExpTree_Repetition) HashCode() int32 {
	greedy := int32(0)
	if r.greedy {
		greedy = 1
	}
	return int32(r.min) + 31*(int32(r.max)+31*(greedy+31*r.body.HashCode()))
}

// RegExpTree_Alternation represents the possibilities ["foo", "bar" ] for a
// RegExp /foo|bar/.
//
// Java: public static final class RegExpTree.Alternation
type RegExpTree_Alternation struct {
	RegExpTreeBase
	alternatives []RegExpTree
}

// NewRegExpTree_Alternation constructs an Alternation over a copy of
// alternatives.
//
// Java: Alternation(List<? extends RegExpTree>)
func NewRegExpTree_Alternation(alternatives []RegExpTree) *RegExpTree_Alternation {
	a := &RegExpTree_Alternation{}
	a.InitRegExpTreeBase(a)
	a.alternatives = slices.Clone(alternatives) // Java: ImmutableList.copyOf
	return a
}

// Simplify flattens, deduplicates and collapses character alternatives.
//
// Java: Alternation.simplify(String)
func (a *RegExpTree_Alternation) Simplify(flags string) RegExpTree {
	alternatives := []RegExpTree{}
	for _, alternative := range a.alternatives {
		alternative = alternative.Simplify(flags)
		if alternation, ok := alternative.(*RegExpTree_Alternation); ok {
			alternatives = append(alternatives, alternation.alternatives...)
		} else {
			alternatives = append(alternatives, alternative)
		}
	}
	// Remove duplicates
	var last RegExpTree
	for i := 0; i < len(alternatives); { // Java: Iterator with it.remove()
		alternative := alternatives[i]
		if alternative.Equals(RegExpTree_NEVER_MATCHES) {
			i++
			continue
		}
		if alternative.Equals(last) && !alternative.HasCapturingGroup() {
			alternatives = slices.Delete(alternatives, i, i+1)
		} else {
			last = alternative
			i++
		}
	}
	// Collapse character alternatives into character sets.
	for i, n := 0, len(alternatives); i < n; i++ {
		alternative := alternatives[i]
		_, isCharset := alternative.(*RegExpTree_Charset)
		if regExpTree_isSingleCharText(alternative) || isCharset {
			end := i
			nCharsets := 0
			for end < n {
				follower := alternatives[end]
				if _, ok := follower.(*RegExpTree_Charset); ok {
					nCharsets++
				} else if !regExpTree_isSingleCharText(follower) {
					break
				}
				end++
			}
			if end-i >= 3 || (nCharsets != 0 && end-i >= 2) {
				members := make([]int, end-i-nCharsets)
				memberIdx := 0
				chars := CharRanges_EMPTY
				ieExplicits := CharRanges_EMPTY
				charAlternatives := alternatives[i:end] // Java: alternatives.subList(i, end)
				for _, charAlternative := range charAlternatives {
					if text, ok := charAlternative.(*RegExpTree_Text); ok {
						ch := jutil.StringCharAt(text.text, 0)
						members[memberIdx] = int(ch)
						memberIdx++
						if RegExpTree_IE_SPEC_ERRORS.Contains(int(ch)) {
							ieExplicits = ieExplicits.Union(CharRanges_Inclusive(int(ch), int(ch)))
						}
					} else if cs, ok := charAlternative.(*RegExpTree_Charset); ok {
						chars = chars.Union(cs.ranges)
						ieExplicits = ieExplicits.Union(cs.ieExplicits)
					}
				}
				chars = chars.Union(CharRanges_WithMembers(members...))
				// Java: charAlternatives.clear(); charAlternatives.add(...)
				alternatives = slices.Replace(alternatives, i, end,
					NewRegExpTree_Charset(chars, ieExplicits).Simplify(flags))
				n = len(alternatives)
			}
		}
	}
	switch len(alternatives) {
	case 0:
		return RegExpTree_Empty_INSTANCE
	case 1:
		return alternatives[0]
	case 2:
		if _, ok := alternatives[1].(*RegExpTree_Empty); ok { // (?:a|) -> a?
			return NewRegExpTree_Repetition(alternatives[0], 0, 1, true)
		} else if _, ok := alternatives[0].(*RegExpTree_Empty); ok {
			return NewRegExpTree_Repetition(alternatives[1], 0, 1, false)
		}
	}
	// TODO: maybe pull out common prefix or suffix
	if regExpTree_listEquals(alternatives, a.alternatives) {
		return a
	}
	return NewRegExpTree_Alternation(alternatives)
}

// regExpTree_isSingleCharText is `x instanceof Text text &&
// text.text.length() == 1`.
//
// Go-only: language-adapter — pattern-matching instanceof as a helper.
func regExpTree_isSingleCharText(t RegExpTree) bool {
	text, ok := t.(*RegExpTree_Text)
	return ok && jutil.StringLength(text.text) == 1
}

// IsCaseSensitive reports whether any alternative is.
//
// Java: Alternation.isCaseSensitive()
func (a *RegExpTree_Alternation) IsCaseSensitive() bool {
	for _, alternative := range a.alternatives {
		if alternative.IsCaseSensitive() {
			return true
		}
	}
	return false
}

// ContainsAnchor reports whether any alternative does.
//
// Java: Alternation.containsAnchor()
func (a *RegExpTree_Alternation) ContainsAnchor() bool {
	for _, alternative := range a.alternatives {
		if alternative.ContainsAnchor() {
			return true
		}
	}
	return false
}

// NumCapturingGroups sums over the alternatives.
//
// Java: Alternation.numCapturingGroups()
func (a *RegExpTree_Alternation) NumCapturingGroups() int {
	n := 0
	for _, alternative := range a.alternatives {
		n += alternative.NumCapturingGroups()
	}
	return n
}

// Children returns the alternatives.
//
// Java: Alternation.children()
func (a *RegExpTree_Alternation) Children() []RegExpTree {
	return a.alternatives
}

// AppendSourceCode emits the alternatives separated by '|'.
//
// Java: Alternation.appendSourceCode(StringBuilder)
func (a *RegExpTree_Alternation) AppendSourceCode(sb *jutil.StringBuilder) {
	for i, n := 0, len(a.alternatives); i < n; i++ {
		if i != 0 {
			sb.AppendChar('|')
		}
		a.alternatives[i].AppendSourceCode(sb)
	}
}

// AppendDebugInfo emits nothing besides children.
//
// Java: Alternation.appendDebugInfo(StringBuilder)
func (a *RegExpTree_Alternation) AppendDebugInfo(sb *jutil.StringBuilder) {
	// Nothing besides children.
}

// Equals compares the alternatives.
//
// Java: Alternation.equals(Object)
func (a *RegExpTree_Alternation) Equals(o any) bool {
	if alternation, ok := o.(*RegExpTree_Alternation); ok {
		return a == alternation || regExpTree_listEquals(a.alternatives, alternation.alternatives)
	}
	return false
}

// HashCode hashes the alternatives.
//
// Java: Alternation.hashCode()
func (a *RegExpTree_Alternation) HashCode() int32 {
	return 0x51b57cd1 ^ regExpTree_listHashCode(a.alternatives)
}

// regExpTree_listEquals is List.equals on lists of RegExpTree.
//
// Go-only: language-adapter — java.util.List.equals.
func regExpTree_listEquals(a, b []RegExpTree) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !a[i].Equals(b[i]) {
			return false
		}
	}
	return true
}

// regExpTree_listHashCode is List.hashCode on lists of RegExpTree.
//
// Go-only: language-adapter — java.util.List.hashCode.
func regExpTree_listHashCode(l []RegExpTree) int32 {
	hashCode := int32(1)
	for _, e := range l {
		hashCode = 31*hashCode + e.HashCode()
	}
	return hashCode
}

// RegExpTree_NEVER_MATCHES is (?!), which matches nothing.
var RegExpTree_NEVER_MATCHES RegExpTree = NewRegExpTree_LookaheadAssertion(RegExpTree_Empty_INSTANCE, false)

// RegExpTree_LookaheadAssertion represents a lookahead assertion such as
// (?=...) or (?!...).
//
// Java: private static final class RegExpTree.LookaheadAssertion
type RegExpTree_LookaheadAssertion struct {
	RegExpTreeBase
	body     RegExpTree
	positive bool
}

// NewRegExpTree_LookaheadAssertion constructs a LookaheadAssertion.
//
// Java: LookaheadAssertion(RegExpTree, boolean)
func NewRegExpTree_LookaheadAssertion(body RegExpTree, positive bool) *RegExpTree_LookaheadAssertion {
	l := &RegExpTree_LookaheadAssertion{}
	l.InitRegExpTreeBase(l)
	l.body = body
	l.positive = positive
	return l
}

// Simplify simplifies the body; (?=) is always true.
//
// Java: LookaheadAssertion.simplify(String)
func (l *RegExpTree_LookaheadAssertion) Simplify(flags string) RegExpTree {
	simpleBody := l.body.Simplify(flags)
	if _, ok := simpleBody.(*RegExpTree_Empty); ok {
		if l.positive { // Always true
			return simpleBody
		}
	}
	return NewRegExpTree_LookaheadAssertion(simpleBody, l.positive)
}

// IsCaseSensitive delegates to the body.
//
// Java: LookaheadAssertion.isCaseSensitive()
func (l *RegExpTree_LookaheadAssertion) IsCaseSensitive() bool {
	return l.body.IsCaseSensitive()
}

// ContainsAnchor delegates to the body.
//
// Java: LookaheadAssertion.containsAnchor()
func (l *RegExpTree_LookaheadAssertion) ContainsAnchor() bool {
	return l.body.ContainsAnchor()
}

// NumCapturingGroups delegates to the body.
//
// Java: LookaheadAssertion.numCapturingGroups()
func (l *RegExpTree_LookaheadAssertion) NumCapturingGroups() int {
	return l.body.NumCapturingGroups()
}

// Children is [body].
//
// Java: LookaheadAssertion.children()
func (l *RegExpTree_LookaheadAssertion) Children() []RegExpTree {
	return []RegExpTree{l.body}
}

// AppendSourceCode emits (?=...) or (?!...).
//
// Java: LookaheadAssertion.appendSourceCode(StringBuilder)
func (l *RegExpTree_LookaheadAssertion) AppendSourceCode(sb *jutil.StringBuilder) {
	if l.positive {
		sb.Append("(?=")
	} else {
		sb.Append("(?!")
	}
	l.body.AppendSourceCode(sb)
	sb.AppendChar(')')
}

// AppendDebugInfo emits positive or negative.
//
// Java: LookaheadAssertion.appendDebugInfo(StringBuilder)
func (l *RegExpTree_LookaheadAssertion) AppendDebugInfo(sb *jutil.StringBuilder) {
	if l.positive {
		sb.Append("positive")
	} else {
		sb.Append("negative")
	}
}

// Equals compares polarity and body.
//
// Java: LookaheadAssertion.equals(Object)
func (l *RegExpTree_LookaheadAssertion) Equals(o any) bool {
	that, ok := o.(*RegExpTree_LookaheadAssertion)
	if !ok {
		return false
	}
	return l.positive == that.positive && l.body.Equals(that.body)
}

// HashCode hashes the body.
//
// Java: LookaheadAssertion.hashCode()
func (l *RegExpTree_LookaheadAssertion) HashCode() int32 {
	return 0x723aba9 ^ l.body.HashCode()
}

// RegExpTree_LookbehindAssertion represents a lookbehind assertion such as
// (?<=...) or (?<!...).
//
// Java: public static final class RegExpTree.LookbehindAssertion
type RegExpTree_LookbehindAssertion struct {
	RegExpTreeBase
	body     RegExpTree
	positive bool
}

// NewRegExpTree_LookbehindAssertion constructs a LookbehindAssertion.
//
// Java: LookbehindAssertion(RegExpTree, boolean)
func NewRegExpTree_LookbehindAssertion(body RegExpTree, positive bool) *RegExpTree_LookbehindAssertion {
	l := &RegExpTree_LookbehindAssertion{}
	l.InitRegExpTreeBase(l)
	l.body = body
	l.positive = positive
	return l
}

// Simplify simplifies the body; (?<=) is always true.
//
// Java: LookbehindAssertion.simplify(String)
func (l *RegExpTree_LookbehindAssertion) Simplify(flags string) RegExpTree {
	simpleBody := l.body.Simplify(flags)
	if _, ok := simpleBody.(*RegExpTree_Empty); ok {
		if l.positive { // Always true
			return simpleBody
		}
	}
	return NewRegExpTree_LookbehindAssertion(simpleBody, l.positive)
}

// IsCaseSensitive delegates to the body.
//
// Java: LookbehindAssertion.isCaseSensitive()
func (l *RegExpTree_LookbehindAssertion) IsCaseSensitive() bool {
	return l.body.IsCaseSensitive()
}

// ContainsAnchor delegates to the body.
//
// Java: LookbehindAssertion.containsAnchor()
func (l *RegExpTree_LookbehindAssertion) ContainsAnchor() bool {
	return l.body.ContainsAnchor()
}

// NumCapturingGroups delegates to the body.
//
// Java: LookbehindAssertion.numCapturingGroups()
func (l *RegExpTree_LookbehindAssertion) NumCapturingGroups() int {
	return l.body.NumCapturingGroups()
}

// Children is [body].
//
// Java: LookbehindAssertion.children()
func (l *RegExpTree_LookbehindAssertion) Children() []RegExpTree {
	return []RegExpTree{l.body}
}

// AppendSourceCode emits (?<=...) or (?<!...).
//
// Java: LookbehindAssertion.appendSourceCode(StringBuilder)
func (l *RegExpTree_LookbehindAssertion) AppendSourceCode(sb *jutil.StringBuilder) {
	if l.positive {
		sb.Append("(?<=")
	} else {
		sb.Append("(?<!")
	}
	l.body.AppendSourceCode(sb)
	sb.AppendChar(')')
}

// AppendDebugInfo emits positive or negative.
//
// Java: LookbehindAssertion.appendDebugInfo(StringBuilder)
func (l *RegExpTree_LookbehindAssertion) AppendDebugInfo(sb *jutil.StringBuilder) {
	if l.positive {
		sb.Append("positive")
	} else {
		sb.Append("negative")
	}
}

// Equals compares polarity and body.
//
// Java: LookbehindAssertion.equals(Object)
func (l *RegExpTree_LookbehindAssertion) Equals(o any) bool {
	that, ok := o.(*RegExpTree_LookbehindAssertion)
	if !ok {
		return false
	}
	return l.positive == that.positive && l.body.Equals(that.body)
}

// HashCode hashes the body.
//
// Java: LookbehindAssertion.hashCode()
func (l *RegExpTree_LookbehindAssertion) HashCode() int32 {
	return 0x723aba9 ^ l.body.HashCode()
}

// RegExpTree_CapturingGroup represents a capturing group such as (asdf).
//
// Java: public static final class RegExpTree.CapturingGroup
type RegExpTree_CapturingGroup struct {
	RegExpTreeBase
	body RegExpTree
}

// NewRegExpTree_CapturingGroup constructs a CapturingGroup.
//
// Java: CapturingGroup(RegExpTree)
func NewRegExpTree_CapturingGroup(body RegExpTree) *RegExpTree_CapturingGroup {
	c := &RegExpTree_CapturingGroup{}
	c.InitRegExpTreeBase(c)
	c.body = body
	return c
}

// Simplify simplifies the body.
//
// Java: CapturingGroup.simplify(String)
func (c *RegExpTree_CapturingGroup) Simplify(flags string) RegExpTree {
	return NewRegExpTree_CapturingGroup(c.body.Simplify(flags))
}

// IsCaseSensitive delegates to the body.
//
// Java: CapturingGroup.isCaseSensitive()
func (c *RegExpTree_CapturingGroup) IsCaseSensitive() bool {
	return c.body.IsCaseSensitive()
}

// ContainsAnchor delegates to the body.
//
// Java: CapturingGroup.containsAnchor()
func (c *RegExpTree_CapturingGroup) ContainsAnchor() bool {
	return c.body.ContainsAnchor()
}

// NumCapturingGroups is one plus the body's.
//
// Java: CapturingGroup.numCapturingGroups()
func (c *RegExpTree_CapturingGroup) NumCapturingGroups() int {
	return 1 + c.body.NumCapturingGroups()
}

// Children is [body].
//
// Java: CapturingGroup.children()
func (c *RegExpTree_CapturingGroup) Children() []RegExpTree {
	return []RegExpTree{c.body}
}

// AppendSourceCode emits (...).
//
// Java: CapturingGroup.appendSourceCode(StringBuilder)
func (c *RegExpTree_CapturingGroup) AppendSourceCode(sb *jutil.StringBuilder) {
	sb.AppendChar('(')
	c.body.AppendSourceCode(sb)
	sb.AppendChar(')')
}

// AppendDebugInfo emits nothing besides children.
//
// Java: CapturingGroup.appendDebugInfo(StringBuilder)
func (c *RegExpTree_CapturingGroup) AppendDebugInfo(sb *jutil.StringBuilder) {
	// Nothing besides children.
}

// Equals compares the body.
//
// Java: CapturingGroup.equals(Object)
func (c *RegExpTree_CapturingGroup) Equals(o any) bool {
	capturingGroup, ok := o.(*RegExpTree_CapturingGroup)
	return ok && c.body.Equals(capturingGroup.body)
}

// HashCode hashes the body.
//
// Java: CapturingGroup.hashCode()
func (c *RegExpTree_CapturingGroup) HashCode() int32 {
	return 0x55781738 ^ c.body.HashCode()
}

// RegExpTree_NamedCaptureGroup represents a named capture group.
//
// Java: public static final class RegExpTree.NamedCaptureGroup
type RegExpTree_NamedCaptureGroup struct {
	RegExpTreeBase
	body RegExpTree
	name string
}

// NewRegExpTree_NamedCaptureGroup constructs a NamedCaptureGroup.
//
// Java: NamedCaptureGroup(RegExpTree, String)
func NewRegExpTree_NamedCaptureGroup(body RegExpTree, name string) *RegExpTree_NamedCaptureGroup {
	n := &RegExpTree_NamedCaptureGroup{}
	n.InitRegExpTreeBase(n)
	n.body = body
	n.name = name
	return n
}

// Simplify simplifies the body.
//
// Java: NamedCaptureGroup.simplify(String)
func (n *RegExpTree_NamedCaptureGroup) Simplify(flags string) RegExpTree {
	return NewRegExpTree_NamedCaptureGroup(n.body.Simplify(flags), n.name)
}

// IsCaseSensitive delegates to the body.
//
// Java: NamedCaptureGroup.isCaseSensitive()
func (n *RegExpTree_NamedCaptureGroup) IsCaseSensitive() bool {
	return n.body.IsCaseSensitive()
}

// ContainsAnchor delegates to the body.
//
// Java: NamedCaptureGroup.containsAnchor()
func (n *RegExpTree_NamedCaptureGroup) ContainsAnchor() bool {
	return n.body.ContainsAnchor()
}

// NumCapturingGroups is one plus the body's.
//
// Java: NamedCaptureGroup.numCapturingGroups()
func (n *RegExpTree_NamedCaptureGroup) NumCapturingGroups() int {
	return 1 + n.body.NumCapturingGroups()
}

// Children is [body].
//
// Java: NamedCaptureGroup.children()
func (n *RegExpTree_NamedCaptureGroup) Children() []RegExpTree {
	return []RegExpTree{n.body}
}

// AppendSourceCode emits (?<name>...).
//
// Java: NamedCaptureGroup.appendSourceCode(StringBuilder)
func (n *RegExpTree_NamedCaptureGroup) AppendSourceCode(sb *jutil.StringBuilder) {
	sb.Append("(?<")
	sb.Append(n.name)
	sb.AppendChar('>')
	n.body.AppendSourceCode(sb)
	sb.AppendChar(')')
}

// AppendDebugInfo emits the name.
//
// Java: NamedCaptureGroup.appendDebugInfo(StringBuilder)
func (n *RegExpTree_NamedCaptureGroup) AppendDebugInfo(sb *jutil.StringBuilder) {
	sb.Append(" name=").Append(n.name)
}

// Equals compares name and body.
//
// Java: NamedCaptureGroup.equals(Object)
func (n *RegExpTree_NamedCaptureGroup) Equals(o any) bool {
	namedCaptureGroup, ok := o.(*RegExpTree_NamedCaptureGroup)
	return ok &&
		n.name == namedCaptureGroup.name &&
		n.body.Equals(namedCaptureGroup.body)
}

// HashCode hashes name and body.
//
// Java: NamedCaptureGroup.hashCode()
func (n *RegExpTree_NamedCaptureGroup) HashCode() int32 {
	return jutil.StringHashCode(n.name) ^ n.body.HashCode() // Java: Objects.hashCode(name)
}

// RegExpTree_UnicodePropertyEscape represents a Unicode Property Escape such
// as in /\p{Script=Greek}/u.
//
// Java: public static final class RegExpTree.UnicodePropertyEscape
type RegExpTree_UnicodePropertyEscape struct {
	RegExpTree_RegExpTreeAtomBase

	propertyName  *string // Java: @Nullable String
	propertyValue string
	negated       bool
}

// NewRegExpTree_UnicodePropertyEscape constructs a UnicodePropertyEscape;
// propertyName is nil for the shorthand form \p{Value}.
//
// Java: UnicodePropertyEscape(@Nullable String, String, boolean)
func NewRegExpTree_UnicodePropertyEscape(propertyName *string, propertyValue string, negated bool) *RegExpTree_UnicodePropertyEscape {
	u := &RegExpTree_UnicodePropertyEscape{}
	u.InitRegExpTree_RegExpTreeAtomBase(u)
	// Java: checkState(propertyValue != null);
	// Go-only: language-adapter — a Go string is never null, so the check
	// always holds.
	jutil.CheckArgumentMsg(
		propertyName == nil || *propertyName != "",
		"if '=' is present in a unicode property escape, the name cannot be empty")
	jutil.CheckArgumentMsg(propertyValue != "", "unicode property escape value cannot be empty")
	u.propertyName = propertyName
	u.propertyValue = propertyValue
	u.negated = negated
	return u
}

// Simplify returns u.
//
// Java: UnicodePropertyEscape.simplify(String)
func (u *RegExpTree_UnicodePropertyEscape) Simplify(flags string) RegExpTree {
	return u
}

// AppendSourceCode emits \p{...} or \P{...}.
//
// Java: UnicodePropertyEscape.appendSourceCode(StringBuilder)
func (u *RegExpTree_UnicodePropertyEscape) AppendSourceCode(sb *jutil.StringBuilder) {
	if u.negated {
		sb.Append("\\P{")
	} else {
		sb.Append("\\p{")
	}
	if u.propertyName != nil {
		sb.Append(*u.propertyName)
		sb.AppendChar('=')
	}
	sb.Append(u.propertyValue)
	sb.Append("}")
}

// AppendDebugInfo emits nothing.
//
// Java: UnicodePropertyEscape.appendDebugInfo(StringBuilder)
func (u *RegExpTree_UnicodePropertyEscape) AppendDebugInfo(sb *jutil.StringBuilder) {
	// Nothing besides properties and possible negation.
}

// Equals compares negation, name and value.
//
// Java: UnicodePropertyEscape.equals(Object)
func (u *RegExpTree_UnicodePropertyEscape) Equals(o any) bool {
	unicodePropertyEscape, ok := o.(*RegExpTree_UnicodePropertyEscape)
	return ok &&
		u.negated == unicodePropertyEscape.negated &&
		regExpTree_nullableStringEquals(u.propertyName, unicodePropertyEscape.propertyName) &&
		u.propertyValue == unicodePropertyEscape.propertyValue
}

// HashCode hashes negation, name and value.
//
// Java: UnicodePropertyEscape.hashCode()
func (u *RegExpTree_UnicodePropertyEscape) HashCode() int32 {
	// Java: Objects.hash(negated, propertyName, propertyValue), i.e.
	// Arrays.hashCode(new Object[] {negated, propertyName, propertyValue}).
	negatedHash := int32(1237) // Boolean.hashCode(false)
	if u.negated {
		negatedHash = 1231
	}
	nameHash := int32(0)
	if u.propertyName != nil {
		nameHash = jutil.StringHashCode(*u.propertyName)
	}
	result := int32(1)
	result = 31*result + negatedHash
	result = 31*result + nameHash
	result = 31*result + jutil.StringHashCode(u.propertyValue)
	return result
}

// regExpTree_nullableStringEquals is Objects.equals on nullable strings.
//
// Go-only: language-adapter — java.util.Objects.equals(String, String).
func regExpTree_nullableStringEquals(a, b *string) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

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
var RegExpTree_NAMED_CHAR_GROUPS = func() *jutil.LinkedHashMap[uint16, *CharRanges] {
	m := jutil.NewLinkedHashMap[uint16, *CharRanges]()
	m.Put('d', RegExpTree_DIGITS)
	m.Put('D', CharRanges_ALL_CODE_UNITS.Difference(RegExpTree_DIGITS))
	m.Put('s', RegExpTree_SPACE_CHARS)
	m.Put('S', CharRanges_ALL_CODE_UNITS.Difference(RegExpTree_SPACE_CHARS))
	m.Put('w', RegExpTree_WORD_CHARS)
	m.Put('W', RegExpTree_INVERSE_WORD_CHARS)
	return m
}()

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
	RegExpTree_RegExpTreeAtomBase
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
	c := &RegExpTree_Charset{}
	c.InitRegExpTree_RegExpTreeAtomBase(c)
	c.ranges = ranges
	c.ieExplicits = ieExplicits
	return c
}

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
		var options []*CharRanges
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
		return NewRegExpTree_Text(jutil.StringFromChar(uint16(best.Start(0))))
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
	namedGroups := jutil.NewStringBuilder()
	rangesInterIeExplicits := ranges.Intersection(c.ieExplicits)
	for {
		groupName := uint16(0)
		var simplest *CharRanges
		minComplexity := RegExpTree_DecomposedCharset_Complexity(ranges)
		for key, group := range RegExpTree_NAMED_CHAR_GROUPS.All() {
			if ranges.ContainsAll(group) {
				withoutGroup := ranges.Difference(group).Union(rangesInterIeExplicits)
				complexity := RegExpTree_DecomposedCharset_Complexity(withoutGroup)
				if complexity < minComplexity {
					simplest = withoutGroup
					groupName = key
					minComplexity = complexity
				}
			}
		}
		if simplest != nil {
			namedGroups.AppendChar('\\').AppendChar(groupName)
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
func (c *RegExpTree_Charset) AppendSourceCode(sb *jutil.StringBuilder) {
	if RegExpTree_DOT_CHARSET.ranges.Equals(c.ranges) {
		sb.AppendChar('.')
		return
	}
	c.Decompose().AppendSourceCode(sb)
}

// AppendDebugInfo emits the ranges.
//
// Java: Charset.appendDebugInfo(StringBuilder)
func (c *RegExpTree_Charset) AppendDebugInfo(sb *jutil.StringBuilder) {
	sb.Append(c.ranges.String())
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
func (d *RegExpTree_DecomposedCharset) AppendSourceCode(sb *jutil.StringBuilder) {
	if d.ranges.IsEmpty() {
		if !d.inverted && jutil.StringLength(d.namedGroups) == 2 {
			sb.Append(d.namedGroups)
			return
		} else if d.ranges.IsEmpty() && d.namedGroups == "" {
			if d.inverted {
				sb.Append("[\\S\\s]")
			} else {
				sb.Append("(?!)")
			}
			return
		}
	}
	sb.AppendChar('[')
	if d.inverted {
		sb.AppendChar('^')
	}
	sb.Append(d.namedGroups)
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
			sb.AppendChar('-')
			RegExpTree_DecomposedCharset_EscapeRangeCharOnto(end, rangesStartCharset, false, true, sb)
		}
	}
	if emitDashAtEnd {
		sb.AppendChar('-')
	}
	sb.AppendChar(']')
}

// RegExpTree_DecomposedCharset_EscapeRangeCharOnto escapes a character that
// appears in a charset.
//
// Java: DecomposedCharset.escapeRangeCharOnto(char, boolean, boolean, boolean, StringBuilder)
func RegExpTree_DecomposedCharset_EscapeRangeCharOnto(
	ch uint16, startIsFlush bool, atStart bool, atEnd bool, sb *jutil.StringBuilder) {
	switch ch {
	case '\b':
		sb.Append("\\b")
	case '^':
		if atStart && startIsFlush {
			sb.Append("\\^")
		} else {
			sb.Append("^")
		}
	case '-':
		if atStart || atEnd {
			sb.Append("-")
		} else {
			sb.Append("\\-")
		}
	case '\\', ']':
		sb.AppendChar('\\').AppendChar(ch)
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

// RegExpTree_Concatenation represents a series of nodes chained one after
// another such as (?:...)[a-z]*(...).
//
// Java: public static final class RegExpTree.Concatenation
type RegExpTree_Concatenation struct {
	RegExpTreeBase
	elements []RegExpTree
}

// NewRegExpTree_ConcatenationOfPair constructs the concatenation of a and b.
//
// Java: Concatenation(RegExpTree, RegExpTree)
func NewRegExpTree_ConcatenationOfPair(a RegExpTree, b RegExpTree) *RegExpTree_Concatenation {
	c := &RegExpTree_Concatenation{}
	c.InitRegExpTreeBase(c)
	c.elements = []RegExpTree{a, b}
	return c
}

// NewRegExpTree_Concatenation constructs a concatenation over a copy of
// elements.
//
// Java: Concatenation(List<? extends RegExpTree>)
func NewRegExpTree_Concatenation(elements []RegExpTree) *RegExpTree_Concatenation {
	c := &RegExpTree_Concatenation{}
	c.InitRegExpTreeBase(c)
	c.elements = slices.Clone(elements) // Java: ImmutableList.copyOf
	return c
}

// Simplify flattens nested concatenations and folds adjacent texts and
// repetitions.
//
// Java: Concatenation.simplify(String)
func (c *RegExpTree_Concatenation) Simplify(flags string) RegExpTree {
	s := NewRegExpTree_Concatenation_Simplifier(flags)
	for _, element := range c.elements {
		s.Simplify(element.Simplify(flags))
	}

	switch len(s.simplified) {
	case 0:
		return RegExpTree_Empty_INSTANCE
	case 1:
		return s.simplified[0]
	default:
		return NewRegExpTree_Concatenation(s.simplified)
	}
}

// RegExpTree_Concatenation_Simplifier accumulates simplified elements.
//
// Java: local class Simplifier in Concatenation.simplify(String); the
// captured flags are a constructor parameter.
type RegExpTree_Concatenation_Simplifier struct {
	flags      string
	simplified []RegExpTree
}

// NewRegExpTree_Concatenation_Simplifier constructs a Simplifier.
//
// Java: new Simplifier() (inside Concatenation.simplify)
func NewRegExpTree_Concatenation_Simplifier(flags string) *RegExpTree_Concatenation_Simplifier {
	s := &RegExpTree_Concatenation_Simplifier{}
	s.flags = flags
	s.simplified = []RegExpTree{}
	return s
}

// Simplify appends t, flattening and folding pairwise.
//
// Java: Simplifier.simplify(RegExpTree)
func (s *RegExpTree_Concatenation_Simplifier) Simplify(t RegExpTree) {
	if concatenation, ok := t.(*RegExpTree_Concatenation); ok {
		for _, child := range concatenation.elements {
			s.Simplify(child)
		}
	} else if _, ok := t.(*RegExpTree_Empty); ok {
		// Do nothing
	} else {
		lastIndex := len(s.simplified) - 1
		if lastIndex >= 0 {
			pairwise := s.SimplifyPairwise(s.simplified[lastIndex], t)
			if pairwise != nil {
				s.simplified[lastIndex] = pairwise
				return
			}
		}
		s.simplified = append(s.simplified, t)
	}
}

// SimplifyPairwise folds two adjacent elements, or returns nil.
//
// Java: Simplifier.simplifyPairwise(RegExpTree, RegExpTree)
func (s *RegExpTree_Concatenation_Simplifier) SimplifyPairwise(before RegExpTree, after RegExpTree) RegExpTree {
	beforeText, ok1 := before.(*RegExpTree_Text)
	afterText, ok2 := after.(*RegExpTree_Text)
	if ok1 && ok2 {
		return NewRegExpTree_Text(beforeText.text + afterText.text).Simplify(s.flags)
	}
	// Fold adjacent repetitions.
	beforeMin := 1
	beforeMax := 1
	beforeBody := before
	beforeGreedy := false
	if r, ok := before.(*RegExpTree_Repetition); ok {
		beforeMin = r.min
		beforeMax = r.max
		beforeBody = r.body
		beforeGreedy = r.greedy
	}
	afterMin := 1
	afterMax := 1
	afterBody := after
	afterGreedy := false
	if r, ok := after.(*RegExpTree_Repetition); ok {
		afterMin = r.min
		afterMax = r.max
		afterBody = r.body
		afterGreedy = r.greedy
	}
	if beforeBody.Equals(afterBody) && !beforeBody.HasCapturingGroup() {
		lmin := int64(beforeMin) + int64(afterMin)
		lmax := int64(beforeMax) + int64(afterMax)
		if lmin < javaIntegerMaxValue {
			min_ := int(lmin)
			max_ := javaIntegerMaxValue
			if lmax < javaIntegerMaxValue {
				max_ = int(lmax)
			}
			return NewRegExpTree_Repetition(
				beforeBody, min_, max_, beforeGreedy || afterGreedy || min_ == max_)
		}
	}
	return nil
}

// IsCaseSensitive reports whether any element is.
//
// Java: Concatenation.isCaseSensitive()
func (c *RegExpTree_Concatenation) IsCaseSensitive() bool {
	for _, element := range c.elements {
		if element.IsCaseSensitive() {
			return true
		}
	}
	return false
}

// ContainsAnchor reports whether any element does.
//
// Java: Concatenation.containsAnchor()
func (c *RegExpTree_Concatenation) ContainsAnchor() bool {
	for _, element := range c.elements {
		if element.ContainsAnchor() {
			return true
		}
	}
	return false
}

// NumCapturingGroups sums over the elements.
//
// Java: Concatenation.numCapturingGroups()
func (c *RegExpTree_Concatenation) NumCapturingGroups() int {
	n := 0
	for _, element := range c.elements {
		n += element.NumCapturingGroups()
	}
	return n
}

// Children returns the elements.
//
// Java: Concatenation.children()
func (c *RegExpTree_Concatenation) Children() []RegExpTree {
	return c.elements
}

// AppendSourceCode emits the elements, guarding against digits that would
// bleed into a preceding back reference or brace.
//
// Java: Concatenation.appendSourceCode(StringBuilder)
func (c *RegExpTree_Concatenation) AppendSourceCode(sb *jutil.StringBuilder) {
	// True if the last content written might consume
	// decimal digits written subsequently.
	digitsMightBleed := false
	for _, element := range c.elements {
		parenthesize := false
		_, isAlternation := element.(*RegExpTree_Alternation)
		_, isConcatenation := element.(*RegExpTree_Concatenation)
		if isAlternation || isConcatenation {
			parenthesize = true
		}
		if parenthesize {
			sb.Append("(?:")
			element.AppendSourceCode(sb)
			sb.AppendChar(')')
		} else {
			start := sb.Length()
			element.AppendSourceCode(sb)
			if digitsMightBleed && sb.Length() > start {
				firstChar := sb.CharAt(start)
				if '0' <= firstChar && firstChar <= '9' {
					// Bleeding happened.
					// If the last character would be ambiguous
					// with a repetition, escape it.
					if sb.CharAt(start-1) == '{' {
						// Concatenation from optimization of
						// /{(?:0,}/ -> /\{0,}/
						sb.InsertChar(start-1, '\\')
					} else {
						// Or parenthesize otherwise.
						// Concatenation from optimization of
						// /(.)\1(?:0)/ -> /(.)\1(?:0)/.
						sb.Insert(start, "(?:").AppendChar(')')
					}
				}
			}
		}
		backReference, isBackReference := element.(*RegExpTree_BackReference)
		text, isText := element.(*RegExpTree_Text)
		digitsMightBleed =
			// \1(?:0) bleeds if there are 10 or more
			// capturing groups preceding.
			(isBackReference && backReference.groupIndex < 10) ||
				// foo{(?:10}) bleeds.
				(isText && strings.HasSuffix(text.text, "{"))
	}
}

// AppendDebugInfo emits nothing besides children.
//
// Java: Concatenation.appendDebugInfo(StringBuilder)
func (c *RegExpTree_Concatenation) AppendDebugInfo(sb *jutil.StringBuilder) {
	// Nothing besides children.
}

// Equals compares the elements.
//
// Java: Concatenation.equals(Object)
func (c *RegExpTree_Concatenation) Equals(o any) bool {
	concatenation, ok := o.(*RegExpTree_Concatenation)
	return ok && regExpTree_listEquals(c.elements, concatenation.elements)
}

// HashCode hashes the elements.
//
// Java: Concatenation.hashCode()
func (c *RegExpTree_Concatenation) HashCode() int32 {
	return 0x20997e3e ^ regExpTree_listHashCode(c.elements)
}

// RegExpTree_EscapeCharOnto escapes a code unit for regular expression
// source.
//
// Java: RegExpTree.escapeCharOnto(char, StringBuilder)
func RegExpTree_EscapeCharOnto(ch uint16, sb *jutil.StringBuilder) {
	switch ch {
	case '\u0000':
		sb.Append("\\0")
	case '\f':
		sb.Append("\\f")
	case '\t':
		sb.Append("\\t")
	case '\n':
		sb.Append("\\n")
	case '\r':
		sb.Append("\\r")
	case '\\':
		sb.Append("\\\\")
	default:
		if ch < 0x20 || ch >= 0x7f {
			if ch >= 0x100 {
				sb.Append("\\u")
				sb.AppendChar(uint16("0123456789abcdef"[(ch>>12)&0xf]))
				sb.AppendChar(uint16("0123456789abcdef"[(ch>>8)&0xf]))
				sb.AppendChar(uint16("0123456789abcdef"[(ch>>4)&0xf]))
				sb.AppendChar(uint16("0123456789abcdef"[(ch)&0xf]))
			} else {
				sb.Append("\\x")
				sb.AppendChar(uint16("0123456789abcdef"[(ch>>4)&0xf]))
				sb.AppendChar(uint16("0123456789abcdef"[(ch)&0xf]))
			}
		} else {
			sb.AppendChar(ch)
		}
	}
}
