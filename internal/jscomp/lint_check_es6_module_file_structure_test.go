// Port of com/google/javascript/jscomp/lint/CheckEs6ModuleFileStructureTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/lint/CheckEs6ModuleFileStructureTest.java)

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

type CheckEs6ModuleFileStructureTest struct {
	CompilerTestCaseBase
}

func NewCheckEs6ModuleFileStructureTest() *CheckEs6ModuleFileStructureTest {
	c := &CheckEs6ModuleFileStructureTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestCheckEs6ModuleFileStructureTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCheckEs6ModuleFileStructureTest() })
}

// Java: CheckEs6ModuleFileStructureTest.setUp()
func (c *CheckEs6ModuleFileStructureTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.IgnoreWarnings(ModuleLoader_INVALID_MODULE_PATH)
}

// Java: CheckEs6ModuleFileStructureTest.getProcessor(Compiler)
func (c *CheckEs6ModuleFileStructureTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewCheckEs6ModuleFileStructure(compiler)
}

// Java: CheckEs6ModuleFileStructureTest.testImportsMustBeFirstIfPresent()
func (c *CheckEs6ModuleFileStructureTest) TestImportsMustBeFirstIfPresent() {
	c.TestSame("var notAModule;")
	c.TestSame("export let noImports;")
	c.TestSame("import importsFirst from 'file'; let notFirst;")
	c.TestWarning("let first; import notFirst from 'file';", CheckEs6ModuleFileStructure_MUST_COME_BEFORE)
}

// Java: CheckEs6ModuleFileStructureTest.testDeclareModuleIdAfterImportsIfPresent()
func (c *CheckEs6ModuleFileStructureTest) TestDeclareModuleIdAfterImportsIfPresent() {
	c.TestSame("import noDeclareNamespace from 'file';")
	c.TestSame("goog.declareModuleId('name.space'); export let noImports;")
	c.TestSame("import first from 'file'; goog.declareModuleId('name.space');")
	c.TestWarning("export let first; goog.declareModuleId('name.space');", CheckEs6ModuleFileStructure_MUST_COME_BEFORE)
	c.TestWarning(
		"import first from 'file'; let second; goog.declareModuleId('name.space');",
		CheckEs6ModuleFileStructure_MUST_COME_BEFORE)
}

// Java: CheckEs6ModuleFileStructureTest.testGoogRequireAfterImportsAndDeclareNamesButBeforeOthers()
func (c *CheckEs6ModuleFileStructureTest) TestGoogRequireAfterImportsAndDeclareNamesButBeforeOthers() {
	c.TestSame("const bar = goog.require('bar');")
	c.TestSame("const bar = goog.require('bar'); let notAModule;")
	c.TestSame("let notAModule; const bar = goog.require('bar');")
	c.TestSame("const bar = goog.require('bar'); export {};")
	c.TestSame("export {}; const bar = goog.require('bar');", CheckEs6ModuleFileStructure_MUST_COME_BEFORE)
	c.TestSame("goog.declareModuleId('name'); const bar = goog.require('bar'); export {};")
	c.TestWarning(
		"const bar = goog.require('bar'); goog.declareModuleId('name'); export {};",
		CheckEs6ModuleFileStructure_MUST_COME_BEFORE)
	c.TestSame("import 'file'; const bar = goog.require('bar');")
	c.TestWarning("const bar = goog.require('bar'); import 'file';", CheckEs6ModuleFileStructure_MUST_COME_BEFORE)
	c.TestSame("import 'file'; goog.require('bar');")
	c.TestWarning("goog.require('bar'); import 'file';", CheckEs6ModuleFileStructure_MUST_COME_BEFORE)
	c.TestSame(
		`import 'file';
goog.declareModuleId('name');
const bar = goog.require('bar');
let rest;
`)
}
