// Port of com/google/gson/stream/JsonReader.java (Gson 2.9.1), with
// com/google/gson/stream/JsonToken.java and JsonScope.java (all three files
// carry the same header)

/*
 * Copyright (C) 2010 Google Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jutil

import (
	"math"
	"slices"
	"strconv"
	"strings"
	"unicode/utf16"
)

// Peeked token kinds produced by JsonScanner.
const (
	JsonPeekNone uint8 = iota
	JsonPeekBeginObject
	JsonPeekEndObject
	JsonPeekBeginArray
	JsonPeekEndArray
	JsonPeekTrue
	JsonPeekFalse
	JsonPeekNull
	JsonPeekSingleQuoted
	JsonPeekDoubleQuoted
	JsonPeekUnquoted
	JsonPeekSingleQuotedName
	JsonPeekDoubleQuotedName
	JsonPeekUnquotedName
	JsonPeekLong
	JsonPeekNumber
	JsonPeekEOF
)

// com.google.gson.stream.JsonScope
const (
	jsonScope_EMPTY_ARRAY       = 1
	jsonScope_NONEMPTY_ARRAY    = 2
	jsonScope_EMPTY_OBJECT      = 3
	jsonScope_DANGLING_NAME     = 4
	jsonScope_NONEMPTY_OBJECT   = 5
	jsonScope_EMPTY_DOCUMENT    = 6
	jsonScope_NONEMPTY_DOCUMENT = 7
	jsonScope_CLOSED            = 8
)

const (
	jsonNumNone = iota
	jsonNumSign
	jsonNumDigit
	jsonNumDecimal
	jsonNumFractionDigit
	jsonNumExpE
	jsonNumExpSign
	jsonNumExpDigit
)

const (
	jsonMinIncompleteInt64 = -0x8000000000000000 / 10
	jsonMaxNumberChars     = 1024
)

// jsonNonLiteralLUT marks the 16 ASCII bytes that terminate unquoted literals in Gson.
var jsonNonLiteralLUT = [256]bool{
	'/': true, '\\': true, ';': true, '#': true, '=': true,
	'{': true, '}': true, '[': true, ']': true, ':': true, ',': true,
	' ': true, '\t': true, '\f': true, '\r': true, '\n': true,
}

type jsonFrame struct {
	name    string
	index   int
	scope   uint8
	hasName bool
}

// JsonScanner is a zero-allocation streaming Gson 2.9.1 JSON token scanner
// operating directly on a UTF-8/WTF-8 string. Unescaped strings are sliced
// directly from json without intermediate buffers, while preserving Gson's
// exact strict/lenient rules, 1024-char number limit, UTF-16 line/column
// locations, and JsonPath tracking.
type JsonScanner struct {
	json      string
	pos       int
	length    int
	bomBytes  int
	nonExecNL int
	lenient   bool

	peeked             uint8
	peekedLong         int64
	peekedNumberLength int

	stack    []jsonFrame
	stackBuf [16]jsonFrame
}

// Init resets s to scan json with the given leniency mode.
func (s *JsonScanner) Init(json string, lenient bool) {
	*s = JsonScanner{
		json:    json,
		length:  len(json),
		lenient: lenient,
	}
	s.stackBuf[0] = jsonFrame{scope: jsonScope_EMPTY_DOCUMENT}
	s.stack = s.stackBuf[:1]
	if strings.HasPrefix(json, "\xef\xbb\xbf") {
		s.pos = 3
		s.bomBytes = 3
	}
}

// Clone returns an independent value snapshot of s.
func (s *JsonScanner) Clone() JsonScanner {
	cp := *s
	if len(s.stack) > 0 && &s.stack[0] == &s.stackBuf[0] {
		cp.stack = cp.stackBuf[:len(s.stack)]
	} else {
		cp.stack = slices.Clone(s.stack)
	}
	return cp
}

func (s *JsonScanner) SetLenient(lenient bool) { s.lenient = lenient }
func (s *JsonScanner) IsLenient() bool         { return s.lenient }
func (s *JsonScanner) PeekedLong() int64       { return s.peekedLong }

// Close marks the scanner closed so subsequent operations fail with IllegalStateException.
func (s *JsonScanner) Close() {
	s.peeked = JsonPeekNone
	s.stackBuf[0] = jsonFrame{scope: jsonScope_CLOSED}
	s.stack = s.stackBuf[:1]
}

func (s *JsonScanner) topScope() uint8     { return s.stack[len(s.stack)-1].scope }
func (s *JsonScanner) setTopScope(v uint8) { s.stack[len(s.stack)-1].scope = v }
func (s *JsonScanner) pushScope(v uint8)   { s.stack = append(s.stack, jsonFrame{scope: v}) }
func (s *JsonScanner) incTopIndex()        { s.stack[len(s.stack)-1].index++ }
func (s *JsonScanner) setTopName(name string) {
	f := &s.stack[len(s.stack)-1]
	f.name = name
	f.hasName = true
}

func (s *JsonScanner) syntaxError(msg string, colDelta int) *malformedJsonException {
	return newMalformedJsonException(msg + s.locationString(s.pos, colDelta))
}

func (s *JsonScanner) eofSyntaxError() *eofException {
	return newEOFException("End of input" + s.locationString(s.pos, 0))
}

func (s *JsonScanner) checkLenient(colDelta int) Throwable {
	if !s.lenient {
		return s.syntaxError("Use JsonReader.setLenient(true) to accept malformed JSON", colDelta)
	}
	return nil
}

func (s *JsonScanner) isLiteralByte(b byte) (bool, Throwable) {
	if b >= 0x80 || !jsonNonLiteralLUT[b] {
		return true, nil
	}
	if !s.lenient {
		switch b {
		case '/', '\\', ';', '#', '=':
			return false, s.checkLenient(0)
		}
	}
	return false, nil
}

// LocationString formats " at line L column C path P" at the current scanner position.
func (s *JsonScanner) LocationString() string {
	return s.locationString(s.pos, 0)
}

func (s *JsonScanner) locationString(bytePos int, colDelta int) string {
	prefix := s.json[:bytePos]
	line := 1
	lineStartByte := s.bomBytes
	if lastNL := strings.LastIndexByte(prefix, '\n'); lastNL >= 0 {
		nlCount := strings.Count(prefix[:lastNL+1], "\n")
		if s.nonExecNL > 0 && s.nonExecNL < bytePos {
			nlCount--
			if lastNL == s.nonExecNL {
				lastNL = strings.LastIndexByte(prefix[:lastNL], '\n')
			}
		}
		line = nlCount + 1
		if lastNL >= 0 {
			lineStartByte = lastNL + 1
		}
	}
	col := 1 + colDelta
	if lineStartByte < bytePos {
		col += StringLength(s.json[lineStartByte:bytePos])
	}
	return " at line " + strconv.Itoa(line) + " column " + strconv.Itoa(col) + " path " + s.getPath(false)
}

func (s *JsonScanner) getPath(usePreviousPath bool) string {
	var sb strings.Builder
	sb.WriteByte('$')
	for i, f := range s.stack {
		switch f.scope {
		case jsonScope_EMPTY_ARRAY, jsonScope_NONEMPTY_ARRAY:
			idx := f.index
			if usePreviousPath && idx > 0 && i == len(s.stack)-1 {
				idx--
			}
			sb.WriteByte('[')
			sb.WriteString(strconv.Itoa(idx))
			sb.WriteByte(']')
		case jsonScope_EMPTY_OBJECT, jsonScope_DANGLING_NAME, jsonScope_NONEMPTY_OBJECT:
			sb.WriteByte('.')
			if f.hasName {
				sb.WriteString(f.name)
			}
		}
	}
	return sb.String()
}

// GetPath returns a JsonPath to the next (to be read) element.
func (s *JsonScanner) GetPath() string { return s.getPath(false) }

// GetPreviousPath returns a JsonPath to the element just read.
func (s *JsonScanner) GetPreviousPath() string { return s.getPath(true) }

func (s *JsonScanner) skipToEndOfLine() {
	for s.pos < s.length {
		c := s.json[s.pos]
		s.pos++
		if c == '\n' || c == '\r' {
			break
		}
	}
}

// posOneUnitBeforeEOF matches JsonReader.skipTo("*/") on unterminated block comments,
// where JsonReader stops with 1 UTF-16 code unit remaining in its buffer (or 0 if "/*" is at EOF).
func (s *JsonScanner) posOneUnitBeforeEOF(afterOpen int) (bytePos int, colDelta int) {
	if afterOpen >= s.length {
		return s.length, 0
	}
	i := s.length - 1
	for i > afterOpen && (s.json[i]&0xC0) == 0x80 {
		i--
	}
	r, _ := decodeWTF8(s.json[i:])
	if r >= 0x10000 {
		return i, 1
	}
	return i, 0
}

func (s *JsonScanner) nextNonWhitespace(throwOnEOF bool, rootInitialPeek bool) (int, int, int, Throwable, bool) {
	for s.pos < s.length {
		b := s.json[s.pos]
		if b == ' ' || b == '\r' || b == '\t' || b == '\n' {
			s.pos++
			continue
		}
		if b == '/' {
			s.pos++
			if s.pos == s.length {
				return '/', 1, 0, nil, false
			}
			if err := s.checkLenient(0); err != nil {
				return 0, 0, 0, err, false
			}
			switch s.json[s.pos] {
			case '*':
				s.pos++
				idx := strings.Index(s.json[s.pos:], "*/")
				if idx < 0 {
					errPos, delta := s.posOneUnitBeforeEOF(s.pos)
					s.pos = errPos
					return 0, 0, 0, s.syntaxError("Unterminated comment", delta), false
				}
				s.pos += idx + 2
				continue
			case '/':
				s.pos++
				s.skipToEndOfLine()
				continue
			default:
				return '/', 1, 0, nil, false
			}
		}
		if b == '#' {
			s.pos++
			if err := s.checkLenient(0); err != nil {
				return 0, 0, 0, err, false
			}
			s.skipToEndOfLine()
			continue
		}
		if b < 0x80 {
			s.pos++
			return int(b), 1, 0, nil, false
		}
		r, w := decodeWTF8(s.json[s.pos:])
		s.pos += w
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
	return 0, 0, 0, s.eofSyntaxError(), false
}

// DoPeek advances to and classifies the next token. When rootInitialPeek is
// true and the input contains only whitespace/comments, emptyEOF is returned
// as true without allocating an EOFException.
func (s *JsonScanner) DoPeek(rootInitialPeek bool) (uint8, bool, Throwable) {
	peekStack := s.topScope()
	switch peekStack {
	case jsonScope_EMPTY_ARRAY:
		s.setTopScope(jsonScope_NONEMPTY_ARRAY)
	case jsonScope_NONEMPTY_ARRAY:
		c, _, delta, err, _ := s.nextNonWhitespace(true, false)
		if err != nil {
			return 0, false, err
		}
		switch c {
		case ']':
			s.peeked = JsonPeekEndArray
			return s.peeked, false, nil
		case ';':
			if err := s.checkLenient(delta); err != nil {
				return 0, false, err
			}
		case ',':
		default:
			return 0, false, s.syntaxError("Unterminated array", delta)
		}
	case jsonScope_EMPTY_OBJECT, jsonScope_NONEMPTY_OBJECT:
		s.setTopScope(jsonScope_DANGLING_NAME)
		if peekStack == jsonScope_NONEMPTY_OBJECT {
			c, _, delta, err, _ := s.nextNonWhitespace(true, false)
			if err != nil {
				return 0, false, err
			}
			switch c {
			case '}':
				s.peeked = JsonPeekEndObject
				return s.peeked, false, nil
			case ';':
				if err := s.checkLenient(delta); err != nil {
					return 0, false, err
				}
			case ',':
			default:
				return 0, false, s.syntaxError("Unterminated object", delta)
			}
		}
		c, w, delta, err, _ := s.nextNonWhitespace(true, false)
		if err != nil {
			return 0, false, err
		}
		switch c {
		case '"':
			s.peeked = JsonPeekDoubleQuotedName
			return s.peeked, false, nil
		case '\'':
			if err := s.checkLenient(delta); err != nil {
				return 0, false, err
			}
			s.peeked = JsonPeekSingleQuotedName
			return s.peeked, false, nil
		case '}':
			if peekStack != jsonScope_NONEMPTY_OBJECT {
				s.peeked = JsonPeekEndObject
				return s.peeked, false, nil
			}
			return 0, false, s.syntaxError("Expected name", delta)
		default:
			if err := s.checkLenient(delta); err != nil {
				return 0, false, err
			}
			s.pos -= w
			lit, err := s.isLiteralByte(s.json[s.pos])
			if err != nil {
				return 0, false, err
			}
			if lit {
				s.peeked = JsonPeekUnquotedName
				return s.peeked, false, nil
			}
			return 0, false, s.syntaxError("Expected name", 0)
		}
	case jsonScope_DANGLING_NAME:
		s.setTopScope(jsonScope_NONEMPTY_OBJECT)
		c, _, delta, err, _ := s.nextNonWhitespace(true, false)
		if err != nil {
			return 0, false, err
		}
		switch c {
		case ':':
		case '=':
			if err := s.checkLenient(delta); err != nil {
				return 0, false, err
			}
			if s.pos < s.length && s.json[s.pos] == '>' {
				s.pos++
			}
		default:
			return 0, false, s.syntaxError("Expected ':'", delta)
		}
	case jsonScope_EMPTY_DOCUMENT:
		if s.lenient {
			_, w, _, err, emptyEOF := s.nextNonWhitespace(true, rootInitialPeek)
			if err != nil || emptyEOF {
				return 0, emptyEOF, err
			}
			s.pos -= w
			if s.pos+5 <= s.length && s.json[s.pos:s.pos+5] == ")]}'\n" {
				s.pos += 5
				s.nonExecNL = s.pos - 1
			}
		}
		s.setTopScope(jsonScope_NONEMPTY_DOCUMENT)
	case jsonScope_NONEMPTY_DOCUMENT:
		c, w, delta, err, _ := s.nextNonWhitespace(false, false)
		if err != nil {
			return 0, false, err
		}
		if c == -1 {
			s.peeked = JsonPeekEOF
			return s.peeked, false, nil
		}
		if err := s.checkLenient(delta); err != nil {
			return 0, false, err
		}
		s.pos -= w
	case jsonScope_CLOSED:
		return 0, false, NewIllegalStateExceptionMsg("JsonReader is closed")
	}

	c, w, delta, err, emptyEOF := s.nextNonWhitespace(true, rootInitialPeek)
	if err != nil || emptyEOF {
		return 0, emptyEOF, err
	}
	switch c {
	case ']':
		if peekStack == jsonScope_EMPTY_ARRAY {
			s.peeked = JsonPeekEndArray
			return s.peeked, false, nil
		}
		fallthrough
	case ';', ',':
		if peekStack == jsonScope_EMPTY_ARRAY || peekStack == jsonScope_NONEMPTY_ARRAY {
			if err := s.checkLenient(delta); err != nil {
				return 0, false, err
			}
			s.pos -= w
			s.peeked = JsonPeekNull
			return s.peeked, false, nil
		}
		return 0, false, s.syntaxError("Unexpected value", delta)
	case '\'':
		if err := s.checkLenient(delta); err != nil {
			return 0, false, err
		}
		s.peeked = JsonPeekSingleQuoted
		return s.peeked, false, nil
	case '"':
		s.peeked = JsonPeekDoubleQuoted
		return s.peeked, false, nil
	case '[':
		s.peeked = JsonPeekBeginArray
		return s.peeked, false, nil
	case '{':
		s.peeked = JsonPeekBeginObject
		return s.peeked, false, nil
	default:
		s.pos -= w
	}

	kw, err := s.peekKeyword()
	if err != nil {
		return 0, false, err
	}
	if kw != JsonPeekNone {
		return kw, false, nil
	}
	num, err := s.peekNumber()
	if err != nil {
		return 0, false, err
	}
	if num != JsonPeekNone {
		return num, false, nil
	}
	lit, err := s.isLiteralByte(s.json[s.pos])
	if err != nil {
		return 0, false, err
	}
	if !lit {
		return 0, false, s.syntaxError("Expected value", 0)
	}
	if err := s.checkLenient(0); err != nil {
		return 0, false, err
	}
	s.peeked = JsonPeekUnquoted
	return s.peeked, false, nil
}

func (s *JsonScanner) peekKeyword() (uint8, Throwable) {
	b := s.json[s.pos]
	var kwLower, kwUpper string
	var peeking uint8
	switch b {
	case 't', 'T':
		kwLower, kwUpper, peeking = "true", "TRUE", JsonPeekTrue
	case 'f', 'F':
		kwLower, kwUpper, peeking = "false", "FALSE", JsonPeekFalse
	case 'n', 'N':
		kwLower, kwUpper, peeking = "null", "NULL", JsonPeekNull
	default:
		return JsonPeekNone, nil
	}
	n := len(kwLower)
	if s.pos+n > s.length {
		return JsonPeekNone, nil
	}
	for i := 1; i < n; i++ {
		c := s.json[s.pos+i]
		if c != kwLower[i] && c != kwUpper[i] {
			return JsonPeekNone, nil
		}
	}
	if s.pos+n < s.length {
		lit, err := s.isLiteralByte(s.json[s.pos+n])
		if err != nil {
			return 0, err
		}
		if lit {
			return JsonPeekNone, nil
		}
	}
	s.pos += n
	s.peeked = peeking
	return peeking, nil
}

func (s *JsonScanner) peekNumber() (uint8, Throwable) {
	str := s.json
	p := s.pos
	end := s.length

	var value int64
	negative := false
	fitsInLong := true
	last := jsonNumNone
	i := 0

loop:
	for ; ; i++ {
		if p+i == end {
			break
		}
		if i == jsonMaxNumberChars {
			return JsonPeekNone, nil
		}
		c := str[p+i]
		switch c {
		case '-':
			if last == jsonNumNone {
				negative = true
				last = jsonNumSign
				continue
			} else if last == jsonNumExpE {
				last = jsonNumExpSign
				continue
			}
			return JsonPeekNone, nil
		case '+':
			if last == jsonNumExpE {
				last = jsonNumExpSign
				continue
			}
			return JsonPeekNone, nil
		case 'e', 'E':
			if last == jsonNumDigit || last == jsonNumFractionDigit {
				last = jsonNumExpE
				continue
			}
			return JsonPeekNone, nil
		case '.':
			if last == jsonNumDigit {
				last = jsonNumDecimal
				continue
			}
			return JsonPeekNone, nil
		default:
			if c < '0' || c > '9' {
				lit, err := s.isLiteralByte(c)
				if err != nil {
					return 0, err
				}
				if !lit {
					break loop
				}
				return JsonPeekNone, nil
			}
			if last == jsonNumSign || last == jsonNumNone {
				value = -int64(c - '0')
				last = jsonNumDigit
			} else if last == jsonNumDigit {
				if value == 0 {
					return JsonPeekNone, nil
				}
				newValue := value*10 - int64(c-'0')
				fitsInLong = fitsInLong && (value > jsonMinIncompleteInt64 ||
					(value == jsonMinIncompleteInt64 && newValue < value))
				value = newValue
			} else if last == jsonNumDecimal {
				last = jsonNumFractionDigit
			} else if last == jsonNumExpE || last == jsonNumExpSign {
				last = jsonNumExpDigit
			}
		}
	}

	if last == jsonNumDigit && fitsInLong && (value != math.MinInt64 || negative) && (value != 0 || !negative) {
		if negative {
			s.peekedLong = value
		} else {
			s.peekedLong = -value
		}
		s.pos += i
		s.peeked = JsonPeekLong
		return JsonPeekLong, nil
	}
	if last == jsonNumDigit || last == jsonNumFractionDigit || last == jsonNumExpDigit {
		s.peekedNumberLength = i
		s.peeked = JsonPeekNumber
		return JsonPeekNumber, nil
	}
	return JsonPeekNone, nil
}

// Peek returns the internal peeked token code without consuming it.
func (s *JsonScanner) Peek() (uint8, Throwable) {
	if s.peeked != JsonPeekNone {
		return s.peeked, nil
	}
	p, _, err := s.DoPeek(false)
	return p, err
}

func (s *JsonScanner) readUnquotedValue() (string, Throwable) {
	start := s.pos
	for i := start; i < s.length; i++ {
		b := s.json[i]
		if b < 0x80 && jsonNonLiteralLUT[b] {
			if _, err := s.isLiteralByte(b); err != nil {
				return "", err
			}
			s.pos = i
			return s.json[start:i], nil
		}
	}
	s.pos = s.length
	return s.json[start:s.length], nil
}

func (s *JsonScanner) skipQuotedValue(quote byte) Throwable {
	for s.pos < s.length {
		c := s.json[s.pos]
		s.pos++
		if c == quote {
			return nil
		}
		if c == '\\' {
			if _, err := s.readEscapeCharacter(); err != nil {
				return err
			}
		}
	}
	return s.syntaxError("Unterminated string", 0)
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
		appendWTF8(sb, rune(hi))
	}
	if u >= 0xD800 && u <= 0xDBFF {
		*pendingHi = u
		return
	}
	appendWTF8(sb, rune(u))
}

func flushPendingHi(sb *strings.Builder, pendingHi *uint16) {
	if *pendingHi != 0 {
		hi := *pendingHi
		*pendingHi = 0
		appendWTF8(sb, rune(hi))
	}
}

func (s *JsonScanner) readQuotedValue(quote byte) (string, Throwable) {
	start := s.pos
	for s.pos < s.length {
		c := s.json[s.pos]
		if c == quote {
			res := s.json[start:s.pos]
			s.pos++
			return res, nil
		}
		if c == '\\' {
			break
		}
		s.pos++
	}
	if s.pos >= s.length {
		return "", s.syntaxError("Unterminated string", 0)
	}

	var sb strings.Builder
	sb.Grow((s.pos - start) + 16)
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
			r, w := decodeWTF8(seg[i:])
			i += w
			if r >= 0x10000 {
				flushPendingHi(&sb, &pendingHi)
				sb.WriteRune(r)
			} else {
				writeWTF8CodeUnit(&sb, &pendingHi, uint16(r))
			}
		}
	}

	appendRawSlice(s.json[start:s.pos])

	for s.pos < s.length {
		c := s.json[s.pos]
		if c == quote {
			s.pos++
			flushPendingHi(&sb, &pendingHi)
			return sb.String(), nil
		}
		if c == '\\' {
			s.pos++
			u, err := s.readEscapeCharacter()
			if err != nil {
				return "", err
			}
			writeWTF8CodeUnit(&sb, &pendingHi, u)
			segStart := s.pos
			for s.pos < s.length {
				b := s.json[s.pos]
				if b == quote || b == '\\' {
					break
				}
				s.pos++
			}
			if segStart < s.pos {
				appendRawSlice(s.json[segStart:s.pos])
			}
			continue
		}
		s.pos++
	}
	return "", s.syntaxError("Unterminated string", 0)
}

func parseHexDigit(ch uint16) (uint16, bool) {
	switch {
	case ch >= '0' && ch <= '9':
		return ch - '0', true
	case ch >= 'a' && ch <= 'f':
		return ch - 'a' + 10, true
	case ch >= 'A' && ch <= 'F':
		return ch - 'A' + 10, true
	default:
		return 0, false
	}
}

func (s *JsonScanner) readEscapeCharacter() (uint16, Throwable) {
	if s.pos >= s.length {
		return 0, s.syntaxError("Unterminated escape sequence", 0)
	}
	b := s.json[s.pos]
	if b >= 0x80 {
		r, w := decodeWTF8(s.json[s.pos:])
		if r >= 0x10000 {
			return 0, s.syntaxError("Invalid escape sequence", 1)
		}
		s.pos += w
		return 0, s.syntaxError("Invalid escape sequence", 0)
	}
	s.pos++
	switch b {
	case 'u':
		if s.pos+4 <= s.length && (s.json[s.pos]|s.json[s.pos+1]|s.json[s.pos+2]|s.json[s.pos+3]) < 0x80 {
			var res uint16
			for i := 0; i < 4; i++ {
				d, ok := parseHexDigit(uint16(s.json[s.pos+i]))
				if !ok {
					return 0, NewNumberFormatException("\\u" + s.json[s.pos:s.pos+4])
				}
				res = (res << 4) | d
			}
			s.pos += 4
			return res, nil
		}
		var units [4]uint16
		n, p := 0, s.pos
		for n < 4 && p < s.length {
			if s.json[p] < 0x80 {
				units[n] = uint16(s.json[p])
				n++
				p++
				continue
			}
			r, w := decodeWTF8(s.json[p:])
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
			return 0, s.syntaxError("Unterminated escape sequence", 0)
		}
		var res uint16
		for i := 0; i < 4; i++ {
			d, ok := parseHexDigit(units[i])
			if !ok {
				return 0, NewNumberFormatException("\\u" + StringFromUTF16(units[:]))
			}
			res = (res << 4) | d
		}
		s.pos = p
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
		return 0, s.syntaxError("Invalid escape sequence", 0)
	}
}
