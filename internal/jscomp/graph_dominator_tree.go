// Port of com/google/javascript/jscomp/graph/DominatorTree.java @ closure-compiler v20261021

/*
 * Copyright 2026 The Closure Compiler Authors.
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

// DominatorTree computes the dominator tree of a directed graph.
//
// A node d dominates a node n if every path from the entry node to n must
// go through d.
//
// The "immediate dominator" of n is the unique dominator d of n such that
// d != n and d is dominated by every other dominator of n (excluding n
// itself).
//
// This implementation uses the iterative algorithm by Cooper, Harvey, and
// Kennedy. See https://www.cs.tufts.edu/~nr/cs257/archive/keith-cooper/dom14.pdf.
//
// N is the value type that the graph node stores.
//
// Go-only: language-adapter — node values are compared with Go == (Java
// equals() for strings, numbers and identity-compared objects).
type DominatorTree[N comparable] struct {
	// Mapping from node value to its immediate dominator's value.
	// entry node maps to null.
	idoms *jutil.LinkedHashMap[N, N]

	// Mapping from node value to the number of nodes it dominates
	// (including itself).
	subtreeSizes *jutil.LinkedHashMap[N, int]
}

// Java: DominatorTree(Map<N, N>, Map<N, Integer>)
func newDominatorTree[N comparable](idoms *jutil.LinkedHashMap[N, N], subtreeSizes *jutil.LinkedHashMap[N, int]) *DominatorTree[N] {
	return &DominatorTree[N]{
		idoms:        jutil.NewLinkedHashMapFrom(idoms),
		subtreeSizes: jutil.NewLinkedHashMapFrom(subtreeSizes),
	}
}

// DominatorTree_Compute computes the dominator tree for the given graph
// starting from the entry node.
//
// graph is the directed graph to analyze; entry is the entry node of the
// graph. Returns a DominatorTree instance.
//
// Java: DominatorTree.compute(DiGraph<N, ?>, N)
func DominatorTree_Compute[N comparable, E any](graph DiGraph[N, E], entry N) *DominatorTree[N] {
	jutil.CheckNotNullMsg(graph.GetDiGraphNode(entry), "Entry node not in graph")

	// 1. Post-order traversal to assign indexes to nodes.
	postOrder := []N{}
	postOrderIndex := jutil.NewLinkedHashMap[N, int]()
	dominatorTree_buildPostOrder(graph, entry, &postOrder, jutil.NewLinkedHashMap[N, bool]())

	numNodes := len(postOrder)
	reversePostOrder := dominatorTree_toReversedArray(postOrder)

	// Map each node to its index in the post-order list.
	for i := 0; i < numNodes; i++ {
		postOrderIndex.Put(postOrder[i], i)
	}

	// 2. Iterative algorithm for immediate dominators.
	idomIndexes := make([]int, numNodes)
	for i := range idomIndexes {
		idomIndexes[i] = -1
	}

	startNodeIndex := postOrderIndex.Get(entry)
	idomIndexes[startNodeIndex] = startNodeIndex

	changed := true
	for changed {
		changed = false
		// Process nodes in reverse post-order (excluding entry node).
		for _, node := range reversePostOrder {
			if node == entry {
				continue
			}

			nodeIdx := postOrderIndex.Get(node)
			newIdomIdx := -1

			// Pick first processed predecessor.
			for _, pred := range graph.GetDirectedPredNodesOfValue(node) {
				predValue := pred.GetValue()
				predIdx, ok := postOrderIndex.Lookup(predValue)
				// predIdx can be null if `predValue` is not in the graph. This
				// is normal - it just means the `predValue` isn't reachable from
				// the given entry point.
				if ok && idomIndexes[predIdx] != -1 {
					newIdomIdx = predIdx
					break
				}
			}

			if newIdomIdx != -1 {
				// Intersect with all other processed predecessors.
				for _, pred := range graph.GetDirectedPredNodesOfValue(node) {
					predValue := pred.GetValue()
					predIdx, ok := postOrderIndex.Lookup(predValue)
					// predIdx can be null if `predValue` is not in the graph.
					// This is normal - it just means the `predValue` isn't
					// reachable from the given entry point.
					if ok && idomIndexes[predIdx] != -1 {
						newIdomIdx = dominatorTree_intersect(idomIndexes, predIdx, newIdomIdx)
					}
				}

				if idomIndexes[nodeIdx] != newIdomIdx {
					idomIndexes[nodeIdx] = newIdomIdx
					changed = true
				}
			}
		}
	}

	// 3. Build the result maps.
	idomsMap := jutil.NewLinkedHashMap[N, N]()
	// ListMultimap<N, N> children = ArrayListMultimap.create(); only get()
	// is used, so a map of lists suffices.
	children := jutil.NewLinkedHashMap[N, []N]()
	for i := 0; i < numNodes; i++ {
		node := postOrder[i]
		idomIdx := idomIndexes[i]
		if idomIdx != -1 && i != startNodeIndex {
			idom := postOrder[idomIdx]
			idomsMap.Put(node, idom)
			children.Put(idom, append(children.Get(idom), node))
		}
	}

	sizesMap := jutil.NewLinkedHashMap[N, int]()
	dominatorTree_computeSubtreeSizes(entry, children, sizesMap)

	return newDominatorTree(idomsMap, sizesMap)
}

// Java: DominatorTree.intersect(int[], int, int)
func dominatorTree_intersect(idomIndexes []int, i int, j int) int {
	finger1 := i
	finger2 := j
	for finger1 != finger2 {
		for finger1 < finger2 {
			finger1 = idomIndexes[finger1]
		}
		for finger2 < finger1 {
			finger2 = idomIndexes[finger2]
		}
	}
	return finger1
}

// Java: DominatorTree.buildPostOrder(DiGraph<N, ?>, N, List<N>, Map<N, Boolean>)
func dominatorTree_buildPostOrder[N comparable, E any](
	graph DiGraph[N, E], node N, postOrder *[]N, visited *jutil.LinkedHashMap[N, bool]) {
	visited.Put(node, true)
	for _, succ := range graph.GetDirectedSuccNodesOfValue(node) {
		succValue := succ.GetValue()
		if !visited.ContainsKey(succValue) {
			dominatorTree_buildPostOrder(graph, succValue, postOrder, visited)
		}
	}
	*postOrder = append(*postOrder, node)
}

// Java: DominatorTree.computeSubtreeSizes(N, ListMultimap<N, N>, Map<N, Integer>)
func dominatorTree_computeSubtreeSizes[N comparable](
	node N, children *jutil.LinkedHashMap[N, []N], sizes *jutil.LinkedHashMap[N, int]) int {
	size := 1
	nodeChildren := children.Get(node)
	for _, child := range nodeChildren {
		size += dominatorTree_computeSubtreeSizes(child, children, sizes)
	}
	sizes.Put(node, size)
	return size
}

// GetImmediateDominator returns node's immediate dominator.
//
// Go-only: language-adapter — Java returns null for the entry node and for
// values not reachable from it; here ok is false then.
//
// Java: DominatorTree.getImmediateDominator(N)
func (t *DominatorTree[N]) GetImmediateDominator(node N) (idom N, ok bool) {
	return t.idoms.Lookup(node)
}

// Java: DominatorTree.getDominatedSubtreeSize(N)
func (t *DominatorTree[N]) GetDominatedSubtreeSize(node N) int {
	return t.subtreeSizes.GetOrDefault(node, 0)
}

// GetAllSubtreeSizes returns the subtree sizes in Java ImmutableMap order.
// Callers must not mutate the map (Java: ImmutableMap).
//
// Java: DominatorTree.getAllSubtreeSizes()
func (t *DominatorTree[N]) GetAllSubtreeSizes() *jutil.LinkedHashMap[N, int] {
	return t.subtreeSizes
}

// Java: DominatorTree.toReversedArray(List<N>)
func dominatorTree_toReversedArray[N any](list []N) []N {
	size := len(list)
	reversed := make([]N, size)
	for i := 0; i < size; i++ {
		reversed[i] = list[size-1-i]
	}
	return reversed
}
