// Port of com/google/javascript/jscomp/CheckUnreachableCodeTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/CheckUnreachableCodeTest.java)

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

// CheckUnreachableCodeTest: Tests for CheckUnreachableCode.
type CheckUnreachableCodeTest struct {
	CompilerTestCaseBase
}

// Java: CheckUnreachableCodeTest()
func NewCheckUnreachableCodeTest() *CheckUnreachableCodeTest {
	c := &CheckUnreachableCodeTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestCheckUnreachableCodeTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCheckUnreachableCodeTest() })
}

// Java: CheckUnreachableCodeTest.getProcessor(Compiler)
func (c *CheckUnreachableCodeTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewCombinedCompilerPass(compiler, NewCheckUnreachableCode(compiler))
}

// Java: CheckUnreachableCodeTest.testCorrectSimple()
func (c *CheckUnreachableCodeTest) TestCorrectSimple() {
	c.TestSame("var x")
	c.TestSame("var x = 1")
	c.TestSame("var x = 1; x = 2;")
	c.TestSame("if (x) { var x = 1 }")
	c.TestSame("if (x) { var x = 1 } else { var y = 2 }")
	c.TestSame("while(x) {}")
}

// Java: CheckUnreachableCodeTest.testIncorrectSimple()
func (c *CheckUnreachableCodeTest) TestIncorrectSimple() {
	c.AssertUnreachable("function f() { return; x=1; }")
	c.AssertUnreachable("function f() { return; x=1; x=1; }")
	c.AssertUnreachable("function f() { return; var x = 1; }")
}

// Java: CheckUnreachableCodeTest.testCorrectIfReturns()
func (c *CheckUnreachableCodeTest) TestCorrectIfReturns() {
	c.TestSame("function f() { if (x) { return } }")
	c.TestSame("function f() { if (x) { return } return }")
	c.TestSame("function f() { if (x) { if (y) { return } } else { return }}")
	c.TestSame("function f() { if (x) { if (y) { return } return } else { return }}")
}

// Java: CheckUnreachableCodeTest.testInCorrectIfReturns()
func (c *CheckUnreachableCodeTest) TestInCorrectIfReturns() {
	c.AssertUnreachable("function f() { if (x) { return } else { return } return }")
}

// Java: CheckUnreachableCodeTest.testCorrectSwitchReturn()
func (c *CheckUnreachableCodeTest) TestCorrectSwitchReturn() {
	c.TestSame("function f() { switch(x) { default: return; case 1: x++; }}")
	c.TestSame("function f() { switch(x) { default: return; case 1: x++; } return }")
	c.TestSame("function f() { switch(x) { default: return; case 1: return; }}")
	c.TestSame("function f() { switch(x) { case 1: return; } return }")
	c.TestSame("function f() { switch(x) { case 1: case 2: return; } return }")
	c.TestSame("function f() { switch(x) { case 1: return; case 2: return; } return }")
	c.TestSame("function f() { switch(x) { case 1 : return; case 2: return; } return }")
}

// Java: CheckUnreachableCodeTest.testInCorrectSwitchReturn()
func (c *CheckUnreachableCodeTest) TestInCorrectSwitchReturn() {
	c.AssertUnreachable("function f() { switch(x) { default: return; case 1: return; } return }")
	c.AssertUnreachable("function f() { switch(x) { default: return; return; case 1: return; } }")
}

// Java: CheckUnreachableCodeTest.testCorrectLoopBreaksAndContinues()
func (c *CheckUnreachableCodeTest) TestCorrectLoopBreaksAndContinues() {
	c.TestSame("while(1) { foo(); break }")
	c.TestSame("while(1) { foo(); continue }")
	c.TestSame("for(;;) { foo(); break }")
	c.TestSame("for(;;) { foo(); continue }")
	c.TestSame("for(;;) { if (x) { break } }")
	c.TestSame("for(;;) { if (x) { continue } }")
	c.TestSame("do { foo(); continue} while(1)")
}

// Java: CheckUnreachableCodeTest.testInCorrectLoopBreaksAndContinues()
func (c *CheckUnreachableCodeTest) TestInCorrectLoopBreaksAndContinues() {
	c.AssertUnreachable("while(1) { foo(); break; bar()}")
	c.AssertUnreachable("while(1) { foo(); continue; bar() }")
	c.AssertUnreachable("for(;;) { foo(); break; bar() }")
	c.AssertUnreachable("for(;;) { foo(); continue; bar() }")
	c.AssertUnreachable("for(;;) { if (x) { break; bar() } }")
	c.AssertUnreachable("for(;;) { if (x) { continue; bar() } }")
	c.AssertUnreachable("do { foo(); continue; bar()} while(1)")
}

// Java: CheckUnreachableCodeTest.testUncheckedWhileInDo()
func (c *CheckUnreachableCodeTest) TestUncheckedWhileInDo() {
	c.AssertUnreachable("do { foo(); break} while(1)")
}

// Java: CheckUnreachableCodeTest.testUncheckedConditionInFor()
func (c *CheckUnreachableCodeTest) TestUncheckedConditionInFor() {
	c.AssertUnreachable("for(var x = 0; x < 100; x++) { break };")
}

// Java: CheckUnreachableCodeTest.testFunctionDeclaration()
func (c *CheckUnreachableCodeTest) TestFunctionDeclaration() {
	// functions are not in our CFG.
	c.TestSame("function f() { return; function ff() { }}")
}

// Java: CheckUnreachableCodeTest.testVarDeclaration()
func (c *CheckUnreachableCodeTest) TestVarDeclaration() {
	c.AssertUnreachable("function f() { return; var x = 1 }")
	// I think the user should fix this as well.
	c.AssertUnreachable("function f() { return; var x }")
}

// Java: CheckUnreachableCodeTest.testReachableTryCatchFinally()
func (c *CheckUnreachableCodeTest) TestReachableTryCatchFinally() {
	c.TestSame("try { } finally {  }")
	c.TestSame("try { foo(); } finally { bar() } ")
	c.TestSame("try { foo() } finally { bar() }")
	c.TestSame("try { foo(); } catch (e) {e()} finally { bar() }")
	c.TestSame("try { foo() } catch (e) {e()} finally { bar() }")
	c.TestSame("try { foo() } catch (e) { throw e; } finally { bar() }")
}

// Java: CheckUnreachableCodeTest.testUnreachableCatch()
func (c *CheckUnreachableCodeTest) TestUnreachableCatch() {
	c.AssertUnreachable("try { var x = 0 } catch (e) { }")
	c.AssertUnreachable("try { } catch (e) { throw e; }")
}

// Java: CheckUnreachableCodeTest.testSpuriousBreak()
func (c *CheckUnreachableCodeTest) TestSpuriousBreak() {
	c.TestSame("switch (x) { default: throw x; break; }")
}

// Java: CheckUnreachableCodeTest.testInstanceOfThrowsException()
func (c *CheckUnreachableCodeTest) TestInstanceOfThrowsException() {
	c.TestSame("function f() {try { if (value instanceof type) return true; } " + "catch (e) { }}")
}

// Java: CheckUnreachableCodeTest.testFalseCondition()
func (c *CheckUnreachableCodeTest) TestFalseCondition() {
	c.AssertUnreachable("if(false) { }")
	c.AssertUnreachable("if(0) { }")
}

// Java: CheckUnreachableCodeTest.testUnreachableLoop()
func (c *CheckUnreachableCodeTest) TestUnreachableLoop() {
	c.AssertUnreachable("while(false) {}")
}

// Java: CheckUnreachableCodeTest.testInfiniteLoop()
func (c *CheckUnreachableCodeTest) TestInfiniteLoop() {
	c.TestSame("while (true) { foo(); break; }")

	// TODO(user): Have a infinite loop warning instead.
	c.AssertUnreachable("while(true) {} foo()")
}

// Java: CheckUnreachableCodeTest.testUnreachableStaticBlocks()
func (c *CheckUnreachableCodeTest) TestUnreachableStaticBlocks() {
	c.AssertUnreachable("class C{ static { function f(){ return; let x; } } }")
	c.AssertUnreachable("class C{ static{ if(false) { var x; } } }")
	c.AssertUnreachable("class C{ static{ while(false){ let x; } } }")
}

// Java: CheckUnreachableCodeTest.testReachableStaticBlocks()
func (c *CheckUnreachableCodeTest) TestReachableStaticBlocks() {
	c.TestSame("class C{ static { let x; } }")
	c.TestSame("class C{ static { function a(){} } }")
	c.TestSame("class C{ static {} a(){} }")
	c.TestSame("class C{ static {} static {} }")
	c.TestSame("class C{ static{} static x; static{} static y;}")
	c.TestSame("class C{ static [l('l1')] = l('r1'); static{ l('s2');} static[l('l3')] = l('r3');}")
}

// Java: CheckUnreachableCodeTest.testSuppression()
func (c *CheckUnreachableCodeTest) TestSuppression() {
	c.AssertUnreachable("if(false) { }")

	c.TestSame(
		`/** @fileoverview
 * @suppress {uselessCode}
 */
if(false) { }
`)

	c.TestSame(
		`/** @fileoverview
 * @suppress {uselessCode}
 */
function f() { if(false) { } }
`)

	c.TestSame(
		`/**
 * @suppress {uselessCode}
 */
function f() { if(false) { } }
`)

	c.AssertUnreachable(
		`/**
 * @suppress {uselessCode}
 */
function f() { if(false) { } }
function g() { if(false) { } }
`)

	c.TestSame(
		`/**
 * @suppress {uselessCode}
 */
function f() {
  function g() { if(false) { } }
  if(false) { } }
`)

	c.AssertUnreachable(
		`function f() {
  /**
   * @suppress {uselessCode}
   */
  function g() { if(false) { } }
  if(false) { } }
`)

	c.TestSame(
		`function f() {
  /**
   * @suppress {uselessCode}
   */
  function g() { if(false) { } }
}
`)
}

// Java: CheckUnreachableCodeTest.testES6FeaturesInIfExpression()
func (c *CheckUnreachableCodeTest) TestES6FeaturesInIfExpression() {
	// class X{} always eval to true by toBoolean();
	c.AssertUnreachable("if(!class {}) x = 1;")
	c.AssertUnreachable("if(!class A{}) x = 1;")

	// Template string with substitution and tagged template will be evaluated
	// to UNKNOWN. Template string with definite length (without substitution)
	// will be evaled like a normal string.
	c.AssertUnreachable("if(!`tempLit`) {x = 1;}")
	c.AssertUnreachable("if(``) {x = 1;}")
	c.TestSame("if(`temp${sub}Lit`) {x = 1;} else {x = 2;}")
	c.TestSame("if(`${sub}`) {x = 1;} else {x = 2;}")
	c.TestSame("if(tagged`tempLit`) {x = 1;} else {x = 2;}")

	// Functions are truthy.
	c.AssertUnreachable("if(()=>true) {x = 1;} else {x = 2;}")
}

// Java: CheckUnreachableCodeTest.testES6FeaturesInTryCatch()
func (c *CheckUnreachableCodeTest) TestES6FeaturesInTryCatch() {
	c.AssertUnreachable("try { let x = 1; } catch(e) {}")
	c.AssertUnreachable("try { const x = 1; } catch(e) {}")
	c.AssertUnreachable("try {()=>42;} catch(e) {}")
	c.AssertUnreachable("try {function *gen(){};} catch(e) {}")
	// Assumed tagged template may throw exception.
	c.TestSame("try {tagged`temp`;} catch(e) {}")

	c.AssertUnreachable("try { var obj = {a(){}};} catch(e) {}")
	c.TestSame("try { var obj = {a(){}}; obj.a();} catch(e) {}")
}

// Java: CheckUnreachableCodeTest.testCorrectForOfBreakAndContinues()
func (c *CheckUnreachableCodeTest) TestCorrectForOfBreakAndContinues() {
	c.TestSame("for(x of [1, 2, 3]) {foo();}")
	c.TestSame("for(x of [1, 2, 3]) {foo(); break;}")
	c.TestSame("for(x of [1, 2, 3]) {foo(); continue;}")
}

// Java: CheckUnreachableCodeTest.testInCorrectForOfBreakAndContinues()
func (c *CheckUnreachableCodeTest) TestInCorrectForOfBreakAndContinues() {
	c.AssertUnreachable("for(x of [1, 2, 3]) {foo(); break; bar();}")
	c.AssertUnreachable("for(x of [1, 2, 3]) {foo(); continue; bar();}")

	c.AssertUnreachable("for(x of [1, 2, 3]) {if(x) {break; bar();}}")
	c.AssertUnreachable("for(x of [1, 2, 3]) {if(x) {continue; bar();}}")
}

// Java: CheckUnreachableCodeTest.testForLoopsEs6()
func (c *CheckUnreachableCodeTest) TestForLoopsEs6() {
	c.AssertUnreachable("for(;;) {if(x) {continue; bar();}}")
	c.AssertUnreachable("for(x in y) {if(x) {continue; bar();}}")
}

// Java: CheckUnreachableCodeTest.testReturnsInShorthandFunctionOfObjLit()
func (c *CheckUnreachableCodeTest) TestReturnsInShorthandFunctionOfObjLit() {
	c.TestSame(
		`var obj = {
  f() {
    switch(x) {
      default: return;
      case 1: x++;
    }
  }
}
`)
	c.AssertUnreachable(
		`var obj = {f() {
  switch(x) {
    default: return;
    case 1: return;
  }
  return;
}}
`)
	c.TestSame("var obj = {f() { if(x) {return;} else {return; }}}")
	c.AssertUnreachable(
		`var obj = {f() {
  if(x) {
    return;
  } else {
    return;
  }
  return;
}}
`)
}

// Java: CheckUnreachableCodeTest.testObjLit()
func (c *CheckUnreachableCodeTest) TestObjLit() {
	c.AssertUnreachable("var a = {c(){if(true){return;}x = 1;}};")
}

// Java: CheckUnreachableCodeTest.testClass()
func (c *CheckUnreachableCodeTest) TestClass() {
	c.TestSame("class C{func(){}}")
	c.AssertUnreachable("class C{func(){if (true){return;} else {return;}}}")
	c.AssertUnreachable("class C{func(){if (true){return;} x = 1;}}")
	c.TestSame("var C = class{func(){}}")
	c.TestSame("let C = class{func(){}}")
	c.TestSame("var C; C = class{func(){}}")
	c.TestSame("let C; C = class{func(){}}")
	c.AssertUnreachable("var C = class{func(){if (true){return;} x = 1;}}")
	c.AssertUnreachable("class C {x = function(){if(true) {return;} x=1;}}")
}

// Java: CheckUnreachableCodeTest.testUnderClass()
func (c *CheckUnreachableCodeTest) TestUnderClass() {
	c.TestSame("class C {} alert(1);")
	c.TestSame("class D {} class C extends D {} alert(1)")
	c.TestSame("class D{} alert(1); class C extends D {}")
}

// Java: CheckUnreachableCodeTest.testFunction()
func (c *CheckUnreachableCodeTest) TestFunction() {
	c.TestSame("function f() {} alert(1);")
}

// Java: CheckUnreachableCodeTest.testSubclass()
func (c *CheckUnreachableCodeTest) TestSubclass() {
	c.TestSame(
		`class D {foo() {if (true) {return;}}}
class C extends D {foo() {super.foo();}}
`)
}

// Java: CheckUnreachableCodeTest.testArrowFunction()
func (c *CheckUnreachableCodeTest) TestArrowFunction() {
	c.TestSame("() => 3")
	c.TestSame("e => e + 1")
	c.TestSame("listen('click', e => onclick(e), true);")
	c.TestSame("listen('click', e => { onclick(e); }, true);")
	c.AssertUnreachable("listen('click', e => { return 0; onclick(e); }, true);")
	c.AssertUnreachable("() => {return 3; x = 1;}")
	c.AssertUnreachable("() => { if (false) x = 1;}")
	c.TestSame("var f = array.filter(g => {if (g % 3) x = 1;});")
	c.AssertUnreachable("var f = array.filter(g => {if (false) g = 1;});")
}

// Java: CheckUnreachableCodeTest.testGenerators()
func (c *CheckUnreachableCodeTest) TestGenerators() {
	c.TestSame(
		`function* f() {
  var i = 0;
  while(true)
    yield i++;
}
`)

	c.AssertUnreachable(
		`function* f() {
  var i = 0;
  while(true) {
    yield i++;
  }
  i = 1;
}
`)

	c.TestSame(
		`function* f() {
  var i = 0;
  while(true) {
    yield i;
    i++;
  }
}
`)

	c.TestSame(
		`function *f() {
  try {
    yield;
  } catch (e) {
    alert(e);
  }
}
`)

	c.TestSame(
		`function *f() {
  try {
    yield 1;
  } catch (e) {
    alert(e);
  }
}
`)
}

// Java: CheckUnreachableCodeTest.testAwaitCanCauseError()
func (c *CheckUnreachableCodeTest) TestAwaitCanCauseError() {
	c.TestSame("async function f(/** !Promise<?> */ p) { try { await p; } catch (e) {} }")
	// Note: for simplicity, we assume that `await <expr>` may always throw an exception, even if
	// <expr> is a literal.
	c.TestSame("async function f() { try { await 3; } catch (e) {} }")
}

// Java: CheckUnreachableCodeTest.testForAwait()
func (c *CheckUnreachableCodeTest) TestForAwait() {
	c.TestSame("async function f(iter) { for await (const item of iter) { item; } }")
	c.AssertUnreachable(
		`async function f(iter) {
  for await (const item of iter) { if (false) { 3; } }
}
`)
}

// Java: CheckUnreachableCodeTest.testCorrectEsModuleExports()
func (c *CheckUnreachableCodeTest) TestCorrectEsModuleExports() {
	c.IgnoreWarnings(DiagnosticGroups_MODULE_LOAD)
	c.TestSame("export var x")
	c.TestSame("export var x = 1")
	c.TestSame("export function f() {}")
	c.TestSame("export class C {};")
	c.TestSame("export default 0;")
	c.TestSame("const x = 0; export {x};")
	c.TestSame("export {x as y} from './mod'")
}

// Java: CheckUnreachableCodeTest.testCorrectEsModuleImports()
func (c *CheckUnreachableCodeTest) TestCorrectEsModuleImports() {
	c.IgnoreWarnings(DiagnosticGroups_MODULE_LOAD)
	c.TestSame("import * as mod from './mod';")
	c.TestSame("import def from './mod';")
	c.TestSame("import {x, y} from './mod';")
	c.TestSame("import {x as y} from './mod';")
}

// Java: CheckUnreachableCodeTest.assertUnreachable(String)
func (c *CheckUnreachableCodeTest) AssertUnreachable(js string) {
	c.Test(js, js, c.Warning(CheckUnreachableCode_UNREACHABLE_CODE))
}
