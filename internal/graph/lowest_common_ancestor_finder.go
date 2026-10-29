// Port of com/google/javascript/jscomp/graph/LowestCommonAncestorFinder.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// LowestCommonAncestorFinder implements a lowest common ancestor search
// algorithm.
//
// The LCA of a set of nodes is the node that is an ancestor to all of them
// but is not an ancestor of any other common ancestor. In a non-tree, There
// may be multiple LCAs for a given set of search nodes.
//
// In a cyclic graph, the LCAs may not be well defined. Within a cycle, all
// elements are both above and below one another, so there is no uniquely
// lowest element. If the set of common ancestors is rooted on a cycle, this
// implementation returns one or more elements of that cycle. Those elements
// are chosen arbitrarily but deterministically (as long as the underlying
// graph has deterministic iteration).
type LowestCommonAncestorFinder[N comparable, E any] interface {
	// FindAll executes a search on all the elements of roots.
	//
	// Java: LowestCommonAncestorFinder.findAll(Set<N>)
	FindAll(roots *jutil.LinkedHashSet[N]) *jutil.LinkedHashSet[N]

	// FindAllSlice executes a search on distinct elements of roots, appending
	// the resulting LCAs to dst without allocating sets.
	FindAllSlice(roots []N, dst []N) []N
}

// LowestCommonAncestorFinder_Factory is an abstraction for
// LowestCommonAncestorFinder::new, allowing injection in tests.
type LowestCommonAncestorFinder_Factory[N comparable, E any] func(graph DiGraph[N, E]) LowestCommonAncestorFinder[N, E]

// lcaColorNotLowest marks common ancestors that have lower common ancestors.
// Because its MSB is set (-1), bitwise OR with any root bitmask stays -1.
const lcaColorNotLowest int32 = -1

type lowestCommonAncestorFinder[N comparable, E any] struct {
	graph       *LinkedDirectedGraph[N, E]
	colors      []int32
	touched     []*LinkedDiGraphNode[N, E]
	searchQueue []*LinkedDiGraphNode[N, E]
}

// Java: LowestCommonAncestorFinder(DiGraph<N, E>)
func NewLowestCommonAncestorFinder[N comparable, E any](graph DiGraph[N, E]) LowestCommonAncestorFinder[N, E] {
	f := &lowestCommonAncestorFinder[N, E]{}
	if graph != nil {
		f.graph = graph.AsLinkedDirectedGraph()
	}
	return f
}

func (f *lowestCommonAncestorFinder[N, E]) ensureColorCapacity() {
	capNeeded := f.graph.NodeIndexCap()
	if len(f.colors) < capNeeded {
		if cap(f.colors) >= capNeeded {
			f.colors = f.colors[:capNeeded]
		} else {
			newCap := max(capNeeded, len(f.colors)*2, 64)
			f.colors = make([]int32, capNeeded, newCap)
		}
	}
}

func (f *lowestCommonAncestorFinder[N, E]) paintRoot(root N, bitForRoot int32) {
	rootNode := f.graph.GetNode(root)
	jutil.CheckNotNullMsg1(rootNode, "Root not present in graph: %s", root)
	idx := rootNode.Index
	if f.colors[idx] == 0 {
		f.touched = append(f.touched, rootNode)
	}
	f.colors[idx] |= bitForRoot
	f.paintAncestors(rootNode, bitForRoot)
}

func (f *lowestCommonAncestorFinder[N, E]) finishSearch(allColor int32, dst []N) []N {
	// Every ancestor of a common ancestor is already painted with allColor
	// (hence already in f.touched in insertion order).
	for i := 0; i < len(f.touched); i++ {
		node := f.touched[i]
		if f.colors[node.Index] == allColor {
			f.paintAncestors(node, lcaColorNotLowest)
		}
	}

	for i, node := range f.touched {
		if f.colors[node.Index] == allColor {
			dst = append(dst, node.Value)
		}
		f.colors[node.Index] = 0
		f.touched[i] = nil
	}
	f.touched = f.touched[:0]
	return dst
}

// Java: LowestCommonAncestorFinder.findAll(Set<N>)
func (f *lowestCommonAncestorFinder[N, E]) FindAll(roots *jutil.LinkedHashSet[N]) *jutil.LinkedHashSet[N] {
	if roots == nil {
		panic(jutil.NewNullPointerException())
	}
	// We reserved the MSB of each Color for bookkeeping.
	jutil.CheckArgumentMsg(roots.Size() <= 32-1, "Too many roots.")
	jutil.CheckState(len(f.touched) == 0)

	results := jutil.NewLinkedHashSet[N]()
	if roots.IsEmpty() {
		return results
	}

	f.ensureColorCapacity()
	allColor := (int32(1) << roots.Size()) - 1
	bitForRoot := int32(1)
	for root := range roots.All() {
		f.paintRoot(root, bitForRoot)
		bitForRoot <<= 1
	}

	for i := 0; i < len(f.touched); i++ {
		node := f.touched[i]
		if f.colors[node.Index] == allColor {
			f.paintAncestors(node, lcaColorNotLowest)
		}
	}
	for i, node := range f.touched {
		if f.colors[node.Index] == allColor {
			results.Add(node.Value)
		}
		f.colors[node.Index] = 0
		f.touched[i] = nil
	}
	f.touched = f.touched[:0]
	return results
}

func (f *lowestCommonAncestorFinder[N, E]) FindAllSlice(roots []N, dst []N) []N {
	jutil.CheckArgumentMsg(len(roots) <= 32-1, "Too many roots.")
	jutil.CheckState(len(f.touched) == 0)
	if len(roots) == 0 {
		return dst
	}

	f.ensureColorCapacity()
	allColor := (int32(1) << len(roots)) - 1
	bitForRoot := int32(1)
	for _, root := range roots {
		f.paintRoot(root, bitForRoot)
		bitForRoot <<= 1
	}
	return f.finishSearch(allColor, dst)
}

// paintAncestors paints all nodes above root with color.
//
// root itself will not have its color changed. color will be mixed with
// all existing colors on ancestor nodes.
//
// Java: LowestCommonAncestorFinder.paintAncestors(DiGraphNode<N, E>, Color)
func (f *lowestCommonAncestorFinder[N, E]) paintAncestors(root *LinkedDiGraphNode[N, E], color int32) {
	jutil.CheckState(len(f.searchQueue) == 0)

	f.searchQueue = append(f.searchQueue, root)

	for head := 0; head < len(f.searchQueue); head++ {
		curr := f.searchQueue[head]
		for _, parentEdge := range curr.InEdgeList.Slice() {
			parent := parentEdge.SourceNode
			if parent == root {
				continue // Don't paint `root`.
			}

			pIdx := parent.Index
			oldColor := f.colors[pIdx]
			if (oldColor & color) != color {
				// Only explore in directions that have not yet been painted.
				f.searchQueue = append(f.searchQueue, parent)
				if oldColor == 0 {
					f.touched = append(f.touched, parent)
				}
				f.colors[pIdx] = oldColor | color
			}
		}
	}
	clear(f.searchQueue)
	f.searchQueue = f.searchQueue[:0]
}
