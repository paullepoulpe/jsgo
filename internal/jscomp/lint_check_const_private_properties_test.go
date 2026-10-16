// Port of com/google/javascript/jscomp/lint/CheckConstPrivatePropertiesTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/lint/CheckConstPrivatePropertiesTest.java)

/*
 * Copyright 2018 The Closure Compiler Authors.
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

type CheckConstPrivatePropertiesTest struct {
	CompilerTestCaseBase
}

func NewCheckConstPrivatePropertiesTest() *CheckConstPrivatePropertiesTest {
	c := &CheckConstPrivatePropertiesTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestCheckConstPrivatePropertiesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCheckConstPrivatePropertiesTest() })
}

// Java: CheckConstPrivatePropertiesTest.getProcessor(Compiler)
func (c *CheckConstPrivatePropertiesTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewCheckConstPrivateProperties(compiler)
}

// Java: CheckConstPrivatePropertiesTest.getOptions()
func (c *CheckConstPrivatePropertiesTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetWarningLevel(DiagnosticGroups_LINT_CHECKS, CheckLevel_WARNING)
	return options
}

// Java: CheckConstPrivatePropertiesTest.testConstructorPropModified1()
func (c *CheckConstPrivatePropertiesTest) TestConstructorPropModified1() {
	c.TestWarning(
		"/** @constructor */ function C() {} /** @private */ C.prop = 1;", CheckConstPrivateProperties_MISSING_CONST_PROPERTY)
	c.TestSame("/** @constructor */ function C() {} /** @private */ C.prop = 1; C.prop = 2;")
}

// Java: CheckConstPrivatePropertiesTest.testConstructorPropModified_es6Class()
func (c *CheckConstPrivatePropertiesTest) TestConstructorPropModified_es6Class() {
	c.TestWarning("class C {} /** @private */ C.prop = 1;", CheckConstPrivateProperties_MISSING_CONST_PROPERTY)
	c.TestSame("class C {} /** @private */ C.prop = 1; C.prop = 2;")
}

// Java: CheckConstPrivatePropertiesTest.testConstructorPropModified_interface()
func (c *CheckConstPrivatePropertiesTest) TestConstructorPropModified_interface() {
	c.TestWarning(
		"/** @interface */ function C() {} /** @private */ C.prop = 1;", CheckConstPrivateProperties_MISSING_CONST_PROPERTY)
	c.TestSame("/** @interface */ function C() {} /** @private */ C.prop = 1; C.prop = 2;")
}

// Java: CheckConstPrivatePropertiesTest.testConstructorPropModified_function()
func (c *CheckConstPrivatePropertiesTest) TestConstructorPropModified_function() {
	c.TestSame("/** @constructor */ function C() {} /** @private */ C.prop = function() {};")
}

// Java: CheckConstPrivatePropertiesTest.testConstructorPropModified_enum()
func (c *CheckConstPrivatePropertiesTest) TestConstructorPropModified_enum() {
	c.TestSame(
		"/** @constructor */ function C() {} /** @enum {number} @private */ C.prop = { A: 1 };")
}

// Java: CheckConstPrivatePropertiesTest.testConstructorPropModified_const()
func (c *CheckConstPrivatePropertiesTest) TestConstructorPropModified_const() {
	c.TestSame("/** @constructor */ function C() {} /** @private @const */ C.prop = 1;")
}

// Java: CheckConstPrivatePropertiesTest.testConstructorPropUnmodified_nonQnameClass()
func (c *CheckConstPrivatePropertiesTest) TestConstructorPropUnmodified_nonQnameClass() {
	c.TestSame(
		`const obj = {};
/** @constructor */ obj['ctor'] = function () {};
/** @private */ obj['ctor'].prop = 1;
`)
}

// Java: CheckConstPrivatePropertiesTest.testConstructorPropModified2()
func (c *CheckConstPrivatePropertiesTest) TestConstructorPropModified2() {
	c.TestWarning(
		"/** @constructor */ function C() { /** @private */ this.foo = 1; } ",
		CheckConstPrivateProperties_MISSING_CONST_PROPERTY)
	c.TestSame("/** @constructor */ function C() { /** @private */ this.foo = 1; this.foo = 2; } ")
	c.TestSame(
		`/** @constructor */ function C() { /** @private */ this.foo = 1; }
C.prototype.bar = function() { this.foo = 2; }
`)
}

// Java: CheckConstPrivatePropertiesTest.testConstructorPropModified_delete()
func (c *CheckConstPrivatePropertiesTest) TestConstructorPropModified_delete() {
	c.TestWarning(
		"/** @constructor */ function C() { /** @private */ this.foo = 1; } ",
		CheckConstPrivateProperties_MISSING_CONST_PROPERTY)
	c.TestSame(
		"/** @constructor */ function C() { /** @private */ this.foo = 1; delete this.foo; } ")
	c.TestSame(
		`/** @constructor */ function C() { /** @private */ this.foo = 1; }
C.prototype.bar = function() { delete this.foo; }
`)
}

// Java: CheckConstPrivatePropertiesTest.testConstructorPropModified_arrayIndexGetProp()
func (c *CheckConstPrivatePropertiesTest) TestConstructorPropModified_arrayIndexGetProp() {
	c.TestWarning(
		"/** @constructor */ function C() { /** @private */ this.foo = [1]; } ",
		CheckConstPrivateProperties_MISSING_CONST_PROPERTY)
	// Even though the underlying data structure is testSame, the reference this['foo'] points to is
	// not modified.
	c.TestWarning(
		"/** @constructor */ function C() { /** @private */ this.foo = [1]; this.foo[0] = 2; } ",
		CheckConstPrivateProperties_MISSING_CONST_PROPERTY)
	c.TestWarning(
		`/** @constructor */ function C() { /** @private */ this.foo = [1]; }
C.prototype.bar = function() { this.foo[0] = 2; }
`,
		CheckConstPrivateProperties_MISSING_CONST_PROPERTY)
}

// Java: CheckConstPrivatePropertiesTest.testConstructorPropModified_objectElement()
func (c *CheckConstPrivatePropertiesTest) TestConstructorPropModified_objectElement() {
	c.TestWarning(
		"/** @constructor */ function C() { /** @private */ this['foo'] = 1; } ",
		CheckConstPrivateProperties_MISSING_CONST_PROPERTY)
	c.TestSame(
		"/** @constructor */ function C() { /** @private */ this['foo'] = 1; this['foo'] = 2; } ")
	c.TestSame(
		`/** @constructor */ function C() { /** @private */ this['foo'] = 1; }
C.prototype.bar = function() { this['foo'] = 2; }
`)
}

// Java: CheckConstPrivatePropertiesTest.testConstructorPropModified_arrayIndexGetElem()
func (c *CheckConstPrivatePropertiesTest) TestConstructorPropModified_arrayIndexGetElem() {
	c.TestWarning(
		"/** @constructor */ function C() { /** @private */ this['foo'] = [1]; } ",
		CheckConstPrivateProperties_MISSING_CONST_PROPERTY)
	// Even though the underlying data structure is testSame, the reference this['foo'] points to is
	// not modified.
	c.TestWarning(
		"/** @constructor */ function C() { /** @private */ this['f'] = [1]; this['f'][0] = 2; } ",
		CheckConstPrivateProperties_MISSING_CONST_PROPERTY)
	c.TestWarning(
		`/** @constructor */ function C() { /** @private */ this['foo'] = [1]; }
C.prototype.bar = function() { this['foo'][0] = 2; }
`,
		CheckConstPrivateProperties_MISSING_CONST_PROPERTY)
}

// Java: CheckConstPrivatePropertiesTest.testConstructorPropModified_mixedObjectAccess()
func (c *CheckConstPrivatePropertiesTest) TestConstructorPropModified_mixedObjectAccess() {
	c.TestSame(
		"/** @constructor */ function C() { /** @private */ this.foo = 1; this['foo'] = 2; } ")
	c.TestSame(
		`/** @constructor */ function C() { /** @private */ this['foo'] = 1; }
C.prototype.bar = function() { this.foo = 2; }
`)
}

// Java: CheckConstPrivatePropertiesTest.testConstructorPropModified_lambda()
func (c *CheckConstPrivatePropertiesTest) TestConstructorPropModified_lambda() {
	c.TestSame(
		`/** @constructor */
function C() {
  /** @private */
  this.foo_ = 2;

  (() => { this.foo_ = 1; })();
}
`)
}

// Java: CheckConstPrivatePropertiesTest.testClassPropModified1()
func (c *CheckConstPrivatePropertiesTest) TestClassPropModified1() {
	c.TestWarning(
		"class C { constructor() { /** @private */ this.a = 2; } }", CheckConstPrivateProperties_MISSING_CONST_PROPERTY)
	c.TestSame("class C { constructor() { /** @private */ this.a = 2; this.a = 3; } }")
	c.TestSame("class C { constructor() { /** @private */ this.a = 2; } foo() { this.a = 3; } }")
}

// Java: CheckConstPrivatePropertiesTest.testClassPropModified_const()
func (c *CheckConstPrivatePropertiesTest) TestClassPropModified_const() {
	c.TestSame("class C { constructor() { /** @private @const */ this.a = 2; } }")
}

// Java: CheckConstPrivatePropertiesTest.testClassPropModified_lambda()
func (c *CheckConstPrivatePropertiesTest) TestClassPropModified_lambda() {
	c.TestSame(
		`class C {
  constructor() {
    /** @private */
    this.foo_ = 2;

    (() => { this.foo_ = 1; })();
  }
}
`)
}

// Java: CheckConstPrivatePropertiesTest.testClassPropModified_assignModify()
func (c *CheckConstPrivatePropertiesTest) TestClassPropModified_assignModify() {
	c.TestWarning(
		"class C { constructor() { /** @private */ this.a = 2; } }", CheckConstPrivateProperties_MISSING_CONST_PROPERTY)
	c.TestSame("class C { constructor() { /** @private */ this.a = 2; this.a += 3; } }")
	c.TestSame("class C { constructor() { /** @private */ this.a = 2; } foo() { this.a += 3; } }")
	c.TestSame("class C { constructor() { /** @private */ this.a = 2; this.a -= 3; } }")
	c.TestSame("class C { constructor() { /** @private */ this.a = 2; } foo() { this.a -= 3; } }")
}

// Java: CheckConstPrivatePropertiesTest.testClassPropModified_increment()
func (c *CheckConstPrivatePropertiesTest) TestClassPropModified_increment() {
	c.TestWarning(
		"class C { constructor() { /** @private */ this.a = 2; } }", CheckConstPrivateProperties_MISSING_CONST_PROPERTY)
	c.TestSame("class C { constructor() { /** @private */ this.a = 2; this.a++; } }")
	c.TestSame("class C { constructor() { /** @private */ this.a = 2; } foo() { this.a++; } }")
}

// Java: CheckConstPrivatePropertiesTest.testClassPropModified_decrement()
func (c *CheckConstPrivatePropertiesTest) TestClassPropModified_decrement() {
	c.TestWarning(
		"class C { constructor() { /** @private */ this.a = 2; } }", CheckConstPrivateProperties_MISSING_CONST_PROPERTY)
	c.TestSame("class C { constructor() { /** @private */ this.a = 2; this.a--; } }")
	c.TestSame("class C { constructor() { /** @private */ this.a = 2; } foo() { this.a--; } }")
}

// Java: CheckConstPrivatePropertiesTest.testClassPropModified_delete()
func (c *CheckConstPrivatePropertiesTest) TestClassPropModified_delete() {
	c.TestWarning(
		"class C { constructor() { /** @private */ this.a = 2; } }", CheckConstPrivateProperties_MISSING_CONST_PROPERTY)
	c.TestSame("class C { constructor() { /** @private */ this.a = 2; delete this.a; } }")
	c.TestSame("class C { constructor() { /** @private */ this.a = 2; } foo() { delete this.a; } }")
}

// Java: CheckConstPrivatePropertiesTest.testClassPropModified_multiFile()
func (c *CheckConstPrivatePropertiesTest) TestClassPropModified_multiFile() {
	c.TestWarning(
		c.Srcs(
			"class A { constructor() { /** @private */ this.a = 2; } }", //
			"class B { constructor() { /** @private */ this.a = 1; this.a = 4; } }"),
		CheckConstPrivateProperties_MISSING_CONST_PROPERTY)
	c.TestSame(
		c.Srcs(
			"class A { constructor() { /** @private */ this.a = 2; this.a = 3; } }", //
			"class B { constructor() { /** @private */ this.a = 2; this.a = 4; } }"))
}

// Java: CheckConstPrivatePropertiesTest.testClassPropModified_destructuring()
func (c *CheckConstPrivatePropertiesTest) TestClassPropModified_destructuring() {
	c.TestSame("class C { constructor() { /** @private */ this.a = 2; [this.a] = []; } }")
	c.TestSame("class C { constructor() { /** @private */ this.a = 2; ({a: this.a} = {}); } }")
	c.TestSame("class C { constructor() { /** @private */ this.a = 2; [{a: this.a}] = [{}]; } }")
	c.TestSame("class C { constructor() { /** @private */ this.a = 2; [this.a = 0] = []; } }")
}

// Java: CheckConstPrivatePropertiesTest.testPrototype_Property()
func (c *CheckConstPrivatePropertiesTest) TestPrototype_Property() {
	c.TestSame("/** @constructor */ function C() {} /** @private */ C.prototype.prop = 1;")
	c.TestSame(
		`/** @constructor */ function C() {} /** @private */ C.prototype.prop = 1;
C.prototype.prop = 2;
`)
}

// Java: CheckConstPrivatePropertiesTest.testPrototype_Method()
func (c *CheckConstPrivatePropertiesTest) TestPrototype_Method() {
	c.TestSame(
		"/** @constructor */ function C() {} /** @private */ C.prototype.method = function() {};")
}

// Java: CheckConstPrivatePropertiesTest.testPropertyAsRvalueOfDeclaration()
func (c *CheckConstPrivatePropertiesTest) TestPropertyAsRvalueOfDeclaration() {
	// Test bug where the rhs of a @private prop declaration was treated like a declaration.
	// See b/110938671
	// In the repro below, there was a spurious warning that Foo.SOME_PROP was an effectively
	// constant private property, and the JSDoc "/** @private {number} */" should be @const.
	c.TestNoWarning(
		`class Foo {}
class C {
  constructor() {
    /** @private {number} */
    this.nonConstProp_ = Foo.SOME_PROP;
    this.nonConstProp_ = 4; // change this.nonConstProp_ to avoid a warning
  }
}
`)
}
