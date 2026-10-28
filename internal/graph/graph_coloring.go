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

// GraphColoring annotates the graph with a color in a way that no
// connected node will have the same color. Nodes of the same color can then
// be partitioned together and be represented by a super node. This type
// annotates the nodes with a color using GraphNode.SetAnnotation and provides
// a node to super node mapping with GetPartitionSuperNode. The given graph
// itself will not be modified.
//
// This algorithm is NOT deterministic by default. Passes that require
// deterministic output should provide a comparator in the constructor as a
// tie-breaker. This tie-break will be used when deciding which node should
// be colored first when multiple nodes have the same degree.
//
// N is the value type that the graph node stores; E is the value type that
// the graph edge stores.
type GraphColoring[N, E any] = *GraphColoring_GreedyGraphColoring[N, E]

// GraphColoring_ColorOf is the unchecked cast (Color) annotation.
func GraphColoring_ColorOf(a Annotation) *GraphColoring_Color {
	if a == nil {
		return nil
	}
	c, ok := a.(*GraphColoring_Color)
	if !ok {
		panic(jutil.NewClassCastException(jutil.ToJavaString(a), "com.google.javascript.jscomp.graph.GraphColoring$Color"))
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

// GraphColoring_GreedyGraphColoring greedily assigns nodes with high degree
// unique colors.
type GraphColoring_GreedyGraphColoring[N, E any] struct {
	graph         AdjacencyGraph[N, E]
	TieBreaker    func(a, b N) int
	superNodes    []N
	hasSuperNodes []bool
	colored       bool
}

// Java: GraphColoring.GreedyGraphColoring(AdjacencyGraph<N, E>)
func NewGraphColoring_GreedyGraphColoring[N, E any](graph AdjacencyGraph[N, E]) *GraphColoring_GreedyGraphColoring[N, E] {
	return NewGraphColoring_GreedyGraphColoringWithTieBreaker[N, E](graph, nil)
}

// NewGraphColoring_GreedyGraphColoringWithTieBreaker: in case of a tie
// between two nodes of the same degree, tieBreaker (a Comparator, may be
// nil) will determine which node should be colored first.
//
// Java: GraphColoring.GreedyGraphColoring(AdjacencyGraph<N, E>, Comparator<N>)
func NewGraphColoring_GreedyGraphColoringWithTieBreaker[N, E any](
	graph AdjacencyGraph[N, E], tieBreaker func(a, b N) int) *GraphColoring_GreedyGraphColoring[N, E] {
	return &GraphColoring_GreedyGraphColoring[N, E]{
		graph:      graph,
		TieBreaker: tieBreaker,
	}
}

// GetPartitionSuperNode: using the coloring as partitions, finds the node
// that represents that partition as the super node. The first to retrieve
// its partition will become the super node.
//
// Java: GraphColoring.getPartitionSuperNode(N)
func (c *GraphColoring_GreedyGraphColoring[N, E]) GetPartitionSuperNode(node N) N {
	if !c.colored {
		panic(jutil.NewNullPointerExceptionMsg("No coloring founded. color() should be called first."))
	}
	color := GraphColoring_ColorOf(c.graph.GetNode(node).GetAnnotation())
	idx := color.Value
	if !c.hasSuperNodes[idx] {
		c.superNodes[idx] = node
		c.hasSuperNodes[idx] = true
		return node
	}
	return c.superNodes[idx]
}

// HaveSameColor: using the coloring as partitions, checks whether two given
// nodes have the same color.
//
// Java: GraphColoring.haveSameColor(N, N)
func (c *GraphColoring_GreedyGraphColoring[N, E]) HaveSameColor(first N, second N) bool {
	if !c.colored {
		panic(jutil.NewNullPointerExceptionMsg("No coloring founded. color() should be called first."))
	}
	colorFirst := GraphColoring_ColorOf(c.graph.GetNode(first).GetAnnotation())
	colorSecond := GraphColoring_ColorOf(c.graph.GetNode(second).GetAnnotation())
	jutil.CheckNotNull(colorFirst)
	jutil.CheckNotNull(colorSecond)
	return colorFirst.Value == colorSecond.Value
}

// Java: GraphColoring.getGraph()
func (c *GraphColoring_GreedyGraphColoring[N, E]) GetGraph() AdjacencyGraph[N, E] {
	return c.graph
}

type coloringEntry[N, E any] struct {
	node   GraphNode[N, E]
	value  N
	weight int
	idx    int
}

// Color annotates the graph with Color objects using GraphNode.SetAnnotation
// and returns the number of unique colors needed.
//
// Java: GraphColoring.GreedyGraphColoring.color()
func (c *GraphColoring_GreedyGraphColoring[N, E]) Color() int {
	nodes := c.graph.GetNodes()
	n := len(nodes)
	if n == 0 {
		c.initSuperNodes(1)
		return 1
	}

	entries := make([]coloringEntry[N, E], n)
	switch nodes[0].(type) {
	case *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]:
		for i, node := range nodes {
			uNode := node.(*LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E])
			entries[i] = coloringEntry[N, E]{
				node:   node,
				value:  uNode.Value,
				weight: len(uNode.NeighborEdges),
				idx:    i,
			}
		}
	case *LinkedDirectedGraph_LinkedDiGraphNode[N, E]:
		for i, node := range nodes {
			dNode := node.(*LinkedDirectedGraph_LinkedDiGraphNode[N, E])
			entries[i] = coloringEntry[N, E]{
				node:   node,
				value:  dNode.Value,
				weight: dNode.InEdgeList.Size() + dNode.OutEdgeList.Size(),
				idx:    int(dNode.Index),
			}
		}
	default:
		for i, node := range nodes {
			val := node.GetValue()
			entries[i] = coloringEntry[N, E]{
				node:   node,
				value:  val,
				weight: c.graph.GetWeight(val),
				idx:    i,
			}
		}
	}

	// Sort nodes by descending weight (degree), breaking ties deterministically
	// when TieBreaker is provided, and stably otherwise.
	tieBreaker := c.TieBreaker
	slices.SortStableFunc(entries, func(left, right coloringEntry[N, E]) int {
		if result := right.weight - left.weight; result != 0 {
			return result
		}
		if tieBreaker != nil {
			return tieBreaker(left.value, right.value)
		}
		return 0
	})

	var count int
	switch entries[0].node.(type) {
	case *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]:
		count = colorUndirectedEntries(entries)
	case *LinkedDirectedGraph_LinkedDiGraphNode[N, E]:
		count = colorDirectedEntries(entries)
	default:
		count = c.colorSubGraphEntries(entries)
	}

	c.initSuperNodes(count)
	return count
}

func (c *GraphColoring_GreedyGraphColoring[N, E]) initSuperNodes(count int) {
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

func colorUndirectedEntries[N, E any](worklist []coloringEntry[N, E]) int {
	n := len(worklist)
	nodeIdx := make(map[*LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]]int, n)
	for i := range worklist {
		uNode := worklist[i].node.(*LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E])
		nodeIdx[uNode] = i
		worklist[i].idx = i
	}
	// forbidden[i] == epoch means worklist node i is adjacent to a node already
	// assigned color (epoch - 1). Epochs start at 1 so zero-initialized slots are
	// never forbidden and require no clearing between color iterations.
	forbidden := make([]int, n)
	count := 0
	for len(worklist) > 0 {
		epoch := count + 1
		color := NewGraphColoring_Color(count)
		kept := worklist[:0]
		for _, entry := range worklist {
			if forbidden[entry.idx] != epoch {
				uNode := entry.node.(*LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E])
				for _, edge := range uNode.NeighborEdges {
					a := edge.GetNodeA().(*LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E])
					if a == uNode {
						b := edge.GetNodeB().(*LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E])
						forbidden[nodeIdx[b]] = epoch
					} else {
						forbidden[nodeIdx[a]] = epoch
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

func colorDirectedEntries[N, E any](worklist []coloringEntry[N, E]) int {
	maxIdx := 0
	for i := range worklist {
		if worklist[i].idx > maxIdx {
			maxIdx = worklist[i].idx
		}
	}
	forbidden := make([]int, maxIdx+1)
	count := 0
	for len(worklist) > 0 {
		epoch := count + 1
		color := NewGraphColoring_Color(count)
		kept := worklist[:0]
		for _, entry := range worklist {
			if forbidden[entry.idx] != epoch {
				dNode := entry.node.(*LinkedDirectedGraph_LinkedDiGraphNode[N, E])
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

func (c *GraphColoring_GreedyGraphColoring[N, E]) colorSubGraphEntries(worklist []coloringEntry[N, E]) int {
	count := 0
	for len(worklist) > 0 {
		color := NewGraphColoring_Color(count)
		subgraph := c.graph.NewSubGraph()
		kept := worklist[:0]
		for _, entry := range worklist {
			if subgraph.IsIndependentOf(entry.value) {
				subgraph.AddNode(entry.value)
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
