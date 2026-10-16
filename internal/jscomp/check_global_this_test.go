// Port of com/google/javascript/jscomp/CheckGlobalThisTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/CheckGlobalThisTest.java)

/*
 * Copyright 2007 The Closure Compiler Authors.
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

// CheckGlobalThisTest: Tests CheckGlobalThis.
type CheckGlobalThisTest struct {
	CompilerTestCaseBase
}

// Java: CheckGlobalThisTest()
func NewCheckGlobalThisTest() *CheckGlobalThisTest {
	c := &CheckGlobalThisTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestCheckGlobalThisTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCheckGlobalThisTest() })
}

// Java: CheckGlobalThisTest.setUp()
func (c *CheckGlobalThisTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableParseTypeInfo()
}

// Java: CheckGlobalThisTest.getProcessor(Compiler)
func (c *CheckGlobalThisTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewCombinedCompilerPass(compiler, NewCheckGlobalThis(compiler))
}

// Java: CheckGlobalThisTest.testFailure(String)
func (c *CheckGlobalThisTest) TestFailure(js string) {
	c.TestWarning(js, CheckGlobalThis_GLOBAL_THIS)
}

// Java: CheckGlobalThisTest.testGlobalThis1()
func (c *CheckGlobalThisTest) TestGlobalThis1() {
	c.TestSame("var a = this;")
}

// Java: CheckGlobalThisTest.testGlobalThis2()
func (c *CheckGlobalThisTest) TestGlobalThis2() {
	c.TestFailure("this.foo = 5;")
}

// Java: CheckGlobalThisTest.testGlobalThis3()
func (c *CheckGlobalThisTest) TestGlobalThis3() {
	c.TestFailure("this[foo] = 5;")
}

// Java: CheckGlobalThisTest.testGlobalThis4()
func (c *CheckGlobalThisTest) TestGlobalThis4() {
	c.TestFailure("this['foo'] = 5;")
}

// Java: CheckGlobalThisTest.testGlobalThis5()
func (c *CheckGlobalThisTest) TestGlobalThis5() {
	c.TestFailure("(a = this).foo = 4;")
}

// Java: CheckGlobalThisTest.testGlobalThis6()
func (c *CheckGlobalThisTest) TestGlobalThis6() {
	c.TestSame("a = this;")
}

// Java: CheckGlobalThisTest.testGlobalThis7()
func (c *CheckGlobalThisTest) TestGlobalThis7() {
	c.TestFailure("var a = this.foo;")
}

// Java: CheckGlobalThisTest.testStaticFunction1()
func (c *CheckGlobalThisTest) TestStaticFunction1() {
	c.TestSame("function a() { return this; }")
}

// Java: CheckGlobalThisTest.testStaticFunction2()
func (c *CheckGlobalThisTest) TestStaticFunction2() {
	c.TestFailure("function a() { this.complex = 5; }")
}

// Java: CheckGlobalThisTest.testStaticFunction3()
func (c *CheckGlobalThisTest) TestStaticFunction3() {
	c.TestSame("var a = function() { return this; }")
}

// Java: CheckGlobalThisTest.testStaticFunction4()
func (c *CheckGlobalThisTest) TestStaticFunction4() {
	c.TestFailure("var a = function() { this.foo.bar = 6; }")
}

// Java: CheckGlobalThisTest.testStaticFunction5()
func (c *CheckGlobalThisTest) TestStaticFunction5() {
	c.TestSame("function a() { return function() { return this; } }")
}

// Java: CheckGlobalThisTest.testStaticFunction6()
func (c *CheckGlobalThisTest) TestStaticFunction6() {
	c.TestSame("function a() { return function() { this.x = 8; } }")
}

// Java: CheckGlobalThisTest.testStaticFunction7()
func (c *CheckGlobalThisTest) TestStaticFunction7() {
	c.TestSame("var a = function() { return function() { this.x = 8; } }")
}

// Java: CheckGlobalThisTest.testStaticFunction8()
func (c *CheckGlobalThisTest) TestStaticFunction8() {
	c.TestFailure("var a = function() { return this.foo; };")
}

// Java: CheckGlobalThisTest.testConstructor1()
func (c *CheckGlobalThisTest) TestConstructor1() {
	c.TestSame("/** @constructor */function A() { this.m2 = 5; }")
}

// Java: CheckGlobalThisTest.testConstructor2()
func (c *CheckGlobalThisTest) TestConstructor2() {
	c.TestSame("/** @constructor */var A = function() { this.m2 = 5; }")
}

// Java: CheckGlobalThisTest.testConstructor3()
func (c *CheckGlobalThisTest) TestConstructor3() {
	c.TestSame("/** @constructor */a.A = function() { this.m2 = 5; }")
}

// Java: CheckGlobalThisTest.testInterface1()
func (c *CheckGlobalThisTest) TestInterface1() {
	c.TestSame("/** @interface */function A() { /** @type {string} */ this.m2; }")
}

// Java: CheckGlobalThisTest.testOverride1()
func (c *CheckGlobalThisTest) TestOverride1() {
	c.TestSame(
		`/** @constructor */function A() { } var a = new A();
/** @override */ a.foo = function() { this.bar = 5; };
`)
}

// Java: CheckGlobalThisTest.testThisJSDoc1()
func (c *CheckGlobalThisTest) TestThisJSDoc1() {
	c.TestSame("/** @this {whatever} */function h() { this.foo = 56; }")
}

// Java: CheckGlobalThisTest.testThisJSDoc2()
func (c *CheckGlobalThisTest) TestThisJSDoc2() {
	c.TestSame("/** @this {whatever} */var h = function() { this.foo = 56; }")
}

// Java: CheckGlobalThisTest.testThisJSDoc3()
func (c *CheckGlobalThisTest) TestThisJSDoc3() {
	c.TestSame("/** @this {whatever} */foo.bar = function() { this.foo = 56; }")
}

// Java: CheckGlobalThisTest.testThisJSDoc4()
func (c *CheckGlobalThisTest) TestThisJSDoc4() {
	c.TestSame("/** @this {whatever} */function f() { this.foo = 56; }")
}

// Java: CheckGlobalThisTest.testThisJSDoc5()
func (c *CheckGlobalThisTest) TestThisJSDoc5() {
	c.TestSame("function a() { /** @this {x} */function f() { this.foo = 56; } }")
}

// Java: CheckGlobalThisTest.testMethod1()
func (c *CheckGlobalThisTest) TestMethod1() {
	c.TestSame("A.prototype.m1 = function() { this.m2 = 5; }")
}

// Java: CheckGlobalThisTest.testMethod2()
func (c *CheckGlobalThisTest) TestMethod2() {
	c.TestSame("a.B.prototype.m1 = function() { this.m2 = 5; }")
}

// Java: CheckGlobalThisTest.testMethod3()
func (c *CheckGlobalThisTest) TestMethod3() {
	c.TestSame("a.b.c.D.prototype.m1 = function() { this.m2 = 5; }")
}

// Java: CheckGlobalThisTest.testMethod4()
func (c *CheckGlobalThisTest) TestMethod4() {
	c.TestSame("a.prototype['x' + 'y'] =  function() { this.foo = 3; };")
}

// Java: CheckGlobalThisTest.testPropertyOfMethod()
func (c *CheckGlobalThisTest) TestPropertyOfMethod() {
	c.TestFailure(
		`a.protoype.b = {};
a.prototype.b.c = function() { this.foo = 3; };
`)
}

// Java: CheckGlobalThisTest.testStaticMethod1()
func (c *CheckGlobalThisTest) TestStaticMethod1() {
	c.TestFailure("a.b = function() { this.m2 = 5; }")
}

// Java: CheckGlobalThisTest.testStaticMethod2()
func (c *CheckGlobalThisTest) TestStaticMethod2() {
	c.TestSame("a.b = function() { return function() { this.m2 = 5; } }")
}

// Java: CheckGlobalThisTest.testStaticMethod3()
func (c *CheckGlobalThisTest) TestStaticMethod3() {
	c.TestSame("a.b.c = function() { return function() { this.m2 = 5; } }")
}

// Java: CheckGlobalThisTest.testMethodInStaticFunction()
func (c *CheckGlobalThisTest) TestMethodInStaticFunction() {
	c.TestSame("function f() { A.prototype.m1 = function() { this.m2 = 5; } }")
}

// Java: CheckGlobalThisTest.testStaticFunctionInMethod1()
func (c *CheckGlobalThisTest) TestStaticFunctionInMethod1() {
	c.TestSame("A.prototype.m1 = function() { function me() { this.m2 = 5; } }")
}

// Java: CheckGlobalThisTest.testStaticFunctionInMethod2()
func (c *CheckGlobalThisTest) TestStaticFunctionInMethod2() {
	c.TestSame(
		`A.prototype.m1 = function() {
  function me() {
    function myself() {
      function andI() { this.m2 = 5; } } } }
`)
}

// Java: CheckGlobalThisTest.testInnerFunction1()
func (c *CheckGlobalThisTest) TestInnerFunction1() {
	c.TestFailure("function f() { function g() { return this.x; } }")
}

// Java: CheckGlobalThisTest.testInnerFunction2()
func (c *CheckGlobalThisTest) TestInnerFunction2() {
	c.TestFailure("function f() { var g = function() { return this.x; } }")
}

// Java: CheckGlobalThisTest.testInnerFunction3()
func (c *CheckGlobalThisTest) TestInnerFunction3() {
	c.TestFailure("function f() { var x = {}; x.y = function() { return this.x; } }")
}

// Java: CheckGlobalThisTest.testInnerFunction4()
func (c *CheckGlobalThisTest) TestInnerFunction4() {
	c.TestSame("function f() { var x = {}; x.y(function() { return this.x; }); }")
}

// Java: CheckGlobalThisTest.testIssue182a()
func (c *CheckGlobalThisTest) TestIssue182a() {
	c.TestFailure("var NS = {read: function() { return this.foo; }};")
}

// Java: CheckGlobalThisTest.testIssue182b()
func (c *CheckGlobalThisTest) TestIssue182b() {
	c.TestFailure("var NS = {write: function() { this.foo = 3; }};")
}

// Java: CheckGlobalThisTest.testIssue182c()
func (c *CheckGlobalThisTest) TestIssue182c() {
	c.TestFailure("var NS = {}; NS.write2 = function() { this.foo = 3; };")
}

// Java: CheckGlobalThisTest.testIssue182d()
func (c *CheckGlobalThisTest) TestIssue182d() {
	c.TestSame(
		`function Foo() {}
Foo.prototype = {write: function() { this.foo = 3; }};
`)
}

// Java: CheckGlobalThisTest.testLendsAnnotation1()
func (c *CheckGlobalThisTest) TestLendsAnnotation1() {
	c.TestFailure(
		`/** @constructor */ function F() {}
dojo.declare(F, {foo: function() { return this.foo; }});
`)
}

// Java: CheckGlobalThisTest.testLendsAnnotation2()
func (c *CheckGlobalThisTest) TestLendsAnnotation2() {
	c.TestFailure(
		`/** @constructor */ function F() {}
dojo.declare(F, /** @lends {F.bar} */ (
    {foo: function() { return this.foo; }}));
`)
}

// Java: CheckGlobalThisTest.testLendsAnnotation3()
func (c *CheckGlobalThisTest) TestLendsAnnotation3() {
	c.TestSame(
		`/** @constructor */ function F() {}
dojo.declare(F, /** @lends {F.prototype} */ (
    {foo: function() { return this.foo; }}));
`)
}

// Java: CheckGlobalThisTest.testSuppressWarning()
func (c *CheckGlobalThisTest) TestSuppressWarning() {
	c.TestFailure("var x = function() { this.complex = 5; };")
}

// Java: CheckGlobalThisTest.testArrowFunction1()
func (c *CheckGlobalThisTest) TestArrowFunction1() {
	c.TestFailure("var a = () => this.foo;")
}

// Java: CheckGlobalThisTest.testArrowFunction2()
func (c *CheckGlobalThisTest) TestArrowFunction2() {
	c.TestFailure("(() => this.foo)();")
}

// Java: CheckGlobalThisTest.testArrowFunction3()
func (c *CheckGlobalThisTest) TestArrowFunction3() {
	c.TestFailure(
		`function Foo() {}
Foo.prototype.getFoo = () => this.foo;
`)
}

// Java: CheckGlobalThisTest.testArrowFunction4()
func (c *CheckGlobalThisTest) TestArrowFunction4() {
	c.TestFailure(
		`function Foo() {}
Foo.prototype.setFoo = (f) => { this.foo = f; };
`)
}

// Java: CheckGlobalThisTest.testInnerFunctionInClassMethod1()
func (c *CheckGlobalThisTest) TestInnerFunctionInClassMethod1() {
	// TODO(user): It would be nice to warn for using 'this' here
	c.TestSame(
		`function Foo() {}
Foo.prototype.init = function() {
  button.addEventListener('click', function () {
    this.click();
  });
}
Foo.prototype.click = function() {}
`)
}

// Java: CheckGlobalThisTest.testInnerFunctionInClassMethod2()
func (c *CheckGlobalThisTest) TestInnerFunctionInClassMethod2() {
	// TODO(user): It would be nice to warn for using 'this' here
	c.TestSame(
		`function Foo() {
  var x = function() {
    button.addEventListener('click', function () {
      this.click();
    });
  }
}
`)
}

// Java: CheckGlobalThisTest.testInnerFunctionInEs6ClassMethod()
func (c *CheckGlobalThisTest) TestInnerFunctionInEs6ClassMethod() {
	// TODO(user): It would be nice to warn for using 'this' here
	c.TestSame(
		`class Foo {
  constructor() {
    button.addEventListener('click', function () {
      this.click();
    });
  }
  click() {}
}
`)
}

// Java: CheckGlobalThisTest.testStaticBlockThis()
func (c *CheckGlobalThisTest) TestStaticBlockThis() {
	c.TestSame("class Foo { static {var x = this.y;} }")
	c.TestSame("class Foo {static {this.x = 2; }}")
	c.TestSame("class Foo {static {this[this.x] = 3;}}")
	c.TestSame(
		`class Foo {
  static {
    function g() {
      return this.f() + 1;
    }
    var y = g() + 1;
  }
  static f() {return 1;}
}
`)
	c.TestSame(
		`class Foo {
  static {
    button.addEventListener('click', function () {
      this.click();
    });
  }
  static click() {}
}
`)
}

// Java: CheckGlobalThisTest.testES6ClassStaticMethodThis()
func (c *CheckGlobalThisTest) TestES6ClassStaticMethodThis() {
	c.TestSame("class Foo { static h() {var x = this.y;} }")
	c.TestSame("class Foo {static h() {this.x = 2; }}")
	c.TestSame("class Foo {static h() {this[this.x] = 3;}}")
	c.TestSame(
		`class Foo {
  static h() {
    function g() {
      return this.f() + 1;
    }
    var y = g() + 1;
  }
  static f() {return 1;}
}
`)
	c.TestSame(
		`class Foo {
  static h() {
    button.addEventListener('click', function () {
      this.click();
    });
  }
  static click() {}
}
`)
}

// Note: No static initializer block tests because they are always removed by Es6NormalizeClasses.

// Java: CheckGlobalThisTest.testES6ClassField()
func (c *CheckGlobalThisTest) TestES6ClassField() {
	c.TestSame(
		`class Clazz {
  x = 1;
  y = this.x;
}
`)

	c.TestSame(
		`class Clazz {
  static x = 1;
  static y = this.x;
}
`)
}

// Java: CheckGlobalThisTest.testES6ClassComputedProperty_thisInComputed()
func (c *CheckGlobalThisTest) TestES6ClassComputedProperty_thisInComputed() {
	c.TestFailure(
		`class Clazz {
  [this.val] = 1;
}
`)
	c.TestFailure(
		`class Clazz {
  static [this.val] = 1;
}
`)

	c.TestFailure(
		`class Clazz {
  [this.val]() {}
}
`)
	c.TestFailure(
		`class Clazz {
  static [this.val]() {}
}
`)

	c.TestFailure(
		`class Clazz {
  get [this.val]() {
    return 1;
  }
}
`)
	c.TestFailure(
		`class Clazz {
  static get [this.val]() {
    return 1;
  }
}
`)

	c.TestFailure(
		`class Clazz {
  set [this.val](x) {}
}
`)
	c.TestFailure(
		`class Clazz {
  static set [this.val](x) {}
}
`)
}

// Java: CheckGlobalThisTest.testES6ClassComputedProperty_thisInBody()
func (c *CheckGlobalThisTest) TestES6ClassComputedProperty_thisInBody() {
	c.TestSame(
		`class Clazz {
  x = 1;
  ['prop'] = this.x;
}
`)
	c.TestSame(
		`class Clazz {
  static x = 1;
  static ['prop'] = this.x;
}
`)

	c.TestSame(
		`class Clazz {
  x = 1;
  ['prop']() {
    this.x = 2;
  }
}
`)
	c.TestSame(
		`class Clazz {
  static x = 1;
  static ['prop']() {
    this.x = 2;
  }
}
`)

	c.TestSame(
		`class Clazz {
  x = 1
  get ['prop']() {
    return this.x;
  }
}
`)
	c.TestSame(
		`class Clazz {
  static x = 1;
  static get ['prop']() {
    return this.x;
  }
}
`)

	c.TestSame(
		`class Clazz {
  x = 1;
  set ['prop'](x) {
    this.x = x;
  }
}
`)
	c.TestSame(
		`class Clazz {
  static x = 1;
  static set ['prop'](x) {
    this.x = x;
  }
}
`)
}

// Java: CheckGlobalThisTest.testFunctionWithThisTypeAnnotated()
func (c *CheckGlobalThisTest) TestFunctionWithThisTypeAnnotated() {
	c.TestSame(
		`/**
 * @type {function(this:{hello:string})}
 */
function test() {
  console.log(this.hello)
}
`)
}

// Java: CheckGlobalThisTest.testFunctionWithoutThisTypeAnnotated()
func (c *CheckGlobalThisTest) TestFunctionWithoutThisTypeAnnotated() {
	c.TestFailure(
		`/**
 * @type {function()}
 */
function test() {
  console.log(this.hello)
}
`)
}
