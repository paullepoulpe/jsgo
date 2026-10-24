// Port of com/google/javascript/jscomp/ImplicitNullabilityCheckTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/ImplicitNullabilityCheckTest.java)

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

// ImplicitNullabilityCheckTest holds the unit tests for
// ImplicitNullabilityCheck.
type ImplicitNullabilityCheckTest struct {
	CompilerTestCaseBase
}

// Java: ImplicitNullabilityCheckTest()
func NewImplicitNullabilityCheckTest() *ImplicitNullabilityCheckTest {
	c := &ImplicitNullabilityCheckTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestImplicitNullabilityCheckTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewImplicitNullabilityCheckTest() })
}

// Java: ImplicitNullabilityCheckTest.getOptions()
func (c *ImplicitNullabilityCheckTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetWarningLevel(DiagnosticGroups_ANALYZER_CHECKS, CheckLevel_WARNING)
	return options
}

// Java: ImplicitNullabilityCheckTest.getProcessor(Compiler)
func (c *ImplicitNullabilityCheckTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewImplicitNullabilityCheck(compiler)
}

// Java: ImplicitNullabilityCheckTest.setUp()
func (c *ImplicitNullabilityCheckTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableCreateModuleMap()
	c.EnableTypeCheck()
}

// Java: ImplicitNullabilityCheckTest.testExplicitJsdocDoesntWarn()
func (c *ImplicitNullabilityCheckTest) TestExplicitJsdocDoesntWarn() {
	c.NoWarning("/** @type {bigint} */ var x;")
	c.NoWarning("/** @type {boolean} */ var x;")
	c.NoWarning("/** @type {symbol} */ var x;")
	c.NoWarning("/** @type {null} */ var x;")
	c.NoWarning("/** @type {!Object} */ var x;")
	c.NoWarning("/** @type {?Object} */ var x;")
	c.NoWarning("/** @type {function(new:Object)} */ function f(){}")
	c.NoWarning("/** @type {function(this:Object)} */ function f(){}")
	c.NoWarning("/** @typedef {!Object} */ var Obj; var /** ?Obj */ x;")

	// Test let and const
	c.NoWarning("/** @type {boolean} */ let x;")
	c.NoWarning("/** @type {!Object} */ let x;")

	c.NoWarning("/** @type {!Object} */ const x = {};")
	c.NoWarning("/** @type {boolean} */ const x = true;")
}

// Java: ImplicitNullabilityCheckTest.testExplicitlyNullableUnion()
func (c *ImplicitNullabilityCheckTest) TestExplicitlyNullableUnion() {
	c.NoWarning("/** @type {(Object|null)} */ var x;")
	c.NoWarning("/** @type {(Object|number)?} */ var x;")
	c.NoWarning("/** @type {?(Object|number)} */ var x;")
	c.NoWarning("/** @type {(Object|?number)} */ var x;")
	c.WarnImplicitlyNullable("/** @type {(Object|number)} */ var x;")

	c.NoWarning("/** @type {(Object|null)} */ let x;")
	c.WarnImplicitlyNullable("/** @type {(Object|number)} */ let x;")

	c.NoWarning("/** @type {(Object|null)} */ const x = null;")
	c.WarnImplicitlyNullable("/** @type {(Object|number)} */ const x = 3;;")
}

// Java: ImplicitNullabilityCheckTest.testJsdocPositions()
func (c *ImplicitNullabilityCheckTest) TestJsdocPositions() {
	c.WarnImplicitlyNullable("/** @type {Object} */ var x;")
	c.WarnImplicitlyNullable("var /** Object */ x;")
	c.WarnImplicitlyNullable("/** @typedef {Object} */ var x;")
	c.WarnImplicitlyNullable("/** @param {Object} x */ function f(x){}")
	c.WarnImplicitlyNullable("/** @return {Object} */ function f(x){ return {}; }")

	c.WarnImplicitlyNullable("/** @type {Object} */ let x;")
	c.WarnImplicitlyNullable("/** @type {Object} */ const x = {};")
}

// Java: ImplicitNullabilityCheckTest.testParameterizedObject()
func (c *ImplicitNullabilityCheckTest) TestParameterizedObject() {
	c.WarnImplicitlyNullable(
		"/** @param {Object<string, string>=} opt_values */\n" +
			"function getMsg(opt_values) {};\n")
}

// Java: ImplicitNullabilityCheckTest.testNullableTypedef()
func (c *ImplicitNullabilityCheckTest) TestNullableTypedef() {
	// Arguable whether or not this deserves a warning
	c.WarnImplicitlyNullable("/** @typedef {?number} */ var Num; var /** Num */ x;")
}

// Java: ImplicitNullabilityCheckTest.testNonnullTypedef()
func (c *ImplicitNullabilityCheckTest) TestNonnullTypedef() {
	c.Test(
		c.Srcs("/** @typedef {number} */ var Num; var /** Num */ x;"),
		c.Warning(ImplicitNullabilityCheck_IMPLICITLY_NONNULL_JSDOC))
}

// Java: ImplicitNullabilityCheckTest.testNonnullEnum()
func (c *ImplicitNullabilityCheckTest) TestNonnullEnum() {
	c.Test(
		c.Srcs("/** @enum {number} */ var Enum = {NUM: 0}; var /** Enum */ x;"),
		c.Warning(ImplicitNullabilityCheck_IMPLICITLY_NONNULL_JSDOC))
}

// Java: ImplicitNullabilityCheckTest.testNonnullRecordType()
func (c *ImplicitNullabilityCheckTest) TestNonnullRecordType() {
	c.NoWarning("var /** {x: number} */ o;")
}

// Java: ImplicitNullabilityCheckTest.testUnknownTypenameDoesntWarn()
func (c *ImplicitNullabilityCheckTest) TestUnknownTypenameDoesntWarn() {
	c.Test(
		c.Externs(CompilerTestCase_DEFAULT_EXTERNS),
		c.Srcs("/** @type {gibberish} */ var x;"),
		c.Warning(RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR))
}

// Java: ImplicitNullabilityCheckTest.testThrowsDoesntWarn()
func (c *ImplicitNullabilityCheckTest) TestThrowsDoesntWarn() {
	c.NoWarning("/** @throws {Error} */ function f() {}")
	c.NoWarning("/** @throws {TypeError}\n * @throws {SyntaxError} */ function f() {}")
}

// Java: ImplicitNullabilityCheckTest.testTypeofDoesntWarn()
func (c *ImplicitNullabilityCheckTest) TestTypeofDoesntWarn() {
	c.NoWarning("/** @type {typeof Object} */ var x;")
}

// Java: ImplicitNullabilityCheckTest.testUserDefinedClass()
func (c *ImplicitNullabilityCheckTest) TestUserDefinedClass() {
	c.WarnImplicitlyNullable(
		"/** @constructor */\n" +
			"function Foo() {}\n" +
			"/** @type {Foo} */ var x;\n")

	c.WarnImplicitlyNullable(
		"function f() {\n" +
			"  /** @constructor */\n" +
			"  function Foo() {}\n" +
			"  /** @type {Foo} */ var x;\n" +
			"}\n")
}

// Java: ImplicitNullabilityCheckTest.testNamespacedTypeDoesntCrash()
func (c *ImplicitNullabilityCheckTest) TestNamespacedTypeDoesntCrash() {
	c.WarnImplicitlyNullable(
		"/** @const */ var a = {};\n" +
			"/** @const */ a.b = {};\n" +
			"/** @constructor */ a.b.Foo = function() {};\n" +
			"/** @type Array<!a.b.Foo> */ var foos = [];\n")

	// in goog.module
	c.Test(
		c.Externs(NewTestExternsBuilder().AddClosureExterns().Build()),
		c.Srcs(
			"goog.module('ns');\n"+
				"/** @const */ var a = {};\n"+
				"/** @const */ a.b = {};\n"+
				"/** @constructor */ a.b.Foo = function() {};\n"+
				"/** @type Array<!a.b.Foo> */ var foos = [];\n"+
				"exports = {a};\n"),
		c.Warning(ImplicitNullabilityCheck_IMPLICITLY_NULLABLE_JSDOC))

	// in goog.module with a goog.require
	c.Test(
		c.Externs(NewTestExternsBuilder().AddClosureExterns().Build()),
		c.Srcs(
			"goog.module('ns');\n"+
				"/** @const */ var a = {};\n"+
				"/** @const */ a.b = {};\n"+
				"exports = {a};\n",
			"goog.module('ns2');\n"+
				"const {a} = goog.require('ns');\n"+
				"/** @constructor */ a.b.Foo = function() {};\n"+
				"/** @type Array<!a.b.Foo> */ var foos = [];\n"+
				"exports = {foos};\n"),
		c.Warning(ImplicitNullabilityCheck_IMPLICITLY_NULLABLE_JSDOC))

	// in goog.provide
	c.Test(
		c.Externs(NewTestExternsBuilder().AddClosureExterns().Build()),
		c.Srcs(
			"goog.provide('ns.a');\n"+
				"/** @const */ ns.a.b = {};\n"+
				"/** @constructor */ ns.a.b.Foo = function() {};\n"+
				"/** @type Array<!ns.a.b.Foo> */ var foos = [];\n"),
		c.Warning(ImplicitNullabilityCheck_IMPLICITLY_NULLABLE_JSDOC))
}

// TestShadowingGlobalTypedefWithClass is the @Test method
// shadowingGlobalTypedefWithClass (Go-only: language-adapter — RunJUnit4
// runs the methods named Test*).
//
// Java: ImplicitNullabilityCheckTest.shadowingGlobalTypedefWithClass()
func (c *ImplicitNullabilityCheckTest) TestShadowingGlobalTypedefWithClass() {
	c.WarnImplicitlyNullable(
		"/** @typedef {string} */ let Type;\n" +
			"function local() {\n" +
			"  class Type {}\n" +
			"  /** @type {Type} */ var x;\n" +
			"}\n")
}

// TestShadowingGlobalClassWithTypedef is the @Test method
// shadowingGlobalClassWithTypedef (Go-only: language-adapter — RunJUnit4
// runs the methods named Test*).
//
// Java: ImplicitNullabilityCheckTest.shadowingGlobalClassWithTypedef()
func (c *ImplicitNullabilityCheckTest) TestShadowingGlobalClassWithTypedef() {
	c.Test(
		c.Srcs(
			"class Type {}\n"+
				"function local() {\n"+
				"  /** @typedef {string} */ let Type;\n"+
				"  /** @type {Type} */ var x;\n"+
				"}\n"),
		c.Warning(ImplicitNullabilityCheck_IMPLICITLY_NONNULL_JSDOC))
}

// Java: ImplicitNullabilityCheckTest.warnImplicitlyNullable(String)
func (c *ImplicitNullabilityCheckTest) WarnImplicitlyNullable(js string) {
	c.Test(
		c.Externs(CompilerTestCase_DEFAULT_EXTERNS),
		c.Srcs(js),
		c.Warning(ImplicitNullabilityCheck_IMPLICITLY_NULLABLE_JSDOC))
}

// Java: ImplicitNullabilityCheckTest.noWarning(String)
func (c *ImplicitNullabilityCheckTest) NoWarning(js string) {
	c.TestSame(c.Externs(CompilerTestCase_DEFAULT_EXTERNS), c.Srcs(js))
}
