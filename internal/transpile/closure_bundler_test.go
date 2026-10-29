// Port of test/com/google/javascript/jscomp/deps/ClosureBundlerTest.java @ closure-compiler v20261021

/*
 * Copyright 2014 The Closure Compiler Authors.
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

package transpile

import (
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/bundle"
	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

var closureBundlerTest_MODULE jscomp.DependencyInfo = jscomp.SimpleDependencyInfo_Builder_("", "").SetGoogModule(true).Build()

var closureBundlerTest_TRADITIONAL jscomp.DependencyInfo = jscomp.SimpleDependencyInfo_Builder_("", "").Build()

func closureBundlerTest_assertEquals(t *testing.T, got string, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got  %q\nwant %q", got, want)
	}
}

func TestClosureBundler_GoogModule(t *testing.T) {
	sb := jutil.NewStringWriter()
	NewClosureBundler().AppendTo(sb, closureBundlerTest_MODULE, "\"a string\"")
	closureBundlerTest_assertEquals(t, sb.String(),
		"goog.loadModule(function(exports) {'use strict';\"a string\"\n;return exports;});\n")
}

func TestClosureBundler_GoogModuleWithSourceURL(t *testing.T) {
	sb := jutil.NewStringWriter()
	NewClosureBundler().UseEval(true).WithSourceUrl("URL").AppendTo(sb, closureBundlerTest_MODULE, "\"a string\"")
	closureBundlerTest_assertEquals(t, sb.String(),
		"goog.loadModule(\"\\x22a string\\x22\\n//# sourceURL\\x3dURL\\n\");\n")
}

func TestClosureBundler_GoogModuleWithEval(t *testing.T) {
	sb := jutil.NewStringWriter()
	NewClosureBundler().UseEval(true).AppendTo(sb, closureBundlerTest_MODULE, "\"a string\"")
	closureBundlerTest_assertEquals(t, sb.String(), "goog.loadModule(\"\\x22a string\\x22\");\n")
}

func TestClosureBundler_GoogModuleWithEvalWithURL(t *testing.T) {
	sb := jutil.NewStringWriter()
	NewClosureBundler().UseEval(true).WithSourceUrl("URL").AppendTo(sb, closureBundlerTest_MODULE, "\"a string\"")
	closureBundlerTest_assertEquals(t, sb.String(),
		"goog.loadModule(\"\\x22a string\\x22\\n//# sourceURL\\x3dURL\\n\");\n")
}

func TestClosureBundler_Traditional(t *testing.T) {
	sb := jutil.NewStringWriter()
	NewClosureBundler().AppendTo(sb, closureBundlerTest_TRADITIONAL, "\"a string\"")
	closureBundlerTest_assertEquals(t, sb.String(), "\"a string\"")
}

func TestClosureBundler_TraditionalWithSourceURL(t *testing.T) {
	sb := jutil.NewStringWriter()
	NewClosureBundler().WithSourceUrl("URL").AppendTo(sb, closureBundlerTest_TRADITIONAL, "\"a string\"")
	closureBundlerTest_assertEquals(t, sb.String(), "\"a string\"\n//# sourceURL=URL\n")
}

func TestClosureBundler_TraditionalWithEval(t *testing.T) {
	sb := jutil.NewStringWriter()
	NewClosureBundler().UseEval(true).AppendTo(sb, closureBundlerTest_TRADITIONAL, "\"a string\"")
	closureBundlerTest_assertEquals(t, sb.String(), "eval(this.CLOSURE_EVAL_PREFILTER(\"\\x22a string\\x22\"));\n")
}

func TestClosureBundler_TraditionalWithEvalWithSourceUrl(t *testing.T) {
	sb := jutil.NewStringWriter()
	NewClosureBundler().
		UseEval(true).
		WithSourceUrl("URL").
		AppendTo(sb, closureBundlerTest_TRADITIONAL, "\"a string\"")
	closureBundlerTest_assertEquals(t, sb.String(),
		"eval(this.CLOSURE_EVAL_PREFILTER(\"\\x22a string\\x22\\n//# sourceURL\\x3dURL\\n\"));\n")
}

// closureBundlerTest_mockTranspiler is the Mockito mock of testTranspilation
// (RETURNS_SMART_NULLS: an unstubbed transpile call returns null).
//
// Go-only: test-helper — Mockito.
type closureBundlerTest_mockTranspiler struct {
	path  string
	input string
}

func (m *closureBundlerTest_mockTranspiler) Runtime() string { return "RUNTIME;" }

func (m *closureBundlerTest_mockTranspiler) Transpile(path string, code string) *TranspileResult {
	if path == m.path && code == m.input {
		return NewTranspileResult(m.path, m.input, "TRANSPILED;", "")
	}
	return nil
}

func TestClosureBundler_Transpilation(t *testing.T) {
	input := "goog.module('Foo');\nclass Foo {}"

	var transpiler Transpiler = &closureBundlerTest_mockTranspiler{path: "foo.js", input: input}

	bundler := NewClosureBundlerWithTranspiler(transpiler).WithPath("foo.js")
	sb := jutil.NewStringWriter()
	bundler.AppendRuntimeTo(sb)
	bundler.AppendTo(sb, closureBundlerTest_MODULE, input)
	if !strings.HasPrefix(sb.String(), "RUNTIME;") {
		t.Errorf("%q does not start with RUNTIME;", sb.String())
	}
	// Call endsWith because the ES6 module runtime is also injected.
	want := "goog.loadModule(function(exports) {'use strict';TRANSPILED;\n;return exports;});\n"
	if !strings.HasSuffix(sb.String(), want) {
		t.Errorf("%q does not end with %q", sb.String(), want)
	}

	// Without calling appendRuntimeTo(), the runtime is not included anymore.
	sb = jutil.NewStringWriter()
	bundler.AppendTo(sb, closureBundlerTest_MODULE, input)
	closureBundlerTest_assertEquals(t, sb.String(), want)
}

type closureBundlerTestEs6Module struct{ jscomp.DependencyInfo }

func (closureBundlerTestEs6Module) IsEs6Module() bool { return true }

func TestClosureBundler_Es6Module(t *testing.T) {
	input := "import {x} from './other.js';\n" +
		"export {x as y};\n" +
		"let local;\n" +
		"export function foo() { return local; }\n"
	bundler := NewClosureBundlerWithTranspiler(BaseTranspiler_LATEST_TRANSPILER).WithPath("nested/path/foo.js")
	sb := jutil.NewStringWriter()
	bundler.AppendRuntimeTo(sb)
	bundler.AppendTo(
		sb,
		closureBundlerTestEs6Module{closureBundlerTest_TRADITIONAL},
		input)
	result := sb.String()
	// ES6 module runtime should be injected.
	if !strings.Contains(result, "$jscomp.require = createRequire();") {
		t.Errorf("no $jscomp.require in %q", result)
	}
	if !strings.HasPrefix(result, "var $jscomp") {
		t.Errorf("%q does not start with var $jscomp", result)
	}
	want := "$jscomp.registerAndLoadModule(function($$require, $$exports, $$module) {\n" +
		"  function foo() {\n" +
		"    return local;\n" +
		"  }\n" +
		"  Object.defineProperties($$exports, {foo:{enumerable:true, get:function() {\n" +
		"    return foo;\n" +
		"  }}, y:{enumerable:true, get:function() {\n" +
		"    return module$nested$path$other.x;\n" +
		"  }}});\n" +
		"  var module$nested$path$other = $$require(\"nested/path/other.js\");\n" +
		"  let local;\n" +
		"}, \"nested/path/foo.js\", [\"nested/path/other.js\"]);\n"
	if !strings.HasSuffix(result, want) {
		t.Errorf("%q does not end with %q", result, want)
	}
}

func TestClosureBundler_PassThroughIfNoTranspilationNeeded(t *testing.T) {
	input := "/** Hello Comments! */ const s = 0;\n  let intended;"
	bundler := NewClosureBundlerWithTranspiler(BaseTranspiler_LATEST_TRANSPILER)
	sb := jutil.NewStringWriter()
	bundler.AppendTo(sb, jscomp.SimpleDependencyInfo_Builder_("", "").Build(), input)
	closureBundlerTest_assertEquals(t, sb.String(), input)
}

func TestClosureBundler_CommentsAndFormattingRemovedWithTranspilation(t *testing.T) {
	input := "/** Hello Comments! */ const s = 0;\n  let intended;"
	bundler := NewClosureBundlerWithTranspiler(BaseTranspiler_ES5_TRANSPILER)
	sb := jutil.NewStringWriter()
	bundler.AppendTo(sb, jscomp.SimpleDependencyInfo_Builder_("", "").Build(), input)
	closureBundlerTest_assertEquals(t, sb.String(), "var s = 0;\nvar intended;\n")
}

func TestClosureBundler_FullWidthLowLineWithDefaultTranspilerIsOkay(t *testing.T) {
	// The last character is something the compiler doesn't handle correctly
	input := "var ａｅｓｔｈｅｔｉｃ＿"
	bundler := NewClosureBundler()
	sb := jutil.NewStringWriter()
	bundler.AppendTo(sb, jscomp.SimpleDependencyInfo_Builder_("", "").Build(), input)
	closureBundlerTest_assertEquals(t, sb.String(), input)
}

// TODO(johnplaisted): If / when the compiler can parse full width low line in
// identifiers this should be okay to be transpiled.
func TestClosureBundler_FullWidthLowLineInTranspiledCodeIsError(t *testing.T) {
	// The last character is something the compiler doesn't handle correctly
	input := "let ａｅｓｔｈｅｔｉｃ＿"
	bundler := NewClosureBundlerWithTranspiler(BaseTranspiler_ES5_TRANSPILER)
	sb := jutil.NewStringWriter()
	thrown := false
	jutil.Try(func() {
		bundler.AppendTo(sb, jscomp.SimpleDependencyInfo_Builder_("", "").Build(), input)
	}, func(e *bundle.TranspilationException) {
		thrown = true
		if !strings.Contains(e.GetMessage(), "Parse error. Character '＿' (U+FF3F) is not a valid identifier start char") {
			t.Errorf("message %q", e.GetMessage())
		}
	})
	if !thrown {
		t.Errorf("Expected an exception")
	}
}
