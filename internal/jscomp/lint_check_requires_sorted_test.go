// Port of com/google/javascript/jscomp/lint/CheckRequiresSortedTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/lint/CheckRequiresSortedTest.java)

/*
 * Copyright 2019 The Closure Compiler Authors.
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

	"github.com/paullepoulpe/jsgo/internal/parsing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Tests for {@link CheckRequiresSorted}.
//
// <p>Note that the sort/deduplication logic is extensively tested by {@link ErrorToFixMapperTest}.
// These tests are only for asserting that warnings are reported correctly.
type CheckRequiresSortedTest struct {
	CompilerTestCaseBase
}

func NewCheckRequiresSortedTest() *CheckRequiresSortedTest {
	c := &CheckRequiresSortedTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestCheckRequiresSortedTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCheckRequiresSortedTest() })
}

// Java: CheckRequiresSortedTest.getProcessor(Compiler)
func (c *CheckRequiresSortedTest) GetProcessor(compiler *Compiler) CompilerPass {
	callback := NewCheckRequiresSorted(CheckRequiresSorted_Mode_COLLECT_AND_REPORT)
	return &checkRequiresSortedTest_processor{compiler: compiler, callback: callback}
}

// checkRequiresSortedTest_processor is the anonymous CompilerPass of getProcessor.
type checkRequiresSortedTest_processor struct {
	compiler *Compiler
	callback *CheckRequiresSorted
}

// Java: CheckRequiresSortedTest.getProcessor(...).new CompilerPass() {...}.process(Node, Node)
func (p *checkRequiresSortedTest_processor) Process(externs *rhino.Node, root *rhino.Node) {
	NodeTraversal_Traverse(p.compiler, root, p.callback)
}

// Java: CheckRequiresSortedTest.getOptions()
func (c *CheckRequiresSortedTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetParseJsDocDocumentation(parsing.Config_JsDocParsing_INCLUDE_DESCRIPTIONS_WITH_WHITESPACE)
	options.SetPreserveDetailedSourceInfo(true)
	options.SetPrettyPrint(true)
	options.SetPreserveTypeAnnotations(true)
	options.SetPreferSingleQuotes(true)
	options.SetEmitUseStrict(false)
	return options
}

// Java: CheckRequiresSortedTest.testNoWarning()
//
// Go-only: test-helper — a trailing underscore: the plain name is the
// inherited CompilerTestCase.testNoWarning(...), which this no-argument @Test method
// overloads in Java.
func (c *CheckRequiresSortedTest) TestNoWarning_() {
	c.TestNoWarning(
		`goog.module('x');

const c = goog.require('c');
const {b} = goog.require('b');
goog.require('a');
`)
}

// Java: CheckRequiresSortedTest.testNoWarning_noRequires()
func (c *CheckRequiresSortedTest) TestNoWarning_noRequires() {
	c.TestNoWarning(
		`goog.module('x');

alert(1);
`)
}

// Java: CheckRequiresSortedTest.testWarning()
//
// Go-only: test-helper — a trailing underscore: the plain name is the
// inherited CompilerTestCase.testWarning(...), which this no-argument @Test method
// overloads in Java.
func (c *CheckRequiresSortedTest) TestWarning_() {
	c.Test(
		c.Srcs(
			`goog.module('x');

const b = goog.require('b');
const a = goog.require('a')

alert(1);
`),
		c.Warning(CheckRequiresSorted_REQUIRES_NOT_SORTED).WithMessageContaining(
			`The correct order is:

const a = goog.require('a');
const b = goog.require('b');
`))
}

// Java: CheckRequiresSortedTest.testWarning_withJsDoc()
func (c *CheckRequiresSortedTest) TestWarning_withJsDoc() {
	c.Test(
		c.Srcs(
			`goog.module('x');

/**
 * @suppress {extraRequire}
 */
goog.require('b');
goog.require('a');

alert(1);
`),
		c.Warning(CheckRequiresSorted_REQUIRES_NOT_SORTED).WithMessageContaining(
			`The correct order is:

goog.require('a');
/**
 * @suppress {extraRequire}
 */
goog.require('b');
`))
}

// Java: CheckRequiresSortedTest.testWarning_destructuringWithNoShorthandProperties()
func (c *CheckRequiresSortedTest) TestWarning_destructuringWithNoShorthandProperties() {
	c.Test(
		c.Srcs(
			`goog.module('x');

const {a: a} = goog.require('a');

alert(1);
`),
		c.Warning(CheckRequiresSorted_REQUIRES_NOT_SORTED).WithMessageContaining(
			`The correct order is:

const {a} = goog.require('a');
`))
}
