// Port of com/google/javascript/jscomp/lint/CheckArrayWithGoogObjectTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/lint/CheckArrayWithGoogObjectTest.java)

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

// Test case for {@link CheckArrayWithGoogObject}.
type CheckArrayWithGoogObjectTest struct {
	CompilerTestCaseBase
}

func NewCheckArrayWithGoogObjectTest() *CheckArrayWithGoogObjectTest {
	c := &CheckArrayWithGoogObjectTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestCheckArrayWithGoogObjectTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCheckArrayWithGoogObjectTest() })
}

var checkArrayWithGoogObjectTest_GOOG_OBJECT = `goog.object = {};
goog.object.forEach = function(obj, f, opt_this) {}
`

// Java: CheckArrayWithGoogObjectTest.getProcessor(Compiler)
func (c *CheckArrayWithGoogObjectTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewCheckArrayWithGoogObject(compiler)
}

// Java: CheckArrayWithGoogObjectTest.getOptions()
func (c *CheckArrayWithGoogObjectTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetWarningLevel(DiagnosticGroups_LINT_CHECKS, CheckLevel_WARNING)
	return options
}

// Java: CheckArrayWithGoogObjectTest.setUp()
func (c *CheckArrayWithGoogObjectTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableTranspile()
	// Transpilation runs normalization, which does renaming on the externs.
	// It also runs a pass that reverts the renaming, but it ignores the externs, leaving them
	// changed.
	c.AllowExternsChanges()
	c.EnableTypeCheck()
}

// Java: CheckArrayWithGoogObjectTest.testGoogObjectForEach1()
func (c *CheckArrayWithGoogObjectTest) TestGoogObjectForEach1() {
	c.TestGoogObjectWarning(
		checkArrayWithGoogObjectTest_GOOG_OBJECT +
			`var arr = [1, 2, 3];
goog.object.forEach(arr, alert);
`)
}

// Java: CheckArrayWithGoogObjectTest.testGoogObjectForEach2()
func (c *CheckArrayWithGoogObjectTest) TestGoogObjectForEach2() {
	c.TestGoogObjectWarning(
		checkArrayWithGoogObjectTest_GOOG_OBJECT +
			`function f(/** Array<number>|number */ n) {
  if (typeof n == 'number')
    alert(n);
  else
    goog.object.forEach(n, alert);
}
`)
}

// Java: CheckArrayWithGoogObjectTest.testGoogObjectForEach3()
func (c *CheckArrayWithGoogObjectTest) TestGoogObjectForEach3() {
	c.TestGoogObjectWarning(
		checkArrayWithGoogObjectTest_GOOG_OBJECT +
			`function f(/** !Array<number> */ arr) {
  goog.object.forEach(arr, alert);
}
`)
}

// Java: CheckArrayWithGoogObjectTest.testGoogObjectForEach4()
func (c *CheckArrayWithGoogObjectTest) TestGoogObjectForEach4() {
	c.TestNoGoogObjectWarning(
		checkArrayWithGoogObjectTest_GOOG_OBJECT +
			`function f(/** Object<string, number> */ obj) {
  goog.object.forEach(obj, alert);
}
`)
}

// Java: CheckArrayWithGoogObjectTest.testGoogObjectWarning(String)
func (c *CheckArrayWithGoogObjectTest) TestGoogObjectWarning(js string) {
	c.Test(c.Externs(CompilerTestCase_DEFAULT_EXTERNS), c.Srcs(js), c.Warning(CheckArrayWithGoogObject_ARRAY_PASSED_TO_GOOG_OBJECT))
}

// Java: CheckArrayWithGoogObjectTest.testNoGoogObjectWarning(String)
func (c *CheckArrayWithGoogObjectTest) TestNoGoogObjectWarning(js string) {
	c.TestSame(c.Externs(CompilerTestCase_DEFAULT_EXTERNS), c.Srcs(js))
}
