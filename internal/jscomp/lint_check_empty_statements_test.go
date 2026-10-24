// Port of com/google/javascript/jscomp/lint/CheckEmptyStatementsTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/lint/CheckEmptyStatementsTest.java)

/*
 * Copyright 2015 The Closure Compiler Authors.
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

// Test case for {@link CheckEmptyStatements}.
type CheckEmptyStatementsTest struct {
	CompilerTestCaseBase
}

func NewCheckEmptyStatementsTest() *CheckEmptyStatementsTest {
	c := &CheckEmptyStatementsTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestCheckEmptyStatementsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCheckEmptyStatementsTest() })
}

// Java: CheckEmptyStatementsTest.getProcessor(Compiler)
func (c *CheckEmptyStatementsTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewCheckEmptyStatements(compiler)
}

// Java: CheckEmptyStatementsTest.getOptions()
func (c *CheckEmptyStatementsTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetWarningLevel(DiagnosticGroups_LINT_CHECKS, CheckLevel_WARNING)
	return options
}

// TestWarning adds the private overload testWarning(String) to the
// inherited testWarning(...) overloads.
//
// Go-only: language-adapter — Java picks the overload statically; the
// harness's TestWarning(args ...any) picks it at run time from the argument
// types, and so does this override: one string argument is this class's
// overload, anything else the inherited ones.
//
// Java: CheckEmptyStatementsTest.testWarning(String)
func (c *CheckEmptyStatementsTest) TestWarning(args ...any) {
	if len(args) == 1 {
		if js, ok := args[0].(string); ok {
			c.TestWarning(js, CheckEmptyStatements_USELESS_EMPTY_STATEMENT)
			return
		}
	}
	c.CompilerTestCaseBase.TestWarning(args...)
}

// Java: CheckEmptyStatementsTest.testWarning()
//
// Go-only: test-helper — a trailing underscore: the plain name is the
// inherited CompilerTestCase.testWarning(...), which this no-argument @Test method
// overloads in Java.
func (c *CheckEmptyStatementsTest) TestWarning_() {
	c.TestWarning("function f() {};")
	c.TestWarning("var x;;")
	c.TestWarning("alert(1);;")
}

// Java: CheckEmptyStatementsTest.testWarning_withES6Modules01()
func (c *CheckEmptyStatementsTest) TestWarning_withES6Modules01() {
	c.TestWarning("export function f() {};;")
}

// Java: CheckEmptyStatementsTest.testWarning_withES6Modules02()
func (c *CheckEmptyStatementsTest) TestWarning_withES6Modules02() {
	c.TestWarning("export var x;;")
}

// Java: CheckEmptyStatementsTest.testNoWarning()
//
// Go-only: test-helper — a trailing underscore: the plain name is the
// inherited CompilerTestCase.testNoWarning(...), which this no-argument @Test method
// overloads in Java.
func (c *CheckEmptyStatementsTest) TestNoWarning_() {
	c.TestSame("function f() {}")
	c.TestSame("var x;")
	c.TestSame("alert(1);")
	c.TestSame("if (x); y;")
}

// Java: CheckEmptyStatementsTest.testNoWarning_withES6Modules01()
func (c *CheckEmptyStatementsTest) TestNoWarning_withES6Modules01() {
	c.TestSame("export function f() {}")
}

// Java: CheckEmptyStatementsTest.testNoWarning_withES6Modules02()
func (c *CheckEmptyStatementsTest) TestNoWarning_withES6Modules02() {
	c.TestSame("export var x;")
}
