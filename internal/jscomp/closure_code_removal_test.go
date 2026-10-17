// Port of com/google/javascript/jscomp/ClosureCodeRemovalTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/ClosureCodeRemovalTest.java)

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

// Tests for ClosureCodeRemoval
type ClosureCodeRemovalTest struct {
	CompilerTestCaseBase
}

const closureCodeRemovalTest_EXTERNS = "var window;"

const closureCodeRemovalTest_ASSERTIONS = `const asserts = {};
/** @closurePrimitive {asserts.truthy} */
asserts.assert = function(...args) {};

/** @closurePrimitive {asserts.fail} */
asserts.fail = function(...args) {};
`

// Java: ClosureCodeRemovalTest()
func NewClosureCodeRemovalTest() *ClosureCodeRemovalTest {
	c := &ClosureCodeRemovalTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, closureCodeRemovalTest_EXTERNS)
	return c
}

func TestClosureCodeRemovalTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewClosureCodeRemovalTest() })
}

// Java: ClosureCodeRemovalTest.getProcessor(Compiler)
func (c *ClosureCodeRemovalTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewClosureCodeRemoval(compiler, true, true)
}

// Java: ClosureCodeRemovalTest.testRemoveAbstract()
func (c *ClosureCodeRemovalTest) TestRemoveAbstract() {
	c.Test(
		"function Foo() {}; Foo.prototype.doSomething = goog.abstractMethod;",
		"function Foo() {};")
}

// Java: ClosureCodeRemovalTest.testRemoveMultiplySetAbstract()
func (c *ClosureCodeRemovalTest) TestRemoveMultiplySetAbstract() {
	c.Test(
		`function Foo() {}; Foo.prototype.doSomething =
Foo.prototype.doSomethingElse = Foo.prototype.oneMore =
goog.abstractMethod;
`,
		"function Foo() {};")
}

// Java: ClosureCodeRemovalTest.testDoNotRemoveNormal()
func (c *ClosureCodeRemovalTest) TestDoNotRemoveNormal() {
	c.TestSame("function Foo() {}; Foo.prototype.doSomething = function() {};")
}

// Java: ClosureCodeRemovalTest.testDoNotRemoveOverride()
func (c *ClosureCodeRemovalTest) TestDoNotRemoveOverride() {
	c.Test(
		`function Foo() {}; Foo.prototype.doSomething = goog.abstractMethod;
function Bar() {}; goog.inherits(Bar, Foo);
Bar.prototype.doSomething = function() {}
`,
		`function Foo() {}; function Bar() {}; goog.inherits(Bar, Foo);
Bar.prototype.doSomething = function() {}
`)
}

// Java: ClosureCodeRemovalTest.testDoNotRemoveNonQualifiedName()
func (c *ClosureCodeRemovalTest) TestDoNotRemoveNonQualifiedName() {
	c.TestSame("document.getElementById('x').y = goog.abstractMethod;")
}

// Java: ClosureCodeRemovalTest.testStopRemovalAtNonQualifiedName()
func (c *ClosureCodeRemovalTest) TestStopRemovalAtNonQualifiedName() {
	c.Test(
		`function Foo() {}; function Bar() {};
Foo.prototype.x = document.getElementById('x').y = Bar.prototype.x
 = goog.abstractMethod;
`,
		`function Foo() {}; function Bar() {};
Foo.prototype.x = document.getElementById('x').y =
goog.abstractMethod;
`)
}

// Java: ClosureCodeRemovalTest.testRemoveAbstractAssignmentOfEmptyFunction()
func (c *ClosureCodeRemovalTest) TestRemoveAbstractAssignmentOfEmptyFunction() {
	c.Test(
		`function Foo() {};
/** @abstract */
Foo.prototype.doSomething = function() {};
`,
		"function Foo() {};")
}

// Java: ClosureCodeRemovalTest.testDoNotRemoveAbstractAssignmentOfFunctionCall()
func (c *ClosureCodeRemovalTest) TestDoNotRemoveAbstractAssignmentOfFunctionCall() {
	c.TestSame(
		`function Foo() {};
/** @abstract */
Foo.prototype.doSomething = (function() { /* return something fancy */ })();
`)
}

// Java: ClosureCodeRemovalTest.testRemoveAbstract_annotation_es6()
func (c *ClosureCodeRemovalTest) TestRemoveAbstract_annotation_es6() {
	c.Test(
		`/** @abstract */
class Foo {
  /** @abstract */
  doSomething() {}
}
`,
		"/** @abstract */ class Foo {}")
}

// Java: ClosureCodeRemovalTest.testDoNotRemoveNormal_es6()
func (c *ClosureCodeRemovalTest) TestDoNotRemoveNormal_es6() {
	c.TestSame(
		`/** @abstract */
class Foo {
  doSomething() {}
}
`)
}

// Java: ClosureCodeRemovalTest.testAssertionRemoval1()
func (c *ClosureCodeRemovalTest) TestAssertionRemoval1() {
	c.Test("var x = goog.asserts.assert(y(), 'message');", "var x = y();")
}

// Java: ClosureCodeRemovalTest.testAssertionRemoval2()
func (c *ClosureCodeRemovalTest) TestAssertionRemoval2() {
	c.Test("goog.asserts.assert(y(), 'message');", "")
}

// Java: ClosureCodeRemovalTest.testAssertionRemoval3()
func (c *ClosureCodeRemovalTest) TestAssertionRemoval3() {
	c.Test("goog.asserts.assert();", "")
}

// Java: ClosureCodeRemovalTest.testAssertionRemoval4()
func (c *ClosureCodeRemovalTest) TestAssertionRemoval4() {
	c.Test("var x = goog.asserts.assert();", "var x = void 0;")
}

// Java: ClosureCodeRemovalTest.testClosurePrimitiveAssertionRemoval1()
func (c *ClosureCodeRemovalTest) TestClosurePrimitiveAssertionRemoval1() {
	c.EnableTypeCheck()
	c.Test(closureCodeRemovalTest_ASSERTIONS+"var x = asserts.assert(y(), 'message');", closureCodeRemovalTest_ASSERTIONS+"var x = y();")
}

// Java: ClosureCodeRemovalTest.testClosurePrimitiveAssertionRemoval2()
func (c *ClosureCodeRemovalTest) TestClosurePrimitiveAssertionRemoval2() {
	c.EnableTypeCheck()
	c.Test(closureCodeRemovalTest_ASSERTIONS+"asserts.assert(y(), 'message');", closureCodeRemovalTest_ASSERTIONS)
}

// Java: ClosureCodeRemovalTest.testClosurePrimitiveAssertionRemoval3()
func (c *ClosureCodeRemovalTest) TestClosurePrimitiveAssertionRemoval3() {
	c.EnableTypeCheck()
	c.Test(closureCodeRemovalTest_ASSERTIONS+"asserts.assert();", closureCodeRemovalTest_ASSERTIONS)
}

// Java: ClosureCodeRemovalTest.testClosurePrimitiveAssertionRemoval4()
func (c *ClosureCodeRemovalTest) TestClosurePrimitiveAssertionRemoval4() {
	c.EnableTypeCheck()
	c.Test(closureCodeRemovalTest_ASSERTIONS+"var x = asserts.assert();", closureCodeRemovalTest_ASSERTIONS+"var x = void 0;")
}

// Java: ClosureCodeRemovalTest.testClosurePrimitiveAssertionRemoval_keepsAssertsFail()
func (c *ClosureCodeRemovalTest) TestClosurePrimitiveAssertionRemoval_keepsAssertsFail() {
	c.EnableTypeCheck()
	c.TestSame(closureCodeRemovalTest_ASSERTIONS + "var x = asserts.fail();")
}

// Java: ClosureCodeRemovalTest.testClosurePrimitiveAssertionRemoval_worksWithColors()
func (c *ClosureCodeRemovalTest) TestClosurePrimitiveAssertionRemoval_worksWithColors() {
	c.EnableTypeCheck()
	c.ReplaceTypesWithColors()
	c.DisableCompareJsDoc()

	c.Test(closureCodeRemovalTest_ASSERTIONS+"var x = asserts.assert(y(), 'message');", closureCodeRemovalTest_ASSERTIONS+"var x = y();")
}

// Java: ClosureCodeRemovalTest.testClosurePrimitiveAssertionRemoval_keepsAssertsFailWithColors()
func (c *ClosureCodeRemovalTest) TestClosurePrimitiveAssertionRemoval_keepsAssertsFailWithColors() {
	c.EnableTypeCheck()
	c.ReplaceTypesWithColors()
	c.DisableCompareJsDoc()

	c.TestSame(closureCodeRemovalTest_ASSERTIONS + "var x = asserts.fail();")
}

// Java: ClosureCodeRemovalTest.testDoNotRemoveAbstractClass()
func (c *ClosureCodeRemovalTest) TestDoNotRemoveAbstractClass() {
	c.TestSame(
		`var ns = {};
/** @abstract */
ns.A = class {};
ns.B = class extends ns.A {}
`)
}
