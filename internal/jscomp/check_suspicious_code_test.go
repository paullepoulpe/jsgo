// Port of com/google/javascript/jscomp/CheckSuspiciousCodeTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/CheckSuspiciousCodeTest.java)

/*
 * Copyright 2012 The Closure Compiler Authors.
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

// CheckSuspiciousCodeTest: Tests for CheckSuspiciousCode
type CheckSuspiciousCodeTest struct {
	CompilerTestCaseBase
}

// Java: CheckSuspiciousCodeTest()
func NewCheckSuspiciousCodeTest() *CheckSuspiciousCodeTest {
	c := &CheckSuspiciousCodeTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestCheckSuspiciousCodeTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCheckSuspiciousCodeTest() })
}

// Java: CheckSuspiciousCodeTest.getProcessor(Compiler)
func (c *CheckSuspiciousCodeTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewCombinedCompilerPass(compiler, NewCheckSuspiciousCode())
}

// Java: CheckSuspiciousCodeTest.setUp()
func (c *CheckSuspiciousCodeTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableParseTypeInfo()
}

// TestSuspiciousBreakingOutOfOptionalChain is the @Test method
// suspiciousBreakingOutOfOptionalChain(); RunJUnit4 runs the methods whose
// name starts with Test (Go-only: language-adapter).
//
// Java: CheckSuspiciousCodeTest.suspiciousBreakingOutOfOptionalChain()
func (c *CheckSuspiciousCodeTest) TestSuspiciousBreakingOutOfOptionalChain() {
	e := CheckSuspiciousCode_SUSPICIOUS_BREAKING_OUT_OF_OPTIONAL_CHAIN

	c.TestSame("a?.b.c")
	c.TestWarning("(a?.b).c", e)

	c.TestSame("a.b?.[c][d]")
	c.TestWarning("(a.b?.[c])[d]", e)

	c.TestSame("a.b?.()()")
	c.TestWarning("(a.b?.())()", e)

	c.TestSame("a(b?.c);")
	c.TestSame("a[b?.c];")
}

// Java: CheckSuspiciousCodeTest.testSuspiciousSemi()
func (c *CheckSuspiciousCodeTest) TestSuspiciousSemi() {
	e := CheckSuspiciousCode_SUSPICIOUS_SEMICOLON

	c.TestSame("if(x()) x = y;")
	c.TestWarning("if(x()); x = y;", e) // I've had this bug, damned ;
	c.TestSame("if(x()){} x = y;")

	c.TestSame("if(x()) x = y; else y=z;")
	c.TestWarning("if(x()); else y=z;", e)
	c.TestSame("if(x()){} else y=z;")
	c.TestWarning("if(x()) x = y; else;", e)
	c.TestSame("if(x()) x = y; else {}")

	c.TestSame("while(x()) x = y;")
	c.TestWarning("while(x()); x = y;", e)
	c.TestSame("while(x()){} x = y;")
	c.TestWarning("while(x()); {x = y}", e)
	c.TestSame("while(x()){} {x = y}")

	c.TestSame("for(;;) x = y;")
	c.TestWarning("for(;;); x = y;", e)
	c.TestSame("for(;;){} x = y;")
	c.TestSame("for(x in y) x = y;")
	c.TestWarning("for(x in y); x = y;", e)
	c.TestSame("for(x in y){} x = y;")

	c.TestSame("var y = [1, 2, 3]; for(x of y) console.log(x);")
	c.TestWarning("var y = [1, 2, 3]; for(x of y); console.log(x);", e)
	c.TestSame("var y = [1, 2, 3]; for(x of y){} console.log(x);")

	c.TestSame("async () => { var y = [1, 2, 3]; for await (x of y) console.log(x); }")
	c.TestWarning("async () => { var y = [1, 2, 3]; for await (x of y); console.log(x); }", e)
	c.TestSame("async () => { var y = [1, 2, 3]; for await (x of y){} console.log(x); }")
}

// Java: CheckSuspiciousCodeTest.testSuspiciousIn()
func (c *CheckSuspiciousCodeTest) TestSuspiciousIn() {
	c.TestWarning("'foo' in 1", CheckSuspiciousCode_SUSPICIOUS_IN_OPERATOR)
	c.TestWarning("'foo' in 'test'", CheckSuspiciousCode_SUSPICIOUS_IN_OPERATOR)
	c.TestWarning("'foo' in NaN", CheckSuspiciousCode_SUSPICIOUS_IN_OPERATOR)
	c.TestWarning("'foo' in undefined", CheckSuspiciousCode_SUSPICIOUS_IN_OPERATOR)
	c.TestWarning("'foo' in Infinity", CheckSuspiciousCode_SUSPICIOUS_IN_OPERATOR)
	c.TestWarning("'foo' in true", CheckSuspiciousCode_SUSPICIOUS_IN_OPERATOR)
	c.TestWarning("'foo' in false", CheckSuspiciousCode_SUSPICIOUS_IN_OPERATOR)
	c.TestWarning("'foo' in null", CheckSuspiciousCode_SUSPICIOUS_IN_OPERATOR)
	c.TestWarning("'foo' in !Object", CheckSuspiciousCode_SUSPICIOUS_IN_OPERATOR)
	c.TestSame("'foo' in Object")
	c.TestSame("'foo' in {}")
}

// Java: CheckSuspiciousCodeTest.testForOf()
func (c *CheckSuspiciousCodeTest) TestForOf() {
	c.TestSame("var y = [1, 2, 3]; for (var x of y) console.log(x);")
	c.TestSame("var y = [1, 2, 3]; for (var x of 'test') console.log(x);")
	c.TestSame("for (var x of 123) console.log(x);")
	c.TestSame("for (var x of false) console.log(x);")
	c.TestSame("for (var x of true) console.log(x);")
	c.TestSame("for (var x of undefined) console.log(x);")
	c.TestSame("for (var x of NaN) console.log(x);")
	c.TestSame("for (var x of Infinity) console.log(x);")
	c.TestSame("for (var x of null) console.log(x);")
}

// Java: CheckSuspiciousCodeTest.testForAwaitOf()
func (c *CheckSuspiciousCodeTest) TestForAwaitOf() {
	c.TestSame("async () => { var y = [1, 2, 3]; for await (var x of y) console.log(x); }")
	c.TestSame("async () => { var y = [1, 2, 3]; for await (var x of 'test') console.log(x); }")
	c.TestSame("async () => { for await (var x of 123) console.log(x); }")
	c.TestSame("async () => { for await (var x of false) console.log(x); }")
	c.TestSame("async () => { for await (var x of true) console.log(x); }")
	c.TestSame("async () => { for await (var x of undefined) console.log(x); }")
	c.TestSame("async () => { for await (var x of NaN) console.log(x); }")
	c.TestSame("async () => { for await (var x of Infinity) console.log(x); }")
	c.TestSame("async () => { for await (var x of null) console.log(x); }")
}

// Java: CheckSuspiciousCodeTest.testReportNaN(String)
func (c *CheckSuspiciousCodeTest) TestReportNaN(js string) {
	c.TestWarning(js, CheckSuspiciousCode_SUSPICIOUS_COMPARISON_WITH_NAN)
}

// Java: CheckSuspiciousCodeTest.testComparison1()
func (c *CheckSuspiciousCodeTest) TestComparison1() {
	c.TestReportNaN("x == NaN")
	c.TestReportNaN("x != NaN")
	c.TestReportNaN("x === NaN")
	c.TestReportNaN("x !== NaN")
	c.TestReportNaN("x < NaN")
	c.TestReportNaN("x <= NaN")
	c.TestReportNaN("x > NaN")
	c.TestReportNaN("x >= NaN")
}

// Java: CheckSuspiciousCodeTest.testComparison2()
func (c *CheckSuspiciousCodeTest) TestComparison2() {
	c.TestReportNaN("NaN == x")
	c.TestReportNaN("NaN != x")
	c.TestReportNaN("NaN === x")
	c.TestReportNaN("NaN !== x")
	c.TestReportNaN("NaN < x")
	c.TestReportNaN("NaN <= x")
	c.TestReportNaN("NaN > x")
	c.TestReportNaN("NaN >= x")
}

// Java: CheckSuspiciousCodeTest.testComparison3()
func (c *CheckSuspiciousCodeTest) TestComparison3() {
	c.TestReportNaN("x == 0/0")
	c.TestReportNaN("x != 0/0")
	c.TestReportNaN("x === 0/0")
	c.TestReportNaN("x !== 0/0")
	c.TestReportNaN("x < 0/0")
	c.TestReportNaN("x <= 0/0")
	c.TestReportNaN("x > 0/0")
	c.TestReportNaN("x >= 0/0")
}

// Java: CheckSuspiciousCodeTest.testComparison4()
func (c *CheckSuspiciousCodeTest) TestComparison4() {
	c.TestReportNaN("0/0 == x")
	c.TestReportNaN("0/0 != x")
	c.TestReportNaN("0/0 === x")
	c.TestReportNaN("0/0 !== x")
	c.TestReportNaN("0/0 < x")
	c.TestReportNaN("0/0 <= x")
	c.TestReportNaN("0/0 > x")
	c.TestReportNaN("0/0 >= x")
}

// Java: CheckSuspiciousCodeTest.testComparison5()
func (c *CheckSuspiciousCodeTest) TestComparison5() {
	c.TestReportNaN("x == Number.NaN")
	c.TestReportNaN("x != Number.NaN")
	c.TestReportNaN("x === Number.NaN")
	c.TestReportNaN("x !== Number.NaN")
	c.TestReportNaN("x < Number.NaN")
	c.TestReportNaN("x <= Number.NaN")
	c.TestReportNaN("x > Number.NaN")
	c.TestReportNaN("x >= Number.NaN")
}

// Java: CheckSuspiciousCodeTest.testComparison6()
func (c *CheckSuspiciousCodeTest) TestComparison6() {
	c.TestReportNaN("Number.NaN == x")
	c.TestReportNaN("Number.NaN != x")
	c.TestReportNaN("Number.NaN === x")
	c.TestReportNaN("Number.NaN !== x")
	c.TestReportNaN("Number.NaN < x")
	c.TestReportNaN("Number.NaN <= x")
	c.TestReportNaN("Number.NaN > x")
	c.TestReportNaN("Number.NaN >= x")
}

// Java: CheckSuspiciousCodeTest.testInstanceOf()
func (c *CheckSuspiciousCodeTest) TestInstanceOf() {
	c.TestReportInstanceOf("''", "String")
	c.TestReportInstanceOf("4", "Number")
	c.TestReportInstanceOf("-4", "Number")
	c.TestReportInstanceOf("null", "Number")
	c.TestReportInstanceOf("true", "Boolean")
	c.TestReportInstanceOf("false", "Boolean")
	c.TestReportInstanceOf("!true", "Boolean")
	c.TestReportInstanceOf("undefined", "Number")
	c.TestReportInstanceOf("Infinity", "Number")
	c.TestReportInstanceOf("NaN", "Number")
	c.TestReportInstanceOf(
		`/** @constructor */ function Foo() {}; var foo = new Foo();
!foo
`,
		"Foo")

	c.TestReportInstanceOf("(4 + 5)", "Number")
	c.TestReportInstanceOf("('a' + 'b')", "String")
	c.TestReportInstanceOf("('a' + 5)", "String")

	// The following five are from "Wat — Destroy All Software Talks"
	c.TestReportInstanceOf("([] + [])", "String")
	c.TestReportInstanceOf("([] + {})", "String")
	c.TestReportInstanceOf("({} + [])", "String")
	c.TestReportInstanceOf("({c:1} + [])", "String")
	c.TestReportInstanceOf("({} + {})", "Number") // NaN

	c.TestReportInstanceOf("!''", "String")
	c.TestReportInstanceOf("!4", "Number")
	c.TestReportInstanceOf("!(new Boolean(true))", "Boolean")
	c.TestReportInstanceOf("!(new Object())", "Object")
	c.TestReportInstanceOf("!Object.prototype", "Object")
	c.TestReportInstanceOf("!Function", "Object")
	c.TestReportInstanceOf("!func()", "String")
	c.TestReportInstanceOf("!({})", "Object")
	c.TestReportInstanceOf(
		`/** @constructor */ function Foo() {
!this
`,
		`Foo;
}
`)

	c.TestSame("new String('') instanceof String")
	c.TestSame("new Number(4) instanceof Number")
	c.TestSame("var a = new Number(4); a instanceof Number")
	c.TestSame("new Boolean(true) instanceof Boolean")
	c.TestSame("new Object() instanceof Object")
	c.TestSame("Object.prototype instanceof Object")
	c.TestSame("Function instanceof Object")
	c.TestSame("func() instanceof String")
	c.TestSame("({}) instanceof Object")
	c.TestSame(
		`/** @constructor */ function Foo() {
 var a = this instanceof Foo; }
`)

	c.TestSame("(()=>42) instanceof Function")
	c.TestSame("class Person{} Person instanceof Function")
	c.TestSame(
		`class Person{}
var peter = new Person();
peter instanceof Person
`)
	c.TestSame("taggedTemplate`${tagged}Temp` instanceof Function")
}

// Java: CheckSuspiciousCodeTest.testReportInstanceOf(String, String)
func (c *CheckSuspiciousCodeTest) TestReportInstanceOf(left string, right string) {
	c.TestWarning(
		left+" instanceof "+right, CheckSuspiciousCode_SUSPICIOUS_INSTANCEOF_LEFT_OPERAND)
}

// Java: CheckSuspiciousCodeTest.testCheckNegatedLeftOperandOfInOperator()
func (c *CheckSuspiciousCodeTest) TestCheckNegatedLeftOperandOfInOperator() {
	c.TestSame("if (!(x in y)) {}")
	c.TestSame("if (('' + !x) in y) {}")
	c.TestWarning(
		"if (!x in y) {}", CheckSuspiciousCode_SUSPICIOUS_NEGATED_LEFT_OPERAND_OF_IN_OPERATOR)
}

// Java: CheckSuspiciousCodeTest.testSuspiciousLeftArgumentOfLogicalOperator()
func (c *CheckSuspiciousCodeTest) TestSuspiciousLeftArgumentOfLogicalOperator() {
	c.TestSame("if (x && y) {}")
	c.TestWarning("if (false && y) {}", CheckSuspiciousCode_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR)
	c.TestWarning("if (true && y) {}", CheckSuspiciousCode_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR)
	c.TestWarning("if (false || y) {}", CheckSuspiciousCode_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR)
	c.TestWarning("if (true || y) {}", CheckSuspiciousCode_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR)
	c.TestWarning("if (0 || y) {}", CheckSuspiciousCode_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR)
	c.TestWarning("if (null && y) {}", CheckSuspiciousCode_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR)
	c.TestWarning("if (NaN || y) {}", CheckSuspiciousCode_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR)
	c.TestWarning("if ({} && y) {}", CheckSuspiciousCode_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR)
	c.TestWarning("if (void x && y) {}", CheckSuspiciousCode_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR)

	c.TestWarning("if (x || true || y) {}", CheckSuspiciousCode_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR)
	c.TestWarning("if (x && false || y) {}", CheckSuspiciousCode_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR)
	c.TestSame("if (x || false || y) {}")
	c.TestSame("if (x && true || y) {}")
}

// Java: CheckSuspiciousCodeTest.testSuspiciousLeftArgumentOfLogicalOperator_propAccesses()
func (c *CheckSuspiciousCodeTest) TestSuspiciousLeftArgumentOfLogicalOperator_propAccesses() {
	c.EnableTypeCheck()
	// Cases that must not warn as LHS is not always falsy.
	c.TestNoWarning("/** @type {null|{b:string}} */ let a; if (a?.b && true) {}")
	c.TestNoWarning("/** @type {null|{b:string}} */ let a; if (a?.[b] && c) {}")
	c.TestNoWarning("/** @type {{b: (null|string)}} */ let a; if (a?.b && true) {}")
	c.TestNoWarning("/** @type {{b: (null|string)}} */ let a; if (a?.[b] && true) {}")

	// Cases that could warn but don't, as we assume that GETPROPs and GETELEMs type information is
	// likely wrong
	// normal prop access
	c.TestNoWarning("/** @type {{b:null}} */ let a; if( a[b] && true) {}")
	c.TestNoWarning("/** @type {{b:null}} */ let a; if( a.b && true) {}")
	// optChain prop access
	c.TestNoWarning("/** @type {{b:null}} */ let a; if( a?.[b] && true) {}")
	c.TestNoWarning("/** @type {{b:null}} */ let a; if( a?.b && true) {}")

	// normal call
	c.TestWarning(
		"let a = {}; /** @return {null} */ a.b = function() {}; if (a.b() && c) {}",
		CheckSuspiciousCode_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR)
	// optChain calls
	c.TestWarning(
		"let a = {}; /** @return {null} */ a.b = function() {}; if (a?.b() && c) {}",
		CheckSuspiciousCode_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR)
	c.TestWarning(
		"/** @type {null} */ let a = null; if (a?.b() && c) {}",
		CheckSuspiciousCode_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR)
}

// Java: CheckSuspiciousCodeTest.testSuspiciousLeftArgumentOfLogicalOperator_typeBased()
func (c *CheckSuspiciousCodeTest) TestSuspiciousLeftArgumentOfLogicalOperator_typeBased() {
	c.EnableTypeCheck()
	prefix :=
		`/** @return {!Object} */ function truthy() { return {}; }
/** @return {null|undefined} */ function falsy() { return null; }
/** @return {number} */ function number() { return 42; }
/** @return {?} */ function unknown() { return 42; }
/** @const */ var ns = {};
/** @type {!Object<!Object>} */ ns.truthy = {};
/** @type {null} */ ns.falsy = null;
`

	c.TestSame("var x = x || {};")
	c.TestSame(prefix + "if (number() && y) {}")
	c.TestSame(prefix + "if (unknown() && y) {}")
	c.TestSame(prefix + "if (ns.truthy && y) {}")
	c.TestSame(prefix + "if (ns.falsy && y) {}")
	c.TestSame(prefix + "if (!ns.truthy && y) {}")
	c.TestSame(prefix + "if (!ns.falsy && y) {}")

	c.TestSame(prefix + "if (number() || y) {}")
	c.TestSame(prefix + "if (unknown() || y) {}")
	c.TestSame(prefix + "if (ns.truthy || y) {}")
	c.TestSame(prefix + "if (ns.falsy || y) {}")
	c.TestSame(prefix + "if (!ns.truthy || y) {}")
	c.TestSame(prefix + "if (!ns.falsy || y) {}")

	c.TestSame(prefix + "if (ns.falsy || ns.truthy || y) {}")
	c.TestSame(prefix + "if (ns.falsy && ns.truthy || y) {}")
	c.TestSame(prefix + "if (ns.truthy[ns.truthy] || y) {}")
	c.TestSame(prefix + "if (ns.truthy[truthy()] || y) {}")
	c.TestSame(prefix + "if ((truthy(), !ns.falsy) || y) {}")

	// It's always okay to second-guess a type-based always-truthy.
	c.TestSame(prefix + "if (truthy() && y) {}")
	c.TestSame(prefix + "if (!truthy() && y) {}")
	c.TestSame(prefix + "if (truthy() || y) {}")
	c.TestSame(prefix + "if (!truthy() || y) {}")

	c.TestWarning(prefix+"if (falsy() && y) {}", CheckSuspiciousCode_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR)
	c.TestWarning(prefix+"if (!falsy() && y) {}", CheckSuspiciousCode_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR)
	c.TestWarning(prefix+"if (falsy() || y) {}", CheckSuspiciousCode_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR)
	c.TestWarning(prefix+"if (!falsy() || y) {}", CheckSuspiciousCode_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR)

	c.TestWarning(
		prefix+"if ((truthy(), !falsy()) || y) {}", CheckSuspiciousCode_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR)
}

// Java: CheckSuspiciousCodeTest.testSuspiciousLeftArgumentOfLogicalOperator_guardedInitialization_noWarning()
func (c *CheckSuspiciousCodeTest) TestSuspiciousLeftArgumentOfLogicalOperator_guardedInitialization_noWarning() {
	// In this example, 'x' is provably falsy, but it's a common pattern to initialize a global
	// variable that may already exist before the script began.
	c.EnableTypeCheck()
	c.TestSame("var x = x || {};")
}

// Java: CheckSuspiciousCodeTest.testSuspiciousLeftArgumentOfLogicalOperator_truthyFunctionCalls_noWarning()
func (c *CheckSuspiciousCodeTest) TestSuspiciousLeftArgumentOfLogicalOperator_truthyFunctionCalls_noWarning() {
	// There are commonplace examples where functions that are annotated to return truthy results
	// are actually lying, so we need to back off from warning about these.
	c.EnableTypeCheck()
	prefix := "/** @return {!Object} */ function truthy() { return {}; }\n"
	c.TestSame(prefix + "if (truthy() && y) {}")
	c.TestSame(prefix + "if (!truthy() && y) {}")
	c.TestSame(prefix + "if (truthy() || y) {}")
	c.TestSame(prefix + "if (!truthy() || y) {}")
}

// Java: CheckSuspiciousCodeTest.testSuspiciousLeftArgumentOfLogicalOperator_qualifiedName_noWarning()
func (c *CheckSuspiciousCodeTest) TestSuspiciousLeftArgumentOfLogicalOperator_qualifiedName_noWarning() {
	// Type information on qualified names can be wrong in both directions (e.g. an extern might
	// assert that a name exists, but code must still verify to be sure).
	c.EnableTypeCheck()
	prefix :=
		`/** @const */ var ns = {};
/** @type {!Object<!Object>} */ ns.truthy = {};
/** @type {null} */ ns.falsy = null;
`
	c.TestSame(prefix + "if (ns.truthy && y) {}")
	c.TestSame(prefix + "if (ns.falsy && y) {}")
	c.TestSame(prefix + "if (!ns.truthy && y) {}")
	c.TestSame(prefix + "if (!ns.falsy && y) {}")
	c.TestSame(prefix + "if (ns.truthy || y) {}")
	c.TestSame(prefix + "if (ns.falsy || y) {}")
	c.TestSame(prefix + "if (!ns.truthy || y) {}")
	c.TestSame(prefix + "if (!ns.falsy || y) {}")
}

// Java: CheckSuspiciousCodeTest.testSuspiciousLeftArgumentOfLogicalOperator_primitiveOrUnknownTypes_noWarning()
func (c *CheckSuspiciousCodeTest) TestSuspiciousLeftArgumentOfLogicalOperator_primitiveOrUnknownTypes_noWarning() {
	// Primitive types (like number) and unknown types can be either true or false, so don't warn.
	c.EnableTypeCheck()
	prefix :=
		`/** @return {number} */ function number() { return 42; }
/** @return {?} */ function unknown() { return 42; }
`
	c.TestSame(prefix + "if (number() && y) {}")
	c.TestSame(prefix + "if (unknown() && y) {}")
	c.TestSame(prefix + "if (number() || y) {}")
	c.TestSame(prefix + "if (unknown() || y) {}")
	c.TestSame(prefix + "if (!number() && y) {}")
	c.TestSame(prefix + "if (!unknown() && y) {}")
	c.TestSame(prefix + "if (!number() || y) {}")
	c.TestSame(prefix + "if (!unknown() || y) {}")
}

// Java: CheckSuspiciousCodeTest.testSuspiciousLeftArgumentOfLogicalOperator_deeperNesting()
func (c *CheckSuspiciousCodeTest) TestSuspiciousLeftArgumentOfLogicalOperator_deeperNesting() {
	c.EnableTypeCheck()
	prefix :=
		`/** @return {!Object} */ function truthy() { return {}; }
function falsy() { return; }
/** @const */ var ns = {};
/** @type {!Object<!Object>} */ ns.truthy = {};
/** @type {null} */ ns.falsy = null;
`

	// When the above constructs (that don't warn) are combined more deeply, they still don't warn.
	c.TestSame(prefix + "if (ns.falsy || ns.truthy || y) {}")
	c.TestSame(prefix + "if (ns.falsy && ns.truthy || y) {}")
	c.TestSame(prefix + "if (ns.truthy[ns.truthy] || y) {}")
	c.TestSame(prefix + "if (ns.truthy[truthy()] || y) {}")

	// Comma operator only recurses into RHS, so the truthy() below is ignored.
	c.TestSame(prefix + "if ((truthy(), !ns.falsy) || y) {}")

	// But if the RHS of the comma would normally warn, then we should warn on the comma, too.
	c.TestWarning(
		prefix+"if ((truthy(), falsy()) || y) {}", CheckSuspiciousCode_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR)
}

// Java: CheckSuspiciousCodeTest.testSuspiciousLeftArgumentOfLogicalOperator_falsyFunctionCalls_warn()
func (c *CheckSuspiciousCodeTest) TestSuspiciousLeftArgumentOfLogicalOperator_falsyFunctionCalls_warn() {
	// Falsy function call results are less likely to be a false positive, since there are far
	// fewer lies about them in the wild.
	c.EnableTypeCheck()
	prefix :=
		`/** @return {!Object} */ function truthy() { return {}; }
function falsy() { return; }
`

	c.TestWarning(prefix+"if (falsy() && y) {}", CheckSuspiciousCode_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR)
	c.TestWarning(prefix+"if (!falsy() && y) {}", CheckSuspiciousCode_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR)
	c.TestWarning(prefix+"if (falsy() || y) {}", CheckSuspiciousCode_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR)
	c.TestWarning(prefix+"if (!falsy() || y) {}", CheckSuspiciousCode_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR)
}
