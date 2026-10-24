// Port of com/google/javascript/jscomp/lint/CheckEs6ModulesTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/lint/CheckEs6ModulesTest.java)

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

type CheckEs6ModulesTest struct {
	CompilerTestCaseBase
}

func NewCheckEs6ModulesTest() *CheckEs6ModulesTest {
	c := &CheckEs6ModulesTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestCheckEs6ModulesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCheckEs6ModulesTest() })
}

// Java: CheckEs6ModulesTest.setUp()
func (c *CheckEs6ModulesTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.IgnoreWarnings(ModuleLoader_INVALID_MODULE_PATH)
}

// Java: CheckEs6ModulesTest.getProcessor(Compiler)
func (c *CheckEs6ModulesTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewCheckEs6Modules(compiler)
}

// Java: CheckEs6ModulesTest.testDuplicateImports()
func (c *CheckEs6ModulesTest) TestDuplicateImports() {
	c.TestSame("import singleImport from 'file';")
	c.TestSame("import first from 'first'; import second from 'second';")
	c.TestWarning("import * as first from 'file'; import {second} from 'file';", CheckEs6Modules_DUPLICATE_IMPORT)
}

// Java: CheckEs6ModulesTest.testNoDefaultExport()
func (c *CheckEs6ModulesTest) TestNoDefaultExport() {
	c.TestWarning("export default 0", CheckEs6Modules_NO_DEFAULT_EXPORT)
}
