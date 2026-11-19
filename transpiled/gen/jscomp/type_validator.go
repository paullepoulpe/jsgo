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
// Port of com/google/javascript/jscomp/TypeValidator.java (generated mechanically by j2go).

package jscomp

import (
	"strconv"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.TypeValidator
type TypeValidator struct {
	compiler                AbstractCompiler                                 `gson:"-"`
	typeRegistry            *pkg_rhino.JSTypeRegistry                        `gson:"typeRegistry"`
	allBitwisableValueTypes pkg_rhino.JSType                                 `gson:"allBitwisableValueTypes"`
	nullOrUndefined         pkg_rhino.JSType                                 `gson:"nullOrUndefined"`
	promiseOfUnknownType    pkg_rhino.JSType                                 `gson:"promiseOfUnknownType"`
	iterableOrAsyncIterable pkg_rhino.JSType                                 `gson:"iterableOrAsyncIterable"`
	mismatches              *TypeMismatch_Accumulator                        `gson:"mismatches"`
	reusableVerifier        *TypeValidator_WellFormedTemplatizedTypeVerifier `gson:"-"`
	idHash_                 atomic.Int32
	subtypingMode           pkg_rhino.JSType_SubtypingMode `gson:"subtypingMode"`
	verifierInUse           bool                           `gson:"-"`
}

func (this_ *TypeValidator) BindSelf_(s any) {
	_ = s
}
func (this_ *TypeValidator) IsJavaNil_() bool { return this_ == nil }
func (*TypeValidator) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TypeValidator", "TypeValidator"
}
func (this_ *TypeValidator) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TypeValidator) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TypeValidator", "com.google.javascript.jscomp.TypeValidator", func(v_ any) bool {
		return j2gort.InstanceOf[*TypeValidator](v_)
	})
}

func StaticGet_TypeValidator_FOUND_REQUIRED() string { return "{0}\nfound   : {1}\nrequired: {2}" }
func StaticGet_TypeValidator_FOUND_REQUIRED_MISSING() string {
	return "{0}\nfound   : {1}\nrequired: {2}\nmissing : [{3}]\nmismatch: [{4}]"
}

var static_TypeValidator_INVALID_CAST *DiagnosticType

func StaticGet_TypeValidator_INVALID_CAST() *DiagnosticType {
	clinit_TypeValidator()
	return static_TypeValidator_INVALID_CAST
}
func StaticSet_TypeValidator_INVALID_CAST(v *DiagnosticType) {
	clinit_TypeValidator()
	static_TypeValidator_INVALID_CAST = v
}

var static_TypeValidator_TYPE_MISMATCH_WARNING *DiagnosticType

func StaticGet_TypeValidator_TYPE_MISMATCH_WARNING() *DiagnosticType {
	clinit_TypeValidator()
	return static_TypeValidator_TYPE_MISMATCH_WARNING
}
func StaticSet_TypeValidator_TYPE_MISMATCH_WARNING(v *DiagnosticType) {
	clinit_TypeValidator()
	static_TypeValidator_TYPE_MISMATCH_WARNING = v
}

var static_TypeValidator_INVALID_ASYNC_RETURN_TYPE *DiagnosticType

func StaticGet_TypeValidator_INVALID_ASYNC_RETURN_TYPE() *DiagnosticType {
	clinit_TypeValidator()
	return static_TypeValidator_INVALID_ASYNC_RETURN_TYPE
}
func StaticSet_TypeValidator_INVALID_ASYNC_RETURN_TYPE(v *DiagnosticType) {
	clinit_TypeValidator()
	static_TypeValidator_INVALID_ASYNC_RETURN_TYPE = v
}

var static_TypeValidator_INVALID_OPERAND_TYPE *DiagnosticType

func StaticGet_TypeValidator_INVALID_OPERAND_TYPE() *DiagnosticType {
	clinit_TypeValidator()
	return static_TypeValidator_INVALID_OPERAND_TYPE
}
func StaticSet_TypeValidator_INVALID_OPERAND_TYPE(v *DiagnosticType) {
	clinit_TypeValidator()
	static_TypeValidator_INVALID_OPERAND_TYPE = v
}

var static_TypeValidator_MISSING_EXTENDS_TAG_WARNING *DiagnosticType

func StaticGet_TypeValidator_MISSING_EXTENDS_TAG_WARNING() *DiagnosticType {
	clinit_TypeValidator()
	return static_TypeValidator_MISSING_EXTENDS_TAG_WARNING
}
func StaticSet_TypeValidator_MISSING_EXTENDS_TAG_WARNING(v *DiagnosticType) {
	clinit_TypeValidator()
	static_TypeValidator_MISSING_EXTENDS_TAG_WARNING = v
}

var static_TypeValidator_DUP_VAR_DECLARATION *DiagnosticType

func StaticGet_TypeValidator_DUP_VAR_DECLARATION() *DiagnosticType {
	clinit_TypeValidator()
	return static_TypeValidator_DUP_VAR_DECLARATION
}
func StaticSet_TypeValidator_DUP_VAR_DECLARATION(v *DiagnosticType) {
	clinit_TypeValidator()
	static_TypeValidator_DUP_VAR_DECLARATION = v
}

var static_TypeValidator_DUP_VAR_DECLARATION_TYPE_MISMATCH *DiagnosticType

func StaticGet_TypeValidator_DUP_VAR_DECLARATION_TYPE_MISMATCH() *DiagnosticType {
	clinit_TypeValidator()
	return static_TypeValidator_DUP_VAR_DECLARATION_TYPE_MISMATCH
}
func StaticSet_TypeValidator_DUP_VAR_DECLARATION_TYPE_MISMATCH(v *DiagnosticType) {
	clinit_TypeValidator()
	static_TypeValidator_DUP_VAR_DECLARATION_TYPE_MISMATCH = v
}

var static_TypeValidator_INTERFACE_METHOD_NOT_IMPLEMENTED *DiagnosticType

func StaticGet_TypeValidator_INTERFACE_METHOD_NOT_IMPLEMENTED() *DiagnosticType {
	clinit_TypeValidator()
	return static_TypeValidator_INTERFACE_METHOD_NOT_IMPLEMENTED
}
func StaticSet_TypeValidator_INTERFACE_METHOD_NOT_IMPLEMENTED(v *DiagnosticType) {
	clinit_TypeValidator()
	static_TypeValidator_INTERFACE_METHOD_NOT_IMPLEMENTED = v
}

var static_TypeValidator_HIDDEN_INTERFACE_PROPERTY_MISMATCH *DiagnosticType

func StaticGet_TypeValidator_HIDDEN_INTERFACE_PROPERTY_MISMATCH() *DiagnosticType {
	clinit_TypeValidator()
	return static_TypeValidator_HIDDEN_INTERFACE_PROPERTY_MISMATCH
}
func StaticSet_TypeValidator_HIDDEN_INTERFACE_PROPERTY_MISMATCH(v *DiagnosticType) {
	clinit_TypeValidator()
	static_TypeValidator_HIDDEN_INTERFACE_PROPERTY_MISMATCH = v
}

var static_TypeValidator_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH *DiagnosticType

func StaticGet_TypeValidator_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH() *DiagnosticType {
	clinit_TypeValidator()
	return static_TypeValidator_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH
}
func StaticSet_TypeValidator_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH(v *DiagnosticType) {
	clinit_TypeValidator()
	static_TypeValidator_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH = v
}

var static_TypeValidator_ABSTRACT_METHOD_NOT_IMPLEMENTED *DiagnosticType

func StaticGet_TypeValidator_ABSTRACT_METHOD_NOT_IMPLEMENTED() *DiagnosticType {
	clinit_TypeValidator()
	return static_TypeValidator_ABSTRACT_METHOD_NOT_IMPLEMENTED
}
func StaticSet_TypeValidator_ABSTRACT_METHOD_NOT_IMPLEMENTED(v *DiagnosticType) {
	clinit_TypeValidator()
	static_TypeValidator_ABSTRACT_METHOD_NOT_IMPLEMENTED = v
}

var static_TypeValidator_UNKNOWN_TYPEOF_VALUE *DiagnosticType

func StaticGet_TypeValidator_UNKNOWN_TYPEOF_VALUE() *DiagnosticType {
	clinit_TypeValidator()
	return static_TypeValidator_UNKNOWN_TYPEOF_VALUE
}
func StaticSet_TypeValidator_UNKNOWN_TYPEOF_VALUE(v *DiagnosticType) {
	clinit_TypeValidator()
	static_TypeValidator_UNKNOWN_TYPEOF_VALUE = v
}

var static_TypeValidator_ILLEGAL_PROPERTY_ACCESS *DiagnosticType

func StaticGet_TypeValidator_ILLEGAL_PROPERTY_ACCESS() *DiagnosticType {
	clinit_TypeValidator()
	return static_TypeValidator_ILLEGAL_PROPERTY_ACCESS
}
func StaticSet_TypeValidator_ILLEGAL_PROPERTY_ACCESS(v *DiagnosticType) {
	clinit_TypeValidator()
	static_TypeValidator_ILLEGAL_PROPERTY_ACCESS = v
}

var static_TypeValidator_ALL_DIAGNOSTICS *DiagnosticGroup

func StaticGet_TypeValidator_ALL_DIAGNOSTICS() *DiagnosticGroup {
	clinit_TypeValidator()
	return static_TypeValidator_ALL_DIAGNOSTICS
}
func StaticSet_TypeValidator_ALL_DIAGNOSTICS(v *DiagnosticGroup) {
	clinit_TypeValidator()
	static_TypeValidator_ALL_DIAGNOSTICS = v
}

var clinitDone_TypeValidator uint32
var clinitState_TypeValidator j2gort.ClinitState

func clinit_TypeValidator() {
	if atomic.LoadUint32(&clinitDone_TypeValidator) != 0 {
		return
	}
	clinitSlow_TypeValidator()
}

func clinitSlow_TypeValidator() {
	if !j2gort.EnterClinit(&clinitState_TypeValidator, &clinitDone_TypeValidator) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_TypeValidator, &clinitDone_TypeValidator)
	static_TypeValidator_INVALID_CAST = DiagnosticType_Warning("JSC_INVALID_CAST", "invalid cast - must be a subtype or supertype\nfrom: {0}\nto  : {1}")
	static_TypeValidator_TYPE_MISMATCH_WARNING = DiagnosticType_Warning("JSC_TYPE_MISMATCH", "{0}")
	static_TypeValidator_INVALID_ASYNC_RETURN_TYPE = DiagnosticType_Warning("JSC_INVALID_ASYNC_RETURN_TYPE", "The return type of an async function must be a supertype of Promise\nfound: {0}")
	static_TypeValidator_INVALID_OPERAND_TYPE = DiagnosticType_Disabled("JSC_INVALID_OPERAND_TYPE", "{0}")
	static_TypeValidator_MISSING_EXTENDS_TAG_WARNING = DiagnosticType_Warning("JSC_MISSING_EXTENDS_TAG", "Missing @extends tag on type {0}")
	static_TypeValidator_DUP_VAR_DECLARATION = DiagnosticType_Warning("JSC_DUP_VAR_DECLARATION", "variable {0} redefined, original definition at {1}:{2}")
	static_TypeValidator_DUP_VAR_DECLARATION_TYPE_MISMATCH = DiagnosticType_Warning("JSC_DUP_VAR_DECLARATION_TYPE_MISMATCH", "variable {0} redefined with type {1}, original definition at {2}:{3} with type {4}")
	static_TypeValidator_INTERFACE_METHOD_NOT_IMPLEMENTED = DiagnosticType_Warning("JSC_INTERFACE_METHOD_NOT_IMPLEMENTED", "property {0} on interface {1} is not implemented by type {2}")
	static_TypeValidator_HIDDEN_INTERFACE_PROPERTY_MISMATCH = DiagnosticType_Warning("JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH", "mismatch of the {0} property on type {4} and the type of the property it overrides from interface {1}\noriginal: {2}\noverride: {3}")
	static_TypeValidator_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH = DiagnosticType_Warning("JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH", "mismatch of the {0} property type and the type of the property it overrides from superclass {1}\noriginal: {2}\noverride: {3}")
	static_TypeValidator_ABSTRACT_METHOD_NOT_IMPLEMENTED = DiagnosticType_Warning("JSC_ABSTRACT_METHOD_NOT_IMPLEMENTED", "property {0} on abstract class {1} is not implemented by type {2}")
	static_TypeValidator_UNKNOWN_TYPEOF_VALUE = DiagnosticType_Warning("JSC_UNKNOWN_TYPEOF_VALUE", "unknown type: {0}")
	static_TypeValidator_ILLEGAL_PROPERTY_ACCESS = DiagnosticType_Warning("JSC_ILLEGAL_PROPERTY_ACCESS", "Cannot do {0} access on a {1}")
	static_TypeValidator_ALL_DIAGNOSTICS = NewDiagnosticGroup_1_DiagnosticType_2([]*DiagnosticType{StaticGet_TypeValidator_ABSTRACT_METHOD_NOT_IMPLEMENTED(), StaticGet_TypeValidator_DUP_VAR_DECLARATION(), StaticGet_TypeValidator_DUP_VAR_DECLARATION_TYPE_MISMATCH(), StaticGet_TypeValidator_HIDDEN_INTERFACE_PROPERTY_MISMATCH(), StaticGet_TypeValidator_ILLEGAL_PROPERTY_ACCESS(), StaticGet_TypeValidator_INTERFACE_METHOD_NOT_IMPLEMENTED(), StaticGet_TypeValidator_INVALID_ASYNC_RETURN_TYPE(), StaticGet_TypeValidator_INVALID_CAST(), StaticGet_TypeValidator_MISSING_EXTENDS_TAG_WARNING(), StaticGet_TypeValidator_TYPE_MISMATCH_WARNING(), StaticGet_TypeValidator_UNKNOWN_TYPEOF_VALUE()})
}

func (this_ *TypeValidator) initFields_TypeValidator() {
	_ = this_
	this_.subtypingMode = pkg_rhino.StaticGet_JSType_SubtypingMode_NORMAL()
	this_.mismatches = NewTypeMismatch_Accumulator()
}

// Java constructor: TypeValidator(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewTypeValidator(compiler AbstractCompiler) *TypeValidator {
	clinit_TypeValidator()
	this_ := &TypeValidator{}
	InitTypeValidator(this_, compiler)
	return this_
}

func InitTypeValidator(this_ *TypeValidator, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_TypeValidator()
	this_.compiler = compiler
	this_.typeRegistry = compiler.GetTypeRegistry()
	this_.allBitwisableValueTypes = this_.typeRegistry.CreateUnionType_1_JSTypeNative([]pkg_rhino.JSTypeNative{pkg_rhino.StaticGet_JSTypeNative_STRING_TYPE(), pkg_rhino.StaticGet_JSTypeNative_NUMBER_TYPE(), pkg_rhino.StaticGet_JSTypeNative_BOOLEAN_TYPE(), pkg_rhino.StaticGet_JSTypeNative_NULL_TYPE(), pkg_rhino.StaticGet_JSTypeNative_VOID_TYPE()})
	this_.nullOrUndefined = this_.typeRegistry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_NULL_VOID())
	this_.promiseOfUnknownType = pkg_rhino.ToIface_JSType(this_.typeRegistry.CreateTemplatizedType_2_ObjectType_JSType(this_.typeRegistry.GetNativeObjectType(pkg_rhino.StaticGet_JSTypeNative_PROMISE_TYPE()), []pkg_rhino.JSType{this_.typeRegistry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())}))
	this_.iterableOrAsyncIterable = this_.typeRegistry.CreateUnionType_2(j2gort.AsIface[pkg_rhino.JSType](this_.typeRegistry.GetNativeObjectType(pkg_rhino.StaticGet_JSTypeNative_ITERATOR_TYPE())), j2gort.AsIface[pkg_rhino.JSType](this_.typeRegistry.GetNativeObjectType(pkg_rhino.StaticGet_JSTypeNative_ASYNC_ITERATOR_TYPE())))
}

// Java static method: TypeValidator.getInstanceOfCtor(com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType t)
func TypeValidator_GetInstanceOfCtor(t pkg_rhino.JSType) pkg_rhino.ObjectType {
	clinit_TypeValidator()
	_ = t
	if t == nil || t.IsJavaNil_() {
		return nil
	}
	var ctor *pkg_rhino.FunctionType = pkg_rhino.JSType_ToMaybeFunctionType(j2gort.AsIface[pkg_rhino.JSType](t.Dereference()))
	_ = ctor
	if (ctor != nil) && ctor.IsConstructor() {
		return ctor.GetInstanceType()
	}
	return nil
}

// Java method: TypeValidator.getMismatches()
func (this_ *TypeValidator) GetMismatches() j2gort.Ext_java_lang_Iterable {
	return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.mismatches.GetMismatches())
}

// Java method: TypeValidator.setSubtypingMode(com.google.javascript.rhino.jstype.JSType.SubtypingMode mode)
func (this_ *TypeValidator) SetSubtypingMode(mode pkg_rhino.JSType_SubtypingMode) {
	_ = mode
	this_.subtypingMode = mode
}

// Java method: TypeValidator.expectValidTypeofName(com.google.javascript.rhino.Node n, java.lang.String found)
func (this_ *TypeValidator) ExpectValidTypeofName(n pkg_rhino.Node, found string) {
	_ = n
	_ = found
	this_.compiler.Report_1(JSError_Make_3(n, StaticGet_TypeValidator_UNKNOWN_TYPEOF_VALUE(), []string{found}))
}

// Java method: TypeValidator.expectObject(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType type, java.lang.String msg)
func (this_ *TypeValidator) ExpectObject(n pkg_rhino.Node, type_ pkg_rhino.JSType, msg string) bool {
	_ = n
	_ = type_
	_ = msg
	if !(type_.MatchesObjectContext()) {
		this_.priv_TypeValidator_mismatch_5652(n, msg, type_, pkg_rhino.StaticGet_JSTypeNative_OBJECT_TYPE())
		return false
	}
	return true
}

// Java method: TypeValidator.expectActualObject(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType type, java.lang.String msg)
func (this_ *TypeValidator) ExpectActualObject(n pkg_rhino.Node, type_ pkg_rhino.JSType, msg string) {
	_ = n
	_ = type_
	_ = msg
	if !(type_.IsObject()) {
		this_.priv_TypeValidator_mismatch_5652(n, msg, type_, pkg_rhino.StaticGet_JSTypeNative_OBJECT_TYPE())
	}
}

// Java method: TypeValidator.expectAnyObject(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType type, java.lang.String msg)
func (this_ *TypeValidator) ExpectAnyObject(n pkg_rhino.Node, type_ pkg_rhino.JSType, msg string) {
	_ = n
	_ = type_
	_ = msg
	var anyObjectType pkg_rhino.JSType = this_.priv_TypeValidator_getNativeType_5657(pkg_rhino.StaticGet_JSTypeNative_NO_OBJECT_TYPE())
	_ = anyObjectType
	if !(anyObjectType.IsSubtypeOf_1(type_)) && !(type_.IsEmptyType()) {
		this_.priv_TypeValidator_mismatch_5653(n, msg, type_, anyObjectType)
	}
}

// Java method: TypeValidator.expectAutoboxesToIterable(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType type, java.lang.String msg)
func (this_ *TypeValidator) ExpectAutoboxesToIterable(n pkg_rhino.Node, type_ pkg_rhino.JSType, msg string) bool {
	_ = n
	_ = type_
	_ = msg
	if type_.IsUnionType() {
		var union *pkg_rhino.UnionType = type_.ToMaybeUnionType()
		_ = union
		{
			var i int32 = int32(0)
			_ = i
			var len_ int32 = union.GetAlternatesCount()
			_ = len_
			for ; i < len_; i = i + 1 {
				var alt pkg_rhino.JSType = union.GetAlternateAt(i)
				_ = alt
				if alt.IsBoxableScalar() {
					alt = alt.AutoboxesTo()
				} else {
					alt = alt
				}
				if !(alt.IsSubtypeOf_1(this_.priv_TypeValidator_getNativeType_5657(pkg_rhino.StaticGet_JSTypeNative_ITERABLE_TYPE()))) {
					this_.priv_TypeValidator_mismatch_5652(n, msg, type_, pkg_rhino.StaticGet_JSTypeNative_ITERABLE_TYPE())
					return false
				}
			}
		}
	} else {
		var autoboxedType pkg_rhino.JSType
		_ = autoboxedType
		if type_.IsBoxableScalar() {
			autoboxedType = type_.AutoboxesTo()
		} else {
			autoboxedType = type_
		}
		if !(autoboxedType.IsSubtypeOf_1(this_.priv_TypeValidator_getNativeType_5657(pkg_rhino.StaticGet_JSTypeNative_ITERABLE_TYPE()))) {
			this_.priv_TypeValidator_mismatch_5652(n, msg, type_, pkg_rhino.StaticGet_JSTypeNative_ITERABLE_TYPE())
			return false
		}
	}
	return true
}

// Java method: TypeValidator.expectAutoboxesToIterableOrAsyncIterable(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType type, java.lang.String msg)
func (this_ *TypeValidator) ExpectAutoboxesToIterableOrAsyncIterable(n pkg_rhino.Node, type_ pkg_rhino.JSType, msg string) *j2gort.Ext_java_util_Optional {
	_ = n
	_ = type_
	_ = msg
	var maybeBoxed *JsIterables_MaybeBoxedType = JsIterables_MaybeBoxIterableOrAsyncIterable(type_, this_.typeRegistry)
	_ = maybeBoxed
	if maybeBoxed.IsMatch() {
		return j2gort.Ext_java_util_Optional_Of(maybeBoxed.GetTemplatedType())
	}
	this_.priv_TypeValidator_mismatch_5653(n, msg, type_, this_.iterableOrAsyncIterable)
	return j2gort.Ext_java_util_Optional_Empty()
}

// Java method: TypeValidator.expectGeneratorSupertype(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType type, java.lang.String msg)
func (this_ *TypeValidator) ExpectGeneratorSupertype(n pkg_rhino.Node, type_ pkg_rhino.JSType, msg string) {
	_ = n
	_ = type_
	_ = msg
	if !(this_.priv_TypeValidator_getNativeType_5657(pkg_rhino.StaticGet_JSTypeNative_GENERATOR_TYPE()).IsSubtypeOf_1(type_)) {
		this_.priv_TypeValidator_mismatch_5652(n, msg, type_, pkg_rhino.StaticGet_JSTypeNative_GENERATOR_TYPE())
	}
}

// Java method: TypeValidator.expectAsyncGeneratorSupertype(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType type, java.lang.String msg)
func (this_ *TypeValidator) ExpectAsyncGeneratorSupertype(n pkg_rhino.Node, type_ pkg_rhino.JSType, msg string) {
	_ = n
	_ = type_
	_ = msg
	if !(this_.priv_TypeValidator_getNativeType_5657(pkg_rhino.StaticGet_JSTypeNative_ASYNC_GENERATOR_TYPE()).IsSubtypeOf_1(type_)) {
		this_.priv_TypeValidator_mismatch_5652(n, msg, type_, pkg_rhino.StaticGet_JSTypeNative_ASYNC_GENERATOR_TYPE())
	}
}

// Java method: TypeValidator.expectValidAsyncReturnType(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType type)
func (this_ *TypeValidator) ExpectValidAsyncReturnType(n pkg_rhino.Node, type_ pkg_rhino.JSType) {
	_ = n
	_ = type_
	if this_.promiseOfUnknownType.IsSubtypeOf_1(type_) {
		return
	}
	var err *JSError = JSError_Make_3(n, StaticGet_TypeValidator_INVALID_ASYNC_RETURN_TYPE(), []string{type_.String()})
	_ = err
	this_.priv_TypeValidator_registerMismatchAndReport_5655(type_, this_.promiseOfUnknownType, err)
}

// Java method: TypeValidator.expectITemplateArraySupertype(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType type, java.lang.String msg)
func (this_ *TypeValidator) ExpectITemplateArraySupertype(n pkg_rhino.Node, type_ pkg_rhino.JSType, msg string) {
	_ = n
	_ = type_
	_ = msg
	if !(this_.priv_TypeValidator_getNativeType_5657(pkg_rhino.StaticGet_JSTypeNative_I_TEMPLATE_ARRAY_TYPE()).IsSubtypeOf_1(type_)) {
		this_.priv_TypeValidator_mismatch_5652(n, msg, type_, pkg_rhino.StaticGet_JSTypeNative_I_TEMPLATE_ARRAY_TYPE())
	}
}

// Java method: TypeValidator.expectString(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType type, java.lang.String msg)
func (this_ *TypeValidator) ExpectString(n pkg_rhino.Node, type_ pkg_rhino.JSType, msg string) {
	_ = n
	_ = type_
	_ = msg
	if !(type_.MatchesStringContext()) {
		this_.priv_TypeValidator_mismatch_5652(n, msg, type_, pkg_rhino.StaticGet_JSTypeNative_STRING_TYPE())
	}
}

// Java method: TypeValidator.expectNumber(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType type, java.lang.String msg)
func (this_ *TypeValidator) ExpectNumber(n pkg_rhino.Node, type_ pkg_rhino.JSType, msg string) {
	_ = n
	_ = type_
	_ = msg
	if !(type_.MatchesNumberContext()) {
		this_.priv_TypeValidator_mismatch_5652(n, msg, type_, pkg_rhino.StaticGet_JSTypeNative_NUMBER_TYPE())
	} else {
		this_.ExpectNumberStrict(n, type_, msg)
	}
}

// Java method: TypeValidator.expectNumberStrict(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType type, java.lang.String msg)
func (this_ *TypeValidator) ExpectNumberStrict(n pkg_rhino.Node, type_ pkg_rhino.JSType, msg string) {
	_ = n
	_ = type_
	_ = msg
	if !(type_.IsSubtypeOf_1(this_.priv_TypeValidator_getNativeType_5657(pkg_rhino.StaticGet_JSTypeNative_NUMBER_TYPE()))) {
		this_.priv_TypeValidator_registerMismatchAndReport_5654(n, StaticGet_TypeValidator_INVALID_OPERAND_TYPE(), msg, type_, this_.priv_TypeValidator_getNativeType_5657(pkg_rhino.StaticGet_JSTypeNative_NUMBER_TYPE()), nil, nil)
	}
}

// Java method: TypeValidator.expectMatchingTypesStrict(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType left, com.google.javascript.rhino.jstype.JSType right, java.lang.String msg)
func (this_ *TypeValidator) ExpectMatchingTypesStrict(n pkg_rhino.Node, left pkg_rhino.JSType, right pkg_rhino.JSType, msg string) {
	_ = n
	_ = left
	_ = right
	_ = msg
	if !(left.IsSubtypeOf_1(right)) && !(right.IsSubtypeOf_1(left)) {
		this_.priv_TypeValidator_registerMismatchAndReport_5654(n, StaticGet_TypeValidator_INVALID_OPERAND_TYPE(), msg, right, left, nil, nil)
	}
}

// Java method: TypeValidator.expectBitwiseable(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType type, java.lang.String msg)
func (this_ *TypeValidator) ExpectBitwiseable(n pkg_rhino.Node, type_ pkg_rhino.JSType, msg string) {
	_ = n
	_ = type_
	_ = msg
	if !(type_.MatchesNumberContext()) && !(type_.IsSubtypeOf_1(this_.allBitwisableValueTypes)) {
		this_.priv_TypeValidator_mismatch_5653(n, msg, type_, this_.allBitwisableValueTypes)
	} else {
		this_.ExpectNumberStrict(n, type_, msg)
	}
}

// Java method: TypeValidator.expectNumberOrSymbol(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType type, java.lang.String msg)
func (this_ *TypeValidator) ExpectNumberOrSymbol(n pkg_rhino.Node, type_ pkg_rhino.JSType, msg string) {
	_ = n
	_ = type_
	_ = msg
	if !(type_.MatchesNumberContext()) && !(type_.MatchesSymbolContext()) {
		this_.priv_TypeValidator_mismatch_5652(n, msg, type_, pkg_rhino.StaticGet_JSTypeNative_NUMBER_SYMBOL())
	}
}

// Java method: TypeValidator.expectStringOrSymbol(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType type, java.lang.String msg)
func (this_ *TypeValidator) ExpectStringOrSymbol(n pkg_rhino.Node, type_ pkg_rhino.JSType, msg string) {
	_ = n
	_ = type_
	_ = msg
	if !(type_.MatchesStringContext()) && !(type_.MatchesSymbolContext()) {
		this_.priv_TypeValidator_mismatch_5652(n, msg, type_, pkg_rhino.StaticGet_JSTypeNative_STRING_SYMBOL())
	}
}

// Java method: TypeValidator.expectUnknownOrComparable(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType type, java.lang.String msg)
func (this_ *TypeValidator) ExpectUnknownOrComparable(n pkg_rhino.Node, type_ pkg_rhino.JSType, msg string) {
	_ = n
	_ = type_
	_ = msg
	if !(type_.IsSubtypeOf_1(this_.priv_TypeValidator_getNativeType_5657(pkg_rhino.StaticGet_JSTypeNative_BIGINT_NUMBER_STRING()))) && !(type_.IsSubtypeOf_1(this_.priv_TypeValidator_getNativeType_5657(pkg_rhino.StaticGet_JSTypeNative_BIGINT_NUMBER_STRING_OBJECT()))) {
		if type_.MatchesNumberContext() || type_.MatchesStringContext() {
			this_.priv_TypeValidator_registerMismatchAndReport_5654(n, StaticGet_TypeValidator_INVALID_OPERAND_TYPE(), msg, type_, this_.priv_TypeValidator_getNativeType_5657(pkg_rhino.StaticGet_JSTypeNative_BIGINT_NUMBER_STRING()), nil, nil)
		} else {
			this_.priv_TypeValidator_mismatch_5652(n, msg, type_, pkg_rhino.StaticGet_JSTypeNative_BIGINT_NUMBER_STRING())
		}
	}
}

// Java method: TypeValidator.expectStringOrNumberOrSymbol(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType type, java.lang.String msg)
func (this_ *TypeValidator) ExpectStringOrNumberOrSymbol(n pkg_rhino.Node, type_ pkg_rhino.JSType, msg string) {
	_ = n
	_ = type_
	_ = msg
	if (!(type_.MatchesNumberContext()) && !(type_.MatchesStringContext())) && !(type_.MatchesSymbolContext()) {
		this_.priv_TypeValidator_mismatch_5652(n, msg, type_, pkg_rhino.StaticGet_JSTypeNative_NUMBER_STRING_SYMBOL())
	} else {
		this_.ExpectStringOrNumberOrSymbolStrict(n, type_, msg)
	}
}

// Java method: TypeValidator.expectStringOrNumberOrSymbolStrict(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType type, java.lang.String msg)
func (this_ *TypeValidator) ExpectStringOrNumberOrSymbolStrict(n pkg_rhino.Node, type_ pkg_rhino.JSType, msg string) {
	_ = n
	_ = type_
	_ = msg
	if !(type_.IsSubtypeOf_1(this_.priv_TypeValidator_getNativeType_5657(pkg_rhino.StaticGet_JSTypeNative_NUMBER_STRING_SYMBOL()))) {
		this_.priv_TypeValidator_registerMismatchAndReport_5654(n, StaticGet_TypeValidator_INVALID_OPERAND_TYPE(), msg, type_, this_.priv_TypeValidator_getNativeType_5657(pkg_rhino.StaticGet_JSTypeNative_NUMBER_STRING_SYMBOL()), nil, nil)
	}
}

// Java method: TypeValidator.expectBigIntOrNumber(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType type, java.lang.String msg)
func (this_ *TypeValidator) ExpectBigIntOrNumber(n pkg_rhino.Node, type_ pkg_rhino.JSType, msg string) {
	_ = n
	_ = type_
	_ = msg
	if !(type_.IsSubtypeOf_1(this_.priv_TypeValidator_getNativeType_5657(pkg_rhino.StaticGet_JSTypeNative_BIGINT_NUMBER()))) && !(type_.IsSubtypeOf_1(this_.priv_TypeValidator_getNativeType_5657(pkg_rhino.StaticGet_JSTypeNative_BIGINT_NUMBER_OBJECT()))) {
		if type_.MatchesNumberContext() {
			this_.priv_TypeValidator_registerMismatchAndReport_5654(n, StaticGet_TypeValidator_INVALID_OPERAND_TYPE(), msg, type_, this_.priv_TypeValidator_getNativeType_5657(pkg_rhino.StaticGet_JSTypeNative_BIGINT_NUMBER()), nil, nil)
		} else {
			this_.priv_TypeValidator_mismatch_5652(n, msg, type_, pkg_rhino.StaticGet_JSTypeNative_BIGINT_NUMBER())
		}
	}
}

// Java method: TypeValidator.expectNotNullOrUndefined(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType type, java.lang.String msg, com.google.javascript.rhino.jstype.JSType expectedType)
func (this_ *TypeValidator) ExpectNotNullOrUndefined(t *NodeTraversal, n pkg_rhino.Node, type_ pkg_rhino.JSType, msg string, expectedType pkg_rhino.JSType) bool {
	_ = t
	_ = n
	_ = type_
	_ = msg
	_ = expectedType
	if ((!(type_.IsNoType()) && !(type_.IsUnknownType())) && type_.IsSubtypeOf_1(this_.nullOrUndefined)) && !(TypeValidator_ContainsForwardDeclaredUnresolvedName(type_)) {
		if (n.IsGetProp() && !(t.InGlobalScope())) && type_.IsNullType() {
			return true
		}
		this_.priv_TypeValidator_mismatch_5653(n, msg, type_, expectedType)
		return false
	}
	return true
}

// Java static method: TypeValidator.containsForwardDeclaredUnresolvedName(com.google.javascript.rhino.jstype.JSType type)
func TypeValidator_ContainsForwardDeclaredUnresolvedName(type_ pkg_rhino.JSType) bool {
	clinit_TypeValidator()
	_ = type_
	if type_.IsUnionType() {
		var union *pkg_rhino.UnionType = type_.ToMaybeUnionType()
		_ = union
		{
			var i int32 = int32(0)
			_ = i
			var len_ int32 = union.GetAlternatesCount()
			_ = len_
			for ; i < len_; i = i + 1 {
				if TypeValidator_ContainsForwardDeclaredUnresolvedName(union.GetAlternateAt(i)) {
					return true
				}
			}
		}
	}
	return type_.IsNoResolvedType()
}

// Java method: TypeValidator.expectSwitchMatchesCase(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType switchType, com.google.javascript.rhino.jstype.JSType caseType)
func (this_ *TypeValidator) ExpectSwitchMatchesCase(n pkg_rhino.Node, switchType pkg_rhino.JSType, caseType pkg_rhino.JSType) {
	_ = n
	_ = switchType
	_ = caseType
	if !(switchType.CanTestForShallowEqualityWith(caseType)) {
		this_.priv_TypeValidator_mismatch_5653(n.GetFirstChild(), "case expression doesn't match switch", caseType, switchType)
	}
}

// Java method: TypeValidator.expectIndexMatch(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType objType, com.google.javascript.rhino.jstype.JSType indexType)
func (this_ *TypeValidator) ExpectIndexMatch(n pkg_rhino.Node, objType pkg_rhino.JSType, indexType pkg_rhino.JSType) {
	_ = n
	_ = objType
	_ = indexType
	if !(((n.IsGetElem() || n.IsOptChainGetElem()) || n.IsComputedProp()) || n.IsComputedFieldDef()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
	}
	var indexNode pkg_rhino.Node
	_ = indexNode
	if n.IsGetElem() {
		indexNode = n.GetLastChild()
	} else {
		indexNode = n.GetFirstChild()
	}
	if indexType.IsSymbolValueType() {
		return
	}
	if objType.IsUnknownType() {
		this_.ExpectStringOrNumberOrSymbol(indexNode, indexType, "property access")
		return
	}
	var dereferenced pkg_rhino.ObjectType = objType.Dereference()
	_ = dereferenced
	if (dereferenced != nil && !dereferenced.IsJavaNil_()) && dereferenced.GetTemplateTypeMap().HasTemplateKey_1(this_.typeRegistry.GetObjectIndexKey()) {
		var templateType pkg_rhino.JSType = dereferenced.GetTemplateTypeMap().GetResolvedTemplateType(this_.typeRegistry.GetObjectIndexKey())
		_ = templateType
		this_.ExpectCanAssignTo(indexNode, indexType, templateType, "restricted index type")
	} else if (dereferenced != nil && !dereferenced.IsJavaNil_()) && dereferenced.IsArrayType() {
		this_.ExpectNumberOrSymbol(indexNode, indexType, "array access")
	} else if objType.IsStruct() {
		this_.compiler.Report_1(JSError_Make_3(indexNode, StaticGet_TypeValidator_ILLEGAL_PROPERTY_ACCESS(), []string{"'[]'", "struct"}))
	} else if objType.MatchesObjectContext() {
		this_.ExpectStringOrSymbol(indexNode, indexType, "property access")
	} else {
		this_.priv_TypeValidator_mismatch_5653(n, "only arrays or objects can be accessed", objType, this_.typeRegistry.CreateUnionType_1_JSTypeNative([]pkg_rhino.JSTypeNative{pkg_rhino.StaticGet_JSTypeNative_ARRAY_TYPE(), pkg_rhino.StaticGet_JSTypeNative_OBJECT_TYPE()}))
	}
}

// Java method: TypeValidator.expectCanAssignToPropertyOf(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType rightType, com.google.javascript.rhino.jstype.JSType leftType, com.google.javascript.rhino.Node owner, java.lang.String propName)
func (this_ *TypeValidator) ExpectCanAssignToPropertyOf_5(n pkg_rhino.Node, rightType pkg_rhino.JSType, leftType pkg_rhino.JSType, owner pkg_rhino.Node, propName string) bool {
	_ = n
	_ = rightType
	_ = leftType
	_ = owner
	_ = propName
	var typeNameSupplier j2gort.Ext_com_google_common_base_Supplier = nil
	_ = typeNameSupplier
	if n.IsMemberFieldDef() {
		var classType *pkg_rhino.FunctionType = n.GetGrandparent().GetJSType().AssertFunctionType()
		_ = classType
		typeNameSupplier = j2gort.NewLambda_Ext_com_google_common_base_Supplier(func() any {
			return j2gort.BoxString(classType.GetInstanceType().String())
		})
	} else {
		typeNameSupplier = j2gort.NewLambda_Ext_com_google_common_base_Supplier(func() any {
			return j2gort.BoxString(this_.typeRegistry.GetReadableTypeName(owner))
		})
	}
	return this_.ExpectCanAssignToPropertyOf_6(n, rightType, leftType, this_.priv_TypeValidator_getJSType_5656(owner), typeNameSupplier, propName)
}

// Java method: TypeValidator.expectCanAssignToPropertyOf(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType rightType, com.google.javascript.rhino.jstype.JSType leftType, com.google.javascript.rhino.jstype.JSType ownerType, com.google.common.base.Supplier<java.lang.String> typeNameSupplier, java.lang.String propName)
func (this_ *TypeValidator) ExpectCanAssignToPropertyOf_6(n pkg_rhino.Node, rightType pkg_rhino.JSType, leftType pkg_rhino.JSType, ownerType pkg_rhino.JSType, typeNameSupplier j2gort.Ext_com_google_common_base_Supplier, propName string) bool {
	_ = n
	_ = rightType
	_ = leftType
	_ = ownerType
	_ = typeNameSupplier
	_ = propName
	if leftType.IsTemplateType() {
		var left *pkg_rhino.TemplateType = leftType.ToMaybeTemplateType()
		_ = left
		if (rightType.ContainsReferenceAncestor(pkg_rhino.ToIface_JSType(left)) || rightType.IsUnknownType()) || left.IsUnknownType() {
			return true
		} else {
			this_.priv_TypeValidator_registerMismatchAndReport_5654(n, StaticGet_TypeValidator_TYPE_MISMATCH_WARNING(), ("assignment to property " + j2gort.StrValString(propName) + " of " + j2gort.StrValString(j2gort.CastString(typeNameSupplier.Get()))), rightType, leftType, j2gort.NewExt_java_util_LinkedHashSet(), j2gort.NewExt_java_util_LinkedHashSet())
			return false
		}
	}
	if !(leftType.IsNoType()) && !(rightType.IsSubtypeOf_1(leftType)) {
		if ownerType.IsFunctionPrototypeType() {
			var ownerFn *pkg_rhino.FunctionType = ownerType.ToObjectType().GetOwnerFunction()
			_ = ownerFn
			if (ownerFn.IsInterface() && rightType.IsFunctionType()) && leftType.IsFunctionType() {
				return true
			}
		}
		this_.priv_TypeValidator_mismatch_5653(n, ("assignment to property " + j2gort.StrValString(propName) + " of " + j2gort.StrValString(j2gort.CastString(typeNameSupplier.Get()))), rightType, leftType)
		return false
	}
	return true
}

// Java method: TypeValidator.expectCanAssignTo(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType rightType, com.google.javascript.rhino.jstype.JSType leftType, java.lang.String msg)
func (this_ *TypeValidator) ExpectCanAssignTo(n pkg_rhino.Node, rightType pkg_rhino.JSType, leftType pkg_rhino.JSType, msg string) bool {
	_ = n
	_ = rightType
	_ = leftType
	_ = msg
	if leftType.IsTemplateType() {
		var left *pkg_rhino.TemplateType = leftType.ToMaybeTemplateType()
		_ = left
		if (rightType.ContainsReferenceAncestor(pkg_rhino.ToIface_JSType(left)) || rightType.IsUnknownType()) || left.IsUnknownType() {
			return true
		} else {
			this_.priv_TypeValidator_registerMismatchAndReport_5654(n, StaticGet_TypeValidator_TYPE_MISMATCH_WARNING(), msg, rightType, leftType, j2gort.NewExt_java_util_LinkedHashSet(), j2gort.NewExt_java_util_LinkedHashSet())
			return false
		}
	}
	if !(rightType.IsSubtypeOf_1(leftType)) {
		this_.priv_TypeValidator_mismatch_5653(n, msg, rightType, leftType)
		return false
	}
	return true
}

// Java method: TypeValidator.expectArgumentMatchesParameter(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType argType, com.google.javascript.rhino.jstype.JSType paramType, com.google.javascript.rhino.Node callNode, int ordinal)
func (this_ *TypeValidator) ExpectArgumentMatchesParameter(n pkg_rhino.Node, argType pkg_rhino.JSType, paramType pkg_rhino.JSType, callNode pkg_rhino.Node, ordinal int32) {
	_ = n
	_ = argType
	_ = paramType
	_ = callNode
	_ = ordinal
	if !(argType.IsSubtypeOf_1(paramType)) {
		this_.priv_TypeValidator_mismatch_5653(n, ("actual parameter " + j2gort.StrValInt32(ordinal) + " of " + j2gort.StrValString(this_.typeRegistry.GetReadableTypeNameNoDeref(callNode.GetFirstChild())) + " does not match formal parameter"), argType, paramType)
	}
}

// Java method: TypeValidator.expectSuperType(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.ObjectType superObject, com.google.javascript.rhino.jstype.ObjectType subObject)
func (this_ *TypeValidator) ExpectSuperType(n pkg_rhino.Node, superObject pkg_rhino.ObjectType, subObject pkg_rhino.ObjectType) {
	_ = n
	_ = superObject
	_ = subObject
	var subCtor *pkg_rhino.FunctionType = subObject.GetConstructor()
	_ = subCtor
	var implicitProto pkg_rhino.ObjectType = subObject.GetImplicitPrototype()
	_ = implicitProto
	var declaredSuper pkg_rhino.ObjectType
	_ = declaredSuper
	if implicitProto == nil || implicitProto.IsJavaNil_() {
		declaredSuper = nil
	} else {
		declaredSuper = implicitProto.GetImplicitPrototype()
	}
	if (declaredSuper != nil && !declaredSuper.IsJavaNil_()) && declaredSuper.IsTemplatizedType() {
		declaredSuper = declaredSuper.ToMaybeTemplatizedType().GetReferencedType()
	}
	if ((declaredSuper != nil && !declaredSuper.IsJavaNil_()) && !(j2gort.InstanceOf[*pkg_rhino.UnknownType](superObject))) && !(declaredSuper.Equals(superObject)) {
		if declaredSuper.Equals(this_.priv_TypeValidator_getNativeType_5657(pkg_rhino.StaticGet_JSTypeNative_OBJECT_TYPE())) {
			this_.priv_TypeValidator_registerMismatchAndReport_5655(j2gort.AsIface[pkg_rhino.JSType](superObject), j2gort.AsIface[pkg_rhino.JSType](declaredSuper), JSError_Make_3(n, StaticGet_TypeValidator_MISSING_EXTENDS_TAG_WARNING(), []string{subObject.String()}))
		} else {
			this_.priv_TypeValidator_mismatch_5653(n, "mismatch in declaration of superclass type", j2gort.AsIface[pkg_rhino.JSType](superObject), j2gort.AsIface[pkg_rhino.JSType](declaredSuper))
		}
		if !(subCtor.HasCachedValues()) {
			subCtor.SetPrototypeBasedOn_1(superObject)
		}
	}
}

// Java method: TypeValidator.expectExtends(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.FunctionType subCtor, com.google.javascript.rhino.jstype.FunctionType astSuperCtor)
func (this_ *TypeValidator) ExpectExtends(n pkg_rhino.Node, subCtor *pkg_rhino.FunctionType, astSuperCtor *pkg_rhino.FunctionType) {
	_ = n
	_ = subCtor
	_ = astSuperCtor
	if (astSuperCtor == nil) || (!(astSuperCtor.IsConstructor()) && !(astSuperCtor.IsInterface())) {
		return
	}
	if astSuperCtor.IsConstructor() != subCtor.IsConstructor() {
		return
	}
	var astSuperInstance pkg_rhino.ObjectType = astSuperCtor.GetInstanceType()
	_ = astSuperInstance
	if subCtor.IsConstructor() {
		var registeredSuperCtor *pkg_rhino.FunctionType = subCtor.GetSuperClassConstructor()
		_ = registeredSuperCtor
		if registeredSuperCtor != nil {
			var registeredSuperInstance pkg_rhino.ObjectType = registeredSuperCtor.GetInstanceType()
			_ = registeredSuperInstance
			if !(astSuperInstance.Equals(registeredSuperInstance)) {
				this_.priv_TypeValidator_mismatch_5653(n, "mismatch in declaration of superclass type", j2gort.AsIface[pkg_rhino.JSType](astSuperInstance), j2gort.AsIface[pkg_rhino.JSType](registeredSuperInstance))
			}
		}
	} else if subCtor.IsInterface() {
	}
}

// Java method: TypeValidator.expectCanAssignToPrototype(com.google.javascript.rhino.jstype.JSType ownerType, com.google.javascript.rhino.Node node, com.google.javascript.rhino.jstype.JSType rightType)
func (this_ *TypeValidator) ExpectCanAssignToPrototype(ownerType pkg_rhino.JSType, node pkg_rhino.Node, rightType pkg_rhino.JSType) {
	_ = ownerType
	_ = node
	_ = rightType
	if ownerType.IsFunctionType() {
		var functionType *pkg_rhino.FunctionType = ownerType.ToMaybeFunctionType()
		_ = functionType
		if functionType.IsConstructor() {
			this_.ExpectObject(node, rightType, "cannot override prototype with non-object")
		}
	}
}

// Java method: TypeValidator.expectCanCast(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType targetType, com.google.javascript.rhino.jstype.JSType sourceType)
func (this_ *TypeValidator) ExpectCanCast(n pkg_rhino.Node, targetType pkg_rhino.JSType, sourceType pkg_rhino.JSType) {
	_ = n
	_ = targetType
	_ = sourceType
	if !(sourceType.CanCastTo(targetType)) {
		this_.priv_TypeValidator_registerMismatchAndReport_5655(sourceType, targetType, JSError_Make_3(n, StaticGet_TypeValidator_INVALID_CAST(), []string{sourceType.String(), targetType.String()}))
	} else if !(sourceType.IsSubtypeOf_1(targetType)) {
		this_.mismatches.RegisterMismatch(n, sourceType, targetType)
	}
}

// Java method: TypeValidator.expectUndeclaredVariable(java.lang.String sourceName, com.google.javascript.jscomp.CompilerInput input, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent, com.google.javascript.jscomp.TypedVar var, java.lang.String variableName, com.google.javascript.rhino.jstype.JSType newType)
func (this_ *TypeValidator) ExpectUndeclaredVariable(sourceName string, input *CompilerInput, n pkg_rhino.Node, parent pkg_rhino.Node, var_ *TypedVar, variableName string, newType pkg_rhino.JSType) *TypedVar {
	_ = sourceName
	_ = input
	_ = n
	_ = parent
	_ = var_
	_ = variableName
	_ = newType
	var newVar *TypedVar = var_
	_ = newVar
	var varType pkg_rhino.JSType = var_.GetType()
	_ = varType
	if (((varType != nil && !varType.IsJavaNil_()) && !(j2gort.IfaceEquals[pkg_rhino.JSType](varType, this_.typeRegistry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())))) && (newType != nil && !newType.IsJavaNil_())) && !(j2gort.IfaceEquals[pkg_rhino.JSType](newType, this_.typeRegistry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE()))) {
		if var_.GetInput() == nil {
			var s *TypedScope = j2gort.CastPtr[TypedScope](var_.GetScope())
			_ = s
			s.Undeclare(ToIface_AbstractVar(var_))
			newVar = s.Declare(variableName, n, varType, input, false)
			n.SetJSType(varType)
			if parent.IsVar() {
				if n.HasChildren() {
					n.GetFirstChild().SetJSType(varType)
				}
			} else if parent.IsExprResult() {
				n.SetJSType(varType)
			} else {
				if !(parent.IsFunction() || parent.IsClass()) {
					j2gort.CheckState(false)
				}
				parent.SetJSType(varType)
			}
		} else {
			var allowDupe bool = TypeValidator_HasDuplicateDeclarationSuppression(this_.compiler, var_.GetNameNode())
			_ = allowDupe
			if !(allowDupe) {
				if !(newType.Equals(varType)) {
					this_.compiler.Report_1(JSError_Make_3(n, StaticGet_TypeValidator_DUP_VAR_DECLARATION_TYPE_MISMATCH(), []string{variableName, newType.String(), var_.GetInputName(), strconv.FormatInt(int64(var_.GetNameNode().GetLineno()), 10), varType.String()}))
				} else if !(var_.GetParentNode().IsExprResult()) {
					this_.compiler.Report_1(JSError_Make_3(n, StaticGet_TypeValidator_DUP_VAR_DECLARATION(), []string{variableName, var_.GetInputName(), strconv.FormatInt(int64(var_.GetNameNode().GetLineno()), 10)}))
				}
			}
		}
	}
	return newVar
}

// Java method: TypeValidator.expectAllInterfaceProperties(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.FunctionType type)
func (this_ *TypeValidator) ExpectAllInterfaceProperties(n pkg_rhino.Node, type_ *pkg_rhino.FunctionType) {
	_ = n
	_ = type_
	var instance pkg_rhino.ObjectType = type_.GetInstanceType()
	_ = instance
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(type_.GetAllImplementedInterfaces()))
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			implemented := j2gort.AsIface[pkg_rhino.ObjectType](rawElem_2)
			_ = implemented
			this_.priv_TypeValidator_expectInterfaceProperties_5647(n, instance, implemented)
		}
	}
	{
		cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(type_.GetExtendedInterfaces()))
		for cur_3.HasNext() {
			rawElem_4 := cur_3.Next()
			extended := j2gort.AsIface[pkg_rhino.ObjectType](rawElem_4)
			_ = extended
			this_.priv_TypeValidator_expectInterfaceProperties_5647(n, instance, extended)
		}
	}
}

// Java method: TypeValidator.expectInterfaceProperties(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.ObjectType instance, com.google.javascript.rhino.jstype.ObjectType ancestorInterface)
func (this_ *TypeValidator) priv_TypeValidator_expectInterfaceProperties_5647(n pkg_rhino.Node, instance pkg_rhino.ObjectType, ancestorInterface pkg_rhino.ObjectType) {
	_ = n
	_ = instance
	_ = ancestorInterface
	{
		cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(ancestorInterface.GetOwnPropertyKeys()))
		for cur_5.HasNext() {
			rawElem_6 := cur_5.Next()
			prop := j2gort.AsIface[pkg_rhino.Property_Key](rawElem_6)
			_ = prop
			this_.priv_TypeValidator_expectInterfaceProperty_5648(n, instance, ancestorInterface, prop)
		}
	}
	if !(j2gort.IfaceEquals[pkg_rhino.ObjectType](ancestorInterface.GetImplicitPrototype(), nil)) {
		{
			cur_7 := j2gort.NewIterCursor(j2gort.CheckNotNull(ancestorInterface.GetImplicitPrototype().GetOwnPropertyKeys()))
			for cur_7.HasNext() {
				rawElem_8 := cur_7.Next()
				prop_2 := j2gort.AsIface[pkg_rhino.Property_Key](rawElem_8)
				_ = prop_2
				this_.priv_TypeValidator_expectInterfaceProperty_5648(n, instance, ancestorInterface, prop_2)
			}
		}
	}
}

// Java method: TypeValidator.expectInterfaceProperty(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.ObjectType instance, com.google.javascript.rhino.jstype.ObjectType implementedInterface, com.google.javascript.rhino.jstype.Property.Key propName)
func (this_ *TypeValidator) priv_TypeValidator_expectInterfaceProperty_5648(n pkg_rhino.Node, instance pkg_rhino.ObjectType, implementedInterface pkg_rhino.ObjectType, propName pkg_rhino.Property_Key) {
	_ = n
	_ = instance
	_ = implementedInterface
	_ = propName
	var propSlot *pkg_rhino.Property_OwnedProperty = instance.FindClosestDefinition_1_Property_Key(propName)
	_ = propSlot
	if (propSlot == nil) || (!(instance.GetConstructor().IsInterface()) && propSlot.IsOwnedByInterface()) {
		if instance.GetConstructor().IsAbstract() || instance.GetConstructor().IsInterface() {
			return
		}
		if implementedInterface.GetPropertyType_1_Property_Key(propName).IsVoidable() {
			return
		}
		this_.priv_TypeValidator_registerMismatchAndReport_5655(j2gort.AsIface[pkg_rhino.JSType](instance), j2gort.AsIface[pkg_rhino.JSType](implementedInterface), JSError_Make_3(n, StaticGet_TypeValidator_INTERFACE_METHOD_NOT_IMPLEMENTED(), []string{propName.HumanReadableName(), implementedInterface.GetReferenceName(), instance.String()}))
	} else {
		var local bool = propSlot.GetOwnerInstanceType().Equals(instance)
		_ = local
		if !(local) && instance.GetConstructor().IsInterface() {
			return
		}
		var prop *pkg_rhino.Property = propSlot.GetValue()
		_ = prop
		var propNode pkg_rhino.Node
		_ = propNode
		if j2gort.CastPtr[pkg_rhino.Property](prop.GetDeclaration()) == nil {
			propNode = pkg_rhino.Node(0)
		} else {
			propNode = j2gort.CastPtr[pkg_rhino.Property](prop.GetDeclaration()).GetNode()
		}
		if propNode == 0 {
			propNode = n
		} else {
			propNode = propNode
		}
		this_.CheckPropertyType_5_Node_JSType_ObjectType_Property_Key_JSType(propNode, j2gort.AsIface[pkg_rhino.JSType](instance), implementedInterface, propName, prop.GetType())
	}
}

// Java method: TypeValidator.checkPropertyType(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType instance, com.google.javascript.rhino.jstype.ObjectType parent, java.lang.String propertyName, com.google.javascript.rhino.jstype.JSType found)
func (this_ *TypeValidator) CheckPropertyType_5_Node_JSType_ObjectType_string_JSType(n pkg_rhino.Node, instance pkg_rhino.JSType, parent pkg_rhino.ObjectType, propertyName string, found pkg_rhino.JSType) {
	_ = n
	_ = instance
	_ = parent
	_ = propertyName
	_ = found
	var required pkg_rhino.JSType = parent.GetPropertyType_1_string(propertyName)
	_ = required
	var typeMap *pkg_rhino.TemplateTypeMap = instance.GetTemplateTypeMap()
	_ = typeMap
	if !(typeMap.IsEmpty()) && required.HasAnyTemplateTypes() {
		required = this_.typeRegistry.ReplaceTemplateTypesPartial(required, typeMap)
	}
	if found.IsSubtypeOf_2(required, this_.subtypingMode) {
		return
	}
	var err *JSError = JSError_Make_3(n, func() *DiagnosticType {
		if parent.GetConstructor().IsInterface() {
			return StaticGet_TypeValidator_HIDDEN_INTERFACE_PROPERTY_MISMATCH()
		}
		return StaticGet_TypeValidator_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH()
	}(), []string{propertyName, parent.GetReferenceName(), required.String(), found.String(), instance.String()})
	_ = err
	this_.priv_TypeValidator_registerMismatchAndReport_5655(found, required, err)
}

// Java method: TypeValidator.checkPropertyType(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType instance, com.google.javascript.rhino.jstype.ObjectType parent, com.google.javascript.rhino.jstype.Property.Key propertyName, com.google.javascript.rhino.jstype.JSType found)
func (this_ *TypeValidator) CheckPropertyType_5_Node_JSType_ObjectType_Property_Key_JSType(n pkg_rhino.Node, instance pkg_rhino.JSType, parent pkg_rhino.ObjectType, propertyName pkg_rhino.Property_Key, found pkg_rhino.JSType) {
	_ = n
	_ = instance
	_ = parent
	_ = propertyName
	_ = found
	var required pkg_rhino.JSType = parent.GetPropertyType_1_Property_Key(propertyName)
	_ = required
	var typeMap *pkg_rhino.TemplateTypeMap = instance.GetTemplateTypeMap()
	_ = typeMap
	if !(typeMap.IsEmpty()) && required.HasAnyTemplateTypes() {
		required = this_.typeRegistry.ReplaceTemplateTypesPartial(required, typeMap)
	}
	if found.IsSubtypeOf_2(required, this_.subtypingMode) {
		return
	}
	var err *JSError = JSError_Make_3(n, func() *DiagnosticType {
		if parent.GetConstructor().IsInterface() {
			return StaticGet_TypeValidator_HIDDEN_INTERFACE_PROPERTY_MISMATCH()
		}
		return StaticGet_TypeValidator_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH()
	}(), []string{propertyName.HumanReadableName(), parent.GetReferenceName(), required.String(), found.String(), instance.String()})
	_ = err
	this_.priv_TypeValidator_registerMismatchAndReport_5655(found, required, err)
}

// Java method: TypeValidator.expectAbstractMethodsImplemented(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.FunctionType ctorType)
func (this_ *TypeValidator) ExpectAbstractMethodsImplemented(n pkg_rhino.Node, ctorType *pkg_rhino.FunctionType) {
	_ = n
	_ = ctorType
	if !(ctorType.IsConstructor()) {
		j2gort.CheckArgument(false)
	}
	var abstractMethodSuperTypeMap j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()
	_ = abstractMethodSuperTypeMap
	var currSuperCtor *pkg_rhino.FunctionType = ctorType.GetSuperClassConstructor()
	_ = currSuperCtor
	if (currSuperCtor == nil) || !(currSuperCtor.IsAbstract()) {
		return
	}
	for (currSuperCtor != nil) && currSuperCtor.IsAbstract() {
		var superType pkg_rhino.ObjectType = currSuperCtor.GetInstanceType()
		_ = superType
		{
			cur_9 := j2gort.NewIterCursor(j2gort.CheckNotNull(currSuperCtor.GetPrototype().GetOwnPropertyKeys()))
			for cur_9.HasNext() {
				rawElem_10 := cur_9.Next()
				prop := j2gort.AsIface[pkg_rhino.Property_Key](rawElem_10)
				_ = prop
				var maybeAbstractMethod *pkg_rhino.FunctionType = superType.FindPropertyType_1_Property_Key(prop).ToMaybeFunctionType()
				_ = maybeAbstractMethod
				if ((maybeAbstractMethod != nil) && maybeAbstractMethod.IsAbstract()) && !(abstractMethodSuperTypeMap.ContainsKey(prop)) {
					abstractMethodSuperTypeMap.Put_2(prop, superType)
				}
			}
		}
		currSuperCtor = currSuperCtor.GetSuperClassConstructor()
	}
	var instance pkg_rhino.ObjectType = ctorType.GetInstanceType()
	_ = instance
	{
		cur_11 := j2gort.NewMapEntryCursor(abstractMethodSuperTypeMap)
		for cur_11.NextEntry() {
			var method pkg_rhino.Property_Key = j2gort.AsIface[pkg_rhino.Property_Key](cur_11.GetKey())
			_ = method
			var superType_2 pkg_rhino.ObjectType = j2gort.AsIface[pkg_rhino.ObjectType](cur_11.GetValue())
			_ = superType_2
			var abstractMethod *pkg_rhino.FunctionType = instance.FindPropertyType_1_Property_Key(method).ToMaybeFunctionType()
			_ = abstractMethod
			if (abstractMethod == nil) || abstractMethod.IsAbstract() {
				this_.priv_TypeValidator_registerMismatchAndReport_5655(j2gort.AsIface[pkg_rhino.JSType](instance), j2gort.AsIface[pkg_rhino.JSType](superType_2), JSError_Make_3(n, StaticGet_TypeValidator_ABSTRACT_METHOD_NOT_IMPLEMENTED(), []string{method.HumanReadableName(), superType_2.String(), instance.String()}))
			}
		}
	}
}

// Java method: TypeValidator.mismatch(com.google.javascript.rhino.Node n, java.lang.String msg, com.google.javascript.rhino.jstype.JSType found, com.google.javascript.rhino.jstype.JSTypeNative required)
func (this_ *TypeValidator) priv_TypeValidator_mismatch_5652(n pkg_rhino.Node, msg string, found pkg_rhino.JSType, required pkg_rhino.JSTypeNative) {
	_ = n
	_ = msg
	_ = found
	_ = required
	this_.priv_TypeValidator_mismatch_5653(n, msg, found, this_.priv_TypeValidator_getNativeType_5657(required))
}

// Java method: TypeValidator.mismatch(com.google.javascript.rhino.Node n, java.lang.String msg, com.google.javascript.rhino.jstype.JSType found, com.google.javascript.rhino.jstype.JSType required)
func (this_ *TypeValidator) priv_TypeValidator_mismatch_5653(n pkg_rhino.Node, msg string, found pkg_rhino.JSType, required pkg_rhino.JSType) {
	_ = n
	_ = msg
	_ = found
	_ = required
	if !(found.IsSubtypeOf_2(required, this_.subtypingMode)) {
		var missing j2gort.Ext_java_util_Set = nil
		_ = missing
		var mismatch j2gort.Ext_java_util_Set = nil
		_ = mismatch
		var structuralType pkg_rhino.JSType = nil
		_ = structuralType
		if required.IsStructuralType() {
			structuralType = required
		} else if required.IsUnionType() {
			var restrictedType pkg_rhino.JSType = required.RestrictByNotNullOrUndefined()
			_ = restrictedType
			if restrictedType.IsStructuralType() {
				structuralType = restrictedType
			}
		}
		if structuralType != nil && !structuralType.IsJavaNil_() {
			missing = j2gort.NewExt_java_util_TreeSet()
			mismatch = j2gort.NewExt_java_util_TreeSet()
			var requiredObject pkg_rhino.ObjectType = structuralType.ToMaybeObjectType()
			_ = requiredObject
			var foundObject pkg_rhino.ObjectType = found.ToMaybeObjectType()
			_ = foundObject
			if (requiredObject != nil && !requiredObject.IsJavaNil_()) && (foundObject != nil && !foundObject.IsJavaNil_()) {
				{
					cur_12 := j2gort.NewIterCursor(j2gort.CheckNotNull(requiredObject.GetPropertyNames()))
					for cur_12.HasNext() {
						property := cur_12.NextStr()
						_ = property
						var propRequired pkg_rhino.JSType = requiredObject.GetPropertyType_1_string(property)
						_ = propRequired
						var hasProperty bool = foundObject.HasProperty_1_string(property)
						_ = hasProperty
						if !(propRequired.IsExplicitlyVoidable()) || hasProperty {
							if hasProperty {
								if !(foundObject.GetPropertyType_1_string(property).IsSubtypeOf_2(propRequired, this_.subtypingMode)) {
									j2gort.SetAddStr(mismatch, property)
								}
							} else {
								j2gort.SetAddStr(missing, property)
							}
						}
					}
				}
			}
		}
		this_.priv_TypeValidator_registerMismatchAndReport_5654(n, StaticGet_TypeValidator_TYPE_MISMATCH_WARNING(), msg, found, required, missing, mismatch)
	}
}

// Java method: TypeValidator.registerMismatchAndReport(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.DiagnosticType diagnostic, java.lang.String msg, com.google.javascript.rhino.jstype.JSType found, com.google.javascript.rhino.jstype.JSType required, java.util.@org.jspecify.annotations.Nullable Set<java.lang.String> missing, java.util.@org.jspecify.annotations.Nullable Set<java.lang.String> mismatch)
func (this_ *TypeValidator) priv_TypeValidator_registerMismatchAndReport_5654(n pkg_rhino.Node, diagnostic *DiagnosticType, msg string, found pkg_rhino.JSType, required pkg_rhino.JSType, missing j2gort.Ext_java_util_Set, mismatch j2gort.Ext_java_util_Set) {
	_ = n
	_ = diagnostic
	_ = msg
	_ = found
	_ = required
	_ = missing
	_ = mismatch
	var foundRequiredFormatted string = TypeValidator_FormatFoundRequired(msg, found, required, missing, mismatch)
	_ = foundRequiredFormatted
	var err *JSError = JSError_Make_3(n, diagnostic, []string{foundRequiredFormatted})
	_ = err
	this_.priv_TypeValidator_registerMismatchAndReport_5655(found, required, err)
}

// Java method: TypeValidator.registerMismatchAndReport(com.google.javascript.rhino.jstype.JSType found, com.google.javascript.rhino.jstype.JSType required, com.google.javascript.jscomp.JSError error)
func (this_ *TypeValidator) priv_TypeValidator_registerMismatchAndReport_5655(found pkg_rhino.JSType, required pkg_rhino.JSType, error_ *JSError) {
	_ = found
	_ = required
	_ = error_
	this_.compiler.Report_1(error_)
	this_.mismatches.RegisterMismatch(error_.Node(), found, required)
}

// Java static method: TypeValidator.formatNodeLocation(com.google.javascript.rhino.Node node)
func TypeValidator_FormatNodeLocation(node pkg_rhino.Node) string {
	clinit_TypeValidator()
	_ = node
	return j2gort.Ext_java_lang_String_Format_2("%s:%s:%s", []any{j2gort.BoxString(node.GetSourceFileName()), j2gort.ValueOfInt(node.GetLineno()), j2gort.ValueOfInt(node.GetCharno())})
}

// Java static method: TypeValidator.formatFoundRequired(java.lang.String description, com.google.javascript.rhino.jstype.JSType found, com.google.javascript.rhino.jstype.JSType required, java.util.Set<java.lang.String> missing, java.util.Set<java.lang.String> mismatch)
func TypeValidator_FormatFoundRequired(description string, found pkg_rhino.JSType, required pkg_rhino.JSType, missing j2gort.Ext_java_util_Set, mismatch j2gort.Ext_java_util_Set) string {
	clinit_TypeValidator()
	_ = description
	_ = found
	_ = required
	_ = missing
	_ = mismatch
	var foundStr string = found.String()
	_ = foundStr
	var requiredStr string = required.String()
	_ = requiredStr
	if j2gort.CheckNotNullString(foundStr) == requiredStr {
		foundStr = found.ToAnnotationString(pkg_rhino.StaticGet_JSType_Nullability_IMPLICIT())
		requiredStr = required.ToAnnotationString(pkg_rhino.StaticGet_JSType_Nullability_IMPLICIT())
	}
	if j2gort.CheckNotNullString(foundStr) == requiredStr {
		if found.IsEnumElementType() {
			foundStr = (j2gort.StrValString(foundStr) + " (enum definition: " + j2gort.StrValString(TypeValidator_FormatNodeLocation(found.ToMaybeEnumElementType().GetEnumType().GetSource())) + ")")
		}
		if required.IsEnumElementType() {
			requiredStr = (j2gort.StrValString(requiredStr) + " (enum definition: " + j2gort.StrValString(TypeValidator_FormatNodeLocation(required.ToMaybeEnumElementType().GetEnumType().GetSource())) + ")")
		}
	}
	var missingStr string = ""
	_ = missingStr
	var mismatchStr string = ""
	_ = mismatchStr
	if (missing != nil) && !(missing.IsEmpty()) {
		missingStr = j2gort.Ext_com_google_common_base_Joiner_On_1_string(",").Join_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](missing))
	}
	if (mismatch != nil) && !(mismatch.IsEmpty()) {
		mismatchStr = j2gort.Ext_com_google_common_base_Joiner_On_1_string(",").Join_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](mismatch))
	}
	if (j2gort.JStrLen(missingStr) > int32(0)) || (j2gort.JStrLen(mismatchStr) > int32(0)) {
		return Platform_FormatMessage("{0}\nfound   : {1}\nrequired: {2}\nmissing : [{3}]\nmismatch: [{4}]", []string{description, foundStr, requiredStr, missingStr, mismatchStr})
	} else {
		return Platform_FormatMessage("{0}\nfound   : {1}\nrequired: {2}", []string{description, foundStr, requiredStr})
	}
}

// Java method: TypeValidator.getJSType(com.google.javascript.rhino.Node n)
func (this_ *TypeValidator) priv_TypeValidator_getJSType_5656(n pkg_rhino.Node) pkg_rhino.JSType {
	_ = n
	return j2gort.CheckNotNullIface(n.GetJSType())
}

// Java method: TypeValidator.getNativeType(com.google.javascript.rhino.jstype.JSTypeNative typeId)
func (this_ *TypeValidator) priv_TypeValidator_getNativeType_5657(typeId pkg_rhino.JSTypeNative) pkg_rhino.JSType {
	_ = typeId
	return this_.typeRegistry.GetNativeType(typeId)
}

// Java static method: TypeValidator.hasDuplicateDeclarationSuppression(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node decl)
func TypeValidator_HasDuplicateDeclarationSuppression(compiler AbstractCompiler, decl pkg_rhino.Node) bool {
	clinit_TypeValidator()
	_ = compiler
	_ = decl
	var originalDeclLevel *CheckLevel = compiler.GetErrorLevel(JSError_Make_3(decl, StaticGet_TypeValidator_DUP_VAR_DECLARATION(), []string{"dummy", "dummy"}))
	_ = originalDeclLevel
	return (originalDeclLevel == StaticGet_CheckLevel_OFF())
}

// Java method: TypeValidator.expectWellFormedTemplatizedType(com.google.javascript.rhino.Node n)
func (this_ *TypeValidator) ExpectWellFormedTemplatizedType(n pkg_rhino.Node) {
	_ = n
	var type_ pkg_rhino.JSType = n.GetJSType()
	_ = type_
	if type_ == nil || type_.IsJavaNil_() {
		return
	}
	if !(this_.verifierInUse) {
		var verifier *TypeValidator_WellFormedTemplatizedTypeVerifier = this_.reusableVerifier
		_ = verifier
		if verifier == nil {
			verifier = NewTypeValidator_WellFormedTemplatizedTypeVerifier(this_, pkg_rhino.Node(0))
			this_.reusableVerifier = verifier
		}
		this_.verifierInUse = true
		{
			r_13 := verifier
			r_13.node = n
		}
		defer func() {
			{
				r_14 := verifier
				r_14.node = pkg_rhino.Node(0)
			}
			this_.verifierInUse = false
		}()
		{
			type_.Visit_1(pkg_rhino.ToIface_Visitor(verifier))
		}
	} else {
		type_.Visit_1(NewTypeValidator_WellFormedTemplatizedTypeVerifier(this_, n))
	}
}

// Java class: com.google.javascript.jscomp.TypeValidator.WellFormedTemplatizedTypeVerifier
type TypeValidator_WellFormedTemplatizedTypeVerifier struct {
	pkg_rhino.Visitor_WithDefaultCaseBase
	outer_ *TypeValidator
	node   pkg_rhino.Node `gson:"node"`
}

func (this_ *TypeValidator_WellFormedTemplatizedTypeVerifier) BindSelf_(s any) {
	_ = s
	this_.Visitor_WithDefaultCaseBase.BindSelf_(s)
}
func (this_ *TypeValidator_WellFormedTemplatizedTypeVerifier) IsJavaNil_() bool { return this_ == nil }
func (*TypeValidator_WellFormedTemplatizedTypeVerifier) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TypeValidator$WellFormedTemplatizedTypeVerifier", "WellFormedTemplatizedTypeVerifier"
}
func (this_ *TypeValidator_WellFormedTemplatizedTypeVerifier) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.Visitor_WithDefaultCaseBase.IdentityHash_()
}
func (*TypeValidator_WellFormedTemplatizedTypeVerifier) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TypeValidator.WellFormedTemplatizedTypeVerifier", "com.google.javascript.jscomp.TypeValidator$WellFormedTemplatizedTypeVerifier", func(v_ any) bool {
		return j2gort.InstanceOf[*TypeValidator_WellFormedTemplatizedTypeVerifier](v_)
	})
}
func (this_ *TypeValidator_WellFormedTemplatizedTypeVerifier) SetOuter_(o *TypeValidator) {
	this_.outer_ = o
}
func (this_ *TypeValidator_WellFormedTemplatizedTypeVerifier) GetOuter_() *TypeValidator {
	return this_.outer_
}
func (*TypeValidator_WellFormedTemplatizedTypeVerifier) Implements_rhino_Visitor_WithDefaultCase() {}
func (*TypeValidator_WellFormedTemplatizedTypeVerifier) Implements_rhino_Visitor()                 {}

func (this_ *TypeValidator_WellFormedTemplatizedTypeVerifier) initFields_TypeValidator_WellFormedTemplatizedTypeVerifier() {
	_ = this_
}

// Java constructor: WellFormedTemplatizedTypeVerifier(com.google.javascript.rhino.Node node)
func NewTypeValidator_WellFormedTemplatizedTypeVerifier(outer_ *TypeValidator, node pkg_rhino.Node) *TypeValidator_WellFormedTemplatizedTypeVerifier {
	this_ := &TypeValidator_WellFormedTemplatizedTypeVerifier{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitTypeValidator_WellFormedTemplatizedTypeVerifier(this_, node)
	return this_
}

func InitTypeValidator_WellFormedTemplatizedTypeVerifier(this_ *TypeValidator_WellFormedTemplatizedTypeVerifier, node pkg_rhino.Node) {
	_ = node
	pkg_rhino.InitVisitor_WithDefaultCase(&this_.Visitor_WithDefaultCaseBase)
	this_.initFields_TypeValidator_WellFormedTemplatizedTypeVerifier()
	this_.node = node
}

// Java method: WellFormedTemplatizedTypeVerifier.caseDefault(com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType type)
func (this_ *TypeValidator_WellFormedTemplatizedTypeVerifier) CaseDefault(type_ pkg_rhino.JSType) any {
	_ = type_
	return j2gort.ValueOfBoolean(true)
}

// Java method: WellFormedTemplatizedTypeVerifier.caseEnumElementType(com.google.javascript.rhino.jstype.EnumElementType type)
func (this_ *TypeValidator_WellFormedTemplatizedTypeVerifier) CaseEnumElementType(type_ *pkg_rhino.EnumElementType) any {
	_ = type_
	return j2gort.ValueOfBoolean((j2gort.IfaceEquals[pkg_rhino.JSType](type_.GetPrimitiveType(), nil) || j2gort.UnboxBoolean(j2gort.CastPtr[j2gort.BoxBoolean](type_.GetPrimitiveType().Visit_1(this_)))))
}

// Java method: WellFormedTemplatizedTypeVerifier.caseFunctionType(com.google.javascript.rhino.jstype.FunctionType type)
func (this_ *TypeValidator_WellFormedTemplatizedTypeVerifier) CaseFunctionType(type_ *pkg_rhino.FunctionType) any {
	_ = type_
	var params []pkg_rhino.FunctionType_Parameter = type_.GetParameterArray()
	_ = params
	{
		var i int32 = int32(0)
		_ = i
		var len_ int32 = int32(len(params))
		_ = len_
		for ; i < len_; i = i + 1 {
			if !(j2gort.UnboxBoolean(j2gort.CastPtr[j2gort.BoxBoolean](params[i].GetJSType().Visit_1(this_)))) {
				return j2gort.ValueOfBoolean(false)
			}
		}
	}
	return j2gort.PtrToAny(j2gort.CastPtr[j2gort.BoxBoolean](type_.GetReturnType().Visit_1(this_)))
}

// Java method: WellFormedTemplatizedTypeVerifier.caseNamedType(com.google.javascript.rhino.jstype.NamedType type)
func (this_ *TypeValidator_WellFormedTemplatizedTypeVerifier) CaseNamedType(type_ *pkg_rhino.NamedType) any {
	_ = type_
	return j2gort.PtrToAny(j2gort.CastPtr[j2gort.BoxBoolean](type_.GetReferencedType().Visit_1(this_)))
}

// Java method: WellFormedTemplatizedTypeVerifier.caseUnionType(com.google.javascript.rhino.jstype.UnionType type)
func (this_ *TypeValidator_WellFormedTemplatizedTypeVerifier) CaseUnionType(type_ *pkg_rhino.UnionType) any {
	_ = type_
	var alternateCount int32 = type_.GetAlternatesCount()
	_ = alternateCount
	{
		var i int32 = int32(0)
		_ = i
		for ; i < alternateCount; i = i + 1 {
			var alternative pkg_rhino.JSType = type_.GetAlternateAt(i)
			_ = alternative
			if !(j2gort.UnboxBoolean(j2gort.CastPtr[j2gort.BoxBoolean](alternative.Visit_1(this_)))) {
				return j2gort.ValueOfBoolean(false)
			}
		}
	}
	return j2gort.ValueOfBoolean(true)
}

// Java method: WellFormedTemplatizedTypeVerifier.caseTemplatizedType(com.google.javascript.rhino.jstype.TemplatizedType type)
func (this_ *TypeValidator_WellFormedTemplatizedTypeVerifier) CaseTemplatizedType(type_ *pkg_rhino.TemplatizedType) any {
	_ = type_
	var referencedTemplates *pkg_rhino.TemplateTypeMap = type_.GetReferencedType().GetTemplateTypeMap()
	_ = referencedTemplates
	var templateTypes []pkg_rhino.JSType = j2gort.CheckNotNullPtr(type_).Slice_GetTemplateTypes()
	_ = templateTypes
	var refCount int32 = referencedTemplates.Size()
	_ = refCount
	{
		var i int32 = int32(0)
		_ = i
		var len_ int32 = int32(len(templateTypes))
		_ = len_
		for ; i < len_; i = i + 1 {
			var assignedType pkg_rhino.JSType = templateTypes[i]
			_ = assignedType
			if i < refCount {
				var templateType *pkg_rhino.TemplateType = referencedTemplates.GetTemplateKeyAt(i)
				_ = templateType
				if !(assignedType.IsSubtype_1(templateType.GetBound())) {
					this_.outer_.priv_TypeValidator_registerMismatchAndReport_5655(assignedType, templateType.GetBound(), JSError_Make_3(this_.node, StaticGet_RhinoErrorReporter_BOUNDED_GENERIC_TYPE_ERROR(), []string{assignedType.String(), templateType.GetReferenceName(), templateType.GetBound().String()}))
					return j2gort.ValueOfBoolean(false)
				}
			}
		}
	}
	return j2gort.ValueOfBoolean(true)
}

// Java method: WellFormedTemplatizedTypeVerifier.caseTemplateType(com.google.javascript.rhino.jstype.TemplateType templateType)
func (this_ *TypeValidator_WellFormedTemplatizedTypeVerifier) CaseTemplateType(templateType *pkg_rhino.TemplateType) any {
	_ = templateType
	return j2gort.ValueOfBoolean((!(templateType.ContainsCycle()) && j2gort.UnboxBoolean(j2gort.CastPtr[j2gort.BoxBoolean](templateType.GetBound().Visit_1(this_)))))
}
