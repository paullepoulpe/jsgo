// Port of com/google/javascript/jscomp/deps/JsFileLineParser.java @ closure-compiler v20261015

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

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// JsFileLineParser is the base class for classes that parse JavaScript
// sources on a line-by-line basis. Strips comments from files and records
// all parsing errors.
//
// The Go interface declares the overridable methods; implementations embed
// JsFileLineParserBase (PORTING.md §4).
//
// Go-only: language-adapter — line and string positions are byte offsets
// into the (WTF-8) Go strings: every index this class computes comes from
// finding an ASCII delimiter ("//", "/*", "*/", quotes, backslashes, ':' and
// ',') or from stepping over one, so byte arithmetic selects exactly the
// characters Java's UTF-16 arithmetic selects.
type JsFileLineParser interface {
	GetJsFileLineParserBase() *JsFileLineParserBase

	// SetShortcutMode: in shortcut mode, the file line parser can stop
	// reading early if it thinks it found enough information.
	//
	// For example, many parsers assume that dependency information never
	// shows up after "real" code.
	//
	// Java: JsFileLineParser.setShortcutMode(boolean)
	SetShortcutMode(mode bool)

	// Java: JsFileLineParser.didParseSucceed()
	DidParseSucceed() bool

	// DoParse performs the line-by-line parsing of the given fileContents.
	// This method strips out JavaScript comments and then uses the abstract
	// parseLine() method to do the line parsing.
	//
	// Java: JsFileLineParser.doParse(String, Reader)
	DoParse(filePath string, fileContents jutil.JavaReader)

	// ParseLine is called for each line of the file being parsed. It
	// returns true to keep going, false otherwise. It panics with
	// *JsFileLineParser_ParseException to signify a problem with the line.
	//
	// Java: JsFileLineParser.parseLine(String)
	ParseLine(line string) bool

	// ParseJsDocCommentLine is called for each JSDoc line of the file being
	// parsed. It returns true to keep going, false otherwise.
	//
	// Java: JsFileLineParser.parseJsDocCommentLine(String)
	ParseJsDocCommentLine(line string) bool

	// Java: JsFileLineParser.parseJsString(String)
	ParseJsString(jsStringLiteral string) string

	// Java: JsFileLineParser.parseJsStringArray(String)
	ParseJsStringArray(input string) []string

	// Java: JsFileLineParser.parseJsStringMap(String)
	ParseJsStringMap(input string) *jutil.LinkedHashMap[string, string]
}

// JsFileLineParserBase holds the fields and base implementations of
// JsFileLineParser.
type JsFileLineParserBase struct {
	self JsFileLineParser

	shortcutMode bool

	// Matcher used in the parsing string literals.
	valueMatcher jsFileLineParser_stringLiteralMatcher

	// Path of the file currently being parsed.
	filePath string
	// The line number of the line currently being parsed.
	lineNum int
	// Handles error messages.
	errorManager ErrorManager
	// Did our parse succeed.
	parseSucceeded bool
}

// Java: JsFileLineParser.PARSE_WARNING
var JsFileLineParser_PARSE_WARNING = DiagnosticType_Warning("DEPS_PARSE_WARNING", "{0}\n{1}")

// Java: JsFileLineParser.PARSE_ERROR
var JsFileLineParser_PARSE_ERROR = DiagnosticType_Error("DEPS_PARSE_ERROR", "{0}\n{1}")

// JsFileLineParser_ParseException is thrown by base classes to signify a
// problem parsing a line.
type JsFileLineParser_ParseException struct {
	jutil.JavaException
	fatal bool
}

func (*JsFileLineParser_ParseException) parseException() {}

// IJsFileLineParser_ParseException matches JsFileLineParser_ParseException.
type IJsFileLineParser_ParseException interface {
	jutil.IException
	parseException()
	IsFatal() bool
}

// Java: JsFileLineParser.ParseException.serialVersionUID
const JsFileLineParser_ParseException_serialVersionUID int64 = 1

// NewJsFileLineParser_ParseException is the constructor: message is a
// description of what caused the exception; fatal is whether the exception
// is recoverable.
//
// Java: JsFileLineParser.ParseException(String, boolean)
func NewJsFileLineParser_ParseException(message string, fatal bool) *JsFileLineParser_ParseException {
	e := &JsFileLineParser_ParseException{fatal: fatal}
	jutil.InitThrowable(e, "com.google.javascript.jscomp.deps.JsFileLineParser$ParseException", &message, nil)
	return e
}

// Java: JsFileLineParser.ParseException.isFatal()
func (e *JsFileLineParser_ParseException) IsFatal() bool {
	return e.fatal
}

// InitJsFileLineParserBase is the constructor; errorManager is the parse
// error handler.
//
// Java: JsFileLineParser(ErrorManager)
func (b *JsFileLineParserBase) InitJsFileLineParserBase(self JsFileLineParser, errorManager ErrorManager) {
	b.self = self
	b.valueMatcher.Reset("")
	b.errorManager = errorManager
}

// GetJsFileLineParserBase gives access to the fields of the base class.
func (b *JsFileLineParserBase) GetJsFileLineParserBase() *JsFileLineParserBase { return b }

// Java: JsFileLineParser.setShortcutMode(boolean)
func (b *JsFileLineParserBase) SetShortcutMode(mode bool) {
	b.shortcutMode = mode
}

// Java: JsFileLineParser.didParseSucceed()
func (b *JsFileLineParserBase) DidParseSucceed() bool {
	return b.parseSucceeded
}

// DoParse performs the line-by-line parsing of the given fileContents. This
// method strips out JavaScript comments and then uses the abstract
// parseLine() method to do the line parsing.
//
// filePath is the path to the file being parsed. Used for reporting parse
// exceptions. fileContents is a reader for the contents of the file.
//
// Java: JsFileLineParser.doParse(String, Reader)
func (b *JsFileLineParserBase) DoParse(filePath string, fileContents jutil.JavaReader) {
	b.filePath = filePath
	b.parseSucceeded = true
	b.lineNum = 0

	var text string
	readOK := true
	jutil.Try(func() {
		text = jutil.CharStreams_ToString(fileContents)
	}, func(e jutil.IIOException) {
		fp := filePath
		b.errorManager.Report(CheckLevel_ERROR.Ptr(),
			JSError_MakeWithSourceLocation(&fp, 0, 0, /* char offset */
				JsFileLineParser_PARSE_ERROR, "Error reading file: "+filePath))
		b.parseSucceeded = false
		readOK = false
	})
	if !readOK {
		return
	}
	b.doParseLoadedString(filePath, text)
}

func (b *JsFileLineParserBase) doParseString(filePath string, text string) {
	b.filePath = filePath
	b.parseSucceeded = true
	b.lineNum = 0
	b.doParseLoadedString(filePath, text)
}

func (b *JsFileLineParserBase) doParseLoadedString(filePath string, text string) {
	inMultilineComment := false
	inJsDocComment := false
	pos := 0
	for pos < len(text) {
		if b.doParseLinesChunk(filePath, text, &pos, &inMultilineComment, &inJsDocComment) {
			break
		}
	}
}

func (b *JsFileLineParserBase) doParseLinesChunk(
	filePath string,
	text string,
	pos *int,
	inMultilineComment *bool,
	inJsDocComment *bool,
) (stopped bool) {
	var curLine string
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(IJsFileLineParser_ParseException); ok {
				level := CheckLevel_WARNING
				diagnostic := JsFileLineParser_PARSE_WARNING
				if e.IsFatal() {
					level = CheckLevel_ERROR
					diagnostic = JsFileLineParser_PARSE_ERROR
				}
				fp := filePath
				b.errorManager.Report(
					level.Ptr(),
					JSError_MakeWithSourceLocation(&fp, b.lineNum, 0, /* char offset */
						diagnostic,
						e.GetMessage(), curLine))
				b.parseSucceeded = b.parseSucceeded && !e.IsFatal()
				stopped = false
				return
			}
			panic(r)
		}
	}()

	for {
		line, nextPos, ok := jutil.NextLineFromString(text, *pos)
		if !ok {
			return true
		}
		*pos = nextPos
		curLine = line
		b.lineNum++

		revisedLine := line
		revisedJsDocCommentLine := ""
		if *inMultilineComment {
			endOfComment := strings.Index(revisedLine, "*/")
			if endOfComment != -1 {
				if *inJsDocComment {
					revisedJsDocCommentLine = revisedLine[:endOfComment+2]
					*inJsDocComment = false
				}
				revisedLine = revisedLine[endOfComment+2:]
				*inMultilineComment = false
			} else {
				if *inJsDocComment {
					revisedJsDocCommentLine = line
				}
				revisedLine = ""
			}
		}

		if !*inMultilineComment {
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
					if b.isCommentQuoted(revisedLine, startOfMultilineComment, '\'') {
						break
					}
					if b.isCommentQuoted(revisedLine, startOfMultilineComment, '"') {
						break
					}
					if len(revisedLine) > startOfMultilineComment+2 && revisedLine[startOfMultilineComment+2] == '*' {
						*inJsDocComment = true
					}
					endOfMultilineComment := jsFileLineParser_indexFrom(revisedLine, "*/", startOfMultilineComment+2)
					if endOfMultilineComment == -1 {
						if *inJsDocComment {
							revisedJsDocCommentLine = revisedLine[startOfMultilineComment:]
						}
						revisedLine = revisedLine[:startOfMultilineComment]
						*inMultilineComment = true
						break
					} else {
						if *inJsDocComment {
							jsDocComment := revisedLine[startOfMultilineComment : endOfMultilineComment+2]
							if !b.self.ParseJsDocCommentLine(jsDocComment) && b.shortcutMode {
								break
							}
							*inJsDocComment = false
						}
						revisedLine = revisedLine[:startOfMultilineComment] +
							revisedLine[endOfMultilineComment+2:]
					}
				} else {
					break
				}
			}
		}

		if revisedJsDocCommentLine != "" {
			if !b.self.ParseJsDocCommentLine(revisedJsDocCommentLine) && b.shortcutMode {
				return true
			}
		}

		if revisedLine != "" {
			// This check for shortcut mode should be redundant, but
			// it's done for safety reasons.
			if !b.self.ParseLine(revisedLine) && b.shortcutMode {
				return true
			}
		}
	}
}

// jsFileLineParser_newBufferedReader is new BufferedReader(fileContents).
//
// Go-only: language-adapter — jutil.BufferedReader reads from a Go string,
// so the Reader is drained first (an IOException surfaces from the first
// read, as from Java's first readLine()). It is drained with
// CharStreams_ToString, which returns an unread StringReader's string
// without converting it to UTF-16 and back.
func jsFileLineParser_newBufferedReader(fileContents jutil.JavaReader) *jutil.BufferedReader {
	return jutil.NewBufferedReaderFromString(jutil.CharStreams_ToString(fileContents))
}

// jsFileLineParser_indexFrom is String.indexOf(String, int) on byte offsets.
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
func (b *JsFileLineParserBase) isCommentQuoted(line string, startOfMultilineComment int, quoteChar byte) bool {
	startQuoteIndex := strings.IndexByte(line, quoteChar)
	// Loop in case there are multiple strings between start of line and start of the comment.
	for startQuoteIndex != -1 && startQuoteIndex < startOfMultilineComment {
		closingQuoteIndex := jsFileLineParser_indexByteFrom(line, quoteChar, startQuoteIndex+1)
		if closingQuoteIndex == -1 {
			return false
		}
		// Skip escaped quotes
		for b.isEscaped(line, closingQuoteIndex) {
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

// jsFileLineParser_indexByteFrom is String.indexOf(int, int) on byte
// offsets.
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
func (b *JsFileLineParserBase) isEscaped(line string, closingQuoteIndex int) bool {
	escapeCharCount := 0
	index := closingQuoteIndex - 1
	for line[index] == '\\' {
		index--
		escapeCharCount++
	}
	// Odd number of backslashes means one is actually escaping the quote. Otherwise they are
	// literal backslashes which are being escaped by another backslash.
	return escapeCharCount%2 == 1
}

// ParseJsDocCommentLine is called for each JSDoc line of the file being
// parsed: line is the JSDoc comment line to parse. It returns true to keep
// going, false otherwise.
//
// Java: JsFileLineParser.parseJsDocCommentLine(String)
func (b *JsFileLineParserBase) ParseJsDocCommentLine(line string) bool {
	return true
}

// ParseJsString parses a JS string literal.
//
// jsStringLiteral is the literal. Must look like "asdf" or 'asdf'. It panics
// with *JsFileLineParser_ParseException if there is a string literal that
// cannot be parsed.
//
// Java: JsFileLineParser.parseJsString(String)
func (b *JsFileLineParserBase) ParseJsString(jsStringLiteral string) string {
	b.valueMatcher.Reset(jsStringLiteral)
	if !b.valueMatcher.Matches() {
		panic(NewJsFileLineParser_ParseException("Syntax error in JS String literal", true /* fatal */))
	}
	return b.valueMatcher.MatchedValue()
}

// ParseJsStringArray parses a JavaScript array of string literals. (eg:
// ['a', 'b', "c"]).
//
// input is a string containing a JavaScript array of string literals; the
// result is a list of parsed string literals. It panics with
// *JsFileLineParser_ParseException if there is a syntax error with the
// input.
//
// Java: JsFileLineParser.parseJsStringArray(String)
func (b *JsFileLineParserBase) ParseJsStringArray(input string) []string {
	results := []string{}
	indexStart := strings.IndexByte(input, '[')
	indexEnd := strings.LastIndexByte(input, ']')
	if (indexStart == -1) || (indexEnd == -1) {
		panic(NewJsFileLineParser_ParseException("Syntax error when parsing JS array", true /* fatal */))
	}
	if indexStart+1 > indexEnd {
		// String.substring(begin, end) with begin > end.
		panic(jutil.NewStringIndexOutOfBoundsExceptionForRange(
			jutil.StringLength(input[:indexStart+1]), jutil.StringLength(input[:indexEnd]), jutil.StringLength(input)))
	}
	innerValues := input[indexStart+1 : indexEnd]

	if jutil.StringTrim(innerValues) != "" {
		b.valueMatcher.Reset(innerValues)
		for {
			// Parse the current string literal.
			if !b.valueMatcher.LookingAt() {
				panic(NewJsFileLineParser_ParseException("Syntax error in JS String literal", true /* fatal */))
			}
			// Add it to the results.
			results = append(results, b.valueMatcher.MatchedValue())
			if b.valueMatcher.HitEnd() {
				break
			}
			// Ensure there is a comma after the value.
			if innerValues[b.valueMatcher.End()] != ',' {
				panic(NewJsFileLineParser_ParseException("Missing comma in string array", true /* fatal */))
			}
			// Move to the next value.
			b.valueMatcher.Region(b.valueMatcher.End()+1, b.valueMatcher.RegionEnd())
		}
	}
	return results
}

// TODO(sdh): Consider simplifying this by reusing the parser or a separate JSON library.

// ParseJsStringMap parses a JavaScript map of string literals. (eg: {'a':
// 'b', "c": "d"}).
//
// input is a string containing a JavaScript map of string literals; the
// result is a map of parsed string literals. It panics with
// *JsFileLineParser_ParseException if there is a syntax error with the
// input.
//
// Java: JsFileLineParser.parseJsStringMap(String)
func (b *JsFileLineParserBase) ParseJsStringMap(input string) *jutil.LinkedHashMap[string, string] {
	input = jutil.CharMatcher_Whitespace_TrimFrom(input)
	jsFileLineParser_check(
		input != "" && input[0] == '{' && input[len(input)-1] == '}',
		"Syntax error when parsing JS object")
	input = jutil.StringTrim(input[1 : len(input)-1])

	results := jutil.NewLinkedHashMap[string, string]()
	done := input == ""
	b.valueMatcher.Reset(input)
	for !done {
		// Parse the next key (TODO(sdh): need to support non-quoted keys?).
		jsFileLineParser_check(b.valueMatcher.LookingAt(), "Bad key in JS object literal")
		key := b.valueMatcher.MatchedValue()
		jsFileLineParser_check(!b.valueMatcher.HitEnd(), "Missing value in JS object literal")
		jsFileLineParser_check(input[b.valueMatcher.End()] == ':', "Missing colon in JS object literal")
		b.valueMatcher.Region(b.valueMatcher.End()+1, b.valueMatcher.RegionEnd())

		// Parse the corresponding value.
		jsFileLineParser_check(b.valueMatcher.LookingAt(), "Bad value in JS object literal")
		val := b.valueMatcher.MatchedValue()
		results.Put(key, val)
		if !b.valueMatcher.HitEnd() {
			jsFileLineParser_check(input[b.valueMatcher.End()] == ',', "Missing comma in JS object literal")
			b.valueMatcher.Region(b.valueMatcher.End()+1, b.valueMatcher.RegionEnd())
		} else {
			done = true
		}
	}
	return results
}

// Java: JsFileLineParser.check(boolean, String)
func jsFileLineParser_check(condition bool, message string) {
	if !condition {
		panic(NewJsFileLineParser_ParseException(message, true /* fatal */))
	}
}

// jsFileLineParser_stringLiteralMatcher is the java.util.regex.Matcher of
// STRING_LITERAL_PATTERN:
//
//	\s*(?:'((?:\\'|[^'])*?)'|"((?:\\"|[^"])*?)")\s*
//
// Go-only: language-adapter — Java regex. JsFileLineParser relies on
// Matcher.hitEnd() and Matcher.region(), which RE2 cannot provide, so this
// adapter runs the pattern with Java's backtracking order (greedy \s*, the
// first alternative before the second, the lazy loop trying the closing
// quote before the escape before any other character) and records hitEnd
// whenever Java's engine would read at the region end. \s is Java's
// [ \t\n\x0B\f\r]. Positions are byte offsets: the pattern only ever
// compares ASCII characters, so byte and UTF-16 matching agree.
type jsFileLineParser_stringLiteralMatcher struct {
	input  string
	from   int // region start
	to     int // region end
	hitEnd bool
	// The last successful match: group 0 end and groups 1 and 2 (-1 when
	// the group did not participate).
	matched                bool
	end                    int
	g1Start, g1End         int
	g2Start, g2End         int
	requireEnd             bool
	curStart, curEndOfLoop int
}

func newJsFileLineParser_stringLiteralMatcher(input string) *jsFileLineParser_stringLiteralMatcher {
	m := &jsFileLineParser_stringLiteralMatcher{}
	m.Reset(input)
	return m
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

// Group is Matcher.group(int) for groups 1 and 2: nil when the group did
// not participate in the match.
func (m *jsFileLineParser_stringLiteralMatcher) Group(group int) *string {
	if !m.matched {
		panic(jutil.NewIllegalStateExceptionMsg("No match found"))
	}
	var s, e int
	switch group {
	case 1:
		s, e = m.g1Start, m.g1End
	case 2:
		s, e = m.g2Start, m.g2End
	default:
		panic(jutil.NewIndexOutOfBoundsException("No group " + jutil.ToJavaString(group)))
	}
	if s < 0 {
		return nil
	}
	str := m.input[s:e]
	return &str
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
		// First alternative: '...'
		if m.quoted(p, '\'', 1) {
			return true
		}
		// Second alternative: "..."
		if m.quoted(p, '"', 2) {
			return true
		}
	}
	return false
}

// quoted matches the alternative quote((?:\\quote|[^quote])*?)quote\s* at p.
func (m *jsFileLineParser_stringLiteralMatcher) quoted(p int, quote byte, group int) bool {
	if p >= m.to {
		m.hitEnd = true
		return false
	}
	if m.input[p] != quote {
		return false
	}
	return m.lazyLoop(p+1, p+1, quote, group)
}

// lazyLoop is the lazy loop (?:\\quote|[^quote])*? at q, followed by the
// closing quote, \s* and the end of the pattern.
func (m *jsFileLineParser_stringLiteralMatcher) lazyLoop(groupStart int, q int, quote byte, group int) bool {
	for {
		// Try to finish the match first.
		if m.closing(groupStart, q, quote, group) {
			return true
		}
		// Body alternative 1: the two-character escape.
		if q >= m.to {
			m.hitEnd = true
		} else if m.input[q] == '\\' {
			if q+1 >= m.to {
				m.hitEnd = true
			} else if m.input[q+1] == quote {
				if m.lazyLoop(groupStart, q+2, quote, group) {
					return true
				}
			}
		}
		// Body alternative 2: any other character.
		if q >= m.to {
			m.hitEnd = true
			return false
		}
		if m.input[q] == quote {
			return false
		}
		q++
	}
}

// closing matches quote\s* at q and accepts.
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
	// Greedy \s* backs off one character at a time; only Matcher.matches()
	// rejects an end before the region end.
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
