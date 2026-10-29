// Port of com/google/javascript/jscomp/TypeValidator.java @ closure-compiler v20261021

/*
 * Copyright 2009 The Closure Compiler Authors.
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
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// TypeValidator is a central reporter for all type violations: places where
// the programmer has annotated a variable (or property) with one type, but
// has assigned another type to it.
//
// Also doubles as a central repository for all type violations, so that
// type-based optimizations (like AmbiguateProperties) can be fault-tolerant.
type TypeValidator struct {
	compiler                AbstractCompiler
	typeRegistry            *rhino.JSTypeRegistry
	allBitwisableValueTypes rhino.JSType
	nullOrUndefined         rhino.JSType
	promiseOfUnknownType    rhino.JSType
	iterableOrAsyncIterable rhino.JSType

	// In TypeCheck, when we are analyzing a file with .java.js suffix, we set
	// this field to IGNORE_NULL_UNDEFINED
	subtypingMode rhino.JSType_SubtypingMode

	mismatches TypeMismatch_Accumulator
}

// User warnings
var (
	TypeValidator_INVALID_CAST = DiagnosticType_Warning(
		"JSC_INVALID_CAST",
		"invalid cast - must be a subtype or supertype\n"+
			"from: {0}\n"+
			"to  : {1}")

	TypeValidator_TYPE_MISMATCH_WARNING = DiagnosticType_Warning("JSC_TYPE_MISMATCH", "{0}")

	TypeValidator_INVALID_ASYNC_RETURN_TYPE = DiagnosticType_Warning(
		"JSC_INVALID_ASYNC_RETURN_TYPE",
		"The return type of an async function must be a supertype of Promise\n"+"found: {0}")

	TypeValidator_INVALID_OPERAND_TYPE = DiagnosticType_Disabled("JSC_INVALID_OPERAND_TYPE", "{0}")

	TypeValidator_MISSING_EXTENDS_TAG_WARNING = DiagnosticType_Warning("JSC_MISSING_EXTENDS_TAG", "Missing @extends tag on type {0}")

	TypeValidator_DUP_VAR_DECLARATION = DiagnosticType_Warning(
		"JSC_DUP_VAR_DECLARATION", "variable {0} redefined, original definition at {1}:{2}")

	TypeValidator_DUP_VAR_DECLARATION_TYPE_MISMATCH = DiagnosticType_Warning(
		"JSC_DUP_VAR_DECLARATION_TYPE_MISMATCH",
		"variable {0} redefined with type {1}, original definition at {2}:{3} with type {4}")

	TypeValidator_INTERFACE_METHOD_NOT_IMPLEMENTED = DiagnosticType_Warning(
		"JSC_INTERFACE_METHOD_NOT_IMPLEMENTED",
		"property {0} on interface {1} is not implemented by type {2}")

	TypeValidator_HIDDEN_INTERFACE_PROPERTY_MISMATCH = DiagnosticType_Warning(
		"JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH",
		"mismatch of the {0} property on type {4} and the type "+
			"of the property it overrides from interface {1}\n"+
			"original: {2}\n"+
			"override: {3}")

	// TODO(goktug): consider updating super class / interface property mismatch to follow the same
	// pattern.
	TypeValidator_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH = DiagnosticType_Warning(
		"JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH",
		"mismatch of the {0} property type and the type "+
			"of the property it overrides from superclass {1}\n"+
			"original: {2}\n"+
			"override: {3}")

	TypeValidator_ABSTRACT_METHOD_NOT_IMPLEMENTED = DiagnosticType_Warning(
		"JSC_ABSTRACT_METHOD_NOT_IMPLEMENTED",
		"property {0} on abstract class {1} is not implemented by type {2}")

	TypeValidator_UNKNOWN_TYPEOF_VALUE = DiagnosticType_Warning("JSC_UNKNOWN_TYPEOF_VALUE", "unknown type: {0}")

	TypeValidator_ILLEGAL_PROPERTY_ACCESS = DiagnosticType_Warning("JSC_ILLEGAL_PROPERTY_ACCESS", "Cannot do {0} access on a {1}")

	TypeValidator_ALL_DIAGNOSTICS = NewDiagnosticGroup(
		TypeValidator_ABSTRACT_METHOD_NOT_IMPLEMENTED,
		TypeValidator_DUP_VAR_DECLARATION,
		TypeValidator_DUP_VAR_DECLARATION_TYPE_MISMATCH,
		TypeValidator_HIDDEN_INTERFACE_PROPERTY_MISMATCH,
		TypeValidator_ILLEGAL_PROPERTY_ACCESS,
		TypeValidator_INTERFACE_METHOD_NOT_IMPLEMENTED,
		TypeValidator_INVALID_ASYNC_RETURN_TYPE,
		TypeValidator_INVALID_CAST,
		TypeValidator_MISSING_EXTENDS_TAG_WARNING,
		TypeValidator_TYPE_MISMATCH_WARNING,
		TypeValidator_UNKNOWN_TYPEOF_VALUE)
)

// Java: TypeValidator(AbstractCompiler)
func NewTypeValidator(compiler AbstractCompiler) *TypeValidator {
	v := &TypeValidator{
		subtypingMode: rhino.JSType_SubtypingMode_NORMAL,
	}
	v.compiler = compiler
	v.typeRegistry = compiler.GetTypeRegistry()
	v.allBitwisableValueTypes = v.typeRegistry.CreateUnionTypeOfNatives(
		rhino.JSTypeNative_STRING_TYPE, rhino.JSTypeNative_NUMBER_TYPE, rhino.JSTypeNative_BOOLEAN_TYPE,
		rhino.JSTypeNative_NULL_TYPE, rhino.JSTypeNative_VOID_TYPE)
	v.nullOrUndefined = v.typeRegistry.GetNativeType(rhino.JSTypeNative_NULL_VOID)
	v.promiseOfUnknownType = v.typeRegistry.CreateTemplatizedType(
		v.typeRegistry.GetNativeObjectType(rhino.JSTypeNative_PROMISE_TYPE),
		v.typeRegistry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE))
	v.iterableOrAsyncIterable = v.typeRegistry.CreateUnionType(
		v.typeRegistry.GetNativeObjectType(rhino.JSTypeNative_ITERATOR_TYPE),
		v.typeRegistry.GetNativeObjectType(rhino.JSTypeNative_ASYNC_ITERATOR_TYPE))
	return v
}

// TypeValidator_GetInstanceOfCtor is a utility function that attempts to
// get an instance type from a potential constructor type
//
// Java: TypeValidator.getInstanceOfCtor(JSType)
func TypeValidator_GetInstanceOfCtor(t rhino.JSType) rhino.ObjectType {
	if t == nil {
		return nil
	}
	ctor := rhino.JSType_ToMaybeFunctionType(t.Dereference())
	if ctor != nil && ctor.IsConstructor() {
		return ctor.GetInstanceType()
	}
	return nil
}

// GetMismatches gets a list of type violations.
//
// For each violation, one element is the expected type and the other is the
// type that is actually found. Order is not significant.
//
// NOTE(dimvar): Even though TypeMismatch is a pair, the passes that call
// this method never use it as a pair; they just add both its elements to a
// set of invalidating types. Consider just maintaining a set of types here
// instead of a set of type pairs.
//
// Java: TypeValidator.getMismatches()
func (v *TypeValidator) GetMismatches() []TypeMismatch {
	return v.mismatches.GetMismatches()
}

// Java: TypeValidator.setSubtypingMode(SubtypingMode)
func (v *TypeValidator) SetSubtypingMode(mode rhino.JSType_SubtypingMode) {
	v.subtypingMode = mode
}

// All non-private methods should have the form:
// expectCondition(Node n, ...);
// If there is a mismatch, the {@code expect} method should issue
// a warning and attempt to correct the mismatch, when possible.

// Java: TypeValidator.expectValidTypeofName(Node, String)
func (v *TypeValidator) ExpectValidTypeofName(n rhino.Node, found string) {
	v.compiler.Report(JSError_MakeWithNode(n, TypeValidator_UNKNOWN_TYPEOF_VALUE, found))
}

// ExpectObject expects the type to be an object, or a type convertible to
// object. If the expectation is not met, issue a warning at the provided
// node's source code position.
//
// Returns true if there was no warning, false if there was a mismatch.
//
// Java: TypeValidator.expectObject(Node, JSType, String)
func (v *TypeValidator) ExpectObject(n rhino.Node, type_ rhino.JSType, msg string) bool {
	if !type_.MatchesObjectContext() {
		v.MismatchNative(n, msg, type_, rhino.JSTypeNative_OBJECT_TYPE)
		return false
	}
	return true
}

// ExpectActualObject expects the type to be an object. Unlike expectObject,
// a type convertible to object is not acceptable.
//
// Java: TypeValidator.expectActualObject(Node, JSType, String)
func (v *TypeValidator) ExpectActualObject(n rhino.Node, type_ rhino.JSType, msg string) {
	if !type_.IsObject() {
		v.MismatchNative(n, msg, type_, rhino.JSTypeNative_OBJECT_TYPE)
	}
}

// ExpectAnyObject expects the type to contain an object sometimes. If the
// expectation is not met, issue a warning at the provided node's source
// code position.
//
// Java: TypeValidator.expectAnyObject(Node, JSType, String)
func (v *TypeValidator) ExpectAnyObject(n rhino.Node, type_ rhino.JSType, msg string) {
	anyObjectType := v.GetNativeType(rhino.JSTypeNative_NO_OBJECT_TYPE)
	if !anyObjectType.IsSubtypeOf(type_) && !type_.IsEmptyType() {
		v.Mismatch(n, msg, type_, anyObjectType)
	}
}

// ExpectAutoboxesToIterable expects the type to autobox to be an Iterable.
//
// Returns true if there was no warning, false if there was a mismatch.
//
// Java: TypeValidator.expectAutoboxesToIterable(Node, JSType, String)
func (v *TypeValidator) ExpectAutoboxesToIterable(n rhino.Node, type_ rhino.JSType, msg string) bool {
	// Note: we don't just use JSType.autobox() here because that removes null and undefined.
	// We want to keep null and undefined around.
	if type_.IsUnionType() {
		for _, alt := range type_.ToMaybeUnionType().GetAlternates() {
			if alt.IsBoxableScalar() {
				alt = alt.AutoboxesTo()
			}
			if !alt.IsSubtypeOf(v.GetNativeType(rhino.JSTypeNative_ITERABLE_TYPE)) {
				v.MismatchNative(n, msg, type_, rhino.JSTypeNative_ITERABLE_TYPE)
				return false
			}
		}

	} else {
		autoboxedType := type_
		if type_.IsBoxableScalar() {
			autoboxedType = type_.AutoboxesTo()
		}
		if !autoboxedType.IsSubtypeOf(v.GetNativeType(rhino.JSTypeNative_ITERABLE_TYPE)) {
			v.MismatchNative(n, msg, type_, rhino.JSTypeNative_ITERABLE_TYPE)
			return false
		}
	}
	return true
}

// ExpectAutoboxesToIterableOrAsyncIterable expects the type to autobox to
// be an Iterable or AsyncIterable.
//
// Returns the unwrapped variants of the iterable(s), or empty if not
// iterable (Go: nil for Optional.empty()).
//
// Java: TypeValidator.expectAutoboxesToIterableOrAsyncIterable(Node, JSType, String)
func (v *TypeValidator) ExpectAutoboxesToIterableOrAsyncIterable(n rhino.Node, type_ rhino.JSType, msg string) rhino.JSType {
	maybeBoxed := JsIterables_MaybeBoxIterableOrAsyncIterable(type_, v.typeRegistry)

	if maybeBoxed.IsMatch() {
		// Optional.of
		return jutil.CheckNotNull(maybeBoxed.GetTemplatedType())
	}

	v.Mismatch(n, msg, type_, v.iterableOrAsyncIterable)

	return nil
}

// ExpectGeneratorSupertype expects the type to be a Generator or supertype
// of Generator.
//
// Java: TypeValidator.expectGeneratorSupertype(Node, JSType, String)
func (v *TypeValidator) ExpectGeneratorSupertype(n rhino.Node, type_ rhino.JSType, msg string) {
	if !v.GetNativeType(rhino.JSTypeNative_GENERATOR_TYPE).IsSubtypeOf(type_) {
		v.MismatchNative(n, msg, type_, rhino.JSTypeNative_GENERATOR_TYPE)
	}
}

// ExpectAsyncGeneratorSupertype expects the type to be a AsyncGenerator or
// supertype of AsyncGenerator.
//
// Java: TypeValidator.expectAsyncGeneratorSupertype(Node, JSType, String)
func (v *TypeValidator) ExpectAsyncGeneratorSupertype(n rhino.Node, type_ rhino.JSType, msg string) {
	if !v.GetNativeType(rhino.JSTypeNative_ASYNC_GENERATOR_TYPE).IsSubtypeOf(type_) {
		v.MismatchNative(n, msg, type_, rhino.JSTypeNative_ASYNC_GENERATOR_TYPE)
	}
}

// ExpectValidAsyncReturnType expects the type to be a supertype of
// `Promise`.
//
// `Promise` is the lower bound of the declared return type, since that's
// what async functions always return; the user can't return an instance of
// a more specific type.
//
// Java: TypeValidator.expectValidAsyncReturnType(Node, JSType)
func (v *TypeValidator) ExpectValidAsyncReturnType(n rhino.Node, type_ rhino.JSType) {
	if v.promiseOfUnknownType.IsSubtypeOf(type_) {
		return
	}

	err := JSError_MakeWithNode(n, TypeValidator_INVALID_ASYNC_RETURN_TYPE, type_.String())
	v.RegisterMismatchAndReport(type_, v.promiseOfUnknownType, err)
}

// ExpectITemplateArraySupertype expects the type to be an ITemplateArray or
// supertype of ITemplateArray.
//
// Java: TypeValidator.expectITemplateArraySupertype(Node, JSType, String)
func (v *TypeValidator) ExpectITemplateArraySupertype(n rhino.Node, type_ rhino.JSType, msg string) {
	if !v.GetNativeType(rhino.JSTypeNative_I_TEMPLATE_ARRAY_TYPE).IsSubtypeOf(type_) {
		v.MismatchNative(n, msg, type_, rhino.JSTypeNative_I_TEMPLATE_ARRAY_TYPE)
	}
}

// ExpectString expects the type to be a string, or a type convertible to
// string. If the expectation is not met, issue a warning at the provided
// node's source code position.
//
// Java: TypeValidator.expectString(Node, JSType, String)
func (v *TypeValidator) ExpectString(n rhino.Node, type_ rhino.JSType, msg string) {
	if !type_.MatchesStringContext() {
		v.MismatchNative(n, msg, type_, rhino.JSTypeNative_STRING_TYPE)
	}
}

// ExpectNumber expects the type to be a number, or a type convertible to
// number. If the expectation is not met, issue a warning at the provided
// node's source code position.
//
// Java: TypeValidator.expectNumber(Node, JSType, String)
func (v *TypeValidator) ExpectNumber(n rhino.Node, type_ rhino.JSType, msg string) {
	if !type_.MatchesNumberContext() {
		v.MismatchNative(n, msg, type_, rhino.JSTypeNative_NUMBER_TYPE)
	} else {
		v.ExpectNumberStrict(n, type_, msg)
	}
}

// ExpectNumberStrict expects the type to be a number or a subtype.
//
// Java: TypeValidator.expectNumberStrict(Node, JSType, String)
func (v *TypeValidator) ExpectNumberStrict(n rhino.Node, type_ rhino.JSType, msg string) {
	if !type_.IsSubtypeOf(v.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)) {
		v.RegisterMismatchAndReportWithDiagnostic(
			n, TypeValidator_INVALID_OPERAND_TYPE, msg, type_, v.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE), nil, nil)
	}
}

// Java: TypeValidator.expectMatchingTypesStrict(Node, JSType, JSType, String)
func (v *TypeValidator) ExpectMatchingTypesStrict(n rhino.Node, left rhino.JSType, right rhino.JSType, msg string) {
	if !left.IsSubtypeOf(right) && !right.IsSubtypeOf(left) {
		v.RegisterMismatchAndReportWithDiagnostic(n, TypeValidator_INVALID_OPERAND_TYPE, msg, right, left, nil, nil)
	}
}

// ExpectBitwiseable expects the type to be a valid operand to a bitwise
// operator. This includes numbers, any type convertible to a number, or any
// other primitive type (undefined|null|boolean|string).
//
// Java: TypeValidator.expectBitwiseable(Node, JSType, String)
func (v *TypeValidator) ExpectBitwiseable(n rhino.Node, type_ rhino.JSType, msg string) {
	if !type_.MatchesNumberContext() && !type_.IsSubtypeOf(v.allBitwisableValueTypes) {
		v.Mismatch(n, msg, type_, v.allBitwisableValueTypes)
	} else {
		v.ExpectNumberStrict(n, type_, msg)
	}
}

// ExpectNumberOrSymbol expects the type to be a number or string, or a type
// convertible to a number or symbol. If the expectation is not met, issue a
// warning at the provided node's source code position.
//
// Java: TypeValidator.expectNumberOrSymbol(Node, JSType, String)
func (v *TypeValidator) ExpectNumberOrSymbol(n rhino.Node, type_ rhino.JSType, msg string) {
	if !type_.MatchesNumberContext() && !type_.MatchesSymbolContext() {
		v.MismatchNative(n, msg, type_, rhino.JSTypeNative_NUMBER_SYMBOL)
	}
}

// ExpectStringOrSymbol expects the type to be a string or symbol, or a type
// convertible to a string. If the expectation is not met, issue a warning at
// the provided node's source code position.
//
// Java: TypeValidator.expectStringOrSymbol(Node, JSType, String)
func (v *TypeValidator) ExpectStringOrSymbol(n rhino.Node, type_ rhino.JSType, msg string) {
	if !type_.MatchesStringContext() && !type_.MatchesSymbolContext() {
		v.MismatchNative(n, msg, type_, rhino.JSTypeNative_STRING_SYMBOL)
	}
}

// ExpectUnknownOrComparable expects the type to be unknown or a comparable
// type (bigint, number, or string)
//
// Java: TypeValidator.expectUnknownOrComparable(Node, JSType, String)
func (v *TypeValidator) ExpectUnknownOrComparable(n rhino.Node, type_ rhino.JSType, msg string) {
	if !type_.IsSubtypeOf(v.GetNativeType(rhino.JSTypeNative_BIGINT_NUMBER_STRING)) &&
		!type_.IsSubtypeOf(v.GetNativeType(rhino.JSTypeNative_BIGINT_NUMBER_STRING_OBJECT)) {
		if type_.MatchesNumberContext() || type_.MatchesStringContext() {
			// INVALID_OPERAND_TYPE is suppressed unless strict type checking is enabled
			v.RegisterMismatchAndReportWithDiagnostic(
				n, TypeValidator_INVALID_OPERAND_TYPE, msg, type_, v.GetNativeType(rhino.JSTypeNative_BIGINT_NUMBER_STRING), nil, nil)
		} else {
			v.MismatchNative(n, msg, type_, rhino.JSTypeNative_BIGINT_NUMBER_STRING)
		}
	}
}

// ExpectStringOrNumberOrSymbol expects the type to be a number or string or
// symbol, or a type convertible to a number or string. If the expectation is
// not met, issue a warning at the provided node's source code position.
//
// Java: TypeValidator.expectStringOrNumberOrSymbol(Node, JSType, String)
func (v *TypeValidator) ExpectStringOrNumberOrSymbol(n rhino.Node, type_ rhino.JSType, msg string) {
	if !type_.MatchesNumberContext() &&
		!type_.MatchesStringContext() &&
		!type_.MatchesSymbolContext() {
		v.MismatchNative(n, msg, type_, rhino.JSTypeNative_NUMBER_STRING_SYMBOL)
	} else {
		v.ExpectStringOrNumberOrSymbolStrict(n, type_, msg)
	}
}

// Java: TypeValidator.expectStringOrNumberOrSymbolStrict(Node, JSType, String)
func (v *TypeValidator) ExpectStringOrNumberOrSymbolStrict(n rhino.Node, type_ rhino.JSType, msg string) {
	if !type_.IsSubtypeOf(v.GetNativeType(rhino.JSTypeNative_NUMBER_STRING_SYMBOL)) {
		v.RegisterMismatchAndReportWithDiagnostic(
			n, TypeValidator_INVALID_OPERAND_TYPE, msg, type_, v.GetNativeType(rhino.JSTypeNative_NUMBER_STRING_SYMBOL), nil, nil)
	}
}

// ExpectBigIntOrNumber expects the type to be a bigint or number, or a type
// convertible to number. If the expectation is not met, issue a warning at
// the provided node's source code position.
//
// Java: TypeValidator.expectBigIntOrNumber(Node, JSType, String)
func (v *TypeValidator) ExpectBigIntOrNumber(n rhino.Node, type_ rhino.JSType, msg string) {
	if !type_.IsSubtypeOf(v.GetNativeType(rhino.JSTypeNative_BIGINT_NUMBER)) &&
		!type_.IsSubtypeOf(v.GetNativeType(rhino.JSTypeNative_BIGINT_NUMBER_OBJECT)) {
		if type_.MatchesNumberContext() {
			// INVALID_OPERAND_TYPE is suppressed unless strict type checking is enabled
			v.RegisterMismatchAndReportWithDiagnostic(
				n, TypeValidator_INVALID_OPERAND_TYPE, msg, type_, v.GetNativeType(rhino.JSTypeNative_BIGINT_NUMBER), nil, nil)
		} else {
			v.MismatchNative(n, msg, type_, rhino.JSTypeNative_BIGINT_NUMBER)
		}
	}
}

// ExpectNotNullOrUndefined expects the type to be anything but the null or
// void type. If the expectation is not met, issue a warning at the provided
// node's source code position. Note that a union type that includes the void
// type and at least one other type meets the expectation.
//
// Returns whether the expectation was met.
//
// Java: TypeValidator.expectNotNullOrUndefined(NodeTraversal, Node, JSType, String, JSType)
func (v *TypeValidator) ExpectNotNullOrUndefined(
	t *NodeTraversal, n rhino.Node, type_ rhino.JSType, msg string, expectedType rhino.JSType) bool {
	if !type_.IsNoType() &&
		!type_.IsUnknownType() &&
		type_.IsSubtypeOf(v.nullOrUndefined) &&
		!TypeValidator_ContainsForwardDeclaredUnresolvedName(type_) {

		// There's one edge case right now that we don't handle well, and
		// that we don't want to warn about.
		// if (this.x == null) {
		//   this.initializeX();
		//   this.x.foo();
		// }
		// In this case, we incorrectly type x because of how we
		// infer properties locally. See issue 109.
		// http://blickly.github.io/closure-compiler-issues/#109
		//
		// We do not do this inference globally.
		if n.IsGetProp() && !t.InGlobalScope() && type_.IsNullType() {
			return true
		}

		v.Mismatch(n, msg, type_, expectedType)
		return false
	}
	return true
}

// Java: TypeValidator.containsForwardDeclaredUnresolvedName(JSType)
func TypeValidator_ContainsForwardDeclaredUnresolvedName(type_ rhino.JSType) bool {
	if type_.IsUnionType() {
		for _, alt := range type_.ToMaybeUnionType().GetAlternates() {
			if TypeValidator_ContainsForwardDeclaredUnresolvedName(alt) {
				return true
			}
		}
	}
	return type_.IsNoResolvedType()
}

// ExpectSwitchMatchesCase expects that the type of a switch condition
// matches the type of its case condition.
//
// Java: TypeValidator.expectSwitchMatchesCase(Node, JSType, JSType)
func (v *TypeValidator) ExpectSwitchMatchesCase(n rhino.Node, switchType rhino.JSType, caseType rhino.JSType) {
	// ECMA-262, page 68, step 3 of evaluation of CaseBlock
	if !switchType.CanTestForShallowEqualityWith(caseType) {
		v.Mismatch(n.GetFirstChild(), "case expression doesn't match switch", caseType, switchType)
	}
}

// ExpectIndexMatch expects that the first type can be addressed with
// GETELEM syntax and that the second type is the right type for an index
// into the first type.
//
// n is the GETELEM or COMPUTED_PROP node to issue warnings on; objType is
// the type we're indexing into (the left side of the GETELEM); indexType is
// the type inside the brackets of the GETELEM/COMPUTED_PROP.
//
// Java: TypeValidator.expectIndexMatch(Node, JSType, JSType)
func (v *TypeValidator) ExpectIndexMatch(n rhino.Node, objType rhino.JSType, indexType rhino.JSType) {
	jutil.CheckStateMsg1(
		n.IsGetElem() || n.IsOptChainGetElem() || n.IsComputedProp() || n.IsComputedFieldDef(), "%s", n)
	var indexNode rhino.Node
	if n.IsGetElem() {
		indexNode = n.GetLastChild()
	} else {
		indexNode = n.GetFirstChild()
	}
	if indexType.IsSymbolValueType() {
		// For now, allow symbols definitions/access on any type. In the future only allow them
		// on the subtypes for which they are defined.
		return
	}
	if objType.IsUnknownType() {
		v.ExpectStringOrNumberOrSymbol(indexNode, indexType, "property access")
		return
	}
	dereferenced := objType.Dereference()
	if dereferenced != nil &&
		dereferenced.GetTemplateTypeMap().HasTemplateKey(v.typeRegistry.GetObjectIndexKey()) {
		templateType := dereferenced.
			GetTemplateTypeMap().
			GetResolvedTemplateType(v.typeRegistry.GetObjectIndexKey())
		v.ExpectCanAssignTo(indexNode, indexType, templateType, "restricted index type")
	} else if dereferenced != nil && dereferenced.IsArrayType() {
		v.ExpectNumberOrSymbol(indexNode, indexType, "array access")
	} else if objType.IsStruct() {
		v.compiler.Report(JSError_MakeWithNode(indexNode, TypeValidator_ILLEGAL_PROPERTY_ACCESS, "'[]'", "struct"))
	} else if objType.MatchesObjectContext() {
		v.ExpectStringOrSymbol(indexNode, indexType, "property access")
	} else {
		v.Mismatch(
			n,
			"only arrays or objects can be accessed",
			objType,
			v.typeRegistry.CreateUnionTypeOfNatives(rhino.JSTypeNative_ARRAY_TYPE, rhino.JSTypeNative_OBJECT_TYPE))
	}
}

// ExpectCanAssignToPropertyOf expects that the first type can be assigned to
// a symbol of the second type.
//
// n is the node to issue warnings on; rightType is the type on the RHS of
// the assign; leftType is the type of the symbol on the LHS of the assign;
// owner is the owner of the property being assigned to; propName is the name
// of the property being assigned to. Returns true if the types matched,
// false otherwise.
//
// Java: TypeValidator.expectCanAssignToPropertyOf(Node, JSType, JSType, Node, String)
func (v *TypeValidator) ExpectCanAssignToPropertyOf(
	n rhino.Node, rightType rhino.JSType, leftType rhino.JSType, owner rhino.Node, propName string) bool {
	var classType rhino.FunctionType
	if n.IsMemberFieldDef() {
		classType = n.GetGrandparent().GetJSType().AssertFunctionType()
	}
	ok, isTemplateMismatch := v.checkCanAssignToPropertyOf(rightType, leftType, v.GetJSType(owner))
	if ok {
		return true
	}
	var ownerTypeName string
	if classType != nil {
		ownerTypeName = classType.GetInstanceType().String()
	} else {
		ownerTypeName = v.typeRegistry.GetReadableTypeName(owner)
	}
	v.reportPropertyAssignmentMismatch(n, rightType, leftType, propName, ownerTypeName, isTemplateMismatch)
	return false
}

// expectCanAssignToPropertyOfQualifiedOwner checks property assignment when the
// owner name is a QualifiedName, avoiding any closure or Join() allocation on
// the valid fast path.
func (v *TypeValidator) expectCanAssignToPropertyOfQualifiedOwner(
	n rhino.Node,
	rightType rhino.JSType,
	leftType rhino.JSType,
	ownerType rhino.JSType,
	ownerQName rhino.QualifiedName,
	propName string,
) bool {
	ok, isTemplateMismatch := v.checkCanAssignToPropertyOf(rightType, leftType, ownerType)
	if ok {
		return true
	}
	v.reportPropertyAssignmentMismatch(n, rightType, leftType, propName, ownerQName.Join(), isTemplateMismatch)
	return false
}

func (v *TypeValidator) checkCanAssignToPropertyOf(
	rightType rhino.JSType, leftType rhino.JSType, ownerType rhino.JSType,
) (ok bool, isTemplateMismatch bool) {
	if leftType.IsTemplateType() {
		left := leftType.ToMaybeTemplateType()
		if rightType.ContainsReferenceAncestor(left) ||
			rightType.IsUnknownType() ||
			left.IsUnknownType() {
			// The only time we can assign to a variable with a template type is if the value assigned
			// has a type that explicitly has it as a supertype.
			// Otherwise, the template type is existential and it is unknown whether or not it is a
			// proper super type.
			return true, false
		}
		return false, true
	}
	// The NoType check is a hack to make typedefs work OK.
	if !leftType.IsNoType() && !rightType.IsSubtypeOf(leftType) {
		// Do not type-check interface methods, because we expect that
		// they will have dummy implementations that do not match the type
		// annotations.
		if ownerType.IsFunctionPrototypeType() {
			ownerFn := ownerType.ToObjectType().GetOwnerFunction()
			if ownerFn.IsInterface() && rightType.IsFunctionType() && leftType.IsFunctionType() {
				return true, false
			}
		}
		return false, false
	}
	return true, false
}

func (v *TypeValidator) reportPropertyAssignmentMismatch(
	n rhino.Node,
	rightType rhino.JSType,
	leftType rhino.JSType,
	propName string,
	ownerTypeName string,
	isTemplateMismatch bool,
) {
	msg := "assignment to property " + propName + " of " + ownerTypeName
	if isTemplateMismatch {
		v.RegisterMismatchAndReportWithDiagnostic(
			n,
			TypeValidator_TYPE_MISMATCH_WARNING,
			msg,
			rightType,
			leftType,
			nil,
			nil)
		return
	}
	v.Mismatch(n, msg, rightType, leftType)
}

// ExpectCanAssignTo expects that the first type can be assigned to a symbol
// of the second type.
//
// n is the node to issue warnings on; rightType is the type on the RHS of
// the assign; leftType is the type of the symbol on the LHS of the assign;
// msg is an extra message for the mismatch warning, if necessary. Returns
// true if the types matched, false otherwise.
//
// Java: TypeValidator.expectCanAssignTo(Node, JSType, JSType, String)
func (v *TypeValidator) ExpectCanAssignTo(n rhino.Node, rightType rhino.JSType, leftType rhino.JSType, msg string) bool {
	if leftType.IsTemplateType() {
		left := leftType.ToMaybeTemplateType()
		if rightType.ContainsReferenceAncestor(left) ||
			rightType.IsUnknownType() ||
			left.IsUnknownType() {
			// The only time we can assign to a variable with a template type is if the value assigned
			// has a type that explicitly has it as a supertype.
			// Otherwise, the template type is existential and it is unknown whether or not it is a
			// proper super type.
			return true
		} else {
			v.RegisterMismatchAndReportWithDiagnostic(
				n,
				TypeValidator_TYPE_MISMATCH_WARNING,
				msg,
				rightType,
				leftType,
				nil,
				nil)
			return false
		}
	}
	if !rightType.IsSubtypeOf(leftType) {
		v.Mismatch(n, msg, rightType, leftType)
		return false
	}
	return true
}

// ExpectArgumentMatchesParameter expects that the type of an argument
// matches the type of the parameter that it's fulfilling.
//
// n is the node to issue warnings on; argType is the type of the argument;
// paramType is the type of the parameter; callNode is the call node, to help
// with the warning message; ordinal is the argument ordinal, to help with
// the warning message.
//
// Java: TypeValidator.expectArgumentMatchesParameter(Node, JSType, JSType, Node, int)
func (v *TypeValidator) ExpectArgumentMatchesParameter(
	n rhino.Node, argType rhino.JSType, paramType rhino.JSType, callNode rhino.Node, ordinal int) {
	if !argType.IsSubtypeOf(paramType) {
		v.Mismatch(
			n,
			"actual parameter "+
				strconv.Itoa(ordinal)+
				" of "+
				v.typeRegistry.GetReadableTypeNameNoDeref(callNode.GetFirstChild())+
				" does not match formal parameter",
			argType,
			paramType)
	}
}

// ExpectSuperType expects that the first type is the direct superclass of
// the second type.
//
// n is the node where warnings should point to; superObject is the expected
// super instance type; subObject is the sub instance type.
//
// Java: TypeValidator.expectSuperType(Node, ObjectType, ObjectType)
func (v *TypeValidator) ExpectSuperType(n rhino.Node, superObject rhino.ObjectType, subObject rhino.ObjectType) {
	subCtor := subObject.GetConstructor()
	implicitProto := subObject.GetImplicitPrototype()
	var declaredSuper rhino.ObjectType
	if implicitProto != nil {
		declaredSuper = implicitProto.GetImplicitPrototype()
	}
	if declaredSuper != nil && declaredSuper.IsTemplatizedType() {
		declaredSuper = declaredSuper.ToMaybeTemplatizedType().GetReferencedType()
	}
	_, superIsUnknown := superObject.(*rhino.UnknownType)
	if declaredSuper != nil &&
		!superIsUnknown &&
		!declaredSuper.Equals(superObject) {
		if declaredSuper.Equals(v.GetNativeType(rhino.JSTypeNative_OBJECT_TYPE)) {
			v.RegisterMismatchAndReport(
				superObject,
				declaredSuper,
				JSError_MakeWithNode(n, TypeValidator_MISSING_EXTENDS_TAG_WARNING, subObject.String()))
		} else {
			v.Mismatch(n, "mismatch in declaration of superclass type", superObject, declaredSuper)
		}

		// Correct the super type.
		if !subCtor.HasCachedValues() {
			subCtor.SetPrototypeBasedOn(superObject)
		}
	}
}

// ExpectExtends expects that an ES6 class's extends clause is actually a
// supertype of the given class. Compares the registered supertype, which is
// taken from the JSDoc if present, otherwise from the AST, with the type in
// the extends node of the AST.
//
// n is the node where warnings should point to; subCtor is the sub
// constructor type; astSuperCtor is the expected super constructor from the
// extends node in the AST.
//
// Java: TypeValidator.expectExtends(Node, FunctionType, FunctionType)
func (v *TypeValidator) ExpectExtends(n rhino.Node, subCtor rhino.FunctionType, astSuperCtor rhino.FunctionType) {
	if astSuperCtor == nil || (!astSuperCtor.IsConstructor() && !astSuperCtor.IsInterface()) {
		// toMaybeFunctionType failed, or we've got a loose type.  Let it go for now.
		return
	}
	if astSuperCtor.IsConstructor() != subCtor.IsConstructor() {
		// Don't bother looking if one is a constructor and the other is an interface.
		// We'll report an error elsewhere.
		return
	}
	astSuperInstance := astSuperCtor.GetInstanceType()
	if subCtor.IsConstructor() {
		// There should be exactly one superclass, and it needs to have this constructor.
		// Note: if the registered supertype (from the @extends jsdoc) was unresolved,
		// then getSuperClassConstructor will be null - make sure not to crash.
		registeredSuperCtor := subCtor.GetSuperClassConstructor()
		if registeredSuperCtor != nil {
			registeredSuperInstance := registeredSuperCtor.GetInstanceType()
			if !astSuperInstance.Equals(registeredSuperInstance) {
				v.Mismatch(
					n,
					"mismatch in declaration of superclass type",
					astSuperInstance,
					registeredSuperInstance)
			}
		}
	} else if subCtor.IsInterface() {
		// We intentionally skip this check for interfaces because they can extend multiple other
		// interfaces.
	}
}

// ExpectCanAssignToPrototype expects that it's valid to assign something to
// a given type's prototype.
//
// Most of these checks occur during TypedScopeCreator, so we just handle
// very basic cases here
//
// For example, assuming `Foo` is a constructor, `Foo.prototype = 3;` will
// warn because `3` is not an object.
//
// ownerType is the type of the object whose prototype is being changed
// (e.g. `Foo` above); node is the node to issue warnings on (e.g. `3`
// above); rightType is the rvalue type being assigned to the prototype (e.g.
// `number` above).
//
// Java: TypeValidator.expectCanAssignToPrototype(JSType, Node, JSType)
func (v *TypeValidator) ExpectCanAssignToPrototype(ownerType rhino.JSType, node rhino.Node, rightType rhino.JSType) {
	if ownerType.IsFunctionType() {
		functionType := ownerType.ToMaybeFunctionType()
		if functionType.IsConstructor() {
			v.ExpectObject(node, rightType, "cannot override prototype with non-object")
		}
	}
}

// ExpectCanCast expects that the first type can be cast to the second type.
// The first type must have some relationship with the second.
//
// n is the node where warnings should point; targetType is the type being
// cast to; sourceType is the type being cast from.
//
// Java: TypeValidator.expectCanCast(Node, JSType, JSType)
func (v *TypeValidator) ExpectCanCast(n rhino.Node, targetType rhino.JSType, sourceType rhino.JSType) {
	if !sourceType.CanCastTo(targetType) {
		v.RegisterMismatchAndReport(
			sourceType,
			targetType,
			JSError_MakeWithNode(n, TypeValidator_INVALID_CAST, sourceType.String(), targetType.String()))
		// The "canCastTo" check is intentionally looser than the subtyping check, but we still want
		// to record potential mismatches for disambiguation safety.
	} else if !sourceType.IsSubtypeOf(targetType) {
		v.mismatches.RegisterMismatch(n, sourceType, targetType)
	}
}
