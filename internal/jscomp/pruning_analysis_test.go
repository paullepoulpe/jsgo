// Port of com/google/javascript/jscomp/PruningAnalysisTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/PruningAnalysisTest.java)

/*
 * Copyright 2026 The Closure Compiler Authors.
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
)

// PruningAnalysisTest holds the unit tests of PruningAnalysis.
type PruningAnalysisTest struct {
	t        *testing.T
	compiler *Compiler
}

func (pt *PruningAnalysisTest) JUnitSetT(t *testing.T) { pt.t = t }
func (pt *PruningAnalysisTest) JUnitAfter()            {}

// Java: PruningAnalysisTest.setUp()
func (pt *PruningAnalysisTest) JUnitBefore() {
	pt.compiler = NewCompiler()
	pt.compiler.InitOptions(NewCompilerOptions())
}

func TestPruningAnalysisTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return &PruningAnalysisTest{} })
}

// Java: PruningAnalysisTest.testNoBloat()
func (pt *PruningAnalysisTest) TestNoBloat() {
	root1 := pt.createInput("root1.js", "goog.provide('root1'); goog.require('dep1');")
	dep1 := pt.createInput("dep1.js", "goog.provide('dep1');")
	root2 := pt.createInput("root2.js", "goog.provide('root2'); goog.require('dep2');")
	dep2 := pt.createInput("dep2.js", "goog.provide('dep2');")

	result := pt.analyze([]*CompilerInput{root1, root2}, root1, dep1, root2, dep2)

	AssertThat(result.EntryPointDependencyCount()).ContainsExactlyEntriesIn(pruningAnalysisTest_mapOf("root1.js", 2, "root2.js", 2))
}

// Java: PruningAnalysisTest.testSharedDependencies()
func (pt *PruningAnalysisTest) TestSharedDependencies() {
	root1 := pt.createInput("root1.js", "goog.provide('root1'); goog.require('s1'); goog.require('s2');")
	root2 := pt.createInput("root2.js", "goog.provide('root2'); goog.require('s1'); goog.require('s3');")
	root3 := pt.createInput("root3.js", "goog.provide('root3'); goog.require('s1'); goog.require('s2'); goog.require('s3');")
	s1 := pt.createInput("s1.js", "goog.provide('s1');")
	s2 := pt.createInput("s2.js", "goog.provide('s2');")
	s3 := pt.createInput("s3.js", "goog.provide('s3');")

	result := pt.analyze([]*CompilerInput{root1, root2, root3}, root1, root2, root3, s1, s2, s3)

	AssertThat(result.EntryPointDependencyCount()).ContainsAtLeastEntriesIn(pruningAnalysisTest_mapOf("root1.js", 3, "root2.js", 3, "root3.js", 4))
}

// Java: PruningAnalysisTest.testWeakDependencies_included()
func (pt *PruningAnalysisTest) TestWeakDependencies_included() {
	root1 := pt.createInput("root1.js", "goog.provide('root1'); goog.requireType('weak1');")
	weak1 := pt.createInput("weak1.js", "goog.provide('weak1');")

	result := pt.analyze([]*CompilerInput{root1}, root1, weak1)

	AssertThat(result.EntryPointDependencyCount()).ContainsAtLeastEntriesIn(pruningAnalysisTest_mapOf("root1.js", 2))
}

// Java: PruningAnalysisTest.testForwardDeclare_notIncluded()
func (pt *PruningAnalysisTest) TestForwardDeclare_notIncluded() {
	// goog.forwardDeclare is NOT captured by JsFileRegexParser as a dependency.
	root1 := pt.createInput("root1.js", "goog.provide('root1'); goog.forwardDeclare('missing');")
	fwdDeclared := pt.createInput("fwdDeclared.js", "goog.provide('fwdDeclared');")

	result := pt.analyze([]*CompilerInput{root1}, root1, fwdDeclared)

	AssertThat(result.EntryPointDependencyCount()).ContainsAtLeastEntriesIn(pruningAnalysisTest_mapOf("root1.js", 1))
}

// Java: PruningAnalysisTest.testOverlappingEntryPoints()
func (pt *PruningAnalysisTest) TestOverlappingEntryPoints() {
	entry1 := pt.createInput("entry1.js", "goog.require('dep1');")
	entry2 := pt.createInput("entry2.js", "goog.require('dep2');")
	dep1 := pt.createInput("dep1.js", "goog.provide('dep1'); goog.require('dep2'); goog.require('excl1');")
	dep2 := pt.createInput("dep2.js", "goog.provide('dep2'); goog.require('excl2');")
	excl1 := pt.createInput("excl1.js", "goog.provide('excl1');")
	excl2 := pt.createInput("excl2.js", "goog.provide('excl2');")

	result := pt.analyze([]*CompilerInput{entry1, entry2}, entry1, entry2, dep1, dep2, excl1, excl2)

	// Symbols are entry1.js and entry2.js.
	// entry1.js closure: {entry1.js, dep1.js, dep2.js excl1.js, excl2.js}
	// entry2.js closure: {entry2.js, dep2.js, excl2.js}
	AssertThat(result.EntryPointDependencyCount()).ContainsAtLeastEntriesIn(pruningAnalysisTest_mapOf("entry1.js", 5, "entry2.js", 3))
}

// Java: PruningAnalysisTest.testMultipleEntryPoints()
func (pt *PruningAnalysisTest) TestMultipleEntryPoints() {
	root1 := pt.createInput("root1.js", "goog.require('dep1');")
	root2 := pt.createInput("root2.js", "goog.require('dep2');")
	dep1 := pt.createInput("dep1.js", "goog.provide('dep1');")
	dep2 := pt.createInput("dep2.js", "goog.provide('dep2');")

	result := pt.analyze([]*CompilerInput{root1, root2}, root1, root2, dep1, dep2)

	AssertThat(result.EntryPointDependencyCount()).ContainsAtLeastEntriesIn(pruningAnalysisTest_mapOf("root1.js", 2, "root2.js", 2))
}

// Java: PruningAnalysisTest.testWeakCycle()
func (pt *PruningAnalysisTest) TestWeakCycle() {
	dep1 := pt.createInput("dep1.js", "goog.provide('dep1'); goog.require('dep2');")
	dep2 := pt.createInput("dep2.js", "goog.provide('dep2'); goog.requireType('root1');")
	root := pt.createInput("root.js", "goog.require('dep1');")

	result := pt.analyze([]*CompilerInput{root}, root, dep1, dep2)

	AssertThat(result.EntryPointDependencyCount()).ContainsAtLeastEntriesIn(pruningAnalysisTest_mapOf("root.js", 3))
}

// Java: PruningAnalysisTest.testBottlenecks()
func (pt *PruningAnalysisTest) TestBottlenecks() {
	// entry -> A -> B -> {C, D, E}
	a := pt.createInput("a.js", "goog.provide('a'); goog.require('b');")
	b := pt.createInput("b.js", "goog.provide('b'); goog.require('c'); goog.require('d'); goog.require('e');")
	c := pt.createInput("c.js", "goog.provide('c');")
	d := pt.createInput("d.js", "goog.provide('d');")
	e := pt.createInput("e.js", "goog.provide('e');")
	entry := pt.createInput("entry.js", "goog.require('a');")

	result := pt.analyze([]*CompilerInput{entry}, a, b, c, d, e, entry)

	// Bottleneck scores (dominated subtree size):
	// a.js: dominates {a, b, c, d, e}. Size = 5.
	// b.js: dominates {b, c, d, e}. Size = 4.
	// c.js, d.js, e.js: Size = 1.
	AssertThat(result.BottleneckBlame()).ContainsEntry("a.js", 5)
	AssertThat(result.BottleneckBlame()).ContainsEntry("b.js", 4)
	AssertThat(result.BottleneckBlame()).DoesNotContainEntry("c.js", 1)
}

// Java: PruningAnalysisTest.testSharedBottleneckInAnalysis()
func (pt *PruningAnalysisTest) TestSharedBottleneckInAnalysis() {
	// entry1 -> B, entry2 -> B, B -> {C, D}
	b := pt.createInput("b.js", "goog.provide('b'); goog.require('c'); goog.require('d');")
	c := pt.createInput("c.js", "goog.provide('c');")
	d := pt.createInput("d.js", "goog.provide('d');")
	entry1 := pt.createInput("entry1.js", "goog.require('b');")
	entry2 := pt.createInput("entry2.js", "goog.require('b');")

	result := pt.analyze([]*CompilerInput{entry1, entry2}, b, c, d, entry1, entry2)

	// entry1.js and entry2.js are roots.
	// Virtual root -> {entry1.js, entry2.js}
	// entry1.js -> b.js
	// entry2.js -> b.js
	// b.js -> {c.js, d.js}
	// IDom of b.js is Virtual Root.
	// b.js dominates {b.js, c.js, d.js}. Size = 3.
	AssertThat(result.BottleneckBlame()).ContainsEntry("b.js", 3)
}

// Java: PruningAnalysisTest.createInput(String, String)
func (pt *PruningAnalysisTest) createInput(name string, code string) *CompilerInput {
	input := NewCompilerInput(SourceFile_FromCode(name, code))
	input.SetCompiler(pt.compiler)
	return input
}

// Java: PruningAnalysisTest.analyze(ImmutableList<CompilerInput>, CompilerInput...)
func (pt *PruningAnalysisTest) analyze(
	entryPoints []*CompilerInput, allInputs ...*CompilerInput) *PruningAnalysis_Result {
	inputList := append([]*CompilerInput(nil), allInputs...)
	sorter := NewSortedDependencies[*CompilerInput](inputList)
	return PruningAnalysis_Create(sorter, entryPoints).Analyze()
}

// pruningAnalysisTest_mapOf builds the expected map of Truth's
// containsExactly(k0, v0, k1, v1, ...) / containsAtLeast(...) on a map.
func pruningAnalysisTest_mapOf(kv ...any) *jutil.LinkedHashMap[string, int] {
	m := jutil.NewLinkedHashMap[string, int]()
	for i := 0; i < len(kv); i += 2 {
		m.Put(kv[i].(string), kv[i+1].(int))
	}
	return m
}
