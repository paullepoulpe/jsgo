// Port of com/google/javascript/jscomp/graph/GraphReachability.java @ closure-compiler v20261015

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// GraphReachability computes all the reachable nodes. Upon execution of
// Compute, the graph nodes will be annotated with REACHABLE if it is
// reachable from the specified entry node.
//
// N is the type of data that the graph node holds; E is the type of data
// that the graph edge holds.
//
// See GraphNode.GetAnnotation.
type GraphReachability[N comparable, E any] struct {
	// TODO(user): This should work for undirected graphs when
	// FixedPointGraphTraversal accepts them.
	graph DiGraph[N, E]

	edgePredicate func(*GraphReachability_EdgeTuple[N, E]) bool
	edgeTuple     GraphReachability_EdgeTuple[N, E]
	queue         []DiGraph_DiGraphNode[N, E]
}

// Java: GraphReachability(DiGraph<N, E>)
func NewGraphReachability[N comparable, E any](graph DiGraph[N, E]) *GraphReachability[N, E] {
	return NewGraphReachabilityWithPredicate[N, E](graph, nil)
}

// NewGraphReachabilityWithPredicate: given the predecessor P of the a node
// S and the edge coming from P to S, edgePredicate should return true if S
// is reachable from P using the edge. It may be nil.
//
// Java: GraphReachability(DiGraph<N, E>, Predicate<EdgeTuple<N, E>>)
func NewGraphReachabilityWithPredicate[N comparable, E any](
	graph DiGraph[N, E], edgePredicate func(*GraphReachability_EdgeTuple[N, E]) bool) *GraphReachability[N, E] {
	return &GraphReachability[N, E]{
		graph:         graph,
		edgePredicate: edgePredicate,
	}
}

// getNode is graph.getNode(value), dereferenced: a missing value is a
// NullPointerException at the Java call sites.
func (r *GraphReachability[N, E]) getNode(value N) DiGraph_DiGraphNode[N, E] {
	n := r.graph.GetDiGraphNode(value)
	if n == nil {
		panic(jutil.NewNullPointerException())
	}
	return n
}

func (r *GraphReachability[N, E]) propagateFrom(startNode DiGraph_DiGraphNode[N, E]) {
	var queue []DiGraph_DiGraphNode[N, E]
	if cap(r.queue) > 0 {
		queue = append(r.queue[:0], startNode)
	} else {
		queue = make([]DiGraph_DiGraphNode[N, E], 1, 32)
		queue[0] = startNode
	}
	pred := r.edgePredicate
	for head := 0; head < len(queue); head++ {
		node := queue[head]
		var sourceValue N
		if pred != nil {
			sourceValue = node.GetValue()
		}
		for _, edge := range node.GetOutEdges() {
			destNode := edge.GetDestination()
			if destNode.GetAnnotation() == GraphReachability_REACHABLE {
				continue
			}
			if pred != nil {
				r.edgeTuple.sourceNode = sourceValue
				r.edgeTuple.edge = edge.GetValue()
				if !pred(&r.edgeTuple) {
					continue
				}
			}
			destNode.SetAnnotation(GraphReachability_REACHABLE)
			queue = append(queue, destNode)
		}
	}
	clear(queue)
	if cap(queue) <= 8192 {
		r.queue = queue[:0]
	} else {
		r.queue = nil
	}
}

// Java: GraphReachability.compute(N)
func (r *GraphReachability[N, E]) Compute(entry N) {
	r.graph.ClearNodeAnnotations()
	entryNode := r.getNode(entry)
	entryNode.SetAnnotation(GraphReachability_REACHABLE)
	r.propagateFrom(entryNode)
}

// Java: GraphReachability.recompute(N)
func (r *GraphReachability[N, E]) Recompute(reachableNode N) {
	newReachable := r.getNode(reachableNode)
	jutil.CheckState(newReachable.GetAnnotation() != GraphReachability_REACHABLE)
	newReachable.SetAnnotation(GraphReachability_REACHABLE)
	r.propagateFrom(newReachable)
}

// Java: GraphReachability.traverseEdge(N, E, N)
func (r *GraphReachability[N, E]) TraverseEdge(source N, e E, destination N) bool {
	if r.getNode(source).GetAnnotation() == GraphReachability_REACHABLE {
		if r.edgePredicate != nil {
			r.edgeTuple.sourceNode = source
			r.edgeTuple.edge = e
			if !r.edgePredicate(&r.edgeTuple) {
				return false
			}
		}
		destNode := r.getNode(destination)
		if destNode.GetAnnotation() != GraphReachability_REACHABLE {
			destNode.SetAnnotation(GraphReachability_REACHABLE)
			return true
		}
	}
	return false
}

// graphReachability_sentinel is the anonymous `new Annotation() {}` class
// (compared by identity; non-zero-sized so its address is unique).
type graphReachability_sentinel struct{ name string }

// Java: GraphReachability.REACHABLE
var GraphReachability_REACHABLE Annotation = &graphReachability_sentinel{"REACHABLE"}

// GraphReachability_EdgeTuple represents a Source Node and an Edge.
type GraphReachability_EdgeTuple[N, E any] struct {
	sourceNode N
	edge       E
}

// NewGraphReachability_EdgeTuple: destNode is accepted but not stored, as
// in Java.
//
// Java: GraphReachability.EdgeTuple(N, E, N)
func NewGraphReachability_EdgeTuple[N, E any](sourceNode N, edge E, destNode N) *GraphReachability_EdgeTuple[N, E] {
	return &GraphReachability_EdgeTuple[N, E]{
		sourceNode: sourceNode,
		edge:       edge,
	}
}

var _ FixedPointGraphTraversal_EdgeCallback[string, string] = (*GraphReachability[string, string])(nil)
