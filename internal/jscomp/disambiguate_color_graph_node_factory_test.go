// Port of com/google/javascript/jscomp/disambiguate/ColorGraphNodeFactoryTest.java @ closure-compiler v20261021

/*
 * Copyright 2019 The Closure Compiler Authors.
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

func newTestColorRegistry() *colors.ColorRegistry {
	return colors.ColorRegistry_Builder_().SetDefaultNativeColorsForTesting().Build()
}

func colorSetOf(cs ...*colors.Color) *jutil.ObjectLinkedHashSet[*colors.Color] {
	return jutil.NewObjectLinkedHashSet(cs...)
}

func TestColorGraphNodeFactory_TopLikeTypesFlattenToTop(t *testing.T) {
	colorRegistry := newTestColorRegistry()
	// Given
	factory := ColorGraphNodeFactory_CreateFactory(colorRegistry)
	flatTop := factory.CreateNode(colors.StandardColors_UNKNOWN)

	// When & Then
	assertEqual(t, "null", factory.CreateNode(nil), flatTop)
	assertEqual(t, "null or void", factory.CreateNode(colors.StandardColors_NULL_OR_VOID), flatTop)
}

func TestColorGraphNodeFactory_PrimitiveTypesFlattenToBoxedType(t *testing.T) {
	colorRegistry := newTestColorRegistry()
	factory := ColorGraphNodeFactory_CreateFactory(colorRegistry)
	flatString := factory.CreateNode(colors.StandardColors_STRING)
	flatStringObject := factory.CreateNode(colorRegistry.Get(colors.StandardColors_STRING_OBJECT_ID))
	flatTop := factory.CreateNode(colors.StandardColors_UNKNOWN)

	if flatString == flatTop {
		t.Error("string flattened to top")
	}
	assertEqual(t, "string object", flatString, flatStringObject)

	assertEqual(t, "union",
		factory.CreateNode(
			colors.Color_CreateUnion(colorSetOf(colors.StandardColors_STRING, colors.StandardColors_NUMBER))),
		factory.CreateNode(
			colors.Color_CreateUnion(
				colorSetOf(
					colorRegistry.Get(colors.StandardColors_STRING_OBJECT_ID),
					colorRegistry.Get(colors.StandardColors_NUMBER_OBJECT_ID)))))
}

func TestColorGraphNodeFactory_UnionTypesWhenFlattenedDropNullAndUndefined(t *testing.T) {
	colorRegistry := newTestColorRegistry()
	// Given
	factory := ColorGraphNodeFactory_CreateFactory(colorRegistry)

	nullOrVoidOrNumberType :=
		colors.Color_CreateUnion(colorSetOf(colors.StandardColors_NULL_OR_VOID, colors.StandardColors_NUMBER))

	flatNumber := factory.CreateNode(colors.StandardColors_NUMBER)

	// Given
	flatNullOrVoidOrNumber := factory.CreateNode(nullOrVoidOrNumberType)

	// Then
	assertEqual(t, "number", flatNullOrVoidOrNumber, flatNumber)
}

func TestColorGraphNodeFactory_ColorGraphNodesAreCached(t *testing.T) {
	colorRegistry := newTestColorRegistry()
	// Given
	factory := ColorGraphNodeFactory_CreateFactory(colorRegistry)

	fooType := colors.Color_SingleBuilder().SetId(colors.ColorId_FromAscii("Foo")).Build()
	flatFoo := factory.CreateNode(fooType)

	// When
	flatFooAgain := factory.CreateNode(fooType)

	// Then
	assertEqual(t, "cached", flatFooAgain, flatFoo)
}

func sampleColorsForFactoryTest() []*colors.Color {
	return colorSetOf(
		colors.StandardColors_NUMBER,
		colors.Color_SingleBuilder().SetId(colors.ColorId_FromAscii("Foo")).Build(),
		colors.Color_SingleBuilder().SetId(colors.ColorId_FromAscii("Bar")).Build()).ToSlice()
}

func flatSamples(factory *ColorGraphNodeFactory, sampleTypes []*colors.Color) []*ColorGraphNode {
	out := jutil.NewLinkedHashSet[*ColorGraphNode]()
	for _, c := range sampleTypes {
		out.Add(factory.CreateNode(c))
	}
	return out.ToSlice()
}

func TestColorGraphNodeFactory_ColorGraphNodesAreTracked(t *testing.T) {
	colorRegistry := newTestColorRegistry()
	// Given
	factory := ColorGraphNodeFactory_CreateFactory(colorRegistry)

	// When
	samples := flatSamples(factory, sampleColorsForFactoryTest())
	allKnownTypes := factory.GetAllKnownTypes()

	// Then
	for _, s := range samples {
		if !allKnownTypes.Contains(s) {
			t.Errorf("%s not known", s)
		}
	}
}

func TestColorGraphNodeFactory_ColorGraphNodeIndicesAreUnique(t *testing.T) {
	colorRegistry := newTestColorRegistry()
	// Given
	factory := ColorGraphNodeFactory_CreateFactory(colorRegistry)

	// When
	samples := flatSamples(factory, sampleColorsForFactoryTest())

	// Then
	seen := map[int]bool{}
	for _, s := range samples {
		if seen[s.GetIndex()] {
			t.Errorf("duplicate index %d", s.GetIndex())
		}
		seen[s.GetIndex()] = true
	}
}

func TestColorGraphNodeFactory_ColorGraphNodeIndicesAreDeterministic(t *testing.T) {
	colorRegistry := newTestColorRegistry()
	// Given
	sampleA := flatSamples(ColorGraphNodeFactory_CreateFactory(colorRegistry), sampleColorsForFactoryTest())
	sampleB := flatSamples(ColorGraphNodeFactory_CreateFactory(colorRegistry), sampleColorsForFactoryTest())

	// Then
	assertEqual(t, "size", len(sampleA), len(sampleB))
	for i := range sampleA {
		assertEqual(t, "index", sampleA[i].GetIndex(), sampleB[i].GetIndex())
		for _, b := range sampleB {
			if sampleA[i] == b {
				t.Errorf("%s is shared", b)
			}
		}
	}
}
