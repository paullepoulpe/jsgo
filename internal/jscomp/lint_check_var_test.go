// Port of com/google/javascript/jscomp/lint/CheckVarTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/lint/CheckVarTest.java)

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

// Test case for {@link CheckVar}.
type CheckVarTest struct {
	CompilerTestCaseBase
}

func NewCheckVarTest() *CheckVarTest {
	c := &CheckVarTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestCheckVarTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCheckVarTest() })
}

// Java: CheckVarTest.getProcessor(Compiler)
func (c *CheckVarTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewCheckVar(compiler)
}

// Java: CheckVarTest.getOptions()
func (c *CheckVarTest) GetOptions() *CompilerOptions {
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
// Java: CheckVarTest.testWarning(String)
func (c *CheckVarTest) TestWarning(args ...any) {
	if len(args) == 1 {
		if js, ok := args[0].(string); ok {
			c.TestWarning(js, CheckVar_VAR)
			return
		}
	}
	c.CompilerTestCaseBase.TestWarning(args...)
}

// Java: CheckVarTest.testWarning()
//
// Go-only: test-helper — a trailing underscore: the plain name is the
// inherited CompilerTestCase.testWarning(...), which this no-argument @Test method
// overloads in Java.
func (c *CheckVarTest) TestWarning_() {
	c.TestWarning("var x;")
	c.TestWarning("var x = 12;")
	c.TestWarning("export var x;")
	c.TestWarning("function f() { var x = 12; return x; }")
}

// Java: CheckVarTest.testSuppressWarning()
func (c *CheckVarTest) TestSuppressWarning() {
	c.TestSame(
		`/**
 * @fileoverview
 * @suppress {lintVarDeclarations}
 */
var x;
var x12 = 12;
export var x;
function f() { var x = 12; return x; }
`)
}

// Java: CheckVarTest.testNoWarning()
//
// Go-only: test-helper — a trailing underscore: the plain name is the
// inherited CompilerTestCase.testNoWarning(...), which this no-argument @Test method
// overloads in Java.
func (c *CheckVarTest) TestNoWarning_() {
	c.TestSame("function f() { return 'var'; }")
	c.TestSame("let x;")
	c.TestSame("let x = 12;")
	c.TestSame("const x = 12;")
	c.TestSame("function f() { let x; x = 12; return x; }")
	c.TestSame("function f() { const x = 12; return x; }")
}
