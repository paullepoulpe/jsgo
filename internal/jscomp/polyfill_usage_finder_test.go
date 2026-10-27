// Port of com/google/javascript/jscomp/PolyfillUsageFinderTest.java @ closure-compiler v20261021
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
	polyfillsTable      *PolyfillUsageFinder_Polyfills
	compiler            *Compiler
	polyfillUsageFinder *PolyfillUsageFinder
	rootNode            rhino.Node
}

// Java: PolyfillUsageFinderTest.TestPolyfillUsageFinder.Builder.withPolyfillTableLines(String)
func (b *PolyfillUsageFinderTest_TestPolyfillUsageFinder_Builder) WithPolyfillTableLines(polyfillTable string) *PolyfillUsageFinderTest_TestPolyfillUsageFinder_Builder {
	b.polyfillsTable = PolyfillUsageFinder_Polyfills_FromTable(polyfillTable)
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
	b.polyfillUsageFinder =
		NewPolyfillUsageFinder(jutil.CheckNotNull(b.compiler), jutil.CheckNotNull(b.polyfillsTable))
	return newPolyfillUsageFinderTest_TestPolyfillUsageFinder(b)
}

// Java: PolyfillUsageFinderTest.TestPolyfillUsageFinder(Builder)
func newPolyfillUsageFinderTest_TestPolyfillUsageFinder(builder *PolyfillUsageFinderTest_TestPolyfillUsageFinder_Builder) *PolyfillUsageFinderTest_TestPolyfillUsageFinder {
	f := &PolyfillUsageFinderTest_TestPolyfillUsageFinder{}
	f.polyfillUsageFinder = builder.polyfillUsageFinder
	f.rootNode = builder.rootNode
	return f
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
