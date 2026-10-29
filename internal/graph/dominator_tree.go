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

package graph

import (
	"iter"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// DominatorSubtreeSizes stores dominated subtree sizes in deterministic post-order.
type DominatorSubtreeSizes[N comparable] struct {
	sizes map[N]int
	order []N
}

// Size returns the number of entries.
func (s *DominatorSubtreeSizes[N]) Size() int { return len(s.order) }

// Lookup returns the subtree size for k and whether k is present.
func (s *DominatorSubtreeSizes[N]) Lookup(k N) (int, bool) {
	v, ok := s.sizes[k]
	return v, ok
}

// All iterates (node, subtreeSize) pairs in insertion order.
func (s *DominatorSubtreeSizes[N]) All() iter.Seq2[N, int] {
	return func(yield func(N, int) bool) {
		for _, k := range s.order {
			if !yield(k, s.sizes[k]) {
				return
			}
		}
	}
}

func (s *DominatorSubtreeSizes[N]) put(k N, size int) {
	if _, exists := s.sizes[k]; !exists {
		s.order = append(s.order, k)
	}
	s.sizes[k] = size
}

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
	// entry node is omitted.
	idoms map[N]N

	// Mapping from node value to the number of nodes it dominates
	// (including itself).
	subtreeSizes DominatorSubtreeSizes[N]
}

// Java: DominatorTree(Map<N, N>, Map<N, Integer>)
func newDominatorTree[N comparable](idoms map[N]N, subtreeSizes DominatorSubtreeSizes[N]) *DominatorTree[N] {
	return &DominatorTree[N]{
		idoms:        idoms,
		subtreeSizes: subtreeSizes,
	}
}

// DominatorTree_Compute computes the dominator tree for the given graph
// starting from the entry node.
//
// graph is the directed graph to analyze; entry is the entry node of the
// graph. Returns a DominatorTree instance.
//
// Java: DominatorTree.compute(DiGraph<N, ?>, N)
func DominatorTree_Compute[N, E comparable](graph DiGraph[N, E], entry N) *DominatorTree[N] {
	ldg := graph.AsLinkedDirectedGraph()
	entryNode := ldg.GetNode(entry)
	jutil.CheckNotNullMsg(entryNode, "Entry node not in graph")

	// 1. Post-order traversal indexed directly by LinkedDiGraphNode.Index:
	// 0 = unvisited/unreachable, -1 = in-progress on DFS stack, >0 = postOrderIdx + 1.
	poPlus1 := make([]int32, ldg.NodeIndexCap())
	postOrder := make([]*LinkedDiGraphNode[N, E], 0, min(ldg.GetNodeCount(), 64))
	dominatorTree_buildPostOrder(entryNode, &postOrder, poPlus1)

	numNodes := len(postOrder)

	// 2. Iterative algorithm for immediate dominators + flat child lists.
	scratch := make([]int32, 3*numNodes)
	for i := range scratch {
		scratch[i] = -1
	}
	idomIndexes := scratch[:numNodes]
	firstChild := scratch[numNodes : 2*numNodes]
	nextSibling := scratch[2*numNodes : 3*numNodes]

	startNodeIndex := int(poPlus1[entryNode.Index] - 1)
	idomIndexes[startNodeIndex] = int32(startNodeIndex)

	changed := true
	for changed {
		changed = false
		// Process nodes in reverse post-order (excluding entry node).
		for nodeIdx := numNodes - 1; nodeIdx >= 0; nodeIdx-- {
			if nodeIdx == startNodeIndex {
				continue
			}

			gNode := postOrder[nodeIdx]
			inEdges := gNode.InEdgeList.Slice()
			newIdomIdx := int32(-1)

			// Pick first processed predecessor.
			for _, edge := range inEdges {
				// predPo == 0 means edge.SourceNode isn't reachable from entry.
				if predPo := poPlus1[edge.SourceNode.Index]; predPo > 0 {
					predIdx := predPo - 1
					if idomIndexes[predIdx] != -1 {
						newIdomIdx = predIdx
						break
					}
				}
			}

			if newIdomIdx != -1 {
				// Intersect with all other processed predecessors.
				for _, edge := range inEdges {
					if predPo := poPlus1[edge.SourceNode.Index]; predPo > 0 {
						predIdx := predPo - 1
						if idomIndexes[predIdx] != -1 {
							newIdomIdx = dominatorTree_intersect(idomIndexes, predIdx, newIdomIdx)
						}
					}
				}

				if idomIndexes[nodeIdx] != newIdomIdx {
					idomIndexes[nodeIdx] = newIdomIdx
					changed = true
				}
			}
		}
	}

	// 3. Build the result maps, prepending children in descending i order so
	// child lists walk in exact ascending post-order.
	idomsMap := make(map[N]N, numNodes)
	for i := numNodes - 1; i >= 0; i-- {
		idomIdx := idomIndexes[i]
		if idomIdx != -1 && i != startNodeIndex {
			idomsMap[postOrder[i].Value] = postOrder[idomIdx].Value
			nextSibling[i] = firstChild[idomIdx]
			firstChild[idomIdx] = int32(i)
		}
	}

	sizesMap := DominatorSubtreeSizes[N]{
		sizes: make(map[N]int, numNodes),
		order: make([]N, 0, numNodes),
	}
	dominatorTree_computeSubtreeSizes(int32(startNodeIndex), postOrder, firstChild, nextSibling, &sizesMap)

	return newDominatorTree(idomsMap, sizesMap)
}

// Java: DominatorTree.intersect(int[], int, int)
func dominatorTree_intersect(idomIndexes []int32, i int32, j int32) int32 {
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
func dominatorTree_buildPostOrder[N, E comparable](
	gNode *LinkedDiGraphNode[N, E], postOrder *[]*LinkedDiGraphNode[N, E], poPlus1 []int32) {
	poPlus1[gNode.Index] = -1
	for _, edge := range gNode.OutEdgeList.Slice() {
		succ := edge.DestNode
		if poPlus1[succ.Index] == 0 {
			dominatorTree_buildPostOrder(succ, postOrder, poPlus1)
		}
	}
	idx := len(*postOrder)
	*postOrder = append(*postOrder, gNode)
	poPlus1[gNode.Index] = int32(idx + 1)
}

// Java: DominatorTree.computeSubtreeSizes(N, ListMultimap<N, N>, Map<N, Integer>)
func dominatorTree_computeSubtreeSizes[N, E comparable](
	nodeIdx int32,
	postOrder []*LinkedDiGraphNode[N, E],
	firstChild []int32,
	nextSibling []int32,
	sizes *DominatorSubtreeSizes[N],
) int {
	size := 1
	for child := firstChild[nodeIdx]; child != -1; child = nextSibling[child] {
		size += dominatorTree_computeSubtreeSizes(child, postOrder, firstChild, nextSibling, sizes)
	}
	sizes.put(postOrder[nodeIdx].Value, size)
	return size
}

// GetImmediateDominator returns node's immediate dominator.
//
// Go-only: language-adapter — Java returns null for the entry node and for
// values not reachable from it; here ok is false then.
//
// Java: DominatorTree.getImmediateDominator(N)
func (t *DominatorTree[N]) GetImmediateDominator(node N) (idom N, ok bool) {
	idom, ok = t.idoms[node]
	return idom, ok
}

// Java: DominatorTree.getDominatedSubtreeSize(N)
func (t *DominatorTree[N]) GetDominatedSubtreeSize(node N) int {
	return t.subtreeSizes.sizes[node]
}

// GetAllSubtreeSizes returns the subtree sizes in Java ImmutableMap order.
// Callers must not mutate the map (Java: ImmutableMap).
//
// Java: DominatorTree.getAllSubtreeSizes()
func (t *DominatorTree[N]) GetAllSubtreeSizes() *DominatorSubtreeSizes[N] {
	return &t.subtreeSizes
}
