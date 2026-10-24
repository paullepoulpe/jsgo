// Port of com/google/javascript/jscomp/CheckAccessControlsTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/CheckAccessControlsTest.java)

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

// CheckAccessControlsTest tests CheckAccessControls.
//
// This file is forked from CheckAccessControlsOldSyntaxTest because nearly
// all cases there use the class syntax from before the `class` keyword was
// added to the language. Test cases using ES5-style classes with the
// `@constructor`, `@interface`, or `@record` annotations should remain in
// that file.
type CheckAccessControlsTest struct {
	CompilerTestCaseBase
}

// Java: CheckAccessControlsTest()
func NewCheckAccessControlsTest() *CheckAccessControlsTest {
	c := &CheckAccessControlsTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, CompilerTypeTestCase_DEFAULT_EXTERNS)
	return c
}

func TestCheckAccessControlsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCheckAccessControlsTest() })
}

// Java: CheckAccessControlsTest.setUp()
func (c *CheckAccessControlsTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableTypeCheck()
	c.EnableParseTypeInfo()
	c.EnableClosurePass()
	c.EnableRewriteClosureCode()
	c.EnableCreateModuleMap()
}

// Java: CheckAccessControlsTest.getProcessor(Compiler)
func (c *CheckAccessControlsTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewCheckAccessControls(compiler)
}

// Java: CheckAccessControlsTest.getOptions()
func (c *CheckAccessControlsTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetWarningLevel(DiagnosticGroups_ACCESS_CONTROLS, CheckLevel_ERROR)
	options.SetWarningLevel(DiagnosticGroups_DEPRECATED, CheckLevel_ERROR)
	options.SetWarningLevel(DiagnosticGroups_CONSTANT_PROPERTY, CheckLevel_ERROR)
	return options
}

// Java: CheckAccessControlsTest.testPrivateConstructor_propertyAccess()
func (c *CheckAccessControlsTest) TestPrivateConstructor_propertyAccess() {
	// AccessControlUtils handles visibility modifiers.
	c.TestError(
		c.Srcs("class Foo { /** @private */ constructor() { this.x = 5; } }", "(new Foo()).x;"),
		CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS)
}

// Java: CheckAccessControlsTest.testOverrideAllowed()
func (c *CheckAccessControlsTest) TestOverrideAllowed() {
	c.TestSame(
		`class Foo {
  /** @final*/ x;
  constructor() {
    this.x = 5;
  }
}
`)
}

// Java: CheckAccessControlsTest.testMultipleFieldsOverrideAllowed()
func (c *CheckAccessControlsTest) TestMultipleFieldsOverrideAllowed() {
	c.TestSame(
		`class Bar {
  /** @final*/ a;
  /** @final */ b;
  constructor() {
    this.a = 1;
    this.b = 2;
  }
}
`)
}

// Java: CheckAccessControlsTest.testFieldOverrideNotAllowed()
func (c *CheckAccessControlsTest) TestFieldOverrideNotAllowed() {
	c.TestError(
		`class Bar {
  /** @final*/ a;
  constructor() {
    this.a = 1;
    this.a = 2;
  }
}
`,
		CheckAccessControls_FINAL_PROPERTY_OVERRIDDEN)
	c.TestError(
		`class Bar {
  /** @final*/ a = 5;
  constructor() {
    this.a = 1;
  }
}
`,
		CheckAccessControls_FINAL_PROPERTY_OVERRIDDEN)
}

// Java: CheckAccessControlsTest.testThisAssignmentInConstructor()
func (c *CheckAccessControlsTest) TestThisAssignmentInConstructor() {
	c.TestSame(
		`class Foo {
  constructor() {
    /** @const {number} */ this.x;
    this.x = 4;
  }
}
`)
}

// Java: CheckAccessControlsTest.testWarningInNormalClass()
func (c *CheckAccessControlsTest) TestWarningInNormalClass() {
	c.Test(
		c.Srcs(
			`/** @deprecated FooBar */ function f() {}

var Foo = class {
  bar() { f(); }
};
`),
		c.DeprecatedName("Variable f has been deprecated: FooBar"))
}

// Java: CheckAccessControlsTest.testWarningForProperty1()
func (c *CheckAccessControlsTest) TestWarningForProperty1() {
	c.Test(
		c.Srcs(
			`class Foo {
  baz() { alert((new Foo()).bar); }
}
/** @deprecated A property is bad */ Foo.prototype.bar = 3;
`),
		c.DeprecatedProp("Property bar of type Foo has been deprecated: A property is bad"))
}

// Java: CheckAccessControlsTest.testWarningForProperty2()
func (c *CheckAccessControlsTest) TestWarningForProperty2() {
	c.Test(
		c.Srcs(
			`class Foo {
  baz() { alert(this.bar); }
}
/** @deprecated Zee prop, it is deprecated! */ Foo.prototype.bar = 3;
`),
		c.DeprecatedProp(
			"Property bar of type Foo has been deprecated: Zee prop, it is deprecated!"))
}

// Java: CheckAccessControlsTest.testWarningForDeprecatedClass()
func (c *CheckAccessControlsTest) TestWarningForDeprecatedClass() {
	c.Test(
		c.Srcs(
			`/** @deprecated Use the class 'Bar' */
class Foo {}

function f() { new Foo(); }
`),
		c.DeprecatedClass("Class Foo has been deprecated: Use the class 'Bar'"))
}

// Java: CheckAccessControlsTest.testWarningForDeprecatedClassNoReason()
func (c *CheckAccessControlsTest) TestWarningForDeprecatedClassNoReason() {
	c.Test(
		c.Srcs(
			`/** @deprecated */
class Foo {}

function f() { new Foo(); }
`),
		c.Error(CheckAccessControls_DEPRECATED_CLASS))
}

// Java: CheckAccessControlsTest.testNoWarningForDeprecatedClassInstance()
func (c *CheckAccessControlsTest) TestNoWarningForDeprecatedClassInstance() {
	c.Test(
		c.Srcs(
			`/** @deprecated */
class Foo {}

/** @param {Foo} x */
function f(x) { return x; }
`))
}

// Java: CheckAccessControlsTest.testNoWarningForDeprecatedSuperClass()
func (c *CheckAccessControlsTest) TestNoWarningForDeprecatedSuperClass() {
	c.TestNoWarning(
		c.Srcs(
			`/** @deprecated Superclass to the rescue! */
class Foo {}

class SubFoo extends Foo {}

function f() { new SubFoo(); }
`))
}

// Java: CheckAccessControlsTest.testNoWarningForDeprecatedSuperClassOnNamespace()
func (c *CheckAccessControlsTest) TestNoWarningForDeprecatedSuperClassOnNamespace() {
	c.TestNoWarning(
		c.Srcs(
			`/** @deprecated Its only weakness is Kryptoclass */
class Foo {}

/** @const */
var namespace = {};

namespace.SubFoo = class extends Foo {};

function f() { new namespace.SubFoo(); }
`))
}

// Java: CheckAccessControlsTest.testWarningForPrototypeProperty()
func (c *CheckAccessControlsTest) TestWarningForPrototypeProperty() {
	c.Test(
		c.Srcs(
			`class Foo {
  baz() { alert(Foo.prototype.bar); }
}

/** @deprecated It is now in production, use that one */ Foo.prototype.bar = 3;
`),
		c.DeprecatedProp(
			"Property bar of type Foo.prototype has been deprecated:"+
				" It is now in production, use that one"))
}

// Java: CheckAccessControlsTest.testNoWarningForNumbers()
func (c *CheckAccessControlsTest) TestNoWarningForNumbers() {
	c.Test(
		c.Srcs(
			`class Foo {
  baz() { alert(3); }
}

/** @deprecated */ Foo.prototype.bar = 3;
`))
}

// Java: CheckAccessControlsTest.testWarningForMethod1()
func (c *CheckAccessControlsTest) TestWarningForMethod1() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @deprecated There is a madness to this method */
  bar() {}

  baz() { this.bar(); }
}
`),
		c.DeprecatedProp(
			"Property bar of type Foo has been deprecated: There is a madness to this method"))
}

// Java: CheckAccessControlsTest.testWarningForMethod2()
func (c *CheckAccessControlsTest) TestWarningForMethod2() {
	c.Test(
		c.Srcs(
			`class Foo {
  baz() { this.bar(); }
}

/** @deprecated Stop the ringing! */
Foo.prototype.bar;
`),
		c.DeprecatedProp("Property bar of type Foo has been deprecated: Stop the ringing!"))
}

// Java: CheckAccessControlsTest.testWarningOnClassField_deprecated()
func (c *CheckAccessControlsTest) TestWarningOnClassField_deprecated() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @deprecated There is a madness to this method */
  bar = 2;
  baz() { return this.bar; }
}
`),
		c.DeprecatedProp(
			"Property bar of type Foo has been deprecated: There is a madness to this method"))
}

// Java: CheckAccessControlsTest.testNoWarningOnClassField_deprecated()
func (c *CheckAccessControlsTest) TestNoWarningOnClassField_deprecated() {
	c.Test(
		c.Srcs(
			`class Foo {
  /**
   * @type {?}
   * @deprecated Use something else.
   */
  x = 2;
}
`))
}

// Java: CheckAccessControlsTest.testNoWarningOnClassField_deprecated1()
func (c *CheckAccessControlsTest) TestNoWarningOnClassField_deprecated1() {
	// TODO(b/239747805): Should throw deprecatedProp warning
	c.Test(
		c.Srcs(
			`class Foo {
  /**
   * @type {number}
   * @deprecated No.
   */
  static x=1;
}
Foo.x;
`))
}

// Java: CheckAccessControlsTest.testNoWarningOnClassField_deprecated2()
func (c *CheckAccessControlsTest) TestNoWarningOnClassField_deprecated2() {
	// TODO(b/239747805): Should say "Property x of type number" not "super"
	c.Test(
		c.Srcs(
			`class Foo {
  /**
   * @type {number}
   * @deprecated No.
   */
  static x = 1;
}
class Bar extends Foo {
  static y = super.x;
}
`),
		c.DeprecatedProp("Property x of type super has been deprecated: No."))
}

// Java: CheckAccessControlsTest.testNoWarningInDeprecatedClass()
func (c *CheckAccessControlsTest) TestNoWarningInDeprecatedClass() {
	c.Test(
		c.Srcs(
			`/** @deprecated */
function f() {}

/** @deprecated */
var Foo = class {
  bar() { f(); }
};
`))
}

// Java: CheckAccessControlsTest.testNoWarningOnDeclaration()
func (c *CheckAccessControlsTest) TestNoWarningOnDeclaration() {
	c.Test(
		c.Srcs(
			`class Foo {
  constructor() {
    /**
     * @type {number}
     * @deprecated Use something else.
     */
    this.code;
  }
}
`))
}

// Java: CheckAccessControlsTest.testNoWarningInDeprecatedClass2()
func (c *CheckAccessControlsTest) TestNoWarningInDeprecatedClass2() {
	c.Test(
		c.Srcs(
			`/** @deprecated */
function f() {}

/** @deprecated */
var Foo = class {
  static bar() { f(); }
};
`))
}

// Java: CheckAccessControlsTest.testNoWarningInDeprecatedStaticMethod()
func (c *CheckAccessControlsTest) TestNoWarningInDeprecatedStaticMethod() {
	c.Test(
		c.Srcs(
			`/** @deprecated */
function f() {}

var Foo = class {
  /** @deprecated */
  static bar() { f(); }
};
`))
}

// Java: CheckAccessControlsTest.testWarningInStaticMethod()
func (c *CheckAccessControlsTest) TestWarningInStaticMethod() {
	c.Test(
		c.Srcs(
			`/** @deprecated crazy! */
function f() {}

var Foo = class {
  static bar() { f(); }
};
`),
		c.DeprecatedName("Variable f has been deprecated: crazy!"))
}

// Java: CheckAccessControlsTest.testWarningInClassStaticBlock()
func (c *CheckAccessControlsTest) TestWarningInClassStaticBlock() {
	c.Test(
		c.Srcs(
			`/** @deprecated Oh No! */
function f() {}

class Foo {
  static { f(); }
}
`),
		c.DeprecatedName("Variable f has been deprecated: Oh No!"))
}

// Java: CheckAccessControlsTest.testWarningInClassStaticBlock1()
func (c *CheckAccessControlsTest) TestWarningInClassStaticBlock1() {
	c.Test(
		c.Srcs(
			`/** @deprecated Woah! */
var x;

class Foo {
  static { x; }
}
`),
		c.DeprecatedName("Variable x has been deprecated: Woah!"))
}

// Java: CheckAccessControlsTest.testWarningInClassStaticBlock2()
func (c *CheckAccessControlsTest) TestWarningInClassStaticBlock2() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @deprecated Bad :D */
   static x = 1;
   static {
      var y = this.x;
   }
}
`),
		c.DeprecatedProp("Property x of type this has been deprecated: Bad :D"))
}

// Java: CheckAccessControlsTest.testNoWarningInClassStaticBlock()
func (c *CheckAccessControlsTest) TestNoWarningInClassStaticBlock() {
	// TODO(b/235871861): Compiler should throw warning as it doesn't make sense to deprecate a
	// static block
	c.Test(
		c.Srcs(
			`class Foo {
  /** @deprecated D: */
  static {
  }
}
`))
}

// Java: CheckAccessControlsTest.testNoWarningInClassStaticBlock1()
func (c *CheckAccessControlsTest) TestNoWarningInClassStaticBlock1() {
	// TODO(b/239747805): Should say "Property x of type number" not "super"
	c.Test(
		c.Srcs(
			`class A {
  /** @deprecated Meow! */
  static x = 1;
}
class B extends A {
  static {
    this.y = super.x;
  }
}
`),
		c.DeprecatedProp("Property x of type super has been deprecated: Meow!"))
}

// Java: CheckAccessControlsTest.testWarningForSubclassMethod()
func (c *CheckAccessControlsTest) TestWarningForSubclassMethod() {
	c.Test(
		c.Srcs(
			`class Foo {
  bar() {}
}

class SubFoo extends Foo {
  /**
   * @override
   * @deprecated I have a parent class!
   */
  bar() {}
}

function f() { (new SubFoo()).bar(); }
`),
		c.DeprecatedProp("Property bar of type SubFoo has been deprecated: I have a parent class!"))
}

// Java: CheckAccessControlsTest.testWarningForSuperClassWithDeprecatedSubclassMethod()
func (c *CheckAccessControlsTest) TestWarningForSuperClassWithDeprecatedSubclassMethod() {
	c.Test(
		c.Srcs(
			`class Foo {
  bar() {}
}

class SubFoo extends Foo {
  /**
   * @override
   * @deprecated I have a parent class!
   */
  bar() {}
}

function f() { (new Foo()).bar(); };
`))
}

// Java: CheckAccessControlsTest.testWarningForSuperclassMethod()
func (c *CheckAccessControlsTest) TestWarningForSuperclassMethod() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @deprecated I have a child class! */
  bar() {}
}

class SubFoo extends Foo {
  /** @override */
  bar() {}
}

function f() { (new SubFoo()).bar(); };
`),
		c.DeprecatedProp("Property bar of type SubFoo has been deprecated: I have a child class!"))
}

// Java: CheckAccessControlsTest.testWarningForSuperclassMethod2()
func (c *CheckAccessControlsTest) TestWarningForSuperclassMethod2() {
	c.Test(
		c.Srcs(
			`class Foo {
  /**
   * @protected
   * @deprecated I have another child class...
   */
  bar() {}
}

class SubFoo extends Foo {
  /**
   * @protected
   * @override
   */
  bar() {}
}

function f() { (new SubFoo()).bar(); };
`),
		c.DeprecatedProp(
			"Property bar of type SubFoo has been deprecated: I have another child class..."))
}

// Java: CheckAccessControlsTest.testWarningForDeprecatedClassInGlobalScope()
func (c *CheckAccessControlsTest) TestWarningForDeprecatedClassInGlobalScope() {
	c.Test(
		c.Srcs(
			`/** @deprecated I'm a very worldly object! */
var Foo = class {};

new Foo();
`),
		c.DeprecatedClass("Class Foo has been deprecated: I'm a very worldly object!"))
}

// Java: CheckAccessControlsTest.testNoWarningForPrototypeCopying()
func (c *CheckAccessControlsTest) TestNoWarningForPrototypeCopying() {
	c.Test(
		c.Srcs(
			`var Foo = class {
  bar() {}
};

/** @deprecated */
Foo.prototype.baz = Foo.prototype.bar;

(new Foo()).bar();
`))
}

// Java: CheckAccessControlsTest.testNoWarningOnDeprecatedPrototype()
func (c *CheckAccessControlsTest) TestNoWarningOnDeprecatedPrototype() {
	c.Test(
		c.Srcs(
			`var Foo = class {};

/** @deprecated */
Foo.prototype;
`))
}

// Java: CheckAccessControlsTest.testPrivateAccessForProperties1()
func (c *CheckAccessControlsTest) TestPrivateAccessForProperties1() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @private */
  bar_() {}

  baz() { this.bar_(); }
}

(new Foo).bar_();
`))
}

// Java: CheckAccessControlsTest.testPrivateAccessForProperties2()
func (c *CheckAccessControlsTest) TestPrivateAccessForProperties2() {
	c.Test(
		c.Srcs(
			"class Foo {}",
			`Foo.prototype.bar_ = function() {};
Foo.prototype.baz = function() { this.bar_(); };

(new Foo).bar_();
`))
}

// Java: CheckAccessControlsTest.testPrivateAccessForProperties3()
func (c *CheckAccessControlsTest) TestPrivateAccessForProperties3() {
	// Even though baz is "part of the Foo class" the access is disallowed since it's
	// not in the same file.
	c.Test(
		c.Srcs(
			`class Foo {
  /** @private */
  bar_() {}
}
`,
			"Foo.prototype.baz = function() { this.bar_(); };"),
		c.Error(CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testPrivateAccessForProperties4()
func (c *CheckAccessControlsTest) TestPrivateAccessForProperties4() {
	c.Test(
		c.Srcs(
			`/** @unrestricted */
class Foo {
  /** @private */
  bar_() {}

  ['baz']() { (new Foo()).bar_(); }
}
`))
}

// Java: CheckAccessControlsTest.testPrivateAccessForProperties5()
func (c *CheckAccessControlsTest) TestPrivateAccessForProperties5() {
	c.Test(
		c.Srcs(
			`class Parent {
  constructor() {
    /** @private */
    this.prop = 'foo';
  }
};
`,
			`class Child extends Parent {
  constructor() {
    this.prop = 'asdf';
  }
}
`),
		c.Error(CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS).
			WithMessage("Access to private property prop of Parent not allowed here."))
}

// Java: CheckAccessControlsTest.testPrivatePropAccess_inSameFile_throughDestructuring()
func (c *CheckAccessControlsTest) TestPrivatePropAccess_inSameFile_throughDestructuring() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @private */
  bar_() { }
}

function f(/** !Foo */ x) {
  const {bar_: bar} = x;
}
`))
}

// Java: CheckAccessControlsTest.testPrivateAccessForProperties6()
func (c *CheckAccessControlsTest) TestPrivateAccessForProperties6() {
	c.Test(
		c.Srcs(
			`goog.provide('x.y.z.Parent');

x.y.z.Parent = class {
  constructor() {
    /** @private */
    this.prop = 'foo';
  }
};
`,
			`goog.require('x.y.z.Parent');

class Child extends x.y.z.Parent {
  constructor() {
    this.prop = 'asdf';
  }
}
`),
		c.Error(CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS).
			WithMessage("Access to private property prop of x.y.z.Parent not allowed here."))
}

// Java: CheckAccessControlsTest.testPrivateAccessToConstructorThroughNew()
func (c *CheckAccessControlsTest) TestPrivateAccessToConstructorThroughNew() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @private */
  constructor() { }
}

new Foo();
`))
}

// Java: CheckAccessControlsTest.testPrivateAccessToConstructorThroughExtendsClause()
func (c *CheckAccessControlsTest) TestPrivateAccessToConstructorThroughExtendsClause() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @private */
  constructor() { }
}

class SubFoo extends Foo { }
`))
}

// Java: CheckAccessControlsTest.testPrivateAccessToClass()
func (c *CheckAccessControlsTest) TestPrivateAccessToClass() {
	c.Test(
		c.Srcs(
			`/** @private */
class Foo { }

Foo;
`))
}

// Java: CheckAccessControlsTest.testPrivateAccess_googModule()
func (c *CheckAccessControlsTest) TestPrivateAccess_googModule() {
	c.Test(
		c.Srcs(
			`goog.module('example.One');

class One {
  /** @private */
  m() {}
}

exports = One;
`,
			`goog.module('example.two');

const One = goog.require('example.One');

(new One()).m();
`),
		c.Error(CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS).
			WithMessage("Access to private property m of One not allowed here."))
}

// Java: CheckAccessControlsTest.testNoPrivateAccessForProperties1()
func (c *CheckAccessControlsTest) TestNoPrivateAccessForProperties1() {
	c.Test(
		c.Srcs(
			`class Foo {};
(new Foo).bar_();
`,
			`/** @private */
Foo.prototype.bar_ = function() {};

Foo.prototype.baz = function() { this.bar_(); };
`),
		c.Error(CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoPrivateAccessForProperties2()
func (c *CheckAccessControlsTest) TestNoPrivateAccessForProperties2() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @private */
  bar_() {}

  baz() { this.bar_(); }
}
`,
			"(new Foo).bar_();"),
		c.Error(CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoPrivateAccessForProperties3()
func (c *CheckAccessControlsTest) TestNoPrivateAccessForProperties3() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @private */
  bar_() {}
}
`,
			`class OtherFoo {
  constructor() {
    (new Foo).bar_();
  }
}
`),
		c.Error(CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoPrivateAccessForProperties4()
func (c *CheckAccessControlsTest) TestNoPrivateAccessForProperties4() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @private */
  bar_() {}
}
`,
			`class SubFoo extends Foo {
  constructor() {
    this.bar_();
  }
}
`),
		c.Error(CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoPrivateAccessForProperties5()
func (c *CheckAccessControlsTest) TestNoPrivateAccessForProperties5() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @private */
  bar_() {}
}
`,
			`class SubFoo extends Foo {
  baz() { this.bar_(); }
}
`),
		c.Error(CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoPrivateAccessForProperties6()
func (c *CheckAccessControlsTest) TestNoPrivateAccessForProperties6() {
	// Overriding a private property with a non-private property
	// in a different file causes problems.
	c.Test(
		c.Srcs(
			`class Foo {
  /** @private */
  bar_() {}
}
`,
			`class SubFoo extends Foo {
  bar_() {}
}
`),
		c.Error(CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoPrivateAccessForProperties6a()
func (c *CheckAccessControlsTest) TestNoPrivateAccessForProperties6a() {
	// Same as above, except with namespaced constructors
	c.Test(
		c.Srcs(
			`/** @const */ var ns = {};

ns.Foo = class {
  /** @private */
  bar_() {}
};
`,
			`ns.SubFoo = class extends ns.Foo {
  bar_() {}
};
`),
		c.Error(CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoPrivateAccessForProperties7()
func (c *CheckAccessControlsTest) TestNoPrivateAccessForProperties7() {
	// It's OK to override a private property with a non-private property
	// in the same file, but you'll get yelled at when you try to use it.
	c.Test(
		c.Srcs(
			`class Foo {
  /** @private */
  bar_() {}
}

class SubFoo extends Foo {
  bar_() {}
}
`,
			"SubFoo.prototype.baz = function() { this.bar_(); }"),
		c.Error(CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoPrivateAccessForProperties8()
func (c *CheckAccessControlsTest) TestNoPrivateAccessForProperties8() {
	c.Test(
		c.Srcs(
			`class Foo {
  constructor() {
    /** @private */
    this.bar_ = 3;
  }
}
`,
			`class SubFoo extends Foo {
  constructor() {
    /** @private */
    this.bar_ = 3;
  }
}
`),
		c.Error(CheckAccessControls_PRIVATE_OVERRIDE))
}

// Java: CheckAccessControlsTest.testNoPrivateAccessForProperties9()
func (c *CheckAccessControlsTest) TestNoPrivateAccessForProperties9() {
	c.Test(
		c.Srcs(
			`class Foo {}

/** @private */
Foo.prototype.bar_ = 3;
`,
			"new Foo().bar_;"),
		c.Error(CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoPrivateAccessForProperties10()
func (c *CheckAccessControlsTest) TestNoPrivateAccessForProperties10() {
	c.Test(
		c.Srcs(
			`class Foo {}

/** @private */
Foo.prototype.bar_ = function() {};
`,
			"new Foo().bar_;"),
		c.Error(CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoPrivateAccessForProperties11()
func (c *CheckAccessControlsTest) TestNoPrivateAccessForProperties11() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @private */
  get bar_() { return 1; }
}
`,
			"var a = new Foo().bar_;"),
		c.Error(CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoPrivateAccessForProperties12()
func (c *CheckAccessControlsTest) TestNoPrivateAccessForProperties12() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @private */
  set bar_(x) { this.barValue = x; }
}
`,
			"new Foo().bar_ = 1;"),
		c.Error(CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoPrivateAccessForProperties13()
func (c *CheckAccessControlsTest) TestNoPrivateAccessForProperties13() {
	// Visibility of getter and setter differ.
	fooFileSrc := `class Foo {
  constructor() {
    /** @private */
    this.barValue_ = 1;
  }

  /**
   * @param {number} x
   * @private
   */
  set bar(x) {
    this.barValue_ = x;
  }

  /** @return {number} */
  get bar() {
    return this.barValue_;
  }
}
`
	// TODO: b/416023147 - This should not be an error.
	c.Test(c.Srcs(fooFileSrc, "const barCopy = new Foo().bar;"), c.Error(CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS))
	// This is fine.
	c.Test(c.Srcs(fooFileSrc, "new Foo().bar = 2;"), c.Error(CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoPrivatePropAccess_inDifferentFile_throughDestructuring()
func (c *CheckAccessControlsTest) TestNoPrivatePropAccess_inDifferentFile_throughDestructuring() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @private */
  bar_() { }
}
`,
			`function f(/** !Foo */ x) {
  const {bar_: bar} = x;
}
`),
		c.Error(CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoPrivateAccessToConstructorThroughNew()
func (c *CheckAccessControlsTest) TestNoPrivateAccessToConstructorThroughNew() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @private */
  constructor() { }
}
`,
			"new Foo();"),
		c.Error(CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoPrivateAccessToConstructorThroughExtendsClause()
func (c *CheckAccessControlsTest) TestNoPrivateAccessToConstructorThroughExtendsClause() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @private */
  constructor() { }
}
`,
			"class SubFoo extends Foo { }"),
		c.Error(CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoPrivateAccessToClass()
func (c *CheckAccessControlsTest) TestNoPrivateAccessToClass() {
	c.Test(
		c.Srcs(
			`/** @private */
class Foo { }
`,
			"Foo"),
		c.Error(CheckAccessControls_BAD_PRIVATE_GLOBAL_ACCESS))
}

// Java: CheckAccessControlsTest.testNoPrivateAccessForFields()
func (c *CheckAccessControlsTest) TestNoPrivateAccessForFields() {
	// Overriding a private property with a non-private property
	// in a different file causes problems.
	c.Test(
		c.Srcs(
			`class Foo {
  /** @private */
  bar
}
`,
			`class SubFoo extends Foo {
  bar
}
`),
		c.Error(CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoPrivateAccessForFields2()
func (c *CheckAccessControlsTest) TestNoPrivateAccessForFields2() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @private */
  bar
}
`,
			"new Foo().bar"),
		c.Error(CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoPrivateAccessForFields3()
func (c *CheckAccessControlsTest) TestNoPrivateAccessForFields3() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @private */
  bar = 2
}
`,
			"new Foo().bar = 4"),
		c.Error(CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testPrivateAccessForStaticBlock()
func (c *CheckAccessControlsTest) TestPrivateAccessForStaticBlock() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @private */
  static bar = 2
  static {
    this.bar = 4;
  }
}
`))
}

// Java: CheckAccessControlsTest.testNoPrivateAccessForStaticBlock()
func (c *CheckAccessControlsTest) TestNoPrivateAccessForStaticBlock() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @private */
  static bar = 2;
}
`,
			`class Bar extends Foo{
  static {
    this.bar = 4;
  }
}
`),
		c.Error(CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testProtectedAccessForProperties1()
func (c *CheckAccessControlsTest) TestProtectedAccessForProperties1() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @protected */
  bar() {}
}

(new Foo).bar();
`,
			"Foo.prototype.baz = function() { this.bar(); };"))
}

// Java: CheckAccessControlsTest.testProtectedAccessForProperties2()
func (c *CheckAccessControlsTest) TestProtectedAccessForProperties2() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @protected */
  bar() {}
}

(new Foo).bar();
`,
			`class SubFoo extends Foo {
  constructor() {
    this.bar();
  }
}
`))
}

// Java: CheckAccessControlsTest.testProtectedAccessForProperties3()
func (c *CheckAccessControlsTest) TestProtectedAccessForProperties3() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @protected */
  bar() {}
}

(new Foo).bar();
`,
			`class SubFoo extends Foo {
  static baz() { (new Foo()).bar(); }
}
`))
}

// Java: CheckAccessControlsTest.testProtectedAccessForProperties4()
func (c *CheckAccessControlsTest) TestProtectedAccessForProperties4() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @protected */
  static bar() {}
}
`,
			`class SubFoo extends Foo {
  static foo() { super.bar(); }
}
`))
}

// Java: CheckAccessControlsTest.testProtectedAccessForProperties5()
func (c *CheckAccessControlsTest) TestProtectedAccessForProperties5() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @protected */ bar() {}
}

(new Foo).bar();
`,
			`var SubFoo = class extends Foo {
  constructor() {
    super();
    this.bar();
  }
}
`))
}

// Java: CheckAccessControlsTest.testProtectedAccessForProperties6()
func (c *CheckAccessControlsTest) TestProtectedAccessForProperties6() {
	c.Test(
		c.Srcs(
			`goog.Foo = class {
  /** @protected */
  bar() {}
}
`,
			`goog.SubFoo = class extends goog.Foo {
  constructor() {
    super();
    this.bar();
  }
};
`))
}

// Java: CheckAccessControlsTest.testProtectedAccessForProperties7()
func (c *CheckAccessControlsTest) TestProtectedAccessForProperties7() {
	c.Test(
		c.Srcs(
			`var Foo = class {};

/** @protected */
Foo.prototype.bar = function() {};
`,
			`var SubFoo = class extends Foo {
  constructor() {
    super();
    this.bar();
  }
};

SubFoo.prototype.moo = function() { this.bar(); };
`))
}

// Java: CheckAccessControlsTest.testProtectedAccessForProperties8()
func (c *CheckAccessControlsTest) TestProtectedAccessForProperties8() {
	c.Test(
		c.Srcs(
			`var Foo = class {};

/** @protected */
Foo.prototype.bar = function() {};
`,
			`var SubFoo = class extends Foo {
  get moo() { this.bar(); }
};
`))
}

// Java: CheckAccessControlsTest.testProtectedAccessForProperties9()
func (c *CheckAccessControlsTest) TestProtectedAccessForProperties9() {
	c.Test(
		c.Srcs(
			`/** @unrestricted */
var Foo = class {};

/** @protected */
Foo.prototype.bar = function() {};
`,
			`/** @unrestricted */
var SubFoo = class extends Foo {
  set moo(val) { this.x = this.bar(); }
};
`))
}

// Java: CheckAccessControlsTest.testProtectedAccessForProperties10()
func (c *CheckAccessControlsTest) TestProtectedAccessForProperties10() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				"foo.js",
				`var Foo = class {
  /** @protected */
  bar() {}
}
`),
			SourceFile_FromCode(
				"sub_foo.js",
				`var SubFoo = class extends Foo {};

(function() {
  SubFoo.prototype.baz = function() { this.bar(); }
})();
`)))
}

// Java: CheckAccessControlsTest.testProtectedAccessForProperties11()
func (c *CheckAccessControlsTest) TestProtectedAccessForProperties11() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				"foo.js",
				`goog.provide('Foo');

/** @interface */
Foo = class {};

/** @protected */
Foo.prop = {};
`),
			SourceFile_FromCode(
				"bar.js",
				`goog.require('Foo');

/** @implements {Foo} */
class Bar {
  constructor() {
    Foo.prop;
  }
};
`)))
}

// Java: CheckAccessControlsTest.testProtectedAccessForProperties12()
func (c *CheckAccessControlsTest) TestProtectedAccessForProperties12() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				"a.js",
				`goog.provide('A');

var A = class {
  constructor() {
    /** @protected {string} */
    this.prop;
  }
};
`),
			SourceFile_FromCode(
				"b.js",
				`goog.require('A');

var B = class extends A {
  constructor() {
    super();

    this.prop.length;
  }
};
`)))
}

// FYI: Java warns for the b1.method access in c.js.
// Instead of following that in NTI, we chose to follow the behavior of
// the old JSCompiler type checker, to make migration easier.
// Java: CheckAccessControlsTest.testProtectedAccessForProperties13()
func (c *CheckAccessControlsTest) TestProtectedAccessForProperties13() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				"a.js",
				`goog.provide('A');

var A = class {
  /** @protected */
  method() {}
}
`),
			SourceFile_FromCode(
				"b1.js",
				`goog.require('A');

goog.provide('B1');

var B1 = class extends A {
  /** @override */
  method() {}
}
`),
			SourceFile_FromCode(
				"b2.js",
				`goog.require('A');

goog.provide('B2');

var B2 = class extends A {
  /** @override */
  method() {}
};
`),
			SourceFile_FromCode(
				"c.js",
				`goog.require('B1');
goog.require('B2');

var C = class extends B2 {
  /** @param {!B1} b1 */
  constructor(b1) {
    var x = b1.method();
  }
};
`)))
}

// Java: CheckAccessControlsTest.testProtectedAccessForProperties14()
func (c *CheckAccessControlsTest) TestProtectedAccessForProperties14() {
	// access in member function
	c.Test(
		c.Srcs(
			`var Foo = class {};

/** @protected */
Foo.prototype.bar = function() {};
`,
			`var OtherFoo = class extends Foo {
  constructor() {
    super();

    this.bar();
  }
};

OtherFoo.prototype.moo = function() { new Foo().bar(); };
`))
}

// Java: CheckAccessControlsTest.testProtectedAccessForProperties15()
func (c *CheckAccessControlsTest) TestProtectedAccessForProperties15() {
	// access in computed member function
	c.Test(
		c.Srcs(
			`/** @unrestricted */
var Foo = class {};

/** @protected */
Foo.prototype.bar = function() {};
`,
			`/** @unrestricted */
var OtherFoo = class extends Foo {
  constructor() {
    super();
    this['bar']();
  };
}

OtherFoo.prototype['bar'] = function() { new Foo().bar(); };
`))
}

// Java: CheckAccessControlsTest.testProtectedAccessForProperties16()
func (c *CheckAccessControlsTest) TestProtectedAccessForProperties16() {
	// access in nested arrow function
	c.Test(
		c.Srcs(
			`class Foo {
  /** @protected */ bar() {}
}
`,
			`class OtherFoo extends Foo {
  constructor() { super(); var f = () => this.bar(); }
  baz() { return () => this.bar(); }
}
`))
}

// Java: CheckAccessControlsTest.testProtectedAccessForProperties17()
func (c *CheckAccessControlsTest) TestProtectedAccessForProperties17() {
	// Access in subclass field initializer where Subclass parent is SCRIPT.
	c.Test(
		c.Srcs(
			`class SuperClass {
  /** @protected @return {number} */
  getNum() { return 1; }
}
`,
			`class Subclass extends SuperClass {
  numField = this.getNum();
}
`))
}

// Java: CheckAccessControlsTest.testProtectedAccessForProperties18()
func (c *CheckAccessControlsTest) TestProtectedAccessForProperties18() {
	// Access in subclass field initializer where Subclass parent is BLOCK.
	c.Test(
		c.Srcs(
			`class SuperClass {
  /** @protected @return {number} */
  getNum() { return 1; }
}
`,
			`{
  class Subclass extends SuperClass {
    numField = this.getNum();
  }
}
`))
}

// Java: CheckAccessControlsTest.testProtectedAccessForProperties19()
func (c *CheckAccessControlsTest) TestProtectedAccessForProperties19() {
	// Access in subclass field initializer where class parent is ASSIGN.
	c.Test(
		c.Srcs(
			`class SuperClass {
  /** @protected @return {number} */
  getNum() { return 1; }
}
`,
			`const Subclass = class extends SuperClass {
  numField = this.getNum();
}
`))
}

// Java: CheckAccessControlsTest.testProtectedAccessForProperties20()
func (c *CheckAccessControlsTest) TestProtectedAccessForProperties20() {
	// Access in subclass field initializer where class parent is FUNCTION.
	c.Test(
		c.Srcs(
			`class SuperClass {
  /** @protected @return {number} */
  getNum() { return 1; }
}
`,
			`const Subclass = (() => class extends SuperClass {
  numField = this.getNum();
})();
`))
}

// Java: CheckAccessControlsTest.testProtectedAccessForProperties21()
func (c *CheckAccessControlsTest) TestProtectedAccessForProperties21() {
	// Access in subclass field initializer where class parent is RETURN.
	c.Test(
		c.Srcs(
			`class SuperClass {
  /** @protected @return {number} */
  getNum() { return 1; }
}
`,
			`function classGenerator() {
  return class extends SuperClass {
    numField = this.getNum();
  };
}
const Subclass = classGenerator();
`))
}

// Java: CheckAccessControlsTest.testProtectedAccessForProperties22()
func (c *CheckAccessControlsTest) TestProtectedAccessForProperties22() {
	// Access in subclass field initializer where class parent is GETPROP.
	c.Test(
		c.Srcs(
			`class SuperClass {
  /** @protected @return {number} */
  getNum() { return 1; }
}
`,
			`function getValue() {
  return class extends SuperClass {
    numField = this.getNum();
    static staticField = 1;
  }.staticField;
}
const value = getValue();
`))
}

// Java: CheckAccessControlsTest.testProtectedPropAccess_inDifferentFile_inSubclass_throughDestructuring()
func (c *CheckAccessControlsTest) TestProtectedPropAccess_inDifferentFile_inSubclass_throughDestructuring() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @protected */
  bar() { }
}
`,
			`class SubFoo extends Foo {
  method(/** !Foo */ x) {
    const {bar: bar} = x;
  }
}
`))
}

// Java: CheckAccessControlsTest.testNoProtectedAccess_forOverriddenProperty_elsewhereInSubclassFile()
func (c *CheckAccessControlsTest) TestNoProtectedAccess_forOverriddenProperty_elsewhereInSubclassFile() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @protected */
  bar() { }
}
`,
			`class Bar extends Foo {
  /** @override */
  bar() { }
}

// TODO(b/113705099): This should be legal.
(new Bar()).bar();
`),
		c.Error(CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testProtectedAccessToConstructorThroughExtendsClause()
func (c *CheckAccessControlsTest) TestProtectedAccessToConstructorThroughExtendsClause() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @protected */ constructor() {}
}
`,
			"class OtherFoo extends Foo { }"))
}

// Java: CheckAccessControlsTest.testProtectedAccessToConstructorThroughSubclassInstanceMethod()
func (c *CheckAccessControlsTest) TestProtectedAccessToConstructorThroughSubclassInstanceMethod() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @protected */ constructor() {}
}
`,
			`class OtherFoo extends Foo {
  bar() { new Foo(); }
}
`))
}

// Java: CheckAccessControlsTest.testProtectedAccessToClassThroughSubclassInstanceMethod()
func (c *CheckAccessControlsTest) TestProtectedAccessToClassThroughSubclassInstanceMethod() {
	c.Test(
		c.Srcs(
			`/** @protected */
class Foo {}
`,
			`class OtherFoo extends Foo {
  bar() { Foo; }
}
`))
}

// Java: CheckAccessControlsTest.testProtectedAccessThroughNestedFunction()
func (c *CheckAccessControlsTest) TestProtectedAccessThroughNestedFunction() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @protected */ bar() {}
}
`,
			`class Bar extends Foo {
  constructor(/** Foo */ foo) {
    function f() {
      foo.bar();
    }
  }
}
`))
}

// Java: CheckAccessControlsTest.testProtectedAccessThroughNestedEs5Class()
func (c *CheckAccessControlsTest) TestProtectedAccessThroughNestedEs5Class() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @protected */ bar() {}
}
`,
			`class Bar extends Foo {
  constructor() {
    /** @constructor */
    const Nested = function() {}

    /** @param {!Foo} foo */
    Nested.prototype.qux = function(foo) {
      foo.bar();
    }
  }
}
`))
}

// Java: CheckAccessControlsTest.testProtectedAccessThroughNestedEs6Class()
func (c *CheckAccessControlsTest) TestProtectedAccessThroughNestedEs6Class() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @protected */ bar() {}
}
`,
			`class Bar extends Foo {
  constructor() {
    class Nested {
      qux(/** Foo */ foo) {
        foo.bar();
      }
    }
  }
}
`))
}

// Java: CheckAccessControlsTest.testNoProtectedAccessForProperties1()
func (c *CheckAccessControlsTest) TestNoProtectedAccessForProperties1() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @protected */
  bar() {}
}
`,
			"(new Foo).bar();"),
		c.Error(CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoProtectedAccessForProperties2()
func (c *CheckAccessControlsTest) TestNoProtectedAccessForProperties2() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @protected */
  bar() {}
}
`,
			`class OtherFoo {
  constructor() {
    (new Foo).bar();
  }
}
`),
		c.Error(CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoProtectedAccessForProperties3()
func (c *CheckAccessControlsTest) TestNoProtectedAccessForProperties3() {
	c.Test(
		c.Srcs(
			`class Foo {}

class SubFoo extends Foo {
  /** @protected */
  bar() {}
}
`,
			`class SubberFoo extends Foo {
  constructor() {
    (new SubFoo).bar();
  }
}
`),
		c.Error(CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoProtectedAccessForProperties4()
func (c *CheckAccessControlsTest) TestNoProtectedAccessForProperties4() {
	c.Test(
		c.Srcs(
			`class Foo {
  constructor() {
    (new SubFoo).bar();
  }
}
`,
			`class SubFoo extends Foo {
  /** @protected */
  bar() {}
}
`),
		c.Error(CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoProtectedAccessForProperties5()
func (c *CheckAccessControlsTest) TestNoProtectedAccessForProperties5() {
	c.Test(
		c.Srcs(
			`goog.Foo = class {
  /** @protected */
  bar() {}
}
`,
			`goog.NotASubFoo = class {
  constructor() {
    (new goog.Foo).bar();
  }
}
`),
		c.Error(CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoProtectedAccessForProperties6()
func (c *CheckAccessControlsTest) TestNoProtectedAccessForProperties6() {
	c.Test(
		c.Srcs(
			`class Foo {}

/** @protected */
Foo.prototype.bar = 3;
`,
			"new Foo().bar;"),
		c.Error(CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoProtectedAccessForProperties7()
func (c *CheckAccessControlsTest) TestNoProtectedAccessForProperties7() {
	c.Test(
		c.Srcs(
			`class Foo {}

/** @protected */
Foo.prototype.bar = function() {};
`,
			"new Foo().bar();"),
		c.Error(CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoProtectedAccessForProperties8()
func (c *CheckAccessControlsTest) TestNoProtectedAccessForProperties8() {
	c.Test(
		c.Srcs(
			`var Foo = class {};

/** @protected */
Foo.prototype.bar = function() {};
`,
			`var OtherFoo = class {
  constructor() {
    this.bar();
  }
}

OtherFoo.prototype.moo = function() { new Foo().bar(); };
`),
		c.Error(CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoProtectedAccessForProperties9()
func (c *CheckAccessControlsTest) TestNoProtectedAccessForProperties9() {
	c.Test(
		c.Srcs(
			`var Foo = class {};

/** @protected */
Foo.prototype.bar = function() {};
`,
			`var OtherFoo = class {
  constructor() { this['bar'](); }
};

OtherFoo.prototype['bar'] = function() { new Foo().bar(); };
`),
		c.Error(CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoProtectedPropAccess_inDifferentFile_throughDestructuring()
func (c *CheckAccessControlsTest) TestNoProtectedPropAccess_inDifferentFile_throughDestructuring() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @protected */
  bar() { }
}
`,
			`function f(/** !Foo */ x) {
  const {bar: bar} = x;
}
`),
		c.Error(CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoProtectedAccess_forInheritedProperty_elsewhereInSubclassFile()
func (c *CheckAccessControlsTest) TestNoProtectedAccess_forInheritedProperty_elsewhereInSubclassFile() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @protected */
  bar() { }
}
`,
			`class Bar extends Foo { }

(new Bar()).bar();
`),
		c.Error(CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoProtectedAccessToConstructorFromUnrelatedClass()
func (c *CheckAccessControlsTest) TestNoProtectedAccessToConstructorFromUnrelatedClass() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @protected */
  constructor() {}
}
`,
			`class OtherFoo {
  bar() { new Foo(); }
}
`),
		c.Error(CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoProtectedAccessForPropertiesWithNoRhs()
func (c *CheckAccessControlsTest) TestNoProtectedAccessForPropertiesWithNoRhs() {
	c.Test(
		c.Srcs(
			`class Foo {}

/** @protected */
Foo.prototype.x;
`,
			`class Bar extends Foo {}

/** @protected */
Bar.prototype.x;
`))
}

// Java: CheckAccessControlsTest.testNoProtectedAccessForFields1()
func (c *CheckAccessControlsTest) TestNoProtectedAccessForFields1() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @protected */
  bar
}
`,
			"(new Foo).bar;"),
		c.Error(CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoProtectedAccessForFields2()
func (c *CheckAccessControlsTest) TestNoProtectedAccessForFields2() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @protected */
  bar
}
`,
			`class OtherFoo {
  constructor() {
    (new Foo).bar;
  }
}
`),
		c.Error(CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoProtectedAccessForFields3()
func (c *CheckAccessControlsTest) TestNoProtectedAccessForFields3() {
	c.Test(
		c.Srcs(
			`class Foo {}

class SubFoo extends Foo {
  /** @protected */
  bar
}
`,
			`class SubberFoo extends Foo {
  constructor() {
    (new SubFoo).bar;
  }
}
`),
		c.Error(CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoProtectedAccessForFields4()
func (c *CheckAccessControlsTest) TestNoProtectedAccessForFields4() {
	c.Test(
		c.Srcs(
			`class Foo {
  constructor() {
    (new SubFoo).bar;
  }
}
`,
			`class SubFoo extends Foo {
  /** @protected */
  bar
}
`),
		c.Error(CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoProtectedAccessForFields5()
func (c *CheckAccessControlsTest) TestNoProtectedAccessForFields5() {
	c.Test(
		c.Srcs(
			`goog.Foo = class {
  /** @protected */
  bar
}
`,
			`goog.NotASubFoo = class {
  constructor() {
    (new goog.Foo).bar;
  }
}
`),
		c.Error(CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testProtectedAccessForStaticBlocks()
func (c *CheckAccessControlsTest) TestProtectedAccessForStaticBlocks() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @protected */
  static bar
  static {
    this.bar;
  }
}
`))
}

// Java: CheckAccessControlsTest.testProtectedAccessForStaticBlocks_sameFile()
func (c *CheckAccessControlsTest) TestProtectedAccessForStaticBlocks_sameFile() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @protected */
  static bar
}
class SubFoo extends Foo {
  static {
    this.bar;
  }
}
`))
}

// Java: CheckAccessControlsTest.testProtectedAccessForStaticBlocks_sameFile1()
func (c *CheckAccessControlsTest) TestProtectedAccessForStaticBlocks_sameFile1() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @protected */
  static bar
}
class Bar {
  static {
    Foo.bar;
  }
}
`))
}

// Java: CheckAccessControlsTest.testProtectedAccessForStatic_differentFile()
func (c *CheckAccessControlsTest) TestProtectedAccessForStatic_differentFile() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @protected */
  static bar
}
`,
			`class SubFoo extends Foo {
  static {
    this.bar;
  }
}
`))
}

// Java: CheckAccessControlsTest.testNoProtectedAccessForStaticBlocks_differentFile()
func (c *CheckAccessControlsTest) TestNoProtectedAccessForStaticBlocks_differentFile() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @protected */
  static bar
}
`,
			`class Bar {
  static {
    Foo.bar;
  }
}
`),
		c.Error(CheckAccessControls_BAD_PROTECTED_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testPackagePrivateAccessForProperties1()
func (c *CheckAccessControlsTest) TestPackagePrivateAccessForProperties1() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @package */
  bar() {}

  baz() { this.bar(); }
}

(new Foo).bar();
`))
}

// Java: CheckAccessControlsTest.testPackagePrivateAccessForProperties2()
func (c *CheckAccessControlsTest) TestPackagePrivateAccessForProperties2() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(Compiler_JoinPathParts("foo", "bar.js"), "class Foo {}"),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"),
				`/** @package */
Foo.prototype.bar = function() {};

Foo.prototype.baz = function() { this.bar(); };

(new Foo).bar();
`)))
}

// Java: CheckAccessControlsTest.testPackagePrivateAccessForProperties3()
func (c *CheckAccessControlsTest) TestPackagePrivateAccessForProperties3() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`class Foo {
  /** @package */
  bar() {}
}

(new Foo).bar();
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "baz.js"),
				"Foo.prototype.baz = function() { this.bar(); };")))
}

// Java: CheckAccessControlsTest.testPackagePrivateAccessForProperties4()
func (c *CheckAccessControlsTest) TestPackagePrivateAccessForProperties4() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`class Foo {
  /** @package */
  bar() {}
}
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "baz.js"),
				"Foo.prototype['baz'] = function() { (new Foo()).bar(); };")))
}

// Java: CheckAccessControlsTest.testPackagePrivateAccessForProperties5()
func (c *CheckAccessControlsTest) TestPackagePrivateAccessForProperties5() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`class Parent {
  constructor() {
    /** @package */
    this.prop = 'foo';
  }
}
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"),
				`class Child extends Parent {
  constructor() {
    this.prop = 'asdf';
  }
}
`)),
		c.Error(CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testPackagePropAccess_inSamePackage_throughDestructuring()
func (c *CheckAccessControlsTest) TestPackagePropAccess_inSamePackage_throughDestructuring() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`class Foo {
  /** @package */
  bar() { }
}
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "quux.js"),
				`function f(/** !Foo */ x) {
  const {bar: bar} = x;
}
`)))
}

// Java: CheckAccessControlsTest.testNoPackagePropAccess_inDifferentPackage_throughDestructuring()
func (c *CheckAccessControlsTest) TestNoPackagePropAccess_inDifferentPackage_throughDestructuring() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`class Foo {
  /** @package */
  bar() { }
}
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"),
				`function f(/** !Foo */ x) {
  const {bar: bar} = x;
}
`)),
		c.Error(CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testPackageAccessToConstructorThroughNew()
func (c *CheckAccessControlsTest) TestPackageAccessToConstructorThroughNew() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`class Foo {
  /** @package */
  constructor() {}
}
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "quux.js"), //
				"new Foo();")))
}

// Java: CheckAccessControlsTest.testPackageAccessToConstructorThroughExtendsClause()
func (c *CheckAccessControlsTest) TestPackageAccessToConstructorThroughExtendsClause() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`class Foo {
  /** @package */
  constructor() {}
}
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "quux.js"), //
				"class SubFoo extends Foo { }")))
}

// Java: CheckAccessControlsTest.testPackageAccessToClass()
func (c *CheckAccessControlsTest) TestPackageAccessToClass() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/** @package */
class Foo {}
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "quux.js"), //
				"Foo;")))
}

// Java: CheckAccessControlsTest.testPackageAccessToConstructorThroughNew_fileOveriew()
func (c *CheckAccessControlsTest) TestPackageAccessToConstructorThroughNew_fileOveriew() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/** @fileoverview @package */

class Foo {
  constructor() {}
}
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "quux.js"), //
				"new Foo();")))
}

// Java: CheckAccessControlsTest.testPackageAccessToConstructorThroughExtendsClause_fileOveriew()
func (c *CheckAccessControlsTest) TestPackageAccessToConstructorThroughExtendsClause_fileOveriew() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/** @fileoverview @package */

class Foo {
  constructor() {}
}
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "quux.js"), //
				"class SubFoo extends Foo { }")))
}

// Java: CheckAccessControlsTest.testPackageAccessToClass_fileOveriew()
func (c *CheckAccessControlsTest) TestPackageAccessToClass_fileOveriew() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/** @fileoverview @package */

class Foo {}
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "quux.js"), //
				"Foo;")))
}

// Java: CheckAccessControlsTest.testNoPackagePrivateAccessForProperties1()
func (c *CheckAccessControlsTest) TestNoPackagePrivateAccessForProperties1() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`class Foo {}

(new Foo).bar();
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"),
				`/** @package */
Foo.prototype.bar = function() {};

Foo.prototype.baz = function() { this.bar(); };
`)),
		c.Error(CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoPackagePrivateAccessForProperties2()
func (c *CheckAccessControlsTest) TestNoPackagePrivateAccessForProperties2() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`class Foo {
  /** @package */
  bar() {}

  baz() { this.bar(); };
}
`),
			SourceFile_FromCode(Compiler_JoinPathParts("baz", "quux.js"), "(new Foo).bar();")),
		c.Error(CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoPackagePrivateAccessForProperties3()
func (c *CheckAccessControlsTest) TestNoPackagePrivateAccessForProperties3() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`class Foo {
  /** @package */
  bar() {};
}
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"),
				`class OtherFoo {
  constructor() {
    (new Foo).bar();
  }
}
`)),
		c.Error(CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoPackagePrivateAccessForProperties4()
func (c *CheckAccessControlsTest) TestNoPackagePrivateAccessForProperties4() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`class Foo {
  /** @package */
  bar() {}
}
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"),
				`class SubFoo extends Foo {
  constructor() {
    this.bar();
  }
}
`)),
		c.Error(CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoPackagePrivateAccessForProperties5()
func (c *CheckAccessControlsTest) TestNoPackagePrivateAccessForProperties5() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`class Foo {
  /** @package */
  bar() {}
}
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"),
				`class SubFoo extends Foo {
  baz() { this.bar(); }
}
`)),
		c.Error(CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoPackagePrivateAccessForProperties6()
func (c *CheckAccessControlsTest) TestNoPackagePrivateAccessForProperties6() {
	// Overriding a private property with a non-package-private property
	// in a different file causes problems.
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`class Foo {
  /** @package */
  bar() {}
}
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"),
				`class SubFoo extends Foo {
  bar() {}
}
`)),
		c.Error(CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoPackagePrivateAccessForProperties7()
func (c *CheckAccessControlsTest) TestNoPackagePrivateAccessForProperties7() {
	// It's OK to override a package-private property with a
	// non-package-private property in the same file, but you'll get
	// yelled at when you try to use it.
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`class Foo {
  /** @package */
  bar() {}
}

class SubFoo extends Foo {
  bar() {}
}
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"),
				"SubFoo.prototype.baz = function() { this.bar(); }")),
		c.Error(CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoPackgeAccessToConstructorThroughNew()
func (c *CheckAccessControlsTest) TestNoPackgeAccessToConstructorThroughNew() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`class Foo {
  /** @package */
  constructor() {}
}
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"), //
				"new Foo();")),
		c.Error(CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoPackageAccessToConstructorThroughExtendsClause()
func (c *CheckAccessControlsTest) TestNoPackageAccessToConstructorThroughExtendsClause() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`class Foo {
  /** @package */
  constructor() {}
}
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"), //
				"class SubFoo extends Foo { }")),
		c.Error(CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoPackageAccessToClass()
func (c *CheckAccessControlsTest) TestNoPackageAccessToClass() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/** @package */
class Foo {}
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"), //
				"Foo;")),
		c.Error(CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoPackgeAccessToConstructorThroughNew_fileOveriew()
func (c *CheckAccessControlsTest) TestNoPackgeAccessToConstructorThroughNew_fileOveriew() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/** @fileoverview @package */

/** @public */
class Foo {
  constructor() {}
}
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"), //
				"new Foo();")),
		c.Error(CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoPackageAccessToConstructorThroughExtendsClause_fileOveriew()
func (c *CheckAccessControlsTest) TestNoPackageAccessToConstructorThroughExtendsClause_fileOveriew() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/** @fileoverview @package */

/** @public */
class Foo {
  constructor() {}
}
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"), //
				"class SubFoo extends Foo { }")),
		c.Error(CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoPackageAccessToClass_fileOveriew()
func (c *CheckAccessControlsTest) TestNoPackageAccessToClass_fileOveriew() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/** @fileoverview @package */

class Foo {}
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"), //
				"Foo;")),
		c.Error(CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testOverrideWithoutVisibilityRedeclInFileWithFileOverviewVisibilityNotAllowed_OneFile()
func (c *CheckAccessControlsTest) TestOverrideWithoutVisibilityRedeclInFileWithFileOverviewVisibilityNotAllowed_OneFile() {
	c.Test(
		c.Srcs(
			`/**
 * @fileoverview
 * @package
 */

Foo = class {
  /** @private */
  privateMethod_() {}
}

Bar = class extends Foo {
  /** @override */
  privateMethod_() {}
}
`),
		c.Error(CheckAccessControls_BAD_PROPERTY_OVERRIDE_IN_FILE_WITH_FILEOVERVIEW_VISIBILITY))
}

// Java: CheckAccessControlsTest.testOverrideWithoutVisibilityRedeclInFileWithFileOverviewVisibilityNotAllowed_TwoFiles()
func (c *CheckAccessControlsTest) TestOverrideWithoutVisibilityRedeclInFileWithFileOverviewVisibilityNotAllowed_TwoFiles() {
	c.Test(
		c.Srcs(
			`Foo = class {
  /** @protected */
  protectedMethod() {}
};
`,
			`/**
 * @fileoverview
 * @package
 */

Bar = class extends Foo {
  /** @override */
  protectedMethod() {}
};
`),
		c.Error(CheckAccessControls_BAD_PROPERTY_OVERRIDE_IN_FILE_WITH_FILEOVERVIEW_VISIBILITY))
}

// Java: CheckAccessControlsTest.testOverrideWithoutVisibilityRedeclInFileWithNoFileOverviewOk()
func (c *CheckAccessControlsTest) TestOverrideWithoutVisibilityRedeclInFileWithNoFileOverviewOk() {
	c.Test(
		c.Srcs(
			`Foo = class {
  /** @private */
  privateMethod_() {}
};

Bar = class extends Foo {
  /** @override */
  privateMethod_() {}
};
`))
}

// Java: CheckAccessControlsTest.testOverrideWithoutVisibilityRedeclInFileWithNoFileOverviewVisibilityOk()
func (c *CheckAccessControlsTest) TestOverrideWithoutVisibilityRedeclInFileWithNoFileOverviewVisibilityOk() {
	c.Test(
		c.Srcs(
			`/**
 * @fileoverview
 */

Foo = class {
  /** @private */
  privateMethod_() {}
};

Bar = class extends Foo {
  /** @override */
  privateMethod_() {}
};
`))
}

// Java: CheckAccessControlsTest.testOverrideWithVisibilityRedeclInFileWithFileOverviewVisibilityOk_OneFile()
func (c *CheckAccessControlsTest) TestOverrideWithVisibilityRedeclInFileWithFileOverviewVisibilityOk_OneFile() {
	c.Test(
		c.Srcs(
			`/**
 * @fileoverview
 * @package
 */

Foo = class {
  /** @private */
  privateMethod_() {};
};

Bar = class extends Foo {
  /** @override @private */
  privateMethod_() {};
};
`))
}

// Java: CheckAccessControlsTest.testOverrideWithVisibilityRedeclInFileWithFileOverviewVisibilityOk_TwoFiles()
func (c *CheckAccessControlsTest) TestOverrideWithVisibilityRedeclInFileWithFileOverviewVisibilityOk_TwoFiles() {
	c.Test(
		c.Srcs(
			`Foo = class {
  /** @protected */
  protectedMethod() {}
};
`,
			`/**
 * @fileoverview
 * @package
 */

Bar = class extends Foo {
  /** @override @protected */
  protectedMethod() {}
};
`))
}

// Java: CheckAccessControlsTest.testConstructorVisibility_canBeNarrowed()
func (c *CheckAccessControlsTest) TestConstructorVisibility_canBeNarrowed() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @public */
  constructor() {}
};
`,
			`class Bar extends Foo {
  /** @private */
  constructor() {}
};
`))
}

// Java: CheckAccessControlsTest.testConstructorVisibility_canBeExpanded()
func (c *CheckAccessControlsTest) TestConstructorVisibility_canBeExpanded() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @protected */
  constructor() {}
};
`,
			`class Bar extends Foo {
  /** @public */
  constructor() {}
};
`))
}

// Java: CheckAccessControlsTest.testPublicFileOverviewVisibilityDoesNotApplyToNameWithExplicitPackageVisibility()
func (c *CheckAccessControlsTest) TestPublicFileOverviewVisibilityDoesNotApplyToNameWithExplicitPackageVisibility() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/**
 * @fileoverview
 * @public
 */

/** @package */
class Foo {}
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"), //
				"Foo;")),
		c.Error(CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testPackageFileOverviewVisibilityDoesNotApplyToNameWithExplicitPublicVisibility()
func (c *CheckAccessControlsTest) TestPackageFileOverviewVisibilityDoesNotApplyToNameWithExplicitPublicVisibility() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/**
 * @fileoverview
 * @package
 */
/** @public */
class Foo {}
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"), //
				"Foo;")))
}

// Java: CheckAccessControlsTest.testPackageFileOverviewVisibilityAppliesToNameWithoutExplicitVisibility()
func (c *CheckAccessControlsTest) TestPackageFileOverviewVisibilityAppliesToNameWithoutExplicitVisibility() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/**
 * @fileoverview
 * @package
 */

var Foo = class {};
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"), //
				"Foo;")),
		c.Error(CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testPackageFileOverviewVisibilityAppliesToNameWithoutExplicitVisibility_googModule()
func (c *CheckAccessControlsTest) TestPackageFileOverviewVisibilityAppliesToNameWithoutExplicitVisibility_googModule() {
	c.DisableRewriteClosureCode()
	c.TestError(
		c.Srcs(
			[]*SourceFile{
				SourceFile_FromCode(
					Compiler_JoinPathParts("foo", "bar.js"),
					`/**
 * @fileoverview
 * @package
 */
goog.module('Foo');
class Foo {}
exports = Foo;
`),
				SourceFile_FromCode(
					Compiler_JoinPathParts("baz", "quux.js"),
					"goog.module('client'); const Foo = goog.require('Foo'); new Foo();")}),
		c.Error(CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testPackageFileOverviewVisibilityAppliesToNameWithoutExplicitVisibility_esModule()
func (c *CheckAccessControlsTest) TestPackageFileOverviewVisibilityAppliesToNameWithoutExplicitVisibility_esModule() {
	c.TestError(
		c.Srcs(
			[]*SourceFile{
				SourceFile_FromCode(
					Compiler_JoinPathParts("foo", "bar.js"),
					`/**
 * @fileoverview
 * @package
 */
class Foo {}
export {Foo};
`),
				SourceFile_FromCode(
					Compiler_JoinPathParts("baz", "quux.js"),
					"import {Foo} from '/foo/bar.js'; new Foo();")}),
		c.Error(CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testPackageFileOverviewVisibilityDoesNotApplyToPropertyWithExplicitPublicVisibility()
func (c *CheckAccessControlsTest) TestPackageFileOverviewVisibilityDoesNotApplyToPropertyWithExplicitPublicVisibility() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/**
 * @fileoverview
 * @package
 */

/** @public */ // The class must be visible.
Foo = class {
  /** @public */ // The constructor must be visible.
  constructor() { }

  /** @public */
  bar() {}
};
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"),
				`var foo = new Foo();
foo.bar();
`)))
}

// Java: CheckAccessControlsTest.testPublicFileOverviewVisibilityDoesNotApplyToPropertyWithExplicitPackageVisibility()
func (c *CheckAccessControlsTest) TestPublicFileOverviewVisibilityDoesNotApplyToPropertyWithExplicitPackageVisibility() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/**
 * @fileoverview
 * @public
 */

Foo = class {
  /** @package */
  bar() {}
};
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"),
				`var foo = new Foo();
foo.bar();
`)),
		c.Error(CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testPublicFileOverviewVisibilityAppliesToPropertyWithoutExplicitVisibility()
func (c *CheckAccessControlsTest) TestPublicFileOverviewVisibilityAppliesToPropertyWithoutExplicitVisibility() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/**
 * @fileoverview
 * @public
 */

Foo = class {
  bar() {}
}
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"),
				`var foo = new Foo();
foo.bar();
`)))
}

// Java: CheckAccessControlsTest.testPackageFileOverviewVisibilityAppliesToPropertyWithoutExplicitVisibility()
func (c *CheckAccessControlsTest) TestPackageFileOverviewVisibilityAppliesToPropertyWithoutExplicitVisibility() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/**
 * @fileoverview
 * @package
 */

/** @public */ // The class must be visible.
Foo = class {
  /** @public */ // The constructor must be visible.
  constructor() { }

  bar() {}
}
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"),
				`var foo = new Foo();
foo.bar();
`)),
		c.Error(CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testFileOverviewVisibilityComesFromDeclarationFileNotUseFile()
func (c *CheckAccessControlsTest) TestFileOverviewVisibilityComesFromDeclarationFileNotUseFile() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				`/**
 * @fileoverview
 * @package
 */

/** @public */ // The class must be visible.
Foo = class {
  /** @public */ // The constructor must be visible.
  constructor() { }

  bar() {}
}
`),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"),
				`/**
 * @fileoverview
 * @public
 */

var foo = new Foo();
foo.bar();
`)),
		c.Error(CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testImplicitSubclassConstructor_doesNotInheritVisibility_andIsPublic()
func (c *CheckAccessControlsTest) TestImplicitSubclassConstructor_doesNotInheritVisibility_andIsPublic() {
	c.Test(
		c.Srcs(
			"class Foo {\n"+
				"  /** @private */\n"+
				"  constructor() { }\n"+
				"}\n"+
				"\n"+
				"// The implict constructor for `SubFoo` should be treated as public.\n"+
				"class SubFoo extends Foo { }\n",
			// So we get no warning for using it here.
			"new SubFoo();"))
}

// Java: CheckAccessControlsTest.testImplicitSubclassConstructor_doesNotInheritVisibility_andUsesFileOverview()
func (c *CheckAccessControlsTest) TestImplicitSubclassConstructor_doesNotInheritVisibility_andUsesFileOverview() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				"/**\n"+
					" * @fileoverview\n"+
					" * @package\n"+
					" */\n"+
					"\n"+
					"class Foo {\n"+
					"  /** @private */\n"+
					"  constructor() { }\n"+
					"}\n"+
					"\n"+
					"// The implict constructor for `SubFoo` should be trated as package.\n"+
					"/** @public */ // The class must be visible.\n"+
					"class SubFoo extends Foo { }\n"),
			SourceFile_FromCode(Compiler_JoinPathParts("baz", "quux.js"), "new SubFoo();")),
		c.Error(CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testUnannotatedSubclassConstructor_doesNotInheritVisibility_andIsPublic()
func (c *CheckAccessControlsTest) TestUnannotatedSubclassConstructor_doesNotInheritVisibility_andIsPublic() {
	c.Test(
		c.Srcs(
			"class Foo {\n"+
				"  /** @private */\n"+
				"  constructor() { }\n"+
				"}\n"+
				"\n"+
				"class SubFoo extends Foo {\n"+
				"// The unannotated constructor for `SubFoo` should be treated as public.\n"+
				"  constructor() {\n"+
				"    super();\n"+
				"  }\n"+
				"}\n",
			// So we get no warning for using it here.
			"new SubFoo();"))
}

// Java: CheckAccessControlsTest.testUnannotatedSubclassConstructor_doesNotInheritVisibility_andUsesFileOverview()
func (c *CheckAccessControlsTest) TestUnannotatedSubclassConstructor_doesNotInheritVisibility_andUsesFileOverview() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				Compiler_JoinPathParts("foo", "bar.js"),
				"/**\n"+
					" * @fileoverview\n"+
					" * @package\n"+
					" */\n"+
					"\n"+
					"class Foo {\n"+
					"  /** @private */\n"+
					"  constructor() { }\n"+
					"}\n"+
					"\n"+
					"/** @public */ // The class must be visible.\n"+
					"class SubFoo extends Foo {\n"+
					"// The unannotated constructor for `SubFoo` should be treated as package.\n"+
					"  constructor() {\n"+
					"    super();\n"+
					"  }\n"+
					"}\n"),
			SourceFile_FromCode(
				Compiler_JoinPathParts("baz", "quux.js"), //
				"new SubFoo();")),
		c.Error(CheckAccessControls_BAD_PACKAGE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoExceptionsWithBadConstructors1()
func (c *CheckAccessControlsTest) TestNoExceptionsWithBadConstructors1() {
	c.Test(
		c.Srcs(
			`class Foo {
  constructor() {
    (new Bar).bar();
  }
}

class Bar {
  /** @protected */
  bar() {}
}
`))
}

// Java: CheckAccessControlsTest.testNoExceptionsWithBadConstructors2()
func (c *CheckAccessControlsTest) TestNoExceptionsWithBadConstructors2() {
	c.Test(
		c.Srcs(
			`class Foo {
  bar() {}
}

class Bar {
  /** @protected */
  bar() { (new Foo).bar(); }
}
`))
}

// Java: CheckAccessControlsTest.testGoodOverrideOfProtectedProperty()
func (c *CheckAccessControlsTest) TestGoodOverrideOfProtectedProperty() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @protected */
  bar() {}
}
`,
			`class SubFoo extends Foo {
  /** @inheritDoc */
  bar() {}
}
`))
}

// Java: CheckAccessControlsTest.testPublicOverrideOfProtectedProperty()
func (c *CheckAccessControlsTest) TestPublicOverrideOfProtectedProperty() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @protected */
  bar() {}
}
`,
			`class SubFoo extends Foo {
  /** @public */
  bar() {}
}
`))
}

// Java: CheckAccessControlsTest.testBadOverrideOfProtectedProperty()
func (c *CheckAccessControlsTest) TestBadOverrideOfProtectedProperty() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @protected */
  bar() {}
}
`,
			`class SubFoo extends Foo {
  /** @private */
  bar() {}
}
`),
		c.Error(CheckAccessControls_VISIBILITY_MISMATCH))
}

// Java: CheckAccessControlsTest.testProtectedOverrideOfPackageProperty()
func (c *CheckAccessControlsTest) TestProtectedOverrideOfPackageProperty() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @package */
  bar() {}
}
`,
			`class SubFoo extends Foo {
  /** @protected */
  bar() {}
}
`))
}

// Java: CheckAccessControlsTest.testBadOverrideOfPrivateProperty()
func (c *CheckAccessControlsTest) TestBadOverrideOfPrivateProperty() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @private */
  bar() {}
}
`,
			`class SubFoo extends Foo {
  /** @protected */
  bar() {}
}
`),
		c.Error(CheckAccessControls_PRIVATE_OVERRIDE))
}

// Java: CheckAccessControlsTest.testAccessOfStaticMethodOnPrivateClass()
func (c *CheckAccessControlsTest) TestAccessOfStaticMethodOnPrivateClass() {
	c.Test(
		c.Srcs(
			`/** @private */
class Foo {
  static create() { return new Foo(); }
}
`,
			"Foo.create()"),
		c.Error(CheckAccessControls_BAD_PRIVATE_GLOBAL_ACCESS))
}

// Java: CheckAccessControlsTest.testAccessOfStaticMethodOnPrivateQualifiedConstructor()
func (c *CheckAccessControlsTest) TestAccessOfStaticMethodOnPrivateQualifiedConstructor() {
	c.Test(
		c.Srcs(
			`/** @private */
goog.Foo = class {
  static create() { return new goog.Foo(); }
}
`,
			"goog.Foo.create()"),
		c.Error(CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testInstanceofOfPrivateConstructor()
func (c *CheckAccessControlsTest) TestInstanceofOfPrivateConstructor() {
	c.Test(
		c.Srcs(
			`goog.Foo = class {
  /** @private */
  constructor() {}
};
`,
			"goog instanceof goog.Foo"))
}

// Java: CheckAccessControlsTest.testOkAssignmentOfDeprecatedProperty()
func (c *CheckAccessControlsTest) TestOkAssignmentOfDeprecatedProperty() {
	c.Test(
		c.Srcs(
			`class Foo {
  constructor() {
    /** @deprecated */
    this.bar = 3;
  }
}
`))
}

// Java: CheckAccessControlsTest.testBadReadOfDeprecatedProperty()
func (c *CheckAccessControlsTest) TestBadReadOfDeprecatedProperty() {
	c.Test(
		c.Srcs(
			`class Foo {
  constructor() {
    /** @deprecated GRR */
    this.bar = 3;

    this.baz = this.bar;
  }
}
`),
		c.DeprecatedProp("Property bar of type Foo has been deprecated: GRR"))
}

// Java: CheckAccessControlsTest.testNullableDeprecatedProperty()
func (c *CheckAccessControlsTest) TestNullableDeprecatedProperty() {
	c.Test(
		c.Srcs(
			`class Foo {}

/** @deprecated */
Foo.prototype.length;

/** @param {?Foo} x */
function f(x) { return x.length; }
`),
		c.Error(CheckAccessControls_DEPRECATED_PROP))
}

// Java: CheckAccessControlsTest.testNullablePrivateProperty()
func (c *CheckAccessControlsTest) TestNullablePrivateProperty() {
	c.Test(
		c.Srcs(
			`class Foo {}

/** @private */
Foo.prototype.length;
`,
			`/** @param {?Foo} x */
function f(x) { return x.length; }
`),
		c.Error(CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS))
}

// Java: CheckAccessControlsTest.testNoPrivatePropertyByConvention1()
func (c *CheckAccessControlsTest) TestNoPrivatePropertyByConvention1() {
	c.Test(
		c.Srcs(
			`class Foo {}

/** @type {number} */
Foo.prototype.length_;
`,
			`/** @param {?Foo} x */
function f(x) { return x.length_; }
`))
}

// Java: CheckAccessControlsTest.testNoPrivatePropertyByConvention2()
func (c *CheckAccessControlsTest) TestNoPrivatePropertyByConvention2() {
	c.Test(
		c.Srcs(
			`class Foo {
  constructor() {
    /** @type {number} */
    this.length_ = 1;
  }
}

/** @type {number} */
 Foo.prototype.length_;
`,
			`/** @param {Foo} x */
function f(x) { return x.length_; }
`))
}

// Java: CheckAccessControlsTest.testNoDeclarationAndConventionNoConflict1()
func (c *CheckAccessControlsTest) TestNoDeclarationAndConventionNoConflict1() {
	c.Test(
		c.Srcs(
			`class Foo {
  constructor() {
    /** @protected */
    this.length_ = 1;
  }
}
`))
}

// Java: CheckAccessControlsTest.testConstantProperty1a()
func (c *CheckAccessControlsTest) TestConstantProperty1a() {
	c.Test(
		c.Srcs(
			`class A {
  constructor() {
    /** @const */
    this.bar = 3;
  }
}

class B {
  constructor() {
    /** @const */
    this.bar = 3;

    this.bar += 4;
  }
}
`),
		c.Error(CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE))
}

// Java: CheckAccessControlsTest.testConstantProperty1aLogicalAssignment()
func (c *CheckAccessControlsTest) TestConstantProperty1aLogicalAssignment() {
	c.Test(
		c.Srcs(
			`class A {
  constructor() {
    /** @const */
    this.bar = null;

    this.bar ??= 4;
  }
}
`),
		c.Error(CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE))
}

// Java: CheckAccessControlsTest.testConstantPropertyReassigned_crossModuleWithCollidingNames()
func (c *CheckAccessControlsTest) TestConstantPropertyReassigned_crossModuleWithCollidingNames() {
	c.DisableRewriteClosureCode()
	c.TestNoWarning(
		c.Srcs(
			`goog.module('mod1');
class A {
  constructor() {
    /** @const */
    this.bar = 3;
  }
}
`,
			`goog.module('mod2');
class A {
  constructor() {
    /** @const */
    this.bar = 3;
  }
}
`))
}

// Java: CheckAccessControlsTest.testConstantProperty_ConventionNotEnforced1()
func (c *CheckAccessControlsTest) TestConstantProperty_ConventionNotEnforced1() {
	c.Test(
		c.Srcs(
			`class A {
  constructor() {
    this.BAR = 3;
  }
}

class B {
  constructor() {
    this.BAR = 3;

    this.BAR += 4;
  }
}
`))
}

// Java: CheckAccessControlsTest.testConstantProperty_conventionNotEnforced1LogicalAssignment()
func (c *CheckAccessControlsTest) TestConstantProperty_conventionNotEnforced1LogicalAssignment() {
	c.Test(
		c.Srcs(
			`class A {
  constructor() {
    this.BAR = null;

    this.BAR ??= 4;
  }
}
`))
}

// Java: CheckAccessControlsTest.testConstantProperty2()
func (c *CheckAccessControlsTest) TestConstantProperty2() {
	c.Test(
		c.Srcs(
			`class  Foo {}

/** @const */
Foo.prototype.prop = 2;

var foo = new Foo();
foo.prop = 3;
`),
		c.Error(CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE))
}

// Java: CheckAccessControlsTest.testConstantProperty_ConventionNotEnforced2()
func (c *CheckAccessControlsTest) TestConstantProperty_ConventionNotEnforced2() {
	c.Test(
		c.Srcs(
			`class  Foo {}

Foo.prototype.PROP = 2;

var foo = new Foo();
foo.PROP = 3;
`))
}

// Java: CheckAccessControlsTest.testConstantProperty4()
func (c *CheckAccessControlsTest) TestConstantProperty4() {
	c.Test(
		c.Srcs(
			`class Cat {}

/** @const */
Cat.test = 1;

Cat.test *= 2;
`),
		c.Error(CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE))
}

// Java: CheckAccessControlsTest.testConstantProperty_ConventionNotEnforced3()
func (c *CheckAccessControlsTest) TestConstantProperty_ConventionNotEnforced3() {
	c.Test(
		c.Srcs(
			`class Cat { }

Cat.TEST = 1;
Cat.TEST *= 2;
`))
}

// Java: CheckAccessControlsTest.testConstantProperty5()
func (c *CheckAccessControlsTest) TestConstantProperty5() {
	c.Test(
		c.Srcs(
			`class Foo {
  constructor() {
    this.prop = 1;
  }
}

/** @const */
Foo.prototype.prop;

Foo.prototype.prop = 2
`),
		c.Error(CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE))
}

// Java: CheckAccessControlsTest.testConstantProperty6()
func (c *CheckAccessControlsTest) TestConstantProperty6() {
	c.Test(
		c.Srcs(
			`class Foo {
  constructor() {
    this.prop = 1;
  }
}

/** @const */
Foo.prototype.prop = 2;
`),
		c.Error(CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE))
}

// Java: CheckAccessControlsTest.testConstantProperty7()
func (c *CheckAccessControlsTest) TestConstantProperty7() {
	c.Test(
		c.Srcs(
			`class Foo {
  bar_() {}
}

class SubFoo extends Foo {
  /**
   * @const
   * @override
   */
  bar_() {}

  baz() { this.bar_(); }
}
`))
}

// Java: CheckAccessControlsTest.testConstantProperty9()
func (c *CheckAccessControlsTest) TestConstantProperty9() {
	c.Test(
		c.Srcs(
			`class A {
  constructor() {
    /** @const */
    this.bar = 3;
  }
}

class B {
  constructor() {
    this.bar = 4;
  }
}
`))
}

// Java: CheckAccessControlsTest.testConstantProperty10a()
func (c *CheckAccessControlsTest) TestConstantProperty10a() {
	c.Test(
		c.Srcs(
			`class Foo {
  constructor() {
    this.prop = 1;
  }
}

/** @const */
Foo.prototype.prop;
`))
}

// Java: CheckAccessControlsTest.testConstantProperty10b()
func (c *CheckAccessControlsTest) TestConstantProperty10b() {
	c.Test(
		c.Srcs(
			`class Foo {
  constructor() {
    this.PROP = 1;
  }
}

Foo.prototype.PROP;
`))
}

// Java: CheckAccessControlsTest.testConstantProperty11()
func (c *CheckAccessControlsTest) TestConstantProperty11() {
	c.Test(
		c.Srcs(
			`class Foo {}

/** @const */
Foo.prototype.bar;

class SubFoo extends Foo {
  constructor() {
    this.bar = 5;
    this.bar = 6;
  }
}
`),
		c.Error(CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE))
}

// Java: CheckAccessControlsTest.testConstantProperty12()
func (c *CheckAccessControlsTest) TestConstantProperty12() {
	c.Test(
		c.Srcs(
			`class Foo {}

/** @const */
Foo.prototype.bar;

class SubFoo extends Foo {
  constructor() {
    this.bar = 5;
  }
}

class SubFoo2 extends Foo {
  constructor() {
    this.bar = 5;
  }
}
`))
}

// Java: CheckAccessControlsTest.testConstantProperty13()
func (c *CheckAccessControlsTest) TestConstantProperty13() {
	c.Test(
		c.Srcs(
			`class Foo {}

/** @const */
Foo.prototype.bar;

class SubFoo extends Foo {
  constructor() {
    this.bar = 5;
  }
}

class SubSubFoo extends SubFoo {
  constructor() {
    this.bar = 5;
  }
}
`),
		c.Error(CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE))
}

// Java: CheckAccessControlsTest.testConstantProperty14()
func (c *CheckAccessControlsTest) TestConstantProperty14() {
	c.Test(
		c.Srcs(
			`class Foo {
  constructor() {
    /** @const */
    this.bar = 3;

    delete this.bar;
  }
}
`),
		c.Error(CheckAccessControls_CONST_PROPERTY_DELETED))
}

// Java: CheckAccessControlsTest.testConstantProperty_recordType()
func (c *CheckAccessControlsTest) TestConstantProperty_recordType() {
	c.Test(
		c.Srcs(
			`/** @record */
class Foo {
  constructor() {
    /** @const {number} */
    this.bar;
  }
}

const /** !Foo */ x = {
  bar: 9,
};
x.bar = 0;
`),
		c.Error(CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE).
			WithMessageContaining("unknown location due to structural typing"))
}

// Java: CheckAccessControlsTest.testConstantProperty_fromExternsOrIjs_duplicateExternOk()
func (c *CheckAccessControlsTest) TestConstantProperty_fromExternsOrIjs_duplicateExternOk() {
	c.TestSame(
		c.Externs(
			`class Foo {}
/** @const */ Foo.prototype.PROP;
/** @const */ Foo.prototype.PROP;
`),
		c.Srcs(""))
}

// Java: CheckAccessControlsTest.testConstantProperty_fromExternsOrIjs()
func (c *CheckAccessControlsTest) TestConstantProperty_fromExternsOrIjs() {
	c.Test(
		c.Externs(
			`class Foo {}
/** @const */ Foo.prototype.PROP;
`),
		c.Srcs(
			`var f = new Foo();
f.PROP = 1;
`),
		c.Error(CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE).WithMessageContaining("at externs:2:"))
}

// Java: CheckAccessControlsTest.testConstantProperty_ConventionNotEnforced15a()
func (c *CheckAccessControlsTest) TestConstantProperty_ConventionNotEnforced15a() {
	c.Test(
		c.Srcs(
			`class Foo {
  constructor() {
    this.CONST = 100;
  }
}

/** @type {Foo} */
var foo = new Foo();

/** @type {number} */
foo.CONST = 0;
`))
}

// Java: CheckAccessControlsTest.testConstantProperty_ConventionNotEnforced15b()
func (c *CheckAccessControlsTest) TestConstantProperty_ConventionNotEnforced15b() {
	c.Test(
		c.Srcs(
			`class Foo {}

Foo.prototype.CONST = 100;

/** @type {Foo} */
var foo = new Foo();

/** @type {number} */
foo.CONST = 0;
`))
}

// Java: CheckAccessControlsTest.testConstantProperty_ConventionNotEnforced15c()
func (c *CheckAccessControlsTest) TestConstantProperty_ConventionNotEnforced15c() {
	c.Test(
		c.Srcs(
			`class Bar {
  constructor() {
    this.CONST = 100;
  }
}

class Foo extends Bar {};

/** @type {Foo} */
var foo = new Foo();

/** @type {number} */
foo.CONST = 0;
`))
}

// Java: CheckAccessControlsTest.testConstantProperty16()
func (c *CheckAccessControlsTest) TestConstantProperty16() {
	c.Test(
		c.Srcs(
			`class Foo {}

Foo.CONST = 100;

class Bar {}

Bar.CONST = 100;
`))
}

// Java: CheckAccessControlsTest.testConstantProperty17()
func (c *CheckAccessControlsTest) TestConstantProperty17() {
	c.Test(
		c.Srcs(
			`class Foo {
/** @const */
  x = 2;
}

/** @const */
new Foo().x = 2;
`),
		c.Error(CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE))
}

// Java: CheckAccessControlsTest.testConstantPropertyReassignmentInClassStaticBlock()
func (c *CheckAccessControlsTest) TestConstantPropertyReassignmentInClassStaticBlock() {
	c.Test(
		c.Srcs(
			`class Foo {
/** @const */
  static x = 2;
  static {
    this.x = 3;
  }
}
`),
		c.Error(CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE))
}

// Java: CheckAccessControlsTest.testFinalClassCannotBeSubclassed()
func (c *CheckAccessControlsTest) TestFinalClassCannotBeSubclassed() {
	c.Test(
		c.Srcs(
			`/** @final */
var Foo = class {};

var Bar = class extends Foo {};
`),
		c.Error(CheckAccessControls_EXTEND_FINAL_CLASS))
	c.Test(
		c.Srcs(
			`/** @final */
class Foo {};

class Bar extends Foo {};
`),
		c.Error(CheckAccessControls_EXTEND_FINAL_CLASS))
	c.Test(
		c.Srcs(
			`/** @final */
class Foo {};

class Bar extends Foo { constructor() {} };
`),
		c.Error(CheckAccessControls_EXTEND_FINAL_CLASS))
	c.Test(
		c.Srcs(
			`/** @const */
var Foo = class {};

var Bar = class extends Foo {};
`))
}

// Java: CheckAccessControlsTest.testFinalAndConstTreatedAsFinal()
func (c *CheckAccessControlsTest) TestFinalAndConstTreatedAsFinal() {
	c.Test(
		c.Srcs(
			`class Foo {}
/** @final @const */
Foo.prop = 0;
Foo.prop = 1;
`),
		c.Error(CheckAccessControls_FINAL_PROPERTY_OVERRIDDEN))
}

// Java: CheckAccessControlsTest.testFinalMethodCanBeCalled()
func (c *CheckAccessControlsTest) TestFinalMethodCanBeCalled() {
	c.TestNoWarning(
		c.Srcs(
			`class Foo {
  /** @final */
  method() {}
}
class Bar extends Foo {}
new Foo().method();
new Bar().method();
`))
}

// Java: CheckAccessControlsTest.testFinalMethodCannotBeOverridden()
func (c *CheckAccessControlsTest) TestFinalMethodCannotBeOverridden() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @final */
  method() {}
}
class Bar extends Foo {
  method() {}
}
`),
		c.Error(CheckAccessControls_FINAL_PROPERTY_OVERRIDDEN))
	c.Test(
		c.Srcs(
			`class Foo {
  /** @final */
  method() {}
}
class Bar extends Foo {
  /** @override */
  method() {}
}
`),
		c.Error(CheckAccessControls_FINAL_PROPERTY_OVERRIDDEN))
	c.Test(
		c.Srcs(
			`class Foo {
  /** @final */
  method() {}
}
class Bar extends Foo {}
class Baz extends Bar {
  /** @override */
  method() {}
}
`),
		c.Error(CheckAccessControls_FINAL_PROPERTY_OVERRIDDEN))
	c.Test(
		c.Srcs(
			`class Foo {
  /** @final */
  method() {}
}
Foo.prototype.method = () => '';
`),
		c.Error(CheckAccessControls_FINAL_PROPERTY_OVERRIDDEN))
}

// Java: CheckAccessControlsTest.testFinalMethodCannotBeDeleted()
func (c *CheckAccessControlsTest) TestFinalMethodCannotBeDeleted() {
	c.Test(
		c.Srcs(
			`class Foo {
  /** @final */
  method() {}
}
delete Foo.prototype.method;
`),
		c.Error(CheckAccessControls_CONST_PROPERTY_DELETED))
	c.Test(
		c.Srcs(
			`class Foo {
  /** @final */
  method() {}
}
class Bar extends Foo {}
delete Bar.prototype.method;
`),
		c.Error(CheckAccessControls_CONST_PROPERTY_DELETED))
}

// Java: CheckAccessControlsTest.testCircularPrototypeLink()
func (c *CheckAccessControlsTest) TestCircularPrototypeLink() {
	// NOTE: this does yield a useful warning, except we don't check for it in this test:
	//      WARNING - Cycle detected in inheritance chain of type Foo
	// This warning already has a test: TypeCheckTest::testPrototypeLoop.
	c.Test(
		c.Srcs(
			`class Foo extends Foo {}

/** @const */
Foo.prop = 1;

Foo.prop = 2;
`),
		c.Error(CheckAccessControls_CONST_PROPERTY_REASSIGNED_VALUE))
}

// Java: CheckAccessControlsTest.deprecatedName(String)
// (static in Java)
func (c *CheckAccessControlsTest) DeprecatedName(errorMessage string) *CompilerTestCase_Diagnostic {
	return c.Error(CheckAccessControls_DEPRECATED_NAME_REASON).WithMessage(errorMessage)
}

// Java: CheckAccessControlsTest.deprecatedProp(String)
// (static in Java)
func (c *CheckAccessControlsTest) DeprecatedProp(errorMessage string) *CompilerTestCase_Diagnostic {
	return c.Error(CheckAccessControls_DEPRECATED_PROP_REASON).WithMessage(errorMessage)
}

// Java: CheckAccessControlsTest.deprecatedClass(String)
// (static in Java)
func (c *CheckAccessControlsTest) DeprecatedClass(errorMessage string) *CompilerTestCase_Diagnostic {
	return c.Error(CheckAccessControls_DEPRECATED_CLASS_REASON).WithMessage(errorMessage)
}
