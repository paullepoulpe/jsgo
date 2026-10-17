// Port of com/google/javascript/jscomp/JsIterables.java @ closure-compiler v20261015

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

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// JsIterables models type transformations of JavaScript `Iterable` and
// `Iterator` types.
//
// These tranformations can be especially helpful when working with
// generator functions. `Generator` is a subtype of `Iterable`.
type JsIterables struct{}

// JsIterables_GetElementType returns the given `Iterable`s element type.
//
// If the given type is not an `Iterator`, `Iterable`, `AsyncIterator`, or
// `AsyncIterable`, returns the unknown type.
//
// Java: JsIterables.getElementType(JSType, JSTypeRegistry)
func JsIterables_GetElementType(iterableOrIterator rhino.JSType, typeRegistry *rhino.JSTypeRegistry) rhino.JSType {
	templateTypeMap := iterableOrIterator.
		// Remember that `string` will box to a `Iterable`.
		Autobox().
		GetTemplateTypeMap()

	if templateTypeMap.HasTemplateKey(typeRegistry.GetIterableValueTemplate()) {
		// `Iterable<SomeElementType>` or `Generator<SomeElementType>`
		return templateTypeMap.GetResolvedTemplateType(typeRegistry.GetIterableValueTemplate())
	} else if templateTypeMap.HasTemplateKey(typeRegistry.GetIteratorValueTemplate()) {
		// `Iterator<SomeElementType>`
		return templateTypeMap.GetResolvedTemplateType(typeRegistry.GetIteratorValueTemplate())
	} else if templateTypeMap.HasTemplateKey(typeRegistry.GetIteratorLikeValueTemplate()) {
		// `IteratorLike<SomeElementType>`
		return templateTypeMap.GetResolvedTemplateType(typeRegistry.GetIteratorLikeValueTemplate())
	} else if templateTypeMap.HasTemplateKey(typeRegistry.GetAsyncIterableValueTemplate()) {
		// `AsyncIterable<SomeElementType>` or `AsyncGenerator<SomeElementType>`
		return templateTypeMap.GetResolvedTemplateType(typeRegistry.GetAsyncIterableValueTemplate())
	} else if templateTypeMap.HasTemplateKey(typeRegistry.GetAsyncIteratorValueTemplate()) {
		// `AsyncIterator<SomeElementType>`
		return templateTypeMap.GetResolvedTemplateType(typeRegistry.GetAsyncIteratorValueTemplate())
	}
	nextProperty := iterableOrIterator.FindPropertyType("next")
	if nextProperty != nil {
		nextFn := nextProperty.ToMaybeFunctionType()
		if nextFn != nil && nextFn.GetReturnType() != nil {
			valueProperty := nextFn.GetReturnType().FindPropertyType("value")
			if valueProperty != nil {
				return valueProperty
			}
		}
	}
	return typeRegistry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
}

// JsIterables_GetReturnElementType returns the given `Iterable`s return
// type.
//
// This corresponds to the `TReturn` template variable in
// Iterable<T, TReturn, TNext>.
//
// If the given type is not an `Iterator`, `Iterable`, `AsyncIterator`, or
// `AsyncIterable`, returns the unknown type.
//
// Java: JsIterables.getReturnElementType(JSType, JSTypeRegistry)
func JsIterables_GetReturnElementType(iterableOrIterator rhino.JSType, typeRegistry *rhino.JSTypeRegistry) rhino.JSType {
	templateTypeMap := iterableOrIterator.
		// Remember that `string` will box to a `Iterable`.
		Autobox().
		GetTemplateTypeMap()

	if templateTypeMap.HasTemplateKey(typeRegistry.GetIterableReturnTemplate()) {
		// `Iterable<?, SomeElementType>` or `Generator<?, SomeElementType>`
		return templateTypeMap.GetResolvedTemplateType(typeRegistry.GetIterableReturnTemplate())
	} else if templateTypeMap.HasTemplateKey(typeRegistry.GetIteratorReturnTemplate()) {
		// `Iterator<?, SomeElementType>`
		return templateTypeMap.GetResolvedTemplateType(typeRegistry.GetIteratorReturnTemplate())
	} else if templateTypeMap.HasTemplateKey(typeRegistry.GetIteratorLikeReturnTemplate()) {
		// `IteratorLike<?, SomeElementType>`
		return templateTypeMap.GetResolvedTemplateType(typeRegistry.GetIteratorLikeReturnTemplate())
	} else if templateTypeMap.HasTemplateKey(typeRegistry.GetAsyncIterableReturnTemplate()) {
		// `AsyncIterable<?, SomeElementType>` or `AsyncGenerator<?, SomeElementType>`
		return templateTypeMap.GetResolvedTemplateType(typeRegistry.GetAsyncIterableReturnTemplate())
	} else if templateTypeMap.HasTemplateKey(typeRegistry.GetAsyncIteratorReturnTemplate()) {
		// `AsyncIterator<?, SomeElementType>`
		return templateTypeMap.GetResolvedTemplateType(typeRegistry.GetAsyncIteratorReturnTemplate())
	}
	return typeRegistry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
}

// JsIterables_CreateIterableTypeOf returns an `Iterable` type templated on
// elementType.
//
// Example: `number' => `Iterable<number>`.
//
// Java: JsIterables.createIterableTypeOf(JSType, JSTypeRegistry)
func JsIterables_CreateIterableTypeOf(elementType rhino.JSType, typeRegistry *rhino.JSTypeRegistry) rhino.JSType {
	return typeRegistry.CreateTemplatizedType(
		typeRegistry.GetNativeObjectType(rhino.JSTypeNative_ITERABLE_TYPE), elementType)
}

// JsIterables_MaybeBoxedType is the template of an Iterable|AsyncIterable,
// or a mismatch if subtype of Iterable|AsyncIterable.
type JsIterables_MaybeBoxedType struct {
	templatedType rhino.JSType
	mismatchType  rhino.JSType
}

// newJsIterables_MaybeBoxedType: templatedType is the unwrapped type that is
// being yieled by the iterable / async iterable. e.g. in
// Iterable<number>|AsyncIterable<Qux>|string, this should be
// number|Qux|string. Null if a mismatch. mismatchType is the type that
// caused the mismatch. This can be the entire type (e.g. number is not
// iterable or async iterable) or a piece of a union that caused the mismatch
// (e.g. number in number|Iterable<Qux>). Null if a match.
//
// Java: JsIterables.MaybeBoxedType(JSType, JSType)
func newJsIterables_MaybeBoxedType(templatedType rhino.JSType, mismatchType rhino.JSType) *JsIterables_MaybeBoxedType {
	m := &JsIterables_MaybeBoxedType{}
	m.templatedType = templatedType
	m.mismatchType = mismatchType
	return m
}

// GetTemplatedType returns the template of the iterable / async iterable;
// throws if not an iterable / async iterable.
//
// Java: JsIterables.MaybeBoxedType.getTemplatedType()
func (m *JsIterables_MaybeBoxedType) GetTemplatedType() rhino.JSType {
	if !m.IsMatch() {
		panic(jutil.NewIllegalStateExceptionMsg("Type was not boxable to iterable or async iterable!"))
	}
	return m.templatedType
}

// GetMismatchType returns the type that caused the mismatch, if any
//
// Java: JsIterables.MaybeBoxedType.getMismatchType()
func (m *JsIterables_MaybeBoxedType) GetMismatchType() rhino.JSType {
	return m.mismatchType
}

// Java: JsIterables.MaybeBoxedType.orElse(JSType)
func (m *JsIterables_MaybeBoxedType) OrElse(type_ rhino.JSType) rhino.JSType {
	if m.IsMatch() {
		return m.templatedType
	}

	return type_
}

// Java: JsIterables.MaybeBoxedType.isMatch()
func (m *JsIterables_MaybeBoxedType) IsMatch() bool {
	return m.templatedType != nil
}

// JsIterables_MaybeBoxIterable returns the template for iterable types.
//
// This does not apply to async-iterables.
//
// If this value is not an Iterable, we return a sentinel.
//
// Java: JsIterables.maybeBoxIterable(JSType, JSTypeRegistry)
func JsIterables_MaybeBoxIterable(type_ rhino.JSType, typeRegistry *rhino.JSTypeRegistry) *JsIterables_MaybeBoxedType {
	return JsIterables_MaybeBoxIterableAndPossiblyAsyncIterable(
		type_, typeRegistry /* includeAsyncIterable= */, false)
}

// JsIterables_MaybeBoxIterableOrAsyncIterable returns the template for
// iterable or async-iterable types.
//
// If this value is not an Iterable or AsyncIterable, we return a sentinel.
//
// Java: JsIterables.maybeBoxIterableOrAsyncIterable(JSType, JSTypeRegistry)
func JsIterables_MaybeBoxIterableOrAsyncIterable(
	type_ rhino.JSType, typeRegistry *rhino.JSTypeRegistry) *JsIterables_MaybeBoxedType {
	return JsIterables_MaybeBoxIterableAndPossiblyAsyncIterable(
		type_, typeRegistry /* includeAsyncIterable= */, true)
}

// JsIterables_MaybeBoxIterableAndPossiblyAsyncIterable: given a type, if it
// is an iterable or async iterable, will return its template. If not a
// subtype of Iterable|AsyncIterable, returns an object that has no match,
// and will indicate the mismatch. e.g. both number and number|Iterable are
// not subtypes of Iterable|AsyncIterable.
//
// includeAsyncIterable is whether to look for AsyncIterable as opposed to
// just Iterable.
//
// Java: JsIterables.maybeBoxIterableAndPossiblyAsyncIterable(JSType, JSTypeRegistry, boolean)
func JsIterables_MaybeBoxIterableAndPossiblyAsyncIterable(
	type_ rhino.JSType, typeRegistry *rhino.JSTypeRegistry, includeAsyncIterable bool) *JsIterables_MaybeBoxedType {
	templatedTypes := []rhino.JSType{}

	// Note: we don't just use JSType.autobox() here because that removes null and undefined.
	// We want to keep null and undefined around because they should cause a mismatch.
	if type_.IsUnionType() {
		for _, alt := range type_.ToMaybeUnionType().GetAlternates() {
			if alt.IsBoxableScalar() {
				alt = alt.AutoboxesTo()
			}
			isIterable := alt.IsSubtypeOf(typeRegistry.GetNativeType(rhino.JSTypeNative_ITERABLE_TYPE))
			isAsyncIterable := includeAsyncIterable &&
				alt.IsSubtypeOf(typeRegistry.GetNativeType(rhino.JSTypeNative_ASYNC_ITERABLE_TYPE))
			if !isIterable && !isAsyncIterable {
				return newJsIterables_MaybeBoxedType(nil, alt)
			}
			var valueTemplate *rhino.TemplateType
			if isAsyncIterable {
				valueTemplate = typeRegistry.GetAsyncIterableValueTemplate()
			} else {
				valueTemplate = typeRegistry.GetIterableValueTemplate()
			}
			templatedTypes = append(templatedTypes, alt.GetTemplateTypeMap().GetResolvedTemplateType(valueTemplate))
		}
	} else {
		autoboxedType := type_
		if type_.IsBoxableScalar() {
			autoboxedType = type_.AutoboxesTo()
		}
		isIterable := autoboxedType.IsSubtypeOf(typeRegistry.GetNativeType(rhino.JSTypeNative_ITERABLE_TYPE))
		isAsyncIterable := includeAsyncIterable &&
			autoboxedType.IsSubtypeOf(typeRegistry.GetNativeType(rhino.JSTypeNative_ASYNC_ITERABLE_TYPE))
		if !isIterable && !isAsyncIterable {
			return newJsIterables_MaybeBoxedType(nil, autoboxedType)
		}
		var templateType *rhino.TemplateType
		if isAsyncIterable {
			templateType = typeRegistry.GetAsyncIterableValueTemplate()
		} else {
			templateType = typeRegistry.GetIterableValueTemplate()
		}
		templatedTypes = append(templatedTypes, autoboxedType.GetTemplateTypeMap().GetResolvedTemplateType(templateType))
	}
	return newJsIterables_MaybeBoxedType(typeRegistry.CreateUnionTypeOfList(templatedTypes), nil)
}

// Java: JsIterables()
func newJsIterables() *JsIterables { return &JsIterables{} }
