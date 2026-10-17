// Port of com/google/javascript/jscomp/lint/CheckUselessBlocksTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/lint/CheckUselessBlocksTest.java)

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

// Test case for {@link CheckUselessBlocks}.
type CheckUselessBlocksTest struct {
	CompilerTestCaseBase
}

func NewCheckUselessBlocksTest() *CheckUselessBlocksTest {
	c := &CheckUselessBlocksTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestCheckUselessBlocksTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCheckUselessBlocksTest() })
}

// Java: CheckUselessBlocksTest.getProcessor(Compiler)
func (c *CheckUselessBlocksTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewCheckUselessBlocks(compiler)
}

// Java: CheckUselessBlocksTest.getOptions()
func (c *CheckUselessBlocksTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetWarningLevel(DiagnosticGroups_LINT_CHECKS, CheckLevel_WARNING)
	return options
}

// Java: CheckUselessBlocksTest.testCheckUselessBlocks_noWarning()
func (c *CheckUselessBlocksTest) TestCheckUselessBlocks_noWarning() {
	c.TestSame("while (foo) { bar(); }")
	c.TestSame("if (true) { var x = 1; }")
	c.TestSame("if (foo) { if (bar) { baz(); } }")
	c.TestSame("function bar() { baz(); }")
	c.TestSame("function f() { switch (x) { case 1: { return 5; } } }")
	c.TestSame("blah: { break blah; }")
	// TODO(moz): For block-scoped function declaration, we should technically
	// warn if we are in non-strict mode and the language mode is ES5 or below.
	c.TestSame("{ function foo() {} }")
	c.TestSame("let x = 1;")
	c.TestSame("{ let x = 1; }")
	c.TestSame("if (true) { let x = 1; }")
	c.TestSame("{ const y = 1; }")
	c.TestSame("{ class Foo {} }")
}

// Java: CheckUselessBlocksTest.testCheckUselessBlocks_withES6Modules_noWarning()
func (c *CheckUselessBlocksTest) TestCheckUselessBlocks_withES6Modules_noWarning() {
	c.TestSame("export function f() { switch (x) { case 1: { return 5; } } }")
}

// Java: CheckUselessBlocksTest.testCheckUselessBlocks_warning()
func (c *CheckUselessBlocksTest) TestCheckUselessBlocks_warning() {
	c.TestWarning("{}", CheckUselessBlocks_USELESS_BLOCK)
	c.TestWarning("{ var f = function() {}; }", CheckUselessBlocks_USELESS_BLOCK)
	c.TestWarning("{ var x = 1; }", CheckUselessBlocks_USELESS_BLOCK)
	c.TestWarning(
		`function f() {
  return
    {foo: 'bar'};
}
`,
		CheckUselessBlocks_USELESS_BLOCK)
	c.TestWarning(
		`if (foo) {
  bar();
  {
    baz();
  }
}
`,
		CheckUselessBlocks_USELESS_BLOCK)
	c.TestWarning(
		`if (foo) {
  bar();
} {
  baz();
}
`,
		CheckUselessBlocks_USELESS_BLOCK)
	c.TestWarning("function bar() { { baz(); } }", CheckUselessBlocks_USELESS_BLOCK)
	c.TestWarning("{ let x = function() {}; {} }", CheckUselessBlocks_USELESS_BLOCK)
	c.TestWarning("{ let x = function() { {} }; }", CheckUselessBlocks_USELESS_BLOCK)
	c.TestWarning("{ var f = class {}; }", CheckUselessBlocks_USELESS_BLOCK)
}

// Java: CheckUselessBlocksTest.testCheckUselessBlocks_withES6Modules_warning()
func (c *CheckUselessBlocksTest) TestCheckUselessBlocks_withES6Modules_warning() {
	c.TestWarning("export function bar() { { baz(); } }", CheckUselessBlocks_USELESS_BLOCK)
}
