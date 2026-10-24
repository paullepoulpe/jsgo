// Port of com/google/javascript/jscomp/StatementFusionTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/StatementFusionTest.java)

/*
 * Copyright 2011 The Closure Compiler Authors.
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

// StatementFusionTest is the unit tests for StatementFusion.
type StatementFusionTest struct {
	CompilerTestCaseBase
}

// Java: StatementFusionTest()
func NewStatementFusionTest() *StatementFusionTest {
	c := &StatementFusionTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestStatementFusionTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewStatementFusionTest() })
}

// Java: StatementFusionTest.getProcessor(Compiler)
func (c *StatementFusionTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewPeepholeOptimizationsPass(compiler, c.GetName(), NewStatementFusion())
}

// Java: StatementFusionTest.testNothingToDo()
func (c *StatementFusionTest) TestNothingToDo() {
	c.FuseSame("")
	c.FuseSame("a")
	c.FuseSame("a()")
	c.FuseSame("if(a()){}")
}

// Java: StatementFusionTest.testFoldBlockWithStatements()
func (c *StatementFusionTest) TestFoldBlockWithStatements() {
	c.Fuse("a;b;c", "a,b,c")
	c.Fuse("a();b();c();", "a(),b(),c()")
	c.Fuse("a(),b();c(),d()", "a(),b(),c(),d()")
	c.Fuse("a();b(),c(),d()", "a(),b(),c(),d()")
	c.Fuse("a(),b(),c();d()", "a(),b(),c(),d()")
}

// Java: StatementFusionTest.testFoldBlockIntoIf()
func (c *StatementFusionTest) TestFoldBlockIntoIf() {
	c.Fuse("a;b;c;if(x){}", "if(a,b,c,x){}")
	c.Fuse("a;b;c;if(x,y){}else{}", "if(a,b,c,x,y){}else{}")
	c.Fuse("a;b;c;if(x,y){}", "if(a,b,c,x,y){}")
	c.Fuse("a;b;c;if(x,y,z){}", "if(a,b,c,x,y,z){}")
	// Can't fuse if there are statements after the IF.
	c.FuseSame("a();if(a()){}a()")
}

// Java: StatementFusionTest.testFoldBlockReturn()
func (c *StatementFusionTest) TestFoldBlockReturn() {
	c.Fuse("a;b;c;return x", "return a,b,c,x")
	c.Fuse("a;b;c;return x+y", "return a,b,c,x+y")
	// DeadAssignmentElimination would have cleaned it up anyways.
	c.FuseSame("a;b;c;return x;a;b;c")
}

// Java: StatementFusionTest.testFoldBlockThrow()
func (c *StatementFusionTest) TestFoldBlockThrow() {
	c.Fuse("a;b;c;throw x", "throw a,b,c,x")
	c.Fuse("a;b;c;throw x+y", "throw a,b,c,x+y")
	c.FuseSame("a;b;c;throw x;a;b;c")
}

// Java: StatementFusionTest.testFoldSwitch()
func (c *StatementFusionTest) TestFoldSwitch() {
	c.Fuse("a;b;c;switch(x){}", "switch(a,b,c,x){}")
}

// Java: StatementFusionTest.testFuseIntoForIn1()
func (c *StatementFusionTest) TestFuseIntoForIn1() {
	c.Fuse("a;b;c;for(x in y){}", "for(x in a,b,c,y){}")
}

// Java: StatementFusionTest.testFuseIntoForIn2()
func (c *StatementFusionTest) TestFuseIntoForIn2() {
	// This test case causes a parse warning in ES5 strict out, but is a parse error in ES6+ out.
	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT5_STRICT)
	c.SetExpectParseWarningsInThisTest()
	c.FuseSame("a();for(var x = b() in y){}")
}

// Java: StatementFusionTest.testFuseIntoVanillaFor1()
func (c *StatementFusionTest) TestFuseIntoVanillaFor1() {
	c.Fuse("a;b;c;for(;g;){}", "for(a,b,c;g;){}")
	c.Fuse("a;b;c;for(d;g;){}", "for(a,b,c,d;g;){}")
	c.Fuse("a;b;c;for(d,e;g;){}", "for(a,b,c,d,e;g;){}")
	c.FuseSame("a();for(var x;g;){}")
}

// Java: StatementFusionTest.testFuseIntoVanillaFor2()
func (c *StatementFusionTest) TestFuseIntoVanillaFor2() {
	c.FuseSame("a;b;c;for(var d;g;){}")
	c.FuseSame("a;b;c;for(let d;g;){}")
	c.FuseSame("a;b;c;for(const d = 5;g;){}")
}

// Java: StatementFusionTest.testFuseIntoLabel()
func (c *StatementFusionTest) TestFuseIntoLabel() {
	c.Fuse("a;b;c;label:for(x in y){}", "label:for(x in a,b,c,y){}")
	c.Fuse("a;b;c;label:for(;g;){}", "label:for(a,b,c;g;){}")
	c.Fuse("a;b;c;l1:l2:l3:for(;g;){}", "l1:l2:l3:for(a,b,c;g;){}")
	c.FuseSame("a;b;c;label:while(true){}")
}

// Java: StatementFusionTest.testFuseIntoBlock()
func (c *StatementFusionTest) TestFuseIntoBlock() {
	c.Fuse("a;b;c;{d;e;f}", "{a,b,c,d,e,f}")
	c.Fuse("a;b; label: { if(q) break label; bar(); }", "label: { if(a,b,q) break label; bar(); }")
	c.FuseSame("a;b;c;{var x;d;e;}")
	c.FuseSame("a;b;c;label:{break label;d;e;}")
}

// Java: StatementFusionTest.testNoFuseIntoWhile()
func (c *StatementFusionTest) TestNoFuseIntoWhile() {
	c.FuseSame("a;b;c;while(x){}")
}

// Java: StatementFusionTest.testNoFuseIntoDo()
func (c *StatementFusionTest) TestNoFuseIntoDo() {
	c.FuseSame("a;b;c;do{}while(x)")
}

// Java: StatementFusionTest.testNoFuseIntoBlock()
func (c *StatementFusionTest) TestNoFuseIntoBlock() {
	// Never fuse a statement into a block that contains let/const/class declarations, or you risk
	// colliding variable names. (unless the AST is normalized).
	c.Fuse("a; {b;}", "{a,b;}")
	c.Fuse("a; {b; var a = 1;}", "{a,b; var a = 1;}")
	c.FuseSame("a; { b; let a = 1; }")
	c.FuseSame("a; { b; const a = 1; }")
	c.FuseSame("a; { b; class a {} }")
	c.FuseSame("a; { b; function a() {} }")
	c.FuseSame("a; { b; const otherVariable = 1; }")
	c.EnableNormalize()
	c.Test(
		"function f(a) { if (COND) { a; { b; let a = 1; } } }",
		"function f(a) { if (COND) { { a,b; let a$jscomp$1 = 1; } } }")
	c.Test(
		"function f(a) { if (COND) { a; { b; let otherVariable = 1; } } }",
		"function f(a) { if (COND) {  { a,b; let otherVariable = 1; } } }")
}

// Java: StatementFusionTest.testNoGlobalScopeChanges()
func (c *StatementFusionTest) TestNoGlobalScopeChanges() {
	c.TestSame("a,b,c")
}

// Java: StatementFusionTest.testNoFunctionBlockChanges()
func (c *StatementFusionTest) TestNoFunctionBlockChanges() {
	c.TestSame("function foo() { a,b,c }")
}

// Java: StatementFusionTest.fuse(String, String)
func (c *StatementFusionTest) Fuse(before string, after string) {
	c.Test(
		"function F(){if(CONDITION){"+before+"}}",
		"function F(){if(CONDITION){"+after+"}}")
}

// Java: StatementFusionTest.fuseSame(String)
func (c *StatementFusionTest) FuseSame(code string) {
	c.TestSame("function F(){if(CONDITION){" + code + "}}")
}
