// Port of com/google/javascript/jscomp/PeepholeReplaceKnownMethods.java @ closure-compiler v20261021

/*
 * Copyright 2011 The Closure Compiler Authors.
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
	"math"
	"strings"
	"unicode/utf8"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

var (
	rkmIDSplit       = uint32(strpool.Intern("split"))
	rkmIDToLowerCase = uint32(strpool.Intern("toLowerCase"))
	rkmIDToUpperCase = uint32(strpool.Intern("toUpperCase"))
	rkmIDTrim        = uint32(strpool.Intern("trim"))
	rkmIDTrimStart   = uint32(strpool.Intern("trimStart"))
	rkmIDTrimEnd     = uint32(strpool.Intern("trimEnd"))
	rkmIDTrimLeft    = uint32(strpool.Intern("trimLeft"))
	rkmIDTrimRight   = uint32(strpool.Intern("trimRight"))
	rkmIDAt          = uint32(strpool.Intern("at"))
	rkmIDIndexOf     = uint32(strpool.Intern("indexOf"))
	rkmIDLastIndexOf = uint32(strpool.Intern("lastIndexOf"))
	rkmIDIncludes    = uint32(strpool.Intern("includes"))
	rkmIDStartsWith  = uint32(strpool.Intern("startsWith"))
	rkmIDEndsWith    = uint32(strpool.Intern("endsWith"))
	rkmIDSubstr      = uint32(strpool.Intern("substr"))
	rkmIDSubstring   = uint32(strpool.Intern("substring"))
	rkmIDSlice       = uint32(strpool.Intern("slice"))
	rkmIDCharAt      = uint32(strpool.Intern("charAt"))
	rkmIDCharCodeAt  = uint32(strpool.Intern("charCodeAt"))
	rkmIDReplace     = uint32(strpool.Intern("replace"))
	rkmIDReplaceAll  = uint32(strpool.Intern("replaceAll"))
)

// Java: PeepholeReplaceKnownMethods.TRIM_BOTH_PATTERN
var peepholeReplaceKnownMethods_TRIM_BOTH_PATTERN = peepholeReplaceKnownMethods_trimPattern{leading: true, trailing: true}

// Java: PeepholeReplaceKnownMethods.TRIM_START_PATTERN
var peepholeReplaceKnownMethods_TRIM_START_PATTERN = peepholeReplaceKnownMethods_trimPattern{leading: true}

// Java: PeepholeReplaceKnownMethods.TRIM_END_PATTERN
var peepholeReplaceKnownMethods_TRIM_END_PATTERN = peepholeReplaceKnownMethods_trimPattern{trailing: true}

// peepholeReplaceKnownMethods_trimPattern is one of the compiled patterns
// TRIM_BOTH_PATTERN (`^W|W\z`), TRIM_START_PATTERN (`^W`) and
// TRIM_END_PATTERN (`W\z`), where W is ECMA_WHITESPACE.
//
// Go-only: language-adapter — `\z` (jutil.Pattern rejects it), so this is an
// equivalent hand matcher for these java.util.regex uses (PORTING.md,
// Regex). With W a greedy run of whitespace characters, `^W` matches the
// maximal leading run, and `W\z` matches only a run that reaches the end of
// the input: the first position from which one does is the start of the
// maximal trailing run. Every character of W is in the BMP and is not a
// surrogate, so matching code points of the WTF-8 string is matching Java's.
type peepholeReplaceKnownMethods_trimPattern struct {
	leading, trailing bool
}

// replaceAll is Pattern.matcher(s).replaceAll("").
func (t peepholeReplaceKnownMethods_trimPattern) replaceAll(s string) string {
	start := 0
	if t.leading {
		for start < len(s) {
			r, size := utf8.DecodeRuneInString(s[start:])
			if !peepholeReplaceKnownMethods_isEcmaWhitespace(r) {
				break
			}
			start += size
		}
	}
	end := len(s)
	if t.trailing {
		for end > start {
			r, size := utf8.DecodeLastRuneInString(s[start:end])
			if !peepholeReplaceKnownMethods_isEcmaWhitespace(r) {
				break
			}
			end -= size
		}
	}
	return s[start:end]
}

// peepholeReplaceKnownMethods_isEcmaWhitespace is the character class of
// ECMA_WHITESPACE: [ \t\n-\r\u00A0\u1680\u2000-\u200A\u2028\u2029\u202F
// \u205F\u3000\uFEFF].
func peepholeReplaceKnownMethods_isEcmaWhitespace(r rune) bool {
	switch {
	case r == ' ' || r == '\t' || (r >= '\n' && r <= '\r'):
		return true
	case r == 0x00A0 || r == 0x1680 || (r >= 0x2000 && r <= 0x200A):
		return true
	case r == 0x2028 || r == 0x2029 || r == 0x202F || r == 0x205F || r == 0x3000 || r == 0xFEFF:
		return true
	}
	return false
}

// stringUpperCaseRootSpecial holds the unconditional Locale.ROOT uppercase
// expansions from Unicode 15.0 SpecialCasing.txt that differ from simple 1:1
// UnicodeData.txt uppercase mapping.
var stringUpperCaseRootSpecial = map[rune]string{
	0xdf: "SS", 0x149: "\u02bcN", 0x1f0: "J\u030c", 0x390: "\u0399\u0308\u0301", 0x3b0: "\u03a5\u0308\u0301",
	0x587: "\u0535\u0552", 0x1e96: "H\u0331", 0x1e97: "T\u0308", 0x1e98: "W\u030a", 0x1e99: "Y\u030a", 0x1e9a: "A\u02be",
	0x1f50: "\u03a5\u0313", 0x1f52: "\u03a5\u0313\u0300", 0x1f54: "\u03a5\u0313\u0301", 0x1f56: "\u03a5\u0313\u0342",
	0x1f80: "\u1f08\u0399", 0x1f81: "\u1f09\u0399", 0x1f82: "\u1f0a\u0399", 0x1f83: "\u1f0b\u0399",
	0x1f84: "\u1f0c\u0399", 0x1f85: "\u1f0d\u0399", 0x1f86: "\u1f0e\u0399", 0x1f87: "\u1f0f\u0399",
	0x1f88: "\u1f08\u0399", 0x1f89: "\u1f09\u0399", 0x1f8a: "\u1f0a\u0399", 0x1f8b: "\u1f0b\u0399",
	0x1f8c: "\u1f0c\u0399", 0x1f8d: "\u1f0d\u0399", 0x1f8e: "\u1f0e\u0399", 0x1f8f: "\u1f0f\u0399",
	0x1f90: "\u1f28\u0399", 0x1f91: "\u1f29\u0399", 0x1f92: "\u1f2a\u0399", 0x1f93: "\u1f2b\u0399",
	0x1f94: "\u1f2c\u0399", 0x1f95: "\u1f2d\u0399", 0x1f96: "\u1f2e\u0399", 0x1f97: "\u1f2f\u0399",
	0x1f98: "\u1f28\u0399", 0x1f99: "\u1f29\u0399", 0x1f9a: "\u1f2a\u0399", 0x1f9b: "\u1f2b\u0399",
	0x1f9c: "\u1f2c\u0399", 0x1f9d: "\u1f2d\u0399", 0x1f9e: "\u1f2e\u0399", 0x1f9f: "\u1f2f\u0399",
	0x1fa0: "\u1f68\u0399", 0x1fa1: "\u1f69\u0399", 0x1fa2: "\u1f6a\u0399", 0x1fa3: "\u1f6b\u0399",
	0x1fa4: "\u1f6c\u0399", 0x1fa5: "\u1f6d\u0399", 0x1fa6: "\u1f6e\u0399", 0x1fa7: "\u1f6f\u0399",
	0x1fa8: "\u1f68\u0399", 0x1fa9: "\u1f69\u0399", 0x1faa: "\u1f6a\u0399", 0x1fab: "\u1f6b\u0399",
	0x1fac: "\u1f6c\u0399", 0x1fad: "\u1f6d\u0399", 0x1fae: "\u1f6e\u0399", 0x1faf: "\u1f6f\u0399",
	0x1fb2: "\u1fba\u0399", 0x1fb3: "\u0391\u0399", 0x1fb4: "\u0386\u0399", 0x1fb6: "\u0391\u0342", 0x1fb7: "\u0391\u0342\u0399", 0x1fbc: "\u0391\u0399",
	0x1fc2: "\u1fca\u0399", 0x1fc3: "\u0397\u0399", 0x1fc4: "\u0389\u0399", 0x1fc6: "\u0397\u0342", 0x1fc7: "\u0397\u0342\u0399", 0x1fcc: "\u0397\u0399",
	0x1fd2: "\u0399\u0308\u0300", 0x1fd3: "\u0399\u0308\u0301", 0x1fd6: "\u0399\u0342", 0x1fd7: "\u0399\u0308\u0342",
	0x1fe2: "\u03a5\u0308\u0300", 0x1fe3: "\u03a5\u0308\u0301", 0x1fe4: "\u03a1\u0313", 0x1fe6: "\u03a5\u0342", 0x1fe7: "\u03a5\u0308\u0342",
	0x1ff2: "\u1ffa\u0399", 0x1ff3: "\u03a9\u0399", 0x1ff4: "\u038f\u0399", 0x1ff6: "\u03a9\u0342", 0x1ff7: "\u03a9\u0342\u0399", 0x1ffc: "\u03a9\u0399",
	0xfb00: "FF", 0xfb01: "FI", 0xfb02: "FL", 0xfb03: "FFI", 0xfb04: "FFL", 0xfb05: "ST", 0xfb06: "ST",
	0xfb13: "\u0544\u0546", 0xfb14: "\u0544\u0535", 0xfb15: "\u0544\u053b", 0xfb16: "\u054e\u0546", 0xfb17: "\u0544\u053d",
}

// stringToUpperCaseRoot implements String.toUpperCase(Locale.ROOT).
func stringToUpperCaseRoot(s string) string {
	if jutil.IsASCII(s) {
		return jutil.StringToUpperCaseASCII(s)
	}
	u := jutil.StringToUTF16(s)
	var out []uint16
	for i := 0; i < len(u); {
		c := u[i]
		n := 1
		cp := rune(c)
		if c >= 0xD800 && c <= 0xDBFF && i+1 < len(u) && u[i+1] >= 0xDC00 && u[i+1] <= 0xDFFF {
			cp = 0x10000 + (rune(c-0xD800)<<10 | rune(u[i+1]-0xDC00))
			n = 2
		}
		if spec, ok := stringUpperCaseRootSpecial[cp]; ok {
			if out == nil {
				out = append(make([]uint16, 0, len(u)+8), u[:i]...)
			}
			for _, r := range spec {
				if r < 0x10000 {
					out = append(out, uint16(r))
				} else {
					v := r - 0x10000
					out = append(out, uint16(0xD800+(v>>10)), uint16(0xDC00+(v&0x3FF)))
				}
			}
			i += n
			continue
		}
		upper := jutil.Unicode15ToUpper(cp)
		if upper != cp && out == nil {
			out = append(make([]uint16, 0, len(u)+8), u[:i]...)
		}
		if out != nil {
			if upper < 0x10000 {
				out = append(out, uint16(upper))
			} else {
				v := upper - 0x10000
				out = append(out, uint16(0xD800+(v>>10)), uint16(0xDC00+(v&0x3FF)))
			}
		}
		i += n
	}
	if out == nil {
		return s
	}
	return jutil.StringFromUTF16(out)
}

// TryFoldKnownStringMethods tries to evaluate known String methods
// .indexOf(), .substr(), .substring()
//
// Java: PeepholeReplaceKnownMethods.tryFoldKnownStringMethods(Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldKnownStringMethods(call ast.Call, callTarget ast.GetProp) rhino.Node {
	subtree := call.Node()

	// check if this is a call on a string method
	// then dispatch to specific folding method.
	stringNode := callTarget.Object()

	str, isStringLiteral := stringNode.AsStringLit()
	if !isStringLiteral && !p.useTypes {
		return subtree
	}
	methodID := callTarget.Node().GetStringID()
	firstArg, hasArgs := call.FirstArg()
	if isStringLiteral {
		if methodID == rkmIDSplit {
			return p.TryFoldStringSplit(call, str)
		} else if !hasArgs {
			switch methodID {
			case rkmIDToLowerCase:
				return p.TryFoldStringToLowerCase(call, str)
			case rkmIDToUpperCase:
				return p.TryFoldStringToUpperCase(call, str)
			case rkmIDTrim, rkmIDTrimStart, rkmIDTrimEnd, rkmIDTrimLeft, rkmIDTrimRight:
				return p.TryFoldStringTrim(call, str, callTarget.Property())
			case rkmIDAt:
				return p.TryFoldStringAt(call, str)
			default:
			}
		} else if NodeUtil_IsImmutableValue(firstArg.Node()) {
			switch methodID {
			case rkmIDIndexOf:
				return p.TryFoldStringIndexOf(call, "indexOf", str, firstArg)
			case rkmIDLastIndexOf:
				return p.TryFoldStringIndexOf(call, "lastIndexOf", str, firstArg)
			case rkmIDIncludes:
				return p.TryFoldStringIncludes(call, str, firstArg)
			case rkmIDStartsWith:
				return p.TryFoldStringStartsWith(call, str, firstArg)
			case rkmIDEndsWith:
				return p.TryFoldStringEndsWith(call, str, firstArg)
			case rkmIDAt:
				return p.TryFoldStringAt(call, str)
			case rkmIDSubstr:
				return p.TryFoldStringSubstr(call, str, firstArg)
			case rkmIDSubstring, rkmIDSlice:
				return p.TryFoldStringSubstringOrSlice(call, str, firstArg)
			case rkmIDCharAt:
				return p.TryFoldStringCharAt(call, str, firstArg)
			case rkmIDCharCodeAt:
				return p.TryFoldStringCharCodeAt(call, str, firstArg)
			case rkmIDReplace:
				return p.TryFoldStringReplace(call, str, firstArg)
			case rkmIDReplaceAll:
				return p.TryFoldStringReplaceAll(call, str, firstArg)
			default:
			}
		}
	}
	if !p.useTypes ||
		!hasArgs ||
		(methodID != rkmIDSubstr && methodID != rkmIDSubstring && methodID != rkmIDSlice) ||
		!(isStringLiteral || colors.StandardColors_STRING.Equals(stringNode.Node().GetColor())) {
		return subtree
	}
	// A call with exactly two arguments: the start, and the length or the end.
	secondArg, ok := firstArg.Next()
	if !ok {
		return subtree
	}
	if _, ok := secondArg.Next(); ok {
		return subtree
	}
	maybeStart, ok := p.GetSideEffectFreeNumberValue(firstArg.Node())
	if !ok {
		return subtree
	}
	start := int(jutil.D2I(maybeStart))
	maybeLengthOrEnd, ok := p.GetSideEffectFreeNumberValue(secondArg.Node())
	if !ok {
		return subtree
	}
	switch methodID {
	case rkmIDSubstr:
		length := int(jutil.D2I(maybeLengthOrEnd))
		if start >= 0 && length == 1 {
			return p.ReplaceWithCharAt(call, callTarget, firstArg, secondArg)
		}
	case rkmIDSubstring, rkmIDSlice:
		end := int(jutil.D2I(maybeLengthOrEnd))
		// unlike slice and substring, chatAt can not be used with negative indexes
		if start >= 0 && jutil.IntSub(end, start) == 1 {
			return p.ReplaceWithCharAt(call, callTarget, firstArg, secondArg)
		}
	}
	return subtree
}

// TryFoldStringToLowerCase returns The lowered string Node.
//
// This method is believed to be correct independent of the locale of the
// compiler and the JSVM executing the compiled code, assuming both are
// implementations of Unicode are correct.
//
// See https://tc39.es/ecma262/#sec-string.prototype.tolowercase,
// https://unicode.org/faq/casemap_charprop.html#5 and
// https://docs.oracle.com/javase/8/docs/api/java/lang/String.html#toLowerCase-java.util.Locale-
//
// Java: PeepholeReplaceKnownMethods.tryFoldStringToLowerCase(Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldStringToLowerCase(call ast.Call, stringNode ast.StringLit) rhino.Node {
	lowered := jutil.StringToLowerCaseRoot(stringNode.Value())
	replacement := rhino.IR_String(lowered)
	call.Node().ReplaceWith(replacement)
	p.ReportChangeToEnclosingScope(replacement)
	return replacement
}

// TryFoldStringToUpperCase returns The upped string Node.
//
// This method is believed to be correct independent of the locale of the
// compiler and the JSVM executing the compiled code, assuming both are
// implementations of Unicode are correct.
//
// See https://tc39.es/ecma262/#sec-string.prototype.touppercase,
// https://unicode.org/faq/casemap_charprop.html#5 and
// https://docs.oracle.com/javase/8/docs/api/java/lang/String.html#toUpperCase-java.util.Locale-
//
// Java: PeepholeReplaceKnownMethods.tryFoldStringToUpperCase(Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldStringToUpperCase(call ast.Call, stringNode ast.StringLit) rhino.Node {
	upped := stringToUpperCaseRoot(stringNode.Value())
	replacement := rhino.IR_String(upped)
	call.Node().ReplaceWith(replacement)
	p.ReportChangeToEnclosingScope(replacement)
	return replacement
}

// TryFoldStringTrim returns The trimmed string Node.
//
// Java: PeepholeReplaceKnownMethods.tryFoldStringTrim(Node, Node, String)
func (p *PeepholeReplaceKnownMethods) TryFoldStringTrim(call ast.Call, stringNode ast.StringLit, functionName string) rhino.Node {
	subtree := call.Node()
	// See ECMA-262 § 22.1.3.32 (trim), § 22.1.3.34 (trimStart), § 22.1.3.33 (trimEnd),
	// Annex B.2.2.15 (trimLeft), Annex B.2.2.16 (trimRight), § 12.2, and § 12.3
	// All Unicode whitespace + LineTerminator + BOM
	var result string
	switch functionName {
	case "trim":
		result = peepholeReplaceKnownMethods_TRIM_BOTH_PATTERN.replaceAll(stringNode.Value())
	case "trimStart", "trimLeft":
		result = peepholeReplaceKnownMethods_TRIM_START_PATTERN.replaceAll(stringNode.Value())
	case "trimEnd", "trimRight":
		result = peepholeReplaceKnownMethods_TRIM_END_PATTERN.replaceAll(stringNode.Value())
	default:
		return subtree
	}
	replacement := rhino.IR_String(result).Srcref(subtree)
	subtree.ReplaceWith(replacement)
	p.ReportChangeToEnclosingScope(replacement)
	return replacement
}

// TryFoldStringIndexOf tries to evaluate String.indexOf/lastIndexOf:
//
//	"abcdef".indexOf("bc") -> 1
//	"abcdefbc".indexOf("bc", 3) -> 6
//
// Java: PeepholeReplaceKnownMethods.tryFoldStringIndexOf(Node, String, Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldStringIndexOf(call ast.Call, functionName string, lstringNode ast.StringLit, firstArg ast.Arg) rhino.Node {
	n := call.Node()
	lstring := lstringNode.Value()
	isIndexOf := functionName == "indexOf"
	searchValue, ok := p.GetSideEffectFreeStringValue(firstArg.Node())
	// searchValue must be a valid string.
	if !ok {
		return n
	}
	fromIndex := 0
	if !isIndexOf {
		fromIndex = jutil.StringLength(lstring)
	}
	if secondArg, ok := firstArg.Next(); ok {
		// Third-argument and non-numeric second arg are problematic. Discard.
		from, isNumber := secondArg.AsNumber()
		if _, ok := secondArg.Next(); ok || !isNumber {
			return n
		}
		fromIndex = int(jutil.D2I(from.Value()))
	}
	var indexVal int
	if isIndexOf {
		indexVal = jutil.StringIndexOfFrom(lstring, searchValue, fromIndex)
	} else {
		indexVal = jutil.StringLastIndexOfFrom(lstring, searchValue, fromIndex)
	}
	newNode := NodeUtil_NumberNode(float64(indexVal), n)
	n.ReplaceWith(newNode)
	p.ReportChangeToEnclosingScope(newNode)

	return newNode
}

// TryFoldStringIncludes tries to evaluate String.prototype.includes:
//
//	"abcdef".includes("bc") -> true
//	"abcdef".includes("bc", 2) -> false
//
// Java: PeepholeReplaceKnownMethods.tryFoldStringIncludes(Node, Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldStringIncludes(call ast.Call, lstringNode ast.StringLit, firstArg ast.Arg) rhino.Node {
	n := call.Node()
	if _, ok := firstArg.AsRegExp(); ok {
		return n
	}

	lstring := lstringNode.Value()
	searchValue, ok := p.GetSideEffectFreeStringValue(firstArg.Node())
	if !ok {
		return n
	}

	fromIndex := 0
	if arg, ok := firstArg.Next(); ok {
		secondArg := arg.Node()
		if _, ok := arg.Next(); ok || p.MayHaveSideEffects(secondArg) {
			return n
		}
		if NodeUtil_IsNumericLiteral(secondArg) {
			posVal := NodeUtil_GetNumericLiteralValue(secondArg)
			fromIndex = int(jutil.D2I(base.MathMin(base.MathMax(posVal, 0), float64(jutil.StringLength(lstring)))))
		} else if !NodeUtil_IsNullOrUndefined(secondArg) {
			return n
		}
	}

	contains := jutil.StringIndexOfFrom(lstring, searchValue, fromIndex) != -1
	newNode := NodeUtil_BooleanNode(contains).Srcref(n)
	n.ReplaceWith(newNode)
	p.ReportChangeToEnclosingScope(newNode)
	return newNode
}

// TryFoldStringStartsWith tries to evaluate String.prototype.startsWith:
//
//	"abcdef".startsWith("abc") -> true
//	"abcdef".startsWith("bc", 1) -> true
//
// Java: PeepholeReplaceKnownMethods.tryFoldStringStartsWith(Node, Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldStringStartsWith(call ast.Call, lstringNode ast.StringLit, firstArg ast.Arg) rhino.Node {
	n := call.Node()
	if _, ok := firstArg.AsRegExp(); ok {
		return n
	}

	lstring := lstringNode.Value()
	searchValue, ok := p.GetSideEffectFreeStringValue(firstArg.Node())
	if !ok {
		return n
	}

	start := 0
	if arg, ok := firstArg.Next(); ok {
		secondArg := arg.Node()
		if _, ok := arg.Next(); ok || p.MayHaveSideEffects(secondArg) {
			return n
		}
		if NodeUtil_IsNumericLiteral(secondArg) {
			posVal := NodeUtil_GetNumericLiteralValue(secondArg)
			start = int(jutil.D2I(base.MathMin(base.MathMax(posVal, 0), float64(jutil.StringLength(lstring)))))
		} else if !NodeUtil_IsNullOrUndefined(secondArg) {
			return n
		}
	}

	// String.startsWith(String, int)
	result := jutil.StringRegionMatches(lstring, start, searchValue, 0, jutil.StringLength(searchValue))
	newNode := NodeUtil_BooleanNode(result).Srcref(n)
	n.ReplaceWith(newNode)
	p.ReportChangeToEnclosingScope(newNode)
	return newNode
}

// TryFoldStringEndsWith tries to evaluate String.prototype.endsWith:
//
//	"abcdef".endsWith("def") -> true
//	"abcdef".endsWith("de", 5) -> true
//
// Java: PeepholeReplaceKnownMethods.tryFoldStringEndsWith(Node, Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldStringEndsWith(call ast.Call, lstringNode ast.StringLit, firstArg ast.Arg) rhino.Node {
	n := call.Node()
	if _, ok := firstArg.AsRegExp(); ok {
		return n
	}

	lstring := lstringNode.Value()
	searchValue, ok := p.GetSideEffectFreeStringValue(firstArg.Node())
	if !ok {
		return n
	}

	end := jutil.StringLength(lstring)
	if arg, ok := firstArg.Next(); ok {
		secondArg := arg.Node()
		if _, ok := arg.Next(); ok || p.MayHaveSideEffects(secondArg) {
			return n
		}
		if NodeUtil_IsNumericLiteral(secondArg) {
			posVal := NodeUtil_GetNumericLiteralValue(secondArg)
			end = int(jutil.D2I(base.MathMin(base.MathMax(posVal, 0), float64(jutil.StringLength(lstring)))))
		} else if secondArg.IsNull() {
			end = 0
		} else if NodeUtil_IsUndefined(secondArg) {
			end = jutil.StringLength(lstring)
		} else {
			return n
		}
	}

	searchLength := jutil.StringLength(searchValue)
	start := end - searchLength
	// String.startsWith(String, int)
	result := start >= 0 && jutil.StringRegionMatches(lstring, start, searchValue, 0, searchLength)
	newNode := NodeUtil_BooleanNode(result).Srcref(n)
	n.ReplaceWith(newNode)
	p.ReportChangeToEnclosingScope(newNode)
	return newNode
}

// TryFoldStringAt tries to evaluate String.prototype.at:
//
//	"abcdef".at(0) -> "a"
//	"abcdef".at(-1) -> "f"
//	"abcdef".at(10) -> void 0
//
// Java: PeepholeReplaceKnownMethods.tryFoldStringAt(Node, Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldStringAt(call ast.Call, stringNode ast.StringLit) rhino.Node {
	n := call.Node()
	var posVal float64
	firstArg, ok := call.FirstArg()
	if !ok {
		posVal = 0.0
	} else {
		if _, ok := firstArg.Next(); ok {
			return n
		}
		maybePos, ok := p.GetSideEffectFreeNumberValue(firstArg.Node())
		if !ok {
			return n
		}
		if math.IsNaN(maybePos) {
			posVal = 0.0
		} else {
			posVal = maybePos
		}
	}

	stringVal := stringNode.Value()
	len_ := jutil.StringLength(stringVal)

	var replacement rhino.Node
	if math.IsInf(posVal, 0) || posVal >= float64(len_) || posVal <= float64(-(len_+1)) {
		replacement = NodeUtil_NewUndefinedNode(n)
	} else {
		intPos := jutil.D2L(posVal)
		var k int64
		if intPos >= 0 {
			k = intPos
		} else {
			k = int64(len_) + intPos
		}
		if k < 0 || k >= int64(len_) {
			replacement = NodeUtil_NewUndefinedNode(n)
		} else if jutil.IsASCII(stringVal) {
			replacement = rhino.IR_String(stringVal[k : k+1]).Srcref(n)
		} else {
			replacement = rhino.IR_String(jutil.StringFromChar(jutil.StringCharAt(stringVal, int(k)))).Srcref(n)
		}
	}

	n.ReplaceWith(replacement)
	p.ReportChangeToEnclosingScope(replacement)
	return replacement
}

// TryFoldStringSubstr tries to fold .substr() calls on strings.
//
// Java: PeepholeReplaceKnownMethods.tryFoldStringSubstr(Node, Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldStringSubstr(call ast.Call, stringNode ast.StringLit, arg1 ast.Arg) rhino.Node {
	n := call.Node()

	var start int
	var length int
	stringAsString := stringNode.Value()

	maybeStart, ok := p.GetSideEffectFreeNumberValue(arg1.Node())
	if ok {
		start = int(jutil.D2I(maybeStart))
	} else {
		return n
	}

	if arg2, ok := arg1.Next(); ok {
		maybeLength, ok := p.GetSideEffectFreeNumberValue(arg2.Node())
		if ok {
			length = int(jutil.D2I(maybeLength))
		} else {
			return n
		}

		if _, ok := arg2.Next(); ok {
			// If we got more args than we expected, bail out.
			return n
		}
	} else {
		// parameter 2 not passed
		length = jutil.IntSub(jutil.StringLength(stringAsString), start)
	}

	// Don't handle these cases. The specification actually does
	// specify the behavior in some of these cases, but we haven't
	// done a thorough investigation that it is correctly implemented
	// in all browsers.
	if jutil.IntAdd(start, length) > jutil.StringLength(stringAsString) || (length < 0) || (start < 0) {
		return n
	}

	result := jutil.StringSubstring(stringAsString, start, jutil.IntAdd(start, length))
	resultNode := rhino.IR_String(result)

	parent := n.GetParent()
	n.ReplaceWith(resultNode)
	p.ReportChangeToEnclosingScope(parent)
	return resultNode
}

// TryFoldStringSubstringOrSlice tries to fold .substring() or .slice() calls
// on strings.
//
// Java: PeepholeReplaceKnownMethods.tryFoldStringSubstringOrSlice(Node, Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldStringSubstringOrSlice(call ast.Call, stringNode ast.StringLit, arg1 ast.Arg) rhino.Node {
	n := call.Node()

	var start int
	var end int
	stringAsString := stringNode.Value()

	maybeStart, ok := p.GetSideEffectFreeNumberValue(arg1.Node())
	if ok {
		start = int(jutil.D2I(maybeStart))
	} else {
		return n
	}

	if arg2, ok := arg1.Next(); ok {
		maybeEnd, ok := p.GetSideEffectFreeNumberValue(arg2.Node())
		if ok {
			end = int(jutil.D2I(maybeEnd))
		} else {
			return n
		}

		if _, ok := arg2.Next(); ok {
			// If we got more args than we expected, bail out.
			return n
		}
	} else {
		// parameter 2 not passed
		end = jutil.StringLength(stringAsString)
	}

	// Don't handle these cases. The specification actually does
	// specify the behavior in some of these cases, but we haven't
	// done a thorough investigation that it is correctly implemented
	// in all browsers.
	if (end > jutil.StringLength(stringAsString)) ||
		(start > jutil.StringLength(stringAsString)) ||
		(start < 0) ||
		(end < 0) ||
		(start > end) {
		return n
	}

	result := jutil.StringSubstring(stringAsString, start, end)
	resultNode := rhino.IR_String(result)

	parent := n.GetParent()
	n.ReplaceWith(resultNode)
	p.ReportChangeToEnclosingScope(parent)
	return resultNode
}

// ReplaceWithCharAt rewrites s.substr(start, 1), s.substring(start, start
// + 1) or s.slice(start, start + 1) in place to s.charAt(start): it renames
// the method and removes secondArg.
//
// Go-only: it takes the second argument (Java: firstArg.getNext()).
//
// Java: PeepholeReplaceKnownMethods.replaceWithCharAt(Node, Node, Node)
func (p *PeepholeReplaceKnownMethods) ReplaceWithCharAt(call ast.Call, callTarget ast.GetProp, firstArg ast.Arg, secondArg ast.Arg) rhino.Node {
	// TODO(moz): Maybe correct the arity of the function type here.
	callTarget.SetProperty("charAt")
	secondArg.Detach()
	p.ReportChangeToEnclosingScope(firstArg.Node())
	return call.Node()
}

// TryFoldStringCharAt tries to fold .charAt() calls on strings.
//
// Java: PeepholeReplaceKnownMethods.tryFoldStringCharAt(Node, Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldStringCharAt(call ast.Call, stringNode ast.StringLit, arg1 ast.Arg) rhino.Node {
	n := call.Node()

	var index int
	stringAsString := stringNode.Value()

	// The only argument, a number.
	num, isNumber := arg1.AsNumber()
	if _, more := arg1.Next(); isNumber && !more {
		index = int(jutil.D2I(num.Value()))
	} else {
		return n
	}

	if index < 0 || jutil.StringLength(stringAsString) <= index {
		// http://es5.github.com/#x15.5.4.4 says "" is returned when index is
		// out of bounds but we bail.
		return n
	}

	resultNode := rhino.IR_String(jutil.StringSubstring(stringAsString, index, index+1))
	parent := n.GetParent()
	n.ReplaceWith(resultNode)
	p.ReportChangeToEnclosingScope(parent)
	return resultNode
}

// TryFoldStringReplace tries to fold .replace() calls on strings.
//
// Java: PeepholeReplaceKnownMethods.tryFoldStringReplace(Node, Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldStringReplace(call ast.Call, stringNode ast.StringLit, arg1 ast.Arg) rhino.Node {
	return p.tryFoldStringReplaceImpl(call, stringNode, arg1, false)
}

// TryFoldStringReplaceAll tries to fold .replaceAll() calls on strings.
//
// Java: PeepholeReplaceKnownMethods.tryFoldStringReplaceAll(Node, Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldStringReplaceAll(call ast.Call, stringNode ast.StringLit, arg1 ast.Arg) rhino.Node {
	return p.tryFoldStringReplaceImpl(call, stringNode, arg1, true)
}

func (p *PeepholeReplaceKnownMethods) tryFoldStringReplaceImpl(call ast.Call, stringNode ast.StringLit, arg1 ast.Arg, replaceAll bool) rhino.Node {
	n := call.Node()
	arg2, ok := arg1.Next()
	if !ok {
		return n
	}
	if _, ok := arg2.Next(); ok {
		return n
	}

	pattern, isString1 := arg1.AsStringLit()
	replacement, isString2 := arg2.AsStringLit()
	if !isString1 || !isString2 {
		// only string literals are supported for folding.
		return n
	}

	lookForPattern := pattern.Value()
	replacementPattern := replacement.Value()
	if strings.IndexByte(replacementPattern, '$') != -1 {
		// 'special' replacements aren't supported yet.
		return n
	}

	original := stringNode.Value()
	var newString string
	if replaceAll {
		if jutil.IsASCII(original) && jutil.IsASCII(lookForPattern) {
			newString = strings.ReplaceAll(original, lookForPattern, replacementPattern)
		} else {
			newString = jutil.StringReplace(original, lookForPattern, replacementPattern)
		}
	} else if jutil.IsASCII(original) && jutil.IsASCII(lookForPattern) {
		index := strings.Index(original, lookForPattern)
		if index == -1 {
			return n
		}
		newString = original[:index] + replacementPattern + original[index+len(lookForPattern):]
	} else {
		index := jutil.StringIndexOf(original, lookForPattern)
		if index == -1 {
			return n
		}
		newString =
			jutil.StringSubstring(original, 0, index) +
				replacementPattern +
				jutil.StringSubstringFrom(original, index+jutil.StringLength(lookForPattern))
	}

	resultNode := rhino.IR_String(newString).Srcref(stringNode.Node())
	parent := n.GetParent()
	n.ReplaceWith(resultNode)
	p.ReportChangeToEnclosingScope(parent)
	return resultNode
}

// TryFoldStringCharCodeAt tries to fold .charCodeAt() calls on strings.
//
// Java: PeepholeReplaceKnownMethods.tryFoldStringCharCodeAt(Node, Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldStringCharCodeAt(call ast.Call, stringNode ast.StringLit, arg1 ast.Arg) rhino.Node {
	n := call.Node()

	var index int
	stringAsString := stringNode.Value()

	// The only argument, a number.
	num, isNumber := arg1.AsNumber()
	if _, more := arg1.Next(); isNumber && !more {
		index = int(jutil.D2I(num.Value()))
	} else {
		return n
	}

	if index < 0 || jutil.StringLength(stringAsString) <= index {
		// http://es5.github.com/#x15.5.4.5 says NaN is returned when index is
		// out of bounds but we bail.
		return n
	}

	resultNode := rhino.IR_Number(float64(jutil.StringCharAt(stringAsString, index)))
	parent := n.GetParent()
	n.ReplaceWith(resultNode)
	p.ReportChangeToEnclosingScope(parent)
	return resultNode
}

// JsSplit implements the JS String.split method using a string separator.
//
// Java: PeepholeReplaceKnownMethods.jsSplit(String, String, int)
func (p *PeepholeReplaceKnownMethods) JsSplit(stringValue string, separator string, hasSeparator bool, limit int) []string {
	jutil.CheckArgument(limit >= 0)

	if limit == 0 {
		return []string{}
	}
	if !hasSeparator {
		return []string{stringValue}
	}

	var splitStrings []string
	if jutil.IsASCII(stringValue) && jutil.IsASCII(separator) {
		if separator == "" {
			n := min(len(stringValue), limit)
			splitStrings = make([]string, n)
			for i := 0; i < n; i++ {
				splitStrings[i] = stringValue[i : i+1]
			}
			return splitStrings
		}
		startIndex := 0
		for len(splitStrings) < limit && startIndex+len(separator) <= len(stringValue) {
			rel := strings.Index(stringValue[startIndex:], separator)
			if rel < 0 {
				break
			}
			matchIndex := startIndex + rel
			splitStrings = append(splitStrings, stringValue[startIndex:matchIndex])
			startIndex = matchIndex + len(separator)
		}
		if len(splitStrings) < limit {
			splitStrings = append(splitStrings, stringValue[startIndex:])
		}
		return splitStrings
	}

	strLen := jutil.StringLength(stringValue)
	if separator == "" {
		for i := 0; i < strLen && i < limit; i++ {
			splitStrings = append(splitStrings, jutil.StringSubstring(stringValue, i, i+1))
		}
		return splitStrings
	}

	sepLen := jutil.StringLength(separator)
	startIndex := 0
	for len(splitStrings) < limit && startIndex+sepLen <= strLen {
		matchIndex := jutil.StringIndexOfFrom(stringValue, separator, startIndex)
		if matchIndex < 0 {
			break
		}
		splitStrings = append(splitStrings, jutil.StringSubstring(stringValue, startIndex, matchIndex))
		startIndex = matchIndex + sepLen
	}
	if len(splitStrings) < limit {
		if startIndex < strLen {
			splitStrings = append(splitStrings, jutil.StringSubstringFrom(stringValue, startIndex))
		} else {
			splitStrings = append(splitStrings, "")
		}
	}
	return splitStrings
}

// TryFoldStringSplit tries to fold .split() calls on strings.
//
// Java: PeepholeReplaceKnownMethods.tryFoldStringSplit(Node, Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldStringSplit(call ast.Call, stringNode ast.StringLit) rhino.Node {
	n := call.Node()
	if p.late {
		return n
	}

	var separator string
	hasSeparator := false
	stringValue := stringNode.Value()

	// Maximum number of possible splits
	limit := jutil.StringLength(stringValue) + 1

	if arg1, ok := call.FirstArg(); ok {
		if sep, ok := arg1.AsStringLit(); ok {
			separator, hasSeparator = sep.Value(), true
		} else if arg1.Token() != rhino.Token_NULL {
			return n
		}

		// Go-only note: like the Java, this does not look past the limit, so
		// "a,b".split(",", 1, f()) folds to ["a"] and drops the call to f.
		// Kept for identical output.
		if arg2, ok := arg1.Next(); ok {
			if NodeUtil_IsNumericLiteral(arg2.Node()) {
				// Integer.toUnsignedLong(int)
				uint32Limit :=
					int64(uint32(base.JSCompDoubles_EcmascriptToUint32(NodeUtil_GetNumericLiteralValue(arg2.Node()))))
				limit = int(min(uint32Limit, int64(limit)))
			} else {
				return n
			}
		}
	}

	// Split the string and convert the returned array into JS nodes
	stringArray := p.JsSplit(stringValue, separator, hasSeparator, limit)
	arrayOfStrings := rhino.IR_Arraylit()
	for _, element := range stringArray {
		arrayOfStrings.AddChildToBack(rhino.IR_String(element).Srcref(stringNode.Node()))
	}

	parent := n.GetParent()
	n.ReplaceWith(arrayOfStrings)
	p.ReportChangeToEnclosingScope(parent)
	return arrayOfStrings
}
