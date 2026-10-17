// Port of com/google/javascript/jscomp/LateEs6ToEs3ConverterTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/LateEs6ToEs3ConverterTest.java)

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

package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// LateEs6ToEs3ConverterTest has the test cases for ES6 transpilation. This
// tests LateEs6ToEs3Converter
type LateEs6ToEs3ConverterTest struct {
	CompilerTestCaseBase
}

// Sample first script for use when testing non-first scripts
//
// Java: LateEs6ToEs3ConverterTest.SCRIPT1
const LateEs6ToEs3ConverterTest_SCRIPT1 = "var x;"

// The pass generates variable names consisting of filepath based uniqueIDs e.g.
// `$jscomp$templatelit$12345`. Use generic name `TAGGED_TEMPLATE_TMP_VAR` in test sources for
// readability.
//
// Java: LateEs6ToEs3ConverterTest.REPLACEMENT_PREFIXES
var LateEs6ToEs3ConverterTest_REPLACEMENT_PREFIXES = func() *jutil.LinkedHashMap[string, string] {
	m := jutil.NewLinkedHashMap[string, string]()
	m.Put("TAGGED_TEMPLATE_TMP_VAR", "$jscomp$templatelit$")
	m.Put("COMP_PROP_VAR", "$jscomp$compprop")
	return m
}()

// Java: LateEs6ToEs3ConverterTest.RUNTIME_STUBS
const LateEs6ToEs3ConverterTest_RUNTIME_STUBS = `/** @const */
var $jscomp = {};
$jscomp.createTemplateTagFirstArg = function(arrayStrings) {};
$jscomp.createTemplateTagFirstArgWithRaw = function(anotherArray, rawArrayStrings) {};
`

// Java: LateEs6ToEs3ConverterTest()
func NewLateEs6ToEs3ConverterTest() *LateEs6ToEs3ConverterTest {
	c := &LateEs6ToEs3ConverterTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, CompilerTestCase_MINIMAL_EXTERNS+LateEs6ToEs3ConverterTest_RUNTIME_STUBS)
	return c
}

func TestLateEs6ToEs3ConverterTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewLateEs6ToEs3ConverterTest() })
}

// JUnitBefore runs the @Before methods, superclass ones first:
// CompilerTestCase.setUp(), then customSetUp().
func (c *LateEs6ToEs3ConverterTest) JUnitBefore() {
	c.CompilerTestCaseBase.JUnitBefore()
	c.CustomSetUp()
}

// Java: LateEs6ToEs3ConverterTest.customSetUp()
func (c *LateEs6ToEs3ConverterTest) CustomSetUp() {
	c.EnableNormalize()
	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_2015)
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT3)
	c.EnableTypeCheck()
	c.EnableTypeInfoValidation()
	c.ReplaceTypesWithColors()
	c.EnableMultistageCompilation()
	c.SetGenericNameReplacements(LateEs6ToEs3ConverterTest_REPLACEMENT_PREFIXES)
}

// Java: LateEs6ToEs3ConverterTest.getProcessor(Compiler)
func (c *LateEs6ToEs3ConverterTest) GetProcessor(compiler *Compiler) CompilerPass {
	optimizer := NewPhaseOptimizer(compiler, nil)
	optimizer.AddOneTimePass(
		CompilerTestCase_MakePassFactory(
			"injectTranspilationRuntimeLibraries", func(abstractCompiler AbstractCompiler) CompilerPass {
				return NewInjectTranspilationRuntimeLibraries(abstractCompiler)
			}))
	optimizer.AddOneTimePass(CompilerTestCase_MakePassFactory("lateEs6ToEs3Converter", func(abstractCompiler AbstractCompiler) CompilerPass {
		return NewLateEs6ToEs3Converter(abstractCompiler)
	}))
	return optimizer
}

// Java: LateEs6ToEs3ConverterTest.testObjectLiteralMemberFunctionDef()
func (c *LateEs6ToEs3ConverterTest) TestObjectLiteralMemberFunctionDef() {
	c.Test(
		"var x = {/** @return {number} */ a() { return 0; } };",
		"var x = {a: function() { return 0; } };")
}

// Java: LateEs6ToEs3ConverterTest.testInitSymbolIterator()
func (c *LateEs6ToEs3ConverterTest) TestInitSymbolIterator() {
	c.Test(
		c.Externs("/** @type {symbol} */ Symbol.iterator;"),
		c.Srcs("var x = {[Symbol.iterator]: function() { return this; }};"),
		c.Expected(
			`var COMP_PROP_VAR$0 = {};
var x = (COMP_PROP_VAR$0[Symbol.iterator] = function() {return this;},
         COMP_PROP_VAR$0)
`))
}

// Java: LateEs6ToEs3ConverterTest.testMethodInObject()
func (c *LateEs6ToEs3ConverterTest) TestMethodInObject() {
	c.Test("var obj = { f() {alert(1); } };", "var obj = { f: function() {alert(1); } };")

	c.Test("var obj = { f() { alert(1); }, x };", "var obj = { f: function() { alert(1); }, x: x };")
}

// Java: LateEs6ToEs3ConverterTest.testComputedPropertiesWithMethod()
func (c *LateEs6ToEs3ConverterTest) TestComputedPropertiesWithMethod() {
	c.Test(
		"var obj = { ['f' + 1]: 1, m() {}, ['g' + 1]: 1, };",
		`var COMP_PROP_VAR$0 = {};
var obj = (COMP_PROP_VAR$0['f' + 1] = 1,
  (COMP_PROP_VAR$0.m = function() {},
     (COMP_PROP_VAR$0['g' + 1] = 1, COMP_PROP_VAR$0)));
`)
}

// Java: LateEs6ToEs3ConverterTest.testComputedProperties()
func (c *LateEs6ToEs3ConverterTest) TestComputedProperties() {
	c.Test(
		"var obj = { ['f' + 1] : 1, ['g' + 1] : 1 };",
		`var COMP_PROP_VAR$0 = {};
var obj = (COMP_PROP_VAR$0['f' + 1] = 1,
  (COMP_PROP_VAR$0['g' + 1] = 1, COMP_PROP_VAR$0));
`)

	c.Test(
		"var obj = { ['f'] : 1};",
		`var COMP_PROP_VAR$0 = {};
var obj = (COMP_PROP_VAR$0['f'] = 1,
  COMP_PROP_VAR$0);
`)

	c.Test(
		"var o = { ['f'] : 1}; var p = { ['g'] : 1};",
		`var COMP_PROP_VAR$0 = {};
var o = (COMP_PROP_VAR$0['f'] = 1,
  COMP_PROP_VAR$0);
var COMP_PROP_VAR$1 = {};
var p = (COMP_PROP_VAR$1['g'] = 1,
  COMP_PROP_VAR$1);
`)

	c.Test(
		"({['f' + 1] : 1})",
		`var COMP_PROP_VAR$0 = {};
(COMP_PROP_VAR$0['f' + 1] = 1,
  COMP_PROP_VAR$0)
`)

	c.Test(
		"({'a' : 2, ['f' + 1] : 1})",
		`var COMP_PROP_VAR$0 = {};
(COMP_PROP_VAR$0['a'] = 2,
  (COMP_PROP_VAR$0['f' + 1] = 1, COMP_PROP_VAR$0));
`)

	c.Test(
		"({['f' + 1] : 1, 'a' : 2})",
		`var COMP_PROP_VAR$0 = {};
(COMP_PROP_VAR$0['f' + 1] = 1,
  (COMP_PROP_VAR$0['a'] = 2, COMP_PROP_VAR$0));
`)

	c.Test(
		"({'a' : 1, ['f' + 1] : 1, 'b' : 1})",
		`var COMP_PROP_VAR$0 = {};
(COMP_PROP_VAR$0['a'] = 1,
  (COMP_PROP_VAR$0['f' + 1] = 1, (COMP_PROP_VAR$0['b'] = 1, COMP_PROP_VAR$0)));
`)

	c.Test(
		"({'a' : x++, ['f' + x++] : 1, 'b' : x++})",
		`var COMP_PROP_VAR$0 = {};
(COMP_PROP_VAR$0['a'] = x++, (COMP_PROP_VAR$0['f' + x++] = 1,
  (COMP_PROP_VAR$0['b'] = x++, COMP_PROP_VAR$0)))
`)

	c.Test(
		"({a : x++, ['f' + x++] : 1, b : x++})",
		`var COMP_PROP_VAR$0 = {};
(COMP_PROP_VAR$0.a = x++, (COMP_PROP_VAR$0['f' + x++] = 1,
  (COMP_PROP_VAR$0.b = x++, COMP_PROP_VAR$0)))
`)

	c.Test(
		"({a, ['f' + 1] : 1})",
		`var COMP_PROP_VAR$0 = {};
  (COMP_PROP_VAR$0.a = a, (COMP_PROP_VAR$0['f' + 1] = 1, COMP_PROP_VAR$0))
`)

	c.Test(
		"({['f' + 1] : 1, a})",
		`var COMP_PROP_VAR$0 = {};
  (COMP_PROP_VAR$0['f' + 1] = 1, (COMP_PROP_VAR$0.a = a, COMP_PROP_VAR$0))
`)

	c.Test(
		"var obj = { [foo]() {}}",
		`var COMP_PROP_VAR$0 = {};
var obj = (COMP_PROP_VAR$0[foo] = function(){}, COMP_PROP_VAR$0)
`)
}

// Java: LateEs6ToEs3ConverterTest.testComputedPropGetterSetter()
func (c *LateEs6ToEs3ConverterTest) TestComputedPropGetterSetter() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)

	c.TestSame("var obj = {get latest () {return undefined;}}")
	c.TestSame("var obj = {set latest (strParam) {}}")
	c.Test(
		"var obj = {'a' : 2, get l () {return null;}, ['f' + 1] : 1}",
		`var COMP_PROP_VAR$0 = {get l () {return null;}};
var obj = (COMP_PROP_VAR$0['a'] = 2, (COMP_PROP_VAR$0['f' + 1] = 1, COMP_PROP_VAR$0));
`)
	c.Test(
		"var obj = {['a' + 'b'] : 2, set l (strParam) {}}",
		`var COMP_PROP_VAR$0 = {set l (strParam) {}};
var obj = (COMP_PROP_VAR$0['a' + 'b'] = 2, COMP_PROP_VAR$0);
`)
}

// Java: LateEs6ToEs3ConverterTest.testComputedPropCannotConvert()
func (c *LateEs6ToEs3ConverterTest) TestComputedPropCannotConvert() {
	c.TestError(
		"var o = { get [foo]() {}}", ReportUntranspilableFeatures_UNTRANSPILABLE_FEATURE_PRESENT)
	c.TestError(
		"var o = { set [foo](val) {}}",
		ReportUntranspilableFeatures_UNTRANSPILABLE_FEATURE_PRESENT)
}

// Java: LateEs6ToEs3ConverterTest.testNoComputedProperties()
func (c *LateEs6ToEs3ConverterTest) TestNoComputedProperties() {
	c.TestSame("({'a' : 1})")
	c.TestSame("({'a' : 1, f : 1, b : 1})")
}

// Java: LateEs6ToEs3ConverterTest.testUntaggedTemplateLiteral()
func (c *LateEs6ToEs3ConverterTest) TestUntaggedTemplateLiteral() {
	c.Test("``", "''")
	c.Test("`\"`", "'\\\"'")
	c.Test("`'`", "\"'\"")
	c.Test("`\\``", "'`'")
	c.Test("`\\\"`", "'\\\"'")
	c.Test("`\\\\\"`", "'\\\\\\\"'")
	c.Test("`\"\\\\`", "'\"\\\\'")
	c.Test("`$$`", "'$$'")
	c.Test("`$$$`", "'$$$'")
	c.Test("`\\$$$`", "'$$$'")
	c.Test("`hello`", "'hello'")
	c.Test("`hello\nworld`", "'hello\\nworld'")
	c.Test("`hello\rworld`", "'hello\\nworld'")
	c.Test("`hello\r\nworld`", "'hello\\nworld'")
	c.Test("`hello\n\nworld`", "'hello\\n\\nworld'")
	c.Test("`hello\\r\\nworld`", "'hello\\r\\nworld'")
	c.Test("`${world}`", "'' + world")
	c.Test("`hello ${world}`", "'hello ' + world")
	c.Test("`${hello} world`", "hello + ' world'")
	c.Test("`${hello}${world}`", "'' + hello + world")
	c.Test("`${a} b ${c} d ${e}`", "a + ' b ' + c + ' d ' + e")
	c.Test("`hello ${a + b}`", "'hello ' + (a + b)")
	c.Test("`hello ${a, b, c}`", "'hello ' + (a, b, c)")
	c.Test("`hello ${a ? b : c}${a * b}`", "'hello ' + (a ? b : c) + (a * b)")
}

// Java: LateEs6ToEs3ConverterTest.testTaggedTemplateLiteral_singleScript()
func (c *LateEs6ToEs3ConverterTest) TestTaggedTemplateLiteral_singleScript() {
	c.Test(
		"tag``",
		`/** @noinline */ var TAGGED_TEMPLATE_TMP_VAR$0 =
    $jscomp.createTemplateTagFirstArg(['']);
tag(TAGGED_TEMPLATE_TMP_VAR$0);
`)

	c.Test(
		"tag`${hello} world`",
		`/** @noinline */ var TAGGED_TEMPLATE_TMP_VAR$0 =
    $jscomp.createTemplateTagFirstArg(['', ' world']);
tag(TAGGED_TEMPLATE_TMP_VAR$0, hello);
`)

	c.Test(
		"tag`${hello} ${world}`",
		`/** @noinline */ var TAGGED_TEMPLATE_TMP_VAR$0 =
    $jscomp.createTemplateTagFirstArg(['', ' ', '']);
tag(TAGGED_TEMPLATE_TMP_VAR$0, hello, world);
`)

	c.Test(
		"tag`\"`",
		`/** @noinline */ var TAGGED_TEMPLATE_TMP_VAR$0 =
    $jscomp.createTemplateTagFirstArg(['"']);
tag(TAGGED_TEMPLATE_TMP_VAR$0);
`)

	// The cooked string and the raw string are different.
	// Note that this test is tricky to read, because any escape sequences will be escaped twice.
	// This table is helpful:
	//
	//     Java String    JavaScript String      JavaScript Value
	//
	//     ----------------------------------------------------------------
	//     \t        ->   <tab character>    -> <tab character> (length: 1)
	//     \\t       ->   \t                 -> <tab character> (length: 1)
	//     \\\t      ->   \<tab character>   -> <tab character> (length: 1)
	//     \\\\t     ->   \\t                -> \t              (length: 2)
	//
	c.Test(
		"tag`a\\tb`",
		`/** @noinline */ var TAGGED_TEMPLATE_TMP_VAR$0 =
    $jscomp.createTemplateTagFirstArgWithRaw(['a\tb'],['a\\tb']);
tag(TAGGED_TEMPLATE_TMP_VAR$0);
`)

	c.Test(
		"tag()`${hello} world`",
		`/** @noinline */ var TAGGED_TEMPLATE_TMP_VAR$0 =
    $jscomp.createTemplateTagFirstArg(['', ' world']);
tag()(TAGGED_TEMPLATE_TMP_VAR$0, hello);
`)

	c.Test(
		c.Externs(LateEs6ToEs3ConverterTest_RUNTIME_STUBS, "var a = {}; a.b;"),
		c.Srcs("a.b`${hello} world`"),
		c.Expected(
			`/** @noinline */ var TAGGED_TEMPLATE_TMP_VAR$0 =
    $jscomp.createTemplateTagFirstArg(['', ' world']);
a.b(TAGGED_TEMPLATE_TMP_VAR$0, hello);
`))

	c.Test(
		c.Externs(LateEs6ToEs3ConverterTest_RUNTIME_STUBS, "var a = {}; var b; a[b];"),
		c.Srcs("a[b]`${hello} world`"),
		c.Expected(
			`/** @noinline */ var TAGGED_TEMPLATE_TMP_VAR$0 =
    $jscomp.createTemplateTagFirstArg(['', ' world']);
a[b](TAGGED_TEMPLATE_TMP_VAR$0, hello);
`))
	// https://github.com/google/closure-compiler/issues/1299
	c.Test(
		"tag`<p class=\"foo\">${x}</p>`",
		`/** @noinline */ var TAGGED_TEMPLATE_TMP_VAR$0 =
    $jscomp.createTemplateTagFirstArg(['<p class="foo">','</p>']);
tag(TAGGED_TEMPLATE_TMP_VAR$0, x);
`)
	c.Test(
		"tag`<p class='foo'>${x}</p>`",
		`/** @noinline */ var TAGGED_TEMPLATE_TMP_VAR$0 =
    $jscomp.createTemplateTagFirstArg(['<p class=\'foo\'>','</p>']);
tag(TAGGED_TEMPLATE_TMP_VAR$0, x);
`)

	// invalid escape sequences result in undefined cooked string
	c.Test(
		"tag`\\unicode`",
		`/** @noinline */ var TAGGED_TEMPLATE_TMP_VAR$0 =
    $jscomp.createTemplateTagFirstArgWithRaw([void 0],['\\unicode']);
tag(TAGGED_TEMPLATE_TMP_VAR$0);
`)
}

// Tests that calls to `createTemplateTagFirstArg` are inserted after its injected declaration
// when input is a single script
// Java: LateEs6ToEs3ConverterTest.testTaggedTemplateLiteral_insertPosition_singleScript()
func (c *LateEs6ToEs3ConverterTest) TestTaggedTemplateLiteral_insertPosition_singleScript() {
	c.Test(
		c.Externs(""), // clear the default externs which contain the runtime stubs
		c.Srcs(
			LateEs6ToEs3ConverterTest_RUNTIME_STUBS+
				"var a = {};\ntag``; var b;\n"),
		c.Expected(
			LateEs6ToEs3ConverterTest_RUNTIME_STUBS+
				`/** @noinline */ var TAGGED_TEMPLATE_TMP_VAR$0 =
    $jscomp.createTemplateTagFirstArg(['']);
var a = {};
tag(TAGGED_TEMPLATE_TMP_VAR$0);
var b;
`))

	c.Test(
		c.Externs(""), // clear the default externs which contain the runtime stubs
		c.Srcs(
			LateEs6ToEs3ConverterTest_RUNTIME_STUBS+
				"var a = {};\nfunction foo() {tag``;}\n"),
		c.Expected(
			LateEs6ToEs3ConverterTest_RUNTIME_STUBS+
				`/** @noinline */ var TAGGED_TEMPLATE_TMP_VAR$0 =
    $jscomp.createTemplateTagFirstArg(['']);
var a = {};
function foo() {tag(TAGGED_TEMPLATE_TMP_VAR$0);}
`))

	c.Test(
		c.Externs(""), // clear the default externs which contain the runtime stubs
		c.Srcs(
			LateEs6ToEs3ConverterTest_RUNTIME_STUBS+
				"var a = {};\nfunction foo() {function bar() {tag``;}}\n"),
		c.Expected(
			LateEs6ToEs3ConverterTest_RUNTIME_STUBS+
				`/** @noinline */ var TAGGED_TEMPLATE_TMP_VAR$0 =
    $jscomp.createTemplateTagFirstArg(['']);
var a = {};
function foo() {function bar() {tag(TAGGED_TEMPLATE_TMP_VAR$0);}}
`))
}

// When input is multiple scripts, tests that calls to `createTemplateTagFirstArg` are inserted at
// the top for any subsequent (non-first) script.
// Java: LateEs6ToEs3ConverterTest.testTaggedTemplateLiteral_insertPosition_multipleScripts()
func (c *LateEs6ToEs3ConverterTest) TestTaggedTemplateLiteral_insertPosition_multipleScripts() {
	c.Test(
		c.Externs(""), // clear the default externs which contain the runtime stubs
		c.Srcs(
			LateEs6ToEs3ConverterTest_RUNTIME_STUBS+LateEs6ToEs3ConverterTest_SCRIPT1,
			"var a = {};\ntag``; var b;\n"),
		c.Expected(
			LateEs6ToEs3ConverterTest_RUNTIME_STUBS+LateEs6ToEs3ConverterTest_SCRIPT1,
			`/** @noinline */ var TAGGED_TEMPLATE_TMP_VAR$0 =
    $jscomp.createTemplateTagFirstArg(['']);
var a = {};
tag(TAGGED_TEMPLATE_TMP_VAR$0); var b;
`))

	c.Test(
		c.Externs(""), // clear the default externs which contain the runtime stubs
		c.Srcs(
			LateEs6ToEs3ConverterTest_RUNTIME_STUBS+LateEs6ToEs3ConverterTest_SCRIPT1,
			"var a = {};\nfunction foo() {tag``;}\n"),
		c.Expected(
			LateEs6ToEs3ConverterTest_RUNTIME_STUBS+LateEs6ToEs3ConverterTest_SCRIPT1,
			`/** @noinline */ var TAGGED_TEMPLATE_TMP_VAR$0 =
    $jscomp.createTemplateTagFirstArg(['']);
var a = {};
function foo() {tag(TAGGED_TEMPLATE_TMP_VAR$0);}
`))

	c.Test(
		c.Externs(""), // clear the default externs which contain the runtime stubs
		c.Srcs(
			LateEs6ToEs3ConverterTest_SCRIPT1+LateEs6ToEs3ConverterTest_RUNTIME_STUBS,
			"var a = {};\nfunction foo() {function bar() {tag``;}}\n"),
		c.Expected(
			LateEs6ToEs3ConverterTest_SCRIPT1+LateEs6ToEs3ConverterTest_RUNTIME_STUBS,
			`/** @noinline */ var TAGGED_TEMPLATE_TMP_VAR$0 =
    $jscomp.createTemplateTagFirstArg(['']);
var a = {};
function foo() {function bar() {tag(TAGGED_TEMPLATE_TMP_VAR$0);}}
`))
}

// Java: LateEs6ToEs3ConverterTest.testComputedPropertyInLoopHeader()
func (c *LateEs6ToEs3ConverterTest) TestComputedPropertyInLoopHeader() {
	c.Test(
		"for (; {[foo]: bar}; ) {}",
		`var COMP_PROP_VAR$0;
for (; (COMP_PROP_VAR$0 = {}, (COMP_PROP_VAR$0[foo] = bar, COMP_PROP_VAR$0)); ) {}
`)

	c.Test(
		"while ({[foo]: bar}) {}",
		`var COMP_PROP_VAR$0;
for (; (COMP_PROP_VAR$0 = {}, (COMP_PROP_VAR$0[foo] = bar, COMP_PROP_VAR$0)); ) {}
`)

	c.Test(
		"do {} while ({[foo]: bar});",
		`var COMP_PROP_VAR$0;
do {} while ((COMP_PROP_VAR$0 = {}, (COMP_PROP_VAR$0[foo] = bar, COMP_PROP_VAR$0)));
`)
}

// Java: LateEs6ToEs3ConverterTest.testComputedPropertyInLoopBody()
func (c *LateEs6ToEs3ConverterTest) TestComputedPropertyInLoopBody() {
	c.Test(
		"while (true) { var x = {[foo]: bar}; }",
		`for (; true;) {
  var COMP_PROP_VAR$0 = {};
  var x = (COMP_PROP_VAR$0[foo] = bar, COMP_PROP_VAR$0);
}
`)
}

// Java: LateEs6ToEs3ConverterTest.testUnicodeEscapes()
func (c *LateEs6ToEs3ConverterTest) TestUnicodeEscapes() {
	c.Test("var \\u{73} = '\\u{2603}'", "var s = '☃'")  // ☃
	c.Test("var \\u{63} = '\\u{1f42a}'", "var c = '🐪'") // 🐪
	c.Test("var str = `begin\\u{2026}end`", "var str = 'begin\\u2026end'")
}
