// Port of com/google/javascript/jscomp/serialization/SerializeAndDeserializeAstTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/serialization/SerializeAndDeserializeAstTest.java)

/*
 * Copyright 2021 The Closure Compiler Authors.
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
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// SerializeAndDeserializeAstTest holds tests that both serialize and then
// deserialize a compiler AST.
//
// Due to the difference from a normal compiler pass, this is not actually
// able to reuse much of the infrastructure inherited from CompilerTestCase,
// and thus it may make sense to separate these tests more fully.
//
// The Java @Test methods not named test* get a "Test" prefix in Go so that
// RunJUnit4 runs them.
//
// Java: public final class SerializeAndDeserializeAstTest extends CompilerTestCase
type SerializeAndDeserializeAstTest struct {
	CompilerTestCaseBase

	consumer                    func(*TypedAst)
	includeTypes                bool
	resolveSourceMapAnnotations bool
	parseInlineSourceMaps       bool
	checkJsDocEquality          bool
	runtimeLibraries            []string
}

// Java: SerializeAndDeserializeAstTest()
func NewSerializeAndDeserializeAstTest() *SerializeAndDeserializeAstTest {
	c := &SerializeAndDeserializeAstTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestSerializeAndDeserializeAstTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewSerializeAndDeserializeAstTest() })
}

// Java: SerializeAndDeserializeAstTest.getProcessor(Compiler)
func (c *SerializeAndDeserializeAstTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewSerializeTypedAstPass(
		compiler,
		c.consumer,
		SerializationOptions_Builder_().
			SetIncludeDebugInfo(false).
			SetRuntimeLibraries(c.runtimeLibraries).
			Build())
}

// Java: SerializeAndDeserializeAstTest.setUp()
func (c *SerializeAndDeserializeAstTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableTypeCheck()
	c.EnableCreateModuleMap()
	c.EnableSourceInformationAnnotator()
	c.includeTypes = true
	c.checkJsDocEquality = true
	c.resolveSourceMapAnnotations = true
	c.parseInlineSourceMaps = true
	c.runtimeLibraries = []string{}
}

// Java: SerializeAndDeserializeAstTest.testConstNumberDeclaration()
func (c *SerializeAndDeserializeAstTest) TestConstNumberDeclaration() {
	c.TestSame("const x = 7;")
}

// Java: SerializeAndDeserializeAstTest.testConstBigIntDeclaration()
func (c *SerializeAndDeserializeAstTest) TestConstBigIntDeclaration() {
	c.TestSame("const x = 7n;")
}

// Java: SerializeAndDeserializeAstTest.testConstStringDeclaration()
func (c *SerializeAndDeserializeAstTest) TestConstStringDeclaration() {
	c.TestSame("const x = 'x';")
}

// Java: SerializeAndDeserializeAstTest.testConstRegexDeclaration()
func (c *SerializeAndDeserializeAstTest) TestConstRegexDeclaration() {
	c.TestSame("const x = /regexp/;")
}

// Java: SerializeAndDeserializeAstTest.testConstObjectDeclaration()
func (c *SerializeAndDeserializeAstTest) TestConstObjectDeclaration() {
	c.TestSame("const obj = {x: 7};")
}

// Java: SerializeAndDeserializeAstTest.testObjectWithShorthandProperty()
func (c *SerializeAndDeserializeAstTest) TestObjectWithShorthandProperty() {
	c.TestSame("const z = 0; const obj = {z};")
}

// Java: SerializeAndDeserializeAstTest.testObjectWithMethod()
func (c *SerializeAndDeserializeAstTest) TestObjectWithMethod() {
	c.TestSame("let obj = {method() {}};")
}

// Java: SerializeAndDeserializeAstTest.testObjectWithQuotedMethod()
func (c *SerializeAndDeserializeAstTest) TestObjectWithQuotedMethod() {
	c.TestSame("let obj = {'method'() {}};")
}

// Java: SerializeAndDeserializeAstTest.testObjectWithGetter()
func (c *SerializeAndDeserializeAstTest) TestObjectWithGetter() {
	c.TestSame("let obj = {get x() {}};")
}

// Java: SerializeAndDeserializeAstTest.testObjectWithQuotedGetter()
func (c *SerializeAndDeserializeAstTest) TestObjectWithQuotedGetter() {
	c.TestSame("let obj = {get 'x'() {}};")
}

// Java: SerializeAndDeserializeAstTest.testObjectWithComputedGetter()
func (c *SerializeAndDeserializeAstTest) TestObjectWithComputedGetter() {
	c.TestSame("let obj = {get [x]() {}};")
}

// Java: SerializeAndDeserializeAstTest.testObjectWithSetter()
func (c *SerializeAndDeserializeAstTest) TestObjectWithSetter() {
	c.TestSame("let obj = {set x(value) {}};")
}

// Java: SerializeAndDeserializeAstTest.testObjectWithQuotedSetter()
func (c *SerializeAndDeserializeAstTest) TestObjectWithQuotedSetter() {
	c.TestSame("let obj = {set 'x'(value) {}};")
}

// Java: SerializeAndDeserializeAstTest.testObjectWithComputedSetter()
func (c *SerializeAndDeserializeAstTest) TestObjectWithComputedSetter() {
	c.TestSame("let obj = {set [x](value) {}};")
}

// Java: SerializeAndDeserializeAstTest.testTemplateLiteral_simple()
func (c *SerializeAndDeserializeAstTest) TestTemplateLiteral_simple() {
	c.TestSame("let obj = `foobar`;")
}

// Java: SerializeAndDeserializeAstTest.testTemplateLiteral_withSubstitution()
func (c *SerializeAndDeserializeAstTest) TestTemplateLiteral_withSubstitution() {
	c.TestSame("let obj = `Hello ${2+3}`;")
}

// Java: SerializeAndDeserializeAstTest.testTemplateLiteral_withIllegalEscape()
func (c *SerializeAndDeserializeAstTest) TestTemplateLiteral_withIllegalEscape() {
	// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Template_literals#es2018_revision_of_illegal_escape_sequences
	c.TestSame("latex`\\unicode`;")
}

// Java: SerializeAndDeserializeAstTest.testConstRegexpDeclaration()
func (c *SerializeAndDeserializeAstTest) TestConstRegexpDeclaration() {
	c.TestSame("const x = /hello world/;")
}

// Java: SerializeAndDeserializeAstTest.testAwait()
func (c *SerializeAndDeserializeAstTest) TestAwait() {
	c.TestSame("async function f() { await 1; }")
}

// Java: SerializeAndDeserializeAstTest.testAwait_topLevel()
func (c *SerializeAndDeserializeAstTest) TestAwait_topLevel() {
	c.TestSame("export const x = await 1;")

	c.checkJsDocEquality = false // @unrestricted
	c.TestSame("export const x = 1; /** @unrestricted */ class C { [await 1] = 2; }")
}

// Java: SerializeAndDeserializeAstTest.testVanillaForLoop()
func (c *SerializeAndDeserializeAstTest) TestVanillaForLoop() {
	c.TestSame("for (let x = 0; x < 10; ++x);")
}

// Java: SerializeAndDeserializeAstTest.testForOfLoop()
func (c *SerializeAndDeserializeAstTest) TestForOfLoop() {
	c.TestSame("for (let elem of []);")
}

// Java: SerializeAndDeserializeAstTest.testForAwaitOfLoop()
func (c *SerializeAndDeserializeAstTest) TestForAwaitOfLoop() {
	c.TestSame("async function f() { for await (let elem of []) {} }")
}

// Java: SerializeAndDeserializeAstTest.testForAwaitOfLoop_topLevel()
func (c *SerializeAndDeserializeAstTest) TestForAwaitOfLoop_topLevel() {
	c.TestSame("export const x = 1; for await (let elem of []) {}")
}

// Java: SerializeAndDeserializeAstTest.testConstructorJsdoc()
func (c *SerializeAndDeserializeAstTest) TestConstructorJsdoc() {
	c.TestSame("/** @constructor */ function Foo() {}")
}

// Java: SerializeAndDeserializeAstTest.testSideEffectsJsdoc()
func (c *SerializeAndDeserializeAstTest) TestSideEffectsJsdoc() {
	c.TestSame("/** @nosideeffects */ function f(x) {}")
	c.TestSame("/** @modifies {arguments} */ function f(x) {}")
	c.TestSame("/** @modifies {this} */ function f(x) {}")
}

// Java: SerializeAndDeserializeAstTest.testCollapsePropertiesJsdoc()
func (c *SerializeAndDeserializeAstTest) TestCollapsePropertiesJsdoc() {
	c.TestSame("const ns = {}; /** @const */ ns.f = (x) => x;")
	c.TestSame("const ns = {}; /** @nocollapse */ ns.f = (x) => x;")
	// Replace types with colors before serializing for this case.
	// 1. It's good to have a test case that covers serializing an AST that already has Colors
	// 2. The test plumbing in this class expects that the AST that was serialized will match the
	//    AST that was deserialized and both will match the expected code.
	c.ReplaceTypesWithColors()
	c.Test(
		"/** @enum {string} */ const Enum = { A: 'string' };",
		"/** @enum {!JSDocSerializer_placeholder_type} */ const Enum = { A: 'string' };")
}

// Java: SerializeAndDeserializeAstTest.testEmptyClassDeclaration()
func (c *SerializeAndDeserializeAstTest) TestEmptyClassDeclaration() {
	c.TestSame("class Foo {}")
}

// Java: SerializeAndDeserializeAstTest.testEmptyClassDeclarationWithExtends()
func (c *SerializeAndDeserializeAstTest) TestEmptyClassDeclarationWithExtends() {
	c.TestSame("class Foo {} class Bar extends Foo {}")
}

// Java: SerializeAndDeserializeAstTest.testClassDeclarationWithMethods()
func (c *SerializeAndDeserializeAstTest) TestClassDeclarationWithMethods() {
	c.TestSame(
		"class Foo {\n" +
			"  a() {}\n" +
			"  get c() {}\n" +
			"  set d(x) {}\n" +
			"}\n")

	c.checkJsDocEquality = false // @unrestricted
	c.TestSame(
		"/** @unrestricted */\n" +
			"class Foo {\n" +
			"  a() {}\n" +
			"  'b'() {}\n" +
			"  get c() {}\n" +
			"  set d(x) {}\n" +
			"  ['e']() {}\n" +
			"}\n")
}

// Java: SerializeAndDeserializeAstTest.testClassDeclarationWithComputedGetter()
func (c *SerializeAndDeserializeAstTest) TestClassDeclarationWithComputedGetter() {
	c.checkJsDocEquality = false // @dict
	c.TestSame(
		"/** @dict */\n" +
			"class Foo {\n" +
			"  get ['a']() {}\n" +
			"}\n")
}

// Java: SerializeAndDeserializeAstTest.testClassDeclarationWithComputedSetter()
func (c *SerializeAndDeserializeAstTest) TestClassDeclarationWithComputedSetter() {
	c.checkJsDocEquality = false // @dict
	c.TestSame(
		"/** @dict */\n" +
			"class Foo {\n" +
			"  set ['a'](val) {}\n" +
			"}\n")
}

// Java: SerializeAndDeserializeAstTest.testClassDeclarationWithFields()
func (c *SerializeAndDeserializeAstTest) TestClassDeclarationWithFields() {
	c.TestSame(
		"class Foo {\n" +
			"  a = 1;\n" +
			"  d;\n" +
			"}\n")

	// Type checking will report computed property accesses as errors for a class,
	// so disable it for this case which contains several.
	c.DisableTypeCheck()
	c.TestSame(
		"class Foo {\n" +
			"  a = 1;\n" +
			"  'b' = 4;\n" +
			"  static ['c'] = 'hi';\n" +
			"  d;\n" +
			"  ['e'];\n" +
			"  1 = 2;\n" +
			"}\n")
}

// Java: SerializeAndDeserializeAstTest.testEmptyClassStaticBlock()
func (c *SerializeAndDeserializeAstTest) TestEmptyClassStaticBlock() {
	c.TestSame(
		"class Foo {\n" +
			"  static {\n" +
			"  }\n" +
			"}\n")
}

// Java: SerializeAndDeserializeAstTest.testClassStaticBlock_variables()
func (c *SerializeAndDeserializeAstTest) TestClassStaticBlock_variables() {
	c.TestSame(
		"class Foo {\n" +
			"  static {\n" +
			"    this.x=1;\n" +
			"    let y =2;\n" +
			"    var z =3;\n" +
			"  }\n" +
			"}\n")
}

// Java: SerializeAndDeserializeAstTest.testClassStaticBlock_function()
func (c *SerializeAndDeserializeAstTest) TestClassStaticBlock_function() {
	c.TestSame(
		"class Foo {\n" +
			"  static {\n" +
			"    function x() {\n" +
			"    }\n" +
			"  }\n" +
			"}\n")
}

// Java: SerializeAndDeserializeAstTest.testMultipleClassStaticBlocks()
func (c *SerializeAndDeserializeAstTest) TestMultipleClassStaticBlocks() {
	c.TestSame(
		"class Foo {\n" +
			"  static {\n" +
			"    this.x=1;\n" +
			"  }\n" +
			"  static {\n" +
			"    this.y=2;\n" +
			"  }\n" +
			"}\n")
}

// Java: SerializeAndDeserializeAstTest.testVanillaFunctionDeclaration()
func (c *SerializeAndDeserializeAstTest) TestVanillaFunctionDeclaration() {
	c.TestSame("function f(x, y) { return x ** y; }")
}

// Java: SerializeAndDeserializeAstTest.testBlockScopedFunctionDeclaration()
func (c *SerializeAndDeserializeAstTest) TestBlockScopedFunctionDeclaration() {
	c.TestSame("if (true) { function f() {} }")
}

// Java: SerializeAndDeserializeAstTest.testAsyncVanillaFunctionDeclaration()
func (c *SerializeAndDeserializeAstTest) TestAsyncVanillaFunctionDeclaration() {
	c.TestSame("async function f() {}")
}

// Java: SerializeAndDeserializeAstTest.testArrowFunctionDeclaration()
func (c *SerializeAndDeserializeAstTest) TestArrowFunctionDeclaration() {
	c.TestSame("let fn = (x) => x >>> 0x80;")
}

// Java: SerializeAndDeserializeAstTest.testAsyncFunctionDeclaration()
func (c *SerializeAndDeserializeAstTest) TestAsyncFunctionDeclaration() {
	c.TestSame("async function f(x) {};")
}

// Java: SerializeAndDeserializeAstTest.testGeneratorFunctionDeclaration()
func (c *SerializeAndDeserializeAstTest) TestGeneratorFunctionDeclaration() {
	c.TestSame("function* f(x) {};")
}

// Java: SerializeAndDeserializeAstTest.testAsyncGeneratorFunctionDeclaration()
func (c *SerializeAndDeserializeAstTest) TestAsyncGeneratorFunctionDeclaration() {
	c.TestSame("async function* f(x) {};")
}

// Java: SerializeAndDeserializeAstTest.testYieldAll()
func (c *SerializeAndDeserializeAstTest) TestYieldAll() {
	c.TestSame("function* f(x) { yield* x; };")
}

// Java: SerializeAndDeserializeAstTest.testFunctionCallRestAndSpread()
func (c *SerializeAndDeserializeAstTest) TestFunctionCallRestAndSpread() {
	c.TestSame("function f(...x) {} f(...[1, 2, 3]);")
}

// Java: SerializeAndDeserializeAstTest.testFunctionDefaultAndDestructuringParameters()
func (c *SerializeAndDeserializeAstTest) TestFunctionDefaultAndDestructuringParameters() {
	c.TestSame("function f([a, b], x = 0, {y, ...z} = {y: 1}) {}")
}

// Java: SerializeAndDeserializeAstTest.testTrailingComma()
func (c *SerializeAndDeserializeAstTest) TestTrailingComma() {
	c.TestSame("function f(x, y,) {}")
	c.TestSame("function f(x) {} f(1,);")
	c.TestSame("class C {constructor(x) {}} new C(1,);")
}

// Java: SerializeAndDeserializeAstTest.testComputedProperty()
func (c *SerializeAndDeserializeAstTest) TestComputedProperty() {
	c.TestSame("const o = {['foo']: 33};")
}

// Java: SerializeAndDeserializeAstTest.testLabledStatement()
func (c *SerializeAndDeserializeAstTest) TestLabledStatement() {
	c.TestSame("label: for (;;);")
}

// Java: SerializeAndDeserializeAstTest.testIdGenerator()
func (c *SerializeAndDeserializeAstTest) TestIdGenerator() {
	c.TestSame("/** @idGenerator {xid} */ function xid(id) {}")
}

// Java: SerializeAndDeserializeAstTest.testUnpairedSurrogateStrings()
func (c *SerializeAndDeserializeAstTest) TestUnpairedSurrogateStrings() {
	c.TestSame("const s = '\xed\xa0\x80';")
}

// Java: SerializeAndDeserializeAstTest.testLogicalAssignment()
func (c *SerializeAndDeserializeAstTest) TestLogicalAssignment() {
	c.TestSame("let x = 0; x ||= 1")
	c.TestSame("function f(x, y) { return x &&= y; }")
	c.TestSame("let y = null; if (y ??= 2) { function f() {} }")
}

// Java: SerializeAndDeserializeAstTest.testExponentAssignment()
func (c *SerializeAndDeserializeAstTest) TestExponentAssignment() {
	c.TestSame("let x = 2; x **= 0;")
}

// Java: SerializeAndDeserializeAstTest.testEsModule()
func (c *SerializeAndDeserializeAstTest) TestEsModule() {
	c.TestSame(
		[]string{
			"export const x = 0; const y = 1; export {y};",
			"import {x, y as z} from './testcode0'; import * as input0 from './testcode0'",
		})
}

// Java: private static final String SerializeAndDeserializeAstTest.BASE64_PREFIX
const serializeAndDeserializeAstTest_BASE64_PREFIX = "data:application/json;base64,"

// Java: private static final String SerializeAndDeserializeAstTest.ENCODED_SOURCE_MAP
const serializeAndDeserializeAstTest_ENCODED_SOURCE_MAP = "eyJ2ZXJzaW9uIjozLCJmaWxlIjoiZm9vLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiZm9vLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiJBQUFBO0lBR0UsV0FBWSxLQUFhO1FBQ3ZCLElBQUksQ0FBQyxDQUFDLEdBQUcsS0FBSyxDQUFDO0lBQ2pCLENBQUM7SUFDSCxRQUFDO0FBQUQsQ0FBQyxBQU5ELElBTUM7QUFFRCxPQUFPLENBQUMsR0FBRyxDQUFDLElBQUksQ0FBQyxDQUFDLENBQUMsQ0FBQyxDQUFDLENBQUMifQ=="

// Java: SerializeAndDeserializeAstTest.testInlineSourceMaps()
func (c *SerializeAndDeserializeAstTest) TestInlineSourceMaps() {
	sourceMapTestCode :=
		"var X = (function () {\n" +
			"    function X(input) {\n" +
			"        this.y = input;\n" +
			"    }\n" +
			"    return X;\n" +
			"}());\n"
	sourceMappingURLComment := "//# sourceMappingURL=" + serializeAndDeserializeAstTest_BASE64_PREFIX + serializeAndDeserializeAstTest_ENCODED_SOURCE_MAP
	code := sourceMapTestCode + "\n" + sourceMappingURLComment

	result := c.testAndReturnResult(c.Srcs(code), c.Expected(code))
	AssertThat(result.compiler.GetBase64SourceMapContents("testcode")).
		IsEqualTo(serializeAndDeserializeAstTest_ENCODED_SOURCE_MAP)
}

// Java: SerializeAndDeserializeAstTest.testSourceMapsInSeparateMapFiles()
func (c *SerializeAndDeserializeAstTest) TestSourceMapsInSeparateMapFiles() {
	// Sourcemap URLs should be a base64 encoded data url, not the name of a .js.map file.
	// If we see a .js.map file, we will not serialize it.
	sourceMapTestCode :=
		"var X = (function () {\n" +
			"    function X(input) {\n" +
			"        this.y = input;\n" +
			"    }\n" +
			"    return X;\n" +
			"}());\n"
	sourceMappingURL := "foo.js.map"
	sourceMappingURLComment := "//# sourceMappingURL=" + sourceMappingURL
	code := sourceMapTestCode + "\n" + sourceMappingURLComment

	result := c.testAndReturnResult(c.Srcs(code), c.Expected(code))
	AssertThat(result.compiler.GetBase64SourceMapContents("testcode")).IsNull()
}

// Java: SerializeAndDeserializeAstTest.testSourceMapsWithoutResolvingSourceMapAnnotations()
func (c *SerializeAndDeserializeAstTest) TestSourceMapsWithoutResolvingSourceMapAnnotations() {
	c.resolveSourceMapAnnotations = false
	sourceMapTestCode :=
		"var X = (function () {\n" +
			"    function X(input) {\n" +
			"        this.y = input;\n" +
			"    }\n" +
			"    return X;\n" +
			"}());\n"
	sourceMappingURLComment := "//# sourceMappingURL=" + serializeAndDeserializeAstTest_BASE64_PREFIX + serializeAndDeserializeAstTest_ENCODED_SOURCE_MAP
	code := sourceMapTestCode + "\n" + sourceMappingURLComment

	result := c.testAndReturnResult(c.Srcs(code), c.Expected(code))
	// Source map not registered because `c.resolveSourceMapAnnotations = false`
	AssertThat(result.compiler.GetBase64SourceMapContents("testcode")).IsNull()
}

// Java: SerializeAndDeserializeAstTest.testSourceMapsWithoutParsingInlineSourceMaps()
func (c *SerializeAndDeserializeAstTest) TestSourceMapsWithoutParsingInlineSourceMaps() {
	c.parseInlineSourceMaps = false
	sourceMapTestCode :=
		"var X = (function () {\n" +
			"    function X(input) {\n" +
			"        this.y = input;\n" +
			"    }\n" +
			"    return X;\n" +
			"}());\n"
	sourceMappingURLComment := "//# sourceMappingURL=" + serializeAndDeserializeAstTest_BASE64_PREFIX + serializeAndDeserializeAstTest_ENCODED_SOURCE_MAP
	code := sourceMapTestCode + "\n" + sourceMappingURLComment

	result := c.testAndReturnResult(c.Srcs(code), c.Expected(code))
	// Source map is registered when `c.parseInlineSourceMaps = false`, but we won't try to
	// parse it as a Base64 encoded source map.
	AssertThat(result.compiler.GetBase64SourceMapContents("testcode")).IsNull()
}

// Java: SerializeAndDeserializeAstTest.testConfiguredDirectorySourceMaps()
func (c *SerializeAndDeserializeAstTest) TestConfiguredDirectorySourceMaps() {
	// We do not allow the TypeScript compiler to set "compilerOptions.sourceRoot" (option to
	// configure a directory to store
	// sourcemaps). Sourcemap URLs should be a base64 encoded data url, not a path to the
	// sourcemap file. If we see a path, we will not serialize anything.
	sourceMapTestCode :=
		"var X = (function () {\n" +
			"    function X(input) {\n" +
			"        this.y = input;\n" +
			"    }\n" +
			"    return X;\n" +
			"}());\n"
	sourceMappingURLPath := "directory/foo.js.map"
	sourceMappingURLComment := "//# sourceMappingURL=" + sourceMappingURLPath
	code := sourceMapTestCode + "\n" + sourceMappingURLComment

	result := c.testAndReturnResult(c.Srcs(code), c.Expected(code))
	AssertThat(result.compiler.GetBase64SourceMapContents("testcode")).IsNull()
}

// Java: SerializeAndDeserializeAstTest.testConvertsNumberTypeToColor()
func (c *SerializeAndDeserializeAstTest) TestConvertsNumberTypeToColor() {
	c.EnableTypeCheck()

	result := c.testAndReturnResult(c.Srcs("3"), c.Expected("3"))
	newScript := result.sourceRoot.GetFirstChild()
	AssertNode(newScript).HasToken(rhino.Token_SCRIPT)
	three := newScript.GetFirstFirstChild()

	AssertNode(three).HasToken(rhino.Token_NUMBER)
	ColorSubject_AssertThat(three.GetColor()).IsSameInstanceAs(colors.StandardColors_NUMBER)
}

// Java: SerializeAndDeserializeAstTest.testConvertsArrayTypeToColor_andStoresInRegistry()
func (c *SerializeAndDeserializeAstTest) TestConvertsArrayTypeToColor_andStoresInRegistry() {
	c.EnableTypeCheck()

	result := c.testAndReturnResult(c.Srcs("[]"), c.Expected("[]"))
	newScript := result.sourceRoot.GetFirstChild()
	AssertNode(newScript).HasToken(rhino.Token_SCRIPT)
	three := newScript.GetFirstFirstChild()

	AssertNode(three).HasToken(rhino.Token_ARRAYLIT)
	ColorSubject_AssertThat(three.GetColor()).IsSameInstanceAs(result.registry.Get(colors.StandardColors_ARRAY_ID))
}

// Java: SerializeAndDeserializeAstTest.testSkipAddingColors()
func (c *SerializeAndDeserializeAstTest) TestSkipAddingColors() {
	c.EnableTypeCheck()
	// the pre-serialization AST includes JSTypes, but the deserialized AST has no Colors
	c.includeTypes = false

	result := c.testAndReturnResult(c.Srcs("3"), c.Expected("3"))
	newScript := result.sourceRoot.GetFirstChild()
	AssertNode(newScript).HasToken(rhino.Token_SCRIPT)
	three := newScript.GetFirstFirstChild()

	AssertNode(three).HasToken(rhino.Token_NUMBER)
	ColorSubject_AssertThat(three.GetColor()).IsNull()
	AssertThat(result.registry).IsNull()
}

// Java: SerializeAndDeserializeAstTest.testOriginalNamePreserved()
func (c *SerializeAndDeserializeAstTest) TestOriginalNamePreserved() {
	newRoot :=
		c.testAndReturnResult(c.Srcs("const x = 0;"), c.Expected("const x = 0;")).
			sourceRoot.
			GetFirstChild()

	constDeclaration := newRoot.GetFirstChild()
	AssertNode(constDeclaration).HasToken(rhino.Token_CONST)

	x := constDeclaration.GetOnlyChild()
	AssertNode(x).HasStringThat().IsEqualTo("x")
	AssertNode(x).HasOriginalNameThat().IsEqualTo("x")
}

// Java: SerializeAndDeserializeAstTest.testOriginalNamePreservedAfterModuleRewriting()
func (c *SerializeAndDeserializeAstTest) TestOriginalNamePreservedAfterModuleRewriting() {
	c.EnableRewriteClosureCode()

	newRoot :=
		c.testAndReturnResultWithExterns(
			c.Externs(NewTestExternsBuilder().AddClosureExterns().Build()),
			c.Srcs("goog.module('a.b.c'); const x = 0;"),
			c.Expected(
				"/** @const */ var module$exports$a$b$c = {};\n"+
					"const module$contents$a$b$c_x = 0;\n")).
			sourceRoot.
			GetFirstChild()

	constDeclaration := newRoot.GetSecondChild()
	AssertNode(constDeclaration).HasToken(rhino.Token_CONST)

	globalizedXName := constDeclaration.GetOnlyChild()
	AssertNode(globalizedXName).HasStringThat().IsEqualTo("module$contents$a$b$c_x")
	AssertNode(globalizedXName).HasOriginalNameThat().IsEqualTo("x")
}

// Java: SerializeAndDeserializeAstTest.serializesFileWithPreloadedCode()
func (c *SerializeAndDeserializeAstTest) TestSerializesFileWithPreloadedCode() {
	a := SourceFile_FromCode("a.js", "const a = 0;")
	b := SourceFile_FromCode("b.js", "const b = a;")

	result :=
		c.testAndReturnResult(c.Srcs([]*SourceFile{a, b}), c.Expected([]*SourceFile{a, b}))
	scriptA := result.sourceRoot.GetFirstChild()
	scriptB := result.sourceRoot.GetSecondChild()

	IsInstanceOf[*SourceFile](AssertThat(scriptA.GetStaticSourceFile()))
	IsInstanceOf[*SourceFile](AssertThat(scriptB.GetStaticSourceFile()))

	AssertThat(scriptA.GetStaticSourceFile().(*SourceFile).GetCode()).IsEqualTo("const a = 0;")
	AssertThat(scriptB.GetStaticSourceFile().(*SourceFile).GetCode()).IsEqualTo("const b = a;")
}

// Java: SerializeAndDeserializeAstTest.serializeAndDeserializeFileOnDiskWithUTF16()
func (c *SerializeAndDeserializeAstTest) TestSerializeAndDeserializeFileOnDiskWithUTF16() {
	pathA := serializeAndDeserializeAstTest_createTempFile(c.T(), "tmp", "a.js")
	serializeAndDeserializeAstTest_writeLines(c.T(), pathA, []string{"const ಠ_ಠ = 0;"}, jutil.StandardCharsets_UTF_16)

	a := SourceFile_FromFileWithCharset(pathA, jutil.StandardCharsets_UTF_16)

	result :=
		c.testAndReturnResult(c.Srcs([]*SourceFile{a}), c.Expected([]*SourceFile{a}))
	scriptA := result.sourceRoot.GetFirstChild()

	IsInstanceOf[*SourceFile](AssertThat(scriptA.GetStaticSourceFile()))
	AssertThat(scriptA.GetStaticSourceFile().(*SourceFile).GetCode()).
		IsEqualTo("const ಠ_ಠ = 0;\n")
}

// Java: SerializeAndDeserializeAstTest.serializeAndDeserializeFileOnDiskWithOriginalName()
func (c *SerializeAndDeserializeAstTest) TestSerializeAndDeserializeFileOnDiskWithOriginalName() {
	pathA := serializeAndDeserializeAstTest_createTempFile(c.T(), "tmp", "a.js")
	serializeAndDeserializeAstTest_writeLines(c.T(), pathA, []string{"const a = 0;"}, jutil.StandardCharsets_UTF_8)

	a :=
		SourceFile_Builder_().WithOriginalPath("original_a.js").WithPath(pathA).Build()

	result :=
		c.testAndReturnResult(c.Srcs([]*SourceFile{a}), c.Expected([]*SourceFile{a}))
	scriptA := result.sourceRoot.GetFirstChild()

	IsInstanceOf[*SourceFile](AssertThat(scriptA.GetStaticSourceFile()))
	AssertThat(scriptA.GetStaticSourceFile().(*SourceFile).GetCode()).IsEqualTo("const a = 0;\n")
	AssertThat(scriptA.GetSourceFileName()).IsEqualTo("original_a.js")
}

// Java: SerializeAndDeserializeAstTest.serializesZipEntries()
func (c *SerializeAndDeserializeAstTest) TestSerializesZipEntries() {
	// Setup environment.
	expectedContent := "const a = 0;"
	// @Rule TemporaryFolder folder: folder.newFile("test.js.zip").
	jsZipPath := filepath.Join(c.T().TempDir(), "test.js.zip")
	serializeAndDeserializeAstTest_createZipWithContent(c.T(), jsZipPath, expectedContent)

	a := SourceFile_FromFileWithCharset(jsZipPath+"!/a.js", jutil.StandardCharsets_UTF_8)

	result :=
		c.testAndReturnResult(c.Srcs([]*SourceFile{a}), c.Expected([]*SourceFile{a}))
	scriptA := result.sourceRoot.GetFirstChild()

	IsInstanceOf[*SourceFile](AssertThat(scriptA.GetStaticSourceFile()))
	AssertThat(scriptA.GetStaticSourceFile().(*SourceFile).GetCode()).IsEqualTo("const a = 0;")
	AssertThat(scriptA.GetStaticSourceFile().GetName()).IsEqualTo(a.GetName())
}

// Java: SerializeAndDeserializeAstTest.setsSourceFileOfSyntheticCode()
func (c *SerializeAndDeserializeAstTest) TestSetsSourceFileOfSyntheticCode() {
	c.EnsureLibraryInjected("base")

	result :=
		c.testAndReturnResult(
			c.Srcs("0;"),
			// the injected "base" library is merged into the first file's script. ensure that
			// SourceFiles are wired up correctly.
			c.Expected(
				"/** @const */ var $jscomp = $jscomp || {};\n"+
					"/** @const */\n"+
					"$jscomp.scope = {};\n"+
					"0;\n"))

	script := result.sourceRoot.GetFirstChild()
	AssertNode(script).HasToken(rhino.Token_SCRIPT)
	AssertThat(script.GetSourceFileName()).IsEqualTo("testcode")

	jscompDeclaration := script.GetFirstChild()
	AssertThat(jscompDeclaration.GetSourceFileName()).EndsWith("js/base.js")
	AssertThat(jscompDeclaration.GetFirstChild().GetSourceFileName()).EndsWith("js/base.js")

	number := script.GetLastChild()
	AssertThat(number.GetSourceFileName()).IsEqualTo("testcode")
}

// Java: SerializeAndDeserializeAstTest.includesExternsSummary()
func (c *SerializeAndDeserializeAstTest) TestIncludesExternsSummary() {
	c.EnableGatherExternProperties()

	result :=
		c.testAndReturnResultWithExterns(
			c.Externs(
				"class Foo {\n"+
					"// Ensure JSDoc properties are included (b/180424427)\n"+
					"  /** @param {{arg: string}} x */\n"+
					"  method(x) { }\n"+
					"}\n"),
			c.Srcs(""),
			c.Expected(""))

	AssertThat(result.ast.GetExternProperties()).ContainsAtLeast("method", "arg")
}

// Java: SerializeAndDeserializeAstTest.includesShadowedCode()
func (c *SerializeAndDeserializeAstTest) TestIncludesShadowedCode() {
	result :=
		c.runWithoutExpected(
			c.Externs(NewTestExternsBuilder().AddClosureExterns().Build()),
			c.Srcs(
				"/** @closureUnaware @fileoverview */\n"+
					"goog.module('test');\n"+
					"\n"+
					"/** @closureUnaware */\n"+
					"(function() {\n"+
					"  window['foo'] = 5;\n"+
					"})();\n"))

	// Also validate that the content of the shadow is actually correct
	shadowHost :=
		CodeSubTree_FindFirstNode(result.sourceRoot, func(n rhino.Node) bool { return !n.GetClosureUnawareShadow().IsNil() })

	shadowedContent := shadowHost.GetClosureUnawareShadow()

	expectedShadowContent :=
		c.ParseExpectedJs(
			"$jscomp_sink_closure_unaware_impl(function() {\n" +
				"  window['foo'] = 5;\n" +
				"});\n")
	AssertNode(shadowedContent).IsEqualIncludingJsDocTo(expectedShadowContent)

	AssertThat(result.compiler.ToSourceOfNode(result.sourceRoot)).
		IsEqualTo(
			"/** @closureUnaware */ goog.module(\"test\");(function(){window[\"foo\"]=5})()")

	var cb NodeTraversal_Callback = &serializeAndDeserializeAstTest_shadowCallback{}

	NodeTraversal_Traverse(result.compiler, shadowedContent, cb)
}

// Java: new NodeTraversal.AbstractPostOrderCallback() in
// SerializeAndDeserializeAstTest.includesShadowedCode()
type serializeAndDeserializeAstTest_shadowCallback struct {
	NodeTraversal_AbstractPostOrderCallbackBase
}

// Java: SerializeAndDeserializeAstTest.includesShadowedCode().new AbstractPostOrderCallback() {...}.visit(NodeTraversal, Node, Node)
func (cb *serializeAndDeserializeAstTest_shadowCallback) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	AssertThat(n.GetIsInClosureUnawareSubtree()).IsTrue()
	if !parent.IsNil() {
		AssertThat(rhino.Node_ValidateMemorySensitivePropertyGuarantees(n, true, parent, true)).
			IsTrue()
	}
}

// Java: SerializeAndDeserializeAstTest.shadowAsts_featuresSetOnMainAstScript()
func (c *SerializeAndDeserializeAstTest) TestShadowAsts_featuresSetOnMainAstScript() {
	result :=
		c.runWithoutExpected(
			c.Externs(NewTestExternsBuilder().AddClosureExterns().Build()),
			c.Srcs(
				"/** @closureUnaware @fileoverview */\n"+
					"goog.module('test');\n"+
					"\n"+
					"/** @closureUnaware */\n"+
					"(function() {\n"+
					"  const x = 5;\n"+
					"})();\n"))

	shadowHost :=
		CodeSubTree_FindFirstNode(result.sourceRoot, func(n rhino.Node) bool { return !n.GetClosureUnawareShadow().IsNil() })
	shadowScript := shadowHost.GetClosureUnawareShadow().GetOnlyChild()
	mainScript := result.sourceRoot.GetOnlyChild()
	jutil.CheckStateMsg1(mainScript.IsScript(), "%s", mainScript)
	jutil.CheckStateMsg1(shadowScript.IsScript(), "%s", shadowScript)

	// The main AST script contains the "const" feature; the shadow AST featureset is null.
	// NOTE: this behavior was selected for consistency with the parser/IRFactory and to address
	// AstValidator breakages. However, we could consider instead setting the featureset on
	// the shadow script directly, and omitting CONST_DECLARATIONS in the main script featureset.
	// At present, the main advantage of doing so seems to be a potential performance gain because
	// transpilation passes on the main AST can skip any features only in the shadow AST, but we
	// don't expect that performance gain to be significant.
	AssertThat(mainScript.GetProp(rhino.Node_FEATURE_SET)).
		IsEqualTo(parser.FeatureSet_BARE_MINIMUM.With(parser.FeatureSet_Feature_CONST_DECLARATIONS))
	AssertThat(shadowScript.GetProp(rhino.Node_FEATURE_SET)).IsNull()
}

// Java: SerializeAndDeserializeAstTest.includesRuntimeLibraryPaths()
func (c *SerializeAndDeserializeAstTest) TestIncludesRuntimeLibraryPaths() {
	c.EnableGatherExternProperties()
	c.runtimeLibraries = []string{"base", "es6/string"}

	result := c.testAndReturnResultWithExterns(c.Externs(""), c.Srcs(""), c.Expected(""))

	AssertThat(result.ast.GetRuntimeLibraries()).ContainsExactly("base", "es6/string")
}

// Java: SerializeAndDeserializeAstTest.wellKnownSymbolsBecomeSymbol()
func (c *SerializeAndDeserializeAstTest) TestWellKnownSymbolsBecomeSymbol() {
	c.EnableTypeCheck()
	result :=
		c.testAndReturnResultWithExterns(
			c.Externs(
				"/** @type {symbol} */\n"+
					"Symbol.foobar;\n"),
			c.Srcs("Symbol.foobar"),
			c.Expected("Symbol.foobar"))

	newScript := result.sourceRoot.GetFirstChild()
	jutil.CheckStateMsg1(newScript.IsScript(), "%s", newScript)

	foobar := newScript.GetFirstFirstChild()

	AssertNode(foobar).MatchesQualifiedName("Symbol.foobar")
	ColorSubject_AssertThat(foobar.GetColor()).IsSameInstanceAs(colors.StandardColors_SYMBOL)
}

// Java: SerializeAndDeserializeAstTest.testSame(String)
// Java: SerializeAndDeserializeAstTest.testSame(String[])
func (c *SerializeAndDeserializeAstTest) TestSame(args ...any) {
	switch {
	case compilerTestCase_argsAre(args, "string"):
		code := args[0].(string)
		c.Test(code, code)
	case compilerTestCase_argsAre(args, "[]string"):
		sources := args[0].([]string)
		c.testAndReturnResult(c.Srcs(sources), c.Expected(sources))
	default:
		// Go-only: language-adapter — the testSame overloads this class does not override.
		c.CompilerTestCaseBase.TestSame(args...)
	}
}

// Java: SerializeAndDeserializeAstTest.test(String, String)
func (c *SerializeAndDeserializeAstTest) Test(args ...any) {
	if compilerTestCase_argsAre(args, "string", "string") {
		code, expected := args[0].(string), args[1].(string)
		c.testAndReturnResult(c.Srcs(code), c.Expected(expected))
		return
	}
	// Go-only: language-adapter — the test overloads this class does not override.
	c.CompilerTestCaseBase.Test(args...)
}

// Java: SerializeAndDeserializeAstTest.testAndReturnResult(Sources, Expected)
func (c *SerializeAndDeserializeAstTest) testAndReturnResult(code CompilerTestCase_Sources, expected *CompilerTestCase_Expected) *serializeAndDeserializeAstTest_Result {
	return c.testAndReturnResultWithExterns(c.Externs([]*SourceFile{}), code, expected)
}

// Java: SerializeAndDeserializeAstTest.testAndReturnResult(Externs, Sources, Expected)
func (c *SerializeAndDeserializeAstTest) testAndReturnResultWithExterns(externs *CompilerTestCase_Externs, code CompilerTestCase_Sources, expected *CompilerTestCase_Expected) *serializeAndDeserializeAstTest_Result {
	serializedStream := c.toInputStream(externs, code, expected)
	result := c.run(externs, code, serializedStream)

	newSourceRoot := result.sourceRoot
	expectedRoot := c.ParseExpectedJs(expected)
	if c.checkJsDocEquality {
		AssertNode(newSourceRoot).IsEqualIncludingJsDocTo(expectedRoot)
	} else {
		AssertNode(newSourceRoot).IsEqualTo(expectedRoot)
	}
	return result
}

// Java: SerializeAndDeserializeAstTest.runWithoutExpected(Externs, Sources)
func (c *SerializeAndDeserializeAstTest) runWithoutExpected(externs *CompilerTestCase_Externs, code CompilerTestCase_Sources) *serializeAndDeserializeAstTest_Result {
	result := make([]*TypedAst, 1)
	c.consumer = func(ast *TypedAst) { result[0] = ast }
	c.CompilerTestCaseBase.TestNoWarning(externs, code)
	serializedStream := bytes.NewBuffer((&TypedAst_List{typedAsts_: result}).ToByteArray())
	return c.run(externs, code, serializedStream)
}

// Java: SerializeAndDeserializeAstTest.run(Externs, Sources, InputStream)
func (c *SerializeAndDeserializeAstTest) run(externs *CompilerTestCase_Externs, code CompilerTestCase_Sources, serializedStream io.Reader) *serializeAndDeserializeAstTest_Result {
	serializingCompiler := c.GetLastCompiler()

	oldRoot := serializingCompiler.GetRoot()

	externFiles := c.collectSourceFilesFromScripts(oldRoot.GetFirstChild())
	codeFiles := c.collectSourceFilesFromScripts(oldRoot.GetSecondChild())

	// NOTE: We need a fresh compiler instance in which to deserialize, because:
	// 1. This is a better representation of what will happen in production use.
	// 2. Deserializing expects to start with a compiler which has no Color-related state.
	//    For example, both serializing and deserializing currently need to call
	//    `compiler.initRuntimeLibraryTypedAsts()`, which may only be called once.
	deserializingCompiler := c.CreateAndInitializeCompiler(externs, code)
	requiredInputFiles := make([]*SourceFile, 0, len(externFiles)+len(codeFiles))
	requiredInputFiles = append(requiredInputFiles, externFiles...)
	requiredInputFiles = append(requiredInputFiles, codeFiles...)
	ast :=
		TypedAstDeserializer_DeserializeFullAst(
			deserializingCompiler,
			SourceFile_FromCodeWithKind("syntheticExterns", "", rhino.StaticSourceFile_SourceKind_EXTERN),
			requiredInputFiles,
			serializedStream,
			c.includeTypes,
			c.resolveSourceMapAnnotations,
			c.parseInlineSourceMaps)

	registry := ast.GetColorRegistry()
	newExternsRoot := rhino.IR_Root()
	newSourceRoot := rhino.IR_Root()

	// this code is a little clunky, but basically the TypedAstDeserializer API assumes you already
	// know what order the sources/externs are expected to be in. The easiest way for us to get that
	// ordering is by looking at the pre-serialized AST, since the Sources/Expected APIs hide the
	// actual file names.
	for oldExtern := serializingCompiler.GetRoot().GetFirstFirstChild(); !oldExtern.IsNil(); oldExtern = oldExtern.GetNext() {
		extern := oldExtern.GetStaticSourceFile().(*SourceFile)
		script := ast.filesystem[extern].DeserializeNew()
		script.SetInputId(rhino.NewInputId(extern.GetName()))
		newExternsRoot.AddChildToBack(script)
	}
	for oldSource := serializingCompiler.GetRoot().GetSecondChild().GetFirstChild(); !oldSource.IsNil(); oldSource = oldSource.GetNext() {
		source := oldSource.GetStaticSourceFile().(*SourceFile)
		script := ast.filesystem[source].DeserializeNew()
		script.SetInputId(rhino.NewInputId(source.GetName()))
		newSourceRoot.AddChildToBack(script)
	}
	newRoot := rhino.IR_Root(newExternsRoot, newSourceRoot)

	c.assertFeatureSetsOfScriptsMatch(oldRoot, newRoot)

	NewAstValidatorWithValidateScriptFeatures(deserializingCompiler /* validateScriptFeatures= */, true).
		ValidateRoot(newRoot)

	c.consumer = nil

	return newSerializeAndDeserializeAstTest_Result(ast, registry, newSourceRoot, deserializingCompiler)
}

// Java: SerializeAndDeserializeAstTest.assertFeatureSetsOfScriptsMatch(Node, Node)
func (c *SerializeAndDeserializeAstTest) assertFeatureSetsOfScriptsMatch(oldRoot rhino.Node, newRoot rhino.Node) {
	featureSetsOfOldScripts := c.getFeatureSetsOfScripts(oldRoot)
	featureSetsOfNewScripts := c.getFeatureSetsOfScripts(newRoot)
	AssertThat(featureSetsOfNewScripts.KeySlice()).IsEqualTo(featureSetsOfOldScripts.KeySlice())
	for sourceFile := range featureSetsOfOldScripts.Keys() {
		oldFeatureSet := featureSetsOfOldScripts.Get(sourceFile)
		newFeatureSet := featureSetsOfNewScripts.Get(sourceFile)
		jutil.CheckStateMsg3(
			jutil.ObjectsEquals(oldFeatureSet, newFeatureSet),
			"Feature sets of scripts do not match for file %s. Old: %s, New: %s",
			sourceFile.GetName(),
			oldFeatureSet,
			newFeatureSet)
	}
}

// Java: SerializeAndDeserializeAstTest.getFeatureSetsOfScripts(Node)
func (c *SerializeAndDeserializeAstTest) getFeatureSetsOfScripts(root rhino.Node) *jutil.LinkedHashMap[*SourceFile, *parser.FeatureSet] {
	builder := jutil.NewImmutableMapBuilder[*SourceFile, *parser.FeatureSet]()
	externsRoot := root.GetFirstChild()
	sourceRoot := root.GetSecondChild()
	for script := externsRoot.GetFirstChild(); !script.IsNil(); script = script.GetNext() {
		jutil.CheckState(script.IsScript())
		builder.Put(
			script.GetStaticSourceFile().(*SourceFile), serializeAndDeserializeAstTest_featureSet(script))
	}
	for script := sourceRoot.GetFirstChild(); !script.IsNil(); script = script.GetNext() {
		jutil.CheckState(script.IsScript())
		builder.Put(
			script.GetStaticSourceFile().(*SourceFile), serializeAndDeserializeAstTest_featureSet(script))
	}
	return builder.BuildOrThrow()
}

// serializeAndDeserializeAstTest_featureSet is the cast
// (FeatureSet) script.getProp(Node.FEATURE_SET), null included.
//
// Go-only: language-adapter — a Go type assertion panics on a nil interface.
func serializeAndDeserializeAstTest_featureSet(script rhino.Node) *parser.FeatureSet {
	featureSet, _ := script.GetProp(rhino.Node_FEATURE_SET).(*parser.FeatureSet)
	return featureSet
}

// Java: SerializeAndDeserializeAstTest.collectSourceFilesFromScripts(Node)
func (c *SerializeAndDeserializeAstTest) collectSourceFilesFromScripts(root rhino.Node) []*SourceFile {
	files := []*SourceFile{}
	for n := root.GetFirstChild(); !n.IsNil(); n = n.GetNext() {
		files = append(files, n.GetStaticSourceFile().(*SourceFile))
	}
	return files
}

// Java: private static class SerializeAndDeserializeAstTest.Result
type serializeAndDeserializeAstTest_Result struct {
	ast        *TypedAstDeserializer_DeserializedAst
	registry   *colors.ColorRegistry
	sourceRoot rhino.Node
	compiler   *Compiler
}

// Java: SerializeAndDeserializeAstTest.Result(DeserializedAst, ColorRegistry, Node, Compiler)
func newSerializeAndDeserializeAstTest_Result(ast *TypedAstDeserializer_DeserializedAst, registry *colors.ColorRegistry, sourceRoot rhino.Node, compiler *Compiler) *serializeAndDeserializeAstTest_Result {
	return &serializeAndDeserializeAstTest_Result{
		ast:        ast,
		registry:   registry,
		sourceRoot: sourceRoot,
		compiler:   compiler,
	}
}

// Java: SerializeAndDeserializeAstTest.toInputStream(Externs, Sources, Expected)
func (c *SerializeAndDeserializeAstTest) toInputStream(externs *CompilerTestCase_Externs, code CompilerTestCase_Sources, expected *CompilerTestCase_Expected) io.Reader {
	result := make([]*TypedAst, 1)
	c.consumer = func(ast *TypedAst) { result[0] = ast }
	// The process of serializing modifies the input code (e.g. removing casts and changing JSDoc)
	// to the expected result.
	c.CompilerTestCaseBase.Test(externs, code, expected)
	return bytes.NewBuffer((&TypedAst_List{typedAsts_: result}).ToByteArray())
}

// Java: SerializeAndDeserializeAstTest.createZipWithContent(Path, String)
func serializeAndDeserializeAstTest_createZipWithContent(t *testing.T, zipFile string, content string) {
	lastModified := time.Now()
	if info, err := os.Stat(zipFile); err == nil {
		// Ensure that file modified date is updated, otherwise could cause flakiness (b/123962282).
		lastModified = info.ModTime().Add(time.Second)
		os.Remove(zipFile)
	}

	f, err := os.Create(zipFile)
	if err != nil {
		t.Fatal(err)
	}
	zos := zip.NewWriter(f)
	w, err := zos.Create("a.js")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := w.Write(jutil.StringGetBytesUTF8(content)); err != nil {
		t.Fatal(err)
	}
	if err := zos.Close(); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(zipFile, lastModified, lastModified); err != nil {
		t.Fatal(err)
	}
}

// serializeAndDeserializeAstTest_createTempFile is
// Files.createTempFile(prefix, suffix).toString().
//
// Go-only: test-helper — java.nio.file.Files.
func serializeAndDeserializeAstTest_createTempFile(t *testing.T, prefix string, suffix string) string {
	f, err := os.CreateTemp(t.TempDir(), prefix+"*"+suffix)
	if err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	return f.Name()
}

// serializeAndDeserializeAstTest_writeLines is Files.write(path, lines,
// charset): each line followed by the line separator ("\n"), encoded with
// the charset (UTF-8, or UTF-16: a big-endian byte order mark, then
// big-endian units).
//
// Go-only: test-helper — java.nio.file.Files.
func serializeAndDeserializeAstTest_writeLines(t *testing.T, path string, lines []string, charset *jutil.Charset) {
	var text strings.Builder
	for _, line := range lines {
		text.WriteString(line + "\n")
	}
	var data []byte
	switch charset {
	case jutil.StandardCharsets_UTF_8:
		data = jutil.StringGetBytesUTF8(text.String())
	case jutil.StandardCharsets_UTF_16:
		data = []byte{0xFE, 0xFF}
		for _, u := range jutil.StringToUTF16(text.String()) {
			data = append(data, byte(u>>8), byte(u))
		}
	default:
		t.Fatalf("unsupported charset %s", charset)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}
