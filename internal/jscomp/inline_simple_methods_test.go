// Port of com/google/javascript/jscomp/InlineSimpleMethodsTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/InlineSimpleMethodsTest.java)

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

// InlineSimpleMethodsTest is the unit test of InlineSimpleMethods.
type InlineSimpleMethodsTest struct {
	CompilerTestCaseBase
}

// Java: InlineSimpleMethodsTest()
func NewInlineSimpleMethodsTest() *InlineSimpleMethodsTest {
	c := &InlineSimpleMethodsTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestInlineSimpleMethodsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewInlineSimpleMethodsTest() })
}

// Java: InlineSimpleMethodsTest.getProcessor(Compiler)
func (c *InlineSimpleMethodsTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewInlineSimpleMethods(compiler)
}

// Java: InlineSimpleMethodsTest.setUp()
func (c *InlineSimpleMethodsTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableGatherExternProperties()
	c.EnableNormalize()
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
}

// Helper for tests that expects definitions to remain unchanged, such that {@code definitions+js}
// is converted to {@code definitions+expected}.
// Java: InlineSimpleMethodsTest.testWithPrefix(String, String, String)
func (c *InlineSimpleMethodsTest) TestWithPrefix(definitions string, js string, expected string) {
	c.Test(definitions+js, definitions+expected)
}

// Java: InlineSimpleMethodsTest.testDoesNotInlineMethodOnBaseClass()
func (c *InlineSimpleMethodsTest) TestDoesNotInlineMethodOnBaseClass() {
	baseClassJs := `class Base {
  constructor() {
    /** @const */
    this.prop_ =
        Math.random() > .5;
  }
  method() {
    return this.prop_;
  }
}
`
	derivedClassJS := `class Derived extends Base {
  constructor() {
    super();
  }
  derivedMethod() {
    super.method();
  }
}

(new Derived()).derivedMethod();
`
	source := baseClassJs + derivedClassJS
	c.TestSame(source)
}

// Java: InlineSimpleMethodsTest.testSimpleInline1()
func (c *InlineSimpleMethodsTest) TestSimpleInline1() {
	c.TestWithPrefix(
		`function Foo(){}
Foo.prototype.bar=function(){return this.baz};
`,
		"var x=(new Foo).bar();var y=(new Foo).bar();",
		"var x=(new Foo).baz;var y=(new Foo).baz")
}

// Java: InlineSimpleMethodsTest.testSimpleInline2()
func (c *InlineSimpleMethodsTest) TestSimpleInline2() {
	c.TestWithPrefix(
		`function Foo(){}
Foo.prototype={bar:function(){return this.baz}};
`,
		"var x=(new Foo).bar();var y=(new Foo).bar();",
		"var x=(new Foo).baz;var y=(new Foo).baz")
}

// Java: InlineSimpleMethodsTest.testSimpleGetterInline1()
func (c *InlineSimpleMethodsTest) TestSimpleGetterInline1() {
	// TODO(johnlenz): Support this case.
	c.TestSame(
		`function Foo(){}
Foo.prototype={get bar(){return this.baz}};
var x=(new Foo).bar;var y=(new Foo).bar
`)
	// Verify we are not confusing calling the result of an ES5 getter
	// with call the getter.
	c.TestSame(
		`function Foo(){}
Foo.prototype={get bar(){return this.baz}};
var x=(new Foo).bar();var y=(new Foo).bar()
`)
}

// Java: InlineSimpleMethodsTest.testSimpleSetterInline1()
func (c *InlineSimpleMethodsTest) TestSimpleSetterInline1() {
	// Verify 'get' and 'set' are not confused.
	c.TestSame(
		`function Foo(){}
Foo.prototype={set bar(a){return this.baz}};
var x=(new Foo).bar;var y=(new Foo).bar
`)
	c.TestSame(
		`function Foo(){}
Foo.prototype={set bar(a){return this.baz}};
var x=(new Foo).bar();var y=(new Foo).bar()
`)
}

// Java: InlineSimpleMethodsTest.testSelfInline()
func (c *InlineSimpleMethodsTest) TestSelfInline() {
	c.TestWithPrefix(
		`function Foo(){}
Foo.prototype.bar=function(){return this.baz};
`,
		"Foo.prototype.meth=function(){this.bar();}",
		"Foo.prototype.meth=function(){this.baz}")
}

// Java: InlineSimpleMethodsTest.testCallWithArgs()
func (c *InlineSimpleMethodsTest) TestCallWithArgs() {
	c.TestWithPrefix(
		`function Foo(){}
Foo.prototype.bar=function(){return this.baz};
`,
		"var x=(new Foo).bar(3,new Foo)",
		"var x=(new Foo).bar(3,new Foo)")
}

// Java: InlineSimpleMethodsTest.testCallWithConstArgs()
func (c *InlineSimpleMethodsTest) TestCallWithConstArgs() {
	c.TestWithPrefix(
		`function Foo(){}
Foo.prototype.bar=function(a){return this.baz};
`,
		"var x=(new Foo).bar(3, 4)",
		"var x=(new Foo).baz")
}

// Java: InlineSimpleMethodsTest.testNestedProperties()
func (c *InlineSimpleMethodsTest) TestNestedProperties() {
	c.TestWithPrefix(
		`function Foo(){}
Foo.prototype.bar=function(){return this.baz.ooka};
`,
		"(new Foo).bar()",
		"(new Foo).baz.ooka")
}

// Java: InlineSimpleMethodsTest.testSkipComplexMethods()
func (c *InlineSimpleMethodsTest) TestSkipComplexMethods() {
	c.TestWithPrefix(
		`function Foo(){}
Foo.prototype.bar=function(){return this.baz};
Foo.prototype.condy=function(){return this.baz?this.baz:1};
`,
		"var x=(new Foo).argy()",
		"var x=(new Foo).argy()")
}

// Java: InlineSimpleMethodsTest.testSkipConflictingMethods()
func (c *InlineSimpleMethodsTest) TestSkipConflictingMethods() {
	c.TestWithPrefix(
		`function Foo(){}
Foo.prototype.bar=function(){return this.baz};
Foo.prototype.bar=function(){return this.bazz};
`,
		"var x=(new Foo).bar()",
		"var x=(new Foo).bar()")
}

// Java: InlineSimpleMethodsTest.testSameNamesDifferentDefinitions()
func (c *InlineSimpleMethodsTest) TestSameNamesDifferentDefinitions() {
	c.TestWithPrefix(
		`function A(){}
A.prototype.g=function(){return this.a};
function B(){}
B.prototype.g=function(){return this.b};
`,
		`var x=(new A).g();
var y=(new B).g();
var a=new A;
var ag=a.g();
`,
		`var x=(new A).g();
var y=(new B).g();
var a=new A;
var ag=a.g()
`)
}

// Java: InlineSimpleMethodsTest.testSameNamesSameDefinitions()
func (c *InlineSimpleMethodsTest) TestSameNamesSameDefinitions() {
	c.TestWithPrefix(
		`function A(){}
A.prototype.g=function(){return this.a};
function B(){}
B.prototype.g=function(){return this.a};
`,
		`var x=(new A).g();
var y=(new B).g();
var a=new A;
var ag=a.g();
`,
		`var x=(new A).a;
var y=(new B).a;
var a=new A;
var ag=a.a
`)
}

// Java: InlineSimpleMethodsTest.testConfusingNames()
func (c *InlineSimpleMethodsTest) TestConfusingNames() {
	c.TestWithPrefix(
		`function Foo(){}
Foo.prototype.bar=function(){return this.baz};
`,
		"function bar(){var bar=function(){};bar()}",
		"function bar(){var bar=function(){};bar()}")
}

// Java: InlineSimpleMethodsTest.testConstantInline()
func (c *InlineSimpleMethodsTest) TestConstantInline() {
	c.TestWithPrefix(
		`function Foo(){}
Foo.prototype.bar=function(){return 3};
`,
		"var f=new Foo;var x=f.bar()",
		"var f=new Foo;var x=3")
}

// Java: InlineSimpleMethodsTest.testConstantArrayInline()
func (c *InlineSimpleMethodsTest) TestConstantArrayInline() {
	c.TestWithPrefix(
		`function Foo(){}
Foo.prototype.bar=function(){return[3,4]};
`,
		"var f=new Foo;var x=f.bar()",
		"var f=new Foo;var x=[3,4]")
}

// Java: InlineSimpleMethodsTest.testConstantInlineWithSideEffects()
func (c *InlineSimpleMethodsTest) TestConstantInlineWithSideEffects() {
	c.TestWithPrefix(
		`function Foo(){}
Foo.prototype.bar=function(){return 3};
`,
		"var x=(new Foo).bar()",
		"var x=(new Foo).bar()")
}

// Java: InlineSimpleMethodsTest.testEmptyMethodInline()
func (c *InlineSimpleMethodsTest) TestEmptyMethodInline() {
	c.TestWithPrefix(
		"function Foo(){} Foo.prototype.bar=function(a){};",
		"var x = new Foo(); x.bar();",
		"var x = new Foo();;")
}

// Java: InlineSimpleMethodsTest.testEmptyMethodInlineWithSideEffects()
func (c *InlineSimpleMethodsTest) TestEmptyMethodInlineWithSideEffects() {
	c.TestWithPrefix(
		`function Foo(){}
Foo.prototype.bar=function(){};
`,
		"(new Foo).bar();var y=new Foo;y.bar(new Foo)",
		"(new Foo).bar();var y=new Foo;y.bar(new Foo)")
}

// Java: InlineSimpleMethodsTest.testEmptyMethodInlineInAssign1()
func (c *InlineSimpleMethodsTest) TestEmptyMethodInlineInAssign1() {
	c.TestWithPrefix(
		`function Foo(){}
Foo.prototype.bar=function(){};
`,
		"var x=new Foo;var y=x.bar()",
		"var x=new Foo;var y=void 0")
}

// Java: InlineSimpleMethodsTest.testEmptyMethodInlineInAssign2()
func (c *InlineSimpleMethodsTest) TestEmptyMethodInlineInAssign2() {
	c.TestWithPrefix(
		`function Foo(){}
Foo.prototype.bar=function(){};
`,
		"var x=new Foo;var y=x.bar().toString()",
		"var x=new Foo;var y=(void 0).toString()")
}

// Java: InlineSimpleMethodsTest.testNormalMethod()
func (c *InlineSimpleMethodsTest) TestNormalMethod() {
	c.TestWithPrefix(
		`function Foo(){}
Foo.prototype.bar=function(){var x=1};
`,
		"var x=new Foo;x.bar()",
		"var x=new Foo;x.bar()")
}

// Java: InlineSimpleMethodsTest.testNoInlineOfExternMethods1()
func (c *InlineSimpleMethodsTest) TestNoInlineOfExternMethods1() {
	c.TestSame(c.Externs("var external={};external.charAt;"), c.Srcs("external.charAt()"))
}

// Java: InlineSimpleMethodsTest.testNoInlineOfExternMethods2()
func (c *InlineSimpleMethodsTest) TestNoInlineOfExternMethods2() {
	c.TestSame(c.Externs("var external={};external.charAt=function(){};"), c.Srcs("external.charAt()"))
}

// Java: InlineSimpleMethodsTest.testNoInlineOfExternMethods3()
func (c *InlineSimpleMethodsTest) TestNoInlineOfExternMethods3() {
	c.TestSame(
		c.Externs("var external={};external.bar=function(){};"),
		c.Srcs("function Foo(){}Foo.prototype.bar=function(){};(new Foo).bar()"))
}

// Java: InlineSimpleMethodsTest.testNoInlineOfDangerousProperty()
func (c *InlineSimpleMethodsTest) TestNoInlineOfDangerousProperty() {
	c.TestSame(
		`function Foo(){this.bar=3}
Foo.prototype.bar=function(){};
var x=new Foo;var y=x.bar()
`)
}

// Don't warn about argument naming conventions (this is done in another pass)
//
//	opt_ parameters must not be followed by non-optional parameters.
//	var_args must be last
//
// Java: InlineSimpleMethodsTest.testNoWarn()
func (c *InlineSimpleMethodsTest) TestNoWarn() {
	c.TestSame(
		`function Foo(){}
Foo.prototype.bar=function(opt_a,b){var x=1};
var x=new Foo;x.bar()
`)
	c.TestSame(
		`function Foo(){}
Foo.prototype.bar=function(var_args,b){var x=1};
var x=new Foo;x.bar()
`)
}

// Java: InlineSimpleMethodsTest.testObjectLit()
func (c *InlineSimpleMethodsTest) TestObjectLit() {
	c.TestSame(
		`Foo.prototype.bar=function(){return this.baz_};
var blah={bar:function(){}};
(new Foo).bar()
`)
}

// Java: InlineSimpleMethodsTest.testObjectLit2()
func (c *InlineSimpleMethodsTest) TestObjectLit2() {
	c.TestSame(
		`var blah={bar:function(){}};
(new Foo).bar()
`)
}

// Java: InlineSimpleMethodsTest.testObjectLit3()
func (c *InlineSimpleMethodsTest) TestObjectLit3() {
	c.TestSame(
		`var blah={bar(){}};
(new Foo).bar()
`)
}

// Java: InlineSimpleMethodsTest.testObjectLit4()
func (c *InlineSimpleMethodsTest) TestObjectLit4() {
	c.TestSame(
		`var key='bar';
var blah={[key]:a};
(new Foo).bar
`)
}

// Java: InlineSimpleMethodsTest.testObjectLitExtern1()
func (c *InlineSimpleMethodsTest) TestObjectLitExtern1() {
	externs := "window.bridge={_sip:function(){}};"
	c.TestSame(c.Externs(externs), c.Srcs("window.bridge._sip()"))
}

// Java: InlineSimpleMethodsTest.testObjectLitExtern2()
func (c *InlineSimpleMethodsTest) TestObjectLitExtern2() {
	externs := "window.bridge={_sip(){}};"
	c.TestSame(c.Externs(externs), c.Srcs("window.bridge._sip()"))
}

// Java: InlineSimpleMethodsTest.testClassExtern()
func (c *InlineSimpleMethodsTest) TestClassExtern() {
	externs := "window.bridge= class { _sip() {} };"
	c.TestSame(c.Externs(externs), c.Srcs("window.bridge._sip()"))
}

// Java: InlineSimpleMethodsTest.testExternFunction()
func (c *InlineSimpleMethodsTest) TestExternFunction() {
	externs := "function emptyFunction() {}"
	c.TestSame(
		c.Externs(externs),
		c.Srcs(
			`function Foo(){this.empty=emptyFunction}
(new Foo).empty()
`))
}

// Java: InlineSimpleMethodsTest.testIssue2508576_1()
func (c *InlineSimpleMethodsTest) TestIssue2508576_1() {
	// Method defined by an extern should be left alone.
	externs := "function alert(a) {}"
	c.TestSame(c.Externs(externs), c.Srcs("({a:alert,b:alert}).a('a')"))
}

// Java: InlineSimpleMethodsTest.testIssue2508576_2()
func (c *InlineSimpleMethodsTest) TestIssue2508576_2() {
	// Anonymous object definition with a side-effect should be left alone.
	c.TestSame("({a:function(){},b:x()}).a('a')")
}

// Java: InlineSimpleMethodsTest.testIssue2508576_3()
func (c *InlineSimpleMethodsTest) TestIssue2508576_3() {
	// Anonymous object definition without side-effect should be removed.
	c.Test("({a:function(){},b:alert}).a('a')", ";")
}

// When there are two methods with the same name, one on an ES3-style class, and one on an
// ES6-style class, make sure the right one gets inlined into the right place, if inlining happens
// at all.
// Java: InlineSimpleMethodsTest.testEs6Issue1()
func (c *InlineSimpleMethodsTest) TestEs6Issue1() {
	c.TestSame(
		`/** @constructor */
function OldClass() {}

OldClass.prototype.foo = function() { return this.oldbar; };

class NewClass {
  foo() { return this.newbar; }
}

var x = new OldClass;
x.foo();
x = new NewClass;
x.foo();
`)
}

// Java: InlineSimpleMethodsTest.testEs6Issue2()
func (c *InlineSimpleMethodsTest) TestEs6Issue2() {
	c.TestSame(
		`/** @constructor */
function OldClass() {}

OldClass.prototype.foo = function() { return this.oldbar; };

class NewClass {
  foo() { return this.newbar; }
}

var x = new OldClass;
x.foo();
var y = new NewClass;
y.foo();
`)
}

// Java: InlineSimpleMethodsTest.testArrowFunction()
func (c *InlineSimpleMethodsTest) TestArrowFunction() {
	// Arrow functions cannot be trivially inlined.
	// This is only an issue if ES6+ is being targeted,
	// because otherwise the arrow function is removed before this pass.
	c.TestSame(
		`class Holder {
    constructor() {
        this.val = {internal: true};
        this.context =  {getVal: () => this.val}
    }
}

console.log(new Holder().context.getVal().internal);
`)
}

// Java: InlineSimpleMethodsTest.testAnonymousGet()
func (c *InlineSimpleMethodsTest) TestAnonymousGet() {
	// Anonymous object definition without side-effect should be removed.
	c.TestSame("({get a(){return function(){}},b:alert}).a('a')")
	c.TestSame("({get a(){},b:alert}).a('a')")
	c.TestSame("({get a(){},b:alert}).a")
}

// Java: InlineSimpleMethodsTest.testAnonymousSet()
func (c *InlineSimpleMethodsTest) TestAnonymousSet() {
	// Anonymous object definition without side-effect should be removed.
	c.TestSame("({set a(b){return function(){}},b:alert}).a('a')")
	c.TestSame("({set a(b){},b:alert}).a('a')")
	c.TestSame("({set a(b){},b:alert}).a")
}

// Java: InlineSimpleMethodsTest.testInlinesEvenIfClassEscapes()
func (c *InlineSimpleMethodsTest) TestInlinesEvenIfClassEscapes() {
	// The purpose of this test is to record an unsafe assumption made by the
	// pass. In practice, it's usually safe to inline even if the class
	// escapes, because method definitions aren't commonly mutated.
	c.Test(
		c.Externs("var esc;"),
		c.Srcs(
			`/** @constructor */
function Foo() {
  this.prop = 123;
}
Foo.prototype.m = function() {
  return this.prop;
}
(new Foo).m();
esc(Foo);
`),
		c.Expected(
			`/** @constructor */
function Foo(){this.prop=123}
Foo.prototype.m=function(){return this.prop}
(new Foo).m();
esc(Foo)
`))
}

// Java: InlineSimpleMethodsTest.testClassFieldDoesntCrash()
func (c *InlineSimpleMethodsTest) TestClassFieldDoesntCrash() {
	c.TestSame("class C { a; };")
	c.TestSame("class C { a = 2; };")
	c.TestSame("class C { ['a'] = 2; }")
	c.TestSame("class C { a = function() { return 4; }; }")
	c.TestSame("class C { static a; };")
	c.TestSame("class C { static a = 2; };")
	c.TestSame("class C { static ['a'] = 2; }")
	c.TestSame("class C { static a = function() { return 4; }; }")
}

// Java: InlineSimpleMethodsTest.testStaticInitializationBlockDoesntCrash()
func (c *InlineSimpleMethodsTest) TestStaticInitializationBlockDoesntCrash() {
	c.TestSame(
		`class C {
  static {
    alert('foo');
  }
}
`)
}

// Java: InlineSimpleMethodsTest.testNonStaticClassFieldArrowFunction()
func (c *InlineSimpleMethodsTest) TestNonStaticClassFieldArrowFunction() {
	c.TestSame(
		`class Foo {
  a = 5;
  b = () => this.a;
}
new Foo().b()
`)
}

// Java: InlineSimpleMethodsTest.testStaticClassFieldFunctionDoesInline()
func (c *InlineSimpleMethodsTest) TestStaticClassFieldFunctionDoesInline() {
	c.Test(
		`class Foo {
  static a = 5;
  static b = function() { return this.a; };
}
Foo.b();
`,
		`class Foo {
  static a = 5;
  static b = function() { return this.a; };
}
Foo.a;
`)
}

// Java: InlineSimpleMethodsTest.testNoInline()
func (c *InlineSimpleMethodsTest) TestNoInline() {
	c.TestWithPrefix(
		`class Foo {
 /** @noinline */ bar() { return 'hi'; }
}
`,
		"var x=new Foo;x.bar()",
		"var x=new Foo;x.bar()")
}

// Java: InlineSimpleMethodsTest.testReflectObjectProperty()
func (c *InlineSimpleMethodsTest) TestReflectObjectProperty() {
	c.TestWithPrefix(
		`class Foo {
 bar() { return 'hi'; }
}
const c = goog.reflect.objectProperty('bar', Foo.prototype);
`,
		"var x=new Foo;x.bar()",
		"var x=new Foo;x.bar()")
}

// Java: InlineSimpleMethodsTest.testNoInlineOfAsyncMethod_es5Style()
func (c *InlineSimpleMethodsTest) TestNoInlineOfAsyncMethod_es5Style() {
	// Don't inline: the bar() call returns a Promise wrapping this.baz.
	c.TestSame(
		`function Foo(){}
Foo.prototype.bar = async function(){return this.baz};
var x = (new Foo).bar();
var y = (new Foo).bar();
`)
}

// Java: InlineSimpleMethodsTest.testNoInlineOfAsyncMethod_es6Style()
func (c *InlineSimpleMethodsTest) TestNoInlineOfAsyncMethod_es6Style() {
	// Don't inline: the bar() call returns a Promise wrapping this.baz.
	c.TestSame(
		`class Foo { async bar(){return this.baz} }
var x = (new Foo).bar();
var y = (new Foo).bar();
`)
}

// Java: InlineSimpleMethodsTest.testNoInlineOfGeneratorMethod_es5Style()
func (c *InlineSimpleMethodsTest) TestNoInlineOfGeneratorMethod_es5Style() {
	// Don't inline: the bar() call returns a generator, not this.baz
	c.TestSame(
		`function Foo(){}
Foo.prototype.bar = function*(){return this.baz};
var x = (new Foo).bar();
var y = (new Foo).bar();
`)
}

// Java: InlineSimpleMethodsTest.testNoInlineOfGeneratorMethod_es6Style()
func (c *InlineSimpleMethodsTest) TestNoInlineOfGeneratorMethod_es6Style() {
	// Don't inline: the bar() call returns a generator, not this.baz
	c.TestSame(
		`class Foo { *bar(){return this.baz} }
var x = (new Foo).bar();
var y = (new Foo).bar();
`)
}

// Java: InlineSimpleMethodsTest.testNoInlineOfAsynceneratorMethod_es5Style()
func (c *InlineSimpleMethodsTest) TestNoInlineOfAsynceneratorMethod_es5Style() {
	// Don't inline: the bar() call returns an async generator, not this.baz
	c.TestSame(
		`function Foo(){}
Foo.prototype.bar = async function*(){return this.baz};
var x = (new Foo).bar();
var y =(new Foo).bar();
`)
}

// Java: InlineSimpleMethodsTest.testNoInlineOfAsyncGeneratorMethod_es6Style()
func (c *InlineSimpleMethodsTest) TestNoInlineOfAsyncGeneratorMethod_es6Style() {
	// Don't inline: the bar() call returns an async generator, not this.baz
	c.TestSame(
		`class Foo { async *bar(){return this.baz} }
var x = (new Foo).bar();
var y = (new Foo).bar();
`)
}

// Java: InlineSimpleMethodsTest.testNoInlineOfMethodWithSideEffectfulDefaultParameter()
func (c *InlineSimpleMethodsTest) TestNoInlineOfMethodWithSideEffectfulDefaultParameter() {
	c.TestSame(
		`class Foo {
  getVal(x = audit()) { return this.val; }
}
var x = (new Foo).getVal();
`)
}

// Java: InlineSimpleMethodsTest.testInlineOfMethodWithSideEffectFreeDefaultParameter()
func (c *InlineSimpleMethodsTest) TestInlineOfMethodWithSideEffectFreeDefaultParameter() {
	c.Test(
		`class Foo {
  getVal(x = 1) { return this.val; }
}
var x = (new Foo).getVal();
`,
		`class Foo {
  getVal(x = 1) { return this.val; }
}
var x = (new Foo).val;
`)
}

// Java: InlineSimpleMethodsTest.testNoInlineOfEmptyMethodWithSideEffectfulDefaultParameter()
func (c *InlineSimpleMethodsTest) TestNoInlineOfEmptyMethodWithSideEffectfulDefaultParameter() {
	c.TestSame(
		`class Foo {
  record(x = audit()) {}
}
var f = new Foo();
var x = f.record();
`)
}

// Java: InlineSimpleMethodsTest.testInlineOfEmptyMethodWithSideEffectFreeDefaultParameter()
func (c *InlineSimpleMethodsTest) TestInlineOfEmptyMethodWithSideEffectFreeDefaultParameter() {
	c.Test(
		`class Foo {
  record(x = 1) {}
}
var f = new Foo();
var x = f.record();
`,
		`class Foo {
  record(x = 1) {}
}
var f = new Foo();
var x = void 0;
`)
}

// Java: InlineSimpleMethodsTest.testNoInlineOfMethodWithDestructuring()
func (c *InlineSimpleMethodsTest) TestNoInlineOfMethodWithDestructuring() {
	c.TestSame(
		`class Foo {
  getVal({x}) { return this.val; }
}
var x = (new Foo).getVal({});
`)
	c.TestSame(
		`class Foo {
  getVal([x]) { return this.val; }
}
var x = (new Foo).getVal([1]);
`)
	c.TestSame(
		`class Foo {
  getVal({x} = {}) { return this.val; }
}
var x = (new Foo).getVal();
`)
}

// Java: InlineSimpleMethodsTest.testNoInlineOfEmptyMethodWithDestructuring()
func (c *InlineSimpleMethodsTest) TestNoInlineOfEmptyMethodWithDestructuring() {
	c.TestSame(
		`class Foo {
  record({x}) {}
}
var f = new Foo();
var x = f.record({});
`)
	c.TestSame(
		`class Foo {
  record([x]) {}
}
var f = new Foo();
var x = f.record([1]);
`)
}

// Java: InlineSimpleMethodsTest.testInlineOfMethodWithRestParameter()
func (c *InlineSimpleMethodsTest) TestInlineOfMethodWithRestParameter() {
	c.Test(
		`class Foo {
  getVal(...x) { return this.val; }
}
var x = (new Foo).getVal();
`,
		`class Foo {
  getVal(...x) { return this.val; }
}
var x = (new Foo).val;
`)
}

// Java: InlineSimpleMethodsTest.testInlineOfEmptyMethodWithRestParameter()
func (c *InlineSimpleMethodsTest) TestInlineOfEmptyMethodWithRestParameter() {
	c.Test(
		`class Foo {
  record(...x) {}
}
var f = new Foo();
var x = f.record();
`,
		`class Foo {
  record(...x) {}
}
var f = new Foo();
var x = void 0;
`)
}
