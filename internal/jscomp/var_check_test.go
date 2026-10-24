// Port of com/google/javascript/jscomp/VarCheckTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/VarCheckTest.java)

/*
 * Copyright 2005 The Closure Compiler Authors.
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

type VarCheckTest struct {
	CompilerTestCaseBase
	strictChunkDepErrorLevel CheckLevel
	validityCheck            bool
	// A Java @Nullable CheckLevel: nil is null.
	externValidationErrorLevel *CheckLevel
	declarationCheck           bool
}

// Java: VarCheckTest()
func NewVarCheckTest() *VarCheckTest {
	c := &VarCheckTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, varCheckTest_EXTERNS)
	c.validityCheck = false
	return c
}

func TestVarCheckTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewVarCheckTest() })
}

// Java: VarCheckTest.EXTERNS
const varCheckTest_EXTERNS = "var window; function alert() {}"

// Java: VarCheckTest.setUp()
func (c *VarCheckTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	// Setup value set by individual tests to the appropriate defaults.
	c.AllowExternsChanges()
	c.strictChunkDepErrorLevel = CheckLevel_OFF
	c.externValidationErrorLevel = nil
	c.validityCheck = false
	c.declarationCheck = false
}

// Java: VarCheckTest.getOptions()
func (c *VarCheckTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetClosurePass(true)
	options.SetWarningLevel(DiagnosticGroups_MODULE_LOAD, CheckLevel_OFF)
	options.SetWarningLevel(DiagnosticGroups_STRICT_MODULE_DEP_CHECK, c.strictChunkDepErrorLevel)
	if c.externValidationErrorLevel != nil {
		options.SetWarningLevel(DiagnosticGroups_EXTERNS_VALIDATION, *c.externValidationErrorLevel)
	}
	return options
}

// Java: VarCheckTest.getProcessor(Compiler)
func (c *VarCheckTest) GetProcessor(compiler *Compiler) CompilerPass {
	return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
		NewVarCheckWithValidityCheck(compiler, c.validityCheck).Process(externs, root)
		if !c.validityCheck && !compiler.HasErrors() {
			// If the original test turned off sanity check, make sure our synthesized
			// code passes it.
			NewVarCheckWithValidityCheck(compiler, true).Process(externs, root)
		}
		if c.declarationCheck {
			NewVarCheckTest_VariableTestCheck(compiler).Process(externs, root)
		}
	})
}

// Java: VarCheckTest.testShorthandObjLit()
func (c *VarCheckTest) TestShorthandObjLit() {
	c.TestError("var x = {y};", VarCheck_UNDEFINED_VAR_ERROR)
	c.TestSame("var {x} = {x: 5}; let y = x;")

	c.TestError("var {...x} = {...y};", VarCheck_UNDEFINED_VAR_ERROR)
	c.TestSame("let y; var {...x} = {...y} = {};")
	c.TestSame("var {...x} = {x: 5}; let y = x;")
}

// Java: VarCheckTest.testBreak()
func (c *VarCheckTest) TestBreak() {
	c.TestSame("a: while(1) break a;")
}

// Java: VarCheckTest.testContinue()
func (c *VarCheckTest) TestContinue() {
	c.TestSame("a: while(1) continue a;")
}

// Java: VarCheckTest.testReferencedVarNotDefined()
func (c *VarCheckTest) TestReferencedVarNotDefined() {
	c.TestError("x = 0;", VarCheck_UNDEFINED_VAR_ERROR)
}

// Java: VarCheckTest.testReferencedVarNotDefined_typeof()
func (c *VarCheckTest) TestReferencedVarNotDefined_typeof() {
	// typeof x !== 'undefined' is the strict mode compliant way to check for a variable's
	// existence.
	c.TestSame("if (typeof undeclaredVariable !== 'undefined') {}")
	// Regression test: the exact pattern emitted by TypeScript to reference possibly declared
	// generic type variables in decorator emit.
	c.TestSame("var _a; typeof (_a = typeof Value !== \"undefined\" && Value) === \"function\"")
}

// Java: VarCheckTest.testReferencedVarNotDefined_arrowFunctionBody()
func (c *VarCheckTest) TestReferencedVarNotDefined_arrowFunctionBody() {
	c.TestError("() => y", VarCheck_UNDEFINED_VAR_ERROR)
}

// Java: VarCheckTest.testReferencedLetNotDefined()
func (c *VarCheckTest) TestReferencedLetNotDefined() {
	c.TestError("{ let x = 1; } var y = x;", VarCheck_UNDEFINED_VAR_ERROR)
}

// Java: VarCheckTest.testReferencedLetNotDefined_withES6Modules()
func (c *VarCheckTest) TestReferencedLetNotDefined_withES6Modules() {
	c.TestError("export function f() { { let x = 1; } var y = x; }", VarCheck_UNDEFINED_VAR_ERROR)
}

// Java: VarCheckTest.testReferencedLetDefined1()
func (c *VarCheckTest) TestReferencedLetDefined1() {
	c.TestSame("let x; x = 1;")
}

// Java: VarCheckTest.testNullishCoalesce()
func (c *VarCheckTest) TestNullishCoalesce() {
	c.TestSame("let x; x = 0 ?? true")
	c.TestError("let x; x = a ?? \"hi\"", VarCheck_UNDEFINED_VAR_ERROR)
}

// Java: VarCheckTest.testLogicalAssignment()
func (c *VarCheckTest) TestLogicalAssignment() {
	c.TestSame("let x; x ||= 2")
	c.TestError("let x; x ||= a", VarCheck_UNDEFINED_VAR_ERROR)
	// References let, not defined
	c.TestError("{ let x = 1; } var y; y ||= x;", VarCheck_UNDEFINED_VAR_ERROR)
	// Multiple declared vars
	c.TestSame("try { var x = 1; x ||= 2; } catch (x) {}")
	// In externs
	c.Test(c.Externs("x ||= new Klass();"), c.Srcs("class Klass{}"), c.Error(VarCheck_UNDEFINED_VAR_ERROR))
}

// Java: VarCheckTest.testReferencedLetDefined1_withES6Modules()
func (c *VarCheckTest) TestReferencedLetDefined1_withES6Modules() {
	c.TestSame("export let x; x = 1;")
}

// Java: VarCheckTest.testReferencedLetDefined2()
func (c *VarCheckTest) TestReferencedLetDefined2() {
	c.TestSame("let x; function y() {x = 1;}")
}

// Java: VarCheckTest.testReferencedConstDefined2()
func (c *VarCheckTest) TestReferencedConstDefined2() {
	c.TestSame("const x = 1; var y = x + 1;")
}

// Java: VarCheckTest.testReferencedVarDefined1()
func (c *VarCheckTest) TestReferencedVarDefined1() {
	c.TestSame("var x, y; x=1;")
}

// Java: VarCheckTest.testReferencedVarDefined2()
func (c *VarCheckTest) TestReferencedVarDefined2() {
	c.TestSame("var x; function y() {x=1;}")
}

// Java: VarCheckTest.testReferencedVarsExternallyDefined()
func (c *VarCheckTest) TestReferencedVarsExternallyDefined() {
	c.TestSame("var x = window; alert(x);")
}

// Java: VarCheckTest.testMultiplyDeclaredVars1()
func (c *VarCheckTest) TestMultiplyDeclaredVars1() {
	c.TestError("var x = 1; var x = 2;", VarCheck_VAR_MULTIPLY_DECLARED_ERROR)
}

// Java: VarCheckTest.testMultiplyDeclaredVars2()
func (c *VarCheckTest) TestMultiplyDeclaredVars2() {
	c.TestSame("var y; try { y=1 } catch (x) {} try { y=1 } catch (x) {}")
}

// Java: VarCheckTest.testMultiplyDeclaredVars3()
func (c *VarCheckTest) TestMultiplyDeclaredVars3() {
	c.TestSame("try { var x = 1; x *=2; } catch (x) {}")
}

// Java: VarCheckTest.testMultiplyDeclaredVars4()
func (c *VarCheckTest) TestMultiplyDeclaredVars4() {
	c.Test(c.Externs("x;"), c.Srcs("var x = 1; var x = 2;"), c.Error(VarCheck_VAR_MULTIPLY_DECLARED_ERROR))
}

// Java: VarCheckTest.testMultiplyDeclaredLets()
func (c *VarCheckTest) TestMultiplyDeclaredLets() {
	c.TestError("let x = 1; let x = 2;", VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR)
	c.TestError("let x = 1; var x = 2;", VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR)
	c.TestError("var x = 1; let x = 2;", VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR)
	c.TestError("let {x} = 1; let {x} = 2;", VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR)
	c.TestError("let {x} = 1; var {x} = 2;", VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR)
	c.TestError("var {x} = 1; let {x} = 2;", VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR)
}

// Java: VarCheckTest.testMultiplyDeclaredConsts()
func (c *VarCheckTest) TestMultiplyDeclaredConsts() {
	c.TestError("const x = 1; const x = 2;", VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR)
	c.TestError("const x = 1; var x = 2;", VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR)
	c.TestError("var x = 1; const x = 2;", VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR)
	c.TestError("const {x} = 1; const {x} = 2;", VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR)
	c.TestError("const {x} = 1; var {x} = 2;", VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR)
	c.TestError("var {x} = 1; const {x} = 2;", VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR)
}

// Java: VarCheckTest.testMultiplyDeclaredConsts_withES6Modules()
func (c *VarCheckTest) TestMultiplyDeclaredConsts_withES6Modules() {
	c.TestError(
		"export function f() { const x = 1; const x = 2; }",
		VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR)

	c.TestError("export const x = 1; export var x = 2;", VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR)

	c.TestError("export const a = 1, a = 2;", VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR)
}

// Java: VarCheckTest.testMultiplyDeclareLetsInDifferentScope()
func (c *VarCheckTest) TestMultiplyDeclareLetsInDifferentScope() {
	c.TestSame("let x = 1; if (123) {let x = 2;}")
	c.TestSame("try {let x = 1;} catch(x){}")
}

// Java: VarCheckTest.testMultiplyDeclaredCatchAndVar()
func (c *VarCheckTest) TestMultiplyDeclaredCatchAndVar() {
	// Note: This is technically valid code, but it's difficult to model the scoping rules in the
	// compiler and inconsistent across browsers. We forbid it in VariableReferenceCheck.
	c.TestSame("try {} catch (x) { var x = 1; }")
}

// Java: VarCheckTest.testMultiplyDeclaredCatchAndLet()
func (c *VarCheckTest) TestMultiplyDeclaredCatchAndLet() {
	c.TestError("try {} catch (x) { let x = 1; }", VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR)
}

// Java: VarCheckTest.testMultiplyDeclaredCatchAndConst()
func (c *VarCheckTest) TestMultiplyDeclaredCatchAndConst() {
	c.TestError("try {} catch (x) { const x = 1; }", VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR)
}

// Java: VarCheckTest.testReferencedVarDefinedClass()
func (c *VarCheckTest) TestReferencedVarDefinedClass() {
	c.TestError("var x; class x{ }", VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR)
	c.TestError("let x; class x{ }", VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR)
	c.TestError("const x = 1; class x{ }", VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR)
	c.TestError("class x{ } let x;", VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR)
	c.TestError("export default class x{ } let x;", VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR)
}

// Java: VarCheckTest.testMessagesPointToFirstOccurance()
func (c *VarCheckTest) TestMessagesPointToFirstOccurance() {
	c.Test(
		c.Srcs("var x = 1; var x = 2;"),
		c.Error(VarCheck_VAR_MULTIPLY_DECLARED_ERROR).
			WithMessageContaining("First occurrence: testcode:1:4"))
	c.Test(
		c.Srcs("var x; class x{ }"),
		c.Error(VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR).
			WithMessageContaining("First occurrence: testcode:1:4"))
}

// Java: VarCheckTest.testNamedClass()
func (c *VarCheckTest) TestNamedClass() {
	c.TestSame("class x {}")
	c.TestSame("var x = class x {};")
	c.TestSame("var y = class x {};")
	c.TestSame("var y = class x { foo() { return new x; } };")
	c.TestError("var Foo = class extends Bar {};", VarCheck_UNDEFINED_VAR_ERROR)
}

// Java: VarCheckTest.testVarReferenceInExterns()
func (c *VarCheckTest) TestVarReferenceInExterns() {
	c.TestSame(
		c.Externs("asdf;"),
		c.Srcs("var /** @suppress {duplicate} */ asdf;"),
		c.Warning(VarCheck_NAME_REFERENCE_IN_EXTERNS_ERROR))
}

// Java: VarCheckTest.testVarReferenceInExterns2()
func (c *VarCheckTest) TestVarReferenceInExterns2() {
	c.TestSame(
		c.Externs("asdf;"),
		c.Srcs(
			`(function() { var asdf; })()
var /** @suppress {duplicate} */ asdf;
`),
		c.Warning(VarCheck_NAME_REFERENCE_IN_EXTERNS_ERROR))
}

// Java: VarCheckTest.testNamespaceDeclarationInExterns()
func (c *VarCheckTest) TestNamespaceDeclarationInExterns() {
	c.TestSame(c.Externs("/** @const */ var $jscomp = $jscomp || {};"), c.Srcs(""))
}

// Java: VarCheckTest.testCallInExterns()
func (c *VarCheckTest) TestCallInExterns() {
	c.TestSame(
		c.Externs("yz();"),
		c.Srcs("function /** @suppress {duplicate} */ yz() {}"),
		c.Warning(VarCheck_NAME_REFERENCE_IN_EXTERNS_ERROR))
}

// Java: VarCheckTest.testDestructuringInExterns()
func (c *VarCheckTest) TestDestructuringInExterns() {
	c.TestSame(c.Externs("function externalFunction({x, y}) {}"), c.Srcs(""))
	c.TestSame(c.Externs("function externalFunction({x, y:{z}}) {}"), c.Srcs(""))
	c.TestSame(c.Externs("function externalFunction({x:localName}) {}"), c.Srcs(""))
	c.TestSame(c.Externs("function externalFunction([a, b, c]) {}"), c.Srcs(""))
	c.TestSame(c.Externs("function externalFunction([[...a], b, c = 5, ...d]) {}"), c.Srcs(""))
}

// Java: VarCheckTest.testVarReferenceInExterns_withEs6Modules()
func (c *VarCheckTest) TestVarReferenceInExterns_withEs6Modules() {
	// vars in ES6 modules are not in global scope, so foo is undefined.
	c.TestError(c.Externs("foo;"), c.Srcs("export var foo;"), c.Error(VarCheck_UNDEFINED_VAR_ERROR))
}

// Java: VarCheckTest.testVarDeclarationInExterns()
func (c *VarCheckTest) TestVarDeclarationInExterns() {
	c.TestSame(c.Externs("var asdf;"), c.Srcs("asdf;"))
}

// Java: VarCheckTest.testVarReferenceInTypeSummary()
func (c *VarCheckTest) TestVarReferenceInTypeSummary() {
	c.TestSame(
		c.Externs(
			`/** @typeSummary */
var goog;
goog.addSingletonGetter;
class Foo {}
goog.addSingletonGetter(Foo);
`),
		c.Srcs("Foo.getInstance();"))
}

// Java: VarCheckTest.testFunctionDeclarationInExterns()
func (c *VarCheckTest) TestFunctionDeclarationInExterns() {
	c.TestSame(c.Externs("function foo(x = 7) {}"), c.Srcs("foo();"))
	c.TestSame(c.Externs("function foo(...rest) {}"), c.Srcs("foo(1,2,3);"))
}

// Java: VarCheckTest.testVarAssignmentInExterns()
func (c *VarCheckTest) TestVarAssignmentInExterns() {
	c.TestSame(c.Externs("/** @type{{foo:string}} */ var foo; var asdf = foo;"), c.Srcs("asdf.foo;"))
}

// Java: VarCheckTest.testAliasesInExterns()
func (c *VarCheckTest) TestAliasesInExterns() {
	c.externValidationErrorLevel = CheckLevel_ERROR.Ptr()

	c.TestSame(c.Externs("var foo; /** @const */ var asdf = foo;"), c.Srcs(""))
	c.TestSame(c.Externs("var Foo; var ns = {}; /** @const */ ns.FooAlias = Foo;"), c.Srcs(""))
	c.TestSame(
		c.Externs(
			`var ns = {}; /** @constructor */ ns.Foo = function() {};
var ns2 = {}; /** @const */ ns2.Bar = ns.Foo;
`),
		c.Srcs(""))
}

// Java: VarCheckTest.testDuplicateNamespaceInExterns()
func (c *VarCheckTest) TestDuplicateNamespaceInExterns() {
	// Compare as tree does not allow multiple externs, but VarCheck forcibly inserts them.
	c.DisableCompareAsTree()
	c.TestExternChanges(
		c.Externs("/** @const */ var ns = {}; /** @const */ var ns = {};"),
		c.Srcs(""),
		c.Expected(CompilerTestCase_VAR_CHECK_EXTERNS+"/** @const */ var ns = {};"))
}

// Java: VarCheckTest.testLetDeclarationInExterns()
func (c *VarCheckTest) TestLetDeclarationInExterns() {
	c.TestSame(c.Externs("let asdf;"), c.Srcs("asdf;"))
}

// Java: VarCheckTest.testConstDeclarationInExterns()
func (c *VarCheckTest) TestConstDeclarationInExterns() {
	c.TestSame(c.Externs("const asdf = 1;"), c.Srcs("asdf;"))
}

// Java: VarCheckTest.testNewInExterns()
func (c *VarCheckTest) TestNewInExterns() {
	// Class is not hoisted.
	c.Test(c.Externs("x = new Klass();"), c.Srcs("class Klass{}"), c.Error(VarCheck_UNDEFINED_VAR_ERROR))
}

// Java: VarCheckTest.testPropReferenceInExterns1()
func (c *VarCheckTest) TestPropReferenceInExterns1() {
	c.TestSame(
		c.Externs("asdf.foo;"),
		c.Srcs("var /** @suppress {duplicate} */ asdf;"),
		c.Warning(VarCheck_UNDEFINED_EXTERN_VAR_ERROR))
}

// Java: VarCheckTest.testPropReferenceInExterns2()
func (c *VarCheckTest) TestPropReferenceInExterns2() {
	c.Test(c.Externs("asdf.foo;"), c.Srcs(""), c.Error(VarCheck_UNDEFINED_VAR_ERROR))
}

// Java: VarCheckTest.testPropReferenceInExterns3()
func (c *VarCheckTest) TestPropReferenceInExterns3() {
	c.TestSame(
		c.Externs("asdf.foo;"),
		c.Srcs("var /** @suppress {duplicate} */ asdf;"),
		c.Warning(VarCheck_UNDEFINED_EXTERN_VAR_ERROR))

	c.externValidationErrorLevel = CheckLevel_ERROR.Ptr()
	c.Test(c.Externs("asdf.foo;"), c.Srcs("var asdf;"), c.Error(VarCheck_UNDEFINED_EXTERN_VAR_ERROR))

	c.externValidationErrorLevel = CheckLevel_OFF.Ptr()
	c.TestSame(c.Externs("asdf.foo;"), c.Srcs("var asdf;"))
}

// Java: VarCheckTest.testPropReferenceInExterns4()
func (c *VarCheckTest) TestPropReferenceInExterns4() {
	c.TestSame(c.Externs("asdf.foo;"), c.Srcs("let asdf;"), c.Warning(VarCheck_UNDEFINED_EXTERN_VAR_ERROR))
}

// Java: VarCheckTest.testPropReferenceInExterns5()
func (c *VarCheckTest) TestPropReferenceInExterns5() {
	c.TestSame(
		c.Externs("asdf.foo;"), c.Srcs("class asdf {}"), c.Warning(VarCheck_UNDEFINED_EXTERN_VAR_ERROR))
}

// Java: VarCheckTest.testVarInWithBlock()
func (c *VarCheckTest) TestVarInWithBlock() {
	c.TestError("var a = {b:5}; with (a){b;}", VarCheck_UNDEFINED_VAR_ERROR)
}

// Java: VarCheckTest.testFunctionDeclaredInBlock()
func (c *VarCheckTest) TestFunctionDeclaredInBlock() {
	c.TestError("if (true) {function foo() {}} foo();", VarCheck_UNDEFINED_VAR_ERROR)
	c.TestError("foo(); if (true) {function foo() {}}", VarCheck_UNDEFINED_VAR_ERROR)

	c.TestSame("if (true) {var foo = ()=>{}} foo();")
	c.TestError("if (true) {let foo = ()=>{}} foo();", VarCheck_UNDEFINED_VAR_ERROR)
	c.TestError("if (true) {const foo = ()=>{}} foo();", VarCheck_UNDEFINED_VAR_ERROR)

	c.TestSame("foo(); if (true) {var foo = ()=>{}}")
	c.TestError("foo(); if (true) {let foo = ()=>{}}", VarCheck_UNDEFINED_VAR_ERROR)
	c.TestError("foo(); if (true) {const foo = ()=>{}}", VarCheck_UNDEFINED_VAR_ERROR)
}

// Java: VarCheckTest.testValidFunctionExpr()
func (c *VarCheckTest) TestValidFunctionExpr() {
	c.TestSame("(function() {});")
}

// Java: VarCheckTest.testRecursiveFunction()
func (c *VarCheckTest) TestRecursiveFunction() {
	c.TestSame("(function a() { return a(); })();")
}

// Java: VarCheckTest.testRecursiveFunction2()
func (c *VarCheckTest) TestRecursiveFunction2() {
	c.TestSame("var a = 3; (function a() { return a(); })();")
}

// Java: VarCheckTest.testParam()
func (c *VarCheckTest) TestParam() {
	c.TestSame("function fn(a){ var b = a; }")
	c.TestSame("function fn(a){ var a = 2; }")
	c.TestError("function fn(){ var b = a; }", VarCheck_UNDEFINED_VAR_ERROR)

	// Default parameters
	c.TestError("function fn(a = b) { function g(a = 3) { var b; } }", VarCheck_UNDEFINED_VAR_ERROR)
	c.TestError("function f(x=a) { let a; }", VarCheck_UNDEFINED_VAR_ERROR)
	c.TestError("function f(x=a) { { let a; } }", VarCheck_UNDEFINED_VAR_ERROR)
	c.TestError("function f(x=b) { function a(x=1) { var b; } }", VarCheck_UNDEFINED_VAR_ERROR)
	c.TestError("function f(x=a) { var a; }", VarCheck_UNDEFINED_VAR_ERROR)
	c.TestError("function f(x=a()) { function a() {} }", VarCheck_UNDEFINED_VAR_ERROR)
	c.TestError("function f(x=[a]) { var a; }", VarCheck_UNDEFINED_VAR_ERROR)
	c.TestError("function f(x = new foo.bar()) {}", VarCheck_UNDEFINED_VAR_ERROR)
	c.TestSame("var foo = {}; foo.bar = class {}; function f(x = new foo.bar()) {}")

	c.TestSame("function fn(a = 2){ var b = a; }")
	c.TestSame("function fn(a = 2){ var a = 3; }")
	c.TestSame("function fn({a, b}){ var c = a; }")
	c.TestSame("function fn({a, b}){ var a = 3; }")
}

// Java: VarCheckTest.testParamArrowFunction()
func (c *VarCheckTest) TestParamArrowFunction() {
	c.TestSame("(a) => { var b = a; }")
	c.TestError("() => { var b = a; }", VarCheck_UNDEFINED_VAR_ERROR)
	c.TestError("(x=a) => { let a; }", VarCheck_UNDEFINED_VAR_ERROR)

	// Arrow function nested
	c.TestError(
		`function FUNC() {
  {
    () => { var b = a; }
  }
}
`,
		VarCheck_UNDEFINED_VAR_ERROR)
}

// Java: VarCheckTest.testLegalVarReferenceBetweenChunks()
func (c *VarCheckTest) TestLegalVarReferenceBetweenChunks() {
	c.TestDependentChunks("var x = 10;", "var y = x++;", nil)
}

// Java: VarCheckTest.testLegalLetReferenceBetweenChunks()
func (c *VarCheckTest) TestLegalLetReferenceBetweenChunks() {
	c.TestDependentChunks("let x = 10;", "let y = x++;", nil)
}

// Java: VarCheckTest.testLegalConstReferenceBetweenChunks()
func (c *VarCheckTest) TestLegalConstReferenceBetweenChunks() {
	c.TestDependentChunks("const x = 10;", "const y = x + 1;", nil)
}

// Java: VarCheckTest.testMissingChunkDependencyDefault()
func (c *VarCheckTest) TestMissingChunkDependencyDefault() {
	c.TestIndependentChunks("var x = 10;", "var y = x++;", nil, VarCheck_MISSING_CHUNK_DEP_ERROR)
}

// Java: VarCheckTest.testMissingChunkDependencyLetAndConst()
func (c *VarCheckTest) TestMissingChunkDependencyLetAndConst() {
	c.TestIndependentChunks("let x = 10;", "let y = x++;", nil, VarCheck_MISSING_CHUNK_DEP_ERROR)
	c.TestIndependentChunks(
		"const x = 10;", "const y = x + 1;", nil, VarCheck_MISSING_CHUNK_DEP_ERROR)
}

// Java: VarCheckTest.testViolatedChunkDependencyDefault()
func (c *VarCheckTest) TestViolatedChunkDependencyDefault() {
	c.TestDependentChunks("var y = x++;", "var x = 10;", VarCheck_VIOLATED_CHUNK_DEP_ERROR)
}

// Java: VarCheckTest.testViolatedChunkDependencyLetAndConst()
func (c *VarCheckTest) TestViolatedChunkDependencyLetAndConst() {
	c.TestDependentChunks("let y = x++;", "let x = 10;", VarCheck_VIOLATED_CHUNK_DEP_ERROR)
	c.TestDependentChunks("const y = x + 1;", "const x = 10;", VarCheck_VIOLATED_CHUNK_DEP_ERROR)
}

// Java: VarCheckTest.testMissingChunkDependencySkipNonStrict()
func (c *VarCheckTest) TestMissingChunkDependencySkipNonStrict() {
	c.validityCheck = true
	c.TestIndependentChunks("var x = 10;", "var y = x++;", nil, nil)
}

// Java: VarCheckTest.testViolatedChunkDependencySkipNonStrict()
func (c *VarCheckTest) TestViolatedChunkDependencySkipNonStrict() {
	c.validityCheck = true
	c.TestDependentChunks("var y = x++;", "var x = 10;", nil)
}

// Java: VarCheckTest.testMissingChunkDependencySkipNonStrictNotPromoted()
func (c *VarCheckTest) TestMissingChunkDependencySkipNonStrictNotPromoted() {
	c.validityCheck = true
	c.strictChunkDepErrorLevel = CheckLevel_ERROR
	c.TestIndependentChunks("var x = 10;", "var y = x++;", nil, nil)
}

// Java: VarCheckTest.testViolatedChunkDependencyNonStrictNotPromoted()
func (c *VarCheckTest) TestViolatedChunkDependencyNonStrictNotPromoted() {
	c.validityCheck = true
	c.strictChunkDepErrorLevel = CheckLevel_ERROR
	c.TestDependentChunks("var y = x++;", "var x = 10;", nil)
}

// Java: VarCheckTest.testDependentStrictChunkDependencyCheck()
func (c *VarCheckTest) TestDependentStrictChunkDependencyCheck() {
	c.strictChunkDepErrorLevel = CheckLevel_ERROR
	c.TestDependentChunks(
		"var f = function() {return new B();};",
		"var B = function() {}",
		VarCheck_STRICT_CHUNK_DEP_ERROR)
}

// Java: VarCheckTest.testIndependentStrictChunkDependencyCheck()
func (c *VarCheckTest) TestIndependentStrictChunkDependencyCheck() {
	c.strictChunkDepErrorLevel = CheckLevel_ERROR
	c.TestIndependentChunks(
		"var f = function() {return new B();};",
		"var B = function() {}",
		VarCheck_STRICT_CHUNK_DEP_ERROR,
		nil)
}

// Java: VarCheckTest.testStarStrictChunkDependencyCheck()
func (c *VarCheckTest) TestStarStrictChunkDependencyCheck() {
	c.strictChunkDepErrorLevel = CheckLevel_WARNING
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk("function a() {}").
				AddChunk("function b() { a(); c(); }").
				AddChunk("function c() { a(); }").
				Build()),
		c.Warning(VarCheck_STRICT_CHUNK_DEP_ERROR))
}

// Java: VarCheckTest.testForwardVarReferenceInLocalScope1()
func (c *VarCheckTest) TestForwardVarReferenceInLocalScope1() {
	c.TestDependentChunks("var x = 10; function a() {y++;}", "var y = 11; a();", nil)
}

// Java: VarCheckTest.testForwardVarReferenceInLocalScope2()
func (c *VarCheckTest) TestForwardVarReferenceInLocalScope2() {
	// It would be nice if this pass could use a call graph to flag this case
	// as an error, but it currently doesn't.
	c.TestDependentChunks("var x = 10; function a() {y++;} a();", "var y = 11;", nil)
}

// Java: VarCheckTest.testDependentChunks(String, String, DiagnosticType)
func (c *VarCheckTest) TestDependentChunks(code1 string, code2 string, error *DiagnosticType) {
	c.TestDependentChunksWithWarning(code1, code2, error, nil)
}

// Java: VarCheckTest.testDependentChunks(String, String, DiagnosticType, DiagnosticType)
func (c *VarCheckTest) TestDependentChunksWithWarning(code1 string, code2 string, error *DiagnosticType, warning *DiagnosticType) {
	c.TestTwoChunks(code1, code2, true, error, warning)
}

// Java: VarCheckTest.testIndependentChunks(String, String, DiagnosticType, DiagnosticType)
func (c *VarCheckTest) TestIndependentChunks(code1 string, code2 string, error *DiagnosticType, warning *DiagnosticType) {
	c.TestTwoChunks(code1, code2, false, error, warning)
}

// Java: VarCheckTest.testTwoChunks(String, String, boolean, DiagnosticType, DiagnosticType)
func (c *VarCheckTest) TestTwoChunks(code1 string, code2 string, m2DependsOnm1 bool, error *DiagnosticType, warning *DiagnosticType) {
	m1 := NewJSChunk("m1")
	m1.AddSourceFile(SourceFile_FromCode("input1", code1))
	m2 := NewJSChunk("m2")
	m2.AddSourceFile(SourceFile_FromCode("input2", code2))
	if m2DependsOnm1 {
		m2.AddDependency(m1)
	}
	if error == nil && warning == nil {
		c.Test(c.Srcs(m1, m2), c.Expected(code1, code2))
	} else if error == nil {
		c.Test(c.Srcs(m1, m2), c.Expected(code1, code2), c.Warning(warning))
	} else {
		c.TestError(c.Srcs(m1, m2), c.Error(error))
	}
}

//////////////////////////////////////////////////////////////////////////////
// Test synthesis of externs

// Java: VarCheckTest.testSimple()
func (c *VarCheckTest) TestSimple() {
	c.CheckSynthesizedExtern(c.Srcs("x"), "var x;"+CompilerTestCase_VAR_CHECK_EXTERNS)
	c.CheckSynthesizedExtern(c.Srcs("var x"), CompilerTestCase_VAR_CHECK_EXTERNS)
}

// Java: VarCheckTest.testSimpleValidityCheck()
func (c *VarCheckTest) TestSimpleValidityCheck() {
	c.validityCheck = true
	jutil.Try(func() {
		c.CheckSynthesizedExtern(c.Srcs("x"), "")
		AssertWithMessage("Expected RuntimeException").Fail()
	}, func(e jutil.IRuntimeException) {
		AssertThat(e.GetMessage()).Contains("Unexpected variable x")
	})
}

// Java: VarCheckTest.testParameter()
func (c *VarCheckTest) TestParameter() {
	c.CheckSynthesizedExtern(c.Srcs("function f(x){}"), CompilerTestCase_VAR_CHECK_EXTERNS)
}

// Java: VarCheckTest.testLocalVar()
func (c *VarCheckTest) TestLocalVar() {
	c.CheckSynthesizedExtern(c.Srcs("function f(){x}"), "var x;"+CompilerTestCase_VAR_CHECK_EXTERNS)
}

// Java: VarCheckTest.testTwoLocalVars()
func (c *VarCheckTest) TestTwoLocalVars() {
	c.CheckSynthesizedExtern(c.Srcs("function f(){x}function g() {x}"), "var x;"+CompilerTestCase_VAR_CHECK_EXTERNS)
}

// Java: VarCheckTest.testInnerFunctionLocalVar()
func (c *VarCheckTest) TestInnerFunctionLocalVar() {
	c.CheckSynthesizedExtern(c.Srcs("function f(){function g() {x}}"), "var x;"+CompilerTestCase_VAR_CHECK_EXTERNS)
}

// Java: VarCheckTest.testNoCreateVarsForLabels()
func (c *VarCheckTest) TestNoCreateVarsForLabels() {
	c.CheckSynthesizedExtern(c.Srcs("x:var y"), CompilerTestCase_VAR_CHECK_EXTERNS)
}

// Java: VarCheckTest.testJSCompiler_renamePropertyInRequiredSymbols()
func (c *VarCheckTest) TestJSCompiler_renamePropertyInRequiredSymbols() {
	AssertThat(VarCheck_REQUIRED_SYMBOLS).Contains("JSCompiler_renameProperty")
}

// Java: VarCheckTest.testJSCompiler_renamePropertySynthesizedExtern()
func (c *VarCheckTest) TestJSCompiler_renamePropertySynthesizedExtern() {
	c.CheckSynthesizedExtern(c.Srcs(""), CompilerTestCase_VAR_CHECK_EXTERNS)
	AssertThat(CompilerTestCase_VAR_CHECK_EXTERNS).Contains("var JSCompiler_renameProperty;\n")
}

// Java: VarCheckTest.testJSCompiler_renamePropertySynthesizedExternAllowsLaterPasses()
func (c *VarCheckTest) TestJSCompiler_renamePropertySynthesizedExternAllowsLaterPasses() {
	compiler := c.CreateAndInitializeCompiler(c.Externs(""), c.Srcs("var placeholder = 1;"))
	compiler.ParseInputs()
	externs := compiler.GetExternsRoot()
	root := compiler.GetJsRoot()
	NewVarCheckWithValidityCheck(compiler /* validityCheck= */, false).Process(externs, root)
	// Simulate a pass like ClosureOptimizePrimitives or DartSuperAccessorsPass adding a call
	// to JSCompiler_renameProperty without an explicit extern declaration.
	syntheticCall :=
		rhino.IR_ExprResult(rhino.IR_Call(rhino.IR_Name("JSCompiler_renameProperty"), rhino.IR_String("foo")))
	root.GetFirstChild().AddChildToBack(syntheticCall)
	// ValidityCheck runs VarCheck with validityCheck = true, which must pass without throwing.
	NewVarCheckWithValidityCheck(compiler /* validityCheck= */, true).Process(externs, root)
}

// Java: VarCheckTest.testVariableInNormalCodeUsedInExterns1()
func (c *VarCheckTest) TestVariableInNormalCodeUsedInExterns1() {
	c.CheckSynthesizedExternWithExterns(
		c.Externs("x.foo;"),
		c.Srcs("var x;"),
		"var x; "+CompilerTestCase_VAR_CHECK_EXTERNS+" x.foo;",
		c.Warning(VarCheck_UNDEFINED_EXTERN_VAR_ERROR))
}

// Java: VarCheckTest.testVariableInNormalCodeUsedInExterns2()
func (c *VarCheckTest) TestVariableInNormalCodeUsedInExterns2() {
	c.CheckSynthesizedExternWithExterns(
		c.Externs("x;"),
		c.Srcs("var x;"),
		"var x; "+CompilerTestCase_VAR_CHECK_EXTERNS+" x;",
		c.Warning(VarCheck_NAME_REFERENCE_IN_EXTERNS_ERROR))
}

// Java: VarCheckTest.testVariableInNormalCodeUsedInExterns3()
func (c *VarCheckTest) TestVariableInNormalCodeUsedInExterns3() {
	c.CheckSynthesizedExternWithExterns(
		c.Externs("x.foo;"),
		c.Srcs("function x() {}"),
		"var x; "+CompilerTestCase_VAR_CHECK_EXTERNS+" x.foo;",
		c.Warning(VarCheck_UNDEFINED_EXTERN_VAR_ERROR))
}

// Java: VarCheckTest.testVariableInNormalCodeUsedInExterns4()
func (c *VarCheckTest) TestVariableInNormalCodeUsedInExterns4() {
	c.CheckSynthesizedExternWithExterns(
		c.Externs("x;"),
		c.Srcs("function x() {}"),
		"var x; "+CompilerTestCase_VAR_CHECK_EXTERNS+" x;",
		c.Warning(VarCheck_NAME_REFERENCE_IN_EXTERNS_ERROR))
}

// Java: VarCheckTest.testRedeclaration1()
func (c *VarCheckTest) TestRedeclaration1() {
	js := "var a; var a;"
	c.TestError(js, VarCheck_VAR_MULTIPLY_DECLARED_ERROR)
}

// Java: VarCheckTest.testRedeclaration2()
func (c *VarCheckTest) TestRedeclaration2() {
	js := "var a; /** @suppress {duplicate} */ var a;"
	c.TestSame(js)
}

// Java: VarCheckTest.testRedeclaration3()
func (c *VarCheckTest) TestRedeclaration3() {
	js := " /** @suppress {duplicate} */ var a; var a; "
	c.TestSame(js)
}

// Java: VarCheckTest.testRedeclaration4()
func (c *VarCheckTest) TestRedeclaration4() {
	js :=
		`/** @fileoverview @suppress {duplicate} */
/** @type {string} */ var a;
var a; 
`
	c.TestSame(js)
}

// Java: VarCheckTest.testSuppressionWithInlineJsDoc()
func (c *VarCheckTest) TestSuppressionWithInlineJsDoc() {
	c.TestSame("/** @suppress {duplicate} */ var /** number */ a; var a;")
}

// Java: VarCheckTest.testDuplicateVar()
func (c *VarCheckTest) TestDuplicateVar() {
	c.TestError(
		"/** @define {boolean} */ var DEF = false; var DEF = true;", VarCheck_VAR_MULTIPLY_DECLARED_ERROR)
}

// Java: VarCheckTest.testDontAllowSuppressDupeOnLet()
func (c *VarCheckTest) TestDontAllowSuppressDupeOnLet() {
	c.TestError(
		"let a; /** @suppress {duplicate} */ let a; ", VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR)

	c.TestError(
		"function f() { let a; /** @suppress {duplicate} */ let a; }",
		VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR)
}

// Java: VarCheckTest.testDuplicateBlockScopedDeclarationInSwitch()
func (c *VarCheckTest) TestDuplicateBlockScopedDeclarationInSwitch() {
	c.TestError(
		`function f(x) {
  switch (x) {
    case 'a':
      let z = 123;
      break;
    case 'b':
      let z = 234;
      break;
  }
}
`,
		VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR)

	c.TestError(
		`function f(x) {
  switch (x) {
    case 'a':
      class C {}
      break;
    case 'b':
      class C {}
      break;
  }
}
`,
		VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR)
}

// Java: VarCheckTest.testLetConstRedeclareWithFunctions_withEs6Modules()
func (c *VarCheckTest) TestLetConstRedeclareWithFunctions_withEs6Modules() {
	c.TestError("function f() {} let f = 1; export {f};", VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR)
	c.TestError("let f = 1; function f() {}  export {f};", VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR)
	c.TestError(
		"const f = 1; function f() {} export {f};", VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR)
	c.TestError(
		"function f() {} const f = 1;  export {f};", VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR)
	c.TestError(
		"export default function f() {}; let f = 5;", VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR)
}

// Java: VarCheckTest.testFunctionScopeArguments()
func (c *VarCheckTest) TestFunctionScopeArguments() {
	// A var declaration doesn't mask arguments
	c.TestSame("function f() {var arguments}")

	c.TestError("var f = function arguments() {}", VarCheck_VAR_ARGUMENTS_SHADOWED_ERROR)
	c.TestError("var f = function (arguments) {}", VarCheck_VAR_ARGUMENTS_SHADOWED_ERROR)
	c.TestSame("function f() {try {} catch(arguments) {}}")

	c.validityCheck = true
	c.TestSame("function f() {var arguments}")
}

// Java: VarCheckTest.testFunctionRedeclared_global()
func (c *VarCheckTest) TestFunctionRedeclared_global() {
	// Redeclaration in global scope is allowed.
	c.TestSame("/** @fileoverview @suppress {duplicate} */ function f() {};function f() {};")
}

// Java: VarCheckTest.testFunctionRedeclared1()
func (c *VarCheckTest) TestFunctionRedeclared1() {
	c.TestError("{ function f() {}; function f() {}; }", VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR)
	c.TestError(
		"if (0) { function f() {}; function f() {}; }", VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR)
	c.TestError(
		"try { function f() {}; function f() {}; } catch (e) {}",
		VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR)
}

// Java: VarCheckTest.testFunctionRedeclared2()
func (c *VarCheckTest) TestFunctionRedeclared2() {
	c.TestError("{ let f = 0; function f() {}; }", VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR)
	c.TestError(
		"if (0) { const f = 1; function f() {}; }", VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR)
	c.TestError(
		"try { class f {}; function f() {}; } catch (e) {}",
		VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR)
}

// Java: VarCheckTest.testFunctionRedeclared3()
func (c *VarCheckTest) TestFunctionRedeclared3() {
	c.TestError("{ function f() {}; let f = 0; }", VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR)
	c.TestError("if (0) { function f() {}; const f = 0;}", VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR)
	c.TestError(
		"try { function f() {}; class f {}; } catch (e) {}",
		VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR)
}

// Java: VarCheckTest.testUndeclaredVarForGoog()
func (c *VarCheckTest) TestUndeclaredVarForGoog() {
	c.EnableClosurePass()
	c.EnableRewriteClosureCode()
	// We don't want to get goog as an undeclared var here.
	c.Test(
		c.Srcs("goog.require('namespace.Class1');\n"),
		c.Error(ClosurePrimitiveErrors_MISSING_MODULE_OR_PROVIDE),
		c.Error(VarCheck_UNDEFINED_VAR_ERROR))
}

// Java: VarCheckTest.testImportedNameCollision()
func (c *VarCheckTest) TestImportedNameCollision() {
	// TODO(tbreisacher): This should throw a duplicate declaration error.
	c.TestSame("import foo from './foo'; foo('hello'); var foo = 5;")
}

// Java: VarCheckTest.testImportStar()
func (c *VarCheckTest) TestImportStar() {
	c.TestSame("import * as foo from './foo.js';")
}

// Java: VarCheckTest.testExportAsAlias()
func (c *VarCheckTest) TestExportAsAlias() {
	c.TestSame("let a = 1; export {a as b};")
	c.TestError("let a = 1; export {b as a};", VarCheck_UNDEFINED_VAR_ERROR)
	c.TestError("export {a as a};", VarCheck_UNDEFINED_VAR_ERROR)

	// Make sure non-aliased exports still work correctly.
	c.TestSame("let a = 1; export {a}")
	c.TestError("let a = 1; export {b};", VarCheck_UNDEFINED_VAR_ERROR)
}

// Java: VarCheckTest.testImportAsAlias()
func (c *VarCheckTest) TestImportAsAlias() {
	c.TestSame("import {b as a} from './foo.js'; let c = a;")
	c.TestError("import {b as a} from './foo.js'; let c = b;", VarCheck_UNDEFINED_VAR_ERROR)
	c.TestSame("import {a} from './foo.js'; let c = a;")
}

// Java: VarCheckTest.testComputedPropertyWithNamedFunction()
func (c *VarCheckTest) TestComputedPropertyWithNamedFunction() {
	c.TestSame("({[0]: function f() {}})")
}

// Java: VarCheckTest.testEsModule_withUndefinedExportsRef()
func (c *VarCheckTest) TestEsModule_withUndefinedExportsRef() {
	c.TestError("exports = function() {}; export {exports};", VarCheck_UNDEFINED_VAR_ERROR)
}

// Java: VarCheckTest.testGoogModule_withDefaultExports()
func (c *VarCheckTest) TestGoogModule_withDefaultExports() {
	c.TestSame(
		c.Srcs(
			testExternsBuilder_GetClosureExternsAsSource(),
			"goog.module('a.b'); exports = function() {};"))
}

// Java: VarCheckTest.testGoogModule_withExportsRefInFunction()
func (c *VarCheckTest) TestGoogModule_withExportsRefInFunction() {
	c.TestSame(
		c.Srcs(
			// Referencing 'exports' inside a function is strange but allowed.
			testExternsBuilder_GetClosureExternsAsSource(),
			"goog.module('a.b'); exports.f = function() { exports.f(); };"))
}

// Java: VarCheckTest.testGoogModule_withNamedExports()
func (c *VarCheckTest) TestGoogModule_withNamedExports() {
	c.TestSame(
		c.Srcs(
			testExternsBuilder_GetClosureExternsAsSource(),
			"goog.module('a.b'); exports.f = function() {}; exports.x = 0;"))
}

// Java: VarCheckTest.testGoogProvide_simpleName()
func (c *VarCheckTest) TestGoogProvide_simpleName() {
	c.TestSame(
		testExternsBuilder_GetClosureExternsAsSource() + "goog.provide('A'); var A = class {};")
	c.TestSame(testExternsBuilder_GetClosureExternsAsSource() + "goog.provide('A'); A = class {};")
}

// Java: VarCheckTest.testGoogProvide_simpleName_earlyReference()
func (c *VarCheckTest) TestGoogProvide_simpleName_earlyReference() {
	c.TestSame(testExternsBuilder_GetClosureExternsAsSource() + "A.B = class {}; goog.provide('A');")
	c.TestSame(
		testExternsBuilder_GetClosureExternsAsSource() +
			"A.B = class {}; goog.provide('A'); var A = class {};")
}

// Java: VarCheckTest.testGoogProvide_multipleRootsInSameFile()
func (c *VarCheckTest) TestGoogProvide_multipleRootsInSameFile() {
	c.TestSame(
		testExternsBuilder_GetClosureExternsAsSource() +
			"goog.provide('A'); goog.provide('B'); var A = class {}; var B = class {};")
	c.TestSame(
		testExternsBuilder_GetClosureExternsAsSource() +
			"goog.provide('A.a'); goog.provide('B.b'); A.a = 0; B.b = 1")
}

// Java: VarCheckTest.testGoogProvide_complexName()
func (c *VarCheckTest) TestGoogProvide_complexName() {
	c.TestSame(
		testExternsBuilder_GetClosureExternsAsSource() +
			"goog.provide('foo.A'); foo.A = class {};")
	c.TestSame(
		testExternsBuilder_GetClosureExternsAsSource() +
			"goog.provide('foo.bar.A'); foo.bar.A = class {};")
}

// Java: VarCheckTest.testGoogLegacyModule()
func (c *VarCheckTest) TestGoogLegacyModule() {
	c.TestSame(
		c.Srcs(
			testExternsBuilder_GetClosureExternsAsSource(),
			"goog.module('foo.A'); goog.module.declareLegacyNamespace(); exports = class {};",
			"new foo.A();"))
}

// Java: VarCheckTest.testGoogNonLegacyModule()
func (c *VarCheckTest) TestGoogNonLegacyModule() {
	c.TestError(
		c.Srcs(testExternsBuilder_GetClosureExternsAsSource(), "goog.module('foo.A');", "foo.A();"),
		c.Error(VarCheck_UNDEFINED_VAR_ERROR))

	c.TestError(
		c.Srcs(
			testExternsBuilder_GetClosureExternsAsSource(),
			"goog.module('foo.A'); exports = class {};",
			"new foo.A();"),
		c.Error(VarCheck_UNDEFINED_VAR_ERROR))
}

// Java: VarCheckTest.testGoogLegacyModule_inLoadModule()
func (c *VarCheckTest) TestGoogLegacyModule_inLoadModule() {
	c.TestSame(
		c.Srcs(
			testExternsBuilder_GetClosureExternsAsSource(),
			`goog.loadModule(function(exports) {
  goog.module('foo.A');
  goog.module.declareLegacyNamespace();
  exports = class {};
  return exports;
});
`,
			"new foo.A();"))
}

// Java: VarCheckTest.testGoogNonLegacyModule_inLoadModule()
func (c *VarCheckTest) TestGoogNonLegacyModule_inLoadModule() {
	c.TestError(
		c.Srcs(
			testExternsBuilder_GetClosureExternsAsSource(),
			`goog.loadModule(function(exports) {
  goog.module('foo.A');
  exports = class {};
  return exports;
});
`,
			"new foo.A();"),
		c.Error(VarCheck_UNDEFINED_VAR_ERROR))
}

// Java: VarCheckTest.testGoogProvide_externs_addsDeclForNamespace()
func (c *VarCheckTest) TestGoogProvide_externs_addsDeclForNamespace() {
	c.CheckSynthesizedExternWithExterns(
		c.Externs("var goog; goog.provide('a.b'); a.b.C = class {};"),
		c.Srcs(""),
		CompilerTestCase_VAR_CHECK_EXTERNS+"var goog;goog.provide('a.b');a.b.C = class {};")
}

// Java: VarCheckTest.testGoogForwardDeclare_canBeReferencedInExternsIfInCode()
func (c *VarCheckTest) TestGoogForwardDeclare_canBeReferencedInExternsIfInCode() {
	c.TestNoWarning(
		// Allow referencing `goog` in externs even though it is actually only defined in code.
		c.Externs("goog.forwardDeclare('a.b.C');"), c.Srcs("var /** !a.b.C */ c; var goog;"))
	c.CheckSynthesizedExternWithExterns(
		// Don't synthesize an extern for 'goog'.
		c.Externs("goog.forwardDeclare('a.b.C');"),
		c.Srcs("var goog;"),
		CompilerTestCase_VAR_CHECK_EXTERNS+"goog.forwardDeclare('a.b.C');")
}

// Java: VarCheckTest.testArbitraryPropertyNamedForwardDeclare_cannotBeReferencedInExterns()
func (c *VarCheckTest) TestArbitraryPropertyNamedForwardDeclare_cannotBeReferencedInExterns() {
	c.TestWarning(
		c.Externs("notGoog.forwardDeclare('a.b.C');"),
		c.Srcs("var notGoog;"),
		c.Warning(VarCheck_UNDEFINED_EXTERN_VAR_ERROR))
}

// Java: VarCheckTest.testGoogProvide_cannotBeReferencedInExterns()
func (c *VarCheckTest) TestGoogProvide_cannotBeReferencedInExterns() {
	c.TestWarning(
		c.Externs("goog.provide('a.b');"), c.Srcs("var goog;"), c.Warning(VarCheck_UNDEFINED_EXTERN_VAR_ERROR))
}

// Java: VarCheckTest.testGoogModule_cannotBeReferencedInExterns()
func (c *VarCheckTest) TestGoogModule_cannotBeReferencedInExterns() {
	c.TestWarning(
		c.Externs("goog.module('a.b');"), c.Srcs("var goog;"), c.Warning(VarCheck_UNDEFINED_EXTERN_VAR_ERROR))
}

// Java: VarCheckTest.testReferenceToWeakVar_fromStrongFile()
func (c *VarCheckTest) TestReferenceToWeakVar_fromStrongFile() {
	weakChunk := NewJSChunk(JSChunk_WEAK_CHUNK_NAME)
	strongChunk := NewJSChunk(JSChunk_STRONG_CHUNK_NAME)
	weakChunk.AddDependency(strongChunk)

	weakChunk.AddSourceFile(SourceFile_FromCodeWithKind("weak.js", "var weakVar = 0;", rhino.StaticSourceFile_SourceKind_WEAK))

	strongChunk.AddSourceFile(SourceFile_FromCodeWithKind("strong.js", "weakVar();", rhino.StaticSourceFile_SourceKind_STRONG))

	c.Test(
		c.Srcs(
			[]*JSChunk{strongChunk, weakChunk}),
		c.Error(VarCheck_VIOLATED_CHUNK_DEP_ERROR),
		c.Error(VarCheck_UNDEFINED_VAR_ERROR))
}

// Java: VarCheckTest.testReferenceToWeakVar_fromWeakFile()
func (c *VarCheckTest) TestReferenceToWeakVar_fromWeakFile() {
	weakChunk := NewJSChunk(JSChunk_WEAK_CHUNK_NAME)
	weakChunk.AddSourceFile(
		SourceFile_FromCodeWithKind(
			"weak.js",
			`var weakVar = 0;
weakVar();
`,
			rhino.StaticSourceFile_SourceKind_WEAK))

	c.TestSame(c.Srcs(weakChunk))
}

// Java: VarCheckTest.testReferenceToWeakNamespaceRoot_fromStrongFile()
func (c *VarCheckTest) TestReferenceToWeakNamespaceRoot_fromStrongFile() {
	weakChunk := NewJSChunk(JSChunk_WEAK_CHUNK_NAME)
	strongChunk := NewJSChunk(JSChunk_STRONG_CHUNK_NAME)
	weakChunk.AddDependency(strongChunk)

	weakChunk.AddSourceFile(
		SourceFile_FromCodeWithKind(
			"weak.js",
			testExternsBuilder_GetClosureExternsAsSource()+"goog.provide('foo.bar');",
			rhino.StaticSourceFile_SourceKind_WEAK))

	strongChunk.AddSourceFile(SourceFile_FromCodeWithKind("strong.js", "foo();", rhino.StaticSourceFile_SourceKind_STRONG))

	c.Test(
		c.Srcs(
			[]*JSChunk{strongChunk, weakChunk}),
		c.Error(VarCheck_UNDEFINED_VAR_ERROR))
}

// Java: VarCheckTest.testReferenceToWeakNamespaceRoot_fromStrongFile_synthesizesExtern()
func (c *VarCheckTest) TestReferenceToWeakNamespaceRoot_fromStrongFile_synthesizesExtern() {
	weakChunk := NewJSChunk(JSChunk_WEAK_CHUNK_NAME)
	strongChunk := NewJSChunk(JSChunk_STRONG_CHUNK_NAME)
	weakChunk.AddDependency(strongChunk)

	weakChunk.AddSourceFile(
		SourceFile_FromCodeWithKind(
			"weak.js",
			testExternsBuilder_GetClosureExternsAsSource()+"goog.provide('foo.bar');",
			rhino.StaticSourceFile_SourceKind_WEAK))

	strongChunk.AddSourceFile(
		SourceFile_FromCodeWithKind(
			"strong.js",
			`/** @suppress {undefinedVars} */
foo();
`,
			rhino.StaticSourceFile_SourceKind_STRONG))

	c.TestExternChanges(
		c.Externs(""), c.Srcs(strongChunk, weakChunk), c.Expected("var foo;"+CompilerTestCase_VAR_CHECK_EXTERNS))
}

// Java: VarCheckTest.testReferenceToweakChunkNamespaceRoot_fromStrongFile()
func (c *VarCheckTest) TestReferenceToweakChunkNamespaceRoot_fromStrongFile() {
	weakChunk := NewJSChunk(JSChunk_WEAK_CHUNK_NAME)
	strongChunk := NewJSChunk(JSChunk_STRONG_CHUNK_NAME)
	weakChunk.AddDependency(strongChunk)

	weakChunk.AddSourceFile(
		SourceFile_FromCodeWithKind(
			"weak.js",
			testExternsBuilder_GetClosureExternsAsSource()+
				"goog.module('foo.bar'); goog.module.declareLegacyNamespace();",
			rhino.StaticSourceFile_SourceKind_WEAK))

	strongChunk.AddSourceFile(SourceFile_FromCodeWithKind("strong.js", "foo();", rhino.StaticSourceFile_SourceKind_STRONG))

	c.Test(c.Srcs(strongChunk, weakChunk), c.Error(VarCheck_UNDEFINED_VAR_ERROR))
}

// Java: VarCheckTest.testReferenceToStrongNamespaceRoot_withAdditionalWeakProvide_fromStrongFile()
func (c *VarCheckTest) TestReferenceToStrongNamespaceRoot_withAdditionalWeakProvide_fromStrongFile() {
	weakChunk := NewJSChunk(JSChunk_WEAK_CHUNK_NAME)
	strongChunk := NewJSChunk(JSChunk_STRONG_CHUNK_NAME)
	weakChunk.AddDependency(strongChunk)

	weakChunk.AddSourceFile(SourceFile_FromCodeWithKind("weak.js", "goog.provide('foo.bar');", rhino.StaticSourceFile_SourceKind_WEAK))

	strongChunk.AddSourceFile(
		SourceFile_FromCodeWithKind(
			"strong0.js",
			testExternsBuilder_GetClosureExternsAsSource()+"goog.provide('foo.qux');",
			rhino.StaticSourceFile_SourceKind_STRONG))
	strongChunk.AddSourceFile(SourceFile_FromCodeWithKind("strong1.js", "foo();", rhino.StaticSourceFile_SourceKind_STRONG))

	c.TestSame(
		c.Srcs(
			[]*JSChunk{strongChunk, weakChunk}))
}

// Java: VarCheckTest.VariableTestCheck
type VarCheckTest_VariableTestCheck struct {
	compiler AbstractCompiler
}

// Java: VarCheckTest.VariableTestCheck(AbstractCompiler)
func NewVarCheckTest_VariableTestCheck(compiler AbstractCompiler) *VarCheckTest_VariableTestCheck {
	v := &VarCheckTest_VariableTestCheck{}
	v.compiler = compiler
	return v
}

// Java: VarCheckTest.VariableTestCheck.process(Node, Node)
func (v *VarCheckTest_VariableTestCheck) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_TraverseRoots(
		v.compiler,
		&varCheckTest_variableTestCheckCallback{},
		externs,
		root)
}

// varCheckTest_variableTestCheckCallback is the anonymous
// AbstractPostOrderCallback of VariableTestCheck.process.
type varCheckTest_variableTestCheckCallback struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	// Go-only: language-adapter — distinct identity.
	_ byte
}

// Java: VarCheckTest.VariableTestCheck.process(Node, Node).new AbstractPostOrderCallback() {...}.visit(NodeTraversal, Node, Node)
func (cb *varCheckTest_variableTestCheckCallback) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsName() && !parent.IsFunction() && !parent.IsLabel() {
		ScopeSubject_AssertScope(t.GetScope()).Declares(n.GetString())
	}
}

// Java: VarCheckTest.checkSynthesizedExtern(Sources, String)
func (c *VarCheckTest) CheckSynthesizedExtern(input CompilerTestCase_Sources, expectedExtern string) {
	c.CheckSynthesizedExternWithExterns(c.Externs(""), input, expectedExtern)
}

// Java: VarCheckTest.checkSynthesizedExtern(Externs, Sources, String, Diagnostic...)
func (c *VarCheckTest) CheckSynthesizedExternWithExterns(extern *CompilerTestCase_Externs, input CompilerTestCase_Sources, expectedExtern string, warnings ...*CompilerTestCase_Diagnostic) {
	c.declarationCheck = !c.validityCheck
	c.DisableCompareAsTree()
	args := []any{extern, input, c.Expected(expectedExtern)}
	for _, w := range warnings {
		args = append(args, w)
	}
	c.TestExternChanges(args...)
}
