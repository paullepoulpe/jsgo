// Port of com/google/javascript/jscomp/lint/CheckGoogModuleTypeScriptNameTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/lint/CheckGoogModuleTypeScriptNameTest.java)

/*
 * Copyright 2021 The Closure Compiler Authors.
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

// Tests for {@link CheckGoogModuleTypeScriptName}.
type CheckGoogModuleTypeScriptNameTest struct {
	CompilerTestCaseBase
}

func NewCheckGoogModuleTypeScriptNameTest() *CheckGoogModuleTypeScriptNameTest {
	c := &CheckGoogModuleTypeScriptNameTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestCheckGoogModuleTypeScriptNameTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCheckGoogModuleTypeScriptNameTest() })
}

// Java: CheckGoogModuleTypeScriptNameTest.getProcessor(Compiler)
func (c *CheckGoogModuleTypeScriptNameTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewCheckGoogModuleTypeScriptName(compiler)
}

// Java: CheckGoogModuleTypeScriptNameTest.getOptions()
func (c *CheckGoogModuleTypeScriptNameTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetWarningLevel(DiagnosticGroups_LINT_CHECKS, CheckLevel_WARNING)
	return options
}

// Java: CheckGoogModuleTypeScriptNameTest.testNoWarning()
//
// Go-only: test-helper — a trailing underscore: the plain name is the
// inherited CompilerTestCase.testNoWarning(...), which this no-argument @Test method
// overloads in Java.
func (c *CheckGoogModuleTypeScriptNameTest) TestNoWarning_() {
	c.TestNoWarning(
		c.Srcs(
			SourceFile_FromCode(
				"gws/js/test.js", `goog.module('google3.gws.js.test');
alert(1);`)))
}

// Java: CheckGoogModuleTypeScriptNameTest.testNoWarning_notOnAllowList()
func (c *CheckGoogModuleTypeScriptNameTest) TestNoWarning_notOnAllowList() {
	c.TestNoWarning(
		c.Srcs(
			SourceFile_FromCode("javascript/apps/foo/test.js", `goog.module('test');
alert(1);`)))
}

// Java: CheckGoogModuleTypeScriptNameTest.testNoWarning_noModule()
func (c *CheckGoogModuleTypeScriptNameTest) TestNoWarning_noModule() {
	c.TestNoWarning(
		`goog.provide('test');
alert(1)
`)
}

// Java: CheckGoogModuleTypeScriptNameTest.testWarning()
//
// Go-only: test-helper — a trailing underscore: the plain name is the
// inherited CompilerTestCase.testWarning(...), which this no-argument @Test method
// overloads in Java.
func (c *CheckGoogModuleTypeScriptNameTest) TestWarning_() {
	c.Test(
		c.Srcs(
			[]*SourceFile{
				SourceFile_FromCode("gws/js/test.js", `goog.module('test');
alert(1);`)}),
		c.Warning(CheckGoogModuleTypeScriptName_MODULE_NAMESPACE_MISMATCHES_TYPESCRIPT_NAMESPACE).WithMessageContaining("The correct namespace is: \"google3.gws.js.test\""))
}

// MOE::begin_strip
// Java: CheckGoogModuleTypeScriptNameTest.testWarningWithFullSourceName()
func (c *CheckGoogModuleTypeScriptNameTest) TestWarningWithFullSourceName() {
	c.Test(
		c.Srcs(
			[]*SourceFile{
				SourceFile_FromCode(
					"/google/src/cloud/someuser/someworkspace/google3/gws/js/test.js",
					`goog.module('test');
alert(1);`)}),
		c.Warning(CheckGoogModuleTypeScriptName_MODULE_NAMESPACE_MISMATCHES_TYPESCRIPT_NAMESPACE).WithMessageContaining("The correct namespace is: \"google3.gws.js.test\""))
}
