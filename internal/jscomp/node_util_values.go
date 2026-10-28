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
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/dtoa"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

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
		if part, ok := ast.MustTemplateLit(n).OnlyPart(); ok {
			cookedString, ok := part.MustTemplateLitString().Cooked()
			return base.Tri_ForBoolean(ok && cookedString != "")
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
		return NodeUtil_GetBooleanValue(ast.MustUnary(n).Operand().Node()).Not()

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
	case rhino.Token_ASSIGN:
		return NodeUtil_GetBooleanValue(ast.MustAssign(n).Value().Node())
	case rhino.Token_COMMA:
		return NodeUtil_GetBooleanValue(ast.MustBinary(n).Right().Node())
	case rhino.Token_AND, rhino.Token_ASSIGN_AND:
		lhs := NodeUtil_GetBooleanValue(n.GetFirstChild())
		rhs := NodeUtil_GetBooleanValue(n.GetLastChild())
		return lhs.And(rhs)
	case rhino.Token_OR, rhino.Token_ASSIGN_OR:
		lhs := NodeUtil_GetBooleanValue(n.GetFirstChild())
		rhs := NodeUtil_GetBooleanValue(n.GetLastChild())
		return lhs.Or(rhs)
	case rhino.Token_HOOK:
		hook := ast.MustHook(n)
		trueValue := NodeUtil_GetBooleanValue(hook.Then().Node())
		falseValue := NodeUtil_GetBooleanValue(hook.Else().Node())
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
	s, ok := nodeUtil_getStringValue(n)
	if !ok {
		return nil
	}
	return &s
}

func nodeUtil_getStringValue(n rhino.Node) (string, bool) {
	// TODO(user): regex literals as well.
	switch n.GetToken() {
	case rhino.Token_STRINGLIT, rhino.Token_STRING_KEY:
		return n.GetString(), true
	case rhino.Token_TEMPLATELIT:
		// Only convert a template literal if all its expressions can be converted.
		var string_ strings.Builder
		for part := range ast.MustTemplateLit(n).Parts() {
			expression := part.Node()
			if sub, ok := part.AsTemplateLitSub(); ok {
				expression = sub.Expr().Node()
			}
			expressionString, ok := nodeUtil_getStringValue(expression)
			if !ok {
				// Cannot convert.
				return "", false
			}
			string_.WriteString(expressionString)
		}
		return string_.String(), true
	case rhino.Token_TEMPLATELIT_STRING:
		return ast.MustTemplateLitString(n).Cooked()
	case rhino.Token_NAME:
		name := n.GetString()
		if "undefined" == name || "Infinity" == name || "NaN" == name {
			return name, true
		}
		return "", false
	case rhino.Token_NEG, rhino.Token_NUMBER:
		if value, ok := nodeUtil_doGetNumberValue(n, true); ok {
			return dtoa.DToA_NumberToString(value), true
		}
		return "", false
	case rhino.Token_BIGINT:
		return n.GetBigInt().String(), true
	case rhino.Token_FALSE:
		return "false", true
	case rhino.Token_TRUE:
		return "true", true
	case rhino.Token_NULL:
		return "null", true
	case rhino.Token_VOID:
		return "undefined", true
	case rhino.Token_NOT:
		child := NodeUtil_GetBooleanValue(ast.MustUnary(n).Operand().Node())
		if child != base.Tri_UNKNOWN {
			if child.ToBoolean(true) { // reversed.
				return "false", true
			}
			return "true", true
		}
		return "", false
	case rhino.Token_ARRAYLIT:
		return nodeUtil_arrayToString(n)
	case rhino.Token_OBJECTLIT:
		if n.HasChildren() {
			return "", false
		}
		return "[object Object]", true
	default:
		return "", false
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
	s, ok := nodeUtil_getArrayElementStringValue(n)
	if !ok {
		return nil
	}
	return &s
}

func nodeUtil_getArrayElementStringValue(n rhino.Node) (string, bool) {
	if NodeUtil_IsNullOrUndefined(n) || n.IsEmpty() {
		return "", true
	}
	return nodeUtil_getStringValue(n)
}

// Java: NodeUtil.arrayToString(Node)
func NodeUtil_ArrayToString(literal rhino.Node) *string {
	s, ok := nodeUtil_arrayToString(literal)
	if !ok {
		return nil
	}
	return &s
}

func nodeUtil_arrayToString(literal rhino.Node) (string, bool) {
	first := literal.GetFirstChild()
	var result strings.Builder
	for n := first; !n.IsNil(); n = n.GetNext() {
		childValue, ok := nodeUtil_getArrayElementStringValue(n)
		if !ok {
			return "", false
		}
		if n != first {
			result.WriteByte(',')
		}
		result.WriteString(childValue)
	}
	return result.String(), true
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
			if len(digits) <= 14 {
				v, err := strconv.ParseInt(digits, radix, 64)
				if err != nil {
					return math.NaN(), true
				}
				return float64(v), true
			}
			v, ok := nodeUtil_bigIntegerDoubleValue(digits, radix)
			if !ok {
				return math.NaN(), true
			}
			return v, true
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

	v, ok := nodeUtil_parseDecimalDouble(s)
	if !ok {
		return math.NaN(), true
	}
	return v, true
}

// nodeUtil_parseDecimalDouble parses a non-empty string restricted to ASCII
// [0-9.eE+-] using Java's Double.parseDouble decimal grammar without panicking.
func nodeUtil_parseDecimalDouble(s string) (float64, bool) {
	i := 0
	negative := false
	if s[0] == '-' || s[0] == '+' {
		negative = s[0] == '-'
		i = 1
	}
	body := s[i:]
	if body == "" {
		return 0, false
	}
	pos, points, digitCount := 0, 0, 0
	for pos < len(body) && ((body[pos] >= '0' && body[pos] <= '9') || body[pos] == '.') {
		if body[pos] == '.' {
			points++
		} else {
			digitCount++
		}
		pos++
	}
	if points > 1 || digitCount == 0 {
		return 0, false
	}
	significandEnd := pos
	var exponent int64
	if pos < len(body) && (body[pos] == 'e' || body[pos] == 'E') {
		pos++
		expNeg := false
		if pos < len(body) && (body[pos] == '+' || body[pos] == '-') {
			expNeg = body[pos] == '-'
			pos++
		}
		expStart := pos
		for pos < len(body) && body[pos] >= '0' && body[pos] <= '9' {
			if exponent < 1e15 {
				exponent = exponent*10 + int64(body[pos]-'0')
			}
			pos++
		}
		if pos == expStart {
			return 0, false
		}
		if expNeg {
			exponent = -exponent
		}
	}
	if pos != len(body) {
		return 0, false
	}
	if exponent > -10000 && exponent < 10000 {
		v, err := strconv.ParseFloat(s, 64)
		if err == nil {
			return v, true
		}
		if ne, ok := err.(*strconv.NumError); ok && ne.Err == strconv.ErrRange {
			return v, true
		}
	}
	sig := make([]byte, 0, digitCount)
	var intDigits int64
	seenPoint := false
	for j := 0; j < significandEnd; j++ {
		switch c := body[j]; {
		case c == '.':
			seenPoint = true
		case c == '0' && len(sig) == 0:
			if seenPoint {
				intDigits--
			}
		default:
			sig = append(sig, c)
			if !seenPoint {
				intDigits++
			}
		}
	}
	var v float64
	if len(sig) > 0 {
		scale10 := intDigits + exponent
		switch {
		case scale10 > 310:
			v = math.Inf(1)
		case scale10 < -400:
			v = 0
		default:
			v, _ = strconv.ParseFloat("0."+string(sig)+"e"+strconv.FormatInt(scale10, 10), 64)
		}
	}
	if negative {
		v = -v
	}
	return v, true
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
func nodeUtil_bigIntegerDoubleValue(digits string, radix int) (float64, bool) {
	x, ok := new(big.Int).SetString(digits, radix)
	if !ok {
		return 0, false
	}
	f, _ := new(big.Float).SetInt(x).Float64()
	return f, true
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
			return nodeUtil_parseBigInteger(jutil.StringSubstringFrom(s, 2), radix)
		}
	}

	// Java: new BigInteger(s), i.e. new BigInteger(s, 10).
	return nodeUtil_parseBigInteger(s, 10)
}

// nodeUtil_parseBigInteger parses s in radix using Java's BigInteger(String, int)
// grammar, returning nil instead of panicking on invalid input.
func nodeUtil_parseBigInteger(s string, radix int) *big.Int {
	units := jutil.StringToUTF16(s)
	start := 0
	negative := false
	if len(units) > 0 && (units[0] == '-' || units[0] == '+') {
		negative = units[0] == '-'
		start = 1
	}
	if start == len(units) {
		return nil
	}
	for _, c := range units[start:] {
		if c == '-' || c == '+' {
			return nil
		}
	}
	for start < len(units) && jutil.CharacterDigit(units[start], radix) == 0 {
		start++
	}
	if start == len(units) {
		return big.NewInt(0)
	}
	digits := units[start:]
	text := make([]byte, len(digits))
	for i, c := range digits {
		d := jutil.CharacterDigit(c, radix)
		if d < 0 {
			return nil
		}
		text[i] = "0123456789abcdefghijklmnopqrstuvwxyz"[d]
	}
	x, ok := new(big.Int).SetString(string(text), radix)
	if !ok {
		return nil
	}
	if negative {
		x.Neg(x)
	}
	return x
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
		for part := range ast.MustTemplateLit(n).Parts() {
			if sub, ok := part.AsTemplateLitSub(); ok {
				if !NodeUtil_IsImmutableValue(sub.Expr().Node()) {
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
