package jscomp

import (
	"sync"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// parallelPrebuild runs every submitted task, retaining failures in input order.
// Each slot is written by one worker and consumed only after the join.
func parallelPrebuild(count, workers int, run func(int)) []jutil.Throwable {
	failures := make([]jutil.Throwable, count)
	work := make(chan int)
	var wg sync.WaitGroup
	for range min(workers, count) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range work {
				jutil.Try(func() { run(index) }, func(t jutil.Throwable) { failures[index] = t })
			}
		}()
	}
	for index := range count {
		work <- index
	}
	close(work)
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
		type dependencyPath struct {
			original     *ModuleLoader_ModulePath
			private      ModuleLoader_ModulePath
			errors       []prebuildDependencyDiagnostic
			moduleErrors prebuildDependencyModuleErrors
		}
		paths := make([]dependencyPath, len(inputs))
		for i, input := range inputs {
			path := &paths[i]
			path.original = input.GetPath()
			path.moduleErrors = prebuildDependencyModuleErrors{
				target: path.original.outer.GetErrorHandler(), events: &path.errors,
			}
			loader := prebuildCloneModuleLoader(path.original.outer, &path.moduleErrors)
			if loader == nil {
				return false
			}
			path.private = *path.original
			path.private.outer = loader
		}
		// Parsing has shared compiler side effects, so first publish independent
		// AST parses in order. Dependency tree walks then touch only their own
		// CompilerInput and its AST.
		if !NewPrebuildAst(compiler, p.numParallelThreads).prebuildParallel(inputs) {
			return false
		}
		failures := parallelPrebuild(len(inputs), p.numParallelThreads, func(i int) {
			input := inputs[i]
			input.modulePath = &paths[i].private
			defer func() { input.modulePath = paths[i].original }()
			inputs[i].GetDependencyInfo()
		})
		var failure jutil.Throwable
		for i := range inputs {
			jutil.Try(func() {
				for _, diagnostic := range paths[i].errors {
					diagnostic.target.Report(diagnostic.level, diagnostic.error)
				}
				if failures[i] != nil {
					panic(failures[i])
				}
			}, func(t jutil.Throwable) {
				if failure == nil {
					failure = t
				}
			})
		}
		if failure != nil {
			panic(jutil.NewRuntimeExceptionWithCause(jutil.NewExecutionException(failure)))
		}
		return true
	}

	type outcome struct {
		info         DependencyInfo
		errors       prebuildDependencyErrors
		loader       *ModuleLoader
		readFailure  bool
		moduleErrors prebuildDependencyModuleErrors
	}
	results := make([]outcome, len(inputs))
	for i := range inputs {
		results[i].errors.ErrorManager = compiler.GetErrorManager()
		results[i].moduleErrors = prebuildDependencyModuleErrors{
			target: compiler.GetModuleLoader().GetErrorHandler(), events: &results[i].errors.events,
		}
		results[i].loader = prebuildCloneModuleLoader(compiler.GetModuleLoader(), &results[i].moduleErrors)
		if results[i].loader == nil {
			return false
		}
	}
	failures := parallelPrebuild(len(inputs), p.numParallelThreads, func(i int) {
		input, result := inputs[i], &results[i]
		if input.dependencyInfo != nil {
			return
		}
		jutil.Try(func() {
			result.info = NewJsFileRegexParser(&result.errors).
				SetModuleLoader(result.loader).SetIncludeGoogBase(true).
				ParseFile(input.GetName(), input.GetName(), input.GetCode())
		}, func(e jutil.IIOException) {
			result.readFailure = true
			result.errors.Report(CheckLevel_ERROR.Ptr(),
				JSError_Make(AbstractCompiler_READ_ERROR, input.GetName(), e.GetMessage()))
			result.info = SimpleDependencyInfo_EMPTY
		})
	})
	var failure jutil.Throwable
	for i, input := range inputs {
		result := &results[i]
		jutil.Try(func() {
			for _, diagnostic := range result.errors.events {
				diagnostic.target.Report(diagnostic.level, diagnostic.error)
			}
			if failures[i] != nil {
				panic(failures[i])
			}
			if result.info != nil {
				// IO failures return EMPTY directly, matching generateDependencyInfo.
				if result.readFailure {
					input.dependencyInfo = result.info
				} else {
					input.dependencyInfo = NewLazyParsedDependencyInfo(result.info, input, compiler)
				}
			}
			input.GetDependencyInfo() // merge extra requires/provides in encounter order
		}, func(t jutil.Throwable) {
			if failure == nil {
				failure = t
			}
		})
	}
	if failure != nil {
		panic(jutil.NewRuntimeExceptionWithCause(jutil.NewExecutionException(failure)))
	}
	return true
}
