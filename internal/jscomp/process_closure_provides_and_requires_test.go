// Port of com/google/javascript/jscomp/ProcessClosureProvidesAndRequiresTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/ProcessClosureProvidesAndRequiresTest.java)

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
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ProcessClosureProvidesAndRequiresTest holds unit tests for
// ProcessClosureProvidesAndRequires
type ProcessClosureProvidesAndRequiresTest struct {
	CompilerTestCaseBase

	preserveGoogProvidesAndRequires bool
	lastProcessor                   *ProcessClosureProvidesAndRequires
}

// Java: ProcessClosureProvidesAndRequiresTest()
func NewProcessClosureProvidesAndRequiresTest() *ProcessClosureProvidesAndRequiresTest {
	c := &ProcessClosureProvidesAndRequiresTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, CompilerTestCase_MINIMAL_EXTERNS+NewTestExternsBuilder().AddClosureExterns().Build())
	return c
}

func TestProcessClosureProvidesAndRequiresTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewProcessClosureProvidesAndRequiresTest() })
}

// Java: ProcessClosureProvidesAndRequiresTest.createClosureProcessor(Compiler)
func (c *ProcessClosureProvidesAndRequiresTest) CreateClosureProcessor(compiler *Compiler) *ProcessClosureProvidesAndRequires {
	return NewProcessClosureProvidesAndRequires(compiler, c.preserveGoogProvidesAndRequires)
}

// Java: ProcessClosureProvidesAndRequiresTest.setUp()
func (c *ProcessClosureProvidesAndRequiresTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_2017)
	c.preserveGoogProvidesAndRequires = false
	c.EnableTypeCheck()
	c.EnableCreateModuleMap() // necessary for the typechecker
	c.ReplaceTypesWithColors()
	c.lastProcessor = nil
}

// Java: ProcessClosureProvidesAndRequiresTest.getProcessor(Compiler)
func (c *ProcessClosureProvidesAndRequiresTest) GetProcessor(compiler *Compiler) CompilerPass {
	return CompilerPassFunc(func(externs *rhino.Node, root *rhino.Node) {
		c.VerifyCollectProvidedNamesDoesntChangeAst(externs, root, compiler)
		c.lastProcessor = c.CreateClosureProcessor(compiler)
		c.lastProcessor.RewriteProvidesAndRequires(externs, root)
	})
}

// Java: ProcessClosureProvidesAndRequiresTest.getProvidedNameCollection(String)
func (c *ProcessClosureProvidesAndRequiresTest) GetProvidedNameCollection(js string) *jutil.LinkedHashMap[string, *ProcessClosureProvidesAndRequires_ProvidedName] {
	compiler := c.CreateCompiler()
	processor := c.CreateClosureProcessor(compiler)
	jsRoot := compiler.ParseTestCode(js)
	scopeRoot := rhino.IR_Root(rhino.IR_Root(), rhino.IR_Root(jsRoot))
	providedNameMap :=
		processor.CollectProvidedNames(scopeRoot.GetFirstChild(), scopeRoot.GetSecondChild())
	AssertThat(compiler.GetWarnings()).IsEmpty()
	AssertThat(compiler.GetErrors()).IsEmpty()
	return providedNameMap
}

// VerifyCollectProvidedNamesDoesntChangeAst validates that running
// ProcessClosureProvidesAndRequires#collectProvidedNames(Node, Node) does
// not modify the AST.
//
// This is important because we want to call this method to gain
// information about goog.provides but preserve the original AST structure
// for future checks.
//
// Java: ProcessClosureProvidesAndRequiresTest.verifyCollectProvidedNamesDoesntChangeAst(Node, Node, Compiler)
func (c *ProcessClosureProvidesAndRequiresTest) VerifyCollectProvidedNamesDoesntChangeAst(
	externs *rhino.Node, root *rhino.Node, compiler *Compiler) {
	// Validate that this does not modify the AST at all!
	originalExterns := externs.CloneTree()
	originalRoot := root.CloneTree()
	processor := c.CreateClosureProcessor(compiler)
	processor.CollectProvidedNames(externs, root)

	AssertNode(externs).IsEqualIncludingJsDocTo(originalExterns)
	AssertNode(root).IsEqualIncludingJsDocTo(originalRoot)
}

// Java: ProcessClosureProvidesAndRequiresTest.testChunk(String[], String[])
func (c *ProcessClosureProvidesAndRequiresTest) TestChunk(chunkInputs []string, expected []string) {
	c.Test(c.Srcs(JSChunkGraphBuilder_ForStar().AddChunks(chunkInputs).Build()), c.Expected(expected))
}

// Java: ProcessClosureProvidesAndRequiresTest.testTypedefProvides_withProvidedParent()
func (c *ProcessClosureProvidesAndRequiresTest) TestTypedefProvides_withProvidedParent() {
	c.Test(
		c.Srcs(
			`goog.provide('ns');
goog.provide('ns.SomeType');
goog.provide('ns.SomeType.EnumValue');
goog.provide('ns.SomeType.defaultName');
goog.provide('ns.SomeType.NestedType');
// subnamespace assignment happens before parent.
/** @enum {number} */
ns.SomeType.EnumValue = { A: 1, B: 2 };
// parent namespace isn't ever actually assigned.
// we're relying on goog.provide to provide it.
/** @typedef {{name: string, value: ns.SomeType.EnumValue}} */
ns.SomeType;
/** @const {string} */
ns.SomeType.defaultName = 'foobarbaz';
/** @typedef {number} */
ns.SomeType.NestedType;
`),
		c.Expected(
			`/** @const */ var ns = {};
// Created from goog.provide.
/** @const */
ns.SomeType = {}; // created from goog.provide
/** @const */
ns.SomeType.NestedType = {}; // created from goog.provide
/** @enum {!JSDocSerializer_placeholder_type} */
ns.SomeType.EnumValue = {A:1, B:2};
ns.SomeType;
/** @const */
ns.SomeType.defaultName = 'foobarbaz';
ns.SomeType.NestedType;
`))
}

// Java: ProcessClosureProvidesAndRequiresTest.testTypedefProvidesWithExplicitParentNamespace_errorInCollectProvidesMode()
func (c *ProcessClosureProvidesAndRequiresTest) TestTypedefProvidesWithExplicitParentNamespace_errorInCollectProvidesMode() {
	c.TestError(
		c.Srcs(
			`goog.provide('foo.bar');
goog.provide('foo.bar.Type');

foo.bar = function() {};
/** @typedef {string} */ foo.bar.Type;
`),
		c.Error(ProcessClosureProvidesAndRequires_TYPEDEF_CHILD_OF_PROVIDE))
}

// Java: ProcessClosureProvidesAndRequiresTest.testSimpleProvides()
func (c *ProcessClosureProvidesAndRequiresTest) TestSimpleProvides() {
	c.Test(c.Srcs("goog.provide('foo');"), c.Expected("/** @const */ var foo={};"))
	c.Test(
		c.Srcs("goog.provide('foo.bar');"),
		c.Expected("/** @const */ var foo={}; /** @const */ foo.bar={};"))
	c.Test(
		c.Srcs("goog.provide('foo.bar.baz');"),
		c.Expected(
			"/** @const */ var foo={}; /** @const */ foo.bar={}; /** @const */ foo.bar.baz={};"))
	c.Test(
		c.Srcs("goog.provide('foo.bar.baz.boo');"),
		c.Expected(
			`/** @const */ var foo={};
/** @const */ foo.bar={};
/** @const */ foo.bar.baz={};
/** @const */ foo.bar.baz.boo={};
`))
	// goog is special-cased
	c.Test(c.Srcs("goog.provide('goog.bar');"), c.Expected("/** @const */ goog.bar={};"))
}

// Java: ProcessClosureProvidesAndRequiresTest.testMultipleProvides()
func (c *ProcessClosureProvidesAndRequiresTest) TestMultipleProvides() {
	c.Test(
		c.Srcs("goog.provide('foo.bar'); goog.provide('foo.baz');"),
		c.Expected("/** @const */ var foo={}; /** @const */ foo.bar={}; /** @const */ foo.baz={};"))

	c.Test(
		c.Srcs("goog.provide('foo.bar.baz'); goog.provide('foo.boo.foo');"),
		c.Expected(
			`/** @const */
var foo = {};
/** @const */
foo.bar={};
/** @const */
foo.bar.baz={};
/** @const */
foo.boo={};
/** @const */
foo.boo.foo={};
`))

	c.Test(
		c.Srcs("goog.provide('foo.bar.baz'); goog.provide('foo.bar.boo');"),
		c.Expected(
			`/** @const */
var foo = {};
/** @const */
foo.bar = {};
/** @const */
foo.bar.baz = {};
/** @const */
foo.bar.boo={};
`))

	c.Test(
		c.Srcs("goog.provide('foo.bar.baz'); goog.provide('goog.bar.boo');"),
		c.Expected(
			`/** @const */
var foo = {};
/** @const */
foo.bar={};
/** @const */
foo.bar.baz={};
/** @const */
goog.bar={};
/** @const */
goog.bar.boo = {};
`))
}

// Java: ProcessClosureProvidesAndRequiresTest.testProvideAlreadyProvided()
func (c *ProcessClosureProvidesAndRequiresTest) TestProvideAlreadyProvided() {
	c.Test(
		c.Srcs(
			`goog.provide('a.b')
/** @provideAlreadyProvided */
goog.provide('a.b.c')
a.b = something;
`),
		c.Expected(
			`/** @const */ var a = {};
a.b = something;
`))
}

// Java: ProcessClosureProvidesAndRequiresTest.testRemovalOfProvidedObjLit()
func (c *ProcessClosureProvidesAndRequiresTest) TestRemovalOfProvidedObjLit() {
	c.Test(c.Srcs("goog.provide('foo'); foo = 0;"), c.Expected("var foo = 0;"))
	c.Test(c.Srcs("goog.provide('foo'); foo = {a: 0};"), c.Expected("var foo = {a: 0};"))

	c.Test(c.Srcs("goog.provide('foo'); foo = function(){};"), c.Expected("var foo = function(){};"))

	c.Test(c.Srcs("goog.provide('foo'); foo = ()=>{};"), c.Expected("var foo = ()=>{};"))

	c.Test(c.Srcs("goog.provide('foo'); var foo = 0;"), c.Expected("var foo = 0;"))

	c.Test(c.Srcs("goog.provide('foo'); let foo = 0;"), c.Expected("let foo = 0;"))

	c.Test(c.Srcs("goog.provide('foo'); const foo = 0;"), c.Expected("const foo = 0;"))

	c.Test(c.Srcs("goog.provide('foo'); var foo = {a: 0};"), c.Expected("var foo = {a: 0};"))

	c.Test(c.Srcs("goog.provide('foo'); var foo = function(){};"), c.Expected("var foo = function(){};"))

	c.Test(c.Srcs("goog.provide('foo'); var foo = ()=>{};"), c.Expected("var foo = ()=>{};"))

	c.Test(
		c.Srcs("goog.provide('foo.bar.Baz'); foo.bar.Baz=function(){};"),
		c.Expected(
			`/** @const */
var foo={};
/** @const */
foo.bar = {};
foo.bar.Baz=function(){};
`))
	c.Test(
		c.Srcs("goog.provide('foo.bar.moo'); foo.bar.moo={E:1,S:2};"),
		c.Expected(
			`/** @const */
var foo={};
/** @const */
foo.bar={};
foo.bar.moo={E:1,S:2};
`))

	c.Test(
		c.Srcs("goog.provide('foo.bar.moo'); foo.bar.moo={E:1}; foo.bar.moo={E:2};"),
		c.Expected(
			`/** @const */
var foo={};
/** @const */
foo.bar={};
foo.bar.moo={E:1};
foo.bar.moo={E:2};
`))

	c.Test(c.Srcs("goog.provide('foo'); var foo = class {}"), c.Expected("var foo = class {}"))
}

// Java: ProcessClosureProvidesAndRequiresTest.testRemovalMultipleAssignment1()
func (c *ProcessClosureProvidesAndRequiresTest) TestRemovalMultipleAssignment1() {
	c.Test(c.Srcs("goog.provide('foo'); foo = 0; foo = 1"), c.Expected("var foo = 0; foo = 1;"))
}

// Java: ProcessClosureProvidesAndRequiresTest.testRemovalMultipleAssignment2()
func (c *ProcessClosureProvidesAndRequiresTest) TestRemovalMultipleAssignment2() {
	c.Test(c.Srcs("goog.provide('foo'); var foo = 0; foo = 1"), c.Expected("var foo = 0; foo = 1;"))

	c.Test(
		c.Srcs("goog.provide('foo'); let foo = 0; let foo = 1"),
		c.Expected("let foo = 0; let foo = 1;"))
}

// Java: ProcessClosureProvidesAndRequiresTest.testRemovalMultipleAssignment3()
func (c *ProcessClosureProvidesAndRequiresTest) TestRemovalMultipleAssignment3() {
	c.Test(c.Srcs("goog.provide('foo'); foo = 0; var foo = 1"), c.Expected("foo = 0; var foo = 1;"))

	c.Test(c.Srcs("goog.provide('foo'); foo = 0; let foo = 1"), c.Expected("foo = 0; let foo = 1;"))
}

// Java: ProcessClosureProvidesAndRequiresTest.testRemovalMultipleAssignment4()
func (c *ProcessClosureProvidesAndRequiresTest) TestRemovalMultipleAssignment4() {
	c.Test(
		c.Srcs("goog.provide('foo.bar'); foo.bar = 0; foo.bar = 1"),
		c.Expected("/** @const */ var foo = {}; foo.bar = 0; foo.bar = 1"))
}

// Java: ProcessClosureProvidesAndRequiresTest.testNoRemovalFunction1()
func (c *ProcessClosureProvidesAndRequiresTest) TestNoRemovalFunction1() {
	c.Test(
		c.Srcs(
			`goog.provide('foo');
function f(){
/** @suppress {checkTypes} */
  foo = 0;
}
`),
		c.Expected(
			`/** @const */ var foo = {};
function f(){
  foo = 0;
}
`))
}

// Java: ProcessClosureProvidesAndRequiresTest.testNoRemovalFunction2()
func (c *ProcessClosureProvidesAndRequiresTest) TestNoRemovalFunction2() {
	c.Test(
		c.Srcs("goog.provide('foo'); function f(){var foo = 0}"),
		c.Expected("/** @const */ var foo = {}; function f(){var foo = 0}"))
}

// Java: ProcessClosureProvidesAndRequiresTest.testNoRemovalFunction3()
func (c *ProcessClosureProvidesAndRequiresTest) TestNoRemovalFunction3() {
	c.Test(
		c.Srcs("goog.provide('foo'); function f(foo = 0){}"),
		c.Expected("/** @const */ var foo = {}; function f(foo = 0){}"))
}

// Java: ProcessClosureProvidesAndRequiresTest.testRemovalMultipleAssignmentInIf1()
func (c *ProcessClosureProvidesAndRequiresTest) TestRemovalMultipleAssignmentInIf1() {
	c.Test(
		c.Srcs("goog.provide('foo'); if (true) { var foo = 0 } else { foo = 1 }"),
		c.Expected("if (true) { var foo = 0 } else { foo = 1 }"))
}

// Java: ProcessClosureProvidesAndRequiresTest.testRemovalMultipleAssignmentInIf2()
func (c *ProcessClosureProvidesAndRequiresTest) TestRemovalMultipleAssignmentInIf2() {
	c.Test(
		c.Srcs("goog.provide('foo'); if (true) { foo = 0 } else { var foo = 1 }"),
		c.Expected("if (true) { foo = 0 } else { var foo = 1 }"))
}

// Java: ProcessClosureProvidesAndRequiresTest.testRemovalMultipleAssignmentInIf3()
func (c *ProcessClosureProvidesAndRequiresTest) TestRemovalMultipleAssignmentInIf3() {
	c.Test(
		c.Srcs("goog.provide('foo'); if (true) { foo = 0 } else { foo = 1 }"),
		c.Expected("if (true) { var foo = 0 } else { foo = 1 }"))
}

// Java: ProcessClosureProvidesAndRequiresTest.testRemovalMultipleAssignmentInIf4()
func (c *ProcessClosureProvidesAndRequiresTest) TestRemovalMultipleAssignmentInIf4() {
	c.Test(
		c.Srcs("goog.provide('foo.bar'); if (true) { foo.bar = 0 } else { foo.bar = 1 }"),
		c.Expected(
			`/** @const */ var foo = {};
if (true) {
  foo.bar = 0;
} else {
  foo.bar = 1;
}
`))
}

// Java: ProcessClosureProvidesAndRequiresTest.testMultipleDeclarationError1()
func (c *ProcessClosureProvidesAndRequiresTest) TestMultipleDeclarationError1() {
	rest := "if (true) { foo.bar = 0 } else { foo.bar = 1 }"
	c.Test(
		c.Srcs(
			`goog.provide('foo.bar');
var foo = {};
`+
				rest),
		c.Expected(
			`var foo = {};
var foo = {};
`+
				rest))
}

// Java: ProcessClosureProvidesAndRequiresTest.testMultipleDeclarationError2()
func (c *ProcessClosureProvidesAndRequiresTest) TestMultipleDeclarationError2() {
	c.Test(
		c.Srcs(
			`goog.provide('foo.bar');
if (true) { var foo = {}; foo.bar = 0 } else { foo.bar = 1 }
`),
		c.Expected(
			`var foo = {};
if (true) {
  var foo = {}; foo.bar = 0
} else {
  foo.bar = 1
}
`),
		// TODO(b/149765184): this warning should not happen
		c.Warning(TypeValidator_TYPE_MISMATCH_WARNING))
}

// Java: ProcessClosureProvidesAndRequiresTest.testMultipleDeclarationError3()
func (c *ProcessClosureProvidesAndRequiresTest) TestMultipleDeclarationError3() {
	c.Test(
		c.Srcs(
			`goog.provide('foo.bar');
if (true) { foo.bar = 0 } else { var foo = {}; foo.bar = 1 }
`),
		c.Expected(
			`var foo = {};
if (true) {
  foo.bar = 0
} else {
  var foo = {}; foo.bar = 1
}
`),
		// TODO(b/149765184): this warning should not happen
		c.Warning(TypeValidator_TYPE_MISMATCH_WARNING))
}

// Java: ProcessClosureProvidesAndRequiresTest.testProvideAfterDeclarationError()
func (c *ProcessClosureProvidesAndRequiresTest) TestProvideAfterDeclarationError() {
	c.Test(c.Srcs("var x = 42; goog.provide('x');"), c.Expected("var x = 42; /** @const */ var x = {}"))
}

// Java: ProcessClosureProvidesAndRequiresTest.testProvideAfterDeclaration_noErrorInExterns()
func (c *ProcessClosureProvidesAndRequiresTest) TestProvideAfterDeclaration_noErrorInExterns() {
	c.Test(
		c.Externs(CompilerTestCase_DEFAULT_EXTERNS, "var x = {};"),
		c.Srcs("goog.provide('x');"),
		c.Expected("/** @const */ var x = {}"))
}

// Java: ProcessClosureProvidesAndRequiresTest.testDuplicateProvideDoesntCrash()
func (c *ProcessClosureProvidesAndRequiresTest) TestDuplicateProvideDoesntCrash() {
	c.IgnoreWarnings(ClosurePrimitiveErrors_DUPLICATE_NAMESPACE)
	c.Test(
		c.Srcs("goog.provide('foo'); goog.provide('foo');"),
		c.Expected("/** @const */ var foo = {}; goog.provide('foo');"))

	c.Test(
		c.Srcs("goog.provide('foo.bar'); goog.provide('foo'); goog.provide('foo');"),
		c.Expected(
			`/** @const */ var foo={};
/** @const */ foo.bar = {};
goog.provide('foo');
`))
}

// Java: ProcessClosureProvidesAndRequiresTest.testRemovalOfRequires()
func (c *ProcessClosureProvidesAndRequiresTest) TestRemovalOfRequires() {
	c.Test(c.Srcs("goog.provide('foo'); goog.require('foo');"), c.Expected("/** @const */ var foo={};"))

	c.Test(
		c.Srcs("goog.provide('foo.bar'); goog.require('foo.bar');"),
		c.Expected("/** @const */ var foo={}; /** @const */ foo.bar={};"))
	c.Test(
		c.Srcs("goog.provide('foo.bar.baz'); goog.require('foo.bar.baz');"),
		c.Expected(
			"/** @const */ var foo={}; /** @const */ foo.bar={}; /** @const */ foo.bar.baz={};"))
	c.Test(
		c.Srcs("goog.provide('foo'); var x = 3; goog.require('foo'); something();"),
		c.Expected("/** @const */ var foo={}; var x = 3; something();"))
	c.TestSame("foo.require('foo.bar');")
}

// Java: ProcessClosureProvidesAndRequiresTest.testRemovalOfRequireType()
func (c *ProcessClosureProvidesAndRequiresTest) TestRemovalOfRequireType() {
	c.Test(
		c.Srcs("goog.provide('foo'); goog.requireType('foo');"),
		c.Expected("/** @const */ var foo={};"))
	c.Test(
		c.Srcs("goog.provide('foo.bar'); goog.requireType('foo.bar');"),
		c.Expected("/** @const */ var foo={}; /** @const */ foo.bar={};"))
	c.Test(
		c.Srcs("goog.provide('foo.bar.baz'); goog.requireType('foo.bar.baz');"),
		c.Expected(
			"/** @const */ var foo={}; /** @const */ foo.bar={}; /** @const */ foo.bar.baz={};"))
	c.Test(
		c.Srcs("goog.provide('foo'); var x = 3; goog.requireType('foo'); something();"),
		c.Expected("/** @const */ var foo={}; var x = 3; something();"))
	c.TestSame("foo.requireType('foo.bar');")
}

// Java: ProcessClosureProvidesAndRequiresTest.testPreserveGoogRequire()
func (c *ProcessClosureProvidesAndRequiresTest) TestPreserveGoogRequire() {
	c.preserveGoogProvidesAndRequires = true
	c.Test(
		c.Srcs("goog.provide('foo'); goog.require('foo');"),
		c.Expected("/** @const */ var foo={}; goog.provide('foo'); goog.require('foo');"))
	c.Test(
		c.Srcs("goog.provide('foo'); goog.require('foo'); var a = {};"),
		c.Expected(
			"/** @const */ var foo = {}; goog.provide('foo'); goog.require('foo'); var a = {};"))
}

// Java: ProcessClosureProvidesAndRequiresTest.testPreserveGoogRequireType()
func (c *ProcessClosureProvidesAndRequiresTest) TestPreserveGoogRequireType() {
	c.preserveGoogProvidesAndRequires = true

	c.Test(
		c.Srcs("goog.provide('foo'); goog.requireType('foo');"),
		c.Expected("/** @const */ var foo={}; goog.provide('foo'); goog.requireType('foo');"))
	c.Test(
		c.Srcs("goog.provide('foo'); goog.requireType('foo'); var a = {};"),
		c.Expected(
			`/** @const */ var foo = {}; goog.provide('foo'); goog.requireType('foo'); var a =
 {};
`))
}

// Java: ProcessClosureProvidesAndRequiresTest.testLateProvideForRequire()
func (c *ProcessClosureProvidesAndRequiresTest) TestLateProvideForRequire() {
	// Error reported elsewhere
	c.Test(
		c.Srcs("goog.require('foo'); goog.provide('foo');"), c.Expected("/** @const */ var foo = {};"))
}

// Java: ProcessClosureProvidesAndRequiresTest.testLateProvideForRequireType()
func (c *ProcessClosureProvidesAndRequiresTest) TestLateProvideForRequireType() {
	c.TestNoWarning("goog.requireType('foo'); goog.provide('foo');")
}

// Java: ProcessClosureProvidesAndRequiresTest.testMissingProvideForRequire()
func (c *ProcessClosureProvidesAndRequiresTest) TestMissingProvideForRequire() {
	c.Test(c.Srcs("goog.require('foo');"), c.Expected(""))
}

// Java: ProcessClosureProvidesAndRequiresTest.testMissingProvideForRequireType()
func (c *ProcessClosureProvidesAndRequiresTest) TestMissingProvideForRequireType() {
	c.Test(c.Srcs("goog.requireType('foo');"), c.Expected(""))
}

// Java: ProcessClosureProvidesAndRequiresTest.testProvideInExterns()
func (c *ProcessClosureProvidesAndRequiresTest) TestProvideInExterns() {
	c.AllowExternsChanges()

	c.Test(
		c.Externs(
			CompilerTestCase_DEFAULT_EXTERNS,
			`/** @externs */ goog.provide('animals.Dog');
/** @constructor */ animals.Dog = function() {}
`),
		c.Srcs("goog.require('animals.Dog'); new animals.Dog()"),
		c.Expected("new animals.Dog();"))
}

// Java: ProcessClosureProvidesAndRequiresTest.testForwardDeclarations()
func (c *ProcessClosureProvidesAndRequiresTest) TestForwardDeclarations() {
	c.Test(c.Srcs("goog.forwardDeclare('A.B')"), c.Expected(""))

	c.TestWarning(
		"goog.forwardDeclare();",
		TypeCheck_WRONG_ARGUMENT_COUNT) // This pass isn't responsible for invalid declares but
	// the typechecker will warn.
	c.TestWarning("goog.forwardDeclare('C', 'D');", TypeCheck_WRONG_ARGUMENT_COUNT)
}

// Java: ProcessClosureProvidesAndRequiresTest.testGoodCrossChunkRequire1()
func (c *ProcessClosureProvidesAndRequiresTest) TestGoodCrossChunkRequire1() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk("goog.provide('goog.ui');").
				AddChunk("").
				AddChunk("goog.require('goog.ui');").
				Build()),
		c.Expected("/** @const */ goog.ui = {};", "", ""))
}

// Java: ProcessClosureProvidesAndRequiresTest.testGoodCrossChunkRequire2()
func (c *ProcessClosureProvidesAndRequiresTest) TestGoodCrossChunkRequire2() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk("").
				AddChunk("").
				AddChunk("").
				AddChunk("goog.provide('goog.ui'); goog.require('goog.ui');").
				Build()),
		c.Expected("", "", "", "/** @const */ goog.ui = {};"))
}

// Java: ProcessClosureProvidesAndRequiresTest.testCrossChunkRequireType()
func (c *ProcessClosureProvidesAndRequiresTest) TestCrossChunkRequireType() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk("").
				AddChunk("goog.requireType('goog.ui');").
				AddChunk("").
				AddChunk("goog.provide('goog.ui')").
				Build()),
		c.Expected("", "", "", "/** @const */ goog.ui = {};"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk("").
				AddChunk("").
				AddChunk("goog.provide('goog.ui');").
				AddChunk("goog.requireType('goog.ui');").
				Build()),
		c.Expected("", "", "/** @const */ goog.ui = {};", ""))
}

// Java: ProcessClosureProvidesAndRequiresTest.testNamespaceInExterns()
func (c *ProcessClosureProvidesAndRequiresTest) TestNamespaceInExterns() {
	// Note: This style is not recommended but the compiler sort-of supports it.
	c.Test(
		c.Externs(CompilerTestCase_DEFAULT_EXTERNS, "var root = {}; /** @type {number} */ root.someProperty;"),
		c.Srcs(
			`goog.provide('root.branch.Leaf')
root.branch.Leaf = class {};
`),
		c.Expected(
			`/** @const */
var root = {};
/** @const */ root.branch = {};
root.branch.Leaf = class {};
`))
}

// Java: ProcessClosureProvidesAndRequiresTest.testNamespaceInExterns_withExplicitNamespaceReinitialization()
func (c *ProcessClosureProvidesAndRequiresTest) TestNamespaceInExterns_withExplicitNamespaceReinitialization() {
	// Note: This style is not recommended but the compiler sort-of supports it.
	c.Test(
		c.Externs(CompilerTestCase_DEFAULT_EXTERNS, "var root = {}; /** @type {number} */ root.someProperty;"),
		c.Srcs(
			`goog.provide('root.branch.Leaf')
var root = {};
root.branch.Leaf = class {};
`),
		c.Expected(
			`var root = {};
/** @const */ root.branch = {};
var root = {};
root.branch.Leaf = class {};
`))
}

// Java: ProcessClosureProvidesAndRequiresTest.testTypedefProvide_withChild()
func (c *ProcessClosureProvidesAndRequiresTest) TestTypedefProvide_withChild() {
	c.Test(
		c.Srcs(
			`goog.provide('foo.Bar');
goog.provide('foo.Bar.Baz');
/** @typedef {!Array<string>} */
foo.Bar;
foo.Bar.Baz = {};
`),
		c.Expected(
			`/** @const */ var foo = {};
/** @const */
foo.Bar = {};
foo.Bar;
foo.Bar.Baz = {}
`))
}

// Provide a name after the definition of the class providing the
// parent namespace.
//
// Java: ProcessClosureProvidesAndRequiresTest.testProvideOrder3a()
func (c *ProcessClosureProvidesAndRequiresTest) TestProvideOrder3a() {
	c.Test(
		c.Srcs(
			`goog.provide('a.b');
a.b = function(x,y) {};
goog.provide('a.b.c');
a.b.c;
`),
		c.Expected(
			`/** @const */
var a = {};
a.b = function(x,y) {};
/** @const */
a.b.c = {};
a.b.c;
`))
}

// Java: ProcessClosureProvidesAndRequiresTest.testProvideOrder4a()
func (c *ProcessClosureProvidesAndRequiresTest) TestProvideOrder4a() {
	c.Test(
		c.Srcs(
			`goog.provide('goog.a');
goog.provide('goog.a.b');
if (x) {
  goog.a.b = 1;
} else {
  goog.a.b = 2;
}
`),
		c.Expected(
			`/** @const */
goog.a={};
if(x)
  goog.a.b=1;
else
  goog.a.b=2;
`))
}

// Java: ProcessClosureProvidesAndRequiresTest.testInvalidRequireScope()
func (c *ProcessClosureProvidesAndRequiresTest) TestInvalidRequireScope() {
	c.Test(
		c.Srcs("goog.provide('a.b'); var x = x || goog.require('a.b');"),
		c.Expected("/** @const */ var a={};/** @const */ a.b={};var x=x||goog.require(\"a.b\")"))
	c.Test(
		c.Srcs("goog.provide('a.b'); x = goog.require('a.b');"),
		c.Expected("/** @const */ var a={};/** @const */ a.b={}; x= goog.require(\"a.b\")"))
	c.Test(
		c.Srcs("goog.provide('a.b'); function f() { goog.require('a.b'); }"),
		c.Expected("/** @const */ var a={};/** @const */ a.b={};function f(){goog.require('a.b')}"))
}

// Java: ProcessClosureProvidesAndRequiresTest.testImplicitAndExplicitProvide()
func (c *ProcessClosureProvidesAndRequiresTest) TestImplicitAndExplicitProvide() {
	c.Test(
		c.Srcs("goog.provide('goog.foo.bar'); goog.provide('goog.foo');"),
		c.Expected(
			`/** @const */
goog.foo = {};
/** @const */
goog.foo.bar = {};
`))
}

// Java: ProcessClosureProvidesAndRequiresTest.testImplicitProvideInIndependentChunks()
func (c *ProcessClosureProvidesAndRequiresTest) TestImplicitProvideInIndependentChunks() {
	c.TestChunk(
		[]string{"", "goog.provide('apps.A');", "goog.provide('apps.B');"},
		[]string{
			"/** @const */ var apps = {};",
			"/** @const */ apps.A = {};",
			"/** @const */ apps.B = {};",
		})
}

// Java: ProcessClosureProvidesAndRequiresTest.testImplicitProvideInIndependentChunks2()
func (c *ProcessClosureProvidesAndRequiresTest) TestImplicitProvideInIndependentChunks2() {
	c.TestChunk(
		[]string{
			"",
			"goog.provide('apps');", //
			"goog.provide('apps.foo.A');",
			"goog.provide('apps.foo.B');",
		},
		[]string{
			"/** @const */ var apps = {}; /** @const */ apps.foo = {};",
			"", // empty file
			"/** @const */ apps.foo.A = {};",
			"/** @const */ apps.foo.B = {};",
		})
}

// Java: ProcessClosureProvidesAndRequiresTest.testImplicitProvideInIndependentChunks_onGoogNamespace()
func (c *ProcessClosureProvidesAndRequiresTest) TestImplicitProvideInIndependentChunks_onGoogNamespace() {
	c.Test(
		c.Externs(""), // exclude the default Closure externs
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk(
					`/** @const */
var goog = {};
goog.provide = function(ns) {}
`).
				AddChunk("goog.provide('goog.foo.A');").
				AddChunk("goog.provide('goog.foo.B');").
				Build()),
		c.Expected(
			`/** @const */
var goog = {};
/** @const */ goog.foo = {};
goog.provide = function(ns) {}
`,
			"/** @const */ goog.foo.A = {};",
			"/** @const */ goog.foo.B = {};"))
}

// Java: ProcessClosureProvidesAndRequiresTest.testProvideInIndependentChunks1()
func (c *ProcessClosureProvidesAndRequiresTest) TestProvideInIndependentChunks1() {
	c.TestChunk(
		[]string{
			"", "goog.provide('apps');", "goog.provide('apps.foo');", "goog.provide('apps.foo.B');",
		},
		[]string{
			"/** @const */ var apps = {}; /** @const */ apps.foo = {};",
			"",
			"",
			"/** @const */ apps.foo.B = {};",
		})
}

// Java: ProcessClosureProvidesAndRequiresTest.testProvideInIndependentChunks2()
func (c *ProcessClosureProvidesAndRequiresTest) TestProvideInIndependentChunks2() {
	// TODO(nicksantos): Make this an error.
	c.TestChunk(
		[]string{
			"",
			"goog.provide('apps');",
			"goog.provide('apps.foo'); apps.foo = {};",
			"goog.provide('apps.foo.B');",
		},
		[]string{
			"/** @const */ var apps = {};", "", "apps.foo = {};", "/** @const */ apps.foo.B = {};",
		})
}

// Java: ProcessClosureProvidesAndRequiresTest.testProvideInIndependentChunks2b()
func (c *ProcessClosureProvidesAndRequiresTest) TestProvideInIndependentChunks2b() {
	// TODO(nicksantos): Make this an error.
	c.TestChunk(
		[]string{
			"",
			"goog.provide('apps');",
			"goog.provide('apps.foo'); apps.foo = function() {};",
			"goog.provide('apps.foo.B');",
		},
		[]string{
			"/** @const */ var apps = {};",
			"",
			"apps.foo = function() {};",
			"/** @const */ apps.foo.B = {};",
		})
}

// Java: ProcessClosureProvidesAndRequiresTest.testProvideInIndependentChunks3()
func (c *ProcessClosureProvidesAndRequiresTest) TestProvideInIndependentChunks3() {
	c.TestChunk(
		[]string{
			"",
			"goog.provide('apps');",
			"goog.provide('apps.foo.B');",
			"goog.provide('apps.foo'); goog.require('apps.foo');",
		},
		[]string{
			"/** @const */ var apps = {}; /** @const */ apps.foo = {};",
			"",
			"/** @const */ apps.foo.B = {};",
			"",
		})
}

// Java: ProcessClosureProvidesAndRequiresTest.testProvideInIndependentChunks3b()
func (c *ProcessClosureProvidesAndRequiresTest) TestProvideInIndependentChunks3b() {
	// TODO(nicksantos): Make this an error.
	c.TestChunk(
		[]string{
			"",
			"goog.provide('apps');",
			"goog.provide('apps.foo.B');",
			`goog.provide('apps.foo'); apps.foo = function() {};
goog.require('apps.foo');
`,
		},
		[]string{
			"/** @const */ var apps = {};",
			"",
			"/** @const */ apps.foo.B = {};",
			"apps.foo = function() {};",
		})
}

// Java: ProcessClosureProvidesAndRequiresTest.testProvideInIndependentChunks4()
func (c *ProcessClosureProvidesAndRequiresTest) TestProvideInIndependentChunks4() {
	// Regression test for bug 261:
	// http://blickly.github.io/closure-compiler-issues/#261
	c.TestChunk(
		[]string{
			"",
			"goog.provide('apps');",
			"goog.provide('apps.foo.bar.B');",
			"goog.provide('apps.foo.bar.C');",
		},
		[]string{
			`/** @const */ var apps={};
/** @const */ apps.foo={};
/** @const */ apps.foo.bar={}
`,
			"",
			"/** @const */ apps.foo.bar.B = {};",
			"/** @const */ apps.foo.bar.C = {};",
		})
}

// Java: ProcessClosureProvidesAndRequiresTest.testSourcePositionPreservation()
func (c *ProcessClosureProvidesAndRequiresTest) TestSourcePositionPreservation() {
	c.Test(
		c.Srcs("goog.provide('foo.bar.baz');"),
		c.Expected(
			`/** @const */ var foo = {};
/** @const */ foo.bar = {};
/** @const */ foo.bar.baz = {};
`))

	root := c.GetLastCompiler().GetRoot()

	fooDecl := c.FindQualifiedNameNode("foo", root)
	fooBarDecl := c.FindQualifiedNameNode("foo.bar", root)
	fooBarBazDecl := c.FindQualifiedNameNode("foo.bar.baz", root)

	AssertThat(fooDecl.GetLineno()).IsEqualTo(1)
	AssertThat(fooDecl.GetCharno()).IsEqualTo(14)

	AssertThat(fooBarDecl.GetLineno()).IsEqualTo(1)
	AssertThat(fooBarDecl.GetCharno()).IsEqualTo(18)

	AssertThat(fooBarBazDecl.GetLineno()).IsEqualTo(1)
	AssertThat(fooBarBazDecl.GetCharno()).IsEqualTo(22)
}

// Java: ProcessClosureProvidesAndRequiresTest.testProvidedTypedef_noParentProvide()
func (c *ProcessClosureProvidesAndRequiresTest) TestProvidedTypedef_noParentProvide() {
	c.Test(c.Srcs("goog.provide('x'); /** @typedef {number} */ var x;"), c.Expected("var x;"))
}

// Java: ProcessClosureProvidesAndRequiresTest.testProvidedTypedef_noParentProvide_doublyNestedNamespace()
func (c *ProcessClosureProvidesAndRequiresTest) TestProvidedTypedef_noParentProvide_doublyNestedNamespace() {
	c.Test(
		c.Srcs("goog.provide('x.y'); /** @typedef {number} */ x.y;"),
		c.Expected("/** @const */ var x = {}; /** @const */ x.y = {}; x.y;"))
}

// Java: ProcessClosureProvidesAndRequiresTest.testProvidedTypedef_noParentProvide_triplyNestedNamespace()
func (c *ProcessClosureProvidesAndRequiresTest) TestProvidedTypedef_noParentProvide_triplyNestedNamespace() {
	c.Test(
		c.Srcs("goog.provide('x.y.z'); /** @typedef {number} */ x.y.z;"),
		c.Expected(
			"/** @const */ var x = {}; /** @const */ x.y = {}; /** @const */  x.y.z = {}; x.y.z;"))
}

// Java: ProcessClosureProvidesAndRequiresTest.testProvideRequireSameFile()
func (c *ProcessClosureProvidesAndRequiresTest) TestProvideRequireSameFile() {
	c.Test(c.Srcs("goog.provide('x');\ngoog.require('x');"), c.Expected("/** @const */ var x = {};"))

	c.Test(c.Srcs("goog.provide('x');\ngoog.requireType('x');"), c.Expected("/** @const */ var x = {};"))
}

// Java: ProcessClosureProvidesAndRequiresTest.testSimpleProvidedNameCollection()
func (c *ProcessClosureProvidesAndRequiresTest) TestSimpleProvidedNameCollection() {
	providedNameMap :=
		c.GetProvidedNameCollection(
			`goog.provide('a.b');
goog.provide('a.b.c');
goog.provide('a.b.d');
`)

	AssertThat(providedNameMap.KeySlice()).ContainsExactly("goog", "a", "a.b", "a.b.c", "a.b.d")
}

// Java: ProcessClosureProvidesAndRequiresTest.testLegacyGoogModule()
func (c *ProcessClosureProvidesAndRequiresTest) TestLegacyGoogModule() {
	providedNameMap :=
		c.GetProvidedNameCollection(
			`goog.module('a.b.c');
goog.module.declareLegacyNamespace();

exports = class {};
`)

	AssertThat(providedNameMap.KeySlice()).ContainsExactly("goog", "a", "a.b", "a.b.c")
}

// Java: ProcessClosureProvidesAndRequiresTest.testLegacyGoogModule_inLoadModuleCall()
func (c *ProcessClosureProvidesAndRequiresTest) TestLegacyGoogModule_inLoadModuleCall() {
	providedNameMap :=
		c.GetProvidedNameCollection(
			`goog.loadModule(function(exports) {
  goog.module('a.b.c');
  goog.module.declareLegacyNamespace();
  exports = class {};
  return exports;
});
`)

	AssertThat(providedNameMap.KeySlice()).ContainsExactly("goog", "a", "a.b", "a.b.c")
}

// Java: ProcessClosureProvidesAndRequiresTest.testEsModule_ignored()
func (c *ProcessClosureProvidesAndRequiresTest) TestEsModule_ignored() {
	providedNameMap :=
		c.GetProvidedNameCollection("goog.declareModuleId('a.b.c'); export const x = 0;")

	AssertThat(providedNameMap.KeySlice()).ContainsExactly("goog")
}

// Java: ProcessClosureProvidesAndRequiresTest.testGoogExportSymbolRecording()
func (c *ProcessClosureProvidesAndRequiresTest) TestGoogExportSymbolRecording() {
	c.TestSame(
		`    goog.exportSymbol('a', 0);
    goog.exportSymbol('foo.bar.baz', 1);
    goog.exportSymbol(notAStringLiteral, 2);
    function f() {
      goog.exportSymbol('b', 3);
    }
`)
	AssertThat(c.lastProcessor.GetExportedVariableNames()).ContainsExactly("a", "foo", "b")
}
