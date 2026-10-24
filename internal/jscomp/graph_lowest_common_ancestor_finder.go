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

package jscomp

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
//
// The class is not final (tests subclass it to stub findAll), so it is an
// interface plus LowestCommonAncestorFinderBase (PORTING.md §4).
type LowestCommonAncestorFinder[N comparable, E any] interface {
	// FindAll executes a search on all the elements of roots.
	//
	// This is a general-purpose, bare-bones implementation. There are lots
	// of special case optimizations that could be applied.
	//
	// Java: LowestCommonAncestorFinder.findAll(Set<N>)
	FindAll(roots jutil.Set[N]) *jutil.LinkedHashSet[N]

	GetLowestCommonAncestorFinderBase() *LowestCommonAncestorFinderBase[N, E]
}

// LowestCommonAncestorFinder_Factory is an abstraction for
// LowestCommonAncestorFinder::new.
//
// This interface allows injection in tests.
//
// Go-only: language-adapter — the @FunctionalInterface is a func type, so
// constructor references and lambdas convert directly.
type LowestCommonAncestorFinder_Factory[N comparable, E any] func(graph DiGraph[N, E]) LowestCommonAncestorFinder[N, E]

// lowestCommonAncestorFinder_Color is a "color" for a node, encoded as a
// combination of other colors using a one-hot scheme.
type lowestCommonAncestorFinder_Color struct {
	bitset int32
}

// A set of commonly used colors to minimize allocations.
//
// An array is used to prevent bounds checking overhead.
var lowestCommonAncestorFinder_Color_COMMON_COLOR_CACHE = func() [2 << 5]*lowestCommonAncestorFinder_Color {
	var cache [2 << 5]*lowestCommonAncestorFinder_Color
	for i := range cache {
		cache[i] = newLowestCommonAncestorFinder_Color(int32(i))
	}
	return cache
}()

// A color that when mixed with any other color x returns x.
//
// Mixing this color is idempotent.
var lowestCommonAncestorFinder_Color_BLANK = jutil.CheckNotNull(lowestCommonAncestorFinder_Color_COMMON_COLOR_CACHE[0])

// A color for common ancestors that indicates that there are lower common
// ancestors.
//
// Because this color sets its MSB high, it can never equal any other color.
// Also notice that mixing this Color with any other produces this Color
// again; mixing it is nullipotent.
var lowestCommonAncestorFinder_Color_NOT_LOWEST = newLowestCommonAncestorFinder_Color(-1)

// Java: LowestCommonAncestorFinder.Color.create(int)
func lowestCommonAncestorFinder_Color_create(bitset int32) *lowestCommonAncestorFinder_Color {
	if bitset < 0 {
		jutil.CheckArgument(bitset == -1)
		return lowestCommonAncestorFinder_Color_NOT_LOWEST
	} else if int(bitset) < len(lowestCommonAncestorFinder_Color_COMMON_COLOR_CACHE) {
		return lowestCommonAncestorFinder_Color_COMMON_COLOR_CACHE[bitset]
	}
	return newLowestCommonAncestorFinder_Color(bitset)
}

// Java: LowestCommonAncestorFinder.Color(int)
func newLowestCommonAncestorFinder_Color(bitset int32) *lowestCommonAncestorFinder_Color {
	return &lowestCommonAncestorFinder_Color{bitset: bitset}
}

// Java: LowestCommonAncestorFinder.Color.mix(Color)
func (c *lowestCommonAncestorFinder_Color) mix(other *lowestCommonAncestorFinder_Color) *lowestCommonAncestorFinder_Color {
	if c.bitset == other.bitset {
		return c
	}
	return lowestCommonAncestorFinder_Color_create(c.bitset | other.bitset)
}

// Java: LowestCommonAncestorFinder.Color.contains(Color)
func (c *lowestCommonAncestorFinder_Color) contains(other *lowestCommonAncestorFinder_Color) bool {
	return (c.bitset & other.bitset) == other.bitset
}

// Java: LowestCommonAncestorFinder.Color.equals(Object)
func (c *lowestCommonAncestorFinder_Color) Equals(other any) bool {
	return c.bitset == other.(*lowestCommonAncestorFinder_Color).bitset
}

// Java: LowestCommonAncestorFinder.Color.hashCode()
func (c *lowestCommonAncestorFinder_Color) HashCode() int32 {
	return c.bitset
}

// LowestCommonAncestorFinderBase holds LowestCommonAncestorFinder's fields
// and implementation.
type LowestCommonAncestorFinderBase[N comparable, E any] struct {
	graph DiGraph[N, E]

	searchColoring *jutil.LinkedHashMap[DiGraph_DiGraphNode[N, E], *lowestCommonAncestorFinder_Color]
	// ArrayDeque<DiGraphNode<N, E>> used as a FIFO queue.
	searchQueue []DiGraph_DiGraphNode[N, E]
}

// Java: LowestCommonAncestorFinder(DiGraph<N, E>)
func NewLowestCommonAncestorFinder[N comparable, E any](graph DiGraph[N, E]) *LowestCommonAncestorFinderBase[N, E] {
	f := &LowestCommonAncestorFinderBase[N, E]{}
	f.InitLowestCommonAncestorFinderBase(graph)
	return f
}

// InitLowestCommonAncestorFinderBase runs the constructor on f (for
// subclasses, which embed LowestCommonAncestorFinderBase).
func (f *LowestCommonAncestorFinderBase[N, E]) InitLowestCommonAncestorFinderBase(graph DiGraph[N, E]) {
	f.graph = graph
	f.searchColoring = jutil.NewLinkedHashMap[DiGraph_DiGraphNode[N, E], *lowestCommonAncestorFinder_Color]()
}

// GetLowestCommonAncestorFinderBase gives access to the base fields
// through the interface.
func (f *LowestCommonAncestorFinderBase[N, E]) GetLowestCommonAncestorFinderBase() *LowestCommonAncestorFinderBase[N, E] {
	return f
}

// Java: LowestCommonAncestorFinder.findAll(Set<N>)
func (f *LowestCommonAncestorFinderBase[N, E]) FindAll(roots jutil.Set[N]) *jutil.LinkedHashSet[N] {
	// We reserved the MSB of each Color for bookkeeping.
	jutil.CheckArgumentMsg(roots.Size() <= 32-1, "Too many roots.")
	jutil.CheckState(f.searchColoring.IsEmpty())

	// In two's-complement, (2^n - 1) sets the lowest n bits high.
	allColor := lowestCommonAncestorFinder_Color_create((int32(1) << roots.Size()) - 1)

	// Paint up from each root using the color associated with that root.
	//
	// When done, the set of common ancestors is the set of nodes painted
	// `allColor`.
	bitForRoot := int32(1)
	for root := range roots.All() {
		rootNode := f.graph.GetDiGraphNode(root)
		jutil.CheckNotNullMsg1(rootNode, "Root not present in graph: %s", root)

		color := lowestCommonAncestorFinder_Color_create(bitForRoot)

		// Preserve any existing colors (Map.merge with Color::mix).
		if old, ok := f.searchColoring.Lookup(rootNode); ok {
			f.searchColoring.Put(rootNode, old.mix(color))
		} else {
			f.searchColoring.Put(rootNode, color)
		}
		f.paintAncestors(rootNode, color)
		bitForRoot <<= 1
	}

	// For every common ancestor, paint all of its ancestors with a color
	// indicating it is not the lowest.
	//
	// (Every ancestor of a common ancestor is already a key, so this only
	// updates values in place: Java's forEach sees the updates and never
	// throws ConcurrentModificationException; All() behaves the same.)
	for node, color := range f.searchColoring.All() {
		if color.Equals(allColor) {
			f.paintAncestors(node, lowestCommonAncestorFinder_Color_NOT_LOWEST)
		}
	}

	results := jutil.NewLinkedHashSet[N]() // ImmutableSet.builder()
	for node, color := range f.searchColoring.All() {
		if color.Equals(allColor) {
			results.Add(node.GetValue())
		}
	}

	f.searchColoring.Clear()
	f.searchQueue = nil
	return results
}

// paintAncestors paints all nodes above root with color.
//
// root itself will not have its color changed. color will be mixed with
// all existing colors on ancestor nodes.
//
// Java: LowestCommonAncestorFinder.paintAncestors(DiGraphNode<N, E>, Color)
func (f *LowestCommonAncestorFinderBase[N, E]) paintAncestors(root DiGraph_DiGraphNode[N, E], color *lowestCommonAncestorFinder_Color) {
	jutil.CheckState(len(f.searchQueue) == 0)

	f.searchQueue = append(f.searchQueue, root)

	for len(f.searchQueue) != 0 {
		curr := f.searchQueue[0]
		f.searchQueue = f.searchQueue[1:]
		parentEdges := curr.GetInEdges()

		for _, parentEdge := range parentEdges {
			parent := parentEdge.GetSource()
			if parent == root {
				continue // Don't paint `root`.
			}

			oldColor := f.searchColoring.GetOrDefault(parent, lowestCommonAncestorFinder_Color_BLANK)
			if !oldColor.contains(color) {
				// Only explore in directions that have not yet been painted.
				f.searchQueue = append(f.searchQueue, parent)
				f.searchColoring.Put(parent, oldColor.mix(color))
			}
		}
	}
}

var _ LowestCommonAncestorFinder[string, string] = (*LowestCommonAncestorFinderBase[string, string])(nil)
