// Port of com/google/debugging/sourcemap/SourceMapJsonLexer.java @ closure-compiler v20261021

/*
 * Copyright 2026 The Closure Compiler Authors.
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

package sourcemap

import (
	"math"
	"slices"
	"strconv"
	"strings"
	"unicode/utf16"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

const (
	smScopeEmptyArray uint8 = iota + 1
	smScopeNonEmptyArray
	smScopeEmptyObject
	smScopeDanglingName
	smScopeNonEmptyObject
	smScopeEmptyDocument
	smScopeNonEmptyDocument
)

const (
	smPeekNone uint8 = iota
	smPeekBeginObject
	smPeekEndObject
	smPeekBeginArray
	smPeekEndArray
	smPeekTrue
	smPeekFalse
	smPeekNull
	smPeekSingleQuoted
	smPeekDoubleQuoted
	smPeekUnquoted
	smPeekSingleQuotedName
	smPeekDoubleQuotedName
	smPeekUnquotedName
	smPeekLong
	smPeekNumber
	smPeekEOF
)

const (
	smNumNone = iota
	smNumSign
	smNumDigit
	smNumDecimal
	smNumFractionDigit
	smNumExpE
	smNumExpSign
	smNumExpDigit
)

const smMinIncompleteInt64 = -0x8000000000000000 / 10

// smNonLiteralLUT marks the 16 ASCII bytes that terminate unquoted literals in Gson.
var smNonLiteralLUT = [256]bool{
	'/': true, '\\': true, ';': true, '#': true, '=': true,
	'{': true, '}': true, '[': true, ']': true, ':': true, ',': true,
	' ': true, '\t': true, '\f': true, '\r': true, '\n': true,
}

type smFrame struct {
	name    string
	index   int
	scope   uint8
	hasName bool
}

// SourceMapJsonLexer is a streaming JSON lexer operating directly on the UTF-8
// source map string. Unescaped strings (including large Base64 VLQ "mappings"
// payloads) are sliced directly from json without allocating UTF-16 buffers.
type SourceMapJsonLexer struct {
	json     string
	pos      int
	length   int
	bomBytes int
	lenient  bool

	peeked             uint8
	peekedLong         int64
	peekedNumberLength int

	stack    []smFrame
	stackBuf [16]smFrame
}

// NewSourceMapJsonLexer constructs a lexer over json.
//
// Java: SourceMapJsonLexer(String)
func NewSourceMapJsonLexer(json string) *SourceMapJsonLexer {
	return &SourceMapJsonLexer{json: json, length: len(json)}
}

func (l *SourceMapJsonLexer) initStream(json string) {
	*l = SourceMapJsonLexer{
		json:    json,
		length:  len(json),
		lenient: true,
	}
	l.stackBuf[0] = smFrame{scope: smScopeEmptyDocument}
	l.stack = l.stackBuf[:1]
	if strings.HasPrefix(json, "\xef\xbb\xbf") {
		l.pos = 3
		l.bomBytes = 3
	}
}

func (l *SourceMapJsonLexer) cloneState() SourceMapJsonLexer {
	cp := *l
	if len(l.stack) > 0 && &l.stack[0] == &l.stackBuf[0] {
		cp.stack = cp.stackBuf[:len(l.stack)]
	} else {
		cp.stack = slices.Clone(l.stack)
	}
	return cp
}

func (l *SourceMapJsonLexer) topScope() uint8     { return l.stack[len(l.stack)-1].scope }
func (l *SourceMapJsonLexer) setTopScope(s uint8) { l.stack[len(l.stack)-1].scope = s }
func (l *SourceMapJsonLexer) pushScope(s uint8)   { l.stack = append(l.stack, smFrame{scope: s}) }
func (l *SourceMapJsonLexer) incTopIndex()        { l.stack[len(l.stack)-1].index++ }
func (l *SourceMapJsonLexer) setTopName(name string) {
	f := &l.stack[len(l.stack)-1]
	f.name = name
	f.hasName = true
}

// Java: SourceMapJsonLexer.skipWhitespace()
func (l *SourceMapJsonLexer) SkipWhitespace() {
	for l.pos < l.length {
		c := l.json[l.pos]
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			l.pos++
		} else {
			break
		}
	}
}

// Java: SourceMapJsonLexer.nextString()
func (l *SourceMapJsonLexer) NextString() string {
	l.SkipWhitespace()
	if l.pos >= l.length || l.json[l.pos] != '"' {
		panic(NewSourceMapParseException("Expected string"))
	}
	l.pos++
	s, err := l.readQuotedValue('"')
	if err != nil {
		if sm, ok := err.(*SourceMapParseException); ok {
			panic(sm)
		}
		panic(NewSourceMapParseException(err.GetMessage()))
	}
	return s
}

// Java: SourceMapJsonLexer.nextInt()
func (l *SourceMapJsonLexer) NextInt() int {
	l.SkipWhitespace()
	start := l.pos
	if l.pos < l.length && l.json[l.pos] == '-' {
		l.pos++
	}
	for l.pos < l.length && l.json[l.pos] >= '0' && l.json[l.pos] <= '9' {
		l.pos++
	}
	if start == l.pos {
		panic(NewSourceMapParseException("Expected integer"))
	}
	result, ok := sourceMapJsonLexer_parseInt(l.json, start, l.pos, 10)
	if !ok {
		panic(NewSourceMapParseException("Invalid integer: " + l.json[start:l.pos]))
	}
	return int(result)
}

// ---------------------------------------------------------------------------
// Gson 2.9.1 lenient streaming token scanner

func wrapGsonSyntax(causeClass, msg string) *SourceMapParseException {
	return NewSourceMapParseException(
		"JSON parse exception: " + jutil.GsonPackage + ".JsonSyntaxException: " + causeClass + ": " + msg)
}

func wrapGsonSyntaxDirect(msg string) *SourceMapParseException {
	return NewSourceMapParseException(
		"JSON parse exception: " + jutil.GsonPackage + ".JsonSyntaxException: " + msg)
}

func (l *SourceMapJsonLexer) syntaxError(msg string, colDelta int) *SourceMapParseException {
	return wrapGsonSyntax(jutil.GsonPackage+".stream.MalformedJsonException", msg+l.locationString(l.pos, colDelta))
}

func (l *SourceMapJsonLexer) eofSyntaxError() *SourceMapParseException {
	return wrapGsonSyntax("java.io.EOFException", "End of input"+l.locationString(l.pos, 0))
}

func (l *SourceMapJsonLexer) checkLenient(colDelta int) jutil.Throwable {
	if !l.lenient {
		return l.syntaxError("Use JsonReader.setLenient(true) to accept malformed JSON", colDelta)
	}
	return nil
}

func (l *SourceMapJsonLexer) locationString(bytePos int, colDelta int) string {
	prefix := l.json[:bytePos]
	line := 1
	lineStartByte := l.bomBytes
	if lastNL := strings.LastIndexByte(prefix, '\n'); lastNL >= 0 {
		line = strings.Count(prefix[:lastNL+1], "\n") + 1
		lineStartByte = lastNL + 1
	}
	col := 1 + colDelta
	if lineStartByte < bytePos {
		col += jutil.StringLength(l.json[lineStartByte:bytePos])
	}
	return " at line " + strconv.Itoa(line) + " column " + strconv.Itoa(col) + " path " + l.getPath()
}

func (l *SourceMapJsonLexer) getPath() string {
	var sb strings.Builder
	sb.WriteByte('$')
	for _, f := range l.stack {
		switch f.scope {
		case smScopeEmptyArray, smScopeNonEmptyArray:
			sb.WriteByte('[')
			sb.WriteString(strconv.Itoa(f.index))
			sb.WriteByte(']')
		case smScopeEmptyObject, smScopeDanglingName, smScopeNonEmptyObject:
			sb.WriteByte('.')
			if f.hasName {
				sb.WriteString(f.name)
			}
		}
	}
	return sb.String()
}

func (l *SourceMapJsonLexer) skipToEndOfLine() {
	for l.pos < l.length {
		c := l.json[l.pos]
		l.pos++
		if c == '\n' || c == '\r' {
			break
		}
	}
}

// posOneUnitBeforeEOF matches JsonReader.skipTo("*/") on unterminated block comments,
// where JsonReader stops with 1 UTF-16 code unit remaining in its buffer (or 0 if "/*" is at EOF).
func (l *SourceMapJsonLexer) posOneUnitBeforeEOF(afterOpen int) (bytePos int, colDelta int) {
	if afterOpen >= l.length {
		return l.length, 0
	}
	// Walk back 1 UTF-8/WTF-8 code point from EOF.
	i := l.length - 1
	for i > afterOpen && (l.json[i]&0xC0) == 0x80 {
		i--
	}
	r, _ := util_decodeWTF8(l.json[i:])
	if r >= 0x10000 {
		// Supplementary rune = 2 UTF-16 units; 1 unit before EOF is between hi and lo surrogates.
		return i, 1
	}
	return i, 0
}

// nextNonWhitespace returns the next UTF-16 code unit, its byte width in l.json,
// any UTF-16 column adjustment (when width == 4 supplementary rune), and any error.
// On EOF: returns (-1, 0, 0, nil) if !throwOnEOF, else returns eofErr sentinel/exception.
func (l *SourceMapJsonLexer) nextNonWhitespace(throwOnEOF bool, rootInitialPeek bool) (int, int, int, jutil.Throwable, bool) {
	for l.pos < l.length {
		b := l.json[l.pos]
		if b == ' ' || b == '\r' || b == '\t' || b == '\n' {
			l.pos++
			continue
		}
		if b == '/' {
			l.pos++
			if l.pos == l.length {
				return '/', 1, 0, nil, false
			}
			if err := l.checkLenient(0); err != nil {
				return 0, 0, 0, err, false
			}
			switch l.json[l.pos] {
			case '*':
				l.pos++
				idx := strings.Index(l.json[l.pos:], "*/")
				if idx < 0 {
					errPos, delta := l.posOneUnitBeforeEOF(l.pos)
					l.pos = errPos
					return 0, 0, 0, l.syntaxError("Unterminated comment", delta), false
				}
				l.pos += idx + 2
				continue
			case '/':
				l.pos++
				l.skipToEndOfLine()
				continue
			default:
				return '/', 1, 0, nil, false
			}
		}
		if b == '#' {
			l.pos++
			if err := l.checkLenient(0); err != nil {
				return 0, 0, 0, err, false
			}
			l.skipToEndOfLine()
			continue
		}
		if b < 0x80 {
			l.pos++
			return int(b), 1, 0, nil, false
		}
		r, w := util_decodeWTF8(l.json[l.pos:])
		l.pos += w
		if r >= 0x10000 {
			hi, _ := utf16.EncodeRune(r)
			return int(hi), w, -1, nil, false
		}
		return int(r), w, 0, nil, false
	}
	if !throwOnEOF {
		return -1, 0, 0, nil, false
	}
	if rootInitialPeek {
		return 0, 0, 0, nil, true
	}
	return 0, 0, 0, l.eofSyntaxError(), false
}

func (l *SourceMapJsonLexer) doPeek(rootInitialPeek bool) (uint8, bool, jutil.Throwable) {
	peekStack := l.topScope()
	switch peekStack {
	case smScopeEmptyArray:
		l.setTopScope(smScopeNonEmptyArray)
	case smScopeNonEmptyArray:
		c, _, delta, err, _ := l.nextNonWhitespace(true, false)
		if err != nil {
			return 0, false, err
		}
		switch c {
		case ']':
			l.peeked = smPeekEndArray
			return l.peeked, false, nil
		case ';', ',':
		default:
			return 0, false, l.syntaxError("Unterminated array", delta)
		}
	case smScopeEmptyObject, smScopeNonEmptyObject:
		l.setTopScope(smScopeDanglingName)
		if peekStack == smScopeNonEmptyObject {
			c, _, delta, err, _ := l.nextNonWhitespace(true, false)
			if err != nil {
				return 0, false, err
			}
			switch c {
			case '}':
				l.peeked = smPeekEndObject
				return l.peeked, false, nil
			case ';', ',':
			default:
				return 0, false, l.syntaxError("Unterminated object", delta)
			}
		}
		c, w, delta, err, _ := l.nextNonWhitespace(true, false)
		if err != nil {
			return 0, false, err
		}
		switch c {
		case '"':
			l.peeked = smPeekDoubleQuotedName
			return l.peeked, false, nil
		case '\'':
			l.peeked = smPeekSingleQuotedName
			return l.peeked, false, nil
		case '}':
			if peekStack != smScopeNonEmptyObject {
				l.peeked = smPeekEndObject
				return l.peeked, false, nil
			}
			return 0, false, l.syntaxError("Expected name", delta)
		default:
			l.pos -= w
			if c >= 0x80 || !smNonLiteralLUT[byte(c)] {
				l.peeked = smPeekUnquotedName
				return l.peeked, false, nil
			}
			return 0, false, l.syntaxError("Expected name", 0)
		}
	case smScopeDanglingName:
		l.setTopScope(smScopeNonEmptyObject)
		c, _, delta, err, _ := l.nextNonWhitespace(true, false)
		if err != nil {
			return 0, false, err
		}
		switch c {
		case ':':
		case '=':
			if l.pos < l.length && l.json[l.pos] == '>' {
				l.pos++
			}
		default:
			return 0, false, l.syntaxError("Expected ':'", delta)
		}
	case smScopeEmptyDocument:
		if l.lenient {
			_, w, _, err, emptyEOF := l.nextNonWhitespace(true, rootInitialPeek)
			if err != nil || emptyEOF {
				return 0, emptyEOF, err
			}
			l.pos -= w
			if l.pos+5 <= l.length && l.json[l.pos:l.pos+5] == ")]}'\n" {
				l.pos += 5
			}
		}
		l.setTopScope(smScopeNonEmptyDocument)
	case smScopeNonEmptyDocument:
		c, w, delta, err, _ := l.nextNonWhitespace(false, false)
		if err != nil {
			return 0, false, err
		}
		if c == -1 {
			l.peeked = smPeekEOF
			return l.peeked, false, nil
		}
		if err := l.checkLenient(delta); err != nil {
			return 0, false, err
		}
		l.pos -= w
	}

	c, w, delta, err, emptyEOF := l.nextNonWhitespace(true, rootInitialPeek)
	if err != nil || emptyEOF {
		return 0, emptyEOF, err
	}
	switch c {
	case ']':
		if peekStack == smScopeEmptyArray {
			l.peeked = smPeekEndArray
			return l.peeked, false, nil
		}
		fallthrough
	case ';', ',':
		if peekStack == smScopeEmptyArray || peekStack == smScopeNonEmptyArray {
			l.pos -= w
			l.peeked = smPeekNull
			return l.peeked, false, nil
		}
		return 0, false, l.syntaxError("Unexpected value", delta)
	case '\'':
		l.peeked = smPeekSingleQuoted
		return l.peeked, false, nil
	case '"':
		l.peeked = smPeekDoubleQuoted
		return l.peeked, false, nil
	case '[':
		l.peeked = smPeekBeginArray
		return l.peeked, false, nil
	case '{':
		l.peeked = smPeekBeginObject
		return l.peeked, false, nil
	default:
		l.pos -= w
	}

	if kw := l.peekKeyword(); kw != smPeekNone {
		return kw, false, nil
	}
	if num := l.peekNumber(); num != smPeekNone {
		return num, false, nil
	}
	if l.json[l.pos] < 0x80 && smNonLiteralLUT[l.json[l.pos]] {
		return 0, false, l.syntaxError("Expected value", 0)
	}
	l.peeked = smPeekUnquoted
	return l.peeked, false, nil
}

func (l *SourceMapJsonLexer) peekKeyword() uint8 {
	b := l.json[l.pos]
	var kwLower, kwUpper string
	var peeking uint8
	switch b {
	case 't', 'T':
		kwLower, kwUpper, peeking = "true", "TRUE", smPeekTrue
	case 'f', 'F':
		kwLower, kwUpper, peeking = "false", "FALSE", smPeekFalse
	case 'n', 'N':
		kwLower, kwUpper, peeking = "null", "NULL", smPeekNull
	default:
		return smPeekNone
	}
	n := len(kwLower)
	if l.pos+n > l.length {
		return smPeekNone
	}
	for i := 1; i < n; i++ {
		c := l.json[l.pos+i]
		if c != kwLower[i] && c != kwUpper[i] {
			return smPeekNone
		}
	}
	if l.pos+n < l.length {
		after := l.json[l.pos+n]
		if after >= 0x80 || !smNonLiteralLUT[after] {
			return smPeekNone
		}
	}
	l.pos += n
	l.peeked = peeking
	return peeking
}

func (l *SourceMapJsonLexer) peekNumber() uint8 {
	s := l.json
	p := l.pos
	end := l.length

	var value int64
	negative := false
	fitsInLong := true
	last := smNumNone
	i := 0

loop:
	for ; ; i++ {
		if i == 1024 {
			return smPeekNone
		}
		if p+i == end {
			break
		}
		c := s[p+i]
		switch c {
		case '-':
			if last == smNumNone {
				negative = true
				last = smNumSign
				continue
			} else if last == smNumExpE {
				last = smNumExpSign
				continue
			}
			return smPeekNone
		case '+':
			if last == smNumExpE {
				last = smNumExpSign
				continue
			}
			return smPeekNone
		case 'e', 'E':
			if last == smNumDigit || last == smNumFractionDigit {
				last = smNumExpE
				continue
			}
			return smPeekNone
		case '.':
			if last == smNumDigit {
				last = smNumDecimal
				continue
			}
			return smPeekNone
		default:
			if c < '0' || c > '9' {
				if c < 0x80 && smNonLiteralLUT[c] {
					break loop
				}
				return smPeekNone
			}
			if last == smNumSign || last == smNumNone {
				value = -int64(c - '0')
				last = smNumDigit
			} else if last == smNumDigit {
				if value == 0 {
					return smPeekNone
				}
				newValue := value*10 - int64(c-'0')
				fitsInLong = fitsInLong && (value > smMinIncompleteInt64 ||
					(value == smMinIncompleteInt64 && newValue < value))
				value = newValue
			} else if last == smNumDecimal {
				last = smNumFractionDigit
			} else if last == smNumExpE || last == smNumExpSign {
				last = smNumExpDigit
			}
		}
	}

	if last == smNumDigit && fitsInLong && (value != math.MinInt64 || negative) && (value != 0 || !negative) {
		if negative {
			l.peekedLong = value
		} else {
			l.peekedLong = -value
		}
		l.pos += i
		l.peeked = smPeekLong
		return smPeekLong
	}
	if last == smNumDigit || last == smNumFractionDigit || last == smNumExpDigit {
		l.peekedNumberLength = i
		l.peeked = smPeekNumber
		return smPeekNumber
	}
	return smPeekNone
}

func (l *SourceMapJsonLexer) peekStream() (uint8, jutil.Throwable) {
	if l.peeked != smPeekNone {
		return l.peeked, nil
	}
	p, _, err := l.doPeek(false)
	return p, err
}

func (l *SourceMapJsonLexer) hasNextStream() (bool, jutil.Throwable) {
	p, err := l.peekStream()
	if err != nil {
		return false, err
	}
	return p != smPeekEndObject && p != smPeekEndArray && p != smPeekEOF, nil
}

func (l *SourceMapJsonLexer) beginObjectStream() jutil.Throwable {
	p, err := l.peekStream()
	if err != nil {
		return err
	}
	if p != smPeekBeginObject {
		return jutil.NewIllegalStateExceptionMsg("Expected BEGIN_OBJECT")
	}
	l.pushScope(smScopeEmptyObject)
	l.peeked = smPeekNone
	return nil
}

func (l *SourceMapJsonLexer) endObjectStream() jutil.Throwable {
	p, err := l.peekStream()
	if err != nil {
		return err
	}
	if p != smPeekEndObject {
		return jutil.NewIllegalStateExceptionMsg("Expected END_OBJECT")
	}
	l.stack[len(l.stack)-1] = smFrame{}
	l.stack = l.stack[:len(l.stack)-1]
	l.incTopIndex()
	l.peeked = smPeekNone
	return nil
}

func (l *SourceMapJsonLexer) beginArrayStream() jutil.Throwable {
	p, err := l.peekStream()
	if err != nil {
		return err
	}
	if p != smPeekBeginArray {
		return jutil.NewIllegalStateExceptionMsg("Expected BEGIN_ARRAY")
	}
	l.pushScope(smScopeEmptyArray)
	l.peeked = smPeekNone
	return nil
}

func (l *SourceMapJsonLexer) endArrayStream() jutil.Throwable {
	p, err := l.peekStream()
	if err != nil {
		return err
	}
	if p != smPeekEndArray {
		return jutil.NewIllegalStateExceptionMsg("Expected END_ARRAY")
	}
	l.stack = l.stack[:len(l.stack)-1]
	l.incTopIndex()
	l.peeked = smPeekNone
	return nil
}

func (l *SourceMapJsonLexer) nextNameStream() (string, jutil.Throwable) {
	p, err := l.peekStream()
	if err != nil {
		return "", err
	}
	var name string
	switch p {
	case smPeekDoubleQuotedName:
		name, err = l.readQuotedValue('"')
	case smPeekSingleQuotedName:
		name, err = l.readQuotedValue('\'')
	case smPeekUnquotedName:
		name = l.readUnquotedValue()
	default:
		return "", jutil.NewIllegalStateExceptionMsg("Expected a name")
	}
	if err != nil {
		return "", err
	}
	l.peeked = smPeekNone
	l.setTopName(name)
	return name, nil
}

func (l *SourceMapJsonLexer) nextStringStream() (string, jutil.Throwable) {
	p, err := l.peekStream()
	if err != nil {
		return "", err
	}
	var res string
	switch p {
	case smPeekDoubleQuoted:
		res, err = l.readQuotedValue('"')
	case smPeekSingleQuoted:
		res, err = l.readQuotedValue('\'')
	case smPeekUnquoted:
		res = l.readUnquotedValue()
	case smPeekLong:
		res = strconv.FormatInt(l.peekedLong, 10)
	case smPeekNumber:
		res = l.json[l.pos : l.pos+l.peekedNumberLength]
		l.pos += l.peekedNumberLength
	default:
		return "", jutil.NewIllegalStateExceptionMsg("Expected a string")
	}
	if err != nil {
		return "", err
	}
	l.peeked = smPeekNone
	l.incTopIndex()
	return res, nil
}

func (l *SourceMapJsonLexer) consumePeekedScalar() {
	if l.peeked == smPeekNumber {
		l.pos += l.peekedNumberLength
	}
	l.peeked = smPeekNone
	l.incTopIndex()
}

func (l *SourceMapJsonLexer) readUnquotedValue() string {
	start := l.pos
	for l.pos < l.length {
		b := l.json[l.pos]
		if b < 0x80 && smNonLiteralLUT[b] {
			break
		}
		l.pos++
	}
	return l.json[start:l.pos]
}

func writeWTF8CodeUnit(sb *strings.Builder, pendingHi *uint16, u uint16) {
	if *pendingHi != 0 {
		if u >= 0xDC00 && u <= 0xDFFF {
			sb.WriteRune(utf16.DecodeRune(rune(*pendingHi), rune(u)))
			*pendingHi = 0
			return
		}
		hi := *pendingHi
		*pendingHi = 0
		sb.WriteByte(0xED)
		sb.WriteByte(0x80 | byte((hi>>6)&0x3F))
		sb.WriteByte(0x80 | byte(hi&0x3F))
	}
	if u >= 0xD800 && u <= 0xDBFF {
		*pendingHi = u
		return
	}
	if u >= 0xDC00 && u <= 0xDFFF {
		sb.WriteByte(0xED)
		sb.WriteByte(0x80 | byte((u>>6)&0x3F))
		sb.WriteByte(0x80 | byte(u&0x3F))
		return
	}
	sb.WriteRune(rune(u))
}

func flushPendingHi(sb *strings.Builder, pendingHi *uint16) {
	if *pendingHi != 0 {
		hi := *pendingHi
		*pendingHi = 0
		sb.WriteByte(0xED)
		sb.WriteByte(0x80 | byte((hi>>6)&0x3F))
		sb.WriteByte(0x80 | byte(hi&0x3F))
	}
}

func (l *SourceMapJsonLexer) readQuotedValue(quote byte) (string, jutil.Throwable) {
	start := l.pos
	for l.pos < l.length {
		c := l.json[l.pos]
		if c == quote {
			s := l.json[start:l.pos]
			l.pos++
			return s, nil
		}
		if c == '\\' {
			break
		}
		l.pos++
	}
	if l.pos >= l.length {
		return "", l.syntaxError("Unterminated string", 0)
	}

	var sb strings.Builder
	sb.Grow((l.pos - start) + 16)
	var pendingHi uint16

	appendRawSlice := func(seg string) {
		if pendingHi == 0 && (len(seg) < 3 || seg[len(seg)-3] != 0xED) {
			sb.WriteString(seg)
			return
		}
		for i := 0; i < len(seg); {
			if seg[i] < 0x80 {
				if pendingHi != 0 {
					flushPendingHi(&sb, &pendingHi)
				}
				j := i + 1
				for j < len(seg) && seg[j] < 0x80 {
					j++
				}
				sb.WriteString(seg[i:j])
				i = j
				continue
			}
			r, w := util_decodeWTF8(seg[i:])
			i += w
			if r >= 0x10000 {
				flushPendingHi(&sb, &pendingHi)
				sb.WriteRune(r)
			} else {
				writeWTF8CodeUnit(&sb, &pendingHi, uint16(r))
			}
		}
	}

	appendRawSlice(l.json[start:l.pos])

	for l.pos < l.length {
		c := l.json[l.pos]
		if c == quote {
			l.pos++
			flushPendingHi(&sb, &pendingHi)
			return sb.String(), nil
		}
		if c == '\\' {
			l.pos++
			u, err := l.readEscapeCharacter()
			if err != nil {
				return "", err
			}
			writeWTF8CodeUnit(&sb, &pendingHi, u)
			segStart := l.pos
			for l.pos < l.length {
				b := l.json[l.pos]
				if b == quote || b == '\\' {
					break
				}
				l.pos++
			}
			if segStart < l.pos {
				appendRawSlice(l.json[segStart:l.pos])
			}
			continue
		}
		l.pos++
	}
	return "", l.syntaxError("Unterminated string", 0)
}

func (l *SourceMapJsonLexer) readEscapeCharacter() (uint16, jutil.Throwable) {
	if l.pos >= l.length {
		return 0, l.syntaxError("Unterminated escape sequence", 0)
	}
	b := l.json[l.pos]
	if b >= 0x80 {
		r, w := util_decodeWTF8(l.json[l.pos:])
		if r >= 0x10000 {
			return 0, l.syntaxError("Invalid escape sequence", 1)
		}
		l.pos += w
		return 0, l.syntaxError("Invalid escape sequence", 0)
	}
	l.pos++
	switch b {
	case 'u':
		if l.pos+4 <= l.length && (l.json[l.pos]|l.json[l.pos+1]|l.json[l.pos+2]|l.json[l.pos+3]) < 0x80 {
			var res uint16
			for i := 0; i < 4; i++ {
				ch := l.json[l.pos+i]
				res <<= 4
				switch {
				case ch >= '0' && ch <= '9':
					res |= uint16(ch - '0')
				case ch >= 'a' && ch <= 'f':
					res |= uint16(ch - 'a' + 10)
				case ch >= 'A' && ch <= 'F':
					res |= uint16(ch - 'A' + 10)
				default:
					return 0, jutil.NewNumberFormatException("\\u" + l.json[l.pos:l.pos+4])
				}
			}
			l.pos += 4
			return res, nil
		}
		var units [4]uint16
		n := 0
		p := l.pos
		for n < 4 && p < l.length {
			if l.json[p] < 0x80 {
				units[n] = uint16(l.json[p])
				n++
				p++
				continue
			}
			r, w := util_decodeWTF8(l.json[p:])
			p += w
			if r >= 0x10000 {
				hi, lo := utf16.EncodeRune(r)
				units[n] = uint16(hi)
				n++
				if n < 4 {
					units[n] = uint16(lo)
					n++
				}
			} else {
				units[n] = uint16(r)
				n++
			}
		}
		if n < 4 {
			return 0, l.syntaxError("Unterminated escape sequence", 0)
		}
		for i := 0; i < 4; i++ {
			ch := units[i]
			if !((ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')) {
				return 0, jutil.NewNumberFormatException("\\u" + jutil.StringFromUTF16(units[:]))
			}
		}
		l.pos = p
		// All 4 units are ASCII hex digits if we got here.
		var res uint16
		for i := 0; i < 4; i++ {
			ch := units[i]
			res <<= 4
			switch {
			case ch >= '0' && ch <= '9':
				res |= ch - '0'
			case ch >= 'a' && ch <= 'f':
				res |= ch - 'a' + 10
			default:
				res |= ch - 'A' + 10
			}
		}
		return res, nil
	case 't':
		return '\t', nil
	case 'b':
		return '\b', nil
	case 'n':
		return '\n', nil
	case 'r':
		return '\r', nil
	case 'f':
		return '\f', nil
	case '\n', '\'', '"', '\\', '/':
		return uint16(b), nil
	default:
		return 0, l.syntaxError("Invalid escape sequence", 0)
	}
}

// sourceMapJsonLexer_parseInt is Integer.parseInt(CharSequence, int, int,
// int) on s[beginIndex:endIndex] (for radix 10 or 16); ok is false where Java
// throws NumberFormatException.
func sourceMapJsonLexer_parseInt(s string, beginIndex, endIndex, radix int) (int32, bool) {
	i := beginIndex
	neg := false
	if i < endIndex && (s[i] == '-' || s[i] == '+') {
		neg = s[i] == '-'
		i++
	}
	if i == endIndex {
		return 0, false
	}
	limit := uint64(math.MaxInt32)
	if neg {
		limit++
	}
	base := uint64(radix)
	var mag uint64
	for i < endIndex {
		c := s[i]
		var d int
		if c >= '0' && c <= '9' {
			d = int(c - '0')
			i++
		} else if radix == 16 && c >= 'a' && c <= 'f' {
			d = int(c - 'a' + 10)
			i++
		} else if radix == 16 && c >= 'A' && c <= 'F' {
			d = int(c - 'A' + 10)
			i++
		} else if c < 0x80 {
			return 0, false
		} else {
			r, width := util_decodeWTF8(s[i:endIndex])
			if r >= 0x10000 {
				return 0, false
			}
			d = jutil.CharacterDigit(uint16(r), radix)
			if d < 0 {
				return 0, false
			}
			i += width
		}
		if mag > (limit-uint64(d))/base {
			return 0, false
		}
		mag = mag*base + uint64(d)
	}
	if neg {
		return int32(-int64(mag)), true
	}
	return int32(mag), true
}
