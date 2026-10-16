// Port of com/google/javascript/jscomp/MaybeReachingVariableUseTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/MaybeReachingVariableUseTest.java)

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
)

// MaybeReachingVariableUseTest: tests for MaybeReachingVariableUse.
//
// Each test case consist of a short code snippet that has an instruction
// labeled with `D:` and one or more with label starting with `U:`. When
// `assertMatch` is called, the test suite verifies that all the uses with
// label starting with `U:` is a may be reaching use of the definition label
// at `D:`.
//
// TODO(rishipal): Consider classifying these tests based on the position of
// `D` and `U` in input.
type MaybeReachingVariableUseTest struct {
	// Go-only: test-helper — the Go test running the current test method.
	t *testing.T
}

// Java: MaybeReachingVariableUseTest()
func NewMaybeReachingVariableUseTest() *MaybeReachingVariableUseTest {
	return &MaybeReachingVariableUseTest{}
}

func TestMaybeReachingVariableUseTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewMaybeReachingVariableUseTest() })
}

// JUnitSetT gives the instance the Go test running its test method.
func (c *MaybeReachingVariableUseTest) JUnitSetT(t *testing.T) { c.t = t }

// JUnitBefore: the class has no @Before method.
func (c *MaybeReachingVariableUseTest) JUnitBefore() {}

// JUnitAfter: the class has no @After method.
func (c *MaybeReachingVariableUseTest) JUnitAfter() {}

// Illustrates that the label D needn't correspond to a definition of `x` in these tests, but is
// symbolic of any program point at which we want to check the use U is upward exposed (reachable)
// Java: MaybeReachingVariableUseTest.testReachableUseUAtProgramPointD()
func (c *MaybeReachingVariableUseTest) TestReachableUseUAtProgramPointD() {
	c.AssertMatch("var x = 1; D: var y = 2; y; U:x")
	c.AssertNotMatch("D: var y = 1; var x = 2; y; U:x")
}

// Java: MaybeReachingVariableUseTest.testStraightLine()
func (c *MaybeReachingVariableUseTest) TestStraightLine() {
	c.AssertMatch("D:var x=1; U: x")
	c.AssertMatch("var x; D:x=1; U: x")
	c.AssertNotMatch("D:var x=1; x = 2; U: x")
	c.AssertMatch("var x=1; D:x=2; U: x")
	c.AssertNotMatch("U:x; D:var x = 1")
	c.AssertMatch("D: var x = 1; var y = 2; y; U:x")
}

// Java: MaybeReachingVariableUseTest.testIf()
func (c *MaybeReachingVariableUseTest) TestIf() {
	c.AssertMatch("var x; if(a){ D:x=1 }else { x=2 }; U:x")
	c.AssertMatch("var x; if(a){ x=1 }else { D:x=2 }; U:x")
	c.AssertMatch("D:var x=1; if(a){ U1: x }else { U2: x };")
	// Original def is redefined along all paths to the use, hence not reachable.
	c.AssertNotMatch("D: var x; if(a){ x=1 }else { x=2 }; U:x")
}

// Java: MaybeReachingVariableUseTest.testLoops()
func (c *MaybeReachingVariableUseTest) TestLoops() {
	c.AssertMatch("var x=0; while(a){ D:x=1 }; U:x")
	c.AssertMatch("var x=0; for(;cond;) { D:x=1 }; U:x")
	// statically infinite loops don't have a CFG edge to the end
	c.AssertNotMatch("var x=0; for(;;) { D:x=1 }; U:x")
	c.AssertNotMatch("var x=0; for(;true;) { D:x=1 }; U:x")
	c.AssertNotMatch("var x=0; while (true) { D:x=1 }; U:x")
	// even un-obscured(exposed) defs before the infinite loop don't reach their uses
	c.AssertNotMatch("D: var x=0; while (true) { y=1 }; U:x")
	c.AssertMatch("D:var x=1; while(a) { U:x }")
	c.AssertMatch("D:var x=1; for(;;)  { U:x }")
}

// This test shows that MaybeReachingVariableUseTest does not use data flow(values) in
// conditionals but relies only on static CFG edges to find whether a use reaches a def.
// TODO(rishipal): Make Control flow analysis smarter about short circuiting and update this test.
// Java: MaybeReachingVariableUseTest.testShortCircuiting_usesOnlyCFGEdges()
func (c *MaybeReachingVariableUseTest) TestShortCircuiting_usesOnlyCFGEdges() {
	// Even though `(x=1)` will never execute at runtime in the following cases, it is conditionally
	// executed in  static analysis (=no dataflow) , i.e. there exists a static CFG path from D->U
	// that does not redefine `x`. Hence we add U to the "may-be" reaching use of D.
	c.AssertMatch("var x=0; D: var y = false && (x=1); U:x")
	c.AssertMatch("var x=0; D: var y = true || (x=1); U:x")
	c.AssertMatch("var x=0; var y=0; D:(y=0)&&(x=1); U:x")
	// Even though `(x=1)` will always execute at runtime, it is conditionally executed in  static
	// analysis (=no dataflow) , i.e. there exists a static CFG path from D->U which adds U to the
	// "may-be" reaching use of D.
	c.AssertMatch("D: var x=0; var y = true && (x=1); U:x")
	c.AssertMatch("D: var x=0; var y = false || (x=1); U:x")
}

// Java: MaybeReachingVariableUseTest.testConditional()
func (c *MaybeReachingVariableUseTest) TestConditional() {
	// Def on LHS is unconditional
	c.AssertMatch("var x=0; var y; D:(x=1)&&y; U:x")
	c.AssertNotMatch("D: var x=0; var y; (x=1)&&(y); U:x")
	// Even though `(x=1)` will always execute at runtime, it is conditionally executed in static
	// analysis (=no dataflow) , i.e. there exists a static CFG path from D->U that does not
	// redefine `x`. Hence we add U to the "may-be" reaching use of D.
	c.AssertMatch("D: var x=0; var y=0; (y=1)&&((y=2)||(x=1)); U:x")
	c.AssertMatch("D: var x=0; var y=1; (y)&&(x=1); U:x")
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: MaybeReachingVariableUseTest.nullishCoalesce()
func (c *MaybeReachingVariableUseTest) TestNullishCoalesce() {
	// LHS always executed
	c.AssertMatch("var x=0; var y; D:(x=1)??y; U:x")
	c.AssertNotMatch("D: var x=0; var y; (x=1)??(y); U:x")
	// Even though `(x=1)` will always execute at runtime as `y` is undefined, it is conditionally
	// executed in  static analysis (=no dataflow) , i.e. there exists a static CFG path from D->U
	// that does not redefine `x`. Hence we add U to the "may-be" reaching use of D.
	c.AssertMatch("var x=0; var y; D:y??(x=1); U:x")
	c.AssertMatch("D: var x=0; var y; y??(x=1); U:x")
	c.AssertMatch("D: var x=0; var y; (y)??((y)||(x=1)); U:x")
	// Even though `(x=1)` will never execute at runtime in the following cases, it is conditionally
	// executed in  static analysis (=no dataflow), i.e. there exists a static CFG path from D->U
	// that does not redefine `x`. Hence we add U to the "may-be" reaching use of D.
	c.AssertMatch("var x=0; var y=1; D:(y=1)??(x=1); U:x")
	c.AssertMatch("D: var x=0; var y; (y)&&((y)??(x=1)); U:x")
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: MaybeReachingVariableUseTest.optChain()
func (c *MaybeReachingVariableUseTest) TestOptChain() {
	// LHS always executed
	c.AssertMatch("var x=0; var y; D:(x=1)?.y; U:x")
	c.AssertMatch("var x=0; var y=0; D:(x=1)?.(y=0); U:x")
	c.AssertNotMatch("D: var x=0; var y=0; (x=1)?.(y=0); U:x")
	// one of the paths will reach the definition
	c.AssertMatch("var x=0; var y=0; D:(y=0)?.(x=1); U:x")
	c.AssertMatch("var x=0; var y; D:y?.(x=1); U:x")
	// Flow analysis isn't smart enough to recognize that `(x=1)` here gets short circuited and
	// never executes at runtime
	c.AssertMatch("var x=0; var y; D: y?.(true||(x=1)); U:x")
	// RHS not always executed due to short circuiting
	c.AssertMatch("D: var x=0; var y; y?.((y=2)||(x=1)); U:x")
	c.AssertMatch("D: var x=0; var y; y?.(x=1); U:x")
}

// Java: MaybeReachingVariableUseTest.testUseAndDefInSameInstruction()
func (c *MaybeReachingVariableUseTest) TestUseAndDefInSameInstruction() {
	c.AssertNotMatch("D:var x=0; U:x=1,x")
	c.AssertMatch("D:var x=0; U:x,x=1")
}

// Java: MaybeReachingVariableUseTest.testAssignmentInExpressions()
func (c *MaybeReachingVariableUseTest) TestAssignmentInExpressions() {
	c.AssertMatch("var x=0; D:foo(bar(x=1)); U:x")
	c.AssertMatch("var x=0; D:foo(bar + (x = 1)); U:x")
}

// Java: MaybeReachingVariableUseTest.testHook()
func (c *MaybeReachingVariableUseTest) TestHook() {
	c.AssertMatch("var x=0; D:foo() ? x=1 : bar(); U:x")
	c.AssertMatch("var x=0; D:foo() ? x=1 : x=2; U:x")
	// TODO(rishipal): Fix this test. The U should not be reachable to D as D is obscured by redef.
	c.AssertMatch("D: var x=0; foo() ? x=1 : x=2; U:x")
}

// Java: MaybeReachingVariableUseTest.testAssignmentOps()
func (c *MaybeReachingVariableUseTest) TestAssignmentOps() {
	c.AssertNotMatch("D: var x = 0; U: x = 100")
	c.AssertMatch("D: var x = 0; U: x += 100")
	c.AssertMatch("D: var x = 0; U: x -= 100")
	c.AssertNotMatch("D: var x = 0; x+=10; U:x")
}

// Java: MaybeReachingVariableUseTest.testInc()
func (c *MaybeReachingVariableUseTest) TestInc() {
	c.AssertMatch("D: var x = 0; U:x++")
	c.AssertMatch("var x = 0; D:x++; U:x")
	// TODO(rishipal): Fix this test. The U should not be reachable to D as D is obscured by redef.
	c.AssertMatch("D: var x = 0; x++; U:x")
}

// Java: MaybeReachingVariableUseTest.testForIn()
func (c *MaybeReachingVariableUseTest) TestForIn() {
	// Uses within FOR-IN header are hard to test. They are covered
	// by the tests in the flow sensitive inliner.
	c.AssertNotMatch("D: var x = [], foo; U: for (x in foo) { }")
	c.AssertNotMatch("D: var x = [], foo; for (x in foo) { U:x }")
	c.AssertMatch("var x = [], foo; D: for (x in foo) { U:x }")
	c.AssertMatch("var foo; D: for (let x in foo) { U:x }")
	c.AssertMatch("var foo; D: for (const x in foo) { U:x }")
	c.AssertMatch("D: var x = 1, foo; U: x; U: for (let [z = x] in foo) {}")
	c.AssertMatch("D: var x = 1, foo; U: x; for (let [x] in foo) {}")
}

// Java: MaybeReachingVariableUseTest.testForOf()
func (c *MaybeReachingVariableUseTest) TestForOf() {
	c.AssertNotMatch("D: var x = [], foo; U: for (x of foo) { }")
	c.AssertNotMatch("D: var x = [], foo; for (x of foo) { U:x }")
	c.AssertMatch("var x = [], foo; D: for (x of foo) { U:x }")
	c.AssertMatch("var foo; D: for (let x of foo) { U:x }")
	c.AssertMatch("var foo; D: for (const x of foo) { U:x }")
	c.AssertMatch("D: var x = 1, foo; U: x; U: for (let [z = x] of foo) {}")
	c.AssertMatch("D: var x = 1, foo; U: x; for (let [x] of foo) {}")
}

// Java: MaybeReachingVariableUseTest.testForAwaitOf()
func (c *MaybeReachingVariableUseTest) TestForAwaitOf() {
	c.AssertNotAsyncMatch("D: var x = [], foo; U: for await (x of foo) { }")
	c.AssertNotAsyncMatch("D: var x = [], foo; for await (x of foo) { U:x }")
	c.AssertAsyncMatch("var x = [], foo; D: for await (x of foo) { U:x }")
	c.AssertAsyncMatch("var foo; D: for await (let x of foo) { U:x }")
	c.AssertAsyncMatch("var foo; D: for await (const x of foo) { U:x }")
	c.AssertAsyncMatch("D: var x = 1, foo; U: x; U: for await (let [z = x] of foo) {}")
	c.AssertAsyncMatch("D: var x = 1, foo; U: x; for await (let [x] of foo) {}")
}

// Java: MaybeReachingVariableUseTest.testTryCatch()
func (c *MaybeReachingVariableUseTest) TestTryCatch() {
	c.AssertMatch(
		`D: var x = 1;
try { U: var y = foo() + x; } catch (e) {}
U: var z = x;
`)
	c.AssertMatch(
		`D: var x = 1;
try { x=2; U: var y = foo() + x; } catch (e) {} 
`)
	// TODO(rishipal): Fix this test. The U should not be reachable to D as D is obscured by redef.
	c.AssertMatch(
		`D: var x = 1;
try { x=2; U: var y = foo() + x; } catch (e) {}
U:x;
`)
}

// Java: MaybeReachingVariableUseTest.testDestructuring()
func (c *MaybeReachingVariableUseTest) TestDestructuring() {
	c.AssertMatch("D: var x = 1; U: var [y = x] = [];")
	c.AssertMatch("D: var x = 1; var y; U: [y = x] = [];")
	c.AssertMatch("D: var [x] = []; U: x;")
	c.AssertMatch("var x; x = 3; D: [x] = 5; U: x;")
	c.AssertNotMatch("D: var x; x = 3; [x] = 5; U: x;")
}

// Java: MaybeReachingVariableUseTest.assertMatch(String)
func (c *MaybeReachingVariableUseTest) AssertMatch(src string) {
	c.AssertMatchWithAsync(src, false)
}

// Java: MaybeReachingVariableUseTest.assertAsyncMatch(String)
func (c *MaybeReachingVariableUseTest) AssertAsyncMatch(src string) {
	c.AssertMatchWithAsync(src, true)
}

// The def of `x` at D: may be used by the read of `x` at U:.
//
// Java: MaybeReachingVariableUseTest.assertMatch(String, boolean)
func (c *MaybeReachingVariableUseTest) AssertMatchWithAsync(src string, async bool) {
	tester := ReachingUseDefTester_Create()
	tester.ComputeReachingUses(src, async)
	tester.ExtractDefAndUsesFromInputLabels()
	AssertThat(tester.GetComputedUses()).ContainsAtLeastElementsIn(tester.GetExtractedUses())
}

// Java: MaybeReachingVariableUseTest.assertNotMatch(String)
func (c *MaybeReachingVariableUseTest) AssertNotMatch(src string) {
	c.AssertNotMatchWithAsync(src, false)
}

// Java: MaybeReachingVariableUseTest.assertNotAsyncMatch(String)
func (c *MaybeReachingVariableUseTest) AssertNotAsyncMatch(src string) {
	c.AssertNotMatchWithAsync(src, true)
}

// The def of `x` at D: is not used by the read of `x` at U:.
//
// Java: MaybeReachingVariableUseTest.assertNotMatch(String, boolean)
func (c *MaybeReachingVariableUseTest) AssertNotMatchWithAsync(src string, async bool) {
	tester := ReachingUseDefTester_Create()
	tester.ComputeReachingUses(src, async)
	tester.ExtractDefAndUsesFromInputLabels()
	computedUses := tester.GetComputedUses()
	containsAll := true
	for _, use := range tester.GetExtractedUses() {
		containsAll = containsAll && computedUses.Contains(use)
	}
	AssertThat(containsAll).IsFalse()
}
