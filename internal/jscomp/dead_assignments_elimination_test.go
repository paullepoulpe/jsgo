// Port of com/google/javascript/jscomp/DeadAssignmentsEliminationTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/DeadAssignmentsEliminationTest.java)

/*
 * Copyright 2008 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"testing"
)

type DeadAssignmentsEliminationTest struct {
	CompilerTestCaseBase
}

// Java: DeadAssignmentsEliminationTest()
func NewDeadAssignmentsEliminationTest() *DeadAssignmentsEliminationTest {
	c := &DeadAssignmentsEliminationTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, "var extern;")
	return c
}

func TestDeadAssignmentsEliminationTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewDeadAssignmentsEliminationTest() })
}

// JUnitBefore runs the @Before methods: CompilerTestCase.setUp() (the
// superclass one) first, then customSetUp().
func (c *DeadAssignmentsEliminationTest) JUnitBefore() {
	c.CompilerTestCaseBase.JUnitBefore()
	c.CustomSetUp()
}

// Java: DeadAssignmentsEliminationTest.customSetUp()
func (c *DeadAssignmentsEliminationTest) CustomSetUp() {
	c.EnableNormalize()
}

// Java: DeadAssignmentsEliminationTest.getProcessor(Compiler)
func (c *DeadAssignmentsEliminationTest) GetProcessor(compiler *Compiler) CompilerPass {
	return CompilerPassFunc(func(externs *rhino.Node, js *rhino.Node) {
		NodeTraversal_Traverse(compiler, js, NewDeadAssignmentsElimination(compiler))
	})
}

// Java: DeadAssignmentsEliminationTest.testSimple()
func (c *DeadAssignmentsEliminationTest) TestSimple() {
	c.InFunctionWithExpected("var a; a=1", "var a; 1")
	c.InFunctionWithExpected("var a; a=1+1", "var a; 1+1")
	c.InFunctionWithExpected("var a; a=foo();", "var a; foo()")
	c.InFunctionWithExpected("var a; a=foo?.();", "var a; foo?.()")
	c.InFunctionWithExpected("a=1; var a; a=foo();", "1; var a; foo();")
	c.InFunctionWithExpected("a=1; var a; a=foo?.();", "1; var a; foo?.();")
	// This should be: "var a; (function f(){})", but we don't mess with
	// functions with inner functions.
	c.InFunction("var a; a=function f(){}")
}

// Java: DeadAssignmentsEliminationTest.testPropAssignmentNotRemoved()
func (c *DeadAssignmentsEliminationTest) TestPropAssignmentNotRemoved() {
	// We only remove dead assignments when lhs is a name node.
	c.InFunction("var a = {b:1}; a.b=1+1")
	c.InFunction("var a = {b:1}; a.b=foo();")
}

// Java: DeadAssignmentsEliminationTest.testArguments()
func (c *DeadAssignmentsEliminationTest) TestArguments() {
	c.Test("function f(a){ a=1; }", "function f(a){ 1; }")
	c.Test("function f(a){ a=1+1; }", "function f(a){ 1+1; }")
	c.Test("function f(a){ a=foo(); }", "function f(a){ foo(); }")
	c.Test("function f(a){ a=foo?.(); }", "function f(a){ foo?.(); }")
	c.Test("function f(a){ a=1; a=foo(); }", "function f(a){ 1; foo(); }")
	c.Test("function f(a){ a=1; a=foo?.(); }", "function f(a){ 1; foo?.(); }")
}

// Java: DeadAssignmentsEliminationTest.testLoops()
func (c *DeadAssignmentsEliminationTest) TestLoops() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.InFunction("for(var a=0; a<10; a++) {}")
	c.InFunction("var x; for(var a=0; a<10; a++) {x=a}; a(x)")
	c.InFunctionWithExpected("var x; for(var a=0; x=a<10; a++) {}", "var x; for(var a=0; a<10; a++) {}")
	c.InFunctionWithExpected("var x; for(var a=0; a<10; x=a) {}", "var x; for(var a=0; a<10; a) {}")
	c.InFunctionWithExpected("var x; for(var a=0; a<10; x=a,a++) {}", "var x; for(var a=0; a<10; a,a++) {}")
	c.InFunctionWithExpected("var x; for(var a=0; a<10; a++,x=a) {}", "var x; for(var a=0; a<10; a++,a) {}")
	c.InFunctionWithExpected("var x;for(var a=0; a<10; a++) {x=1}", "var x;for(var a=0; a<10; a++) {1}")
	c.InFunctionWithExpected("var x; x=1; do{x=2}while(0); x", "var x; 1; do{x=2}while(0); x")
	c.InFunction("var x; x=1; while(1){x=2}; x")
}

// Java: DeadAssignmentsEliminationTest.testMultiPaths()
func (c *DeadAssignmentsEliminationTest) TestMultiPaths() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.InFunctionWithExpected("var x,y; if(x)y=1;", "var x,y; if(x)1;")
	c.InFunctionWithExpected("var x,y; if(x)y=1; y=2; x(y)", "var x,y; if(x)1; y=2; x(y)")
	c.InFunction("var x; switch(x) { case(1): x=1; break; } x")
	c.InFunctionWithExpected(
		"var x; switch(x) { case(1): x=1; break; }", "var x; switch(x) { case(1): 1; break; }")
}

// Java: DeadAssignmentsEliminationTest.testUsedAsConditions()
func (c *DeadAssignmentsEliminationTest) TestUsedAsConditions() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.InFunctionWithExpected("var x; while(x=1){}", "var x; while(1){}")
	c.InFunctionWithExpected("var x; if(x=1){}", "var x; if(1){}")
	c.InFunctionWithExpected("var x; do{}while(x=1)", "var x; do{}while(1)")
	c.InFunctionWithExpected("var x; if(x=1==4&&1){}", "var x; if(1==4&&1) {}")
	c.InFunctionWithExpected("var x; if(0&&(x=1)){}", "var x; if(0&&1){}")
	c.InFunctionWithExpected("var x; if((x=2)&&(x=1)){}", "var x; if(2&&1){}")
	c.InFunction("var x; x=2; if(0&&(x=1)){}; x")

	c.InFunctionWithExpected("var x,y; if( (x=1)+(y=2) > 3){}", "var x,y; if( 1+2 > 3){}")
}

// Java: DeadAssignmentsEliminationTest.nullishCoalesce()
func (c *DeadAssignmentsEliminationTest) TestNullishCoalesce() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.InFunctionWithExpected("var x; if(x=1==4??1){}", "var x; if(1==4??1) {}")
	c.InFunctionWithExpected("var x; if(0??(x=1)){}", "var x; if(0??1){}")
	c.InFunctionWithExpected("var x; if((x=2)??(x=1)){}", "var x; if(2??1){}")
	c.InFunction("var x; x=2; if(0??(x=1)){}; x")
	c.InFunction("var a, b; if ((a = 1) ?? (b = a)) {b}")
	c.InFunctionWithExpected("var a, b; if ((b = a) ?? (a = 1)) {b}", "var a, b; if ((b = a) ?? (1)) {b}")
	c.InFunctionWithExpected("var a; (a = 1) ?? (a = 2)", "var a; 1 ?? 2")
}

// Java: DeadAssignmentsEliminationTest.testUsedAsConditionsInSwitchStatements()
func (c *DeadAssignmentsEliminationTest) TestUsedAsConditionsInSwitchStatements() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.InFunctionWithExpected("var x; switch(x=1){}", "var x; switch(1){}")
	c.InFunctionWithExpected("var x; switch(x){case(x=1):break;}", "var x; switch(x){case(1):break;}")

	c.InFunction("var x,y; switch(y) { case (x += 1): break; case (x): break;}")

	c.InFunctionWithExpected(
		"var x,y; switch(y) { case (x = 1): break; case (2): break;}",
		"var x,y; switch(y) { case (1): break; case (2): break;}")
	c.InFunctionWithExpected(
		"var x,y; switch(y) { case (x+=1): break; case (x=2): break;}",
		"var x,y; switch(y) { case (x+1): break; case (2): break;}")
}

// Java: DeadAssignmentsEliminationTest.testAssignmentInReturn()
func (c *DeadAssignmentsEliminationTest) TestAssignmentInReturn() {
	c.InFunctionWithExpected("var x; return x = 1;", "var x; return 1")
	c.InFunction("var x; return")
}

// Java: DeadAssignmentsEliminationTest.testAssignmentSamples()
func (c *DeadAssignmentsEliminationTest) TestAssignmentSamples() {
	// We want this to be "var x" in these cases.
	c.InFunction("var x = 2;")
	c.InFunctionWithExpected("var x = 2; x++;", "var x=2; void 0")
	c.InFunctionWithExpected("var x; x=x++;", "var x;x++")
	c.InFunctionWithExpected("var x; x+=1;", "var x;x+1")
}

// Java: DeadAssignmentsEliminationTest.testAssignmentInArgs()
func (c *DeadAssignmentsEliminationTest) TestAssignmentInArgs() {
	c.InFunctionWithExpected("var x; foo(x = 1);", "var x; foo(1);")
	c.InFunctionWithExpected("var x; return foo(x = 1);", "var x; return foo(1);")
}

/** BUG #1358904 */
// Java: DeadAssignmentsEliminationTest.testAssignAndReadInCondition()
func (c *DeadAssignmentsEliminationTest) TestAssignAndReadInCondition() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.InFunction("var a, b; if ((a = 1) && (b = a)) {b}")
	c.InFunctionWithExpected("var a, b; if ((b = a) && (a = 1)) {b}", "var a, b; if ((b = a) && (1)) {b}")
}

// Java: DeadAssignmentsEliminationTest.testParameters()
func (c *DeadAssignmentsEliminationTest) TestParameters() {
	c.InFunctionWithExpected("param1=1; param1=2; param2(param1)", "1; param1=2; param2(param1)")
	c.InFunctionWithExpected("param1=param2()", "param2()")
}

// Java: DeadAssignmentsEliminationTest.testErrorHandling()
func (c *DeadAssignmentsEliminationTest) TestErrorHandling() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.InFunction("var x; try{ x=1 } catch(e){ x=2 }; x")
	c.InFunctionWithExpected("var x; try{ x=1 } catch(e){ x=2 }", "var x;try{ 1 } catch(e) { 2 }")
	c.InFunctionWithExpected("var x; try{ x=1 } finally { x=2 }; x", "var x;try{ 1 } finally{ x=2 }; x")
	c.InFunction("var x; while(1) { try{x=1;break}finally{x} }")
	c.InFunction("var x; try{throw 1} catch(e){x=2} finally{x}")
	c.InFunctionWithExpected(
		"var x; try{x=1;throw 1;x} finally{x=2}; x", "var x; try{1;throw 1;x} finally{x=2}; x")
}

// Java: DeadAssignmentsEliminationTest.testErrorHandling2()
func (c *DeadAssignmentsEliminationTest) TestErrorHandling2() {
	c.InFunction(
		"try {\n" +
			"} catch (e) {\n" +
			"  e = 1;\n" +
			"  let g = e;\n" +
			"  print(g)\n" +
			"}\n")

	c.InFunction(
		"try {\n" +
			"} catch (e) {\n" +
			"    e = 1;\n" +
			"    {\n" +
			"      let g = e;\n" +
			"      print(g)\n" +
			"    }\n" +
			"}\n")
}

// Java: DeadAssignmentsEliminationTest.testDeadVarDeclarations1()
func (c *DeadAssignmentsEliminationTest) TestDeadVarDeclarations1() {
	c.InFunctionWithExpected("var x=1; x=2; x", "var x; 1; x=2; x")
}

// Java: DeadAssignmentsEliminationTest.testDeadVarDeclarations2()
func (c *DeadAssignmentsEliminationTest) TestDeadVarDeclarations2() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.InFunction("var x=1;")
	c.InFunctionWithExpected("var x=1; x=2; x", "var x; 1; x=2; x")
	c.InFunctionWithExpected("var x=1, y=10; x=2; x", "var x; 1; var y; 10; x=2; x")
	c.InFunction("var x=1, y=x; y")
	c.InFunctionWithExpected("var x=1, y=x; x=2; x", "var x = 1; var y; x; x=2; x;")
}

// Java: DeadAssignmentsEliminationTest.testDeadVarDeclarations_forLoop()
func (c *DeadAssignmentsEliminationTest) TestDeadVarDeclarations_forLoop() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.InFunction("for(var x=1;;);")
	c.InFunction("for(var x=1,y=x;;);")
	c.InFunction("for(var x=1;10;);")
}

// Java: DeadAssignmentsEliminationTest.testGlobal()
func (c *DeadAssignmentsEliminationTest) TestGlobal() {
	// Doesn't do any work on global scope yet.
	c.Test("var x; x=1; x=2; x=3;", "var x; x=1; x=2; x=3;")
}

// Java: DeadAssignmentsEliminationTest.testInnerFunctions()
func (c *DeadAssignmentsEliminationTest) TestInnerFunctions() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.InFunctionWithExpected("var x = function() { var x; x=1; }", "var x = function() { var x; 1; }")
}

// Java: DeadAssignmentsEliminationTest.testInnerFunctions2()
func (c *DeadAssignmentsEliminationTest) TestInnerFunctions2() {
	// Give up DCE if there is a inner function.
	c.InFunction("var x = 0; print(x); x = 1; var y = function(){}; y()")
}

// Java: DeadAssignmentsEliminationTest.testSelfReAssignment()
func (c *DeadAssignmentsEliminationTest) TestSelfReAssignment() {
	c.InFunctionWithExpected("var x; x = x;", "var x; x")
}

// Java: DeadAssignmentsEliminationTest.testSelfIncrement()
func (c *DeadAssignmentsEliminationTest) TestSelfIncrement() {
	c.InFunctionWithExpected("var x; x = x + 1;", "var x; x + 1")
}

// Java: DeadAssignmentsEliminationTest.testAssignmentOp()
func (c *DeadAssignmentsEliminationTest) TestAssignmentOp() {
	// We have remove constant expressions that cleans this one up.
	c.InFunctionWithExpected("var x; x += foo()", "var x; x + foo()")
}

// Java: DeadAssignmentsEliminationTest.testAssignmentOpUsedAsLhs()
func (c *DeadAssignmentsEliminationTest) TestAssignmentOpUsedAsLhs() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.InFunctionWithExpected("var x,y; y = x += foo(); print(y)", "var x,y; y = x +  foo(); print(y)")
	c.InFunctionWithExpected("var x,y; y = x += foo?.(); print(y)", "var x,y; y = x +  foo?.(); print(y)")
}

// Java: DeadAssignmentsEliminationTest.testAssignmentOpUsedAsCondition()
func (c *DeadAssignmentsEliminationTest) TestAssignmentOpUsedAsCondition() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.InFunctionWithExpected("var x; if(x += foo()) {}", "var x; if(x +  foo()) {}")
	c.InFunctionWithExpected("var x; if(x += foo?.()) {}", "var x; if(x +  foo?.()) {}")

	c.InFunctionWithExpected("var x; if((x += foo()) > 1) {}", "var x; if((x +  foo()) > 1) {}")

	// Not in a while because this happens every loop.
	c.InFunction("var x; while((x += foo()) > 1) {}")

	c.InFunction("var x; for(;--x;){}")
	c.InFunction("var x; for(;x--;){}")
	c.InFunction("var x; for(;x -= 1;){}")
	c.InFunctionWithExpected("var x; for(;x = 0;){}", "var x; for(;0;){}")

	c.InFunction("var x; for(;;--x){}")
	c.InFunction("var x; for(;;x--){}")
	c.InFunction("var x; for(;;x -= 1){}")
	c.InFunctionWithExpected("var x; for(;;x = 0){}", "var x; for(;;0){}")

	c.InFunctionWithExpected("var x; for(--x;;){}", "var x; void 0; for(;;){}")
	c.InFunctionWithExpected("var x; for(x--;;){}", "var x; void 0; for(;;){}")
	c.InFunctionWithExpected("var x; for(x -= 1;;){}", "var x; for(x - 1;;){}")
	c.InFunctionWithExpected("var x; for(x = 0;;){}", "var x; for(0;;){}")
}

// Java: DeadAssignmentsEliminationTest.testDeadIncrement()
func (c *DeadAssignmentsEliminationTest) TestDeadIncrement() {
	// TODO(user): Optimize this.
	c.InFunctionWithExpected("var x; x ++", "var x; void 0")
	c.InFunctionWithExpected("var x; x --", "var x; void 0")
}

// Java: DeadAssignmentsEliminationTest.testDeadButAlivePartiallyWithinTheExpression()
func (c *DeadAssignmentsEliminationTest) TestDeadButAlivePartiallyWithinTheExpression() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.InFunctionWithExpected("var x; x = 100, print(x), x = 101;", "var x; x = 100, print(x),     101;")
	c.InFunctionWithExpected(
		"var x; x = 100, print(x), print(x), x = 101;",
		"var x; x = 100, print(x), print(x),     101;")
	c.InFunctionWithExpected(
		"var x; x = 100, print(x), x = 0, print(x), x = 101;",
		"var x; x = 100, print(x), x = 0, print(x),     101;")

	// Here, `a=C` is removed as it is dead. `X=a` is removed as it is dead.
	c.InFunctionWithExpected(
		"var a, C, X, S; if ((X = a) && (a = C)) {}; a = S;", //
		"var a, C, X, S; if (a&&C) {}; S;")

	// Here, `a=C` is preserved as it is NOT dead. `X=a` is removed as it is dead.
	c.InFunctionWithExpected(
		"var a, C, X, S; if ((a = C) && (X = a)) {}; a = S;",
		"var a, C, X, S; if ((a = C) && a) {}; S;")
}

// Java: DeadAssignmentsEliminationTest.testMutipleDeadAssignmentsButAlivePartiallyWithinTheExpression()
func (c *DeadAssignmentsEliminationTest) TestMutipleDeadAssignmentsButAlivePartiallyWithinTheExpression() {
	c.InFunctionWithExpected(
		"var x; x = 1, x = 2, x = 3, x = 4, x = 5,\n"+
			"  print(x), x = 0, print(x), x = 101;\n",
		"var x; 1, 2, 3, 4, x = 5, print(x), x = 0, print(x), 101;")
}

// Java: DeadAssignmentsEliminationTest.testDeadPartiallyWithinTheExpression()
func (c *DeadAssignmentsEliminationTest) TestDeadPartiallyWithinTheExpression() {
	// Sadly, this is not covered. We don't suspect this would happen too
	// often.
	c.InFunction("var x; x = 100, x = 101; print(x);")
}

// Java: DeadAssignmentsEliminationTest.testAssignmentChain()
func (c *DeadAssignmentsEliminationTest) TestAssignmentChain() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.InFunctionWithExpected("var a,b,c,d,e; a = b = c = d = e = 1", "var a,b,c,d,e; 1")
	c.InFunctionWithExpected(
		"var a,b,c,d,e; a = b = c = d = e = 1; print(c)",
		"var a,b,c,d,e;         c = 1        ; print(c)")
	c.InFunctionWithExpected(
		"var a,b,c,d,e; a = b = c = d = e = 1; print(a + e)",
		"var a,b,c,d,e; a =             e = 1; print(a + e)")
	c.InFunctionWithExpected(
		"var a,b,c,d,e; a = b = c = d = e = 1; print(b + d)",
		"var a,b,c,d,e;     b =     d     = 1; print(b + d)")
	c.InFunctionWithExpected(
		"var a,b,c,d,e; a = b = c = d = e = 1; print(a + b + d + e)",
		"var a,b,c,d,e; a = b =     d = e = 1; print(a + b + d + e)")
	c.InFunction("var a,b,c,d,e; a = b = c = d = e = 1; print(a+b+c+d+e)")
}

// Java: DeadAssignmentsEliminationTest.testAssignmentOpChain()
func (c *DeadAssignmentsEliminationTest) TestAssignmentOpChain() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.InFunctionWithExpected("var a,b,c,d,e; a = b = c += d = e = 1", "var a,b,c,d,e;         c + 1")
	c.InFunctionWithExpected(
		"var a,b,c,d,e; a = b = c += d = e = 1;  print(e)",
		"var a,b,c,d,e;         c +     (e = 1); print(e)")
	c.InFunctionWithExpected(
		"var a,b,c,d,e; a = b = c += d = e = 1;  print(d)",
		"var a,b,c,d,e;         c + (d = 1)  ;   print(d)")
	c.InFunctionWithExpected(
		"var a,b,c,d,e; a = b = c += d = e = 1;  print(a)",
		"var a,b,c,d,e; a =     c +          1;  print(a)")
}

// Java: DeadAssignmentsEliminationTest.testIncDecInSubExpressions()
func (c *DeadAssignmentsEliminationTest) TestIncDecInSubExpressions() {
	c.InFunction("var a; a = 1, a++; a")
	c.InFunction("var a; a = 1, ++a; a")
	c.InFunction("var a; a = 1, a--; a")
	c.InFunction("var a; a = 1, --a; a")

	c.InFunction("var a; a = 1, a++, print(a)")
	c.InFunction("var a; a = 1, ++a, print(a)")
	c.InFunction("var a; a = 1, a--, print(a)")
	c.InFunction("var a; a = 1, --a, print(a)")

	c.InFunction("var a; a = 1, print(a++)")
	c.InFunction("var a; a = 1, print(++a)")

	c.InFunction("var a; a = 1, print(a++)")
	c.InFunction("var a; a = 1, print(++a)")

	c.InFunction("var a; a = 1, print(a--)")
	c.InFunction("var a; a = 1, print(--a)")
}

// Java: DeadAssignmentsEliminationTest.testNestedReassignments()
func (c *DeadAssignmentsEliminationTest) TestNestedReassignments() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.InFunctionWithExpected("var a; a = (a = 1)", "var a; 1")
	c.InFunctionWithExpected("var a; a = (a *= 2)", "var a; a*2")

	// Note a = (a++) is not same as a++. Only if 'a' is dead.
	c.InFunctionWithExpected("var a; a = (a++)", "var a; a++") // Preferred: "var a"
	c.InFunctionWithExpected("var a; a = (++a)", "var a; ++a") // Preferred: "var a"

	c.InFunctionWithExpected("var a; a = (b = (a = 1))", "var a; b = 1")
	c.InFunctionWithExpected("var a; a = (b = (a *= 2))", "var a; b = a * 2")
	c.InFunctionWithExpected("var a; a = (b = (a++))", "var a; b=a++")
	c.InFunctionWithExpected("var a; a = (b = (++a))", "var a; b=++a")

	// Include b as local.
	c.InFunctionWithExpected("var a,b; a = (b = (a = 1))", "var a,b; 1")
	c.InFunctionWithExpected("var a,b; a = (b = (a *= 2))", "var a,b; a * 2")
	c.InFunctionWithExpected("var a,b; a = (b = (a++))", "var a,b; a++") // Preferred: "var a,b"
	c.InFunctionWithExpected("var a,b; a = (b = (++a))", "var a,b; ++a") // Preferred: "var a,b"

	c.InFunctionWithExpected("var a; a += (a++)", "var a; a + a++")
	c.InFunctionWithExpected("var a; a += (++a)", "var a; a+ (++a)")

	// Include b as local.
	c.InFunctionWithExpected("var a,b; a += (b = (a = 1))", "var a,b; a + 1")
	c.InFunctionWithExpected("var a,b; a += (b = (a *= 2))", "var a,b; a + (a * 2)")
	c.InFunctionWithExpected("var a,b; a += (b = (a++))", "var a,b; a + a++")
	c.InFunctionWithExpected("var a,b; a += (b = (++a))", "var a,b; a+(++a)")
}

// Java: DeadAssignmentsEliminationTest.testIncrementalReassignmentInForLoops()
func (c *DeadAssignmentsEliminationTest) TestIncrementalReassignmentInForLoops() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.InFunction("for(;x+=1;x+=1) {}")
	c.InFunction("for(;x;x+=1){}")
	c.InFunction("for(;x+=1;){foo(x)}")
	c.InFunction("for(;1;x+=1){foo(x)}")
}

// Java: DeadAssignmentsEliminationTest.testIdentityAssignments()
func (c *DeadAssignmentsEliminationTest) TestIdentityAssignments() {
	c.InFunctionWithExpected("var x; x=x", "var x; x")
	c.InFunction("var x; x.y=x.y")
}

// Java: DeadAssignmentsEliminationTest.testBug8730257()
func (c *DeadAssignmentsEliminationTest) TestBug8730257() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.InFunction(
		"try {\n" +
			"   var sortIndices = {};\n" +
			"   sortIndices = bar();\n" +
			"   for (var i = 0; i < 100; i++) {\n" +
			"     var sortIndex = sortIndices[i];\n" +
			"     bar(sortIndex);\n" +
			"   }\n" +
			" } finally {\n" +
			"   bar();\n" +
			" }\n")
}

// Java: DeadAssignmentsEliminationTest.testAssignToExtern()
func (c *DeadAssignmentsEliminationTest) TestAssignToExtern() {
	c.InFunction("extern = true;")
}

// Java: DeadAssignmentsEliminationTest.testIssue297a()
func (c *DeadAssignmentsEliminationTest) TestIssue297a() {
	c.TestSame(
		"function f(p) {\n" +
			" var x;\n" +
			" return ((x=p.id) && (x=parseInt(x.substr(1))) && x>0);\n" +
			"}; f('');\n")
}

// Java: DeadAssignmentsEliminationTest.testIssue297b()
func (c *DeadAssignmentsEliminationTest) TestIssue297b() {
	c.Test(
		"function f() {\n"+
			" var x;\n"+
			" return (x='') && (x = x.substr(1));\n"+
			"};\n",
		"function f() {\n"+
			" var x;\n"+
			" return (x='') && (x.substr(1));\n"+
			"};\n")
}

// Java: DeadAssignmentsEliminationTest.testIssue297c()
func (c *DeadAssignmentsEliminationTest) TestIssue297c() {
	c.Test(
		"function f() {\n"+
			" var x;\n"+
			" return (x=1) && (x = f(x));\n"+
			"};\n",
		"function f() {\n"+
			" var x;\n"+
			" return (x=1) && f(x);\n"+
			"};\n")
}

// Java: DeadAssignmentsEliminationTest.testIssue297d()
func (c *DeadAssignmentsEliminationTest) TestIssue297d() {
	c.Test(
		"function f(a) {\n"+
			" return (a=1) && (a = f(a));\n"+
			"};\n",
		"function f(a) {\n"+
			" return (a=1) && (f(a));\n"+
			"};\n")
}

// Java: DeadAssignmentsEliminationTest.testIssue297e()
func (c *DeadAssignmentsEliminationTest) TestIssue297e() {
	c.Test(
		"function f(a) {\n"+
			" return (a=1) - (a = g(a));\n"+
			"};\n",
		"function f(a) {\n"+
			" return (a=1) - (g(a));\n"+
			"};\n")
}

// Java: DeadAssignmentsEliminationTest.testIssue297f()
func (c *DeadAssignmentsEliminationTest) TestIssue297f() {
	c.Test(
		"function f(a) {\n"+
			" h((a=1) - (a = g(a)));\n"+
			"};\n",
		"function f(a) {\n"+
			" h((a=1) - (g(a)));\n"+
			"};\n")
}

// Java: DeadAssignmentsEliminationTest.testIssue297g()
func (c *DeadAssignmentsEliminationTest) TestIssue297g() {
	c.Test(
		"function f(a) {\n"+
			" var b = h((b=1) - (b = g(b)));\n"+
			" return b;\n"+
			"};\n",
		// The last assignment in the initializer should be eliminated
		"function f(a) {\n"+
			" var b = h((b=1) - (b = g(b)));\n"+
			" return b;\n"+
			"};\n")
}

// Java: DeadAssignmentsEliminationTest.testIssue297h()
func (c *DeadAssignmentsEliminationTest) TestIssue297h() {
	c.Test(
		"function f(a) {\n"+
			" var b = b=1;\n"+
			" return b;\n"+
			"};\n",
		// The assignment in the initializer should be eliminated
		"function f(a) {\n"+
			" var b = b = 1;\n"+
			" return b;\n"+
			"};\n")
}

// Java: DeadAssignmentsEliminationTest.testInExpression0()
func (c *DeadAssignmentsEliminationTest) TestInExpression0() {
	c.InFunctionWithExpected("var a; return a=(a=(a=a));", "var a; return a;")
}

// Java: DeadAssignmentsEliminationTest.testInExpression1()
func (c *DeadAssignmentsEliminationTest) TestInExpression1() {
	c.InFunctionWithExpected("var a; return a=(a=(a=3));", "var a; return 3;")
	c.InFunctionWithExpected("var a; return a=(a=(a=a));", "var a; return a;")
	c.InFunctionWithExpected("var a; return a=(a=(a=a+1)+1);", "var a; return a+1+1;")
	c.InFunctionWithExpected("var a; return a=(a=(a=f(a)+1)+1);", "var a; return f(a)+1+1;")
	c.InFunctionWithExpected("var a; return a=f(a=f(a=f(a)));", "var a; return f(f(f(a)));")
}

// Java: DeadAssignmentsEliminationTest.testInExpression2()
func (c *DeadAssignmentsEliminationTest) TestInExpression2() {
	c.InFunctionWithExpected("var a; (a = 1) || (a = 2)", "var a; 1 || 2")

	c.InFunction("var a; (a = 1) || (a = 2); return a")

	c.InFunctionWithExpected("var a; a = 1; a ? a = 2 : a;", "var a; a = 1; a ?     2 : a;")

	c.InFunction("var a; a = 1; a ? a = 2 : a; return a")

	c.InFunctionWithExpected("var a; a = 1; a ? a : a = 2;", "var a; a = 1; a ? a : 2;")

	c.InFunction("var a; a = 1; a ? a : a =2; return a")

	c.InFunctionWithExpected("var a; (a = 1) ? a = 2 : a = 3;", "var a;      1  ?     2 :     3;")

	// This can be improved.  "a = 1" is dead but "a" is read in the following
	// expression.
	c.InFunction("var a; (a = 1) ? a = 2 : a = 3; return a")
}

// Java: DeadAssignmentsEliminationTest.testIssue384a()
func (c *DeadAssignmentsEliminationTest) TestIssue384a() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.InFunction(
		"var a, b;\n" +
			"if (f(b = true) || f(b = false))\n" +
			"  a = b;\n" +
			"else\n" +
			"  a = null;\n" +
			"return a;\n")
}

// Java: DeadAssignmentsEliminationTest.testIssue384b()
func (c *DeadAssignmentsEliminationTest) TestIssue384b() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.InFunction(
		"var a, b;\n" +
			"(f(b = true) || f(b = false)) ? (a = b) : (a = null);\n" +
			"return a;\n")
}

// Java: DeadAssignmentsEliminationTest.testIssue384c()
func (c *DeadAssignmentsEliminationTest) TestIssue384c() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.InFunction(
		"var a, b;\n" +
			"(a ? f(b = true) : f(b = false)) && (a = b);\n" +
			"return a;\n")
}

// Java: DeadAssignmentsEliminationTest.testIssue384d()
func (c *DeadAssignmentsEliminationTest) TestIssue384d() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.InFunction(
		"var a, b;\n" +
			"(f(b = true) || f(b = false)) && (a = b);\n" +
			"return a;\n")
}

// Java: DeadAssignmentsEliminationTest.testForIn()
func (c *DeadAssignmentsEliminationTest) TestForIn() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.InFunction("var x = {}; for (var y in x) { y() }")
	c.InFunctionWithExpected(
		"var x, y, z; x = {}; z = {}; for (y in x = z) { y() }",
		"var x, y, z;   ({}); z = {}; for (y in z)     { y() }")
	c.InFunctionWithExpected(
		"var x, y, z; x = {}; z = {}; for (y[z=1] in z) { y() }",
		"var x, y, z;   ({}); z = {}; for (y[z=1] in z) { y() }")

	// "x in z" doesn't overwrite x if z is empty.
	// TODO(user): If you look outside of just liveness, x = {} is dead.
	// That probably requires value numbering or SSA to detect that case.
	c.InFunction("var x, y, z; x = {}; z = {}; for (x in z) { x() }")
}

// Java: DeadAssignmentsEliminationTest.testArrowFunction()
func (c *DeadAssignmentsEliminationTest) TestArrowFunction() {
	c.Test("() => {var x; x = 1}", "() => {var x; 1}")

	c.Test("(a) => {a = foo()}", "(a) => {foo()}")
}

// Java: DeadAssignmentsEliminationTest.testClassMethods()
func (c *DeadAssignmentsEliminationTest) TestClassMethods() {
	c.Test(
		"class C{\n"+
			"  func() {\n"+
			"    var x;\n"+
			"    x = 1;\n"+
			"  }\n"+
			"}\n",
		"class C{\n"+
			"  func() {\n"+
			"    var x;\n"+
			"    1;\n"+
			"  }\n"+
			"}\n")

	c.Test(
		"class C{\n"+
			"  constructor(x, y) {\n"+
			"    this.x = x;\n"+
			"    this.y = y;\n"+
			"  }\n"+
			"  func() {\n"+
			"    var z;\n"+
			"    z = 1;\n"+
			"    this.x = 3\n"+
			"  }\n"+
			"}\n",
		"class C{\n"+
			"  constructor(x, y) {\n"+
			"    this.x = x;\n"+
			"    this.y = y;\n"+
			"  }\n"+
			"  func() {\n"+
			"    var z;\n"+
			"    1;\n"+
			"    this.x = 3\n"+
			"  }\n"+
			"}\n")
}

// Java: DeadAssignmentsEliminationTest.testClassStaticBlocks()
func (c *DeadAssignmentsEliminationTest) TestClassStaticBlocks() {
	// TODO(b/240443227): Improve ClassStaticBlock optimization, dead code is not removed in
	// expression.
	c.TestSame(
		"class C{\n" +
			"  static{\n" +
			"    var x;\n" +
			"    x = 1;\n" +
			"  }\n" +
			"}\n")

	c.TestSame(
		" var x = 0;\n" +
			" print(x);\n" +
			"  x = 1;\n" +
			"  class C {\n" +
			"   static {\n" +
			"    print(x);\n" +
			"    }\n" +
			"  }\n")
}

// Java: DeadAssignmentsEliminationTest.testGenerators()
func (c *DeadAssignmentsEliminationTest) TestGenerators() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.Test(
		"function* f() {\n"+
			"  var x, y;\n"+
			"  x = 1; y = 2;\n"+
			"  yield y;\n"+
			"}\n",
		"function* f() {\n"+
			"  var x, y;\n"+
			"  1; y = 2;\n"+
			"  yield y;\n"+
			"}\n")
}

// Java: DeadAssignmentsEliminationTest.testForOf()
func (c *DeadAssignmentsEliminationTest) TestForOf() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.InFunction("var x = {}; for (var y of x) { y() }")

	c.InFunctionWithExpected(
		"var x, y, z; x = {}; z = {}; for (y of x = z) {}",
		"var x, y, z;   ({}); z = {}; for (y of z)     {}")
}

// Java: DeadAssignmentsEliminationTest.testForAwaitOf()
func (c *DeadAssignmentsEliminationTest) TestForAwaitOf() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.InAsyncFunction("var x = {}; for await (var y of x) { y() }")

	c.InAsyncFunctionWithExpected(
		"var x, y, z; x = {}; z = {}; for await (y of x = z) {}",
		"var x, y, z;   ({}); z = {}; for await (y of z)     {}")
}

// Java: DeadAssignmentsEliminationTest.testTemplateStrings()
func (c *DeadAssignmentsEliminationTest) TestTemplateStrings() {
	c.InFunction("var name; name = 'Foo'; `Hello ${name}`")

	c.InFunctionWithExpected(
		"var name; name = 'Foo'; name = 'Bar'; `Hello ${name}`",
		"var name; 'Foo'; name = 'Bar'; `Hello ${name}`")
}

// Java: DeadAssignmentsEliminationTest.testDestructuring()
func (c *DeadAssignmentsEliminationTest) TestDestructuring() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.InFunction("var a, b, c; [a, b, c] = [1, 2, 3];")

	c.InFunction("var a, b, c; [a, b, c] = [1, 2, 3]; return a + c;")

	c.InFunctionWithExpected(
		"var a, b; a = 1; b = 2; [a, b] = [3, 4]; return a + b;",
		"var a, b; 1; 2; [a, b] = [3, 4]; return a + b;")

	c.InFunction("var x; x = {}; [x.a] = [3];")
}

// Java: DeadAssignmentsEliminationTest.testDestructuringDeclarationRvalue()
func (c *DeadAssignmentsEliminationTest) TestDestructuringDeclarationRvalue() {
	// Test array destructuring
	c.InFunction(
		"let arr = []\n" +
			"if (CONDITION) {\n" +
			"  arr = [3];\n" +
			"}\n" +
			"let [foo] = arr;\n" +
			"use(foo);\n")

	// Test object destructuring
	c.InFunction(
		"let obj = {}\n" +
			"if (CONDITION) {\n" +
			"  obj = {foo: 3};\n" +
			"}\n" +
			"let {foo} = obj;\n" +
			"use(foo);\n")
}

// Java: DeadAssignmentsEliminationTest.testDestructuringAssignmentRValue()
func (c *DeadAssignmentsEliminationTest) TestDestructuringAssignmentRValue() {
	// Test array destructuring
	c.InFunction(
		"let arr = []\n" +
			"if (CONDITION) {\n" +
			"  arr = [3];\n" +
			"}\n" +
			"let foo;\n" +
			"[foo] = arr;\n" +
			"use(foo);\n")

	// Test object destructuring
	c.InFunction(
		"let obj = {}\n" +
			"if (CONDITION) {\n" +
			"  obj = {foo: 3};\n" +
			"}\n" +
			"let foo;\n" +
			"({foo} = obj);\n" +
			"use(foo);\n")
}

// Java: DeadAssignmentsEliminationTest.testForOfWithDestructuring()
func (c *DeadAssignmentsEliminationTest) TestForOfWithDestructuring() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.InFunction(
		"let x;\n" +
			"x = [];\n" +
			"var y = 5; // Don't eliminate because if arr is empty, y will remain 5.\n" +
			"for ([y = x] of arr) { y; }\n" +
			"y;\n")

	c.InFunction(
		"let x;\n" +
			"x = [];\n" +
			"for (let [y = x] of arr) { y; }\n")

	c.InFunction("for (let [key, value] of arr) {}")
	c.InFunction("for (let [key, value] of arr) { key; value; }")
	c.InFunctionWithExpected(
		"var a; a = 3; for (let [a] of arr) { a; }", "var a; 3; for (let [a] of arr) { a; }")
}

// Java: DeadAssignmentsEliminationTest.testReferenceInDestructuringPatternDefaultValue()
func (c *DeadAssignmentsEliminationTest) TestReferenceInDestructuringPatternDefaultValue() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.InFunction(
		"let bar = [];\n" +
			"const {foo = bar} = obj;\n" +
			"foo;\n")

	c.InFunction(
		"let bar;\n" +
			"bar = [];\n" +
			"const {foo = bar} = obj;\n" +
			"foo;\n")

	c.InFunction("let bar; bar = 3; const [foo = bar] = arr; foo;")
	c.InFunction("let foo, bar; bar = 3; [foo = bar] = arr; foo;")
}

// Java: DeadAssignmentsEliminationTest.testReferenceInDestructuringPatternComputedProperty()
func (c *DeadAssignmentsEliminationTest) TestReferenceInDestructuringPatternComputedProperty() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.InFunction("let str; str = 'bar'; const {[str + 'baz']: foo} = obj; foo;")

	c.InFunction(
		"let obj = {};\n" +
			"let str, foo;\n" +
			"str = 'bar';\n" +
			"({[str + 'baz']: foo} = obj);\n" +
			"foo;\n")
}

// Java: DeadAssignmentsEliminationTest.testDefaultParameter()
func (c *DeadAssignmentsEliminationTest) TestDefaultParameter() {
	c.Test(
		"function f(x, y = 12) {\n"+
			"  var z;\n"+
			"  z = y;\n"+
			"}\n",
		"function f(x, y = 12) {\n"+
			"  var z;\n"+
			"  y;\n"+
			"}\n")
}

// Java: DeadAssignmentsEliminationTest.testObjectLiterals()
func (c *DeadAssignmentsEliminationTest) TestObjectLiterals() {
	c.Test(
		"var obj = {\n"+
			"  f() {\n"+
			"  var x;\n"+
			"  x = 2;\n"+
			"  }\n"+
			"}\n",
		"var obj = {\n"+
			"  f() {\n"+
			"  var x;\n"+
			"  2;\n"+
			"  }\n"+
			"}\n")
}

// Java: DeadAssignmentsEliminationTest.testObjectLiteralsComputedProperties()
func (c *DeadAssignmentsEliminationTest) TestObjectLiteralsComputedProperties() {
	c.InFunction("let a; a = 2; let obj = {[a]: 3}; obj")
}

// Java: DeadAssignmentsEliminationTest.testSpread_consideredRead()
func (c *DeadAssignmentsEliminationTest) TestSpread_consideredRead() {
	c.InFunction(
		"var a;\n" +
			"a = [];\n" +
			"[...a];\n")

	c.InFunction(
		"var a;\n" +
			"a = {};\n" +
			"({...a});\n")
}

// Java: DeadAssignmentsEliminationTest.testRest_notConsideredWrite()
func (c *DeadAssignmentsEliminationTest) TestRest_notConsideredWrite() {
	// TODO(b/126441776): The initial writes are dead. The pass should rewrite to the commented
	// code.

	c.InFunction(
		"var a = 9;\n" +
			"[...a] = itr;\n" +
			"return a;\n",
	// ,
	// """
	// var a;
	// [...a] = itr;
	// return a;
	// """
	)

	c.InFunction(
		"var a = 9;\n" +
			"({...a} = obj);\n" +
			"return a;\n",
	// ,
	// """
	// var a;
	// ({...a} = obj);
	// return a;
	// """
	)
}

// Java: DeadAssignmentsEliminationTest.testDestructuring_notConsideredWrite()
func (c *DeadAssignmentsEliminationTest) TestDestructuring_notConsideredWrite() {
	// TODO(b/126441776): The initial writes are dead. The pass should rewrite to the commented
	// code.

	c.InFunction(
		"var a = 9;\n" +
			"[a] = itr;\n" +
			"return a;\n",
	// ,
	// """
	// var a;
	// [a] = itr;
	// return a;
	// """
	)

	c.InFunction(
		"var a = 9;\n" +
			"({a} = obj);\n" +
			"return a;\n",
	// ,
	// """
	// var a;
	// ({a} = obj);
	// return a;
	// """
	)
}

// Java: DeadAssignmentsEliminationTest.testRest_isNotRemovable()
func (c *DeadAssignmentsEliminationTest) TestRest_isNotRemovable() {
	// TODO(b/126441776): Elimination is possible here under getter/setter assumptions. Determine if
	// this is the correct behaviour.

	c.InFunction(
		"var a;\n" +
			"[...a] = itr;\n")

	c.InFunction(
		"var a;\n" +
			"({...a} = obj);\n")
}

// Java: DeadAssignmentsEliminationTest.testDestructuring_isNotRemovable()
func (c *DeadAssignmentsEliminationTest) TestDestructuring_isNotRemovable() {
	// TODO(b/126441776): Elimination is possible here under getter/setter assumptions. Determine if
	// this is the correct behaviour.

	c.InFunction(
		"var a;\n" +
			"[a] = itr;\n")

	c.InFunction(
		"var a;\n" +
			"({a} = obj);\n")
}

// Java: DeadAssignmentsEliminationTest.testLet()
func (c *DeadAssignmentsEliminationTest) TestLet() {
	c.InFunctionWithExpected("let a; a = 2;", "let a; 2;")

	c.InFunctionWithExpected(
		"let a; let b; a = foo(); b = 2; return b;", "let a; let b; foo(); b = 2; return b;")

	c.InFunctionWithExpected(
		"let a; let b; a = foo?.(); b = 2; return b;", //
		"let a; let b; foo?.(); b = 2; return b;")
}

// Java: DeadAssignmentsEliminationTest.testConst1()
func (c *DeadAssignmentsEliminationTest) TestConst1() {
	c.InFunction("const a = 1;")
}

// Java: DeadAssignmentsEliminationTest.testConst2()
func (c *DeadAssignmentsEliminationTest) TestConst2() {
	c.Test(
		"async function f(d) { if (d) { d = 5; } const a = 1; const b = 2; const [x, y] = b; }",
		"async function f(d) { if (d) {     5; } const a = 1; const b = 2; const [x, y] = b; }")
}

// Java: DeadAssignmentsEliminationTest.testBlockScoping()
func (c *DeadAssignmentsEliminationTest) TestBlockScoping() {
	c.InFunctionWithExpected(
		"let x;\n"+
			"{\n"+
			"  let x;\n"+
			"  x = 1;\n"+
			"}\n"+
			"x = 2;\n"+
			"return x;\n",
		"let x;\n"+
			"{\n"+
			"  let x$jscomp$1;\n"+
			"  1;\n"+
			"}\n"+
			"x = 2;\n"+
			"return x;\n")

	c.InFunctionWithExpected(
		"let x;\n"+
			"x = 2\n"+
			"{\n"+
			"  let x;\n"+
			"  x = 1;\n"+
			"}\n"+
			"print(x);\n",
		"let x;\n"+
			"x = 2;\n"+
			"{\n"+
			"  let x$jscomp$1;\n"+
			"  1;\n"+
			"}\n"+
			"print(x);\n")
}

// Java: DeadAssignmentsEliminationTest.testComputedClassField()
func (c *DeadAssignmentsEliminationTest) TestComputedClassField() {
	c.InFunctionWithExpected(
		"let x;\n"+
			"class C {\n"+
			"  static [x = 'field1'] = (x = 5);\n"+
			"  [x = 'field2'] = 7;\n"+
			"}\n",
		"let x;\n"+
			"class C {\n"+
			"  static ['field1'] = 5;\n"+
			"  ['field2'] = 7;\n"+
			"}\n")

	c.InFunction(
		"let x;\n" +
			"class C {\n" +
			"  static [x = 'field1'] = (x = 5);\n" +
			"  [x = 'field2'] = 7;\n" +
			"}\n" +
			"use(x);\n")

	c.InFunctionWithExpected(
		"let x;\n"+
			"class C {\n"+
			"  static field1 = x;\n"+
			"  [x = 'field2'] = 7;\n"+
			"}\n"+
			"use(C.field1);\n",
		"let x;\n"+
			"class C {\n"+
			"  static field1 = x;\n"+
			"// TODO(b/189993301): don't remove 'x = field2' because it's read by 'field1 = x'\n"+
			"  ['field2'] = 7;\n"+
			"}\n"+
			"use(C.field1);\n")
}

// Java: DeadAssignmentsEliminationTest.testComputedClassMethod()
func (c *DeadAssignmentsEliminationTest) TestComputedClassMethod() {
	c.InFunction(
		"let x;\n" +
			"class C {\n" +
			"// NOTE: it would be correct to eliminate the following two assignments\n" +
			"  static [x = 'field1']() {}\n" +
			"  [x = 'field2']() {}\n" +
			"}\n")

	c.InFunction(
		"let x;\n" +
			"class C {\n" +
			"  static [x = 'field1']() {};\n" +
			"  [x = 'field2']() {}\n" +
			"}\n" +
			"use(x);\n")

	c.InFunction(
		"let x;\n" +
			"class C {\n" +
			"  static field1 = x;\n" +
			"  [x = 'field2']() {}\n" +
			"}\n" +
			"use(C.field1);\n")
}

// Java: DeadAssignmentsEliminationTest.inFunction(String)
func (c *DeadAssignmentsEliminationTest) InFunction(src string) {
	c.InFunctionWithExpected(src, src)
}

// Java: DeadAssignmentsEliminationTest.inFunction(String, String)
func (c *DeadAssignmentsEliminationTest) InFunctionWithExpected(src string, expected string) {
	c.Test(
		"function FUNC(param1, param2){"+src+"}",
		"function FUNC(param1, param2){"+expected+"}")
}

// Java: DeadAssignmentsEliminationTest.inAsyncFunction(String)
func (c *DeadAssignmentsEliminationTest) InAsyncFunction(src string) {
	c.InAsyncFunctionWithExpected(src, src)
}

// Java: DeadAssignmentsEliminationTest.inAsyncFunction(String, String)
func (c *DeadAssignmentsEliminationTest) InAsyncFunctionWithExpected(src string, expected string) {
	c.Test(
		"async function FUNC(param1, param2){"+src+"}",
		"async function FUNC(param1, param2){"+expected+"}")
}
