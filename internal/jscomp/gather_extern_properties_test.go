// Port of com/google/javascript/jscomp/GatherExternPropertiesTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/GatherExternPropertiesTest.java)

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

// GatherExternPropertiesTest is the test case for GatherExternProperties.
type GatherExternPropertiesTest struct {
	CompilerTestCaseBase

	// nil is Java null.
	mode *GatherExternProperties_Mode
}

// Java: GatherExternPropertiesTest.EXTERNS
const GatherExternPropertiesTest_EXTERNS = `/**
 * @constructor
 * @param {*=} opt_value
 * @return {!Object}
 */
function Object(opt_value) {}
/**
 * @constructor
 * @param {...*} var_args
 */
function Function(var_args) {}
/**
 * @constructor
 * @param {*=} arg
 * @return {string}
 */
function String(arg) {}
/**
 * @template T
 * @constructor
 * @param {...*} var_args
 * @return {!Array<?>}
 */
function Array(var_args) {}
`

// Java: GatherExternPropertiesTest()
func NewGatherExternPropertiesTest() *GatherExternPropertiesTest {
	c := &GatherExternPropertiesTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, GatherExternPropertiesTest_EXTERNS)
	return c
}

func TestGatherExternPropertiesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewGatherExternPropertiesTest() })
}

// Java: GatherExternPropertiesTest.getProcessor(Compiler)
func (c *GatherExternPropertiesTest) GetProcessor(compiler *Compiler) CompilerPass {
	if c.mode == nil {
		c.mode = gatherExternPropertiesTest_modeOf(GatherExternProperties_Mode_OPTIMIZE)
	}
	return NewGatherExternProperties(compiler, *c.mode)
}

// gatherExternPropertiesTest_modeOf boxes a Mode for the nullable field.
//
// Go-only: test-helper.
func gatherExternPropertiesTest_modeOf(m GatherExternProperties_Mode) *GatherExternProperties_Mode {
	return &m
}

// Java: GatherExternPropertiesTest.testGatherExternPropertiesNullishCoalesce()
func (c *GatherExternPropertiesTest) TestGatherExternPropertiesNullishCoalesce() {
	c.AssertExternProperties("foo = {bar: {} ?? 2}", "bar")
}

// Java: GatherExternPropertiesTest.testGatherExternProperties()
func (c *GatherExternPropertiesTest) TestGatherExternProperties() {
	// Properties.
	c.AssertExternProperties("foo.bar;", "bar")

	// Object literals.
	c.AssertExternProperties("foo = {bar: null, 'baz': {foobar: null}};", "bar", "baz", "foobar")
	// Object literal with numeric propertic.
	c.AssertExternProperties("foo = {0: null};", "0")

	// Top-level variables do not count.
	c.AssertExternProperties("var foo;")

	// String-key access does not count.
	c.AssertExternProperties("foo['bar'] = {};")
}

// Java: GatherExternPropertiesTest.testGatherExternTypedefProperties()
func (c *GatherExternPropertiesTest) TestGatherExternTypedefProperties() {
	c.mode = gatherExternPropertiesTest_modeOf(GatherExternProperties_Mode_CHECK_AND_OPTIMIZE)
	typedefExtern :=
		// Quotes around subTypedefProp should be stripped
		`/**
 * @typedef {{
 *    typedefPropA: { 'subTypedefProp': string },
 *  }}
 */
var TypedefExtern;
/**
 * @param {!{ paramProp1, paramProp2: number }} p
 */
function externFunction(p) {}
`
	c.AssertExternProperties(
		typedefExtern, "typedefPropA", "subTypedefProp", "paramProp1", "paramProp2")
}

// Java: GatherExternPropertiesTest.testGatherExternPropertiesIncludingRecordTypes()
func (c *GatherExternPropertiesTest) TestGatherExternPropertiesIncludingRecordTypes() {
	c.mode = gatherExternPropertiesTest_modeOf(GatherExternProperties_Mode_CHECK_AND_OPTIMIZE)
	// Properties.
	c.AssertExternProperties("foo.bar;", "bar")

	// Object literals.
	c.AssertExternProperties("foo = {bar: null, 'baz': {foobar: null}};", "bar", "baz", "foobar")

	// Object literal with numeric propertic.
	c.AssertExternProperties("foo = {0: null};", "0")

	// Top-level variables do not count.
	c.AssertExternProperties("var foo;")

	// String-key access does not count.
	c.AssertExternProperties("foo['bar'] = {};")

	// Record types on properties.
	c.AssertExternProperties("/** @type {{bar: string, baz: string}} */ var foo;", "bar", "baz")

	// Record types in typedef.
	c.AssertExternProperties(
		"/** @typedef {{bar: string, baz: string}} */ var FooType;", "bar", "baz")

	// Record types in type unions.
	c.AssertExternProperties(
		"/** @type {string|{bar: string}|{baz: string}} */ var foo;", "bar", "baz")

	// Record types in function parameters and return types.
	c.AssertExternProperties(
		`/** @type {function(string, {bar: string}): {baz: string}} */
var foo;
`,
		"bar",
		"baz")

	// Record types as template arguments.
	c.AssertExternProperties(
		"/** @type {Array<{bar: string, baz: string}>} */ var foo;", "bar", "baz")

	// Record types in implemented interfaces.
	c.AssertExternProperties(
		`/**
 * @interface
 * @template T
 */
var Foo = function() {};
/**
 * @constructor
 * @implements {Foo<{bar: string, baz: string}>}
 */
var Bar;
`,
		"bar",
		"baz")

	// Record types in extended class.
	c.AssertExternProperties(
		`/**
 * @constructor
 * @template T
 */
var Foo = function() {};
/**
 * @constructor
 * @extends {Foo<{bar: string, baz: string}>}
 */
var Bar = function() {};
`,
		"bar",
		"baz")

	// Record types in enum.
	// Note that "baz" exists only in the type of the enum,
	// but it is still picked up.
	c.AssertExternProperties(
		`/** @enum {{bar: string, baz: (string|undefined)}} */
var FooEnum = {VALUE: {bar: ''}};
`,
		"VALUE",
		"bar",
		"baz")

	// Nested record types.
	c.AssertExternProperties(
		"/** @type {{bar: string, baz: {foobar: string}}} */ var foo;", "bar", "baz", "foobar")

	// Recursive @record types.
	c.AssertExternProperties(
		`/** @record */
function D1() { /** @type {D2} */ this.a; }

/** @record */
function D2() { /** @type {D1} */ this.b; }
`,
		"a",
		"b")
	c.AssertExternProperties(
		`/** @record */
function D1() { /** @type {function(D2)} */ this.a; }

/** @record */
function D2() { /** @type {D1} */ this.b; }
`,
		"a",
		"b")

	// Recursive types
	c.AssertExternProperties(
		`/** @typedef {{a: D2}} */
var D1;

/** @typedef {{b: D1}} */
var D2;
`,
		"a",
		"b")
	c.AssertExternProperties(
		`/** @typedef {{a: function(D2)}} */
var D1;

/** @typedef {{b: D1}} */
var D2;
`,
		"a",
		"b")

	// Record types defined in normal code and referenced in externs should
	// not bleed-through.
	c.TestSame(
		c.Externs(GatherExternPropertiesTest_EXTERNS+"/** @type {NonExternType} */ var foo;"),
		c.Srcs("/** @typedef {{bar: string, baz: string}} */ var NonExternType;"),
		// Check that no properties were found.
		GatherExternPropertiesTest_ExpectExterns())
}

// Java: GatherExternPropertiesTest.testExternClass()
func (c *GatherExternPropertiesTest) TestExternClass() {
	c.AssertExternProperties(
		`class Foo {
  bar() {
    return this;
  }
}
var baz = new Foo();
var bar = baz.bar;
`,
		"bar")
}

// Java: GatherExternPropertiesTest.testExternWithMethod()
func (c *GatherExternPropertiesTest) TestExternWithMethod() {
	c.AssertExternProperties(
		`foo = {
  method() {}
}
`,
		"method")
}

// Java: GatherExternPropertiesTest.testGatherExternsInCheckMode()
func (c *GatherExternPropertiesTest) TestGatherExternsInCheckMode() {
	c.mode = gatherExternPropertiesTest_modeOf(GatherExternProperties_Mode_CHECK)
	c.AssertExternProperties(
		`/** @fileoverview @externs */
var ns = {};
ns.x;
/** @type {{y: string}} */
ns.yObj;
`,
		"y")
}

// Java: GatherExternPropertiesTest.testGatherExternsInOptimizeMode()
func (c *GatherExternPropertiesTest) TestGatherExternsInOptimizeMode() {
	c.mode = gatherExternPropertiesTest_modeOf(GatherExternProperties_Mode_OPTIMIZE)
	c.AssertExternProperties(
		`/** @fileoverview @externs */
var ns = {};
ns.x;
/** @type {{y: string}} */
ns.yObj;
`,
		"x",
		"yObj")
}

// Java: GatherExternPropertiesTest.testExternClassField()
func (c *GatherExternPropertiesTest) TestExternClassField() {
	c.AssertExternProperties(
		`class Foo {
  /** @type {string} */
  helloWorld;
  constructor() {
    this.goodbye;
  }
}
`,
		"helloWorld",
		"goodbye",
		"constructor")
}

// Java: GatherExternPropertiesTest.testExternClassGetterSetter()
func (c *GatherExternPropertiesTest) TestExternClassGetterSetter() {
	c.AssertExternProperties(
		`class Foo {
  get bar() { return 1; }
  set baz(v) {}
}
`,
		"bar",
		"baz")
}

// Java: GatherExternPropertiesTest.expectExterns(String...)
func GatherExternPropertiesTest_ExpectExterns(properties ...string) CompilerTestCase_Postcondition {
	return CompilerTestCase_PostconditionFunc(func(compiler *Compiler) {
		AssertThat(compiler.GetExternProperties()).ContainsExactlyElementsIn(properties)
	})
}

// Java: GatherExternPropertiesTest.assertExternProperties(String, String...)
func (c *GatherExternPropertiesTest) AssertExternProperties(externs string, properties ...string) {
	c.TestSame(c.Externs(GatherExternPropertiesTest_EXTERNS+externs), c.Srcs(""), GatherExternPropertiesTest_ExpectExterns(properties...))
}
