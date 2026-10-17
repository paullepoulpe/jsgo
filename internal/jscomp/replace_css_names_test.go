// Port of com/google/javascript/jscomp/ReplaceCssNamesTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/ReplaceCssNamesTest.java)

/*
 * Copyright 2009 The Closure Compiler Authors.
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

// ReplaceCssNamesTest tests ReplaceCssNames.
//
// Java: public final class ReplaceCssNamesTest extends CompilerTestCase
type ReplaceCssNamesTest struct {
	CompilerTestCaseBase

	// Whether to pass the map of replacements as opposed to null
	useReplacementMap bool

	// Map of replacements to use during the test.
	replacementMap *jutil.LinkedHashMap[string, string]

	replacementMapFull *jutil.LinkedHashMap[string, string]

	renamingMap CssRenamingMap
	skiplist    jutil.Set[string]

	// Go-only: language-adapter — Sets.newHashSet() is only compared
	// without order (containsExactly), so insertion order is used.
	cssNames *jutil.LinkedHashSet[string]
}

// Java: ReplaceCssNamesTest()
func NewReplaceCssNamesTest() *ReplaceCssNamesTest {
	c := &ReplaceCssNamesTest{}
	c.InitCompilerTestCaseBaseWithExterns(c,
		CompilerTestCase_DEFAULT_EXTERNS+
			`Object.prototype.getClass;
goog.getCssName;
goog.setCssNameMapping;
`)
	c.replacementMap = jutil.NewImmutableMapBuilder[string, string]().
		Put("active", "a").
		Put("buttonbar", "b").
		Put("colorswatch", "c").
		Put("disabled", "d").
		Put("elephant", "e").
		Put("footer", "f").
		Put("goog", "g").
		Put("fooStylesBar", "fsr").
		Put("fooStylesBaz", "fsz").
		Put("--foo-bar", "--fb").
		Put("---foo-baz", "--fbz").
		Put("--foo-bar-baz--qux", "--fbzq").
		BuildOrThrow()
	c.replacementMapFull = jutil.NewImmutableMapBuilder[string, string]().
		Put("long-prefix", "h").
		Put("suffix1", "i").
		Put("unrelated-word", "k").
		Put("unrelated", "l").
		Put("long-suffix", "m").
		Put("long-prefix-suffix1", "h-i").
		Put("--foo-bar", "--fb").
		Put("---foo-baz", "--fbz").
		Put("--foo-bar-baz--qux", "--fbzq").
		BuildOrThrow()
	return c
}

func TestReplaceCssNamesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewReplaceCssNamesTest() })
}

// Java: ReplaceCssNamesTest.getProcessor(Compiler)
func (c *ReplaceCssNamesTest) GetProcessor(compiler *Compiler) CompilerPass {
	var symbolMap CssRenamingMap
	if c.useReplacementMap {
		symbolMap = c.renamingMap
	}
	cssNames := c.cssNames
	return NewReplaceCssNames(
		compiler, symbolMap, ReplaceCssNames_CssNameCollectorFunc(func(name string) { cssNames.Add(name) }), c.skiplist)
}

// replaceCssNamesTest_partialMap is the anonymous CssRenamingMap.ByPart of
// getPartialMap.
type replaceCssNamesTest_partialMap struct {
	CssRenamingMap_ByPartBase
	outer *ReplaceCssNamesTest
}

// Java: ReplaceCssNamesTest.getPartialMap().new ByPart() {...}.get(String)
func (m *replaceCssNamesTest_partialMap) Get(value string) *string {
	if v, ok := m.outer.replacementMap.Lookup(value); ok {
		return &v
	}
	return nil
}

// Java: ReplaceCssNamesTest.getPartialMap()
func (c *ReplaceCssNamesTest) GetPartialMap() CssRenamingMap {
	return &replaceCssNamesTest_partialMap{outer: c}
}

// replaceCssNamesTest_fullMap is the anonymous CssRenamingMap.ByWhole of
// getFullMap.
type replaceCssNamesTest_fullMap struct {
	CssRenamingMap_ByWholeBase
	outer *ReplaceCssNamesTest
}

// Java: ReplaceCssNamesTest.getFullMap().new ByWhole() {...}.get(String)
func (m *replaceCssNamesTest_fullMap) Get(value string) *string {
	if v, ok := m.outer.replacementMapFull.Lookup(value); ok {
		return &v
	}
	return nil
}

// Java: ReplaceCssNamesTest.getFullMap()
func (c *ReplaceCssNamesTest) GetFullMap() CssRenamingMap {
	return &replaceCssNamesTest_fullMap{outer: c}
}

// Java: ReplaceCssNamesTest.setUp()
func (c *ReplaceCssNamesTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableTypeCheck()
	c.EnableRewriteClosureCode()
	c.cssNames = jutil.NewLinkedHashSet[string]()
	c.useReplacementMap = true
	c.renamingMap = c.GetPartialMap()
}

// replaceCssNamesTest_silentErrorManager is the anonymous BasicErrorManager
// of testNoSymbolMapStripsCallAndDoesntIssueWarnings.
type replaceCssNamesTest_silentErrorManager struct {
	BasicErrorManagerBase
}

func newReplaceCssNamesTest_silentErrorManager() *replaceCssNamesTest_silentErrorManager {
	m := &replaceCssNamesTest_silentErrorManager{}
	m.InitBasicErrorManagerBase(m)
	return m
}

func (m *replaceCssNamesTest_silentErrorManager) PrintSummary() {}

func (m *replaceCssNamesTest_silentErrorManager) Println(level *CheckLevel, error *JSError) {}

// Java: ReplaceCssNamesTest.testSimpleValidMapping()
func (c *ReplaceCssNamesTest) TestSimpleValidMapping() {
	c.Test("function f() { goog.setCssNameMapping({}); }", "function f() {}")
}

// Java: ReplaceCssNamesTest.testValidSetCssNameMapping()
func (c *ReplaceCssNamesTest) TestValidSetCssNameMapping() {
	c.Test(
		`goog.setCssNameMapping({foo:'bar',"biz":'baz'});
const x = goog.getCssName('foo'),
  y = goog.getCssName('biz'),
  z = goog.getCssName('foo-biz');
`,
		"const x = 'bar', y = 'baz', z = 'bar-baz';")
}

// Java: ReplaceCssNamesTest.testValidSetCssNameMapping_byPart()
func (c *ReplaceCssNamesTest) TestValidSetCssNameMapping_byPart() {
	c.Test(
		`goog.setCssNameMapping({foo:'bar',"biz":'baz'}, 'BY_PART');
const x = goog.getCssName('foo'),
  y = goog.getCssName('biz'),
  z = goog.getCssName('foo-biz');
`,
		"const x = 'bar', y = 'baz', z = 'bar-baz';")
}

// Java: ReplaceCssNamesTest.testValidSetCssNameMapping_byPart_withCssVariable()
func (c *ReplaceCssNamesTest) TestValidSetCssNameMapping_byPart_withCssVariable() {
	c.Test(
		`goog.setCssNameMapping({classy: 'c',selector: 's',"--foo":'--f',"--biz-baz":'--b'}, 'BY_PART');
const x = goog.getCssName('--foo'),
  y = goog.getCssName('--biz-baz'),
  z = goog.getCssName('classy-selector');
`,
		"const x = '--f', y = '--b', z = 'c-s';")
}

// Java: ReplaceCssNamesTest.testValidSetCssNameMapping_byWhole()
func (c *ReplaceCssNamesTest) TestValidSetCssNameMapping_byWhole() {
	c.Test(
		`goog.setCssNameMapping({foo:'bar',"biz":'baz'}, 'BY_WHOLE');
const x = goog.getCssName('foo'),
  y = goog.getCssName('biz'),
  z = goog.getCssName('foo-biz');
`,
		"const x = 'bar', y = 'baz', z = 'foo-biz';")

	c.Test(
		`goog.setCssNameMapping({foo:'bar',"biz":'baz','biz-foo':'baz-bar'}, 'BY_WHOLE');
const x = goog.getCssName('foo'),
  y = goog.getCssName('biz'),
  z = goog.getCssName('foo-biz');
`,
		"const x = 'bar', y = 'baz', z = 'foo-biz';")
}

// Java: ReplaceCssNamesTest.testValidSetCssNameMapping_byWhole_withCssVariable()
func (c *ReplaceCssNamesTest) TestValidSetCssNameMapping_byWhole_withCssVariable() {
	c.Test(
		`goog.setCssNameMapping({classy: 'c',"--foo":'--f',"--biz-baz":'--b'}, 'BY_WHOLE');
const x = goog.getCssName('--foo'),
  y = goog.getCssName('--biz-baz'),
  z = goog.getCssName('classy');
`,
		"const x = '--f', y = '--b', z = 'c';")
}

// Java: ReplaceCssNamesTest.testDoNotUseReplacementMap()
func (c *ReplaceCssNamesTest) TestDoNotUseReplacementMap() {
	c.useReplacementMap = false
	c.Test(
		"var x = goog.getCssName('goog-footer-active')", //
		"var x = 'goog-footer-active'")
	c.Test(
		"el.className = goog.getCssName('goog-colorswatch-disabled')",
		"el.className = 'goog-colorswatch-disabled'")
	c.Test(
		"setClass(goog.getCssName('active-buttonbar'))", //
		"setClass('active-buttonbar')")
	AssertThat(c.cssNames).
		ContainsExactly("goog-footer-active", "goog-colorswatch-disabled", "active-buttonbar")
}

// Java: ReplaceCssNamesTest.testOneArgWithUnknownStringLiterals()
func (c *ReplaceCssNamesTest) TestOneArgWithUnknownStringLiterals() {
	c.Test(
		"var x = goog.getCssName('unknown')", //
		"var x = 'unknown'",
		c.Warning(ReplaceCssNames_UNKNOWN_SYMBOL_WARNING))
	c.Test(
		"el.className = goog.getCssName('ooo')",
		"el.className = 'ooo'",
		c.Warning(ReplaceCssNames_UNKNOWN_SYMBOL_WARNING))
	c.Test(
		"setClass(goog.getCssName('ab'))", //
		"setClass('ab')",
		c.Warning(ReplaceCssNames_UNKNOWN_SYMBOL_WARNING))
}

// Java: ReplaceCssNamesTest.oneArgWithSimpleStringLiterals()
func (c *ReplaceCssNamesTest) OneArgWithSimpleStringLiterals() {
	c.Test(
		"var x = goog.getCssName('buttonbar')", //
		"var x = 'b'")
	c.Test(
		"el.className = goog.getCssName('colorswatch')", //
		"el.className = 'c'")
	c.Test(
		"setClass(goog.getCssName('elephant'))", //
		"setClass('e')")
	AssertThat(c.cssNames).ContainsExactly("buttonbar", "colorswatch", "elephant")
}

// Java: ReplaceCssNamesTest.testOneArgWithSimpleStringLiterals()
func (c *ReplaceCssNamesTest) TestOneArgWithSimpleStringLiterals() {
	c.OneArgWithSimpleStringLiterals()
}

// Java: ReplaceCssNamesTest.oneArgWithCompositeClassNames()
func (c *ReplaceCssNamesTest) OneArgWithCompositeClassNames() {
	c.Test(
		"var x = goog.getCssName('goog-footer-active')", //
		"var x = 'g-f-a'")
	c.Test(
		"el.className = goog.getCssName('goog-colorswatch-disabled')", //
		"el.className = 'g-c-d'")
	c.Test(
		"setClass(goog.getCssName('active-buttonbar'))", //
		"setClass('a-b')")
	AssertThat(c.cssNames).
		ContainsExactly("goog-footer-active", "goog-colorswatch-disabled", "active-buttonbar")
}

// Java: ReplaceCssNamesTest.testOneArgWithCompositeClassNames()
func (c *ReplaceCssNamesTest) TestOneArgWithCompositeClassNames() {
	c.OneArgWithCompositeClassNames()
}

// Java: ReplaceCssNamesTest.testOneArgWithCompositeClassNamesFull()
func (c *ReplaceCssNamesTest) TestOneArgWithCompositeClassNamesFull() {
	c.renamingMap = c.GetFullMap()

	c.Test(
		"var x = goog.getCssName('long-prefix')", //
		"var x = 'h'")
	c.Test(
		"var x = goog.getCssName('long-prefix-suffix1')", //
		"var x = 'h-i'")
	c.Test(
		"var x = goog.getCssName('unrelated')", //
		"var x = 'l'")
	c.Test(
		"var x = goog.getCssName('unrelated-word')", //
		"var x = 'k'")
}

// Java: ReplaceCssNamesTest.testOneArgWithCssVariables_byWhole()
func (c *ReplaceCssNamesTest) TestOneArgWithCssVariables_byWhole() {
	c.renamingMap = c.GetFullMap()

	c.Test(
		"var x = goog.getCssName('--foo-bar')", //
		"var x = '--fb'")

	c.Test(
		"var x = goog.getCssName('---foo-baz')", //
		"var x = '--fbz'")

	c.Test(
		"var x = goog.getCssName('--foo-bar-baz--qux')", //
		"var x = '--fbzq'")
}

// Java: ReplaceCssNamesTest.testOneArgWithCompositeClassNamesWithUnknownParts()
func (c *ReplaceCssNamesTest) TestOneArgWithCompositeClassNamesWithUnknownParts() {
	c.Test(
		"var x = goog.getCssName('goog-header-active')",
		"var x = 'goog-header-active'",
		c.Warning(ReplaceCssNames_UNKNOWN_SYMBOL_WARNING))
	c.Test(
		"el.className = goog.getCssName('goog-colorswatch-focussed')",
		"el.className = 'goog-colorswatch-focussed'",
		c.Warning(ReplaceCssNames_UNKNOWN_SYMBOL_WARNING))
	c.Test(
		"setClass(goog.getCssName('inactive-buttonbar'))",
		"setClass('inactive-buttonbar')",
		c.Warning(ReplaceCssNames_UNKNOWN_SYMBOL_WARNING))
}

// Java: ReplaceCssNamesTest.testOneArgWithCssVariable()
func (c *ReplaceCssNamesTest) TestOneArgWithCssVariable() {
	c.Test(
		"var x = goog.getCssName('--foo-bar')", //
		"var x = '--fb'")

	c.Test(
		"var x = goog.getCssName('---foo-baz')", //
		"var x = '--fbz'")

	c.Test(
		"var x = goog.getCssName('--foo-bar-baz--qux')", //
		"var x = '--fbzq'")
}

// Java: ReplaceCssNamesTest.testTwoArgsWithStringLiterals()
func (c *ReplaceCssNamesTest) TestTwoArgsWithStringLiterals() {
	c.TestError("var x = goog.getCssName('header', 'active')", ReplaceCssNames_UNEXPECTED_STRING_LITERAL_ERROR)
	c.TestError(
		"el.className = goog.getCssName('footer', window)",
		ReplaceCssNames_STRING_LITERAL_EXPECTED_ERROR)
	c.TestError(
		"setClass(goog.getCssName('buttonbar', 'disabled'))", ReplaceCssNames_UNEXPECTED_STRING_LITERAL_ERROR)
}

// Java: ReplaceCssNamesTest.testNestedCall()
func (c *ReplaceCssNamesTest) TestNestedCall() {
	c.TestError(
		"setClass(goog.getCssName(goog.getCssName('buttonbar'), 'active'))",
		ReplaceCssNames_NESTED_CALL_ERROR)
}

// Java: ReplaceCssNamesTest.testTwoArsWithVariableFirstArg()
func (c *ReplaceCssNamesTest) TestTwoArsWithVariableFirstArg() {
	c.Test(
		"var x = goog.getCssName(baseClass, 'active')", //
		"var x = baseClass + '-a'")
	c.Test(
		"el.className = goog.getCssName((new Object).getClass(), 'disabled')",
		"el.className = (new Object).getClass() + '-d'")
	c.Test(
		"setClass(goog.getCssName(BASE_CLASS, 'disabled'))", //
		"setClass(BASE_CLASS + '-d')")
}

// Java: ReplaceCssNamesTest.testTwoArgsWithVariableFirstArgFull()
func (c *ReplaceCssNamesTest) TestTwoArgsWithVariableFirstArgFull() {
	c.renamingMap = c.GetFullMap()

	c.Test(
		"var x = goog.getCssName(baseClass, 'long-suffix')", //
		"var x = baseClass + '-m'")
}

// Java: ReplaceCssNamesTest.testZeroArguments()
func (c *ReplaceCssNamesTest) TestZeroArguments() {
	c.TestError("goog.getCssName()", ReplaceCssNames_INVALID_NUM_ARGUMENTS_ERROR)
}

// Java: ReplaceCssNamesTest.testManyArguments()
func (c *ReplaceCssNamesTest) TestManyArguments() {
	c.TestError("goog.getCssName('a', 'b', 'c')", ReplaceCssNames_INVALID_NUM_ARGUMENTS_ERROR)
	c.TestError("goog.getCssName('a', 'b', 'c', 'd')", ReplaceCssNames_INVALID_NUM_ARGUMENTS_ERROR)
	c.TestError(
		"goog.getCssName('a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i')",
		ReplaceCssNames_INVALID_NUM_ARGUMENTS_ERROR)
}

// Java: ReplaceCssNamesTest.testNonStringArgument()
func (c *ReplaceCssNamesTest) TestNonStringArgument() {
	c.TestError("goog.getCssName(window);", ReplaceCssNames_STRING_LITERAL_EXPECTED_ERROR)
	c.TestError("goog.getCssName(555);", ReplaceCssNames_STRING_LITERAL_EXPECTED_ERROR)
	c.TestError("goog.getCssName([]);", ReplaceCssNames_STRING_LITERAL_EXPECTED_ERROR)
	c.TestError("goog.getCssName({});", ReplaceCssNames_STRING_LITERAL_EXPECTED_ERROR)
	c.TestError("goog.getCssName(null);", ReplaceCssNames_STRING_LITERAL_EXPECTED_ERROR)
	c.TestError("goog.getCssName(undefined);", ReplaceCssNames_STRING_LITERAL_EXPECTED_ERROR)

	c.TestError("goog.getCssName(baseClass, window);", ReplaceCssNames_STRING_LITERAL_EXPECTED_ERROR)
	c.TestError("goog.getCssName(baseClass, 555);", ReplaceCssNames_STRING_LITERAL_EXPECTED_ERROR)
	c.TestError("goog.getCssName(baseClass, []);", ReplaceCssNames_STRING_LITERAL_EXPECTED_ERROR)
	c.TestError("goog.getCssName(baseClass, {});", ReplaceCssNames_STRING_LITERAL_EXPECTED_ERROR)
	c.TestError("goog.getCssName(baseClass, null);", ReplaceCssNames_STRING_LITERAL_EXPECTED_ERROR)
	c.TestError(
		"goog.getCssName(baseClass, undefined);", ReplaceCssNames_STRING_LITERAL_EXPECTED_ERROR)
	c.TestError("goog.getCssName('foo', 3);", ReplaceCssNames_STRING_LITERAL_EXPECTED_ERROR)
}

// Java: ReplaceCssNamesTest.testNoSymbolMapStripsCallAndDoesntIssueWarnings()
func (c *ReplaceCssNamesTest) TestNoSymbolMapStripsCallAndDoesntIssueWarnings() {
	input := "[goog.getCssName('test'), goog.getCssName(base, 'active')]"
	compiler := NewCompiler()
	options := NewCompilerOptions()
	options.SetEmitUseStrict(false)
	compiler.InitOptions(options)
	errorMan := newReplaceCssNamesTest_silentErrorManager()
	compiler.SetErrorManager(errorMan)
	root := compiler.ParseTestCode(input)
	c.useReplacementMap = false
	replacer := NewReplaceCssNames(compiler, nil, ReplaceCssNames_CssNameCollectorFunc(func(unused string) {}), nil)
	replacer.Process(nil, root)
	AssertThat(compiler.ToSourceOfNode(root)).IsEqualTo("[\"test\",base+\"-active\"]")
	AssertWithMessage("There should be no errors").That(errorMan.GetErrorCount()).IsEqualTo(0)
	AssertWithMessage("There should be no warnings").That(errorMan.GetWarningCount()).IsEqualTo(0)
}

// Java: ReplaceCssNamesTest.testWhitelistByPart()
func (c *ReplaceCssNamesTest) TestWhitelistByPart() {
	c.skiplist = jutil.NewLinkedHashSet[string]("goog", "elephant")
	c.Test(
		"var x = goog.getCssName('goog')", //
		"var x = 'goog'")
	c.Test(
		"var x = goog.getCssName('elephant')", //
		"var x = 'elephant'")
	// Whitelisting happens before splitting, not after.
	c.Test(
		"var x = goog.getCssName('goog-elephant')", //
		"var x = 'g-e'")
}

// Java: ReplaceCssNamesTest.testWhitelistByWhole()
func (c *ReplaceCssNamesTest) TestWhitelistByWhole() {
	c.skiplist = jutil.NewLinkedHashSet[string]("long-prefix")
	c.renamingMap = c.GetFullMap()
	c.Test(
		"var x = goog.getCssName('long-prefix')", //
		"var x = 'long-prefix'")
}

// Java: ReplaceCssNamesTest.testWhitelistWithDashes()
func (c *ReplaceCssNamesTest) TestWhitelistWithDashes() {
	c.skiplist = jutil.NewLinkedHashSet[string]("goog-elephant")
	c.Test(
		"var x = goog.getCssName('goog')", //
		"var x = 'g'")
	c.Test(
		"var x = goog.getCssName('elephant')", //
		"var x = 'e'")
	c.Test(
		"var x = goog.getCssName('goog-elephant')", //
		"var x = 'goog-elephant'")
}

// Java: ReplaceCssNamesTest.testVariableReferencesToCssClasses()
func (c *ReplaceCssNamesTest) TestVariableReferencesToCssClasses() {
	cssVarsDefinition :=
		SourceFile_FromCode(
			"foo/styles.css.closure.js",
			`/**
 * @fileoverview generated from foo/styles.css
 * @sassGeneratedCssTs
 */
goog.module('foo.styles$2ecss');
// These files will be automatically generated, so we know the
// exported 'classes' property will always be declared like this.
/** @type {{Bar: string, Baz: string}} */
exports.classes = {
  'Bar': goog.getCssName('fooStylesBar'),
  'Baz': goog.getCssName('fooStylesBaz'),
}
`)
	cssVarsExpected :=
		SourceFile_FromCode(
			"foo/styles.css.closure.js",
			`/**
 * @fileoverview generated from foo/styles.css
 * @sassGeneratedCssTs
 */
goog.module('foo.styles$2ecss');
/** @type {{Bar: string, Baz: string}} */
exports.classes = {
  'Bar': 'fsr',
  'Baz': 'fsz',
}
`)
	importer :=
		SourceFile_FromCode(
			"foo/importer.closure.js",
			"goog.module('foo.importer');\n/** @type {string} */\nconst foo_styles_css = goog.require('foo.styles$2ecss')\n// Even when the original TS import was\n// `import {classes as foo} from './path/to/file';`\n// tsickle will convert references to `foo` into a fully qualified\n// name like this rather than creating an alias variable named `foo`.\nvar x = foo_styles_css.classes.Bar\n")
	c.Test(c.Srcs(cssVarsDefinition, importer), c.Expected(cssVarsExpected, importer))
	AssertThat(c.cssNames).ContainsExactly("fooStylesBar")
}

// Java: ReplaceCssNamesTest.testVariableReferencesToXidCssClasses()
func (c *ReplaceCssNamesTest) TestVariableReferencesToXidCssClasses() {
	cssVarsDefinition :=
		SourceFile_FromCode(
			"foo/styles.css.closure.js",
			`/**
 * @fileoverview generated from foo/styles.css
 * @sassGeneratedCssTs
 */
goog.module('foo.styles$2ecss');
// These files will be automatically generated, so we know the
// exported 'classes' property will always be declared like this.
/** @type {{Bar: string, Baz: string}} */
exports.classes = {
  'Bar': obfuscate ? xid('fooStylesBar') : goog.getCssName('fooStylesBar'),
  'Baz': obfuscate ? xid('fooStylesBaz') : goog.getCssName('fooStylesBaz'),
}
`)
	cssVarsExpected :=
		SourceFile_FromCode(
			"foo/styles.css.closure.js",
			`/**
 * @fileoverview generated from foo/styles.css
 * @sassGeneratedCssTs
 */
goog.module('foo.styles$2ecss');
/** @type {{Bar: string, Baz: string}} */
exports.classes = {
  'Bar': obfuscate ? xid('fooStylesBar') : 'fsr',
  'Baz': obfuscate ? xid('fooStylesBaz') : 'fsz',
}
`)
	importer :=
		SourceFile_FromCode(
			"foo/importer.closure.js",
			"goog.module('foo.importer');\n/** @type {string} */\nconst foo_styles_css = goog.require('foo.styles$2ecss')\n// Even when the original TS import was\n// `import {classes as foo} from './path/to/file';`\n// tsickle will convert references to `foo` into a fully qualified\n// name like this rather than creating an alias variable named `foo`.\nvar x = foo_styles_css.classes.Bar\n")
	c.Test(c.Srcs(cssVarsDefinition, importer), c.Expected(cssVarsExpected, importer))
	AssertThat(c.cssNames).ContainsExactly("fooStylesBar")
}

// Java: ReplaceCssNamesTest.testMixedVariableAndStringReferences()
func (c *ReplaceCssNamesTest) TestMixedVariableAndStringReferences() {
	cssVarsDefinition :=
		SourceFile_FromCode(
			"foo/styles.css.closure.js",
			`/**
 * @fileoverview generated from foo/styles.css
 * @sassGeneratedCssTs
 */
goog.module('foo.styles$2ecss');
/** @type {{Bar: string, Baz: string}} */
exports.classes = {
  'Bar': goog.getCssName('fooStylesBar'),
  'Baz': goog.getCssName('fooStylesBaz'),
}
`)
	cssVarsExpected :=
		SourceFile_FromCode(
			"foo/styles.css.closure.js",
			`/**
 * @fileoverview generated from foo/styles.css
 * @sassGeneratedCssTs
 */
goog.module('foo.styles$2ecss');
/** @type {{Bar: string, Baz: string}} */
exports.classes = {
  'Bar': 'fsr',
  'Baz': 'fsz',
}
`)
	importer :=
		SourceFile_FromCode(
			"foo/importer.closure.js",
			`goog.module('foo.importer');
/** @type {string} */
const foo_styles_css = goog.require('foo.styles$2ecss');
var x = foo_styles_css.classes.Bar;
var y = goog.getCssName('active');
`)
	importerExpected :=
		SourceFile_FromCode(
			"foo/importer.closure.js",
			`goog.module('foo.importer');
/** @type {string} */
const foo_styles_css = goog.require('foo.styles$2ecss');
var x = foo_styles_css.classes.Bar;
var y = 'a';
`)
	c.Test(c.Srcs(cssVarsDefinition, importer), c.Expected(cssVarsExpected, importerExpected))
	AssertThat(c.cssNames).ContainsExactly("fooStylesBar", "active")
}

// Java: ReplaceCssNamesTest.testIgnoreReferencesInCssTsFiles()
func (c *ReplaceCssNamesTest) TestIgnoreReferencesInCssTsFiles() {
	cssVarsDefinition :=
		SourceFile_FromCode(
			"foo/styles.css.closure.js",
			`/**
 * @fileoverview generated from foo/styles.css
 * @sassGeneratedCssTs
 */
goog.module('foo.styles$2ecss');
/** @type {{Bar: string, Baz: string}} */
exports.classes = {
  'Bar': goog.getCssName('fooStylesBar'),
  'Baz': goog.getCssName('fooStylesBaz'),
}
`)
	cssVarsExpected :=
		SourceFile_FromCode(
			"foo/styles.css.closure.js",
			`/**
 * @fileoverview generated from foo/styles.css
 * @sassGeneratedCssTs
 */
goog.module('foo.styles$2ecss');
/** @type {{Bar: string, Baz: string}} */
exports.classes = {
  'Bar': 'fsr',
  'Baz': 'fsz',
}
`)
	c.Test(c.Srcs(cssVarsDefinition), c.Expected(cssVarsExpected))
	AssertThat(c.cssNames).IsEmpty()
}

// Java: ReplaceCssNamesTest.testUnexpectedSassGeneratedCssTsError()
func (c *ReplaceCssNamesTest) TestUnexpectedSassGeneratedCssTsError() {
	nonCssClosureFile :=
		SourceFile_FromCode(
			"foo/invalid.closure.js",
			`/**
 * @fileoverview contains an invalid sassGeneratedCssTs annotation
// This annotation is invalid since this is not a .css.closure.js file
 * @sassGeneratedCssTs
 */
var x = 'foo';
`)
	c.Test(c.Srcs(nonCssClosureFile), c.Error(ReplaceCssNames_UNEXPECTED_SASS_GENERATED_CSS_TS_ERROR))
}

// Java: ReplaceCssNamesTest.testInvalidUseOfClassesObjectError()
func (c *ReplaceCssNamesTest) TestInvalidUseOfClassesObjectError() {
	cssVarsDefinition :=
		SourceFile_FromCode(
			"foo/styles.css.closure.js",
			`/**
 * @fileoverview generated from foo/styles.css
 * @sassGeneratedCssTs
 */
goog.module('foo.styles$2ecss');
/** @type {{Bar: string, Baz: string}} */
exports.classes = {
  'Bar': goog.getCssName('fooStylesBar'),
  'Baz': goog.getCssName('fooStylesBaz'),
}
`)
	importer :=
		SourceFile_FromCode(
			"foo/importer.closure.js",
			`goog.module('foo.importer');
/** @type {string} */
const foo_styles_css = goog.require('foo.styles$2ecss')
var x = foo_styles_css.classes;
`)
	c.Test(c.Srcs(cssVarsDefinition, importer), c.Error(ReplaceCssNames_INVALID_USE_OF_CLASSES_OBJECT_ERROR))
}

// Java: ReplaceCssNamesTest.testUnknownSymbolError()
func (c *ReplaceCssNamesTest) TestUnknownSymbolError() {
	cssVarsDefinition :=
		SourceFile_FromCode(
			"foo/styles.css.closure.js",
			`/**
 * @fileoverview generated from foo/styles.css
 * @sassGeneratedCssTs
 */
goog.module('foo.styles$2ecss');
/** @type {{Bar: string, Baz: string}} */
exports.classes = {
  'Bar': goog.getCssName('fooStylesBar'),
  'Baz': goog.getCssName('fooStylesBaz'),
}
`)
	importer :=
		SourceFile_FromCode(
			"foo/importer.closure.js",
			`goog.module('foo.importer');
/** @type {string} */
const foo_styles_css = goog.require('foo.styles$2ecss')
var x = foo_styles_css.classes.Quux;
`)
	c.Test(c.Srcs(cssVarsDefinition, importer), c.Error(ReplaceCssNames_UNKNOWN_SYMBOL_ERROR))
}
