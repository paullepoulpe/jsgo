// Port of com/google/javascript/jscomp/CheckTypeImportCodeReferencesTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/CheckTypeImportCodeReferencesTest.java)

/*
 * Copyright 2019 The Closure Compiler Authors.
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

// CheckTypeImportCodeReferencesTest: Tests for CheckTypeImportCodeReferences.
type CheckTypeImportCodeReferencesTest struct {
	CompilerTestCaseBase
}

// Java: CheckTypeImportCodeReferencesTest()
func NewCheckTypeImportCodeReferencesTest() *CheckTypeImportCodeReferencesTest {
	c := &CheckTypeImportCodeReferencesTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestCheckTypeImportCodeReferencesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCheckTypeImportCodeReferencesTest() })
}

// Java: CheckTypeImportCodeReferencesTest.getProcessor(Compiler)
func (c *CheckTypeImportCodeReferencesTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewCheckTypeImportCodeReferences(compiler)
}

// Java: CheckTypeImportCodeReferencesTest.testRequireTypeCodeReference_namespace()
func (c *CheckTypeImportCodeReferencesTest) TestRequireTypeCodeReference_namespace() {
	c.TestSame(
		c.Srcs(
			`goog.module('test');
const alias = goog.requireType('namespace');
/** @type {namespace} */ let x;
namespace;
`))
}

// Java: CheckTypeImportCodeReferencesTest.testRequireTypeCodeReference_alias()
func (c *CheckTypeImportCodeReferencesTest) TestRequireTypeCodeReference_alias() {
	c.TestSame(
		c.Srcs(
			`goog.module('test');
const alias = goog.requireType('namespace');
/** @type {alias} */ let x;
`))

	c.Test(
		c.Srcs(
			`goog.module('test');
const alias = goog.requireType('namespace');
alias;
`),
		c.Error(CheckTypeImportCodeReferences_TYPE_IMPORT_CODE_REFERENCE))

	c.Test(
		c.Srcs(
			`goog.module('test');
let alias = goog.requireType('namespace');
alias;
`),
		c.Error(CheckTypeImportCodeReferences_TYPE_IMPORT_CODE_REFERENCE))

	c.Test(
		c.Srcs(
			`goog.module('test');
var alias = goog.requireType('namespace');
alias;
`),
		c.Error(CheckTypeImportCodeReferences_TYPE_IMPORT_CODE_REFERENCE))
}

// Java: CheckTypeImportCodeReferencesTest.testRequireTypeCodeReference_destructuring()
func (c *CheckTypeImportCodeReferencesTest) TestRequireTypeCodeReference_destructuring() {
	c.TestSame(
		c.Srcs(
			`goog.module('test');
const {alias} = goog.requireType('namespace');
/** @type {alias} */ let x;
`))

	c.Test(
		c.Srcs(
			`goog.module('test');
const {alias} = goog.requireType('namespace');
alias;
`),
		c.Error(CheckTypeImportCodeReferences_TYPE_IMPORT_CODE_REFERENCE))
}

// Java: CheckTypeImportCodeReferencesTest.testRequireTypeCodeReference_destructuring_shortname()
func (c *CheckTypeImportCodeReferencesTest) TestRequireTypeCodeReference_destructuring_shortname() {
	c.TestSame(
		c.Srcs(
			`goog.module('test');
const {alias: myAlias} = goog.requireType('namespace');
alias;
`))

	c.TestSame(
		c.Srcs(
			`goog.module('test');
const {alias: myAlias} = goog.requireType('namespace');
/** @type {myAlias} */ let x;
`))

	c.Test(
		c.Srcs(
			`goog.module('test');
const {alias: myAlias} = goog.requireType('namespace');
myAlias;
`),
		c.Error(CheckTypeImportCodeReferences_TYPE_IMPORT_CODE_REFERENCE))
}

// Java: CheckTypeImportCodeReferencesTest.testRequireTypeCodeReference_destructuring_multiple()
func (c *CheckTypeImportCodeReferencesTest) TestRequireTypeCodeReference_destructuring_multiple() {
	c.Test(
		c.Srcs(
			`goog.module('test');
const {otherAlias, alias: myAlias, anotherAlias} = goog.requireType('namespace');
myAlias;
`),
		c.Error(CheckTypeImportCodeReferences_TYPE_IMPORT_CODE_REFERENCE))
}

// Java: CheckTypeImportCodeReferencesTest.testRequireTypeCodeReference_shadowing()
func (c *CheckTypeImportCodeReferencesTest) TestRequireTypeCodeReference_shadowing() {
	c.TestSame(
		c.Srcs(
			`goog.module('test');
const alias = goog.requireType('namespace');
function foo(alias) { alias; }
`))

	c.TestSame(
		c.Srcs(
			`goog.module('test');
const alias = goog.requireType('namespace');
function foo() { let alias; }
`))

	c.Test(
		c.Srcs(
			`goog.module('test');
const alias = goog.requireType('namespace');
function foo() { alias; }
`),
		c.Error(CheckTypeImportCodeReferences_TYPE_IMPORT_CODE_REFERENCE))
}

// Java: CheckTypeImportCodeReferencesTest.testRequireTypeCodeReference_getprop_left()
func (c *CheckTypeImportCodeReferencesTest) TestRequireTypeCodeReference_getprop_left() {
	c.TestSame(
		c.Srcs(
			`goog.module('test');
const alias = goog.requireType('namespace');
/** @type {alias.prop} */ let x;
`))

	c.Test(
		c.Srcs(
			`goog.module('test');
const alias = goog.requireType('namespace');
alias.prop;
`),
		c.Error(CheckTypeImportCodeReferences_TYPE_IMPORT_CODE_REFERENCE))
}

// Java: CheckTypeImportCodeReferencesTest.testRequireTypeCodeReference_getprop_right()
func (c *CheckTypeImportCodeReferencesTest) TestRequireTypeCodeReference_getprop_right() {
	c.TestSame(
		c.Srcs(
			`goog.module('test');
const alias = goog.requireType('namespace');
obj.alias;
`))
}

// Java: CheckTypeImportCodeReferencesTest.testRequireTypeCodeReference_implicit_var()
func (c *CheckTypeImportCodeReferencesTest) TestRequireTypeCodeReference_implicit_var() {
	c.TestSame(
		c.Srcs(
			`goog.module('test');
const alias = goog.requireType('namespace');
function f() { arguments; }
`))
}
