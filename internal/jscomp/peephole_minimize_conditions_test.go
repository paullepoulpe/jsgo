// Port of com/google/javascript/jscomp/PeepholeMinimizeConditionsTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/PeepholeMinimizeConditionsTest.java)

/*
 * Copyright 2004 The Closure Compiler Authors.
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

// PeepholeMinimizeConditionsTest is the tests for PeepholeMinimizeConditions
// in isolation. Tests for the interaction of multiple peephole passes are in
// PeepholeIntegrationTest.
type PeepholeMinimizeConditionsTest struct {
	CompilerTestCaseBase

	late bool
}

// Java: PeepholeMinimizeConditionsTest()
func NewPeepholeMinimizeConditionsTest() *PeepholeMinimizeConditionsTest {
	c := &PeepholeMinimizeConditionsTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, CompilerTestCase_DEFAULT_EXTERNS)
	c.late = true
	return c
}

func TestPeepholeMinimizeConditionsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewPeepholeMinimizeConditionsTest() })
}

// Java: PeepholeMinimizeConditionsTest.setUp()
func (c *PeepholeMinimizeConditionsTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.late = true
	c.DisableTypeCheck()
}

// Java: PeepholeMinimizeConditionsTest.getProcessor(Compiler)
func (c *PeepholeMinimizeConditionsTest) GetProcessor(compiler *Compiler) CompilerPass {
	peepholePass := NewPeepholeOptimizationsPass(compiler, c.GetName(), NewPeepholeMinimizeConditions(c.late))
	peepholePass.SetRetraverseOnChange(false)
	return peepholePass
}

// Java: PeepholeMinimizeConditionsTest.foldSame(String)
func (c *PeepholeMinimizeConditionsTest) FoldSame(js string) {
	c.TestSame(js)
}

// Java: PeepholeMinimizeConditionsTest.fold(String, String)
func (c *PeepholeMinimizeConditionsTest) Fold(js string, expected string) {
	c.Test(js, expected)
}

// Check that removing blocks with 1 child works
// Java: PeepholeMinimizeConditionsTest.testFoldOneChildBlocks()
func (c *PeepholeMinimizeConditionsTest) TestFoldOneChildBlocks() {
	c.late = false
	c.Fold("function f(){if(x)a();x=3}", "function f(){x&&a();x=3}")
	c.Fold("function f(){if(x)a?.();x=3}", "function f(){x&&a?.();x=3}")
	c.Fold("function f(){if(x){a()}x=3}", "function f(){x&&a();x=3}")
	c.Fold("function f(){if(x){a?.()}x=3}", "function f(){x&&a?.();x=3}")
	c.Fold("function f(){if(x){return 3}}", "function f(){if(x)return 3}")
	c.Fold("function f(){if(x){a()}}", "function f(){x&&a()}")
	c.Fold("function f(){if(x){throw 1}}", "function f(){if(x)throw 1;}")
	// Try it out with functions
	c.Fold("function f(){if(x){foo()}}", "function f(){x&&foo()}")
	c.Fold("function f(){if(x){foo()}else{bar()}}", "function f(){x?foo():bar()}")
	// Try it out with properties and methods
	c.Fold("function f(){if(x){a.b=1}}", "function f(){if(x)a.b=1}")
	c.Fold("function f(){if(x){a.b*=1}}", "function f(){x&&(a.b*=1)}")
	c.Fold("function f(){if(x){a.b+=1}}", "function f(){x&&(a.b+=1)}")
	c.Fold("function f(){if(x){++a.b}}", "function f(){x&&++a.b}")
	c.Fold("function f(){if(x){a.foo()}}", "function f(){x&&a.foo()}")
	c.Fold("function f(){if(x){a?.foo()}}", "function f(){x&&a?.foo()}")
	// Try it out with throw/catch/finally [which should not change]
	c.FoldSame("function f(){try{foo()}catch(e){bar(e)}finally{baz()}}")
	// Try it out with switch statements
	c.FoldSame("function f(){switch(x){case 1:break}}")
	// Do while loops stay in a block if that's where they started
	c.FoldSame("function f(){if(e1){do foo();while(e2)}else foo2()}")
	// Test an obscure case with do and while
	c.Fold("if(x){do{foo()}while(y)}else bar()", "if(x){do foo();while(y)}else bar()")
	// Play with nested IFs
	c.Fold("function f(){if(x){if(y)foo()}}", "function f(){x && (y && foo())}")
	c.Fold("function f(){if(x){if(y)foo();else bar()}}", "function f(){x&&(y?foo():bar())}")
	c.Fold("function f(){if(x){if(y)foo()}else bar()}", "function f(){x?y&&foo():bar()}")
	c.Fold(
		"function f(){if(x){if(y)foo();else bar()}else{baz()}}",
		"function f(){x?y?foo():bar():baz()}")
	c.Fold("if(e1){while(e2){if(e3){foo()}}}else{bar()}", "if(e1)while(e2)e3&&foo();else bar()")
	c.Fold("if(e1){with(e2){if(e3){foo()}}}else{bar()}", "if(e1)with(e2)e3&&foo();else bar()")
	c.Fold("if(a||b){if(c||d){var x;}}", "if(a||b)if(c||d)var x")
	c.Fold("if(x){ if(y){var x;}else{var z;} }", "if(x)if(y)var x;else var z")
	// NOTE - technically we can remove the blocks since both the parent
	// and child have elses. But we don't since it causes ambiguities in
	// some cases where not all descendent ifs having elses
	c.Fold("if(x){ if(y){var x;}else{var z;} }else{var w}", "if(x)if(y)var x;else var z;else var w")
	c.Fold("if (x) {var x;}else { if (y) { var y;} }", "if(x)var x;else if(y)var y")
	// Here's some of the ambiguous cases
	c.Fold(
		"if(a){if(b){f1();f2();}else if(c){f3();}}else {if(d){f4();}}",
		"if(a)if(b){f1();f2()}else c&&f3();else d&&f4()")
	c.Fold("function f(){foo()}", "function f(){foo()}")
	c.Fold("switch(x){case y: foo()}", "switch(x){case y:foo()}")
	c.Fold("try{foo()}catch(ex){bar()}finally{baz()}", "try{foo()}catch(ex){bar()}finally{baz()}")
}

// Try to minimize returns
// Java: PeepholeMinimizeConditionsTest.testFoldReturns()
func (c *PeepholeMinimizeConditionsTest) TestFoldReturns() {
	c.Fold("function f(){if(x)return 1;else return 2}", "function f(){return x?1:2}")
	c.Fold("function f(){if(x)return 1;return 2}", "function f(){return x?1:2}")
	c.Fold("function f(){if(x)return;return 2}", "function f(){return x?void 0:2}")
	c.Fold("function f(){if(x)return 1+x;else return 2-x}", "function f(){return x?1+x:2-x}")
	c.Fold("function f(){if(x)return 1+x;return 2-x}", "function f(){return x?1+x:2-x}")
	c.Fold(
		"function f(){if(x)return y += 1;else return y += 2}",
		"function f(){return x?(y+=1):(y+=2)}")
	c.Fold("function f(){if(x)return;else return 2-x}", "function f(){if(x);else return 2-x}")
	c.Fold("function f(){if(x)return;return 2-x}", "function f(){return x?void 0:2-x}")
	c.Fold("function f(){if(x)return x;else return}", "function f(){if(x)return x;{}}")
	c.Fold("function f(){if(x)return x;return}", "function f(){if(x)return x}")
	c.FoldSame("function f(){for(var x in y) { return x.y; } return k}")
}

// Java: PeepholeMinimizeConditionsTest.testCombineIfs1()
func (c *PeepholeMinimizeConditionsTest) TestCombineIfs1() {
	c.Fold("function f() {if (x) return 1; if (y) return 1}", "function f() {if (x||y) return 1;}")
	c.Fold(
		"function f() {if (x) return 1; if (y) foo(); else return 1}",
		"function f() {if ((!x)&&y) foo(); else return 1;}")
}

// Java: PeepholeMinimizeConditionsTest.testCombineIfs2()
func (c *PeepholeMinimizeConditionsTest) TestCombineIfs2() {
	// combinable but not yet done
	c.FoldSame("function f() {if (x) throw 1; if (y) throw 1}")
	// Can't combine, side-effect
	c.Fold("function f(){ if (x) g(); if (y) g() }", "function f(){ x&&g(); y&&g() }")
	c.Fold("function f(){ if (x) g?.(); if (y) g?.() }", "function f(){ x&&g?.(); y&&g?.() }")
	// Can't combine, side-effect
	c.Fold("function f(){ if (x) y = 0; if (y) y = 0; }", "function f(){ x&&(y = 0); y&&(y = 0); }")
}

// Java: PeepholeMinimizeConditionsTest.testCombineIfs3()
func (c *PeepholeMinimizeConditionsTest) TestCombineIfs3() {
	c.FoldSame("function f() {if (x) return 1; if (y) {g();f()}}")
}

// Try to minimize assignments
// Java: PeepholeMinimizeConditionsTest.testFoldAssignments()
func (c *PeepholeMinimizeConditionsTest) TestFoldAssignments() {
	c.Fold("function f(){if(x)y=3;else y=4;}", "function f(){y=x?3:4}")
	c.Fold("function f(){if(x)y=1+a;else y=2+a;}", "function f(){y=x?1+a:2+a}")
	// and operation assignments
	c.Fold("function f(){if(x)y+=1;else y+=2;}", "function f(){y+=x?1:2}")
	c.Fold("function f(){if(x)y-=1;else y-=2;}", "function f(){y-=x?1:2}")
	c.Fold("function f(){if(x)y%=1;else y%=2;}", "function f(){y%=x?1:2}")
	c.Fold("function f(){if(x)y|=1;else y|=2;}", "function f(){y|=x?1:2}")
	// Don't fold if the 2 ops don't match.
	c.FoldSame("function f(){x ? y-=1 : y+=2}")
	// Don't fold if the 2 LHS don't match.
	c.FoldSame("function f(){x ? y-=1 : z-=1}")
	// Don't fold if there are potential effects.
	c.FoldSame("function f(){x ? y().a=3 : y().a=4}")
}

// Java: PeepholeMinimizeConditionsTest.testRemoveDuplicateStatements()
func (c *PeepholeMinimizeConditionsTest) TestRemoveDuplicateStatements() {
	c.EnableNormalize()
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.Fold("if (a) { x = 1; x++ } else { x = 2; x++ }", "x=(a) ? 1 : 2; x++")
	c.Fold(
		"if (a) { x = 1; x++; y += 1; z = pi; }\n"+
			" else  { x = 2; x++; y += 1; z = pi; }\n",
		"x=(a) ? 1 : 2; x++; y += 1; z = pi;")
	c.Fold(
		"function z() {\n"+
			"if (a) { foo(); return !0 } else { goo(); return !0 }\n"+
			"}\n",
		"function z() {(a) ? foo() : goo(); return !0}")
	c.Fold(
		"function z() {if (a) { foo(); x = true; return true\n"+
			"} else { goo(); x = true; return true }}\n",
		"function z() {(a) ? foo() : goo(); x = true; return true}")
	c.Fold(
		"function z() {\n"+
			"  if (a) { bar(); foo(); return true }\n"+
			"    else { bar(); goo(); return true }\n"+
			"}\n",
		"function z() {\n"+
			"  if (a) { bar(); foo(); }\n"+
			"    else { bar(); goo(); }\n"+
			"  return true;\n"+
			"}\n")
}

// Java: PeepholeMinimizeConditionsTest.testFoldReturnsIntegration2()
func (c *PeepholeMinimizeConditionsTest) TestFoldReturnsIntegration2() {
	c.late = true
	c.DisableNormalize()
	// if-then-else duplicate statement removal handles this case:
	c.TestSame("function test(a) {if (a) {const a = Math.random();if(a) {return a;}} return a; }")
}

// Java: PeepholeMinimizeConditionsTest.testDontRemoveDuplicateStatementsWithoutNormalization()
func (c *PeepholeMinimizeConditionsTest) TestDontRemoveDuplicateStatementsWithoutNormalization() {
	// In the following test case, we can't remove the duplicate "alert(x);" lines since each "x"
	// refers to a different variable.
	// We only try removing duplicate statements if the AST is normalized and names are unique.
	c.TestSame("if (Math.random() < 0.5) { const x = 3; alert(x); } else { const x = 5; alert(x); }")
}

// Java: PeepholeMinimizeConditionsTest.testDontFoldDirectAndIndirectEval()
func (c *PeepholeMinimizeConditionsTest) TestDontFoldDirectAndIndirectEval() {
	c.Fold("if (a) { eval(x); } else { (0, eval)(x); }", "a ? eval(x) : (0, eval)(x);")
}

// Java: PeepholeMinimizeConditionsTest.testNotCond()
func (c *PeepholeMinimizeConditionsTest) TestNotCond() {
	c.Fold("function f(){if(!x)foo()}", "function f(){x||foo()}")
	c.Fold("function f(){if(!x)b=1}", "function f(){x||(b=1)}")
	c.Fold("if(!x)z=1;else if(y)z=2", "if(x){y&&(z=2);}else{z=1;}")
	c.Fold("if(x)y&&(z=2);else z=1;", "x ? y&&(z=2) : z=1")
	c.Fold("function f(){if(!(x=1))a.b=1}", "function f(){(x=1)||(a.b=1)}")
}

// Java: PeepholeMinimizeConditionsTest.testAndParenthesesCount()
func (c *PeepholeMinimizeConditionsTest) TestAndParenthesesCount() {
	c.Fold("function f(){if(x||y)a.foo()}", "function f(){(x||y)&&a.foo()}")
	c.Fold("function f(){if(x.a)x.a=0}", "function f(){x.a&&(x.a=0)}")
	c.Fold("function f(){if(x?.a)x.a=0}", "function f(){x?.a&&(x.a=0)}")
	c.FoldSame("function f(){if(x()||y()){x()||y()}}")
}

// Java: PeepholeMinimizeConditionsTest.testFoldLogicalOpStringCompare()
func (c *PeepholeMinimizeConditionsTest) TestFoldLogicalOpStringCompare() {
	// side-effects
	// There is two way to parse two &&'s and both are correct.
	c.Fold("if (foo() && false) z()", "(foo(), 0) && z()")
}

// Java: PeepholeMinimizeConditionsTest.testFoldNot()
func (c *PeepholeMinimizeConditionsTest) TestFoldNot() {
	c.Fold("while(!(x==y)){a=b;}", "while(x!=y){a=b;}")
	c.Fold("while(!(x!=y)){a=b;}", "while(x==y){a=b;}")
	c.Fold("while(!(x===y)){a=b;}", "while(x!==y){a=b;}")
	c.Fold("while(!(x!==y)){a=b;}", "while(x===y){a=b;}")
	// Because !(x<NaN) != x>=NaN don't fold < and > cases.
	c.FoldSame("while(!(x>y)){a=b;}")
	c.FoldSame("while(!(x>=y)){a=b;}")
	c.FoldSame("while(!(x<y)){a=b;}")
	c.FoldSame("while(!(x<=y)){a=b;}")
	c.FoldSame("while(!(x<=NaN)){a=b;}")

	// NOT forces a boolean context
	c.Fold("x = !(y() && true)", "x = !y()")
	// This will be further optimized by PeepholeFoldConstants.
	c.Fold("x = !true", "x = !1")
	c.Fold("x = !!(a < b)", "x = a < b")
	c.FoldSame("x = !!a")
}

// Java: PeepholeMinimizeConditionsTest.testMinimizeExprCondition()
func (c *PeepholeMinimizeConditionsTest) TestMinimizeExprCondition() {
	c.Fold("(x ? true : false) && y()", "x&&y()")
	c.Fold("(x ? false : true) && y()", "(!x)&&y()")
	c.Fold("(x ? true : y) && y()", "(x || y)&&y()")
	c.Fold("(x ? y : false) && y()", "(x && y)&&y()")
	c.Fold("(x && true) && y()", "x && y()")
	c.Fold("(x && false) && y()", "0&&y()")
	c.Fold("(x || true) && y()", "1&&y()")
	c.Fold("(x || false) && y()", "x&&y()")
}

// Java: PeepholeMinimizeConditionsTest.testMinimizeWhileCondition()
func (c *PeepholeMinimizeConditionsTest) TestMinimizeWhileCondition() {
	// This test uses constant folding logic, so is only here for completeness.
	c.Fold("while(!!true) foo()", "while(1) foo()")
	// These test tryMinimizeCondition
	c.Fold("while(!!x) foo()", "while(x) foo()")
	c.Fold("while(!(!x&&!y)) foo()", "while(x||y) foo()")
	c.Fold("while(x||!!y) foo()", "while(x||y) foo()")
	c.Fold("while(!(!!x&&y)) foo()", "while(!x||!y) foo()")
	c.Fold("while(!(!x&&y)) foo()", "while(x||!y) foo()")
	c.Fold("while(!(x||!y)) foo()", "while(!x&&y) foo()")
	c.Fold("while(!(x||y)) foo()", "while(!x&&!y) foo()")
	c.Fold("while(!(!x||y-z)) foo()", "while(x&&!(y-z)) foo()")
	c.Fold("while(!(!(x/y)||z+w)) foo()", "while(x/y&&!(z+w)) foo()")
	c.FoldSame("while(!(x+y||z)) foo()")
	c.FoldSame("while(!(x&&y*z)) foo()")
	c.Fold("while(!(!!x&&y)) foo()", "while(!x||!y) foo()")
	c.Fold("while(x&&!0) foo()", "while(x) foo()")
	c.Fold("while(x||!1) foo()", "while(x) foo()")
	c.Fold("while(!((x,y)&&z)) foo()", "while((x,!y)||!z) foo()")
}

// Java: PeepholeMinimizeConditionsTest.testMinimizeDemorganRemoveLeadingNot()
func (c *PeepholeMinimizeConditionsTest) TestMinimizeDemorganRemoveLeadingNot() {
	c.Fold("if(!(!a||!b)&&c) foo()", "((a&&b)&&c)&&foo()")
	c.Fold("if(!(x&&y)) foo()", "x&&y||foo()")
	c.Fold("if(!(x||y)) foo()", "(x||y)||foo()")
}

// Java: PeepholeMinimizeConditionsTest.testMinimizeDemorgan1()
func (c *PeepholeMinimizeConditionsTest) TestMinimizeDemorgan1() {
	c.Fold("if(!a&&!b)foo()", "(a||b)||foo()")
}

// Java: PeepholeMinimizeConditionsTest.testMinimizeDemorgan2()
func (c *PeepholeMinimizeConditionsTest) TestMinimizeDemorgan2() {
	// Make sure trees with cloned functions are marked as changed
	c.Fold("(!(a&&!((function(){})())))||foo()", "!a||(function(){})()||foo()")
}

// Java: PeepholeMinimizeConditionsTest.testMinimizeDemorgan2b()
func (c *PeepholeMinimizeConditionsTest) TestMinimizeDemorgan2b() {
	// Make sure unchanged trees with functions are not marked as changed
	c.FoldSame("!a||(function(){})()||foo()")
}

// Java: PeepholeMinimizeConditionsTest.testMinimizeDemorgan3()
func (c *PeepholeMinimizeConditionsTest) TestMinimizeDemorgan3() {
	c.Fold("if((!a||!b)&&(c||d)) foo()", "(a&&b||!c&&!d)||foo()")
}

// Java: PeepholeMinimizeConditionsTest.testMinimizeDemorgan5()
func (c *PeepholeMinimizeConditionsTest) TestMinimizeDemorgan5() {
	c.Fold("if((!a||!b)&&c) foo()", "(a&&b||!c)||foo()")
}

// Java: PeepholeMinimizeConditionsTest.testMinimizeDemorgan11()
func (c *PeepholeMinimizeConditionsTest) TestMinimizeDemorgan11() {
	c.Fold(
		"if (x && (y===2 || !f()) && (y===3 || !h())) foo()",
		"(!x || y!==2 && f() || y!==3 && h()) || foo()")
}

// Java: PeepholeMinimizeConditionsTest.testMinimizeDemorgan20a()
func (c *PeepholeMinimizeConditionsTest) TestMinimizeDemorgan20a() {
	c.Fold(
		"if (0===c && (2===a || 1===a)) f(); else g()",
		"if (0!==c || 2!==a && 1!==a) g(); else f()")
}

// Java: PeepholeMinimizeConditionsTest.testMinimizeDemorgan20b()
func (c *PeepholeMinimizeConditionsTest) TestMinimizeDemorgan20b() {
	c.Fold("if (0!==c || 2!==a && 1!==a) g(); else f()", "(0!==c || 2!==a && 1!==a) ? g() : f()")
}

// Java: PeepholeMinimizeConditionsTest.testPreserveIf()
func (c *PeepholeMinimizeConditionsTest) TestPreserveIf() {
	c.FoldSame("if(!a&&!b)for(;f(););")
}

// Java: PeepholeMinimizeConditionsTest.testNoSwapWithDanglingElse()
func (c *PeepholeMinimizeConditionsTest) TestNoSwapWithDanglingElse() {
	c.FoldSame("if(!x) {for(;;)foo(); for(;;)bar()} else if(y) for(;;) f()")
	c.FoldSame("if(!a&&!b) {for(;;)foo(); for(;;)bar()} else if(y) for(;;) f()")
}

// Java: PeepholeMinimizeConditionsTest.testMinimizeHook()
func (c *PeepholeMinimizeConditionsTest) TestMinimizeHook() {
	c.Fold("x ? x : y", "x || y")
	// We assume GETPROPs don't have side effects.
	c.Fold("x.y ? x.y : x.z", "x.y || x.z")
	c.Fold("x?.y ? x?.y : x.z", "x?.y || x.z")
	c.Fold("x?.y ? x?.y : x?.z", "x?.y || x?.z")

	// This can be folded if x() does not have side effects.
	c.FoldSame("x() ? x() : y()")
	c.FoldSame("x?.() ? x?.() : y()")

	c.Fold("!x ? foo() : bar()", "x ? bar() : foo()")
	c.Fold("while(!(x ? y : z)) foo();", "while(x ? !y : !z) foo();")
	c.Fold("(x ? !y : !z) ? foo() : bar()", "(x ? y : z) ? bar() : foo()")

	c.Fold("y = x ? true : false", "y = !!x")
	c.Fold("y = x ? !0 : !1", "y = !!x")
	c.Fold("y = x ? false : true", "y = !x")
	c.Fold("y = x ? !1 : !0", "y = !x")
	c.Fold("y = !x ? true : false", "y = !x")
	c.Fold("y = !x ? false : true", "y = !!x")
	c.Fold("y = (a === b) ? true : false", "y = a === b")
	c.Fold("y = (a === b) ? false : true", "y = a !== b")
	c.Fold("y = !(a === b) ? true : false", "y = a !== b")
	c.Fold("y = (a < b) ? true : false", "y = a < b")
	c.Fold("y = (a < b) ? false : true", "y = !(a < b)")
	c.Fold("y = !(a < b) ? true : false", "y = !(a < b)")
	c.Fold("y = !(a < b) ? false : true", "y = a < b")
	c.Fold("y = !!(a < b) ? true : false", "y = a < b")
	c.Fold("y = (!a || !b) ? true : false", "y = !a || !b")
	c.Fold("y = !(a && !b) ? true : false", "y = !(a && !b)")
	c.Fold("y = !(a && !b) ? false : true", "y = !!(a && !b)")
	c.Fold("y = !(a === 1 && !b) ? false : true", "y = a === 1 && !b")
	c.Fold("y = (inc(), !test()) ? true : false", "y = (inc(), !test())")
	c.Fold("y = x ? (function(){}, true) : false", "y = !!x")
	c.Fold("y = (a === b) ? (function(){}, true) : false", "y = a === b")
	c.FoldSame("y = x ? 1 : 0")
	c.FoldSame("y = x ? 'a' : ''")
}

// Java: PeepholeMinimizeConditionsTest.testMinimizeComma()
func (c *PeepholeMinimizeConditionsTest) TestMinimizeComma() {
	c.Fold("while(!(inc(), test())) foo();", "while(inc(), !test()) foo();")
	c.Fold("(inc(), !test()) ? foo() : bar()", "(inc(), test()) ? bar() : foo()")
}

// Java: PeepholeMinimizeConditionsTest.testMinimizeExprResult()
func (c *PeepholeMinimizeConditionsTest) TestMinimizeExprResult() {
	c.Fold("!x||!y", "x&&y")
	c.Fold("if(!(x&&!y)) foo()", "(!x||y)&&foo()")
	c.Fold("if(!x||y) foo()", "(!x||y)&&foo()")
	c.Fold("(!x||y)&&foo()", "x&&!y||!foo()")
}

// Java: PeepholeMinimizeConditionsTest.testMinimizeDemorgan21()
func (c *PeepholeMinimizeConditionsTest) TestMinimizeDemorgan21() {
	c.Fold("if (0===c && (2===a || 1===a)) f()", "(0!==c || 2!==a && 1!==a) || f()")
}

// Java: PeepholeMinimizeConditionsTest.testMinimizeAndOr1()
func (c *PeepholeMinimizeConditionsTest) TestMinimizeAndOr1() {
	c.Fold("if ((!a || !b) && (d || e)) f()", "(a&&b || !d&&!e) || f()")
}

// Java: PeepholeMinimizeConditionsTest.testMinimizeForCondition()
func (c *PeepholeMinimizeConditionsTest) TestMinimizeForCondition() {
	// This test uses constant folding logic, so is only here for completeness.
	// These could be simplified to "for(;;) ..."
	c.Fold("for(;!!true;) foo()", "for(;1;) foo()")
	// Verify function deletion tracking.
	c.Fold("if(!!true||function(){}) {}", "if(1) {}")
	// Don't bother with FOR inits as there are normalized out.
	c.Fold("for(!!true;;) foo()", "for(!0;;) foo()")
	// These test tryMinimizeCondition
	c.Fold("for(;!!x;) foo()", "for(;x;) foo()")
	c.FoldSame("for(a in b) foo()")
	c.FoldSame("for(a in {}) foo()")
	c.FoldSame("for(a in []) foo()")
	c.Fold("for(a in !!true) foo()", "for(a in !0) foo()")
	c.FoldSame("for(a of b) foo()")
	c.FoldSame("for(a of {}) foo()")
	c.FoldSame("for(a of []) foo()")
	c.Fold("for(a of !!true) foo()", "for(a of !0) foo()")
}

// Java: PeepholeMinimizeConditionsTest.testMinimizeCondition_example1()
func (c *PeepholeMinimizeConditionsTest) TestMinimizeCondition_example1() {
	// Based on a real failing code sample.
	c.Fold("if(!!(f() > 20)) {foo();foo()}", "if(f() > 20){foo();foo()}")
}

// Java: PeepholeMinimizeConditionsTest.testFoldLoopBreakLate()
func (c *PeepholeMinimizeConditionsTest) TestFoldLoopBreakLate() {
	c.late = true
	c.Fold("for(;;) if (a) break", "for(;!a;);")
	c.FoldSame("for(;;) if (a) { f(); break }")
	c.Fold("for(;;) if (a) break; else f()", "for(;!a;) { { f(); } }")
	c.Fold("for(;a;) if (b) break", "for(;a && !b;);")
	c.Fold("for(;a;) { if (b) break; if (c) break; }", "for(;(a && !b);) if (c) break;")
	c.Fold("for(;(a && !b);) if (c) break;", "for(;(a && !b) && !c;);")
	c.Fold("for(;;) { if (foo) { break; var x; } } x;", "var x; for(;!foo;) {} x;")
	// 'while' is normalized to 'for'
	c.EnableNormalize()
	c.Fold("while(true) if (a) break", "for(;1&&!a;);")
	c.DisableNormalize()
}

// Java: PeepholeMinimizeConditionsTest.testFoldLoopBreakEarly()
func (c *PeepholeMinimizeConditionsTest) TestFoldLoopBreakEarly() {
	c.late = false
	c.FoldSame("for(;;) if (a) break")
	c.FoldSame("for(;;) if (a) { f(); break }")
	c.FoldSame("for(;;) if (a) break; else f()")
	c.FoldSame("for(;a;) if (b) break")
	c.FoldSame("for(;a;) { if (b) break; if (c) break; }")
	c.FoldSame("while(1) if (a) break")
	c.EnableNormalize()
	c.FoldSame("for (; 1; ) if (a) break")
}

// Java: PeepholeMinimizeConditionsTest.testFoldConditionalVarDeclaration()
func (c *PeepholeMinimizeConditionsTest) TestFoldConditionalVarDeclaration() {
	c.Fold("if(x) var y=1;else y=2", "var y=x?1:2")
	c.Fold("if(x) y=1;else var y=2", "var y=x?1:2")
	c.FoldSame("if(x) var y = 1; z = 2")
	c.FoldSame("if(x||y) y = 1; var z = 2")
	c.FoldSame("if(x) { var y = 1; print(y)} else y = 2 ")
	c.FoldSame("if(x) var y = 1; else {y = 2; print(y)}")
}

// Java: PeepholeMinimizeConditionsTest.testFoldIfWithLowerOperatorsInside()
func (c *PeepholeMinimizeConditionsTest) TestFoldIfWithLowerOperatorsInside() {
	c.Fold("if (x + (y=5)) z && (w,z);", "x + (y=5) && (z && (w,z))")
	c.Fold("if (!(x+(y=5))) z && (w,z);", "x + (y=5) || z && (w,z)")
	c.Fold(
		"if (x + (y=5)) if (z && (w,z)) for(;;) foo();",
		"if (x + (y=5) && (z && (w,z))) for(;;) foo();")
}

// Java: PeepholeMinimizeConditionsTest.testSubsituteReturn()
func (c *PeepholeMinimizeConditionsTest) TestSubsituteReturn() {
	c.late = false
	c.EnableNormalize()
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()

	c.Fold("function f() { while(x) { return }}", "function f() { while(x) { break }}")

	c.FoldSame("function f() { while(x) { return 5 } }")

	c.FoldSame("function f() { a: { return 5 } }")

	c.Fold(
		"function f() { while(x) { return 5}  return 5}",
		"function f() { while(x) { break }    return 5}")

	c.Fold(
		"function f() { while(x) { return x}  return x}",
		"function f() { while(x) { break }    return x}")

	c.Fold(
		"function f() { while(x) { if (y) { return }}}",
		"function f() { while(x) { if (y) { break  }}}")

	c.Fold(
		"function f() { while(x) { if (y) { return }} return}",
		"function f() { while(x) { if (y) { break  }}}")

	c.Fold(
		"function f() { while(x) { if (y) { return 5 }} return 5}",
		"function f() { while(x) { if (y) { break    }} return 5}")

	// It doesn't matter if x is changed between them. We are still returning
	// x at whatever x value current holds. The whole x = 1 is skipped.
	c.Fold(
		"function f() { while(x) { if (y) { return x } x = 1} return x}",
		"function f() { while(x) { if (y) { break    } x = 1} return x}")

	c.Fold(
		"function f() { while(x) { if (y) { return x } return x} return x}",
		"function f() { while(x) { if (y) {} break }return x}")

	// A break here only breaks out of the inner loop.
	c.FoldSame("function f() { while(x) { while (y) { return } } }")

	c.FoldSame("function f() { while(1) { return 7}  return 5}")

	c.FoldSame(
		"function f() {\n" +
			"  try { while(x) {return f()}} catch (e) { } return f()}\n")

	c.FoldSame(
		"function f() {\n" +
			"  try { while(x) {return f()}} finally {alert(1)} return f()}\n")

	// Both returns has the same handler
	c.Fold(
		"function f() {\n"+
			"  try { while(x) { return f() } return f() } catch (e) { } }\n",
		"function f() {\n"+
			"  try { while(x) { break } return f() } catch (e) { } }\n")

	// We can't fold this because it'll change the order of when foo is called.
	c.FoldSame(
		"function f() {\n" +
			"  try { while(x) { return foo() } } finally { alert(1) }\n" +
			"  return foo()}\n")

	// This is fine, we have no side effect in the return value.
	c.Fold(
		"function f() {\n"+
			"  try { while(x) { return 1 } } finally { alert(1) } return 1}\n",
		"function f() {\n"+
			"  try { while(x) { break    } } finally { alert(1) } return 1}\n")

	c.FoldSame("function f() { try{ return a } finally { a = 2 } return a; }")

	c.Fold(
		"function f() { switch(a){ case 1: return a; default: g();} return a;}",
		"function f() { switch(a){ case 1: break; default: g();} return a; }")
	c.FoldSame(
		"function f(flag) {\n" +
			"  var ok = false;\n" +
			"  while (flag) {\n" +
			"    try { return ok; } finally { ok = true; }\n" +
			"  }\n" +
			"  return ok;\n" +
			"}\n")

	c.Fold(
		"function f(flag) {\n"+
			"  var ok = false;\n"+
			"  while (flag) {\n"+
			"    try { return 1; } finally { ok = true; }\n"+
			"  }\n"+
			"  return 1;\n"+
			"}\n",
		"function f(flag) {\n"+
			"  var ok = false;\n"+
			"  while (flag) {\n"+
			"    try { break; } finally { ok = true; }\n"+
			"  }\n"+
			"  return 1;\n"+
			"}\n")
}

// Java: PeepholeMinimizeConditionsTest.testSubsituteBreakForThrow()
func (c *PeepholeMinimizeConditionsTest) TestSubsituteBreakForThrow() {
	c.late = false
	c.EnableNormalize()
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()

	c.FoldSame("function f() { while(x) { throw Error }}")

	c.Fold(
		"function f() { while(x) { throw Error } throw Error }",
		"function f() { while(x) { break } throw Error}")
	c.FoldSame("function f() { while(x) { throw Error(1) } throw Error(2)}")
	c.FoldSame("function f() { while(x) { throw Error(1) } return Error(2)}")

	c.FoldSame("function f() { while(x) { throw 5 } }")

	c.FoldSame("function f() { a: { throw 5 } }")

	c.Fold(
		"function f() { while(x) { throw 5}  throw 5}",
		"function f() { while(x) { break }   throw 5}")

	c.Fold(
		"function f() { while(x) { throw x}  throw x}",
		"function f() { while(x) { break }   throw x}")

	c.FoldSame("function f() { while(x) { if (y) { throw Error }}}")

	c.Fold(
		"function f() { while(x) { if (y) { throw Error }} throw Error}",
		"function f() { while(x) { if (y) { break }} throw Error}")

	c.Fold(
		"function f() { while(x) { if (y) { throw 5 }} throw 5}",
		"function f() { while(x) { if (y) { break    }} throw 5}")

	// It doesn't matter if x is changed between them. We are still throwing
	// x at whatever x value current holds. The whole x = 1 is skipped.
	c.Fold(
		"function f() { while(x) { if (y) { throw x } x = 1} throw x}",
		"function f() { while(x) { if (y) { break    } x = 1} throw x}")

	c.Fold(
		"function f() { while(x) { if (y) { throw x } throw x} throw x}",
		"function f() { while(x) { if (y) {} break }throw x}")

	// A break here only breaks out of the inner loop.
	c.FoldSame("function f() { while(x) { while (y) { throw Error } } }")

	c.FoldSame("function f() { while(1) { throw 7}  throw 5}")

	c.FoldSame(
		"function f() {\n" +
			"  try { while(x) {throw f()}} catch (e) { } throw f()}\n")

	c.FoldSame(
		"function f() {\n" +
			"  try { while(x) {throw f()}} finally {alert(1)} throw f()}\n")

	// Both throws has the same handler
	c.Fold(
		"function f() {\n"+
			"  try { while(x) { throw f() } throw f() } catch (e) { } }\n",
		"function f() {\n"+
			"  try { while(x) { break } throw f() } catch (e) { } }\n")

	// We can't fold this because it'll change the order of when foo is called.
	c.FoldSame(
		"function f() {\n" +
			"  try { while(x) { throw foo() } } finally { alert(1) }\n" +
			"  throw foo()}\n")

	// This is fine, we have no side effect in the throw value.
	c.Fold(
		"function f() {\n"+
			"  try { while(x) { throw 1 } } finally { alert(1) } throw 1}\n",
		"function f() {\n"+
			"  try { while(x) { break    } } finally { alert(1) } throw 1}\n")

	c.FoldSame("function f() { try{ throw a } finally { a = 2 } throw a; }")

	c.Fold(
		"function f() { switch(a){ case 1: throw a; default: g();} throw a;}",
		"function f() { switch(a){ case 1: break; default: g();} throw a; }")
	c.FoldSame(
		"function f(flag, e) {\n" +
			"  while (flag) {\n" +
			"    try { throw e; } finally { e = null; }\n" +
			"  }\n" +
			"  throw e;\n" +
			"}\n")
}

// Java: PeepholeMinimizeConditionsTest.testRemoveDuplicateReturn()
func (c *PeepholeMinimizeConditionsTest) TestRemoveDuplicateReturn() {
	c.late = false
	c.EnableNormalize()
	c.Fold("function f() { return; }", "function f(){}")
	c.FoldSame("function f() { return a; }")
	c.Fold("function f() { if (x) { return a } return a; }", "function f() { if (x) {} return a; }")
	c.FoldSame("function f() { try { if (x) { return a } } catch(e) {} return a; }")
	c.FoldSame("function f() { try { if (x) {} } catch(e) {} return 1; }")
	// finally clauses may have side effects
	c.FoldSame("function f() { try { if (x) { return a } } finally { a++ } return a; }")
	// but they don't matter if the result doesn't have side effects and can't
	// be affect by side-effects.
	c.Fold(
		"function f() { try { if (x) { return 1 } } finally {} return 1; }",
		"function f() { try { if (x) {} } finally {} return 1; }")
	c.Fold(
		"function f() { switch(a){ case 1: return a; } return a; }",
		"function f() { switch(a){ case 1: } return a; }")
	c.Fold(
		"function f() { switch(a){\n"+
			"  case 1: return a; case 2: return a; } return a; }\n",
		"function f() { switch(a){\n"+
			"  case 1: break; case 2: } return a; }\n")
}

// Java: PeepholeMinimizeConditionsTest.testRemoveDuplicateThrow()
func (c *PeepholeMinimizeConditionsTest) TestRemoveDuplicateThrow() {
	c.late = false
	c.EnableNormalize()
	c.FoldSame("function f() { throw a; }")
	c.Fold("function f() { if (x) { throw a } throw a; }", "function f() { if (x) {} throw a; }")
	c.FoldSame("function f() { try { if (x) {throw a} } catch(e) {} throw a; }")
	c.FoldSame("function f() { try { if (x) {throw 1} } catch(e) {f()} throw 1; }")
	c.FoldSame("function f() { try { if (x) {throw 1} } catch(e) {f()} throw 1; }")
	c.FoldSame("function f() { try { if (x) {throw 1} } catch(e) {throw 1}}")
	c.Fold(
		"function f() { try { if (x) {throw 1} } catch(e) {throw 1} throw 1; }",
		"function f() { try { if (x) {throw 1} } catch(e) {} throw 1; }")
	// finally clauses may have side effects
	c.FoldSame("function f() { try { if (x) { throw a } } finally { a++ } throw a; }")
	// but they don't matter if the result doesn't have side effects and can't
	// be affect by side-effects.
	c.Fold(
		"function f() { try { if (x) { throw 1 } } finally {} throw 1; }",
		"function f() { try { if (x) {} } finally {} throw 1; }")
	c.Fold(
		"function f() { switch(a){ case 1: throw a; } throw a; }",
		"function f() { switch(a){ case 1: } throw a; }")
	c.Fold(
		"function f() { switch(a){\n"+
			"case 1: throw a; case 2: throw a; } throw a; }\n",
		"function f() { switch(a){ case 1: break; case 2: } throw a; }")
}

// Java: PeepholeMinimizeConditionsTest.testNestedIfCombine()
func (c *PeepholeMinimizeConditionsTest) TestNestedIfCombine() {
	c.Fold("if(x)if(y){while(1){}}", "if(x&&y){while(1){}}")
	c.Fold("if(x||z)if(y){while(1){}}", "if((x||z)&&y){while(1){}}")
	c.Fold("if(x)if(y||z){while(1){}}", "if((x)&&(y||z)){while(1){}}")
	c.FoldSame("if(x||z)if(y||z){while(1){}}")
	c.Fold("if(x)if(y){if(z){while(1){}}}", "if(x&&(y&&z)){while(1){}}")
}

// See: http://blickly.github.io/closure-compiler-issues/#291
// Java: PeepholeMinimizeConditionsTest.testIssue291()
func (c *PeepholeMinimizeConditionsTest) TestIssue291() {
	c.Fold("if (true) { f.onchange(); }", "if (1) f.onchange();")
	c.FoldSame("if (f) { f.onchange(); }")
	c.FoldSame("if (f) { f.bar(); } else { f.onchange(); }")
	c.Fold("if (f) { f.bonchange(); }", "f && f.bonchange();")
	c.FoldSame("if (f) { f['x'](); }")
	// optional versions
	c.Fold("if (true) { f?.onchange(); }", "if (1) f?.onchange();")
	c.FoldSame("if (f) { f?.onchange(); }")
	c.FoldSame("if (f) { f?.bar(); } else { f?.onchange(); }")
	c.Fold("if (f) { f?.bonchange(); }", "f && f?.bonchange();")
	c.FoldSame("if (f) { f?.['x'](); }")
}

// Java: PeepholeMinimizeConditionsTest.testObjectLiteral()
func (c *PeepholeMinimizeConditionsTest) TestObjectLiteral() {
	c.Test("({})", "1")
	c.Test("({a:1})", "1")
	c.TestSame("({a:foo()})")
	c.TestSame("({'a':foo()})")
}

// Java: PeepholeMinimizeConditionsTest.testArrayLiteral()
func (c *PeepholeMinimizeConditionsTest) TestArrayLiteral() {
	c.Test("([])", "1")
	c.Test("([1])", "1")
	c.Test("([a])", "1")
	c.TestSame("([foo()])")
}

// Java: PeepholeMinimizeConditionsTest.testRemoveElseCause()
func (c *PeepholeMinimizeConditionsTest) TestRemoveElseCause() {
	c.Test(
		"function f() {\n"+
			" if(x) return 1;\n"+
			" else if(x) return 2;\n"+
			" else if(x) return 3 }\n",
		"function f() {\n"+
			" if(x) return 1;\n"+
			"{ if(x) return 2;\n"+
			"{ if(x) return 3 } } }\n")
}

// Java: PeepholeMinimizeConditionsTest.testRemoveElseCause1()
func (c *PeepholeMinimizeConditionsTest) TestRemoveElseCause1() {
	c.Test("function f() { if (x) throw 1; else f() }", "function f() { if (x) throw 1; { f() } }")
}

// Java: PeepholeMinimizeConditionsTest.testRemoveElseCause2()
func (c *PeepholeMinimizeConditionsTest) TestRemoveElseCause2() {
	c.Test("function f() { if (x) return 1; else f() }", "function f() { if (x) return 1; { f() } }")
	c.Test("function f() { if (x) return; else f() }", "function f() { if (x) {} else { f() } }")
	// This case is handled by minimize exit points.
	c.TestSame("function f() { if (x) return; f() }")
}

// Java: PeepholeMinimizeConditionsTest.testRemoveElseCause3()
func (c *PeepholeMinimizeConditionsTest) TestRemoveElseCause3() {
	c.TestSame("function f() { a:{if (x) break a; else f() } }")
	c.TestSame("function f() { if (x) { a:{ break a } } else f() }")
	c.TestSame("function f() { if (x) a:{ break a } else f() }")
}

// Java: PeepholeMinimizeConditionsTest.testRemoveElseCause4()
func (c *PeepholeMinimizeConditionsTest) TestRemoveElseCause4() {
	c.TestSame("function f() { if (x) { if (y) { return 1; } } else f() }")
}

// Java: PeepholeMinimizeConditionsTest.testIssue925()
func (c *PeepholeMinimizeConditionsTest) TestIssue925() {
	c.Test(
		"if (x[--y] === 1) {\n"+
			"    x[y] = 0;\n"+
			"} else {\n"+
			"    x[y] = 1;\n"+
			"}\n",
		"(x[--y] === 1) ? x[y] = 0 : x[y] = 1;")
	c.Test(
		"if (x[--y]) {\n"+
			"    a = 0;\n"+
			"} else {\n"+
			"    a = 1;\n"+
			"}\n",
		"a = (x[--y]) ? 0 : 1;")
	c.Test(
		"if (x?.[--y]) {\n"+
			"    a = 0;\n"+
			"} else {\n"+
			"    a = 1;\n"+
			"}\n",
		"a = (x?.[--y]) ? 0 : 1;")
	c.Test("if (x++) { x += 2 } else { x += 3 }", "x++ ? x += 2 : x += 3")
	c.Test("if (x++) { x = x + 2 } else { x = x + 3 }", "x = x++ ? x + 2 : x + 3")
}

// Java: PeepholeMinimizeConditionsTest.testCoercionSubstitution_disabled()
func (c *PeepholeMinimizeConditionsTest) TestCoercionSubstitution_disabled() {
	c.EnableTypeCheck()
	c.TestSame("var x = {}; if (x != null) throw 'a';")
	c.TestSame("var x = {}; var y = x != null;")
	c.TestSame("var x = 1; if (x != 0) throw 'a';")
	c.TestSame("var x = 1; var y = x != 0;")
}

// Java: PeepholeMinimizeConditionsTest.testCoercionSubstitution_booleanResult0()
func (c *PeepholeMinimizeConditionsTest) TestCoercionSubstitution_booleanResult0() {
	c.EnableTypeCheck()
	c.TestSame("var x = {}; var y = x != null;")
}

// Java: PeepholeMinimizeConditionsTest.testCoercionSubstitution_booleanResult1()
func (c *PeepholeMinimizeConditionsTest) TestCoercionSubstitution_booleanResult1() {
	c.EnableTypeCheck()
	c.TestSame("var x = {}; var y = x == null;")
	c.TestSame("var x = {}; var y = x !== null;")
	c.TestSame("var x = undefined; var y = x !== null;")
	c.TestSame("var x = {}; var y = x === null;")
	c.TestSame("var x = undefined; var y = x === null;")
	c.TestSame("var x = 1; var y = x != 0;")
	c.TestSame("var x = 1; var y = x == 0;")
	c.TestSame("var x = 1; var y = x !== 0;")
	c.TestSame("var x = 1; var y = x === 0;")
}

// Java: PeepholeMinimizeConditionsTest.testCoercionSubstitution_if()
func (c *PeepholeMinimizeConditionsTest) TestCoercionSubstitution_if() {
	c.EnableTypeCheck()
	c.Test("var x = {};\nif (x != null) throw 'a';\n", "var x={}; if (x!=null) throw 'a'")
	c.TestSame("var x = {};\nif (x == null) throw 'a';\n")
	c.TestSame("var x = {};\nif (x !== null) throw 'a';\n")
	c.TestSame("var x = {};\nif (x === null) throw 'a';\n")
	c.TestSame("var x = {};\nif (null != x) throw 'a';\n")
	c.TestSame("var x = {};\nif (null == x) throw 'a';\n")
	c.TestSame("var x = {};\nif (null !== x) throw 'a';\n")
	c.TestSame("var x = {};\nif (null === x) throw 'a';\n")
	c.TestSame("var x = 1;\nif (x != 0) throw 'a';\n")
	c.TestSame("var x = 1;\nif (x != 0) throw 'a';\n")
	c.TestSame("var x = 1;\nif (x == 0) throw 'a';\n")
	c.TestSame("var x = 1;\nif (x !== 0) throw 'a';\n")
	c.TestSame("var x = 1;\nif (x === 0) throw 'a';\n")
	c.TestSame("var x = 1;\nif (0 != x) throw 'a';\n")
	c.TestSame("var x = 1;\nif (0 == x) throw 'a';\n")
	c.TestSame("var x = 1;\nif (0 !== x) throw 'a';\n")
	c.TestSame("var x = 1;\nif (0 === x) throw 'a';\n")
	c.TestSame("var x = NaN;\nif (0 === x) throw 'a';\n")
	c.TestSame("var x = NaN;\nif (x === 0) throw 'a';\n")
}

// Java: PeepholeMinimizeConditionsTest.testCoercionSubstitution_expression()
func (c *PeepholeMinimizeConditionsTest) TestCoercionSubstitution_expression() {
	c.EnableTypeCheck()
	c.TestSame("var x = {}; x != null && alert('b');")
	c.TestSame("var x = 1; x != 0 && alert('b');")
}

// Java: PeepholeMinimizeConditionsTest.testCoercionSubstitution_hook()
func (c *PeepholeMinimizeConditionsTest) TestCoercionSubstitution_hook() {
	c.EnableTypeCheck()
	c.TestSame(
		"var x = {};\n" +
			"var y = x != null ? 1 : 2;\n")
	c.TestSame(
		"var x = 1;\n" +
			"var y = x != 0 ? 1 : 2;\n")
}

// Java: PeepholeMinimizeConditionsTest.testCoercionSubstitution_not()
func (c *PeepholeMinimizeConditionsTest) TestCoercionSubstitution_not() {
	c.EnableTypeCheck()
	c.Test(
		"var x = {};\nvar y = !(x != null) ? 1 : 2;\n",
		"var x = {};\nvar y = (x == null) ? 1 : 2;\n")
	c.Test("var x = 1;\nvar y = !(x != 0) ? 1 : 2;\n", "var x = 1;\nvar y = x == 0 ? 1 : 2;\n")
}

// Java: PeepholeMinimizeConditionsTest.testCoercionSubstitution_while()
func (c *PeepholeMinimizeConditionsTest) TestCoercionSubstitution_while() {
	c.EnableTypeCheck()
	c.TestSame("var x = {}; while (x != null) throw 'a';")
	c.TestSame("var x = 1; while (x != 0) throw 'a';")
}

// Java: PeepholeMinimizeConditionsTest.testCoercionSubstitution_unknownType()
func (c *PeepholeMinimizeConditionsTest) TestCoercionSubstitution_unknownType() {
	c.EnableTypeCheck()
	c.TestSame("var x = /** @type {?} */ ({});\nif (x != null) throw 'a';\n")
	c.TestSame("var x = /** @type {?} */ (1);\nif (x != 0) throw 'a';\n")
}

// Java: PeepholeMinimizeConditionsTest.testCoercionSubstitution_allType()
func (c *PeepholeMinimizeConditionsTest) TestCoercionSubstitution_allType() {
	c.EnableTypeCheck()
	c.TestSame("var x = /** @type {*} */ ({});\nif (x != null) throw 'a';\n")
	c.TestSame("var x = /** @type {*} */ (1);\nif (x != 0) throw 'a';\n")
}

// Java: PeepholeMinimizeConditionsTest.testCoercionSubstitution_primitivesVsNull()
func (c *PeepholeMinimizeConditionsTest) TestCoercionSubstitution_primitivesVsNull() {
	c.EnableTypeCheck()
	c.TestSame("var x = 0;\nif (x != null) throw 'a';\n")
	c.TestSame("var x = '';\nif (x != null) throw 'a';\n")
	c.TestSame("var x = false;\nif (x != null) throw 'a';\n")
}

// Java: PeepholeMinimizeConditionsTest.testCoercionSubstitution_nonNumberVsZero()
func (c *PeepholeMinimizeConditionsTest) TestCoercionSubstitution_nonNumberVsZero() {
	c.EnableTypeCheck()
	c.TestSame("var x = {};\nif (x != 0) throw 'a';\n")
	c.TestSame("var x = '';\nif (x != 0) throw 'a';\n")
	c.TestSame("var x = false;\nif (x != 0) throw 'a';\n")
}

// Java: PeepholeMinimizeConditionsTest.testCoercionSubstitution_boxedNumberVsZero()
func (c *PeepholeMinimizeConditionsTest) TestCoercionSubstitution_boxedNumberVsZero() {
	c.EnableTypeCheck()
	c.TestSame("var x = new Number(0);\nif (x != 0) throw 'a';\n")
}

// Java: PeepholeMinimizeConditionsTest.testCoercionSubstitution_boxedPrimitives()
func (c *PeepholeMinimizeConditionsTest) TestCoercionSubstitution_boxedPrimitives() {
	c.EnableTypeCheck()
	c.TestSame("var x = new Number(); if (x != null) throw 'a';")
	c.TestSame("var x = new String(); if (x != null) throw 'a';")
	c.TestSame("var x = new Boolean();\nif (x != null) throw 'a';")
}

// Java: PeepholeMinimizeConditionsTest.testMinimizeIfWithNewTargetCondition()
func (c *PeepholeMinimizeConditionsTest) TestMinimizeIfWithNewTargetCondition() {
	// Related to https://github.com/google/closure-compiler/issues/3097
	c.Test(
		"function x() {\n"+
			"  if (new.target) {\n"+
			"    return 1;\n"+
			"  } else {\n"+
			"    return 2;\n"+
			"  }\n"+
			"}\n",
		"function x() {\n"+
			"  return new.target ? 1 : 2;\n"+
			"}\n")
}

// Java: PeepholeMinimizeConditionsTest.testBatchE_booleanAbsorptionAndContextSimplifications()
func (c *PeepholeMinimizeConditionsTest) TestBatchE_booleanAbsorptionAndContextSimplifications() {
	// OPP-022: Boolean Absorption & Context Simplifications
	// Hook inversion on negated condition
	c.Test("var res = !x ? a : b;", "var res = x ? b : a;")

	// if-statement minimization preserves Boolean(x) wrapper in current pass
	c.Test("if (Boolean(x)) { foo(); }", "Boolean(x) && foo();")

	// Guard cases: non-constant / side-effectful / shadowed expressions
	c.TestSame("while (Boolean(x)) foo();")
	c.TestSame("var res = Boolean(x) ? a : b;")
	c.TestSame("var res = !Boolean(x);")
	c.TestSame("x = a || (a && b);")
	c.TestSame("x = a && (a || b);")
	c.TestSame("x = foo() || (foo() && b);")
	c.TestSame("function f(Boolean) { Boolean(x) ? a : b; }")
}

// Go-only: PeepholeOptimizationsPass's prefilter skips assignment statements
// that certainly change state, and must let the others through: Java folds
// a side-effect-free assignment to a property of a literal (the test262
// files pass/d96153b59454dddd.js and pass-explicit/687f678cde900411.js);
// the expectations are PeepholeDump's output with the v20261015 jar.
func (c *PeepholeMinimizeConditionsTest) TestGoOnlyPrefilterLetsSideEffectFreeAssignmentsThrough() {
	c.Fold("[1].a = 2", "1")
	c.Fold("([1]).a = (2)", "1")
	c.Fold("/x/.lastIndex = 0", "0")
	c.FoldSame("({x: 1}).y += 3")
	c.FoldSame("x.y.z = 1")
	c.FoldSame("this.q = 2")
	c.FoldSame("[1][0] = f()")
}
