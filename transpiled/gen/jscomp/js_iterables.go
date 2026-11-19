/*
 * Copyright 2014 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/JsIterables.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.JsIterables
type JsIterables struct {
	idHash_ atomic.Int32
}

func (this_ *JsIterables) BindSelf_(s any) {
	_ = s
}
func (this_ *JsIterables) IsJavaNil_() bool { return this_ == nil }
func (*JsIterables) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.JsIterables", "JsIterables"
}
func (this_ *JsIterables) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JsIterables) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.JsIterables", "com.google.javascript.jscomp.JsIterables", func(v_ any) bool {
		return j2gort.InstanceOf[*JsIterables](v_)
	})
}

func (this_ *JsIterables) initFields_JsIterables() {
	_ = this_
}

// Java constructor: JsIterables()
func NewJsIterables() *JsIterables {
	this_ := &JsIterables{}
	InitJsIterables(this_)
	return this_
}

func InitJsIterables(this_ *JsIterables) {
	this_.initFields_JsIterables()
}

// Java static method: JsIterables.getElementType(com.google.javascript.rhino.jstype.JSType iterableOrIterator, com.google.javascript.rhino.jstype.JSTypeRegistry typeRegistry)
func JsIterables_GetElementType(iterableOrIterator pkg_rhino.JSType, typeRegistry *pkg_rhino.JSTypeRegistry) pkg_rhino.JSType {
	_ = iterableOrIterator
	_ = typeRegistry
	var templateTypeMap *pkg_rhino.TemplateTypeMap = iterableOrIterator.Autobox().GetTemplateTypeMap()
	_ = templateTypeMap
	if templateTypeMap.HasTemplateKey_1(typeRegistry.GetIterableValueTemplate()) {
		return templateTypeMap.GetResolvedTemplateType(typeRegistry.GetIterableValueTemplate())
	} else if templateTypeMap.HasTemplateKey_1(typeRegistry.GetIteratorValueTemplate()) {
		return templateTypeMap.GetResolvedTemplateType(typeRegistry.GetIteratorValueTemplate())
	} else if templateTypeMap.HasTemplateKey_1(typeRegistry.GetIteratorLikeValueTemplate()) {
		return templateTypeMap.GetResolvedTemplateType(typeRegistry.GetIteratorLikeValueTemplate())
	} else if templateTypeMap.HasTemplateKey_1(typeRegistry.GetAsyncIterableValueTemplate()) {
		return templateTypeMap.GetResolvedTemplateType(typeRegistry.GetAsyncIterableValueTemplate())
	} else if templateTypeMap.HasTemplateKey_1(typeRegistry.GetAsyncIteratorValueTemplate()) {
		return templateTypeMap.GetResolvedTemplateType(typeRegistry.GetAsyncIteratorValueTemplate())
	}
	var nextProperty pkg_rhino.JSType = iterableOrIterator.FindPropertyType_1_string("next")
	_ = nextProperty
	if nextProperty != nil && !nextProperty.IsJavaNil_() {
		var nextFn *pkg_rhino.FunctionType = nextProperty.ToMaybeFunctionType()
		_ = nextFn
		if (nextFn != nil) && !(j2gort.IfaceEquals[pkg_rhino.JSType](nextFn.GetReturnType(), nil)) {
			var valueProperty pkg_rhino.JSType = nextFn.GetReturnType().FindPropertyType_1_string("value")
			_ = valueProperty
			if valueProperty != nil && !valueProperty.IsJavaNil_() {
				return valueProperty
			}
		}
	}
	return typeRegistry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
}

// Java static method: JsIterables.getReturnElementType(com.google.javascript.rhino.jstype.JSType iterableOrIterator, com.google.javascript.rhino.jstype.JSTypeRegistry typeRegistry)
func JsIterables_GetReturnElementType(iterableOrIterator pkg_rhino.JSType, typeRegistry *pkg_rhino.JSTypeRegistry) pkg_rhino.JSType {
	_ = iterableOrIterator
	_ = typeRegistry
	var templateTypeMap *pkg_rhino.TemplateTypeMap = iterableOrIterator.Autobox().GetTemplateTypeMap()
	_ = templateTypeMap
	if templateTypeMap.HasTemplateKey_1(typeRegistry.GetIterableReturnTemplate()) {
		return templateTypeMap.GetResolvedTemplateType(typeRegistry.GetIterableReturnTemplate())
	} else if templateTypeMap.HasTemplateKey_1(typeRegistry.GetIteratorReturnTemplate()) {
		return templateTypeMap.GetResolvedTemplateType(typeRegistry.GetIteratorReturnTemplate())
	} else if templateTypeMap.HasTemplateKey_1(typeRegistry.GetIteratorLikeReturnTemplate()) {
		return templateTypeMap.GetResolvedTemplateType(typeRegistry.GetIteratorLikeReturnTemplate())
	} else if templateTypeMap.HasTemplateKey_1(typeRegistry.GetAsyncIterableReturnTemplate()) {
		return templateTypeMap.GetResolvedTemplateType(typeRegistry.GetAsyncIterableReturnTemplate())
	} else if templateTypeMap.HasTemplateKey_1(typeRegistry.GetAsyncIteratorReturnTemplate()) {
		return templateTypeMap.GetResolvedTemplateType(typeRegistry.GetAsyncIteratorReturnTemplate())
	}
	return typeRegistry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
}

// Java static method: JsIterables.createIterableTypeOf(com.google.javascript.rhino.jstype.JSType elementType, com.google.javascript.rhino.jstype.JSTypeRegistry typeRegistry)
func JsIterables_CreateIterableTypeOf(elementType pkg_rhino.JSType, typeRegistry *pkg_rhino.JSTypeRegistry) pkg_rhino.JSType {
	_ = elementType
	_ = typeRegistry
	return pkg_rhino.ToIface_JSType(typeRegistry.CreateTemplatizedType_2_ObjectType_JSType(typeRegistry.GetNativeObjectType(pkg_rhino.StaticGet_JSTypeNative_ITERABLE_TYPE()), []pkg_rhino.JSType{elementType}))
}

// Java static method: JsIterables.maybeBoxIterable(com.google.javascript.rhino.jstype.JSType type, com.google.javascript.rhino.jstype.JSTypeRegistry typeRegistry)
func JsIterables_MaybeBoxIterable(type_ pkg_rhino.JSType, typeRegistry *pkg_rhino.JSTypeRegistry) *JsIterables_MaybeBoxedType {
	_ = type_
	_ = typeRegistry
	return JsIterables_MaybeBoxIterableAndPossiblyAsyncIterable(type_, typeRegistry, false)
}

// Java static method: JsIterables.maybeBoxIterableOrAsyncIterable(com.google.javascript.rhino.jstype.JSType type, com.google.javascript.rhino.jstype.JSTypeRegistry typeRegistry)
func JsIterables_MaybeBoxIterableOrAsyncIterable(type_ pkg_rhino.JSType, typeRegistry *pkg_rhino.JSTypeRegistry) *JsIterables_MaybeBoxedType {
	_ = type_
	_ = typeRegistry
	return JsIterables_MaybeBoxIterableAndPossiblyAsyncIterable(type_, typeRegistry, true)
}

// Java static method: JsIterables.maybeBoxIterableAndPossiblyAsyncIterable(com.google.javascript.rhino.jstype.JSType type, com.google.javascript.rhino.jstype.JSTypeRegistry typeRegistry, boolean includeAsyncIterable)
func JsIterables_MaybeBoxIterableAndPossiblyAsyncIterable(type_ pkg_rhino.JSType, typeRegistry *pkg_rhino.JSTypeRegistry, includeAsyncIterable bool) *JsIterables_MaybeBoxedType {
	_ = type_
	_ = typeRegistry
	_ = includeAsyncIterable
	var templatedTypes j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
	_ = templatedTypes
	if type_.IsUnionType() {
		var union *pkg_rhino.UnionType = type_.ToMaybeUnionType()
		_ = union
		{
			var i int32 = int32(0)
			_ = i
			var n int32 = union.GetAlternatesCount()
			_ = n
			for ; i < n; i = i + 1 {
				var alt pkg_rhino.JSType = union.GetAlternateAt(i)
				_ = alt
				if alt.IsBoxableScalar() {
					alt = alt.AutoboxesTo()
				} else {
					alt = alt
				}
				var isIterable bool = alt.IsSubtypeOf_1(typeRegistry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_ITERABLE_TYPE()))
				_ = isIterable
				var isAsyncIterable bool = (includeAsyncIterable && alt.IsSubtypeOf_1(typeRegistry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_ASYNC_ITERABLE_TYPE())))
				_ = isAsyncIterable
				if !(isIterable) && !(isAsyncIterable) {
					return NewJsIterables_MaybeBoxedType(nil, alt)
				}
				var valueTemplate *pkg_rhino.TemplateType
				_ = valueTemplate
				if isAsyncIterable {
					valueTemplate = typeRegistry.GetAsyncIterableValueTemplate()
				} else {
					valueTemplate = typeRegistry.GetIterableValueTemplate()
				}
				templatedTypes.Add_1(alt.GetTemplateTypeMap().GetResolvedTemplateType(valueTemplate))
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
		var isIterable_2 bool = autoboxedType.IsSubtypeOf_1(typeRegistry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_ITERABLE_TYPE()))
		_ = isIterable_2
		var isAsyncIterable_2 bool = (includeAsyncIterable && autoboxedType.IsSubtypeOf_1(typeRegistry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_ASYNC_ITERABLE_TYPE())))
		_ = isAsyncIterable_2
		if !(isIterable_2) && !(isAsyncIterable_2) {
			return NewJsIterables_MaybeBoxedType(nil, autoboxedType)
		}
		var templateType *pkg_rhino.TemplateType
		_ = templateType
		if isAsyncIterable_2 {
			templateType = typeRegistry.GetAsyncIterableValueTemplate()
		} else {
			templateType = typeRegistry.GetIterableValueTemplate()
		}
		templatedTypes.Add_1(autoboxedType.GetTemplateTypeMap().GetResolvedTemplateType(templateType))
	}
	return NewJsIterables_MaybeBoxedType(typeRegistry.CreateUnionType_1_Ext_java_util_List(templatedTypes), nil)
}

// Java class: com.google.javascript.jscomp.JsIterables.MaybeBoxedType
type JsIterables_MaybeBoxedType struct {
	templatedType pkg_rhino.JSType `gson:"templatedType"`
	mismatchType  pkg_rhino.JSType `gson:"mismatchType"`
	idHash_       atomic.Int32
}

func (this_ *JsIterables_MaybeBoxedType) BindSelf_(s any) {
	_ = s
}
func (this_ *JsIterables_MaybeBoxedType) IsJavaNil_() bool { return this_ == nil }
func (*JsIterables_MaybeBoxedType) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.JsIterables$MaybeBoxedType", "MaybeBoxedType"
}
func (this_ *JsIterables_MaybeBoxedType) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JsIterables_MaybeBoxedType) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.JsIterables.MaybeBoxedType", "com.google.javascript.jscomp.JsIterables$MaybeBoxedType", func(v_ any) bool {
		return j2gort.InstanceOf[*JsIterables_MaybeBoxedType](v_)
	})
}

func (this_ *JsIterables_MaybeBoxedType) initFields_JsIterables_MaybeBoxedType() {
	_ = this_
}

// Java constructor: MaybeBoxedType(com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType templatedType, com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType mismatchType)
func NewJsIterables_MaybeBoxedType(templatedType pkg_rhino.JSType, mismatchType pkg_rhino.JSType) *JsIterables_MaybeBoxedType {
	this_ := &JsIterables_MaybeBoxedType{}
	InitJsIterables_MaybeBoxedType(this_, templatedType, mismatchType)
	return this_
}

func InitJsIterables_MaybeBoxedType(this_ *JsIterables_MaybeBoxedType, templatedType pkg_rhino.JSType, mismatchType pkg_rhino.JSType) {
	_ = templatedType
	_ = mismatchType
	this_.initFields_JsIterables_MaybeBoxedType()
	this_.templatedType = templatedType
	this_.mismatchType = mismatchType
}

// Java method: MaybeBoxedType.getTemplatedType()
func (this_ *JsIterables_MaybeBoxedType) GetTemplatedType() pkg_rhino.JSType {
	if !(this_.IsMatch()) {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Type was not boxable to iterable or async iterable!")))
	}
	return this_.templatedType
}

// Java method: MaybeBoxedType.getMismatchType()
func (this_ *JsIterables_MaybeBoxedType) GetMismatchType() pkg_rhino.JSType {
	return this_.mismatchType
}

// Java method: MaybeBoxedType.orElse(com.google.javascript.rhino.jstype.JSType type)
func (this_ *JsIterables_MaybeBoxedType) OrElse(type_ pkg_rhino.JSType) pkg_rhino.JSType {
	_ = type_
	if this_.IsMatch() {
		return this_.templatedType
	}
	return type_
}

// Java method: MaybeBoxedType.isMatch()
func (this_ *JsIterables_MaybeBoxedType) IsMatch() bool {
	return this_.templatedType != nil && !this_.templatedType.IsJavaNil_()
}
