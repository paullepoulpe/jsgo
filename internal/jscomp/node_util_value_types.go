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
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

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
		switch n.GetStringID() {
		case nodeUtilIDUndefined:
			return NodeUtil_ValueType_VOID
		case nodeUtilIDNaN, nodeUtilIDInfinity:
			return NodeUtil_ValueType_NUMBER
		default:
			return NodeUtil_ValueType_UNDETERMINED
		}
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
		id := value.GetStringID()
		// We assume here that programs don't change the value of the keyword
		// undefined to something other than the value undefined.
		return id == nodeUtilIDInfinity || id == nodeUtilIDNaN
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
	return (n.IsName() && n.GetStringID() == nodeUtilIDNaN) ||
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
