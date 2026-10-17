// Port of com/google/javascript/jscomp/MinimizeExitPointsTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/MinimizeExitPointsTest.java)

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

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// MinimizeExitPointsTest is the unit tests for MinimizeExitPoints.
type MinimizeExitPointsTest struct {
	CompilerTestCaseBase
}

// Java: MinimizeExitPointsTest()
func NewMinimizeExitPointsTest() *MinimizeExitPointsTest {
	c := &MinimizeExitPointsTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestMinimizeExitPointsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewMinimizeExitPointsTest() })
}

// Java: MinimizeExitPointsTest.setUp()
func (c *MinimizeExitPointsTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableNormalize()
}

// Java: MinimizeExitPointsTest.getProcessor(Compiler)
func (c *MinimizeExitPointsTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewPeepholeOptimizationsPass(compiler, c.GetName(), NewMinimizeExitPoints())
}

// Java: MinimizeExitPointsTest.foldSame(String)
func (c *MinimizeExitPointsTest) FoldSame(js string) {
	c.TestSame(js)
}

// Java: MinimizeExitPointsTest.fold(String, String)
func (c *MinimizeExitPointsTest) Fold(js string, expected string) {
	c.Test(js, expected)
}

// Java: MinimizeExitPointsTest.testBreakOptimization()
func (c *MinimizeExitPointsTest) TestBreakOptimization() {
	c.Fold("f:{if(true){a();break f;}else;b();}", "f:{if(true){a()}else{b()}}")
	c.Fold("f:{if(false){a();break f;}else;b();break f;}", "f:{if(false){a()}else{b()}}")
	c.Fold("f:{if(a()){b();break f;}else;c();}", "f:{if(a()){b();}else{c();}}")
	c.Fold("f:{if(a()){b()}else{c();break f;}}", "f:{if(a()){b()}else{c();}}")
	c.Fold("f:{if(a()){b();break f;}else;}", "f:{if(a()){b();}else;}")
	c.Fold("f:{if(a()){break f;}else;}", "f:{if(a()){}else;}")
	c.Fold("f:for(;a();)break f;", "f:for(;a();)break f")
	c.FoldSame("f:for(x in a())break f")
	c.FoldSame("f:{for(;a();)break;}")
	c.FoldSame("f:{for(x in a())break}")
	c.Fold("f:try{break f;}catch(e){break f;}", "f: { try{}catch(e){} }")
	c.Fold(
		"f:try{if(a()){break f;}else{break f;} break f;}catch(e){}",
		"f:{ try{if(a()){}else{}}catch(e){} }")
	c.Fold("f:g:break f", "f: g: {}")
	c.Fold("f:g:{if(a()){break f;}else{break f;} break f;}", "f:g:{if(a()){}else{}}")
	c.Fold("function f() { a: break a; }", "function f() { a: {} }")
	c.Fold("function f() { a: { break a; } }", "function f() { a: {} }")
}

// Java: MinimizeExitPointsTest.testFunctionReturnOptimization1()
func (c *MinimizeExitPointsTest) TestFunctionReturnOptimization1() {
	c.Fold("function f(){return}", "function f(){}")
}

// Java: MinimizeExitPointsTest.testFunctionReturnOptimization2()
func (c *MinimizeExitPointsTest) TestFunctionReturnOptimization2() {
	c.Fold("function f(){if(a()){b();if(c())return;}}", "function f(){if(a()){b();if(c());}}")
	c.Fold("function f(){if(x)return; x=3; return; }", "function f(){if(x); else x=3}")
	c.Fold("function f(){if(true){a();return;}else;b();}", "function f(){if(true){a();}else{b();}}")
	c.Fold(
		"function f(){if(false){a();return;}else;b();return;}",
		"function f(){if(false){a();}else{b();}}")
	c.Fold("function f(){if(a()){b();return;}else;c();}", "function f(){if(a()){b();}else{c();}}")
	c.Fold("function f(){if(a()){b()}else{c();return;}}", "function f(){if(a()){b()}else{c();}}")
	c.Fold("function f(){if(a()){b();return;}else;}", "function f(){if(a()){b();}else;}")
	c.Fold("function f(){if(a()){return;}else{return;} return;}", "function f(){if(a()){}else{}}")
	c.Fold(
		"function f(){if(a()){return;}else{return;} b();}",
		"function f(){if(a()){}else{return;b()}}")
	c.Fold(
		"function f(){ if (x) return; if (y) return; if (z) return; w(); }",
		"function f() {\n"+
			"  if (x) {} else { if (y) {} else { if (z) {} else w(); }}\n"+
			"}\n")
	c.Fold("function f(){for(;a();)return;}", "function f(){for(;a();)return}")
	c.FoldSame("function f(){for(x in a())return}")
	c.Fold("function f(){for(;a();)break;}", "function f(){for(;a();)break}")
	c.FoldSame("function f(){for(x in a())break}")
	c.Fold(
		"function f(){try{return;}catch(e){throw 9;}finally{return}}",
		"function f(){try{}catch(e){throw 9;}finally{return}}")
	c.FoldSame("function f(){try{throw 9;}finally{return;}}")
	c.Fold("function f(){try{return;}catch(e){return;}}", "function f(){try{}catch(e){}}")
	c.Fold(
		"function f(){try{if(a()){return;}else{return;} return;}catch(e){}}",
		"function f(){try{if(a()){}else{}}catch(e){}}")
	c.Fold("function f(){g:return}", "function f(){g: {}}")
	c.Fold(
		"function f(){g:if(a()){return;}else{return;} return;}",
		"function f(){g:{if(a()){}else{}}}")
	c.Fold(
		"function f() {\n"+
			"  try {\n"+
			"    g: if (a()) {\n"+
			"      throw 9;\n"+
			"    }\n"+
			"    return;\n"+
			"  } finally {\n"+
			"    return;\n"+
			"  }\n"+
			"}\n",
		"function f() {\n"+
			"  try {\n"+
			"    g: {\n"+
			"      if (a()) {\n"+
			"        throw 9;\n"+
			"      }\n"+
			"    }\n"+
			"  } finally {\n"+
			"    return;\n"+
			"  }\n"+
			"}\n")
}

// Java: MinimizeExitPointsTest.testFunctionReturnScoped()
func (c *MinimizeExitPointsTest) TestFunctionReturnScoped() {
	c.TestSame(
		"function f(a) {\n" +
			"  if (a) {\n" +
			"    const aInner = Math.random();\n" +
			"\n" +
			"    if (aInner < 0.5) {\n" +
			"        return aInner;\n" +
			"    }\n" +
			"  }\n" +
			"\n" +
			"  return a;\n" +
			"}\n")
}

// Java: MinimizeExitPointsTest.testWhileContinueOptimization()
func (c *MinimizeExitPointsTest) TestWhileContinueOptimization() {
	// Normalization should convert all WHILE loops to FOR loops, so just have a simple test
	// verifying that happens & we don't need explicit WHILE loop support.
	c.Fold("while(true){if(x)continue; x=3; continue; }", "for(;true;)if(x);else x=3")
}

// Java: MinimizeExitPointsTest.testDoContinueOptimization()
func (c *MinimizeExitPointsTest) TestDoContinueOptimization() {
	c.Fold("do{if(x)continue; x=3; continue; }while(true)", "do if(x); else x=3; while(true)")
	c.FoldSame("do{a();continue;b()}while(true)")
	c.Fold(
		"do{if(true){a();continue;}else;b();}while(true)",
		"do{if(true){a();}else{b();}}while(true)")
	c.Fold(
		"do{if(false){a();continue;}else;b();continue;}while(true)",
		"do{if(false){a();}else{b();}}while(true)")
	c.Fold("do{if(a()){b();continue;}else;c();}while(true)", "do{if(a()){b();}else{c()}}while(true)")
	c.Fold(
		"do{if(a()){b();}else{c();continue;}}while(true)",
		"do{if(a()){b();}else{c();}}while(true)")
	c.Fold("do{if(a()){b();continue;}else;}while(true)", "do{if(a()){b();}else;}while(true)")
	c.Fold(
		"do{if(a()){continue;}else{continue;} continue;}while(true)",
		"do{if(a()){}else{}}while(true)")
	c.Fold(
		"do{if(a()){continue;}else{continue;} b();}while(true)",
		"do{if(a()){}else{continue; b();}}while(true)")
	c.Fold("do{for(;a();)continue;}while(true)", "do for(;a(););while(true)")
	c.Fold("do{for(x in a())continue}while(true)", "do for(x in a());while(true)")
	c.Fold("do{for(;a();)break;}while(true)", "do for(;a();)break;while(true)")
	c.FoldSame("do for(x in a())break;while(true)")
	c.Fold("do{try{continue;}catch(e){continue;}}while(true)", "do{try{}catch(e){}}while(true)")
	c.Fold(
		"do{try{if(a()){continue;}else{continue;} continue;}catch(e){}}while(true)",
		"do{try{if(a()){}else{}}catch(e){}}while(true)")
	c.Fold("do{g:continue}while(true)", "do{g: {}}while(true)")
	// This case could be improved.
	c.Fold(
		"do{g:if(a()){continue;}else{continue;} continue;}while(true)",
		"do{g: { if(a());else; } }while(true)")
	c.Fold("do { foo(); continue; } while(false)", "do { foo(); } while(false)")
	c.Fold("do { foo(); break; } while(false)", "do { foo(); } while(false)")
	c.Fold("do{break}while(!new Date());", "do{}while(!new Date());")
	c.FoldSame("do { foo(); switch (x) { case 1: break; default: f()}; } while(false)")
}

// Java: MinimizeExitPointsTest.testForContinueOptimization()
func (c *MinimizeExitPointsTest) TestForContinueOptimization() {
	c.Fold("for(x in y){if(x)continue; x=3; continue; }", "for(x in y)if(x);else x=3")
	c.FoldSame("for(x in y){a();continue;b()}")
	c.Fold("for(x in y){if(true){a();continue;}else;b();}", "for(x in y){if(true)a();else b();}")
	c.Fold(
		"for(x in y){if(false){a();continue;}else;b();continue;}",
		"for(x in y){if(false){a();}else{b()}}")
	c.Fold("for(x in y){if(a()){b();continue;}else;c();}", "for(x in y){if(a()){b();}else{c();}}")
	c.Fold("for(x in y){if(a()){b();}else{c();continue;}}", "for(x in y){if(a()){b();}else{c();}}")
	c.Fold("for(x of y){if(x)continue; x=3; continue; }", "for(x of y)if(x);else x=3")
	c.FoldSame("for(x of y){a();continue;b()}")
	c.Fold("for(x of y){if(true){a();continue;}else;b();}", "for(x of y){if(true)a();else b();}")
	c.Fold(
		"for(x of y){if(false){a();continue;}else;b();continue;}",
		"for(x of y){if(false){a();}else{b()}}")
	c.Fold("for(x of y){if(a()){b();continue;}else;c();}", "for(x of y){if(a()){b();}else{c();}}")
	c.Fold("for(x of y){if(a()){b();}else{c();continue;}}", "for(x of y){if(a()){b();}else{c();}}")
	c.Fold(
		"async () => { for await (x of y){if(x)continue; x=3; continue; }}",
		"async () => { for await (x of y)if(x);else x=3 }")
	c.FoldSame("async () => { for await (x of y){a();continue;b()}}")
	c.Fold(
		"async () => { for await (x of y){if(true){a();continue;}else;b();}}",
		"async () => { for await (x of y){if(true)a();else b();}}")
	c.Fold(
		"async () => { for await (x of y){if(false){a();continue;}else;b();continue;}}",
		"async () => { for await (x of y){if(false){a();}else{b()}}}")
	c.Fold(
		"async () => { for await (x of y){if(a()){b();continue;}else;c();}}",
		"async () => { for await (x of y){if(a()){b();}else{c();}}}")
	c.Fold(
		"async () => { for await (x of y){if(a()){b();}else{c();continue;}}}",
		"async () => { for await (x of y){if(a()){b();}else{c();}}}")
	c.Fold("for(x=0;x<y;x++){if(a()){b();continue;}else;}", "x=0;for(;x<y;x++){if(a()){b();}else;}")
	c.Fold(
		"for(x=0;x<y;x++){if(a()){continue;}else{continue;} continue;}",
		"x=0;for(;x<y;x++){if(a()){}else{}}")
	c.Fold(
		"for(x=0;x<y;x++){if(a()){continue;}else{continue;} b();}",
		"x=0;for(;x<y;x++){if(a()){}else{continue; b();}}")
	c.Fold("for(x=0;x<y;x++)for(;a();)continue;", "x=0;for(;x<y;x++)for(;a(););")
	c.Fold("for(x=0;x<y;x++)for(x in a())continue", "x=0;for(;x<y;x++)for(x in a());")
	c.Fold("for(x=0;x<y;x++)for(;a();)break;", "x=0;for(;x<y;x++)for(;a();)break")
	c.FoldSame("x=0;for(;x<y;x++)for(x in a())break")
	c.Fold(
		"for(x=0;x<y;x++){try{continue;}catch(e){continue;}}",
		"x=0;for(;x<y;x++){try{}catch(e){}}")
	c.Fold(
		"for(x=0;x<y;x++){try{if(a()){continue;}else{continue;} continue;}catch(e){}}",
		"x=0;for(;x<y;x++){try{if(a()){}else{}}catch(e){}}")
	c.Fold("for(x=0;x<y;x++){g:continue}", "x=0;for(;x<y;x++){ g: {} }")
	c.Fold(
		"for(x=0;x<y;x++){g:if(a()){continue;}else{continue;} continue;}",
		"x=0;for(;x<y;x++){g:{if(a());else;}}")
}

// Java: MinimizeExitPointsTest.testCodeMotionDoesntBreakFunctionHoisting()
func (c *MinimizeExitPointsTest) TestCodeMotionDoesntBreakFunctionHoisting() {
	c.Fold(
		"function f() { if (x) return; foo(); function foo() {} }",
		// Note: normalization moves all function declarations to the beginning of a block.
		// This test cases was originally added before requiring normalization.
		"function f() { function foo() {} if (x) {} else { foo(); } }")
}

// Java: MinimizeExitPointsTest.testDontRemoveBreakInTryFinally()
func (c *MinimizeExitPointsTest) TestDontRemoveBreakInTryFinally() {
	c.FoldSame("function f() {b: {try{throw 9} finally {break b}} return 1;}")
}

// The 'break' prevents the 'b=false' from being evaluated. If we fold the do-while to
// 'do;while(b=false)' the code will be incorrect.
//
// @see https://github.com/google/closure-compiler/issues/554
// Java: MinimizeExitPointsTest.testDontFoldBreakInDoWhileIfConditionHasSideEffects()
func (c *MinimizeExitPointsTest) TestDontFoldBreakInDoWhileIfConditionHasSideEffects() {
	c.FoldSame("var b=true;do{break}while(b=false);")
}

// Java: MinimizeExitPointsTest.testSwitchExitPoints1()
func (c *MinimizeExitPointsTest) TestSwitchExitPoints1() {
	c.Fold("switch (x) { case 1: f(); break; }", "switch (x) { case 1: f();        }")
	c.Fold(
		"switch (x) { case 1: f(); break; case 2: g(); break; }",
		"switch (x) { case 1: f(); break; case 2: g();        }")
	c.Fold(
		"switch (x) { case 1: if (x) { f(); break; } break; default: g(); break; }",
		"switch (x) { case 1: if (x) { f();        } break; default: g();        }")
}

// Java: MinimizeExitPointsTest.testThrowsExceptionIfAstNotNormalized()
func (c *MinimizeExitPointsTest) TestThrowsExceptionIfAstNotNormalized() {
	c.DisableNormalize()
	AssertThrows[jutil.IRuntimeException](func() { c.FoldSame("") })
}

// Java: MinimizeExitPointsTest.testFoldBlockScopedVariables()
func (c *MinimizeExitPointsTest) TestFoldBlockScopedVariables() {
	// When moving block-scoped variable declarations into inner blocks, first convert them to
	// "var" declarations to avoid breaking any references in inner functions.
	// For example, in the following test case, moving "let c = 3;" directly inside the else block
	// would break the function "g"'s reference to "c".
	c.Fold(
		"function f() { function g() { return c; } if (x) {return;} let c = 3; }",
		"function f() { function g() { return c; } if (x){} else {var c = 3;} }")
	c.Fold(
		"function f() { function g() { return c; } if (x) {return;} const c = 3; }",
		"function f() { function g() { return c; } if (x) {} else {var c = 3;} }")
	// Convert let and const even they're if not referenced by any functions.
	c.Fold(
		"function f() { if (x) {return;} const c = 3; }",
		"function f() { if (x) {} else { var c = 3; } }")
	c.Fold(
		"function f() { if (x) {return;} let a = 3; let b = () => a; }",
		"function f() { if (x) {} else { var a = 3; var b = () => { return a; }} }")
	c.Fold(
		"function f() { if (x) { if (y) {return;} let c = 3; } }",
		"function f() { if (x) { if (y) {} else { var c = 3; } } }")
}

// Java: MinimizeExitPointsTest.testDontFoldBlockScopedVariablesInLoops()
func (c *MinimizeExitPointsTest) TestDontFoldBlockScopedVariablesInLoops() {
	// Don't move block-scoped declarations into inner blocks inside a loop, since converting
	// let/const declarations to vars in a loop can cause incorrect semantics.
	// See the following test case for an example.
	c.FoldSame(
		"function f(param) {\n" +
			"  let arr = [];\n" +
			"  for (let x of param) {\n" +
			"    if (x < 0) continue;\n" +
			"    let y = x * 2;\n" +
			"    arr.push(() => {\n" +
			"      return y;\n" +
			"    }); // If y were a var, this would capture the wrong value.\n" +
			"  }\n" +
			"  return arr;\n" +
			"}\n")
	// Additional tests for different kinds of loops.
	c.FoldSame("function f() { for (;true;) { if (true) {return;} let c = 3; } }")
	c.FoldSame("function f() { do { if (true) {return;} let c = 3; } while (x); }")
	c.FoldSame("function f() { for (;;) { if (true) { return; } let c = 3; } }")
	c.FoldSame("function f(y) { for(x in []){ if(x) { return; } let c = 3; } }")
	c.FoldSame("async function f(y) { for await (x in []){ if(x) { return; } let c = 3; } }")
}
