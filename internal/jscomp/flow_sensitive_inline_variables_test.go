// Port of com/google/javascript/jscomp/FlowSensitiveInlineVariablesTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/FlowSensitiveInlineVariablesTest.java)

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
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// FlowSensitiveInlineVariablesTest: unit tests for
// FlowSensitiveInlineVariables.
type FlowSensitiveInlineVariablesTest struct {
	CompilerTestCaseBase
}

// Java: FlowSensitiveInlineVariablesTest.EXTERN_FUNCTIONS
const FlowSensitiveInlineVariablesTest_EXTERN_FUNCTIONS = `var print;
var alert;
/** @nosideeffects */ function noSFX() {}
                      function hasSFX() {}
`

// Java: FlowSensitiveInlineVariablesTest()
func NewFlowSensitiveInlineVariablesTest() *FlowSensitiveInlineVariablesTest {
	c := &FlowSensitiveInlineVariablesTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestFlowSensitiveInlineVariablesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewFlowSensitiveInlineVariablesTest() })
}

// Java: FlowSensitiveInlineVariablesTest.setUp()
func (c *FlowSensitiveInlineVariablesTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableNormalize()
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
}

// Java: FlowSensitiveInlineVariablesTest.getNumRepetitions()
func (c *FlowSensitiveInlineVariablesTest) GetNumRepetitions() int {
	// Test repeatedly inline.
	return 3
}

// Java: FlowSensitiveInlineVariablesTest.getProcessor(Compiler)
func (c *FlowSensitiveInlineVariablesTest) GetProcessor(compiler *Compiler) CompilerPass {
	return CompilerPassFunc(func(externs *rhino.Node, root *rhino.Node) {
		NewPureFunctionIdentifier_Driver(compiler).Process(externs, root)
		NewFlowSensitiveInlineVariables(compiler).Process(externs, root)
	})
}

// Java: FlowSensitiveInlineVariablesTest.testSimpleAssign()
func (c *FlowSensitiveInlineVariablesTest) TestSimpleAssign() {
	c.Inline("var x; x = 1; print(x)", "var x; print(1)")
	c.Inline("var x; x = 1; x", "var x; 1")
	c.Inline("var x; x = 1; var a = x", "var x; var a = 1")
	c.Inline("var x; x = 1; x = x + 1", "var x; x = 1 + 1")
}

// Java: FlowSensitiveInlineVariablesTest.testSimpleVar()
func (c *FlowSensitiveInlineVariablesTest) TestSimpleVar() {
	c.Inline("var x = 1; print(x)", "var x; print(1)")
	c.Inline("var x = 1; x", "var x; 1")
	c.Inline("var x = 1; var a = x", "var x; var a = 1")
	c.Inline("var x = 1; x = x + 1", "var x; x = 1 + 1")
}

// Java: FlowSensitiveInlineVariablesTest.testSimpleLet()
func (c *FlowSensitiveInlineVariablesTest) TestSimpleLet() {
	c.Inline("let x = 1; print(x)", "let x; print(1)")
	c.Inline("let x = 1; x", "let x; 1")
	c.Inline("let x = 1; let a = x", "let x; let a = 1")
	c.Inline("let x = 1; x = x + 1", "let x; x = 1 + 1")
}

// Java: FlowSensitiveInlineVariablesTest.testSimpleConst()
func (c *FlowSensitiveInlineVariablesTest) TestSimpleConst() {
	c.Inline("const x = 1; print(x)", "const x = undefined; print(1)")
	c.Inline("const x = 1; x", "const x = undefined; 1")
	c.Inline("const x = 1; const a = x", "const x = undefined; const a = 1")
}

// Java: FlowSensitiveInlineVariablesTest.testSimpleForIn()
func (c *FlowSensitiveInlineVariablesTest) TestSimpleForIn() {
	c.Inline(
		"var a,b,x = a in b; x", //
		"var a,b,x;  a in b;  ")
	c.NoInline("var a, b; var x = a in b; print(1); x")
	c.NoInline("var a,b,x = a in b; delete a[b]; x")
}

// Java: FlowSensitiveInlineVariablesTest.testExported()
func (c *FlowSensitiveInlineVariablesTest) TestExported() {
	c.NoInline("var _x = 1; print(_x)")
}

// Java: FlowSensitiveInlineVariablesTest.testDoNotInlineIncrement()
func (c *FlowSensitiveInlineVariablesTest) TestDoNotInlineIncrement() {
	c.NoInline("var x = 1; x++;")
	c.NoInline("var x = 1; x--;")
}

// Java: FlowSensitiveInlineVariablesTest.testMultiUse()
func (c *FlowSensitiveInlineVariablesTest) TestMultiUse() {
	c.NoInline("var x; x = 1; print(x); print (x);")
}

// Java: FlowSensitiveInlineVariablesTest.testMultiUseInSameCfgNode()
func (c *FlowSensitiveInlineVariablesTest) TestMultiUseInSameCfgNode() {
	c.NoInline("var x; x = 1; print(x) || print (x);")
}

// Java: FlowSensitiveInlineVariablesTest.testMultiUseInTwoDifferentPath()
func (c *FlowSensitiveInlineVariablesTest) TestMultiUseInTwoDifferentPath() {
	c.NoInline("var x = 1; if (print) { print(x) } else { alert(x) }")
}

// Java: FlowSensitiveInlineVariablesTest.testAssignmentBeforeDefinition()
func (c *FlowSensitiveInlineVariablesTest) TestAssignmentBeforeDefinition() {
	c.Inline(
		"x = 1; var x = 0; print(x)", //
		"x = 1; var x    ; print(0)")
}

// Java: FlowSensitiveInlineVariablesTest.testVarInConditionPath()
func (c *FlowSensitiveInlineVariablesTest) TestVarInConditionPath() {
	c.NoInline("if (foo) { var x = 0 } print(x)")
}

// Java: FlowSensitiveInlineVariablesTest.testMultiDefinitionsBeforeUse()
func (c *FlowSensitiveInlineVariablesTest) TestMultiDefinitionsBeforeUse() {
	c.Inline("var x = 0; x = 1; print(x)", "var x = 0; print(1)")
}

// Java: FlowSensitiveInlineVariablesTest.testMultiDefinitionsInSameCfgNode()
func (c *FlowSensitiveInlineVariablesTest) TestMultiDefinitionsInSameCfgNode() {
	c.NoInline("var x; (x = 1) || (x = 2); print(x)")
	c.NoInline("var x; x = (1 || (x = 2)); print(x)")
	c.NoInline("var x;(x = 1) && (x = 2); print(x)")
	c.NoInline("var x;x = (1 && (x = 2)); print(x)")
	c.NoInline("var x; x = 1 , x = 2; print(x)")
}

// Java: FlowSensitiveInlineVariablesTest.testNotReachingDefinitions()
func (c *FlowSensitiveInlineVariablesTest) TestNotReachingDefinitions() {
	c.NoInline("var x; if (foo) { x = 0 } print (x)")
}

// Java: FlowSensitiveInlineVariablesTest.testNoInlineLoopCarriedDefinition()
func (c *FlowSensitiveInlineVariablesTest) TestNoInlineLoopCarriedDefinition() {
	// First print is undefined instead.
	c.NoInline("var x; while(true) { print(x); x = 1; }")
	// Prints 0 1 1 1 1....
	c.NoInline("var x = 0; while(true) { print(x); x = 1; }")
}

// Java: FlowSensitiveInlineVariablesTest.testDoNotExitLoop()
func (c *FlowSensitiveInlineVariablesTest) TestDoNotExitLoop() {
	c.NoInline("while (z) { var x = 3; } var y = x;")
}

// Java: FlowSensitiveInlineVariablesTest.testDoNotInlineWithinLoop()
func (c *FlowSensitiveInlineVariablesTest) TestDoNotInlineWithinLoop() {
	c.NoInline("var y = noSFX(); do { var z = y.foo(); } while (true);")
}

// Java: FlowSensitiveInlineVariablesTest.testDoNotInlineCatchExpression1()
func (c *FlowSensitiveInlineVariablesTest) TestDoNotInlineCatchExpression1() {
	c.NoInline(
		`var a;
try {
  throw Error("");
}catch(err) {
   a = err;
}
return a.stack
`)
}

// Java: FlowSensitiveInlineVariablesTest.testDoNotInlineCatchExpression1a()
func (c *FlowSensitiveInlineVariablesTest) TestDoNotInlineCatchExpression1a() {
	c.NoInline(
		`var a;
try {
  throw Error("");
} catch(err) {
   a = err + 1;
}
return a.stack
`)
}

// Java: FlowSensitiveInlineVariablesTest.testDoNotInlineCatchExpression2()
func (c *FlowSensitiveInlineVariablesTest) TestDoNotInlineCatchExpression2() {
	c.NoInline(
		`var a;
try {
  if (x) {throw Error("");}
} catch(err) {
   a = err;
}
return a.stack
`)
}

// Java: FlowSensitiveInlineVariablesTest.testDoNotInlineCatchExpression3()
func (c *FlowSensitiveInlineVariablesTest) TestDoNotInlineCatchExpression3() {
	c.NoInline(
		`var a;
try {
  throw Error("");
} catch(err) {
  err = x;
  a = err;
}
return a.stack
`)
}

// Java: FlowSensitiveInlineVariablesTest.testDoNotInlineCatchExpression4()
func (c *FlowSensitiveInlineVariablesTest) TestDoNotInlineCatchExpression4() {
	// Note: it is valid to inline "x" here but we currently don't.
	c.NoInline(
		`try {
  stuff();
} catch (e) {
  x = e;
  print(x);
}
`)
}

// Java: FlowSensitiveInlineVariablesTest.testDefinitionAfterUse()
func (c *FlowSensitiveInlineVariablesTest) TestDefinitionAfterUse() {
	c.Inline("var x = 0; print(x); x = 1", "var x; print(0); x = 1")
}

// Java: FlowSensitiveInlineVariablesTest.testInlineSameVariableInStraightLine()
func (c *FlowSensitiveInlineVariablesTest) TestInlineSameVariableInStraightLine() {
	c.Inline(
		"var x; x = 1; print(x); x = 2; print(x)", //
		"var x;        print(1);        print(2)")
}

// Java: FlowSensitiveInlineVariablesTest.testInlineInDifferentPaths()
func (c *FlowSensitiveInlineVariablesTest) TestInlineInDifferentPaths() {
	c.Inline(
		"var x; if (print) {x = 1; print(x)} else {x = 2; print(x)}",
		"var x; if (print) {       print(1)} else {       print(2)}")
}

// Java: FlowSensitiveInlineVariablesTest.testNoInlineInMergedPath()
func (c *FlowSensitiveInlineVariablesTest) TestNoInlineInMergedPath() {
	c.NoInline("var x,y;x = 1;while(y) { if(y){ print(x) } else { x = 1 } } print(x)")
}

// Java: FlowSensitiveInlineVariablesTest.testInlineIntoExpressions()
func (c *FlowSensitiveInlineVariablesTest) TestInlineIntoExpressions() {
	c.Inline("var x = 1; print(x + 1);", "var x; print(1 + 1)")
}

// Java: FlowSensitiveInlineVariablesTest.testInlineExpressions1()
func (c *FlowSensitiveInlineVariablesTest) TestInlineExpressions1() {
	c.Inline("var a, b; var x = a+b; print(x)", "var a, b; var x; print(a+b)")
}

// Java: FlowSensitiveInlineVariablesTest.testInlineExpressions2()
func (c *FlowSensitiveInlineVariablesTest) TestInlineExpressions2() {
	// We can't inline because of the redefinition of "a".
	c.NoInline("var a, b; var x = a + b; a = 1; print(x)")
}

// Java: FlowSensitiveInlineVariablesTest.testInlineExpressions3()
func (c *FlowSensitiveInlineVariablesTest) TestInlineExpressions3() {
	c.Inline(
		"var a,b,x; x=a+b; x=a-b; print(  x)", //
		"var a,b,x; x=a+b;        print(a-b)")
}

// Java: FlowSensitiveInlineVariablesTest.testInlineExpressions4()
func (c *FlowSensitiveInlineVariablesTest) TestInlineExpressions4() {
	// Precision is lost due to comma's.
	c.NoInline("var a,b,x; x=a+b, x=a-b; print(x)")
}

// Java: FlowSensitiveInlineVariablesTest.testInlineExpressions5()
func (c *FlowSensitiveInlineVariablesTest) TestInlineExpressions5() {
	c.NoInline("var a; var x = a = 1; print(x)")
}

// Java: FlowSensitiveInlineVariablesTest.testInlineExpressions6()
func (c *FlowSensitiveInlineVariablesTest) TestInlineExpressions6() {
	c.NoInline("var a, x; a = 1 + (x = 1); print(x)")
}

// Java: FlowSensitiveInlineVariablesTest.testInlineExpression7()
func (c *FlowSensitiveInlineVariablesTest) TestInlineExpression7() {
	// Possible side effects in foo() that might conflict with bar();
	c.NoInline("var x = foo() + 1; bar(); print(x)")
	// This is a possible case but we don't have analysis to prove this yet.
	// TODO(user): It is possible to cover this case with the same algorithm
	//                as the missing return check.
	c.NoInline("var x = foo() + 1; print(x)")
}

// Java: FlowSensitiveInlineVariablesTest.testInlineExpression8()
func (c *FlowSensitiveInlineVariablesTest) TestInlineExpression8() {
	// The same variable inlined twice.
	c.Inline(
		"var a,b; var x = a + b; print(    x); x = a - b; print(    x)",
		"var a,b; var x        ; print(a + b);            print(a - b)")
}

// Java: FlowSensitiveInlineVariablesTest.testInlineExpression9()
func (c *FlowSensitiveInlineVariablesTest) TestInlineExpression9() {
	// Check for actual control flow sensitivity.
	c.Inline(
		"var a,b; var x; if (g) { x= a + b; print(    x)}  x = a - b; print(    x)",
		"var a,b; var x; if (g) {           print(a + b)}             print(a - b)")
}

// Java: FlowSensitiveInlineVariablesTest.testInlineExpression10()
func (c *FlowSensitiveInlineVariablesTest) TestInlineExpression10() {
	// The DFA is not fine grain enough for this.
	c.NoInline("var x, y; x = ((y = 1), print(y))")
}

// Java: FlowSensitiveInlineVariablesTest.testInlineExpressions11()
func (c *FlowSensitiveInlineVariablesTest) TestInlineExpressions11() {
	c.Inline("var x; x = x + 1; print(x)", "var x; print(x + 1)")
	c.NoInline("var x; x = x + 1; print(x); print(x)")
}

// Java: FlowSensitiveInlineVariablesTest.testInlineExpressions12()
func (c *FlowSensitiveInlineVariablesTest) TestInlineExpressions12() {
	// ++ is an assignment and considered to modify state so it will not be
	// inlined.
	c.NoInline("var x = 10; x = c++; print(x)")
}

// Java: FlowSensitiveInlineVariablesTest.testInlineExpressions13()
func (c *FlowSensitiveInlineVariablesTest) TestInlineExpressions13() {
	c.Inline(
		`var a = 1, b = 2;
var x = a;
var y = b;
var z = x + y;
var i = z;
var j = z + y;
var k = i;
`,
		`var a, b;
var x;
var y = 2;
var z = 1 + y;
var i;
var j = z + y;
var k = z;
`)
}

// Java: FlowSensitiveInlineVariablesTest.testInlineExpressions14()
func (c *FlowSensitiveInlineVariablesTest) TestInlineExpressions14() {
	c.Inline("var a = function() {}; var b = a;", "var a; var b = function() {}")
}

// Java: FlowSensitiveInlineVariablesTest.testNoInlineIfDefinitionMayNotReach()
func (c *FlowSensitiveInlineVariablesTest) TestNoInlineIfDefinitionMayNotReach() {
	c.NoInline("var x; if (x=1) {} x;")
}

// Java: FlowSensitiveInlineVariablesTest.testNoInlineEscapedToInnerFunction()
func (c *FlowSensitiveInlineVariablesTest) TestNoInlineEscapedToInnerFunction() {
	c.NoInline("var x = 1; function foo() { x = 2 }; print(x)")
}

// Java: FlowSensitiveInlineVariablesTest.testNoInlineLValue()
func (c *FlowSensitiveInlineVariablesTest) TestNoInlineLValue() {
	c.NoInline("var x; if (x = 1) { print(x) }")
}

// Java: FlowSensitiveInlineVariablesTest.testSwitchCase()
func (c *FlowSensitiveInlineVariablesTest) TestSwitchCase() {
	c.Inline("var x = 1; switch(x) { }", "var x; switch(1) { }")
}

// Java: FlowSensitiveInlineVariablesTest.testShadowedVariableInnerFunction()
func (c *FlowSensitiveInlineVariablesTest) TestShadowedVariableInnerFunction() {
	c.Inline(
		"var x = 1; print(x) || (function() {  var x; x = 1; print(x)})()",
		"var x    ; print(1) || (function() {  var x;        print(1)})()")
}

// Java: FlowSensitiveInlineVariablesTest.testCatch()
func (c *FlowSensitiveInlineVariablesTest) TestCatch() {
	c.NoInline("var x = 0; try { } catch (x) { }")
	c.NoInline("try { } catch (x) { print(x) }")
}

// Java: FlowSensitiveInlineVariablesTest.testNoInlineGetProp1()
func (c *FlowSensitiveInlineVariablesTest) TestNoInlineGetProp1() {
	// We don't know if j aliases a.b
	c.NoInline("var x = a.b.c; j.c = 1; print(x);")
}

// Java: FlowSensitiveInlineVariablesTest.testNoInlineGetProp2()
func (c *FlowSensitiveInlineVariablesTest) TestNoInlineGetProp2() {
	c.NoInline("var x = 1 * a.b.c; j.c = 1; print(x);")
}

// Java: FlowSensitiveInlineVariablesTest.testNoInlineGetProp3()
func (c *FlowSensitiveInlineVariablesTest) TestNoInlineGetProp3() {
	// Anything inside a function is fine.
	c.Inline(
		"var a = {b: {}}; var x = function(){1 * a.b.c}; print(x);",
		"var a = {b: {}}; var x; print(function(){1 * a.b.c});")
}

// Java: FlowSensitiveInlineVariablesTest.testNoInlineGetElem()
func (c *FlowSensitiveInlineVariablesTest) TestNoInlineGetElem() {
	// Again we don't know if i = j
	c.NoInline("var x = a[i]; a[j] = 2; print(x); ")
}

// TODO(user): These should be inlinable.
// Java: FlowSensitiveInlineVariablesTest.testNoInlineConstructors()
func (c *FlowSensitiveInlineVariablesTest) TestNoInlineConstructors() {
	c.NoInline("var x = new Iterator(); x.next();")
}

// TODO(user): These should be inlinable.
// Java: FlowSensitiveInlineVariablesTest.testNoInlineArrayLits()
func (c *FlowSensitiveInlineVariablesTest) TestNoInlineArrayLits() {
	c.NoInline("var x = []; print(x)")
}

// TODO(user): These should be inlinable.
// Java: FlowSensitiveInlineVariablesTest.testNoInlineObjectLits()
func (c *FlowSensitiveInlineVariablesTest) TestNoInlineObjectLits() {
	c.NoInline("var x = {}; print(x)")
}

// TODO(user): These should be inlinable after the REGEX checks.
// Java: FlowSensitiveInlineVariablesTest.testNoInlineRegExpLits()
func (c *FlowSensitiveInlineVariablesTest) TestNoInlineRegExpLits() {
	c.NoInline("var x = /y/; print(x)")
}

// Java: FlowSensitiveInlineVariablesTest.testInlineConstructorCallsIntoLoop()
func (c *FlowSensitiveInlineVariablesTest) TestInlineConstructorCallsIntoLoop() {
	// Don't inline construction into loops.
	c.NoInline(
		`var x = new Iterator();
for(i = 0; i < 10; i++) {
  j = x.next();
}
`)
}

// Java: FlowSensitiveInlineVariablesTest.testRemoveWithLabels()
func (c *FlowSensitiveInlineVariablesTest) TestRemoveWithLabels() {
	c.Inline(
		"var x = 1; L: x = 2; print(x)", //
		"var x = 1; L: {    } print(2)")
	c.Inline(
		"var x = 1; L: M: x = 2; print(x)", //
		"var x = 1; L: M: {    } print(2)")
	c.Inline(
		"var x = 1; L: M: N: x = 2; print(x)", //
		"var x = 1; L: M: N: {    } print(2)")
}

// Java: FlowSensitiveInlineVariablesTest.testInlineAcrossSideEffect1()
func (c *FlowSensitiveInlineVariablesTest) TestInlineAcrossSideEffect1() {
	// This can't be inlined because print() has side-effects and might change
	// the definition of noSFX.
	//
	// noSFX must be both const and pure in order to inline it.
	c.NoInline("var y; var x = noSFX(y); print(x)")
	// inline("var y; var x = noSFX(y); print(x)", "var y;var x;print(noSFX(y))");
}

// Java: FlowSensitiveInlineVariablesTest.testInlineAcrossSideEffect2()
func (c *FlowSensitiveInlineVariablesTest) TestInlineAcrossSideEffect2() {
	// Think noSFX() as a function that reads y.foo and return it
	// and SFX() write some new value of y.foo. If that's the case,
	// inlining across hasSFX() is not valid.
	// This is a case where hasSFX is right of the source of the inlining.
	c.NoInline("var y; var x = noSFX(y), z = hasSFX(y); print(x)")
	c.NoInline("var y; var x = noSFX(y), z = new hasSFX(y); print(x)")
	c.NoInline("var y; var x = new noSFX(y), z = new hasSFX(y); print(x)")
}

// Java: FlowSensitiveInlineVariablesTest.testInlineAcrossSideEffect3()
func (c *FlowSensitiveInlineVariablesTest) TestInlineAcrossSideEffect3() {
	// This is a case where hasSFX is left of the destination of the inlining.
	c.NoInline("var y; var x = noSFX(y); hasSFX(y), print(x)")
	c.NoInline("var y; var x = noSFX(y); new hasSFX(y), print(x)")
	c.NoInline("var y; var x = new noSFX(y); new hasSFX(y), print(x)")
}

// Java: FlowSensitiveInlineVariablesTest.testInlineAcrossSideEffect4()
func (c *FlowSensitiveInlineVariablesTest) TestInlineAcrossSideEffect4() {
	// This is a case where hasSFX is some control flow path between the
	// source and its destination.
	c.NoInline("var y; var x = noSFX(y); hasSFX(y); print(x)")
	c.NoInline("var y; var x = noSFX(y); new hasSFX(y); print(x)")
	c.NoInline("var y; var x = new noSFX(y); new hasSFX(y); print(x)")
}

// Java: FlowSensitiveInlineVariablesTest.testCanInlineAcrossNoSideEffect()
func (c *FlowSensitiveInlineVariablesTest) TestCanInlineAcrossNoSideEffect() {
	// This can't be inlined because print() has side-effects and might change
	// the definition of noSFX. We should be able to mark noSFX as const
	// in some way.
	c.NoInline(
		`var y;
var x = noSFX(y),
z = noSFX();
noSFX();
noSFX(),
print(x)
`)
	// inline(
	//    "var y; var x = noSFX(y), z = noSFX(); noSFX(); noSFX(), print(x)",
	//    "var y; var x, z = noSFX(); noSFX(); noSFX(), print(noSFX(y))");
}

// Java: FlowSensitiveInlineVariablesTest.testDependOnOuterScopeVariables()
func (c *FlowSensitiveInlineVariablesTest) TestDependOnOuterScopeVariables() {
	c.NoInline("var x; function foo() { var y = x; x = 0; print(y) }")
	c.NoInline("var x; function foo() { var y = x; x++; print(y) }")
	// Sadly, we don't understand the data flow of outer scoped variables as
	// it can be modified by code outside of this scope. We can't inline
	// at all if the definition has dependence on such variable.
	c.NoInline("var x; function foo() { var y = x; print(y) }")
}

// Java: FlowSensitiveInlineVariablesTest.testInlineIfNameIsLeftSideOfAssign()
func (c *FlowSensitiveInlineVariablesTest) TestInlineIfNameIsLeftSideOfAssign() {
	c.Inline(
		"var x = 1; x = print(x) + 1", //
		"var x    ; x = print(1) + 1")
	c.Inline(
		"var x = 1; L: x = x + 2", //
		"var x    ; L: x = 1 + 2")
	c.Inline(
		"var x = 1; x = (x = x + 1)", //
		"var x    ; x = (x = 1 + 1)")
	c.Inline(
		"// Create a block scope within the function\n"+
			"{\n"+
			"  const C1 = 1;\n"+
			"  const C2 = 2;\n"+
			"// `var` gives `x` a larger scope than `C1`\n"+
			"  var x = C1;\n"+
			"  x = x == C1 ? C1 * 2 : C2 * 2;\n"+
			"}\n"+
			"// `x` still exists here\n"+
			"console.log(x);\n",
		"{\n"+
			"  const C1 = 1;\n"+
			"  const C2 = undefined; // C2 was inlined\n"+
			"  var x = C1;\n"+
			"  x = x == C1 ? C1 * 2 : 2 * 2;\n"+
			"}\n"+
			"// Inlining `C1` to replace `x` here would not work, since `C1` is out of scope here.\n"+
			"console.log(x);\n")
	c.Inline(
		"// Create a block scope within the function\n"+
			"{\n"+
			"  const C1 = 1;\n"+
			"  const C2 = 2;\n"+
			"// `let` gives `x` the same scope as `C1`\n"+
			"  let x = C1;\n"+
			"  x = x == C1 ? C1 * 2 : C2 * 2;\n"+
			"}\n",
		`{
  const C1 = 1;
  const C2 = undefined; // C2 was inlined
  let x; // x was inlined
  x = C1 == C1 ? C1 * 2 : 2 * 2;
}
`)
	c.NoInline("var x = 1; x = (x = (x = 10) + x)")
	c.NoInline("var x = 1; x = (f(x) + (x = 10) + x);")
	c.NoInline("var x = 1; x=-1,foo(x)")
	c.NoInline("var x = 1; x-=1,foo(x)")
}

// Java: FlowSensitiveInlineVariablesTest.testInlineArguments()
func (c *FlowSensitiveInlineVariablesTest) TestInlineArguments() {
	c.TestSame("function _func(x) { print(x) }")
	c.TestSame("function _func(x,y) { if(y) { x = 1 }; print(x) }")
	c.Test(
		"function f(x, y) { x = 1; print(x) }", //
		"function f(x, y) {        print(1) }")
	c.Test(
		"function f(x, y) { if (y) { x = 1; print(x) }}",
		"function f(x, y) { if (y) {        print(1) }}")
}

// Java: FlowSensitiveInlineVariablesTest.testInvalidInlineArguments1()
func (c *FlowSensitiveInlineVariablesTest) TestInvalidInlineArguments1() {
	c.TestSame("function f(x, y) { x = 1; arguments[0] = 2; print(x) }")
	c.TestSame(
		`function f(x, y) {
  x = 1;
  var z = arguments;
  z[0] = 2;
  z[1] = 3;
  print(x);
}
`)
	c.TestSame("function g(a){a[0]=2} function f(x){x=1;g(arguments);print(x)}")
}

// Java: FlowSensitiveInlineVariablesTest.testInvalidInlineArguments2()
func (c *FlowSensitiveInlineVariablesTest) TestInvalidInlineArguments2() {
	c.TestSame(
		`function f(c) {
  var f = c;
  arguments[0] = this;
  f.apply(this, arguments);
  return this;
}
`)
}

// Java: FlowSensitiveInlineVariablesTest.testForIn()
func (c *FlowSensitiveInlineVariablesTest) TestForIn() {
	c.NoInline("var x; var y = {}; for(x in y){}")
	c.NoInline("var x; var y = {}; var z; for(x in z = y){print(z)}")
	c.NoInline("var x; var y = {}; var z; for(x in y){print(z)}")
}

// Java: FlowSensitiveInlineVariablesTest.testForInDestructuring()
func (c *FlowSensitiveInlineVariablesTest) TestForInDestructuring() {
	c.NoInline("var x = 1, y = [], z; for ({z = x} in y) {}")
	c.NoInline("var x = 1, y = [], z; for ([z = x] in y) {}")
	c.NoInline("var x = 1, y = [], z; print(x); for ({z = x} in y) {}")
	c.NoInline("var x = 1, y = [], z; print(x); for ([z = x] in y) {}")
	c.NoInline("var x = 1, y = [], z; print(x); for (let {z = x} in y) {}")
	c.NoInline("var x = 1, y = [], z; print(x); for (const {z = x} in y) {}")
	c.NoInline("var x = 1; if (true) { x = 3; } var y = [[0]], z = x; for ([x] in y) {}; alert(z);")
}

// Java: FlowSensitiveInlineVariablesTest.testNotOkToSkipCheckPathBetweenNodes()
func (c *FlowSensitiveInlineVariablesTest) TestNotOkToSkipCheckPathBetweenNodes() {
	c.NoInline("var x; for(x = 1; foo(x);) {}")
	c.NoInline("var x; for(; x = 1;foo(x)) {}")
}

// Java: FlowSensitiveInlineVariablesTest.testIssue698()
func (c *FlowSensitiveInlineVariablesTest) TestIssue698() {
	// Most of the flow algorithms operate on Vars. We want to make
	// sure the algorithm bails out appropriately if it sees
	// a var that it doesn't know about.
	c.Inline(
		`var x = '';
unknown.length < 2 && (unknown='0' + unknown);
x = x + unknown;
unknown.length < 3 && (unknown='0' + unknown);
x = x + unknown;
return x;
`,
		`var x;
unknown.length < 2 && (unknown='0' + unknown);
x = '' + unknown;
unknown.length < 3 && (unknown='0' + unknown);
x = x + unknown;
return x;
`)
}

// Java: FlowSensitiveInlineVariablesTest.testIssue777()
func (c *FlowSensitiveInlineVariablesTest) TestIssue777() {
	c.Test(
		`function f(cmd, ta) {
  var temp = cmd;
  var temp2 = temp >> 2;
  cmd = STACKTOP;
  for (var src = temp2, dest = cmd >> 2, stop = src + 37;
       src < stop;
       src++, dest++) {
    HEAP32[dest] = HEAP32[src];
  }
  temp = ta;
  temp2 = temp >> 2;
  ta = STACKTOP;
  STACKTOP += 8;
  HEAP32[ta >> 2] = HEAP32[temp2];
  HEAP32[ta + 4 >> 2] = HEAP32[temp2 + 1];
}
`,
		`function f(cmd, ta){
  var temp;
  var temp2 = cmd >> 2;
  cmd = STACKTOP;
  var src = temp2;
  var dest = cmd >> 2;
  var stop = src + 37;
  for(;src<stop;src++,dest++)HEAP32[dest]=HEAP32[src];
  temp2 = ta >> 2;
  ta = STACKTOP;
  STACKTOP += 8;
  HEAP32[ta>>2] = HEAP32[temp2];
  HEAP32[ta+4>>2] = HEAP32[temp2+1];
}
`)
}

// Java: FlowSensitiveInlineVariablesTest.testTransitiveDependencies1()
func (c *FlowSensitiveInlineVariablesTest) TestTransitiveDependencies1() {
	c.Test(
		"function f(x) { var a = x; var b = a; x = 3; return b; }",
		"function f(x) { var a;     var b = x; x = 3; return b; }")
}

// Java: FlowSensitiveInlineVariablesTest.testTransitiveDependencies2()
func (c *FlowSensitiveInlineVariablesTest) TestTransitiveDependencies2() {
	c.Test(
		"function f(x) { var a = x; var b = a; var c = b; x = 3; return c; }",
		"function f(x) { var a    ; var b = x; var c    ; x = 3; return b; }")
}

// Java: FlowSensitiveInlineVariablesTest.testIssue794a()
func (c *FlowSensitiveInlineVariablesTest) TestIssue794a() {
	c.NoInline(
		`var x = 1;
try { x += someFunction(); } catch (e) {}
x += 1;
try { x += someFunction(); } catch (e) {}
return x;
`)
}

// Java: FlowSensitiveInlineVariablesTest.testIssue794b()
func (c *FlowSensitiveInlineVariablesTest) TestIssue794b() {
	c.NoInline(
		`var x = 1;
try { x = x + someFunction(); } catch (e) {}
x = x + 1;
try { x = x + someFunction(); } catch (e) {}
return x;
`)
}

// Java: FlowSensitiveInlineVariablesTest.testVarAssignInsideHookIssue965()
func (c *FlowSensitiveInlineVariablesTest) TestVarAssignInsideHookIssue965() {
	c.NoInline("var i = 0; return 1 ? (i = 5) : 0, i;")
	c.NoInline("var i = 0; return (1 ? (i = 5) : 0) ? i : 0;")
	c.NoInline("var i = 0; return (1 ? (i = 5) : 0) || i;")
	c.NoInline("var i = 0; return (1 ? (i = 5) : 0) * i;")
}

// GitHub issue #250: https://github.com/google/closure-compiler/issues/250
// Java: FlowSensitiveInlineVariablesTest.testInlineStringConcat()
func (c *FlowSensitiveInlineVariablesTest) TestInlineStringConcat() {
	c.Test(
		`function f() {
  var x = '';
  x = x + '1';
  x = x + '2';
  x = x + '3';
  x = x + '4';
  x = x + '5';
  x = x + '6';
  x = x + '7';
  return x;
}
`,
		"function f() { var x; return '' + '1' + '2' + '3' + '4' + '5' + '6' + '7'; }")
}

// Java: FlowSensitiveInlineVariablesTest.testInlineInArrowFunctions()
func (c *FlowSensitiveInlineVariablesTest) TestInlineInArrowFunctions() {
	c.Test(
		"() => {var v; v = 1; return v;} ", //
		"() => {var v       ; return 1;}")
	c.Test(
		"(v) => {v = 1; return v;}", //
		"(v) => {       return 1;}")
}

// Java: FlowSensitiveInlineVariablesTest.testInlineInClassMemberFunctions()
func (c *FlowSensitiveInlineVariablesTest) TestInlineInClassMemberFunctions() {
	c.Test(
		`class C {
  func() {
    var x;
    x = 1;
    return x;
  }
}
`,
		`class C {
  func() {
    var x;
    return 1;
  }
}
`)
}

// Java: FlowSensitiveInlineVariablesTest.testInlineLet()
func (c *FlowSensitiveInlineVariablesTest) TestInlineLet() {
	c.Inline(
		"let a = 1; print(a + 1)", //
		"let a;     print(1 + 1)")
	c.Inline(
		"let a; a = 1; print(a + 1)", //
		"let a;        print(1 + 1)")
	c.NoInline("let a = noSFX(); print(a)")
}

// Java: FlowSensitiveInlineVariablesTest.testInlineConst()
func (c *FlowSensitiveInlineVariablesTest) TestInlineConst() {
	c.Inline(
		"const a = 1        ; print(a + 1)", //
		"const a = undefined; print(1 + 1)")
	c.Inline(
		"const a =         1; const b =         a; print(b + 1)",
		"const a = undefined; const b = undefined; print(1 + 1)")
	c.NoInline("const a = noSFX(); print(a)")
}

// Java: FlowSensitiveInlineVariablesTest.testSpecific()
func (c *FlowSensitiveInlineVariablesTest) TestSpecific() {
	c.Inline(
		"let a = 1; print(a + 1)", //
		"let a    ; print(1 + 1)")
}

// Java: FlowSensitiveInlineVariablesTest.testBlockScoping()
func (c *FlowSensitiveInlineVariablesTest) TestBlockScoping() {
	c.Inline(
		`let a = 1
print(a + 1);
{
  let b = 2;
  print(b + 1);
}
`,
		`let a;
print(1 + 1);
{
  let b;
  print(2 + 1);
}
`)
	c.Inline(
		`let a = 1
{
  let a = 2;
  print(a + 1);
}
print(a + 1);
`,
		`let a = 1
{
  let a;
  print(2 + 1);
}
print(a + 1);
`)
	c.Inline(
		`let a = 1;
  {let b;}
print(a)
`,
		`let a;
  {let b;}
print(1)
`)
	// This test fails to inline due to CheckPathsBetweenNodes analysis in the canInline function
	// in FlowSensitiveInlineVariables.
	c.NoInline(
		`let a = 1;
{
  let b;
  f(b);
}
return(a)
`)
}

// Java: FlowSensitiveInlineVariablesTest.testBlockScoping_shouldntInline()
func (c *FlowSensitiveInlineVariablesTest) TestBlockScoping_shouldntInline() {
	c.NoInline(
		`var JSCompiler_inline_result;
{
  let a = 1;
  if (3 < 4) {
    a = 2;
  }
  JSCompiler_inline_result = a;
}
alert(JSCompiler_inline_result);
`)
	// test let/const shadowing of a var
	c.NoInline(
		`var JSCompiler_inline_result;
var a = 0;
{
  let a = 1;
  if (3 < 4) {
    a = 2;
  }
  JSCompiler_inline_result = a;
}
alert(JSCompiler_inline_result);
`)
	c.NoInline("{ let value = 1; var g = () => value; } return g;")
}

// Java: FlowSensitiveInlineVariablesTest.testInlineInGenerators()
func (c *FlowSensitiveInlineVariablesTest) TestInlineInGenerators() {
	c.Test(
		`function* f() {
  var x = 1;
  return x + 1;
}
`,
		`function* f() {
  var x;
  return 1 + 1;
}
`)
}

// Java: FlowSensitiveInlineVariablesTest.testNoInlineForOf()
func (c *FlowSensitiveInlineVariablesTest) TestNoInlineForOf() {
	c.NoInline("for (var x of n){} ")
	c.NoInline("var x = 1; var n = {}; for(x of n) {}")
}

// Java: FlowSensitiveInlineVariablesTest.testForOfDestructuring()
func (c *FlowSensitiveInlineVariablesTest) TestForOfDestructuring() {
	c.NoInline("var x = 1, y = [], z; for ({z = x} of y) {}")
	c.NoInline("var x = 1, y = [], z; for ([z = x] of y) {}")
	c.NoInline("var x = 1, y = [], z; print(x); for ({z = x} of y) {}")
	c.NoInline("var x = 1, y = [], z; print(x); for ([z = x] of y) {}")
	c.NoInline("var x = 1, y = [], z; print(x); for (let [z = x] of y) {}")
	c.NoInline("var x = 1, y = [], z; print(x); for (const [z = x] of y) {}")
	c.NoInline("var x = 1; if (true) { x = 3; } var y = [[0]], z = x; for ([x] of y) {}; alert(z);")
}

// Java: FlowSensitiveInlineVariablesTest.testTemplateStrings()
func (c *FlowSensitiveInlineVariablesTest) TestTemplateStrings() {
	c.Inline(
		"var name = 'Foo'; `Hello ${ name}`", //
		"var name        ; `Hello ${'Foo'}`")
	c.Inline(
		"var name = 'Foo'; var foo = name; `Hello ${ foo }`", //
		"var name        ; var foo       ; `Hello ${'Foo'}`")
	c.Inline(
		"var age = 3; `Age: ${age}`", //
		"var age    ; `Age: ${  3}`")
}

// Java: FlowSensitiveInlineVariablesTest.testArrayDestructuring()
func (c *FlowSensitiveInlineVariablesTest) TestArrayDestructuring() {
	c.NoInline("var [a, b, c] = [1, 2, 3]; print(a + b + c);")
	c.NoInline("var arr = [1, 2, 3, 4]; var [a, b, ,d] = arr;")
	c.NoInline("var x = 3; [x] = 4; print(x);")
	c.Inline("var [x] = []; x = 3; print(x);", "var [x] = []; print(3);")
}

// Java: FlowSensitiveInlineVariablesTest.testObjectDestructuring()
func (c *FlowSensitiveInlineVariablesTest) TestObjectDestructuring() {
	c.NoInline("var {a, b} = {a: 3, b: 4}; print(a + b);")
	c.NoInline("var obj = {a: 3, b: 4}; var {a, b} = obj;")
}

// Java: FlowSensitiveInlineVariablesTest.testDontInlineOverChangingRvalue_destructuring()
func (c *FlowSensitiveInlineVariablesTest) TestDontInlineOverChangingRvalue_destructuring() {
	c.NoInline("var x = 1; if (true) { x = 2; } var y = x; var [z = (x = 3, 4)] = []; print(y);")
	c.NoInline("var x = 1; if (true) { x = 2; } var y = x; [x] = []; print(y);")
	c.NoInline("var x = 1; if (true) { x = 2; } var y = x; ({x} = {}); print(y);")
	c.NoInline("var x = 1; if (true) { x = 2; } var y = x; var [z] = [x = 3]; print(y);")
}

// Java: FlowSensitiveInlineVariablesTest.testDestructuringDefaultValue()
func (c *FlowSensitiveInlineVariablesTest) TestDestructuringDefaultValue() {
	c.Inline("var x = 1; var [y = x] = [];", "var x; var [y = 1] = [];")
	c.Inline("var x = 1; var {y = x} = {};", "var x; var {y = 1} = {};")
	c.Inline("var x = 1; var {[3]: y = x} = {};", "var x; var {[3]: y = 1} = {};")
	c.NoInline("var x = 1; var {[x]: y = x} = {};")
	c.NoInline("var x = 1; var [y = x] = []; print(x);")
	// don't inline because x is only conditionally reassigned to 2.
	c.NoInline("var x = 1; var [y = (x = 2, 4)] = []; print(x);")
	c.NoInline("var x = 1; print(x); var [y = (x = 2, 4)] = [0]; print(x);")
	// x = 2 is executed before reading x in the default value.
	c.Inline(
		"var x = 1; print(x); var obj = {}; [obj[x = 2] = x] = [];",
		"var x    ; print(1); var obj = {}; [obj[x = 2] = x] = [];")
	// [x] is evaluated before obj[x = 2] is executed
	c.NoInline("var x = 1; print(x); var obj = {}; [[obj[x = 2]] = [x]] = [];")
	c.NoInline("var x = 1; alert(x); ({x = x * 2} = {});")
	c.NoInline("var x = 1; alert(x); [x = x * 2] = [];")
}

// Java: FlowSensitiveInlineVariablesTest.testDestructuringComputedProperty()
func (c *FlowSensitiveInlineVariablesTest) TestDestructuringComputedProperty() {
	c.Inline("var x = 1; var {[x]: y} = {};", "var x; var {[1]: y} = {};")
	c.NoInline("var x = 1; var {[x]: y} = {}; print(x);")
	c.NoInline("var x = 1; alert(x); ({[x]: x} = {}); alert(x);")
	c.Inline("var x = 1; var y = x; ({[y]: x} = {});", "var x; var y; ({[1]: x} = {});")
}

// Java: FlowSensitiveInlineVariablesTest.testDeadAssignments()
func (c *FlowSensitiveInlineVariablesTest) TestDeadAssignments() {
	c.Inline(
		"let a = 3; if (3 < 4) { a = 8; } else { print(a); }",
		"let a; if (3 < 4) { a = 8 } else { print(3); }")
	c.Inline(
		"let a = 3; if (3 < 4) { [a] = 8; } else { print(a); }",
		"let a; if (3 < 4) { [a] = 8 } else { print(3); }")
}

// Java: FlowSensitiveInlineVariablesTest.testDestructuringEvaluationOrder()
func (c *FlowSensitiveInlineVariablesTest) TestDestructuringEvaluationOrder() {
	// Should not inline "x = 2" in these cases because x is changed beforehand
	c.NoInline("var x = 2; var {x, y = x} = {x: 3};")
	c.NoInline("var x = 2; var {y = (x = 3), z = x} = {};")
	// These examples are safe to inline, but FlowSensitiveInlineVariables never inlines variables
	// used twice in the same CFG node even when safe to do so.
	c.NoInline("var x = 2; var {a: y = (x = 3)} = {a: x};")
	c.NoInline("var x = 1; var {a = x} = {a: (x = 2, 3)};")
	c.NoInline("var x = 2; var {a: x = 3} = {a: x};")
	c.NoInline("var x = 1; print(x); var {a: x} = {a: x};")
	c.NoInline("var x = 1; print(x); ({a: x} = {a: x});")
	c.NoInline("var x = 1; print(x); var y; [y = x, x] = [];")
	c.Inline(
		"var x = 1; print(x); var y; [x, y = x] = [2];",
		"var x    ; print(1); var y; [x, y = x] = [2];")
}

// Java: FlowSensitiveInlineVariablesTest.testDestructuringWithSideEffects()
func (c *FlowSensitiveInlineVariablesTest) TestDestructuringWithSideEffects() {
	c.NoInline("function f() { x++; }  var x = 2; var {y = x} = {key: f()}")
	c.NoInline("function f() { x++; } var x = 2; var y = x; var {z = y} = {a: f()}")
	c.NoInline("function f() { x++; } var x = 2; var y = x; var {a = f(), b = y} = {}")
	c.NoInline("function f() { x++; } var x = 2; var y = x; var {[f()]: z = y} = {}")
	c.NoInline("function f() { x++; } var x = 2; var y = x; var {a: {b: z = y} = f()} = {};")
	c.NoInline("function f() { x++; } var x = 2; var y; var {z = (y = x, 3)} = {a: f()}; print(y);")
	c.Inline(
		"function f() { x++; }  var x = 2; var y = x; var {z = f()} = {a: y}",
		"function f() { x++; }  var x = 2; var y    ; var {z = f()} = {a: x}")
	c.Inline(
		"function f() { x++; }  var x = 2; var y = x; var {a = y, b = f()} = {}",
		"function f() { x++; }  var x = 2; var y    ; var {a = x, b = f()} = {}")
	c.Inline(
		"function f() { x++; }  var x = 2; var y = x; var {[y]: z = f()} = {}",
		"function f() { x++; }  var x = 2; var y    ; var {[x]: z = f()} = {}")
	c.Inline(
		"function f() { x++; } var x = 2; var y = x; var {a: {b: z = f()} = y} = {};",
		"function f() { x++; } var x = 2; var y    ; var {a: {b: z = f()} = x} = {};")
}

// Java: FlowSensitiveInlineVariablesTest.testGithubIssue2818()
func (c *FlowSensitiveInlineVariablesTest) TestGithubIssue2818() {
	c.NoInline("var x = 1; var y = x; print(x++, y);")
	c.NoInline("var x = 1; var y = x; print(x = x + 3, y);")
	c.NoInline("var x = 1; var y = x; print(({x} = {x: x * 2}), y); print(x);")
}

// Java: FlowSensitiveInlineVariablesTest.testNoInlineOnOptionalGetProp()
func (c *FlowSensitiveInlineVariablesTest) TestNoInlineOnOptionalGetProp() {
	// b/360959953 - github issue 4187
	// `const t2 = y?.left` should not get inlined into `node.right = t2` (to be
	// `node.right = y?.left`), because the `y?.left` value needs to be stored here
	// in t2, before it is rewritten by `y.left = node`.
	c.TestSame(
		`function swap(node) {
const y = node.right;
const t2 = y?.left;
y.left = node;
node.right = t2;
return node;
}
`)
}

// Java: FlowSensitiveInlineVariablesTest.testNoInlineOnOptionalGetElem()
func (c *FlowSensitiveInlineVariablesTest) TestNoInlineOnOptionalGetElem() {
	// b/360959953 - github issue 4187
	c.TestSame(
		`function swap(node) {
const y = node.right;
const t2 = y?.['foo'];
y['foo'] = node;
node.right = t2;
return node;
}
`)
}

// Java: FlowSensitiveInlineVariablesTest.testNoInlineOnAwait()
func (c *FlowSensitiveInlineVariablesTest) TestNoInlineOnAwait() {
	c.Test("async function f() {var x = 1; print(x) }", "async function f() { var x; print(1) }")
	c.TestSame("async function f() {var x = await 1; print(x) }")
}

// Java: FlowSensitiveInlineVariablesTest.testNoInlineOnYeild()
func (c *FlowSensitiveInlineVariablesTest) TestNoInlineOnYeild() {
	c.Test("function *f() {var x = 1; print(x) }", "function *f() { var x; print(1) }")
	c.TestSame("function *f() {var x = yield 1; print(x) }")
}

// Java: FlowSensitiveInlineVariablesTest.testNoInlineOnClass()
func (c *FlowSensitiveInlineVariablesTest) TestNoInlineOnClass() {
	c.TestSame(
		`function f() {
const x = class {};
const y = x;
}
`)
}

// Java: FlowSensitiveInlineVariablesTest.testNoInlineOnTaggedTemplate()
func (c *FlowSensitiveInlineVariablesTest) TestNoInlineOnTaggedTemplate() {
	c.TestSame(
		"function f() {\n" +
			"var f = (a)=>{};\n" +
			"const x = f`tagged`;\n" +
			"const y = x;\n" +
			"}\n")
}

// Java: FlowSensitiveInlineVariablesTest.testInlineOnOptionalCall()
func (c *FlowSensitiveInlineVariablesTest) TestInlineOnOptionalCall() {
	c.Inline("let x = 1; const y = print(x)", "let x; const y = print(1)")
	c.Inline("let x = 1; const y = print?.(x)", "let x; const y = print?.(1)")
}

// Java: FlowSensitiveInlineVariablesTest.testOkayToInlineWithSideEffects()
func (c *FlowSensitiveInlineVariablesTest) TestOkayToInlineWithSideEffects() {
	c.Inline(
		"var x = 1; var y = x; var z = 1; print(z++, y);",
		"var x    ; var y    ; var z = 1; print(z++, 1);")
	c.Inline(
		"var x = 1; var y = x; var z = 1; print([z] = [], y);",
		"var x    ; var y    ; var z = 1; print([z] = [], 1);")
	c.Inline("var x = 1; var y = x; print(x = 3, y);", "var x; var y; print(x = 3, 1);")
	c.Inline(
		"var x = 1; if (true) { x = 2; } var y = x; var z; z = x = y + 1;",
		"var x = 1; if (true) { x = 2; } var y    ; var z; z = x = x + 1;")
}

// Java: FlowSensitiveInlineVariablesTest.noInline(String)
func (c *FlowSensitiveInlineVariablesTest) NoInline(input string) {
	c.Inline(input, input)
}

// Java: FlowSensitiveInlineVariablesTest.inline(String, String)
func (c *FlowSensitiveInlineVariablesTest) Inline(input string, expected string) {
	c.Test(
		c.Externs(FlowSensitiveInlineVariablesTest_EXTERN_FUNCTIONS),
		c.Srcs(
			jutil.StringReplace(`function _func() {
INPUT
}
`, "INPUT", input)),
		c.Expected(
			jutil.StringReplace(`function _func() {
EXPECTED
}
`, "EXPECTED", expected)))
}
