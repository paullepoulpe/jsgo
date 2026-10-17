// Port of com/google/javascript/jscomp/InvalidatingTypesTest.java @ closure-compiler v20261015

/*
 * Copyright 2020 The Closure Compiler Authors.
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
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/rhinotesting"
)

func newInvalidatingTypesTest(t *testing.T) (*rhino.JSTypeRegistry, rhinotesting.Truth) {
	return rhino.NewJSTypeRegistry(rhinotesting.NewTestErrorReporter()), rhinotesting.NewTruth(t)
}

func TestInvalidatingTypes_ObjectIsInvalidating(t *testing.T) {
	registry, truth := newInvalidatingTypesTest(t)
	invalidatingTypes := NewInvalidatingTypes_Builder(registry).Build()

	truth.AssertThat(invalidatingTypes.IsInvalidating(registry.GetNativeObjectType(rhino.JSTypeNative_OBJECT_TYPE))).
		IsTrue()
}

func TestInvalidatingTypes_TemplatizedObjectIsInvalidating(t *testing.T) {
	registry, truth := newInvalidatingTypesTest(t)
	invalidatingTypes := NewInvalidatingTypes_Builder(registry).Build()

	objectOfString := registry.CreateTemplatizedType(
		registry.GetNativeObjectType(rhino.JSTypeNative_OBJECT_TYPE), registry.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE))
	truth.AssertThat(invalidatingTypes.IsInvalidating(objectOfString)).IsTrue()
}

func TestInvalidatingTypes_InvaliatingRawTypeAlsoInvalidatesTemplatization(t *testing.T) {
	registry, truth := newInvalidatingTypesTest(t)
	invalidatingTypes := NewInvalidatingTypes_Builder(registry).
		AddAllTypeMismatches(
			[]*TypeMismatch{
				TypeMismatch_CreateForTesting(
					registry.GetNativeObjectType(rhino.JSTypeNative_ARRAY_TYPE),
					registry.GetNativeType(rhino.JSTypeNative_STRING_TYPE)),
			}).
		Build()

	arrayOfNumber := registry.CreateTemplatizedType(
		registry.GetNativeObjectType(rhino.JSTypeNative_ARRAY_TYPE), registry.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE))
	truth.AssertThat(invalidatingTypes.IsInvalidating(arrayOfNumber)).IsTrue()
}

func TestInvalidatingTypes_InvalidatingOneTemplatizationInvalidatesAll(t *testing.T) {
	registry, truth := newInvalidatingTypesTest(t)
	arrayOfString := registry.CreateTemplatizedType(
		registry.GetNativeObjectType(rhino.JSTypeNative_ARRAY_TYPE), registry.GetNativeType(rhino.JSTypeNative_STRING_TYPE))
	invalidatingTypes := NewInvalidatingTypes_Builder(registry).
		AddAllTypeMismatches(
			[]*TypeMismatch{
				TypeMismatch_CreateForTesting(
					arrayOfString, registry.GetNativeType(rhino.JSTypeNative_STRING_TYPE)),
			}).
		Build()

	arrayOfNumber := registry.CreateTemplatizedType(
		registry.GetNativeObjectType(rhino.JSTypeNative_ARRAY_TYPE), registry.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE))
	truth.AssertThat(invalidatingTypes.IsInvalidating(arrayOfNumber)).IsTrue()
}

func TestInvalidatingTypes_InvalidatingOneTemplatizationInvalidatesRawType(t *testing.T) {
	registry, truth := newInvalidatingTypesTest(t)
	arrayOfString := registry.CreateTemplatizedType(
		registry.GetNativeObjectType(rhino.JSTypeNative_ARRAY_TYPE), registry.GetNativeType(rhino.JSTypeNative_STRING_TYPE))
	invalidatingTypes := NewInvalidatingTypes_Builder(registry).
		AddAllTypeMismatches(
			[]*TypeMismatch{
				TypeMismatch_CreateForTesting(
					arrayOfString, registry.GetNativeType(rhino.JSTypeNative_STRING_TYPE)),
			}).
		Build()

	truth.AssertThat(invalidatingTypes.IsInvalidating(registry.GetNativeObjectType(rhino.JSTypeNative_ARRAY_TYPE))).IsTrue()
}

func TestInvalidatingTypes_ArrayNotInvalidating(t *testing.T) {
	registry, truth := newInvalidatingTypesTest(t)
	invalidatingTypes := NewInvalidatingTypes_Builder(registry).Build()

	truth.AssertThat(
		invalidatingTypes.IsInvalidating(
			registry.CreateTemplatizedType(
				registry.GetNativeObjectType(rhino.JSTypeNative_ARRAY_TYPE), registry.GetNativeType(rhino.JSTypeNative_STRING_TYPE)))).
		IsFalse()
}

func TestInvalidatingTypes_InvalidatingInstanceTypeInvalidatesCtor(t *testing.T) {
	registry, truth := newInvalidatingTypesTest(t)
	invalidatingTypes := NewInvalidatingTypes_Builder(registry).
		AddAllTypeMismatches(
			[]*TypeMismatch{
				TypeMismatch_CreateForTesting(
					registry.GetNativeObjectType(rhino.JSTypeNative_ARRAY_TYPE),
					registry.GetNativeType(rhino.JSTypeNative_STRING_TYPE)),
			}).
		Build()

	truth.AssertThat(
		invalidatingTypes.IsInvalidating(
			registry.GetNativeObjectType(rhino.JSTypeNative_ARRAY_TYPE).GetConstructor())).
		IsTrue()
}

func TestInvalidatingTypes_InvalidatingTemplatizedInstanceTypeInvalidatesCtor(t *testing.T) {
	registry, truth := newInvalidatingTypesTest(t)
	arrayOfString := registry.CreateTemplatizedType(
		registry.GetNativeObjectType(rhino.JSTypeNative_ARRAY_TYPE), registry.GetNativeType(rhino.JSTypeNative_STRING_TYPE))
	invalidatingTypes := NewInvalidatingTypes_Builder(registry).
		AddAllTypeMismatches(
			[]*TypeMismatch{
				TypeMismatch_CreateForTesting(
					arrayOfString, registry.GetNativeType(rhino.JSTypeNative_STRING_TYPE)),
			}).
		Build()

	truth.AssertThat(
		invalidatingTypes.IsInvalidating(
			registry.GetNativeObjectType(rhino.JSTypeNative_ARRAY_TYPE).GetConstructor())).
		IsTrue()
}
