// Port of com/google/protobuf/TextFormat.java @ protobuf-java 4.30.2 (the text-format parser)

// Protocol Buffers - Google's data interchange format
// Copyright 2008 Google Inc.  All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd
//
// Here that LICENSE file is LICENSES/BSD-3-Clause-protobuf.txt.

// Go-only: language-adapter — protobuf is a leaf package (PORTING.md §2) and
// its messages are reconstructed (proto_runtime.go, which already prints the
// text format). This file translates Tokenizer, ParseException, and
// unescapeBytes from com.google.protobuf.TextFormat (the version shaded into
// the pinned closure-compiler.jar).

package jscomp

import (
	"math"
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// TextFormat_Tokenizer represents a stream of tokens parsed from a String.
//
// The Java standard library provides many classes that you might think would
// be useful for implementing this, but aren't. For example:
//
//   - java.io.StreamTokenizer: This almost does what we want -- or, at least,
//     something that would get us close to what we want -- except for one
//     fatal flaw: It automatically un-escapes strings using Java escape
//     sequences, which do not include all the escape sequences we need to
//     support (e.g. '\x').
//   - java.util.Scanner: This seems like a great way at least to parse
//     regular expressions out of a stream (so we wouldn't have to load the
//     entire input into a single string before parsing). Sadly, Scanner
//     requires that tokens be delimited with some delimiter. Thus, although
//     the text "foo:" should parse to two tokens ("foo" and ":"), Scanner
//     would recognize it only as a single token. Furthermore, Scanner
//     provides no way to inspect the contents of delimiters, making it
//     impossible to keep track of line and column numbers.
//
// Luckily, Java's regular expression support does manage to be useful to us.
// (Barely: We need Matcher.usePattern(), which is new in Java 1.5.) So, we
// can use it to implement a custom Tokenizer without too much trouble.
type TextFormat_Tokenizer struct {
	// Go-only: language-adapter — the CharSequence as UTF-16 code units.
	text         []uint16
	currentToken string

	// The character index within this.text at which the current token begins.
	pos int

	// The line and column numbers of the current token.
	line                int
	column              int
	lineInfoTrackingPos int

	// The line and column numbers of the previous token (allows throwing
	// errors *after* consuming).
	previousLine   int
	previousColumn int
}

// NewTextFormat_Tokenizer constructs a tokenizer that parses tokens from the
// given text.
//
// Java: TextFormat.Tokenizer(CharSequence)
func NewTextFormat_Tokenizer(text string) *TextFormat_Tokenizer {
	t := &TextFormat_Tokenizer{}
	t.text = jutil.StringToUTF16(text)
	t.SkipWhitespace()
	t.NextToken()
	return t
}

// Java: TextFormat.Tokenizer.getPreviousLine()
func (t *TextFormat_Tokenizer) GetPreviousLine() int { return t.previousLine }

// Java: TextFormat.Tokenizer.getPreviousColumn()
func (t *TextFormat_Tokenizer) GetPreviousColumn() int { return t.previousColumn }

// AtEnd: are we at the end of the input?
//
// Java: TextFormat.Tokenizer.atEnd()
func (t *TextFormat_Tokenizer) AtEnd() bool {
	return len(t.currentToken) == 0
}

// NextToken advances to the next token.
//
// Java: TextFormat.Tokenizer.nextToken()
func (t *TextFormat_Tokenizer) NextToken() {
	t.previousLine = t.line
	t.previousColumn = t.column

	// Advance the line counter to the current position.
	for t.lineInfoTrackingPos < t.pos {
		if t.text[t.lineInfoTrackingPos] == '\n' {
			t.line++
			t.column = 0
		} else {
			t.column++
		}
		t.lineInfoTrackingPos++
	}

	// Match the next token.
	if t.pos == len(t.text) {
		t.currentToken = "" // EOF
	} else {
		t.currentToken = t.NextTokenInternal()
		t.SkipWhitespace()
	}
}

// Java: TextFormat.Tokenizer.nextTokenInternal()
func (t *TextFormat_Tokenizer) NextTokenInternal() string {
	textLength := len(t.text)
	startPos := t.pos
	startChar := t.text[startPos]

	endPos := t.pos
	if TextFormat_Tokenizer_IsAlphaUnder(startChar) { // Identifier
		for {
			endPos++
			if endPos == textLength {
				break
			}
			c := t.text[endPos]
			if !(TextFormat_Tokenizer_IsAlphaUnder(c) || TextFormat_Tokenizer_IsDigitPlusMinus(c)) {
				break
			}
		}
	} else if TextFormat_Tokenizer_IsDigitPlusMinus(startChar) || startChar == '.' { // Number
		if startChar == '.' { // Optional leading dot
			endPos++
			if endPos == textLength {
				return t.NextTokenSingleChar()
			}

			if !TextFormat_Tokenizer_IsDigitPlusMinus(t.text[endPos]) { // Mandatory first digit
				return t.NextTokenSingleChar()
			}
		}

		for {
			endPos++
			if endPos == textLength {
				break
			}
			c := t.text[endPos]
			if !(TextFormat_Tokenizer_IsDigitPlusMinus(c) || TextFormat_Tokenizer_IsAlphaUnder(c) || c == '.') {
				break
			}
		}
	} else if startChar == '"' || startChar == '\'' { // String
		for {
			endPos++
			if endPos == textLength {
				break
			}
			c := t.text[endPos]
			if c == startChar {
				endPos++
				break // Quote terminates
			} else if c == '\n' {
				break // Newline terminates (error during parsing) (not consumed)
			} else if c == '\\' {
				endPos++
				if endPos == textLength {
					break // Escape into end-of-text terminates (error during parsing)
				} else if t.text[endPos] == '\n' {
					break // Escape into newline terminates (error during parsing) (not consumed)
				} else {
					// Otherwise the escaped char is legal and consumed
				}
			} else {
				// Otherwise the char is a legal and consumed
			}
		}
	} else {
		return t.NextTokenSingleChar() // Unrecognized start character
	}

	t.pos = endPos
	return jutil.StringFromUTF16(t.text[startPos:endPos])
}

// Java: TextFormat.Tokenizer.isAlphaUnder(char)
func TextFormat_Tokenizer_IsAlphaUnder(c uint16) bool {
	// Defining this char-class with numeric comparisons is much faster than using a regex.
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || c == '_'
}

// Java: TextFormat.Tokenizer.isDigitPlusMinus(char)
func TextFormat_Tokenizer_IsDigitPlusMinus(c uint16) bool {
	// Defining this char-class with numeric comparisons is much faster than using a regex.
	return ('0' <= c && c <= '9') || c == '+' || c == '-'
}

// Java: TextFormat.Tokenizer.isWhitespace(char)
func TextFormat_Tokenizer_IsWhitespace(c uint16) bool {
	// Defining this char-class with numeric comparisons is much faster than using a regex.
	return c == ' ' || c == '\f' || c == '\n' || c == '\r' || c == '\t'
}

// NextTokenSingleChar produces a token for the single char at the current
// position.
//
// This method must not be called if the current position is after the
// end-of-text.
//
// Java: TextFormat.Tokenizer.nextTokenSingleChar()
func (t *TextFormat_Tokenizer) NextTokenSingleChar() string {
	c := t.text[t.pos]
	t.pos++
	return jutil.StringFromUTF16([]uint16{c})
}

// SkipWhitespace skips over any whitespace so that the matcher region starts
// at the next token.
//
// Java: TextFormat.Tokenizer.skipWhitespace()
func (t *TextFormat_Tokenizer) SkipWhitespace() {
	textLength := len(t.text)

	endPos := t.pos - 1
	for {
		endPos++
		if endPos == textLength {
			break
		}
		c := t.text[endPos]
		if c == '#' {
			for {
				endPos++
				if endPos == textLength {
					break
				}
				if t.text[endPos] == '\n' {
					break // Consume the newline as whitespace.
				}
			}
			if endPos == textLength {
				break
			}
		} else if TextFormat_Tokenizer_IsWhitespace(c) {
			// OK
		} else {
			break
		}
	}

	t.pos = endPos
}

// TryConsume: if the next token exactly matches token, consume it and return
// true. Otherwise, return false without doing anything.
//
// Java: TextFormat.Tokenizer.tryConsume(String)
func (t *TextFormat_Tokenizer) TryConsume(token string) bool {
	if t.currentToken == token {
		t.NextToken()
		return true
	} else {
		return false
	}
}

// Consume: if the next token exactly matches token, consume it. Otherwise,
// return a ParseException.
//
// Java: TextFormat.Tokenizer.consume(String)
func (t *TextFormat_Tokenizer) Consume(token string) *TextFormat_ParseException {
	if !t.TryConsume(token) {
		return t.ParseException("Expected \"" + token + "\".")
	}
	return nil
}

// LookingAtInteger returns true if the next token is an integer, but does not
// consume it.
//
// Java: TextFormat.Tokenizer.lookingAtInteger()
func (t *TextFormat_Tokenizer) LookingAtInteger() bool {
	if len(t.currentToken) == 0 {
		return false
	}

	return TextFormat_Tokenizer_IsDigitPlusMinus(uint16(t.currentToken[0]))
}

// LookingAt returns true if the current token's text is equal to that
// specified.
//
// Java: TextFormat.Tokenizer.lookingAt(String)
func (t *TextFormat_Tokenizer) LookingAt(text string) bool {
	return t.currentToken == text
}

func (t *TextFormat_Tokenizer) isIdentifier() bool {
	for _, c := range jutil.StringToUTF16(t.currentToken) {
		if !(TextFormat_Tokenizer_IsAlphaUnder(c) || ('0' <= c && c <= '9') || (c == '.')) {
			return false
		}
	}
	return true
}

// ConsumeIdentifier: if the next token is an identifier, consume it and
// return its value. Otherwise, return a ParseException.
//
// Java: TextFormat.Tokenizer.consumeIdentifier()
func (t *TextFormat_Tokenizer) ConsumeIdentifier() (string, *TextFormat_ParseException) {
	if !t.isIdentifier() {
		return "", t.ParseException("Expected identifier. Found '" + t.currentToken + "'")
	}

	result := t.currentToken
	t.NextToken()
	return result, nil
}

// TryConsumeIdentifier: if the next token is an identifier, consume it and
// return true. Otherwise, return false without doing anything.
//
// Java: TextFormat.Tokenizer.tryConsumeIdentifier()
func (t *TextFormat_Tokenizer) TryConsumeIdentifier() bool {
	if !t.isIdentifier() {
		return false
	}
	t.NextToken()
	return true
}

// ConsumeInt32: if the next token is a 32-bit signed integer, consume it and
// return its value. Otherwise, return a ParseException.
//
// Java: TextFormat.Tokenizer.consumeInt32()
func (t *TextFormat_Tokenizer) ConsumeInt32() (int32, *TextFormat_ParseException) {
	result, errMsg, ok := textFormat_parseInteger(t.currentToken, true, false)
	if !ok {
		return 0, t.ParseException("Couldn't parse integer: " + errMsg)
	}
	t.NextToken()
	return int32(result), nil
}

// TryConsumeInt64: if the next token is a 64-bit signed integer, consume it
// and return true. Otherwise, return false without doing anything.
//
// Java: TextFormat.Tokenizer.tryConsumeInt64()
func (t *TextFormat_Tokenizer) TryConsumeInt64() bool {
	if _, _, ok := textFormat_parseInteger(t.currentToken, true, true); !ok {
		return false
	}
	t.NextToken()
	return true
}

// TryConsumeUInt64: if the next token is a 64-bit unsigned integer, consume
// it and return true. Otherwise, return false without doing anything.
//
// Java: TextFormat.Tokenizer.tryConsumeUInt64()
func (t *TextFormat_Tokenizer) TryConsumeUInt64() bool {
	if _, _, ok := textFormat_parseInteger(t.currentToken, false, true); !ok {
		return false
	}
	t.NextToken()
	return true
}

func (t *TextFormat_Tokenizer) parseDoubleToken() (float64, string, bool) {
	// We need to parse infinity and nan separately because
	// Double.parseDouble() does not accept "inf", "infinity", or "nan".
	switch jutil.StringToLowerCaseRoot(t.currentToken) {
	case "-inf", "-infinity":
		return math.Inf(-1), "", true
	case "inf", "infinity":
		return math.Inf(1), "", true
	case "nan":
		return jutil.Double_NaN, "", true
	}
	return textFormat_parseJavaDouble(t.currentToken)
}

// TryConsumeDouble: if the next token is a double, consume it and return
// true. Otherwise, return false without doing anything.
//
// Java: TextFormat.Tokenizer.tryConsumeDouble()
func (t *TextFormat_Tokenizer) TryConsumeDouble() bool {
	if _, _, ok := t.parseDoubleToken(); !ok {
		return false
	}
	t.NextToken()
	return true
}

func (t *TextFormat_Tokenizer) parseFloatToken() (float32, string, bool) {
	// We need to parse infinity and nan separately because
	// Float.parseFloat() does not accept "inf", "infinity", or "nan".
	switch jutil.StringToLowerCaseRoot(t.currentToken) {
	case "-inf", "-inff", "-infinity", "-infinityf":
		return float32(math.Inf(-1)), "", true
	case "inf", "inff", "infinity", "infinityf":
		return float32(math.Inf(1)), "", true
	case "nan", "nanf":
		return float32(jutil.Double_NaN), "", true
	}
	v, errMsg, ok := textFormat_parseJavaDouble(t.currentToken)
	if !ok {
		return 0, errMsg, false
	}
	return float32(v), "", true
}

// TryConsumeFloat: if the next token is a float, consume it and return true.
// Otherwise, return false without doing anything.
//
// Java: TextFormat.Tokenizer.tryConsumeFloat()
func (t *TextFormat_Tokenizer) TryConsumeFloat() bool {
	if _, _, ok := t.parseFloatToken(); !ok {
		return false
	}
	t.NextToken()
	return true
}

// ConsumeBoolean: if the next token is a boolean, consume it and return its
// value. Otherwise, return a ParseException.
//
// Java: TextFormat.Tokenizer.consumeBoolean()
func (t *TextFormat_Tokenizer) ConsumeBoolean() (bool, *TextFormat_ParseException) {
	if t.currentToken == "true" ||
		t.currentToken == "True" ||
		t.currentToken == "t" ||
		t.currentToken == "1" {
		t.NextToken()
		return true, nil
	} else if t.currentToken == "false" ||
		t.currentToken == "False" ||
		t.currentToken == "f" ||
		t.currentToken == "0" {
		t.NextToken()
		return false, nil
	}
	return false, t.ParseException("Expected \"true\" or \"false\". Found \"" + t.currentToken + "\".")
}

// ConsumeString: if the next token is a string, consume it and return its
// (unescaped) value. Otherwise, return a ParseException.
//
// Java: TextFormat.Tokenizer.consumeString()
func (t *TextFormat_Tokenizer) ConsumeString() (string, *TextFormat_ParseException) {
	b, err := t.ConsumeByteString()
	if err != nil {
		return "", err
	}
	return jutil.NewStringFromUTF8Bytes(b), nil
}

// ConsumeByteString: if the next token is a string, consume it, unescape it
// as a ByteString, and return it. Otherwise, return a ParseException.
//
// Java: TextFormat.Tokenizer.consumeByteString()
func (t *TextFormat_Tokenizer) ConsumeByteString() ([]byte, *TextFormat_ParseException) {
	list := [][]byte{}
	if err := t.consumeByteStringInto(&list); err != nil {
		return nil, err
	}
	for strings.HasPrefix(t.currentToken, "'") || strings.HasPrefix(t.currentToken, "\"") {
		if err := t.consumeByteStringInto(&list); err != nil {
			return nil, err
		}
	}
	var result []byte
	for _, b := range list {
		result = append(result, b...)
	}
	return result, nil
}

// TryConsumeByteString: if the next token is a string, consume it and return
// true. Otherwise, return false.
//
// Java: TextFormat.Tokenizer.tryConsumeByteString()
func (t *TextFormat_Tokenizer) TryConsumeByteString() bool {
	if _, _, ok := t.parseByteStringToken(); !ok {
		return false
	}
	t.NextToken()
	for strings.HasPrefix(t.currentToken, "'") || strings.HasPrefix(t.currentToken, "\"") {
		if _, _, ok := t.parseByteStringToken(); !ok {
			return false
		}
		t.NextToken()
	}
	return true
}

// consumeByteStringInto is like ConsumeByteString but adds each token of the
// string to the given list. String literals (whether bytes or text) may come
// in multiple adjacent tokens which are automatically concatenated, like in
// C or Python.
//
// Java: TextFormat.Tokenizer.consumeByteString(List<ByteString>)
func (t *TextFormat_Tokenizer) consumeByteStringInto(list *[][]byte) *TextFormat_ParseException {
	result, errMsg, ok := t.parseByteStringToken()
	if !ok {
		return t.ParseException(errMsg)
	}
	t.NextToken()
	*list = append(*list, result)
	return nil
}

func (t *TextFormat_Tokenizer) parseByteStringToken() ([]byte, string, bool) {
	var quote byte
	if len(t.currentToken) > 0 {
		quote = t.currentToken[0]
	}
	if quote != '"' && quote != '\'' {
		return nil, "Expected string.", false
	}

	if len(t.currentToken) < 2 || t.currentToken[len(t.currentToken)-1] != quote {
		return nil, "String missing ending quote.", false
	}

	escaped := t.currentToken[1 : len(t.currentToken)-1]
	return textFormat_unescapeBytes(escaped)
}

func textFormat_unescapeBytes(charString string) ([]byte, string, bool) {
	input := jutil.StringGetBytesUTF8(charString)
	result := make([]byte, 0, len(input))
	for i := 0; i < len(input); i++ {
		c := input[i]
		if c == '\\' {
			if i+1 < len(input) {
				i++
				c = input[i]
				if textFormat_isOctal(c) {
					code := textFormat_digitValue(c)
					if i+1 < len(input) && textFormat_isOctal(input[i+1]) {
						i++
						code = code*8 + textFormat_digitValue(input[i])
					}
					if i+1 < len(input) && textFormat_isOctal(input[i+1]) {
						i++
						code = code*8 + textFormat_digitValue(input[i])
					}
					result = append(result, byte(code))
				} else {
					switch c {
					case 'a':
						result = append(result, 0x07)
					case 'b':
						result = append(result, '\b')
					case 'f':
						result = append(result, '\f')
					case 'n':
						result = append(result, '\n')
					case 'r':
						result = append(result, '\r')
					case 't':
						result = append(result, '\t')
					case 'v':
						result = append(result, 0x0b)
					case '\\':
						result = append(result, '\\')
					case '\'':
						result = append(result, '\'')
					case '"':
						result = append(result, '"')
					case '?':
						result = append(result, '?')

					case 'x':
						code := 0
						if i+1 < len(input) && textFormat_isHex(input[i+1]) {
							i++
							code = textFormat_digitValue(input[i])
						} else {
							return nil, "Invalid escape sequence: '\\x' with no digits", false
						}
						if i+1 < len(input) && textFormat_isHex(input[i+1]) {
							i++
							code = code*16 + textFormat_digitValue(input[i])
						}
						result = append(result, byte(code))

					case 'u':
						i++
						if i+3 < len(input) &&
							textFormat_isHex(input[i]) &&
							textFormat_isHex(input[i+1]) &&
							textFormat_isHex(input[i+2]) &&
							textFormat_isHex(input[i+3]) {
							ch := uint16(textFormat_digitValue(input[i])<<12 |
								textFormat_digitValue(input[i+1])<<8 |
								textFormat_digitValue(input[i+2])<<4 |
								textFormat_digitValue(input[i+3]))

							if ch >= 0xD800 && ch <= 0xDFFF {
								return nil, "Invalid escape sequence: '\\u' refers to a surrogate", false
							}
							result = append(result, jutil.StringGetBytesUTF8(jutil.StringFromUTF16([]uint16{ch}))...)
							i += 3
						} else {
							return nil, "Invalid escape sequence: '\\u' with too few hex chars", false
						}

					case 'U':
						i++
						if i+7 >= len(input) {
							return nil, "Invalid escape sequence: '\\U' with too few hex chars", false
						}
						codepoint := 0
						for offset := i; offset < i+8; offset++ {
							b := input[offset]
							if !textFormat_isHex(b) {
								return nil, "Invalid escape sequence: '\\U' with too few hex chars", false
							}
							codepoint = (codepoint << 4) | textFormat_digitValue(b)
						}
						// Character.isValidCodePoint; the hex digits make it non-negative
						// as an int only below 0x80000000.
						if int32(codepoint) < 0 || codepoint > 0x10FFFF {
							return nil, "Invalid escape sequence: '\\U" +
								jutil.NewStringFromUTF8Bytes(input[i:i+8]) +
								"' is not a valid code point value", false
						}
						// Character.UnicodeBlock LOW_SURROGATES, HIGH_SURROGATES,
						// HIGH_PRIVATE_USE_SURROGATES: U+D800..U+DFFF.
						if codepoint >= 0xD800 && codepoint <= 0xDFFF {
							return nil, "Invalid escape sequence: '\\U" +
								jutil.NewStringFromUTF8Bytes(input[i:i+8]) +
								"' refers to a surrogate code unit", false
						}
						result = append(result, []byte(string(rune(codepoint)))...)
						i += 7

					default:
						return nil, "Invalid escape sequence: '\\" + string(rune(c)) + "'", false
					}
				}
			} else {
				return nil, "Invalid escape sequence: '\\' at end of string.", false
			}
		} else {
			result = append(result, c)
		}
	}

	return result, "", true
}

// Java: TextFormat.isOctal(byte)
func textFormat_isOctal(c byte) bool {
	return '0' <= c && c <= '7'
}

// Java: TextFormat.isHex(byte)
func textFormat_isHex(c byte) bool {
	return ('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}

// textFormat_digitValue interprets a character as a digit (in any base up to
// 36) and returns the numeric value. This is like Character.digit() but we
// don't accept non-ASCII digits.
//
// Java: TextFormat.digitValue(byte)
func textFormat_digitValue(c byte) int {
	if '0' <= c && c <= '9' {
		return int(c - '0')
	} else if 'a' <= c && c <= 'z' {
		return int(c-'a') + 10
	} else {
		return int(c-'A') + 10
	}
}

// ParseException returns a ParseException with the current line and column
// numbers in the description, suitable for throwing.
//
// Java: TextFormat.Tokenizer.parseException(String)
func (t *TextFormat_Tokenizer) ParseException(description string) *TextFormat_ParseException {
	// Note:  People generally prefer one-based line and column numbers.
	return NewTextFormat_ParseExceptionAt(t.line+1, t.column+1, description)
}

// ParseExceptionPreviousToken returns a ParseException with the line and
// column numbers of the previous token in the description, suitable for
// throwing.
//
// Java: TextFormat.Tokenizer.parseExceptionPreviousToken(String)
func (t *TextFormat_Tokenizer) ParseExceptionPreviousToken(description string) *TextFormat_ParseException {
	// Note:  People generally prefer one-based line and column numbers.
	return NewTextFormat_ParseExceptionAt(t.previousLine+1, t.previousColumn+1, description)
}

// TextFormat_ParseException is thrown when parsing an invalid text format
// message (a java.io.IOException).
type TextFormat_ParseException struct {
	jutil.IOException
}

// NewTextFormat_ParseExceptionAt creates a new instance; line is the line
// number where the parse error occurred, using 1-offset; column the column
// number where the parser error occurred, using 1-offset.
//
// Java: TextFormat.ParseException(int, int, String)
func NewTextFormat_ParseExceptionAt(line int, column int, message string) *TextFormat_ParseException {
	e := &TextFormat_ParseException{}
	jutil.InitThrowable(e, "com.google.javascript.jscomp.jarjar.com.google.protobuf.TextFormat$ParseException",
		jutil.StringPtr(strconv.Itoa(line)+":"+strconv.Itoa(column)+": "+message), nil)
	return e
}
