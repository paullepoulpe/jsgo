// Port of com/google/javascript/jscomp/OptimizeConstructorsTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/OptimizeConstructorsTest.java)

/*
 * Copyright 2021 The Closure Compiler Authors.
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

type OptimizeConstructorsTest struct {
	CompilerTestCaseBase
}

// Java: OptimizeConstructorsTest()
func NewOptimizeConstructorsTest() *OptimizeConstructorsTest {
	c := &OptimizeConstructorsTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, CompilerTestCase_DEFAULT_EXTERNS+"var alert;var use;")
	return c
}

func TestOptimizeConstructorsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewOptimizeConstructorsTest() })
}

// Java: OptimizeConstructorsTest.getProcessor(Compiler)
func (c *OptimizeConstructorsTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewOptimizeConstructors(compiler)
}

// Java: OptimizeConstructorsTest.setUp()
func (c *OptimizeConstructorsTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableNormalize()
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.EnableGatherExternProperties()
}

// Test case: reassignment:
//   class A {};  A = ...

// Test case: reassignment with destructuring:
//   class A {};  {A} = {A:x}

// Java: OptimizeConstructorsTest.testSimple()
func (c *OptimizeConstructorsTest) TestSimple() {
	// As simple a test case as I can come up with...
	c.Test(
		"class S { constructor() {} }\n"+
			"class C extends S { constructor() { super(); } }\n"+
			"let c = new C();\n",
		"class S { }\n"+
			"class C extends S {}\n"+
			"let c = new C();\n")
}

// Java: OptimizeConstructorsTest.testSimpleClassExpression()
func (c *OptimizeConstructorsTest) TestSimpleClassExpression() {
	// As simple a test case as I can come up with...
	c.Test(
		"const S = class { constructor() {} }\n"+
			"let C = class extends S { constructor() { super(); } };\n"+
			"let instance = new C();\n",
		"const S = class { };\n"+
			"let C = class extends S {};\n"+
			"let instance = new C();\n")
}

// Java: OptimizeConstructorsTest.testClassExpressionWithInComma()
func (c *OptimizeConstructorsTest) TestClassExpressionWithInComma() {
	// The use of the class being modified doesn't matter so
	// definitions in complex expressions is ok.
	c.Test(
		"class S { constructor() {} }\n"+
			"let C = use ?? class extends S { constructor() { super(); } };\n"+
			"let c = new C();\n",
		"class S { }\n"+
			"let C = use ?? class extends S { };\n"+
			"let c = new C();\n")
}

// Java: OptimizeConstructorsTest.testClassExpressionDefinedWithComma()
func (c *OptimizeConstructorsTest) TestClassExpressionDefinedWithComma() {
	// As simple a test case as I can come up with...
	c.Test(
		"const S = (0, class { constructor() {} });\n"+
			"let C = class extends S { constructor() { super(); } };\n"+
			"let instance = new C();\n",
		"const S = (0, class { });\n"+
			"let C = class extends S {};\n"+
			"let instance = new C();\n")
}

// Java: OptimizeConstructorsTest.testES5SuperClass()
func (c *OptimizeConstructorsTest) TestES5SuperClass() {
	// NOTE: we can remove subclasses of well defined ES5 classes
	c.Test(
		"/** @constructor */ let S = function() {};\n"+
			"class C extends S { constructor() { super(); } }\n"+
			"let c = new C();\n",
		"/** @constructor */ let S = function() {};\n"+
			"class C extends S { }\n"+
			"let c = new C();\n")
}

// Java: OptimizeConstructorsTest.testES5SubClass()
func (c *OptimizeConstructorsTest) TestES5SubClass() {
	c.Test(
		"class S { constructor() {} }\n"+
			"class C extends S { constructor() { super(); } }\n"+
			"function E() { return Reflect.construct(C); }\n"+
			"let c = new C();\n",
		"class S { }\n"+
			"class C extends S { }\n"+
			"function E() { return Reflect.construct(C); }\n"+
			"let c = new C();\n")
}

// Java: OptimizeConstructorsTest.testParameterMismatch1()
func (c *OptimizeConstructorsTest) TestParameterMismatch1() {
	c.TestSame(
		"class S { constructor(a=undefined) {use(a);} }\n" +
			"class C extends S { constructor() { super(); } }\n" +
			"let c = new C(1);\n")
}

// Java: OptimizeConstructorsTest.testParameterMismatch2()
func (c *OptimizeConstructorsTest) TestParameterMismatch2() {
	// NOTE: super class asks for "rest" so we don't remove subclass constructors
	// that don't provide all of them.  We can do better here though because
	// we know the constructor is removed here because it isn't doing anything interesting.
	c.Test(
		"class S { constructor(...rest) {} }\n"+
			"class C extends S { constructor() { super(); } }\n"+
			"let c = new C(1);\n",
		"class S { }\n"+
			"class C extends S { constructor() { super(); } }\n"+
			"let c = new C(1);\n")
}

// Java: OptimizeConstructorsTest.testParameterMismatch3()
func (c *OptimizeConstructorsTest) TestParameterMismatch3() {
	c.TestSame(
		"class S { constructor() { use(arguments); } }\n" +
			"class C extends S { constructor() { super(); } }\n" +
			"let c = new C(1);\n")
}

// Java: OptimizeConstructorsTest.testOptimize_emptyConstructor()
func (c *OptimizeConstructorsTest) TestOptimize_emptyConstructor() {
	c.Test("class A {  constructor() {} }", "class A {}")
	c.Test("class A {  constructor(a,b,c) {} }", "class A {}")
}

// Java: OptimizeConstructorsTest.testOptimize_noArgs()
func (c *OptimizeConstructorsTest) TestOptimize_noArgs() {
	c.Test(
		"class Super {\n"+
			"  constructor() {\n"+
			"    this.a = 1;\n"+
			"  }\n"+
			"}\n"+
			"class A extends Super {\n"+
			"  constructor() {\n"+
			"    super();\n"+
			"  }\n"+
			"}\n",
		"class Super {\n"+
			"  constructor() {\n"+
			"    this.a = 1;\n"+
			"  }\n"+
			"}\n"+
			"class A extends Super {}\n")
}

// Java: OptimizeConstructorsTest.testOptimize_matchingArgs()
func (c *OptimizeConstructorsTest) TestOptimize_matchingArgs() {
	c.Test(
		"class B {\n"+
			"  constructor(a,b,c) {\n"+
			"    this.a=a;\n"+
			"    this.b=b;\n"+
			"    this.c=c;\n"+
			"  }\n"+
			"}\n"+
			"class A extends B {\n"+
			"  constructor(a,b,c) {\n"+
			"    super(a,b,c);\n"+
			"  }\n"+
			"}\n",
		"class B{\n"+
			"  constructor(a,b,c) {\n"+
			"    this.a=a;\n"+
			"    this.b=b;\n"+
			"    this.c=c;\n"+
			"  }\n"+
			"}\n"+
			"class A extends B {}\n")
}

// Java: OptimizeConstructorsTest.testOptimize_explicitObjectSuper()
func (c *OptimizeConstructorsTest) TestOptimize_explicitObjectSuper() {
	// NOTE: It would be valid to remove the constructor in these cases.
	c.TestSame("class A extends Object { constructor(a) { super(); } }")
	c.TestSame("class A extends Object { constructor(a) { super(a); } }")
	c.TestSame("class A extends Object { constructor() { super(A); }}")
}

// Java: OptimizeConstructorsTest.testOptimize_es5super()
func (c *OptimizeConstructorsTest) TestOptimize_es5super() {
	c.Test(
		"/** @constructor */ function B(a,b,c) {\n"+
			"  this.a=a;\n"+
			"  this.b=b;\n"+
			"  this.c=c;\n"+
			"}\n"+
			"class A extends B {\n"+
			"  constructor(a,b,c) {\n"+
			"    super(a,b,c);\n"+
			"  }\n"+
			"}\n",
		"/** @constructor */ function B(a,b,c) {\n"+
			"  this.a=a;\n"+
			"  this.b=b;\n"+
			"  this.c=c;\n"+
			"}\n"+
			"class A extends B {}\n")
}

// Java: OptimizeConstructorsTest.testOptimize_externSuper()
func (c *OptimizeConstructorsTest) TestOptimize_externSuper() {
	// NOTE: to optimize this we need to assume the external definition is accurately describes
	// the behavior.
	c.TestSame(
		c.Externs(SourceFile_FromCode("externs", "/** @constructor */ function Error(x, y, z) {}")),
		c.Srcs(
			"class A extends Error {\n"+
				"  constructor(a, b, c) {\n"+
				"    super(a, b, c);\n"+
				"  }\n"+
				"}\n"))
}

// Java: OptimizeConstructorsTest.testOptimize_varArgs()
func (c *OptimizeConstructorsTest) TestOptimize_varArgs() {
	c.Test(
		"class Super {\n"+
			"  constructor(...a) {\n"+
			"    this.a=a;\n"+
			"  }\n"+
			"}\n"+
			"class A extends Super {\n"+
			"  constructor(...a) {\n"+
			"    super(...a);\n"+
			"  }\n"+
			"}\n",
		"class Super {\n"+
			"  constructor(...a) {\n"+
			"    this.a=a;\n"+
			"  }\n"+
			"}\n"+
			"class A extends Super {}\n")
}

// Java: OptimizeConstructorsTest.testOptimize_syntheticConstructor()
func (c *OptimizeConstructorsTest) TestOptimize_syntheticConstructor() {
	// NOTE: This test demonstrates the current behavior but
	// the subclass constructor here is safely removable.
	c.TestSame(
		"class Super {\n" +
			"  constructor(a, b, c) {\n" +
			"    this.a=a;\n" +
			"  }\n" +
			"}\n" +
			"class A extends Super {\n" +
			"  constructor() {\n" +
			"    super(...arguments);\n" +
			"  }\n" +
			"}\n")
}

// Java: OptimizeConstructorsTest.testOptimize_implicitConstructor()
func (c *OptimizeConstructorsTest) TestOptimize_implicitConstructor() {
	// NOTE: extend the handling to allow for
	// intermediate class with implicit constructors
	c.TestSame(
		"class Top {\n" +
			"  constructor(a) {\n" +
			"    this.a=1;\n" +
			"  }\n" +
			"}\n" +
			"class Super extends Top {}\n" +
			"class A extends Super {\n" +
			"  constructor(b) {\n" +
			"    super(b);\n" +
			"  }\n" +
			"}\n")
}

// Java: OptimizeConstructorsTest.testNoOptimize_wrongParameterOrder()
func (c *OptimizeConstructorsTest) TestNoOptimize_wrongParameterOrder() {
	c.TestSame(
		"class B {\n" +
			"  constructor(a,b,c) { this.a = a; }\n" +
			"}\n" +
			"class A extends B {\n" +
			"  constructor(a,b,c) {\n" +
			"    super(b,a,c);\n" +
			"  }\n" +
			"}\n")
}

// Java: OptimizeConstructorsTest.testNoOptimize_notEveryParameterUsed()
func (c *OptimizeConstructorsTest) TestNoOptimize_notEveryParameterUsed() {
	c.TestSame(
		"class B {\n" +
			"  constructor(a,b) {\n" +
			"    this.a=a;\n" +
			"    this.b=b;\n" +
			"  }\n" +
			"}\n" +
			"class A extends B {\n" +
			"  constructor(a,b,c) {\n" +
			"    super(a,b);\n" +
			"  }\n" +
			"}\n")
}

// Java: OptimizeConstructorsTest.testNoOptimize_nonEmptyConstructor()
func (c *OptimizeConstructorsTest) TestNoOptimize_nonEmptyConstructor() {
	c.TestSame(
		"class A {\n" +
			"  constructor() {\n" +
			"    A.$clinit();\n" +
			"  }\n" +
			"  static $clinit() {}\n" +
			"}\n" +
			"class C {\n" +
			"  constructor(a, b, c) {\n" +
			"    this.a = a;\n" +
			"  }\n" +
			"}\n" +
			"class B extends C {\n" +
			"  constructor(a, b, c) {\n" +
			"    B.$clinit();\n" +
			"    super(a, b, c);\n" +
			"  }\n" +
			"  static $clinit() {}\n" +
			"}\n")
}

// Java: OptimizeConstructorsTest.testNoOptimize_superNotMatching_es6()
func (c *OptimizeConstructorsTest) TestNoOptimize_superNotMatching_es6() {
	c.TestSame(
		"class B {\n" +
			"  constructor(opt_a,opt_b) {\n" +
			"    this.a=opt_a;\n" +
			"    this.b=opt_b;\n" +
			"  }\n" +
			"}\n" +
			"class A extends B { constructor() { super(); } }\n")
}

// Java: OptimizeConstructorsTest.testNoOptimize_superNotMatching_es5()
func (c *OptimizeConstructorsTest) TestNoOptimize_superNotMatching_es5() {
	c.TestSame(
		"/** @constructor */ function B(opt_a,opt_b) {\n" +
			"  this.a=opt_a;\n" +
			"  this.b=opt_b;\n" +
			"}\n" +
			"class A extends B { constructor() { super(); } }\n")
}

// Java: OptimizeConstructorsTest.testNoOptimize_superNotMatching_extern()
func (c *OptimizeConstructorsTest) TestNoOptimize_superNotMatching_extern() {
	c.TestSame(
		c.Externs(
			SourceFile_FromCode(
				"externs", "/** @constructor */ function Error(opt_a, opt_b, opt_c) {}")),
		c.Srcs("class A extends Error { constructor() { super(); } }"))
}

// Java: OptimizeConstructorsTest.testNoOptimize_superAcceptsRest()
func (c *OptimizeConstructorsTest) TestNoOptimize_superAcceptsRest() {
	c.TestSame(
		"class Super {\n" +
			"  constructor(...a) {\n" +
			"    this.a=a;\n" +
			"  }\n" +
			"}\n" +
			"class A extends Super {\n" +
			"  constructor(a) {\n" +
			"    super(a);\n" +
			"  }\n" +
			"}\n")
}

// Java: OptimizeConstructorsTest.testNoOptimize_restNotSread()
func (c *OptimizeConstructorsTest) TestNoOptimize_restNotSread() {
	c.TestSame(
		"class Super {\n" +
			"  constructor(...a) {\n" +
			"    this.a=a;\n" +
			"  }\n" +
			"}\n" +
			"class A extends Super {\n" +
			"  constructor(...a) {\n" +
			"    super(a);\n" +
			"  }\n" +
			"}\n")
}

// Java: OptimizeConstructorsTest.testNoOptimize_unknownParent()
func (c *OptimizeConstructorsTest) TestNoOptimize_unknownParent() {
	c.TestSame(
		"var Super;\n" +
			"class A extends Super {\n" +
			"  constructor() {\n" +
			"    super();\n" +
			"  }\n" +
			"}\n")
}

// Java: OptimizeConstructorsTest.testNoOptimize_extendsExpression()
func (c *OptimizeConstructorsTest) TestNoOptimize_extendsExpression() {
	c.TestSame(
		"function f() { return undefined; }\n" +
			"class A extends f() {\n" +
			"  constructor() {\n" +
			"    super();\n" +
			"  }\n" +
			"}\n")
}

// Java: OptimizeConstructorsTest.testNoOptimize_implicitConstructor()
func (c *OptimizeConstructorsTest) TestNoOptimize_implicitConstructor() {
	c.TestSame(
		"class Top {\n" +
			"  constructor(opt_a) {\n" +
			"    this.a=1;\n" +
			"  }\n" +
			"}\n" +
			"class Super extends Top {}\n" +
			"class A extends Super {\n" +
			"  constructor() {\n" +
			"    super();\n" +
			"  }\n" +
			"}\n")
}

// Java: OptimizeConstructorsTest.testNoOptimize_sideEffectObjectSuperCall()
func (c *OptimizeConstructorsTest) TestNoOptimize_sideEffectObjectSuperCall() {
	c.TestSame("class Top extends Object { constructor(a) { super(a++); }}")
}

// Java: OptimizeConstructorsTest.testNoOptimize_reassignedSuperclass()
func (c *OptimizeConstructorsTest) TestNoOptimize_reassignedSuperclass() {
	c.TestSame(
		"class Super {\n" +
			"  constructor() {}\n" +
			"}\n" +
			"class A extends Super {\n" +
			"  constructor() {\n" +
			"    super();\n" +
			"  }\n" +
			"}\n" +
			"Super = class {};\n" +
			"let a = new A();\n")
}

// Java: OptimizeConstructorsTest.testNoOptimize_reassignedConstructor()
func (c *OptimizeConstructorsTest) TestNoOptimize_reassignedConstructor() {
	c.TestSame(
		"class A {\n" +
			"  constructor() {}\n" +
			"}\n" +
			"A = class {};\n" +
			"let a = new A();\n")
}
