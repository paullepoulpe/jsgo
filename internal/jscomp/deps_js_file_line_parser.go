// Port of com/google/javascript/jscomp/deps/JsFileLineParser.java @ closure-compiler v20261021

/*
 * Copyright 2008 The Closure Compiler Authors.
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

package jscomp

import (
	"strings"
	"unicode/utf8"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// jsLineHandler is implemented by concrete line-based JS/deps parsers.
type jsLineHandler interface {
	parseLine(line string) (keepGoing bool, errMsg string)
	parseJsDocCommentLine(line string) bool
}

// jsFileLineParser strips JS comments line-by-line and parses string, array,
// and object literals for JsFileRegexParser and DepsFileRegexParser.
//
// Go-only: language-adapter — line and string positions are byte offsets
// into the (WTF-8) Go strings: every index this struct computes comes from
// finding an ASCII delimiter ("//", "/*", "*/", quotes, backslashes, ':' and
// ',') or from stepping over one, so byte arithmetic selects exactly the
// characters Java's UTF-16 arithmetic selects.
type jsFileLineParser struct {
	shortcutMode   bool
	valueMatcher   jsFileLineParser_stringLiteralMatcher
	filePath       string
	lineNum        int
	errorManager   ErrorManager
	parseSucceeded bool
}

// Java: JsFileLineParser.PARSE_WARNING
var JsFileLineParser_PARSE_WARNING = DiagnosticType_Warning("DEPS_PARSE_WARNING", "{0}\n{1}")

// Java: JsFileLineParser.PARSE_ERROR
var JsFileLineParser_PARSE_ERROR = DiagnosticType_Error("DEPS_PARSE_ERROR", "{0}\n{1}")

// JsFileLineParser_ParseException signifies a problem parsing a JS literal
// (retained for direct ParseJs* test/oracle entrypoints; production line
// parsing returns error strings without panicking).
type JsFileLineParser_ParseException struct {
	jutil.JavaException
}

func newJsFileLineParser_ParseException(message string) *JsFileLineParser_ParseException {
	e := &JsFileLineParser_ParseException{}
	jutil.InitThrowable(e, "com.google.javascript.jscomp.deps.JsFileLineParser$ParseException", &message, nil)
	return e
}

// Java: JsFileLineParser.setShortcutMode(boolean)
func (b *jsFileLineParser) SetShortcutMode(mode bool) {
	b.shortcutMode = mode
}

// Java: JsFileLineParser.didParseSucceed()
func (b *jsFileLineParser) DidParseSucceed() bool {
	return b.parseSucceeded
}

func (*jsFileLineParser) parseJsDocCommentLine(string) bool {
	return true
}

func readJavaReaderToString(r jutil.JavaReader) (text string, err error) {
	defer func() {
		if p := recover(); p != nil {
			if e, ok := p.(error); ok {
				if _, isIO := p.(jutil.IIOException); isIO {
					err = e
					return
				}
			}
			panic(p)
		}
	}()
	return jutil.CharStreams_ToString(r), nil
}

func parseJsReader[H jsLineHandler](b *jsFileLineParser, h H, filePath string, fileContents jutil.JavaReader) {
	b.filePath = filePath
	b.parseSucceeded = true
	b.lineNum = 0

	text, err := readJavaReaderToString(fileContents)
	if err != nil {
		fp := filePath
		b.errorManager.Report(CheckLevel_ERROR.Ptr(),
			JSError_MakeWithSourceLocation(&fp, 0, 0,
				JsFileLineParser_PARSE_ERROR, "Error reading file: "+filePath))
		b.parseSucceeded = false
		return
	}
	parseJsLines(b, h, filePath, text)
}

func parseJsString[H jsLineHandler](b *jsFileLineParser, h H, filePath string, text string) {
	b.filePath = filePath
	b.parseSucceeded = true
	b.lineNum = 0
	parseJsLines(b, h, filePath, text)
}

func parseJsLines[H jsLineHandler](b *jsFileLineParser, h H, filePath string, text string) {
	inMultilineComment := false
	inJsDocComment := false
	pos := 0
	for {
		line, nextPos, ok := jutil.NextLineFromString(text, pos)
		if !ok {
			return
		}
		pos = nextPos
		b.lineNum++

		revisedLine := line
		revisedJsDocCommentLine := ""
		if inMultilineComment {
			endOfComment := strings.Index(revisedLine, "*/")
			if endOfComment != -1 {
				if inJsDocComment {
					revisedJsDocCommentLine = revisedLine[:endOfComment+2]
					inJsDocComment = false
				}
				revisedLine = revisedLine[endOfComment+2:]
				inMultilineComment = false
			} else {
				if inJsDocComment {
					revisedJsDocCommentLine = line
				}
				revisedLine = ""
			}
		}

		if !inMultilineComment {
			for {
				startOfLineComment := strings.Index(revisedLine, "//")
				startOfMultilineComment := strings.Index(revisedLine, "/*")
				if startOfLineComment != -1 &&
					(startOfMultilineComment == -1 ||
						startOfLineComment < startOfMultilineComment) {
					revisedLine = revisedLine[:startOfLineComment]
					break
				} else if startOfMultilineComment != -1 {
					// If comment is in a string (single or double quoted), don't parse as a comment.
					if jsFileLineParser_isCommentQuoted(revisedLine, startOfMultilineComment, '\'') ||
						jsFileLineParser_isCommentQuoted(revisedLine, startOfMultilineComment, '"') {
						break
					}
					if len(revisedLine) > startOfMultilineComment+2 && revisedLine[startOfMultilineComment+2] == '*' {
						inJsDocComment = true
					}
					endOfMultilineComment := jsFileLineParser_indexFrom(revisedLine, "*/", startOfMultilineComment+2)
					if endOfMultilineComment == -1 {
						if inJsDocComment {
							revisedJsDocCommentLine = revisedLine[startOfMultilineComment:]
						}
						revisedLine = revisedLine[:startOfMultilineComment]
						inMultilineComment = true
						break
					}
					if inJsDocComment {
						jsDocComment := revisedLine[startOfMultilineComment : endOfMultilineComment+2]
						if !h.parseJsDocCommentLine(jsDocComment) && b.shortcutMode {
							break
						}
						inJsDocComment = false
					}
					revisedLine = revisedLine[:startOfMultilineComment] +
						revisedLine[endOfMultilineComment+2:]
				} else {
					break
				}
			}
		}

		if revisedJsDocCommentLine != "" {
			if !h.parseJsDocCommentLine(revisedJsDocCommentLine) && b.shortcutMode {
				return
			}
		}

		if revisedLine != "" {
			keepGoing, errMsg := h.parseLine(revisedLine)
			if errMsg != "" {
				fp := filePath
				b.errorManager.Report(
					CheckLevel_ERROR.Ptr(),
					JSError_MakeWithSourceLocation(&fp, b.lineNum, 0,
						JsFileLineParser_PARSE_ERROR, errMsg, line))
				b.parseSucceeded = false
				continue
			}
			if !keepGoing && b.shortcutMode {
				return
			}
		}
	}
}

func jsFileLineParser_indexFrom(s string, sub string, from int) int {
	if from > len(s) {
		return -1
	}
	if i := strings.Index(s[from:], sub); i >= 0 {
		return from + i
	}
	return -1
}

// Java: JsFileLineParser.isCommentQuoted(String, int, char)
func jsFileLineParser_isCommentQuoted(line string, startOfMultilineComment int, quoteChar byte) bool {
	startQuoteIndex := strings.IndexByte(line, quoteChar)
	// Loop in case there are multiple strings between start of line and start of the comment.
	for startQuoteIndex != -1 && startQuoteIndex < startOfMultilineComment {
		closingQuoteIndex := jsFileLineParser_indexByteFrom(line, quoteChar, startQuoteIndex+1)
		if closingQuoteIndex == -1 {
			return false
		}
		// Skip escaped quotes
		for jsFileLineParser_isEscaped(line, closingQuoteIndex) {
			closingQuoteIndex = jsFileLineParser_indexByteFrom(line, quoteChar, closingQuoteIndex+1)
			if closingQuoteIndex == -1 {
				return false
			}
		}
		if closingQuoteIndex > startOfMultilineComment {
			return true
		}
		startQuoteIndex = jsFileLineParser_indexByteFrom(line, quoteChar, closingQuoteIndex+1)
	}
	return false
}

func jsFileLineParser_indexByteFrom(s string, c byte, from int) int {
	if from >= len(s) {
		return -1
	}
	if i := strings.IndexByte(s[from:], c); i >= 0 {
		return from + i
	}
	return -1
}

// Java: JsFileLineParser.isEscaped(String, int)
func jsFileLineParser_isEscaped(line string, closingQuoteIndex int) bool {
	escapeCharCount := 0
	for index := closingQuoteIndex - 1; line[index] == '\\'; index-- {
		escapeCharCount++
	}
	// Odd number of backslashes means one is actually escaping the quote. Otherwise they are
	// literal backslashes which are being escaped by another backslash.
	return escapeCharCount%2 == 1
}

// Java: JsFileLineParser.parseJsString(String)
func (b *jsFileLineParser) ParseJsString(jsStringLiteral string) string {
	val, errMsg, ok := b.parseJsString(jsStringLiteral)
	if !ok {
		panic(newJsFileLineParser_ParseException(errMsg))
	}
	return val
}

func (b *jsFileLineParser) parseJsString(jsStringLiteral string) (string, string, bool) {
	b.valueMatcher.Reset(jsStringLiteral)
	if !b.valueMatcher.Matches() {
		return "", "Syntax error in JS String literal", false
	}
	return b.valueMatcher.MatchedValue(), "", true
}

// Java: JsFileLineParser.parseJsStringArray(String)
func (b *jsFileLineParser) ParseJsStringArray(input string) []string {
	vals, errMsg, ok := b.parseJsStringArray(input)
	if !ok {
		panic(newJsFileLineParser_ParseException(errMsg))
	}
	return vals
}

func (b *jsFileLineParser) parseJsStringArray(input string) ([]string, string, bool) {
	indexStart := strings.IndexByte(input, '[')
	indexEnd := strings.LastIndexByte(input, ']')
	if indexStart == -1 || indexEnd == -1 {
		return nil, "Syntax error when parsing JS array", false
	}
	if indexStart+1 > indexEnd {
		// String.substring(begin, end) with begin > end.
		panic(jutil.NewStringIndexOutOfBoundsExceptionForRange(
			jutil.StringLength(input[:indexStart+1]), jutil.StringLength(input[:indexEnd]), jutil.StringLength(input)))
	}
	innerValues := input[indexStart+1 : indexEnd]
	results := []string{}
	if jsFileLineParser_hasNonControlSpace(innerValues) {
		b.valueMatcher.Reset(innerValues)
		for {
			if !b.valueMatcher.LookingAt() {
				return nil, "Syntax error in JS String literal", false
			}
			results = append(results, b.valueMatcher.MatchedValue())
			if b.valueMatcher.HitEnd() {
				break
			}
			if innerValues[b.valueMatcher.End()] != ',' {
				return nil, "Missing comma in string array", false
			}
			b.valueMatcher.Region(b.valueMatcher.End()+1, b.valueMatcher.RegionEnd())
		}
	}
	return results, "", true
}

// Java: JsFileLineParser.parseJsStringMap(String)
func (b *jsFileLineParser) ParseJsStringMap(input string) *jutil.LinkedHashMap[string, string] {
	m, errMsg, ok := b.parseJsStringMap(input)
	if !ok {
		panic(newJsFileLineParser_ParseException(errMsg))
	}
	return m
}

func (b *jsFileLineParser) parseJsStringMap(input string) (*jutil.LinkedHashMap[string, string], string, bool) {
	input = jsFileLineParser_trimCharMatcherWhitespace(input)
	if input == "" || input[0] != '{' || input[len(input)-1] != '}' {
		return nil, "Syntax error when parsing JS object", false
	}
	input = jsFileLineParser_trimJava(input[1 : len(input)-1])

	results := jutil.NewLinkedHashMap[string, string]()
	b.valueMatcher.Reset(input)
	for input != "" {
		// Parse the next key (TODO(sdh): need to support non-quoted keys?).
		if !b.valueMatcher.LookingAt() {
			return nil, "Bad key in JS object literal", false
		}
		key := b.valueMatcher.MatchedValue()
		if b.valueMatcher.HitEnd() {
			return nil, "Missing value in JS object literal", false
		}
		if input[b.valueMatcher.End()] != ':' {
			return nil, "Missing colon in JS object literal", false
		}
		b.valueMatcher.Region(b.valueMatcher.End()+1, b.valueMatcher.RegionEnd())

		// Parse the corresponding value.
		if !b.valueMatcher.LookingAt() {
			return nil, "Bad value in JS object literal", false
		}
		val := b.valueMatcher.MatchedValue()
		results.Put(key, val)
		if b.valueMatcher.HitEnd() {
			break
		}
		if input[b.valueMatcher.End()] != ',' {
			return nil, "Missing comma in JS object literal", false
		}
		b.valueMatcher.Region(b.valueMatcher.End()+1, b.valueMatcher.RegionEnd())
	}
	return results, "", true
}

func jsFileLineParser_hasNonControlSpace(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > ' ' {
			return true
		}
	}
	return false
}

// jsFileLineParser_trimJava is Java's String.trim() without UTF-16 allocation:
// every UTF-16 code unit <= 0x20 is an ASCII byte <= 0x20 in WTF-8.
func jsFileLineParser_trimJava(s string) string {
	start := 0
	for start < len(s) && s[start] <= ' ' {
		start++
	}
	end := len(s)
	for end > start && s[end-1] <= ' ' {
		end--
	}
	return s[start:end]
}

func jsFileLineParser_isCharMatcherSpace(r rune, size int) bool {
	if r < utf8.RuneSelf {
		return r == ' ' || (r >= 0x09 && r <= 0x0D)
	}
	if (r == utf8.RuneError && size == 1) || r > 0xFFFF {
		return false
	}
	return jutil.CharMatcher_Whitespace_Matches(uint16(r))
}

func jsFileLineParser_trimCharMatcherWhitespace(s string) string {
	start := 0
	for start < len(s) {
		c := s[start]
		if c < utf8.RuneSelf {
			if c != ' ' && (c < 0x09 || c > 0x0D) {
				break
			}
			start++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[start:])
		if !jsFileLineParser_isCharMatcherSpace(r, size) {
			break
		}
		start += size
	}
	end := len(s)
	for end > start {
		c := s[end-1]
		if c < utf8.RuneSelf {
			if c != ' ' && (c < 0x09 || c > 0x0D) {
				break
			}
			end--
			continue
		}
		r, size := utf8.DecodeLastRuneInString(s[start:end])
		if !jsFileLineParser_isCharMatcherSpace(r, size) {
			break
		}
		end -= size
	}
	return s[start:end]
}

// jsFileLineParser_stringLiteralMatcher is the java.util.regex.Matcher of
// STRING_LITERAL_PATTERN:
//
//	\s*(?:'((?:\\.|[^'\\])*)'|"((?:\\.|[^"\\])*)")\s*
//
// Go-only: language-adapter — Java regex. JsFileLineParser relies on
// Matcher.hitEnd() and Matcher.region(), which RE2 cannot provide, so this
// adapter runs the pattern with Java's backtracking order (greedy \s*, the
// first alternative before the second, the greedy loop trying the escape
// before any other character, then the closing quote) and records hitEnd
// whenever Java's engine would read at the region end. \s is Java's
// [ \t\n\x0B\f\r], and . any character but a line terminator (\n, \r,
// U+0085, U+2028, U+2029). Positions are byte offsets: apart from those
// line terminators the pattern only ever compares ASCII characters, and a
// region always starts and ends next to one, so byte and UTF-16 matching
// agree.
type jsFileLineParser_stringLiteralMatcher struct {
	input      string
	from       int // region start
	to         int // region end
	hitEnd     bool
	matched    bool
	end        int
	g1Start    int
	g1End      int
	g2Start    int
	g2End      int
	requireEnd bool
}

// Reset is Matcher.reset(CharSequence): new input, whole-input region.
func (m *jsFileLineParser_stringLiteralMatcher) Reset(input string) {
	m.input = input
	m.from = 0
	m.to = len(input)
	m.matched = false
	m.hitEnd = false
}

// Region is Matcher.region(int, int): it resets the matcher and limits it
// to [start, end).
func (m *jsFileLineParser_stringLiteralMatcher) Region(start int, end int) {
	if start < 0 || start > len(m.input) {
		panic(jutil.NewIndexOutOfBoundsException("start"))
	}
	if end < 0 || end > len(m.input) {
		panic(jutil.NewIndexOutOfBoundsException("end"))
	}
	if start > end {
		panic(jutil.NewIndexOutOfBoundsException("start > end"))
	}
	m.matched = false
	m.hitEnd = false
	m.from = start
	m.to = end
}

// RegionEnd is Matcher.regionEnd().
func (m *jsFileLineParser_stringLiteralMatcher) RegionEnd() int { return m.to }

// HitEnd is Matcher.hitEnd().
func (m *jsFileLineParser_stringLiteralMatcher) HitEnd() bool { return m.hitEnd }

// Matches is Matcher.matches(): the whole region must match.
func (m *jsFileLineParser_stringLiteralMatcher) Matches() bool {
	return m.match(true)
}

// LookingAt is Matcher.lookingAt(): a match must start at the region start.
func (m *jsFileLineParser_stringLiteralMatcher) LookingAt() bool {
	return m.match(false)
}

// End is Matcher.end(): it panics with IllegalStateException when there is
// no match.
func (m *jsFileLineParser_stringLiteralMatcher) End() int {
	if !m.matched {
		panic(jutil.NewIllegalStateExceptionMsg("No match available"))
	}
	return m.end
}

// MatchedValue returns the matched string from group 1 or group 2 without allocating.
func (m *jsFileLineParser_stringLiteralMatcher) MatchedValue() string {
	if !m.matched {
		panic(jutil.NewIllegalStateExceptionMsg("No match found"))
	}
	if m.g1Start >= 0 {
		return m.input[m.g1Start:m.g1End]
	}
	return m.input[m.g2Start:m.g2End]
}

func jsFileLineParser_isJavaSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == 0x0B || c == '\f' || c == '\r'
}

func (m *jsFileLineParser_stringLiteralMatcher) match(requireEnd bool) bool {
	m.hitEnd = false
	m.matched = false
	m.requireEnd = requireEnd
	m.g1Start, m.g1End, m.g2Start, m.g2End = -1, -1, -1, -1
	// \s* (greedy). Backing off cannot help: the next element is a quote.
	i := m.from
	for i < m.to && jsFileLineParser_isJavaSpace(m.input[i]) {
		i++
	}
	if i >= m.to {
		m.hitEnd = true
	}
	for p := i; p >= m.from; p-- {
		if m.quoted(p, '\'', 1) || m.quoted(p, '"', 2) {
			return true
		}
	}
	return false
}

func (m *jsFileLineParser_stringLiteralMatcher) quoted(p int, quote byte, group int) bool {
	if p >= m.to {
		m.hitEnd = true
		return false
	}
	if m.input[p] != quote {
		return false
	}
	return m.greedyLoop(p+1, p+1, quote, group)
}

func (m *jsFileLineParser_stringLiteralMatcher) greedyLoop(groupStart int, q int, quote byte, group int) bool {
	for {
		if q >= m.to {
			m.hitEnd = true
			break
		}
		c := m.input[q]
		if c == '\\' {
			if q+1 >= m.to {
				m.hitEnd = true
				break
			}
			if jsFileLineParser_isJavaLineTerminatorAt(m.input, q+1) {
				break
			}
			q += 2
			continue
		}
		if c == quote {
			break
		}
		q++
	}
	return m.closing(groupStart, q, quote, group)
}

func jsFileLineParser_isJavaLineTerminatorAt(s string, i int) bool {
	switch s[i] {
	case '\n', '\r':
		return true
	case 0xC2: // U+0085 is C2 85.
		return i+1 < len(s) && s[i+1] == 0x85
	case 0xE2: // U+2028 and U+2029 are E2 80 A8 and E2 80 A9.
		return i+2 < len(s) && s[i+1] == 0x80 && (s[i+2] == 0xA8 || s[i+2] == 0xA9)
	}
	return false
}

func (m *jsFileLineParser_stringLiteralMatcher) closing(groupStart int, q int, quote byte, group int) bool {
	if q >= m.to {
		m.hitEnd = true
		return false
	}
	if m.input[q] != quote {
		return false
	}
	i := q + 1
	for i < m.to && jsFileLineParser_isJavaSpace(m.input[i]) {
		i++
	}
	if i >= m.to {
		m.hitEnd = true
	}
	if m.requireEnd && i != m.to {
		return false
	}
	m.matched = true
	m.end = i
	if group == 1 {
		m.g1Start, m.g1End = groupStart, q
	} else {
		m.g2Start, m.g2End = groupStart, q
	}
	return true
}
