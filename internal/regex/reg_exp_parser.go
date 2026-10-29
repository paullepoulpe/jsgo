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
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// RegExpParseErrorKind classifies a regular expression parse failure.
type RegExpParseErrorKind uint8

const (
	RegExpParseErrorIllegalArgument RegExpParseErrorKind = iota
	RegExpParseErrorIndexOutOfBounds
	RegExpParseErrorStringIndexOutOfBounds
)

// RegExpParseError describes a syntax or bounds error encountered while parsing
// a regular expression pattern.
type RegExpParseError struct {
	Kind    RegExpParseErrorKind
	Message string
	isRange bool
	index   int
	end     int
	length  int
}

func (e *RegExpParseError) Error() string {
	return e.Message
}

func (e *RegExpParseError) toThrowable() jutil.IRuntimeException {
	switch e.Kind {
	case RegExpParseErrorIndexOutOfBounds:
		return jutil.NewIndexOutOfBoundsException(e.Message)
	case RegExpParseErrorStringIndexOutOfBounds:
		if e.isRange {
			return jutil.NewStringIndexOutOfBoundsExceptionForRange(e.index, e.end, e.length)
		}
		return jutil.NewStringIndexOutOfBoundsExceptionForIndex(e.index, e.length)
	default:
		return jutil.NewIllegalArgumentExceptionMsg(e.Message)
	}
}

func illegalArgError(msg string) *RegExpParseError {
	return &RegExpParseError{
		Kind:    RegExpParseErrorIllegalArgument,
		Message: msg,
	}
}

func stringIndexError(index, length int) *RegExpParseError {
	return &RegExpParseError{
		Kind:    RegExpParseErrorStringIndexOutOfBounds,
		Message: "Index " + strconv.Itoa(index) + " out of bounds for length " + strconv.Itoa(length),
		index:   index,
		length:  length,
	}
}

func stringRangeError(begin, end, length int) *RegExpParseError {
	return &RegExpParseError{
		Kind:    RegExpParseErrorStringIndexOutOfBounds,
		Message: "Range [" + strconv.Itoa(begin) + ", " + strconv.Itoa(end) + ") out of bounds for length " + strconv.Itoa(length),
		isRange: true,
		index:   begin,
		end:     end,
		length:  length,
	}
}

func charRangeOrderError(start, end int) *RegExpParseError {
	return &RegExpParseError{
		Kind:    RegExpParseErrorIndexOutOfBounds,
		Message: strconv.Itoa(start) + " > " + strconv.Itoa(end),
	}
}

type parentheticalType uint8

const (
	parentheticalCapturing parentheticalType = iota
	parentheticalNonCapturing
	parentheticalPositiveLookahead
	parentheticalNegativeLookahead
	parentheticalPositiveLookbehind
	parentheticalNegativeLookbehind
	parentheticalNamedGroups
)

// ParseRegExp parses a regular expression pattern and flags into a RegExpTree,
// returning a structured *RegExpParseError instead of panicking on malformed input.
//
// pattern is the "foo" from /foo/i; flags is the "i" from /foo/i.
func ParseRegExp(pattern string, flags string) (RegExpTree, *RegExpParseError) {
	var p regExpParser
	p.init(pattern, flags)
	return p.parseTopLevel()
}

// RegExpTree_ParseRegExp parses a regular expression to an AST, panicking with
// the matching Java RuntimeException on malformed input.
//
// Java: RegExpTree.parseRegExp(String, String)
func RegExpTree_ParseRegExp(pattern string, flags string) RegExpTree {
	tree, err := ParseRegExp(pattern, flags)
	if err != nil {
		panic(err.toThrowable())
	}
	return tree
}

// regExpParser is a recursive descent parser over pattern and flags.
//
// Java: local class Parser in RegExpTree.parseRegExp(String, String).
type regExpParser struct {
	pattern                           string
	patternChars                      []uint16
	pos                               int
	numCapturingGroups                int
	capturingGroupNames               map[string]struct{}
	limit                             int
	isASCII                           bool
	isCaseInsensitive                 bool
	hasUFlag                          bool
	lookForNamedCaptureBackreferences bool
}

func (p *regExpParser) init(pattern, flags string) {
	p.pattern = pattern
	p.isCaseInsensitive = strings.IndexByte(flags, 'i') >= 0
	p.hasUFlag = strings.IndexByte(flags, 'u') >= 0
	if jutil.IsASCII(pattern) {
		p.isASCII = true
		p.limit = len(pattern)
	} else {
		p.patternChars = jutil.StringToUTF16(pattern)
		p.limit = len(p.patternChars)
	}
}

// codeUnitAt returns the UTF-16 code unit at i when 0 <= i < p.limit is already known.
func (p *regExpParser) codeUnitAt(i int) uint16 {
	if p.isASCII {
		return uint16(p.pattern[i])
	}
	return p.patternChars[i]
}

func (p *regExpParser) charAt(i int) (uint16, *RegExpParseError) {
	if i < 0 || i >= p.limit {
		return 0, stringIndexError(i, p.limit)
	}
	return p.codeUnitAt(i), nil
}

func (p *regExpParser) substring(begin, end int) (string, *RegExpParseError) {
	if begin < 0 || end > p.limit || begin > end {
		return "", stringRangeError(begin, end, p.limit)
	}
	if p.isASCII {
		return p.pattern[begin:end], nil
	}
	return jutil.StringFromUTF16(p.patternChars[begin:end]), nil
}

func (p *regExpParser) sliceError(prefix string, begin, end int) *RegExpParseError {
	sub, err := p.substring(begin, end)
	if err != nil {
		return err
	}
	return illegalArgError(prefix + sub)
}

func (p *regExpParser) tailError(prefix string, begin int) *RegExpParseError {
	return p.sliceError(prefix, begin, p.limit)
}

// parseTopLevel parses the whole pattern.
//
// Java: Parser.parseTopLevel()
func (p *regExpParser) parseTopLevel() (RegExpTree, *RegExpParseError) {
	// First assume there are no named capture backreferences,
	// because the spec says they should only be recognized
	// if the pattern contains at least one named capture group.
	p.pos = 0
	p.numCapturingGroups = 0
	p.lookForNamedCaptureBackreferences = false
	out, err := p.parse()
	if err != nil {
		return nil, err
	}

	// If there were named capture groups, we must parse the pattern string again
	// so we can check for any backreferences to them.
	if len(p.capturingGroupNames) != 0 {
		p.pos = 0
		p.numCapturingGroups = 0
		p.lookForNamedCaptureBackreferences = true
		out, err = p.parse()
		if err != nil {
			return nil, err
		}
	}

	if p.pos < p.limit { // Unmatched closed paren maybe.
		return nil, p.tailError("", p.pos)
	}
	return out, nil
}

// parse parses a disjunction up to the end of the pattern or an unmatched
// close paren.
//
// Java: Parser.parse()
func (p *regExpParser) parse() (RegExpTree, *RegExpParseError) {
	// Collects ["foo", "bar", "baz"] for /foo|bar|baz/.
	var alternatives []RegExpTree
	// The last item parsed within an alternation.
	var preceder RegExpTree

topLoop:
	for p.pos < p.limit {
		ch := p.codeUnitAt(p.pos)
		var atom RegExpTree
		var err *RegExpParseError
		switch ch {
		case '[':
			atom, err = p.parseCharset()
		case '(':
			atom, err = p.parseParenthetical()
		case ')':
			break topLoop
		case '\\':
			atom, err = p.parseEscape()
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
				switch p.codeUnitAt(end) {
				case '[', '(', ')', '\\', '^', '$', '|', '.', '*', '+', '?', '{':
					break charsLoop
				default:
					// Repetition binds more tightly than concatenation.
					// Only consume up to "foo" in /foob*/ so that the suffix
					// operator parser below has the right precedence.
					if end+1 >= p.limit || !isRepetitionStart(p.codeUnitAt(end+1)) {
						end++
					} else {
						break charsLoop
					}
				}
			}
			var text string
			text, err = p.substring(start, end)
			if err == nil {
				atom = NewRegExpTree_Text(text)
				p.pos = end
			}
		}
		if err != nil {
			return nil, err
		}
		if p.pos < p.limit && isRepetitionStart(p.codeUnitAt(p.pos)) {
			atom = p.parseRepetition(atom)
		}
		if preceder == nil {
			preceder = atom
		} else {
			preceder = NewRegExpTree_ConcatenationOfPair(preceder, atom)
		}
		// If this is an alternative in an alternation, then add it to the
		// list of complete alternatives, and reset the parser state for the
		// next alternative.
		if p.pos < p.limit && p.codeUnitAt(p.pos) == '|' {
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
		return NewRegExpTree_Alternation(alternatives), nil
	}
	return preceder, nil
}

// parseParenthetical handles capturing groups (...), non-capturing groups
// (?:...), and lookahead/lookbehind assertions.
//
// Java: Parser.parseParenthetical()
func (p *regExpParser) parseParenthetical() (RegExpTree, *RegExpParseError) {
	start := p.pos
	p.pos++
	var kind parentheticalType
	var captureName string
	if p.pos < p.limit && p.codeUnitAt(p.pos) == '?' {
		if p.pos+1 < p.limit {
			ch := p.codeUnitAt(p.pos + 1)
			switch ch {
			case ':':
				p.pos += 2
				kind = parentheticalNonCapturing
			case '=':
				p.pos += 2
				kind = parentheticalPositiveLookahead
			case '!':
				p.pos += 2
				kind = parentheticalNegativeLookahead
			case '<':
				if p.pos+2 < p.limit && p.codeUnitAt(p.pos+2) == '=' {
					p.pos += 3
					kind = parentheticalPositiveLookbehind
				} else if p.pos+2 < p.limit && p.codeUnitAt(p.pos+2) == '!' {
					p.pos += 3
					kind = parentheticalNegativeLookbehind
				} else {
					p.pos += 2
					name, err := p.scanNamedGroupName()
					if err != nil {
						return nil, err
					}
					captureName = name
					if p.capturingGroupNames == nil {
						p.capturingGroupNames = make(map[string]struct{})
					}
					p.capturingGroupNames[captureName] = struct{}{}
					kind = parentheticalNamedGroups
				}
			default:
				return nil, p.tailError("Malformed parenthetical: ", start)
			}
		} else {
			return nil, p.tailError("Malformed parenthetical: ", start)
		}
	} else {
		kind = parentheticalCapturing
	}
	body, err := p.parse()
	if err != nil {
		return nil, err
	}
	if p.pos < p.limit && p.codeUnitAt(p.pos) == ')' {
		p.pos++
	} else {
		return nil, p.tailError("Unclosed parenthetical group: ", start)
	}
	switch kind {
	case parentheticalCapturing:
		p.numCapturingGroups++
		return NewRegExpTree_CapturingGroup(body), nil
	case parentheticalNonCapturing:
		return body, nil
	case parentheticalPositiveLookahead:
		return NewRegExpTree_LookaheadAssertion(body, true), nil
	case parentheticalNegativeLookahead:
		return NewRegExpTree_LookaheadAssertion(body, false), nil
	case parentheticalPositiveLookbehind:
		return NewRegExpTree_LookbehindAssertion(body, true), nil
	case parentheticalNegativeLookbehind:
		return NewRegExpTree_LookbehindAssertion(body, false), nil
	default:
		p.numCapturingGroups++
		return NewRegExpTree_NamedCaptureGroup(body, captureName), nil
	}
}

// scanNamedGroupName scans the pattern for a named group name. Assumes that
// pos points to the character after '<'.
//
// Java: Parser.scanNamedGroupName()
func (p *regExpParser) scanNamedGroupName() (string, *RegExpParseError) {
	start := p.pos
	end := p.pos
	first, err := p.charAt(start)
	if err != nil {
		return "", err
	}
	if !RegExpTree_IsIdentifierStart(first) {
		return "", p.tailError("Invalid capture group name: <", start)
	}
	end++
	for end < p.limit {
		ch := p.codeUnitAt(end)
		if ch == '>' {
			p.pos = end + 1
			return p.substring(start, end)
		} else if RegExpTree_IsIdentifierPart(ch) {
			end++
		} else {
			return "", p.tailError("Invalid capture group name: <", start)
		}
	}
	return "", p.tailError("Malformed named capture group: <", start)
}

// parseCharset parses a square bracketed character set. Standalone character
// groups /\d/ are handled by parseEscape.
//
// Java: Parser.parseCharset()
func (p *regExpParser) parseCharset() (RegExpTree, *RegExpParseError) {
	p.pos++
	inverse := p.pos < p.limit && p.codeUnitAt(p.pos) == '^'
	if inverse {
		p.pos++
	}
	ranges := CharRanges_EMPTY
	ieExplicits := CharRanges_EMPTY
	for p.pos < p.limit && p.codeUnitAt(p.pos) != ']' {
		ch := p.codeUnitAt(p.pos)
		var start int
		if ch == '\\' {
			p.pos++
			possibleGroupName, err := p.charAt(p.pos)
			if err != nil {
				return nil, err
			}
			group := regExpTree_namedCharGroup(possibleGroupName)
			if group != nil {
				p.pos++
				ranges = ranges.Union(group)
				continue
			}
			start, err = p.parseEscapeChar()
			if err != nil {
				return nil, err
			}
		} else {
			start = int(ch)
			p.pos++
		}
		end := start
		if p.pos+1 < p.limit && p.codeUnitAt(p.pos) == '-' && p.codeUnitAt(p.pos+1) != ']' {
			p.pos++
			ch = p.codeUnitAt(p.pos)
			if ch == '\\' {
				p.pos++
				var err *RegExpParseError
				end, err = p.parseEscapeChar()
				if err != nil {
					return nil, err
				}
			} else {
				end = int(ch)
				p.pos++
			}
		}
		if start > end {
			return nil, charRangeOrderError(start, end)
		}
		range_ := CharRanges_Inclusive(start, end)
		ranges = ranges.Union(range_)
		if RegExpTree_IE_SPEC_ERRORS.Contains(start) && RegExpTree_IE_SPEC_ERRORS.Contains(end) {
			ieExplicits = ieExplicits.Union(range_.Intersection(RegExpTree_IE_SPEC_ERRORS))
		}
		if p.isCaseInsensitive {
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

	return NewRegExpTree_Charset(ranges, ieExplicits), nil
}

// parseEscapeChar parses an escape to a code point.
//
// Java: Parser.parseEscapeChar()
func (p *regExpParser) parseEscapeChar() (int, *RegExpParseError) {
	p.pos++
	ch, err := p.charAt(p.pos - 1)
	if err != nil {
		return 0, err
	}
	switch ch {
	case 'b':
		return '\b', nil
	case 'f':
		return '\f', nil
	case 'n':
		return '\n', nil
	case 'r':
		return '\r', nil
	case 't':
		return '\t', nil
	case 'u':
		if p.hasUFlag && p.pos < p.limit && p.codeUnitAt(p.pos) == '{' {
			return p.parseBracedUnicodeEscape()
		}
		return p.parseHex(4)
	case 'v':
		return '\u000b', nil
	case 'x':
		return p.parseHex(2)
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
				ch = p.codeUnitAt(p.pos)
				if '0' <= ch && ch <= '7' {
					codeUnit = (codeUnit << 3) + (ch - '0')
					p.pos++
				} else {
					break
				}
			}
			return int(codeUnit), nil
		}
		return int(ch), nil
	}
}

func isUnicodePropertyIdentChar(ch uint16) bool {
	return ch == '_' ||
		('a' <= ch && ch <= 'z') ||
		('A' <= ch && ch <= 'Z') ||
		('0' <= ch && ch <= '9')
}

// parseEscape parses an escape that appears outside a charset.
//
// Java: Parser.parseEscape()
func (p *regExpParser) parseEscape() (RegExpTree, *RegExpParseError) {
	start := p.pos
	p.pos++
	ch, err := p.charAt(p.pos)
	if err != nil {
		return nil, err
	}
	if ch == 'b' || ch == 'B' {
		p.pos++
		return NewRegExpTree_WordBoundary(ch), nil
	} else if (ch == 'p' || ch == 'P') && p.hasUFlag {
		// handle ES2018 unicode property tests, e.g.
		// /\p{ASCII_Hex_Digit=true}/ only hex digits
		// /\P{Script=Greek}/  no greek letters
		negated := ch == 'P'
		p.pos++
		if p.pos < p.limit && p.codeUnitAt(p.pos) == '{' {
			lhsStart := p.pos + 1
			for {
				p.pos++
				if !(p.pos < p.limit) {
					break
				}
				ch = p.codeUnitAt(p.pos)
				if !isUnicodePropertyIdentChar(ch) {
					break
				}
			}
			lhs, err := p.substring(lhsStart, p.pos)
			if err != nil {
				return nil, err
			}
			if p.pos < p.limit && ch == '}' {
				// Case of shorthand like /\p{ASCII_Hex_Digit}/u
				p.pos++
				return newUnicodePropertyEscape("", false, lhs, negated)
			} else if p.pos < p.limit && ch == '=' {
				// Case of having '=' like /\p{Script=Greek}/u
				rhsStart := p.pos + 1
				for {
					p.pos++
					if !(p.pos < p.limit) {
						break
					}
					ch = p.codeUnitAt(p.pos)
					if !isUnicodePropertyIdentChar(ch) {
						break
					}
				}
				if p.pos < p.limit && ch == '}' {
					rhs, err := p.substring(rhsStart, p.pos)
					if err != nil {
						return nil, err
					}
					p.pos++
					return newUnicodePropertyEscape(lhs, true, rhs, negated)
				}
				return nil, p.sliceError("Malformed Unicode Property Escape: expected '}' after ", start, p.pos)
			}
			return nil, p.sliceError("Malformed Unicode Property Escape: expected '=' or '}' after ", start, p.pos)
		}
		return nil, p.sliceError("Malformed Unicode Property Escape: expected '{' after ", start, p.pos)
	} else if '1' <= ch && ch <= '9' {
		p.pos++
		possibleGroupIndex := int(ch - '0')
		if p.numCapturingGroups >= possibleGroupIndex {
			if p.pos < p.limit {
				next := p.codeUnitAt(p.pos)
				if '0' <= next && next <= '9' {
					twoDigitGroupIndex := possibleGroupIndex*10 + int(next-'0')
					if p.numCapturingGroups >= twoDigitGroupIndex {
						p.pos++
						possibleGroupIndex = twoDigitGroupIndex
					}
				}
			}
			return NewRegExpTree_BackReference(possibleGroupIndex), nil
		}
		// \1 - \7 are octal escapes if there is no such group.
		// \8 and \9 are the literal characters '8' and '9' if there
		// is no such group.
		c := ch
		if possibleGroupIndex <= 7 {
			c = uint16(possibleGroupIndex)
		}
		return NewRegExpTree_Text(stringFromUTF16Char(c)), nil
	} else if p.lookForNamedCaptureBackreferences &&
		ch == 'k' &&
		p.pos+1 < p.limit &&
		p.codeUnitAt(p.pos+1) == '<' &&
		// According to the spec
		// https://github.com/tc39/proposal-regexp-named-groups#backwards-compatibility-of-new-syntax
		// we want to treat \k as a normal string if there are no named capturing groups present
		len(p.capturingGroupNames) != 0 {
		p.pos += 2
		potentialName, err := p.scanNamedGroupName()
		if err != nil {
			return nil, err
		}
		if _, ok := p.capturingGroupNames[potentialName]; !ok {
			return nil, p.tailError("Invalid named capture referenced: ", start)
		}
		return NewRegExpTree_NamedBackReference(potentialName), nil
	}
	charGroup := regExpTree_namedCharGroup(ch)
	if charGroup != nil { // Handle \d, etc.
		p.pos++
		return NewRegExpTree_Charset(charGroup, CharRanges_EMPTY), nil
	}
	cp, err := p.parseEscapeChar()
	if err != nil {
		return nil, err
	}
	return NewRegExpTree_Text(stringFromCodePoint(cp)), nil
}

// parseHex parses n hex digits to a code-unit.
//
// Java: Parser.parseHex(int)
func (p *regExpParser) parseHex(n int) (int, *RegExpParseError) {
	if p.pos+n > p.limit {
		return 0, p.tailError("Abbreviated hex escape ", p.pos)
	}
	if n > 7 {
		// We need to guard the MSB to prevent overflow.
		return 0, p.sliceError("Cannot parse hexadecimal encoding wider than 28 bits: ", p.pos, p.pos+n)
	}

	result := 0
	for n--; n >= 0; n-- {
		ch := p.codeUnitAt(p.pos)
		var digit int
		if '0' <= ch && ch <= '9' {
			digit = int(ch - '0')
		} else if 'a' <= ch && ch <= 'f' {
			digit = int(ch) + (10 - 'a')
		} else if 'A' <= ch && ch <= 'F' {
			digit = int(ch) + (10 - 'A')
		} else {
			return 0, p.tailError("", p.pos)
		}
		p.pos++
		result = (result << 4) | digit
	}
	return result, nil
}

// parseBracedUnicodeEscape parses \u{...} (with the 'u' flag).
//
// Java: Parser.parseBracedUnicodeEscape()
func (p *regExpParser) parseBracedUnicodeEscape() (int, *RegExpParseError) {
	openBrace := p.pos
	p.pos++

	closeBrace := p.pos
	for closeBrace < p.limit && p.codeUnitAt(closeBrace) != '}' {
		closeBrace++
	}
	if closeBrace == p.limit {
		return 0, p.tailError("Malformed unicode escape: expected '}' after ", openBrace)
	} else if closeBrace == p.pos {
		return 0, illegalArgError("Empty unicode escape")
	}

	result, err := p.parseHex(closeBrace - p.pos)
	if err != nil {
		return 0, err
	}
	if result > 0x10FFFF {
		return 0, p.sliceError("Unicode must be at most 0x10FFFF: ", openBrace+1, p.pos)
	}
	p.pos++ // Consume the close brace.
	return result, nil
}

func isRepetitionStart(ch uint16) bool {
	switch ch {
	case '?', '*', '+', '{':
		return true
	default:
		return false
	}
}

// parseRepetition parses a repetition. x? is treated as a repetition -- an
// optional production can be matched 0 or 1 time.
//
// Java: Parser.parseRepetition(RegExpTree)
func (p *regExpParser) parseRepetition(body RegExpTree) RegExpTree {
	if p.pos >= p.limit {
		return body
	}
	var min_ int
	var max_ int
	switch p.codeUnitAt(p.pos) {
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
			if p.codeUnitAt(i) == '}' {
				end = i
				break
			}
		}
		if end < 0 {
			p.pos = start - 1
			return body
		}
		p.pos = end + 1
		comma := -1
		for i := start; i < end; i++ {
			if p.codeUnitAt(i) == ',' {
				comma = i
				break
			}
		}
		var ok bool
		if comma >= 0 {
			min_, ok = p.parseDecimalInt(start, comma)
			if ok {
				if comma+1 != end {
					max_, ok = p.parseDecimalInt(comma+1, end)
				} else {
					max_ = javaIntegerMaxValue
				}
			}
		} else {
			min_, ok = p.parseDecimalInt(start, end)
			max_ = min_
		}
		if !ok || min_ < 0 || min_ > max_ {
			// Treat the open curly bracket literally.
			p.pos = start - 1
			return body
		}
	default:
		return body
	}
	greedy := true
	if p.pos < p.limit && p.codeUnitAt(p.pos) == '?' {
		greedy = false
		p.pos++
	}
	return NewRegExpTree_Repetition(body, min_, max_, greedy)
}

// parseDecimalInt parses a signed 32-bit decimal integer from pattern[start:end],
// matching Java Integer.parseInt(s, 10) without exceptions.
func (p *regExpParser) parseDecimalInt(start, end int) (int, bool) {
	i := start
	neg := false
	if i < end && (p.codeUnitAt(i) == '+' || p.codeUnitAt(i) == '-') {
		neg = p.codeUnitAt(i) == '-'
		i++
	}
	if i == end {
		return 0, false
	}
	limit := int64(javaIntegerMaxValue)
	if neg {
		limit++
	}
	var v int64
	for ; i < end; i++ {
		c := p.codeUnitAt(i)
		var d int
		if c >= '0' && c <= '9' {
			d = int(c - '0')
		} else if c < 0x80 {
			return 0, false
		} else {
			d = jutil.CharacterDigit(c, 10)
			if d < 0 {
				return 0, false
			}
		}
		v = v*10 + int64(d)
		if v > limit {
			return 0, false
		}
	}
	if neg {
		return int(-v), true
	}
	return int(v), true
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
