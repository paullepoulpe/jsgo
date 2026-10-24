// Port of test/com/google/javascript/jscomp/bundle/CoverageInstrumenterTest.java @ closure-compiler v20261021

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
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/sourcemap"
)

// Tests for CoverageInstrumenter

func TestCoverageInstrumenter_CoverageInstrumenter_instrument(t *testing.T) {
	mockCompiler := &bundleTest_mockCompiler{runtimes: map[string]string{}}
	var instrumenter Source_Transformer = NewCoverageInstrumenter(mockCompiler)
	mockCompiler.compilePath, mockCompiler.compileCode = "foo.js", "bar"
	mockCompiler.compileResult = NewCompilerBasedTransformer_CompileResult("result", transpilerTest_NO_ERRORS, true, "srcmap")
	bundleTest_assertSource(t, instrumenter.Transform(transpilerTest_source(transpilerTest_FOO_JS, "bar")),
		Source_Builder_().
			SetPath(transpilerTest_FOO_JS).
			SetOriginalCode("bar").
			SetCode("result").
			SetSourceMap("srcmap").
			Build())
}

func TestCoverageInstrumenter_CoverageInstrumenter_noInstrumentation(t *testing.T) {
	mockCompiler := &bundleTest_mockCompiler{runtimes: map[string]string{}}
	var instrumenter Source_Transformer = NewCoverageInstrumenter(mockCompiler)
	mockCompiler.compilePath, mockCompiler.compileCode = "foo.js", "bar"
	mockCompiler.compileResult = NewCompilerBasedTransformer_CompileResult("result", transpilerTest_NO_ERRORS, false, "srcmap")
	bundleTest_assertSource(t, instrumenter.Transform(transpilerTest_source(transpilerTest_FOO_JS, "bar")),
		transpilerTest_source(transpilerTest_FOO_JS, "bar"))
}

// Tests for CompilerSupplier

func TestCoverageInstrumenter_CompilerSupplier_instruments(t *testing.T) {
	compiler := CoverageInstrumenter_CompilerSupplier_()
	result := compiler.Compile(transpilerTest_SOURCE_JS, "var x = 42;")
	expected := []string{
		"(function(self){if(!self.window){self.window=self;self.window.top=self}})(typeof " +
			"self!==\"undefined\"?self:globalThis);",
		"var __jscov=window.top[\"__jscov\"]||",
		"(window.top[\"__jscov\"]={\"fileNames\":[],\"instrumentedLines\":[]," +
			"\"executedLines\":[]});",
		"var JSCompiler_lcov_data_source_js=[];",
		"__jscov[\"executedLines\"].push(JSCompiler_lcov_data_source_js);",
		"__jscov[\"instrumentedLines\"].push(\"01\");",
		"__jscov[\"fileNames\"].push(\"source.js\");",
		"JSCompiler_lcov_data_source_js[0]=true;var x=42;",
	}

	if want := strings.Join(expected, ""); result.source != want {
		t.Errorf("source = %q, want %q", result.source, want)
	}
	if len(result.errors) != 0 || !result.transformed {
		t.Errorf("errors %v, transformed %v", result.errors, result.transformed)
	}

	// Ensure that source map for "x" is correct.
	sourceMap := sourcemap.NewSourceMapConsumerV3()
	sourceMap.Parse(result.sourceMap)
	mappingFoX := sourceMap.GetMappingForLine(1, jutil.StringIndexOf(result.source, "x=42")+1)
	want := sourcemap.Mapping_OriginalMapping_NewBuilder().
		SetOriginalFile(transpilerTest_SOURCE_JS.String()).
		SetLineNumber(1).
		SetColumnPosition(5).
		SetIdentifier("x").
		SetPrecision(sourcemap.Mapping_OriginalMapping_Precision_EXACT).
		Build()
	if !jutil.ObjectsEquals(mappingFoX, want) {
		t.Errorf("mapping = %v, want %v", mappingFoX, want)
	}
}
