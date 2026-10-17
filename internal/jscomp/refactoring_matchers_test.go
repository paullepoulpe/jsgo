// Port of com/google/javascript/refactoring/MatchersTest.java @ closure-compiler v20261015
// (test/com/google/javascript/refactoring/MatchersTest.java)

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
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// MatchersTest holds unit tests for JsFlume Matchers.
//
// Java: public class MatchersTest
type MatchersTest struct {
	// Go-only: test-helper — the Go test running the current test method.
	t *testing.T
}

// Java: MatchersTest()
func NewMatchersTest() *MatchersTest {
	return &MatchersTest{}
}

func TestMatchersTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewMatchersTest() })
}

// JUnitSetT gives the instance the Go test running its test method.
func (c *MatchersTest) JUnitSetT(t *testing.T) { c.t = t }

// JUnitBefore: the class has no @Before method.
func (c *MatchersTest) JUnitBefore() {}

// JUnitAfter: the class has no @After method.
func (c *MatchersTest) JUnitAfter() {}

// Java: MatchersTest.testAnything()
func (c *MatchersTest) TestAnything() {
	input := "goog.base(this);"
	root := matchersTest_compileToScriptRoot(matchersTest_getCompiler(input))
	child := root.GetFirstChild()
	AssertThat(Matchers_Anything().Matches(nil, nil)).IsTrue()
	AssertThat(Matchers_Anything().Matches(root, nil)).IsTrue()
	AssertThat(Matchers_Anything().Matches(child, nil)).IsTrue()
	AssertThat(Matchers_Anything().Matches(child, NewNodeMetadata(nil))).IsTrue()
}

// Java: MatchersTest.testAllOf()
func (c *MatchersTest) TestAllOf() {
	input := "goog.require('goog.dom');"
	root := matchersTest_compileToScriptRoot(matchersTest_getCompiler(input))
	fnCall := root.GetFirstFirstChild()
	AssertThat(fnCall.IsCall()).IsTrue()

	notMatcher := Matchers_Not(Matchers_Anything())
	AssertThat(Matchers_AllOf(notMatcher).Matches(root, nil)).IsFalse()
	AssertThat(Matchers_AllOf(notMatcher, Matchers_FunctionCall()).Matches(fnCall, nil)).IsFalse()
	AssertThat(Matchers_AllOf(
		Matchers_Anything(), Matchers_FunctionCallWithName(jutil.StringPtr("goog.require"))).Matches(fnCall, nil)).IsTrue()
}

// Java: MatchersTest.testAnyOf()
func (c *MatchersTest) TestAnyOf() {
	input := "goog.require('goog.dom');"
	root := matchersTest_compileToScriptRoot(matchersTest_getCompiler(input))
	fnCall := root.GetFirstFirstChild()
	AssertThat(fnCall.IsCall()).IsTrue()

	notMatcher := Matchers_Not(Matchers_Anything())
	AssertThat(Matchers_AnyOf(notMatcher).Matches(root, nil)).IsFalse()
	AssertThat(Matchers_AnyOf(notMatcher, Matchers_FunctionCall()).Matches(fnCall, nil)).IsTrue()
	AssertThat(Matchers_AnyOf(Matchers_FunctionCallWithName(jutil.StringPtr("goog.require"))).Matches(fnCall, nil)).IsTrue()
}

// Java: MatchersTest.testNot()
func (c *MatchersTest) TestNot() {
	AssertThat(Matchers_Not(Matchers_Anything()).Matches(nil, nil)).IsFalse()

	input := "goog.require('goog.dom');"
	root := matchersTest_compileToScriptRoot(matchersTest_getCompiler(input))
	fnCall := root.GetFirstFirstChild()
	AssertThat(fnCall.IsCall()).IsTrue()
	AssertThat(Matchers_Not(Matchers_FunctionCall()).Matches(fnCall, nil)).IsFalse()
	AssertThat(Matchers_Not(Matchers_FunctionCallWithName(jutil.StringPtr("goog.require"))).Matches(fnCall, nil)).IsFalse()
}

// Java: MatchersTest.testConstructor_any()
func (c *MatchersTest) TestConstructor_any() {
	input := "/** @constructor */ var Foo = function() {};"
	root := matchersTest_compileToScriptRoot(matchersTest_getCompiler(input))
	AssertThat(Matchers_Constructor().Matches(root.GetFirstChild(), nil)).IsTrue()
}

// Java: MatchersTest.testConstructor_specificClass()
func (c *MatchersTest) TestConstructor_specificClass() {
	input := "/** @constructor */ var Foo = function() {};"
	root := matchersTest_compileToScriptRoot(matchersTest_getCompiler(input))
	AssertThat(Matchers_ConstructorWithName(jutil.StringPtr("Foo")).Matches(root.GetFirstChild(), nil)).IsTrue()
	AssertThat(Matchers_ConstructorWithName(jutil.StringPtr("Bar")).Matches(root.GetFirstChild(), nil)).IsFalse()
}

// Java: MatchersTest.testConstructor_differentConstructorTypes()
func (c *MatchersTest) TestConstructor_differentConstructorTypes() {
	input := "/** @constructor */ var Foo = function() {};"
	root := matchersTest_compileToScriptRoot(matchersTest_getCompiler(input))
	ctorNode := root.GetFirstChild()
	AssertThat(Matchers_ConstructorWithName(jutil.StringPtr("Foo")).Matches(ctorNode, nil)).IsTrue()

	input = "/** @constructor */ bar.Foo = function() {};"
	root = matchersTest_compileToScriptRoot(matchersTest_getCompiler(input))
	ctorNode = root.GetFirstFirstChild()
	AssertThat(Matchers_ConstructorWithName(jutil.StringPtr("bar.Foo")).Matches(ctorNode, nil)).IsTrue()

	input = "/** @constructor */ function Foo() {};"
	root = matchersTest_compileToScriptRoot(matchersTest_getCompiler(input))
	ctorNode = root.GetFirstChild()
	AssertThat(Matchers_ConstructorWithName(jutil.StringPtr("Foo")).Matches(ctorNode, nil)).IsTrue()

	// TODO(mknichel): Make this test case work.
	// input =
	//     """
	//     ns = {
	//       /** @constructor */
	//       Foo: function() {}
	//     }
	//     """
	// root = matchersTest_compileToScriptRoot(matchersTest_getCompiler(input))
	// ctorNode = root.GetFirstFirstChild().GetLastChild().GetFirstChild()
	// assertTrue(Matchers_ConstructorWithName(jutil.StringPtr("ns.Foo")).Matches(ctorNode, nil))
}

// Java: MatchersTest.testNewClass()
func (c *MatchersTest) TestNewClass() {
	input := "new Object()"
	root := matchersTest_compileToScriptRoot(matchersTest_getCompiler(input))
	newNode := root.GetFirstFirstChild()
	AssertThat(newNode.IsNew()).IsTrue()
	AssertThat(Matchers_NewClass().Matches(newNode, nil)).IsTrue()
	AssertThat(Matchers_NewClass().Matches(newNode.GetFirstChild(), nil)).IsFalse()
}

// Java: MatchersTest.testNewClass_specificClass()
func (c *MatchersTest) TestNewClass_specificClass() {
	externs :=
		"/** @constructor */\n" +
			"function Foo() {};\n" +
			"/** @constructor */\n" +
			"function Bar() {};\n"
	input := "var foo = new Foo();"
	compiler := matchersTest_getCompilerWithExterns(externs, input)
	metadata := NewNodeMetadata(compiler)
	root := matchersTest_compileToScriptRoot(compiler)
	varNode := root.GetFirstChild()
	newNode := varNode.GetFirstFirstChild()
	AssertThat(newNode.IsNew()).IsTrue()
	AssertThat(Matchers_NewClassWithName("Foo").Matches(newNode, metadata)).IsTrue()
	AssertThat(Matchers_NewClassWithName("Bar").Matches(newNode, metadata)).IsFalse()
	AssertThat(Matchers_NewClassWithName("Foo").Matches(newNode.GetFirstChild(), metadata)).IsFalse()
}

// Java: MatchersTest.testFunctionCall_any()
func (c *MatchersTest) TestFunctionCall_any() {
	input := "goog.base(this)"
	root := matchersTest_compileToScriptRoot(matchersTest_getCompiler(input))
	fnCall := root.GetFirstFirstChild()
	AssertThat(fnCall.IsCall()).IsTrue()
	AssertThat(Matchers_FunctionCall().Matches(fnCall, nil)).IsTrue()
}

// Java: MatchersTest.testFunctionCall_numArgs()
func (c *MatchersTest) TestFunctionCall_numArgs() {
	input := "goog.base(this)"
	root := matchersTest_compileToScriptRoot(matchersTest_getCompiler(input))
	fnCall := root.GetFirstFirstChild()
	AssertThat(fnCall.IsCall()).IsTrue()
	AssertThat(Matchers_FunctionCallWithNumArgs(1).Matches(fnCall, nil)).IsTrue()
	AssertThat(Matchers_FunctionCallWithNumArgs(2).Matches(fnCall, nil)).IsFalse()
	AssertThat(Matchers_FunctionCallWithNumArgsAndName("goog.base", 1).Matches(fnCall, nil)).IsTrue()
	AssertThat(Matchers_FunctionCallWithNumArgsAndName("goog.base", 2).Matches(fnCall, nil)).IsFalse()
	AssertThat(Matchers_FunctionCallWithNumArgsAndName("goog.require", 1).Matches(fnCall, nil)).IsFalse()
	AssertThat(Matchers_FunctionCallWithNumArgsAndName("goog.require", 2).Matches(fnCall, nil)).IsFalse()
}

// Java: MatchersTest.testFunctionCall_static()
func (c *MatchersTest) TestFunctionCall_static() {
	input := "goog.require('goog.dom');"
	root := matchersTest_compileToScriptRoot(matchersTest_getCompiler(input))
	fnCall := root.GetFirstFirstChild()
	AssertThat(fnCall.IsCall()).IsTrue()
	AssertThat(Matchers_FunctionCallWithName(jutil.StringPtr("goog.require")).Matches(fnCall, nil)).IsTrue()
	AssertThat(Matchers_FunctionCallWithName(jutil.StringPtr("goog.provide")).Matches(fnCall, nil)).IsFalse()
}

// Java: MatchersTest.testFunctionCall_prototype()
func (c *MatchersTest) TestFunctionCall_prototype() {
	externs :=
		"/** @constructor */\n" +
			"function Foo() {};\n" +
			"Foo.prototype.bar = function() {};\n" +
			"Foo.prototype.baz = function() {};\n"
	input := "var foo = new Foo(); foo.bar();"
	compiler := matchersTest_getCompilerWithExterns(externs, input)
	metadata := NewNodeMetadata(compiler)
	root := matchersTest_compileToScriptRoot(compiler)
	fnCall := root.GetLastChild().GetFirstChild()
	AssertThat(fnCall.IsCall()).IsTrue()
	AssertThat(Matchers_FunctionCallWithName(jutil.StringPtr("Foo.prototype.bar")).Matches(fnCall, metadata)).IsTrue()
	AssertThat(Matchers_FunctionCallWithName(jutil.StringPtr("Foo.prototype.baz")).Matches(fnCall, metadata)).IsFalse()
}

// Java: MatchersTest.testGoogModule()
func (c *MatchersTest) TestGoogModule() {
	input := "goog.module('testcase');"
	compiler := matchersTest_getCompiler(input)
	root := matchersTest_compileToScriptRoot(compiler)
	fnCall := root.GetFirstFirstChild()
	metadata := NewNodeMetadata(compiler)
	AssertThat(Matchers_GoogModule().Matches(fnCall, metadata)).IsTrue()
	AssertThat(Matchers_GoogModuleOrProvide().Matches(fnCall, metadata)).IsTrue()
}

// Java: MatchersTest.testEnum()
func (c *MatchersTest) TestEnum() {
	input := "/** @enum {string} */ var foo = {BAR: 'baz'};"
	root := matchersTest_compileToScriptRoot(matchersTest_getCompiler(input))
	enumNode := root.GetFirstFirstChild()
	AssertThat(Matchers_EnumDefinition().Matches(enumNode, nil)).IsTrue()
}

// Java: MatchersTest.testEnumOfType()
func (c *MatchersTest) TestEnumOfType() {
	input := "/** @enum {string} */ var foo = {BAR: 'baz'};"
	compiler := matchersTest_getCompiler(input)
	root := matchersTest_compileToScriptRoot(compiler)
	enumNode := root.GetFirstFirstChild()
	AssertThat(Matchers_EnumDefinitionOfType("string").Matches(
		enumNode, NewNodeMetadata(compiler))).IsTrue()
	AssertThat(Matchers_EnumDefinitionOfType("number").Matches(
		enumNode, NewNodeMetadata(compiler))).IsFalse()
}

// Java: MatchersTest.testAssignmentWithRhs()
func (c *MatchersTest) TestAssignmentWithRhs() {
	externs :=
		"var goog = {};\n" +
			"goog.number = function() {};\n" +
			"var someObj = {};\n"
	input := "someObj.foo = goog.number();"
	compiler := matchersTest_getCompilerWithExterns(externs, input)
	root := matchersTest_compileToScriptRoot(compiler)
	assignNode := root.GetFirstFirstChild()
	AssertThat(
		Matchers_AssignmentWithRhs(Matchers_FunctionCallWithName(jutil.StringPtr("goog.number"))).Matches(assignNode, nil)).IsTrue()
	AssertThat(
		Matchers_AssignmentWithRhs(Matchers_FunctionCallWithName(jutil.StringPtr("goog.base"))).Matches(assignNode, nil)).IsFalse()
	AssertThat(Matchers_AssignmentWithRhs(Matchers_FunctionCallWithName(jutil.StringPtr("goog.number"))).
		Matches(assignNode.GetFirstChild(), nil)).IsFalse()
}

// Java: MatchersTest.testPrototypeDeclarations()
func (c *MatchersTest) TestPrototypeDeclarations() {
	input :=
		"/** @constructor */\n" +
			"function Foo() {}\n" +
			"Foo.prototype.bar = 3;\n" +
			"Foo.prototype.baz = function() {};\n"
	compiler := matchersTest_getCompiler(input)
	metadata := NewNodeMetadata(compiler)
	root := matchersTest_compileToScriptRoot(compiler)
	prototypeVarAssign := root.GetSecondChild().GetFirstChild()
	prototypeFnAssign := root.GetLastChild().GetFirstChild()

	AssertThat(Matchers_PrototypeVariableDeclaration().Matches(
		prototypeVarAssign.GetFirstChild(), metadata)).IsTrue()
	AssertThat(Matchers_PrototypeMethodDeclaration().Matches(
		prototypeVarAssign.GetFirstChild(), metadata)).IsFalse()

	AssertThat(Matchers_PrototypeMethodDeclaration().Matches(
		prototypeFnAssign.GetFirstChild(), metadata)).IsTrue()
	AssertThat(Matchers_PrototypeVariableDeclaration().Matches(
		prototypeFnAssign.GetFirstChild(), metadata)).IsFalse()
}

// Java: MatchersTest.testJsDocType1()
func (c *MatchersTest) TestJsDocType1() {
	input := "/** @type {number} */ var foo = 1;"
	compiler := matchersTest_getCompiler(input)
	root := matchersTest_compileToScriptRoot(compiler)
	node := root.GetFirstFirstChild()
	AssertThat(Matchers_JsDocType("number").Matches(node, NewNodeMetadata(compiler))).IsTrue()
	AssertThat(Matchers_JsDocType("string").Matches(node, NewNodeMetadata(compiler))).IsFalse()
}

// Java: MatchersTest.testJsDocType2()
func (c *MatchersTest) TestJsDocType2() {
	input := "/** @type {number} */ let foo = 1;"
	compiler := matchersTest_getCompiler(input)
	root := matchersTest_compileToScriptRoot(compiler)
	node := root.GetFirstFirstChild()
	AssertThat(Matchers_JsDocType("number").Matches(node, NewNodeMetadata(compiler))).IsTrue()
	AssertThat(Matchers_JsDocType("string").Matches(node, NewNodeMetadata(compiler))).IsFalse()
}

// Java: MatchersTest.testJsDocType3()
func (c *MatchersTest) TestJsDocType3() {
	input := "/** @type {number} */ const foo = 1;"
	compiler := matchersTest_getCompiler(input)
	root := matchersTest_compileToScriptRoot(compiler)
	node := root.GetFirstFirstChild()
	AssertThat(Matchers_JsDocType("number").Matches(node, NewNodeMetadata(compiler))).IsTrue()
	AssertThat(Matchers_JsDocType("string").Matches(node, NewNodeMetadata(compiler))).IsFalse()
}

// Java: MatchersTest.testJsDocTypeNoMatch1()
func (c *MatchersTest) TestJsDocTypeNoMatch1() {
	input := "/** @const */ var foo = 1;"
	compiler := matchersTest_getCompiler(input)
	root := matchersTest_compileToScriptRoot(compiler)
	node := root.GetFirstFirstChild()
	AssertThat(Matchers_JsDocType("number").Matches(node, NewNodeMetadata(compiler))).IsFalse()
	AssertThat(Matchers_JsDocType("string").Matches(node, NewNodeMetadata(compiler))).IsFalse()
}

// Java: MatchersTest.testJsDocTypeNoMatch2()
func (c *MatchersTest) TestJsDocTypeNoMatch2() {
	input := "const foo = 1;"
	compiler := matchersTest_getCompiler(input)
	root := matchersTest_compileToScriptRoot(compiler)
	node := root.GetFirstFirstChild()
	AssertThat(Matchers_JsDocType("number").Matches(node, NewNodeMetadata(compiler))).IsFalse()
	AssertThat(Matchers_JsDocType("string").Matches(node, NewNodeMetadata(compiler))).IsFalse()
}

// Java: MatchersTest.testPropertyAccess()
func (c *MatchersTest) TestPropertyAccess() {
	input := "foo.bar.method();"
	compiler := matchersTest_getCompiler(input)
	root := matchersTest_compileToScriptRoot(compiler)
	methodNode := root.GetFirstFirstChild().GetFirstChild()
	barNode := methodNode.GetFirstChild()
	AssertThat(Matchers_PropertyAccess().Matches(methodNode, NewNodeMetadata(compiler))).IsTrue()
	AssertThat(Matchers_PropertyAccess().Matches(barNode, NewNodeMetadata(compiler))).IsTrue()
	AssertThat(Matchers_PropertyAccessWithName(jutil.StringPtr("foo.bar.method")).Matches(
		methodNode, NewNodeMetadata(compiler))).IsTrue()
	AssertThat(Matchers_PropertyAccessWithName(jutil.StringPtr("foo.bar")).Matches(
		barNode, NewNodeMetadata(compiler))).IsTrue()
	AssertThat(Matchers_PropertyAccessWithName(jutil.StringPtr("foo")).Matches(
		barNode.GetFirstChild(), NewNodeMetadata(compiler))).IsFalse()
}

// Java: MatchersTest.testPropertyAccess_instance()
func (c *MatchersTest) TestPropertyAccess_instance() {
	externs :=
		"/** @constructor */\n" +
			"function Foo() {};\n" +
			"Foo.prototype.bar = 3;\n"
	input := "var foo = new Foo(); foo.bar;"
	compiler := matchersTest_getCompilerWithExterns(externs, input)
	root := matchersTest_compileToScriptRoot(compiler)
	node := root.GetLastChild().GetFirstChild()
	AssertThat(Matchers_PropertyAccess().Matches(node, NewNodeMetadata(compiler))).IsTrue()
	AssertThat(Matchers_PropertyAccessWithName(jutil.StringPtr("Foo.prototype.bar")).Matches(
		node, NewNodeMetadata(compiler))).IsTrue()
	AssertThat(Matchers_PropertyAccessWithName(jutil.StringPtr("foo.bar")).Matches(
		node, NewNodeMetadata(compiler))).IsTrue()
}

// Java: MatchersTest.testConstructorPropertyDeclaration()
func (c *MatchersTest) TestConstructorPropertyDeclaration() {
	externs := ""
	input :=
		"/** @constructor */\n" +
			"function MyClass() {\n" +
			"  this.foo = 5;\n" +
			"  var bar = 10;\n" +
			"}\n"
	compiler := matchersTest_getCompilerWithExterns(externs, input)
	root := matchersTest_compileToScriptRoot(compiler)
	// The ASSIGN node
	node := root.GetFirstChild().GetLastChild().GetFirstFirstChild()
	AssertThat(
		Matchers_ConstructorPropertyDeclaration().Matches(node, NewNodeMetadata(compiler))).IsTrue()

	// The VAR node
	node = root.GetFirstChild().GetLastChild().GetLastChild().GetFirstChild()
	AssertThat(
		Matchers_ConstructorPropertyDeclaration().Matches(node, NewNodeMetadata(compiler))).IsFalse()
}

// Java: MatchersTest.testIsPrivate()
func (c *MatchersTest) TestIsPrivate() {
	externs := ""
	input :=
		"/** @private */\n" +
			"var foo = 3;\n"
	compiler := matchersTest_getCompilerWithExterns(externs, input)
	root := matchersTest_compileToScriptRoot(compiler)
	node := root.GetFirstChild()
	AssertThat(Matchers_IsPrivate().Matches(node, NewNodeMetadata(compiler))).IsTrue()

	input =
		"/** @package */\n" +
			"var foo = 3;\n"
	compiler = matchersTest_getCompilerWithExterns(externs, input)
	root = matchersTest_compileToScriptRoot(compiler)
	node = root.GetFirstChild()
	AssertThat(Matchers_IsPrivate().Matches(node, NewNodeMetadata(compiler))).IsFalse()
}

// matchersTest_compileToScriptRoot returns the root script node produced
// from the compiled JS input.
//
// Java: MatchersTest.compileToScriptRoot(Compiler)
func matchersTest_compileToScriptRoot(compiler *Compiler) *rhino.Node {
	root := compiler.GetRoot()
	// The last child of the compiler root is a Block node, and the first child
	// of that is the Script node.
	return root.GetLastChild().GetFirstChild()
}

// Java: MatchersTest.getCompiler(String)
func matchersTest_getCompiler(jsInput string) *Compiler {
	return matchersTest_getCompilerWithExterns("", jsInput)
}

// Java: MatchersTest.getCompiler(String, String)
func matchersTest_getCompilerWithExterns(externs string, jsInput string) *Compiler {
	compiler := NewCompilerWithErrorManager(NewBlackHoleErrorManager())
	compiler.DisableThreads()
	options := RefactoringDriver_GetCompilerOptions()
	compiler.Compile(
		[]*SourceFile{SourceFile_FromCode("externs", externs)},
		[]*SourceFile{SourceFile_FromCode("test", jsInput)},
		options)
	return compiler
}
