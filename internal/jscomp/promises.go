// Port of com/google/javascript/jscomp/Promises.java @ closure-compiler v20261021

/*
 * Copyright 2018 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Promises models different Javascript Promise-related operations
type Promises struct{}

// Promises_GetTemplateTypeOfThenable: if this object is known to be an
// IThenable, returns the type it resolves to.
//
// Returns unknown otherwise.
//
// (This is different from getResolvedType, which will attempt to model the
// then type of an expression after calling Promise.resolve() on it.
//
// Java: Promises.getTemplateTypeOfThenable(JSTypeRegistry, JSType)
func Promises_GetTemplateTypeOfThenable(registry *rhino.JSTypeRegistry, maybeThenable rhino.JSType) rhino.JSType {
	// Without ".restrictByNotNullOrUndefined" we'd get the unknown type for "?IThenable<null>"
	templateType := registry.GetIThenableTemplate()
	return maybeThenable.
		// Without ".restrictByNotNullOrUndefined" we'd get the unknown type for "?IThenable<null>"
		RestrictByNotNullOrUndefined().
		GetTemplateTypeMap().
		GetResolvedTemplateType(templateType)
}

// Promises_GetResolvedType returns the type of `await [expr]`.
//
// This is equivalent to the type of `result` in
// `Promise.resolve([expr]).then(result => `
//
// For example:
//
// !Promise<number> becomes number
//
// !IThenable<number> becomes number
//
// string becomes string
//
// (!Promise<number>|string) becomes (number|string)
//
// ?Promise<number> becomes (null|number)
//
// Java: Promises.getResolvedType(JSTypeRegistry, JSType)
func Promises_GetResolvedType(registry *rhino.JSTypeRegistry, type_ rhino.JSType) rhino.JSType {
	if type_.IsUnknownType() {
		return type_
	}

	if type_.IsUnionType() {
		unionTypeBuilder := rhino.UnionType_Builder_(registry)
		for _, alternate := range type_.ToMaybeUnionType().GetAlternates() {
			unionTypeBuilder.AddAlternate(Promises_GetResolvedType(registry, alternate))
		}
		return unionTypeBuilder.Build()
	}

	// If we can find the "IThenable" template key (which is true for Promise and IThenable), return
	// the resolved value. e.g. for "!Promise<string>" return "string".
	templates := type_.GetTemplateTypeMap()
	if templates.HasTemplateKey(registry.GetIThenableTemplate()) {
		// Call getResolvedPromiseType again in case someone does something unusual like
		// !Promise<!Promise<number>>
		// TODO(lharker): we don't need to handle this case and should report an error for this in a
		// type annotation (not here, maybe in TypeCheck). A Promise cannot resolve to another Promise
		return Promises_GetResolvedType(
			registry, templates.GetResolvedTemplateType(registry.GetIThenableTemplate()))
	}

	// Awaiting anything with a ".then" property (other than IThenable, handled above) should return
	// unknown, rather than the type itself.
	if type_.IsSubtypeOf(registry.GetNativeType(rhino.JSTypeNative_THENABLE_TYPE)) {
		return registry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
	}

	return type_
}

// Promises_WrapInIThenable wraps the given type in an IThenable.
//
// If the given type is already IThenable it is first unwrapped. For example:
//
// number becomes IThenable<number>
//
// IThenable<number> becomes IThenable<number>
//
// Promise<number> becomes IThenable<number>
//
// IThenable<number>|string becomes IThenable<number|string>
//
// IThenable<number>|IThenable<string> becomes IThenable<number|string>
//
// Java: Promises.wrapInIThenable(JSTypeRegistry, JSType)
func Promises_WrapInIThenable(registry *rhino.JSTypeRegistry, maybeThenable rhino.JSType) rhino.JSType {
	// Unwrap for simplicity first in the event it is a thenable.
	unwrapped := Promises_GetResolvedType(registry, maybeThenable)
	return registry.CreateTemplatizedType(
		registry.GetNativeObjectType(rhino.JSTypeNative_I_THENABLE_TYPE), unwrapped)
}

// Promises_CreateAsyncReturnableType synthesizes a type representing the
// legal types of a return expression within async code (i.e.`Promise`
// callbacks, async functions) based on the expected return type of that
// code.
//
// The return type will generally be a union but may not be in the case of
// top-like types. If the expected return type is a union, any synchronous
// elements will be dropped, since they can never occur. For example:
//
//   - `!Promise<number>` => `number|!IThenable<number>`
//   - `number` => `?`
//   - `number|!Promise<string>` => `string|!IThenable<string>`
//   - `!IThenable<number>|!Promise<string>` => `number|string|!IThenable<number|string>`
//   - `!IThenable<number|string>` => `number|string|!IThenable<number|string>`
//   - `?` => `?`
//   - `*` => `?`
//
// Java: Promises.createAsyncReturnableType(JSTypeRegistry, JSType)
func Promises_CreateAsyncReturnableType(registry *rhino.JSTypeRegistry, maybeThenable rhino.JSType) rhino.JSType {
	unknownType := registry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
	iThenableType := registry.GetNativeObjectType(rhino.JSTypeNative_I_THENABLE_TYPE)

	var iThenableOfUnknownType rhino.JSType = registry.CreateTemplatizedType(iThenableType, unknownType)

	var alternates []rhino.JSType
	if maybeThenable.IsUnionType() {
		alternates = maybeThenable.ToMaybeUnionType().GetAlternates()
	} else {
		alternates = []rhino.JSType{maybeThenable}
	}
	asyncTemplateAlternates := []rhino.JSType{}
	for _, t := range alternates {
		if !t.IsSubtypeOf(iThenableOfUnknownType) { // Discard "synchronous" types.
			continue
		}
		// Unwrap "asynchronous" types.
		asyncTemplateAlternates = append(asyncTemplateAlternates, Promises_GetTemplateTypeOfThenable(registry, t))
	}

	if len(asyncTemplateAlternates) == 0 {
		return unknownType
	}

	asyncTemplateUnion := registry.CreateUnionTypeOfList(asyncTemplateAlternates)
	return registry.CreateUnionType(
		asyncTemplateUnion, registry.CreateTemplatizedType(iThenableType, asyncTemplateUnion))
}
