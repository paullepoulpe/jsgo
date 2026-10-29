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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

var (
	foldIDDefineProperties = uint32(strpool.Intern("defineProperties"))
	foldIDString           = uint32(strpool.Intern("String"))
	foldIDLength           = uint32(strpool.Intern("length"))
	foldIDNaN              = uint32(strpool.Intern("NaN"))
	foldIDObject           = uint32(strpool.Intern("Object"))
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

var peepholeFoldConstants_MAX_FOLD_NUMBER = base.MathPow(2, 53)

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
	case rhino.Token_NEG:
		if c := n.GetFirstChild(); c.IsNumber() && !math.IsInf(c.GetDouble(), 0) && !math.IsNaN(c.GetDouble()) {
			return
		}
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
		if argumentNode.Node().GetStringID() == nodeUtilIDUndefined {
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
		} else if right.IsName() && right.GetStringID() == foldIDObject {
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
		strings.HasPrefix(firstRight.GetRawString(), "{") {
		return oldNode
	}
	cooked := peepholeCookedString(lastLeft) + peepholeCookedString(firstRight)
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
	for i := 0; i < len(s); i++ {
		switch s[i] {
		// See: https://tc39.es/ecma262/#prod-LineTerminator
		case '`', '$', '{', '\\', '\n', '\r':
			return false
		case 0xe2:
			// U+2028 (E2 80 A8) and U+2029 (E2 80 A9)
			if i+2 < len(s) && s[i+1] == 0x80 && (s[i+2] == 0xa8 || s[i+2] == 0xa9) {
				return false
			}
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
	for i := len(raw) - 1; i >= 0 && raw[i] == '\\'; i-- {
		backslashCount++
	}
	return backslashCount%2 != 0
}

// EndsWithUnescapedDollar returns true if the raw string ends with an
// unescaped dollar sign.
//
// Java: PeepholeFoldConstants.endsWithUnescapedDollar(String)
func (p *PeepholeFoldConstants) EndsWithUnescapedDollar(raw string) bool {
	if len(raw) == 0 || raw[len(raw)-1] != '$' {
		return false
	}
	backslashCount := 0
	for i := len(raw) - 2; i >= 0 && raw[i] == '\\'; i-- {
		backslashCount++
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
						curRaw := cur.GetRawString()
						nextNextRaw := nextNext.GetRawString()

						// Avoid creating a new `${` sequence when folding.
						if strValue == "" &&
							p.EndsWithUnescapedDollar(curRaw) &&
							strings.HasPrefix(nextNextRaw, "{") {
							cur = next
							continue
						}

						// Avoid creating a restricted octal escape sequence (e.g. \07)
						if p.IsRestrictedOctalEscape(curRaw, strValue) ||
							(len(nextNextRaw) > 0 && nextNextRaw[0] >= '0' && nextNextRaw[0] <= '9' &&
								p.IsRestrictedOctalEscape(curRaw+strValue, nextNextRaw)) {
							cur = next
							continue
						}

						combinedCooked := *cur.GetCookedString() + strValue + *nextNext.GetCookedString()
						combinedRaw := curRaw + strValue + nextNextRaw

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
	if len(suffixRaw) == 0 || suffixRaw[0] < '0' || suffixRaw[0] > '9' {
		return false
	}

	// Check if prefix ends with an unescaped \0
	if len(prefixRaw) > 0 && prefixRaw[len(prefixRaw)-1] == '0' {
		return p.EndsWithUnescapedBackslash(prefixRaw[:len(prefixRaw)-1])
	}
	return false
}

func peepholeCookedString(n rhino.Node) string {
	if s := n.GetCookedString(); s != nil {
		return *s
	}
	return "null"
}

// Java: PeepholeFoldConstants.foldTemplateLiteralAndConstant(Node, Node, String)
func (p *PeepholeFoldConstants) FoldTemplateLiteralAndConstant(oldNode rhino.Node, left rhino.Node, rightStr string) rhino.Node {
	lastLeft := left.GetLastChild()
	jutil.CheckState(lastLeft.IsTemplateLitString())

	if p.IsRestrictedOctalEscape(lastLeft.GetRawString(), rightStr) {
		return oldNode
	}

	cooked := peepholeCookedString(lastLeft) + rightStr
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

	cooked := leftStr + peepholeCookedString(firstRight)
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
