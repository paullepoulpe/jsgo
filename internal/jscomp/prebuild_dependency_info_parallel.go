// Port of com/google/javascript/jscomp/PrebuildDependencyInfo.java @ closure-compiler v20261021

/*
 * Copyright 2019 The Closure Compiler Authors.
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
	"iter"
	"slices"
	"sync"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// PrebuildDependencyInfo is a helper class to prebuild DependencyInfo from
// a list of CompilerInput. Dependency info (requires, annotations, etc.) is
// parsed and generated the first time we try to get it for an input. Get
// them all now using multiple threads, so they'll be parsed in parallel and
// already available during the rest of the compilation.
type PrebuildDependencyInfo struct {
	numParallelThreads int
}

// Java: PrebuildDependencyInfo(int)
func NewPrebuildDependencyInfo(numParalleThreads int) *PrebuildDependencyInfo {
	return &PrebuildDependencyInfo{numParallelThreads: numParalleThreads}
}

func prebuildCallDependencyInfo(input *CompilerInput) (failure jutil.Throwable) {
	defer func() {
		if r := recover(); r != nil {
			if t, ok := r.(jutil.Throwable); ok {
				failure = t
				return
			}
			panic(r)
		}
	}()
	input.GetDependencyInfo()
	return nil
}

func prebuildRunTask(run func(int), index int) (failure jutil.Throwable) {
	defer func() {
		if r := recover(); r != nil {
			if t, ok := r.(jutil.Throwable); ok {
				failure = t
				return
			}
			panic(r)
		}
	}()
	run(index)
	return nil
}

func prebuildRunIndexedTask(run func(worker, index int), worker, index int) (failure jutil.Throwable) {
	defer func() {
		if r := recover(); r != nil {
			if t, ok := r.(jutil.Throwable); ok {
				failure = t
				return
			}
			panic(r)
		}
	}()
	run(worker, index)
	return nil
}

func prebuildPublishDiagnostics(events []prebuildDependencyDiagnostic) (failure jutil.Throwable) {
	if len(events) == 0 {
		return nil
	}
	defer func() {
		if r := recover(); r != nil {
			if t, ok := r.(jutil.Throwable); ok {
				failure = t
				return
			}
			panic(r)
		}
	}()
	for _, diagnostic := range events {
		diagnostic.target.Report(diagnostic.level, diagnostic.error)
	}
	return nil
}

// Prebuild gets the dependency info of every input.
//
// Go-only: language-adapter — dependency scans run in a bounded worker pool.
// Per-file diagnostics and dependency caches are published in input order.
// Unsupported compiler/source/resolver extensions retain the sequential path.
// Every task runs even after one fails, and the first failure is wrapped as in
// Java's Futures.allAsList(...).get().
//
// Java: PrebuildDependencyInfo.prebuild(Iterable<CompilerInput>)
func (p *PrebuildDependencyInfo) Prebuild(allInputs iter.Seq[*CompilerInput]) {
	if p.numParallelThreads > 1 {
		inputs := slices.Collect(allInputs)
		if p.prebuildParallel(inputs) {
			return
		}
		allInputs = slices.Values(inputs)
	}
	var failure jutil.Throwable
	// TODO(moz): Support canceling all parsing on the first halting error
	for input := range allInputs {
		if t := prebuildCallDependencyInfo(input); t != nil && failure == nil {
			failure = t
		}
	}

	if failure != nil {
		panic(jutil.NewRuntimeExceptionWithCause(jutil.NewExecutionException(failure)))
	}
}

// parallelPrebuild runs every submitted task, retaining failures in input order.
// Each slot is written by one worker and consumed only after the join.
func parallelPrebuild(count, workers int, run func(int)) []jutil.Throwable {
	return parallelPrebuildWorkers(count, workers, func(claim func() (int, bool), recordFailure func(int, jutil.Throwable)) {
		for {
			index, ok := claim()
			if !ok {
				return
			}
			if t := prebuildRunTask(run, index); t != nil {
				recordFailure(index, t)
			}
		}
	})
}

func parallelPrebuildIndexed(count, workers int, run func(worker, index int)) []jutil.Throwable {
	var nextWorker atomic.Int32
	return parallelPrebuildWorkers(count, workers, func(claim func() (int, bool), recordFailure func(int, jutil.Throwable)) {
		worker := int(nextWorker.Add(1) - 1)
		for {
			index, ok := claim()
			if !ok {
				return
			}
			if t := prebuildRunIndexedTask(run, worker, index); t != nil {
				recordFailure(index, t)
			}
		}
	})
}

func parallelPrebuildWorkers(
	count, workers int,
	workerFn func(claim func() (int, bool), recordFailure func(int, jutil.Throwable)),
) []jutil.Throwable {
	failures := make([]jutil.Throwable, count)
	var next atomic.Int64
	claim := func() (int, bool) {
		index := int(next.Add(1) - 1)
		if index >= count {
			return 0, false
		}
		return index, true
	}
	recordFailure := func(index int, t jutil.Throwable) {
		failures[index] = t
	}
	var wg sync.WaitGroup
	for range min(workers, count) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			workerFn(claim, recordFailure)
		}()
	}
	wg.Wait()
	return failures
}

type prebuildDependencyDiagnostic struct {
	target ErrorHandler
	level  *CheckLevel
	error  *JSError
}

type prebuildDependencyErrors struct {
	ErrorManager
	events []prebuildDependencyDiagnostic
}

func (r *prebuildDependencyErrors) Report(level *CheckLevel, error *JSError) {
	r.events = append(r.events, prebuildDependencyDiagnostic{r.ErrorManager, level, error})
}

// Resolver diagnostics normally pass through the loader's ErrorHandler (often
// Compiler.Report and its warning guards), not directly to the ErrorManager.
// Retain that route while buffering both diagnostic sources in encounter order.
type prebuildDependencyModuleErrors struct {
	target ErrorHandler
	events *[]prebuildDependencyDiagnostic
}

func (r *prebuildDependencyModuleErrors) Report(level *CheckLevel, error *JSError) {
	*r.events = append(*r.events, prebuildDependencyDiagnostic{r.target, level, error})
}

// Each scan gets its own diagnostic target. Resolver tables are immutable
// during prebuilding; copy only the resolver/loader headers, sharing those
// read-only tables. Custom resolver implementations use the serialized path.
func prebuildCloneModuleLoader(source *ModuleLoader, errors ErrorHandler) *ModuleLoader {
	loader := *source
	switch resolver := source.moduleResolver.(type) {
	case *BrowserModuleResolver:
		copy := *resolver
		loader.moduleResolver = &copy
	case *BrowserWithTransformedPrefixesModuleResolver:
		copy := *resolver
		loader.moduleResolver = &copy
	case *NodeModuleResolver:
		copy := *resolver
		loader.moduleResolver = &copy
	case *WebpackModuleResolver:
		copy := *resolver
		loader.moduleResolver = &copy
	default:
		return nil
	}
	loader.SetErrorHandler(errors)
	return &loader
}

func (p *PrebuildDependencyInfo) prebuildParallel(inputs []*CompilerInput) bool {
	if len(inputs) < 2 || inputs[0] == nil {
		return false
	}
	compiler, ok := inputs[0].compiler.(*Compiler)
	if !ok || compiler == nil || compiler.typedAstFilesystem != nil {
		return false
	}
	seen := make(map[*SourceFile]bool, len(inputs))
	for _, input := range inputs {
		if input == nil || input.compiler != compiler || input.ast == nil || seen[input.sourceFile] {
			return false
		}
		seen[input.sourceFile] = true
	}
	if !JsFileRegexParser_IsSupported() || !compiler.PreferRegexParser() {
		origPaths := make([]*ModuleLoader_ModulePath, len(inputs))
		var baseLoader *ModuleLoader
		for i, input := range inputs {
			orig := input.GetPath()
			origPaths[i] = orig
			if orig.outer != baseLoader {
				if prebuildCloneModuleLoader(orig.outer, nil) == nil {
					return false
				}
				baseLoader = orig.outer
			}
		}
		// Parsing has shared compiler side effects, so first publish independent
		// AST parses in order. Dependency tree walks then touch only their own
		// CompilerInput and its AST.
		if !NewPrebuildAst(compiler, p.numParallelThreads).prebuildParallel(inputs) {
			return false
		}
		errorsByInput := make([][]prebuildDependencyDiagnostic, len(inputs))
		failures := parallelPrebuildWorkers(len(inputs), p.numParallelThreads, func(claim func() (int, bool), recordFailure func(int, jutil.Throwable)) {
			var workerEvents []prebuildDependencyDiagnostic
			var workerModuleErrors prebuildDependencyModuleErrors
			workerModuleErrors.events = &workerEvents
			var cachedSource *ModuleLoader
			var workerLoader *ModuleLoader
			var privatePath ModuleLoader_ModulePath
			scanOne := func(i int) (failure jutil.Throwable) {
				defer func() {
					if r := recover(); r != nil {
						if t, ok := r.(jutil.Throwable); ok {
							failure = t
							return
						}
						panic(r)
					}
				}()
				input := inputs[i]
				orig := origPaths[i]
				workerEvents = nil
				workerModuleErrors.target = orig.outer.GetErrorHandler()
				if orig.outer != cachedSource {
					cachedSource = orig.outer
					workerLoader = prebuildCloneModuleLoader(orig.outer, &workerModuleErrors)
				}
				privatePath = *orig
				privatePath.outer = workerLoader
				input.modulePath = &privatePath
				defer func() {
					input.modulePath = orig
					if len(workerEvents) > 0 {
						errorsByInput[i] = workerEvents
					}
				}()
				input.GetDependencyInfo()
				return nil
			}
			for {
				i, ok := claim()
				if !ok {
					return
				}
				if t := scanOne(i); t != nil {
					recordFailure(i, t)
				}
			}
		})
		var failure jutil.Throwable
		for i := range inputs {
			if t := prebuildPublishDiagnostics(errorsByInput[i]); t != nil {
				if failure == nil {
					failure = t
				}
				continue
			}
			if failures[i] != nil && failure == nil {
				failure = failures[i]
			}
		}
		if failure != nil {
			panic(jutil.NewRuntimeExceptionWithCause(jutil.NewExecutionException(failure)))
		}
		return true
	}

	type outcome struct {
		info        DependencyInfo
		events      []prebuildDependencyDiagnostic
		readFailure bool
	}
	if prebuildCloneModuleLoader(compiler.GetModuleLoader(), nil) == nil {
		return false
	}
	results := make([]outcome, len(inputs))
	errorManager := compiler.GetErrorManager()
	moduleErrorTarget := compiler.GetModuleLoader().GetErrorHandler()
	failures := parallelPrebuildWorkers(len(inputs), p.numParallelThreads, func(claim func() (int, bool), recordFailure func(int, jutil.Throwable)) {
		workerErrors := prebuildDependencyErrors{ErrorManager: errorManager}
		workerModuleErrors := prebuildDependencyModuleErrors{
			target: moduleErrorTarget,
			events: &workerErrors.events,
		}
		loader := prebuildCloneModuleLoader(compiler.GetModuleLoader(), &workerModuleErrors)
		parser := NewJsFileRegexParser(&workerErrors).SetModuleLoader(loader).SetIncludeGoogBase(true)
		parseOne := func(i int) (failure jutil.Throwable) {
			defer func() {
				if r := recover(); r != nil {
					if t, ok := r.(jutil.Throwable); ok {
						failure = t
						return
					}
					panic(r)
				}
			}()
			input, result := inputs[i], &results[i]
			if input.dependencyInfo != nil {
				return nil
			}
			workerErrors.events = nil
			if code, err := input.getCode(); err != nil {
				result.readFailure = true
				workerErrors.Report(CheckLevel_ERROR.Ptr(),
					JSError_Make(AbstractCompiler_READ_ERROR, input.GetName(), ioErrorMessage(err)))
				result.info = SimpleDependencyInfo_EMPTY
			} else {
				result.info = parser.ParseFile(input.GetName(), input.GetName(), code)
			}
			result.events = workerErrors.events
			return nil
		}
		for {
			i, ok := claim()
			if !ok {
				return
			}
			if t := parseOne(i); t != nil {
				recordFailure(i, t)
			}
		}
	})
	var failure jutil.Throwable
	for i, input := range inputs {
		result := &results[i]
		if t := prebuildPublishDiagnostics(result.events); t != nil {
			if failure == nil {
				failure = t
			}
			continue
		}
		if failures[i] != nil {
			if failure == nil {
				failure = failures[i]
			}
			continue
		}
		if result.info != nil {
			// IO failures return EMPTY directly, matching generateDependencyInfo.
			if result.readFailure {
				input.dependencyInfo = result.info
			} else {
				input.dependencyInfo = NewLazyParsedDependencyInfo(result.info, input, compiler)
			}
		}
		if t := prebuildCallDependencyInfo(input); t != nil && failure == nil { // merge extra requires/provides in encounter order
			failure = t
		}
	}
	if failure != nil {
		panic(jutil.NewRuntimeExceptionWithCause(jutil.NewExecutionException(failure)))
	}
	return true
}
