// Port of com/google/javascript/jscomp/type/SemanticReverseAbstractInterpreter.java @ closure-compiler v20261015

/*
 * Copyright 2007 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// SemanticReverseAbstractInterpreter is a reverse abstract interpreter using
// the semantics of the JavaScript language as a means to reverse interpret
// computations. This interpreter expects the parse tree inputs to be typed.
type SemanticReverseAbstractInterpreter struct {
	ChainableReverseAbstractInterpreterBase
	trackingUsed []*SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope
	trackingFree []*SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope
	trackingSlab jutil.Slab[SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope]
}

// SemanticReverseAbstractInterpreter_Eq is the merging function for
// equality between types.
//
// Java: SemanticReverseAbstractInterpreter.eq(TypePair)
func SemanticReverseAbstractInterpreter_Eq(p *rhino.JSType_TypePair) *rhino.JSType_TypePair {
	if p.TypeA == nil || p.TypeB == nil {
		return nil
	}
	return p.TypeA.GetTypesUnderEquality(p.TypeB)
}

// SemanticReverseAbstractInterpreter_Ne is the merging function for
// non-equality between types.
//
// Java: SemanticReverseAbstractInterpreter.ne(TypePair)
func SemanticReverseAbstractInterpreter_Ne(p *rhino.JSType_TypePair) *rhino.JSType_TypePair {
	if p.TypeA == nil || p.TypeB == nil {
		return nil
	}
	return p.TypeA.GetTypesUnderInequality(p.TypeB)
}

// SemanticReverseAbstractInterpreter_Sheq is the merging function for strict
// equality between types.
//
// Java: SemanticReverseAbstractInterpreter.sheq(TypePair)
func SemanticReverseAbstractInterpreter_Sheq(p *rhino.JSType_TypePair) *rhino.JSType_TypePair {
	if p.TypeA == nil || p.TypeB == nil {
		return nil
	}
	return p.TypeA.GetTypesUnderShallowEquality(p.TypeB)
}

// SemanticReverseAbstractInterpreter_Shne is the merging function for strict
// non-equality between types.
//
// Java: SemanticReverseAbstractInterpreter.shne(TypePair)
func SemanticReverseAbstractInterpreter_Shne(p *rhino.JSType_TypePair) *rhino.JSType_TypePair {
	if p.TypeA == nil || p.TypeB == nil {
		return nil
	}
	return p.TypeA.GetTypesUnderShallowInequality(p.TypeB)
}

// SemanticReverseAbstractInterpreter_Ineq is the merging function for
// inequality comparisons between types.
//
// Java: SemanticReverseAbstractInterpreter.ineq(TypePair)
func SemanticReverseAbstractInterpreter_Ineq(p *rhino.JSType_TypePair) *rhino.JSType_TypePair {
	var typeA, typeB rhino.JSType
	if p.TypeA != nil {
		typeA = p.TypeA.RestrictByNotUndefined()
	}
	if p.TypeB != nil {
		typeB = p.TypeB.RestrictByNotUndefined()
	}
	return rhino.NewJSType_TypePair(typeA, typeB)
}

// NewSemanticReverseAbstractInterpreter creates a semantic reverse abstract
// interpreter.
//
// Java: SemanticReverseAbstractInterpreter(JSTypeRegistry)
func NewSemanticReverseAbstractInterpreter(typeRegistry *rhino.JSTypeRegistry) *SemanticReverseAbstractInterpreter {
	s := &SemanticReverseAbstractInterpreter{}
	s.InitChainableReverseAbstractInterpreterBase(s, typeRegistry)
	return s
}

// Java: SemanticReverseAbstractInterpreter.getPreciserScopeKnowingConditionOutcome(Node, FlowScope, Outcome)
func (s *SemanticReverseAbstractInterpreter) GetPreciserScopeKnowingConditionOutcome(
	condition *rhino.Node, blindScope FlowScope, outcome rhino.Outcome) FlowScope {
	// Check for the typeof operator.
	operatorToken := condition.GetToken()
	switch operatorToken {
	case rhino.Token_EQ, rhino.Token_NE, rhino.Token_SHEQ, rhino.Token_SHNE, rhino.Token_CASE:
		var left *rhino.Node
		var right *rhino.Node
		if operatorToken == rhino.Token_CASE {
			left = condition.GetParent().GetPrevious() // the switch condition
			jutil.CheckStateMsg1(condition.GetParent().IsSwitchBody(), "%s", condition.GetParent())
			right = condition.GetFirstChild()
		} else {
			left = condition.GetFirstChild()
			right = condition.GetLastChild()
		}

		var typeOfNode *rhino.Node
		var stringNode *rhino.Node
		if left.IsTypeOf() && right.IsStringLit() {
			typeOfNode = left
			stringNode = right
		} else if right.IsTypeOf() && left.IsStringLit() {
			typeOfNode = right
			stringNode = left
		}
		if typeOfNode != nil && stringNode != nil {
			operandNode := typeOfNode.GetFirstChild()
			operandType := s.GetTypeIfRefinable(operandNode, blindScope)
			if operandType != nil {
				resultEqualsValue := operatorToken == rhino.Token_EQ ||
					operatorToken == rhino.Token_SHEQ ||
					operatorToken == rhino.Token_CASE
				if !outcome.IsTruthy() {
					resultEqualsValue = !resultEqualsValue
				}
				return s.CaseTypeOf(
					operandNode, operandType, stringNode.GetString(), resultEqualsValue, blindScope)
			}
		}
	default:
	}
	switch operatorToken {
	case rhino.Token_AND:
		if outcome.IsTruthy() {
			return s.CaseAndOrNotShortCircuiting(
				condition.GetFirstChild(), condition.GetLastChild(), blindScope, rhino.Outcome_TRUE)
		} else {
			return s.CaseAndOrMaybeShortCircuiting(
				condition.GetFirstChild(), condition.GetLastChild(), blindScope, rhino.Outcome_TRUE)
		}
	case rhino.Token_OR:
		if !outcome.IsTruthy() {
			return s.CaseAndOrNotShortCircuiting(
				condition.GetFirstChild(), condition.GetLastChild(), blindScope, rhino.Outcome_FALSE)
		} else {
			return s.CaseAndOrMaybeShortCircuiting(
				condition.GetFirstChild(), condition.GetLastChild(), blindScope, rhino.Outcome_FALSE)
		}
	case rhino.Token_EQ:
		if outcome.IsTruthy() {
			return s.caseEqualityMode(condition.GetFirstChild(), condition.GetLastChild(), blindScope, sraiModeEq)
		} else {
			return s.caseEqualityMode(condition.GetFirstChild(), condition.GetLastChild(), blindScope, sraiModeNe)
		}
	case rhino.Token_NE:
		if outcome.IsTruthy() {
			return s.caseEqualityMode(condition.GetFirstChild(), condition.GetLastChild(), blindScope, sraiModeNe)
		} else {
			return s.caseEqualityMode(condition.GetFirstChild(), condition.GetLastChild(), blindScope, sraiModeEq)
		}
	case rhino.Token_SHEQ:
		if outcome.IsTruthy() {
			return s.caseEqualityMode(condition.GetFirstChild(), condition.GetLastChild(), blindScope, sraiModeSheq)
		} else {
			return s.caseEqualityMode(condition.GetFirstChild(), condition.GetLastChild(), blindScope, sraiModeShne)
		}
	case rhino.Token_SHNE:
		if outcome.IsTruthy() {
			return s.caseEqualityMode(condition.GetFirstChild(), condition.GetLastChild(), blindScope, sraiModeShne)
		} else {
			return s.caseEqualityMode(condition.GetFirstChild(), condition.GetLastChild(), blindScope, sraiModeSheq)
		}
	case rhino.Token_NAME, rhino.Token_GETPROP:
		return s.CaseNameOrGetProp(condition, blindScope, outcome)
	case rhino.Token_ASSIGN:
		return s.FirstPreciserScopeKnowingConditionOutcome(
			condition.GetFirstChild(),
			s.FirstPreciserScopeKnowingConditionOutcome(
				condition.GetSecondChild(), blindScope, outcome),
			outcome)
	case rhino.Token_NOT:
		return s.FirstPreciserScopeKnowingConditionOutcome(
			condition.GetFirstChild(), blindScope, outcome.Not())
	case rhino.Token_LE, rhino.Token_LT, rhino.Token_GE, rhino.Token_GT:
		if outcome.IsTruthy() {
			return s.caseEqualityMode(condition.GetFirstChild(), condition.GetLastChild(), blindScope, sraiModeIneq)
		}
	case rhino.Token_INSTANCEOF:
		return s.CaseInstanceOf(
			condition.GetFirstChild(), condition.GetLastChild(), blindScope, outcome)
	case rhino.Token_IN:
		if outcome.IsTruthy() && condition.GetFirstChild().IsStringLit() {
			return s.CaseIn(
				condition.GetLastChild(), condition.GetFirstChild().GetString(), blindScope)
		}
	case rhino.Token_CASE:
		left := condition.GetParent().GetPrevious() // the switch condition
		right := condition.GetFirstChild()
		if outcome.IsTruthy() {
			return s.caseEqualityMode(left, right, blindScope, sraiModeSheq)
		} else {
			return s.caseEqualityMode(left, right, blindScope, sraiModeShne)
		}
	case rhino.Token_CALL:
		left := condition.GetFirstChild()
		leftName := left.GetQualifiedName()
		if leftName != nil && *leftName == "Array.isArray" && left.GetNext() != nil {
			return s.CaseIsArray(left.GetNext(), blindScope, outcome)
		}
	default:
	}

	return s.NextPreciserScopeKnowingConditionOutcome(condition, blindScope, outcome)
}

// Java: SemanticReverseAbstractInterpreter.caseIsArray(Node, FlowScope, Outcome)
func (s *SemanticReverseAbstractInterpreter) CaseIsArray(value *rhino.Node, blindScope FlowScope, outcome rhino.Outcome) FlowScope {
	type_ := s.GetTypeIfRefinable(value, blindScope)
	if type_ != nil {
		var visitor rhino.Visitor
		if outcome.IsTruthy() {
			visitor = s.restrictToArrayVisitor
		} else {
			visitor = s.restrictToNotArrayVisitor
		}
		return s.MaybeRestrictName(blindScope, value, type_, rhino.JSType_VisitForJSType(type_, visitor))
	}
	return blindScope
}

type sraiEqualityMode uint8

const (
	sraiModeEq sraiEqualityMode = iota
	sraiModeNe
	sraiModeSheq
	sraiModeShne
	sraiModeIneq
)

func (s *SemanticReverseAbstractInterpreter) caseEqualityMode(
	left *rhino.Node, right *rhino.Node, blindScope FlowScope, mode sraiEqualityMode) FlowScope {
	leftType := s.GetTypeIfRefinable(left, blindScope)
	leftIsRefineable := leftType != nil
	if !leftIsRefineable {
		leftType = left.GetJSType()
	}

	rightType := s.GetTypeIfRefinable(right, blindScope)
	rightIsRefineable := rightType != nil
	if !rightIsRefineable {
		rightType = right.GetJSType()
	}

	var mergedA, mergedB rhino.JSType
	var ok bool
	switch mode {
	case sraiModeEq:
		if leftType != nil && rightType != nil {
			mergedA, mergedB = leftType.GetTypesUnderEqualityPair(rightType)
			ok = true
		}
	case sraiModeNe:
		if leftType != nil && rightType != nil {
			mergedA, mergedB = leftType.GetTypesUnderInequalityPair(rightType)
			ok = true
		}
	case sraiModeSheq:
		if leftType != nil && rightType != nil {
			mergedA, mergedB = leftType.GetTypesUnderShallowEqualityPair(rightType)
			ok = true
		}
	case sraiModeShne:
		if leftType != nil && rightType != nil {
			mergedA, mergedB = leftType.GetTypesUnderShallowInequalityPair(rightType)
			ok = true
		}
	case sraiModeIneq:
		if leftType != nil {
			mergedA = leftType.RestrictByNotUndefined()
		}
		if rightType != nil {
			mergedB = rightType.RestrictByNotUndefined()
		}
		ok = true
	}

	if ok {
		var restrictedLeftType, restrictedRightType rhino.JSType
		if leftIsRefineable {
			restrictedLeftType = mergedA
		}
		if rightIsRefineable {
			restrictedRightType = mergedB
		}
		return s.MaybeRestrictTwoNames(
			blindScope,
			left,
			leftType,
			restrictedLeftType,
			right,
			rightType,
			restrictedRightType)
	}
	return blindScope
}

// Java: SemanticReverseAbstractInterpreter.caseEquality(Node, FlowScope, Function<TypePair, TypePair>)
func (s *SemanticReverseAbstractInterpreter) CaseEquality(
	condition *rhino.Node, blindScope FlowScope, merging func(*rhino.JSType_TypePair) *rhino.JSType_TypePair) FlowScope {
	return s.CaseEqualityOfLeftAndRight(condition.GetFirstChild(), condition.GetLastChild(), blindScope, merging)
}

// Java: SemanticReverseAbstractInterpreter.caseEquality(Node, Node, FlowScope, Function<TypePair, TypePair>)
func (s *SemanticReverseAbstractInterpreter) CaseEqualityOfLeftAndRight(
	left *rhino.Node, right *rhino.Node, blindScope FlowScope,
	merging func(*rhino.JSType_TypePair) *rhino.JSType_TypePair) FlowScope {
	// left type
	leftType := s.GetTypeIfRefinable(left, blindScope)
	var leftIsRefineable bool
	if leftType != nil {
		leftIsRefineable = true
	} else {
		leftIsRefineable = false
		leftType = left.GetJSType()
	}

	// right type
	rightType := s.GetTypeIfRefinable(right, blindScope)
	var rightIsRefineable bool
	if rightType != nil {
		rightIsRefineable = true
	} else {
		rightIsRefineable = false
		rightType = right.GetJSType()
	}

	// merged types
	merged := merging(rhino.NewJSType_TypePair(leftType, rightType))

	// creating new scope
	if merged != nil {
		var restrictedLeftType, restrictedRightType rhino.JSType
		if leftIsRefineable {
			restrictedLeftType = merged.TypeA
		}
		if rightIsRefineable {
			restrictedRightType = merged.TypeB
		}
		return s.MaybeRestrictTwoNames(
			blindScope,
			left,
			leftType,
			restrictedLeftType,
			right,
			rightType,
			restrictedRightType)
	}
	return blindScope
}

// Java: SemanticReverseAbstractInterpreter.caseAndOrNotShortCircuiting(Node, Node, FlowScope, Outcome)
func (s *SemanticReverseAbstractInterpreter) CaseAndOrNotShortCircuiting(
	left *rhino.Node, right *rhino.Node, blindScope FlowScope, outcome rhino.Outcome) FlowScope {
	// left type
	leftType := s.GetTypeIfRefinable(left, blindScope)
	var leftIsRefineable bool
	if leftType != nil {
		leftIsRefineable = true
	} else {
		leftIsRefineable = false
		leftType = left.GetJSType()
		blindScope = s.FirstPreciserScopeKnowingConditionOutcome(left, blindScope, outcome)
	}

	// restricting left type
	var restrictedLeftType rhino.JSType
	if leftType != nil {
		restrictedLeftType = leftType.GetRestrictedTypeGivenOutcome(outcome)
	}
	if restrictedLeftType == nil {
		return s.FirstPreciserScopeKnowingConditionOutcome(right, blindScope, outcome)
	}
	var leftRestriction rhino.JSType
	if leftIsRefineable {
		leftRestriction = restrictedLeftType
	}
	blindScope = s.MaybeRestrictName(blindScope, left, leftType, leftRestriction)

	// right type
	rightType := s.GetTypeIfRefinable(right, blindScope)
	var rightIsRefineable bool
	if rightType != nil {
		rightIsRefineable = true
	} else {
		rightIsRefineable = false
		rightType = right.GetJSType()
		blindScope = s.FirstPreciserScopeKnowingConditionOutcome(right, blindScope, outcome)
	}

	if outcome.IsTruthy() {
		var restrictedRightType rhino.JSType
		if rightType != nil {
			restrictedRightType = rightType.GetRestrictedTypeGivenOutcome(outcome)
		}
		// creating new scope
		var rightRestriction rhino.JSType
		if rightIsRefineable {
			rightRestriction = restrictedRightType
		}
		return s.MaybeRestrictName(blindScope, right, rightType, rightRestriction)
	}
	return blindScope
}

// Java: SemanticReverseAbstractInterpreter.caseAndOrMaybeShortCircuiting(Node, Node, FlowScope, Outcome)
func (s *SemanticReverseAbstractInterpreter) CaseAndOrMaybeShortCircuiting(
	left *rhino.Node, right *rhino.Node, blindScope FlowScope, outcome rhino.Outcome) FlowScope {
	// Perform two separate refinements, one for if short-circuiting occurred, and one for if it did
	// not.  Because it's not clear whether short-circuiting occurred, we actually have to ignore
	// both separate result flow scopes individually, but if they both refined the same slot, we
	// can join the two refinements.  TODO(sdh): look into simplifying this.  If joining were
	// more efficient, we should just be able to join the scopes unconditionally?
	mark := len(s.trackingUsed)
	rootTracking := s.acquireTrackingScope(blindScope, nil)
	tracker := rootTracking.tracker
	blindScope = rootTracking
	leftScope := s.FirstPreciserScopeKnowingConditionOutcome(left, blindScope, outcome.Not())
	var leftVar rhino.StaticTypedSlot
	if tracker.count == 1 {
		leftVar, _ = leftScope.GetSlot(tracker.first).(rhino.StaticTypedSlot)
	}
	if leftVar == nil {
		// If we did create a more precise scope, blindScope has a child and
		// it is frozen. We can't just throw it away to return it. So we
		// must create a child instead.
		res := SemanticReverseAbstractInterpreter_Unwrap(blindScope)
		s.releaseTrackingScopes(mark)
		return res
	}
	*tracker = sraiRefinementTracker{}
	// Note: re-wrap the scope, in case it was unwrapped by a nested call to this method.
	var rightScope FlowScope = s.acquireTrackingScope(
		s.FirstPreciserScopeKnowingConditionOutcome(left, blindScope, outcome), tracker)
	rightScope = s.FirstPreciserScopeKnowingConditionOutcome(right, rightScope, outcome.Not())
	var rightVar rhino.StaticTypedSlot
	if tracker.count == 1 {
		rightVar, _ = rightScope.GetSlot(tracker.first).(rhino.StaticTypedSlot)
	}
	if rightVar == nil || leftVar.GetName() != rightVar.GetName() {
		res := SemanticReverseAbstractInterpreter_Unwrap(blindScope)
		s.releaseTrackingScopes(mark)
		return res
	}
	type_ := leftVar.GetType().GetLeastSupertype(rightVar.GetType())
	unwrapped := SemanticReverseAbstractInterpreter_Unwrap(blindScope)
	s.releaseTrackingScopes(mark)
	return unwrapped.InferSlotType(leftVar.GetName(), type_)
}

// MaybeRestrictName: if the restrictedType differs from the originalType,
// then we should branch the current flow scope and create a new flow scope
// with the name declared with the new type.
//
// We try not to create spurious child flow scopes as this makes type
// inference slower.
//
// We also do not want spurious slots around in type inference, because we
// use these as a signal for "checked unknown" types. A "checked unknown"
// type is a symbol that the programmer has already checked and verified that
// it's defined, even if we don't know what it is.
//
// It is OK to pass non-name nodes into this method, as long as you pass in
// null for a restricted type.
//
// Java: SemanticReverseAbstractInterpreter.maybeRestrictName(FlowScope, Node, JSType, JSType)
func (s *SemanticReverseAbstractInterpreter) MaybeRestrictName(
	blindScope FlowScope, node *rhino.Node, originalType rhino.JSType, restrictedType rhino.JSType) FlowScope {
	if restrictedType != nil && !base.JSCompObjects_Identical(restrictedType, originalType) {
		return s.DeclareNameInScope(blindScope, node, restrictedType)
	}
	return blindScope
}

// MaybeRestrictTwoNames: see MaybeRestrictName.
//
// Java: SemanticReverseAbstractInterpreter.maybeRestrictTwoNames(FlowScope, Node, JSType, JSType, Node, JSType, JSType)
func (s *SemanticReverseAbstractInterpreter) MaybeRestrictTwoNames(
	blindScope FlowScope,
	left *rhino.Node,
	originalLeftType rhino.JSType,
	restrictedLeftType rhino.JSType,
	right *rhino.Node,
	originalRightType rhino.JSType,
	restrictedRightType rhino.JSType) FlowScope {
	shouldRefineLeft := restrictedLeftType != nil && !base.JSCompObjects_Identical(restrictedLeftType, originalLeftType)
	shouldRefineRight := restrictedRightType != nil && !base.JSCompObjects_Identical(restrictedRightType, originalRightType)
	if shouldRefineLeft || shouldRefineRight {
		informed := blindScope
		if shouldRefineLeft {
			informed = s.DeclareNameInScope(informed, left, restrictedLeftType)
		}
		if shouldRefineRight {
			informed = s.DeclareNameInScope(informed, right, restrictedRightType)
		}
		return informed
	}
	return blindScope
}

// Java: SemanticReverseAbstractInterpreter.caseNameOrGetProp(Node, FlowScope, Outcome)
func (s *SemanticReverseAbstractInterpreter) CaseNameOrGetProp(name *rhino.Node, blindScope FlowScope, outcome rhino.Outcome) FlowScope {
	type_ := s.GetTypeIfRefinable(name, blindScope)
	if type_ != nil {
		restrictedType := type_.GetRestrictedTypeGivenOutcome(outcome)
		return s.MaybeRestrictName(blindScope, name, type_, restrictedType)
	}
	return blindScope
}

// Java: SemanticReverseAbstractInterpreter.caseTypeOf(Node, JSType, String, boolean, FlowScope)
func (s *SemanticReverseAbstractInterpreter) CaseTypeOf(
	node *rhino.Node, type_ rhino.JSType, value string, resultEqualsValue bool, blindScope FlowScope) FlowScope {
	return s.MaybeRestrictName(
		blindScope, node, type_, s.GetRestrictedByTypeOfResult(type_, value, resultEqualsValue))
}

// Java: SemanticReverseAbstractInterpreter.caseInstanceOf(Node, Node, FlowScope, Outcome)
func (s *SemanticReverseAbstractInterpreter) CaseInstanceOf(
	left *rhino.Node, right *rhino.Node, blindScope FlowScope, outcome rhino.Outcome) FlowScope {
	leftType := s.GetTypeIfRefinable(left, blindScope)
	if leftType == nil {
		return blindScope
	}
	rightType := right.GetJSType()
	targetType := s.typeRegistry.GetNativeObjectType(rhino.JSTypeNative_UNKNOWN_TYPE)
	if rightType != nil && rightType.IsFunctionType() {
		targetType = rightType.ToMaybeFunctionType()
	}
	var visitor rhino.Visitor
	if outcome.IsTruthy() {
		visitor = newSemanticReverseAbstractInterpreter_RestrictByTrueInstanceOfResultVisitor(s, targetType)
	} else {
		visitor = newSemanticReverseAbstractInterpreter_RestrictByFalseInstanceOfResultVisitor(s, targetType)
	}
	return s.MaybeRestrictName(blindScope, left, leftType, rhino.JSType_VisitForJSType(leftType, visitor))
}

// CaseIn: given 'property in object', ensures that the object has the
// property in the informed scope by defining it as a qualified name if the
// object type lacks the property and it's not in the blind scope.
//
// object is the node of the right-side of the in; propertyName is the
// string of the left-side of the in.
//
// Java: SemanticReverseAbstractInterpreter.caseIn(Node, String, FlowScope)
func (s *SemanticReverseAbstractInterpreter) CaseIn(object *rhino.Node, propertyName string, blindScope FlowScope) FlowScope {
	jsType := object.GetJSType()
	if jsType != nil {
		jsType = jsType.RestrictByNotNullOrUndefined()
	}

	hasProperty := false
	objectType := rhino.ObjectType_Cast(jsType)
	if objectType != nil {
		hasProperty = objectType.HasProperty(propertyName)
	}
	if !hasProperty {
		qualifiedName := object.GetQualifiedName()
		if qualifiedName != nil {
			propertyQualifiedName := *qualifiedName + "." + propertyName
			if blindScope.GetSlot(propertyQualifiedName) == nil {
				unknownType := s.typeRegistry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
				return blindScope.InferQualifiedSlot(
					object, propertyQualifiedName, unknownType, unknownType, false)
			}
		}
	}
	return blindScope
}

// SemanticReverseAbstractInterpreter_RestrictByTrueInstanceOfResultVisitor:
// see SemanticReverseAbstractInterpreter.CaseInstanceOf.
//
// Java: private class SemanticReverseAbstractInterpreter.RestrictByTrueInstanceOfResultVisitor
// extends RestrictByTrueTypeOfResultVisitor
type SemanticReverseAbstractInterpreter_RestrictByTrueInstanceOfResultVisitor struct {
	ChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase

	target rhino.ObjectType
}

// Java: RestrictByTrueInstanceOfResultVisitor(ObjectType)
func newSemanticReverseAbstractInterpreter_RestrictByTrueInstanceOfResultVisitor(
	outer *SemanticReverseAbstractInterpreter, target rhino.ObjectType) *SemanticReverseAbstractInterpreter_RestrictByTrueInstanceOfResultVisitor {
	v := &SemanticReverseAbstractInterpreter_RestrictByTrueInstanceOfResultVisitor{}
	v.InitChainableReverseAbstractInterpreter_RestrictByTrueTypeOfResultVisitorBase(v, outer)
	v.target = target
	return v
}

// Java: RestrictByTrueInstanceOfResultVisitor.caseTopType(JSType)
func (v *SemanticReverseAbstractInterpreter_RestrictByTrueInstanceOfResultVisitor) CaseTopType(type_ rhino.JSType) rhino.JSType {
	return v.ApplyCommonRestriction(type_)
}

// Java: RestrictByTrueInstanceOfResultVisitor.caseUnknownType()
func (v *SemanticReverseAbstractInterpreter_RestrictByTrueInstanceOfResultVisitor) CaseUnknownType() any {
	funcTarget := rhino.JSType_ToMaybeFunctionType(v.target)
	if funcTarget != nil && funcTarget.HasInstanceType() {
		return funcTarget.GetInstanceType()
	}
	return v.outer.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
}

// Java: RestrictByTrueInstanceOfResultVisitor.caseObjectType(ObjectType)
func (v *SemanticReverseAbstractInterpreter_RestrictByTrueInstanceOfResultVisitor) CaseObjectType(type_ rhino.ObjectType) any {
	return v.ApplyCommonRestriction(type_)
}

// Java: RestrictByTrueInstanceOfResultVisitor.caseUnionType(UnionType)
func (v *SemanticReverseAbstractInterpreter_RestrictByTrueInstanceOfResultVisitor) CaseUnionType(type_ *rhino.UnionType) any {
	return v.ApplyCommonRestriction(type_)
}

// Java: RestrictByTrueInstanceOfResultVisitor.caseFunctionType(FunctionType)
func (v *SemanticReverseAbstractInterpreter_RestrictByTrueInstanceOfResultVisitor) CaseFunctionType(type_ rhino.FunctionType) any {
	return v.self.CaseObjectType(type_)
}

// Java: RestrictByTrueInstanceOfResultVisitor.applyCommonRestriction(JSType)
func (v *SemanticReverseAbstractInterpreter_RestrictByTrueInstanceOfResultVisitor) ApplyCommonRestriction(type_ rhino.JSType) rhino.JSType {
	if v.target.IsUnknownType() {
		return type_
	}

	funcTarget := v.target.ToMaybeFunctionType()
	if funcTarget.HasInstanceType() {
		return type_.GetGreatestSubtype(funcTarget.GetInstanceType())
	}

	return nil
}

// SemanticReverseAbstractInterpreter_RestrictByFalseInstanceOfResultVisitor:
// see SemanticReverseAbstractInterpreter.CaseInstanceOf.
//
// Java: private class SemanticReverseAbstractInterpreter.RestrictByFalseInstanceOfResultVisitor
// extends RestrictByFalseTypeOfResultVisitor
type SemanticReverseAbstractInterpreter_RestrictByFalseInstanceOfResultVisitor struct {
	ChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase

	target rhino.ObjectType
}

// Java: RestrictByFalseInstanceOfResultVisitor(ObjectType)
func newSemanticReverseAbstractInterpreter_RestrictByFalseInstanceOfResultVisitor(
	outer *SemanticReverseAbstractInterpreter, target rhino.ObjectType) *SemanticReverseAbstractInterpreter_RestrictByFalseInstanceOfResultVisitor {
	v := &SemanticReverseAbstractInterpreter_RestrictByFalseInstanceOfResultVisitor{}
	v.InitChainableReverseAbstractInterpreter_RestrictByFalseTypeOfResultVisitorBase(v, outer)
	v.target = target
	return v
}

// Java: RestrictByFalseInstanceOfResultVisitor.caseObjectType(ObjectType)
func (v *SemanticReverseAbstractInterpreter_RestrictByFalseInstanceOfResultVisitor) CaseObjectType(type_ rhino.ObjectType) any {
	if v.target.IsUnknownType() {
		return type_
	}

	funcTarget := v.target.ToMaybeFunctionType()
	if funcTarget.HasInstanceType() {
		if type_.IsSubtypeOf(funcTarget.GetInstanceType()) {
			return nil
		}

		return type_
	}

	return nil
}

// Java: RestrictByFalseInstanceOfResultVisitor.caseUnionType(UnionType)
func (v *SemanticReverseAbstractInterpreter_RestrictByFalseInstanceOfResultVisitor) CaseUnionType(type_ *rhino.UnionType) any {
	if v.target.IsUnknownType() {
		return type_
	}

	funcTarget := v.target.ToMaybeFunctionType()
	if funcTarget.HasInstanceType() {
		return type_.GetRestrictedUnion(funcTarget.GetInstanceType())
	}

	return nil
}

// Java: RestrictByFalseInstanceOfResultVisitor.caseFunctionType(FunctionType)
func (v *SemanticReverseAbstractInterpreter_RestrictByFalseInstanceOfResultVisitor) CaseFunctionType(type_ rhino.FunctionType) any {
	return v.self.CaseObjectType(type_)
}

// SemanticReverseAbstractInterpreter_Unwrap unwraps any
// RefinementTrackingFlowScopes.
//
// Java: SemanticReverseAbstractInterpreter.unwrap(FlowScope)
func SemanticReverseAbstractInterpreter_Unwrap(scope FlowScope) FlowScope {
	for {
		r, ok := scope.(*SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope)
		if !ok {
			break
		}
		scope = r.delegate
	}
	return scope
}

type sraiRefinementTracker struct {
	first string
	count int
}

func (t *sraiRefinementTracker) add(symbol string) {
	if t.count == 0 {
		t.first = symbol
		t.count = 1
	} else if t.count == 1 && t.first != symbol {
		t.count = 2
	}
}

// SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope is a
// wrapper around FlowScope that keeps track of which vars were refined.
type SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope struct {
	owner      *SemanticReverseAbstractInterpreter
	delegate   FlowScope
	tracker    *sraiRefinementTracker
	ownTracker sraiRefinementTracker
}

func (s *SemanticReverseAbstractInterpreter) acquireTrackingScope(
	delegate FlowScope, tracker *sraiRefinementTracker,
) *SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope {
	var r *SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope
	if n := len(s.trackingFree); n > 0 {
		r = s.trackingFree[n-1]
		s.trackingFree[n-1] = nil
		s.trackingFree = s.trackingFree[:n-1]
	} else {
		r = s.trackingSlab.New()
	}
	r.owner = s
	r.delegate = delegate
	if tracker == nil {
		r.ownTracker = sraiRefinementTracker{}
		r.tracker = &r.ownTracker
	} else {
		r.tracker = tracker
	}
	s.trackingUsed = append(s.trackingUsed, r)
	return r
}

func (s *SemanticReverseAbstractInterpreter) releaseTrackingScopes(mark int) {
	for i := len(s.trackingUsed) - 1; i >= mark; i-- {
		r := s.trackingUsed[i]
		s.trackingUsed[i] = nil
		r.owner = nil
		r.delegate = nil
		r.tracker = nil
		r.ownTracker = sraiRefinementTracker{}
		s.trackingFree = append(s.trackingFree, r)
	}
	s.trackingUsed = s.trackingUsed[:mark]
}

// Java: RefinementTrackingFlowScope.withSyntacticScope(StaticTypedScope)
func (r *SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope) WithSyntacticScope(scope rhino.StaticTypedScope) FlowScope {
	panic(jutil.NewUnsupportedOperationException())
}

// Java: RefinementTrackingFlowScope.inferSlotType(String, JSType)
func (r *SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope) InferSlotType(symbol string, type_ rhino.JSType) FlowScope {
	r.tracker.add(symbol)
	return r.Wrap(r.delegate.InferSlotType(symbol, type_))
}

// Java: RefinementTrackingFlowScope.inferQualifiedSlot(Node, String, JSType, JSType, boolean)
func (r *SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope) InferQualifiedSlot(
	node *rhino.Node, symbol string, bottomType rhino.JSType, inferredType rhino.JSType, declare bool) FlowScope {
	r.tracker.add(symbol)
	return r.Wrap(r.delegate.InferQualifiedSlot(node, symbol, bottomType, inferredType, declare))
}

// Java: RefinementTrackingFlowScope.wrap(FlowScope)
func (r *SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope) Wrap(scope FlowScope) FlowScope {
	if scope != r.delegate {
		return r.owner.acquireTrackingScope(scope, r.tracker)
	}
	return r
}

// Java: RefinementTrackingFlowScope.getDeclarationScope()
func (r *SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope) GetDeclarationScope() rhino.StaticTypedScope {
	return r.delegate.GetDeclarationScope()
}

// Java: RefinementTrackingFlowScope.getRootNode()
func (r *SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope) GetRootNode() *rhino.Node {
	return r.delegate.GetRootNode()
}

// Java: RefinementTrackingFlowScope.getParentScope()
func (r *SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope) GetParentScope() rhino.StaticScope {
	panic(jutil.NewUnsupportedOperationException())
}

// Java: RefinementTrackingFlowScope.getSlot(String)
func (r *SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope) GetSlot(name string) rhino.StaticSlot {
	return r.delegate.GetSlot(name)
}

// Java: RefinementTrackingFlowScope.getOwnSlot(String)
func (r *SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope) GetOwnSlot(name string) rhino.StaticSlot {
	return r.delegate.GetOwnSlot(name)
}

// Java: RefinementTrackingFlowScope.getTypeOfThis()
func (r *SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope) GetTypeOfThis() rhino.JSType {
	return r.delegate.GetTypeOfThis()
}

// GetTopmostScopeOfEventualDeclaration is the inherited default method.
//
// Java: StaticScope.getTopmostScopeOfEventualDeclaration(String) (default method)
func (r *SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope) GetTopmostScopeOfEventualDeclaration(name string) rhino.StaticScope {
	return rhino.StaticScope_GetTopmostScopeOfEventualDeclaration(r, name)
}

// LookupQualifiedName is the inherited default method.
//
// Java: StaticTypedScope.lookupQualifiedName(QualifiedName) (default method)
func (r *SemanticReverseAbstractInterpreter_RefinementTrackingFlowScope) LookupQualifiedName(qname rhino.QualifiedName) rhino.JSType {
	return rhino.StaticTypedScope_LookupQualifiedName(r, qname)
}
