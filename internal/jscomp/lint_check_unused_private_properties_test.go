// Port of com/google/javascript/jscomp/lint/CheckUnusedPrivatePropertiesTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/lint/CheckUnusedPrivatePropertiesTest.java)

/*
 * Copyright 2015 The Closure Compiler Authors.
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

// Tests for {@link CheckUnusedPrivateProperties}
type CheckUnusedPrivatePropertiesTest struct {
	CompilerTestCaseBase
}

func NewCheckUnusedPrivatePropertiesTest() *CheckUnusedPrivatePropertiesTest {
	c := &CheckUnusedPrivatePropertiesTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, checkUnusedPrivatePropertiesTest_EXTERNS)
	return c
}

func TestCheckUnusedPrivatePropertiesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCheckUnusedPrivatePropertiesTest() })
}

var checkUnusedPrivatePropertiesTest_EXTERNS = CompilerTestCase_DEFAULT_EXTERNS +
	`/** @const */ goog.reflect = {};
goog.reflect.object;
/** @constructor */
function Window() {}
Window.prototype.x;
Window.prototype.a;
Window.prototype.ext;
/** @type {!Window} */ var window;
function alert(a) {}
var EXT = {};
EXT.ext;
`

// Java: CheckUnusedPrivatePropertiesTest.setUp()
func (c *CheckUnusedPrivatePropertiesTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableGatherExternProperties()
}

// Java: CheckUnusedPrivatePropertiesTest.getProcessor(Compiler)
func (c *CheckUnusedPrivatePropertiesTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewCheckUnusedPrivateProperties(compiler)
}

// Java: CheckUnusedPrivatePropertiesTest.getOptions()
func (c *CheckUnusedPrivatePropertiesTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetWarningLevel(DiagnosticGroups_LINT_CHECKS, CheckLevel_WARNING)
	return options
}

// Java: CheckUnusedPrivatePropertiesTest.unused(String)
func (c *CheckUnusedPrivatePropertiesTest) Unused(code string) {
	c.TestWarning(code, CheckUnusedPrivateProperties_UNUSED_PRIVATE_PROPERTY)
}

// Java: CheckUnusedPrivatePropertiesTest.used(String)
func (c *CheckUnusedPrivatePropertiesTest) Used(code string) {
	c.TestSame(code)
}

// Java: CheckUnusedPrivatePropertiesTest.testSimpleUnused1()
func (c *CheckUnusedPrivatePropertiesTest) TestSimpleUnused1() {
	c.Unused("/** @private */ this.a = 2")
	c.Used("/** @private */ this.a = 2; alert(this.a);")
}

// Java: CheckUnusedPrivatePropertiesTest.testConstructorPropertyUsed1()
func (c *CheckUnusedPrivatePropertiesTest) TestConstructorPropertyUsed1() {
	c.Unused("/** @constructor */ function C() {} /** @private */ C.prop = 1;")
	c.Used("/** @constructor */ function C() {} /** @private */ C.prop = 1; alert(C.prop);")
	c.Used("/** @constructor */ function C() {} /** @private */ C.prop = 1; function f(a) {a.prop}")
}

// Java: CheckUnusedPrivatePropertiesTest.testClassConstructorPropertyUnused1()
func (c *CheckUnusedPrivatePropertiesTest) TestClassConstructorPropertyUnused1() {
	// Don't ever warn about unused private constructors
	c.Used("class C { /** @private */ constructor() {} }; ")
}

// Java: CheckUnusedPrivatePropertiesTest.testClassPropUnused1()
func (c *CheckUnusedPrivatePropertiesTest) TestClassPropUnused1() {
	// A property defined on "this" can be removed
	c.Unused("class C { constructor() { /** @private */ this.a = 2 } }")
}

// Java: CheckUnusedPrivatePropertiesTest.testClassPropUnused2()
func (c *CheckUnusedPrivatePropertiesTest) TestClassPropUnused2() {
	c.Unused("/** @constructor */ function C() { /** @private */ this.prop = 1; }")
	c.Used("/** @constructor */ function C() {/** @private */ this.prop = 1;} alert(some.prop);")
}

// Java: CheckUnusedPrivatePropertiesTest.testClassMethodUnused1()
func (c *CheckUnusedPrivatePropertiesTest) TestClassMethodUnused1() {
	c.Unused("class C { constructor() {}  /** @private */ method() {} }")
	c.Used(`class C { constructor() {}  /** @private */ method() {} }
 new C().method();`)
}

// Java: CheckUnusedPrivatePropertiesTest.testClassQuotedMethodUnused_noWarning()
func (c *CheckUnusedPrivatePropertiesTest) TestClassQuotedMethodUnused_noWarning() {
	c.Used(`class C { constructor() {}
  /** @private */ ['method']() {} }`)
	c.Used(`class C { constructor() {}
  /** @private */ ['method']() {} }
 new C()['method']();`)
}

// Java: CheckUnusedPrivatePropertiesTest.testSimple2()
func (c *CheckUnusedPrivatePropertiesTest) TestSimple2() {
	// A property defined on "this" can be removed, even when defined
	// as part of an expression
	c.Unused("/** @private */ this.a = 1; this.a = 2, f()")
	c.Unused("/** @private */ this.a = 1; x = (this.a = 2, f())")
	c.Unused("/** @private */ this.a = 1; x = (f(), this.a = 2)")
}

// Java: CheckUnusedPrivatePropertiesTest.testSimple3()
func (c *CheckUnusedPrivatePropertiesTest) TestSimple3() {
	// A property defined on an object other than "this" can not be removed.
	c.Used("y.a = 2")
	// and is seen as a use.
	c.Used("y.a = 2; /** @private */ this.a = 2")
	// Some use of the property "a" appears as a use.
	c.Used("y.a = 2; /** @private */ this.a = 1; alert(x.a)")
}

// Java: CheckUnusedPrivatePropertiesTest.testObjLit()
func (c *CheckUnusedPrivatePropertiesTest) TestObjLit() {
	// A property defined on an object other than "this" is considered a use.
	c.Used("({a:2})")
	// and is seen as a use on 'this'.
	c.Used("({a:0}); /** @private */ this.a = 1;")
	// Some use of the property "a" anywhere is considered a use
	c.Used("x = ({a:0}); /** @private */ this.a = 1; alert(x.a)")
}

// Java: CheckUnusedPrivatePropertiesTest.testExtern()
func (c *CheckUnusedPrivatePropertiesTest) TestExtern() {
	// A property defined in the externs and isn't a warning.
	c.TestSame("this.ext = 2")
}

// Java: CheckUnusedPrivatePropertiesTest.testExport()
func (c *CheckUnusedPrivatePropertiesTest) TestExport() {
	// An exported property can not be removed.
	c.TestSame("this.ext = 2; window['export'] = this.ext;")
	c.TestSame("function f() { this.ext = 2; } window['export'] = this.ext;")
}

// Java: CheckUnusedPrivatePropertiesTest.testAssignOp1()
func (c *CheckUnusedPrivatePropertiesTest) TestAssignOp1() {
	// Properties defined using a compound assignment can be removed if the
	// result of the assignment expression is not immediately used.
	c.Unused("/** @private */ this.x; this.x += 2")
	c.Used("/** @private */ this.x; x = (this.x += 2)")
	c.Used("/** @private */ this.x; this.x += 2; x = this.x;")
	// But, of course, a later use prevents its removal.
	c.Used("/** @private */ this.x; this.x += 2; x.x;")
}

// Java: CheckUnusedPrivatePropertiesTest.testAssignOp2()
func (c *CheckUnusedPrivatePropertiesTest) TestAssignOp2() {
	// Properties defined using a compound assignment can be removed if the
	// result of the assignment expression is not immediately used.
	c.Unused("/** @private */ this.a; this.a += 2, f()")
	c.Unused("/** @private */ this.a; x = (this.a += 2, f())")
	c.Used("/** @private */ this.a; x = (f(), this.a += 2)")
}

// Java: CheckUnusedPrivatePropertiesTest.testInc1()
func (c *CheckUnusedPrivatePropertiesTest) TestInc1() {
	// Increments and Decrements are handled similarly to compound assignments
	// but need a placeholder value when replaced.
	c.Unused("/** @private */ this.x; this.x++")
	c.Used("/** @private */ this.x; x = (this.x++)")
	c.Used("/** @private */ this.x; this.x++; x = this.x;")
	c.Unused("/** @private */ this.x; --this.x")
	c.Used("/** @private */ this.x; x = (--this.x)")
	c.Used("/** @private */ this.x; --this.x; x = this.x;")
}

// Java: CheckUnusedPrivatePropertiesTest.testInc2()
func (c *CheckUnusedPrivatePropertiesTest) TestInc2() {
	// Increments and Decrements are handled similarly to compound assignments
	c.Unused("/** @private */ this.a; this.a++, f()")
	c.Unused("/** @private */ this.a;x = (this.a++, f())")
	c.Used("/** @private */ this.a;x = (f(), this.a++)")
	c.Unused("/** @private */ this.a; --this.a, f()")
	c.Unused("/** @private */ this.a; x = (--this.a, f())")
	c.Used("/** @private */ this.a; x = (f(), --this.a)")
}

// Java: CheckUnusedPrivatePropertiesTest.testJSCompiler_renameProperty()
func (c *CheckUnusedPrivatePropertiesTest) TestJSCompiler_renameProperty() {
	// JSCompiler_renameProperty introduces a use of the property
	c.Used("/** @private */ this.a = 2; x[JSCompiler_renameProperty('a')]")
	c.Used("/** @private */ this.a = 2; JSCompiler_renameProperty('a')")
}

// Java: CheckUnusedPrivatePropertiesTest.testForIn()
func (c *CheckUnusedPrivatePropertiesTest) TestForIn() {
	// This is the basic assumption that this pass makes:
	// that it can warn even if it is used indirectly in a for-in loop
	c.Unused(
		`/** @constructor */ var X = function() {
  /** @private */ this.y = 1;
}
for (var a in new X()) { alert(x[a]) }
`)
}

// Java: CheckUnusedPrivatePropertiesTest.testObjectReflection1()
func (c *CheckUnusedPrivatePropertiesTest) TestObjectReflection1() {
	// Verify reflection prevents warning.
	c.Used(
		`/** @constructor */ function A() {/** @private */ this.foo = 1;}
use(goog.reflect.object(A, {foo: 'foo'}));
`)
	// Verify reflection prevents warning.
	c.Used(
		`/** @const */ var $jscomp = {};
/** @const */ $jscomp.scope = {};
/**
 * @param {!Function} type
 * @param {Object} object
 * @return {Object}
 */
$jscomp.reflectObject = function (type, object) { return object; };
/** @constructor */ function A() {/** @private */ this.foo = 1;}
use($jscomp.reflectObject(A, {foo: 'foo'}));
`)
}

// Java: CheckUnusedPrivatePropertiesTest.testObjectReflection2()
func (c *CheckUnusedPrivatePropertiesTest) TestObjectReflection2() {
	// Any object literal definition prevents warning.
	c.Used(
		`/** @constructor */ function A() {/** @private */  this.foo = 1;}
use({foo: 'foo'});
`)
	// member functions prevent renaming (since we allow them in goog.reflect.object)
	c.Used(
		`/** @constructor */ function A() {/** @private */  this.foo = 1;}
use({foo() {}});
`)
	// computed property doesn't prevent warning
	c.Unused(
		`/** @constructor */ function A() {/** @private */  this.foo = 1;}
use({['foo']: 'foo'});
`)
}

// Java: CheckUnusedPrivatePropertiesTest.testPrototypeProps1()
func (c *CheckUnusedPrivatePropertiesTest) TestPrototypeProps1() {
	c.Unused(
		`/** @constructor */ function A() {this.foo = 1;}
/** @private */ A.prototype.foo = 0;
A.prototype.method = function() {this.foo++};
new A().method()
`)
}

// Java: CheckUnusedPrivatePropertiesTest.testPrototypeProps2()
func (c *CheckUnusedPrivatePropertiesTest) TestPrototypeProps2() {
	// warn about all private properties
	c.Unused(
		`/** @constructor */ function A() {this._foo = 1;}
/** @private */ A.prototype._foo = 0;
A.prototype.method = function() {this._foo++};
new A().method()
`)
}

// Java: CheckUnusedPrivatePropertiesTest.testTypedef()
func (c *CheckUnusedPrivatePropertiesTest) TestTypedef() {
	c.Used(
		`/** @constructor */ function A() {}
/** @private @typedef {string} */ A.typedef_;
`)
}

// Java: CheckUnusedPrivatePropertiesTest.testInterface()
func (c *CheckUnusedPrivatePropertiesTest) TestInterface() {
	c.Used(
		`/** @constructor */ function A() {}
/**
 * @interface
 * @private
 */
A.Interface = function() {};
`)
}

// Java: CheckUnusedPrivatePropertiesTest.testConstructorProperty1()
func (c *CheckUnusedPrivatePropertiesTest) TestConstructorProperty1() {
	c.Unused(
		`/** @constructor */
function C() {}
/** @private */ C.prop = 1;
`)
}

// Java: CheckUnusedPrivatePropertiesTest.testConstructorProperty2()
func (c *CheckUnusedPrivatePropertiesTest) TestConstructorProperty2() {
	c.Used(
		`/** @constructor */
function C() {}
/** @private */ C.prop = 1;
function use(a) {
  alert(a.prop)
};
use(C)
`)
}

// Java: CheckUnusedPrivatePropertiesTest.testInterfaceProperty()
func (c *CheckUnusedPrivatePropertiesTest) TestInterfaceProperty() {
	c.Unused(
		`/** @interface */
function C() {}
/** @private */ C.prop = 1;
`)
}

// Java: CheckUnusedPrivatePropertiesTest.testConstructorProperty_nonQnameClass()
func (c *CheckUnusedPrivatePropertiesTest) TestConstructorProperty_nonQnameClass() {
	c.TestSame(
		`const obj = {};
/** @constructor */ obj['ctor'] = function () {};
/** @private */ obj['ctor'].prop = 1;
`)
}

// Java: CheckUnusedPrivatePropertiesTest.testEs6ClassProperty()
func (c *CheckUnusedPrivatePropertiesTest) TestEs6ClassProperty() {
	c.Unused("class C {} /** @private */ C.prop = 1;")
}
