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
	compiler                AbstractCompiler          `gson:"-"`
	typeRegistry            *pkg_rhino.JSTypeRegistry `gson:"typeRegistry"`
	allBitwisableValueTypes pkg_rhino.JSType          `gson:"allBitwisableValueTypes"`
	nullOrUndefined         pkg_rhino.JSType          `gson:"nullOrUndefined"`
	promiseOfUnknownType    pkg_rhino.JSType          `gson:"promiseOfUnknownType"`
	iterableOrAsyncIterable pkg_rhino.JSType          `gson:"iterableOrAsyncIterable"`
	mismatches              *TypeMismatch_Accumulator `gson:"mismatches"`
	idHash_                 atomic.Int32
	subtypingMode           pkg_rhino.JSType_SubtypingMode `gson:"subtypingMode"`
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
	this_.iterableOrAsyncIterable = this_.typeRegistry.CreateUnionType_2(pkg_rhino.Upcast_JSType(this_.typeRegistry.GetNativeObjectType(pkg_rhino.StaticGet_JSTypeNative_ITERATOR_TYPE())), pkg_rhino.Upcast_JSType(this_.typeRegistry.GetNativeObjectType(pkg_rhino.StaticGet_JSTypeNative_ASYNC_ITERATOR_TYPE())))
}

// Java static method: TypeValidator.getInstanceOfCtor(com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType t)
func TypeValidator_GetInstanceOfCtor(t pkg_rhino.JSType) pkg_rhino.ObjectType {
	clinit_TypeValidator()
	_ = t
	if t == nil || t.IsJavaNil_() {
		return nil
	}
	var ctor *pkg_rhino.FunctionType = pkg_rhino.JSType_ToMaybeFunctionType(pkg_rhino.Upcast_JSType(t.GetJSTypeBase().Dereference()))
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
		this_.priv_TypeValidator_mismatch_5663(n, msg, type_, pkg_rhino.StaticGet_JSTypeNative_OBJECT_TYPE())
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
		this_.priv_TypeValidator_mismatch_5663(n, msg, type_, pkg_rhino.StaticGet_JSTypeNative_OBJECT_TYPE())
	}
}

// Java method: TypeValidator.expectAnyObject(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType type, java.lang.String msg)
func (this_ *TypeValidator) ExpectAnyObject(n pkg_rhino.Node, type_ pkg_rhino.JSType, msg string) {
	_ = n
	_ = type_
	_ = msg
	var anyObjectType pkg_rhino.JSType = this_.priv_TypeValidator_getNativeType_5668(pkg_rhino.StaticGet_JSTypeNative_NO_OBJECT_TYPE())
	_ = anyObjectType
	if !(anyObjectType.GetJSTypeBase().IsSubtypeOf_1(type_)) && !(type_.GetJSTypeBase().IsEmptyType()) {
		this_.priv_TypeValidator_mismatch_5664(n, msg, type_, anyObjectType)
	}
}

// Java method: TypeValidator.expectAutoboxesToIterable(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType type, java.lang.String msg)
func (this_ *TypeValidator) ExpectAutoboxesToIterable(n pkg_rhino.Node, type_ pkg_rhino.JSType, msg string) bool {
	_ = n
	_ = type_
	_ = msg
	base_type__ := type_.GetJSTypeBase()
	if base_type__.IsUnionType() {
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
				if alt.GetJSTypeBase().IsBoxableScalar() {
					alt = alt.AutoboxesTo()
				} else {
					alt = alt
				}
				if !(alt.GetJSTypeBase().IsSubtypeOf_1(this_.priv_TypeValidator_getNativeType_5668(pkg_rhino.StaticGet_JSTypeNative_ITERABLE_TYPE()))) {
					this_.priv_TypeValidator_mismatch_5663(n, msg, type_, pkg_rhino.StaticGet_JSTypeNative_ITERABLE_TYPE())
					return false
				}
			}
		}
	} else {
		var autoboxedType pkg_rhino.JSType
		_ = autoboxedType
		if base_type__.IsBoxableScalar() {
			autoboxedType = type_.AutoboxesTo()
		} else {
			autoboxedType = type_
		}
		if !(autoboxedType.GetJSTypeBase().IsSubtypeOf_1(this_.priv_TypeValidator_getNativeType_5668(pkg_rhino.StaticGet_JSTypeNative_ITERABLE_TYPE()))) {
			this_.priv_TypeValidator_mismatch_5663(n, msg, type_, pkg_rhino.StaticGet_JSTypeNative_ITERABLE_TYPE())
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
	this_.priv_TypeValidator_mismatch_5664(n, msg, type_, this_.iterableOrAsyncIterable)
	return j2gort.Ext_java_util_Optional_Empty()
}

// Java method: TypeValidator.expectGeneratorSupertype(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType type, java.lang.String msg)
func (this_ *TypeValidator) ExpectGeneratorSupertype(n pkg_rhino.Node, type_ pkg_rhino.JSType, msg string) {
	_ = n
	_ = type_
	_ = msg
	if !(this_.priv_TypeValidator_getNativeType_5668(pkg_rhino.StaticGet_JSTypeNative_GENERATOR_TYPE()).GetJSTypeBase().IsSubtypeOf_1(type_)) {
		this_.priv_TypeValidator_mismatch_5663(n, msg, type_, pkg_rhino.StaticGet_JSTypeNative_GENERATOR_TYPE())
	}
}

// Java method: TypeValidator.expectAsyncGeneratorSupertype(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType type, java.lang.String msg)
func (this_ *TypeValidator) ExpectAsyncGeneratorSupertype(n pkg_rhino.Node, type_ pkg_rhino.JSType, msg string) {
	_ = n
	_ = type_
	_ = msg
	if !(this_.priv_TypeValidator_getNativeType_5668(pkg_rhino.StaticGet_JSTypeNative_ASYNC_GENERATOR_TYPE()).GetJSTypeBase().IsSubtypeOf_1(type_)) {
		this_.priv_TypeValidator_mismatch_5663(n, msg, type_, pkg_rhino.StaticGet_JSTypeNative_ASYNC_GENERATOR_TYPE())
	}
}

// Java method: TypeValidator.expectValidAsyncReturnType(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType type)
func (this_ *TypeValidator) ExpectValidAsyncReturnType(n pkg_rhino.Node, type_ pkg_rhino.JSType) {
	_ = n
	_ = type_
	if this_.promiseOfUnknownType.GetJSTypeBase().IsSubtypeOf_1(type_) {
		return
	}
	var err *JSError = JSError_Make_3(n, StaticGet_TypeValidator_INVALID_ASYNC_RETURN_TYPE(), []string{type_.GetJSTypeBase().String()})
	_ = err
	this_.priv_TypeValidator_registerMismatchAndReport_5666(type_, this_.promiseOfUnknownType, err)
}

// Java method: TypeValidator.expectITemplateArraySupertype(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType type, java.lang.String msg)
func (this_ *TypeValidator) ExpectITemplateArraySupertype(n pkg_rhino.Node, type_ pkg_rhino.JSType, msg string) {
	_ = n
	_ = type_
	_ = msg
	if !(this_.priv_TypeValidator_getNativeType_5668(pkg_rhino.StaticGet_JSTypeNative_I_TEMPLATE_ARRAY_TYPE()).GetJSTypeBase().IsSubtypeOf_1(type_)) {
		this_.priv_TypeValidator_mismatch_5663(n, msg, type_, pkg_rhino.StaticGet_JSTypeNative_I_TEMPLATE_ARRAY_TYPE())
	}
}

// Java method: TypeValidator.expectString(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType type, java.lang.String msg)
func (this_ *TypeValidator) ExpectString(n pkg_rhino.Node, type_ pkg_rhino.JSType, msg string) {
	_ = n
	_ = type_
	_ = msg
	if !(type_.MatchesStringContext()) {
		this_.priv_TypeValidator_mismatch_5663(n, msg, type_, pkg_rhino.StaticGet_JSTypeNative_STRING_TYPE())
	}
}

// Java method: TypeValidator.expectNumber(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType type, java.lang.String msg)
func (this_ *TypeValidator) ExpectNumber(n pkg_rhino.Node, type_ pkg_rhino.JSType, msg string) {
	_ = n
	_ = type_
	_ = msg
	if !(type_.MatchesNumberContext()) {
		this_.priv_TypeValidator_mismatch_5663(n, msg, type_, pkg_rhino.StaticGet_JSTypeNative_NUMBER_TYPE())
	} else {
		this_.ExpectNumberStrict(n, type_, msg)
	}
}

// Java method: TypeValidator.expectNumberStrict(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType type, java.lang.String msg)
func (this_ *TypeValidator) ExpectNumberStrict(n pkg_rhino.Node, type_ pkg_rhino.JSType, msg string) {
	_ = n
	_ = type_
	_ = msg
	if !(type_.GetJSTypeBase().IsSubtypeOf_1(this_.priv_TypeValidator_getNativeType_5668(pkg_rhino.StaticGet_JSTypeNative_NUMBER_TYPE()))) {
		this_.priv_TypeValidator_registerMismatchAndReport_5665(n, StaticGet_TypeValidator_INVALID_OPERAND_TYPE(), msg, type_, this_.priv_TypeValidator_getNativeType_5668(pkg_rhino.StaticGet_JSTypeNative_NUMBER_TYPE()), nil, nil)
	}
}

// Java method: TypeValidator.expectMatchingTypesStrict(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType left, com.google.javascript.rhino.jstype.JSType right, java.lang.String msg)
func (this_ *TypeValidator) ExpectMatchingTypesStrict(n pkg_rhino.Node, left pkg_rhino.JSType, right pkg_rhino.JSType, msg string) {
	_ = n
	_ = left
	_ = right
	_ = msg
	if !(left.GetJSTypeBase().IsSubtypeOf_1(right)) && !(right.GetJSTypeBase().IsSubtypeOf_1(left)) {
		this_.priv_TypeValidator_registerMismatchAndReport_5665(n, StaticGet_TypeValidator_INVALID_OPERAND_TYPE(), msg, right, left, nil, nil)
	}
}

// Java method: TypeValidator.expectBitwiseable(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType type, java.lang.String msg)
func (this_ *TypeValidator) ExpectBitwiseable(n pkg_rhino.Node, type_ pkg_rhino.JSType, msg string) {
	_ = n
	_ = type_
	_ = msg
	if !(type_.MatchesNumberContext()) && !(type_.GetJSTypeBase().IsSubtypeOf_1(this_.allBitwisableValueTypes)) {
		this_.priv_TypeValidator_mismatch_5664(n, msg, type_, this_.allBitwisableValueTypes)
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
		this_.priv_TypeValidator_mismatch_5663(n, msg, type_, pkg_rhino.StaticGet_JSTypeNative_NUMBER_SYMBOL())
	}
}

// Java method: TypeValidator.expectStringOrSymbol(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType type, java.lang.String msg)
func (this_ *TypeValidator) ExpectStringOrSymbol(n pkg_rhino.Node, type_ pkg_rhino.JSType, msg string) {
	_ = n
	_ = type_
	_ = msg
	if !(type_.MatchesStringContext()) && !(type_.MatchesSymbolContext()) {
		this_.priv_TypeValidator_mismatch_5663(n, msg, type_, pkg_rhino.StaticGet_JSTypeNative_STRING_SYMBOL())
	}
}

// Java method: TypeValidator.expectUnknownOrComparable(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType type, java.lang.String msg)
func (this_ *TypeValidator) ExpectUnknownOrComparable(n pkg_rhino.Node, type_ pkg_rhino.JSType, msg string) {
	_ = n
	_ = type_
	_ = msg
	if !(type_.GetJSTypeBase().IsSubtypeOf_1(this_.priv_TypeValidator_getNativeType_5668(pkg_rhino.StaticGet_JSTypeNative_BIGINT_NUMBER_STRING()))) && !(type_.GetJSTypeBase().IsSubtypeOf_1(this_.priv_TypeValidator_getNativeType_5668(pkg_rhino.StaticGet_JSTypeNative_BIGINT_NUMBER_STRING_OBJECT()))) {
		if type_.MatchesNumberContext() || type_.MatchesStringContext() {
			this_.priv_TypeValidator_registerMismatchAndReport_5665(n, StaticGet_TypeValidator_INVALID_OPERAND_TYPE(), msg, type_, this_.priv_TypeValidator_getNativeType_5668(pkg_rhino.StaticGet_JSTypeNative_BIGINT_NUMBER_STRING()), nil, nil)
		} else {
			this_.priv_TypeValidator_mismatch_5663(n, msg, type_, pkg_rhino.StaticGet_JSTypeNative_BIGINT_NUMBER_STRING())
		}
	}
}

// Java method: TypeValidator.expectStringOrNumberOrSymbol(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType type, java.lang.String msg)
func (this_ *TypeValidator) ExpectStringOrNumberOrSymbol(n pkg_rhino.Node, type_ pkg_rhino.JSType, msg string) {
	_ = n
	_ = type_
	_ = msg
	if (!(type_.MatchesNumberContext()) && !(type_.MatchesStringContext())) && !(type_.MatchesSymbolContext()) {
		this_.priv_TypeValidator_mismatch_5663(n, msg, type_, pkg_rhino.StaticGet_JSTypeNative_NUMBER_STRING_SYMBOL())
	} else {
		this_.ExpectStringOrNumberOrSymbolStrict(n, type_, msg)
	}
}

// Java method: TypeValidator.expectStringOrNumberOrSymbolStrict(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType type, java.lang.String msg)
func (this_ *TypeValidator) ExpectStringOrNumberOrSymbolStrict(n pkg_rhino.Node, type_ pkg_rhino.JSType, msg string) {
	_ = n
	_ = type_
	_ = msg
	if !(type_.GetJSTypeBase().IsSubtypeOf_1(this_.priv_TypeValidator_getNativeType_5668(pkg_rhino.StaticGet_JSTypeNative_NUMBER_STRING_SYMBOL()))) {
		this_.priv_TypeValidator_registerMismatchAndReport_5665(n, StaticGet_TypeValidator_INVALID_OPERAND_TYPE(), msg, type_, this_.priv_TypeValidator_getNativeType_5668(pkg_rhino.StaticGet_JSTypeNative_NUMBER_STRING_SYMBOL()), nil, nil)
	}
}

// Java method: TypeValidator.expectBigIntOrNumber(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType type, java.lang.String msg)
func (this_ *TypeValidator) ExpectBigIntOrNumber(n pkg_rhino.Node, type_ pkg_rhino.JSType, msg string) {
	_ = n
	_ = type_
	_ = msg
	if !(type_.GetJSTypeBase().IsSubtypeOf_1(this_.priv_TypeValidator_getNativeType_5668(pkg_rhino.StaticGet_JSTypeNative_BIGINT_NUMBER()))) && !(type_.GetJSTypeBase().IsSubtypeOf_1(this_.priv_TypeValidator_getNativeType_5668(pkg_rhino.StaticGet_JSTypeNative_BIGINT_NUMBER_OBJECT()))) {
		if type_.MatchesNumberContext() {
			this_.priv_TypeValidator_registerMismatchAndReport_5665(n, StaticGet_TypeValidator_INVALID_OPERAND_TYPE(), msg, type_, this_.priv_TypeValidator_getNativeType_5668(pkg_rhino.StaticGet_JSTypeNative_BIGINT_NUMBER()), nil, nil)
		} else {
			this_.priv_TypeValidator_mismatch_5663(n, msg, type_, pkg_rhino.StaticGet_JSTypeNative_BIGINT_NUMBER())
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
	base_type__ := type_.GetJSTypeBase()
	if ((!(base_type__.IsNoType()) && !(base_type__.IsUnknownType())) && base_type__.IsSubtypeOf_1(this_.nullOrUndefined)) && !(TypeValidator_ContainsForwardDeclaredUnresolvedName(type_)) {
		if (n.IsGetProp() && !(t.InGlobalScope())) && base_type__.IsNullType() {
			return true
		}
		this_.priv_TypeValidator_mismatch_5664(n, msg, type_, expectedType)
		return false
	}
	return true
}

// Java static method: TypeValidator.containsForwardDeclaredUnresolvedName(com.google.javascript.rhino.jstype.JSType type)
func TypeValidator_ContainsForwardDeclaredUnresolvedName(type_ pkg_rhino.JSType) bool {
	clinit_TypeValidator()
	_ = type_
	base_type__ := type_.GetJSTypeBase()
	if base_type__.IsUnionType() {
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
	return base_type__.IsNoResolvedType()
}

// Java method: TypeValidator.expectSwitchMatchesCase(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType switchType, com.google.javascript.rhino.jstype.JSType caseType)
func (this_ *TypeValidator) ExpectSwitchMatchesCase(n pkg_rhino.Node, switchType pkg_rhino.JSType, caseType pkg_rhino.JSType) {
	_ = n
	_ = switchType
	_ = caseType
	if !(switchType.GetJSTypeBase().CanTestForShallowEqualityWith(caseType)) {
		this_.priv_TypeValidator_mismatch_5664(n.GetFirstChild(), "case expression doesn't match switch", caseType, switchType)
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
	if indexType.GetJSTypeBase().IsSymbolValueType() {
		return
	}
	base_objType_ := objType.GetJSTypeBase()
	if base_objType_.IsUnknownType() {
		this_.ExpectStringOrNumberOrSymbol(indexNode, indexType, "property access")
		return
	}
	var dereferenced pkg_rhino.ObjectType = base_objType_.Dereference()
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
		this_.priv_TypeValidator_mismatch_5664(n, "only arrays or objects can be accessed", objType, this_.typeRegistry.CreateUnionType_1_JSTypeNative([]pkg_rhino.JSTypeNative{pkg_rhino.StaticGet_JSTypeNative_ARRAY_TYPE(), pkg_rhino.StaticGet_JSTypeNative_OBJECT_TYPE()}))
	}
}

// Java method: TypeValidator.expectCanAssignToPropertyOf(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType rightType, com.google.javascript.rhino.jstype.JSType leftType, com.google.javascript.rhino.Node owner, java.lang.String propName)
func (this_ *TypeValidator) ExpectCanAssignToPropertyOf_5(n pkg_rhino.Node, rightType pkg_rhino.JSType, leftType pkg_rhino.JSType, owner pkg_rhino.Node, propName string) bool {
	_ = n
	_ = rightType
	_ = leftType
	_ = owner
	_ = propName
	base_leftType_ := leftType.GetJSTypeBase()
	if base_leftType_.IsTemplateType() {
		var left *pkg_rhino.TemplateType = leftType.ToMaybeTemplateType()
		_ = left
		if (rightType.GetJSTypeBase().ContainsReferenceAncestor(pkg_rhino.ToIface_JSType(left)) || rightType.GetJSTypeBase().IsUnknownType()) || left.IsUnknownType() {
			return true
		}
		var ownerTypeName string
		_ = ownerTypeName
		if n.IsMemberFieldDef() {
			ownerTypeName = n.GetGrandparent().GetJSType().GetJSTypeBase().AssertFunctionType().GetInstanceType().GetObjectTypeBase().String()
		} else {
			ownerTypeName = this_.typeRegistry.GetReadableTypeName(owner)
		}
		this_.priv_TypeValidator_registerMismatchAndReport_5665(n, StaticGet_TypeValidator_TYPE_MISMATCH_WARNING(), ("assignment to property " + j2gort.StrValString(propName) + " of " + j2gort.StrValString(ownerTypeName)), rightType, leftType, j2gort.NewExt_java_util_LinkedHashSet(), j2gort.NewExt_java_util_LinkedHashSet())
		return false
	}
	if !(base_leftType_.IsNoType()) && !(rightType.GetJSTypeBase().IsSubtypeOf_1(leftType)) {
		var ownerType pkg_rhino.JSType = this_.priv_TypeValidator_getJSType_5667(owner)
		_ = ownerType
		if ownerType.IsFunctionPrototypeType() {
			var ownerFn *pkg_rhino.FunctionType = ownerType.GetJSTypeBase().ToObjectType().GetOwnerFunction()
			_ = ownerFn
			if (ownerFn.IsInterface() && rightType.GetJSTypeBase().IsFunctionType()) && base_leftType_.IsFunctionType() {
				return true
			}
		}
		var ownerTypeName_2 string
		_ = ownerTypeName_2
		if n.IsMemberFieldDef() {
			ownerTypeName_2 = n.GetGrandparent().GetJSType().GetJSTypeBase().AssertFunctionType().GetInstanceType().GetObjectTypeBase().String()
		} else {
			ownerTypeName_2 = this_.typeRegistry.GetReadableTypeName(owner)
		}
		this_.priv_TypeValidator_mismatch_5664(n, ("assignment to property " + j2gort.StrValString(propName) + " of " + j2gort.StrValString(ownerTypeName_2)), rightType, leftType)
		return false
	}
	return true
}

// Java method: TypeValidator.expectCanAssignToPropertyOf(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType rightType, com.google.javascript.rhino.jstype.JSType leftType, com.google.javascript.rhino.jstype.JSType ownerType, com.google.javascript.rhino.QualifiedName ownerName, java.lang.String propName)
func (this_ *TypeValidator) ExpectCanAssignToPropertyOf_6_Node_JSType_JSType_JSType_QualifiedName_string(n pkg_rhino.Node, rightType pkg_rhino.JSType, leftType pkg_rhino.JSType, ownerType pkg_rhino.JSType, ownerName pkg_rhino.QualifiedName, propName string) bool {
	_ = n
	_ = rightType
	_ = leftType
	_ = ownerType
	_ = ownerName
	_ = propName
	base_leftType_ := leftType.GetJSTypeBase()
	if base_leftType_.IsTemplateType() {
		var left *pkg_rhino.TemplateType = leftType.ToMaybeTemplateType()
		_ = left
		if (rightType.GetJSTypeBase().ContainsReferenceAncestor(pkg_rhino.ToIface_JSType(left)) || rightType.GetJSTypeBase().IsUnknownType()) || left.IsUnknownType() {
			return true
		}
		this_.priv_TypeValidator_registerMismatchAndReport_5665(n, StaticGet_TypeValidator_TYPE_MISMATCH_WARNING(), ("assignment to property " + j2gort.StrValString(propName) + " of " + j2gort.StrValString(ownerName.Join_0())), rightType, leftType, j2gort.NewExt_java_util_LinkedHashSet(), j2gort.NewExt_java_util_LinkedHashSet())
		return false
	}
	if !(base_leftType_.IsNoType()) && !(rightType.GetJSTypeBase().IsSubtypeOf_1(leftType)) {
		if ownerType.IsFunctionPrototypeType() {
			var ownerFn *pkg_rhino.FunctionType = ownerType.GetJSTypeBase().ToObjectType().GetOwnerFunction()
			_ = ownerFn
			if (ownerFn.IsInterface() && rightType.GetJSTypeBase().IsFunctionType()) && base_leftType_.IsFunctionType() {
				return true
			}
		}
		this_.priv_TypeValidator_mismatch_5664(n, ("assignment to property " + j2gort.StrValString(propName) + " of " + j2gort.StrValString(ownerName.Join_0())), rightType, leftType)
		return false
	}
	return true
}

// Java method: TypeValidator.expectCanAssignToPropertyOf(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType rightType, com.google.javascript.rhino.jstype.JSType leftType, com.google.javascript.rhino.jstype.JSType ownerType, com.google.common.base.Supplier<java.lang.String> typeNameSupplier, java.lang.String propName)
func (this_ *TypeValidator) ExpectCanAssignToPropertyOf_6_Node_JSType_JSType_JSType_Ext_com_google_common_base_Supplier_string(n pkg_rhino.Node, rightType pkg_rhino.JSType, leftType pkg_rhino.JSType, ownerType pkg_rhino.JSType, typeNameSupplier j2gort.Ext_com_google_common_base_Supplier, propName string) bool {
	_ = n
	_ = rightType
	_ = leftType
	_ = ownerType
	_ = typeNameSupplier
	_ = propName
	base_leftType_ := leftType.GetJSTypeBase()
	if base_leftType_.IsTemplateType() {
		var left *pkg_rhino.TemplateType = leftType.ToMaybeTemplateType()
		_ = left
		if (rightType.GetJSTypeBase().ContainsReferenceAncestor(pkg_rhino.ToIface_JSType(left)) || rightType.GetJSTypeBase().IsUnknownType()) || left.IsUnknownType() {
			return true
		} else {
			this_.priv_TypeValidator_registerMismatchAndReport_5665(n, StaticGet_TypeValidator_TYPE_MISMATCH_WARNING(), ("assignment to property " + j2gort.StrValString(propName) + " of " + j2gort.StrValString(j2gort.CastString(typeNameSupplier.Get()))), rightType, leftType, j2gort.NewExt_java_util_LinkedHashSet(), j2gort.NewExt_java_util_LinkedHashSet())
			return false
		}
	}
	if !(base_leftType_.IsNoType()) && !(rightType.GetJSTypeBase().IsSubtypeOf_1(leftType)) {
		if ownerType.IsFunctionPrototypeType() {
			var ownerFn *pkg_rhino.FunctionType = ownerType.GetJSTypeBase().ToObjectType().GetOwnerFunction()
			_ = ownerFn
			if (ownerFn.IsInterface() && rightType.GetJSTypeBase().IsFunctionType()) && base_leftType_.IsFunctionType() {
				return true
			}
		}
		this_.priv_TypeValidator_mismatch_5664(n, ("assignment to property " + j2gort.StrValString(propName) + " of " + j2gort.StrValString(j2gort.CastString(typeNameSupplier.Get()))), rightType, leftType)
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
	if leftType.GetJSTypeBase().IsTemplateType() {
		var left *pkg_rhino.TemplateType = leftType.ToMaybeTemplateType()
		_ = left
		if (rightType.GetJSTypeBase().ContainsReferenceAncestor(pkg_rhino.ToIface_JSType(left)) || rightType.GetJSTypeBase().IsUnknownType()) || left.IsUnknownType() {
			return true
		} else {
			this_.priv_TypeValidator_registerMismatchAndReport_5665(n, StaticGet_TypeValidator_TYPE_MISMATCH_WARNING(), msg, rightType, leftType, j2gort.NewExt_java_util_LinkedHashSet(), j2gort.NewExt_java_util_LinkedHashSet())
			return false
		}
	}
	if !(rightType.GetJSTypeBase().IsSubtypeOf_1(leftType)) {
		this_.priv_TypeValidator_mismatch_5664(n, msg, rightType, leftType)
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
	if !(argType.GetJSTypeBase().IsSubtypeOf_1(paramType)) {
		this_.priv_TypeValidator_mismatch_5664(n, ("actual parameter " + j2gort.StrValInt32(ordinal) + " of " + j2gort.StrValString(this_.typeRegistry.GetReadableTypeNameNoDeref(callNode.GetFirstChild())) + " does not match formal parameter"), argType, paramType)
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
	if (declaredSuper != nil && !declaredSuper.IsJavaNil_()) && declaredSuper.GetObjectTypeBase().IsTemplatizedType() {
		declaredSuper = declaredSuper.ToMaybeTemplatizedType().GetReferencedType()
	}
	if ((declaredSuper != nil && !declaredSuper.IsJavaNil_()) && !(j2gort.InstanceOf[*pkg_rhino.UnknownType](superObject))) && !(declaredSuper.GetObjectTypeBase().EqualsJSType(pkg_rhino.Upcast_JSType(superObject))) {
		if declaredSuper.GetObjectTypeBase().EqualsJSType(this_.priv_TypeValidator_getNativeType_5668(pkg_rhino.StaticGet_JSTypeNative_OBJECT_TYPE())) {
			this_.priv_TypeValidator_registerMismatchAndReport_5666(pkg_rhino.Upcast_JSType(superObject), pkg_rhino.Upcast_JSType(declaredSuper), JSError_Make_3(n, StaticGet_TypeValidator_MISSING_EXTENDS_TAG_WARNING(), []string{subObject.GetObjectTypeBase().String()}))
		} else {
			this_.priv_TypeValidator_mismatch_5664(n, "mismatch in declaration of superclass type", pkg_rhino.Upcast_JSType(superObject), pkg_rhino.Upcast_JSType(declaredSuper))
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
			if !(astSuperInstance.GetObjectTypeBase().EqualsJSType(pkg_rhino.Upcast_JSType(registeredSuperInstance))) {
				this_.priv_TypeValidator_mismatch_5664(n, "mismatch in declaration of superclass type", pkg_rhino.Upcast_JSType(astSuperInstance), pkg_rhino.Upcast_JSType(registeredSuperInstance))
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
	if ownerType.GetJSTypeBase().IsFunctionType() {
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
	base_sourceType_ := sourceType.GetJSTypeBase()
	if !(base_sourceType_.CanCastTo(targetType)) {
		this_.priv_TypeValidator_registerMismatchAndReport_5666(sourceType, targetType, JSError_Make_3(n, StaticGet_TypeValidator_INVALID_CAST(), []string{base_sourceType_.String(), targetType.GetJSTypeBase().String()}))
	} else if !(base_sourceType_.IsSubtypeOf_1(targetType)) {
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
				if !(newType.GetJSTypeBase().EqualsJSType(varType)) {
					this_.compiler.Report_1(JSError_Make_3(n, StaticGet_TypeValidator_DUP_VAR_DECLARATION_TYPE_MISMATCH(), []string{variableName, newType.GetJSTypeBase().String(), var_.GetInputName(), strconv.FormatInt(int64(var_.GetNameNode().GetLineno()), 10), varType.GetJSTypeBase().String()}))
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
			this_.priv_TypeValidator_expectInterfaceProperties_5658(n, instance, implemented)
		}
	}
	{
		cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(type_.GetExtendedInterfaces()))
		for cur_3.HasNext() {
			rawElem_4 := cur_3.Next()
			extended := j2gort.AsIface[pkg_rhino.ObjectType](rawElem_4)
			_ = extended
			this_.priv_TypeValidator_expectInterfaceProperties_5658(n, instance, extended)
		}
	}
}

// Java method: TypeValidator.expectInterfaceProperties(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.ObjectType instance, com.google.javascript.rhino.jstype.ObjectType ancestorInterface)
func (this_ *TypeValidator) priv_TypeValidator_expectInterfaceProperties_5658(n pkg_rhino.Node, instance pkg_rhino.ObjectType, ancestorInterface pkg_rhino.ObjectType) {
	_ = n
	_ = instance
	_ = ancestorInterface
	{
		cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(ancestorInterface.GetObjectTypeBase().GetOwnPropertyKeys()))
		for cur_5.HasNext() {
			rawElem_6 := cur_5.Next()
			prop := j2gort.AsIface[pkg_rhino.Property_Key](rawElem_6)
			_ = prop
			this_.priv_TypeValidator_expectInterfaceProperty_5659(n, instance, ancestorInterface, prop)
		}
	}
	if !(j2gort.IfaceEquals[pkg_rhino.ObjectType](ancestorInterface.GetImplicitPrototype(), nil)) {
		{
			cur_7 := j2gort.NewIterCursor(j2gort.CheckNotNull(ancestorInterface.GetImplicitPrototype().GetObjectTypeBase().GetOwnPropertyKeys()))
			for cur_7.HasNext() {
				rawElem_8 := cur_7.Next()
				prop_2 := j2gort.AsIface[pkg_rhino.Property_Key](rawElem_8)
				_ = prop_2
				this_.priv_TypeValidator_expectInterfaceProperty_5659(n, instance, ancestorInterface, prop_2)
			}
		}
	}
}

// Java method: TypeValidator.expectInterfaceProperty(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.ObjectType instance, com.google.javascript.rhino.jstype.ObjectType implementedInterface, com.google.javascript.rhino.jstype.Property.Key propName)
func (this_ *TypeValidator) priv_TypeValidator_expectInterfaceProperty_5659(n pkg_rhino.Node, instance pkg_rhino.ObjectType, implementedInterface pkg_rhino.ObjectType, propName pkg_rhino.Property_Key) {
	_ = n
	_ = instance
	_ = implementedInterface
	_ = propName
	base_instance_ := instance.GetObjectTypeBase()
	var propSlot *pkg_rhino.Property_OwnedProperty = base_instance_.FindClosestDefinition_1_Property_Key(propName)
	_ = propSlot
	if (propSlot == nil) || (!(instance.GetConstructor().IsInterface()) && propSlot.IsOwnedByInterface()) {
		if instance.GetConstructor().IsAbstract() || instance.GetConstructor().IsInterface() {
			return
		}
		if implementedInterface.GetPropertyType_1_Property_Key(propName).IsVoidable() {
			return
		}
		this_.priv_TypeValidator_registerMismatchAndReport_5666(pkg_rhino.Upcast_JSType(instance), pkg_rhino.Upcast_JSType(implementedInterface), JSError_Make_3(n, StaticGet_TypeValidator_INTERFACE_METHOD_NOT_IMPLEMENTED(), []string{propName.HumanReadableName(), implementedInterface.GetReferenceName(), base_instance_.String()}))
	} else {
		var local bool = propSlot.GetOwnerInstanceType().GetObjectTypeBase().EqualsJSType(pkg_rhino.Upcast_JSType(instance))
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
		this_.CheckPropertyType_5_Node_JSType_ObjectType_Property_Key_JSType(propNode, pkg_rhino.Upcast_JSType(instance), implementedInterface, propName, prop.GetType())
	}
}

// Java method: TypeValidator.checkPropertyType(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType instance, com.google.javascript.rhino.jstype.ObjectType parent, java.lang.String propertyName, com.google.javascript.rhino.jstype.JSType found)
func (this_ *TypeValidator) CheckPropertyType_5_Node_JSType_ObjectType_string_JSType(n pkg_rhino.Node, instance pkg_rhino.JSType, parent pkg_rhino.ObjectType, propertyName string, found pkg_rhino.JSType) {
	_ = n
	_ = instance
	_ = parent
	_ = propertyName
	_ = found
	var required pkg_rhino.JSType = parent.GetObjectTypeBase().GetPropertyType_1_string(propertyName)
	_ = required
	var typeMap *pkg_rhino.TemplateTypeMap = instance.GetTemplateTypeMap()
	_ = typeMap
	if !(typeMap.IsEmpty()) && required.GetJSTypeBase().HasAnyTemplateTypes() {
		required = this_.typeRegistry.ReplaceTemplateTypesPartial(required, typeMap)
	}
	if found.GetJSTypeBase().IsSubtypeOf_2(required, this_.subtypingMode) {
		return
	}
	var err *JSError = JSError_Make_3(n, func() *DiagnosticType {
		if parent.GetConstructor().IsInterface() {
			return StaticGet_TypeValidator_HIDDEN_INTERFACE_PROPERTY_MISMATCH()
		}
		return StaticGet_TypeValidator_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH()
	}(), []string{propertyName, parent.GetReferenceName(), required.GetJSTypeBase().String(), found.GetJSTypeBase().String(), instance.GetJSTypeBase().String()})
	_ = err
	this_.priv_TypeValidator_registerMismatchAndReport_5666(found, required, err)
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
	if !(typeMap.IsEmpty()) && required.GetJSTypeBase().HasAnyTemplateTypes() {
		required = this_.typeRegistry.ReplaceTemplateTypesPartial(required, typeMap)
	}
	if found.GetJSTypeBase().IsSubtypeOf_2(required, this_.subtypingMode) {
		return
	}
	var err *JSError = JSError_Make_3(n, func() *DiagnosticType {
		if parent.GetConstructor().IsInterface() {
			return StaticGet_TypeValidator_HIDDEN_INTERFACE_PROPERTY_MISMATCH()
		}
		return StaticGet_TypeValidator_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH()
	}(), []string{propertyName.HumanReadableName(), parent.GetReferenceName(), required.GetJSTypeBase().String(), found.GetJSTypeBase().String(), instance.GetJSTypeBase().String()})
	_ = err
	this_.priv_TypeValidator_registerMismatchAndReport_5666(found, required, err)
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
			cur_9 := j2gort.NewIterCursor(j2gort.CheckNotNull(currSuperCtor.GetPrototype().GetObjectTypeBase().GetOwnPropertyKeys()))
			for cur_9.HasNext() {
				rawElem_10 := cur_9.Next()
				prop := j2gort.AsIface[pkg_rhino.Property_Key](rawElem_10)
				_ = prop
				var maybeAbstractMethod *pkg_rhino.FunctionType = superType.GetObjectTypeBase().FindPropertyType_1_Property_Key(prop).ToMaybeFunctionType()
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
			var abstractMethod *pkg_rhino.FunctionType = instance.GetObjectTypeBase().FindPropertyType_1_Property_Key(method).ToMaybeFunctionType()
			_ = abstractMethod
			if (abstractMethod == nil) || abstractMethod.IsAbstract() {
				this_.priv_TypeValidator_registerMismatchAndReport_5666(pkg_rhino.Upcast_JSType(instance), pkg_rhino.Upcast_JSType(superType_2), JSError_Make_3(n, StaticGet_TypeValidator_ABSTRACT_METHOD_NOT_IMPLEMENTED(), []string{method.HumanReadableName(), superType_2.GetObjectTypeBase().String(), instance.GetObjectTypeBase().String()}))
			}
		}
	}
}

// Java method: TypeValidator.mismatch(com.google.javascript.rhino.Node n, java.lang.String msg, com.google.javascript.rhino.jstype.JSType found, com.google.javascript.rhino.jstype.JSTypeNative required)
func (this_ *TypeValidator) priv_TypeValidator_mismatch_5663(n pkg_rhino.Node, msg string, found pkg_rhino.JSType, required pkg_rhino.JSTypeNative) {
	_ = n
	_ = msg
	_ = found
	_ = required
	this_.priv_TypeValidator_mismatch_5664(n, msg, found, this_.priv_TypeValidator_getNativeType_5668(required))
}

// Java method: TypeValidator.mismatch(com.google.javascript.rhino.Node n, java.lang.String msg, com.google.javascript.rhino.jstype.JSType found, com.google.javascript.rhino.jstype.JSType required)
func (this_ *TypeValidator) priv_TypeValidator_mismatch_5664(n pkg_rhino.Node, msg string, found pkg_rhino.JSType, required pkg_rhino.JSType) {
	_ = n
	_ = msg
	_ = found
	_ = required
	if !(found.GetJSTypeBase().IsSubtypeOf_2(required, this_.subtypingMode)) {
		var missing j2gort.Ext_java_util_Set = nil
		_ = missing
		var mismatch j2gort.Ext_java_util_Set = nil
		_ = mismatch
		var structuralType pkg_rhino.JSType = nil
		_ = structuralType
		if required.IsStructuralType() {
			structuralType = required
		} else if required.GetJSTypeBase().IsUnionType() {
			var restrictedType pkg_rhino.JSType = required.RestrictByNotNullOrUndefined()
			_ = restrictedType
			if restrictedType.IsStructuralType() {
				structuralType = restrictedType
			}
		}
		if structuralType != nil && !structuralType.IsJavaNil_() {
			missing = j2gort.NewExt_java_util_TreeSet()
			mismatch = j2gort.NewExt_java_util_TreeSet()
			var requiredObject pkg_rhino.ObjectType = structuralType.GetJSTypeBase().ToMaybeObjectType()
			_ = requiredObject
			var foundObject pkg_rhino.ObjectType = found.GetJSTypeBase().ToMaybeObjectType()
			_ = foundObject
			if (requiredObject != nil && !requiredObject.IsJavaNil_()) && (foundObject != nil && !foundObject.IsJavaNil_()) {
				{
					cur_12 := j2gort.NewIterCursor(j2gort.CheckNotNull(requiredObject.GetObjectTypeBase().GetPropertyNames()))
					for cur_12.HasNext() {
						property := cur_12.NextStr()
						_ = property
						var propRequired pkg_rhino.JSType = requiredObject.GetObjectTypeBase().GetPropertyType_1_string(property)
						_ = propRequired
						var hasProperty bool = foundObject.GetObjectTypeBase().HasProperty_1_string(property)
						_ = hasProperty
						if !(propRequired.IsExplicitlyVoidable()) || hasProperty {
							if hasProperty {
								if !(foundObject.GetObjectTypeBase().GetPropertyType_1_string(property).GetJSTypeBase().IsSubtypeOf_2(propRequired, this_.subtypingMode)) {
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
		this_.priv_TypeValidator_registerMismatchAndReport_5665(n, StaticGet_TypeValidator_TYPE_MISMATCH_WARNING(), msg, found, required, missing, mismatch)
	}
}

// Java method: TypeValidator.registerMismatchAndReport(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.DiagnosticType diagnostic, java.lang.String msg, com.google.javascript.rhino.jstype.JSType found, com.google.javascript.rhino.jstype.JSType required, java.util.@org.jspecify.annotations.Nullable Set<java.lang.String> missing, java.util.@org.jspecify.annotations.Nullable Set<java.lang.String> mismatch)
func (this_ *TypeValidator) priv_TypeValidator_registerMismatchAndReport_5665(n pkg_rhino.Node, diagnostic *DiagnosticType, msg string, found pkg_rhino.JSType, required pkg_rhino.JSType, missing j2gort.Ext_java_util_Set, mismatch j2gort.Ext_java_util_Set) {
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
	this_.priv_TypeValidator_registerMismatchAndReport_5666(found, required, err)
}

// Java method: TypeValidator.registerMismatchAndReport(com.google.javascript.rhino.jstype.JSType found, com.google.javascript.rhino.jstype.JSType required, com.google.javascript.jscomp.JSError error)
func (this_ *TypeValidator) priv_TypeValidator_registerMismatchAndReport_5666(found pkg_rhino.JSType, required pkg_rhino.JSType, error_ *JSError) {
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
	base_found_ := found.GetJSTypeBase()
	var foundStr string = base_found_.String()
	_ = foundStr
	base_required_ := required.GetJSTypeBase()
	var requiredStr string = base_required_.String()
	_ = requiredStr
	if j2gort.CheckNotNullString(foundStr) == requiredStr {
		foundStr = base_found_.ToAnnotationString(pkg_rhino.StaticGet_JSType_Nullability_IMPLICIT())
		requiredStr = base_required_.ToAnnotationString(pkg_rhino.StaticGet_JSType_Nullability_IMPLICIT())
	}
	if j2gort.CheckNotNullString(foundStr) == requiredStr {
		if base_found_.IsEnumElementType() {
			foundStr = (j2gort.StrValString(foundStr) + " (enum definition: " + j2gort.StrValString(TypeValidator_FormatNodeLocation(found.ToMaybeEnumElementType().GetEnumType().GetSource())) + ")")
		}
		if base_required_.IsEnumElementType() {
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
func (this_ *TypeValidator) priv_TypeValidator_getJSType_5667(n pkg_rhino.Node) pkg_rhino.JSType {
	_ = n
	return j2gort.CheckNotNullIface(n.GetJSType())
}

// Java method: TypeValidator.getNativeType(com.google.javascript.rhino.jstype.JSTypeNative typeId)
func (this_ *TypeValidator) priv_TypeValidator_getNativeType_5668(typeId pkg_rhino.JSTypeNative) pkg_rhino.JSType {
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
	var kind int32 = type_.GetJSTypeBase().GetWellFormedTemplatizedKind()
	_ = kind
	if kind != int32(0) {
		this_.priv_TypeValidator_checkWellFormedTemplatizedType_5670(n, type_, kind)
	}
}

// Java method: TypeValidator.checkWellFormedTemplatizedType(com.google.javascript.rhino.Node node, com.google.javascript.rhino.jstype.JSType type, int kind)
func (this_ *TypeValidator) priv_TypeValidator_checkWellFormedTemplatizedType_5670(node pkg_rhino.Node, type_ pkg_rhino.JSType, kind int32) bool {
	_ = node
	_ = type_
	_ = kind
	var ok bool = func() bool {
		switch kind {
		case int32(128):
			var fn *pkg_rhino.FunctionType = j2gort.CastPtr[pkg_rhino.FunctionType](type_)
			_ = fn
			var params []pkg_rhino.FunctionType_Parameter = fn.GetParameterArray()
			_ = params
			var paramsOk bool = true
			_ = paramsOk
			{
				var i int32 = int32(0)
				_ = i
				var len_ int32 = int32(len(params))
				_ = len_
				for ; i < len_; i = i + 1 {
					var paramType pkg_rhino.JSType = params[i].GetJSType()
					_ = paramType
					var paramKind int32 = paramType.GetJSTypeBase().GetWellFormedTemplatizedKind()
					_ = paramKind
					if (paramKind != int32(0)) && !(this_.priv_TypeValidator_checkWellFormedTemplatizedType_5670(node, paramType, paramKind)) {
						paramsOk = false
						break
					}
				}
			}
			if !(paramsOk) {
				return false
			}
			var retType pkg_rhino.JSType = fn.GetReturnType()
			_ = retType
			var retKind int32 = retType.GetJSTypeBase().GetWellFormedTemplatizedKind()
			_ = retKind
			return ((retKind == int32(0)) || this_.priv_TypeValidator_checkWellFormedTemplatizedType_5670(node, retType, retKind))
		case int32(64):
			var union *pkg_rhino.UnionType = j2gort.CastPtr[pkg_rhino.UnionType](type_)
			_ = union
			var alternateCount int32 = union.GetAlternatesCount()
			_ = alternateCount
			var unionOk bool = true
			_ = unionOk
			{
				var i_2 int32 = int32(0)
				_ = i_2
				for ; i_2 < alternateCount; i_2 = i_2 + 1 {
					var alternative pkg_rhino.JSType = union.GetAlternateAt(i_2)
					_ = alternative
					var altKind int32 = alternative.GetJSTypeBase().GetWellFormedTemplatizedKind()
					_ = altKind
					if (altKind != int32(0)) && !(this_.priv_TypeValidator_checkWellFormedTemplatizedType_5670(node, alternative, altKind)) {
						unionOk = false
						break
					}
				}
			}
			return unionOk
		case int32(4096):
			var templatized *pkg_rhino.TemplatizedType = j2gort.CastPtr[pkg_rhino.TemplatizedType](type_)
			_ = templatized
			var referencedTemplates *pkg_rhino.TemplateTypeMap = templatized.GetReferencedType().GetTemplateTypeMap()
			_ = referencedTemplates
			var templateTypes []pkg_rhino.JSType = j2gort.CheckNotNullPtr(templatized).Slice_GetTemplateTypes()
			_ = templateTypes
			var refCount int32 = referencedTemplates.Size()
			_ = refCount
			var templatizedOk bool = true
			_ = templatizedOk
			{
				var i_3 int32 = int32(0)
				_ = i_3
				var len__2 int32 = int32(len(templateTypes))
				_ = len__2
				for ; i_3 < len__2; i_3 = i_3 + 1 {
					var assignedType pkg_rhino.JSType = templateTypes[i_3]
					_ = assignedType
					if i_3 < refCount {
						var templateType *pkg_rhino.TemplateType = referencedTemplates.GetTemplateKeyAt(i_3)
						_ = templateType
						if !(assignedType.GetJSTypeBase().IsSubtype_1(templateType.GetBound())) {
							this_.priv_TypeValidator_registerMismatchAndReport_5666(assignedType, templateType.GetBound(), JSError_Make_3(node, StaticGet_RhinoErrorReporter_BOUNDED_GENERIC_TYPE_ERROR(), []string{assignedType.GetJSTypeBase().String(), templateType.GetReferenceName(), templateType.GetBound().GetJSTypeBase().String()}))
							templatizedOk = false
							break
						}
					}
				}
			}
			return templatizedOk
		case int32(1024):
			var named *pkg_rhino.NamedType = j2gort.CastPtr[pkg_rhino.NamedType](type_)
			_ = named
			var ref pkg_rhino.JSType = named.GetReferencedType()
			_ = ref
			var refKind int32 = ref.GetJSTypeBase().GetWellFormedTemplatizedKind()
			_ = refKind
			return ((refKind == int32(0)) || this_.priv_TypeValidator_checkWellFormedTemplatizedType_5670(node, ref, refKind))
		case int32(8192):
			var templateType_2 *pkg_rhino.TemplateType = j2gort.CastPtr[pkg_rhino.TemplateType](type_)
			_ = templateType_2
			if templateType_2.ContainsCycle() {
				return false
			}
			var bound pkg_rhino.JSType = templateType_2.GetBound()
			_ = bound
			var boundKind int32 = bound.GetJSTypeBase().GetWellFormedTemplatizedKind()
			_ = boundKind
			return ((boundKind == int32(0)) || this_.priv_TypeValidator_checkWellFormedTemplatizedType_5670(node, bound, boundKind))
		case int32(256):
			var enumElement *pkg_rhino.EnumElementType = j2gort.CastPtr[pkg_rhino.EnumElementType](type_)
			_ = enumElement
			var prim pkg_rhino.JSType = enumElement.GetPrimitiveType()
			_ = prim
			if prim == nil || prim.IsJavaNil_() {
				return true
			}
			var primKind int32 = prim.GetJSTypeBase().GetWellFormedTemplatizedKind()
			_ = primKind
			return ((primKind == int32(0)) || this_.priv_TypeValidator_checkWellFormedTemplatizedType_5670(node, prim, primKind))
		default:
			return true
		}
		panic("unreachable switch expression")
	}()
	_ = ok
	if ok {
		type_.GetJSTypeBase().MarkWellFormedTemplatizedOk()
	}
	return ok
}
