// Port of com/google/javascript/jscomp/graph/GraphColoring.java @ closure-compiler v20261104

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
type GraphColoring[N, E comparable] = *GreedyGraphColoring[N, E]

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
type GreedyGraphColoring[N, E comparable] struct {
	graph         Graph[N, E]
	TieBreaker    func(a, b N) int
	superNodes    []N
	hasSuperNodes []bool
	colored       bool
}

// Java: GraphColoring.GreedyGraphColoring(AdjacencyGraph<N, E>)
func NewGraphColoring_GreedyGraphColoring[N, E comparable](graph Graph[N, E]) *GreedyGraphColoring[N, E] {
	return NewGraphColoring_GreedyGraphColoringWithTieBreaker[N, E](graph, nil)
}

// NewGraphColoring_GreedyGraphColoringWithTieBreaker constructs a greedy graph
// coloring with an optional tie-breaking comparator for equal-degree nodes.
//
// Java: GraphColoring.GreedyGraphColoring(AdjacencyGraph<N, E>, Comparator<N>)
func NewGraphColoring_GreedyGraphColoringWithTieBreaker[N, E comparable](
	graph Graph[N, E], tieBreaker func(a, b N) int) *GreedyGraphColoring[N, E] {
	return &GreedyGraphColoring[N, E]{
		graph:      graph,
		TieBreaker: tieBreaker,
	}
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

type coloringEntry[N comparable, Node interface{ SetAnnotation(Annotation) }] struct {
	node   Node
	value  N
	idx    int
	weight int
}

func colorWorklist[N comparable, Node interface{ SetAnnotation(Annotation) }](
	worklist []coloringEntry[N, Node],
	forbiddenLen int,
	tieBreaker func(a, b N) int,
	markNeighbors func(node Node, forbidden []int, epoch int),
) int {
	slices.SortStableFunc(worklist, func(left, right coloringEntry[N, Node]) int {
		if result := right.weight - left.weight; result != 0 {
			return result
		}
		if tieBreaker != nil {
			return tieBreaker(left.value, right.value)
		}
		return 0
	})

	forbidden := make([]int, forbiddenLen)
	count := 0
	for len(worklist) > 0 {
		epoch := count + 1
		color := NewGraphColoring_Color(count)
		kept := worklist[:0]
		for _, entry := range worklist {
			if forbidden[entry.idx] != epoch {
				markNeighbors(entry.node, forbidden, epoch)
				entry.node.SetAnnotation(color)
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

func colorUndirectedGraph[N, E comparable](g *LinkedUndirectedGraph[N, E], tieBreaker func(a, b N) int) int {
	uNodes := g.Nodes.Slice()
	worklist := make([]coloringEntry[N, *LinkedUndirectedGraphNode[N, E]], len(uNodes))
	for i, uNode := range uNodes {
		uNode.colorIdx = int32(i)
		worklist[i] = coloringEntry[N, *LinkedUndirectedGraphNode[N, E]]{
			node:   uNode,
			value:  uNode.Value,
			idx:    i,
			weight: uNode.NeighborEdges.Size(),
		}
	}
	return colorWorklist(worklist, len(uNodes), tieBreaker, func(uNode *LinkedUndirectedGraphNode[N, E], forbidden []int, epoch int) {
		for _, edge := range uNode.NeighborEdges.Slice() {
			if edge.nodeA == uNode {
				forbidden[edge.nodeB.colorIdx] = epoch
			} else {
				forbidden[edge.nodeA.colorIdx] = epoch
			}
		}
	})
}

func colorDirectedGraph[N, E comparable](g *LinkedDirectedGraph[N, E], tieBreaker func(a, b N) int) int {
	dNodes := g.Nodes.Slice()
	worklist := make([]coloringEntry[N, *LinkedDiGraphNode[N, E]], len(dNodes))
	maxIdx := 0
	for i, dNode := range dNodes {
		idx := int(dNode.Index)
		if idx > maxIdx {
			maxIdx = idx
		}
		worklist[i] = coloringEntry[N, *LinkedDiGraphNode[N, E]]{
			node:   dNode,
			value:  dNode.Value,
			idx:    idx,
			weight: dNode.InEdgeList.Size() + dNode.OutEdgeList.Size(),
		}
	}
	return colorWorklist(worklist, maxIdx+1, tieBreaker, func(dNode *LinkedDiGraphNode[N, E], forbidden []int, epoch int) {
		for _, edge := range dNode.InEdgeList.Slice() {
			forbidden[edge.SourceNode.Index] = epoch
		}
		for _, edge := range dNode.OutEdgeList.Slice() {
			forbidden[edge.DestNode.Index] = epoch
		}
	})
}
