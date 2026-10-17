// Port of com/google/javascript/jscomp/graph/GraphColoringTest.java @ closure-compiler v20261015

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
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// The Java tests use Graph<String, String> with null edge values; E is
// `any` here so that null is representable (nil).

func TestGraphColoring_NoEdge(t *testing.T) {
	var graph Graph[string, any] = LinkedUndirectedGraph_Create[string, any]()
	for i := 0; i < 5; i++ {
		graph.CreateNode("Node " + strconv.Itoa(i))
		// All node with same color.
		var coloring GraphColoring[string, any] = NewGraphColoring_GreedyGraphColoring[string, any](graph)
		assertEqual(t, "colors", coloring.Color(), 1)
		graphColoringTest_validateColoring(t, graph)
		for j := 0; j < i; j++ {
			assertEqual(t, "super node", coloring.GetPartitionSuperNode("Node 0"), "Node 0")
		}
	}
}

func TestGraphColoring_TwoNodesConnected(t *testing.T) {
	var graph Graph[string, any] = LinkedUndirectedGraph_Create[string, any]()
	graph.CreateNode("A")
	graph.CreateNode("B")
	graph.Connect("A", "--", "B")
	var coloring GraphColoring[string, any] = NewGraphColoring_GreedyGraphColoring[string, any](graph)
	assertEqual(t, "colors", coloring.Color(), 2)
	graphColoringTest_validateColoring(t, graph)
	assertEqual(t, "super(A)", coloring.GetPartitionSuperNode("A"), "A")
	assertEqual(t, "super(B)", coloring.GetPartitionSuperNode("B"), "B")
}

func TestGraphColoring_Greedy(t *testing.T) {
	var graph Graph[string, any] = LinkedUndirectedGraph_Create[string, any]()
	graph.CreateNode("A")
	graph.CreateNode("B")
	graph.CreateNode("C")
	graph.CreateNode("D")
	graph.Connect("A", "--", "C")
	graph.Connect("B", "--", "C")
	graph.Connect("B", "--", "D")
	var coloring GraphColoring[string, any] = NewGraphColoring_GreedyGraphColoring[string, any](graph)
	assertEqual(t, "colors", coloring.Color(), 2)
	graphColoringTest_validateColoring(t, graph)
	assertEqual(t, "super(A)", coloring.GetPartitionSuperNode("A"), "A")
	assertEqual(t, "super(B)", coloring.GetPartitionSuperNode("B"), "A")
	assertEqual(t, "super(C)", coloring.GetPartitionSuperNode("C"), "C")
}

func TestGraphColoring_FullyConnected(t *testing.T) {
	const count = 100
	var graph Graph[string, any] = LinkedUndirectedGraph_Create[string, any]()
	for i := 0; i < count; i++ {
		graph.CreateNode("Node " + strconv.Itoa(i))
		for j := 0; j < count; j++ {
			graph.CreateNode("Node " + strconv.Itoa(j))
			if i != j {
				graph.Connect("Node "+strconv.Itoa(i), nil, "Node "+strconv.Itoa(j))
			}
		}
	}
	var coloring GraphColoring[string, any] = NewGraphColoring_GreedyGraphColoring[string, any](graph)
	assertEqual(t, "colors", coloring.Color(), count)
	graphColoringTest_validateColoring(t, graph)
	for i := 0; i < count; i++ {
		assertEqual(t, "super node", coloring.GetPartitionSuperNode("Node "+strconv.Itoa(i)), "Node "+strconv.Itoa(i))
	}
}

func TestGraphColoring_AllConnectedToOneNode(t *testing.T) {
	const count = 10
	var graph Graph[string, any] = LinkedUndirectedGraph_Create[string, any]()
	graph.CreateNode("Center")
	for i := 0; i < count; i++ {
		graph.CreateNode("Node " + strconv.Itoa(i))
		graph.Connect("Center", nil, "Node "+strconv.Itoa(i))
	}
	var coloring GraphColoring[string, any] = NewGraphColoring_GreedyGraphColoring[string, any](graph)
	assertEqual(t, "colors", coloring.Color(), 2)
	graphColoringTest_validateColoring(t, graph)
	assertEqual(t, "super(Center)", coloring.GetPartitionSuperNode("Center"), "Center")
	for i := 0; i < count; i++ {
		assertEqual(t, "super node", coloring.GetPartitionSuperNode("Node "+strconv.Itoa(i)), "Node 0")
	}
}

func TestGraphColoring_TwoFullyConnected(t *testing.T) {
	const count = 100
	// A graph with two disconnected disjunct cliques.
	var graph Graph[string, any] = LinkedUndirectedGraph_Create[string, any]()
	for i := 0; i < count; i++ {
		graph.CreateNode("Node Left " + strconv.Itoa(i))
		graph.CreateNode("Node Right " + strconv.Itoa(i))
		for j := 0; j < count; j++ {
			graph.CreateNode("Node Left " + strconv.Itoa(j))
			graph.CreateNode("Node Right " + strconv.Itoa(j))
			if i != j {
				graph.Connect("Node Left "+strconv.Itoa(i), nil, "Node Left "+strconv.Itoa(j))
				graph.Connect("Node Right "+strconv.Itoa(i), nil, "Node Right "+strconv.Itoa(j))
			}
		}
	}
	assertEqual(t, "colors", NewGraphColoring_GreedyGraphColoring[string, any](graph).Color(), count)
	graphColoringTest_validateColoring(t, graph)

	// Connect the two cliques.
	for i := 0; i < count; i++ {
		graph.Connect("Node Left "+strconv.Itoa(i), nil, "Node Right "+strconv.Itoa(i))
	}
	// Think of two exactly same graph with the same coloring side by side.
	// If we circularly shift the colors of one of the graph by 1, we can
	// connect the isomorphic nodes and still have a valid coloring in the
	// resulting graph.
	assertEqual(t, "colors", NewGraphColoring_GreedyGraphColoring[string, any](graph).Color(), count)
	graphColoringTest_validateColoring(t, graph)
}

func TestGraphColoring_Deterministic(t *testing.T) {
	// A pentagon.
	var graph Graph[string, any] = LinkedUndirectedGraph_Create[string, any]()
	graph.CreateNode("A")
	graph.CreateNode("B")
	graph.CreateNode("C")
	graph.CreateNode("D")
	graph.CreateNode("E")
	graph.Connect("A", "-->", "B")
	graph.Connect("B", "-->", "C")
	graph.Connect("C", "-->", "D")
	graph.Connect("D", "-->", "E")
	graph.Connect("E", "-->", "A")

	var coloring GraphColoring[string, any] = NewGraphColoring_GreedyGraphColoringWithTieBreaker[string, any](graph, jutil.StringCompareTo)
	assertEqual(t, "colors", coloring.Color(), 3)
	graphColoringTest_validateColoring(t, graph)
	assertEqual(t, "super(A)", coloring.GetPartitionSuperNode("A"), "A")
	assertEqual(t, "super(C)", coloring.GetPartitionSuperNode("C"), "A")

	biasD := func(a, b string) int {
		return jutil.StringCompareTo(strings.ReplaceAll(a, "D", "@"), strings.ReplaceAll(b, "D", "@"))
	}

	coloring = NewGraphColoring_GreedyGraphColoringWithTieBreaker[string, any](graph, biasD)
	assertEqual(t, "colors", coloring.Color(), 3)
	graphColoringTest_validateColoring(t, graph)
	assertEqual(t, "super(A)", coloring.GetPartitionSuperNode("A"), "A")
	assertFalse(t, "super(C) is A", coloring.GetPartitionSuperNode("C") == "A")
}

// graphColoringTest_validateColoring validates that each node has been
// colored and connected nodes have different coloring.
func graphColoringTest_validateColoring[N, E any](t *testing.T, graph Graph[N, E]) {
	t.Helper()
	for _, node := range graph.GetNodes() {
		if node.GetAnnotation() == nil {
			t.Errorf("node %v has no color", node.GetValue())
		}
	}
	for _, edge := range graph.GetEdges() {
		c1, _ := edge.GetNodeA().GetAnnotation().(*GraphColoring_Color)
		c2, _ := edge.GetNodeB().GetAnnotation().(*GraphColoring_Color)
		if c1 == nil || c2 == nil {
			t.Errorf("edge %v has an uncolored node", edge)
			continue
		}
		if c1.Equals(c2) {
			t.Errorf("edge %v connects nodes of the same color", edge)
		}
	}
}
