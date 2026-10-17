// Port of com/google/javascript/jscomp/ExportTestFunctionsTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/ExportTestFunctionsTest.java)

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
	"testing"
)

// ExportTestFunctionsTest holds tests for ExportTestFunctions.
type ExportTestFunctionsTest struct {
	CompilerTestCaseBase
}

const exportTestFunctionsTest_EXTERNS = "function google_exportSymbol(a, b) {}; function google_exportProperty(a, b, c) {};"

const exportTestFunctionsTest_TEST_FUNCTIONS_WITH_NAMES = `function Foo(arg) {};
function setUp(arg3) {};
function tearDown(arg, arg2) {};
function testBar(arg) {};
function test$(arg) {};
function test$foo(arg) {}
`

// Java: ExportTestFunctionsTest()
func NewExportTestFunctionsTest() *ExportTestFunctionsTest {
	c := &ExportTestFunctionsTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, exportTestFunctionsTest_EXTERNS)
	return c
}

func TestExportTestFunctionsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewExportTestFunctionsTest() })
}

// Java: ExportTestFunctionsTest.getProcessor(Compiler)
func (c *ExportTestFunctionsTest) GetProcessor(compiler *Compiler) CompilerPass {
	exportSymbolFunction := "google_exportSymbol"
	exportPropertyFunction := "google_exportProperty"
	return NewExportTestFunctions(compiler, &exportSymbolFunction, &exportPropertyFunction)
}

// Java: ExportTestFunctionsTest.testFunctionsAreExported()
func (c *ExportTestFunctionsTest) TestFunctionsAreExported() {
	c.Test(
		exportTestFunctionsTest_TEST_FUNCTIONS_WITH_NAMES,
		`function Foo(arg){};
function setUp(arg3){} google_exportSymbol("setUp",setUp);;
function tearDown(arg,arg2) {}
google_exportSymbol("tearDown",tearDown);;
function testBar(arg){} google_exportSymbol("testBar",testBar);;
function test$(arg){} google_exportSymbol("test$",test$);;
function test$foo(arg){} google_exportSymbol("test$foo",test$foo)
`)
}

// Helper functions
//
// Java: ExportTestFunctionsTest.testBasicTestFunctionsAreExported()
func (c *ExportTestFunctionsTest) TestBasicTestFunctionsAreExported() {
	c.TestSame("function Foo() {function testA(){}}")
	c.Test("function setUp() {}", "function setUp(){} google_exportSymbol('setUp',setUp)")
	c.Test(
		"function setUpPage() {}",
		"function setUpPage(){} google_exportSymbol('setUpPage',setUpPage)")
	c.Test(
		"function shouldRunTests() {}",
		`function shouldRunTests(){}
google_exportSymbol('shouldRunTests',shouldRunTests)
`)
	c.Test(
		"function tearDown() {}", "function tearDown(){} google_exportSymbol('tearDown',tearDown)")
	c.Test(
		"function tearDownPage() {}",
		`function tearDownPage(){} google_exportSymbol('tearDownPage',
tearDownPage)
`)
	c.Test(
		"function testBar() { function testB() {}}",
		`function testBar(){function testB(){}}
google_exportSymbol('testBar',testBar)
`)
	c.TestSame("var testCase = {}; testCase.setUpPage = function() {}")
}

// Java: ExportTestFunctionsTest.testBasicTestFunctionsAreExported_inGoogModule()
func (c *ExportTestFunctionsTest) TestBasicTestFunctionsAreExported_inGoogModule() {
	c.TestSame("goog.module('test'); function Foo() {function testA(){}}")
	c.Test(
		"goog.module('test'); function setUp() {}",
		"goog.module('test'); function setUp(){} google_exportSymbol('setUp',setUp)")
}

// Java: ExportTestFunctionsTest.testFunctionExpressionsAreExported()
func (c *ExportTestFunctionsTest) TestFunctionExpressionsAreExported() {
	c.TestSame("var Foo = function() {var testA = function() {}}")
	c.Test(
		"var setUp = function() {}",
		`var setUp = function() {};
google_exportSymbol('setUp',setUp)
`)
	c.Test(
		"var setUpPage = function() {}",
		`var setUpPage = function() {};
google_exportSymbol('setUpPage',setUpPage)
`)
	c.Test(
		"var shouldRunTests = function() {}",
		`var shouldRunTests = function() {};
google_exportSymbol('shouldRunTests',shouldRunTests)
`)
	c.Test(
		"var tearDown = function() {}",
		`var tearDown = function() {};
google_exportSymbol('tearDown',tearDown)
`)
	c.Test(
		"var tearDownPage = function() {}",
		`var tearDownPage = function() {};
google_exportSymbol('tearDownPage', tearDownPage)
`)
	c.Test(
		"var testBar = function() { var testB = function() {}}",
		`var testBar = function(){ var testB = function() {}};
google_exportSymbol('testBar',testBar)
`)
}

// Java: ExportTestFunctionsTest.testFunctionExpressionsAreExported_inGoogModule()
func (c *ExportTestFunctionsTest) TestFunctionExpressionsAreExported_inGoogModule() {
	c.TestSame("goog.module('test'); var Foo = function() {var testA = function() {}}")
	c.Test(
		"goog.module('test'); var setUp = function() {}",
		"goog.module('test'); var setUp = function() {}; google_exportSymbol('setUp',setUp)")
}

// https://github.com/google/closure-compiler/issues/2563
//
// Java: ExportTestFunctionsTest.testFunctionExpressionsInAssignAreExported()
func (c *ExportTestFunctionsTest) TestFunctionExpressionsInAssignAreExported() {
	c.Test(
		"testBar = function() {};",
		"testBar = function() {}; google_exportSymbol('testBar',testBar)")
}

// Java: ExportTestFunctionsTest.testFunctionExpressionsInAssignAreExported_inGoogModule()
func (c *ExportTestFunctionsTest) TestFunctionExpressionsInAssignAreExported_inGoogModule() {
	c.Test(
		"goog.module('test'); testBar = function() {};",
		"goog.module('test'); testBar = function() {}; google_exportSymbol('testBar',testBar)")
}

// Java: ExportTestFunctionsTest.testFunctionExpressionsByLetAreExported()
func (c *ExportTestFunctionsTest) TestFunctionExpressionsByLetAreExported() {
	c.TestSame("let Foo = function() {var testA = function() {}}")
	c.Test(
		"let setUp = function() {}",
		"let setUp = function() {}; google_exportSymbol('setUp', setUp)")
	c.Test(
		"let testBar = function() {}",
		"let testBar = function() {}; google_exportSymbol('testBar', testBar)")
	c.Test(
		"let tearDown = function() {}",
		`let tearDown = function() {};
google_exportSymbol('tearDown', tearDown)
`)
}

// Java: ExportTestFunctionsTest.testFunctionExpressionsByConstAreExported()
func (c *ExportTestFunctionsTest) TestFunctionExpressionsByConstAreExported() {
	c.TestSame("const Foo = function() {var testA = function() {}}")
	c.Test(
		"const setUp = function() {}",
		`const setUp = function() {};
google_exportSymbol('setUp', setUp)
`)
	c.Test(
		"const testBar = function() {}",
		`const testBar = function() {};
google_exportSymbol('testBar', testBar)
`)
	c.Test(
		"const tearDown = function() {}",
		`const tearDown = function() {};
google_exportSymbol('tearDown', tearDown)
`)
}

// Java: ExportTestFunctionsTest.testArrowFunctionExpressionsAreExported()
func (c *ExportTestFunctionsTest) TestArrowFunctionExpressionsAreExported() {
	c.TestSame("var Foo = ()=>{var testA = function() {}}")
	c.Test(
		"var setUp = ()=>{}",
		`var setUp = ()=>{};
google_exportSymbol('setUp', setUp)
`)
	c.Test(
		"var testBar = ()=>{}",
		`var testBar = ()=>{};
google_exportSymbol('testBar', testBar)
`)
	c.Test(
		"var tearDown = ()=>{}",
		`var tearDown = ()=>{};
google_exportSymbol('tearDown', tearDown)
`)
}

// Java: ExportTestFunctionsTest.testFunctionAssignmentsAreExported()
func (c *ExportTestFunctionsTest) TestFunctionAssignmentsAreExported() {
	c.TestSame("Foo = {}; Foo.prototype.bar = function() {};")

	c.Test(
		"Foo = {}; Foo.prototype.setUpPage = function() {};",
		`Foo = {}; Foo.prototype.setUpPage = function() {};
google_exportProperty(Foo.prototype, 'setUpPage',
Foo.prototype.setUpPage);
`)

	c.Test(
		"Foo = {}; Foo.prototype.shouldRunTests = function() {};",
		`Foo = {}; Foo.prototype.shouldRunTests = function() {};
google_exportProperty(Foo.prototype, 'shouldRunTests',
Foo.prototype.shouldRunTests);
`)

	c.Test(
		"Foo = {}; Foo.prototype.testBar = function() {};",
		`Foo = {}; Foo.prototype.testBar = function() {};
google_exportProperty(Foo.prototype, 'testBar',
Foo.prototype.testBar);
`)

	c.Test(
		"window.testBar = function() {};",
		`window.testBar = function() {};
google_exportProperty(window, 'testBar',
window.testBar);
`)

	c.Test(
		`Foo = {}; Foo.prototype.testBar = function()
{ var testBaz = function() {}};
`,
		`Foo = {}; Foo.prototype.testBar = function()
{ var testBaz = function() {}};
google_exportProperty(Foo.prototype, 'testBar',
Foo.prototype.testBar);
`)

	c.Test(
		`Foo = {}; Foo.baz.prototype.testBar = function()
{ var testBaz = function() {}};
`,
		`Foo = {}; Foo.baz.prototype.testBar = function()
{ var testBaz = function() {}};
google_exportProperty(Foo.baz.prototype, 'testBar',
Foo.baz.prototype.testBar);
`)
}

// Java: ExportTestFunctionsTest.testExportTestSuite()
func (c *ExportTestFunctionsTest) TestExportTestSuite() {
	c.TestSame("goog.testing.testSuite({'a': function() {}, 'b': function() {}});")
	c.Test(
		"goog.testing.testSuite({a: function() {}, b: function() {}});",
		"goog.testing.testSuite({'a': function() {}, 'b': function() {}});")
}

// Java: ExportTestFunctionsTest.testExportTestSuite_whereGoogIsDefined()
func (c *ExportTestFunctionsTest) TestExportTestSuite_whereGoogIsDefined() {
	c.TestSame("var goog = {}; goog.testing.testSuite({'a': function() {}, 'b': function() {}});")
	c.Test(
		"var goog = {}; goog.testing.testSuite({a: function() {}, b: function() {}});",
		"var goog = {}; goog.testing.testSuite({'a': function() {}, 'b': function() {}});")
}

// Java: ExportTestFunctionsTest.testExportTestSuite_inGoogModule()
func (c *ExportTestFunctionsTest) TestExportTestSuite_inGoogModule() {
	c.Test(
		`goog.module('my.test');
const testSuite = goog.require('goog.testing.testSuite');
testSuite({a: function() {}, b: function() {}});
`,
		`goog.module('my.test');
const testSuite = goog.require('goog.testing.testSuite');
testSuite({'a': function() {}, 'b': function() {}});
`)
}

// Java: ExportTestFunctionsTest.testNonGoogTestingTestSuiteNotRewrittenInModule()
func (c *ExportTestFunctionsTest) TestNonGoogTestingTestSuiteNotRewrittenInModule() {
	c.TestSame(
		`goog.module('my.test');
const testSuite = goog.require('some.other.testSuite');
testSuite({a: function() {}, b: function() {}});
`)
}

// Java: ExportTestFunctionsTest.testExportTestSuite_differentNameInGoogModule()
func (c *ExportTestFunctionsTest) TestExportTestSuite_differentNameInGoogModule() {
	c.Test(
		`goog.module('my.test');
const someRandomName = goog.require('goog.testing.testSuite');
someRandomName({a: function() {}, b: function() {}});
`,
		`goog.module('my.test');
const someRandomName = goog.require('goog.testing.testSuite');
someRandomName({'a': function() {}, 'b': function() {}});
`)
}

// Java: ExportTestFunctionsTest.testExportTestSuite_inEsModule()
func (c *ExportTestFunctionsTest) TestExportTestSuite_inEsModule() {
	c.Test(
		`const testSuite = goog.require('goog.testing.testSuite');
testSuite({a: function() {}, b: function() {}});
export {};
`, // Add `export {}` to make this an ES module.
		`const testSuite = goog.require('goog.testing.testSuite');
testSuite({'a': function() {}, 'b': function() {}});
export {};
`)
}

// Java: ExportTestFunctionsTest.testMemberDefInObjLitInTestSuite_becomesStringKey()
func (c *ExportTestFunctionsTest) TestMemberDefInObjLitInTestSuite_becomesStringKey() {
	c.Test(
		"goog.testing.testSuite({a() {}, b() {}});",
		"goog.testing.testSuite({'a': function() {}, 'b': function() {}});")
}

// Java: ExportTestFunctionsTest.testMemberDefInObjLitInTestSuite_becomesStringKey_withSameJSDoc()
func (c *ExportTestFunctionsTest) TestMemberDefInObjLitInTestSuite_becomesStringKey_withSameJSDoc() {
	c.Test(
		`goog.testing.testSuite({
  /** @suppress {checkTypes} */ a() {},
  b() {}
});
`,
		`goog.testing.testSuite({
  /** @suppress {checkTypes} */
  'a': function() {},
  'b': function() {}
});
`)
}

// Java: ExportTestFunctionsTest.testComputedPropInObjLitInTestSuite_doesNotChange()
func (c *ExportTestFunctionsTest) TestComputedPropInObjLitInTestSuite_doesNotChange() {
	c.TestSame(
		`goog.testing.testSuite({
  /** @suppress {checkTypes} */ ['a']() {},
  ['b']() {}
});
`)
}

// Java: ExportTestFunctionsTest.testEs6Class_testMethod()
func (c *ExportTestFunctionsTest) TestEs6Class_testMethod() {
	c.Test(
		"class MyTest {testFoo() {}} goog.testing.testSuite(new MyTest());",
		`class MyTest {testFoo() {}}
google_exportProperty(MyTest.prototype, 'testFoo', MyTest.prototype.testFoo);
goog.testing.testSuite(new MyTest());
`)
}

// Java: ExportTestFunctionsTest.testEs6Class_testMethod_inGoogModule()
func (c *ExportTestFunctionsTest) TestEs6Class_testMethod_inGoogModule() {
	c.Test(
		"goog.module('test'); class MyTest {testFoo() {}} goog.testing.testSuite(new MyTest());",
		`goog.module('test');
class MyTest {testFoo() {}}
google_exportProperty(MyTest.prototype, 'testFoo', MyTest.prototype.testFoo);
goog.testing.testSuite(new MyTest());
`)
}

// Java: ExportTestFunctionsTest.testEs6Class_lifeCycleMethods()
func (c *ExportTestFunctionsTest) TestEs6Class_lifeCycleMethods() {
	c.Test(
		`class MyTest {
testFoo(){} setUp(){} tearDown(){} setUpPage(){} tearDownPage(){} notExported(){}
}
goog.testing.testSuite(new MyTest());
`,
		`class MyTest {
testFoo(){} setUp(){} tearDown(){} setUpPage(){} tearDownPage(){} notExported(){}
}
google_exportProperty(MyTest.prototype, 'testFoo', MyTest.prototype.testFoo);
google_exportProperty(MyTest.prototype, 'setUp', MyTest.prototype.setUp);
google_exportProperty(MyTest.prototype, 'tearDown', MyTest.prototype.tearDown);
google_exportProperty(MyTest.prototype, 'setUpPage', MyTest.prototype.setUpPage);
google_exportProperty(MyTest.prototype, 'tearDownPage',
MyTest.prototype.tearDownPage);
goog.testing.testSuite(new MyTest());
`)
}

// https://github.com/google/closure-compiler/issues/2563
//
// Java: ExportTestFunctionsTest.testES6ClassAssignmentsAreExported()
func (c *ExportTestFunctionsTest) TestES6ClassAssignmentsAreExported() {
	c.TestSame("Foo = class {bar() {}}")

	c.Test(
		"Foo = class {testBar() {}}",
		`Foo = class {testBar() {}};
google_exportProperty(Foo.prototype, 'testBar', Foo.prototype.testBar);
`)
}

// Java: ExportTestFunctionsTest.testEs6Class_testClassExpressionMethod()
func (c *ExportTestFunctionsTest) TestEs6Class_testClassExpressionMethod() {
	c.Test(
		"var MyTest=class{testFoo() {}}; goog.testing.testSuite(new MyTest());",
		`var MyTest=class{testFoo() {}};
google_exportProperty(MyTest.prototype, 'testFoo', MyTest.prototype.testFoo);
goog.testing.testSuite(new MyTest());
`)
}

// Java: ExportTestFunctionsTest.testEs6Class_testClassExpressionByLetMethod()
func (c *ExportTestFunctionsTest) TestEs6Class_testClassExpressionByLetMethod() {
	c.Test(
		"let MyTest=class{testFoo() {}}; goog.testing.testSuite(new MyTest());",
		`let MyTest=class{testFoo() {}};
google_exportProperty(MyTest.prototype, 'testFoo', MyTest.prototype.testFoo);
goog.testing.testSuite(new MyTest());
`)
}

// Java: ExportTestFunctionsTest.testEs6Class_testClassExpressionByConstMethod()
func (c *ExportTestFunctionsTest) TestEs6Class_testClassExpressionByConstMethod() {
	c.Test(
		"const MyTest=class{testFoo() {}}; goog.testing.testSuite(new MyTest());",
		`const MyTest=class{testFoo() {}};
google_exportProperty(MyTest.prototype, 'testFoo', MyTest.prototype.testFoo);
goog.testing.testSuite(new MyTest());
`)
}

// Java: ExportTestFunctionsTest.testCommaOperatorRecognized()
func (c *ExportTestFunctionsTest) TestCommaOperatorRecognized() {
	c.Test(
		c.Srcs("(0, goog.testing.testSuite)({setUp(){console.error('hello');}});"),
		c.Expected("(0, goog.testing.testSuite)({'setUp': function (){console.error('hello');}});"))

	c.Test(
		c.Srcs("(0, goog.testing.testSuite)({testFunction(){console.error('hello');}});"),
		c.Expected(
			"(0, goog.testing.testSuite)({'testFunction': function (){console.error('hello');}});"))
}
