// Port of com/google/javascript/jscomp/TypeCheck.java @ closure-compiler v20261021

/*
 * Copyright 2006 The Closure Compiler Authors.
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

// Java: TypeCheck.visitUnaryPlus(Node)
func (tc *TypeCheck) VisitUnaryPlus(op ast.Unary) {
	n := op.Node()
	if tc.GetJSType(n).IsNoType() {
		tc.Report(
			n,
			TypeCheck_UNARY_OPERATION,
			stringPtrOrNull(NodeUtil_OpToStr(op.Token())),
			tc.GetJSType(op.Operand().Node()).String())
	} else {
		tc.EnsureTypedWithNative(n, rhino.JSTypeNative_NUMBER_TYPE)
	}
}

// Java: TypeCheck.visitUnaryMinus(Node)
func (tc *TypeCheck) VisitUnaryMinus(op ast.Unary) {
	n := op.Node()
	operatorType := tc.GetJSType(n)
	childType := tc.GetJSType(op.Operand().Node())
	if operatorType.IsNumber() {
		// This condition is meant to catch any old cases (where bigint isn't involved)
		tc.validator.ExpectNumber(n, childType, "sign operator")
		tc.EnsureTypedWithNative(n, rhino.JSTypeNative_NUMBER_TYPE)
	} else {
		tc.validator.ExpectBigIntOrNumber(n, childType, "unary minus operator")
	}
}

// Java: TypeCheck.visitBitwiseNOT(Node)
func (tc *TypeCheck) VisitBitwiseNOT(op ast.Unary) {
	n := op.Node()
	operatorType := tc.GetJSType(n)
	childType := tc.GetJSType(op.Operand().Node())
	if operatorType.IsNumber() {
		// This condition is meant to catch any old cases (where bigint isn't involved)
		if !childType.MatchesNumberContext() {
			tc.Report(n, TypeCheck_BIT_OPERATION, stringPtrOrNull(NodeUtil_OpToStr(op.Token())), childType.String())
		} else {
			tc.validator.ExpectNumberStrict(n, childType, "bitwise NOT")
		}
		tc.EnsureTypedWithNative(n, rhino.JSTypeNative_NUMBER_TYPE)
	} else {
		tc.validator.ExpectBigIntOrNumber(n, childType, "bitwise NOT")
	}
}

// Java: TypeCheck.checkSpread(Node)
func (tc *TypeCheck) CheckSpread(spreadNode rhino.Node) {
	target := spreadNode.GetOnlyChild()
	tc.EnsureTyped(target)
	targetType := tc.GetJSType(target)

	switch spreadNode.GetParent().GetToken() {
	case rhino.Token_OBJECTLIT:
		// Case: `var x = {A: a, B: b, ...obj}`.
		// Nothing to check about object spread.
	case rhino.Token_ARRAYLIT, rhino.Token_CALL, rhino.Token_OPTCHAIN_CALL, rhino.Token_NEW:
		// Case: `var x = [a, b, ...itr]`
		// Case: `var x = fn(a, b, ...itr)`
		// Case: `var x = fn?.(a, b, ...itr)`
		tc.validator.ExpectAutoboxesToIterable(
			target, targetType, "Spread operator only applies to Iterable types")
	default:
		panic(jutil.NewIllegalStateExceptionMsg(
			"Unexpected parent of SPREAD: " + spreadNode.GetParent().ToStringTree()))
	}
}

// Java: TypeCheck.checkTypeofString(Node, String)
func (tc *TypeCheck) CheckTypeofString(n rhino.Node, s string) {
	if !(s == "number" ||
		s == "string" ||
		s == "boolean" ||
		s == "undefined" ||
		s == "function" ||
		s == "object" ||
		s == "symbol" ||
		s == "unknown" ||
		s == "bigint") {
		tc.validator.ExpectValidTypeofName(n, s)
	}
}

// VisitAssign visits an assignment `lvalue = rvalue`. If the `lvalue` is a
// prototype modification, we change the schema of the object type it is
// referring to.
//
// t is the traversal; assign the assign node (an `=`, `||=`, `&&=` or
// `??=`).
//
// Java: TypeCheck.visitAssign(NodeTraversal, Node)
func (tc *TypeCheck) VisitAssign(t *NodeTraversal, assign ast.Assignment) {
	n := assign.Node()
	info := n.GetJSDocInfo()
	lvalue := assign.Target()
	rvalue := assign.Value()

	rightType := tc.GetJSType(rvalue.Node())
	tc.CheckCanAssignToWithScope(t, n, lvalue, rightType, info, "assignment")
	tc.EnsureTypedWithType(n, rightType)
}

// CheckCanAssignToWithScope checks that we can assign the given right type to
// the given lvalue or destructuring pattern.
//
// See CheckCanAssignToNameGetpropOrGetelem for more details.
//
// nodeToWarn is a node to report type mismatch warnings on; lvalue the lvalue
// to which we're assigning or a destructuring pattern; rightType the type
// we're assigning to the lvalue; msg a message to report along with any type
// mismatch warnings.
//
// Java: TypeCheck.checkCanAssignToWithScope(NodeTraversal, Node, Node, JSType, JSDocInfo, String)
func (tc *TypeCheck) CheckCanAssignToWithScope(
	t *NodeTraversal,
	nodeToWarn rhino.Node,
	lvalue ast.Target,
	rightType rhino.JSType,
	info *rhino.JSDocInfo,
	msg string) {
	if pattern, ok := lvalue.AsPattern(); ok {
		tc.CheckDestructuringAssignment(t, nodeToWarn, pattern, rightType, msg)
	} else {
		tc.CheckCanAssignToNameGetpropOrGetelem(t, nodeToWarn, lvalue, rightType, info, msg)
	}
}

// CheckDestructuringAssignment recursively checks that an assignment to a
// destructuring pattern is valid for all the lvalues contained in the pattern
// (including in nested patterns).
//
// Java: TypeCheck.checkDestructuringAssignment(NodeTraversal, Node, Node, JSType, String)
func (tc *TypeCheck) CheckDestructuringAssignment(
	t *NodeTraversal, nodeToWarn rhino.Node, pattern ast.Pattern, rightType rhino.JSType, msg string) {
	ForEachDestructuredTarget(tc.typeRegistry, rightType, pattern.Node(), func(target DestructuredTarget) {
		// TODO(b/77597706): this is not very efficient because it re-infers the types below,
		// which we already did once in TypeInference. don't repeat the work.
		tc.CheckCanAssignToWithScope(
			t, nodeToWarn, ast.MustTarget(target.GetNode()), target.InferType() /* info= */, nil, msg)
	})
}

// CheckCanAssignToNameGetpropOrGetelem checks that we can assign the given
// right type to the given lvalue.
//
// If the lvalue is a qualified name, and has a declared type in the given
// scope, uses the declared type of the qualified name instead of the type on
// the node.
//
// nodeToWarn is a node to report type mismatch warnings on; lvalue the lvalue
// to which we're assigning - a NAME, GETELEM, or GETPROP; rightType the type
// we're assigning to the lvalue; msg a message to report along with any type
// mismatch warnings.
//
// Java: TypeCheck.checkCanAssignToNameGetpropOrGetelem(NodeTraversal, Node, Node, JSType, JSDocInfo, String)
func (tc *TypeCheck) CheckCanAssignToNameGetpropOrGetelem(
	t *NodeTraversal, nodeToWarn rhino.Node, lvalue ast.Target, rightType rhino.JSType, info *rhino.JSDocInfo, msg string) {
	_, isPattern := lvalue.AsPattern()
	jutil.CheckArgumentMsg1(!isPattern, "%s", lvalue.Node())

	// Ensure our LHS is not readonly.
	tc.CheckNotReadonlyPropertyAssignment(lvalue)

	if getProp, ok := lvalue.AsGetProp(); ok {
		object := getProp.Object()
		objectJsType := tc.GetJSType(object.Node())
		getPropNode := getProp.Node()

		// the first name in this getprop refers to an interface
		// we perform checks in addition to the ones below
		if objectGetProp, ok := object.AsGetProp(); ok {
			jsType := tc.GetJSType(objectGetProp.Object().Node())
			if jsType.IsInterface() && objectGetProp.Node().GetStringID() == typeCheckPrototypeID {
				tc.VisitInterfacePropertyAssignment(object, getProp)
			}
		}

		tc.CheckEnumAlias(t, info, rightType, nodeToWarn)
		tc.CheckPropCreation(lvalue)

		// Prototype assignments are special, because they actually affect
		// the definition of a class. These are mostly validated
		// during TypedScopeCreator, and we only look for the "dumb" cases here.
		// object.prototype = ...;
		if getPropNode.GetStringID() == typeCheckPrototypeID {
			tc.validator.ExpectCanAssignToPrototype(objectJsType, nodeToWarn, rightType)
			return
		}

		pname := getProp.Property()

		// The generic checks for 'object.property' when 'object' is known,
		// and 'property' is declared on it.
		// object.property = ...;
		objectCastType := rhino.ObjectType_Cast(objectJsType.RestrictByNotNullOrUndefined())
		expectedPropertyType := tc.GetPropertyTypeIfDeclared(objectCastType, rhino.NewProperty_StringKey(pname))

		tc.CheckPropertyInheritanceOnAssignment(
			nodeToWarn, object, rhino.NewProperty_StringKey(pname), info, expectedPropertyType)

		// If we successfully found a non-unknown declared type, validate the assignment and don't do
		// any further checks.
		if !expectedPropertyType.IsUnknownType() {
			// Note: if the property has @implicitCast at its declaration, we don't check any
			// assignments to it.
			if !TypeCheck_PropertyIsImplicitCast(objectCastType, pname) {
				tc.validator.ExpectCanAssignToPropertyOf(
					nodeToWarn, rightType, expectedPropertyType, object.Node(), pname)
			}
			return
		}
	} else if getElem, ok := lvalue.AsGetElem(); ok && getElem.Index().Node().GetJSType().IsKnownSymbolValueType() {
		object := getElem.Object()
		objectJsType := tc.GetJSType(object.Node())
		objectCastType := rhino.ObjectType_Cast(objectJsType.RestrictByNotNullOrUndefined())
		property := getElem.Index().Node().GetJSType().ToMaybeKnownSymbolType()
		expectedPropertyType :=
			tc.GetPropertyTypeIfDeclared(objectCastType, rhino.NewProperty_SymbolKey(property))
		tc.CheckPropertyInheritanceOnAssignment(
			nodeToWarn, object, rhino.NewProperty_SymbolKey(property), info, expectedPropertyType)
	}

	// Check qualified name sets to 'object' and 'object.property'.
	// This can sometimes handle cases when the type of 'object' is not known.
	// e.g.,
	// var obj = createUnknownType();
	// /** @type {number} */ obj.foo = true;
	leftType := tc.GetJSType(lvalue.Node())
	if qName, ok := lvalue.Node().GetQualifiedNameString(); ok {
		// variable with inferred type case
		var_ := t.GetTypedScope().GetVar(qName)
		if var_ != nil {
			if var_.IsTypeInferred() {
				return
			}

			if NodeUtil_GetRootOfQualifiedName(lvalue.Node()).IsThis() &&
				rhino.StaticScope(t.GetTypedScope()) != var_.GetScope() {
				// Don't look at "this.foo" variables from other scopes.
				return
			}

			if var_.GetType() != nil {
				leftType = var_.GetType()
			}
		}
	} // Fall through case for arbitrary LHS and arbitrary RHS.

	tc.validator.ExpectCanAssignTo(nodeToWarn, rightType, leftType, msg)
}

// Java: TypeCheck.checkPropCreation(Node)
func (tc *TypeCheck) CheckPropCreation(lvalue ast.Target) {
	if getProp, ok := lvalue.AsGetProp(); ok {
		objType := tc.GetJSType(getProp.Object().Node())
		if !objType.IsEmptyType() && !objType.IsUnknownType() {
			propNode := getProp.Node()
			kind := tc.typeRegistry.CanPropertyBeDefinedByID(objType, propNode.GetStringID())
			if kind != rhino.JSTypeRegistry_PropDefinitionKind_KNOWN {

				if objType.IsStruct() {
					if objType.RestrictByNotNullOrUndefined().IsUnionType() {
						tc.Report(propNode, TypeCheck_ILLEGAL_PROPERTY_CREATION_ON_UNION_TYPE)
					} else {
						tc.Report(propNode, TypeCheck_ILLEGAL_PROPERTY_CREATION)
					}
				} else {
					// null checks are reported elsewhere
					if !objType.IsNoType() &&
						!objType.IsUnknownType() &&
						objType.IsSubtypeOf(tc.GetNativeType(rhino.JSTypeNative_NULL_VOID)) {
						return
					}

					tc.ReportMissingProperty(getProp.Object().Node(), objType, propNode, kind, true)
				}
			}
		}
	}
}

// VisitObjectPattern validates all keys in an object pattern.
//
// Validating the types assigned to any lhs nodes in the pattern is done at
// the ASSIGN/VAR/PARAM_LIST/etc. node.
//
// Java: TypeCheck.visitObjectPattern(Node)
func (tc *TypeCheck) VisitObjectPattern(pattern ast.ObjectPattern) {
	n := pattern.Node()
	patternType := tc.GetJSType(n)
	tc.validator.ExpectObject(n, patternType, "cannot destructure 'null' or 'undefined'")
	for prop := range pattern.Props() {
		// Go-only: the Java creates the DestructuredTarget of every prop and asks it for its
		// computed property or string key, which is the prop itself; here the prop is asked
		// for its kind, and the target is created where it is used.
		if computedProperty, ok := prop.AsComputedProp(); ok {
			tc.validator.ExpectIndexMatch(
				computedProperty.Node(), patternType, tc.GetJSType(computedProperty.Key().Node()))
		} else if stringKey, ok := prop.AsStringKey(); ok {
			if !stringKey.IsQuoted() {
				if patternType.IsDict() {
					tc.Report(stringKey.Node(), TypeValidator_ILLEGAL_PROPERTY_ACCESS, "unquoted", "dict")
				}
				// check for missing properties given `const {a} = obj;` but not `const {'a': a} = obj;`
				target := DestructuredTarget_CreateTarget(tc.typeRegistry, patternType, prop.Node())
				tc.CheckPropertyAccessForDestructuring(
					pattern, patternType, stringKey, tc.GetJSType(target.GetNode()))
			} else if patternType.IsStruct() {
				// check that we are not accessing a struct with a quoted string
				tc.Report(stringKey.Node(), TypeValidator_ILLEGAL_PROPERTY_ACCESS, "quoted", "struct")
			}
		}
	}
	tc.EnsureTyped(n)
}

// VisitObjectOrClassLiteralKey visits an object literal field definition
// `key : value`, or a class member definition `key() { ... }` If the
// `lvalue` is a prototype modification, we change the schema of the object
// type it is referring to.
//
// key is the ASSIGN, STRING_KEY, MEMBER_FUNCTION_DEF, SPREAD,
// COMPUTED_PROPERTY, MEMBER_FIELD_DEF, or COMPUTED_FIELD_DEF node; owner the
// parent node, either OBJECTLIT or CLASS_MEMBERS; ownerType the instance type
// of the enclosing object/class.
//
// Java: TypeCheck.visitObjectOrClassLiteralKey(Node, Node, JSType)
func (tc *TypeCheck) VisitObjectOrClassLiteralKey(key rhino.Node, owner rhino.Node, ownerType rhino.JSType) {
	// Do not validate object lit value types in externs. We don't really care,
	// and it makes it easier to generate externs.
	if owner.IsFromExterns() {
		tc.EnsureTyped(key)
		return
	}

	// Validate computed properties similarly to how we validate GETELEMs.
	if key.IsComputedProp() || key.IsComputedFieldDef() {
		keyType := tc.GetJSType(key.GetFirstChild())
		tc.validator.ExpectIndexMatch(key, ownerType, keyType)
		if keyType.IsKnownSymbolValueType() {
			if owner.IsClass() {
				classConstructorType := owner.GetJSType().AssertFunctionType()
				tc.CheckPropertyInheritanceOnClassMember(
					key, rhino.NewProperty_SymbolKey(keyType.ToMaybeKnownSymbolType()), classConstructorType)
			} else {
				jutil.CheckStateMsg1(owner.IsObjectLit(), "Unexpected owner %s", owner)
				// Check if this property has an expected declared type on the owner.
				objectCastType := rhino.ObjectType_Cast(ownerType.RestrictByNotNullOrUndefined())
				expectedPropertyType :=
					tc.GetPropertyTypeIfDeclared(
						objectCastType, rhino.NewProperty_SymbolKey(keyType.ToMaybeKnownSymbolType()))
				tc.validator.ExpectCanAssignToPropertyOf(
					key,
					expectedPropertyType,
					tc.GetJSType(key.GetSecondChild()),
					owner,
					stringPtrOrNull(keyType.ToMaybeKnownSymbolType().GetDisplayName()))
			}
		}
		return
	}

	if key.IsBlock() {
		return
	}

	if key.IsQuotedStringKey() {
		// NB: this case will never be triggered for member functions, since we store quoted
		// member functions as computed properties. This case does apply to regular string key
		// properties, getters, and setters.
		// See also https://github.com/google/closure-compiler/issues/3071
		if ownerType.IsStruct() {
			if owner.IsClass() {
				tc.Report(key, TypeCheck_ILLEGAL_CLASS_KEY, "struct")
			} else {
				tc.Report(key, TypeCheck_ILLEGAL_OBJLIT_KEY, "struct")
			}
		}
	} else {
		// we have either a non-quoted string or a member function def
		// Neither is allowed for an @dict type except for "constructor" as a special case.
		if ownerType.IsDict() && !NodeUtil_IsEs6ConstructorMemberFunctionDef(key) {
			// Object literals annotated as @dict may only have
			// If you annotate a class with @dict, only the constructor can be a non-computed property.
			if owner.IsClass() {
				tc.Report(key, TypeCheck_ILLEGAL_CLASS_KEY, "dict")
			} else {
				tc.Report(key, TypeCheck_ILLEGAL_OBJLIT_KEY, "dict")
			}
		}
	}

	if key.IsSpread() {
		// Rely on type inference to figure out what the key/types this adds to this object.
		return
	}

	// TODO(johnlenz): Validate get and set function declarations are valid
	// as is the functions can have "extraneous" bits.

	// For getter and setter property definitions the
	// r-value type != the property type.
	rvalue := key.GetFirstChild()
	if rvalue.IsNil() {
		tc.EnsureTyped(key)
		return
	}
	rightType := TypeCheck_GetObjectLitKeyTypeFromValueType(key, tc.GetJSType(rvalue))
	if rightType == nil {
		rightType = tc.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
	}

	// Validate value is assignable to the key type.
	keyType := tc.GetJSType(key)
	propertyName := NodeUtil_GetObjectOrClassLitKeyName(key)
	allowedValueType := keyType
	if allowedValueType.IsEnumElementType() {
		allowedValueType = allowedValueType.ToMaybeEnumElementType().GetPrimitiveType()
	}

	valid :=
		tc.validator.ExpectCanAssignToPropertyOf(
			key, rightType, allowedValueType, owner, propertyName)
	if valid {
		tc.EnsureTypedWithType(key, rightType)
	} else {
		tc.EnsureTyped(key)
	}

	// Validate inheritance for classes and object literals used as prototypes
	propertyKey := rhino.NewProperty_StringKey(propertyName)
	if owner.IsClass() {
		classConstructorType := owner.GetJSType().AssertFunctionType()
		tc.CheckPropertyInheritanceOnClassMember(key, propertyKey, classConstructorType)
	} else if ownerType.ToMaybeObjectType() != nil {
		tc.CheckPropertyInheritanceOnPrototypeLitKey(key, propertyKey, ownerType.ToMaybeObjectType())
	}
}

// TypeCheck_PropertyIsImplicitCast returns true if any type in the chain has
// an implicitCast annotation for the given property.
//
// Java: TypeCheck.propertyIsImplicitCast(ObjectType, String)
func TypeCheck_PropertyIsImplicitCast(type_ rhino.ObjectType, prop string) bool {
	for ; type_ != nil; type_ = type_.GetImplicitPrototype() {
		docInfo := type_.GetOwnPropertyJSDocInfo(prop)
		if docInfo != nil && docInfo.IsImplicitCast() {
			return true
		}
	}
	return false
}

// TypeCheck_GetObjectLitKeyTypeFromValueType returns the type expected when
// using the key.
//
// key is a OBJECTLIT key node.
//
// Java: TypeCheck.getObjectLitKeyTypeFromValueType(Node, JSType)
func TypeCheck_GetObjectLitKeyTypeFromValueType(key rhino.Node, valueType rhino.JSType) rhino.JSType {
	if valueType != nil {
		switch key.GetToken() {
		case rhino.Token_GETTER_DEF:
			// GET must always return a function type.
			if valueType.IsFunctionType() {
				fntype := valueType.ToMaybeFunctionType()
				valueType = fntype.GetReturnType()
			} else {
				return nil
			}
		case rhino.Token_SETTER_DEF:
			if valueType.IsFunctionType() {
				// SET must always return a function type.
				fntype := valueType.ToMaybeFunctionType()
				param := fntype.GetParameters()[0]
				// SET function must always have one parameter.
				valueType = param.GetJSType()
			} else {
				return nil
			}
		default:
		}
	}
	return valueType
}

// VisitName visits a NAME node.
//
// t is the node traversal object that supplies context, such as the scope
// chain to use in name lookups as well as error reporting; n the node being
// visited. Returns whether the node is typeable or not.
//
// Java: TypeCheck.visitName(NodeTraversal, Node, Node)
func (tc *TypeCheck) VisitName(t *NodeTraversal, name ast.Name, parent rhino.Node) bool {
	n := name.Node()
	// Skip empty function expression names. They don't need a type.
	// (Go-only: the rhino API, because Function.Name() reports an anonymous function's
	// empty NAME as no name.)
	isFunctionName := n.GetParent().IsFunction() && n.IsFirstChildOf(parent)
	if isFunctionName && name.Ident() == "" {
		return false
	}

	// the compiler's unknown type reporting allows certain names to be untyped if declared but
	// never used in an expression. The associated AST nodes still need to be typed to prevent
	// future passes from crashing, though.
	reportIfMissingType := true

	// At this stage, we need to determine whether this is a leaf
	// node in an expression (which therefore needs to have a type
	// assigned for it) versus some other decorative node that we
	// can safely ignore.
	parentNodeType := parent.GetToken()
	if isFunctionName ||
		parentNodeType == rhino.Token_CATCH ||
		parentNodeType == rhino.Token_PARAM_LIST ||
		NodeUtil_IsNameDeclaration(parent) {
		reportIfMissingType = false
	}

	// Not need to type first key in for-in or for-of.
	if loop, ok := ast.AsEnhancedFor(parent); ok && loop.Target().Is(n) {
		reportIfMissingType = false
	}

	type_ := n.GetJSType()
	if type_ == nil {
		type_ = tc.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
		// TODO(b/149843534): crash instead of defaulting to '?' when the var is null.
		var_ := t.GetTypedScope().GetVar(name.Ident())
		if var_ != nil {
			type_ = var_.GetType()
			if type_ == nil {
				panic(jutil.NewNullPointerException())
			}
		}
	}
	tc.EnsureTypedWithType(n, type_)
	return reportIfMissingType
}

// CheckForOfTypes visits the loop variable of a FOR_OF and FOR_AWAIT_OF and
// verifies the type being assigned to it.
//
// Java: TypeCheck.checkForOfTypes(NodeTraversal, Node)
func (tc *TypeCheck) CheckForOfTypes(t *NodeTraversal, forOf ast.ForOfLoop) {
	iterable := forOf.Iterable().Node()
	iterableType := tc.GetJSType(iterable)
	var actualType rhino.JSType

	if forOf.Token() == rhino.Token_FOR_AWAIT_OF {
		maybeType :=
			tc.validator.ExpectAutoboxesToIterableOrAsyncIterable(
				iterable,
				iterableType,
				"Can only async iterate over a (non-null) Iterable or AsyncIterable type")

		if maybeType == nil {
			// Not iterable or async iterable, error reported by
			// expectAutoboxesToIterableOrAsyncIterable.
			return
		}

		actualType = maybeType
	} else {
		tc.validator.ExpectAutoboxesToIterable(
			iterable, iterableType, "Can only iterate over a (non-null) Iterable type")

		actualType =
			// Convert primitives to their wrapper type and remove null/undefined
			// If iterable is a union type, autoboxes each member of the union.
			iterableType.
				Autobox().
				GetTemplateTypeMap().
				GetResolvedTemplateType(tc.typeRegistry.GetIterableValueTemplate())
	}

	lhs := typeCheck_forOfTarget(forOf.Target())
	tc.CheckCanAssignToWithScope(
		t,
		forOf.Node(),
		lhs,
		actualType,
		lhs.Node().GetJSDocInfo(),
		"declared type of for-of loop variable does not match inferred type")
}

// typeCheck_forOfTarget returns what a for-of loop assigns each value to:
// its target, or what its declaration declares.
//
// Go-only: the Java steps down from the declaration by position.
func typeCheck_forOfTarget(target ast.ForTarget) ast.Target {
	decl, ok := target.AsNameDeclaration()
	if !ok {
		return ast.MustTarget(target.Node())
	}
	// e.g. get "x" given the VAR in "for (var x of arr) {"
	declarator, ok := decl.OnlyDeclarator()
	jutil.CheckStateMsg1(ok, "%s", decl.Node())
	if destructuring, ok := declarator.AsDestructuringLhs(); ok {
		// e.g. get `[x, y]` given the VAR in `for (var [x, y] of arr) {`
		return destructuring.Pattern().AsTarget()
	}
	return ast.MustName(declarator.Node()).AsTarget()
}

// CheckNotReadonlyPropertyAssignment: given the LHS of a property or element
// assignment, checks that the type that we're assigning into is not readonly
// (ReadonlyArray, in particular).
//
// This is basically unsound since it checks just "ReadonlyArray" so casting
// up to "Iterable" or down to "Array" will defeat it; but it should catch
// basic errors.
//
// Java: TypeCheck.checkNotReadonlyPropertyAssignment(Node)
func (tc *TypeCheck) CheckNotReadonlyPropertyAssignment(lhs ast.Target) {
	// We only care about element or property assignments.
	get, ok := lhs.AsNormalGet()
	if !ok {
		return
	}

	// If we do not have type information, drop out.
	lhsType := get.Object().Node().GetJSType()
	if lhsType == nil {
		return
	}

	// We could be a reference to "ReadonlyArray" or a templatized wrapper.
	roArray := tc.GetNativeType(rhino.JSTypeNative_READONLY_ARRAY_TYPE)
	if !lhsType.IsUnionType() {
		if roArray.Equals(lhsType) || roArray.Equals(TypeCheck_MaybeReferencedType(lhsType)) {
			tc.compiler.Report(JSError_MakeWithNode(get.Node(), TypeCheck_PROPERTY_ASSIGNMENT_TO_READONLY_VALUE, lhsType.String()))
		}
		return
	}
	alternates := TypeCheck_FlattenUnion(lhsType)
	for i := 0; i < len(alternates); i++ {
		type_ := alternates[i]
		if roArray.Equals(type_) || roArray.Equals(TypeCheck_MaybeReferencedType(type_)) {
			tc.compiler.Report(JSError_MakeWithNode(get.Node(), TypeCheck_PROPERTY_ASSIGNMENT_TO_READONLY_VALUE, type_.String()))
			break
		}
	}
}
