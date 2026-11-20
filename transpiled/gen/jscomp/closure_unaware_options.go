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
// Port of com/google/javascript/jscomp/ClosureUnawareOptions.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.ClosureUnawareOptions
type ClosureUnawareOptions struct {
	shadowOptions *CompilerOptions `gson:"shadowOptions"`
	original *CompilerOptions `gson:"original"`
	idHash_ atomic.Int32
}

func (this_ *ClosureUnawareOptions) BindSelf_(s any) {
	_ = s
}
func (this_ *ClosureUnawareOptions) IsJavaNil_() bool { return this_ == nil }
func (*ClosureUnawareOptions) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.ClosureUnawareOptions", "ClosureUnawareOptions" }
func (this_ *ClosureUnawareOptions) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ClosureUnawareOptions) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ClosureUnawareOptions", "com.google.javascript.jscomp.ClosureUnawareOptions", func(v_ any) bool {
		return j2gort.InstanceOf[*ClosureUnawareOptions](v_)
	})
}

func (this_ *ClosureUnawareOptions) initFields_ClosureUnawareOptions() {
	_ = this_
	this_.shadowOptions = NewCompilerOptions()
}

// Java constructor: ClosureUnawareOptions(com.google.javascript.jscomp.CompilerOptions original)
func NewClosureUnawareOptions(original *CompilerOptions) *ClosureUnawareOptions {
	this_ := &ClosureUnawareOptions{}
	InitClosureUnawareOptions(this_, original)
	return this_
}

func InitClosureUnawareOptions(this_ *ClosureUnawareOptions, original *CompilerOptions) {
	_ = original
	this_.initFields_ClosureUnawareOptions()
	this_.original = original
}

// Java static method: ClosureUnawareOptions.convert(com.google.javascript.jscomp.CompilerOptions original)
func ClosureUnawareOptions_Convert(original *CompilerOptions) *CompilerOptions {
	_ = original
	var val_1 ClosureUnawareOptions
	InitClosureUnawareOptions(&val_1, original)
	var options *ClosureUnawareOptions = &val_1
	_ = options
	return options.priv_ClosureUnawareOptions_toCompilerOptions_1247()
}

// Java method: ClosureUnawareOptions.toCompilerOptions()
func (this_ *ClosureUnawareOptions) priv_ClosureUnawareOptions_toCompilerOptions_1247() *CompilerOptions {
	this_.priv_ClosureUnawareOptions_setTranspilationOptions_1248()
	this_.priv_ClosureUnawareOptions_setSafeOptimizationAssumptions_1249()
	this_.priv_ClosureUnawareOptions_copyOutputOptions_1250()
	this_.priv_ClosureUnawareOptions_copyDebugOptions_1251()
	return this_.shadowOptions
}

// Java method: ClosureUnawareOptions.setTranspilationOptions()
func (this_ *ClosureUnawareOptions) priv_ClosureUnawareOptions_setTranspilationOptions_1248() {
	this_.shadowOptions.SetRewritePolyfills(false)
	switch (this_.original.GetRuntimeLibraryMode()).Ordinal() {
	case int32(1) /* RuntimeLibraryMode.NO_OP */:
		this_.shadowOptions.SetRuntimeLibraryMode(StaticGet_RuntimeJsLibManager_RuntimeLibraryMode_NO_OP())
	case int32(2) /* RuntimeLibraryMode.RECORD_ONLY */, int32(3) /* RuntimeLibraryMode.RECORD_AND_VALIDATE_FIELDS */, int32(0) /* RuntimeLibraryMode.INJECT */:
		this_.shadowOptions.SetRuntimeLibraryMode(StaticGet_RuntimeJsLibManager_RuntimeLibraryMode_EXTERN_FIELD_NAMES())
	case int32(4) /* RuntimeLibraryMode.EXTERN_FIELD_NAMES */:
		panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError")))
	}
	this_.shadowOptions.SetSkipNonTranspilationPasses(this_.original.GetSkipNonTranspilationPasses())
	this_.shadowOptions.SetEs6SubclassTranspilation(StaticGet_CompilerOptions_Es6SubclassTranspilation_SAFE_REFLECT_CONSTRUCT())
	this_.shadowOptions.SetOutputFeatureSet(this_.original.GetOutputFeatureSet())
	this_.shadowOptions.SetInstrumentAsyncContext(this_.original.GetInstrumentAsyncContext())
	this_.shadowOptions.SetWarningLevel(StaticGet_DiagnosticGroups_UNTRANSPILABLE_FEATURES(), StaticGet_CheckLevel_OFF())
}

// Java method: ClosureUnawareOptions.setSafeOptimizationAssumptions()
func (this_ *ClosureUnawareOptions) priv_ClosureUnawareOptions_setSafeOptimizationAssumptions_1249() {
	StaticGet_CompilationLevel_SIMPLE_OPTIMIZATIONS().SetOptionsForCompilationLevel(this_.shadowOptions)
	this_.shadowOptions.SetRemoveUnusedVariables(StaticGet_CompilerOptions_Reach_ALL())
	this_.shadowOptions.SetComputeFunctionSideEffects(false)
	this_.shadowOptions.SetCodingConvention(CodingConventions_GetDefault())
	this_.shadowOptions.SetRemoveClosureAsserts(false)
	this_.shadowOptions.SetAssumeGettersArePure(false)
	this_.shadowOptions.SetAssumePropertiesAreStaticallyAnalyzable(false)
	this_.shadowOptions.SetAssumeStaticInheritanceIsNotUsed(false)
	this_.shadowOptions.SetAssumeStrictThis(false)
	if ((this_.original.GetMaxFunctionSizeAfterInlining() != int32(-1))) {
		this_.shadowOptions.SetMaxFunctionSizeAfterInlining(this_.original.GetMaxFunctionSizeAfterInlining())
	}
	this_.shadowOptions.SetUseSizeHeuristicToStopOptimizationLoop(this_.original.ShouldUseSizeHeuristicToStopOptimizationLoop())
	this_.shadowOptions.SetMaxOptimizationLoopIterations(this_.original.GetMaxOptimizationLoopIterations())
	this_.shadowOptions.SetPreferStableNames(this_.original.ShouldPreferStableNames())
	this_.shadowOptions.SetInputPropertyMap(this_.original.GetInputPropertyMap())
	this_.shadowOptions.SetInputVariableMap(this_.original.GetInputVariableMap())
}

// Java method: ClosureUnawareOptions.copyOutputOptions()
func (this_ *ClosureUnawareOptions) priv_ClosureUnawareOptions_copyOutputOptions_1250() {
	this_.shadowOptions.SetEmitUseStrict(this_.original.ShouldEmitUseStrict())
	this_.shadowOptions.SetOutputCharset(this_.original.GetOutputCharset())
	this_.shadowOptions.SetGeneratePseudoNames(this_.original.ShouldGeneratePseudoNames())
	this_.shadowOptions.SetLineBreak(this_.original.ShouldAddLineBreak())
	this_.shadowOptions.SetLineLengthThreshold(this_.original.GetLineLengthThreshold())
	this_.shadowOptions.SetPreferSingleQuotes(this_.original.ShouldPreferSingleQuotes())
	this_.shadowOptions.SetPrettyPrint(this_.original.IsPrettyPrint())
	this_.shadowOptions.SetErrorFormat(this_.original.GetErrorFormat())
	this_.shadowOptions.SetErrorHandler(this_.original.GetErrorHandler())
	this_.shadowOptions.SetColorizeErrorOutput(this_.original.ShouldColorizeErrorOutput())
}

// Java method: ClosureUnawareOptions.copyDebugOptions()
func (this_ *ClosureUnawareOptions) priv_ClosureUnawareOptions_copyDebugOptions_1251() {
	this_.shadowOptions.SetTracerMode(this_.original.GetTracerMode())
	this_.shadowOptions.SetTracerOutput(this_.original.GetTracerOutput())
	this_.shadowOptions.SetDevMode(this_.original.GetDevMode())
	this_.shadowOptions.SetPrintSourceAfterEachPass(this_.original.ShouldPrintSourceAfterEachPass())
	this_.shadowOptions.SetFilesToPrintAfterEachPassRegexList(this_.original.GetFilesToPrintAfterEachPassRegexList())
	this_.shadowOptions.SetChunksToPrintAfterEachPassRegexList(this_.original.GetChunksToPrintAfterEachPassRegexList())
	this_.shadowOptions.SetPrintInputDelimiter(this_.original.ShouldPrintInputDelimiter())
	this_.shadowOptions.SetQnameUsesToPrintAfterEachPassList(this_.original.GetQnameUsesToPrintAfterEachPassList())
	this_.shadowOptions.SetInputDelimiter(this_.original.GetInputDelimiter())
	this_.shadowOptions.SetUseOriginalNamesInOutput(this_.original.GetUseOriginalNamesInOutput())
	this_.shadowOptions.SetDebugLogFilter(this_.original.GetDebugLogFilter())
	var debugLogDirectory j2gort.Ext_java_nio_file_Path = this_.original.GetDebugLogDirectory()
	_ = debugLogDirectory
	if ((debugLogDirectory != nil)) {
		this_.shadowOptions.SetDebugLogDirectory(debugLogDirectory.Resolve_1_string("./TranspileAndOptimizeClosureUnaware"))
	}
}

