// Port of com/google/javascript/jscomp/graph/LowestCommonAncestorFinderTest.java @ closure-compiler v20261015

/*
 * Copyright 2008 The Closure Compiler Authors.
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
	"strconv"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Java's edge type is Void (only null); E is `any` holding nil here.

func TestLowestCommonAncestorFinder_FindAll_onTree_withCommonParent(t *testing.T) {
	newLcaCaseBuilder(t). //
				edge(3, 1).
				edge(3, 2).
				roots(1, 2).
				expect(3).
				testAndRenderGraph()
}

func TestLowestCommonAncestorFinder_FindAll_onTree_withCommonGrandParent(t *testing.T) {
	newLcaCaseBuilder(t). //
				edge(5, 4).
				edge(5, 3).
				edge(4, 2).
				edge(3, 1).
				roots(1, 2).
				expect(5).
				testAndRenderGraph()
}

func TestLowestCommonAncestorFinder_FindAll_onTree_whenLcaHasParent(t *testing.T) {
	newLcaCaseBuilder(t). //
				edge(4, 3).
				edge(3, 1).
				edge(3, 2).
				roots(1, 2).
				expect(3).
				testAndRenderGraph()
}

func TestLowestCommonAncestorFinder_FindAll_onDag_withNestedVs_findsSingleLca(t *testing.T) {
	newLcaCaseBuilder(t). //
				edge(4, 3).
				edge(4, 2).
				edge(4, 1).
				edge(3, 2).
				edge(3, 1).
				roots(1, 2).
				expect(3).
				testAndRenderGraph()
}

func TestLowestCommonAncestorFinder_FindAll_onDag_findsMultipleLcas(t *testing.T) {
	newLcaCaseBuilder(t). //
				edge(4, 2).
				edge(4, 1).
				edge(3, 2).
				edge(3, 1).
				roots(1, 2).
				expect(3, 4).
				testAndRenderGraph()
}

func TestLowestCommonAncestorFinder_FindAll_onDag_withCommonParent_findsMultipleLcas(t *testing.T) {
	newLcaCaseBuilder(t). //
				edge(5, 4).
				edge(5, 3).
				edge(4, 2).
				edge(4, 1).
				edge(3, 2).
				edge(3, 1).
				roots(1, 2).
				expect(3, 4).
				testAndRenderGraph()
}

func TestLowestCommonAncestorFinder_FindAll_onDag_withCommonParent_findsMultipleLcas_atDifferentDistances(t *testing.T) {
	newLcaCaseBuilder(t).
		edge(10, 3).
		edge(10, 9).
		edge(9, 8).
		edge(9, 7).
		edge(8, 6).
		edge(7, 5).
		edge(6, 1).
		edge(5, 2).
		edge(3, 2).
		edge(3, 1).
		roots(1, 2).
		expect(3, 9).
		testAndRenderGraph()
}

func TestLowestCommonAncestorFinder_FindAll_onDag_sidePathsIntoNonLowest(t *testing.T) {
	newLcaCaseBuilder(t).
		edge(20, 18).
		edge(20, 19).
		edge(19, 17).
		edge(19, 12).
		edge(18, 11).
		edge(17, 15).
		edge(16, 15).
		edge(15, 14).
		edge(15, 2).
		edge(14, 1).
		edge(12, 9).
		edge(11, 10).
		edge(10, 1).
		edge(9, 2).
		roots(1, 2).
		expect(15).
		testAndRenderGraph()
}

func TestLowestCommonAncestorFinder_FindAll_onCyclic_whenLcaInsideCycle_findsSomeMember_deterministically(t *testing.T) {
	newLcaCaseBuilder(t).
		edge(1, 2).
		edge(2, 3).
		edge(3, 4).
		edge(4, 1).
		edge(1, 5).
		edge(2, 6).
		roots(5, 6).
		expect(1).
		testAndRenderGraph()
}

func TestLowestCommonAncestorFinder_FindAll_onCyclic_whenLcaInsideCycle_appearsUnique_mayFindOtherMember(t *testing.T) {
	newLcaCaseBuilder(t).
		edge(1, 2).
		edge(2, 3).
		edge(3, 4).
		edge(4, 1).
		edge(1, 5).
		// Notice how 3 is a parent of both 5 and 6, yet isn't selected.
		edge(3, 5).
		edge(3, 6).
		roots(5, 6).
		expect(1).
		testAndRenderGraph()
}

func TestLowestCommonAncestorFinder_FindAll_onCyclic_whenLcaBelowCycle_ignoresCycle(t *testing.T) {
	newLcaCaseBuilder(t).
		edge(1, 2).
		edge(2, 3).
		edge(3, 1).
		edge(3, 4).
		edge(4, 5).
		edge(4, 6).
		roots(5, 6).
		expect(4).
		testAndRenderGraph()
}

func TestLowestCommonAncestorFinder_FindAll_onDigitLattice_withUnrelatedRoots_findsMultipleLcas(t *testing.T) {
	newLcaCaseBuilder(t). //
				copyGraph(lcaTest_DIGIT_SUBSTRING_LATTICE).
				roots(1, 3, 5).
				expect(135, 153, 315, 351, 531, 513).
				testAndRenderGraph()
}

func TestLowestCommonAncestorFinder_FindAll_onDigitLattice_withRelatedRoots_findsMultipleLcas(t *testing.T) {
	newLcaCaseBuilder(t). //
				copyGraph(lcaTest_DIGIT_SUBSTRING_LATTICE).
				roots(1, 13, 5).
				expect(135, 513).
				testAndRenderGraph()
}

func TestLowestCommonAncestorFinder_FindAll_onDigitLattice_withOverlappingRoots_findsSingleLca(t *testing.T) {
	newLcaCaseBuilder(t). //
				copyGraph(lcaTest_DIGIT_SUBSTRING_LATTICE).
				roots(7, 4, 8, 48, 74).
				expect(748).
				testAndRenderGraph()
}

func TestLowestCommonAncestorFinder_FindAll_onDigitLattice_noLca(t *testing.T) {
	newLcaCaseBuilder(t).
		copyGraph(lcaTest_DIGIT_SUBSTRING_LATTICE). //
		roots(2, 4, 7, 8).
		expect().
		testAndRenderGraph()
}

func TestLowestCommonAncestorFinder_FindAll_onDisjoint_noSolution(t *testing.T) {
	newLcaCaseBuilder(t). //
				edge(4, 2).
				edge(3, 1).
				roots(1, 2).
				expect().
				testAndRenderGraph()
}

func TestLowestCommonAncestorFinder_FindAll_rejectsTooManyRoots(t *testing.T) {
	// Given
	builder := newLcaCaseBuilder(t).expect(40)
	roots := make([]int, 32) // Integer.SIZE
	for i := range roots {
		roots[i] = i
		builder.edge(40, roots[i])
	}
	builder.roots(roots...)

	// Then
	assertThrows[jutil.Throwable](t, builder.testAndRenderGraph)
}

type lcaCaseBuilder struct {
	t        *testing.T
	graph    *LinkedDirectedGraph[int, any]
	roots_   *jutil.LinkedHashSet[int]
	expected *jutil.LinkedHashSet[int]
	actual   *jutil.LinkedHashSet[int]
}

func newLcaCaseBuilder(t *testing.T) *lcaCaseBuilder {
	return &lcaCaseBuilder{t: t, graph: LinkedDirectedGraph_Create[int, any]()}
}

func (b *lcaCaseBuilder) edge(src int, dest int) *lcaCaseBuilder {
	b.graph.CreateNode(src)
	b.graph.CreateNode(dest)
	b.graph.Connect(src, nil, dest)
	if n := len(b.graph.GetEdgesBetween(src, dest)); n != 1 {
		b.t.Fatalf("edges(%d, %d): got %d, want 1", src, dest, n)
	}
	return b
}

func (b *lcaCaseBuilder) copyGraph(src *LinkedDirectedGraph[int, any]) *lcaCaseBuilder {
	for _, edge := range src.GetDiGraphEdges() {
		b.edge(edge.GetSource().GetValue(), edge.GetDestination().GetValue())
	}
	return b
}

func (b *lcaCaseBuilder) roots(roots ...int) *lcaCaseBuilder {
	jutil.CheckState(b.roots_ == nil)
	b.roots_ = jutil.NewLinkedHashSet[int]()
	b.roots_.AddAll(roots...)
	return b
}

func (b *lcaCaseBuilder) expect(expected ...int) *lcaCaseBuilder {
	jutil.CheckState(b.expected == nil)
	b.expected = jutil.NewLinkedHashSet[int]()
	b.expected.AddAll(expected...)
	return b
}

func (b *lcaCaseBuilder) testAndRenderGraph() {
	b.t.Helper()
	jutil.CheckState(b.actual == nil)
	var finder LowestCommonAncestorFinder[int, any] = NewLowestCommonAncestorFinder[int, any](b.graph)
	b.actual = finder.FindAll(b.roots_)

	message := ""
	if b.renderGraph() {
		message = "Look for an SVG of the test graph in the test artifacts"
	}
	assertSameElements(b.t, "findAll "+message, b.actual.ToSlice(), b.expected.ToSlice())
}

func (b *lcaCaseBuilder) renderGraph() bool {
	if b.graph.GetNodeCount() > 100 {
		return false // This graph won't be readable anyway.
	}

	for _, node := range b.graph.GetDiGraphNodes() {
		purpose := lcaTest_NodePurpose_NONE
		if b.roots_.Contains(node.GetValue()) {
			purpose = purpose.mix(lcaTest_NodePurpose_ROOT)
		}
		if b.expected.Contains(node.GetValue()) {
			purpose = purpose.mix(lcaTest_NodePurpose_EXPECTED)
		}
		if b.actual.Contains(node.GetValue()) {
			purpose = purpose.mix(lcaTest_NodePurpose_ACTUAL)
		}
		node.SetAnnotation(purpose)
	}

	return true
}

// lcaTest_NodePurpose is the enum NodePurpose (an Annotation).
type lcaTest_NodePurpose int32

// The order of these elements it loadbearing.
const (
	lcaTest_NodePurpose_NONE lcaTest_NodePurpose = iota
	lcaTest_NodePurpose_ACTUAL
	lcaTest_NodePurpose_EXPECTED
	lcaTest_NodePurpose_ACTUAL_EXPECTED
	lcaTest_NodePurpose_ROOT
	lcaTest_NodePurpose_ACTUAL_ROOT
	lcaTest_NodePurpose_EXPECTED_ROOT
	lcaTest_NodePurpose_ACTUAL_EXPECTED_ROOT
)

func (p lcaTest_NodePurpose) mix(x lcaTest_NodePurpose) lcaTest_NodePurpose {
	return p | x // values()[this.ordinal() | x.ordinal()]
}

func (p lcaTest_NodePurpose) String() string {
	return [...]string{"NONE", "ACTUAL", "EXPECTED", "ACTUAL_EXPECTED", "ROOT", "ACTUAL_ROOT",
		"EXPECTED_ROOT", "ACTUAL_EXPECTED_ROOT"}[p]
}

// lcaTest_DIGIT_SUBSTRING_LATTICE is a graph of integers in which each value
// (X) connects to a value (Y) if the decimal representation of (Y) appears
// as a substring in the decimal representation of (X).
//
// The largest element is 999. There are no transitive edges.
//
// TODO(nickreid): Switch to a graph that's simpler to describe but has the
// same non-trivial structural aspects.
var lcaTest_DIGIT_SUBSTRING_LATTICE = lcaTest_makeDigitLattice()

func lcaTest_makeDigitLattice() *LinkedDirectedGraph[int, any] {
	result := LinkedDirectedGraph_CreateWithoutAnnotations[int, any]()

	for i := 0; i < 1000; i++ {
		iStr := strconv.Itoa(i)
		windowSize := len(iStr) - 1

		result.CreateNode(i)

		if windowSize == 0 {
			continue
		}
		for offset := 0; offset+windowSize <= len(iStr); offset++ {
			descStr := iStr[offset : offset+windowSize]
			desc, _ := strconv.Atoi(descStr)
			result.ConnectIfNotConnectedInDirection(i, nil, desc)
		}
	}

	return result
}
