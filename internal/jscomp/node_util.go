// Port of com/google/javascript/jscomp/NodeUtil.java @ closure-compiler v20261021

/*
 * Copyright 2004 The Closure Compiler Authors.
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
	"iter"
	"math"
	"math/big"
	"slices"
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/dtoa"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// NodeUtil contains generally useful AST utilities.
type NodeUtil struct{}

// Value of JavaScript's Number.MAX_SAFE_INTEGER
//
// Java: NodeUtil.MAX_POSITIVE_INTEGER_NUMBER
const NodeUtil_MAX_POSITIVE_INTEGER_NUMBER int64 = (1 << 53) - 1

// Java: NodeUtil.JSC_PROPERTY_NAME_FN
const NodeUtil_JSC_PROPERTY_NAME_FN = "JSCompiler_renameProperty"

// Java: NodeUtil.LARGEST_BASIC_LATIN
const NodeUtil_LARGEST_BASIC_LATIN uint16 = 0x7f

// Java: NodeUtil.GOOG_MODULE_DECLARE_LEGACY_NAMESPACE
var NodeUtil_GOOG_MODULE_DECLARE_LEGACY_NAMESPACE = rhino.QualifiedName_Of("goog.module.declareLegacyNamespace")

// Java: NodeUtil.GOOG_SET_TEST_ONLY
var NodeUtil_GOOG_SET_TEST_ONLY = rhino.QualifiedName_Of("goog.setTestOnly")

// Java: NodeUtil.GOOG_PROVIDE
var NodeUtil_GOOG_PROVIDE = rhino.QualifiedName_Of("goog.provide")

// Java: NodeUtil.GOOG_MODULE
var NodeUtil_GOOG_MODULE = rhino.QualifiedName_Of("goog.module")

// Java: NodeUtil.GOOG_MODULE_GET
var NodeUtil_GOOG_MODULE_GET = rhino.QualifiedName_Of("goog.module.get")

// Java: NodeUtil.GOOG_REQUIRE
var NodeUtil_GOOG_REQUIRE = rhino.QualifiedName_Of("goog.require")

// Java: NodeUtil.GOOG_REQUIRE_TYPE
var NodeUtil_GOOG_REQUIRE_TYPE = rhino.QualifiedName_Of("goog.requireType")

// Java: NodeUtil.GOOG_FORWARD_DECLARE
var NodeUtil_GOOG_FORWARD_DECLARE = rhino.QualifiedName_Of("goog.forwardDeclare")

// Java: NodeUtil.GOOG_REQUIRE_DYNAMIC
var NodeUtil_GOOG_REQUIRE_DYNAMIC = rhino.QualifiedName_Of("goog.requireDynamic")

// Java: NodeUtil.GOOG_WEAK_USAGE
var NodeUtil_GOOG_WEAK_USAGE = rhino.QualifiedName_Of("goog.weakUsage")

// Java: NodeUtil.GOOG_WEAK_USAGE_MANGLED
var NodeUtil_GOOG_WEAK_USAGE_MANGLED = rhino.QualifiedName_Of("goog$weakUsage")

// Utility class; do not instantiate.
//
// Java: NodeUtil()
func NewNodeUtil() *NodeUtil {
	return &NodeUtil{}
}

// NodeUtil_GetBooleanValue gets the boolean value of a node that represents
// an expression, or Tri.UNKNOWN if no such value can be determined by static
// analysis.
//
// This method does not consider whether the node may have side-effects.
//
// Java: NodeUtil.getBooleanValue(Node)
func NodeUtil_GetBooleanValue(n rhino.Node) base.Tri {
	// This switch consists of cases that are not supported by getLiteralBooleanValue(),
	// which we will call if none of these match.
	switch n.GetToken() {
	case rhino.Token_NULL, rhino.Token_FALSE, rhino.Token_VOID:
		return base.Tri_FALSE
	case rhino.Token_TRUE, rhino.Token_REGEXP, rhino.Token_FUNCTION, rhino.Token_CLASS, rhino.Token_NEW,
		rhino.Token_ARRAYLIT, rhino.Token_OBJECTLIT:
		return base.Tri_TRUE

	case rhino.Token_TEMPLATELIT:
		if n.HasOneChild() {
			templateLitString := n.GetOnlyChild()
			jutil.CheckStateMsg1(templateLitString.IsTemplateLitString(), "%s", templateLitString)
			cookedString := templateLitString.GetCookedString()
			return base.Tri_ForBoolean(cookedString != nil && *cookedString != "")
		} else {
			return base.Tri_UNKNOWN
		}
	case rhino.Token_STRINGLIT:
		return base.Tri_ForBoolean(jutil.StringLength(n.GetString()) > 0)
	case rhino.Token_NUMBER:
		return base.Tri_ForBoolean(n.GetDouble() != 0)
	case rhino.Token_BIGINT:
		return base.Tri_ForBoolean(n.GetBigInt().Sign() != 0)
	case rhino.Token_NOT:
		return NodeUtil_GetBooleanValue(n.GetLastChild()).Not()

	case rhino.Token_NAME:
		// We assume here that programs don't change the value of these global variables.
		switch n.GetString() {
		case "undefined", "NaN":
			return base.Tri_FALSE
		case "Infinity":
			return base.Tri_TRUE
		default:
			return base.Tri_UNKNOWN
		}

	case rhino.Token_BITNOT, rhino.Token_POS, rhino.Token_NEG:
		doubleVal := NodeUtil_GetNumberValue(n)
		if doubleVal != nil {
			isFalsey := math.IsNaN(*doubleVal) || base.JSCompDoubles_IsEitherZero(*doubleVal)
			return base.Tri_ForBoolean(!isFalsey)
		}

		bigintVal := NodeUtil_GetBigIntValue(n)
		if bigintVal != nil {
			isFalsey := bigintVal.Sign() == 0
			return base.Tri_ForBoolean(!isFalsey)
		}

		return base.Tri_UNKNOWN

	// For ASSIGN and COMMA the value is the value of the RHS.
	case rhino.Token_ASSIGN, rhino.Token_COMMA:
		return NodeUtil_GetBooleanValue(n.GetLastChild())
	case rhino.Token_AND, rhino.Token_ASSIGN_AND:
		lhs := NodeUtil_GetBooleanValue(n.GetFirstChild())
		rhs := NodeUtil_GetBooleanValue(n.GetLastChild())
		return lhs.And(rhs)
	case rhino.Token_OR, rhino.Token_ASSIGN_OR:
		lhs := NodeUtil_GetBooleanValue(n.GetFirstChild())
		rhs := NodeUtil_GetBooleanValue(n.GetLastChild())
		return lhs.Or(rhs)
	case rhino.Token_HOOK:
		trueValue := NodeUtil_GetBooleanValue(n.GetSecondChild())
		falseValue := NodeUtil_GetBooleanValue(n.GetLastChild())
		if trueValue == falseValue {
			return trueValue
		} else {
			return base.Tri_UNKNOWN
		}
	case rhino.Token_COALESCE, rhino.Token_ASSIGN_COALESCE:
		lhs := NodeUtil_GetBooleanValue(n.GetFirstChild())
		rhs := NodeUtil_GetBooleanValue(n.GetLastChild())
		if lhs == base.Tri_TRUE || lhs == rhs {
			return lhs
		} else {
			return base.Tri_UNKNOWN
		}
	default:
		return base.Tri_UNKNOWN
	}
}

// NodeUtil_GetStringValue gets the value of a node as a String, or null if it
// cannot be converted. When it returns a non-null String, this method
// effectively emulates the String() JavaScript cast function.
//
// IMPORTANT: This method does not consider whether n may have side effects.
//
// Java: NodeUtil.getStringValue(Node)
func NodeUtil_GetStringValue(n rhino.Node) *string {
	// TODO(user): regex literals as well.
	switch n.GetToken() {
	case rhino.Token_STRINGLIT, rhino.Token_STRING_KEY:
		return jutil.StringPtr(n.GetString())
	case rhino.Token_TEMPLATELIT:
		// Only convert a template literal if all its expressions can be converted.
		string_ := jutil.NewStringBuilder()
		for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
			expression := child
			if child.IsTemplateLitSub() {
				expression = child.GetFirstChild()
			}
			expressionString := NodeUtil_GetStringValue(expression)
			if expressionString == nil {
				// Cannot convert.
				return nil
			}
			string_.Append(*expressionString)
		}
		return jutil.StringPtr(string_.String())
	case rhino.Token_TEMPLATELIT_STRING:
		return n.GetCookedString()
	case rhino.Token_NAME:
		name := n.GetString()
		if "undefined" == name || "Infinity" == name || "NaN" == name {
			return jutil.StringPtr(name)
		}
		return nil
	case rhino.Token_NEG, rhino.Token_NUMBER:
		value := NodeUtil_GetNumberValue(n)
		if value != nil {
			return jutil.StringPtr(dtoa.DToA_NumberToString(*value))
		}
		return nil
	case rhino.Token_BIGINT:
		return jutil.StringPtr(n.GetBigInt().String())
	case rhino.Token_FALSE:
		return jutil.StringPtr("false")
	case rhino.Token_TRUE:
		return jutil.StringPtr("true")
	case rhino.Token_NULL:
		return jutil.StringPtr("null")
	case rhino.Token_VOID:
		return jutil.StringPtr("undefined")
	case rhino.Token_NOT:
		child := NodeUtil_GetBooleanValue(n.GetFirstChild())
		if child != base.Tri_UNKNOWN {
			if child.ToBoolean(true) { // reversed.
				return jutil.StringPtr("false")
			}
			return jutil.StringPtr("true")
		}
		return nil
	case rhino.Token_ARRAYLIT:
		return NodeUtil_ArrayToString(n)
	case rhino.Token_OBJECTLIT:
		if n.HasChildren() {
			return nil
		}
		return jutil.StringPtr("[object Object]")
	default:
		return nil
	}
}

// NodeUtil_GetArrayElementStringValue: when converting arrays to string using
// Array.prototype.toString or Array.prototype.join, the rules for conversion
// to String are different than converting each element individually.
// Specifically, "null" and "undefined" are converted to an empty string.
//
// n is a node that is a member of an Array; returns the string
// representation.
//
// Java: NodeUtil.getArrayElementStringValue(Node)
func NodeUtil_GetArrayElementStringValue(n rhino.Node) *string {
	if NodeUtil_IsNullOrUndefined(n) || n.IsEmpty() {
		return jutil.StringPtr("")
	}
	return NodeUtil_GetStringValue(n)
}

// Java: NodeUtil.arrayToString(Node)
func NodeUtil_ArrayToString(literal rhino.Node) *string {
	first := literal.GetFirstChild()
	var result strings.Builder
	for n := first; !n.IsNil(); n = n.GetNext() {
		childValue := NodeUtil_GetArrayElementStringValue(n)
		if childValue == nil {
			return nil
		}
		if n != first {
			result.WriteByte(',')
		}
		result.WriteString(*childValue)
	}
	return jutil.StringPtr(result.String())
}

// NodeUtil_GetNumberValue gets the value of a node as a Number, or null if it
// cannot be converted. When it returns a non-null Double, this method
// effectively emulates the Number() JavaScript cast function.
//
// IMPORTANT: This method does not consider whether n may have side effects.
//
// Returns the value of a node as a Number, or null if it cannot be converted.
//
// Java: NodeUtil.getNumberValue(Node)
func NodeUtil_GetNumberValue(n rhino.Node) *float64 {
	v, ok := nodeUtil_doGetNumberValue(n, true /* numberConversions */)
	if !ok {
		return nil
	}
	return &v
}

// NodeUtil_GetNumberValueNoConversions gets the value of a node as a Number,
// or null if it isn't a number. This method is similar to
// `typeof v == "number" ? v : null` as it does not perform any automatic
// conversions.
//
// IMPORTANT: This method does not consider whether n may have side effects.
//
// Java: NodeUtil.getNumberValueNoConversions(Node)
func NodeUtil_GetNumberValueNoConversions(n rhino.Node) *float64 {
	v, ok := nodeUtil_doGetNumberValue(n, false /* numberConversions */)
	if !ok {
		return nil
	}
	return &v
}

// nodeUtil_double boxes a double (Java's autoboxing to Double).
func nodeUtil_double(d float64) *float64 {
	return &d
}

// Java: NodeUtil.doGetNumberValue(Node, boolean)
func NodeUtil_DoGetNumberValue(n rhino.Node, numberConversions bool) *float64 {
	v, ok := nodeUtil_doGetNumberValue(n, numberConversions)
	if !ok {
		return nil
	}
	return &v
}

func nodeUtil_doGetNumberValue(n rhino.Node, numberConversions bool) (float64, bool) {
	switch n.GetToken() {
	case rhino.Token_NUMBER:
		return n.GetDouble(), true
	case rhino.Token_BIGINT:
		// When this call returns non-null, it is an assertion that JavaScript automatic conversion
		// to Number (e.g. during arithmetic operations) would produce the value returned here.
		// The spec does not allow automatic conversion from BigInt to Number, since that would
		// likely result in incorrect computation results.
		return 0, false
	case rhino.Token_VOID:
		if !numberConversions {
			return 0, false
		}
		return math.NaN(), true
	case rhino.Token_NAME:
		switch n.GetString() {
		case "undefined":
			if numberConversions {
				return math.NaN(), true
			}
			return 0, false
		case "NaN":
			return math.NaN(), true
		case "Infinity":
			return math.Inf(1), true
		default:
			return 0, false
		}
	case rhino.Token_POS:
		// unary plus triggers numeric conversions only when numberConversions is true
		return nodeUtil_doGetNumberValue(n.GetOnlyChild(), numberConversions)
	case rhino.Token_NEG:
		val, ok := nodeUtil_doGetNumberValue(n.GetOnlyChild(), numberConversions)
		if !ok {
			return 0, false
		}
		return -val, true
	case rhino.Token_BITNOT:
		val, ok := nodeUtil_doGetNumberValue(n.GetOnlyChild(), numberConversions)
		if !ok {
			return 0, false
		}
		return float64(^base.JSCompDoubles_EcmascriptToInt32(val)), true
	case rhino.Token_FALSE, rhino.Token_NOT, rhino.Token_NULL, rhino.Token_TRUE:
		if !numberConversions {
			return 0, false
		}
		switch NodeUtil_GetBooleanValue(n) {
		case base.Tri_TRUE:
			return 1.0, true
		case base.Tri_FALSE:
			return 0.0, true
		default: // case UNKNOWN
			return 0, false
		}
	case rhino.Token_TEMPLATELIT:
		if !numberConversions {
			return 0, false
		}
		string_ := NodeUtil_GetStringValue(n)
		if string_ == nil {
			return 0, false
		}
		return nodeUtil_getStringNumberValue(*string_)
	case rhino.Token_STRINGLIT:
		if !numberConversions {
			return 0, false
		}
		return nodeUtil_getStringNumberValue(n.GetString())
	case rhino.Token_ARRAYLIT, rhino.Token_OBJECTLIT:
		if !numberConversions {
			return 0, false
		}
		value := NodeUtil_GetStringValue(n)
		if value != nil {
			return nodeUtil_getStringNumberValue(*value)
		}
		return 0, false
	default:
	}

	return 0, false
}

// Java: NodeUtil.getStringNumberValue(String)
func NodeUtil_GetStringNumberValue(rawJsString string) *float64 {
	v, ok := nodeUtil_getStringNumberValue(rawJsString)
	if !ok {
		return nil
	}
	return &v
}

func nodeUtil_getStringNumberValue(rawJsString string) (float64, bool) {
	if strings.Contains(rawJsString, "\u000b") {
		// vertical tab is not always whitespace
		return 0, false
	}

	s := NodeUtil_TrimJsWhiteSpace(rawJsString)
	// return ScriptRuntime.toNumber(s);
	if s == "" {
		return 0.0, true
	}

	if jutil.StringLength(s) >= 2 && jutil.StringCharAt(s, 0) == '0' {
		radix := 0
		prefix := jutil.StringCharAt(s, 1)
		if prefix == 'x' || prefix == 'X' {
			radix = 16
		} else if prefix == 'b' || prefix == 'B' {
			radix = 2
		} else if prefix == 'o' || prefix == 'O' {
			radix = 8
		}
		if radix != 0 {
			digits := jutil.StringSubstringFrom(s, 2)
			if digits == "" {
				return math.NaN(), true
			}
			for i := 0; i < jutil.StringLength(digits); i++ {
				if !nodeUtil_isValidDigit(jutil.StringCharAt(digits, i), radix) {
					return math.NaN(), true
				}
			}
			var result float64
			jutil.Try(func() {
				if jutil.StringLength(digits) <= 14 {
					result = float64(jutil.LongParseLong(digits, radix))
				} else {
					result = nodeUtil_bigIntegerDoubleValue(digits, radix)
				}
			}, func(e jutil.INumberFormatException) {
				result = math.NaN()
			})
			return result, true
		}
	}

	if jutil.StringLength(s) > 3 &&
		(jutil.StringCharAt(s, 0) == '-' || jutil.StringCharAt(s, 0) == '+') &&
		jutil.StringCharAt(s, 1) == '0' &&
		(jutil.StringCharAt(s, 2) == 'x' || jutil.StringCharAt(s, 2) == 'X') {
		// hex numbers with explicit signs vary between browsers.
		return 0, false
	}

	// Firefox and IE treat the "Infinity" differently. Firefox is case
	// insensitive, but IE treats "infinity" as NaN.  So leave it alone.
	if s == "infinity" || s == "-infinity" || s == "+infinity" {
		return 0, false
	}

	if s == "Infinity" || s == "+Infinity" {
		return math.Inf(1), true
	}

	if s == "-Infinity" {
		return math.Inf(-1), true
	}

	for i := 0; i < len(s); i++ {
		// Go-only: language-adapter — a byte loop over the WTF-8 string: every
		// accepted char is ASCII, and a non-ASCII char's first byte is not one.
		c := s[i]
		if ('0' <= c && c <= '9') || c == '.' || c == 'e' || c == 'E' || c == '+' || c == '-' {
			continue
		}
		return math.NaN(), true
	}

	var result float64
	jutil.Try(func() {
		result = jutil.DoubleParseDouble(s)
	}, func(e jutil.INumberFormatException) {
		result = math.NaN()
	})
	return result, true
}

// Java: NodeUtil.isValidDigit(char, int)
func nodeUtil_isValidDigit(c uint16, radix int) bool {
	switch radix {
	case 2:
		return c == '0' || c == '1'
	case 8:
		return '0' <= c && c <= '7'
	case 16:
		return ('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
	default:
		return false
	}
}

// nodeUtil_bigIntegerDoubleValue is new BigInteger(digits, radix).doubleValue()
// for digits that are all valid in radix: the double nearest to the value,
// ties to even, or Infinity if it is too large.
//
// Go-only: language-adapter — java.math.BigInteger.
func nodeUtil_bigIntegerDoubleValue(digits string, radix int) float64 {
	x, ok := new(big.Int).SetString(digits, radix)
	if !ok {
		panic(jutil.NewNumberFormatException("For input string: \"" + digits + "\""))
	}
	f, _ := new(big.Float).SetInt(x).Float64()
	return f
}

const (
	nodeUtil_smallBigIntMin = -128
	nodeUtil_smallBigIntMax = 383
)

var nodeUtil_smallBigInts = func() [nodeUtil_smallBigIntMax - nodeUtil_smallBigIntMin + 1]big.Int {
	var arr [nodeUtil_smallBigIntMax - nodeUtil_smallBigIntMin + 1]big.Int
	for i := range arr {
		arr[i].SetInt64(int64(i + nodeUtil_smallBigIntMin))
	}
	return arr
}()

func nodeUtil_bigIntFromInt64(v int64) *big.Int {
	if v >= nodeUtil_smallBigIntMin && v <= nodeUtil_smallBigIntMax {
		return &nodeUtil_smallBigInts[v-nodeUtil_smallBigIntMin]
	}
	return big.NewInt(v)
}

// NodeUtil_GetBigIntValue gets the value of a node as a BigInt, or null if it
// cannot be converted. When it returns a non-null BigInteger, this method
// effectively emulates the BigInt() JavaScript cast function.
//
// IMPORTANT: This method does not consider whether n may have side effects.
//
// Java: NodeUtil.getBigIntValue(Node)
func NodeUtil_GetBigIntValue(n rhino.Node) *big.Int {
	switch n.GetToken() {
	case rhino.Token_NUMBER:
		val := n.GetDouble()
		if base.JSCompDoubles_IsAtLeastIntegerPrecision(val) && base.JSCompDoubles_IsExactInt64(val) {
			return nodeUtil_bigIntFromInt64(jutil.D2L(val))
		}
		return nil

	case rhino.Token_BIGINT:
		return n.GetBigInt()

	case rhino.Token_FALSE, rhino.Token_NOT, rhino.Token_TRUE:
		switch NodeUtil_GetBooleanValue(n) {
		case base.Tri_TRUE:
			return &nodeUtil_smallBigInts[1-nodeUtil_smallBigIntMin]
		case base.Tri_FALSE:
			return &nodeUtil_smallBigInts[0-nodeUtil_smallBigIntMin]
		default: // case UNKNOWN
			return nil
		}

	case rhino.Token_TEMPLATELIT:
		string_ := NodeUtil_GetStringValue(n)
		if string_ == nil {
			return nil
		}
		return NodeUtil_GetStringBigIntValue(*string_)

	case rhino.Token_STRINGLIT:
		return NodeUtil_GetStringBigIntValue(n.GetString())

	case rhino.Token_NEG:
		result := NodeUtil_GetBigIntValue(n.GetOnlyChild())
		if result == nil {
			return nil
		}
		return new(big.Int).Neg(result)
	case rhino.Token_BITNOT:
		result := NodeUtil_GetBigIntValue(n.GetOnlyChild())
		if result == nil {
			return nil
		}
		return new(big.Int).Not(result)

	case rhino.Token_ARRAYLIT, rhino.Token_OBJECTLIT:
		value := NodeUtil_GetStringValue(n)
		if value != nil {
			return NodeUtil_GetStringBigIntValue(*value)
		}
		return nil
	case rhino.Token_VOID, rhino.Token_NAME, rhino.Token_NULL:
		return nil
	default:
		return nil
	}
}

// Java: NodeUtil.getStringBigIntValue(String)
func NodeUtil_GetStringBigIntValue(rawJsString string) *big.Int {
	if strings.Contains(rawJsString, "\u000b") {
		// vertical tab is not always whitespace
		return nil
	}

	s := NodeUtil_TrimJsWhiteSpace(rawJsString)
	if s == "" {
		return big.NewInt(0)
	}

	if jutil.StringLength(s) > 2 && jutil.StringCharAt(s, 0) == '0' {
		// Attempt to convert hex, octal, and binary formats.
		var radix int
		switch jutil.StringCharAt(s, 1) {
		case 'x', 'X':
			radix = 16
		case 'o', 'O':
			radix = 8
		case 'b', 'B':
			radix = 2
		default:
			radix = 0
		}
		if radix != 0 {
			var result *big.Int
			jutil.Try(func() {
				result = jutil.NewBigIntegerFromStringRadix(jutil.StringSubstringFrom(s, 2), radix)
			}, func(e jutil.INumberFormatException) {
				result = nil
			})
			return result
		}
	}

	var result *big.Int
	jutil.Try(func() {
		// Java: new BigInteger(s), i.e. new BigInteger(s, 10).
		result = jutil.NewBigIntegerFromStringRadix(s, 10)
	}, func(e jutil.INumberFormatException) {
		result = nil
	})
	return result
}

// Java: NodeUtil.trimJsWhiteSpace(String)
func NodeUtil_TrimJsWhiteSpace(s string) string {
	units := jutil.StringToUTF16(s)
	start := 0
	end := len(units)
	for end > 0 && rhino.TokenUtil_IsStrWhiteSpaceChar(int(units[end-1])) == base.Tri_TRUE {
		end--
	}
	for start < end && rhino.TokenUtil_IsStrWhiteSpaceChar(int(units[start])) == base.Tri_TRUE {
		start++
	}
	return jutil.StringFromUTF16(units[start:end])
}

// NodeUtil_GetName returns the name of the given function or class, if it
// has one. n is a function or class node.
//
// Java: NodeUtil.getName(Node)
func NodeUtil_GetName(n rhino.Node) *string {
	nameNode := NodeUtil_GetNameNode(n)
	if nameNode.IsNil() {
		return nil
	}
	return nameNode.GetQualifiedName()
}

func NodeUtil_GetNameString(n rhino.Node) (string, bool) {
	nameNode := NodeUtil_GetNameNode(n)
	if nameNode.IsNil() {
		return "", false
	}
	return nameNode.GetQualifiedNameString()
}

// NodeUtil_GetNameNode gets the node of a function or class's name. This
// method recognizes five forms:
//
//   - class name {...}
//   - var name = class {...}
//   - qualified.name = class {...}
//   - var name2 = class name1 {...}
//   - qualified.name2 = class name1 {...}
//
// In two last cases with named function expressions, the second name is
// returned (the variable or qualified name).
//
// n is a function or class node; returns the node best representing the
// class's name.
//
// Java: NodeUtil.getNameNode(Node)
func NodeUtil_GetNameNode(n rhino.Node) rhino.Node {
	jutil.CheckStateMsg1(n.IsFunction() || n.IsClass(), "%s", n)
	parent := n.GetParent()
	switch parent.GetToken() {
	case rhino.Token_NAME:
		// var name = function() ...
		// var name2 = function name1() ...
		return parent
	case rhino.Token_ASSIGN:
		// qualified.name = function() ...
		// qualified.name2 = function name1() ...
		firstChild := parent.GetFirstChild()
		if firstChild.IsQualifiedName() {
			return firstChild
		}
		return rhino.NilNode
	default:
		// function name() ...
		//   or
		// class Name ...
		funNameNode := n.GetFirstChild()
		// Don't return the name node for anonymous functions/classes.
		// TODO(tbreisacher): Currently we do two kinds of "empty" checks because
		// anonymous classes have an EMPTY name node while anonymous functions
		// have a STRING node with an empty string. Consider making these the same.
		if funNameNode.IsEmpty() || funNameNode.GetString() == "" {
			return rhino.NilNode
		}
		return funNameNode
	}
}

// NodeUtil_RemoveName sets the given function/class node to an empty name.
//
// Java: NodeUtil.removeName(Node)
func NodeUtil_RemoveName(n rhino.Node) {
	jutil.CheckState(n.IsFunction() || n.IsClass())
	originalName := n.GetFirstChild()
	var emptyName rhino.Node
	if n.IsFunction() {
		emptyName = rhino.IR_Name("")
	} else {
		emptyName = rhino.IR_Empty()
	}
	originalName.ReplaceWith(emptyName.Srcref(originalName))
}

// NodeUtil_GetNearestFunctionName gets the function's name. This method
// recognizes the forms:
//
//   - {'name': function() ...}
//   - {name: function() ...}
//   - function name() ...
//   - var name = function() ...
//   - var obj = {name() {} ...}
//   - qualified.name = function() ...
//   - var name2 = function name1() ...
//   - qualified.name2 = function name1() ...
//
// n is a node whose type is Token#FUNCTION; returns the function's name, or
// null if it has no name.
//
// Java: NodeUtil.getNearestFunctionName(Node)
func NodeUtil_GetNearestFunctionName(n rhino.Node) *string {
	if !n.IsFunction() {
		return nil
	}

	name := NodeUtil_GetName(n)
	if name != nil {
		return name
	}

	// Check for the form { 'x' : function() { }} and {x() {}}
	parent := n.GetParent()
	switch parent.GetToken() {
	// Return the name of the literal's key.
	case rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_SETTER_DEF, rhino.Token_GETTER_DEF, rhino.Token_STRING_KEY:
		return jutil.StringPtr(parent.GetString())
	case rhino.Token_NUMBER:
		return NodeUtil_GetStringValue(parent)
	default:
		return nil
	}
}

// Java: NodeUtil.getClassMembers(Node)
func NodeUtil_GetClassMembers(n rhino.Node) rhino.Node {
	jutil.CheckArgument(n.IsClass())
	return n.GetLastChild()
}

// Java: NodeUtil.getEs6ClassConstructorMemberFunctionDef(Node)
func NodeUtil_GetEs6ClassConstructorMemberFunctionDef(classNode rhino.Node) rhino.Node {
	jutil.CheckArgumentMsg1(classNode.IsClass(), "%s", classNode)
	classMembers := rhino.CheckNotNullMsg1(classNode.GetLastChild(), "%s", classNode)
	for memberFunctionDef := classMembers.GetFirstChild(); !memberFunctionDef.IsNil(); memberFunctionDef = memberFunctionDef.GetNext() {
		if NodeUtil_IsEs6ConstructorMemberFunctionDef(memberFunctionDef) {
			return memberFunctionDef
		}
	}

	return rhino.NilNode
}

// NodeUtil_IsImmutableValue returns true if this is an immutable value.
//
// Java: NodeUtil.isImmutableValue(Node)
func NodeUtil_IsImmutableValue(n rhino.Node) bool {
	// TODO(johnlenz): rename this function.  It is currently being used
	// in two disjoint cases:
	// 1) We only care about the result of the expression
	//    (in which case NOT here should return true)
	// 2) We care that expression is a side-effect free and can't
	//    be side-effected by other expressions.
	// This should only be used to say the value is immutable and
	// hasSideEffects and canBeSideEffected should be used for the other case.

	switch n.GetToken() {
	case rhino.Token_STRINGLIT, rhino.Token_NUMBER, rhino.Token_BIGINT, rhino.Token_NULL, rhino.Token_TRUE,
		rhino.Token_FALSE:
		return true
	case rhino.Token_CAST, rhino.Token_NOT, rhino.Token_VOID, rhino.Token_NEG:
		return NodeUtil_IsImmutableValue(n.GetFirstChild())
	case rhino.Token_NAME:
		name := n.GetString()
		// We assume here that programs don't change the value of the keyword
		// undefined to something other than the value undefined.
		return "undefined" == name || "Infinity" == name || "NaN" == name
	case rhino.Token_TEMPLATELIT:
		for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
			if child.IsTemplateLitSub() {
				if !NodeUtil_IsImmutableValue(child.GetFirstChild()) {
					return false
				}
			}
		}
		return true
	// TODO(yitingwang) There are probably other tokens that shouldn't get to the default branch
	case rhino.Token_TEMPLATELIT_STRING:
		panic(jutil.NewIllegalArgumentExceptionMsg("Invalid argument " + n.String()))
	default:
		return false
	}
}

// NodeUtil_IsNumericLiteral returns true if n is a numeric literal or a unary
// +/- applied directly to a numeric literal (e.g. 12.3, -12.3, +12.3).
//
// Java: NodeUtil.isNumericLiteral(Node)
func NodeUtil_IsNumericLiteral(n rhino.Node) bool {
	if n.IsNil() {
		return false
	}
	return n.IsNumber() ||
		((n.IsNeg() || n.IsPos()) && n.HasOneChild() && n.GetFirstChild().IsNumber())
}

// NodeUtil_GetNumericLiteralValue returns the numeric value of n, which must
// satisfy isNumericLiteral(Node); it throws IllegalArgumentException if n
// does not.
//
// Java: NodeUtil.getNumericLiteralValue(Node)
func NodeUtil_GetNumericLiteralValue(n rhino.Node) float64 {
	jutil.CheckArgumentMsg1(NodeUtil_IsNumericLiteral(n), "%s", n)
	switch n.GetToken() {
	case rhino.Token_NUMBER:
		return n.GetDouble()
	case rhino.Token_POS:
		return n.GetFirstChild().GetDouble()
	case rhino.Token_NEG:
		return -n.GetFirstChild().GetDouble()
	default:
		panic(jutil.NewAssertionError("Unreachable: " + n.String()))
	}
}

// NodeUtil_IsSymmetricOperation returns true if the operator on this node is
// symmetric.
//
// Java: NodeUtil.isSymmetricOperation(Node)
func NodeUtil_IsSymmetricOperation(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_EQ, // equal
		rhino.Token_NE,   // not equal
		rhino.Token_SHEQ, // exactly equal
		rhino.Token_SHNE, // exactly not equal
		rhino.Token_MUL:  // multiply, unlike add it only works on numbers
		// or results NaN if any of the operators is not a number
		return true
	default:
		return false
	}
}

// NodeUtil_IsRelationalOperation returns true if the operator on this node is
// relational. the returned set does not include the equalities.
//
// Java: NodeUtil.isRelationalOperation(Node)
func NodeUtil_IsRelationalOperation(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_GT, rhino.Token_GE, rhino.Token_LT, rhino.Token_LE:
		return true
	default:
		return false
	}
}

// NodeUtil_GetInverseOperator returns the inverse of an operator if it is
// invertible. ex. '>' ==> '<'
//
// Java: NodeUtil.getInverseOperator(Token)
func NodeUtil_GetInverseOperator(type_ rhino.Token) rhino.Token {
	switch type_ {
	case rhino.Token_GT:
		return rhino.Token_LT
	case rhino.Token_LT:
		return rhino.Token_GT
	case rhino.Token_GE:
		return rhino.Token_LE
	case rhino.Token_LE:
		return rhino.Token_GE
	default:
		panic(jutil.NewIllegalArgumentExceptionMsg("Unexpected token: " + type_.String()))
	}
}

// NodeUtil_IsLiteralValue returns true if this is a literal value. We define
// a literal value as any node that evaluates to the same thing regardless of
// when or where it is evaluated. So /xyz/ and [3, 5] are literals, but the
// name a is not.
//
// Function literals do not meet this definition, because they lexically
// capture variables. For example, if you have
// `function() { return a; }`. If it is evaluated in a different scope, then
// it captures a different variable. Even if the function did not read any
// captured variables directly, it would still fail this definition, because
// it affects the lifecycle of variables in the enclosing scope.
//
// However, a function literal with respect to a particular scope is a
// literal.
//
// includeFunctions: if true, all function expressions will be treated as
// literals.
//
// Java: NodeUtil.isLiteralValue(Node, boolean)
func NodeUtil_IsLiteralValue(n rhino.Node, includeFunctions bool) bool {
	switch n.GetToken() {
	case rhino.Token_CAST:
		return NodeUtil_IsLiteralValue(n.GetFirstChild(), includeFunctions)
	case rhino.Token_ARRAYLIT:
		for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
			if (!child.IsEmpty()) && !NodeUtil_IsLiteralValue(child, includeFunctions) {
				return false
			}
		}
		return true
	case rhino.Token_REGEXP:
		// Return true only if all descendants are const.
		for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
			if !NodeUtil_IsLiteralValue(child, includeFunctions) {
				return false
			}
		}
		return true
	case rhino.Token_OBJECTLIT:
		for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
			switch child.GetToken() {
			case rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF:
				// { methodName() {...} }
				// { get propertyName() {...} }
				// { set propertyName(value) {...} }
				if !includeFunctions {
					return false
				}
				// keep going

			case rhino.Token_COMPUTED_PROP:
				// { [key_expression]: value, ... }
				// { [key_expression](args) {...}, ... }
				// { get [key_expression]() {...}, ... }
				// { set [key_expression](args) {...}, ... }
				if !NodeUtil_IsLiteralValue(child.GetFirstChild(), includeFunctions) ||
					!NodeUtil_IsLiteralValue(child.GetLastChild(), includeFunctions) {
					return false
				}
				// keep going

			case rhino.Token_OBJECT_SPREAD:
				if !NodeUtil_IsLiteralValue(child.GetOnlyChild(), includeFunctions) {
					return false
				}
				// keep going

			case rhino.Token_STRING_KEY:
				// { key: value, ... }
				// { "quoted_key": value, ... }
				if !NodeUtil_IsLiteralValue(child.GetOnlyChild(), includeFunctions) {
					return false
				}
				// keep going

			default:
				panic(jutil.NewIllegalArgumentExceptionMsg(
					"Unexpected child of OBJECTLIT: " + child.ToStringTree()))
			}
		}
		return true
	case rhino.Token_FUNCTION:
		return includeFunctions && !NodeUtil_IsFunctionDeclaration(n)
	case rhino.Token_TEMPLATELIT:
		for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
			if child.IsTemplateLitSub() {
				if !NodeUtil_IsLiteralValue(child.GetFirstChild(), includeFunctions) {
					return false
				}
			}
		}
		return true
	default:
		return NodeUtil_IsImmutableValue(n)
	}
}

// NodeUtil_FindClosureUnawareScriptRoot returns the closureUnaware function
// body within a script modelling a shadow AST.
//
// TranspileAndOptimizeClosureUnaware pass models each @closureUnaware block as
// a function body, and will throw a runtime exception if any pass inserts new
// code outside that function body. So we sometimes need to map from an actual
// script node to the corresponding @closureUnaware script-esque block.
//
// TODO: lharker - if we end up needing this special logic often, consider
// instead making TranspileAndOptimizeclosureUnaware support inserting code at
// the script level.
//
// Java: NodeUtil.findClosureUnawareScriptRoot(Node)
func NodeUtil_FindClosureUnawareScriptRoot(script rhino.Node) rhino.Node {
	jutil.CheckArgumentMsg1(script.IsScript() && script.GetIsInClosureUnawareSubtree(), "%s", script)
	return NodeUtil_FindPreorder(script, (rhino.Node).IsBlock, func(c rhino.Node) bool { return true })
}

// NodeUtil_GetInsertionPointAfterAllInnerFunctionDeclarations returns the
// node within the given block at which something can be inserted such that
// it's after all inner function declarations in that block. We want this
// because normalization expects all inner function declarations to be
// hoisted. If there's not good insertion point (e.g. the block is empty or
// only contains inner function declarations), return null.
//
// Java: NodeUtil.getInsertionPointAfterAllInnerFunctionDeclarations(Node)
func NodeUtil_GetInsertionPointAfterAllInnerFunctionDeclarations(block rhino.Node) rhino.Node {
	jutil.CheckState(block.IsBlock())
	current := block.GetFirstChild()

	for !current.IsNil() && NodeUtil_IsFunctionDeclaration(current) {
		current = current.GetNext()
	}
	return current
}

// NodeUtil_IsSomeCompileTimeConstStringValue returns true iff the value
// associated with the node is a JS string literal, a concatenation thereof or
// a ternary operator choosing between string literals.
//
// Java: NodeUtil.isSomeCompileTimeConstStringValue(Node)
func NodeUtil_IsSomeCompileTimeConstStringValue(node rhino.Node) bool {
	// TODO(bangert): Support constants, using a Scope argument. See ConstParamCheck
	if node.IsStringLit() || (node.IsTemplateLit() && node.HasOneChild()) {
		return true
	} else if node.IsAdd() {
		jutil.CheckStateMsg1(node.HasTwoChildren(), "%s", node)
		left := node.GetFirstChild()
		right := node.GetLastChild()
		return NodeUtil_IsSomeCompileTimeConstStringValue(left) && NodeUtil_IsSomeCompileTimeConstStringValue(right)
	} else if node.IsHook() {
		// Ternary operator a ? b : c
		left := node.GetSecondChild()
		right := node.GetLastChild()
		return NodeUtil_IsSomeCompileTimeConstStringValue(left) && NodeUtil_IsSomeCompileTimeConstStringValue(right)
	}
	return false
}

// NodeUtil_IsEmptyBlock returns whether this a BLOCK node with no children.
//
// Java: NodeUtil.isEmptyBlock(Node)
func NodeUtil_IsEmptyBlock(block rhino.Node) bool {
	if !block.IsBlock() {
		return false
	}

	for n := block.GetFirstChild(); !n.IsNil(); n = n.GetNext() {
		if !n.IsEmpty() {
			return false
		}
	}
	return true
}

// Java: NodeUtil.isBinaryOperator(Node)
func NodeUtil_IsBinaryOperator(n rhino.Node) bool {
	return NodeUtil_IsBinaryOperatorType(n.GetToken())
}

// NodeUtil_IsBinaryOperatorType: an operator with two operands that does not
// assign a value to either. Once you cut through the layers of rules, these
// all parse similarly, taking LeftHandSideExpression operands on either side.
// Comma is not included, because it takes AssignmentExpression operands,
// making its syntax different.
//
// Java: NodeUtil.isBinaryOperatorType(Token)
func NodeUtil_IsBinaryOperatorType(type_ rhino.Token) bool {
	switch type_ {
	case rhino.Token_OR,
		rhino.Token_AND,
		rhino.Token_COALESCE,
		rhino.Token_BITOR,
		rhino.Token_BITXOR,
		rhino.Token_BITAND,
		rhino.Token_EQ,
		rhino.Token_NE,
		rhino.Token_SHEQ,
		rhino.Token_SHNE,
		rhino.Token_LT,
		rhino.Token_GT,
		rhino.Token_LE,
		rhino.Token_GE,
		rhino.Token_INSTANCEOF,
		rhino.Token_IN,
		rhino.Token_LSH,
		rhino.Token_RSH,
		rhino.Token_URSH,
		rhino.Token_ADD,
		rhino.Token_SUB,
		rhino.Token_MUL,
		rhino.Token_DIV,
		rhino.Token_MOD,
		rhino.Token_EXPONENT:
		return true
	default:
		return false
	}
}

// Java: NodeUtil.isUnaryOperator(Node)
func NodeUtil_IsUnaryOperator(n rhino.Node) bool {
	return NodeUtil_IsUnaryOperatorType(n.GetToken())
}

// NodeUtil_IsUnaryOperatorType: an operator taking only one operand. These
// all parse very similarly, taking LeftHandSideExpression operands.
//
// Java: NodeUtil.isUnaryOperatorType(Token)
func NodeUtil_IsUnaryOperatorType(type_ rhino.Token) bool {
	switch type_ {
	case rhino.Token_DELPROP, rhino.Token_VOID, rhino.Token_TYPEOF, rhino.Token_POS, rhino.Token_NEG,
		rhino.Token_BITNOT, rhino.Token_NOT:
		return true
	default:
		return false
	}
}

// Java: NodeUtil.isUpdateOperator(Node)
func NodeUtil_IsUpdateOperator(n rhino.Node) bool {
	return NodeUtil_IsUpdateOperatorType(n.GetToken())
}

// Java: NodeUtil.isUpdateOperatorType(Token)
func NodeUtil_IsUpdateOperatorType(type_ rhino.Token) bool {
	switch type_ {
	case rhino.Token_INC, rhino.Token_DEC:
		return true
	default:
		return false
	}
}

// Java: NodeUtil.isSimpleOperator(Node)
func NodeUtil_IsSimpleOperator(n rhino.Node) bool {
	return NodeUtil_IsSimpleOperatorType(n.GetToken())
}

// NodeUtil_IsSimpleOperatorType: a "simple" operator is one whose children
// are expressions, has no direct side-effects (unlike '+='), and has no
// conditional aspects (unlike '||').
//
// Java: NodeUtil.isSimpleOperatorType(Token)
func NodeUtil_IsSimpleOperatorType(type_ rhino.Token) bool {
	switch type_ {
	case rhino.Token_ADD,
		rhino.Token_BITAND,
		rhino.Token_BITNOT,
		rhino.Token_BITOR,
		rhino.Token_BITXOR,
		rhino.Token_COMMA,
		rhino.Token_DIV,
		rhino.Token_EQ,
		rhino.Token_EXPONENT,
		rhino.Token_GE,
		rhino.Token_GT,
		rhino.Token_IN,
		rhino.Token_INSTANCEOF,
		rhino.Token_LE,
		rhino.Token_LSH,
		rhino.Token_LT,
		rhino.Token_MOD,
		rhino.Token_MUL,
		rhino.Token_NE,
		rhino.Token_NOT,
		rhino.Token_RSH,
		rhino.Token_SHEQ,
		rhino.Token_SHNE,
		rhino.Token_SUB,
		rhino.Token_TYPEOF,
		rhino.Token_VOID,
		rhino.Token_POS,
		rhino.Token_NEG,
		rhino.Token_URSH:
		return true
	default:
		return false
	}
}

// NodeUtil_IsNamespaceDecl returns true iff this node defines a namespace,
// e.g.,
//
//	/** @const */ var goog = {};
//	/** @const */ var goog = goog || {};
//	/** @const */ goog.math = goog.math || {};
//
// Java: NodeUtil.isNamespaceDecl(Node)
func NodeUtil_IsNamespaceDecl(n rhino.Node) bool {
	jsdoc := NodeUtil_GetBestJSDocInfo(n)
	if jsdoc != nil && jsdoc.HasTypeNodes() && !jsdoc.HasTypedefType() {
		return false
	}
	// In externs, we allow namespace definitions without @const.
	// This is a worse design than always requiring @const, but it helps with
	// namespaces that are defined in many places, such as gapi.
	// Also, omitting @const in externs is not as confusing as in source code,
	// because assigning an object literal in externs only makes sense when
	// defining a namespace or enum.
	isMarkedConst := n.GetParent().IsConst() || (jsdoc != nil && jsdoc.IsConstant())
	if !n.IsFromExterns() && !isMarkedConst {
		return false
	}
	var qnameNode rhino.Node
	var initializer rhino.Node
	if NodeUtil_IsNameDeclaration(n.GetParent()) {
		qnameNode = n
		initializer = n.GetFirstChild()
	} else if n.IsExprResult() {
		expr := n.GetFirstChild()
		if !expr.IsAssign() || !expr.GetFirstChild().IsGetProp() {
			return false
		}
		qnameNode = expr.GetFirstChild()
		initializer = expr.GetLastChild()
	} else if n.IsGetProp() {
		parent := n.GetParent()
		if !parent.IsAssign() || !parent.GetParent().IsExprResult() {
			return false
		}
		qnameNode = n
		initializer = parent.GetLastChild()
	} else {
		return false
	}
	if initializer.IsNil() || qnameNode.IsNil() {
		return false
	}
	if initializer.IsObjectLit() {
		return true
	}
	return initializer.IsOr() &&
		qnameNode.MatchesQualifiedNameNode(initializer.GetFirstChild()) &&
		initializer.GetLastChild().IsObjectLit()
}

// NodeUtil_IsFromTypeSummary determines if the given SCRIPT is a
// @typeSummary file, like an i.js file.
//
// Java: NodeUtil.isFromTypeSummary(Node)
func NodeUtil_IsFromTypeSummary(n rhino.Node) bool {
	jutil.CheckArgumentMsg1(n.IsScript(), "%s", n)
	info := n.GetJSDocInfo()
	return info != nil && info.IsTypeSummary()
}

// NodeUtil_NewExpr creates an EXPR_RESULT.
//
// child is the expression itself; returns a newly created EXPR node with the
// child as subexpression.
//
// Java: NodeUtil.newExpr(Node)
func NodeUtil_NewExpr(child rhino.Node) rhino.Node {
	return rhino.IR_ExprResult(child).Srcref(child)
}

// NodeUtil_IteratesImpureIterable returns true if node might execute an
// `Iterable` iteration that has side-effects, false if there are definitely
// no such side-effects.
//
// This function only considers purity of the iteration. Other expressions
// within the node subtree may still have side-effects.
//
// Throws IllegalStateException if node is of a kind that does not trigger
// iteration. An explicit goal of this function is to record all the kinds of
// nodes that do.
//
// Java: NodeUtil.iteratesImpureIterable(Node)
func NodeUtil_IteratesImpureIterable(node rhino.Node) bool {
	parent := node.GetParent()

	var iterable rhino.Node
	switch node.GetToken() {
	case rhino.Token_ITER_SPREAD:
		iterable = node.GetOnlyChild()
	case rhino.Token_YIELD:
		if !node.IsYieldAll() {
			return false // Regular `yield` does not iterate, only `yield*`.
		}
		iterable = node.GetOnlyChild()
	case rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
		iterable = node.GetSecondChild()
	case rhino.Token_ITER_REST:
		switch parent.GetToken() {
		// Rest arguments are flat at the call-site.
		case rhino.Token_PARAM_LIST:
			return false
		// TODO(b/127862986): We assume the r-value to be an impure iterable.
		case rhino.Token_ARRAY_PATTERN:
			return true
		default:
			panic(jutil.NewIllegalStateExceptionMsg(
				"Unexpected parent of ITRE_REST: " + parent.ToStringTree()))
		}
	default:
		panic(jutil.NewIllegalStateExceptionMsg(
			"Expected a kind of node that may trigger iteration: " + node.ToStringTree()))
	}

	return !NodeUtil_IsPureIterable(iterable)
}

// NodeUtil_IsPureIterable returns true if node is guaranteed to be an
// `Iterable` that causes no side-effects during iteration, false otherwise.
//
// Java: NodeUtil.isPureIterable(Node)
func NodeUtil_IsPureIterable(node rhino.Node) bool {
	// TODO(b/127862986): The type of the iterable should also allow us to say it's pure.
	switch node.GetToken() {
	// These iterables are known to be pure.
	case rhino.Token_ARRAYLIT, rhino.Token_STRINGLIT, rhino.Token_TEMPLATELIT:
		return true
	// Anything else, including a non-iterable (e.g. `null`), would be impure.
	default:
		return false
	}
}

// NodeUtil_NewHasLocalResult returns whether the new has a local result.
//
// Java: NodeUtil.newHasLocalResult(Node)
func NodeUtil_NewHasLocalResult(n rhino.Node) bool {
	jutil.CheckStateMsg1(n.IsNew(), "%s", n)
	return n.IsOnlyModifiesThisCall()
}

// Java: NodeUtil.allArgsUnescapedLocal(Node)
func NodeUtil_AllArgsUnescapedLocal(callOrNew rhino.Node) bool {
	for arg := callOrNew.GetSecondChild(); !arg.IsNil(); arg = arg.GetNext() {
		if !NodeUtil_EvaluatesToLocalValue(arg) {
			return false
		}
	}
	return true
}

// NodeUtil_KNOWN_CONSTANTS: we will assume that side effects never change
// the values of these global variables.
//
// Java: NodeUtil.KNOWN_CONSTANTS
var NodeUtil_KNOWN_CONSTANTS = map[string]bool{
	"undefined": true,
	"Infinity":  true,
	"NaN":       true,
}

// NodeUtil_CanBeSideEffected returns whether the tree can be affected by
// side-effects or has side-effects.
//
// Java: NodeUtil.canBeSideEffected(Node)
func NodeUtil_CanBeSideEffected(n rhino.Node) bool {
	return NodeUtil_CanBeSideEffectedWithScope(n, NodeUtil_KNOWN_CONSTANTS, nil)
}

// NodeUtil_CanBeSideEffectedWithScope: knownConstants is a set of names
// known to be constant value at node 'n' (such as locals that are last
// written before n can execute). Returns whether the tree can be affected by
// side-effects or has side-effects.
//
// Java: NodeUtil.canBeSideEffected(Node, Set<String>, Scope)
//
// TODO(nick): Get rid of the knownConstants argument in favor of using
// scope with InferConsts.
func NodeUtil_CanBeSideEffectedWithScope(n rhino.Node, knownConstants map[string]bool, scope *Scope) bool {
	switch n.GetToken() {
	case rhino.Token_YIELD, rhino.Token_CALL, rhino.Token_OPTCHAIN_CALL, rhino.Token_NEW, rhino.Token_TAGGED_TEMPLATELIT, rhino.Token_AWAIT, rhino.Token_DYNAMIC_IMPORT:
		// Function calls or constructor can reference changed values.
		// TODO(johnlenz): Add some mechanism for determining that functions
		// are unaffected by side effects.
		return true
	case rhino.Token_NAME:
		// Non-constant names values may have been changed.
		return !NodeUtil_IsConstantVar(n, scope) && !knownConstants[n.GetString()]

	// Properties on constant NAMEs can still be side-effected.
	case rhino.Token_GETPROP:
		// Assume well-known symbols will not be modified.
		return !n.GetFirstChild().MatchesName("Symbol")
	// Properties on constant NAMEs can still be side-effected.
	case rhino.Token_GETELEM, rhino.Token_OPTCHAIN_GETPROP, rhino.Token_OPTCHAIN_GETELEM:
		return true
	case rhino.Token_FUNCTION:
		// Function expression are not changed by side-effects,
		// and function declarations are not part of expressions.
		// TODO(bradfordcsmith): Do we need to add a case for CLASS here?
		//     This checkState currently does not exclude class methods.
		jutil.CheckStateMsg1(!NodeUtil_IsFunctionDeclaration(n), "%s", n)
		return false
	default:
	}

	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if NodeUtil_CanBeSideEffectedWithScope(c, knownConstants, scope) {
			return true
		}
	}

	return false
}

// NodeUtil_Precedence: the comma operator has the lowest precedence, 0,
// followed by the assignment operators (=, &=, +=, etc.) which have
// precedence of 1, and so on.
//
// See
// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Operators/Operator_Precedence
//
// Java: NodeUtil.precedence(Token)
func NodeUtil_Precedence(type_ rhino.Token) int {
	switch type_ {
	case rhino.Token_COMMA:
		return 0
	case rhino.Token_ASSIGN_BITOR,
		rhino.Token_ASSIGN_BITXOR,
		rhino.Token_ASSIGN_BITAND,
		rhino.Token_ASSIGN_LSH,
		rhino.Token_ASSIGN_RSH,
		rhino.Token_ASSIGN_URSH,
		rhino.Token_ASSIGN_ADD,
		rhino.Token_ASSIGN_SUB,
		rhino.Token_ASSIGN_MUL,
		rhino.Token_ASSIGN_EXPONENT,
		rhino.Token_ASSIGN_DIV,
		rhino.Token_ASSIGN_MOD,
		rhino.Token_ASSIGN_OR,
		rhino.Token_ASSIGN_AND,
		rhino.Token_ASSIGN_COALESCE,
		rhino.Token_ASSIGN:
		return 1
	case rhino.Token_YIELD:
		return 2
	// ?: operator
	case rhino.Token_HOOK:
		return 3
	case rhino.Token_OR:
		return 4
	case rhino.Token_AND:
		return 5
	case rhino.Token_COALESCE:
		return 6
	case rhino.Token_BITOR:
		return 7
	case rhino.Token_BITXOR:
		return 8
	case rhino.Token_BITAND:
		return 9
	case rhino.Token_EQ, rhino.Token_NE, rhino.Token_SHEQ, rhino.Token_SHNE:
		return 10
	case rhino.Token_LT, rhino.Token_GT, rhino.Token_LE, rhino.Token_GE, rhino.Token_INSTANCEOF, rhino.Token_IN:
		return 11
	case rhino.Token_LSH, rhino.Token_RSH, rhino.Token_URSH:
		return 12
	case rhino.Token_SUB, rhino.Token_ADD:
		return 13
	case rhino.Token_MUL, rhino.Token_MOD, rhino.Token_DIV:
		return 14
	case rhino.Token_EXPONENT:
		return 15
	// Unary operators
	case rhino.Token_AWAIT, rhino.Token_NEW, rhino.Token_DELPROP, rhino.Token_TYPEOF, rhino.Token_VOID,
		rhino.Token_NOT, rhino.Token_BITNOT, rhino.Token_POS, rhino.Token_NEG:
		return 16
	// Update operators
	case rhino.Token_INC, rhino.Token_DEC:
		return 17
	case rhino.Token_CALL,
		rhino.Token_GETELEM,
		rhino.Token_GETPROP,
		rhino.Token_OPTCHAIN_CALL,
		rhino.Token_OPTCHAIN_GETELEM,
		rhino.Token_OPTCHAIN_GETPROP,
		rhino.Token_NEW_TARGET,
		rhino.Token_IMPORT_META,
		// Data values
		rhino.Token_ARRAYLIT,
		rhino.Token_ARRAY_PATTERN,
		rhino.Token_DEFAULT_VALUE,
		rhino.Token_DESTRUCTURING_LHS,
		rhino.Token_EMPTY, // TODO(johnlenz): remove this.
		rhino.Token_FALSE,
		rhino.Token_FUNCTION,
		rhino.Token_CLASS,
		rhino.Token_INTERFACE,
		rhino.Token_NAME,
		rhino.Token_NULL,
		rhino.Token_NUMBER,
		rhino.Token_BIGINT,
		rhino.Token_OBJECTLIT,
		rhino.Token_OBJECT_PATTERN,
		rhino.Token_REGEXP,
		rhino.Token_ITER_REST,
		rhino.Token_OBJECT_REST,
		rhino.Token_ITER_SPREAD,
		rhino.Token_OBJECT_SPREAD,
		rhino.Token_STRINGLIT,
		rhino.Token_STRING_KEY,
		rhino.Token_MEMBER_VARIABLE_DEF,
		rhino.Token_INDEX_SIGNATURE,
		rhino.Token_CALL_SIGNATURE,
		rhino.Token_THIS,
		rhino.Token_SUPER,
		rhino.Token_TRUE,
		rhino.Token_TAGGED_TEMPLATELIT,
		rhino.Token_TEMPLATELIT,
		rhino.Token_DYNAMIC_IMPORT,
		// Tokens from the type declaration AST
		rhino.Token_UNION_TYPE:
		return 18
	case rhino.Token_FUNCTION_TYPE:
		return 19
	case rhino.Token_ARRAY_TYPE, rhino.Token_PARAMETERIZED_TYPE:
		return 20
	case rhino.Token_STRING_TYPE,
		rhino.Token_NUMBER_TYPE,
		rhino.Token_BOOLEAN_TYPE,
		rhino.Token_ANY_TYPE,
		rhino.Token_RECORD_TYPE,
		rhino.Token_NULLABLE_TYPE,
		rhino.Token_NAMED_TYPE,
		rhino.Token_UNDEFINED_TYPE,
		rhino.Token_VOID_TYPE,
		rhino.Token_GENERIC_TYPE:
		return 21
	case rhino.Token_CAST:
		return 22
	default:
		jutil.CheckArgument(type_ != rhino.Token_TEMPLATELIT_STRING)
		panic(jutil.NewIllegalStateExceptionMsg("Unknown precedence for " + type_.String()))
	}
}

// Java: NodeUtil.isUndefined(Node)
func NodeUtil_IsUndefined(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_VOID:
		return true
	case rhino.Token_NAME:
		return n.GetString() == "undefined"
	default:
		return false
	}
}

// Java: NodeUtil.isNullOrUndefined(Node)
func NodeUtil_IsNullOrUndefined(n rhino.Node) bool {
	return n.IsNull() || NodeUtil_IsUndefined(n)
}

// NodeUtil_ValueType: see NodeUtil_GetKnownValueType.
//
// Java: NodeUtil.ValueType
type NodeUtil_ValueType int32

const (
	NodeUtil_ValueType_UNDETERMINED NodeUtil_ValueType = iota
	NodeUtil_ValueType_NULL
	NodeUtil_ValueType_VOID
	NodeUtil_ValueType_NUMBER
	NodeUtil_ValueType_BIGINT
	NodeUtil_ValueType_STRING
	NodeUtil_ValueType_BOOLEAN
	NodeUtil_ValueType_OBJECT
)

var nodeUtil_ValueType_names = [...]string{
	"UNDETERMINED", "NULL", "VOID", "NUMBER", "BIGINT", "STRING", "BOOLEAN", "OBJECT",
}

// Java: NodeUtil.ValueType.name()
func (t NodeUtil_ValueType) Name() string { return nodeUtil_ValueType_names[t] }

// Java: NodeUtil.ValueType.toString()
func (t NodeUtil_ValueType) String() string { return nodeUtil_ValueType_names[t] }

// Java: NodeUtil.ValueType.ordinal()
func (t NodeUtil_ValueType) Ordinal() int { return int(t) }

// NodeUtil_GetKnownValueType evaluates a node's token and attempts to
// determine which primitive value type it could resolve to. Without proper
// type information some assumptions had to be made for operations that could
// result in a BigInt or a Number. If there is not enough information
// available to determine one or the other then we assume Number in order to
// maintain historical behavior of the compiler and avoid breaking projects
// that relied on this behavior.
//
// Java: NodeUtil.getKnownValueType(Node)
func NodeUtil_GetKnownValueType(n rhino.Node) NodeUtil_ValueType {
	switch n.GetToken() {
	case rhino.Token_CAST:
		return NodeUtil_GetKnownValueType(n.GetFirstChild())
	case rhino.Token_ASSIGN, rhino.Token_COMMA:
		return NodeUtil_GetKnownValueType(n.GetLastChild())
	case rhino.Token_AND, rhino.Token_OR, rhino.Token_COALESCE, rhino.Token_ASSIGN_OR, rhino.Token_ASSIGN_AND,
		rhino.Token_ASSIGN_COALESCE:
		return NodeUtil_And(NodeUtil_GetKnownValueType(n.GetFirstChild()), NodeUtil_GetKnownValueType(n.GetLastChild()))
	case rhino.Token_HOOK:
		return NodeUtil_And(NodeUtil_GetKnownValueType(n.GetSecondChild()), NodeUtil_GetKnownValueType(n.GetLastChild()))
	case rhino.Token_ADD:
		last := NodeUtil_GetKnownValueType(n.GetLastChild())
		if last == NodeUtil_ValueType_STRING {
			return NodeUtil_ValueType_STRING
		}
		first := NodeUtil_GetKnownValueType(n.GetFirstChild())
		if first == NodeUtil_ValueType_STRING {
			return NodeUtil_ValueType_STRING
		}

		// There are some pretty weird cases for object types:
		//   {} + [] === "0"
		//   [] + {} === "[object Object]"
		if first == NodeUtil_ValueType_OBJECT || last == NodeUtil_ValueType_OBJECT {
			return NodeUtil_ValueType_UNDETERMINED
		}

		if !NodeUtil_MayBeStringValueType(first) && !NodeUtil_MayBeStringValueType(last) {
			if first == NodeUtil_ValueType_BIGINT || last == NodeUtil_ValueType_BIGINT {
				// If one operand is a BigInt, then the result is a BigInt or there's a type error
				return NodeUtil_ValueType_BIGINT
			} else {
				// ADD used with compilations of null, undefined, boolean and number always result
				// in numbers.
				return NodeUtil_ValueType_NUMBER
			}
		}

		return NodeUtil_ValueType_UNDETERMINED
	case rhino.Token_ASSIGN_ADD:
		last := NodeUtil_GetKnownValueType(n.GetLastChild())
		if last == NodeUtil_ValueType_STRING {
			return NodeUtil_ValueType_STRING
		}
		return NodeUtil_ValueType_UNDETERMINED
	case rhino.Token_NAME:
		name := n.GetString()
		if name == "undefined" {
			return NodeUtil_ValueType_VOID
		}
		if name == "NaN" {
			return NodeUtil_ValueType_NUMBER
		}
		if name == "Infinity" {
			return NodeUtil_ValueType_NUMBER
		}
		return NodeUtil_ValueType_UNDETERMINED
	case rhino.Token_ASSIGN_BITOR,
		rhino.Token_ASSIGN_BITXOR,
		rhino.Token_ASSIGN_BITAND,
		rhino.Token_ASSIGN_LSH,
		rhino.Token_ASSIGN_RSH,
		rhino.Token_ASSIGN_URSH,
		rhino.Token_ASSIGN_SUB,
		rhino.Token_ASSIGN_MUL,
		rhino.Token_ASSIGN_EXPONENT,
		rhino.Token_ASSIGN_DIV,
		rhino.Token_ASSIGN_MOD:
		// assign operators could be using BIGINT or NUMBER
		if NodeUtil_GetKnownValueType(n.GetLastChild()) == NodeUtil_ValueType_BIGINT {
			return NodeUtil_ValueType_BIGINT
		} else {
			return NodeUtil_ValueType_NUMBER
		}
	case rhino.Token_BIGINT:
		return NodeUtil_ValueType_BIGINT
	case rhino.Token_BITOR, rhino.Token_BITXOR, rhino.Token_BITAND, rhino.Token_LSH, rhino.Token_RSH, rhino.Token_SUB,
		rhino.Token_MUL, rhino.Token_MOD, rhino.Token_DIV, rhino.Token_EXPONENT:
		// binary arithmetic operators could result in BIGINT or NUMBER
		first := NodeUtil_GetKnownValueType(n.GetFirstChild())
		last := NodeUtil_GetKnownValueType(n.GetLastChild())
		if first == NodeUtil_ValueType_BIGINT || last == NodeUtil_ValueType_BIGINT {
			return NodeUtil_ValueType_BIGINT
		} else {
			return NodeUtil_ValueType_NUMBER
		}
	case rhino.Token_BITNOT, rhino.Token_NEG:
		// unary negation (bitwise or arithmetic) could be using BIGINT or NUMBER
		if NodeUtil_GetKnownValueType(n.GetOnlyChild()) == NodeUtil_ValueType_BIGINT {
			return NodeUtil_ValueType_BIGINT
		} else {
			return NodeUtil_ValueType_NUMBER
		}
	case rhino.Token_INC, rhino.Token_DEC:
		// increment and decrement can only be used on variables, so we assume they're numbers
		return NodeUtil_ValueType_NUMBER
	case rhino.Token_URSH, rhino.Token_POS, rhino.Token_NUMBER:
		// unary + and unsigned right shift don't apply to bigint
		return NodeUtil_ValueType_NUMBER
	case
		// Primitives
		rhino.Token_TRUE,
		rhino.Token_FALSE,
		// Comparisons
		rhino.Token_EQ,
		rhino.Token_NE,
		rhino.Token_SHEQ,
		rhino.Token_SHNE,
		rhino.Token_LT,
		rhino.Token_GT,
		rhino.Token_LE,
		rhino.Token_GE,
		// Queries
		rhino.Token_IN,
		rhino.Token_INSTANCEOF,
		// Inversion
		rhino.Token_NOT,
		// delete operator returns a boolean.
		rhino.Token_DELPROP:
		return NodeUtil_ValueType_BOOLEAN
	case rhino.Token_TYPEOF, rhino.Token_STRINGLIT, rhino.Token_TEMPLATELIT:
		return NodeUtil_ValueType_STRING
	case rhino.Token_NULL:
		return NodeUtil_ValueType_NULL
	case rhino.Token_VOID:
		return NodeUtil_ValueType_VOID
	case rhino.Token_FUNCTION, rhino.Token_NEW, rhino.Token_ARRAYLIT, rhino.Token_OBJECTLIT, rhino.Token_REGEXP:
		return NodeUtil_ValueType_OBJECT
	default:
		jutil.CheckArgument(!n.IsTemplateLitString())
		return NodeUtil_ValueType_UNDETERMINED
	}
}

// Java: NodeUtil.and(ValueType, ValueType)
func NodeUtil_And(a NodeUtil_ValueType, b NodeUtil_ValueType) NodeUtil_ValueType {
	if a == b {
		return a
	}
	return NodeUtil_ValueType_UNDETERMINED
}

// NodeUtil_IsNumericResult returns true if the result of node evaluation is
// always a number.
//
// Java: NodeUtil.isNumericResult(Node)
func NodeUtil_IsNumericResult(n rhino.Node) bool {
	return NodeUtil_GetKnownValueType(n) == NodeUtil_ValueType_NUMBER
}

// NodeUtil_IsBigIntResult returns true if the result of node evaluation is
// always a bigint.
//
// Java: NodeUtil.isBigIntResult(Node)
func NodeUtil_IsBigIntResult(n rhino.Node) bool {
	return NodeUtil_GetKnownValueType(n) == NodeUtil_ValueType_BIGINT
}

// NodeUtil_IsBooleanResult returns whether the result of node evaluation is
// always a boolean.
//
// Java: NodeUtil.isBooleanResult(Node)
func NodeUtil_IsBooleanResult(n rhino.Node) bool {
	return NodeUtil_GetKnownValueType(n) == NodeUtil_ValueType_BOOLEAN
}

// NodeUtil_IsStringResult returns whether the result of node evaluation is
// always a string.
//
// Java: NodeUtil.isStringResult(Node)
func NodeUtil_IsStringResult(n rhino.Node) bool {
	return NodeUtil_GetKnownValueType(n) == NodeUtil_ValueType_STRING
}

// NodeUtil_IsObjectResult returns whether the result of node evaluation is
// always an object.
//
// Java: NodeUtil.isObjectResult(Node)
func NodeUtil_IsObjectResult(n rhino.Node) bool {
	return NodeUtil_GetKnownValueType(n) == NodeUtil_ValueType_OBJECT
}

// Java: NodeUtil.mayBeString(Node)
func NodeUtil_MayBeString(n rhino.Node) bool {
	return NodeUtil_MayBeStringWithUseType(n, false)
}

// NodeUtil_MayBeStringWithUseType returns if the node is possibly a string.
//
// useType: if true and the node has a primitive type, return true if that
// type is string and false otherwise. Returns whether the results is possibly
// a string.
//
// Java: NodeUtil.mayBeString(Node, boolean)
func NodeUtil_MayBeStringWithUseType(n rhino.Node, useType bool) bool {
	if useType {
		color := n.GetColor()
		if color != nil {
			if color.Equals(colors.StandardColors_STRING) {
				return true
			} else if color.Equals(colors.StandardColors_NUMBER) ||
				color.Equals(colors.StandardColors_BIGINT) ||
				color.Equals(colors.StandardColors_BOOLEAN) ||
				color.Equals(colors.StandardColors_NULL_OR_VOID) {
				return false
			}
		}
		type_ := n.GetJSType()
		if type_ != nil {
			if type_.IsStringValueType() {
				return true
			} else if type_.IsNumberValueType() ||
				type_.IsBigIntValueType() ||
				type_.IsBooleanValueType() ||
				type_.IsNullType() ||
				type_.IsVoidType() {
				return false
			}
		}
	}
	return NodeUtil_MayBeStringValueType(NodeUtil_GetKnownValueType(n))
}

// NodeUtil_MayBeStringValueType returns whether the results is possibly a
// string, this includes Objects which may implicitly be converted to a
// string.
//
// Java: NodeUtil.mayBeString(ValueType)
func NodeUtil_MayBeStringValueType(type_ NodeUtil_ValueType) bool {
	switch type_ {
	case NodeUtil_ValueType_BOOLEAN, NodeUtil_ValueType_NULL, NodeUtil_ValueType_NUMBER, NodeUtil_ValueType_BIGINT,
		NodeUtil_ValueType_VOID:
		return false
	default: // case OBJECT, STRING, UNDETERMINED
		return true
	}
}

// Java: NodeUtil.mayBeObject(Node)
func NodeUtil_MayBeObject(n rhino.Node) bool {
	return NodeUtil_MayBeObjectValueType(NodeUtil_GetKnownValueType(n))
}

// Java: NodeUtil.mayBeObject(ValueType)
func NodeUtil_MayBeObjectValueType(type_ NodeUtil_ValueType) bool {
	switch type_ {
	case NodeUtil_ValueType_BOOLEAN, NodeUtil_ValueType_NULL, NodeUtil_ValueType_NUMBER, NodeUtil_ValueType_BIGINT,
		NodeUtil_ValueType_STRING, NodeUtil_ValueType_VOID:
		return false
	default: // case OBJECT, UNDETERMINED
		return true
	}
}

// NodeUtil_IsAssociative returns true if the operator is associative. e.g.
// (a * b) * c = a * (b * c) Note: "+" is not associative because it is also
// the concatenation for strings. e.g. "a" + (1 + 2) is not "a" + 1 + 2.
// Multilication is not associative because it can include floating point
// numbers.
//
// e.g. 1e-300 * 1e300 * 1e9 does not equal 1e-300 * (1e300 * 1e9).
//
// Java: NodeUtil.isAssociative(Token)
func NodeUtil_IsAssociative(type_ rhino.Token) bool {
	switch type_ {
	case rhino.Token_AND, rhino.Token_OR, rhino.Token_COALESCE, rhino.Token_BITOR, rhino.Token_BITXOR,
		rhino.Token_BITAND:
		return true
	default:
		return false
	}
}

// NodeUtil_IsCommutative returns true if the operator is commutative. e.g.
// (a * b) * c = c * (b * a) Note 1: "+" is not commutative because it is also
// the concatenation for strings. e.g. "a" + (1 + 2) is not "a" + 1 + 2 Note
// 2: only operations on literals and pure functions are commutative.
//
// Java: NodeUtil.isCommutative(Token)
func NodeUtil_IsCommutative(type_ rhino.Token) bool {
	switch type_ {
	case rhino.Token_MUL, rhino.Token_BITOR, rhino.Token_BITXOR, rhino.Token_BITAND:
		return true
	default:
		return false
	}
}

// NodeUtil_IsAssignmentOp returns true if the operator is an assignment type
// operator. Note: The logical assignments (i.e. ASSIGN_OR, ASSIGN_AND,
// ASSIGN_COALESCE) follow short-circuiting behavior, and the RHS may not
// always be evaluated. They are still considered AssignmentOps (may be
// optimized).
//
// Java: NodeUtil.isAssignmentOp(Node)
func NodeUtil_IsAssignmentOp(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_ASSIGN,
		rhino.Token_ASSIGN_BITOR,
		rhino.Token_ASSIGN_BITXOR,
		rhino.Token_ASSIGN_BITAND,
		rhino.Token_ASSIGN_LSH,
		rhino.Token_ASSIGN_RSH,
		rhino.Token_ASSIGN_URSH,
		rhino.Token_ASSIGN_ADD,
		rhino.Token_ASSIGN_SUB,
		rhino.Token_ASSIGN_MUL,
		rhino.Token_ASSIGN_EXPONENT,
		rhino.Token_ASSIGN_DIV,
		rhino.Token_ASSIGN_MOD,
		rhino.Token_ASSIGN_OR,
		rhino.Token_ASSIGN_AND,
		rhino.Token_ASSIGN_COALESCE:
		return true
	default:
		return false
	}
}

// NodeUtil_IsLogicalAssignmentOp returns true if the operator is a logical
// assignment type operator.
//
// Java: NodeUtil.isLogicalAssignmentOp(Node)
func NodeUtil_IsLogicalAssignmentOp(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_ASSIGN_OR, rhino.Token_ASSIGN_AND, rhino.Token_ASSIGN_COALESCE:
		return true
	default:
		return false
	}
}

// Java: NodeUtil.isCompoundAssignmentOp(Node)
func NodeUtil_IsCompoundAssignmentOp(n rhino.Node) bool {
	return NodeUtil_IsAssignmentOp(n) && !n.IsAssign()
}

// Java: NodeUtil.getOpFromAssignmentOp(Node)
func NodeUtil_GetOpFromAssignmentOp(n rhino.Node) rhino.Token {
	switch n.GetToken() {
	case rhino.Token_ASSIGN_BITOR:
		return rhino.Token_BITOR
	case rhino.Token_ASSIGN_BITXOR:
		return rhino.Token_BITXOR
	case rhino.Token_ASSIGN_BITAND:
		return rhino.Token_BITAND
	case rhino.Token_ASSIGN_LSH:
		return rhino.Token_LSH
	case rhino.Token_ASSIGN_RSH:
		return rhino.Token_RSH
	case rhino.Token_ASSIGN_URSH:
		return rhino.Token_URSH
	case rhino.Token_ASSIGN_ADD:
		return rhino.Token_ADD
	case rhino.Token_ASSIGN_SUB:
		return rhino.Token_SUB
	case rhino.Token_ASSIGN_MUL:
		return rhino.Token_MUL
	case rhino.Token_ASSIGN_EXPONENT:
		return rhino.Token_EXPONENT
	case rhino.Token_ASSIGN_DIV:
		return rhino.Token_DIV
	case rhino.Token_ASSIGN_MOD:
		return rhino.Token_MOD
	case rhino.Token_ASSIGN_OR:
		return rhino.Token_OR
	case rhino.Token_ASSIGN_AND:
		return rhino.Token_AND
	case rhino.Token_ASSIGN_COALESCE:
		return rhino.Token_COALESCE
	default:
		panic(jutil.NewIllegalArgumentExceptionMsg("Not an assignment op:" + n.String()))
	}
}

// NodeUtil_GetEnclosingType gets the closest ancestor to the given node of
// the provided type.
//
// Java: NodeUtil.getEnclosingType(Node, Token)
func NodeUtil_GetEnclosingType(n rhino.Node, type_ rhino.Token) rhino.Node {
	return NodeUtil_GetEnclosingNode(n, func(n1 rhino.Node) bool { return n1.GetToken() == type_ })
}

// Java: NodeUtil.getEnclosingNonArrowFunction(Node)
func NodeUtil_GetEnclosingNonArrowFunction(n rhino.Node) rhino.Node {
	return NodeUtil_GetEnclosingNode(n, NodeUtil_IsNonArrowFunction)
}

// NodeUtil_GetEnclosingClass finds the class containing the given node.
//
// Java: NodeUtil.getEnclosingClass(Node)
func NodeUtil_GetEnclosingClass(n rhino.Node) rhino.Node {
	return NodeUtil_GetEnclosingNode(n, (rhino.Node).IsClass)
}

// Java: NodeUtil.getEnclosingModuleIfPresent(Node)
func NodeUtil_GetEnclosingModuleIfPresent(n rhino.Node) rhino.Node {
	return NodeUtil_GetEnclosingNode(n, (rhino.Node).IsModuleBody)
}

// NodeUtil_GetEnclosingFunction finds the function containing the given
// node.
//
// Java: NodeUtil.getEnclosingFunction(Node)
func NodeUtil_GetEnclosingFunction(n rhino.Node) rhino.Node {
	return NodeUtil_GetEnclosingNode(n, (rhino.Node).IsFunction)
}

// NodeUtil_GetEnclosingScript finds the script containing the given node.
//
// Java: NodeUtil.getEnclosingScript(Node)
func NodeUtil_GetEnclosingScript(n rhino.Node) rhino.Node {
	return NodeUtil_GetEnclosingNode(n, (rhino.Node).IsScript)
}

// NodeUtil_GetEnclosingBlock finds the block containing the given node.
//
// Java: NodeUtil.getEnclosingBlock(Node)
func NodeUtil_GetEnclosingBlock(n rhino.Node) rhino.Node {
	return NodeUtil_GetEnclosingNode(n, (rhino.Node).IsBlock)
}

// Java: NodeUtil.getEnclosingBlockScopeRoot(Node)
func NodeUtil_GetEnclosingBlockScopeRoot(n rhino.Node) rhino.Node {
	return NodeUtil_GetEnclosingNode(n, NodeUtil_CreatesBlockScope)
}

// Java: NodeUtil.getEnclosingScopeRoot(Node)
func NodeUtil_GetEnclosingScopeRoot(n rhino.Node) rhino.Node {
	return NodeUtil_GetEnclosingNode(n, NodeUtil_CreatesScope)
}

// NodeUtil_GetEnclosingHoistScopeRoot returns the nearest enclosing hoist
// scope root node, null for the global scope. There are currently 4 such
// roots to consider: the global, module, function, and class static blocks.
//
// Java: NodeUtil.getEnclosingHoistScopeRoot(Node)
func NodeUtil_GetEnclosingHoistScopeRoot(n rhino.Node) rhino.Node {
	return NodeUtil_GetEnclosingNode(n, NodeUtil_IsHoistScopeRoot)
}

// Java: NodeUtil.isHoistScopeRoot(Node)
func NodeUtil_IsHoistScopeRoot(n rhino.Node) bool {
	// The "global" hoist scope has multiple meaning so
	return n.IsFunction() || n.IsModuleBody() || NodeUtil_IsClassStaticBlock(n)
}

// Java: NodeUtil.isInFunction(Node)
func NodeUtil_IsInFunction(n rhino.Node) bool {
	return !NodeUtil_GetEnclosingFunction(n).IsNil()
}

// Java: NodeUtil.getEnclosingStatement(Node)
func NodeUtil_GetEnclosingStatement(n rhino.Node) rhino.Node {
	return NodeUtil_GetEnclosingNode(n, NodeUtil_IsStatement)
}

// Java: NodeUtil.getEnclosingNode(Node, Predicate<Node>)
func NodeUtil_GetEnclosingNode(n rhino.Node, pred func(rhino.Node) bool) rhino.Node {
	curr := n
	for !curr.IsNil() && !pred(curr) {
		curr = curr.GetParent()
	}
	return curr
}

// NodeUtil_GetFirstPropMatchingKey returns the first property in the objlit
// or class members, that matches the key.
//
// Java: NodeUtil.getFirstPropMatchingKey(Node, String)
func NodeUtil_GetFirstPropMatchingKey(n rhino.Node, keyName string) rhino.Node {
	jutil.CheckState(n.IsObjectLit() || n.IsClassMembers())
	for keyNode := n.GetFirstChild(); !keyNode.IsNil(); keyNode = keyNode.GetNext() {
		if (keyNode.IsStringKey() || keyNode.IsMemberFunctionDef()) &&
			keyNode.GetString() == keyName {
			return keyNode.GetFirstChild()
		} else if keyNode.IsComputedProp() &&
			keyNode.GetFirstChild().IsStringLit() &&
			keyNode.GetFirstChild().GetString() == keyName {
			return keyNode.GetLastChild()
		}
	}
	return rhino.NilNode
}

// NodeUtil_GetFirstGetterMatchingKey returns the first getter in the class
// members that matches the key.
//
// Java: NodeUtil.getFirstGetterMatchingKey(Node, String)
func NodeUtil_GetFirstGetterMatchingKey(n rhino.Node, keyName string) rhino.Node {
	jutil.CheckStateMsg1(n.IsClassMembers() || n.IsObjectLit(), "%s", n)
	for keyNode := n.GetFirstChild(); !keyNode.IsNil(); keyNode = keyNode.GetNext() {
		if keyNode.IsGetterDef() && keyNode.GetString() == keyName {
			return keyNode
		}
	}
	return rhino.NilNode
}

// NodeUtil_ReferencesOwnReceiver returns true if this function references
// its receiver object.
//
// We define this function in terms of "receiver" rather than specific
// syntax, however internally we search for `this` and `super`.
//
// Arrow functions return false. They do not have their own receiver, but
// rather capture the receiver of the enclosing scope.
//
// Java: NodeUtil.referencesOwnReceiver(Node)
func NodeUtil_ReferencesOwnReceiver(fn rhino.Node) bool {
	jutil.CheckState(fn.IsFunction())
	if fn.IsArrowFunction() {
		return false
	}

	return NodeUtil_ReferencesEnclosingReceiver(NodeUtil_GetFunctionParameters(fn)) ||
		NodeUtil_ReferencesEnclosingReceiver(NodeUtil_GetFunctionBody(fn))
}

// NodeUtil_ReferencesEnclosingReceiver returns true if this subtree
// references the receiver object from its enclosing function scope.
//
// We define this function in terms of "receiver" rather than specific
// syntax, however internally we search for `this` and `super`.
//
// Arrow functions may return true. They capture the receiver of the
// enclosing scope, rather than having their own.
//
// Java: NodeUtil.referencesEnclosingReceiver(Node)
func NodeUtil_ReferencesEnclosingReceiver(n rhino.Node) bool {
	return NodeUtil_Has(n, func(c rhino.Node) bool { return c.IsThis() || c.IsSuper() },
		NodeUtil_MATCH_ANYTHING_BUT_NON_ARROW_FUNCTION)
}

// NodeUtil_ReferencesSuper returns true if the current scope contains
// references to the 'super' keyword. Note that if there are classes declared
// inside the current class, super calls which reference those classes are
// not reported.
//
// Java: NodeUtil.referencesSuper(Node)
func NodeUtil_ReferencesSuper(n rhino.Node) bool {
	curr := n.GetFirstChild()
	for !curr.IsNil() {
		if NodeUtil_Has(curr, (rhino.Node).IsSuper, func(node rhino.Node) bool { return !node.IsClass() }) {
			return true
		}
		curr = curr.GetNext()
	}
	return false
}

// NodeUtil_IsNormalOrOptChainGet: is this a GETPROP, OPTCHAIN_GETPROP,
// GETELEM, or OPTCHAIN_GETELEM?
//
// Java: NodeUtil.isNormalOrOptChainGet(Node)
func NodeUtil_IsNormalOrOptChainGet(n rhino.Node) bool {
	return NodeUtil_IsNormalGet(n) || NodeUtil_IsOptChainGet(n)
}

// NodeUtil_IsNormalOrOptChainGetProp: is this a GETPROP or OPTCHAIN_GETPROP?
//
// Java: NodeUtil.isNormalOrOptChainGetProp(Node)
func NodeUtil_IsNormalOrOptChainGetProp(n rhino.Node) bool {
	return n.IsGetProp() || n.IsOptChainGetProp()
}

// NodeUtil_IsNormalOrOptChainCall: is this a CALL or OPTCHAIN_CALL?
//
// Java: NodeUtil.isNormalOrOptChainCall(Node)
func NodeUtil_IsNormalOrOptChainCall(n rhino.Node) bool {
	return n.IsCall() || n.IsOptChainCall()
}

// NodeUtil_IsNormalGet: is this a GETPROP or GETELEM node?
//
// Java: NodeUtil.isNormalGet(Node)
func NodeUtil_IsNormalGet(n rhino.Node) bool {
	return n.IsGetProp() || n.IsGetElem()
}

// NodeUtil_IsOptChainGet: is this an OPTCHAIN_GETPROP or OPTCHAIN_GETELEM
// node?
//
// Java: NodeUtil.isOptChainGet(Node)
func NodeUtil_IsOptChainGet(n rhino.Node) bool {
	return n.IsOptChainGetProp() || n.IsOptChainGetElem()
}

// NodeUtil_IsOptChainNode: is this a OPTCHAIN_GETPROP, OPTCHAIN_GETELEM,
// OPTCHAIN_CALL node?
//
// Java: NodeUtil.isOptChainNode(Node)
func NodeUtil_IsOptChainNode(n rhino.Node) bool {
	return n.IsOptChainGetProp() || n.IsOptChainGetElem() || n.IsOptChainCall()
}

// NodeUtil_GetStartOfOptChainSegment finds the start of the optional chain.
// E.g Find the `a?. ...` node in `a?.b.c.d` given any other node.
//
// n is a node in an optional chain; returns the start of the optional chain
// that n is part of.
//
// Java: NodeUtil.getStartOfOptChainSegment(Node)
func NodeUtil_GetStartOfOptChainSegment(n rhino.Node) rhino.Node {
	jutil.CheckStateMsg1(NodeUtil_IsOptChainNode(n), "%s", n)
	if n.IsOptionalChainStart() {
		return n
	}
	return NodeUtil_GetStartOfOptChainSegment(n.GetFirstChild())
}

// NodeUtil_GetEndOfOptChainSegment finds the end of an optional chain
// segment.
//
// Each `?.` ends one segment and starts another.
//
// Examples
//
//	a?.b.c.d   // end is the node with children `a?.b.c` and `d`
//	a?.b.c?.d  // given a?.b end is the node with children `a?.b` and `c`
//	a?.b.c?.d  // given a?.b.c end is the node with children `a?.b.c` and `d`
//	(a?.b.c).d // given a?.b end is the node with children `a?.b` and `c`
//
// n is a node in an optional chain; returns the end of the optional chain
// that n is part of.
//
// Java: NodeUtil.getEndOfOptChainSegment(Node)
func NodeUtil_GetEndOfOptChainSegment(n rhino.Node) rhino.Node {
	jutil.CheckStateMsg1(NodeUtil_IsOptChainNode(n), "%s", n)
	if NodeUtil_IsEndOfOptChainSegment(n) {
		return n
	} else {
		return NodeUtil_GetEndOfOptChainSegment(n.GetParent())
	}
}

// NodeUtil_IsEndOfFullOptChain: is this node the final node of a full
// optional chain?
//
// e.g. for `a?.b.c?.d` this method returns true only for the Node with
// children `a?.b.c` and `d`. That node is the end of the whole chain, and
// also represents the whole chain in the AST.
//
// Java: NodeUtil.isEndOfFullOptChain(Node)
func NodeUtil_IsEndOfFullOptChain(n rhino.Node) bool {
	if NodeUtil_IsOptChainNode(n) {
		parent := n.GetParent()
		// the chain continues if this is the first child of another optional chain node
		return !(NodeUtil_IsOptChainNode(parent) && n.IsFirstChildOf(parent))
	} else {
		// not even an optional chain node
		return false
	}
}

// NodeUtil_IsEndOfOptChainSegment: is this node the end of an optional chain
// segment?
//
// Each `?.` begins a new segment and ends the previous one, if any. The end
// of the whole chain is also the end of its final segment.
//
// Java: NodeUtil.isEndOfOptChainSegment(Node)
func NodeUtil_IsEndOfOptChainSegment(n rhino.Node) bool {
	if !NodeUtil_IsOptChainNode(n) {
		return false
	} else {
		parent := n.GetParent()
		// Check for null so this method will work for a disconnected Node.
		if !parent.IsNil() && n.IsFirstChildOf(parent) && NodeUtil_IsOptChainNode(parent) {
			// The parent is a continuation of this optional chain.
			// If it starts a new segment, then this node is the end of a segment
			return parent.IsOptionalChainStart()
		} else {
			// The parent doesn't continue this node's optional chain, though it might be part of a
			// different one.
			// e.g. in `a?.(x?.y.z)` the parent of `x?.y.z` is part of a different optional chain.
			return true
		}
	}
}

// NodeUtil_ConvertToNonOptionalChainSegment: given the end of an optional
// chain segment changes all nodes from the end down to the start into
// non-optional nodes. e.g `({a})?.a.b.c.d()?.x.y.z` gets converted to
// `({a}).a.b.c.d()?.x.y.z` when the passed in endOfOptChainSegment is
// `({a})?.a.b.c.d()`.
//
// Java: NodeUtil.convertToNonOptionalChainSegment(Node)
func NodeUtil_ConvertToNonOptionalChainSegment(endOfOptChainSegment rhino.Node) {
	start := NodeUtil_GetStartOfOptChainSegment(endOfOptChainSegment)
	NodeUtil_ConvertToNonOptionalChainSegmentDownTo(endOfOptChainSegment, start.GetFirstChild())
}

// NodeUtil_ConvertToNonOptionalChainSegmentDownTo: given the end of an
// optional chain segment, changes all nodes from endOfOptChainSegment down
// to (but excluding) stopNode into non-optional nodes.
//
// Java: NodeUtil.convertToNonOptionalChainSegmentDownTo(Node, Node)
func NodeUtil_ConvertToNonOptionalChainSegmentDownTo(endOfOptChainSegment rhino.Node, stopNode rhino.Node) {
	jutil.CheckArgumentMsg1(NodeUtil_IsEndOfOptChainSegment(endOfOptChainSegment), "%s", endOfOptChainSegment)
	var segmentNodes []rhino.Node
	segmentNode := endOfOptChainSegment
	for !segmentNode.IsNil() && segmentNode != stopNode {
		if NodeUtil_IsOptChainNode(segmentNode) {
			segmentNodes = append(segmentNodes, segmentNode)
		}
		segmentNode = segmentNode.GetFirstChild()
	}
	for _, n := range segmentNodes {
		n.SetIsOptionalChainStart(false)
		n.SetToken(NodeUtil_GetNonOptChainToken(n.GetToken()))
	}
}

// Java: NodeUtil.getNonOptChainToken(Token)
func NodeUtil_GetNonOptChainToken(optChainToken rhino.Token) rhino.Token {
	switch optChainToken {
	case rhino.Token_OPTCHAIN_CALL:
		return rhino.Token_CALL
	case rhino.Token_OPTCHAIN_GETELEM:
		return rhino.Token_GETELEM
	case rhino.Token_OPTCHAIN_GETPROP:
		return rhino.Token_GETPROP
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Should be an OPTCHAIN token: " + optChainToken.String()))
	}
}

// NodeUtil_IsBlockScopedDeclaration: is this node a block-scoped
// declaration? Checks for let, const, class, or block-scoped function
// declarations.
//
// Returns true if n is a block-scoped declaration statement.
//
// Java: NodeUtil.isBlockScopedDeclaration(Node)
func NodeUtil_IsBlockScopedDeclaration(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_LET, rhino.Token_CONST, rhino.Token_CATCH:
		return true
	case rhino.Token_CLASS:
		return NodeUtil_IsClassDeclaration(n)
	case rhino.Token_FUNCTION:
		return NodeUtil_IsBlockScopedFunctionDeclaration(n)
	default:
		return false
	}
}

// NodeUtil_IsNameDeclaration: is this node a name declaration?
//
// Returns true if n is VAR, LET or CONST.
//
// Java: NodeUtil.isNameDeclaration(Node)
func NodeUtil_IsNameDeclaration(n rhino.Node) bool {
	return !n.IsNil() && (n.IsVar() || n.IsLet() || n.IsConst())
}

// NodeUtil_IsDestructuringDeclaration returns true if n is a VAR, LET or
// CONST that contains a destructuring pattern.
//
// Java: NodeUtil.isDestructuringDeclaration(Node)
func NodeUtil_IsDestructuringDeclaration(n rhino.Node) bool {
	if NodeUtil_IsNameDeclaration(n) {
		for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			if c.IsDestructuringLhs() {
				return true
			}
		}
	}
	return false
}

// NodeUtil_GetAssignedValue: for an assignment or variable declaration get
// the assigned value.
//
// Returns the value node representing the new value.
//
// Java: NodeUtil.getAssignedValue(Node)
func NodeUtil_GetAssignedValue(n rhino.Node) rhino.Node {
	jutil.CheckStateMsg1(n.IsName() || n.IsGetProp(), "%s", n)
	parent := n.GetParent()
	if NodeUtil_IsNameDeclaration(parent) {
		return n.GetFirstChild()
	} else if parent.IsAssign() && parent.GetFirstChild() == n {
		return n.GetNext()
	} else {
		return rhino.NilNode
	}
}

// NodeUtil_IsExprAssign: is this node an assignment expression statement?
//
// Returns true if n is EXPR_RESULT and n's first child is ASSIGN.
//
// Java: NodeUtil.isExprAssign(Node)
func NodeUtil_IsExprAssign(n rhino.Node) bool {
	return n.IsExprResult() && n.GetFirstChild().IsAssign()
}

// NodeUtil_IsExprCall: is this node a call expression statement?
//
// Returns true if n is EXPR_RESULT and n's first child is CALL.
//
// Java: NodeUtil.isExprCall(Node)
func NodeUtil_IsExprCall(n rhino.Node) bool {
	return n.IsExprResult() && n.GetFirstChild().IsCall()
}

// Java: NodeUtil.isNonArrowFunction(Node)
func NodeUtil_IsNonArrowFunction(n rhino.Node) bool {
	return n.IsFunction() && !n.IsArrowFunction()
}

// Java: NodeUtil.isEnhancedFor(Node)
func NodeUtil_IsEnhancedFor(n rhino.Node) bool {
	return n.IsForOf() || n.IsForAwaitOf() || n.IsForIn()
}

// Java: NodeUtil.isAnyFor(Node)
func NodeUtil_IsAnyFor(n rhino.Node) bool {
	return n.IsVanillaFor() || n.IsForIn() || n.IsForOf() || n.IsForAwaitOf()
}

// NodeUtil_IsLoopStructure determines whether the given node is a FOR, DO,
// or WHILE node.
//
// Java: NodeUtil.isLoopStructure(Node)
func NodeUtil_IsLoopStructure(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF, rhino.Token_DO,
		rhino.Token_WHILE:
		return true
	default:
		return false
	}
}

// NodeUtil_GetLoopCodeBlock: n is the node to inspect. If the node, is a
// FOR, WHILE, or DO, it returns the node for the code BLOCK, null otherwise.
//
// Java: NodeUtil.getLoopCodeBlock(Node)
func NodeUtil_GetLoopCodeBlock(n rhino.Node) rhino.Node {
	switch n.GetToken() {
	case rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF, rhino.Token_WHILE:
		return n.GetLastChild()
	case rhino.Token_DO:
		return n.GetFirstChild()
	default:
		return rhino.NilNode
	}
}

// NodeUtil_IsWithinLoop returns whether the specified node has a loop parent
// that is within the current scope.
//
// Java: NodeUtil.isWithinLoop(Node)
func NodeUtil_IsWithinLoop(n rhino.Node) bool {
	for parent := range n.GetAncestors().All() {
		if NodeUtil_IsLoopStructure(parent) {
			return true
		}

		if parent.IsFunction() {
			break
		}
	}
	return false
}

// NodeUtil_IsControlStructure determines whether the given node is a FOR,
// DO, WHILE, WITH, or IF node.
//
// Java: NodeUtil.isControlStructure(Node)
func NodeUtil_IsControlStructure(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_FOR,
		rhino.Token_FOR_IN,
		rhino.Token_FOR_OF,
		rhino.Token_FOR_AWAIT_OF,
		rhino.Token_DO,
		rhino.Token_WHILE,
		rhino.Token_WITH,
		rhino.Token_IF,
		rhino.Token_LABEL,
		rhino.Token_TRY,
		rhino.Token_CATCH,
		rhino.Token_SWITCH,
		rhino.Token_CASE,
		rhino.Token_DEFAULT_CASE:
		return true
	default:
		return false
	}
}

// NodeUtil_IsControlStructureCodeBlock determines whether the given node is
// code node for FOR, DO, WHILE, WITH, or IF node.
//
// Java: NodeUtil.isControlStructureCodeBlock(Node, Node)
func NodeUtil_IsControlStructureCodeBlock(parent rhino.Node, n rhino.Node) bool {
	switch parent.GetToken() {
	case rhino.Token_DO:
		return parent.GetFirstChild() == n
	case rhino.Token_TRY:
		return parent.GetFirstChild() == n || parent.GetLastChild() == n
	case rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF, rhino.Token_WHILE,
		rhino.Token_LABEL, rhino.Token_WITH, rhino.Token_CATCH:
		return parent.GetLastChild() == n
	case rhino.Token_IF, rhino.Token_SWITCH, rhino.Token_CASE:
		return parent.GetFirstChild() != n // exclude the condition
	case rhino.Token_DEFAULT_CASE:
		return true
	default:
		jutil.CheckStateMsg1(NodeUtil_IsControlStructure(parent), "%s", parent)
		return false
	}
}

// NodeUtil_GetConditionExpression gets the condition of an ON_TRUE /
// ON_FALSE CFG edge.
//
// n is a node with an outgoing conditional CFG edge; returns the condition
// node or null if the condition is not obviously a node.
//
// Java: NodeUtil.getConditionExpression(Node)
func NodeUtil_GetConditionExpression(n rhino.Node) rhino.Node {
	switch n.GetToken() {
	case rhino.Token_IF, rhino.Token_WHILE:
		return n.GetFirstChild()
	case rhino.Token_DO:
		return n.GetLastChild()
	case rhino.Token_FOR:
		return n.GetSecondChild()
	case rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF, rhino.Token_CASE:
		return rhino.NilNode
	default:
		panic(jutil.NewIllegalArgumentExceptionMsg(n.String() + " does not have a condition."))
	}
}

// NodeUtil_IsStatementBlock returns whether the node is of a type that
// contain other statements.
//
// Java: NodeUtil.isStatementBlock(Node)
func NodeUtil_IsStatementBlock(n rhino.Node) bool {
	return n.IsRoot() || n.IsScript() || n.IsBlock() || n.IsModuleBody()
}

// NodeUtil_CreatesBlockScope: a block scope is created by a non-synthetic
// block node, a for loop node, or a for-of loop node.
//
// Note: for functions, we use two separate scopes for parameters and
// declarations in the body. We need to make sure default parameters cannot
// reference var / function declarations in the body.
//
// Returns whether the node creates a block scope.
//
// Java: NodeUtil.createsBlockScope(Node)
func NodeUtil_CreatesBlockScope(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_BLOCK:
		parent := n.GetParent()
		if !parent.IsNil() && (parent.IsBlock() || parent.IsScript()) && n.IsSyntheticBlock() {
			// Don't create block scope for synthetic blocks.
			return false
		}
		// Don't create block scope for switch cases or catch blocks.
		return !parent.IsNil() && !NodeUtil_IsSwitchCase(parent) && !parent.IsCatch()
	case rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF, rhino.Token_SWITCH_BODY,
		rhino.Token_CLASS:
		return true
	default:
		return false
	}
}

// Java: NodeUtil.createsScope(Node)
func NodeUtil_CreatesScope(n rhino.Node) bool {
	return NodeUtil_CreatesBlockScope(n) ||
		n.IsFunction() ||
		n.IsModuleBody() ||
		n.IsMemberFieldDef() ||
		n.IsComputedFieldDef() ||
		// The ROOT nodes that are the root of the externs tree or main JS tree do not
		// create scopes. The parent of those two, which is the root of the entire AST and
		// therefore has no parent, is the only ROOT node that creates a scope.
		(n.IsRoot() && n.GetParent().IsNil())
}

// Java: NodeUtil.isValidCfgRoot(Node)
func NodeUtil_IsValidCfgRoot(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_FUNCTION, rhino.Token_SCRIPT, rhino.Token_MODULE_BODY, rhino.Token_ROOT:
		return true
	default:
		return NodeUtil_IsClassStaticBlock(n)
	}
}

// NodeUtil_IsStatement returns whether the node is used as a statement.
//
// Java: NodeUtil.isStatement(Node)
func NodeUtil_IsStatement(n rhino.Node) bool {
	return !n.IsModuleBody() && !n.IsScript() && !n.IsRoot() && NodeUtil_IsStatementParent(n.GetParent())
}

// Java: NodeUtil.isStatementParent(Node)
func NodeUtil_IsStatementParent(parent rhino.Node) bool {
	// It is not possible to determine definitely if a node is a statement
	// or not if it is not part of the AST.  A FUNCTION node can be
	// either part of an expression or a statement.
	switch parent.GetToken() {
	case rhino.Token_SCRIPT,
		rhino.Token_MODULE_BODY,
		rhino.Token_BLOCK,
		rhino.Token_LABEL,
		rhino.Token_NAMESPACE_ELEMENTS,
		rhino.Token_INTERFACE_MEMBERS:
		return true
	default:
		return false
	}
}

// Java: NodeUtil.isDeclarationParent(Node)
func NodeUtil_IsDeclarationParent(parent rhino.Node) bool {
	switch parent.GetToken() {
	case rhino.Token_DECLARE, rhino.Token_EXPORT:
		return true
	default:
		return NodeUtil_IsStatementParent(parent)
	}
}

// NodeUtil_IsSwitchCase: whether the node is part of a switch statement.
//
// Java: NodeUtil.isSwitchCase(Node)
func NodeUtil_IsSwitchCase(n rhino.Node) bool {
	return n.IsCase() || n.IsDefaultCase()
}

// NodeUtil_IsReferenceName returns whether the node is a reference to a
// variable, function, class or function parameter (not a label or an empty
// function expression name).
//
// Java: NodeUtil.isReferenceName(Node)
func NodeUtil_IsReferenceName(n rhino.Node) bool {
	return n.IsName() && n.GetString() != ""
}

// NodeUtil_IsNonlocalModuleExportName returns whether the given name in an
// import or export spec is not defined within the module, but is an exported
// name from this or another module.
//
// Examples include `nonlocal` in:
//
//   - export {a as nonlocal};
//   - import {nonlocal} from './foo.js';
//   - import {nonlocal as a} from './foo.js';
//   - export {nonlocal as a} from './foo.js';
//   - export {a as nonlocal} from './foo.js';
//
// n is a NAME node.
//
// Java: NodeUtil.isNonlocalModuleExportName(Node)
func NodeUtil_IsNonlocalModuleExportName(n rhino.Node) bool {
	jutil.CheckArgumentMsg1(n.IsName(), "%s", n)
	parent := n.GetParent()
	if parent.IsImportSpec() && n.IsFirstChildOf(parent) {
		// import {nonlocal as x} from './foo.js'
		return true
	} else if parent.IsExportSpec() {
		if n.IsFirstChildOf(parent) {
			// export {nonlocal as b} from './foo.js';
			return NodeUtil_IsExportFrom(parent.GetGrandparent())
		} else {
			// export {local as nonlocal};
			return true
		}
	}
	return false
}

// NodeUtil_IsTryFinallyNode: whether the child node is the FINALLY block of
// a try.
//
// Java: NodeUtil.isTryFinallyNode(Node, Node)
func NodeUtil_IsTryFinallyNode(parent rhino.Node, child rhino.Node) bool {
	return parent.IsTry() && parent.HasXChildren(3) && child == parent.GetLastChild()
}

// NodeUtil_IsTryCatchNodeContainer: whether the node is a CATCH container
// BLOCK.
//
// Java: NodeUtil.isTryCatchNodeContainer(Node)
func NodeUtil_IsTryCatchNodeContainer(n rhino.Node) bool {
	parent := n.GetParent()
	return parent.IsTry() && parent.GetSecondChild() == n
}

// Java: NodeUtil.isInSyntheticScript(Node)
//
// TODO(tbreisacher): Add a method for detecting nodes injected as runtime
// libraries.
func NodeUtil_IsInSyntheticScript(n rhino.Node) bool {
	sourceFileName := n.GetSourceFileName()
	return sourceFileName != nil &&
		(strings.HasPrefix(*sourceFileName, " [synthetic:") ||
			strings.HasPrefix(*sourceFileName, AbstractCompiler_RUNTIME_LIB_DIR))
}

// NodeUtil_DeleteNode permanently deletes the given node from the AST, as
// well as report the related AST changes/deletions to the given compiler.
//
// Java: NodeUtil.deleteNode(Node, AbstractCompiler)
func NodeUtil_DeleteNode(n rhino.Node, compiler AbstractCompiler) {
	parent := n.GetParent()
	NodeUtil_MarkFunctionsDeleted(n, compiler)
	n.Detach()
	compiler.ReportChangeToEnclosingScope(parent)
}

// NodeUtil_DeleteFunctionCall permanently deletes the given call from the
// AST while maintaining a valid node structure, as well as report the
// related AST changes to the given compiler. In some cases, this is done by
// deleting the parent from the AST and is come cases expression is replaced
// by undefined.
//
// Java: NodeUtil.deleteFunctionCall(Node, AbstractCompiler)
func NodeUtil_DeleteFunctionCall(n rhino.Node, compiler AbstractCompiler) {
	jutil.CheckState(n.IsCall())

	parent := n.GetParent()
	if parent.IsExprResult() {
		grandParent := parent.GetParent()
		parent.Detach()
		parent = grandParent
	} else {
		// Seems like part of more complex expression, fallback to replacing with no-op.
		n.ReplaceWith(NodeUtil_NewUndefinedNode(n))
	}

	NodeUtil_MarkFunctionsDeleted(n, compiler)
	compiler.ReportChangeToEnclosingScope(parent)
}

// NodeUtil_DeleteChildren permanently deletes all the children of the given
// node, including reporting changes.
//
// Java: NodeUtil.deleteChildren(Node, AbstractCompiler)
func NodeUtil_DeleteChildren(n rhino.Node, compiler AbstractCompiler) {
	for n.HasChildren() {
		NodeUtil_DeleteNode(n.GetFirstChild(), compiler)
	}
}

// NodeUtil_RemoveChild safely removes children while maintaining a valid
// node structure. In some cases, this is done by removing the parent from
// the AST as well.
//
// Java: NodeUtil.removeChild(Node, Node)
func NodeUtil_RemoveChild(parent rhino.Node, node rhino.Node) {
	if NodeUtil_IsTryFinallyNode(parent, node) {
		if NodeUtil_HasCatchHandler(NodeUtil_GetCatchBlock(parent)) {
			// A finally can only be removed if there is a catch.
			node.Detach()
		} else {
			// Otherwise, only its children can be removed.
			node.DetachChildren()
		}
	} else if node.IsCatch() {
		// The CATCH can can only be removed if there is a finally clause.
		tryNode := node.GetGrandparent()
		jutil.CheckState(NodeUtil_HasFinally(tryNode))
		node.Detach()
	} else if NodeUtil_IsTryCatchNodeContainer(node) {
		// The container node itself can't be removed, but the contained CATCH
		// can if there is a 'finally' clause
		tryNode := node.GetParent()
		jutil.CheckState(NodeUtil_HasFinally(tryNode))
		node.DetachChildren()
	} else if node.IsBlock() {
		// Simply empty the block.  This maintains source location and
		// "synthetic"-ness.
		node.DetachChildren()
	} else if NodeUtil_IsStatementBlock(parent) || NodeUtil_IsSwitchCase(node) || node.IsMemberFunctionDef() {
		// A statement in a block or a member function can simply be removed
		node.Detach()
	} else if NodeUtil_IsNameDeclaration(parent) || parent.IsExprResult() {
		if parent.HasMoreThanOneChild() {
			node.Detach()
		} else {
			// Remove the node from the parent, so it can be reused.
			node.Detach()
			// This would leave an empty VAR, remove the VAR itself.
			NodeUtil_RemoveChild(parent.GetParent(), parent)
		}
	} else if parent.IsLabel() && node == parent.GetLastChild() {
		// Remove the node from the parent, so it can be reused.
		node.Detach()
		// A LABEL without children can not be referred to, remove it.
		NodeUtil_RemoveChild(parent.GetParent(), parent)
	} else if parent.IsVanillaFor() {
		// Only Token.FOR can have an Token.EMPTY other control structure
		// need something for the condition. Others need to be replaced
		// or the structure removed.
		node.ReplaceWith(rhino.IR_Empty())
	} else if parent.IsObjectPattern() {
		// Remove the name from the object pattern
		node.Detach()
	} else if parent.IsArrayPattern() {
		if node == parent.GetLastChild() {
			node.Detach()
		} else {
			node.ReplaceWith(rhino.IR_Empty())
		}
	} else if parent.IsDestructuringLhs() {
		// Destructuring is empty so we should remove the node
		node.Detach()
		if parent.GetParent().HasChildren() {
			// removing the destructuring could leave an empty variable declaration node, so we would
			// want to remove it from the AST
			NodeUtil_RemoveChild(parent.GetParent(), parent)
		}
	} else if parent.IsRest() {
		// Rest params can only ever have one child node
		parent.Detach()
	} else if parent.IsParamList() {
		node.Detach()
	} else if parent.IsImport() {
		// An import node must always have three child nodes. Only the first can be safely removed.
		if node == parent.GetFirstChild() {
			node.ReplaceWith(rhino.IR_Empty())
		} else {
			panic(jutil.NewIllegalStateExceptionMsg(
				"Invalid attempt to remove: " + node.String() + " from " + parent.String()))
		}
	} else {
		panic(jutil.NewIllegalStateExceptionMsg(
			"Invalid attempt to remove node: " + node.String() + " of " + parent.String()))
	}
}

// NodeUtil_ReplaceDeclarationChild replaces the child of a var/let/const
// declaration (usually a name) with a new statement. Preserves the order of
// side effects for all the other declaration children.
//
// declChild is the name node to be replaced; newStatement is the statement to
// replace with.
//
// Java: NodeUtil.replaceDeclarationChild(Node, Node)
func NodeUtil_ReplaceDeclarationChild(declChild rhino.Node, newStatement rhino.Node) {
	jutil.CheckArgument(NodeUtil_IsNameDeclaration(declChild.GetParent()))
	jutil.CheckArgument(newStatement.GetParent().IsNil())

	decl := declChild.GetParent()
	if decl.HasOneChild() {
		decl.ReplaceWith(newStatement)
	} else if declChild.GetNext().IsNil() {
		declChild.Detach()
		newStatement.InsertAfter(decl)
	} else if declChild.GetPrevious().IsNil() {
		declChild.Detach()
		newStatement.InsertBefore(decl)
	} else {
		jutil.CheckState(decl.HasMoreThanOneChild())
		newDecl := rhino.NewNode(decl.GetToken()).Srcref(decl)
		var next rhino.Node
		for after := declChild.GetNext(); !after.IsNil(); after = next {
			next = after.GetNext()
			newDecl.AddChildToBack(after.Detach())
		}
		declChild.Detach()
		newStatement.InsertAfter(decl)
		newDecl.InsertAfter(newStatement)
	}
}

// NodeUtil_MaybeAddFinally adds a finally block if one does not exist.
//
// Java: NodeUtil.maybeAddFinally(Node)
func NodeUtil_MaybeAddFinally(tryNode rhino.Node) {
	jutil.CheckState(tryNode.IsTry())
	if !NodeUtil_HasFinally(tryNode) {
		tryNode.AddChildToBack(rhino.IR_Block().Srcref(tryNode))
	}
}

// NodeUtil_TryMergeBlock merges a block with its parent block.
//
// ignoreBlockScopedDeclarations: merge the block regardless of any inner
// block-scoped declarations that may cause name collisions. use if e.g. the
// AST is normalized. Returns whether the block was removed.
//
// Java: NodeUtil.tryMergeBlock(Node, boolean)
func NodeUtil_TryMergeBlock(block rhino.Node, ignoreBlockScopedDeclarations bool) bool {
	jutil.CheckState(block.IsBlock())
	parent := block.GetParent()
	canMerge := ignoreBlockScopedDeclarations || NodeUtil_CanMergeBlock(block)
	// Try to remove the block if its parent is a block/script or if its
	// parent is label and it has exactly one child.
	if NodeUtil_IsStatementBlock(parent) && canMerge {
		previous := block
		for block.HasChildren() {
			child := block.RemoveFirstChild()
			child.InsertAfter(previous)
			previous = child
		}
		block.Detach()
		return true
	} else {
		return false
	}
}

// NodeUtil_CanMergeBlock: a check inside a block to see if there are const,
// let, class, or function declarations to be safe and not hoist them into
// the upper block.
//
// Returns whether the block can be removed.
//
// Java: NodeUtil.canMergeBlock(Node)
func NodeUtil_CanMergeBlock(block rhino.Node) bool {
	for c := block.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		switch c.GetToken() {
		case rhino.Token_LABEL:
			if NodeUtil_CanMergeBlock(c) {
				continue
			} else {
				return false
			}
		case rhino.Token_CONST, rhino.Token_LET, rhino.Token_CLASS, rhino.Token_FUNCTION:
			return false
		default:
			continue
		}
	}
	return true
}

// NodeUtil_IsCallOrNew returns whether the call is a NEW or CALL node.
//
// Java: NodeUtil.isCallOrNew(Node)
func NodeUtil_IsCallOrNew(node rhino.Node) bool {
	return node.IsCall() || node.IsNew() || node.IsOptChainCall()
}

// NodeUtil_GetFunctionBody returns a BLOCK node for the given FUNCTION node.
//
// Java: NodeUtil.getFunctionBody(Node)
func NodeUtil_GetFunctionBody(fn rhino.Node) rhino.Node {
	jutil.CheckArgumentMsg1(fn.IsFunction(), "%s", fn)
	return fn.GetLastChild()
}

// NodeUtil_GetCallTargetResolvingIndirectCalls returns the call target for a
// call expression, resolving an indirected call (`(0, foo)()`) if present.
//
// Java: NodeUtil.getCallTargetResolvingIndirectCalls(Node)
func NodeUtil_GetCallTargetResolvingIndirectCalls(call rhino.Node) rhino.Node {
	jutil.CheckArgumentMsg1(
		call.IsCall() || call.IsNew() || call.IsTaggedTemplateLit(),
		"must be call, new or tagged-template-literal expression, got %s",
		call)
	target := call.GetFirstChild()
	if target.IsComma() && target.GetSecondChild().IsQualifiedName() {
		return target.GetSecondChild()
	}
	return target
}

// NodeUtil_IsDeclaration: is the node a var, const, let, function, or class
// declaration? See NodeUtil_IsFunctionDeclaration, NodeUtil_IsClassDeclaration,
// and NodeUtil_IsNameDeclaration.
//
// Java: NodeUtil.isDeclaration(Node)
func NodeUtil_IsDeclaration(n rhino.Node) bool {
	return NodeUtil_IsNameDeclaration(n) || NodeUtil_IsFunctionDeclaration(n) || NodeUtil_IsClassDeclaration(n)
}

// NodeUtil_IsNamedExportsLiteral: whether this is an assignment to 'exports'
// that creates named exports.
//
//   - exports = {a, b}; // named export, returns true.
//   - exports = 0; // namespace export, returns false.
//   - exports = {a: 0, b}; // namespace export, returns false.
//
// Java: NodeUtil.isNamedExportsLiteral(Node)
func NodeUtil_IsNamedExportsLiteral(objectLiteral rhino.Node) bool {
	if !objectLiteral.IsObjectLit() || !objectLiteral.HasChildren() {
		return false
	}
	for key := objectLiteral.GetFirstChild(); !key.IsNil(); key = key.GetNext() {
		if !key.IsStringKey() || key.IsQuotedStringKey() {
			return false
		}
		if !key.GetFirstChild().IsName() {
			return false
		}
	}
	return true
}

// NodeUtil_IsFunctionDeclaration: is this node a function declaration? A
// function declaration is a function that has a name that is added to the
// current scope (i.e. a function that is not part of a expression; see
// NodeUtil_IsFunctionExpression).
//
// Java: NodeUtil.isFunctionDeclaration(Node)
func NodeUtil_IsFunctionDeclaration(n rhino.Node) bool {
	// Note: There is currently one case where an unnamed function has a declaration parent.
	// `export default function() {...}`
	// In this case we consider the function to be an expression.
	return n.IsFunction() && NodeUtil_IsDeclarationParent(n.GetParent()) && NodeUtil_IsNamedFunction(n)
}

// NodeUtil_IsMethodDeclaration: is this node a class or object literal
// member function?
//
// examples:
//
//	class C {
//	  f() {}
//	  get x() { return this.x_; }
//	  set x(v) { this.x_ = v; }
//	  [someExpr]() {}
//	}
//	obj = {
//	  f() {}
//	  get x() { return this.x_; }
//	  set x(v) { this.x_ = v; }
//	  [someExpr]() {}
//	}
//
// Java: NodeUtil.isMethodDeclaration(Node)
func NodeUtil_IsMethodDeclaration(n rhino.Node) bool {
	if n.IsFunction() {
		parent := n.GetParent()
		switch parent.GetToken() {
		case rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_MEMBER_FUNCTION_DEF:
			// `({ get x() {} })`
			// `({ set x(v) {} })`
			// `({ f() {} })`
			return true
		case rhino.Token_COMPUTED_PROP:
			// `({ [expression]() {} })`
			// `({ get [expression]() {} })`
			// `({ set [expression](x) {} })`
			// (but not `({ [expression]: function() {} })`
			// The first child is the expression, and could possibly be a function.
			return parent.GetLastChild() == n &&
				(parent.GetBooleanProp(rhino.Node_COMPUTED_PROP_METHOD) ||
					parent.GetBooleanProp(rhino.Node_COMPUTED_PROP_GETTER) ||
					parent.GetBooleanProp(rhino.Node_COMPUTED_PROP_SETTER))
		default:
			return false
		}
	} else {
		return false
	}
}

// NodeUtil_IsClassDeclaration: is this a class declaration.
//
// Java: NodeUtil.isClassDeclaration(Node)
func NodeUtil_IsClassDeclaration(n rhino.Node) bool {
	return n.IsClass() && NodeUtil_IsDeclarationParent(n.GetParent()) && NodeUtil_IsNamedClass(n)
}

// NodeUtil_IsHoistedFunctionDeclaration: is this node a hoisted function
// declaration? A function declaration in the scope root is hoisted to the
// top of the scope. See NodeUtil_IsFunctionDeclaration.
//
// Java: NodeUtil.isHoistedFunctionDeclaration(Node)
func NodeUtil_IsHoistedFunctionDeclaration(n rhino.Node) bool {
	if NodeUtil_IsFunctionDeclaration(n) {
		parent := n.GetParent()
		return parent.IsScript() ||
			parent.IsModuleBody() ||
			parent.GetParent().IsFunction() ||
			parent.IsExport()
	}
	return false
}

// Java: NodeUtil.isBlockScopedFunctionDeclaration(Node)
func NodeUtil_IsBlockScopedFunctionDeclaration(n rhino.Node) bool {
	if !NodeUtil_IsFunctionDeclaration(n) {
		return false
	}
	current := n.GetParent()
	for !current.IsNil() {
		switch current.GetToken() {
		case rhino.Token_BLOCK:
			return !current.GetParent().IsFunction()
		case rhino.Token_FUNCTION, rhino.Token_SCRIPT, rhino.Token_DECLARE, rhino.Token_EXPORT,
			rhino.Token_MODULE_BODY:
			return false
		default:
			jutil.CheckStateMsg1(current.IsLabel(), "%s", current)
			current = current.GetParent()
		}
	}
	return false
}

// Java: NodeUtil.isFunctionBlock(Node)
func NodeUtil_IsFunctionBlock(n rhino.Node) bool {
	return n.IsBlock() && n.HasParent() && n.GetParent().IsFunction()
}

// Java: NodeUtil.isClassStaticBlock(Node)
func NodeUtil_IsClassStaticBlock(n rhino.Node) bool {
	return n.IsBlock() && n.HasParent() && n.GetParent().IsClassMembers()
}

// NodeUtil_IsNonStaticClassFieldInitializer returns whether n is the
// initializer expression of a non-static class field (e.g. a in
// `class { x = a; }` or `class { [k] = a; }`).
//
// Unlike static field initializers and computed field keys, which are
// evaluated eagerly when the class definition is evaluated, non-static field
// initializers are evaluated once per new call (like an inner function
// body).
//
// Java: NodeUtil.isNonStaticClassFieldInitializer(Node)
func NodeUtil_IsNonStaticClassFieldInitializer(n rhino.Node) bool {
	parent := n.GetParent()
	// Note: COMPUTED_FIELD_DEF has only one child (the key) if there is no initializer.
	if parent.IsMemberFieldDef() ||
		(parent.IsComputedFieldDef() && n == parent.GetSecondChild()) {
		return !parent.IsStaticMember()
	}
	return false
}

// NodeUtil_IsFunctionOrNonStaticClassFieldInitializer returns whether n is a
// function or a non-static class field initializer: both are code whose
// evaluation is deferred until later (and possibly repeated), and which may
// capture variables from the enclosing scope.
//
// Java: NodeUtil.isFunctionOrNonStaticClassFieldInitializer(Node)
func NodeUtil_IsFunctionOrNonStaticClassFieldInitializer(n rhino.Node) bool {
	return n.IsFunction() || NodeUtil_IsNonStaticClassFieldInitializer(n)
}

// NodeUtil_IsFunctionExpression: is a FUNCTION node a function expression?
//
// A function expression is a function that:
//
//   - has either no name or a name that is not added to the current scope
//   - AND can be manipulated as an expression (assigned to variables, passed
//     to functions, etc.) i.e. It is not a method declaration on a class or
//     object literal.
//
// Some examples of function expressions:
//
//	(function () {})
//	(function f() {})()
//	[ function f() {} ]
//	var f = function f() {};
//	for (function f() {};;) {}
//	export default function() {}
//	() => 1
//
// Some examples of functions that are not expressions:
//
//	function f() {}
//	if (x); else function f() {}
//	for (;;) { function f() {} }
//	export default function f() {}
//	({
//	  f() {},
//	  set x(v) {},
//	  get x() {},
//	  [expr]() {}
//	})
//	class {
//	  f() {}
//	  set x(v) {}
//	  get x() {}
//	  [expr]() {}
//	}
//
// Returns whether n is a function used within an expression.
//
// Java: NodeUtil.isFunctionExpression(Node)
func NodeUtil_IsFunctionExpression(n rhino.Node) bool {
	return n.IsFunction() && !NodeUtil_IsFunctionDeclaration(n) && !NodeUtil_IsMethodDeclaration(n)
}

// NodeUtil_IsNamedFunctionExpression returns whether the node is both a
// function expression and the function is named.
//
// Java: NodeUtil.isNamedFunctionExpression(Node)
func NodeUtil_IsNamedFunctionExpression(n rhino.Node) bool {
	return NodeUtil_IsFunctionExpression(n) && n.GetFirstChild().GetString() != ""
}

// NodeUtil_IsClassExpression: see NodeUtil_IsFunctionExpression. Returns
// whether n is a class used within an expression.
//
// Java: NodeUtil.isClassExpression(Node)
func NodeUtil_IsClassExpression(n rhino.Node) bool {
	return n.IsClass() && (!NodeUtil_IsNamedClass(n) || !NodeUtil_IsDeclarationParent(n.GetParent()))
}

// NodeUtil_IsNamedClassExpression returns whether the node is both a class
// expression and the class is named.
//
// Java: NodeUtil.isNamedClassExpression(Node)
func NodeUtil_IsNamedClassExpression(n rhino.Node) bool {
	return NodeUtil_IsClassExpression(n) && n.GetFirstChild().IsName()
}

// NodeUtil_IsNamedFunction returns whether n is a function with a nonempty
// name. Differs from NodeUtil_IsFunctionDeclaration because the name might
// in a function expression and not be added to the current scope.
//
// Some named functions include
//
//	(function f() {})();
//	export default function f() {};
//	function f() {};
//	var f = function f() {};
//
// Java: NodeUtil.isNamedFunction(Node)
func NodeUtil_IsNamedFunction(n rhino.Node) bool {
	return n.IsFunction() && NodeUtil_IsReferenceName(n.GetFirstChild())
}

// NodeUtil_IsNamedClass: see NodeUtil_IsNamedFunction. Returns whether n is
// a named class.
//
// Java: NodeUtil.isNamedClass(Node)
func NodeUtil_IsNamedClass(n rhino.Node) bool {
	return n.IsClass() && NodeUtil_IsReferenceName(n.GetFirstChild())
}

// NodeUtil_IsBleedingFunctionName returns whether this is a bleeding
// function (an anonymous named function that bleeds into the inner scope).
//
// Java: NodeUtil.isBleedingFunctionName(Node)
func NodeUtil_IsBleedingFunctionName(n rhino.Node) bool {
	if !n.IsName() || n.GetString() == "" {
		return false
	}
	parent := n.GetParent()
	return NodeUtil_IsFunctionExpression(parent) && n == parent.GetFirstChild()
}

// NodeUtil_IsEmptyFunctionExpression determines if a node is a function
// expression that has an empty body.
//
// Returns whether the given node is a function expression that is empty.
//
// Java: NodeUtil.isEmptyFunctionExpression(Node)
func NodeUtil_IsEmptyFunctionExpression(node rhino.Node) bool {
	return NodeUtil_IsFunctionExpression(node) && NodeUtil_IsEmptyBlock(node.GetLastChild())
}

// NodeUtil_DoesFunctionReferenceOwnArgumentsObject returns whether a
// function has a reference to its own "arguments" object.
//
// Java: NodeUtil.doesFunctionReferenceOwnArgumentsObject(Node)
func NodeUtil_DoesFunctionReferenceOwnArgumentsObject(fn rhino.Node) bool {
	jutil.CheckArgument(fn.IsFunction())
	if fn.IsArrowFunction() {
		return false
	}
	return NodeUtil_ReferencesArgumentsHelper(fn.GetSecondChild()) ||
		NodeUtil_ReferencesArgumentsHelper(fn.GetLastChild())
}

// NodeUtil_ReferencesArgumentsHelper returns whether any child is a
// reference to the "arguments" object of the root. Effectively, this
// includes arrow method bodies (which don't have their own) and excludes
// other functions which shadow the "arguments" value with their own.
//
// Java: NodeUtil.referencesArgumentsHelper(Node)
func NodeUtil_ReferencesArgumentsHelper(node rhino.Node) bool {
	if node.IsName() && node.GetString() == "arguments" {
		return true
	}

	if NodeUtil_IsNonArrowFunction(node) {
		return false
	}

	for c := node.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if NodeUtil_ReferencesArgumentsHelper(c) {
			return true
		}
	}

	return false
}

// NodeUtil_IsObjectCallMethod returns whether node is a call to methodName:
//
//	a.f(...)
//	a?.f(...)
//	a['f'](...)
//
// Java: NodeUtil.isObjectCallMethod(Node, String)
func NodeUtil_IsObjectCallMethod(callNode rhino.Node, methodName string) bool {
	if callNode.IsCall() || callNode.IsOptChainCall() {
		callee := callNode.GetFirstChild()
		if NodeUtil_IsNormalOrOptChainGetProp(callee) {
			return callee.GetString() == methodName
		} else if NodeUtil_IsNormalOrOptChainGet(callee) {
			last := callee.GetLastChild()
			if !last.IsNil() && last.IsStringLit() {
				propName := last.GetString()
				return propName == methodName
			}
		}
	}
	return false
}

// NodeUtil_IsFunctionObjectCall returns whether the callNode represents an
// expression in the form of: x.call(...) x['call'](...)
//
// Java: NodeUtil.isFunctionObjectCall(Node)
func NodeUtil_IsFunctionObjectCall(callNode rhino.Node) bool {
	return NodeUtil_IsObjectCallMethod(callNode, "call")
}

// NodeUtil_IsFunctionObjectApply returns whether the callNode represents an
// expression in the form of: x.apply(...) x['apply'](...)
//
// Java: NodeUtil.isFunctionObjectApply(Node)
func NodeUtil_IsFunctionObjectApply(callNode rhino.Node) bool {
	return NodeUtil_IsObjectCallMethod(callNode, "apply")
}

// NodeUtil_IsNameDeclOrSimpleAssignLhs determines whether this node is
// strictly on the left hand side of an assign or var initialization.
// Notably, this does not include all L-values, only statements where the
// node is used only as an L-value.
//
// n is the node, parent the parent of the node. Returns true if n is the
// left hand of an assign.
//
// Java: NodeUtil.isNameDeclOrSimpleAssignLhs(Node, Node)
func NodeUtil_IsNameDeclOrSimpleAssignLhs(n rhino.Node, parent rhino.Node) bool {
	return (parent.IsAssign() && parent.GetFirstChild() == n) || NodeUtil_IsNameDeclaration(parent)
}

// NodeUtil_IsLValue determines whether this node is used as an L-value.
// Notice that sometimes names are used as both L-values and R-values.
//
// We treat "var x;" and "let x;" as an L-value because it's syntactically
// similar to "var x = undefined", even though it's technically not an
// L-value. But it kind of makes sense if you treat it as "assignment to
// 'undefined' at the top of the scope".
//
// Returns true if n is an L-value.
//
// Java: NodeUtil.isLValue(Node)
func NodeUtil_IsLValue(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_NAME, rhino.Token_GETPROP, rhino.Token_GETELEM:
	default:
		return false
	}

	parent := n.GetParent()
	if parent.IsNil() {
		return false
	}

	switch parent.GetToken() {
	case rhino.Token_IMPORT_SPEC:
		return parent.GetLastChild() == n
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST, rhino.Token_ITER_REST, rhino.Token_OBJECT_REST,
		rhino.Token_PARAM_LIST, rhino.Token_IMPORT, rhino.Token_INC, rhino.Token_DEC, rhino.Token_CATCH:
		return true
	case rhino.Token_CLASS, rhino.Token_FUNCTION, rhino.Token_DEFAULT_VALUE, rhino.Token_FOR, rhino.Token_FOR_IN,
		rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
		return parent.GetFirstChild() == n
	case rhino.Token_ARRAY_PATTERN, rhino.Token_STRING_KEY, rhino.Token_COMPUTED_PROP:
		return NodeUtil_IsLhsByDestructuring(n)
	default:
		return NodeUtil_IsAssignmentOp(parent) && parent.GetFirstChild() == n
	}
}

// NodeUtil_IsDeclarationLValue determines whether this node is used as an
// L-value that is a declaration.
//
// `x = 5;` is an L-value but does not declare a variable.
//
// Java: NodeUtil.isDeclarationLValue(Node)
func NodeUtil_IsDeclarationLValue(n rhino.Node) bool {
	isLValue := NodeUtil_IsLValue(n)

	if !isLValue {
		return false
	}

	parent := n.GetParent()

	switch parent.GetToken() {
	case rhino.Token_IMPORT_SPEC, rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST, rhino.Token_PARAM_LIST,
		rhino.Token_IMPORT, rhino.Token_CATCH, rhino.Token_CLASS, rhino.Token_FUNCTION:
		return true
	case rhino.Token_STRING_KEY,
		rhino.Token_OBJECT_PATTERN,
		rhino.Token_ARRAY_PATTERN,
		rhino.Token_COMPUTED_PROP,
		rhino.Token_DEFAULT_VALUE,
		rhino.Token_ITER_REST,
		rhino.Token_OBJECT_REST:
		rootTargetParent := NodeUtil_GetRootTarget(n).GetParent()
		if rootTargetParent.IsRest() || rootTargetParent.IsDefaultValue() {
			return rootTargetParent.GetParent().IsParamList()
		} else if rootTargetParent.IsDestructuringLhs() {
			return NodeUtil_IsNameDeclaration(rootTargetParent.GetParent())
		}
		return rootTargetParent.IsParamList() || rootTargetParent.IsCatch()
	default:
		return false
	}
}

// Java: NodeUtil.isLhsOfAssign(Node)
func NodeUtil_IsLhsOfAssign(n rhino.Node) bool {
	parent := n.GetParent()
	return !parent.IsNil() && parent.IsAssign() && parent.GetFirstChild() == n
}

// Java: NodeUtil.isImportedName(Node)
func NodeUtil_IsImportedName(n rhino.Node) bool {
	parent := n.GetParent()
	return parent.IsImport() || (parent.IsImportSpec() && parent.GetLastChild() == n)
}

// NodeUtil_GetDeclaringParent returns the node that is effectively declaring
// the given target.
//
// Examples:
//
//	const a = 1; // getDeclaringParent(a) returns CONST
//	let {[expression]: [x = 3]} = obj; // getDeclaringParent(x) returns LET
//	function foo({a, b}) {}; // getDeclaringParent(a) returns PARAM_LIST
//	function foo(a = 1) {}; // getDeclaringParent(a) returns PARAM_LIST
//	function foo({a, b} = obj) {}; // getDeclaringParent(a) returns PARAM_LIST
//	function foo(...a) {}; // getDeclaringParent(a) returns PARAM_LIST
//	function foo() {}; // gotRootTarget(foo) returns FUNCTION
//	class foo {}; // gotRootTarget(foo) returns CLASS
//	import foo from './foo'; // getDeclaringParent(foo) returns IMPORT
//	import {foo} from './foo'; // getDeclaringParent(foo) returns IMPORT
//	import {foo as bar} from './foo'; // getDeclaringParent(bar) returns IMPORT
//	} catch (err) { // getDeclaringParent(err) returns CATCH
//
// targetNode is a NAME, OBJECT_PATTERN, or ARRAY_PATTERN; returns a node of
// type LET, CONST, VAR, FUNCTION, CLASS, PARAM_LIST, CATCH, or IMPORT. Throws
// IllegalStateException if targetNode is not actually used as a declaration
// target.
//
// Java: NodeUtil.getDeclaringParent(Node)
func NodeUtil_GetDeclaringParent(targetNode rhino.Node) rhino.Node {
	rootTarget := NodeUtil_GetRootTarget(targetNode)
	parent := rootTarget.GetParent()
	if parent.IsRest() || parent.IsDefaultValue() {
		// e.g. `function foo(targetNode1 = default, ...targetNode2) {}`
		parent = parent.GetParent()
		jutil.CheckStateMsg1(parent.IsParamList(), "%s", parent)
	} else if parent.IsDestructuringLhs() {
		// e.g. `let [a, b] = something;` targetNode is `[a, b]`
		parent = parent.GetParent()
		jutil.CheckStateMsg1(NodeUtil_IsNameDeclaration(parent), "%s", parent)
	} else if parent.IsClass() || parent.IsFunction() {
		// e.g. `function targetNode() {}`
		// e.g. `class targetNode {}`
		jutil.CheckStateMsg1(targetNode == parent.GetFirstChild(), "%s", targetNode)
	} else if parent.IsImportSpec() {
		// e.g. `import {foo as targetNode} from './foo';
		jutil.CheckStateMsg1(targetNode == parent.GetSecondChild(), "%s", targetNode)
		// import -> import_specs -> import_spec
		// we want import
		parent = parent.GetGrandparent()
		jutil.CheckStateMsg1(parent.IsImport(), "%s", parent)
	} else {
		// e.g. `function foo(targetNode) {};`
		// e.g. `let targetNode = something;`
		// e.g. `import targetNode from './foo';
		// e.g. `} catch (foo) {`
		jutil.CheckStateMsg1(
			parent.IsParamList() ||
				NodeUtil_IsNameDeclaration(parent) ||
				parent.IsImport() ||
				parent.IsCatch(),
			"%s",
			parent)
	}
	return parent
}

// NodeUtil_GetRootTarget returns the outermost target enclosing the given
// assignment target.
//
// Returns targetNode itself if there is no enclosing target.
//
// Examples:
//
//	const a = 1; // getRootTarget(a) returns a
//	let {[expression]: [x = 3]} = obj; // getRootTarget(x) returns {[expression]: [x = 3]}
//	{a = 1} = obj; // getRootTarget(a) returns {a = 1}
//	{[expression]: [x = 3]} = obj; // getRootTarget(x) returns {[expression]: [x = 3]}
//	function foo({a, b}) {}; // getRootTarget(a) returns {a, b}
//	function foo(a = 1) {}; // getRootTarget(a) returns a
//	function foo({a, b} = obj) {}; // getRootTarget(a) returns a
//	function foo(...a) {}; // getRootTarget(a) returns a
//	function foo() {}; // gotRootTarget(foo) returns foo
//	class foo {}; // gotRootTarget(foo) returns foo
//	import foo from './foo'; // getRootTarget(foo) returns foo
//	import {foo} from './foo'; // getRootTarget(foo) returns foo
//	import {foo as bar} from './foo'; // getRootTarget(bar) returns bar
//
// Throws IllegalStateException if targetNode is not actually used as a
// target.
//
// Java: NodeUtil.getRootTarget(Node)
func NodeUtil_GetRootTarget(targetNode rhino.Node) rhino.Node {
	enclosingTarget := targetNode
	for nextTarget := NodeUtil_GetEnclosingTarget(enclosingTarget); !nextTarget.IsNil(); nextTarget = NodeUtil_GetEnclosingTarget(enclosingTarget) {
		enclosingTarget = nextTarget
	}
	return enclosingTarget
}

// NodeUtil_GetEnclosingTarget returns the immediately enclosing target node
// for a given target node, or null if none found.
//
// See NodeUtil_GetRootTarget for examples.
//
// Java: NodeUtil.getEnclosingTarget(Node)
func NodeUtil_GetEnclosingTarget(targetNode rhino.Node) rhino.Node {
	jutil.CheckStateMsg1(rhino.CheckNotNull(targetNode).IsValidAssignmentTarget(), "%s", targetNode)
	parent := rhino.CheckNotNullMsg1(targetNode.GetParent(), "%s", targetNode)
	targetIsFirstChild := parent.GetFirstChild() == targetNode
	if parent.IsDefaultValue() || parent.IsRest() {
		// in `([something = targetNode] = x)` targetNode isn't actually acting
		// as a target.
		jutil.CheckStateMsg1(targetIsFirstChild, "%s", parent)
		// The DEFAULT_VALUE or REST occupies the place where the assignment target it contains would
		// otherwise be in the AST, so pretend it is the target for the logic below.
		targetNode = parent
		parent = rhino.CheckNotNull(targetNode.GetParent())
		targetIsFirstChild = targetNode == parent.GetFirstChild()
	}
	switch parent.GetToken() {
	case rhino.Token_ARRAY_PATTERN:
		// e.g. ([targetNode] = something)
		return parent

	case rhino.Token_OBJECT_PATTERN:
		// e.g. ({...rest} = something);
		return parent

	case rhino.Token_COMPUTED_PROP:
		// e.g. ({[expression]: targetNode} = something)
		// e.g. ({[expression]: targetNode = default} = something)
		// make sure the effective target (targetNode or DEFAULT_VALUE containing it)
		// isn't the expression part
		jutil.CheckStateMsg1(!targetIsFirstChild, "%s", parent)
		grandparent := rhino.CheckNotNullMsg1(parent.GetParent(), "%s", parent)
		jutil.CheckStateMsg1(grandparent.IsObjectPattern(), "%s", grandparent)
		return grandparent
	case rhino.Token_STRING_KEY:
		// e.g. ({parent: targetNode} = something)
		grandparent := rhino.CheckNotNullMsg1(parent.GetParent(), "%s", parent)
		jutil.CheckStateMsg1(grandparent.IsObjectPattern(), "%s", grandparent)
		return grandparent

	case rhino.Token_PARAM_LIST, rhino.Token_LET, rhino.Token_CONST, rhino.Token_VAR:
		// e.g. `function foo(targetNode) {}`
		// or non-destructured declarations
		// e.g. `let targetNode = 3;`
		return rhino.NilNode

	case rhino.Token_FUNCTION, rhino.Token_CLASS:
		// e.g. `function targetNode() {}`
		// e.g. `class targetNode {}`
		jutil.CheckStateMsg1(targetIsFirstChild, "%s", targetNode)
		return rhino.NilNode

	case rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
		// e.g. `for ({length} in obj) {}` // targetNode is `{length}`
		// e.g. `for ({length} of obj) {}` // targetNode is `{length}`
		jutil.CheckStateMsg1(targetIsFirstChild, "%s", targetNode)
		return rhino.NilNode

	case rhino.Token_DESTRUCTURING_LHS:
		// destructured declarations
		// e.g. `let [a] = 3`; // targetNode is `[a]`
		jutil.CheckStateMsg1(targetIsFirstChild, "%s", targetNode)
		return rhino.NilNode

	case rhino.Token_IMPORT:
		// e.g. `import targetNode from './foo/bar';`
		return rhino.NilNode

	case rhino.Token_IMPORT_SPEC:
		// e.g. `import {bar as targetNode} from './foo/bar';`
		// e.g. `import {targetNode} from './foo/bar';` // AST will have {targetNode as targetNode}
		jutil.CheckStateMsg1(!targetIsFirstChild, "%s", parent)
		return rhino.NilNode

	case rhino.Token_CATCH:
		// e.g. `try {} catch (foo) {}`
		return rhino.NilNode

	default:
		// e.g. targetNode = something
		jutil.CheckStateMsg1(NodeUtil_IsAssignmentOp(parent) && targetIsFirstChild, "%s", parent)
		return rhino.NilNode
	}
}

// NodeUtil_IsLhsByDestructuring returns true if the node is a lhs value of a
// destructuring assignment.
//
// For example, x in `var [x] = [1];`, `var [...x] = [1];`, and
// `var {a: x} = {a: 1}` or a.b in `([a.b] = [1]);` or
// `({key: a.b} = {key: 1});`
//
// Java: NodeUtil.isLhsByDestructuring(Node)
func NodeUtil_IsLhsByDestructuring(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_NAME, rhino.Token_GETPROP, rhino.Token_GETELEM:
		return NodeUtil_IsLhsByDestructuringHelper(n)
	default:
		return false
	}
}

// NodeUtil_IsLhsByDestructuringHelper returns true if the given node is
// either an LHS node in a destructuring pattern or if one of its descendants
// contains an LHS node in a destructuring pattern. For example, in
// `var {a: b = 3}}`, this returns true given the NAME b or the DEFAULT_VALUE
// node containing b.
//
// Java: NodeUtil.isLhsByDestructuringHelper(Node)
func NodeUtil_IsLhsByDestructuringHelper(n rhino.Node) bool {
	parent := n.GetParent()
	grandparent := n.GetGrandparent()

	switch parent.GetToken() {
	// `b` in `var [b] = ...`
	// `b` in `var [...b] = ...`
	case rhino.Token_ARRAY_PATTERN, rhino.Token_ITER_REST, rhino.Token_OBJECT_REST:
		return true
	case rhino.Token_COMPUTED_PROP:
		return !n.IsFirstChildOf(parent) && grandparent.IsObjectPattern()
	// the "b" in "var {a: b} = ..."
	case rhino.Token_STRING_KEY:
		return grandparent.IsObjectPattern()
	case rhino.Token_DEFAULT_VALUE:
		return n.IsFirstChildOf(parent) &&
			// The first child of a DEFAULT_VALUE is a NAME node and a potential LHS.
			// The second child is the value, so never a LHS node.
			NodeUtil_IsLhsByDestructuringHelper(parent)

	default:
		return false
	}
}

// NodeUtil_MayBeObjectLitKey determines whether a node represents a possible
// object literal key (e.g. key1 in {key1: value1, key2: value2}). Computed
// properties are excluded here (see b/111621528). This method does not check
// whether the node is actually in an object literal! it also returns true
// for object pattern keys, and member functions/getters in ES6 classes.
//
// Java: NodeUtil.mayBeObjectLitKey(Node)
//
// TODO(b/189993301): should fields be added to this method?
func NodeUtil_MayBeObjectLitKey(node rhino.Node) bool {
	switch node.GetToken() {
	case rhino.Token_STRING_KEY, rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_MEMBER_FUNCTION_DEF:
		return true
	default:
		return false
	}
}

// NodeUtil_IsObjectLitKey determines whether a node represents an object
// literal key (e.g. key1 in {key1: value1, key2: value2}) and is in an
// object literal. Computed properties are excluded here (see b/111621528).
//
// Java: NodeUtil.isObjectLitKey(Node)
func NodeUtil_IsObjectLitKey(node rhino.Node) bool {
	return node.GetParent().IsObjectLit() && NodeUtil_MayBeObjectLitKey(node)
}

// NodeUtil_GetObjectOrClassLitKeyName gets the name of an object literal
// key.
//
// Java: NodeUtil.getObjectOrClassLitKeyName(Node)
func NodeUtil_GetObjectOrClassLitKeyName(key rhino.Node) string {
	keyNode := NodeUtil_GetObjectOrClassLitKeyNode(key)
	if !keyNode.IsNil() {
		return keyNode.GetString()
	}
	panic(jutil.NewIllegalStateExceptionMsg("Unexpected node type: " + key.String()))
}

// NodeUtil_GetObjectOrClassLitKeyNode gets the Node that defines the name of
// an object literal key.
//
// Java: NodeUtil.getObjectOrClassLitKeyNode(Node)
func NodeUtil_GetObjectOrClassLitKeyNode(key rhino.Node) rhino.Node {
	switch key.GetToken() {
	case rhino.Token_STRING_KEY, rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_MEMBER_FUNCTION_DEF,
		rhino.Token_MEMBER_FIELD_DEF:
		return key
	case rhino.Token_COMPUTED_PROP, rhino.Token_COMPUTED_FIELD_DEF:
		if key.GetFirstChild().
			IsStringLit() { // TODO(b/189993301): may be an issue with non string lits
			return key.GetFirstChild()
		}
		return rhino.NilNode
	default:
	}
	panic(jutil.NewIllegalStateExceptionMsg("Unexpected node type: " + key.String()))
}

// NodeUtil_IsGetOrSetKey determines whether a node represents an object
// literal get or set key (e.g. key1 in {get key1() {}, set key2(a){}).
//
// Java: NodeUtil.isGetOrSetKey(Node)
func NodeUtil_IsGetOrSetKey(node rhino.Node) bool {
	switch node.GetToken() {
	case rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF:
		return true
	case rhino.Token_COMPUTED_PROP:
		return node.GetBooleanProp(rhino.Node_COMPUTED_PROP_GETTER) ||
			node.GetBooleanProp(rhino.Node_COMPUTED_PROP_SETTER)
	default:
		return false
	}
}

// NodeUtil_OpToStr converts an operator's token value (see Token) to a
// string representation.
//
// operator is the operator's token value to convert; returns the string
// representation or null if the token value is not an operator.
//
// Java: NodeUtil.opToStr(Token)
func NodeUtil_OpToStr(operator rhino.Token) *string {
	// Go-only: language-adapter — Java returns interned string constants; the
	// results of the switch below are computed once per token so that the
	// *string result does not allocate on every call (the code printer calls
	// this per operator).
	return nodeUtil_opToStrTable[operator]
}

var nodeUtil_opToStrTable = func() []*string {
	values := rhino.Token_values()
	table := make([]*string, len(values))
	for _, t := range values {
		if s, ok := nodeUtil_opToStrSwitch(t); ok {
			table[t] = jutil.StringPtr(s)
		}
	}
	return table
}()

// nodeUtil_opToStrSwitch is the switch of NodeUtil.opToStr(Token); ok is
// false where Java returns null.
func nodeUtil_opToStrSwitch(operator rhino.Token) (string, bool) {
	switch operator {
	case rhino.Token_COALESCE:
		return "??", true
	case rhino.Token_BITOR:
		return "|", true
	case rhino.Token_OR:
		return "||", true
	case rhino.Token_BITXOR:
		return "^", true
	case rhino.Token_AND:
		return "&&", true
	case rhino.Token_BITAND:
		return "&", true
	case rhino.Token_SHEQ:
		return "===", true
	case rhino.Token_EQ:
		return "==", true
	case rhino.Token_NOT:
		return "!", true
	case rhino.Token_NE:
		return "!=", true
	case rhino.Token_SHNE:
		return "!==", true
	case rhino.Token_LSH:
		return "<<", true
	case rhino.Token_IN:
		return "in", true
	case rhino.Token_LE:
		return "<=", true
	case rhino.Token_LT:
		return "<", true
	case rhino.Token_URSH:
		return ">>>", true
	case rhino.Token_RSH:
		return ">>", true
	case rhino.Token_GE:
		return ">=", true
	case rhino.Token_GT:
		return ">", true
	case rhino.Token_MUL:
		return "*", true
	case rhino.Token_DIV:
		return "/", true
	case rhino.Token_MOD:
		return "%", true
	case rhino.Token_EXPONENT:
		return "**", true
	case rhino.Token_BITNOT:
		return "~", true
	case rhino.Token_ADD, rhino.Token_POS:
		return "+", true
	case rhino.Token_SUB, rhino.Token_NEG:
		return "-", true
	case rhino.Token_ASSIGN:
		return "=", true
	case rhino.Token_ASSIGN_BITOR:
		return "|=", true
	case rhino.Token_ASSIGN_BITXOR:
		return "^=", true
	case rhino.Token_ASSIGN_BITAND:
		return "&=", true
	case rhino.Token_ASSIGN_LSH:
		return "<<=", true
	case rhino.Token_ASSIGN_RSH:
		return ">>=", true
	case rhino.Token_ASSIGN_URSH:
		return ">>>=", true
	case rhino.Token_ASSIGN_ADD:
		return "+=", true
	case rhino.Token_ASSIGN_SUB:
		return "-=", true
	case rhino.Token_ASSIGN_MUL:
		return "*=", true
	case rhino.Token_ASSIGN_EXPONENT:
		return "**=", true
	case rhino.Token_ASSIGN_DIV:
		return "/=", true
	case rhino.Token_ASSIGN_MOD:
		return "%=", true
	case rhino.Token_ASSIGN_OR:
		return "||=", true
	case rhino.Token_ASSIGN_AND:
		return "&&=", true
	case rhino.Token_ASSIGN_COALESCE:
		return "??=", true
	case rhino.Token_VOID:
		return "void", true
	case rhino.Token_TYPEOF:
		return "typeof", true
	case rhino.Token_INSTANCEOF:
		return "instanceof", true
	default:
		return "", false
	}
}

// NodeUtil_OpToStrNoFail converts an operator's token value (see Token) to a
// string representation or fails.
//
// operator is the operator's token value to convert; returns the string
// representation. Throws Error if the token value is not an operator.
//
// Java: NodeUtil.opToStrNoFail(Token)
func NodeUtil_OpToStrNoFail(operator rhino.Token) string {
	res := NodeUtil_OpToStr(operator)
	if res == nil {
		panic(jutil.NewJavaLangError("Unknown op " + operator.String()))
	}
	return *res
}

// NodeUtil_RedeclareVarsInsideBranch: given a node tree, finds all the VAR
// declarations in that tree that are not in an inner scope. Then adds a new
// VAR node at the top of the current scope that redeclares them, if
// necessary.
//
// Java: NodeUtil.redeclareVarsInsideBranch(Node)
func NodeUtil_RedeclareVarsInsideBranch(branch rhino.Node) {
	if !NodeUtil_CanContainHoistedVarsDecls(branch) {
		return
	}

	vars := NodeUtil_GetVarsDeclaredInBranch(branch)
	if len(vars) == 0 {
		return
	}

	parent := NodeUtil_GetAddingRoot(branch)
	for _, nameNode := range vars {
		var_ := rhino.IR_Var(rhino.IR_Name(nameNode.GetString()).Srcref(nameNode)).Srcref(nameNode)
		NodeUtil_CopyNameAnnotations(nameNode, var_.GetFirstChild())
		parent.AddChildToFront(var_)
	}
}

// Java: NodeUtil.canContainHoistedVarsDecls(Node)
func NodeUtil_CanContainHoistedVarsDecls(n rhino.Node) bool {
	if NodeUtil_IsStatement(n) {
		switch n.GetToken() {
		case rhino.Token_EXPR_RESULT, rhino.Token_RETURN, rhino.Token_THROW, rhino.Token_BREAK, rhino.Token_CONTINUE,
			rhino.Token_EMPTY, rhino.Token_DEBUGGER:
			return false
		default:
			return true // assume anything else can
		}
	}
	return true
}

// NodeUtil_CopyNameAnnotations copies any annotations that follow a named
// value.
//
// Java: NodeUtil.copyNameAnnotations(Node, Node)
func NodeUtil_CopyNameAnnotations(source rhino.Node, destination rhino.Node) {
	if source.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME) {
		destination.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
	}
}

// NodeUtil_GetAddingRoot gets a Node at the top of the current scope where
// we can add new var declarations as children.
//
// Java: NodeUtil.getAddingRoot(Node)
func NodeUtil_GetAddingRoot(n rhino.Node) rhino.Node {
	var addingRoot rhino.Node
	ancestor := n
crawl_ancestors:
	for {
		ancestor = ancestor.GetParent()
		if ancestor.IsNil() {
			break
		}
		switch ancestor.GetToken() {
		case rhino.Token_SCRIPT, rhino.Token_MODULE_BODY:
			addingRoot = ancestor
			break crawl_ancestors
		case rhino.Token_FUNCTION:
			addingRoot = ancestor.GetLastChild()
			break crawl_ancestors
		default:
			continue crawl_ancestors
		}
	}

	// make sure that the adding root looks ok
	jutil.CheckState(addingRoot.IsBlock() || addingRoot.IsModuleBody() || addingRoot.IsScript())
	jutil.CheckState(!addingRoot.HasChildren() || !addingRoot.GetFirstChild().IsScript())
	return addingRoot
}

// Java: NodeUtil.newDeclaration(Node, Node, Token)
func NodeUtil_NewDeclaration(lhs rhino.Node, rhs rhino.Node, declarationType rhino.Token) rhino.Node {
	if rhs.IsNil() {
		return rhino.IR_Declaration(lhs, declarationType)
	}
	return rhino.IR_DeclarationWithValue(lhs, rhs, declarationType)
}

// NodeUtil_NewQName creates a node representing a qualified name.
//
// name is a qualified name (e.g. "foo" or "foo.bar.baz"); returns a NAME or
// GETPROP node.
//
// Java: NodeUtil.newQName(AbstractCompiler, String)
func NodeUtil_NewQName(compiler AbstractCompiler, name string) rhino.Node {
	endPos := jutil.StringIndexOfChar(name, '.')
	if endPos == -1 {
		endPos = jutil.StringLength(name)
	}

	var qname rhino.Node
	nodeName := jutil.StringSubstring(name, 0, endPos)
	if "this" == nodeName {
		qname = rhino.IR_ThisNode()
	} else if "super" == nodeName {
		qname = rhino.IR_SuperNode()
	} else {
		qname = NodeUtil_NewName(compiler, nodeName)
	}
	qname.SetLength(endPos)

	for startPos := endPos + 1; endPos < jutil.StringLength(name); startPos = endPos + 1 {
		endPos = jutil.StringIndexOfCharFrom(name, '.', startPos)
		if endPos == -1 {
			endPos = jutil.StringLength(name)
		}

		part := jutil.StringSubstring(name, startPos, endPos)
		qname = rhino.IR_Getprop(qname, part)
		if compiler.GetCodingConvention().IsConstantKey(part) {
			qname.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
		}
		qname.SetLength(jutil.StringLength(part))
	}

	return qname
}

// NodeUtil_NewQNameWithOriginalName creates a node representing a qualified
// name, copying over the source location information from the basis node and
// assigning the given original name to the node.
//
// name is a qualified name (e.g. "foo" or "foo.bar.baz"); basisNode is the
// node that represents the name as currently found in the AST; originalName
// is the original name of the item being represented by the NAME node. Used
// for debugging information. Returns a NAME or GETPROP node.
//
// Java: NodeUtil.newQName(AbstractCompiler, String, Node, String)
func NodeUtil_NewQNameWithOriginalName(
	compiler AbstractCompiler, name string, basisNode rhino.Node, originalName string) rhino.Node {
	node := NodeUtil_NewQName(compiler, name)
	node.SrcrefTreeIfMissing(basisNode)
	if nodeOriginalName := node.GetOriginalName(); nodeOriginalName == nil || originalName != *nodeOriginalName {
		// If basisNode already had the correct original name, then it will already be set correctly.
		// Setting it again will force the QName node to have a different property list from all of
		// its children, causing greater memory consumption.
		node.SetOriginalName(jutil.StringPtr(originalName))
	}
	return node
}

// NodeUtil_GetDeclarationFromName attaches nameNode to a new qualified name
// declaration and returns the new qualified declaration.
//
// Returns a new qualified name declaration.
//
// Java: NodeUtil.getDeclarationFromName(Node, Node, Token, JSDocInfo)
func NodeUtil_GetDeclarationFromName(nameNode rhino.Node, value rhino.Node, type_ rhino.Token, info *rhino.JSDocInfo) rhino.Node {
	var result rhino.Node
	if nameNode.IsName() {
		if value.IsNil() {
			result = rhino.IR_Declaration(nameNode, type_)
		} else {
			result = rhino.IR_DeclarationWithValue(nameNode, value, type_)
		}
		result.SetJSDocInfo(info)
	} else if !value.IsNil() {
		result = rhino.IR_ExprResult(rhino.IR_Assign(nameNode, value))
		result.GetFirstChild().SetJSDocInfo(info)
	} else {
		result = rhino.IR_ExprResult(nameNode)
		result.GetFirstChild().SetJSDocInfo(info)
	}
	return result
}

// NodeUtil_NewPropertyAccess creates a property access on the context tree.
//
// Java: NodeUtil.newPropertyAccess(AbstractCompiler, Node, String)
func NodeUtil_NewPropertyAccess(compiler AbstractCompiler, context rhino.Node, name string) rhino.Node {
	propNode := rhino.IR_Getprop(context, name)
	if compiler.GetCodingConvention().IsConstantKey(name) {
		propNode.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
	}
	return propNode
}

// NodeUtil_NewQNameDeclaration creates a node representing a qualified name.
//
// name is a qualified name (e.g. "foo" or "foo.bar.baz"); returns a VAR
// node, or an EXPR_RESULT node containing an ASSIGN or NAME node.
//
// Java: NodeUtil.newQNameDeclaration(AbstractCompiler, String, Node, JSDocInfo)
func NodeUtil_NewQNameDeclaration(
	compiler AbstractCompiler, name string, value rhino.Node, info *rhino.JSDocInfo) rhino.Node {
	return NodeUtil_NewQNameDeclarationWithType(compiler, name, value, info, rhino.Token_VAR)
}

// NodeUtil_NewQNameDeclarationWithType creates a node representing a
// qualified name.
//
// name is a qualified name (e.g. "foo" or "foo.bar.baz"); type_ must be VAR,
// CONST, or LET. Ignored if name is dotted. Returns a VAR/CONST/LET node, or
// an EXPR_RESULT node containing an ASSIGN or NAME node.
//
// Java: NodeUtil.newQNameDeclaration(AbstractCompiler, String, Node, JSDocInfo, Token)
func NodeUtil_NewQNameDeclarationWithType(
	compiler AbstractCompiler, name string, value rhino.Node, info *rhino.JSDocInfo, type_ rhino.Token) rhino.Node {
	jutil.CheckStateMsg1(type_ == rhino.Token_VAR || type_ == rhino.Token_LET || type_ == rhino.Token_CONST, "%s", type_)
	nameNode := NodeUtil_NewQName(compiler, name)
	return NodeUtil_GetDeclarationFromName(nameNode, value, type_, info)
}

// NodeUtil_GetRootOfQualifiedName gets the root node of a qualified name.
// Must be either NAME, THIS or SUPER.
//
// Java: NodeUtil.getRootOfQualifiedName(Node)
func NodeUtil_GetRootOfQualifiedName(qName rhino.Node) rhino.Node {
	for current := qName; ; current = current.GetFirstChild() {
		if current.IsName() || current.IsThis() || current.IsSuper() {
			return current
		}
		jutil.CheckStateMsg1(current.IsGetProp(), "Not a getprop node:  (%s)", current)
	}
}

// NodeUtil_GetRootOfQualifiedNameString gets the root node of a string
// representing a qualified name.
//
// Java: NodeUtil.getRootOfQualifiedName(String)
func NodeUtil_GetRootOfQualifiedNameString(qName string) string {
	dot := jutil.StringIndexOfChar(qName, '.')
	if dot == -1 {
		return qName
	}
	return jutil.StringSubstring(qName, 0, dot)
}

// Java: NodeUtil.newName(AbstractCompiler, String)
func NodeUtil_NewName(compiler AbstractCompiler, name string) rhino.Node {
	nameNode := rhino.IR_Name(name)
	nameNode.SetLength(jutil.StringLength(name))
	if compiler.GetCodingConvention().IsConstant(name) {
		nameNode.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
	}
	return nameNode
}

// NodeUtil_NewNameWithSrcref creates a new node representing an *existing*
// name, copying over the source location information from the basis node.
//
// name is the name for the new NAME node; srcref is the node that represents
// the name as currently found in the AST. Returns the node created.
//
// Java: NodeUtil.newName(AbstractCompiler, String, Node)
func NodeUtil_NewNameWithSrcref(compiler AbstractCompiler, name string, srcref rhino.Node) rhino.Node {
	return NodeUtil_NewName(compiler, name).Srcref(srcref)
}

// NodeUtil_NewNameWithOriginalName creates a new node representing an
// *existing* name, copying over the source location information from the
// basis node and assigning the given original name to the node.
//
// name is the name for the new NAME node; basisNode is the node that
// represents the name as currently found in the AST; originalName is the
// original name of the item being represented by the NAME node. Used for
// debugging information. Returns the node created.
//
// Java: NodeUtil.newName(AbstractCompiler, String, Node, String)
func NodeUtil_NewNameWithOriginalName(
	compiler AbstractCompiler, name string, basisNode rhino.Node, originalName string) rhino.Node {
	nameNode := NodeUtil_NewNameWithSrcref(compiler, name, basisNode)
	nameNode.SetOriginalNameString(originalName)
	return nameNode
}

// NodeUtil_IsLatin tests if all characters in the string are in the Basic
// Latin (aka ASCII) character set - that they have UTF-16 values equal to or
// below 0x7f. This check can find which identifiers with Unicode characters
// need to be escaped in order to allow resulting files to be processed by
// non-Unicode aware UNIX tools and editors. See
// http://en.wikipedia.org/wiki/Latin_characters_in_Unicode for more on Basic
// Latin.
//
// s is the string to be checked for ASCII-goodness. Returns true if all
// characters in the string are in Basic Latin set.
//
// Java: NodeUtil.isLatin(String)
func NodeUtil_IsLatin(s string) bool {
	// Go-only: language-adapter — a UTF-16 char above 0x7f exists iff the
	// WTF-8 encoding has a byte above 0x7f, so the bytes are scanned directly.
	len_ := len(s)
	for index := 0; index < len_; index++ {
		c := s[index]
		if uint16(c) > NodeUtil_LARGEST_BASIC_LATIN {
			return false
		}
	}
	return true
}

// NodeUtil_IsValidSimpleName determines whether the given name is a valid
// variable name.
//
// Java: NodeUtil.isValidSimpleName(String)
func NodeUtil_IsValidSimpleName(name string) bool {
	return rhino.TokenStream_IsJSIdentifier(name) &&
		!rhino.TokenStream_IsKeyword(name) &&
		// no Unicode escaped characters - some browsers are less tolerant
		// of Unicode characters that might be valid according to the
		// language spec.
		// Note that by this point, Unicode escapes have been converted
		// to UTF-16 characters, so we're only searching for character
		// values, not escapes.
		NodeUtil_IsLatin(name)
}

// Deprecated: use NodeUtil_IsValidQualifiedNameFeatureSet(mode.ToFeatureSet(), name).
//
// Java: NodeUtil.isValidQualifiedName(LanguageMode, String)
func NodeUtil_IsValidQualifiedName(mode CompilerOptions_LanguageMode, name string) bool {
	return NodeUtil_IsValidQualifiedNameFeatureSet(mode.ToFeatureSet(), name)
}

// NodeUtil_IsValidQualifiedNameFeatureSet determines whether the given name
// is a valid qualified name.
//
// Java: NodeUtil.isValidQualifiedName(FeatureSet, String)
func NodeUtil_IsValidQualifiedNameFeatureSet(mode *parser.FeatureSet, name string) bool {
	if strings.HasSuffix(name, ".") || strings.HasPrefix(name, ".") {
		return false
	}

	// Go-only: language-adapter — Splitter.on('.').splitToList(name) keeps
	// empty parts, exactly like strings.Split with a one-character separator.
	parts := strings.Split(name, ".")
	for _, part := range parts {
		if !NodeUtil_IsValidPropertyName(mode, part) {
			return false
		}
	}
	return NodeUtil_IsValidSimpleName(parts[0])
}

// NodeUtil_IsValidPropertyName determines whether the given name can appear
// on the right side of the dot operator. Many properties (like reserved
// words) cannot, in ES3.
//
// Java: NodeUtil.isValidPropertyName(FeatureSet, String)
func NodeUtil_IsValidPropertyName(mode *parser.FeatureSet, name string) bool {
	if NodeUtil_IsValidSimpleName(name) {
		return true
	} else {
		return mode.Has(parser.FeatureSet_Feature_KEYWORDS_AS_PROPERTIES) && rhino.TokenStream_IsKeyword(name)
	}
}

// NodeUtil_VarCollector collects the names declared by VARs.
//
// Java: NodeUtil.VarCollector
type NodeUtil_VarCollector struct {
	vars *jutil.LinkedHashMap[string, rhino.Node]
}

// Java: NodeUtil.VarCollector()
func NewNodeUtil_VarCollector() *NodeUtil_VarCollector {
	return &NodeUtil_VarCollector{vars: jutil.NewLinkedHashMap[string, rhino.Node]()}
}

// Java: NodeUtil.VarCollector.visit(Node)
func (c *NodeUtil_VarCollector) Visit(n rhino.Node) {
	if n.IsVar() {
		NodeUtil_VisitLhsNodesInNode(n, func(name rhino.Node) { c.vars.PutIfAbsent(name.GetString(), name) })
	}
}

// NodeUtil_GetVarsDeclaredInBranch retrieves vars declared in the current
// node tree, excluding descent scopes.
//
// Java: NodeUtil.getVarsDeclaredInBranch(Node)
func NodeUtil_GetVarsDeclaredInBranch(root rhino.Node) []rhino.Node {
	collector := NewNodeUtil_VarCollector()
	NodeUtil_VisitPreOrderWithPred(root, collector, NodeUtil_MATCH_NOT_FUNCTION)
	return slices.Collect(collector.vars.Values())
}

// Java: NodeUtil.getLhsNodesHelper(Node, Consumer<Node>)
func NodeUtil_GetLhsNodesHelper(n rhino.Node, consumer func(rhino.Node)) {
	switch n.GetToken() {
	case rhino.Token_IMPORT:
		NodeUtil_GetLhsNodesHelper(n.GetFirstChild(), consumer)
		NodeUtil_GetLhsNodesHelper(n.GetSecondChild(), consumer)
		return
	case rhino.Token_VAR, rhino.Token_CONST, rhino.Token_LET, rhino.Token_OBJECT_PATTERN, rhino.Token_ARRAY_PATTERN,
		rhino.Token_PARAM_LIST, rhino.Token_IMPORT_SPECS:
		for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
			NodeUtil_GetLhsNodesHelper(child, consumer)
		}
		return
	case rhino.Token_DESTRUCTURING_LHS, rhino.Token_DEFAULT_VALUE, rhino.Token_CATCH, rhino.Token_ITER_REST,
		rhino.Token_OBJECT_REST, rhino.Token_CAST:
		NodeUtil_GetLhsNodesHelper(n.GetFirstChild(), consumer)
	case rhino.Token_IMPORT_SPEC, rhino.Token_COMPUTED_PROP, rhino.Token_STRING_KEY:
		NodeUtil_GetLhsNodesHelper(n.GetLastChild(), consumer)
	case rhino.Token_NAME, rhino.Token_IMPORT_STAR:
		consumer(n)
	case rhino.Token_GETPROP, rhino.Token_GETELEM:
		// Not valid in declarations but may appear in assignments.
		consumer(n)
	case rhino.Token_EMPTY: // nothing to do.
	case rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
		// Enhanced for loops assign to variables in their first child
		// e.g.
		// for (some.prop in someObj) {...
		// for ({a, b} of someIterable) {...
		NodeUtil_GetLhsNodesHelper(n.GetFirstChild(), consumer)
	default:
		if NodeUtil_IsAssignmentOp(n) {
			NodeUtil_GetLhsNodesHelper(n.GetFirstChild(), consumer)
		} else {
			panic(jutil.NewIllegalStateExceptionMsg("Invalid node in lhs: " + n.String()))
		}
	}
}

// NodeUtil_VisitLhsNodesInNode retrieves lhs nodes declared or assigned in a
// given assigning parent node.
//
// An assigning parent node is one that assigns a value to one or more LHS
// nodes.
//
// Java: NodeUtil.visitLhsNodesInNode(Node, Consumer<Node>)
func NodeUtil_VisitLhsNodesInNode(assigningParent rhino.Node, consumer func(rhino.Node)) {
	jutil.CheckArgumentMsg1(
		NodeUtil_IsNameDeclaration(assigningParent) || // e.g. `var [a,b,c];`
			assigningParent.IsParamList() || // e.g. `function foo([a,b,c]) {}`
			NodeUtil_IsAssignmentOp(assigningParent) || // e.g. `[a,b,c] = [4,5,6]`
			assigningParent.IsCatch() ||
			assigningParent.IsDestructuringLhs() ||
			assigningParent.IsDefaultValue() ||
			assigningParent.IsImport() ||
			// enhanced for loops assign to loop variables
			NodeUtil_IsEnhancedFor(assigningParent),
		"%s",
		assigningParent)
	NodeUtil_GetLhsNodesHelper(assigningParent, consumer)
}

// NodeUtil_VisitLhsNodesInDestructuringPattern retrieves lhs nodes declared
// or assigned in a given destructuring pattern node.
//
// Java: NodeUtil.visitLhsNodesInDestructuringPattern(Node, Consumer<Node>)
func NodeUtil_VisitLhsNodesInDestructuringPattern(
	destructuringPattern rhino.Node, consumer func(rhino.Node)) {
	jutil.CheckArgumentMsg(
		destructuringPattern.IsDestructuringPattern(), "Must be a destructuring pattern node.")
	NodeUtil_GetLhsNodesHelper(destructuringPattern, consumer)
}

// NodeUtil_IsObjectDefinePropertiesDefinition returns true if the node is a
// definition with Object.defineProperties.
//
// Java: NodeUtil.isObjectDefinePropertiesDefinition(Node)
func NodeUtil_IsObjectDefinePropertiesDefinition(n rhino.Node) bool {
	// We intentionally don't check optional Object.defineProperties?.() because we expect it to be
	// defined or polyfill-ed, and to avoid changing the current invariant that this method returns
	// true only for non-optional `Object.defineProperties()`.
	if !n.IsCall() || !n.HasXChildren(3) {
		return false
	}
	getprop := n.GetFirstChild()
	if !getprop.IsGetProp() {
		return false
	}
	return getprop.GetString() == "defineProperties" &&
		NodeUtil_IsKnownGlobalObjectReference(getprop.GetFirstChild())
}

// NodeUtil_IsObjectCreateWithPropertiesDefinition returns true if the node is
// a 2-argument Object.create(proto, descriptors).
//
// Java: NodeUtil.isObjectCreateWithPropertiesDefinition(Node)
func NodeUtil_IsObjectCreateWithPropertiesDefinition(n rhino.Node) bool {
	if !n.IsCall() || !n.HasXChildren(3) {
		return false
	}
	getprop := n.GetFirstChild()
	return getprop.IsGetProp() &&
		getprop.GetString() == "create" &&
		NodeUtil_IsKnownGlobalObjectReference(getprop.GetFirstChild())
}

// Java: NodeUtil.GLOBAL_OBJECT
var NodeUtil_GLOBAL_OBJECT = rhino.QualifiedName_Of("$jscomp.global.Object")

// Java: NodeUtil.GLOBAL_OBJECT_MANGLED
var NodeUtil_GLOBAL_OBJECT_MANGLED = rhino.QualifiedName_Of("$jscomp$global.Object")

// Java: NodeUtil.isKnownGlobalObjectReference(Node)
func NodeUtil_IsKnownGlobalObjectReference(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_NAME:
		return n.GetString() == "Object"
	case rhino.Token_GETPROP:
		return NodeUtil_GLOBAL_OBJECT.Matches(n) || NodeUtil_GLOBAL_OBJECT_MANGLED.Matches(n)
	default:
		return false
	}
}

// NodeUtil_IsObjectDefinePropertyDefinition returns true if the node is a
// definition with Object.defineProperty.
//
// Java: NodeUtil.isObjectDefinePropertyDefinition(Node)
func NodeUtil_IsObjectDefinePropertyDefinition(n rhino.Node) bool {
	if !n.IsCall() || !n.HasXChildren(4) {
		return false
	}
	getprop := n.GetFirstChild()
	if !getprop.IsGetProp() {
		return false
	}
	return getprop.GetString() == "defineProperty" &&
		NodeUtil_IsKnownGlobalObjectReference(getprop.GetFirstChild())
}

// NodeUtil_GetObjectDefinedPropertiesKeys returns a list of STRING_KEY
// properties defined by a Object.defineProperties(o, {...}) call.
//
// Java: NodeUtil.getObjectDefinedPropertiesKeys(Node)
func NodeUtil_GetObjectDefinedPropertiesKeys(definePropertiesCall rhino.Node) []rhino.Node {
	jutil.CheckArgument(NodeUtil_IsObjectDefinePropertiesDefinition(definePropertiesCall))
	var properties []rhino.Node
	objectLiteral := definePropertiesCall.GetLastChild()
	for key := objectLiteral.GetFirstChild(); !key.IsNil(); key = key.GetNext() {
		if !key.IsStringKey() {
			continue
		}
		properties = append(properties, key)
	}
	return properties
}

// NodeUtil_IsPrototypePropertyDeclaration returns true if the node an
// assignment to a prototype property of some constructor.
//
// Java: NodeUtil.isPrototypePropertyDeclaration(Node)
func NodeUtil_IsPrototypePropertyDeclaration(n rhino.Node) bool {
	return NodeUtil_IsExprAssign(n) && NodeUtil_IsPrototypeProperty(n.GetFirstFirstChild())
}

// NodeUtil_IsPrototypeProperty returns whether the node represents a
// qualified prototype property.
//
// Java: NodeUtil.isPrototypeProperty(Node)
func NodeUtil_IsPrototypeProperty(n rhino.Node) bool {
	if !n.IsGetProp() {
		return false
	}
	recv := n.GetFirstChild()
	return recv.IsGetProp() && recv.GetString() == "prototype"
}

// NodeUtil_IsPrototypeMethod returns whether the node represents a prototype
// method.
//
// Java: NodeUtil.isPrototypeMethod(Node)
func NodeUtil_IsPrototypeMethod(n rhino.Node) bool {
	if !n.IsFunction() {
		return false
	}
	assignNode := n.GetParent()
	if !assignNode.IsAssign() {
		return false
	}
	return NodeUtil_IsPrototypePropertyDeclaration(assignNode.GetParent())
}

// Java: NodeUtil.isPrototypeAssignment(Node)
func NodeUtil_IsPrototypeAssignment(getProp rhino.Node) bool {
	if !getProp.IsGetProp() {
		return false
	}
	parent := getProp.GetParent()
	return parent.IsAssign() &&
		getProp.IsFirstChildOf(parent) &&
		getProp.GetString() == "prototype"
}

// NodeUtil_IsPropertyTest determines whether this node is testing for the
// existence of a property. If true, we will not emit warnings about a
// missing property.
//
// propAccess is the GETPROP or GETELEM being tested.
//
// Java: NodeUtil.isPropertyTest(AbstractCompiler, Node)
func NodeUtil_IsPropertyTest(compiler AbstractCompiler, propAccess rhino.Node) bool {
	parent := propAccess.GetParent()
	switch parent.GetToken() {
	case rhino.Token_CALL:
		return parent.GetFirstChild() != propAccess &&
			compiler.GetCodingConvention().IsPropertyTestFunction(parent)
	case rhino.Token_IF, rhino.Token_WHILE, rhino.Token_DO, rhino.Token_FOR, rhino.Token_FOR_IN:
		return NodeUtil_GetConditionExpression(parent) == propAccess

	case rhino.Token_INSTANCEOF, rhino.Token_TYPEOF, rhino.Token_AND, rhino.Token_OR, rhino.Token_COALESCE,
		rhino.Token_OPTCHAIN_GETPROP:
		return true

	case rhino.Token_NE, rhino.Token_SHNE:
		var other rhino.Node
		if parent.GetFirstChild() == propAccess {
			other = parent.GetSecondChild()
		} else {
			other = parent.GetFirstChild()
		}
		return NodeUtil_IsUndefined(other) || (parent.IsNE() && other.IsNull())
	case rhino.Token_OPTCHAIN_CALL, rhino.Token_OPTCHAIN_GETELEM, rhino.Token_HOOK:
		return parent.GetFirstChild() == propAccess
	case rhino.Token_NOT:
		return parent.GetParent().IsOr() && parent.GetParent().GetFirstChild() == parent
	case rhino.Token_CAST:
		return NodeUtil_IsPropertyTest(compiler, parent)
	default:
		return false
	}
}

// Java: NodeUtil.isPropertyAbsenceTest(Node)
func NodeUtil_IsPropertyAbsenceTest(propAccess rhino.Node) bool {
	parent := propAccess.GetParent()
	switch parent.GetToken() {
	case rhino.Token_EQ, rhino.Token_SHEQ:
		var other rhino.Node
		if parent.GetFirstChild() == propAccess {
			other = parent.GetSecondChild()
		} else {
			other = parent.GetFirstChild()
		}
		return NodeUtil_IsUndefined(other) || (parent.IsEQ() && other.IsNull())
	default:
		return false
	}
}

// NodeUtil_GetPrototypeClassName: qName is a qualified name node
// representing a class prototype, or a property on that prototype, e.g.
// foo.Bar.prototype, or foo.Bar.prototype.toString. Returns the class name
// part of a qualified prototype name, e.g. foo.Bar.
//
// Java: NodeUtil.getPrototypeClassName(Node)
func NodeUtil_GetPrototypeClassName(qName rhino.Node) rhino.Node {
	if !qName.IsGetProp() {
		return rhino.NilNode
	}
	if qName.GetString() == "prototype" {
		return qName.GetFirstChild()
	}
	recv := qName.GetFirstChild()
	if recv.IsGetProp() && recv.GetString() == "prototype" {
		return recv.GetFirstChild()
	}
	return rhino.NilNode
}

// NodeUtil_GetPrototypePropertyName returns the string property name part of
// a qualified prototype name.
//
// Java: NodeUtil.getPrototypePropertyName(Node)
func NodeUtil_GetPrototypePropertyName(qName rhino.Node) string {
	qNameStr, ok := qName.GetQualifiedNameString()
	if !ok {
		panic(jutil.NewNullPointerException())
	}
	prototypeIdx := jutil.StringLastIndexOf(qNameStr, ".prototype.")
	memberIndex := prototypeIdx + jutil.StringLength(".prototype") + 1
	return jutil.StringSubstringFrom(qNameStr, memberIndex)
}

// NodeUtil_NewUndefinedNode creates a node for an empty result expression:
// "void 0".
//
// Java: NodeUtil.newUndefinedNode(Node)
func NodeUtil_NewUndefinedNode(srcReferenceNode rhino.Node) rhino.Node {
	node := rhino.IR_VoidNode(rhino.IR_Number(0))
	if !srcReferenceNode.IsNil() {
		node.SrcrefTree(srcReferenceNode)
	}
	return node
}

// NodeUtil_NewVarNode creates a VAR node containing the given name and
// initial value expression.
//
// Java: NodeUtil.newVarNode(String, Node)
func NodeUtil_NewVarNode(name string, value rhino.Node) rhino.Node {
	lhs := rhino.IR_Name(name)
	if !value.IsNil() {
		lhs.Srcref(value)
	}
	return NodeUtil_NewVarNodeWithLhs(lhs, value)
}

// NodeUtil_NewVarNodeWithLhs creates a VAR node containing the given lhs
// (name or destructuring pattern) and initial value expression.
//
// Java: NodeUtil.newVarNode(Node, Node)
func NodeUtil_NewVarNodeWithLhs(lhs rhino.Node, value rhino.Node) rhino.Node {
	if lhs.IsDestructuringPattern() {
		// TODO(b/314005766): Ensure that destructuring declarations are generated as normalized (i.e
		// individual names get their own stub declarations)
		rhino.CheckNotNull(value)
		return rhino.IR_Var(rhino.NewNodeWith2Children(rhino.Token_DESTRUCTURING_LHS, lhs, value).Srcref(lhs)).Srcref(lhs)
	} else {
		jutil.CheckState(lhs.IsName() && !lhs.HasChildren())
		if !value.IsNil() {
			lhs.AddChildToBack(value)
		}
		return rhino.IR_Var(lhs).Srcref(lhs)
	}
}

// Java: NodeUtil.emptyFunction()
func NodeUtil_EmptyFunction() rhino.Node {
	return rhino.IR_Function(rhino.IR_Name(""), rhino.IR_ParamList(), rhino.IR_Block())
}

// NodeUtil_MatchNameNode is a predicate for matching name nodes with the
// specified node.
//
// Java: NodeUtil.MatchNameNode
type NodeUtil_MatchNameNode struct {
	name string
}

// Java: NodeUtil.MatchNameNode(String)
func NewNodeUtil_MatchNameNode(name string) *NodeUtil_MatchNameNode {
	return &NodeUtil_MatchNameNode{name: name}
}

// Java: NodeUtil.MatchNameNode.apply(Node)
func (m *NodeUtil_MatchNameNode) Apply(n rhino.Node) bool {
	return n.IsName() && n.GetString() == m.name
}

// NodeUtil_MatchNodeType is a predicate for matching nodes with the
// specified type.
//
// Java: NodeUtil.MatchNodeType
type NodeUtil_MatchNodeType struct {
	type_ rhino.Token
}

// Java: NodeUtil.MatchNodeType(Token)
func NewNodeUtil_MatchNodeType(type_ rhino.Token) *NodeUtil_MatchNodeType {
	return &NodeUtil_MatchNodeType{type_: type_}
}

// Java: NodeUtil.MatchNodeType.apply(Node)
func (m *NodeUtil_MatchNodeType) Apply(n rhino.Node) bool {
	return n.GetToken() == m.type_
}

// NodeUtil_MatchDeclaration is a predicate for matching var, let, const,
// class or function declarations.
//
// Java: NodeUtil.MatchDeclaration
type NodeUtil_MatchDeclaration struct {
	_ byte // Go-only: language-adapter — distinct identity
}

// Java: NodeUtil.MatchDeclaration()
func NewNodeUtil_MatchDeclaration() *NodeUtil_MatchDeclaration {
	return &NodeUtil_MatchDeclaration{}
}

// Java: NodeUtil.MatchDeclaration.apply(Node)
func (m *NodeUtil_MatchDeclaration) Apply(n rhino.Node) bool {
	return NodeUtil_IsDeclaration(n)
}

// Java: NodeUtil.MATCH_NOT_FUNCTION
var NodeUtil_MATCH_NOT_FUNCTION = func(n rhino.Node) bool { return !n.IsFunction() }

// NodeUtil_MATCH_ANYTHING_BUT_NON_ARROW_FUNCTION is a predicate for matching
// anything except for a non-arrow function.
//
// Useful to avoid traversing into scopes that don't share the same values
// for this, super, or arguments.
//
// Java: NodeUtil.MATCH_ANYTHING_BUT_NON_ARROW_FUNCTION
var NodeUtil_MATCH_ANYTHING_BUT_NON_ARROW_FUNCTION = func(n rhino.Node) bool { return !NodeUtil_IsNonArrowFunction(n) }

// NodeUtil_IsShallowStatementTree: whether the given node's subtree may
// contain statements, excepting nested functions.
//
// This is useful for traversing all statements in a given script or
// function without traversing into nested functions, so it will return true
// for e.g. a BLOCK or SWITCH statement.
//
// Java: NodeUtil.isShallowStatementTree(Node)
func NodeUtil_IsShallowStatementTree(n rhino.Node) bool {
	return n.IsNil() ||
		NodeUtil_IsControlStructure(n) ||
		NodeUtil_IsStatementBlock(n) ||
		n.IsSwitchBody()
}

// NodeUtil_MatchShallowStatement is a predicate for matching statements
// without exiting the current scope.
//
// Java: NodeUtil.MatchShallowStatement
type NodeUtil_MatchShallowStatement struct {
	_ byte // Go-only: language-adapter — distinct identity
}

// Java: NodeUtil.MatchShallowStatement()
func NewNodeUtil_MatchShallowStatement() *NodeUtil_MatchShallowStatement {
	return &NodeUtil_MatchShallowStatement{}
}

// Java: NodeUtil.MatchShallowStatement.apply(Node)
func (m *NodeUtil_MatchShallowStatement) Apply(n rhino.Node) bool {
	parent := n.GetParent()
	return n.IsRoot() || n.IsBlock() || (!n.IsFunction() && NodeUtil_IsShallowStatementTree(parent))
}

// NodeUtil_GetNodeTypeReferenceCount finds the number of times a type is
// referenced within the node tree.
//
// Java: NodeUtil.getNodeTypeReferenceCount(Node, Token, Predicate<Node>)
func NodeUtil_GetNodeTypeReferenceCount(
	node rhino.Node, type_ rhino.Token, traverseChildrenPred func(rhino.Node) bool) int {
	return NodeUtil_GetCount(node, NewNodeUtil_MatchNodeType(type_).Apply, traverseChildrenPred)
}

// NodeUtil_IsNameReferencedWithPred: whether a simple name is referenced
// within the node tree.
//
// Java: NodeUtil.isNameReferenced(Node, String, Predicate<Node>)
func NodeUtil_IsNameReferencedWithPred(node rhino.Node, name string, traverseChildrenPred func(rhino.Node) bool) bool {
	return NodeUtil_Has(node, NewNodeUtil_MatchNameNode(name).Apply, traverseChildrenPred)
}

// NodeUtil_IsNameReferenced: whether a simple name is referenced within the
// node tree.
//
// Java: NodeUtil.isNameReferenced(Node, String)
func NodeUtil_IsNameReferenced(node rhino.Node, name string) bool {
	return NodeUtil_IsNameReferencedWithPred(node, name, nodeUtil_alwaysTrue)
}

// NodeUtil_GetNameReferenceCount finds the number of times a simple name is
// referenced within the node tree.
//
// Java: NodeUtil.getNameReferenceCount(Node, String)
func NodeUtil_GetNameReferenceCount(node rhino.Node, name string) int {
	return NodeUtil_GetCount(node, NewNodeUtil_MatchNameNode(name).Apply, nodeUtil_alwaysTrue)
}

// nodeUtil_alwaysTrue is Guava's Predicates.alwaysTrue() for nodes.
func nodeUtil_alwaysTrue(rhino.Node) bool { return true }

// NodeUtil_Has returns whether the predicate is true for the node or any of
// its descendants.
//
// Java: NodeUtil.has(Node, Predicate<Node>, Predicate<Node>)
func NodeUtil_Has(node rhino.Node, pred func(rhino.Node) bool, traverseChildrenPred func(rhino.Node) bool) bool {
	if pred(node) {
		return true
	}

	if !traverseChildrenPred(node) {
		return false
	}

	for c := node.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if NodeUtil_Has(c, pred, traverseChildrenPred) {
			return true
		}
	}

	return false
}

// NodeUtil_FindPreorder returns the first Node matching the given pred via a
// pre-order traversal.
//
// Java: NodeUtil.findPreorder(Node, Predicate<Node>, Predicate<Node>)
func NodeUtil_FindPreorder(
	node rhino.Node, pred func(rhino.Node) bool, traverseChildrenPred func(rhino.Node) bool) rhino.Node {
	if pred(node) {
		return node
	}

	if !traverseChildrenPred(node) {
		return rhino.NilNode
	}

	for c := node.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		result := NodeUtil_FindPreorder(c, pred, traverseChildrenPred)
		if !result.IsNil() {
			return result
		}
	}

	return rhino.NilNode
}

// NodeUtil_GetCount returns the number of times the predicate is true for
// the node or any of its descendants.
//
// Java: NodeUtil.getCount(Node, Predicate<Node>, Predicate<Node>)
func NodeUtil_GetCount(n rhino.Node, pred func(rhino.Node) bool, traverseChildrenPred func(rhino.Node) bool) int {
	total := 0

	if pred(n) {
		total++
	}

	if traverseChildrenPred(n) {
		for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			total += NodeUtil_GetCount(c, pred, traverseChildrenPred)
		}
	}

	return total
}

// NodeUtil_Visitor is the interface for use with the visit method.
//
// See NodeUtil_VisitPreOrder.
//
// Java: NodeUtil.Visitor
type NodeUtil_Visitor interface {
	// Java: NodeUtil.Visitor.visit(Node)
	Visit(node rhino.Node)
}

// NodeUtil_VisitorFunc adapts a function (a Java lambda) to NodeUtil_Visitor.
//
// Go-only: language-adapter — Java lambdas implementing the functional
// interface.
type NodeUtil_VisitorFunc func(node rhino.Node)

// Visit calls f.
func (f NodeUtil_VisitorFunc) Visit(node rhino.Node) { f(node) }

// NodeUtil_VisitPreOrder is a pre-order traversal, calling Visitor.visit for
// each decendent.
//
// Java: NodeUtil.visitPreOrder(Node, Visitor)
func NodeUtil_VisitPreOrder(node rhino.Node, visitor NodeUtil_Visitor) {
	NodeUtil_VisitPreOrderWithPred(node, visitor, nodeUtil_alwaysTrue)
}

// NodeUtil_VisitPreOrderWithPred is a pre-order traversal, calling
// Visitor.visit for each node in the tree. Children of nodes that do not
// match the predicate will not be visited.
//
// Java: NodeUtil.visitPreOrder(Node, Visitor, Predicate<Node>)
func NodeUtil_VisitPreOrderWithPred(
	node rhino.Node, visitor NodeUtil_Visitor, traverseChildrenPred func(rhino.Node) bool) {
	visitor.Visit(node)

	if traverseChildrenPred(node) {
		for c := node.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			NodeUtil_VisitPreOrderWithPred(c, visitor, traverseChildrenPred)
		}
	}
}

// NodeUtil_VisitPostOrder is a post-order traversal, calling Visitor.visit
// for each decendent.
//
// Java: NodeUtil.visitPostOrder(Node, Visitor)
func NodeUtil_VisitPostOrder(node rhino.Node, visitor NodeUtil_Visitor) {
	NodeUtil_VisitPostOrderWithPred(node, visitor, nodeUtil_alwaysTrue)
}

// NodeUtil_VisitPostOrderWithPred is a post-order traversal, calling
// Visitor.visit for each node in the tree. Children of nodes that do not
// match the predicate will not be visited.
//
// Java: NodeUtil.visitPostOrder(Node, Visitor, Predicate<Node>)
func NodeUtil_VisitPostOrderWithPred(
	node rhino.Node, visitor NodeUtil_Visitor, traverseChildrenPred func(rhino.Node) bool) {
	if traverseChildrenPred(node) {
		for c := node.GetFirstChild(); !c.IsNil(); {
			next := c.GetNext()
			NodeUtil_VisitPostOrderWithPred(c, visitor, traverseChildrenPred)
			c = next
		}
	}

	visitor.Visit(node)
}

// NodeUtil_PreOrderIterableWithPred creates an Iterable over the given node
// and all descendents. Nodes are given in depth-first pre-order (document
// source order,
// https://www.w3.org/TR/DOM-Level-2-Traversal-Range/glossary.html#dt-documentorder).
//
// This has the benefit over the visitor patten that it is iterative and can
// be used with Java 8 Stream API for searching, filtering, transforms, and
// other functional processing.
//
// The given predicate determines whether a node's children will be iterated
// over. If a node does not match the predicate, none of its children will be
// visited.
//
// root is the root of the tree; travserseNodePredicate matches nodes in the
// tree whose children should be traversed.
//
// Go-only: language-adapter — the Iterable is an iter.Seq; each traversal
// starts a new PreOrderIterator, like each iterator() call of the Java
// lambda.
//
// Java: NodeUtil.preOrderIterable(Node, Predicate<Node>)
func NodeUtil_PreOrderIterableWithPred(root rhino.Node, travserseNodePredicate func(rhino.Node) bool) iter.Seq[rhino.Node] {
	return func(yield func(rhino.Node) bool) {
		if root.IsNil() {
			panic(jutil.NewNullPointerException())
		}
		it := NodeUtil_PreOrderIterator{traverseNodePredicate: travserseNodePredicate, current: root}
		for it.HasNext() {
			if !yield(it.Next()) {
				return
			}
		}
	}
}

// NodeUtil_PreOrderIterable is the same as NodeUtil_PreOrderIterableWithPred
// but iterates over all nodes in the tree without exception.
//
// Java: NodeUtil.preOrderIterable(Node)
func NodeUtil_PreOrderIterable(root rhino.Node) iter.Seq[rhino.Node] {
	return NodeUtil_PreOrderIterableWithPred(
		root,
		// Traverse all nodes.
		nodeUtil_alwaysTrue)
}

// NodeUtil_PreOrderIterator is the utility class for preOrderIterable.
// Iterates over nodes in tree in depth-first pre-order.
//
// Java: NodeUtil.PreOrderIterator
type NodeUtil_PreOrderIterator struct {
	traverseNodePredicate func(rhino.Node) bool

	current rhino.Node

	// Go-only: language-adapter — AbstractIterator's state: whether next
	// holds a computed element, and whether endOfData() was reached.
	next    rhino.Node
	hasNext bool
	done    bool
}

// Java: NodeUtil.PreOrderIterator(Node, Predicate<Node>)
func NewNodeUtil_PreOrderIterator(root rhino.Node, traverseNodePredicate func(rhino.Node) bool) *NodeUtil_PreOrderIterator {
	if root.IsNil() {
		panic(jutil.NewNullPointerException())
	}
	return &NodeUtil_PreOrderIterator{traverseNodePredicate: traverseNodePredicate, current: root}
}

// HasNext is AbstractIterator.hasNext().
func (it *NodeUtil_PreOrderIterator) HasNext() bool {
	if it.done {
		return false
	}
	if !it.hasNext {
		next, ok := it.ComputeNext()
		if !ok {
			it.done = true
			return false
		}
		it.next = next
		it.hasNext = true
	}
	return true
}

// Next is AbstractIterator.next().
func (it *NodeUtil_PreOrderIterator) Next() rhino.Node {
	if !it.HasNext() {
		panic(jutil.NewNoSuchElementException())
	}
	it.hasNext = false
	result := it.next
	it.next = rhino.NilNode
	return result
}

// ComputeNext returns the next node, or ok=false for endOfData().
//
// Java: NodeUtil.PreOrderIterator.computeNext()
func (it *NodeUtil_PreOrderIterator) ComputeNext() (next rhino.Node, ok bool) {
	if it.current.IsNil() {
		return rhino.NilNode, false
	}

	returnValue := it.current
	it.current = it.CalculateNextNode(returnValue)
	return returnValue, true
}

// Java: NodeUtil.PreOrderIterator.calculateNextNode(Node)
func (it *NodeUtil_PreOrderIterator) CalculateNextNode(currentNode rhino.Node) rhino.Node {
	if currentNode.IsNil() {
		panic(jutil.NewNullPointerException())
	}

	// If node does not match the predicate, do not descend into it.
	if it.traverseNodePredicate(currentNode) {
		// In prefix order, the next node is the leftmost child.
		if currentNode.HasChildren() {
			return currentNode.GetFirstChild()
		}
	}

	// If currentNode doesn't have children, it is a leaf node.

	// To find the next node, walk up the ancestry chain (including current node) and return the
	// first sibling we see.
	// If we don't find one, we're done.

	for !currentNode.IsNil() {
		next := currentNode.GetNext()
		if !next.IsNil() {
			return next
		}

		currentNode = currentNode.GetParent()
	}

	return rhino.NilNode
}

// NodeUtil_IsExportFrom returns whether an EXPORT node has a from clause.
//
// Java: NodeUtil.isExportFrom(Node)
func NodeUtil_IsExportFrom(n rhino.Node) bool {
	jutil.CheckArgument(n.IsExport())
	return n.HasTwoChildren()
}

// NodeUtil_HasFinally returns whether a TRY node has a finally block.
//
// Java: NodeUtil.hasFinally(Node)
func NodeUtil_HasFinally(n rhino.Node) bool {
	jutil.CheckArgument(n.IsTry())
	return n.HasXChildren(3)
}

// NodeUtil_GetCatchBlock returns the BLOCK node containing the CATCH node
// (if any) of a TRY.
//
// Java: NodeUtil.getCatchBlock(Node)
func NodeUtil_GetCatchBlock(n rhino.Node) rhino.Node {
	jutil.CheckArgument(n.IsTry())
	return n.GetSecondChild()
}

// NodeUtil_HasCatchHandler returns whether BLOCK (from a TRY node) contains
// a CATCH.
//
// See NodeUtil_GetCatchBlock.
//
// Java: NodeUtil.hasCatchHandler(Node)
func NodeUtil_HasCatchHandler(n rhino.Node) bool {
	jutil.CheckArgument(n.IsBlock())
	return n.HasChildren() && n.GetFirstChild().IsCatch()
}

// NodeUtil_GetFunctionParameters: fnNode is the function. Returns the Node
// containing the Function parameters.
//
// Java: NodeUtil.getFunctionParameters(Node)
func NodeUtil_GetFunctionParameters(fnNode rhino.Node) rhino.Node {
	jutil.CheckArgument(fnNode.IsFunction())
	return fnNode.GetSecondChild()
}

// NodeUtil_HasNonSimpleParameters returns whether the given function node
// has any non-simple parameters (e.g. default values, rest parameters, or
// destructuring patterns).
//
// fnNode is a FUNCTION node. Returns true if function parameters are
// non-simple; false otherwise.
//
// Java: NodeUtil.hasNonSimpleParameters(Node)
func NodeUtil_HasNonSimpleParameters(fnNode rhino.Node) bool {
	jutil.CheckArgument(fnNode.IsFunction())
	paramList := NodeUtil_GetFunctionParameters(fnNode)
	for param := paramList.GetFirstChild(); !param.IsNil(); param = param.GetNext() {
		if !param.IsName() {
			return true
		}
	}
	return false
}

// NodeUtil_FunctionParametersMayHaveSideEffects returns whether any parameter
// of the given function node may have side effects when evaluated, or may
// throw an exception during parameter pattern matching / destructuring.
//
// fnNode is a FUNCTION node; astAnalyzer is the AstAnalyzer to check side
// effects of default parameter initializers. Returns true if function
// parameters may have side effects or throw; false otherwise.
//
// Java: NodeUtil.functionParametersMayHaveSideEffects(Node, AstAnalyzer)
func NodeUtil_FunctionParametersMayHaveSideEffects(fnNode rhino.Node, astAnalyzer *AstAnalyzer) bool {
	jutil.CheckArgument(fnNode.IsFunction())
	paramList := NodeUtil_GetFunctionParameters(fnNode)
	for param := paramList.GetFirstChild(); !param.IsNil(); param = param.GetNext() {
		if param.IsName() {
			continue
		}
		if param.IsRest() {
			if param.GetFirstChild().IsName() {
				continue
			}
			return true
		}
		if param.IsDefaultValue() {
			target := param.GetFirstChild()
			initializer := param.GetSecondChild()
			if !target.IsName() {
				return true
			}
			if astAnalyzer.MayHaveSideEffects(initializer) {
				return true
			}
			continue
		}
		return true
	}
	return false
}

// Java: NodeUtil.isConstantVar(Node, Scope)
func NodeUtil_IsConstantVar(node rhino.Node, scope *Scope) bool {
	if NodeUtil_IsConstantName(node) {
		return true
	}

	if !node.IsName() || scope == nil {
		return false
	}

	var_ := scope.GetVar(node.GetString())
	return var_ != nil && (var_.IsDeclaredOrInferredConst() || var_.IsConst())
}

// NodeUtil_IsConstantName determines whether a variable is constant:
//
//  1. In Normalize, any name that matches the CodingConvention#isConstant(String)
//     is annotated with an IS_CONSTANT_NAME property.
//
// node is a NAME or STRING node. Returns true if a name node represents a
// constant variable.
//
// TODO(dimvar): this method and the next two do similar but not quite
// identical things. Clean up
//
// Java: NodeUtil.isConstantName(Node)
func NodeUtil_IsConstantName(node rhino.Node) bool {
	return node.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME)
}

// NodeUtil_IsConstantByConvention returns whether the given name is constant
// by coding convention.
//
// Deprecated: we want to delete the constant by convention logic - see
// http://b/135755127
//
// Java: NodeUtil.isConstantByConvention(CodingConvention, Node)
func NodeUtil_IsConstantByConvention(convention CodingConvention, node rhino.Node) bool {
	if NodeUtil_IsNormalOrOptChainGetProp(node) {
		return convention.IsConstantKey(node.GetString())
	} else if NodeUtil_MayBeObjectLitKey(node) {
		return convention.IsConstantKey(node.GetString())
	} else if node.IsName() {
		return convention.IsConstant(node.GetString())
	}
	return false
}

// NodeUtil_IsConstantDeclaration determines whether the given lvalue is
// declared constant or is a name assigned exactly once.
//
// Note that this intentionally excludes variables that are constant
// according to the coding convention.
//
// node is some lvalue node. Throws IllegalStateException if the given node
// is not an lvalue.
//
// Java: NodeUtil.isConstantDeclaration(JSDocInfo, Node)
func NodeUtil_IsConstantDeclaration(info *rhino.JSDocInfo, node rhino.Node) bool {
	if NodeUtil_IsObjectLitKey(node) ||
		(node.GetParent().IsAssign() && node.IsFirstChildOf(node.GetParent())) ||
		(node.GetParent().IsExprResult() && NodeUtil_IsNormalGet(node)) ||
		node.IsMemberFieldDef() ||
		node.IsComputedFieldDef() ||
		node.IsComputedProp() {
		return info != nil && info.IsConstant()
	}
	jutil.CheckArgumentMsg1(node.IsName(), "%s", node)
	declaringParent := NodeUtil_GetDeclaringParent(node) // throws an error if `node` is not an lvalue
	if declaringParent.IsConst() {
		return true
	} else if info != nil && info.IsConstant() {
		return true
	}
	return node.IsInferredConstantVar()
}

// Java: NodeUtil.functionHasInlineJsdocs(Node)
func NodeUtil_FunctionHasInlineJsdocs(fn rhino.Node) bool {
	if !fn.IsFunction() {
		return false
	}
	// Check inline return annotation
	if fn.GetFirstChild().GetJSDocInfo() != nil {
		return true
	}
	// Check inline parameter annotations
	param := fn.GetSecondChild().GetFirstChild()
	for !param.IsNil() {
		if param.GetJSDocInfo() != nil {
			return true
		}
		param = param.GetNext()
	}
	return false
}

// NodeUtil_GetSourceName returns the source name property on the node or its
// ancestors.
//
// Java: NodeUtil.getSourceName(Node)
func NodeUtil_GetSourceName(n rhino.Node) *string {
	var sourceName *string
	for sourceName == nil && !n.IsNil() {
		sourceName = n.GetSourceFileName()
		n = n.GetParent()
	}
	return sourceName
}

// NodeUtil_GetSourceFile returns the source name property on the node or its
// ancestors.
//
// Java: NodeUtil.getSourceFile(Node)
func NodeUtil_GetSourceFile(n rhino.Node) rhino.StaticSourceFile {
	var sourceName rhino.StaticSourceFile
	for sourceName == nil && !n.IsNil() {
		sourceName = n.GetStaticSourceFile()
		n = n.GetParent()
	}
	return sourceName
}

// NodeUtil_GetInputId returns the InputId property on the node or its
// ancestors.
//
// Java: NodeUtil.getInputId(Node)
func NodeUtil_GetInputId(n rhino.Node) *rhino.InputId {
	for !n.IsNil() && !n.IsScript() {
		n = n.GetParent()
	}

	if !n.IsNil() && n.IsScript() {
		return n.GetInputId()
	}
	return nil
}

// NodeUtil_NewCallNode returns a new CALL node with the "FREE_CALL" set
// based on call target.
//
// Java: NodeUtil.newCallNode(Node, Node...)
func NodeUtil_NewCallNode(callTarget rhino.Node, parameters ...rhino.Node) rhino.Node {
	isFreeCall := !NodeUtil_IsNormalGet(callTarget)
	call := rhino.IR_Call(callTarget)
	call.PutBooleanProp(rhino.Node_FREE_CALL, isFreeCall)
	for _, parameter := range parameters {
		call.AddChildToBack(parameter)
	}
	return call
}

// NodeUtil_EvaluatesToLocalValue: whether the result of the expression node
// is known to be a primitive value or an object that has not yet escaped.
//
// This guarantee is different than that provided by isLiteralValue (where
// literal values are immune to side-effects if unescaped) or
// isImmutableValue (which can be safely aliased).
//
// The concept of "local values" allow for the containment of side-effect
// operations. For example, setting a property on a local value does not
// produce a global side-effect.
//
// Note that the concept of "local value" is not deep, it does not say
// anything about the properties of the "local value" (all class instances
// have "constructor" properties that are not local values for instance).
//
// Note that this method only provides the starting state of the expression
// result, it does not guarantee that the value is forever a local value. If
// the containing method has any non-local side-effect, "local values" may
// escape.
//
// Java: NodeUtil.evaluatesToLocalValue(Node)
func NodeUtil_EvaluatesToLocalValue(value rhino.Node) bool {
	switch value.GetToken() {
	case rhino.Token_ASSIGN:
		// A result that is aliased by a non-local name, is the effectively the
		// same as returning a non-local name, but this doesn't matter if the
		// value is immutable.
		return NodeUtil_IsImmutableValue(value.GetLastChild())
	case rhino.Token_COMMA:
		return NodeUtil_EvaluatesToLocalValue(value.GetLastChild())
	case rhino.Token_AND, rhino.Token_OR, rhino.Token_COALESCE:
		return NodeUtil_EvaluatesToLocalValue(value.GetFirstChild()) &&
			NodeUtil_EvaluatesToLocalValue(value.GetLastChild())
	case rhino.Token_HOOK:
		return NodeUtil_EvaluatesToLocalValue(value.GetSecondChild()) &&
			NodeUtil_EvaluatesToLocalValue(value.GetLastChild())
	case rhino.Token_DYNAMIC_IMPORT:
		// Dynamic import always returns a newly created Promise.
		return true
	case rhino.Token_THIS, rhino.Token_SUPER:
		return false
	case rhino.Token_NAME:
		return NodeUtil_IsImmutableValue(value)
	case rhino.Token_GETELEM, rhino.Token_GETPROP, rhino.Token_OPTCHAIN_GETELEM, rhino.Token_OPTCHAIN_GETPROP:
		// There is no information about the locality of object properties.
		return false
	case rhino.Token_CALL, rhino.Token_OPTCHAIN_CALL:
		return NodeUtil_IsToStringMethodCall(value)
	case rhino.Token_TAGGED_TEMPLATELIT:
		// No information about local values for tagged template literals
		return false
	case rhino.Token_NEW:
		return NodeUtil_NewHasLocalResult(value)
	case rhino.Token_DELPROP, rhino.Token_INC, rhino.Token_DEC, rhino.Token_CLASS, rhino.Token_FUNCTION,
		rhino.Token_REGEXP, rhino.Token_EMPTY, rhino.Token_ARRAYLIT, rhino.Token_OBJECTLIT, rhino.Token_TEMPLATELIT:
		return true
	case rhino.Token_CAST:
		return NodeUtil_EvaluatesToLocalValue(value.GetFirstChild())
	case rhino.Token_ITER_SPREAD, rhino.Token_OBJECT_SPREAD, rhino.Token_NEW_TARGET:
		// TODO(johnlenz): remove this case.
		// Returns an alias of a constructor (current or subclass).
		return false
	case rhino.Token_YIELD, rhino.Token_AWAIT:
		// TODO(johnlenz): we can do better for await if we use type information.  That is,
		// if we know the promise being awaited on is a immutable value type (string, etc)
		// we could return true here.
		return false
	default:
		// A logical assignment could evaluate to either the assignment target or the
		// right-hand side, and we generally have no information about the locality of the lhs.
		if NodeUtil_IsLogicalAssignmentOp(value) {
			return false
		}
		// Other non-logical assignment ops force a local value:
		//  '' + g (a local string)
		//  x -= g (x is now an local number)
		if NodeUtil_IsAssignmentOp(value) || NodeUtil_IsSimpleOperator(value) || NodeUtil_IsImmutableValue(value) {
			return true
		}

		panic(jutil.NewIllegalStateExceptionMsg(
			"Unexpected expression node: " + value.String() + "\n parent:" + jutil.ToJavaString(value.GetParent())))
	}
}

// NodeUtil_MayBeUndefined returns whether the provided expression is may
// evaluate to 'undefined'.
//
// Java: NodeUtil.mayBeUndefined(Node)
func NodeUtil_MayBeUndefined(n rhino.Node) bool {
	return !NodeUtil_IsDefinedValue(n)
}

// NodeUtil_IsDefinedValue returns whether the provided expression is known
// not to evaluate to 'undefined'.
//
// Similar to #getKnownValueType only for 'undefined'. This is useful for
// simplifying default value expressions.
//
// Java: NodeUtil.isDefinedValue(Node)
func NodeUtil_IsDefinedValue(value rhino.Node) bool {
	switch value.GetToken() {
	case rhino.Token_ASSIGN, rhino.Token_CAST, rhino.Token_COMMA:
		// Only the assigned value matters here.
		return NodeUtil_IsDefinedValue(value.GetLastChild())
	case rhino.Token_COALESCE:
		// 'null' is a "defined" value so we can only trust the RHS.
		// NOTE: consider creating and using a "isDefinedAndNotNull" that would allow us to
		// trust the tested value.
		return NodeUtil_IsDefinedValue(value.GetSecondChild())
	case rhino.Token_AND, rhino.Token_OR:
		return NodeUtil_IsDefinedValue(value.GetFirstChild()) && NodeUtil_IsDefinedValue(value.GetLastChild())
	case rhino.Token_HOOK:
		return NodeUtil_IsDefinedValue(value.GetSecondChild()) && NodeUtil_IsDefinedValue(value.GetLastChild())
		// Assume undefined leaks in this and call results.
	case rhino.Token_CALL,
		rhino.Token_OPTCHAIN_CALL,
		rhino.Token_GETELEM,
		rhino.Token_GETPROP,
		rhino.Token_OPTCHAIN_GETELEM,
		rhino.Token_OPTCHAIN_GETPROP,
		rhino.Token_TAGGED_TEMPLATELIT,
		rhino.Token_THIS,
		rhino.Token_YIELD,
		rhino.Token_AWAIT,
		rhino.Token_VOID:
		return false
	case rhino.Token_DELPROP,
		rhino.Token_INC,
		rhino.Token_DEC,
		rhino.Token_CLASS,
		rhino.Token_FUNCTION,
		rhino.Token_REGEXP,
		rhino.Token_EMPTY,
		rhino.Token_ARRAYLIT,
		rhino.Token_OBJECTLIT,
		rhino.Token_TEMPLATELIT,
		rhino.Token_STRINGLIT,
		rhino.Token_NUMBER,
		rhino.Token_BIGINT,
		rhino.Token_NULL,
		rhino.Token_TRUE,
		rhino.Token_FALSE,
		rhino.Token_NEW:
		return true
	case rhino.Token_TEMPLATELIT_STRING:
		return value.GetCookedString() != nil
	case rhino.Token_NAME:
		name := value.GetString()
		// We assume here that programs don't change the value of the keyword
		// undefined to something other than the value undefined.
		return "Infinity" == name || "NaN" == name
	default:
		// Other op force a local value:
		//  '' + g (a  string)
		//  x -= g (x is now an number)
		if NodeUtil_IsAssignmentOp(value) || NodeUtil_IsSimpleOperator(value) {
			return true
		}

		panic(jutil.NewIllegalStateExceptionMsg(
			"Unexpected expression node: " + value.String() + "\n parent:" + jutil.ToJavaString(value.GetParent())))
	}
}

// NodeUtil_GetNthSibling: given the first sibling, this returns the nth
// sibling or null if no such sibling exists. This is like "getChildAtIndex"
// but returns null for non-existent indexes.
//
// Java: NodeUtil.getNthSibling(Node, int)
func NodeUtil_GetNthSibling(first rhino.Node, index int) rhino.Node {
	sibling := first
	for index != 0 && !sibling.IsNil() {
		sibling = sibling.GetNext()
		index--
	}
	return sibling
}

// NodeUtil_GetArgumentForFunction: given the function, this returns the nth
// argument or null if no such parameter exists.
//
// Java: NodeUtil.getArgumentForFunction(Node, int)
func NodeUtil_GetArgumentForFunction(function rhino.Node, index int) rhino.Node {
	jutil.CheckState(function.IsFunction())
	return NodeUtil_GetNthSibling(function.GetSecondChild().GetFirstChild(), index)
}

// NodeUtil_GetArgumentForCallOrNew: given the new or call, this returns the
// nth argument of the call or null if no such argument exists.
//
// Java: NodeUtil.getArgumentForCallOrNew(Node, int)
func NodeUtil_GetArgumentForCallOrNew(call rhino.Node, index int) rhino.Node {
	jutil.CheckState(NodeUtil_IsCallOrNew(call))
	return NodeUtil_GetNthSibling(call.GetSecondChild(), index)
}

// NodeUtil_IsInvocationTarget returns whether this is a target of a call or
// new.
//
// Java: NodeUtil.isInvocationTarget(Node)
func NodeUtil_IsInvocationTarget(n rhino.Node) bool {
	parent := n.GetParent()
	return !parent.IsNil() &&
		(NodeUtil_IsCallOrNew(parent) || parent.IsTaggedTemplateLit()) &&
		parent.GetFirstChild() == n
}

// NodeUtil_IsInvocation returns whether this is a call (including tagged
// template lits) or new.
//
// Java: NodeUtil.isInvocation(Node)
func NodeUtil_IsInvocation(n rhino.Node) bool {
	return NodeUtil_IsCallOrNew(n) || n.IsTaggedTemplateLit()
}

// Java: NodeUtil.isCallOrNewArgument(Node)
func NodeUtil_IsCallOrNewArgument(n rhino.Node) bool {
	parent := n.GetParent()
	return !parent.IsNil() && NodeUtil_IsCallOrNew(parent) && parent.GetFirstChild() != n
}

// Java: NodeUtil.isToStringMethodCall(Node)
func NodeUtil_IsToStringMethodCall(call rhino.Node) bool {
	getNode := call.GetFirstChild()
	var name string
	if NodeUtil_IsNormalOrOptChainGetProp(getNode) {
		name = getNode.GetString()
	} else if NodeUtil_IsNormalOrOptChainGet(getNode) {
		propNode := getNode.GetLastChild()
		if !propNode.IsStringLit() {
			return false
		}
		name = propNode.GetString()
	} else {
		return false
	}

	return "toString" == name
}

// NodeUtil_GetDeclaredTypeExpression returns declared JSDoc type for the
// given name declaration, or null if none present.
//
// Java: NodeUtil.getDeclaredTypeExpression(Node)
func NodeUtil_GetDeclaredTypeExpression(declaration rhino.Node) *rhino.JSTypeExpression {
	jutil.CheckArgument(declaration.IsName() || declaration.IsStringKey())
	nameJsdoc := NodeUtil_GetBestJSDocInfo(declaration)
	if nameJsdoc != nil {
		return nameJsdoc.GetType()
	}
	parent := declaration.GetParent()
	if parent.IsRest() || parent.IsDefaultValue() {
		parent = parent.GetParent()
	}
	if parent.IsParamList() {
		functionJsdoc := NodeUtil_GetBestJSDocInfo(parent.GetParent())
		if functionJsdoc != nil {
			return functionJsdoc.GetParameterType(declaration.GetString())
		}
	}
	return nil
}

// NodeUtil_GetBestJSDocInfo finds the best JSDoc for the given node.
//
// Java: NodeUtil.getBestJSDocInfo(Node)
func NodeUtil_GetBestJSDocInfo(n rhino.Node) *rhino.JSDocInfo {
	jsdocNode := NodeUtil_GetBestJSDocInfoNode(n)
	if jsdocNode.IsNil() {
		return nil
	}
	return jsdocNode.GetJSDocInfo()
}

// NodeUtil_GetBestJsDocInfoNodeStrict finds the best JSDocInfo node for the
// given node. This version is stricter than getBestJSDocInfoNode in that it
// only accepts a node that is either a lhs "let, const, var" name node, or a
// RHS class/function, or a class/function declaration, or an export.
//
// Throws IllegalStateException if any other input node token is provided.
//
// Java: NodeUtil.getBestJsDocInfoNodeStrict(Node)
func NodeUtil_GetBestJsDocInfoNodeStrict(n rhino.Node) rhino.Node {
	isDeclaredName := n.IsName() && NodeUtil_IsNameDeclaration(n.GetParent())
	if isDeclaredName {
		return NodeUtil_GetBestJsDocInfoNodeInternal(n)
	}
	isDeclaredClass := NodeUtil_IsClassDeclaration(n)
	if isDeclaredClass {
		return NodeUtil_GetBestJsDocInfoNodeInternal(n)
	}
	isDeclaredFunction := NodeUtil_IsFunctionDeclaration(n)
	if isDeclaredFunction {
		return NodeUtil_GetBestJsDocInfoNodeInternal(n)
	}
	isRhsClass := NodeUtil_IsRhsClass(n)
	if isRhsClass {
		return NodeUtil_GetBestJsDocInfoNodeInternal(n)
	}
	isRhsFunction := NodeUtil_IsRhsFunction(n)
	if isRhsFunction {
		return NodeUtil_GetBestJsDocInfoNodeInternal(n)
	}
	isRhsClassName := n.IsName() && NodeUtil_IsRhsClass(n.GetParent())
	if isRhsClassName {
		return NodeUtil_GetBestJsDocInfoNodeInternal(n)
	}
	isRhsFunctionName := n.IsName() && NodeUtil_IsRhsFunction(n.GetParent())
	if isRhsFunctionName {
		return NodeUtil_GetBestJsDocInfoNodeInternal(n)
	}
	isExport := n.IsExport()
	if isExport {
		return NodeUtil_GetBestJsDocInfoNodeInternal(n)
	}
	panic(jutil.NewIllegalStateExceptionMsg("Not allowed to get JSDocInfo node for node: " + n.String()))
}

// NodeUtil_IsRhsClass: e.g. `class A` in `/** some */ var x = class A {};`
//
// Java: NodeUtil.isRhsClass(Node)
func NodeUtil_IsRhsClass(n rhino.Node) bool {
	return NodeUtil_IsNamedClassExpression(n) &&
		(n.GetParent().IsName() || n.GetParent().IsAssign())
}

// NodeUtil_IsRhsFunction: e.g. `function A` in
// `/** some */ var x = function A() {};`
//
// Java: NodeUtil.isRhsFunction(Node)
func NodeUtil_IsRhsFunction(n rhino.Node) bool {
	return NodeUtil_IsNamedFunctionExpression(n) &&
		(n.GetParent().IsName() || n.GetParent().IsAssign())
}

// NodeUtil_GetBestJSDocInfoNode finds the best JSDocInfo node for the given
// node.
//
// Deprecated: because we want to control the input node tokens accepted by
// this function. Use getBestJSDocInfoNodeStrict instead.
//
// Java: NodeUtil.getBestJSDocInfoNode(Node)
func NodeUtil_GetBestJSDocInfoNode(n rhino.Node) rhino.Node {
	return NodeUtil_GetBestJsDocInfoNodeInternal(n)
}

// Java: NodeUtil.getBestJsDocInfoNodeInternal(Node)
func NodeUtil_GetBestJsDocInfoNodeInternal(n rhino.Node) rhino.Node {
	if n.IsExprResult() {
		return NodeUtil_GetBestJSDocInfoNode(n.GetFirstChild())
	}
	info := n.GetJSDocInfo()
	if info == nil {
		parent := n.GetParent()
		if parent.IsNil() || n.IsExprResult() {
			return rhino.NilNode
		}

		if parent.IsName() ||
			parent.IsExport() ||
			parent.IsVar() ||
			parent.IsLet() ||
			parent.IsConst() ||
			parent.IsDeclare() {
			return NodeUtil_GetBestJSDocInfoNode(parent)
		} else if parent.IsAssign() {
			return NodeUtil_GetBestJSDocInfoNode(parent)
		} else if NodeUtil_MayBeObjectLitKey(parent) || parent.IsComputedProp() {
			return parent
		} else if (parent.IsFunction() || parent.IsClass()) && n == parent.GetFirstChild() {
			// n is the NAME node of the function/class.
			return NodeUtil_GetBestJSDocInfoNode(parent)
		} else if NodeUtil_IsNameDeclaration(parent) && parent.HasOneChild() {
			return parent
		} else if (parent.IsHook() && parent.GetFirstChild() != n) ||
			parent.IsOr() ||
			parent.IsAnd() ||
			(parent.IsComma() && parent.GetFirstChild() != n) {
			return NodeUtil_GetBestJSDocInfoNode(parent)
		}
	}
	return n
}

// NodeUtil_GetBestLValue finds the l-value that the given r-value is being
// assigned to.
//
// Java: NodeUtil.getBestLValue(Node)
func NodeUtil_GetBestLValue(n rhino.Node) rhino.Node {
	parent := n.GetParent()
	if NodeUtil_IsFunctionDeclaration(n) || NodeUtil_IsClassDeclaration(n) {
		return n.GetFirstChild()
	} else if n.IsClassMembers() {
		return NodeUtil_GetBestLValue(parent)
	} else if parent.IsName() {
		return parent
	} else if parent.IsAssign() {
		return parent.GetFirstChild()
	} else if NodeUtil_MayBeObjectLitKey(parent) || parent.IsComputedProp() {
		return parent
	} else if (parent.IsHook() && parent.GetFirstChild() != n) ||
		parent.IsOr() ||
		parent.IsAnd() ||
		(parent.IsComma() && parent.GetFirstChild() != n) {
		return NodeUtil_GetBestLValue(parent)
	} else if parent.IsCast() {
		return NodeUtil_GetBestLValue(parent)
	}
	return rhino.NilNode
}

// NodeUtil_GetRValueOfLValue gets the r-value (or initializer) of a node
// returned by getBestLValue.
//
// Java: NodeUtil.getRValueOfLValue(Node)
func NodeUtil_GetRValueOfLValue(n rhino.Node) rhino.Node {
	parent := n.GetParent()
	switch parent.GetToken() {
	case rhino.Token_ASSIGN,
		rhino.Token_ASSIGN_BITOR,
		rhino.Token_ASSIGN_BITXOR,
		rhino.Token_ASSIGN_BITAND,
		rhino.Token_ASSIGN_LSH,
		rhino.Token_ASSIGN_RSH,
		rhino.Token_ASSIGN_URSH,
		rhino.Token_ASSIGN_ADD,
		rhino.Token_ASSIGN_SUB,
		rhino.Token_ASSIGN_MUL,
		rhino.Token_ASSIGN_EXPONENT,
		rhino.Token_ASSIGN_DIV,
		rhino.Token_ASSIGN_MOD,
		rhino.Token_ASSIGN_OR,
		rhino.Token_ASSIGN_AND,
		rhino.Token_ASSIGN_COALESCE,
		rhino.Token_DESTRUCTURING_LHS:
		return n.GetNext()
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		return n.GetLastChild()
	case rhino.Token_OBJECTLIT:
		jutil.CheckStateMsg1(
			n.IsStringKey() ||
				n.IsComputedProp() ||
				n.IsMemberFunctionDef() ||
				n.IsGetterDef() ||
				n.IsSetterDef(),
			"%s",
			n)
		return n.GetLastChild()
	case rhino.Token_CLASS_MEMBERS:
		jutil.CheckStateMsg1(
			n.IsMemberFunctionDef() ||
				n.IsMemberFieldDef() ||
				n.IsComputedFieldDef() ||
				n.IsGetterDef() ||
				n.IsSetterDef(),
			"%s",
			n)
		return n.GetLastChild()
	case rhino.Token_FUNCTION, rhino.Token_CLASS:
		return parent
	default:
	}
	return rhino.NilNode
}

// NodeUtil_GetBestLValueOwner gets the owner of the given l-value node.
//
// Java: NodeUtil.getBestLValueOwner(Node)
//
// TODO(b/189993301): do I have to add logic for class fields in these LValue
// functions?
func NodeUtil_GetBestLValueOwner(lValue rhino.Node) rhino.Node {
	if lValue.IsNil() || lValue.GetParent().IsNil() {
		return rhino.NilNode
	}
	if NodeUtil_MayBeObjectLitKey(lValue) || lValue.IsComputedProp() {
		return NodeUtil_GetBestLValue(lValue.GetParent())
	} else if NodeUtil_IsNormalGet(lValue) {
		return lValue.GetFirstChild()
	}

	return rhino.NilNode
}

// NodeUtil_GetBestLValueName gets the name of the given l-value node.
//
// Go-only: performance — the body is in nodeUtil_getBestLValueNameString, so
// that this wrapper inlines and a caller that only reads the result keeps
// the string header on its stack.
//
// Java: NodeUtil.getBestLValueName(Node)
func NodeUtil_GetBestLValueName(lValue rhino.Node) *string {
	s, ok := nodeUtil_getBestLValueNameString(lValue)
	if !ok {
		return nil
	}
	return &s
}

// nodeUtil_getBestLValueNameString is the body of getBestLValueName; ok is
// false where Java returns null.
func nodeUtil_getBestLValueNameString(lValue rhino.Node) (string, bool) {
	if lValue.IsNil() || lValue.GetParent().IsNil() {
		return "", false
	}
	if lValue.GetParent().IsClassMembers() && !lValue.IsComputedProp() {
		className := NodeUtil_GetName(lValue.GetGrandparent())
		if className == nil { // Anonymous class
			return "", false
		}
		methodName := lValue.GetString()
		if lValue.IsStaticMember() {
			return *className + "." + methodName, true
		}
		return *className + ".prototype." + methodName, true
	}
	// TODO(sdh): Tighten this to simply require !lValue.isQuotedString()
	// Could get rid of the isJSIdentifier check, but may need to fix depot.
	if NodeUtil_MayBeObjectLitKey(lValue) {
		owner := NodeUtil_GetBestLValue(lValue.GetParent())
		if !owner.IsNil() {
			key := NodeUtil_GetObjectOrClassLitKeyName(lValue)
			if !rhino.TokenStream_IsJSIdentifier(key) {
				return "", false
			}
			ownerName, ok := nodeUtil_getBestLValueNameString(owner)
			if ok {
				return ownerName + "." + key, true
			}
		}
		return "", false
	}
	return lValue.GetQualifiedNameString()
}

// NodeUtil_GetBestLValueRoot gets the root of a qualified name l-value.
//
// Java: NodeUtil.getBestLValueRoot(Node)
func NodeUtil_GetBestLValueRoot(lValue rhino.Node) rhino.Node {
	if lValue.IsNil() {
		return rhino.NilNode
	}
	switch lValue.GetToken() {
	case rhino.Token_STRING_KEY:
		// NOTE: beware of getBestLValue returning null (or be null-permissive?)
		return NodeUtil_GetBestLValueRoot(NodeUtil_GetBestLValue(lValue.GetParent()))
	case rhino.Token_GETPROP, rhino.Token_GETELEM:
		return NodeUtil_GetBestLValueRoot(lValue.GetFirstChild())
	case rhino.Token_THIS, rhino.Token_SUPER, rhino.Token_NAME:
		return lValue
	default:
		return rhino.NilNode
	}
}

// NodeUtil_IsExpressionResultUsed returns true iff the result of the
// expression is consumed.
//
// Java: NodeUtil.isExpressionResultUsed(Node)
func NodeUtil_IsExpressionResultUsed(expr rhino.Node) bool {
	parent := expr.GetParent()
	switch parent.GetToken() {
	case rhino.Token_BLOCK, rhino.Token_EXPR_RESULT:
		return false
	case rhino.Token_CAST:
		return NodeUtil_IsExpressionResultUsed(parent)
	case rhino.Token_HOOK, rhino.Token_AND, rhino.Token_OR, rhino.Token_COALESCE:
		return (expr == parent.GetFirstChild()) || NodeUtil_IsExpressionResultUsed(parent)
	case rhino.Token_COMMA:
		grandparent := parent.GetParent()
		if (grandparent.IsCall() || grandparent.IsTaggedTemplateLit()) &&
			parent == grandparent.GetFirstChild() &&
			expr == parent.GetFirstChild() &&
			parent.HasTwoChildren() {
			// Special case the indirect function call pattern, e.g. (0, myFn)(arg1, arg2).
			// The indirect call pattern has two use cases:
			calledFn := parent.GetLastChild()
			if calledFn.IsName() && calledFn.MatchesName("eval") {
				// 1) eval
				//    Semantically, a direct call to eval is different from an indirect
				//    call to an eval. See ECMA-262 S15.1.2.1.
				return true
			}
			if NodeUtil_IsNormalOrOptChainGet(calledFn) {
				// 2) Removing "this" context
				//    When calling prefix.myFn(), myFn receives "prefix" as its this context.
				//    Calling "(0, prefix.myFn)()" removes the this context, so this is useful code.
				return true
			}
		}
		return expr != parent.GetFirstChild() && NodeUtil_IsExpressionResultUsed(parent)
	case rhino.Token_FOR:
		// Only an expression whose result is in the condition part of the
		// expression is used.
		return parent.GetSecondChild() == expr
	default:
	}
	return true
}

// NodeUtil_IsExecutedExactlyOnce: n is the expression to check. Returns
// whether the expression is unconditionally executed only once in the
// containing execution scope.
//
// Java: NodeUtil.isExecutedExactlyOnce(Node)
func NodeUtil_IsExecutedExactlyOnce(n rhino.Node) bool {
	// Go-only: language-adapter — Java's do-while; `continue inspect` jumps
	// to the loop condition, which ends each iteration below.
inspect:
	for {
		parent := n.GetParent()
		switch parent.GetToken() {
		case rhino.Token_IF, rhino.Token_HOOK, rhino.Token_AND, rhino.Token_OR, rhino.Token_COALESCE,
			rhino.Token_ASSIGN_OR, rhino.Token_ASSIGN_AND, rhino.Token_ASSIGN_COALESCE:
			if parent.GetFirstChild() != n {
				return false
			}
			// other ancestors may be conditional
			// continue inspect
		case rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
			if NodeUtil_IsEnhancedFor(parent) {
				if parent.GetSecondChild() != n {
					return false
				}
			} else {
				if parent.GetFirstChild() != n {
					return false
				}
			}
			// other ancestors may be conditional
			// continue inspect
		case rhino.Token_TRY:
			// Consider all code under a try/catch to be conditionally executed.
			if !NodeUtil_HasFinally(parent) || parent.GetLastChild() != n {
				return false
			}
			// continue inspect
		case rhino.Token_WHILE, rhino.Token_DO, rhino.Token_CASE, rhino.Token_DEFAULT_CASE:
			return false
		case rhino.Token_SCRIPT, rhino.Token_FUNCTION:
			// Done, we've reached the scope root.
			break inspect
		default:
		}
		n = n.GetParent()
		if n.IsNil() {
			break
		}
	}
	return true
}

// NodeUtil_BooleanNode returns an appropriate AST node for the boolean
// value.
//
// Java: NodeUtil.booleanNode(boolean)
func NodeUtil_BooleanNode(value bool) rhino.Node {
	if value {
		return rhino.IR_TrueNode()
	}
	return rhino.IR_FalseNode()
}

// NodeUtil_NumberNode returns an appropriate AST node for the double value.
//
// Java: NodeUtil.numberNode(double, Node)
func NodeUtil_NumberNode(value float64, srcref rhino.Node) rhino.Node {
	var result rhino.Node
	if math.IsNaN(value) {
		result = rhino.IR_Name("NaN")
		result.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
	} else {
		if math.IsInf(value, 0) {
			result = rhino.IR_Name("Infinity")
			result.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
		} else {
			result = rhino.IR_Number(math.Abs(value))
		}
		if base.JSCompDoubles_IsNegative(value) {
			result = rhino.IR_Neg(result)
		}
	}
	if !srcref.IsNil() {
		result.SrcrefTree(srcref)
	}
	return result
}

// Java: NodeUtil.isNaN(Node)
func NodeUtil_IsNaN(n rhino.Node) bool {
	return (n.IsName() && n.GetString() == "NaN") ||
		(n.GetToken() == rhino.Token_DIV &&
			n.GetFirstChild().IsNumber() &&
			n.GetFirstChild().GetDouble() == 0 &&
			n.GetLastChild().IsNumber() &&
			n.GetLastChild().GetDouble() == 0) ||
		n.MatchesQualifiedNameNode(NodeUtil_NUMBER_NAN)
}

// Java: NodeUtil.NUMBER_NAN
var NodeUtil_NUMBER_NAN = rhino.IR_Getprop(rhino.IR_Name("Number"), "NaN")

// Java: NodeUtil.countAstSizeUpToLimit(Node, int)
func NodeUtil_CountAstSizeUpToLimit(n rhino.Node, limit int) int {
	count := 1
	if limit > 1 {
		for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			if c.GetFirstChild().IsNil() || limit-count <= 1 {
				count++
			} else {
				count += NodeUtil_CountAstSizeUpToLimit(c, limit-count)
			}
		}
	}
	return count
}

// Java: NodeUtil.countAstSize(Node)
func NodeUtil_CountAstSize(n rhino.Node) int {
	count := 1
	if n.GetFirstChild().IsNil() {
		return 1
	}
	var stackBuf [64]rhino.Node
	stack := stackBuf[:0]
	cur := n
	for {
		for c := cur.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			count++
			if gc := c.GetFirstChild(); !gc.IsNil() {
				if gc.GetNext().IsNil() && gc.GetFirstChild().IsNil() {
					count++
				} else {
					stack = append(stack, c)
				}
			}
		}
		if len(stack) == 0 {
			return count
		}
		cur = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
	}
}

// Java: NodeUtil.createConstantJsDoc()
func NodeUtil_CreateConstantJsDoc() *rhino.JSDocInfo {
	builder := rhino.JSDocInfo_Builder_()
	builder.RecordConstancy()
	return builder.Build()
}

// Java: NodeUtil.isGoogProvideCall(Node)
func NodeUtil_IsGoogProvideCall(n rhino.Node) bool {
	if NodeUtil_IsExprCall(n) {
		target := n.GetFirstFirstChild()
		return NodeUtil_GOOG_PROVIDE.Matches(target)
	}
	return false
}

// Java: NodeUtil.isGoogModuleCall(Node)
func NodeUtil_IsGoogModuleCall(n rhino.Node) bool {
	if NodeUtil_IsExprCall(n) {
		target := n.GetFirstFirstChild()
		return NodeUtil_GOOG_MODULE.Matches(target)
	}
	return false
}

// Java: NodeUtil.isGoogModuleGetCall(Node)
func NodeUtil_IsGoogModuleGetCall(callNode rhino.Node) bool {
	if !callNode.IsCall() {
		return false
	}
	return NodeUtil_GOOG_MODULE_GET.Matches(callNode.GetFirstChild()) &&
		callNode.HasTwoChildren() &&
		callNode.GetSecondChild().IsStringLit()
}

// Java: NodeUtil.isGoogRequireCall(Node)
func NodeUtil_IsGoogRequireCall(call rhino.Node) bool {
	if call.IsCall() {
		target := call.GetFirstChild()
		return NodeUtil_GOOG_REQUIRE.Matches(target)
	}
	return false
}

// Java: NodeUtil.isGoogRequireTypeCall(Node)
func NodeUtil_IsGoogRequireTypeCall(call rhino.Node) bool {
	if call.IsCall() {
		target := call.GetFirstChild()
		return NodeUtil_GOOG_REQUIRE_TYPE.Matches(target)
	}
	return false
}

// Java: NodeUtil.isGoogForwardDeclareCall(Node)
func NodeUtil_IsGoogForwardDeclareCall(call rhino.Node) bool {
	if call.IsCall() {
		target := call.GetFirstChild()
		return NodeUtil_GOOG_FORWARD_DECLARE.Matches(target)
	}
	return false
}

// Java: NodeUtil.isGoogRequireDynamicCall(Node)
func NodeUtil_IsGoogRequireDynamicCall(call rhino.Node) bool {
	if call.IsCall() {
		target := call.GetFirstChild()
		return NodeUtil_GOOG_REQUIRE_DYNAMIC.Matches(target)
	}
	return false
}

// Java: NodeUtil.isGoogWeakUsageCall(Node)
func NodeUtil_IsGoogWeakUsageCall(call rhino.Node) bool {
	if call.IsCall() {
		target := call.GetFirstChild()
		return NodeUtil_GOOG_WEAK_USAGE.Matches(target) || NodeUtil_GOOG_WEAK_USAGE_MANGLED.Matches(target)
	}
	return false
}

// Java: NodeUtil.isModuleScopeRoot(Node)
func NodeUtil_IsModuleScopeRoot(n rhino.Node) bool {
	return n.IsModuleBody() || NodeUtil_IsBundledGoogModuleScopeRoot(n)
}

// NodeUtil_IsGoogModuleExportsReference: whether the given node is
// referencing the `exports` object in some goog.module.
//
// Java: NodeUtil.isGoogModuleExportsReference(Scope, Node)
func NodeUtil_IsGoogModuleExportsReference(scope *Scope, possibleName rhino.Node) bool {
	if !possibleName.IsName() || possibleName.GetString() != "exports" {
		return false
	}
	var_ := scope.GetVar(possibleName.GetString())
	var scopeRoot rhino.Node
	if var_ != nil {
		scopeRoot = var_.GetScopeRoot()
	}
	return !scopeRoot.IsNil() &&
		(scopeRoot.IsModuleBody() ||
			(scopeRoot.IsFunction() &&
				NodeUtil_IsBundledGoogModuleScopeRoot(NodeUtil_GetFunctionBody(scopeRoot))))
}

// Java: NodeUtil.GOOG_LOADMODULE
var NodeUtil_GOOG_LOADMODULE = rhino.QualifiedName_Of("goog.loadModule")

// Java: NodeUtil.isBundledGoogModuleCall(Node)
func NodeUtil_IsBundledGoogModuleCall(n rhino.Node) bool {
	// TODO(lharker): take an EXPR_RESULT to align with NodeUtil.isGoogModuleCall and
	// NodeUtil.isGoogProvideCall.
	if !(n.IsCall() && n.HasTwoChildren() && NodeUtil_GOOG_LOADMODULE.Matches(n.GetFirstChild())) {
		return false
	}
	return n.HasParent() &&
		n.GetParent().IsExprResult() &&
		!n.GetGrandparent().IsNil() &&
		n.GetGrandparent().IsScript()
}

// Java: NodeUtil.isBundledGoogModuleScopeRoot(Node)
func NodeUtil_IsBundledGoogModuleScopeRoot(n rhino.Node) bool {
	if !n.IsBlock() || !n.HasChildren() || !NodeUtil_IsGoogModuleCall(n.GetFirstChild()) {
		return false
	}
	function := n.GetParent()
	if function.IsNil() ||
		!function.IsFunction() ||
		!NodeUtil_GetFunctionParameters(function).HasOneChild() ||
		!NodeUtil_GetFunctionParameters(function).GetFirstChild().MatchesName("exports") {
		return false
	}
	call := function.GetParent()
	if !call.IsCall() ||
		!call.HasTwoChildren() ||
		!NodeUtil_GOOG_LOADMODULE.Matches(call.GetFirstChild()) {
		return false
	}
	return call.GetParent().IsExprResult() && call.GetGrandparent().IsScript()
}

// Java: NodeUtil.isGoogModuleDeclareLegacyNamespaceCall(Node)
func NodeUtil_IsGoogModuleDeclareLegacyNamespaceCall(n rhino.Node) bool {
	if NodeUtil_IsExprCall(n) {
		target := n.GetFirstFirstChild()
		return NodeUtil_GOOG_MODULE_DECLARE_LEGACY_NAMESPACE.Matches(target)
	}
	return false
}

// Java: NodeUtil.isGoogSetTestOnlyCall(Node)
func NodeUtil_IsGoogSetTestOnlyCall(n rhino.Node) bool {
	if NodeUtil_IsExprCall(n) {
		target := n.GetFirstFirstChild()
		return NodeUtil_GOOG_SET_TEST_ONLY.Matches(target)
	}
	return false
}

// Java: NodeUtil.isTopLevel(Node)
func NodeUtil_IsTopLevel(n rhino.Node) bool {
	return n.IsScript() || n.IsModuleBody()
}

// NodeUtil_IsGoogModuleFile returns whether the node is a goog.module file's
// SCRIPT node.
//
// Java: NodeUtil.isGoogModuleFile(Node)
func NodeUtil_IsGoogModuleFile(n rhino.Node) bool {
	return n.IsScript() &&
		n.HasChildren() &&
		n.GetFirstChild().IsModuleBody() &&
		NodeUtil_IsGoogModuleCall(n.GetFirstFirstChild())
}

// NodeUtil_IsLegacyGoogModuleFile returns whether the node is a SCRIPT node
// for a goog.module that has a declareLegacyNamespace call.
//
// Java: NodeUtil.isLegacyGoogModuleFile(Node)
func NodeUtil_IsLegacyGoogModuleFile(n rhino.Node) bool {
	return NodeUtil_IsGoogModuleFile(n) &&
		NodeUtil_IsGoogModuleDeclareLegacyNamespaceCall(n.GetFirstChild().GetSecondChild())
}

// Java: NodeUtil.isConstructor(Node)
func NodeUtil_IsConstructor(fnNode rhino.Node) bool {
	if fnNode.IsNil() || !fnNode.IsFunction() {
		return false
	}

	type_ := fnNode.GetJSType()
	jsDocInfo := NodeUtil_GetBestJSDocInfo(fnNode)
	color := fnNode.GetColor()

	return (type_ != nil && type_.IsConstructor()) ||
		(jsDocInfo != nil && jsDocInfo.IsConstructor()) ||
		(color != nil && color.IsConstructor()) ||
		NodeUtil_IsEs6Constructor(fnNode)
}

// Java: NodeUtil.isEs6ConstructorMemberFunctionDef(Node)
func NodeUtil_IsEs6ConstructorMemberFunctionDef(memberFunctionDef rhino.Node) bool {
	if !memberFunctionDef.IsMemberFunctionDef() {
		return false // not a member function at all
	}
	return memberFunctionDef.GetParent().IsClassMembers() && // is in a class
		!memberFunctionDef.IsStaticMember() && // constructors aren't static
		memberFunctionDef.GetString() == "constructor"
}

// Java: NodeUtil.isEs6Constructor(Node)
func NodeUtil_IsEs6Constructor(fnNode rhino.Node) bool {
	if !fnNode.IsFunction() {
		return false
	}
	memberFunctionDef := fnNode.GetParent()
	return !memberFunctionDef.IsNil() && NodeUtil_IsEs6ConstructorMemberFunctionDef(memberFunctionDef)
}

// Java: NodeUtil.isGetterOrSetter(Node)
func NodeUtil_IsGetterOrSetter(propNode rhino.Node) bool {
	if NodeUtil_IsGetOrSetKey(propNode) {
		return true
	}
	if !propNode.IsStringKey() || !propNode.GetFirstChild().IsFunction() {
		return false
	}
	keyName := propNode.GetString()
	return keyName == "get" || keyName == "set"
}

// Java: NodeUtil.isCallTo(Node, String)
func NodeUtil_IsCallTo(n rhino.Node, qualifiedName string) bool {
	return n.IsCall() && n.GetFirstChild().MatchesQualifiedName(qualifiedName)
}

// Java: NodeUtil.isCallTo(Node, QualifiedName)
func NodeUtil_IsCallToQualifiedName(n rhino.Node, qualifiedName rhino.QualifiedName) bool {
	return n.IsCall() && qualifiedName.Matches(n.GetFirstChild())
}

// NodeUtil_IsCallToNode is a faster version of NodeUtil_IsCallTo.
//
// n is the node to check if a call; targetMethod is the prebuilt AST getprop
// node that represents the method to check.
//
// Java: NodeUtil.isCallTo(Node, Node)
func NodeUtil_IsCallToNode(n rhino.Node, targetMethod rhino.Node) bool {
	if !n.IsCall() {
		return false
	}

	return n.GetFirstChild().MatchesQualifiedNameNode(targetMethod)
}

// Java: NodeUtil.collectExternVariableNames(AbstractCompiler, Node)
func NodeUtil_CollectExternVariableNames(
	compiler AbstractCompiler, externs rhino.Node) *jutil.LinkedHashSet[string] {
	externsRefs := NewReferenceCollector(
		compiler, ReferenceCollector_DO_NOTHING_BEHAVIOR, NewSyntacticScopeCreator(compiler))
	externsRefs.ProcessRoot(externs)
	externsNames := jutil.NewLinkedHashSet[string]()
	for v := range externsRefs.GetAllSymbols() {
		if !v.IsParam() {
			externsNames.Add(v.GetName())
		}
	}
	return externsNames
}

// Java: NodeUtil.createSynthesizedExternsSymbol(AbstractCompiler, String)
func NodeUtil_CreateSynthesizedExternsSymbol(compiler AbstractCompiler, nameToAdd string) {
	name := rhino.IR_Name(nameToAdd)
	name.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
	var_ := rhino.IR_Var(name)
	input := compiler.GetSynthesizedExternsInput()
	root := input.GetAstRoot(compiler)
	name.SetStaticSourceFileFrom(root)
	var_.SetStaticSourceFileFrom(root)
	root.AddChildToBack(var_)
	compiler.ReportChangeToEnclosingScope(var_)
}

// NodeUtil_MarkNewScopesChanged recurses through a tree, marking all
// function nodes as changed.
//
// Java: NodeUtil.markNewScopesChanged(Node, AbstractCompiler)
func NodeUtil_MarkNewScopesChanged(node rhino.Node, compiler AbstractCompiler) {
	if node.IsFunction() {
		compiler.ReportChangeToChangeScope(node)
	}
	for child := node.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		NodeUtil_MarkNewScopesChanged(child, compiler)
	}
}

// NodeUtil_MarkFunctionsDeleted recurses through a tree, marking all
// function nodes deleted.
//
// Java: NodeUtil.markFunctionsDeleted(Node, AbstractCompiler)
func NodeUtil_MarkFunctionsDeleted(node rhino.Node, compiler AbstractCompiler) {
	if node.IsFunction() {
		compiler.ReportFunctionDeleted(node)
	}

	for child := node.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		NodeUtil_MarkFunctionsDeleted(child, compiler)
	}

	shadowed := node.GetClosureUnawareShadow()
	if !shadowed.IsNil() {
		NodeUtil_MarkFunctionsDeleted(shadowed, compiler)
	}
}

// NodeUtil_GetParentChangeScopeNodes returns the list of scope nodes which
// are parents of the provided list of scope nodes.
//
// Java: NodeUtil.getParentChangeScopeNodes(List<Node>)
func NodeUtil_GetParentChangeScopeNodes(scopeNodes []rhino.Node) []rhino.Node {
	parentScopeNodes := jutil.NewLinkedHashSet(scopeNodes...)
	for _, scopeNode := range scopeNodes {
		parentScopeNodes.Add(ChangeTracker_GetEnclosingChangeScopeRoot(scopeNode))
	}
	return nodeUtil_newArrayList(parentScopeNodes)
}

// NodeUtil_RemoveNestedChangeScopeNodes removes any scope nodes from the
// provided list that are nested within some other scope node also in the
// list. Returns the modified list.
//
// Java: NodeUtil.removeNestedChangeScopeNodes(List<Node>)
func NodeUtil_RemoveNestedChangeScopeNodes(scopeNodes []rhino.Node) []rhino.Node {
	uniqueScopeNodes := jutil.NewLinkedHashSet(scopeNodes...)
	for _, scopeNode := range scopeNodes {
		for ancestor := scopeNode.GetParent(); !ancestor.IsNil(); ancestor = ancestor.GetParent() {
			if ChangeTracker_IsChangeScopeRoot(ancestor) && uniqueScopeNodes.Contains(ancestor) {
				uniqueScopeNodes.Remove(scopeNode)
				break
			}
		}
	}
	return nodeUtil_newArrayList(uniqueScopeNodes)
}

// nodeUtil_newArrayList is new ArrayList<>(set).
//
// Go-only: language-adapter — the new list is never null, also when it is
// empty, and a nil slice is Java null to the callers (J2clClinitPrunerPass
// reads a null list of changed scopes as "traverse everything").
func nodeUtil_newArrayList(set *jutil.LinkedHashSet[rhino.Node]) []rhino.Node {
	list := make([]rhino.Node, 0, set.Size())
	for n := range set.All() {
		list = append(list, n)
	}
	return list
}

// Java: NodeUtil.getInvocationArgsAsIterable(Node)
func NodeUtil_GetInvocationArgsAsIterable(invocation rhino.Node) iter.Seq[rhino.Node] {
	if invocation.IsTaggedTemplateLit() {
		return NewNodeUtil_TemplateArgsIterable(invocation.GetLastChild()).All()
	}

	jutil.CheckStateMsg1(NodeUtil_IsCallOrNew(invocation), "%s", invocation)
	if invocation.HasOneChild() {
		return slices.Values([]rhino.Node(nil))
	}

	var list []rhino.Node
	for arg := invocation.GetSecondChild(); !arg.IsNil(); arg = arg.GetNext() {
		list = append(list, arg)
	}
	return slices.Values(list)
}

// NodeUtil_GetInvocationArgsCount returns the number of arguments in this
// invocation. For template literals it takes into account the implicit first
// argument of ITemplateArray.
//
// Java: NodeUtil.getInvocationArgsCount(Node)
func NodeUtil_GetInvocationArgsCount(invocation rhino.Node) int {
	if invocation.IsTaggedTemplateLit() {
		args := NewNodeUtil_TemplateArgsIterable(invocation.GetLastChild())
		size := 0
		for range args.All() {
			size++
		}
		return size + 1
	} else {
		return invocation.GetChildCount() - 1
	}
}

// NodeUtil_TemplateArgsIterable represents an iterable of the children of
// templatelit_sub nodes of a template lit node This iterable will skip over
// the String children of the template lit node.
//
// Java: NodeUtil.TemplateArgsIterable
type NodeUtil_TemplateArgsIterable struct {
	templateLit rhino.Node
}

// Java: NodeUtil.TemplateArgsIterable(Node)
func NewNodeUtil_TemplateArgsIterable(templateLit rhino.Node) *NodeUtil_TemplateArgsIterable {
	jutil.CheckState(templateLit.IsTemplateLit())
	return &NodeUtil_TemplateArgsIterable{templateLit: templateLit}
}

// Iterator returns the anonymous AbstractIterator<Node> of the Java class.
//
// Java: NodeUtil.TemplateArgsIterable.iterator()
func (t *NodeUtil_TemplateArgsIterable) Iterator() *nodeUtil_templateArgsIterator {
	return &nodeUtil_templateArgsIterator{nextChild: t.templateLit.GetFirstChild()}
}

// All iterates the iterable.
//
// Go-only: language-adapter — Java's for-each over the Iterable.
func (t *NodeUtil_TemplateArgsIterable) All() iter.Seq[rhino.Node] {
	return func(yield func(rhino.Node) bool) {
		it := t.Iterator()
		for it.HasNext() {
			if !yield(it.Next()) {
				return
			}
		}
	}
}

// nodeUtil_templateArgsIterator is the anonymous AbstractIterator<Node>
// returned by TemplateArgsIterable.iterator().
type nodeUtil_templateArgsIterator struct {
	nextChild rhino.Node

	// Go-only: language-adapter — AbstractIterator's state.
	next    rhino.Node
	hasNext bool
	done    bool
}

// HasNext is AbstractIterator.hasNext().
func (it *nodeUtil_templateArgsIterator) HasNext() bool {
	if it.done {
		return false
	}
	if !it.hasNext {
		next, ok := it.ComputeNext()
		if !ok {
			it.done = true
			return false
		}
		it.next = next
		it.hasNext = true
	}
	return true
}

// Next is AbstractIterator.next().
func (it *nodeUtil_templateArgsIterator) Next() rhino.Node {
	if !it.HasNext() {
		panic(jutil.NewNoSuchElementException())
	}
	it.hasNext = false
	result := it.next
	it.next = rhino.NilNode
	return result
}

// ComputeNext returns the next argument, or ok=false for endOfData().
//
// Java: NodeUtil.TemplateArgsIterable.iterator().computeNext()
func (it *nodeUtil_templateArgsIterator) ComputeNext() (next rhino.Node, ok bool) {
	for !it.nextChild.IsNil() && !it.nextChild.IsTemplateLitSub() {
		it.nextChild = it.nextChild.GetNext()
	}
	if it.nextChild.IsNil() {
		return rhino.NilNode, false
	} else {
		result := it.nextChild.GetFirstChild()
		it.nextChild = it.nextChild.GetNext()
		return result, true
	}
}

// NodeUtil_GetAllVarNamesDeclaredInModule records a mapping of names to vars
// of everything reachable in a module. Should only be called with a module
// scope.
//
// nameVarMap: an empty map that gets populated with the keys being variable
// names and values being variable objects. orderedVars: an empty list that
// gets populated with variable objects in the order that they appear in the
// module.
//
// Java: NodeUtil.getAllVarNamesDeclaredInModule(Node, AbstractCompiler, ScopeCreator, Scope)
func NodeUtil_GetAllVarNamesDeclaredInModule(
	moduleNode rhino.Node,
	compiler AbstractCompiler,
	scopeCreator ScopeCreator,
	globalScope *Scope) *jutil.LinkedHashSet[string] {

	jutil.CheckStateMsg(moduleNode.IsModuleBody(), "getAllVarsDeclaredInModule expects a module body node")
	jutil.CheckStateMsg1(globalScope.IsGlobal(), "%s", globalScope)

	nameVars := jutil.NewLinkedHashSet[string]()

	finder := &nodeUtil_getAllVarNamesDeclaredInModule_finder{nameVars: nameVars}
	NodeTraversal_Builder_().
		SetCompiler(compiler).
		SetCallback(finder).
		SetScopeCreator(scopeCreator).
		TraverseWithScope(moduleNode, globalScope)
	return nameVars
}

// nodeUtil_getAllVarNamesDeclaredInModule_finder is the anonymous
// ScopedCallback of getAllVarNamesDeclaredInModule.
type nodeUtil_getAllVarNamesDeclaredInModule_finder struct {
	nameVars *jutil.LinkedHashSet[string]
}

// Java: NodeUtil.getAllVarNamesDeclaredInModule(...).finder.enterScope(NodeTraversal)
func (f *nodeUtil_getAllVarNamesDeclaredInModule_finder) EnterScope(t *NodeTraversal) {
	currentScope := t.GetScope()
	if currentScope.IsModuleScope() {
		for v := range currentScope.GetVarIterable() {
			f.nameVars.Add(v.GetName())
		}
	}
}

// Java: NodeUtil.getAllVarNamesDeclaredInModule(...).finder.exitScope(NodeTraversal)
func (f *nodeUtil_getAllVarNamesDeclaredInModule_finder) ExitScope(t *NodeTraversal) {}

// Java: NodeUtil.getAllVarNamesDeclaredInModule(...).finder.shouldTraverse(NodeTraversal, Node, Node)
func (f *nodeUtil_getAllVarNamesDeclaredInModule_finder) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	return n.IsModuleBody()
}

// Java: NodeUtil.getAllVarNamesDeclaredInModule(...).finder.visit(NodeTraversal, Node, Node)
func (f *nodeUtil_getAllVarNamesDeclaredInModule_finder) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
}

// NodeUtil_GetAllVarsDeclaredInFunction returns a mapping of names to vars
// of everything reachable in a function. Should only be called with a
// function scope. Does not enter new control flow areas aka embedded
// functions.
//
// Java: NodeUtil.getAllVarsDeclaredInFunction(AbstractCompiler, ScopeCreator, Scope)
func NodeUtil_GetAllVarsDeclaredInFunction(
	compiler AbstractCompiler, scopeCreator ScopeCreator, scope *Scope) *NodeUtil_AllVarsDeclaredInFunction {
	jutil.CheckStateMsg1(scope.IsFunctionScope(), "%s", scope)

	finder := &nodeUtil_getAllVarsDeclaredInFunction_finder{
		scope: scope,
	}
	finder.result.allVarsInFn = &finder.result.mapStorage

	b := NodeTraversal_Builder{
		compiler:     compiler,
		callback:     finder,
		scopeCreator: scopeCreator,
	}
	b.TraverseAtScope(scope)

	jutil.CheckState(finder.result.allVarsInFn.IsEmpty() == (len(finder.result.orderedVars) == 0))
	return &finder.result
}

// nodeUtil_getAllVarsDeclaredInFunction_finder is the anonymous
// ScopedCallback of getAllVarsDeclaredInFunction.
type nodeUtil_getAllVarsDeclaredInFunction_finder struct {
	result NodeUtil_AllVarsDeclaredInFunction
	scope  *Scope
}

func nodeUtilSubtreeHasBlockScopeDecls(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_LET, rhino.Token_CONST, rhino.Token_CLASS, rhino.Token_CATCH, rhino.Token_IMPORT:
		return true
	case rhino.Token_FUNCTION:
		return !NodeUtil_IsFunctionExpression(n)
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if nodeUtilSubtreeHasBlockScopeDecls(c) {
			return true
		}
	}
	return false
}

func nodeUtilStmtHasBlockScopeDecls(stmt rhino.Node) bool {
	switch stmt.GetToken() {
	case rhino.Token_FUNCTION:
		return false
	case rhino.Token_CLASS, rhino.Token_CATCH, rhino.Token_IMPORT:
		return true
	default:
		for c := stmt.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			if nodeUtilSubtreeHasBlockScopeDecls(c) {
				return true
			}
		}
		return false
	}
}

func nodeUtilBlockScopeMayHaveDecls(n, root rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_LET, rhino.Token_CONST, rhino.Token_CATCH, rhino.Token_IMPORT:
		return true
	case rhino.Token_EXPORT:
		first := n.GetFirstChild()
		return !first.IsNil() && nodeUtilBlockScopeMayHaveDecls(first, root)
	case rhino.Token_FUNCTION:
		return !NodeUtil_IsFunctionExpression(n) && n.GetFirstChild().GetString() != ""
	case rhino.Token_CLASS:
		return !NodeUtil_IsClassExpression(n) && n.GetFirstChild().GetString() != ""
	case rhino.Token_BLOCK:
		if n != root && scanVarsBlockCreatesScope(n) {
			return false
		}
	case rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF, rhino.Token_SWITCH_BODY:
		if n != root {
			return false
		}
	case rhino.Token_ROOT, rhino.Token_DO, rhino.Token_WHILE, rhino.Token_WITH, rhino.Token_IF,
		rhino.Token_LABEL, rhino.Token_TRY, rhino.Token_SWITCH, rhino.Token_CASE, rhino.Token_DEFAULT_CASE:
	default:
		return false
	}
	for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		if scanVarsTokenKind[uint8(child.GetToken())] == 1 {
			if nodeUtilBlockScopeMayHaveDecls(child, root) {
				return true
			}
		}
	}
	return false
}

func nodeUtilScopeRootMayDeclareVars(root, fnRoot rhino.Node) bool {
	switch root.GetToken() {
	case rhino.Token_COMPUTED_FIELD_DEF, rhino.Token_MEMBER_FIELD_DEF:
		return false
	case rhino.Token_CLASS:
		return !root.GetFirstChild().IsEmpty() && NodeUtil_IsClassExpression(root)
	case rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF, rhino.Token_SWITCH_BODY:
		return nodeUtilBlockScopeMayHaveDecls(root, root)
	case rhino.Token_BLOCK:
		parent := root.GetParent()
		if parent == fnRoot || (!parent.IsNil() && parent.IsClassMembers()) {
			return root.HasChildren()
		}
		return nodeUtilBlockScopeMayHaveDecls(root, root)
	default:
		return true
	}
}

func (f *nodeUtil_getAllVarsDeclaredInFunction_finder) appendScopeVars(currentScope *Scope) {
	n := currentScope.GetVarCount()
	if n == 0 {
		return
	}
	if need := len(f.result.orderedVars) + n; need > cap(f.result.orderedVars) {
		newCap := max(need, 2*cap(f.result.orderedVars))
		buf := make([]*Var, len(f.result.orderedVars), newCap)
		copy(buf, f.result.orderedVars)
		f.result.orderedVars = buf
	}
	for _, v := range currentScope.varList {
		if v == nil {
			continue
		}
		f.result.mapStorage.Put(v.GetName(), v)
		f.result.orderedVars = append(f.result.orderedVars, v)
	}
}

// Java: NodeUtil.getAllVarsDeclaredInFunction(...).finder.enterScope(NodeTraversal)
func (f *nodeUtil_getAllVarsDeclaredInFunction_finder) EnterScope(t *NodeTraversal) {
	root := t.GetScopeRoot()
	var currentScope *Scope
	if root == f.scope.GetRootNode() {
		currentScope = f.scope
	} else {
		if !nodeUtilScopeRootMayDeclareVars(root, f.scope.GetRootNode()) {
			return
		}
		currentScope = t.GetScope()
	}
	f.appendScopeVars(currentScope)
}

func (*nodeUtil_getAllVarsDeclaredInFunction_finder) nodeTraversalNoVisit() {}

// Java: NodeUtil.getAllVarsDeclaredInFunction(...).finder.exitScope(NodeTraversal)
func (f *nodeUtil_getAllVarsDeclaredInFunction_finder) ExitScope(t *NodeTraversal) {}

// Java: NodeUtil.getAllVarsDeclaredInFunction(...).finder.shouldTraverse(NodeTraversal, Node, Node)
func (f *nodeUtil_getAllVarsDeclaredInFunction_finder) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	fnRoot := f.scope.GetRootNode()
	if n == fnRoot {
		return true
	}
	if n.IsFunction() {
		// Don't enter any new functions
		return false
	}
	if parent == fnRoot {
		if n.IsBlock() {
			return n.HasChildren()
		}
		if n.IsParamList() {
			for p := n.GetFirstChild(); !p.IsNil(); p = p.GetNext() {
				if !p.IsName() && nodeUtilSubtreeHasBlockScopeDecls(p) {
					return true
				}
			}
			return false
		}
		if n == fnRoot.GetFirstChild() {
			return false
		}
		return nodeUtilSubtreeHasBlockScopeDecls(n)
	}
	if !parent.IsNil() && parent.GetParent() == fnRoot && parent.IsBlock() {
		return nodeUtilStmtHasBlockScopeDecls(n)
	}
	return true
}

// Java: NodeUtil.getAllVarsDeclaredInFunction(...).finder.visit(NodeTraversal, Node, Node)
func (f *nodeUtil_getAllVarsDeclaredInFunction_finder) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
}

// NodeUtil_AllVarsDeclaredInFunction represents a mapping of names to vars
// of everything reachable in a function.
//
// Java: NodeUtil.AllVarsDeclaredInFunction
type NodeUtil_AllVarsDeclaredInFunction struct {
	allVarsInFn *jutil.LinkedHashMap[string, *Var]
	orderedVars []*Var
	mapStorage  jutil.LinkedHashMap[string, *Var]
}

// Java: NodeUtil.AllVarsDeclaredInFunction(Map<String, Var>, List<Var>)
func NewNodeUtil_AllVarsDeclaredInFunction(allVarsInFn *jutil.LinkedHashMap[string, *Var], orderedVars []*Var) *NodeUtil_AllVarsDeclaredInFunction {
	jutil.CheckState(allVarsInFn.IsEmpty() == (len(orderedVars) == 0))
	return &NodeUtil_AllVarsDeclaredInFunction{allVarsInFn: allVarsInFn, orderedVars: orderedVars}
}

// Java: NodeUtil.AllVarsDeclaredInFunction.getAllVariables()
func (a *NodeUtil_AllVarsDeclaredInFunction) GetAllVariables() *jutil.LinkedHashMap[string, *Var] {
	return a.allVarsInFn
}

// Java: NodeUtil.AllVarsDeclaredInFunction.getAllVariablesInOrder()
func (a *NodeUtil_AllVarsDeclaredInFunction) GetAllVariablesInOrder() []*Var {
	return a.orderedVars
}

// NodeUtil_IsObjLitProperty returns true if the node is a property of an
// object literal.
//
// Java: NodeUtil.isObjLitProperty(Node)
func NodeUtil_IsObjLitProperty(node rhino.Node) bool {
	return node.IsStringKey() ||
		node.IsGetterDef() ||
		node.IsSetterDef() ||
		node.IsMemberFunctionDef() ||
		node.IsComputedProp()
}

// NodeUtil_IsBlocklessArrowFunctionResult returns whether the node
// represents the return value of a blockless Arrow function.
//
// Java: NodeUtil.isBlocklessArrowFunctionResult(Node)
func NodeUtil_IsBlocklessArrowFunctionResult(n rhino.Node) bool {
	parent := n.GetParent()
	return !parent.IsNil() && parent.IsFunction() && n == parent.GetLastChild() && !n.IsBlock()
}

// NodeUtil_GetFeatureSetOfScript returns a script node's FeatureSet, which
// is set at parse-time. This may not be up-to-date as passes can add/remove
// features from a script node's descendants.
//
// The feature set will be null if the script node was created artificially
// or if the parser didn't detect any interesting features.
//
// Java: NodeUtil.getFeatureSetOfScript(Node)
func NodeUtil_GetFeatureSetOfScript(scriptNode rhino.Node) *parser.FeatureSet {
	jutil.CheckStateMsg1(scriptNode.IsScript(), "%s", scriptNode)
	prop := scriptNode.GetProp(rhino.Node_FEATURE_SET)
	if prop == nil {
		return nil
	}
	return prop.(*parser.FeatureSet)
}

// NodeUtil_AddFeatureToScript adds the given features to a SCRIPT node's
// FeatureSet property.
//
// Also updates the compiler's FeatureSet.
//
// Java: NodeUtil.addFeatureToScript(Node, Feature, AbstractCompiler)
func NodeUtil_AddFeatureToScript(scriptNode rhino.Node, feature parser.FeatureSet_Feature, compiler AbstractCompiler) {
	jutil.CheckStateMsg1(scriptNode.IsScript(), "%s", scriptNode)
	currentFeatures := NodeUtil_GetFeatureSetOfScript(scriptNode)
	var newFeatures *parser.FeatureSet
	if currentFeatures != nil {
		newFeatures = currentFeatures.With(feature)
	} else {
		newFeatures = parser.FeatureSet_BARE_MINIMUM.With(feature)
	}
	scriptNode.PutProp(rhino.Node_FEATURE_SET, newFeatures)
	if feature != parser.FeatureSet_Feature_MODULES {
		// Except MODULES (which can get reintroduced later in ConvertChunksToEsModules pass),
		// prohibit passes from reintroducing any feature that already got transpiled.
		jutil.CheckStateMsg1(
			compiler.GetAllowableFeatures().ContainsFeature(feature),
			"Cannot add feature: %s. It is not supported in the output language, and either 1) its"+
				" corresponding transpilation pass has already run or 2) transpilation of this"+
				" feature is unsupported entirely",
			feature)
	} else {
		compiler.SetAllowableFeatures(compiler.GetAllowableFeatures().With(feature))
	}
}

// NodeUtil_RemoveFeatureFromScript removes the given feature from the SCRIPT
// node's FeatureSet property.
//
// Removing a feature from a single script does not mean that it's removed
// from all scripts. Hence this method does not remove it from the compiler's
// featureSet. The caller must remove it from the compiler's featureSet if
// it's calling this method for each script.
//
// Java: NodeUtil.removeFeatureFromScript(Node, Feature)
func NodeUtil_RemoveFeatureFromScript(scriptNode rhino.Node, feature parser.FeatureSet_Feature) {
	currentFeatures := NodeUtil_GetFeatureSetOfScript(scriptNode)
	if currentFeatures == nil {
		scriptNode.PutProp(rhino.Node_FEATURE_SET, nil)
		return
	}
	updated := currentFeatures.WithoutFeatures(feature)
	if updated != currentFeatures {
		scriptNode.PutProp(rhino.Node_FEATURE_SET, updated)
	}
}

// NodeUtil_RemoveFeaturesFromScript removes the given featureSet from the
// SCRIPT node's FeatureSet property.
//
// Removing a feature from a single script does not mean that it's removed
// from all scripts. Hence this method does not remove it from the compiler's
// featureSet. The caller must remove it from the compiler's featureSet using
// AbstractCompiler.markFeatureSetNotAllowed or
// AbstractCompiler.markFeatureNotAllowed if it's calling this method for each
// script.
//
// Java: NodeUtil.removeFeaturesFromScript(Node, FeatureSet)
func NodeUtil_RemoveFeaturesFromScript(scriptNode rhino.Node, featureSet *parser.FeatureSet) {
	currentFeatures := NodeUtil_GetFeatureSetOfScript(scriptNode)
	if currentFeatures != nil {
		updated := currentFeatures.Without(featureSet)
		if updated != currentFeatures {
			scriptNode.PutProp(rhino.Node_FEATURE_SET, updated)
		}
	}
}

// NodeUtil_RemoveFeatureFromAllScripts removes the given feature from the
// FEATURE_SET prop of all SCRIPT nodes under root.
//
// Also updates the compiler's FeatureSet.
//
// Java: NodeUtil.removeFeatureFromAllScripts(Node, Feature, AbstractCompiler)
func NodeUtil_RemoveFeatureFromAllScripts(root rhino.Node, feature parser.FeatureSet_Feature, compiler AbstractCompiler) {
	jutil.CheckArgumentMsg1(root.IsRoot(), "%s", root)
	for childNode := root.GetFirstChild(); !childNode.IsNil(); childNode = childNode.GetNext() {
		jutil.CheckState(childNode.IsScript())
		NodeUtil_RemoveFeatureFromScript(childNode, feature)
	}
	compiler.MarkFeatureNotAllowed(feature)
}

// NodeUtil_RemoveFeaturesFromAllScripts removes the given feature from the
// FEATURE_SET prop of all SCRIPT nodes under root.
//
// Also updates the compiler's FeatureSet.
//
// Java: NodeUtil.removeFeaturesFromAllScripts(Node, FeatureSet, AbstractCompiler)
func NodeUtil_RemoveFeaturesFromAllScripts(
	root rhino.Node, featureSet *parser.FeatureSet, compiler AbstractCompiler) {
	jutil.CheckArgumentMsg1(root.IsRoot(), "%s", root)
	for childNode := root.GetFirstChild(); !childNode.IsNil(); childNode = childNode.GetNext() {
		jutil.CheckState(childNode.IsScript())
		NodeUtil_RemoveFeaturesFromScript(childNode, featureSet)
	}
	compiler.MarkFeatureSetNotAllowed(featureSet)
}

// NodeUtil_AddFeatureToAllScripts adds the given feature to the FEATURE_SET
// prop of all scripts under externs and root.
//
// Java: NodeUtil.addFeatureToAllScripts(Node, Feature, AbstractCompiler)
func NodeUtil_AddFeatureToAllScripts(root rhino.Node, feature parser.FeatureSet_Feature, compiler AbstractCompiler) {
	jutil.CheckArgumentMsg1(root.IsRoot(), "%s", root)
	for childNode := root.GetFirstChild(); !childNode.IsNil(); childNode = childNode.GetNext() {
		jutil.CheckState(childNode.IsScript())
		NodeUtil_AddFeatureToScript(childNode, feature, compiler)
	}
}

// NodeUtil_AddFeaturesToScript adds the given features to a SCRIPT node's
// FeatureSet property.
//
// Also updates the compiler's FeatureSet.
//
// Java: NodeUtil.addFeaturesToScript(Node, FeatureSet, AbstractCompiler)
func NodeUtil_AddFeaturesToScript(scriptNode rhino.Node, features *parser.FeatureSet, compiler AbstractCompiler) {
	jutil.CheckStateMsg1(scriptNode.IsScript(), "%s", scriptNode)
	for feature := range features.All() {
		NodeUtil_AddFeatureToScript(scriptNode, feature, compiler)
	}
}

// NodeUtil_GetParamOrPatternNames calls cb with all NAMEs declared in a
// PARAM_LIST or destructuring pattern.
//
// Java: NodeUtil.getParamOrPatternNames(Node, Consumer<Node>)
func NodeUtil_GetParamOrPatternNames(n rhino.Node, cb func(rhino.Node)) {
	parsing.ParsingUtil_GetParamOrPatternNames(n, cb)
}

// NodeUtil_GoogRequire represents a goog.require'd namespace and property
// inside a module.
//
// namespace is the Closure namespace inside the require call; property is
// the property off the required namespace. Needed only for destructuring
// requires. isStrongRequire is true if the require is a goog.require, false
// if it is a goog.requireType.
//
// Java: NodeUtil.GoogRequire
type NodeUtil_GoogRequire struct {
	namespace       string
	property        *string
	isStrongRequire bool
}

// Java: NodeUtil.GoogRequire(String, String, boolean)
func NewNodeUtil_GoogRequire(namespace string, property *string, isStrongRequire bool) *NodeUtil_GoogRequire {
	// Go-only: language-adapter — requireNonNull(namespace, "namespace"): a Go
	// string cannot be null.
	return &NodeUtil_GoogRequire{namespace: namespace, property: property, isStrongRequire: isStrongRequire}
}

// Java: NodeUtil.GoogRequire.namespace()
func (r *NodeUtil_GoogRequire) Namespace() string { return r.namespace }

// Java: NodeUtil.GoogRequire.property()
func (r *NodeUtil_GoogRequire) Property() *string { return r.property }

// Java: NodeUtil.GoogRequire.isStrongRequire()
func (r *NodeUtil_GoogRequire) IsStrongRequire() bool { return r.isStrongRequire }

// Equals is the record's component-wise equals.
//
// Java: NodeUtil.GoogRequire.equals(Object)
func (r *NodeUtil_GoogRequire) Equals(o any) bool {
	other, ok := o.(*NodeUtil_GoogRequire)
	if !ok || other == nil {
		return false
	}
	if r == other {
		return true
	}
	if (r.property == nil) != (other.property == nil) ||
		(r.property != nil && *r.property != *other.property) {
		return false
	}
	return r.namespace == other.namespace && r.isStrongRequire == other.isStrongRequire
}

// HashCode is the record's hashCode: 31 * h + hash(component), in
// declaration order.
//
// Java: NodeUtil.GoogRequire.hashCode()
func (r *NodeUtil_GoogRequire) HashCode() int32 {
	var h int32
	h = 31*h + jutil.StringHashCode(r.namespace)
	var propertyHash int32
	if r.property != nil {
		propertyHash = jutil.StringHashCode(*r.property)
	}
	h = 31*h + propertyHash
	var strongHash int32 = 1237
	if r.isStrongRequire {
		strongHash = 1231
	}
	h = 31*h + strongHash
	return h
}

// String is the record's toString.
//
// Java: NodeUtil.GoogRequire.toString()
func (r *NodeUtil_GoogRequire) String() string {
	return "GoogRequire[namespace=" + r.namespace + ", property=" + jutil.ToJavaString(r.property) +
		", isStrongRequire=" + jutil.ToJavaString(r.isStrongRequire) + "]"
}

// Java: NodeUtil.GoogRequire.fromNamespace(String, boolean)
func NodeUtil_GoogRequire_FromNamespace(namespace string, isStrongRequire bool) *NodeUtil_GoogRequire {
	return NewNodeUtil_GoogRequire(namespace /* property= */, nil, isStrongRequire)
}

// Java: NodeUtil.GoogRequire.fromNamespaceAndProperty(String, String, boolean)
func NodeUtil_GoogRequire_FromNamespaceAndProperty(
	namespace string, property string, isStrongRequire bool) *NodeUtil_GoogRequire {
	return NewNodeUtil_GoogRequire(namespace, jutil.StringPtr(property), isStrongRequire)
}

// Java: NodeUtil.getGoogRequireInfo(String, Scope)
func NodeUtil_GetGoogRequireInfoWithScope(name string, scope *Scope) *NodeUtil_GoogRequire {
	var_ := scope.GetVar(name)
	if var_ == nil {
		return nil
	}
	return NodeUtil_GetGoogRequireInfo(var_)
}

// Java: NodeUtil.getGoogRequireInfo(Var)
func NodeUtil_GetGoogRequireInfo(var_ *Var) *NodeUtil_GoogRequire {
	jutil.CheckNotNull(var_)
	if !var_.GetScopeRoot().IsModuleBody() || var_.GetNameNode().IsNil() {
		return nil
	}

	return NodeUtil_GetGoogRequireInfoNode(var_.GetNameNode())
}

// Java: NodeUtil.getGoogRequireInfo(Node)
func NodeUtil_GetGoogRequireInfoNode(nameNode rhino.Node) *NodeUtil_GoogRequire {
	// IMPORT_STAR is not part of a GoogRequire call.
	if nameNode.IsImportStar() {
		return nil
	}
	jutil.CheckStateMsg1(nameNode.IsName(), "unexpected node type: %s", nameNode)

	if NodeUtil_IsNameDeclaration(nameNode.GetParent()) {
		requireCall := nameNode.GetFirstChild()
		if requireCall.IsNil() {
			return nil
		}
		isStrongRequire := NodeUtil_IsGoogRequireCall(requireCall)
		if !(isStrongRequire || NodeUtil_IsGoogRequireTypeCall(requireCall)) {
			return nil
		}
		namespace := requireCall.GetSecondChild().GetString()
		return NodeUtil_GoogRequire_FromNamespace(namespace, isStrongRequire)
	} else if nameNode.GetParent().IsStringKey() && nameNode.GetGrandparent().IsObjectPattern() {
		requireCall := nameNode.GetGrandparent().GetNext()
		if requireCall.IsNil() {
			return nil
		}
		isStrongRequire := NodeUtil_IsGoogRequireCall(requireCall)
		if !(isStrongRequire || NodeUtil_IsGoogRequireTypeCall(requireCall)) {
			return nil
		}
		property := nameNode.GetParent().GetString()
		namespace := requireCall.GetSecondChild().GetString()
		return NodeUtil_GoogRequire_FromNamespaceAndProperty(namespace, property, isStrongRequire)
	}
	return nil
}

// NodeUtil_EstimateNumLines estimates the number of lines in the file of
// this script node.
//
// This method returns the line number of the last node in the script +1.
// This is not strictly the number of lines in the file (consider trailing
// comments or whitespace), but should be strongly correlated with it. If
// perfect accuracy is desired the original source file will have to be read.
//
// Java: NodeUtil.estimateNumLines(Node)
func NodeUtil_EstimateNumLines(scriptNode rhino.Node) int {
	jutil.CheckArgument(scriptNode.IsScript())
	current := scriptNode
	for current.HasChildren() {
		current = current.GetLastChild()
	}

	// +1 since this is the start line of the last AST node and we can expect at least one
	// trailing newline
	return current.GetLineno() + 1
}

// NodeUtil_GetOriginalName returns the original name of a renamed variable
// by stripping unique suffixes and module prefixes.
//
// Java: NodeUtil.getOriginalName(String)
func NodeUtil_GetOriginalName(name string) string {
	// Strip unique name suffix: "foo$$1" -> "foo"
	name = jutil.StringReplaceAll(name, "\\$\\$\\d+$", "")

	// Strip contextual rename suffix: "foo$jscomp$1" -> "foo"
	name = MakeDeclaredNamesUnique_ContextualRenameInverter_GetOriginalName(name)

	// Strip module exports prefix: "module$exports$foo$Bar" -> "Bar"
	if strings.HasPrefix(name, "module$exports$") {
		lastDollar := jutil.StringLastIndexOf(name, "$")
		if lastDollar != -1 {
			name = jutil.StringSubstringFrom(name, lastDollar+1)
		}
	} else if strings.HasPrefix(name, "module$contents$") {
		// Strip module contents prefix: "module$contents$foo$Bar_baz" -> "baz"
		lastUnderscore := jutil.StringLastIndexOf(name, "_")
		if lastUnderscore != -1 {
			name = jutil.StringSubstringFrom(name, lastUnderscore+1)
		}
	}
	return name
}

// nodeUtilFormatInt formats v in base 10 (the uniqueNameIdSupplier hot path).
func nodeUtilFormatInt(v int) string {
	return strconv.Itoa(v)
}
