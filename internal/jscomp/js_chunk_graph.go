// Port of com/google/javascript/jscomp/JSChunkGraph.java @ closure-compiler v20261021

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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// JSChunkGraph is a JSChunk dependency graph that assigns a depth to each
// chunk and can answer depth-related queries about them. For the purposes of
// this class, a chunk's depth is defined as the number of hops in the longest
// (non cyclic) path from the chunk to a chunk with no dependencies.
type JSChunkGraph struct {
	chunks []*JSChunk

	// selfPlusTransitiveDeps[i] = indices of all chunks that chunks[i]
	// depends on, including itself.
	selfPlusTransitiveDeps []inlineBitSet

	// subtreeSize[i] = Number of chunks that transitively depend on
	// chunks[i], including itself.
	subtreeSize []int

	// Lists of chunks at each depth. chunksByDepth.get(3) is a list of the
	// chunks at depth 3, for example.
	chunksByDepth [][]*JSChunk

	// transitiveDeps[i] caches the deepest-first transitive dependencies of
	// chunks[i] (excluding chunks[i] itself). A nil entry means not yet cached.
	transitiveDeps [][]*JSChunk

	// Go-only: reusable candidate/dependent BitSets for GetSmallestCoveringSubtree.
	scratchCandidates inlineBitSet
	scratchDependent  inlineBitSet
}

// Java: JSChunkGraph.WEAK_FILE_REACHABLE_FROM_ENTRY_POINT_ERROR
var JSChunkGraph_WEAK_FILE_REACHABLE_FROM_ENTRY_POINT_ERROR = DiagnosticType_Error(
	"JSC_WEAK_FILE_REACHABLE_FROM_ENTRY_POINT_ERROR",
	"File strongly reachable from an entry point must not be weak: {0}")

// Java: JSChunkGraph.EXPLICIT_WEAK_ENTRY_POINT_ERROR
var JSChunkGraph_EXPLICIT_WEAK_ENTRY_POINT_ERROR = DiagnosticType_Error(
	"JSC_EXPLICIT_WEAK_ENTRY_POINT_ERROR",
	"Explicit entry point input must not be weak: {0}")

// Java: JSChunkGraph.IMPLICIT_WEAK_ENTRY_POINT_ERROR
var JSChunkGraph_IMPLICIT_WEAK_ENTRY_POINT_ERROR = DiagnosticType_Warning(
	"JSC_IMPLICIT_WEAK_ENTRY_POINT_ERROR",
	"Implicit entry point input should not be weak: {0}")

// NewJSChunkGraphFromArray creates a chunk graph from a list of chunks in
// dependency order.
//
// Java: JSChunkGraph(JSChunk[])
func NewJSChunkGraphFromArray(chunksInDepOrder []*JSChunk) *JSChunkGraph {
	return NewJSChunkGraph(chunksInDepOrder)
}

// NewJSChunkGraph creates a chunk graph from a list of chunks in dependency
// order.
//
// Java: JSChunkGraph(List<JSChunk>)
func NewJSChunkGraph(chunksInDepOrder []*JSChunk) *JSChunkGraph {
	g, err := tryNewJSChunkGraph(chunksInDepOrder)
	if err != nil {
		panic(err)
	}
	return g
}

func tryNewJSChunkGraph(chunksInDepOrder []*JSChunk) (*JSChunkGraph, *JSChunkGraph_ChunkDependenceException) {
	g := &JSChunkGraph{}
	jutil.CheckState(len(chunksInDepOrder) != 0)
	chunksInDepOrder = g.makeWeakChunk(chunksInDepOrder)
	g.chunks = make([]*JSChunk, len(chunksInDepOrder))

	// n = number of chunks
	// Populate chunks O(n)
	for chunkIndex := 0; chunkIndex < len(g.chunks); chunkIndex++ {
		chunk := chunksInDepOrder[chunkIndex]
		jutil.CheckStateMsg1(chunk.GetIndex() == -1, "Chunk index already set: %s", chunk)
		chunk.SetIndex(chunkIndex)
		g.chunks[chunkIndex] = chunk
	}

	// Determine depth for all chunks.
	// m = number of edges in the graph
	// O(n*m)
	chunksByDepth, err := g.tryInitChunksByDepth()
	if err != nil {
		return nil, err
	}
	g.chunksByDepth = chunksByDepth

	// Determine transitive deps for all chunks.
	// O(n*m * log(n)) (probably a bit better than that)
	g.selfPlusTransitiveDeps = g.initTransitiveDepsBitSets()

	// O(n*m)
	g.subtreeSize = g.initSubtreeSize()

	// Move all sources marked as weak by outside sources (e.g. flags) into the weak chunk.
	jsChunkGraph_moveMarkedWeakSources(g.GetChunkByName(JSChunk_WEAK_CHUNK_NAME), g.GetAllInputs())
	return g, nil
}

// Java: JSChunkGraph.initChunksByDepth()
func (g *JSChunkGraph) initChunksByDepth() [][]*JSChunk {
	chunksByDepth, err := g.tryInitChunksByDepth()
	if err != nil {
		panic(err)
	}
	return chunksByDepth
}

func (g *JSChunkGraph) tryInitChunksByDepth() ([][]*JSChunk, *JSChunkGraph_ChunkDependenceException) {
	tmpChunksByDepth := [][]*JSChunk{}
	for chunkIndex := 0; chunkIndex < len(g.chunks); chunkIndex++ {
		chunk := g.chunks[chunkIndex]
		jutil.CheckStateMsg1(chunk.GetDepth() == -1, "Chunk depth already set: %s", chunk)
		depth := 0
		for _, dep := range chunk.GetDependencies() {
			depDepth := dep.GetDepth()
			if depDepth < 0 {
				return nil, NewJSChunkGraph_ChunkDependenceException(
					"Chunks not in dependency order: "+chunk.GetName()+" preceded "+dep.GetName(),
					chunk,
					dep)
			}
			depth = max(depth, depDepth+1)
		}

		chunk.SetDepth(depth)
		if depth == len(tmpChunksByDepth) {
			tmpChunksByDepth = append(tmpChunksByDepth, []*JSChunk{})
		}
		tmpChunksByDepth[depth] = append(tmpChunksByDepth[depth], chunk)
	}
	return tmpChunksByDepth, nil
}

// makeWeakChunk: if a weak chunk doesn't already exist, creates a weak chunk
// depending on every other chunk.
//
// Does not move any sources into the weak chunk.
//
// Returns a new list of chunks that includes the weak chunk, if it was newly
// created, or the same list if the weak chunk already existed. It panics
// with IllegalStateException if a weak chunk already exists but doesn't
// fulfill the above conditions
//
// Java: JSChunkGraph.makeWeakChunk(List<JSChunk>)
func (g *JSChunkGraph) makeWeakChunk(chunksInDepOrder []*JSChunk) []*JSChunk {
	hasWeakChunk := false
	for _, chunk := range chunksInDepOrder {
		if chunk.GetName() == JSChunk_WEAK_CHUNK_NAME {
			hasWeakChunk = true
			weakDeps := make(map[*JSChunk]struct{}, len(chunksInDepOrder))
			stack := slices.Clone(chunk.GetDependencies())
			for len(stack) > 0 {
				dep := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if _, seen := weakDeps[dep]; !seen {
					weakDeps[dep] = struct{}{}
					stack = append(stack, dep.GetDependencies()...)
				}
			}
			otherChunks := make(map[*JSChunk]struct{}, len(chunksInDepOrder))
			containsAll := true
			for _, other := range chunksInDepOrder {
				if other == chunk {
					continue
				}
				otherChunks[other] = struct{}{}
				if _, ok := weakDeps[other]; !ok {
					containsAll = false
					break
				}
			}
			jutil.CheckStateMsg(
				containsAll,
				"A weak chunk already exists but it does not depend on every other chunk.")
			jutil.CheckStateMsg(
				len(weakDeps) == len(otherChunks),
				"The weak chunk cannot have extra dependencies.")
			break
		}
	}
	if hasWeakChunk {
		// All weak files (and only weak files) should be in the weak chunk.
		misplacedWeakFiles := []string{}
		misplacedStrongFiles := []string{}
		for _, chunk := range chunksInDepOrder {
			isWeakChunk := chunk.GetName() == JSChunk_WEAK_CHUNK_NAME
			for _, input := range chunk.GetInputs() {
				if isWeakChunk && !input.GetSourceFile().IsWeak() {
					misplacedStrongFiles = append(misplacedStrongFiles, input.GetSourceFile().GetName())
				} else if !isWeakChunk && input.GetSourceFile().IsWeak() {
					misplacedWeakFiles = append(misplacedWeakFiles,
						input.GetSourceFile().GetName()+" (in chunk "+chunk.GetName()+")")
				}
			}
		}
		if !(len(misplacedStrongFiles) == 0 && len(misplacedWeakFiles) == 0) {
			var sb strings.Builder
			sb.WriteString("A weak chunk exists but some sources are misplaced.")
			if len(misplacedStrongFiles) != 0 {
				sb.WriteString("\nFound these strong sources in the weak chunk:\n  ")
				sb.WriteString(strings.Join(misplacedStrongFiles, "\n  "))
			}
			if len(misplacedWeakFiles) != 0 {
				sb.WriteString("\nFound these weak sources in other chunks:\n  ")
				sb.WriteString(strings.Join(misplacedWeakFiles, "\n  "))
			}
			panic(jutil.NewIllegalStateExceptionMsg(sb.String()))
		}
	} else {
		weakChunk := NewJSChunk(JSChunk_WEAK_CHUNK_NAME)
		for _, chunk := range chunksInDepOrder {
			weakChunk.AddDependency(chunk)
		}
		chunksInDepOrder = append(append([]*JSChunk{}, chunksInDepOrder...), weakChunk)
	}
	return chunksInDepOrder
}

// Java: JSChunkGraph.initTransitiveDepsBitSets()
func (g *JSChunkGraph) initTransitiveDepsBitSets() []inlineBitSet {
	array := make([]inlineBitSet, len(g.chunks))
	for chunkIndex := 0; chunkIndex < len(g.chunks); chunkIndex++ {
		chunk := g.chunks[chunkIndex]
		selfPlusTransitiveDeps := &array[chunkIndex]
		selfPlusTransitiveDeps.Set(chunkIndex)
		// O(chunkIndex * log64(chunkIndex))
		for _, dep := range chunk.GetDependencies() {
			// Add this dependency and all of its dependencies to the current chunk.
			// O(log64(chunkIndex))
			selfPlusTransitiveDeps.Or(&array[dep.GetIndex()])
		}
	}
	return array
}

// Java: JSChunkGraph.initSubtreeSize()
func (g *JSChunkGraph) initSubtreeSize() []int {
	subtreeSize := make([]int, len(g.chunks))
	for dependentIndex := 0; dependentIndex < len(g.chunks); dependentIndex++ {
		dependencies := &g.selfPlusTransitiveDeps[dependentIndex]
		// Iterating backward through the bitset is slightly more efficient, since it avoids
		// considering later chunks, which this one cannot depend on.
		for requiredIndex := dependentIndex; requiredIndex >= 0; requiredIndex = dependencies.PreviousSetBit(requiredIndex - 1) {
			subtreeSize[requiredIndex] += 1 // Count dependent in required chunk's subtree.
		}
	}
	return subtreeSize
}

// GetAllInputs gets an iterable over all input source files in dependency
// order.
//
// Java: JSChunkGraph.getAllInputs()
func (g *JSChunkGraph) GetAllInputs() []*CompilerInput {
	out := []*CompilerInput{}
	for _, chunk := range g.chunks {
		out = append(out, chunk.GetInputs()...)
	}
	return out
}

// GetInputCount gets the total number of input source files.
//
// Java: JSChunkGraph.getInputCount()
func (g *JSChunkGraph) GetInputCount() int {
	count := 0
	for _, chunk := range g.chunks {
		count += chunk.GetInputCount()
	}
	return count
}

// GetAllChunks gets an iterable over all chunks in dependency order.
//
// Java: JSChunkGraph.getAllChunks()
func (g *JSChunkGraph) GetAllChunks() []*JSChunk {
	return g.chunks
}

// GetIntegralChunkArrayForSerialization gets all chunks in dependency order.
// Note: Do not modify the returned array.
//
// Java: JSChunkGraph.getIntegralChunkArrayForSerialization()
func (g *JSChunkGraph) GetIntegralChunkArrayForSerialization() []*JSChunk {
	return g.chunks
}

// GetChunkByName gets a single chunk by name; nil if no such chunk exists.
//
// Java: JSChunkGraph.getChunkByName(String)
func (g *JSChunkGraph) GetChunkByName(name string) *JSChunk {
	if n := len(g.chunks); n > 0 && g.chunks[n-1].GetName() == name {
		return g.chunks[n-1]
	}
	for _, m := range g.chunks {
		if m.GetName() == name {
			return m
		}
	}
	return nil
}

// GetChunksByName gets all chunks indexed by name.
//
// Java: JSChunkGraph.getChunksByName()
func (g *JSChunkGraph) GetChunksByName() *jutil.LinkedHashMap[string, *JSChunk] {
	result := jutil.NewLinkedHashMap[string, *JSChunk]()
	for _, m := range g.chunks {
		result.Put(m.GetName(), m)
	}
	return result
}

// GetChunkCount gets the total number of chunks.
//
// Java: JSChunkGraph.getChunkCount()
func (g *JSChunkGraph) GetChunkCount() int {
	return len(g.chunks)
}

// GetRootChunk gets the root chunk.
//
// Java: JSChunkGraph.getRootChunk()
func (g *JSChunkGraph) GetRootChunk() *JSChunk {
	return colorPool_getOnlyElement(g.chunksByDepth[0])
}

// ToJson returns a JSON representation of the JSChunkGraph. Specifically a
// JsonArray of "Chunks" where each chunk has a - "name" - "dependencies"
// (list of chunk names) - "transitive-dependencies" (list of chunk names,
// deepest first) - "inputs" (list of file names)
//
// Java: JSChunkGraph.toJson()
func (g *JSChunkGraph) ToJson() *jutil.JsonArray {
	chunks := jutil.NewJsonArray()
	for _, chunk := range g.GetAllChunks() {
		node := jutil.NewJsonObject()
		node.Add("name", jutil.NewJsonPrimitiveString(chunk.GetName()))
		deps := jutil.NewJsonArray()
		node.Add("dependencies", deps)
		for _, m := range chunk.GetDependencies() {
			deps.Add(jutil.NewJsonPrimitiveString(m.GetName()))
		}
		transitiveDeps := jutil.NewJsonArray()
		node.Add("transitive-dependencies", transitiveDeps)
		for _, m := range g.getTransitiveDepsDeepestFirstCached(chunk) {
			transitiveDeps.Add(jutil.NewJsonPrimitiveString(m.GetName()))
		}
		inputs := jutil.NewJsonArray()
		node.Add("inputs", inputs)
		for _, input := range chunk.GetInputs() {
			inputs.Add(jutil.NewJsonPrimitiveString(input.GetSourceFile().GetName()))
		}
		chunks.Add(node)
	}
	return chunks
}

// DependsOn determines whether this chunk depends on a given chunk. Note
// that a chunk never depends on itself, as that dependency would be cyclic.
//
// Java: JSChunkGraph.dependsOn(JSChunk, JSChunk)
func (g *JSChunkGraph) DependsOn(src *JSChunk, m *JSChunk) bool {
	return src != m && g.selfPlusTransitiveDeps[src.GetIndex()].Get(m.GetIndex())
}

// GetSmallestCoveringSubtree finds the chunk with the fewest transitive
// dependents on which all of the given chunks depend and that is a subtree
// of the given parent chunk tree.
//
// If no such subtree can be found, the parent chunk is returned.
//
// If multiple candidates have the same number of dependents, the chunk
// farthest down in the total ordering of chunks will be chosen.
//
// parentTree is the chunk on which the result must depend; dependentChunks
// are the indices of chunks to consider. The result is a chunk on which all
// of the argument chunks depend
//
// Java: JSChunkGraph.getSmallestCoveringSubtree(JSChunk, BitSet)
func (g *JSChunkGraph) GetSmallestCoveringSubtree(parentTree *JSChunk, dependentChunks *jutil.BitSet) *JSChunk {
	jutil.CheckState(!dependentChunks.IsEmpty())
	g.scratchDependent.CopyFromJutil(dependentChunks)
	return g.getSmallestCoveringSubtree(parentTree, &g.scratchDependent)
}

func (g *JSChunkGraph) getSmallestCoveringSubtree(parentTree *JSChunk, dependentChunks *inlineBitSet) *JSChunk {
	jutil.CheckState(!dependentChunks.IsEmpty())

	// Candidate chunks are those that all of the given dependent chunks depend on, including
	// themselves. The dependent chunk with the smallest index might be our answer, if all
	// the other chunks depend on it.
	minDependentChunkIndex := dependentChunks.NextSetBit(0)
	nextDependent := dependentChunks.NextSetBit(minDependentChunkIndex + 1)
	if nextDependent < 0 {
		// Fast path for a single dependent chunk: by transitivity of selfPlusTransitiveDeps,
		// if minDependentChunkIndex depends on parentTree, it is itself the smallest covering
		// subtree; otherwise no chunk in its dependency set depends on parentTree.
		if g.selfPlusTransitiveDeps[minDependentChunkIndex].Get(parentTree.GetIndex()) {
			return g.chunks[minDependentChunkIndex]
		}
		return parentTree
	}

	candidates := &g.scratchCandidates
	candidates.CopyFrom(&g.selfPlusTransitiveDeps[minDependentChunkIndex])
	for dependentIndex := nextDependent; dependentIndex >= 0; dependentIndex = dependentChunks.NextSetBit(dependentIndex + 1) {
		candidates.And(&g.selfPlusTransitiveDeps[dependentIndex])
	}
	jutil.CheckStateMsg1(!candidates.IsEmpty(), "No common dependency found for %s", dependentChunks)

	// All candidates must have an index <= the smallest dependent chunk index.
	// Work backwards through the candidates starting with the dependent chunk with the smallest
	// index. For each candidate, we'll remove all of the chunks it depends on from consideration,
	// since they must all have larger subtrees than the one we're considering.
	parentTreeIndex := parentTree.GetIndex()
	// default to parent tree if we don't find anything better
	bestCandidateIndex := parentTreeIndex
	for candidateIndex := candidates.PreviousSetBit(minDependentChunkIndex); candidateIndex >= 0; candidateIndex = candidates.PreviousSetBit(candidateIndex - 1) {

		candidatePlusTransitiveDeps := &g.selfPlusTransitiveDeps[candidateIndex]
		if candidatePlusTransitiveDeps.Get(parentTreeIndex) {
			// candidate is a subtree of parentTree
			candidates.AndNot(candidatePlusTransitiveDeps)
			if g.subtreeSize[candidateIndex] < g.subtreeSize[bestCandidateIndex] {
				bestCandidateIndex = candidateIndex
			}
		} // eliminate candidates that are not a subtree of parentTree
	}
	return g.chunks[bestCandidateIndex]
}

// GetDeepestCommonDependency finds the deepest common dependency of two
// chunks, not including the two chunks themselves.
//
// m1 and m2 are chunks in this graph. The result is the deepest common dep
// of m1 and m2, or nil if they have no common dependencies
//
// Java: JSChunkGraph.getDeepestCommonDependency(JSChunk, JSChunk)
func (g *JSChunkGraph) GetDeepestCommonDependency(m1 *JSChunk, m2 *JSChunk) *JSChunk {
	m1Depth := m1.GetDepth()
	m2Depth := m2.GetDepth()
	// According our definition of depth, the result must have a strictly
	// smaller depth than either m1 or m2.
	for depth := min(m1Depth, m2Depth) - 1; depth >= 0; depth-- {
		chunksAtDepth := g.chunksByDepth[depth]
		// Look at the chunks at this depth in reverse order, so that we use the
		// original ordering of the chunks to break ties (later meaning deeper).
		for i := len(chunksAtDepth) - 1; i >= 0; i-- {
			m := chunksAtDepth[i]
			if g.DependsOn(m1, m) && g.DependsOn(m2, m) {
				return m
			}
		}
	}
	return nil
}

// GetDeepestCommonDependencyInclusive finds the deepest common dependency of
// two chunks, including the chunks themselves.
//
// m1 and m2 are chunks in this graph. The result is the deepest common dep
// of m1 and m2, or nil if they have no common dependencies
//
// Java: JSChunkGraph.getDeepestCommonDependencyInclusive(JSChunk, JSChunk)
func (g *JSChunkGraph) GetDeepestCommonDependencyInclusive(m1 *JSChunk, m2 *JSChunk) *JSChunk {
	if m2 == m1 || g.DependsOn(m2, m1) {
		return m1
	} else if g.DependsOn(m1, m2) {
		return m2
	}

	return g.GetDeepestCommonDependency(m1, m2)
}

// GetDeepestCommonDependencyInclusiveOfCollection returns the deepest common
// dependency of the given chunks.
//
// Java: JSChunkGraph.getDeepestCommonDependencyInclusive(Collection<JSChunk>)
func (g *JSChunkGraph) GetDeepestCommonDependencyInclusiveOfCollection(chunks []*JSChunk) *JSChunk {
	if len(chunks) == 0 {
		// Iterator.next() on an empty collection.
		panic(jutil.NewNoSuchElementException())
	}
	dep := chunks[0]
	for _, next := range chunks[1:] {
		dep = g.GetDeepestCommonDependencyInclusive(dep, next)
	}
	return dep
}

// GetTransitiveDepsDeepestFirst creates an iterable over the transitive
// dependencies of chunk m in a non-increasing depth ordering. The result
// does not include the chunk m.
//
// m is a chunk in this graph. The result is the transitive dependencies of
// chunk m
//
// Java: JSChunkGraph.getTransitiveDepsDeepestFirst(JSChunk)
func (g *JSChunkGraph) GetTransitiveDepsDeepestFirst(m *JSChunk) []*JSChunk {
	return slices.Clone(g.getTransitiveDepsDeepestFirstCached(m))
}

// getTransitiveDepsDeepestFirstCached returns the cached deepest-first
// transitive dependencies of chunk m. Callers must not mutate the returned slice.
func (g *JSChunkGraph) getTransitiveDepsDeepestFirstCached(m *JSChunk) []*JSChunk {
	idx := m.GetIndex()
	if idx >= 0 && idx < len(g.chunks) && g.chunks[idx] == m {
		if g.transitiveDeps == nil {
			g.transitiveDeps = make([][]*JSChunk, len(g.chunks))
		}
		if cached := g.transitiveDeps[idx]; cached != nil {
			return cached
		}
		bits := &g.selfPlusTransitiveDeps[idx]
		deps := []*JSChunk{}
		for i := bits.NextSetBit(0); i >= 0; i = bits.NextSetBit(i + 1) {
			if i != idx {
				deps = append(deps, g.chunks[i])
			}
		}
		slices.SortStableFunc(deps, JSChunkGraph_InverseDepthComparator_INSTANCE.Compare)
		g.transitiveDeps[idx] = deps
		return deps
	}

	// Fallback if m is not indexed in this graph.
	seen := make(map[*JSChunk]struct{})
	seen[m] = struct{}{}
	deps := []*JSChunk{}
	var visit func(cur *JSChunk)
	visit = func(cur *JSChunk) {
		for _, dep := range cur.GetDependencies() {
			if _, ok := seen[dep]; !ok {
				seen[dep] = struct{}{}
				deps = append(deps, dep)
				visit(dep)
			}
		}
	}
	visit(m)
	slices.SortStableFunc(deps, JSChunkGraph_InverseDepthComparator_INSTANCE.Compare)
	return deps
}
