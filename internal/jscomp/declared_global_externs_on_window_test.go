// Port of com/google/javascript/jscomp/DeclaredGlobalExternsOnWindowTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/DeclaredGlobalExternsOnWindowTest.java)

/*
 * Copyright 2014 The Closure Compiler Authors.
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

type DeclaredGlobalExternsOnWindowTest struct {
	CompilerTestCaseBase
}

// Java: DeclaredGlobalExternsOnWindowTest()
func NewDeclaredGlobalExternsOnWindowTest() *DeclaredGlobalExternsOnWindowTest {
	c := &DeclaredGlobalExternsOnWindowTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestDeclaredGlobalExternsOnWindowTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewDeclaredGlobalExternsOnWindowTest() })
}

// Java: DeclaredGlobalExternsOnWindowTest.getProcessor(Compiler)
func (c *DeclaredGlobalExternsOnWindowTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewDeclaredGlobalExternsOnWindow(compiler)
}

// Java: DeclaredGlobalExternsOnWindowTest.setUp()
func (c *DeclaredGlobalExternsOnWindowTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableTypeCheck()
	c.AllowExternsChanges()
	c.EnableRunTypeCheckAfterProcessing()
}

// Java: DeclaredGlobalExternsOnWindowTest.testWindowProperty1a()
func (c *DeclaredGlobalExternsOnWindowTest) TestWindowProperty1a() {
	c.TestExternChanges(
		c.Externs("var window; var a;"), c.Srcs(""), c.Expected("var window;var a;window.a"))
}

// No "var window;" so use "this" instead.
// Java: DeclaredGlobalExternsOnWindowTest.testWindowProperty1b()
func (c *DeclaredGlobalExternsOnWindowTest) TestWindowProperty1b() {
	c.TestExternChanges(c.Externs("var a;"), c.Srcs(""), c.Expected("var a;this.a"))
}

// Java: DeclaredGlobalExternsOnWindowTest.testWindowProperty2()
func (c *DeclaredGlobalExternsOnWindowTest) TestWindowProperty2() {
	c.TestExternChanges(c.Externs(""), c.Srcs("var a"), c.Expected(""))
}

// Java: DeclaredGlobalExternsOnWindowTest.testWindowProperty3a()
func (c *DeclaredGlobalExternsOnWindowTest) TestWindowProperty3a() {
	c.TestExternChanges(
		c.Externs("var window; function f() {}"),
		c.Srcs("var b"),
		c.Expected("var window;function f(){}window.f;"))
}

// No "var window;" so use "this" instead.
// Java: DeclaredGlobalExternsOnWindowTest.testWindowProperty3b()
func (c *DeclaredGlobalExternsOnWindowTest) TestWindowProperty3b() {
	c.TestExternChanges(c.Externs("function f() {}"), c.Srcs("var b"), c.Expected("function f(){}this.f"))
}

// Java: DeclaredGlobalExternsOnWindowTest.testWindowProperty4()
func (c *DeclaredGlobalExternsOnWindowTest) TestWindowProperty4() {
	c.TestExternChanges(c.Externs(""), c.Srcs("function f() {}"), c.Expected(""))
}

// Java: DeclaredGlobalExternsOnWindowTest.testWindowProperty5a()
func (c *DeclaredGlobalExternsOnWindowTest) TestWindowProperty5a() {
	c.TestExternChanges(
		c.Externs("var window; var x = function f() {}"),
		c.Srcs("var b"),
		c.Expected("var window;var x=function f(){};window.x;"))
	c.TestExternChanges(
		c.Externs("var window; var x = function () {}"),
		c.Srcs("var b"),
		c.Expected("var window;var x=function(){};window.x;"))
}

// No "var window;" so use "this" instead.
// Java: DeclaredGlobalExternsOnWindowTest.testWindowProperty5b()
func (c *DeclaredGlobalExternsOnWindowTest) TestWindowProperty5b() {
	c.TestExternChanges(
		c.Externs("var x = function f() {};"),
		c.Srcs("var b"),
		c.Expected("var x=function f(){};this.x"))
}

// Java: DeclaredGlobalExternsOnWindowTest.testWindowProperty5c()
func (c *DeclaredGlobalExternsOnWindowTest) TestWindowProperty5c() {
	c.TestExternChanges(
		c.Externs("var window; var x = ()=>{}"),
		c.Srcs("var b"),
		c.Expected("var window;var x=()=>{};window.x;"))
}

// Java: DeclaredGlobalExternsOnWindowTest.testWindowProperty6()
func (c *DeclaredGlobalExternsOnWindowTest) TestWindowProperty6() {
	c.TestExternChanges(
		c.Externs("var window; /** @const {number} */ var n;"),
		c.Srcs(""),
		c.Expected(
			`var window;
/** @const {number} */ var n;
/** @const {number} @suppress {const,duplicate} */ window.n;
`))
}

// Java: DeclaredGlobalExternsOnWindowTest.testWindowProperty7()
func (c *DeclaredGlobalExternsOnWindowTest) TestWindowProperty7() {
	c.TestExternChanges(
		c.Externs("var window; /** @const */ var ns = {}"),
		c.Srcs(""),
		c.Expected(
			`var window;
/** @const */ var ns = {};
/** @suppress {const,duplicate} @const */ window.ns = ns;
`))
}

// Java: DeclaredGlobalExternsOnWindowTest.testNameAliasing()
func (c *DeclaredGlobalExternsOnWindowTest) TestNameAliasing() {
	c.TestExternChanges(
		c.Externs(
			`var window;
/** @const */
var ns = {};
/** @const */
var ns2 = ns;
`),
		c.Srcs(""),
		c.Expected(
			`var window;
/** @const */
var ns = {};
/** @const */
var ns2 = ns;
/** @suppress {const,duplicate} @const */
window.ns = ns;
/** @suppress {const,duplicate} @const */
window.ns2 = ns;
`))
}

// Java: DeclaredGlobalExternsOnWindowTest.testQualifiedNameAliasing()
func (c *DeclaredGlobalExternsOnWindowTest) TestQualifiedNameAliasing() {
	c.TestExternChanges(
		c.Externs(
			`var window;
/** @const */
var ns = {};
/** @type {number} A very important constant */
ns.THE_NUMBER;
/** @const */
var num = ns.THE_NUMBER;
`),
		c.Srcs(""),
		c.Expected(
			`var window;
/** @const */
var ns = {};
/** @type {number} A very important constant */
ns.THE_NUMBER;
/** @const */
var num = ns.THE_NUMBER;
/** @suppress {const,duplicate} @const */
window.ns=ns;
/** @suppress {const,duplicate} @const */
window.num = ns.THE_NUMBER;
`))
}

// Java: DeclaredGlobalExternsOnWindowTest.testWindowProperty8()
func (c *DeclaredGlobalExternsOnWindowTest) TestWindowProperty8() {
	c.TestExternChanges(
		c.Externs("var window; /** @constructor */ function Foo() {}"),
		c.Srcs(""),
		c.Expected(
			`var window;
/** @constructor */ function Foo(){}
/** @constructor @suppress {const,duplicate} */ window.Foo = Foo;
`))
}

// Java: DeclaredGlobalExternsOnWindowTest.testEnumWindowProperty()
func (c *DeclaredGlobalExternsOnWindowTest) TestEnumWindowProperty() {
	c.TestExternChanges(
		c.Externs("var window; /** @enum {string} */ var Enum = { A: 'str' };"),
		c.Srcs(""),
		c.Expected(
			`var window;
/** @enum {string} */ var Enum = { A: 'str' };
/** @enum {string} @suppress {const,duplicate} */ window.Enum = Enum;
`))
}

/** Test to make sure the compiler knows the type of "window.x" is the same as that of "x". */
// Java: DeclaredGlobalExternsOnWindowTest.testWindowPropertyWithJsDoc()
func (c *DeclaredGlobalExternsOnWindowTest) TestWindowPropertyWithJsDoc() {
	c.TestSame(
		c.Externs(
			CompilerTestCase_MINIMAL_EXTERNS+
				`var window;
/** @type {string} */ var x;
`),
		c.Srcs(
			`/** @param {number} n*/
function f(n) {}
f(window.x);
`),
		c.Warning(TypeValidator_TYPE_MISMATCH_WARNING))
}

// Java: DeclaredGlobalExternsOnWindowTest.testEnum()
func (c *DeclaredGlobalExternsOnWindowTest) TestEnum() {
	c.TestSame(
		c.Externs(CompilerTestCase_MINIMAL_EXTERNS+"/** @enum {string} */ var Enum = {FOO: 'foo', BAR: 'bar'};"),
		c.Srcs(
			`/** @param {Enum} e*/
function f(e) {}
f(window.Enum.FOO);
`))
}

/**
 * Test to make sure that if Foo is a constructor, Foo is considered to be the same type as
 * window.Foo.
 */
// Java: DeclaredGlobalExternsOnWindowTest.testConstructorIsSameType()
func (c *DeclaredGlobalExternsOnWindowTest) TestConstructorIsSameType() {
	c.TestSame(
		c.Externs(
			CompilerTestCase_MINIMAL_EXTERNS+
				`var window;
/** @constructor */ function Foo() {}
`),
		c.Srcs(
			`/** @param {!window.Foo} f*/
function bar(f) {}
bar(new Foo());
`))

	c.TestSame(
		c.Externs(CompilerTestCase_MINIMAL_EXTERNS+"/** @constructor */ function Foo() {}"),
		c.Srcs(
			`/** @param {!Foo} f*/
function bar(f) {}
bar(new window.Foo());
`))
}
