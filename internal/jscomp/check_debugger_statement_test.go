// Port of com/google/javascript/jscomp/CheckDebuggerStatementTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/CheckDebuggerStatementTest.java)

/*
 * Copyright 2011 The Closure Compiler Authors.
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

// CheckDebuggerStatementTest: CheckDebuggerStatementTest is a unit test for CheckDebuggerStatement.
type CheckDebuggerStatementTest struct {
	CompilerTestCaseBase
	// A Java @Nullable CheckLevel: nil is null.
	checkLevel *CheckLevel
}

// Java: CheckDebuggerStatementTest()
func NewCheckDebuggerStatementTest() *CheckDebuggerStatementTest {
	c := &CheckDebuggerStatementTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestCheckDebuggerStatementTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCheckDebuggerStatementTest() })
}

// Java: CheckDebuggerStatementTest.tearDown()
func (c *CheckDebuggerStatementTest) TearDown() {
	c.checkLevel = nil
}

// Java: CheckDebuggerStatementTest.getProcessor(Compiler)
func (c *CheckDebuggerStatementTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewCheckDebuggerStatement(compiler)
}

// Java: CheckDebuggerStatementTest.getOptions()
func (c *CheckDebuggerStatementTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	if c.checkLevel != nil {
		options.SetWarningLevel(DiagnosticGroups_DEBUGGER_STATEMENT_PRESENT, *c.checkLevel)
	}
	return options
}

// Java: CheckDebuggerStatementTest.testCheckDebuggerStatement()
func (c *CheckDebuggerStatementTest) TestCheckDebuggerStatement() {
	c.checkLevel = CheckLevel_WARNING.Ptr()

	c.TestWarning("debugger;", CheckDebuggerStatement_DEBUGGER_STATEMENT_PRESENT)
	c.TestWarning("function foo() { debugger; }", CheckDebuggerStatement_DEBUGGER_STATEMENT_PRESENT)
}

// Java: CheckDebuggerStatementTest.testCheckIsDisabledByDefault()
func (c *CheckDebuggerStatementTest) TestCheckIsDisabledByDefault() {
	c.checkLevel = nil

	c.TestSame("debugger;")
	c.TestSame("function foo() { debugger; }")
}

// Java: CheckDebuggerStatementTest.testNoWarningWhenExplicitlyDisabled()
func (c *CheckDebuggerStatementTest) TestNoWarningWhenExplicitlyDisabled() {
	c.checkLevel = CheckLevel_OFF.Ptr()

	c.TestSame("debugger;")
	c.TestSame("function foo() { debugger; }")
}

// Java: CheckDebuggerStatementTest.testCheckDebuggerKeywordMayAppearInComments()
func (c *CheckDebuggerStatementTest) TestCheckDebuggerKeywordMayAppearInComments() {
	c.checkLevel = CheckLevel_WARNING.Ptr()

	c.Test("// I like the debugger; it is helpful.", "")
}

// Java: CheckDebuggerStatementTest.testCheckDebuggerStatementInEval()
func (c *CheckDebuggerStatementTest) TestCheckDebuggerStatementInEval() {
	c.checkLevel = CheckLevel_WARNING.Ptr()

	c.TestSame("eval('debugger');")
}
