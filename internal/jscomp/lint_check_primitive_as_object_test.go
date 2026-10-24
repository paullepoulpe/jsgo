// Port of com/google/javascript/jscomp/lint/CheckPrimitiveAsObjectTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/lint/CheckPrimitiveAsObjectTest.java)

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

type CheckPrimitiveAsObjectTest struct {
	CompilerTestCaseBase
}

func NewCheckPrimitiveAsObjectTest() *CheckPrimitiveAsObjectTest {
	c := &CheckPrimitiveAsObjectTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestCheckPrimitiveAsObjectTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCheckPrimitiveAsObjectTest() })
}

// Java: CheckPrimitiveAsObjectTest.getProcessor(Compiler)
func (c *CheckPrimitiveAsObjectTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewCheckPrimitiveAsObject(compiler)
}

// Java: CheckPrimitiveAsObjectTest.testWarningForBooleanObjectCreation()
func (c *CheckPrimitiveAsObjectTest) TestWarningForBooleanObjectCreation() {
	c.TestWarning("new Boolean(false)", CheckPrimitiveAsObject_NEW_PRIMITIVE_OBJECT)
}

// Java: CheckPrimitiveAsObjectTest.testWarningForNumberObjectCreation()
func (c *CheckPrimitiveAsObjectTest) TestWarningForNumberObjectCreation() {
	c.TestWarning("new Number(5)", CheckPrimitiveAsObject_NEW_PRIMITIVE_OBJECT)
}

// Java: CheckPrimitiveAsObjectTest.testWarningForStringObjectCreation()
func (c *CheckPrimitiveAsObjectTest) TestWarningForStringObjectCreation() {
	c.TestWarning("new String(\"hello\")", CheckPrimitiveAsObject_NEW_PRIMITIVE_OBJECT)
}

// Java: CheckPrimitiveAsObjectTest.testNoWarningForObjectCreation()
func (c *CheckPrimitiveAsObjectTest) TestNoWarningForObjectCreation() {
	c.TestSame("new Object()")
}

// Java: CheckPrimitiveAsObjectTest.testNoWarningForQualifiedClassCreation()
func (c *CheckPrimitiveAsObjectTest) TestNoWarningForQualifiedClassCreation() {
	c.TestSame("new my.qualified.ClassName()")
}

// Java: CheckPrimitiveAsObjectTest.testWarningForBooleanTypeDeclaration()
func (c *CheckPrimitiveAsObjectTest) TestWarningForBooleanTypeDeclaration() {
	c.TestWarning("/** @type {Boolean} */ var x;", CheckPrimitiveAsObject_PRIMITIVE_OBJECT_DECLARATION)
}

// Java: CheckPrimitiveAsObjectTest.testWarningForBooleanTypeDeclaration_withES6Modules()
func (c *CheckPrimitiveAsObjectTest) TestWarningForBooleanTypeDeclaration_withES6Modules() {
	c.TestWarning("export /** @type {Boolean} */ var x;", CheckPrimitiveAsObject_PRIMITIVE_OBJECT_DECLARATION)
}

// Java: CheckPrimitiveAsObjectTest.testWarningForNumberTypeDeclaration()
func (c *CheckPrimitiveAsObjectTest) TestWarningForNumberTypeDeclaration() {
	c.TestWarning("/** @type {Number} */ var x;", CheckPrimitiveAsObject_PRIMITIVE_OBJECT_DECLARATION)
}

// Java: CheckPrimitiveAsObjectTest.testWarningForNumberTypeDeclaration_withES6Modules()
func (c *CheckPrimitiveAsObjectTest) TestWarningForNumberTypeDeclaration_withES6Modules() {
	c.TestWarning("export /** @type {Number} */ var x;", CheckPrimitiveAsObject_PRIMITIVE_OBJECT_DECLARATION)
}

// Java: CheckPrimitiveAsObjectTest.testWarningForStringTypeDeclaration()
func (c *CheckPrimitiveAsObjectTest) TestWarningForStringTypeDeclaration() {
	c.TestWarning("/** @type {String} */ var x;", CheckPrimitiveAsObject_PRIMITIVE_OBJECT_DECLARATION)
}

// Java: CheckPrimitiveAsObjectTest.testWarningForStringTypeDeclaration_withES6Modules()
func (c *CheckPrimitiveAsObjectTest) TestWarningForStringTypeDeclaration_withES6Modules() {
	c.TestWarning("export /** @type {String} */ var x;", CheckPrimitiveAsObject_PRIMITIVE_OBJECT_DECLARATION)
}

// Java: CheckPrimitiveAsObjectTest.testWarningForBooleanInsideTypeDeclaration()
func (c *CheckPrimitiveAsObjectTest) TestWarningForBooleanInsideTypeDeclaration() {
	c.TestWarning("/** @type {function(): Boolean} */ var x;", CheckPrimitiveAsObject_PRIMITIVE_OBJECT_DECLARATION)
	c.TestWarning("/** @type {function(Boolean)} */ var x;", CheckPrimitiveAsObject_PRIMITIVE_OBJECT_DECLARATION)
	c.TestWarning("/** @type {{b: Boolean}} */ var x;", CheckPrimitiveAsObject_PRIMITIVE_OBJECT_DECLARATION)
}

// Java: CheckPrimitiveAsObjectTest.testWarningForBooleanInsideTypeDeclaration_withES6Modules()
func (c *CheckPrimitiveAsObjectTest) TestWarningForBooleanInsideTypeDeclaration_withES6Modules() {
	c.TestWarning("export /** @type {function(): Boolean} */ var x;", CheckPrimitiveAsObject_PRIMITIVE_OBJECT_DECLARATION)
}

// Java: CheckPrimitiveAsObjectTest.testWarningForNumberParameterDeclaration()
func (c *CheckPrimitiveAsObjectTest) TestWarningForNumberParameterDeclaration() {
	c.TestWarning(
		`/**
 * @param {Number=} x
 * @return {number}
 */
function f(x) {
  return x + 1;
}
`,
		CheckPrimitiveAsObject_PRIMITIVE_OBJECT_DECLARATION)
}

// Java: CheckPrimitiveAsObjectTest.testWarningForNumberParameterDeclaration_withES6Modules()
func (c *CheckPrimitiveAsObjectTest) TestWarningForNumberParameterDeclaration_withES6Modules() {
	c.TestWarning(
		`export
/**
 * @param {Number=} x
 * @return {number}
 */
function f(x) {
  return x + 1;
}
`,
		CheckPrimitiveAsObject_PRIMITIVE_OBJECT_DECLARATION)
}

// Java: CheckPrimitiveAsObjectTest.testWarningForBooleanParameterDeclarationInTypedef()
func (c *CheckPrimitiveAsObjectTest) TestWarningForBooleanParameterDeclarationInTypedef() {
	c.TestWarning(
		`/**
 * @typedef {function(Boolean=)}
 */
var takesOptionalBoolean;
`,
		CheckPrimitiveAsObject_PRIMITIVE_OBJECT_DECLARATION)
}

// Java: CheckPrimitiveAsObjectTest.testWarningForBooleanParameterDeclarationInTypedef_withES6Modules()
func (c *CheckPrimitiveAsObjectTest) TestWarningForBooleanParameterDeclarationInTypedef_withES6Modules() {
	c.TestWarning(
		`export
/**
 * @typedef {function(Boolean=)}
 */
var takesOptionalBoolean;
`,
		CheckPrimitiveAsObject_PRIMITIVE_OBJECT_DECLARATION)
}

// Java: CheckPrimitiveAsObjectTest.testWarningForNumberReturnDeclaration()
func (c *CheckPrimitiveAsObjectTest) TestWarningForNumberReturnDeclaration() {
	c.TestWarning(
		`/**
 * @return {Number}
 */
function f() {
  return 5;
}
`,
		CheckPrimitiveAsObject_PRIMITIVE_OBJECT_DECLARATION)
}

// Java: CheckPrimitiveAsObjectTest.testWarningForNumberReturnDeclaration_withES6Modules()
func (c *CheckPrimitiveAsObjectTest) TestWarningForNumberReturnDeclaration_withES6Modules() {
	c.TestWarning(
		`export
/**
 * @return {Number}
 */
function f() {
  return 5;
}
`,
		CheckPrimitiveAsObject_PRIMITIVE_OBJECT_DECLARATION)
}
