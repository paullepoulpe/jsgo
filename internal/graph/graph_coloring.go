// Port of com/google/javascript/jscomp/graph/GraphColoring.java @ closure-compiler v20261021

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
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

var preallocatedColors = func() [256]GraphColoring_Color {
	var arr [256]GraphColoring_Color
	for i := range arr {
		arr[i].Value = i
	}
	return arr
}()

// GraphColoring aliases *GreedyGraphColoring, the sole graph coloring implementation.
type (
	GraphColoring[N comparable, E any]                     = *GreedyGraphColoring[N, E]
	GraphColoring_GreedyGraphColoring[N comparable, E any] = GreedyGraphColoring[N, E]
)

// GraphColoring_ColorOf is the unchecked cast (Color) annotation.
func GraphColoring_ColorOf(a Annotation) *GraphColoring_Color {
	if a == nil {
		return nil
	}
	c, ok := a.(*GraphColoring_Color)
	if !ok {
		panic(jutil.NewClassCastException(formatGraphValue(a), "com.google.javascript.jscomp.graph.GraphColoring$Color"))
	}
	return c
}

// GraphColoring_Color is the color of a node.
type GraphColoring_Color struct {
	Value int
}

// Java: GraphColoring.Color(int)
func NewGraphColoring_Color(value int) *GraphColoring_Color {
	if value >= 0 && value < len(preallocatedColors) {
		return &preallocatedColors[value]
	}
	return &GraphColoring_Color{Value: value}
}

// Java: GraphColoring.Color.equals(Object)
func (c *GraphColoring_Color) Equals(other any) bool {
	color, ok := other.(*GraphColoring_Color)
	return ok && color != nil && c.Value == color.Value
}

// Java: GraphColoring.Color.hashCode()
func (c *GraphColoring_Color) HashCode() int32 {
	return int32(c.Value)
}

// GreedyGraphColoring greedily assigns nodes with high degree unique colors
// such that no connected nodes share a color. Nodes of the same color can be
// partitioned together and represented by a super node via GetPartitionSuperNode.
// The given graph itself is not modified beyond node annotations.
type GreedyGraphColoring[N comparable, E any] struct {
	graph         Graph[N, E]
	TieBreaker    func(a, b N) int
	superNodes    []N
	hasSuperNodes []bool
	colored       bool
}

// Java: GraphColoring.GreedyGraphColoring(AdjacencyGraph<N, E>)
func NewGreedyGraphColoring[N comparable, E any](graph Graph[N, E]) *GreedyGraphColoring[N, E] {
	return NewGreedyGraphColoringWithTieBreaker[N, E](graph, nil)
}

// NewGraphColoring_GreedyGraphColoring preserves the legacy constructor name.
func NewGraphColoring_GreedyGraphColoring[N comparable, E any](graph Graph[N, E]) *GreedyGraphColoring[N, E] {
	return NewGreedyGraphColoringWithTieBreaker[N, E](graph, nil)
}

// NewGreedyGraphColoringWithTieBreaker constructs a greedy graph coloring with
// an optional tie-breaking comparator for equal-degree nodes.
//
// Java: GraphColoring.GreedyGraphColoring(AdjacencyGraph<N, E>, Comparator<N>)
func NewGreedyGraphColoringWithTieBreaker[N comparable, E any](
	graph Graph[N, E], tieBreaker func(a, b N) int) *GreedyGraphColoring[N, E] {
	return &GreedyGraphColoring[N, E]{
		graph:      graph,
		TieBreaker: tieBreaker,
	}
}

// NewGraphColoring_GreedyGraphColoringWithTieBreaker preserves the legacy constructor name.
func NewGraphColoring_GreedyGraphColoringWithTieBreaker[N comparable, E any](
	graph Graph[N, E], tieBreaker func(a, b N) int) *GreedyGraphColoring[N, E] {
	return NewGreedyGraphColoringWithTieBreaker[N, E](graph, tieBreaker)
}

func (c *GreedyGraphColoring[N, E]) nodeAnnotation(node N) Annotation {
	switch g := c.graph.(type) {
	case *LinkedUndirectedGraph[N, E]:
		uNode := g.Nodes.Get(node)
		if uNode == nil {
			panic(jutil.NewNullPointerException())
		}
		return uNode.GetAnnotation()
	case DiGraph[N, E]:
		dNode := g.AsLinkedDirectedGraph().Nodes.Get(node)
		if dNode == nil {
			panic(jutil.NewNullPointerException())
		}
		return dNode.GetAnnotation()
	default:
		panic(jutil.NewUnsupportedOperationExceptionMsg("unsupported graph implementation"))
	}
}

// GetPartitionSuperNode finds the representative super node for node's color
// partition. The first node queried in a partition becomes its super node.
//
// Java: GraphColoring.getPartitionSuperNode(N)
func (c *GreedyGraphColoring[N, E]) GetPartitionSuperNode(node N) N {
	if !c.colored {
		panic(jutil.NewNullPointerExceptionMsg("No coloring founded. color() should be called first."))
	}
	color := GraphColoring_ColorOf(c.nodeAnnotation(node))
	idx := color.Value
	if !c.hasSuperNodes[idx] {
		c.superNodes[idx] = node
		c.hasSuperNodes[idx] = true
		return node
	}
	return c.superNodes[idx]
}

// HaveSameColor reports whether first and second were assigned the same color.
//
// Java: GraphColoring.haveSameColor(N, N)
func (c *GreedyGraphColoring[N, E]) HaveSameColor(first N, second N) bool {
	if !c.colored {
		panic(jutil.NewNullPointerExceptionMsg("No coloring founded. color() should be called first."))
	}
	colorFirst := GraphColoring_ColorOf(c.nodeAnnotation(first))
	colorSecond := GraphColoring_ColorOf(c.nodeAnnotation(second))
	jutil.CheckNotNull(colorFirst)
	jutil.CheckNotNull(colorSecond)
	return colorFirst.Value == colorSecond.Value
}

// Java: GraphColoring.getGraph()
func (c *GreedyGraphColoring[N, E]) GetGraph() Graph[N, E] {
	return c.graph
}

// Color annotates the graph nodes with Color objects and returns the number of
// unique colors needed.
//
// Java: GraphColoring.GreedyGraphColoring.color()
func (c *GreedyGraphColoring[N, E]) Color() int {
	n := c.graph.GetNodeCount()
	if n == 0 {
		c.initSuperNodes(1)
		return 1
	}

	var count int
	switch g := c.graph.(type) {
	case *LinkedUndirectedGraph[N, E]:
		count = colorUndirectedGraph(g, c.TieBreaker)
	case DiGraph[N, E]:
		count = colorDirectedGraph(g.AsLinkedDirectedGraph(), c.TieBreaker)
	default:
		panic(jutil.NewUnsupportedOperationExceptionMsg("unsupported graph implementation"))
	}

	c.initSuperNodes(count)
	return count
}

func (c *GreedyGraphColoring[N, E]) initSuperNodes(count int) {
	if cap(c.superNodes) >= count {
		clear(c.superNodes[:count])
		c.superNodes = c.superNodes[:count]
	} else {
		c.superNodes = make([]N, count)
	}
	if cap(c.hasSuperNodes) >= count {
		clear(c.hasSuperNodes[:count])
		c.hasSuperNodes = c.hasSuperNodes[:count]
	} else {
		c.hasSuperNodes = make([]bool, count)
	}
	c.colored = true
}

type undirectedColoringEntry[N, E any] struct {
	node   *LinkedUndirectedGraphNode[N, E]
	weight int
	idx    int
}

func colorUndirectedGraph[N comparable, E any](g *LinkedUndirectedGraph[N, E], tieBreaker func(a, b N) int) int {
	n := g.Nodes.Size()
	worklist := make([]undirectedColoringEntry[N, E], 0, n)
	for uNode := range g.Nodes.Values() {
		worklist = append(worklist, undirectedColoringEntry[N, E]{
			node:   uNode,
			weight: len(uNode.NeighborEdges),
		})
	}
	slices.SortStableFunc(worklist, func(left, right undirectedColoringEntry[N, E]) int {
		if result := right.weight - left.weight; result != 0 {
			return result
		}
		if tieBreaker != nil {
			return tieBreaker(left.node.Value, right.node.Value)
		}
		return 0
	})

	nodeIdx := make(map[*LinkedUndirectedGraphNode[N, E]]int, n)
	for i := range worklist {
		nodeIdx[worklist[i].node] = i
		worklist[i].idx = i
	}
	forbidden := make([]int, n)
	count := 0
	for len(worklist) > 0 {
		epoch := count + 1
		color := NewGraphColoring_Color(count)
		kept := worklist[:0]
		for _, entry := range worklist {
			if forbidden[entry.idx] != epoch {
				uNode := entry.node
				for _, edge := range uNode.NeighborEdges {
					if edge.nodeA == uNode {
						forbidden[nodeIdx[edge.nodeB]] = epoch
					} else {
						forbidden[nodeIdx[edge.nodeA]] = epoch
					}
				}
				uNode.SetAnnotation(color)
			} else {
				kept = append(kept, entry)
			}
		}
		clear(worklist[len(kept):])
		worklist = kept
		count++
	}
	return count
}

type directedColoringEntry[N, E any] struct {
	node   *LinkedDiGraphNode[N, E]
	weight int
}

func colorDirectedGraph[N comparable, E any](g *LinkedDirectedGraph[N, E], tieBreaker func(a, b N) int) int {
	n := g.Nodes.Size()
	worklist := make([]directedColoringEntry[N, E], 0, n)
	maxIdx := uint32(0)
	for dNode := range g.Nodes.Values() {
		if dNode.Index > maxIdx {
			maxIdx = dNode.Index
		}
		worklist = append(worklist, directedColoringEntry[N, E]{
			node:   dNode,
			weight: dNode.InEdgeList.Size() + dNode.OutEdgeList.Size(),
		})
	}
	slices.SortStableFunc(worklist, func(left, right directedColoringEntry[N, E]) int {
		if result := right.weight - left.weight; result != 0 {
			return result
		}
		if tieBreaker != nil {
			return tieBreaker(left.node.Value, right.node.Value)
		}
		return 0
	})

	forbidden := make([]int, int(maxIdx)+1)
	count := 0
	for len(worklist) > 0 {
		epoch := count + 1
		color := NewGraphColoring_Color(count)
		kept := worklist[:0]
		for _, entry := range worklist {
			dNode := entry.node
			if forbidden[dNode.Index] != epoch {
				for _, edge := range dNode.InEdgeList.Slice() {
					forbidden[edge.SourceNode.Index] = epoch
				}
				for _, edge := range dNode.OutEdgeList.Slice() {
					forbidden[edge.DestNode.Index] = epoch
				}
				dNode.SetAnnotation(color)
			} else {
				kept = append(kept, entry)
			}
		}
		clear(worklist[len(kept):])
		worklist = kept
		count++
	}
	return count
}
