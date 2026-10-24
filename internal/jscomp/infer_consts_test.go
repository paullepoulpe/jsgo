// Port of com/google/javascript/jscomp/InferConstsTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/InferConstsTest.java)

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

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// InferConstsTest: tests for InferConsts.
type InferConstsTest struct {
	CompilerTestCaseBase

	constFinder *InferConstsTest_FindConstants

	names []string
}

// Java: InferConstsTest()
func NewInferConstsTest() *InferConstsTest {
	c := &InferConstsTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestInferConstsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewInferConstsTest() })
}

// Java: InferConstsTest.getProcessor(Compiler)
func (c *InferConstsTest) GetProcessor(compiler *Compiler) CompilerPass {
	c.constFinder = NewInferConstsTest_FindConstants(c.names)
	return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
		NewInferConsts(compiler).Process(externs, root)
		NodeTraversal_Traverse(compiler, root, c.constFinder)
	})
}

// Java: InferConstsTest.testSimple()
func (c *InferConstsTest) TestSimple() {
	c.AssertConsts("var x = 3;", InferConstsTest_Inferred("x"))
	c.AssertConsts("/** @const */ var x;", InferConstsTest_Declared("x"), InferConstsTest_NotInferred("x"))
	c.AssertConsts("var x = 3, y = 4;", InferConstsTest_Inferred("x", "y"))
	c.AssertConsts("var x = 3, y;", InferConstsTest_Inferred("x"), InferConstsTest_NotInferred("y"))
	c.AssertConsts("var x = 3;  function f(){x;}", InferConstsTest_Inferred("x"))
	c.AssertConsts("var x = 3, y; y ||= 4;", InferConstsTest_Inferred("x"), InferConstsTest_NotInferred("y"))
}

// Java: InferConstsTest.testSimpleLetConst()
func (c *InferConstsTest) TestSimpleLetConst() {
	c.AssertConsts("let x = 3, y", InferConstsTest_Inferred("x"))
	c.AssertConsts("let x;", InferConstsTest_NotInferred("x"))
	c.AssertConsts("let x; x = 0;", InferConstsTest_NotInferred("x"))
	c.AssertConsts("let x = 3; let y = 4;", InferConstsTest_Inferred("x", "y"))
	c.AssertConsts("let x = 3, y = 4; x++;", InferConstsTest_Inferred("y"))
	c.AssertConsts("let x = 3;  function f(){let x = 4;}", InferConstsTest_Inferred("x"))
	c.AssertConsts("let x; x ||= 0;", InferConstsTest_NotInferred("x"))
	c.AssertConsts("/** @const */ let x;", InferConstsTest_Declared("x"), InferConstsTest_NotInferred("x"))
	c.AssertConsts("const x = 1;", InferConstsTest_Declared("x"), InferConstsTest_Inferred("x"))
}

// Java: InferConstsTest.testPossibleEarlyRefInFunction()
func (c *InferConstsTest) TestPossibleEarlyRefInFunction() {
	c.AssertConsts("function f() { x; } var x = 0;", InferConstsTest_NotInferred("x"))
	c.AssertConsts("function f() { x; } let x = 0;", InferConstsTest_Inferred("x"))
	c.AssertConsts("function f() { x; } let [x] = [];", InferConstsTest_Inferred("x"))
	c.AssertConsts("function f() { x; } const x = 0;", InferConstsTest_Declared("x"), InferConstsTest_Inferred("x"))
	c.AssertConsts("function f() { Foo; } class Foo {}", InferConstsTest_Inferred("Foo"))
	c.AssertConsts("function f() { g; } function g() {}", InferConstsTest_Inferred("g"))
	c.AssertConsts("function f([y = () => x], x) {}", InferConstsTest_Inferred("x"))
	c.AssertConsts("function f() { x ||= 1; } let x = 0;", InferConstsTest_NotInferred("x"))
}

// Java: InferConstsTest.testUnfound()
func (c *InferConstsTest) TestUnfound() {
	c.AssertNotConsts("var x = 2; x++;", "x")
	c.AssertNotConsts("var x = 2; x = 3;", "x")
	c.AssertNotConsts("var x = 3;  function f(){x++;}", "x")
	c.AssertNotConsts("let x = 3; x++;", "x")
	c.AssertNotConsts("let x = 3; x = 2;", "x", "y")
	c.AssertNotConsts("/** @const */let x; let y;", "y")
	c.AssertNotConsts("let x = 3;  function f() {let x = 4; x++;} x++;", "x")
	c.AssertNotConsts("var x = 2; x ||= 1;", "x")
	c.AssertNotConsts("let x = 3; x ||= 1;", "x")
}

// Java: InferConstsTest.testForOf()
func (c *InferConstsTest) TestForOf() {
	c.AssertNotConsts("var x = 0; for (x of [1, 2, 3]) {}", "x")
	c.AssertNotConsts("var x = 0; for (x of {a, b, c}) {}", "x")
}

// Java: InferConstsTest.testForIn()
func (c *InferConstsTest) TestForIn() {
	c.AssertNotConsts("var x = 0; for (x in {a, b}) {}", "x")
}

// Java: InferConstsTest.testForVar()
func (c *InferConstsTest) TestForVar() {
	c.AssertNotConsts("for (var x = 0; x < 2; x++) {}", "x")
	c.AssertNotConsts("for (var x in [1, 2, 3]) {}", "x")
	c.AssertNotConsts("for (var x of {a, b, c}) {}", "x")
}

// Java: InferConstsTest.testForLet()
func (c *InferConstsTest) TestForLet() {
	c.AssertNotConsts("for (let x = 0; x < 2; x++) {}", "x")
	c.AssertNotConsts("for (let x in [1, 2, 3]) {}", "x")
	c.AssertNotConsts("for (let x of {a, b, c}) {}", "x")
}

// Java: InferConstsTest.testForConst()
func (c *InferConstsTest) TestForConst() {
	// Using 'const' here is not allowed, and ConstCheck will warn for this
	c.AssertConsts("for (const x = 0; x < 2; x++) {}", InferConstsTest_NotInferred("x"), InferConstsTest_Declared("x"))
}

// Java: InferConstsTest.testForConst1()
func (c *InferConstsTest) TestForConst1() {
	c.AssertConsts("for (const x in [1, 2, 3]) {}", InferConstsTest_Declared("x"), InferConstsTest_NotInferred("x"))
	c.AssertConsts("for (const x of {a, b, c}) {}", InferConstsTest_Declared("x"), InferConstsTest_NotInferred("x"))
}

// Java: InferConstsTest.testForConstJSDoc()
func (c *InferConstsTest) TestForConstJSDoc() {
	c.AssertConsts("for (/** @const */ let x = 0; x < 2; x++) {}", InferConstsTest_Declared("x"), InferConstsTest_NotInferred("x"))
	c.AssertConsts("for (/** @const */ let x in [1, 2, 3]) {}", InferConstsTest_Declared("x"), InferConstsTest_NotInferred("x"))
	c.AssertConsts("for (/** @const */ let x of {a, b, c}) {}", InferConstsTest_Declared("x"), InferConstsTest_NotInferred("x"))
}

// Java: InferConstsTest.testFunctionParam()
func (c *InferConstsTest) TestFunctionParam() {
	c.AssertConsts("var x = function(){};", InferConstsTest_Inferred("x"))
	c.AssertConsts("var x = ()=>{};", InferConstsTest_Inferred("x"))
	c.AssertConsts("const x = ()=>{};", InferConstsTest_Inferred("x"))
	c.AssertConsts("/** @const */ let x = ()=>{};", InferConstsTest_Inferred("x"), InferConstsTest_Declared("x"))
	c.AssertConsts("function fn(a){var b = a + 1}; ", InferConstsTest_Inferred("a", "b"))
	c.AssertConsts("function fn(a = 1){var b = a + 1}; ", InferConstsTest_Inferred("a", "b"))
	c.AssertConsts("function fn(a, {b, c}){var d = a + 1}; ", InferConstsTest_Inferred("a", "b", "c", "d"))
}

// Java: InferConstsTest.testClass()
func (c *InferConstsTest) TestClass() {
	c.AssertConsts("var Foo = class {}", InferConstsTest_Inferred("Foo"))
	c.AssertConsts("const Foo = class {}", InferConstsTest_Declared("Foo"), InferConstsTest_Inferred("Foo"))
	c.AssertConsts("class Foo {}", InferConstsTest_NotDeclared("Foo"), InferConstsTest_Inferred("Foo"))
	c.AssertConsts("var Foo = function() {};", InferConstsTest_NotDeclared("Foo"), InferConstsTest_Inferred("Foo"))
	c.AssertConsts("const Foo = function() {}", InferConstsTest_Declared("Foo"), InferConstsTest_Inferred("Foo"))
	c.AssertConsts("function Foo() {}", InferConstsTest_NotDeclared("Foo"), InferConstsTest_Inferred("Foo"))
}

// Java: InferConstsTest.testVarArguments()
func (c *InferConstsTest) TestVarArguments() {
	c.AssertConsts("var arguments = 3;", InferConstsTest_Inferred("arguments"))
}

// Java: InferConstsTest.testConstArguments()
func (c *InferConstsTest) TestConstArguments() {
	c.AssertConsts("const arguments = 4;", InferConstsTest_Declared("arguments"), InferConstsTest_Inferred("arguments"))
}

// Java: InferConstsTest.testArgumentsJSDoc()
func (c *InferConstsTest) TestArgumentsJSDoc() {
	c.AssertConsts("/** @const */let arguments = 5;", InferConstsTest_Declared("arguments"), InferConstsTest_Inferred("arguments"))
}

// Java: InferConstsTest.testDestructuring()
func (c *InferConstsTest) TestDestructuring() {
	c.AssertConsts("var [a, b, c] = [1, 2, 3];", InferConstsTest_Inferred("a", "b", "c"))
	c.AssertConsts("const [a, b, c] = [1, 2, 3];", InferConstsTest_Inferred("a", "b", "c"))
	c.AssertNotConsts("var [a, b, c] = obj;", "obj")
	c.AssertNotConsts(
		"var [a, b, c] = [1, 2, 3];\n"+
			"[a, b, c] = [1, 2, 3];\n",
		"a",
		"b",
		"c")
	c.AssertConsts(
		"var [a, b, c] = [1, 2, 3];\n"+
			"[a, b]= [1, 2];\n",
		InferConstsTest_Inferred("c"))

	c.AssertNotConsts("var [a, b, c] = [1, 2, 3]; [a, b] = [1, 2];", "a", "b")
	c.AssertConsts("var [a, b, c] = [1, 2, 3]; [a, b] = [1, 2];", InferConstsTest_Inferred("c"))

	c.AssertConsts("var {a: b} = {a: 1}", InferConstsTest_Inferred("b"))
	c.AssertNotConsts("var {a: b} = {a: 1}", "a")
	// Note that this "a" looks for the destructured "a"
	c.AssertConsts("var obj = {a: 1}; var {a} = obj", InferConstsTest_Inferred("a"))

	c.AssertNotConsts(
		"var [{a: x} = {a: 'y'}] = [{a: 'x'}];\n"+
			"[{a: x} = {a: 'x'}] = {};\n",
		"x")
	c.AssertNotConsts(
		"let fg = '', bg = '';\n"+
			"({fg, bg} = pal[val - 1]);\n",
		"fg",
		"bg")

	c.AssertConsts("var [a, , b] = [1, 2, 3];", InferConstsTest_Inferred("a", "b"))
	c.AssertConsts("const [a, , b] = [1, 2, 3];", InferConstsTest_Inferred("a", "b"))
}

// Java: InferConstsTest.testDefaultValue()
func (c *InferConstsTest) TestDefaultValue() {
	c.AssertConsts("function fn(a = 1){}", InferConstsTest_Inferred("a"))
	c.AssertNotConsts("function fn(a = 1){a = 2}", "a")

	c.AssertConsts("function fn({b, c} = {b:1, c:2}){}", InferConstsTest_Inferred("b", "c"))
	c.AssertNotConsts("function fn({b, c} = {b:1, c:2}){c = 1}", "c")
}

// Java: InferConstsTest.testVarInBlock()
func (c *InferConstsTest) TestVarInBlock() {
	c.AssertConsts(
		"function f() { if (true) { var x = function() {}; x(); } }",
		InferConstsTest_NotDeclared("x"),
		InferConstsTest_Inferred("x"))
}

// Java: InferConstsTest.testGeneratorFunctionVar()
func (c *InferConstsTest) TestGeneratorFunctionVar() {
	c.AssertNotConsts(
		"function *gen() {\n"+
			"  var x = 0;\n"+
			"  while (x < 3)\n"+
			"    yield x++;\n"+
			"}\n",
		"x")
}

// Java: InferConstsTest.testGeneratorFunctionConst()
func (c *InferConstsTest) TestGeneratorFunctionConst() {
	c.AssertConsts(
		"function *gen() {\n"+
			"  var x = 0;\n"+
			"  yield x;\n"+
			"}\n",
		InferConstsTest_NotDeclared("x"),
		InferConstsTest_Inferred("x"))
}

// Java: InferConstsTest.testGoogModuleExports()
func (c *InferConstsTest) TestGoogModuleExports() {
	c.AssertConsts("goog.module('m'); exports = {};")
}

// Java: InferConstsTest.testEsModuleImports()
func (c *InferConstsTest) TestEsModuleImports() {
	c.IgnoreWarnings(DiagnosticGroups_MODULE_LOAD)
	c.AssertConsts("import x from './mod';", InferConstsTest_NotInferred("x"))
	c.AssertConsts("import {x as y, z} from './mod';", InferConstsTest_NotInferred("y", "z"))
	c.AssertConsts("import * as x from './mod';", InferConstsTest_NotInferred("x"))
	c.AssertNotConsts("import * as CONST_NAME from './mod';", "CONST_NAME")
}

// Java: InferConstsTest.testConstantByConventionButNotInPractice()
func (c *InferConstsTest) TestConstantByConventionButNotInPractice() {
	c.AssertNotConsts("let CONST_NAME = 0; CONST_NAME++;", "CONST_NAME")
	c.AssertNotConsts("var CONST_NAME = 0; CONST_NAME++;", "CONST_NAME")
	c.AssertConsts("var CONST_NAME = 0;", InferConstsTest_Inferred("CONST_NAME"))
}

// Java: InferConstsTest.assertNotConsts(String, String...)
func (c *InferConstsTest) AssertNotConsts(js string, names ...string) {
	c.AssertConsts(js, InferConstsTest_NotDeclared(names...), InferConstsTest_NotInferred(names...))
}

// Java: InferConstsTest.assertConsts(String, Expectation...)
func (c *InferConstsTest) AssertConsts(js string, expectations ...*InferConstsTest_Expectation) {
	// stream(expectations).flatMap(e -> e.names.stream()).distinct().collect(toImmutableList())
	distinct := jutil.NewLinkedHashSet[string]()
	for _, e := range expectations {
		distinct.AddAll(e.names...)
	}
	c.names = distinct.ToSlice()

	c.TestSame(js)

	for _, expectation := range expectations {
		namesToCheck := c.constFinder.declaredNodes
		if expectation.isInferred {
			namesToCheck = c.constFinder.inferredNodes
		}
		setName := "declared constant names"
		if expectation.isInferred {
			setName = "inferred constant names"
		}
		for _, name := range expectation.names {
			if expectation.isConst {
				AssertWithMessage(setName).That(namesToCheck).Contains(name)
			} else {
				AssertWithMessage(setName).That(namesToCheck).DoesNotContain(name)
			}
		}
	}
}

// Java: InferConstsTest.declared(String...)
func InferConstsTest_Declared(names ...string) *InferConstsTest_Expectation {
	return NewInferConstsTest_Expectation(false /* isInferred= */, true /* isConst= */, names...)
}

// Java: InferConstsTest.inferred(String...)
func InferConstsTest_Inferred(names ...string) *InferConstsTest_Expectation {
	return NewInferConstsTest_Expectation(true /* isInferred= */, true /* isConst= */, names...)
}

// Java: InferConstsTest.notInferred(String...)
func InferConstsTest_NotInferred(names ...string) *InferConstsTest_Expectation {
	return NewInferConstsTest_Expectation(true /* isInferred= */, false /* isConst= */, names...)
}

// Java: InferConstsTest.notDeclared(String...)
func InferConstsTest_NotDeclared(names ...string) *InferConstsTest_Expectation {
	return NewInferConstsTest_Expectation(false /* isInferred= */, false /* isConst= */, names...)
}

// InferConstsTest_Expectation is the nested class Expectation.
type InferConstsTest_Expectation struct {
	names      []string
	isInferred bool
	isConst    bool
}

// Java: InferConstsTest.Expectation(boolean, boolean, String...)
func NewInferConstsTest_Expectation(isInferred bool, isConst bool, names ...string) *InferConstsTest_Expectation {
	return &InferConstsTest_Expectation{names: append([]string(nil), names...), isInferred: isInferred, isConst: isConst}
}

// InferConstsTest_FindConstants is the nested class FindConstants.
type InferConstsTest_FindConstants struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	names         []string
	declaredNodes *jutil.LinkedHashSet[string]
	inferredNodes *jutil.LinkedHashSet[string]
}

// Java: InferConstsTest.FindConstants(ImmutableList<String>)
func NewInferConstsTest_FindConstants(names []string) *InferConstsTest_FindConstants {
	return &InferConstsTest_FindConstants{
		names:         names,
		declaredNodes: jutil.NewLinkedHashSet[string](),
		inferredNodes: jutil.NewLinkedHashSet[string](),
	}
}

// Java: InferConstsTest.FindConstants.visit(NodeTraversal, Node, Node)
func (f *InferConstsTest_FindConstants) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	for _, name := range f.names {
		if (n.IsName() || n.IsImportStar()) && n.MatchesQualifiedName(name) {
			if n.IsDeclaredConstantVar() {
				f.declaredNodes.Add(name)
			}
			if n.IsInferredConstantVar() {
				f.inferredNodes.Add(name)
			}
		}
	}
}
