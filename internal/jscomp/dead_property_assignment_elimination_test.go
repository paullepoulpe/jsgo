// Port of com/google/javascript/jscomp/DeadPropertyAssignmentEliminationTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/DeadPropertyAssignmentEliminationTest.java)

/*
 * Copyright 2016 The Closure Compiler Authors.
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

type DeadPropertyAssignmentEliminationTest struct {
	CompilerTestCaseBase
}

// Java: DeadPropertyAssignmentEliminationTest()
func NewDeadPropertyAssignmentEliminationTest() *DeadPropertyAssignmentEliminationTest {
	c := &DeadPropertyAssignmentEliminationTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestDeadPropertyAssignmentEliminationTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewDeadPropertyAssignmentEliminationTest() })
}

// Java: DeadPropertyAssignmentEliminationTest.setUp()
func (c *DeadPropertyAssignmentEliminationTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableGatherExternProperties()
}

// Java: DeadPropertyAssignmentEliminationTest.testBasic()
func (c *DeadPropertyAssignmentEliminationTest) TestBasic() {
	c.TestSame(
		"var foo = function() {\n" +
			"  this.a = 20;\n" +
			"}\n")

	c.Test(
		"var foo = function() {\n"+
			"  this.a = 10;\n"+
			"  this.a = 20;\n"+
			"}\n",
		"var foo = function() {\n"+
			"  10;\n"+
			"  this.a = 20;\n"+
			"}\n")

	c.TestSame(
		"var foo = function() {\n" +
			"  this.a = 20;\n" +
			"  this.a = this.a + 20;\n" +
			"}\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testMultipleProperties()
func (c *DeadPropertyAssignmentEliminationTest) TestMultipleProperties() {
	c.Test(
		"var foo = function() {\n"+
			"  this.a = 10;\n"+
			"  this.b = 15;\n"+
			"  this.a = 20;\n"+
			"}\n",
		"var foo = function() {\n"+
			"  10;\n"+
			"  this.b = 15;\n"+
			"  this.a = 20;\n"+
			"}\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testNonStandardAssign()
func (c *DeadPropertyAssignmentEliminationTest) TestNonStandardAssign() {
	c.Test(
		"var foo = function() {\n"+
			"  this.a = 10;\n"+
			"  this.a += 15;\n"+
			"  this.a = 20;\n"+
			"}\n",
		"var foo = function() {\n"+
			"  this.a = 10;\n"+
			"  this.a + 15;\n"+
			"  this.a = 20;\n"+
			"}\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testChainingPropertiesAssignments()
func (c *DeadPropertyAssignmentEliminationTest) TestChainingPropertiesAssignments() {
	c.Test(
		"var foo = function() {\n"+
			"  this.a = this.b = this.c = 10;\n"+
			"  this.b = 15;\n"+
			"}\n",
		"var foo = function() {\n"+
			"  this.a = this.c = 10;\n"+
			"  this.b = 15;\n"+
			"}\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testConditionalProperties()
func (c *DeadPropertyAssignmentEliminationTest) TestConditionalProperties() {
	// We don't handle conditionals at all.
	c.TestSame(
		"var foo = function() {\n" +
			"  this.a = 10;\n" +
			"  if (true) { this.a = 20; } else { this.a = 30; }\n" +
			"}\n")

	// We don't handle conditionals at all.
	c.TestSame(
		"var bar = function(x)  {};\n" +
			"var foo = function() {\n" +
			"  this.a = 10; // must preserve this assignment\n" +
			"  bar?.(this.a = 20);\n" +
			"}\n")

	// However, we do handle everything up until the conditional.
	c.Test(
		"var foo = function() {\n"+
			"  this.a = 10;\n"+
			"  this.a = 20;\n"+
			"  if (true) { this.a = 20; } else { this.a = 30; }\n"+
			"}\n",
		"var foo = function() {\n"+
			"  10;\n"+
			"  this.a = 20;\n"+
			"  if (true) { this.a = 20; } else { this.a = 30; }\n"+
			"}\n")

	// However, we do handle everything up until the conditional.
	c.Test(
		"var bar = function(x)  {};\n"+
			"var foo = function() {\n"+
			"  this.a = 10; // this is a dead assignment and must be eliminated\n"+
			"  this.a = 20;\n"+
			"  bar?.(this.a = 30);\n"+
			"}\n",
		"var bar = function(x)  {};\n"+
			"var foo = function() {\n"+
			"  10;\n"+
			"  this.a = 20;\n"+
			"  bar?.(this.a = 30);\n"+
			"}\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testQualifiedNamePrefixAssignment()
func (c *DeadPropertyAssignmentEliminationTest) TestQualifiedNamePrefixAssignment() {
	c.TestSame(
		"var foo = function() {\n" +
			"  a.b.c = 20;\n" +
			"  a.b = other;\n" +
			"  a.b.c = 30;\n" +
			"}\n")

	c.TestSame(
		"var foo = function() {\n" +
			"  a.b = 20;\n" +
			"  a = other;\n" +
			"  a.b = 30;\n" +
			"}\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testCall()
func (c *DeadPropertyAssignmentEliminationTest) TestCall() {
	c.TestSame(
		"var foo = function() {\n" +
			"  a.b.c = 20;\n" +
			"  doSomething();\n" +
			"  a.b.c = 30;\n" +
			"}\n")

	c.TestSame(
		"var foo = function() {\n" +
			"  a.b.c = 20;\n" +
			"  doSomething().b.c = 25;\n" +
			"  a.b.c = 30;\n" +
			"}\n")

	c.TestSame(
		"var foo = function() { // to preserve newlines\n" +
			"  a.b.c = 20;\n" +
			"  doSomething?.();\n" +
			"  a.b.c = 30;\n" +
			"}\n")

	c.TestSame(
		"/** @constructor */\n" +
			"var foo = function() {\n" +
			"  this.c = 20;\n" +
			"  doSomething(this);\n" +
			"  this.c = 30;\n" +
			"}\n")

	c.TestSame(
		"/** @constructor */\n" +
			"var foo = function() {\n" +
			"  this.c = 20;\n" +
			"  this.doSomething();\n" +
			"  this.c = 30;\n" +
			"}\n")

	c.TestSame(
		"/** @constructor */\n" +
			"var foo = function() {\n" +
			"  this.c = 20;\n" +
			"  doSomething(this.c);\n" +
			"  this.c = 30;\n" +
			"}\n")

	c.Test(
		"var foo = function() {\n"+
			"  a.b.c = 20;\n"+
			"  doSomething(a.b.c = 25); // unconditional assignment in arg renders above\n"+
			"// assignment dead.\n"+
			"  a.b.c = 30;\n"+
			"}\n",
		"var foo = function() {\n"+
			"  20;\n"+
			"  doSomething(a.b.c = 25);\n"+
			"  a.b.c = 30;\n"+
			"}\n")

	c.TestSame(
		"var foo = function() {\n" +
			"  a.b.c = 20;\n" +
			"  doSomething?.(a.b.c = 25); // conditional assignment in optional chain arg keeps\n" +
			"// above assignment alive.\n" +
			"  a.b.c = 30;\n" +
			"}\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testYield()
func (c *DeadPropertyAssignmentEliminationTest) TestYield() {
	// Assume that properties may be read during a yield
	c.TestSame(
		"var foo = function*() {\n" +
			"  a.b.c = 20;\n" +
			"  yield;\n" +
			"  a.b.c = 30;\n" +
			"}\n")

	c.TestSame(
		"/** @constructor */\n" +
			"var foo = function*() {\n" +
			"  this.c = 20;\n" +
			"  yield;\n" +
			"  this.c = 30;\n" +
			"}\n")

	c.TestSame(
		"var obj = {\n" +
			"  *gen() {\n" +
			"    this.c = 20;\n" +
			"    yield;\n" +
			"    this.c = 30;\n" +
			"  }\n" +
			"}\n")

	c.Test(
		"var foo = function*() {\n"+
			"  a.b.c = 20;\n"+
			"  yield a.b.c = 25;\n"+
			"  a.b.c = 30;\n"+
			"}\n",
		"var foo = function*() {\n"+
			"  20;\n"+
			"  yield a.b.c = 25;\n"+
			"  a.b.c = 30;\n"+
			"}\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testNew()
func (c *DeadPropertyAssignmentEliminationTest) TestNew() {
	// Assume that properties may be read during a constructor call
	c.TestSame(
		"var foo = function() {\n" +
			"  a.b.c = 20;\n" +
			"  new C;\n" +
			"  a.b.c = 30;\n" +
			"}\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testTaggedTemplateLit()
func (c *DeadPropertyAssignmentEliminationTest) TestTaggedTemplateLit() {
	// Assume that properties may be read during a tagged template lit invocation
	c.TestSame(
		"var foo = function() {\n" +
			"  a.b.c = 20;\n" +
			"  doSomething`foo`;\n" +
			"  a.b.c = 30;\n" +
			"}\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testAwait()
func (c *DeadPropertyAssignmentEliminationTest) TestAwait() {
	// Assume that properties may be read while waiting for "await"
	c.TestSame(
		"async function foo() {\n" +
			"  a.b.c = 20;\n" +
			"  await bar;\n" +
			"  a.b.c = 30;\n" +
			"}\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testUnknownLookup()
func (c *DeadPropertyAssignmentEliminationTest) TestUnknownLookup() {
	c.TestSame(
		"/** @constructor */\n" +
			"var foo = function(str) {\n" +
			"  this.x = 5;\n" +
			"  var y = this[str];\n" +
			"  this.x = 10;\n" +
			"}\n")

	c.TestSame(
		"/** @constructor */\n" +
			"var foo = function(str) {\n" +
			"  this.x = 5;\n" +
			"  var y = this?.[str];\n" +
			"  this.x = 10;\n" +
			"}\n")

	c.TestSame(
		"/** @constructor */\n" +
			"var foo = function(x, str) {\n" +
			"  x.y = 5;\n" +
			"  var y = x[str];\n" +
			"  x.y = 10;\n" +
			"}\n")

	c.TestSame(
		"/** @constructor */\n" +
			"var foo = function(x, str) {\n" +
			"  x.y = 5;\n" +
			"  var y = x?.[str];\n" +
			"  x.y = 10;\n" +
			"}\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testName()
func (c *DeadPropertyAssignmentEliminationTest) TestName() {
	c.TestSame(
		"function f(x) {\n" +
			"  var y = { a: 0 };\n" +
			"  x.a = 123;\n" +
			"  y = x;\n" +
			"  x.a = 234;\n" +
			"  return x.a + y.a;\n" +
			"}\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testName2()
func (c *DeadPropertyAssignmentEliminationTest) TestName2() {
	c.TestSame(
		"function f(x) {\n" +
			"  var y = x;\n" +
			"  x.a = 123;\n" +
			"  x = {};\n" +
			"  x.a = 234;\n" +
			"  return x.a + y.a;\n" +
			"}\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testAliasing()
func (c *DeadPropertyAssignmentEliminationTest) TestAliasing() {
	c.TestSame(
		"function f(x) {\n" +
			"  x.b.c = 1;\n" +
			"  var y = x.a.c; // x.b.c is read here\n" +
			"  x.b.c = 2;\n" +
			"  return x.b.c + y;\n" +
			"}\n" +
			"var obj = { c: 123 };\n" +
			"f({a: obj, b: obj});\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testHook()
func (c *DeadPropertyAssignmentEliminationTest) TestHook() {
	c.TestSame(
		"function f(x, pred) {\n" +
			"  var y;\n" +
			"  x.p = 234;\n" +
			"  y = pred ? (x.p = 123) : x.p;\n" +
			"}\n")

	c.TestSame(
		"function f(x, pred) {\n" +
			"  var y;\n" +
			"  x.p = 234;\n" +
			"  y = pred ? (x.p = 123) : 123;\n" +
			"  return x.p;\n" +
			"}\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testConditionalExpression()
func (c *DeadPropertyAssignmentEliminationTest) TestConditionalExpression() {
	c.TestSame(
		"function f(x) {\n" +
			"  return (x.p = 2) || (x.p = 3); // Second assignment will never execute.\n" +
			"}\n")

	c.TestSame(
		"function f(x) {\n" +
			"  return (x.p = 0) && (x.p = 3); // Second assignment will never execute.\n" +
			"}\n")
}

// Java: DeadPropertyAssignmentEliminationTest.nullishCoalesce()
func (c *DeadPropertyAssignmentEliminationTest) TestNullishCoalesce() {
	c.TestSame(
		"function f(x) {\n" +
			"  return (x.p = 2) ?? (x.p = 3); // Second assignment will never execute.\n" +
			"}\n")
}

// Java: DeadPropertyAssignmentEliminationTest.optChain()
func (c *DeadPropertyAssignmentEliminationTest) TestOptChain() {
	c.TestSame(
		"function f(x) {\n" +
			"  return bar(x.p = 2)?.(x.p = 3); // second assignment conditionally executes, hence\n" +
			"// marks property `p` as read.\n" +
			"}\n")

	c.TestSame(
		"function foo() {\n" +
			"    this.p = 123;\n" +
			"    var z = this?.p; // marks property `p` as read\n" +
			"    this.p = 234;\n" +
			"}\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testBrackets()
func (c *DeadPropertyAssignmentEliminationTest) TestBrackets() {
	c.TestSame(
		"function f(x, p) {\n" +
			"  x.prop = 123;\n" +
			"  x[p] = 234;\n" +
			"  return x.prop;\n" +
			"}\n")
}

// This pass does not currently use the control-flow graph. Hence, dead assignments in some of the
// following tests get preserved even when they are safe to delete.
// Java: DeadPropertyAssignmentEliminationTest.testFor()
func (c *DeadPropertyAssignmentEliminationTest) TestFor() {
	c.TestSame(
		"function f(x) {\n" +
			"  x.p = 1;\n" +
			"  for(;x;) {}\n" +
			"  x.p = 2;\n" +
			"}\n")

	c.TestSame(
		"function f(x) {\n" +
			"  for(;x;) {\n" +
			"    x.p = 1;\n" +
			"  }\n" +
			"  x.p = 2;\n" +
			"}\n")

	c.TestSame(
		"function f(x) {\n" +
			"  x.p = 1;\n" +
			"  for(;;) {\n" +
			"    x.p = 2;\n" +
			"  }\n" +
			"}\n")

	c.TestSame(
		"function f(x) {\n" +
			"  x.p = 1;\n" +
			"  for(x.p = 2;;) {\n" +
			"  }\n" +
			"}\n")

	c.TestSame(
		"function f(x) {\n" +
			"  x.p = 1;\n" +
			"  for(x.p = 2;;x.p=3) {\n" +
			"    return x.p; // Reads the \"x.p = 2\" assignment.\n" +
			"  }\n" +
			"}\n")

	c.TestSame(
		"function f(x) {\n" +
			"  for(;;) {\n" +
			"    x.p = 1;\n" +
			"    x.p = 2;\n" +
			"  }\n" +
			"}\n")

	c.TestSame(
		"function f(x) {\n" +
			"  x.p = 1;\n" +
			"  for(;;) {\n" +
			"  }\n" +
			"  x.p = 2;\n" +
			"}\n")
}

// This pass does not currently use the control-flow graph. Hence, dead assignments in some of the
// following tests get preserved even when they are safe to delete.
// Java: DeadPropertyAssignmentEliminationTest.testWhile()
func (c *DeadPropertyAssignmentEliminationTest) TestWhile() {
	c.TestSame(
		"function f(x) {\n" +
			"  x.p = 1;\n" +
			"  while(x);\n" +
			"  x.p = 2;\n" +
			"}\n")

	c.TestSame(
		"function f(x) {\n" +
			"  x.p = 1;\n" +
			"  while(1) {\n" +
			"    x.p = 2;\n" +
			"  }\n" +
			"}\n")

	c.TestSame(
		"function f(x) {\n" +
			"  while(true) {\n" +
			"    x.p = 1;\n" +
			"    if (random()) continue;\n" +
			"    x.p = 2;\n" +
			"  }\n" +
			"}\n")

	c.TestSame(
		"function f(x) {\n" +
			"  while(true) {\n" +
			"    x.p = 1;\n" +
			"    if (random()) break;\n" +
			"    x.p = 2;\n" +
			"  }\n" +
			"}\n")

	c.Test(
		"function f(x) {\n"+
			"  x.p = 1;\n"+
			"  while(x.p = 2) {\n"+
			"  }\n"+
			"}\n",
		"function f(x) {\n"+
			"  1;\n"+
			"  while(x.p = 2) {\n"+
			"  }\n"+
			"}\n")

	c.Test(
		"function f(x) {\n"+
			"  while(true) {\n"+
			"    x.p = 1;\n"+
			"    x.p = 2;\n"+
			"  }\n"+
			"}\n",
		"function f(x) {\n"+
			"  while(true) {\n"+
			"    1;\n"+
			"    x.p = 2;\n"+
			"  }\n"+
			"}\n")

	c.Test(
		"function f(x) {\n"+
			"  x.p = 1;\n"+
			"  while(1) {}\n"+
			"  x.p = 2;\n"+
			"}\n",
		"function f(x) {\n"+
			"  1;\n"+
			"  while(1) {}\n"+
			"  x.p = 2;\n"+
			"}\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testTry()
func (c *DeadPropertyAssignmentEliminationTest) TestTry() {
	c.TestSame(
		"function f(x) {\n" +
			"  x.p = 1;\n" +
			"  try {\n" +
			"    x.p = 2;\n" +
			"  } catch (e) {}\n" +
			"}\n")

	c.TestSame(
		"function f(x) {\n" +
			"  x.p = 1;\n" +
			"  try {\n" +
			"  } catch (e) {\n" +
			"    x.p = 2;\n" +
			"  }\n" +
			"}\n")

	c.TestSame(
		"function f(x) {\n" +
			"  try {\n" +
			"    x.p = 1;\n" +
			"  } catch (e) {\n" +
			"    x.p = 2;\n" +
			"  }\n" +
			"}\n")

	c.Test(
		"function f(x) {\n"+
			"  try {\n"+
			"    x.p = 1;\n"+
			"    x.p = 2;\n"+
			"  } catch (e) {\n"+
			"  }\n"+
			"}\n",
		"function f(x) {\n"+
			"  try {\n"+
			"    1;\n"+
			"    x.p = 2;\n"+
			"  } catch (e) {\n"+
			"  }\n"+
			"}\n")

	c.TestSame(
		"function f(x) {\n" +
			"  try {\n" +
			"    x.p = 1;\n" +
			"    maybeThrow();\n" +
			"    x.p = 2;\n" +
			"  } catch (e) {\n" +
			"  }\n" +
			"}\n")

	c.TestSame(
		"function f(x) {\n" +
			"  try {\n" +
			"    x.p = 1;\n" +
			"    maybeThrow?.();\n" +
			"    x.p = 2;\n" +
			"  } catch (e) {\n" +
			"  }\n" +
			"}\n")

	c.TestSame(
		"/** @constructor */\n" +
			"function f() {\n" +
			"  try {\n" +
			"    this.p = 1;\n" +
			"    maybeThrow();\n" +
			"    this.p = 2;\n" +
			"  } catch (e) {\n" +
			"  }\n" +
			"}\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testThrow()
func (c *DeadPropertyAssignmentEliminationTest) TestThrow() {
	c.TestSame(
		"function f(x) {\n" +
			"  x.p = 10\n" +
			"  if (random) throw err;\n" +
			"  x.p = 20;\n" +
			"}\n")

	c.TestSame(
		"function f(x) {\n" +
			"  x.p = 10\n" +
			"  throw err;\n" +
			"  x.p = 20;\n" +
			"}\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testSwitch()
func (c *DeadPropertyAssignmentEliminationTest) TestSwitch() {
	c.TestSame(
		"function f(x, pred) {\n" +
			"  x.p = 1;\n" +
			"  switch (pred) {\n" +
			"    case 1:\n" +
			"      x.p = 2;\n" +
			"    case 2:\n" +
			"      x.p = 3;\n" +
			"      break;\n" +
			"    default:\n" +
			"      return x.p;\n" +
			"  }\n" +
			"}\n")

	c.TestSame(
		"function f(x, pred) {\n" +
			"  x.p = 1;\n" +
			"  switch (pred) {\n" +
			"    default:\n" +
			"      x.p = 2;\n" +
			"  }\n" +
			"  x.p = 3;\n" +
			"}\n")

	c.TestSame(
		"function f(x, pred) {\n" +
			"  x.p = 1;\n" +
			"  switch (pred) {\n" +
			"    default:\n" +
			"      return;\n" +
			"  }\n" +
			"  x.p = 2;\n" +
			"}\n")

	// For now we don't enter switch statements.
	c.TestSame(
		"function f(x, pred) {\n" +
			"  switch (pred) {\n" +
			"    default:\n" +
			"      x.p = 2;\n" +
			"      x.p = 3;\n" +
			"  }\n" +
			"}\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testIf()
func (c *DeadPropertyAssignmentEliminationTest) TestIf() {
	c.Test(
		"function f(x, pred) {\n"+
			"  if (pred) {\n"+
			"    x.p = 1;\n"+
			"    x.p = 2;\n"+
			"    return x.p;\n"+
			"  }\n"+
			"}\n",
		"function f(x, pred) {\n"+
			"  if (pred) {\n"+
			"    1;\n"+
			"    x.p = 2;\n"+
			"    return x.p;\n"+
			"  }\n"+
			"}\n")

	c.Test(
		"function f(x, pred) {\n"+
			"  x.p = 1;\n"+
			"  if (pred) {}\n"+
			"  x.p = 2;\n"+
			"  return x.p;\n"+
			"}\n",
		"function f(x, pred) {\n"+
			"  1;\n"+
			"  if (pred) {}\n"+
			"  x.p = 2;\n"+
			"  return x.p;\n"+
			"}\n")

	c.TestSame(
		"function f(x, pred) {\n" +
			"  if (pred) {\n" +
			"    x.p = 1;\n" +
			"  }\n" +
			"  x.p = 2;\n" +
			"}\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testCircularPropChain()
func (c *DeadPropertyAssignmentEliminationTest) TestCircularPropChain() {
	c.TestSame(
		"function f(x, y) {\n" +
			"  x.p = {};\n" +
			"  x.p.y.p.z = 10;\n" +
			"  x.p = {};\n" +
			"}\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testDifferentQualifiedNames()
func (c *DeadPropertyAssignmentEliminationTest) TestDifferentQualifiedNames() {
	c.TestSame(
		"function f(x, y) {\n" +
			"  x.p = 10;\n" +
			"  y.p = 11;\n" +
			"}\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testGetPropContainsNonQualifiedNames()
func (c *DeadPropertyAssignmentEliminationTest) TestGetPropContainsNonQualifiedNames() {
	c.TestSame(
		"function f(x) {\n" +
			"  foo(x).p = 10;\n" +
			"  foo(x).p = 11;\n" +
			"}\n")

	c.TestSame(
		"function f(x) {\n" +
			"  (x = 10).p = 10;\n" +
			"  (x = 10).p = 11;\n" +
			"}\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testEs6Constructor()
func (c *DeadPropertyAssignmentEliminationTest) TestEs6Constructor() {
	c.TestSame(
		"class Foo {\n" +
			"  constructor() {\n" +
			"    this.p = 123;\n" +
			"    var z = this.p;\n" +
			"    this.p = 234;\n" +
			"  }\n" +
			"}\n")

	c.TestSame(
		"class Foo {\n" +
			"  constructor() {\n" +
			"    this.p = 123;\n" +
			"    var z = this?.p;\n" +
			"    this.p = 234;\n" +
			"  }\n" +
			"}\n")

	c.Test(
		"class Foo {\n"+
			"  constructor() {\n"+
			"    this.p = 123;\n"+
			"    this.p = 234;\n"+
			"  }\n"+
			"}\n",
		"class Foo {\n"+
			"  constructor() {\n"+
			"    123;\n"+
			"    this.p = 234;\n"+
			"  }\n"+
			"}\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testES6ClassExtends()
func (c *DeadPropertyAssignmentEliminationTest) TestES6ClassExtends() {
	c.TestSame(
		"class C {\n" +
			"  constructor() {\n" +
			"    this.x = 20;\n" +
			"  }\n" +
			"}\n" +
			"class D extends C {\n" +
			"  constructor() {\n" +
			"    super();\n" +
			"    this.x = 40;\n" +
			"  }\n" +
			"}\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testStaticBlock()
func (c *DeadPropertyAssignmentEliminationTest) TestStaticBlock() {
	c.TestSame(
		"class Foo {\n" +
			"  static {\n" +
			"    this.p = 123;\n" +
			"    var z = this.p;\n" +
			"    this.p = 234;\n" +
			"  }\n" +
			"}\n")

	c.TestSame(
		"class C {\n" +
			"  static {\n" +
			"    this.x = 20;\n" +
			"  }\n" +
			"}\n" +
			"class D extends C {\n" +
			"  static {\n" +
			"    this.x = 40;\n" +
			"  }\n" +
			"}\n")

	// TODO(b/235871861): Ideally the first assignment would be removed
	c.TestSame(
		"class Foo {\n" +
			"  static {\n" +
			"    this.p = 123; // should get removed\n" +
			"    this.p = 234;\n" +
			"  }\n" +
			"}\n")

	// TODO(b/235871861): Ideally the first assignment would be removed
	c.TestSame(
		"class Foo {\n" +
			"  static {\n" +
			"    let o = { prop : 2};\n" +
			"    o.prop = 123; // should get removed\n" +
			"    o.prop = 234;\n" +
			"  }\n" +
			"}\n")

	c.TestSame(
		"let o = { prop : 2};\n" +
			"class C {\n" +
			"  static {\n" +
			"    o.prop = 20; // should not get removed\n" +
			"  }\n" +
			"}\n" +
			"o.prop;\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testGetter()
func (c *DeadPropertyAssignmentEliminationTest) TestGetter() {
	c.TestSame(
		"/** @constructor */ function Foo() { this.enabled = false; };\n" +
			"Object.defineProperties(Foo.prototype, {bar: {\n" +
			"  get: function () { return this.enabled ? 'enabled' : 'disabled'; }\n" +
			"}});\n" +
			"function f() {\n" +
			"  var foo = new Foo()\n" +
			"  foo.enabled = true;\n" +
			"  var f = foo.bar;\n" +
			"  foo.enabled = false;\n" +
			"}\n")

	c.TestSame(
		"/** @constructor */ function Foo() { this.enabled = false; };\n" +
			"Object.defineProperties(Foo.prototype, {bar: {\n" +
			"  get: function () { return this.enabled ? 'enabled' : 'disabled'; }\n" +
			"}});\n" +
			"function f() {\n" +
			"  var foo = new Foo()\n" +
			"  foo.enabled = true;\n" +
			"  var f = foo?.bar;\n" +
			"  foo.enabled = false;\n" +
			"}\n")

	c.TestSame(
		"/** @constructor */ function Foo() { this.enabled = false; };\n" +
			"Object.defineProperty(Foo, 'bar', {\n" +
			"  get: function () { return this.enabled ? 'enabled' : 'disabled'; }\n" +
			"});\n" +
			"function f() {\n" +
			"  var foo = new Foo()\n" +
			"  foo.enabled = true;\n" +
			"  var f = foo.bar;\n" +
			"  foo.enabled = false;\n" +
			"}\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testGetter_afterDeadAssignment()
func (c *DeadPropertyAssignmentEliminationTest) TestGetter_afterDeadAssignment() {
	c.TestSame(
		"function f() {\n" +
			"  var foo = new Foo()\n" +
			"  foo.enabled = true;\n" +
			"  var f = foo.bar;\n" +
			"  foo.enabled = false;\n" +
			"}\n" +
			"/** @constructor */ function Foo() { this.enabled = false; };\n" +
			"Object.defineProperties(Foo.prototype, {bar: {\n" +
			"  get: function () { return this.enabled ? 'enabled' : 'disabled'; }\n" +
			"}});\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testGetter_onDifferentType()
func (c *DeadPropertyAssignmentEliminationTest) TestGetter_onDifferentType() {
	c.TestSame(
		"/** @constructor */\n" +
			"function Foo() {\n" +
			"  this.enabled = false;\n" +
			"};\n" +
			"Object.defineProperties(\n" +
			"    Foo.prototype, {\n" +
			"      baz: {\n" +
			"        get: function () { return this.enabled ? 'enabled' : 'disabled'; }\n" +
			"      }\n" +
			"    });\n" +
			"/** @constructor */\n" +
			"function Bar() {\n" +
			"  this.enabled = false;\n" +
			"  this.baz = 123;\n" +
			"};\n" +
			"function f() {\n" +
			"  var bar = new Bar();\n" +
			"  bar.enabled = true;\n" +
			"  var ret = bar.baz;\n" +
			"  bar.enabled = false;\n" +
			"  return ret;\n" +
			"};\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testSetter()
func (c *DeadPropertyAssignmentEliminationTest) TestSetter() {
	c.TestSame(
		"/** @constructor */ function Foo() { this.enabled = false; this.x = null; };\n" +
			"Object.defineProperties(Foo.prototype, {bar: {\n" +
			"  set: function (x) { this.x = this.enabled ? x * 2 : x; }\n" +
			"}});\n" +
			"function f() {\n" +
			"  var foo = new Foo()\n" +
			"  foo.enabled = true;\n" +
			"  foo.bar = 10;\n" +
			"  foo.enabled = false;\n" +
			"}\n")

	c.TestSame(
		"/** @constructor */ function Foo() { this.enabled = false; this.x = null; };\n" +
			"Object.defineProperty(Foo, 'bar', {\n" +
			"  set: function (x) { this.x = this.enabled ? x * 2 : x; }\n" +
			"});\n" +
			"function f() {\n" +
			"  var foo = new Foo()\n" +
			"  foo.enabled = true;\n" +
			"  foo.bar = 10;\n" +
			"  foo.enabled = false;\n" +
			"}\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testEs5Getter()
func (c *DeadPropertyAssignmentEliminationTest) TestEs5Getter() {
	c.TestSame(
		"var bar = {\n" +
			"  enabled: false,\n" +
			"  get baz() {\n" +
			"    return this.enabled ? 'enabled' : 'disabled';\n" +
			"  }\n" +
			"};\n" +
			"function f() {\n" +
			"  bar.enabled = true;\n" +
			"  var ret = bar.baz;\n" +
			"  bar.enabled = false;\n" +
			"  return ret;\n" +
			"};\n")

	c.TestSame(
		"var bar = {\n" +
			"  enabled: false,\n" +
			"  get baz() {\n" +
			"    return this.enabled ? 'enabled' : 'disabled';\n" +
			"  }\n" +
			"};\n" +
			"function f() {\n" +
			"  bar.enabled = true;\n" +
			"  var ret = bar?.baz;\n" +
			"  bar.enabled = false;\n" +
			"  return ret;\n" +
			"};\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testEs5Getter_computed()
func (c *DeadPropertyAssignmentEliminationTest) TestEs5Getter_computed() {
	c.Test(
		"var bar = {\n"+
			"  enabled: false,\n"+
			"  get ['baz']() {\n"+
			"    return this.enabled ? 'enabled' : 'disabled';\n"+
			"  }\n"+
			"};\n"+
			"function f() {\n"+
			"  bar.enabled = true;\n"+
			"  var ret = bar.baz;\n"+
			"  bar.enabled = false;\n"+
			"  return ret;\n"+
			"};\n",
		"var bar = {\n"+
			"  enabled: false,\n"+
			"  get ['baz']() {\n"+
			"    return this.enabled ? 'enabled' : 'disabled';\n"+
			"  }\n"+
			"};\n"+
			"function f() {\n"+
			"  true;\n"+
			"  var ret = bar.baz;\n"+
			"  bar.enabled = false;\n"+
			"  return ret;\n"+
			"};\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testEs5Setter()
func (c *DeadPropertyAssignmentEliminationTest) TestEs5Setter() {
	c.TestSame(
		"var bar = {\n" +
			"  enabled: false,\n" +
			"  set baz(x) {\n" +
			"    this.x = this.enabled ? x * 2 : x;\n" +
			"  }\n" +
			"};\n" +
			"function f() {\n" +
			"  bar.enabled = true;\n" +
			"  bar.baz = 10;\n" +
			"  bar.enabled = false;\n" +
			"};\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testEs5Setter_computed()
func (c *DeadPropertyAssignmentEliminationTest) TestEs5Setter_computed() {
	c.Test(
		"var bar = {\n"+
			"  enabled: false,\n"+
			"  set ['baz'](x) {\n"+
			"    this.x = this.enabled ? x * 2 : x;\n"+
			"  }\n"+
			"};\n"+
			"function f() {\n"+
			"  bar.enabled = true;\n"+
			"  bar.baz = 10;\n"+
			"  bar.enabled = false;\n"+
			"};\n",
		"var bar = {\n"+
			"  enabled: false,\n"+
			"  set ['baz'](x) {\n"+
			"    this.x = this.enabled ? x * 2 : x;\n"+
			"  }\n"+
			"};\n"+
			"function f() {\n"+
			"  true;\n"+
			"  bar.baz = 10;\n"+
			"  bar.enabled = false;\n"+
			"};\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testObjectDefineProperty_aliasedParams()
func (c *DeadPropertyAssignmentEliminationTest) TestObjectDefineProperty_aliasedParams() {
	c.Test(
		"function addGetter(obj, propName) {\n"+
			"  Object.defineProperty(obj, propName, {\n"+
			"    get: function() { return this[propName]; }\n"+
			"  });\n"+
			"};\n"+
			"/** @constructor */ function Foo() { this.enabled = false; this.x = null; };\n"+
			"addGetter(Foo.prototype, 'x');\n"+
			"function f() {\n"+
			"  var foo = new Foo()\n"+
			"  foo.enabled = true;\n"+
			"  foo.bar = 10;\n"+
			"  foo.enabled = false;\n"+
			"}\n"+
			"function z() {\n"+
			"  var x = {};\n"+
			"  x.bar = 10;\n"+
			"  x.bar = 20;\n"+
			"}\n",
		"function addGetter(obj, propName) {\n"+
			"  Object.defineProperty(obj, propName, {\n"+
			"    get: function() { return this[propName]; }\n"+
			"  });\n"+
			"};\n"+
			"/** @constructor */ function Foo() { this.enabled = false; this.x = null; };\n"+
			"addGetter(Foo.prototype, 'x');\n"+
			"function f() {\n"+
			"  var foo = new Foo()\n"+
			"  true;\n"+
			"  foo.bar = 10;\n"+
			"  foo.enabled = false;\n"+
			"}\n"+
			"function z() {\n"+
			"  var x = {};\n"+
			"  10;\n"+
			"  x.bar = 20;\n"+
			"}\n")

	c.Test(
		"function f() {\n"+
			"  var foo = new Foo()\n"+
			"  foo.enabled = true;\n"+
			"  foo.bar = 10;\n"+
			"  foo.enabled = false;\n"+
			"}\n"+
			"function addGetter(obj, propName) {\n"+
			"  Object.defineProperty(obj, propName, {\n"+
			"    get: function() { return this[propName]; }\n"+
			"  });\n"+
			"};\n"+
			"/** @constructor */ function Foo() { this.enabled = false; this.x = null; };\n"+
			"addGetter(Foo.prototype, 'x');\n"+
			"function z() {\n"+
			"  var x = {};\n"+
			"  x.bar = 10;\n"+
			"  x.bar = 20;\n"+
			"}\n",
		"function f() {\n"+
			"  var foo = new Foo()\n"+
			"  true;\n"+
			"  foo.bar = 10;\n"+
			"  foo.enabled = false;\n"+
			"}\n"+
			"function addGetter(obj, propName) {\n"+
			"  Object.defineProperty(obj, propName, {\n"+
			"    get: function() { return this[propName]; }\n"+
			"  });\n"+
			"};\n"+
			"/** @constructor */ function Foo() { this.enabled = false; this.x = null; };\n"+
			"addGetter(Foo.prototype, 'x');\n"+
			"function z() {\n"+
			"  var x = {};\n"+
			"  10;\n"+
			"  x.bar = 20;\n"+
			"}\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testObjectDefineProperty_aliasedObject()
func (c *DeadPropertyAssignmentEliminationTest) TestObjectDefineProperty_aliasedObject() {
	c.TestSame(
		"/** @constructor */ function Foo() { this.enabled = false; this.x = null; };\n" +
			"var x = Foo.prototype;\n" +
			"Object.defineProperty(x, 'bar', {\n" +
			"  set: function (x) { this.x = this.enabled ? x * 2 : x; }\n" +
			"});\n" +
			"function f() {\n" +
			"  var foo = new Foo()\n" +
			"  foo.enabled = true;\n" +
			"  foo.bar = 10;\n" +
			"  foo.enabled = false;\n" +
			"}\n" +
			"function z() {\n" +
			"  var x = {};\n" +
			"  x.bar = 10;\n" +
			"  x.bar = 20;\n" +
			"}\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testObjectDefineProperty_aliasedPropName()
func (c *DeadPropertyAssignmentEliminationTest) TestObjectDefineProperty_aliasedPropName() {
	c.Test(
		"/** @constructor */ function Foo() { this.enabled = false; this.x = null; };\n"+
			"var x = 'bar';\n"+
			"Object.defineProperty(Foo.prototype, x, {\n"+
			"  set: function (x) { this.x = this.enabled ? x * 2 : x; }\n"+
			"});\n"+
			"function f() {\n"+
			"  var foo = new Foo()\n"+
			"  foo.enabled = true;\n"+
			"  foo.bar = 10;\n"+
			"  foo.enabled = false;\n"+
			"}\n"+
			"function z() {\n"+
			"  var x = {};\n"+
			"  x.bar = 10;\n"+
			"  x.bar = 20;\n"+
			"}\n",
		"/** @constructor */ function Foo() { this.enabled = false; this.x = null; };\n"+
			"var x = 'bar';\n"+
			"Object.defineProperty(Foo.prototype, x, {\n"+
			"  set: function (x) { this.x = this.enabled ? x * 2 : x; }\n"+
			"});\n"+
			"function f() {\n"+
			"  var foo = new Foo()\n"+
			"  true;\n"+
			"  foo.bar = 10;\n"+
			"  foo.enabled = false;\n"+
			"}\n"+
			"function z() {\n"+
			"  var x = {};\n"+
			"  10;\n"+
			"  x.bar = 20;\n"+
			"}\n")

	c.Test(
		"/** @constructor */ function Foo() { this.enabled = false; this.x = null; };\n"+
			"var x = 'bar';\n"+
			"Object.defineProperty(Foo.prototype, x, {\n"+
			"  value: 10\n"+
			"});\n"+
			"function f() {\n"+
			"  var foo = new Foo()\n"+
			"  foo.enabled = true;\n"+
			"  foo.bar = 10;\n"+
			"  foo.enabled = false;\n"+
			"}\n"+
			"function z() {\n"+
			"  var x = {};\n"+
			"  x.bar = 10;\n"+
			"  x.bar = 20;\n"+
			"}\n",
		"/** @constructor */ function Foo() { this.enabled = false; this.x = null; };\n"+
			"var x = 'bar';\n"+
			"Object.defineProperty(Foo.prototype, x, {\n"+
			"  value: 10\n"+
			"});\n"+
			"function f() {\n"+
			"  var foo = new Foo()\n"+
			"  true;\n"+
			"  foo.bar = 10;\n"+
			"  foo.enabled = false;\n"+
			"}\n"+
			"function z() {\n"+
			"  var x = {};\n"+
			"  10;\n"+
			"  x.bar = 20;\n"+
			"}\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testObjectDefineProperty_aliasedPropSet()
func (c *DeadPropertyAssignmentEliminationTest) TestObjectDefineProperty_aliasedPropSet() {
	c.Test(
		"/** @constructor */ function Foo() { this.enabled = false; this.x = null; };\n"+
			"var x = {\n"+
			"  set: function (x) { this.x = this.enabled ? x * 2 : x; }\n"+
			"};\n"+
			"Object.defineProperty(Foo.prototype, 'bar', x);\n"+
			"function f() {\n"+
			"  var foo = new Foo()\n"+
			"  foo.enabled = true;\n"+
			"  foo.bar = 10;\n"+
			"  foo.enabled = false;\n"+
			"}\n"+
			"function z() {\n"+
			"  var x = {};\n"+
			"  x.bar = 10;\n"+
			"  x.bar = 20;\n"+
			"}\n",
		"/** @constructor */ function Foo() { this.enabled = false; this.x = null; };\n"+
			"var x = {\n"+
			"  set: function (x) { this.x = this.enabled ? x * 2 : x; }\n"+
			"};\n"+
			"Object.defineProperty(Foo.prototype, 'bar', x);\n"+
			"function f() {\n"+
			"  var foo = new Foo()\n"+
			"  true;\n"+
			"  foo.bar = 10;\n"+
			"  foo.enabled = false;\n"+
			"}\n"+
			"function z() {\n"+
			"  var x = {};\n"+
			"  10;\n"+
			"  x.bar = 20;\n"+
			"}\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testObjectDefineProperties_aliasedPropertyMap()
func (c *DeadPropertyAssignmentEliminationTest) TestObjectDefineProperties_aliasedPropertyMap() {
	c.Test(
		"/** @constructor */ function Foo() { this.enabled = false; this.x = null; };\n"+
			"var properties = {bar: {\n"+
			"  set: function (x) { this.x = this.enabled ? x * 2 : x; }\n"+
			"}};\n"+
			"Object.defineProperties(Foo.prototype, properties);\n"+
			"function f() {\n"+
			"  var foo = new Foo()\n"+
			"  foo.enabled = true;\n"+
			"  foo.bar = 10;\n"+
			"  foo.enabled = false;\n"+
			"}\n"+
			"function z() {\n"+
			"  var x = {};\n"+
			"  x.bar = 10;\n"+
			"  x.bar = 20;\n"+
			"}\n",
		"/** @constructor */ function Foo() { this.enabled = false; this.x = null; };\n"+
			"var properties = {bar: {\n"+
			"  set: function (x) { this.x = this.enabled ? x * 2 : x; }\n"+
			"}};\n"+
			"Object.defineProperties(Foo.prototype, properties);\n"+
			"function f() {\n"+
			"  var foo = new Foo()\n"+
			"  true;\n"+
			"  foo.bar = 10;\n"+
			"  foo.enabled = false;\n"+
			"}\n"+
			"function z() {\n"+
			"  var x = {};\n"+
			"  10;\n"+
			"  x.bar = 20;\n"+
			"}\n")

	c.Test(
		"/** @constructor */ function Foo() { this.enabled = false; this.x = null; };\n"+
			"var properties = {\n"+
			"  set: function (x) { this.x = this.enabled ? x * 2 : x; }\n"+
			"};\n"+
			"Object.defineProperties(Foo.prototype, {bar: properties});\n"+
			"function f() {\n"+
			"  var foo = new Foo()\n"+
			"  foo.enabled = true;\n"+
			"  foo.bar = 10;\n"+
			"  foo.enabled = false;\n"+
			"}\n"+
			"function z() {\n"+
			"  var x = {};\n"+
			"  x.bar = 10;\n"+
			"  x.bar = 20;\n"+
			"}\n",
		"/** @constructor */ function Foo() { this.enabled = false; this.x = null; };\n"+
			"var properties = {\n"+
			"  set: function (x) { this.x = this.enabled ? x * 2 : x; }\n"+
			"};\n"+
			"Object.defineProperties(Foo.prototype, {bar: properties});\n"+
			"function f() {\n"+
			"  var foo = new Foo()\n"+
			"  true;\n"+
			"  foo.bar = 10;\n"+
			"  foo.enabled = false;\n"+
			"}\n"+
			"function z() {\n"+
			"  var x = {};\n"+
			"  10;\n"+
			"  x.bar = 20;\n"+
			"}\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testObjectDefineProperties_aliasedObject()
func (c *DeadPropertyAssignmentEliminationTest) TestObjectDefineProperties_aliasedObject() {
	c.Test(
		"/** @constructor */ function Foo() { this.enabled = false; this.x = null; };\n"+
			"var properties = {bar: {\n"+
			"  set: function (x) { this.x = this.enabled ? x * 2 : x; }\n"+
			"}};\n"+
			"var x = Foo.prototype;\n"+
			"Object.defineProperties(x, {bar: {\n"+
			"  set: function (x) { this.x = this.enabled ? x * 2 : x; }\n"+
			"}});\n"+
			"function f() {\n"+
			"  var foo = new Foo()\n"+
			"  foo.enabled = true;\n"+
			"  foo.bar = 10;\n"+
			"  foo.enabled = false;\n"+
			"  foo.enabled = true;\n"+
			"}\n"+
			"function z() {\n"+
			"  var x = {};\n"+
			"  x.bar = 10;\n"+
			"  x.bar = 20;\n"+
			"}\n",
		"/** @constructor */ function Foo() { this.enabled = false; this.x = null; };\n"+
			"var properties = {bar: {\n"+
			"  set: function (x) { this.x = this.enabled ? x * 2 : x; }\n"+
			"}};\n"+
			"var x = Foo.prototype;\n"+
			"Object.defineProperties(x, {bar: {\n"+
			"  set: function (x) { this.x = this.enabled ? x * 2 : x; }\n"+
			"}});\n"+
			"function f() {\n"+
			"  var foo = new Foo()\n"+
			"  foo.enabled = true;\n"+
			"  foo.bar = 10;\n"+
			"  false;\n"+
			"  foo.enabled = true;\n"+
			"}\n"+
			"function z() {\n"+
			"  var x = {};\n"+
			"  x.bar = 10;\n"+
			"  x.bar = 20;\n"+
			"}\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testPropertyDefinedInExterns()
func (c *DeadPropertyAssignmentEliminationTest) TestPropertyDefinedInExterns() {
	externs :=
		"var window = {};\n" +
			"/** @type {number} */ window.innerWidth\n" +
			"/** @constructor */\n" +
			"var Image = function() {};\n" +
			"/** @type {string} */ Image.prototype.src;\n"

	c.TestSame(
		c.Externs(externs),
		c.Srcs(
			"function z() {\n"+
				"  window.innerWidth = 10;\n"+
				"  window.innerWidth = 20;\n"+
				"}\n"))

	c.TestSame(
		c.Externs(externs),
		c.Srcs(
			"function z() {\n"+
				"  var img = new Image();\n"+
				"  img.src = '';\n"+
				"  img.src = 'foo.bar';\n"+
				"}\n"))

	c.TestSame(
		c.Externs(externs),
		c.Srcs(
			"function z(x) {\n"+
				"  x.src = '';\n"+
				"  x.src = 'foo.bar';\n"+
				"}\n"))
}

// Java: DeadPropertyAssignmentEliminationTest.testJscompInherits()
func (c *DeadPropertyAssignmentEliminationTest) TestJscompInherits() {
	c.Test(
		"/** @constructor */ function Foo() { this.bar = null; };\n"+
			"var $jscomp = {};\n"+
			"$jscomp.inherits = function(x) {\n"+
			"  Object.defineProperty(x, x, x);\n"+
			"};\n"+
			"function f() {\n"+
			"  var foo = new Foo()\n"+
			"  foo.bar = 10;\n"+
			"  foo.bar = 20;\n"+
			"}\n",
		"/** @constructor */ function Foo() { this.bar = null; };\n"+
			"var $jscomp = {};\n"+
			"$jscomp.inherits = function(x) {\n"+
			"  Object.defineProperty(x, x, x);\n"+
			"};\n"+
			"function f() {\n"+
			"  var foo = new Foo()\n"+
			"  10;\n"+
			"  foo.bar = 20;\n"+
			"}\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testGithubIssue2874()
func (c *DeadPropertyAssignmentEliminationTest) TestGithubIssue2874() {
	c.TestSame(
		"var globalObj = {i0:0};\n" +
			"\n" +
			"function func(b) {\n" +
			"  var g = globalObj;\n" +
			"  var f = b;\n" +
			"  g.i0 = f.i0;\n" +
			"  g = b;\n" +
			"  g.i0 = 0;\n" +
			"}\n" +
			"func({i0:2});\n" +
			"alert(globalObj);\n")
}

// Java: DeadPropertyAssignmentEliminationTest.testReplaceShorthandAssignmentOpWithRegularOp()
func (c *DeadPropertyAssignmentEliminationTest) TestReplaceShorthandAssignmentOpWithRegularOp() {
	// See https://github.com/google/closure-compiler/issues/3017
	c.Test(
		"function f(obj) { // preserve newlines\n"+
			"  obj.a = (obj.a |= 2) | 8;\n"+
			"  obj.a = (obj.a |= 16) | 32;\n"+
			"}\n",
		"function f(obj) { // preserve newlines\n"+
			"  obj.a = (obj.a | 2) | 8;\n"+
			"  obj.a = (obj.a | 16) | 32;\n"+
			"}\n")
}

// Java: DeadPropertyAssignmentEliminationTest.getProcessor(Compiler)
func (c *DeadPropertyAssignmentEliminationTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewDeadPropertyAssignmentElimination(compiler)
}
