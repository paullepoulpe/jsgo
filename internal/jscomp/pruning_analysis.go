// Port of com/google/javascript/jscomp/PruningAnalysis.java @ closure-compiler v20261104

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
	"cmp"
	"slices"

	"github.com/paullepoulpe/jsgo/internal/graph"
)

// PruningCountMap is an insertion-ordered map of file names to dependency counts.
type PruningCountMap = OrderedStringMap[int]

// PruningAnalysis performs analysis on pruned dependencies to identify which
// root requirements are pulling in the most transitive files.
//
// This is intended to help developers understand the size of their
// dependency graph better, and make decisions about how to optimize
// compilation performance by reducing the number of inputs.
type PruningAnalysis struct {
	sorter      *SortedDependencies[*CompilerInput]
	entryPoints OrderedSet[*CompilerInput]
}

// Limit the size of the "bottleneck" report to the top N files.
//
// Java: PruningAnalysis.BOTTLENECK_THRESHOLD
const pruningBottleneckThreshold = 25

// PruningAnalysis_Result holds the results of the pruning analysis.
//
// entryPointDependencyCount is a map of entry point file name to the number
// of files that are pulled in by that root, including shared dependencies.
// bottleneckBlame is a map of the top N files that are considered
// "bottlenecks" in the dependency graph, along with how many files are
// pulled in by that bottleneck (i.e. how many files are "dominated" by that
// bottleneck in a dominator tree analysis.)
//
// Java: record PruningAnalysis.Result
type PruningAnalysis_Result struct {
	entryPointDependencyCount *PruningCountMap
	bottleneckBlame           *PruningCountMap
}

// Java: PruningAnalysis.Result(ImmutableMap<String, Integer>, ImmutableMap<String, Integer>)
func NewPruningAnalysis_Result(
	entryPointDependencyCount *PruningCountMap,
	bottleneckBlame *PruningCountMap) *PruningAnalysis_Result {
	return &PruningAnalysis_Result{
		entryPointDependencyCount: entryPointDependencyCount,
		bottleneckBlame:           bottleneckBlame,
	}
}

// Java: PruningAnalysis.Result.entryPointDependencyCount()
func (r *PruningAnalysis_Result) EntryPointDependencyCount() *PruningCountMap {
	return r.entryPointDependencyCount
}

// Java: PruningAnalysis.Result.bottleneckBlame()
func (r *PruningAnalysis_Result) BottleneckBlame() *PruningCountMap {
	return r.bottleneckBlame
}

// Java: PruningAnalysis(SortedDependencies<CompilerInput>, Collection<CompilerInput>)
func newPruningAnalysis(
	sorter *SortedDependencies[*CompilerInput], entryPoints []*CompilerInput) *PruningAnalysis {
	return &PruningAnalysis{
		sorter:      sorter,
		entryPoints: *NewOrderedSet(entryPoints...),
	}
}

// Java: PruningAnalysis.create(SortedDependencies<CompilerInput>, Collection<CompilerInput>)
func PruningAnalysis_Create(
	sorter *SortedDependencies[*CompilerInput], entryPoints []*CompilerInput) *PruningAnalysis {
	return newPruningAnalysis(sorter, entryPoints)
}

// Analyze analyzes the dependency graph to identify possible areas of code
// reduction by pruning dependencies.
//
// Java: PruningAnalysis.analyze()
func (a *PruningAnalysis) Analyze() *PruningAnalysis_Result {
	total := a.CalculateTotalBlame()
	bottlenecks := a.CalculateBottlenecks()

	return NewPruningAnalysis_Result(total, bottlenecks)
}

// Java: PruningAnalysis.calculateTotalBlame()
func (a *PruningAnalysis) CalculateTotalBlame() *PruningCountMap {
	reachabilityPerEntryPoint := NewOrderedStringMapWithCapacity[int](a.entryPoints.Size())
	seen := make(map[*CompilerInput]struct{})
	for provider := range a.entryPoints.All() {
		clear(seen)
		for _, in := range a.sorter.GetStrongDependenciesOf([]*CompilerInput{provider} /* sorted= */, false) {
			seen[in] = struct{}{}
		}
		for _, in := range a.sorter.GetSortedWeakDependenciesOf([]*CompilerInput{provider}) {
			seen[in] = struct{}{}
		}
		reachabilityPerEntryPoint.Put(provider.GetName(), len(seen))
	}
	return reachabilityPerEntryPoint
}

// Java: PruningAnalysis.VIRTUAL_ROOT
const pruningVirtualRoot = "ENTRY POINTS ROOT"

// CalculateBottlenecks calculates the top N files that are considered
// "bottlenecks" in the dependency graph, along with how many files are
// pulled in by that bottleneck.
//
// Java: PruningAnalysis.calculateBottlenecks()
func (a *PruningAnalysis) CalculateBottlenecks() *PruningCountMap {
	depGraph := a.buildGraph()

	dominatorTree := graph.DominatorTree_Compute[string, struct{}](depGraph, pruningVirtualRoot)

	// bottleneckBlame.entrySet().stream()
	//     .sorted(comparingInt(Map.Entry::getValue).reversed())
	//     .limit(BOTTLENECK_THRESHOLD)
	//     .collect(toImmutableMap(Map.Entry::getKey, Map.Entry::getValue))
	type entry struct {
		key   string
		value int
	}
	var entries []entry
	for node, retainedCount := range dominatorTree.GetAllSubtreeSizes().All() {
		if node != pruningVirtualRoot && retainedCount > 1 {
			entries = append(entries, entry{node, retainedCount})
		}
	}
	// Stream.sorted is stable for ordered streams.
	slices.SortStableFunc(entries, func(x, y entry) int { return cmp.Compare(y.value, x.value) })
	if len(entries) > pruningBottleneckThreshold {
		entries = entries[:pruningBottleneckThreshold]
	}
	result := NewOrderedStringMapWithCapacity[int](len(entries))
	for _, e := range entries {
		result.Put(e.key, e.value)
	}
	return result
}

// Java: PruningAnalysis.buildGraph()
func (a *PruningAnalysis) buildGraph() graph.DiGraph[string, struct{}] {
	depGraph := graph.LinkedDirectedGraph_Create[string, struct{}]()
	depGraph.CreateNode(pruningVirtualRoot)

	sortedInputs := a.sorter.GetSortedList()
	for _, input := range sortedInputs {
		depGraph.CreateNode(input.GetName())
	}

	for input := range a.entryPoints.All() {
		depGraph.Connect(pruningVirtualRoot, struct{}{}, input.GetName())
	}

	var depSymbols OrderedSet[string]
	for _, input := range sortedInputs {
		name := input.GetName()
		depSymbols.Clear()
		for _, req := range input.GetRequires() {
			depSymbols.Add(req.GetSymbol())
		}
		depSymbols.AddAll(input.GetTypeRequires()...)

		for symbol := range depSymbols.All() {
			if dep := a.sorter.MaybeGetInputProviding(symbol); dep != nil && dep != input {
				depGraph.ConnectIfNotFound(name, struct{}{}, dep.GetName())
			}
		}
	}
	return depGraph
}
