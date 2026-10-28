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
	typePool := TypePool_NewBuilder().AddType(TypeProto_NewBuilder().SetObjectFromBuilder(ObjectTypeProto_NewBuilder().SetUuid(colors.StandardColors_NUMBER_OBJECT_ID.AsByteString()).SetIsInvalidating(true)).Build()).AddType(TypeProto_NewBuilder().SetObjectFromBuilder(ObjectTypeProto_NewBuilder().SetUuid([]byte("Num.pro"))).Build()).AddDisambiguationEdgesFromBuilder(
		// Number is a subtype of Number.prototype
		SubtypingEdge_NewBuilder().SetSubtype(jsaPoolPointer(0)).SetSupertype(jsaPoolPointer(1))).Build()
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
	typePool := TypePool_NewBuilder().AddTypeFromBuilder(TypeProto_NewBuilder().SetObjectFromBuilder(ObjectTypeProto_NewBuilder().SetUuid([]byte("Foo")))).Build()
	colorPool := ColorPool_FromOnlyShardForTesting(typePool, StringPool_Empty()).GetOnlyShardForTesting()
	jsaJavaEq(t, colorPool.GetColor(jsaInt(jsaPoolPointer(0))), colors.Color_SingleBuilder().SetId(colors.ColorId_FromAscii("Foo")).Build())
}

func TestColorPool_DeserializesObjectWithPrototypeAndInstanceType(t *testing.T) {
	typePool := TypePool_NewBuilder().AddTypeFromBuilder(TypeProto_NewBuilder().SetObjectFromBuilder(ObjectTypeProto_NewBuilder().SetUuid([]byte("Foo.pro")))).AddTypeFromBuilder(TypeProto_NewBuilder().SetObjectFromBuilder(ObjectTypeProto_NewBuilder().SetUuid([]byte("Foo.ins")))).AddTypeFromBuilder(TypeProto_NewBuilder().SetObjectFromBuilder(ObjectTypeProto_NewBuilder().SetUuid([]byte("Foo")).AddPrototype(jsaPoolPointer(0)).AddInstanceType(jsaPoolPointer(1)).SetMarkedConstructor(true))).Build()
	colorPool := ColorPool_FromOnlyShardForTesting(typePool, StringPool_Empty()).GetOnlyShardForTesting()
	jsaJavaEq(t, colorPool.GetColor(jsaInt(jsaPoolPointer(2))), jsaCreateObjectColorBuilder().SetId(colors.ColorId_FromAscii("Foo")).SetInstanceColor(jsaCreateObjectColorBuilder().SetId(colors.ColorId_FromAscii("Foo.ins")).Build()).SetPrototype(jsaCreateObjectColorBuilder().SetId(colors.ColorId_FromAscii("Foo.pro")).Build()).SetConstructor(true).Build())
}

func TestColorPool_DeserializesPropertiesFromStringPool(t *testing.T) {
	typePool := TypePool_NewBuilder().AddTypeFromBuilder(TypeProto_NewBuilder().SetObjectFromBuilder(ObjectTypeProto_NewBuilder().SetUuid([]byte("Foo")).AddOwnProperty(1).AddOwnProperty(2))).Build()
	stringPool := StringPool_Builder_().PutAnd("x").PutAnd("y").Build()
	colorPool := ColorPool_FromOnlyShardForTesting(typePool, stringPool).GetOnlyShardForTesting()
	jsaJavaEq(t, colorPool.GetColor(jsaInt(jsaPoolPointer(0))), jsaCreateObjectColorBuilder().SetId(colors.ColorId_FromAscii("Foo")).SetOwnProperties([]string{"x", "y"}).Build())
}

func TestColorPool_MarksInvalidatingObject(t *testing.T) {
	typePool := TypePool_NewBuilder().AddTypeFromBuilder(TypeProto_NewBuilder().SetObjectFromBuilder(ObjectTypeProto_NewBuilder().SetUuid([]byte("Foo")).SetIsInvalidating(true))).Build()
	colorPool := ColorPool_FromOnlyShardForTesting(typePool, StringPool_Empty()).GetOnlyShardForTesting()
	jsaTrue(t, colorPool.GetColor(jsaInt(jsaPoolPointer(0))).IsInvalidating())
	jsaJavaEq(t, colorPool.GetColor(jsaInt(jsaPoolPointer(0))).GetId(), colors.ColorId_FromAscii("Foo"))
}

func TestColorPool_MarksClosureAssert(t *testing.T) {
	typePool := TypePool_NewBuilder().AddTypeFromBuilder(TypeProto_NewBuilder().SetObjectFromBuilder(ObjectTypeProto_NewBuilder().SetUuid([]byte("Foo")).SetClosureAssert(true))).Build()
	colorPool := ColorPool_FromOnlyShardForTesting(typePool, StringPool_Empty()).GetOnlyShardForTesting()
	jsaTrue(t, colorPool.GetColor(jsaInt(jsaPoolPointer(0))).IsClosureAssert())
	jsaJavaEq(t, colorPool.GetColor(jsaInt(jsaPoolPointer(0))).GetId(), colors.ColorId_FromAscii("Foo"))
}

func TestColorPool_RecordsDisambiguationsSupertypes_multipleSupertypesPerSubtype(t *testing.T) {
	typePool := TypePool_NewBuilder().AddTypeFromBuilder(TypeProto_NewBuilder().SetObjectFromBuilder(ObjectTypeProto_NewBuilder().SetUuid([]byte("Foo")))).AddTypeFromBuilder(TypeProto_NewBuilder().SetObjectFromBuilder(ObjectTypeProto_NewBuilder().SetUuid([]byte("Bar")))).AddTypeFromBuilder(TypeProto_NewBuilder().SetObjectFromBuilder(ObjectTypeProto_NewBuilder().SetUuid([]byte("Baz")))).
		// Bar is a subtype of Foo
		AddDisambiguationEdgesFromBuilder(SubtypingEdge_NewBuilder().SetSubtype(jsaPoolPointer(1)).SetSupertype(jsaPoolPointer(0))).
		// Bar is also a subtype of Baz
		AddDisambiguationEdgesFromBuilder(SubtypingEdge_NewBuilder().SetSubtype(jsaPoolPointer(1)).SetSupertype(jsaPoolPointer(2))).Build()
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
	typePool := TypePool_NewBuilder().AddTypeFromBuilder(TypeProto_NewBuilder().SetObjectFromBuilder(ObjectTypeProto_NewBuilder().SetUuid([]byte("Foo")))).AddTypeFromBuilder(TypeProto_NewBuilder().SetObjectFromBuilder(ObjectTypeProto_NewBuilder().SetUuid([]byte("Bar")))).
		// Foo is a subtype of Bar
		AddDisambiguationEdgesFromBuilder(SubtypingEdge_NewBuilder().SetSubtype(jsaPoolPointer(0)).SetSupertype(jsaPoolPointer(1))).
		// Bar is also a subtype of Foo
		AddDisambiguationEdgesFromBuilder(SubtypingEdge_NewBuilder().SetSubtype(jsaPoolPointer(1)).SetSupertype(jsaPoolPointer(0))).Build()
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
	typePool := TypePool_NewBuilder().AddTypeFromBuilder(TypeProto_NewBuilder().SetObjectFromBuilder(ObjectTypeProto_NewBuilder().SetUuid([]byte("Foo")))).AddDisambiguationEdgesFromBuilder(SubtypingEdge_NewBuilder().SetSubtype(jsaPoolPointer(0)).SetSupertype(jsaPoolPointer(1))).Build()
	jsaThrows[IMalformedTypedAstException](t, func() { ColorPool_FromOnlyShardForTesting(typePool, StringPool_Empty()) })
}

func TestColorPool_DeserializesMultipleUnionsOfNativeTypes(t *testing.T) {
	typePool := TypePool_NewBuilder().AddTypeFromBuilder(TypeProto_NewBuilder().SetUnionFromBuilder(UnionTypeProto_NewBuilder().AddUnionMember(jsaPrim(PrimitiveType_NUMBER_TYPE)).AddUnionMember(jsaPrim(PrimitiveType_STRING_TYPE)))).AddTypeFromBuilder(TypeProto_NewBuilder().SetUnionFromBuilder(UnionTypeProto_NewBuilder().AddUnionMember(jsaPrim(PrimitiveType_NUMBER_TYPE)).AddUnionMember(jsaPrim(PrimitiveType_BIGINT_TYPE)))).Build()
	colorPool := ColorPool_FromOnlyShardForTesting(typePool, StringPool_Empty())
	view := colorPool.GetOnlyShardForTesting()
	jsaHasAlternates(t, view.GetColor(jsaInt(jsaPoolPointer(0))), colors.StandardColors_STRING, colors.StandardColors_NUMBER)
	jsaHasAlternates(t, view.GetColor(jsaInt(jsaPoolPointer(1))), colors.StandardColors_BIGINT, colors.StandardColors_NUMBER)
}

func TestColorPool_DeserializesUnionReferencingOtherUnion(t *testing.T) {
	typePool := TypePool_NewBuilder().
		// U0 := (number, string)
		AddTypeFromBuilder(TypeProto_NewBuilder().SetUnionFromBuilder(UnionTypeProto_NewBuilder().AddUnionMember(jsaPrim(PrimitiveType_NUMBER_TYPE)).AddUnionMember(jsaPrim(PrimitiveType_STRING_TYPE)))).
		// U1 := (U0, bigint)
		AddTypeFromBuilder(TypeProto_NewBuilder().SetUnionFromBuilder(UnionTypeProto_NewBuilder().AddUnionMember(jsaPoolPointer(0)).AddUnionMember(jsaPrim(PrimitiveType_BIGINT_TYPE)))).Build()
	jsaThrows[IMalformedTypedAstException](t, func() { ColorPool_FromOnlyShardForTesting(typePool, StringPool_Empty()) })
}

func TestColorPool_DeserializingUnionsInCycleThrowsErrors(t *testing.T) {
	// Create two union types U0 and U1. U0 := (U1 | NUMBER), U1 := (U0 | number)
	// These cycles are at least possible to construct in the Closure type system.
	// (Note - cycles are much more likely in object types)
	typePool := TypePool_NewBuilder().AddTypeFromBuilder(TypeProto_NewBuilder().SetUnionFromBuilder(UnionTypeProto_NewBuilder().AddUnionMember(jsaPrim(PrimitiveType_NUMBER_TYPE)).AddUnionMember(jsaPoolPointer(1)))).AddTypeFromBuilder(TypeProto_NewBuilder().SetUnionFromBuilder(UnionTypeProto_NewBuilder().AddUnionMember(jsaPrim(PrimitiveType_NUMBER_TYPE)).AddUnionMember(jsaPoolPointer(0)))).Build()
	// Eventually we may need to support this case, but for now throwing an explicit exception is
	// better than infinite recursion.

	jsaThrows[IMalformedTypedAstException](t, func() {
		ColorPool_FromOnlyShardForTesting(typePool, StringPool_Empty()).GetOnlyShardForTesting().GetColor(jsaInt(jsaPoolPointer(0)))
	})
}

func TestColorPool_ThrowsException_onTypeWithoutKindCase(t *testing.T) {
	typePool := TypePool_NewBuilder().AddType(TypeProto_GetDefaultInstance()).Build()
	jsaThrows[IMalformedTypedAstException](t, func() {
		ColorPool_FromOnlyShardForTesting(typePool, StringPool_Empty()).GetOnlyShardForTesting().GetColor(jsaInt(jsaPoolPointer(0)))
	})
}

func TestColorPool_ThrowsException_onSerializedUnionOfSingleElement(t *testing.T) {
	typePool := TypePool_NewBuilder().AddTypeFromBuilder(TypeProto_NewBuilder().SetUnionFromBuilder(UnionTypeProto_NewBuilder().AddUnionMember(jsaPrim(PrimitiveType_NUMBER_TYPE)))).Build()
	jsaThrows[IMalformedTypedAstException](t, func() {
		ColorPool_FromOnlyShardForTesting(typePool, StringPool_Empty()).GetOnlyShardForTesting().GetColor(jsaInt(jsaPoolPointer(0)))
	})
}

func TestColorPool_ThrowsException_onOutOfBoundsStringPoolOffset(t *testing.T) {
	typePool := TypePool_NewBuilder().AddTypeFromBuilder(TypeProto_NewBuilder().SetObjectFromBuilder(ObjectTypeProto_NewBuilder().SetUuid([]byte("Foo")).AddOwnProperty(1001))).Build()
	jsaThrows[jutil.IIndexOutOfBoundsException](t, func() { ColorPool_FromOnlyShardForTesting(typePool, StringPool_Empty()) })
}

func TestColorPool_ThrowsException_onIdenticalPools(t *testing.T) {
	typePool := TypePool_NewBuilder().AddTypeFromBuilder(TypeProto_NewBuilder().SetObjectFromBuilder(ObjectTypeProto_NewBuilder().SetUuid([]byte("Foo")))).Build()
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
	typePool := TypePool_NewBuilder().AddTypeFromBuilder(TypeProto_NewBuilder().SetObject(ObjectTypeProto_NewBuilder().SetUuid(jsaTEST_ID.AsByteString()).Build())).AddTypeFromBuilder(TypeProto_NewBuilder().SetObject(ObjectTypeProto_NewBuilder().SetUuid(jsaTEST_ID.AsByteString()).Build())).Build()
	// When & Then
	jsaThrows[IMalformedTypedAstException](t, func() { ColorPool_FromOnlyShardForTesting(typePool, StringPool_Empty()) })
}

func TestColorPool_Uuid_mustBeSet(t *testing.T) {
	// Given
	typePool := jsaSingleObjectPool(ObjectTypeProto_NewBuilder())
	// When & Then
	jsaThrows[IMalformedTypedAstException](t, func() { ColorPool_Builder_().AddShardAnd(typePool, StringPool_Empty()).Build() })
}

func TestColorPool_Uuid_mustNotBeAxiomatic(t *testing.T) {
	for id := range colors.StandardColors_AXIOMATIC_COLORS.Keys() {
		// Given
		typePool := jsaSingleObjectPool(ObjectTypeProto_NewBuilder().SetUuid(id.AsByteString()))
		// When & Then
		jsaThrows[IMalformedTypedAstException](t, func() { ColorPool_Builder_().AddShardAnd(typePool, StringPool_Empty()).Build() })
	}
}

func TestColorPool_Reconcile_defaultValues(t *testing.T) {
	// Given
	typePool := jsaSingleObjectPool(ObjectTypeProto_NewBuilder().SetUuid(jsaTEST_ID.AsByteString()))
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
	typePool0 := TypePool_NewBuilder().AddTypeFromBuilder(TypeProto_NewBuilder().SetObject(ObjectTypeProto_NewBuilder().SetUuid(jsaTEST_ID.AsByteString()).Build())).AddDisambiguationEdges(SubtypingEdge_NewBuilder().SetSubtype(jsaPoolPointer(0)).SetSupertype(jsaPrim(PrimitiveType_NUMBER_TYPE)).Build()).Build()
	typePool1 := TypePool_NewBuilder().AddTypeFromBuilder(TypeProto_NewBuilder().SetObject(ObjectTypeProto_NewBuilder().SetUuid(jsaTEST_ID.AsByteString()).Build())).AddDisambiguationEdges(SubtypingEdge_NewBuilder().SetSubtype(jsaPoolPointer(0)).SetSupertype(jsaPrim(PrimitiveType_STRING_TYPE)).Build()).Build()
	// When
	colorPool := ColorPool_Builder_().AddShardAnd(typePool0, StringPool_Empty()).AddShardAnd(typePool1, StringPool_Empty()).Build()
	// Then
	jsaColorsExactly(t, colorPool.GetRegistry().GetDisambiguationSupertypes(colorPool.GetColor(jsaTEST_ID)), colors.StandardColors_NUMBER, colors.StandardColors_STRING)
}

func TestColorPool_Reconcile_setInstanceColors(t *testing.T) {
	// Given
	typePool0 := jsaSingleObjectPool(ObjectTypeProto_NewBuilder().SetUuid(jsaTEST_ID.AsByteString()).AddInstanceType(jsaPrim(PrimitiveType_NUMBER_TYPE)))
	typePool1 := jsaSingleObjectPool(ObjectTypeProto_NewBuilder().SetUuid(jsaTEST_ID.AsByteString()).AddInstanceType(jsaPrim(PrimitiveType_STRING_TYPE)))
	// When
	colorPool := ColorPool_Builder_().AddShardAnd(typePool0, StringPool_Empty()).AddShardAnd(typePool1, StringPool_Empty()).Build()
	// Then
	jsaColorSliceExactly(t, colorPool.GetColor(jsaTEST_ID).GetInstanceColors(), colors.StandardColors_NUMBER, colors.StandardColors_STRING)
}

func TestColorPool_Reconcile_addPrototypes(t *testing.T) {
	// Given
	typePool0 := jsaSingleObjectPool(ObjectTypeProto_NewBuilder().SetUuid(jsaTEST_ID.AsByteString()).AddPrototype(jsaPrim(PrimitiveType_NUMBER_TYPE)))
	typePool1 := jsaSingleObjectPool(ObjectTypeProto_NewBuilder().SetUuid(jsaTEST_ID.AsByteString()).AddPrototype(jsaPrim(PrimitiveType_STRING_TYPE)))
	// When
	colorPool := ColorPool_Builder_().AddShardAnd(typePool0, StringPool_Empty()).AddShardAnd(typePool1, StringPool_Empty()).Build()
	// Then
	jsaColorSliceExactly(t, colorPool.GetColor(jsaTEST_ID).GetPrototypes(), colors.StandardColors_NUMBER, colors.StandardColors_STRING)
}

func TestColorPool_Reconcile_setOwnProperties(t *testing.T) {
	// Given
	stringPool0 := StringPool_Builder_().PutAnd("hello").Build()
	typePool0 := jsaSingleObjectPool(ObjectTypeProto_NewBuilder().SetUuid(jsaTEST_ID.AsByteString()).AddOwnProperty(1))
	stringPool1 := StringPool_Builder_().PutAnd("world").Build()
	typePool1 := jsaSingleObjectPool(ObjectTypeProto_NewBuilder().SetUuid(jsaTEST_ID.AsByteString()).AddOwnProperty(1))
	// When
	colorPool := ColorPool_Builder_().AddShardAnd(typePool0, stringPool0).AddShardAnd(typePool1, stringPool1).Build()
	// Then
	jsaStringsExactly(t, colorPool.GetColor(jsaTEST_ID).GetOwnProperties(), "hello", "world")
}

func TestColorPool_Reconcile_setClosureAssert_true(t *testing.T) {
	// Given
	typePool0 := jsaSingleObjectPool(ObjectTypeProto_NewBuilder().SetUuid(jsaTEST_ID.AsByteString()).SetClosureAssert(true))
	typePool1 := jsaSingleObjectPool(ObjectTypeProto_NewBuilder().SetUuid(jsaTEST_ID.AsByteString()).SetClosureAssert(true))
	// When
	colorPool := ColorPool_Builder_().AddShardAnd(typePool0, StringPool_Empty()).AddShardAnd(typePool1, StringPool_Empty()).Build()
	// Then
	jsaTrue(t, colorPool.GetColor(jsaTEST_ID).IsClosureAssert())
}

func TestColorPool_Reconcile_setClosureAssert_false(t *testing.T) {
	// Given
	typePool0 := jsaSingleObjectPool(ObjectTypeProto_NewBuilder().SetUuid(jsaTEST_ID.AsByteString()).SetClosureAssert(false))
	typePool1 := jsaSingleObjectPool(ObjectTypeProto_NewBuilder().SetUuid(jsaTEST_ID.AsByteString()).SetClosureAssert(false))
	// When
	colorPool := ColorPool_Builder_().AddShardAnd(typePool0, StringPool_Empty()).AddShardAnd(typePool1, StringPool_Empty()).Build()
	// Then
	jsaFalse(t, colorPool.GetColor(jsaTEST_ID).IsClosureAssert())
}

func TestColorPool_Reconcile_setClosureAssert_mixed(t *testing.T) {
	// Given
	typePool0 := jsaSingleObjectPool(ObjectTypeProto_NewBuilder().SetUuid(jsaTEST_ID.AsByteString()).SetClosureAssert(true))
	typePool1 := jsaSingleObjectPool(ObjectTypeProto_NewBuilder().SetUuid(jsaTEST_ID.AsByteString()).SetClosureAssert(false))
	// When & Then
	jsaThrows[IMalformedTypedAstException](t, func() {
		ColorPool_Builder_().AddShardAnd(typePool0, StringPool_Empty()).AddShardAnd(typePool1, StringPool_Empty()).Build()
	})
}

func TestColorPool_Reconcile_setConstructor(t *testing.T) {
	// Given
	typePool0 := jsaSingleObjectPool(ObjectTypeProto_NewBuilder().SetUuid(jsaTEST_ID.AsByteString()).SetMarkedConstructor(false))
	typePool1 := jsaSingleObjectPool(ObjectTypeProto_NewBuilder().SetUuid(jsaTEST_ID.AsByteString()).SetMarkedConstructor(true))
	// When
	colorPool := ColorPool_Builder_().AddShardAnd(typePool0, StringPool_Empty()).AddShardAnd(typePool1, StringPool_Empty()).Build()
	// Then
	jsaTrue(t, colorPool.GetColor(jsaTEST_ID).IsConstructor())
}

func TestColorPool_Reconcile_setInvalidating(t *testing.T) {
	// Given
	typePool0 := jsaSingleObjectPool(ObjectTypeProto_NewBuilder().SetUuid(jsaTEST_ID.AsByteString()).SetIsInvalidating(false))
	typePool1 := jsaSingleObjectPool(ObjectTypeProto_NewBuilder().SetUuid(jsaTEST_ID.AsByteString()).SetIsInvalidating(true))
	// When
	colorPool := ColorPool_Builder_().AddShardAnd(typePool0, StringPool_Empty()).AddShardAnd(typePool1, StringPool_Empty()).Build()
	// Then
	jsaTrue(t, colorPool.GetColor(jsaTEST_ID).IsInvalidating())
}

func TestColorPool_Reconcile_setPropertiesKeepOriginalName(t *testing.T) {
	// Given
	typePool0 := jsaSingleObjectPool(ObjectTypeProto_NewBuilder().SetUuid(jsaTEST_ID.AsByteString()).SetPropertiesKeepOriginalName(false))
	typePool1 := jsaSingleObjectPool(ObjectTypeProto_NewBuilder().SetUuid(jsaTEST_ID.AsByteString()).SetPropertiesKeepOriginalName(true))
	// When
	colorPool := ColorPool_Builder_().AddShardAnd(typePool0, StringPool_Empty()).AddShardAnd(typePool1, StringPool_Empty()).Build()
	// Then
	jsaTrue(t, colorPool.GetColor(jsaTEST_ID).GetPropertiesKeepOriginalName())
}

func TestColorPool_Reconcile_debugInfo_mismatches(t *testing.T) {
	// Given
	typePool0 := TypePool_NewBuilder().SetDebugInfoFromBuilder(TypePool_DebugInfo_NewBuilder().AddMismatch(TypePool_DebugInfo_Mismatch_NewBuilder().SetSourceRef("location_0").AddInvolvedColor(jsaPrim(PrimitiveType_SYMBOL_TYPE)).AddInvolvedColor(jsaPrim(PrimitiveType_STRING_TYPE)).Build()).AddMismatch(TypePool_DebugInfo_Mismatch_NewBuilder().SetSourceRef("location_1").AddInvolvedColor(jsaPrim(PrimitiveType_NUMBER_TYPE)).Build())).Build()
	typePool1 := TypePool_NewBuilder().SetDebugInfoFromBuilder(TypePool_DebugInfo_NewBuilder().AddMismatch(TypePool_DebugInfo_Mismatch_NewBuilder().SetSourceRef("location_0").AddInvolvedColor(jsaPrim(PrimitiveType_SYMBOL_TYPE)).AddInvolvedColor(jsaPrim(PrimitiveType_BOOLEAN_TYPE)).Build()).AddMismatch(TypePool_DebugInfo_Mismatch_NewBuilder().SetSourceRef("location_2").AddInvolvedColor(jsaPrim(PrimitiveType_BIGINT_TYPE)).Build())).Build()
	// When
	colorPool := ColorPool_Builder_().AddShardAnd(typePool0, StringPool_Empty()).AddShardAnd(typePool1, StringPool_Empty()).Build()
	// Then
	expected := jutil.NewSetMultimap[colors.ColorId, string]()
	expected.Put(colors.StandardColors_BOOLEAN.GetId(), "location_0")
	expected.Put(colors.StandardColors_STRING.GetId(), "location_0")
	expected.Put(colors.StandardColors_SYMBOL.GetId(), "location_0")
	expected.Put(colors.StandardColors_NUMBER.GetId(), "location_1")
	expected.Put(colors.StandardColors_BIGINT.GetId(), "location_2")
	jsaMultimapEq(t, colorPool.GetRegistry().GetMismatchLocationsForDebugging(), expected)
}

func TestColorPool_Reconcile_shards_preserveOriginalOffsets(t *testing.T) {
	// Given
	otherId := colors.ColorId_FromAscii("other")
	typePool0 := TypePool_NewBuilder().AddTypeFromBuilder(TypeProto_NewBuilder().SetObject(ObjectTypeProto_NewBuilder().SetUuid(otherId.AsByteString()).Build())).AddTypeFromBuilder(TypeProto_NewBuilder().SetObject(ObjectTypeProto_NewBuilder().SetUuid(jsaTEST_ID.AsByteString()).Build())).Build()
	typePool1 := TypePool_NewBuilder().AddTypeFromBuilder(TypeProto_NewBuilder().SetObject(ObjectTypeProto_NewBuilder().SetUuid(jsaTEST_ID.AsByteString()).Build())).AddTypeFromBuilder(TypeProto_NewBuilder().SetObject(ObjectTypeProto_NewBuilder().SetUuid(otherId.AsByteString()).Build())).Build()
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

func jsaSingleObjectPool(builder *ObjectTypeProto_Builder) *TypePool {
	return TypePool_NewBuilder().AddTypeFromBuilder(TypeProto_NewBuilder().SetObject(builder.Build())).Build()
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

func jsaMultimapEq[K any, V any](t *testing.T, got, want *jutil.SetMultimap[K, V]) {
	t.Helper()
	ok := got.Size() == want.Size()
	for k, v := range want.Entries() {
		ok = ok && got.Get(k).Contains(v)
	}
	if !ok {
		t.Errorf("got %v, want %v", got, want)
	}
}
