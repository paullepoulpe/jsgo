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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CheckGetterOrSetterType verifies that a user did not give a getter and
// setter different types, as the type system will arbitrarily take the
// first of the types if different.
//
// Java: TypeCheck.checkGetterOrSetterType(Node, FunctionType)
func (tc *TypeCheck) CheckGetterOrSetterType(child ast.GetterOrSetterDef, classType rhino.FunctionType) {
	propertyName := child.Key()

	methodType := child.Function().Node().GetJSType().ToMaybeFunctionType()
	_, isGetter := child.AsGetterDef()
	var propertyType rhino.JSType
	if isGetter {
		propertyType = tc.DetermineGetterType(methodType)
	} else {
		propertyType = typeCheck_iterablesGetOnlyElement(methodType.GetParameters()).GetJSType()
	}
	var officialPropertyType rhino.JSType
	if child.Node().IsStaticMember() {
		officialPropertyType = classType.GetPropertyType(propertyName)
	} else {
		officialPropertyType = classType.GetPrototype().GetPropertyType(propertyName)
	}
	if !propertyType.Equals(officialPropertyType) {
		// TODO(b/116797078): make this not an error - instead, store the getter and setter types
		// separately
		var getterType, setterType string
		if isGetter {
			getterType = propertyType.String()
			setterType = officialPropertyType.String()
		} else {
			getterType = officialPropertyType.String()
			setterType = propertyType.String()
		}
		tc.Report(
			child.Node(),
			TypeCheck_CONFLICTING_GETTER_SETTER_TYPE,
			propertyName,
			getterType,
			setterType)
	}
}

// typeCheck_iterablesGetOnlyElement is Guava's
// Iterables.getOnlyElement(Iterable): NoSuchElementException when empty,
// IllegalArgumentException listing up to five elements when longer.
//
// Go-only: language-adapter — Guava helper.
func typeCheck_iterablesGetOnlyElement(params []*rhino.FunctionType_Parameter) *rhino.FunctionType_Parameter {
	if len(params) == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	if len(params) == 1 {
		return params[0]
	}
	var sb strings.Builder
	sb.WriteString("expected one element but was: <")
	sb.WriteString(params[0].String())
	i := 1
	for ; i < 5 && i < len(params); i++ {
		sb.WriteString(", ")
		sb.WriteString(params[i].String())
	}
	if i < len(params) {
		sb.WriteString(", ...")
	}
	sb.WriteByte('>')
	panic(jutil.NewIllegalArgumentExceptionMsg(sb.String()))
}

// Java: TypeCheck.determineGetterType(FunctionType)
func (tc *TypeCheck) DetermineGetterType(methodType rhino.FunctionType) rhino.JSType {
	// TODO(sdh): consider only falling back on unknown if the function body is empty?  But we
	// need to not report a conflicting type error if there's different unknowns.
	if !methodType.IsReturnTypeInferred() {
		return methodType.GetReturnType()
	}
	return tc.typeRegistry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
}

// Java: TypeCheck.checkPropertyInheritanceOnAssignment(Node, Node, Property.Key, JSDocInfo, JSType)
func (tc *TypeCheck) CheckPropertyInheritanceOnAssignment(
	assign rhino.Node, object ast.Expr, property rhino.Property_Key, info *rhino.JSDocInfo, propertyType rhino.JSType) {
	// Inheritance checks for prototype properties.
	//
	// TODO(nicksantos): This isn't the right place to do this check. We
	// really want to do this when we're looking at the constructor.
	// We'd find all its properties and make sure they followed inheritance
	// rules, like we currently do for @implements to make sure
	// all the methods are implemented.
	//
	// As-is, this misses many other ways to override a property.

	if prototype, ok := object.AsGetProp(); ok && prototype.Node().GetStringID() == typeCheckPrototypeID {
		// ASSIGN = assign
		//   GETPROP
		//     GETPROP = object
		//       ? = preObject
		//       STRING = "prototype"
		//     STRING = property

		preObject := prototype.Object()
		ctorType := tc.GetJSType(preObject.Node()).ToMaybeFunctionType()
		if ctorType == nil || !ctorType.HasInstanceType() {
			return
		}

		tc.CheckDeclaredPropertyAgainstNominalInheritance(
			assign, ctorType, property, info, propertyType)
		tc.CheckAbstractMethodInConcreteClass(assign, ctorType, info)
	} else {
		// ASSIGN = assign
		//   GETPROP
		//     ? = object
		//     STRING = property
		// or
		//  ASSIGN = assign
		//    GETELEM
		//      ? = object
		//      ? = property

		// We only care about checking a static property assignment.
		ctorType := tc.GetJSType(object.Node()).ToMaybeFunctionType()
		if ctorType == nil || !ctorType.HasInstanceType() {
			return
		}

		tc.CheckDeclaredPropertyAgainstPrototypalInheritance(
			assign, ctorType, property, info, propertyType)
	}
}

// Java: TypeCheck.checkPropertyInheritanceOnPrototypeLitKey(Node, Property.Key, ObjectType)
func (tc *TypeCheck) CheckPropertyInheritanceOnPrototypeLitKey(
	key rhino.Node, propertyName rhino.Property_Key, type_ rhino.ObjectType) {
	// Inheritance checks for prototype objlit properties.
	//
	// TODO(nicksantos): This isn't the right place to do this check. We
	// really want to do this when we're looking at the constructor.
	// We'd find all its properties and make sure they followed inheritance
	// rules, like we currently do for @implements to make sure
	// all the methods are implemented.
	//
	// As-is, this misses many other ways to override a property.
	//
	// object.prototype = { key: function() {} };
	tc.CheckPropertyInheritance(key, propertyName, type_.GetOwnerFunction(), type_)
}

// Java: TypeCheck.checkPropertyInheritanceOnClassMember(Node, Property.Key, FunctionType)
func (tc *TypeCheck) CheckPropertyInheritanceOnClassMember(
	key rhino.Node, propertyName rhino.Property_Key, ctorType rhino.FunctionType) {
	if key.IsStaticMember() {
		tc.CheckDeclaredPropertyAgainstPrototypalInheritance(
			key, ctorType, propertyName, key.GetJSDocInfo(), ctorType.GetPropertyTypeByKey(propertyName))
	} else {
		tc.CheckPropertyInheritance(key, propertyName, ctorType, ctorType.GetInstanceType())
	}
}

// Java: TypeCheck.checkPropertyInheritance(Node, Property.Key, FunctionType, ObjectType)
func (tc *TypeCheck) CheckPropertyInheritance(
	key rhino.Node, propertyName rhino.Property_Key, ctorType rhino.FunctionType, type_ rhino.ObjectType) {
	if ctorType == nil || !ctorType.HasInstanceType() {
		return
	}

	// The node reported on is the first child of key: the value of a string
	// key or a field, the function of a method, getter or setter, the key of
	// a computed property.
	tc.CheckDeclaredPropertyAgainstNominalInheritance(
		key.GetFirstChild(),
		ctorType,
		propertyName,
		key.GetJSDocInfo(),
		type_.GetPropertyTypeByKey(propertyName))
	tc.CheckAbstractMethodInConcreteClass(key, ctorType, key.GetJSDocInfo())
}

// CheckDeclaredPropertyAgainstNominalInheritance: given a ctorType, check that
// the property (propertyName), on the corresponding instance type
// (receiverType), conforms to inheritance rules.
//
// This method only checks nominal inheritance (via extends and implements
// declarations). Compare to CheckDeclaredPropertyAgainstPrototypalInheritance.
//
// To be conformant, the propertyName must
//   - Carry the @override annotation iff it is an override.
//   - Be typed as a subtype of the type of propertyName on all supertypes of
//     receiverType.
//
// Java: TypeCheck.checkDeclaredPropertyAgainstNominalInheritance(Node, FunctionType, Property.Key, JSDocInfo, JSType)
func (tc *TypeCheck) CheckDeclaredPropertyAgainstNominalInheritance(
	n rhino.Node,
	ctorType rhino.FunctionType,
	propertyName rhino.Property_Key,
	info *rhino.JSDocInfo,
	propertyType rhino.JSType) {

	// No need to check special properties; @override is not required for them, nor they are
	// manually typed by the developers.
	if propertyName.Matches("__proto__") || propertyName.Matches("constructor") {
		return
	}

	// Interfaces are checked elsewhere.
	if ctorType.IsInterface() {
		return
	}

	// If the supertype doesn't resolve correctly, we've warned about this already.
	if TypeCheck_HasUnknownOrEmptySupertype(ctorType) {
		return
	}

	foundProperty := false

	superCtor := ctorType.GetSuperClassConstructor()
	if superCtor != nil {
		propSlot := superCtor.GetInstanceType().FindClosestDefinitionByKey(propertyName)
		superClassHasProperty := propSlot != nil && !propSlot.IsOwnedByInterface()
		foundProperty = foundProperty || superClassHasProperty

		if superClassHasProperty {
			superClass := propSlot.GetOwnerInstanceType()
			superClassHasDeclaredProperty := !propSlot.GetValue().IsTypeInferred()
			if superClassHasDeclaredProperty {
				if TypeCheck_IsDeclaredLocally(ctorType, propertyName) && !TypeCheck_DeclaresOverride(info) {
					tc.compiler.Report(
						JSError_MakeWithNode(
							n,
							TypeCheck_HIDDEN_SUPERCLASS_PROPERTY,
							propertyName.HumanReadableName(),
							stringPtrOrNull(superClass.GetReferenceName())))
				}
				tc.validator.CheckPropertyType(
					n, ctorType.GetTypeOfThis(), superClass, propertyName, propertyType)
			}
		}
	}

	for _, implementedInterface := range ctorType.GetAllImplementedInterfaces() {
		if implementedInterface.IsUnknownType() || implementedInterface.IsEmptyType() {
			continue
		}
		jutil.CheckStateMsg1(implementedInterface.IsInstanceType(), "%s", implementedInterface)
		propSlot := implementedInterface.FindClosestDefinitionByKey(propertyName)
		interfaceHasProperty := propSlot != nil
		foundProperty = foundProperty || interfaceHasProperty

		if interfaceHasProperty && !TypeCheck_DeclaresOverride(info) {
			tc.compiler.Report(
				JSError_MakeWithNode(
					n,
					TypeCheck_HIDDEN_INTERFACE_PROPERTY,
					propertyName.HumanReadableName(),
					stringPtrOrNull(propSlot.GetOwnerInstanceType().GetReferenceName())))
		}
	}

	if !foundProperty && TypeCheck_DeclaresOverride(info) {
		tc.compiler.Report(
			JSError_MakeWithNode(
				n,
				TypeCheck_UNKNOWN_OVERRIDE,
				propertyName.HumanReadableName(),
				stringPtrOrNull(ctorType.GetInstanceType().GetReferenceName())))
	}
}

// Java: TypeCheck.isDeclaredLocally(FunctionType, Property.Key)
func TypeCheck_IsDeclaredLocally(ctorType rhino.FunctionType, propertyName rhino.Property_Key) bool {
	jutil.CheckState(ctorType.IsConstructor())
	return ctorType.GetPrototype().HasOwnPropertyByKey(propertyName) ||
		ctorType.GetInstanceType().HasOwnPropertyByKey(propertyName)
}

// CheckDeclaredPropertyAgainstPrototypalInheritance: given a receiverType,
// check that the property (propertyName) conforms to inheritance rules.
//
// This method only checks prototypal inheritance (via the prototype chain).
// Compare to CheckDeclaredPropertyAgainstNominalInheritance.
//
// To be conformant, the propertyName must
//   - Carry the @override annotation iff it is an override.
//   - Be typed as a subtype of the type of propertyName on all supertypes of
//     receiverType.
//
// Java: TypeCheck.checkDeclaredPropertyAgainstPrototypalInheritance(Node, ObjectType, Property.Key, JSDocInfo, JSType)
func (tc *TypeCheck) CheckDeclaredPropertyAgainstPrototypalInheritance(
	n rhino.Node,
	receiverType rhino.ObjectType,
	propertyName rhino.Property_Key,
	info *rhino.JSDocInfo,
	propertyType rhino.JSType) {
	// TODO(nickreid): Right now this is only expected to run on ctors. However, it wouldn't be bad
	// if it ran on more things. Consider a precondition on the arguments.

	declaredOverride := TypeCheck_DeclaresOverride(info)
	var supertypeWithProperty rhino.ObjectType
	for type_ := receiverType.GetImplicitPrototype(); type_ != nil; type_ = type_.GetImplicitPrototype() {
		// We want to report the supertype that actually had the overridden declaration.
		if type_.HasOwnPropertyByKey(propertyName) {
			// We only care about the lowest match in the chain because it must be the most
			// specific.
			supertypeWithProperty = type_
			break
		}
	}

	if supertypeWithProperty == nil {
		// TODO(b/144327372): stop loosening typechecking for forward supertype references once
		// the type system correctly models this case.
		if declaredOverride && !receiverType.LoosenTypecheckingDueToForwardReferencedSupertype() {
			tc.compiler.Report(
				JSError_MakeWithNode(
					n, //
					TypeCheck_UNKNOWN_PROTOTYPAL_OVERRIDE,
					propertyName.HumanReadableName(),
					receiverType.String()))
		}
	} else {
		if !declaredOverride {
			tc.compiler.Report(
				JSError_MakeWithNode(
					n, //
					TypeCheck_HIDDEN_PROTOTYPAL_SUPERTYPE_PROPERTY,
					propertyName.HumanReadableName(),
					supertypeWithProperty.String()))
		}

		overriddenPropertyType := supertypeWithProperty.GetPropertyTypeByKey(propertyName)
		if !propertyType.IsSubtypeOf(overriddenPropertyType) {
			tc.compiler.Report(
				JSError_MakeWithNode(
					n,
					TypeCheck_HIDDEN_PROTOTYPAL_SUPERTYPE_PROPERTY_MISMATCH,
					propertyName.HumanReadableName(),
					supertypeWithProperty.String(),
					overriddenPropertyType.String(),
					propertyType.String()))
		}
	}
}

// Java: TypeCheck.checkAbstractMethodInConcreteClass(Node, FunctionType, JSDocInfo)
func (tc *TypeCheck) CheckAbstractMethodInConcreteClass(n rhino.Node, ctorType rhino.FunctionType, info *rhino.JSDocInfo) {
	if info == nil || !info.IsAbstract() {
		return
	}

	if ctorType.IsConstructor() && !ctorType.IsAbstract() {
		tc.Report(n, TypeCheck_ABSTRACT_METHOD_IN_CONCRETE_CLASS)
	}
}

// TypeCheck_HasUnknownOrEmptySupertype: given a constructor or an interface
// type, find out whether the unknown type is a supertype of the current type.
//
// Java: TypeCheck.hasUnknownOrEmptySupertype(FunctionType)
func TypeCheck_HasUnknownOrEmptySupertype(ctor rhino.FunctionType) bool {
	jutil.CheckArgument(ctor.IsConstructor() || ctor.IsInterface())
	jutil.CheckArgument(!ctor.IsUnknownType())

	// The type system should notice inheritance cycles on its own
	// and break the cycle.
	for {
		maybeSuperInstanceType := ctor.GetPrototype().GetImplicitPrototype()
		if maybeSuperInstanceType == nil {
			return false
		}
		if maybeSuperInstanceType.IsUnknownType() || maybeSuperInstanceType.IsEmptyType() {
			return true
		}
		ctor = maybeSuperInstanceType.GetConstructor()
		if ctor == nil {
			return false
		}
		jutil.CheckState(ctor.IsConstructor() || ctor.IsInterface())
	}
}

// VisitInterfacePropertyAssignment visits an lvalue node for cases such as
//
//	interface.prototype.property = ...;
//
// Java: TypeCheck.visitInterfacePropertyAssignment(Node, Node)
func (tc *TypeCheck) VisitInterfacePropertyAssignment(object ast.Expr, lvalue ast.GetProp) {
	assign, ok := ast.AsAssign(lvalue.Node().GetParent())
	if !ok {
		// assignments to interface properties cannot be in destructuring patterns or for-of loops
		tc.ReportInvalidInterfaceMemberDeclaration(object.Node())
		return
	}
	rvalue := assign.Value()
	rvalueType := tc.GetJSType(rvalue.Node())

	// Only 2 values are allowed for interface methods:
	//    goog.abstractMethod
	//    function () {};
	// Other (non-method) interface properties must be stub declarations without assignments, e.g.
	//     someinterface.prototype.nonMethodProperty;
	// which is why we enforce that `rvalueType.isFunctionType()`.
	if !rvalueType.IsFunctionType() {
		tc.ReportInvalidInterfaceMemberDeclaration(object.Node())
	}

	if fn, ok := rvalue.AsFunction(); ok && !NodeUtil_IsEmptyBlock(fn.Body().Node()) {
		abstractMethodName := tc.compiler.GetCodingConvention().GetAbstractMethodName()
		tc.compiler.Report(JSError_MakeWithNode(object.Node(), TypeCheck_INTERFACE_METHOD_NOT_EMPTY, stringPtrOrNull(abstractMethodName)))
	}
}

// Java: TypeCheck.reportInvalidInterfaceMemberDeclaration(Node)
func (tc *TypeCheck) ReportInvalidInterfaceMemberDeclaration(interfaceNode rhino.Node) {
	abstractMethodName := tc.compiler.GetCodingConvention().GetAbstractMethodName()
	// This is bad i18n style but we don't localize our compiler errors.
	abstractMethodMessage := ""
	if abstractMethodName != nil {
		abstractMethodMessage = ", or " + *abstractMethodName
	}
	tc.compiler.Report(
		JSError_MakeWithNode(interfaceNode, TypeCheck_INVALID_INTERFACE_MEMBER_DECLARATION, abstractMethodMessage))
}

// CheckEnumAlias checks enum aliases.
//
// We verify that the enum element type of the enum used for initialization
// is a subtype of the enum element type of the enum the value is being copied
// in.
//
// Example:
//
//	var myEnum = myOtherEnum;
//
// Enum aliases are irregular, so we need special code for this :(
//
// valueType is the type of the value used for initialization of the enum;
// nodeToWarn the node on which to issue warnings on.
//
// Java: TypeCheck.checkEnumAlias(NodeTraversal, JSDocInfo, JSType, Node)
func (tc *TypeCheck) CheckEnumAlias(
	t *NodeTraversal, declInfo *rhino.JSDocInfo, valueType rhino.JSType, nodeToWarn rhino.Node) {
	if declInfo == nil || !declInfo.HasEnumParameterType() {
		return
	}

	if !valueType.IsEnumType() {
		return
	}

	valueEnumType := valueType.ToMaybeEnumType()
	valueEnumPrimitiveType := valueEnumType.GetElementsType().GetPrimitiveType()
	tc.validator.ExpectCanAssignTo(
		nodeToWarn,
		valueEnumPrimitiveType,
		declInfo.GetEnumParameterType().Evaluate(t.GetTypedScope(), tc.typeRegistry),
		"incompatible enum element types")
}

// CheckJsdocInfoContainsObjectWithBadKey checks if current node contains js
// docs and checks all types specified in the js doc whether they have Objects
// with potentially invalid keys. For example: `Object<!Object, number>`. If
// such type is found, a warning is reported for the current node.
//
// Java: TypeCheck.checkJsdocInfoContainsObjectWithBadKey(Node)
func (tc *TypeCheck) CheckJsdocInfoContainsObjectWithBadKey(n rhino.Node) {
	if info := n.GetJSDocInfo(); info != nil {
		tc.checkJsdocInfoContainsObjectWithBadKeySlow(n, info)
	}
}

func (tc *TypeCheck) checkJsdocInfoContainsObjectWithBadKeySlow(n rhino.Node, info *rhino.JSDocInfo) {
	tc.CheckTypeContainsObjectWithBadKey(n, info.GetType())
	tc.CheckTypeContainsObjectWithBadKey(n, info.GetReturnType())
	tc.CheckTypeContainsObjectWithBadKey(n, info.GetTypedefType())
	paramCount := info.GetParameterCount()
	for i := 0; i < paramCount; i++ {
		tc.CheckTypeContainsObjectWithBadKey(n, info.GetParameterTypeAt(i))
	}
}

// Java: TypeCheck.checkTypeContainsObjectWithBadKey(Node, JSTypeExpression)
func (tc *TypeCheck) CheckTypeContainsObjectWithBadKey(n rhino.Node, type_ *rhino.JSTypeExpression) {
	if type_ != nil && type_.GetRoot().GetJSType() != nil {
		realType := type_.GetRoot().GetJSType()
		if !realType.IsTemplatizedType() && !realType.IsUnionType() && !realType.IsOrdinaryFunction() {
			return
		}
		if tc.badKeyCheckedTypes == nil {
			tc.badKeyCheckedTypes = make(map[rhino.JSType]struct{}, 8)
		} else {
			clear(tc.badKeyCheckedTypes)
		}
		objectWithBadKey :=
			tc.FindObjectWithNonStringifiableKey(realType, tc.badKeyCheckedTypes)
		clear(tc.badKeyCheckedTypes)
		if objectWithBadKey != nil {
			tc.compiler.Report(JSError_MakeWithNode(n, TypeCheck_NON_STRINGIFIABLE_OBJECT_KEY, objectWithBadKey.String()))
		}
	}
}

// IsReasonableObjectPropertyKey checks whether type is useful as the key of
// an object property access. This means it should be either stringifiable or
// a symbol. Stringifiable types are types that can be converted to string and
// give unique results for different objects. For example objects have native
// toString() method that on chrome returns "[object Object]" for all objects
// making it useless when used as keys. At the same time native types like
// numbers can be safely converted to strings and used as keys. Also user
// might have provided custom toString() methods for a class making it
// suitable for using as key.
//
// Java: TypeCheck.isReasonableObjectPropertyKey(JSType)
func (tc *TypeCheck) IsReasonableObjectPropertyKey(type_ rhino.JSType) bool {
	// Check built-in types
	if type_.IsUnknownType() ||
		type_.IsNumber() ||
		type_.IsString() ||
		type_.IsSymbol() ||
		type_.IsBooleanObjectType() ||
		type_.IsBooleanValueType() ||
		type_.IsDateType() ||
		type_.IsRegexpType() ||
		type_.IsInterface() ||
		type_.IsRecordType() ||
		type_.IsNullType() ||
		type_.IsVoidType() {
		return true
	}

	// For enums check that underlying type is stringifiable.
	if type_.ToMaybeEnumElementType() != nil {
		return tc.IsReasonableObjectPropertyKey(type_.ToMaybeEnumElementType().GetPrimitiveType())
	}

	// Array is stringifiable if it doesn't have template type or if it does have it, the template
	// type must be also stringifiable.
	// Good: Array, Array.<number>
	// Bad: Array.<!Object>
	if type_.IsArrayType() {
		return true
	}
	if type_.IsTemplatizedType() {
		templatizedType := type_.ToMaybeTemplatizedType()
		if templatizedType.GetReferencedType().IsArrayType() {
			return tc.IsReasonableObjectPropertyKey(templatizedType.GetTemplateTypes()[0])
		}
	}

	// Named types are usually @typedefs. For such types we need to check underlying type specified
	// in @typedef annotation.
	if namedType, ok := type_.(*rhino.NamedType); ok {
		return tc.IsReasonableObjectPropertyKey(namedType.GetReferencedType())
	}

	// For union type every alternate must be stringifiable.
	if type_.IsUnionType() {
		for _, alternateType := range type_.ToMaybeUnionType().GetAlternates() {
			if !tc.IsReasonableObjectPropertyKey(alternateType) {
				return false
			}
		}
		return true
	}

	// Handle interfaces and classes.
	if type_.IsObject() {
		objectType := type_.ToMaybeObjectType()
		constructor := objectType.GetConstructor()
		// Interfaces considered stringifiable as user might implement toString() method in
		// classes-implementations.
		if constructor != nil && constructor.IsInterface() {
			return true
		}
		// This is user-defined class so check if it has custom toString() method.
		return tc.ClassHasToString(objectType)
	}
	return false
}

// IsObjectTypeWithNonStringifiableKey checks whether current type is Object
// type with non-stringifable key.
//
// Java: TypeCheck.isObjectTypeWithNonStringifiableKey(JSType)
func (tc *TypeCheck) IsObjectTypeWithNonStringifiableKey(type_ rhino.JSType) bool {
	if !type_.IsTemplatizedType() {
		// TODO(nickreid): Why don't we care about types like `Foo extends Object<Bar, Qux>`?
		return false
	}

	templateTypeMap := type_.GetTemplateTypeMap()
	objectIndexKey := tc.typeRegistry.GetObjectIndexKey()
	if templateTypeMap.HasTemplateKey(objectIndexKey) {
		return !tc.IsReasonableObjectPropertyKey(
			templateTypeMap.GetResolvedTemplateType(objectIndexKey))
	} else {
		return false
	}
}

// FindObjectWithNonStringifiableKey checks whether type (or one of its
// component if is composed type like union or templatized type) has Object
// with non-stringifiable key. For example `Object.<!Object, number>`.
//
// Returns non-stringifiable type which is used as key or null if all there
// are no such types.
//
// Java: TypeCheck.findObjectWithNonStringifiableKey(JSType, Set<JSType>)
func (tc *TypeCheck) FindObjectWithNonStringifiableKey(
	type_ rhino.JSType, alreadyCheckedTypes map[rhino.JSType]struct{}) rhino.JSType {
	if !type_.IsTemplatizedType() && !type_.IsUnionType() && !type_.IsOrdinaryFunction() {
		return nil
	}
	if _, seen := alreadyCheckedTypes[type_]; seen {
		// This can happen in recursive types. Current type already being checked earlier in
		// stacktrace so now we just skip it.
		return nil
	}
	alreadyCheckedTypes[type_] = struct{}{}
	if tc.IsObjectTypeWithNonStringifiableKey(type_) {
		return type_
	}
	if type_.IsUnionType() {
		for _, alternateType := range type_.ToMaybeUnionType().GetAlternates() {
			result := tc.FindObjectWithNonStringifiableKey(alternateType, alreadyCheckedTypes)
			if result != nil {
				return result
			}
		}
	}
	if type_.IsTemplatizedType() {
		for _, templateType := range type_.ToMaybeTemplatizedType().GetTemplateTypes() {
			result := tc.FindObjectWithNonStringifiableKey(templateType, alreadyCheckedTypes)
			if result != nil {
				return result
			}
		}
	}
	if type_.IsOrdinaryFunction() {
		function := type_.ToMaybeFunctionType()
		for _, parameter := range function.GetParameters() {
			result :=
				tc.FindObjectWithNonStringifiableKey(parameter.GetJSType(), alreadyCheckedTypes)
			if result != nil {
				return result
			}
		}
		return tc.FindObjectWithNonStringifiableKey(function.GetReturnType(), alreadyCheckedTypes)
	}
	return nil
}

// ClassHasToString checks whether class has overridden toString() method.
// All objects has native toString() method but we ignore it as it is not
// useful so we need user-provided toString() method.
//
// Java: TypeCheck.classHasToString(ObjectType)
func (tc *TypeCheck) ClassHasToString(type_ rhino.ObjectType) bool {
	toStringProperty := type_.GetOwnSlot("toString")
	if toStringProperty != nil {
		return toStringProperty.GetType().IsFunctionType()
	}
	parent := type_.GetImplicitPrototype()
	if parent != nil && !parent.IsNativeObjectType() {
		return tc.ClassHasToString(parent)
	}
	return false
}
