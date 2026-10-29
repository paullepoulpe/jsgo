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
// text format). This file translates the parsing half of
// com.google.protobuf.TextFormat (the version shaded into the pinned
// closure-compiler.jar): Tokenizer, ParseException,
// InvalidEscapeSequenceException, Parser (merge, the skipping of unknown
// fields) and unescapeBytes / parseInteger, over the generated builders'
// MessageReflection.MergeTarget (proto_descriptors.go). The reconstructed
// messages have no extensions, no google.protobuf.Any and no groups, and no
// caller asks for a parse-info tree, so those branches are not translated;
// an extension name always resolves to "unknown", as with an empty
// ExtensionRegistry.

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

// Java: TextFormat.Tokenizer.getLine()
func (t *TextFormat_Tokenizer) GetLine() int { return t.line }

// Java: TextFormat.Tokenizer.getColumn()
func (t *TextFormat_Tokenizer) GetColumn() int { return t.column }

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
// throw a ParseException.
//
// Java: TextFormat.Tokenizer.consume(String)
func (t *TextFormat_Tokenizer) Consume(token string) {
	if !t.TryConsume(token) {
		panic(t.ParseException("Expected \"" + token + "\"."))
	}
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
// return its value. Otherwise, throw a ParseException.
//
// Java: TextFormat.Tokenizer.consumeIdentifier()
func (t *TextFormat_Tokenizer) ConsumeIdentifier() string {
	if !t.isIdentifier() {
		panic(t.ParseException("Expected identifier. Found '" + t.currentToken + "'"))
	}

	result := t.currentToken
	t.NextToken()
	return result
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
// return its value. Otherwise, throw a ParseException.
//
// Java: TextFormat.Tokenizer.consumeInt32()
func (t *TextFormat_Tokenizer) ConsumeInt32() int32 {
	result, errMsg, ok := textFormat_parseInteger(t.currentToken, true, false)
	if !ok {
		panic(t.IntegerParseException(jutil.NewNumberFormatException(errMsg)))
	}
	t.NextToken()
	return int32(result)
}

// ConsumeUInt32: if the next token is a 32-bit unsigned integer, consume it
// and return its value. Otherwise, throw a ParseException.
//
// Java: TextFormat.Tokenizer.consumeUInt32()
func (t *TextFormat_Tokenizer) ConsumeUInt32() int32 {
	result, errMsg, ok := textFormat_parseInteger(t.currentToken, false, false)
	if !ok {
		panic(t.IntegerParseException(jutil.NewNumberFormatException(errMsg)))
	}
	t.NextToken()
	return int32(result)
}

// ConsumeInt64: if the next token is a 64-bit signed integer, consume it and
// return its value. Otherwise, throw a ParseException.
//
// Java: TextFormat.Tokenizer.consumeInt64()
func (t *TextFormat_Tokenizer) ConsumeInt64() int64 {
	result, errMsg, ok := textFormat_parseInteger(t.currentToken, true, true)
	if !ok {
		panic(t.IntegerParseException(jutil.NewNumberFormatException(errMsg)))
	}
	t.NextToken()
	return result
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

// ConsumeUInt64: if the next token is a 64-bit unsigned integer, consume it
// and return its value. Otherwise, throw a ParseException.
//
// Java: TextFormat.Tokenizer.consumeUInt64()
func (t *TextFormat_Tokenizer) ConsumeUInt64() int64 {
	result, errMsg, ok := textFormat_parseInteger(t.currentToken, false, true)
	if !ok {
		panic(t.IntegerParseException(jutil.NewNumberFormatException(errMsg)))
	}
	t.NextToken()
	return result
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

// ConsumeDouble: if the next token is a double, consume it and return its
// value. Otherwise, throw a ParseException.
//
// Java: TextFormat.Tokenizer.consumeDouble()
func (t *TextFormat_Tokenizer) ConsumeDouble() float64 {
	result, errMsg, ok := t.parseDoubleToken()
	if !ok {
		panic(t.FloatParseException(jutil.NewNumberFormatException(errMsg)))
	}
	t.NextToken()
	return result
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

// ConsumeFloat: if the next token is a float, consume it and return its
// value. Otherwise, throw a ParseException.
//
// Go-only: language-adapter — Float.parseFloat accepts exactly the strings
// that Double.parseDouble accepts; the value is computed by parsing as a
// double and rounding to float32 (no reconstructed message has a float
// field, and the parser only uses this to skip unknown fields).
//
// Java: TextFormat.Tokenizer.consumeFloat()
func (t *TextFormat_Tokenizer) ConsumeFloat() float32 {
	result, errMsg, ok := t.parseFloatToken()
	if !ok {
		panic(t.FloatParseException(jutil.NewNumberFormatException(errMsg)))
	}
	t.NextToken()
	return result
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
// value. Otherwise, throw a ParseException.
//
// Java: TextFormat.Tokenizer.consumeBoolean()
func (t *TextFormat_Tokenizer) ConsumeBoolean() bool {
	if t.currentToken == "true" ||
		t.currentToken == "True" ||
		t.currentToken == "t" ||
		t.currentToken == "1" {
		t.NextToken()
		return true
	} else if t.currentToken == "false" ||
		t.currentToken == "False" ||
		t.currentToken == "f" ||
		t.currentToken == "0" {
		t.NextToken()
		return false
	} else {
		panic(t.ParseException("Expected \"true\" or \"false\". Found \"" + t.currentToken + "\"."))
	}
}

// ConsumeString: if the next token is a string, consume it and return its
// (unescaped) value. Otherwise, throw a ParseException.
//
// Java: TextFormat.Tokenizer.consumeString()
func (t *TextFormat_Tokenizer) ConsumeString() string {
	return jutil.NewStringFromUTF8Bytes(t.ConsumeByteString())
}

// ConsumeByteString: if the next token is a string, consume it, unescape it
// as a ByteString, and return it. Otherwise, throw a ParseException.
//
// Java: TextFormat.Tokenizer.consumeByteString()
func (t *TextFormat_Tokenizer) ConsumeByteString() []byte {
	list := [][]byte{}
	t.consumeByteStringInto(&list)
	for strings.HasPrefix(t.currentToken, "'") || strings.HasPrefix(t.currentToken, "\"") {
		t.consumeByteStringInto(&list)
	}
	var result []byte
	for _, b := range list {
		result = append(result, b...)
	}
	return result
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
func (t *TextFormat_Tokenizer) consumeByteStringInto(list *[][]byte) {
	result, errMsg, ok := t.parseByteStringToken()
	if !ok {
		panic(t.ParseException(errMsg))
	}
	t.NextToken()
	*list = append(*list, result)
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

// IntegerParseException constructs an appropriate ParseException for the
// given NumberFormatException when trying to parse an integer.
//
// Java: TextFormat.Tokenizer.integerParseException(NumberFormatException)
func (t *TextFormat_Tokenizer) IntegerParseException(e jutil.INumberFormatException) *TextFormat_ParseException {
	return t.ParseException("Couldn't parse integer: " + e.GetMessage())
}

// FloatParseException constructs an appropriate ParseException for the given
// NumberFormatException when trying to parse a float or double.
//
// Java: TextFormat.Tokenizer.floatParseException(NumberFormatException)
func (t *TextFormat_Tokenizer) FloatParseException(e jutil.INumberFormatException) *TextFormat_ParseException {
	return t.ParseException("Couldn't parse number: " + e.GetMessage())
}

// TextFormat_ParseException is thrown when parsing an invalid text format
// message (a java.io.IOException).
type TextFormat_ParseException struct {
	jutil.IOException

	line   int
	column int
}

// NewTextFormat_ParseException creates a new instance, with -1 as the line
// and column numbers.
//
// Java: TextFormat.ParseException(String)
func NewTextFormat_ParseException(message string) *TextFormat_ParseException {
	return NewTextFormat_ParseExceptionAt(-1, -1, message)
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
	e.line = line
	e.column = column
	return e
}

// GetLine returns the line where the parse exception occurred, or -1 when
// none is provided. The value is specified as 1-offset, so the first line is
// line 1.
//
// Java: TextFormat.ParseException.getLine()
func (e *TextFormat_ParseException) GetLine() int { return e.line }

// GetColumn returns the column where the parse exception occurred, or -1
// when none is provided. The value is specified as 1-offset, so the first
// line is line 1.
//
// Java: TextFormat.ParseException.getColumn()
func (e *TextFormat_ParseException) GetColumn() int { return e.column }
