// Port of com/google/javascript/jscomp/lint/CheckNullabilityModifiersTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/lint/CheckNullabilityModifiersTest.java)

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

// Unit tests for {@link CheckNullabilityModifiers}.
type CheckNullabilityModifiersTest struct {
	CompilerTestCaseBase
}

func NewCheckNullabilityModifiersTest() *CheckNullabilityModifiersTest {
	c := &CheckNullabilityModifiersTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestCheckNullabilityModifiersTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCheckNullabilityModifiersTest() })
}

// Java: CheckNullabilityModifiersTest.getProcessor(Compiler)
func (c *CheckNullabilityModifiersTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewCheckNullabilityModifiers(compiler)
}

// Java: CheckNullabilityModifiersTest.getOptions()
func (c *CheckNullabilityModifiersTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetWarningLevel(DiagnosticGroups_LINT_CHECKS, CheckLevel_WARNING)
	return options
}

// Java: CheckNullabilityModifiersTest.testPrimitiveType()
func (c *CheckNullabilityModifiersTest) TestPrimitiveType() {
	c.CheckRedundantWarning("/** @type {!boolean} */ var x;")
	c.CheckRedundantWarning("/** @type {!number} */ var x;")
	c.CheckRedundantWarning("/** @type {!bigint} */ var x;")
	c.CheckRedundantWarning("/** @type {!string} */ var x;")
	c.CheckRedundantWarning("/** @type {!symbol} */ var x;")
	c.CheckRedundantWarning("/** @type {!undefined} */ var x;")
	c.CheckRedundantWarning("/** @type {!void} */ var x;")
	c.CheckNoWarning("/** @type {?boolean} */ var x;")
	c.CheckNoWarning("/** @type {boolean} */ var x;")
}

// Java: CheckNullabilityModifiersTest.testReferenceType()
func (c *CheckNullabilityModifiersTest) TestReferenceType() {
	c.CheckMissingWarning("/** @type {Object} */ var x;")
	c.CheckMissingWarning("/** @type {Function} */ var x;")
	c.CheckMissingWarning("/** @type {Symbol} */ var x;")
	c.CheckNoWarning("/** @type {?Object} */ var x;")
	c.CheckNoWarning("/** @type {!Object} */ var x;")
}

// Java: CheckNullabilityModifiersTest.testRecordType()
func (c *CheckNullabilityModifiersTest) TestRecordType() {
	c.CheckRedundantWarning("/** @type {!{foo: string}} */ var x;")
	c.CheckRedundantWarning("/** @type {{foo: !string}} */ var x;")
	c.CheckMissingWarning("/** @type {{foo: Object}} */ var x;")
	c.CheckNoWarning("/** @type {?{foo: string}} */ var x;")
	c.CheckNoWarning("/** @type {{foo: string}} */ var x;")
	c.CheckNoWarning("/** @type {{foo: ?string}} */ var x;")
	c.CheckNoWarning("/** @type {{foo: !Object}} */ var x;")
	c.CheckNoWarning("/** @type {{foo: ?Object}} */ var x;")
}

// Java: CheckNullabilityModifiersTest.testFunctionType()
func (c *CheckNullabilityModifiersTest) TestFunctionType() {
	c.CheckRedundantWarning("/** @type {!function()} */ function f(){}")
	c.CheckRedundantWarning("/** @type {function(!string)} */ function f(x){}")
	c.CheckRedundantWarning("/** @type {function(): !string} */ function f(x){}")
	c.CheckMissingWarning("/** @type {function(Object)} */ function f(x){}")
	c.CheckMissingWarning("/** @type {function(): Object} */ function f(x){}")
	c.CheckNoWarning("/** @type {function()} */ function f(){}")
	c.CheckNoWarning("/** @type {function(?string): ?string} */ function f(x){}")
	c.CheckNoWarning("/** @type {function(string): string} */ function f(x){}")
	c.CheckNoWarning("/** @type {function(!Object): ?Object} */ function f(x){}")
	c.CheckNoWarning("/** @type {function(new:Object)} */ function f(){}")
	c.CheckNoWarning("/** @type {function(this:Object)} */ function f(){}")
}

// Java: CheckNullabilityModifiersTest.testUnionType()
func (c *CheckNullabilityModifiersTest) TestUnionType() {
	c.CheckRedundantWarning("/** @type {!Object|!string} */ var x;")
	c.CheckMissingWarning("/** @type {Object|string} */ var x;")
	c.CheckNoWarning("/** @type {?Object|string} */ var x;")
	c.CheckNoWarning("/** @type {!Object|string} */ var x;")
}

// Java: CheckNullabilityModifiersTest.testEnumType()
func (c *CheckNullabilityModifiersTest) TestEnumType() {
	c.CheckRedundantWarning("/** @enum {!boolean} */ var x;")
	c.CheckRedundantWarning("/** @enum {!number} */ var x;")
	c.CheckRedundantWarning("/** @enum {!bigint} */ var x;")
	c.CheckRedundantWarning("/** @enum {!string} */ var x;")
	c.CheckRedundantWarning("/** @enum {!symbol} */ var x;")
	c.CheckNoWarning("/** @enum {?string} */ var x;")
	c.CheckNoWarning("/** @enum {string} */ var x;")
	c.CheckNoWarning("/** @enum {Object} */ var x;")
	c.CheckNoWarning("/** @enum {?Object} */ var x;")
	c.CheckNoWarning("/** @enum {!Object} */ var x;")
}

// Java: CheckNullabilityModifiersTest.testTemplateDefinitionType()
func (c *CheckNullabilityModifiersTest) TestTemplateDefinitionType() {
	c.CheckNoWarning("/** @param {T} x @template T */", "function f(x){}")
	c.CheckNoWarning("/** @param {S} x @return {T} @template S,T */", "function f(x){}")
	c.CheckNoWarning(
		`/** @constructor @template T */ function Foo(){}
/** @param {T} x */ Foo.prototype.bar = function(x){};
`)
}

// Java: CheckNullabilityModifiersTest.testTemplateDefinitionTypeWithTtl()
func (c *CheckNullabilityModifiersTest) TestTemplateDefinitionTypeWithTtl() {
	c.CheckNoWarning(
		`/**
 * @param {T} x
 * @param {U} y
 * @template T
 * @template U := T =:
 */
function f(x, y) {}
`)
	c.CheckNoWarning(
		`/**
 * @param {T} x
 * @param {U} y
 * @template T
// note that "Object" below is /not/ nullable
 * @template U := cond(isUnknown(T), "Object", T) =:
 */
function f(x, y) {}
`)
}

// Java: CheckNullabilityModifiersTest.testTemplateInstantiationType()
func (c *CheckNullabilityModifiersTest) TestTemplateInstantiationType() {
	c.CheckRedundantWarning("/** @type {!Array<!string>} */ var x;")
	c.CheckMissingWarning("/** @type {Array<string>} */ var x;")
	c.CheckMissingWarning("/** @type {!Array<Object>} */ var x;")
	c.CheckNoWarning("/** @type {!Array<?string>} */ var x;")
	c.CheckNoWarning("/** @type {!Array<string>} */ var x;")
	c.CheckNoWarning("/** @type {!Array<?Object>} */ var x;")
	c.CheckNoWarning("/** @type {!Array<!Object>} */ var x;")
}

// Java: CheckNullabilityModifiersTest.testParamType()
func (c *CheckNullabilityModifiersTest) TestParamType() {
	c.CheckRedundantWarning("/** @param {!string} x */ function f(x){}")
	c.CheckMissingWarning("/** @param {Object} x */ function f(x){}")
	c.CheckNoWarning("/** @param {?string} x */ function f(x){}")
	c.CheckNoWarning("/** @param {string} x */ function f(x){}")
	c.CheckNoWarning("/** @param {?Object} x */ function f(x){}")
	c.CheckNoWarning("/** @param {!Object} x */ function f(x){}")
}

// Java: CheckNullabilityModifiersTest.testParamMissingType()
func (c *CheckNullabilityModifiersTest) TestParamMissingType() {
	c.CheckNoWarning("/** @param x */ function f(x){}")
}

// Java: CheckNullabilityModifiersTest.testReturnType()
func (c *CheckNullabilityModifiersTest) TestReturnType() {
	c.CheckRedundantWarning("/** @return {!string} */ function f(){}")
	c.CheckMissingWarning("/** @return {Object} */ function f(){}")
	c.CheckNoWarning("/** @return {?string} */ function f(){}")
	c.CheckNoWarning("/** @return {string} */ function f(){}")
	c.CheckNoWarning("/** @return {?Object} */ function f(){}")
	c.CheckNoWarning("/** @return {!Object} */ function f(){}")
}

// Java: CheckNullabilityModifiersTest.testTypedefType()
func (c *CheckNullabilityModifiersTest) TestTypedefType() {
	c.CheckRedundantWarning("/** @typedef {!string} */ var x;")
	c.CheckMissingWarning("/** @typedef {Object} */ var x;")
	c.CheckNoWarning("/** @typedef {?string} */ var x;")
	c.CheckNoWarning("/** @typedef {string} */ var x;")
	c.CheckNoWarning("/** @typedef {?Object} */ var x;")
	c.CheckNoWarning("/** @typedef {!Object} */ var x;")
}

// Java: CheckNullabilityModifiersTest.testThisType()
func (c *CheckNullabilityModifiersTest) TestThisType() {
	c.CheckRedundantWarning("/** @this {!string} */ function f(){}")
	c.CheckNoWarning("/** @this {?string} */ function f(){}")
	c.CheckNoWarning("/** @this {string} */ function f(){}")
	c.CheckNoWarning("/** @this {Object} */ function f(){}")
	c.CheckNoWarning("/** @this {?Object} */ function f(){}")
	c.CheckNoWarning("/** @this {!Object} */ function f(){}")
}

// Java: CheckNullabilityModifiersTest.testBaseType()
func (c *CheckNullabilityModifiersTest) TestBaseType() {
	c.CheckNoWarning("/** @extends {Object} */ function f(){}")
	c.CheckNoWarning("/** @implements {Object} */ function f(){}")
}

// Java: CheckNullabilityModifiersTest.testTypeOf()
func (c *CheckNullabilityModifiersTest) TestTypeOf() {
	c.CheckNoWarning("/** @type {typeof Object} */ var x;")
}

// Java: CheckNullabilityModifiersTest.testEndPosition()
func (c *CheckNullabilityModifiersTest) TestEndPosition() {
	c.CheckRedundantWarning("/** @type {string!} */ var x;")
	c.CheckNoWarning("/** @type {string?} */ var x;")
	c.CheckNoWarning("/** @type {Object!} */ var x;")
	c.CheckNoWarning("/** @type {Object?} */ var x;")
}

// Java: CheckNullabilityModifiersTest.testMultipleFiles()
func (c *CheckNullabilityModifiersTest) TestMultipleFiles() {
	c.CheckMissingWarning(
		"/** @param {T} x @return {T} @template T */ function f(x){}",
		"/** @param {T} x */ function g(x){}")
}

// Java: CheckNullabilityModifiersTest.testSetToNull()
func (c *CheckNullabilityModifiersTest) TestSetToNull() {
	c.CheckNullMissingWarning("/** @type {Object} */ var x = null;")
	c.CheckNullMissingWarning(
		"/** @constructor */ function C() {} /** @private {Object} */ C.prop = null;")
	c.CheckNullMissingWarning(
		"/** @constructor */ function C() { /** @private {Object} */ this.foo = null; }")
	c.CheckNullMissingWarning(
		"/** @constructor */ function C() { /** @private {Object|string} */ this.foo = null; }")
	c.CheckNullMissingWarning(
		"/** @constructor */ function C() { /** @private {!String|Symbol} */ this.foo = null; }")
	c.CheckNoWarning("/** @type {?Object} */ var x = null;")
	c.CheckNoWarning("/** @constructor */ function C() {} /** @private {?Symbol} */ C.prop = null;")
	c.CheckNoWarning(
		"/** @constructor */ function C() { /** @private {?Object} */ this.foo = null; }")
	// don't recommend making 'Type' nullable since it's not the root of the type expression
	c.CheckMissingWarning("/** @type {!Array<Type>} */ let arr = null;")
	c.CheckMissingWarning("/** @type {?{prop: Type}} */ let o = null;")
}

// Java: CheckNullabilityModifiersTest.checkNoWarning(String...)
func (c *CheckNullabilityModifiersTest) CheckNoWarning(js ...string) {
	c.TestSame(c.Srcs(js))
}

// Java: CheckNullabilityModifiersTest.checkMissingWarning(String...)
func (c *CheckNullabilityModifiersTest) CheckMissingWarning(js ...string) {
	c.TestWarning(c.Srcs(js), CheckNullabilityModifiers_MISSING_NULLABILITY_MODIFIER_JSDOC)
}

// Java: CheckNullabilityModifiersTest.checkNullMissingWarning(String...)
func (c *CheckNullabilityModifiersTest) CheckNullMissingWarning(js ...string) {
	c.TestWarning(c.Srcs(js), CheckNullabilityModifiers_NULL_MISSING_NULLABILITY_MODIFIER_JSDOC)
}

// Java: CheckNullabilityModifiersTest.checkRedundantWarning(String...)
func (c *CheckNullabilityModifiersTest) CheckRedundantWarning(js ...string) {
	c.TestWarning(c.Srcs(js), CheckNullabilityModifiers_REDUNDANT_NULLABILITY_MODIFIER_JSDOC)
}
