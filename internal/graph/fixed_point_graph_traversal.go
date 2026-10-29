// Port of com/google/javascript/jscomp/graph/FixedPointGraphTraversal.java @ closure-compiler v20261021

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
	"math"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// FixedPointGraphTraversal is a utility class for doing fixed-point
// computations. We traverse the edges over the given directed graph until
// the graph reaches a steady state.
//
// N is the value type that the graph node stores; E is the value type that
// the graph edge stores.
//
// TODO(nicksantos): Generalize the algorithm for undirected graphs, if we
// need it.
type FixedPointGraphTraversal[N comparable, E any] struct {
	callback FixedPointGraphTraversal_EdgeCallback[N, E]

	traversalDirection fixedPointGraphTraversal_TraversalDirection
}

// fixedPointGraphTraversal_TraversalDirection is the private enum
// TraversalDirection.
type fixedPointGraphTraversal_TraversalDirection int32

const (
	fixedPointGraphTraversal_TraversalDirection_INWARDS  fixedPointGraphTraversal_TraversalDirection = iota // from a node to its incoming edges
	fixedPointGraphTraversal_TraversalDirection_OUTWARDS                                                    // from a node to its outgoing edges
)

func (d fixedPointGraphTraversal_TraversalDirection) String() string {
	switch d {
	case fixedPointGraphTraversal_TraversalDirection_INWARDS:
		return "INWARDS"
	case fixedPointGraphTraversal_TraversalDirection_OUTWARDS:
		return "OUTWARDS"
	}
	return "?"
}

// Java: FixedPointGraphTraversal.NON_HALTING_ERROR_MSG
const FixedPointGraphTraversal_NON_HALTING_ERROR_MSG = "Fixed point computation not halting"

// newFixedPointGraphTraversal creates a new traversal. callback is a
// callback for updating the state of the graph each time an edge is
// traversed.
//
// Java: FixedPointGraphTraversal(EdgeCallback<N, E>, TraversalDirection)
func newFixedPointGraphTraversal[N comparable, E any](
	callback FixedPointGraphTraversal_EdgeCallback[N, E], traversalDirection fixedPointGraphTraversal_TraversalDirection) *FixedPointGraphTraversal[N, E] {
	return &FixedPointGraphTraversal[N, E]{
		callback:           callback,
		traversalDirection: traversalDirection,
	}
}

// FixedPointGraphTraversal_NewTraversal is a helper method for creating new
// traversals that traverse from parent to child.
//
// Java: FixedPointGraphTraversal.newTraversal(EdgeCallback<NodeT, EdgeT>)
func FixedPointGraphTraversal_NewTraversal[NodeT comparable, EdgeT any](
	callback FixedPointGraphTraversal_EdgeCallback[NodeT, EdgeT]) *FixedPointGraphTraversal[NodeT, EdgeT] {
	return newFixedPointGraphTraversal(callback, fixedPointGraphTraversal_TraversalDirection_OUTWARDS)
}

// FixedPointGraphTraversal_NewReverseTraversal is a helper method for
// creating new traversals that traverse from child to parent.
//
// Java: FixedPointGraphTraversal.newReverseTraversal(EdgeCallback<NodeT, EdgeT>)
func FixedPointGraphTraversal_NewReverseTraversal[NodeT comparable, EdgeT any](
	callback FixedPointGraphTraversal_EdgeCallback[NodeT, EdgeT]) *FixedPointGraphTraversal[NodeT, EdgeT] {
	return newFixedPointGraphTraversal(callback, fixedPointGraphTraversal_TraversalDirection_INWARDS)
}

// ComputeFixedPoint computes a fixed point for the given graph.
//
// Java: FixedPointGraphTraversal.computeFixedPoint(DiGraph<N, E>)
func (t *FixedPointGraphTraversal[N, E]) ComputeFixedPoint(graph DiGraph[N, E]) {
	var zeroN N
	t.computeFixedPointLinked(graph.AsLinkedDirectedGraph(), true, zeroN, nil)
}

// ComputeFixedPointFrom computes a fixed point for the given graph,
// entering from the given node. entry is the node to begin traversing
// from.
//
// Java: FixedPointGraphTraversal.computeFixedPoint(DiGraph<N, E>, N)
func (t *FixedPointGraphTraversal[N, E]) ComputeFixedPointFrom(graph DiGraph[N, E], entry N) {
	t.computeFixedPointLinked(graph.AsLinkedDirectedGraph(), false, entry, nil)
}

// We cube the number of nodes to estimate the largest number of iterations
// we should allow before deciding that we aren't reaching a fixed point.
// We have to make sure that we don't hit integer overflow when calculating
// this value. We also don't want to wait longer than a full minute for any
// fixed point calculation. We'll be generous and assume each iteration
// takes only a nanosecond. That's 60*10^9 iterations, so we must cap the
// node count we use for calculation at the cube root of that number.
var fixedPointGraphTraversal_MAX_NODE_COUNT_FOR_ITERATION_LIMIT = int64(math.Floor(math.Cbrt(60e9)))

// ComputeFixedPointFromSet computes a fixed point for the given graph,
// entering from the given nodes. entrySet is the nodes to begin traversing
// from.
//
// Java: FixedPointGraphTraversal.computeFixedPoint(DiGraph<N, E>, Set<N>)
func (t *FixedPointGraphTraversal[N, E]) ComputeFixedPointFromSet(graph DiGraph[N, E], entrySet *jutil.LinkedHashSet[N]) {
	if entrySet == nil {
		panic(jutil.NewNullPointerException())
	}
	var zeroN N
	t.computeFixedPointLinked(graph.AsLinkedDirectedGraph(), false, zeroN, entrySet)
}

func (t *FixedPointGraphTraversal[N, E]) computeFixedPointLinked(
	ldg *LinkedDirectedGraph[N, E], initAll bool, singleEntry N, entrySet *jutil.LinkedHashSet[N]) {
	nodeCountInt := ldg.Nodes.Size()
	ringCap := 16
	for ringCap <= nodeCountInt+1 {
		ringCap <<= 1
	}
	var ringBuf [64]*LinkedDiGraphNode[N, E]
	var ring []*LinkedDiGraphNode[N, E]
	if ringCap <= len(ringBuf) {
		ring = ringBuf[:ringCap]
	} else {
		ring = make([]*LinkedDiGraphNode[N, E], ringCap)
	}
	var head, tail uint32
	mask := uint32(ringCap - 1)
	defer func() {
		for head != tail {
			if n := ring[head&mask]; n != nil {
				n.inWorkSet = false
			}
			head++
		}
	}()

	var nilInWorkSet bool
	if initAll {
		for node := range ldg.Nodes.Values() {
			if !node.inWorkSet {
				node.inWorkSet = true
				ring[tail&mask] = node
				tail++
			}
		}
	} else if entrySet != nil {
		for v := range entrySet.All() {
			node := ldg.GetLinkedDiGraphNode(v)
			if node == nil {
				if !nilInWorkSet {
					nilInWorkSet = true
					ring[tail&mask] = nil
					tail++
				}
				continue
			}
			if !node.inWorkSet {
				node.inWorkSet = true
				ring[tail&mask] = node
				tail++
			}
		}
	} else {
		node := ldg.GetLinkedDiGraphNode(singleEntry)
		if node == nil {
			panic(jutil.NewNullPointerException())
		}
		node.inWorkSet = true
		ring[tail&mask] = node
		tail++
	}

	cycleCount := int64(0)
	nodeCount := min(int64(nodeCountInt), fixedPointGraphTraversal_MAX_NODE_COUNT_FOR_ITERATION_LIMIT)
	maxIterations := max(nodeCount*nodeCount*nodeCount, 100)

	outwards := t.traversalDirection == fixedPointGraphTraversal_TraversalDirection_OUTWARDS
	inwards := t.traversalDirection == fixedPointGraphTraversal_TraversalDirection_INWARDS

	for ; head != tail && cycleCount < maxIterations; cycleCount++ {
		idx := head & mask
		node := ring[idx]
		ring[idx] = nil
		head++
		if node == nil {
			// An entry value that is not in the graph: node.getValue() on null.
			panic(jutil.NewNullPointerException())
		}
		node.inWorkSet = false

		if outwards {
			sourceValue := node.Value
			for _, edge := range node.OutEdgeList.Slice() {
				destNode := edge.DestNode
				if t.callback.TraverseEdge(sourceValue, edge.Value, destNode.Value) {
					if !destNode.inWorkSet {
						destNode.inWorkSet = true
						ring[tail&mask] = destNode
						tail++
					}
				}
			}
		} else if inwards {
			revSourceValue := node.Value
			for _, edge := range node.InEdgeList.Slice() {
				srcNode := edge.SourceNode
				if t.callback.TraverseEdge(revSourceValue, edge.Value, srcNode.Value) {
					if !srcNode.inWorkSet {
						srcNode.inWorkSet = true
						ring[tail&mask] = srcNode
						tail++
					}
				}
			}
		} else {
			panic(jutil.NewAssertionError("Unrecognized direction " + t.traversalDirection.String()))
		}
	}

	if cycleCount == maxIterations {
		panic(jutil.NewIllegalStateExceptionMsg(FixedPointGraphTraversal_NON_HALTING_ERROR_MSG))
	}
}

// FixedPointGraphTraversal_EdgeCallback is the edge callback.
type FixedPointGraphTraversal_EdgeCallback[NodeT, EdgeT any] interface {
	// TraverseEdge updates the state of the destination node when the given
	// edge is traversed.
	//
	// Recall that depending on the direction of the traversal, source and
	// destination may be swapped compared to the orientation of the edge in
	// the graph. In either case, only the destination parameter may be
	// mutated.
	//
	// source is the start node, e the edge, destination the end node.
	// Returns whether the state of the destination node changed.
	//
	// Java: FixedPointGraphTraversal.EdgeCallback.traverseEdge(NodeT, EdgeT, NodeT)
	TraverseEdge(source NodeT, e EdgeT, destination NodeT) bool
}

// FixedPointGraphTraversal_EdgeCallbackFunc lets a closure (a Java lambda)
// be used as an EdgeCallback.
//
// Go-only: language-adapter — Java lambdas implement the functional
// interface directly.
type FixedPointGraphTraversal_EdgeCallbackFunc[NodeT, EdgeT any] func(source NodeT, e EdgeT, destination NodeT) bool

// TraverseEdge calls f.
func (f FixedPointGraphTraversal_EdgeCallbackFunc[NodeT, EdgeT]) TraverseEdge(source NodeT, e EdgeT, destination NodeT) bool {
	return f(source, e, destination)
}
