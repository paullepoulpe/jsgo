// Port of com/google/javascript/jscomp/Compiler.java @ closure-compiler v20261021

/*
 * Copyright 2004 The Closure Compiler Authors.
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
	"os"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CompileOfSourceFiles compiles a single source file and a single externs
// file.
//
// Java: Compiler.compile(SourceFile, SourceFile, CompilerOptions)
func (c *Compiler) CompileOfSourceFiles(extern *SourceFile, input *SourceFile, options *CompilerOptions) *Result {
	return c.Compile([]*SourceFile{extern}, []*SourceFile{input}, options)
}

// Compile compiles a list of inputs.
//
// This is a convenience method to wrap up all the work of compilation,
// including generating the error and warning report.
//
// NOTE: All methods called here must be public, because client code must
// be able to replicate and customize this.
//
// Java: Compiler.compile(List<SourceFile>, List<SourceFile>, CompilerOptions)
func (c *Compiler) Compile(externs []*SourceFile, inputs []*SourceFile, options *CompilerOptions) *Result {
	// The compile method should only be called once.
	jutil.CheckState(c.jsRoot.IsNil())

	func() {
		defer c.GenerateReport()
		c.Init(externs, inputs, options)
		if !c.HasErrors() {
			c.ParseForCompilation()
		}
		if !c.HasErrors() {
			if options.GetInstrumentForCoverageOnly() {
				// TODO(bradfordcsmith): The option to instrument for coverage only should belong to the
				//     runner, not the compiler.
				c.InstrumentForCoverage()
			} else {
				c.Stage1Passes()
				if !c.HasErrors() {
					c.Stage2Passes(CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS)
					if !c.HasErrors() {
						c.Stage3Passes()
					}
				}
			}
			c.PerformPostCompilationTasks()
		}
	}()
	return c.GetResult()
}

// GenerateReport generates a report of all warnings and errors found during
// compilation to stderr.
//
// Client code must call this method explicitly if it doesn't use one of
// the convenience methods that do so automatically.
//
// Always call this method, even if the compiler throws an exception. The
// report will include information about the exception.
//
// Java: Compiler.generateReport()
func (c *Compiler) GenerateReport() {
	t := c.NewTracer("generateReport")
	c.errorManager.GenerateReport()
	c.clearSourcesCachedForReport()
	c.StopTracer(t, "generateReport")
}

// clearSourcesCachedForReport drops the code that the report read into the
// SourceFiles to quote the lines of its warnings, so that it is not kept
// through the output. A later reader loads it again.
//
// Go-only: performance — Java keeps it cached until the Compiler is gone.
func (c *Compiler) clearSourcesCachedForReport() {
	for _, input := range c.inputsById { // order-independent: ClearCachedSource only drops each SourceFile's independent in-memory code cache
		input.GetSourceFile().ClearCachedSource()
	}
	for f := range c.sourceMapOriginalSources.Values() {
		f.ClearCachedSource()
	}
}

// CompileChunks compiles a list of chunks.
//
// This is a convenience method to wrap up all the work of compilation,
// including generating the error and warning report.
//
// NOTE: All methods called here must be public, because client code must
// be able to replicate and customize this.
//
// Java: Compiler.compileChunks(List<SourceFile>, List<JSChunk>, CompilerOptions)
func (c *Compiler) CompileChunks(externs []*SourceFile, chunks []*JSChunk, options *CompilerOptions) *Result {
	// The compile method should only be called once.
	jutil.CheckState(c.jsRoot.IsNil())

	func() {
		defer c.GenerateReport()
		c.InitChunks(externs, chunks, options)
		if !c.HasErrors() {
			c.ParseForCompilation()
		}
		if !c.HasErrors() {
			// TODO(bradfordcsmith): The option to instrument for coverage only should belong to the
			//     runner, not the compiler.
			if options.GetInstrumentForCoverageOnly() {
				c.InstrumentForCoverage()
			} else {
				c.Stage1Passes()
				if !c.HasErrors() {
					c.Stage2Passes(CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS)
					if !c.HasErrors() {
						c.Stage3Passes()
					}
				}
			}
			c.PerformPostCompilationTasks()
		}
	}()
	return c.GetResult()
}

// Stage1Passes performs compiler passes for stage 1 of compilation.
//
// Stage 1 consists primarily of error and type checking passes.
//
// parseForCompilation() must be called before this method is called.
//
// The caller is responsible for also calling generateReport() to generate
// a report of warnings and errors to stderr. See the invocation in compile
// for a good example.
//
// Java: Compiler.stage1Passes()
func (c *Compiler) Stage1Passes() {
	jutil.CheckStateMsg(c.chunkGraph != nil, "No inputs. Did you call init() or initChunks()?")
	jutil.CheckState(!c.HasErrors())
	jutil.CheckState(!c.options.GetInstrumentForCoverageOnly())
	c.runInCompilerThread(c.performChecks)
}

// Stage2Passes performs compiler passes for stage 2 of compilation.
//
// Stage 2 consists primarily of transpilation and optimization passes.
//
// stage1Passes() must be called before this method is called.
//
// The caller is responsible for also calling generateReport() to generate
// a report of warnings and errors to stderr. See the invocation in compile
// for a good example.
//
// Java: Compiler.stage2Passes(SegmentOfCompilationToRun)
func (c *Compiler) Stage2Passes(segmentOfCompilationToRun CompilerOptions_SegmentOfCompilationToRun) {
	jutil.CheckStateMsg1(
		segmentOfCompilationToRun == CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS ||
			segmentOfCompilationToRun == CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS_FIRST_HALF ||
			segmentOfCompilationToRun == CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS_SECOND_HALF,
		"Unsupported segment of compilation to run in stage 2: %s",
		segmentOfCompilationToRun)
	jutil.CheckStateMsg(c.chunkGraph != nil, "No inputs. Did you call init() or initChunks()?")
	jutil.CheckState(!c.HasErrors())
	jutil.CheckState(!c.options.GetInstrumentForCoverageOnly())
	weakChunk := c.chunkGraph.GetChunkByName(JSChunk_WEAK_CHUNK_NAME)
	if weakChunk != nil {
		for _, i := range c.chunkGraph.GetAllInputs() {
			if i.GetSourceFile().IsWeak() {
				jutil.CheckStateMsg(i.GetChunk() == weakChunk, "Expected all weak files to be in the weak chunk.")
			}
		}
	}
	c.runInCompilerThread(func() {
		if c.options.ShouldOptimize() {
			c.PerformTranspilationAndOptimizations(segmentOfCompilationToRun)
		}
	})
}

// Stage3Passes performs compiler passes for stage 3 of compilation.
//
// Stage 3 consists primarily of localization passes.
//
// stage2Passes() must be called before this method is called.
//
// The caller is responsible for also calling generateReport() to generate
// a report of warnings and errors to stderr. See the invocation in compile
// for a good example.
//
// Java: Compiler.stage3Passes()
func (c *Compiler) Stage3Passes() {
	jutil.CheckStateMsg(c.chunkGraph != nil, "No inputs. Did you call init() or initChunks()?")
	jutil.CheckStateMsg1(!c.HasErrors(), "Cannot run stage3Passes with errors present: %s", c.GetErrors())
	jutil.CheckState(!c.options.GetInstrumentForCoverageOnly())
	c.runInCompilerThread(func() {
		if c.options.ShouldOptimize() {
			c.PerformFinalizations()
		}
	})
}

// DisableThreads disables threads. This is for clients that run on
// AppEngine and don't have threads.
//
// Java: Compiler.disableThreads()
func (c *Compiler) DisableThreads() {
	c.compilerExecutor.DisableThreads()
}

// SetTimeout sets the timeout when Compiler is run in a thread
//
// timeout: seconds to wait before timeout
//
// Java: Compiler.setTimeout(int)
func (c *Compiler) SetTimeout(timeout int) {
	c.compilerExecutor.SetTimeout(timeout)
}

// Compiler_RunInCompilerThread: the primary purpose of this method is to
// run the provided code with a larger than standard stack.
//
// Go-only: language-adapter — Java's generic instance method
// <T> T runInCompilerThread(Callable<T>); Go methods cannot have type
// parameters.
//
// Java: Compiler.runInCompilerThread(Callable<T>)
func Compiler_RunInCompilerThread[T any](c *Compiler, callable func() T) T {
	return CompilerExecutor_RunInCompilerThread(
		c.compilerExecutor, callable, c.options != nil && c.options.GetTracerMode().IsOn())
}

func (c *Compiler) runInCompilerThread(fn func()) {
	Compiler_RunInCompilerThread(c, func() struct{} {
		fn()
		return struct{}{}
	})
}

// Java: Compiler.performChecks()
func (c *Compiler) performChecks() {
	if c.options.GetSkipNonTranspilationPasses() {
		// i.e. whitespace-only mode, which will not work with goog.module without:
		c.WhitespaceOnlyPasses()
		if c.options.NeedsTranspilationFrom(c.options.GetLanguageIn().ToFeatureSet()) {
			c.TranspileAndDontCheck()
		}
	} else {
		c.Check()
	}
}

// PerformPostCompilationTasks performs all the bookkeeping required at the
// end of a compilation.
//
// This method must be called if the compilation makes it as far as doing
// checks.
//
// DON'T call it if the compiler threw an exception.
//
// DO call it even when hasErrors() returns true.
//
// Java: Compiler.performPostCompilationTasks()
func (c *Compiler) PerformPostCompilationTasks() {
	c.runInCompilerThread(c.performPostCompilationTasksInternal)
}

// performPostCompilationTasksInternal performs all the bookkeeping
// required at the end of a compilation.
//
// Java: Compiler.performPostCompilationTasksInternal()
func (c *Compiler) performPostCompilationTasksInternal() {
	if c.options.GetDevMode() == CompilerOptions_DevMode_START_AND_END {
		c.runValidityCheck()
	}

	// TODO: b/501060517 - consider deleting this. The code printer will already print out
	// @closureUnaware contents as if they were unwrapped, but the explicit wrapping might be
	// more maintainable.
	ManageClosureUnawareCode_Unwrap(c).Process(c.externsRoot, c.jsRoot)

	if c.tracker != nil {
		if c.options.GetTracerOutput() == nil {
			c.tracker.OutputTracerReport(c.outStream)
		} else {
			jutil.Try(func() {
				stream, err := os.OpenFile(c.options.GetTracerOutput().String(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
				if err != nil {
					panic(jutil.NewIOException(err.Error()))
				}
				out := jutil.NewPrintStream(stream)
				defer func() {
					// PrintStream.close: flush, then close the stream.
					out.Flush()
					_ = stream.Close()
				}()
				c.tracker.OutputTracerReport(out)
			}, func(e jutil.Throwable) {
				// Java: catch (Exception e), which does not see java.lang.Errors.
				if _, isError := e.(jutil.IError); isError {
					panic(e)
				}
				panic(jutil.NewRuntimeExceptionWithCause(e))
			})
		}
	}
}

// InstrumentForCoverage instruments code for coverage.
//
// parseForCompilation() must be called before this method is called.
//
// The caller is responsible for also calling generateReport() to generate
// a report of warnings and errors to stderr. See the invocation in compile
// for a good example.
//
// This method is mutually exclusive with stage1Passes() and
// stage2Passes(). Either call those two methods or this one, but not both.
//
// Java: Compiler.instrumentForCoverage()
func (c *Compiler) InstrumentForCoverage() {
	jutil.CheckStateMsg(c.chunkGraph != nil, "No inputs. Did you call init() or initChunks()?")
	jutil.CheckState(!c.HasErrors())
	c.runInCompilerThread(func() {
		jutil.CheckState(c.options.GetInstrumentForCoverageOnly())
		jutil.CheckState(!c.HasErrors())
		if c.options.GetInstrumentForCoverageOption() != CompilerOptions_InstrumentOption_NONE {
			c.instrumentForCoverageInternal(c.options.GetInstrumentForCoverageOption())
		}
	})
}

// Java: Compiler.instrumentForCoverageInternal(InstrumentOption)
func (c *Compiler) instrumentForCoverageInternal(instrumentForCoverageOption CompilerOptions_InstrumentOption) {
	tracer := c.NewTracer("instrumentationPass")
	c.Process(
		NewCoverageInstrumentationPassWithInstrumentOption(
			c,
			CoverageInstrumentationPass_CoverageReach_ALL,
			instrumentForCoverageOption,
			c.options.GetProductionInstrumentationArrayName()))
	c.StopTracer(tracer, "instrumentationPass")
}

// ParseForCompilation parses input files in preparation for compilation.
//
// Either init() or initChunks() must be called first to set up the input
// files to be read.
//
// TODO(bradfordcsmith): Rename this to parse()
//
// Java: Compiler.parseForCompilation()
func (c *Compiler) ParseForCompilation() {
	jutil.CheckStateMsg(
		c.typedAstFilesystem == nil, "Unnecessary if initWithTypedAstFilesystem was called")
	c.runInCompilerThread(c.parseForCompilationInternal)
}

// parseForCompilationInternal parses input files in preparation for
// compilation.
//
// Either init() or initChunks() must be called first to set up the input
// files to be read.
//
// TODO(bradfordcsmith): Rename this to parse()
//
// Java: Compiler.parseForCompilationInternal()
func (c *Compiler) parseForCompilationInternal() {
	CompilerOptionsPreprocessor_Preprocess(c.options)
	c.MaybeSetTracker()
	c.ParseInputs()
}

// Parse parses input files without doing progress tracking that is part of
// a full compile.
//
// Either init() or initChunks() must be called first to set up the input
// files to be read.
//
// TODO(bradfordcsmith): Rename this to parseIndependentOfCompilation() or
// similar.
//
// Java: Compiler.parse()
func (c *Compiler) Parse() {
	c.ParseInputs()
}

// Java: Compiler.parse(SourceFile)
func (c *Compiler) ParseOfSourceFile(file *SourceFile) rhino.Node {
	c.InitCompilerOptionsIfTesting()
	Compiler_logger.Finest("Parsing: " + file.GetName())
	return NewCompilerInput(file).GetAstRoot(c)
}

// Java: Compiler.getPassConfig()
func (c *Compiler) GetPassConfig() PassConfig {
	if c.passes == nil {
		c.passes = c.CreatePassConfigInternal()
	}
	return c.passes
}

// CreatePassConfigInternal creates the passes object. Clients should use
// setPassConfig instead of overriding this.
//
// Java: Compiler.createPassConfigInternal()
func (c *Compiler) CreatePassConfigInternal() PassConfig {
	return NewDefaultPassConfig(c.options)
}

// SetPassConfig sets the PassConfig to use with this Compiler.
//
// Panics with NullPointerException if passes is null, and with
// IllegalStateException if this.passes has already been assigned.
//
// Java: Compiler.setPassConfig(PassConfig)
func (c *Compiler) SetPassConfig(passes PassConfig) {
	// Important to check for null because if setPassConfig(null) is
	// called before this.passes is set, getPassConfig() will create a
	// new PassConfig object and use that, which is probably not what
	// the client wanted since they probably meant to use their
	// own PassConfig object.
	jutil.CheckNotNull(passes)
	jutil.CheckStateMsg(c.passes == nil, "setPassConfig was already called")
	c.passes = passes
}

// Java: Compiler.whitespaceOnlyPasses()
func (c *Compiler) WhitespaceOnlyPasses() {
	c.runCustomPasses(CustomPassExecutionTime_BEFORE_CHECKS)

	t := c.NewTracer("runWhitespaceOnlyPasses")
	defer c.StopTracer(t, "runWhitespaceOnlyPasses")
	for _, pf := range c.GetPassConfig().GetWhitespaceOnlyPasses().Build() {
		pf.Create(c).Process(c.externsRoot, c.jsRoot)
	}
}

// Java: Compiler.markTranspiledFiles()
func (c *Compiler) markTranspiledFiles() {
	// note whether any script will get transpiled
	for child := c.jsRoot.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		jutil.CheckState(child.IsScript())
		featureSet := NodeUtil_GetFeatureSetOfScript(child)
		if !c.options.GetOutputFeatureSet().Contains(featureSet) {
			c.SetTranspiledFiles(true)
			return // marked atleast 1 script that needs tranpile
		}
	}
}

// Java: Compiler.transpileAndDontCheck()
func (c *Compiler) TranspileAndDontCheck() {
	t := c.NewTracer("runTranspileOnlyPasses")
	c.markTranspiledFiles()
	defer c.StopTracer(t, "runTranspileOnlyPasses")
	for _, pf := range c.GetPassConfig().GetTranspileOnlyPasses().Build() {
		if c.HasErrors() {
			return
		}
		c.BeforePass(pf.GetName())
		pf.Create(c).Process(c.externsRoot, c.jsRoot)
		c.AfterPass(pf.GetName())
	}
}

// Java: Compiler.createPhaseOptimizer()
func (c *Compiler) createPhaseOptimizer() *PhaseOptimizer {
	phaseOptimizer := NewPhaseOptimizer(c, c.tracker)
	if c.options.GetDevMode() == CompilerOptions_DevMode_EVERY_PASS {
		phaseOptimizer.SetValidityCheck(c.validityCheck)
	}
	if c.options.GetCheckDeterminism() {
		phaseOptimizer.SetPrintAstHashcodes(true)
	}
	return phaseOptimizer
}

// Java: Compiler.check()
func (c *Compiler) Check() {
	c.runCustomPasses(CustomPassExecutionTime_BEFORE_CHECKS)

	c.phaseOptimizer = c.createPhaseOptimizer()
	c.phaseOptimizer.Consume(c.GetPassConfig().GetChecks().Build())
	c.phaseOptimizer.Process(c.externsRoot, c.jsRoot)
	if c.HasErrors() {
		return
	}

	c.runCustomPasses(CustomPassExecutionTime_BEFORE_OPTIMIZATIONS)
	c.phaseOptimizer = nil
}

// Java: Compiler.setExternExports(String)
func (c *Compiler) SetExternExports(externExports string) {
	c.externExports = &externExports
}

// Java: Compiler.process(CompilerPass)
func (c *Compiler) Process(p CompilerPass) {
	p.Process(c.externsRoot, c.jsRoot)
}

// Java: Compiler.runValidityCheck()
func (c *Compiler) runValidityCheck() {
	c.validityCheck.Create(c).Process(c.externsRoot, c.jsRoot)
}

// runCustomPasses runs custom passes that are designated to run at a
// particular time.
//
// Java: Compiler.runCustomPasses(CustomPassExecutionTime)
func (c *Compiler) runCustomPasses(executionTime CustomPassExecutionTime) {
	customPasses := c.options.GetCustomPassesAt(executionTime)
	t := c.NewTracer("runCustomPasses")
	defer c.StopTracer(t, "runCustomPasses")
	for _, p := range customPasses {
		c.Process(p)
	}
}

// Java: Compiler.beforePass(String)
func (c *Compiler) BeforePass(passName string) {
	c.currentPassIndex++
}

// Java: Compiler.afterPass(String)
func (c *Compiler) AfterPass(passName string) {
	c.maybePrintSourceAfterEachPass(passName)
	if DebugAfterPassHook != nil {
		DebugAfterPassHook(c, passName)
	}
}

// DebugAfterPassHook, when set, is called after every pass with the pass
// name. Go-only: lets cmd/jsgo write a heap profile after a named pass
// (JSGO_HEAP_PROFILE_AFTER_PASS) for memory investigations.
var DebugAfterPassHook func(c *Compiler, passName string)

// GetLifeCycleStage returns the current life-cycle stage of the AST we're
// working on.
//
// Java: Compiler.getLifeCycleStage()
func (c *Compiler) GetLifeCycleStage() AbstractCompiler_LifeCycleStage {
	return c.stage
}

// SetLifeCycleStage sets the current life-cycle state.
//
// Java: Compiler.setLifeCycleStage(LifeCycleStage)
func (c *Compiler) SetLifeCycleStage(stage AbstractCompiler_LifeCycleStage) {
	c.stage = stage
}

// Java: Compiler.performFinalizations()
func (c *Compiler) PerformFinalizations() {
	finalizations := c.GetPassConfig().GetFinalizations().Build()
	if len(finalizations) == 0 {
		return
	}

	c.phaseOptimizer = c.createPhaseOptimizer()
	c.phaseOptimizer.Consume(finalizations)
	c.phaseOptimizer.Process(c.externsRoot, c.jsRoot)
	c.phaseOptimizer = nil
}

// getOptimizationPassesToRunInCurrentSegment takes the entire list of
// optimization passes to run in the current build, and returns the subset
// of passes that should be run in this current segment of optimizations.
//
// The 3 cases are:
//
//  1. If we are running `SegmentOfCompilationToRun.OPTIMIZATIONS`, this
//     will return the entire list of optimization passes.
//  2. If we are running `SegmentOfCompilationToRun.OPTIMIZATIONS_FIRST_HALF`,
//     this will return all passes up to the half-way point marker pass
//     `OPTIMIZATIONS_HALFWAY_POINT`.
//  3. If we are running `SegmentOfCompilationToRun.OPTIMIZATIONS_SECOND_HALF`,
//     this will return all passes after the half-way point marker pass
//     `OPTIMIZATIONS_HALFWAY_POINT`.
//
// Java: Compiler.getOptimizationPassesToRunInCurrentSegment(SegmentOfCompilationToRun, List<PassFactory>)
func (c *Compiler) getOptimizationPassesToRunInCurrentSegment(
	segmentOfCompilationToRun CompilerOptions_SegmentOfCompilationToRun,
	allOptimizationPassesToRun []*PassFactory) []*PassFactory {

	// Create a list of passes based on which segment of optimizations we are running.
	switch segmentOfCompilationToRun {
	case CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS:
		return allOptimizationPassesToRun
	case CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS_FIRST_HALF:
		return c.passListUntil(allOptimizationPassesToRun, PassNames_OPTIMIZATIONS_HALFWAY_POINT)
	case CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS_SECOND_HALF:
		return c.passListAfter(allOptimizationPassesToRun, PassNames_OPTIMIZATIONS_HALFWAY_POINT)
	default:
		panic(jutil.NewIllegalArgumentExceptionMsg(
			"Unsupported segment of optimizations to run: " + segmentOfCompilationToRun.String()))
	}
}

// passListUntil returns a subset of passes from the given list of passes,
// up to but not including the given pass name.
//
// Java: Compiler.passListUntil(List<PassFactory>, String)
func (c *Compiler) passListUntil(passes []*PassFactory, passName string) []*PassFactory {
	passList := []*PassFactory{}
	for _, pass := range passes {
		if pass.GetName() == passName {
			return passList
		}
		passList = append(passList, pass)
	}
	return passList
}

// passListAfter returns a subset of passes from the given list of passes,
// starting right after the given pass name.
//
// Java: Compiler.passListAfter(List<PassFactory>, String)
func (c *Compiler) passListAfter(passes []*PassFactory, passName string) []*PassFactory {
	passList := []*PassFactory{}
	for _, pass := range passes {
		if pass.GetName() == passName {
			passList = passList[:0]
			continue
		}
		passList = append(passList, pass)
	}
	return passList
}
