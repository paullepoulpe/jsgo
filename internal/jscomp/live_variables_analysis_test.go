// Port of com/google/javascript/jscomp/LiveVariablesAnalysisTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/LiveVariablesAnalysisTest.java)

/*
 * Copyright 2017 The Closure Compiler Authors.
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
	"fmt"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Tests for LiveVariablesAnalysis. Test cases are snippets of a function and
// assertions are made at the instruction labeled with X.
type LiveVariablesAnalysisTest struct {
	t        *testing.T
	liveness *LiveVariablesAnalysis
}

func (c *LiveVariablesAnalysisTest) JUnitSetT(t *testing.T) { c.t = t }
func (c *LiveVariablesAnalysisTest) JUnitBefore()           {}
func (c *LiveVariablesAnalysisTest) JUnitAfter()            {}

// Java: LiveVariablesAnalysisTest()
func NewLiveVariablesAnalysisTest() *LiveVariablesAnalysisTest {
	c := &LiveVariablesAnalysisTest{}
	c.liveness = nil
	return c
}

func TestLiveVariablesAnalysisTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewLiveVariablesAnalysisTest() })
}

// Java: LiveVariablesAnalysisTest.testStraightLine()
func (c *LiveVariablesAnalysisTest) TestStraightLine() {
	// A sample of simple straight line of code with different liveness changes.
	c.AssertNotLiveBeforeX("X:var a;", "a")
	c.AssertNotLiveAfterX("X:var a;", "a")
	c.AssertNotLiveAfterX("X:var a=1;", "a")
	c.AssertLiveAfterX("X:var a=1; a()", "a")
	c.AssertNotLiveBeforeX("X:var a=1; a()", "a")
	c.AssertLiveBeforeX("var a;X:a;", "a")
	c.AssertLiveBeforeX("var a;X:a=a+1;", "a")
	c.AssertLiveBeforeX("var a;X:a+=1;", "a")
	c.AssertLiveBeforeX("var a;X:a++;", "a")
	c.AssertNotLiveAfterX("var a,b;X:b();", "a")
	c.AssertNotLiveBeforeX("var a,b;X:b();", "a")
	c.AssertLiveBeforeX("var a,b;X:b(a);", "a")
	c.AssertLiveBeforeX("var a,b;X:b(1,2,3,b(a + 1));", "a")
	c.AssertNotLiveBeforeX("var a,b;X:a=1;b(a)", "a")
	c.AssertNotLiveAfterX("var a,b;X:b(a);b()", "a")
	c.AssertLiveBeforeX("var a,b;X:b();b=1;a()", "b")
	c.AssertLiveAfterX("X:a();var a;a()", "a")
	c.AssertNotLiveAfterX("X:a();var a=1;a()", "a")
	c.AssertLiveBeforeX("var a,b;X:a,b=1", "a")
}

// Java: LiveVariablesAnalysisTest.testProperties()
func (c *LiveVariablesAnalysisTest) TestProperties() {
	// Reading property of a local variable makes that variable live.
	c.AssertLiveBeforeX("var a,b;X:a.P;", "a")

	// Assigning to a property doesn't kill "a". It makes it live instead.
	c.AssertLiveBeforeX("var a,b;X:a.P=1;b()", "a")
	c.AssertLiveBeforeX("var a,b;X:a.P.Q=1;b()", "a")

	// An "a" in a different context.
	c.AssertNotLiveAfterX("var a,b;X:b.P.Q.a=1;", "a")

	c.AssertLiveBeforeX("var a,b;X:b.P.Q=a;", "a")
}

// Java: LiveVariablesAnalysisTest.testConditions()
func (c *LiveVariablesAnalysisTest) TestConditions() {
	// Reading the condition makes the variable live.
	c.AssertLiveBeforeX("var a,b;X:if(a){}", "a")
	c.AssertLiveBeforeX("var a,b;X:if(a||b) {}", "a")
	c.AssertLiveBeforeX("var a,b;X:if(b||a) {}", "a")
	c.AssertLiveBeforeX("var a,b;X:if(b||b(a)) {}", "a")
	c.AssertNotLiveAfterX("var a,b;X:b();if(a) {}", "b")

	// We can kill within a condition as well.
	c.AssertNotLiveAfterX("var a,b;X:a();if(a=b){}a()", "a")
	c.AssertNotLiveAfterX("var a,b;X:a();while(a=b){}a()", "a")

	// The kill can be "conditional" due to short circuit.
	c.AssertNotLiveAfterX("var a,b;X:a();if((a=b)&&b){}a()", "a")
	c.AssertNotLiveAfterX("var a,b;X:a();while((a=b)&&b){}a()", "a")
	c.AssertLiveBeforeX("var a,b;a();X:if(b&&(a=b)){}a()", "a") // Assumed live.
	c.AssertLiveBeforeX("var a,b;a();X:if(a&&(a=b)){}a()", "a")
	c.AssertLiveBeforeX("var a,b;a();X:while(b&&(a=b)){}a()", "a")
	c.AssertLiveBeforeX("var a,b;a();X:while(a&&(a=b)){}a()", "a")
}

// Java: LiveVariablesAnalysisTest.nullishCoalesce()
func (c *LiveVariablesAnalysisTest) TestNullishCoalesce() {
	// Reading the condition makes the variable live.
	c.AssertLiveBeforeX("var a,b;X:if(a??b) {}", "a")
	c.AssertLiveBeforeX("var a,b;X:if(b??a) {}", "a")
	c.AssertLiveBeforeX("var a,b;X:if(b??b(a)) {}", "a")

	// Unconditionally killed on lhs of ??
	c.AssertNotLiveAfterX("var a,b;X:a();if((a=b)??b){}a()", "a")
	c.AssertNotLiveAfterX("var a,b;X:a();while((a=b)??b){}a()", "a")

	// The kill can be "conditional" due to short circuit.
	c.AssertLiveBeforeX("var a,b; X:if(b??(a=b)){}a()", "a") // Assumed live.
	c.AssertLiveBeforeX("var a,b; X:if(a??(a=b)){}a()", "a")
	c.AssertLiveBeforeX("var a,b; X:while(b??(a=b)){}a()", "a")
	c.AssertLiveBeforeX("var a,b; X:while(a??(a=b)){}a()", "a")
}

// Java: LiveVariablesAnalysisTest.logicalAssignment()
func (c *LiveVariablesAnalysisTest) TestLogicalAssignment() {
	// This pattern is normalized away
	c.AssertLiveBeforeX("var a,b;X:a??=b", "a")
	c.AssertLiveBeforeX("var a,b;X:a??=b", "b")
}

// Java: LiveVariablesAnalysisTest.optChainGetProp()
func (c *LiveVariablesAnalysisTest) TestOptChainGetProp() {
	// Reading the var on lhs of opt chain makes the variable live.
	c.AssertNotLiveBeforeX("var a,b; X:if(b) {}", "a")
	c.AssertLiveBeforeX("var a,b; X:if(a?.b) {}", "a")

	// Reading a prop with the same name as var does not make the var live
	c.AssertNotLiveBeforeX("var a,b;X:if(b?.a) {}", "a")

	// unconditional kill on lhs of ?.
	c.AssertNotLiveAfterX("var a,b;X:a();if((a=c)?.b){} a()", "a")
	c.AssertNotLiveAfterX("var a,b;X:a();while((a=b)?.b){} a()", "a")
}

// Java: LiveVariablesAnalysisTest.optChainCall()
func (c *LiveVariablesAnalysisTest) TestOptChainCall() {
	// conditionally accessing var keeps it live
	c.AssertLiveBeforeX("var a,b; X:if(b?.(a)){}", "a")

	// unconditionally overwriting a var kills it
	c.AssertNotLiveAfterX("var a,b; X:a(); if((a=b)?.b()){} a()", "a")

	// conditionally overwriting var does not kill it
	c.AssertLiveBeforeX("var a,b; X:if(b?.(a=c)){} a();", "a")

	// conditional overwrite on rhs of ?. does not kill the var
	c.AssertLiveBeforeX("var a,b; X:if(b?.(a=b)){}a()", "a") // Assumed live.
	c.AssertLiveBeforeX("var a,b; X:if(a?.(a=b)){}a()", "a")
	c.AssertLiveBeforeX("var a,b; X:while(b?.(a=b)){}a()", "a")
	c.AssertLiveBeforeX("var a,b; X:while(a?.(a=b)){}a()", "a")
}

// Java: LiveVariablesAnalysisTest.optChainGetElem()
func (c *LiveVariablesAnalysisTest) TestOptChainGetElem() {
	// conditionally accessing var keeps it live
	c.AssertLiveBeforeX("var a,b; X:if(b?.[a]) {}", "a")

	// unconditionally overwriting a var kills it
	c.AssertNotLiveAfterX("var a,b; X:a(); if((a=b)?.[b]){} a()", "a")

	// conditionally overwriting var does not kill it
	c.AssertLiveBeforeX("var a,b; X:if(b?.[a=c]) {} a();", "a")

	// conditional overwrite on rhs of ?. does not kill the var
	c.AssertLiveBeforeX("var a,b; X:if(b?.[a=b]){}a()", "a") // Assumed live.
	c.AssertLiveBeforeX("var a,b; X:if(a?.[a=b]){}a()", "a")
	c.AssertLiveBeforeX("var a,b; X:while(b?.[a=b]){}a()", "a")
	c.AssertLiveBeforeX("var a,b; X:while(a?.[a=b]){}a()", "a")
}

// Java: LiveVariablesAnalysisTest.testArrays()
func (c *LiveVariablesAnalysisTest) TestArrays() {
	c.AssertLiveBeforeX("var a;X:a[1]", "a")
	c.AssertLiveBeforeX("var a,b;X:b[a]", "a")
	c.AssertLiveBeforeX("var a,b;X:b[1,2,3,4,b(a)]", "a")
	c.AssertLiveBeforeX("var a,b;X:b=[a,'a']", "a")
	c.AssertNotLiveBeforeX("var a,b;X:a=[];b(a)", "a")

	// Element assignment doesn't kill the array.
	c.AssertLiveBeforeX("var a;X:a[1]=1", "a")
}

// Java: LiveVariablesAnalysisTest.testTwoPaths()
func (c *LiveVariablesAnalysisTest) TestTwoPaths() {
	// Both Paths.
	c.AssertLiveBeforeX("var a,b;X:if(b){b(a)}else{b(a)};", "a")

	// Only one path.
	c.AssertLiveBeforeX("var a,b;X:if(b){b(b)}else{b(a)};", "a")
	c.AssertLiveBeforeX("var a,b;X:if(b){b(a)}else{b(b)};", "a")

	// None of the paths.
	c.AssertNotLiveAfterX("var a,b;X:if(b){b(b)}else{b(b)};", "a")

	// At the very end.
	c.AssertLiveBeforeX("var a,b;X:if(b){b(b)}else{b(b)}a();", "a")

	// The loop might or might not be executed.
	c.AssertLiveBeforeX("var a;X:while(param1){a()};", "a")
	c.AssertLiveBeforeX("var a;X:while(param1){a=1};a()", "a")

	// Same idea with if.
	c.AssertLiveBeforeX("var a;X:if(param1){a()};", "a")
	c.AssertLiveBeforeX("var a;X:if(param1){a=1};a()", "a")

	// This is different in DO. We know for sure at least one iteration is
	// executed.
	c.AssertNotLiveAfterX("X:var a;do{a=1}while(param1);a()", "a")
}

// Java: LiveVariablesAnalysisTest.testThreePaths()
func (c *LiveVariablesAnalysisTest) TestThreePaths() {
	c.AssertLiveBeforeX("var a;X:if(1){}else if(2){}else{a()};", "a")
	c.AssertLiveBeforeX("var a;X:if(1){}else if(2){a()}else{};", "a")
	c.AssertLiveBeforeX("var a;X:if(1){a()}else if(2){}else{};", "a")
	c.AssertLiveBeforeX("var a;X:if(1){}else if(2){}else{};a()", "a")
}

// Java: LiveVariablesAnalysisTest.testHooks()
func (c *LiveVariablesAnalysisTest) TestHooks() {
	c.AssertLiveBeforeX("var a;X:1?a=1:1;a()", "a")

	// Unfortunately, we cannot prove the following because we assume there is
	// no control flow within a hook (i.e. no joins / set unions).
	// assertNotLiveAfterX("var a;X:1?a=1:a=2;a", "a");
	c.AssertLiveBeforeX("var a,b;X:b=1?a:2", "a")
}

// Java: LiveVariablesAnalysisTest.testForLoops()
func (c *LiveVariablesAnalysisTest) TestForLoops() {
	// Induction variable should not be live after the loop.
	c.AssertNotLiveBeforeX("var a,b;for(a=0;a<9;a++){b(a)};X:b", "a")
	c.AssertNotLiveBeforeX("var a,b;for(a in b){a()};X:b", "a")
	c.AssertNotLiveBeforeX("var a,b;for(a in b){a()};X:a", "b")
	c.AssertLiveBeforeX("var b;for(var a in b){X:a()};", "a")

	// It should be live within the loop even if it is not used.
	c.AssertLiveBeforeX("var a,b;for(a=0;a<9;a++){X:1}", "a")
	c.AssertLiveAfterX("var a,b;for(a in b){X:b};", "a")
	// For-In should serve as a gen as well.
	c.AssertLiveBeforeX("var a,b; X:for(a in b){ }", "a")

	// "a in b" should kill "a" before it.
	// Can't prove this unless we have branched backward DFA.
	// assertNotLiveAfterX("var a,b;X:b;for(a in b){a()};", "a");

	// Unless it is used before.
	c.AssertLiveBeforeX("var a,b;X:a();b();for(a in b){a()};", "a")

	// Initializer
	c.AssertLiveBeforeX("var a,b;X:b;for(b=a;;){};", "a")
	c.AssertNotLiveBeforeX("var a,b;X:a;for(b=a;;){b()};b();", "b")
}

// Java: LiveVariablesAnalysisTest.testForOfLoopsVar()
func (c *LiveVariablesAnalysisTest) TestForOfLoopsVar() {
	c.AssertLiveBeforeX("var a; for (a of [1, 2, 3]) {X:{}}", "a")
	c.AssertLiveAfterX("for (var a of [1, 2, 3]) {X:{}}", "a")
	c.AssertLiveBeforeX("var a,b; for (var y of a = [0, 1, 2]) { X:a[y] }", "a")
}

// Java: LiveVariablesAnalysisTest.testForOfLoopsDestructuring()
func (c *LiveVariablesAnalysisTest) TestForOfLoopsDestructuring() {
	c.AssertLiveBeforeX("var key, value; X:for ([key, value] of arr) {value;} value;", "value")
	c.AssertLiveBeforeX("let x = 3; X:for (var [y = x] of arr) { y; }", "x")
	c.AssertLiveBeforeX("for (let [key, value] of arr) { X: key; value; }", "key")
}

// Java: LiveVariablesAnalysisTest.testForAwaitOfLoopsVar()
func (c *LiveVariablesAnalysisTest) TestForAwaitOfLoopsVar() {
	c.AssertLiveBeforeXWithAsync("var a; for await (a of [1, 2, 3]) {X:{}}", "a", true)
	c.AssertLiveAfterXWithAsync("for await (var a of [1, 2, 3]) {X:{}}", "a", true)
	c.AssertLiveBeforeXWithAsync("var a,b; for await (var y of a = [0, 1, 2]) { X:a[y] }", "a", true)
}

// Java: LiveVariablesAnalysisTest.testForAwaitOfLoopsDestructuring()
func (c *LiveVariablesAnalysisTest) TestForAwaitOfLoopsDestructuring() {
	c.AssertLiveBeforeXWithAsync(
		"var key, value; X:for await ([key, value] of arr) {value;} value;", "value", true)
	c.AssertLiveBeforeXWithAsync("let x = 3; X:for await (var [y = x] of arr) { y; }", "x", true)
	c.AssertLiveBeforeXWithAsync("for await (let [key, value] of arr) { X: key; value; }", "key", true)
}

// Java: LiveVariablesAnalysisTest.testNestedLoops()
func (c *LiveVariablesAnalysisTest) TestNestedLoops() {
	c.AssertLiveBeforeX("var a;X:while(1){while(1){a()}}", "a")
	c.AssertLiveBeforeX("var a;X:while(1){while(1){while(1){a()}}}", "a")
	c.AssertLiveBeforeX("var a;X:while(1){while(1){a()};a=1}", "a")
	c.AssertLiveAfterX("var a;while(1){while(1){a()};X:a=1;}", "a")
	c.AssertLiveAfterX("var a;while(1){X:a=1;while(1){a()}}", "a")
	c.AssertNotLiveBeforeX("var a;X:1;do{do{do{a=1;}while(1)}while(1)}while(1);a()", "a")
}

// Java: LiveVariablesAnalysisTest.testSwitches()
func (c *LiveVariablesAnalysisTest) TestSwitches() {
	c.AssertLiveBeforeX("var a,b;X:switch(a){}", "a")
	c.AssertLiveBeforeX("var a,b;X:switch(b){case(a):break;}", "a")
	c.AssertLiveBeforeX("var a,b;X:switch(b){case(b):case(a):break;}", "a")
	c.AssertNotLiveBeforeX("var a,b;X:switch(b){case 1:a=1;break;default:a=2;break};a()", "a")

	c.AssertLiveBeforeX("var a,b;X:switch(b){default:a();break;}", "a")
}

// Java: LiveVariablesAnalysisTest.testAssignAndReadInCondition()
func (c *LiveVariablesAnalysisTest) TestAssignAndReadInCondition() {
	// BUG #1358904
	// Technically, this isn't exactly true....but we haven't model control flow
	// within an instruction.
	c.AssertLiveBeforeX("var a, b; X: if ((a = this) && (b = a)) {}", "a")
	c.AssertNotLiveBeforeX("var a, b; X: a = 1, b = 1;", "a")
	c.AssertNotLiveBeforeX("var a; X: a = 1, a = 1;", "a")
}

// Java: LiveVariablesAnalysisTest.testParam()
func (c *LiveVariablesAnalysisTest) TestParam() {
	// Unused parameter should not be live.
	c.AssertNotLiveAfterX("var a;X:a()", "param1")
	c.AssertLiveBeforeX("var a;X:a(param1)", "param1")
	c.AssertNotLiveAfterX("var a;X:a();a(param2)", "param1")
}

// Java: LiveVariablesAnalysisTest.testExpressionInForIn()
func (c *LiveVariablesAnalysisTest) TestExpressionInForIn() {
	c.AssertLiveBeforeX("var a = [0]; X:for (a[1] in foo) { }", "a")
}

// Java: LiveVariablesAnalysisTest.testArgumentsArray()
func (c *LiveVariablesAnalysisTest) TestArgumentsArray() {
	// Check that use of arguments forces the parameters into the
	// escaped set.
	LiveVariablesAnalysisTest_AssertEscaped("arguments[0]", "param1")
	LiveVariablesAnalysisTest_AssertNotEscaped("arguments[0]", "param2")
	LiveVariablesAnalysisTest_AssertNotEscaped("arguments[0]", "param3")

	LiveVariablesAnalysisTest_AssertEscaped("var args = arguments", "param1")
	LiveVariablesAnalysisTest_AssertNotEscaped("var args = arguments", "param2")
	LiveVariablesAnalysisTest_AssertNotEscaped("var args = arguments", "param3")

	LiveVariablesAnalysisTest_AssertNotEscaped("arguments = []", "param1")
	LiveVariablesAnalysisTest_AssertNotEscaped("arguments = []", "param2")
	LiveVariablesAnalysisTest_AssertNotEscaped("arguments = []", "param3")

	LiveVariablesAnalysisTest_AssertEscaped("arguments[0] = 1", "param1")
	LiveVariablesAnalysisTest_AssertNotEscaped("arguments[0] = 1", "param2")
	LiveVariablesAnalysisTest_AssertNotEscaped("arguments[0] = 1", "param3")

	LiveVariablesAnalysisTest_AssertEscaped("arguments[arguments[0]] = 1", "param1")
	LiveVariablesAnalysisTest_AssertNotEscaped("arguments[arguments[0]] = 1", "param2")
	LiveVariablesAnalysisTest_AssertNotEscaped("arguments[arguments[0]] = 1", "param3")
}

// Java: LiveVariablesAnalysisTest.testArgumentsArray_doesNotEscape_destructuredParams()
func (c *LiveVariablesAnalysisTest) TestArgumentsArray_doesNotEscape_destructuredParams() {
	// These cases also cover a crash related to assuming all RESTs have a NAME child.
	LiveVariablesAnalysisTest_AssertNotEscapedWithWrapper("function f([a]) { arguments; }", "a", LiveVariablesAnalysisTest_Wrapper_NONE)
	LiveVariablesAnalysisTest_AssertNotEscapedWithWrapper("function f([a] = []) { arguments; }", "a", LiveVariablesAnalysisTest_Wrapper_NONE)
	LiveVariablesAnalysisTest_AssertNotEscapedWithWrapper("function f(...[a]) { arguments; }", "a", LiveVariablesAnalysisTest_Wrapper_NONE)
	LiveVariablesAnalysisTest_AssertNotEscapedWithWrapper("function f({a}) { arguments; }", "a", LiveVariablesAnalysisTest_Wrapper_NONE)
	LiveVariablesAnalysisTest_AssertNotEscapedWithWrapper("function f({a} = {}) { arguments; }", "a", LiveVariablesAnalysisTest_Wrapper_NONE)
}

// Java: LiveVariablesAnalysisTest.testTryCatchFinally()
func (c *LiveVariablesAnalysisTest) TestTryCatchFinally() {
	c.AssertLiveAfterX("var a; try {X:a=1} finally {a}", "a")
	c.AssertLiveAfterX("var a; try {a()} catch(e) {X:a=1} finally {a}", "a")
	// Because the outer catch doesn't catch any exceptions at all, the read of
	// "a" within the catch block should not make "a" live.
	c.AssertNotLiveAfterX(
		"var a = 1; try {\n"+
			"try {a()} catch(e) {X:1} } catch(E) {a}\n",
		"a")
	c.AssertLiveAfterX("var a; while(1) { try {X:a=1;break} finally {a}}", "a")
}

// Java: LiveVariablesAnalysisTest.testForInAssignment()
func (c *LiveVariablesAnalysisTest) TestForInAssignment() {
	c.AssertLiveBeforeX("var a,b; for (var y in a = b) { X:a[y] }", "a")
	// No one refers to b after the first iteration.
	c.AssertNotLiveBeforeX("var a,b; for (var y in a = b) { X:a[y] }", "b")
	c.AssertLiveBeforeX("var a,b; for (var y in a = b) { X:a[y] }", "y")
	c.AssertLiveAfterX("var a,b; for (var y in a = b) { a[y]; X: y();}", "a")
}

// Java: LiveVariablesAnalysisTest.testExceptionThrowingAssignments()
func (c *LiveVariablesAnalysisTest) TestExceptionThrowingAssignments() {
	c.AssertLiveBeforeX("try{var a; X:a=foo();a} catch(e) {e()}", "a")
	c.AssertLiveBeforeX("try{X:var a=foo();a} catch(e) {e()}", "a")
	c.AssertLiveBeforeX("try{X:var a=foo()} catch(e) {e(a)}", "a")
}

// Java: LiveVariablesAnalysisTest.testInnerFunctions()
func (c *LiveVariablesAnalysisTest) TestInnerFunctions() {
	c.AssertLiveBeforeX("function a() {}; X: a()", "a")
	c.AssertNotLiveBeforeX("X:; function a() {}", "a")
	c.AssertLiveBeforeX("a = function(){}; function a() {}; X: a()", "a")
	// NOTE: function a() {} has no CFG node representation since it is not
	// part of the control execution.
	c.AssertLiveAfterX("X: a = function(){}; function a() {}; a()", "a")
	c.AssertNotLiveBeforeX("X: a = function(){}; function a() {}; a()", "a")
}

// Java: LiveVariablesAnalysisTest.testEscaped()
func (c *LiveVariablesAnalysisTest) TestEscaped() {
	LiveVariablesAnalysisTest_AssertEscaped("var a;function b(){a()}", "a")
	LiveVariablesAnalysisTest_AssertEscaped("var a;function b(){param1()}", "param1")
	LiveVariablesAnalysisTest_AssertEscaped("var a;function b(){function c(){a()}}", "a")
	LiveVariablesAnalysisTest_AssertEscaped("var a;function b(){param1.x = function() {a()}}", "a")
	LiveVariablesAnalysisTest_AssertNotEscaped("var a;function b(){var c; c()}", "c")
	LiveVariablesAnalysisTest_AssertNotEscaped("var a;function f(){function b(){var c;c()}}", "c")
	LiveVariablesAnalysisTest_AssertNotEscaped("var a;function b(){};a()", "a")
	LiveVariablesAnalysisTest_AssertNotEscaped("var a;function f(){function b(){}}a()", "a")
	LiveVariablesAnalysisTest_AssertNotEscaped("var a;function b(){var a;a()};a()", "a")

	// Escaped by exporting.
	LiveVariablesAnalysisTest_AssertEscaped("var _x", "_x")
}

// Java: LiveVariablesAnalysisTest.testEscapedInClassStaticBlock()
func (c *LiveVariablesAnalysisTest) TestEscapedInClassStaticBlock() {
	LiveVariablesAnalysisTest_AssertEscaped("var a; class C{ static{a()}}", "a")
	LiveVariablesAnalysisTest_AssertEscaped("var a; class C{static{ param1() } }", "param1")
}

// Java: LiveVariablesAnalysisTest.testNotEscapedInClassStaticBlock()
func (c *LiveVariablesAnalysisTest) TestNotEscapedInClassStaticBlock() {
	LiveVariablesAnalysisTest_AssertNotEscaped("var a; class C{static{}} a()", "a")
	LiveVariablesAnalysisTest_AssertNotEscaped("let a; class C{static{let a;a()}}a()", "a")
	LiveVariablesAnalysisTest_AssertNotEscaped("var a; class C{static{var c; c()}}", "c")
}

// ES6 does not require separate handling for catch because the catch block is already recognized
// by the scope creator
// Java: LiveVariablesAnalysisTest.testNotEscapedWithCatch()
func (c *LiveVariablesAnalysisTest) TestNotEscapedWithCatch() {
	LiveVariablesAnalysisTest_AssertEscaped("try{} catch(e){}", "e")
}

// Java: LiveVariablesAnalysisTest.testEscapedLiveness()
func (c *LiveVariablesAnalysisTest) TestEscapedLiveness() {
	c.AssertNotLiveBeforeX("var a;X:a();function b(){a()}", "a")
}

// Java: LiveVariablesAnalysisTest.testBug1449316()
func (c *LiveVariablesAnalysisTest) TestBug1449316() {
	c.AssertLiveBeforeX("try {var x=[]; X:var y=x[0]} finally {foo()}", "x")
}

// Java: LiveVariablesAnalysisTest.testSimpleLet()
func (c *LiveVariablesAnalysisTest) TestSimpleLet() {
	// a is defined after X and not used
	c.AssertNotLiveBeforeDecl("let a;", "a")
	c.AssertNotLiveAfterDecl("let a;", "a")
	c.AssertNotLiveAfterDecl("let a=1;", "a")

	// a is used and defined after X
	c.AssertLiveAfterDecl("let a=1; a()", "a")
	c.AssertNotLiveBeforeDecl("let a=1; a()", "a")

	// no assignment to x; let is initialized with undefined
	c.AssertLiveBeforeX("let a;X:a;", "a")
	c.AssertNotLiveAfterX("let a,b;X:b();", "a")
	c.AssertLiveBeforeX("let a,b;X:b(a);", "a")
	c.AssertNotLiveBeforeX("let a,b;X:a=1;b(a)", "a")
	c.AssertNotLiveAfterX("let a,b;X:b(a);b()", "a")
	c.AssertLiveBeforeX("let a,b;X:b();b=1;a()", "b")

	// let initialized afterX
	c.AssertLiveAfterX("X:a();let a;a()", "a")
	c.AssertNotLiveAfterX("X:a();let a=1;a()", "a")
}

// Java: LiveVariablesAnalysisTest.testLetInnerBlock()
func (c *LiveVariablesAnalysisTest) TestLetInnerBlock() {
	c.AssertNotLiveAfterX("let x; { X:x = 2; let y; }", "x")
}

// Java: LiveVariablesAnalysisTest.testSimpleConst()
func (c *LiveVariablesAnalysisTest) TestSimpleConst() {
	// a is defined after X and not used
	c.AssertLiveBeforeX("const a = 4; X:a;", "a")
	c.AssertNotLiveBeforeDecl("let a = 1;", "a")
	c.AssertNotLiveBeforeDecl("const a = 1;", "a")
	c.AssertNotLiveAfterDecl("const a = 1;", "a")
}

// Java: LiveVariablesAnalysisTest.testArrayDestructuring()
func (c *LiveVariablesAnalysisTest) TestArrayDestructuring() {
	c.AssertLiveBeforeX("var [a, b] = [1, 2]; X:a;", "a")
	c.AssertNotLiveBeforeX("X: var [...a] = f();", "a")
	LiveVariablesAnalysisTest_AssertNotEscaped("var [a, ...b] = [1, 2];", "b")
	LiveVariablesAnalysisTest_AssertNotEscaped("var [a, ...b] = [1, 2];", "a")
	LiveVariablesAnalysisTest_AssertNotEscaped("var [a, ,b] = [1, 2, 3];", "a")
	LiveVariablesAnalysisTest_AssertNotEscaped("var [a, ,b] = [1, 2, 3];", "b")
	c.AssertNotLiveBeforeX("var x = 3; X: [x] = [4]; x;", "x")
	c.AssertLiveBeforeX("var x = {}; X: [x.a] = [3]; x.a;", "x")
	c.AssertLiveBeforeX("var x = []; X: var [c] = x;", "x")
}

// Java: LiveVariablesAnalysisTest.testObjectDestructuring()
func (c *LiveVariablesAnalysisTest) TestObjectDestructuring() {
	c.AssertLiveBeforeX("var {a: x, b: y} = g(); X:x", "x")
	c.AssertNotLiveBeforeX("X: var {a: x, b: y} = g();", "y")
	LiveVariablesAnalysisTest_AssertNotEscaped("var {a: x, b: y} = g()", "x")
	LiveVariablesAnalysisTest_AssertNotEscaped("var {a: x, b: y} = g()", "y")
	LiveVariablesAnalysisTest_AssertNotEscaped("var {a: x = 3, b: y} = g();", "x")
	c.AssertNotLiveBeforeX("var x = {}; X: ({x} = {}); x;", "x")
	c.AssertLiveBeforeX("var x = {}; X: ({a: x.a} = {}); x.a;", "x")
	c.AssertLiveBeforeX("var x = {}; X: var {c} = x;", "x")
}

// Java: LiveVariablesAnalysisTest.testComplexDestructuringPattern()
func (c *LiveVariablesAnalysisTest) TestComplexDestructuringPattern() {
	c.AssertLiveBeforeX("var x = 3; X: var [y = x] = [];", "x")
	c.AssertLiveBeforeX("var x = 3, y; X: [y = x] = [];", "x")
	c.AssertLiveBeforeX("var x = 3; X: var {y = x} = {};", "x")
	c.AssertLiveBeforeX("var x = 3; X: var {key: y = x} = {};", "x")
	c.AssertLiveBeforeX("var x = 3; X: var {[x + x]: foo} = obj; x;", "x")
	c.AssertLiveBeforeX("var x = 3; X: var {[x + x]: x} = obj; x;", "x")

	// TODO(b/538153547): Fix LiveVariablesAnalysis to handle destructuring default values
	// conditionally
	// Nested assignments in default values are conditional and should not kill the variable
	c.AssertNotLiveBeforeX("var a, x, obj; a = 1; X: [x = (a = 2)] = obj; a();", "a")
	c.AssertNotLiveBeforeX("var a, x, obj; a = 1; X: ({x = (a = 2)} = obj); a();", "a")
}

// Java: LiveVariablesAnalysisTest.testComplicatedDeclaration()
func (c *LiveVariablesAnalysisTest) TestComplicatedDeclaration() {
	LiveVariablesAnalysisTest_AssertNotEscaped("var a = 1, {b: b} = f(), c = g()", "a")
	LiveVariablesAnalysisTest_AssertNotEscaped("var a = 1, {b: b} = f(), c = g()", "b")
	LiveVariablesAnalysisTest_AssertNotEscaped("var a = 1, {b: b} = f(), c = g()", "c")
}

func (c *LiveVariablesAnalysisTest) TestLattice128AndVarTableThresholds() {
	// > 16 variables exercises the reused idToIndex map; > 64 variables exercises
	// the second 64-bit word of [2]uint64 in LiveVariableLattice / escaped / untracked.
	var decls string
	for i := 0; i < 80; i++ {
		if i > 0 {
			decls += ", "
		}
		decls += fmt.Sprintf("v%d = %d", i, i)
	}
	src := "var " + decls + "; X: f(); use(v0, v16, v17, v63, v64, v79);"
	state := c.GetFlowStateAtX(src, false)
	AssertThat(state).IsNotNull()
	inLat := state.GetIn()
	for _, liveName := range []string{"v0", "v16", "v17", "v63", "v64", "v79"} {
		idx := c.liveness.GetVarIndex(liveName)
		AssertWithMessage("%s (idx=%d) should be live", liveName, idx).That(inLat.IsLive(idx)).IsTrue()
	}
	for _, deadName := range []string{"v1", "v15", "v18", "v62", "v65", "v78"} {
		idx := c.liveness.GetVarIndex(deadName)
		AssertWithMessage("%s (idx=%d) should be dead", deadName, idx).That(inLat.IsLive(idx)).IsFalse()
	}
	v64Idx := c.liveness.GetVarIndex("v64")
	AssertThat(v64Idx >= 64).IsTrue()
	AssertThat(inLat.NextSetBit(v64Idx)).IsEqualTo(v64Idx)
	copyLat := NewLiveVariablesAnalysis_LiveVariableLatticeOf(inLat)
	AssertThat(inLat.Equals(copyLat)).IsTrue()
	AssertThat(inLat.HashCode()).IsEqualTo(copyLat.HashCode())
	AssertThat(inLat.String()).IsEqualTo(copyLat.String())

	// > 256 block-scoped variables exercises multi-word extra []uint64 and int var indices,
	// including through DeadAssignmentsElimination (where function+body var count <= 100).
	var blockSrc string
	for i := 0; i < 300; i++ {
		blockSrc += fmt.Sprintf("{ let bv%d = %d; if (cond) use(bv%d); } ", i, i, i)
	}
	blockSrc += "{ let tailLive = 1, tailDead = 2; X: f(); use(tailLive); }"
	overState := c.GetFlowStateAtX(blockSrc, false)
	AssertThat(overState).IsNotNull()
	overIn := overState.GetIn()
	liveIdx := c.liveness.GetVarIndex("tailLive")
	deadIdx := c.liveness.GetVarIndex("tailDead")
	AssertThat(liveIdx >= 300).IsTrue()
	AssertThat(overIn.IsLive(liveIdx)).IsTrue()
	AssertThat(overIn.IsLive(deadIdx)).IsFalse()
	AssertThat(overIn.NextSetBit(liveIdx)).IsEqualTo(liveIdx)
	AssertThat(overIn.NextSetBit(liveIdx + 1)).IsEqualTo(-1)
	overCopy := NewLiveVariablesAnalysis_LiveVariableLatticeOf(overIn)
	AssertThat(overIn.Equals(overCopy)).IsTrue()
	AssertThat(overIn.HashCode()).IsEqualTo(overCopy.HashCode())
	AssertThat(overIn.String()).IsEqualTo(overCopy.String())
}

// Java: LiveVariablesAnalysisTest.assertLiveBeforeX(String, String)
func (c *LiveVariablesAnalysisTest) AssertLiveBeforeX(src string, var_ string) {
	c.AssertLiveBeforeXWithAsync(src, var_, false)
}

// Java: LiveVariablesAnalysisTest.assertLiveBeforeX(String, String, boolean)
func (c *LiveVariablesAnalysisTest) AssertLiveBeforeXWithAsync(src string, var_ string, async bool) {
	state := c.GetFlowStateAtX(src, async)
	AssertWithMessage("%s should contain a label 'X:'", src).That(state).IsNotNull()
	AssertWithMessage("Variable %s should be live before X", var_).
		That(state.GetIn().IsLive(c.liveness.GetVarIndex(var_))).
		IsTrue()
}

// Java: LiveVariablesAnalysisTest.assertLiveAfterX(String, String)
func (c *LiveVariablesAnalysisTest) AssertLiveAfterX(src string, var_ string) {
	c.AssertLiveAfterXWithAsync(src, var_, false)
}

// Java: LiveVariablesAnalysisTest.assertLiveAfterX(String, String, boolean)
func (c *LiveVariablesAnalysisTest) AssertLiveAfterXWithAsync(src string, var_ string, async bool) {
	state := c.GetFlowStateAtX(src, async)
	AssertWithMessage("Label X should be in the input program.").That(state).IsNotNull()
	AssertWithMessage("Variable %s should be live after X", var_).
		That(state.GetOut().IsLive(c.liveness.GetVarIndex(var_))).
		IsTrue()
}

// Java: LiveVariablesAnalysisTest.assertNotLiveAfterX(String, String)
func (c *LiveVariablesAnalysisTest) AssertNotLiveAfterX(src string, var_ string) {
	state := c.GetFlowStateAtX(src, false)
	AssertWithMessage("Label X should be in the input program.").That(state).IsNotNull()
	AssertWithMessage("Variable %s should not be live after X", var_).
		That(state.GetOut().IsLive(c.liveness.GetVarIndex(var_))).
		IsFalse()
}

// Java: LiveVariablesAnalysisTest.assertNotLiveBeforeX(String, String)
func (c *LiveVariablesAnalysisTest) AssertNotLiveBeforeX(src string, var_ string) {
	state := c.GetFlowStateAtX(src, false)
	AssertWithMessage("Label X should be in the input program.").That(state).IsNotNull()
	AssertWithMessage("Variable %s should not be live before X", var_).
		That(state.GetIn().IsLive(c.liveness.GetVarIndex(var_))).
		IsFalse()
}

// Java: LiveVariablesAnalysisTest.assertLiveAfterDecl(String, String)
func (c *LiveVariablesAnalysisTest) AssertLiveAfterDecl(src string, var_ string) {
	state := c.GetFlowStateAtDeclaration(src, var_)
	AssertWithMessage("Variable %s should be declared", var_).That(state).IsNotNull()
	AssertWithMessage("Variable%s should be live after its declaration", var_).
		That(state.GetOut().IsLive(c.liveness.GetVarIndex(var_))).
		IsTrue()
}

// Java: LiveVariablesAnalysisTest.assertNotLiveAfterDecl(String, String)
func (c *LiveVariablesAnalysisTest) AssertNotLiveAfterDecl(src string, var_ string) {
	state := c.GetFlowStateAtDeclaration(src, var_)
	AssertWithMessage("Variable %s should be declared", var_).That(state).IsNotNull()
	AssertWithMessage("Variable %s should not be live after its declaration", var_).
		That(state.GetOut().IsLive(c.liveness.GetVarIndex(var_))).
		IsFalse()
}

// Java: LiveVariablesAnalysisTest.assertNotLiveBeforeDecl(String, String)
func (c *LiveVariablesAnalysisTest) AssertNotLiveBeforeDecl(src string, var_ string) {
	state := c.GetFlowStateAtDeclaration(src, var_)
	AssertWithMessage("Variable %s should be declared", var_).That(state).IsNotNull()
	AssertWithMessage("Variable %s should not be live before its declaration", var_).
		That(state.GetIn().IsLive(c.liveness.GetVarIndex(var_))).
		IsFalse()
}

// Java: LiveVariablesAnalysisTest.getFlowStateAtX(String, boolean)
func (c *LiveVariablesAnalysisTest) GetFlowStateAtX(src string, async bool) *DataFlowAnalysis_LinearFlowState[*LiveVariablesAnalysis_LiveVariableLattice] {
	if async {
		c.liveness = LiveVariablesAnalysisTest_ComputeLivenessWithWrapper(src, LiveVariablesAnalysisTest_Wrapper_ASYNC_FUNCTION)
	} else {
		c.liveness = LiveVariablesAnalysisTest_ComputeLivenessWithWrapper(src, LiveVariablesAnalysisTest_Wrapper_FUNCTION)
	}

	return c.GetFlowStateAtXOfNode(c.liveness.GetCfg().GetEntry().GetValue(), c.liveness.GetCfg())
}

// Java: LiveVariablesAnalysisTest.getFlowStateAtX(Node, ControlFlowGraph<Node>)
func (c *LiveVariablesAnalysisTest) GetFlowStateAtXOfNode(node rhino.Node, cfg *ControlFlowGraph[rhino.Node]) *DataFlowAnalysis_LinearFlowState[*LiveVariablesAnalysis_LiveVariableLattice] {
	if node.IsLabel() {
		if node.GetFirstChild().GetString() == "X" {
			return c.liveness.FlowStateAt(cfg.GetLinkedDiGraphNode(node.GetLastChild()).Index)
		}
	}
	for child := node.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		state := c.GetFlowStateAtXOfNode(child, cfg)
		if state != nil {
			return state
		}
	}
	return nil
}

// Java: LiveVariablesAnalysisTest.getFlowStateAtDeclaration(String, String)
func (c *LiveVariablesAnalysisTest) GetFlowStateAtDeclaration(src string, name string) *DataFlowAnalysis_LinearFlowState[*LiveVariablesAnalysis_LiveVariableLattice] {
	c.liveness = LiveVariablesAnalysisTest_ComputeLiveness(src)
	return c.GetFlowStateAtDeclarationOfNode(
		c.liveness.GetCfg().GetEntry().GetValue(), c.liveness.GetCfg(), name)
}

// Use this for lexical declarations which can't be labelled; e.g. `LABEL: let x = 0;` is invalid
// syntax.
//
// Java: LiveVariablesAnalysisTest.getFlowStateAtDeclaration(Node, ControlFlowGraph<Node>, String)
func (c *LiveVariablesAnalysisTest) GetFlowStateAtDeclarationOfNode(node rhino.Node, cfg *ControlFlowGraph[rhino.Node], name string) *DataFlowAnalysis_LinearFlowState[*LiveVariablesAnalysis_LiveVariableLattice] {
	if NodeUtil_IsNameDeclaration(node) {
		if node.GetFirstChild().GetString() == name {
			return c.liveness.FlowStateAt(cfg.GetLinkedDiGraphNode(node).Index)
		}
	}
	for child := node.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		state := c.GetFlowStateAtDeclarationOfNode(child, cfg, name)
		if state != nil {
			return state
		}
	}
	return nil
}

// Java: LiveVariablesAnalysisTest.assertEscaped(String, String)
func LiveVariablesAnalysisTest_AssertEscaped(src string, name string) {
	for var_ := range LiveVariablesAnalysisTest_ComputeLiveness(src).GetEscapedLocals().All() {
		if var_.GetName() == name {
			return
		}
	}
	AssertWithMessage("Variable %s should be in the escaped local list.", name).Fail()
}

// Java: LiveVariablesAnalysisTest.assertNotEscaped(String, String)
func LiveVariablesAnalysisTest_AssertNotEscaped(src string, name string) {
	LiveVariablesAnalysisTest_AssertNotEscapedWithWrapper(src, name, LiveVariablesAnalysisTest_Wrapper_FUNCTION)
}

// wrapper: the kind of wrapper function the src will be placed in.
//
// Java: LiveVariablesAnalysisTest.assertNotEscaped(String, String, Wrapper)
func LiveVariablesAnalysisTest_AssertNotEscapedWithWrapper(src string, name string, wrapper LiveVariablesAnalysisTest_Wrapper) {
	for var_ := range LiveVariablesAnalysisTest_ComputeLivenessWithWrapper(src, wrapper).GetEscapedLocals().All() {
		AssertThat(var_.GetName()).IsNotEqualTo(name)
	}
}

// Wraps src in a function and computes a LiveVariablesAnalysis.
//
// Java: LiveVariablesAnalysisTest.computeLiveness(String)
func LiveVariablesAnalysisTest_ComputeLiveness(src string) *LiveVariablesAnalysis {
	return LiveVariablesAnalysisTest_ComputeLivenessWithWrapper(src, LiveVariablesAnalysisTest_Wrapper_FUNCTION)
}

// Optionally wraps the src in a function and computes a LiveVariablesAnalysis.
//
// Java: LiveVariablesAnalysisTest.computeLiveness(String, Wrapper)
func LiveVariablesAnalysisTest_ComputeLivenessWithWrapper(src string, wrapper LiveVariablesAnalysisTest_Wrapper) *LiveVariablesAnalysis {
	// Set up compiler
	compiler := NewCompiler()
	options := NewCompilerOptions()
	options.SetLanguage(CompilerOptions_LanguageMode_UNSUPPORTED)
	options.SetCodingConvention(NewGoogleCodingConvention())
	compiler.InitOptions(options)
	compiler.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_NORMALIZED)

	// Set up test case
	switch wrapper {
	case LiveVariablesAnalysisTest_Wrapper_FUNCTION:
		src = "function _FUNCTION(param1, param2 = 1, ...param3){" + src + "}"
	case LiveVariablesAnalysisTest_Wrapper_ASYNC_FUNCTION:
		src = "async function _FUNCTION(param1, param2 = 1, ...param3){" + src + "}"
	default:
	}
	n := compiler.ParseTestCode(src).RemoveFirstChild()
	jutil.CheckStateMsg1(n.IsFunction(), "%s", n)
	script := rhino.NewNodeWithChild(rhino.Token_SCRIPT, n)
	script.SetInputId(rhino.NewInputId("test"))
	AssertThat(compiler.GetErrors()).IsEmpty()

	// Create scopes
	scopeCreator := NewSyntacticScopeCreator(compiler)
	scope := scopeCreator.CreateScopeOfScope(n, Scope_CreateGlobalScope(script))
	childScope := scopeCreator.CreateScopeOfScope(NodeUtil_GetFunctionBody(n), scope)

	// Control flow graph
	cfg := ComputeControlFlowGraph(compiler, n, false, true)

	// All variables declared in function
	allVarsDeclaredInFunction := NodeUtil_GetAllVarsDeclaredInFunction(compiler, scopeCreator, scope)

	// Compute liveness of variables
	analysis := NewLiveVariablesAnalysis(
		cfg, scope, childScope, compiler, scopeCreator, allVarsDeclaredInFunction)
	analysis.Analyze()
	return analysis
}

func NewLiveVariablesAnalysis(
	cfg *ControlFlowGraph[rhino.Node],
	jsScope *Scope,
	jsScopeChild *Scope,
	compiler AbstractCompiler,
	scopeCreator ScopeCreator,
	allVarsDeclaredInFunction *NodeUtil_AllVarsDeclaredInFunction,
) *LiveVariablesAnalysis {
	l := &LiveVariablesAnalysis{}
	l.Reset(cfg, jsScope, jsScopeChild, compiler, scopeCreator, allVarsDeclaredInFunction)
	return l
}

func (l *LiveVariablesAnalysis) GetEscapedLocals() *jutil.ObjectLinkedHashSet[*Var] {
	out := jutil.NewObjectLinkedHashSet[*Var]()
	l.escaped.ForEach(func(idx int) {
		if idx < len(l.table.vars) {
			out.Add(l.table.vars[idx])
		}
	})
	if l.paramsMarkedEscaped {
		params := ast.MustFunction(l.jsScope.GetRootNode()).Params()
		for param, ok := params.FirstParam(); ok; param, ok = param.Next() {
			if param.Token() == rhino.Token_NAME {
				if pv := l.jsScope.GetVarByID(param.Node().GetStringID()); pv != nil {
					out.Add(pv)
				}
			}
		}
	}
	return out
}

// LiveVariablesAnalysisTest_Wrapper: the kinds of functions that source code
// can be wrapped in for testing.
type LiveVariablesAnalysisTest_Wrapper int32

const (
	LiveVariablesAnalysisTest_Wrapper_NONE LiveVariablesAnalysisTest_Wrapper = iota
	// e.g. <code>function f(){ ... }</code>
	LiveVariablesAnalysisTest_Wrapper_FUNCTION
	// e.g. <code>async function f(){ ... }</code>
	LiveVariablesAnalysisTest_Wrapper_ASYNC_FUNCTION
)
