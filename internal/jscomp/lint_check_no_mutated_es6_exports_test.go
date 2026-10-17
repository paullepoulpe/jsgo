// Port of com/google/javascript/jscomp/lint/CheckNoMutatedEs6ExportsTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/lint/CheckNoMutatedEs6ExportsTest.java)

/*
 * Copyright 2018 The Closure Compiler Authors.
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

type CheckNoMutatedEs6ExportsTest struct {
	CompilerTestCaseBase
}

func NewCheckNoMutatedEs6ExportsTest() *CheckNoMutatedEs6ExportsTest {
	c := &CheckNoMutatedEs6ExportsTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestCheckNoMutatedEs6ExportsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCheckNoMutatedEs6ExportsTest() })
}

// Java: CheckNoMutatedEs6ExportsTest.getProcessor(Compiler)
func (c *CheckNoMutatedEs6ExportsTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewCheckNoMutatedEs6Exports(compiler)
}

// Java: CheckNoMutatedEs6ExportsTest.testNeverMutatedExportIsOk()
func (c *CheckNoMutatedEs6ExportsTest) TestNeverMutatedExportIsOk() {
	c.TestSame("export let x = 0;")
}

// Java: CheckNoMutatedEs6ExportsTest.testLocalWithExportedNameCanBeMutated()
func (c *CheckNoMutatedEs6ExportsTest) TestLocalWithExportedNameCanBeMutated() {
	c.TestSame("export let x = 0; () => { let x = 0; x++; }")
}

// Java: CheckNoMutatedEs6ExportsTest.testMutatedDuringInitializationIsOk()
func (c *CheckNoMutatedEs6ExportsTest) TestMutatedDuringInitializationIsOk() {
	c.TestSame("export let x = 0; x++;")
	c.TestSame("let x = 0; export {x as y}; x++;")
}

// Java: CheckNoMutatedEs6ExportsTest.testMutatedInInnerScopeIsError()
func (c *CheckNoMutatedEs6ExportsTest) TestMutatedInInnerScopeIsError() {
	c.TestWarning("export let x = 0; () => x++;", CheckNoMutatedEs6Exports_MUTATED_EXPORT)
	c.TestWarning("let x = 0; export {x as y}; () => x++;", CheckNoMutatedEs6Exports_MUTATED_EXPORT)
	c.TestWarning("export function foo() {}; () => foo = 0;", CheckNoMutatedEs6Exports_MUTATED_EXPORT)
	c.TestWarning("export default function foo() {}; () => foo = 0;", CheckNoMutatedEs6Exports_MUTATED_EXPORT)
	c.TestWarning("export default class Foo {}; () => Foo = 0;", CheckNoMutatedEs6Exports_MUTATED_EXPORT)
	c.TestWarning("export function foo() { foo = 0; };", CheckNoMutatedEs6Exports_MUTATED_EXPORT)
	c.TestWarning("export let x = 0; export function foo() { x++; };", CheckNoMutatedEs6Exports_MUTATED_EXPORT)
}

// Java: CheckNoMutatedEs6ExportsTest.testExportDefaultExpressionIsOk()
func (c *CheckNoMutatedEs6ExportsTest) TestExportDefaultExpressionIsOk() {
	c.TestSame("export default 0;")
	c.TestSame("export default {prop: 0}")
	c.TestSame("export default function() {};")
	c.TestSame("export default class {};")
}
