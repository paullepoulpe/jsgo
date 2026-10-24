// Port of com/google/javascript/jscomp/PeepholeCollectPropertyAssignmentsTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/PeepholeCollectPropertyAssignmentsTest.java)

/*
 * Copyright 2010 The Closure Compiler Authors.
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

// PeepholeCollectPropertyAssignmentsTest is the unit tests for
// PeepholeCollectPropertyAssignments.
type PeepholeCollectPropertyAssignmentsTest struct {
	CompilerTestCaseBase
}

// Java: PeepholeCollectPropertyAssignmentsTest()
func NewPeepholeCollectPropertyAssignmentsTest() *PeepholeCollectPropertyAssignmentsTest {
	c := &PeepholeCollectPropertyAssignmentsTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestPeepholeCollectPropertyAssignmentsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewPeepholeCollectPropertyAssignmentsTest() })
}

// Java: PeepholeCollectPropertyAssignmentsTest.getProcessor(Compiler)
func (c *PeepholeCollectPropertyAssignmentsTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewPeepholeOptimizationsPass(
		compiler, c.GetName(), NewPeepholeCollectPropertyAssignments())
}

// Java: PeepholeCollectPropertyAssignmentsTest.test36122565a()
func (c *PeepholeCollectPropertyAssignmentsTest) Test36122565a() {
	c.TestSame(
		"var foo = { bar: g(), baz: 4 };\n" +
			"foo.bar = 3;\n" +
			"foo.baz = 3;\n" +
			"console.log(foo.bar);\n" +
			"console.log(foo.baz);\n")
	c.Test(
		"var foo = { bar: g(), baz: 4 };\n"+
			"foo.baz = 3;\n"+
			"foo.bar = 3;\n"+
			"console.log(foo.bar);\n"+
			"console.log(foo.baz);\n",
		"var foo = { bar: g(), baz: 3 };\n"+
			"foo.bar = 3;\n"+
			"console.log(foo.bar);\n"+
			"console.log(foo.baz);\n")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testComputedPropertyShadowing()
func (c *PeepholeCollectPropertyAssignmentsTest) TestComputedPropertyShadowing() {
	c.TestSame(
		`var called = false;
var x = { set ['a'](v) { called = true; } };
x.a = 2;
`)
	c.TestSame(
		`var called = false;
var x = { get ['a']() { return 1; } };
x.a = 2;
`)
}

// Java: PeepholeCollectPropertyAssignmentsTest.testObjectSpreadDoesNotBlockFolding()
func (c *PeepholeCollectPropertyAssignmentsTest) TestObjectSpreadDoesNotBlockFolding() {
	c.Test("var x = { ...c }; x.a = 1;", "var x = { ...c, a: 1 };")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testArrayOptimization1()
func (c *PeepholeCollectPropertyAssignmentsTest) TestArrayOptimization1() {
	c.Test("var a = []; a[0] = 1; a[1] = 2; a[2] = 3;", "var a = [1, 2, 3];")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testArray_nonConsecutiveAssignments()
func (c *PeepholeCollectPropertyAssignmentsTest) TestArray_nonConsecutiveAssignments() {
	c.Test(
		"var a = []; a[2] = 1; a[0] = 2; a[1] = 3;", //
		"var a = [2, 3, 1];")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testArrayOptimization2()
func (c *PeepholeCollectPropertyAssignmentsTest) TestArrayOptimization2() {
	c.Test("var a; a = []; a[0] = 1; a[1] = 2; a[2] = 3;", "var a; a = [1, 2, 3];")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testArrayOptimization3()
func (c *PeepholeCollectPropertyAssignmentsTest) TestArrayOptimization3() {
	c.TestSame("var a; a.b = []; a.b[0] = 1; a.b[1] = 2; a.b[2] = 3;")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testArrayOptimizationWithLet()
func (c *PeepholeCollectPropertyAssignmentsTest) TestArrayOptimizationWithLet() {
	c.Test("let a = []; a[0] = 1; a[1] = 2; a[2] = 3;", "let a = [1, 2, 3];")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testArrayOptimizationWithConst()
func (c *PeepholeCollectPropertyAssignmentsTest) TestArrayOptimizationWithConst() {
	c.Test("const a = []; a[0] = 1; a[1] = 2; a[2] = 3;", "const a = [1, 2, 3];")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testCompoundAssignment()
func (c *PeepholeCollectPropertyAssignmentsTest) TestCompoundAssignment() {
	c.TestSame("var x, a; a = []; a[0] *= x;")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testNegativeArrayIndex1()
func (c *PeepholeCollectPropertyAssignmentsTest) TestNegativeArrayIndex1() {
	c.TestSame("var a = []; a[-1] = 1;")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testNegativeArrayIndex2()
func (c *PeepholeCollectPropertyAssignmentsTest) TestNegativeArrayIndex2() {
	c.TestSame("var a; a = []; a[-1] = 1;")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testFractionalArrayIndex1()
func (c *PeepholeCollectPropertyAssignmentsTest) TestFractionalArrayIndex1() {
	c.TestSame("var a = []; a[0.5] = 1;")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testFractionalArrayIndex2()
func (c *PeepholeCollectPropertyAssignmentsTest) TestFractionalArrayIndex2() {
	c.TestSame("var a; a = []; a[0.5] = 1;")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testArrayOptimizationOfPartiallyBuiltArray1()
func (c *PeepholeCollectPropertyAssignmentsTest) TestArrayOptimizationOfPartiallyBuiltArray1() {
	c.Test("var a = [1, 2]; a[2] = 3;", "var a = [1, 2, 3];")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testArrayOptimizationOfPartiallyBuiltArray2()
func (c *PeepholeCollectPropertyAssignmentsTest) TestArrayOptimizationOfPartiallyBuiltArray2() {
	c.Test("var a; a = [1, 2]; a[2] = 3;", "var a; a = [1, 2, 3];")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testArrayOptimizationWithAHole1()
func (c *PeepholeCollectPropertyAssignmentsTest) TestArrayOptimizationWithAHole1() {
	c.Test("var a = []; a[0] = 1; a[1] = 2; a[3] = 4;", "var a = [1, 2, , 4];")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testArrayOptimizationWithAHole2()
func (c *PeepholeCollectPropertyAssignmentsTest) TestArrayOptimizationWithAHole2() {
	c.Test("var a; a = []; a[0] = 1; a[1] = 2; a[3] = 4;", "var a; a = [1, 2, , 4];")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testEarlyUsage1()
func (c *PeepholeCollectPropertyAssignmentsTest) TestEarlyUsage1() {
	c.TestSame(
		"function c() {return sum(a)};\n" +
			"var a = [1,2,3];\n" +
			"a[4] = c();\n")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testEarlyUsage2()
func (c *PeepholeCollectPropertyAssignmentsTest) TestEarlyUsage2() {
	c.TestSame(
		"function c() {return sum(a)};\n" +
			"var a; a = [1,2,3];\n" +
			"a[4] = c();\n")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testArrayTooSparseOptimization1()
func (c *PeepholeCollectPropertyAssignmentsTest) TestArrayTooSparseOptimization1() {
	c.Test("var a = []; a[0] = 1; a[1] = 2; a[100] = 4;", "var a = [1, 2]; a[100] = 4;")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testArrayTooSparseOptimization2()
func (c *PeepholeCollectPropertyAssignmentsTest) TestArrayTooSparseOptimization2() {
	c.Test("var a; a = []; a[0] = 1; a[1] = 2; a[100] = 4;", "var a; a = [1, 2]; a[100] = 4;")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testArrayOutOfOrder()
func (c *PeepholeCollectPropertyAssignmentsTest) TestArrayOutOfOrder() {
	c.Test("var a = []; a[1] = 1; a[0] = 0;", "var a = [0, 1];")
	c.Test("var a; a = []; a[1] = 1; a[0] = 0;", "var a; a = [0, 1];")
	// We cannot change the order of side-effects.
	// The below should not be
	//   var x = 0; var a = [x++, x++]
	// since that would produce
	//   var a = [0, 1], x = 2;
	// instead of
	//   var a = [1, 0], x = 2;
	c.TestSame("var x = 0; var a = []; a[1] = x++; a[0] = x++;")
	c.TestSame("var x; x = 0; var a = []; a[1] = x++; a[0] = x++;")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testMultipleNames1()
func (c *PeepholeCollectPropertyAssignmentsTest) TestMultipleNames1() {
	c.Test("var b = []; b[0] = 2; var a = []; a[0] = 1;", "var b = [2]; var a = [1];")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testMultipleNames2()
func (c *PeepholeCollectPropertyAssignmentsTest) TestMultipleNames2() {
	c.Test("var b; b = []; b[0] = 2; var a = []; a[0] = 1;", "var b; b = [2]; var a = [1];")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testArrayReassignedInValue1()
func (c *PeepholeCollectPropertyAssignmentsTest) TestArrayReassignedInValue1() {
	c.Test(
		"var a = []; a[0] = 1; a[1] = (a = []); a[3] = 4;",
		"var a = [1]; a[1] = (a = []); a[3] = 4;")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testArrayReassignedInValue2()
func (c *PeepholeCollectPropertyAssignmentsTest) TestArrayReassignedInValue2() {
	c.Test(
		"var a; a = []; a[0] = 1; a[1] = (a = []); a[3] = 4;",
		"var a; a = [1]; a[1] = (a = []); a[3] = 4;")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testArrayReassignedInSubsequentVar1()
func (c *PeepholeCollectPropertyAssignmentsTest) TestArrayReassignedInSubsequentVar1() {
	c.TestSame("var a = []; a[0] = a = []; a[1] = 2;")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testArrayReassignedInSubsequentVar2()
func (c *PeepholeCollectPropertyAssignmentsTest) TestArrayReassignedInSubsequentVar2() {
	c.TestSame("var a; a = []; a[0] = a = []; a[1] = 2;")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testForwardReference1()
func (c *PeepholeCollectPropertyAssignmentsTest) TestForwardReference1() {
	c.Test("var a; a = []; a[0] = 1; a[1] = a;", "var a; a = [1]; a[1] = a;")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testForwardReference2()
func (c *PeepholeCollectPropertyAssignmentsTest) TestForwardReference2() {
	c.Test("var a; a = []; a[0] = 1; a[1] = a;", "var a; a = [1]; a[1] = a;")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testArray_usingSpread_diablesOptimization_leadingSpread()
func (c *PeepholeCollectPropertyAssignmentsTest) TestArray_usingSpread_diablesOptimization_leadingSpread() {
	// TODO(nickreid): A future improvement could generate `var a = [3, ...b, 2];`
	c.TestSame("var a = [1, ...b, 2]; a[0] = 3;")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testArray_usingSpread_diablesOptimization_withinSpread()
func (c *PeepholeCollectPropertyAssignmentsTest) TestArray_usingSpread_diablesOptimization_withinSpread() {
	c.TestSame("var a = [1, ...b, 2]; a[1] = 3;")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testArray_usingSpread_diablesOptimization_followingSpread()
func (c *PeepholeCollectPropertyAssignmentsTest) TestArray_usingSpread_diablesOptimization_followingSpread() {
	c.TestSame("var a = [1, ...b, 2]; a[2] = 3;")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testObjectOptimization1()
func (c *PeepholeCollectPropertyAssignmentsTest) TestObjectOptimization1() {
	c.Test("var o = {}; o.x = 0; o['y'] = 1; o[2] = 2;", "var o = { x: 0, \"y\": 1, \"2\": 2 };")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testObjectOptimization2()
func (c *PeepholeCollectPropertyAssignmentsTest) TestObjectOptimization2() {
	c.Test(
		"var o; o = {}; o.x = 0; o['y'] = 1; o[2] = 2;",
		"var o; o = { x: 0, \"y\": 1, \"2\": 2 };")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testObjectOptimizationWithLet()
func (c *PeepholeCollectPropertyAssignmentsTest) TestObjectOptimizationWithLet() {
	c.Test("let o = {}; o.x = 0; o['y'] = 1; o[2] = 2;", "let o = { x: 0, 'y': 1, '2': 2 };")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testObjectOptimizationWithConst()
func (c *PeepholeCollectPropertyAssignmentsTest) TestObjectOptimizationWithConst() {
	c.Test("const o = {}; o.x = 0; o['y'] = 1; o[2] = 2;", "const o = { x: 0, 'y': 1, '2': 2 };")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testObjectReassignedInValue1()
func (c *PeepholeCollectPropertyAssignmentsTest) TestObjectReassignedInValue1() {
	c.Test(
		"var o = {}; o.x = 1; o.y = (o = {}); o.z = 4;", "var o = {x:1}; o.y = (o = {}); o.z = 4;")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testObjectReassignedInValue2()
func (c *PeepholeCollectPropertyAssignmentsTest) TestObjectReassignedInValue2() {
	c.Test(
		"var o; o = {}; o.x = 1; o.y = (o = {}); o.z = 4;",
		"var o; o = {x:1}; o.y = (o = {}); o.z = 4;")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testObjectFunctionRollup1()
func (c *PeepholeCollectPropertyAssignmentsTest) TestObjectFunctionRollup1() {
	c.Test(
		"var o; o = {};\n"+
			"o.x = function() {};\n",
		"var o; o = {x:function () {}};")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testObjectFunctionRollup2()
func (c *PeepholeCollectPropertyAssignmentsTest) TestObjectFunctionRollup2() {
	c.TestSame(
		"var o; o = {};\n" +
			"o.x = (function() {return o})();\n")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testObjectFunctionRollup3()
func (c *PeepholeCollectPropertyAssignmentsTest) TestObjectFunctionRollup3() {
	c.Test(
		"var o; o = {};\n"+
			"o.x = function() {return o};\n",
		"var o; o = {x:function () {return o}};")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testObjectFunctionRollup4()
func (c *PeepholeCollectPropertyAssignmentsTest) TestObjectFunctionRollup4() {
	c.TestSame(
		"function f() {return o};\n" +
			"var o; o = {};\n" +
			"o.x = f();\n")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testObjectFunctionRollup5()
func (c *PeepholeCollectPropertyAssignmentsTest) TestObjectFunctionRollup5() {
	c.Test(
		"var o; o = {};\n"+
			"o.x = function() {return o};\n"+
			"o.y = [function() {return o}];\n"+
			"o.z = {a:function() {return o}};\n",
		"var o; o = {\n"+
			"x:function () {return o},\n"+
			"y:[function () {return o}],\n"+
			"z:{a:function () {return o}}};\n")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testObjectPropertyReassigned()
func (c *PeepholeCollectPropertyAssignmentsTest) TestObjectPropertyReassigned() {
	c.Test(
		"var a = {b:''};\n"+
			"a.b='c';\n",
		"var a={b:'c'};")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testObjectPropertyReassigned2()
func (c *PeepholeCollectPropertyAssignmentsTest) TestObjectPropertyReassigned2() {
	c.Test(
		"var a = {b:'', x:10};\n"+
			"a.b='c';\n",
		"var a={x:10, b:'c'};")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testObjectPropertyReassigned3()
func (c *PeepholeCollectPropertyAssignmentsTest) TestObjectPropertyReassigned3() {
	c.Test(
		"var a = {x:10};\n"+
			"a.b = 'c';\n",
		"var a = {x:10, b:'c'};")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testObjectPropertyReassigned4()
func (c *PeepholeCollectPropertyAssignmentsTest) TestObjectPropertyReassigned4() {
	c.TestSame(
		"var a = {b:10};\n" +
			"var x = 1;\n" +
			"a.b = x+10;\n")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testObject_usingSpread_diablesOptimization_leadingSpread()
func (c *PeepholeCollectPropertyAssignmentsTest) TestObject_usingSpread_diablesOptimization_leadingSpread() {
	c.Test(
		"var a = {a: 1, ...b, c: 2}; a.a = 3;", //
		"var a = {...b, c: 2, a: 3};")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testObject_usingSpread_diablesOptimization_withinSpread()
func (c *PeepholeCollectPropertyAssignmentsTest) TestObject_usingSpread_diablesOptimization_withinSpread() {
	c.TestSame("var a = {a: 1, ...b, c: 2, b: 3};")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testObject_usingSpread_diablesOptimization_followingSpread()
func (c *PeepholeCollectPropertyAssignmentsTest) TestObject_usingSpread_diablesOptimization_followingSpread() {
	c.TestSame("var a = {a: 1, ...b, c: 3};")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testObjectComputedProp1()
func (c *PeepholeCollectPropertyAssignmentsTest) TestObjectComputedProp1() {
	c.TestSame(
		"var a = {['computed']: 10};\n" +
			"var alsoComputed = 'someValue';\n" +
			"a[alsoComputed] = 20;\n")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testObjectComputedProp2()
func (c *PeepholeCollectPropertyAssignmentsTest) TestObjectComputedProp2() {
	c.Test(
		"var a = {['computed']: 10};\n"+
			"a.prop = 20;\n",
		"var a = {\n"+
			"  ['computed']: 10,\n"+
			"  prop: 20,\n"+
			"};\n")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testObjectMemberFunction1()
func (c *PeepholeCollectPropertyAssignmentsTest) TestObjectMemberFunction1() {
	c.Test(
		"var a = { member() {} };\n"+
			"a.prop = 20;\n",
		"var a = {\n"+
			"  member() {},\n"+
			"  prop: 20,\n"+
			"};\n")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testObjectMemberFunction2()
func (c *PeepholeCollectPropertyAssignmentsTest) TestObjectMemberFunction2() {
	c.Test(
		"var a = { member() {} };\n"+
			"a.member = 20;\n",
		"var a = {\n"+
			"  member: 20,\n"+
			"};\n")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testObjectGetter()
func (c *PeepholeCollectPropertyAssignmentsTest) TestObjectGetter() {
	c.TestSame(
		"var a = { get x() {} };\n" +
			"a.x = 20;\n")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testObjectSetter()
func (c *PeepholeCollectPropertyAssignmentsTest) TestObjectSetter() {
	c.TestSame(
		"var a = { set x(value) {} };\n" +
			"a.x = 20;\n")
}

// Java: PeepholeCollectPropertyAssignmentsTest.testArraySpread()
func (c *PeepholeCollectPropertyAssignmentsTest) TestArraySpread() {
	c.TestSame("var a = [1, ...spread]; a[2] = 3;")
	c.TestSame("var a = [...spread]; a[0] = 1;")
	c.TestSame("var a = [1, ...spread]; a[0] = 2;")
}
