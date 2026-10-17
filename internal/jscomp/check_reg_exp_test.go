// Port of com/google/javascript/jscomp/CheckRegExpTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/CheckRegExpTest.java)

/*
 * Copyright 2010 The Closure Compiler Authors.
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

// CheckRegExpTest: @author johnlenz@google.com (John Lenz)
type CheckRegExpTest struct {
	CompilerTestCaseBase
	last                               *CheckRegExp
	reportErrors                       bool
	assumeAllGlobalRegexpUsagesVisible bool
}

// Java: CheckRegExpTest()
func NewCheckRegExpTest() *CheckRegExpTest {
	c := &CheckRegExpTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, "var RegExp;")
	c.last = nil
	return c
}

func TestCheckRegExpTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCheckRegExpTest() })
}

// Java: CheckRegExpTest.setUp()
func (c *CheckRegExpTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.reportErrors = true
	c.assumeAllGlobalRegexpUsagesVisible = true
}

// Java: CheckRegExpTest.getOptions()
func (c *CheckRegExpTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetWarningLevel(DiagnosticGroups_CHECK_REGEXP, CheckLevel_WARNING)
	return options
}

// Java: CheckRegExpTest.getProcessor(Compiler)
func (c *CheckRegExpTest) GetProcessor(compiler *Compiler) CompilerPass {
	c.last = NewCheckRegExp(compiler, c.assumeAllGlobalRegexpUsagesVisible, c.reportErrors)
	return c.last
}

// Java: CheckRegExpTest.testReference(String, boolean)
func (c *CheckRegExpTest) TestReference(code string, expected bool) {
	if expected && c.reportErrors {
		c.TestWarning(code, CheckRegExp_REGEXP_REFERENCE)
	} else {
		c.TestSame(code)
	}
	AssertThat(c.last.IsGlobalRegExpPropertiesUsed()).IsEqualTo(expected)
}

// Java: CheckRegExpTest.testRegExp()
func (c *CheckRegExpTest) TestRegExp() {
	// Creating RegExp's is OK.
	c.TestReference("RegExp();", false)
	c.TestReference("var x = RegExp();", false)
	c.TestReference("new RegExp();", false)
	c.TestReference("var x = new RegExp();", false)

	// Checking for RegExp instances is OK, as well.
	c.TestReference("x instanceof RegExp;", false)

	// Comparing for equality with RegExp is OK.
	c.TestReference("x === RegExp;", false)
	c.TestReference("x !== RegExp;", false)
	c.TestReference("switch (x) { case RegExp: }", false)
	c.TestReference("x == RegExp;", false)
	c.TestReference("x != RegExp;", false)

	// Access to non-magical properties is OK.
	c.TestReference("RegExp.test();", false)
	c.TestReference("var x = RegExp.test();", false)
	c.TestReference("RegExp.exec();", false)
	c.TestReference("RegExp.foobar;", false)

	// Magical properties aren't allowed.
	c.TestReference("RegExp.$1;", true)
	c.TestReference("RegExp.$_;", true)
	c.TestReference("RegExp.$input;", true)
	c.TestReference("RegExp.rightContext;", true)
	c.TestReference("RegExp.multiline;", true)

	// Any other reference isn't allowed.
	c.TestReference("delete RegExp;", true)
	c.TestReference("RegExp;", true)
	c.TestReference("if (RegExp);", true)
	c.TestReference("if (!RegExp);", true)

	// Aliases aren't allowed.
	c.TestReference("var x = RegExp;", true)
	c.TestReference("f(RegExp);", true)
	c.TestReference("new f(RegExp);", true)
	c.TestReference("var x = RegExp; x.test()", true)
	c.TestReference("let x = RegExp;", true)
	c.TestReference("const x = RegExp;", true)

	// No RegExp reference is OK.
	c.TestReference("var x;", false)

	// Local RegExp is OK.
	c.TestReference("function f() {var RegExp; RegExp.test();}", false)
	c.TestReference("function f() {let RegExp; RegExp.test();}", false)
	c.TestReference("function *gen() {var RegExp; yield RegExp.test();}", false)

	// Property named 'RegExp' is OK.
	c.TestReference("var x = {RegExp: {}}; x.RegExp.$1;", false)

	// Class property is also OK.
	c.TestReference(
		`class x {
  constructor() {this.RegExp = {};}
  method() {
    this.RegExp.$1;
    this.RegExp.test();
  }
}
`,
		false)
}

// Java: CheckRegExpTest.testDisableErrorReporting()
func (c *CheckRegExpTest) TestDisableErrorReporting() {
	c.reportErrors = true

	c.TestWarning("RegExp.$1;", CheckRegExp_REGEXP_REFERENCE)

	c.reportErrors = false

	c.TestNoWarning("RegExp.$1;")
}

// Java: CheckRegExpTest.testInvalidRange()
func (c *CheckRegExpTest) TestInvalidRange() {
	c.TestWarning("\"asdf\".match(/[z-a]/)", CheckRegExp_MALFORMED_REGEXP)
}

// Java: CheckRegExpTest.testCanAssumeAllGlobalRegexpUsagesVisible()
func (c *CheckRegExpTest) TestCanAssumeAllGlobalRegexpUsagesVisible() {
	c.assumeAllGlobalRegexpUsagesVisible = true

	c.TestSame("")

	AssertThat(c.last.IsGlobalRegExpPropertiesUsed()).IsFalse()
}

// Java: CheckRegExpTest.testCannotAssumeAllGlobalRegexpUsagesVisible()
func (c *CheckRegExpTest) TestCannotAssumeAllGlobalRegexpUsagesVisible() {
	c.assumeAllGlobalRegexpUsagesVisible = false

	c.TestSame("")

	AssertThat(c.last.IsGlobalRegExpPropertiesUsed()).IsTrue()
}
