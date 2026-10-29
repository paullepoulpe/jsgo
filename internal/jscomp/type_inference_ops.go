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
func (ti *TypeInference) TraverseUnaryPlus(n rhino.Node, scope *LinkedFlowScope) *LinkedFlowScope {
	scope = ti.TraverseChildren(n, scope) // Find types.
	bigintPresenceInOperand := TypeInference_GetBigIntPresence(ti.GetJSType(ast.MustUnary(n).Operand().Node()))
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
func (ti *TypeInference) TraverseBigIntCompatibleUnaryOperator(n rhino.Node, scope *LinkedFlowScope) *LinkedFlowScope {
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
func (ti *TypeInference) TraverseBigIntCompatibleBinaryOperator(n rhino.Node, scope *LinkedFlowScope) *LinkedFlowScope {
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
func (ti *TypeInference) TraverseClass(n rhino.Node, scope *LinkedFlowScope) *LinkedFlowScope {
	// The name already has a type applied (from TypedScopeCreator) if it's non-empty, and the
	// members are traversed in the class scope (and in their own function scopes).  But the extends
	// clause and computed property keys are in the outer scope and must be traversed here.
	cls := ast.MustClass(n)
	scope = ti.Traverse(cls.ExtendsNode(), scope)

	for member := range cls.Body().Members() {
		// Computed properties LHS need to happen before any RHS values
		if cm, ok := member.AsComputedMember(); ok {
			scope = ti.Traverse(cm.Key().Node(), scope)
		}
	}
	for member := range cls.Body().Members() {
		scope = ti.TraverseClassMemberRhs(member.Node(), scope)
	}
	return scope
}

// Java: TypeInference.traverseClassMemberRhs(Node, FlowScope)
func (ti *TypeInference) TraverseClassMemberRhs(member rhino.Node, scope *LinkedFlowScope) *LinkedFlowScope {
	switch member.GetToken() {
	case rhino.Token_MEMBER_FIELD_DEF, rhino.Token_COMPUTED_FIELD_DEF:
		rhs := TypeInference_GetRhsOfField(member)
		if !rhs.IsNil() {
			computedFieldDefTypedScope := ti.scopeCreator.CreateScopeOfNode(member)
			computedFieldDefFlowScope := scope.withTypedScope(computedFieldDefTypedScope)
			rhsScope := ti.Traverse(rhs, computedFieldDefFlowScope).withTypedScope(scope.syntacticScope)
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
	val, _ := ast.MustFieldDef(fieldNode).Value()
	return val.Node()
}

// TraverseArrayLiteral traverses each element of the array.
//
// Java: TypeInference.traverseArrayLiteral(Node, FlowScope)
func (ti *TypeInference) TraverseArrayLiteral(n rhino.Node, scope *LinkedFlowScope) *LinkedFlowScope {
	scope = ti.TraverseChildren(n, scope)
	if ti.arrayOfUnknownType == nil {
		ti.arrayOfUnknownType = ti.registry.CreateTemplatizedType(
			ti.registry.GetNativeObjectType(rhino.JSTypeNative_ARRAY_TYPE), ti.unknownType)
	}
	n.SetJSType(ti.arrayOfUnknownType)
	return scope
}

// Java: TypeInference.traverseObjectLiteral(Node, FlowScope)
func (ti *TypeInference) TraverseObjectLiteral(n rhino.Node, scope *LinkedFlowScope) *LinkedFlowScope {
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
func (ti *TypeInference) TraverseAdd(n rhino.Node, scope *LinkedFlowScope) *LinkedFlowScope {
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

// Java: TypeInference.traverseChildren(Node, FlowScope)
func (ti *TypeInference) TraverseChildren(n rhino.Node, scope *LinkedFlowScope) *LinkedFlowScope {
	for el := n.GetFirstChild(); !el.IsNil(); el = el.GetNext() {
		scope = ti.Traverse(el, scope)
	}
	return scope
}

// Java: TypeInference.traverseGetElem(Node, FlowScope)
func (ti *TypeInference) TraverseGetElem(n rhino.Node, originalScope *LinkedFlowScope) *LinkedFlowScope {
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
func (ti *TypeInference) SetGetElemNodeTypeAfterChildrenTraversed(n rhino.Node, scopeAfterChildren *LinkedFlowScope) *LinkedFlowScope {
	elem := ast.MustNormalOrOptChainGetElem(n)
	ti.InferGetElemType(n)
	scopeAfterChildren = ti.TightenTypeAfterDereference(elem.Object().Node(), scopeAfterChildren)
	return scopeAfterChildren
}

// Java: TypeInference.inferGetElemType(Node)
func (ti *TypeInference) InferGetElemType(n rhino.Node) {
	elem := ast.MustNormalOrOptChainGetElem(n)
	objType := ti.GetJSType(elem.Object().Node()).RestrictByNotNullOrUndefined()
	indexKey := elem.Index().Node()
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
func (ti *TypeInference) TraverseGetProp(n rhino.Node, scope *LinkedFlowScope) *LinkedFlowScope {
	scope = ti.TraverseChildren(n, scope)
	return ti.SetGetPropNodeTypeAfterChildrenTraversed(n, scope)
}

// SetGetPropNodeTypeAfterChildrenTraversed sets the appropriate type on the
// GETPROP node `n` after its children have been traversed.
//
// Java: TypeInference.setGetPropNodeTypeAfterChildrenTraversed(Node, FlowScope)
func (ti *TypeInference) SetGetPropNodeTypeAfterChildrenTraversed(n rhino.Node, scopeAfterChildren *LinkedFlowScope) *LinkedFlowScope {
	gp := ast.MustNormalOrOptChainGetProp(n)
	objNode := gp.Object().Node()
	n.SetJSType(ti.GetPropertyType(objNode.GetJSType(), gp.Property(), n, scopeAfterChildren))
	return ti.TightenTypeAfterDereference(objNode, scopeAfterChildren)
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
func (ti *TypeInference) GetPropertyType(objType rhino.JSType, propName string, n rhino.Node, scope *LinkedFlowScope) rhino.JSType {
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
	if hasQName {
		if declaredVar, slotVar, varType, isInferred, ok := scope.LookupSlotType(qualifiedName); ok && varType != nil {
			isDeclared := !isInferred
			if !isDeclared || varType.Equals(ti.unknownType) {
				// Go-only: language-adapter — Java's `var != getDeclaredVar(...)` is an
				// identity comparison; a null declared var differs from var.
				isLocallyInferred = declaredVar == nil || slotVar != declaredVar
			}
			if isDeclared || isLocallyInferred {
				propertyType = varType
			}
		}
	}

	if propertyType == nil && objType != nil {
		var foundType rhino.JSType
		if !n.IsNil() && (n.IsGetProp() || n.IsOptChainGetProp()) && n.GetString() == propName {
			foundType = rhino.FindPropertyTypeByID(objType, n.GetStringID())
		} else {
			foundType = objType.FindPropertyType(propName)
		}
		if foundType != nil {
			propertyType = foundType
		}
	}

	if (propertyType == nil || propertyType.IsUnknownType()) && hasQName {
		// If we find this node in the registry, then we can infer its type.
		regType :=
			rhino.ObjectType_Cast(ti.registry.GetType(scope.GetDeclarationScope(), qualifiedName))
		if regType != nil {
			if ctor := regType.GetConstructor(); ctor != nil {
				propertyType = ctor
			} else {
				propertyType = nil
			}
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

// Java: TypeInference.traverseAwait(Node, FlowScope)
func (ti *TypeInference) TraverseAwait(await rhino.Node, scope *LinkedFlowScope) *LinkedFlowScope {
	scope = ti.TraverseChildren(await, scope)

	expr := ast.MustAwait(await).Value().Node()
	exprType := ti.GetJSType(expr)
	await.SetJSType(Promises_GetResolvedType(ti.registry, exprType))

	return scope
}

// Java: TypeInference.traverseDynamicImport(Node, FlowScope)
func (ti *TypeInference) TraverseDynamicImport(dynamicImport rhino.Node, scope *LinkedFlowScope) *LinkedFlowScope {
	templateType := ti.registry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)

	// If the module specifier is a string, attempt to resolve the module
	moduleMap := ti.compiler.GetModuleMap()
	importSpecifier := ast.MustDynamicImport(dynamicImport).Specifier().Node()
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

// Java: TypeInference.traverseUnsignedRightShift(Node, FlowScope)
func (ti *TypeInference) TraverseUnsignedRightShift(n rhino.Node, scope *LinkedFlowScope) *LinkedFlowScope {
	scope = ti.TraverseChildren(n, scope)
	var left, right rhino.Node
	if bin, ok := ast.AsBinary(n); ok {
		left = bin.Left().Node()
		right = bin.Right().Node()
	} else {
		aop := ast.MustAssignOp(n)
		left = aop.Target().Node()
		right = aop.Value().Node()
	}
	if TypeInference_GetBigIntPresence(ti.GetJSType(left)) != TypeInference_BigIntPresence_NO_BIGINT ||
		TypeInference_GetBigIntPresence(ti.GetJSType(right)) != TypeInference_BigIntPresence_NO_BIGINT {
		// The spec does not allow for BigInts in an unsigned right shift, so we will report an error
		n.SetJSType(ti.GetNativeType(rhino.JSTypeNative_NO_TYPE))
	} else {
		n.SetJSType(ti.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE))
	}
	return scope
}

// Java: TypeInference.traverseYieldAll(Node, FlowScope)
func (ti *TypeInference) TraverseYieldAll(n rhino.Node, scope *LinkedFlowScope) *LinkedFlowScope {
	// A yield* expression will first yield all the elements of the given iterable, and then
	// evaluate to whatever the iterable returns when done.
	// The yielded type and done return type are not necessarily the same. Here, we look for the
	// done type - TReturn in Iterable<T, TReturn, TNext>.
	scope = ti.TraverseChildren(n, scope)
	val, _ := ast.MustYield(n).Value()
	innerType := ti.GetJSType(val.Node())
	yieldAllResult := JsIterables_GetReturnElementType(innerType, ti.registry)
	n.SetJSType(yieldAllResult)
	return scope
}
