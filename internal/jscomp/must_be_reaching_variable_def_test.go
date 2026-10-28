// Port of com/google/javascript/jscomp/MustBeReachingVariableDefTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/MustBeReachingVariableDefTest.java)

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
	"fmt"
	"strings"
	"testing"
)

// MustBeReachingVariableDefTest: tests for MustBeReachingVariableDef.
//
// Tests in this class look for the labels 'D:' and 'U:' in the input code and
// check whether the definition of a var `x` at label `D:` is the
// must-reaching definition for its use at label `U:`.
type MustBeReachingVariableDefTest struct {
	// Go-only: test-helper — the Go test running the current test method.
	t *testing.T
}

// Java: MustBeReachingVariableDefTest()
func NewMustBeReachingVariableDefTest() *MustBeReachingVariableDefTest {
	return &MustBeReachingVariableDefTest{}
}

func TestMustBeReachingVariableDefTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewMustBeReachingVariableDefTest() })
}

// JUnitSetT gives the instance the Go test running its test method.
func (c *MustBeReachingVariableDefTest) JUnitSetT(t *testing.T) { c.t = t }

// JUnitBefore: the class has no @Before method.
func (c *MustBeReachingVariableDefTest) JUnitBefore() {}

// JUnitAfter: the class has no @After method.
func (c *MustBeReachingVariableDefTest) JUnitAfter() {}

// Java: MustBeReachingVariableDefTest.testStraightLine()
func (c *MustBeReachingVariableDefTest) TestStraightLine() {
	c.AssertMatch("D:var x=1; U: x")
	c.AssertMatch("var x; D:x=1; U: x")
	c.AssertNotMatch("D:var x=1; x = 2; U: x")
	c.AssertMatch("var x=1; D:x=2; U: x")
	c.AssertNotMatch("U:x; D:var x = 1")
	c.AssertNotMatch("D:var x; U:x; x=1")
	c.AssertNotMatch("D:var x; U:x; x=1; x")
	c.AssertMatch("D: var x = 1; var y = 2; y; U:x")
	c.AssertMatch("let x; D:x=1; U: x")
}

// Java: MustBeReachingVariableDefTest.testIf()
func (c *MustBeReachingVariableDefTest) TestIf() {
	c.AssertNotMatch("var x; if(a){ D:x=1 } else { x=2 }; U:x")
	c.AssertNotMatch("var x; if(a){ x=1 } else { D:x=2 }; U:x")
	c.AssertMatch("D:var x=1; if(a){ U:x } else { x };")
	c.AssertMatch("D:var x=1; if(a){ x } else { U:x };")
	c.AssertNotMatch("var x; if(a) { D: x = 1 }; U:x;")
}

// Java: MustBeReachingVariableDefTest.testLoops()
func (c *MustBeReachingVariableDefTest) TestLoops() {
	c.AssertNotMatch("var x=0; while(a){ D:x=1 }; U:x")
	c.AssertNotMatch("var x=0; for(;;) { D:x=1 }; U:x")
	c.AssertMatch("D:var x=1; while(a) { U:x }")
	c.AssertMatch("D:var x=1; for(;;)  { U:x }")
}

// Java: MustBeReachingVariableDefTest.testConditional()
func (c *MustBeReachingVariableDefTest) TestConditional() {
	c.AssertMatch("var x=0,y; D:(x=1)&&y; U:x")
	c.AssertNotMatch("var x=0,y; D:y&&(x=1); U:x")
	c.AssertNotMatch("D:var x=0; var y; y&&(x=1); U:x")
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: MustBeReachingVariableDefTest.nullishCoalesce()
func (c *MustBeReachingVariableDefTest) TestNullishCoalesce() {
	// LHS is always executed so the definition of x = 1 must be reached
	c.AssertMatch("var x=0,y; D:(x=1)??y; U:x")
	// definitions in RHS are not always executed
	c.AssertNotMatch("var x=0; var y; D:y??(x=1); U:x")
	c.AssertNotMatch("D:var x=0; var y; y??(x=1); U:x")
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: MustBeReachingVariableDefTest.optChain()
func (c *MustBeReachingVariableDefTest) TestOptChain() {
	// LHS is always executed so the definition of x = 1 must be reached
	c.AssertMatch("var x=0,y; D:(x=1)?.y; U:x")
	c.AssertMatch("var x=0,y; D:(x=1)?.[y]; U:x")
	c.AssertMatch("var x=0,y; D:(x=1)?.(y); U:x")
	c.AssertMatch("var x=0,y,z; D:z(x=1)?.y; U:x")
	c.AssertMatch("var x=0,y,z; D:z[x=1]?.y; U:x")
	// definitions in RHS are not always executed
	c.AssertNotMatch("var x = 0,y; D:y?.(x=1); U:x")
	c.AssertNotMatch("D:var x = 0; var y; y?.(x=1); U:x")
	c.AssertNotMatch("var x=0,y; D:y?.[x=1]; U:x")
}

// Java: MustBeReachingVariableDefTest.testUseAndDefInSameInstruction()
func (c *MustBeReachingVariableDefTest) TestUseAndDefInSameInstruction() {
	c.AssertMatch("D:var x=0; U:x=1,x")
	c.AssertMatch("D:var x=0; U:x,x=1")
}

// Java: MustBeReachingVariableDefTest.testAssignmentInExpressions()
func (c *MustBeReachingVariableDefTest) TestAssignmentInExpressions() {
	c.AssertMatch("var x=0; D:foo(bar(x=1)); U:x")
	c.AssertMatch("var x=0; D:foo(bar + (x = 1)); U:x")
}

// Java: MustBeReachingVariableDefTest.testHook()
func (c *MustBeReachingVariableDefTest) TestHook() {
	c.AssertNotMatch("var x=0; D:foo() ? x=1 : bar(); U:x")
	c.AssertNotMatch("var x=0; D:foo() ? x=1 : x=2; U:x")
}

// Java: MustBeReachingVariableDefTest.testExpressionVariableReassignment()
func (c *MustBeReachingVariableDefTest) TestExpressionVariableReassignment() {
	c.AssertMatch("var a,b; D: var x = a + b; U:x")
	c.AssertNotMatch("var a,b,c; D: var x = a + b; a = 1; U:x")
	c.AssertNotMatch("var a,b,c; D: var x = a + b; f(b = 1); U:x")
	c.AssertMatch("var a,b,c; D: var x = a + b; c = 1; U:x")
	// Even if the sub-expression is change conditionally
	c.AssertNotMatch("var a,b,c; D: var x = a + b; c ? a = 1 : 0; U:x")
}

// Java: MustBeReachingVariableDefTest.testMergeDefinitions()
func (c *MustBeReachingVariableDefTest) TestMergeDefinitions() {
	c.AssertNotMatch("var x,y; D: y = x + x; if(x) { x = 1 }; U:y")
}

// Java: MustBeReachingVariableDefTest.testMergesWithOneDefinition()
func (c *MustBeReachingVariableDefTest) TestMergesWithOneDefinition() {
	c.AssertNotMatch("var x,y; while(y) { if (y) { print(x) } else { D: x = 1 } } U:x")
}

// Java: MustBeReachingVariableDefTest.testRedefinitionUsingItself()
func (c *MustBeReachingVariableDefTest) TestRedefinitionUsingItself() {
	c.AssertMatch("var x = 1; D: x = x + 1; U:x;")
	c.AssertNotMatch("var x = 1; D: x = x + 1; x = 1; U:x;")
}

// Java: MustBeReachingVariableDefTest.testMultipleDefinitionsWithDependence()
func (c *MustBeReachingVariableDefTest) TestMultipleDefinitionsWithDependence() {
	c.AssertMatch("var x, a, b; D: x = a, x = b; U: x")
	c.AssertMatch("var x, a, b; D: x = a, x = b; a = 1; U: x")
	c.AssertNotMatch("var x, a, b; D: x = a, x = b; b = 1; U: x")
}

// Java: MustBeReachingVariableDefTest.testExterns()
func (c *MustBeReachingVariableDefTest) TestExterns() {
	c.AssertNotMatch("var goog = {}; D:goog = {}; U: goog")
}

// Java: MustBeReachingVariableDefTest.testAssignmentOp()
func (c *MustBeReachingVariableDefTest) TestAssignmentOp() {
	c.AssertMatch("var x = 0; D: x += 1; U: x")
	c.AssertMatch("var x = 0; D: x *= 1; U: x")
	c.AssertNotMatch("D: var x = 0; x += 1; U: x")
}

// Java: MustBeReachingVariableDefTest.testIncAndDec()
func (c *MustBeReachingVariableDefTest) TestIncAndDec() {
	c.AssertMatch("var x; D: x++; U: x")
	c.AssertMatch("var x; D: x--; U: x")
}

// Java: MustBeReachingVariableDefTest.testFunctionParams1()
func (c *MustBeReachingVariableDefTest) TestFunctionParams1() {
	c.AssertNotMatch("if (x) { D: x = 1; U: x }")
}

// Java: MustBeReachingVariableDefTest.testFunctionParams2()
func (c *MustBeReachingVariableDefTest) TestFunctionParams2() {
	c.AssertNotMatch("if (y) { D: x = 1} U: x")
}

// Java: MustBeReachingVariableDefTest.testArgumentsObjectModifications()
func (c *MustBeReachingVariableDefTest) TestArgumentsObjectModifications() {
	c.AssertNotMatch("D: x = 1; arguments[0] = 2; U: x")
}

// Java: MustBeReachingVariableDefTest.testArgumentsObjectEscaped()
func (c *MustBeReachingVariableDefTest) TestArgumentsObjectEscaped() {
	c.AssertNotMatch("D: x = 1; var y = arguments; y[0] = 2; U: x")
}

// Java: MustBeReachingVariableDefTest.testArgumentsObjectEscapedDependents()
func (c *MustBeReachingVariableDefTest) TestArgumentsObjectEscapedDependents() {
	c.AssertNotMatch("param1=1; var x; D:x=param1; var y=arguments; U:x")
}

// Java: MustBeReachingVariableDefTest.testSideEffects()
func (c *MustBeReachingVariableDefTest) TestSideEffects() {
	c.AssertNotMatch("var a = 1; D: var x = a; a++; U: print(x);")
	// FlowSensitiveInlineVariables needs to handle this case, where a subexpression in the same CFG
	// node as the read of x makes it unsafe to inline "x = a" into print(a++, a);
	c.AssertMatch("var a = 1; D: var x = a; U: print(a++, x);")
}

// Java: MustBeReachingVariableDefTest.testDestructuringDefinitions()
func (c *MustBeReachingVariableDefTest) TestDestructuringDefinitions() {
	c.AssertMatch("D: var [x] = [1]; U: x;")
	c.AssertMatch("D: var x = [1]; U: var [y] = [x];")
	c.AssertMatch("var x = 1; D: [x] = [2]; U: x;")
	c.AssertNotMatch("D: var x = 1; [x] = [2]; U: x;")
	c.AssertNotMatch("var y = 1; D: var x = y; [y] = [2]; U: x;")
	c.AssertMatch("var x; var y; D: [x] = [x = y]; y = 1; U: x;")
	c.AssertMatch("var x; D: var [y] = [x = 1]; U: x;")
}

// Java: MustBeReachingVariableDefTest.testDestructuringDefaultValue()
func (c *MustBeReachingVariableDefTest) TestDestructuringDefaultValue() {
	// conditional definitions of x do not match the usage of x.
	c.AssertNotMatch("var x; D: var [y = x = 3] = []; U: x;")
	c.AssertNotMatch("var x, obj = {}; D: [obj[x = 1] = x = 2] = []; U: x;")
	c.AssertMatch("var x, obj = {}; D: [[obj[x = 1]] = [x = 2]] = []; U: x;")
	c.AssertMatch("var x; D: [x = x = 3] = []; U: x")
}

// Asserts that the use of x at U: is the definition of x at D:.
//
// Java: MustBeReachingVariableDefTest.assertMatch(String)
func (c *MustBeReachingVariableDefTest) AssertMatch(src string) {
	tester := ReachingUseDefTester_Create()
	tester.ComputeReachingDef(src)
	tester.ExtractDefAndUsesFromInputLabels()
	AssertThat(tester.GetComputedDef()).IsSameInstanceAs(tester.GetExtractedDef())
}

// Asserts that the use of x at U: is not the definition of x at D:.
//
// Java: MustBeReachingVariableDefTest.assertNotMatch(String)
func (c *MustBeReachingVariableDefTest) AssertNotMatch(src string) {
	tester := ReachingUseDefTester_Create()
	tester.ComputeReachingDef(src)
	tester.ExtractDefAndUsesFromInputLabels()
	AssertThat(tester.GetComputedDef()).IsNotSameInstanceAs(tester.GetExtractedDef())
}

// Go-only: test — verifies dense variable table and multi-word dependency
// bitsets across <=16, <=64, <=128, >128, and >512 local variables, plus
// escaped-variable dependency invalidation.
func TestMustBeReachingVariableDefWideAndEscaped(t *testing.T) {
	c := &MustBeReachingVariableDefTest{t: t}
	for _, count := range []int{8, 32, 65, 130, 520} {
		t.Run(fmt.Sprintf("vars_%d", count), func(t *testing.T) {
			var decls strings.Builder
			decls.WriteString("var ")
			for i := 0; i < count; i++ {
				if i > 0 {
					decls.WriteByte(',')
				}
				fmt.Fprintf(&decls, "v%d=%d", i, i)
			}
			decls.WriteString(",x;")
			last := count - 1
			mid := count / 2

			// Unmodified high-index dependencies still reach U:x.
			c.AssertMatch(fmt.Sprintf("%s D: x = v0 + v%d + v%d; v1 = 99; U: x", decls.String(), mid, last))
			// Reassigning a high-index dependency invalidates D:x.
			c.AssertNotMatch(fmt.Sprintf("%s D: x = v0 + v%d + v%d; v%d = 99; U: x", decls.String(), mid, last, last))
			// Branch join on x with wide variable table.
			c.AssertNotMatch(fmt.Sprintf("%s if (v0) { D: x = v%d; } else { x = v0; } U: x", decls.String(), last))
			// Parameter escape via arguments invalidates x depending on param1.
			c.AssertNotMatch(fmt.Sprintf("%s D: x = param1 + v%d; arguments[0] = 1; U: x", decls.String(), last))
		})
	}

	// Reassigning an escaped variable invalidates local definitions that read it,
	// both across statements and within the same var declaration.
	c.AssertNotMatch("var esc = 1; function inner() { return esc; } D: var x = esc; esc = 2; U: x")
	c.AssertNotMatch("var esc = 1, x; function inner() { return esc; } D: (x = esc, esc = 2); U: x")
}
