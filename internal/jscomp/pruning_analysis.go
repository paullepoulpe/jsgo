// Port of com/google/javascript/jscomp/PruningAnalysis.java @ closure-compiler v20261021

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

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// PruningAnalysis performs analysis on pruned dependencies to identify which
// root requirements are pulling in the most transitive files.
//
// This is intended to help developers understand the size of their
// dependency graph better, and make decisions about how to optimize
// compilation performance by reducing the number of inputs.
type PruningAnalysis struct {
	sorter      *SortedDependencies[*CompilerInput]
	entryPoints *jutil.LinkedHashSet[*CompilerInput]
}

// Limit the size of the "bottleneck" report to the top N files.
//
// Java: PruningAnalysis.BOTTLENECK_THRESHOLD
const pruningAnalysis_BOTTLENECK_THRESHOLD = 25

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
	entryPointDependencyCount *jutil.LinkedHashMap[string, int]
	bottleneckBlame           *jutil.LinkedHashMap[string, int]
}

// Java: PruningAnalysis.Result(ImmutableMap<String, Integer>, ImmutableMap<String, Integer>)
func NewPruningAnalysis_Result(
	entryPointDependencyCount *jutil.LinkedHashMap[string, int],
	bottleneckBlame *jutil.LinkedHashMap[string, int]) *PruningAnalysis_Result {
	return &PruningAnalysis_Result{
		entryPointDependencyCount: entryPointDependencyCount,
		bottleneckBlame:           bottleneckBlame,
	}
}

// Java: PruningAnalysis.Result.entryPointDependencyCount()
func (r *PruningAnalysis_Result) EntryPointDependencyCount() *jutil.LinkedHashMap[string, int] {
	return r.entryPointDependencyCount
}

// Java: PruningAnalysis.Result.bottleneckBlame()
func (r *PruningAnalysis_Result) BottleneckBlame() *jutil.LinkedHashMap[string, int] {
	return r.bottleneckBlame
}

// Java: PruningAnalysis(SortedDependencies<CompilerInput>, Collection<CompilerInput>)
func newPruningAnalysis(
	sorter *SortedDependencies[*CompilerInput], entryPoints []*CompilerInput) *PruningAnalysis {
	return &PruningAnalysis{
		sorter:      sorter,
		entryPoints: jutil.NewLinkedHashSet(entryPoints...),
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
func (a *PruningAnalysis) CalculateTotalBlame() *jutil.LinkedHashMap[string, int] {
	reachabilityPerEntryPoint := jutil.NewImmutableMapBuilder[string, int]()
	for provider := range a.entryPoints.All() {
		transitiveClosure := jutil.NewLinkedHashSet[*CompilerInput]()
		transitiveClosure.AddAll(
			a.sorter.GetStrongDependenciesOf([]*CompilerInput{provider} /* sorted= */, false)...)
		transitiveClosure.AddAll(a.sorter.GetSortedWeakDependenciesOf([]*CompilerInput{provider})...)
		reachabilityPerEntryPoint.Put(provider.GetName(), transitiveClosure.Size())
	}
	return reachabilityPerEntryPoint.BuildOrThrow()
}

// Java: PruningAnalysis.VIRTUAL_ROOT
const pruningAnalysis_VIRTUAL_ROOT = "ENTRY POINTS ROOT"

// CalculateBottlenecks calculates the top N files that are considered
// "bottlenecks" in the dependency graph, along with how many files are
// pulled in by that bottleneck.
//
// Java: PruningAnalysis.calculateBottlenecks()
func (a *PruningAnalysis) CalculateBottlenecks() *jutil.LinkedHashMap[string, int] {
	graph := a.BuildGraph()

	dominatorTree := DominatorTree_Compute[string, *string](graph, pruningAnalysis_VIRTUAL_ROOT)

	bottleneckBlame := jutil.NewLinkedHashMap[string, int]()
	for node, retainedCount := range dominatorTree.GetAllSubtreeSizes().All() {
		if node != pruningAnalysis_VIRTUAL_ROOT && retainedCount > 1 {
			bottleneckBlame.Put(node, retainedCount)
		}
	}

	// bottleneckBlame.entrySet().stream()
	//     .sorted(comparingInt(Map.Entry::getValue).reversed())
	//     .limit(BOTTLENECK_THRESHOLD)
	//     .collect(toImmutableMap(Map.Entry::getKey, Map.Entry::getValue))
	type entry struct {
		key   string
		value int
	}
	var entries []entry
	for k, v := range bottleneckBlame.All() {
		entries = append(entries, entry{k, v})
	}
	// Stream.sorted is stable for ordered streams.
	slices.SortStableFunc(entries, func(x, y entry) int { return cmp.Compare(y.value, x.value) })
	if len(entries) > pruningAnalysis_BOTTLENECK_THRESHOLD {
		entries = entries[:pruningAnalysis_BOTTLENECK_THRESHOLD]
	}
	result := jutil.NewImmutableMapBuilder[string, int]()
	for _, e := range entries {
		result.Put(e.key, e.value)
	}
	return result.BuildOrThrow()
}

// Java: PruningAnalysis.buildGraph()
func (a *PruningAnalysis) BuildGraph() DiGraph[string, *string] {
	graph := LinkedDirectedGraph_Create[string, *string]()
	graph.CreateNode(pruningAnalysis_VIRTUAL_ROOT)

	for _, input := range a.sorter.GetSortedList() {
		graph.CreateNode(input.GetName())
	}

	for input := range a.entryPoints.All() {
		graph.Connect(pruningAnalysis_VIRTUAL_ROOT, nil, input.GetName())
	}

	for _, input := range a.sorter.GetSortedList() {
		name := input.GetName()
		depSymbols := jutil.NewLinkedHashSet[string]()
		for _, req := range input.GetRequires() {
			depSymbols.Add(req.GetSymbol())
		}
		depSymbols.AddAll(input.GetTypeRequires()...)

		for symbol := range depSymbols.All() {
			jutil.Try(func() {
				dep := a.sorter.GetInputProviding(symbol)
				if dep != input {
					graph.ConnectIfNotFound(name, nil, dep.GetName())
				}
			}, func(e ISortedDependencies_MissingProvideException) {
				// Skip missing symbols.
			})
		}
	}
	return graph
}
