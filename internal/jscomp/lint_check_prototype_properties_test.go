// Port of com/google/javascript/jscomp/lint/CheckPrototypePropertiesTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/lint/CheckPrototypePropertiesTest.java)

/*
 * Copyright 2014 The Closure Compiler Authors.
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

type CheckPrototypePropertiesTest struct {
	CompilerTestCaseBase
}

func NewCheckPrototypePropertiesTest() *CheckPrototypePropertiesTest {
	c := &CheckPrototypePropertiesTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestCheckPrototypePropertiesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCheckPrototypePropertiesTest() })
}

// Java: CheckPrototypePropertiesTest.getProcessor(Compiler)
func (c *CheckPrototypePropertiesTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewCheckPrototypeProperties(compiler)
}

// Java: CheckPrototypePropertiesTest.getOptions()
func (c *CheckPrototypePropertiesTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetWarningLevel(DiagnosticGroups_LINT_CHECKS, CheckLevel_WARNING)
	return options
}

// Java: CheckPrototypePropertiesTest.testNoWarning()
//
// Go-only: test-helper — a trailing underscore: the plain name is the
// inherited CompilerTestCase.testNoWarning(...), which this no-argument @Test method
// overloads in Java.
func (c *CheckPrototypePropertiesTest) TestNoWarning_() {
	c.TestSame("function C() {}; C.prototype.foo = null;")
	c.TestSame("function C() {}; C.prototype.foo = undefined;")
	c.TestSame("function C() {}; C.prototype.foo;")
	c.TestSame("function C() {}; C.prototype.foo = 0;")
	c.TestSame("function C() {}; C.prototype.foo = 'someString';")
	c.TestSame("function C() {}; C.prototype.foo = function() {};")
	c.TestSame("function C() {}; /** @enum {number} */ C.prototype.foo = { BAR: 0 };")
}

// Java: CheckPrototypePropertiesTest.testNoWarning_withES6Modules()
func (c *CheckPrototypePropertiesTest) TestNoWarning_withES6Modules() {
	c.TestSame("export function C() {}; C.prototype.foo = null;")
}

// Java: CheckPrototypePropertiesTest.testWarnings()
func (c *CheckPrototypePropertiesTest) TestWarnings() {
	c.TestSame("function C() {}; C.prototype.foo = [];", CheckPrototypeProperties_ILLEGAL_PROTOTYPE_MEMBER)
	c.TestSame("function C() {}; C.prototype.foo = {};", CheckPrototypeProperties_ILLEGAL_PROTOTYPE_MEMBER)
	c.TestSame("function C() {}; C.prototype.foo = { BAR: 0 };", CheckPrototypeProperties_ILLEGAL_PROTOTYPE_MEMBER)
}

// Java: CheckPrototypePropertiesTest.testWarnings_withES6Modules()
func (c *CheckPrototypePropertiesTest) TestWarnings_withES6Modules() {
	c.TestSame("export function C() {}; C.prototype.foo = [];", CheckPrototypeProperties_ILLEGAL_PROTOTYPE_MEMBER)
}
