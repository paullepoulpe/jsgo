// Port of com/google/javascript/jscomp/TypeInference.java @ closure-compiler v20261021

/*
 * Copyright 2008 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Java: TypeInference.traverseHook(Node, FlowScope)
func (ti *TypeInference) TraverseHook(n rhino.Node, scope *LinkedFlowScope) *LinkedFlowScope {
	hook := ast.MustHook(n)
	condition := hook.Cond().Node()
	trueNode := hook.Then().Node()
	falseNode := hook.Else().Node()

	// verify the condition
	scope = ti.Traverse(condition, scope)

	// reverse abstract interpret the condition to produce two new scopes
	trueScope :=
		ti.reverseInterpreter.GetPreciserScopeKnowingConditionOutcome(condition, scope, rhino.Outcome_TRUE)
	falseScope :=
		ti.reverseInterpreter.GetPreciserScopeKnowingConditionOutcome(condition, scope, rhino.Outcome_FALSE)

	// traverse the true node with the trueScope
	ti.Traverse(trueNode, trueScope)

	// traverse the false node with the falseScope
	ti.Traverse(falseNode, falseScope)

	// meet true and false nodes' types and assign
	trueType := trueNode.GetJSType()
	falseType := falseNode.GetJSType()
	if trueType != nil && falseType != nil {
		n.SetJSType(trueType.GetLeastSupertype(falseType))
	} else {
		n.SetJSType(ti.unknownType)
	}

	return scope
}

// Java: TypeInference.traverseNullishCoalesce(Node, FlowScope)
func (ti *TypeInference) TraverseNullishCoalesce(n rhino.Node, scope *LinkedFlowScope) *LinkedFlowScope {
	jutil.CheckArgument(n.IsNullishCoalesce() || n.IsAssignNullishCoalesce())
	left := n.GetFirstChild()
	right := n.GetLastChild()

	scope = ti.Traverse(left, scope)

	rightScope :=
		ti.reverseInterpreter.GetPreciserScopeKnowingConditionOutcome(left, scope, rhino.Outcome_NULLISH)

	scopeAfterTraverseRight := ti.Traverse(right, rightScope)

	leftType := left.GetJSType()
	rightType := right.GetJSType()

	var type_ rhino.JSType = ti.unknownType
	if leftType != nil {
		if !leftType.IsNullable() && !leftType.IsVoidable() {
			type_ = leftType
		} else if rightType != nil {
			type_ = ti.registry.CreateUnionType(leftType.RestrictByNotNullOrUndefined(), rightType)
			scope = ti.Join(scope, scopeAfterTraverseRight)
			// Assignment occurs if lhs is null
			if n.IsAssignNullishCoalesce() {
				scope = ti.UpdateScopeForAssignment(scope, left, type_, TypeInference_AssignmentType_ASSIGN)
			}
		}
	}
	n.SetJSType(type_)
	return scope
}

// Java: TypeInference.traverseAnd(Node, FlowScope)
func (ti *TypeInference) TraverseAnd(n rhino.Node, scope *LinkedFlowScope) TypeInference_BooleanOutcomePair {
	left := ast.MustBinary(n).Left().Node()
	leftOutcome := ti.TraverseWithinShortCircuitingBinOp(left, scope)
	return ti.TraverseShortCircuitingBinOp(n, left, leftOutcome)
}

// Java: TypeInference.traverseOr(Node, FlowScope)
func (ti *TypeInference) TraverseOr(n rhino.Node, scope *LinkedFlowScope) TypeInference_BooleanOutcomePair {
	left := ast.MustBinary(n).Left().Node()
	leftOutcome := ti.TraverseWithinShortCircuitingBinOp(left, scope)
	return ti.TraverseShortCircuitingBinOp(n, left, leftOutcome)
}

// Java: TypeInference.traverseShortCircuitingBinOpAssignment(Node, FlowScope)
func (ti *TypeInference) TraverseShortCircuitingBinOpAssignment(n rhino.Node, scope *LinkedFlowScope) *LinkedFlowScope {
	left := ast.MustAssignOp(n).Target().Node()
	nIsAnd := n.IsAssignAnd()
	leftOutcome := ti.TraverseWithinShortCircuitingBinOp(left, scope)
	outcome := ti.TraverseShortCircuitingBinOp(n, left, leftOutcome)

	outcomeJoinedFlowScope := outcome.GetJoinedFlowScope()

	if leftOutcome.toBooleanOutcomes == rhino.BooleanLiteralSet_Get(!nIsAnd) {
		// Either n is && and lhs has a toBooleanOutcome of false,
		// or n is || and lhs has a toBooleanOutcome of true, so assignment does not occur
		// The scope is otherwise updated according to the result of an assignment.
		return outcomeJoinedFlowScope
	}
	return ti.UpdateScopeForAssignment(
		outcomeJoinedFlowScope, left, n.GetJSType(), TypeInference_AssignmentType_ASSIGN)
}

// Java: TypeInference.traverseShortCircuitingBinOp(Node, Node, BooleanOutcomePair)
func (ti *TypeInference) TraverseShortCircuitingBinOp(
	n rhino.Node, left rhino.Node, leftOutcome TypeInference_BooleanOutcomePair) TypeInference_BooleanOutcomePair {
	jutil.CheckArgument(n.IsAnd() || n.IsOr() || n.IsAssignAnd() || n.IsAssignOr())
	nIsAnd := n.IsAnd() || n.IsAssignAnd()
	right := n.GetLastChild()

	// type the left node
	leftType := left.GetJSType()

	// reverse abstract interpret the left node to produce the correct
	// scope in which to verify the right node
	rightScope :=
		ti.reverseInterpreter.GetPreciserScopeKnowingConditionOutcome(
			left,
			leftOutcome.GetOutcomeFlowScope(left.GetToken(), nIsAnd),
			rhino.Outcome_ForBoolean(nIsAnd))

	// type the right node
	rightOutcome := ti.TraverseWithinShortCircuitingBinOp(right, rightScope)
	rightType := right.GetJSType()

	var type_ rhino.JSType
	var outcome TypeInference_BooleanOutcomePair
	if leftType != nil && rightType != nil {
		leftType = leftType.GetRestrictedTypeGivenOutcome(rhino.Outcome_ForBoolean(!nIsAnd))
		if leftOutcome.toBooleanOutcomes == rhino.BooleanLiteralSet_Get(!nIsAnd) {
			// Either n is && and lhs is false, or n is || and lhs is true.
			// Use the restricted left type; the right side never gets evaluated.
			type_ = leftType
			outcome = leftOutcome
		} else {
			// Use the join of the restricted left type knowing the outcome of the
			// ToBoolean predicate and of the right type.
			type_ = leftType.GetLeastSupertype(rightType)
			outcome =
				NewTypeInference_BooleanOutcomePair(
					ti,
					TypeInference_JoinBooleanOutcomes(
						nIsAnd, leftOutcome.toBooleanOutcomes, rightOutcome.toBooleanOutcomes),
					TypeInference_JoinBooleanOutcomes(nIsAnd, leftOutcome.booleanValues, rightOutcome.booleanValues),
					leftOutcome.GetJoinedFlowScope(),
					rightOutcome.GetJoinedFlowScope())
		}
		// Exclude the boolean type if the literal set is empty because a boolean
		// can never actually be returned.
		if outcome.booleanValues == rhino.BooleanLiteralSet_EMPTY &&
			ti.GetNativeType(rhino.JSTypeNative_BOOLEAN_TYPE).IsSubtypeOf(type_) {
			// Exclusion only makes sense for a union type.
			if type_.IsUnionType() {
				type_ = type_.ToMaybeUnionType().GetRestrictedUnion(ti.GetNativeType(rhino.JSTypeNative_BOOLEAN_TYPE))
			}
		}
	} else {
		type_ = ti.unknownType
		outcome =
			NewTypeInference_BooleanOutcomePair(
				ti,
				rhino.BooleanLiteralSet_BOTH,
				rhino.BooleanLiteralSet_BOTH,
				leftOutcome.GetJoinedFlowScope(),
				rightOutcome.GetJoinedFlowScope())
	}
	n.SetJSType(type_)
	return outcome
}

// Java: TypeInference.traverseWithinShortCircuitingBinOp(Node, FlowScope)
func (ti *TypeInference) TraverseWithinShortCircuitingBinOp(n rhino.Node, scope *LinkedFlowScope) TypeInference_BooleanOutcomePair {
	switch n.GetToken() {
	case rhino.Token_AND:
		return ti.TraverseAnd(n, scope)
	case rhino.Token_OR:
		return ti.TraverseOr(n, scope)
	default:
		scope = ti.Traverse(n, scope)
		return ti.NewBooleanOutcomePair(n.GetJSType(), scope)
	}
}

// Java: TypeInference.joinBooleanOutcomes(boolean, BooleanLiteralSet, BooleanLiteralSet)
func TypeInference_JoinBooleanOutcomes(
	isAnd bool, left rhino.BooleanLiteralSet, right rhino.BooleanLiteralSet) rhino.BooleanLiteralSet {
	// A truthy value on the lhs of an {@code &&} can never make it to the
	// result. Same for a falsy value on the lhs of an {@code ||}.
	// Hence the intersection.
	return right.Union(left.Intersection(rhino.BooleanLiteralSet_Get(!isAnd)))
}

// TypeInference_BooleanOutcomePair: when traversing short-circuiting binary
// operations, we need to keep track of two sets of boolean literals: 1.
// toBooleanOutcomes: boolean literals as converted from any types, 2.
// booleanValues: boolean literals from just boolean types.
type TypeInference_BooleanOutcomePair struct {
	outer *TypeInference

	toBooleanOutcomes rhino.BooleanLiteralSet
	booleanValues     rhino.BooleanLiteralSet

	// The scope if only half of the expression executed, when applicable.
	leftScope *LinkedFlowScope

	// The scope when the whole expression executed.
	rightScope *LinkedFlowScope

	// The scope when we don't know how much of the expression is executed.
	joinedScope *LinkedFlowScope
}

// Java: TypeInference.BooleanOutcomePair(BooleanLiteralSet, BooleanLiteralSet, FlowScope, FlowScope)
func NewTypeInference_BooleanOutcomePair(
	outer *TypeInference,
	toBooleanOutcomes rhino.BooleanLiteralSet,
	booleanValues rhino.BooleanLiteralSet,
	leftScope *LinkedFlowScope,
	rightScope *LinkedFlowScope) TypeInference_BooleanOutcomePair {
	return TypeInference_BooleanOutcomePair{
		outer:             outer,
		toBooleanOutcomes: toBooleanOutcomes,
		booleanValues:     booleanValues,
		leftScope:         leftScope,
		rightScope:        rightScope,
	}
}

// GetJoinedFlowScope gets the safe estimated scope without knowing if all of
// the subexpressions will be evaluated.
//
// Java: TypeInference.BooleanOutcomePair.getJoinedFlowScope()
func (p *TypeInference_BooleanOutcomePair) GetJoinedFlowScope() *LinkedFlowScope {
	if p.joinedScope == nil {
		if p.leftScope == p.rightScope {
			p.joinedScope = p.rightScope
		} else {
			p.joinedScope = p.outer.Join(p.leftScope, p.rightScope)
		}
	}
	return p.joinedScope
}

// GetOutcomeFlowScope gets the outcome scope if we do know the outcome of
// the entire expression.
//
// Java: TypeInference.BooleanOutcomePair.getOutcomeFlowScope(Token, boolean)
func (p *TypeInference_BooleanOutcomePair) GetOutcomeFlowScope(nodeType rhino.Token, outcome bool) *LinkedFlowScope {
	if (nodeType == rhino.Token_AND && outcome) || (nodeType == rhino.Token_OR && !outcome) {
		// We know that the whole expression must have executed.
		return p.rightScope
	} else {
		return p.GetJoinedFlowScope()
	}
}

// Java: TypeInference.newBooleanOutcomePair(JSType, FlowScope)
func (ti *TypeInference) NewBooleanOutcomePair(jsType rhino.JSType, flowScope *LinkedFlowScope) TypeInference_BooleanOutcomePair {
	if jsType == nil {
		return NewTypeInference_BooleanOutcomePair(
			ti, rhino.BooleanLiteralSet_BOTH, rhino.BooleanLiteralSet_BOTH, flowScope, flowScope)
	}
	var booleanValues rhino.BooleanLiteralSet
	if ti.registry.GetNativeType(rhino.JSTypeNative_BOOLEAN_TYPE).IsSubtypeOf(jsType) {
		booleanValues = rhino.BooleanLiteralSet_BOTH
	} else {
		booleanValues = rhino.BooleanLiteralSet_EMPTY
	}
	return NewTypeInference_BooleanOutcomePair(
		ti,
		jsType.GetPossibleToBooleanOutcomes(),
		booleanValues,
		flowScope,
		flowScope)
}

// SetOptChainNodeTypeAfterChildrenTraversed sets the appropriate type on the
// OptChain node `n` after its children have been traversed.
//
// Java: TypeInference.setOptChainNodeTypeAfterChildrenTraversed(Node, FlowScope)
func (ti *TypeInference) SetOptChainNodeTypeAfterChildrenTraversed(
	n rhino.Node, scopeAfterChildren *LinkedFlowScope) *LinkedFlowScope {
	switch n.GetToken() {
	case rhino.Token_OPTCHAIN_GETPROP:
		return ti.SetGetPropNodeTypeAfterChildrenTraversed(n, scopeAfterChildren)
	case rhino.Token_OPTCHAIN_GETELEM:
		return ti.SetGetElemNodeTypeAfterChildrenTraversed(n, scopeAfterChildren)
	case rhino.Token_OPTCHAIN_CALL:
		return ti.SetCallNodeTypeAfterChildrenTraversed(n, scopeAfterChildren)
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Illegal token inside finishTraversingOptChain"))
	}
}

// TypeInference_OptChainInfo holds flow scopes for conditional and
// unconditional parts during traversal of an optional chain.
type TypeInference_OptChainInfo struct {
	endOfChain         rhino.Node
	startOfChain       rhino.Node
	unconditionalScope *LinkedFlowScope
}

// Java: TypeInference.OptChainInfo(Node, Node)
func NewTypeInference_OptChainInfo(endOfChain rhino.Node, startOfChain rhino.Node) TypeInference_OptChainInfo {
	return TypeInference_OptChainInfo{
		endOfChain:   endOfChain,
		startOfChain: startOfChain,
	}
}

// TraverseOptChain: traversal requires holding scopes from the
// unconditional start (lhs of the `?.`) and the conditional part (rhs of the
// `?.`), and using the startType to determine the resulting scope.
//
// Java: TypeInference.traverseOptChain(Node, FlowScope)
func (ti *TypeInference) TraverseOptChain(n rhino.Node, scope *LinkedFlowScope) *LinkedFlowScope {
	jutil.CheckArgument(NodeUtil_IsOptChainNode(n))

	if NodeUtil_IsEndOfOptChainSegment(n) {
		// Create new optional chain tracking object and push it onto the stack.
		startOfChain := NodeUtil_GetStartOfOptChainSegment(n)
		optChainInfo := NewTypeInference_OptChainInfo(n, startOfChain)
		// optChainArrayDeque.addFirst(optChainInfo)
		ti.optChainStack = append(ti.optChainStack, optChainInfo)
	}

	lhsScope := ti.Traverse(n.GetFirstChild(), scope)

	if n.IsOptionalChainStart() {
		// Store lhsScope into top-of-stack as unexecuted (unconditional) scope.
		if len(ti.optChainStack) == 0 {
			panic(jutil.NewNullPointerException())
		}
		ti.optChainStack[len(ti.optChainStack)-1].unconditionalScope = lhsScope
	}

	// Traverse the remaining children and capture their changes into a new FlowScope var
	// `aboutToExecuteScope`. This FlowScope must be constructed on top of the `lhsScope` and not
	// the original scope `scope`, otherwise changes to outer variable that are preserved in the
	// lhsScope would not be captured in the `aboutToExecuteScope`.
	aboutToExecuteScope := lhsScope
	nextChild := n.GetSecondChild()
	for !nextChild.IsNil() {
		aboutToExecuteScope = ti.Traverse(nextChild, aboutToExecuteScope)
		nextChild = nextChild.GetNext()
	}

	// Assigns the type to `n` assuming the entire chain executes and returns the the executed
	// scope.
	executedScope := ti.SetOptChainNodeTypeAfterChildrenTraversed(n, aboutToExecuteScope)

	// Unlike CALL, the OPTCHAIN_CALL nodes must not remain untyped when left child is untyped.
	if n.GetJSType() == nil {
		n.SetJSType(ti.unknownType)
	}

	if NodeUtil_IsEndOfOptChainSegment(n) {
		// Use the startNode's type to selectively join the executed scope with the unexecuted scope,
		// and update the type assigned to `n` in `setXAfterChildrenTraversed()`
		startOfChain := NodeUtil_GetStartOfOptChainSegment(n)

		// Pop the stack to obtain the current chain.
		if len(ti.optChainStack) == 0 {
			panic(jutil.NewNoSuchElementException())
		}
		top := len(ti.optChainStack) - 1
		currentChain := ti.optChainStack[top]
		ti.optChainStack[top] = TypeInference_OptChainInfo{}
		ti.optChainStack = ti.optChainStack[:top]

		// Sanity check that the popped chain correctly corresponds to the current optional chain
		jutil.CheckState(currentChain.endOfChain == n)
		jutil.CheckState(currentChain.startOfChain == startOfChain)

		startNode := rhino.CheckNotNull(startOfChain.GetFirstChild())
		return ti.UpdateTypeWhenEndOfOptChain(n, startNode, currentChain, executedScope)
	} else {
		// `n` is just an inner node (i.e. not the end of the current chain). Simply return the
		// executedScope.
		return executedScope
	}
}

// UpdateTypeWhenEndOfOptChain uses the startNode's type to selectively join
// the executed scope with the unexecuted scope, and updates the type assigned
// to `optChain` in `setXAfterChildrenTraversed()`.
//
// Java: TypeInference.updateTypeWhenEndOfOptChain(Node, Node, OptChainInfo, FlowScope)
func (ti *TypeInference) UpdateTypeWhenEndOfOptChain(
	optChain rhino.Node, startNode rhino.Node, currentChain TypeInference_OptChainInfo, executedScope *LinkedFlowScope) *LinkedFlowScope {

	startType := ti.GetJSType(startNode)
	if startType.IsUnknownType() {
		// Unknown startType: Conditional part may execute. Result type must be unknown.
		optChain.SetJSType(ti.unknownType)
		return ti.Join(executedScope, currentChain.unconditionalScope)
	} else if startType.IsNullType() || startType.IsVoidType() {
		// Conditional part will not execute.
		optChain.SetJSType(ti.registry.GetNativeType(rhino.JSTypeNative_VOID_TYPE))
		return currentChain.unconditionalScope
	} else if !startType.IsNullable() {
		// Conditional part will execute; `optChain` was assigned the right type in
		// `setXAfterChildrenTraversed()`.
		return executedScope
	} else {
		// Nullable start: Conditional part may execute. Need to add a VOID_TYPE to the type
		// assigned to `optChain` within `setXAfterChildrenTraversed()`.
		optChain.SetJSType(
			ti.registry.CreateUnionType(ti.registry.GetNativeType(rhino.JSTypeNative_VOID_TYPE), optChain.GetJSType()))
		unexecutedScope :=
			ti.reverseInterpreter.GetPreciserScopeKnowingConditionOutcome(
				startNode, currentChain.unconditionalScope, rhino.Outcome_NULLISH)
		return ti.Join(executedScope, unexecutedScope)
	}
}

// TypeInference_InferPropertyTypesToMatchConstraint: suppose X is an object
// with inferred properties. Suppose also that X is used in a way where it
// would only type-check correctly if some of those properties are widened.
// Then we should be polite and automatically widen X's properties.
//
// For a concrete example, consider: param x {{prop: (number|undefined)}}
// function f(x) {} f({});
//
// If we give the anonymous object an inferred property of
// (number|undefined), then this code will type-check appropriately.
//
// Java: TypeInference.inferPropertyTypesToMatchConstraint(JSType, JSType)
func TypeInference_InferPropertyTypesToMatchConstraint(type_ rhino.JSType, constraint rhino.JSType) {
	if type_ == nil || constraint == nil {
		return
	}

	type_.MatchConstraint(constraint)
}

// TightenTypeAfterDereference: if we access a property of a symbol, then that
// symbol is not null or undefined.
//
// Java: TypeInference.tightenTypeAfterDereference(Node, FlowScope)
func (ti *TypeInference) TightenTypeAfterDereference(n rhino.Node, scope *LinkedFlowScope) *LinkedFlowScope {
	switch n.GetToken() {
	case rhino.Token_NAME:
		name := n.GetString()
		if name == "" {
			return scope
		}
		type_ := ti.GetJSType(n)
		narrowed := type_.RestrictByNotNullOrUndefined()
		if type_ != narrowed && !type_.Equals(narrowed) {
			if narrowed == nil {
				narrowed = ti.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
			}
			var declaredVar *TypedVar
			if scope.syntacticScope != nil {
				declaredVar = scope.syntacticScope.GetVar(name)
			}
			if !ti.IsUnflowable(declaredVar) {
				scope = scope.InferSlotType(name, narrowed)
			}
		}
	case rhino.Token_GETPROP:
		type_ := ti.GetJSType(n)
		narrowed := type_.RestrictByNotNullOrUndefined()
		if type_ != narrowed && !type_.Equals(narrowed) {
			if qName, ok := n.GetQualifiedNameString(); ok {
				scope = scope.InferQualifiedSlot(n, qName, type_, narrowed, false)
			}
		}
	case rhino.Token_SUPER:
		type_ := ti.GetJSType(n)
		narrowed := type_.RestrictByNotNullOrUndefined()
		if type_ != narrowed && !type_.Equals(narrowed) {
			scope = ti.RedeclareSimpleVar(scope, n, narrowed)
		}
	}
	return scope
}
