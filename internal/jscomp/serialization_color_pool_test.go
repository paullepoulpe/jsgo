// Port of test/com/google/javascript/jscomp/serialization/ColorPoolTest.java @ closure-compiler v20261021

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

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func TestColorPool_DeserializesJavaScriptPrimitivesFromEmptyTypePool(t *testing.T) {
	colorPool := ColorPool_FromOnlyShardForTesting(TypePool_GetDefaultInstance(), StringPool_Empty()).GetOnlyShardForTesting()
	jsaJavaEq(t, colorPool.GetColor(jsaInt(jsaPrim(PrimitiveType_NUMBER_TYPE))), colors.StandardColors_NUMBER)
	jsaFalse(t, colorPool.GetColor(jsaInt(jsaPrim(PrimitiveType_NUMBER_TYPE))).IsInvalidating())
	jsaJavaEq(t, colorPool.GetColor(jsaInt(jsaPrim(PrimitiveType_STRING_TYPE))), colors.StandardColors_STRING)
	jsaJavaEq(t, colorPool.GetColor(jsaInt(jsaPrim(PrimitiveType_UNKNOWN_TYPE))), colors.StandardColors_UNKNOWN)
}

func TestColorPool_DeserializesTopObjectTypeFromEmptyTypePool(t *testing.T) {
	colorPool := ColorPool_FromOnlyShardForTesting(TypePool_GetDefaultInstance(), StringPool_Empty()).GetOnlyShardForTesting()
	object := colorPool.GetColor(jsaInt(jsaPrim(PrimitiveType_TOP_OBJECT)))
	jsaTrue(t, object.IsInvalidating())
	jsaTrue(t, len(object.GetInstanceColors()) == 0)
	jsaTrue(t, len(object.GetPrototypes()) == 0)
}

func TestColorPool_DeserializesNativeObjectTableIntoNativeColor(t *testing.T) {
	stringPool := StringPool_Builder_()
	typePool := &TypePool{
		type_: []*TypeProto{
			testObjectTypeProto(&ObjectTypeProto{uuid_: colors.StandardColors_NUMBER_OBJECT_ID.AsByteString(), isInvalidating_: true}),
			testObjectTypeProto(&ObjectTypeProto{uuid_: []byte("Num.pro")}),
		},
		// Number is a subtype of Number.prototype
		disambiguationEdges_: []*SubtypingEdge{{subtype_: jsaPoolPointer(0), supertype_: jsaPoolPointer(1)}},
	}
	colorPool := ColorPool_FromOnlyShardForTesting(typePool, stringPool.Build())
	view := colorPool.GetOnlyShardForTesting()
	numberObject := colorPool.GetColor(colors.StandardColors_NUMBER_OBJECT_ID)
	jsaTrue(t, numberObject.IsInvalidating())
	jsaColorsExactly(t, colorPool.GetRegistry().GetDisambiguationSupertypes(numberObject), view.GetColor(jsaInt(jsaPoolPointer(1))))
	jsaTrue(t, numberObject == view.GetColor(jsaInt(jsaPoolPointer(0))))
	jsaJavaEq(t, numberObject.GetId(), colors.StandardColors_NUMBER_OBJECT_ID)
}

func TestColorPool_TestSynthesizesMissingStandardColors(t *testing.T) {
	colorPool := ColorPool_FromOnlyShardForTesting(TypePool_GetDefaultInstance(), StringPool_Empty())
	for id := range colors.ColorRegistry_REQUIRED_IDS.All() {
		jsaJavaEq(t, colorPool.GetColor(id).GetId(), id)
	}
}

func TestColorPool_DeserializesSimpleObject(t *testing.T) {
	typePool := jsaSingleObjectPool(&ObjectTypeProto{uuid_: []byte("Foo")})
	colorPool := ColorPool_FromOnlyShardForTesting(typePool, StringPool_Empty()).GetOnlyShardForTesting()
	jsaJavaEq(t, colorPool.GetColor(jsaInt(jsaPoolPointer(0))), colors.Color_SingleBuilder().SetId(colors.ColorId_FromAscii("Foo")).Build())
}

func TestColorPool_DeserializesObjectWithPrototypeAndInstanceType(t *testing.T) {
	typePool := &TypePool{
		type_: []*TypeProto{
			testObjectTypeProto(&ObjectTypeProto{uuid_: []byte("Foo.pro")}),
			testObjectTypeProto(&ObjectTypeProto{uuid_: []byte("Foo.ins")}),
			testObjectTypeProto(&ObjectTypeProto{
				uuid_:              []byte("Foo"),
				prototype_:         []int32{jsaPoolPointer(0)},
				instanceType_:      []int32{jsaPoolPointer(1)},
				markedConstructor_: true,
			}),
		},
	}
	colorPool := ColorPool_FromOnlyShardForTesting(typePool, StringPool_Empty()).GetOnlyShardForTesting()
	jsaJavaEq(t, colorPool.GetColor(jsaInt(jsaPoolPointer(2))), jsaCreateObjectColorBuilder().SetId(colors.ColorId_FromAscii("Foo")).SetInstanceColor(jsaCreateObjectColorBuilder().SetId(colors.ColorId_FromAscii("Foo.ins")).Build()).SetPrototype(jsaCreateObjectColorBuilder().SetId(colors.ColorId_FromAscii("Foo.pro")).Build()).SetConstructor(true).Build())
}

func TestColorPool_DeserializesPropertiesFromStringPool(t *testing.T) {
	typePool := jsaSingleObjectPool(&ObjectTypeProto{uuid_: []byte("Foo"), ownProperty_: []int32{1, 2}})
	stringPool := StringPool_Builder_().PutAnd("x").PutAnd("y").Build()
	colorPool := ColorPool_FromOnlyShardForTesting(typePool, stringPool).GetOnlyShardForTesting()
	jsaJavaEq(t, colorPool.GetColor(jsaInt(jsaPoolPointer(0))), jsaCreateObjectColorBuilder().SetId(colors.ColorId_FromAscii("Foo")).SetOwnProperties([]string{"x", "y"}).Build())
}

func TestColorPool_MarksInvalidatingObject(t *testing.T) {
	typePool := jsaSingleObjectPool(&ObjectTypeProto{uuid_: []byte("Foo"), isInvalidating_: true})
	colorPool := ColorPool_FromOnlyShardForTesting(typePool, StringPool_Empty()).GetOnlyShardForTesting()
	jsaTrue(t, colorPool.GetColor(jsaInt(jsaPoolPointer(0))).IsInvalidating())
	jsaJavaEq(t, colorPool.GetColor(jsaInt(jsaPoolPointer(0))).GetId(), colors.ColorId_FromAscii("Foo"))
}

func TestColorPool_MarksClosureAssert(t *testing.T) {
	typePool := jsaSingleObjectPool(&ObjectTypeProto{uuid_: []byte("Foo"), closureAssert_: true})
	colorPool := ColorPool_FromOnlyShardForTesting(typePool, StringPool_Empty()).GetOnlyShardForTesting()
	jsaTrue(t, colorPool.GetColor(jsaInt(jsaPoolPointer(0))).IsClosureAssert())
	jsaJavaEq(t, colorPool.GetColor(jsaInt(jsaPoolPointer(0))).GetId(), colors.ColorId_FromAscii("Foo"))
}

func TestColorPool_RecordsDisambiguationsSupertypes_multipleSupertypesPerSubtype(t *testing.T) {
	typePool := &TypePool{
		type_: []*TypeProto{
			testObjectTypeProto(&ObjectTypeProto{uuid_: []byte("Foo")}),
			testObjectTypeProto(&ObjectTypeProto{uuid_: []byte("Bar")}),
			testObjectTypeProto(&ObjectTypeProto{uuid_: []byte("Baz")}),
		},
		disambiguationEdges_: []*SubtypingEdge{
			// Bar is a subtype of Foo
			{subtype_: jsaPoolPointer(1), supertype_: jsaPoolPointer(0)},
			// Bar is also a subtype of Baz
			{subtype_: jsaPoolPointer(1), supertype_: jsaPoolPointer(2)},
		},
	}
	// When
	colorPoolBuilder := ColorPool_Builder_()
	shard := colorPoolBuilder.AddShard(typePool, StringPool_Empty())
	colorPool := colorPoolBuilder.Build()
	// Then
	foo := shard.GetColor(jsaInt(jsaPoolPointer(0)))
	bar := shard.GetColor(jsaInt(jsaPoolPointer(1)))
	baz := shard.GetColor(jsaInt(jsaPoolPointer(2)))
	jsaColorsExactly(t, colorPool.GetRegistry().GetDisambiguationSupertypes(bar), foo, baz)
}

func TestColorPool_RecordsDisambiguationsSupertypes_cylceInSupertypeGraph(t *testing.T) {
	typePool := &TypePool{
		type_: []*TypeProto{
			testObjectTypeProto(&ObjectTypeProto{uuid_: []byte("Foo")}),
			testObjectTypeProto(&ObjectTypeProto{uuid_: []byte("Bar")}),
		},
		disambiguationEdges_: []*SubtypingEdge{
			// Foo is a subtype of Bar
			{subtype_: jsaPoolPointer(0), supertype_: jsaPoolPointer(1)},
			// Bar is also a subtype of Foo
			{subtype_: jsaPoolPointer(1), supertype_: jsaPoolPointer(0)},
		},
	}
	// When
	colorPoolBuilder := ColorPool_Builder_()
	shard := colorPoolBuilder.AddShard(typePool, StringPool_Empty())
	colorPool := colorPoolBuilder.Build()
	// Then
	foo := shard.GetColor(jsaInt(jsaPoolPointer(0)))
	bar := shard.GetColor(jsaInt(jsaPoolPointer(1)))
	jsaColorsExactly(t, colorPool.GetRegistry().GetDisambiguationSupertypes(foo), bar)
	jsaColorsExactly(t, colorPool.GetRegistry().GetDisambiguationSupertypes(bar), foo)
}

func TestColorPool_ThrowsErrorIfDisambiguationEdgesContainsInvalidOffset(t *testing.T) {
	typePool := &TypePool{
		type_:                []*TypeProto{testObjectTypeProto(&ObjectTypeProto{uuid_: []byte("Foo")})},
		disambiguationEdges_: []*SubtypingEdge{{subtype_: jsaPoolPointer(0), supertype_: jsaPoolPointer(1)}},
	}
	jsaThrows[IMalformedTypedAstException](t, func() { ColorPool_FromOnlyShardForTesting(typePool, StringPool_Empty()) })
}

func TestColorPool_DeserializesMultipleUnionsOfNativeTypes(t *testing.T) {
	typePool := &TypePool{
		type_: []*TypeProto{
			testUnionTypeProto(jsaPrim(PrimitiveType_NUMBER_TYPE), jsaPrim(PrimitiveType_STRING_TYPE)),
			testUnionTypeProto(jsaPrim(PrimitiveType_NUMBER_TYPE), jsaPrim(PrimitiveType_BIGINT_TYPE)),
		},
	}
	colorPool := ColorPool_FromOnlyShardForTesting(typePool, StringPool_Empty())
	view := colorPool.GetOnlyShardForTesting()
	jsaHasAlternates(t, view.GetColor(jsaInt(jsaPoolPointer(0))), colors.StandardColors_STRING, colors.StandardColors_NUMBER)
	jsaHasAlternates(t, view.GetColor(jsaInt(jsaPoolPointer(1))), colors.StandardColors_BIGINT, colors.StandardColors_NUMBER)
}

func TestColorPool_DeserializesUnionReferencingOtherUnion(t *testing.T) {
	typePool := &TypePool{
		type_: []*TypeProto{
			// U0 := (number, string)
			testUnionTypeProto(jsaPrim(PrimitiveType_NUMBER_TYPE), jsaPrim(PrimitiveType_STRING_TYPE)),
			// U1 := (U0, bigint)
			testUnionTypeProto(jsaPoolPointer(0), jsaPrim(PrimitiveType_BIGINT_TYPE)),
		},
	}
	jsaThrows[IMalformedTypedAstException](t, func() { ColorPool_FromOnlyShardForTesting(typePool, StringPool_Empty()) })
}

func TestColorPool_DeserializingUnionsInCycleThrowsErrors(t *testing.T) {
	// Create two union types U0 and U1. U0 := (U1 | NUMBER), U1 := (U0 | number)
	// These cycles are at least possible to construct in the Closure type system.
	// (Note - cycles are much more likely in object types)
	typePool := &TypePool{
		type_: []*TypeProto{
			testUnionTypeProto(jsaPrim(PrimitiveType_NUMBER_TYPE), jsaPoolPointer(1)),
			testUnionTypeProto(jsaPrim(PrimitiveType_NUMBER_TYPE), jsaPoolPointer(0)),
		},
	}
	// Eventually we may need to support this case, but for now throwing an explicit exception is
	// better than infinite recursion.

	jsaThrows[IMalformedTypedAstException](t, func() {
		ColorPool_FromOnlyShardForTesting(typePool, StringPool_Empty()).GetOnlyShardForTesting().GetColor(jsaInt(jsaPoolPointer(0)))
	})
}

func TestColorPool_ThrowsException_onTypeWithoutKindCase(t *testing.T) {
	typePool := &TypePool{type_: []*TypeProto{{}}}
	jsaThrows[IMalformedTypedAstException](t, func() {
		ColorPool_FromOnlyShardForTesting(typePool, StringPool_Empty()).GetOnlyShardForTesting().GetColor(jsaInt(jsaPoolPointer(0)))
	})
}

func TestColorPool_ThrowsException_onSerializedUnionOfSingleElement(t *testing.T) {
	typePool := &TypePool{type_: []*TypeProto{testUnionTypeProto(jsaPrim(PrimitiveType_NUMBER_TYPE))}}
	jsaThrows[IMalformedTypedAstException](t, func() {
		ColorPool_FromOnlyShardForTesting(typePool, StringPool_Empty()).GetOnlyShardForTesting().GetColor(jsaInt(jsaPoolPointer(0)))
	})
}

func TestColorPool_ThrowsException_onOutOfBoundsStringPoolOffset(t *testing.T) {
	typePool := jsaSingleObjectPool(&ObjectTypeProto{uuid_: []byte("Foo"), ownProperty_: []int32{1001}})
	jsaThrows[jutil.IIndexOutOfBoundsException](t, func() { ColorPool_FromOnlyShardForTesting(typePool, StringPool_Empty()) })
}

func TestColorPool_ThrowsException_onIdenticalPools(t *testing.T) {
	typePool := jsaSingleObjectPool(&ObjectTypeProto{uuid_: []byte("Foo")})
	builder := ColorPool_Builder_().AddShardAnd(typePool, StringPool_Empty())
	jsaThrows[jutil.Throwable](t, func() { builder.AddShard(typePool, StringPool_Empty()) })
}

func TestColorPool_ThrowsException_onIdenticalPools_exceptDefaultInstance(t *testing.T) {
	// Given
	builder := ColorPool_Builder_()
	shard := builder.AddShard(TypePool_GetDefaultInstance(), StringPool_Empty())
	// When
	otherShard := builder.AddShard(TypePool_GetDefaultInstance(), StringPool_Empty())
	// Then
	jsaTrue(t, shard == otherShard)
}

func TestColorPool_ThrowsException_onDuplicateIdsInSinglePool(t *testing.T) {
	// Given
	typePool := &TypePool{
		type_: []*TypeProto{
			testObjectTypeProto(&ObjectTypeProto{uuid_: jsaTEST_ID.AsByteString()}),
			testObjectTypeProto(&ObjectTypeProto{uuid_: jsaTEST_ID.AsByteString()}),
		},
	}
	// When & Then
	jsaThrows[IMalformedTypedAstException](t, func() { ColorPool_FromOnlyShardForTesting(typePool, StringPool_Empty()) })
}

func TestColorPool_Uuid_mustBeSet(t *testing.T) {
	// Given
	typePool := jsaSingleObjectPool(&ObjectTypeProto{})
	// When & Then
	jsaThrows[IMalformedTypedAstException](t, func() { ColorPool_Builder_().AddShardAnd(typePool, StringPool_Empty()).Build() })
}

func TestColorPool_Uuid_mustNotBeAxiomatic(t *testing.T) {
	for _, c := range TypePointers_OFFSET_TO_AXIOMATIC_COLOR {
		id := c.GetId()
		// Given
		typePool := jsaSingleObjectPool(&ObjectTypeProto{uuid_: id.AsByteString()})
		// When & Then
		jsaThrows[IMalformedTypedAstException](t, func() { ColorPool_Builder_().AddShardAnd(typePool, StringPool_Empty()).Build() })
	}
}

func TestColorPool_Reconcile_defaultValues(t *testing.T) {
	// Given
	typePool := jsaSingleObjectPool(&ObjectTypeProto{uuid_: jsaTEST_ID.AsByteString()})
	// When
	colorPool := ColorPool_Builder_().AddShardAnd(typePool, StringPool_Empty()).Build()
	// Then
	defaultColor := colorPool.GetColor(jsaTEST_ID)
	jsaTrue(t, len(defaultColor.GetInstanceColors()) == 0)
	jsaTrue(t, len(defaultColor.GetPrototypes()) == 0)
	jsaTrue(t, len(defaultColor.GetOwnProperties()) == 0)
	jsaFalse(t, defaultColor.IsClosureAssert())
	jsaFalse(t, defaultColor.IsConstructor())
	jsaFalse(t, defaultColor.IsInvalidating())
	jsaFalse(t, defaultColor.GetPropertiesKeepOriginalName())
}

func TestColorPool_Reconcile_setDisambiguationSupertypes(t *testing.T) {
	// Given
	typePool0 := &TypePool{
		type_:                []*TypeProto{testObjectTypeProto(&ObjectTypeProto{uuid_: jsaTEST_ID.AsByteString()})},
		disambiguationEdges_: []*SubtypingEdge{{subtype_: jsaPoolPointer(0), supertype_: jsaPrim(PrimitiveType_NUMBER_TYPE)}},
	}
	typePool1 := &TypePool{
		type_:                []*TypeProto{testObjectTypeProto(&ObjectTypeProto{uuid_: jsaTEST_ID.AsByteString()})},
		disambiguationEdges_: []*SubtypingEdge{{subtype_: jsaPoolPointer(0), supertype_: jsaPrim(PrimitiveType_STRING_TYPE)}},
	}
	// When
	colorPool := ColorPool_Builder_().AddShardAnd(typePool0, StringPool_Empty()).AddShardAnd(typePool1, StringPool_Empty()).Build()
	// Then
	jsaColorsExactly(t, colorPool.GetRegistry().GetDisambiguationSupertypes(colorPool.GetColor(jsaTEST_ID)), colors.StandardColors_NUMBER, colors.StandardColors_STRING)
}

func TestColorPool_Reconcile_setInstanceColors(t *testing.T) {
	// Given
	typePool0 := jsaSingleObjectPool(&ObjectTypeProto{uuid_: jsaTEST_ID.AsByteString(), instanceType_: []int32{jsaPrim(PrimitiveType_NUMBER_TYPE)}})
	typePool1 := jsaSingleObjectPool(&ObjectTypeProto{uuid_: jsaTEST_ID.AsByteString(), instanceType_: []int32{jsaPrim(PrimitiveType_STRING_TYPE)}})
	// When
	colorPool := ColorPool_Builder_().AddShardAnd(typePool0, StringPool_Empty()).AddShardAnd(typePool1, StringPool_Empty()).Build()
	// Then
	jsaColorSliceExactly(t, colorPool.GetColor(jsaTEST_ID).GetInstanceColors(), colors.StandardColors_NUMBER, colors.StandardColors_STRING)
}

func TestColorPool_Reconcile_addPrototypes(t *testing.T) {
	// Given
	typePool0 := jsaSingleObjectPool(&ObjectTypeProto{uuid_: jsaTEST_ID.AsByteString(), prototype_: []int32{jsaPrim(PrimitiveType_NUMBER_TYPE)}})
	typePool1 := jsaSingleObjectPool(&ObjectTypeProto{uuid_: jsaTEST_ID.AsByteString(), prototype_: []int32{jsaPrim(PrimitiveType_STRING_TYPE)}})
	// When
	colorPool := ColorPool_Builder_().AddShardAnd(typePool0, StringPool_Empty()).AddShardAnd(typePool1, StringPool_Empty()).Build()
	// Then
	jsaColorSliceExactly(t, colorPool.GetColor(jsaTEST_ID).GetPrototypes(), colors.StandardColors_NUMBER, colors.StandardColors_STRING)
}

func TestColorPool_Reconcile_setOwnProperties(t *testing.T) {
	// Given
	stringPool0 := StringPool_Builder_().PutAnd("hello").Build()
	typePool0 := jsaSingleObjectPool(&ObjectTypeProto{uuid_: jsaTEST_ID.AsByteString(), ownProperty_: []int32{1}})
	stringPool1 := StringPool_Builder_().PutAnd("world").Build()
	typePool1 := jsaSingleObjectPool(&ObjectTypeProto{uuid_: jsaTEST_ID.AsByteString(), ownProperty_: []int32{1}})
	// When
	colorPool := ColorPool_Builder_().AddShardAnd(typePool0, stringPool0).AddShardAnd(typePool1, stringPool1).Build()
	// Then
	jsaStringsExactly(t, colorPool.GetColor(jsaTEST_ID).GetOwnProperties(), "hello", "world")
}

func TestColorPool_Reconcile_setClosureAssert_true(t *testing.T) {
	// Given
	typePool0 := jsaSingleObjectPool(&ObjectTypeProto{uuid_: jsaTEST_ID.AsByteString(), closureAssert_: true})
	typePool1 := jsaSingleObjectPool(&ObjectTypeProto{uuid_: jsaTEST_ID.AsByteString(), closureAssert_: true})
	// When
	colorPool := ColorPool_Builder_().AddShardAnd(typePool0, StringPool_Empty()).AddShardAnd(typePool1, StringPool_Empty()).Build()
	// Then
	jsaTrue(t, colorPool.GetColor(jsaTEST_ID).IsClosureAssert())
}

func TestColorPool_Reconcile_setClosureAssert_false(t *testing.T) {
	// Given
	typePool0 := jsaSingleObjectPool(&ObjectTypeProto{uuid_: jsaTEST_ID.AsByteString(), closureAssert_: false})
	typePool1 := jsaSingleObjectPool(&ObjectTypeProto{uuid_: jsaTEST_ID.AsByteString(), closureAssert_: false})
	// When
	colorPool := ColorPool_Builder_().AddShardAnd(typePool0, StringPool_Empty()).AddShardAnd(typePool1, StringPool_Empty()).Build()
	// Then
	jsaFalse(t, colorPool.GetColor(jsaTEST_ID).IsClosureAssert())
}

func TestColorPool_Reconcile_setClosureAssert_mixed(t *testing.T) {
	// Given
	typePool0 := jsaSingleObjectPool(&ObjectTypeProto{uuid_: jsaTEST_ID.AsByteString(), closureAssert_: true})
	typePool1 := jsaSingleObjectPool(&ObjectTypeProto{uuid_: jsaTEST_ID.AsByteString(), closureAssert_: false})
	// When & Then
	jsaThrows[IMalformedTypedAstException](t, func() {
		ColorPool_Builder_().AddShardAnd(typePool0, StringPool_Empty()).AddShardAnd(typePool1, StringPool_Empty()).Build()
	})
}

func TestColorPool_Reconcile_setConstructor(t *testing.T) {
	// Given
	typePool0 := jsaSingleObjectPool(&ObjectTypeProto{uuid_: jsaTEST_ID.AsByteString(), markedConstructor_: false})
	typePool1 := jsaSingleObjectPool(&ObjectTypeProto{uuid_: jsaTEST_ID.AsByteString(), markedConstructor_: true})
	// When
	colorPool := ColorPool_Builder_().AddShardAnd(typePool0, StringPool_Empty()).AddShardAnd(typePool1, StringPool_Empty()).Build()
	// Then
	jsaTrue(t, colorPool.GetColor(jsaTEST_ID).IsConstructor())
}

func TestColorPool_Reconcile_setInvalidating(t *testing.T) {
	// Given
	typePool0 := jsaSingleObjectPool(&ObjectTypeProto{uuid_: jsaTEST_ID.AsByteString(), isInvalidating_: false})
	typePool1 := jsaSingleObjectPool(&ObjectTypeProto{uuid_: jsaTEST_ID.AsByteString(), isInvalidating_: true})
	// When
	colorPool := ColorPool_Builder_().AddShardAnd(typePool0, StringPool_Empty()).AddShardAnd(typePool1, StringPool_Empty()).Build()
	// Then
	jsaTrue(t, colorPool.GetColor(jsaTEST_ID).IsInvalidating())
}

func TestColorPool_Reconcile_setPropertiesKeepOriginalName(t *testing.T) {
	// Given
	typePool0 := jsaSingleObjectPool(&ObjectTypeProto{uuid_: jsaTEST_ID.AsByteString(), propertiesKeepOriginalName_: false})
	typePool1 := jsaSingleObjectPool(&ObjectTypeProto{uuid_: jsaTEST_ID.AsByteString(), propertiesKeepOriginalName_: true})
	// When
	colorPool := ColorPool_Builder_().AddShardAnd(typePool0, StringPool_Empty()).AddShardAnd(typePool1, StringPool_Empty()).Build()
	// Then
	jsaTrue(t, colorPool.GetColor(jsaTEST_ID).GetPropertiesKeepOriginalName())
}

func TestColorPool_Reconcile_debugInfo_mismatches(t *testing.T) {
	// Given
	typePool0 := &TypePool{
		debugInfo_: &TypePool_DebugInfo{
			mismatch_: []*TypePool_DebugInfo_Mismatch{
				{sourceRef_: "location_0", involvedColor_: []int32{jsaPrim(PrimitiveType_SYMBOL_TYPE), jsaPrim(PrimitiveType_STRING_TYPE)}},
				{sourceRef_: "location_1", involvedColor_: []int32{jsaPrim(PrimitiveType_NUMBER_TYPE)}},
			},
		},
	}
	typePool1 := &TypePool{
		debugInfo_: &TypePool_DebugInfo{
			mismatch_: []*TypePool_DebugInfo_Mismatch{
				{sourceRef_: "location_0", involvedColor_: []int32{jsaPrim(PrimitiveType_SYMBOL_TYPE), jsaPrim(PrimitiveType_BOOLEAN_TYPE)}},
				{sourceRef_: "location_2", involvedColor_: []int32{jsaPrim(PrimitiveType_BIGINT_TYPE)}},
			},
		},
	}
	// When
	colorPool := ColorPool_Builder_().AddShardAnd(typePool0, StringPool_Empty()).AddShardAnd(typePool1, StringPool_Empty()).Build()
	// Then
	expected := colors.NewColorMismatchLocations()
	expected.Put(colors.StandardColors_BOOLEAN.GetId(), "location_0")
	expected.Put(colors.StandardColors_STRING.GetId(), "location_0")
	expected.Put(colors.StandardColors_SYMBOL.GetId(), "location_0")
	expected.Put(colors.StandardColors_NUMBER.GetId(), "location_1")
	expected.Put(colors.StandardColors_BIGINT.GetId(), "location_2")
	jsaColorMismatchEq(t, colorPool.GetRegistry().GetMismatchLocationsForDebugging(), expected)
}

func TestColorPool_Reconcile_shards_preserveOriginalOffsets(t *testing.T) {
	// Given
	otherId := colors.ColorId_FromAscii("other")
	typePool0 := &TypePool{
		type_: []*TypeProto{
			testObjectTypeProto(&ObjectTypeProto{uuid_: otherId.AsByteString()}),
			testObjectTypeProto(&ObjectTypeProto{uuid_: jsaTEST_ID.AsByteString()}),
		},
	}
	typePool1 := &TypePool{
		type_: []*TypeProto{
			testObjectTypeProto(&ObjectTypeProto{uuid_: jsaTEST_ID.AsByteString()}),
			testObjectTypeProto(&ObjectTypeProto{uuid_: otherId.AsByteString()}),
		},
	}
	// When
	colorPoolBuilder := ColorPool_Builder_()
	shard0 := colorPoolBuilder.AddShard(typePool0, StringPool_Empty())
	shard1 := colorPoolBuilder.AddShard(typePool1, StringPool_Empty())
	colorPoolBuilder.Build()
	// Then
	jsaJavaEq(t, shard0.GetColor(jsaInt(jsaPoolPointer(1))).GetId(), jsaTEST_ID)
	jsaTrue(t, shard0.GetColor(jsaInt(jsaPoolPointer(1))) == shard1.GetColor(jsaInt(jsaPoolPointer(0))))
}

var jsaTEST_ID = colors.ColorId_FromUnsigned(100)

func jsaCreateObjectColorBuilder() *colors.Color_Builder {
	return colors.Color_SingleBuilder().SetId(colors.ColorId_FromAscii(""))
}

func testObjectTypeProto(obj *ObjectTypeProto) *TypeProto {
	return &TypeProto{object_: obj, kindCase_: TypeProto_OBJECT_FIELD_NUMBER}
}

func testUnionTypeProto(members ...int32) *TypeProto {
	return &TypeProto{union_: &UnionTypeProto{unionMember_: members}, kindCase_: TypeProto_UNION_FIELD_NUMBER}
}

func jsaSingleObjectPool(obj *ObjectTypeProto) *TypePool {
	return &TypePool{type_: []*TypeProto{testObjectTypeProto(obj)}}
}

func jsaPoolPointer(offset int) int32 {
	return int32(TypePointers_UntrimOffset(offset))
}

func jsaPrim(p PrimitiveType) int32 { return p.GetNumber() }

func jsaInt(x int32) int { return int(x) }

// jsaColorsExactly is assertThat(colors).containsExactly(expected...).
func jsaColorsExactly(t *testing.T, got *jutil.ObjectLinkedHashSet[*colors.Color], expected ...*colors.Color) {
	t.Helper()
	ok := got.Size() == len(expected)
	for _, e := range expected {
		ok = ok && got.Contains(e)
	}
	if !ok {
		t.Errorf("got %v, want exactly %v", got, expected)
	}
}

func jsaColorSliceExactly(t *testing.T, got []*colors.Color, expected ...*colors.Color) {
	t.Helper()
	ok := len(got) == len(expected)
	for _, e := range expected {
		found := false
		for _, g := range got {
			if g == e || g.Equals(e) {
				found = true
				break
			}
		}
		ok = ok && found
	}
	if !ok {
		t.Errorf("got %v, want exactly %v", got, expected)
	}
}

// jsaHasAlternates is ColorSubject.hasAlternates(Color...).
func jsaHasAlternates(t *testing.T, c *colors.Color, alternates ...*colors.Color) {
	t.Helper()
	jsaTrue(t, c.IsUnion(), c, " is a union")
	jsaColorSliceExactly(t, c.GetUnionElements(), alternates...)
}

func jsaStringsExactly(t *testing.T, got []string, expected ...string) {
	t.Helper()
	ok := len(got) == len(expected)
	for _, e := range expected {
		found := false
		for _, g := range got {
			if g == e {
				found = true
				break
			}
		}
		ok = ok && found
	}
	if !ok {
		t.Errorf("got %v, want exactly %v", got, expected)
	}
}

func jsaColorMismatchEq(t *testing.T, got, want *colors.ColorMismatchLocations) {
	t.Helper()
	ok := got.Size() == want.Size()
	for k, v := range want.Entries() {
		ok = ok && got.Get(k).Contains(v)
	}
	if !ok {
		t.Errorf("got %v, want %v", got, want)
	}
}
