// Port of com/google/javascript/jscomp/CheckAccessControlsPolymerTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/CheckAccessControlsPolymerTest.java)

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

// CheckAccessControlsPolymerTest is an integration test to check that
// PolymerPass and CheckAccessControls work together as expected.
type CheckAccessControlsPolymerTest struct {
	CompilerTestCaseBase
}

// Java: CheckAccessControlsPolymerTest.EXTERNS
var checkAccessControlsPolymerTest_EXTERNS = CompilerTypeTestCase_DEFAULT_EXTERNS +
	`var Polymer = function(descriptor) {};
/** @constructor */
var PolymerElement = function() {};
`

// Java: CheckAccessControlsPolymerTest()
func NewCheckAccessControlsPolymerTest() *CheckAccessControlsPolymerTest {
	c := &CheckAccessControlsPolymerTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, checkAccessControlsPolymerTest_EXTERNS)
	return c
}

func TestCheckAccessControlsPolymerTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCheckAccessControlsPolymerTest() })
}

// Java: CheckAccessControlsPolymerTest.setUp()
func (c *CheckAccessControlsPolymerTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableTypeCheck()
	c.EnableParseTypeInfo()
	c.EnablePolymerPass()
	c.AllowExternsChanges()
	c.EnableCreateModuleMap()
}

// Java: CheckAccessControlsPolymerTest.getProcessor(Compiler)
func (c *CheckAccessControlsPolymerTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewCheckAccessControls(compiler)
}

// Java: CheckAccessControlsPolymerTest.getOptions()
func (c *CheckAccessControlsPolymerTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetWarningLevel(DiagnosticGroups_ACCESS_CONTROLS, CheckLevel_ERROR)
	options.SetWarningLevel(DiagnosticGroups_CONSTANT_PROPERTY, CheckLevel_ERROR)
	return options
}

// Java: CheckAccessControlsPolymerTest.testPrivateMethodInElement()
func (c *CheckAccessControlsPolymerTest) TestPrivateMethodInElement() {
	c.TestNoWarning(
		`var AnElement = Polymer({
  is: 'an-element',

  /** @private */
  foo_: function() {},
  bar: function() { this.foo_(); },
});
`)
}

// Java: CheckAccessControlsPolymerTest.testPrivateMethodInBehavior()
func (c *CheckAccessControlsPolymerTest) TestPrivateMethodInBehavior() {
	c.Test(
		c.Srcs(
			`/** @polymerBehavior */
var Behavior = {
  /** @private */
  foo_: function() {},
  bar: function() { this.foo_(); },
};
`,
			`var AnElement = Polymer({
  is: 'an-element',
  behaviors: [Behavior],
});
`))
}

// Java: CheckAccessControlsPolymerTest.testPrivateMethodFromBehaviorUsedInElement()
func (c *CheckAccessControlsPolymerTest) TestPrivateMethodFromBehaviorUsedInElement() {
	c.TestError(
		c.Srcs(
			`/** @polymerBehavior */
var Behavior = {
  /** @private */
  foo_: function() {},
};
`,
			`var AnElement = Polymer({
  is: 'an-element',
  behaviors: [Behavior],
  bar: function() { this.foo_(); },
});
`),
		CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS)
}

// Java: CheckAccessControlsPolymerTest.testPrivatePropertyInBehavior()
func (c *CheckAccessControlsPolymerTest) TestPrivatePropertyInBehavior() {
	c.Test(
		c.Srcs(
			`/** @polymerBehavior */
var Behavior = {
  /** @private */
  foo_: 'foo',
  bar: function() { alert(this.foo_); },
};
`,
			`var AnElement = Polymer({
  is: 'an-element',
  behaviors: [Behavior],
});
`))
}

// Java: CheckAccessControlsPolymerTest.testPrivatePropertyFromBehaviorUsedInElement()
func (c *CheckAccessControlsPolymerTest) TestPrivatePropertyFromBehaviorUsedInElement() {
	c.TestError(
		c.Srcs(
			`/** @polymerBehavior */
var Behavior = {
  /** @private */
  foo_: 'foo',
};
`,
			`var AnElement = Polymer({
  is: 'an-element',
  behaviors: [Behavior],
  bar: function() { alert(this.foo_); },
});
`),
		CheckAccessControls_BAD_PRIVATE_PROPERTY_ACCESS)
}
