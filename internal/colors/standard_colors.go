// Port of com/google/javascript/jscomp/colors/StandardColors.java @ closure-compiler v20261015

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

import "github.com/paullepoulpe/jsgo/internal/jutil"

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
type StandardColors struct{}

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

	StandardColors_BIGINT_OBJECT_ID  = ColorId_FromUnsignedInt(int32(-0x56265293)) // 0xa9d9ad6d
	StandardColors_BOOLEAN_OBJECT_ID = ColorId_FromUnsignedInt(int32(-0x6dfa23fa)) // 0x9205dc06
	StandardColors_NUMBER_OBJECT_ID  = ColorId_FromUnsignedInt(0x34ba2fb1)
	StandardColors_STRING_OBJECT_ID  = ColorId_FromUnsignedInt(0x186008a9)
	StandardColors_SYMBOL_OBJECT_ID  = ColorId_FromUnsignedInt(0x5e514f7e)
)

func colorIdPtr(id ColorId) *ColorId { return &id }

var StandardColors_BIGINT = Color_SingleBuilder().
	SetId(ColorId_FromUnsignedInt(0x234eb61a)).
	SetBoxId(colorIdPtr(StandardColors_BIGINT_OBJECT_ID)).
	SetInvalidating(false).
	BuildAxiomatic()

var StandardColors_BOOLEAN = Color_SingleBuilder().
	SetId(ColorId_FromUnsignedInt(0x126812ee)).
	SetBoxId(colorIdPtr(StandardColors_BOOLEAN_OBJECT_ID)).
	SetInvalidating(false).
	BuildAxiomatic()

var StandardColors_NULL_OR_VOID = Color_SingleBuilder().
	SetBoxId(nil).
	SetId(ColorId_FromUnsignedInt(0x22b49f69)).
	SetInvalidating(false).
	BuildAxiomatic()

var StandardColors_NUMBER = Color_SingleBuilder().
	SetId(ColorId_FromUnsignedInt(int32(-0x2f7e8dd4))). // 0xd081722c
	SetBoxId(colorIdPtr(StandardColors_NUMBER_OBJECT_ID)).
	SetInvalidating(false).
	BuildAxiomatic()

var StandardColors_STRING = Color_SingleBuilder().
	SetId(ColorId_FromUnsignedInt(int32(-0x73b2709b))). // 0x8c4d8f65
	SetBoxId(colorIdPtr(StandardColors_STRING_OBJECT_ID)).
	SetInvalidating(false).
	BuildAxiomatic()

var StandardColors_SYMBOL = Color_SingleBuilder().
	SetId(ColorId_FromUnsignedInt(0x759f2066)).
	SetBoxId(colorIdPtr(StandardColors_SYMBOL_OBJECT_ID)).
	SetInvalidating(false).
	BuildAxiomatic()

// StandardColors_TOP_OBJECT is the supertype of all objects but not
// primitives.
//
// Separate from UNKNOWN because some optimizations back off on any
// non-object primitives + unknown but operate on the top object.
var StandardColors_TOP_OBJECT = Color_SingleBuilder().
	SetId(ColorId_FromUnsignedInt(int32(-0x776497c8))). // 0x889b6838
	SetBoxId(nil).
	SetInvalidating(true).
	BuildAxiomatic()

// StandardColors_TOP_FUNCTION is the supertype of all functions.
var StandardColors_TOP_FUNCTION = Color_SingleBuilder().
	SetId(ColorId_FromUnsignedInt(0x283b5838)).
	SetBoxId(nil).
	SetInvalidating(true).
	BuildAxiomatic()

// StandardColors_UNKNOWN is analagous to Closure '*'/'?' and TS unknown/any.
var StandardColors_UNKNOWN = Color_SingleBuilder().
	SetBoxId(nil).
	SetId(ColorId_FromUnsignedInt(0)). // Make UNKNOWN the "default" numerical value.
	SetInvalidating(true).
	BuildAxiomatic()

// Java: StandardColors.GBIGINT
var StandardColors_GBIGINT = Color_SingleBuilder().
	SetId(ColorId_FromUnsignedInt(int32(-0x5551c988))). // 0xaaae3678
	SetBoxId(nil).
	SetInvalidating(false).
	BuildAxiomatic()

// StandardColors_AXIOMATIC_COLORS is the set of Colors whose definitions are
// independent of all JS code.
//
// Because they have no dependencies, they are instantated once, here, and
// made globally available. No other Color object may ever have these IDs.
//
// Some axiomatic Colors do not have an associated runtime type or set of
// runtime instances. These represent abstract concepts, rather than
// anything in a JSVM.
//
// Go-only: language-adapter — assigned in init() rather than by an
// initializer: Color.Builder.buildAxiomatic reads it (and requires it to
// still be null, as during Java's static initialization) while the
// axiomatic colors above are built, which a Go initializer would reject as
// an initialization cycle.
var StandardColors_AXIOMATIC_COLORS *jutil.LinkedHashMap[ColorId, *Color]

func init() {
	m := jutil.NewLinkedHashMap[ColorId, *Color]()
	for _, c := range []*Color{
		StandardColors_BIGINT,
		StandardColors_BOOLEAN,
		StandardColors_NULL_OR_VOID,
		StandardColors_NUMBER,
		StandardColors_STRING,
		StandardColors_SYMBOL,
		StandardColors_TOP_OBJECT,
		StandardColors_UNKNOWN,
		StandardColors_GBIGINT,
		StandardColors_TOP_FUNCTION,
	} {
		immutableMapBuilderPut(m, c.GetId(), c)
	}
	StandardColors_AXIOMATIC_COLORS = m
}

// immutableMapBuilderPut is ImmutableMap.Builder.put followed by
// buildOrThrow's duplicate-key check.
func immutableMapBuilderPut(m *jutil.LinkedHashMap[ColorId, *Color], k ColorId, v *Color) {
	if m.ContainsKey(k) {
		panic(jutil.NewIllegalArgumentExceptionMsg("Multiple entries with same key: " + k.String()))
	}
	m.Put(k, v)
}

// StandardColors_PRIMITIVE_COLORS is the set of Colors that have associated
// runtime values but are not objects.
//
// In JS, instances of primitve Colors are deeply immutable.
var StandardColors_PRIMITIVE_COLORS = func() *jutil.LinkedHashMap[ColorId, *Color] {
	m := jutil.NewLinkedHashMap[ColorId, *Color]()
	for _, c := range []*Color{
		StandardColors_BIGINT,
		StandardColors_BOOLEAN,
		StandardColors_NULL_OR_VOID,
		StandardColors_NUMBER,
		StandardColors_STRING,
		StandardColors_SYMBOL,
		StandardColors_GBIGINT,
	} {
		immutableMapBuilderPut(m, c.GetId(), c)
	}
	return m
}()

// StandardColors_PRIMITIVE_BOX_IDS is the set of ColorIds for object Colors
// that "box" primitive Colors.
//
// Boxing is when the JSVM automatically wraps a primitive value in a
// temporary object so it can participate in a method call or other
// object-only operation.
//
// The eaxct definition of each box Color may be altered by JS code.
var StandardColors_PRIMITIVE_BOX_IDS = jutil.NewLinkedHashSet(
	StandardColors_BIGINT_OBJECT_ID,
	StandardColors_BOOLEAN_OBJECT_ID,
	StandardColors_NUMBER_OBJECT_ID,
	StandardColors_STRING_OBJECT_ID,
	StandardColors_SYMBOL_OBJECT_ID)

// StandardColors_STANDARD_OBJECT_IDS is the set of ColorIds for object
// Colors that need to be referenced in optimizations.
//
// The exact definition of each box Color may be altered by JS code.
var StandardColors_STANDARD_OBJECT_IDS = func() *jutil.LinkedHashSet[ColorId] {
	s := jutil.NewLinkedHashSet[ColorId]()
	s.AddAll(StandardColors_PRIMITIVE_BOX_IDS.ToSlice()...)
	s.Add(StandardColors_ARRAY_ID)
	s.Add(StandardColors_READONLY_ARRAY_ID)
	s.Add(StandardColors_ARGUMENTS_ID)
	s.Add(StandardColors_ASYNC_ITERATOR_ITERABLE_ID)
	s.Add(StandardColors_GENERATOR_ID)
	s.Add(StandardColors_I_TEMPLATE_ARRAY_ID)
	s.Add(StandardColors_ITERATOR_ID)
	s.Add(StandardColors_PROMISE_ID)
	return s
}()

// Java: StandardColors()
func NewStandardColors() *StandardColors {
	panic(jutil.NewAssertionErrorNoMessage())
}
