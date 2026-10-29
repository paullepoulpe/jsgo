// Port of com/google/javascript/jscomp/graph/CheckPathsBetweenNodes.java @ closure-compiler v20261021

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

// CheckPathsBetweenNodes: see the constructor, NewCheckPathsBetweenNodes,
// for a description of this algorithm.
//
// N is the node type; E is the edge type.
//
// Go-only: performance — instead of mutating graph node/edge annotations via
// PushNodeAnnotations/PushEdgeAnnotations on every query, visited edges are
// tracked in an epoch-stamped uint32 array indexed by LinkedDiGraphEdge.Index,
// and traversal walks concrete LinkedDiGraphNodes directly.
type CheckPathsBetweenNodes[N, E comparable] struct {
	nodePredicate func(N) bool
	edgePredicate func(*LinkedDiGraphEdge[N, E]) bool
	allEdges      bool
	inclusive     bool

	linkedGraph *LinkedDirectedGraph[N, E]
	linkedStart *LinkedDiGraphNode[N, E]
	linkedEnd   *LinkedDiGraphNode[N, E]

	epoch     uint32
	edgeState []uint32
	edgeBuf   [64]uint32
}

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
func NewCheckPathsBetweenNodes[N, E comparable](
	graph DiGraph[N, E],
	a *LinkedDiGraphNode[N, E],
	b *LinkedDiGraphNode[N, E],
	nodePredicate func(N) bool,
	edgePredicate func(*LinkedDiGraphEdge[N, E]) bool,
	inclusive bool,
) *CheckPathsBetweenNodes[N, E] {
	c := &CheckPathsBetweenNodes[N, E]{}
	c.Reset(graph, a, b, nodePredicate, edgePredicate, inclusive)
	return c
}

// NewCheckPathsBetweenNodesInclusive is the inclusive check.
//
// Java: CheckPathsBetweenNodes(DiGraph<N, E>, DiGraphNode<N, E>, DiGraphNode<N, E>, Predicate<N>, Predicate<DiGraphEdge<N, E>>)
func NewCheckPathsBetweenNodesInclusive[N, E comparable](
	graph DiGraph[N, E],
	a *LinkedDiGraphNode[N, E],
	b *LinkedDiGraphNode[N, E],
	nodePredicate func(N) bool,
	edgePredicate func(*LinkedDiGraphEdge[N, E]) bool,
) *CheckPathsBetweenNodes[N, E] {
	return NewCheckPathsBetweenNodes(graph, a, b, nodePredicate, edgePredicate, true)
}

// NewCheckPathsBetweenNodesInclusiveValue returns a stack-allocatable CheckPathsBetweenNodes value.
func NewCheckPathsBetweenNodesInclusiveValue[N, E comparable](
	graph DiGraph[N, E],
	a *LinkedDiGraphNode[N, E],
	b *LinkedDiGraphNode[N, E],
	nodePredicate func(N) bool,
	edgePredicate func(*LinkedDiGraphEdge[N, E]) bool,
) CheckPathsBetweenNodes[N, E] {
	var c CheckPathsBetweenNodes[N, E]
	c.Reset(graph, a, b, nodePredicate, edgePredicate, true)
	return c
}

// Reset configures c for a new query while reusing its edge-visited state buffer.
func (c *CheckPathsBetweenNodes[N, E]) Reset(
	graph DiGraph[N, E],
	a *LinkedDiGraphNode[N, E],
	b *LinkedDiGraphNode[N, E],
	nodePredicate func(N) bool,
	edgePredicate func(*LinkedDiGraphEdge[N, E]) bool,
	inclusive bool,
) {
	c.linkedGraph = graph.AsLinkedDirectedGraph()
	c.linkedStart = a
	c.linkedEnd = b
	c.nodePredicate = nodePredicate
	c.edgePredicate = edgePredicate
	c.allEdges = false
	c.inclusive = inclusive
}

// ResetLinkedAllEdges configures c for a query on a concrete LinkedDirectedGraph
// where all edges are traversed without calling an edge predicate.
func (c *CheckPathsBetweenNodes[N, E]) ResetLinkedAllEdges(
	graph *LinkedDirectedGraph[N, E],
	a *LinkedDiGraphNode[N, E],
	b *LinkedDiGraphNode[N, E],
	nodePredicate func(N) bool,
	inclusive bool,
) {
	c.linkedGraph = graph
	c.linkedStart = a
	c.linkedEnd = b
	c.nodePredicate = nodePredicate
	c.edgePredicate = nil
	c.allEdges = true
	c.inclusive = inclusive
}

func (c *CheckPathsBetweenNodes[N, E]) prepare() uint32 {
	edgeCap := c.linkedGraph.EdgeIndexCap()
	c.epoch++
	if c.epoch == 0 {
		clear(c.edgeState)
		clear(c.edgeBuf[:])
		c.epoch = 1
	}
	if cap(c.edgeState) == 0 && edgeCap <= len(c.edgeBuf) {
		c.edgeState = c.edgeBuf[:]
	} else if cap(c.edgeState) < edgeCap {
		c.edgeState = make([]uint32, max(edgeCap, 128))
	} else if len(c.edgeState) < edgeCap {
		c.edgeState = c.edgeState[:edgeCap]
	}
	return c.epoch
}

// AllPathsSatisfyPredicate returns true iff all paths contain at least one
// node that satisfy the predicate.
//
// Java: CheckPathsBetweenNodes.allPathsSatisfyPredicate()
func (c *CheckPathsBetweenNodes[N, E]) AllPathsSatisfyPredicate() bool {
	epoch := c.prepare()
	return c.checkAllPaths(c.linkedStart, c.linkedEnd, epoch)
}

// SomePathsSatisfyPredicate returns true iff some paths contain at least
// one node that satisfy the predicate.
//
// Java: CheckPathsBetweenNodes.somePathsSatisfyPredicate()
func (c *CheckPathsBetweenNodes[N, E]) SomePathsSatisfyPredicate() bool {
	epoch := c.prepare()
	return c.checkSomePaths(c.linkedStart, c.linkedEnd, epoch)
}

func (c *CheckPathsBetweenNodes[N, E]) checkAllPaths(
	a *LinkedDiGraphNode[N, E],
	b *LinkedDiGraphNode[N, E],
	epoch uint32,
) bool {
	if (c.inclusive || (a != c.linkedStart && a != c.linkedEnd)) && c.nodePredicate(a.Value) {
		return true
	}
	if a == b {
		return false
	}
	for _, e := range a.OutEdgeList.Slice() {
		if c.edgeState[e.Index] == epoch {
			continue
		}
		c.edgeState[e.Index] = epoch
		if !c.allEdges && !c.edgePredicate(e) {
			continue
		}
		if !c.checkAllPaths(e.DestNode, b, epoch) {
			return false
		}
	}
	return true
}

func (c *CheckPathsBetweenNodes[N, E]) checkSomePaths(
	a *LinkedDiGraphNode[N, E],
	b *LinkedDiGraphNode[N, E],
	epoch uint32,
) bool {
	if (c.inclusive || (a != c.linkedStart && a != c.linkedEnd)) && c.nodePredicate(a.Value) {
		return true
	}
	if a == b {
		return false
	}
	for _, e := range a.OutEdgeList.Slice() {
		if c.edgeState[e.Index] == epoch {
			continue
		}
		c.edgeState[e.Index] = epoch
		if !c.allEdges && !c.edgePredicate(e) {
			continue
		}
		if c.checkSomePaths(e.DestNode, b, epoch) {
			return true
		}
	}
	return false
}
