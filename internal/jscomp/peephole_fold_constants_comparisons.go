// Port of com/google/javascript/jscomp/PeepholeFoldConstants.java @ closure-compiler v20261021

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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// TryFoldComparison tries to fold comparison nodes, e.g ==
//
// Java: PeepholeFoldConstants.tryFoldComparison(Node, Node, Node)
func (p *PeepholeFoldConstants) TryFoldComparison(n rhino.Node, left rhino.Node, right rhino.Node) rhino.Node {
	result := PeepholeFoldConstants_EvaluateComparison(&p.AbstractPeepholeOptimizationBase, n.GetToken(), left, right)
	if result == base.Tri_UNKNOWN {
		return n
	}

	newNode := NodeUtil_BooleanNode(result.ToBoolean(true))
	p.ReportChangeToEnclosingScope(n)
	n.ReplaceWith(newNode)
	p.MarkFunctionsDeleted(n)

	return newNode
}

// PeepholeFoldConstants_TryAbstractRelationalComparison is
// https://tc39.es/ecma262/#sec-abstract-relational-comparison
//
// Java: PeepholeFoldConstants.tryAbstractRelationalComparison(AbstractPeepholeOptimization, Node, Node, boolean)
func PeepholeFoldConstants_TryAbstractRelationalComparison(
	peepholeOptimization *AbstractPeepholeOptimizationBase,
	left rhino.Node,
	right rhino.Node,
	willNegate bool) base.Tri {
	leftValueType := NodeUtil_GetKnownValueType(left)
	rightValueType := NodeUtil_GetKnownValueType(right)
	// First, check for a string comparison.
	if leftValueType == NodeUtil_ValueType_STRING && rightValueType == NodeUtil_ValueType_STRING {
		lvStr, lvOk := peepholeOptimization.GetSideEffectFreeStringValue(left)
		rvStr, rvOk := peepholeOptimization.GetSideEffectFreeStringValue(right)
		if lvOk && rvOk {
			// In JS, browsers parse \v differently. So do not compare strings if one contains \v.
			if strings.IndexByte(lvStr, '\v') != -1 || strings.IndexByte(rvStr, '\v') != -1 {
				return base.Tri_UNKNOWN
			}
			return base.Tri_ForBoolean(jutil.StringCompareTo(lvStr, rvStr) < 0)
		} else if peepholeFoldConstants_isTypeofOfSameName(left, right) {
			// Special case: `typeof a < typeof a` is always false.
			return base.Tri_FALSE
		}
	}

	// Next, handle comparisons between BigInt and String.
	// Per ECMA-262, comparing a BigInt with a String converts the String to a BigInt via
	// StringToBigInt(string). If StringToBigInt returns undefined (invalid BigInt syntax like
	// "foo", "1e2", "1.5"), all relational comparisons (<, >, <=, >=) return false.
	if leftValueType == NodeUtil_ValueType_BIGINT && rightValueType == NodeUtil_ValueType_STRING {
		if rvStr, ok := peepholeOptimization.GetSideEffectFreeStringValue(right); ok {
			rvBig := NodeUtil_GetStringBigIntValue(rvStr)
			if rvBig == nil {
				return base.Tri_ForBoolean(willNegate)
			}
			lvBig := peepholeOptimization.GetSideEffectFreeBigIntValue(left)
			if lvBig != nil {
				return base.Tri_ForBoolean(lvBig.Cmp(rvBig) < 0)
			}
		}
		return base.Tri_UNKNOWN
	}
	if leftValueType == NodeUtil_ValueType_STRING && rightValueType == NodeUtil_ValueType_BIGINT {
		if lvStr, ok := peepholeOptimization.GetSideEffectFreeStringValue(left); ok {
			lvBig := NodeUtil_GetStringBigIntValue(lvStr)
			if lvBig == nil {
				return base.Tri_ForBoolean(willNegate)
			}
			rvBig := peepholeOptimization.GetSideEffectFreeBigIntValue(right)
			if rvBig != nil {
				return base.Tri_ForBoolean(lvBig.Cmp(rvBig) < 0)
			}
		}
		return base.Tri_UNKNOWN
	}

	// Next, try to evaluate based on the value of the node. Try comparing as BigInts first.
	lvBig := peepholeOptimization.GetSideEffectFreeBigIntValue(left)
	var rvBig *big.Int
	if lvBig != nil {
		rvBig = peepholeOptimization.GetSideEffectFreeBigIntValue(right)
		if rvBig != nil {
			return base.Tri_ForBoolean(lvBig.Cmp(rvBig) < 0)
		}
	}

	// Then, try comparing as Numbers.
	lvNum, lvOk := peepholeOptimization.GetSideEffectFreeNumberValue(left)
	if lvBig != nil || lvOk {
		rvNum, rvOk := peepholeOptimization.GetSideEffectFreeNumberValue(right)
		if lvOk && rvOk {
			if math.IsNaN(lvNum) || math.IsNaN(rvNum) {
				return base.Tri_ForBoolean(willNegate)
			}
			return base.Tri_ForBoolean(lvNum < rvNum)
		}

		// Finally, try comparisons between BigInt and Number.
		if lvBig != nil && rvOk {
			return PeepholeFoldConstants_BigintLessThanDouble(lvBig, rvNum, base.Tri_FALSE, willNegate)
		}
		if lvOk {
			rvBig = peepholeOptimization.GetSideEffectFreeBigIntValue(right)
			if rvBig != nil {
				return PeepholeFoldConstants_BigintLessThanDouble(rvBig, lvNum, base.Tri_TRUE, willNegate)
			}
		}
	}

	// Special case: `x < x` is always false.
	// TODO(moz): If we knew the named value wouldn't be NaN, it would be nice to handle
	// LE and GE. We should use type information if available here.
	if !willNegate && left.IsName() && right.IsName() {
		if left.GetStringID() == right.GetStringID() {
			return base.Tri_FALSE
		}
	}

	return base.Tri_UNKNOWN
}

// peepholeFoldConstants_isTypeofOfSameName reports whether left and right
// are both `typeof a`, for the same name a.
//
// Go-only: the Java spells out this test in both comparisons.
func peepholeFoldConstants_isTypeofOfSameName(left rhino.Node, right rhino.Node) bool {
	leftName, ok := peepholeFoldConstants_typeofName(left)
	if !ok {
		return false
	}
	rightName, ok := peepholeFoldConstants_typeofName(right)
	return ok && leftName == rightName
}

// peepholeFoldConstants_typeofName returns the string ID of a, if n is `typeof a`.
func peepholeFoldConstants_typeofName(n rhino.Node) (uint32, bool) {
	if typeof, ok := ast.AsUnary(n); ok && typeof.Token() == rhino.Token_TYPEOF {
		if name, ok := typeof.Operand().AsName(); ok {
			return name.Node().GetStringID(), true
		}
	}
	return 0, false
}

// Java: PeepholeFoldConstants.bigintLessThanDouble(BigInteger, double, Tri, boolean)
func PeepholeFoldConstants_BigintLessThanDouble(
	bigint *big.Int, number float64, invert base.Tri, willNegate bool) base.Tri {
	// if invert is false, then the number is on the right in tryAbstractRelationalComparison
	// if it's true, then the number is on the left
	if math.IsNaN(number) {
		return base.Tri_ForBoolean(willNegate)
	} else if number == math.Inf(1) {
		return base.Tri_TRUE.Xor(invert)
	} else if number == math.Inf(-1) {
		return base.Tri_FALSE.Xor(invert)
	} else if !base.JSCompDoubles_IsAtLeastIntegerPrecision(number) {
		return base.Tri_UNKNOWN
	}

	// long can hold all values within [-2^53, 2^53]
	var numberAsBigInt big.Int
	numberAsBigInt.SetInt64(jutil.D2L(number))
	negativeMeansBigintSmaller := bigint.Cmp(&numberAsBigInt)
	if negativeMeansBigintSmaller < 0 {
		return base.Tri_TRUE.Xor(invert)
	} else if negativeMeansBigintSmaller > 0 {
		return base.Tri_FALSE.Xor(invert)
	} else if base.JSCompDoubles_IsExactInt64(number) {
		return base.Tri_FALSE // This is the == case, don't invert.
	} else {
		return base.Tri_ForBoolean(base.JSCompDoubles_IsPositive(number)).Xor(invert)
	}
}

// PeepholeFoldConstants_TryAbstractEqualityComparison is
// http://www.ecma-international.org/ecma-262/6.0/#sec-abstract-equality-comparison
//
// Java: PeepholeFoldConstants.tryAbstractEqualityComparison(AbstractPeepholeOptimization, Node, Node)
func PeepholeFoldConstants_TryAbstractEqualityComparison(
	peepholeOptimization *AbstractPeepholeOptimizationBase, left rhino.Node, right rhino.Node) base.Tri {
	// Evaluate based on the general type.
	leftValueType := NodeUtil_GetKnownValueType(left)
	rightValueType := NodeUtil_GetKnownValueType(right)
	if leftValueType != NodeUtil_ValueType_UNDETERMINED && rightValueType != NodeUtil_ValueType_UNDETERMINED {
		// Delegate to strict equality comparison for values of the same type.
		if leftValueType == rightValueType {
			return PeepholeFoldConstants_TryStrictEqualityComparison(peepholeOptimization, left, right)
		}

		if (leftValueType == NodeUtil_ValueType_NULL && rightValueType == NodeUtil_ValueType_VOID) ||
			(leftValueType == NodeUtil_ValueType_VOID && rightValueType == NodeUtil_ValueType_NULL) {
			return base.Tri_TRUE
		}

		if (leftValueType == NodeUtil_ValueType_NUMBER && rightValueType == NodeUtil_ValueType_STRING) ||
			rightValueType == NodeUtil_ValueType_BOOLEAN {
			rv, ok := peepholeOptimization.GetSideEffectFreeNumberValue(right)
			if !ok {
				return base.Tri_UNKNOWN
			}
			return PeepholeFoldConstants_TryAbstractEqualityComparison(
				peepholeOptimization, left, NodeUtil_NumberNode(rv, right))
		}
		if (leftValueType == NodeUtil_ValueType_STRING && rightValueType == NodeUtil_ValueType_NUMBER) ||
			leftValueType == NodeUtil_ValueType_BOOLEAN {
			lv, ok := peepholeOptimization.GetSideEffectFreeNumberValue(left)
			if !ok {
				return base.Tri_UNKNOWN
			}
			return PeepholeFoldConstants_TryAbstractEqualityComparison(
				peepholeOptimization, NodeUtil_NumberNode(lv, left), right)
		}

		if leftValueType == NodeUtil_ValueType_BIGINT && rightValueType == NodeUtil_ValueType_STRING {
			if rvStr, ok := peepholeOptimization.GetSideEffectFreeStringValue(right); ok {
				rvBig := NodeUtil_GetStringBigIntValue(rvStr)
				if rvBig == nil {
					return base.Tri_FALSE
				}
				lvBig := peepholeOptimization.GetSideEffectFreeBigIntValue(left)
				if lvBig != nil {
					return base.Tri_ForBoolean(jutil.BigIntegerEquals(lvBig, rvBig))
				}
				return base.Tri_UNKNOWN
			}
			return base.Tri_UNKNOWN
		}
		if leftValueType == NodeUtil_ValueType_STRING && rightValueType == NodeUtil_ValueType_BIGINT {
			if lvStr, ok := peepholeOptimization.GetSideEffectFreeStringValue(left); ok {
				lvBig := NodeUtil_GetStringBigIntValue(lvStr)
				if lvBig == nil {
					return base.Tri_FALSE
				}
				rvBig := peepholeOptimization.GetSideEffectFreeBigIntValue(right)
				if rvBig != nil {
					return base.Tri_ForBoolean(jutil.BigIntegerEquals(lvBig, rvBig))
				}
				return base.Tri_UNKNOWN
			}
			return base.Tri_UNKNOWN
		}

		if leftValueType == NodeUtil_ValueType_BIGINT || rightValueType == NodeUtil_ValueType_BIGINT {
			lv := peepholeOptimization.GetSideEffectFreeBigIntValue(left)
			rv := peepholeOptimization.GetSideEffectFreeBigIntValue(right)
			if lv != nil && rv != nil {
				return base.Tri_ForBoolean(jutil.BigIntegerEquals(lv, rv))
			}
			if leftValueType == NodeUtil_ValueType_NULL ||
				leftValueType == NodeUtil_ValueType_VOID ||
				rightValueType == NodeUtil_ValueType_NULL ||
				rightValueType == NodeUtil_ValueType_VOID {
				return base.Tri_FALSE
			}
			if leftValueType == NodeUtil_ValueType_STRING && lv == nil {
				if _, ok := peepholeOptimization.GetSideEffectFreeStringValue(left); ok {
					return base.Tri_FALSE
				}
			}
			if rightValueType == NodeUtil_ValueType_STRING && rv == nil {
				if _, ok := peepholeOptimization.GetSideEffectFreeStringValue(right); ok {
					return base.Tri_FALSE
				}
			}
			return base.Tri_UNKNOWN
		}

		if (leftValueType == NodeUtil_ValueType_STRING || leftValueType == NodeUtil_ValueType_NUMBER) &&
			rightValueType == NodeUtil_ValueType_OBJECT {
			return base.Tri_UNKNOWN
		}
		if leftValueType == NodeUtil_ValueType_OBJECT &&
			(rightValueType == NodeUtil_ValueType_STRING || rightValueType == NodeUtil_ValueType_NUMBER) {
			return base.Tri_UNKNOWN
		}

		return base.Tri_FALSE
	}
	// In general, the rest of the cases cannot be folded.
	return base.Tri_UNKNOWN
}

// PeepholeFoldConstants_TryStrictEqualityComparison is
// http://www.ecma-international.org/ecma-262/6.0/#sec-strict-equality-comparison
//
// Java: PeepholeFoldConstants.tryStrictEqualityComparison(AbstractPeepholeOptimization, Node, Node)
func PeepholeFoldConstants_TryStrictEqualityComparison(
	peepholeOptimization *AbstractPeepholeOptimizationBase, left rhino.Node, right rhino.Node) base.Tri {
	// First, try to evaluate based on the general type.
	leftValueType := NodeUtil_GetKnownValueType(left)
	rightValueType := NodeUtil_GetKnownValueType(right)
	if leftValueType != NodeUtil_ValueType_UNDETERMINED && rightValueType != NodeUtil_ValueType_UNDETERMINED {
		// Strict equality can only be true for values of the same type.
		if leftValueType != rightValueType {
			return base.Tri_FALSE
		}
		switch leftValueType {
		case NodeUtil_ValueType_VOID, NodeUtil_ValueType_NULL:
			return base.Tri_TRUE
		case NodeUtil_ValueType_NUMBER:
			if NodeUtil_IsNaN(left) || NodeUtil_IsNaN(right) {
				return base.Tri_FALSE
			}
			lv, lOk := peepholeOptimization.GetSideEffectFreeNumberValue(left)
			rv, rOk := peepholeOptimization.GetSideEffectFreeNumberValue(right)
			if lOk && rOk {
				return base.Tri_ForBoolean(lv == rv)
			}
		case NodeUtil_ValueType_STRING:
			lv, lOk := peepholeOptimization.GetSideEffectFreeStringValue(left)
			rv, rOk := peepholeOptimization.GetSideEffectFreeStringValue(right)
			if lOk && rOk {
				// In JS, browsers parse \v differently. So do not consider strings
				// equal if one contains \v.
				if strings.IndexByte(lv, '\v') != -1 || strings.IndexByte(rv, '\v') != -1 {
					return base.Tri_UNKNOWN
				}
				return base.Tri_ForBoolean(lv == rv)
			} else if peepholeFoldConstants_isTypeofOfSameName(left, right) {
				// Special case, typeof a == typeof a is always true.
				return base.Tri_TRUE
			}
		case NodeUtil_ValueType_BOOLEAN:
			lv := peepholeOptimization.GetSideEffectFreeBooleanValue(left)
			rv := peepholeOptimization.GetSideEffectFreeBooleanValue(right)
			return lv.And(rv).Or(lv.Not().And(rv.Not()))
		case NodeUtil_ValueType_BIGINT:
			lv := peepholeOptimization.GetSideEffectFreeBigIntValue(left)
			rv := peepholeOptimization.GetSideEffectFreeBigIntValue(right)
			if lv != nil && rv != nil {
				return base.Tri_ForBoolean(jutil.BigIntegerEquals(lv, rv))
			}
		default:
			// Symbol and Object cannot be folded in the general case.
			return base.Tri_UNKNOWN
		}
	}

	// Then, try to evaluate based on the value of the node. There's only one special case:
	// Any strict equality comparison against NaN returns false.
	if NodeUtil_IsNaN(left) || NodeUtil_IsNaN(right) {
		return base.Tri_FALSE
	}
	return base.Tri_UNKNOWN
}

// Java: PeepholeFoldConstants.evaluateComparison(AbstractPeepholeOptimization, Token, Node, Node)
func PeepholeFoldConstants_EvaluateComparison(
	peepholeOptimization *AbstractPeepholeOptimizationBase, op rhino.Token, left rhino.Node, right rhino.Node) base.Tri {
	// Don't try to minimize side-effects here.
	if peepholeOptimization.MayHaveSideEffects(left) ||
		peepholeOptimization.MayHaveSideEffects(right) {
		return base.Tri_UNKNOWN
	}

	switch op {
	case rhino.Token_EQ:
		return PeepholeFoldConstants_TryAbstractEqualityComparison(peepholeOptimization, left, right)
	case rhino.Token_NE:
		return PeepholeFoldConstants_TryAbstractEqualityComparison(peepholeOptimization, left, right).Not()
	case rhino.Token_SHEQ:
		return PeepholeFoldConstants_TryStrictEqualityComparison(peepholeOptimization, left, right)
	case rhino.Token_SHNE:
		return PeepholeFoldConstants_TryStrictEqualityComparison(peepholeOptimization, left, right).Not()
	case rhino.Token_LT:
		return PeepholeFoldConstants_TryAbstractRelationalComparison(peepholeOptimization, left, right, false)
	case rhino.Token_GT:
		return PeepholeFoldConstants_TryAbstractRelationalComparison(peepholeOptimization, right, left, false)
	case rhino.Token_LE:
		return PeepholeFoldConstants_TryAbstractRelationalComparison(peepholeOptimization, right, left, true).Not()
	case rhino.Token_GE:
		return PeepholeFoldConstants_TryAbstractRelationalComparison(peepholeOptimization, left, right, true).Not()
	default:
	}
	panic(jutil.NewIllegalStateExceptionMsg("Unexpected operator for comparison"))
}
