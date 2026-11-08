// Port of com/google/javascript/jscomp/colors/StandardColors.java @ closure-compiler v20261104

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

package colors

// StandardColors holds the Colors that are expected to be referencable as
// part of any compilation.
//
// This set describes any Color that needs to be used a priori by a compiler
// pass. That includes:
//
//   - "axiomatic" Colors (those with constant definitions, e.g. UNKNOWN)
//   - primitive Colors (e.g. `number`, `null`)
//   - box Colors (e.g. `Number`, `String`)
//   - Colors used in AstFactory to type synthetic nodes (e.g. `Array`)
func colorIdFromUint64(x uint64) ColorId {
	return ColorId_FromUnsigned(int64(x))
}

// These IDs are randomly chosen.
var (
	StandardColors_ARGUMENTS_ID               = ColorId_FromUnsignedInt(0x1939a66d)
	StandardColors_ARRAY_ID                   = ColorId_FromUnsignedInt(0x79d4a603)
	StandardColors_READONLY_ARRAY_ID          = ColorId_FromUnsignedInt(0x5627ffb0)
	StandardColors_GENERATOR_ID               = ColorId_FromUnsignedInt(int32(-0x644ecfc1)) // 0x9bb1303f
	StandardColors_I_TEMPLATE_ARRAY_ID        = ColorId_FromUnsignedInt(0x46ab3f0e)
	StandardColors_ITERATOR_ID                = ColorId_FromUnsignedInt(0x417ed2ab)
	StandardColors_ASYNC_ITERATOR_ITERABLE_ID = ColorId_FromUnsignedInt(int32(-0x34c7d1f6)) // 0xcb382e0a
	StandardColors_PROMISE_ID                 = ColorId_FromUnsignedInt(0x39581abf)

	// Deterministic ColorIds produced by JSTypeColorIdHasher for global externs collection types.
	StandardColors_ARRAY_PROTOTYPE_ID            = colorIdFromUint64(0x64c23e6d334e56f7)
	StandardColors_READONLY_ARRAY_PROTOTYPE_ID   = colorIdFromUint64(0x3589d35cec22a47f)
	StandardColors_I_TEMPLATE_ARRAY_PROTOTYPE_ID = colorIdFromUint64(0x7a40cd76f33e9198)
	StandardColors_MAP_ID                        = colorIdFromUint64(0x33a2e1724e099aa4)
	StandardColors_MAP_PROTOTYPE_ID              = colorIdFromUint64(0x6d01d1a4e84831d6)
	StandardColors_READONLY_MAP_ID               = colorIdFromUint64(0x03a6d92ab09bc918)
	StandardColors_READONLY_MAP_PROTOTYPE_ID     = colorIdFromUint64(0x158c76decaa1fa9e)
	StandardColors_SET_ID                        = colorIdFromUint64(0x4e48c8722ab53ebd)
	StandardColors_SET_PROTOTYPE_ID              = colorIdFromUint64(0xfe841a0bed231439)
	StandardColors_READONLY_SET_ID               = colorIdFromUint64(0x30860b55852d1f9d)
	StandardColors_READONLY_SET_PROTOTYPE_ID     = colorIdFromUint64(0x8d0f691546b43933)

	StandardColors_BIGINT_OBJECT_ID  = ColorId_FromUnsignedInt(int32(-0x56265293)) // 0xa9d9ad6d
	StandardColors_BOOLEAN_OBJECT_ID = ColorId_FromUnsignedInt(int32(-0x6dfa23fa)) // 0x9205dc06
	StandardColors_NUMBER_OBJECT_ID  = ColorId_FromUnsignedInt(0x34ba2fb1)
	StandardColors_STRING_OBJECT_ID  = ColorId_FromUnsignedInt(0x186008a9)
	StandardColors_SYMBOL_OBJECT_ID  = ColorId_FromUnsignedInt(0x5e514f7e)
)

const (
	axiomaticBigintId      int64 = 0x234eb61a
	axiomaticBooleanId     int64 = 0x126812ee
	axiomaticNullOrVoidId  int64 = 0x22b49f69
	axiomaticNumberId      int64 = 0xd081722c
	axiomaticStringId      int64 = 0x8c4d8f65
	axiomaticSymbolId      int64 = 0x759f2066
	axiomaticTopObjectId   int64 = 0x889b6838
	axiomaticTopFunctionId int64 = 0x283b5838
	axiomaticUnknownId     int64 = 0
	axiomaticGBigintId     int64 = 0xaaae3678
)

var StandardColors_BIGINT = Color_SingleBuilder().
	SetId(ColorId{rightAligned: axiomaticBigintId}).
	SetBoxIdValue(StandardColors_BIGINT_OBJECT_ID).
	SetInvalidating(false).
	buildAxiomatic()

var StandardColors_BOOLEAN = Color_SingleBuilder().
	SetId(ColorId{rightAligned: axiomaticBooleanId}).
	SetBoxIdValue(StandardColors_BOOLEAN_OBJECT_ID).
	SetInvalidating(false).
	buildAxiomatic()

var StandardColors_NULL_OR_VOID = Color_SingleBuilder().
	ClearBoxId().
	SetId(ColorId{rightAligned: axiomaticNullOrVoidId}).
	SetInvalidating(false).
	buildAxiomatic()

var StandardColors_NUMBER = Color_SingleBuilder().
	SetId(ColorId{rightAligned: axiomaticNumberId}).
	SetBoxIdValue(StandardColors_NUMBER_OBJECT_ID).
	SetInvalidating(false).
	buildAxiomatic()

var StandardColors_STRING = Color_SingleBuilder().
	SetId(ColorId{rightAligned: axiomaticStringId}).
	SetBoxIdValue(StandardColors_STRING_OBJECT_ID).
	SetInvalidating(false).
	buildAxiomatic()

var StandardColors_SYMBOL = Color_SingleBuilder().
	SetId(ColorId{rightAligned: axiomaticSymbolId}).
	SetBoxIdValue(StandardColors_SYMBOL_OBJECT_ID).
	SetInvalidating(false).
	buildAxiomatic()

// StandardColors_TOP_OBJECT is the supertype of all objects but not
// primitives.
//
// Separate from UNKNOWN because some optimizations back off on any
// non-object primitives + unknown but operate on the top object.
var StandardColors_TOP_OBJECT = Color_SingleBuilder().
	SetId(ColorId{rightAligned: axiomaticTopObjectId}).
	ClearBoxId().
	SetInvalidating(true).
	buildAxiomatic()

// StandardColors_TOP_FUNCTION is the supertype of all functions.
var StandardColors_TOP_FUNCTION = Color_SingleBuilder().
	SetId(ColorId{rightAligned: axiomaticTopFunctionId}).
	ClearBoxId().
	SetInvalidating(true).
	buildAxiomatic()

// StandardColors_UNKNOWN is analagous to Closure '*'/'?' and TS unknown/any.
var StandardColors_UNKNOWN = Color_SingleBuilder().
	ClearBoxId().
	SetId(ColorId{rightAligned: axiomaticUnknownId}). // Make UNKNOWN the "default" numerical value.
	SetInvalidating(true).
	buildAxiomatic()

// Java: StandardColors.GBIGINT
var StandardColors_GBIGINT = Color_SingleBuilder().
	SetId(ColorId{rightAligned: axiomaticGBigintId}).
	ClearBoxId().
	SetInvalidating(false).
	buildAxiomatic()

func lookupAxiomaticColor(id ColorId) *Color {
	switch id.rightAligned {
	case axiomaticBigintId:
		return StandardColors_BIGINT
	case axiomaticBooleanId:
		return StandardColors_BOOLEAN
	case axiomaticNullOrVoidId:
		return StandardColors_NULL_OR_VOID
	case axiomaticNumberId:
		return StandardColors_NUMBER
	case axiomaticStringId:
		return StandardColors_STRING
	case axiomaticSymbolId:
		return StandardColors_SYMBOL
	case axiomaticTopObjectId:
		return StandardColors_TOP_OBJECT
	case axiomaticTopFunctionId:
		return StandardColors_TOP_FUNCTION
	case axiomaticUnknownId:
		return StandardColors_UNKNOWN
	case axiomaticGBigintId:
		return StandardColors_GBIGINT
	default:
		return nil
	}
}

func isPrimitiveColorId(id ColorId) bool {
	switch id.rightAligned {
	case axiomaticBigintId,
		axiomaticBooleanId,
		axiomaticNullOrVoidId,
		axiomaticNumberId,
		axiomaticStringId,
		axiomaticSymbolId,
		axiomaticGBigintId:
		return true
	default:
		return false
	}
}

// StandardColors_PRIMITIVE_BOX_IDS is the set of ColorIds for object Colors
// that "box" primitive Colors.
//
// Java: StandardColors.PRIMITIVE_BOX_IDS
var StandardColors_PRIMITIVE_BOX_IDS = ColorIdSet{
	StandardColors_BIGINT_OBJECT_ID,
	StandardColors_BOOLEAN_OBJECT_ID,
	StandardColors_NUMBER_OBJECT_ID,
	StandardColors_STRING_OBJECT_ID,
	StandardColors_SYMBOL_OBJECT_ID,
}

// StandardColors_STANDARD_OBJECT_IDS is the set of ColorIds for object Colors
// that need to be referenced in optimizations.
//
// Java: StandardColors.STANDARD_OBJECT_IDS
var StandardColors_STANDARD_OBJECT_IDS = ColorIdSet{
	StandardColors_BIGINT_OBJECT_ID,
	StandardColors_BOOLEAN_OBJECT_ID,
	StandardColors_NUMBER_OBJECT_ID,
	StandardColors_STRING_OBJECT_ID,
	StandardColors_SYMBOL_OBJECT_ID,
	StandardColors_ARRAY_ID,
	StandardColors_ARRAY_PROTOTYPE_ID,
	StandardColors_READONLY_ARRAY_ID,
	StandardColors_READONLY_ARRAY_PROTOTYPE_ID,
	StandardColors_ARGUMENTS_ID,
	StandardColors_ASYNC_ITERATOR_ITERABLE_ID,
	StandardColors_GENERATOR_ID,
	StandardColors_I_TEMPLATE_ARRAY_ID,
	StandardColors_I_TEMPLATE_ARRAY_PROTOTYPE_ID,
	StandardColors_ITERATOR_ID,
	StandardColors_MAP_ID,
	StandardColors_MAP_PROTOTYPE_ID,
	StandardColors_READONLY_MAP_ID,
	StandardColors_READONLY_MAP_PROTOTYPE_ID,
	StandardColors_SET_ID,
	StandardColors_SET_PROTOTYPE_ID,
	StandardColors_READONLY_SET_ID,
	StandardColors_READONLY_SET_PROTOTYPE_ID,
	StandardColors_PROMISE_ID,
}
