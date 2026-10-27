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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// TypeInference_BigIntPresence: now that BigInt is a factor, we need a
// better understanding of which types are involved in any given operation.
// In most cases, simply knowing that BigInt is involved is not enough
// information to determine that we should report an error, for example. So,
// we have designed this enumeration (and a respective utility function
// below) to aid in providing more context for cases that are more complex.
// Currently, there are four possibilities that concern us. This is subject to
// change as we add more support for type inference with BigInts.
type TypeInference_BigIntPresence int32

const (
	TypeInference_BigIntPresence_NO_BIGINT TypeInference_BigIntPresence = iota
	TypeInference_BigIntPresence_ALL_BIGINT
	TypeInference_BigIntPresence_BIGINT_OR_NUMBER
	TypeInference_BigIntPresence_BIGINT_OR_OTHER
)

var typeInference_BigIntPresence_names = [...]string{
	"NO_BIGINT", "ALL_BIGINT", "BIGINT_OR_NUMBER", "BIGINT_OR_OTHER",
}

// Name is Java's Enum.name().
func (p TypeInference_BigIntPresence) Name() string { return typeInference_BigIntPresence_names[p] }

// String is Java's Enum.toString().
func (p TypeInference_BigIntPresence) String() string { return p.Name() }

// Ordinal is Java's Enum.ordinal().
func (p TypeInference_BigIntPresence) Ordinal() int { return int(p) }

// TypeInference_GetBigIntPresence is a utility function for determining the
// BigIntPresence for any type.
//
// Java: TypeInference.getBigIntPresence(JSType)
func TypeInference_GetBigIntPresence(type_ rhino.JSType) TypeInference_BigIntPresence {
	// Base case
	if type_.IsOnlyBigInt() {
		return TypeInference_BigIntPresence_ALL_BIGINT
	}

	// Checking enum case
	typeAsEnumElement := type_.ToMaybeEnumElementType()
	if typeAsEnumElement != nil {
		// No matter what type the enum element is, this function can resolve it to a BigIntPresence
		return TypeInference_GetBigIntPresence(typeAsEnumElement.GetPrimitiveType())
	}

	// Union case
	typeAsUnion := type_.ToMaybeUnionType()
	if typeAsUnion != nil {
		containsBigInt := false
		containsNumber := false
		containsOther := false
		for _, alternate := range typeAsUnion.GetAlternates() {
			if TypeInference_GetBigIntPresence(alternate) != TypeInference_BigIntPresence_NO_BIGINT {
				containsBigInt = true
			} else if alternate.IsNumber() && !alternate.IsUnknownType() {
				containsNumber = true
			} else {
				containsOther = true
			}
		}
		if containsBigInt {
			if containsOther {
				return TypeInference_BigIntPresence_BIGINT_OR_OTHER
			} else if containsNumber {
				return TypeInference_BigIntPresence_BIGINT_OR_NUMBER
			} else {
				return TypeInference_BigIntPresence_ALL_BIGINT
			}
		}
	}

	// If it’s not a bigint object, bigint value, enum containing either, or a union, then we can
	// safely assume that it’s not a bigint in anyway
	return TypeInference_BigIntPresence_NO_BIGINT
}

// TraverseUnaryPlus checks for BigInt with a unary plus.
//
// Java: TypeInference.traverseUnaryPlus(Node, FlowScope)
func (ti *TypeInference) TraverseUnaryPlus(n rhino.Node, scope FlowScope) FlowScope {
	scope = ti.TraverseChildren(n, scope) // Find types.
	bigintPresenceInOperand := TypeInference_GetBigIntPresence(ti.GetJSType(n.GetFirstChild()))
	if bigintPresenceInOperand != TypeInference_BigIntPresence_NO_BIGINT {
		// Unary plus throws an exception when applied to a bigint
		n.SetJSType(ti.GetNativeType(rhino.JSTypeNative_NO_TYPE))
	} else {
		n.SetJSType(ti.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE))
	}
	return scope
}

// TraverseBigIntCompatibleUnaryOperator traverses unary minus, bitwise NOT,
// increment, and decrement.
//
// Java: TypeInference.traverseBigIntCompatibleUnaryOperator(Node, FlowScope)
func (ti *TypeInference) TraverseBigIntCompatibleUnaryOperator(n rhino.Node, scope FlowScope) FlowScope {
	scope = ti.TraverseChildren(n, scope) // Find types.
	switch TypeInference_GetBigIntPresence(ti.GetJSType(n.GetFirstChild())) {
	case TypeInference_BigIntPresence_ALL_BIGINT:
		// BigIntPresence in operand
		n.SetJSType(ti.GetNativeType(rhino.JSTypeNative_BIGINT_TYPE))
	case TypeInference_BigIntPresence_NO_BIGINT:
		n.SetJSType(ti.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE))
	case TypeInference_BigIntPresence_BIGINT_OR_NUMBER, TypeInference_BigIntPresence_BIGINT_OR_OTHER:
		n.SetJSType(ti.GetNativeType(rhino.JSTypeNative_BIGINT_NUMBER))
	}
	return scope
}

// TraverseBigIntCompatibleBinaryOperator traverses all binary numeric
// operations (except for URSH and string concatenation with ADD).
//
// Java: TypeInference.traverseBigIntCompatibleBinaryOperator(Node, FlowScope)
func (ti *TypeInference) TraverseBigIntCompatibleBinaryOperator(n rhino.Node, scope FlowScope) FlowScope {
	scope = ti.TraverseChildren(n, scope) // Find types.
	leftBigIntPresence := TypeInference_GetBigIntPresence(ti.GetJSType(n.GetFirstChild()))
	rightBigIntPresence := TypeInference_GetBigIntPresence(ti.GetJSType(n.GetLastChild()))
	if leftBigIntPresence != rightBigIntPresence {
		// disallowed mixture of bigint and non-bigint operands
		n.SetJSType(ti.GetNativeType(rhino.JSTypeNative_NO_TYPE))
	} else {
		switch leftBigIntPresence {
		case TypeInference_BigIntPresence_NO_BIGINT:
			n.SetJSType(ti.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE))
		case TypeInference_BigIntPresence_ALL_BIGINT:
			n.SetJSType(ti.GetNativeType(rhino.JSTypeNative_BIGINT_TYPE))
		case TypeInference_BigIntPresence_BIGINT_OR_NUMBER:
			n.SetJSType(ti.GetNativeType(rhino.JSTypeNative_BIGINT_NUMBER))
		case TypeInference_BigIntPresence_BIGINT_OR_OTHER:
			// In the case of arithmetic operations, BigInts are only compatible with other BigInts.
			// So if bigint is in a union with anything but number (and even then they both have to
			// be
			// {bigint|number}), then an error is reported.
			n.SetJSType(ti.GetNativeType(rhino.JSTypeNative_NO_TYPE))
		}
	}
	return scope
}

// Java: TypeInference.traverseClass(Node, FlowScope)
func (ti *TypeInference) TraverseClass(n rhino.Node, scope FlowScope) FlowScope {
	// The name already has a type applied (from TypedScopeCreator) if it's non-empty, and the
	// members are traversed in the class scope (and in their own function scopes).  But the extends
	// clause and computed property keys are in the outer scope and must be traversed here.
	scope = ti.Traverse(n.GetSecondChild(), scope)
	classMembers := NodeUtil_GetClassMembers(n)

	for member := classMembers.GetFirstChild(); !member.IsNil(); member = member.GetNext() {
		// Computed properties LHS need to happen before any RHS values
		if member.IsComputedProp() || member.IsComputedFieldDef() {
			scope = ti.Traverse(member.GetFirstChild(), scope)
		}
	}
	for member := classMembers.GetFirstChild(); !member.IsNil(); member = member.GetNext() {
		scope = ti.TraverseClassMemberRhs(member, scope)
	}
	return scope
}

// Java: TypeInference.traverseClassMemberRhs(Node, FlowScope)
func (ti *TypeInference) TraverseClassMemberRhs(member rhino.Node, scope FlowScope) FlowScope {
	switch member.GetToken() {
	case rhino.Token_MEMBER_FIELD_DEF, rhino.Token_COMPUTED_FIELD_DEF:
		rhs := TypeInference_GetRhsOfField(member)
		if !rhs.IsNil() {
			computedFieldDefTypedScope := ti.scopeCreator.CreateScopeOfNode(member)
			computedFieldDefFlowScope :=
				scope.WithSyntacticScope(computedFieldDefTypedScope)
			rhsScope :=
				ti.Traverse(rhs, computedFieldDefFlowScope).
					WithSyntacticScope(scope.GetDeclarationScope())
			if member.IsStaticMember() {
				return rhsScope
			}
		}
		return scope
	case rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_COMPUTED_PROP, rhino.Token_BLOCK,
		rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF:
		return scope
	default:
		panic(jutil.NewAssertionErrorNoMessage())
	}
}

// Java: TypeInference.getRhsOfField(Node)
func TypeInference_GetRhsOfField(fieldNode rhino.Node) rhino.Node {
	switch fieldNode.GetToken() {
	case rhino.Token_MEMBER_FIELD_DEF:
		if fieldNode.HasOneChild() {
			return fieldNode.GetFirstChild()
		}
		return rhino.NilNode
	case rhino.Token_COMPUTED_FIELD_DEF:
		if fieldNode.HasTwoChildren() {
			return fieldNode.GetSecondChild()
		}
		return rhino.NilNode
	default:
		panic(jutil.NewAssertionErrorNoMessage())
	}
}

// TraverseArrayLiteral traverses each element of the array.
//
// Java: TypeInference.traverseArrayLiteral(Node, FlowScope)
func (ti *TypeInference) TraverseArrayLiteral(n rhino.Node, scope FlowScope) FlowScope {
	scope = ti.TraverseChildren(n, scope)
	if ti.arrayOfUnknownType == nil {
		ti.arrayOfUnknownType = ti.registry.CreateTemplatizedType(
			ti.registry.GetNativeObjectType(rhino.JSTypeNative_ARRAY_TYPE), ti.unknownType)
	}
	n.SetJSType(ti.arrayOfUnknownType)
	return scope
}

// Java: TypeInference.traverseObjectLiteral(Node, FlowScope)
func (ti *TypeInference) TraverseObjectLiteral(n rhino.Node, scope FlowScope) FlowScope {
	type_ := n.GetJSType()
	jutil.CheckNotNull(type_)

	spreadOperatorSeen := false
	for name := n.GetFirstChild(); !name.IsNil(); name = name.GetNext() {
		spreadOperatorSeen = spreadOperatorSeen || name.IsSpread()
		scope = ti.TraverseChildren(name, scope)
	}

	// Object literals can be reflected on other types.
	// See CodingConvention#getObjectLiteralCast and goog.reflect.object
	// Ignore these types of literals.
	objectType := rhino.ObjectType_Cast(type_)
	if objectType == nil || n.GetBooleanProp(rhino.Node_REFLECTED_OBJECT) || objectType.IsEnumType() {
		return scope
	}

	var qObjName string
	var hasQObjName bool
	if !n.GetFirstChild().IsNil() {
		qObjName, hasQObjName = ti.getBestLValueNameString(NodeUtil_GetBestLValue(n))
	}
	for key := n.GetFirstChild(); !key.IsNil(); key = key.GetNext() {
		if key.IsComputedProp() {
			// Don't define computed properties as inferred properties on the object
			continue
		}

		if key.IsSpread() {
			// TODO(b/128355893): Do smarter inferrence. There are a lot of potential issues with
			// inference on object-spread, so for now we just give up and say `Object`.
			n.SetJSType(ti.registry.GetNativeType(rhino.JSTypeNative_OBJECT_TYPE))
			break
		}

		// Java tests `memberName != null` here and otherwise sets the unknown type
		// on n; NodeUtil.getObjectOrClassLitKeyName never returns null (it throws).
		memberName := NodeUtil_GetObjectOrClassLitKeyName(key)
		rawValueType := key.GetFirstChild().GetJSType()
		valueType := TypeCheck_GetObjectLitKeyTypeFromValueType(key, rawValueType)
		if valueType == nil {
			valueType = ti.unknownType
		}

		// See b/260837012. traverseObjectLiteral() can be invoked multiple times for the same
		// literal. When a spread operator is encountered, the type of the literal is changed to
		// OBJECT_TYPE. A second invocation of this code may then erroneously define new properties
		// on the global OBJECT_TYPE. So, if there are spread operators in the literal and the type
		// has already been propagated to OBJECT_TYPE, do not add newly inferred properties to the
		// type.
		if !spreadOperatorSeen || !n.GetJSType().IsNativeObjectType() {
			objectType.DefineInferredProperty(memberName, valueType, key)
		}

		// Do normal flow inference if this is a direct property assignment.
		if hasQObjName && key.IsStringKey() {
			qKeyName := qObjName + "." + memberName
			var_ := TypeInference_GetDeclaredVar(scope, qKeyName)
			var oldType rhino.JSType
			if var_ != nil {
				oldType = var_.GetType()
			}
			if var_ != nil && var_.IsTypeInferred() {
				if oldType == nil {
					var_.SetType(valueType)
				} else {
					var_.SetType(oldType.GetLeastSupertype(oldType))
				}
			}

			var bottomType rhino.JSType
			if oldType == nil {
				bottomType = ti.unknownType
			} else {
				bottomType = oldType
			}
			scope =
				scope.InferQualifiedSlot(
					key, qKeyName, bottomType, valueType, false)
		}
	}
	return scope
}

// Java: TypeInference.traverseAdd(Node, FlowScope)
func (ti *TypeInference) TraverseAdd(n rhino.Node, scope FlowScope) FlowScope {
	left := n.GetFirstChild()
	right := left.GetNext()
	scope = ti.TraverseChildren(n, scope)

	leftType := left.GetJSType()
	rightType := right.GetJSType()

	var type_ rhino.JSType = ti.unknownType
	if leftType != nil && rightType != nil {
		leftIsUnknown := leftType.IsUnknownType()
		rightIsUnknown := rightType.IsUnknownType()
		if (!leftIsUnknown && leftType.IsString()) || (!rightIsUnknown && rightType.IsString()) {
			type_ = ti.GetNativeType(rhino.JSTypeNative_STRING_TYPE)
		} else if TypeInference_GetBigIntPresence(leftType) != TypeInference_BigIntPresence_NO_BIGINT ||
			TypeInference_GetBigIntPresence(rightType) != TypeInference_BigIntPresence_NO_BIGINT {
			return ti.TraverseBigIntCompatibleBinaryOperator(n, scope)
		} else if leftIsUnknown || rightIsUnknown {
			type_ = ti.unknownType
		} else if ti.IsAddedAsNumber(leftType) && ti.IsAddedAsNumber(rightType) {
			type_ = ti.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)
		} else {
			type_ = ti.GetNativeType(rhino.JSTypeNative_NUMBER_STRING)
		}
	}
	n.SetJSType(type_)

	if n.IsAssignAdd() {
		scope =
			ti.UpdateScopeForAssignmentWithUpdateNode(
				scope, left, type_, rhino.NilNode /* updateNode= */, TypeInference_AssignmentType_ASSIGN)
	}

	return scope
}

// Java: TypeInference.isAddedAsNumber(JSType)
func (ti *TypeInference) IsAddedAsNumber(type_ rhino.JSType) bool {
	return type_.IsSubtypeOf(ti.numberAdditionSupertype)
}

// Java: TypeInference.traverseHook(Node, FlowScope)
func (ti *TypeInference) TraverseHook(n rhino.Node, scope FlowScope) FlowScope {
	condition := n.GetFirstChild()
	trueNode := condition.GetNext()
	falseNode := n.GetLastChild()

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
func (ti *TypeInference) TraverseNullishCoalesce(n rhino.Node, scope FlowScope) FlowScope {
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
func (ti *TypeInference) TraverseAnd(n rhino.Node, scope FlowScope) *TypeInference_BooleanOutcomePair {
	left := n.GetFirstChild()
	leftOutcome := ti.TraverseWithinShortCircuitingBinOp(left, scope)
	return ti.TraverseShortCircuitingBinOp(n, left, leftOutcome)
}

// Java: TypeInference.traverseChildren(Node, FlowScope)
func (ti *TypeInference) TraverseChildren(n rhino.Node, scope FlowScope) FlowScope {
	for el := n.GetFirstChild(); !el.IsNil(); el = el.GetNext() {
		scope = ti.Traverse(el, scope)
	}
	return scope
}

// Java: TypeInference.traverseGetElem(Node, FlowScope)
func (ti *TypeInference) TraverseGetElem(n rhino.Node, originalScope FlowScope) FlowScope {
	executedScope := ti.TraverseChildren(n, originalScope)
	return ti.SetGetElemNodeTypeAfterChildrenTraversed(n, executedScope)
}

// SetGetElemNodeTypeAfterChildrenTraversed sets the appropriate type on the
// GETELEM node `n` after its children have been traversed.
//
// n is the node that we want to finishTraversing; scopeAfterChildren the
// scope after children are traversed.
//
// Java: TypeInference.setGetElemNodeTypeAfterChildrenTraversed(Node, FlowScope)
func (ti *TypeInference) SetGetElemNodeTypeAfterChildrenTraversed(n rhino.Node, scopeAfterChildren FlowScope) FlowScope {
	jutil.CheckArgument(n.IsGetElem() || n.IsOptChainGetElem())
	ti.InferGetElemType(n)
	scopeAfterChildren = ti.TightenTypeAfterDereference(n.GetFirstChild(), scopeAfterChildren)
	return scopeAfterChildren
}

// Java: TypeInference.inferGetElemType(Node)
func (ti *TypeInference) InferGetElemType(n rhino.Node) {
	objType := ti.GetJSType(n.GetFirstChild()).RestrictByNotNullOrUndefined()
	indexKey := n.GetLastChild()
	indexType := ti.GetJSType(indexKey)

	var inferredType rhino.JSType
	if indexType.IsKnownSymbolValueType() {
		inferredType = ti.DereferenceKnownSymbolProp(objType, indexType.ToMaybeKnownSymbolType())
	} else if indexType.IsSymbolValueType() {
		// For now, allow symbols definitions/access on any type. In the future only allow them
		// on the subtypes for which they are defined.
		// TODO(b/77474174): be stricter about accesses for non-well-known symbols
		inferredType = ti.unknownType
	} else {
		inferredType = ti.DereferenceIndexSignature(objType)
	}
	if inferredType != nil {
		n.SetJSType(inferredType)
	} else {
		n.SetJSType(ti.unknownType)
	}
}

// Java: TypeInference.dereferenceKnownSymbolProp(JSType, KnownSymbolType)
func (ti *TypeInference) DereferenceKnownSymbolProp(obj rhino.JSType, indexSymbol *rhino.KnownSymbolType) rhino.JSType {
	if !obj.IsUnionType() {
		objType := obj.ToMaybeObjectType()
		if objType == nil {
			return nil
		}
		propertyType := objType.FindPropertyTypeByKey(rhino.NewProperty_SymbolKey(indexSymbol))
		if propertyType == nil {
			return nil
		}
		if !propertyType.IsUnionType() && !propertyType.IsNoType() {
			return propertyType
		}
		argumentTypes := rhino.UnionType_Builder_(ti.registry)
		argumentTypes.AddAlternate(propertyType)
		return argumentTypes.Build()
	}

	// If this is a union type, then we must extract type arguments from each option.
	argumentTypes := rhino.UnionType_Builder_(ti.registry)
	alternates := obj.ToMaybeUnionType().GetAlternates()
	key := rhino.NewProperty_SymbolKey(indexSymbol)
	for _, option := range alternates {
		if option.ToMaybeObjectType() == nil {
			// This isn't an array or object, drop out.
			argumentTypes = nil
			break
		}

		// Extract the element type and add all options to our set of alternates.
		propertyType := option.ToMaybeObjectType().FindPropertyTypeByKey(key)
		if propertyType == nil {
			// this union member doesn't have the property. just make it unknown.
			argumentTypes = nil
			break
		}
		argumentTypes.AddAlternate(propertyType)
	}

	// Unwrap the union if possible, and fail if we had no alternates.
	if argumentTypes == nil {
		return nil
	}
	return argumentTypes.Build()
}

// Java: TypeInference.dereferenceIndexSignature(JSType)
func (ti *TypeInference) DereferenceIndexSignature(obj rhino.JSType) rhino.JSType {
	if !obj.IsUnionType() {
		typeMap := obj.GetTemplateTypeMap()
		if !typeMap.HasTemplateType(ti.registry.GetObjectElementKey()) {
			return nil
		}
		resolved := typeMap.GetResolvedTemplateType(ti.registry.GetObjectElementKey())
		if !resolved.IsUnionType() && !resolved.IsNoType() {
			return resolved
		}
		argumentTypes := rhino.UnionType_Builder_(ti.registry)
		argumentTypes.AddAlternate(resolved)
		return argumentTypes.Build()
	}

	// If this is a union type, then we must extract type arguments from each option.
	argumentTypes := rhino.UnionType_Builder_(ti.registry)
	alternates := obj.ToMaybeUnionType().GetAlternates()
	for _, option := range alternates {
		typeMap := option.GetTemplateTypeMap()
		if !typeMap.HasTemplateType(ti.registry.GetObjectElementKey()) {
			// This isn't an array or object, drop out.
			argumentTypes = nil
			break
		}

		// Extract the element type and add all options to our set of alternates.
		argumentTypes.AddAlternate(typeMap.GetResolvedTemplateType(ti.registry.GetObjectElementKey()))
	}

	// Unwrap the union if possible, and fail if we had no alternates.
	if argumentTypes == nil {
		return nil
	}
	return argumentTypes.Build()
}

// Java: TypeInference.traverseGetProp(Node, FlowScope)
func (ti *TypeInference) TraverseGetProp(n rhino.Node, scope FlowScope) FlowScope {
	scope = ti.TraverseChildren(n, scope)
	return ti.SetGetPropNodeTypeAfterChildrenTraversed(n, scope)
}

// SetGetPropNodeTypeAfterChildrenTraversed sets the appropriate type on the
// GETPROP node `n` after its children have been traversed.
//
// Java: TypeInference.setGetPropNodeTypeAfterChildrenTraversed(Node, FlowScope)
func (ti *TypeInference) SetGetPropNodeTypeAfterChildrenTraversed(n rhino.Node, scopeAfterChildren FlowScope) FlowScope {
	jutil.CheckArgument(n.IsGetProp() || n.IsOptChainGetProp())
	objNode := n.GetFirstChild()
	n.SetJSType(ti.GetPropertyType(objNode.GetJSType(), n.GetString(), n, scopeAfterChildren))
	return ti.TightenTypeAfterDereference(n.GetFirstChild(), scopeAfterChildren)
}

// SetOptChainNodeTypeAfterChildrenTraversed sets the appropriate type on the
// OptChain node `n` after its children have been traversed.
//
// Java: TypeInference.setOptChainNodeTypeAfterChildrenTraversed(Node, FlowScope)
func (ti *TypeInference) SetOptChainNodeTypeAfterChildrenTraversed(
	n rhino.Node, scopeAfterChildren FlowScope) FlowScope {
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
	unconditionalScope FlowScope
}

// Java: TypeInference.OptChainInfo(Node, Node)
func NewTypeInference_OptChainInfo(endOfChain rhino.Node, startOfChain rhino.Node) *TypeInference_OptChainInfo {
	o := &TypeInference_OptChainInfo{}
	o.endOfChain = endOfChain
	o.startOfChain = startOfChain
	o.unconditionalScope = nil
	return o
}

// TraverseOptChain: traversal requires holding scopes from the
// unconditional start (lhs of the `?.`) and the conditional part (rhs of the
// `?.`), and using the startType to determine the resulting scope.
//
// Java: TypeInference.traverseOptChain(Node, FlowScope)
func (ti *TypeInference) TraverseOptChain(n rhino.Node, scope FlowScope) FlowScope {
	jutil.CheckArgument(NodeUtil_IsOptChainNode(n))

	if NodeUtil_IsEndOfOptChainSegment(n) {
		// Create new optional chain tracking object and push it onto the stack.
		startOfChain := NodeUtil_GetStartOfOptChainSegment(n)
		optChainInfo := NewTypeInference_OptChainInfo(n, startOfChain)
		// optChainArrayDeque.addFirst(optChainInfo)
		ti.optChainArrayDeque = append(ti.optChainArrayDeque, optChainInfo)
	}

	lhsScope := ti.Traverse(n.GetFirstChild(), scope)

	if n.IsOptionalChainStart() {
		// Store lhsScope into top-of-stack as unexecuted (unconditional) scope.
		ti.optChainArrayDequePeekFirst().unconditionalScope = lhsScope
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
		currentChain := ti.optChainArrayDequeRemoveFirst()

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

// optChainArrayDequePeekFirst is ArrayDeque.peekFirst(): null when empty.
//
// Go-only: language-adapter — see optChainArrayDeque.
func (ti *TypeInference) optChainArrayDequePeekFirst() *TypeInference_OptChainInfo {
	if len(ti.optChainArrayDeque) == 0 {
		return nil
	}
	return ti.optChainArrayDeque[len(ti.optChainArrayDeque)-1]
}

// optChainArrayDequeRemoveFirst is ArrayDeque.removeFirst(): throws
// NoSuchElementException when empty.
//
// Go-only: language-adapter — see optChainArrayDeque.
func (ti *TypeInference) optChainArrayDequeRemoveFirst() *TypeInference_OptChainInfo {
	if len(ti.optChainArrayDeque) == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	first := ti.optChainArrayDeque[len(ti.optChainArrayDeque)-1]
	ti.optChainArrayDeque[len(ti.optChainArrayDeque)-1] = nil
	ti.optChainArrayDeque = ti.optChainArrayDeque[:len(ti.optChainArrayDeque)-1]
	return first
}

// UpdateTypeWhenEndOfOptChain uses the startNode's type to selectively join
// the executed scope with the unexecuted scope, and updates the type assigned
// to `optChain` in `setXAfterChildrenTraversed()`.
//
// Java: TypeInference.updateTypeWhenEndOfOptChain(Node, Node, OptChainInfo, FlowScope)
func (ti *TypeInference) UpdateTypeWhenEndOfOptChain(
	optChain rhino.Node, startNode rhino.Node, currentChain *TypeInference_OptChainInfo, executedScope FlowScope) FlowScope {

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
func (ti *TypeInference) TightenTypeAfterDereference(n rhino.Node, scope FlowScope) FlowScope {
	switch n.GetToken() {
	case rhino.Token_NAME:
		if n.GetString() == "" {
			return scope
		}
		type_ := ti.GetJSType(n)
		narrowed := type_.RestrictByNotNullOrUndefined()
		if type_ != narrowed && !type_.Equals(narrowed) {
			scope = ti.RedeclareSimpleVar(scope, n, narrowed)
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

func (ti *TypeInference) getBestLValueNameString(lValue rhino.Node) (string, bool) {
	if lValue.IsNil() || lValue.GetParent().IsObjectLit() {
		return "", false
	}
	if NodeUtil_IsGetOrSetKey(lValue) || lValue.IsComputedProp() {
		return "", false
	}
	if lValue.IsStringKey() || lValue.IsMemberFunctionDef() || lValue.IsMemberFieldDef() {
		objectLit := lValue.GetParent()
		owner, ok := ti.getBestLValueNameString(NodeUtil_GetBestLValue(objectLit))
		if ok {
			return owner + "." + NodeUtil_GetObjectOrClassLitKeyName(lValue), true
		}
		return "", false
	}
	return lValue.GetQualifiedNameString()
}

// Java: TypeInference.getPropertyType(JSType, String, Node, FlowScope)
func (ti *TypeInference) GetPropertyType(objType rhino.JSType, propName string, n rhino.Node, scope FlowScope) rhino.JSType {
	// We often have a couple of different types to choose from for the
	// property. Ordered by accuracy, we have
	// 1) A locally inferred qualified name (which is in the FlowScope)
	// 2) A globally declared qualified name (which is in the FlowScope)
	// 3) A property on the owner type (which is on objType)
	// 4) A name in the type registry (as a last resort)
	var propertyType rhino.JSType
	isLocallyInferred := false

	// Scopes sometimes contain inferred type info about qualified names.
	qualifiedName, hasQName := n.GetQualifiedNameString()
	var var_ rhino.StaticTypedSlot
	if hasQName {
		var_, _ = scope.GetSlot(qualifiedName).(rhino.StaticTypedSlot)
	}
	if var_ != nil {
		varType := var_.GetType()
		if varType != nil {
			isDeclared := !var_.IsTypeInferred()
			if !isDeclared || varType.Equals(ti.unknownType) {
				declaredVar := TypeInference_GetDeclaredVar(scope, qualifiedName)
				// Go-only: language-adapter — Java's `var != getDeclaredVar(...)` is an
				// identity comparison; a null declared var differs from var.
				isLocallyInferred = declaredVar == nil || var_ != rhino.StaticTypedSlot(declaredVar)
			}
			if isDeclared || isLocallyInferred {
				propertyType = varType
			}
		}
	}

	if propertyType == nil && objType != nil {
		foundType := objType.FindPropertyType(propName)
		if foundType != nil {
			propertyType = foundType
		}
	}

	if (propertyType == nil || propertyType.IsUnknownType()) && hasQName {
		// If we find this node in the registry, then we can infer its type.
		regType :=
			rhino.ObjectType_Cast(ti.registry.GetType(scope.GetDeclarationScope(), qualifiedName))
		if regType != nil {
			propertyType = regType.GetConstructor()
		}
	}

	if propertyType == nil {
		return ti.unknownType
	} else if propertyType.Equals(ti.unknownType) && isLocallyInferred {
		// If the type has been checked in this scope,
		// then use CHECKED_UNKNOWN_TYPE instead to indicate that.
		return ti.GetNativeType(rhino.JSTypeNative_CHECKED_UNKNOWN_TYPE)
	} else {
		return propertyType
	}
}

// Java: TypeInference.traverseOr(Node, FlowScope)
func (ti *TypeInference) TraverseOr(n rhino.Node, scope FlowScope) *TypeInference_BooleanOutcomePair {
	left := n.GetFirstChild()
	leftOutcome := ti.TraverseWithinShortCircuitingBinOp(left, scope)
	return ti.TraverseShortCircuitingBinOp(n, left, leftOutcome)
}

// Java: TypeInference.traverseShortCircuitingBinOpAssignment(Node, FlowScope)
func (ti *TypeInference) TraverseShortCircuitingBinOpAssignment(n rhino.Node, scope FlowScope) FlowScope {
	left := n.GetFirstChild()
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
		outcomeJoinedFlowScope, n.GetFirstChild(), n.GetJSType(), TypeInference_AssignmentType_ASSIGN)
}

// Java: TypeInference.traverseShortCircuitingBinOp(Node, Node, BooleanOutcomePair)
func (ti *TypeInference) TraverseShortCircuitingBinOp(
	n rhino.Node, left rhino.Node, leftOutcome *TypeInference_BooleanOutcomePair) *TypeInference_BooleanOutcomePair {
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
	var outcome *TypeInference_BooleanOutcomePair
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
func (ti *TypeInference) TraverseWithinShortCircuitingBinOp(n rhino.Node, scope FlowScope) *TypeInference_BooleanOutcomePair {
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

// Java: TypeInference.traverseAwait(Node, FlowScope)
func (ti *TypeInference) TraverseAwait(await rhino.Node, scope FlowScope) FlowScope {
	scope = ti.TraverseChildren(await, scope)

	expr := await.GetFirstChild()
	exprType := ti.GetJSType(expr)
	await.SetJSType(Promises_GetResolvedType(ti.registry, exprType))

	return scope
}

// Java: TypeInference.traverseDynamicImport(Node, FlowScope)
func (ti *TypeInference) TraverseDynamicImport(dynamicImport rhino.Node, scope FlowScope) FlowScope {
	templateType := ti.registry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)

	// If the module specifier is a string, attempt to resolve the module
	moduleMap := ti.compiler.GetModuleMap()
	importSpecifier := dynamicImport.GetFirstChild()
	input := ti.compiler.GetInput(NodeUtil_GetInputId(dynamicImport))
	if importSpecifier.IsStringLit() && moduleMap != nil && input != nil {
		targetPath :=
			input.
				GetPath().
				ResolveJsModule(
					importSpecifier.GetString(),
					*importSpecifier.GetSourceFileName(),
					importSpecifier.GetLineno(),
					importSpecifier.GetCharno())
		var targetModule *Module
		if targetPath != nil {
			targetModule = moduleMap.GetModuleByPath(targetPath)
		}
		if targetModule != nil {
			// TypedScopeCreator ensures that the MODULE_BODY type is the export namespace type
			scriptNode := targetModule.Metadata().RootNode()
			if scriptNode.HasOneChild() && scriptNode.GetFirstChild().IsModuleBody() {
				var exportNamespaceType rhino.JSType
				if !scriptNode.IsNil() {
					exportNamespaceType = scriptNode.GetOnlyChild().GetJSType()
				}
				if exportNamespaceType != nil {
					templateType = exportNamespaceType
				}
			} else {
				// Module transpilation has occurred before type inference so a MODULE_BODY node
				// no longer exists. Traverse the script to locate the module namespace variable
				// and retrieve the type from it.
				moduleName :=
					NodeUtil_FindPreorder(
						scriptNode,
						func(node rhino.Node) bool { return node.MatchesQualifiedName(targetPath.ToModuleName()) },
						func(rhino.Node) bool { return true })
				if !moduleName.IsNil() && moduleName.GetJSType() != nil {
					templateType = moduleName.GetJSType()
				}
			}
		}
	}

	dynamicImport.SetJSType(
		ti.registry.CreateTemplatizedType(
			ti.registry.GetNativeObjectType(rhino.JSTypeNative_PROMISE_TYPE), templateType))

	return ti.TraverseChildren(dynamicImport, scope)
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
	leftScope FlowScope

	// The scope when the whole expression executed.
	rightScope FlowScope

	// The scope when we don't know how much of the expression is executed.
	joinedScope FlowScope
}

// Java: TypeInference.BooleanOutcomePair(BooleanLiteralSet, BooleanLiteralSet, FlowScope, FlowScope)
func NewTypeInference_BooleanOutcomePair(
	outer *TypeInference,
	toBooleanOutcomes rhino.BooleanLiteralSet,
	booleanValues rhino.BooleanLiteralSet,
	leftScope FlowScope,
	rightScope FlowScope) *TypeInference_BooleanOutcomePair {
	return &TypeInference_BooleanOutcomePair{
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
func (p *TypeInference_BooleanOutcomePair) GetJoinedFlowScope() FlowScope {
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
func (p *TypeInference_BooleanOutcomePair) GetOutcomeFlowScope(nodeType rhino.Token, outcome bool) FlowScope {
	if (nodeType == rhino.Token_AND && outcome) || (nodeType == rhino.Token_OR && !outcome) {
		// We know that the whole expression must have executed.
		return p.rightScope
	} else {
		return p.GetJoinedFlowScope()
	}
}

// Java: TypeInference.newBooleanOutcomePair(JSType, FlowScope)
func (ti *TypeInference) NewBooleanOutcomePair(jsType rhino.JSType, flowScope FlowScope) *TypeInference_BooleanOutcomePair {
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
