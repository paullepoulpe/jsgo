// Port of com/google/javascript/jscomp/ForbidDynamicImportUsageTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/ForbidDynamicImportUsageTest.java)

/*
 * Copyright 2020 The Closure Compiler Authors.
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

// ForbidDynamicImportUsageTest: Test Dynamic Import usage issues warning
type ForbidDynamicImportUsageTest struct {
	CompilerTestCaseBase
}

// Java: ForbidDynamicImportUsageTest()
func NewForbidDynamicImportUsageTest() *ForbidDynamicImportUsageTest {
	c := &ForbidDynamicImportUsageTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestForbidDynamicImportUsageTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewForbidDynamicImportUsageTest() })
}

// Java: ForbidDynamicImportUsageTest.setUp()
func (c *ForbidDynamicImportUsageTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()

	c.EnableTypeInfoValidation()
	c.EnableTypeCheck()
}

// Java: ForbidDynamicImportUsageTest.getProcessor(Compiler)
func (c *ForbidDynamicImportUsageTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewForbidDynamicImportUsage(compiler)
}

// Java: ForbidDynamicImportUsageTest.testDynamicImportError()
func (c *ForbidDynamicImportUsageTest) TestDynamicImportError() {
	c.IgnoreWarnings(DiagnosticGroups_MODULE_LOAD)
	c.TestError("import('./foo.js')", ForbidDynamicImportUsage_DYNAMIC_IMPORT_USAGE)
}
