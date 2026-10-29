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
	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

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
