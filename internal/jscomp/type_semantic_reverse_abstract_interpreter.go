// Port of com/google/javascript/jscomp/type/SemanticReverseAbstractInterpreter.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// SemanticReverseAbstractInterpreter is a reverse abstract interpreter using
// the semantics of the JavaScript language as a means to reverse interpret
// computations. This interpreter expects the parse tree inputs to be typed.
type SemanticReverseAbstractInterpreter struct {
	ChainableReverseAbstractInterpreterBase
	trackingUsed []*sraiTrackingFlowScope
	trackingFree []*sraiTrackingFlowScope
	trackingSlab jutil.Slab[sraiTrackingFlowScope]
}

// NewSemanticReverseAbstractInterpreter creates a semantic reverse abstract
// interpreter.
//
// Java: SemanticReverseAbstractInterpreter(JSTypeRegistry)
func NewSemanticReverseAbstractInterpreter(typeRegistry *rhino.JSTypeRegistry) *SemanticReverseAbstractInterpreter {
	s := &SemanticReverseAbstractInterpreter{}
	s.initChainableBase(s, typeRegistry)
	return s
}

// Java: SemanticReverseAbstractInterpreter.getPreciserScopeKnowingConditionOutcome(Node, FlowScope, Outcome)
func (s *SemanticReverseAbstractInterpreter) GetPreciserScopeKnowingConditionOutcome(
	condition rhino.Node, blindScope FlowScope, outcome rhino.Outcome) FlowScope {
	// Check for the typeof operator.
	operatorToken := condition.GetToken()
	switch operatorToken {
	case rhino.Token_EQ, rhino.Token_NE, rhino.Token_SHEQ, rhino.Token_SHNE, rhino.Token_CASE:
		var left rhino.Node
		var right rhino.Node
		if operatorToken == rhino.Token_CASE {
			jutil.CheckStateMsg1(condition.GetParent().IsSwitchBody(), "%s", condition.GetParent())
			left = ast.MustSwitch(condition.GetGrandparent()).Discriminant().Node() // the switch condition
			right = ast.MustCase(condition).Value().Node()
		} else {
			bin := ast.MustBinary(condition)
			left = bin.Left().Node()
			right = bin.Right().Node()
		}

		var typeOfNode rhino.Node
		var stringNode rhino.Node
		if left.IsTypeOf() && right.IsStringLit() {
			typeOfNode = left
			stringNode = right
		} else if right.IsTypeOf() && left.IsStringLit() {
			typeOfNode = right
			stringNode = left
		}
		if !typeOfNode.IsNil() && !stringNode.IsNil() {
			operandNode := ast.MustUnary(typeOfNode).Operand().Node()
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
		bin := ast.MustBinary(condition)
		if outcome.IsTruthy() {
			return s.CaseAndOrNotShortCircuiting(
				bin.Left().Node(), bin.Right().Node(), blindScope, rhino.Outcome_TRUE)
		} else {
			return s.CaseAndOrMaybeShortCircuiting(
				bin.Left().Node(), bin.Right().Node(), blindScope, rhino.Outcome_TRUE)
		}
	case rhino.Token_OR:
		bin := ast.MustBinary(condition)
		if !outcome.IsTruthy() {
			return s.CaseAndOrNotShortCircuiting(
				bin.Left().Node(), bin.Right().Node(), blindScope, rhino.Outcome_FALSE)
		} else {
			return s.CaseAndOrMaybeShortCircuiting(
				bin.Left().Node(), bin.Right().Node(), blindScope, rhino.Outcome_FALSE)
		}
	case rhino.Token_EQ:
		bin := ast.MustBinary(condition)
		if outcome.IsTruthy() {
			return s.caseEqualityMode(bin.Left().Node(), bin.Right().Node(), blindScope, sraiModeEq)
		} else {
			return s.caseEqualityMode(bin.Left().Node(), bin.Right().Node(), blindScope, sraiModeNe)
		}
	case rhino.Token_NE:
		bin := ast.MustBinary(condition)
		if outcome.IsTruthy() {
			return s.caseEqualityMode(bin.Left().Node(), bin.Right().Node(), blindScope, sraiModeNe)
		} else {
			return s.caseEqualityMode(bin.Left().Node(), bin.Right().Node(), blindScope, sraiModeEq)
		}
	case rhino.Token_SHEQ:
		bin := ast.MustBinary(condition)
		if outcome.IsTruthy() {
			return s.caseEqualityMode(bin.Left().Node(), bin.Right().Node(), blindScope, sraiModeSheq)
		} else {
			return s.caseEqualityMode(bin.Left().Node(), bin.Right().Node(), blindScope, sraiModeShne)
		}
	case rhino.Token_SHNE:
		bin := ast.MustBinary(condition)
		if outcome.IsTruthy() {
			return s.caseEqualityMode(bin.Left().Node(), bin.Right().Node(), blindScope, sraiModeShne)
		} else {
			return s.caseEqualityMode(bin.Left().Node(), bin.Right().Node(), blindScope, sraiModeSheq)
		}
	case rhino.Token_NAME, rhino.Token_GETPROP:
		return s.CaseNameOrGetProp(condition, blindScope, outcome)
	case rhino.Token_ASSIGN:
		assign := ast.MustAssign(condition)
		return s.FirstPreciserScopeKnowingConditionOutcome(
			assign.Target().Node(),
			s.FirstPreciserScopeKnowingConditionOutcome(
				assign.Value().Node(), blindScope, outcome),
			outcome)
	case rhino.Token_NOT:
		return s.FirstPreciserScopeKnowingConditionOutcome(
			ast.MustUnary(condition).Operand().Node(), blindScope, outcome.Not())
	case rhino.Token_LE, rhino.Token_LT, rhino.Token_GE, rhino.Token_GT:
		if outcome.IsTruthy() {
			bin := ast.MustBinary(condition)
			return s.caseEqualityMode(bin.Left().Node(), bin.Right().Node(), blindScope, sraiModeIneq)
		}
	case rhino.Token_INSTANCEOF:
		bin := ast.MustBinary(condition)
		return s.CaseInstanceOf(
			bin.Left().Node(), bin.Right().Node(), blindScope, outcome)
	case rhino.Token_IN:
		bin := ast.MustBinary(condition)
		if outcome.IsTruthy() && bin.Left().Node().IsStringLit() {
			return s.CaseIn(
				bin.Right().Node(), bin.Left().Node().GetString(), blindScope)
		}
	case rhino.Token_CASE:
		left := ast.MustSwitch(condition.GetGrandparent()).Discriminant().Node() // the switch condition
		right := ast.MustCase(condition).Value().Node()
		if outcome.IsTruthy() {
			return s.caseEqualityMode(left, right, blindScope, sraiModeSheq)
		} else {
			return s.caseEqualityMode(left, right, blindScope, sraiModeShne)
		}
	case rhino.Token_CALL:
		call := ast.MustCall(condition)
		if firstArg, ok := call.FirstArg(); ok && call.Callee().Node().MatchesQualifiedName("Array.isArray") {
			return s.CaseIsArray(firstArg.Node(), blindScope, outcome)
		}
	default:
	}

	return s.NextPreciserScopeKnowingConditionOutcome(condition, blindScope, outcome)
}

// Java: SemanticReverseAbstractInterpreter.caseIsArray(Node, FlowScope, Outcome)
func (s *SemanticReverseAbstractInterpreter) CaseIsArray(value rhino.Node, blindScope FlowScope, outcome rhino.Outcome) FlowScope {
	type_ := s.GetTypeIfRefinable(value, blindScope)
	if type_ != nil {
		var restricted rhino.JSType
		if outcome.IsTruthy() {
			restricted = s.restrictToArray(type_)
		} else {
			restricted = s.restrictToNotArray(type_)
		}
		return s.MaybeRestrictName(blindScope, value, type_, restricted)
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
	left rhino.Node, right rhino.Node, blindScope FlowScope, mode sraiEqualityMode) FlowScope {
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

// Java: SemanticReverseAbstractInterpreter.caseAndOrNotShortCircuiting(Node, Node, FlowScope, Outcome)
func (s *SemanticReverseAbstractInterpreter) CaseAndOrNotShortCircuiting(
	left rhino.Node, right rhino.Node, blindScope FlowScope, outcome rhino.Outcome) FlowScope {
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
	left rhino.Node, right rhino.Node, blindScope FlowScope, outcome rhino.Outcome) FlowScope {
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
	blindScope FlowScope, node rhino.Node, originalType rhino.JSType, restrictedType rhino.JSType) FlowScope {
	if restrictedType != nil && restrictedType != originalType {
		return s.DeclareNameInScope(blindScope, node, restrictedType)
	}
	return blindScope
}

// MaybeRestrictTwoNames: see MaybeRestrictName.
//
// Java: SemanticReverseAbstractInterpreter.maybeRestrictTwoNames(FlowScope, Node, JSType, JSType, Node, JSType, JSType)
func (s *SemanticReverseAbstractInterpreter) MaybeRestrictTwoNames(
	blindScope FlowScope,
	left rhino.Node,
	originalLeftType rhino.JSType,
	restrictedLeftType rhino.JSType,
	right rhino.Node,
	originalRightType rhino.JSType,
	restrictedRightType rhino.JSType) FlowScope {
	shouldRefineLeft := restrictedLeftType != nil && restrictedLeftType != originalLeftType
	shouldRefineRight := restrictedRightType != nil && restrictedRightType != originalRightType
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
func (s *SemanticReverseAbstractInterpreter) CaseNameOrGetProp(name rhino.Node, blindScope FlowScope, outcome rhino.Outcome) FlowScope {
	type_ := s.GetTypeIfRefinable(name, blindScope)
	if type_ != nil {
		restrictedType := type_.GetRestrictedTypeGivenOutcome(outcome)
		return s.MaybeRestrictName(blindScope, name, type_, restrictedType)
	}
	return blindScope
}

// Java: SemanticReverseAbstractInterpreter.caseTypeOf(Node, JSType, String, boolean, FlowScope)
func (s *SemanticReverseAbstractInterpreter) CaseTypeOf(
	node rhino.Node, type_ rhino.JSType, value string, resultEqualsValue bool, blindScope FlowScope) FlowScope {
	return s.MaybeRestrictName(
		blindScope, node, type_, s.GetRestrictedByTypeOfResult(type_, value, resultEqualsValue))
}

// Java: SemanticReverseAbstractInterpreter.caseInstanceOf(Node, Node, FlowScope, Outcome)
func (s *SemanticReverseAbstractInterpreter) CaseInstanceOf(
	left rhino.Node, right rhino.Node, blindScope FlowScope, outcome rhino.Outcome) FlowScope {
	leftType := s.GetTypeIfRefinable(left, blindScope)
	if leftType == nil {
		return blindScope
	}
	rightType := right.GetJSType()
	targetType := s.typeRegistry.GetNativeObjectType(rhino.JSTypeNative_UNKNOWN_TYPE)
	if rightType != nil && rightType.IsFunctionType() {
		targetType = rightType.ToMaybeFunctionType()
	}
	var restricted rhino.JSType
	if outcome.IsTruthy() {
		restricted = s.restrictByTrueInstanceOf(leftType, targetType)
	} else {
		restricted = s.restrictByFalseInstanceOf(leftType, targetType)
	}
	return s.MaybeRestrictName(blindScope, left, leftType, restricted)
}

// CaseIn: given 'property in object', ensures that the object has the
// property in the informed scope by defining it as a qualified name if the
// object type lacks the property and it's not in the blind scope.
//
// object is the node of the right-side of the in; propertyName is the
// string of the left-side of the in.
//
// Java: SemanticReverseAbstractInterpreter.caseIn(Node, String, FlowScope)
func (s *SemanticReverseAbstractInterpreter) CaseIn(object rhino.Node, propertyName string, blindScope FlowScope) FlowScope {
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
		if qualifiedName, ok := object.GetQualifiedNameString(); ok {
			propertyQualifiedName := qualifiedName + "." + propertyName
			if blindScope.GetSlot(propertyQualifiedName) == nil {
				unknownType := s.typeRegistry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
				return blindScope.InferQualifiedSlot(
					object, propertyQualifiedName, unknownType, unknownType, false)
			}
		}
	}
	return blindScope
}

// restrictByTrueInstanceOf restricts t when `x instanceof target` is truthy.
//
// Java: SemanticReverseAbstractInterpreter.RestrictByTrueInstanceOfResultVisitor
func (s *SemanticReverseAbstractInterpreter) restrictByTrueInstanceOf(t rhino.JSType, target rhino.ObjectType) rhino.JSType {
	switch x := t.(type) {
	case *rhino.NoType:
		return x
	case *rhino.AllType:
		return applyCommonInstanceOfRestriction(s.GetNativeType(rhino.JSTypeNative_ALL_TYPE), target)
	case *rhino.UnknownType:
		funcTarget := rhino.JSType_ToMaybeFunctionType(target)
		if funcTarget != nil && funcTarget.HasInstanceType() {
			return funcTarget.GetInstanceType()
		}
		return s.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
	case *rhino.UnionType:
		return applyCommonInstanceOfRestriction(x, target)
	case *rhino.FunctionTypeBase, *rhino.InstanceObjectType, *rhino.PrototypeObjectTypeBase,
		*rhino.RecordType, *rhino.EnumType, *rhino.NoResolvedType, *rhino.TemplatizedType, *rhino.TemplateType:
		return applyCommonInstanceOfRestriction(t, target)
	case *rhino.EnumElementType:
		prim := x.GetPrimitiveType()
		r := s.restrictByTrueInstanceOf(prim, target)
		if r != nil && prim.Equals(r) {
			return x
		}
		return r
	case *rhino.NamedType:
		return s.restrictByTrueInstanceOf(x.GetReferencedTypeInternal(), target)
	case *rhino.ProxyObjectTypeBase:
		return s.restrictByTrueInstanceOf(x.GetReferencedTypeInternal(), target)
	case *rhino.NoObjectTypeBase, *rhino.BooleanType, *rhino.NullType,
		*rhino.NumberType, *rhino.BigIntType, *rhino.StringType,
		*rhino.SymbolTypeBase, *rhino.KnownSymbolType, *rhino.VoidType:
		return nil
	default:
		if _, ok := t.(rhino.NoObjectType); ok {
			return nil
		}
		if _, ok := t.(rhino.FunctionType); ok {
			return applyCommonInstanceOfRestriction(t, target)
		}
		if proxy, ok := t.(rhino.ProxyObjectType); ok {
			return s.restrictByTrueInstanceOf(proxy.GetReferencedTypeInternal(), target)
		}
		if _, ok := t.(rhino.ObjectType); ok {
			return applyCommonInstanceOfRestriction(t, target)
		}
		if _, ok := t.(rhino.SymbolType); ok {
			return nil
		}
		panic(jutil.NewUnsupportedOperationException())
	}
}

// Java: RestrictByTrueInstanceOfResultVisitor.applyCommonRestriction(JSType)
func applyCommonInstanceOfRestriction(type_ rhino.JSType, target rhino.ObjectType) rhino.JSType {
	if target.IsUnknownType() {
		return type_
	}
	funcTarget := target.ToMaybeFunctionType()
	if funcTarget.HasInstanceType() {
		return type_.GetGreatestSubtype(funcTarget.GetInstanceType())
	}
	return nil
}

// restrictByFalseInstanceOf restricts t when `x instanceof target` is falsy.
//
// Java: SemanticReverseAbstractInterpreter.RestrictByFalseInstanceOfResultVisitor
func (s *SemanticReverseAbstractInterpreter) restrictByFalseInstanceOf(t rhino.JSType, target rhino.ObjectType) rhino.JSType {
	switch x := t.(type) {
	case *rhino.NoType:
		return x
	case *rhino.AllType:
		return s.GetNativeType(rhino.JSTypeNative_ALL_TYPE)
	case *rhino.UnknownType:
		return s.GetNativeType(rhino.JSTypeNative_CHECKED_UNKNOWN_TYPE)
	case *rhino.NoObjectTypeBase:
		return s.GetNativeType(rhino.JSTypeNative_NO_OBJECT_TYPE)
	case *rhino.BooleanType:
		return s.GetNativeType(rhino.JSTypeNative_BOOLEAN_TYPE)
	case *rhino.NullType:
		return s.GetNativeType(rhino.JSTypeNative_NULL_TYPE)
	case *rhino.NumberType:
		return s.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)
	case *rhino.BigIntType:
		return s.GetNativeType(rhino.JSTypeNative_BIGINT_TYPE)
	case *rhino.StringType:
		return s.GetNativeType(rhino.JSTypeNative_STRING_TYPE)
	case *rhino.VoidType:
		return s.GetNativeType(rhino.JSTypeNative_VOID_TYPE)
	case *rhino.SymbolTypeBase, *rhino.KnownSymbolType:
		return t
	case *rhino.UnionType:
		if target.IsUnknownType() {
			return x
		}
		funcTarget := target.ToMaybeFunctionType()
		if funcTarget.HasInstanceType() {
			return x.GetRestrictedUnion(funcTarget.GetInstanceType())
		}
		return nil
	case *rhino.FunctionTypeBase, *rhino.InstanceObjectType, *rhino.PrototypeObjectTypeBase,
		*rhino.RecordType, *rhino.EnumType, *rhino.NoResolvedType, *rhino.TemplatizedType, *rhino.TemplateType:
		return restrictObjectByFalseInstanceOf(t.(rhino.ObjectType), target)
	case *rhino.EnumElementType:
		prim := x.GetPrimitiveType()
		r := s.restrictByFalseInstanceOf(prim, target)
		if r != nil && prim.Equals(r) {
			return x
		}
		return r
	case *rhino.NamedType:
		return s.restrictByFalseInstanceOf(x.GetReferencedTypeInternal(), target)
	case *rhino.ProxyObjectTypeBase:
		return s.restrictByFalseInstanceOf(x.GetReferencedTypeInternal(), target)
	default:
		if _, ok := t.(rhino.NoObjectType); ok {
			return s.GetNativeType(rhino.JSTypeNative_NO_OBJECT_TYPE)
		}
		if fn, ok := t.(rhino.FunctionType); ok {
			return restrictObjectByFalseInstanceOf(fn, target)
		}
		if proxy, ok := t.(rhino.ProxyObjectType); ok {
			return s.restrictByFalseInstanceOf(proxy.GetReferencedTypeInternal(), target)
		}
		if obj, ok := t.(rhino.ObjectType); ok {
			return restrictObjectByFalseInstanceOf(obj, target)
		}
		if _, ok := t.(rhino.SymbolType); ok {
			return t
		}
		panic(jutil.NewUnsupportedOperationException())
	}
}

func restrictObjectByFalseInstanceOf(type_ rhino.ObjectType, target rhino.ObjectType) rhino.JSType {
	if target.IsUnknownType() {
		return type_
	}
	funcTarget := target.ToMaybeFunctionType()
	if funcTarget.HasInstanceType() {
		if type_.IsSubtypeOf(funcTarget.GetInstanceType()) {
			return nil
		}
		return type_
	}
	return nil
}

// SemanticReverseAbstractInterpreter_Unwrap unwraps any
// RefinementTrackingFlowScopes.
//
// Java: SemanticReverseAbstractInterpreter.unwrap(FlowScope)
func SemanticReverseAbstractInterpreter_Unwrap(scope FlowScope) FlowScope {
	for {
		r, ok := scope.(*sraiTrackingFlowScope)
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

// sraiTrackingFlowScope is a wrapper around FlowScope that keeps track of
// which vars were refined.
type sraiTrackingFlowScope struct {
	owner      *SemanticReverseAbstractInterpreter
	delegate   FlowScope
	tracker    *sraiRefinementTracker
	ownTracker sraiRefinementTracker
}

func (s *SemanticReverseAbstractInterpreter) acquireTrackingScope(
	delegate FlowScope, tracker *sraiRefinementTracker,
) *sraiTrackingFlowScope {
	var r *sraiTrackingFlowScope
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
func (r *sraiTrackingFlowScope) WithSyntacticScope(scope rhino.StaticTypedScope) FlowScope {
	panic(jutil.NewUnsupportedOperationException())
}

// Java: RefinementTrackingFlowScope.inferSlotType(String, JSType)
func (r *sraiTrackingFlowScope) InferSlotType(symbol string, type_ rhino.JSType) FlowScope {
	r.tracker.add(symbol)
	return r.Wrap(r.delegate.InferSlotType(symbol, type_))
}

// Java: RefinementTrackingFlowScope.inferQualifiedSlot(Node, String, JSType, JSType, boolean)
func (r *sraiTrackingFlowScope) InferQualifiedSlot(
	node rhino.Node, symbol string, bottomType rhino.JSType, inferredType rhino.JSType, declare bool) FlowScope {
	r.tracker.add(symbol)
	return r.Wrap(r.delegate.InferQualifiedSlot(node, symbol, bottomType, inferredType, declare))
}

// Java: RefinementTrackingFlowScope.wrap(FlowScope)
func (r *sraiTrackingFlowScope) Wrap(scope FlowScope) FlowScope {
	if scope != r.delegate {
		return r.owner.acquireTrackingScope(scope, r.tracker)
	}
	return r
}

// Java: RefinementTrackingFlowScope.getDeclarationScope()
func (r *sraiTrackingFlowScope) GetDeclarationScope() rhino.StaticTypedScope {
	return r.delegate.GetDeclarationScope()
}

// Java: RefinementTrackingFlowScope.getRootNode()
func (r *sraiTrackingFlowScope) GetRootNode() rhino.Node {
	return r.delegate.GetRootNode()
}

// Java: RefinementTrackingFlowScope.getParentScope()
func (r *sraiTrackingFlowScope) GetParentScope() rhino.StaticScope {
	panic(jutil.NewUnsupportedOperationException())
}

// Java: RefinementTrackingFlowScope.getSlot(String)
func (r *sraiTrackingFlowScope) GetSlot(name string) rhino.StaticSlot {
	return r.delegate.GetSlot(name)
}

// Java: RefinementTrackingFlowScope.getOwnSlot(String)
func (r *sraiTrackingFlowScope) GetOwnSlot(name string) rhino.StaticSlot {
	return r.delegate.GetOwnSlot(name)
}

// Java: RefinementTrackingFlowScope.getTypeOfThis()
func (r *sraiTrackingFlowScope) GetTypeOfThis() rhino.JSType {
	return r.delegate.GetTypeOfThis()
}

// GetTopmostScopeOfEventualDeclaration is the inherited default method.
//
// Java: StaticScope.getTopmostScopeOfEventualDeclaration(String) (default method)
func (r *sraiTrackingFlowScope) GetTopmostScopeOfEventualDeclaration(name string) rhino.StaticScope {
	return rhino.StaticScope_GetTopmostScopeOfEventualDeclaration(r, name)
}

// LookupQualifiedName is the inherited default method.
//
// Java: StaticTypedScope.lookupQualifiedName(QualifiedName) (default method)
func (r *sraiTrackingFlowScope) LookupQualifiedName(qname rhino.QualifiedName) rhino.JSType {
	return rhino.StaticTypedScope_LookupQualifiedName(r, qname)
}
