// Port of com/google/javascript/jscomp/lint/CheckMissingSemicolonTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/lint/CheckMissingSemicolonTest.java)

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

type CheckMissingSemicolonTest struct {
	CompilerTestCaseBase
}

func NewCheckMissingSemicolonTest() *CheckMissingSemicolonTest {
	c := &CheckMissingSemicolonTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestCheckMissingSemicolonTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCheckMissingSemicolonTest() })
}

// Java: CheckMissingSemicolonTest.getProcessor(Compiler)
func (c *CheckMissingSemicolonTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewCheckMissingSemicolon(compiler)
}

// Java: CheckMissingSemicolonTest.getOptions()
func (c *CheckMissingSemicolonTest) GetOptions() *CompilerOptions {
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
// Java: CheckMissingSemicolonTest.testWarning(String)
func (c *CheckMissingSemicolonTest) TestWarning(args ...any) {
	if len(args) == 1 {
		if js, ok := args[0].(string); ok {
			c.TestWarning(js, CheckMissingSemicolon_MISSING_SEMICOLON)
			return
		}
	}
	c.CompilerTestCaseBase.TestWarning(args...)
}

// Java: CheckMissingSemicolonTest.testWarning()
//
// Go-only: test-helper — a trailing underscore: the plain name is the
// inherited CompilerTestCase.testWarning(...), which this no-argument @Test method
// overloads in Java.
func (c *CheckMissingSemicolonTest) TestWarning_() {
	c.TestWarning("var x")
	c.TestWarning("alert(1)")
	c.TestWarning("do { things; } while (true)")
}

// Java: CheckMissingSemicolonTest.testNoWarning_withSemicolon()
func (c *CheckMissingSemicolonTest) TestNoWarning_withSemicolon() {
	c.TestSame("var x;")
	c.TestSame("alert(1);")
}

// Java: CheckMissingSemicolonTest.testNoWarning_controlStructure()
func (c *CheckMissingSemicolonTest) TestNoWarning_controlStructure() {
	c.TestSame("if (true) {}")
	c.TestSame("while (true) {}")
	c.TestSame("do { things; } while (true);")
	c.TestSame("switch (n) { case 3: alert(4); }")
	c.TestSame("switch (n) { case 5: { alert(6); } }")
	c.TestSame("alert(1); LABEL: while (true) { alert(2); }")
	c.TestSame("for (;;) {}")
	c.TestSame("for (x of y) {}")
	c.TestSame("for (x in y) {}")
}

// Java: CheckMissingSemicolonTest.testNoWarning_functionOrClass()
func (c *CheckMissingSemicolonTest) TestNoWarning_functionOrClass() {
	c.TestSame("function f() {}")
	c.TestSame("function* f() {}")
	c.TestSame("class Example {}")
}

// Java: CheckMissingSemicolonTest.testWarning_export()
func (c *CheckMissingSemicolonTest) TestWarning_export() {
	c.IgnoreWarnings(DiagnosticGroups_MODULE_LOAD)
	c.TestWarning("export var x = 3")
	c.TestWarning("export * from './other'")
}

// Java: CheckMissingSemicolonTest.testNoWarning_export()
func (c *CheckMissingSemicolonTest) TestNoWarning_export() {
	c.TestSame("export function f() {}")
	c.TestSame("export class C {}")
}
