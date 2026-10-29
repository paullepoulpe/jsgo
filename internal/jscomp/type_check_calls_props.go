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
	"math"
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// VisitGetProp visits a GETPROP node.
//
// t is the node traversal object that supplies context, such as the scope
// chain to use in name lookups as well as error reporting; n the node being
// visited.
//
// Java: TypeCheck.visitGetProp(NodeTraversal, Node)
func (tc *TypeCheck) VisitGetProp(t *NodeTraversal, getProp ast.GetProp) {
	// obj.prop or obj.method()
	// Lots of types can appear on the left, a call to a void function can
	// never be on the left. getPropertyType will decide what is acceptable
	// and what isn't.
	n := getProp.Node()
	objNode := getProp.Object().Node()
	childType := tc.GetJSType(objNode)

	if childType.IsDict() {
		tc.Report(n, TypeValidator_ILLEGAL_PROPERTY_ACCESS, "'.'", "dict")
	} else if tc.validator.ExpectNotNullOrUndefined(
		t, n, childType, "No properties on this expression", tc.GetNativeType(rhino.JSTypeNative_OBJECT_TYPE)) {
		tc.CheckPropertyAccessForGetProp(getProp.AsNormalOrOptChainGetProp())
	}
	tc.EnsureTyped(n)
}

// VisitOptChainGetProp visits a OPTCHAIN_GETPROP node.
//
// Java: TypeCheck.visitOptChainGetProp(Node)
func (tc *TypeCheck) VisitOptChainGetProp(optChainGetProp ast.OptChainGetProp) {
	// obj?.prop
	// Unlike GETPROP, a call to a void function can also be on the lhs for an OPTCHAIN_GETPROP.
	objNode := optChainGetProp.Object().Node()
	childType := tc.GetJSType(objNode)

	if childType.IsDict() {
		// Go-only: the Java reports on `property`, optChainGetProp.getLastChild(), which is the
		// object: the property name is the node's string, not a child.
		tc.Report(objNode, TypeValidator_ILLEGAL_PROPERTY_ACCESS, "'?.'", "dict")
	} else if !childType.IsUnknownType() {
		tc.CheckPropertyAccessForGetProp(optChainGetProp.AsNormalOrOptChainGetProp())
	}
	tc.EnsureTyped(optChainGetProp.Node())
}

// Java: TypeCheck.checkPropertyAccessForGetProp(Node)
func (tc *TypeCheck) CheckPropertyAccessForGetProp(getProp ast.NormalOrOptChainGetProp) {
	objNode := getProp.Object().Node()
	objType := tc.GetJSType(objNode)
	propType := tc.GetJSType(getProp.Node())

	tc.CheckAbstractPropertyAccess(getProp)
	tc.CheckPropertyAccess(objType, getProp.Node(), propType, objNode)
}

// Java: TypeCheck.checkPropertyAccessForDestructuring(Node, JSType, Node, JSType)
func (tc *TypeCheck) CheckPropertyAccessForDestructuring(
	pattern ast.ObjectPattern, objectType rhino.JSType, stringKey ast.StringKey, inferredPropType rhino.JSType) {

	// Get the object node being destructured when it exists. These cases have an actual node `obj`:
	//   const {a} = obj;
	//   ({a} = obj);
	// while these do not:
	//    for (const {a} of arrayOfObjects) {
	//    const {{a}} = obj;
	var objNode rhino.Node
	patternParent := pattern.Node().GetParent()
	if assign, ok := ast.AsAssign(patternParent); ok {
		objNode = assign.Value().Node()
	} else if declaration, ok := ast.AsDestructuringLhs(patternParent); ok {
		if value, ok := declaration.Init(); ok {
			objNode = value.Node()
		}
	}
	tc.CheckPropertyAccess(objectType, stringKey.Node(), inferredPropType, objNode)
}

// CheckAbstractPropertyAccess warns if @abstract methods are dereferenced,
// with some false negatives.
//
// This method only handles some cases that we are certain are incorrect.
// e.g. we are lenient about union types, and we don't track an abstract
// methods once they are reassigned to new variables.
//
// This method's logic is complicated to avoid spurious warnings. Sometimes a
// reference to something typed @abstract is okay. For example, don't warn on
// `this.foo();` in
// `/** @abstract * / class Base { /** @abstract * / foo() {} bar() { this.foo(); } } }`
//
// The `this` object with which `Base.prototype.bar` is called must be a
// concrete subclass of Base. To avoid false positives, we warn only in the
// following cases:
//   - (a) the function is accessed off `super`
//   - (b) the function is accessed off a .prototype
//   - (c) the function is transpiled from a goog.base superclass reference
//
// Java: TypeCheck.checkAbstractPropertyAccess(Node)
func (tc *TypeCheck) CheckAbstractPropertyAccess(method ast.NormalOrOptChainGetProp) {

	if NodeUtil_IsLhsOfAssign(method.Node()) {
		// Allow declaring abstract methods. (This assumes they are never re-assigned)
		//   /** @abstract */ Foo.prototype.bar = function() {}
		return
	}

	methodType := tc.GetJSType(method.Node()).ToMaybeFunctionType()
	if methodType == nil || !methodType.IsAbstract() || methodType.IsConstructor() {
		// Ignore non-abstract methods and @abstract constructors. An @abstract constructor is still
		// callable.
		return
	}

	object := method.Object()
	if object.Node().IsSuper() {
		// case (a)
		// `super.foo()` definitely refers to `Superclass.prototype.foo`, not an override.
		// At parse time, `Subclass.prototype` becomes a lower bound for what `super` evaluates to,
		// even if the `this` object changes. So `super` will never resolve to a concrete subclass.
		tc.Report(method.Node(), TypeCheck_ABSTRACT_SUPER_METHOD_NOT_USABLE, stringPtrOrNull(methodType.GetDisplayName()))
	} else if objectGetProp, ok := object.AsGetProp(); ok {
		if objectGetProp.Node().GetStringID() == typeCheckPrototypeID || // case (b), e.g. `Foo.prototype.bar`
			tc.compiler.GetCodingConvention().IsSuperClassReference(objectGetProp.Property()) { // case (c)
			tc.Report(method.Node(), TypeCheck_ABSTRACT_SUPER_METHOD_NOT_USABLE, stringPtrOrNull(methodType.GetDisplayName()))
		}
	}
}

// CheckPropertyAccess emits a warning if we can prove that a property cannot
// possibly be defined on an object. Note the difference between JS and a
// strictly statically typed language: we're checking if the property *cannot
// be defined*, whereas a java compiler would check if the property *can be
// undefined.
//
// This method handles property access in both GETPROPs and object
// destructuring. Consequentially some of its arguments are optional - the
// actual object node and the getprop - while others are required.
//
// objNode is the actual node representing the object we're accessing.
// optional because destructuring accesses MAY not have an actual object node.
//
// Java: TypeCheck.checkPropertyAccess(JSType, Node, JSType, Node)
func (tc *TypeCheck) CheckPropertyAccess(
	childType rhino.JSType, propNode rhino.Node, propType rhino.JSType, objNode rhino.Node) {
	isGetprop := NodeUtil_IsNormalOrOptChainGetProp(propNode)

	unknownType := tc.typeRegistry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
	if propType == unknownType || propType.Equals(unknownType) {
		childType = childType.Autobox()
		objectType := rhino.ObjectType_Cast(childType)
		if objectType != nil {
			// We special-case object types so that checks on enums can be
			// much stricter, and so that we can use hasProperty (which is much
			// faster in most cases).
			if !rhino.HasPropertyByID(objectType, propNode.GetStringID()) ||
				objectType == unknownType ||
				objectType.Equals(unknownType) {
				if _, isEnum := objectType.(*rhino.EnumType); isEnum {
					tc.Report(propNode, TypeCheck_INEXISTENT_ENUM_ELEMENT, propNode.GetString())
				} else {
					tc.CheckPropertyAccessHelper(objectType, propNode, objNode, false)
				}
			}
		} else {
			tc.CheckPropertyAccessHelper(childType, propNode, objNode, false)
		}
	} else if childType.IsUnionType() && isGetprop && !tc.IsLValueGetProp(propNode) {
		// NOTE: strict property assignment checks are done on assignment.
		tc.CheckPropertyAccessHelper(childType, propNode, objNode, true)
	}
}

// Java: TypeCheck.isLValueGetProp(Node)
func (tc *TypeCheck) IsLValueGetProp(n rhino.Node) bool {
	parent := n.GetParent()
	// TODO(b/77597706): this won't work for destructured lvalues
	if incDec, ok := ast.AsIncDec(parent); ok {
		return incDec.Target().Is(n)
	}
	assign, ok := ast.AsAssignment(parent)
	return ok && assign.Target().Is(n)
}

// CheckPropertyAccessHelper: strictCheck is whether this is a check that is
// only performed when "strict missing properties" checks are enabled.
//
// Java: TypeCheck.checkPropertyAccessHelper(JSType, Node, Node, boolean)
func (tc *TypeCheck) CheckPropertyAccessHelper(
	objectType rhino.JSType, propNode rhino.Node, objNode rhino.Node, strictCheck bool) {
	isGetprop := NodeUtil_IsNormalOrOptChainGetProp(propNode)

	if !tc.reportMissingProperties ||
		objectType.IsEmptyType() ||
		(isGetprop && tc.AllowStrictPropertyAccessOnNode(propNode)) {
		return
	}

	kind := tc.typeRegistry.CanPropertyBeDefinedByID(objectType, propNode.GetStringID())
	if kind == rhino.JSTypeRegistry_PropDefinitionKind_KNOWN {
		return
	}
	// If the property definition is known, but only loosely associated,
	// only report a "strict error" which can be optional as code is migrated.
	isLooselyAssociated :=
		kind == rhino.JSTypeRegistry_PropDefinitionKind_LOOSE || kind == rhino.JSTypeRegistry_PropDefinitionKind_LOOSE_UNION
	isUnknownType := objectType.IsUnknownType()
	if isLooselyAssociated && isUnknownType {
		// We still don't want to report this.
		return
	}

	isStruct := objectType.IsStruct()
	loosePropertyDeclaration := !isStruct && isGetprop && tc.IsQNameAssignmentTarget(propNode)
	// always false for destructuring
	maybePropExistenceCheck :=
		!isStruct && isGetprop && tc.AllowLoosePropertyAccessOnNode(propNode)
	// Traditionally, we would not report a warning for "loose" properties, but we want to be
	// able to be more strict, so introduce an optional warning.
	strictReport :=
		strictCheck || isLooselyAssociated || loosePropertyDeclaration || maybePropExistenceCheck

	tc.ReportMissingProperty(objNode, objectType, propNode, kind, strictReport)
}

// Java: TypeCheck.reportMissingProperty(Node, JSType, Node, PropDefinitionKind, boolean)
func (tc *TypeCheck) ReportMissingProperty(
	objNode rhino.Node,
	objectType rhino.JSType,
	propNode rhino.Node,
	kind rhino.JSTypeRegistry_PropDefinitionKind,
	strictReport bool) {
	propName := propNode.GetString()

	isObjectType := objectType.Equals(tc.GetNativeType(rhino.JSTypeNative_OBJECT_TYPE))
	lowConfidence := objectType.IsUnknownType() || objectType.IsAllType() || isObjectType

	isKnownToUnionMember := kind == rhino.JSTypeRegistry_PropDefinitionKind_LOOSE_UNION

	var pair *TypeCheck_SuggestionPair
	if !lowConfidence && !isKnownToUnionMember {
		pair = TypeCheck_GetClosestPropertySuggestion(objectType, propName, (jutil.StringLength(propName)-1)/4)
	}
	var objectTypeName string
	if !objNode.IsNil() {
		objectTypeName = tc.typeRegistry.GetReadableTypeName(objNode)
	} else {
		objectTypeName = objectType.String()
	}
	if pair != nil {
		var reportType *DiagnosticType
		if strictReport {
			reportType = TypeCheck_STRICT_INEXISTENT_PROPERTY_WITH_SUGGESTION
		} else {
			reportType = TypeCheck_INEXISTENT_PROPERTY_WITH_SUGGESTION
		}
		tc.Report(
			propNode,
			reportType,
			propName,
			objectTypeName,
			pair.suggestion)
	} else {
		var reportType *DiagnosticType
		if strictReport {
			if isKnownToUnionMember {
				reportType = TypeCheck_STRICT_INEXISTENT_UNION_PROPERTY
			} else {
				reportType = TypeCheck_STRICT_INEXISTENT_PROPERTY
			}
		} else if lowConfidence {
			reportType = TypeCheck_POSSIBLE_INEXISTENT_PROPERTY
		} else {
			reportType = TypeCheck_INEXISTENT_PROPERTY
		}
		tc.Report(
			propNode,
			reportType,
			propName,
			objectTypeName)
	}
}

// Java: TypeCheck.allowStrictPropertyAccessOnNode(Node)
func (tc *TypeCheck) AllowStrictPropertyAccessOnNode(n rhino.Node) bool {
	return n.GetParent().IsTypeOf()
}

// AllowLoosePropertyAccessOnNode checks whether a property access is either
// (1) a stub declaration, or (2) a property presence or absence test.
//
// Presence and absence are both allowed here because both are valid for
// conditional use of a property (e.g. `if (x.y != null) use(x.y);` but also
// the opposite, `if (x.y == null) return; use(x.y);`).
//
// Java: TypeCheck.allowLoosePropertyAccessOnNode(Node)
func (tc *TypeCheck) AllowLoosePropertyAccessOnNode(n rhino.Node) bool {
	parent := n.GetParent()
	return NodeUtil_IsPropertyTest(tc.compiler, n) ||
		NodeUtil_IsPropertyAbsenceTest(n) ||
		// Stub property declaration
		(n.IsQualifiedName() && parent.IsExprResult())
}

// Java: TypeCheck.isQNameAssignmentTarget(Node)
func (tc *TypeCheck) IsQNameAssignmentTarget(n rhino.Node) bool {
	assign, ok := ast.AsAssign(n.GetParent())
	return n.IsQualifiedName() && ok && assign.Target().Is(n)
}

// Java: TypeCheck.getClosestPropertySuggestion(JSType, String, int)
func TypeCheck_GetClosestPropertySuggestion(
	objectType rhino.JSType, propName string, maxDistance int) *TypeCheck_SuggestionPair {
	return nil
}

// VisitGetElem visits a GETELEM node.
//
// Java: TypeCheck.visitGetElem(Node)
func (tc *TypeCheck) VisitGetElem(getElem ast.NormalOrOptChainGetElem) {
	n := getElem.Node()
	tc.validator.ExpectIndexMatch(n, tc.GetJSType(getElem.Object().Node()), tc.GetJSType(getElem.Index().Node()))
	tc.EnsureTyped(n)
}

// VisitOptChainGetElem visits a OPTCHAIN_GETELEM node.
//
// Java: TypeCheck.visitOptChainGetElem(Node)
func (tc *TypeCheck) VisitOptChainGetElem(optChainGetElem ast.OptChainGetElem) {
	obj := optChainGetElem.Object().Node()
	objType := tc.GetJSType(obj)
	if objType.IsNullType() || objType.IsVoidType() {
		// no error reported when conditionally checking a null or void lhs object.
		// e.g. `a?.[b]` should not report if `a` is null`
		tc.EnsureTypedWithNative(optChainGetElem.Node(), rhino.JSTypeNative_VOID_TYPE)
	} else {
		tc.VisitGetElem(optChainGetElem.AsNormalOrOptChainGetElem())
	}
}

// VisitVar visits a VAR node.
//
// t is the node traversal object that supplies context, such as the scope
// chain to use in name lookups as well as error reporting; n the node being
// visited.
//
// Java: TypeCheck.visitVar(NodeTraversal, Node)
func (tc *TypeCheck) VisitVar(t *NodeTraversal, decl ast.NameDeclaration) {
	n := decl.Node()
	// Handle var declarations in for-of loops separately from regular var declarations.
	if _, inLoop := ast.AsEnhancedFor(n.GetParent()); inLoop {
		return
	}

	// TODO(nicksantos): Fix this so that the doc info always shows up
	// on the NAME node. We probably want to wait for the parser
	// merge to fix this.
	var varInfo *rhino.JSDocInfo
	if _, ok := decl.OnlyDeclarator(); ok {
		varInfo = n.GetJSDocInfo()
	}
	for declarator := range decl.Declarators() {
		if name, ok := declarator.AsName(); ok {
			if value, ok := name.Init(); ok {
				valueType := tc.GetJSType(value.Node())
				info := name.Node().GetJSDocInfo()
				if info == nil {
					info = varInfo
				}

				tc.CheckEnumAlias(t, info, valueType, value.Node())
				tc.CheckCanAssignToWithScope(t, value.Node(), name.AsTarget(), valueType, info, "initializing variable")
			}
		} else {
			destructuring := ast.MustDestructuringLhs(declarator.Node())
			// Outside a for-in or for-of loop, a destructuring declaration has a value.
			value, ok := destructuring.Init()
			jutil.CheckStateMsg1(ok, "%s", destructuring.Node())
			valueType := tc.GetJSType(value.Node())
			tc.CheckCanAssignToWithScope(
				t, destructuring.Node(), destructuring.Pattern().AsTarget(), valueType /* info= */, nil, "initializing variable")
		}
	}
}

// VisitNew visits a NEW node.
//
// Java: TypeCheck.visitNew(Node)
func (tc *TypeCheck) VisitNew(newNode ast.New) {
	n := newNode.Node()
	constructor := newNode.Callee().Node()
	type_ := tc.GetJSType(constructor).RestrictByNotNullOrUndefined()
	if !tc.CouldBeAConstructor(type_) ||
		type_.Equals(tc.typeRegistry.GetNativeType(rhino.JSTypeNative_SYMBOL_OBJECT_FUNCTION_TYPE)) ||
		type_.Equals(tc.typeRegistry.GetNativeType(rhino.JSTypeNative_BIGINT_OBJECT_FUNCTION_TYPE)) {
		tc.Report(n, TypeCheck_NOT_A_CONSTRUCTOR, type_.String())
		tc.EnsureTyped(n)
		return
	}

	fnType := type_.ToMaybeFunctionType()
	if fnType != nil && fnType.HasInstanceType() {
		if fnType.IsAbstract() {
			tc.Report(n, TypeCheck_INSTANTIATE_ABSTRACT_CLASS)
		}

		tc.VisitArgumentList(newNode.AsCallOrNew(), fnType)

		objType := fnType.GetInstanceType()
		if objType != nil {
			tc.EnsureTypedWithType(n, objType)
		} else {
			tc.EnsureTypedWithType(n, tc.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE))
		}
	} else {
		tc.EnsureTyped(n)
	}
}

// Java: TypeCheck.couldBeAConstructor(JSType)
func (tc *TypeCheck) CouldBeAConstructor(type_ rhino.JSType) bool {
	return type_.IsConstructor() || type_.IsEmptyType() || type_.IsUnknownType()
}

// CheckInterfaceConflictProperties checks whether there's any property
// conflict for a particular super interface.
//
// n is the node being visited; functionName the function name being checked;
// properties the property names in the super interfaces that have been
// visited; currentProperties the property names in the super interface that
// have been visited; interfaceType the super interface that is being visited.
//
// Java: TypeCheck.checkInterfaceConflictProperties(Node, String, Map<Property.Key, ObjectType>, Map<Property.Key, ObjectType>, ObjectType)
func (tc *TypeCheck) CheckInterfaceConflictProperties(
	n rhino.Node,
	functionName string,
	properties *jutil.ObjectLinkedHashMap[rhino.Property_Key, rhino.ObjectType],
	currentProperties *jutil.ObjectLinkedHashMap[rhino.Property_Key, rhino.ObjectType],
	interfaceType rhino.ObjectType) {
	implicitProto := interfaceType.GetImplicitPrototype()
	var currentPropertyNames []rhino.Property_Key
	if implicitProto == nil {
		// This can be the case if interfaceType is proxy to a non-existent
		// object (which is a bad type annotation, but shouldn't crash).
		currentPropertyNames = nil
	} else {
		currentPropertyNames = implicitProto.GetOwnPropertyKeys()
	}
	for _, name := range currentPropertyNames {
		oType := properties.Get(name)
		currentProperties.Put(name, interfaceType)
		if oType != nil {
			thisPropType := interfaceType.GetPropertyTypeByKey(name)
			oPropType := oType.GetPropertyTypeByKey(name)
			if thisPropType.IsSubtypeOfWithMode(oPropType, tc.subtypingMode) ||
				oPropType.IsSubtypeOfWithMode(thisPropType, tc.subtypingMode) ||
				(thisPropType.IsFunctionType() &&
					oPropType.IsFunctionType() &&
					thisPropType.
						ToMaybeFunctionType().
						HasEqualCallType(oPropType.ToMaybeFunctionType())) {
				continue
			}
			tc.compiler.Report(
				JSError_MakeWithNode(
					n,
					TypeCheck_INCOMPATIBLE_EXTENDED_PROPERTY_TYPE,
					functionName,
					name.HumanReadableName(),
					oType.String(),
					interfaceType.String()))
		}
	}
	for _, iType := range interfaceType.GetCtorExtendedInterfaces() {
		tc.CheckInterfaceConflictProperties(n, functionName, properties, currentProperties, iType)
	}
}

// VisitFunction visits a FUNCTION node.
//
// Java: TypeCheck.visitFunction(Node)
func (tc *TypeCheck) VisitFunction(n rhino.Node) {
	if NodeUtil_IsEs6Constructor(n) {
		return // These will be checked via the CLASS node.
	}

	functionType := rhino.JSType_ToMaybeFunctionType(n.GetJSType())
	if functionType.IsConstructor() {
		tc.CheckConstructor(n, functionType)
	} else if functionType.IsInterface() {
		tc.CheckInterface(n, functionType)
	} else if n.IsAsyncGeneratorFunction() {
		// An async generator function must return a AsyncGenerator or supertype of AsyncGenerator
		returnType := functionType.GetReturnType()
		tc.validator.ExpectAsyncGeneratorSupertype(
			n, returnType, "An async generator function must return a (supertype of) AsyncGenerator")
	} else if n.IsGeneratorFunction() {
		// A generator function must return a Generator or supertype of Generator
		returnType := functionType.GetReturnType()
		tc.validator.ExpectGeneratorSupertype(
			n, returnType, "A generator function must return a (supertype of) Generator")

	} else if n.IsAsyncFunction() {
		// An async function must return a Promise or supertype of Promise
		returnType := functionType.GetReturnType()
		tc.validator.ExpectValidAsyncReturnType(n, returnType)
	}
}

// VisitClass visits a CLASS node.
//
// Java: TypeCheck.visitClass(Node)
func (tc *TypeCheck) VisitClass(class ast.Class) {
	n := class.Node()
	functionType := rhino.JSType_ToMaybeFunctionType(n.GetJSType())
	if extendsClause, ok := class.Extends(); ok {
		// Ensure that the `extends` clause is actually a constructor or interface.  If it is, but
		// it's the wrong one then checkConstructor or checkInterface will warn.
		superType := extendsClause.Node().GetJSType()
		if superType.IsConstructor() || superType.IsInterface() {
			tc.validator.ExpectExtends(n, functionType, superType.ToMaybeFunctionType())
		} else if !superType.IsUnknownType() {
			// Only give this error for supertypes *known* to be wrong - unresolved types are OK here.
			kind := "interface"
			if functionType.IsConstructor() {
				kind = "constructor"
			}
			tc.compiler.Report(
				JSError_MakeWithNode(
					n,
					TypeCheck_CONFLICTING_EXTENDED_TYPE,
					kind,
					tc.GetBestFunctionName(n)))
		}
	}
	if functionType.IsConstructor() {
		tc.CheckConstructor(n, functionType)
	} else if functionType.IsInterface() {
		tc.CheckInterface(n, functionType)
	} else {
		panic(jutil.NewIllegalStateExceptionMsg(
			"CLASS node's type must be either constructor or interface: " + jsTypeStringOrNull(functionType)))
	}
}

// CheckConstructor checks a constructor, which may be either an ES5-style
// FUNCTION node, or a CLASS node.
//
// Java: TypeCheck.checkConstructor(Node, FunctionType)
func (tc *TypeCheck) CheckConstructor(n rhino.Node, functionType rhino.FunctionType) {
	baseConstructor := functionType.GetSuperClassConstructor()
	if baseConstructor != nil &&
		baseConstructor.IsInterface() &&
		!baseConstructor.Equals(tc.GetNativeType(rhino.JSTypeNative_OBJECT_FUNCTION_TYPE)) {
		// Warn if a class extends an interface.
		tc.compiler.Report(
			JSError_MakeWithNode(n, TypeCheck_CONFLICTING_EXTENDED_TYPE, "constructor", tc.GetBestFunctionName(n)))
	} else {
		if n.IsFunction() &&
			baseConstructor != nil &&
			!baseConstructor.GetSource().IsNil() &&
			baseConstructor.GetSource().IsClass() &&
			!functionType.GetSource().IsClass() {
			// Warn if an ES5 class extends an ES6 class.
			tc.compiler.Report(
				JSError_MakeWithNode(
					n,
					TypeCheck_ES5_CLASS_EXTENDING_ES6_CLASS,
					stringPtrOrNull(functionType.GetDisplayName()),
					stringPtrOrNull(baseConstructor.GetDisplayName())))
		}
		if baseConstructor != nil {
			tc.CheckStructDictSubtyping(n, functionType, baseConstructor)
		}

		// Warn if any @implemented types are not interfaces or if there are any duplicates
		implementedInterfaces := functionType.GetOwnImplementedInterfaces()
		if len(implementedInterfaces) > 0 {
			var seenBuf [8]rhino.JSType
			seen := seenBuf[:0]

			for _, baseInterfaceObjectType := range implementedInterfaces {
				var baseInterface rhino.JSType = baseInterfaceObjectType
				badImplementedType := false
				baseInterfaceObj := rhino.ObjectType_Cast(baseInterface)
				if baseInterfaceObj != nil {
					interfaceConstructor := baseInterfaceObj.GetConstructor()
					if interfaceConstructor != nil && !interfaceConstructor.IsInterface() {
						badImplementedType = true
					}
				} else {
					badImplementedType = true
				}
				if badImplementedType {
					tc.Report(n, TypeCheck_BAD_IMPLEMENTED_TYPE, tc.GetBestFunctionName(n))
				}
				if len(implementedInterfaces) > 1 {
					// Disallow implementing both `Foo<string>` and `Foo<number>`, for example.
					baseInterface = tc.NormalizeTemplatizedType(baseInterface)
					duplicate := false
					for _, prev := range seen {
						if prev == baseInterface || prev.Equals(baseInterface) {
							duplicate = true
							break
						}
					}
					if duplicate {
						tc.Report(n, TypeCheck_SAME_INTERFACE_MULTIPLE_IMPLEMENTS, baseInterface.String())
					} else {
						seen = append(seen, baseInterface)
					}
				}
			}
		}
		// check properties
		tc.validator.ExpectAllInterfaceProperties(n, functionType)
		if !functionType.IsAbstract() {
			tc.validator.ExpectAbstractMethodsImplemented(n, functionType)
		}
	}
}

// NormalizeTemplatizedType normalizes `Foo<string>` to just `Foo` and is the
// identify function for other types.
//
// Java: TypeCheck.normalizeTemplatizedType(JSType)
func (tc *TypeCheck) NormalizeTemplatizedType(maybeTemplatizedType rhino.JSType) rhino.JSType {
	if !maybeTemplatizedType.IsTemplatizedType() {
		return maybeTemplatizedType
	}
	return maybeTemplatizedType.ToMaybeTemplatizedType().GetRawType()
}

// Java: TypeCheck.checkStructDictSubtyping(Node, FunctionType, FunctionType)
func (tc *TypeCheck) CheckStructDictSubtyping(n rhino.Node, subtype rhino.FunctionType, supertype rhino.FunctionType) {
	if subtype.MakesDicts() && supertype.MakesStructs() {
		tc.compiler.Report(
			JSError_MakeWithNode(
				n, TypeCheck_DICT_EXTEND_STRUCT_TYPE, stringPtrOrNull(subtype.GetDisplayName()), stringPtrOrNull(supertype.GetDisplayName())))
	} else if subtype.MakesStructs() && supertype.MakesDicts() {
		tc.compiler.Report(
			JSError_MakeWithNode(
				n, TypeCheck_STRUCT_EXTEND_DICT_TYPE, stringPtrOrNull(subtype.GetDisplayName()), stringPtrOrNull(supertype.GetDisplayName())))
	}
}

// CheckInterface checks an interface, which may be either an ES5-style
// FUNCTION node, or a CLASS node.
//
// Java: TypeCheck.checkInterface(Node, FunctionType)
func (tc *TypeCheck) CheckInterface(n rhino.Node, functionType rhino.FunctionType) {
	// Interface must extend only interfaces
	for _, extInterface := range functionType.GetExtendedInterfaces() {
		if extInterface.GetConstructor() != nil && !extInterface.GetConstructor().IsInterface() {
			tc.compiler.Report(
				JSError_MakeWithNode(n, TypeCheck_CONFLICTING_EXTENDED_TYPE, "interface", tc.GetBestFunctionName(n)))
		}
	}

	// Check whether the extended interfaces have any conflicts
	if functionType.GetExtendedInterfacesCount() > 1 {
		// Only check when extending more than one interfaces
		properties := jutil.NewObjectLinkedHashMap[rhino.Property_Key, rhino.ObjectType]()
		currentProperties := jutil.NewObjectLinkedHashMap[rhino.Property_Key, rhino.ObjectType]()
		for _, interfaceType := range functionType.GetExtendedInterfaces() {
			currentProperties.Clear()
			tc.CheckInterfaceConflictProperties(
				n, tc.GetBestFunctionName(n), properties, currentProperties, interfaceType)
			for k, v := range currentProperties.All() {
				properties.Put(k, v)
			}
		}
	}

	loopPath := functionType.CheckExtendsLoop()
	if loopPath != nil {
		names := make([]string, len(loopPath))
		for i, fn := range loopPath {
			names[i] = stringPtrOrNull(fn.GetDisplayName())
		}
		strPath := strings.Join(names, " -> ")
		tc.compiler.Report(
			JSError_MakeWithNode(n, TypeCheck_INTERFACE_EXTENDS_LOOP, stringPtrOrNull(loopPath[0].GetDisplayName()), strPath))
	}

	tc.validator.ExpectAllInterfaceProperties(n, functionType)
}

// Java: TypeCheck.getBestFunctionName(Node)
func (tc *TypeCheck) GetBestFunctionName(n rhino.Node) string {
	jutil.CheckState(n.IsClass() || n.IsFunction())
	name := NodeUtil_GetBestLValueName(NodeUtil_GetBestLValue(n))
	if name != nil {
		return *name
	}
	return "<anonymous@" + stringPtrOrNull(n.GetSourceFileName()) + ":" + strconv.Itoa(n.GetLineno()) + ">"
}

var TypeCheck_GOOG_INHERITS = rhino.QualifiedName_Of("goog.inherits")

// CheckCallConventions validates class-defining calls. Because JS has no
// 'native' syntax for defining classes, we need to do this manually.
//
// Java: TypeCheck.checkCallConventions(NodeTraversal, Node)
func (tc *TypeCheck) CheckCallConventions(t *NodeTraversal, call ast.NormalOrOptChainCall) {
	n := call.Node()
	relationship := tc.compiler.GetCodingConvention().GetClassesDefinedByCall(n)
	scope := t.GetTypedScope()
	if relationship != nil {
		superClass := rhino.StaticTypedScope_LookupQualifiedNameString(scope, relationship.superclassName)
		superClassInstanceType := TypeValidator_GetInstanceOfCtor(superClass)
		subClass := rhino.StaticTypedScope_LookupQualifiedNameString(scope, relationship.subclassName)
		subClassInstance := TypeValidator_GetInstanceOfCtor(subClass)
		if relationship.type_ == CodingConvention_SubclassType_INHERITS &&
			superClassInstanceType != nil &&
			!superClassInstanceType.IsEmptyType() &&
			subClassInstance != nil &&
			!subClassInstance.IsEmptyType() {
			if callee := call.Callee().Node(); callee.IsQualifiedName() &&
				TypeCheck_GOOG_INHERITS.Matches(callee) &&
				subClass.ToMaybeFunctionType() != nil &&
				!subClass.ToMaybeFunctionType().GetSource().IsNil() &&
				subClass.ToMaybeFunctionType().GetSource().IsClass() {
				tc.compiler.Report(JSError_MakeWithNode(n, TypeCheck_ES6_CLASS_EXTENDING_CLASS_WITH_GOOG_INHERITS))
			}
			tc.validator.ExpectSuperType(n, superClassInstanceType, subClassInstance)
		}
	}
}

// VisitCall visits a CALL node.
//
// t is the node traversal object that supplies context, such as the scope
// chain to use in name lookups as well as error reporting; n the node being
// visited.
//
// Java: TypeCheck.visitCall(NodeTraversal, Node)
func (tc *TypeCheck) VisitCall(t *NodeTraversal, call ast.NormalOrOptChainCall) {
	n := call.Node()
	tc.CheckCallConventions(t, call)

	callee := call.Callee()
	childType := tc.GetJSType(callee.Node()).RestrictByNotNullOrUndefined()

	if !childType.CanBeCalled() {
		tc.Report(n, TypeCheck_NOT_CALLABLE, childType.String())
		tc.EnsureTyped(n)
		return
	}

	// A couple of types can be called as if they were functions.
	// If it is a function type, then validate parameters.
	if childType.IsFunctionType() {
		functionType := childType.ToMaybeFunctionType()

		// Non-native constructors should not be called directly
		// unless they specify a return type
		if functionType.IsConstructor() &&
			!functionType.IsNativeObjectType() &&
			(functionType.GetReturnType().IsUnknownType() ||
				functionType.GetReturnType().IsVoidType()) &&
			!callee.Node().IsSuper() {
			tc.Report(n, TypeCheck_CONSTRUCTOR_NOT_CALLABLE, childType.String())
		}

		// Functions with explicit 'this' types must be called in a GETPROP or GETELEM.
		_, calleeIsGet := callee.AsNormalOrOptChainGet()
		if functionType.IsOrdinaryFunction() && !calleeIsGet {
			receiverType := functionType.GetTypeOfThis()
			if receiverType.IsUnknownType() ||
				receiverType.IsAllType() ||
				receiverType.IsVoidType() ||
				(receiverType.IsObjectType() && receiverType.ToObjectType().IsNativeObjectType()) {
				// Allow these special cases.
			} else {
				tc.Report(n, TypeCheck_EXPECTED_THIS_TYPE, functionType.String())
			}
		}

		tc.VisitArgumentList(call.AsCallOrNew(), functionType)
		tc.EnsureTypedWithType(n, functionType.GetReturnType())
	} else {
		tc.EnsureTyped(n)
	}

	// TODO(nicksantos): Add something to check for calls of RegExp objects,
	// which is not supported by IE. Either say something about the return type
	// or warn about the non-portability of the call or both.
}

// VisitArgumentList visits the parameters of a CALL or a NEW node.
//
// Java: TypeCheck.visitArgumentList(Node, FunctionType)
func (tc *TypeCheck) VisitArgumentList(call ast.CallOrNew, functionType rhino.FunctionType) {
	tc.CheckArgumentsMatchParameters(call.Node(), functionType, functionType.GetParameters(), 0)
}

// CheckArgumentsMatchParameters checks that a list of arguments match a list
// of formal parameters.
//
// If given a TAGGED_TEMPLATE_LIT, the given Iterator should only contain the
// parameters corresponding to the actual template lit sub arguments, skipping
// over the first parameter.
//
// firstParameterIndex is the index of the first parameter in the given
// Iterator in the function type's parameter list.
//
// Go-only: the Java gets the arguments as an Iterator over
// NodeUtil.getInvocationArgsAsIterable(call); here they are read from call,
// a CallOrNew (its arguments) or a TaggedTemplateLit (the substitutions of
// its template), and checkArgument is the body of the Java loop. The
// parameters are a slice (the elements the Java Iterator has left to
// yield).
//
// Java: TypeCheck.checkArgumentsMatchParameters(Node, FunctionType, Iterator<Node>, Iterator<Parameter>, int)
func (tc *TypeCheck) CheckArgumentsMatchParameters(
	call rhino.Node,
	functionType rhino.FunctionType,
	parameters []*rhino.FunctionType_Parameter,
	firstParameterIndex int) {

	spreadArgumentCount := 0
	normalArgumentCount := firstParameterIndex
	checkArgumentTypeAgainstParameter := true
	var parameter *rhino.FunctionType_Parameter
	nextParameter := 0
	checkArgument := func(argument ast.Arg) {
		// Count normal & spread arguments.
		if _, isSpread := argument.AsIterSpread(); isSpread {
			// we have some form of this case
			// someCall(arg1, arg2, ...firstSpreadExpression, argN, ...secondSpreadExpression)
			spreadArgumentCount++
			// Once we see a spread parameter, we can no longer match up arguments with parameters.
			checkArgumentTypeAgainstParameter = false
		} else {
			normalArgumentCount++
		}

		// Get the next parameter, if we're still matching parameters and arguments.
		if checkArgumentTypeAgainstParameter {
			if nextParameter < len(parameters) {
				parameter = parameters[nextParameter]
				nextParameter++
			} else if parameter != nil && parameter.IsVariadic() {
				// use varargs for all remaining parameters
			} else {
				// else we ran out of parameters and will report that after this loop
				parameter = nil
				checkArgumentTypeAgainstParameter = false
			}
		}

		if checkArgumentTypeAgainstParameter {
			tc.validator.ExpectArgumentMatchesParameter(
				argument.Node(), tc.GetJSType(argument.Node()), parameter.GetJSType(), call, normalArgumentCount)
		}
	}

	var callee ast.Expr
	if tagged, ok := ast.AsTaggedTemplateLit(call); ok {
		callee = tagged.Tag()
		for part := range tagged.Template().Parts() {
			if sub, ok := part.AsTemplateLitSub(); ok {
				checkArgument(sub.Expr().AsArg())
			}
		}
	} else {
		invocation := ast.MustCallOrNew(call)
		callee = invocation.Callee()
		for argument := range invocation.Args() {
			checkArgument(argument)
		}
	}

	minArity := functionType.GetMinArity()
	maxArity := functionType.GetMaxArity()

	if spreadArgumentCount > 0 {
		if normalArgumentCount > maxArity {
			// We cannot reliably check whether the total argument count is wrong, but we can at
			// least tell if there are more arguments than the function can handle even ignoring the
			// spreads.
			maxArityText := ""
			if maxArity != math.MaxInt32 {
				maxArityText = " and no more than " + strconv.Itoa(maxArity) + " argument(s)"
			}
			tc.Report(
				call,
				TypeCheck_WRONG_ARGUMENT_COUNT,
				tc.typeRegistry.GetReadableTypeNameNoDeref(callee.Node()),
				"at least "+strconv.Itoa(normalArgumentCount),
				strconv.Itoa(minArity),
				maxArityText)
		}
	} else {
		if minArity > normalArgumentCount || maxArity < normalArgumentCount {
			maxArityText := ""
			if maxArity != math.MaxInt32 {
				maxArityText = " and no more than " + strconv.Itoa(maxArity) + " argument(s)"
			}
			tc.Report(
				call,
				TypeCheck_WRONG_ARGUMENT_COUNT,
				tc.typeRegistry.GetReadableTypeNameNoDeref(callee.Node()),
				strconv.Itoa(normalArgumentCount),
				strconv.Itoa(minArity),
				maxArityText)
		}
	}
}

// VisitImplicitReturnExpression visits an arrow function expression body.
//
// Java: TypeCheck.visitImplicitReturnExpression(NodeTraversal, Node)
func (tc *TypeCheck) VisitImplicitReturnExpression(t *NodeTraversal, exprNode rhino.Node) {
	enclosingFunction := t.GetEnclosingFunction()
	jsType := tc.GetJSType(enclosingFunction)
	if jsType.IsFunctionType() {
		functionType := jsType.ToMaybeFunctionType()

		expectedReturnType := functionType.GetReturnType()
		// if no return type is specified, undefined must be returned
		// (it's a void function)
		if expectedReturnType == nil {
			expectedReturnType = tc.GetNativeType(rhino.JSTypeNative_VOID_TYPE)
		} else if enclosingFunction.IsAsyncFunction() {
			// Unwrap the async function's declared return type.
			expectedReturnType = Promises_CreateAsyncReturnableType(tc.typeRegistry, expectedReturnType)
		}

		// Fetch the returned value's type
		actualReturnType := tc.GetJSType(exprNode)

		tc.validator.ExpectCanAssignTo(
			exprNode, actualReturnType, expectedReturnType, "inconsistent return type")
	}
}

// VisitReturn visits a RETURN node.
//
// t is the node traversal object that supplies context, such as the scope
// chain to use in name lookups as well as error reporting; n the node being
// visited.
//
// Java: TypeCheck.visitReturn(NodeTraversal, Node)
func (tc *TypeCheck) VisitReturn(t *NodeTraversal, ret ast.Return) {
	enclosingFunction := t.GetEnclosingFunction()

	jsType := tc.GetJSType(enclosingFunction)

	if jsType.IsFunctionType() {
		functionType := jsType.ToMaybeFunctionType()

		returnType := functionType.GetReturnType()
		// if no return type is specified, undefined must be returned
		// (it's a void function)
		if returnType == nil {
			returnType = tc.GetNativeType(rhino.JSTypeNative_VOID_TYPE)
		} else if enclosingFunction.IsGeneratorFunction() {
			// Unwrap the template variable from a generator function's declared return type.
			// e.g. if returnType is "Generator<string, number, void>", the generator should return
			// "number".
			returnType = JsIterables_GetReturnElementType(returnType, tc.typeRegistry)

			if enclosingFunction.IsAsyncGeneratorFunction() {
				// Can return x|IThenable<x> in an AsyncGenerator<x>, no await needed. Note that we must
				// first wrap the type in IThenable as createAsyncReturnableType will map a non-IThenable
				// to `?`.
				returnType =
					Promises_CreateAsyncReturnableType(
						tc.typeRegistry, Promises_WrapInIThenable(tc.typeRegistry, returnType))
			}
		} else if enclosingFunction.IsAsyncFunction() {
			// e.g. `!Promise<string>` => `string|!IThenable<string>`
			// We transform the expected return type rather than the actual return type so that the
			// extual return type is always reported to the user. This was felt to be clearer.
			returnType = Promises_CreateAsyncReturnableType(tc.typeRegistry, returnType)
		} else if returnType.IsVoidType() && functionType.IsConstructor() {
			// Allow constructors to use empty returns for flow control.
			if _, hasValue := ret.Value(); !hasValue {
				return
			}

			// Allow constructors to return its own instance type
			returnType = functionType.GetInstanceType()
		}

		// fetching the returned value's type
		var valueNode rhino.Node
		var actualReturnType rhino.JSType
		if value, ok := ret.Value(); ok {
			valueNode = value.Node()
			actualReturnType = tc.GetJSType(valueNode)
		} else {
			actualReturnType = tc.GetNativeType(rhino.JSTypeNative_VOID_TYPE)
			valueNode = ret.Node()
		}

		// verifying
		tc.validator.ExpectCanAssignTo(
			valueNode, actualReturnType, returnType, "inconsistent return type")
	}
}

// VisitYield visits a YIELD node.
//
// Java: TypeCheck.visitYield(NodeTraversal, Node)
func (tc *TypeCheck) VisitYield(t *NodeTraversal, yield ast.Yield) {
	n := yield.Node()
	jsType := tc.GetJSType(t.GetEnclosingFunction())

	declaredYieldType := tc.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
	if jsType.IsFunctionType() {
		functionType := jsType.ToMaybeFunctionType()
		returnType := functionType.GetReturnType()
		declaredYieldType = JsIterables_GetElementType(returnType, tc.typeRegistry)

		if t.GetEnclosingFunction().IsAsyncGeneratorFunction() {
			// Can yield x|IThenable<x> in an AsyncGenerator<x>, no await needed. Note that we must
			// first wrap the type in IThenable as createAsyncReturnableType will map a non-IThenable to
			// `?`.
			declaredYieldType =
				Promises_CreateAsyncReturnableType(
					tc.typeRegistry, Promises_WrapInIThenable(tc.typeRegistry, declaredYieldType))
		}
	}

	// fetching the yielded value's type
	var valueNode rhino.Node
	var actualYieldType rhino.JSType
	if value, ok := yield.Value(); ok {
		valueNode = value.Node()
		actualYieldType = tc.GetJSType(valueNode)
	} else {
		actualYieldType = tc.GetNativeType(rhino.JSTypeNative_VOID_TYPE)
		valueNode = n
	}

	if n.IsYieldAll() {
		if t.GetEnclosingFunction().IsAsyncGeneratorFunction() {
			maybeActualYieldType :=
				tc.validator.ExpectAutoboxesToIterableOrAsyncIterable(
					n, actualYieldType, "Expression yield* expects an iterable or async iterable")
			if maybeActualYieldType == nil {
				// don't do any further typechecking of the yield* type.
				return
			}
			actualYieldType = maybeActualYieldType
		} else {
			if !tc.validator.ExpectAutoboxesToIterable(
				n, actualYieldType, "Expression yield* expects an iterable") {
				// don't do any further typechecking of the yield* type.
				return
			}
			templateType := tc.typeRegistry.GetIterableValueTemplate()
			actualYieldType =
				actualYieldType.Autobox().GetTemplateTypeMap().GetResolvedTemplateType(templateType)
		}
	}

	// verifying
	tc.validator.ExpectCanAssignTo(
		valueNode,
		actualYieldType,
		declaredYieldType,
		"Yielded type does not match declared return type.")
}

// Java: TypeCheck.visitTaggedTemplateLit(Node)
func (tc *TypeCheck) VisitTaggedTemplateLit(tagged ast.TaggedTemplateLit) {
	n := tagged.Node()
	tag := tagged.Tag().Node()
	tagType := tag.GetJSType().RestrictByNotNullOrUndefined()

	if !tagType.CanBeCalled() {
		tc.Report(n, TypeCheck_NOT_CALLABLE, tagType.String())
		return
	} else if !tagType.IsFunctionType() {
		// A few types, like the unknown, regexp, and bottom types, can be called as if they are
		// functions. Return if we have one of those types that is not actually a known function.
		return
	}

	tagFnType := tagType.ToMaybeFunctionType()
	parameters := tagFnType.GetParameters()

	// The tag function gets an array of all the template lit substitutions as its first argument,
	// but there's no actual AST node representing that array so we typecheck it separately from
	// the other tag arguments.

	// Validate that the tag function takes at least one parameter
	if len(parameters) == 0 {
		tc.Report(
			n,
			TypeCheck_WRONG_ARGUMENT_COUNT,
			tc.typeRegistry.GetReadableTypeNameNoDeref(tag),
			strconv.Itoa(NodeUtil_GetInvocationArgsCount(n)),
			"0",
			" and no more than 0 argument(s)")
		return
	}

	// Validate that the first parameter is a supertype of ITemplateArray
	firstParameter := parameters[0]
	parameters = parameters[1:]
	parameterType := firstParameter.GetJSType().RestrictByNotNullOrUndefined()
	if parameterType != nil {
		tc.validator.ExpectITemplateArraySupertype(
			tag, parameterType, "Invalid type for the first parameter of tag function")
	}

	// Validate the remaining parameters (the template literal substitutions)
	tc.CheckArgumentsMatchParameters(n, tagFnType, parameters, 1)
}

// VisitBinaryOperator unifies the type checking involved in the core binary
// operators and the corresponding assignment operators. The representation
// used internally is such that common code can handle both kinds of
// operators easily.
//
// op is the operator; n the node being checked.
//
// Go-only: n is an ast.Binary or an ast.AssignOp, which share no view, so
// the caller passes the operands: left and right, or target and value.
//
// Java: TypeCheck.visitBinaryOperator(Token, Node)
func (tc *TypeCheck) VisitBinaryOperator(op rhino.Token, n rhino.Node, left rhino.Node, right rhino.Node) {
	operatorType := tc.GetJSType(n)
	leftType := tc.GetJSType(left)
	rightType := tc.GetJSType(right)
	if operatorType.IsNoType() {
		// An attempt to mix bigint with other types
		tc.Report(
			n,
			TypeCheck_BINARY_OPERATION,
			stringPtrOrNull(NodeUtil_OpToStr(n.GetToken())),
			leftType.String(),
			rightType.String())
		return
	}
	switch op {
	case rhino.Token_ASSIGN_LSH, rhino.Token_ASSIGN_RSH, rhino.Token_LSH, rhino.Token_RSH, rhino.Token_ASSIGN_URSH,
		rhino.Token_URSH:
		if operatorType.IsNumber() {
			// TypeInference set the operator type to 'number', so we know bigint isn't involved.
			// NOTE: >>> and >>>= aren't valid operations for bigint, so TypeInference ignores the
			//     operand types for it and always sets the operator's type to 'number'. Thus,
			//     if one of the operands is a bigint, we'll end up reporting that as an error here.
			opStr := stringPtrOrNull(NodeUtil_OpToStr(n.GetToken()))
			if !leftType.MatchesNumberContext() {
				tc.Report(left, TypeCheck_BIT_OPERATION, opStr, leftType.String())
			} else {
				tc.validator.ExpectNumberStrict(n, leftType, "operator "+opStr)
			}
			if !rightType.MatchesNumberContext() {
				tc.Report(right, TypeCheck_BIT_OPERATION, opStr, rightType.String())
			} else {
				tc.validator.ExpectNumberStrict(n, rightType, "operator "+opStr)
			}
		} else {
			tc.validator.ExpectBigIntOrNumber(left, leftType, "left operand")
			tc.validator.ExpectBigIntOrNumber(right, rightType, "right operand")
		}
	case rhino.Token_ASSIGN_DIV,
		rhino.Token_ASSIGN_MOD,
		rhino.Token_ASSIGN_MUL,
		rhino.Token_ASSIGN_SUB,
		rhino.Token_ASSIGN_EXPONENT,
		rhino.Token_DIV,
		rhino.Token_MOD,
		rhino.Token_MUL,
		rhino.Token_SUB,
		rhino.Token_EXPONENT:
		if operatorType.IsNumber() {
			// TypeInference set the operator type to 'number', so we know bigint isn't involved.
			tc.validator.ExpectNumber(left, leftType, "left operand")
			tc.validator.ExpectNumber(right, rightType, "right operand")
		} else {
			tc.validator.ExpectBigIntOrNumber(left, leftType, "left operand")
			tc.validator.ExpectBigIntOrNumber(right, rightType, "right operand")
		}
	case rhino.Token_ASSIGN_BITAND, rhino.Token_ASSIGN_BITXOR, rhino.Token_ASSIGN_BITOR, rhino.Token_BITAND,
		rhino.Token_BITXOR, rhino.Token_BITOR:
		if operatorType.IsNumber() {
			// This condition is meant to catch any old cases (where bigint isn't involved)
			tc.validator.ExpectBitwiseable(left, leftType, "bad left operand to bitwise operator")
			tc.validator.ExpectBitwiseable(right, rightType, "bad right operand to bitwise operator")
		} else {
			tc.validator.ExpectBigIntOrNumber(left, leftType, "bad left operand to bitwise operator")
			tc.validator.ExpectBigIntOrNumber(right, rightType, "bad right operand to bitwise operator")
		}
	case rhino.Token_ASSIGN_ADD, rhino.Token_ADD:
	default:
		tc.Report(n, TypeCheck_UNEXPECTED_TOKEN, op.String())
	}
	tc.EnsureTyped(n)
}

// typeCheck_inputProvider is the lambda `(b) -> t.getInput()`, a
// CompilerInputProvider.
//
// Go-only: language-adapter — a lambda implementing a Go interface.
type typeCheck_inputProvider struct {
	t *NodeTraversal
}

func (p typeCheck_inputProvider) GetInput(*rhino.InputId) *CompilerInput {
	return p.t.GetInput()
}

// VisitModuleBody validates the implicit assignment to the global for a
// legacy goog.module.
//
// Java: TypeCheck.visitModuleBody(NodeTraversal, Node)
func (tc *TypeCheck) VisitModuleBody(t *NodeTraversal, body ast.ModuleBody) {
	moduleBody := body.Node()
	associatedModule :=
		ModuleImportResolver_GetModuleFromScopeRoot(
			tc.compiler.GetModuleMap(), typeCheck_inputProvider{t}, moduleBody)
	if !associatedModule.Metadata().IsLegacyGoogModule() {
		return
	}
	moduleName := rhino.QualifiedName_Of(*associatedModule.ClosureNamespace())
	// The first statement of a legacy goog.module is its goog.module call.
	googModuleCall, ok := body.FirstStatement()
	jutil.CheckState(ok)
	if moduleName.IsSimple() {
		globalVar := tc.topScope.GetVar(moduleName.Join())
		tc.validator.ExpectCanAssignTo(
			googModuleCall.Node(), moduleBody.GetJSType(), globalVar.GetType(), "legacy goog.module export")
	} else {
		parentType := tc.topScope.LookupQualifiedName(moduleName.GetOwner())
		var parentObjectType rhino.ObjectType
		if parentType != nil {
			parentObjectType = parentType.ToMaybeObjectType()
		}
		if parentObjectType == nil {
			return
		}
		tc.validator.expectCanAssignToPropertyOfQualifiedOwner(
			googModuleCall.Node(),
			moduleBody.GetJSType(),
			parentObjectType.GetPropertyType(moduleName.GetComponent()),
			parentObjectType,
			moduleName.GetOwner(),
			moduleName.GetComponent())
	}
}

// VisitDynamicImport validates that a dynamic import statement has a single
// child of type string.
//
// Java: TypeCheck.visitDynamicImport(NodeTraversal, Node)
func (tc *TypeCheck) VisitDynamicImport(t *NodeTraversal, dynamicImport ast.DynamicImport) {
	tc.EnsureTypedWithType(dynamicImport.Node(), tc.GetNativeType(rhino.JSTypeNative_PROMISE_TYPE))

	importSpecifier := dynamicImport.Specifier().Node()
	importSpecifierType := importSpecifier.GetJSType()
	if importSpecifierType == nil {
		tc.EnsureTypedWithNative(importSpecifier, rhino.JSTypeNative_STRING_TYPE)
	} else {
		tc.validator.ExpectNotNullOrUndefined(
			t,
			importSpecifier,
			importSpecifierType,
			"dynamic import specifier",
			tc.GetNativeType(rhino.JSTypeNative_STRING_TYPE))
	}
}
