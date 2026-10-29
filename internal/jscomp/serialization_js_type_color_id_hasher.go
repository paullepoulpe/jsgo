// Port of com/google/javascript/jscomp/serialization/JSTypeColorIdHasher.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// JSTypeColorIdHasher computes the ColorId of an ObjectType.
type JSTypeColorIdHasher struct {
	registry *rhino.JSTypeRegistry
	// Go-only: language-adapter — an ImmutableMap keyed by ObjectType, whose
	// equals/hashCode are JSType.equals/hashCode, so it is an
	// ObjectLinkedHashMap rather than a Go-== keyed map.
	standardColorObjectIds *jutil.ObjectLinkedHashMap[rhino.ObjectType, colors.ColorId]
	standardColorByPtr     map[*rhino.JSTypeBase]colors.ColorId
	hasher                 jutil.FarmHashFingerprint64Hasher
	hashOwnPropFn          func(string)
}

func (h *JSTypeColorIdHasher) onOwnProp(prop string) {
	JSTypeColorIdHasher_PutUtf8(&h.hasher, prop)
}

// Java: JSTypeColorIdHasher(JSTypeRegistry)
func NewJSTypeColorIdHasher(registry *rhino.JSTypeRegistry) *JSTypeColorIdHasher {
	if registry == nil {
		panic(jutil.NewNullPointerException())
	}
	return &JSTypeColorIdHasher{registry: registry}
}

func (h *JSTypeColorIdHasher) initStandardColorByPtr() {
	standardColorByPtr := make(map[*rhino.JSTypeBase]colors.ColorId, len(jsTypeColorIdHasherNativeEntries))
	for _, entry := range jsTypeColorIdHasherNativeEntries {
		k := h.registry.GetNativeObjectType(entry.native)
		standardColorByPtr[k.GetJSTypeBase()] = entry.id
	}
	h.standardColorByPtr = standardColorByPtr
}

func (h *JSTypeColorIdHasher) initStandardColorObjectIds() {
	standardColorObjectIds := jutil.NewObjectLinkedHashMap[rhino.ObjectType, colors.ColorId]()
	for _, entry := range jsTypeColorIdHasherNativeEntries {
		k := h.registry.GetNativeObjectType(entry.native)
		if old, ok := standardColorObjectIds.Lookup(k); ok {
			panic(jutil.NewIllegalArgumentExceptionMsg("Multiple entries with same key: " +
				k.String() + "=" + entry.id.String() + " and " + k.String() + "=" + old.String()))
		}
		standardColorObjectIds.Put(k, entry.id)
	}
	h.standardColorObjectIds = standardColorObjectIds
}

// Java: JSTypeColorIdHasher.hashObjectType(ObjectType)
func (h *JSTypeColorIdHasher) HashObjectType(type_ rhino.ObjectType) colors.ColorId {
	jutil.CheckStateMsg1(
		type_ != nil &&
			!type_.IsEnumElementType() &&
			!type_.IsNoResolvedType() &&
			!type_.IsTemplatizedType() &&
			!type_.IsUnknownType() &&
			!type_.IsTemplateType(),
		"%s",
		type_)

	if type_.IsNominalType() {
		if h.standardColorByPtr == nil {
			h.initStandardColorByPtr()
		}
		if boxId, ok := h.standardColorByPtr[type_.GetJSTypeBase()]; ok {
			return boxId
		}
		if type_.IsNamedType() {
			if h.standardColorObjectIds == nil {
				h.initStandardColorObjectIds()
			}
			if boxId, ok := h.standardColorObjectIds.Lookup(type_); ok {
				return boxId
			}
		}
	}

	hasher := &h.hasher
	hasher.Reset()

	if type_.IsFunctionType() {
		fnType := type_.ToMaybeFunctionType()
		if fnType.HasInstanceType() {
			hasher.PutInt(jsTypeColorIdHasher_Marker_HAS_INSTANCE_TYPE)
		}
	}

	if type_.HasReferenceName() {
		h.PutReferenceNameType(hasher, type_)
	} else {
		if h.hashOwnPropFn == nil {
			h.hashOwnPropFn = h.onOwnProp
		}
		type_.ForEachOwnPropertyName(h.hashOwnPropFn)
	}

	return colors.ColorId_FromUnsigned(hasher.HashAsLong())
}

// Java: JSTypeColorIdHasher.putReferenceNameType(Hasher, ObjectType)
func (h *JSTypeColorIdHasher) PutReferenceNameType(hasher *jutil.FarmHashFingerprint64Hasher, type_ rhino.ObjectType) {
	JSTypeColorIdHasher_PutUtf8(hasher, *type_.GetReferenceName())

	sourceRef := JSTypeColorIdHasher_SourceRefFor(type_)
	if sourceRef == nil {
		// TODO(b/185519307): This is a hack to work around bugs in the typesystem.
		hasher.PutInt(jsTypeColorIdHasher_Marker_UNKNOWN_SOURCEREF)
		return
	}

	moduleId := sourceRef.GetGoogModuleId()
	if moduleId == nil {
		hasher.PutInt(jsTypeColorIdHasher_Marker_NO_GOOG_MODULE_ID)
	} else {
		JSTypeColorIdHasher_PutUtf8(hasher, *moduleId)
	}
}

// Java: JSTypeColorIdHasher.sourceRefFor(ObjectType)
func JSTypeColorIdHasher_SourceRefFor(type_ rhino.ObjectType) rhino.JSType_WithSourceRef {
	if type_.IsEnumType() {
		return type_.ToMaybeEnumType()
	} else if type_.IsEnumElementType() {
		// EnumElementTypes are proxied to their underlying type when transformed to colors.
		panic(jutil.NewAssertionError(type_.String()))
	}

	if type_.IsFunctionType() {
		return type_.ToMaybeFunctionType()
	} else if type_.IsFunctionPrototypeType() {
		return type_.GetOwnerFunction()
	} else if type_.GetConstructor() != nil {
		return type_.GetConstructor()
	}

	return nil
}

// Java: JSTypeColorIdHasher.FARM_64 (newHasher())
var jsTypeColorIdHasher_FARM_64 = jutil.Hashing_FarmHashFingerprint64NewHasher

// Java: JSTypeColorIdHasher.putUtf8(Hasher, String)
func JSTypeColorIdHasher_PutUtf8(hasher *jutil.FarmHashFingerprint64Hasher, s string) {
	hasher.PutStringUTF8(s)
}

// Arbitrary constants that may get hashed into ColorId.
//
// Since a hash may end up including combinations of markers, using int
// values reduces that chance that two combinations will concatenate to the
// same bit sequence.
//
// Java: JSTypeColorIdHasher.Marker
const (
	jsTypeColorIdHasher_Marker_HAS_INSTANCE_TYPE int32 = -0x73748f25 // 0x8c8b70db
	jsTypeColorIdHasher_Marker_NO_GOOG_MODULE_ID int32 = 0x2593c5ff
	jsTypeColorIdHasher_Marker_UNKNOWN_SOURCEREF int32 = 0x660be782
)

const jsTypeNativeCount = int(rhino.JSTypeNative_NUMBER_ADDITION_SUPERTYPE) + 1

type jsTypeNativeColorIdEntry struct {
	native rhino.JSTypeNative
	id     colors.ColorId
}

var jsTypeColorIdHasherNativeEntries = [...]jsTypeNativeColorIdEntry{
	{rhino.JSTypeNative_ARGUMENTS_TYPE, colors.StandardColors_ARGUMENTS_ID},
	{rhino.JSTypeNative_ARRAY_TYPE, colors.StandardColors_ARRAY_ID},
	{rhino.JSTypeNative_READONLY_ARRAY_TYPE, colors.StandardColors_READONLY_ARRAY_ID},
	{rhino.JSTypeNative_ASYNC_ITERATOR_ITERABLE_TYPE, colors.StandardColors_ASYNC_ITERATOR_ITERABLE_ID},
	{rhino.JSTypeNative_GENERATOR_TYPE, colors.StandardColors_GENERATOR_ID},
	{rhino.JSTypeNative_I_TEMPLATE_ARRAY_TYPE, colors.StandardColors_I_TEMPLATE_ARRAY_ID},
	{rhino.JSTypeNative_ITERATOR_TYPE, colors.StandardColors_ITERATOR_ID},
	{rhino.JSTypeNative_PROMISE_TYPE, colors.StandardColors_PROMISE_ID},
	{rhino.JSTypeNative_BIGINT_OBJECT_TYPE, colors.StandardColors_BIGINT_OBJECT_ID},
	{rhino.JSTypeNative_BOOLEAN_OBJECT_TYPE, colors.StandardColors_BOOLEAN_OBJECT_ID},
	{rhino.JSTypeNative_NUMBER_OBJECT_TYPE, colors.StandardColors_NUMBER_OBJECT_ID},
	{rhino.JSTypeNative_STRING_OBJECT_TYPE, colors.StandardColors_STRING_OBJECT_ID},
	{rhino.JSTypeNative_SYMBOL_OBJECT_TYPE, colors.StandardColors_SYMBOL_OBJECT_ID},
}

// Java: JSTypeColorIdHasher.NATIVE_TYPE_TO_ID
var JSTypeColorIdHasher_NATIVE_TYPE_TO_ID = [jsTypeNativeCount]colors.ColorId{
	rhino.JSTypeNative_ARGUMENTS_TYPE:               colors.StandardColors_ARGUMENTS_ID,
	rhino.JSTypeNative_ARRAY_TYPE:                   colors.StandardColors_ARRAY_ID,
	rhino.JSTypeNative_READONLY_ARRAY_TYPE:          colors.StandardColors_READONLY_ARRAY_ID,
	rhino.JSTypeNative_ASYNC_ITERATOR_ITERABLE_TYPE: colors.StandardColors_ASYNC_ITERATOR_ITERABLE_ID,
	rhino.JSTypeNative_GENERATOR_TYPE:               colors.StandardColors_GENERATOR_ID,
	rhino.JSTypeNative_I_TEMPLATE_ARRAY_TYPE:        colors.StandardColors_I_TEMPLATE_ARRAY_ID,
	rhino.JSTypeNative_ITERATOR_TYPE:                colors.StandardColors_ITERATOR_ID,
	rhino.JSTypeNative_PROMISE_TYPE:                 colors.StandardColors_PROMISE_ID,
	rhino.JSTypeNative_BIGINT_OBJECT_TYPE:           colors.StandardColors_BIGINT_OBJECT_ID,
	rhino.JSTypeNative_BOOLEAN_OBJECT_TYPE:          colors.StandardColors_BOOLEAN_OBJECT_ID,
	rhino.JSTypeNative_NUMBER_OBJECT_TYPE:           colors.StandardColors_NUMBER_OBJECT_ID,
	rhino.JSTypeNative_STRING_OBJECT_TYPE:           colors.StandardColors_STRING_OBJECT_ID,
	rhino.JSTypeNative_SYMBOL_OBJECT_TYPE:           colors.StandardColors_SYMBOL_OBJECT_ID,
}
