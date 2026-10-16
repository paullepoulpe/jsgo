// Port of com/google/javascript/jscomp/ReplaceTogglesTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/ReplaceTogglesTest.java)

/*
 * Copyright 2023 The Closure Compiler Authors.
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

// ReplaceTogglesTest holds tests for ReplaceToggles.
type ReplaceTogglesTest struct {
	CompilerTestCaseBase
}

// Java: ReplaceTogglesTest()
func NewReplaceTogglesTest() *ReplaceTogglesTest {
	c := &ReplaceTogglesTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestReplaceTogglesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewReplaceTogglesTest() })
}

// Java: ReplaceTogglesTest.testBootstrapOrdinals_notAnObject()
func (c *ReplaceTogglesTest) TestBootstrapOrdinals_notAnObject() {
	c.Test(
		c.Srcs("var CLOSURE_TOGGLE_ORDINALS = 1;"),
		c.Error(ReplaceToggles_INVALID_ORDINAL_MAPPING).
			WithMessageContaining("not an object literal"))
	c.Test(
		c.Srcs("var CLOSURE_TOGGLE_ORDINALS = [];"),
		c.Error(ReplaceToggles_INVALID_ORDINAL_MAPPING).
			WithMessageContaining("not an object literal"))
}

// Java: ReplaceTogglesTest.testBootstrapOrdinals_badKey()
func (c *ReplaceTogglesTest) TestBootstrapOrdinals_badKey() {
	c.Test(
		c.Srcs("var CLOSURE_TOGGLE_ORDINALS = {[x]: 1};"),
		c.Error(ReplaceToggles_INVALID_ORDINAL_MAPPING).WithMessageContaining("non-string key"))
	c.Test(
		c.Srcs("var CLOSURE_TOGGLE_ORDINALS = {['x']: 1};"),
		c.Error(ReplaceToggles_INVALID_ORDINAL_MAPPING).WithMessageContaining("non-string key"))
	c.Test(
		c.Srcs("var CLOSURE_TOGGLE_ORDINALS = {x() { return 1; }};"),
		c.Error(ReplaceToggles_INVALID_ORDINAL_MAPPING).WithMessageContaining("non-string key"))
}

// Java: ReplaceTogglesTest.testBootstrapOrdinals_duplicateKey()
func (c *ReplaceTogglesTest) TestBootstrapOrdinals_duplicateKey() {
	c.Test(
		c.Srcs("var CLOSURE_TOGGLE_ORDINALS = {x: 1, y: 2, x: 3};"),
		c.Error(ReplaceToggles_INVALID_ORDINAL_MAPPING).WithMessageContaining("duplicate key: x"))
}

// Java: ReplaceTogglesTest.testBootstrapOrdinals_badValue()
func (c *ReplaceTogglesTest) TestBootstrapOrdinals_badValue() {
	c.Test(
		c.Srcs("var CLOSURE_TOGGLE_ORDINALS = {x: 'y'};"),
		c.Error(ReplaceToggles_INVALID_ORDINAL_MAPPING).
			WithMessageContaining("value not a boolean or whole number literal"))
	c.Test(
		c.Srcs("var CLOSURE_TOGGLE_ORDINALS = {x};"),
		c.Error(ReplaceToggles_INVALID_ORDINAL_MAPPING).
			WithMessageContaining("value not a boolean or whole number literal"))
	c.Test(
		c.Srcs("var CLOSURE_TOGGLE_ORDINALS = {x: 1 + 2};"),
		c.Error(ReplaceToggles_INVALID_ORDINAL_MAPPING).
			WithMessageContaining("value not a boolean or whole number literal"))
	c.Test(
		c.Srcs("var CLOSURE_TOGGLE_ORDINALS = {x: NaN};"),
		c.Error(ReplaceToggles_INVALID_ORDINAL_MAPPING).
			WithMessageContaining("value not a boolean or whole number literal"))
	c.Test(
		c.Srcs("var CLOSURE_TOGGLE_ORDINALS = {x: Infinity};"),
		c.Error(ReplaceToggles_INVALID_ORDINAL_MAPPING).
			WithMessageContaining("value not a boolean or whole number literal"))
	c.Test(
		c.Srcs("var CLOSURE_TOGGLE_ORDINALS = {x: -1};"),
		c.Error(ReplaceToggles_INVALID_ORDINAL_MAPPING).
			WithMessageContaining("value not a boolean or whole number literal"))
	c.Test(
		c.Srcs("var CLOSURE_TOGGLE_ORDINALS = {x: 1.2};"),
		c.Error(ReplaceToggles_INVALID_ORDINAL_MAPPING).
			WithMessageContaining("value not a boolean or whole number literal"))
	c.Test(
		c.Srcs("var CLOSURE_TOGGLE_ORDINALS = {x: null};"),
		c.Error(ReplaceToggles_INVALID_ORDINAL_MAPPING).
			WithMessageContaining("value not a boolean or whole number literal"))
}

// Java: ReplaceTogglesTest.testBootstrapOrdinals_duplicateOrdinal()
func (c *ReplaceTogglesTest) TestBootstrapOrdinals_duplicateOrdinal() {
	c.Test(
		c.Srcs("var CLOSURE_TOGGLE_ORDINALS = {x: 1, y: 2, z: 1};"),
		c.Error(ReplaceToggles_INVALID_ORDINAL_MAPPING).
			WithMessageContaining("duplicate ordinal: 1"))
}

// Java: ReplaceTogglesTest.testBootstrapOrdinals_twoInitializers()
func (c *ReplaceTogglesTest) TestBootstrapOrdinals_twoInitializers() {
	c.Test(
		c.Srcs(
			`var CLOSURE_TOGGLE_ORDINALS = {x: 1, y: 2};
var CLOSURE_TOGGLE_ORDINALS = {x: 1, y: 2};
`),
		c.Error(ReplaceToggles_INVALID_ORDINAL_MAPPING).
			WithMessageContaining("multiple initialized copies"))
}

// Java: ReplaceTogglesTest.testNoTogglesGiven()
func (c *ReplaceTogglesTest) TestNoTogglesGiven() {
	c.Test(
		"const foo = goog.readToggleInternalDoNotCallDirectly('foo');", //
		"const foo = false;")
}

// Java: ReplaceTogglesTest.testSimpleToggles()
func (c *ReplaceTogglesTest) TestSimpleToggles() {
	c.Test(
		`var CLOSURE_TOGGLE_ORDINALS = {'foo': 1, 'bar': 30, 'baz': 59};
const foo = goog.readToggleInternalDoNotCallDirectly('foo');
const bar = goog.readToggleInternalDoNotCallDirectly('bar');
const baz = goog.readToggleInternalDoNotCallDirectly('baz');
`,
		`var CLOSURE_TOGGLE_ORDINALS = {'foo': 1, 'bar': 30, 'baz': 59};
const foo = !!(goog.TOGGLES_[0] & 2);
const bar = !!(goog.TOGGLES_[1] & 1);
const baz = !!(goog.TOGGLES_[1] >> 29 & 1);
`)
}

// Java: ReplaceTogglesTest.testHardcodedToggles()
func (c *ReplaceTogglesTest) TestHardcodedToggles() {
	c.Test(
		`var CLOSURE_TOGGLE_ORDINALS = {'foo': false, 'bar': true};
const foo = goog.readToggleInternalDoNotCallDirectly('foo');
const bar = goog.readToggleInternalDoNotCallDirectly('bar');
`,
		`var CLOSURE_TOGGLE_ORDINALS = {'foo': false, 'bar': true};
const foo = false;
const bar = true;
`)
}

// Java: ReplaceTogglesTest.testBootstrapOrdinals_ignoresExtraUninitializedDefinitions()
func (c *ReplaceTogglesTest) TestBootstrapOrdinals_ignoresExtraUninitializedDefinitions() {
	c.Test(
		`var CLOSURE_TOGGLE_ORDINALS; // ignored
var CLOSURE_TOGGLE_ORDINALS = {'foo': 1, 'bar': 30, 'baz': 59};
var CLOSURE_TOGGLE_ORDINALS; // ignored
const foo = goog.readToggleInternalDoNotCallDirectly('foo');
const bar = goog.readToggleInternalDoNotCallDirectly('bar');
const baz = goog.readToggleInternalDoNotCallDirectly('baz');
`,
		`var CLOSURE_TOGGLE_ORDINALS;
var CLOSURE_TOGGLE_ORDINALS = {'foo': 1, 'bar': 30, 'baz': 59};
var CLOSURE_TOGGLE_ORDINALS;
const foo = !!(goog.TOGGLES_[0] & 2);
const bar = !!(goog.TOGGLES_[1] & 1);
const baz = !!(goog.TOGGLES_[1] >> 29 & 1);
`)
}

// Java: ReplaceTogglesTest.testBootstrapOrdinals_booleanAllowsDuplicates()
func (c *ReplaceTogglesTest) TestBootstrapOrdinals_booleanAllowsDuplicates() {
	c.Test(
		`var CLOSURE_TOGGLE_ORDINALS =
    {'foo_bar': false, 'baz': false, 'qux': true, 'corge': true};
const fooBar = goog.readToggleInternalDoNotCallDirectly('foo_bar');
const baz = goog.readToggleInternalDoNotCallDirectly('baz');
const qux = goog.readToggleInternalDoNotCallDirectly('qux');
const corge = goog.readToggleInternalDoNotCallDirectly('corge');
`,
		`var CLOSURE_TOGGLE_ORDINALS =
    {'foo_bar': false, 'baz': false, 'qux': true, 'corge': true};
const fooBar = false;
const baz = false;
const qux = true;
const corge = true;
`)
}

// Java: ReplaceTogglesTest.testUnknownToggle()
func (c *ReplaceTogglesTest) TestUnknownToggle() {
	c.Test(
		`var CLOSURE_TOGGLE_ORDINALS = {'foo': 0};
const bar = goog.readToggleInternalDoNotCallDirectly('bar');
`,
		`var CLOSURE_TOGGLE_ORDINALS = {'foo': 0};
const bar = false;
`)
}

// Java: ReplaceTogglesTest.testBadArgument()
func (c *ReplaceTogglesTest) TestBadArgument() {
	c.TestError(
		`const fooString = 'foo';
const foo = goog.readToggleInternalDoNotCallDirectly(fooString);
`,
		ReplaceToggles_INVALID_TOGGLE_PARAMETER)
	c.TestError(
		"const foo = goog.readToggleInternalDoNotCallDirectly('foo', 'bar');", //
		ReplaceToggles_INVALID_TOGGLE_PARAMETER)
	c.TestError(
		"const foo = goog.readToggleInternalDoNotCallDirectly();", //
		ReplaceToggles_INVALID_TOGGLE_PARAMETER)
}

// Java: ReplaceTogglesTest.getProcessor(Compiler)
func (c *ReplaceTogglesTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewReplaceToggles(compiler)
}
