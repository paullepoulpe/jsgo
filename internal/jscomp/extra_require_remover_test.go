// Port of com/google/javascript/jscomp/ExtraRequireRemoverTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/ExtraRequireRemoverTest.java)

/*
 * Copyright 2024 The Closure Compiler Authors.
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

// ExtraRequireRemoverTest holds tests for ExtraRequireRemover.
type ExtraRequireRemoverTest struct {
	CompilerTestCaseBase
}

func NewExtraRequireRemoverTest() *ExtraRequireRemoverTest {
	c := &ExtraRequireRemoverTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestExtraRequireRemoverTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewExtraRequireRemoverTest() })
}

// Java: ExtraRequireRemoverTest.getProcessor(Compiler)
func (c *ExtraRequireRemoverTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewExtraRequireRemover(compiler)
}

// Java: ExtraRequireRemoverTest.getOptions()
func (c *ExtraRequireRemoverTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetWarningLevel(DiagnosticGroups_MODULE_LOAD, CheckLevel_OFF)
	return options
}

// Java: ExtraRequireRemoverTest.testRequireTypeAlias_inGoogModule_unusedAliasedImportsArePruned()
func (c *ExtraRequireRemoverTest) TestRequireTypeAlias_inGoogModule_unusedAliasedImportsArePruned() {
	c.Test(
		`goog.module('x.y.z');
// keep these because they are used
const Foo1 = goog.require('Foo');
const Bar1 = goog.require('Bar');
// remove these because they are unused
const B1 = goog.requireType('B');
const C1 = goog.require('C');
/** @type {!Foo1} */
let foo;
/** @type {!Bar1} */
let bar;
`,
		`goog.module('x.y.z');
const Foo1 = goog.require('Foo');
const Bar1 = goog.require('Bar');
/** @type {!Foo1} */
let foo;
/** @type {!Bar1} */
let bar;
`)
}

// Java: ExtraRequireRemoverTest.testRequireType_inGoogModule_allUnaliasedImportsAreNotPruned()
func (c *ExtraRequireRemoverTest) TestRequireType_inGoogModule_allUnaliasedImportsAreNotPruned() {
	// Unaliased imports in goog.module files may be TypeScript side effects imports (modules that
	// are imported with an import statement that only references the module). We want to preserve
	// these because these imports may be needed, but we won’t be able to find a reference to them
	// in the .i.js file.
	c.TestSame(
		`goog.module('x.y.z');

goog.require('D');

goog.requireType('a.b.c');

/** @const */ var CAlias = a.b.c;
`)
}

// Java: ExtraRequireRemoverTest.testRequireTypeAlias_inGoogModule_partialUnusedDestructuredImportsAreNotPruned()
func (c *ExtraRequireRemoverTest) TestRequireTypeAlias_inGoogModule_partialUnusedDestructuredImportsAreNotPruned() {
	c.Test(
		`goog.module('x.y.z');

// keep A1 which is a used destructured import, remove A2 since it is unused
const {A1, A2} = goog.require('A');

let a = A1();
`,
		`goog.module('x.y.z');

const {A1} = goog.require('A');

let a = A1();
`)
}

// Java: ExtraRequireRemoverTest.testRequireTypeAlias_inGoogModule_unusedDestructuredImportsArePruned()
func (c *ExtraRequireRemoverTest) TestRequireTypeAlias_inGoogModule_unusedDestructuredImportsArePruned() {
	c.Test(
		`goog.module('x.y.z');

const {A1, A2} = goog.require('A');
`,
		"goog.module('x.y.z');")
}

// Java: ExtraRequireRemoverTest.testRequireType_inGoogProvide_usedUnaliasedImportsAreNotPruned()
func (c *ExtraRequireRemoverTest) TestRequireType_inGoogProvide_usedUnaliasedImportsAreNotPruned() {
	c.TestSame(
		`goog.provide('FooAlias');

goog.requireType('a.b.c');

goog.requireType('a.b.c.Foo');

/** @const */ var FooAlias = a.b.c.Foo;
`)
}

// Java: ExtraRequireRemoverTest.testRequireType_inGoogProvide_unusedUnaliasedImportsArePruned()
func (c *ExtraRequireRemoverTest) TestRequireType_inGoogProvide_unusedUnaliasedImportsArePruned() {
	c.Test(
		`goog.provide('FooAlias');

goog.requireType('a.b.c');

goog.requireType('a.b.c.Foo'); // this gets pruned because a.b.c.Foo is not used

/** @const */ var FooAlias = a.b.c;
`,
		`goog.provide('FooAlias');

goog.requireType('a.b.c');

/** @const */ var FooAlias = a.b.c;
`)
}

// Java: ExtraRequireRemoverTest.testRequireType_inGoogProvide_subtypeOfImportUsed()
func (c *ExtraRequireRemoverTest) TestRequireType_inGoogProvide_subtypeOfImportUsed() {
	c.TestSame(
		`goog.provide('FooAlias');

// not pruned because subtype a.b.c.Bar is used
goog.requireType('a.b.c');

/** @type {a.b.c.Bar} */
let bar;
`)
}

// Java: ExtraRequireRemoverTest.testForwardDeclare_inGoogModule_unusedUnaliasedImportsAreNotPruned()
func (c *ExtraRequireRemoverTest) TestForwardDeclare_inGoogModule_unusedUnaliasedImportsAreNotPruned() {
	c.TestSame(
		`goog.module('x.y.z');
// keep all imports because because aliased goog.forwardDeclare's will not be pruned

goog.forwardDeclare('a.b.c');

goog.forwardDeclare('a.b.c.Foo');

/** @const */ var FooAlias = a.b.c;
`)
}

// Java: ExtraRequireRemoverTest.testForwardDeclare_inGoogModule_unusedAliasedImportsAreNotPruned()
func (c *ExtraRequireRemoverTest) TestForwardDeclare_inGoogModule_unusedAliasedImportsAreNotPruned() {
	c.TestSame(
		`goog.module('x.y.z');
// keep all imports because because aliased goog.forwardDeclare's will not be pruned

const foo = goog.forwardDeclare('foo');

const bar = goog.forwardDeclare('bar');

/** @const */ var FooAlias = foo;
`)
}

// Java: ExtraRequireRemoverTest.testForwardDeclare_inGoogProvide_unusedUnliasedImportsArePruned()
func (c *ExtraRequireRemoverTest) TestForwardDeclare_inGoogProvide_unusedUnliasedImportsArePruned() {
	c.Test(
		`goog.provide('FooAlias');

goog.forwardDeclare('a.b.c');

// this gets pruned because bar is not used
goog.forwardDeclare('bar');

/** @const */ var FooAlias = a.b.c;
`,
		`goog.provide('FooAlias');

goog.forwardDeclare('a.b.c');

/** @const */ var FooAlias = a.b.c;
`)
}
