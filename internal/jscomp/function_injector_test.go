// Port of com/google/javascript/jscomp/FunctionInjectorTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/FunctionInjectorTest.java)

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
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// FunctionInjectorTest: inline function tests.
type FunctionInjectorTest struct {
	// Go-only: test-helper — the Go test running the current test method.
	t *testing.T

	assumeStrictThis   bool
	allowDecomposition bool
}

// Java: FunctionInjectorTest.INLINE_DIRECT
var FunctionInjectorTest_INLINE_DIRECT = FunctionInjector_InliningMode_DIRECT

// Java: FunctionInjectorTest.INLINE_BLOCK
var FunctionInjectorTest_INLINE_BLOCK = FunctionInjector_InliningMode_BLOCK

// Java: FunctionInjectorTest.ASSUME_MINIMUM_CAPTURE
const FunctionInjectorTest_ASSUME_MINIMUM_CAPTURE = false

// TODO(johnlenz): remove this constant once this has been proven in
// production code.
//
// Java: FunctionInjectorTest.NEW_VARS_IN_GLOBAL_SCOPE
var FunctionInjectorTest_NEW_VARS_IN_GLOBAL_SCOPE = FunctionInjector_CanInlineResult_YES

// Java: FunctionInjectorTest()
func NewFunctionInjectorTest() *FunctionInjectorTest {
	return &FunctionInjectorTest{assumeStrictThis: false}
}

func TestFunctionInjectorTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewFunctionInjectorTest() })
}

// JUnitSetT gives the instance the Go test running its test method.
func (c *FunctionInjectorTest) JUnitSetT(t *testing.T) { c.t = t }

// T is the Go test running the current test method.
//
// Go-only: test-helper — used by the @Ignore'd methods to skip.
func (c *FunctionInjectorTest) T() *testing.T { return c.t }

// JUnitBefore runs the @Before method setUp().
func (c *FunctionInjectorTest) JUnitBefore() { c.SetUp() }

// JUnitAfter: the class has no @After method.
func (c *FunctionInjectorTest) JUnitAfter() {}

// Java: FunctionInjectorTest.setUp()
func (c *FunctionInjectorTest) SetUp() {
	c.assumeStrictThis = false
	c.allowDecomposition = false
}

// Java: FunctionInjectorTest.testIsSimpleFunction1()
func (c *FunctionInjectorTest) TestIsSimpleFunction1() {
	AssertThat(FunctionInjector_IsDirectCallNodeReplacementPossible(FunctionInjectorTest_Prep("function f(){}"))).IsTrue()
}

// Java: FunctionInjectorTest.testIsSimpleFunction2()
func (c *FunctionInjectorTest) TestIsSimpleFunction2() {
	AssertThat(FunctionInjector_IsDirectCallNodeReplacementPossible(FunctionInjectorTest_Prep("function f(){return 0;}"))).IsTrue()
}

// Java: FunctionInjectorTest.testIsSimpleFunction3()
func (c *FunctionInjectorTest) TestIsSimpleFunction3() {
	AssertThat(FunctionInjector_IsDirectCallNodeReplacementPossible(FunctionInjectorTest_Prep("function f(){return x ? 0 : 1}"))).
		IsTrue()
}

// Java: FunctionInjectorTest.testIsSimpleFunction4()
func (c *FunctionInjectorTest) TestIsSimpleFunction4() {
	AssertThat(FunctionInjector_IsDirectCallNodeReplacementPossible(FunctionInjectorTest_Prep("function f(){return;}"))).IsFalse()
}

// Java: FunctionInjectorTest.testIsSimpleFunction5()
func (c *FunctionInjectorTest) TestIsSimpleFunction5() {
	AssertThat(FunctionInjector_IsDirectCallNodeReplacementPossible(FunctionInjectorTest_Prep("function f(){return 0; return 0;}"))).
		IsFalse()
}

// Java: FunctionInjectorTest.testIsSimpleFunction6()
func (c *FunctionInjectorTest) TestIsSimpleFunction6() {
	AssertThat(
		FunctionInjector_IsDirectCallNodeReplacementPossible(FunctionInjectorTest_Prep("function f(){var x=true;return x ? 0 : 1}"))).
		IsFalse()
}

// Java: FunctionInjectorTest.testIsSimpleFunction7()
func (c *FunctionInjectorTest) TestIsSimpleFunction7() {
	AssertThat(
		FunctionInjector_IsDirectCallNodeReplacementPossible(
			FunctionInjectorTest_Prep("function f(){if (x) return 0; else return 1}"))).
		IsFalse()
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction1()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction1() {
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_YES, "function foo(){}; foo();", "foo", FunctionInjectorTest_INLINE_DIRECT)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction2()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction2() {
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_YES, "function foo(){}; foo();", "foo", FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction3()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction3() {
	// NOTE: FoldConstants will convert this to a empty function,
	// so there is no need to explicitly support it.
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_NO, "function foo(){return;}; foo();", "foo", FunctionInjectorTest_INLINE_DIRECT)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction4()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction4() {
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_YES, "function foo(){return;}; foo();", "foo", FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction5()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction5() {
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_YES, "function foo(){return true;}; foo();", "foo", FunctionInjectorTest_INLINE_DIRECT)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction6()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction6() {
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_YES, "function foo(){return true;}; foo();", "foo", FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction7()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction7() {
	// In var initialization.
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_YES, "function foo(){return true;}; var x=foo();", "foo", FunctionInjectorTest_INLINE_DIRECT)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction8()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction8() {
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_YES, "function foo(){return true;}; var x=foo();", "foo", FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction9()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction9() {
	// In assignment.
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_YES, "function foo(){return true;}; var x; x=foo();", "foo", FunctionInjectorTest_INLINE_DIRECT)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction10()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction10() {
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_YES, "function foo(){return true;}; var x; x=foo();", "foo", FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction11()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction11() {
	// In expression.
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_YES,
		"function foo(){return true;}; var x; x=x+foo();",
		"foo",
		FunctionInjectorTest_INLINE_DIRECT)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction12()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction12() {
	// "foo" is not known to be side-effect free, it might change the value
	// of "x", so it can't be inlined.
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_NO, "function foo(){return true;}; var x; x=x+foo();", "foo", FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction12b()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction12b() {
	// "foo" is not known to be side-effect free, it might change the value
	// of "x", so it can't be inlined.
	c.allowDecomposition = true
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_AFTER_PREPARATION,
		"function foo(){return true;}; var x; x=x+foo();",
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// TODO(nicksantos): Re-enable with side-effect detection.
// Java: FunctionInjectorTest.testCanInlineReferenceToFunction13()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction13() {
	c.T().Skip("@Ignore") // @Ignore
	// ... if foo is side-effect free we can inline here.
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_YES,
		`/** @nosideeffects */ function foo(){return true;};
var x; x=x+foo();
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction14()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction14() {
	// Simple call with parameters
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_YES, "function foo(a){return true;}; foo(x);", "foo", FunctionInjectorTest_INLINE_DIRECT)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction15()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction15() {
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_YES, "function foo(a){return true;}; foo(x);", "foo", FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction16()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction16() {
	// Function "foo" as it contains "var b" which
	// must be brought into the global scope.
	c.HelperCanInlineReferenceToFunction(
		FunctionInjectorTest_NEW_VARS_IN_GLOBAL_SCOPE,
		"function foo(a){var b;return a;}; foo(goo());",
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction17()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction17() {
	// This doesn't bring names into the global name space.
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_YES,
		`function foo(a){return a;};
function x() { foo(goo()); }
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction18()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction18() {
	// Parameter has side-effects.
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_YES, "function foo(a){return a;} foo(x++);", "foo", FunctionInjectorTest_INLINE_DIRECT)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction18b()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction18b() {
	// Parameter has side-effects.
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_NO,
		"function foo(a,b){return b,a;} foo(x++,use(x));",
		"foo",
		FunctionInjectorTest_INLINE_DIRECT)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction19()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction19() {
	// Parameter has mutable parameter referenced more than once.
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_NO, "function foo(a){return a+a} foo([]);", "foo", FunctionInjectorTest_INLINE_DIRECT)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction20()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction20() {
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_NO, "function foo(a){return a+a} foo({});", "foo", FunctionInjectorTest_INLINE_DIRECT)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction21()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction21() {
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_NO, "function foo(a){return a+a} foo(new Date);", "foo", FunctionInjectorTest_INLINE_DIRECT)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction22()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction22() {
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_NO,
		"function foo(a){return a+a} foo(true && new Date);",
		"foo",
		FunctionInjectorTest_INLINE_DIRECT)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction23()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction23() {
	// variables to global scope.
	c.HelperCanInlineReferenceToFunction(
		FunctionInjectorTest_NEW_VARS_IN_GLOBAL_SCOPE, "function foo(a){return a;}; foo(x++);", "foo", FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction24()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction24() {
	// ... this is OK, because it doesn't introduce a new global name.
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_YES,
		`function foo(a){return a;};
function x() { foo(x++); }
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction25()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction25() {
	// Parameter has side-effects.
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_NO, "function foo(a){return a+a;}; foo(x++);", "foo", FunctionInjectorTest_INLINE_DIRECT)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction26()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction26() {
	c.HelperCanInlineReferenceToFunction(
		FunctionInjectorTest_NEW_VARS_IN_GLOBAL_SCOPE, "function foo(a){return a+a;}; foo(x++);", "foo", FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction27()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction27() {
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_YES,
		`function foo(a){return a+a;};
function x() { foo(x++); }
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction28()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction28() {
	// Parameter has side-effects.
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_NO, "function foo(a){return true;}; foo(goo());", "foo", FunctionInjectorTest_INLINE_DIRECT)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction29()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction29() {
	c.HelperCanInlineReferenceToFunction(
		FunctionInjectorTest_NEW_VARS_IN_GLOBAL_SCOPE,
		"function foo(a){return true;}; foo(goo());",
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction30()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction30() {
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_YES,
		`function foo(a){return true;};
function x() { foo(goo()); }
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction31()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction31() {
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_YES,
		`function foo(a) {return true;};
function x() {foo.call(this, 1);}
`,
		"foo",
		FunctionInjectorTest_INLINE_DIRECT)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction32()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction32() {
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_NO,
		`function foo(a){return true;};
function x() { foo.apply(this, [1]); }
`,
		"foo",
		FunctionInjectorTest_INLINE_DIRECT)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction33()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction33() {
	// No special handling is required for method calls passing this.
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_YES,
		`function foo(a){return true;};
function x() { foo.bar(this, 1); }
`,
		"foo",
		FunctionInjectorTest_INLINE_DIRECT)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction34()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction34() {
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_YES,
		`function foo(a){return true;};
function x() { foo.call(this, goo()); }
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction35()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction35() {
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_NO,
		`function foo(a){return true;};
function x() { foo.apply(this, goo()); }
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction36()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction36() {
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_YES,
		`function foo(a){return true;};
function x() { foo.bar(this, goo()); }
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction37()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction37() {
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_NO,
		`function foo(a){return true;};
function x() { foo.call(null, 1); }
`,
		"foo",
		FunctionInjectorTest_INLINE_DIRECT)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction38()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction38() {
	c.assumeStrictThis = false
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_NO,
		`function foo(a){return true;};
function x() { foo.call(null, goo()); }
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
	c.assumeStrictThis = true
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_YES,
		`function foo(a){return true;};
function x() { foo.call(null, goo()); }
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction39()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction39() {
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_NO,
		`function foo(a){return true;};
function x() { foo.call(bar, 1); }
`,
		"foo",
		FunctionInjectorTest_INLINE_DIRECT)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction40()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction40() {
	c.assumeStrictThis = false
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_NO,
		`function foo(a){return true;};
function x() { foo.call(bar, goo()); }
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
	c.assumeStrictThis = true
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_YES,
		`function foo(a){return true;};
function x() { foo.call(bar, goo()); }
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction41()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction41() {
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_NO,
		`function foo(a){return true;};
function x() { foo.call(new bar(), 1); }
`,
		"foo",
		FunctionInjectorTest_INLINE_DIRECT)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction42()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction42() {
	c.assumeStrictThis = false
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_NO,
		`function foo(a){return true;};
function x() { foo.call(new bar(), goo()); }
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
	c.assumeStrictThis = true
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_YES,
		`function foo(a){return true;};
function x() { foo.call(new bar(), goo()); }
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction43()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction43() {
	// Handle the case of a missing 'this' value in a call.
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_NO,
		`function foo(){return true;};
function x() { foo.call(); }
`,
		"foo",
		FunctionInjectorTest_INLINE_DIRECT)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction44()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction44() {
	c.assumeStrictThis = false
	// Handle the case of a missing 'this' value in a call.
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_NO,
		`function foo(){return true;};
function x() { foo.call(); }
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
	c.assumeStrictThis = true
	// Handle the case of a missing 'this' value in a call.
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_YES,
		`function foo(){return true;};
function x() { foo.call(); }
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction45()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction45() {
	// Call with inner function expression.
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_YES,
		"function foo(){return function() {return true;}}; foo();",
		"foo",
		FunctionInjectorTest_INLINE_DIRECT)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction46()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction46() {
	// Call with inner function expression.
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_YES,
		"function foo(){return function() {return true;}}; foo();",
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction47()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction47() {
	// Call with inner function expression and variable decl.
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_NO,
		"function foo(){var a; return function() {return true;}}; foo();",
		"foo",
		FunctionInjectorTest_INLINE_DIRECT)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction48()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction48() {
	// Call with inner function expression and variable decl.
	// TODO(johnlenz): should we validate no values in scope?
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_YES,
		"function foo(){var a; return function() {return true;}}; foo();",
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction49()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction49() {
	// Call with inner function expression.
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_YES,
		"function foo(){return function() {var a; return true;}}; foo();",
		"foo",
		FunctionInjectorTest_INLINE_DIRECT)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction50()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction50() {
	// Call with inner function expression.
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_YES,
		"function foo(){return function() {var a; return true;}}; foo();",
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction51()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction51() {
	// Call with inner function statement.
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_YES,
		"function foo(){function x() {var a; return true;} return x}; foo();",
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunction52()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunction52() {
	// Don't inline functions with var declarations into a scope with inner functions
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_NO,
		`function foo() { var a = 3; return a; }
function bar() { function baz() {} if (true) { foo(); } }
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunctionInExpression1()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunctionInExpression1() {
	// Call in if condition
	c.allowDecomposition = true
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_AFTER_PREPARATION,
		"function foo(a){return true;}; function x() { if (foo(1)) throw 'test'; }",
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunctionInExpression2()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunctionInExpression2() {
	// Call in return expression
	c.allowDecomposition = true
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_AFTER_PREPARATION,
		"function foo(a){return true;}; function x() { return foo(1); }",
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunctionInExpression3()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunctionInExpression3() {
	// Call in switch expression
	c.allowDecomposition = true
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_AFTER_PREPARATION,
		"function foo(a){return true;}; function x() { switch(foo(1)) { default:break; } }",
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunctionInExpression4()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunctionInExpression4() {
	// Call in hook condition
	c.allowDecomposition = true
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_AFTER_PREPARATION,
		"function foo(a){return true;}; function x() {foo(1)?0:1 }",
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunctionInExpression5()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunctionInExpression5() {
	// Call in hook side-effect free condition
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_NO,
		`function foo(a){return true;};
function x() {true?foo(1):1 }
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunctionInExpression5a()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunctionInExpression5a() {
	// Call in hook side-effect free condition
	c.allowDecomposition = true
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_AFTER_PREPARATION,
		"function foo(a){return true;}; function x() {true?foo(1):1 }",
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunctionInExpression6()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunctionInExpression6() {
	// Call in expression statement "condition"
	c.allowDecomposition = true
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_AFTER_PREPARATION,
		"function foo(a){return true;}; function x() {foo(1) && 1 }",
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunctionInExpression7()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunctionInExpression7() {
	// Call in expression statement after side-effect free "condition"
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_NO,
		`function foo(a){return true;};
function x() {1 && foo(1) }
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunctionInExpression7a()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunctionInExpression7a() {
	// Call in expression statement after side-effect free "condition"
	c.allowDecomposition = true
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_AFTER_PREPARATION,
		"function foo(a){return true;}; function x() {1 && foo(1) }",
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunctionInExpression8()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunctionInExpression8() {
	// Call in expression statement after side-effect free operator
	c.allowDecomposition = true
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_AFTER_PREPARATION,
		"function foo(a){return true;}; function x() {1 + foo(1) }",
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunctionInExpression9()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunctionInExpression9() {
	// Call in VAR expression.
	c.allowDecomposition = true
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_AFTER_PREPARATION,
		"function foo(a){return true;}; function x() {var b = 1 + foo(1)}",
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunctionInExpression10()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunctionInExpression10() {
	// Call in assignment expression.
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_NO,
		`function foo(a){return true;};
function x() {var b; b += 1 + foo(1) }
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunctionInExpression10a()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunctionInExpression10a() {
	// Call in assignment expression.
	c.allowDecomposition = true
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_AFTER_PREPARATION,
		"function foo(a){return true;}; function x() {var b; b += 1 + foo(1) }",
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// TODO(nicksantos): Re-enable with side-effect detection.
// Java: FunctionInjectorTest.testCanInlineReferenceToFunctionInExpression11()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunctionInExpression11() {
	c.T().Skip("@Ignore") // @Ignore
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_YES,
		`/** @nosideeffects */ function foo(a){return true;};
function x() {var b; b += 1 + foo(1) }
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunctionInExpression12()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunctionInExpression12() {
	c.allowDecomposition = true
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_AFTER_PREPARATION,
		"function foo(a){return true;}; function x() {var a,b,c; a = b = c = foo(1) }",
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunctionInExpression13()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunctionInExpression13() {
	c.allowDecomposition = true
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_AFTER_PREPARATION,
		"function foo(a){return true;}; function x() {var a,b,c; a = b = c = 1 + foo(1) }",
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunctionInExpression14()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunctionInExpression14() {
	// ... foo can not be inlined because of possible changes to "c".
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_NO,
		`var a = {}, b = {}, c;
a.test = 'a';
b.test = 'b';
c = a;
function foo(){c = b; return 'foo'};
c.test=foo();
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunctionInExpression14a()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunctionInExpression14a() {
	// ... foo can be inlined despite possible changes to "c".
	c.allowDecomposition = true
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_AFTER_PREPARATION,
		`var a = {}, b = {}, c;
a.test = 'a';
b.test = 'b';
c = a;
function foo(){c = b; return 'foo'};
c.test=foo();
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// TODO(nicksantos): Re-enable with side-effect detection.
// Java: FunctionInjectorTest.testCanInlineReferenceToFunctionInExpression15()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunctionInExpression15() {
	c.T().Skip("@Ignore") // @Ignore
	// ... foo can be inlined as it is side-effect free.
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_YES,
		`var a = {}, b = {}, c;
a.test = 'a';
b.test = 'b';
c = a;
/** @nosideeffects */ function foo(){return 'foo'};
c.test=foo();
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunctionInExpression16()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunctionInExpression16() {
	c.T().Skip("@Ignore") // @Ignore
	// ... foo can not be inlined because of possible side-effects of x()
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_NO,
		`var a = {}, b = {}, c;
a.test = 'a';
b.test = 'b';
c = a;
function x(){return c};
/** @nosideeffects */ function foo(){return 'foo'};
x().test=foo();
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunctionInExpression17()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunctionInExpression17() {
	c.T().Skip("@Ignore") // @Ignore
	// ... foo can be inlined because of x() is side-effect free.
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_YES,
		`var a = {}, b = {}, c;
a.test = 'a';
b.test = 'b';
c = a;
/** @nosideeffects */ function x(){return c};
/** @nosideeffects */ function foo(){return 'foo'};
x().test=foo();
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunctionInExpression18()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunctionInExpression18() {
	// Call in within a call
	c.allowDecomposition = true
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_AFTER_PREPARATION,
		"function foo(){return _g();}; function x() {1 + foo()() }",
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunctionInExpression19()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunctionInExpression19() {
	// ... unless foo is known to be side-effect free, it might actually
	// change the value of "_g" which would unfortunately change the behavior,
	// so we can't inline here.
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_NO,
		`function foo(){return a;};
function x() {1 + _g(foo()) }
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunctionInExpression19a()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunctionInExpression19a() {
	// ... unless foo is known to be side-effect free, it might actually
	// change the value of "_g" which would unfortunately change the behavior,
	// so we can't inline here.
	c.allowDecomposition = true
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_AFTER_PREPARATION,
		"function foo(){return a;}; function x() {1 + _g(foo()) }",
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// TODO(nicksantos): Re-enable with side-effect detection.
// Java: FunctionInjectorTest.testCanInlineReferenceToFunctionInExpression20()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunctionInExpression20() {
	c.T().Skip("@Ignore") // @Ignore
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_YES,
		`/** @nosideeffects */ function foo(){return a;};
function x() {1 + _g(foo()) }
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunctionInExpression21()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunctionInExpression21() {
	// Assignments to object are problematic if the call has side-effects,
	// as the object that is being referred to can change.
	// Note: This could be changed be inlined if we in some way make "z"
	// as not escaping from the local scope.
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_NO,
		`var z = {};
function foo(a){z = {};return true;};
function x() { z.gack = foo(1) }
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunctionInExpression21a()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunctionInExpression21a() {
	// Assignments to object are problematic if the call has side-effects,
	// as the object that is being referred to can change.
	// Note: This could be changed be inlined if we in some way make "z"
	// as not escaping from the local scope.
	c.allowDecomposition = true
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_AFTER_PREPARATION,
		`var z = {};
function foo(a){z = {};return true;};
function x() { z.gack = foo(1) }
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunctionInExpression22()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunctionInExpression22() {
	// ... foo() is after a side-effect
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_NO,
		`function foo(){return a;};
function x() {1 + _g(_a(), foo()) }
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunctionInExpression22a()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunctionInExpression22a() {
	// ... foo() is after a side-effect
	c.allowDecomposition = true
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_AFTER_PREPARATION,
		"function foo(){return a;}; function x() {1 + _g(_a(), foo()) }",
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunctionInExpression23()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunctionInExpression23() {
	// ... foo() is after a side-effect
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_NO,
		`function foo(){return a;};
function x() {1 + _g(_a(), foo.call(this)) }
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunctionInExpression23a()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunctionInExpression23a() {
	// ... foo() is after a side-effect
	c.allowDecomposition = true
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_AFTER_PREPARATION,
		"function foo(){return a;}; function x() {1 + _g(_a(), foo.call(this)) }",
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunctionInLoop1()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunctionInLoop1() {
	c.allowDecomposition = true
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_YES, "function foo(){return a;}; while(1) { foo(); }", "foo", FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReferenceToFunctionInLoop2()
func (c *FunctionInjectorTest) TestCanInlineReferenceToFunctionInLoop2() {
	// If function contains function, don't inline it into a loop.
	// TODO(johnlenz): this can be improved by looking to see
	// if the inner function contains any references to values defined
	// in the outer function.
	c.allowDecomposition = true
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_NO,
		"function foo(){return function() {};}; while(1) { foo(); }",
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineFunctionWithInnerArrowFunction1()
func (c *FunctionInjectorTest) TestCanInlineFunctionWithInnerArrowFunction1() {
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_YES, "function foo(){ () => { alert(1); }; } foo();", "foo", FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineFunctionWithInnerArrowFunction2()
func (c *FunctionInjectorTest) TestCanInlineFunctionWithInnerArrowFunction2() {
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_NO, "function foo(){ () => { this; }; } foo();", "foo", FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInline1()
func (c *FunctionInjectorTest) TestInline1() {
	c.HelperInlineReferenceToFunction(
		"function foo(){}; foo();", "function foo(){}; void 0", "foo", FunctionInjectorTest_INLINE_DIRECT)
}

// Java: FunctionInjectorTest.testInline2()
func (c *FunctionInjectorTest) TestInline2() {
	c.HelperInlineReferenceToFunction(
		"function foo(){}; foo();", "function foo(){}; {}", "foo", FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInline3()
func (c *FunctionInjectorTest) TestInline3() {
	c.HelperInlineReferenceToFunction(
		"function foo(){return;}; foo();", "function foo(){return;}; {}", "foo", FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInline4()
func (c *FunctionInjectorTest) TestInline4() {
	c.HelperInlineReferenceToFunction(
		"function foo(){return true;}; foo();",
		"function foo(){return true;}; true;",
		"foo",
		FunctionInjectorTest_INLINE_DIRECT)
}

// Java: FunctionInjectorTest.testInline5()
func (c *FunctionInjectorTest) TestInline5() {
	c.HelperInlineReferenceToFunction(
		"function foo(){return true;}; foo();",
		"function foo(){return true;}; {true;}",
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInline6()
func (c *FunctionInjectorTest) TestInline6() {
	// In var initialization.
	c.HelperInlineReferenceToFunction(
		"function foo(){return true;}; var x=foo();",
		"function foo(){return true;}; var x=true;",
		"foo",
		FunctionInjectorTest_INLINE_DIRECT)
}

// Java: FunctionInjectorTest.testInline7()
func (c *FunctionInjectorTest) TestInline7() {
	c.HelperInlineReferenceToFunction(
		"function foo(){return true;}; var x=foo();",
		`function foo(){return true;}; var x;
{x=true}
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInline8()
func (c *FunctionInjectorTest) TestInline8() {
	// In assignment.
	c.HelperInlineReferenceToFunction(
		"function foo(){return true;}; var x; x=foo();",
		"function foo(){return true;}; var x; x=true;",
		"foo",
		FunctionInjectorTest_INLINE_DIRECT)
}

// Java: FunctionInjectorTest.testInline9()
func (c *FunctionInjectorTest) TestInline9() {
	c.HelperInlineReferenceToFunction(
		"function foo(){return true;}; var x; x=foo();",
		"function foo(){return true;}; var x;{x=true}",
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInline10()
func (c *FunctionInjectorTest) TestInline10() {
	// In expression.
	c.HelperInlineReferenceToFunction(
		"function foo(){return true;}; var x; x=x+foo();",
		"function foo(){return true;}; var x; x=x+true;",
		"foo",
		FunctionInjectorTest_INLINE_DIRECT)
}

// Java: FunctionInjectorTest.testInline11()
func (c *FunctionInjectorTest) TestInline11() {
	// Simple call with parameters
	c.HelperInlineReferenceToFunction(
		"function foo(a){return true;}; foo(x);",
		"function foo(a){return true;}; true;",
		"foo",
		FunctionInjectorTest_INLINE_DIRECT)
}

// Java: FunctionInjectorTest.testInline12()
func (c *FunctionInjectorTest) TestInline12() {
	c.HelperInlineReferenceToFunction(
		"function foo(a){return true;}; foo(x);",
		"function foo(a){return true;}; {true}",
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInline13()
func (c *FunctionInjectorTest) TestInline13() {
	// Parameter has side-effects.
	c.HelperInlineReferenceToFunction(
		`function foo(a){return a;};
function x() { foo(x++); }
`,
		`function foo(a){return a;};
function x(){{x++}}
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInline14()
func (c *FunctionInjectorTest) TestInline14() {
	// Parameter has side-effects.
	c.HelperInlineReferenceToFunction(
		"function foo(a){return a+a;}; foo(x++);",
		`function foo(a){return a+a;};
{var a$jscomp$inline_0=x++;
 a$jscomp$inline_0+
a$jscomp$inline_0;}
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInline15()
func (c *FunctionInjectorTest) TestInline15() {
	// Parameter has mutable, references more than once.
	c.HelperInlineReferenceToFunction(
		"function foo(a){return a+a;}; foo(new Date());",
		`function foo(a){return a+a;};
{var a$jscomp$inline_0=new Date();
 a$jscomp$inline_0+
a$jscomp$inline_0;}
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInline16()
func (c *FunctionInjectorTest) TestInline16() {
	// Parameter is large, references more than once.
	c.HelperInlineReferenceToFunction(
		"function foo(a){return a+a;}; foo(function(){});",
		`function foo(a){return a+a;};
{var a$jscomp$inline_0=function(){};
 a$jscomp$inline_0+
a$jscomp$inline_0;}
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInline17()
func (c *FunctionInjectorTest) TestInline17() {
	// Parameter has side-effects.
	c.HelperInlineReferenceToFunction(
		"function foo(a){return true;}; foo(goo());",
		`function foo(a){return true;};
{var a$jscomp$inline_0=goo();true}
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInline18()
func (c *FunctionInjectorTest) TestInline18() {
	// This doesn't bring names into the global name space.
	c.HelperInlineReferenceToFunction(
		"function foo(a){var b;return a;} function x() { foo(goo()); }",
		`function foo(a) {
  var b;
  return a;
}
function x() {
  {
    var a$jscomp$inline_0 = goo();
    var b$jscomp$inline_1;
    a$jscomp$inline_0;
  }
}
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInline19()
func (c *FunctionInjectorTest) TestInline19() {
	// Properly alias.
	c.HelperInlineReferenceToFunction(
		`var x = 1; var y = 2;
function foo(a,b){x = b; y = a;};
function bar() { foo(x,y); }
`,
		`var x = 1; var y = 2;
function foo(a,b){x = b; y = a;};
function bar() {
{var a$jscomp$inline_0=x;
x = y;
y = a$jscomp$inline_0;}
}
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInline19b()
func (c *FunctionInjectorTest) TestInline19b() {
	c.HelperInlineReferenceToFunction(
		`var x = 1; var y = 2;
function foo(a, b) {
  y = a; x = b;
};
function bar() { foo(x,y); }
`,
		`var x = 1; var y = 2;
function foo(a,b){
  y = a; x = b;
};
function bar() {
{
  var a$jscomp$inline_0=x;
  var b$jscomp$inline_1=y;
  y=a$jscomp$inline_0;
  x=b$jscomp$inline_1}
}
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInline20()
func (c *FunctionInjectorTest) TestInline20() {
	// cloned FUNCTION node must be reported as being added.
	c.HelperInlineReferenceToFunction(
		"function foo(a){return a}; foo(function(){});",
		"function foo(a){return a}; (function(){})",
		"foo",
		FunctionInjectorTest_INLINE_DIRECT)
}

// Java: FunctionInjectorTest.testInline21()
func (c *FunctionInjectorTest) TestInline21() {
	// cloned FUNCTION node must be reported as being added.
	c.HelperInlineReferenceToFunction(
		"function foo(a){return a}; foo(function(){});",
		"function foo(a){return a}; {(function(){})}",
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInlineIntoLoop()
func (c *FunctionInjectorTest) TestInlineIntoLoop() {
	c.HelperInlineReferenceToFunction(
		"function foo(a){var b;return a;}; for(;1;){ foo(1); }",
		"function foo(a){var b;return a;}; for(;1;){ { var b$jscomp$inline_1=void 0;1} }",
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
	c.HelperInlineReferenceToFunction(
		`function foo(a){var b;return a;};
do{ foo(1); } while(1)
`,
		`function foo(a){var b;return a;};
do{ {
var b$jscomp$inline_1=void 0;1}}while(1)
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
	c.HelperInlineReferenceToFunction(
		`function foo(a){for(var b in c)return a;};
for(;1;){ foo(1); }
`,
		`function foo(a){var b;for(b in c)return a;};
for(;1;){ {JSCompiler_inline_label_foo_2:{
var b$jscomp$inline_1=void 0;for(b$jscomp$inline_1 in c){
1;break JSCompiler_inline_label_foo_2
}}}}
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInlineFunctionWithInnerFunction1()
func (c *FunctionInjectorTest) TestInlineFunctionWithInnerFunction1() {
	// Call with inner function expression.
	c.HelperInlineReferenceToFunction(
		"function foo(){return function() {return true;}}; foo();",
		`function foo(){return function() {return true;}};
(function() {return true;})
`,
		"foo",
		FunctionInjectorTest_INLINE_DIRECT)
}

// Java: FunctionInjectorTest.testInlineFunctionWithInnerFunction2()
func (c *FunctionInjectorTest) TestInlineFunctionWithInnerFunction2() {
	// Call with inner function expression.
	c.HelperInlineReferenceToFunction(
		"function foo(){return function() {return true;}}; foo();",
		`function foo(){return function() {return true;}};
{(function() {return true;})}
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInlineFunctionWithInnerFunction3()
func (c *FunctionInjectorTest) TestInlineFunctionWithInnerFunction3() {
	// Call with inner function expression.
	c.HelperInlineReferenceToFunction(
		"function foo(){return function() {var a; return true;}}; foo();",
		`function foo(){return function() {var a; return true;}};
(function() {var a; return true;});
`,
		"foo",
		FunctionInjectorTest_INLINE_DIRECT)
}

// Java: FunctionInjectorTest.testInlineFunctionWithInnerFunction4()
func (c *FunctionInjectorTest) TestInlineFunctionWithInnerFunction4() {
	// Call with inner function expression.
	c.HelperInlineReferenceToFunction(
		"function foo(){return function() {var a; return true;}}; foo();",
		`function foo(){return function() {var a; return true;}};
{(function() {var a$jscomp$inline_0; return true;});}
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInlineFunctionWithInnerFunction5()
func (c *FunctionInjectorTest) TestInlineFunctionWithInnerFunction5() {
	// Call with inner function statement.
	c.HelperInlineReferenceToFunction(
		"function foo(){function x() {var a; return true;} return x} foo();",
		`function foo() {
  function x() { var a; return true; }
  return x;
}
{
  var x$jscomp$inline_0 = function(){
    var a$jscomp$inline_1;
    return true;
  };
  x$jscomp$inline_0;
}
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInlineFunctionWithInnerArrowFunction1()
func (c *FunctionInjectorTest) TestInlineFunctionWithInnerArrowFunction1() {
	c.HelperInlineReferenceToFunction(
		"function foo(){ () => { alert(1); }; } foo();",
		`function foo(){ () => { alert(1); }; }
{ () => { alert(1); }; }
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInlineReferenceInExpression1()
func (c *FunctionInjectorTest) TestInlineReferenceInExpression1() {
	// Call in if condition
	c.allowDecomposition = true
	c.HelperInlineReferenceToFunction(
		"function foo(a){return true;}; function x() { if (foo(1)) throw 'test'; }",
		`function foo(a){return true;};
function x() { var JSCompiler_inline_result$jscomp$0;
{JSCompiler_inline_result$jscomp$0=true;}
if (JSCompiler_inline_result$jscomp$0) throw 'test'; }
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInlineReferenceInExpression2()
func (c *FunctionInjectorTest) TestInlineReferenceInExpression2() {
	// Call in return expression
	c.allowDecomposition = true
	c.HelperInlineReferenceToFunction(
		"function foo(a){return true;}; function x() { return foo(1); }",
		`function foo(a){return true;};
function x() { var JSCompiler_inline_result$jscomp$0;
{JSCompiler_inline_result$jscomp$0=true;}
return JSCompiler_inline_result$jscomp$0; }
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInlineReferenceInExpression3()
func (c *FunctionInjectorTest) TestInlineReferenceInExpression3() {
	// Call in switch expression
	c.allowDecomposition = true
	c.HelperInlineReferenceToFunction(
		"function foo(a){return true;}; function x() { switch(foo(1)) { default:break; } }",
		`function foo(a){return true;};
function x() { var JSCompiler_inline_result$jscomp$0;
{JSCompiler_inline_result$jscomp$0=true;}
switch(JSCompiler_inline_result$jscomp$0) { default:break; } }
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInlineReferenceInExpression4()
func (c *FunctionInjectorTest) TestInlineReferenceInExpression4() {
	// Call in hook condition
	c.allowDecomposition = true
	c.HelperInlineReferenceToFunction(
		"function foo(a){return true;}; function x() {foo(1)?0:1 }",
		`function foo(a){return true;};
function x() { var JSCompiler_inline_result$jscomp$0;
{JSCompiler_inline_result$jscomp$0=true;}
JSCompiler_inline_result$jscomp$0?0:1 }
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInlineReferenceInExpression5()
func (c *FunctionInjectorTest) TestInlineReferenceInExpression5() {
	// Call in expression statement "condition"
	c.allowDecomposition = true
	c.HelperInlineReferenceToFunction(
		"function foo(a){return true;}; function x() {foo(1)&&1 }",
		`function foo(a){return true;};
function x() { var JSCompiler_inline_result$jscomp$0;
{JSCompiler_inline_result$jscomp$0=true;}
JSCompiler_inline_result$jscomp$0&&1 }
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInlineReferenceInExpression6()
func (c *FunctionInjectorTest) TestInlineReferenceInExpression6() {
	// Call in expression statement after side-effect free "condition"
	c.allowDecomposition = true
	c.HelperInlineReferenceToFunction(
		"function foo(a){return true;}; function x() {1 + foo(1) }",
		`function foo(a){return true;};
function x() { var JSCompiler_inline_result$jscomp$0;
{JSCompiler_inline_result$jscomp$0=true;}
1 + JSCompiler_inline_result$jscomp$0 }
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInlineReferenceInExpression7()
func (c *FunctionInjectorTest) TestInlineReferenceInExpression7() {
	// Call in expression statement "condition"
	c.allowDecomposition = true
	c.HelperInlineReferenceToFunction(
		"function foo(a){return true;}; function x() {foo(1) && 1 }",
		`function foo(a){return true;};
function x() { var JSCompiler_inline_result$jscomp$0;
{JSCompiler_inline_result$jscomp$0=true;}
JSCompiler_inline_result$jscomp$0&&1 }
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInlineReferenceInExpression8()
func (c *FunctionInjectorTest) TestInlineReferenceInExpression8() {
	// Call in expression statement after side-effect free operator
	c.allowDecomposition = true
	c.HelperInlineReferenceToFunction(
		"function foo(a){return true;}; function x() {1 + foo(1) }",
		`function foo(a){return true;};
function x() { var JSCompiler_inline_result$jscomp$0;
{JSCompiler_inline_result$jscomp$0=true;}
1 + JSCompiler_inline_result$jscomp$0 }
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInlineReferenceInExpression9()
func (c *FunctionInjectorTest) TestInlineReferenceInExpression9() {
	// Call in VAR expression.
	c.allowDecomposition = true
	c.HelperInlineReferenceToFunction(
		"function foo(a){return true;}; function x() {var b = 1 + foo(1)}",
		`function foo(a){return true;};
function x() {
var JSCompiler_inline_result$jscomp$0;
{JSCompiler_inline_result$jscomp$0=true;}
var b = 1 + JSCompiler_inline_result$jscomp$0
}
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// TODO(nicksantos): Re-enable with side-effect detection.
// Java: FunctionInjectorTest.testInlineReferenceInExpression10()
func (c *FunctionInjectorTest) TestInlineReferenceInExpression10() {
	c.T().Skip("@Ignore") // @Ignore
	// Call in assignment expression.
	c.HelperInlineReferenceToFunction(
		`/** @nosideeffects */ function foo(a){return true;};
function x() {var b; b += 1 + foo(1) }
`,
		`function foo(a){return true;};
function x() {var b;
{var JSCompiler_inline_result$jscomp$0;
JSCompiler_inline_result$jscomp$0=true;}
b += 1 + JSCompiler_inline_result$jscomp$0 }
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInlineReferenceInExpression11()
func (c *FunctionInjectorTest) TestInlineReferenceInExpression11() {
	// Call under label
	c.allowDecomposition = true
	c.HelperInlineReferenceToFunction(
		"function foo(a){return true;}; function x() {a:foo(1)?0:1 }",
		`function foo(a){return true;};
function x() {
  a:{
    var JSCompiler_inline_result$jscomp$0;
    {JSCompiler_inline_result$jscomp$0=true;}
    JSCompiler_inline_result$jscomp$0?0:1
  }
}
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInlineReferenceInExpression12()
func (c *FunctionInjectorTest) TestInlineReferenceInExpression12() {
	c.allowDecomposition = true
	c.HelperInlineReferenceToFunction(
		"function foo(a){return true;} function x() { 1?foo(1):1; }",
		`function foo(a){return true}
function x() {
  if(1) {
    {true;}
  } else {
    1;
  }
}
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInlineReferenceInExpression13()
func (c *FunctionInjectorTest) TestInlineReferenceInExpression13() {
	c.allowDecomposition = true
	c.HelperInlineReferenceToFunction(
		"function foo(a){return true;}; function x() { goo() + (1?foo(1):1) }",
		`function foo(a){return true;};
function x() { var JSCompiler_temp_const$jscomp$0=goo();
var JSCompiler_temp$jscomp$1;
if(1) {
  {JSCompiler_temp$jscomp$1=true;}
} else {
  JSCompiler_temp$jscomp$1=1;
}
JSCompiler_temp_const$jscomp$0 + JSCompiler_temp$jscomp$1
}
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInlineReferenceInExpression14()
func (c *FunctionInjectorTest) TestInlineReferenceInExpression14() {
	c.allowDecomposition = true
	c.HelperInlineReferenceToFunction(
		`var z = {};
function foo(a){z = {};return true;};
function x() { z.gack = foo(1) }
`,
		`var z = {};
function foo(a){z = {};return true;};
function x() {
var JSCompiler_temp_const$jscomp$0=z;
var JSCompiler_inline_result$jscomp$1;
{
z= {};
JSCompiler_inline_result$jscomp$1 = true;
}
JSCompiler_temp_const$jscomp$0.gack = JSCompiler_inline_result$jscomp$1;
}
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInlineReferenceInExpression15()
func (c *FunctionInjectorTest) TestInlineReferenceInExpression15() {
	c.allowDecomposition = true
	c.HelperInlineReferenceToFunction(
		`var z = {};
function foo(a){z = {};return true;};
function x() { z.gack = foo.call(this,1) }
`,
		`var z = {};
function foo(a){z = {};return true;};
function x() {
var JSCompiler_temp_const$jscomp$0=z;
var JSCompiler_inline_result$jscomp$1;
{
z= {};
JSCompiler_inline_result$jscomp$1 = true;
}
JSCompiler_temp_const$jscomp$0.gack = JSCompiler_inline_result$jscomp$1;
}
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInlineReferenceInExpression16()
func (c *FunctionInjectorTest) TestInlineReferenceInExpression16() {
	c.allowDecomposition = true
	c.HelperInlineReferenceToFunction(
		`var z = {};
function foo(a){z = {};return true;};
function x() { z[bar()] = foo(1) }
`,
		`var z = {};
function foo(a){z = {};return true;};
function x() {
var JSCompiler_temp_const$jscomp$1=z;
var JSCompiler_temp_const$jscomp$0=bar();
var JSCompiler_inline_result$jscomp$2;
{
z= {};
JSCompiler_inline_result$jscomp$2 = true;
}
JSCompiler_temp_const$jscomp$1[JSCompiler_temp_const$jscomp$0] =
JSCompiler_inline_result$jscomp$2;
}
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInlineReferenceToFunction_passingLiteralThis_andSideEffectfulArg()
func (c *FunctionInjectorTest) TestInlineReferenceToFunction_passingLiteralThis_andSideEffectfulArg() {
	c.HelperInlineReferenceToFunction(
		`function foo(y) { bar(this, y, y); }
function x() {foo.call(this, sideEffects());}
`,
		`function foo(y) { bar(this, y, y); }
function x() {
  {
    var y$jscomp$inline_0 = sideEffects();
    bar(this, y$jscomp$inline_0, y$jscomp$inline_0);
  }
}
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInlineReferenceToFunction_passingSideEffectfulThis()
func (c *FunctionInjectorTest) TestInlineReferenceToFunction_passingSideEffectfulThis() {
	c.assumeStrictThis = true
	c.HelperInlineReferenceToFunction(
		`function foo() { bar.call(this, this); }
function x() {foo.call(computeThis());}
`,
		`function foo() { bar.call(this, this); }
function x(){
  {
    var JSCompiler_inline_this_0 = computeThis();
    bar.call(JSCompiler_inline_this_0, JSCompiler_inline_this_0);
  }
}
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInlineReferenceInExpression17()
func (c *FunctionInjectorTest) TestInlineReferenceInExpression17() {
	c.allowDecomposition = true
	c.HelperInlineReferenceToFunction(
		`var z = {};
function foo(a){z = {};return true;};
function x() { z.y.x.gack = foo(1) }
`,
		`var z = {};
function foo(a){z = {};return true;};
function x() {
var JSCompiler_temp_const$jscomp$0=z.y.x;
var JSCompiler_inline_result$jscomp$1;
{
z= {};
JSCompiler_inline_result$jscomp$1 = true;
}
JSCompiler_temp_const$jscomp$0.gack = JSCompiler_inline_result$jscomp$1;
}
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInlineWithinCalls1()
func (c *FunctionInjectorTest) TestInlineWithinCalls1() {
	// Call in within a call
	c.allowDecomposition = true
	c.HelperInlineReferenceToFunction(
		"function foo(){return _g;}; function x() {1 + foo()() }",
		`function foo(){return _g;};
function x() { var JSCompiler_inline_result$jscomp$0;
{JSCompiler_inline_result$jscomp$0=_g;}
1 + JSCompiler_inline_result$jscomp$0() }
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// TODO(nicksantos): Re-enable with side-effect detection.
// Java: FunctionInjectorTest.testInlineWithinCalls2()
func (c *FunctionInjectorTest) TestInlineWithinCalls2() {
	c.T().Skip("@Ignore") // @Ignore
	c.allowDecomposition = true
	c.HelperInlineReferenceToFunction(
		`/** @nosideeffects */ function foo(){return true;};
function x() {1 + _g(foo()) }
`,
		`function foo(){return true;};
function x() { {var JSCompiler_inline_result$jscomp$0;
JSCompiler_inline_result$jscomp$0=true;}
1 + _g(JSCompiler_inline_result$jscomp$0) }
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInlineWithinSuperCall()
func (c *FunctionInjectorTest) TestInlineWithinSuperCall() {
	// Call in within a call
	c.allowDecomposition = true
	c.HelperInlineReferenceToFunction(
		`class A { constructor(g) {} }
class B extends A {
  constructor() { super(foo()); }
}
function foo() {
  return '';
}
`,
		`class A { constructor(g) {} }
class B extends A {
  constructor() {
    var JSCompiler_inline_result$jscomp$0;
    {
      JSCompiler_inline_result$jscomp$0 = '';
    }
    super(JSCompiler_inline_result$jscomp$0)
  }
}
function foo() {
  return '';
}
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInlineWithinSuperCall_followingCall()
func (c *FunctionInjectorTest) TestInlineWithinSuperCall_followingCall() {
	// Call in within a call
	c.allowDecomposition = true
	c.HelperInlineReferenceToFunction(
		`class A { constructor(g) {} }
class B extends A {
  constructor() { super(goo(), foo()); }
}
function foo() {
  return '';
}
`,
		`class A { constructor(g) {} }
class B extends A {
  constructor() {
    var JSCompiler_temp_const$jscomp$0=goo();
    var JSCompiler_inline_result$jscomp$1;
    {
      JSCompiler_inline_result$jscomp$1 = '';
    }
    super(JSCompiler_temp_const$jscomp$0, JSCompiler_inline_result$jscomp$1)
  }
}
function foo() {
  return '';
}
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testInlineAssignmentToConstant()
func (c *FunctionInjectorTest) TestInlineAssignmentToConstant() {
	// Call in within a call
	c.allowDecomposition = true
	c.HelperInlineReferenceToFunction(
		"function foo(){return _g;}; function x(){var CONSTANT_RESULT = foo(); }",
		`function foo(){return _g;};
function x() {
  var JSCompiler_inline_result$jscomp$0;
  {JSCompiler_inline_result$jscomp$0=_g;}
  var CONSTANT_RESULT = JSCompiler_inline_result$jscomp$0;
}
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testBug1897706()
func (c *FunctionInjectorTest) TestBug1897706() {
	c.HelperInlineReferenceToFunction(
		"function foo(a){}; foo(x())",
		"function foo(a){}; {var a$jscomp$inline_0=x()}",
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
	c.HelperInlineReferenceToFunction(
		"function foo(a){bar()}; foo(x())",
		"function foo(a){bar()}; {var a$jscomp$inline_0=x();bar()}",
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
	c.HelperInlineReferenceToFunction(
		"function foo(a,b){bar()}; foo(x(),y())",
		`function foo(a,b){bar()};
{var a$jscomp$inline_0=x();var b$jscomp$inline_1=y();bar()}
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testIssue1101a()
func (c *FunctionInjectorTest) TestIssue1101a() {
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_NO,
		"function foo(a){return modifiyX() + a;} foo(x);",
		"foo",
		FunctionInjectorTest_INLINE_DIRECT)
}

// Java: FunctionInjectorTest.testIssue1101b()
func (c *FunctionInjectorTest) TestIssue1101b() {
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_NO,
		"function foo(a){return (x.prop = 2),a;} foo(x.prop);",
		"foo",
		FunctionInjectorTest_INLINE_DIRECT)
}

// Java: FunctionInjectorTest.testCanInlineReference_isNo_ifCalledWithSpread()
func (c *FunctionInjectorTest) TestCanInlineReference_isNo_ifCalledWithSpread() {
	c.HelperCanInlineReferenceToFunction(
		FunctionInjector_CanInlineResult_NO, //
		`function foo(a) { return a; };
foo(...b);
`,
		"foo",
		FunctionInjectorTest_INLINE_DIRECT)
}

// Java: FunctionInjectorTest.testCanInlineReference_direct_ifArgExpression_containsSpread()
func (c *FunctionInjectorTest) TestCanInlineReference_direct_ifArgExpression_containsSpread() {
	c.HelperInlineReferenceToFunction(
		`function foo(a) { return a; };
foo([...b]);
`,
		`function foo(a) { return a; };
[...b]
`,
		"foo",
		FunctionInjectorTest_INLINE_DIRECT)
}

// Java: FunctionInjectorTest.testCanInlineReference_direct_ifResultIsSpread()
func (c *FunctionInjectorTest) TestCanInlineReference_direct_ifResultIsSpread() {
	c.HelperInlineReferenceToFunction(
		`function foo(b) { return [1, 2, b]; }
bar(...foo(5));
`,
		`function foo(b) { return [1, 2, b]; }
bar(...[1, 2, 5]);
`,
		"foo",
		FunctionInjectorTest_INLINE_DIRECT)
}

// Java: FunctionInjectorTest.testCanInlineReference_block_ifResultIsSpread()
func (c *FunctionInjectorTest) TestCanInlineReference_block_ifResultIsSpread() {
	c.allowDecomposition = true
	c.HelperInlineReferenceToFunction(
		`function foo(b) {
  return [1, 2, b];
}
const bar = function() { };

bar(...foo(5));
`,
		`function foo(b) {
  return [1, 2, b];
}
const bar = function() { };

var JSCompiler_inline_result$jscomp$0
{
  JSCompiler_inline_result$jscomp$0 = [1, 2, 5];
}
bar(...JSCompiler_inline_result$jscomp$0);
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testCanInlineReference_direct_ifPreviousSibling_isSpread()
func (c *FunctionInjectorTest) TestCanInlineReference_direct_ifPreviousSibling_isSpread() {
	c.HelperInlineReferenceToFunction(
		`function foo(b) { return b + 1; }
bar(...qux(), foo(5));
`,
		`function foo(b) { return b + 1; }
bar(...qux(), 5 + 1);
`,
		"foo",
		FunctionInjectorTest_INLINE_DIRECT)
}

// Java: FunctionInjectorTest.testCanInlineReference_block_ifPreviousSibling_isSpread()
func (c *FunctionInjectorTest) TestCanInlineReference_block_ifPreviousSibling_isSpread() {
	c.allowDecomposition = true
	c.HelperInlineReferenceToFunction(
		`function foo(b) {
  return [1, 2, b];
}
const bar = function() { };

bar(...qux(), foo(5));
`,
		`function foo(b) {
  return [1, 2, b];
}
const bar = function() { };

var JSCompiler_temp_const$jscomp$0 = [...qux()];
var JSCompiler_inline_result$jscomp$1
{
  JSCompiler_inline_result$jscomp$1 = [1, 2, 5];
}
bar(...JSCompiler_temp_const$jscomp$0, JSCompiler_inline_result$jscomp$1);
`,
		"foo",
		FunctionInjectorTest_INLINE_BLOCK)
}

// Java: FunctionInjectorTest.testArgumentsReferenceInArrowFunction()
func (c *FunctionInjectorTest) TestArgumentsReferenceInArrowFunction() {
	AssertThat(
		c.DoesFunctionMeetMinimumRequirements(
			"function foo() { return () => arguments; }", "foo")).
		IsFalse()
}

// Java: FunctionInjectorTest.testArgumentsReferenceInNestedVanillaFunction()
func (c *FunctionInjectorTest) TestArgumentsReferenceInNestedVanillaFunction() {
	AssertThat(
		c.DoesFunctionMeetMinimumRequirements(
			"function foo() { return function() { return arguments; }; }", "foo")).
		IsTrue()
}

// HelperCanInlineReferenceToFunction:
//
// # Test case
//
// <p>var a = {}, b = {} a.test = "a", b.test = "b" c = a; foo() { c=b; return "a" } c.teste
//
// Java: FunctionInjectorTest.helperCanInlineReferenceToFunction(CanInlineResult, String, String, InliningMode)
func (c *FunctionInjectorTest) HelperCanInlineReferenceToFunction(
	expectedResult FunctionInjector_CanInlineResult,
	code string,
	fnName string,
	mode FunctionInjector_InliningMode) {
	compiler := NewCompiler()
	compiler.InitOptions(NewCompilerOptions())
	functionArgumentInjector :=
		NewFunctionArgumentInjector(compiler.GetAstAnalyzer())
	injector :=
		NewFunctionInjector_Builder(compiler).
			AllowDecomposition(c.allowDecomposition).
			AssumeStrictThis(c.assumeStrictThis).
			AssumeMinimumCapture(FunctionInjectorTest_ASSUME_MINIMUM_CAPTURE).
			FunctionArgumentInjector(functionArgumentInjector).
			Build()
	tree := FunctionInjectorTest_Parse(compiler, code)

	fnNode := FunctionInjectorTest_FindFunction(tree, fnName)
	unsafe := functionArgumentInjector.FindModifiedParameters(fnNode)

	// can-inline tester
	tester :=
		func(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
			ref := NewFunctionInjector_Reference(n, t.GetScope(), t.GetChunk(), mode)
			fnBody := NodeUtil_GetFunctionBody(fnNode)
			result :=
				injector.CanInlineReferenceToFunction(
					ref,
					fnNode,
					unsafe,
					NodeUtil_ReferencesOwnReceiver(fnNode),
					NodeUtil_Has(
						fnBody,
						NodeUtil_IsFunctionOrNonStaticClassFieldInitializer,
						functionInjectorTest_alwaysTrue))
			AssertThat(result).IsEqualTo(expectedResult)
			return true
		}

	test := NewFunctionInjectorTest_TestCallback(fnName, tester)
	NodeTraversal_Traverse(compiler, tree, test)
}

// Java: FunctionInjectorTest.validateSourceInfo(Compiler, Node)
func (c *FunctionInjectorTest) ValidateSourceInfo(compiler *Compiler, subtree *rhino.Node) {
	NewSourceInfoCheck(compiler).SetCheckSubTree(subtree)
	// Source information problems are reported as compiler errors.
	if compiler.GetErrorCount() != 0 {
		msg := "Error encountered: "
		for _, err := range compiler.GetErrors() {
			msg += err.String() + "\n"
		}
		AssertWithMessage(msg).That(compiler.GetErrorCount()).IsEqualTo(0)
	}
}

// Java: FunctionInjectorTest.helperInlineReferenceToFunction(String, String, String, InliningMode)
func (c *FunctionInjectorTest) HelperInlineReferenceToFunction(
	code string, expectedResult string, fnName string, mode FunctionInjector_InliningMode) {
	compiler := NewCompiler()
	externsInputs := []*SourceFile{SourceFile_FromCode("externs", "")}

	options := NewCompilerOptions()
	options.SetCodingConvention(NewGoogleCodingConvention())
	compiler.Init(externsInputs, []*SourceFile{SourceFile_FromCode("code", code)}, options)

	functionArgumentInjector :=
		NewFunctionArgumentInjector(compiler.GetAstAnalyzer())
	injector :=
		NewFunctionInjector_Builder(compiler).
			AllowDecomposition(c.allowDecomposition).
			AssumeStrictThis(c.assumeStrictThis).
			AssumeMinimumCapture(FunctionInjectorTest_ASSUME_MINIMUM_CAPTURE).
			FunctionArgumentInjector(functionArgumentInjector).
			Build()

	parseRoot := compiler.ParseInputs()
	externsRoot := parseRoot.GetFirstChild()
	tree := parseRoot.GetLastChild()
	AssertThat(tree).IsNotNull()
	AssertThat(tree).IsNotSameInstanceAs(externsRoot)

	expectedRoot := FunctionInjectorTest_ParseExpected(NewCompiler(), expectedResult)

	mainRoot := tree
	Normalize_CreateNormalizeForOptimizations(compiler).Process(externsRoot, mainRoot)
	GatherGetterAndSetterProperties_Update(compiler, externsRoot, mainRoot)
	NewPureFunctionIdentifier_Driver(compiler).Process(externsRoot, mainRoot)

	normalize := Normalize_CreateNormalizeForOptimizations(compiler)
	normalize.Process(externsRoot, mainRoot)
	compiler.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_NORMALIZED)

	fnNode := FunctionInjectorTest_FindFunction(tree, fnName)
	AssertThat(fnNode).IsNotNull()
	unsafe := functionArgumentInjector.FindModifiedParameters(fnNode)
	AssertThat(fnNode).IsNotNull()

	// inline tester
	tester :=
		func(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
			ref := NewFunctionInjector_Reference(n, t.GetScope(), t.GetChunk(), mode)
			fnBody := NodeUtil_GetFunctionBody(fnNode)
			canInline :=
				injector.CanInlineReferenceToFunction(
					ref,
					fnNode,
					unsafe,
					NodeUtil_ReferencesOwnReceiver(fnNode),
					NodeUtil_Has(
						fnBody,
						NodeUtil_IsFunctionOrNonStaticClassFieldInitializer,
						functionInjectorTest_alwaysTrue))
			AssertWithMessage("canInlineReferenceToFunction should not be CAN_NOT_INLINE").
				That(canInline).
				IsNotEqualTo(FunctionInjector_CanInlineResult_NO)
			if c.allowDecomposition {
				AssertWithMessage("canInlineReferenceToFunction should be AFTER_PREPARATION").
					That(canInline).
					IsEqualTo(FunctionInjector_CanInlineResult_AFTER_PREPARATION)

				injector.MaybePrepareCall(ref)

				AssertWithMessage("canInlineReferenceToFunction should be CAN_INLINE").
					That(canInline).
					IsNotEqualTo(FunctionInjector_CanInlineResult_YES)
			}

			result := injector.Inline(ref, fnName, fnNode)
			c.ValidateSourceInfo(compiler, result)

			AssertNode(tree.GetFirstChild()).
				UsingSerializer(FunctionInjectorTest_ToSource).
				IsEqualTo(expectedRoot)

			return true
		}

	verifier := NewChangeVerifier(compiler).Snapshot(mainRoot)
	test := NewFunctionInjectorTest_TestCallback(fnName, tester)
	NodeTraversal_Traverse(compiler, tree, test)
	verifier.CheckRecordedChangesWithPassName("helperInlineReferenceToFunction", mainRoot)
}

// DoesFunctionMeetMinimumRequirements calls
// FunctionInjector#doesFunctionMeetMinimumRequirements(String, Node).
//
// This method is called as a prerequisite to checking if a particular
// reference is inlinable
//
// Java: FunctionInjectorTest.doesFunctionMeetMinimumRequirements(String, String)
func (c *FunctionInjectorTest) DoesFunctionMeetMinimumRequirements(code string, fnName string) bool {
	compiler := NewCompiler()
	compiler.InitOptions(NewCompilerOptions())
	injector :=
		NewFunctionInjector_Builder(compiler).
			AllowDecomposition(c.allowDecomposition).
			AssumeStrictThis(c.assumeStrictThis).
			AssumeMinimumCapture(FunctionInjectorTest_ASSUME_MINIMUM_CAPTURE).
			Build()
	tree := FunctionInjectorTest_Parse(compiler, code)

	fnNode := FunctionInjectorTest_FindFunction(tree, fnName)
	return injector.DoesFunctionMeetMinimumRequirements(fnName, fnNode)
}

// FunctionInjectorTest_Method is the interface Method.
type FunctionInjectorTest_Method func(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool

// FunctionInjectorTest_TestCallback is the static class TestCallback.
type FunctionInjectorTest_TestCallback struct {
	callname string
	method   FunctionInjectorTest_Method
	complete bool
}

// Java: FunctionInjectorTest.TestCallback(String, Method)
func NewFunctionInjectorTest_TestCallback(callname string, method FunctionInjectorTest_Method) *FunctionInjectorTest_TestCallback {
	return &FunctionInjectorTest_TestCallback{callname: callname, method: method}
}

// Java: FunctionInjectorTest.TestCallback.shouldTraverse(NodeTraversal, Node, Node)
func (cb *FunctionInjectorTest_TestCallback) ShouldTraverse(nodeTraversal *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	return !cb.complete
}

// Java: FunctionInjectorTest.TestCallback.visit(NodeTraversal, Node, Node)
func (cb *FunctionInjectorTest_TestCallback) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if n.IsCall() {
		var callee *rhino.Node
		if NodeUtil_IsNormalGet(n.GetFirstChild()) {
			callee = n.GetFirstFirstChild()
		} else {
			callee = n.GetFirstChild()
		}

		if callee.IsName() && callee.GetString() == cb.callname {
			cb.complete = cb.method(t, n, parent)
		}
	}

	if parent == nil {
		AssertThat(cb.complete).IsTrue()
	}
}

// Java: FunctionInjectorTest.findFunction(Node, String)
func FunctionInjectorTest_FindFunction(n *rhino.Node, name string) *rhino.Node {
	if n.IsFunction() {
		if n.GetFirstChild().GetString() == name {
			return n
		}
	}

	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		result := FunctionInjectorTest_FindFunction(c, name)
		if result != nil {
			return result
		}
	}

	return nil
}

// Java: FunctionInjectorTest.prep(String)
func FunctionInjectorTest_Prep(js string) *rhino.Node {
	compiler := NewCompiler()
	n := compiler.ParseTestCode(js)
	AssertThat(compiler.GetErrorCount()).IsEqualTo(0)
	return n.GetFirstChild()
}

// Java: FunctionInjectorTest.parse(Compiler, String)
func FunctionInjectorTest_Parse(compiler *Compiler, js string) *rhino.Node {
	n := compiler.ParseTestCode(js)
	AssertThat(compiler.GetErrorCount()).IsEqualTo(0)
	return n
}

// Java: FunctionInjectorTest.parseExpected(Compiler, String)
func FunctionInjectorTest_ParseExpected(compiler *Compiler, js string) *rhino.Node {
	n := compiler.ParseTestCode(js)
	message := "Unexpected errors: "
	errs := compiler.GetErrors()
	for _, element := range errs {
		message += "\n" + element.String()
	}
	AssertWithMessage(message).That(compiler.GetErrorCount()).IsEqualTo(0)
	return n
}

// functionInjectorTest_alwaysTrue is Guava's Predicates.alwaysTrue().
func functionInjectorTest_alwaysTrue(*rhino.Node) bool { return true }

// Java: FunctionInjectorTest.toSource(Node)
func FunctionInjectorTest_ToSource(n *rhino.Node) string {
	return NewCodePrinter_Builder(n).SetPrettyPrint(true).SetLineBreak(true).Build()
}
