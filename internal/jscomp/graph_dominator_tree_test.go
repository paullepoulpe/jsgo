// Port of com/google/javascript/jscomp/graph/DominatorTreeTest.java @ closure-compiler v20261015

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
)

func dominatorTreeTest_assertIdom(t *testing.T, tree *DominatorTree[string], node string, want string) {
	t.Helper()
	got, ok := tree.GetImmediateDominator(node)
	if !ok || got != want {
		t.Errorf("idom(%s): got %q (present=%v), want %q", node, got, ok, want)
	}
}

func dominatorTreeTest_assertIdomNull(t *testing.T, tree *DominatorTree[string], node string) {
	t.Helper()
	if got, ok := tree.GetImmediateDominator(node); ok {
		t.Errorf("idom(%s): got %q, want null", node, got)
	}
}

// dominatorTreeTest_assertSizes is assertThat(map).containsExactly(k, v, ...).
func dominatorTreeTest_assertSizes(t *testing.T, tree *DominatorTree[string], kv ...any) {
	t.Helper()
	got := tree.GetAllSubtreeSizes()
	if got.Size() != len(kv)/2 {
		t.Errorf("getAllSubtreeSizes: got %v, want %v", got, kv)
		return
	}
	for i := 0; i < len(kv); i += 2 {
		if v, ok := got.Lookup(kv[i].(string)); !ok || v != kv[i+1].(int) {
			t.Errorf("getAllSubtreeSizes: got %v, want %v", got, kv)
			return
		}
	}
}

func TestDominatorTree_SimpleChain(t *testing.T) {
	// A -> B -> C
	var graph DiGraph[string, string] = LinkedDirectedGraph_Create[string, string]()
	graph.CreateNode("A")
	graph.CreateNode("B")
	graph.CreateNode("C")
	graph.Connect("A", "->", "B")
	graph.Connect("B", "->", "C")

	tree := DominatorTree_Compute(graph, "A")

	dominatorTreeTest_assertIdomNull(t, tree, "A")
	dominatorTreeTest_assertIdom(t, tree, "B", "A")
	dominatorTreeTest_assertIdom(t, tree, "C", "B")

	assertEqual(t, "size(A)", tree.GetDominatedSubtreeSize("A"), 3)
	assertEqual(t, "size(B)", tree.GetDominatedSubtreeSize("B"), 2)
	assertEqual(t, "size(C)", tree.GetDominatedSubtreeSize("C"), 1)

	dominatorTreeTest_assertSizes(t, tree, "A", 3, "B", 2, "C", 1)
}

func TestDominatorTree_Diamond(t *testing.T) {
	// A -> {B, C} -> D
	var graph DiGraph[string, string] = LinkedDirectedGraph_Create[string, string]()
	graph.CreateNode("A")
	graph.CreateNode("B")
	graph.CreateNode("C")
	graph.CreateNode("D")
	graph.Connect("A", "->", "B")
	graph.Connect("A", "->", "C")
	graph.Connect("B", "->", "D")
	graph.Connect("C", "->", "D")

	tree := DominatorTree_Compute(graph, "A")

	dominatorTreeTest_assertIdom(t, tree, "B", "A")
	dominatorTreeTest_assertIdom(t, tree, "C", "A")
	dominatorTreeTest_assertIdom(t, tree, "D", "A")

	assertEqual(t, "size(A)", tree.GetDominatedSubtreeSize("A"), 4)
	assertEqual(t, "size(B)", tree.GetDominatedSubtreeSize("B"), 1)
	assertEqual(t, "size(C)", tree.GetDominatedSubtreeSize("C"), 1)
	assertEqual(t, "size(D)", tree.GetDominatedSubtreeSize("D"), 1)
}

func TestDominatorTree_Cycle(t *testing.T) {
	// A -> B -> C -> B
	// Arbitrarily chooses 'B' as the dominator of both 'B' and 'C'.
	// NOTE: in practice, we added DominatorTree for analyzing dependencies in
	// JSCompiler. Cycles between goog.require/goog.requireTyped files can
	// only occur within a single library, where it is valid to have
	// goog.requireType cycles.
	// So this result is not entirely precise, but is still useful for
	// practical purposes.
	var graph DiGraph[string, string] = LinkedDirectedGraph_Create[string, string]()
	graph.CreateNode("A")
	graph.CreateNode("B")
	graph.CreateNode("C")
	graph.Connect("A", "->", "B")
	graph.Connect("B", "->", "C")
	graph.Connect("C", "->", "B")

	tree := DominatorTree_Compute(graph, "A")

	dominatorTreeTest_assertIdom(t, tree, "B", "A")
	dominatorTreeTest_assertIdom(t, tree, "C", "B")

	assertEqual(t, "size(A)", tree.GetDominatedSubtreeSize("A"), 3)
	assertEqual(t, "size(B)", tree.GetDominatedSubtreeSize("B"), 2)
	assertEqual(t, "size(C)", tree.GetDominatedSubtreeSize("C"), 1)
}

func TestDominatorTree_Bottleneck(t *testing.T) {
	// Entry -> A -> B -> {C, D, E}
	var graph DiGraph[string, string] = LinkedDirectedGraph_Create[string, string]()
	graph.CreateNode("Entry")
	graph.CreateNode("A")
	graph.CreateNode("B")
	graph.CreateNode("C")
	graph.CreateNode("D")
	graph.CreateNode("E")
	graph.Connect("Entry", "->", "A")
	graph.Connect("A", "->", "B")
	graph.Connect("B", "->", "C")
	graph.Connect("B", "->", "D")
	graph.Connect("B", "->", "E")

	tree := DominatorTree_Compute(graph, "Entry")

	dominatorTreeTest_assertIdom(t, tree, "B", "A")
	assertEqual(t, "size(B)", tree.GetDominatedSubtreeSize("B"), 4) // B, C, D, E
	assertEqual(t, "size(A)", tree.GetDominatedSubtreeSize("A"), 5) // A, B, C, D, E
}

func TestDominatorTree_SharedBottleneck(t *testing.T) {
	// Entry -> {R1, R2}
	// R1 -> B
	// R2 -> B
	// B -> {C, D}
	var graph DiGraph[string, string] = LinkedDirectedGraph_Create[string, string]()
	graph.CreateNode("Entry")
	graph.CreateNode("R1")
	graph.CreateNode("R2")
	graph.CreateNode("B")
	graph.CreateNode("C")
	graph.CreateNode("D")
	graph.Connect("Entry", "->", "R1")
	graph.Connect("Entry", "->", "R2")
	graph.Connect("R1", "->", "B")
	graph.Connect("R2", "->", "B")
	graph.Connect("B", "->", "C")
	graph.Connect("B", "->", "D")

	tree := DominatorTree_Compute(graph, "Entry")

	dominatorTreeTest_assertIdom(t, tree, "B", "Entry")
	// Neither R1 nor R2 dominate "B", as there is always a path from Entry to
	// B going through the other.
	assertEqual(t, "size(R1)", tree.GetDominatedSubtreeSize("R1"), 1) // R1
	assertEqual(t, "size(R2)", tree.GetDominatedSubtreeSize("R2"), 1) // R2
	assertEqual(t, "size(B)", tree.GetDominatedSubtreeSize("B"), 3)   // B, C, D
}

func TestDominatorTree_NodeNotReachableFromEntryPoint_doesNotCountAsDominator(t *testing.T) {
	// A -> B
	// OTHER -> B
	var graph DiGraph[string, string] = LinkedDirectedGraph_Create[string, string]()
	graph.CreateNode("A")
	graph.CreateNode("B")
	graph.CreateNode("OTHER")
	graph.Connect("A", "->", "B")
	graph.Connect("OTHER", "->", "B")

	tree := DominatorTree_Compute(graph, "A")

	// The dominator tree does not consider "OTHER" to dominate "B", as it is
	// only looking at paths *from the specified entry point* "A" to "B".
	dominatorTreeTest_assertIdom(t, tree, "B", "A")
	dominatorTreeTest_assertIdomNull(t, tree, "OTHER")
	assertEqual(t, "size(OTHER)", tree.GetDominatedSubtreeSize("OTHER"), 0)

	dominatorTreeTest_assertSizes(t, tree, "A", 2, "B", 1)
}

func TestDominatorTree_GetImmediateDominator_nonExistentElement_returnsNull(t *testing.T) {
	var graph DiGraph[string, string] = LinkedDirectedGraph_Create[string, string]()
	graph.CreateNode("A")

	tree := DominatorTree_Compute(graph, "A")

	dominatorTreeTest_assertIdomNull(t, tree, "FOO")
}

func TestDominatorTree_GetSubtreeSize_nonExistentElement_defaultsToZero(t *testing.T) {
	var graph DiGraph[string, string] = LinkedDirectedGraph_Create[string, string]()
	graph.CreateNode("A")

	tree := DominatorTree_Compute(graph, "A")

	assertEqual(t, "size(FOO)", tree.GetDominatedSubtreeSize("FOO"), 0)
}
