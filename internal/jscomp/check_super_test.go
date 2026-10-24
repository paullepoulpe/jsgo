// Port of com/google/javascript/jscomp/CheckSuperTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/CheckSuperTest.java)

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

type CheckSuperTest struct {
	CompilerTestCaseBase
}

// Java: CheckSuperTest()
func NewCheckSuperTest() *CheckSuperTest {
	c := &CheckSuperTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestCheckSuperTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCheckSuperTest() })
}

// Java: CheckSuperTest.getProcessor(Compiler)
func (c *CheckSuperTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewCheckSuper(compiler)
}

// Java: CheckSuperTest.setUp()
func (c *CheckSuperTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	// AST validation will throw exceptions for the cases we want to report as errors here.
	c.DisableAstValidation()
}

// Java: CheckSuperTest.testInvalidUsage()
func (c *CheckSuperTest) TestInvalidUsage() {
	// 'super' may only appear directly within a CALL, a GETPROP, or a GETELEM.
	c.TestError("class Foo extends Bar { constructor() { super(); super; } }", CheckSuper_INVALID_SUPER_USAGE)
	c.TestError("class Foo extends Bar { foo() { super; } }", CheckSuper_INVALID_SUPER_USAGE)
	c.TestError("function f() { super; }", CheckSuper_INVALID_SUPER_USAGE)
	c.TestError("super;", CheckSuper_INVALID_SUPER_USAGE)
}

// Java: CheckSuperTest.testMissingSuper()
func (c *CheckSuperTest) TestMissingSuper() {
	c.TestError("class C extends D { constructor() {} }", CheckSuper_MISSING_CALL_TO_SUPER)
	// Without a `super()` call or a returned value, the constructor always throws, so
	// `missingCallToSuper` doesn't suppress this error.
	c.TestError(
		"class C extends D { /** @suppress {missingCallToSuper} */ constructor() {} }",
		CheckSuper_MISSING_CALL_TO_SUPER)
	c.Test(
		c.Srcs("class C extends D { constructor() { super.foo(); } }"),
		c.Error(CheckSuper_MISSING_CALL_TO_SUPER),
		c.Error(CheckSuper_SUPER_ACCESS_BEFORE_SUPER_CONSTRUCTOR))
}

// Java: CheckSuperTest.testFieldSuperCall()
func (c *CheckSuperTest) TestFieldSuperCall() {
	c.TestSame(
		`class Foo {
  x() {
    return 5;
  }
}
class Bar extends Foo {
  y = () => super.x()
}
`)

	c.TestError(
		`class A {}
class B extends A {
  x = super();
}
`,
		CheckSuper_INVALID_SUPER_CALL)
}

// Java: CheckSuperTest.testComputedSuperFields()
func (c *CheckSuperTest) TestComputedSuperFields() {
	c.TestError(
		`class A {
  x = 'fieldOne'
}
class B extends A {
  [super.x] = 2;
}
`,
		CheckSuper_INVALID_SUPER_ACCESS)

	c.TestError(
		`class A {
  x = 'fieldOne'
}
class B extends A {
  static [super.x] = 2;
}
`,
		CheckSuper_INVALID_SUPER_ACCESS)
}

// Java: CheckSuperTest.testMissingSuper_nestedClass()
func (c *CheckSuperTest) TestMissingSuper_nestedClass() {
	// Note that super is only called for anonymous class "E", not C.
	c.TestError(
		`class C extends D {
  constructor() {
    const E = class extends D { constructor() { super(); } };
  }
}
`,
		CheckSuper_MISSING_CALL_TO_SUPER)
}

// Java: CheckSuperTest.testNoWarning_derivedConstructorCompletesWithUndefinedBeforeSuperCall()
func (c *CheckSuperTest) TestNoWarning_derivedConstructorCompletesWithUndefinedBeforeSuperCall() {
	// Won't fix (b/573571176): a derived class constructor that completes with `undefined` before
	// calling `super()` throws a ReferenceError natively, but for a class without instance fields
	// or private members, `MISSING_CALL_TO_SUPER` is only reported for a constructor that has
	// neither a `super()` call nor a `return` with a value. We decided not to report these
	// constructors, since the compiler isn't required to preserve errors in code that is already
	// broken.
	c.TestNoWarning("class C extends D { constructor() { return undefined; } }")
	c.TestNoWarning("class C extends D { constructor(w) { if (w) return; super(); } }")
	c.TestNoWarning("class C extends D { constructor(w) { if (w) return {}; } }")
}

// TestNoWarning_ is the @Test method testNoWarning(); the trailing _ keeps it
// apart from CompilerTestCase's testNoWarning(...) overloads, which Go cannot
// overload (Go-only: language-adapter).
//
// Java: CheckSuperTest.testNoWarning()
func (c *CheckSuperTest) TestNoWarning_() {
	c.TestSame("class C extends D { constructor() { super(); } }")
	c.TestSame("class C { constructor() {} }")
	c.TestSame("class C extends D {}")
}

// just verifying that import.meta doesn't cause problems with super
// Java: CheckSuperTest.testImportMeta()
func (c *CheckSuperTest) TestImportMeta() {
	c.TestSame("class C extends D { constructor() { super(import.meta); } }")
}

// Java: CheckSuperTest.testThisBeforeSuper()
func (c *CheckSuperTest) TestThisBeforeSuper() {
	c.TestError("class C extends D { constructor() { this.foo(); super(); } }", CheckSuper_THIS_BEFORE_SUPER)
}

// Java: CheckSuperTest.testSuperPropAccessBeforeSuperCall()
func (c *CheckSuperTest) TestSuperPropAccessBeforeSuperCall() {
	c.TestError(
		"class C extends D { constructor() { super.foo(); super(); } }",
		CheckSuper_SUPER_ACCESS_BEFORE_SUPER_CONSTRUCTOR)
}

// Java: CheckSuperTest.testThisAndSuperPropAccessBeforeSuperCall()
func (c *CheckSuperTest) TestThisAndSuperPropAccessBeforeSuperCall() {
	c.Test(
		c.Srcs("class C extends D { constructor() { this.foo(); super.foo(); super(); } }"),
		c.Error(CheckSuper_SUPER_ACCESS_BEFORE_SUPER_CONSTRUCTOR),
		c.Error(CheckSuper_THIS_BEFORE_SUPER))
}

// Java: CheckSuperTest.testThisAndSuperPropAccessBeforeSuperCall_inSuperConstructorArgs()
func (c *CheckSuperTest) TestThisAndSuperPropAccessBeforeSuperCall_inSuperConstructorArgs() {
	c.Test(
		c.Srcs("class C extends D { constructor() { super(this.foo(), super.foo()); } }"),
		c.Error(CheckSuper_SUPER_ACCESS_BEFORE_SUPER_CONSTRUCTOR),
		c.Error(CheckSuper_THIS_BEFORE_SUPER))
}

// Java: CheckSuperTest.testThisBeforeSuper_inParameterDefault()
func (c *CheckSuperTest) TestThisBeforeSuper_inParameterDefault() {
	// Parameter defaults are evaluated before the body, so before super() binds `this`.
	c.TestError("class C extends D { constructor(x = this.foo) { super(); } }", CheckSuper_THIS_BEFORE_SUPER)
}

// We could require that the super() call is the first statement in the constructor, except that
// doing so breaks J2CL-compiled code, which needs to do the static initialization for the class
// before anything else.
// Java: CheckSuperTest.testNoWarning_J2CL()
func (c *CheckSuperTest) TestNoWarning_J2CL() {
	c.TestSame("class C extends D { constructor() { C.init(); super(); } }")
}

// Referencing this within a function before calling super is acceptable at runtime so long as
// it's never executed before the super() has returned. It's also acceptable since this might
// be bound to something other than the class being constructed.
// Java: CheckSuperTest.testNoWarning_thisWithinFunction()
func (c *CheckSuperTest) TestNoWarning_thisWithinFunction() {
	c.TestSame("class C extends D { constructor() { const c = () => this; super(); } }")
	c.TestSame("class C extends D { constructor() { const c = function() { this; }; super(); } }")
}

// Java: CheckSuperTest.testError_superPropWithinConstructorArrowFunction()
func (c *CheckSuperTest) TestError_superPropWithinConstructorArrowFunction() {
	c.TestError(
		"class C extends D { constructor() { const c = () => super.foo; super(); } }",
		CheckSuper_SUPER_ACCESS_BEFORE_SUPER_CONSTRUCTOR)
}

// Java: CheckSuperTest.testOutsideClass()
func (c *CheckSuperTest) TestOutsideClass() {
	c.TestError("var i = super();", CheckSuper_INVALID_SUPER_CALL)
	c.TestError("var i = super.i;", CheckSuper_INVALID_SUPER_ACCESS)
	c.TestError("var i = super[x];", CheckSuper_INVALID_SUPER_ACCESS)
}

// Java: CheckSuperTest.testInOrdinaryFunction()
func (c *CheckSuperTest) TestInOrdinaryFunction() {
	c.TestError("function f() { super(); }", CheckSuper_INVALID_SUPER_CALL)
	c.TestError("function f() { var i = super.i; }", CheckSuper_INVALID_SUPER_ACCESS)
	c.TestError("function f() { var i = super[x]; }", CheckSuper_INVALID_SUPER_ACCESS)
}

// Java: CheckSuperTest.testSuperAccessInNestedOrdinaryFunction()
func (c *CheckSuperTest) TestSuperAccessInNestedOrdinaryFunction() {
	c.TestError(
		`class C {
  method() {}
}
class Sub extends C {
  method() {
    (function() { super.method() }).call(this);
  }
}
`,
		CheckSuper_INVALID_SUPER_ACCESS)
}

// Java: CheckSuperTest.testSuperAccessInNestedArrowFunction()
func (c *CheckSuperTest) TestSuperAccessInNestedArrowFunction() {
	c.TestSame(
		`class C {
  method() {}
}
class Sub extends C {
  method() {
    (() => { super.method() })();
  }
}
`)
}

// Java: CheckSuperTest.testSuperAccessInNestedOrdinaryFunctionInGetter()
func (c *CheckSuperTest) TestSuperAccessInNestedOrdinaryFunctionInGetter() {
	c.TestError(
		`class C {
  get value() { return 1; }
}
class Sub extends C {
  get value() {
    return (function() { super.value }).call(this);
  }
}
`,
		CheckSuper_INVALID_SUPER_ACCESS)
}

// Java: CheckSuperTest.testSuperAccessInNestedArrowFunctionInGetter()
func (c *CheckSuperTest) TestSuperAccessInNestedArrowFunctionInGetter() {
	c.TestSame(
		`class C {
  get value() { return 1; }
}
class Sub extends C {
  get value() {
    return (() => super.value)();
  }
}
`)
}

// Java: CheckSuperTest.testSuperCallInNestedOrdinaryFunction()
func (c *CheckSuperTest) TestSuperCallInNestedOrdinaryFunction() {
	c.TestError(
		`class C {
  constructor() {}
}
class Sub extends C {
  constructor() {
    super(); // avoid getting missing super call warning
    (function() { super() }).call(this);
  }
}
`,
		CheckSuper_INVALID_SUPER_CALL)
}

// Java: CheckSuperTest.testSuperCallInNestedObjectLiteralMethod()
func (c *CheckSuperTest) TestSuperCallInNestedObjectLiteralMethod() {
	c.TestError(
		`class C {
  constructor() {}
}
class Sub extends C {
  constructor() {
    super(); // avoid getting missing super call warning
    this.objLit = {
// Object literals don't have special constructor methods, so calling super()
// makes no sense. We nest this one inside of a method named constructor
// to make sure that doesn't fool our logic.
      constructor() { super(); }
    };
  }
}
`,
		CheckSuper_INVALID_SUPER_CALL)
}

// Java: CheckSuperTest.testObjLitMethodCanAccessSuperProps()
func (c *CheckSuperTest) TestObjLitMethodCanAccessSuperProps() {
	c.TestSame(
		`const objLit = {
// Object literals don't have special constructor methods, so calling super()
// makes no sense. We nest this one inside of a method named constructor
// to make sure that doesn't fool our logic.
      method() { super.toString(); }
};
`)
}

// Java: CheckSuperTest.testObjLitPropDefinedWithArrowFunctionCannotAccessSuperProps()
func (c *CheckSuperTest) TestObjLitPropDefinedWithArrowFunctionCannotAccessSuperProps() {
	c.TestError(
		"const objLit = {\n// This is not a proper object literal method, so it cannot refer to properties using\n// `super`\n      method: () => { super.toString(); }\n};\n",
		CheckSuper_INVALID_SUPER_ACCESS)
}

// Java: CheckSuperTest.testObjLitDoesNotFoolArrowFunctionAccessLogicInMethod()
func (c *CheckSuperTest) TestObjLitDoesNotFoolArrowFunctionAccessLogicInMethod() {
	c.TestSame(
		"class C {\n  method() {\n    this.objLit = {\n// This is not a proper object literal method, but the arrow function is wrapped in\n// a class method, so reference to `super` properties is allowed here.\n      method: () => { super.toString(); }\n    };\n  }\n}\n")
}

// Java: CheckSuperTest.testSuperCallInNestedArrowFunction()
func (c *CheckSuperTest) TestSuperCallInNestedArrowFunction() {
	c.TestError(
		`class C {
  constructor() {}
}
class Sub extends C {
  constructor() {
    (() => { super() })();
  }
}
`,
		CheckSuper_SUPER_CALL_IN_ARROW)
}

// Java: CheckSuperTest.testCallWithNoBaseClass()
func (c *CheckSuperTest) TestCallWithNoBaseClass() {
	c.TestError("class C { constructor() { super(); }}", CheckSuper_INVALID_SUPER_CALL)
	c.TestError("class C { constructor() { super(1); }}", CheckSuper_INVALID_SUPER_CALL)
}

// Java: CheckSuperTest.testCallInConstructor()
func (c *CheckSuperTest) TestCallInConstructor() {
	c.TestSame("class C extends D { constructor() { super(); }}")
	c.TestSame("class C extends D { constructor() { super(1); }}")
}

// Java: CheckSuperTest.testNestedCallInConstructor()
func (c *CheckSuperTest) TestNestedCallInConstructor() {
	// Technically it is legal to call `super()` in an arrow function within a constructor,
	// but doing so makes it harder to tell that the constructor is calling `super()` correctly and
	// there's really no good reason to do it. Such code should be refactored.
	c.TestError("class C extends D { constructor() { (()=>{ super(); })(); }}", CheckSuper_SUPER_CALL_IN_ARROW)
}

// Java: CheckSuperTest.testOutOfOrderNestedCallInConstructor()
func (c *CheckSuperTest) TestOutOfOrderNestedCallInConstructor() {
	c.Test(
		c.Srcs("class C extends D { constructor() { this.foo(); (()=>{ super(); })(); }}"),
		c.Error(CheckSuper_SUPER_CALL_IN_ARROW),
		c.Error(CheckSuper_THIS_BEFORE_SUPER))
}

// Java: CheckSuperTest.testCallInMethod()
func (c *CheckSuperTest) TestCallInMethod() {
	c.Test(
		c.Srcs("class C extends D { foo() { super(); }}"),
		c.Error(CheckSuper_INVALID_SUPER_CALL_WITH_SUGGESTION).
			WithMessage("super() not allowed here. Did you mean super.foo?"))
	c.Test(
		c.Srcs("class C extends D { foo() { super(1); }}"),
		c.Error(CheckSuper_INVALID_SUPER_CALL_WITH_SUGGESTION).
			WithMessage("super() not allowed here. Did you mean super.foo?"))
	c.Test(
		c.Srcs("class C extends D { static foo() { super(1); }}"),
		c.Error(CheckSuper_INVALID_SUPER_CALL_WITH_SUGGESTION).
			WithMessage("super() not allowed here. Did you mean super.foo?"))
	c.TestError("class C { static foo() { super(); }}", CheckSuper_INVALID_SUPER_CALL_WITH_SUGGESTION)
	c.TestError(
		"class C extends D { foo() { (()=>{ super(); })(); }}", CheckSuper_INVALID_SUPER_CALL_WITH_SUGGESTION)
	c.TestError(
		"class C extends D { static foo() { (()=>{ super(); })(); }}",
		CheckSuper_INVALID_SUPER_CALL_WITH_SUGGESTION)
}

// Java: CheckSuperTest.testSuperCallInStaticBlock()
func (c *CheckSuperTest) TestSuperCallInStaticBlock() {
	c.TestError("class C extends D {static {super();}}", CheckSuper_INVALID_SUPER_CALL)
	c.TestError("class C extends D {static {super(1);}}", CheckSuper_INVALID_SUPER_CALL)
	c.TestError("class C extends D { static { (()=>{ super(); })(); }}", CheckSuper_INVALID_SUPER_CALL)
}

// Java: CheckSuperTest.testPropertyInStaticBlock()
func (c *CheckSuperTest) TestPropertyInStaticBlock() {
	c.TestSame("class C extends D { static { super.foo(); }}")
	c.TestSame("class C extends D { static { super.foo(1); }}")

	// TODO(tbreisacher): Consider warning for this. It's valid but likely indicates a mistake.
	c.TestSame("class C extends D { static { var x = super.bar; }}")
	c.TestSame("class C extends D { static { var x = super[y]; }}")
	c.TestSame("class C extends D { static { var x = super.bar(); }}")

	c.TestSame("class C extends D { static { this[super.x] = super.bar; }}")
	c.TestSame("class C extends D { static { this[super.x] = super[y]; }}")
	c.TestSame("class C extends D { static { this[super.x] = super.bar(); }}")
}

// Java: CheckSuperTest.testPropertyInMethod()
func (c *CheckSuperTest) TestPropertyInMethod() {
	c.TestSame("class C extends D { foo() { super.foo(); }}")
	c.TestSame("class C extends D { foo() { super.foo(1); }}")

	// TODO(tbreisacher): Consider warning for this. It's valid but likely indicates a mistake.
	c.TestSame("class C extends D { foo() { var x = super.bar; }}")
	c.TestSame("class C extends D { foo() { var x = super[y]; }}")
	c.TestSame("class C extends D { foo() { var x = super.bar(); }}")

	c.TestSame("class C extends D { foo() { this[super.x] = super.bar; }}")
	c.TestSame("class C extends D { foo() { this[super.x] = super[y]; }}")
	c.TestSame("class C extends D { foo() { this[super.x] = super.bar(); }}")
}

// Java: CheckSuperTest.testPropertyInStaticMethod()
func (c *CheckSuperTest) TestPropertyInStaticMethod() {
	c.TestSame("class C extends D { static foo() { super.foo(); }}")
	c.TestSame("class C extends D { static foo() { super.foo(1); }}")

	// TODO(tbreisacher): Consider warning for this. It's valid but likely indicates a mistake.
	c.TestSame("class C extends D { static foo() { var x = super.bar; }}")
	c.TestSame("class C extends D { static foo() { var x = super[y]; }}")
	c.TestSame("class C extends D { static foo() { var x = super.bar(); }}")

	c.TestSame("class C extends D { static foo() { this[super.x] = super.bar; }}")
	c.TestSame("class C extends D { static foo() { this[super.x] = super[y]; }}")
	c.TestSame("class C extends D { static foo() { this[super.x] = super.bar(); }}")
}

// Java: CheckSuperTest.testNoWarning_withExplicitReturnInConstructor()
func (c *CheckSuperTest) TestNoWarning_withExplicitReturnInConstructor() {
	// TODO(b/573596312): Require `super()` unless `@suppress {missingCallToSuper}` is used.
	c.TestNoWarning("class C extends D { constructor() { return {}; } }")
	c.TestNoWarning("class C extends D { constructor() { return f(); } }")
	// the following cases will error at runtime but we don't bother checking them.
	c.TestNoWarning("class C extends D { constructor() { return 3; } }")
	c.TestNoWarning("class C extends D { constructor() { if (false) return {}; } }")
}

// Java: CheckSuperTest.testReturnWithoutSuperCall()
func (c *CheckSuperTest) TestReturnWithoutSuperCall() {
	// TODO(b/573596312): A derived class constructor should be required to call `super()` as a
	// root-level statement, before any `return`, even when it returns a value. When class fields
	// are transpiled, they are installed at the start of a constructor that never calls `super()`,
	// where reading `this` throws a ReferenceError.
	c.TestNoWarning("class C extends D { x = 1; constructor() { return {}; } }")
	c.TestNoWarning("class C extends D { [k] = 1; constructor() { return {}; } }")
	c.TestNoWarning("class C extends D { #x = 1; constructor() { return {}; } }")
	c.TestNoWarning("class C extends D { #m() {} constructor() { return {}; } }")
	c.TestNoWarning("class C extends D { static x = 1; constructor() { return {}; } }")
	c.TestNoWarning(
		"class C extends D { constructor(a) { if (a) { super(); } else { return {}; } } }")
	c.TestNoWarning("class C extends D { constructor(a) { if (a) return {}; super(); } }")
}

// Java: CheckSuperTest.testNoWarning_rootLevelSuperCall()
func (c *CheckSuperTest) TestNoWarning_rootLevelSuperCall() {
	// `super()` is a root-level statement before any `return`, so every path that returns calls it.
	c.TestNoWarning(
		"class C extends D { x = 1; constructor(a) { if (!a) throw new Error(); super(); } }")
	c.TestNoWarning(
		"class C extends D { x = 1; constructor() { super(); return new Proxy(this, {}); } }")
}

// Java: CheckSuperTest.testWarning_withInvalidReturnInConstructor()
func (c *CheckSuperTest) TestWarning_withInvalidReturnInConstructor() {
	// empty return
	c.TestError("class C extends D { constructor() { return; } }", CheckSuper_MISSING_CALL_TO_SUPER)

	// return in arrow function
	c.TestError(
		"class C extends D { constructor() { () => { return {}; } } }", CheckSuper_MISSING_CALL_TO_SUPER)
}

// Java: CheckSuperTest.testInvalidThisReference_withExplicitReturnInConstructor()
func (c *CheckSuperTest) TestInvalidThisReference_withExplicitReturnInConstructor() {
	c.TestError(
		"class C extends D { constructor() { this.x = 3; super(); return {}; } }",
		CheckSuper_THIS_BEFORE_SUPER)
}

// Java: CheckSuperTest.testThisReference_withExplicitReturnInConstructor()
func (c *CheckSuperTest) TestThisReference_withExplicitReturnInConstructor() {
	// You don't have to call `super()` when the constructor returns a value, but then you also
	// cannot have references to `this` either.
	c.TestError("class C extends D { constructor() { this.x = 3; return {}; } }", CheckSuper_THIS_BEFORE_SUPER)
}

// Java: CheckSuperTest.testPropertyInConstructor()
func (c *CheckSuperTest) TestPropertyInConstructor() {
	// TODO(sdh): See note in testPropertyInMethod - these are valid but questionable.
	c.TestSame("class C extends D { constructor() { super(); super.foo(); }}")
	c.TestSame("class C extends D { constructor() { super(); super.foo(1); }}")
	c.TestSame("class C extends D { constructor() { super(); var x = super.bar; }}")
}

// Java: CheckSuperTest.testNestedProperty()
func (c *CheckSuperTest) TestNestedProperty() {
	c.TestSame("class C extends D { constructor() { super(); (()=>{ super.foo(); })(); }}")
	c.TestSame("class C extends D { foo() { (()=>{ super.foo(); })(); }}")
	c.TestSame("class C extends D { static foo() { (()=>{ super.foo(); })(); }}")
	c.TestSame("class C extends D { static { (()=>{ super.foo(); })(); }}")
}

// Java: CheckSuperTest.testPropertyNoBaseClass()
func (c *CheckSuperTest) TestPropertyNoBaseClass() {
	// TODO(sdh): See note in testPropertyInMethod - these are valid but questionable.
	c.TestSame("class C { constructor() { super.foo(); }}")
	c.TestSame("class C { foo() { super.foo(1); }}")
	c.TestSame("class C { foo() { super.bar(); }}")
	c.TestSame("class C { foo() { super.baz; }}")
	c.TestSame("class C { foo() { super[x](); }}")
}
