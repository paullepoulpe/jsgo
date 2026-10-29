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

	"github.com/paullepoulpe/jsgo/internal/graph"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// JSChunkGraph is a JSChunk dependency graph that assigns a depth to each
// chunk and can answer depth-related queries about them. For the purposes of
// this class, a chunk's depth is defined as the number of hops in the longest
// (non cyclic) path from the chunk to a chunk with no dependencies.
type JSChunkGraph struct {
	chunks []*JSChunk

	// selfPlusTransitiveDeps[i] = indices of all chunks that chunks[i]
	// depends on, including itself.
	selfPlusTransitiveDeps []jutil.BitSet

	// subtreeSize[i] = Number of chunks that transitively depend on
	// chunks[i], including itself.
	subtreeSize []int

	// Lists of chunks at each depth. chunksByDepth.get(3) is a list of the
	// chunks at depth 3, for example.
	chunksByDepth [][]*JSChunk

	// transitiveDeps[i] caches the deepest-first transitive dependencies of
	// chunks[i] (excluding chunks[i] itself). A nil entry means not yet cached.
	transitiveDeps [][]*JSChunk

	// Go-only: reusable candidate BitSet for GetSmallestCoveringSubtree.
	scratchCandidates jutil.BitSet
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
	g.chunksByDepth = g.initChunksByDepth()

	// Determine transitive deps for all chunks.
	// O(n*m * log(n)) (probably a bit better than that)
	g.selfPlusTransitiveDeps = g.initTransitiveDepsBitSets()

	// O(n*m)
	g.subtreeSize = g.initSubtreeSize()

	// Move all sources marked as weak by outside sources (e.g. flags) into the weak chunk.
	jsChunkGraph_moveMarkedWeakSources(g.GetChunkByName(JSChunk_WEAK_CHUNK_NAME), g.GetAllInputs())
	return g
}

// Java: JSChunkGraph.initChunksByDepth()
func (g *JSChunkGraph) initChunksByDepth() [][]*JSChunk {
	tmpChunksByDepth := [][]*JSChunk{}
	for chunkIndex := 0; chunkIndex < len(g.chunks); chunkIndex++ {
		chunk := g.chunks[chunkIndex]
		jutil.CheckStateMsg1(chunk.GetDepth() == -1, "Chunk depth already set: %s", chunk)
		depth := 0
		for _, dep := range chunk.GetDependencies() {
			depDepth := dep.GetDepth()
			if depDepth < 0 {
				panic(NewJSChunkGraph_ChunkDependenceException(
					"Chunks not in dependency order: "+chunk.GetName()+" preceded "+dep.GetName(),
					chunk,
					dep))
			}
			depth = max(depth, depDepth+1)
		}

		chunk.SetDepth(depth)
		if depth == len(tmpChunksByDepth) {
			tmpChunksByDepth = append(tmpChunksByDepth, []*JSChunk{})
		}
		tmpChunksByDepth[depth] = append(tmpChunksByDepth[depth], chunk)
	}
	return tmpChunksByDepth
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
func (g *JSChunkGraph) initTransitiveDepsBitSets() []jutil.BitSet {
	array := make([]jutil.BitSet, len(g.chunks))
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

// jsChunkGraph_moveMarkedWeakSources moves all sources that have
// SourceKind.WEAK into the weak chunk so that they may be pruned later.
//
// Java: JSChunkGraph.moveMarkedWeakSources(JSChunk, Iterable<CompilerInput>)
func jsChunkGraph_moveMarkedWeakSources(weakChunk *JSChunk, inputs []*CompilerInput) {
	jutil.CheckNotNull(weakChunk)
	allInputs := append([]*CompilerInput{}, inputs...)
	for _, i := range allInputs {
		if i.GetSourceFile().IsWeak() {
			existingChunk := i.GetChunk()
			if existingChunk == weakChunk {
				continue
			}
			if existingChunk != nil {
				existingChunk.Remove(i)
			}
			weakChunk.Add(i)
		}
	}
}

// ManageDependencies applies the dependency options to the list of sources,
// returning a new source list re-ordering and dropping files as necessary.
// This chunk graph will be updated to reflect the new list.
//
// See DependencyOptions for more information on how this works.
//
// It panics with *SortedDependencies_MissingProvideException if an entry
// point was not provided by any of the inputs, and with
// *JSChunkGraph_MissingChunkException if an entry point names an unknown
// chunk.
//
// Java: JSChunkGraph.manageDependencies(AbstractCompiler, DependencyOptions)
func (g *JSChunkGraph) ManageDependencies(
	compiler AbstractCompiler, dependencyOptions *DependencyOptions) *JSChunkGraph_DependencyManagementResult {

	// Make a copy since we're going to mutate the graph below.
	originalInputs := g.GetAllInputs()

	inputIndex := make(map[*CompilerInput]int, len(originalInputs))
	for _, input := range originalInputs {
		if _, exists := inputIndex[input]; !exists {
			inputIndex[input] = len(inputIndex)
		}
	}
	numOriginalInputs := len(inputIndex)

	sorter := NewSortedDependencies[*CompilerInput](originalInputs)

	entryPointInputs := g.createEntryPointInputs(compiler, dependencyOptions, originalInputs, sorter, inputIndex)

	// Build a map of symbols to their source file(s). While having multiple source files is invalid
	// we leave that up to typechecking so that we avoid arbitrarily picking a file.
	inputsByProvide := make(map[string][]*CompilerInput, len(originalInputs)*2)
	addProvide := func(symbol string, input *CompilerInput) {
		existing := inputsByProvide[symbol]
		if !slices.Contains(existing, input) {
			inputsByProvide[symbol] = append(existing, input)
		}
	}
	for _, input := range originalInputs {
		for _, provide := range input.GetKnownProvides() {
			addProvide(provide, input)
		}
		addProvide(input.GetPath().ToModuleName(), input)
	}

	// Dynamically imported files must be added to the chunk graph, but
	// they should not be ordered ahead of the files that import them.
	// We add them as entry points to ensure they get included.
	for _, input := range originalInputs {
		for _, require := range input.GetDynamicRequires() {
			entryPointInputs.addAll(inputsByProvide[require])
		}
	}

	// For goog.requireDynamic() imported files.
	for _, input := range originalInputs {
		for _, require := range input.GetRequireDynamicImports() {
			entryPointInputs.addAll(inputsByProvide[require])
		}
	}

	// The order of inputs, sorted independently of chunks.
	absoluteOrder := sorter.GetStrongDependenciesOf(originalInputs, dependencyOptions.ShouldSort())

	// Figure out which sources *must* be in each chunk.
	entryPointInputsPerChunk := make([][]*CompilerInput, len(g.chunks))
	for _, input := range entryPointInputs.items {
		chunk := input.GetChunk()
		jutil.CheckNotNull(chunk)
		chunkIdx := chunk.GetIndex()
		jutil.CheckState(chunkIdx >= 0 && chunkIdx < len(g.chunks) && g.chunks[chunkIdx] == chunk)
		entryPointInputsPerChunk[chunkIdx] = append(entryPointInputsPerChunk[chunkIdx], input)
	}

	// Clear the chunks of their inputs. This also nulls out the input's reference to its chunk.
	for _, chunk := range g.GetAllChunks() {
		chunk.RemoveAll()
	}

	// Figure out which sources *must* be in each chunk, or in one
	// of that chunk's dependencies.
	orderedInputs := []*CompilerInput{}
	sortAndPrune := dependencyOptions.ShouldSort() && dependencyOptions.ShouldPrune()
	var reachedInputs []bool
	var unreachedInputs []bool
	var scratchClosure []*CompilerInput
	if sortAndPrune {
		reachedInputs = make([]bool, len(inputIndex))
		unreachedInputs = make([]bool, len(inputIndex))
	}

	for chunkIdx, chunk := range g.chunks {
		var transitiveClosure []*CompilerInput
		// Prefer a depth first ordering of dependencies from entry points.
		// Always orders in a deterministic fashion regardless of the order of provided inputs
		// given the same entry points in the same order.
		if sortAndPrune {
			transitiveClosure = scratchClosure[:0]
			chunkEntryPoints := entryPointInputsPerChunk[chunkIdx]
			if len(chunkEntryPoints) > 0 {
				// We need the full set of dependencies for each chunk, so start with the full input set
				for i := range unreachedInputs {
					unreachedInputs[i] = i < numOriginalInputs
				}
				for _, entryPoint := range chunkEntryPoints {
					transitiveClosure = g.appendDepthFirstDependenciesOf(
						transitiveClosure, entryPoint, unreachedInputs, inputIndex, inputsByProvide)
				}
				scratchClosure = transitiveClosure
				// For any input we have not yet reached, add them to the ordered list
				for _, orderedInput := range transitiveClosure {
					idx := inputIndex[orderedInput]
					if !reachedInputs[idx] {
						reachedInputs[idx] = true
						orderedInputs = append(orderedInputs, orderedInput)
					}
				}
			}
		} else {
			// Simply order inputs so that any required namespace comes before it's usage.
			// Ordered result varies based on the original order of inputs.
			chunkEntryPoints := entryPointInputsPerChunk[chunkIdx]
			if chunkEntryPoints == nil {
				chunkEntryPoints = []*CompilerInput{}
			}
			transitiveClosure =
				sorter.GetStrongDependenciesOf(chunkEntryPoints, dependencyOptions.ShouldSort())
		}
		for _, input := range transitiveClosure {
			if dependencyOptions.ShouldPrune() &&
				input.GetSourceFile().IsWeak() &&
				!entryPointInputs.contains(input) {
				compiler.Report(
					JSError_Make(
						JSChunkGraph_WEAK_FILE_REACHABLE_FROM_ENTRY_POINT_ERROR, input.GetSourceFile().GetName()))
			}
			oldChunk := input.GetChunk()
			if oldChunk == nil {
				input.SetChunk(chunk)
			} else {
				input.SetChunk(nil)
				input.SetChunk(g.GetDeepestCommonDependencyInclusive(oldChunk, chunk))
			}
		}
	}
	if !sortAndPrune || len(entryPointInputs.items) == 0 {
		orderedInputs = absoluteOrder
	}

	weakChunk := g.GetChunkByName(JSChunk_WEAK_CHUNK_NAME)
	jutil.CheckNotNull(weakChunk)
	// Mark all sources that are detected as weak.
	if dependencyOptions.ShouldPrune() {
		weakInputs := sorter.GetSortedWeakDependenciesOf(orderedInputs)
		for _, i := range weakInputs {
			// Add weak inputs to the weak chunk in dependency order. moveMarkedWeakSources will move
			// in command line flag order.
			jutil.CheckState(i.GetChunk() == nil)
			i.GetSourceFile().SetKind(rhino.StaticSourceFile_SourceKind_WEAK)
			i.SetChunk(weakChunk)
			weakChunk.Add(i)
		}
	} else {
		// Only move sourced marked as weak if the compiler isn't doing its own detection.
		jsChunkGraph_moveMarkedWeakSources(weakChunk, originalInputs)
	}

	// All the inputs are pointing to the chunks that own them. Yeah!
	// Update the chunks to reflect this.
	for _, input := range orderedInputs {
		chunk := input.GetChunk()
		if chunk != nil && chunk.GetByName(input.GetName()) == nil {
			chunk.Add(input)
		}
	}

	return NewJSChunkGraph_DependencyManagementResult(
		slices.Clone(orderedInputs), sorter, jutil.NewLinkedHashSet(entryPointInputs.items...))
}

// JSChunkGraph_DependencyManagementResult is the result of
// JSChunkGraph.manageDependencies.
//
// Java: JSChunkGraph.DependencyManagementResult
type JSChunkGraph_DependencyManagementResult struct {
	orderedInputs    []*CompilerInput
	sorter           *SortedDependencies[*CompilerInput]
	entryPointInputs *jutil.LinkedHashSet[*CompilerInput]
}

// Java: JSChunkGraph.DependencyManagementResult(ImmutableList<CompilerInput>, SortedDependencies<CompilerInput>, ImmutableSet<CompilerInput>)
func NewJSChunkGraph_DependencyManagementResult(
	orderedInputs []*CompilerInput,
	sorter *SortedDependencies[*CompilerInput],
	entryPointInputs *jutil.LinkedHashSet[*CompilerInput]) *JSChunkGraph_DependencyManagementResult {
	return &JSChunkGraph_DependencyManagementResult{
		orderedInputs:    orderedInputs,
		sorter:           sorter,
		entryPointInputs: entryPointInputs,
	}
}

// Java: JSChunkGraph.DependencyManagementResult.orderedInputs()
func (r *JSChunkGraph_DependencyManagementResult) OrderedInputs() []*CompilerInput {
	return r.orderedInputs
}

// Java: JSChunkGraph.DependencyManagementResult.sorter()
func (r *JSChunkGraph_DependencyManagementResult) Sorter() *SortedDependencies[*CompilerInput] {
	return r.sorter
}

// Java: JSChunkGraph.DependencyManagementResult.entryPointInputs()
func (r *JSChunkGraph_DependencyManagementResult) EntryPointInputs() *jutil.LinkedHashSet[*CompilerInput] {
	return r.entryPointInputs
}

// appendDepthFirstDependenciesOf appends rootInput and its unreached strong
// dependencies in post-order depth-first traversal order onto dst.
//
// Java: JSChunkGraph.getDepthFirstDependenciesOf(CompilerInput, Set<CompilerInput>, Map<String, Set<CompilerInput>>)
func (g *JSChunkGraph) appendDepthFirstDependenciesOf(
	dst []*CompilerInput,
	rootInput *CompilerInput,
	unreachedInputs []bool,
	inputIndex map[*CompilerInput]int,
	inputsByProvide map[string][]*CompilerInput) []*CompilerInput {
	rootIdx, ok := inputIndex[rootInput]
	if !ok || !unreachedInputs[rootIdx] {
		return dst
	}
	unreachedInputs[rootIdx] = false

	for _, importedNamespace := range rootInput.GetRequiredSymbols() {
		for _, input := range inputsByProvide[importedNamespace] {
			if idx, ok := inputIndex[input]; ok && unreachedInputs[idx] {
				dst = g.appendDepthFirstDependenciesOf(dst, input, unreachedInputs, inputIndex, inputsByProvide)
			}
		}
	}

	return append(dst, rootInput)
}

type entryPointInputSet struct {
	inputIndex map[*CompilerInput]int
	present    []bool
	items      []*CompilerInput
}

func newEntryPointInputSet(inputIndex map[*CompilerInput]int) *entryPointInputSet {
	return &entryPointInputSet{
		inputIndex: inputIndex,
		present:    make([]bool, len(inputIndex)),
	}
}

func (s *entryPointInputSet) add(input *CompilerInput) {
	idx, ok := s.inputIndex[input]
	if !ok {
		idx = len(s.inputIndex)
		s.inputIndex[input] = idx
		s.present = append(s.present, false)
	}
	if !s.present[idx] {
		s.present[idx] = true
		s.items = append(s.items, input)
	}
}

func (s *entryPointInputSet) addAll(inputs []*CompilerInput) {
	for _, input := range inputs {
		s.add(input)
	}
}

func (s *entryPointInputSet) contains(input *CompilerInput) bool {
	idx, ok := s.inputIndex[input]
	return ok && idx < len(s.present) && s.present[idx]
}

// Java: JSChunkGraph.createEntryPointInputs(AbstractCompiler, DependencyOptions, Iterable<CompilerInput>, SortedDependencies<CompilerInput>)
func (g *JSChunkGraph) createEntryPointInputs(
	compiler AbstractCompiler,
	dependencyOptions *DependencyOptions,
	inputs []*CompilerInput,
	sorter *SortedDependencies[*CompilerInput],
	inputIndex map[*CompilerInput]int) *entryPointInputSet {
	entryPointInputs := newEntryPointInputSet(inputIndex)
	if dependencyOptions.ShouldPrune() {
		// Some files implicitly depend on base.js without actually requiring anything.
		// So we always treat it as the first entry point to ensure it's ordered correctly.
		baseJs := sorter.MaybeGetInputProviding("goog")
		if baseJs != nil {
			entryPointInputs.add(baseJs)
		}

		if !dependencyOptions.ShouldDropMoochers() {
			for _, entryPointInput := range sorter.GetInputsWithoutProvides() {
				if entryPointInput.GetSourceFile().IsWeak() {
					compiler.Report(
						JSError_Make(
							JSChunkGraph_IMPLICIT_WEAK_ENTRY_POINT_ERROR, entryPointInput.GetSourceFile().GetName()))
				} else {
					entryPointInputs.add(entryPointInput)
				}
			}
		}

		var chunksByName map[string]*JSChunk
		for _, entryPoint := range dependencyOptions.EntryPoints() {
			var entryPointInput *CompilerInput
			jutil.Try(func() {
				if entryPoint.ClosureNamespace() == entryPoint.ModuleName() {
					entryPointInput = sorter.MaybeGetInputProviding(entryPoint.ClosureNamespace())
					// Check to see if we can find the entry point as an ES6 and CommonJS module
					// ES6 and CommonJS entry points may not provide any symbols
					if entryPointInput == nil {
						entryPointInput = sorter.GetInputProviding(entryPoint.Name())
					}
				} else {
					if chunksByName == nil {
						chunksByName = make(map[string]*JSChunk, len(g.chunks))
						for _, m := range g.chunks {
							chunksByName[m.GetName()] = m
						}
					}
					chunk := chunksByName[entryPoint.ModuleName()]
					if chunk == nil {
						panic(NewJSChunkGraph_MissingChunkException(entryPoint.ModuleName()))
					} else {
						entryPointInput = sorter.GetInputProviding(entryPoint.ClosureNamespace())
						entryPointInput.OverrideChunk(chunk)
					}
				}
			}, func(e ISortedDependencies_MissingProvideException) {
				panic(NewSortedDependencies_MissingProvideExceptionWithCause(entryPoint.Name(), e))
			})

			if entryPointInput.GetSourceFile().IsWeak() {
				compiler.Report(
					JSError_Make(
						JSChunkGraph_EXPLICIT_WEAK_ENTRY_POINT_ERROR, entryPointInput.GetSourceFile().GetName()))
			} else {
				entryPointInputs.add(entryPointInput)
			}
		}
	} else {
		entryPointInputs.addAll(inputs)
	}
	return entryPointInputs
}

// Java: JSChunkGraph.toGraphvizGraph()
func (g *JSChunkGraph) ToGraphvizGraph() *graph.LinkedDirectedGraph[*JSChunk, string] {
	graphViz := graph.LinkedDirectedGraph_Create[*JSChunk, string]()
	for _, chunk := range g.GetAllChunks() {
		graphViz.CreateNode(chunk)
		for _, dep := range chunk.GetDependencies() {
			graphViz.CreateNode(dep)
			graphViz.Connect(chunk, "->", dep)
		}
	}
	return graphViz
}

// JSChunkGraph_InverseDepthComparator is a chunk depth comparator that
// considers a deeper chunk to be "less than" a shallower chunk. Uses chunk
// names to consistently break ties.
type JSChunkGraph_InverseDepthComparator struct {
	// Go-only: language-adapter — distinct identity
	_ byte
}

// Java: JSChunkGraph.InverseDepthComparator.INSTANCE
var JSChunkGraph_InverseDepthComparator_INSTANCE = &JSChunkGraph_InverseDepthComparator{}

// Java: JSChunkGraph.InverseDepthComparator.compare(JSChunk, JSChunk)
func (c *JSChunkGraph_InverseDepthComparator) Compare(m1 *JSChunk, m2 *JSChunk) int {
	return jsChunkGraph_depthCompare(m2, m1)
}

// Java: JSChunkGraph.depthCompare(JSChunk, JSChunk)
func jsChunkGraph_depthCompare(m1 *JSChunk, m2 *JSChunk) int {
	if m1 == m2 {
		return 0
	}
	d1 := m1.GetDepth()
	d2 := m2.GetDepth()
	if d1 < d2 {
		return -1
	} else if d2 == d1 {
		return jutil.StringCompareTo(m1.GetName(), m2.GetName())
	}
	return 1
}

// JSChunkGraph_ChunkDependenceException is an exception class for declaring
// when the chunks being fed into a JSChunkGraph as input aren't in
// dependence order, and so can't be processed for caching of various
// dependency-related queries.
type JSChunkGraph_ChunkDependenceException struct {
	jutil.IllegalArgumentException

	chunk          *JSChunk
	dependentChunk *JSChunk
}

func (*JSChunkGraph_ChunkDependenceException) chunkDependenceException() {}

// IJSChunkGraph_ChunkDependenceException matches
// JSChunkGraph_ChunkDependenceException.
type IJSChunkGraph_ChunkDependenceException interface {
	jutil.IIllegalArgumentException
	chunkDependenceException()
}

// Java: JSChunkGraph.ChunkDependenceException.serialVersionUID
const jsChunkGraph_ChunkDependenceException_serialVersionUID int64 = 1

// Java: JSChunkGraph.ChunkDependenceException(String, JSChunk, JSChunk)
func NewJSChunkGraph_ChunkDependenceException(
	message string, chunk *JSChunk, dependentChunk *JSChunk) *JSChunkGraph_ChunkDependenceException {
	e := &JSChunkGraph_ChunkDependenceException{chunk: chunk, dependentChunk: dependentChunk}
	jutil.InitThrowable(e, "com.google.javascript.jscomp.JSChunkGraph$ChunkDependenceException", &message, nil)
	return e
}

// Java: JSChunkGraph.ChunkDependenceException.getChunk()
func (e *JSChunkGraph_ChunkDependenceException) GetChunk() *JSChunk {
	return e.chunk
}

// Java: JSChunkGraph.ChunkDependenceException.getDependentChunk()
func (e *JSChunkGraph_ChunkDependenceException) GetDependentChunk() *JSChunk {
	return e.dependentChunk
}

// JSChunkGraph_MissingChunkException is another exception class (a checked
// exception).
type JSChunkGraph_MissingChunkException struct{ jutil.JavaException }

func (*JSChunkGraph_MissingChunkException) missingChunkException() {}

// IJSChunkGraph_MissingChunkException matches
// JSChunkGraph_MissingChunkException.
type IJSChunkGraph_MissingChunkException interface {
	jutil.IException
	missingChunkException()
}

// Java: JSChunkGraph.MissingChunkException(String)
func NewJSChunkGraph_MissingChunkException(chunkName string) *JSChunkGraph_MissingChunkException {
	e := &JSChunkGraph_MissingChunkException{}
	jutil.InitThrowable(e, "com.google.javascript.jscomp.JSChunkGraph$MissingChunkException", &chunkName, nil)
	return e
}
