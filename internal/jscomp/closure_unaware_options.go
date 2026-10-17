// Port of com/google/javascript/jscomp/ClosureUnawareOptions.java @ closure-compiler v20260415

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

import "github.com/paullepoulpe/jsgo/internal/jutil"

// ClosureUnawareOptions creates a partial copy of a CompilerOptions object
// for use in TranspileAndOptimizeClosureUnawareCode.
//
// During a regular compilation (i.e. "main compilation"), we have a special
// pass that will instantiate a new Compiler object with a new
// CompilerOptions, created by this class. We then run
// optimizations/finalizations in that pass. This is called a "nested
// compilation" or "sub-compilation".
//
// As a general principle:
//
//   - always copy options that are related to code printing or debugging.
//   - sometimes copy options related to transpilation: the nested
//     compilation should try to match the main compilation's target language
//     out, but we don't actually support automatic polyfill injection and we
//     have some unique runtime library handling.
//   - never copy options that are related to optimization. The purpose of
//     @closureUnaware code blocks is to protect arbitrary 3P code from
//     breaking under advanced JSCompiler optimizations. Instead, we
//     explicitly set the relevant, safe optimization options a preset list
//     of options, ignoring the main compilation's settings.
//   - never copy options related to check passes or parsing: the nested
//     compilation doesn't actually run the parser or any checks (besides
//     VarCheck), so there's no need to explicitly enumerate and copy these.
type ClosureUnawareOptions struct {
	shadowOptions *CompilerOptions
	original      *CompilerOptions
}

// Java: ClosureUnawareOptions.convert(CompilerOptions)
func ClosureUnawareOptions_Convert(original *CompilerOptions) *CompilerOptions {
	options := NewClosureUnawareOptions(original)
	return options.ToCompilerOptions()
}

// Java: ClosureUnawareOptions(CompilerOptions)
func NewClosureUnawareOptions(original *CompilerOptions) *ClosureUnawareOptions {
	c := &ClosureUnawareOptions{}
	c.shadowOptions = NewCompilerOptions()
	c.original = original
	return c
}

// Java: ClosureUnawareOptions.toCompilerOptions()
func (c *ClosureUnawareOptions) ToCompilerOptions() *CompilerOptions {
	c.SetTranspilationOptions()
	c.SetSafeOptimizationAssumptions()
	c.CopyOutputOptions()
	c.CopyDebugOptions()
	return c.shadowOptions
}

// Java: ClosureUnawareOptions.setTranspilationOptions()
func (c *ClosureUnawareOptions) SetTranspilationOptions() {
	c.shadowOptions.SetRewritePolyfills(false) // for now, we ignore polyfills that may be needed.
	switch c.original.GetRuntimeLibraryMode() {
	case RuntimeJsLibManager_RuntimeLibraryMode_NO_OP:
		c.shadowOptions.SetRuntimeLibraryMode(RuntimeJsLibManager_RuntimeLibraryMode_NO_OP)
	case RuntimeJsLibManager_RuntimeLibraryMode_RECORD_ONLY, RuntimeJsLibManager_RuntimeLibraryMode_RECORD_AND_VALIDATE_FIELDS,
		RuntimeJsLibManager_RuntimeLibraryMode_INJECT:
		c.shadowOptions.SetRuntimeLibraryMode(
			RuntimeJsLibManager_RuntimeLibraryMode_EXTERN_FIELD_NAMES)
	case RuntimeJsLibManager_RuntimeLibraryMode_EXTERN_FIELD_NAMES:
		panic(jutil.NewAssertionErrorNoMessage())
	}
	c.shadowOptions.SetSkipNonTranspilationPasses(c.original.GetSkipNonTranspilationPasses())
	c.shadowOptions.SetEs6SubclassTranspilation(
		CompilerOptions_Es6SubclassTranspilation_SAFE_REFLECT_CONSTRUCT)
	c.shadowOptions.SetOutputFeatureSet(c.original.GetOutputFeatureSet())
	c.shadowOptions.SetInstrumentAsyncContext(c.original.GetInstrumentAsyncContext())

	// To aid rollout of @closureUnaware transpilation, ignore warnings about untranspilable
	// features for now. (Any project using @closureUnaware today (i.e. early 2026) gets zero
	// transpilation, so rolling out with these suppressions is an incremental improvement.)

	// DiagnosticGroups.UNTRANSPILABLE_FEATURES: features JSCompiler will never transpile but
	// can pass through unchanged to the output, such as newer regex syntax.
	c.shadowOptions.SetWarningLevel(DiagnosticGroups_UNTRANSPILABLE_FEATURES, CheckLevel_OFF)
}

// Java: ClosureUnawareOptions.setSafeOptimizationAssumptions()
func (c *ClosureUnawareOptions) SetSafeOptimizationAssumptions() {
	CompilationLevel_SIMPLE_OPTIMIZATIONS.SetOptionsForCompilationLevel(c.shadowOptions)
	c.shadowOptions.SetRemoveUnusedVariables(CompilerOptions_Reach_ALL)

	c.shadowOptions.SetComputeFunctionSideEffects(false) // unsafe for properties.
	c.shadowOptions.SetCodingConvention(CodingConventions_GetDefault())
	c.shadowOptions.SetRemoveClosureAsserts(false)

	c.shadowOptions.SetAssumeClosuresOnlyCaptureReferences(false)
	c.shadowOptions.SetAssumeGettersArePure(false)
	c.shadowOptions.SetAssumePropertiesAreStaticallyAnalyzable(false)
	c.shadowOptions.SetAssumeStaticInheritanceIsNotUsed(false)
	c.shadowOptions.SetAssumeStrictThis(false)

	if c.original.GetMaxFunctionSizeAfterInlining() !=
		CompilerOptions_UNLIMITED_FUN_SIZE_AFTER_INLINING {
		c.shadowOptions.SetMaxFunctionSizeAfterInlining(c.original.GetMaxFunctionSizeAfterInlining())
	}
	c.shadowOptions.SetUseSizeHeuristicToStopOptimizationLoop(
		c.original.ShouldUseSizeHeuristicToStopOptimizationLoop())
	c.shadowOptions.SetMaxOptimizationLoopIterations(c.original.GetMaxOptimizationLoopIterations())

	c.shadowOptions.SetPreferStableNames(c.original.ShouldPreferStableNames())
	c.shadowOptions.SetInputPropertyMap(c.original.GetInputPropertyMap())
	c.shadowOptions.SetInputVariableMap(c.original.GetInputVariableMap())
}

// Java: ClosureUnawareOptions.copyOutputOptions()
func (c *ClosureUnawareOptions) CopyOutputOptions() {
	c.shadowOptions.SetEmitUseStrict(c.original.ShouldEmitUseStrict())
	c.shadowOptions.SetOutputCharset(c.original.GetOutputCharset())

	c.shadowOptions.SetGeneratePseudoNames(c.original.ShouldGeneratePseudoNames())
	c.shadowOptions.SetLineBreak(c.original.ShouldAddLineBreak())
	c.shadowOptions.SetLineLengthThreshold(c.original.GetLineLengthThreshold())
	c.shadowOptions.SetPreferSingleQuotes(c.original.ShouldPreferSingleQuotes())
	c.shadowOptions.SetPrettyPrint(c.original.IsPrettyPrint())

	c.shadowOptions.SetErrorFormat(c.original.GetErrorFormat())
	c.shadowOptions.SetErrorHandler(c.original.GetErrorHandler())
	c.shadowOptions.SetColorizeErrorOutput(c.original.ShouldColorizeErrorOutput())
}

// Java: ClosureUnawareOptions.copyDebugOptions()
func (c *ClosureUnawareOptions) CopyDebugOptions() {
	c.shadowOptions.SetTracerMode(c.original.GetTracerMode())
	c.shadowOptions.SetTracerOutput(c.original.GetTracerOutput())
	c.shadowOptions.SetDevMode(c.original.GetDevMode())

	c.shadowOptions.SetPrintSourceAfterEachPass(c.original.ShouldPrintSourceAfterEachPass())
	c.shadowOptions.SetFilesToPrintAfterEachPassRegexList(
		c.original.GetFilesToPrintAfterEachPassRegexList())
	c.shadowOptions.SetChunksToPrintAfterEachPassRegexList(
		c.original.GetChunksToPrintAfterEachPassRegexList())
	c.shadowOptions.SetPrintInputDelimiter(c.original.ShouldPrintInputDelimiter())
	c.shadowOptions.SetQnameUsesToPrintAfterEachPassList(
		c.original.GetQnameUsesToPrintAfterEachPassList())

	c.shadowOptions.SetInputDelimiter(c.original.GetInputDelimiter())
	c.shadowOptions.SetUseOriginalNamesInOutput(c.original.GetUseOriginalNamesInOutput())

	c.shadowOptions.SetDebugLogFilter(c.original.GetDebugLogFilter())
	debugLogDirectory := c.original.GetDebugLogDirectory()
	if debugLogDirectory != nil {
		c.shadowOptions.SetDebugLogDirectory(
			debugLogDirectory.ResolveString("./TranspileAndOptimizeClosureUnaware"))
	}
}
