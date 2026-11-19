// Clean-room Gson 2.9.1 streaming JsonReader and JsonToken (Apache 2.0).

package gson

import (
	"math"
	"strconv"
	"strings"
	"unicode/utf16"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// JsonToken implements com.google.gson.stream.JsonToken.
type JsonToken int32

const (
	JsonToken_BEGIN_ARRAY JsonToken = iota
	JsonToken_END_ARRAY
	JsonToken_BEGIN_OBJECT
	JsonToken_END_OBJECT
	JsonToken_NAME
	JsonToken_STRING
	JsonToken_NUMBER
	JsonToken_BOOLEAN
	JsonToken_NULL
	JsonToken_END_DOCUMENT
)

var jsonTokenNames = [...]string{
	"BEGIN_ARRAY", "END_ARRAY", "BEGIN_OBJECT", "END_OBJECT", "NAME",
	"STRING", "NUMBER", "BOOLEAN", "NULL", "END_DOCUMENT",
}

func (t JsonToken) String() string { return jsonTokenNames[t] }

const (
	peekNone uint8 = iota
	peekBeginObject
	peekEndObject
	peekBeginArray
	peekEndArray
	peekTrue
	peekFalse
	peekNull
	peekSingleQuoted
	peekDoubleQuoted
	peekUnquoted
	peekSingleQuotedName
	peekDoubleQuotedName
	peekUnquotedName
	peekLong
	peekNumber
	peekEOF
)

const (
	numNone = iota
	numSign
	numDigit
	numDecimal
	numFractionDigit
	numExpE
	numExpSign
	numExpDigit
)

const (
	minIncompleteInt64 = -0x8000000000000000 / 10
	maxNumberChars     = 1024
)

var nonLiteralLUT = [256]bool{
	'/': true, '\\': true, ';': true, '#': true, '=': true,
	'{': true, '}': true, '[': true, ']': true, ':': true, ',': true,
	' ': true, '\t': true, '\f': true, '\r': true, '\n': true,
}

type readerFrame struct {
	name    string
	index   int
	scope   uint8
	hasName bool
}

// MalformedJsonException implements com.google.gson.stream.MalformedJsonException.
type MalformedJsonException struct {
	javaio.IOException
}

func NewMalformedJsonException(msg string) *MalformedJsonException {
	e := &MalformedJsonException{}
	javalang.InitThrowable(e, "com.google.gson.stream.MalformedJsonException", msg, nil)
	return e
}

// JsonReader implements com.google.gson.stream.JsonReader.
type JsonReader struct {
	closer    javaio.Reader
	json      string
	pos       int
	length    int
	bomBytes  int
	nonExecNL int
	lenient   bool

	peeked             uint8
	peekedLong         int64
	peekedNumberLength int

	stack    []readerFrame
	stackBuf [16]readerFrame
}

// NewJsonReader reads all characters from in and initializes a streaming JsonReader.
func NewJsonReader(in javaio.Reader) *JsonReader {
	if in == nil {
		panic(javalang.NewNullPointerExceptionMsg("in == null"))
	}
	var all []uint16
	var buf [4096]uint16
	for {
		n := in.Read(buf[:], 0, len(buf))
		if n <= 0 {
			break
		}
		all = append(all, buf[:n]...)
	}
	r := NewJsonReaderString(javalang.StringFromUTF16(all))
	r.closer = in
	return r
}

// NewJsonReaderString creates a JsonReader scanning json directly.
func NewJsonReaderString(json string) *JsonReader {
	r := &JsonReader{
		json:   json,
		length: len(json),
	}
	r.stackBuf[0] = readerFrame{scope: jsonScopeEmptyDocument}
	r.stack = r.stackBuf[:1]
	if strings.HasPrefix(json, "\xef\xbb\xbf") {
		r.pos = 3
		r.bomBytes = 3
	}
	return r
}

func (r *JsonReader) SetLenient(lenient bool) { r.lenient = lenient }
func (r *JsonReader) IsLenient() bool         { return r.lenient }

func (r *JsonReader) Close() {
	r.peeked = peekNone
	r.stackBuf[0] = readerFrame{scope: jsonScopeClosed}
	r.stack = r.stackBuf[:1]
	if r.closer != nil {
		_ = r.closer.Close()
	}
}

func (r *JsonReader) topScope() uint8     { return r.stack[len(r.stack)-1].scope }
func (r *JsonReader) setTopScope(v uint8) { r.stack[len(r.stack)-1].scope = v }
func (r *JsonReader) pushScope(v uint8)   { r.stack = append(r.stack, readerFrame{scope: v}) }
func (r *JsonReader) incTopIndex()        { r.stack[len(r.stack)-1].index++ }
func (r *JsonReader) setTopName(name string) {
	f := &r.stack[len(r.stack)-1]
	f.name = name
	f.hasName = true
}

func (r *JsonReader) GetPath() string         { return r.getPath(false) }
func (r *JsonReader) GetPreviousPath() string { return r.getPath(true) }

func (r *JsonReader) getPath(usePrevious bool) string {
	var sb strings.Builder
	sb.WriteByte('$')
	for i, f := range r.stack {
		switch f.scope {
		case jsonScopeEmptyArray, jsonScopeNonemptyArray:
			idx := f.index
			if usePrevious && idx > 0 && i == len(r.stack)-1 {
				idx--
			}
			sb.WriteByte('[')
			sb.WriteString(strconv.Itoa(idx))
			sb.WriteByte(']')
		case jsonScopeEmptyObject, jsonScopeDanglingName, jsonScopeNonemptyObject:
			sb.WriteByte('.')
			if f.hasName {
				sb.WriteString(f.name)
			}
		}
	}
	return sb.String()
}

func (r *JsonReader) locationString(bytePos, colDelta int) string {
	prefix := r.json[:bytePos]
	line := 1
	lineStartByte := r.bomBytes
	if lastNL := strings.LastIndexByte(prefix, '\n'); lastNL >= 0 {
		nlCount := strings.Count(prefix[:lastNL+1], "\n")
		if r.nonExecNL > 0 && r.nonExecNL < bytePos {
			nlCount--
			if lastNL == r.nonExecNL {
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
		col += javalang.StringLength(r.json[lineStartByte:bytePos])
	}
	return " at line " + strconv.Itoa(line) + " column " + strconv.Itoa(col) + " path " + r.getPath(false)
}

func (r *JsonReader) syntaxError(msg string, colDelta int) javalang.Throwable {
	return NewMalformedJsonException(msg + r.locationString(r.pos, colDelta))
}

func (r *JsonReader) eofError() javalang.Throwable {
	return javaio.NewEOFException("End of input" + r.locationString(r.pos, 0))
}

func (r *JsonReader) checkLenient(colDelta int) javalang.Throwable {
	if !r.lenient {
		return r.syntaxError("Use JsonReader.setLenient(true) to accept malformed JSON", colDelta)
	}
	return nil
}

func (r *JsonReader) isLiteralByte(b byte) (bool, javalang.Throwable) {
	if b >= 0x80 || !nonLiteralLUT[b] {
		return true, nil
	}
	if !r.lenient {
		switch b {
		case '/', '\\', ';', '#', '=':
			return false, r.checkLenient(0)
		}
	}
	return false, nil
}

func decodeWTF8At(s string) (rune, int) {
	b0 := s[0]
	if b0 < 0x80 {
		return rune(b0), 1
	}
	if b0 < 0xE0 && len(s) >= 2 {
		return rune(b0&0x1F)<<6 | rune(s[1]&0x3F), 2
	}
	if b0 < 0xF0 && len(s) >= 3 {
		return rune(b0&0x0F)<<12 | rune(s[1]&0x3F)<<6 | rune(s[2]&0x3F), 3
	}
	if len(s) >= 4 {
		return rune(b0&0x07)<<18 | rune(s[1]&0x3F)<<12 | rune(s[2]&0x3F)<<6 | rune(s[3]&0x3F), 4
	}
	return rune(b0), 1
}

func appendWTF8CodePoint(sb *strings.Builder, r rune) {
	switch {
	case r < 0x80:
		sb.WriteByte(byte(r))
	case r < 0x800:
		sb.WriteByte(byte(0xC0 | (r >> 6)))
		sb.WriteByte(byte(0x80 | (r & 0x3F)))
	case r < 0x10000:
		sb.WriteByte(byte(0xE0 | (r >> 12)))
		sb.WriteByte(byte(0x80 | ((r >> 6) & 0x3F)))
		sb.WriteByte(byte(0x80 | (r & 0x3F)))
	default:
		sb.WriteByte(byte(0xF0 | (r >> 18)))
		sb.WriteByte(byte(0x80 | ((r >> 12) & 0x3F)))
		sb.WriteByte(byte(0x80 | ((r >> 6) & 0x3F)))
		sb.WriteByte(byte(0x80 | (r & 0x3F)))
	}
}

func (r *JsonReader) skipToEndOfLine() {
	for r.pos < r.length {
		c := r.json[r.pos]
		r.pos++
		if c == '\n' || c == '\r' {
			break
		}
	}
}

func (r *JsonReader) nextNonWhitespace(throwOnEOF, rootInitialPeek bool) (int, int, int, javalang.Throwable, bool) {
	for r.pos < r.length {
		b := r.json[r.pos]
		if b == ' ' || b == '\r' || b == '\t' || b == '\n' {
			r.pos++
			continue
		}
		if b == '/' {
			r.pos++
			if r.pos == r.length {
				return '/', 1, 0, nil, false
			}
			if err := r.checkLenient(0); err != nil {
				return 0, 0, 0, err, false
			}
			switch r.json[r.pos] {
			case '*':
				r.pos++
				idx := strings.Index(r.json[r.pos:], "*/")
				if idx < 0 {
					r.pos = r.length
					return 0, 0, 0, r.syntaxError("Unterminated comment", 0), false
				}
				r.pos += idx + 2
				continue
			case '/':
				r.pos++
				r.skipToEndOfLine()
				continue
			default:
				return '/', 1, 0, nil, false
			}
		}
		if b == '#' {
			r.pos++
			if err := r.checkLenient(0); err != nil {
				return 0, 0, 0, err, false
			}
			r.skipToEndOfLine()
			continue
		}
		if b < 0x80 {
			r.pos++
			return int(b), 1, 0, nil, false
		}
		cp, w := decodeWTF8At(r.json[r.pos:])
		r.pos += w
		if cp >= 0x10000 {
			hi, _ := utf16.EncodeRune(cp)
			return int(hi), w, -1, nil, false
		}
		return int(cp), w, 0, nil, false
	}
	if !throwOnEOF {
		return -1, 0, 0, nil, false
	}
	if rootInitialPeek {
		return 0, 0, 0, nil, true
	}
	return 0, 0, 0, r.eofError(), false
}

func (r *JsonReader) doPeek(rootInitialPeek bool) (uint8, bool, javalang.Throwable) {
	peekStack := r.topScope()
	switch peekStack {
	case jsonScopeEmptyArray:
		r.setTopScope(jsonScopeNonemptyArray)
	case jsonScopeNonemptyArray:
		c, _, delta, err, _ := r.nextNonWhitespace(true, false)
		if err != nil {
			return 0, false, err
		}
		switch c {
		case ']':
			r.peeked = peekEndArray
			return r.peeked, false, nil
		case ';':
			if err := r.checkLenient(delta); err != nil {
				return 0, false, err
			}
		case ',':
		default:
			return 0, false, r.syntaxError("Unterminated array", delta)
		}
	case jsonScopeEmptyObject, jsonScopeNonemptyObject:
		r.setTopScope(jsonScopeDanglingName)
		if peekStack == jsonScopeNonemptyObject {
			c, _, delta, err, _ := r.nextNonWhitespace(true, false)
			if err != nil {
				return 0, false, err
			}
			switch c {
			case '}':
				r.peeked = peekEndObject
				return r.peeked, false, nil
			case ';':
				if err := r.checkLenient(delta); err != nil {
					return 0, false, err
				}
			case ',':
			default:
				return 0, false, r.syntaxError("Unterminated object", delta)
			}
		}
		c, w, delta, err, _ := r.nextNonWhitespace(true, false)
		if err != nil {
			return 0, false, err
		}
		switch c {
		case '"':
			r.peeked = peekDoubleQuotedName
			return r.peeked, false, nil
		case '\'':
			if err := r.checkLenient(delta); err != nil {
				return 0, false, err
			}
			r.peeked = peekSingleQuotedName
			return r.peeked, false, nil
		case '}':
			if peekStack != jsonScopeNonemptyObject {
				r.peeked = peekEndObject
				return r.peeked, false, nil
			}
			return 0, false, r.syntaxError("Expected name", delta)
		default:
			if err := r.checkLenient(delta); err != nil {
				return 0, false, err
			}
			r.pos -= w
			lit, err := r.isLiteralByte(r.json[r.pos])
			if err != nil {
				return 0, false, err
			}
			if lit {
				r.peeked = peekUnquotedName
				return r.peeked, false, nil
			}
			return 0, false, r.syntaxError("Expected name", 0)
		}
	case jsonScopeDanglingName:
		r.setTopScope(jsonScopeNonemptyObject)
		c, _, delta, err, _ := r.nextNonWhitespace(true, false)
		if err != nil {
			return 0, false, err
		}
		switch c {
		case ':':
		case '=':
			if err := r.checkLenient(delta); err != nil {
				return 0, false, err
			}
			if r.pos < r.length && r.json[r.pos] == '>' {
				r.pos++
			}
		default:
			return 0, false, r.syntaxError("Expected ':'", delta)
		}
	case jsonScopeEmptyDocument:
		if r.lenient {
			_, w, _, err, emptyEOF := r.nextNonWhitespace(true, rootInitialPeek)
			if err != nil || emptyEOF {
				return 0, emptyEOF, err
			}
			r.pos -= w
			if r.pos+5 <= r.length && r.json[r.pos:r.pos+5] == ")]}'\n" {
				r.pos += 5
				r.nonExecNL = r.pos - 1
			}
		}
		r.setTopScope(jsonScopeNonemptyDocument)
	case jsonScopeNonemptyDocument:
		c, w, delta, err, _ := r.nextNonWhitespace(false, false)
		if err != nil {
			return 0, false, err
		}
		if c == -1 {
			r.peeked = peekEOF
			return r.peeked, false, nil
		}
		if err := r.checkLenient(delta); err != nil {
			return 0, false, err
		}
		r.pos -= w
	case jsonScopeClosed:
		return 0, false, javalang.NewIllegalStateException("JsonReader is closed")
	}

	c, w, delta, err, emptyEOF := r.nextNonWhitespace(true, rootInitialPeek)
	if err != nil || emptyEOF {
		return 0, emptyEOF, err
	}
	switch c {
	case ']':
		if peekStack == jsonScopeEmptyArray {
			r.peeked = peekEndArray
			return r.peeked, false, nil
		}
		fallthrough
	case ';', ',':
		if peekStack == jsonScopeEmptyArray || peekStack == jsonScopeNonemptyArray {
			if err := r.checkLenient(delta); err != nil {
				return 0, false, err
			}
			r.pos -= w
			r.peeked = peekNull
			return r.peeked, false, nil
		}
		return 0, false, r.syntaxError("Unexpected value", delta)
	case '\'':
		if err := r.checkLenient(delta); err != nil {
			return 0, false, err
		}
		r.peeked = peekSingleQuoted
		return r.peeked, false, nil
	case '"':
		r.peeked = peekDoubleQuoted
		return r.peeked, false, nil
	case '[':
		r.peeked = peekBeginArray
		return r.peeked, false, nil
	case '{':
		r.peeked = peekBeginObject
		return r.peeked, false, nil
	default:
		r.pos -= w
	}

	kw, err := r.peekKeyword()
	if err != nil || kw != peekNone {
		return kw, false, err
	}
	num, err := r.peekNumber()
	if err != nil || num != peekNone {
		return num, false, err
	}
	lit, err := r.isLiteralByte(r.json[r.pos])
	if err != nil {
		return 0, false, err
	}
	if !lit {
		return 0, false, r.syntaxError("Expected value", 0)
	}
	if err := r.checkLenient(0); err != nil {
		return 0, false, err
	}
	r.peeked = peekUnquoted
	return r.peeked, false, nil
}

func (r *JsonReader) peekKeyword() (uint8, javalang.Throwable) {
	b := r.json[r.pos]
	var kwLower, kwUpper string
	var peeking uint8
	switch b {
	case 't', 'T':
		kwLower, kwUpper, peeking = "true", "TRUE", peekTrue
	case 'f', 'F':
		kwLower, kwUpper, peeking = "false", "FALSE", peekFalse
	case 'n', 'N':
		kwLower, kwUpper, peeking = "null", "NULL", peekNull
	default:
		return peekNone, nil
	}
	n := len(kwLower)
	if r.pos+n > r.length {
		return peekNone, nil
	}
	for i := 1; i < n; i++ {
		c := r.json[r.pos+i]
		if c != kwLower[i] && c != kwUpper[i] {
			return peekNone, nil
		}
	}
	if r.pos+n < r.length {
		lit, err := r.isLiteralByte(r.json[r.pos+n])
		if err != nil || lit {
			return peekNone, err
		}
	}
	r.pos += n
	r.peeked = peeking
	return peeking, nil
}

func (r *JsonReader) peekNumber() (uint8, javalang.Throwable) {
	str := r.json
	p := r.pos
	end := r.length

	var value int64
	negative := false
	fitsInLong := true
	last := numNone
	i := 0

loop:
	for ; ; i++ {
		if p+i == end {
			break
		}
		if i == maxNumberChars {
			return peekNone, nil
		}
		c := str[p+i]
		switch c {
		case '-':
			if last == numNone {
				negative = true
				last = numSign
				continue
			} else if last == numExpE {
				last = numExpSign
				continue
			}
			return peekNone, nil
		case '+':
			if last == numExpE {
				last = numExpSign
				continue
			}
			return peekNone, nil
		case 'e', 'E':
			if last == numDigit || last == numFractionDigit {
				last = numExpE
				continue
			}
			return peekNone, nil
		case '.':
			if last == numDigit {
				last = numDecimal
				continue
			}
			return peekNone, nil
		default:
			if c < '0' || c > '9' {
				lit, err := r.isLiteralByte(c)
				if err != nil {
					return 0, err
				}
				if !lit {
					break loop
				}
				return peekNone, nil
			}
			if last == numSign || last == numNone {
				value = -int64(c - '0')
				last = numDigit
			} else if last == numDigit {
				if value == 0 {
					return peekNone, nil
				}
				newValue := value*10 - int64(c-'0')
				fitsInLong = fitsInLong && (value > minIncompleteInt64 ||
					(value == minIncompleteInt64 && newValue < value))
				value = newValue
			} else if last == numDecimal {
				last = numFractionDigit
			} else if last == numExpE || last == numExpSign {
				last = numExpDigit
			}
		}
	}

	if last == numDigit && fitsInLong && (value != math.MinInt64 || negative) && (value != 0 || !negative) {
		if negative {
			r.peekedLong = value
		} else {
			r.peekedLong = -value
		}
		r.pos += i
		r.peeked = peekLong
		return peekLong, nil
	}
	if last == numDigit || last == numFractionDigit || last == numExpDigit {
		r.peekedNumberLength = i
		r.peeked = peekNumber
		return peekNumber, nil
	}
	return peekNone, nil
}

func (r *JsonReader) peekCode() (uint8, javalang.Throwable) {
	if r.peeked != peekNone {
		return r.peeked, nil
	}
	p, _, err := r.doPeek(false)
	return p, err
}

func (r *JsonReader) readUnquotedValue() (string, javalang.Throwable) {
	start := r.pos
	for i := start; i < r.length; i++ {
		b := r.json[i]
		if b < 0x80 && nonLiteralLUT[b] {
			if _, err := r.isLiteralByte(b); err != nil {
				return "", err
			}
			r.pos = i
			return r.json[start:i], nil
		}
	}
	r.pos = r.length
	return r.json[start:r.length], nil
}

func writeWTF8Unit(sb *strings.Builder, pendingHi *uint16, u uint16) {
	if *pendingHi != 0 {
		if u >= 0xDC00 && u <= 0xDFFF {
			sb.WriteRune(utf16.DecodeRune(rune(*pendingHi), rune(u)))
			*pendingHi = 0
			return
		}
		hi := *pendingHi
		*pendingHi = 0
		appendWTF8CodePoint(sb, rune(hi))
	}
	if u >= 0xD800 && u <= 0xDBFF {
		*pendingHi = u
		return
	}
	appendWTF8CodePoint(sb, rune(u))
}

func flushPendingSurrogate(sb *strings.Builder, pendingHi *uint16) {
	if *pendingHi != 0 {
		hi := *pendingHi
		*pendingHi = 0
		appendWTF8CodePoint(sb, rune(hi))
	}
}

func (r *JsonReader) readQuotedValue(quote byte) (string, javalang.Throwable) {
	start := r.pos
	for r.pos < r.length {
		c := r.json[r.pos]
		if c == quote {
			res := r.json[start:r.pos]
			r.pos++
			return res, nil
		}
		if c == '\\' {
			break
		}
		r.pos++
	}
	if r.pos >= r.length {
		return "", r.syntaxError("Unterminated string", 0)
	}

	var sb strings.Builder
	sb.Grow((r.pos - start) + 16)
	var pendingHi uint16

	appendSeg := func(seg string) {
		if pendingHi == 0 && (len(seg) < 3 || seg[len(seg)-3] != 0xED) {
			sb.WriteString(seg)
			return
		}
		for i := 0; i < len(seg); {
			if seg[i] < 0x80 {
				flushPendingSurrogate(&sb, &pendingHi)
				j := i + 1
				for j < len(seg) && seg[j] < 0x80 {
					j++
				}
				sb.WriteString(seg[i:j])
				i = j
				continue
			}
			cp, w := decodeWTF8At(seg[i:])
			i += w
			if cp >= 0x10000 {
				flushPendingSurrogate(&sb, &pendingHi)
				sb.WriteRune(cp)
			} else {
				writeWTF8Unit(&sb, &pendingHi, uint16(cp))
			}
		}
	}

	appendSeg(r.json[start:r.pos])
	for r.pos < r.length {
		c := r.json[r.pos]
		if c == quote {
			r.pos++
			flushPendingSurrogate(&sb, &pendingHi)
			return sb.String(), nil
		}
		if c == '\\' {
			r.pos++
			u, err := r.readEscapeChar()
			if err != nil {
				return "", err
			}
			writeWTF8Unit(&sb, &pendingHi, u)
			segStart := r.pos
			for r.pos < r.length {
				b := r.json[r.pos]
				if b == quote || b == '\\' {
					break
				}
				r.pos++
			}
			if segStart < r.pos {
				appendSeg(r.json[segStart:r.pos])
			}
			continue
		}
		r.pos++
	}
	return "", r.syntaxError("Unterminated string", 0)
}

func parseHex4(s string) (uint16, bool) {
	var res uint16
	for i := 0; i < 4; i++ {
		c := s[i]
		var d uint16
		switch {
		case c >= '0' && c <= '9':
			d = uint16(c - '0')
		case c >= 'a' && c <= 'f':
			d = uint16(c-'a') + 10
		case c >= 'A' && c <= 'F':
			d = uint16(c-'A') + 10
		default:
			return 0, false
		}
		res = (res << 4) | d
	}
	return res, true
}

func (r *JsonReader) readEscapeChar() (uint16, javalang.Throwable) {
	if r.pos >= r.length {
		return 0, r.syntaxError("Unterminated escape sequence", 0)
	}
	b := r.json[r.pos]
	if b >= 0x80 {
		return 0, r.syntaxError("Invalid escape sequence", 0)
	}
	r.pos++
	switch b {
	case 'u':
		if r.pos+4 > r.length {
			return 0, r.syntaxError("Unterminated escape sequence", 0)
		}
		hexSeg := r.json[r.pos : r.pos+4]
		res, ok := parseHex4(hexSeg)
		if !ok {
			return 0, javalang.NewNumberFormatException("\\u" + hexSeg)
		}
		r.pos += 4
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
		return 0, r.syntaxError("Invalid escape sequence", 0)
	}
}

func mustNoErr(err javalang.Throwable) {
	if err != nil {
		panic(err)
	}
}

func (r *JsonReader) Peek() JsonToken {
	tok, err := r.peekTokenErr()
	mustNoErr(err)
	return tok
}

func (r *JsonReader) peekTokenErr() (JsonToken, javalang.Throwable) {
	p, err := r.peekCode()
	if err != nil {
		return 0, err
	}
	switch p {
	case peekBeginObject:
		return JsonToken_BEGIN_OBJECT, nil
	case peekEndObject:
		return JsonToken_END_OBJECT, nil
	case peekBeginArray:
		return JsonToken_BEGIN_ARRAY, nil
	case peekEndArray:
		return JsonToken_END_ARRAY, nil
	case peekSingleQuotedName, peekDoubleQuotedName, peekUnquotedName:
		return JsonToken_NAME, nil
	case peekTrue, peekFalse:
		return JsonToken_BOOLEAN, nil
	case peekNull:
		return JsonToken_NULL, nil
	case peekSingleQuoted, peekDoubleQuoted, peekUnquoted:
		return JsonToken_STRING, nil
	case peekLong, peekNumber:
		return JsonToken_NUMBER, nil
	case peekEOF:
		return JsonToken_END_DOCUMENT, nil
	default:
		return 0, javalang.NewAssertionError(nil)
	}
}

func (r *JsonReader) unexpectedErr(expected string) javalang.Throwable {
	tok, err := r.peekTokenErr()
	if err != nil {
		return err
	}
	return javalang.NewIllegalStateException("Expected " + expected + " but was " + tok.String() + r.locationString(r.pos, 0))
}

func (r *JsonReader) HasNext() bool {
	p, err := r.peekCode()
	mustNoErr(err)
	return p != peekEndObject && p != peekEndArray && p != peekEOF
}

func (r *JsonReader) BeginArray() {
	p, err := r.peekCode()
	mustNoErr(err)
	if p != peekBeginArray {
		panic(r.unexpectedErr("BEGIN_ARRAY"))
	}
	r.pushScope(jsonScopeEmptyArray)
	r.peeked = peekNone
}

func (r *JsonReader) EndArray() {
	p, err := r.peekCode()
	mustNoErr(err)
	if p != peekEndArray {
		panic(r.unexpectedErr("END_ARRAY"))
	}
	r.stack = r.stack[:len(r.stack)-1]
	r.incTopIndex()
	r.peeked = peekNone
}

func (r *JsonReader) BeginObject() {
	p, err := r.peekCode()
	mustNoErr(err)
	if p != peekBeginObject {
		panic(r.unexpectedErr("BEGIN_OBJECT"))
	}
	r.pushScope(jsonScopeEmptyObject)
	r.peeked = peekNone
}

func (r *JsonReader) EndObject() {
	p, err := r.peekCode()
	mustNoErr(err)
	if p != peekEndObject {
		panic(r.unexpectedErr("END_OBJECT"))
	}
	r.stack[len(r.stack)-1] = readerFrame{}
	r.stack = r.stack[:len(r.stack)-1]
	r.incTopIndex()
	r.peeked = peekNone
}

func (r *JsonReader) NextName() string {
	p, err := r.peekCode()
	mustNoErr(err)
	var name string
	switch p {
	case peekDoubleQuotedName:
		name, err = r.readQuotedValue('"')
	case peekSingleQuotedName:
		name, err = r.readQuotedValue('\'')
	case peekUnquotedName:
		name, err = r.readUnquotedValue()
	default:
		panic(r.unexpectedErr("a name"))
	}
	mustNoErr(err)
	r.peeked = peekNone
	r.setTopName(name)
	return name
}

func (r *JsonReader) NextString() string {
	p, err := r.peekCode()
	mustNoErr(err)
	var res string
	switch p {
	case peekDoubleQuoted:
		res, err = r.readQuotedValue('"')
	case peekSingleQuoted:
		res, err = r.readQuotedValue('\'')
	case peekUnquoted:
		res, err = r.readUnquotedValue()
	case peekLong:
		res = strconv.FormatInt(r.peekedLong, 10)
	case peekNumber:
		res = r.json[r.pos : r.pos+r.peekedNumberLength]
		r.pos += r.peekedNumberLength
	default:
		panic(r.unexpectedErr("a string"))
	}
	mustNoErr(err)
	r.peeked = peekNone
	r.incTopIndex()
	return res
}

func (r *JsonReader) NextBoolean() bool {
	p, err := r.peekCode()
	mustNoErr(err)
	switch p {
	case peekTrue:
		r.peeked = peekNone
		r.incTopIndex()
		return true
	case peekFalse:
		r.peeked = peekNone
		r.incTopIndex()
		return false
	default:
		panic(r.unexpectedErr("a boolean"))
	}
}

func (r *JsonReader) NextNull() {
	p, err := r.peekCode()
	mustNoErr(err)
	if p != peekNull {
		panic(r.unexpectedErr("null"))
	}
	r.peeked = peekNone
	r.incTopIndex()
}

func (r *JsonReader) NextLong() int64 {
	p, err := r.peekCode()
	mustNoErr(err)
	if p == peekLong {
		r.peeked = peekNone
		r.incTopIndex()
		return r.peekedLong
	}
	s := r.NextString()
	v, parseErr := strconv.ParseInt(s, 10, 64)
	if parseErr == nil {
		return v
	}
	d, fErr := strconv.ParseFloat(s, 64)
	if fErr != nil || float64(int64(d)) != d {
		panic(javalang.NewNumberFormatException("Expected a long but was " + s + r.locationString(r.pos, 0)))
	}
	return int64(d)
}

func (r *JsonReader) NextInt() int32 {
	v := r.NextLong()
	if v < math.MinInt32 || v > math.MaxInt32 {
		panic(javalang.NewNumberFormatException("Expected an int but was " + strconv.FormatInt(v, 10) + r.locationString(r.pos, 0)))
	}
	return int32(v)
}

func (r *JsonReader) NextDouble() float64 {
	p, err := r.peekCode()
	mustNoErr(err)
	if p == peekLong {
		r.peeked = peekNone
		r.incTopIndex()
		return float64(r.peekedLong)
	}
	s := r.NextString()
	d, fErr := strconv.ParseFloat(s, 64)
	if fErr != nil {
		panic(javalang.NewNumberFormatException("Expected a double but was " + s + r.locationString(r.pos, 0)))
	}
	if !r.lenient && (math.IsNaN(d) || math.IsInf(d, 0)) {
		panic(NewMalformedJsonException("JSON forbids NaN and infinities: " + strconv.FormatFloat(d, 'g', -1, 64) + r.locationString(r.pos, 0)))
	}
	return d
}

func (r *JsonReader) SkipValue() {
	count := 0
	for {
		p, err := r.peekCode()
		mustNoErr(err)
		switch p {
		case peekBeginArray:
			r.pushScope(jsonScopeEmptyArray)
			count++
		case peekBeginObject:
			r.pushScope(jsonScopeEmptyObject)
			count++
		case peekEndArray, peekEndObject:
			r.stack = r.stack[:len(r.stack)-1]
			count--
		case peekUnquotedName, peekUnquoted:
			_, err = r.readUnquotedValue()
			mustNoErr(err)
		case peekSingleQuoted, peekSingleQuotedName:
			_, err = r.readQuotedValue('\'')
			mustNoErr(err)
		case peekDoubleQuoted, peekDoubleQuotedName:
			_, err = r.readQuotedValue('"')
			mustNoErr(err)
		case peekNumber:
			r.pos += r.peekedNumberLength
		}
		r.peeked = peekNone
		if count == 0 {
			break
		}
	}
	r.incTopIndex()
	r.setTopName("null")
}
