// Port of test/com/google/javascript/jscomp/deps/JsFileFullParserTest.java @ closure-compiler v20261021

/*
 * Copyright 2019 The Closure Compiler Authors.
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
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// jsFileFullParserTest_failReporter is the lambda of parse(String): it fails
// (throws an AssertionError) on any report.
type jsFileFullParserTest_failReporter struct{}

func (jsFileFullParserTest_failReporter) Report(fatal bool, message string, sourceName string, line int, lineOffset int) {
	panic(jutil.NewAssertionError(fmt.Sprintf("%s:%d:%d: %s", sourceName, line, lineOffset, message)))
}

// Java: JsFileFullParserTest.parse(String)
func jsFileFullParserTest_parse(content string) *JsFileFullParser_FileInfo {
	return JsFileFullParser_Parse(
		content,
		"file.js",
		jsFileFullParserTest_failReporter{})
}

func jsFileFullParserTest_assertModuleType(t *testing.T, info *JsFileFullParser_FileInfo, want JsFileFullParser_FileInfo_ModuleType) {
	t.Helper()
	if info.moduleType != want {
		t.Errorf("moduleType = %v, want %v", info.moduleType, want)
	}
}

// jsFileFullParserTest_containsExactly is Truth's containsExactly (in any
// order, with multiplicities).
func jsFileFullParserTest_containsExactly(t *testing.T, name string, got []string, want ...string) {
	t.Helper()
	g, w := slices.Clone(got), slices.Clone(want)
	slices.Sort(g)
	slices.Sort(w)
	if !slices.Equal(g, w) {
		t.Errorf("%s = %q, want exactly %q", name, got, want)
	}
}

func TestJsFileFullParser_ModuleType_script(t *testing.T) {
	info := jsFileFullParserTest_parse("alert(42);")
	jsFileFullParserTest_assertModuleType(t, info, JsFileFullParser_FileInfo_ModuleType_UNKNOWN)
}

func TestJsFileFullParser_ModuleType_cjs(t *testing.T) {
	info := jsFileFullParserTest_parse("const fs = require('fs');")
	jsFileFullParserTest_assertModuleType(t, info, JsFileFullParser_FileInfo_ModuleType_UNKNOWN)
}

func TestJsFileFullParser_ModuleType_googProvide(t *testing.T) {
	info := jsFileFullParserTest_parse("goog.provide('foo');")
	jsFileFullParserTest_assertModuleType(t, info, JsFileFullParser_FileInfo_ModuleType_GOOG_PROVIDE)
}

func TestJsFileFullParser_ModuleType_googModule(t *testing.T) {
	info := jsFileFullParserTest_parse("goog.module('foo');")
	jsFileFullParserTest_assertModuleType(t, info, JsFileFullParser_FileInfo_ModuleType_GOOG_MODULE)
}

func TestJsFileFullParser_ModuleType_legacyGoogModule(t *testing.T) {
	info := jsFileFullParserTest_parse(
		"goog.module('foo');\n" +
			"goog.module.declareLegacyNamespace();\n")
	jsFileFullParserTest_assertModuleType(t, info, JsFileFullParser_FileInfo_ModuleType_GOOG_MODULE)
}

func TestJsFileFullParser_ModuleType_esModule_export(t *testing.T) {
	info := jsFileFullParserTest_parse("export function f() {};")
	jsFileFullParserTest_assertModuleType(t, info, JsFileFullParser_FileInfo_ModuleType_ES_MODULE)
}

func TestJsFileFullParser_ModuleType_esModule_import(t *testing.T) {
	info := jsFileFullParserTest_parse("import * as x from './bar';")
	jsFileFullParserTest_assertModuleType(t, info, JsFileFullParser_FileInfo_ModuleType_ES_MODULE)
}

func TestJsFileFullParser_ProvidesRequires(t *testing.T) {
	info := jsFileFullParserTest_parse(
		"goog.provide('providedSymbol');\n" +
			"goog.require('stronglyRequiredSymbol');\n" +
			"goog.requireType('weaklyRequiredSymbol');\n" +
			"async function test() {await goog.requireDynamic('dynamicallyRequiredSymbol');}\n")

	jsFileFullParserTest_containsExactly(t, "provides", info.provides, "providedSymbol")
	jsFileFullParserTest_containsExactly(t, "requires", info.requires, "stronglyRequiredSymbol")
	jsFileFullParserTest_containsExactly(t, "typeRequires", info.typeRequires, "weaklyRequiredSymbol")
	jsFileFullParserTest_containsExactly(t, "dynamicRequires", info.dynamicRequires.ToSlice(), "dynamicallyRequiredSymbol")
}

// jsFileFullParserTest_ErrorCollector is the local class ErrorCollector of
// testProvidesRequires_syntaxError.
type jsFileFullParserTest_ErrorCollector struct {
	lastErrorMessage *string
}

func (c *jsFileFullParserTest_ErrorCollector) Report(fatal bool, message string, sourceName string, line int, lineOffset int) {
	c.lastErrorMessage = &message
}

func TestJsFileFullParser_ProvidesRequires_syntaxError(t *testing.T) {
	src := "goog.provide('providedSymbol');\n" +
		"goog.require('stronglyRequiredSymbol');\n" +
		"goog.requireType('weaklyRequiredSymbol');\n" +
		"\n" +
		"syntax error;\n"

	errorCollector := &jsFileFullParserTest_ErrorCollector{}
	info := JsFileFullParser_Parse(src, "file.js", errorCollector)

	jsFileFullParserTest_containsExactly(t, "provides", info.provides)
	jsFileFullParserTest_containsExactly(t, "requires", info.requires)
	jsFileFullParserTest_containsExactly(t, "typeRequires", info.typeRequires)
	if errorCollector.lastErrorMessage == nil || *errorCollector.lastErrorMessage != "Semi-colon expected" {
		t.Errorf("lastErrorMessage = %v", errorCollector.lastErrorMessage)
	}
}

func TestJsFileFullParser_GoogLoadModule(t *testing.T) {
	info := jsFileFullParserTest_parse(
		"goog.loadModule((exports) => {\n" +
			"  goog.module('inner');\n" +
			"  const i = goog.require('inner_require');\n" +
			"  exports = {inner: 1};\n" +
			"  return exports;\n" +
			"});\n" +
			"goog.loadModule((exports) => {\n" +
			"  goog.module('inner2');\n" +
			"  // intentional self-edge\n" +
			"  const i = goog.require('inner');\n" +
			"  exports = {inner2: 1};\n" +
			"  return exports;\n" +
			"});\n")

	jsFileFullParserTest_containsExactly(t, "provides", info.provides, "inner", "inner2")
	jsFileFullParserTest_containsExactly(t, "requires", info.requires, "inner", "inner_require")
}

func TestJsFileFullParser_GoogLoadModule_inModule(t *testing.T) {
	info := jsFileFullParserTest_parse(
		"goog.module('outer');\n" +
			"const i = goog.require('outer_require');\n" +
			"goog.loadModule((exports) => {\n" +
			"  goog.module('inner');\n" +
			"  const i = goog.require('inner_require');\n" +
			"  // intentional self-edge\n" +
			"  const i2 = goog.require('outer');\n" +
			"  exports = {inner: 1};\n" +
			"  return exports;\n" +
			"});\n" +
			"goog.loadModule((exports) => {\n" +
			"  goog.module('inner2');\n" +
			"  // intentional self-edge\n" +
			"  const i = goog.require('inner');\n" +
			"  exports = {inner2: 1};\n" +
			"  return exports;\n" +
			"});\n" +
			"exports = {outer: 1};\n")

	jsFileFullParserTest_containsExactly(t, "provides", info.provides, "outer", "inner", "inner2")
	jsFileFullParserTest_containsExactly(t, "requires", info.requires, "inner", "inner_require", "outer", "outer_require")
}

func TestJsFileFullParser_GoogLoadModule_doubleNested(t *testing.T) {
	// NB: we cannot use assertThrows: the bazel CI setup doe not have JUnit 5.
	caught := false
	jutil.Try(func() {
		jsFileFullParserTest_parse(
			"goog.module('outer');\n" +
				"const i = goog.require('outer_require');\n" +
				"goog.loadModule((exports) => {\n" +
				"  goog.module('inner');\n" +
				"  goog.loadModule((exports) => {\n" +
				"    goog.module('inner.inner');\n" +
				"    exports = {inner2: 1};\n" +
				"    return exports;\n" +
				"  });\n" +
				"  const i = goog.require('inner_require');\n" +
				"  // intentional self-edge\n" +
				"  const i2 = goog.require('outer');\n" +
				"  exports = {inner: 1};\n" +
				"  return exports;\n" +
				"});\n" +
				"exports = {outer: 1};\n")
		panic(jutil.NewRuntimeException("expected an AssertionError"))
	}, func(e *jutil.AssertionError) {
		caught = true
		if !strings.Contains(e.GetMessage(), "goog.loadModule cannot be nested") {
			t.Errorf("message %q", e.GetMessage())
		}
	})
	if !caught {
		t.Errorf("expected an AssertionError")
	}
}

func TestJsFileFullParser_SoyDeltemplate(t *testing.T) {
	info := jsFileFullParserTest_parse(
		"/**\n" +
			" * @deltemplate {a}\n" +
			" * @modName {m}\n" +
			" */\n")
	jsFileFullParserTest_containsExactly(t, "deltemplates", info.deltemplates.ToSlice(), "a")
}

func TestJsFileFullParser_SoyDeltemplate_legacy(t *testing.T) {
	info := jsFileFullParserTest_parse(
		"/**\n" +
			" * @hassoydeltemplate {a}\n" +
			" * @modName {m}\n" +
			" */\n")
	jsFileFullParserTest_containsExactly(t, "deltemplates", info.deltemplates.ToSlice(), "a")
}

func TestJsFileFullParser_SoyDelcall(t *testing.T) {
	info := jsFileFullParserTest_parse(
		"/**\n" +
			" * @delcall {a}\n" +
			" */\n")
	jsFileFullParserTest_containsExactly(t, "delcalls", info.delcalls.ToSlice(), "a")
}

func TestJsFileFullParser_SoyDelcall_legacy(t *testing.T) {
	info := jsFileFullParserTest_parse(
		"/**\n" +
			" * @hassoydelcall {a}\n" +
			" */\n")
	jsFileFullParserTest_containsExactly(t, "delcalls", info.delcalls.ToSlice(), "a")
}

func TestJsFileFullParser_ReadToggle(t *testing.T) {
	info := jsFileFullParserTest_parse("const {TOGGLE_foo_bar} = goog.require('foo$2etoggles');")
	jsFileFullParserTest_containsExactly(t, "readToggles", info.readToggles.ToSlice(), "foo_bar")
}

func TestJsFileFullParser_TopLevelJsDocComments_beforeBetweenAndAfterStatements(t *testing.T) {
	info := jsFileFullParserTest_parse(
		"/** @fileoverview Header comment. */\n" +
			"/** @requirecss {css.before.module} */\n" +
			"goog.module('foo');\n" +
			"/** @requirecss {css.between.statements} */\n" +
			"class Foo {}\n" +
			"/** @requirecss {css.after.statements} */\n")
	jsFileFullParserTest_containsExactly(t, "requiresCss", info.requiresCss,
		"css.before.module", "css.between.statements", "css.after.statements")
}

func TestJsFileFullParser_NestedJsDocCommentsAreIgnored(t *testing.T) {
	info := jsFileFullParserTest_parse(
		"goog.module('foo');\n" +
			"class Foo {\n" +
			"  /** @requirecss {css.in.method} */\n" +
			"  bar() {}\n" +
			"}\n" +
			"function baz() {\n" +
			"  /** @requirecss {css.in.function} */\n" +
			"  return 1;\n" +
			"}\n" +
			"goog.scope(function() {\n" +
			"  /** @requirecss {css.in.goog_scope} */\n" +
			"  const x = 1;\n" +
			"});\n" +
			"{\n" +
			"  /** @requirecss {css.in.block} */\n" +
			"}\n")
	jsFileFullParserTest_containsExactly(t, "requiresCss", info.requiresCss)
}

func TestJsFileFullParser_TopLevelJsDocComments_moreLinesThanNodeCanRepresent(t *testing.T) {
	// Node encodes line numbers in 20 bits, so a statement on line 2^20 wraps to line 0.
	code := "/** @requirecss {css.before} */\n" +
		"goog.module('foo');\n" +
		strings.Repeat("\n", (1<<20)-3) +
		"const x = 1;\n" +
		"function f() {\n" +
		"  /** @requirecss {css.nested} */\n" +
		"}\n" +
		"/** @requirecss {css.after} */\n"
	info := jsFileFullParserTest_parse(code)
	jsFileFullParserTest_containsExactly(t, "requiresCss", info.requiresCss, "css.before", "css.after")
}

func TestJsFileFullParser_TopLevelJsDocComments_loneCarriageReturns(t *testing.T) {
	// The parser counts each lone \r as a line break, but rhino's SourceFile doesn't.
	info := jsFileFullParserTest_parse(
		"a();\r\r\rfoo();\n" +
			"function f() {\n" +
			"  /** @requirecss {css.nested} */\n" +
			"}\n" +
			"/** @requirecss {css.after} */\n")
	jsFileFullParserTest_containsExactly(t, "requiresCss", info.requiresCss, "css.after")
}

func TestJsFileFullParser_TopLevelJsDocComments_unicodeLineSeparators(t *testing.T) {
	// The parser counts U+2028 and U+2029 as line breaks, even inside strings, but rhino's
	// SourceFile doesn't.
	info := jsFileFullParserTest_parse(
		"a();\n'  ';\nfoo();\nfunction f() {\n/** @requirecss {css.nested} */\n}\n" +
			"/** @requirecss {css.after} */")
	jsFileFullParserTest_containsExactly(t, "requiresCss", info.requiresCss, "css.after")
}

func TestJsFileFullParser_TopLevelJsDocComments_statementPastMaxColumn(t *testing.T) {
	// Node clamps column numbers at 4095, so Node-based offsets for `var b` and `f` are wrong.
	code := "var a = 1;" +
		strings.Repeat(" ", 4100) +
		"/** @requirecss {css.before.b} */ var b = '" +
		strings.Repeat("x", 5000) +
		"'; function f() { /** @requirecss {css.nested} */ }"
	info := jsFileFullParserTest_parse(code)
	jsFileFullParserTest_containsExactly(t, "requiresCss", info.requiresCss, "css.before.b")
}
