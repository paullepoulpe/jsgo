// Port of com/google/javascript/jscomp/RemoveUnnecessarySyntheticExternsTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/RemoveUnnecessarySyntheticExternsTest.java)

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"testing"
)

type RemoveUnnecessarySyntheticExternsTest struct {
	CompilerTestCaseBase

	// use this set to simulate an earlier compiler pass declaring a synthetic extern
	syntheticExternsToAdd *jutil.LinkedHashSet[rhino.Node]
}

// Java: RemoveUnnecessarySyntheticExternsTest()
func NewRemoveUnnecessarySyntheticExternsTest() *RemoveUnnecessarySyntheticExternsTest {
	c := &RemoveUnnecessarySyntheticExternsTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestRemoveUnnecessarySyntheticExternsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewRemoveUnnecessarySyntheticExternsTest() })
}

// Java: RemoveUnnecessarySyntheticExternsTest.setUp()
func (c *RemoveUnnecessarySyntheticExternsTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableMultistageCompilation()
	c.syntheticExternsToAdd = jutil.NewLinkedHashSet[rhino.Node]()
}

// Java: RemoveUnnecessarySyntheticExternsTest.getProcessor(Compiler)
func (c *RemoveUnnecessarySyntheticExternsTest) GetProcessor(compiler *Compiler) CompilerPass {
	return CompilerPassFunc(func(externs rhino.Node, js rhino.Node) {
		syntheticRoot := compiler.GetSynthesizedExternsInput().GetAstRoot(compiler)
		for declaration := range c.syntheticExternsToAdd.All() {
			syntheticRoot.AddChildToBack(declaration.CloneTree())
		}
		NewRemoveUnnecessarySyntheticExterns(compiler).Process(externs, js)
	})
}

// Java: RemoveUnnecessarySyntheticExternsTest.doesntChangeSyntheticExternThatIsNotDeclaredInCode()
func (c *RemoveUnnecessarySyntheticExternsTest) TestDoesntChangeSyntheticExternThatIsNotDeclaredInCode() {
	c.syntheticExternsToAdd.Add(c.CreateUnfulfilledDeclaration("x"))
	c.syntheticExternsToAdd.Add(c.CreateUnfulfilledDeclaration("y"))

	c.TestExternChanges(c.Srcs("x;"), c.Expected("var x; var y;"))
}

// Java: RemoveUnnecessarySyntheticExternsTest.doesntRemoveSyntheticExternIfShadowedInCode()
func (c *RemoveUnnecessarySyntheticExternsTest) TestDoesntRemoveSyntheticExternIfShadowedInCode() {
	c.syntheticExternsToAdd.Add(c.CreateUnfulfilledDeclaration("x"))

	c.TestExternChanges(c.Srcs("function fn() { var x; }"), c.Expected("var x;"))
	c.TestExternChanges(c.Srcs("{ let x; }"), c.Expected("var x;"))
	c.TestExternChanges(c.Srcs("{ const x = 0; }"), c.Expected("var x;"))
	c.TestExternChanges(c.Srcs("{ class x {} }"), c.Expected("var x;"))
}

// Java: RemoveUnnecessarySyntheticExternsTest.removesSyntheticExternDeclaredInCode()
func (c *RemoveUnnecessarySyntheticExternsTest) TestRemovesSyntheticExternDeclaredInCode() {
	c.syntheticExternsToAdd.Add(c.CreateUnfulfilledDeclaration("x"))
	c.syntheticExternsToAdd.Add(c.CreateUnfulfilledDeclaration("y"))

	c.TestExternChanges(c.Srcs("var x;"), c.Expected("var y;"))
	c.TestExternChanges(c.Srcs("{ var x; }"), c.Expected("var y;"))
	c.TestExternChanges(c.Srcs("let x;"), c.Expected("var y;"))
	c.TestExternChanges(c.Srcs("const x = 0;"), c.Expected("var y;"))
	c.TestExternChanges(c.Srcs("function x() {}"), c.Expected("var y;"))
	c.TestExternChanges(c.Srcs("class x {}"), c.Expected("var y;"))
}

// Java: RemoveUnnecessarySyntheticExternsTest.removesSyntheticExternDeclaredInGoogProvide()
func (c *RemoveUnnecessarySyntheticExternsTest) TestRemovesSyntheticExternDeclaredInGoogProvide() {
	c.syntheticExternsToAdd.Add(c.CreateUnfulfilledDeclaration("x"))
	c.syntheticExternsToAdd.Add(c.CreateUnfulfilledDeclaration("y"))

	c.TestExternChanges(c.Srcs("goog.provide('x');"), c.Expected("var y;"))
	c.TestExternChanges(c.Srcs("goog.provide('x.y.z');"), c.Expected("var y;"))
	c.TestExternChanges(c.Srcs("goog.provide('other'); goog.provide('x.y');"), c.Expected("var y;"))
}

// Java: RemoveUnnecessarySyntheticExternsTest.doesNotRemovesSyntheticExternForLegacyGoogModule()
func (c *RemoveUnnecessarySyntheticExternsTest) TestDoesNotRemovesSyntheticExternForLegacyGoogModule() {
	c.syntheticExternsToAdd.Add(c.CreateUnfulfilledDeclaration("x"))
	c.syntheticExternsToAdd.Add(c.CreateUnfulfilledDeclaration("y"))

	c.TestExternChanges(c.Srcs("goog.module('x');"), c.Expected("var x; var y;"))
	c.TestExternChanges(c.Srcs("goog.module('x.y');"), c.Expected("var x; var y;"))
}

// Java: RemoveUnnecessarySyntheticExternsTest.removesSyntheticExternDeclaredInLegacyGoogModuleNamespace()
func (c *RemoveUnnecessarySyntheticExternsTest) TestRemovesSyntheticExternDeclaredInLegacyGoogModuleNamespace() {
	c.syntheticExternsToAdd.Add(c.CreateUnfulfilledDeclaration("x"))
	c.syntheticExternsToAdd.Add(c.CreateUnfulfilledDeclaration("y"))

	c.TestExternChanges(
		c.Srcs("goog.module('x'); goog.module.declareLegacyNamespace();"), c.Expected("var y;"))
	c.TestExternChanges(
		c.Srcs("goog.module('x.y'); goog.module.declareLegacyNamespace();"), c.Expected("var y;"))
}

// Java: RemoveUnnecessarySyntheticExternsTest.removesSyntheticExternDeclaredInOtherExterns()
func (c *RemoveUnnecessarySyntheticExternsTest) TestRemovesSyntheticExternDeclaredInOtherExterns() {
	c.syntheticExternsToAdd.Add(c.CreateUnfulfilledDeclaration("x"))
	c.syntheticExternsToAdd.Add(c.CreateUnfulfilledDeclaration("y"))

	c.TestExternChanges(c.Externs("var x;"), c.Srcs(""), c.Expected("var y;", "var x;"))
	c.TestExternChanges(c.Externs("var x;"), c.Srcs("var y;"), c.Expected("var x;"))
}

// Java: RemoveUnnecessarySyntheticExternsTest.removesDistinctSyntheticExternsDeclaredInCode()
func (c *RemoveUnnecessarySyntheticExternsTest) TestRemovesDistinctSyntheticExternsDeclaredInCode() {
	c.syntheticExternsToAdd.Add(c.CreateUnfulfilledDeclaration("x"))
	c.syntheticExternsToAdd.Add(c.CreateUnfulfilledDeclaration("y"))

	c.TestExternChanges(c.Srcs("var x; var y;"), c.Expected(""))
}

// Java: RemoveUnnecessarySyntheticExternsTest.removesDuplicateUnfulfilledSyntheticDeclarations()
func (c *RemoveUnnecessarySyntheticExternsTest) TestRemovesDuplicateUnfulfilledSyntheticDeclarations() {
	c.syntheticExternsToAdd.Add(c.CreateUnfulfilledDeclaration("x"))
	c.syntheticExternsToAdd.Add(c.CreateUnfulfilledDeclaration("y"))
	c.syntheticExternsToAdd.Add(c.CreateUnfulfilledDeclaration("x"))
	c.syntheticExternsToAdd.Add(c.CreateUnfulfilledDeclaration("x"))

	c.TestExternChanges(c.Srcs(""), c.Expected("var x; var y;"))
	c.TestExternChanges(c.Srcs("var x;"), c.Expected("var y;"))
}

// Java: RemoveUnnecessarySyntheticExternsTest.removesSyntheticExternDeclaredInCode_multipleCodeDeclarations()
func (c *RemoveUnnecessarySyntheticExternsTest) TestRemovesSyntheticExternDeclaredInCode_multipleCodeDeclarations() {
	c.syntheticExternsToAdd.Add(c.CreateUnfulfilledDeclaration("x"))
	c.syntheticExternsToAdd.Add(c.CreateUnfulfilledDeclaration("y"))

	c.TestExternChanges(c.Srcs("var x; function x() {}"), c.Expected("var y;"))
}

// Java: RemoveUnnecessarySyntheticExternsTest.doesntChangeSyntheticExtern_declaredInCodeNotMarkedUnfulfilled()
func (c *RemoveUnnecessarySyntheticExternsTest) TestDoesntChangeSyntheticExtern_declaredInCodeNotMarkedUnfulfilled() {
	// add an extern that is not marked as an 'unfulfilled declaration'
	// we assume this extern was added to prevent renaming, not just enforce that all referenced
	// names are declared.
	c.syntheticExternsToAdd.Add(rhino.IR_Var(rhino.IR_Name("x")))

	c.TestExternChanges(c.Srcs("x;"), c.Expected("var x;"))
}

// Java: RemoveUnnecessarySyntheticExternsTest.removesOnlyUnfulfilledSyntheticExterns_ifMixOfFulfilledAndUnfulfilled()
func (c *RemoveUnnecessarySyntheticExternsTest) TestRemovesOnlyUnfulfilledSyntheticExterns_ifMixOfFulfilledAndUnfulfilled() {
	// add an extern that is not marked as an 'unfulfilled declaration'
	// we assume this extern was added to actually prevent renaming, and so must not be removed
	// even if a duplicate of a non-synthetic extern.
	c.syntheticExternsToAdd.Add(rhino.IR_Var(rhino.IR_Name("x")))
	// add a duplicate declaration of 'x', where the second and third only are unfulfilled.
	c.syntheticExternsToAdd.Add(c.CreateUnfulfilledDeclaration("x"))
	c.syntheticExternsToAdd.Add(c.CreateUnfulfilledDeclaration("x"))

	c.TestExternChanges(c.Srcs(""), c.Expected("var x;"))
	c.TestExternChanges(c.Srcs("var x;"), c.Expected("var x;"))
}

// Java: RemoveUnnecessarySyntheticExternsTest.createUnfulfilledDeclaration(String)
func (c *RemoveUnnecessarySyntheticExternsTest) CreateUnfulfilledDeclaration(name string) rhino.Node {
	// only VAR nodes are allowed to be marked "synthesized unfulfilled" so no need to test other
	// kinds of declarations.
	declaration := rhino.IR_Var(rhino.IR_Name(name))
	declaration.SetIsSynthesizedUnfulfilledNameDeclaration(true)
	return declaration
}
