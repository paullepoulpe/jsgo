// Port of com/google/javascript/jscomp/graph/CheckPathsBetweenNodes.java @ closure-compiler v20260415

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

// CheckPathsBetweenNodes: see the constructor, NewCheckPathsBetweenNodes,
// for a description of this algorithm.
//
// N is the node type; E is the edge type.
type CheckPathsBetweenNodes[N, E any] struct {
	nodePredicate func(N) bool
	edgePredicate func(DiGraph_DiGraphEdge[N, E]) bool
	inclusive     bool

	graph DiGraph[N, E]
	start DiGraph_DiGraphNode[N, E]
	end   DiGraph_DiGraphNode[N, E]
}

// This algorithm works in two stages. First, the depth-first search (DFS)
// tree is calculated with A as the root. During when constructing the DFS
// tree, back edges are recorded. A back edge is a non-tree edge (X -> Y)
// where X is an descendant of Y in the DFS tree. The second step does a
// recursive traversal of the graph. Back edges are ignored during the
// recursive traversal, thus no cycles are encountered. Any recursive branch
// that encounters B without yet satisfying the predicate represents a path
// from the entry node to the exit without any nodes that satisfy the
// predicate.
//
// The implementation of discoverBackEdges follows the DFS-Visit algorithm
// in "Introduction to Algorithms" by Cormen, Leiseron, Rivest, and Stein,
// 2nd ed., on page 541. The calculation of back edges is described on page
// 546.

// checkPathsBetweenNodes_sentinel is the anonymous `new Annotation() {}`
// class. Instances are compared by identity; the field keeps the type
// non-zero-sized so that distinct sentinels have distinct addresses.
type checkPathsBetweenNodes_sentinel struct{ name string }

var (
	// A non-tree edge in the DFS that connects a node to one of its
	// ancestors.
	checkPathsBetweenNodes_BACK_EDGE    Annotation = &checkPathsBetweenNodes_sentinel{"BACK_EDGE"}
	checkPathsBetweenNodes_VISITED_EDGE Annotation = &checkPathsBetweenNodes_sentinel{"VISITED_EDGE"}

	// Not yet visited.
	checkPathsBetweenNodes_WHITE Annotation = nil
	// Being visited.
	checkPathsBetweenNodes_GRAY Annotation = &checkPathsBetweenNodes_sentinel{"GRAY"}
	// Finished visiting.
	checkPathsBetweenNodes_BLACK Annotation = &checkPathsBetweenNodes_sentinel{"BLACK"}
)

// NewCheckPathsBetweenNodes: given a graph G with nodes A and B, this
// algorithm determines if all paths from A to B contain at least one node
// satisfying a given predicate.
//
// Note that nodePredicate is not necessarily called for all nodes in G nor
// is edgePredicate called for all edges in G.
//
// graph is the graph G to analyze; a the node A; b the node B.
// nodePredicate is the predicate which at least one node on each path from
// an A node to B (inclusive) must match. edgePredicate selects the edges to
// consider as part of the graph: edges in graph that don't match
// edgePredicate will be ignored. inclusive includes node A and B in the
// test for the node predicate.
//
// Java: CheckPathsBetweenNodes(DiGraph<N, E>, DiGraphNode<N, E>, DiGraphNode<N, E>, Predicate<N>, Predicate<DiGraphEdge<N, E>>, boolean)
func NewCheckPathsBetweenNodes[N, E any](
	graph DiGraph[N, E],
	a DiGraph_DiGraphNode[N, E],
	b DiGraph_DiGraphNode[N, E],
	nodePredicate func(N) bool,
	edgePredicate func(DiGraph_DiGraphEdge[N, E]) bool,
	inclusive bool) *CheckPathsBetweenNodes[N, E] {
	return &CheckPathsBetweenNodes[N, E]{
		graph:         graph,
		start:         a,
		end:           b,
		nodePredicate: nodePredicate,
		edgePredicate: edgePredicate,
		inclusive:     inclusive,
	}
}

// NewCheckPathsBetweenNodesInclusive is the inclusive check.
//
// Java: CheckPathsBetweenNodes(DiGraph<N, E>, DiGraphNode<N, E>, DiGraphNode<N, E>, Predicate<N>, Predicate<DiGraphEdge<N, E>>)
func NewCheckPathsBetweenNodesInclusive[N, E any](
	graph DiGraph[N, E],
	a DiGraph_DiGraphNode[N, E],
	b DiGraph_DiGraphNode[N, E],
	nodePredicate func(N) bool,
	edgePredicate func(DiGraph_DiGraphEdge[N, E]) bool) *CheckPathsBetweenNodes[N, E] {
	return NewCheckPathsBetweenNodes(graph, a, b, nodePredicate, edgePredicate, true)
}

// AllPathsSatisfyPredicate returns true iff all paths contain at least one
// node that satisfy the predicate.
//
// Java: CheckPathsBetweenNodes.allPathsSatisfyPredicate()
func (c *CheckPathsBetweenNodes[N, E]) AllPathsSatisfyPredicate() bool {
	c.setUp()
	result := c.checkAllPathsWithoutBackEdges(c.start, c.end)
	c.tearDown()
	return result
}

// SomePathsSatisfyPredicate returns true iff some paths contain at least
// one node that satisfy the predicate.
//
// Java: CheckPathsBetweenNodes.somePathsSatisfyPredicate()
func (c *CheckPathsBetweenNodes[N, E]) SomePathsSatisfyPredicate() bool {
	c.setUp()
	result := c.checkSomePathsWithoutBackEdges(c.start, c.end)
	c.tearDown()
	return result
}

// Java: CheckPathsBetweenNodes.setUp()
func (c *CheckPathsBetweenNodes[N, E]) setUp() {
	c.graph.PushNodeAnnotations()
	c.graph.PushEdgeAnnotations()
	c.discoverBackEdges(c.start)
}

// Java: CheckPathsBetweenNodes.tearDown()
func (c *CheckPathsBetweenNodes[N, E]) tearDown() {
	c.graph.PopNodeAnnotations()
	c.graph.PopEdgeAnnotations()
}

// Java: CheckPathsBetweenNodes.discoverBackEdges(DiGraphNode<N, E>)
func (c *CheckPathsBetweenNodes[N, E]) discoverBackEdges(u DiGraph_DiGraphNode[N, E]) {
	u.SetAnnotation(checkPathsBetweenNodes_GRAY)
	for _, e := range u.GetOutEdges() {
		if c.ignoreEdge(e) {
			continue
		}
		v := e.GetDestination()
		if v.GetAnnotation() == checkPathsBetweenNodes_WHITE {
			c.discoverBackEdges(v)
		} else if v.GetAnnotation() == checkPathsBetweenNodes_GRAY {
			e.SetAnnotation(checkPathsBetweenNodes_BACK_EDGE)
		}
	}
	u.SetAnnotation(checkPathsBetweenNodes_BLACK)
}

// Java: CheckPathsBetweenNodes.ignoreEdge(DiGraphEdge<N, E>)
func (c *CheckPathsBetweenNodes[N, E]) ignoreEdge(e DiGraph_DiGraphEdge[N, E]) bool {
	return !c.edgePredicate(e)
}

// checkAllPathsWithoutBackEdges verifies that all non-looping paths from a
// to b pass through at least one node where nodePredicate is true.
//
// Java: CheckPathsBetweenNodes.checkAllPathsWithoutBackEdges(DiGraphNode<N, E>, DiGraphNode<N, E>)
func (c *CheckPathsBetweenNodes[N, E]) checkAllPathsWithoutBackEdges(a DiGraph_DiGraphNode[N, E], b DiGraph_DiGraphNode[N, E]) bool {
	if c.nodePredicate(a.GetValue()) &&
		(c.inclusive || (a != c.start && a != c.end)) {
		return true
	}
	if a == b {
		return false
	}
	for _, e := range a.GetOutEdges() {
		// Once we visited that edge once, we no longer need to
		// re-visit it again.
		if e.GetAnnotation() == checkPathsBetweenNodes_VISITED_EDGE {
			continue
		}
		e.SetAnnotation(checkPathsBetweenNodes_VISITED_EDGE)

		if c.ignoreEdge(e) {
			continue
		}
		if e.GetAnnotation() == checkPathsBetweenNodes_BACK_EDGE {
			continue
		}

		next := e.GetDestination()
		if !c.checkAllPathsWithoutBackEdges(next, b) {
			return false
		}
	}
	return true
}

// checkSomePathsWithoutBackEdges verifies that some non-looping paths from
// a to b pass through at least one node where nodePredicate is true.
//
// Java: CheckPathsBetweenNodes.checkSomePathsWithoutBackEdges(DiGraphNode<N, E>, DiGraphNode<N, E>)
func (c *CheckPathsBetweenNodes[N, E]) checkSomePathsWithoutBackEdges(a DiGraph_DiGraphNode[N, E], b DiGraph_DiGraphNode[N, E]) bool {
	if c.nodePredicate(a.GetValue()) &&
		(c.inclusive || (a != c.start && a != c.end)) {
		return true
	}
	if a == b {
		return false
	}
	for _, e := range a.GetOutEdges() {
		// Once we visited that edge once, we no longer need to
		// re-visit it again.
		if e.GetAnnotation() == checkPathsBetweenNodes_VISITED_EDGE {
			continue
		}
		e.SetAnnotation(checkPathsBetweenNodes_VISITED_EDGE)

		if c.ignoreEdge(e) {
			continue
		}
		if e.GetAnnotation() == checkPathsBetweenNodes_BACK_EDGE {
			continue
		}

		next := e.GetDestination()
		if c.checkSomePathsWithoutBackEdges(next, b) {
			return true
		}
	}
	return false
}
