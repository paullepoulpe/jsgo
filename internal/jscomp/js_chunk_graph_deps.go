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

	"github.com/paullepoulpe/jsgo/internal/graph"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

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
			if entryPoint.ClosureNamespace() == entryPoint.ModuleName() {
				entryPointInput = sorter.MaybeGetInputProviding(entryPoint.ClosureNamespace())
				// Check to see if we can find the entry point as an ES6 and CommonJS module
				// ES6 and CommonJS entry points may not provide any symbols
				if entryPointInput == nil {
					entryPointInput = sorter.MaybeGetInputProviding(entryPoint.Name())
					if entryPointInput == nil {
						panic(NewSortedDependencies_MissingProvideExceptionWithCause(
							entryPoint.Name(), NewSortedDependencies_MissingProvideException(entryPoint.Name())))
					}
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
				}
				entryPointInput = sorter.MaybeGetInputProviding(entryPoint.ClosureNamespace())
				if entryPointInput == nil {
					panic(NewSortedDependencies_MissingProvideExceptionWithCause(
						entryPoint.Name(), NewSortedDependencies_MissingProvideException(entryPoint.ClosureNamespace())))
				}
				entryPointInput.OverrideChunk(chunk)
			}

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

// Java: JSChunkGraph.MissingChunkException(String)
func NewJSChunkGraph_MissingChunkException(chunkName string) *JSChunkGraph_MissingChunkException {
	e := &JSChunkGraph_MissingChunkException{}
	jutil.InitThrowable(e, "com.google.javascript.jscomp.JSChunkGraph$MissingChunkException", &chunkName, nil)
	return e
}
