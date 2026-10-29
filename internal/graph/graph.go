// Port of com/google/javascript/jscomp/graph/AdjacencyGraph.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/graph/Annotatable.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/graph/Annotation.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/graph/DiGraph.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/graph/Graph.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/graph/GraphNode.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/graph/GraphvizGraph.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/graph/LatticeElement.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/graph/StandardUnionFind.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/graph/SubGraph.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/graph/UndiGraph.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/graph/UnionFind.java @ closure-compiler v20261021

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

/*
 * Copyright 2010 The Closure Compiler Authors.
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

// Package graph provides generic directed and undirected graph data structures
// and algorithms used by control-flow, data-flow, dependency, and
// register/property coloring passes in jscomp.
//
// Core types and algorithms:
//   - [LinkedDirectedGraph] and [LinkedUndirectedGraph] — insertion-ordered
//     graphs backed by slab-allocated nodes and inline 2-edge lists
//   - [DiGraph], [Graph], [SubGraph], [GraphvizGraph]
//   - [FixedPointGraphTraversal], [GraphReachability], [CheckPathsBetweenNodes]
//   - [DominatorTree], [LowestCommonAncestorFinder]
//   - [GreedyGraphColoring]
package graph

import (
	"fmt"
	"iter"
	"reflect"
	"slices"
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// GraphNodeMap stores graph nodes keyed by value N in deterministic insertion order.
type GraphNodeMap[N comparable, NodePtr interface {
	comparable
	GetValue() N
}] struct {
	nodes     map[N]NodePtr
	nodeOrder []NodePtr
}

// Get returns the node for k, or the zero NodePtr (nil) if absent.
func (m *GraphNodeMap[N, NodePtr]) Get(k N) NodePtr {
	return m.nodes[k]
}

// Size returns the number of nodes in m.
func (m *GraphNodeMap[N, NodePtr]) Size() int {
	return len(m.nodeOrder)
}

// Slice returns the nodes in deterministic insertion order, aliasing m's backing slice.
func (m *GraphNodeMap[N, NodePtr]) Slice() []NodePtr {
	return m.nodeOrder
}

// PutNew inserts node for absent key k, appending it to nodeOrder.
func (m *GraphNodeMap[N, NodePtr]) PutNew(k N, node NodePtr) {
	if m.nodes == nil {
		m.nodes = make(map[N]NodePtr, 16)
	}
	m.nodes[k] = node
	m.nodeOrder = append(m.nodeOrder, node)
}

// Remove deletes k if present while preserving exact insertion order of remaining nodes.
func (m *GraphNodeMap[N, NodePtr]) Remove(k N) bool {
	node, ok := m.nodes[k]
	if !ok {
		return false
	}
	delete(m.nodes, k)
	if i := slices.Index(m.nodeOrder, node); i >= 0 {
		m.nodeOrder = slices.Delete(m.nodeOrder, i, i+1)
	}
	return true
}

// Clear removes all entries while keeping backing storage for reuse.
func (m *GraphNodeMap[N, NodePtr]) Clear() {
	clear(m.nodes)
	clear(m.nodeOrder)
	m.nodeOrder = m.nodeOrder[:0]
}

// Reset drops backing storage when a graph grew too large to retain for reuse.
func (m *GraphNodeMap[N, NodePtr]) Reset() {
	m.nodes = nil
	m.nodeOrder = nil
}

// Values iterates nodes in deterministic insertion order.
func (m *GraphNodeMap[N, NodePtr]) Values() iter.Seq[NodePtr] {
	return slices.Values(m.nodeOrder)
}

// All iterates (value, node) pairs in deterministic insertion order.
func (m *GraphNodeMap[N, NodePtr]) All() iter.Seq2[N, NodePtr] {
	return func(yield func(N, NodePtr) bool) {
		for _, node := range m.nodeOrder {
			if !yield(node.GetValue(), node) {
				return
			}
		}
	}
}

// Annotation is information that can be annotated to a graph node or edge.
type Annotation interface{}

// Annotatable is an object that carries an annotation.
type Annotatable interface {
	// SetAnnotation annotates a piece of information to the object.
	//
	// Java: Annotatable.setAnnotation(Annotation)
	SetAnnotation(data Annotation)

	// GetAnnotation retrieves a piece of information that has been annotated,
	// or nil if the object has not been annotated.
	//
	// Java: Annotatable.getAnnotation()
	GetAnnotation() Annotation
}

// LatticeElement is a lattice element used in data-flow annotations.
type LatticeElement interface {
	Annotation
}

// GraphNode is a generic node carrying a value and an optional annotation.
type GraphNode[N, E any] interface {
	Annotatable

	// GetValue retrieves the node's value.
	//
	// Java: GraphNode.getValue()
	GetValue() N
}

// SubGraph represents a subgraph that provides adjacency checks against a node.
type SubGraph[N, E any] interface {
	// IsIndependentOf returns true if node is not adjacent to any node in this SubGraph.
	//
	// Java: SubGraph.isIndependentOf(N)
	IsIndependentOf(node N) bool

	// AddNode adds value into this subgraph.
	//
	// Java: SubGraph.addNode(N)
	AddNode(value N)
}

// Graph is the common interface shared by [LinkedDirectedGraph] and
// [LinkedUndirectedGraph] for algorithms such as [GreedyGraphColoring].
type Graph[N comparable, E any] interface {
	Connect(n1 N, edge E, n2 N)
	Disconnect(n1 N, n2 N)
	ConnectIfNotFound(n1 N, edge E, n2 N)
	GetNodeCount() int
	GetNodeDegree(value N) int
	GetWeight(value N) int
	HasNode(n N) bool
	IsConnected(n1 N, n2 N) bool
	IsConnectedWithEdge(n1 N, e E, n2 N) bool
	NewSubGraph() SubGraph[N, E]
	ClearNodeAnnotations()
	ClearEdgeAnnotations()
	PushNodeAnnotations()
	PopNodeAnnotations()
	PushEdgeAnnotations()
	PopEdgeAnnotations()
}

// DiGraph is implemented by [*LinkedDirectedGraph] and types that embed it
// (such as jscomp.ControlFlowGraph).
type DiGraph[N comparable, E any] interface {
	AsLinkedDirectedGraph() *LinkedDirectedGraph[N, E]
}

// GraphvizGraph is a graph that can be dumped to a Graphviz DOT file.
type GraphvizGraph interface {
	// Java: GraphvizGraph.getName()
	GetName() string

	// Java: GraphvizGraph.isDirected()
	IsDirected() bool

	// Java: GraphvizGraph.getGraphvizNodes()
	GetGraphvizNodes() []GraphvizGraph_GraphvizNode

	// Java: GraphvizGraph.getGraphvizEdges()
	GetGraphvizEdges() []GraphvizGraph_GraphvizEdge
}

// GraphvizGraph_GraphvizNode is a Graphviz node.
type GraphvizGraph_GraphvizNode interface {
	// Java: GraphvizGraph.GraphvizNode.getId()
	GetId() string

	// Java: GraphvizGraph.GraphvizNode.getColor()
	GetColor() string

	// Java: GraphvizGraph.GraphvizNode.getLabel()
	GetLabel() string
}

// GraphvizGraph_GraphvizEdge is a Graphviz edge.
type GraphvizGraph_GraphvizEdge interface {
	// Java: GraphvizGraph.GraphvizEdge.getNode1Id()
	GetNode1Id() string

	// Java: GraphvizGraph.GraphvizEdge.getNode2Id()
	GetNode2Id() string

	// Java: GraphvizGraph.GraphvizEdge.getColor()
	GetColor() string

	// Java: GraphvizGraph.GraphvizEdge.getLabel()
	GetLabel() string
}

// annotationState records an object's saved annotation on a push/pop stack.
type annotationState struct {
	first  Annotatable
	second Annotation
}

// annotationStack is a reusable stack of saved node or edge annotation frames.
type annotationStack struct {
	states [][]annotationState
}

func pushAnnotationFrame(stack *annotationStack, capHint int) *[]annotationState {
	idx := len(stack.states)
	if idx < cap(stack.states) {
		stack.states = stack.states[:idx+1]
		top := &stack.states[idx]
		*top = (*top)[:0]
		return top
	}
	stack.states = append(stack.states, make([]annotationState, 0, capHint))
	return &stack.states[idx]
}

func popAnnotations(stack *annotationStack) {
	if len(stack.states) == 0 {
		panic(jutil.NewNoSuchElementException()) // ArrayDeque.pop()
	}
	topIdx := len(stack.states) - 1
	top := stack.states[topIdx]
	for i := range top {
		top[i].first.SetAnnotation(top[i].second)
		top[i] = annotationState{}
	}
	stack.states[topIdx] = top[:0]
	stack.states = stack.states[:topIdx]
}

func formatGraphValue(v any) string {
	if v == nil {
		return "null"
	}
	switch x := v.(type) {
	case string:
		return x
	case *string:
		if x == nil {
			return "null"
		}
		return *x
	case int:
		return strconv.Itoa(x)
	case int32:
		return strconv.FormatInt(int64(x), 10)
	case int64:
		return strconv.FormatInt(x, 10)
	case bool:
		return strconv.FormatBool(x)
	case *int:
		if x == nil {
			return "null"
		}
		return strconv.Itoa(*x)
	case *int32:
		if x == nil {
			return "null"
		}
		return strconv.FormatInt(int64(*x), 10)
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Pointer, reflect.Slice, reflect.Map, reflect.Func, reflect.Chan, reflect.Interface:
		if rv.IsNil() {
			return "null"
		}
	}
	if s, ok := v.(fmt.Stringer); ok {
		return s.String()
	}
	return fmt.Sprint(v)
}
