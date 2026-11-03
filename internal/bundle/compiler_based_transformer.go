// Port of com/google/javascript/jscomp/bundle/CompilerBasedTransformer.java @ closure-compiler v20261028
// Port of com/google/javascript/jscomp/bundle/CoverageInstrumenter.java @ closure-compiler v20261028

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

package bundle

import (
	"strings"
	"sync"

	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CompilerBasedTransformer is a source transformer base struct. May also
// include a runtime that needs to be shipped with the final bundle.
//
// Java: abstract class bundle.CompilerBasedTransformer implements Source.Transformer
type CompilerBasedTransformer struct {
	compilerSupplier   CompilerSupplier
	transformationName string
	runtimeLibraryName string
	hasRuntime         bool
}

var _ Source_Transformer = (*CompilerBasedTransformer)(nil)

// GetRuntime returns the runtime library name and whether one is present.
//
// Java: CompilerBasedTransformer.getRuntime()
func (t *CompilerBasedTransformer) GetRuntime() (string, bool) {
	return t.runtimeLibraryName, t.hasRuntime
}

// Java: CompilerBasedTransformer.transform(Source)
func (t *CompilerBasedTransformer) Transform(input *Source) *Source {
	result := t.compilerSupplier.Compile(input.Path().String(), input.Code())
	if len(result.Errors) != 0 {
		// TODO(sdh): how to handle this?  Currently we throw an ISE with the
		// message, but this may not be the most appropriate option.  It
		// might make sense to add console.log() statements to any JS that
		// comes out, particularly for warnings.
		formatter := jscomp.ErrorFormat_SOURCELESS.ToFormatter(nil, false)
		var message strings.Builder
		message.WriteString(t.transformationName)
		message.WriteString(" failed.\n")
		for _, err := range result.Errors {
			message.WriteString(formatter.FormatError(err))
		}
		panic(jutil.NewIllegalStateExceptionMsg(message.String()))
	}
	if !result.Transformed {
		return input
	}
	if t.hasRuntime {
		return input.WithTransformed(result.Source, result.SourceMap, t.compilerSupplier.Runtime(t.runtimeLibraryName))
	}
	return input.WithTransformed(result.Source, result.SourceMap)
}

// CompilerSupplier wraps the Compiler into a small interface, making it easy to
// test per-file transformers and transpilers without depending on implementation
// details of the Compiler itself, while supplying a fresh Compiler instance on
// each call in single-file mode.
type CompilerSupplier interface {
	Compile(path string, code string) *CompileResult
	Runtime(library string) string
}

// CompileResult holds the emitted JS source together with single-file
// compilation metadata.
type CompileResult struct {
	Source      string
	Errors      []*jscomp.JSError
	Transformed bool
	SourceMap   string
}

// NewCompileResult constructs a CompileResult.
func NewCompileResult(source string, errors []*jscomp.JSError, transformed bool, sourceMap string) *CompileResult {
	return &CompileResult{
		Source:      source,
		Errors:      errors,
		Transformed: transformed,
		SourceMap:   sourceMap,
	}
}

const devNullSourceMapPath = "/dev/null"

var (
	compilerSupplierExterns = jscomp.SourceFile_FromCode("externs.js", "function Symbol() {}")
	compilerSupplierEmpty   = jscomp.SourceFile_FromCode("empty.js", "")
	es5WarningsGroup        = jscomp.NewDiagnosticGroup(jscomp.DiagnosticType_Error("JSC_CANNOT_CONVERT", ""))
)

// OptionsCompilerSupplier implements CompilerSupplier by configuring a fresh
// CompilerOptions instance per compilation or runtime library request.
type OptionsCompilerSupplier struct {
	SetOptions        func(*jscomp.CompilerOptions)
	AlwaysTransformed bool
	ThrowOnError      bool
}

var _ CompilerSupplier = (*OptionsCompilerSupplier)(nil)

func (s *OptionsCompilerSupplier) options() *jscomp.CompilerOptions {
	options := jscomp.NewCompilerOptions()
	if s.SetOptions != nil {
		s.SetOptions(options)
	}
	return options
}

// compileMu serializes standalone single-file compilations so concurrent
// CachedTransformer / CachingTranspiler cache misses never race on the
// process-wide rhino NodeStore.
var compileMu sync.Mutex

// Compile compiles a single JS source file with the shared Symbol() externs.
func (s *OptionsCompilerSupplier) Compile(path string, code string) *CompileResult {
	compileMu.Lock()
	defer compileMu.Unlock()
	if rhino.CurrentNodeStore().IsStatic() {
		defer rhino.ResetNodeStore()
	}
	compiler := jscomp.NewCompiler()
	result := compiler.CompileOfSourceFiles(
		compilerSupplierExterns,
		jscomp.SourceFile_FromCode(path, code),
		s.options())
	if s.ThrowOnError && len(result.Errors) != 0 {
		panic(NewTranspilationException(compiler, result.Errors, result.Warnings))
	}
	source := compiler.ToSource()
	var smBuilder strings.Builder
	if result.SourceMap != nil {
		result.SourceMap.AppendTo(&smBuilder, path)
	}
	transformed := s.AlwaysTransformed || result.TranspiledFiles
	var sourceMap string
	if transformed {
		sourceMap = smBuilder.String()
	}
	return NewCompileResult(source, result.Errors, transformed, sourceMap)
}

// Runtime compiles an empty input with library force-injected and returns the
// emitted runtime JS code.
func (s *OptionsCompilerSupplier) Runtime(library string) string {
	compileMu.Lock()
	defer compileMu.Unlock()
	if rhino.CurrentNodeStore().IsStatic() {
		defer rhino.ResetNodeStore()
	}
	compiler := jscomp.NewCompiler()
	options := s.options()
	options.SetForceLibraryInjection([]string{library})
	compiler.CompileOfSourceFiles(
		compilerSupplierExterns,
		compilerSupplierEmpty,
		options)
	return compiler.ToSource()
}

// SetCommonTranspilationOptions configures the CompilerOptions shared by
// single-file ES transpilers in bundle and transpile.
func SetCommonTranspilationOptions(options *jscomp.CompilerOptions) {
	options.SetQuoteKeywordProperties(true)
	options.SetSkipNonTranspilationPasses(true)
	options.SetVariableRenaming(jscomp.VariableRenamingPolicy_OFF)
	options.SetPropertyRenaming(jscomp.PropertyRenamingPolicy_OFF)
	options.SetWrapGoogModulesForWhitespaceOnly(false)
	options.SetPrettyPrint(true)
	options.SetSourceMapOutputPath(devNullSourceMapPath)
	options.SetSourceMapIncludeSourcesContent(true)
}

// NewCompilerSupplier returns the default ES5 CompilerSupplier.
func NewCompilerSupplier() CompilerSupplier {
	return &OptionsCompilerSupplier{
		SetOptions: setEs5TranspilerOptions,
	}
}

func setEs5TranspilerOptions(options *jscomp.CompilerOptions) {
	// TODO(sdh): It would be nice to allow people to output code in
	// strict mode.  But currently we swallow all the input language
	// strictness checks, and there are various tests that are never
	// compiled and so are broken when we output 'use strict'.  We
	// could consider adding some sort of logging/warning/error in
	// cases where the input was not strict, though there could still
	// be semantic differences even if syntax is strict.  Possibly
	// the first step would be to allow the option of outputting strict
	// and then change the default and see what breaks.  b/33005948
	options.SetLanguageOut(jscomp.CompilerOptions_LanguageMode_ECMASCRIPT5)
	SetCommonTranspilationOptions(options)
	options.SetWarningLevel(es5WarningsGroup, jscomp.CheckLevel_OFF)
}

// CoverageInstrumenter is a source transformer for instrumenting code for
// coverage data collection.
//
// Java: class bundle.CoverageInstrumenter extends CompilerBasedTransformer
type CoverageInstrumenter struct {
	CompilerBasedTransformer
}

var _ Source_Transformer = (*CoverageInstrumenter)(nil)

// Java: CoverageInstrumenter(CompilerBasedTransformer.CompilerSupplier)
func NewCoverageInstrumenter(compilerSupplier CompilerSupplier) *CoverageInstrumenter {
	return &CoverageInstrumenter{
		CompilerBasedTransformer: CompilerBasedTransformer{
			compilerSupplier:   jutil.CheckNotNull(compilerSupplier),
			transformationName: "Coverage Instrumentation",
		},
	}
}

// NewCoverageCompilerSupplier supplies options for coverage instrumentation.
func NewCoverageCompilerSupplier() CompilerSupplier {
	return &OptionsCompilerSupplier{
		SetOptions:        setCoverageOptions,
		AlwaysTransformed: true,
	}
}

func setCoverageOptions(options *jscomp.CompilerOptions) {
	options.SetCoalesceVariableNames(false)
	options.SetLanguageOut(jscomp.CompilerOptions_LanguageMode_ECMASCRIPT5)
	options.SetLanguageIn(jscomp.CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	options.SetStrictModeInput(false)
	// Setting the path to any non-null value will trigger source map
	// generation. CompilerBasedTransformer attaches the sourcemap to the
	// result.
	options.SetSourceMapOutputPath(devNullSourceMapPath)
	options.SetVariableRenaming(jscomp.VariableRenamingPolicy_OFF)
	options.SetInstrumentForCoverageOption(jscomp.CompilerOptions_InstrumentOption_LINE_ONLY)
	options.SetInstrumentForCoverageOnly(true)
}
