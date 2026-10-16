// Port of com/google/javascript/jscomp/graph/GraphTest.java @ closure-compiler v20260415

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
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func assertTrue(t *testing.T, what string, b bool) {
	t.Helper()
	if !b {
		t.Errorf("%s: expected true", what)
	}
}

func assertFalse(t *testing.T, what string, b bool) {
	t.Helper()
	if b {
		t.Errorf("%s: expected false", what)
	}
}

func TestGraph_DirectedSimple(t *testing.T) {
	var graph DiGraph[string, string] = LinkedDirectedGraph_Create[string, string]()
	graph.CreateNode("a")
	graph.CreateNode("b")
	graph.CreateNode("c")
	graph.Connect("a", "->", "b")
	assertTrue(t, "hasNode(a)", graph.HasNode("a"))
	assertTrue(t, "hasNode(b)", graph.HasNode("b"))
	assertTrue(t, "hasNode(c)", graph.HasNode("c"))
	assertFalse(t, "hasNode(d)", graph.HasNode("d"))
	assertTrue(t, "a-b", graph.IsConnected("a", "b"))
	assertTrue(t, "b-a", graph.IsConnected("b", "a"))
	assertFalse(t, "a-c", graph.IsConnected("a", "c"))
	assertFalse(t, "b-c", graph.IsConnected("b", "c"))
	assertFalse(t, "c-a", graph.IsConnected("c", "a"))
	assertFalse(t, "c-b", graph.IsConnected("c", "b"))
	assertFalse(t, "a-a", graph.IsConnected("a", "a"))
	assertFalse(t, "b-b", graph.IsConnected("b", "b"))
	assertFalse(t, "b-c", graph.IsConnected("b", "c"))
	assertTrue(t, "a->b", graph.IsConnectedInDirection("a", "b"))
	assertFalse(t, "b->a", graph.IsConnectedInDirection("b", "a"))
	assertFalse(t, "a->c", graph.IsConnectedInDirection("a", "c"))
	assertFalse(t, "b->c", graph.IsConnectedInDirection("b", "c"))
	assertFalse(t, "c->a", graph.IsConnectedInDirection("c", "a"))
	assertFalse(t, "c->b", graph.IsConnectedInDirection("c", "b"))

	// Removal.
	graph.Disconnect("a", "b")
	assertFalse(t, "a-b", graph.IsConnected("a", "b"))
	assertFalse(t, "b-a", graph.IsConnected("b", "a"))

	// Disconnect both ways.
	graph.Connect("a", "->", "b")
	graph.Connect("b", "->", "a")
	graph.Disconnect("a", "b")
	assertFalse(t, "a-b", graph.IsConnected("a", "b"))
	assertFalse(t, "b-a", graph.IsConnected("b", "a"))

	// Disconnect one way.
	graph.Connect("a", "->", "b")
	graph.Connect("b", "->", "a")
	graph.DisconnectInDirection("a", "b")
	assertTrue(t, "b-a", graph.IsConnected("b", "a"))
	assertTrue(t, "a-b", graph.IsConnected("a", "b"))
	assertFalse(t, "a->b", graph.IsConnectedInDirection("a", "b"))
	assertTrue(t, "b->a", graph.IsConnectedInDirection("b", "a"))
}

func TestGraph_UndirectedSimple(t *testing.T) {
	var graph UndiGraph[string, string] = LinkedUndirectedGraph_Create[string, string]()
	graph.CreateNode("a")
	graph.CreateNode("b")
	graph.CreateNode("c")
	graph.Connect("a", "--", "b")
	assertTrue(t, "hasNode(a)", graph.HasNode("a"))
	assertTrue(t, "hasNode(b)", graph.HasNode("b"))
	assertTrue(t, "hasNode(c)", graph.HasNode("c"))
	assertFalse(t, "hasNode(d)", graph.HasNode("d"))
	assertTrue(t, "a-b", graph.IsConnected("a", "b"))
	assertTrue(t, "b-a", graph.IsConnected("b", "a"))
	assertFalse(t, "a-c", graph.IsConnected("a", "c"))
	assertFalse(t, "b-c", graph.IsConnected("b", "c"))
	assertFalse(t, "c-a", graph.IsConnected("c", "a"))
	assertFalse(t, "c-b", graph.IsConnected("c", "b"))
	assertFalse(t, "a-a", graph.IsConnected("a", "a"))
	assertFalse(t, "b-b", graph.IsConnected("b", "b"))
	assertFalse(t, "b-c", graph.IsConnected("b", "c"))

	// Removal.
	graph.Disconnect("a", "b")
	assertFalse(t, "a-b", graph.IsConnected("a", "b"))
	assertFalse(t, "b-a", graph.IsConnected("b", "a"))
}

func TestGraph_DirectedSelfLoop(t *testing.T) {
	var graph DiGraph[string, string] = LinkedDirectedGraph_Create[string, string]()
	graph.CreateNode("a")
	graph.CreateNode("b")
	graph.Connect("a", "->", "a")
	assertTrue(t, "a-a", graph.IsConnected("a", "a"))
	assertFalse(t, "a-b", graph.IsConnected("a", "b"))
	assertFalse(t, "b-a", graph.IsConnected("b", "a"))
	assertTrue(t, "a->a", graph.IsConnectedInDirection("a", "a"))
	assertFalse(t, "a->b", graph.IsConnectedInDirection("a", "b"))
	assertFalse(t, "b->a", graph.IsConnectedInDirection("b", "a"))

	// Removal.
	graph.Disconnect("a", "a")
	assertFalse(t, "a-a", graph.IsConnected("a", "a"))

	// Disconnect both ways.
	graph.Connect("a", "->", "a")
	graph.Disconnect("a", "a")
	assertFalse(t, "a-a", graph.IsConnected("a", "a"))
	assertFalse(t, "a-a", graph.IsConnected("a", "a"))

	// Disconnect one way.
	graph.Connect("a", "->", "a")
	graph.DisconnectInDirection("a", "a")
	assertFalse(t, "a-a", graph.IsConnected("a", "a"))
}

func TestGraph_UndirectedSelfLoop(t *testing.T) {
	var graph UndiGraph[string, string] = LinkedUndirectedGraph_Create[string, string]()
	graph.CreateNode("a")
	graph.CreateNode("b")
	graph.Connect("a", "--", "a")
	assertTrue(t, "a-a", graph.IsConnected("a", "a"))
	assertFalse(t, "a-b", graph.IsConnected("a", "b"))
	assertFalse(t, "b-a", graph.IsConnected("b", "a"))

	// Removal.
	graph.Disconnect("a", "a")
	assertFalse(t, "a-a", graph.IsConnected("a", "a"))
}

func TestGraph_DirectedInAndOutEdges(t *testing.T) {
	var graph DiGraph[string, string] = LinkedDirectedGraph_Create[string, string]()
	graph.CreateNode("a")
	graph.CreateNode("b")
	graph.CreateNode("c")
	graph.CreateNode("d")
	graph.Connect("a", "->", "b")
	graph.Connect("a", "-->", "b")
	graph.Connect("a", "--->", "b")
	graph.Connect("a", "->", "c")
	graph.Connect("c", "->", "d")
	graphTest_assertSetEquals(t, graph.GetDirectedSuccNodesOfValue("a"), "b", "c")
	graphTest_assertSetEquals(t, graph.GetDirectedPredNodesOfValue("b"), "a")
	graphTest_assertSetEquals(t, graph.GetDirectedPredNodesOfValue("c"), "a")
	graphTest_assertListCount(t, graph.GetDirectedSuccNodesOfValue("a"), "b", 3)

	// Removal.
	graph.Disconnect("a", "b")
	assertFalse(t, "a-b", graph.IsConnected("a", "b"))
}

func TestGraph_UndirectedNeighbors(t *testing.T) {
	var graph UndiGraph[string, string] = LinkedUndirectedGraph_Create[string, string]()
	graph.CreateNode("a")
	graph.CreateNode("b")
	graph.CreateNode("c")
	graph.CreateNode("d")
	graph.Connect("a", "-", "b")
	graph.Connect("a", "--", "b")
	graph.Connect("a", "---", "b")
	graph.Connect("a", "-", "c")
	graph.Connect("c", "-", "d")
	graphTest_assertSetEquals(t, graph.GetNeighborNodes("a"), "b", "c")
	graphTest_assertSetEquals(t, graph.GetNeighborNodes("b"), "a")
	graphTest_assertSetEquals(t, graph.GetNeighborNodes("c"), "a", "d")
	graphTest_assertListCount(t, graph.GetNeighborNodes("a"), "b", 3)

	// Removal.
	graph.Disconnect("a", "b")
	assertFalse(t, "a-b", graph.IsConnected("a", "b"))
}

func TestGraph_DirectedGetFirstEdge(t *testing.T) {
	var graph DiGraph[string, string] = LinkedDirectedGraph_Create[string, string]()
	graph.CreateNode("a")
	graph.CreateNode("b")
	graph.CreateNode("c")
	graph.Connect("a", "-", "b")
	assertEqual(t, "first(a,b)", graph.GetFirstEdge("a", "b").GetValue(), "-")
	assertEqual(t, "first(b,a)", graph.GetFirstEdge("b", "a").GetValue(), "-")
	assertTrue(t, "first(a,c) is null", graph.GetFirstEdge("a", "c") == nil)
}

func TestGraph_UndirectedGetFirstEdge(t *testing.T) {
	var graph UndiGraph[string, string] = LinkedUndirectedGraph_Create[string, string]()
	graph.CreateNode("a")
	graph.CreateNode("b")
	graph.CreateNode("c")
	graph.Connect("a", "-", "b")
	assertEqual(t, "first(a,b)", graph.GetFirstEdge("a", "b").GetValue(), "-")
	assertEqual(t, "first(b,a)", graph.GetFirstEdge("b", "a").GetValue(), "-")
	assertTrue(t, "first(a,c) is null", graph.GetFirstEdge("a", "c") == nil)
}

func TestGraph_NodeAnnotations(t *testing.T) {
	var graph Graph[string, string] = LinkedUndirectedGraph_Create[string, string]()
	a := graph.CreateNode("a")
	b := graph.CreateNode("b")
	graphTest_checkAnnotations(t, graph, a, b)
}

func TestGraph_EdgeAnnotations(t *testing.T) {
	var graph Graph[string, string] = LinkedUndirectedGraph_Create[string, string]()
	graph.CreateNode("1")
	graph.CreateNode("2")
	graph.CreateNode("3")
	graph.Connect("1", "a", "2")
	graph.Connect("2", "b", "3")
	a := graph.GetEdgesBetween("1", "2")[0]
	b := graph.GetEdgesBetween("2", "3")[0]
	graphTest_checkAnnotations(t, graph, a, b)
}

// graphTest_annotation is the anonymous `new Annotation() {}` class.
type graphTest_annotation struct{ name string }

func graphTest_checkAnnotations(t *testing.T, graph Graph[string, string], a Annotatable, b Annotatable) {
	t.Helper()
	var A Annotation = &graphTest_annotation{"A"}
	var B Annotation = &graphTest_annotation{"B"}

	// Initially null.
	assertTrue(t, "a null", a.GetAnnotation() == nil)
	assertTrue(t, "b null", b.GetAnnotation() == nil)

	// Test basic setting.
	a.SetAnnotation(A)
	b.SetAnnotation(B)
	assertTrue(t, "a is A", a.GetAnnotation() == A)
	assertTrue(t, "b is B", b.GetAnnotation() == B)

	// Test clearing.
	graph.ClearEdgeAnnotations()
	graph.ClearNodeAnnotations()
	assertTrue(t, "a null", a.GetAnnotation() == nil)
	assertTrue(t, "b null", b.GetAnnotation() == nil)

	a.SetAnnotation(A)
	b.SetAnnotation(B)
	// Pushing clears.
	graph.PushEdgeAnnotations()
	graph.PushNodeAnnotations()
	assertTrue(t, "a null", a.GetAnnotation() == nil)
	assertTrue(t, "b null", b.GetAnnotation() == nil)
	a.SetAnnotation(B)
	b.SetAnnotation(B)
	graph.PushEdgeAnnotations()
	graph.PushNodeAnnotations()
	a.SetAnnotation(B)
	b.SetAnnotation(A)

	// Test restoring then restoring old values with pop.
	assertTrue(t, "a is B", a.GetAnnotation() == B)
	assertTrue(t, "b is A", b.GetAnnotation() == A)
	graph.PopEdgeAnnotations()
	graph.PopNodeAnnotations()
	assertTrue(t, "a is B", a.GetAnnotation() == B)
	assertTrue(t, "b is B", b.GetAnnotation() == B)
	graph.PopEdgeAnnotations()
	graph.PopNodeAnnotations()
	assertTrue(t, "a is A", a.GetAnnotation() == A)
	assertTrue(t, "b is B", b.GetAnnotation() == B)
}

func TestGraph_Degree(t *testing.T) {
	graphTest_testDirectedDegree(t, LinkedDirectedGraph_Create[string, string]())
	graphTest_testDirectedDegree(t, LinkedUndirectedGraph_Create[string, string]())
}

func graphTest_testDirectedDegree(t *testing.T, graph Graph[string, string]) {
	t.Helper()
	graph.CreateNode("a")
	graph.CreateNode("b")
	graph.CreateNode("c")
	graph.CreateNode("d")
	assertEqual(t, "degree", graph.GetNodeDegree("a"), 0)
	graph.Connect("a", "-", "b")
	assertEqual(t, "degree", graph.GetNodeDegree("a"), 1)
	graph.Connect("b", "-", "c")
	assertEqual(t, "degree", graph.GetNodeDegree("a"), 1)
	graph.Connect("a", "-", "c")
	assertEqual(t, "degree", graph.GetNodeDegree("a"), 2)
	graph.Connect("d", "-", "a")
	assertEqual(t, "degree", graph.GetNodeDegree("a"), 3)
}

func TestGraph_DirectedConnectIfNotFound(t *testing.T) {
	// E is `any` so that the Java null edge value is representable (nil).
	graphTest_testDirectedConnectIfNotFound(t, LinkedDirectedGraph_Create[string, any]())
	graphTest_testDirectedConnectIfNotFound(t, LinkedUndirectedGraph_Create[string, any]())
}

func graphTest_testDirectedConnectIfNotFound(t *testing.T, graph Graph[string, any]) {
	t.Helper()
	graph.CreateNode("a")
	graph.CreateNode("b")
	graph.ConnectIfNotFound("a", "-", "b")
	assertEqual(t, "degree", graph.GetNodeDegree("a"), 1)
	graph.ConnectIfNotFound("a", "-", "b")
	assertEqual(t, "degree", graph.GetNodeDegree("a"), 1)
	graph.ConnectIfNotFound("a", nil, "b")
	assertEqual(t, "degree", graph.GetNodeDegree("a"), 2)
	graph.ConnectIfNotFound("a", nil, "b")
	assertEqual(t, "degree", graph.GetNodeDegree("a"), 2)
}

func TestGraph_SimpleSubGraph(t *testing.T) {
	var graph UndiGraph[string, string] = LinkedUndirectedGraph_Create[string, string]()
	graph.CreateNode("a")
	graph.CreateNode("b")
	graph.CreateNode("c")
	graph.Connect("a", "--", "b")

	subGraph := graph.NewSubGraph()
	subGraph.AddNode("a")
	subGraph.AddNode("b")

	// SubGraph should not allow add for node that is not in graph.
	assertThrows[jutil.IIllegalArgumentException](t, func() { subGraph.AddNode("d") })

	assertFalse(t, "independent of a", subGraph.IsIndependentOf("a"))
	assertFalse(t, "independent of b", subGraph.IsIndependentOf("b"))
	assertTrue(t, "independent of c", subGraph.IsIndependentOf("c"))
}

func graphTest_assertListCount[T GraphNode[string, string]](t *testing.T, list []T, target string, count int) {
	t.Helper()
	for _, node := range list {
		if node.GetValue() == target {
			count--
		}
	}
	assertEqual(t, "count", count, 0)
}

func graphTest_assertSetEquals[T GraphNode[string, string]](t *testing.T, actual []T, expected ...string) {
	t.Helper()
	got := jutil.NewLinkedHashSet[string]()
	for _, n := range actual {
		got.Add(n.GetValue())
	}
	want := jutil.NewLinkedHashSet[string]()
	want.AddAll(expected...)
	assertSameElements(t, "values", got.ToSlice(), want.ToSlice())
}
