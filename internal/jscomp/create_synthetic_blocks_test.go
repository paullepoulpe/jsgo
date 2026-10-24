// Port of com/google/javascript/jscomp/CreateSyntheticBlocksTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/CreateSyntheticBlocksTest.java)

/*
 * Copyright 2009 The Closure Compiler Authors.
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

	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CreateSyntheticBlocksTest holds tests for CreateSyntheticBlocks.
type CreateSyntheticBlocksTest struct {
	CompilerTestCaseBase
}

// Java: CreateSyntheticBlocksTest.START_MARKER
const createSyntheticBlocksTest_START_MARKER = "startMarker"

// Java: CreateSyntheticBlocksTest.END_MARKER
const createSyntheticBlocksTest_END_MARKER = "endMarker"

func NewCreateSyntheticBlocksTest() *CreateSyntheticBlocksTest {
	c := &CreateSyntheticBlocksTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestCreateSyntheticBlocksTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCreateSyntheticBlocksTest() })
}

// Java: CreateSyntheticBlocksTest.setUp()
func (c *CreateSyntheticBlocksTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	// Can't use compare as a tree because of the added synthetic blocks.
	c.DisableCompareAsTree()
}

// Java: CreateSyntheticBlocksTest.getOptions()
func (c *CreateSyntheticBlocksTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	// This is used by NodeTraversal and not by the pass under test here
	startMarker := createSyntheticBlocksTest_START_MARKER
	options.SetSyntheticBlockStartMarker(&startMarker)
	return options
}

// createSyntheticBlocksTest_pass is the anonymous CompilerPass of
// getProcessor.
type createSyntheticBlocksTest_pass struct {
	c        *CreateSyntheticBlocksTest
	compiler *Compiler
}

// Java: CreateSyntheticBlocksTest.getProcessor(Compiler).new CompilerPass() {...}.process(Node, Node)
func (p *createSyntheticBlocksTest_pass) Process(externs rhino.Node, js rhino.Node) {
	NewCreateSyntheticBlocks(p.compiler, createSyntheticBlocksTest_START_MARKER, createSyntheticBlocksTest_END_MARKER).Process(externs, js)
	Normalize_CreateNormalizeForOptimizations(p.compiler).Process(externs, js)
	NewPeepholeOptimizationsPass(
		p.compiler,
		p.c.GetName(),
		NewMinimizeExitPoints(),
		NewPeepholeRemoveDeadCode(),
		NewPeepholeMinimizeConditions( /* late= */ true),
		NewPeepholeFoldConstants(true, false /* useTypes */)).
		Process(externs, js)
	NewDenormalize(p.compiler, parser.FeatureSet_BARE_MINIMUM).Process(externs, js)
}

// Java: CreateSyntheticBlocksTest.getProcessor(Compiler)
func (c *CreateSyntheticBlocksTest) GetProcessor(compiler *Compiler) CompilerPass {
	return &createSyntheticBlocksTest_pass{c: c, compiler: compiler}
}

// TODO(johnlenz): Add tests to the IntegrationTest.

// Java: CreateSyntheticBlocksTest.testFold1()
func (c *CreateSyntheticBlocksTest) TestFold1() {
	c.Test("function f() { if (x) return; y(); }", "function f(){x||y()}")
}

// Java: CreateSyntheticBlocksTest.testFoldWithMarkers1()
func (c *CreateSyntheticBlocksTest) TestFoldWithMarkers1() {
	c.TestSame("function f(){startMarker();if(x)return;endMarker();y()}")
}

// Java: CreateSyntheticBlocksTest.testFoldWithMarkers1a()
func (c *CreateSyntheticBlocksTest) TestFoldWithMarkers1a() {
	c.TestSame("function f(){startMarker();if(x)return;endMarker()}")
}

// Java: CreateSyntheticBlocksTest.testFold2()
func (c *CreateSyntheticBlocksTest) TestFold2() {
	c.Test(
		"function f() { if (x) return; y(); if (a) return; b(); }",
		"function f(){if(!x){y();a||b()}}")
}

// Java: CreateSyntheticBlocksTest.testFoldWithMarkers2()
func (c *CreateSyntheticBlocksTest) TestFoldWithMarkers2() {
	c.TestSame(
		`function f(){startMarker("FOO");startMarker("BAR");` +
			`if(x)return;endMarker("BAR");y();if(a)return;` +
			`endMarker("FOO");b()}`)
}

// Java: CreateSyntheticBlocksTest.testUnmatchedStartMarker()
func (c *CreateSyntheticBlocksTest) TestUnmatchedStartMarker() {
	c.TestError("startMarker()", CreateSyntheticBlocks_UNMATCHED_START_MARKER)
}

// Java: CreateSyntheticBlocksTest.testUnmatchedEndMarker1()
func (c *CreateSyntheticBlocksTest) TestUnmatchedEndMarker1() {
	c.TestError("endMarker()", CreateSyntheticBlocks_UNMATCHED_END_MARKER)
}

// Java: CreateSyntheticBlocksTest.testUnmatchedEndMarker2()
func (c *CreateSyntheticBlocksTest) TestUnmatchedEndMarker2() {
	c.TestError("if(y){startMarker();x()}endMarker()", CreateSyntheticBlocks_UNMATCHED_END_MARKER)
}

// Java: CreateSyntheticBlocksTest.testInvalid1()
func (c *CreateSyntheticBlocksTest) TestInvalid1() {
	c.TestError("startMarker() && true", CreateSyntheticBlocks_INVALID_MARKER_USAGE)
}

// Java: CreateSyntheticBlocksTest.testInvalid2()
func (c *CreateSyntheticBlocksTest) TestInvalid2() {
	c.TestError("false && endMarker()", CreateSyntheticBlocks_INVALID_MARKER_USAGE)
}

// Java: CreateSyntheticBlocksTest.testDenormalize()
func (c *CreateSyntheticBlocksTest) TestDenormalize() {
	c.TestSame("startMarker();for(;;);endMarker()")
}

// Java: CreateSyntheticBlocksTest.testNonMarkingUse()
func (c *CreateSyntheticBlocksTest) TestNonMarkingUse() {
	c.TestSame("function foo(endMarker){}")
	c.TestSame("function foo(){startMarker:foo()}")
}

// Java: CreateSyntheticBlocksTest.testContainingBlockPreservation()
func (c *CreateSyntheticBlocksTest) TestContainingBlockPreservation() {
	c.TestSame("if(y){startMarker();x();endMarker()}")
}

// Java: CreateSyntheticBlocksTest.testArrowFunction()
func (c *CreateSyntheticBlocksTest) TestArrowFunction() {
	c.TestSame("var y=()=>{startMarker();x();endMarker()}")
	c.TestError(
		"var y=()=>{startMarker();x();};endMarker()", CreateSyntheticBlocks_UNMATCHED_END_MARKER)
	c.TestError("var y=()=>startMarker();", CreateSyntheticBlocks_INVALID_MARKER_USAGE)
}

// Java: CreateSyntheticBlocksTest.testFunctionDeclaration1()
func (c *CreateSyntheticBlocksTest) TestFunctionDeclaration1() {
	c.Test(
		"startMarker(); a(); function fn() {}; b(); endMarker()",
		"startMarker();var fn=function(){};a();b();endMarker()")
}

// Java: CreateSyntheticBlocksTest.testFunctionDeclaration2()
func (c *CreateSyntheticBlocksTest) TestFunctionDeclaration2() {
	c.TestSame("startMarker();a();var fn=function(){};b();endMarker()")
}

// Java: CreateSyntheticBlocksTest.testClassDeclaration1()
func (c *CreateSyntheticBlocksTest) TestClassDeclaration1() {
	// Document that classes are mishandled with regard to block scoping
	c.TestSame("startMarker();class C{}endMarker()")
}

// Java: CreateSyntheticBlocksTest.testVariableDeclaration1()
func (c *CreateSyntheticBlocksTest) TestVariableDeclaration1() {
	// Document that let and const are mishandled with regard to block
	// scoping
	c.TestSame("startMarker();var x=1;endMarker()")
	c.TestSame("startMarker();let x=1;endMarker()")
	c.TestSame("startMarker();const x=1;endMarker()")
}

// Java: CreateSyntheticBlocksTest.testSyntheticBlock_doesNotCreateNewScope()
func (c *CreateSyntheticBlocksTest) TestSyntheticBlock_doesNotCreateNewScope() {
	c.TestSame("startMarker();var x=1;endMarker()")
	script := c.GetLastCompiler().GetJsRoot().GetOnlyChild()
	AssertNode(script).IsScript()
	synctheticBlock := script.GetFirstChild()
	AssertNode(synctheticBlock).IsBlock()
	AssertThat(synctheticBlock.IsSyntheticBlock()).IsTrue()
	// confirm that synthetic block does not create a new block scope
	AssertThat(NodeUtil_CreatesBlockScope(synctheticBlock)).IsFalse()
}
