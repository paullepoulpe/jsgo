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
	"bytes"
	"math"
	"math/big"
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Java: PeepholeFoldConstants.tryFoldUnaryOperator(Node)
func (p *PeepholeFoldConstants) TryFoldUnaryOperator(n rhino.Node) rhino.Node {
	unary, ok := ast.AsUnary(n)
	jutil.CheckStateMsg1(ok, "%s", n)

	left := unary.Operand().Node()
	parent := n.GetParent()

	leftVal, ok := p.GetSideEffectFreeBooleanValue(left)
	if !ok {
		return n
	}

	switch unary.Token() {
	case rhino.Token_NOT:
		// Don't fold !0 and !1 back to false.
		if p.late && left.IsNumber() {
			numValue := left.GetDouble()
			if numValue == 0 || numValue == 1 {
				return n
			}
		}
		replacementNode := NodeUtil_BooleanNode(!leftVal)
		n.ReplaceWith(replacementNode)
		p.ReportChangeToEnclosingScope(parent)
		return replacementNode
	case rhino.Token_POS:
		if NodeUtil_IsNumericResult(left) {
			// POS does nothing to numeric values.
			n.ReplaceWith(left.Detach())
			p.ReportChangeToEnclosingScope(parent)
			return left
		}
		return n
	case rhino.Token_NEG:
		var result rhino.Node
		if left.IsName() && left.GetStringID() == foldIDNaN {
			result = left.Detach() // "-NaN" is "NaN".
		} else if neg, ok := ast.AsUnary(left); ok && neg.Token() == rhino.Token_NEG {
			leftLeft := neg.Operand().Node()
			if leftLeft.IsBigInt() || leftLeft.IsNumber() {
				result = leftLeft.Detach() // `-(-4)` is `4`
			}
		}

		if !result.IsNil() {
			n.ReplaceWith(result)
			p.ReportChangeToEnclosingScope(parent)
			return result
		}
		return n
	case rhino.Token_BITNOT:
		doubleVal, ok := p.GetSideEffectFreeNumberValue(left)
		if ok {
			if base.JSCompDoubles_IsMathematicalInteger(doubleVal) {
				intVal := base.JSCompDoubles_EcmascriptToInt32(doubleVal)
				notIntValNode := NodeUtil_NumberNode(float64(^intVal), left)
				n.ReplaceWith(notIntValNode)
				p.ReportChangeToEnclosingScope(parent)
				return notIntValNode
			} else {
				p.Report(PeepholeFoldConstants_FRACTIONAL_BITWISE_OPERAND, left)
				return n
			}
		}

		bigintVal := p.GetSideEffectFreeBigIntValue(n)
		if bigintVal != nil {
			bigintNotNode := PeepholeFoldConstants_BigintNode(bigintVal, n)
			n.ReplaceWith(bigintNotNode)
			p.ReportChangeToEnclosingScope(parent)
			return bigintNotNode
		}

		return n
	default:
		return n
	}
}

// Java: PeepholeFoldConstants.tryFoldAdd(Node, Node, Node)
func (p *PeepholeFoldConstants) TryFoldAdd(node rhino.Node, left rhino.Node, right rhino.Node) rhino.Node {
	jutil.CheckArgument(node.IsAdd())

	if NodeUtil_MayBeStringWithUseType(node, p.shouldUseTypes) {
		if left.IsTemplateLit() && right.IsTemplateLit() {
			return p.FoldAddTemplateLiterals(node, left, right)
		} else if left.IsTemplateLit() {
			if rightStr, ok := p.GetSideEffectFreeStringValue(right); ok && p.IsSimpleTemplateStringValue(rightStr) {
				return p.FoldTemplateLiteralAndConstant(node, left, rightStr)
			}
		} else if right.IsTemplateLit() {
			if leftStr, ok := p.GetSideEffectFreeStringValue(left); ok && p.IsSimpleTemplateStringValue(leftStr) {
				return p.FoldConstantAndTemplateLiteral(node, leftStr, right)
			}
		}

		if NodeUtil_IsLiteralValue(left, false) && NodeUtil_IsLiteralValue(right, false) {
			// '6' + 7
			return p.TryFoldAddConstantString(node, left, right)
		} else if left.IsStringLit() && left.GetStringID() == rhino.RhinoStringPool_EmptyID && p.IsStringTyped(right) {
			return p.Replace(node, right.CloneTreeWithTypeExprs(true))
		} else if right.IsStringLit() && right.GetStringID() == rhino.RhinoStringPool_EmptyID && p.IsStringTyped(left) {
			return p.Replace(node, left.CloneTreeWithTypeExprs(true))
		} else {
			// a + 7 or 6 + a
			return p.TryFoldAdjacentLiteralLeaves(node, left, right)
		}
	} else {
		// Try arithmetic add
		return p.TryFoldArithmeticOp(node, left, right)
	}
}

// TryFoldAdjacentLiteralLeaves takes a subtree representing an expression of
// chained addition between expressions and folds adjacent string addition
// when possible and safe.
//
// n is the root node of ADD expression to be optimized; left is the left
// child of n being added; right is the right child of n being added. Returns
// the AST subtree starting from n that has been optimized to collapse the
// rightmost left child leaf node which must be a string literal and the
// leftmost right child leaf node if possible to do so in a type safe way.
//
// Go-only: the walks down chains of '+' test each node with IsAdd and step
// to a child, which the typed view would spell as a conversion per step;
// this keeps the rhino API.
//
// Java: PeepholeFoldConstants.tryFoldAdjacentLiteralLeaves(Node, Node, Node)
func (p *PeepholeFoldConstants) TryFoldAdjacentLiteralLeaves(n rhino.Node, left rhino.Node, right rhino.Node) rhino.Node {

	// Find left child's rightmost leaf
	leftParent := n
	for left.IsAdd() {
		// This had better be in a chain of '+' operations
		leftParent = left
		left = left.GetSecondChild()
	}

	// Find right child's leftmost leaf
	rightParent := n
	for right.IsAdd() {
		// This had better be in a chain of '+' operations
		rightParent = right
		right = right.GetFirstChild()
	}

	// Try to fold if of the form:
	// ... + <STRINGLITERAL> + <LITERAL>
	// ... + <STRINGLITERAL> + (<LITERAL> + ...)
	// aka. when the literal might be collapsible into the string
	if leftParent.IsAdd() &&
		left.IsStringLit() &&
		rightParent.IsAdd() &&
		NodeUtil_IsLiteralValue(right /* includeFunctions= */, false) {
		rightGrandparent := rightParent.GetParent()
		// `rr` is righthand side term that `right` is being added to and is null if right is still
		// the second child of `n` being added to `left` and non-null if it is in a nested add expr
		rr := right.GetNext()
		foldIsTypeSafe :=
			// If right.getNext() is already a string, folding won't disturb any typecasting
			(!rr.IsNil() && NodeUtil_IsStringResult(rr)) ||
				(!rr.IsNil() &&
					// If right.getNext() isn't a string result, right must be a string to be folded
					right.IsStringLit() &&
					// Access the right grandparent safely...
					!rightGrandparent.IsNil() &&
					rightGrandparent.IsAdd() &&
					// The second child of `rightGrandparent` is what `right + rr` is being added to.
					// If it exists, `right` can only be safely removed if `rr` is still treated as a
					// string in the resulting addition.
					NodeUtil_IsStringResult(rightGrandparent.GetSecondChild())) ||
				// Dangling literal term has no typecasting side effects; fold it!
				(rr.IsNil())
		if foldIsTypeSafe {
			rightStr, ok := nodeUtil_getStringValue(right)
			if !ok {
				rightStr = "null"
			}
			result := left.GetString() + rightStr
			// If the right parent is the root, shift the left parent up so as not to overwrite the tree
			// Otherwise, shift the right parent up
			if rightParent.GetSecondChild() == right {
				left.ReplaceWith(rhino.IR_String(result))
				p.Replace(rightParent, rightParent.GetFirstChild().CloneTreeWithTypeExprs(true))
			} else {
				left.ReplaceWith(rhino.IR_String(result))
				p.Replace(rightParent, rightParent.GetSecondChild().CloneTreeWithTypeExprs(true))
			}
		}
	}
	return n
}

// TryFoldAddConstantString tries to fold an ADD node with constant operands.
//
// Java: PeepholeFoldConstants.tryFoldAddConstantString(Node, Node, Node)
func (p *PeepholeFoldConstants) TryFoldAddConstantString(n rhino.Node, left rhino.Node, right rhino.Node) rhino.Node {
	if left.IsStringLit() || right.IsStringLit() || left.IsArrayLit() || right.IsArrayLit() {
		// Add strings.
		leftString, lOk := p.GetSideEffectFreeStringValue(left)
		rightString, rOk := p.GetSideEffectFreeStringValue(right)
		if lOk && rOk {
			newStringNode := rhino.IR_String(leftString + rightString)
			n.ReplaceWith(newStringNode)
			p.ReportChangeToEnclosingScope(newStringNode)
			return newStringNode
		}
	}

	return n
}

// TryFoldArithmeticOp tries to fold arithmetic binary operators.
//
// Java: PeepholeFoldConstants.tryFoldArithmeticOp(Node, Node, Node)
func (p *PeepholeFoldConstants) TryFoldArithmeticOp(n rhino.Node, left rhino.Node, right rhino.Node) rhino.Node {
	result := p.PerformArithmeticOp(n, left, right)
	if !result.IsNil() {
		result.SrcrefTreeIfMissing(n)
		p.ReportChangeToEnclosingScope(n)
		n.ReplaceWith(result)
		return result
	}
	return n
}

// PerformArithmeticOp tries to fold arithmetic binary operators.
//
// Java: PeepholeFoldConstants.performArithmeticOp(Node, Node, Node)
func (p *PeepholeFoldConstants) PerformArithmeticOp(n rhino.Node, left rhino.Node, right rhino.Node) rhino.Node {
	// Unlike other operations, ADD operands are not always converted
	// to Number.
	if n.IsAdd() &&
		(NodeUtil_MayBeStringWithUseType(left, p.shouldUseTypes) ||
			NodeUtil_MayBeStringWithUseType(right, p.shouldUseTypes)) {
		return rhino.NilNode
	}

	if p.IsBigInt(left) && p.IsBigInt(right) {
		return p.PerformBigIntArithmeticOp(n, left, right)
	}

	var result float64

	// TODO(johnlenz): Handle NaN with unknown value. BIT ops convert NaN
	// to zero so this is a little awkward here.

	lval, lOk := p.GetSideEffectFreeNumberValue(left)
	rval, rOk := p.GetSideEffectFreeNumberValue(right)
	// at least one of the two operands must have a value and both must be numeric
	if (!lOk && !rOk) || !p.IsNumeric(left) || !p.IsNumeric(right) {
		return rhino.NilNode
	}
	// handle the operations that have algebraic identities, since we can simplify the tree without
	// actually knowing the value statically.
	switch n.GetToken() {
	case rhino.Token_ADD:
		if lOk && rOk {
			return p.MaybeReplaceBinaryOpWithNumericResult(lval+rval, lval, rval)
		}
		if lOk && lval == 0 {
			return right.CloneTreeWithTypeExprs(true)
		} else if rOk && rval == 0 {
			return left.CloneTreeWithTypeExprs(true)
		}
		return rhino.NilNode
	case rhino.Token_SUB:
		if lOk && rOk {
			return p.MaybeReplaceBinaryOpWithNumericResult(lval-rval, lval, rval)
		}
		if lOk && lval == 0 {
			// 0 - x -> -x
			// NOTE: this optimization has the subtle side effect of changing `0` to `-0` because
			// `0-0 -> 0` but `-0 -> -0`.
			return rhino.IR_Neg(right.CloneTreeWithTypeExprs(true))
		} else if rOk && rval == 0 {
			// x - 0 -> x
			return left.CloneTreeWithTypeExprs(true)
		}
		return rhino.NilNode
	case rhino.Token_MUL:
		if lOk && rOk {
			return p.MaybeReplaceBinaryOpWithNumericResult(lval*rval, lval, rval)
		}
		// NOTE: 0*x != 0 for all x, if x==0, then it is NaN.  So we can't take advantage of that
		// without some kind of non-NaN proof.  So the special cases here only deal with 1*x
		if lOk {
			if lval == 1 {
				return right.CloneTreeWithTypeExprs(true)
			}
		} else {
			if rval == 1 {
				return left.CloneTreeWithTypeExprs(true)
			}
		}
		return rhino.NilNode
	case rhino.Token_DIV:
		if lOk && rOk {
			if rval == 0 {
				return rhino.NilNode
			}
			return p.MaybeReplaceBinaryOpWithNumericResult(lval/rval, lval, rval)
		}
		// NOTE: 0/x != 0 for all x, if x==0, then it is NaN
		if rOk {
			if rval == 1 {
				// x/1->x
				return left.CloneTreeWithTypeExprs(true)
			}
		}
		return rhino.NilNode
	case rhino.Token_EXPONENT:
		if lOk && rOk {
			return p.MaybeReplaceBinaryOpWithNumericResult(
				base.MathPow(lval, rval), lval, rval)
		}
		return rhino.NilNode
	default:
	}
	if !lOk || !rOk {
		return rhino.NilNode
	}
	switch n.GetToken() {
	case rhino.Token_BITAND:
		result = float64(base.JSCompDoubles_EcmascriptToInt32(lval) & base.JSCompDoubles_EcmascriptToInt32(rval))
	case rhino.Token_BITOR:
		result = float64(base.JSCompDoubles_EcmascriptToInt32(lval) | base.JSCompDoubles_EcmascriptToInt32(rval))
	case rhino.Token_BITXOR:
		result = float64(base.JSCompDoubles_EcmascriptToInt32(lval) ^ base.JSCompDoubles_EcmascriptToInt32(rval))
	case rhino.Token_MOD:
		if rval == 0 {
			return rhino.NilNode
		}
		result = jutil.DRem(lval, rval)
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Unexpected arithmetic operator: " + n.GetToken().String()))
	}
	return p.MaybeReplaceBinaryOpWithNumericResult(result, lval, rval)
}

// Java: PeepholeFoldConstants.performBigIntArithmeticOp(Node, Node, Node)
func (p *PeepholeFoldConstants) PerformBigIntArithmeticOp(n rhino.Node, left rhino.Node, right rhino.Node) rhino.Node {
	lVal := p.GetSideEffectFreeBigIntValue(left)
	rVal := p.GetSideEffectFreeBigIntValue(right)
	if lVal != nil && rVal != nil {
		switch n.GetToken() {
		case rhino.Token_ADD:
			return PeepholeFoldConstants_BigintNode(new(big.Int).Add(lVal, rVal), n)
		case rhino.Token_SUB:
			return PeepholeFoldConstants_BigintNode(new(big.Int).Sub(lVal, rVal), n)
		case rhino.Token_MUL:
			return PeepholeFoldConstants_BigintNode(jutil.BigIntegerMultiply(lVal, rVal), n)
		case rhino.Token_DIV:
			if !jutil.BigIntegerEquals(rVal, jutil.BigInteger_ZERO) {
				return PeepholeFoldConstants_BigintNode(new(big.Int).Quo(lVal, rVal), n)
			} else {
				return rhino.NilNode
			}
		case rhino.Token_EXPONENT:
			if pow, ok := peepholeFoldConstants_tryBigIntPow(lVal, rVal); ok {
				return PeepholeFoldConstants_BigintNode(pow, n)
			}
			return rhino.NilNode
		case rhino.Token_MOD:
			return PeepholeFoldConstants_BigintNode(jutil.BigIntegerMod(lVal, rVal), n)
		case rhino.Token_BITAND:
			return PeepholeFoldConstants_BigintNode(new(big.Int).And(lVal, rVal), n)
		case rhino.Token_BITOR:
			return PeepholeFoldConstants_BigintNode(new(big.Int).Or(lVal, rVal), n)
		case rhino.Token_BITXOR:
			return PeepholeFoldConstants_BigintNode(new(big.Int).Xor(lVal, rVal), n)
		default:
			return rhino.NilNode
		}
	}

	// handle the operations that have algebraic identities, since we can simplify the tree without
	// actually knowing the value statically.
	switch n.GetToken() {
	case rhino.Token_ADD:
		if lVal != nil && jutil.BigIntegerEquals(lVal, jutil.BigInteger_ZERO) {
			return right.CloneTreeWithTypeExprs( /* cloneTypeExprs= */ true)
		} else if rVal != nil && jutil.BigIntegerEquals(rVal, jutil.BigInteger_ZERO) {
			return left.CloneTreeWithTypeExprs( /* cloneTypeExprs= */ true)
		}
		return rhino.NilNode
	case rhino.Token_SUB:
		if lVal != nil && jutil.BigIntegerEquals(lVal, jutil.BigInteger_ZERO) {
			// 0n - x -> -x
			return rhino.IR_Neg(right.CloneTreeWithTypeExprs( /* cloneTypeExprs= */ true))
		} else if rVal != nil && jutil.BigIntegerEquals(rVal, jutil.BigInteger_ZERO) {
			// x - 0n -> x
			return left.CloneTreeWithTypeExprs( /* cloneTypeExprs= */ true)
		}
		return rhino.NilNode
	case rhino.Token_MUL:
		if lVal != nil && jutil.BigIntegerEquals(lVal, jutil.BigInteger_ONE) {
			return right.CloneTreeWithTypeExprs( /* cloneTypeExprs= */ true)
		} else if rVal != nil && jutil.BigIntegerEquals(rVal, jutil.BigInteger_ONE) {
			return left.CloneTreeWithTypeExprs( /* cloneTypeExprs= */ true)
		}
		return rhino.NilNode
	case rhino.Token_DIV:
		if rVal != nil && jutil.BigIntegerEquals(rVal, jutil.BigInteger_ONE) {
			return left.CloneTreeWithTypeExprs( /* cloneTypeExprs= */ true)
		}
		return rhino.NilNode
	default:
		return rhino.NilNode
	}
}

// peepholeFoldConstants_tryBigIntPow computes baseVal^expVal with Java's
// BigInteger.pow(rVal.intValueExact()) range checks, returning (nil, false)
// where Java throws ArithmeticException.
func peepholeFoldConstants_tryBigIntPow(baseVal, expVal *big.Int) (*big.Int, bool) {
	if expVal.Sign() < 0 || !expVal.IsInt64() {
		return nil, false
	}
	exp64 := expVal.Int64()
	if exp64 > math.MaxInt32 {
		return nil, false
	}
	exponent := int32(exp64)
	if exponent == 0 {
		return big.NewInt(1), true
	}
	if baseVal.Sign() == 0 {
		return big.NewInt(0), true
	}
	negative := baseVal.Sign() < 0 && exponent&1 == 1
	magnitude := new(big.Int).Abs(baseVal)
	var result *big.Int
	if twos := magnitude.TrailingZeroBits(); uint(magnitude.BitLen()) == twos+1 {
		bits := int64(twos) * int64(exponent)
		if bits+1 > math.MaxInt32 {
			return nil, false
		}
		result = new(big.Int).Lsh(big.NewInt(1), uint(bits))
	} else {
		if int64(magnitude.BitLen())*int64(exponent)/32 > 1<<26 {
			return nil, false
		}
		result = new(big.Int).Exp(magnitude, big.NewInt(int64(exponent)), nil)
		if result.BitLen() > math.MaxInt32 {
			return nil, false
		}
	}
	if negative {
		result.Neg(result)
	}
	return result, true
}

// Java: PeepholeFoldConstants.isNumeric(Node)
func (p *PeepholeFoldConstants) IsNumeric(n rhino.Node) bool {
	if NodeUtil_IsNumericResult(n) {
		return true
	}
	if p.shouldUseTypes {
		return n.GetColor() != nil && n.GetColor().Equals(colors.StandardColors_NUMBER)
	}
	return false
}

// Java: PeepholeFoldConstants.isBigInt(Node)
func (p *PeepholeFoldConstants) IsBigInt(n rhino.Node) bool {
	if NodeUtil_IsBigIntResult(n) {
		return true
	}
	if p.shouldUseTypes {
		return n.GetColor() != nil && n.GetColor().Equals(colors.StandardColors_BIGINT)
	}
	return false
}

// peepholeIntStringLen returns len(strconv.Itoa(v)) without allocating.
func peepholeIntStringLen(v int) int {
	n := 0
	if v < 0 {
		n = 1
		v = -v
	}
	switch {
	case v < 10:
		return n + 1
	case v < 100:
		return n + 2
	case v < 1000:
		return n + 3
	case v < 10000:
		return n + 4
	case v < 100000:
		return n + 5
	case v < 1000000:
		return n + 6
	default:
		return n + 7
	}
}

// peepholeDoubleToStringLen returns len(jutil.DoubleToString(d)) without
// allocating temporary strings on the hot arithmetic-folding path.
func peepholeDoubleToStringLen(d float64) int {
	switch {
	case math.IsNaN(d):
		return 3
	case math.IsInf(d, 1):
		return 8
	case math.IsInf(d, -1):
		return 9
	case d == 0:
		if math.Signbit(d) {
			return 4
		}
		return 3
	}
	if iv := int64(d); float64(iv) == d && iv > -10000000 && iv < 10000000 {
		return peepholeIntStringLen(int(iv)) + 2 // ".0"
	}
	a := math.Abs(d)
	var buf [32]byte
	b := strconv.AppendFloat(buf[:0], a, 'e', -1, 64)
	ePos := bytes.IndexByte(b, 'e')
	digitsLen := 1
	if ePos > 1 {
		digitsLen = ePos - 1
	}
	if digitsLen == 1 {
		return len(jutil.DoubleToString(d))
	}
	expSign := 1
	i := ePos + 1
	if b[i] == '-' {
		expSign = -1
		i++
	} else if b[i] == '+' {
		i++
	}
	exp := 0
	for ; i < len(b); i++ {
		exp = exp*10 + int(b[i]-'0')
	}
	exp *= expSign
	n := 0
	if d < 0 {
		n = 1
	}
	if a >= 1e-3 && a < 1e7 {
		if exp >= 0 {
			intLen := exp + 1
			if digitsLen <= intLen {
				return n + intLen + 2
			}
			return n + digitsLen + 1
		}
		return n + 2 + (-exp - 1) + digitsLen
	}
	return n + digitsLen + 2 + peepholeIntStringLen(exp)
}

// Java: PeepholeFoldConstants.maybeReplaceBinaryOpWithNumericResult(double, double, double)
func (p *PeepholeFoldConstants) MaybeReplaceBinaryOpWithNumericResult(
	result float64, lval float64, rval float64) rhino.Node {
	if math.IsNaN(result) || math.IsInf(result, 0) {
		return NodeUtil_NumberNode(result, rhino.NilNode)
	}
	// Do not try to fold arithmetic for numbers > 2^53. After that
	// point, fixed-point math starts to break down and become inaccurate.
	if math.Abs(result) > peepholeFoldConstants_MAX_FOLD_NUMBER {
		return rhino.NilNode
	}
	// TODO(johnlenz): consider removing the result length check.
	// length of the left and right value plus 1 byte for the operator.
	// Double.toString is always ASCII.
	if peepholeDoubleToStringLen(result) <=
		peepholeDoubleToStringLen(lval)+peepholeDoubleToStringLen(rval)+1 {
		return NodeUtil_NumberNode(result, rhino.NilNode)
	}
	return rhino.NilNode
}

// TryFoldLeftChildOp folds expressions such as [foo() && 10 && 20] that
// generate parse trees where no node has two const children ((foo() && 10)
// && 20), so performArithmeticOp() won't fold it: tryFoldLeftChildOp() will.
//
// Specifically, this folds associative expressions where:
//   - The left child is also an associative expression of the same type.
//   - The right child is a BIGINT or NUMBER constant.
//   - The left child's right child is a BIGINT or NUMBER constant.
//
// Java: PeepholeFoldConstants.tryFoldLeftChildOp(Node, Node, Node)
func (p *PeepholeFoldConstants) TryFoldLeftChildOp(n rhino.Node, left rhino.Node, right rhino.Node) rhino.Node {
	opType := n.GetToken()
	jutil.CheckState(NodeUtil_IsAssociative(opType) && NodeUtil_IsCommutative(opType))

	// Use getNumberValue to handle constants like "NaN" and "Infinity"
	// other values are converted to numbers elsewhere.
	_, rightNumOk := p.GetSideEffectFreeNumberValue(right)
	rightBigInt := p.GetSideEffectFreeBigIntValue(right)
	// (Go: opType is a binary operator, so a left of the same type is one
	// too, with two operands: the Java checks it has two children.)
	leftOp, leftIsOp := ast.AsBinary(left)
	if (rightNumOk || rightBigInt != nil) && leftIsOp && leftOp.Token() == opType {
		ll, lr := leftOp.Left().Node(), leftOp.Right().Node()

		valueToCombine := ll
		replacement := p.PerformArithmeticOp(n, valueToCombine, right)
		if replacement.IsNil() {
			valueToCombine = lr
			replacement = p.PerformArithmeticOp(n, valueToCombine, right)
		}
		if !replacement.IsNil() {
			// Remove the child that has been combined
			valueToCombine.Detach()
			// Replace the left op with the remaining child.
			left.ReplaceWith(left.RemoveFirstChild())
			// New "-Infinity" node need location info explicitly
			// added.
			replacement.SrcrefTreeIfMissing(right)
			right.ReplaceWith(replacement)
			p.ReportChangeToEnclosingScope(n)
		}
	}

	return n
}

// TryFoldShift tries to fold shift operations.
//
// Java: PeepholeFoldConstants.tryFoldShift(Node, Node, Node)
func (p *PeepholeFoldConstants) TryFoldShift(n rhino.Node, left rhino.Node, right rhino.Node) rhino.Node {
	leftVal, lOk := p.GetSideEffectFreeNumberValue(left)
	rightVal, rOk := p.GetSideEffectFreeNumberValue(right)

	if !p.isReasonableDoubleVal(leftVal, lOk) || !p.isReasonableDoubleVal(rightVal, rOk) {
		return n
	}
	if !base.JSCompDoubles_IsMathematicalInteger(leftVal) {
		p.Report(PeepholeFoldConstants_FRACTIONAL_BITWISE_OPERAND, left)
		return n
	}
	if !base.JSCompDoubles_IsMathematicalInteger(rightVal) {
		p.Report(PeepholeFoldConstants_FRACTIONAL_BITWISE_OPERAND, right)
		return n
	}

	// only the lower 5 bits are used when shifting, so don't do anything
	// if the shift amount is outside [0,32)
	if !(0 <= rightVal && rightVal < 32) {
		return n
	}

	rvalInt := jutil.D2I(rightVal)
	bits := base.JSCompDoubles_EcmascriptToInt32(leftVal)

	var result float64
	switch n.GetToken() {
	case rhino.Token_LSH:
		result = float64(jutil.Shl32(bits, int64(rvalInt)))
	case rhino.Token_RSH:
		result = float64(jutil.Shr32(bits, int64(rvalInt)))
	case rhino.Token_URSH:
		// JavaScript always treats the result of >>> as unsigned.
		// We must force Java to do the same here.
		result = float64(0xffffffff & int64(jutil.Ushr32(bits, int64(rvalInt))))
	default:
		panic(jutil.NewAssertionError("Unknown shift operator: " + n.GetToken().String()))
	}

	newNumber := NodeUtil_NumberNode(result, n)
	p.ReportChangeToEnclosingScope(n)
	n.ReplaceWith(newNumber)

	return newNumber
}

// Java: PeepholeFoldConstants.bigintNode(BigInteger, Node)
func PeepholeFoldConstants_BigintNode(val *big.Int, srcref rhino.Node) rhino.Node {
	var tree rhino.Node
	if val.Sign() < 0 {
		tree = rhino.IR_Neg(rhino.IR_Bigint(new(big.Int).Neg(val)))
	} else {
		tree = rhino.IR_Bigint(val)
	}
	return tree.SrcrefTree(srcref)
}
