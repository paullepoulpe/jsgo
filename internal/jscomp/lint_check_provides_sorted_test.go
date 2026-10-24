// Port of com/google/javascript/jscomp/lint/CheckProvidesSortedTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/lint/CheckProvidesSortedTest.java)

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

// Tests for {@link CheckProvidesSorted}.
//
// <p>Note that the sort/deduplication logic is extensively tested by {@link ErrorToFixMapperTest}.
// These tests are only for asserting that warnings are reported correctly.
type CheckProvidesSortedTest struct {
	CompilerTestCaseBase
}

func NewCheckProvidesSortedTest() *CheckProvidesSortedTest {
	c := &CheckProvidesSortedTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestCheckProvidesSortedTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCheckProvidesSortedTest() })
}

// Java: CheckProvidesSortedTest.getProcessor(Compiler)
func (c *CheckProvidesSortedTest) GetProcessor(compiler *Compiler) CompilerPass {
	callback := NewCheckProvidesSorted(CheckProvidesSorted_Mode_COLLECT_AND_REPORT)
	return &checkProvidesSortedTest_processor{compiler: compiler, callback: callback}
}

// checkProvidesSortedTest_processor is the anonymous CompilerPass of getProcessor.
type checkProvidesSortedTest_processor struct {
	compiler *Compiler
	callback *CheckProvidesSorted
}

// Java: CheckProvidesSortedTest.getProcessor(...).new CompilerPass() {...}.process(Node, Node)
func (p *checkProvidesSortedTest_processor) Process(externs *rhino.Node, root *rhino.Node) {
	NodeTraversal_Traverse(p.compiler, root, p.callback)
}

// Java: CheckProvidesSortedTest.getOptions()
func (c *CheckProvidesSortedTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetParseJsDocDocumentation(parsing.Config_JsDocParsing_INCLUDE_DESCRIPTIONS_WITH_WHITESPACE)
	options.SetPreserveDetailedSourceInfo(true)
	options.SetPrettyPrint(true)
	options.SetPreserveTypeAnnotations(true)
	options.SetPreferSingleQuotes(true)
	options.SetEmitUseStrict(false)
	return options
}

// Java: CheckProvidesSortedTest.testNoWarning()
//
// Go-only: test-helper — a trailing underscore: the plain name is the
// inherited CompilerTestCase.testNoWarning(...), which this no-argument @Test method
// overloads in Java.
func (c *CheckProvidesSortedTest) TestNoWarning_() {
	c.TestNoWarning(
		`/** @fileoverview foo */

goog.provide('a');
goog.provide('b');
goog.provide('c');

alert(1);
`)
}

// Java: CheckProvidesSortedTest.testNoWarning_noProvides()
func (c *CheckProvidesSortedTest) TestNoWarning_noProvides() {
	c.TestNoWarning(
		`/** @fileoverview foo */

goog.module('m');

goog.require('x');
`)
}

// Java: CheckProvidesSortedTest.testWarning()
//
// Go-only: test-helper — a trailing underscore: the plain name is the
// inherited CompilerTestCase.testWarning(...), which this no-argument @Test method
// overloads in Java.
func (c *CheckProvidesSortedTest) TestWarning_() {
	c.Test(
		c.Srcs(
			`/** @fileoverview foo */

goog.provide('b');
goog.provide('a');
goog.provide('c');

alert(1);
`),
		c.Warning(CheckProvidesSorted_PROVIDES_NOT_SORTED).WithMessageContaining(
			`The correct order is:

goog.provide('a');
goog.provide('b');
goog.provide('c');
`))
}
