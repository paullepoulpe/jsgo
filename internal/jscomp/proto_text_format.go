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
	"math/big"
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

// ConsumeIdentifier: if the next token is an identifier, consume it and
// return its value. Otherwise, throw a ParseException.
//
// Java: TextFormat.Tokenizer.consumeIdentifier()
func (t *TextFormat_Tokenizer) ConsumeIdentifier() string {
	for _, c := range jutil.StringToUTF16(t.currentToken) {
		if TextFormat_Tokenizer_IsAlphaUnder(c) || ('0' <= c && c <= '9') || (c == '.') {
			// OK
		} else {
			panic(t.ParseException("Expected identifier. Found '" + t.currentToken + "'"))
		}
	}

	result := t.currentToken
	t.NextToken()
	return result
}

// textFormat_tryParse runs a consume* call, reporting whether it did not
// throw a ParseException.
//
// Go-only: language-adapter — try { consumeX(); return true; } catch
// (ParseException e) { return false; }
func textFormat_tryParse(consume func()) (ok bool) {
	ok = true
	jutil.Try(consume, func(e ITextFormat_ParseException) { ok = false })
	return ok
}

// TryConsumeIdentifier: if the next token is an identifier, consume it and
// return true. Otherwise, return false without doing anything.
//
// Java: TextFormat.Tokenizer.tryConsumeIdentifier()
func (t *TextFormat_Tokenizer) TryConsumeIdentifier() bool {
	return textFormat_tryParse(func() { t.ConsumeIdentifier() })
}

// ConsumeInt32: if the next token is a 32-bit signed integer, consume it and
// return its value. Otherwise, throw a ParseException.
//
// Java: TextFormat.Tokenizer.consumeInt32()
func (t *TextFormat_Tokenizer) ConsumeInt32() int32 {
	var result int32
	jutil.Try(func() {
		result = TextFormat_ParseInt32(t.currentToken)
	}, func(e jutil.INumberFormatException) {
		panic(t.IntegerParseException(e))
	})
	t.NextToken()
	return result
}

// ConsumeUInt32: if the next token is a 32-bit unsigned integer, consume it
// and return its value. Otherwise, throw a ParseException.
//
// Java: TextFormat.Tokenizer.consumeUInt32()
func (t *TextFormat_Tokenizer) ConsumeUInt32() int32 {
	var result int32
	jutil.Try(func() {
		result = TextFormat_ParseUInt32(t.currentToken)
	}, func(e jutil.INumberFormatException) {
		panic(t.IntegerParseException(e))
	})
	t.NextToken()
	return result
}

// ConsumeInt64: if the next token is a 64-bit signed integer, consume it and
// return its value. Otherwise, throw a ParseException.
//
// Java: TextFormat.Tokenizer.consumeInt64()
func (t *TextFormat_Tokenizer) ConsumeInt64() int64 {
	var result int64
	jutil.Try(func() {
		result = TextFormat_ParseInt64(t.currentToken)
	}, func(e jutil.INumberFormatException) {
		panic(t.IntegerParseException(e))
	})
	t.NextToken()
	return result
}

// TryConsumeInt64: if the next token is a 64-bit signed integer, consume it
// and return true. Otherwise, return false without doing anything.
//
// Java: TextFormat.Tokenizer.tryConsumeInt64()
func (t *TextFormat_Tokenizer) TryConsumeInt64() bool {
	return textFormat_tryParse(func() { t.ConsumeInt64() })
}

// ConsumeUInt64: if the next token is a 64-bit unsigned integer, consume it
// and return its value. Otherwise, throw a ParseException.
//
// Java: TextFormat.Tokenizer.consumeUInt64()
func (t *TextFormat_Tokenizer) ConsumeUInt64() int64 {
	var result int64
	jutil.Try(func() {
		result = TextFormat_ParseUInt64(t.currentToken)
	}, func(e jutil.INumberFormatException) {
		panic(t.IntegerParseException(e))
	})
	t.NextToken()
	return result
}

// TryConsumeUInt64: if the next token is a 64-bit unsigned integer, consume
// it and return true. Otherwise, return false without doing anything.
//
// Java: TextFormat.Tokenizer.tryConsumeUInt64()
func (t *TextFormat_Tokenizer) TryConsumeUInt64() bool {
	return textFormat_tryParse(func() { t.ConsumeUInt64() })
}

// ConsumeDouble: if the next token is a double, consume it and return its
// value. Otherwise, throw a ParseException.
//
// Java: TextFormat.Tokenizer.consumeDouble()
func (t *TextFormat_Tokenizer) ConsumeDouble() float64 {
	// We need to parse infinity and nan separately because
	// Double.parseDouble() does not accept "inf", "infinity", or "nan".
	switch jutil.StringToLowerCaseRoot(t.currentToken) {
	case "-inf", "-infinity":
		t.NextToken()
		return jutil.DoubleParseDouble("-Infinity")
	case "inf", "infinity":
		t.NextToken()
		return jutil.DoubleParseDouble("Infinity")
	case "nan":
		t.NextToken()
		return jutil.DoubleParseDouble("NaN")
	default:
		// fall through
	}

	var result float64
	jutil.Try(func() {
		result = jutil.DoubleParseDouble(t.currentToken)
	}, func(e jutil.INumberFormatException) {
		panic(t.FloatParseException(e))
	})
	t.NextToken()
	return result
}

// TryConsumeDouble: if the next token is a double, consume it and return
// true. Otherwise, return false without doing anything.
//
// Java: TextFormat.Tokenizer.tryConsumeDouble()
func (t *TextFormat_Tokenizer) TryConsumeDouble() bool {
	return textFormat_tryParse(func() { t.ConsumeDouble() })
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
	// We need to parse infinity and nan separately because
	// Float.parseFloat() does not accept "inf", "infinity", or "nan".
	switch jutil.StringToLowerCaseRoot(t.currentToken) {
	case "-inf", "-inff", "-infinity", "-infinityf":
		t.NextToken()
		return float32(jutil.DoubleParseDouble("-Infinity"))
	case "inf", "inff", "infinity", "infinityf":
		t.NextToken()
		return float32(jutil.DoubleParseDouble("Infinity"))
	case "nan", "nanf":
		t.NextToken()
		return float32(jutil.DoubleParseDouble("NaN"))
	default:
		// fall through
	}

	var result float32
	jutil.Try(func() {
		result = float32(jutil.DoubleParseDouble(t.currentToken))
	}, func(e jutil.INumberFormatException) {
		panic(t.FloatParseException(e))
	})
	t.NextToken()
	return result
}

// TryConsumeFloat: if the next token is a float, consume it and return true.
// Otherwise, return false without doing anything.
//
// Java: TextFormat.Tokenizer.tryConsumeFloat()
func (t *TextFormat_Tokenizer) TryConsumeFloat() bool {
	return textFormat_tryParse(func() { t.ConsumeFloat() })
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
	return textFormat_tryParse(func() { t.ConsumeByteString() })
}

// consumeByteStringInto is like ConsumeByteString but adds each token of the
// string to the given list. String literals (whether bytes or text) may come
// in multiple adjacent tokens which are automatically concatenated, like in
// C or Python.
//
// Java: TextFormat.Tokenizer.consumeByteString(List<ByteString>)
func (t *TextFormat_Tokenizer) consumeByteStringInto(list *[][]byte) {
	var quote byte
	if len(t.currentToken) > 0 {
		quote = t.currentToken[0]
	}
	if quote != '"' && quote != '\'' {
		panic(t.ParseException("Expected string."))
	}

	if len(t.currentToken) < 2 || t.currentToken[len(t.currentToken)-1] != quote {
		panic(t.ParseException("String missing ending quote."))
	}

	jutil.Try(func() {
		escaped := t.currentToken[1 : len(t.currentToken)-1]
		result := TextFormat_UnescapeBytes(escaped)
		t.NextToken()
		*list = append(*list, result)
	}, func(e *TextFormat_InvalidEscapeSequenceException) {
		panic(t.ParseException(e.GetMessage()))
	})
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

func (*TextFormat_ParseException) textFormatParseException() {}

// ITextFormat_ParseException matches TextFormat.ParseException.
type ITextFormat_ParseException interface {
	jutil.IIOException
	textFormatParseException()
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

// Java: TextFormat.PARSER
var textFormat_PARSER = TextFormat_Parser_NewBuilder().Build()

// TextFormat_GetParser returns a Parser instance which can parse text-format
// messages. The returned instance is thread-safe.
//
// Java: TextFormat.getParser()
func TextFormat_GetParser() *TextFormat_Parser {
	return textFormat_PARSER
}

// TextFormat_Merge parses a text-format message from input and merges the
// contents into builder. It panics with TextFormat_ParseException (a checked
// exception).
//
// Java: TextFormat.merge(CharSequence, Message.Builder)
func TextFormat_Merge(input string, builder TextFormat_MergeTarget) {
	textFormat_PARSER.Merge(input, builder)
}

// TextFormat_Parser is the parser for text-format proto2 instances. This
// class is thread-safe. The implementation largely follows
// google/protobuf/text_format.cc.
//
// Use TextFormat#getParser() to obtain the default parser, or Builder to
// control the parser behavior.
type TextFormat_Parser struct {
	allowUnknownFields      bool
	allowUnknownEnumValues  bool
	allowUnknownExtensions  bool
	singularOverwritePolicy TextFormat_Parser_SingularOverwritePolicy
	recursionLimit          int
}

// TextFormat_Parser_SingularOverwritePolicy determines if repeated values for
// non-repeated fields and oneofs are permitted. For example, given
// required/optional field "foo" and a oneof containing "baz" and "moo":
//
//   - "foo: 1 foo: 2"
//   - "baz: 1 moo: 2"
//   - merging "foo: 2" into a proto in which foo is already set, or
//   - merging "moo: 2" into a proto in which baz is already set.
type TextFormat_Parser_SingularOverwritePolicy int32

const (
	// Later values are merged with earlier values. For primitive fields or
	// conflicting oneofs, the last value is retained.
	TextFormat_Parser_SingularOverwritePolicy_ALLOW_SINGULAR_OVERWRITES TextFormat_Parser_SingularOverwritePolicy = iota
	// An error is issued.
	TextFormat_Parser_SingularOverwritePolicy_FORBID_SINGULAR_OVERWRITES
)

// Java: TextFormat.Parser.newBuilder()
func TextFormat_Parser_NewBuilder() *TextFormat_Parser_Builder {
	return &TextFormat_Parser_Builder{
		singularOverwritePolicy: TextFormat_Parser_SingularOverwritePolicy_ALLOW_SINGULAR_OVERWRITES,
		recursionLimit:          100,
	}
}

// TextFormat_Parser_Builder can be used to obtain new instances of Parser.
type TextFormat_Parser_Builder struct {
	allowUnknownFields      bool
	allowUnknownEnumValues  bool
	allowUnknownExtensions  bool
	singularOverwritePolicy TextFormat_Parser_SingularOverwritePolicy
	recursionLimit          int
}

// SetAllowUnknownFields sets whether this parser will allow unknown fields.
// By default, an exception is thrown if an unknown field is encountered. If
// this is set, the parser will only log a warning. Allow unknown fields will
// also allow unknown extensions.
//
// Java: TextFormat.Parser.Builder.setAllowUnknownFields(boolean)
func (b *TextFormat_Parser_Builder) SetAllowUnknownFields(allowUnknownFields bool) *TextFormat_Parser_Builder {
	b.allowUnknownFields = allowUnknownFields
	return b
}

// SetAllowUnknownExtensions sets whether this parser will allow unknown
// extensions. By default, an exception is thrown if unknown extension is
// encountered. If this is set true, the parser will only log a warning.
// Allow unknown extensions does not mean allow normal unknown fields.
//
// Java: TextFormat.Parser.Builder.setAllowUnknownExtensions(boolean)
func (b *TextFormat_Parser_Builder) SetAllowUnknownExtensions(allowUnknownExtensions bool) *TextFormat_Parser_Builder {
	b.allowUnknownExtensions = allowUnknownExtensions
	return b
}

// SetSingularOverwritePolicy sets parser behavior when a non-repeated field
// appears more than once.
//
// Java: TextFormat.Parser.Builder.setSingularOverwritePolicy(SingularOverwritePolicy)
func (b *TextFormat_Parser_Builder) SetSingularOverwritePolicy(p TextFormat_Parser_SingularOverwritePolicy) *TextFormat_Parser_Builder {
	b.singularOverwritePolicy = p
	return b
}

// SetRecursionLimit sets the maximum recursion limit that the parser will
// allow. If the depth of the message exceeds this limit then the parser will
// stop and throw an exception.
//
// Java: TextFormat.Parser.Builder.setRecursionLimit(int)
func (b *TextFormat_Parser_Builder) SetRecursionLimit(recursionLimit int) *TextFormat_Parser_Builder {
	b.recursionLimit = recursionLimit
	return b
}

// Java: TextFormat.Parser.Builder.build()
func (b *TextFormat_Parser_Builder) Build() *TextFormat_Parser {
	return &TextFormat_Parser{
		allowUnknownFields:      b.allowUnknownFields,
		allowUnknownEnumValues:  b.allowUnknownEnumValues,
		allowUnknownExtensions:  b.allowUnknownExtensions,
		singularOverwritePolicy: b.singularOverwritePolicy,
		recursionLimit:          b.recursionLimit,
	}
}

// TextFormat_Parser_UnknownField_Type is TextFormat.Parser.UnknownField.Type.
type TextFormat_Parser_UnknownField_Type int32

const (
	TextFormat_Parser_UnknownField_Type_FIELD TextFormat_Parser_UnknownField_Type = iota
	TextFormat_Parser_UnknownField_Type_EXTENSION
)

// Java: TextFormat.Parser.UnknownField
type TextFormat_Parser_UnknownField struct {
	message string
	type_   TextFormat_Parser_UnknownField_Type
}

// CheckUnknownFields checks both unknown fields and unknown extensions and
// logs warning messages or throws exceptions according to the flag.
//
// Go-only: language-adapter — the java.util.logging warnings of the
// allow-unknown modes are not emitted (no caller enables them).
//
// Java: TextFormat.Parser.checkUnknownFields(List<UnknownField>)
func (p *TextFormat_Parser) CheckUnknownFields(unknownFields []*TextFormat_Parser_UnknownField) {
	if len(unknownFields) == 0 {
		return
	}

	var msg strings.Builder
	msg.WriteString("Input contains unknown fields and/or extensions:")
	for _, field := range unknownFields {
		msg.WriteByte('\n')
		msg.WriteString(field.message)
	}

	if p.allowUnknownFields {
		return
	}

	firstErrorIndex := 0
	if p.allowUnknownExtensions {
		allUnknownExtensions := true
		for _, field := range unknownFields {
			if field.type_ == TextFormat_Parser_UnknownField_Type_FIELD {
				allUnknownExtensions = false
				break
			}
			firstErrorIndex++
		}
		if allUnknownExtensions {
			return
		}
	}

	lineColumn := strings.Split(unknownFields[firstErrorIndex].message, ":")
	line, _ := strconv.Atoi(lineColumn[0])
	column, _ := strconv.Atoi(lineColumn[1])
	panic(NewTextFormat_ParseExceptionAt(line, column, msg.String()))
}

// Merge parses a text-format message from input and merges the contents
// into builder.
//
// Java: TextFormat.Parser.merge(CharSequence, ExtensionRegistry, Message.Builder)
func (p *TextFormat_Parser) Merge(input string, target TextFormat_MergeTarget) {
	tokenizer := NewTextFormat_Tokenizer(input)
	unknownFields := []*TextFormat_Parser_UnknownField{}

	for !tokenizer.AtEnd() {
		p.MergeField(tokenizer, target, &unknownFields, p.recursionLimit)
	}
	p.CheckUnknownFields(unknownFields)
}

// MergeField parses a single field from tokenizer and merges it into target.
//
// Java: TextFormat.Parser.mergeField(Tokenizer, ExtensionRegistry, MergeTarget, TextFormatParseInfoTree.Builder, List<UnknownField>, int)
func (p *TextFormat_Parser) MergeField(
	tokenizer *TextFormat_Tokenizer,
	target TextFormat_MergeTarget,
	unknownFields *[]*TextFormat_Parser_UnknownField,
	recursionLimit int) {
	var field *Descriptors_FieldDescriptor
	var name string
	type_ := target.GetDescriptorForType()

	// (A google.protobuf.Any type URL "[...]" needs the Any descriptor; no
	// reconstructed message is an Any.)

	if tokenizer.TryConsume("[") {
		// An extension.
		var nameBuilder strings.Builder
		nameBuilder.WriteString(tokenizer.ConsumeIdentifier())
		for tokenizer.TryConsume(".") {
			nameBuilder.WriteByte('.')
			nameBuilder.WriteString(tokenizer.ConsumeIdentifier())
		}
		name = nameBuilder.String()

		// target.findExtensionByName: the reconstructed messages declare no
		// extensions (an empty ExtensionRegistry).
		message := strconv.Itoa(tokenizer.GetPreviousLine()+1) +
			":" +
			strconv.Itoa(tokenizer.GetPreviousColumn()+1) +
			":\t" +
			type_.GetFullName() +
			".[" +
			name +
			"]"
		*unknownFields = append(*unknownFields, &TextFormat_Parser_UnknownField{message, TextFormat_Parser_UnknownField_Type_EXTENSION})

		tokenizer.Consume("]")
	} else {
		name = tokenizer.ConsumeIdentifier()
		field = type_.FindFieldByName(name)

		// Group names are expected to be capitalized as they appear in the
		// .proto file, which actually matches their type names, not their field
		// names.
		if field == nil {
			// Explicitly specify US locale so that this code does not break when
			// executing in Turkey.
			lowerName := jutil.StringToLowerCaseRoot(name)
			field = type_.FindFieldByName(lowerName)
			// If the case-insensitive match worked but the field is NOT a group,
			if field != nil && !field.IsGroupLike() {
				field = nil
			}
			if field != nil && field.GetMessageType().GetName() != name {
				field = nil
			}
		}

		if field == nil {
			message := strconv.Itoa(tokenizer.GetPreviousLine()+1) +
				":" +
				strconv.Itoa(tokenizer.GetPreviousColumn()+1) +
				":\t" +
				type_.GetFullName() +
				"." +
				name
			*unknownFields = append(*unknownFields, &TextFormat_Parser_UnknownField{message, TextFormat_Parser_UnknownField_Type_FIELD})
		}
	}

	// Skips unknown fields.
	if field == nil {
		p.GuessFieldTypeAndSkip(tokenizer, type_, recursionLimit)
		return
	}

	// Handle potential ':'.
	if field.GetJavaType() == Descriptors_FieldDescriptor_JavaType_MESSAGE {
		tokenizer.TryConsume(":") // optional
		p.ConsumeFieldValues(tokenizer, target, field, unknownFields, recursionLimit)
	} else {
		tokenizer.Consume(":") // required
		p.ConsumeFieldValues(tokenizer, target, field, unknownFields, recursionLimit)
	}

	// For historical reasons, fields may optionally be separated by commas or
	// semicolons.
	if !tokenizer.TryConsume(";") {
		tokenizer.TryConsume(",")
	}
}

// Java: TextFormat.Parser.consumeFullTypeName(Tokenizer)
func (p *TextFormat_Parser) ConsumeFullTypeName(tokenizer *TextFormat_Tokenizer) {
	// If there is not a leading `[`, this is just a type name.
	if !tokenizer.TryConsume("[") {
		tokenizer.ConsumeIdentifier()
		return
	}

	// Otherwise, this is an extension or google.protobuf.Any type URL: we consume proto path
	// elements until we've addressed the type.
	tokenizer.ConsumeIdentifier()
	for tokenizer.TryConsume(".") {
		tokenizer.ConsumeIdentifier()
	}
	if tokenizer.TryConsume("/") {
		tokenizer.ConsumeIdentifier()
		for tokenizer.TryConsume(".") {
			tokenizer.ConsumeIdentifier()
		}
	}
	tokenizer.Consume("]")
}

// ConsumeFieldValues parses one or more field values from tokenizer and
// merges them into target.
//
// Java: TextFormat.Parser.consumeFieldValues(Tokenizer, ExtensionRegistry, MergeTarget, FieldDescriptor, ExtensionInfo, TextFormatParseInfoTree.Builder, List<UnknownField>, int)
func (p *TextFormat_Parser) ConsumeFieldValues(
	tokenizer *TextFormat_Tokenizer,
	target TextFormat_MergeTarget,
	field *Descriptors_FieldDescriptor,
	unknownFields *[]*TextFormat_Parser_UnknownField,
	recursionLimit int) {
	// Support specifying repeated field values as a comma-separated list.
	// Ex."foo: [1, 2, 3]"
	if field.IsRepeated() && tokenizer.TryConsume("[") {
		if !tokenizer.TryConsume("]") { // Allow "foo: []" to be treated as empty.
			for {
				p.ConsumeFieldValue(tokenizer, target, field, unknownFields, recursionLimit)
				if tokenizer.TryConsume("]") {
					// End of list.
					break
				}
				tokenizer.Consume(",")
			}
		}
	} else {
		p.ConsumeFieldValue(tokenizer, target, field, unknownFields, recursionLimit)
	}
}

// ConsumeFieldValue parses a single field value from tokenizer and merges it
// into target.
//
// Java: TextFormat.Parser.consumeFieldValue(Tokenizer, ExtensionRegistry, MergeTarget, FieldDescriptor, ExtensionInfo, TextFormatParseInfoTree.Builder, List<UnknownField>, int)
func (p *TextFormat_Parser) ConsumeFieldValue(
	tokenizer *TextFormat_Tokenizer,
	target TextFormat_MergeTarget,
	field *Descriptors_FieldDescriptor,
	unknownFields *[]*TextFormat_Parser_UnknownField,
	recursionLimit int) {
	if p.singularOverwritePolicy == TextFormat_Parser_SingularOverwritePolicy_FORBID_SINGULAR_OVERWRITES &&
		!field.IsRepeated() {
		if target.protoHasField(field) {
			panic(tokenizer.ParseExceptionPreviousToken(
				"Non-repeated field \"" + field.GetFullName() + "\" cannot be overwritten."))
		}
		// (The reconstructed messages Closure parses have no oneofs.)
	}

	var value any

	if field.GetJavaType() == Descriptors_FieldDescriptor_JavaType_MESSAGE {
		if recursionLimit < 1 {
			panic(tokenizer.ParseException("Message is nested too deep"))
		}

		var endToken string
		if tokenizer.TryConsume("<") {
			endToken = ">"
		} else {
			tokenizer.Consume("{")
			endToken = "}"
		}

		subField := target.protoNewMergeTargetForField(field)

		for !tokenizer.TryConsume(endToken) {
			if tokenizer.AtEnd() {
				panic(tokenizer.ParseException("Expected \"" + endToken + "\"."))
			}
			p.MergeField(tokenizer, subField, unknownFields, recursionLimit-1)
		}

		value = subField.protoFinish()
	} else {
		switch field.GetType() {
		case Descriptors_FieldDescriptor_Type_INT32,
			Descriptors_FieldDescriptor_Type_SINT32,
			Descriptors_FieldDescriptor_Type_SFIXED32:
			value = tokenizer.ConsumeInt32()

		case Descriptors_FieldDescriptor_Type_INT64,
			Descriptors_FieldDescriptor_Type_SINT64,
			Descriptors_FieldDescriptor_Type_SFIXED64:
			value = tokenizer.ConsumeInt64()

		case Descriptors_FieldDescriptor_Type_UINT32,
			Descriptors_FieldDescriptor_Type_FIXED32:
			value = tokenizer.ConsumeUInt32()

		case Descriptors_FieldDescriptor_Type_UINT64,
			Descriptors_FieldDescriptor_Type_FIXED64:
			value = tokenizer.ConsumeUInt64()

		case Descriptors_FieldDescriptor_Type_FLOAT:
			value = tokenizer.ConsumeFloat()

		case Descriptors_FieldDescriptor_Type_DOUBLE:
			value = tokenizer.ConsumeDouble()

		case Descriptors_FieldDescriptor_Type_BOOL:
			value = tokenizer.ConsumeBoolean()

		case Descriptors_FieldDescriptor_Type_STRING:
			value = tokenizer.ConsumeString()

		case Descriptors_FieldDescriptor_Type_BYTES:
			value = tokenizer.ConsumeByteString()

		case Descriptors_FieldDescriptor_Type_ENUM:
			enumType := field.GetEnumType()

			if tokenizer.LookingAtInteger() {
				number := tokenizer.ConsumeInt32()
				// (Every reconstructed enum that a parsed message uses is closed.)
				var enumValue *Descriptors_EnumValueDescriptor
				if enumType.IsClosed() {
					enumValue = enumType.FindValueByNumber(number)
				} else {
					panic(jutil.NewUnsupportedOperationExceptionMsg("open enum " + enumType.GetFullName()))
				}
				if enumValue == nil {
					unknownValueMsg := "Enum type \"" +
						enumType.GetFullName() +
						"\" has no value with number " +
						strconv.Itoa(int(number)) +
						"."
					if p.allowUnknownEnumValues {
						return
					} else {
						panic(tokenizer.ParseExceptionPreviousToken(unknownValueMsg))
					}
				}
				value = enumValue
			} else {
				id := tokenizer.ConsumeIdentifier()
				enumValue := enumType.FindValueByName(id)
				if enumValue == nil {
					unknownValueMsg := "Enum type \"" +
						enumType.GetFullName() +
						"\" has no value named \"" +
						id +
						"\"."
					if p.allowUnknownEnumValues {
						return
					} else {
						panic(tokenizer.ParseExceptionPreviousToken(unknownValueMsg))
					}
				}
				value = enumValue
			}

		case Descriptors_FieldDescriptor_Type_MESSAGE, Descriptors_FieldDescriptor_Type_GROUP:
			panic(jutil.NewRuntimeException("Can't get here."))
		}
	}

	if field.IsRepeated() {
		// TODO: If field.isMapField() and FORBID_SINGULAR_OVERWRITES mode,
		//     check for duplicate map keys here.
		target.protoAddRepeatedField(field, value)
	} else {
		target.protoSetField(field, value)
	}
}

// SkipField skips the next field including the field's name and value.
//
// Java: TextFormat.Parser.skipField(Tokenizer, Descriptor, int)
func (p *TextFormat_Parser) SkipField(tokenizer *TextFormat_Tokenizer, type_ *Descriptors_Descriptor, recursionLimit int) {
	p.ConsumeFullTypeName(tokenizer)
	p.GuessFieldTypeAndSkip(tokenizer, type_, recursionLimit)

	if !tokenizer.TryConsume(";") {
		tokenizer.TryConsume(",")
	}
}

// SkipFieldMessage skips the whole body of a message including the beginning
// delimiter and the ending delimiter.
//
// Java: TextFormat.Parser.skipFieldMessage(Tokenizer, Descriptor, int)
func (p *TextFormat_Parser) SkipFieldMessage(tokenizer *TextFormat_Tokenizer, type_ *Descriptors_Descriptor, recursionLimit int) {
	var delimiter string
	if tokenizer.TryConsume("<") {
		delimiter = ">"
	} else {
		tokenizer.Consume("{")
		delimiter = "}"
	}
	for !tokenizer.LookingAt(">") && !tokenizer.LookingAt("}") {
		p.SkipField(tokenizer, type_, recursionLimit)
	}
	tokenizer.Consume(delimiter)
}

// SkipFieldValue skips a field value.
//
// Java: TextFormat.Parser.skipFieldValue(Tokenizer)
func (p *TextFormat_Parser) SkipFieldValue(tokenizer *TextFormat_Tokenizer) {
	if !tokenizer.TryConsumeByteString() &&
		!tokenizer.TryConsumeIdentifier() && // includes enum & boolean
		!tokenizer.TryConsumeInt64() && // includes int32
		!tokenizer.TryConsumeUInt64() && // includes uint32
		!tokenizer.TryConsumeDouble() &&
		!tokenizer.TryConsumeFloat() {
		panic(tokenizer.ParseException("Invalid field value: " + tokenizer.currentToken))
	}
}

// GuessFieldTypeAndSkip tries to guess the type of this field and skips it.
//
// Java: TextFormat.Parser.guessFieldTypeAndSkip(Tokenizer, Descriptor, int)
func (p *TextFormat_Parser) GuessFieldTypeAndSkip(tokenizer *TextFormat_Tokenizer, type_ *Descriptors_Descriptor, recursionLimit int) {
	semicolonConsumed := tokenizer.TryConsume(":")
	if tokenizer.LookingAt("[") {
		p.SkipFieldShortFormedRepeated(tokenizer, semicolonConsumed, type_, recursionLimit)
	} else if semicolonConsumed && !tokenizer.LookingAt("{") && !tokenizer.LookingAt("<") {
		p.SkipFieldValue(tokenizer)
	} else {
		if recursionLimit < 1 {
			panic(tokenizer.ParseException("Message is nested too deep"))
		}
		p.SkipFieldMessage(tokenizer, type_, recursionLimit-1)
	}
}

// SkipFieldShortFormedRepeated skips a short-formed repeated field value.
//
// Reports an error if scalar type is not allowed but showing up inside "[]".
//
// Java: TextFormat.Parser.skipFieldShortFormedRepeated(Tokenizer, boolean, Descriptor, int)
func (p *TextFormat_Parser) SkipFieldShortFormedRepeated(
	tokenizer *TextFormat_Tokenizer, scalarAllowed bool, type_ *Descriptors_Descriptor, recursionLimit int) {
	if !tokenizer.TryConsume("[") || tokenizer.TryConsume("]") {
		return
	}

	for {
		if tokenizer.LookingAt("{") || tokenizer.LookingAt("<") {
			if recursionLimit < 1 {
				panic(tokenizer.ParseException("Message is nested too deep"))
			}
			p.SkipFieldMessage(tokenizer, type_, recursionLimit-1)
		} else if scalarAllowed {
			p.SkipFieldValue(tokenizer)
		} else {
			panic(tokenizer.ParseException(
				"Invalid repeated scalar field: missing \":\" before \"[\"."))
		}
		if tokenizer.TryConsume("]") {
			break
		}
		tokenizer.Consume(",")
	}
}

// TextFormat_UnescapeBytes un-escapes a byte sequence as escaped using
// escapeBytes(ByteString). Two-digit hex escapes (starting with "\x") are
// also recognized. It panics with InvalidEscapeSequenceException.
//
// Java: TextFormat.unescapeBytes(CharSequence)
func TextFormat_UnescapeBytes(charString string) []byte {
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
							panic(NewTextFormat_InvalidEscapeSequenceException(
								"Invalid escape sequence: '\\x' with no digits"))
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
								panic(NewTextFormat_InvalidEscapeSequenceException(
									"Invalid escape sequence: '\\u' refers to a surrogate"))
							}
							result = append(result, jutil.StringGetBytesUTF8(jutil.StringFromUTF16([]uint16{ch}))...)
							i += 3
						} else {
							panic(NewTextFormat_InvalidEscapeSequenceException(
								"Invalid escape sequence: '\\u' with too few hex chars"))
						}

					case 'U':
						i++
						if i+7 >= len(input) {
							panic(NewTextFormat_InvalidEscapeSequenceException(
								"Invalid escape sequence: '\\U' with too few hex chars"))
						}
						codepoint := 0
						for offset := i; offset < i+8; offset++ {
							b := input[offset]
							if !textFormat_isHex(b) {
								panic(NewTextFormat_InvalidEscapeSequenceException(
									"Invalid escape sequence: '\\U' with too few hex chars"))
							}
							codepoint = (codepoint << 4) | textFormat_digitValue(b)
						}
						// Character.isValidCodePoint; the hex digits make it non-negative
						// as an int only below 0x80000000.
						if int32(codepoint) < 0 || codepoint > 0x10FFFF {
							panic(NewTextFormat_InvalidEscapeSequenceException(
								"Invalid escape sequence: '\\U" +
									jutil.NewStringFromUTF8Bytes(input[i:i+8]) +
									"' is not a valid code point value"))
						}
						// Character.UnicodeBlock LOW_SURROGATES, HIGH_SURROGATES,
						// HIGH_PRIVATE_USE_SURROGATES: U+D800..U+DFFF.
						if codepoint >= 0xD800 && codepoint <= 0xDFFF {
							panic(NewTextFormat_InvalidEscapeSequenceException(
								"Invalid escape sequence: '\\U" +
									jutil.NewStringFromUTF8Bytes(input[i:i+8]) +
									"' refers to a surrogate code unit"))
						}
						result = append(result, []byte(string(rune(codepoint)))...)
						i += 7

					default:
						panic(NewTextFormat_InvalidEscapeSequenceException(
							"Invalid escape sequence: '\\" + string(rune(c)) + "'"))
					}
				}
			} else {
				panic(NewTextFormat_InvalidEscapeSequenceException(
					"Invalid escape sequence: '\\' at end of string."))
			}
		} else {
			result = append(result, c)
		}
	}

	return result
}

// TextFormat_InvalidEscapeSequenceException is thrown by unescape methods
// upon encountering an invalid escape sequence (an IOException).
type TextFormat_InvalidEscapeSequenceException struct {
	jutil.IOException
}

func (*TextFormat_InvalidEscapeSequenceException) invalidEscapeSequenceException() {}

// Java: TextFormat.InvalidEscapeSequenceException(String)
func NewTextFormat_InvalidEscapeSequenceException(description string) *TextFormat_InvalidEscapeSequenceException {
	e := &TextFormat_InvalidEscapeSequenceException{}
	jutil.InitThrowable(e, "com.google.javascript.jscomp.jarjar.com.google.protobuf.TextFormat$InvalidEscapeSequenceException",
		jutil.StringPtr(description), nil)
	return e
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

// TextFormat_ParseInt32 parses a 32-bit signed integer from the text. Unlike
// the Java standard Integer.parseInt(), this function recognizes the prefixes
// "0x" and "0" to signify hexadecimal and octal numbers, respectively.
//
// Java: TextFormat.parseInt32(String)
func TextFormat_ParseInt32(text string) int32 {
	return int32(TextFormat_ParseInteger(text, true, false))
}

// TextFormat_ParseUInt32 parses a 32-bit unsigned integer from the text.
// Unlike the Java standard Integer.parseInt(), this function recognizes the
// prefixes "0x" and "0" to signify hexadecimal and octal numbers,
// respectively. The result is coerced to a (signed) int when returned since
// Java has no unsigned integer type.
//
// Java: TextFormat.parseUInt32(String)
func TextFormat_ParseUInt32(text string) int32 {
	return int32(TextFormat_ParseInteger(text, false, false))
}

// TextFormat_ParseInt64 parses a 64-bit signed integer from the text. Unlike
// the Java standard Integer.parseInt(), this function recognizes the prefixes
// "0x" and "0" to signify hexadecimal and octal numbers, respectively.
//
// Java: TextFormat.parseInt64(String)
func TextFormat_ParseInt64(text string) int64 {
	return TextFormat_ParseInteger(text, true, true)
}

// TextFormat_ParseUInt64 parses a 64-bit unsigned integer from the text.
// Unlike the Java standard Integer.parseInt(), this function recognizes the
// prefixes "0x" and "0" to signify hexadecimal and octal numbers,
// respectively. The result is coerced to a (signed) long when returned since
// Java has no unsigned long type.
//
// Java: TextFormat.parseUInt64(String)
func TextFormat_ParseUInt64(text string) int64 {
	return TextFormat_ParseInteger(text, false, true)
}

// Java: TextFormat.parseInteger(String, boolean, boolean)
func TextFormat_ParseInteger(text string, isSigned bool, isLong bool) int64 {
	pos := 0

	negative := false
	if strings.HasPrefix(text[pos:], "-") {
		if !isSigned {
			panic(jutil.NewNumberFormatException("Number must be positive: " + text))
		}
		pos++
		negative = true
	}

	radix := 10
	if strings.HasPrefix(text[pos:], "0x") {
		pos += 2
		radix = 16
	} else if strings.HasPrefix(text[pos:], "0") {
		radix = 8
	}

	numberText := jutil.StringSubstringFrom(text, pos)

	var result int64
	if jutil.StringLength(numberText) < 16 {
		result = jutil.LongParseLong(numberText, radix)
		if negative {
			result = -result
		}

		if !isLong {
			if isSigned {
				if result > 2147483647 || result < -2147483648 {
					panic(jutil.NewNumberFormatException(
						"Number out of range for 32-bit signed integer: " + text))
				}
			} else {
				if result >= (int64(1)<<32) || result < 0 {
					panic(jutil.NewNumberFormatException(
						"Number out of range for 32-bit unsigned integer: " + text))
				}
			}
		}
	} else {
		bigValue := jutil.NewBigIntegerFromStringRadix(numberText, radix)
		if negative {
			bigValue = new(big.Int).Neg(bigValue)
		}

		// BigInteger.bitLength(): the minimal two's-complement width
		// excluding the sign bit.
		bitLength := bigValue.BitLen()
		if bigValue.Sign() < 0 {
			bitLength = new(big.Int).Not(bigValue).BitLen()
		}
		if !isLong {
			if isSigned {
				if bitLength > 31 {
					panic(jutil.NewNumberFormatException(
						"Number out of range for 32-bit signed integer: " + text))
				}
			} else {
				if bitLength > 32 {
					panic(jutil.NewNumberFormatException(
						"Number out of range for 32-bit unsigned integer: " + text))
				}
			}
		} else {
			if isSigned {
				if bitLength > 63 {
					panic(jutil.NewNumberFormatException(
						"Number out of range for 64-bit signed integer: " + text))
				}
			} else {
				if bitLength > 64 {
					panic(jutil.NewNumberFormatException(
						"Number out of range for 64-bit unsigned integer: " + text))
				}
			}
		}

		// BigInteger.longValue(): the low-order 64 bits.
		result = int64(new(big.Int).And(bigValue, new(big.Int).SetUint64(^uint64(0))).Uint64())
	}

	return result
}
