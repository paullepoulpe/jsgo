// Port of com/google/javascript/jscomp/instrumentation/CoverageInstrumentationPassTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/instrumentation/CoverageInstrumentationPassTest.java)

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

import "testing"

// CoverageInstrumentationPassTest tests CoverageInstrumentationPass.
//
// Java: public final class CoverageInstrumentationPassTest
type CoverageInstrumentationPassTest struct {
	t *testing.T
}

func (c *CoverageInstrumentationPassTest) JUnitSetT(t *testing.T) { c.t = t }
func (c *CoverageInstrumentationPassTest) JUnitBefore()           {}
func (c *CoverageInstrumentationPassTest) JUnitAfter()            {}

// Java: CoverageInstrumentationPassTest()
func NewCoverageInstrumentationPassTest() *CoverageInstrumentationPassTest {
	return &CoverageInstrumentationPassTest{}
}

func TestCoverageInstrumentationPassTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCoverageInstrumentationPassTest() })
}

// Java: CoverageInstrumentationPassTest.options(LanguageMode, LanguageMode, boolean)
func (c *CoverageInstrumentationPassTest) optionsWithOut(inMode CompilerOptions_LanguageMode, outMode CompilerOptions_LanguageMode, coverageOnly bool) *CompilerOptions {
	options := GoldenFileComparer_Options()
	options.SetInstrumentForCoverageOption(CompilerOptions_InstrumentOption_LINE_ONLY)
	options.SetInstrumentForCoverageOnly(coverageOnly)
	options.SetLanguageIn(inMode)
	options.SetLanguageOut(outMode)
	options.SetWarningLevel(DiagnosticGroups_MODULE_LOAD, CheckLevel_OFF)
	return options
}

// Java: CoverageInstrumentationPassTest.options(LanguageMode)
func (c *CoverageInstrumentationPassTest) options(inMode CompilerOptions_LanguageMode) *CompilerOptions {
	return c.optionsWithOut(inMode, CompilerOptions_LanguageMode_ECMASCRIPT5 /* coverageOnly= */, false)
}

// Java: CoverageInstrumentationPassTest.branchOptions(LanguageMode)
func (c *CoverageInstrumentationPassTest) branchOptions(inMode CompilerOptions_LanguageMode) *CompilerOptions {
	options := GoldenFileComparer_Options()
	options.SetInstrumentForCoverageOption(CompilerOptions_InstrumentOption_BRANCH_ONLY)
	options.SetLanguageIn(inMode)
	options.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	return options
}

// Java: CoverageInstrumentationPassTest.compareFunctionOneMode(LanguageMode)
func (c *CoverageInstrumentationPassTest) compareFunctionOneMode(mode CompilerOptions_LanguageMode) {
	GoldenFileComparer_CompileAndCompare(
		"CoverageInstrumentationPassTest/FunctionGolden.jsdata",
		c.options(mode),
		"CoverageInstrumentationPassTest/Function.jsdata")
}

// Java: CoverageInstrumentationPassTest.testFunction()
func (c *CoverageInstrumentationPassTest) TestFunction() {
	c.compareFunctionOneMode(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.compareFunctionOneMode(CompilerOptions_LanguageMode_ECMASCRIPT_2015)
}

// Java: CoverageInstrumentationPassTest.compareArrowOneMode(LanguageMode, String)
func (c *CoverageInstrumentationPassTest) compareArrowOneMode(mode CompilerOptions_LanguageMode, prefix string) {
	GoldenFileComparer_CompileAndCompare(
		prefix+"Golden.jsdata", c.options(mode), prefix+".jsdata")
}

// If the body of the arrow function is a block, it is instrumented.
//
// Java: CoverageInstrumentationPassTest.testArrowFunction_block()
func (c *CoverageInstrumentationPassTest) TestArrowFunction_block() {
	c.compareArrowOneMode(CompilerOptions_LanguageMode_ECMASCRIPT_2015, "CoverageInstrumentationPassTest/ArrowBlock")
}

// If the body of the arrow function is an expression, it is converted to a
// block, then instrumented.
//
// Java: CoverageInstrumentationPassTest.testArrowFunction_expression()
func (c *CoverageInstrumentationPassTest) TestArrowFunction_expression() {
	c.compareArrowOneMode(
		CompilerOptions_LanguageMode_ECMASCRIPT_2015, "CoverageInstrumentationPassTest/ArrowExpression")
}

// Java: CoverageInstrumentationPassTest.compareIfBranch(LanguageMode)
func (c *CoverageInstrumentationPassTest) compareIfBranch(mode CompilerOptions_LanguageMode) {
	GoldenFileComparer_CompileAndCompare(
		"CoverageInstrumentationPassTest/IfBranchGolden.jsdata",
		c.branchOptions(mode),
		"CoverageInstrumentationPassTest/IfBranch.jsdata")
}

// Java: CoverageInstrumentationPassTest.testIfBranch()
func (c *CoverageInstrumentationPassTest) TestIfBranch() {
	c.compareIfBranch(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.compareIfBranch(CompilerOptions_LanguageMode_ECMASCRIPT_2015)
}

// Java: CoverageInstrumentationPassTest.compareIfElseBranch(LanguageMode)
func (c *CoverageInstrumentationPassTest) compareIfElseBranch(mode CompilerOptions_LanguageMode) {
	GoldenFileComparer_CompileAndCompare(
		"CoverageInstrumentationPassTest/IfElseBranchGolden.jsdata",
		c.branchOptions(mode),
		"CoverageInstrumentationPassTest/IfElseBranch.jsdata")
}

// Java: CoverageInstrumentationPassTest.testIfElseBranch()
func (c *CoverageInstrumentationPassTest) TestIfElseBranch() {
	c.compareIfElseBranch(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.compareIfElseBranch(CompilerOptions_LanguageMode_ECMASCRIPT_2015)
}

// Java: CoverageInstrumentationPassTest.compareForLoopBranch(LanguageMode)
func (c *CoverageInstrumentationPassTest) compareForLoopBranch(mode CompilerOptions_LanguageMode) {
	GoldenFileComparer_CompileAndCompare(
		"CoverageInstrumentationPassTest/ForLoopBranchGolden.jsdata",
		c.branchOptions(mode),
		"CoverageInstrumentationPassTest/ForLoopBranch.jsdata")
}

// Java: CoverageInstrumentationPassTest.testForLoopBranch()
func (c *CoverageInstrumentationPassTest) TestForLoopBranch() {
	c.compareForLoopBranch(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.compareForLoopBranch(CompilerOptions_LanguageMode_ECMASCRIPT_2015)
}

// Java: CoverageInstrumentationPassTest.compareDoWhileLoopBranch(LanguageMode)
func (c *CoverageInstrumentationPassTest) compareDoWhileLoopBranch(mode CompilerOptions_LanguageMode) {
	GoldenFileComparer_CompileAndCompare(
		"CoverageInstrumentationPassTest/DoWhileLoopBranchGolden.jsdata",
		c.branchOptions(mode),
		"CoverageInstrumentationPassTest/DoWhileLoopBranch.jsdata")
}

// Java: CoverageInstrumentationPassTest.testDoWhileLoopBranch()
func (c *CoverageInstrumentationPassTest) TestDoWhileLoopBranch() {
	c.compareDoWhileLoopBranch(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.compareDoWhileLoopBranch(CompilerOptions_LanguageMode_ECMASCRIPT_2015)
}

// Java: CoverageInstrumentationPassTest.compareDoWhileLoopMultiLineBranch(LanguageMode)
func (c *CoverageInstrumentationPassTest) compareDoWhileLoopMultiLineBranch(mode CompilerOptions_LanguageMode) {
	GoldenFileComparer_CompileAndCompare(
		"CoverageInstrumentationPassTest/DoWhileLoopMultiLineBranchGolden.jsdata",
		c.branchOptions(mode),
		"CoverageInstrumentationPassTest/DoWhileLoopMultiLineBranch.jsdata")
}

// Java: CoverageInstrumentationPassTest.testDoWhileLoopMultiLineBranch()
func (c *CoverageInstrumentationPassTest) TestDoWhileLoopMultiLineBranch() {
	c.compareDoWhileLoopMultiLineBranch(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.compareDoWhileLoopMultiLineBranch(CompilerOptions_LanguageMode_ECMASCRIPT_2015)
}

// Java: CoverageInstrumentationPassTest.compareWhileLoopBranch(LanguageMode)
func (c *CoverageInstrumentationPassTest) compareWhileLoopBranch(mode CompilerOptions_LanguageMode) {
	GoldenFileComparer_CompileAndCompare(
		"CoverageInstrumentationPassTest/WhileLoopBranchGolden.jsdata",
		c.branchOptions(mode),
		"CoverageInstrumentationPassTest/WhileLoopBranch.jsdata")
}

// Java: CoverageInstrumentationPassTest.testWhileLoopBranch()
func (c *CoverageInstrumentationPassTest) TestWhileLoopBranch() {
	c.compareWhileLoopBranch(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.compareWhileLoopBranch(CompilerOptions_LanguageMode_ECMASCRIPT_2015)
}

// Java: CoverageInstrumentationPassTest.testEsModule()
func (c *CoverageInstrumentationPassTest) TestEsModule() {
	GoldenFileComparer_CompileAndCompare(
		"CoverageInstrumentationPassTest/EsModuleGolden.jsdata",
		c.optionsWithOut(
			CompilerOptions_LanguageMode_ECMASCRIPT_NEXT, CompilerOptions_LanguageMode_ECMASCRIPT_NEXT /* coverageOnly= */, true),
		"CoverageInstrumentationPassTest/EsModule.jsdata")
}

// Java: CoverageInstrumentationPassTest.testNestedForLoop()
func (c *CoverageInstrumentationPassTest) TestNestedForLoop() {
	GoldenFileComparer_CompileAndCompare(
		"CoverageInstrumentationPassTest/NestedForLoopBranchGolden.jsdata",
		c.branchOptions(CompilerOptions_LanguageMode_STABLE),
		"CoverageInstrumentationPassTest/NestedForLoopBranch.jsdata")
}

// Java: CoverageInstrumentationPassTest.testForLoopNoEndingBlock()
func (c *CoverageInstrumentationPassTest) TestForLoopNoEndingBlock() {
	GoldenFileComparer_CompileAndCompare(
		"CoverageInstrumentationPassTest/ForLoopBranchNoEndingBlockGolden.jsdata",
		c.branchOptions(CompilerOptions_LanguageMode_STABLE),
		"CoverageInstrumentationPassTest/ForLoopBranchNoEndingBlock.jsdata")
}

// The unit test covers a case where a branch condition on line x of the
// source code matches a branching condition on the same line number X of an
// injected library and will cause a NPE. In this case, the if statement on
// line 21 of ForLoopBranchConflictWithPolyfill.jsdata matches an if statement
// on line 21 of find.js. If the branch in find.js is moved, this test case is
// no longer useful.
//
// Java: CoverageInstrumentationPassTest.testForLoopConflictWithPolyfill()
func (c *CoverageInstrumentationPassTest) TestForLoopConflictWithPolyfill() {
	options := c.branchOptions(CompilerOptions_LanguageMode_STABLE)
	options.SetForceLibraryInjection([]string{"es6/array/find"})
	GoldenFileComparer_CompileAndCompareSubsetOfActualToExpected(
		"CoverageInstrumentationPassTest/ForLoopBranchConflictWithPolyfillGolden.jsdata",
		options,
		"CoverageInstrumentationPassTest/ForLoopBranchConflictWithPolyfill.jsdata")
}

// Java: CoverageInstrumentationPassTest.testNoCoverage()
func (c *CoverageInstrumentationPassTest) TestNoCoverage() {
	GoldenFileComparer_CompileAndCompare(
		"CoverageInstrumentationPassTest/NoCoverageGolden.jsdata",
		c.options(CompilerOptions_LanguageMode_STABLE),
		"CoverageInstrumentationPassTest/NoCoverage.jsdata")
}
