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

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: this pass uses the typed AST view (internal/ast) where it reads
// shapes: the operands of the operators it folds (op.Left(), op.Right(),
// assign.Value(), unary.Operand(), get.Index()...), the callee and
// arguments of calls, the members of literals. The Java reads the first two
// children of any node and checks there are two; here each kind is named,
// so the checks of child counts and of missing operands are gone. The
// numeric and string value logic, the template literal joins (raw and
// cooked strings), the chains of '+' and the edits use the rhino API, in
// the Java's order.

// PeepholeFoldConstants is a peephole optimization to fold constants (e.g.
// x + 1 + 7 --> x + 8).
type PeepholeFoldConstants struct {
	AbstractPeepholeOptimizationBase
	late           bool
	shouldUseTypes bool
}

// TODO(johnlenz): optimizations should not be emiting errors. Move these to
// a check pass.
var PeepholeFoldConstants_INVALID_GETELEM_INDEX_ERROR = DiagnosticType_Warning("JSC_INVALID_GETELEM_INDEX_ERROR", "Array index not integer: {0}")

var PeepholeFoldConstants_FRACTIONAL_BITWISE_OPERAND = DiagnosticType_Warning("JSC_FRACTIONAL_BITWISE_OPERAND", "Fractional bitwise operand: {0}")

var peepholeFoldConstants_MAX_FOLD_NUMBER = jutil.MathPow(2, 53)

// NewPeepholeFoldConstants creates the optimization.
//
// late: When late is false, this mean we are currently running before most
// of the other optimizations. In this case we would avoid optimizations that
// would make the code harder to analyze. When this is true, we would do
// anything to minimize for size.
//
// Java: PeepholeFoldConstants(boolean, boolean)
func NewPeepholeFoldConstants(late bool, shouldUseTypes bool) *PeepholeFoldConstants {
	return &PeepholeFoldConstants{
		late:           late,
		shouldUseTypes: shouldUseTypes,
	}
}

// Java: PeepholeFoldConstants.optimizeSubtree(Node)
func (p *PeepholeFoldConstants) OptimizeSubtree(subtree rhino.Node) rhino.Node {
	switch subtree.GetToken() {
	case rhino.Token_OPTCHAIN_CALL, rhino.Token_CALL:
		return p.TryFoldUselessObjectDotDefinePropertiesCall(subtree)
	case rhino.Token_NEW:
		return p.TryFoldCtorCall(subtree)
	case rhino.Token_TYPEOF:
		return p.TryFoldTypeof(subtree)
	case rhino.Token_ITER_SPREAD:
		return p.TryFoldSpread(subtree)
	case rhino.Token_ARRAYLIT, rhino.Token_OBJECTLIT:
		return p.TryFlattenArrayOrObjectLit(subtree)
	case rhino.Token_NOT, rhino.Token_POS, rhino.Token_NEG, rhino.Token_BITNOT:
		p.TryReduceOperandsForOp(subtree)
		return p.TryFoldUnaryOperator(subtree)
	case rhino.Token_VOID:
		return p.TryReduceVoid(subtree)
	case rhino.Token_OPTCHAIN_GETPROP, rhino.Token_GETPROP:
		return p.TryFoldGetProp(subtree)
	case rhino.Token_TEMPLATELIT:
		return p.TryFoldTemplateLiteralSubstitutions(subtree)
	default:
		p.TryReduceOperandsForOp(subtree)
		return p.TryFoldBinaryOperator(subtree)
	}
}

// Java: PeepholeFoldConstants.tryFoldBinaryOperator(Node)
//
// Go-only: the Java reads the first two children of any node and goes on
// only if there are two ("node is truly a binary operator"). Here the kinds
// with two operands are named: element accesses and assignments, whose
// helpers read their operands by role, and the binary operators, read
// below. Every one of them always has both operands.
func (p *PeepholeFoldConstants) TryFoldBinaryOperator(subtree rhino.Node) rhino.Node {
	token := subtree.GetToken()
	switch token {
	case rhino.Token_GETELEM, rhino.Token_OPTCHAIN_GETELEM:
		return p.TryFoldGetElem(subtree)
	case rhino.Token_ASSIGN:
		return p.TryFoldAssign(subtree)
	case rhino.Token_ASSIGN_BITOR,
		rhino.Token_ASSIGN_BITXOR,
		rhino.Token_ASSIGN_BITAND,
		rhino.Token_ASSIGN_LSH,
		rhino.Token_ASSIGN_RSH,
		rhino.Token_ASSIGN_URSH,
		rhino.Token_ASSIGN_ADD,
		rhino.Token_ASSIGN_SUB,
		rhino.Token_ASSIGN_MUL,
		rhino.Token_ASSIGN_DIV,
		rhino.Token_ASSIGN_MOD,
		rhino.Token_ASSIGN_EXPONENT:
		return p.TryUnfoldAssignOp(subtree)
	}

	op, ok := ast.AsBinary(subtree)
	if !ok {
		return subtree
	}
	left, right := op.Left().Node(), op.Right().Node()
	switch token {
	case rhino.Token_INSTANCEOF:
		return p.TryFoldInstanceof(subtree, left, right)
	case rhino.Token_AND, rhino.Token_OR:
		return p.TryFoldAndOr(subtree, left, right)
	case rhino.Token_COALESCE:
		return p.TryFoldCoalesce(subtree, left, right)
	case rhino.Token_LSH, rhino.Token_RSH, rhino.Token_URSH:
		return p.TryFoldShift(subtree, left, right)
	case rhino.Token_ADD:
		return p.TryFoldAdd(subtree, left, right)
	case rhino.Token_SUB, rhino.Token_DIV, rhino.Token_MOD, rhino.Token_MUL, rhino.Token_EXPONENT:
		return p.TryFoldArithmeticOp(subtree, left, right)
	case rhino.Token_BITAND, rhino.Token_BITOR, rhino.Token_BITXOR:
		result := p.TryFoldArithmeticOp(subtree, left, right)
		if result != subtree {
			return result
		}
		return p.TryFoldLeftChildOp(subtree, left, right)
	case rhino.Token_LT, rhino.Token_GT, rhino.Token_LE, rhino.Token_GE, rhino.Token_EQ, rhino.Token_NE, rhino.Token_SHEQ, rhino.Token_SHNE:
		return p.TryFoldComparison(subtree, left, right)
	default:
		return subtree
	}
}

// Java: PeepholeFoldConstants.tryReduceVoid(Node)
func (p *PeepholeFoldConstants) TryReduceVoid(n rhino.Node) rhino.Node {
	void, ok := ast.AsUnary(n)
	jutil.CheckArgument(ok && void.Token() == rhino.Token_VOID)

	child := void.Operand()
	if num, ok := child.AsNumber(); ok && num.Value() == 0.0 {
		return n // already `void 0`
	}
	if !p.MayHaveSideEffects(n) {
		child.Node().ReplaceWith(rhino.IR_Number(0))
		p.ReportChangeToEnclosingScope(n)
	}
	return n
}

// Java: PeepholeFoldConstants.tryReduceOperandsForOp(Node)
func (p *PeepholeFoldConstants) TryReduceOperandsForOp(n rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_ADD:
		add := ast.MustBinary(n)
		if !NodeUtil_MayBeStringWithUseType(add.Left().Node(), p.shouldUseTypes) &&
			!NodeUtil_MayBeStringWithUseType(add.Right().Node(), p.shouldUseTypes) {
			p.TryConvertOperandsToNumber(n)
		}
	case rhino.Token_ASSIGN_BITOR,
		rhino.Token_ASSIGN_BITXOR,
		rhino.Token_ASSIGN_BITAND,
		rhino.Token_ASSIGN_LSH,
		rhino.Token_ASSIGN_RSH,
		rhino.Token_ASSIGN_URSH,
		rhino.Token_ASSIGN_SUB,
		rhino.Token_ASSIGN_MUL,
		rhino.Token_ASSIGN_MOD,
		rhino.Token_ASSIGN_DIV:
		assignOp := ast.MustAssignOp(n)
		// TODO(johnlenz): convert these to integers.
		p.TryConvertToNumber(assignOp.Value().Node())
	case rhino.Token_BITNOT, rhino.Token_BITOR, rhino.Token_BITXOR, rhino.Token_BITAND, rhino.Token_LSH, rhino.Token_RSH, rhino.Token_URSH, rhino.Token_SUB, rhino.Token_MUL, rhino.Token_MOD, rhino.Token_DIV, rhino.Token_POS, rhino.Token_NEG, rhino.Token_EXPONENT:
		p.TryConvertOperandsToNumber(n)
	default:
	}
}

// Go-only: n is a unary or a binary operator, whose operands this treats
// alike, so it uses the rhino API.
//
// Java: PeepholeFoldConstants.tryConvertOperandsToNumber(Node)
func (p *PeepholeFoldConstants) TryConvertOperandsToNumber(n rhino.Node) {
	var next rhino.Node
	for c := n.GetFirstChild(); !c.IsNil(); c = next {
		next = c.GetNext()
		p.TryConvertToNumber(c)
	}
}

// Java: PeepholeFoldConstants.tryConvertToNumber(Node)
func (p *PeepholeFoldConstants) TryConvertToNumber(n rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_NUMBER:
		// Nothing to do
		return
	case rhino.Token_AND, rhino.Token_OR, rhino.Token_COMMA, rhino.Token_COALESCE:
		op := ast.MustBinary(n)
		p.TryConvertToNumber(op.Right().Node())
		return
	case rhino.Token_HOOK:
		hook := ast.MustHook(n)
		p.TryConvertToNumber(hook.Then().Node())
		p.TryConvertToNumber(hook.Else().Node())
		return
	case rhino.Token_NAME:
		if !NodeUtil_IsUndefined(n) {
			return
		}
	default:
	}

	value, ok := p.GetSideEffectFreeNumberValue(n)
	if !ok {
		return
	}

	replacement := NodeUtil_NumberNode(value, n)
	if replacement.IsEquivalentTo(n) {
		return
	}

	n.ReplaceWith(replacement)
	p.ReportChangeToEnclosingScope(replacement)
}

// TryFoldTypeof folds 'typeof(foo)' if foo is a literal, e.g. typeof("bar")
// --> "string" typeof(6) --> "number"
//
// Java: PeepholeFoldConstants.tryFoldTypeof(Node)
func (p *PeepholeFoldConstants) TryFoldTypeof(originalTypeofNode rhino.Node) rhino.Node {
	typeof, ok := ast.AsUnary(originalTypeofNode)
	jutil.CheckArgument(ok && typeof.Token() == rhino.Token_TYPEOF)

	argumentNode := typeof.Operand()
	if !NodeUtil_IsLiteralValue(argumentNode.Node(), true) {
		return originalTypeofNode
	}

	typeNameString := "" // (Go: "" for the Java's null)

	switch argumentNode.Token() {
	case rhino.Token_FUNCTION:
		typeNameString = "function"
	case rhino.Token_STRINGLIT:
		typeNameString = "string"
	case rhino.Token_NUMBER:
		typeNameString = "number"
	case rhino.Token_TRUE, rhino.Token_FALSE:
		typeNameString = "boolean"
	case rhino.Token_NULL, rhino.Token_OBJECTLIT, rhino.Token_ARRAYLIT:
		typeNameString = "object"
	case rhino.Token_VOID:
		typeNameString = "undefined"
	case rhino.Token_NAME:
		// We assume here that programs don't change the value of the
		// keyword undefined to something other than the value undefined.
		if name, ok := argumentNode.AsName(); ok && name.Ident() == "undefined" {
			typeNameString = "undefined"
		}
	default:
	}

	if typeNameString != "" {
		newNode := rhino.IR_String(typeNameString)
		p.ReportChangeToEnclosingScope(originalTypeofNode)
		originalTypeofNode.ReplaceWith(newNode)
		p.MarkFunctionsDeleted(originalTypeofNode)

		return newNode
	}

	return originalTypeofNode
}

// Java: PeepholeFoldConstants.tryFoldUnaryOperator(Node)
func (p *PeepholeFoldConstants) TryFoldUnaryOperator(n rhino.Node) rhino.Node {
	unary, ok := ast.AsUnary(n)
	jutil.CheckStateMsg1(ok, "%s", n)

	left := unary.Operand().Node()
	parent := n.GetParent()

	leftVal := p.GetSideEffectFreeBooleanValue(left)
	if leftVal == base.Tri_UNKNOWN {
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
		replacementNode := NodeUtil_BooleanNode(!leftVal.ToBoolean(true))
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
		if left.IsName() && left.GetString() == "NaN" {
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

func (p *PeepholeFoldConstants) isReasonableDoubleVal(x float64, ok bool) bool {
	return ok && !math.IsInf(x, 0) && !math.IsNaN(x)
}

// TryFoldInstanceof tries to fold `left instanceof right` into true or
// false.
//
// Java: PeepholeFoldConstants.tryFoldInstanceof(Node, Node, Node)
func (p *PeepholeFoldConstants) TryFoldInstanceof(n rhino.Node, left rhino.Node, right rhino.Node) rhino.Node {
	jutil.CheckArgument(n.IsInstanceOf())

	// TODO(johnlenz) Use type information if available to fold
	// instanceof.
	if NodeUtil_IsLiteralValue(left, true) && !p.MayHaveSideEffects(right) {

		var replacementNode rhino.Node

		if NodeUtil_IsImmutableValue(left) {
			// Non-object types are never instances.
			replacementNode = rhino.IR_FalseNode()
		} else if left.IsTemplateLit() {
			// Template literals always evaluate to strings. If they do not meet the "immutable value"
			// condition above, verify their substitutions have no side effects before folding.
			if !p.MayHaveSideEffects(left) {
				replacementNode = rhino.IR_FalseNode()
			}
		} else if right.IsName() && right.GetString() == "Object" {
			replacementNode = rhino.IR_TrueNode()
		}

		if !replacementNode.IsNil() {
			n.ReplaceWith(replacementNode)
			p.ReportChangeToEnclosingScope(replacementNode)
			p.MarkFunctionsDeleted(n)
			return replacementNode
		}
	}

	return n
}

// Java: PeepholeFoldConstants.tryFoldAssign(Node, Node, Node)
func (p *PeepholeFoldConstants) TryFoldAssign(n rhino.Node) rhino.Node {
	assign, ok := ast.AsAssign(n)
	jutil.CheckArgument(ok)

	if !p.late {
		return n
	}

	// Tries to convert x = x + y -> x += y;
	// Go-only: the operators below are all binary operators, so the RHS is
	// read as one; the Java checks after the switch that it has two
	// children.
	left := assign.Target().Node()
	right, ok := assign.Value().AsBinary()
	if !ok {
		return n
	}

	// Only certain RHS operators are supported.
	operator := right.Token()
	var newType rhino.Token
	switch operator {
	case rhino.Token_ADD:
		newType = rhino.Token_ASSIGN_ADD
	case rhino.Token_BITAND:
		newType = rhino.Token_ASSIGN_BITAND
	case rhino.Token_BITOR:
		newType = rhino.Token_ASSIGN_BITOR
	case rhino.Token_BITXOR:
		newType = rhino.Token_ASSIGN_BITXOR
	case rhino.Token_DIV:
		newType = rhino.Token_ASSIGN_DIV
	case rhino.Token_LSH:
		newType = rhino.Token_ASSIGN_LSH
	case rhino.Token_MOD:
		newType = rhino.Token_ASSIGN_MOD
	case rhino.Token_MUL:
		newType = rhino.Token_ASSIGN_MUL
	case rhino.Token_RSH:
		newType = rhino.Token_ASSIGN_RSH
	case rhino.Token_SUB:
		newType = rhino.Token_ASSIGN_SUB
	case rhino.Token_URSH:
		newType = rhino.Token_ASSIGN_URSH
	case rhino.Token_EXPONENT:
		newType = rhino.Token_ASSIGN_EXPONENT
	default:
		return n
	}

	if p.MayHaveSideEffects(left) {
		return n
	}

	// We can fold any 'leaf' that
	// 1. is equal to the left side of the assignment
	// 2. is either the first child of the operator or any child of the operator if the operator
	//    is commutative
	// 3. Also consider recursively discoverable children of the operator if the operator is
	//    associative.
	//
	// For example, consider x = x + y + z;
	// because + is left associative, the tree looks like `(x + y) + z` so we need to look deeper
	// into the tree to find x.
	// Right now we only recurse one level deep which allows use to handle x = x + y + z but not
	// x = x + y + z + w.

	rightLeft, rightRight := right.Left().Node(), right.Right().Node()
	var newRight rhino.Node
	if p.AreNodesEqualForInlining(left, rightLeft) {
		newRight = rightRight
	} else if NodeUtil_IsCommutative(operator) &&
		!p.MayHaveSideEffects(rightLeft) &&
		p.AreNodesEqualForInlining(left, rightRight) {
		newRight = rightLeft
	} else if inner, ok := ast.AsBinary(rightLeft); ok &&
		NodeUtil_IsAssociative(operator) &&
		inner.Token() == operator &&
		p.AreNodesEqualForInlining(left, inner.Left().Node()) {
		// x = (x op y) op z: the RHS becomes y op z.
		newLeft := inner.Right().Node().Detach()
		right.Node().RemoveFirstChild()
		right.Node().AddChildToFront(newLeft)
		newRight = right.Node()

	} else {
		return n
	}

	newNode := rhino.NewNodeWith2Children(newType, left.Detach(), newRight.Detach())
	n.ReplaceWith(newNode)

	p.ReportChangeToEnclosingScope(newNode)

	return newNode
}

// Java: PeepholeFoldConstants.tryUnfoldAssignOp(Node, Node, Node)
//
// Go-only: the Java checks that n has two children; an AssignOp always has
// its target and value.
func (p *PeepholeFoldConstants) TryUnfoldAssignOp(n rhino.Node) rhino.Node {
	assignOp, ok := ast.AsAssignOp(n)
	jutil.CheckArgument(ok)

	if p.late {
		return n
	}

	left, right := assignOp.Target().Node(), assignOp.Value().Node()
	if p.MayHaveSideEffects(left) {
		return n
	}

	// Tries to convert x += y -> x = x + y;
	op := NodeUtil_GetOpFromAssignmentOp(n)
	replacement :=
		rhino.IR_Assign(left.Detach(), rhino.NewNodeWith2Children(op, left.CloneTree(), right.Detach()).Srcref(n))
	n.ReplaceWith(replacement)
	p.ReportChangeToEnclosingScope(replacement)
	return replacement
}

// TryFoldAndOr tries to fold a AND/OR node.
//
// Java: PeepholeFoldConstants.tryFoldAndOr(Node, Node, Node)
func (p *PeepholeFoldConstants) TryFoldAndOr(n rhino.Node, left rhino.Node, right rhino.Node) rhino.Node {
	parent := n.GetParent()

	var result rhino.Node
	var dropped rhino.Node

	type_ := n.GetToken()

	leftVal := NodeUtil_GetBooleanValue(left)

	if leftVal != base.Tri_UNKNOWN {
		lval := leftVal.ToBoolean(true)

		// (TRUE || x) => TRUE (also, (3 || x) => 3)
		// (FALSE && x) => FALSE
		var keepLeft bool
		if lval {
			keepLeft = type_ == rhino.Token_OR
		} else {
			keepLeft = type_ == rhino.Token_AND
		}
		if keepLeft {
			result = left
			dropped = right
		} else if !p.MayHaveSideEffects(left) {
			// (FALSE || x) => x
			// (TRUE && x) => x
			result = right
			dropped = left
		} else {
			// Left side may have side effects, but we know its boolean value.
			// e.g. true_with_sideeffects || foo() => true_with_sideeffects, foo()
			// or: false_with_sideeffects && foo() => false_with_sideeffects, foo()
			// This, combined with PeepholeRemoveDeadCode, helps reduce expressions
			// like "x() || false || z()".
			n.DetachChildren()
			result = rhino.IR_Comma(left, right)
			dropped = rhino.NilNode
		}
	} else if parentOp, ok := ast.AsBinary(parent); ok && parentOp.Token() == type_ && parentOp.Left().Is(n) {
		rightValue := NodeUtil_GetBooleanValue(right)
		if !p.MayHaveSideEffects(right) {
			if (rightValue == base.Tri_FALSE && type_ == rhino.Token_OR) ||
				(rightValue == base.Tri_TRUE && type_ == rhino.Token_AND) {
				result = left
				dropped = right
			}
		}
	}

	// Note: Right hand side folding is handled by
	// PeepholeMinimizeConditions#tryMinimizeCondition

	if !result.IsNil() {
		// Fold it!
		n.DetachChildren()
		n.ReplaceWith(result)
		p.ReportChangeToEnclosingScope(result)
		if !dropped.IsNil() {
			p.MarkFunctionsDeleted(dropped)
		}
		return result
	} else {
		return n
	}
}

// TryFoldCoalesce tries to fold a COALESCE node.
//
// Java: PeepholeFoldConstants.tryFoldCoalesce(Node, Node, Node)
func (p *PeepholeFoldConstants) TryFoldCoalesce(n rhino.Node, left rhino.Node, right rhino.Node) rhino.Node {

	var result rhino.Node

	leftVal := NodeUtil_GetKnownValueType(left)

	switch leftVal {
	case NodeUtil_ValueType_NULL, NodeUtil_ValueType_VOID:
		// nullish condition => this expression evaluates to the right side.
		if !p.MayHaveSideEffects(left) {
			result = right
			p.MarkFunctionsDeleted(left)
		} else {
			// e.g. `(a(), null) ?? 1` => `(a(), null, 1)`
			n.DetachChildren()
			result = rhino.IR_Comma(left, right)
		}
	case NodeUtil_ValueType_NUMBER, NodeUtil_ValueType_BIGINT, NodeUtil_ValueType_STRING, NodeUtil_ValueType_BOOLEAN, NodeUtil_ValueType_OBJECT:
		// non-nullish condition => this expression evaluates to the left side.
		result = left
		p.MarkFunctionsDeleted(right)
	case NodeUtil_ValueType_UNDETERMINED:
	}

	if !result.IsNil() {
		// Fold!
		n.DetachChildren()
		n.ReplaceWith(result)
		p.ReportChangeToEnclosingScope(result)
		return result
	} else {
		return n
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
				jutil.MathPow(lval, rval), lval, rval)
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

// Java: PeepholeFoldConstants.maybeReplaceBinaryOpWithNumericResult(double, double, double)
func (p *PeepholeFoldConstants) MaybeReplaceBinaryOpWithNumericResult(
	result float64, lval float64, rval float64) rhino.Node {
	// TODO(johnlenz): consider removing the result length check.
	// length of the left and right value plus 1 byte for the operator.
	// Double.toString is always ASCII.
	if (len(jutil.DoubleToString(result)) <=
		len(jutil.DoubleToString(lval))+len(jutil.DoubleToString(rval))+1 &&

		// Do not try to fold arithmetic for numbers > 2^53. After that
		// point, fixed-point math starts to break down and become inaccurate.
		math.Abs(result) <= peepholeFoldConstants_MAX_FOLD_NUMBER) ||
		math.IsNaN(result) ||
		result == math.Inf(1) ||
		result == math.Inf(-1) {
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

// FoldAddTemplateLiterals folds additions of template literals with
// placeholders/substitutions like `${foo} and ${bar}` + `${baz}`.
//
// This code path is only hit if the left and right are both template
// literals WITH placeholders/substitutions. If one of them did not have
// placeholders/substitutions, it would have already have been optimized to a
// string.
//
// Java: PeepholeFoldConstants.foldAddTemplateLiterals(Node, Node, Node)
func (p *PeepholeFoldConstants) FoldAddTemplateLiterals(oldNode rhino.Node, left rhino.Node, right rhino.Node) rhino.Node {
	// All template literal nodes with placeholders/substitutions will have their first child
	// and last child as a string node.
	// Merge the last left child string node with the first right child node.
	lastLeft := left.GetLastChild()
	firstRight := right.GetFirstChild()
	if p.EndsWithUnescapedDollar(lastLeft.GetRawString()) &&
		jutil.StringStartsWith(firstRight.GetRawString(), "{") {
		return oldNode
	}
	cooked := jutil.ToJavaString(lastLeft.GetCookedString()) + jutil.ToJavaString(firstRight.GetCookedString())
	mergeNode :=
		rhino.IR_TemplateLiteralString(
			&cooked,
			lastLeft.GetRawString()+firstRight.GetRawString())
	mergeNode.SrcrefTreeIfMissing(lastLeft)

	// Detach the lastLeft and firstRight children and append all right children to the left node.
	lastLeft.Detach()
	firstRight.Detach()
	left.AddChildToBack(mergeNode)
	left.AddChildrenToBack(right.RemoveChildren())

	// Detach the updated left node and replace the old node with it.
	left.Detach()
	return p.Replace(oldNode, left)
}

// IsSimpleTemplateStringValue simplifies escaping and value calculations by
// avoiding problematic characters when merging strings into template
// literals.
//
// Java: PeepholeFoldConstants.isSimpleTemplateStringValue(String)
func (p *PeepholeFoldConstants) IsSimpleTemplateStringValue(s string) bool {
	chars := jutil.StringToUTF16(s)
	for i := 0; i < len(chars); i++ {
		switch chars[i] {
		// See: https://tc39.es/ecma262/#prod-LineTerminator
		case '`', '$', '{', '\\', '\n', '\r', 0x2028, 0x2029:
			return false
		default:
		}
	}
	return true
}

// EndsWithUnescapedBackslash returns true if the raw string ends with an
// unescaped backslash.
//
// Java: PeepholeFoldConstants.endsWithUnescapedBackslash(String)
func (p *PeepholeFoldConstants) EndsWithUnescapedBackslash(raw string) bool {
	backslashCount := 0
	chars := jutil.StringToUTF16(raw)
	for i := len(chars) - 1; i >= 0; i-- {
		if chars[i] == '\\' {
			backslashCount++
		} else {
			break
		}
	}
	return backslashCount%2 != 0
}

// EndsWithUnescapedDollar returns true if the raw string ends with an
// unescaped dollar sign.
//
// Java: PeepholeFoldConstants.endsWithUnescapedDollar(String)
func (p *PeepholeFoldConstants) EndsWithUnescapedDollar(raw string) bool {
	if !jutil.StringEndsWith(raw, "$") {
		return false
	}
	backslashCount := 0
	chars := jutil.StringToUTF16(raw)
	for i := len(chars) - 2; i >= 0; i-- {
		if chars[i] == '\\' {
			backslashCount++
		} else {
			break
		}
	}
	return backslashCount%2 == 0
}

// Java: PeepholeFoldConstants.tryFoldTemplateLiteralSubstitutions(Node)
func (p *PeepholeFoldConstants) TryFoldTemplateLiteralSubstitutions(n rhino.Node) rhino.Node {
	if !n.GetParent().IsNil() && n.GetParent().IsTaggedTemplateLit() {
		return n
	}

	changed := false
	cur := n.GetFirstChild()
	if cur.IsNil() {
		return n
	}

	for !cur.IsNil() {
		next := cur.GetNext()
		if next.IsNil() {
			break
		}

		sub, nextIsSub := ast.AsTemplateLitSub(next)
		if cur.IsTemplateLitString() && nextIsSub {
			expr := sub.Expr().Node()
			if expr.IsStringLit() || expr.IsNumber() || expr.IsTrue() || expr.IsFalse() {
				strValue, ok := p.GetSideEffectFreeStringValue(expr)
				if ok && p.IsSimpleTemplateStringValue(strValue) {
					nextNext := next.GetNext()
					if !nextNext.IsNil() && nextNext.IsTemplateLitString() {
						curCooked := *cur.GetCookedString()
						nextNextCooked := *nextNext.GetCookedString()

						// Avoid creating a new `${` sequence when folding.
						if strValue == "" &&
							p.EndsWithUnescapedDollar(cur.GetRawString()) &&
							jutil.StringStartsWith(nextNext.GetRawString(), "{") {
							cur = next
							continue
						}

						combinedCooked := curCooked + strValue + nextNextCooked
						combinedRaw := cur.GetRawString() + strValue + nextNext.GetRawString()

						// Avoid creating a restricted octal escape sequence (e.g. \07)
						if p.IsRestrictedOctalEscape(cur.GetRawString(), strValue) ||
							p.IsRestrictedOctalEscape(
								cur.GetRawString()+strValue, nextNext.GetRawString()) {
							cur = next
							continue
						}

						newStringNode := rhino.IR_TemplateLiteralString(&combinedCooked, combinedRaw)
						newStringNode.SrcrefTree(cur)

						cur.ReplaceWith(newStringNode)
						next.Detach()
						nextNext.Detach()

						cur = newStringNode
						changed = true
						continue
					}
				}
			}
		}
		cur = cur.GetNext()
	}

	if changed {
		p.ReportChangeToEnclosingScope(n)
	}
	return n
}

// Java: PeepholeFoldConstants.isRestrictedOctalEscape(String, String)
func (p *PeepholeFoldConstants) IsRestrictedOctalEscape(prefixRaw string, suffixRaw string) bool {
	if suffixRaw == "" {
		return false
	}
	firstChar := jutil.StringCharAt(suffixRaw, 0)
	if firstChar < '0' || firstChar > '9' {
		return false
	}

	// Check if prefix ends with an unescaped \0
	if jutil.StringEndsWith(prefixRaw, "0") {
		beforeZero := jutil.StringSubstring(prefixRaw, 0, jutil.StringLength(prefixRaw)-1)
		return p.EndsWithUnescapedBackslash(beforeZero)
	}
	return false
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
		} else if left.IsStringLit() && left.GetString() == "" && p.IsStringTyped(right) {
			return p.Replace(node, right.CloneTreeWithTypeExprs(true))
		} else if right.IsStringLit() && right.GetString() == "" && p.IsStringTyped(left) {
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

// Java: PeepholeFoldConstants.foldTemplateLiteralAndConstant(Node, Node, String)
func (p *PeepholeFoldConstants) FoldTemplateLiteralAndConstant(oldNode rhino.Node, left rhino.Node, rightStr string) rhino.Node {
	lastLeft := left.GetLastChild()
	jutil.CheckState(lastLeft.IsTemplateLitString())

	if p.IsRestrictedOctalEscape(lastLeft.GetRawString(), rightStr) {
		return oldNode
	}

	cooked := jutil.ToJavaString(lastLeft.GetCookedString()) + rightStr
	mergeNode :=
		rhino.IR_TemplateLiteralString(
			&cooked, lastLeft.GetRawString()+rightStr)
	mergeNode.SrcrefTreeIfMissing(lastLeft)

	lastLeft.ReplaceWith(mergeNode)
	left.Detach()
	return p.Replace(oldNode, left)
}

// Java: PeepholeFoldConstants.foldConstantAndTemplateLiteral(Node, String, Node)
func (p *PeepholeFoldConstants) FoldConstantAndTemplateLiteral(oldNode rhino.Node, leftStr string, right rhino.Node) rhino.Node {
	firstRight := right.GetFirstChild()
	jutil.CheckState(firstRight.IsTemplateLitString())

	if p.IsRestrictedOctalEscape(leftStr, firstRight.GetRawString()) {
		return oldNode
	}

	cooked := leftStr + jutil.ToJavaString(firstRight.GetCookedString())
	mergeNode :=
		rhino.IR_TemplateLiteralString(
			&cooked, leftStr+firstRight.GetRawString())
	mergeNode.SrcrefTreeIfMissing(firstRight)

	firstRight.ReplaceWith(mergeNode)
	right.Detach()
	return p.Replace(oldNode, right)
}

// Java: PeepholeFoldConstants.replace(Node, Node)
func (p *PeepholeFoldConstants) Replace(oldNode rhino.Node, newNode rhino.Node) rhino.Node {
	oldNode.ReplaceWith(newNode)
	p.ReportChangeToEnclosingScope(newNode)
	return newNode
}

// Java: PeepholeFoldConstants.isStringTyped(Node)
func (p *PeepholeFoldConstants) IsStringTyped(n rhino.Node) bool {
	// We could also accept !String, but it is unlikely to be very common.
	if NodeUtil_IsStringResult(n) {
		return true
	}
	if p.shouldUseTypes {
		return n.GetColor() != nil && n.GetColor().Equals(colors.StandardColors_STRING)
	}
	return false
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
