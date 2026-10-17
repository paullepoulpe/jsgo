// Port of com/google/javascript/jscomp/lint/CheckDuplicateCaseTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/lint/CheckDuplicateCaseTest.java)

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

// Test case for {@link CheckDuplicateCase}.
type CheckDuplicateCaseTest struct {
	CompilerTestCaseBase
}

func NewCheckDuplicateCaseTest() *CheckDuplicateCaseTest {
	c := &CheckDuplicateCaseTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestCheckDuplicateCaseTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCheckDuplicateCaseTest() })
}

// Java: CheckDuplicateCaseTest.getProcessor(Compiler)
func (c *CheckDuplicateCaseTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewCheckDuplicateCase(compiler)
}

// Java: CheckDuplicateCaseTest.testCheckDuplicateCase_noWarning()
func (c *CheckDuplicateCaseTest) TestCheckDuplicateCase_noWarning() {
	c.TestSame("switch (foo) { case 1: break; case 2: break; }")
	c.TestSame("switch (foo) { case 1: break; case '1': break; }")
	c.TestSame("switch (foo) { case bar: break; case 'bar': break; }")
	c.TestSame("switch (foo) { case 1: break; } switch (bar) { case 1: break; }")
}

// Java: CheckDuplicateCaseTest.testCheckDuplicateCase_warning()
func (c *CheckDuplicateCaseTest) TestCheckDuplicateCase_warning() {
	c.TestWarning("switch (foo) { case 1: var a = 3; case 1: break; }", CheckDuplicateCase_DUPLICATE_CASE)
	c.TestWarning("switch (foo) { case 1: break; case 2: break; case 1: break; }", CheckDuplicateCase_DUPLICATE_CASE)
	c.TestWarning("switch (foo) { case '1': break; case '1': break; }", CheckDuplicateCase_DUPLICATE_CASE)
	c.TestWarning("switch (foo) { case bar: break; case bar: break; }", CheckDuplicateCase_DUPLICATE_CASE)
	c.TestWarning("switch (foo) { case i++: break; case i++: break; }", CheckDuplicateCase_DUPLICATE_CASE)
}
