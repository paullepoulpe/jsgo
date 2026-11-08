// Port of com/google/javascript/jscomp/PolyfillUsageFinderTest.java @ closure-compiler v20261104
// (test/com/google/javascript/jscomp/PolyfillUsageFinderTest.java)

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
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/testjunit"
)

// PolyfillUsageFinderTest holds the unit tests for PolyfillUsageFinder.
type PolyfillUsageFinderTest struct {
	t *testing.T
}

func (c *PolyfillUsageFinderTest) JUnitSetT(t *testing.T) { c.t = t }
func (c *PolyfillUsageFinderTest) JUnitBefore()           {}
func (c *PolyfillUsageFinderTest) JUnitAfter()            {}

func TestPolyfillUsageFinderTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return &PolyfillUsageFinderTest{} })
}

// TestBasicCoverageTest covers testing of most common cases.
//
//   - guarded & unguarded references
//   - global symbols & prototype methods
//   - non-polyfill references
//   - references via a named global object
//
// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: PolyfillUsageFinderTest.basicCoverageTest()
func (c *PolyfillUsageFinderTest) TestBasicCoverageTest() {
	testPolyfillUsageFinder :=
		PolyfillUsageFinderTest_TestPolyfillUsageFinder_Builder_().
			WithPolyfillTableLines(
				`Array.prototype.fill es6 es3 es6/array/fill
Array.from es6 es3 es6/array/from
Map es6 es3 es6/map
`).
			ForSourceLines(
				`use(new Map()); // 1
globalThis.Array.from // 2
  ? globalThis.Array.from([]) // 3
  : []; // 4
[1, 2, 3].fill(0); // 5
notAPolyfill(); // 6
`).
			Build()

	unguardedPolyfillUsages := testPolyfillUsageFinder.GetUnguardedUsages()

	// Map() and [].fill() references
	AssertThat(unguardedPolyfillUsages).HasSize(2)

	// Source line 1: `use(new Map())`
	mapUsageSubject :=
		PolyfillUsageFinderTest_PolyfillUsageSubject_AssertPolyfillUsage(unguardedPolyfillUsages[0])
	mapUsageSubject.HasNodeThat().IsName("Map").HasLineno(1)
	mapUsageSubject.HasName("Map").IsNotExplicitGlobal()
	mapUsageSubject.
		HasPolyfillThat().
		HasNativeSymbol("Map").
		HasNativeVersionStr("es6").
		HasPolyfillVersionStr("es3").
		HasLibrary("es6/map").
		HasKind(PolyfillUsageFinder_Polyfill_Kind_STATIC)

	// Source line 5: `[1, 2, 3].fill(0);`
	arrayFillUsageSubject :=
		PolyfillUsageFinderTest_PolyfillUsageSubject_AssertPolyfillUsage(unguardedPolyfillUsages[1])
	arrayFillUsageSubject.HasNodeThat().IsGetProp().HasLineno(5)
	arrayFillUsageSubject.HasName("fill").IsNotExplicitGlobal()
	arrayFillUsageSubject.
		HasPolyfillThat().
		HasNativeSymbol("Array.prototype.fill").
		HasNativeVersionStr("es6").
		HasPolyfillVersionStr("es3").
		HasLibrary("es6/array/fill").
		HasKind(PolyfillUsageFinder_Polyfill_Kind_METHOD)

	// The 2 references to `Array.from` on line 2 and 3
	// Note that the guard expression itself counts as a guarded use.
	// Source line 2: `globalThis.Array.from ? globalThis.Array.from([]) : [];`
	guardedPolyfillUsages := testPolyfillUsageFinder.GetGuardedUsages()
	AssertThat(guardedPolyfillUsages).HasSize(2)

	hookConditionUseSubject :=
		PolyfillUsageFinderTest_PolyfillUsageSubject_AssertPolyfillUsage(guardedPolyfillUsages[0])
	hookConditionUseSubject.HasNodeThat().IsGetProp().HasLineno(2)

	arrayDotFromUseSubject :=
		PolyfillUsageFinderTest_PolyfillUsageSubject_AssertPolyfillUsage(guardedPolyfillUsages[1])
	arrayDotFromUseSubject.HasNodeThat().IsGetProp().HasLineno(3)

	// Other than the node they reference, both usage objects are the same
	for _, arrayDotFromUsageSubject := range []*PolyfillUsageFinderTest_PolyfillUsageSubject{hookConditionUseSubject, arrayDotFromUseSubject} {
		arrayDotFromUsageSubject.HasName("Array.from")
		// reference via `globalThis` makes these explicit global references
		arrayDotFromUsageSubject.IsExplicitGlobal()
		arrayDotFromUsageSubject.
			HasPolyfillThat().
			HasNativeSymbol("Array.from").
			HasNativeVersionStr("es6").
			HasPolyfillVersionStr("es3").
			HasLibrary("es6/array/from").
			HasKind(PolyfillUsageFinder_Polyfill_Kind_STATIC)
	}
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: PolyfillUsageFinderTest.shadowedPolyfillTest()
func (c *PolyfillUsageFinderTest) TestShadowedPolyfillTest() {
	testPolyfillUsageFinder :=
		PolyfillUsageFinderTest_TestPolyfillUsageFinder_Builder_().
			WithPolyfillTableLines(
				`Promise es6 es3 es6/promise/promise
Promise.allSettled es_2020 es3 es6/promise/allsettled
`).
			ForSourceLines(
				`// usage of Promise and of Promise.allSettled
Promise.allSettled([])(result => console.log(x));
function foo() {
  // These aren't references to Promise or Promise.allSettled because of the
  // shadowing variable
  let Promise = { allSettled: 'I solemnly swear that I am up to no good.' };
  console.log(Promise.allSettled);
}
`).
			Build()

	allUsages := testPolyfillUsageFinder.GetAllUsages()
	AssertThat(allUsages).HasSize(2)

	promiseUsageSubject :=
		PolyfillUsageFinderTest_PolyfillUsageSubject_AssertPolyfillUsage(allUsages[0])
	promiseUsageSubject. //
				HasNodeThat().
				IsName("Promise").
				HasLineno(2)
	promiseUsageSubject.
		HasPolyfillThat().
		HasNativeSymbol("Promise").
		HasNativeVersionStr("es6").
		HasPolyfillVersionStr("es3").
		HasLibrary("es6/promise/promise").
		HasKind(PolyfillUsageFinder_Polyfill_Kind_STATIC)

	allSettledUsageSubject :=
		PolyfillUsageFinderTest_PolyfillUsageSubject_AssertPolyfillUsage(allUsages[1])

	allSettledUsageSubject. //
				HasNodeThat().
				MatchesQualifiedName("Promise.allSettled").
				HasLineno(2)
	allSettledUsageSubject.
		HasPolyfillThat().
		HasNativeSymbol("Promise.allSettled").
		HasNativeVersionStr("es_2020").
		HasPolyfillVersionStr("es3").
		HasLibrary("es6/promise/allsettled").
		HasKind(PolyfillUsageFinder_Polyfill_Kind_STATIC)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: PolyfillUsageFinderTest.optChainCoverageTest()
func (c *PolyfillUsageFinderTest) TestOptChainCoverageTest() {
	testPolyfillUsageFinder :=
		PolyfillUsageFinderTest_TestPolyfillUsageFinder_Builder_().
			WithPolyfillTableLines("Array.from es6 es3 es6/array/from").
			ForSourceLines(
				"// Checks for the existence of `Array`, but not `.from`, so not guarded\n" +
					"Array?.from([]) ?? [];\n" +
					"// Correctly guarded on existence of `Array.from`\n" +
					"Array.from?.([]) ?? [];\n").
			Build()

	// First usage is unguarded
	unguardedPolyfillUsages := testPolyfillUsageFinder.GetUnguardedUsages()
	AssertThat(unguardedPolyfillUsages).HasSize(1)
	unguardedUsageSubject :=
		PolyfillUsageFinderTest_PolyfillUsageSubject_AssertPolyfillUsage(unguardedPolyfillUsages[0])
	unguardedUsageSubject.HasNodeThat().HasToken(rhino.Token_OPTCHAIN_GETPROP).HasLineno(2)

	// Second usage is guarded
	guardedPolyfillUsages := testPolyfillUsageFinder.GetGuardedUsages()
	AssertThat(guardedPolyfillUsages).HasSize(1)
	guardedUsageSubject :=
		PolyfillUsageFinderTest_PolyfillUsageSubject_AssertPolyfillUsage(guardedPolyfillUsages[0])
	guardedUsageSubject.HasNodeThat().HasToken(rhino.Token_GETPROP).HasLineno(4)

	// Other than the node they reference, both usage objects are the same
	for _, arrayDotFromUsageSubject := range []*PolyfillUsageFinderTest_PolyfillUsageSubject{unguardedUsageSubject, guardedUsageSubject} {
		arrayDotFromUsageSubject.HasName("Array.from")
		// reference via `globalThis` makes these explicit global references
		arrayDotFromUsageSubject.IsNotExplicitGlobal()
		arrayDotFromUsageSubject.
			HasPolyfillThat().
			HasNativeSymbol("Array.from").
			HasNativeVersionStr("es6").
			HasPolyfillVersionStr("es3").
			HasLibrary("es6/array/from").
			HasKind(PolyfillUsageFinder_Polyfill_Kind_STATIC)
	}
}

func polyfillUsageFinderTest_createFinderWithRegistry(registry *colors.ColorRegistry) *PolyfillUsageFinder {
	compiler := NewCompiler()
	compiler.colorRegistry = registry
	return NewPolyfillUsageFinder(
		compiler,
		PolyfillUsageFinder_Polyfills_FromTable("Iterator.prototype.toArray es_next es3 es6/iterator/toarray"),
		/* disambiguateProperties= */ true)
}

// Java: PolyfillUsageFinderTest.isDisjointFromIterator_unknownAndTopObject_returnsFalse()
func (c *PolyfillUsageFinderTest) TestIsDisjointFromIterator_unknownAndTopObject_returnsFalse() {
	registry := colors.NewColorRegistryBuilder().SetDefaultNativeColorsForTesting().Build()
	finder := polyfillUsageFinderTest_createFinderWithRegistry(registry)

	AssertThat(finder.IsDisjointFromIterator(colors.StandardColors_UNKNOWN)).IsFalse()
	AssertThat(finder.IsDisjointFromIterator(colors.StandardColors_TOP_OBJECT)).IsFalse()
	AssertThat(finder.IsDisjointFromIterator(colors.StandardColors_TOP_FUNCTION)).IsFalse()
}

// Java: PolyfillUsageFinderTest.isDisjointFromIterator_sameColor_returnsFalse()
func (c *PolyfillUsageFinderTest) TestIsDisjointFromIterator_sameColor_returnsFalse() {
	registry := colors.NewColorRegistryBuilder().SetDefaultNativeColorsForTesting().Build()
	finder := polyfillUsageFinderTest_createFinderWithRegistry(registry)
	iterator := registry.Get(colors.StandardColors_ITERATOR_ID)

	AssertThat(finder.IsDisjointFromIterator(iterator)).IsFalse()
}

// Java: PolyfillUsageFinderTest.isDisjointFromIterator_unrelatedColors_returnsTrue()
func (c *PolyfillUsageFinderTest) TestIsDisjointFromIterator_unrelatedColors_returnsTrue() {
	registry := colors.NewColorRegistryBuilder().SetDefaultNativeColorsForTesting().Build()
	finder := polyfillUsageFinderTest_createFinderWithRegistry(registry)
	array := registry.Get(colors.StandardColors_ARRAY_ID)
	customFoo := colors.Color_SingleBuilder().SetId(colors.ColorId_FromAscii("Foo")).Build()

	AssertThat(finder.IsDisjointFromIterator(array)).IsTrue()
	AssertThat(finder.IsDisjointFromIterator(customFoo)).IsTrue()
}

// Java: PolyfillUsageFinderTest.isDisjointFromIterator_generatorAndIterator_returnsFalse()
func (c *PolyfillUsageFinderTest) TestIsDisjointFromIterator_generatorAndIterator_returnsFalse() {
	registry := colors.NewColorRegistryBuilder().SetDefaultNativeColorsForTesting().Build()
	finder := polyfillUsageFinderTest_createFinderWithRegistry(registry)
	generator := registry.Get(colors.StandardColors_GENERATOR_ID)

	// Generator is a sibling of Iterator in externs, but inherits from %IteratorPrototype%
	// at runtime in ES2025, so Generator and Iterator are known not to be disjoint.
	AssertThat(finder.IsDisjointFromIterator(generator)).IsFalse()
}

// Java: PolyfillUsageFinderTest.isDisjointFromIterator_subtypesAndSupertypes_returnsFalse()
func (c *PolyfillUsageFinderTest) TestIsDisjointFromIterator_subtypesAndSupertypes_returnsFalse() {
	iterator := colors.Color_SingleBuilder().SetId(colors.StandardColors_ITERATOR_ID).Build()
	subIterator := colors.Color_SingleBuilder().SetId(colors.ColorId_FromAscii("SubIter")).Build()
	subSubIterator := colors.Color_SingleBuilder().SetId(colors.ColorId_FromAscii("SubSub")).Build()
	superIterable := colors.Color_SingleBuilder().SetId(colors.ColorId_FromAscii("SuperIt")).Build()
	unrelated := colors.Color_SingleBuilder().SetId(colors.ColorId_FromAscii("Unrel")).Build()

	registry := colors.NewColorRegistryBuilder().
		SetDefaultNativeColorsForTesting().
		AddDisambiguationEdge(subIterator, iterator).
		AddDisambiguationEdge(subSubIterator, subIterator).
		AddDisambiguationEdge(iterator, superIterable).
		Build()
	finder := polyfillUsageFinderTest_createFinderWithRegistry(registry)

	// Direct subtype: SubIterator extends Iterator
	AssertThat(finder.IsDisjointFromIterator(subIterator)).IsFalse()

	// Indirect subtype: SubSubIterator extends SubIterator extends Iterator
	AssertThat(finder.IsDisjointFromIterator(subSubIterator)).IsFalse()

	// Supertype: Iterator extends SuperIterable (receiver is supertype)
	AssertThat(finder.IsDisjointFromIterator(superIterable)).IsFalse()

	// Unrelated type with no edges
	AssertThat(finder.IsDisjointFromIterator(unrelated)).IsTrue()
}

// Java: PolyfillUsageFinderTest.isDisjointFromIterator_branchingAndSiblingHierarchies()
func (c *PolyfillUsageFinderTest) TestIsDisjointFromIterator_branchingAndSiblingHierarchies() {
	iterator := colors.Color_SingleBuilder().SetId(colors.StandardColors_ITERATOR_ID).Build()
	subIteratorA := colors.Color_SingleBuilder().SetId(colors.ColorId_FromAscii("SubItA")).Build()
	subIteratorB := colors.Color_SingleBuilder().SetId(colors.ColorId_FromAscii("SubItB")).Build()
	baseClass := colors.Color_SingleBuilder().SetId(colors.ColorId_FromAscii("Base")).Build()
	childClass := colors.Color_SingleBuilder().SetId(colors.ColorId_FromAscii("Child")).Build()

	registry := colors.NewColorRegistryBuilder().
		SetDefaultNativeColorsForTesting().
		AddDisambiguationEdge(subIteratorA, iterator).
		AddDisambiguationEdge(subIteratorB, iterator).
		AddDisambiguationEdge(childClass, baseClass).
		Build()
	finder := polyfillUsageFinderTest_createFinderWithRegistry(registry)

	// Siblings under Iterator are both not disjoint from Iterator
	AssertThat(finder.IsDisjointFromIterator(subIteratorA)).IsFalse()
	AssertThat(finder.IsDisjointFromIterator(subIteratorB)).IsFalse()

	// Separate class hierarchy is disjoint from Iterator
	AssertThat(finder.IsDisjointFromIterator(childClass)).IsTrue()
	AssertThat(finder.IsDisjointFromIterator(baseClass)).IsTrue()
}

// Java: PolyfillUsageFinderTest.isDisjointFromIterator_cyclesInSupertypeGraph_terminatesGracefully()
func (c *PolyfillUsageFinderTest) TestIsDisjointFromIterator_cyclesInSupertypeGraph_terminatesGracefully() {
	iterator := colors.Color_SingleBuilder().SetId(colors.StandardColors_ITERATOR_ID).Build()
	classA := colors.Color_SingleBuilder().SetId(colors.ColorId_FromAscii("ClassA")).Build()
	classB := colors.Color_SingleBuilder().SetId(colors.ColorId_FromAscii("ClassB")).Build()

	// Cycle A <-> B with no connection to Iterator
	registryWithCycle := colors.NewColorRegistryBuilder().
		SetDefaultNativeColorsForTesting().
		AddDisambiguationEdge(classA, classB).
		AddDisambiguationEdge(classB, classA).
		Build()
	finderWithCycle := polyfillUsageFinderTest_createFinderWithRegistry(registryWithCycle)

	AssertThat(finderWithCycle.IsDisjointFromIterator(classA)).IsTrue()
	AssertThat(finderWithCycle.IsDisjointFromIterator(classB)).IsTrue()

	// Cycle A <-> B where B also extends Iterator
	registryWithCycleAndIterator := colors.NewColorRegistryBuilder().
		SetDefaultNativeColorsForTesting().
		AddDisambiguationEdge(classA, classB).
		AddDisambiguationEdge(classB, classA).
		AddDisambiguationEdge(classB, iterator).
		Build()
	finderWithCycleAndIterator := polyfillUsageFinderTest_createFinderWithRegistry(registryWithCycleAndIterator)

	AssertThat(finderWithCycleAndIterator.IsDisjointFromIterator(classA)).IsFalse()
	AssertThat(finderWithCycleAndIterator.IsDisjointFromIterator(classB)).IsFalse()
}

// Java: PolyfillUsageFinderTest.isDisjointFromIterator_unionSupertypesInHierarchy()
func (c *PolyfillUsageFinderTest) TestIsDisjointFromIterator_unionSupertypesInHierarchy() {
	iterator := colors.Color_SingleBuilder().SetId(colors.StandardColors_ITERATOR_ID).Build()
	classFoo := colors.Color_SingleBuilder().SetId(colors.ColorId_FromAscii("Foo")).Build()
	classBar := colors.Color_SingleBuilder().SetId(colors.ColorId_FromAscii("Bar")).Build()
	classSub1 := colors.Color_SingleBuilder().SetId(colors.ColorId_FromAscii("Sub1")).Build()
	classSub2 := colors.Color_SingleBuilder().SetId(colors.ColorId_FromAscii("Sub2")).Build()

	unionWithIterator := colors.Color_CreateUnion2(classFoo, iterator)
	unionWithoutIterator := colors.Color_CreateUnion2(classFoo, classBar)

	registry := colors.NewColorRegistryBuilder().
		SetDefaultNativeColorsForTesting().
		AddDisambiguationEdge(classSub1, unionWithIterator).
		AddDisambiguationEdge(classSub2, unionWithoutIterator).
		Build()
	finder := polyfillUsageFinderTest_createFinderWithRegistry(registry)

	// Sub1 has a union supertype that includes Iterator
	AssertThat(finder.IsDisjointFromIterator(classSub1)).IsFalse()

	// Sub2 has a union supertype that only includes disjoint types
	AssertThat(finder.IsDisjointFromIterator(classSub2)).IsTrue()
}

// Java: PolyfillUsageFinderTest.isDisjointFromIterator_primitives_returnsTrue()
func (c *PolyfillUsageFinderTest) TestIsDisjointFromIterator_primitives_returnsTrue() {
	registry := colors.NewColorRegistryBuilder().SetDefaultNativeColorsForTesting().Build()
	finder := polyfillUsageFinderTest_createFinderWithRegistry(registry)

	// Primitives are disjoint from Iterator
	AssertThat(finder.IsDisjointFromIterator(colors.StandardColors_STRING)).IsTrue()
	AssertThat(finder.IsDisjointFromIterator(colors.StandardColors_NUMBER)).IsTrue()
	AssertThat(finder.IsDisjointFromIterator(colors.StandardColors_BOOLEAN)).IsTrue()
	AssertThat(finder.IsDisjointFromIterator(colors.StandardColors_BIGINT)).IsTrue()
	AssertThat(finder.IsDisjointFromIterator(colors.StandardColors_SYMBOL)).IsTrue()
	AssertThat(finder.IsDisjointFromIterator(colors.StandardColors_NULL_OR_VOID)).IsTrue()
}

// Java: PolyfillUsageFinderTest.isDisjointFromIterator_nullAndUndefinedUnions()
func (c *PolyfillUsageFinderTest) TestIsDisjointFromIterator_nullAndUndefinedUnions() {
	registry := colors.NewColorRegistryBuilder().SetDefaultNativeColorsForTesting().Build()
	finder := polyfillUsageFinderTest_createFinderWithRegistry(registry)
	iterator := registry.Get(colors.StandardColors_ITERATOR_ID)
	array := registry.Get(colors.StandardColors_ARRAY_ID)

	nullableArray := colors.Color_CreateUnion2(array, colors.StandardColors_NULL_OR_VOID)
	nullableIterator := colors.Color_CreateUnion2(iterator, colors.StandardColors_NULL_OR_VOID)

	// Nullable/undef Array is still disjoint from Iterator after subtraction
	AssertThat(finder.IsDisjointFromIterator(nullableArray)).IsTrue()

	// Nullable/undef Iterator is not disjoint from Iterator
	AssertThat(finder.IsDisjointFromIterator(nullableIterator)).IsFalse()
}

// Java: PolyfillUsageFinderTest.isDisjointFromIterator_unionReceivers()
func (c *PolyfillUsageFinderTest) TestIsDisjointFromIterator_unionReceivers() {
	iterator := colors.Color_SingleBuilder().SetId(colors.StandardColors_ITERATOR_ID).Build()
	subIterator := colors.Color_SingleBuilder().SetId(colors.ColorId_FromAscii("SubIter")).Build()
	classFoo := colors.Color_SingleBuilder().SetId(colors.ColorId_FromAscii("Foo")).Build()
	classBar := colors.Color_SingleBuilder().SetId(colors.ColorId_FromAscii("Bar")).Build()

	registry := colors.NewColorRegistryBuilder().
		SetDefaultNativeColorsForTesting().
		AddDisambiguationEdge(subIterator, iterator).
		Build()
	finder := polyfillUsageFinderTest_createFinderWithRegistry(registry)

	// All elements disjoint -> disjoint
	unionAllDisjoint := colors.Color_CreateUnion2(classFoo, classBar)
	AssertThat(finder.IsDisjointFromIterator(unionAllDisjoint)).IsTrue()

	// One element is Iterator -> not disjoint
	unionWithIterator := colors.Color_CreateUnion2(classFoo, iterator)
	AssertThat(finder.IsDisjointFromIterator(unionWithIterator)).IsFalse()

	// One element is SubIterator -> not disjoint
	unionWithSubIterator := colors.Color_CreateUnion2(classFoo, subIterator)
	AssertThat(finder.IsDisjointFromIterator(unionWithSubIterator)).IsFalse()

	// One element is UNKNOWN -> not disjoint
	unionWithUnknown := colors.Color_CreateUnion2(classFoo, colors.StandardColors_UNKNOWN)
	AssertThat(finder.IsDisjointFromIterator(unionWithUnknown)).IsFalse()

	// One element is TOP_OBJECT -> not disjoint
	unionWithTopObject := colors.Color_CreateUnion2(classFoo, colors.StandardColors_TOP_OBJECT)
	AssertThat(finder.IsDisjointFromIterator(unionWithTopObject)).IsFalse()
}

// Java: PolyfillUsageFinderTest.isDisjointFromIterator_invalidatingColors_structuralRecordsVsNativeArrays()
func (c *PolyfillUsageFinderTest) TestIsDisjointFromIterator_invalidatingColors_structuralRecordsVsNativeArrays() {
	registry := colors.NewColorRegistryBuilder().SetDefaultNativeColorsForTesting().Build()
	finder := polyfillUsageFinderTest_createFinderWithRegistry(registry)

	structuralRecord := colors.Color_SingleBuilder().SetId(colors.ColorId_FromAscii("Struct")).SetInvalidating(true).Build()
	normalClass := colors.Color_SingleBuilder().SetId(colors.ColorId_FromAscii("Normal")).SetInvalidating(false).Build()

	invalidatingArray := colors.Color_SingleBuilder().SetId(colors.StandardColors_ARRAY_ID).SetInvalidating(true).Build()
	invalidatingReadonlyArray := colors.Color_SingleBuilder().SetId(colors.StandardColors_READONLY_ARRAY_ID).SetInvalidating(true).Build()
	invalidatingTemplateArray := colors.Color_SingleBuilder().SetId(colors.StandardColors_I_TEMPLATE_ARRAY_ID).SetInvalidating(true).Build()

	invalidatingMap := colors.Color_SingleBuilder().SetId(colors.StandardColors_MAP_ID).SetInvalidating(true).Build()
	invalidatingMapPrototype := colors.Color_SingleBuilder().SetId(colors.StandardColors_MAP_PROTOTYPE_ID).SetInvalidating(true).Build()
	invalidatingReadonlyMap := colors.Color_SingleBuilder().SetId(colors.StandardColors_READONLY_MAP_ID).SetInvalidating(true).Build()
	invalidatingSet := colors.Color_SingleBuilder().SetId(colors.StandardColors_SET_ID).SetInvalidating(true).Build()
	invalidatingSetPrototype := colors.Color_SingleBuilder().SetId(colors.StandardColors_SET_PROTOTYPE_ID).SetInvalidating(true).Build()
	invalidatingReadonlySet := colors.Color_SingleBuilder().SetId(colors.StandardColors_READONLY_SET_ID).SetInvalidating(true).Build()

	// Structural interface / record: isInvalidating causes disjoint to return false (b/571244650)
	AssertThat(finder.IsDisjointFromIterator(structuralRecord)).IsFalse()

	// Normal non-invalidating class: is disjoint
	AssertThat(finder.IsDisjointFromIterator(normalClass)).IsTrue()

	// Native Array/Map/Set types marked invalidating are exempted
	AssertThat(finder.IsDisjointFromIterator(invalidatingArray)).IsTrue()
	AssertThat(finder.IsDisjointFromIterator(invalidatingReadonlyArray)).IsTrue()
	AssertThat(finder.IsDisjointFromIterator(invalidatingTemplateArray)).IsTrue()
	AssertThat(finder.IsDisjointFromIterator(invalidatingMap)).IsTrue()
	AssertThat(finder.IsDisjointFromIterator(invalidatingMapPrototype)).IsTrue()
	AssertThat(finder.IsDisjointFromIterator(invalidatingReadonlyMap)).IsTrue()
	AssertThat(finder.IsDisjointFromIterator(invalidatingSet)).IsTrue()
	AssertThat(finder.IsDisjointFromIterator(invalidatingSetPrototype)).IsTrue()
	AssertThat(finder.IsDisjointFromIterator(invalidatingReadonlySet)).IsTrue()
}

// Java: PolyfillUsageFinderTest.isDisjointFromIterator_nodeAndOwnerChecks()
func (c *PolyfillUsageFinderTest) TestIsDisjointFromIterator_nodeAndOwnerChecks() {
	registry := colors.NewColorRegistryBuilder().SetDefaultNativeColorsForTesting().Build()
	array := registry.Get(colors.StandardColors_ARRAY_ID)
	iterator := registry.Get(colors.StandardColors_ITERATOR_ID)

	finderWithRegistry := polyfillUsageFinderTest_createFinderWithRegistry(registry)

	compilerWithoutRegistry := NewCompiler()
	finderWithoutRegistry := NewPolyfillUsageFinder(
		compilerWithoutRegistry,
		PolyfillUsageFinder_Polyfills_FromTable("Iterator.prototype.toArray es_next es3 es6/iterator/toarray"),
		/* disambiguateProperties= */ true)

	receiverNode := rhino.IR_Name("x")
	getPropNode := rhino.IR_Getprop(receiverNode, "toArray")

	// Receiver has Array color, owner is Iterator -> disjoint
	receiverNode.SetColor(array)
	AssertThat(finderWithRegistry.IsDisjointFromIteratorNode(getPropNode, "Iterator")).IsTrue()

	// Receiver has Iterator color, owner is Iterator -> not disjoint
	receiverNode.SetColor(iterator)
	AssertThat(finderWithRegistry.IsDisjointFromIteratorNode(getPropNode, "Iterator")).IsFalse()

	// Receiver has null color -> not disjoint
	receiverNode.SetColor(nil)
	AssertThat(finderWithRegistry.IsDisjointFromIteratorNode(getPropNode, "Iterator")).IsFalse()

	// Owner is null ("") or not Iterator -> not disjoint
	receiverNode.SetColor(array)
	AssertThat(finderWithRegistry.IsDisjointFromIteratorNode(getPropNode, "")).IsFalse()
	AssertThat(finderWithRegistry.IsDisjointFromIteratorNode(getPropNode, "Array")).IsFalse()
	AssertThat(finderWithRegistry.IsDisjointFromIteratorNode(getPropNode, "UnknownType")).IsFalse()

	// Compiler without registry -> always not disjoint
	receiverNode.SetColor(array)
	AssertThat(finderWithoutRegistry.IsDisjointFromIteratorNode(getPropNode, "Iterator")).IsFalse()
}

// Java: PolyfillUsageFinderTest.isDisjointFromIterator_receiverWithNullColor_returnsFalse()
func (c *PolyfillUsageFinderTest) TestIsDisjointFromIterator_receiverWithNullColor_returnsFalse() {
	registry := colors.NewColorRegistryBuilder().SetDefaultNativeColorsForTesting().Build()
	finder := polyfillUsageFinderTest_createFinderWithRegistry(registry)

	receiverNode := rhino.IR_Name("x")
	AssertThat(receiverNode.GetColor()).IsNull()
	getPropNode := rhino.IR_Getprop(receiverNode, "toArray")

	AssertThat(finder.IsDisjointFromIteratorNode(getPropNode, "Iterator")).IsFalse()
}

// Java: PolyfillUsageFinderTest.traverse_withColorRegistry_filtersDisjointUsages()
func (c *PolyfillUsageFinderTest) TestTraverse_withColorRegistry_filtersDisjointUsages() {
	registry := colors.NewColorRegistryBuilder().SetDefaultNativeColorsForTesting().Build()
	array := registry.Get(colors.StandardColors_ARRAY_ID)
	iterator := registry.Get(colors.StandardColors_ITERATOR_ID)

	finder := PolyfillUsageFinderTest_TestPolyfillUsageFinder_Builder_().
		WithPolyfillTableLines(
			"Iterator.prototype.toArray es_next es3 es6/iterator/toarray\n" +
				"Array.prototype.fill es6 es3 es6/array/fill\n").
		WithDisambiguateProperties(true).
		ForSourceLines(
			"arr.toArray(); // 1: receiver is Array, disjoint from Iterator -> filtered!\n" +
				"iter.toArray(); // 2: receiver is Iterator -> kept!\n" +
				"untyped.toArray(); // 3: untyped -> kept!\n" +
				"arr.fill(0); // 4: owner is Array (not Iterator) -> kept!\n").
		WithColorRegistry(registry).
		Build()

	// Annotate AST nodes with colors
	rootNode := finder.GetRootNode()
	arrNodes := polyfillUsageFinderTest_findAllNodes(rootNode, func(n rhino.Node) bool {
		return n.IsName() && n.GetString() == "arr"
	})
	AssertThat(arrNodes).HasSize(2)
	for _, arrNode := range arrNodes {
		arrNode.SetColor(array)
	}
	rhino.CheckNotNull(polyfillUsageFinderTest_findFirstNode(rootNode, func(n rhino.Node) bool {
		return n.IsName() && n.GetString() == "iter"
	})).SetColor(iterator)

	usages := finder.GetUnguardedUsages()

	// Should include: iter.toArray (line 2), untyped.toArray (line 3), arr.fill (line 4)
	// Should NOT include: arr.toArray (line 1)
	AssertThat(usages).HasSize(3)

	PolyfillUsageFinderTest_PolyfillUsageSubject_AssertPolyfillUsage(usages[0]).
		HasName("toArray").
		HasNodeThat().
		HasLineno(2)

	PolyfillUsageFinderTest_PolyfillUsageSubject_AssertPolyfillUsage(usages[1]).
		HasName("toArray").
		HasNodeThat().
		HasLineno(3)

	PolyfillUsageFinderTest_PolyfillUsageSubject_AssertPolyfillUsage(usages[2]).
		HasName("fill").
		HasNodeThat().
		HasLineno(4)
}

// Java: PolyfillUsageFinderTest.traverse_withColorRegistry_disambiguatePropertiesDisabled_doesNotFilterUsages()
func (c *PolyfillUsageFinderTest) TestTraverse_withColorRegistry_disambiguatePropertiesDisabled_doesNotFilterUsages() {
	registry := colors.NewColorRegistryBuilder().SetDefaultNativeColorsForTesting().Build()
	array := registry.Get(colors.StandardColors_ARRAY_ID)

	finder := PolyfillUsageFinderTest_TestPolyfillUsageFinder_Builder_().
		WithPolyfillTableLines("Iterator.prototype.toArray es_next es3 es6/iterator/toarray\n").
		WithDisambiguateProperties(false).
		ForSourceLines("arr.toArray();\n").
		WithColorRegistry(registry).
		Build()

	rootNode := finder.GetRootNode()
	rhino.CheckNotNull(polyfillUsageFinderTest_findFirstNode(rootNode, func(n rhino.Node) bool {
		return n.IsName() && n.GetString() == "arr"
	})).SetColor(array)

	usages := finder.GetUnguardedUsages()
	AssertThat(usages).HasSize(1)
	PolyfillUsageFinderTest_PolyfillUsageSubject_AssertPolyfillUsage(usages[0]).HasName("toArray")
}

// Java: PolyfillUsageFinderTest.traverse_withoutColorRegistry_doesNotFilterUsages()
func (c *PolyfillUsageFinderTest) TestTraverse_withoutColorRegistry_doesNotFilterUsages() {
	registry := colors.NewColorRegistryBuilder().SetDefaultNativeColorsForTesting().Build()
	array := registry.Get(colors.StandardColors_ARRAY_ID)

	finder := PolyfillUsageFinderTest_TestPolyfillUsageFinder_Builder_().
		WithPolyfillTableLines("Iterator.prototype.toArray es_next es3 es6/iterator/toarray\n").
		ForSourceLines("arr.toArray();\n").
		Build() // No color registry

	// Even if arr has a color, without ColorRegistry on compiler it should not filter
	rootNode := finder.GetRootNode()
	rhino.CheckNotNull(polyfillUsageFinderTest_findFirstNode(rootNode, func(n rhino.Node) bool {
		return n.IsName() && n.GetString() == "arr"
	})).SetColor(array)

	usages := finder.GetUnguardedUsages()
	AssertThat(usages).HasSize(1)
	PolyfillUsageFinderTest_PolyfillUsageSubject_AssertPolyfillUsage(usages[0]).HasName("toArray")
}

func polyfillUsageFinderTest_findFirstNode(root rhino.Node, predicate func(rhino.Node) bool) rhino.Node {
	if predicate(root) {
		return root
	}
	for child := root.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		found := polyfillUsageFinderTest_findFirstNode(child, predicate)
		if !found.IsNil() {
			return found
		}
	}
	return rhino.NilNode
}

func polyfillUsageFinderTest_findAllNodes(root rhino.Node, predicate func(rhino.Node) bool) []rhino.Node {
	var out []rhino.Node
	polyfillUsageFinderTest_findAllNodesInternal(root, predicate, &out)
	return out
}

func polyfillUsageFinderTest_findAllNodesInternal(root rhino.Node, predicate func(rhino.Node) bool, out *[]rhino.Node) {
	if predicate(root) {
		*out = append(*out, root)
	}
	for child := root.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		polyfillUsageFinderTest_findAllNodesInternal(child, predicate, out)
	}
}

// PolyfillUsageFinderTest_TestPolyfillUsageFinder uses PolyfillUsageFinder
// to find polyfill usages in source code for testing.
type PolyfillUsageFinderTest_TestPolyfillUsageFinder struct {
	polyfillUsageFinder *PolyfillUsageFinder
	rootNode            rhino.Node
}

// Java: PolyfillUsageFinderTest.TestPolyfillUsageFinder.builder()
func PolyfillUsageFinderTest_TestPolyfillUsageFinder_Builder_() *PolyfillUsageFinderTest_TestPolyfillUsageFinder_Builder {
	return &PolyfillUsageFinderTest_TestPolyfillUsageFinder_Builder{}
}

// PolyfillUsageFinderTest_TestPolyfillUsageFinder_Builder is
// TestPolyfillUsageFinder.Builder.
type PolyfillUsageFinderTest_TestPolyfillUsageFinder_Builder struct {
	polyfillsTable         *PolyfillUsageFinder_Polyfills
	compiler               *Compiler
	polyfillUsageFinder    *PolyfillUsageFinder
	rootNode               rhino.Node
	colorRegistry          *colors.ColorRegistry
	disambiguateProperties bool
}

// Java: PolyfillUsageFinderTest.TestPolyfillUsageFinder.Builder.withPolyfillTableLines(String)
func (b *PolyfillUsageFinderTest_TestPolyfillUsageFinder_Builder) WithPolyfillTableLines(polyfillTable string) *PolyfillUsageFinderTest_TestPolyfillUsageFinder_Builder {
	b.polyfillsTable = PolyfillUsageFinder_Polyfills_FromTable(polyfillTable)
	return b
}

// Java: PolyfillUsageFinderTest.TestPolyfillUsageFinder.Builder.withColorRegistry(ColorRegistry)
func (b *PolyfillUsageFinderTest_TestPolyfillUsageFinder_Builder) WithColorRegistry(colorRegistry *colors.ColorRegistry) *PolyfillUsageFinderTest_TestPolyfillUsageFinder_Builder {
	b.colorRegistry = colorRegistry
	return b
}

// Java: PolyfillUsageFinderTest.TestPolyfillUsageFinder.Builder.withDisambiguateProperties(boolean)
func (b *PolyfillUsageFinderTest_TestPolyfillUsageFinder_Builder) WithDisambiguateProperties(disambiguateProperties bool) *PolyfillUsageFinderTest_TestPolyfillUsageFinder_Builder {
	b.disambiguateProperties = disambiguateProperties
	return b
}

// Java: PolyfillUsageFinderTest.TestPolyfillUsageFinder.Builder.forSourceLines(String)
func (b *PolyfillUsageFinderTest_TestPolyfillUsageFinder_Builder) ForSourceLines(src string) *PolyfillUsageFinderTest_TestPolyfillUsageFinder_Builder {
	srcFile := SourceFile_Builder_().WithPath("src.js").WithContent(src).Build()
	options := NewCompilerOptions()
	// Don't include `"use strict";` when printing the AST as source text
	options.SetEmitUseStrict(false)
	b.compiler = NewCompiler()
	b.compiler.Init(
		[]*SourceFile{NewTestExternsBuilder().BuildExternsFile("externs.js")},
		[]*SourceFile{srcFile},
		options)
	b.compiler.Parse()
	b.rootNode = b.compiler.GetJsRoot()
	return b
}

// Java: PolyfillUsageFinderTest.TestPolyfillUsageFinder.Builder.build()
func (b *PolyfillUsageFinderTest_TestPolyfillUsageFinder_Builder) Build() *PolyfillUsageFinderTest_TestPolyfillUsageFinder {
	jutil.CheckNotNull(b.compiler)
	if b.colorRegistry != nil {
		// Go-only: test-helper — equivalent to TestCompiler.setTestColorRegistry in Java.
		b.compiler.colorRegistry = b.colorRegistry
	}
	b.polyfillUsageFinder =
		NewPolyfillUsageFinder(b.compiler, jutil.CheckNotNull(b.polyfillsTable), b.disambiguateProperties)
	return newPolyfillUsageFinderTest_TestPolyfillUsageFinder(b)
}

// Java: PolyfillUsageFinderTest.TestPolyfillUsageFinder(Builder)
func newPolyfillUsageFinderTest_TestPolyfillUsageFinder(builder *PolyfillUsageFinderTest_TestPolyfillUsageFinder_Builder) *PolyfillUsageFinderTest_TestPolyfillUsageFinder {
	f := &PolyfillUsageFinderTest_TestPolyfillUsageFinder{}
	f.polyfillUsageFinder = builder.polyfillUsageFinder
	f.rootNode = builder.rootNode
	return f
}

// Java: PolyfillUsageFinderTest.TestPolyfillUsageFinder.getRootNode()
func (f *PolyfillUsageFinderTest_TestPolyfillUsageFinder) GetRootNode() rhino.Node {
	return f.rootNode
}

// Java: PolyfillUsageFinderTest.TestPolyfillUsageFinder.getUnguardedUsages()
func (f *PolyfillUsageFinderTest_TestPolyfillUsageFinder) GetUnguardedUsages() []*PolyfillUsageFinder_PolyfillUsage {
	consumer := newPolyfillUsageFinderTest_PolyfillUsageCollectingConsumer()
	f.polyfillUsageFinder.TraverseExcludingGuarded(f.rootNode, consumer.Accept)
	return consumer.GetUsages()
}

// Java: PolyfillUsageFinderTest.TestPolyfillUsageFinder.getGuardedUsages()
func (f *PolyfillUsageFinderTest_TestPolyfillUsageFinder) GetGuardedUsages() []*PolyfillUsageFinder_PolyfillUsage {
	consumer := newPolyfillUsageFinderTest_PolyfillUsageCollectingConsumer()
	f.polyfillUsageFinder.TraverseOnlyGuarded(f.rootNode, consumer.Accept)
	return consumer.GetUsages()
}

// Java: PolyfillUsageFinderTest.TestPolyfillUsageFinder.getAllUsages()
func (f *PolyfillUsageFinderTest_TestPolyfillUsageFinder) GetAllUsages() []*PolyfillUsageFinder_PolyfillUsage {
	consumer := newPolyfillUsageFinderTest_PolyfillUsageCollectingConsumer()
	f.polyfillUsageFinder.TraverseIncludingGuarded(f.rootNode, consumer.Accept)
	return consumer.GetUsages()
}

// polyfillUsageFinderTest_PolyfillUsageCollectingConsumer consumes
// PolyfillUsage objects by storing them into a retrievable list.
type polyfillUsageFinderTest_PolyfillUsageCollectingConsumer struct {
	polyfillUsageListBuilder []*PolyfillUsageFinder_PolyfillUsage
}

func newPolyfillUsageFinderTest_PolyfillUsageCollectingConsumer() *polyfillUsageFinderTest_PolyfillUsageCollectingConsumer {
	return &polyfillUsageFinderTest_PolyfillUsageCollectingConsumer{}
}

// Java: PolyfillUsageFinderTest.PolyfillUsageCollectingConsumer.accept(PolyfillUsage)
func (c *polyfillUsageFinderTest_PolyfillUsageCollectingConsumer) Accept(polyfillUsage *PolyfillUsageFinder_PolyfillUsage) {
	c.polyfillUsageListBuilder = append(c.polyfillUsageListBuilder, jutil.CheckNotNull(polyfillUsage))
}

// Java: PolyfillUsageFinderTest.PolyfillUsageCollectingConsumer.getUsages()
func (c *polyfillUsageFinderTest_PolyfillUsageCollectingConsumer) GetUsages() []*PolyfillUsageFinder_PolyfillUsage {
	return append([]*PolyfillUsageFinder_PolyfillUsage(nil), c.polyfillUsageListBuilder...)
}

// PolyfillUsageFinderTest_PolyfillSubject is the Truth Subject for Polyfill.
type PolyfillUsageFinderTest_PolyfillSubject struct {
	*Subject
	actual *PolyfillUsageFinder_Polyfill
}

// Java: PolyfillUsageFinderTest.PolyfillSubject(FailureMetadata, Polyfill)
func newPolyfillUsageFinderTest_PolyfillSubject(s *Subject, polyfill *PolyfillUsageFinder_Polyfill) *PolyfillUsageFinderTest_PolyfillSubject {
	return &PolyfillUsageFinderTest_PolyfillSubject{Subject: s, actual: polyfill}
}

// Java: PolyfillUsageFinderTest.PolyfillSubject.assertPolyfill(Polyfill)
func PolyfillUsageFinderTest_PolyfillSubject_AssertPolyfill(polyfill *PolyfillUsageFinder_Polyfill) *PolyfillUsageFinderTest_PolyfillSubject {
	return newPolyfillUsageFinderTest_PolyfillSubject(testjunit.NewSubject(polyfill), polyfill)
}

// check is Subject.check(name): a subject for a property of the polyfill.
func (s *PolyfillUsageFinderTest_PolyfillSubject) check(what string) *StandardSubjectBuilder {
	msg := "value of: polyfill." + what + "\npolyfill was: " + show(s.actual)
	if s.Message() != "" {
		msg = s.Message() + "\n" + msg
	}
	return testjunit.NewStandardSubjectBuilder(msg)
}

// Java: PolyfillUsageFinderTest.PolyfillSubject.hasNativeSymbol(String)
func (s *PolyfillUsageFinderTest_PolyfillSubject) HasNativeSymbol(expected string) *PolyfillUsageFinderTest_PolyfillSubject {
	s.check("nativeSymbol").That(s.actual.nativeSymbol).IsEqualTo(expected)
	return s
}

// Java: PolyfillUsageFinderTest.PolyfillSubject.hasNativeVersionStr(String)
func (s *PolyfillUsageFinderTest_PolyfillSubject) HasNativeVersionStr(expected string) *PolyfillUsageFinderTest_PolyfillSubject {
	s.check("nativeVersion").That(s.actual.nativeVersion).IsEqualTo(expected)
	return s
}

// Java: PolyfillUsageFinderTest.PolyfillSubject.hasPolyfillVersionStr(String)
func (s *PolyfillUsageFinderTest_PolyfillSubject) HasPolyfillVersionStr(expected string) *PolyfillUsageFinderTest_PolyfillSubject {
	s.check("polyfillVersion").That(s.actual.polyfillVersion).IsEqualTo(expected)
	return s
}

// Java: PolyfillUsageFinderTest.PolyfillSubject.hasLibrary(String)
func (s *PolyfillUsageFinderTest_PolyfillSubject) HasLibrary(expectedLibraryPath string) *PolyfillUsageFinderTest_PolyfillSubject {
	s.check("library").That(s.actual.library).IsEqualTo(expectedLibraryPath)
	return s
}

// Java: PolyfillUsageFinderTest.PolyfillSubject.hasKind(Kind)
func (s *PolyfillUsageFinderTest_PolyfillSubject) HasKind(expectedKind PolyfillUsageFinder_Polyfill_Kind) *PolyfillUsageFinderTest_PolyfillSubject {
	s.check("kind").That(s.actual.kind).IsEqualTo(expectedKind)
	return s
}

// PolyfillUsageFinderTest_PolyfillUsageSubject is the Truth Subject for
// PolyfillUsage.
type PolyfillUsageFinderTest_PolyfillUsageSubject struct {
	*Subject
	actual *PolyfillUsageFinder_PolyfillUsage
}

// Java: PolyfillUsageFinderTest.PolyfillUsageSubject(FailureMetadata, PolyfillUsage)
func newPolyfillUsageFinderTest_PolyfillUsageSubject(s *Subject, polyfillUsage *PolyfillUsageFinder_PolyfillUsage) *PolyfillUsageFinderTest_PolyfillUsageSubject {
	return &PolyfillUsageFinderTest_PolyfillUsageSubject{Subject: s, actual: polyfillUsage}
}

// Java: PolyfillUsageFinderTest.PolyfillUsageSubject.assertPolyfillUsage(PolyfillUsage)
func PolyfillUsageFinderTest_PolyfillUsageSubject_AssertPolyfillUsage(polyfillUsage *PolyfillUsageFinder_PolyfillUsage) *PolyfillUsageFinderTest_PolyfillUsageSubject {
	return newPolyfillUsageFinderTest_PolyfillUsageSubject(testjunit.NewSubject(polyfillUsage), polyfillUsage)
}

// check is Subject.check(name): a subject for a property of the usage.
func (s *PolyfillUsageFinderTest_PolyfillUsageSubject) check(what string) *StandardSubjectBuilder {
	msg := "value of: polyfillUsage." + what + "\npolyfillUsage was: " + show(s.actual)
	if s.Message() != "" {
		msg = s.Message() + "\n" + msg
	}
	return testjunit.NewStandardSubjectBuilder(msg)
}

// Java: PolyfillUsageFinderTest.PolyfillUsageSubject.hasName(String)
func (s *PolyfillUsageFinderTest_PolyfillUsageSubject) HasName(expectedName string) *PolyfillUsageFinderTest_PolyfillUsageSubject {
	s.check("name()").That(s.actual.Name()).IsEqualTo(expectedName)
	return s
}

// Java: PolyfillUsageFinderTest.PolyfillUsageSubject.isExplicitGlobal()
func (s *PolyfillUsageFinderTest_PolyfillUsageSubject) IsExplicitGlobal() *PolyfillUsageFinderTest_PolyfillUsageSubject {
	s.check("isExplicitGlobal()").That(s.actual.IsExplicitGlobal()).IsTrue()
	return s
}

// Java: PolyfillUsageFinderTest.PolyfillUsageSubject.isNotExplicitGlobal()
func (s *PolyfillUsageFinderTest_PolyfillUsageSubject) IsNotExplicitGlobal() *PolyfillUsageFinderTest_PolyfillUsageSubject {
	s.check("isExplicitGlobal()").That(s.actual.IsExplicitGlobal()).IsFalse()
	return s
}

// Java: PolyfillUsageFinderTest.PolyfillUsageSubject.hasNodeThat()
func (s *PolyfillUsageFinderTest_PolyfillUsageSubject) HasNodeThat() *NodeSubject {
	return NodeSubject_AssertNode(s.actual.Node())
}

// Java: PolyfillUsageFinderTest.PolyfillUsageSubject.hasPolyfillThat()
func (s *PolyfillUsageFinderTest_PolyfillUsageSubject) HasPolyfillThat() *PolyfillUsageFinderTest_PolyfillSubject {
	return PolyfillUsageFinderTest_PolyfillSubject_AssertPolyfill(s.actual.Polyfill())
}
