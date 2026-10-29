// Port of com/google/javascript/rhino/TokenStream.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/TokenUtil.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/JSIdentifier.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/ClosurePrimitive.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/Outcome.java @ closure-compiler v20261021

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Roger Lawrence
 *   Mike McCabe
 *   Igor Bukanov
 *   Ethan Hugg
 *   Bob Jervis
 *   Terry Lucas
 *   Milen Nankov
 *   Pascal-Louis Perez
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Bob Jervis
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Annie Wang
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhino

import (
	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// TokenStream implements the JavaScript scanner.
//
// It is based on the C source files jsscan.c and jsscan.h in the jsref
// package.
//
// IMPORTANT: As of 2018-03-09 it is still not possible to use Java 8
// features in this file due to limitations on some internal Google projects
// that depend on it.
type TokenStream struct{}

// TokenStream_IsKeyword reports whether the string is an ES3 keyword.
//
// Keywords for versions of JavaScript after ES3 are not included, because
// the parser would reject valid ES3 code that happened to newer keywords if
// we did that.
//
// Java: TokenStream.isKeyword(String)
func TokenStream_IsKeyword(name string) bool {
	switch name {
	case "do", "if", "in",
		"for", "int", "new", "try", "var",
		"byte", "case", "char", "else", "enum", "goto", "long", "null", "this", "true", "void", "with",
		"break", "catch", "class", "const", "false", "final", "float", "short", "super", "throw", "while",
		"delete", "double", "export", "import", "native", "public", "return", "static", "switch", "throws", "typeof",
		"boolean", "default", "extends", "finally", "package", "private",
		"abstract", "continue", "debugger", "function", "volatile",
		"interface", "protected", "transient",
		"implements", "instanceof",
		"synchronized":
		return true
	default:
		return false
	}
}

// Java: TokenStream.isJSIdentifier(String)
func TokenStream_IsJSIdentifier(s string) bool {
	return JSIdentifier_IsJSIdentifier(s)
}

// Java: TokenStream()
func NewTokenStream() *TokenStream { return &TokenStream{} }

// TokenUtil holds helper methods for parsing JavaScript.
type TokenUtil struct{}

// TokenUtil_IsJSSpace is as defined in ECMA. jsscan.c uses C isspace()
// (which allows \v, I think.) note that code in getChar() implicitly accepts
// '\r' == \u000D as well.
//
// Java: TokenUtil.isJSSpace(int)
func TokenUtil_IsJSSpace(c int) bool {
	if c <= 127 {
		return c == 0x20 || c == 0x9 || c == 0xC || c == 0xB
	} else {
		return c == 0xA0 ||
			jutil.CharacterGetType(uint16(c)) == jutil.Character_SPACE_SEPARATOR
	}
}

// Java: TokenUtil.isJSFormatChar(int)
func TokenUtil_IsJSFormatChar(c int) bool {
	return c > 127 && jutil.CharacterGetTypeOfCodePoint(c) == jutil.Character_FORMAT
}

// Java: TokenUtil.isWhitespace(int)
func TokenUtil_IsWhitespace(c int) bool {
	return jutil.CharacterIsWhitespaceCodePoint(c)
}

// TokenUtil_IsStrWhiteSpaceChar is copied from Rhino's ScriptRuntime.
//
// Java: TokenUtil.isStrWhiteSpaceChar(int)
func TokenUtil_IsStrWhiteSpaceChar(c int) base.Tri {
	switch c {
	case 0x0B: // <VT>
		return base.Tri_UNKNOWN // IE says "no", ECMAScript says "yes"
	case ' ', // <SP>
		'\n',     // <LF>
		'\r',     // <CR>
		'\t',     // <TAB>
		0xA0,     // <NBSP>
		'\u000C', // <FF>
		0x2028,   // <LS>
		0x2029,   // <PS>
		0xFEFF:   // <BOM>
		return base.Tri_TRUE
	default:
		if jutil.CharacterGetTypeOfCodePoint(c) == jutil.Character_SPACE_SEPARATOR {
			return base.Tri_TRUE
		}
		return base.Tri_FALSE
	}
}

// Java: TokenUtil()
func NewTokenUtil() *TokenUtil { return &TokenUtil{} }

// JSIdentifier is a utility class to hold isJSIdentifier.
//
// IMPORTANT: As of 2018-03-09 it is still not possible to use Java 8
// features in this file due to limitations on some internal Google projects
// that depend on it.
type JSIdentifier struct{}

var asciiIdentStart = [128]uint8{
	'$': 1,
	'A': 1, 'B': 1, 'C': 1, 'D': 1, 'E': 1, 'F': 1, 'G': 1, 'H': 1, 'I': 1, 'J': 1, 'K': 1, 'L': 1, 'M': 1,
	'N': 1, 'O': 1, 'P': 1, 'Q': 1, 'R': 1, 'S': 1, 'T': 1, 'U': 1, 'V': 1, 'W': 1, 'X': 1, 'Y': 1, 'Z': 1,
	'_': 1,
	'a': 1, 'b': 1, 'c': 1, 'd': 1, 'e': 1, 'f': 1, 'g': 1, 'h': 1, 'i': 1, 'j': 1, 'k': 1, 'l': 1, 'm': 1,
	'n': 1, 'o': 1, 'p': 1, 'q': 1, 'r': 1, 's': 1, 't': 1, 'u': 1, 'v': 1, 'w': 1, 'x': 1, 'y': 1, 'z': 1,
}

var asciiIdentPart = [128]uint8{
	'$': 1,
	'0': 1, '1': 1, '2': 1, '3': 1, '4': 1, '5': 1, '6': 1, '7': 1, '8': 1, '9': 1,
	'A': 1, 'B': 1, 'C': 1, 'D': 1, 'E': 1, 'F': 1, 'G': 1, 'H': 1, 'I': 1, 'J': 1, 'K': 1, 'L': 1, 'M': 1,
	'N': 1, 'O': 1, 'P': 1, 'Q': 1, 'R': 1, 'S': 1, 'T': 1, 'U': 1, 'V': 1, 'W': 1, 'X': 1, 'Y': 1, 'Z': 1,
	'_': 1,
	'a': 1, 'b': 1, 'c': 1, 'd': 1, 'e': 1, 'f': 1, 'g': 1, 'h': 1, 'i': 1, 'j': 1, 'k': 1, 'l': 1, 'm': 1,
	'n': 1, 'o': 1, 'p': 1, 'q': 1, 'r': 1, 's': 1, 't': 1, 'u': 1, 'v': 1, 'w': 1, 'x': 1, 'y': 1, 'z': 1,
}

// Java: JSIdentifier.isJSIdentifier(String)
func JSIdentifier_IsJSIdentifier(s string) bool {
	if len(s) == 0 {
		return false
	}
	c0 := s[0]
	if c0 >= 0x80 {
		return isJSIdentifierUTF16(s)
	}
	if asciiIdentStart[c0] == 0 {
		return false
	}
	for i := 1; i < len(s); i++ {
		c := s[i]
		if c >= 0x80 {
			return isJSIdentifierUTF16(s)
		}
		if asciiIdentPart[c] == 0 {
			return false
		}
	}
	return true
}

//go:noinline
func isJSIdentifierUTF16(s string) bool {
	chars := jutil.StringToUTF16(s)
	length := len(chars)
	if length == 0 ||
		jutil.CharacterIsIdentifierIgnorable(chars[0]) ||
		!jutil.CharacterIsJavaIdentifierStart(chars[0]) {
		return false
	}
	for i := 1; i < length; i++ {
		if jutil.CharacterIsIdentifierIgnorable(chars[i]) ||
			!jutil.CharacterIsJavaIdentifierPart(chars[i]) {
			return false
		}
	}
	return true
}

// Java: JSIdentifier()
func NewJSIdentifier() *JSIdentifier { return &JSIdentifier{} }

// ClosurePrimitive enumerates the primitive functions that the compiler
// recognizes.
//
// These correspond to the @closurePrimitive tag in code; in order to parse
// new primitives, add any entry to the list in
// parsing/ParserConfig.properties, then map it to an enum member in the
// idToEnum map.
//
// After typechecking is done, all calls to one of these primitive types
// should have their FunctionType annotated with the corresponding enum
// member. This is intended to make identifying these calls more accurate
// than previous methods of finding primitives by qualified name.
type ClosurePrimitive int32

const (
	ClosurePrimitive_ASSERTS_FAIL           ClosurePrimitive = iota // A function that always throws an error
	ClosurePrimitive_ASSERTS_MATCHES_RETURN                         // A function that asserts its first parameter matches the return type
	ClosurePrimitive_ASSERTS_TRUTHY                                 // A function that asserts its first parameter is truthy and returns the param
)

var closurePrimitive_names = [...]string{"ASSERTS_FAIL", "ASSERTS_MATCHES_RETURN", "ASSERTS_TRUTHY"}

var closurePrimitive_values = [...]ClosurePrimitive{
	ClosurePrimitive_ASSERTS_FAIL,
	ClosurePrimitive_ASSERTS_MATCHES_RETURN,
	ClosurePrimitive_ASSERTS_TRUTHY,
}

// String is Enum.name().
func (p ClosurePrimitive) String() string { return closurePrimitive_names[p] }

// ClosurePrimitive_ID_TO_ENUM maps human-readable ids to enum members.
//
// The expected mapping (although not enforced) of keys -> values is that
// the enum member maps to a lowercase string with "_" replaced with "."
//
// Java: ClosurePrimitive.ID_TO_ENUM
var ClosurePrimitive_ID_TO_ENUM = func() *jutil.LinkedHashMap[string, ClosurePrimitive] {
	m := jutil.NewLinkedHashMap[string, ClosurePrimitive]()
	m.Put("asserts.fail", ClosurePrimitive_ASSERTS_FAIL)
	m.Put("asserts.truthy", ClosurePrimitive_ASSERTS_TRUTHY)
	m.Put("asserts.matchesReturn", ClosurePrimitive_ASSERTS_MATCHES_RETURN)
	return m
}()

// ClosurePrimitive_FromStringId returns the ClosurePrimitive corresponding to
// the given string id.
//
// This is to make reading @closurePrimitive easier in code. Using
// Enum.valueOf to parse closure primitive identifiers from JSDoc directly
// would require code like `closurePrimitive {ASSERTS_FAIL}`; instead we
// separate the string ids from the enum names.
//
// Returns null if the argument is null, otherwise the corresponding enum
// member. Throws IllegalArgumentException if the id is non-null but does not
// match an enum member.
//
// Java: ClosurePrimitive.fromStringId(String)
func ClosurePrimitive_FromStringId(id *string) *ClosurePrimitive {
	if id == nil {
		return nil
	}
	switch *id {
	case "asserts.fail":
		return &closurePrimitive_values[ClosurePrimitive_ASSERTS_FAIL]
	case "asserts.truthy":
		return &closurePrimitive_values[ClosurePrimitive_ASSERTS_TRUTHY]
	case "asserts.matchesReturn":
		return &closurePrimitive_values[ClosurePrimitive_ASSERTS_MATCHES_RETURN]
	default:
		panic(jutil.NewIllegalArgumentExceptionMsg(
			"String id " + *id + " does not match any ClosurePrimitive"))
	}
}

// Outcome is an enum for representing truthiness and nullishness of
// outcomes. The values are TRUE, FALSE, FALSE_NOT_NULL, and NULLISH. FALSE
// represents all falsy values so it contains both FALSE_NOT_NULL and
// NULLISH. It specifically differentiates between falsy and explicitly
// nullish values. e.g. 0 is falsy but not nullish.
type Outcome int32

const (
	// Represents truthy values. For example: {}, true, 1, etc.
	Outcome_TRUE Outcome = iota
	// Represents falsy values. For examples: '', 0, false, null, etc.
	Outcome_FALSE
	Outcome_FALSE_NOT_NULL
	Outcome_NULLISH
)

var outcome_names = [...]string{"TRUE", "FALSE", "FALSE_NOT_NULL", "NULLISH"}

// Outcome_values is Outcome.values().
func Outcome_values() []Outcome {
	return []Outcome{Outcome_TRUE, Outcome_FALSE, Outcome_FALSE_NOT_NULL, Outcome_NULLISH}
}

// String is Enum.name().
func (o Outcome) String() string { return outcome_names[o] }

// IsTruthy determines whether an Outcome enum value is truthy.
//
// Java: Outcome.isTruthy()
func (o Outcome) IsTruthy() bool {
	switch o {
	case Outcome_TRUE:
		return true
	}
	return false
}

// IsNullish determines whether an Outcome enum value is nullish. Using Tri
// instead of a boolean because 0 is Outcome.FALSE but not nullish so
// sometimes it is unclear.
//
// Java: Outcome.isNullish()
func (o Outcome) IsNullish() base.Tri {
	switch o {
	case Outcome_TRUE:
		return base.Tri_FALSE
	case Outcome_FALSE:
		return base.Tri_UNKNOWN
	case Outcome_FALSE_NOT_NULL:
		return base.Tri_FALSE
	}
	return base.Tri_TRUE // NULLISH
}

// Not gets the `not` of `this`.
//
// Java: Outcome.not()
func (o Outcome) Not() Outcome {
	switch o {
	case Outcome_TRUE:
		return Outcome_FALSE
	}
	return Outcome_TRUE // FALSE, FALSE_NOT_NULL, NULLISH
}

// Outcome_ForBoolean gets the Outcome for the given boolean.
//
// Java: Outcome.forBoolean(boolean)
func Outcome_ForBoolean(val bool) Outcome {
	if val {
		return Outcome_TRUE
	}
	return Outcome_FALSE
}
