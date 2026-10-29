// Port of com/google/javascript/jscomp/bundle/CompilerBasedTransformer.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/bundle/CoverageInstrumenter.java @ closure-compiler v20261021

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

	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// CompilerBasedTransformer is a source transformer base struct. May also
// include a runtime that needs to be shipped with the final bundle.
//
// Java: abstract class bundle.CompilerBasedTransformer implements Source.Transformer
type CompilerBasedTransformer struct {
	compilerSupplier   CompilerBasedTransformer_CompilerSupplier
	transformationName string
	runtimeLibraryName string
	hasRuntime         bool
}

var _ Source_Transformer = (*CompilerBasedTransformer)(nil)

// GetRuntime returns nil for Optional.empty().
//
// Java: CompilerBasedTransformer.getRuntime()
func (t *CompilerBasedTransformer) GetRuntime() *string {
	if !t.hasRuntime {
		return nil
	}
	s := t.runtimeLibraryName
	return &s
}

// Java: CompilerBasedTransformer.getTranformationName()
func (t *CompilerBasedTransformer) GetTranformationName() string {
	return t.transformationName
}

// Java: CompilerBasedTransformer.transform(Source)
func (t *CompilerBasedTransformer) Transform(input *Source) *Source {
	result := t.compilerSupplier.Compile(input.Path(), input.Code())
	if len(result.errors) != 0 {
		// TODO(sdh): how to handle this?  Currently we throw an ISE with the
		// message, but this may not be the most appropriate option.  It
		// might make sense to add console.log() statements to any JS that
		// comes out, particularly for warnings.
		formatter := jscomp.ErrorFormat_SOURCELESS.ToFormatter(nil, false)
		var message strings.Builder
		message.WriteString(t.transformationName)
		message.WriteString(" failed.\n")
		for _, error := range result.errors {
			message.WriteString(formatter.FormatError(error))
		}
		panic(jutil.NewIllegalStateExceptionMsg(message.String()))
	}
	if !result.transformed {
		return input
	}
	builder := input.ToBuilder().
		SetCode(result.source).
		SetSourceMap(result.sourceMap)
	if t.hasRuntime {
		builder.AddRuntime(t.compilerSupplier.Runtime(t.runtimeLibraryName))
	}
	return builder.Build()
}

// CompilerBasedTransformer_CompilerSupplier wraps the Compiler into a small
// interface, making it easy to test the CompilerBasedTransformer without
// depending on implementation details of the Compiler itself. Also works
// around the fact that the Compiler is not thread-safe (since we may do
// multiple transpiles concurrently), so we supply a fresh instance each time
// when we're in single-file mode.
//
// Java: static class CompilerBasedTransformer.CompilerSupplier
type CompilerBasedTransformer_CompilerSupplier interface {
	// Java: CompilerSupplier.compile(Path, String)
	Compile(path *jutil.Path, code string) *CompilerBasedTransformer_CompileResult
	// Java: CompilerSupplier.runtime(String)
	Runtime(library string) string
}

// CompilerBasedTransformer_DefaultCompilerSupplier is the default ES5 CompilerSupplier.
type CompilerBasedTransformer_DefaultCompilerSupplier struct{}

var _ CompilerBasedTransformer_CompilerSupplier = (*CompilerBasedTransformer_DefaultCompilerSupplier)(nil)

// Java: CompilerSupplier()
func NewCompilerBasedTransformer_CompilerSupplier() *CompilerBasedTransformer_DefaultCompilerSupplier {
	return &CompilerBasedTransformer_DefaultCompilerSupplier{}
}

func compileSingleFileWithOptions(
	path *jutil.Path,
	code string,
	options *jscomp.CompilerOptions,
	alwaysTransformed bool) *CompilerBasedTransformer_CompileResult {
	compiler := jscomp.NewCompiler()
	result := compiler.CompileOfSourceFiles(
		CompilerBasedTransformer_CompilerSupplier_EXTERNS,
		jscomp.SourceFile_FromCode(path.String(), code),
		options)
	source := compiler.ToSource()
	var sourceMap strings.Builder
	if result.SourceMap != nil {
		name := path.String()
		result.SourceMap.AppendTo(&sourceMap, &name)
	}
	transformed := alwaysTransformed || result.TranspiledFiles
	sm := ""
	if transformed {
		sm = sourceMap.String()
	}
	return NewCompilerBasedTransformer_CompileResult(
		source, result.Errors, transformed, sm)
}

func runtimeSingleFileWithOptions(options *jscomp.CompilerOptions, library string) string {
	compiler := jscomp.NewCompiler()
	options.SetForceLibraryInjection([]string{library})
	compiler.CompileOfSourceFiles(
		CompilerBasedTransformer_CompilerSupplier_EXTERNS,
		CompilerBasedTransformer_CompilerSupplier_EMPTY,
		options)
	return compiler.ToSource()
}

// Java: CompilerSupplier.compile(Path, String)
func (s *CompilerBasedTransformer_DefaultCompilerSupplier) Compile(path *jutil.Path, code string) *CompilerBasedTransformer_CompileResult {
	return compileSingleFileWithOptions(path, code, s.Options(), false)
}

// Java: CompilerSupplier.transformed(Result)
func (s *CompilerBasedTransformer_DefaultCompilerSupplier) Transformed(result *jscomp.Result) bool {
	return result.TranspiledFiles
}

// Java: CompilerSupplier.runtime(String)
func (s *CompilerBasedTransformer_DefaultCompilerSupplier) Runtime(library string) string {
	return runtimeSingleFileWithOptions(s.Options(), library)
}

// Java: CompilerSupplier.compiler()
func (s *CompilerBasedTransformer_DefaultCompilerSupplier) Compiler() *jscomp.Compiler {
	return jscomp.NewCompiler()
}

// Java: CompilerSupplier.options()
func (s *CompilerBasedTransformer_DefaultCompilerSupplier) Options() *jscomp.CompilerOptions {
	options := jscomp.NewCompilerOptions()
	s.SetOptions(options)
	return options
}

// Java: CompilerSupplier.setOptions(CompilerOptions)
func (s *CompilerBasedTransformer_DefaultCompilerSupplier) SetOptions(options *jscomp.CompilerOptions) {
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
	options.SetQuoteKeywordProperties(true)
	options.SetSkipNonTranspilationPasses(true)
	options.SetVariableRenaming(jscomp.VariableRenamingPolicy_OFF)
	options.SetPropertyRenaming(jscomp.PropertyRenamingPolicy_OFF)
	options.SetWrapGoogModulesForWhitespaceOnly(false)
	options.SetPrettyPrint(true)
	devNull := "/dev/null"
	options.SetSourceMapOutputPath(&devNull)
	options.SetSourceMapIncludeSourcesContent(true)
	options.SetWarningLevel(CompilerBasedTransformer_CompilerSupplier_ES5_WARNINGS, jscomp.CheckLevel_OFF)
}

// Java: CompilerSupplier.EXTERNS
var CompilerBasedTransformer_CompilerSupplier_EXTERNS = jscomp.SourceFile_FromCode("externs.js", "function Symbol() {}")

// Java: CompilerSupplier.EMPTY
var CompilerBasedTransformer_CompilerSupplier_EMPTY = jscomp.SourceFile_FromCode("empty.js", "")

// Java: CompilerSupplier.ES5_WARNINGS
var CompilerBasedTransformer_CompilerSupplier_ES5_WARNINGS = jscomp.NewDiagnosticGroup(jscomp.DiagnosticType_Error("JSC_CANNOT_CONVERT", ""))

// CompilerBasedTransformer_CompileResult is the source together with the
// additional compilation results.
//
// Java: static class CompilerBasedTransformer.CompileResult
type CompilerBasedTransformer_CompileResult struct {
	source      string
	errors      []*jscomp.JSError
	transformed bool
	sourceMap   string
}

// Java: CompileResult(String, ImmutableList<JSError>, boolean, String)
func NewCompilerBasedTransformer_CompileResult(
	source string, errors []*jscomp.JSError, transformed bool, sourceMap string) *CompilerBasedTransformer_CompileResult {
	return &CompilerBasedTransformer_CompileResult{
		source:      source,
		errors:      errors, // checkNotNull: a Go slice is never null (nil is an empty list)
		transformed: transformed,
		sourceMap:   sourceMap,
	}
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
func NewCoverageInstrumenter(compilerSupplier CompilerBasedTransformer_CompilerSupplier) *CoverageInstrumenter {
	return &CoverageInstrumenter{
		CompilerBasedTransformer: CompilerBasedTransformer{
			compilerSupplier:   jutil.CheckNotNull(compilerSupplier),
			transformationName: "Coverage Instrumentation",
		},
	}
}

// CoverageInstrumenter_CompilerSupplier supplies options for coverage.
//
// Java: static class CoverageInstrumenter.CompilerSupplier extends CompilerBasedTransformer.CompilerSupplier
type CoverageInstrumenter_CompilerSupplier struct{}

var _ CompilerBasedTransformer_CompilerSupplier = (*CoverageInstrumenter_CompilerSupplier)(nil)

// Java: CoverageInstrumenter.CompilerSupplier()
func NewCoverageInstrumenter_CompilerSupplier() *CoverageInstrumenter_CompilerSupplier {
	return &CoverageInstrumenter_CompilerSupplier{}
}

// Java: CompilerSupplier.compile(Path, String)
func (s *CoverageInstrumenter_CompilerSupplier) Compile(path *jutil.Path, code string) *CompilerBasedTransformer_CompileResult {
	return compileSingleFileWithOptions(path, code, s.Options(), true)
}

// Java: CompilerSupplier.runtime(String)
func (s *CoverageInstrumenter_CompilerSupplier) Runtime(library string) string {
	return runtimeSingleFileWithOptions(s.Options(), library)
}

// Java: CompilerSupplier.options()
func (s *CoverageInstrumenter_CompilerSupplier) Options() *jscomp.CompilerOptions {
	options := jscomp.NewCompilerOptions()
	s.SetOptions(options)
	return options
}

// Java: CoverageInstrumenter.CompilerSupplier.setOptions(CompilerOptions)
func (s *CoverageInstrumenter_CompilerSupplier) SetOptions(options *jscomp.CompilerOptions) {
	options.SetCoalesceVariableNames(false)
	options.SetLanguageOut(jscomp.CompilerOptions_LanguageMode_ECMASCRIPT5)
	options.SetLanguageIn(jscomp.CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	options.SetStrictModeInput(false)
	// Setting the path to any non-null value will trigger source map
	// generation. CompilerBasedTransformer attachs the sourcemap to the
	// result.
	devNull := "/dev/null"
	options.SetSourceMapOutputPath(&devNull)
	options.SetVariableRenaming(jscomp.VariableRenamingPolicy_OFF)
	options.SetInstrumentForCoverageOption(jscomp.CompilerOptions_InstrumentOption_LINE_ONLY)
	options.SetInstrumentForCoverageOnly(true)
}

// Java: CoverageInstrumenter.CompilerSupplier.transformed(Result)
func (s *CoverageInstrumenter_CompilerSupplier) Transformed(result *jscomp.Result) bool {
	return true
}

// CoverageInstrumenter_CompilerSupplier_ is CoverageInstrumenter's static
// compilerSupplier() (the trailing underscore keeps it apart from the type
// CoverageInstrumenter_CompilerSupplier).
//
// Java: CoverageInstrumenter.compilerSupplier()
func CoverageInstrumenter_CompilerSupplier_() *CoverageInstrumenter_CompilerSupplier {
	return NewCoverageInstrumenter_CompilerSupplier()
}

// Java: CoverageInstrumenter.INSTRUMENTER
var CoverageInstrumenter_INSTRUMENTER = NewCoverageInstrumenter(CoverageInstrumenter_CompilerSupplier_())
