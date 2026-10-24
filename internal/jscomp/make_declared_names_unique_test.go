// Port of com/google/javascript/jscomp/MakeDeclaredNamesUniqueTest.java @ closure-compiler v20261021

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
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// The Java @Test methods not named test* get a "Test" prefix so that
// RunJUnit4 runs them.
//
// Java: public final class MakeDeclaredNamesUniqueTest extends CompilerTestCase
type MakeDeclaredNamesUniqueTest struct {
	CompilerTestCaseBase

	// this.useDefaultRenamer = true; invokes the ContextualRenamer
	// this.useDefaultRenamer = false; invokes the InlineRenamer
	useDefaultRenamer bool
	// invert = true; treats JavaScript input as normalized code and inverts the renaming
	// invert = false; conducts renaming
	invert bool
	// removeConst = true; removes const-ness of a name (e.g. If the variable name is CONST)
	removeConst bool
	// whether to throw an exception on any names newly made unique.
	assertOnChange bool
}

// Java: MakeDeclaredNamesUniqueTest.LOCAL_NAME_PREFIX
const makeDeclaredNamesUniqueTest_LOCAL_NAME_PREFIX = "unique_"

// Java: MakeDeclaredNamesUniqueTest()
func NewMakeDeclaredNamesUniqueTest() *MakeDeclaredNamesUniqueTest {
	c := &MakeDeclaredNamesUniqueTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestMakeDeclaredNamesUniqueTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewMakeDeclaredNamesUniqueTest() })
}

// makeDeclaredNamesUniqueTest_pass is the anonymous CompilerPass of
// getProcessor.
type makeDeclaredNamesUniqueTest_pass struct {
	c        *MakeDeclaredNamesUniqueTest
	compiler *Compiler
}

// Java: MakeDeclaredNamesUniqueTest.getProcessor(...).new CompilerPass() {...}.process(Node, Node)
func (p *makeDeclaredNamesUniqueTest_pass) Process(externs rhino.Node, root rhino.Node) {
	renamer := MakeDeclaredNamesUnique_Builder_().WithAssertOnChange(p.c.assertOnChange)
	if !p.c.useDefaultRenamer {
		renamer =
			renamer.WithRenamer(
				NewMakeDeclaredNamesUnique_InlineRenamer(
					p.compiler.GetCodingConvention(),
					p.compiler.GetUniqueNameIdSupplier(),
					makeDeclaredNamesUniqueTest_LOCAL_NAME_PREFIX,
					p.c.removeConst,
					true,
					nil))
	}
	NodeTraversal_TraverseRoots(p.compiler, renamer.Build(), externs, root)
}

// Java: MakeDeclaredNamesUniqueTest.getProcessor(Compiler)
func (c *MakeDeclaredNamesUniqueTest) GetProcessor(compiler *Compiler) CompilerPass {
	if !c.invert {
		return &makeDeclaredNamesUniqueTest_pass{c: c, compiler: compiler}
	} else {
		return MakeDeclaredNamesUnique_GetContextualRenameInverter(compiler)
	}
}

// Java: MakeDeclaredNamesUniqueTest.getOptions()
func (c *MakeDeclaredNamesUniqueTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetWarningLevel(DiagnosticGroups_MODULE_LOAD, CheckLevel_OFF)
	return options
}

// JUnitBefore runs CompilerTestCase's @Before setUp, then the class's own
// @Before customSetUp (JUnit runs superclass @Before methods first).
func (c *MakeDeclaredNamesUniqueTest) JUnitBefore() {
	c.CompilerTestCaseBase.JUnitBefore()
	c.CustomSetUp()
}

// Java: MakeDeclaredNamesUniqueTest.customSetUp()
func (c *MakeDeclaredNamesUniqueTest) CustomSetUp() {
	c.removeConst = false
	c.invert = false
	c.useDefaultRenamer = false
	c.assertOnChange = false
}

// Java: MakeDeclaredNamesUniqueTest.testWithInversion(String, String)
func (c *MakeDeclaredNamesUniqueTest) testWithInversion(original string, expected string) {
	c.invert = false
	c.Test(original, expected)
	c.invert = true
	c.Test(expected, original)
	c.invert = false
}

// Java: MakeDeclaredNamesUniqueTest.testWithInversion(String[], String[])
func (c *MakeDeclaredNamesUniqueTest) testWithInversionOfArrays(original []string, expected []string) {
	c.invert = false
	c.Test(c.Srcs(original), c.Expected(expected))
	c.invert = true
	c.Test(c.Srcs(expected), c.Expected(original))
	c.invert = false
}

// Java: MakeDeclaredNamesUniqueTest.testSameWithInversion(String, String)
func (c *MakeDeclaredNamesUniqueTest) testSameWithInversionWithExterns(externs string, original string) {
	c.invert = false
	c.TestSame(c.Externs(externs), c.Srcs(original))
	c.invert = true
	c.TestSame(c.Externs(externs), c.Srcs(original))
	c.invert = false
}

// Java: MakeDeclaredNamesUniqueTest.testSameWithInversion(String)
func (c *MakeDeclaredNamesUniqueTest) testSameWithInversion(original string) {
	c.testSameWithInversionWithExterns("", original)
}

// Java: MakeDeclaredNamesUniqueTest.wrapInFunction(String)
func makeDeclaredNamesUniqueTest_wrapInFunction(s string) string {
	return "function f(){" + s + "}"
}

// Java: MakeDeclaredNamesUniqueTest.testInFunctionWithInversion(String, String)
func (c *MakeDeclaredNamesUniqueTest) testInFunctionWithInversion(original string, expected string) {
	c.testWithInversion(
		makeDeclaredNamesUniqueTest_wrapInFunction(original), makeDeclaredNamesUniqueTest_wrapInFunction(expected))
}

// Java: MakeDeclaredNamesUniqueTest.testMakeDeclaredNamesUniqueNullishCoalesce()
func (c *MakeDeclaredNamesUniqueTest) TestMakeDeclaredNamesUniqueNullishCoalesce() {
	c.useDefaultRenamer = true

	c.Test(
		"var foo; var x = function foo(){var foo = false ?? {};}",
		"var foo; var x = function foo$jscomp$1(){var foo$jscomp$2 = false ?? {}}")
	c.testSameWithInversion("var a = b ?? c;")
}

// Java: MakeDeclaredNamesUniqueTest.testShadowedBleedingName()
func (c *MakeDeclaredNamesUniqueTest) TestShadowedBleedingName() {
	c.useDefaultRenamer = true

	c.Test(
		"var foo; var x = function foo(){var foo;}",
		"var foo; var x = function foo$jscomp$1(){var foo$jscomp$2}")
}

// Java: MakeDeclaredNamesUniqueTest.testMakeLocalNamesUniqueWithContext1()
func (c *MakeDeclaredNamesUniqueTest) TestMakeLocalNamesUniqueWithContext1() {
	c.useDefaultRenamer = true

	c.invert = true
	c.Test(
		"var a;function foo(){var a$jscomp$inline_1; a = 1}",
		"var a;function foo(){var a$jscomp$0       ; a = 1}")
	c.Test(
		"var a;function foo(){var a$jscomp$inline_1;}", //
		"var a;function foo(){var a                ;}")

	c.Test(
		"let a;function foo(){let a$jscomp$inline_1; a = 1}",
		"let a;function foo(){let a$jscomp$0       ; a = 1}")
	c.Test(
		"const a = 1;function foo(){let a$jscomp$inline_1;}", //
		"const a = 1;function foo(){let a                ;}")
	c.Test(
		"class A {} function foo(){class A$jscomp$inline_1 {}}",
		"class A {} function foo(){class A {}}")
}

// Java: MakeDeclaredNamesUniqueTest.testMakeLocalNamesUniqueWithContext2()
func (c *MakeDeclaredNamesUniqueTest) TestMakeLocalNamesUniqueWithContext2() {
	// Set the test type
	c.useDefaultRenamer = true

	// Verify global names are untouched.
	c.testSameWithInversion("var a;")
	c.testSameWithInversion("let a;")
	c.testSameWithInversion("const a = 0;")

	// Verify global names are untouched.
	c.testSameWithInversion("a;")

	// Local names are made unique.
	c.testWithInversion(
		"var a;function foo(a){var b;a}", "var a;function foo(a$jscomp$1){var b;a$jscomp$1}")
	c.testWithInversion(
		"var a;function foo(){var b;a}function boo(){var b;a}",
		"var a;function foo(){var b;a}function boo(){var b$jscomp$1;a}")
	c.testWithInversion(
		`function foo(a){var b}
function boo(a){var b}
`,
		`function foo(a){var b}
function boo(a$jscomp$1){var b$jscomp$1}
`)
	// variable b is left untouched because it is only declared once
	c.testWithInversion(
		"let a;function foo(a){let b;a}", "let a;function foo(a$jscomp$1){let b;a$jscomp$1}")
	c.testWithInversion(
		"let a;function foo(){let b;a}function boo(){let b;a}",
		"let a;function foo(){let b;a}function boo(){let b$jscomp$1;a}")
	c.testWithInversion(
		`function foo(a){let b}
function boo(a){let b}
`,
		`function foo(a){let b}
function boo(a$jscomp$1){let b$jscomp$1}
`)

	// Verify functions expressions are renamed.
	c.testWithInversion(
		"var a = function foo(){foo()};var b = function foo(){foo()};",
		"var a = function foo(){foo()};var b = function foo$jscomp$1(){foo$jscomp$1()};")
	c.testWithInversion(
		"let a = function foo(){foo()};let b = function foo(){foo()};",
		"let a = function foo(){foo()};let b = function foo$jscomp$1(){foo$jscomp$1()};")

	// Verify catch exceptions names are made unique
	c.testSameWithInversion("try { } catch(e) {e;}")

	// Inversion does not handle exceptions correctly.
	c.Test(
		"try { } catch(e) {e;}; try { } catch(e) {e;}",
		"try { } catch(e) {e;}; try { } catch(e$jscomp$1) {e$jscomp$1;}")
	c.Test(
		"try { } catch(e) {e; try { } catch(e) {e;}};",
		"try { } catch(e) {e; try { } catch(e$jscomp$1) {e$jscomp$1;} }; ")
}

// Java: MakeDeclaredNamesUniqueTest.testMakeLocalNamesUniqueWithContext3()
func (c *MakeDeclaredNamesUniqueTest) TestMakeLocalNamesUniqueWithContext3() {
	// Set the test type
	c.useDefaultRenamer = true

	externs := "var extern1 = {};"

	// Verify global names are untouched.
	c.testSameWithInversionWithExterns(externs, "var extern1 = extern1 || {};")

	// Verify global names are untouched.
	c.TestSame(c.Externs(externs), c.Srcs("var extern1 = extern1 || {};"))
}

// Java: MakeDeclaredNamesUniqueTest.testMakeLocalNamesUniqueWithContext4()
func (c *MakeDeclaredNamesUniqueTest) TestMakeLocalNamesUniqueWithContext4() {
	// Set the test type
	c.useDefaultRenamer = true

	c.testInFunctionWithInversion(
		"var e; try { } catch(e) {e;}; try { } catch(e) {e;}",
		"var e; try { } catch(e$jscomp$1) {e$jscomp$1;}; try { } catch(e$jscomp$2) {e$jscomp$2;}")
	c.testInFunctionWithInversion(
		"var e; try { } catch(e         ) {         e; try { } catch(e         ) {e;         } };",
		"var e; try { } catch(e$jscomp$1) {e$jscomp$1; try { } catch(e$jscomp$2) {e$jscomp$2;} };")
	c.testInFunctionWithInversion(
		"try { } catch(e         ) {e         ;}; try { } catch(e         ) {e         ;}; var e;",
		"try { } catch(e$jscomp$1) {e$jscomp$1;}; try { } catch(e$jscomp$2) {e$jscomp$2;}; var e;")
	c.testInFunctionWithInversion(
		"try { } catch(e         ) {e         ; try { } catch(e         ) {e         ;} }; var e;",
		"try { } catch(e$jscomp$1) {e$jscomp$1; try { } catch(e$jscomp$2) {e$jscomp$2;} }; var e;")
}

// Java: MakeDeclaredNamesUniqueTest.testMakeLocalNamesUniqueWithContext5()
func (c *MakeDeclaredNamesUniqueTest) TestMakeLocalNamesUniqueWithContext5() {
	c.useDefaultRenamer = true
	c.testWithInversion("function f(){var f; f = 1}", "function f(){var f$jscomp$1; f$jscomp$1 = 1}")
}

// Java: MakeDeclaredNamesUniqueTest.testMakeLocalNamesUniqueWithContext6()
func (c *MakeDeclaredNamesUniqueTest) TestMakeLocalNamesUniqueWithContext6() {
	c.useDefaultRenamer = true
	c.testWithInversion("function f(f){f = 1}", "function f(f$jscomp$1){f$jscomp$1 = 1}")
}

// Java: MakeDeclaredNamesUniqueTest.testMakeLocalNamesUniqueWithContext7()
func (c *MakeDeclaredNamesUniqueTest) TestMakeLocalNamesUniqueWithContext7() {
	c.useDefaultRenamer = true
	c.testWithInversion(
		"function f(f){var f; f = 1}", "function f(f$jscomp$1){var f$jscomp$1; f$jscomp$1 = 1}")
}

// Java: MakeDeclaredNamesUniqueTest.testMakeLocalNamesUniqueWithContext8()
func (c *MakeDeclaredNamesUniqueTest) TestMakeLocalNamesUniqueWithContext8() {
	c.useDefaultRenamer = true
	c.Test(
		"var fn = function f(){var f; f = 1}",
		"var fn = function f(){var f$jscomp$1; f$jscomp$1 = 1}")
}

// Java: MakeDeclaredNamesUniqueTest.testMakeLocalNamesUniqueWithContext9()
func (c *MakeDeclaredNamesUniqueTest) TestMakeLocalNamesUniqueWithContext9() {
	c.useDefaultRenamer = true
	c.TestSame("var fn = function f(f){f = 1}")
}

// Java: MakeDeclaredNamesUniqueTest.testMakeLocalNamesUniqueWithContext10()
func (c *MakeDeclaredNamesUniqueTest) TestMakeLocalNamesUniqueWithContext10() {
	c.useDefaultRenamer = true
	c.TestSame("var fn = function f(f){var f; f = 1}")
}

// Java: MakeDeclaredNamesUniqueTest.testMakeLocalNamesUniqueWithContext11()
func (c *MakeDeclaredNamesUniqueTest) TestMakeLocalNamesUniqueWithContext11() {
	c.useDefaultRenamer = true
	// Changes the parameter name if it's the same as loop object name
	c.Test(
		`var loopObjName = {};
for(;;) {
   var fn = (function f(loopObjName){ return function() {}; })(loopObjName);
}
`,
		`var loopObjName = {};
for(;;) {
   var fn = (function f(loopObjName$jscomp$1){ return function() {}; })(loopObjName);
}
`)

	// parameter name is already unique; no change
	c.TestSame(
		`var foo = {};
for(;;) {
   var fn = (function f(bar){ return function() {}; })(foo);}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testMakeFunctionsUniqueWithContext()
func (c *MakeDeclaredNamesUniqueTest) TestMakeFunctionsUniqueWithContext() {
	c.useDefaultRenamer = true
	c.TestSame("function f(){} function f(){}")
	c.TestSame("var x = function() {function f(){} function f(){}};")
}

// Java: MakeDeclaredNamesUniqueTest.testMakeFunctionsUniqueWithContext1()
func (c *MakeDeclaredNamesUniqueTest) TestMakeFunctionsUniqueWithContext1() {
	c.useDefaultRenamer = true
	c.Test(
		"if (1) { function f(){} } else { function f(){} }",
		"if (1) { function f(){} } else { function f$jscomp$1(){} }")
}

// Java: MakeDeclaredNamesUniqueTest.testMakeFunctionsUniqueWithContext2()
func (c *MakeDeclaredNamesUniqueTest) TestMakeFunctionsUniqueWithContext2() {
	c.useDefaultRenamer = true
	c.TestSame("if (1) { function f(){} function f(){} }")
}

// Java: MakeDeclaredNamesUniqueTest.testMakeFunctionsUniqueWithContext3()
func (c *MakeDeclaredNamesUniqueTest) TestMakeFunctionsUniqueWithContext3() {
	c.useDefaultRenamer = true
	c.Test(
		"function f() {} if (1) { function f(){} function f(){} }",
		"function f() {} if (1) { function f$jscomp$1(){} function f$jscomp$1(){} }")
}

// Java: MakeDeclaredNamesUniqueTest.testArguments()
func (c *MakeDeclaredNamesUniqueTest) TestArguments() {
	// Set the test type
	c.useDefaultRenamer = true

	// Don't distinguish between "arguments", it can't be made unique.
	c.testSameWithInversion("function foo(){var arguments;function bar(){var arguments;}}")

	c.invert = true

	// Don't introduce new references to arguments, it is special.
	// Still try to rename it to a name that depends on the shape of the AST rather than
	// the process we happened to follow to reach that shape.
	c.Test(
		"function foo(){var arguments$jscomp$1;}", //
		"function foo(){var arguments$jscomp$0;}")
}

// Java: MakeDeclaredNamesUniqueTest.testClassInForLoop()
func (c *MakeDeclaredNamesUniqueTest) TestClassInForLoop() {
	c.useDefaultRenamer = true
	c.TestSame("for (class a {};;) { break; }")
}

// Java: MakeDeclaredNamesUniqueTest.testFunctionInForLoop()
func (c *MakeDeclaredNamesUniqueTest) TestFunctionInForLoop() {
	c.useDefaultRenamer = true
	c.TestSame("for (function a() {};;) { break; }")
}

// Java: MakeDeclaredNamesUniqueTest.testLetsInSeparateBlocks()
func (c *MakeDeclaredNamesUniqueTest) TestLetsInSeparateBlocks() {
	c.useDefaultRenamer = true
	c.Test(
		`if (x) {
  let e;
  alert(e);
}
if (y) {
  let e;
  alert(e);
}
`,
		`if (x) {
  let e;
  alert(e);
}
if (y) {
  let e$jscomp$1;
  alert(e$jscomp$1);
}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testConstInGlobalHoistScope()
func (c *MakeDeclaredNamesUniqueTest) TestConstInGlobalHoistScope() {
	c.useDefaultRenamer = true
	c.TestSame(
		`if (true) {
  const x = 1; alert(x);
}
`)

	c.Test(
		`if (true) {
  const x = 1; alert(x);
} else {
  const x = 1; alert(x);
}
`,
		`if (true) {
  const x = 1; alert(x);
} else {
  const x$jscomp$1 = 1; alert(x$jscomp$1);
}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testMakeLocalNamesUniqueWithoutContext()
func (c *MakeDeclaredNamesUniqueTest) TestMakeLocalNamesUniqueWithoutContext() {
	c.useDefaultRenamer = false

	c.Test("var a;", "var a$jscomp$unique_0")
	c.Test("let a;", "let a$jscomp$unique_0")

	// Verify undeclared names are untouched.
	c.TestSame("a;")

	// Local names are made unique.
	c.Test(
		`var a;
function foo(a){var b;a}
`,
		`var a$jscomp$unique_0;
function foo$jscomp$unique_1(a$jscomp$unique_2){
  var b$jscomp$unique_3;a$jscomp$unique_2}
`)
	c.Test(
		`var a;
function foo(){var b;a}
function boo(){var b;a}
`,
		`var a$jscomp$unique_0;
function foo$jscomp$unique_1(){var b$jscomp$unique_3;a$jscomp$unique_0}
function boo$jscomp$unique_2(){var b$jscomp$unique_4;a$jscomp$unique_0}
`)

	c.Test(
		"let a; function foo(a) {let b; a; }",
		`let a$jscomp$unique_0;
function foo$jscomp$unique_1(a$jscomp$unique_2) {
  let b$jscomp$unique_3;
  a$jscomp$unique_2;
}
`)

	c.Test(
		`let a;
function foo() { let b; a; }
function boo() { let b; a; }
`,
		`let a$jscomp$unique_0;
function foo$jscomp$unique_1() {
  let b$jscomp$unique_3;
  a$jscomp$unique_0;
}
function boo$jscomp$unique_2() {
  let b$jscomp$unique_4;
  a$jscomp$unique_0;
}
`)

	// Verify function expressions are renamed.
	c.Test(
		"var a = function foo(){foo()};",
		"var a$jscomp$unique_0 = function foo$jscomp$unique_1(){foo$jscomp$unique_1()};")
	c.Test(
		"const a = function foo(){foo()};",
		"const a$jscomp$unique_0 = function foo$jscomp$unique_1(){foo$jscomp$unique_1()};")

	// Verify catch exceptions names are made unique
	c.Test("try { } catch(e) {e;}", "try { } catch(e$jscomp$unique_0) {e$jscomp$unique_0;}")
	c.Test(
		`try { } catch(e) {e;};
try { } catch(e) {e;}
`,
		`try { } catch(e$jscomp$unique_0) {e$jscomp$unique_0;};
try { } catch(e$jscomp$unique_1) {e$jscomp$unique_1;}
`)
	c.Test(
		`try { } catch(e) {e;
try { } catch(e) {e;}};
`,
		`try { } catch(e$jscomp$unique_0) {e$jscomp$unique_0;
try { } catch(e$jscomp$unique_1) {e$jscomp$unique_1;} }; 
`)
}

// Java: MakeDeclaredNamesUniqueTest.testMakeLocalNamesUniqueWithoutContext2()
func (c *MakeDeclaredNamesUniqueTest) TestMakeLocalNamesUniqueWithoutContext2() {
	// Set the test type
	c.useDefaultRenamer = false

	c.Test("var _a;", "var JSCompiler__a$jscomp$unique_0")
	c.Test(
		"var _a = function _b(_c) { var _d; };",
		`var JSCompiler__a$jscomp$unique_0 = function JSCompiler__b$jscomp$unique_2(
JSCompiler__c$jscomp$unique_1) { var JSCompiler__d$jscomp$unique_3; };
`)

	c.Test("let _a;", "let JSCompiler__a$jscomp$unique_0")
	c.Test(
		"const _a = function _b(_c) { let _d; };",
		`const JSCompiler__a$jscomp$unique_0 = function JSCompiler__b$jscomp$unique_2(
JSCompiler__c$jscomp$unique_1) { let JSCompiler__d$jscomp$unique_3; };
`)
}

// Java: MakeDeclaredNamesUniqueTest.testOnlyInversion()
func (c *MakeDeclaredNamesUniqueTest) TestOnlyInversion() {
	c.invert = true
	c.Test(
		"function f(a, a$jscomp$1) {}", //
		"function f(a, a$jscomp$0) {}")
	c.Test(
		"function f(a$jscomp$1, b$jscomp$2) {}", //
		"function f(a         , b         ) {}")
	c.Test(
		"function f(a$jscomp$1, a$jscomp$2) {}", //
		"function f(a         , a$jscomp$0) {}")
	c.Test(
		"try { } catch(e) {e; try { } catch(e$jscomp$1) {e$jscomp$1;} }; ",
		"try { } catch(e) {e; try { } catch(e) {e;} }; ")
	c.TestSame("var a$jscomp$1;")
	c.TestSame("const a$jscomp$1 = 1;")
	c.TestSame("function f() { var $jscomp$; }")
	c.TestSame("var CONST = 3; var b = CONST;")
	c.Test(
		"function f() {var CONST = 3; var ACONST$jscomp$1 = 2;}",
		"function f() {var CONST = 3; var ACONST = 2;}")
	c.Test(
		"function f() {const CONST = 3; const ACONST$jscomp$1 = 2;}",
		"function f() {const CONST = 3; const ACONST = 2;}")
}

// Java: MakeDeclaredNamesUniqueTest.testOnlyInversion2()
func (c *MakeDeclaredNamesUniqueTest) TestOnlyInversion2() {
	c.invert = true
	c.Test(
		"function f() {try { } catch(e) {e;}; try { } catch(e$jscomp$0) {e$jscomp$0;}}",
		"function f() {try { } catch(e) {e;}; try { } catch(e) {e;}}")
}

// Java: MakeDeclaredNamesUniqueTest.testOnlyInversion3()
func (c *MakeDeclaredNamesUniqueTest) TestOnlyInversion3() {
	c.invert = true
	c.Test(
		`function x1() {
  var a$jscomp$1;
  function x2() {
    var a$jscomp$2;
  }
  function x3() {
    var a$jscomp$3;
  }
}
`,
		`function x1() {
  var a;
  function x2() {
    var a;
  }
  function x3() {
    var a;
  }
}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testOnlyInversion4()
func (c *MakeDeclaredNamesUniqueTest) TestOnlyInversion4() {
	c.invert = true
	c.TestSame(
		`function x1() {
// The attempt to rename will re-generate this same exact name.
// The purpose of this test is to make sure we don't accidentally report
// this "renaming" to the same name as a change.
  var a$jscomp$0;
  function x2() {
    var a;a$jscomp$0++
  }
}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testOnlyInversion5()
func (c *MakeDeclaredNamesUniqueTest) TestOnlyInversion5() {
	c.invert = true
	c.Test(
		`function x1() {
  const a$jscomp$1 = 0;
  function x2() {
    const b$jscomp$1 = 0;
  }
}
`,
		`function x1() {
  const a = 0;
  function x2() {
    const b = 0;
  }
}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testConstRemovingRename1()
func (c *MakeDeclaredNamesUniqueTest) TestConstRemovingRename1() {
	c.removeConst = true
	c.Test(
		"(function () {var CONST = 3; var ACONST$jscomp$1 = 2;})",
		"(function () {var CONST$jscomp$unique_0 = 3; var ACONST$jscomp$unique_1 = 2;})")
}

// Java: MakeDeclaredNamesUniqueTest.testConstRemovingRename2()
func (c *MakeDeclaredNamesUniqueTest) TestConstRemovingRename2() {
	c.removeConst = true
	c.Test(
		"var CONST = 3; var b = CONST;",
		"var CONST$jscomp$unique_0 = 3; var b$jscomp$unique_1 = CONST$jscomp$unique_0;")
}

// Java: MakeDeclaredNamesUniqueTest.testConstRemovingRenameAlsoRemovesAnnotation()
func (c *MakeDeclaredNamesUniqueTest) TestConstRemovingRenameAlsoRemovesAnnotation() {
	c.removeConst = true
	c.Test(
		"/** @const */ var c = 3; var b = c;",
		"/** blank */ var c$jscomp$unique_0 = 3; var b$jscomp$unique_1 = c$jscomp$unique_0;")
}

// Java: MakeDeclaredNamesUniqueTest.testRestParamWithoutContext()
func (c *MakeDeclaredNamesUniqueTest) TestRestParamWithoutContext() {
	c.Test(
		"function f(...x) { x; }",
		"function f$jscomp$unique_0(...x$jscomp$unique_1) { x$jscomp$unique_1; }")
}

// Java: MakeDeclaredNamesUniqueTest.testRestParamWithContextWithInversion()
func (c *MakeDeclaredNamesUniqueTest) TestRestParamWithContextWithInversion() {
	c.useDefaultRenamer = true
	c.testWithInversion(
		`let x = 0;
function foo(...x) {
  return x[0];
}
`,
		`let x = 0;
function foo(...x$jscomp$1) {
  return x$jscomp$1[0]
}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testVarParamSameName0()
func (c *MakeDeclaredNamesUniqueTest) TestVarParamSameName0() {
	c.Test(
		`function f(x) {
  if (!x) var x = 6;
}
`,
		`function f$jscomp$unique_0(x$jscomp$unique_1) {
  if (!x$jscomp$unique_1) var x$jscomp$unique_1 = 6;
}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testVarParamSameName1()
func (c *MakeDeclaredNamesUniqueTest) TestVarParamSameName1() {
	c.Test(
		`function f(x) {
  if (!x) x = 6;
}
`,
		`function f$jscomp$unique_0(x$jscomp$unique_1) {
  if (!x$jscomp$unique_1) x$jscomp$unique_1 = 6;
}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testVarParamSameAsLet0()
func (c *MakeDeclaredNamesUniqueTest) TestVarParamSameAsLet0() {
	c.Test(
		`function f(x) {
  if (!x) { let x = 6; }
}
`,
		`function f$jscomp$unique_0(x$jscomp$unique_1) {
  if (!x$jscomp$unique_1) { let x$jscomp$unique_2 = 6; }
}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testObjectProperties()
func (c *MakeDeclaredNamesUniqueTest) TestObjectProperties() {
	c.Test("var a = {x : 'a'};", "var a$jscomp$unique_0 = {x : 'a'};")
	c.Test("let a = {x : 'a'};", "let a$jscomp$unique_0 = {x : 'a'};")
	c.Test("const a = {x : 'a'};", "const a$jscomp$unique_0 = {x : 'a'};")
	c.Test("var a = {x : 'a'}; a.x", "var a$jscomp$unique_0 = {x : 'a'}; a$jscomp$unique_0.x")
}

// Java: MakeDeclaredNamesUniqueTest.testClassesWithContextWithInversion()
func (c *MakeDeclaredNamesUniqueTest) TestClassesWithContextWithInversion() {
	c.useDefaultRenamer = true
	c.testWithInversion(
		`var a;
class Foo {
  constructor(a) {
    this.a = a;
  }
  f() {
    var x = 1;
    return a + x;
  }
}
`,
		`var a;
class Foo {
  constructor(a$jscomp$1) {
    this.a = a$jscomp$1;
  }
  f() {
    var x = 1;
    return a + x;
  }
}
`)

	// class declarations are block-scoped but not hoisted.
	c.testSameWithInversion(
		`{
  let x = new Foo(); // ReferenceError
  class Foo {}
}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testBlockScopesWithContextWithInversion1()
func (c *MakeDeclaredNamesUniqueTest) TestBlockScopesWithContextWithInversion1() {
	c.useDefaultRenamer = true
	c.testWithInversion(
		`{let a;
  {
    let a;
  }}
`,
		`{let a;
  {
  let a$jscomp$1;
  }}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testBlockScopesWithContextWithInversion2()
func (c *MakeDeclaredNamesUniqueTest) TestBlockScopesWithContextWithInversion2() {
	c.useDefaultRenamer = true
	// function declarations are block-scoped
	c.testWithInversion(
		`function foo() {
  function bar() {
    return 1;
  }
}
function boo() {
  function bar() {
    return 2;
  }
}
`,
		`function foo() {
  function bar() {
    return 1;
  }
}
function boo() {
  function bar$jscomp$1() {
    return 2;
  }
}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testBlockScopesWithContextWithInversion3()
func (c *MakeDeclaredNamesUniqueTest) TestBlockScopesWithContextWithInversion3() {
	c.useDefaultRenamer = true
	c.Test(
		`function foo() {
  function bar() {
    return 1;
  }
  if (true) {
    function bar() {
      return 2;
    }
  }
}
`,
		`function foo() {
  function bar() {
    return 1;
  }
  if (true) {
    function bar$jscomp$1() {
      return 2;
    }
  }
}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testBlockScopesWithContextWithInversion4()
func (c *MakeDeclaredNamesUniqueTest) TestBlockScopesWithContextWithInversion4() {
	c.useDefaultRenamer = true
	c.Test(
		`var f1=function(){
  var x
};
(function() {
  function f2() {
    alert(x)
  }
  {
    var x=0
  }
  f2()
})()
`,
		`var f1=function(){
  var x
};
(function() {
  function f2() {
    alert(x$jscomp$1)
  }
  {
    var x$jscomp$1=0
  }
  f2()
})()
`)
}

// Java: MakeDeclaredNamesUniqueTest.testBlockScopesWithContextWithInversion5()
func (c *MakeDeclaredNamesUniqueTest) TestBlockScopesWithContextWithInversion5() {
	c.useDefaultRenamer = true
	c.TestSame(
		`if (true) {
  function f(){};
}
f();
`)
}

// Java: MakeDeclaredNamesUniqueTest.testBlockScopesWithoutContext()
func (c *MakeDeclaredNamesUniqueTest) TestBlockScopesWithoutContext() {
	c.useDefaultRenamer = false
	c.Test(
		`{
  function foo() {return 1;}
  if (true) {
    function foo() {return 2;}
  }
}
`,
		`{
  function foo$jscomp$unique_0() {return 1;}
  if (true) {
    function foo$jscomp$unique_1() {return 2;}
  }
}
`)

	c.Test(
		`function foo(x) {
  return foo(x) - 1;
}
`,
		`function foo$jscomp$unique_0(x$jscomp$unique_1) {
  return foo$jscomp$unique_0(x$jscomp$unique_1) - 1;
}
`)

	c.Test(
		`export function foo(x) {
  return foo(x) - 1;
}
`,
		`export function foo$jscomp$unique_0(x$jscomp$unique_1) {
  return foo$jscomp$unique_0(x$jscomp$unique_1) - 1;
}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testRecursiveFunctionsWithContextWithInversion()
func (c *MakeDeclaredNamesUniqueTest) TestRecursiveFunctionsWithContextWithInversion() {
	c.useDefaultRenamer = true
	c.testSameWithInversion(
		`function foo(x) {
  return foo(x) - 1;
}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testInvertShadowedParameterNames()
func (c *MakeDeclaredNamesUniqueTest) TestInvertShadowedParameterNames() {
	c.useDefaultRenamer = true
	c.testWithInversion(
		`var p;
function f(p) {
  return function g(p) {
    return p;
  }
}
`,
		`var p;
function f(p$jscomp$1) {
  return function g(p$jscomp$2) {
    return p$jscomp$2;
  }
}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testArrowFunctionWithContextWithInversion()
func (c *MakeDeclaredNamesUniqueTest) TestArrowFunctionWithContextWithInversion() {
	c.useDefaultRenamer = true
	c.testWithInversion(
		`function foo() {
  var f = (x) => x;
  return f(1);
}
function boo() {
  var f = (x) => x;
  return f(2);
}
`,
		`function foo() {
  var f = (x) => x;
  return f(1);
}
function boo() {
  var f$jscomp$1 = (x$jscomp$1) => x$jscomp$1;
  return f$jscomp$1(2);
}
`)

	c.testWithInversion(
		`function foo() {
  var f = (x, ...y) => x + y[0];
  return f(1, 2);
}
function boo() {
  var f = (x, ...y) => x + y[0];
  return f(1, 2);
}
`,
		`function foo() {
  var f = (x, ...y) => x + y[0];
  return f(1, 2);
}
function boo() {
  var f$jscomp$1 = (x$jscomp$1, ...y$jscomp$1) => x$jscomp$1 + y$jscomp$1[0];
  return f$jscomp$1(1, 2);
}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testDefaultParameterWithContextWithInversion1()
func (c *MakeDeclaredNamesUniqueTest) TestDefaultParameterWithContextWithInversion1() {
	c.useDefaultRenamer = true
	c.testWithInversion(
		`function foo(x = 1) {
  return x;
}
function boo(x = 1) {
  return x;
}
`,
		`function foo(x = 1) {
  return x;
}
function boo(x$jscomp$1 = 1) {
  return x$jscomp$1;
}
`)

	c.testSameWithInversion(
		`function foo(x = 1, y = x) {
  return x + y;
}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testDefaultParameterWithContextWithInversion2()
func (c *MakeDeclaredNamesUniqueTest) TestDefaultParameterWithContextWithInversion2() {
	c.useDefaultRenamer = true

	// Parameter default values don't see the scope of the body
	// Methods or functions defined "inside" parameter default values don't see the local variables
	// of the body.
	c.testWithInversion(
		`let x = 'outer';
function foo(bar = baz => x) {
  let x = 'inner';
  console.log(bar());
}
`,
		`let x = 'outer';
function foo(bar = baz => x) {
  let x$jscomp$1 = 'inner';
  console.log(bar());
}
`)

	c.testWithInversion(
		`const x = 'outer';
function foo(a = x) {
  const x = 'inner';
  return a;
}
`,
		`const x = 'outer';
function foo(a = x) {
  const x$jscomp$1 = 'inner';
  return a;
}
`)

	c.testWithInversion(
		`const x = 'outerouter';
{
  const x = 'outer';
  function foo(a = x) {
    return a;
  }
foo();
}
`,
		`const x = 'outerouter';
{
  const x$jscomp$1 = 'outer';
  function foo(a = x$jscomp$1) {
    return a;
  }
foo();
}
`)

	c.testSameWithInversion(
		`function foo(x, y = x) {
  return x + y;
}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testObjectLiteralsWithContextWithInversion()
func (c *MakeDeclaredNamesUniqueTest) TestObjectLiteralsWithContextWithInversion() {
	c.useDefaultRenamer = true
	c.testWithInversion(
		`function foo({x:y}) {
  return y;
}
function boo({x:y}) {
  return y;
}
`,
		`function foo({x:y}) {
  return y;
}
function boo({x:y$jscomp$1}) {
  return y$jscomp$1
}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testExportedOrImportedNamesAreUntouched()
func (c *MakeDeclaredNamesUniqueTest) TestExportedOrImportedNamesAreUntouched() {
	// The eventual desired behavior is that none of the 'a's in the following test cases
	// are renamed to a$jscomp$1. Rewrite this test after that behavior is implemented.
	c.useDefaultRenamer = true
	c.Test(
		c.Srcs("var a;", "let a; export {a as a};"),
		c.Expected("var a;", "let a$jscomp$1; export {a$jscomp$1 as a};"))

	c.Test(
		c.Srcs("var a;", "import {a as a} from './foo.js'; let b = a;"),
		c.Expected("var a;", "import {a as a$jscomp$1} from './foo.js'; let b = a$jscomp$1;"))
}

// Java: MakeDeclaredNamesUniqueTest.testTwoMethodsInTheSameFileWithSameLocalNames()
func (c *MakeDeclaredNamesUniqueTest) TestTwoMethodsInTheSameFileWithSameLocalNames() {
	c.useDefaultRenamer = true
	// Verify same local names in 2 different files get different new names.
	// The ContextualRenamer renames an "oldName" by adding "$jscomp$\n" + "id" as a suffix string.
	// So, when another declaration containing "$jscomp$id" exists at any other source location in
	// the entire JS program (due to a prior renaming), the ContextualRenamer should not generate
	// that same name when renaming this declaration.
	c.Test(
		"function foo() {var a; a;} function bar() {let a; let a$jscomp$1; a + a$jscomp$1;}",
		`function foo() {var a; a;}
function bar() {
  let a$jscomp$1; let a$jscomp$1$jscomp$1;
  a$jscomp$1 + a$jscomp$1$jscomp$1;
}
`)

	c.Test(
		"function bar() {let a; let a$jscomp$1; a + a$jscomp$1;} function foo() {var a; a;}",
		`function bar() {
  let a; let a$jscomp$1;
  a + a$jscomp$1;
}
function foo() {var a$jscomp$2; a$jscomp$2;}
`)

	// tests when name with $jscomp$1 suffix comes first
	c.Test(
		"function bar() {let a$jscomp$1; let a; a + a$jscomp$1;} function foo() {var a; a;}",
		`function bar() {
  let a$jscomp$1; let a;
  a + a$jscomp$1;
}
function foo() {var a$jscomp$2; a$jscomp$2;}
`)

	c.Test(
		`function bar() {
  let a; let a$jscomp$1;
  a + a$jscomp$1;
}
function foo() {
// tests when a$jscomp$2 declared later in the same scope
  var a; a; var a$jscomp$2; a$jscomp$2;
}
`,
		`function bar() {
  let a; let a$jscomp$1; a + a$jscomp$1;
}
function foo() {
  var a$jscomp$2; a$jscomp$2;
  var a$jscomp$2$jscomp$1; a$jscomp$2$jscomp$1;
}
`)

	c.Test(
		`function bar() {
  let a; let a$jscomp$1; a + a$jscomp$1;
}
function foo() {
// tests when a$jscomp$2 declared first in the same scope
  var a$jscomp$2; a$jscomp$2; var a; a;
}
`,
		`function bar() {
  let a; let a$jscomp$1; a + a$jscomp$1;
}
function foo() {
  var a$jscomp$2; a$jscomp$2; var a$jscomp$3; a$jscomp$3;
}
`)

	c.Test(
		`function bar() {
  let a; let a$jscomp$1; a + a$jscomp$1;
}
function foo() {
// tests when a$jscomp$2 is declared in another scope
  var a; a;
}
function baz() {
  var a$jscomp$2; a$jscomp$2;
}
`,
		`function bar() {
  let a; let a$jscomp$1; a + a$jscomp$1;
}
function foo() {
  var a$jscomp$2; a$jscomp$2;
}
function baz() {
  var a$jscomp$2$jscomp$1; a$jscomp$2$jscomp$1;
}
`)
}

// Java: MakeDeclaredNamesUniqueTest.testTwoFilesWithSameLocalNames()
func (c *MakeDeclaredNamesUniqueTest) TestTwoFilesWithSameLocalNames() {
	c.useDefaultRenamer = true
	// Verify same local names in 2 different files get different new names
	c.Test(
		c.Srcs(
			"function foo() {var a; a;}",
			"function bar() {let a; let a$jscomp$1; a + a$jscomp$1;}"),
		c.Expected(
			"function foo() {var a; a;}",
			`function bar() {let a$jscomp$1; let a$jscomp$1$jscomp$1; a$jscomp$1 +
 a$jscomp$1$jscomp$1;}
`))
}

// Java: MakeDeclaredNamesUniqueTest.testImportStarWithInversion()
func (c *MakeDeclaredNamesUniqueTest) TestImportStarWithInversion() {
	c.useDefaultRenamer = true
	c.testWithInversionOfArrays(
		[]string{
			"let a = 5;", //
			"import * as a          from './a.js'; const TAU = 2 * a.PI;",
		},
		[]string{
			"let a = 5;", //
			"import * as a$jscomp$1 from './a.js'; const TAU = 2 * a$jscomp$1.PI",
		})
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: MakeDeclaredNamesUniqueTest.assertOnChange_throwsException()
func (c *MakeDeclaredNamesUniqueTest) TestAssertOnChange_throwsException() {
	c.useDefaultRenamer = true
	c.assertOnChange = true

	e :=
		AssertThrows[jutil.IRuntimeException](
			func() { c.TestNoWarning("var a; function foo() { var a = 1; } ") })
	AssertThat(e.GetMessage()).Contains("NAME a")
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: MakeDeclaredNamesUniqueTest.assertOnChange_noExceptionIfNothingChanges()
func (c *MakeDeclaredNamesUniqueTest) TestAssertOnChange_noExceptionIfNothingChanges() {
	c.useDefaultRenamer = true
	c.assertOnChange = true

	c.TestSame("const x = 1; function foo() { const y = 2; }")
}

// Go-only: test-helper — exercise the renamers directly (no Java test).
func TestMakeDeclaredNamesUnique_renamers(t *testing.T) {
	root := NewMakeDeclaredNamesUnique_ContextualRenamer()
	root.AddDeclaredName("a", false)
	if got := root.GetReplacementName("a"); got != nil {
		t.Errorf("global a renamed to %q", *got)
	}
	fn := rhino.IR_Function(rhino.IR_Name("f"), rhino.IR_ParamList(), rhino.IR_Block())
	child := root.CreateForChildScope(fn, false)
	child.AddDeclaredName("a", false)
	child.AddDeclaredName("arguments", false)
	if got := child.GetReplacementName("a"); got == nil || *got != "a$jscomp$1" {
		t.Errorf("local a: got %v", got)
	}
	if got := child.GetReplacementName("arguments"); got != nil {
		t.Errorf("arguments renamed to %q", *got)
	}
	if !strings.HasPrefix(root.String(), "ContextualRenamer{scopeRoot=null, nameUsage=[") {
		t.Errorf("toString: %s", root.String())
	}
	if got := MakeDeclaredNamesUnique_ContextualRenameInverter_GetOriginalName("x$jscomp$12"); got != "x" {
		t.Errorf("getOriginalName: %q", got)
	}
	if got := MakeDeclaredNamesUnique_ContextualRenameInverter_GetOriginalName("$jscomp$12"); got != "$jscomp$12" {
		t.Errorf("getOriginalName: %q", got)
	}
	targetted := NewMakeDeclaredNamesUnique_TargettedRenamer(
		NewMakeDeclaredNamesUnique_ContextualRenamer(), jutil.NewLinkedHashSet("t"))
	targetted.AddDeclaredName("t", false)
	targetted.AddDeclaredName("u", false)
	if targetted.GetReplacementName("u") != nil {
		t.Errorf("untargetted name replaced")
	}
}
