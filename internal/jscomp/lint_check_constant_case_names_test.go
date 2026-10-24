// Port of com/google/javascript/jscomp/lint/CheckConstantCaseNamesTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/lint/CheckConstantCaseNamesTest.java)

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

type CheckConstantCaseNamesTest struct {
	CompilerTestCaseBase
}

func NewCheckConstantCaseNamesTest() *CheckConstantCaseNamesTest {
	c := &CheckConstantCaseNamesTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestCheckConstantCaseNamesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCheckConstantCaseNamesTest() })
}

// Java: CheckConstantCaseNamesTest.getProcessor(Compiler)
func (c *CheckConstantCaseNamesTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewCheckConstantCaseNames(compiler)
}

// Java: CheckConstantCaseNamesTest.getOptions()
func (c *CheckConstantCaseNamesTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetWarningLevel(DiagnosticGroups_LINT_CHECKS, CheckLevel_WARNING)
	return options
}

// Java: CheckConstantCaseNamesTest.emptyScriptIsOk()
//
// Go-only: test-helper — the Test prefix makes RunJUnit4 run this @Test method.
func (c *CheckConstantCaseNamesTest) TestEmptyScriptIsOk() {
	c.TestNoWarning("")
}

// Java: CheckConstantCaseNamesTest.emptyModuleIsOk()
//
// Go-only: test-helper — the Test prefix makes RunJUnit4 run this @Test method.
func (c *CheckConstantCaseNamesTest) TestEmptyModuleIsOk() {
	c.TestNoWarning("goog.module('mod');")
	c.TestNoWarning("export {};")
}

// Java: CheckConstantCaseNamesTest.constantCaseNameDeclarations_requireExplicitConst()
//
// Go-only: test-helper — the Test prefix makes RunJUnit4 run this @Test method.
func (c *CheckConstantCaseNamesTest) TestConstantCaseNameDeclarations_requireExplicitConst() {
	c.TestWarning("goog.module('m'); var ABC;", CheckConstantCaseNames_MISSING_CONST_PROPERTY)
	c.TestWarning("goog.module('m'); var ABC = 0;", CheckConstantCaseNames_MISSING_CONST_PROPERTY)
	c.TestWarning("goog.module('m'); let ABC = 0;", CheckConstantCaseNames_MISSING_CONST_PROPERTY)
	c.TestWarning("goog.module('m'); var ABC_DEF_0 = 0;", CheckConstantCaseNames_MISSING_CONST_PROPERTY)
	c.TestWarning("goog.module('m'); var abc = 0, DEF = 1;", CheckConstantCaseNames_MISSING_CONST_PROPERTY)
	c.Test(
		c.Srcs("goog.module('m'); var ABC = 0, DEF = 1;"),
		c.Warning(CheckConstantCaseNames_MISSING_CONST_PROPERTY),
		c.Warning(CheckConstantCaseNames_MISSING_CONST_PROPERTY))
	c.TestWarning("goog.module('m'); var {ABC} = obj;", CheckConstantCaseNames_MISSING_CONST_PROPERTY)
	c.TestWarning("goog.module('m'); var {abc, DEF} = obj;", CheckConstantCaseNames_MISSING_CONST_PROPERTY)
	c.TestWarning("let ABC = 0; export {ABC};", CheckConstantCaseNames_MISSING_CONST_PROPERTY)
}

// Java: CheckConstantCaseNamesTest.constantCaseNameDeclarations_whenReassigned_requireCamelCase()
//
// Go-only: test-helper — the Test prefix makes RunJUnit4 run this @Test method.
func (c *CheckConstantCaseNamesTest) TestConstantCaseNameDeclarations_whenReassigned_requireCamelCase() {
	c.TestWarning("goog.module('m'); var ABC = 0; ABC = 1;", CheckConstantCaseNames_REASSIGNED_CONSTANT_CASE_NAME)
	c.TestWarning("goog.module('m'); let ABC = 0; ABC = 1;", CheckConstantCaseNames_REASSIGNED_CONSTANT_CASE_NAME)
	c.TestWarning("goog.module('m'); let ABC = 0; [[ABC]] = [[1]];", CheckConstantCaseNames_REASSIGNED_CONSTANT_CASE_NAME)
}

// Java: CheckConstantCaseNamesTest.constantCaseNameDeclarations_whenReassignedInFunction_requireCamelCase()
//
// Go-only: test-helper — the Test prefix makes RunJUnit4 run this @Test method.
func (c *CheckConstantCaseNamesTest) TestConstantCaseNameDeclarations_whenReassignedInFunction_requireCamelCase() {
	c.TestWarning(
		"goog.module('m'); var ABC = 0; exports = function() { ABC = 1; };",
		CheckConstantCaseNames_REASSIGNED_CONSTANT_CASE_NAME)
	c.TestWarning(
		"goog.module('m'); let ABC = 0; exports = function() { ABC = 1; };",
		CheckConstantCaseNames_REASSIGNED_CONSTANT_CASE_NAME)
}

// Java: CheckConstantCaseNamesTest.constantCaseNameDeclarations_localShadowAreNotReassignments()
//
// Go-only: test-helper — the Test prefix makes RunJUnit4 run this @Test method.
func (c *CheckConstantCaseNamesTest) TestConstantCaseNameDeclarations_localShadowAreNotReassignments() {
	c.TestWarning(
		// the local 'ABC' does not refer to the module-level 'ABC'.
		"goog.module('m'); let ABC = 0; if (cond) { let ABC = 0; ABC = 1; }",
		CheckConstantCaseNames_MISSING_CONST_PROPERTY)
	c.TestWarning(
		// the local 'ABC' does not refer to the module-level 'ABC'.
		"goog.module('m'); let ABC = 0; exports = function(ABC) { ABC = 1; };",
		CheckConstantCaseNames_MISSING_CONST_PROPERTY)
}

// Java: CheckConstantCaseNamesTest.constantCaseNameDeclarations_usagesAreNotReassignments()
//
// Go-only: test-helper — the Test prefix makes RunJUnit4 run this @Test method.
func (c *CheckConstantCaseNamesTest) TestConstantCaseNameDeclarations_usagesAreNotReassignments() {
	c.TestWarning("goog.module('m'); let ABC = 0; alert(ABC);", CheckConstantCaseNames_MISSING_CONST_PROPERTY)
	c.TestWarning("goog.module('m'); let ABC = 0; const {d = ABC} = {};", CheckConstantCaseNames_MISSING_CONST_PROPERTY)
	c.TestWarning(
		"goog.module('m'); let ABC = 0; let d = 1; export {d as ABC};", CheckConstantCaseNames_MISSING_CONST_PROPERTY)
}

// Java: CheckConstantCaseNamesTest.constantCaseNameDeclarations_appearingInMultipleFiles_getDistinctErrors()
//
// Go-only: test-helper — the Test prefix makes RunJUnit4 run this @Test method.
func (c *CheckConstantCaseNamesTest) TestConstantCaseNameDeclarations_appearingInMultipleFiles_getDistinctErrors() {
	c.TestWarning(
		c.Srcs(
			"goog.module('m'); let ABC = 0;",
			"let ABC = 0; ABC = 1;", // not a reassignment of the module local.
		),
		c.Warning(CheckConstantCaseNames_MISSING_CONST_PROPERTY))
	c.Test(
		c.Srcs(
			// Recognize that only one out of the two 'ABC's is mutated.
			"goog.module('m1'); let ABC = 0;", "goog.module('m2'); let ABC = 0; ABC = 1;"),
		c.Warning(CheckConstantCaseNames_MISSING_CONST_PROPERTY),
		c.Warning(CheckConstantCaseNames_REASSIGNED_CONSTANT_CASE_NAME))
}

// Java: CheckConstantCaseNamesTest.constantCaseNameDeclarations_okIfAlreadyConst()
//
// Go-only: test-helper — the Test prefix makes RunJUnit4 run this @Test method.
func (c *CheckConstantCaseNamesTest) TestConstantCaseNameDeclarations_okIfAlreadyConst() {
	c.TestSame("goog.module('m'); const ABC = 0;")
	c.TestSame("goog.module('m'); const ABC = 0;")
	c.TestSame("goog.module('m'); const ABC_DEF_0 = 0;")
	c.TestSame("goog.module('m'); /** @const */ var ABC = 0;")
	c.TestSame("goog.module('m'); /** @const */ let ABC = 0;")
	c.TestSame("goog.module('m'); /** @const */ var ABC_DEF_0 = 0;")
	c.TestSame("goog.module('m'); const {ABC} = obj;")
}

// Java: CheckConstantCaseNamesTest.nonConstantCaseNameDeclarations_dontRequireConst()
//
// Go-only: test-helper — the Test prefix makes RunJUnit4 run this @Test method.
func (c *CheckConstantCaseNamesTest) TestNonConstantCaseNameDeclarations_dontRequireConst() {
	c.TestSame("goog.module('m'); var abc = 0;")
	c.TestSame("goog.module('m'); let abc = 0;")
	c.TestSame("goog.module('m'); var AbcDef = 0;")
	c.TestSame("goog.module('m'); var abc;")
}

// Java: CheckConstantCaseNamesTest.nonModuleLevelNamesDontCauseWarning()
//
// Go-only: test-helper — the Test prefix makes RunJUnit4 run this @Test method.
func (c *CheckConstantCaseNamesTest) TestNonModuleLevelNamesDontCauseWarning() {
	c.TestSame("var ABC = 0;")
	c.TestSame("goog.module('m'); function fn() { var ABC = 0; }")
}

// Java: CheckConstantCaseNamesTest.constantCasePropertyDeclarations_dontRequireExplicitConst()
//
// Go-only: test-helper — the Test prefix makes RunJUnit4 run this @Test method.
func (c *CheckConstantCaseNamesTest) TestConstantCasePropertyDeclarations_dontRequireExplicitConst() {
	c.TestSame("goog.module('m'); /** @private */ a.NAME = 'name';")
}

// Java: CheckConstantCaseNamesTest.constantCaseObjectLiteralKeys_dontRequireExplicitConst()
//
// Go-only: test-helper — the Test prefix makes RunJUnit4 run this @Test method.
func (c *CheckConstantCaseNamesTest) TestConstantCaseObjectLiteralKeys_dontRequireExplicitConst() {
	c.TestSame("goog.module('m'); const Colors = {RED: 0, YELLOW: 1, GREEN: 2}")
	c.TestSame("goog.module('m'); ns.Colors = {RED: 0, YELLOW: 1, GREEN: 2}")
}
