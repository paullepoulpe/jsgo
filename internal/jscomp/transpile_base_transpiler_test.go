// Port of test/com/google/javascript/jscomp/transpile/BaseTranspilerTest.java @ closure-compiler v20260415

/*
 * Copyright 2016 The Closure Compiler Authors.
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
)

// baseTranspilerTest_mockCompiler is the Mockito mock of
// BaseTranspiler.CompilerSupplier (RETURNS_SMART_NULLS: an unstubbed call
// returns the empty value, "" for a String, null otherwise).
//
// Go-only: test-helper — Mockito.
type baseTranspilerTest_mockCompiler struct {
	*BaseTranspiler_CompilerSupplierBase
	compileResult *BaseTranspiler_CompileResult
	runtimes      map[string]string
}

func (m *baseTranspilerTest_mockCompiler) Compile(path *jutil.URI, code string) *BaseTranspiler_CompileResult {
	if path.String() == baseTranspilerTest_FOO_JS.String() && code == "bar" {
		return m.compileResult
	}
	return nil
}

func (m *baseTranspilerTest_mockCompiler) Runtime(library string) string {
	return m.runtimes[library]
}

var (
	baseTranspilerTest_FOO_JS    = jutil.NewURI("foo.js")
	baseTranspilerTest_SOURCE_JS = jutil.NewURI("source.js")
)

type baseTranspilerTest struct {
	transpiler   Transpile_Transpiler
	compiler     BaseTranspiler_CompilerSupplier
	mockCompiler *baseTranspilerTest_mockCompiler
}

// Java: BaseTranspilerTest.setUp()
func newBaseTranspilerTest() *baseTranspilerTest {
	x := &baseTranspilerTest{}
	x.mockCompiler = &baseTranspilerTest_mockCompiler{runtimes: map[string]string{}}
	x.transpiler = NewBaseTranspiler(x.mockCompiler, "es6_runtime")
	x.compiler = NewBaseTranspiler_CompilerSupplier()
	return x
}

func baseTranspilerTest_assertResult(t *testing.T, got *TranspileResult, want *TranspileResult) {
	t.Helper()
	if !jutil.ObjectsEquals(got, want) {
		t.Errorf("transpile = %v, want %v", got, want)
	}
}

// Tests for BaseTranspiler

func TestBaseTranspiler_Transpiler_transpile(t *testing.T) {
	x := newBaseTranspilerTest()
	x.mockCompiler.compileResult = NewBaseTranspiler_CompileResult("result", true, "srcmap")
	baseTranspilerTest_assertResult(t, x.transpiler.Transpile(baseTranspilerTest_FOO_JS, "bar"),
		NewTranspileResult(baseTranspilerTest_FOO_JS, "bar", "result", "srcmap"))
}

func TestBaseTranspiler_Transpiler_noTranspilation(t *testing.T) {
	x := newBaseTranspilerTest()
	x.mockCompiler.compileResult = NewBaseTranspiler_CompileResult("result", false, "srcmap")
	baseTranspilerTest_assertResult(t, x.transpiler.Transpile(baseTranspilerTest_FOO_JS, "bar"),
		NewTranspileResult(baseTranspilerTest_FOO_JS, "bar", "bar", ""))
}

func TestBaseTranspiler_Transpiler_runtime(t *testing.T) {
	x := newBaseTranspilerTest()
	x.mockCompiler.runtimes["es6_runtime"] = "$jscomp.es6();"
	if got := x.transpiler.Runtime(); got != "$jscomp.es6();" {
		t.Errorf("runtime() = %q", got)
	}
}

// Tests for CompilerSupplier

func TestBaseTranspiler_CompilerSupplier_compileChanged(t *testing.T) {
	x := newBaseTranspilerTest()
	result := x.compiler.Compile(baseTranspilerTest_SOURCE_JS, "const x = () => 42;")
	if result.source != "var x = function() {\n  return 42;\n};\n" || !result.transpiled ||
		!strings.Contains(result.sourceMap, "\"mappings\":\"AAAA,IAAMA,IAAIA,QAAA,EAAM;AAAA,SAAA,EAAA;AAAA,CAAhB;;\"") {
		t.Errorf("compile = %q %v %q", result.source, result.transpiled, result.sourceMap)
	}
}

func TestBaseTranspiler_CompilerSupplier_compileNoChange(t *testing.T) {
	x := newBaseTranspilerTest()
	result := x.compiler.Compile(baseTranspilerTest_SOURCE_JS, "var x = 42;")
	if result.source != "var x = 42;\n" || result.transpiled || result.sourceMap != "" {
		t.Errorf("compile = %q %v %q", result.source, result.transpiled, result.sourceMap)
	}
}

func TestBaseTranspiler_CompilerSupplier_error(t *testing.T) {
	x := newBaseTranspilerTest()
	thrown := false
	jutil.Try(func() {
		x.compiler.Compile(baseTranspilerTest_SOURCE_JS, "cons x = () => 42;")
	}, func(expected *TranspilationException) {
		thrown = true
	})
	if !thrown {
		t.Errorf("Expected an exception.")
	}
}

func TestBaseTranspiler_CompilerSupplier_runtime(t *testing.T) {
	x := newBaseTranspilerTest()
	runtime := x.compiler.Runtime("es6_runtime")
	for _, want := range []string{"$jscomp.polyfill(\"Map\"", "$jscomp.makeIterator", "$jscomp.inherits"} {
		if !strings.Contains(runtime, want) {
			t.Errorf("runtime does not contain %q", want)
		}
	}
}

func TestBaseTranspiler_CompilerSupplier_inlineSourceMap(t *testing.T) {
	x := newBaseTranspilerTest()
	// This is a source map generated by
	// tsc --inlineSourceMap --inlineSources -t ES2017 source.ts
	// and the content of source.ts is "const x = () => 42".
	result := x.compiler.Compile(
		baseTranspilerTest_SOURCE_JS,
		"const x = () => 42;\n"+
			"//# sourceMappingURL=data:application/json;base64,"+
			"eyJ2ZXJzaW9uIjozLCJmaWxlIjoic291cmNlLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsi"+
			"c291cmNlLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiJBQUFBLE1BQU0sQ0FBQyxHQUFHLEdBQUcs"+
			"RUFBRSxDQUFDLEVBQUUsQ0FBQyIsInNvdXJjZXNDb250ZW50IjpbImNvbnN0IHggPSAoKSA9PiA0Mjtc"+
			"biJdfQ==")
	if result.source != "var x = function() {\n  return 42;\n};\n" || !result.transpiled {
		t.Errorf("compile = %q %v", result.source, result.transpiled)
	}
	// This should map back to the typescript source file.
	want := "{\n" +
		"\"version\":3,\n" +
		"\"file\":\"/source.js\",\n" +
		"\"lineCount\":4,\n" +
		"\"mappings\":\"AAAA,IAAMA,IAAIA,QAAA,EAAM;AAAA,SAAA,EAAA;AAAA,CAAhB;;\",\n" +
		"\"sources\":[\"/source.ts\"],\n" +
		"\"sourcesContent\":[\"const x = () => 42;\\n\"],\n" +
		"\"names\":[\"x\"]\n" +
		"}\n"
	if result.sourceMap != want {
		t.Errorf("sourceMap = %q, want %q", result.sourceMap, want)
	}
}
