// Port of com/google/javascript/jscomp/CheckMissingReturnTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/CheckMissingReturnTest.java)

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
	"strings"
	"testing"
)

// CheckMissingReturnTest: Tests for CheckMissingReturn.
type CheckMissingReturnTest struct {
	CompilerTestCaseBase
}

// Java: CheckMissingReturnTest()
func NewCheckMissingReturnTest() *CheckMissingReturnTest {
	c := &CheckMissingReturnTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestCheckMissingReturnTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCheckMissingReturnTest() })
}

// Java: CheckMissingReturnTest.setUp()
func (c *CheckMissingReturnTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableTypeCheck()
}

// Java: CheckMissingReturnTest.getProcessor(Compiler)
func (c *CheckMissingReturnTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewCombinedCompilerPass(compiler, NewCheckMissingReturn(compiler))
}

// Java: CheckMissingReturnTest.testMissingReturn()
func (c *CheckMissingReturnTest) TestMissingReturn() {
	// Requires control flow analysis.
	c.TestMissing("if (a) { return 1; }")

	// Switch statement.
	c.TestMissing("switch(1) { case 12: return 5; }")

	// Test try catch finally.
	c.TestMissing("try { foo() } catch (e) { return 5; } finally { }")

	// Nested scope.
	c.TestMissing("/** @return {number} */ function f() { var x; }; return 1;")
	c.TestMissing("/** @return {number} */ function f() { return 1; };")
}

// Java: CheckMissingReturnTest.testReturnNotMissing()
func (c *CheckMissingReturnTest) TestReturnNotMissing() {
	// Empty function body. Ignore this case. The remainder of the functions in
	// this test have non-empty bodies.
	c.TestNotMissing("")

	// Simple cases.
	c.TestSame("function f() { var x; }")
	c.TestNotMissing("return 1;")

	// Returning void and undefined.
	c.TestNotMissingWithReturnType("void", "var x;")
	c.TestNotMissingWithReturnType("undefined", "var x;")

	// Returning a union that includes void or undefined.
	c.TestNotMissingWithReturnType("number|undefined", "var x;")
	c.TestNotMissingWithReturnType("number|void", "var x;")
	c.TestNotMissingWithReturnType("*", "var x;")

	// Test try catch finally.
	c.TestNotMissing("try { return foo() } catch (e) { } finally { }")
	c.TestNotMissing("try {throw e;} catch (e) { return foo() } finally { }")
	c.TestNotMissing("try {} catch (e) {} finally {return foo()};")

	// Nested function.
	c.TestNotMissing("/** @return {number} */ function f() { return 1; }; return 1;")

	// Strange tests that come up when reviewing closure code.
	c.TestNotMissing("try { return 12; } finally { return 62; }")
	c.TestNotMissing("try { } finally { return 1; }")
	c.TestNotMissing("switch(1) { default: return 1; }")
	c.TestNotMissing("switch(g) { case 1: return 1; default: return 2; }")
}

// Java: CheckMissingReturnTest.testFinallyStatements()
func (c *CheckMissingReturnTest) TestFinallyStatements() {
	// The control flow analysis (CFA) treats finally blocks somewhat strangely.
	// The CFA might indicate that a finally block implicitly returns. However,
	// if entry into the finally block is normally caused by an explicit return
	// statement, then a return statement isn't missing:
	//
	// try {
	//   return 1;
	// } finally {
	//   // CFA determines implicit return. However, return not missing
	//   // because of return statement in try block.
	// }
	//
	// Hence extra tests are warranted for various cases involving finally
	// blocks.

	// Simple finally case.
	c.TestNotMissing("try { return 1; } finally { }")
	c.TestNotMissing("try { } finally { return 1; }")
	c.TestMissing("try { } finally { }")

	// Cycles in the CFG within the finally block were causing problems before.
	c.TestNotMissing("try { return 1; } finally { while (true) { } }")
	c.TestMissing("try { } finally { while (x) { } }")
	c.TestMissing("try { } finally { while (x) { if (x) { break; } } }")
	c.TestNotMissing("try { return 2; } finally { while (x) { if (x) { break; } } }")

	// Test various cases with nested try statements.
	c.TestMissing("try { } finally { try { } finally { } }")
	c.TestNotMissing("try { } finally { try { return 1; } finally { } }")
	c.TestNotMissing("try { return 1; } finally { try { } finally { } }")

	// Calling a function potentially causes control flow to transfer to finally
	// block. However, the function will not return in this case as the
	// exception will unwind the stack. Hence this function isn't missing a
	// return statement (i.e., the running program will not expect a return
	// value from the function if an exception is thrown).
	c.TestNotMissing("try { g(); return 1; } finally { }")

	// Closures within try ... finally affect missing return statement analysis
	// because of the nested scopes. The following tests check for missing
	// return statements in the three possible configurations: both scopes
	// return; enclosed doesn't return; enclosing doesn't return.
	c.TestNotMissing(
		`try {
   /** @return {number} */ function f() {
       try { return 1; }
       finally { }
   };
   return 1;
}
finally { }
`)
	c.TestMissing(
		`try {
   /** @return {number} */ function f() {
       try { }
       finally { }
   };
   return 1;
}
finally { }
`)
	c.TestMissing(
		`try {
   /** @return {number} */ function f() {
       try { return 1; }
       finally { }
   };
}
finally { }
`)
}

// Java: CheckMissingReturnTest.testKnownConditions()
func (c *CheckMissingReturnTest) TestKnownConditions() {
	c.TestNotMissing("if (true) return 1")
	c.TestMissing("if (true) {} else {return 1}")

	c.TestMissing("if (false) return 1")
	c.TestNotMissing("if (false) {} else {return 1}")

	c.TestNotMissing("if (1) return 1")
	c.TestMissing("if (1) {} else {return 1}")

	c.TestMissing("if (0) return 1")
	c.TestNotMissing("if (0) {} else {return 1}")

	c.TestNotMissing("if (3) return 1")
	c.TestMissing("if (3) {} else {return 1}")
}

// Java: CheckMissingReturnTest.testKnownWhileLoop()
func (c *CheckMissingReturnTest) TestKnownWhileLoop() {
	c.TestNotMissing("while (1) return 1")
	c.TestNotMissing("while (1) { if (x) {return 1} else {return 1}}")
	c.TestNotMissing("while (0) {} return 1")

	// TODO(user): The current algorithm will not detect this case. It is
	// still computable in most cases.
	c.TestNotMissing("while (1) {} return 0")
	c.TestMissing("while (false) return 1")

	// Not known.
	c.TestMissing("while(x) { return 1 }")
}

// Java: CheckMissingReturnTest.testMultiConditions()
func (c *CheckMissingReturnTest) TestMultiConditions() {
	c.TestMissing("if (a) { } else { while (1) {return 1} }")
	c.TestNotMissing("if (a) { return 1} else { while (1) {return 1} }")
}

// Java: CheckMissingReturnTest.testIssue779()
func (c *CheckMissingReturnTest) TestIssue779() {
	c.TestNotMissing(
		`var a = f(); try { alert(); if (a > 0) return 1; }
finally { a = 5; } return 2;
`)
}

// Java: CheckMissingReturnTest.testConstructors()
func (c *CheckMissingReturnTest) TestConstructors() {
	c.TestSame("/** @constructor */ function foo() {} ")

	constructorWithReturn :=
		`/** @constructor 
 * @return {!foo} */ function foo() {
 if (!(this instanceof foo)) { return new foo; } }
`
	c.TestSame(constructorWithReturn)
}

// Java: CheckMissingReturnTest.testClosureAsserts()
func (c *CheckMissingReturnTest) TestClosureAsserts() {
	closureDefs :=
		`/** @const */ var goog = {};
goog.asserts = {};
goog.asserts.fail = function(x) {};
`

	c.TestNotMissing(closureDefs + "goog.asserts.fail('');")

	c.TestNotMissing(
		closureDefs + "switch (x) { case 1: return 1; default: goog.asserts.fail(''); }")
}

// Java: CheckMissingReturnTest.testInfiniteLoops()
func (c *CheckMissingReturnTest) TestInfiniteLoops() {
	c.TestNotMissing("while (true) { x = y; if (x === 0) { return 1; } }")
	c.TestNotMissing("for (;true;) { x = y; if (x === 0) { return 1; } }")
	c.TestNotMissing("for (;;) { x = y; if (x === 0) { return 1; } }")
}

// Java: CheckMissingReturnTest.createFunction(String, String)
func CheckMissingReturnTest_CreateFunction(returnType string, body string) string {
	return "/** @return {" + returnType + "} */ function foo() {" + body + "}"
}

// Java: CheckMissingReturnTest.createShorthandFunctionInObjLit(String, String)
func CheckMissingReturnTest_CreateShorthandFunctionInObjLit(returnType string, body string) string {
	return strings.ReplaceAll(
		strings.ReplaceAll(`var obj = {
  /** @return {RETURN_TYPE} */
  foo() {BODY}
}
`, "RETURN_TYPE", returnType),
		"BODY", body)
}

// Java: CheckMissingReturnTest.testMissingInTraditionalFunction(String, String)
func (c *CheckMissingReturnTest) TestMissingInTraditionalFunction(returnType string, body string) {
	js := CheckMissingReturnTest_CreateFunction(returnType, body)
	c.TestWarning(js, CheckMissingReturn_MISSING_RETURN_STATEMENT)
}

// Java: CheckMissingReturnTest.testNotMissingInTraditionalFunction(String, String)
func (c *CheckMissingReturnTest) TestNotMissingInTraditionalFunction(returnType string, body string) {
	c.TestSame(CheckMissingReturnTest_CreateFunction(returnType, body))
}

// Java: CheckMissingReturnTest.testMissingInShorthandFunction(String, String)
func (c *CheckMissingReturnTest) TestMissingInShorthandFunction(returnType string, body string) {
	js := CheckMissingReturnTest_CreateShorthandFunctionInObjLit(returnType, body)
	c.TestWarning(js, CheckMissingReturn_MISSING_RETURN_STATEMENT)
}

// Java: CheckMissingReturnTest.testNotMissingInShorthandFunction(String, String)
func (c *CheckMissingReturnTest) TestNotMissingInShorthandFunction(returnType string, body string) {
	c.TestSame(CheckMissingReturnTest_CreateShorthandFunctionInObjLit(returnType, body))
}

// Java: CheckMissingReturnTest.testMissing(String, String)
func (c *CheckMissingReturnTest) TestMissingWithReturnType(returnType string, body string) {
	c.TestMissingInTraditionalFunction(returnType, body)
	c.TestMissingInShorthandFunction(returnType, body)
}

// Creates function with return type {number}
//
// Java: CheckMissingReturnTest.testMissing(String)
func (c *CheckMissingReturnTest) TestMissing(body string) {
	c.TestMissingWithReturnType("number", body)
}

// Java: CheckMissingReturnTest.testNotMissing(String, String)
func (c *CheckMissingReturnTest) TestNotMissingWithReturnType(returnType string, body string) {
	c.TestNotMissingInTraditionalFunction(returnType, body)
	c.TestNotMissingInShorthandFunction(returnType, body)
}

// Creates function with return type {number}
//
// Java: CheckMissingReturnTest.testNotMissing(String)
func (c *CheckMissingReturnTest) TestNotMissing(body string) {
	c.TestNotMissingWithReturnType("number", body)
}

// Java: CheckMissingReturnTest.testArrowFunctions_noReturn()
func (c *CheckMissingReturnTest) TestArrowFunctions_noReturn() {
	c.TestNoWarning(
		`/** @return {undefined} */
() => {}
`)
}

// Java: CheckMissingReturnTest.testArrowFunctions_expressionBody1()
func (c *CheckMissingReturnTest) TestArrowFunctions_expressionBody1() {
	c.TestSame(
		`/** @return {number} */
() => 1
`)
}

// Java: CheckMissingReturnTest.testArrowFunctions_expressionBody2()
func (c *CheckMissingReturnTest) TestArrowFunctions_expressionBody2() {
	c.TestSame(
		`/** @return {number} */
(a) => (a > 3) ? 1 : 0
`)
}

// Java: CheckMissingReturnTest.testArrowFunctions_block()
func (c *CheckMissingReturnTest) TestArrowFunctions_block() {
	c.TestSame(
		`/** @return {number} */
(a) => { if (a > 3) { return 1; } else { return 0; }}
`)
}

// Java: CheckMissingReturnTest.testArrowFunctions_blockMissingReturn()
func (c *CheckMissingReturnTest) TestArrowFunctions_blockMissingReturn() {
	c.TestWarning(
		`/** @return {number} */
(a) => { if (a > 3) { return 1; } else { } }
`,
		CheckMissingReturn_MISSING_RETURN_STATEMENT)
}

// Java: CheckMissingReturnTest.testArrowFunctions_objectLiteralExpression()
func (c *CheckMissingReturnTest) TestArrowFunctions_objectLiteralExpression() {
	c.TestSame("(a) => ({foo: 1});")
}

// Java: CheckMissingReturnTest.testGeneratorFunctionWithoutSpecifiedReturnType()
func (c *CheckMissingReturnTest) TestGeneratorFunctionWithoutSpecifiedReturnType() {
	c.TestNoWarning("function *gen() {}")

	c.TestNoWarning(
		`/** @return {!Generator<number>} */ // no yields is OK
function *gen() {}
`)

	c.TestNoWarning(
		`/** @return {!Generator<number>} */ // no yields is OK
function *gen() { return; }
`)

	c.TestNoWarning(
		`/** @return {!Generator<number>} */ // one yield is OK
function *gen() {
 yield 1;
}
`)

	c.TestNoWarning(
		`/** @return {!Object} */ // Return type more vague than Generator is also OK
function *gen() {}
`)
}

// Java: CheckMissingReturnTest.testGeneratorFunctionWithSpecifiedReturnType()
func (c *CheckMissingReturnTest) TestGeneratorFunctionWithSpecifiedReturnType() {
	c.TestWarning(
		`/** @return {!Iterable<number, number>} */
function *gen() {
 yield 1;
}
`,
		CheckMissingReturn_MISSING_RETURN_STATEMENT)
}

// Java: CheckMissingReturnTest.testGeneratorFunctionWithIteratorLikeSpecifiedReturnType()
func (c *CheckMissingReturnTest) TestGeneratorFunctionWithIteratorLikeSpecifiedReturnType() {
	c.TestWarning(
		`/** @return {!IteratorLike<number, number>} */
function *gen() {
 yield 1;
}
`,
		CheckMissingReturn_MISSING_RETURN_STATEMENT)
}

// Java: CheckMissingReturnTest.testAsyncFunction_noJSDoc()
func (c *CheckMissingReturnTest) TestAsyncFunction_noJSDoc() {
	// Note: we add the alert because CheckMissingReturn never warns on functions with empty bodies.
	c.TestNoWarning("async function foo() { alert(1); }")
}

// Java: CheckMissingReturnTest.testAsyncFunction_returnsUndefined()
func (c *CheckMissingReturnTest) TestAsyncFunction_returnsUndefined() {
	c.TestNoWarning("/** @return {!Promise<undefined>} */ async function foo() { alert(1); }")
}

// Java: CheckMissingReturnTest.testAsyncFunction_returnsUnknown()
func (c *CheckMissingReturnTest) TestAsyncFunction_returnsUnknown() {
	c.TestNoWarning("/** @return {!Promise<?>} */ async function foo() { alert(1); }")
}

// Java: CheckMissingReturnTest.testAsyncFunction_returnsNumber()
func (c *CheckMissingReturnTest) TestAsyncFunction_returnsNumber() {
	c.TestNoWarning("/** @return {!Promise<number>} */ async function foo() { return 1; }")

	c.TestWarning(
		"/** @return {!Promise<string>} */ async function foo() { alert(1); }",
		CheckMissingReturn_MISSING_RETURN_STATEMENT)
}

// Java: CheckMissingReturnTest.testClosureAssertsFailPreventsWarning()
func (c *CheckMissingReturnTest) TestClosureAssertsFailPreventsWarning() {
	input :=
		`/** @return {string} */
function foo(param) {
  if (param) {
    return 'success';
  }
  fail();
}
`
	c.TestWarning("function fail() {}"+input, CheckMissingReturn_MISSING_RETURN_STATEMENT)

	c.TestNoWarning("/** @closurePrimitive {asserts.fail} */ function fail() {}" + input)
}
