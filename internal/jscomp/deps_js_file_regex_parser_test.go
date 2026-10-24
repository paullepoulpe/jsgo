// Port of test/com/google/javascript/jscomp/deps/JsFileRegexParserTest.java @ closure-compiler v20261021

/*
 * Copyright 2008 The Closure Compiler Authors.
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
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

const (
	jsFileRegexParserTest_SRC_PATH     = "a"
	jsFileRegexParserTest_CLOSURE_PATH = "b"
)

type jsFileRegexParserTest struct {
	t            *testing.T
	parser       *JsFileRegexParser
	errorManager ErrorManager
}

func newJsFileRegexParserTest(t *testing.T) *jsFileRegexParserTest {
	errorManager := NewPrintStreamErrorManager(jutil.System_Err())
	parser := NewJsFileRegexParser(errorManager)
	parser.SetShortcutMode(true)
	return &jsFileRegexParserTest{t: t, parser: parser, errorManager: errorManager}
}

func jsFileRegexParserTest_lines(lines ...string) string {
	s := ""
	for _, l := range lines {
		s += l + "\n"
	}
	return s
}

func jsFileRegexParserTest_builder() *SimpleDependencyInfo_Builder {
	return SimpleDependencyInfo_Builder_(jsFileRegexParserTest_CLOSURE_PATH, jsFileRegexParserTest_SRC_PATH)
}

func jsFileRegexParserTest_goog(s string) *DependencyInfo_Require {
	return DependencyInfo_Require_GoogRequireSymbol(s)
}

func jsFileRegexParserTest_es6(symbol, raw string) *DependencyInfo_Require {
	return DependencyInfo_Require_Es6Import(symbol, raw)
}

func jsFileRegexParserTest_flags(kv ...string) *jutil.LinkedHashMap[string, string] {
	m := jutil.NewLinkedHashMap[string, string]()
	for i := 0; i < len(kv); i += 2 {
		m.Put(kv[i], kv[i+1])
	}
	return m
}

// assertDeps asserts the deps match without errors
func (x *jsFileRegexParserTest) assertDeps(expected DependencyInfo, actual DependencyInfo) {
	x.t.Helper()
	if !jutil.ObjectsEquals(actual, expected) {
		x.t.Errorf("deps:\n got %s\nwant %s", jutil.ToJavaString(actual), jutil.ToJavaString(expected))
	}
	if len(x.errorManager.GetErrors()) != 0 {
		x.t.Errorf("errors: %v", x.errorManager.GetErrors())
	}
	if len(x.errorManager.GetWarnings()) != 0 {
		x.t.Errorf("warnings: %v", x.errorManager.GetWarnings())
	}
}

// Tests: parsing of comments, parsing of different styles of quotes,
// correct recording of what was parsed.
func TestJsFileRegexParser_ParseFile(t *testing.T) {
	x := newJsFileRegexParserTest(t)
	contents := jsFileRegexParserTest_lines(
		"/*goog.provide('no1');*//*",
		"goog.provide('no2');",
		"*/goog.provide('yes1');",
		"/* blah */goog.provide(\"yes2\")/* blah*/",
		"goog.require('yes3'); // goog.provide('no3');",
		"// goog.provide('no4');",
		"goog.require(\"bar.data.SuperstarAddStarThreadActionRequestDelegate\"); //no new line at EOF")

	expected := jsFileRegexParserTest_builder().
		SetProvides([]string{"yes1", "yes2"}).
		SetRequiresArray(
			jsFileRegexParserTest_goog("yes3"),
			jsFileRegexParserTest_goog("bar.data.SuperstarAddStarThreadActionRequestDelegate")).
		Build()

	result := x.parser.ParseFile(jsFileRegexParserTest_SRC_PATH, jsFileRegexParserTest_CLOSURE_PATH, contents)

	x.assertDeps(expected, result)
}

// Tests correct recording of what was parsed.
func TestJsFileRegexParser_ParseFile2(t *testing.T) {
	x := newJsFileRegexParserTest(t)
	contents := jsFileRegexParserTest_lines(
		"goog.module('yes1');",
		"var yes2 = goog.require('yes2');",
		"var C = goog.require(\"a.b.C\");",
		"let {D, E} = goog.require('a.b.d');")

	expected := jsFileRegexParserTest_builder().
		SetProvides([]string{"yes1"}).
		SetRequires([]*DependencyInfo_Require{jsFileRegexParserTest_goog("yes2"), jsFileRegexParserTest_goog("a.b.C"), jsFileRegexParserTest_goog("a.b.d")}).
		SetLoadFlags(jsFileRegexParserTest_flags("module", "goog")).
		Build()

	result := x.parser.ParseFile(jsFileRegexParserTest_SRC_PATH, jsFileRegexParserTest_CLOSURE_PATH, contents)

	x.assertDeps(expected, result)
}

// Tests correct recording of what was parsed.
func TestJsFileRegexParser_ParseFile3(t *testing.T) {
	x := newJsFileRegexParserTest(t)
	contents := jsFileRegexParserTest_lines(
		"goog.module('yes1');",
		"var yes2=goog.require('yes2');",
		"var C=goog.require(\"a.b.C\");",
		"const {",
		"  D,",
		"  E",
		"}=goog.require(\"a.b.d\");")

	expected := jsFileRegexParserTest_builder().
		SetProvides([]string{"yes1"}).
		SetRequires([]*DependencyInfo_Require{jsFileRegexParserTest_goog("yes2"), jsFileRegexParserTest_goog("a.b.C"), jsFileRegexParserTest_goog("a.b.d")}).
		SetLoadFlags(jsFileRegexParserTest_flags("module", "goog")).
		Build()

	result := x.parser.ParseFile(jsFileRegexParserTest_SRC_PATH, jsFileRegexParserTest_CLOSURE_PATH, contents)

	x.assertDeps(expected, result)
}

// Tests correct recording of what was parsed.
func TestJsFileRegexParser_ParseFileWithMultiLineRequires(t *testing.T) {
	x := newJsFileRegexParserTest(t)
	x.parser.SetShortcutMode(false)

	contents := jsFileRegexParserTest_lines(
		"goog.module('yes1');",
		"var fakerequire = 5;",
		"var yes2=goog.require(",
		"'yes2');",
		"var C=",
		"goog.require(\"a.b.C\");",
		"const {",
		"  D,",
		"  E",
		"}=goog.require(",
		"\"a.b.d\");")

	expected := jsFileRegexParserTest_builder().
		SetProvides([]string{"yes1"}).
		SetRequires([]*DependencyInfo_Require{jsFileRegexParserTest_goog("yes2"), jsFileRegexParserTest_goog("a.b.C"), jsFileRegexParserTest_goog("a.b.d")}).
		SetLoadFlags(jsFileRegexParserTest_flags("module", "goog")).
		Build()

	result := x.parser.ParseFile(jsFileRegexParserTest_SRC_PATH, jsFileRegexParserTest_CLOSURE_PATH, contents)

	x.assertDeps(expected, result)
}

func TestJsFileRegexParser_UnicodeCharacters(t *testing.T) {
	x := newJsFileRegexParserTest(t)
	contents := jsFileRegexParserTest_lines(
		"const FɵɵΔBar = goog.require('imaginary.test.namespace');")

	expected := jsFileRegexParserTest_builder().
		SetRequires([]*DependencyInfo_Require{jsFileRegexParserTest_goog("imaginary.test.namespace")}).
		Build()

	result := x.parser.ParseFile(jsFileRegexParserTest_SRC_PATH, jsFileRegexParserTest_CLOSURE_PATH, contents)

	x.assertDeps(expected, result)
}

func TestJsFileRegexParser_ParseGoogModuleWithRequireType(t *testing.T) {
	x := newJsFileRegexParserTest(t)
	contents := jsFileRegexParserTest_lines(
		"goog.module('yes1');",
		"var yes2=goog.requireType('yes2');",
		"var C=goog.requireType(\"a.b.C\");",
		"const {",
		"  D,",
		"  E",
		"}=goog.requireType(\"a.b.d\");")

	expected := jsFileRegexParserTest_builder().
		SetProvides([]string{"yes1"}).
		SetTypeRequires([]string{"yes2", "a.b.C", "a.b.d"}).
		SetLoadFlags(jsFileRegexParserTest_flags("module", "goog")).
		Build()

	result := x.parser.ParseFile(jsFileRegexParserTest_SRC_PATH, jsFileRegexParserTest_CLOSURE_PATH, contents)

	x.assertDeps(expected, result)
}

func TestJsFileRegexParser_ParseScriptWithRequireType(t *testing.T) {
	x := newJsFileRegexParserTest(t)
	contents := jsFileRegexParserTest_lines(
		"goog.provide('yes1');",
		"goog.requireType('a.b.C');")

	expected := jsFileRegexParserTest_builder().
		SetProvides([]string{"yes1"}).
		SetTypeRequires([]string{"a.b.C"}).
		Build()

	result := x.parser.ParseFile(jsFileRegexParserTest_SRC_PATH, jsFileRegexParserTest_CLOSURE_PATH, contents)

	x.assertDeps(expected, result)
}

// Tests correct recording of what was parsed.
func TestJsFileRegexParser_ParseWrappedGoogModule(t *testing.T) {
	x := newJsFileRegexParserTest(t)
	contents := jsFileRegexParserTest_lines(
		"goog.loadModule(function(){\"use strict\";goog.module('yes1');",
		"var yes2=goog.require('yes2');",
		"var C=goog.require(\"a.b.C\");",
		"const {",
		"  D,",
		"  E",
		"}=goog.require(\"a.b.d\");});")

	expected := jsFileRegexParserTest_builder().
		SetProvides([]string{"yes1"}).
		SetRequires([]*DependencyInfo_Require{jsFileRegexParserTest_goog("yes2"), jsFileRegexParserTest_goog("a.b.C"), jsFileRegexParserTest_goog("a.b.d")}).
		SetLoadFlags(jsFileRegexParserTest_flags()).
		Build() // wrapped modules aren't marked as modules

	result := x.parser.ParseFile(jsFileRegexParserTest_SRC_PATH, jsFileRegexParserTest_CLOSURE_PATH, contents)

	x.assertDeps(expected, result)
}

// TODO(sdh): Add a test for import with .js suffix once #1897 is fixed.

// Tests ES6 modules parsed correctly, particularly the various formats.
func TestJsFileRegexParser_ParseEs6Module(t *testing.T) {
	x := newJsFileRegexParserTest(t)
	contents := jsFileRegexParserTest_lines(
		"import def, {yes2} from './yes2';",
		"import C from './a/b/C';",
		"import * as d from './a/b/d';",
		"import \"./dquote\";",
		"export * from './exported';")

	expected := SimpleDependencyInfo_Builder_("a.js", "b.js").
		SetProvides([]string{"module$b"}).
		SetRequiresArray(
			jsFileRegexParserTest_es6("module$yes2", "./yes2"),
			jsFileRegexParserTest_es6("module$a$b$C", "./a/b/C"),
			jsFileRegexParserTest_es6("module$a$b$d", "./a/b/d"),
			jsFileRegexParserTest_es6("module$dquote", "./dquote"),
			jsFileRegexParserTest_es6("module$exported", "./exported")).
		SetLoadFlags(jsFileRegexParserTest_flags("module", "es6")).
		Build()

	result := x.parser.ParseFile("b.js", "a.js", contents)

	x.assertDeps(expected, result)
}

// Tests relative paths resolved correctly.
func TestJsFileRegexParser_ParseEs6Module2(t *testing.T) {
	x := newJsFileRegexParserTest(t)
	contents := jsFileRegexParserTest_lines(
		"import './x';",
		"import '../y';",
		"import '../a/z';",
		"import '../c/w';")

	expected := SimpleDependencyInfo_Builder_("../../a/b.js", "/foo/bar/a/b.js").
		SetProvides([]string{"module$foo$bar$a$b"}).
		SetRequires([]*DependencyInfo_Require{
			jsFileRegexParserTest_es6("module$foo$bar$a$x", "./x"), jsFileRegexParserTest_es6("module$foo$bar$y", "../y"),
			jsFileRegexParserTest_es6("module$foo$bar$a$z", "../a/z"),
			jsFileRegexParserTest_es6("module$foo$bar$c$w", "../c/w")}).
		SetLoadFlags(jsFileRegexParserTest_flags("module", "es6")).
		Build()

	result := x.parser.ParseFile("/foo/bar/a/b.js", "../../a/b.js", contents)

	x.assertDeps(expected, result)
}

// Tests handles goog.require and import 'goog:...'.
func TestJsFileRegexParser_ParseEs6Module3(t *testing.T) {
	x := newJsFileRegexParserTest(t)
	contents := jsFileRegexParserTest_lines(
		"import 'goog:foo.bar.baz';",
		"goog.require('baz.qux');")

	expected := SimpleDependencyInfo_Builder_("b.js", "a.js").
		SetProvides([]string{"module$a"}).
		SetRequires([]*DependencyInfo_Require{jsFileRegexParserTest_goog("foo.bar.baz"), jsFileRegexParserTest_goog("baz.qux")}).
		SetLoadFlags(jsFileRegexParserTest_flags("module", "es6")).
		Build()

	result := x.parser.ParseFile("a.js", "b.js", contents)

	x.assertDeps(expected, result)
}

func jsFileRegexParserTest_fooLoader() *ModuleLoader {
	return ModuleLoader_Builder_().
		SetModuleRoots([]string{"/foo"}).
		SetInputs([]DependencyInfo{}).
		SetFactory(BrowserModuleResolver_FACTORY).
		Build()
}

// Tests setModuleLoader taken into account
func TestJsFileRegexParser_ParseEs6Module4(t *testing.T) {
	x := newJsFileRegexParserTest(t)
	loader := jsFileRegexParserTest_fooLoader()

	contents := jsFileRegexParserTest_lines(
		"import './a';",
		"import './qux/b';",
		"import '../closure/c';",
		"import '../closure/d/e';",
		"import '../../corge/f';")
	expected := SimpleDependencyInfo_Builder_("../bar/baz.js", "/foo/js/bar/baz.js").
		SetProvides([]string{"module$js$bar$baz"}).
		SetRequires([]*DependencyInfo_Require{
			jsFileRegexParserTest_es6("module$js$bar$a", "./a"),
			jsFileRegexParserTest_es6("module$js$bar$qux$b", "./qux/b"),
			jsFileRegexParserTest_es6("module$js$closure$c", "../closure/c"),
			jsFileRegexParserTest_es6("module$js$closure$d$e", "../closure/d/e"),
			jsFileRegexParserTest_es6("module$corge$f", "../../corge/f")}).
		SetLoadFlags(jsFileRegexParserTest_flags("module", "es6")).
		Build()

	result := x.parser.SetModuleLoader(loader).ParseFile("/foo/js/bar/baz.js", "../bar/baz.js", contents)

	x.assertDeps(expected, result)
}

// TODO(johnplaisted): This should eventually be an error. For now people are relying on this
// behavior for interop / ordering. Until we have official channels for these allow this behavior,
// but don't encourage it.
func TestJsFileRegexParser_ParseEs6ModuleWithGoogProvide(t *testing.T) {
	x := newJsFileRegexParserTest(t)
	loader := jsFileRegexParserTest_fooLoader()

	contents := "goog.provide('my.namespace');\nexport {};"

	expected := SimpleDependencyInfo_Builder_("../bar/baz.js", "/foo/js/bar/baz.js").
		SetProvides([]string{"my.namespace"}).
		Build()

	result := x.parser.SetModuleLoader(loader).ParseFile("/foo/js/bar/baz.js", "../bar/baz.js", contents)

	if !jutil.ObjectsEquals(result, expected) {
		t.Errorf("got %s, want %s", jutil.ToJavaString(result), jutil.ToJavaString(expected))
	}
	assertEqual(t, "errors", x.errorManager.GetErrorCount(), 0)
	assertEqual(t, "warnings", x.errorManager.GetWarningCount(), 1)
	if x.errorManager.GetWarningCount() == 1 {
		assertEqual(t, "warning type", x.errorManager.GetWarnings()[0].Type(), ModuleLoader_MODULE_CONFLICT)
	}
}

func TestJsFileRegexParser_Es6ModuleWithDeclareModuleId(t *testing.T) {
	x := newJsFileRegexParserTest(t)
	loader := jsFileRegexParserTest_fooLoader()

	contents := "goog.declareModuleId('my.namespace');\nexport {};"

	expected := SimpleDependencyInfo_Builder_("../bar/baz.js", "/foo/js/bar/baz.js").
		SetProvides([]string{"my.namespace", "module$js$bar$baz"}).
		SetLoadFlags(jsFileRegexParserTest_flags("module", "es6")).
		Build()

	result := x.parser.SetModuleLoader(loader).ParseFile("/foo/js/bar/baz.js", "../bar/baz.js", contents)

	x.assertDeps(expected, result)
}

func TestJsFileRegexParser_Es6ModuleWithBrowserTransformedPrefixResolver(t *testing.T) {
	x := newJsFileRegexParserTest(t)
	loader := ModuleLoader_Builder_().
		SetModuleRoots([]string{}).
		SetInputs([]DependencyInfo{}).
		SetFactory(
			NewBrowserWithTransformedPrefixesModuleResolver_Factory(
				jsFileRegexParserTest_flags("@root/", "/path/to/project/"))).
		Build()

	contents := "import '@root/my/file.js';"

	expected := SimpleDependencyInfo_Builder_("../bar/baz.js", "/foo/js/bar/baz.js").
		SetProvides([]string{"module$foo$js$bar$baz"}).
		SetRequiresArray(jsFileRegexParserTest_es6("module$path$to$project$my$file", "@root/my/file.js")).
		SetLoadFlags(jsFileRegexParserTest_flags("module", "es6")).
		Build()

	result := x.parser.SetModuleLoader(loader).ParseFile("/foo/js/bar/baz.js", "../bar/baz.js", contents)

	x.assertDeps(expected, result)
}

// Tests shortcut mode doesn't stop at setTestOnly() or declareLegacyNamespace().
func TestJsFileRegexParser_NoShortcutForCommonModuleModifiers(t *testing.T) {
	x := newJsFileRegexParserTest(t)
	contents := jsFileRegexParserTest_lines(
		"goog.module('yes1');",
		"goog.module.declareLegacyNamespace();",
		"goog.setTestOnly();",
		"var yes2=goog.require('yes2');",
		"var C=goog.require(\"a.b.C\");",
		"const {",
		"  D,",
		"  E",
		"}=goog.require(\"a.b.d\");")

	expected := jsFileRegexParserTest_builder().
		SetProvides([]string{"yes1"}).
		SetRequiresArray(jsFileRegexParserTest_goog("yes2"), jsFileRegexParserTest_goog("a.b.C"), jsFileRegexParserTest_goog("a.b.d")).
		SetGoogModule(true).
		Build()

	result := x.parser.ParseFile(jsFileRegexParserTest_SRC_PATH, jsFileRegexParserTest_CLOSURE_PATH, contents)

	x.assertDeps(expected, result)
}

func TestJsFileRegexParser_MultiplePerLine(t *testing.T) {
	x := newJsFileRegexParserTest(t)
	contents := jsFileRegexParserTest_lines(
		"goog.provide('yes1');goog.provide('yes2');/*goog.provide('no1');*/goog.provide('yes3');//goog.provide('no2');")

	expected := jsFileRegexParserTest_builder().
		SetProvides([]string{"yes1", "yes2", "yes3"}).
		SetGoogModule(false).
		Build()

	result := x.parser.ParseFile(jsFileRegexParserTest_SRC_PATH, jsFileRegexParserTest_CLOSURE_PATH, contents)

	x.assertDeps(expected, result)
}

func TestJsFileRegexParser_ShortcutMode1(t *testing.T) {
	x := newJsFileRegexParserTest(t)
	// For efficiency reasons, we stop reading after the ctor.
	contents := jsFileRegexParserTest_lines(
		"// hi !",
		"/* this is a comment */",
		"goog.provide('yes1');",
		"/* and another comment */",
		"goog.provide('yes2'); // include this",
		"foo = function() {};",
		"goog.provide('no1');")

	expected := jsFileRegexParserTest_builder().
		SetProvides([]string{"yes1", "yes2"}).
		SetGoogModule(false).
		Build()
	result := x.parser.ParseFile(jsFileRegexParserTest_SRC_PATH, jsFileRegexParserTest_CLOSURE_PATH, contents)

	x.assertDeps(expected, result)
}

func TestJsFileRegexParser_ShortcutMode2(t *testing.T) {
	x := newJsFileRegexParserTest(t)
	contents := jsFileRegexParserTest_lines(
		"/** goog.provide('no1'); ",
		" * goog.provide('no2');",
		" */",
		"goog.provide('yes1');")

	expected := jsFileRegexParserTest_builder().
		SetProvides([]string{"yes1"}).
		SetGoogModule(false).
		Build()
	result := x.parser.ParseFile(jsFileRegexParserTest_SRC_PATH, jsFileRegexParserTest_CLOSURE_PATH, contents)

	x.assertDeps(expected, result)
}

func TestJsFileRegexParser_ShortcutMode3(t *testing.T) {
	x := newJsFileRegexParserTest(t)
	contents := jsFileRegexParserTest_lines(
		"/**",
		" * goog.provide('no1');",
		" */",
		"goog.provide('yes1');")

	expected := jsFileRegexParserTest_builder().
		SetProvides([]string{"yes1"}).
		SetGoogModule(false).
		Build()
	result := x.parser.ParseFile(jsFileRegexParserTest_SRC_PATH, jsFileRegexParserTest_CLOSURE_PATH, contents)

	x.assertDeps(expected, result)
}

func TestJsFileRegexParser_IncludeGoog1(t *testing.T) {
	x := newJsFileRegexParserTest(t)
	contents := jsFileRegexParserTest_lines(
		"/**",
		" * @provideGoog",
		" */")

	expected := jsFileRegexParserTest_builder().
		SetProvides([]string{"goog"}).
		SetGoogModule(false).
		Build()
	result := x.parser.SetIncludeGoogBase(true).ParseFile(jsFileRegexParserTest_SRC_PATH, jsFileRegexParserTest_CLOSURE_PATH, contents)
	x.assertDeps(expected, result)
}

func TestJsFileRegexParser_IncludeGoog1_quotes(t *testing.T) {
	x := newJsFileRegexParserTest(t)
	contents := "var x = \"/**\\\\n * @provideGoog\\\\n */\\\\n\";"

	expected := jsFileRegexParserTest_builder().
		SetProvides([]string{}).
		SetGoogModule(false).
		Build()
	result := x.parser.SetIncludeGoogBase(true).ParseFile(jsFileRegexParserTest_SRC_PATH, jsFileRegexParserTest_CLOSURE_PATH, contents)
	x.assertDeps(expected, result)
}

func TestJsFileRegexParser_IncludeGoog1_quotesSingleLine(t *testing.T) {
	x := newJsFileRegexParserTest(t)
	contents := "var x = \"/** @provideGoog */\";"

	expected := jsFileRegexParserTest_builder().
		SetProvides([]string{}).
		SetGoogModule(false).
		Build()
	result := x.parser.SetIncludeGoogBase(true).ParseFile(jsFileRegexParserTest_SRC_PATH, jsFileRegexParserTest_CLOSURE_PATH, contents)
	x.assertDeps(expected, result)
}

func TestJsFileRegexParser_IncludeGoog1_quotesBeforeComment(t *testing.T) {
	x := newJsFileRegexParserTest(t)
	contents := "var x = \"foo\"; /** @provideGoog */"

	expected := jsFileRegexParserTest_builder().
		SetProvides([]string{"goog"}).
		SetGoogModule(false).
		Build()
	result := x.parser.SetIncludeGoogBase(true).ParseFile(jsFileRegexParserTest_SRC_PATH, jsFileRegexParserTest_CLOSURE_PATH, contents)
	x.assertDeps(expected, result)
}

func TestJsFileRegexParser_IncludeGoog1_multipleQuotes(t *testing.T) {
	x := newJsFileRegexParserTest(t)
	contents := "var x = \"foo\"; var y = \"/** @provideGoog */\";"

	expected := jsFileRegexParserTest_builder().
		SetProvidesArray().
		SetGoogModule(false).
		Build()
	result := x.parser.SetIncludeGoogBase(true).ParseFile(jsFileRegexParserTest_SRC_PATH, jsFileRegexParserTest_CLOSURE_PATH, contents)
	x.assertDeps(expected, result)
}

func TestJsFileRegexParser_IncludeGoog2(t *testing.T) {
	x := newJsFileRegexParserTest(t)
	contents := "goog.require('bar');"

	expected := jsFileRegexParserTest_builder().
		SetRequires([]*DependencyInfo_Require{jsFileRegexParserTest_goog("goog"), jsFileRegexParserTest_goog("bar")}).
		SetGoogModule(false).
		Build()
	result := x.parser.SetIncludeGoogBase(true).ParseFile(jsFileRegexParserTest_SRC_PATH, jsFileRegexParserTest_CLOSURE_PATH, contents)
	x.assertDeps(expected, result)
}

func TestJsFileRegexParser_IncludeGoog3(t *testing.T) {
	x := newJsFileRegexParserTest(t)
	// This is pretending to provide goog, but it really doesn't.
	contents := jsFileRegexParserTest_lines(
		"goog.provide('x');",
		"/**",
		" * the first constant in base.js",
		" */",
		"var COMPILED = false;")

	expected := jsFileRegexParserTest_builder().
		SetProvides([]string{"x"}).
		SetRequires([]*DependencyInfo_Require{jsFileRegexParserTest_goog("goog")}).
		SetGoogModule(false).
		Build()
	result := x.parser.SetIncludeGoogBase(true).ParseFile(jsFileRegexParserTest_SRC_PATH, jsFileRegexParserTest_CLOSURE_PATH, contents)
	x.assertDeps(expected, result)
}

func TestJsFileRegexParser_IncludeGoog4(t *testing.T) {
	x := newJsFileRegexParserTest(t)
	contents := "goog.addDependency('foo', [], []);\n"

	expected := jsFileRegexParserTest_builder().
		SetRequires([]*DependencyInfo_Require{jsFileRegexParserTest_goog("goog")}).
		SetGoogModule(false).
		Build()
	result := x.parser.SetIncludeGoogBase(true).ParseFile(jsFileRegexParserTest_SRC_PATH, jsFileRegexParserTest_CLOSURE_PATH, contents)
	x.assertDeps(expected, result)
}

func TestJsFileRegexParser_ExternsAnnotation_basic_multiline(t *testing.T) {
	x := newJsFileRegexParserTest(t)
	contents := jsFileRegexParserTest_lines(
		"/**",
		" * @externs",
		" */")

	expected := jsFileRegexParserTest_builder().SetHasExternsAnnotation(true).Build()
	result := x.parser.SetIncludeGoogBase(true).ParseFile(jsFileRegexParserTest_SRC_PATH, jsFileRegexParserTest_CLOSURE_PATH, contents)
	x.assertDeps(expected, result)
}

func TestJsFileRegexParser_ExternsAnnotation_basic_oneLine(t *testing.T) {
	x := newJsFileRegexParserTest(t)
	contents := "/** @externs */\n"

	expected := jsFileRegexParserTest_builder().SetHasExternsAnnotation(true).Build()
	result := x.parser.SetIncludeGoogBase(true).ParseFile(jsFileRegexParserTest_SRC_PATH, jsFileRegexParserTest_CLOSURE_PATH, contents)
	x.assertDeps(expected, result)
}

func TestJsFileRegexParser_ExternsAnnotation_blockComment(t *testing.T) {
	x := newJsFileRegexParserTest(t)
	contents := "/* @externs */\n"

	expected := jsFileRegexParserTest_builder().SetHasExternsAnnotation(false).Build()
	result := x.parser.SetIncludeGoogBase(true).ParseFile(jsFileRegexParserTest_SRC_PATH, jsFileRegexParserTest_CLOSURE_PATH, contents)
	x.assertDeps(expected, result)
}

func TestJsFileRegexParser_ExternsAnnotation_blockComment_multiline(t *testing.T) {
	x := newJsFileRegexParserTest(t)
	contents := "/*\n @externs */\n"

	expected := jsFileRegexParserTest_builder().SetHasExternsAnnotation(false).Build()
	result := x.parser.SetIncludeGoogBase(true).ParseFile(jsFileRegexParserTest_SRC_PATH, jsFileRegexParserTest_CLOSURE_PATH, contents)
	x.assertDeps(expected, result)
}

func TestJsFileRegexParser_NoCompileAnnotation_basic_multiline(t *testing.T) {
	x := newJsFileRegexParserTest(t)
	contents := jsFileRegexParserTest_lines(
		"/**",
		" * @nocompile",
		" */")

	expected := jsFileRegexParserTest_builder().
		SetHasNoCompileAnnotation(true).
		Build()
	result := x.parser.SetIncludeGoogBase(true).ParseFile(jsFileRegexParserTest_SRC_PATH, jsFileRegexParserTest_CLOSURE_PATH, contents)
	x.assertDeps(expected, result)
}

func TestJsFileRegexParser_NoCompileAnnotation_basic_oneLine(t *testing.T) {
	x := newJsFileRegexParserTest(t)
	contents := "/** @nocompile */\n"

	expected := jsFileRegexParserTest_builder().
		SetHasNoCompileAnnotation(true).
		Build()
	result := x.parser.SetIncludeGoogBase(true).ParseFile(jsFileRegexParserTest_SRC_PATH, jsFileRegexParserTest_CLOSURE_PATH, contents)
	x.assertDeps(expected, result)
}

func TestJsFileRegexParser_NoCompileAnnotation_blockComment(t *testing.T) {
	x := newJsFileRegexParserTest(t)
	contents := "/* @nocompile */\n"

	expected := jsFileRegexParserTest_builder().
		SetHasNoCompileAnnotation(false).
		Build()
	result := x.parser.SetIncludeGoogBase(true).ParseFile(jsFileRegexParserTest_SRC_PATH, jsFileRegexParserTest_CLOSURE_PATH, contents)
	x.assertDeps(expected, result)
}

func TestJsFileRegexParser_ParseProvidesAndWrappedGoogModule(t *testing.T) {
	x := newJsFileRegexParserTest(t)
	contents := jsFileRegexParserTest_lines(
		"goog.loadModule(function(){\"use strict\";goog.module('yes1');",
		"goog.provide('my.provide');",
		"var yes2=goog.require('yes2');",
		"var C=goog.require(\"a.b.C\");",
		"const {",
		"  D,",
		"  E",
		"}=goog.require(\"a.b.d\");});")

	expected := jsFileRegexParserTest_builder().
		SetProvides([]string{"yes1", "my.provide"}).
		SetRequires([]*DependencyInfo_Require{jsFileRegexParserTest_goog("yes2"), jsFileRegexParserTest_goog("a.b.C"), jsFileRegexParserTest_goog("a.b.d")}).
		SetLoadFlags(jsFileRegexParserTest_flags()).
		Build() // wrapped modules aren't marked as modules

	result := x.parser.ParseFile(jsFileRegexParserTest_SRC_PATH, jsFileRegexParserTest_CLOSURE_PATH, contents)

	x.assertDeps(expected, result)
}

func TestJsFileRegexParser_Es6AndWrappedGoogModuleIsError(t *testing.T) {
	x := newJsFileRegexParserTest(t)
	contents := jsFileRegexParserTest_lines(
		"goog.loadModule(function(){\"use strict\";goog.module('yes1');});",
		"export {};")

	x.parser.ParseFile(jsFileRegexParserTest_SRC_PATH, jsFileRegexParserTest_CLOSURE_PATH, contents)

	assertEqual(t, "errors", len(x.errorManager.GetErrors()), 0)
	warnings := x.errorManager.GetWarnings()
	if len(warnings) != 1 || warnings[0].Type() != ModuleLoader_MODULE_CONFLICT {
		t.Errorf("warnings: got %v, want exactly one MODULE_CONFLICT", warnings)
	}
}
