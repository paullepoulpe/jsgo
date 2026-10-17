// Port of com/google/javascript/jscomp/lint/CheckUnusedLabelsTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/lint/CheckUnusedLabelsTest.java)

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

// Test case for {@link CheckUnusedLabels}.
type CheckUnusedLabelsTest struct {
	CompilerTestCaseBase
}

func NewCheckUnusedLabelsTest() *CheckUnusedLabelsTest {
	c := &CheckUnusedLabelsTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestCheckUnusedLabelsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCheckUnusedLabelsTest() })
}

// Java: CheckUnusedLabelsTest.getProcessor(Compiler)
func (c *CheckUnusedLabelsTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewCheckUnusedLabels(compiler)
}

// Java: CheckUnusedLabelsTest.getOptions()
func (c *CheckUnusedLabelsTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetWarningLevel(DiagnosticGroups_LINT_CHECKS, CheckLevel_WARNING)
	return options
}

// Java: CheckUnusedLabelsTest.testCheckUnusedLabels_noWarning()
func (c *CheckUnusedLabelsTest) TestCheckUnusedLabels_noWarning() {
	c.TestSame("L: if (true) { break L; }")
	c.TestSame("L: for (;;) { if (true) { break L; } }")
	c.TestSame("L1: L2: if (true) { if (true) { break L1; } break L2; }")
}

// Java: CheckUnusedLabelsTest.testCheckUnusedLabels_warning()
func (c *CheckUnusedLabelsTest) TestCheckUnusedLabels_warning() {
	c.TestWarning("L: var x = 0;", CheckUnusedLabels_UNUSED_LABEL)
	c.TestWarning("L: { f(); }", CheckUnusedLabels_UNUSED_LABEL)
	c.TestWarning("L: for (;;) {}", CheckUnusedLabels_UNUSED_LABEL)
	c.TestWarning("L1: for (;;) { L2: if (true) { break L2; } }", CheckUnusedLabels_UNUSED_LABEL)
	c.TestWarning("() => {a: 2}", CheckUnusedLabels_UNUSED_LABEL)
}
