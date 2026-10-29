// Port of test/com/google/javascript/jscomp/bundle/TranspilerTest.java @ closure-compiler v20261021

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
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// bundleTest_mockCompiler is the Mockito mock of
// CompilerBasedTransformer.CompilerSupplier used by TranspilerTest and
// CoverageInstrumenterTest (RETURNS_SMART_NULLS: an unstubbed call returns
// the empty value, "" for a String, null otherwise).
//
// Go-only: test-helper — Mockito.
type bundleTest_mockCompiler struct {
	compilePath   string
	compileCode   string
	compileResult *CompilerBasedTransformer_CompileResult
	runtimes      map[string]string
}

func (m *bundleTest_mockCompiler) Compile(path *jutil.Path, code string) *CompilerBasedTransformer_CompileResult {
	if path.String() == m.compilePath && code == m.compileCode {
		return m.compileResult
	}
	return nil
}

func (m *bundleTest_mockCompiler) Runtime(library string) string { return m.runtimes[library] }

var (
	transpilerTest_FOO_JS    = jutil.Path_Of("foo.js")
	transpilerTest_SOURCE_JS = jutil.Path_Of("source.js")
	transpilerTest_NO_ERRORS = []*jscomp.JSError{}
)

// Java: TranspilerTest.source(Path, String)
func transpilerTest_source(path *jutil.Path, code string) *Source {
	return Source_Builder_().SetPath(path).SetCode(code).Build()
}

func bundleTest_assertSource(t *testing.T, got *Source, want *Source) {
	t.Helper()
	if !jutil.ObjectsEquals(got, want) {
		t.Errorf("transform = %v, want %v", got, want)
	}
}

// Tests for Transpiler

func TestTranspiler_Transpiler_transpile(t *testing.T) {
	mockCompiler := &bundleTest_mockCompiler{runtimes: map[string]string{}}
	var transpiler Source_Transformer = NewTranspiler(mockCompiler, "es6_runtime")
	mockCompiler.runtimes["es6_runtime"] = "$jscomp.es6();"
	mockCompiler.compilePath, mockCompiler.compileCode = "foo.js", "bar"
	mockCompiler.compileResult = NewCompilerBasedTransformer_CompileResult("result", transpilerTest_NO_ERRORS, true, "srcmap")
	bundleTest_assertSource(t, transpiler.Transform(transpilerTest_source(transpilerTest_FOO_JS, "bar")),
		Source_Builder_().
			SetPath(transpilerTest_FOO_JS).
			SetOriginalCode("bar").
			SetCode("result").
			SetSourceMap("srcmap").
			AddRuntime("$jscomp.es6();").
			Build())
}

func TestTranspiler_Transpiler_noTranspilation(t *testing.T) {
	mockCompiler := &bundleTest_mockCompiler{runtimes: map[string]string{}}
	var transpiler Source_Transformer = NewTranspiler(mockCompiler, "es6_runtime")
	mockCompiler.compilePath, mockCompiler.compileCode = "foo.js", "bar"
	mockCompiler.compileResult = NewCompilerBasedTransformer_CompileResult("result", transpilerTest_NO_ERRORS, false, "srcmap")
	bundleTest_assertSource(t, transpiler.Transform(transpilerTest_source(transpilerTest_FOO_JS, "bar")),
		transpilerTest_source(transpilerTest_FOO_JS, "bar"))
}

// Tests for CompilerSupplier

func TestTranspiler_CompilerSupplier_compileChanged(t *testing.T) {
	compiler := Transpiler_CompilerSupplier()
	result := compiler.Compile(transpilerTest_SOURCE_JS, "const x = () => 42;")
	if result.source != "var x = function() {\n  return 42;\n};\n" || len(result.errors) != 0 || !result.transformed ||
		!strings.Contains(result.sourceMap, "\"mappings\":\"AAAA,IAAMA,IAAIA,QAAA,EAAM;AAAA,SAAA,EAAA;AAAA,CAAhB;;\"") {
		t.Errorf("compile = %q %v %v %q", result.source, result.errors, result.transformed, result.sourceMap)
	}
}

func TestTranspiler_CompilerSupplier_compileNoChange(t *testing.T) {
	compiler := Transpiler_CompilerSupplier()
	result := compiler.Compile(transpilerTest_SOURCE_JS, "var x = 42;")
	if result.source != "var x = 42;\n" || len(result.errors) != 0 || result.transformed || result.sourceMap != "" {
		t.Errorf("compile = %q %v %v %q", result.source, result.errors, result.transformed, result.sourceMap)
	}
}

func TestTranspiler_CompilerSupplier_runtime(t *testing.T) {
	compiler := Transpiler_CompilerSupplier()
	runtime := compiler.Runtime("es6_runtime")
	for _, want := range []string{"$jscomp.polyfill(\"Map\"", "$jscomp.makeIterator", "$jscomp.inherits"} {
		if !strings.Contains(runtime, want) {
			t.Errorf("runtime does not contain %q", want)
		}
	}
}
