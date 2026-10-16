// Port of com/google/javascript/jscomp/PeepholeFoldConstants.java @ closure-compiler v20260415

package jscomp

import (
	"math"
	"math/big"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

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
	p := &PeepholeFoldConstants{}
	p.InitAbstractPeepholeOptimizationBase(p) // the implicit super()
	p.late = late
	p.shouldUseTypes = shouldUseTypes
	return p
}

// Java: PeepholeFoldConstants.optimizeSubtree(Node)
func (p *PeepholeFoldConstants) OptimizeSubtree(subtree *rhino.Node) *rhino.Node {
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
func (p *PeepholeFoldConstants) TryFoldBinaryOperator(subtree *rhino.Node) *rhino.Node {
	left := subtree.GetFirstChild()

	if left == nil {
		return subtree
	}

	right := left.GetNext()

	if right == nil {
		return subtree
	}

	// If we've reached here, node is truly a binary operator.
	switch subtree.GetToken() {
	case rhino.Token_GETELEM, rhino.Token_OPTCHAIN_GETELEM:
		return p.TryFoldGetElem(subtree, left, right)
	case rhino.Token_INSTANCEOF:
		return p.TryFoldInstanceof(subtree, left, right)
	case rhino.Token_AND, rhino.Token_OR:
		return p.TryFoldAndOr(subtree, left, right)
	case rhino.Token_COALESCE:
		return p.TryFoldCoalesce(subtree, left, right)
	case rhino.Token_LSH, rhino.Token_RSH, rhino.Token_URSH:
		return p.TryFoldShift(subtree, left, right)
	case rhino.Token_ASSIGN:
		return p.TryFoldAssign(subtree, left, right)
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
		return p.TryUnfoldAssignOp(subtree, left, right)
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
func (p *PeepholeFoldConstants) TryReduceVoid(n *rhino.Node) *rhino.Node {
	child := n.GetFirstChild()
	if (!child.IsNumber() || child.GetDouble() != 0.0) && !p.MayHaveSideEffects(n) {
		child.ReplaceWith(rhino.IR_Number(0))
		p.ReportChangeToEnclosingScope(n)
	}
	return n
}

// Java: PeepholeFoldConstants.tryReduceOperandsForOp(Node)
func (p *PeepholeFoldConstants) TryReduceOperandsForOp(n *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_ADD:
		left := n.GetFirstChild()
		right := n.GetLastChild()
		if !NodeUtil_MayBeStringWithUseType(left, p.shouldUseTypes) &&
			!NodeUtil_MayBeStringWithUseType(right, p.shouldUseTypes) {
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
		// TODO(johnlenz): convert these to integers.
		p.TryConvertToNumber(n.GetLastChild())
	case rhino.Token_BITNOT, rhino.Token_BITOR, rhino.Token_BITXOR, rhino.Token_BITAND, rhino.Token_LSH, rhino.Token_RSH, rhino.Token_URSH, rhino.Token_SUB, rhino.Token_MUL, rhino.Token_MOD, rhino.Token_DIV, rhino.Token_POS, rhino.Token_NEG, rhino.Token_EXPONENT:
		p.TryConvertOperandsToNumber(n)
	default:
	}
}

// Java: PeepholeFoldConstants.tryConvertOperandsToNumber(Node)
func (p *PeepholeFoldConstants) TryConvertOperandsToNumber(n *rhino.Node) {
	var next *rhino.Node
	for c := n.GetFirstChild(); c != nil; c = next {
		next = c.GetNext()
		p.TryConvertToNumber(c)
	}
}

// Java: PeepholeFoldConstants.tryConvertToNumber(Node)
func (p *PeepholeFoldConstants) TryConvertToNumber(n *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_NUMBER:
		// Nothing to do
		return
	case rhino.Token_AND, rhino.Token_OR, rhino.Token_COMMA, rhino.Token_COALESCE:
		p.TryConvertToNumber(n.GetLastChild())
		return
	case rhino.Token_HOOK:
		p.TryConvertToNumber(n.GetSecondChild())
		p.TryConvertToNumber(n.GetLastChild())
		return
	case rhino.Token_NAME:
		if !NodeUtil_IsUndefined(n) {
			return
		}
	default:
	}

	result := p.GetSideEffectFreeNumberValue(n)
	if result == nil {
		return
	}

	value := *result

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
func (p *PeepholeFoldConstants) TryFoldTypeof(originalTypeofNode *rhino.Node) *rhino.Node {
	jutil.CheckArgument(originalTypeofNode.IsTypeOf())

	argumentNode := originalTypeofNode.GetFirstChild()
	if argumentNode == nil || !NodeUtil_IsLiteralValue(argumentNode, true) {
		return originalTypeofNode
	}

	var typeNameString *string

	switch argumentNode.GetToken() {
	case rhino.Token_FUNCTION:
		s := "function"
		typeNameString = &s
	case rhino.Token_STRINGLIT:
		s := "string"
		typeNameString = &s
	case rhino.Token_NUMBER:
		s := "number"
		typeNameString = &s
	case rhino.Token_TRUE, rhino.Token_FALSE:
		s := "boolean"
		typeNameString = &s
	case rhino.Token_NULL, rhino.Token_OBJECTLIT, rhino.Token_ARRAYLIT:
		s := "object"
		typeNameString = &s
	case rhino.Token_VOID:
		s := "undefined"
		typeNameString = &s
	case rhino.Token_NAME:
		// We assume here that programs don't change the value of the
		// keyword undefined to something other than the value undefined.
		if argumentNode.GetString() == "undefined" {
			s := "undefined"
			typeNameString = &s
		}
	default:
	}

	if typeNameString != nil {
		newNode := rhino.IR_String(*typeNameString)
		p.ReportChangeToEnclosingScope(originalTypeofNode)
		originalTypeofNode.ReplaceWith(newNode)
		p.MarkFunctionsDeleted(originalTypeofNode)

		return newNode
	}

	return originalTypeofNode
}

// Java: PeepholeFoldConstants.tryFoldUnaryOperator(Node)
func (p *PeepholeFoldConstants) TryFoldUnaryOperator(n *rhino.Node) *rhino.Node {
	jutil.CheckStateMsg1(n.HasOneChild(), "%s", n)

	left := n.GetFirstChild()
	parent := n.GetParent()

	if left == nil {
		return n
	}

	leftVal := p.GetSideEffectFreeBooleanValue(left)
	if leftVal == base.Tri_UNKNOWN {
		return n
	}

	switch n.GetToken() {
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
		var result *rhino.Node
		if left.IsName() && left.GetString() == "NaN" {
			result = left.Detach() // "-NaN" is "NaN".
		} else if left.IsNeg() {
			leftLeft := left.GetOnlyChild()
			if leftLeft.IsBigInt() || leftLeft.IsNumber() {
				result = leftLeft.Detach() // `-(-4)` is `4`
			}
		}

		if result != nil {
			n.ReplaceWith(result)
			p.ReportChangeToEnclosingScope(parent)
			return result
		}
		return n
	case rhino.Token_BITNOT:
		doubleVal := p.GetSideEffectFreeNumberValue(left)
		if doubleVal != nil {
			if base.JSCompDoubles_IsMathematicalInteger(*doubleVal) {
				intVal := base.JSCompDoubles_EcmascriptToInt32(*doubleVal)
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

// Java: PeepholeFoldConstants.isReasonableDoubleValue(Double)
func (p *PeepholeFoldConstants) IsReasonableDoubleValue(x *float64) bool {
	return x != nil && !math.IsInf(*x, 0) && !math.IsNaN(*x)
}

// TryFoldInstanceof tries to fold `left instanceof right` into true or
// false.
//
// Java: PeepholeFoldConstants.tryFoldInstanceof(Node, Node, Node)
func (p *PeepholeFoldConstants) TryFoldInstanceof(n *rhino.Node, left *rhino.Node, right *rhino.Node) *rhino.Node {
	jutil.CheckArgument(n.IsInstanceOf())

	// TODO(johnlenz) Use type information if available to fold
	// instanceof.
	if NodeUtil_IsLiteralValue(left, true) && !p.MayHaveSideEffects(right) {

		var replacementNode *rhino.Node

		if NodeUtil_IsImmutableValue(left) {
			// Non-object types are never instances.
			replacementNode = rhino.IR_FalseNode()
		} else if right.IsName() && right.GetString() == "Object" {
			replacementNode = rhino.IR_TrueNode()
		}

		if replacementNode != nil {
			n.ReplaceWith(replacementNode)
			p.ReportChangeToEnclosingScope(replacementNode)
			p.MarkFunctionsDeleted(n)
			return replacementNode
		}
	}

	return n
}

// Java: PeepholeFoldConstants.tryFoldAssign(Node, Node, Node)
func (p *PeepholeFoldConstants) TryFoldAssign(n *rhino.Node, left *rhino.Node, right *rhino.Node) *rhino.Node {
	jutil.CheckArgument(n.IsAssign())

	if !p.late {
		return n
	}

	// Only certain RHS operators are supported.
	operator := right.GetToken()
	var newType rhino.Token
	hasNewType := true
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
		hasNewType = false
	}
	if !hasNewType {
		return n
	}

	// Tries to convert x = x + y -> x += y;
	if !right.HasChildren() || right.GetSecondChild() != right.GetLastChild() {
		// RHS must have two children.
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

	var newRight *rhino.Node
	if p.AreNodesEqualForInlining(left, right.GetFirstChild()) {
		newRight = right.GetLastChild()
	} else if NodeUtil_IsCommutative(operator) &&
		p.AreNodesEqualForInlining(left, right.GetLastChild()) {
		newRight = right.GetFirstChild()
	} else if NodeUtil_IsAssociative(operator) &&
		right.GetFirstChild().GetToken() == operator &&
		p.AreNodesEqualForInlining(left, right.GetFirstFirstChild()) {

		newLeft := right.GetFirstChild().GetLastChild().Detach()
		right.RemoveFirstChild()
		right.AddChildToFront(newLeft)
		newRight = right

	} else {
		return n
	}

	newNode := rhino.NewNodeWith2Children(newType, left.Detach(), newRight.Detach())
	n.ReplaceWith(newNode)

	p.ReportChangeToEnclosingScope(newNode)

	return newNode
}

// Java: PeepholeFoldConstants.tryUnfoldAssignOp(Node, Node, Node)
func (p *PeepholeFoldConstants) TryUnfoldAssignOp(n *rhino.Node, left *rhino.Node, right *rhino.Node) *rhino.Node {
	if p.late {
		return n
	}

	if !n.HasChildren() || n.GetSecondChild() != n.GetLastChild() {
		return n
	}

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
func (p *PeepholeFoldConstants) TryFoldAndOr(n *rhino.Node, left *rhino.Node, right *rhino.Node) *rhino.Node {
	parent := n.GetParent()

	var result *rhino.Node
	var dropped *rhino.Node

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
			dropped = nil
		}
	} else if parent.GetToken() == type_ && n == parent.GetFirstChild() {
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

	if result != nil {
		// Fold it!
		n.DetachChildren()
		n.ReplaceWith(result)
		p.ReportChangeToEnclosingScope(result)
		if dropped != nil {
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
func (p *PeepholeFoldConstants) TryFoldCoalesce(n *rhino.Node, left *rhino.Node, right *rhino.Node) *rhino.Node {

	var result *rhino.Node

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

	if result != nil {
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
// Java: PeepholeFoldConstants.tryFoldAdjacentLiteralLeaves(Node, Node, Node)
func (p *PeepholeFoldConstants) TryFoldAdjacentLiteralLeaves(n *rhino.Node, left *rhino.Node, right *rhino.Node) *rhino.Node {

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
			(rr != nil && NodeUtil_IsStringResult(rr)) ||
				(rr != nil &&
					// If right.getNext() isn't a string result, right must be a string to be folded
					right.IsStringLit() &&
					// Access the right grandparent safely...
					rightGrandparent != nil &&
					rightGrandparent.IsAdd() &&
					// The second child of `rightGrandparent` is what `right + rr` is being added to.
					// If it exists, `right` can only be safely removed if `rr` is still treated as a
					// string in the resulting addition.
					NodeUtil_IsStringResult(rightGrandparent.GetSecondChild())) ||
				// Dangling literal term has no typecasting side effects; fold it!
				(rr == nil)
		if foldIsTypeSafe {
			result := left.GetString() + jutil.ToJavaString(NodeUtil_GetStringValue(right))
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
func (p *PeepholeFoldConstants) TryFoldAddConstantString(n *rhino.Node, left *rhino.Node, right *rhino.Node) *rhino.Node {
	if left.IsStringLit() || right.IsStringLit() || left.IsArrayLit() || right.IsArrayLit() {
		// Add strings.
		leftString := p.GetSideEffectFreeStringValue(left)
		rightString := p.GetSideEffectFreeStringValue(right)
		if leftString != nil && rightString != nil {
			newStringNode := rhino.IR_String(*leftString + *rightString)
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
func (p *PeepholeFoldConstants) TryFoldArithmeticOp(n *rhino.Node, left *rhino.Node, right *rhino.Node) *rhino.Node {
	result := p.PerformArithmeticOp(n, left, right)
	if result != nil {
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
func (p *PeepholeFoldConstants) PerformArithmeticOp(n *rhino.Node, left *rhino.Node, right *rhino.Node) *rhino.Node {
	// Unlike other operations, ADD operands are not always converted
	// to Number.
	if n.IsAdd() &&
		(NodeUtil_MayBeStringWithUseType(left, p.shouldUseTypes) ||
			NodeUtil_MayBeStringWithUseType(right, p.shouldUseTypes)) {
		return nil
	}

	if p.IsBigInt(left) && p.IsBigInt(right) {
		return p.PerformBigIntArithmeticOp(n, left, right)
	}

	var result float64

	// TODO(johnlenz): Handle NaN with unknown value. BIT ops convert NaN
	// to zero so this is a little awkward here.

	lValObj := p.GetSideEffectFreeNumberValue(left)
	rValObj := p.GetSideEffectFreeNumberValue(right)
	// at least one of the two operands must have a value and both must be numeric
	if (lValObj == nil && rValObj == nil) || !p.IsNumeric(left) || !p.IsNumeric(right) {
		return nil
	}
	// handle the operations that have algebraic identities, since we can simplify the tree without
	// actually knowing the value statically.
	switch n.GetToken() {
	case rhino.Token_ADD:
		if lValObj != nil && rValObj != nil {
			return p.MaybeReplaceBinaryOpWithNumericResult(*lValObj+*rValObj, *lValObj, *rValObj)
		}
		if lValObj != nil && *lValObj == 0 {
			return right.CloneTreeWithTypeExprs(true)
		} else if rValObj != nil && *rValObj == 0 {
			return left.CloneTreeWithTypeExprs(true)
		}
		return nil
	case rhino.Token_SUB:
		if lValObj != nil && rValObj != nil {
			return p.MaybeReplaceBinaryOpWithNumericResult(*lValObj-*rValObj, *lValObj, *rValObj)
		}
		if lValObj != nil && *lValObj == 0 {
			// 0 - x -> -x
			// NOTE: this optimization has the subtle side effect of changing `0` to `-0` because
			// `0-0 -> 0` but `-0 -> -0`.
			return rhino.IR_Neg(right.CloneTreeWithTypeExprs(true))
		} else if rValObj != nil && *rValObj == 0 {
			// x - 0 -> x
			return left.CloneTreeWithTypeExprs(true)
		}
		return nil
	case rhino.Token_MUL:
		if lValObj != nil && rValObj != nil {
			return p.MaybeReplaceBinaryOpWithNumericResult(*lValObj**rValObj, *lValObj, *rValObj)
		}
		// NOTE: 0*x != 0 for all x, if x==0, then it is NaN.  So we can't take advantage of that
		// without some kind of non-NaN proof.  So the special cases here only deal with 1*x
		if lValObj != nil {
			if *lValObj == 1 {
				return right.CloneTreeWithTypeExprs(true)
			}
		} else {
			if *rValObj == 1 {
				return left.CloneTreeWithTypeExprs(true)
			}
		}
		return nil
	case rhino.Token_DIV:
		if lValObj != nil && rValObj != nil {
			if *rValObj == 0 {
				return nil
			}
			return p.MaybeReplaceBinaryOpWithNumericResult(*lValObj / *rValObj, *lValObj, *rValObj)
		}
		// NOTE: 0/x != 0 for all x, if x==0, then it is NaN
		if rValObj != nil {
			if *rValObj == 1 {
				// x/1->x
				return left.CloneTreeWithTypeExprs(true)
			}
		}
		return nil
	case rhino.Token_EXPONENT:
		if lValObj != nil && rValObj != nil {
			return p.MaybeReplaceBinaryOpWithNumericResult(
				jutil.MathPow(*lValObj, *rValObj), *lValObj, *rValObj)
		}
		return nil
	default:
	}
	if lValObj == nil || rValObj == nil {
		return nil
	}
	lval := *lValObj
	rval := *rValObj
	switch n.GetToken() {
	case rhino.Token_BITAND:
		result = float64(base.JSCompDoubles_EcmascriptToInt32(lval) & base.JSCompDoubles_EcmascriptToInt32(rval))
	case rhino.Token_BITOR:
		result = float64(base.JSCompDoubles_EcmascriptToInt32(lval) | base.JSCompDoubles_EcmascriptToInt32(rval))
	case rhino.Token_BITXOR:
		result = float64(base.JSCompDoubles_EcmascriptToInt32(lval) ^ base.JSCompDoubles_EcmascriptToInt32(rval))
	case rhino.Token_MOD:
		if rval == 0 {
			return nil
		}
		result = jutil.DRem(lval, rval)
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Unexpected arithmetic operator: " + n.GetToken().String()))
	}
	return p.MaybeReplaceBinaryOpWithNumericResult(result, lval, rval)
}

// Java: PeepholeFoldConstants.performBigIntArithmeticOp(Node, Node, Node)
func (p *PeepholeFoldConstants) PerformBigIntArithmeticOp(n *rhino.Node, left *rhino.Node, right *rhino.Node) *rhino.Node {
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
				return nil
			}
		case rhino.Token_EXPONENT:
			var result *rhino.Node
			jutil.Try(func() {
				result = PeepholeFoldConstants_BigintNode(jutil.BigIntegerPow(lVal, jutil.BigIntegerIntValueExact(rVal)), n)
			}, func(exception jutil.IArithmeticException) {
				result = nil
			})
			return result
		case rhino.Token_MOD:
			return PeepholeFoldConstants_BigintNode(jutil.BigIntegerMod(lVal, rVal), n)
		case rhino.Token_BITAND:
			return PeepholeFoldConstants_BigintNode(new(big.Int).And(lVal, rVal), n)
		case rhino.Token_BITOR:
			return PeepholeFoldConstants_BigintNode(new(big.Int).Or(lVal, rVal), n)
		case rhino.Token_BITXOR:
			return PeepholeFoldConstants_BigintNode(new(big.Int).Xor(lVal, rVal), n)
		default:
			return nil
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
		return nil
	case rhino.Token_SUB:
		if lVal != nil && jutil.BigIntegerEquals(lVal, jutil.BigInteger_ZERO) {
			// 0n - x -> -x
			return rhino.IR_Neg(right.CloneTreeWithTypeExprs( /* cloneTypeExprs= */ true))
		} else if rVal != nil && jutil.BigIntegerEquals(rVal, jutil.BigInteger_ZERO) {
			// x - 0n -> x
			return left.CloneTreeWithTypeExprs( /* cloneTypeExprs= */ true)
		}
		return nil
	case rhino.Token_MUL:
		if lVal != nil && jutil.BigIntegerEquals(lVal, jutil.BigInteger_ONE) {
			return right.CloneTreeWithTypeExprs( /* cloneTypeExprs= */ true)
		} else if rVal != nil && jutil.BigIntegerEquals(rVal, jutil.BigInteger_ONE) {
			return left.CloneTreeWithTypeExprs( /* cloneTypeExprs= */ true)
		}
		return nil
	case rhino.Token_DIV:
		if rVal != nil && jutil.BigIntegerEquals(rVal, jutil.BigInteger_ONE) {
			return left.CloneTreeWithTypeExprs( /* cloneTypeExprs= */ true)
		}
		return nil
	default:
		return nil
	}
}

// Java: PeepholeFoldConstants.isNumeric(Node)
func (p *PeepholeFoldConstants) IsNumeric(n *rhino.Node) bool {
	if NodeUtil_IsNumericResult(n) {
		return true
	}
	if p.shouldUseTypes {
		return n.GetColor() != nil && n.GetColor().Equals(colors.StandardColors_NUMBER)
	}
	return false
}

// Java: PeepholeFoldConstants.isBigInt(Node)
func (p *PeepholeFoldConstants) IsBigInt(n *rhino.Node) bool {
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
	result float64, lval float64, rval float64) *rhino.Node {
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
		return NodeUtil_NumberNode(result, nil)
	}
	return nil
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
func (p *PeepholeFoldConstants) TryFoldLeftChildOp(n *rhino.Node, left *rhino.Node, right *rhino.Node) *rhino.Node {
	opType := n.GetToken()
	jutil.CheckState(NodeUtil_IsAssociative(opType) && NodeUtil_IsCommutative(opType))

	// Use getNumberValue to handle constants like "NaN" and "Infinity"
	// other values are converted to numbers elsewhere.
	rightValObj := p.GetSideEffectFreeNumberValue(right)
	rightBigInt := p.GetSideEffectFreeBigIntValue(right)
	if (rightValObj != nil || rightBigInt != nil) && left.GetToken() == opType {
		jutil.CheckState(left.HasTwoChildren())

		ll := left.GetFirstChild()
		lr := ll.GetNext()

		valueToCombine := ll
		replacement := p.PerformArithmeticOp(n, valueToCombine, right)
		if replacement == nil {
			valueToCombine = lr
			replacement = p.PerformArithmeticOp(n, valueToCombine, right)
		}
		if replacement != nil {
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
func (p *PeepholeFoldConstants) FoldAddTemplateLiterals(oldNode *rhino.Node, left *rhino.Node, right *rhino.Node) *rhino.Node {
	// All template literal nodes with placeholders/substitutions will have their first child
	// and last child as a string node.
	// Merge the last left child string node with the first right child node.
	lastLeft := left.GetLastChild()
	firstRight := right.GetFirstChild()
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

// Java: PeepholeFoldConstants.tryFoldTemplateLiteralSubstitutions(Node)
func (p *PeepholeFoldConstants) TryFoldTemplateLiteralSubstitutions(n *rhino.Node) *rhino.Node {
	if n.GetParent() != nil && n.GetParent().IsTaggedTemplateLit() {
		return n
	}

	changed := false
	cur := n.GetFirstChild()
	if cur == nil {
		return n
	}

	for cur != nil {
		next := cur.GetNext()
		if next == nil {
			break
		}

		if cur.IsTemplateLitString() && next.IsTemplateLitSub() {
			expr := next.GetFirstChild()
			if expr != nil &&
				(expr.IsStringLit() || expr.IsNumber() || expr.IsTrue() || expr.IsFalse()) {
				strValue := p.GetSideEffectFreeStringValue(expr)
				if strValue != nil && p.IsSimpleTemplateStringValue(*strValue) {
					nextNext := next.GetNext()
					if nextNext != nil && nextNext.IsTemplateLitString() {
						curCooked := *cur.GetCookedString()
						nextNextCooked := *nextNext.GetCookedString()

						// Avoid creating a new `${` sequence when folding.
						if *strValue == "" && jutil.StringEndsWith(curCooked, "$") && jutil.StringStartsWith(nextNextCooked, "{") {
							cur = next
							continue
						}

						combinedCooked := curCooked + *strValue + nextNextCooked
						combinedRaw := cur.GetRawString() + *strValue + nextNext.GetRawString()

						// Avoid creating a restricted octal escape sequence (e.g. \07)
						if p.IsRestrictedOctalEscape(cur.GetRawString(), *strValue) ||
							p.IsRestrictedOctalEscape(
								cur.GetRawString()+*strValue, nextNext.GetRawString()) {
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
func (p *PeepholeFoldConstants) TryFoldAdd(node *rhino.Node, left *rhino.Node, right *rhino.Node) *rhino.Node {
	jutil.CheckArgument(node.IsAdd())

	if NodeUtil_MayBeStringWithUseType(node, p.shouldUseTypes) {
		if left.IsTemplateLit() && right.IsTemplateLit() {
			return p.FoldAddTemplateLiterals(node, left, right)
		} else if left.IsTemplateLit() {
			rightStr := p.GetSideEffectFreeStringValue(right)
			if rightStr != nil && p.IsSimpleTemplateStringValue(*rightStr) {
				return p.FoldTemplateLiteralAndConstant(node, left, *rightStr)
			}
		} else if right.IsTemplateLit() {
			leftStr := p.GetSideEffectFreeStringValue(left)
			if leftStr != nil && p.IsSimpleTemplateStringValue(*leftStr) {
				return p.FoldConstantAndTemplateLiteral(node, *leftStr, right)
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
func (p *PeepholeFoldConstants) FoldTemplateLiteralAndConstant(oldNode *rhino.Node, left *rhino.Node, rightStr string) *rhino.Node {
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
func (p *PeepholeFoldConstants) FoldConstantAndTemplateLiteral(oldNode *rhino.Node, leftStr string, right *rhino.Node) *rhino.Node {
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
func (p *PeepholeFoldConstants) Replace(oldNode *rhino.Node, newNode *rhino.Node) *rhino.Node {
	oldNode.ReplaceWith(newNode)
	p.ReportChangeToEnclosingScope(newNode)
	return newNode
}

// Java: PeepholeFoldConstants.isStringTyped(Node)
func (p *PeepholeFoldConstants) IsStringTyped(n *rhino.Node) bool {
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
func (p *PeepholeFoldConstants) TryFoldShift(n *rhino.Node, left *rhino.Node, right *rhino.Node) *rhino.Node {
	leftVal := p.GetSideEffectFreeNumberValue(left)
	rightVal := p.GetSideEffectFreeNumberValue(right)

	if !p.IsReasonableDoubleValue(leftVal) || !p.IsReasonableDoubleValue(rightVal) {
		return n
	}
	if !base.JSCompDoubles_IsMathematicalInteger(*leftVal) {
		p.Report(PeepholeFoldConstants_FRACTIONAL_BITWISE_OPERAND, left)
		return n
	}
	if !base.JSCompDoubles_IsMathematicalInteger(*rightVal) {
		p.Report(PeepholeFoldConstants_FRACTIONAL_BITWISE_OPERAND, right)
		return n
	}

	// only the lower 5 bits are used when shifting, so don't do anything
	// if the shift amount is outside [0,32)
	if !(0 <= *rightVal && *rightVal < 32) {
		return n
	}

	rvalInt := jutil.D2I(*rightVal)
	bits := base.JSCompDoubles_EcmascriptToInt32(*leftVal)

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

// TryFoldComparison tries to fold comparison nodes, e.g ==
//
// Java: PeepholeFoldConstants.tryFoldComparison(Node, Node, Node)
func (p *PeepholeFoldConstants) TryFoldComparison(n *rhino.Node, left *rhino.Node, right *rhino.Node) *rhino.Node {
	result := PeepholeFoldConstants_EvaluateComparison(p, n.GetToken(), left, right)
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
	peepholeOptimization AbstractPeepholeOptimization,
	left *rhino.Node,
	right *rhino.Node,
	willNegate bool) base.Tri {
	leftValueType := NodeUtil_GetKnownValueType(left)
	rightValueType := NodeUtil_GetKnownValueType(right)
	// First, check for a string comparison.
	if leftValueType == NodeUtil_ValueType_STRING && rightValueType == NodeUtil_ValueType_STRING {
		lvStr := peepholeOptimization.GetSideEffectFreeStringValue(left)
		rvStr := peepholeOptimization.GetSideEffectFreeStringValue(right)
		if lvStr != nil && rvStr != nil {
			// In JS, browsers parse \v differently. So do not compare strings if one contains \v.
			if jutil.StringIndexOfChar(*lvStr, 0x000B) != -1 || jutil.StringIndexOfChar(*rvStr, 0x000B) != -1 {
				return base.Tri_UNKNOWN
			} else {
				return base.Tri_ForBoolean(jutil.StringCompareTo(*lvStr, *rvStr) < 0)
			}
		} else if left.IsTypeOf() &&
			right.IsTypeOf() &&
			left.GetFirstChild().IsName() &&
			right.GetFirstChild().IsName() &&
			left.GetFirstChild().GetString() == right.GetFirstChild().GetString() {
			// Special case: `typeof a < typeof a` is always false.
			return base.Tri_FALSE
		}
	}

	// Next, try to evaluate based on the value of the node. Try comparing as BigInts first.
	lvBig := peepholeOptimization.GetSideEffectFreeBigIntValue(left)
	rvBig := peepholeOptimization.GetSideEffectFreeBigIntValue(right)
	if lvBig != nil && rvBig != nil {
		return base.Tri_ForBoolean(lvBig.Cmp(rvBig) < 0)
	}

	// Then, try comparing as Numbers.
	lvNum := peepholeOptimization.GetSideEffectFreeNumberValue(left)
	rvNum := peepholeOptimization.GetSideEffectFreeNumberValue(right)
	if lvNum != nil && rvNum != nil {
		if math.IsNaN(*lvNum) || math.IsNaN(*rvNum) {
			return base.Tri_ForBoolean(willNegate)
		} else {
			return base.Tri_ForBoolean(*lvNum < *rvNum)
		}
	}

	// Finally, try comparisons between BigInt and Number.
	if lvBig != nil && rvNum != nil {
		return PeepholeFoldConstants_BigintLessThanDouble(lvBig, *rvNum, base.Tri_FALSE, willNegate)
	}
	if lvNum != nil && rvBig != nil {
		return PeepholeFoldConstants_BigintLessThanDouble(rvBig, *lvNum, base.Tri_TRUE, willNegate)
	}

	// Special case: `x < x` is always false.
	// TODO(moz): If we knew the named value wouldn't be NaN, it would be nice to handle
	// LE and GE. We should use type information if available here.
	if !willNegate && left.IsName() && right.IsName() {
		if left.GetString() == right.GetString() {
			return base.Tri_FALSE
		}
	}

	return base.Tri_UNKNOWN
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
	numberAsBigInt := big.NewInt(jutil.D2L(number))
	negativeMeansBigintSmaller := bigint.Cmp(numberAsBigInt)
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
	peepholeOptimization AbstractPeepholeOptimization, left *rhino.Node, right *rhino.Node) base.Tri {
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
			rv := peepholeOptimization.GetSideEffectFreeNumberValue(right)
			if rv == nil {
				return base.Tri_UNKNOWN
			}
			return PeepholeFoldConstants_TryAbstractEqualityComparison(
				peepholeOptimization, left, NodeUtil_NumberNode(*rv, right))
		}
		if (leftValueType == NodeUtil_ValueType_STRING && rightValueType == NodeUtil_ValueType_NUMBER) ||
			leftValueType == NodeUtil_ValueType_BOOLEAN {
			lv := peepholeOptimization.GetSideEffectFreeNumberValue(left)
			if lv == nil {
				return base.Tri_UNKNOWN
			}
			return PeepholeFoldConstants_TryAbstractEqualityComparison(
				peepholeOptimization, NodeUtil_NumberNode(*lv, left), right)
		}

		if leftValueType == NodeUtil_ValueType_BIGINT || rightValueType == NodeUtil_ValueType_BIGINT {
			lv := peepholeOptimization.GetSideEffectFreeBigIntValue(left)
			rv := peepholeOptimization.GetSideEffectFreeBigIntValue(right)
			if lv != nil && rv != nil {
				return base.Tri_ForBoolean(jutil.BigIntegerEquals(lv, rv))
			}
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
	peepholeOptimization AbstractPeepholeOptimization, left *rhino.Node, right *rhino.Node) base.Tri {
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
			if NodeUtil_IsNaN(left) {
				return base.Tri_FALSE
			}
			if NodeUtil_IsNaN(right) {
				return base.Tri_FALSE
			}
			lv := peepholeOptimization.GetSideEffectFreeNumberValue(left)
			rv := peepholeOptimization.GetSideEffectFreeNumberValue(right)
			if lv != nil && rv != nil {
				return base.Tri_ForBoolean(*lv == *rv)
			}
		case NodeUtil_ValueType_STRING:
			lv := peepholeOptimization.GetSideEffectFreeStringValue(left)
			rv := peepholeOptimization.GetSideEffectFreeStringValue(right)
			if lv != nil && rv != nil {
				// In JS, browsers parse \v differently. So do not consider strings
				// equal if one contains \v.
				if jutil.StringIndexOfChar(*lv, 0x000B) != -1 || jutil.StringIndexOfChar(*rv, 0x000B) != -1 {
					return base.Tri_UNKNOWN
				} else {
					if *lv == *rv {
						return base.Tri_TRUE
					}
					return base.Tri_FALSE
				}
			} else if left.IsTypeOf() &&
				right.IsTypeOf() &&
				left.GetFirstChild().IsName() &&
				right.GetFirstChild().IsName() &&
				left.GetFirstChild().GetString() == right.GetFirstChild().GetString() {
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
	peepholeOptimization AbstractPeepholeOptimization, op rhino.Token, left *rhino.Node, right *rhino.Node) base.Tri {
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

// TryFoldCtorCall tries to fold away unnecessary object instantiation. e.g.
// this[new String('eval')] -> this.eval
//
// Java: PeepholeFoldConstants.tryFoldCtorCall(Node)
func (p *PeepholeFoldConstants) TryFoldCtorCall(n *rhino.Node) *rhino.Node {
	jutil.CheckArgument(n.IsNew())

	// we can remove this for GETELEM calls (anywhere else?)
	if PeepholeFoldConstants_InForcedStringContext(n) {
		return p.TryFoldInForcedStringContext(n)
	}
	return n
}

// TryFoldUselessObjectDotDefinePropertiesCall removes useless calls:
// Object.defineProperties(o, {}) -> o
//
// Java: PeepholeFoldConstants.tryFoldUselessObjectDotDefinePropertiesCall(Node)
func (p *PeepholeFoldConstants) TryFoldUselessObjectDotDefinePropertiesCall(n *rhino.Node) *rhino.Node {
	jutil.CheckArgument(n.IsCall() || n.IsOptChainCall())

	if NodeUtil_IsObjectDefinePropertiesDefinition(n) {
		srcObj := n.GetLastChild()
		if srcObj.IsObjectLit() && !srcObj.HasChildren() {
			parent := n.GetParent()
			destObj := n.GetSecondChild().Detach()
			n.ReplaceWith(destObj)
			p.ReportChangeToEnclosingScope(parent)
		}
	}
	return n
}

// PeepholeFoldConstants_InForcedStringContext returns whether this node must
// be coerced to a string.
//
// Java: PeepholeFoldConstants.inForcedStringContext(Node)
func PeepholeFoldConstants_InForcedStringContext(n *rhino.Node) bool {
	if n.GetParent().IsGetElem() && n.GetParent().GetLastChild() == n {
		return true
	}

	// we can fold in the case "" + new String("")
	return n.GetParent().IsAdd()
}

// Java: PeepholeFoldConstants.tryFoldInForcedStringContext(Node)
func (p *PeepholeFoldConstants) TryFoldInForcedStringContext(n *rhino.Node) *rhino.Node {
	// For now, we only know how to fold ctors.
	jutil.CheckArgument(n.IsNew())

	objectType := n.GetFirstChild()
	if !objectType.IsName() {
		return n
	}

	if objectType.GetString() == "String" {
		value := objectType.GetNext()
		var stringValue *string
		if value == nil {
			empty := ""
			stringValue = &empty
		} else {
			stringValue = p.GetSideEffectFreeStringValue(value)
		}

		if stringValue == nil {
			return n
		}

		parent := n.GetParent()
		newString := rhino.IR_String(*stringValue)

		n.ReplaceWith(newString)
		newString.SrcrefIfMissing(parent)
		p.ReportChangeToEnclosingScope(parent)

		return newString
	}
	return n
}

// TryFoldGetElem tries, for element access using GETLEM/OPTCHAIN_GETELEM on
// object literals, arrays or strings, to fold the prop access. e.g. folds
// array-element [1, 2, 3][1];
//
// Java: PeepholeFoldConstants.tryFoldGetElem(Node, Node, Node)
func (p *PeepholeFoldConstants) TryFoldGetElem(n *rhino.Node, left *rhino.Node, right *rhino.Node) *rhino.Node {
	jutil.CheckArgument(n.IsGetElem() || n.IsOptChainGetElem())

	if left.IsObjectLit() {
		if right.IsStringLit() {
			return p.TryFoldObjectPropAccess(n, left, right.GetString())
		}
	} else if left.IsArrayLit() {
		return p.TryFoldArrayAccess(n, left, right)
	} else if left.IsStringLit() {
		return p.TryFoldStringArrayAccess(n, left, right)
	}
	return n
}

// TryFoldGetProp tries, for prop access using GETPROP/OPTCHAIN_GETPROP on
// object literals, to fold their property access. For prop access on arrays,
// only tries to fold array-length. e.g [1, 2, 3].length ==> 3, [x, y].length
// ==> 2
//
// Java: PeepholeFoldConstants.tryFoldGetProp(Node)
func (p *PeepholeFoldConstants) TryFoldGetProp(n *rhino.Node) *rhino.Node {
	jutil.CheckArgument(n.IsGetProp() || n.IsOptChainGetProp())
	left := n.GetFirstChild()
	name := n.GetString()

	if left.IsObjectLit() {
		return p.TryFoldObjectPropAccess(n, left, name)
	}

	if name == "length" {
		knownLength := -1
		switch left.GetToken() {
		case rhino.Token_ARRAYLIT:
			if p.MayHaveSideEffects(left) {
				// Nope, can't fold this, without handling the side-effects.
				return n
			}
			knownLength = left.GetChildCount()
		case rhino.Token_STRINGLIT:
			knownLength = jutil.StringLength(left.GetString())
		default:
			// Not a foldable case, forget it.
			return n
		}

		jutil.CheckState(knownLength != -1)
		lengthNode := rhino.IR_Number(float64(knownLength))
		p.ReportChangeToEnclosingScope(n)
		n.ReplaceWith(lengthNode)

		return lengthNode
	}

	return n
}

// Java: PeepholeFoldConstants.tryFoldArrayAccess(Node, Node, Node)
func (p *PeepholeFoldConstants) TryFoldArrayAccess(n *rhino.Node, left *rhino.Node, right *rhino.Node) *rhino.Node {
	// If GETPROP/GETELEM is used as assignment target the array literal is
	// acting as a temporary we can't fold it here:
	//    "[][0] += 1"
	if NodeUtil_IsLValue(n) {
		return n
	}

	index := p.GetSideEffectFreeNumberValue(right)
	if !p.IsReasonableDoubleValue(index) {
		// Sometimes people like to use complex expressions to index into
		// arrays, or strings to index into array methods.
		return n
	}

	if !base.JSCompDoubles_IsExactInt32(*index) {
		// Ideally this should be caught in the check passes.
		p.Report(PeepholeFoldConstants_INVALID_GETELEM_INDEX_ERROR, right)
		return n
	}

	intIndex := int(jutil.D2I(*index))
	var current *rhino.Node
	if intIndex >= 0 {
		current = left.GetFirstChild()
	}
	var elem *rhino.Node
	for i := 0; current != nil; i++ {
		if current.IsSpread() {
			// The only time we can fold getelems with spread is for spread arrays literals, and
			// `tryFlattenArray` already flattens those.
			return n
		}
		if i != intIndex {
			if p.MayHaveSideEffects(current) {
				return n
			}
		} else {
			elem = current
		}

		current = current.GetNext()
	}

	if elem == nil { // If the index was out of bounds
		elem = NodeUtil_NewUndefinedNode(left)
	} else if elem.IsEmpty() {
		elem = NodeUtil_NewUndefinedNode(elem)
	} else {
		elem.Detach()
	}

	// Replace the entire GETELEM with the value
	n.ReplaceWith(elem)
	p.ReportChangeToEnclosingScope(elem)
	return elem
}

// TryFoldSpread folds any occurrences of spread of array literals e.g.
// ...[1,2,3] to 1,2,3
//
// Java: PeepholeFoldConstants.tryFoldSpread(Node)
func (p *PeepholeFoldConstants) TryFoldSpread(spread *rhino.Node) *rhino.Node {
	jutil.CheckState(spread.IsSpread())
	parent := spread.GetParent()
	child := spread.GetOnlyChild()
	if child.IsArrayLit() {
		for n := child.GetFirstChild(); n != nil; n = n.GetNext() {
			if n.IsEmpty() {
				// Do not fold if the array literal has any holes as it's a syntax error to have holes
				// (empty arg) if the spread happens to be a call arg
				return spread
			}
		}
		parent.AddChildrenAfter(child.RemoveChildren(), spread)
		spread.Detach()
		p.ReportChangeToEnclosingScope(parent)
	}
	return parent
}

// TryFlattenArrayOrObjectLit flattens array- or object-literals that contain
// spreads of other literals.
//
// Does not recurse into nested spreads because this method is already called
// as part of a postorder traversal and nested spreads will already have been
// flattened.
//
// Example: `[0, ...[1, 2, 3], 4, ...[5]]` => `[0, 1, 2, 3, 4, 5]`
//
// Java: PeepholeFoldConstants.tryFlattenArrayOrObjectLit(Node)
func (p *PeepholeFoldConstants) TryFlattenArrayOrObjectLit(parentLit *rhino.Node) *rhino.Node {
	for child := parentLit.GetFirstChild(); child != nil; {
		// We have to store the next element here because nodes may be inserted below.
		spread := child
		child = child.GetNext()

		if !spread.IsSpread() {
			continue
		}

		innerLit := spread.GetOnlyChild()
		if parentLit.GetToken() != innerLit.GetToken() {
			continue // We only want to inline arrays into arrays and objects into objects.
		}

		if parentLit.IsObjectLit() && !PeepholeFoldConstants_IsSafeToFlattenObjectLitSpread(innerLit) {
			continue
		}

		parentLit.AddChildrenAfter(innerLit.RemoveChildren(), spread)
		spread.Detach()
		p.ReportChangeToEnclosingScope(parentLit)
	}
	return parentLit
}

// Java: PeepholeFoldConstants.tryFoldStringArrayAccess(Node, Node, Node)
func (p *PeepholeFoldConstants) TryFoldStringArrayAccess(n *rhino.Node, left *rhino.Node, right *rhino.Node) *rhino.Node {
	// If GETPROP/GETELEM is used as assignment target the array literal is
	// acting as a temporary we can't fold it here:
	//    "[][0] += 1"
	if NodeUtil_IsLValue(n) {
		return n
	}

	index := p.GetSideEffectFreeNumberValue(right)
	if !p.IsReasonableDoubleValue(index) {
		// Sometimes people like to use complex expressions to index into
		// arrays, or strings to index into array methods.
		return n
	}

	intIndex := jutil.D2I(*index)
	if float64(intIndex) != *index {
		p.Report(PeepholeFoldConstants_INVALID_GETELEM_INDEX_ERROR, right)
		return n
	}

	jutil.CheckState(left.IsStringLit())
	value := left.GetString()
	if intIndex < 0 || int(intIndex) >= jutil.StringLength(value) {
		undefined := NodeUtil_NewUndefinedNode(left)
		n.ReplaceWith(undefined)
		p.ReportChangeToEnclosingScope(undefined)
		return undefined
	}

	var c uint16
	// Note: For now skip the strings with unicode
	// characters as I don't understand the differences
	// between Java and JavaScript.
	for i := 0; i <= int(intIndex); i++ {
		c = jutil.StringCharAt(value, i)
		if c < 32 || c > 127 {
			return n
		}
	}
	elem := rhino.IR_String(jutil.StringFromChar(c))

	// Replace the entire GETELEM with the value
	n.ReplaceWith(elem)
	p.ReportChangeToEnclosingScope(elem)
	return elem
}

// TryFoldObjectPropAccess tries to fold the node `n` that's accessing an
// object literal's property. Bails out (skips folding and returns the same
// `n` node) with certainty when any of the following holds true:
//   - the access `n` is L-value
//   - the property references super
//   - the object has side-effects other than those preserved by folding the
//     property
//   - property accessed does not exist on the object (might exist on
//     prototype)
//
// Examples of folding:
//   - `({a() { return 1; }})?.a` ---> `(function() { return 1; })`
//   - `({a() { return 1; }}).a()` ---> `(function() { return 1; }())`
//
// Java: PeepholeFoldConstants.tryFoldObjectPropAccess(Node, Node, String)
func (p *PeepholeFoldConstants) TryFoldObjectPropAccess(n *rhino.Node, left *rhino.Node, name string) *rhino.Node {
	jutil.CheckArgument(NodeUtil_IsNormalOrOptChainGet(n))

	if !left.IsObjectLit() {
		return n
	}

	if NodeUtil_IsLValue(n) {
		// If GETPROP/GETELEM is used as assignment target the object literal is
		// acting as a temporary we can't fold it here:
		//    "{a:x}.a += 1" is not "x += 1"
		jutil.CheckState(!NodeUtil_IsOptChainNode(n)) // optional chains can not be targets of assign
		return n
	}

	// find the last definition in the object literal
	var key *rhino.Node
	var value *rhino.Node
	for c := left.GetFirstChild(); c != nil; c = c.GetNext() {
		switch c.GetToken() {
		case rhino.Token_SETTER_DEF:
		case rhino.Token_OBJECT_SPREAD:
			// Reset the search because spread could overwrite any previous result.
			key = nil
			value = nil
		case rhino.Token_COMPUTED_PROP:
			// don't handle computed properties unless the input is a simple string
			prop := c.GetFirstChild()
			if !prop.IsStringLit() {
				return n
			}
			if prop.GetString() == name {
				key = c
				value = c.GetSecondChild()
			}
		case rhino.Token_GETTER_DEF, rhino.Token_STRING_KEY, rhino.Token_MEMBER_FUNCTION_DEF:
			if c.GetString() == name {
				key = c
				value = key.GetFirstChild()
			}
		default:
			panic(jutil.NewIllegalStateException())
		}
	}

	// Didn't find a definition of the name in the object literal, it might
	// be coming from the Object prototype.
	if value == nil {
		return n
	}

	/* `super` captures a hidden reference to the declaring objectlit, so we can't fold it away. */
	if NodeUtil_ReferencesSuper(value) {
		return n
	}

	/*
	 * Check to see if there are any side-effects to this object-literal.
	 *
	 * <p>We remove the value we're going to use because its side-effects will be preserved.
	 */
	tempValue := rhino.IR_NullNode()
	value.ReplaceWith(tempValue)
	hasSideEffectBesidesValue := p.MayHaveSideEffects(left)
	tempValue.ReplaceWith(value)
	if hasSideEffectBesidesValue {
		return n
	}

	keyIsGetter := NodeUtil_IsGetOrSetKey(key)
	nIsInvoked := NodeUtil_IsInvocationTarget(n)

	if keyIsGetter || nIsInvoked {
		// When the code looks like:
		//   {x: f}.x();
		//   {get x() {...}}.x;
		//   {get ['x']() {...}}.x;

		/*
		 * It's not safe, in general, to convert that to just a function call, because the receiver
		 * value will be wrong.
		 *
		 * <p>However, there are some cases where it's ok which we check here.
		 */
		if !value.IsFunction() || NodeUtil_ReferencesOwnReceiver(value) {
			return n
		}
	}

	value.Detach()

	parent := n.GetParent()

	if NodeUtil_IsOptChainNode(parent) {

		/*
		 * If the chain continues after `n`, simply doing `n.replaceWith(value)` below would leave the
		 * subsequent nodes in the current chain segment optional, with their start `n` replaced.
		 *
		 * <p>So, we must ensure that all nodes in the chain's current segment are made non-optional.
		 *
		 * <p>This can happen for e.g.
		 *
		 * <ul>
		 *   <li>`({a() { return 1; }})?.a()` ---> `(function() { return 1; })()`. Here, parent
		 *       OPTCHAIN_CALL node must be converted to CALL.
		 *   <li>`({a() { return 1; }})?.a().b.c?.d` ---> `(function() { return 1; })().b.c?.d`. Here,
		 *       all nodes upto `({a() { return 1; }})?.a().b.c` must become non-optional
		 * </ul>
		 */
		endOfCurrentChain := NodeUtil_GetEndOfOptChainSegment(parent)
		NodeUtil_ConvertToNonOptionalChainSegment(endOfCurrentChain)
	}
	if keyIsGetter {
		value = rhino.IR_Call(value)
		value.PutBooleanProp(rhino.Node_FREE_CALL, true)
	} else if nIsInvoked {
		n.GetParent().PutBooleanProp(rhino.Node_FREE_CALL, true)
	}

	n.ReplaceWith(value)
	p.ReportChangeToEnclosingScope(value)
	p.MarkFunctionsDeleted(n)
	return n
}

// Java: PeepholeFoldConstants.bigintNode(BigInteger, Node)
func PeepholeFoldConstants_BigintNode(val *big.Int, srcref *rhino.Node) *rhino.Node {
	var tree *rhino.Node
	if val.Sign() < 0 {
		tree = rhino.IR_Neg(rhino.IR_Bigint(new(big.Int).Neg(val)))
	} else {
		tree = rhino.IR_Bigint(val)
	}
	return tree.SrcrefTree(srcref)
}

// Java: PeepholeFoldConstants.isSafeToFlattenObjectLitSpread(Node)
func PeepholeFoldConstants_IsSafeToFlattenObjectLitSpread(innerLit *rhino.Node) bool {
	for prop := innerLit.GetFirstChild(); prop != nil; prop = prop.GetNext() {
		// Object spread converts all properties to regular data properties. This step evaluates
		// any getters and omits any setters. See https://tc39.es/ecma262/#sec-copydataproperties.
		// For simplicity, we just back off and do not flatten spread in this case.
		if NodeUtil_IsGetOrSetKey(prop) {
			return false
		}
	}
	return true
}
