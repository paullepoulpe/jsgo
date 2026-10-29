// Port of com/google/javascript/jscomp/graph/CheckPathsBetweenNodesTest.java @ closure-compiler v20261021

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

package graph

import (
	"strings"
	"testing"
)

// checkPathsTest_prefixPredicate is a predicate satisfied by strings with a
// given prefix.
func checkPathsTest_prefixPredicate(prefix string) func(string) bool {
	return func(input string) bool { return strings.HasPrefix(input, prefix) }
}

func checkPathsTest_equalTo(v string) func(string) bool {
	return func(input string) bool { return input == v }
}

var (
	checkPathsTest_FALSE    = func(string) bool { return false }
	checkPathsTest_ALL_EDGE = func(DiGraph_DiGraphEdge[string, string]) bool { return true }
	checkPathsTest_NO_EDGE  = func(DiGraph_DiGraphEdge[string, string]) bool { return false }
)

// Tests straight-line graphs.
func TestCheckPathsBetweenNodes_Simple(t *testing.T) {
	g := LinkedDirectedGraph_Create[string, string]()
	g.CreateNode("a")
	g.CreateNode("b")
	g.CreateNode("c")
	g.CreateNode("d")

	g.Connect("a", "-", "b")
	g.Connect("b", "-", "c")
	g.Connect("c", "-", "d")
	g.Connect("a", "x", "d")

	// Simple case: the sole path from a to d has a matching node.
	checkPathsTest_assertGood(t, checkPathsTest_createTest(g, "a", "d", checkPathsTest_equalTo("b"), checkPathsTest_edgeIs("-")))
	// Test two edge cases where satisfying node is the first and last node on
	// the path.
	checkPathsTest_assertGood(t, checkPathsTest_createTest(g, "a", "d", checkPathsTest_equalTo("a"), checkPathsTest_edgeIs("-")))
	checkPathsTest_assertGood(t, checkPathsTest_createTest(g, "a", "d", checkPathsTest_equalTo("d"), checkPathsTest_edgeIs("-")))

	// Traverse no edges, so no paths.
	checkPathsTest_assertGood(t, checkPathsTest_createTest(g, "a", "d", checkPathsTest_FALSE, checkPathsTest_NO_EDGE))

	// No path with matching edges contains b.
	checkPathsTest_assertBad(t, checkPathsTest_createTest(g, "a", "d", checkPathsTest_equalTo("b"), checkPathsTest_edgeIs("x")))
}

// Tests a graph where some paths between the nodes are valid and others are
// invalid.
func TestCheckPathsBetweenNodes_SomeValidPaths(t *testing.T) {
	g := LinkedDirectedGraph_Create[string, string]()
	g.CreateNode("a")
	g.CreateNode("b")
	g.CreateNode("c")
	g.CreateNode("d")
	g.CreateNode("e")

	g.Connect("a", "1", "b")
	g.Connect("b", "2", "c")
	g.Connect("b", "3", "e")
	g.Connect("e", "4", "d")
	g.Connect("c", "5", "d")

	checkPathsTest_assertBad(t, checkPathsTest_createTest(g, "a", "d", checkPathsTest_equalTo("c"), checkPathsTest_ALL_EDGE))
	checkPathsTest_assertBad(t, checkPathsTest_createTest(g, "a", "d", checkPathsTest_equalTo("z"), checkPathsTest_ALL_EDGE))
}

// Tests a graph with many valid paths.
func TestCheckPathsBetweenNodes_ManyValidPaths(t *testing.T) {
	g := LinkedDirectedGraph_Create[string, string]()
	g.CreateNode("a")
	g.CreateNode("b")
	g.CreateNode("c1")
	g.CreateNode("c2")
	g.CreateNode("c3")
	g.CreateNode("d")

	g.Connect("a", "-", "b")
	g.Connect("b", "-", "c1")
	g.Connect("b", "-", "c2")
	g.Connect("c2", "-", "d")
	g.Connect("c1", "-", "d")
	g.Connect("a", "-", "c3")
	g.Connect("c3", "-", "d")

	checkPathsTest_assertGood(t, checkPathsTest_createTest(g, "a", "d", checkPathsTest_prefixPredicate("c"), checkPathsTest_ALL_EDGE))
}

// Tests a graph with some cycles.
func TestCheckPathsBetweenNodes_Cycles1(t *testing.T) {
	g := LinkedDirectedGraph_Create[string, string]()
	g.CreateNode("a")
	g.CreateNode("b")
	g.CreateNode("c")
	g.CreateNode("d")
	g.CreateNode("e")
	g.CreateNode("f")

	g.Connect("a", "-", "b")
	g.Connect("b", "-", "c")
	g.Connect("c", "-", "d")
	g.Connect("d", "-", "e")
	g.Connect("e", "-", "f")
	g.Connect("f", "-", "b")

	checkPathsTest_assertGood(t, checkPathsTest_createTest(g, "a", "e", checkPathsTest_equalTo("b"), checkPathsTest_ALL_EDGE))
	checkPathsTest_assertGood(t, checkPathsTest_createTest(g, "a", "e", checkPathsTest_equalTo("c"), checkPathsTest_ALL_EDGE))
	checkPathsTest_assertGood(t, checkPathsTest_createTest(g, "a", "e", checkPathsTest_equalTo("d"), checkPathsTest_ALL_EDGE))
	checkPathsTest_assertGood(t, checkPathsTest_createTest(g, "a", "e", checkPathsTest_equalTo("e"), checkPathsTest_ALL_EDGE))
	checkPathsTest_assertBad(t, checkPathsTest_createTest(g, "a", "e", checkPathsTest_equalTo("f"), checkPathsTest_ALL_EDGE))
}

// Tests another graph with cycles. The topology of this graph was inspired
// by a control flow graph that was being incorrectly analyzed by an early
// version of CheckPathsBetweenNodes.
func TestCheckPathsBetweenNodes_Cycles2(t *testing.T) {
	g := LinkedDirectedGraph_Create[string, string]()
	g.CreateNode("a")
	g.CreateNode("b")
	g.CreateNode("c")
	g.CreateNode("d")

	g.Connect("a", "-", "b")
	g.Connect("b", "-", "c")
	g.Connect("c", "-", "b")
	g.Connect("b", "-", "d")

	checkPathsTest_assertGood(t, checkPathsTest_createTest(g, "a", "d", checkPathsTest_equalTo("a"), checkPathsTest_ALL_EDGE))
	checkPathsTest_assertBad(t, checkPathsTest_createTest(g, "a", "d", checkPathsTest_equalTo("z"), checkPathsTest_ALL_EDGE))
}

// Tests another graph with cycles. The topology of this graph was inspired
// by a control flow graph that was being incorrectly analyzed by an early
// version of CheckPathsBetweenNodes.
func TestCheckPathsBetweenNodes_Cycles3(t *testing.T) {
	g := LinkedDirectedGraph_Create[string, string]()
	g.CreateNode("a")
	g.CreateNode("b")
	g.CreateNode("c")
	g.CreateNode("d")

	g.Connect("a", "-", "b")
	g.Connect("b", "-", "c")
	g.Connect("c", "-", "b")
	g.Connect("b", "-", "d")
	g.Connect("c", "-", "d")

	checkPathsTest_assertGood(t, checkPathsTest_createTest(g, "a", "d", checkPathsTest_equalTo("a"), checkPathsTest_ALL_EDGE))
	checkPathsTest_assertBad(t, checkPathsTest_createTest(g, "a", "d", checkPathsTest_equalTo("z"), checkPathsTest_ALL_EDGE))
}

// Much of the tests are done by testing all paths. We quickly verified that
// some paths are indeed correct for the some path case.
func TestCheckPathsBetweenNodes_SomePath1(t *testing.T) {
	g := LinkedDirectedGraph_Create[string, string]()
	g.CreateNode("a")
	g.CreateNode("b")
	g.CreateNode("c")
	g.CreateNode("d")

	g.Connect("a", "-", "b")
	g.Connect("a", "-", "c")
	g.Connect("b", "-", "d")
	g.Connect("c", "-", "d")

	assertTrue(t, "b", checkPathsTest_createTest(g, "a", "d", checkPathsTest_equalTo("b"), checkPathsTest_ALL_EDGE).SomePathsSatisfyPredicate())
	assertTrue(t, "c", checkPathsTest_createTest(g, "a", "d", checkPathsTest_equalTo("c"), checkPathsTest_ALL_EDGE).SomePathsSatisfyPredicate())
	assertTrue(t, "a", checkPathsTest_createTest(g, "a", "d", checkPathsTest_equalTo("a"), checkPathsTest_ALL_EDGE).SomePathsSatisfyPredicate())
	assertTrue(t, "d", checkPathsTest_createTest(g, "a", "d", checkPathsTest_equalTo("d"), checkPathsTest_ALL_EDGE).SomePathsSatisfyPredicate())
	assertFalse(t, "NONE", checkPathsTest_createTest(g, "a", "d", checkPathsTest_equalTo("NONE"), checkPathsTest_ALL_EDGE).SomePathsSatisfyPredicate())
}

func TestCheckPathsBetweenNodes_SomePath2(t *testing.T) {
	// No Paths between nodes, by definition, always false.
	g := LinkedDirectedGraph_Create[string, string]()
	g.CreateNode("a")
	g.CreateNode("b")

	assertFalse(t, "b", checkPathsTest_createTest(g, "a", "b", checkPathsTest_equalTo("b"), checkPathsTest_ALL_EDGE).SomePathsSatisfyPredicate())
	assertFalse(t, "d", checkPathsTest_createTest(g, "a", "b", checkPathsTest_equalTo("d"), checkPathsTest_ALL_EDGE).SomePathsSatisfyPredicate())
	assertTrue(t, "a", checkPathsTest_createTest(g, "a", "b", checkPathsTest_equalTo("a"), checkPathsTest_ALL_EDGE).SomePathsSatisfyPredicate())
}

func TestCheckPathsBetweenNodes_SomePathRevisiting(t *testing.T) {
	g := LinkedDirectedGraph_Create[string, string]()
	g.CreateNode("1")
	g.CreateNode("2a")
	g.CreateNode("2b")
	g.CreateNode("3")
	g.CreateNode("4a")
	g.CreateNode("4b")
	g.CreateNode("5")
	g.Connect("1", "-", "2a")
	g.Connect("1", "-", "2b")
	g.Connect("2a", "-", "3")
	g.Connect("2b", "-", "3")
	g.Connect("3", "-", "4a")
	g.Connect("3", "-", "4b")
	g.Connect("4a", "-", "5")
	g.Connect("4b", "-", "5")

	// CountingPredicate<String>(Predicates.equalTo("4a"))
	count := 0
	delegate := checkPathsTest_equalTo("4a")
	p := func(input string) bool {
		count++
		return delegate(input)
	}

	assertTrue(t, "some", checkPathsTest_createTest(g, "1", "5", p, checkPathsTest_ALL_EDGE).SomePathsSatisfyPredicate())

	// Make sure we are not doing more traversals than we have to.
	assertEqual(t, "count", count, 4)
}

func TestCheckPathsBetweenNodes_NonInclusive(t *testing.T) {
	// No Paths between nodes, by definition, always false.
	g := LinkedDirectedGraph_Create[string, string]()
	g.CreateNode("a")
	g.CreateNode("b")
	g.CreateNode("c")
	g.Connect("a", "-", "b")
	g.Connect("b", "-", "c")
	assertFalse(t, "a", checkPathsTest_createNonInclusiveTest(g, "a", "b", checkPathsTest_equalTo("a"), checkPathsTest_ALL_EDGE).SomePathsSatisfyPredicate())
	assertFalse(t, "b", checkPathsTest_createNonInclusiveTest(g, "a", "b", checkPathsTest_equalTo("b"), checkPathsTest_ALL_EDGE).SomePathsSatisfyPredicate())
	assertTrue(t, "b", checkPathsTest_createNonInclusiveTest(g, "a", "c", checkPathsTest_equalTo("b"), checkPathsTest_ALL_EDGE).SomePathsSatisfyPredicate())
}

func TestCheckPathsBetweenNodes_WithoutAnnotationsAndSideEffectFree(t *testing.T) {
	// Works on graphs with annotations disabled.
	noAnn := LinkedDirectedGraph_CreateWithoutAnnotations[string, string]()
	a := noAnn.CreateLinkedDiGraphNode("a")
	b := noAnn.CreateLinkedDiGraphNode("b")
	c := noAnn.CreateLinkedDiGraphNode("c")
	d := noAnn.CreateLinkedDiGraphNode("d")
	p := noAnn.CreateLinkedDiGraphNode("p")
	// Multi-entry cycle where one entry passes through p and the other bypasses p:
	// a -> p -> b <-> c, b -> d, a -> c.
	noAnn.ConnectLinkedNodes(a, "-", p)
	noAnn.ConnectLinkedNodes(p, "-", b)
	noAnn.ConnectLinkedNodes(b, "-", c)
	noAnn.ConnectLinkedNodes(c, "-", b)
	noAnn.ConnectLinkedNodes(b, "-", d)
	noAnn.ConnectLinkedNodes(a, "-", c)

	var checker CheckPathsBetweenNodes[string, string]
	checker.ResetLinkedAllEdges(noAnn, a, d, checkPathsTest_equalTo("p"), true)
	assertFalse(t, "allPaths bypass via a->c->b->d", checker.AllPathsSatisfyPredicate())
	assertTrue(t, "somePaths via a->p->b->d", checker.SomePathsSatisfyPredicate())

	checker.ResetLinkedAllEdges(noAnn, a, d, checkPathsTest_equalTo("b"), false)
	assertTrue(t, "allPaths must pass through b", checker.AllPathsSatisfyPredicate())
	assertTrue(t, "somePaths pass through b", checker.SomePathsSatisfyPredicate())

	// Preserves existing node and edge annotations on annotated graphs.
	ann := LinkedDirectedGraph_Create[string, string]()
	na := ann.CreateLinkedDiGraphNode("a")
	nb := ann.CreateLinkedDiGraphNode("b")
	nc := ann.CreateLinkedDiGraphNode("c")
	ann.ConnectLinkedNodes(na, "-", nb)
	ann.ConnectLinkedNodes(nb, "-", nb) // self-loop
	ann.ConnectLinkedNodes(nb, "-", nc)

	type marker struct{ id int }
	nodeMarker := &marker{id: 42}
	edgeMarker := &marker{id: 99}
	na.SetAnnotation(nodeMarker)
	nb.SetAnnotation(nodeMarker)
	for _, e := range nb.OutEdgeList.Slice() {
		e.SetAnnotation(edgeMarker)
	}

	checker.ResetLinkedAllEdges(ann, na, nc, checkPathsTest_equalTo("b"), false)
	assertTrue(t, "self-loop allPaths", checker.AllPathsSatisfyPredicate())
	assertTrue(t, "self-loop somePaths", checker.SomePathsSatisfyPredicate())
	assertEqual(t, "node annotation preserved", na.GetAnnotation(), Annotation(nodeMarker))
	assertEqual(t, "node annotation preserved", nb.GetAnnotation(), Annotation(nodeMarker))
	for _, e := range nb.OutEdgeList.Slice() {
		assertEqual(t, "edge annotation preserved", e.GetAnnotation(), Annotation(edgeMarker))
	}
}

func TestCheckPathsBetweenNodes_LargeGraphReuse(t *testing.T) {
	lg := LinkedDirectedGraph_CreateWithoutAnnotations[int, int]()
	const n = 100
	nodes := make([]*LinkedDirectedGraph_LinkedDiGraphNode[int, int], n)
	for i := range n {
		nodes[i] = lg.CreateLinkedDiGraphNode(i)
	}
	for i := range n - 1 {
		lg.ConnectLinkedNodes(nodes[i], 0, nodes[i+1])
		if i > 0 {
			lg.ConnectLinkedNodes(nodes[i], 1, nodes[i-1])
		}
	}
	var checker CheckPathsBetweenNodes[int, int]
	for target := 1; target < n-1; target++ {
		want := target
		checker.ResetLinkedAllEdges(lg, nodes[0], nodes[n-1], func(v int) bool { return v == want }, false)
		assertTrue(t, "chain allPaths", checker.AllPathsSatisfyPredicate())
		assertTrue(t, "chain somePaths", checker.SomePathsSatisfyPredicate())
	}
	checker.ResetLinkedAllEdges(lg, nodes[0], nodes[n-1], func(v int) bool { return v == -1 }, false)
	assertFalse(t, "chain missing allPaths", checker.AllPathsSatisfyPredicate())
	assertFalse(t, "chain missing somePaths", checker.SomePathsSatisfyPredicate())
}

func checkPathsTest_assertGood[N comparable, E any](t *testing.T, test *CheckPathsBetweenNodes[N, E]) {
	t.Helper()
	assertTrue(t, "allPathsSatisfyPredicate", test.AllPathsSatisfyPredicate())
}

func checkPathsTest_assertBad[N comparable, E any](t *testing.T, test *CheckPathsBetweenNodes[N, E]) {
	t.Helper()
	assertFalse(t, "allPathsSatisfyPredicate", test.AllPathsSatisfyPredicate())
}

func checkPathsTest_createTest(
	graph *LinkedDirectedGraph[string, string],
	entry string,
	exit string,
	nodePredicate func(string) bool,
	edgePredicate func(*LinkedDiGraphEdge[string, string]) bool) *CheckPathsBetweenNodes[string, string] {
	return NewCheckPathsBetweenNodesInclusive(
		graph, graph.GetDiGraphNode(entry), graph.GetDiGraphNode(exit), nodePredicate, edgePredicate)
}

func checkPathsTest_createNonInclusiveTest(
	graph *LinkedDirectedGraph[string, string],
	entry string,
	exit string,
	nodePredicate func(string) bool,
	edgePredicate func(*LinkedDiGraphEdge[string, string]) bool) *CheckPathsBetweenNodes[string, string] {
	return NewCheckPathsBetweenNodes(
		graph, graph.GetDiGraphNode(entry), graph.GetDiGraphNode(exit), nodePredicate, edgePredicate, false)
}

func checkPathsTest_edgeIs(val string) func(*LinkedDiGraphEdge[string, string]) bool {
	return func(input *LinkedDiGraphEdge[string, string]) bool { return input.GetValue() == val }
}
