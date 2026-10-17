// Port of com/google/javascript/jscomp/ConcretizeStaticInheritanceForInliningTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/ConcretizeStaticInheritanceForInliningTest.java)

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

type ConcretizeStaticInheritanceForInliningTest struct {
	CompilerTestCaseBase
}

// Java: ConcretizeStaticInheritanceForInliningTest()
func NewConcretizeStaticInheritanceForInliningTest() *ConcretizeStaticInheritanceForInliningTest {
	c := &ConcretizeStaticInheritanceForInliningTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestConcretizeStaticInheritanceForInliningTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewConcretizeStaticInheritanceForInliningTest() })
}

// Java: ConcretizeStaticInheritanceForInliningTest.getProcessor(Compiler)
func (c *ConcretizeStaticInheritanceForInliningTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewConcretizeStaticInheritanceForInlining(compiler)
}

// Java: ConcretizeStaticInheritanceForInliningTest.testSimple()
func (c *ConcretizeStaticInheritanceForInliningTest) TestSimple() {
	c.Test(
		"/** @constructor */\n"+
			"function Example() {}\n"+
			"Example.staticMethod = function() { alert(1); }\n"+
			"\n"+
			"/** @constructor @extends {Example} */\n"+
			"function Subclass() {}\n"+
			"$jscomp.inherits(Subclass, Example);\n",
		"/** @constructor */\n"+
			"function Example() {}\n"+
			"Example.staticMethod = function() { alert(1); }\n"+
			"\n"+
			"/** @constructor @extends {Example} */\n"+
			"function Subclass() {}\n"+
			"$jscomp.inherits(Subclass, Example);\n"+
			"\n"+
			"Subclass.staticMethod = Example.staticMethod;\n")
}

// Java: ConcretizeStaticInheritanceForInliningTest.testTyped()
func (c *ConcretizeStaticInheritanceForInliningTest) TestTyped() {
	c.Test(
		"/** @constructor */\n"+
			"function Example() {}\n"+
			"\n"+
			"/** @return {string} */\n"+
			"Example.staticMethod = function() { return ''; }\n"+
			"\n"+
			"/** @constructor @extends {Example} */\n"+
			"function Subclass() {}\n"+
			"$jscomp.inherits(Subclass, Example);\n",
		"/** @constructor */\n"+
			"function Example() {}\n"+
			"\n"+
			"/** @return {string} */\n"+
			"Example.staticMethod = function() { return ''; }\n"+
			"\n"+
			"/** @constructor @extends {Example} */\n"+
			"function Subclass() {}\n"+
			"$jscomp.inherits(Subclass,Example);\n"+
			"\n"+
			"/** @return {string} */\n"+
			"Subclass.staticMethod = Example.staticMethod;\n")
}

// Java: ConcretizeStaticInheritanceForInliningTest.testNestedSubclass()
func (c *ConcretizeStaticInheritanceForInliningTest) TestNestedSubclass() {
	c.Test(
		"/** @constructor */\n"+
			"function A() {}\n"+
			"A.a = 19;\n"+
			"/** @constructor */\n"+
			"A.B = function() {};\n"+
			"/** @constructor @extends {A} */\n"+
			"A.C = function() {};\n"+
			"$jscomp.inherits(A.C, A);\n"+
			"/** @constructor */\n"+
			"A.D = function() {};\n"+
			"A.e = 42;\n",
		"/** @constructor */\n"+
			"function A() {}\n"+
			"A.a = 19;\n"+
			"/** @constructor */\n"+
			"A.B = function() {}\n"+
			"/** @constructor @extends {A} */\n"+
			"A.C = function() {};\n"+
			"$jscomp.inherits(A.C, A);\n"+
			"/** @constructor @extends {A} */\n"+
			"A.C.C = A.C;\n"+
			"/** @constructor */\n"+
			"A.C.B = A.B;\n"+
			"A.C.a = A.a;\n"+
			"/** @constructor */\n"+
			"A.D = function() {};\n"+
			"A.e = 42;\n")
}

// Java: ConcretizeStaticInheritanceForInliningTest.testOverride()
func (c *ConcretizeStaticInheritanceForInliningTest) TestOverride() {
	c.TestSame(
		"/** @constructor */\n" +
			"function Example() {}\n" +
			"\n" +
			"/** @return {string} */\n" +
			"Example.staticMethod = function() { return ''; }\n" +
			"\n" +
			"/** @constructor @extends {Example} */\n" +
			"function Subclass() {}\n" +
			"$jscomp.inherits(Subclass, Example);\n" +
			"\n" +
			"Subclass.staticMethod = function() { return 5; };\n")

	c.TestSame(
		"/** @constructor */\n" +
			"function Example() {}\n" +
			"Example.staticProp = 5;\n" +
			"/** @constructor @extends {Example} */\n" +
			"function Subclass() {}\n" +
			"$jscomp.inherits(Subclass, Example);\n" +
			"Subclass.staticProp = 6;\n")
}

/** In this example, the base class has a static field which is not a function. */
// Java: ConcretizeStaticInheritanceForInliningTest.testStaticNonMethod()
func (c *ConcretizeStaticInheritanceForInliningTest) TestStaticNonMethod() {
	c.Test(
		"/** @constructor */\n"+
			"function Example() {}\n"+
			"\n"+
			"/** @type {number} */\n"+
			"Example.staticField = 5;\n"+
			"\n"+
			"/** @constructor @extends {Example} */\n"+
			"function Subclass() {}\n"+
			"$jscomp.inherits(Subclass, Example);\n",
		"/** @constructor */\n"+
			"function Example() {}\n"+
			"\n"+
			"/** @type {number} */\n"+
			"Example.staticField = 5;\n"+
			"\n"+
			"/** @constructor @extends {Example} */\n"+
			"function Subclass() {}\n"+
			"$jscomp.inherits(Subclass, Example);\n"+
			"\n"+
			"/** @type {number} */\n"+
			"Subclass.staticField = Example.staticField;\n")
}

// Java: ConcretizeStaticInheritanceForInliningTest.testGetterSetterSimple()
func (c *ConcretizeStaticInheritanceForInliningTest) TestGetterSetterSimple() {
	// This is what the Es6ToEs3Converter produces for:
	//
	//   class Example {
	//     static get property() {}
	//   }
	//
	// or
	//
	//   class Example {
	//     static set property(x) {}
	//   }
	c.Test(
		"/** @constructor */\n"+
			"function Example() {}\n"+
			"\n"+
			"/** @type {string} */\n"+
			"Example.property;\n"+
			"Object.defineProperties(Example, {property: { configurable:true, enumerable:true,\n"+
			"    get:function() { return 1; },\n"+
			"    set:function(a) {}\n"+
			"}});\n"+
			"\n"+
			"/** @constructor @extends {Example} */\n"+
			"function Subclass() {}\n"+
			"$jscomp.inherits(Subclass, Example);\n",
		"/** @constructor */\n"+
			"function Example() {}\n"+
			"\n"+
			"/** @type {string} */\n"+
			"Example.property;\n"+
			"Object.defineProperties(Example, {property:{configurable:true, enumerable:true,\n"+
			"  get:function() { return 1; },\n"+
			"  set:function(a) {}\n"+
			"}});\n"+
			"\n"+
			"/** @constructor @extends {Example} */\n"+
			"function Subclass() {}\n"+
			"\n"+
			"$jscomp.inherits(Subclass, Example);\n")
}

// Java: ConcretizeStaticInheritanceForInliningTest.testGetterSetterQualifiedClassName()
func (c *ConcretizeStaticInheritanceForInliningTest) TestGetterSetterQualifiedClassName() {
	c.Test(
		"let TestCase = {};\n"+
			"TestCase.A = /** @constructor */function() {};\n"+
			"\n"+
			"/** @type {string} */\n"+
			"TestCase.A.property;\n"+
			"Object.defineProperties(TestCase.A, {property: { configurable:true, enumerable:true,\n"+
			"    get:function() { return 1; },\n"+
			"    set:function(a) {}\n"+
			"}});\n"+
			"\n"+
			"/** @constructor @extends {TestCase.A} */\n"+
			"function Subclass() {}\n"+
			"$jscomp.inherits(Subclass, TestCase.A);\n",
		"let TestCase = {};\n"+
			"TestCase.A = /** @constructor */function() {};\n"+
			"\n"+
			"/** @type {string} */\n"+
			"TestCase.A.property;\n"+
			"Object.defineProperties(TestCase.A, {property: { configurable:true, enumerable:true,\n"+
			"    get:function() { return 1; },\n"+
			"    set:function(a) {}\n"+
			"}});\n"+
			"\n"+
			"/** @constructor @extends {TestCase.A} */\n"+
			"function Subclass() {}\n"+
			"$jscomp.inherits(Subclass, TestCase.A);\n")
}

/**
 * In this case the stub is not really a stub. It's just a no-op getter, we would be able to
 * detect this and not copy the stub since there is a member with this name.
 */
// Java: ConcretizeStaticInheritanceForInliningTest.testGetterSetterFakeStub()
func (c *ConcretizeStaticInheritanceForInliningTest) TestGetterSetterFakeStub() {
	c.Test(
		"/** @constructor */\n"+
			"function A() {}\n"+
			"\n"+
			"/** @type {string} */\n"+
			"A.property;\n"+
			"A.property = 'string'\n"+
			"\n"+
			"/** @constructor @extends {A} */\n"+
			"function B() {}\n"+
			"$jscomp.inherits(B, A);\n",
		"/** @constructor */\n"+
			"function A() {}\n"+
			"\n"+
			"/** @type {string} */\n"+
			"A.property;\n"+
			"A.property = 'string'\n"+
			"\n"+
			"/** @constructor @extends {A} */\n"+
			"function B() {}\n"+
			"$jscomp.inherits(B, A);\n"+
			"B.property = A.property;\n")
}

// Java: ConcretizeStaticInheritanceForInliningTest.testGetterSetterSubclassSubclass()
func (c *ConcretizeStaticInheritanceForInliningTest) TestGetterSetterSubclassSubclass() {
	c.Test(
		"/** @constructor */\n"+
			"function A() {}\n"+
			"\n"+
			"/** @type {string} */\n"+
			"A.property;\n"+
			"Object.defineProperties(A, {property: { configurable:true, enumerable:true,\n"+
			"    get:function() { return 1; },\n"+
			"    set:function(a) {}\n"+
			"}});\n"+
			"\n"+
			"/** @constructor @extends {A} */\n"+
			"function B() {}\n"+
			"$jscomp.inherits(B, A);\n"+
			"\n"+
			"/** @constructor @extends {B} */\n"+
			"function C() {}\n"+
			"$jscomp.inherits(C, B);\n",
		"/** @constructor */\n"+
			"function A() {}\n"+
			"\n"+
			"/** @type {string} */\n"+
			"A.property;\n"+
			"Object.defineProperties(A, {property: { configurable:true, enumerable:true,\n"+
			"    get:function() { return 1; },\n"+
			"    set:function(a) {}\n"+
			"}});\n"+
			"\n"+
			"/** @constructor @extends {A} */\n"+
			"function B() {}\n"+
			"$jscomp.inherits(B, A);\n"+
			"\n"+
			"/** @constructor @extends {B} */\n"+
			"function C() {}\n"+
			"$jscomp.inherits(C, B);\n")
}

/** If the subclass overrides the property we don't want to redeclare the stub. */
// Java: ConcretizeStaticInheritanceForInliningTest.testGetterSetterSubclassOverride()
func (c *ConcretizeStaticInheritanceForInliningTest) TestGetterSetterSubclassOverride() {
	c.TestSame(
		"/** @constructor */\n" +
			"function A() {}\n" +
			"\n" +
			"/** @type {string} */\n" +
			"A.property;\n" +
			"Object.defineProperties(A, {property: { configurable:true, enumerable:true,\n" +
			"    get:function() { return 1; },\n" +
			"    set:function(a) {}\n" +
			"}});\n" +
			"\n" +
			"/** @constructor @extends {A} */\n" +
			"function B() {}\n" +
			"/** @type {string} */\n" +
			"B.property;\n" +
			"Object.defineProperties(B, {property: { configurable:true, enumerable:true,\n" +
			"    get:function() { return 2; },\n" +
			"    set:function(a) {}\n" +
			"}});\n" +
			"$jscomp.inherits(B, A);\n")

	c.TestSame(
		"/** @constructor */\n" +
			"function A() {}\n" +
			"\n" +
			"/** @type {string} */\n" +
			"A.property;\n" +
			"Object.defineProperties(A, {property: { configurable:true, enumerable:true,\n" +
			"    get:function() { return 1; },\n" +
			"    set:function(a) {}\n" +
			"}});\n" +
			"\n" +
			"/** @constructor @extends {A} */\n" +
			"function B() {}\n" +
			"/** @type {string} */\n" +
			"$jscomp.inherits(B, A);\n" +
			"B.property = 'asdf';\n")
}

// Java: ConcretizeStaticInheritanceForInliningTest.testGetterSetter_noType()
func (c *ConcretizeStaticInheritanceForInliningTest) TestGetterSetter_noType() {
	c.Test(
		"/** @constructor */\n"+
			"function Example() {}\n"+
			"\n"+
			"Example.property;\n"+
			"Object.defineProperties(Example, {property: { configurable:true, enumerable:true,\n"+
			"    get:function() { return 1; },\n"+
			"    set:function(a) {}\n"+
			"}});\n"+
			"\n"+
			"/** @constructor @extends {Example} */\n"+
			"function Subclass() {}\n"+
			"$jscomp.inherits(Subclass, Example);\n",
		"/** @constructor */\n"+
			"function Example() {}\n"+
			"\n"+
			"Example.property;\n"+
			"Object.defineProperties(Example, {property: { configurable:true, enumerable:true,\n"+
			"    get:function() { return 1; },\n"+
			"    set:function(a) {}\n"+
			"}});\n"+
			"\n"+
			"/** @constructor @extends {Example} */\n"+
			"function Subclass() {}\n"+
			"\n"+
			"$jscomp.inherits(Subclass, Example);\n")
}

// Java: ConcretizeStaticInheritanceForInliningTest.testInheritFromExterns()
func (c *ConcretizeStaticInheritanceForInliningTest) TestInheritFromExterns() {
	// Since externs definitions can't be inlined into the code, there is no inlining benefit
	// from running this pass on externs.
	c.TestSame(
		c.Externs(
			"/** @constructor */ function ExternsClass() {}\n"+
				"ExternsClass.m = function() {};\n"),
		c.Srcs(
			"/** @constructor @struct @extends {ExternsClass} */\n"+
				"var CodeClass = function(var_args) {\n"+
				"  ExternsClass.apply(this,arguments)\n"+
				"};\n"+
				"$jscomp.inherits(CodeClass,ExternsClass)\n"))
}

// Java: ConcretizeStaticInheritanceForInliningTest.testAliasing()
func (c *ConcretizeStaticInheritanceForInliningTest) TestAliasing() {
	c.Test(
		"/** @constructor */\n"+
			"function Foo() {}\n"+
			"Foo.prop = 123;\n"+
			"var aliasFoo = Foo;\n"+
			"/** @constructor @extends {aliasFoo} */\n"+
			"function Bar() {}\n"+
			"$jscomp.inherits(Bar, aliasFoo);\n",
		"/** @constructor */\n"+
			"function Foo() {}\n"+
			"Foo.prop = 123;\n"+
			"var aliasFoo = Foo;\n"+
			"/** @constructor @extends {aliasFoo} */\n"+
			"function Bar() {}\n"+
			"$jscomp.inherits(Bar, aliasFoo);\n"+
			"Bar.prop = aliasFoo.prop;\n")

	c.Test(
		"/** @constructor */\n"+
			"function Foo() {}\n"+
			"Foo.prop = 123;\n"+
			"const aliasFoo = Foo; // make sure a const alias works\n"+
			"/** @constructor @extends {aliasFoo} */\n"+
			"function Bar() {}\n"+
			"$jscomp.inherits(Bar, aliasFoo);\n",
		"/** @constructor */\n"+
			"function Foo() {}\n"+
			"Foo.prop = 123;\n"+
			"const aliasFoo = Foo; // make sure a const alias works\n"+
			"/** @constructor @extends {aliasFoo} */\n"+
			"function Bar() {}\n"+
			"$jscomp.inherits(Bar, aliasFoo);\n"+
			"Bar.prop = aliasFoo.prop;\n")

	c.Test(
		"/** @constructor */\n"+
			"function Foo() {}\n"+
			"var aliasFoo = Foo;\n"+
			"aliasFoo.prop = 123;\n"+
			"/** @constructor @extends {Foo} */\n"+
			"function Bar() {}\n"+
			"$jscomp.inherits(Bar, Foo);\n",
		"/** @constructor */\n"+
			"function Foo() {}\n"+
			"var aliasFoo = Foo;\n"+
			"aliasFoo.prop = 123;\n"+
			"/** @constructor @extends {Foo} */\n"+
			"function Bar() {}\n"+
			"$jscomp.inherits(Bar, Foo);\n"+
			"Bar.prop = Foo.prop;\n")
}

// Java: ConcretizeStaticInheritanceForInliningTest.testScopeHandling()
func (c *ConcretizeStaticInheritanceForInliningTest) TestScopeHandling() {
	c.TestSame(
		"/** @constructor */\n" +
			"function Foo() {}\n" +
			"\n" +
			"function f() {\n" +
			"  let Foo = {};\n" +
			"  Foo.prop = 123; // Not a reference to the Foo class, so no change.\n" +
			"}\n" +
			"/** @constructor @extends {Foo} */\n" +
			"function Bar() {}\n" +
			"$jscomp.inherits(Bar, Foo);\n")

	c.TestSame(
		"/** @constructor */\n" +
			"let Foo = function() {} // make sure let works\n" +
			"\n" +
			"function f() {\n" +
			"  let Foo = {};\n" +
			"  Foo.prop = 123; // Not a reference to the Foo class, so no change.\n" +
			"}\n" +
			"/** @constructor @extends {Foo} */\n" +
			"function Bar() {}\n" +
			"$jscomp.inherits(Bar, Foo);\n")

	c.TestSame(
		"/** @constructor */\n" +
			"function Foo() {}\n" +
			"\n" +
			"function f() {\n" +
			"  let Foo = {};\n" +
			"  function g() {\n" +
			"    Foo.prop = 123; // Not a reference to the Foo class, so no change.\n" +
			"  }\n" +
			"}\n" +
			"/** @constructor @extends {Foo} */\n" +
			"function Bar() {}\n" +
			"$jscomp.inherits(Bar, Foo);\n")
}

// Java: ConcretizeStaticInheritanceForInliningTest.testInlineTypes()
func (c *ConcretizeStaticInheritanceForInliningTest) TestInlineTypes() {
	c.Test(
		"/** @constructor @struct */\n"+
			"let A = function() {}\n"+
			"A.foo = function(/** number */ x) {}\n"+
			"\n"+
			"/** @constructor @struct @extends {A} */\n"+
			"var B = function(var_args) { A.apply(this,arguments); };\n"+
			"$jscomp.inherits(B, A);\n",
		"/** @constructor @struct */\n"+
			"let A = function() {}\n"+
			"A.foo = function(/** number */ x) {}\n"+
			"\n"+
			"/** @constructor @struct @extends {A} */\n"+
			"var B = function(var_args) { A.apply(this,arguments); };\n"+
			"$jscomp.inherits(B, A);\n"+
			"B.foo = A.foo;\n")
}

/** Examples which are handled incorrectly but are unlikely to come up in practice. */
// Java: ConcretizeStaticInheritanceForInliningTest.testIncorrectScopeHandling()
func (c *ConcretizeStaticInheritanceForInliningTest) TestIncorrectScopeHandling() {
	c.Test(
		"let example = {};\n"+
			"/** @constructor */\n"+
			"example.Foo = function() {};\n"+
			"\n"+
			"function f() {\n"+
			"  var example = {};\n"+
			"  example.Foo = {};\n"+
			"  // Not a reference to the example.Foo class, so there should be no change.\n"+
			"  example.Foo.prop = 123;\n"+
			"}\n"+
			"/** @constructor @extends {example.Foo} */\n"+
			"function Bar() {}\n"+
			"$jscomp.inherits(Bar, example.Foo);\n",
		"let example = {};\n"+
			"/** @constructor */\n"+
			"example.Foo = function() {};\n"+
			"\n"+
			"function f() {\n"+
			"  var example = {};\n"+
			"  example.Foo = {};\n"+
			"  example.Foo.prop = 123;\n"+
			"}\n"+
			"/** @constructor @extends {example.Foo} */\n"+
			"function Bar() {}\n"+
			"$jscomp.inherits(Bar, example.Foo);\n"+
			"Bar.prop = example.Foo.prop\n")

	c.TestSame(
		"function a() {\n" +
			"  /** @constructor */\n" +
			"  function Foo() {}\n" +
			"  Foo.bar = function() {};\n" +
			"\n" +
			"  /** @constructor @extends {Foo} */\n" +
			"  function Bar() {}\n" +
			"  $jscomp.inherits(Bar, Foo);\n" +
			"  // There should be a declaration for Bar.bar.\n" +
			"}\n")
}

// Java: ConcretizeStaticInheritanceForInliningTest.testPrototypeAssignments()
func (c *ConcretizeStaticInheritanceForInliningTest) TestPrototypeAssignments() {
	c.TestSame(
		"/** @constructor */\n" +
			"function Example() {}\n" +
			"Example.prototype = {}\n" +
			"\n" +
			"/** @constructor @extends {Example} */\n" +
			"function Subclass() {}\n" +
			"$jscomp.inherits(Subclass, Example);\n")
}

// Java: ConcretizeStaticInheritanceForInliningTest.testAddSingletonGetter()
func (c *ConcretizeStaticInheritanceForInliningTest) TestAddSingletonGetter() {
	// Even though concretizing in this case would stil technically be correct, as the call to
	// `addSingletonGetter` would overwrite the concretized property, doing so can cause bad
	// interactions with the inliner to inline the wrong defintion, so backing off is safer.
	// See b/182154150 for details
	c.TestSame(
		"/** @constructor */\n" +
			"function Example() {}\n" +
			"Example.getInstance = function() {}\n" +
			"\n" +
			"/** @constructor @extends {Example} */\n" +
			"function Subclass() {}\n" +
			"$jscomp.inherits(Subclass, Example);\n" +
			"goog.addSingletonGetter(Subclass);\n")
}

// Java: ConcretizeStaticInheritanceForInliningTest.testNonQnameConstructor_doesntPolluteListOfAssignments()
func (c *ConcretizeStaticInheritanceForInliningTest) TestNonQnameConstructor_doesntPolluteListOfAssignments() {
	// Reproduce a bug that once created a nonsensical assignment:
	//   Subclass.staticMethod = Example.staticMethod;
	c.TestSame(
		"const ns = {};\n" +
			"/** @constructor */\n" +
			"ns['NOT_A_NAME'] = function() {};\n" +
			"ns['NOT_A_NAME'].staticMethod = function() { alert(1); }\n" +
			"\n" +
			"/** @constructor */\n" +
			"const Example = function() {}\n" +
			"\n" +
			"/** @constructor @extends {Example} */\n" +
			"function Subclass() {}\n" +
			"$jscomp.inherits(Subclass, Example);\n")
}
