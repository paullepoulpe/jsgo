// Port of com/google/javascript/jscomp/VariableReferenceCheckTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/VariableReferenceCheckTest.java)

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
)

// VariableReferenceCheckTest tests that warnings are generated in
// appropriate cases and appropriate cases only by VariableReferenceCheck.
type VariableReferenceCheckTest struct {
	CompilerTestCaseBase
}

// Java: VariableReferenceCheckTest.LET_RUN
const variableReferenceCheckTest_LET_RUN = "let a = 1; let b = 2; let c = a + b, d = c;"

// Java: VariableReferenceCheckTest.VARIABLE_RUN
const variableReferenceCheckTest_VARIABLE_RUN = "var a = 1; var b = 2; var c = a + b, d = c;"

// Java: VariableReferenceCheckTest()
func NewVariableReferenceCheckTest() *VariableReferenceCheckTest {
	c := &VariableReferenceCheckTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestVariableReferenceCheckTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewVariableReferenceCheckTest() })
}

// Java: VariableReferenceCheckTest.getProcessor(Compiler)
func (c *VariableReferenceCheckTest) GetProcessor(compiler *Compiler) CompilerPass {
	// Treats bad reads as errors, and reports bad write warnings.
	return NewVariableReferenceCheck(compiler)
}

// Java: VariableReferenceCheckTest.testWithImportMeta()
func (c *VariableReferenceCheckTest) TestWithImportMeta() {
	// just to confirm that presence of import.meta does not cause a compiler crash
	c.TestSame("export function g() { return import.meta; }")
}

// Java: VariableReferenceCheckTest.testDoubleTryCatch()
func (c *VariableReferenceCheckTest) TestDoubleTryCatch() {
	c.TestSame(
		`function g() {
  return f;

  function f() {
    try {
    } catch (e) {
      alert(e);
    }
    try {
    } catch (e) {
      alert(e);
    }
  }
}
`)
}

// Java: VariableReferenceCheckTest.testDoubleTryCatch_withES6Modules()
func (c *VariableReferenceCheckTest) TestDoubleTryCatch_withES6Modules() {
	c.TestSame(
		`export
function g() {
  return f;

  function f() {
    try {
    } catch (e) {
      alert(e);
    }
    try {
    } catch (e) {
      alert(e);
    }
  }
}
`)
}

// Java: VariableReferenceCheckTest.testCorrectCode()
func (c *VariableReferenceCheckTest) TestCorrectCode() {
	c.AssertNoWarning("function foo(d) { (function() { d.foo(); }); d.bar(); } ")
	c.AssertNoWarning("function foo() { bar(); } function bar() { foo(); } ")
	c.AssertNoWarning("function f(d) { d = 3; }")
	c.AssertNoWarning(variableReferenceCheckTest_VARIABLE_RUN)
	c.AssertNoWarning("if (a) { var x; }")
	c.AssertNoWarning("function f() { " + variableReferenceCheckTest_VARIABLE_RUN + "}")
	c.AssertNoWarning(variableReferenceCheckTest_LET_RUN)
	c.AssertNoWarning("function f() { " + variableReferenceCheckTest_LET_RUN + "}")
	c.AssertNoWarning("try { let e; } catch (e) { let x; }")
}

// Java: VariableReferenceCheckTest.testCorrectCode_withES6Modules()
func (c *VariableReferenceCheckTest) TestCorrectCode_withES6Modules() {
	c.AssertNoWarning("export function foo(d) { (function() { d.foo(); }); d.bar(); } ")
}

// Java: VariableReferenceCheckTest.testCorrectShadowing()
func (c *VariableReferenceCheckTest) TestCorrectShadowing() {
	c.AssertNoWarning(variableReferenceCheckTest_VARIABLE_RUN + "function f() { " + variableReferenceCheckTest_VARIABLE_RUN + "}")
}

// Java: VariableReferenceCheckTest.testCorrectShadowing_withES6Modules()
func (c *VariableReferenceCheckTest) TestCorrectShadowing_withES6Modules() {
	c.AssertNoWarning(variableReferenceCheckTest_VARIABLE_RUN + "export function f() { " + variableReferenceCheckTest_VARIABLE_RUN + "}")
}

// Java: VariableReferenceCheckTest.testCorrectRedeclare()
func (c *VariableReferenceCheckTest) TestCorrectRedeclare() {
	c.AssertNoWarning("function f() { if (1) { var a = 2; } else { var a = 3; } }")
}

// Java: VariableReferenceCheckTest.testCorrectRedeclare_withES6Modules()
func (c *VariableReferenceCheckTest) TestCorrectRedeclare_withES6Modules() {
	c.AssertNoWarning("export function f() { if (1) { var a = 2; } else { var a = 3; } }")
}

// Java: VariableReferenceCheckTest.testCorrectRecursion()
func (c *VariableReferenceCheckTest) TestCorrectRecursion() {
	c.AssertNoWarning("function f() { var x = function() { x(); }; }")
}

// Java: VariableReferenceCheckTest.testCorrectRecursion_withES6Modules()
func (c *VariableReferenceCheckTest) TestCorrectRecursion_withES6Modules() {
	c.AssertNoWarning("export function f() { var x = function() { x(); }; }")
}

// Java: VariableReferenceCheckTest.testCorrectCatch()
func (c *VariableReferenceCheckTest) TestCorrectCatch() {
	c.AssertNoWarning("function f() { try { var x = 2; } catch (x) {} }")
	c.AssertNoWarning("function f(e) { e = 3; try {} catch (e) {} }")
}

// Java: VariableReferenceCheckTest.testCorrectCatch_withES6Modules()
func (c *VariableReferenceCheckTest) TestCorrectCatch_withES6Modules() {
	c.AssertNoWarning("export function f() { try { var x = 2; } catch (x) {} }")
}

// Java: VariableReferenceCheckTest.testRedeclare()
func (c *VariableReferenceCheckTest) TestRedeclare() {
	// Only test local scope since global scope is covered elsewhere
	c.AssertRedeclare("function f() { var a = 2; var a = 3; }")
	c.AssertRedeclare("function f(a) { var a = 2; }")
	c.AssertRedeclare("function f(a) { if (!a) var a = 6; }")
	// NOTE: We decided to not give warnings to the following cases. The function won't be
	// overwritten at runtime anyway.
	c.AssertNoWarning("function f() { var f = 1; }")
	c.AssertNoWarning("function f() { let f = 1; }")
}

// Java: VariableReferenceCheckTest.testRedeclare_withES6Modules()
func (c *VariableReferenceCheckTest) TestRedeclare_withES6Modules() {
	c.AssertRedeclare("export function f() { var a = 2; var a = 3; }")
	c.AssertNoWarning("export function f() { let f = 1; }")
	// In an ES6 module vars are in the module scope, not global, so they are covered here.
	c.AssertRedeclare("export var a = 2; var a = 3;")
	c.AssertRedeclare("export var a = 2; if (a) var a = 3;")
	c.AssertRedeclare("function f() {} function f() {} export {f};")
}

// Java: VariableReferenceCheckTest.testIssue166a()
func (c *VariableReferenceCheckTest) TestIssue166a() {
	c.AssertRedeclareError("try { throw 1 } catch(e) { /** @suppress {duplicate} */ var e=2 }")
}

// Java: VariableReferenceCheckTest.testIssue166b()
func (c *VariableReferenceCheckTest) TestIssue166b() {
	c.AssertRedeclareError(
		"function a() { try { throw 1 } catch(e) { /** @suppress {duplicate} */ var e=2 } };")
}

// Java: VariableReferenceCheckTest.testIssue166b_withES6Modules()
func (c *VariableReferenceCheckTest) TestIssue166b_withES6Modules() {
	c.AssertRedeclareError(
		`export function a() {
  try {
    throw 1
  } catch (e) {
      /** @suppress {duplicate} */
      var e = 2
  }
};
`)
}

// Java: VariableReferenceCheckTest.testIssue166c()
func (c *VariableReferenceCheckTest) TestIssue166c() {
	c.AssertRedeclareError(
		"var e = 0; try { throw 1 } catch(e) { /** @suppress {duplicate} */ var e=2 }")
}

// Java: VariableReferenceCheckTest.testIssue166d()
func (c *VariableReferenceCheckTest) TestIssue166d() {
	c.AssertRedeclareError(
		`function a() {
  var e = 0; try { throw 1 } catch(e) {
    /** @suppress {duplicate} */ var e = 2;
  }
};
`)
}

// Java: VariableReferenceCheckTest.testIssue166e()
func (c *VariableReferenceCheckTest) TestIssue166e() {
	c.TestSame("var e = 2; try { throw 1 } catch(e) {}")
}

// Java: VariableReferenceCheckTest.testIssue166e_withES6Modules()
func (c *VariableReferenceCheckTest) TestIssue166e_withES6Modules() {
	c.TestSame("export var e = 2; try { throw 1 } catch(e) {}")
}

// Java: VariableReferenceCheckTest.testIssue166f()
func (c *VariableReferenceCheckTest) TestIssue166f() {
	c.TestSame(
		`function a() {
  var e = 2;
  try { throw 1 } catch(e) {}
}
`)
}

// Java: VariableReferenceCheckTest.testEarlyReference()
func (c *VariableReferenceCheckTest) TestEarlyReference() {
	c.AssertEarlyReferenceWarning("function f() { a = 2; var a = 3; }")
}

// Java: VariableReferenceCheckTest.testEarlyReference_withES6Modules()
func (c *VariableReferenceCheckTest) TestEarlyReference_withES6Modules() {
	c.AssertEarlyReferenceWarning("export function f() { a = 2; var a = 3; }")
}

// Java: VariableReferenceCheckTest.testCorrectEarlyReference()
func (c *VariableReferenceCheckTest) TestCorrectEarlyReference() {
	c.AssertNoWarning("var goog = goog || {}")
	c.AssertNoWarning("var google = google || window['google'] || {}")
	c.AssertNoWarning("function f() { a = 2; } var a = 2;")
}

// Java: VariableReferenceCheckTest.testCorrectEarlyReferenceLogicalAssignment()
func (c *VariableReferenceCheckTest) TestCorrectEarlyReferenceLogicalAssignment() {
	// These patterns are normalized away
	c.AssertNoWarning("function f() { a ||= {}; } let a;")
	c.AssertNoWarning("function f() { a &&= {}; } let a;")
	c.AssertNoWarning("function f() { a ??= {}; } let a;")
}

// Java: VariableReferenceCheckTest.testCorrectEarlyReference_withES6Modules()
func (c *VariableReferenceCheckTest) TestCorrectEarlyReference_withES6Modules() {
	c.AssertNoWarning("export function f() { a = 2; } var a = 2;")
}

// Java: VariableReferenceCheckTest.testUnreferencedBleedingFunction()
func (c *VariableReferenceCheckTest) TestUnreferencedBleedingFunction() {
	c.AssertNoWarning("var x = function y() {}")
	c.AssertNoWarning("var x = function y() {}; var y = 1;")
}

// Java: VariableReferenceCheckTest.testUnreferencedBleedingFunction_withES6Modules()
func (c *VariableReferenceCheckTest) TestUnreferencedBleedingFunction_withES6Modules() {
	c.AssertNoWarning("export var x = function y() {}")
}

// Java: VariableReferenceCheckTest.testReferencedBleedingFunction()
func (c *VariableReferenceCheckTest) TestReferencedBleedingFunction() {
	c.AssertNoWarning("var x = function y() { return y(); }")
}

// Java: VariableReferenceCheckTest.testReferencedBleedingFunction_withES6Modules()
func (c *VariableReferenceCheckTest) TestReferencedBleedingFunction_withES6Modules() {
	c.AssertNoWarning("export var x = function y() { return y(); }")
}

// Java: VariableReferenceCheckTest.testVarShadowsFunctionName()
func (c *VariableReferenceCheckTest) TestVarShadowsFunctionName() {
	c.AssertNoWarning("var x = function y() { var y; }")
	c.AssertNoWarning("var x = function y() { let y; }")
}

// Java: VariableReferenceCheckTest.testVarShadowsFunctionName_withES6Modules()
func (c *VariableReferenceCheckTest) TestVarShadowsFunctionName_withES6Modules() {
	c.AssertNoWarning("export var x = function y() { var y; }")
	c.AssertNoWarning("export var x = function y() { let y; }")
}

// Java: VariableReferenceCheckTest.testDoubleDeclaration()
func (c *VariableReferenceCheckTest) TestDoubleDeclaration() {
	c.AssertRedeclare("function x(y) { if (true) { var y; } }")
}

// Java: VariableReferenceCheckTest.testDoubleDeclaration2()
func (c *VariableReferenceCheckTest) TestDoubleDeclaration2() {
	c.AssertRedeclare("function x() { var y; if (true) { var y; } }")
}

// Java: VariableReferenceCheckTest.testDoubleDeclaration_withES6Modules()
func (c *VariableReferenceCheckTest) TestDoubleDeclaration_withES6Modules() {
	c.AssertRedeclare("export function x(y) { if (true) { var y; } }")
}

// Java: VariableReferenceCheckTest.testHoistedFunction1()
func (c *VariableReferenceCheckTest) TestHoistedFunction1() {
	c.AssertNoWarning("f(); function f() {}")
}

// Java: VariableReferenceCheckTest.testHoistedFunction2()
func (c *VariableReferenceCheckTest) TestHoistedFunction2() {
	c.AssertNoWarning("function g() { f(); function f() {} }")
}

// Java: VariableReferenceCheckTest.testHoistedFunction_withES6Modules()
func (c *VariableReferenceCheckTest) TestHoistedFunction_withES6Modules() {
	c.AssertNoWarning("export function g() { f(); function f() {} }")
}

// Java: VariableReferenceCheckTest.testNonHoistedFunction()
func (c *VariableReferenceCheckTest) TestNonHoistedFunction() {
	c.AssertEarlyReferenceWarning("if (true) { f(); function f() {} }")
}

// Java: VariableReferenceCheckTest.testNonHoistedFunction2()
func (c *VariableReferenceCheckTest) TestNonHoistedFunction2() {
	c.AssertNoWarning("if (false) { function f() {} f(); }")
}

// Java: VariableReferenceCheckTest.testNonHoistedFunction3()
func (c *VariableReferenceCheckTest) TestNonHoistedFunction3() {
	c.AssertNoWarning("function g() { if (false) { function f() {} f(); }}")
}

// Java: VariableReferenceCheckTest.testNonHoistedFunction4()
func (c *VariableReferenceCheckTest) TestNonHoistedFunction4() {
	c.AssertNoWarning("if (false) { function f() {} }  f();")
}

// Java: VariableReferenceCheckTest.testNonHoistedFunction5()
func (c *VariableReferenceCheckTest) TestNonHoistedFunction5() {
	c.AssertNoWarning("function g() { if (false) { function f() {} }  f(); }")
}

// Java: VariableReferenceCheckTest.testNonHoistedFunction6()
func (c *VariableReferenceCheckTest) TestNonHoistedFunction6() {
	c.AssertEarlyReferenceWarning("if (false) { f(); function f() {} }")
}

// Java: VariableReferenceCheckTest.testNonHoistedFunction7()
func (c *VariableReferenceCheckTest) TestNonHoistedFunction7() {
	c.AssertEarlyReferenceWarning("function g() { if (false) { f(); function f() {} }}")
}

// Java: VariableReferenceCheckTest.testNonHoistedFunction_withES6Modules()
func (c *VariableReferenceCheckTest) TestNonHoistedFunction_withES6Modules() {
	c.AssertEarlyReferenceWarning("export function g() { if (false) { f(); function f() {} }}")
}

// Java: VariableReferenceCheckTest.testNonHoistedRecursiveFunction1()
func (c *VariableReferenceCheckTest) TestNonHoistedRecursiveFunction1() {
	c.AssertNoWarning("if (false) { function f() { f(); }}")
}

// Java: VariableReferenceCheckTest.testNonHoistedRecursiveFunction2()
func (c *VariableReferenceCheckTest) TestNonHoistedRecursiveFunction2() {
	c.AssertNoWarning("function g() { if (false) { function f() { f(); }}}")
}

// Java: VariableReferenceCheckTest.testNonHoistedRecursiveFunction3()
func (c *VariableReferenceCheckTest) TestNonHoistedRecursiveFunction3() {
	c.AssertNoWarning("function g() { if (false) { function f() { f(); g(); }}}")
}

// Java: VariableReferenceCheckTest.testNonHoistedRecursiveFunction_withES6Modules()
func (c *VariableReferenceCheckTest) TestNonHoistedRecursiveFunction_withES6Modules() {
	c.AssertNoWarning("export function g() { if (false) { function f() { f(); g(); }}}")
}

// Java: VariableReferenceCheckTest.testForOf()
func (c *VariableReferenceCheckTest) TestForOf() {
	c.AssertEarlyReferenceError("for (let x of []) { console.log(x); let x = 123; }")
	c.AssertNoWarning("for (let x of []) { let x; }")
}

// Java: VariableReferenceCheckTest.testForAwaitOf()
func (c *VariableReferenceCheckTest) TestForAwaitOf() {
	c.AssertEarlyReferenceError(
		"async () => { for await (let x of []) { console.log(x); let x = 123; } }")
	c.AssertNoWarning("async () => { for (let x of []) { let x; } }")
}

// Java: VariableReferenceCheckTest.testDestructuringInFor()
func (c *VariableReferenceCheckTest) TestDestructuringInFor() {
	c.TestSame("for (let [key, val] of X){}")
	c.TestSame("for (let [key, [nestKey, nestVal], val] of X){}")
	c.TestSame("var {x: a, y: b} = {x: 1, y: 2}; a++; b++;")
	c.TestWarning("a++; var {x: a} = {x: 1};", VariableReferenceCheck_EARLY_REFERENCE)
}

// Java: VariableReferenceCheckTest.testSuppressDuplicate_first()
func (c *VariableReferenceCheckTest) TestSuppressDuplicate_first() {
	code := "/** @suppress {duplicate} */ var google; var google"
	c.TestSame(code)
}

// Java: VariableReferenceCheckTest.testSuppressDuplicate_second()
func (c *VariableReferenceCheckTest) TestSuppressDuplicate_second() {
	code := "var google; /** @suppress {duplicate} */ var google"
	c.TestSame(code)
}

// Java: VariableReferenceCheckTest.testSuppressDuplicate_fileoverview()
func (c *VariableReferenceCheckTest) TestSuppressDuplicate_fileoverview() {
	code := `/** @fileoverview @suppress {duplicate} */
/** @type {?} */ var google;
 var google
`
	c.TestSame(code)
}

// Java: VariableReferenceCheckTest.testNoWarnDuplicateInExterns2()
func (c *VariableReferenceCheckTest) TestNoWarnDuplicateInExterns2() {
	// Verify we don't complain about early references in externs
	externs := "window; var window;"
	code := ""
	c.TestSame(c.Externs(externs), c.Srcs(code))
}

// Java: VariableReferenceCheckTest.testNoWarnDuplicateInExterns_withES6Modules()
func (c *VariableReferenceCheckTest) TestNoWarnDuplicateInExterns_withES6Modules() {
	externs := "export var google; /** @suppress {duplicate} */ var google"
	code := ""
	c.TestSame(c.Externs(externs), c.Srcs(code))
}

// Java: VariableReferenceCheckTest.testImportStar()
func (c *VariableReferenceCheckTest) TestImportStar() {
	c.TestSame(
		c.Srcs(
			SourceFile_FromCode("foo.js", ""),
			SourceFile_FromCode("bar.js", "import * as ns from './foo.js'")))
}

// Java: VariableReferenceCheckTest.testUnusedCompoundAssign_withES6Modules()
func (c *VariableReferenceCheckTest) TestUnusedCompoundAssign_withES6Modules() {
	c.AssertNoWarning(
		`export function f(elapsed) {
  let fakeMs = 0;
  stubs.replace(Date, 'now', () => fakeMs -= elapsed);
}
`)
}

// Java: VariableReferenceCheckTest.testGoogModule_duplicateRequire()
func (c *VariableReferenceCheckTest) TestGoogModule_duplicateRequire() {
	c.AssertRedeclareError(
		"goog.module('bar'); const X = goog.require('foo.X'); const X = goog.require('foo.X');")
	c.AssertRedeclareError(
		"goog.module('bar'); let X = goog.require('foo.X'); let X = goog.require('foo.X');")
	c.AssertRedeclareError(
		"goog.module('bar'); const X = goog.require('foo.X'); let X = goog.require('foo.X');")
	c.AssertRedeclareError(
		"goog.module('bar'); let X = goog.require('foo.X'); const X = goog.require('foo.X');")
}

// Java: VariableReferenceCheckTest.testGoogProvide_ok()
func (c *VariableReferenceCheckTest) TestGoogProvide_ok() {
	c.AssertNoWarning("goog.provide('foo');")
	c.AssertNoWarning("goog.provide('foo'); foo = 0;")
	c.AssertNoWarning("goog.provide('foo'); var foo = 0;")
	c.AssertNoWarning("goog.provide('foo.bar');")
	c.AssertNoWarning("goog.provide('foo.bar'); foo.bar = 0;")
}

// Java: VariableReferenceCheckTest.testUndeclaredLet()
func (c *VariableReferenceCheckTest) TestUndeclaredLet() {
	c.AssertEarlyReferenceError("if (a) { x = 3; let x;}")
	c.AssertEarlyReferenceError(
		`var x = 1;
if (true) {
  x++;
  let x = 3;
}
`)
}

// Java: VariableReferenceCheckTest.testUndeclaredLet_withES6Modules()
func (c *VariableReferenceCheckTest) TestUndeclaredLet_withES6Modules() {
	c.AssertEarlyReferenceError(
		`export var x = 1;
if (true) {
  x++;
  let x = 3;
}
`)
}

// Java: VariableReferenceCheckTest.testUndeclaredConst()
func (c *VariableReferenceCheckTest) TestUndeclaredConst() {
	c.AssertEarlyReferenceError("if (a) { x = 3; const x = 3;}")
	// For the following, IE 11 gives "Assignment to const", but technically
	// they are also undeclared references, which get caught in the first place.
	c.AssertEarlyReferenceError(
		`var x = 1;
if (true) {
  x++;
  const x = 3;
}
`)
	c.AssertEarlyReferenceError("a = 1; const a = 0;")
	c.AssertEarlyReferenceError("a++; const a = 0;")
}

// Java: VariableReferenceCheckTest.testIllegalLetShadowing()
func (c *VariableReferenceCheckTest) TestIllegalLetShadowing() {
	c.AssertRedeclareError("if (a) { let x; var x;}")
	c.AssertRedeclareError("if (a) { let x; let x;}")
	c.AssertRedeclareError(
		`function f() {
  let x;
  if (a) {
    var x;
  }
}
`)
	c.AssertNoWarning(
		`function f() {
  if (a) {
    let x;
  }
  var x;
}
`)
	c.AssertNoWarning(
		`function f() {
  if (a) { let x; }
  if (b) { var x; }
}
`)
	c.AssertRedeclareError("let x; var x;")
	c.AssertRedeclareError("var x; let x;")
	c.AssertRedeclareError("let x; let x;")
}

// Java: VariableReferenceCheckTest.testIllegalLetShadowing_withES6Modules()
func (c *VariableReferenceCheckTest) TestIllegalLetShadowing_withES6Modules() {
	c.AssertRedeclareError(
		`export function f() {
  let x;
  if (a) {
    var x;
  }
}
`)
	c.AssertNoWarning(
		`export function f() {
  if (a) {
    let x;
  }
  var x;
}
`)
	c.AssertRedeclareError("export let x; var x;")
}

// Java: VariableReferenceCheckTest.testDuplicateLetConst()
func (c *VariableReferenceCheckTest) TestDuplicateLetConst() {
	c.AssertRedeclareError("let x, x;")
	c.AssertRedeclareError("const x = 0, x = 0;")
}

// Java: VariableReferenceCheckTest.testRedeclareInLabel()
func (c *VariableReferenceCheckTest) TestRedeclareInLabel() {
	c.AssertRedeclareGlobal("a: var x, x;")
}

// Java: VariableReferenceCheckTest.testIllegalBlockScopedEarlyReference()
func (c *VariableReferenceCheckTest) TestIllegalBlockScopedEarlyReference() {
	c.AssertEarlyReferenceError("let x = x")
	c.AssertEarlyReferenceError("let [x] = x")
	c.AssertEarlyReferenceError("const x = x")
	c.AssertEarlyReferenceError("let x = x || 0")
	c.AssertEarlyReferenceError("const x = x || 0")
	// In the following cases, "x" might not be reachable but we warn anyways
	c.AssertEarlyReferenceError("let x = expr || x")
	c.AssertEarlyReferenceError("const x = expr || x")
	c.AssertEarlyReferenceError("X; class X {};")
}

// Java: VariableReferenceCheckTest.testIllegalConstShadowing()
func (c *VariableReferenceCheckTest) TestIllegalConstShadowing() {
	c.AssertRedeclareError("if (a) { const x = 3; var x;}")
	c.AssertRedeclareError(
		`function f() {
  const x = 3;
  if (a) {
    var x;
  }
}
`)
}

// Java: VariableReferenceCheckTest.testIllegalConstShadowing_withES6Modules()
func (c *VariableReferenceCheckTest) TestIllegalConstShadowing_withES6Modules() {
	c.AssertRedeclareError(
		`export function f() {
  const x = 3;
  if (a) {
    var x;
  }
}
`)
}

// Java: VariableReferenceCheckTest.testVarShadowing()
func (c *VariableReferenceCheckTest) TestVarShadowing() {
	c.AssertRedeclareGlobal("if (a) { var x; var x;}")
	c.AssertRedeclareError("if (a) { var x; let x;}")
	c.AssertRedeclare("function f() { var x; if (a) { var x; }}")
	c.AssertRedeclareError("function f() { if (a) { var x; } let x;}")
	c.AssertNoWarning("function f() { var x; if (a) { let x; }}")
	c.AssertNoWarning(
		`function f() {
  if (a) { var x; }
  if (b) { let x; }
}
`)
}

// Java: VariableReferenceCheckTest.testVarShadowing_withES6Modules01()
func (c *VariableReferenceCheckTest) TestVarShadowing_withES6Modules01() {
	c.AssertRedeclare("export function f() { var x; if (a) { var x; }}")
}

// Java: VariableReferenceCheckTest.testVarShadowing_withES6Modules02()
func (c *VariableReferenceCheckTest) TestVarShadowing_withES6Modules02() {
	c.AssertRedeclareError("export function f() { if (a) { var x; } let x;}")
}

// Java: VariableReferenceCheckTest.testVarShadowing_withES6Modules03()
func (c *VariableReferenceCheckTest) TestVarShadowing_withES6Modules03() {
	c.AssertNoWarning("export function f() { var x; if (a) { let x; }}")
}

// Java: VariableReferenceCheckTest.testVarShadowing_withES6Modules04()
func (c *VariableReferenceCheckTest) TestVarShadowing_withES6Modules04() {
	c.AssertNoWarning(
		`function f() {
  if (a) { var x; }
  if (b) { let x; }
}
`)
}

// Java: VariableReferenceCheckTest.testParameterShadowing()
func (c *VariableReferenceCheckTest) TestParameterShadowing() {
	c.AssertRedeclareError("function f(x) { let x; }")
	c.AssertRedeclareError("function f(x) { const x = 3; }")
	c.AssertRedeclareError("function f(X) { class X {} }")
	c.AssertRedeclare("function f(x) { function x() {} }")
	c.AssertRedeclare("function f(x) { var x; }")
	c.AssertRedeclare("function f(x=3) { var x; }")
	c.AssertNoWarning("function f(...x) {}")
	c.AssertRedeclare("function f(...x) { var x; }")
	c.AssertRedeclare("function f(...x) { function x() {} }")
	c.AssertRedeclare("function f(x=3) { function x() {} }")
	c.AssertNoWarning("function f(x) { if (true) { let x; } }")
	c.AssertNoWarning(
		`function outer(x) {
  function inner() {
    let x = 1;
  }
}
`)
	c.AssertNoWarning(
		`function outer(x) {
  function inner() {
    var x = 1;
  }
}
`)
	c.AssertRedeclare("function f({a, b}) { var a = 2 }")
	c.AssertRedeclare("function f({a, b}) { if (!a) var a = 6; }")
}

// Java: VariableReferenceCheckTest.testParameterShadowing_withES6Modules()
func (c *VariableReferenceCheckTest) TestParameterShadowing_withES6Modules() {
	c.AssertRedeclareError("export function f(x) { let x; }")
	c.AssertRedeclare("export function f(x) { function x() {} }")
	c.AssertRedeclare("export function f(x=3) { var x; }")
	c.AssertNoWarning("export function f(...x) {}")
	c.AssertNoWarning(
		`export function outer(x) {
  function inner() {
    var x = 1;
  }
}
`)
}

// Java: VariableReferenceCheckTest.testReassignedConst()
func (c *VariableReferenceCheckTest) TestReassignedConst() {
	c.AssertReassign("const a = 0; a = 1;")
	c.AssertReassign("const a = 0; a++;")
}

// Java: VariableReferenceCheckTest.testLogicalReassignedConst()
func (c *VariableReferenceCheckTest) TestLogicalReassignedConst() {
	// These patterns are normalized away
	c.AssertReassign("const a = 0; a ||= 1;")
	c.AssertReassign("const a = 1; a &&= 1;")
	c.AssertReassign("const a = null; a ??= 1;")
}

// Java: VariableReferenceCheckTest.testLetConstNotDirectlyInBlock()
func (c *VariableReferenceCheckTest) TestLetConstNotDirectlyInBlock() {
	c.TestSame("if (true) var x = 3;")
	c.TestError("if (true) let x = 3;", VariableReferenceCheck_DECLARATION_NOT_DIRECTLY_IN_BLOCK)
	c.TestError("if (true) const x = 3;", VariableReferenceCheck_DECLARATION_NOT_DIRECTLY_IN_BLOCK)
	c.TestError("if (true) class C {}", VariableReferenceCheck_DECLARATION_NOT_DIRECTLY_IN_BLOCK)
	c.TestError("if (true) function f() {}", VariableReferenceCheck_DECLARATION_NOT_DIRECTLY_IN_BLOCK)
}

// Java: VariableReferenceCheckTest.testFunctionHoisting()
func (c *VariableReferenceCheckTest) TestFunctionHoisting() {
	c.AssertEarlyReferenceWarning("if (true) { f(); function f() {} }")
}

// Java: VariableReferenceCheckTest.testFunctionHoistingRedeclaration1()
func (c *VariableReferenceCheckTest) TestFunctionHoistingRedeclaration1() {
	js := []string{
		"var x;", "function x() {}",
	}
	message := "Variable x declared more than once. First occurrence: testcode0:1:4"
	c.TestError(c.Srcs(js), c.Error(VarCheck_VAR_MULTIPLY_DECLARED_ERROR).WithMessage(message))
}

// Java: VariableReferenceCheckTest.testFunctionHoistingRedeclaration2()
func (c *VariableReferenceCheckTest) TestFunctionHoistingRedeclaration2() {
	js := []string{
		"function x() {}", "var x;",
	}
	message := "Variable x declared more than once. First occurrence: testcode0:1:9"
	c.TestError(c.Srcs(js), c.Error(VarCheck_VAR_MULTIPLY_DECLARED_ERROR).WithMessage(message))
}

// Java: VariableReferenceCheckTest.testArrowFunction()
func (c *VariableReferenceCheckTest) TestArrowFunction() {
	c.AssertNoWarning("var f = x => { return x+1; };")
	c.AssertNoWarning("var odds = [1,2,3,4].filter((n) => n%2 == 1)")
	c.AssertRedeclare("var f = x => {var x;}")
	c.AssertRedeclareError("var f = x => {let x;}")
}

// Java: VariableReferenceCheckTest.testArrowFunction_withES6Modules()
func (c *VariableReferenceCheckTest) TestArrowFunction_withES6Modules() {
	c.AssertNoWarning("export var f = x => { return x+1; };")
	c.AssertRedeclare("export var f = x => {var x;}")
	c.AssertRedeclareError("export var f = x => {let x;}")
}

// Java: VariableReferenceCheckTest.testTryCatch()
func (c *VariableReferenceCheckTest) TestTryCatch() {
	c.AssertRedeclareError(
		`function f() {
  try {
    let e = 0;
    if (true) {
      let e = 1;
    }
  } catch (e) {
    let e;
  }
}
`)
	c.AssertRedeclareError(
		`function f() {
  try {
    let e = 0;
    if (true) {
      let e = 1;
    }
  } catch (e) {
      var e;
  }
}
`)
	c.AssertRedeclareError(
		`function f() {
  try {
    let e = 0;
    if (true) {
      let e = 1;
    }
  } catch (e) {
    function e() {
      var e;
    }
  }
}
`)
}

// Java: VariableReferenceCheckTest.testTryCatch_withES6Modules()
func (c *VariableReferenceCheckTest) TestTryCatch_withES6Modules() {
	c.AssertRedeclareError(
		`export function f() {
  try {
    let e = 0;
    if (true) {
      let e = 1;
    }
  } catch (e) {
    let e;
  }
}
`)
}

// Java: VariableReferenceCheckTest.testClass()
func (c *VariableReferenceCheckTest) TestClass() {
	c.AssertNoWarning("class A { f() { return 1729; } }")
}

// Java: VariableReferenceCheckTest.testClass_withES6Modules()
func (c *VariableReferenceCheckTest) TestClass_withES6Modules() {
	c.AssertNoWarning("export class A { f() { return 1729; } }")
}

// Java: VariableReferenceCheckTest.testRedeclareClassName()
func (c *VariableReferenceCheckTest) TestRedeclareClassName() {
	c.AssertNoWarning("var Clazz = class Foo {}; var Foo = 3;")
}

// Java: VariableReferenceCheckTest.testRedeclareClassName_withES6Modules()
func (c *VariableReferenceCheckTest) TestRedeclareClassName_withES6Modules() {
	c.AssertNoWarning("export var Clazz = class Foo {}; var Foo = 3;")
}

// Java: VariableReferenceCheckTest.testClassExtend()
func (c *VariableReferenceCheckTest) TestClassExtend() {
	c.AssertNoWarning("class A {} class C extends A {} C = class extends A {}")
}

// Java: VariableReferenceCheckTest.testClassExtend_withES6Modules()
func (c *VariableReferenceCheckTest) TestClassExtend_withES6Modules() {
	c.AssertNoWarning("export class A {} class C extends A {} C = class extends A {}")
}

// Variable reference before declaration error should not appear for non-static public fields
// Java: VariableReferenceCheckTest.testNonStaticPublicFields()
func (c *VariableReferenceCheckTest) TestNonStaticPublicFields() {
	c.AssertNoWarning("class Foo { x = bar;} let bar = 1;")
	c.AssertNoWarning("class Foo { x = Enum.X; } const Enum = { X: 1 }")
	c.AssertNoWarning("class Foo { x = new Bar(); } class Bar {}")
}

// Java: VariableReferenceCheckTest.testStaticPublicFields()
func (c *VariableReferenceCheckTest) TestStaticPublicFields() {
	c.AssertEarlyReferenceError("class Bar { static x = y; } const y = 3;")
	c.AssertEarlyReferenceError("class Foo { static x = new Bar(); } class Bar {}")
	c.AssertEarlyReferenceError("class Bar { static x = Enum.A; } let Enum = { A: 'str' }")
}

// Java: VariableReferenceCheckTest.testArrayPattern()
func (c *VariableReferenceCheckTest) TestArrayPattern() {
	c.AssertNoWarning("var [a] = [1];")
	c.AssertNoWarning("var [a, b] = [1, 2];")
	c.AssertEarlyReferenceWarning("alert(a); var [a] = [1];")
	c.AssertEarlyReferenceWarning("alert(b); var [a, b] = [1, 2];")
	c.AssertEarlyReferenceWarning("[a] = [1]; var a;")
	c.AssertEarlyReferenceWarning("[a, b] = [1]; var b;")
}

// Java: VariableReferenceCheckTest.testArrayPattern_withES6Modules01()
func (c *VariableReferenceCheckTest) TestArrayPattern_withES6Modules01() {
	c.AssertNoWarning("export var [a] = [1];")
}

// Java: VariableReferenceCheckTest.testArrayPattern_defaultValue()
func (c *VariableReferenceCheckTest) TestArrayPattern_defaultValue() {
	c.AssertNoWarning("var [a = 1] = [2];")
	c.AssertNoWarning("var [a = 1] = [];")
	c.AssertEarlyReferenceWarning("alert(a); var [a = 1] = [2];")
	c.AssertEarlyReferenceWarning("alert(a); var [a = 1] = [];")
	c.AssertEarlyReferenceWarning("alert(a); var [a = b] = [1];")
	c.AssertEarlyReferenceWarning("alert(a); var [a = b] = [];")
}

// Java: VariableReferenceCheckTest.testArrayPattern_defaultValue_withES6Modules01()
func (c *VariableReferenceCheckTest) TestArrayPattern_defaultValue_withES6Modules01() {
	c.AssertNoWarning("export var [a = 1] = [2];")
}

// Java: VariableReferenceCheckTest.testObjectPattern()
func (c *VariableReferenceCheckTest) TestObjectPattern() {
	c.AssertNoWarning("var {a: b} = {a: 1};")
	c.AssertNoWarning("var {a: b} = {};")
	c.AssertNoWarning("var {a} = {a: 1};")
	// 'a' is not declared at all, so the 'a' passed to alert() references
	// the global variable 'a', and there is no warning.
	c.AssertNoWarning("alert(a); var {a: b} = {};")
	c.AssertNoWarning("alert(a); var {a: {a: b}} = {};")
	c.AssertEarlyReferenceWarning("alert(b); var {a: b} = {a: 1};")
	c.AssertEarlyReferenceWarning("alert(a); var {a} = {a: 1};")
	c.AssertEarlyReferenceWarning("({a: b} = {}); var a, b;")
}

// Java: VariableReferenceCheckTest.testObjectPatternRest()
func (c *VariableReferenceCheckTest) TestObjectPatternRest() {
	c.AssertNoWarning("var {a: b, ...r} = {a: 1};")
	c.AssertNoWarning("var {a: b, ...r} = {};")
	c.AssertNoWarning("var {a, ...r} = {a: 1};")
	c.AssertNoWarning("alert(r);")
	c.AssertEarlyReferenceWarning("alert(r); var {...r} = {a: 1};")
	c.AssertNoWarning("({...a} = {});")
	c.AssertEarlyReferenceWarning("({...a} = {}); var a;")
}

// Java: VariableReferenceCheckTest.testObjectPattern_withES6Modules01()
func (c *VariableReferenceCheckTest) TestObjectPattern_withES6Modules01() {
	c.AssertNoWarning("export var {a: b} = {a: 1};")
}

// Java: VariableReferenceCheckTest.testObjectPattern_defaultValue()
func (c *VariableReferenceCheckTest) TestObjectPattern_defaultValue() {
	c.AssertEarlyReferenceWarning("alert(b); var {a: b = c} = {a: 1};")
	c.AssertEarlyReferenceWarning("alert(b); var c; var {a: b = c} = {a: 1};")
	c.AssertEarlyReferenceWarning("var {a: b = c} = {a: 1}; var c;")
	c.AssertEarlyReferenceWarning("alert(b); var {a: b = c} = {};")
	c.AssertEarlyReferenceWarning("alert(a); var {a = c} = {a: 1};")
	c.AssertEarlyReferenceWarning("alert(a); var {a = c} = {};")
}

// Java: VariableReferenceCheckTest.testObjectPattern_defaultValue_withES6Modules()
func (c *VariableReferenceCheckTest) TestObjectPattern_defaultValue_withES6Modules() {
	c.AssertEarlyReferenceWarning("export var {a: b = c} = {a: 1}; var c;")
}

// We can't catch all possible runtime errors but it's useful to have some basic checks.
// Java: VariableReferenceCheckTest.testDefaultParam()
func (c *VariableReferenceCheckTest) TestDefaultParam() {
	c.AssertEarlyReferenceError("function f(x=a) { let a; }")
	c.AssertEarlyReferenceError(
		`function f(x=a) { let a; }
function g(x=1) { var a; }
`)
	c.AssertEarlyReferenceError("function f(x=a) { var a; }")
	c.AssertEarlyReferenceError("function f(x=a()) { function a() {} }")
	c.AssertEarlyReferenceError("function f(x=[a]) { var a; }")
	c.AssertEarlyReferenceError("function f(x={a}) { let a; }")
	c.AssertEarlyReferenceError("function f(x=y, y=2) {}")
	c.AssertEarlyReferenceError("function f(x={y}, y=2) {}")
	c.AssertEarlyReferenceError("function f(x=x) {}")
	c.AssertEarlyReferenceError("function f([x]=x) {}")
	// x within a function isn't referenced at the time the default value for x is evaluated.
	c.AssertNoWarning("function f(x=()=>x) {}")
	c.AssertNoWarning("function f(x=a) {}")
	c.AssertNoWarning("function f(x=a) {} var a;")
	c.AssertNoWarning("let b; function f(x=b) { var b; }")
	c.AssertNoWarning("function f(y = () => x, x = 5) { return y(); }")
	c.AssertNoWarning("function f(x = new foo.bar()) {}")
	c.AssertNoWarning("var foo = {}; foo.bar = class {}; function f(x = new foo.bar()) {}")
}

// Java: VariableReferenceCheckTest.testDefaultParam_withES6Modules()
func (c *VariableReferenceCheckTest) TestDefaultParam_withES6Modules() {
	c.AssertEarlyReferenceError("export function f(x=a) { let a; }")
	c.AssertNoWarning("export function f(x=()=>x) {}")
}

// Java: VariableReferenceCheckTest.testDestructuring()
func (c *VariableReferenceCheckTest) TestDestructuring() {
	c.TestSame(
		`function f() {
  var obj = {a:1, b:2};
  var {a:c, b:d} = obj;
}
`)
	c.TestSame(
		`function f() {
  var obj = {a:1, b:2};
  var {a, b} = obj;
}
`)
	c.AssertRedeclare(
		`function f() {
  var obj = {a:1, b:2};
  var {a:c, b:d} = obj;
  var c = b;
}
`)
	c.AssertEarlyReferenceWarning(
		`function f() {
  var {a:c, b:d} = obj;
  var obj = {a:1, b:2};
}
`)
	c.AssertEarlyReferenceWarning(
		`function f() {
  var {a, b} = obj;
  var obj = {a:1, b:2};
}
`)
	c.AssertEarlyReferenceWarning(
		`function f() {
  var e = c;
  var {a:c, b:d} = {a:1, b:2};
}
`)
}

// Java: VariableReferenceCheckTest.testDestructuring_withES6Modules()
func (c *VariableReferenceCheckTest) TestDestructuring_withES6Modules() {
	c.TestSame(
		`export function f() {
  var obj = {a:1, b:2};
  var {a:c, b:d} = obj;
}
`)
	c.AssertRedeclare(
		`export function f() {
  var obj = {a:1, b:2};
  var {a:c, b:d} = obj;
  var c = b;
}
`)
	c.AssertEarlyReferenceWarning(
		`export function f() {
  var {a:c, b:d} = obj;
  var obj = {a:1, b:2};
}
`)
}

// Java: VariableReferenceCheckTest.testDestructuringInLoop()
func (c *VariableReferenceCheckTest) TestDestructuringInLoop() {
	c.TestSame("for (let {length: x} in obj) {}")
	c.TestSame("for (let [{length: z}, w] in obj) {}")
}

// Java: VariableReferenceCheckTest.testReferencingPreviouslyDeclaredVariableInConst()
func (c *VariableReferenceCheckTest) TestReferencingPreviouslyDeclaredVariableInConst() {
	c.TestSame("const [a, b = a] = [];")
	// TODO(b/111441110): don't error on this. it's valid code.
	c.AssertEarlyReferenceError("for (const [a, b = a] of []);")
}

// Java: VariableReferenceCheckTest.testEarlyReferenceInInnerBlock()
func (c *VariableReferenceCheckTest) TestEarlyReferenceInInnerBlock() {
	c.AssertEarlyReferenceError("for (x of [1, 2, 3]) {} let x;")
	c.AssertEarlyReferenceError("{ x; } let x;")
	c.AssertEarlyReferenceError("{ C; } class C {}")
	c.AssertEarlyReferenceWarning("{ x; }  var x;")
}

// Java: VariableReferenceCheckTest.testEarlyVariableReferenceInsideFunction()
func (c *VariableReferenceCheckTest) TestEarlyVariableReferenceInsideFunction() {
	c.TestSame("function f() { x; } let x; f(); ")
	c.TestSame("function f() { const f = () => x; let x = 3; return f; }")
	// NOTE: this will cause an error at runtime, but we don't report it because we don't track
	// where `f` is being called.
	c.TestSame("function f() { x; } f(); let x;")
	c.TestSame("function f() { x; } f(); var x;")
}

// Java: VariableReferenceCheckTest.testEnhancedForLoopTemporalDeadZone()
func (c *VariableReferenceCheckTest) TestEnhancedForLoopTemporalDeadZone() {
	c.AssertEarlyReferenceError("for (let x of [x]);")
	c.AssertEarlyReferenceError("for (let x in [x]);")
	c.AssertEarlyReferenceError("for (const x of [x]);")
	c.TestSame("for (var x of [x]);")
	c.TestSame("for (let x of [() => x]);")
	c.TestSame("let x = 1; for (let y of [x]);")
}

// Java: VariableReferenceCheckTest.testEnhancedForLoopTemporalDeadZone_withES6Modules()
func (c *VariableReferenceCheckTest) TestEnhancedForLoopTemporalDeadZone_withES6Modules() {
	c.TestSame("export let x = 1; for (let y of [x]);")
}

// Java: VariableReferenceCheckTest.testRedeclareVariableFromImport()
func (c *VariableReferenceCheckTest) TestRedeclareVariableFromImport() {
	c.IgnoreWarnings(ModuleLoader_INVALID_MODULE_PATH)
	c.AssertRedeclareError("import {x} from 'whatever'; let x = 0;")
	c.AssertRedeclareError("import {x} from 'whatever'; const x = 0;")
	c.AssertRedeclareError("import {x} from 'whatever'; var x = 0;")
	c.AssertRedeclareError("import {x} from 'whatever'; function x() {}")
	c.AssertRedeclareError("import {x} from 'whatever'; class x {}")
	c.AssertRedeclareError("import x from 'whatever'; let x = 0;")
	c.AssertRedeclareError("import * as ns from 'whatever'; let ns = 0;")
	c.AssertRedeclareError("import {y as x} from 'whatever'; let x = 0;")
	c.AssertRedeclareError("import {x} from 'whatever'; let {x} = {};")
	c.AssertRedeclareError("import {x} from 'whatever'; let [x] = [];")
	c.AssertRedeclareError("import {x, x} from 'whatever';")
	c.AssertRedeclareError("import {x, y as x} from 'whatever';")
	c.AssertRedeclareError("import {z as x, y as x} from 'whatever';")
	c.AssertRedeclareError("import {x} from 'first'; import {x} from 'second';")
	c.AssertRedeclareError("import {x} from 'first'; import {a as x} from 'second';")
	c.AssertRedeclareError("import {b as x} from 'first'; import {a as x} from 'second';")
	c.TestSame("import {x} from 'whatever'; function f() { let x = 0; }")
	c.TestSame("import {x as x} from 'whatever'; function f() { let x = 0; }")
	c.TestSame("import {y as x} from 'whatever'; function f() { let x = 0; }")
}

// Java: VariableReferenceCheckTest.testOkExportsRefInGoogModule()
func (c *VariableReferenceCheckTest) TestOkExportsRefInGoogModule() {
	c.TestSame("goog.module('m');")
	c.TestSame("goog.module('m'); exports.Foo = 0; exports.Bar = 0;")
	c.TestSame("goog.module('m'); exports = 0;")
	c.TestSame("goog.module('m'); exports = class {}; exports.Foo = class {};")
	c.TestSame("goog.module('m'); function f() { exports = 0; }")
	// Bad style but warn elsewhere
	c.TestSame("goog.module('m'); function f() { return exports; } exports = 1;")
}

// Java: VariableReferenceCheckTest.testBadEarlyExportsRefInGoogModule()
func (c *VariableReferenceCheckTest) TestBadEarlyExportsRefInGoogModule() {
	c.TestError("goog.module('m'); exports.x = 0; exports = {};", VariableReferenceCheck_EARLY_EXPORTS_REFERENCE)
	c.TestError("goog.module('m'); exports.x = 0; exports = class Bar {};", VariableReferenceCheck_EARLY_EXPORTS_REFERENCE)
	c.TestError(
		"goog.module('m'); /** @typedef {string} */ exports.x; exports = {};",
		VariableReferenceCheck_EARLY_EXPORTS_REFERENCE)
}

// Java: VariableReferenceCheckTest.testReferenceInSwitchCondition_shadowedInSwitchBody()
func (c *VariableReferenceCheckTest) TestReferenceInSwitchCondition_shadowedInSwitchBody() {
	c.AssertNoWarning(
		`const x = 0;
switch (x) {
  case 0:
  let x;
  break;
  default: break;
}
`)
}

// Expects the JS to generate one bad-read error.
// Java: VariableReferenceCheckTest.assertRedeclare(String)
func (c *VariableReferenceCheckTest) AssertRedeclare(js string) {
	c.TestWarning(js, VariableReferenceCheck_REDECLARED_VARIABLE)
}

// Java: VariableReferenceCheckTest.assertRedeclareError(String)
func (c *VariableReferenceCheckTest) AssertRedeclareError(js string) {
	c.TestError(js, VariableReferenceCheck_REDECLARED_VARIABLE_ERROR)
}

// Java: VariableReferenceCheckTest.assertReassign(String)
func (c *VariableReferenceCheckTest) AssertReassign(js string) {
	c.TestError(js, VariableReferenceCheck_REASSIGNED_CONSTANT)
}

// Java: VariableReferenceCheckTest.assertRedeclareGlobal(String)
func (c *VariableReferenceCheckTest) AssertRedeclareGlobal(js string) {
	c.TestError(js, VarCheck_VAR_MULTIPLY_DECLARED_ERROR)
}

// Expects the JS to generate one bad-write warning.
// Java: VariableReferenceCheckTest.assertEarlyReferenceWarning(String)
func (c *VariableReferenceCheckTest) AssertEarlyReferenceWarning(js string) {
	c.TestWarning(js, VariableReferenceCheck_EARLY_REFERENCE)
}

// Java: VariableReferenceCheckTest.assertEarlyReferenceError(String)
func (c *VariableReferenceCheckTest) AssertEarlyReferenceError(js string) {
	c.TestError(js, VariableReferenceCheck_EARLY_REFERENCE_ERROR)
}

// Expects the JS to generate no errors or warnings.
// Java: VariableReferenceCheckTest.assertNoWarning(String)
func (c *VariableReferenceCheckTest) AssertNoWarning(js string) {
	c.TestSame(js)
}
