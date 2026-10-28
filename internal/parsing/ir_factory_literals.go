// Port of com/google/javascript/jscomp/parsing/IRFactory.java @ closure-compiler v20261021

/*
 * Copyright 2013 The Closure Compiler Authors.
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

package parsing

import (
	"math/big"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

// irFactory_substring is String.substring(begin, end) on a string given as
// its UTF-16 code units.
//
// Go-only: language-adapter — Java String API on UTF-16 code units.
func irFactory_substring(units []uint16, begin, end int) string {
	if begin < 0 || end > len(units) || begin > end {
		panic(jutil.NewStringIndexOutOfBoundsException(begin, end))
	}
	return jutil.StringFromUTF16(units[begin:end])
}

// irFactory_charAt is String.charAt(int) on a string given as its UTF-16
// code units.
//
// Go-only: language-adapter — Java String API on UTF-16 code units.
func irFactory_charAt(units []uint16, index int) uint16 {
	if index < 0 || index >= len(units) {
		panic(jutil.NewStringIndexOutOfBoundsException(index, len(units)))
	}
	return units[index]
}

// irFactory_indexOf is String.indexOf(int ch, int fromIndex) on a string
// given as its UTF-16 code units.
//
// Go-only: language-adapter — Java String API on UTF-16 code units.
func irFactory_indexOf(units []uint16, ch uint16, fromIndex int) int {
	if fromIndex < 0 {
		fromIndex = 0
	}
	for i := fromIndex; i < len(units); i++ {
		if units[i] == ch {
			return i
		}
	}
	return -1
}

func (f *IRFactory) rawTokenValue(token parser.Token) string {
	checkTokenNotNull(token)
	if id := token.MaybePrivateValueID(); id != 0 {
		return strpool.String(id)
	}
	start, end := token.ValueRange()
	if s, ok := f.fileWithContent.ASCIISlice(start, end); ok {
		return s
	}
	return f.fileWithContent.Substring(start, end)
}

// Java: IRFactory.normalizeRegex(LiteralToken)
func (f *IRFactory) NormalizeRegex(token parser.Token) string {
	value := f.rawTokenValue(token)
	if strings.IndexByte(value, '\\') == -1 {
		// Go-only: language-adapter — the delimiting '/' are ASCII, so their
		// byte indices delimit the same substring as Java's UTF-16 indices.
		lastSlash := strings.LastIndexByte(value, '/')
		if lastSlash < 1 {
			panic(jutil.NewStringIndexOutOfBoundsException(1, lastSlash))
		}
		return value[1:lastSlash]
	}
	// Go-only: language-adapter — the String and StringBuilder operations of
	// the escape path work on the UTF-16 code units of the value.
	units := jutil.StringToUTF16(value)
	lastSlash := -1
	for i := len(units) - 1; i >= 0; i-- {
		if units[i] == '/' {
			lastSlash = i
			break
		}
	}
	cur := irFactory_indexOf(units, '\\', 0)

	result := make([]uint16, 0, len(units))
	start := 1
	for cur != -1 {
		result = append(result, units[start:cur]...)

		cur++ // skip the escape char.
		c := irFactory_charAt(units, cur)
		switch c {
		// Characters for which the backslash is semantically important.
		case '^', '$', '\\', '/', '.', '*', '+', '?', '(', ')', '[', ']', '{', '}', '|', '-',
			'b', 'B', 'c', 'd', 'D', 'f', 'n',
			'k', // 2018 named backreferences
			'p', // 2018 unicode property escapes
			'P', // 2018 unicode property escapes
			'r', 's', 'S', 't', 'u', 'v', 'w', 'W', 'x',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			result = append(result, '\\')
			fallthrough
		default:
			// For all other characters, the backslash has no effect, so just
			// append the next char.
			result = append(result, c)
		}

		start = cur + 1
		cur = irFactory_indexOf(units, '\\', start)
	}
	if start > lastSlash {
		panic(jutil.NewStringIndexOutOfBoundsException(start, lastSlash))
	}
	result = append(result, units[start:lastSlash]...)
	return jutil.StringFromUTF16(result)
}

func (f *IRFactory) normalizeStringToID(token parser.Token) strpool.ID {
	checkTokenNotNull(token)
	if token.MaybePrivateValueID() == 0 {
		start, end := token.ValueRange()
		if s, ok := f.fileWithContent.ASCIISlice(start, end); ok {
			if strings.IndexByte(s, '\\') == -1 && strings.IndexByte(s, '\r') == -1 {
				return f.internString(s[1 : len(s)-1])
			}
		} else if units, ok := f.fileWithContent.UTF16Slice(start, end); ok && len(units) >= 2 {
			inner := units[1 : len(units)-1]
			if len(inner) <= 64 {
				var buf [64]byte
				simpleASCII := true
				for i, u := range inner {
					if u >= 0x80 || u == '\\' || u == '\r' {
						simpleASCII = false
						break
					}
					buf[i] = byte(u)
				}
				if simpleASCII {
					return f.internStringBytes(buf[:len(inner)])
				}
			}
		}
	}
	return f.internString(f.NormalizeString(token, false))
}

// Java: IRFactory.normalizeString(LiteralToken, boolean)
func (f *IRFactory) NormalizeString(token parser.Token, templateLiteral bool) string {
	checkTokenNotNull(token)
	if token.MaybePrivateValueID() == 0 {
		startOff, endOff := token.ValueRange()
		if s, ok := f.fileWithContent.ASCIISlice(startOff, endOff); ok {
			if strings.IndexByte(s, '\r') == -1 && strings.IndexByte(s, '\\') == -1 {
				if templateLiteral {
					return s
				}
				return s[1 : len(s)-1]
			}
		} else if units, ok := f.fileWithContent.UTF16Slice(startOff, endOff); ok {
			hasCR := false
			firstBackslash := -1
			for i, u := range units {
				if u == '\r' {
					hasCR = true
					break
				}
				if u == '\\' && firstBackslash < 0 {
					firstBackslash = i
				}
			}
			if !hasCR {
				if firstBackslash < 0 {
					if templateLiteral {
						return jutil.StringFromUTF16(units)
					}
					return jutil.StringFromUTF16(units[1 : len(units)-1])
				}
				return f.normalizeStringEscapes(units, firstBackslash, token, templateLiteral)
			}
		}
	}

	value := f.rawTokenValue(token)
	// <CR><LF> and <CR> are normalized as <LF>. For raw template literal
	// string values: this is the spec behaviour. For regular string literals:
	// they can only be part of a line continuation, which we want to scrub.
	if strings.IndexByte(value, '\r') >= 0 {
		// Go-only: language-adapter — value.replaceAll("\r\n?", "\n"); the
		// pattern has no regex features beyond the optional '\n'.
		value = strings.ReplaceAll(strings.ReplaceAll(value, "\r\n", "\n"), "\r", "\n")
	}

	if strings.IndexByte(value, '\\') == -1 {
		// short circuit no escapes.
		if templateLiteral {
			return value
		}
		return irFactory_substringWithoutQuotes(value)
	}
	// Go-only: language-adapter — the String and StringBuilder operations of
	// the escape path work on the UTF-16 code units of the value.
	units := jutil.StringToUTF16(value)
	return f.normalizeStringEscapes(units, irFactory_indexOf(units, '\\', 0), token, templateLiteral)
}

func (f *IRFactory) normalizeStringEscapes(units []uint16, cur int, token parser.Token, templateLiteral bool) string {
	start := 1 // skip the leading quote
	if templateLiteral {
		start = 0
	}
	result := make([]uint16, 0, len(units))
	tokenStart := token.Location.Start
	for cur != -1 {
		result = append(result, units[start:cur]...)

		cur += 1 // skip the escape char.
		c := irFactory_charAt(units, cur)
		switch c {
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
			result = append(result, '\u000B')
		case '\n':
			// line continuation, skip the line break
			f.MaybeWarnForFeatureToken(token, parser.FeatureSet_Feature_STRING_CONTINUATION)
			f.errorReporter.Warning(
				IRFactory_STRING_CONTINUATION_WARNING,
				f.sourceName,
				f.LinenoOfPosition(tokenStart),
				f.CharnoOfPosition(tokenStart))
		case '0', '1', '2', '3', '4', '5', '6', '7':
			var numDigits int

			if cur+1 < len(units) && IRFactory_IsOctalDigit(units[cur+1]) {
				// A three digit escape is only legal when the leading digit is 0-3; with a leading
				// digit of 4-7 the value would exceed \377, so only two digits are consumed and the
				// third octal digit is a literal character (ECMAScript Annex B.1.2
				// LegacyOctalEscapeSequence).
				if c <= '3' && cur+2 < len(units) && IRFactory_IsOctalDigit(units[cur+2]) {
					numDigits = 3
				} else {
					numDigits = 2
				}
			} else {
				numDigits = 1
			}

			if f.InStrictContext() || templateLiteral {
				if c == '0' && numDigits == 1 {
					// No warning: "\0" followed by a character which is not an
					// octal digit is allowed in strict mode.
				} else {
					f.errorReporter.Warning(
						IRFactory_OCTAL_STRING_LITERAL_WARNING,
						f.sourceName,
						f.LinenoOfPosition(tokenStart),
						f.CharnoOfPosition(tokenStart))
				}
			}

			result = append(result, uint16(jutil.IntegerParseInt(irFactory_substring(units, cur, cur+numDigits), 8)))
			cur += numDigits - 1
		case 'x':
			result = append(result,
				uint16(IRFactory_Hexdigit(irFactory_charAt(units, cur+1))*0x10+IRFactory_Hexdigit(irFactory_charAt(units, cur+2))))
			cur += 2
		case 'u':
			var escapeEnd int
			var hexDigits string
			if irFactory_charAt(units, cur+1) != '{' {
				// Simple escape with exactly four hex digits: \\uXXXX
				escapeEnd = cur + 5
				hexDigits = irFactory_substring(units, cur+1, escapeEnd)
			} else {
				// Escape with braces can have any number of hex digits:
				// \\u{XXXXXXX}
				escapeEnd = cur + 2
				for jutil.CharacterDigit(irFactory_charAt(units, escapeEnd), 0x10) >= 0 {
					escapeEnd++
				}
				hexDigits = irFactory_substring(units, cur+2, escapeEnd)
				escapeEnd++
			}
			codePointValue := jutil.IntegerParseInt(hexDigits, 0x10)
			if codePointValue > 0x10ffff {
				f.errorReporter.Error(
					"Undefined Unicode code-point",
					f.sourceName,
					f.LinenoOfPosition(tokenStart),
					f.CharnoOfPosition(tokenStart))

				// Compilation should stop, but we should finish the string and
				// find more errors. These appends are just to have a
				// placeholder for the errored normalization.
				result = append(result, '\\', 'u', '{')
				result = append(result, jutil.StringToUTF16(hexDigits)...)
				result = append(result, '}')
			} else {
				result = irFactory_appendCodePoint(result, codePointValue)
			}
			cur = escapeEnd - 1
		case '\'', '"', '\\':
			result = append(result, c)
		default:
			result = append(result, c)
		}
		start = cur + 1
		cur = irFactory_indexOf(units, '\\', start)
	}
	// skip the trailing quote.
	end := len(units) - 1
	if templateLiteral {
		end = len(units)
	}
	if start > end {
		panic(jutil.NewStringIndexOutOfBoundsException(start, end))
	}
	result = append(result, units[start:end]...)

	return jutil.StringFromUTF16(result)
}

// irFactory_substringWithoutQuotes is value.substring(1, value.length() - 1).
//
// Go-only: language-adapter — the first char of a string token is an ASCII
// quote; when the last byte is ASCII too, the UTF-16 range [1, length - 1)
// is the byte range [1, len - 1).
func irFactory_substringWithoutQuotes(value string) string {
	if n := len(value); n >= 2 && value[0] < 0x80 && value[n-1] < 0x80 {
		return value[1 : n-1]
	}
	return jutil.StringSubstring(value, 1, jutil.StringLength(value)-1)
}

// irFactory_appendCodePoint is StringBuilder.append(Character.toChars(int)).
//
// Go-only: language-adapter — JDK methods on UTF-16 code units.
func irFactory_appendCodePoint(units []uint16, codePoint int32) []uint16 {
	if codePoint < 0 || codePoint > 0x10ffff {
		panic(jutil.NewIllegalArgumentExceptionMsg(
			"Not a valid Unicode code point: 0x" + strings.ToUpper(jutil.IntegerToHexString(codePoint))))
	}
	if codePoint < 0x10000 {
		return append(units, uint16(codePoint))
	}
	cp := codePoint - 0x10000
	return append(units, uint16(0xD800+(cp>>10)), uint16(0xDC00+(cp&0x3FF)))
}

// Java: IRFactory.isSupportedForInputLanguageMode(Feature)
func (f *IRFactory) IsSupportedForInputLanguageMode(feature parser.FeatureSet_Feature) bool {
	return f.config.LanguageMode().FeatureSet().Has(feature)
}

// Java: IRFactory.inStrictContext()
func (f *IRFactory) InStrictContext() bool {
	// TODO(johnlenz): in ECMASCRIPT5/6 is a "mixed" mode and we should track
	// the context that we are in, if we want to support it.
	return f.config.StrictMode().IsStrict()
}

// NormalizeNumber is the value of a NUMBER token.
//
// Go-only: language-adapter — a NUMBER token's value only holds ASCII
// characters (digits, '_', '.', exponent and radix letters, signs; see
// Scanner), so byte indices are Java char indices.
//
// Java: IRFactory.normalizeNumber(LiteralToken)
func (f *IRFactory) NormalizeNumber(token parser.Token) float64 {
	if !token.HasNumericSeparator() && token.MaybePrivateValueID() == 0 {
		startOff := int(token.Location.Start.Offset)
		endOff := int(token.Location.End.Offset)
		if units, ok := f.fileWithContent.UTF16Slice(startOff, endOff); ok {
			length := len(units)
			if length == 1 && units[0] == '0' {
				return 0
			}
			if length <= 15 && units[0] >= '1' && units[0] <= '9' {
				v := uint64(units[0] - '0')
				allDigits := true
				for i := 1; i < length; i++ {
					d := units[i] - '0'
					if d > 9 {
						allDigits = false
						break
					}
					v = v*10 + uint64(d)
				}
				if allDigits {
					return float64(v)
				}
			} else if length > 2 && units[0] == '0' && (units[1] == 'x' || units[1] == 'X') {
				v := 0.0
				for c := 2; c < length; c++ {
					v = (v * 0x10) + float64(IRFactory_Hexdigit(units[c]))
				}
				return v
			}
		}
	}

	value := f.rawTokenValue(token)
	if token.HasNumericSeparator() {
		value = f.RemoveNumericSeparators(value, token)
	}
	location := token.Location
	length := len(value)
	jutil.CheckState(length > 0)
	jutil.CheckState(value[0] != '-' && value[0] != '+')
	if value[0] == '.' {
		return jutil.DoubleParseDouble("0" + value)
	} else if value[0] == '0' && length > 1 {
		switch value[1] {
		case '.', 'e', 'E':
			return jutil.DoubleParseDouble(value)
		case 'b', 'B':
			f.MaybeWarnForFeatureToken(token, parser.FeatureSet_Feature_BINARY_LITERALS)
			v := 0.0
			c := 1
			for c++; c < length; c++ {
				v = (v * 2) + float64(IRFactory_Binarydigit(uint16(value[c])))
			}
			return v
		case 'o', 'O':
			f.MaybeWarnForFeatureToken(token, parser.FeatureSet_Feature_OCTAL_LITERALS)
			v := 0.0
			c := 1
			for c++; c < length; c++ {
				v = (v * 8) + float64(IRFactory_Octaldigit(uint16(value[c])))
			}
			return v
		case 'x', 'X':
			v := 0.0
			c := 1
			for c++; c < length; c++ {
				v = (v * 0x10) + float64(IRFactory_Hexdigit(uint16(value[c])))
			}
			return v
		case '0', '1', '2', '3', '4', '5', '6', '7':
			v := 0.0
			c := 0
			for c++; c < length; c++ {
				digit := uint16(value[c])
				if IRFactory_IsOctalDigit(digit) {
					v = (v * 8) + float64(IRFactory_Octaldigit(digit))
				} else {
					f.errorReporter.Error(
						IRFactory_INVALID_OCTAL_DIGIT, f.sourceName, f.LinenoOfPosition(location.Start), f.CharnoOfPosition(location.Start))
					return 0
				}
			}
			if f.InStrictContext() {
				f.errorReporter.Error(
					IRFactory_INVALID_ES5_STRICT_OCTAL,
					f.sourceName,
					f.LinenoOfPosition(location.Start),
					f.CharnoOfPosition(location.Start))
			} else {
				f.errorReporter.Warning(
					IRFactory_INVALID_ES5_STRICT_OCTAL,
					f.sourceName,
					f.LinenoOfPosition(location.Start),
					f.CharnoOfPosition(location.Start))
			}
			return v
		case '8', '9':
			f.errorReporter.Error(
				IRFactory_INVALID_OCTAL_DIGIT, f.sourceName, f.LinenoOfPosition(location.Start), f.CharnoOfPosition(location.Start))
			return 0
		default:
			panic(jutil.NewIllegalStateExceptionMsg(
				"Unexpected character in number literal: " + value[1:2]))
		}
	} else if value == "0" {
		return 0
	} else {
		if length <= 15 && value[0] >= '1' && value[0] <= '9' {
			v := uint64(value[0] - '0')
			allDigits := true
			for i := 1; i < length; i++ {
				d := value[i] - '0'
				if d > 9 {
					allDigits = false
					break
				}
				v = v*10 + uint64(d)
			}
			if allDigits {
				return float64(v)
			}
		}
		return jutil.DoubleParseDouble(value)
	}
}

// NormalizeBigInt is the value of a BIGINT token.
//
// Go-only: language-adapter — a BIGINT token's value only holds ASCII
// characters (see NormalizeNumber).
//
// Java: IRFactory.normalizeBigInt(LiteralToken)
func (f *IRFactory) NormalizeBigInt(token parser.Token) *big.Int {
	value := f.rawTokenValue(token)
	n := strings.IndexByte(value, 'n')
	if n < 0 {
		panic(jutil.NewStringIndexOutOfBoundsException(0, n))
	}
	value = value[0:n]
	if token.HasNumericSeparator() {
		value = f.RemoveNumericSeparators(value, token)
	}
	length := len(value)
	jutil.CheckState(length > 0)
	jutil.CheckState(value[0] != '-' && value[0] != '+')
	if value[0] == '0' && length > 1 {
		switch value[1] {
		case 'b', 'B':
			f.MaybeWarnForFeatureToken(token, parser.FeatureSet_Feature_BINARY_LITERALS)
			return jutil.NewBigIntegerFromStringRadix(value[2:], 2)
		case 'o', 'O':
			f.MaybeWarnForFeatureToken(token, parser.FeatureSet_Feature_OCTAL_LITERALS)
			return jutil.NewBigIntegerFromStringRadix(value[2:], 8)
		case 'x', 'X':
			return jutil.NewBigIntegerFromStringRadix(value[2:], 16)
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			panic(jutil.NewIllegalStateExceptionMsg("Nonzero BigInts can't have a leading zero"))
		default:
			panic(jutil.NewIllegalStateExceptionMsg(
				"Unexpected character in bigint literal: " + value[1:2]))
		}
	} else {
		return jutil.NewBigIntegerFromString(value)
	}
}

// Java: IRFactory.removeNumericSeparators(String, LiteralToken)
func (f *IRFactory) RemoveNumericSeparators(value string, token parser.Token) string {
	f.MaybeWarnForFeatureToken(token, parser.FeatureSet_Feature_NUMERIC_SEPARATOR)
	return jutil.StringReplace(value, "_", "")
}

// Java: IRFactory.binarydigit(char)
func IRFactory_Binarydigit(c uint16) int {
	if c >= '0' && c <= '1' {
		return int(c - '0')
	}
	panic(jutil.NewIllegalStateExceptionMsg("unexpected: " + jutil.StringFromChar(c)))
}

// Java: IRFactory.isOctalDigit(char)
func IRFactory_IsOctalDigit(c uint16) bool {
	return c >= '0' && c <= '7'
}

// Java: IRFactory.octaldigit(char)
func IRFactory_Octaldigit(c uint16) int {
	if IRFactory_IsOctalDigit(c) {
		return int(c - '0')
	}
	panic(jutil.NewIllegalStateExceptionMsg("unexpected: " + jutil.StringFromChar(c)))
}

// Java: IRFactory.hexdigit(char)
func IRFactory_Hexdigit(c uint16) int {
	switch c {
	case '0':
		return 0
	case '1':
		return 1
	case '2':
		return 2
	case '3':
		return 3
	case '4':
		return 4
	case '5':
		return 5
	case '6':
		return 6
	case '7':
		return 7
	case '8':
		return 8
	case '9':
		return 9
	case 'a', 'A':
		return 10
	case 'b', 'B':
		return 11
	case 'c', 'C':
		return 12
	case 'd', 'D':
		return 13
	case 'e', 'E':
		return 14
	case 'f', 'F':
		return 15
	default:
		panic(jutil.NewIllegalStateExceptionMsg("unexpected: " + jutil.StringFromChar(c)))
	}
}

// Java: IRFactory.transformBooleanTokenType(TokenType)
func IRFactory_TransformBooleanTokenType(token parser.TokenType) rhino.Token {
	switch token {
	case parser.TokenType_TRUE:
		return rhino.Token_TRUE
	case parser.TokenType_FALSE:
		return rhino.Token_FALSE
	default:
		panic(jutil.NewIllegalStateExceptionMsg(token.String()))
	}
}

// Java: IRFactory.transformUpdateTokenType(TokenType)
func IRFactory_TransformUpdateTokenType(token parser.TokenType) rhino.Token {
	switch token {
	case parser.TokenType_PLUS_PLUS:
		return rhino.Token_INC
	case parser.TokenType_MINUS_MINUS:
		return rhino.Token_DEC
	default:
		panic(jutil.NewIllegalStateExceptionMsg(token.String()))
	}
}

// Java: IRFactory.transformUnaryTokenType(TokenType)
func IRFactory_TransformUnaryTokenType(token parser.TokenType) rhino.Token {
	switch token {
	case parser.TokenType_BANG:
		return rhino.Token_NOT
	case parser.TokenType_TILDE:
		return rhino.Token_BITNOT
	case parser.TokenType_PLUS:
		return rhino.Token_POS
	case parser.TokenType_MINUS:
		return rhino.Token_NEG
	case parser.TokenType_DELETE:
		return rhino.Token_DELPROP
	case parser.TokenType_TYPEOF:
		return rhino.Token_TYPEOF
	case parser.TokenType_VOID:
		return rhino.Token_VOID
	default:
		panic(jutil.NewIllegalStateExceptionMsg(token.String()))
	}
}

// Java: IRFactory.transformBinaryTokenType(TokenType)
func IRFactory_TransformBinaryTokenType(token parser.TokenType) rhino.Token {
	switch token {
	case parser.TokenType_BAR:
		return rhino.Token_BITOR
	case parser.TokenType_CARET:
		return rhino.Token_BITXOR
	case parser.TokenType_AMPERSAND:
		return rhino.Token_BITAND
	case parser.TokenType_EQUAL_EQUAL:
		return rhino.Token_EQ
	case parser.TokenType_NOT_EQUAL:
		return rhino.Token_NE
	case parser.TokenType_OPEN_ANGLE:
		return rhino.Token_LT
	case parser.TokenType_LESS_EQUAL:
		return rhino.Token_LE
	case parser.TokenType_CLOSE_ANGLE:
		return rhino.Token_GT
	case parser.TokenType_GREATER_EQUAL:
		return rhino.Token_GE
	case parser.TokenType_LEFT_SHIFT:
		return rhino.Token_LSH
	case parser.TokenType_RIGHT_SHIFT:
		return rhino.Token_RSH
	case parser.TokenType_UNSIGNED_RIGHT_SHIFT:
		return rhino.Token_URSH
	case parser.TokenType_PLUS:
		return rhino.Token_ADD
	case parser.TokenType_MINUS:
		return rhino.Token_SUB
	case parser.TokenType_STAR:
		return rhino.Token_MUL
	case parser.TokenType_SLASH:
		return rhino.Token_DIV
	case parser.TokenType_PERCENT:
		return rhino.Token_MOD
	case parser.TokenType_STAR_STAR:
		return rhino.Token_EXPONENT
	case parser.TokenType_EQUAL_EQUAL_EQUAL:
		return rhino.Token_SHEQ
	case parser.TokenType_NOT_EQUAL_EQUAL:
		return rhino.Token_SHNE
	case parser.TokenType_IN:
		return rhino.Token_IN
	case parser.TokenType_INSTANCEOF:
		return rhino.Token_INSTANCEOF
	case parser.TokenType_COMMA:
		return rhino.Token_COMMA
	case parser.TokenType_EQUAL:
		return rhino.Token_ASSIGN
	case parser.TokenType_BAR_EQUAL:
		return rhino.Token_ASSIGN_BITOR
	case parser.TokenType_CARET_EQUAL:
		return rhino.Token_ASSIGN_BITXOR
	case parser.TokenType_AMPERSAND_EQUAL:
		return rhino.Token_ASSIGN_BITAND
	case parser.TokenType_LEFT_SHIFT_EQUAL:
		return rhino.Token_ASSIGN_LSH
	case parser.TokenType_RIGHT_SHIFT_EQUAL:
		return rhino.Token_ASSIGN_RSH
	case parser.TokenType_UNSIGNED_RIGHT_SHIFT_EQUAL:
		return rhino.Token_ASSIGN_URSH
	case parser.TokenType_PLUS_EQUAL:
		return rhino.Token_ASSIGN_ADD
	case parser.TokenType_MINUS_EQUAL:
		return rhino.Token_ASSIGN_SUB
	case parser.TokenType_STAR_EQUAL:
		return rhino.Token_ASSIGN_MUL
	case parser.TokenType_STAR_STAR_EQUAL:
		return rhino.Token_ASSIGN_EXPONENT
	case parser.TokenType_SLASH_EQUAL:
		return rhino.Token_ASSIGN_DIV
	case parser.TokenType_PERCENT_EQUAL:
		return rhino.Token_ASSIGN_MOD
	case parser.TokenType_OR:
		return rhino.Token_OR
	case parser.TokenType_AND:
		return rhino.Token_AND
	case parser.TokenType_QUESTION_QUESTION:
		return rhino.Token_COALESCE
	case parser.TokenType_OR_EQUAL:
		return rhino.Token_ASSIGN_OR
	case parser.TokenType_AND_EQUAL:
		return rhino.Token_ASSIGN_AND
	case parser.TokenType_QUESTION_QUESTION_EQUAL:
		return rhino.Token_ASSIGN_COALESCE
	default:
		panic(jutil.NewIllegalStateExceptionMsg(token.String()))
	}
}
