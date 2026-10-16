// Port of com/google/javascript/jscomp/NestedCompilerRunner.java @ closure-compiler v20260415

/*
 * Copyright 2025 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// NestedCompilerRunner is a class for taking individual @closureUnaware
// "shadow AST" scripts, which are expected to have been initialized as part
// of a main compilation, and running a nested compilation over just those
// scripts.
//
// This class isn't tightly coupled with the exact @closureUnaware
// implementation in JSCompiler, but is still designed to only work for
// @closureUnaware scripts and will throw an exception if asked to compile
// any other code.
//
// A "nested" compilation is distinguished from a regular compilation mainly
// us knowing we are in the middle of the optimization pass loop of a regular
// compilation; so we want to delegate certain things like error reporting or
// debug logging to the main compiler.
type NestedCompilerRunner struct {
	shadowCompiler *Compiler
	inputs         *jutil.ArrayList[*CompilerInput]
	baseChunk      *JSChunk

	original         AbstractCompiler
	shadowOptions    *CompilerOptions
	mode             NestedCompilerRunner_Mode
	shadowPassConfig PassConfig
}

// Java: NestedCompilerRunner.Mode
type NestedCompilerRunner_Mode int32

const (
	NestedCompilerRunner_Mode_TRANSPILE_AND_OPTIMIZE NestedCompilerRunner_Mode = iota
	NestedCompilerRunner_Mode_TRANSPILE_ONLY
)

// Java: NestedCompilerRunner.Mode.values()
func NestedCompilerRunner_Mode_values() []NestedCompilerRunner_Mode {
	return []NestedCompilerRunner_Mode{
		NestedCompilerRunner_Mode_TRANSPILE_AND_OPTIMIZE,
		NestedCompilerRunner_Mode_TRANSPILE_ONLY,
	}
}

// Ordinal is Enum.ordinal().
func (m NestedCompilerRunner_Mode) Ordinal() int { return int(m) }

// String is Enum.name().
func (m NestedCompilerRunner_Mode) String() string {
	switch m {
	case NestedCompilerRunner_Mode_TRANSPILE_AND_OPTIMIZE:
		return "TRANSPILE_AND_OPTIMIZE"
	case NestedCompilerRunner_Mode_TRANSPILE_ONLY:
		return "TRANSPILE_ONLY"
	}
	panic(jutil.NewIllegalArgumentException())
}

// Java: NestedCompilerRunner(AbstractCompiler, CompilerOptions, Mode, PassConfig)
func newNestedCompilerRunner(
	original AbstractCompiler,
	shadowOptions *CompilerOptions,
	mode NestedCompilerRunner_Mode,
	shadowPassConfig PassConfig) *NestedCompilerRunner {
	r := &NestedCompilerRunner{
		shadowCompiler: NewCompiler(),
		inputs:         jutil.NewArrayList[*CompilerInput](),
		baseChunk:      NewJSChunk("shadow_base"),
	}
	r.original = original
	r.shadowOptions = shadowOptions
	r.mode = mode
	r.shadowPassConfig = shadowPassConfig
	return r
}

// Java: NestedCompilerRunner.create(AbstractCompiler, CompilerOptions, Mode, PassConfig)
func NestedCompilerRunner_Create(
	original AbstractCompiler,
	shadowOptions *CompilerOptions,
	mode NestedCompilerRunner_Mode,
	shadowPassConfig PassConfig) *NestedCompilerRunner {
	return newNestedCompilerRunner(original, shadowOptions, mode, shadowPassConfig)
}

// Java: NestedCompilerRunner.addScript(Node, String)
func (r *NestedCompilerRunner) AddScript(shadowScript *rhino.Node, uniqueShadowId string) *NestedCompilerRunner {
	jutil.CheckArgumentMsg1(
		shadowScript.GetIsInClosureUnawareSubtree(),
		"Expected closureUnaware script, found %s",
		shadowScript)

	shadowFile := SourceFile_FromCodeWithKind(uniqueShadowId, "", rhino.StaticSourceFile_SourceKind_STRONG)

	shadowInput := NewCompilerInputWithInputId(shadowFile, rhino.NewInputId(uniqueShadowId))
	shadowInput.InitShadowAst(shadowScript)

	shadowScript.SetInputId(shadowInput.GetInputId())
	shadowScript.SetStaticSourceFile(shadowFile)
	// TODO: b/421971366 - should we also update the shadow script children with this source file?

	r.inputs.Add(shadowInput)
	return r
}

// Compile initializes a new Compiler with all shadow ASTs and runs
// transpilation/optimizations.
//
// Java: NestedCompilerRunner.compile()
func (r *NestedCompilerRunner) Compile() {
	r.initializeCompiler()
	r.runCompilation()
	r.checkInvariants()
}

// initializeCompiler creates a new Compiler and initializes it with a chunk
// graph holding all shadow ASTs.
//
// Java: NestedCompilerRunner.initializeCompiler()
func (r *NestedCompilerRunner) initializeCompiler() {
	externNames :=
		NodeUtil_CollectExternVariableNames(r.original, r.original.GetRoot().GetFirstChild())
	r.shadowCompiler.AddExportedNames(externNames)
	r.shadowCompiler.SetExternProperties(r.original.GetExternProperties())
	debugMessage := "<shadow AST compilation>"
	r.shadowCompiler.SetDebugMessage(&debugMessage)

	chunks := r.createChunks()
	r.shadowCompiler.InitChunks([]*SourceFile{}, chunks, r.shadowOptions)
	r.shadowCompiler.SetPassConfig(r.shadowPassConfig)
	r.shadowCompiler.ParseForCompilation()
}

// Java: NestedCompilerRunner.createChunks()
func (r *NestedCompilerRunner) createChunks() []*JSChunk {
	r.baseChunk.Add(NewCompilerInput(SourceFile_FromCodeWithKind("synthetic_base", "", rhino.StaticSourceFile_SourceKind_STRONG)))
	chunks := []*JSChunk{r.baseChunk}
	for input := range r.inputs.All() {
		chunk := NewJSChunk("shadow_chunk_" + input.GetSourceFile().GetName())
		chunk.Add(input)
		chunk.AddDependency(r.baseChunk)
		chunks = append(chunks, chunk)
	}
	return chunks
}

// runCompilation runs compiler passes over the shadow ASTs.
//
// Java: NestedCompilerRunner.runCompilation()
func (r *NestedCompilerRunner) runCompilation() {
	switch r.mode {
	case NestedCompilerRunner_Mode_TRANSPILE_AND_OPTIMIZE:
		// Run VarCheck to comply with a compiler invariant: when running optimizations, all NAME
		// nodes in the AST must have a corresponding declaration in source or externs. The
		// optimizer will either crash or silently misoptimize code otherwise.
		// (VarCheck not only reports errors for referencing undefined variables, but also adds
		// "synthetic extern declarations" of undefined variables, in case the errors are suppressed
		// and we want to continue compilation)
		NewVarCheck(r.shadowCompiler).
			Process(r.shadowCompiler.GetExternsRoot(), r.shadowCompiler.GetJsRoot())
		r.delegateDiagnosticsToOriginal()
		if r.shadowCompiler.HasHaltingErrors() {
			return
		}
		r.shadowCompiler.Stage2Passes(CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS)
		r.delegateDiagnosticsToOriginal()
		if r.shadowCompiler.HasHaltingErrors() {
			return
		}
		r.shadowCompiler.Stage3Passes()
		r.delegateDiagnosticsToOriginal()
	case NestedCompilerRunner_Mode_TRANSPILE_ONLY:
		r.shadowCompiler.TranspileAndDontCheck()
		r.delegateDiagnosticsToOriginal()
	}
}

// Java: NestedCompilerRunner.delegateDiagnosticsToOriginal()
func (r *NestedCompilerRunner) delegateDiagnosticsToOriginal() {
	for _, error_ := range r.shadowCompiler.GetErrors() {
		r.original.Report(error_)
	}
	for _, error_ := range r.shadowCompiler.GetWarnings() {
		r.original.Report(error_)
	}
}

// Java: NestedCompilerRunner.checkInvariants()
func (r *NestedCompilerRunner) checkInvariants() {
	jutil.CheckStateMsg1(
		len(r.baseChunk.GetInputs()) == 1,
		"Expected exactly 1 base chunk input, got %s",
		r.baseChunk.GetInputs())

	// Check that compilation did not introduce any new code into the base chunk: we don't
	// have a place to put that code in the main AST, once we move out of this nested compilation.
	// Note: if it turns out to be useful to have some shared common base code for shadow ASTs, we
	// could consider designing a "common" closureUnaware shadow AST that goes in the actual main
	// AST root chunk. But we don't have a use case for that yet.
	baseInputNode := r.baseChunk.GetInputs()[0].GetAstRoot(r.shadowCompiler)
	jutil.CheckStateMsg1(
		!baseInputNode.HasChildren(),
		"Expected synthetic base input chunk to have zero children, but found %s",
		baseInputNode.ToStringTree())
}
