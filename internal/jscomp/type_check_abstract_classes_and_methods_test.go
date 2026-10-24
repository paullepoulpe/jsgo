// Port of com/google/javascript/jscomp/TypeCheckAbstractClassesAndMethodsTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/TypeCheckAbstractClassesAndMethodsTest.java)

/*
 * Copyright 2006 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"testing"
)

// TypeCheckAbstractClassesAndMethodsTest tests the type checker.
type TypeCheckAbstractClassesAndMethodsTest struct {
	typeCheckTests_noFixture
}

// Java: TypeCheckAbstractClassesAndMethodsTest()
func NewTypeCheckAbstractClassesAndMethodsTest() *TypeCheckAbstractClassesAndMethodsTest {
	c := &TypeCheckAbstractClassesAndMethodsTest{}
	return c
}

func TestTypeCheckAbstractClassesAndMethodsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewTypeCheckAbstractClassesAndMethodsTest() })
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodInAbstractClass()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodInAbstractClass() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @abstract @constructor */ var C = function() {};
/** @abstract */ C.prototype.foo = function() {};
`).
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodInAbstractEs6Class()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodInAbstractEs6Class() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @abstract */ class C {
  /** @abstract */ foo() {}
}
`).
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodInConcreteClass()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodInConcreteClass() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ var C = function() {};
/** @abstract */ C.prototype.foo = function() {};
`).
		AddDiagnostic(`Abstract methods can only appear in abstract classes. Please declare the class as @abstract
`).
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodInConcreteEs6Class()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodInConcreteEs6Class() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class C {
  /** @abstract */ foo() {}
}
`).
		AddDiagnostic(`Abstract methods can only appear in abstract classes. Please declare the class as @abstract
`).
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodInConcreteClassExtendingAbstractClass()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodInConcreteClassExtendingAbstractClass() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @abstract @constructor */ var A = function() {};
/** @constructor @extends {A} */ var B = function() {};
/** @abstract */ B.prototype.foo = function() {};
`).
		AddDiagnostic(`Abstract methods can only appear in abstract classes. Please declare the class as @abstract
`).
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodInConcreteEs6ClassExtendingAbstractEs6Class()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodInConcreteEs6ClassExtendingAbstractEs6Class() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @abstract */ class A {}
/** @extends {A} */ class B {
  /** @abstract */ foo() {}
}
`).
		AddDiagnostic("Abstract methods can only appear in abstract classes. Please declare the class as @abstract").
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testConcreteMethodOverridingAbstractMethod()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestConcreteMethodOverridingAbstractMethod() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @abstract @constructor */ var A = function() {};
/** @abstract */ A.prototype.foo = function() {};
/** @constructor @extends {A} */ var B = function() {};
/** @override */ B.prototype.foo = function() {};
`).
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testConcreteMethodOverridingAbstractMethodInEs6()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestConcreteMethodOverridingAbstractMethodInEs6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @abstract */ class A {
  /** @abstract*/ foo() {}
}
/** @extends {A} */ class B {
  /** @override */ foo() {}
}
`).
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testConcreteMethodInAbstractClass1()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestConcreteMethodInAbstractClass1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @abstract @constructor */ var A = function() {};
A.prototype.foo = function() {};
/** @constructor @extends {A} */ var B = function() {};
`).
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testConcreteMethodInAbstractEs6Class1()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestConcreteMethodInAbstractEs6Class1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @abstract */ class A {
  foo() {}
}
class B extends A {}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testConcreteMethodInAbstractClass2()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestConcreteMethodInAbstractClass2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(jutil.ToJavaString(CompilerTypeTestCase_CLOSURE_DEFS) + `/** @abstract @constructor */ var A = function() {};
A.prototype.foo = goog.abstractMethod;
/** @constructor @extends {A} */ var B = function() {};
`).
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodInInterface()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodInInterface() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ var I = function() {};
/** @abstract */ I.prototype.foo = function() {};
`).
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodInEs6Interface()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodInEs6Interface() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ class I {
  /** @abstract */ foo() {}
};
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodNotImplemented1()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodNotImplemented1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @abstract @constructor */ var A = function() {};
/** @abstract */ A.prototype.foo = function() {};
/** @constructor @extends {A} */ var B = function() {};
`).
		AddDiagnostic("property foo on abstract class A is not implemented by type B").
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodInEs6NotImplemented1()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodInEs6NotImplemented1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @abstract */ class A {
  /** @abstract */ foo() {}
}
class B extends A {}
`).
		AddDiagnostic("property foo on abstract class A is not implemented by type B").
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodNotImplemented2()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodNotImplemented2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @abstract @constructor */ var A = function() {};
/** @abstract */ A.prototype.foo = function() {};
/** @abstract */ A.prototype.bar = function() {};
/** @constructor @extends {A} */ var B = function() {};
/** @override */ B.prototype.foo = function() {};
`).
		AddDiagnostic("property bar on abstract class A is not implemented by type B").
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodInEs6NotImplemented2()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodInEs6NotImplemented2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @abstract */ class A {
  /** @abstract */ foo() {}
  /** @abstract */ bar() {}
}
class B extends A {
  /** @override */ foo() {}
}
`).
		AddDiagnostic("property bar on abstract class A is not implemented by type B").
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodNotImplemented3()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodNotImplemented3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @abstract @constructor */ var A = function() {};
/** @abstract */ A.prototype.foo = function() {};
/** @abstract @constructor @extends {A} */ var B = function() {};
/** @abstract @override */ B.prototype.foo = function() {};
/** @constructor @extends {B} */ var C = function() {};
`).
		AddDiagnostic("property foo on abstract class B is not implemented by type C").
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodInEs6NotImplemented3()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodInEs6NotImplemented3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @abstract */ class A {
  /** @abstract */ foo() {}
}
/** @abstract */ class B extends A {
  /** @abstract @override */ foo() {}
}
class C extends B {}
`).
		AddDiagnostic("property foo on abstract class B is not implemented by type C").
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodNotImplemented4()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodNotImplemented4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @abstract @constructor */ var A = function() {};
/** @abstract */ A.prototype.foo = function() {};
/** @abstract @constructor @extends {A} */ var B = function() {};
/** @constructor @extends {B} */ var C = function() {};
`).
		AddDiagnostic("property foo on abstract class A is not implemented by type C").
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodInEs6NotImplemented4()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodInEs6NotImplemented4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @abstract */ class A {
  /** @abstract */ foo() {}
}
/** @abstract */ class B extends A {}
class C extends B {}
`).
		AddDiagnostic("property foo on abstract class A is not implemented by type C").
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodNotImplemented5()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodNotImplemented5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ var I = function() {};
I.prototype.foo = function() {};
/** @abstract @constructor @implements {I} */ var A = function() {};
/** @abstract @override */ A.prototype.foo = function() {};
/** @constructor @extends {A} */ var B = function() {};
`).
		AddDiagnostic("property foo on abstract class A is not implemented by type B").
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodInEs6NotImplemented5()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodInEs6NotImplemented5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ class I {
  foo() {}
  bar() {} // Not overridden by abstract class
}
/** @abstract @implements {I} */ class A {
  /** @abstract @override */ foo() {}
}
class B extends A {}
`).
		AddDiagnostic("property bar on interface I is not implemented by type B").
		AddDiagnostic("property foo on abstract class A is not implemented by type B").
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodNotImplemented6()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodNotImplemented6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @abstract @constructor */ var A = function() {};
/** @abstract */ A.prototype.foo = function() {};
/** @constructor @extends {A} */ var B = function() {};
/** @override @type {number} */ B.prototype.foo;
`).
		AddDiagnostic("property foo on abstract class A is not implemented by type B").
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodImplemented1()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodImplemented1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @abstract @constructor */ var A = function() {};
/** @abstract */ A.prototype.foo = function() {};
/** @abstract */ A.prototype.bar = function() {};
/** @constructor @extends {A} */ var B = function() {};
/** @override */ B.prototype.foo = function() {};
/** @override */ B.prototype.bar = function() {};
/** @constructor @extends {B} */ var C = function() {};
`).
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodInEs6Implemented1()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodInEs6Implemented1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @abstract */ class A {
  /** @abstract */ foo() {}
  /** @abstract */ bar() {}
}
class B extends A {
  /** @override */ foo() {}
  /** @override */ bar() {}
}
class C extends B {}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodImplemented2()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodImplemented2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @abstract @constructor */ var A = function() {};
/** @abstract */ A.prototype.foo = function() {};
/** @abstract */ A.prototype.bar = function() {};
/** @abstract @constructor @extends {A} */ var B = function() {};
/** @override */ B.prototype.foo = function() {};
/** @constructor @extends {B} */ var C = function() {};
/** @override */ C.prototype.bar = function() {};
`).
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodInEs6Implemented2()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodInEs6Implemented2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @abstract */ class A {
  /** @abstract */ foo() {}
  /** @abstract */ bar() {}
}
/** @abstract */ class B extends A {
  /** @override */ foo() {}
}
class C extends B {
  /** @override */ bar() {}
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodInEs6NotImplemented_symbolProp()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodInEs6NotImplemented_symbolProp() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @abstract */ class A {
  /** @abstract */ [Symbol.iterator]() {}
}
class B extends A {}
`).
		AddDiagnostic("property Symbol.iterator on abstract class A is not implemented by type B").
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodHandling1()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodHandling1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {Function} */ var abstractFn = function() {};
abstractFn(1);
`).
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodHandling2()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodHandling2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var abstractFn = function() {};
abstractFn(1);
`).
		AddDiagnostic(`Function abstractFn: called with 1 argument(s). Function requires at least 0 argument(s) and no more than 0 argument(s).
`).
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodHandling3()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodHandling3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = {};
/** @type {Function} */ goog.abstractFn = function() {};
goog.abstractFn(1);
`).
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodHandling4()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodHandling4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = {};
goog.abstractFn = function() {};
goog.abstractFn(1);
`).
		AddDiagnostic("Function goog.abstractFn: called with 1 argument(s). Function requires at least 0 argument(s) and no more than 0 argument(s).").
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractFunctionHandling()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractFunctionHandling() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {!Function} */ var abstractFn = function() {};
// the type of 'f' will become 'Function'
/** @param {number} x */ var f = abstractFn;
f('x');
`).
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodHandling6()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodHandling6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = {};
/** @type {Function} */ goog.abstractFn = function() {};
/** @param {number} x */ goog.f = abstractFn;
goog.f('x');
`).
		AddDiagnostic(`actual parameter 1 of goog.f does not match formal parameter
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractSpread()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractSpread() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @abstract */
class X {
  /** @abstract */
  m1() {}

  m2() {
    return () => this.m1(...[]);
  }
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testOverriddenReturnOnAbstractClass()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestOverriddenReturnOnAbstractClass() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ function IFoo() {}
/** @return {*} */ IFoo.prototype.foo = function() {}
/** @constructor */ function Foo() {}
/** @return {string} */ Foo.prototype.foo = function() {}
/** @constructor @extends {Foo} */ function Bar() {}
/**
 * @constructor @abstract
 * @extends {Bar} @implements {IFoo}
 */
function Baz() {}
// Even there is a closer definition in IFoo, Foo should be still the source of truth.
/** @return {string} */
function test() { return (/** @type {Baz} */ (null)).foo(); }
`).
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testOverriddenReturnDoesntMatchOnAbstractClass()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestOverriddenReturnDoesntMatchOnAbstractClass() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ function IFoo() {}
/** @return {number} */ IFoo.prototype.foo = function() {}
/** @constructor */ function Foo() {}
/** @return {string} */ Foo.prototype.foo = function() {}
/** @constructor @extends {Foo} */ function Bar() {}
/**
 * @constructor @abstract
 * @extends {Bar} @implements {IFoo}
 */
function Baz() {}
/** @return {string} */
function test() { return (/** @type {Baz} */ (null)).foo(); }
`).
		AddDiagnostic(`mismatch of the foo property on type Baz and the type of the property it overrides from interface IFoo
original: function(this:IFoo): number
override: function(this:Foo): string
`).
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testOverriddenReturnDoesntMatchOnAbstractClass_symbolProp()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestOverriddenReturnDoesntMatchOnAbstractClass_symbolProp() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ function IFoo() {}
/** @return {number} */ IFoo.prototype[Symbol.iterator] = function() {}
/** @constructor */ function Foo() {}
/** @return {string} */ Foo.prototype[Symbol.iterator] = function() {}
/** @constructor @extends {Foo} */ function Bar() {}
/**
 * @constructor
 * @abstract
 * @extends {Bar}
 * @implements {IFoo}
 */
function Baz() {}
/** @return {string} */
function test() { return (/** @type {Baz} */ (null))[Symbol.iterator](); }
`).
		AddDiagnostic(`mismatch of the Symbol.iterator property on type Baz and the type of the property it overrides from interface IFoo
original: function(this:IFoo): number
override: function(this:Foo): string
`).
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodCall1()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodCall1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor @abstract */ var A = function() {};
/** @abstract */ A.prototype.foo = function() {};
/** @constructor @extends {A} */ var B = function() {};
B.superClass_ = A.prototype
/** @override */ B.prototype.foo = function() { B.superClass_.foo.call(this); };
`).
		AddDiagnostic("Abstract super method A.prototype.foo cannot be dereferenced").
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodCall2()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodCall2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const */ var ns = {};
/** @constructor @abstract */ ns.A = function() {};
/** @abstract */ ns.A.prototype.foo = function() {};
/** @constructor @extends {ns.A} */ ns.B = function() {};
ns.B.superClass_ = ns.A.prototype
/** @override */ ns.B.prototype.foo = function() {
  ns.B.superClass_.foo.call(this);
};
`).
		AddDiagnostic("Abstract super method ns.A.prototype.foo cannot be dereferenced").
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodCall3()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodCall3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor @abstract */ var A = function() {};
/** @abstract */ A.prototype.foo = function() {};
/** @constructor @extends {A} */ var B = function() {};
/** @override */ B.prototype.foo = function() { A.prototype.foo.call(this); };
`).
		AddDiagnostic("Abstract super method A.prototype.foo cannot be dereferenced").
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodCall4()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodCall4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const */ var ns = {};
/** @constructor @abstract */ ns.A = function() {};
ns.A.prototype.foo = function() {};
/** @constructor @extends {ns.A} */ ns.B = function() {};
ns.B.superClass_ = ns.A.prototype
/** @override */ ns.B.prototype.foo = function() {
  ns.B.superClass_.foo.call(this);
};
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodCall5()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodCall5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor @abstract */ var A = function() {};
A.prototype.foo = function() {};
/** @constructor @extends {A} */ var B = function() {};
/** @override */ B.prototype.foo = function() { A.prototype.foo.call(this); };
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodCall6()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodCall6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const */ var ns = {};
/** @constructor @abstract */ ns.A = function() {};
ns.A.prototype.foo = function() {};
ns.A.prototype.foo.bar = function() {};
/** @constructor @extends {ns.A} */ ns.B = function() {};
ns.B.superClass_ = ns.A.prototype
/** @override */ ns.B.prototype.foo = function() {
  ns.B.superClass_.foo.bar.call(this);
};
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodCall7()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodCall7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor @abstract */ var A = function() {};
A.prototype.foo = function() {};
A.prototype.foo.bar = function() {};
/** @constructor @extends {A} */ var B = function() {};
/** @override */ B.prototype.foo = function() { A.prototype.foo.bar.call(this); };
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodCall8()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodCall8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor @abstract */ var A = function() {};
A.prototype.foo = function() {};
/** @constructor @extends {A} */ var B = function() {};
/** @override */ B.prototype.foo = function() { A.prototype.foo['call'](this); };
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodCall9()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodCall9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @struct @constructor */ var A = function() {};
A.prototype.foo = function() {};
/** @struct @constructor @extends {A} */ var B = function() {};
/** @override */ B.prototype.foo = function() {
  (function() {
    return A.prototype.foo.call($jscomp$this);
  })();
};
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodCall10()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodCall10() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor @abstract */ var A = function() {};
/** @abstract */ A.prototype.foo = function() {};
A.prototype.foo.call(new Subtype);
`).
		AddDiagnostic("Abstract super method A.prototype.foo cannot be dereferenced").
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodCall11()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodCall11() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor @abstract */ function A() {};
/** @abstract */ A.prototype.foo = function() {};
/** @constructor @extends {A} */ function B() {};
/** @override */ B.prototype.foo = function() {};
var abstractMethod = A.prototype.foo;
abstractMethod.call(new B);
`).
		AddDiagnostic("Abstract super method A.prototype.foo cannot be dereferenced").
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodCall12()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodCall12() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor @abstract */ var A = function() {};
/** @abstract */ A.prototype.foo = function() {};
/** @constructor @extends {A} */ var B = function() {};
B.superClass_ = A.prototype
/** @override */ B.prototype.foo = function() { B.superClass_.foo.apply(this); };
`).
		AddDiagnostic("Abstract super method A.prototype.foo cannot be dereferenced").
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodCall13()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodCall13() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor @abstract */ var A = function() {};
/** @constructor @extends {A} */ var B = function() { A.call(this); };
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodCall_Indirect1()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodCall_Indirect1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor @abstract */ function A() {};
/** @abstract */ A.prototype.foo = function() {};
/** @constructor @extends {A} */ function B() {};
/** @override */ B.prototype.foo = function() {};
var abstractMethod = A.prototype.foo;
(0, abstractMethod).call(new B);
`).
		AddDiagnostic("Abstract super method A.prototype.foo cannot be dereferenced").
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodCall_Indirect2()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodCall_Indirect2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor @abstract */ function A() {};
/** @abstract */ A.prototype.foo = function() {};
/** @constructor @extends {A} */ function B() {};
/** @override */ B.prototype.foo = function() {};
var abstractMethod = A.prototype.foo;
(abstractMethod = abstractMethod).call(new B);
`).
		AddDiagnostic("Abstract super method A.prototype.foo cannot be dereferenced").
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testDefiningPropOnAbstractMethodForbidden()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestDefiningPropOnAbstractMethodForbidden() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor @abstract */ function A() {};
/** @abstract */ A.prototype.foo = function() {};
A.prototype.foo.callFirst = true;
`).
		AddDiagnostic("Abstract super method A.prototype.foo cannot be dereferenced").
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testPassingAbstractMethodAsArgForbidden()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestPassingAbstractMethodAsArgForbidden() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns("function externsFn(callback) {}").
		AddSource(`/** @constructor @abstract */ function A() {};
/** @abstract */ A.prototype.foo = function() {};
externsFn(A.prototype.foo);
`).
		AddDiagnostic("Abstract super method A.prototype.foo cannot be dereferenced").
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodCall_Es6Class()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodCall_Es6Class() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @abstract */
class Base {
  /** @abstract */
  foo() {}
  bar() {
    this.foo();
  }
}
class Sub extends Base {
  /** @override */
  foo() {}
  /** @override */
  bar() {
    this.foo();
  }
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodCall_Es6Class_prototype()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodCall_Es6Class_prototype() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @abstract */
class Base {
  /** @abstract */
  foo() {}
}
class Sub extends Base {
  /** @override */
  foo() {}
  bar() {
    Sub.prototype.foo();
  }
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodCall_Es6Class_prototype_warning()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodCall_Es6Class_prototype_warning() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @abstract */
class Base {
  /** @abstract */
  foo() {}
}
class Sub extends Base {
  /** @override */
  foo() {}
  bar() {
    Base.prototype.foo();
  }
}
`).
		AddDiagnostic("Abstract super method Base.prototype.foo cannot be dereferenced").
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodCall_Es6Class_abstractSubclass_warns()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodCall_Es6Class_abstractSubclass_warns() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @abstract */
class Base {
  /** @abstract */
  foo() {}
}
/** @abstract */
class Sub extends Base {
  bar() {
    Sub.prototype.foo();
  }
}
`).
		AddDiagnostic("Abstract super method Base.prototype.foo cannot be dereferenced").
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodCall_Es6Class_onAbstractSubclassPrototype_warns()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodCall_Es6Class_onAbstractSubclassPrototype_warns() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @abstract */
class Base {
  /** @abstract */
  foo() {}
}
/** @abstract */
class Sub extends Base {
  bar() {
    Base.prototype.foo();
  }
}
`).
		AddDiagnostic("Abstract super method Base.prototype.foo cannot be dereferenced").
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodCall_Es6Class_concreteSubclassMissingImplementation_warns()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodCall_Es6Class_concreteSubclassMissingImplementation_warns() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @abstract */
class Base {
  /** @abstract */
  foo() {}
}
class Sub extends Base {
  bar() {
    Sub.prototype.foo();
  }
}
`).
		IncludeDefaultExterns().
		AddDiagnostic("property foo on abstract class Base is not implemented by type Sub").
		AddDiagnostic("Abstract super method Base.prototype.foo cannot be dereferenced").
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodCall_Es6Class_concreteSubclassWithImplementation_noWarning()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodCall_Es6Class_concreteSubclassWithImplementation_noWarning() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @abstract */
class Base {
  /** @abstract */
  foo() {}
}
class Sub extends Base {
  /** @override */
  foo() {}
  bar() {
    Sub.prototype.foo();
  }
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodCall_NamespacedEs6Class_prototype_warns()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodCall_NamespacedEs6Class_prototype_warns() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`const ns = {};
/** @abstract */
ns.Base = class {
  /** @abstract */
  foo() {}
}
class Sub extends ns.Base {
  /** @override */
  foo() {}
  bar() {
    ns.Base.prototype.foo();
  }
}
`).
		AddDiagnostic("Abstract super method ns.Base.prototype.foo cannot be dereferenced").
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testNonAbstractMethodCall_Es6Class_prototype()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestNonAbstractMethodCall_Es6Class_prototype() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @abstract */
class Base {
  /** @abstract */
  foo() {}
  bar() {}
}
class Sub extends Base {
  /** @override */
  foo() {}
  /** @override */
  bar() {
    Base.prototype.bar();
  }
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testAbstractMethodCall_Es6ClassWithSpread()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestAbstractMethodCall_Es6ClassWithSpread() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddObject().AddArray().
		AddArguments().
		Build()).
		AddSource(`/** @abstract */
class Base {
  /** @abstract */
  foo() {}
}
class Sub extends Base {
  /** @override */
  foo() {}
  /** @param {!Array} arr */
  bar(arr) {
    this.foo.apply(this, [].concat(arr));
  }
}
`).
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testImplementedInterfacePropertiesShouldFailOnConflictForAbstractClasses()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestImplementedInterfacePropertiesShouldFailOnConflictForAbstractClasses() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */function Int0() {};
/** @interface */function Int1() {};
/** @type {number} */
Int0.prototype.foo;
/** @type {string} */
Int1.prototype.foo;
/** @constructor @abstract @implements {Int0} @implements {Int1} */
function Foo() {};
`).
		Run()
}

// Java: TypeCheckAbstractClassesAndMethodsTest.testImplementedInterfacePropertiesShouldFailOnConflictForAbstractClasses2()
func (c *TypeCheckAbstractClassesAndMethodsTest) TestImplementedInterfacePropertiesShouldFailOnConflictForAbstractClasses2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */function Int0() {};
/** @interface */function Int1() {};
/** @type {number} */
Int0.prototype.foo;
/** @type {string} */
Int1.prototype.foo;
/** @constructor @abstract @implements {Int0} */
function Foo() {};
/** @constructor @abstract @extends {Foo} @implements {Int1} */
function Zoo() {};
`).
		Run()
}
