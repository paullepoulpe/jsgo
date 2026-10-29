// Port of com/google/javascript/jscomp/PrebuildAst.java @ closure-compiler v20261021

/*
 * Copyright 2017 The Closure Compiler Authors.
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

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// PrebuildAst is a helper class to prebuild ASTs from a list of
// CompilerInput. Inputs are parsed into ASTs the first time we try to get
// them. Get them all now using multiple threads, so they'll be parsed in
// parallel and already available during the rest of the compilation.
type PrebuildAst struct {
	compiler           AbstractCompiler
	numParallelThreads int
}

// Java: PrebuildAst(AbstractCompiler, int)
func NewPrebuildAst(compiler AbstractCompiler, numParalleThreads int) *PrebuildAst {
	return &PrebuildAst{compiler: compiler, numParallelThreads: numParalleThreads}
}

// Prebuild gets the AST root of every input.
//
// Go-only: language-adapter — independent files and TypedAST scripts are
// materialized by a bounded worker pool. Compiler state and diagnostics are
// published in input order after the workers finish. Custom compilers and
// shared sources use the sequential compatibility path. Every task completes
// even if another fails; the first failure in input order has Java's
// ExecutionException wrapping.
//
// Java: PrebuildAst.prebuild(Iterable<CompilerInput>)
func (p *PrebuildAst) Prebuild(allInputs iter.Seq[*CompilerInput]) {
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
		if t := p.getAstRootCatching(input); t != nil && failure == nil {
			failure = t
		}
	}

	if failure != nil {
		panic(jutil.NewRuntimeExceptionWithCause(jutil.NewExecutionException(failure)))
	}
}

func (p *PrebuildAst) getAstRootCatching(input *CompilerInput) (caught jutil.Throwable) {
	defer func() {
		if r := recover(); r != nil {
			if t, ok := r.(jutil.Throwable); ok {
				caught = t
				return
			}
			caught = jutil.AsThrowable(r)
		}
	}()
	input.GetAstRoot(p.compiler)
	return nil
}

type prebuildOutcome struct {
	scriptDeserializer *ScriptNodeDeserializer
	deserializer       func() rhino.Node
	deserializedRoot   rhino.Node
	parsed             *parsing.ParserRunner_ParseResult
	errors             parseBufferedErrors
	readErr            jutil.IIOException
	failure            jutil.Throwable
	wasDeserialized    bool
}

func takeDeserializerCatching(compiler *Compiler, sf *SourceFile, out *prebuildOutcome) (needParser bool) {
	defer func() {
		if r := recover(); r != nil {
			if t, ok := r.(jutil.Throwable); ok {
				out.failure = t
				return
			}
			out.failure = jutil.AsThrowable(r)
		}
	}()
	sd, fn := compiler.takeTypedAstDeserializer(sf)
	if fn != nil {
		if sd != nil {
			out.scriptDeserializer = sd
		} else {
			out.deserializer = fn
		}
		out.wasDeserialized = true
		return false
	}
	return true
}

func prebuildWorkerStep(
	input *CompilerInput,
	result *prebuildOutcome,
	normal, extern *parsing.Config,
	alloc *rhino.NodeAllocator,
	index int,
	recordFailure func(int, jutil.Throwable),
) {
	defer func() {
		if r := recover(); r != nil {
			result.scriptDeserializer = nil
			result.deserializer = nil
			if t, ok := r.(jutil.Throwable); ok {
				recordFailure(index, t)
				return
			}
			recordFailure(index, jutil.AsThrowable(r))
		}
	}()
	if input.ast.isParsed() || result.failure != nil {
		result.scriptDeserializer = nil
		result.deserializer = nil
		return
	}
	if result.wasDeserialized {
		if sd := result.scriptDeserializer; sd != nil {
			result.scriptDeserializer = nil
			result.deserializedRoot = sd.DeserializeNew()
			return
		}
		fn := result.deserializer
		result.deserializer = nil
		result.deserializedRoot = fn()
		return
	}
	code, readErr := input.sourceFile.getCode()
	if readErr != nil {
		if ioex, ok := readErr.(jutil.IIOException); ok {
			result.readErr = ioex
		} else {
			result.readErr = jutil.NewIOException(readErr.Error())
		}
		return
	}
	config := normal
	if input.sourceFile.IsExtern() {
		config = extern
	}
	result.parsed = parsing.ParserRunner_ParseWithAllocator(
		input.sourceFile, code, config, &result.errors, alloc)
}

func publishPrebuiltInput(
	compiler *Compiler,
	reporter rhino.ErrorReporter,
	input *CompilerInput,
	result *prebuildOutcome,
) (caught jutil.Throwable) {
	defer func() {
		if r := recover(); r != nil {
			if t, ok := r.(jutil.Throwable); ok {
				caught = t
				return
			}
			caught = jutil.AsThrowable(r)
		}
	}()
	if result.wasDeserialized {
		if result.failure != nil {
			return result.failure
		}
		input.ast.installDeserializedRoot(result.deserializedRoot)
	} else {
		result.errors.replay(reporter)
		if result.readErr != nil {
			compiler.Report(JSError_Make(AbstractCompiler_READ_ERROR, input.sourceFile.GetName(), result.readErr.GetMessage()))
		} else if result.failure != nil {
			if ioerr, ok := result.failure.(jutil.IIOException); ok {
				compiler.Report(JSError_Make(AbstractCompiler_READ_ERROR, input.sourceFile.GetName(), ioerr.GetMessage()))
			} else {
				return result.failure
			}
		} else {
			input.ast.installParseResult(compiler, result.parsed)
		}
		input.ast.ensureScriptRoot()
	}
	jutil.CheckState(input.ast.root.GetStaticSourceFile() == rhino.StaticSourceFile(input.sourceFile))
	input.ast.root.SetInputId(input.id)
	input.sourceFile.ClearCachedSource()
	return nil
}

// Go-only: language-adapter — parse files or deserialize TypedAST script roots
// independently, then publish roots, diagnostics, comments and source maps in
// submission order. Compiler maps and error managers are touched only by the
// calling goroutine. Unusual/shared SourceFiles keep the original sequential
// path.
func (p *PrebuildAst) prebuildParallel(inputs []*CompilerInput) bool {
	compiler, ok := p.compiler.(*Compiler)
	if !ok || compiler == nil || len(inputs) < 2 {
		return false
	}
	seen := make(map[*SourceFile]bool, len(inputs))
	for _, input := range inputs {
		if input == nil || input.ast == nil || seen[input.sourceFile] {
			return false
		}
		seen[input.sourceFile] = true
	}
	var normal, extern *parsing.Config
	needParserConfig := compiler.typedAstFilesystem == nil
	outcomes := make([]prebuildOutcome, len(inputs))
	if compiler.typedAstFilesystem != nil {
		for index, input := range inputs {
			if input.ast.isParsed() {
				continue
			}
			if takeDeserializerCatching(compiler, input.sourceFile, &outcomes[index]) {
				needParserConfig = true
			}
		}
	}
	if needParserConfig {
		normal = compiler.GetParserConfig(AbstractCompiler_ConfigContext_DEFAULT)
		extern = compiler.GetParserConfig(AbstractCompiler_ConfigContext_EXTERNS)
	}
	failures := parallelPrebuildWorkers(len(inputs), p.numParallelThreads, func(claim func() (int, bool), recordFailure func(int, jutil.Throwable)) {
		alloc := rhino.NewNodeAllocator() // this worker's node ids
		for {
			index, ok := claim()
			if !ok {
				return
			}
			prebuildWorkerStep(inputs[index], &outcomes[index], normal, extern, alloc, index, recordFailure)
		}
	})
	for i, f := range failures {
		if f != nil && outcomes[i].failure == nil {
			outcomes[i].failure = f
		}
	}
	var failure jutil.Throwable
	reporter := compiler.GetDefaultErrorReporter()
	for index, input := range inputs {
		if input.ast.isParsed() {
			continue
		}
		if t := publishPrebuiltInput(compiler, reporter, input, &outcomes[index]); t != nil && failure == nil {
			failure = t
		}
	}
	if failure != nil {
		panic(jutil.NewRuntimeExceptionWithCause(jutil.NewExecutionException(failure)))
	}
	return true
}

type parseBufferedError struct {
	warning         bool
	message, source string
	line, offset    int
}
type parseBufferedErrors struct{ events []parseBufferedError }

func (r *parseBufferedErrors) Error(message, source string, line, offset int) {
	r.events = append(r.events, parseBufferedError{false, message, source, line, offset})
}
func (r *parseBufferedErrors) Warning(message, source string, line, offset int) {
	r.events = append(r.events, parseBufferedError{true, message, source, line, offset})
}
func (r *parseBufferedErrors) replay(target rhino.ErrorReporter) {
	for _, e := range r.events {
		if e.warning {
			target.Warning(e.message, e.source, e.line, e.offset)
		} else {
			target.Error(e.message, e.source, e.line, e.offset)
		}
	}
}
