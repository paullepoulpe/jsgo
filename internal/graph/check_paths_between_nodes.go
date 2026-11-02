// Port of com/google/javascript/jscomp/graph/CheckPathsBetweenNodes.java @ closure-compiler v20261028

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
// PushNodeAnnotations/PushEdgeAnnotations on every query, DFS node colors and
// back/visited edges are tracked in epoch-stamped uint32 arrays indexed by
// LinkedDiGraphNode.Index and LinkedDiGraphEdge.Index, and traversal walks
// concrete LinkedDiGraphNodes directly.
type CheckPathsBetweenNodes[N, E comparable] struct {
	nodePredicate func(N) bool
	edgePredicate func(*LinkedDiGraphEdge[N, E]) bool
	allEdges      bool
	inclusive     bool

	linkedGraph *LinkedDirectedGraph[N, E]
	linkedStart *LinkedDiGraphNode[N, E]
	linkedEnd   *LinkedDiGraphNode[N, E]

	epoch     uint32
	nodeState []uint32
	nodeBuf   [64]uint32
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
	nodeCap := c.linkedGraph.NodeIndexCap()
	edgeCap := c.linkedGraph.EdgeIndexCap()
	c.epoch += 2
	if c.epoch < 2 {
		clear(c.nodeState)
		clear(c.nodeBuf[:])
		clear(c.edgeState)
		clear(c.edgeBuf[:])
		c.epoch = 2
	}
	if cap(c.nodeState) == 0 && nodeCap <= len(c.nodeBuf) {
		c.nodeState = c.nodeBuf[:]
	} else if cap(c.nodeState) < nodeCap {
		c.nodeState = make([]uint32, max(nodeCap, 128))
	} else if len(c.nodeState) < nodeCap {
		c.nodeState = c.nodeState[:nodeCap]
	}
	if cap(c.edgeState) == 0 && edgeCap <= len(c.edgeBuf) {
		c.edgeState = c.edgeBuf[:]
	} else if cap(c.edgeState) < edgeCap {
		c.edgeState = make([]uint32, max(edgeCap, 128))
	} else if len(c.edgeState) < edgeCap {
		c.edgeState = c.edgeState[:edgeCap]
	}
	base := c.epoch
	c.discoverBackEdges(c.linkedStart, base)
	return base
}

func (c *CheckPathsBetweenNodes[N, E]) discoverBackEdges(u *LinkedDiGraphNode[N, E], base uint32) {
	c.nodeState[u.Index] = base // GRAY
	for _, e := range u.OutEdgeList.Slice() {
		if !c.allEdges && !c.edgePredicate(e) {
			continue
		}
		v := e.DestNode
		vs := c.nodeState[v.Index]
		if vs < base { // WHITE
			c.discoverBackEdges(v, base)
		} else if vs == base { // GRAY
			c.edgeState[e.Index] = base // BACK_EDGE
		}
	}
	c.nodeState[u.Index] = base + 1 // BLACK
}

// AllPathsSatisfyPredicate returns true iff all paths contain at least one
// node that satisfy the predicate.
//
// Java: CheckPathsBetweenNodes.allPathsSatisfyPredicate()
func (c *CheckPathsBetweenNodes[N, E]) AllPathsSatisfyPredicate() bool {
	base := c.prepare()
	return c.checkAllPaths(c.linkedStart, c.linkedEnd, base)
}

// SomePathsSatisfyPredicate returns true iff some paths contain at least
// one node that satisfy the predicate.
//
// Java: CheckPathsBetweenNodes.somePathsSatisfyPredicate()
func (c *CheckPathsBetweenNodes[N, E]) SomePathsSatisfyPredicate() bool {
	base := c.prepare()
	return c.checkSomePaths(c.linkedStart, c.linkedEnd, base)
}

func (c *CheckPathsBetweenNodes[N, E]) checkAllPaths(
	a *LinkedDiGraphNode[N, E],
	b *LinkedDiGraphNode[N, E],
	base uint32,
) bool {
	if (c.inclusive || (a != c.linkedStart && a != c.linkedEnd)) && c.nodePredicate(a.Value) {
		return true
	}
	if a == b {
		return false
	}
	visited := base + 1
	for _, e := range a.OutEdgeList.Slice() {
		es := c.edgeState[e.Index]
		if es == visited {
			continue
		}
		if !c.allEdges && !c.edgePredicate(e) {
			continue
		}
		if es == base { // BACK_EDGE
			continue
		}
		c.edgeState[e.Index] = visited
		if !c.checkAllPaths(e.DestNode, b, base) {
			return false
		}
	}
	return true
}

func (c *CheckPathsBetweenNodes[N, E]) checkSomePaths(
	a *LinkedDiGraphNode[N, E],
	b *LinkedDiGraphNode[N, E],
	base uint32,
) bool {
	if (c.inclusive || (a != c.linkedStart && a != c.linkedEnd)) && c.nodePredicate(a.Value) {
		return true
	}
	if a == b {
		return false
	}
	visited := base + 1
	for _, e := range a.OutEdgeList.Slice() {
		es := c.edgeState[e.Index]
		if es == visited {
			continue
		}
		if !c.allEdges && !c.edgePredicate(e) {
			continue
		}
		if es == base { // BACK_EDGE
			continue
		}
		c.edgeState[e.Index] = visited
		if c.checkSomePaths(e.DestNode, b, base) {
			return true
		}
	}
	return false
}
