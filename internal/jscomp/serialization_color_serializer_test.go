// Port of test/com/google/javascript/jscomp/serialization/ColorSerializerTest.java @ closure-compiler v20261021

/*
 * Copyright 2021 The Closure Compiler Authors.
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
	"bytes"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func jsaCsAssertTypePool(t *testing.T, got, want *TypePool) {
	t.Helper()
	if got.String() != want.String() || !bytes.Equal(got.ToByteArray(), want.ToByteArray()) {
		t.Errorf("TypePool\n got:\n%s\nwant:\n%s", got, want)
	}
}

func TestColorSerializer_GenerateEmptyTypePool(t *testing.T) {
	// Expect an empty TypePool when no colors are added.
	jsaCsAssertTypePool(t,
		newJsaCsTester().
			Init().
			// don't add anything
			GenerateTypePool(t),
		// DebugInfo is present, but empty
		&TypePool{debugInfo_: &TypePool_DebugInfo{}})

	jsaCsAssertTypePool(t,
		newJsaCsTester().
			SetSerializationMode(SerializationOptions_Builder_().SetIncludeDebugInfo(false).Build()).
			Init().
			// don't add anything
			GenerateTypePool(t),
		TypePool_GetDefaultInstance())
}

func TestColorSerializer_OnlyAddAxiomaticColors(t *testing.T) {
	// These are the colors that are never actually serialized, even when requested.
	// They are always just assumed to be there when deserializing and to have fixed indices as
	// indicated by their index in TypePointers.OFFSET_TO_AXIOMATIC_COLOR.
	var axiomaticTestColors []*jsaCsTestColor
	for _, c := range TypePointers_OFFSET_TO_AXIOMATIC_COLOR {
		axiomaticTestColors = append(axiomaticTestColors, jsaCsGetAxiomaticTestColor(c))
	}

	// Expect an empty TypePool when only axiomatic colors are added.
	jsaCsAssertTypePool(t,
		newJsaCsTester().
			Init().
			// add all the axiomatic colors, but only those
			AddColorsList(t, axiomaticTestColors).
			GenerateTypePool(t),
		// No `type` values, because axiomatic colors are never serialized.
		// DebugInfo is empty, but present.
		&TypePool{debugInfo_: &TypePool_DebugInfo{}})
}

func TestColorSerializer_AddASuperTypeRelationship(t *testing.T) {
	superTypeColor := newJsaCsTestObjectColorBuilder().SetColorId("sup_type").SetTrimmedPoolOffset(0).Build()
	subTypeColor := newJsaCsTestObjectColorBuilder().SetColorId("sub_type").SetTrimmedPoolOffset(1).Build()
	expectedTypePool := &TypePool{
		type_: []*TypeProto{
			superTypeColor.GetExpectedTypeProto(),
			subTypeColor.GetExpectedTypeProto(),
		},
		disambiguationEdges_: []*SubtypingEdge{{
			subtype_:   int32(subTypeColor.expectedTypePointer),
			supertype_: int32(superTypeColor.expectedTypePointer),
		}},
		// empty DebugInfo
		debugInfo_: &TypePool_DebugInfo{},
	}

	jsaCsAssertTypePool(t,
		newJsaCsTester().
			Init().
			LinkSubColorToSuperColor(subTypeColor, superTypeColor).
			AddColor(t, superTypeColor).
			AddColor(t, subTypeColor).
			GenerateTypePool(t),
		expectedTypePool)
}

func TestColorSerializer_ImplicitAdditionOfSuperType(t *testing.T) {
	subTypeColor := newJsaCsTestObjectColorBuilder().
		SetColorId("sub_type").
		// The subtype will be serialized first, because it is explicitly added.
		SetTrimmedPoolOffset(0).
		Build()
	superTypeColor := newJsaCsTestObjectColorBuilder().
		SetColorId("sup_type").
		// The supertype will be serialized second, because it gets added implicitly
		// by its relationship to the subtype.
		SetTrimmedPoolOffset(1).
		Build()
	expectedTypePool := &TypePool{
		type_: []*TypeProto{
			subTypeColor.GetExpectedTypeProto(),
			superTypeColor.GetExpectedTypeProto(),
		},
		disambiguationEdges_: []*SubtypingEdge{{
			subtype_:   int32(subTypeColor.expectedTypePointer),
			supertype_: int32(superTypeColor.expectedTypePointer),
		}},
		// DebugInfo is present but empty.
		debugInfo_: &TypePool_DebugInfo{},
	}

	jsaCsAssertTypePool(t,
		newJsaCsTester().
			Init().
			LinkSubColorToSuperColor(subTypeColor, superTypeColor).
			// Only explicitly add the subtype
			AddColor(t, subTypeColor).
			GenerateTypePool(t),
		expectedTypePool)
}

func TestColorSerializer_SimulatedClasses(t *testing.T) {
	// Simulate Colors for this class:
	//
	// class Base {
	//   static baseStaticField;
	//   baseField;
	// }
	//
	// We will explicitly add the color representing `typeof Base`.
	// The colors for `Base.prototype` and the instances of `Base` will be added implicitly
	// in that order.
	baseFieldPooledString := jsaCsPooledString_Create("baseField", 1)
	basePooledString := jsaCsPooledString_Create("Base", 2)
	typeofBasePrototypePooledString := jsaCsPooledString_Create("typeof Base.prototype", 3)
	typeofBasePooledString := jsaCsPooledString_Create("typeof Base", 4)
	baseClassInstanceTestColor := newJsaCsTestObjectColorBuilder().
		SetColorId("Base.ins").
		// 0 => explicitly added base class constructor
		// 1 => explicitly added child class constructor
		// Then this instance color gets added implicitly
		SetTrimmedPoolOffset(2).
		AddOwnProperty(baseFieldPooledString).
		Build()
	baseClassPrototypeTestColor := newJsaCsTestObjectColorBuilder().
		SetColorId("Base.pro").
		// This is implicitly added after the instance color above.
		SetTrimmedPoolOffset(3).
		Build()
	baseStaticFieldPooledString := jsaCsPooledString_Create("baseStaticField", 0)
	baseClassConstructorTestColor := newJsaCsTestObjectColorBuilder().
		SetColorId("Base.con").
		// This is the first explicitly-added type.
		SetTrimmedPoolOffset(0).
		SetConstructor(true).
		AddPrototypeTestColor(baseClassPrototypeTestColor).
		AddInstanceTestColor(baseClassInstanceTestColor).
		AddOwnProperty(baseStaticFieldPooledString).
		Build()

	// class Child extends Base {
	//   static childStaticField;
	//   childField;
	// }
	// We will explicitly add the color representing `typeof Child`.
	// The colors for `Child.prototype` and the instances of `Child` will be added implicitly
	// in that order.
	childFieldPooledString := jsaCsPooledString_Create("childField", 5)
	childPooledString := jsaCsPooledString_Create("Child", 6)
	typeofChildPrototypePooledString := jsaCsPooledString_Create("typeof Child.prototype", 7)
	typeofChildPooledString := jsaCsPooledString_Create("typeof Child", 8)
	childClassInstanceTestColor := newJsaCsTestObjectColorBuilder().
		SetColorId("Child.in").
		// 0 => explicitly added base class constructor color
		// 1 => explicitly added child class constructor color
		// 2 => implicitly added base class instance color
		// 3 => implicitly added base class prototype color
		// Then this instance color gets added implicitly
		SetTrimmedPoolOffset(4).
		AddOwnProperty(childFieldPooledString).
		Build()
	childClassPrototypeTestColor := newJsaCsTestObjectColorBuilder().
		SetColorId("Child.pr").
		// This is implicitly added after the instance color above.
		SetTrimmedPoolOffset(5).
		Build()
	childStaticFieldPooledString := jsaCsPooledString_Create("childStaticField", 2)
	childClassConstructorTestColor := newJsaCsTestObjectColorBuilder().
		SetColorId("Child.co").
		// This is the second explicitly-added type
		SetTrimmedPoolOffset(1).
		SetConstructor(true).
		AddPrototypeTestColor(childClassPrototypeTestColor).
		AddInstanceTestColor(childClassInstanceTestColor).
		AddOwnProperty(childStaticFieldPooledString).
		Build()

	expectedTypePool := &TypePool{
		type_: []*TypeProto{
			baseClassConstructorTestColor.GetExpectedTypeProto(),
			childClassConstructorTestColor.GetExpectedTypeProto(),
			baseClassInstanceTestColor.GetExpectedTypeProto(),
			baseClassPrototypeTestColor.GetExpectedTypeProto(),
			childClassInstanceTestColor.GetExpectedTypeProto(),
			childClassPrototypeTestColor.GetExpectedTypeProto(),
		},
		disambiguationEdges_: []*SubtypingEdge{
			{
				subtype_:   int32(childClassConstructorTestColor.expectedTypePointer),
				supertype_: int32(baseClassConstructorTestColor.expectedTypePointer),
			},
			{
				subtype_:   int32(childClassInstanceTestColor.expectedTypePointer),
				supertype_: int32(baseClassInstanceTestColor.expectedTypePointer),
			},
			{
				subtype_:   int32(childClassPrototypeTestColor.expectedTypePointer),
				supertype_: int32(baseClassPrototypeTestColor.expectedTypePointer),
			},
		},
		// DebugInfo is present but empty.
		debugInfo_: &TypePool_DebugInfo{},
	}

	jsaCsAssertTypePool(t,
		newJsaCsTester().
			Init().
			AddPooledString(baseStaticFieldPooledString).
			AddPooledString(baseFieldPooledString).
			AddPooledString(childStaticFieldPooledString).
			AddPooledString(childFieldPooledString).
			AddPooledString(basePooledString).
			AddPooledString(typeofBasePrototypePooledString).
			AddPooledString(typeofBasePooledString).
			AddPooledString(childPooledString).
			AddPooledString(typeofChildPrototypePooledString).
			AddPooledString(typeofChildPooledString).
			// Only explicitly add the constructor types
			AddColor(t, baseClassConstructorTestColor).
			AddColor(t, childClassConstructorTestColor).
			LinkSubColorToSuperColor(childClassConstructorTestColor, baseClassConstructorTestColor).
			LinkSubColorToSuperColor(childClassInstanceTestColor, baseClassInstanceTestColor).
			LinkSubColorToSuperColor(childClassPrototypeTestColor, baseClassPrototypeTestColor).
			GenerateTypePool(t),
		expectedTypePool)
}

func TestColorSerializer_IncludeMismatches(t *testing.T) {
	testColor1 := newJsaCsTestObjectColorBuilder().SetColorId("color001").SetTrimmedPoolOffset(0).Build()
	testColor2 := newJsaCsTestObjectColorBuilder().SetColorId("color002").SetTrimmedPoolOffset(1).Build()
	testMismatch1 := jsaCsTestMismatch_Create("location1", testColor1, testColor2)

	testColor3 := newJsaCsTestObjectColorBuilder().SetColorId("color003").SetTrimmedPoolOffset(2).Build()
	testColor4 := newJsaCsTestObjectColorBuilder().SetColorId("color004").SetTrimmedPoolOffset(3).Build()
	testMismatch2 := jsaCsTestMismatch_Create("location2", testColor3, testColor4)

	expectedTypePool := &TypePool{
		type_: []*TypeProto{
			testColor1.GetExpectedTypeProto(),
			testColor2.GetExpectedTypeProto(),
			testColor3.GetExpectedTypeProto(),
			testColor4.GetExpectedTypeProto(),
		},
		debugInfo_: &TypePool_DebugInfo{
			mismatch_: []*TypePool_DebugInfo_Mismatch{
				testMismatch1.GetExpectedMismatch(),
				testMismatch2.GetExpectedMismatch(),
			},
		},
	}

	jsaCsAssertTypePool(t,
		newJsaCsTester().
			Init().
			AddColors(t, testColor1, testColor2, testColor3, testColor4).
			AddMismatch(testMismatch1).
			AddMismatch(testMismatch2).
			GenerateTypePool(t),
		expectedTypePool)
}

func TestColorSerializer_SkipMismatches(t *testing.T) {
	testColor1 := newJsaCsTestObjectColorBuilder().SetColorId("color001").SetTrimmedPoolOffset(0).Build()
	testColor2 := newJsaCsTestObjectColorBuilder().SetColorId("color002").SetTrimmedPoolOffset(1).Build()
	testMismatch1 := jsaCsTestMismatch_Create("location1", testColor1, testColor2)

	testColor3 := newJsaCsTestObjectColorBuilder().SetColorId("color003").SetTrimmedPoolOffset(2).Build()
	testColor4 := newJsaCsTestObjectColorBuilder().SetColorId("color004").SetTrimmedPoolOffset(3).Build()
	testMismatch2 := jsaCsTestMismatch_Create("location2", testColor3, testColor4)

	expectedTypePool := &TypePool{
		type_: []*TypeProto{
			testColor1.GetExpectedTypeProto(),
			testColor2.GetExpectedTypeProto(),
			testColor3.GetExpectedTypeProto(),
			testColor4.GetExpectedTypeProto(),
		},
		// No debug info
	}

	jsaCsAssertTypePool(t,
		newJsaCsTester().
			SetSerializationMode(SerializationOptions_Builder_().SetIncludeDebugInfo(false).Build()).
			Init().
			AddColors(t, testColor1, testColor2, testColor3, testColor4).
			AddMismatch(testMismatch1).
			AddMismatch(testMismatch2).
			GenerateTypePool(t),
		expectedTypePool)
}

func TestColorSerializer_FilterPropertyNames(t *testing.T) {
	// Tell ColorSerializer to only keep property names beginning with "kept"
	propertyNameFilter := func(propName string) bool { return strings.HasPrefix(propName, "kept") }

	// Build a type that has a property name that should be kept and one that should be
	// dropped.
	keptPropertyName := jsaCsPooledString_Create("keptProperty", 0)
	droppedPropertyName := jsaCsPooledString_Create("droppedProperty", 1)
	testColor := newJsaCsTestObjectColorBuilder().
		SetColorId("color001").
		SetTrimmedPoolOffset(0).
		AddOwnProperty(keptPropertyName).
		AddOwnProperty(droppedPropertyName).
		Build()

	actualTypePool := newJsaCsTester().
		SetPropertyNameFilter(propertyNameFilter).
		Init().
		AddPooledString(keptPropertyName).
		AddPooledString(droppedPropertyName).
		AddColor(t, testColor).
		GenerateTypePool(t)

	ownPropertyList := actualTypePool.GetType(0).GetObject().GetOwnPropertyList()
	jsaDeepEq(t, ownPropertyList, []int32{int32(keptPropertyName.poolOffset)})
}

func TestColorSerializer_AddUnionImplicitlyAddsMembers(t *testing.T) {
	testColor1 := newJsaCsTestObjectColorBuilder().
		SetColorId("color001").
		// We will explicitly add this color first
		SetTrimmedPoolOffset(1).
		Build()
	testColor2 := newJsaCsTestObjectColorBuilder().SetColorId("color002").SetTrimmedPoolOffset(2).Build()
	unionTestColor := newJsaCsTestUnionColorBuilder().
		SetTrimmedPoolOffset(0).
		AddTestColor(testColor1).
		AddTestColor(testColor2).
		Build()

	expectedTypePool := &TypePool{
		type_: []*TypeProto{
			unionTestColor.GetExpectedTypeProto(),
			testColor1.GetExpectedTypeProto(),
			testColor2.GetExpectedTypeProto(),
		},
		// empty DebugInfo
		debugInfo_: &TypePool_DebugInfo{},
	}

	jsaCsAssertTypePool(t,
		newJsaCsTester().
			Init().
			AddColor(t, unionTestColor).
			GenerateTypePool(t),
		expectedTypePool)
}

func TestColorSerializer_AvoidAddingDuplicateUnions(t *testing.T) {
	testColor1 := newJsaCsTestObjectColorBuilder().
		SetColorId("color001").
		// We will explicitly add this color first
		SetTrimmedPoolOffset(1).
		Build()
	testColor2 := newJsaCsTestObjectColorBuilder().SetColorId("color002").SetTrimmedPoolOffset(2).Build()
	unionTestColor := newJsaCsTestUnionColorBuilder().
		SetTrimmedPoolOffset(0).
		AddTestColor(testColor1).
		AddTestColor(testColor2).
		Build()
	// Creating a second union containing the same types will create a distinct
	// Java object, but it will have the same ColorId.
	duplicateUnionTestColor := newJsaCsTestUnionColorBuilder().
		SetTrimmedPoolOffset(0). // should end up getting the same offset
		AddTestColor(testColor1).
		AddTestColor(testColor2).
		Build()

	expectedTypePool := &TypePool{
		type_: []*TypeProto{
			unionTestColor.GetExpectedTypeProto(),
			// duplicate Union color should not be serialized
			testColor1.GetExpectedTypeProto(),
			testColor2.GetExpectedTypeProto(),
		},
		// empty DebugInfo
		debugInfo_: &TypePool_DebugInfo{},
	}

	jsaCsAssertTypePool(t,
		newJsaCsTester().
			Init().
			AddColor(t, unionTestColor).
			// Try to add the duplicate union color. It should be ignored.
			AddColor(t, duplicateUnionTestColor).
			GenerateTypePool(t),
		expectedTypePool)
}

// jsaCsPooledString represents a string that is stored in a StringPool,
// recording both its value and its offset.
type jsaCsPooledString struct {
	value      string
	poolOffset int
}

func jsaCsPooledString_Create(value string, poolOffset int) *jsaCsPooledString {
	return &jsaCsPooledString{value: value, poolOffset: poolOffset}
}

// jsaCsTestObjectColorBuilder builds objects to represent both an object
// Color to be added and the Integer and TypeProto we expect to be generated
// for it.
type jsaCsTestObjectColorBuilder struct {
	instanceTestColors  []*jsaCsTestColor
	prototypeTestColors []*jsaCsTestColor
	ownProperties       []*jsaCsPooledString

	colorId           colors.ColorId
	trimmedPoolOffset int
	isConstructor     bool
}

func newJsaCsTestObjectColorBuilder() *jsaCsTestObjectColorBuilder {
	return &jsaCsTestObjectColorBuilder{trimmedPoolOffset: -1}
}

func (b *jsaCsTestObjectColorBuilder) SetColorId(colorIdString string) *jsaCsTestObjectColorBuilder {
	jutil.CheckArgumentMsg1(len(colorIdString) <= 8, "color ID string too long: %s", colorIdString)
	b.colorId = colors.ColorId_FromAscii(colorIdString)
	return b
}

func (b *jsaCsTestObjectColorBuilder) SetTrimmedPoolOffset(trimmedPoolOffset int) *jsaCsTestObjectColorBuilder {
	jutil.CheckStateMsg1(trimmedPoolOffset >= 0, "invalid trimmedPoolOffset: %s", trimmedPoolOffset)
	b.trimmedPoolOffset = trimmedPoolOffset
	return b
}

func (b *jsaCsTestObjectColorBuilder) SetConstructor(constructor bool) *jsaCsTestObjectColorBuilder {
	b.isConstructor = constructor
	return b
}

func (b *jsaCsTestObjectColorBuilder) AddInstanceTestColor(instanceTestColor *jsaCsTestColor) *jsaCsTestObjectColorBuilder {
	b.instanceTestColors = append(b.instanceTestColors, instanceTestColor)
	return b
}

func (b *jsaCsTestObjectColorBuilder) AddPrototypeTestColor(prototypeTestColor *jsaCsTestColor) *jsaCsTestObjectColorBuilder {
	b.prototypeTestColors = append(b.prototypeTestColors, prototypeTestColor)
	return b
}

func (b *jsaCsTestObjectColorBuilder) AddOwnProperty(ownProperty *jsaCsPooledString) *jsaCsTestObjectColorBuilder {
	b.ownProperties = append(b.ownProperties, ownProperty)
	return b
}

func (b *jsaCsTestObjectColorBuilder) Build() *jsaCsTestColor {
	jutil.CheckStateMsg(b.trimmedPoolOffset >= 0, "call setTrimmedPoolOffset() first")
	// There's no testing benefit to varying these values
	propertiesKeepOriginalName := false
	isInvalidating := false
	isClosureAssert := false

	instanceColors := make([]*colors.Color, 0, len(b.instanceTestColors))
	var instancePointers []int32
	for _, c := range b.instanceTestColors {
		instanceColors = append(instanceColors, c.color)
		instancePointers = append(instancePointers, int32(c.expectedTypePointer))
	}
	prototypes := make([]*colors.Color, 0, len(b.prototypeTestColors))
	var prototypePointers []int32
	for _, c := range b.prototypeTestColors {
		prototypes = append(prototypes, c.color)
		prototypePointers = append(prototypePointers, int32(c.expectedTypePointer))
	}
	ownProperties := make([]string, 0, len(b.ownProperties))
	var ownPropertyOffsets []int32
	for _, p := range b.ownProperties {
		ownProperties = append(ownProperties, p.value)
		ownPropertyOffsets = append(ownPropertyOffsets, int32(p.poolOffset))
	}
	colorBuilder := colors.Color_SingleBuilder().
		SetId(b.colorId).
		SetInvalidating(isInvalidating).
		SetPropertiesKeepOriginalName(propertiesKeepOriginalName).
		SetConstructor(b.isConstructor).
		SetClosureAssert(isClosureAssert).
		SetInstanceColors(instanceColors).
		SetPrototypes(prototypes).
		SetOwnProperties(ownProperties)
	typePointer := TypePointers_UntrimOffset(b.trimmedPoolOffset)
	typeProto := testObjectTypeProto(&ObjectTypeProto{
		uuid_:                       b.colorId.AsByteString(),
		isInvalidating_:             isInvalidating,
		propertiesKeepOriginalName_: propertiesKeepOriginalName,
		markedConstructor_:          b.isConstructor,
		closureAssert_:              isClosureAssert,
		instanceType_:               instancePointers,
		prototype_:                  prototypePointers,
		ownProperty_:                ownPropertyOffsets,
	})
	return jsaCsTestColor_Create(colorBuilder.Build(), typeProto, typePointer)
}

type jsaCsTestUnionColorBuilder struct {
	memberTestColors  []*jsaCsTestColor
	trimmedPoolOffset int
}

func newJsaCsTestUnionColorBuilder() *jsaCsTestUnionColorBuilder {
	return &jsaCsTestUnionColorBuilder{trimmedPoolOffset: -1}
}

func (b *jsaCsTestUnionColorBuilder) SetTrimmedPoolOffset(trimmedPoolOffset int) *jsaCsTestUnionColorBuilder {
	jutil.CheckStateMsg1(trimmedPoolOffset >= 0, "invalid trimmedPoolOffset: %s", trimmedPoolOffset)
	b.trimmedPoolOffset = trimmedPoolOffset
	return b
}

func (b *jsaCsTestUnionColorBuilder) AddTestColor(testColor *jsaCsTestColor) *jsaCsTestUnionColorBuilder {
	b.memberTestColors = append(b.memberTestColors, testColor)
	return b
}

func (b *jsaCsTestUnionColorBuilder) Build() *jsaCsTestColor {
	jutil.CheckStateMsg(b.trimmedPoolOffset >= 0, "call setTrimmedPoolOffset() first")
	memberColors := make([]*colors.Color, 0, len(b.memberTestColors))
	var memberTypePointers []int32
	for _, c := range b.memberTestColors {
		memberColors = append(memberColors, c.color)
		memberTypePointers = append(memberTypePointers, int32(c.expectedTypePointer))
	}
	color := colors.Color_CreateUnionSlice(memberColors)

	typePointer := TypePointers_UntrimOffset(b.trimmedPoolOffset)
	typeProto := testUnionTypeProto(memberTypePointers...)
	return jsaCsTestColor_Create(color, typeProto, typePointer)
}

func jsaCsGetAxiomaticTestColor(axiomaticColor *colors.Color) *jsaCsTestColor {
	poolOffset := -1
	for i, c := range TypePointers_OFFSET_TO_AXIOMATIC_COLOR {
		if c.Equals(axiomaticColor) {
			poolOffset = i
			break
		}
	}
	jutil.CheckArgumentMsg1(poolOffset >= 0, "Not an axiomatic color: %s", axiomaticColor)
	return jsaCsTestColor_Create(axiomaticColor, nil, poolOffset)
}

// jsaCsTestColor represents a Color that has been or will be added to the
// ColorSerializer, with the TypeProto (nil for an axiomatic color) and the
// pointer we expect ColorSerializer to create for it.
type jsaCsTestColor struct {
	color                     *colors.Color
	nullableExpectedTypeProto *TypeProto
	expectedTypePointer       int
}

func (c *jsaCsTestColor) GetExpectedTypeProto() *TypeProto {
	return jutil.CheckNotNull(c.nullableExpectedTypeProto)
}

func jsaCsTestColor_Create(color *colors.Color, expectedTypeProto *TypeProto, expectedTypePointer int) *jsaCsTestColor {
	return &jsaCsTestColor{color: color, nullableExpectedTypeProto: expectedTypeProto, expectedTypePointer: expectedTypePointer}
}

// jsaCsTestMismatch represents a color/type mismatch both as it would appear
// in a ColorRegistry and in a TypePool.
type jsaCsTestMismatch struct {
	locationString string
	testColors     []*jsaCsTestColor
}

func (m *jsaCsTestMismatch) GetColors() []*colors.Color {
	var out []*colors.Color
	for _, c := range m.testColors {
		out = append(out, c.color)
	}
	return out
}

func (m *jsaCsTestMismatch) GetExpectedMismatch() *TypePool_DebugInfo_Mismatch {
	var involved []int32
	for _, c := range m.testColors {
		involved = append(involved, int32(c.expectedTypePointer))
	}
	return &TypePool_DebugInfo_Mismatch{
		sourceRef_:     m.locationString,
		involvedColor_: involved,
	}
}

func jsaCsTestMismatch_Create(locationString string, testColors ...*jsaCsTestColor) *jsaCsTestMismatch {
	return &jsaCsTestMismatch{locationString: locationString, testColors: testColors}
}

// jsaCsTester drives a ColorSerializer.
type jsaCsTester struct {
	// Simulates ColorRegistry#getDisambiguationSuperTypes()
	subColorToSuperColorsMap *jutil.ObjectLinkedHashMap[*colors.Color, *jutil.ObjectLinkedHashSet[*colors.Color]]
	// Simulates ColorRegistry#getMismatchLocationsForDebugging()
	colorToMismatchLocationStringsMap *jutil.ObjectLinkedHashMap[*colors.Color, *jutil.LinkedHashSet[string]]
	// Simulates StringPool.Builder#put()
	stringToPoolOffsetMap map[string]int
	// Test object
	colorSerializer *ColorSerializer
	// Default to the serialization mode that includes the most information.
	serializationMode *SerializationOptions
	// Tells ColorSerializer which property names to serialize
	propertyNameFilter func(string) bool
}

func newJsaCsTester() *jsaCsTester {
	return &jsaCsTester{
		subColorToSuperColorsMap:          jutil.NewObjectLinkedHashMap[*colors.Color, *jutil.ObjectLinkedHashSet[*colors.Color]](),
		colorToMismatchLocationStringsMap: jutil.NewObjectLinkedHashMap[*colors.Color, *jutil.LinkedHashSet[string]](),
		stringToPoolOffsetMap:             map[string]int{},
		serializationMode:                 SerializationOptions_Builder_().SetIncludeDebugInfo(true).SetRunValidation(true).Build(),
		propertyNameFilter:                func(string) bool { return true },
	}
}

func (s *jsaCsTester) SetSerializationMode(serializationMode *SerializationOptions) *jsaCsTester {
	jutil.CheckStateMsg(s.colorSerializer == nil, "call this method before init()")
	s.serializationMode = serializationMode
	return s
}

func (s *jsaCsTester) SetPropertyNameFilter(propertyNameFilter func(string) bool) *jsaCsTester {
	jutil.CheckStateMsg(s.colorSerializer == nil, "call this method before init()")
	s.propertyNameFilter = propertyNameFilter
	return s
}

func (s *jsaCsTester) Init() *jsaCsTester {
	jutil.CheckStateMsg(s.colorSerializer == nil, "init() called more than once")
	s.colorSerializer = NewColorSerializer(
		s.serializationMode,
		func(str string) int {
			offset, ok := s.stringToPoolOffsetMap[str]
			jutil.CheckStateMsg1(ok, "unexpected string: %s", str)
			return offset
		},
		s.propertyNameFilter)
	return s
}

func (s *jsaCsTester) AddPooledString(pooledString *jsaCsPooledString) *jsaCsTester {
	_, dup := s.stringToPoolOffsetMap[pooledString.value]
	jutil.CheckStateMsg1(!dup, "duplicate string added: %s", pooledString.value)
	s.stringToPoolOffsetMap[pooledString.value] = pooledString.poolOffset
	return s
}

func (s *jsaCsTester) LinkSubColorToSuperColor(subColor, superColor *jsaCsTestColor) *jsaCsTester {
	set, ok := s.subColorToSuperColorsMap.Lookup(subColor.color)
	if !ok {
		set = jutil.NewObjectLinkedHashSet[*colors.Color]()
		s.subColorToSuperColorsMap.Put(subColor.color, set)
	}
	set.Add(superColor.color)
	return s
}

func (s *jsaCsTester) AddMismatch(testMismatch *jsaCsTestMismatch) *jsaCsTester {
	for _, color := range testMismatch.GetColors() {
		set, ok := s.colorToMismatchLocationStringsMap.Lookup(color)
		if !ok {
			set = jutil.NewLinkedHashSet[string]()
			s.colorToMismatchLocationStringsMap.Put(color, set)
		}
		set.Add(testMismatch.locationString)
	}
	return s
}

func (s *jsaCsTester) AddColor(t *testing.T, testColor *jsaCsTestColor) *jsaCsTester {
	t.Helper()
	typePointer := s.colorSerializer.AddColor(testColor.color)
	jsaEq(t, typePointer, testColor.expectedTypePointer)
	return s
}

func (s *jsaCsTester) AddColors(t *testing.T, testColors ...*jsaCsTestColor) *jsaCsTester {
	t.Helper()
	return s.AddColorsList(t, testColors)
}

func (s *jsaCsTester) AddColorsList(t *testing.T, testColorList []*jsaCsTestColor) *jsaCsTester {
	t.Helper()
	jutil.CheckNotNullMsg(s.colorSerializer, "call init() first")
	var cs []*colors.Color
	for _, c := range testColorList {
		cs = append(cs, c.color)
	}
	typePointers := s.colorSerializer.AddColors(cs)
	for i := 0; i < len(testColorList); i++ {
		jsaEq(t, typePointers[i], testColorList[i].expectedTypePointer)
	}
	return s
}

func (s *jsaCsTester) GenerateTypePool(t *testing.T) *TypePool {
	jutil.CheckNotNullMsg(s.colorSerializer, "call init() first")
	return s.colorSerializer.GenerateTypePool(
		func(subtypeColor *colors.Color) []*colors.Color {
			if set, ok := s.subColorToSuperColorsMap.Lookup(subtypeColor); ok {
				return set.ToSlice()
			}
			return nil
		},
		func(color *colors.Color) *jutil.LinkedHashSet[string] {
			if set, ok := s.colorToMismatchLocationStringsMap.Lookup(color); ok {
				return set
			}
			return jutil.NewLinkedHashSet[string]()
		})
}
