// Port of com/google/javascript/jscomp/graph/GraphColoring.java @ closure-compiler v20260415

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
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// GraphColoring annotates the graph with a color in a way that no
// connected node will have the same color. Nodes of the same color can then
// be partitioned together and be represented by a super node. This class
// will merely annotate the nodes with a color using
// GraphNode#setAnnotation(Annotation) and provide a node to super node
// mapping with getPartitionSuperNode(Object). The given graph itself will
// not be modified.
//
// This algorithm is NOT deterministic by default. Passes that requires
// deterministic output should provide a Comparator in the constructor as a
// tie-breaker. This tie-break will be used when deciding which node should
// be colored first when multiple nodes have the same degree.
//
// N is the value type that the graph node stores; E is the value type that
// the graph edge stores.
type GraphColoring[N, E any] interface {
	// Color annotates the graph with Color objects using
	// GraphNode#setAnnotation(Annotation).
	//
	// Returns the number of unique colors need.
	//
	// Java: GraphColoring.color()
	Color() int

	// Java: GraphColoring.getPartitionSuperNode(N)
	GetPartitionSuperNode(node N) N

	// Java: GraphColoring.haveSameColor(N, N)
	HaveSameColor(first N, second N) bool

	// Java: GraphColoring.getGraph()
	GetGraph() AdjacencyGraph[N, E]

	GetGraphColoringBase() *GraphColoringBase[N, E]
}

// GraphColoringBase holds GraphColoring's fields and base implementations.
type GraphColoringBase[N, E any] struct {
	// Maps a color (represented by an integer) to a variable. If, for
	// example, the color 5 is mapped to "foo". Then any other variables
	// colored with the color 5 will now use the name "foo".
	//
	// Go-only: language-adapter — an N[] whose slots may be null: a nil
	// *N is Java's null slot (N may be a non-pointer type).
	colorToNodeMap []*N
	graph          AdjacencyGraph[N, E]
}

// InitGraphColoringBase is the GraphColoring(AdjacencyGraph) constructor.
//
// Java: GraphColoring(AdjacencyGraph<N, E>)
func (c *GraphColoringBase[N, E]) InitGraphColoringBase(graph AdjacencyGraph[N, E]) {
	c.graph = graph
}

// GetGraphColoringBase gives access to the base through the interface.
func (c *GraphColoringBase[N, E]) GetGraphColoringBase() *GraphColoringBase[N, E] { return c }

// graphColoring_colorOf is the unchecked cast (Color) annotation.
func graphColoring_colorOf(a Annotation) *GraphColoring_Color {
	if a == nil {
		return nil
	}
	c, ok := a.(*GraphColoring_Color)
	if !ok {
		panic(jutil.NewClassCastException(jutil.ToJavaString(a), "com.google.javascript.jscomp.graph.GraphColoring$Color"))
	}
	return c
}

// GetPartitionSuperNode: using the coloring as partitions, finds the node
// that represents that partition as the super node. The first to retrieve
// its partition will become the super node.
//
// Java: GraphColoring.getPartitionSuperNode(N)
func (c *GraphColoringBase[N, E]) GetPartitionSuperNode(node N) N {
	if c.colorToNodeMap == nil {
		panic(jutil.NewNullPointerExceptionMsg("No coloring founded. color() should be called first."))
	}
	color := graphColoring_colorOf(c.graph.GetNode(node).GetAnnotation())
	headNode := c.colorToNodeMap[color.value]
	if headNode == nil {
		c.colorToNodeMap[color.value] = &node
		return node
	} else {
		return *headNode
	}
}

// HaveSameColor: using the coloring as partitions, checks whether two given
// nodes have the same color.
//
// Java: GraphColoring.haveSameColor(N, N)
func (c *GraphColoringBase[N, E]) HaveSameColor(first N, second N) bool {
	if c.colorToNodeMap == nil {
		panic(jutil.NewNullPointerExceptionMsg("No coloring founded. color() should be called first."))
	}
	colorFirst := graphColoring_colorOf(c.graph.GetNode(first).GetAnnotation())
	colorSecond := graphColoring_colorOf(c.graph.GetNode(second).GetAnnotation())
	jutil.CheckNotNull(colorFirst)
	jutil.CheckNotNull(colorSecond)
	return colorFirst.Equals(colorSecond)
}

// Java: GraphColoring.getGraph()
func (c *GraphColoringBase[N, E]) GetGraph() AdjacencyGraph[N, E] {
	return c.graph
}

// GraphColoring_Color is the color of a node.
type GraphColoring_Color struct {
	value int
}

// Java: GraphColoring.Color(int)
func NewGraphColoring_Color(value int) *GraphColoring_Color {
	return &GraphColoring_Color{value: value}
}

// Java: GraphColoring.Color.equals(Object)
func (c *GraphColoring_Color) Equals(other any) bool {
	color, ok := other.(*GraphColoring_Color)
	if !ok || color == nil {
		return false
	} else {
		return c.value == color.value
	}
}

// Java: GraphColoring.Color.hashCode()
func (c *GraphColoring_Color) HashCode() int32 {
	return int32(c.value)
}

// GraphColoring_GreedyGraphColoring greedily assign nodes with high degree
// unique colors.
type GraphColoring_GreedyGraphColoring[N, E any] struct {
	GraphColoringBase[N, E]

	tieBreaker func(a, b N) int
}

// Java: GraphColoring.GreedyGraphColoring(AdjacencyGraph<N, E>)
func NewGraphColoring_GreedyGraphColoring[N, E any](graph AdjacencyGraph[N, E]) *GraphColoring_GreedyGraphColoring[N, E] {
	return NewGraphColoring_GreedyGraphColoringWithTieBreaker[N, E](graph, nil)
}

// NewGraphColoring_GreedyGraphColoringWithTieBreaker: in case of a tie
// between two nodes of the same degree, tieBreaker (a Comparator, may be
// null) will determine which node should be colored first.
//
// Java: GraphColoring.GreedyGraphColoring(AdjacencyGraph<N, E>, Comparator<N>)
func NewGraphColoring_GreedyGraphColoringWithTieBreaker[N, E any](
	graph AdjacencyGraph[N, E], tieBreaker func(a, b N) int) *GraphColoring_GreedyGraphColoring[N, E] {
	c := &GraphColoring_GreedyGraphColoring[N, E]{}
	c.InitGraphColoringBase(graph)
	c.tieBreaker = tieBreaker
	return c
}

// Java: GraphColoring.GreedyGraphColoring.color()
func (c *GraphColoring_GreedyGraphColoring[N, E]) Color() int {
	worklist := slices.Clone(c.graph.GetNodes())

	// Sort nodes by degree.
	// Collections.sort is a stable merge sort.
	slices.SortStableFunc(worklist, func(left, right GraphNode[N, E]) int {
		// comparing(GraphNode::getValue, ...)
		leftProperty, rightProperty := left.GetValue(), right.GetValue()
		// TODO(b/28382956): Take better advantage of Java8 comparing() to
		// simplify this
		result := c.graph.GetWeight(rightProperty) - c.graph.GetWeight(leftProperty)
		if result == 0 && c.tieBreaker != nil {
			return c.tieBreaker(leftProperty, rightProperty)
		}
		return result
	})

	// Idea: From the highest to lowest degree, assign any uncolored node
	// with a unique color if none of its neighbors has been assigned that
	// color.
	count := 0
	for {
		color := NewGraphColoring_Color(count)
		subgraph := c.graph.NewSubGraph()
		kept := worklist[:0]
		for _, node := range worklist {
			if subgraph.IsIndependentOf(node.GetValue()) {
				subgraph.AddNode(node.GetValue())
				node.SetAnnotation(color)
				// i.remove()
			} else {
				kept = append(kept, node)
			}
		}
		clear(worklist[len(kept):])
		worklist = kept
		count++
		if len(worklist) == 0 {
			break
		}
	}
	c.colorToNodeMap = make([]*N, count)
	return count
}
