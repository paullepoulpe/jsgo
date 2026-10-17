// Port of com/google/javascript/jscomp/ExpressionDecomposerTest.java @ closure-compiler v20261015

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
	"fmt"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Unit tests for ExpressionDecomposer
// Note: functions "foo" and "goo" are external functions in the helper.
//
// The Java @Test methods not named test* get a "Test" prefix so that
// RunJUnit4 runs them.
//
// Java: public final class ExpressionDecomposerTest
type ExpressionDecomposerTest struct {
	// Go-only: test-helper — the Go test running the current test method.
	t *testing.T

	// Go-only: language-adapter — a HashSet<String> that is only copied
	// into the decomposer's lookup set, so its order is not observable.
	knownConstants *jutil.LinkedHashSet[string]

	// The language out to set in the compiler options. If null, use the default.
	languageOut *CompilerOptions_LanguageMode // @Nullable

	// Whether we should run type checking and test the type information in the output expression
	shouldTestTypes bool
}

func TestExpressionDecomposerTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase {
		return &ExpressionDecomposerTest{knownConstants: jutil.NewLinkedHashSet[string]()}
	})
}

// JUnitSetT gives the instance the Go test running its test method.
func (c *ExpressionDecomposerTest) JUnitSetT(t *testing.T) { c.t = t }

// JUnitBefore runs the @Before method.
func (c *ExpressionDecomposerTest) JUnitBefore() { c.SetUp() }

// JUnitAfter: the class has no @After method.
func (c *ExpressionDecomposerTest) JUnitAfter() {}

// Java: ExpressionDecomposerTest.setUp()
func (c *ExpressionDecomposerTest) SetUp() {
	c.knownConstants.Clear()
	c.shouldTestTypes = true
	c.languageOut = nil
}

// expressionDecomposerTest_nodeFinder is the tests'
// Function<AbstractCompiler, Function<Node, Node>>.
type expressionDecomposerTest_nodeFinder func(compiler AbstractCompiler) func(root *rhino.Node) *rhino.Node

// Java: ExpressionDecomposerTest.testWindowLocationAssign()
func (c *ExpressionDecomposerTest) TestWindowLocationAssign() {
	// avoid decomposing `window.location.assign` when the output code could
	// end up running on IE11. See more explanation in ExpressionDecomposer.java.
	c.languageOut = CompilerOptions_LanguageMode_ECMASCRIPT5.Ptr()
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "window.location.assign(foo())", exprMatchesStr("foo()"))
	c.helperMoveExpression(
		"window.location.assign(foo())",
		exprMatchesStr("foo()"),
		`var result$jscomp$0 = foo();
window.location.assign(result$jscomp$0)
`)

	// confirm that the default behavior does not treat window.location.assign
	// specially
	c.languageOut = nil
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "window.location.assign(foo())", exprMatchesStr("foo()"))
	c.helperExposeExpression(
		"window.location.assign(foo())",
		exprMatchesStr("foo()"),
		`var temp_const$jscomp$1 = window.location;
var temp_const$jscomp$0 = temp_const$jscomp$1.assign;
temp_const$jscomp$0.call(temp_const$jscomp$1, foo());
`)
}

// Java: ExpressionDecomposerTest.testObjectDestructuring_withComputedKey_doesNotCrash()
func (c *ExpressionDecomposerTest) TestObjectDestructuring_withComputedKey_doesNotCrash() {
	// computed prop is found to be decomposable
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "var a; ({ [foo()]: a} = obj);", exprMatchesStr("foo()"))

	// TODO(b/339040894): Fix this crash.
	ex :=
		AssertThrows[*jutil.IllegalStateException](func() {
			c.helperExposeExpression(
				"var a; ({ [foo()]: a} = obj);", //
				exprMatchesStr("foo()"),
				`var a;
var temp_const$jscomp$0 = obj;
var temp_const$jscomp$1 = foo();
({ [temp_const$jscomp$1]: a} = temp_const$jscomp$0);
`)
		})
	AssertThat(ex.GetMessage()).Contains("exposeExpression exposed nothing")
}

// Java: ExpressionDecomposerTest.testCannotExpose_expression1()
func (c *ExpressionDecomposerTest) TestCannotExpose_expression1() {
	// Can't move or decompose some classes of expressions.
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE, "while(foo());", exprMatchesStr("foo()"))
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE, "while(x = goo()&&foo()){}", exprMatchesStr("foo()"))
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE, "while(x += goo()&&foo()){}", exprMatchesStr("foo()"))

	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE, "do{}while(foo());", exprMatchesStr("foo()"))
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE, "for(;foo(););", exprMatchesStr("foo()"))
	// This case could be supported for loops without conditional continues
	// by moving the increment into the loop body.
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE, "for(;;foo());", exprMatchesStr("foo()"))
	c.helperCanExposeExpression(ExpressionDecomposer_DecompositionType_MOVABLE, "for(foo();;);", exprMatchesStr("foo()"))

	// This is potentially doable but a bit too complex currently.
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE, "switch(1){case foo():;}", exprMatchesStr("foo()"))
}

// Java: ExpressionDecomposerTest.testCanExposeExpression2()
func (c *ExpressionDecomposerTest) TestCanExposeExpression2() {
	c.helperCanExposeExpression(ExpressionDecomposer_DecompositionType_MOVABLE, "foo()", exprMatchesStr("foo()"))
	c.helperCanExposeExpression(ExpressionDecomposer_DecompositionType_MOVABLE, "x = foo()", exprMatchesStr("foo()"))
	c.helperCanExposeExpression(ExpressionDecomposer_DecompositionType_MOVABLE, "var x = foo()", exprMatchesStr("foo()"))
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "const x = foo()", exprMatchesStr("foo()"))
	c.helperCanExposeExpression(ExpressionDecomposer_DecompositionType_MOVABLE, "let x = foo()", exprMatchesStr("foo()"))
	c.helperCanExposeExpression(ExpressionDecomposer_DecompositionType_MOVABLE, "if(foo()){}", exprMatchesStr("foo()"))
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "switch(foo()){}", exprMatchesStr("foo()"))
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "switch(foo()){}", exprMatchesStr("foo()"))
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "function f(){ return foo();}", exprMatchesStr("foo()"))

	c.helperCanExposeExpression(ExpressionDecomposer_DecompositionType_MOVABLE, "x = foo() && 1", exprMatchesStr("foo()"))
	c.helperCanExposeExpression(ExpressionDecomposer_DecompositionType_MOVABLE, "x = foo() || 1", exprMatchesStr("foo()"))
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "x = foo() ? 0 : 1", exprMatchesStr("foo()"))
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "(function(a){b = a})(foo())", exprMatchesStr("foo()"))
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "function f(){ throw foo();}", exprMatchesStr("foo()"))
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: ExpressionDecomposerTest.nullishCoalesceMovable()
func (c *ExpressionDecomposerTest) TestNullishCoalesceMovable() {
	c.helperCanExposeExpression(ExpressionDecomposer_DecompositionType_MOVABLE, "x = foo() ?? 1", exprMatchesStr("foo()"))
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: ExpressionDecomposerTest.nullishCoalesceDecomposable()
func (c *ExpressionDecomposerTest) TestNullishCoalesceDecomposable() {
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "var x = null ?? foo()", exprMatchesStr("foo()"))
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: ExpressionDecomposerTest.nullishCoalesceUnDecomposable()
func (c *ExpressionDecomposerTest) TestNullishCoalesceUnDecomposable() {
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE, "while(x = goo()??foo()){}", exprMatchesStr("foo()"))
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: ExpressionDecomposerTest.assignCoalesceMovable()
func (c *ExpressionDecomposerTest) TestAssignCoalesceMovable() {
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "x ??= goo() + foo()", exprMatchesStr("foo()"))
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: ExpressionDecomposerTest.optChainMovable()
func (c *ExpressionDecomposerTest) TestOptChainMovable() {
	c.helperCanExposeExpression(ExpressionDecomposer_DecompositionType_MOVABLE, "foo()?.x", exprMatchesStr("foo()"))
	c.helperCanExposeExpression(ExpressionDecomposer_DecompositionType_MOVABLE, "foo()?.[x]", exprMatchesStr("foo()"))
	c.helperCanExposeExpression(ExpressionDecomposer_DecompositionType_MOVABLE, "foo()?.()", exprMatchesStr("foo()"))
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: ExpressionDecomposerTest.optChainDecomposable()
func (c *ExpressionDecomposerTest) TestOptChainDecomposable() {
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "x?.[foo()]", exprMatchesStr("foo()"))
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "x?.(foo())", exprMatchesStr("foo()"))
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: ExpressionDecomposerTest.optChainAllowMethodCallDecomposable()
func (c *ExpressionDecomposerTest) TestOptChainAllowMethodCallDecomposable() {
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "x?.y(foo())", exprMatchesStr("foo()"))
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: ExpressionDecomposerTest.optChainUnDecomposable()
func (c *ExpressionDecomposerTest) TestOptChainUnDecomposable() {
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE, "while(x = y?.[foo()]){}", exprMatchesStr("foo()"))
}

// Java: ExpressionDecomposerTest.testCanExposeExpression3()
func (c *ExpressionDecomposerTest) TestCanExposeExpression3() {
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "x = 0 && foo()", exprMatchesStr("foo()"))
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "x = 1 || foo()", exprMatchesStr("foo()"))
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "var x = 1 ? foo() : 0", exprMatchesStr("foo()"))
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "const x = 1 ? foo() : 0", exprMatchesStr("foo()"))
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "let x = 1 ? foo() : 0", exprMatchesStr("foo()"))

	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "goo() && foo()", exprMatchesStr("foo()"))
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "x = goo() && foo()", exprMatchesStr("foo()"))
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "x += goo() && foo()", exprMatchesStr("foo()"))
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "var x = goo() && foo()", exprMatchesStr("foo()"))
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "const x = goo() && foo()", exprMatchesStr("foo()"))
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "let x = goo() && foo()", exprMatchesStr("foo()"))
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "if(goo() && foo()){}", exprMatchesStr("foo()"))
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "switch(goo() && foo()){}", exprMatchesStr("foo()"))
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "switch(goo() && foo()){}", exprMatchesStr("foo()"))
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "switch(x = goo() && foo()){}", exprMatchesStr("foo()"))
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE,
		"function f(){ return goo() && foo();}",
		exprMatchesStr("foo()"))
}

// Java: ExpressionDecomposerTest.testCanExposeExpression_compoundDeclaration_inForInitializer_firstElement()
func (c *ExpressionDecomposerTest) TestCanExposeExpression_compoundDeclaration_inForInitializer_firstElement() {
	// VAR will already be hoisted by `Normalize`.
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "for (var x = foo(), y = 5;;) {}", exprMatchesStr("foo()"))

	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "for (let x = foo(), y = 5;;) {}", exprMatchesStr("foo()"))
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "for (const x = foo(), y = 5;;) {}", exprMatchesStr("foo()"))
}

// Java: ExpressionDecomposerTest.testCanExposeExpression_compoundDeclaration_inForInitializer_nthElement()
func (c *ExpressionDecomposerTest) TestCanExposeExpression_compoundDeclaration_inForInitializer_nthElement() {
	// TODO(b/121157467) FOR introduces complex scoping that isn't currently `Normalize`d.
	// Since in some cases we'd effectively end up having to `Normalize` these, decomposition just
	// bails for now.

	// VAR will already be hoisted by `Normalize`.
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "for (var x = 8, y = foo();;) {}", exprMatchesStr("foo()"))

	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE,
		"for (let x = 8, y = foo();;) {}",
		exprMatchesStr("foo()"))
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE,
		"for (const x = 8, y = foo();;) {}",
		exprMatchesStr("foo()"))
}

// Java: ExpressionDecomposerTest.testCannotExpose_expression4()
func (c *ExpressionDecomposerTest) TestCannotExpose_expression4() {
	// 'this' must be preserved in call.
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "if (goo.a(1, foo()));", exprMatchesStr("foo()"))
}

// Java: ExpressionDecomposerTest.testCannotExpose_expression5()
func (c *ExpressionDecomposerTest) TestCannotExpose_expression5() {
	// 'this' must be preserved in call.
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "if (goo['a'](foo()));", exprMatchesStr("foo()"))
}

// Java: ExpressionDecomposerTest.testCannotExpose_expression6()
func (c *ExpressionDecomposerTest) TestCannotExpose_expression6() {
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "z:if (goo.a(1, foo()));", exprMatchesStr("foo()"))
}

// Java: ExpressionDecomposerTest.testCanExposeExpression7()
func (c *ExpressionDecomposerTest) TestCanExposeExpression7() {
	// Verify calls to function expressions are movable.
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE,
		`(function(map){descriptions_=map})(
  function(){
    var ret={};
    ret[INIT]='a';
    ret[MIGRATION_BANNER_DISMISS]='b';
    return ret
  }());
`,
		func(compiler AbstractCompiler) func(*rhino.Node) *rhino.Node {
			return func(rootNode *rhino.Node) *rhino.Node {
				// Dig out the inner IIFE call and return it as the expression we want to ensure is
				// movable.
				AssertNode(rootNode).HasToken(rhino.Token_SCRIPT)
				exprResult := rootNode.GetOnlyChild()
				AssertNode(exprResult).HasToken(rhino.Token_EXPR_RESULT)
				outerCall := exprResult.GetOnlyChild()
				AssertNode(outerCall).HasToken(rhino.Token_CALL)
				innerIifeCall := outerCall.GetSecondChild()
				AssertNode(innerIifeCall).HasToken(rhino.Token_CALL)
				return innerIifeCall
			}
		})
}

// Java: ExpressionDecomposerTest.testCanExposeExpression8()
func (c *ExpressionDecomposerTest) TestCanExposeExpression8() {
	// Can it be decompose?
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE,
		`HangoutStarter.prototype.launchHangout = function() {
  var self = a.b;
  var myUrl = new goog.Uri(
      getDomServices_(self).getDomHelper().getWindow().location.href);
};
`,
		exprMatchesStr("getDomServices_(self)"))

	// Verify it is properly expose the target expression.
	c.helperExposeExpression(
		`HangoutStarter.prototype.launchHangout = function() {
  var self = a.b;
  var myUrl =
      new goog.Uri(getDomServices_(self).getDomHelper().getWindow().location.href);
};
`,
		exprMatchesStr("getDomServices_(self)"),
		`HangoutStarter.prototype.launchHangout = function() {
  var self = a.b;
  var temp_const$jscomp$0 = goog.Uri;
  var myUrl = new temp_const$jscomp$0(
      getDomServices_(self).getDomHelper().getWindow().location.href);
}
`)

	// Verify the results can be properly moved.
	c.helperMoveExpression(
		`HangoutStarter.prototype.launchHangout = function() {
  var self = a.b;
  var temp_const$jscomp$0 = goog.Uri;
  var myUrl = new temp_const$jscomp$0(
      getDomServices_(self).getDomHelper().getWindow().location.href);
}
`,
		exprMatchesStr("getDomServices_(self)"),
		`HangoutStarter.prototype.launchHangout = function() {
  var self=a.b;
  var temp_const$jscomp$0=goog.Uri;
  var result$jscomp$0=getDomServices_(self);
  var myUrl=new temp_const$jscomp$0(
      result$jscomp$0.getDomHelper().getWindow().location.href);
}
`)
}

// Java: ExpressionDecomposerTest.testCannotExpose_expression9()
func (c *ExpressionDecomposerTest) TestCannotExpose_expression9() {
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE,
		"function *f() { for (let x of yield y) {} }",
		exprMatchesStr("yield y"))
}

// Java: ExpressionDecomposerTest.testCannotExpose_forAwaitOf()
func (c *ExpressionDecomposerTest) TestCannotExpose_forAwaitOf() {
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE,
		"async function *f() { for await (let x of yield y) {} }",
		exprMatchesStr("yield y"))
}

// Java: ExpressionDecomposerTest.testCanExposeExpression10()
func (c *ExpressionDecomposerTest) TestCanExposeExpression10() {
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE,
		"function *f() { for (let x in yield y) {} }",
		exprMatchesStr("yield y"))
}

// Java: ExpressionDecomposerTest.testCannotExpose_expression11()
func (c *ExpressionDecomposerTest) TestCannotExpose_expression11() {
	// expressions in parameter lists
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE, "function f(x = foo()) {}", exprMatchesStr("foo()"))

	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE, "function f({[foo()]: x}) {}", exprMatchesStr("foo()"))

	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE, "(function (x = foo()) {})()", exprMatchesStr("foo()"))

	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE,
		"(function ({[foo()]: x}) {})()",
		exprMatchesStr("foo()"))
}

// Java: ExpressionDecomposerTest.testCanExpose_aCall_withSpreadSibling()
func (c *ExpressionDecomposerTest) TestCanExpose_aCall_withSpreadSibling() {
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "f(...x, y());", exprMatchesStr("y()"))
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "f(y(), ...x);", exprMatchesStr("y()"))

	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "new D(...x, y());", exprMatchesStr("y()"))
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "new D(y(), ...x);", exprMatchesStr("y()"))

	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "[...x, y()];", exprMatchesStr("y()"))
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "({...x, z: y()});", exprMatchesStr("y()"))

	// Array- and object-literal instantiations cannot be side-effected.
	c.helperCanExposeExpression(ExpressionDecomposer_DecompositionType_MOVABLE, "[y(), ...x];", exprMatchesStr("y()"))
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "({z: y(), ...x});", exprMatchesStr("y()"))

	c.helperCanExposeExpression(ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "f(...y());", exprMatchesStr("y()"))
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "f(...y(), x);", exprMatchesStr("y()"))
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "f(x, ...y());", exprMatchesStr("y()"))

	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "f(...x, x, y());", exprMatchesStr("y()"))
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "f(...x, ...x, y());", exprMatchesStr("y()"))
}

// Java: ExpressionDecomposerTest.testCanExpose_anExpression_withSpreadRelative_ifInDifferentFunction()
func (c *ExpressionDecomposerTest) TestCanExpose_anExpression_withSpreadRelative_ifInDifferentFunction() {
	// TODO(b/121004488): There are potential decompositions that weren't implemented.
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "f(function() { [...x]; }, y());", exprMatchesStr("y()"))
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "f(function() { ({...x}); }, y());", exprMatchesStr("y()"))

	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "f(y(), () => [...x]);", exprMatchesStr("y()"))
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "f(y(), () => ({...x}));", exprMatchesStr("y()"))

	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "[() => f(...x), y()];", exprMatchesStr("y()"))

	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE,
		`[
   class {
     f(x) { return [...x]; }
   },
  y()
];
`,
		exprMatchesStr("y()"))
}

// Java: ExpressionDecomposerTest.testCanExposeExpression12()
func (c *ExpressionDecomposerTest) TestCanExposeExpression12() {
	// Test destructuring rhs is evaluated before the lhs
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "const {a, b = goo()} = foo();", exprMatchesStr("foo()"))

	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "const [a, b = goo()] = foo();", exprMatchesStr("foo()"))

	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "({a, b = goo()} = foo());", exprMatchesStr("foo()"))

	// Default value expressions are conditional, which would make the expressions complex.
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE,
		// default value inside array pattern
		"[{ [foo()]: a } = goo()] = arr;",
		exprMatchesStr("foo()"))

	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, //
		// computed property inside object pattern; decomposed
		"({ [foo()]: a = goo()} = arr);",
		exprMatchesStr("foo()"))

	// default value expressions are conditional, which would make the expressions complex
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE, //
		// default value inside object pattern
		"({ [foo()]: a = goo()} = arr);",
		exprMatchesStr("goo()"))

	c.helperExposeExpression(
		`var Di = I(() => {
  function zv() {
    JSCOMPILER_PRESERVE(e), [getObj().propName] = CN();
  }
  function CN() {
    return [1];
  }
});
`,
		exprMatchesStr("CN()"),
		`var Di = I(() => {
  function zv() {
    var temp_const$jscomp$0 = JSCOMPILER_PRESERVE(e);
// TODO(b/339701959): We decided to back off here because of b/338660589. But we should
// optimize this to:
//   var temp_const$jscomp$1 = CN();
//   var temp_const$jscomp$2 = getObj();
//   temp_const$jscomp$0, [temp_const$jscomp$2.propName] = temp_const$jscomp$1;
    temp_const$jscomp$0, [getObj().propName] = CN();
  }
  function CN() {
    return [1];
  }
});
`)

	c.helperExposeExpression(
		`var Di = I(() => {
  function zv() {
    JSCOMPILER_PRESERVE(e), ({x: getObj().propName} = CN());
  }
  function CN() {
    return {x: 1};
  }
});
`,
		exprMatchesStr("CN()"),
		`var Di = I(() => {
  function zv() {
    var temp_const$jscomp$0 = JSCOMPILER_PRESERVE(e);
// TODO(b/339701959): We decided to back off here because of b/338246627. But we should
// optimize this to:
//   var temp_const$jscomp$1 = CN();
//   var temp_const$jscomp$2 = getObj();
//   temp_const$jscomp$0, {x: temp_const$jscomp$2.propName} = temp_const$jscomp$1;
    temp_const$jscomp$0, {x:getObj().propName} = CN();
  }
  function CN() {
    return {x: 1};
  }
});
`)

	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE,
		`var Di = I(() => {
  function zv() {
    JSCOMPILER_PRESERVE(e), [f] = CN();
  }
  function CN() {
    return [1];
  }
});
`,
		exprMatchesStr("CN()"))

	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE,
		`var Di = I(() => {
  function zv() {
    JSCOMPILER_PRESERVE(e), [f = foo()] = CN();
  }
  function CN() {
    return [1];
  }
});
`,
		exprMatchesStr("foo()"))

	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE,
		`var Di = I(() => {
  function zv() {
    JSCOMPILER_PRESERVE(e), ({f: g} = CN());
  }
  function CN() {
    return {f: 1};
  }
});
`,
		exprMatchesStr("CN()"))

	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE,
		`var Di = I(() => {
  function zv() {
    JSCOMPILER_PRESERVE(e), ({f: g = goo()} = CN());
  }
  function CN() {
    return {f: 1};
  }
});
`,
		exprMatchesStr("goo()"))
}

// Java: ExpressionDecomposerTest.testObjectDestructuring_withDefaultValue_generatesValidAST()
func (c *ExpressionDecomposerTest) TestObjectDestructuring_withDefaultValue_generatesValidAST() {
	c.helperExposeExpression(
		"var d; ({c: d = 4} = condition ? y() :  {c: 1});",
		exprMatchesStr("y()"),
		`var d;
var temp$jscomp$0;
if (condition) {
  temp$jscomp$0 = y();
} else {
  temp$jscomp$0 = {c: 1};
}
({c: d = 4} = temp$jscomp$0);
`)
}

// Java: ExpressionDecomposerTest.testObjectDestructuring_withDefaultValue_withComputedKey()
func (c *ExpressionDecomposerTest) TestObjectDestructuring_withDefaultValue_withComputedKey() {
	// default value expressions are conditional, which would make the expressions complex
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE,
		"var a; ({ [foo()]: a = bar()} = baz());",
		exprMatchesStr("bar()"))

	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE,
		"var a; ({ [foo()]: a = bar()} = baz());",
		exprMatchesStr("baz()"))

	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE,
		"var a; ({ [foo()]: a = bar()} = baz());",
		exprMatchesStr("foo()"))
}

// Java: ExpressionDecomposerTest.testArrayDestructuring_withDefaultValue_generatesValidAST()
func (c *ExpressionDecomposerTest) TestArrayDestructuring_withDefaultValue_generatesValidAST() {
	c.helperExposeExpression(
		"var [c = 4] = condition ? y() :  [c = 2];",
		exprMatchesStr("y()"),
		`var c;
var temp$jscomp$0;
if (condition) {
  temp$jscomp$0 = y();
} else {
  temp$jscomp$0 = [c = 2];
}
[c = 4] = temp$jscomp$0;
`)
}

// Java: ExpressionDecomposerTest.testCanExposeExpressionInTemplateLiteralSubstitution()
func (c *ExpressionDecomposerTest) TestCanExposeExpressionInTemplateLiteralSubstitution() {
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "const result = `${foo()}`;", exprMatchesStr("foo()"))

	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE,
		"const obj = {f(x) {}}; obj.f(`${foo()}`);",
		exprMatchesStr("foo()"))

	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "const result = `${foo()} ${goo()}`;", exprMatchesStr("foo()"))

	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE,
		"const result = `${foo()} ${goo()}`;",
		exprMatchesStr("goo()"))
}

// Java: ExpressionDecomposerTest.testCannotExpose_defaultValueInParamList()
func (c *ExpressionDecomposerTest) TestCannotExpose_defaultValueInParamList() {
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE, "function fn(a = g()) {}", exprMatchesStr("g()"))
}

// Java: ExpressionDecomposerTest.testCannotExpose_defaultValueInDestructuring()
func (c *ExpressionDecomposerTest) TestCannotExpose_defaultValueInDestructuring() {
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE, "let {x = fn()} = y;", exprMatchesStr("fn()"))

	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE, "let [x = fn()] = y;", exprMatchesStr("fn()"))
}

// Java: ExpressionDecomposerTest.testMoveExpression1()
func (c *ExpressionDecomposerTest) TestMoveExpression1() {
	// There isn't a reason to do this, but it works.
	c.helperMoveExpression(
		"foo()", exprMatchesStr("foo()"), "var result$jscomp$0 = foo(); result$jscomp$0;")
}

// Java: ExpressionDecomposerTest.testMoveExpression2()
func (c *ExpressionDecomposerTest) TestMoveExpression2() {
	c.helperMoveExpression(
		"x = foo()", exprMatchesStr("foo()"), "var result$jscomp$0 = foo(); x = result$jscomp$0;")
}

// Java: ExpressionDecomposerTest.testMoveExpression3()
func (c *ExpressionDecomposerTest) TestMoveExpression3() {
	c.helperMoveExpression(
		"var x = foo()",
		exprMatchesStr("foo()"),
		"var result$jscomp$0 = foo(); var x = result$jscomp$0;")
}

// Java: ExpressionDecomposerTest.testMoveExpression4()
func (c *ExpressionDecomposerTest) TestMoveExpression4() {
	c.helperMoveExpression(
		"const x = foo()",
		exprMatchesStr("foo()"),
		"var result$jscomp$0 = foo(); const x = result$jscomp$0;")
}

// Java: ExpressionDecomposerTest.testMoveExpression5()
func (c *ExpressionDecomposerTest) TestMoveExpression5() {
	c.helperMoveExpression(
		"let x = foo()",
		exprMatchesStr("foo()"),
		"var result$jscomp$0 = foo(); let x = result$jscomp$0;")
}

// Java: ExpressionDecomposerTest.testMoveExpression6()
func (c *ExpressionDecomposerTest) TestMoveExpression6() {
	c.helperMoveExpression(
		"if(foo()){}",
		exprMatchesStr("foo()"),
		"var result$jscomp$0 = foo(); if (result$jscomp$0);")
}

// Java: ExpressionDecomposerTest.testMoveExpression7()
func (c *ExpressionDecomposerTest) TestMoveExpression7() {
	c.helperMoveExpression(
		"switch(foo()){}",
		exprMatchesStr("foo()"),
		"var result$jscomp$0 = foo(); switch(result$jscomp$0){}")
}

// Java: ExpressionDecomposerTest.testMoveExpression8()
func (c *ExpressionDecomposerTest) TestMoveExpression8() {
	c.helperMoveExpression(
		"switch(1 + foo()){}",
		exprMatchesStr("foo()"),
		"var result$jscomp$0 = foo(); switch(1 + result$jscomp$0){}")
}

// Java: ExpressionDecomposerTest.testMoveExpression9()
func (c *ExpressionDecomposerTest) TestMoveExpression9() {
	c.helperMoveExpression(
		"function f(){ return foo();}",
		exprMatchesStr("foo()"),
		"function f(){ var result$jscomp$0 = foo(); return result$jscomp$0;}")
}

// Java: ExpressionDecomposerTest.testMoveExpression10()
func (c *ExpressionDecomposerTest) TestMoveExpression10() {
	c.helperMoveExpression(
		"x = foo() && 1",
		exprMatchesStr("foo()"),
		"var result$jscomp$0 = foo(); x = result$jscomp$0 && 1")
}

// Java: ExpressionDecomposerTest.testMoveExpression11()
func (c *ExpressionDecomposerTest) TestMoveExpression11() {
	c.helperMoveExpression(
		"x = foo() || 1",
		exprMatchesStr("foo()"),
		"var result$jscomp$0 = foo(); x = result$jscomp$0 || 1")
}

// Java: ExpressionDecomposerTest.testMoveExpression12()
func (c *ExpressionDecomposerTest) TestMoveExpression12() {
	c.helperMoveExpression(
		"x = foo() ? 0 : 1",
		exprMatchesStr("foo()"),
		"var result$jscomp$0 = foo(); x = result$jscomp$0 ? 0 : 1")
}

// Java: ExpressionDecomposerTest.testMoveExpressionNullishCoalesce()
func (c *ExpressionDecomposerTest) TestMoveExpressionNullishCoalesce() {
	c.helperMoveExpression(
		"x = foo() ?? 0",
		exprMatchesStr("foo()"),
		"var result$jscomp$0 = foo(); x = result$jscomp$0 ?? 0")
}

// Java: ExpressionDecomposerTest.testMoveExpressionOptionalChain()
func (c *ExpressionDecomposerTest) TestMoveExpressionOptionalChain() {
	c.helperMoveExpression(
		"foo()?.x", exprMatchesStr("foo()"), "var result$jscomp$0 = foo(); result$jscomp$0?.x")
}

// Java: ExpressionDecomposerTest.testMoveExpression13()
func (c *ExpressionDecomposerTest) TestMoveExpression13() {
	c.helperMoveExpression(
		"const {a, b} = foo();",
		exprMatchesStr("foo()"),
		"var result$jscomp$0 = foo(); const {a, b} = result$jscomp$0;")
}

// Java: ExpressionDecomposerTest.testMoveExpression14()
func (c *ExpressionDecomposerTest) TestMoveExpression14() {
	c.helperMoveExpression(
		"({a, b} = foo());",
		exprMatchesStr("foo()"),
		"var result$jscomp$0 = foo(); ({a, b} = result$jscomp$0);")
}

// Java: ExpressionDecomposerTest.testExposeExpression1()
func (c *ExpressionDecomposerTest) TestExposeExpression1() {
	c.helperExposeExpression(
		"x = 0 && foo()",
		exprMatchesStr("foo()"),
		"var temp$jscomp$0; if (temp$jscomp$0 = 0) temp$jscomp$0 = foo(); x = temp$jscomp$0;")
}

// Java: ExpressionDecomposerTest.testExposeExpression2()
func (c *ExpressionDecomposerTest) TestExposeExpression2() {
	c.helperExposeExpression(
		"x = 1 || foo()",
		exprMatchesStr("foo()"),
		"var temp$jscomp$0; if (temp$jscomp$0 = 1); else temp$jscomp$0=foo(); x = temp$jscomp$0;")
}

// Java: ExpressionDecomposerTest.testExposeExpression3()
func (c *ExpressionDecomposerTest) TestExposeExpression3() {
	c.helperExposeExpression(
		"var x = 1 ? foo() : 0",
		exprMatchesStr("foo()"),
		`var temp$jscomp$0;
if (1) temp$jscomp$0 = foo(); else temp$jscomp$0 = 0;
var x = temp$jscomp$0;
`)
}

// Java: ExpressionDecomposerTest.testExposeExpression4()
func (c *ExpressionDecomposerTest) TestExposeExpression4() {
	c.helperExposeExpression(
		"const x = 1 ? foo() : 0",
		exprMatchesStr("foo()"),
		`var temp$jscomp$0;
if (1) temp$jscomp$0 = foo(); else temp$jscomp$0 = 0;
const x = temp$jscomp$0;
`)
}

// Java: ExpressionDecomposerTest.testExposeExpression5()
func (c *ExpressionDecomposerTest) TestExposeExpression5() {
	c.helperExposeExpression(
		"let x = 1 ? foo() : 0",
		exprMatchesStr("foo()"),
		`var temp$jscomp$0;
if (1) temp$jscomp$0 = foo(); else temp$jscomp$0 = 0;
let x = temp$jscomp$0;
`)
}

// Java: ExpressionDecomposerTest.testExposeExpression6()
func (c *ExpressionDecomposerTest) TestExposeExpression6() {
	c.helperExposeExpression("goo() && foo()", exprMatchesStr("foo()"), "if (goo()) foo();")
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: ExpressionDecomposerTest.exposeExpressionNullishCoalesceNoResult()
func (c *ExpressionDecomposerTest) TestExposeExpressionNullishCoalesceNoResult() {
	c.helperExposeExpression(
		"goo() ?? foo()",
		exprMatchesStr("foo()"),
		`var temp$jscomp$1;
if((temp$jscomp$1 = goo()) != null) temp$jscomp$1;
else foo()
`)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: ExpressionDecomposerTest.exposeExpressionOptionalGetElem()
func (c *ExpressionDecomposerTest) TestExposeExpressionOptionalGetElem() {
	c.helperExposeExpression(
		"a = x?.[foo()]",
		exprMatchesStr("foo()"),
		`let temp_const$jscomp$0;
var temp$jscomp$1;
if ((temp_const$jscomp$0 = x) == null) {
  temp$jscomp$1 = void 0;
} else {
  temp$jscomp$1 = temp_const$jscomp$0[foo()];
}
a = temp$jscomp$1;
`)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: ExpressionDecomposerTest.exposeExpressionOptChainCallChain()
func (c *ExpressionDecomposerTest) TestExposeExpressionOptChainCallChain() {
	c.helperExposeExpression(
		"a = x?.(a).y.z[foo()]",
		exprMatchesStr("foo()"),
		`let temp_const$jscomp$0;
var temp$jscomp$1;
if ((temp_const$jscomp$0 = x) == null) {
  temp$jscomp$1 = void 0;
} else {
  var temp_const$jscomp$2 = temp_const$jscomp$0(a).y.z;
  temp$jscomp$1 = temp_const$jscomp$2[foo()];
}
a = temp$jscomp$1;
`)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: ExpressionDecomposerTest.exposeExpressionOptChainCallChainNoResult()
func (c *ExpressionDecomposerTest) TestExposeExpressionOptChainCallChainNoResult() {
	c.helperExposeExpression(
		"x?.(a)[y].z[foo()]",
		exprMatchesStr("foo()"),
		`let temp_const$jscomp$0;
if ((temp_const$jscomp$0 = x) == null) {
  void 0;
} else {
  var temp_const$jscomp$2 = temp_const$jscomp$0(a)[y].z;
  temp_const$jscomp$2[foo()];
}
`)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: ExpressionDecomposerTest.exposeExpressionOptionalGetPropChain()
func (c *ExpressionDecomposerTest) TestExposeExpressionOptionalGetPropChain() {
	c.helperExposeExpression(
		"a = x?.y.z[foo()]",
		exprMatchesStr("foo()"),
		`let temp_const$jscomp$0;
var temp$jscomp$1;
if ((temp_const$jscomp$0 = x) == null) {
  temp$jscomp$1 = void 0;
} else {
  var temp_const$jscomp$2 = temp_const$jscomp$0.y.z;
  temp$jscomp$1 = temp_const$jscomp$2[foo()];
}
a = temp$jscomp$1;
`)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: ExpressionDecomposerTest.exposeExpressionOptionalGetPropChainNoResult()
func (c *ExpressionDecomposerTest) TestExposeExpressionOptionalGetPropChainNoResult() {
	c.helperExposeExpression(
		"x?.y.z[foo()]",
		exprMatchesStr("foo()"),
		`let temp_const$jscomp$0;
if ((temp_const$jscomp$0 = x) == null) {
  void 0;
} else {
  var temp_const$jscomp$2 = temp_const$jscomp$0.y.z;
  temp_const$jscomp$2[foo()];
}
`)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: ExpressionDecomposerTest.exposeExpressionOptionalGetElemChain()
func (c *ExpressionDecomposerTest) TestExposeExpressionOptionalGetElemChain() {
	c.helperExposeExpression(
		"a = x?.[y].z[foo()];",
		exprMatchesStr("foo()"),
		`let temp_const$jscomp$0;
var temp$jscomp$1;
if ((temp_const$jscomp$0 = x) == null) {
  temp$jscomp$1 = void 0;
} else {
  var temp_const$jscomp$2 = temp_const$jscomp$0[y].z;
  temp$jscomp$1 = temp_const$jscomp$2[foo()];
}
a = temp$jscomp$1;
`)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: ExpressionDecomposerTest.exposeExpressionOptionalGetElemChainNoResult()
func (c *ExpressionDecomposerTest) TestExposeExpressionOptionalGetElemChainNoResult() {
	c.helperExposeExpression(
		"x?.[y].z[foo()]",
		exprMatchesStr("foo()"),
		`let temp_const$jscomp$0;
if ((temp_const$jscomp$0 = x) == null) {
  void 0;
} else {
  var temp_const$jscomp$2 = temp_const$jscomp$0[y].z;
  temp_const$jscomp$2[foo()];
}
`)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: ExpressionDecomposerTest.exposeExpressionOptionalGetElemWithCall()
func (c *ExpressionDecomposerTest) TestExposeExpressionOptionalGetElemWithCall() {
	c.helperExposeExpression(
		"a = x.y?.[z](foo())",
		exprMatchesStr("foo()"),
		`let temp_const$jscomp$0;
var temp$jscomp$1;
if ((temp_const$jscomp$0 = x.y) == null) {
  temp$jscomp$1 = void 0;
} else {
  var temp_const$jscomp$3 = temp_const$jscomp$0;
  var temp_const$jscomp$2 = temp_const$jscomp$3[z];
  temp$jscomp$1 = temp_const$jscomp$2.call(temp_const$jscomp$3, foo());
}
a = temp$jscomp$1;
`)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: ExpressionDecomposerTest.exposeExpressionGetElemWithOptChainCall()
func (c *ExpressionDecomposerTest) TestExposeExpressionGetElemWithOptChainCall() {
	c.helperExposeExpression(
		"a = x.y[z]?.(foo(), d)",
		exprMatchesStr("foo()"),
		`let temp_const$jscomp$0;
let temp_const$jscomp$1;
var temp$jscomp$2;
if ((temp_const$jscomp$1 = (temp_const$jscomp$0 = x.y)[z]) == null) {
  temp$jscomp$2 = void 0;
} else {
  temp$jscomp$2 = temp_const$jscomp$1.call(temp_const$jscomp$0, foo(), d);
}
a = temp$jscomp$2;
`)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: ExpressionDecomposerTest.exposeExpressionOptionalGetPropWithCall()
func (c *ExpressionDecomposerTest) TestExposeExpressionOptionalGetPropWithCall() {
	c.helperExposeExpression(
		"a = x.y?.z(foo(1))",
		exprMatchesStr("foo(1)"),
		`let temp_const$jscomp$0;
var temp$jscomp$1;
if ((temp_const$jscomp$0 = x.y) == null) {
  temp$jscomp$1 = void 0;
} else {
  var temp_const$jscomp$3 = temp_const$jscomp$0;
  var temp_const$jscomp$2 = temp_const$jscomp$3.z;
  temp$jscomp$1 = temp_const$jscomp$2.call(temp_const$jscomp$3, foo(1));
}
a = temp$jscomp$1;
`)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: ExpressionDecomposerTest.exposeExpressionOptionalGetPropWithCallTwiceRewriteCall()
func (c *ExpressionDecomposerTest) TestExposeExpressionOptionalGetPropWithCallTwiceRewriteCall() {
	c.helperExposeExpression(
		"a = x.y?.z(foo(1))",
		exprMatchesStr("foo(1)"),
		`let temp_const$jscomp$0;
var temp$jscomp$1;
if ((temp_const$jscomp$0 = x.y) == null) {
  temp$jscomp$1 = void 0;
} else {
  var temp_const$jscomp$3 = temp_const$jscomp$0;
  var temp_const$jscomp$2 = temp_const$jscomp$3.z;
  temp$jscomp$1 = temp_const$jscomp$2.call(temp_const$jscomp$3, foo(1));
}
a = temp$jscomp$1;
`)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: ExpressionDecomposerTest.exposeExpressionGetPropWithOptChainCall()
func (c *ExpressionDecomposerTest) TestExposeExpressionGetPropWithOptChainCall() {
	c.helperExposeExpression(
		"a = x.y.z?.(foo())",
		exprMatchesStr("foo()"),
		`let temp_const$jscomp$0;
let temp_const$jscomp$1;
var temp$jscomp$2;
if ((temp_const$jscomp$1 = (temp_const$jscomp$0 = x.y).z) == null) {
  temp$jscomp$2 = void 0;
} else {
  temp$jscomp$2 = temp_const$jscomp$1.call(temp_const$jscomp$0, foo());
}
a = temp$jscomp$2;
`)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: ExpressionDecomposerTest.exposeExpressionNewOptChainAfterRewriteCall()
func (c *ExpressionDecomposerTest) TestExposeExpressionNewOptChainAfterRewriteCall() {
	c.helperExposeExpression(
		"a = x?.y(foo())?.z.q",
		exprMatchesStr("foo()"),
		`let temp_const$jscomp$0;
let temp_const$jscomp$1;
var temp$jscomp$2;
if ((temp_const$jscomp$0 = x) == null) {
  temp$jscomp$2 = void 0;
} else {
  var temp_const$jscomp$4 = temp_const$jscomp$0;
  var temp_const$jscomp$3 = temp_const$jscomp$4.y;
  temp$jscomp$2 =
      (temp_const$jscomp$1 =
          temp_const$jscomp$3.call(temp_const$jscomp$4, foo())) == null
              ? void 0 : temp_const$jscomp$1.z.q;
}
a = temp$jscomp$2;
`)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: ExpressionDecomposerTest.exposeExpressionNewOptChainAfter()
func (c *ExpressionDecomposerTest) TestExposeExpressionNewOptChainAfter() {
	c.helperExposeExpression(
		"a = x?.y[foo()]?.z.q",
		exprMatchesStr("foo()"),
		`let temp_const$jscomp$0;
let temp_const$jscomp$1;
var temp$jscomp$2;
if ((temp_const$jscomp$0 = x) == null) {
  temp$jscomp$2 = void 0;
} else {
  var temp_const$jscomp$3 = temp_const$jscomp$0.y;
  temp$jscomp$2 = (temp_const$jscomp$1 = temp_const$jscomp$3[foo()]) == null ? void 0 : temp_const$jscomp$1.z.q;
}
a = temp$jscomp$2;
`)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: ExpressionDecomposerTest.exposeExpressionNotImmediatelyFollowedByNewChain()
func (c *ExpressionDecomposerTest) TestExposeExpressionNotImmediatelyFollowedByNewChain() {
	c.helperExposeExpression(
		"a = x?.y[foo()].z.q?.b.c",
		exprMatchesStr("foo()"),
		`let temp_const$jscomp$0;
let temp_const$jscomp$1;
var temp$jscomp$2;
if ((temp_const$jscomp$0 = x) == null) {
  temp$jscomp$2 = void 0;
} else {
  var temp_const$jscomp$3 = temp_const$jscomp$0.y;
  temp$jscomp$2 = (temp_const$jscomp$1 = temp_const$jscomp$3[foo()].z.q) == null
      ? void 0 : temp_const$jscomp$1.b.c;
}
a = temp$jscomp$2;
`)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: ExpressionDecomposerTest.exposeExpressionBreakingOutOfOptionalChain()
func (c *ExpressionDecomposerTest) TestExposeExpressionBreakingOutOfOptionalChain() {
	c.helperExposeExpression(
		"a = (x?.y[foo()]).z.q",
		exprMatchesStr("foo()"),
		`let temp_const$jscomp$0;
var temp$jscomp$1;
if ((temp_const$jscomp$0 = x) == null) {
  temp$jscomp$1 = void 0;
} else {
  var temp_const$jscomp$2 = temp_const$jscomp$0.y;
  temp$jscomp$1 = temp_const$jscomp$2[foo()];
}
a = temp$jscomp$1.z.q;
`)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: ExpressionDecomposerTest.exposeExpressionCallBreakingOutOfOptionalChain()
func (c *ExpressionDecomposerTest) TestExposeExpressionCallBreakingOutOfOptionalChain() {
	c.shouldTestTypes = false
	// Performing a non-optional call on an optional chain is not good coding, because you could
	// end up trying to call `undefined` as a function, but it is allowed and must be supported.
	c.helperExposeExpression(
		"a = (x?.y.z)(foo())",
		exprMatchesStr("foo()"),
		"let temp_const$jscomp$0;\nlet temp_const$jscomp$1;\nvar temp_const$jscomp$3 =\n    (temp_const$jscomp$0 = x) == null\n        ? void 0 : (temp_const$jscomp$1 = temp_const$jscomp$0.y).z;\n// This double .call is unfortunate but not incorrect.\n// Maybe we should make the ExpressionDecomposer recognize `.call` and avoid creating\n// another one? We'd run the risk of breaking \"real\" methods called `.call`, which\n// are allowed.\nvar temp_const$jscomp$2 = temp_const$jscomp$3.call;\n// The temp_const$jscomp$1 argument gets type '?' when we decompose,\n// which is correct, but TypeInference on this expected code appears to give it\n// `undefined`. This is why we've set `shouldTestTypes = false;` above\na = temp_const$jscomp$2.call(temp_const$jscomp$3, temp_const$jscomp$1, foo());\n")
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: ExpressionDecomposerTest.exposeExpressionFreeCallBreakingOutOfOptionalChain()
func (c *ExpressionDecomposerTest) TestExposeExpressionFreeCallBreakingOutOfOptionalChain() {
	// Performing a non-optional call on an optional chain is not good coding, because you could
	// end up trying to call `undefined` as a function, but it is allowed and must be supported.
	c.helperExposeExpression(
		"a = (x?.y.z())(foo())",
		exprMatchesStr("foo()"),
		`var temp_const$jscomp$0 = x?.y.z();
a = temp_const$jscomp$0(foo());
`)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: ExpressionDecomposerTest.exposeExpressionCallAtEndOfOptionalChain()
func (c *ExpressionDecomposerTest) TestExposeExpressionCallAtEndOfOptionalChain() {
	c.helperExposeExpression(
		"a = x?.y.z(foo())",
		exprMatchesStr("foo()"),
		`let temp_const$jscomp$0;
var temp$jscomp$1;
if ((temp_const$jscomp$0 = x) == null) {
  temp$jscomp$1 = void 0;
} else {
  var temp_const$jscomp$3 = temp_const$jscomp$0.y;
  var temp_const$jscomp$2 = temp_const$jscomp$3.z;
  temp$jscomp$1 = temp_const$jscomp$2.call(temp_const$jscomp$3, foo());
}
a = temp$jscomp$1;
`)
}

// Java: ExpressionDecomposerTest.testBug117935266_expose_call_target()
func (c *ExpressionDecomposerTest) TestBug117935266_expose_call_target() {
	c.helperExposeExpression(
		`function first() {
  alert('first');
  return '';
}
// alert must be preserved before the first side-effect
alert(first().method(alert('second')).method(alert('third')));
`,
		exprMatchesStr("first()"),
		`function first() {
  alert('first');
      return '';
}
var temp_const$jscomp$0 = alert;
temp_const$jscomp$0(first().method(
    alert('second')).method(alert('third')));
`)
}

// Java: ExpressionDecomposerTest.testBug117935266_move_call_target()
func (c *ExpressionDecomposerTest) TestBug117935266_move_call_target() {
	c.helperMoveExpression(
		`function first() {
  alert('first');
      return '';
}
var temp_const$jscomp$0 = alert;
temp_const$jscomp$0(first().toString(
    alert('second')).toString(alert('third')));
`,
		exprMatchesStr("first()"),
		`function first() {
  alert('first');
      return '';
}
var temp_const$jscomp$0 = alert;
var result$jscomp$0 = first();
temp_const$jscomp$0(result$jscomp$0.toString(
    alert('second')).toString(alert('third')));
`)
}

// Java: ExpressionDecomposerTest.testBug117935266_expose_call_parameters()
func (c *ExpressionDecomposerTest) TestBug117935266_expose_call_parameters() {
	c.helperExposeExpression(
		"alert(fn(first(), second(), third()));",
		exprMatchesStr("first()"),
		`var temp_const$jscomp$1 = alert;
var temp_const$jscomp$0 = fn;
temp_const$jscomp$1(temp_const$jscomp$0(first(), second(), third()));
`)

	c.helperExposeExpression(
		"alert(fn(first(), second(), third()));",
		exprMatchesStr("second()"),
		`var temp_const$jscomp$2 = alert;
var temp_const$jscomp$1 = fn;
var temp_const$jscomp$0 = first();
temp_const$jscomp$2(temp_const$jscomp$1(temp_const$jscomp$0, second(), third()));
`)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: ExpressionDecomposerTest.exposeExpressionAfterTwoOptionalChains()
func (c *ExpressionDecomposerTest) TestExposeExpressionAfterTwoOptionalChains() {
	c.helperExposeExpression(
		"a = x?.y.z?.q(foo());",
		exprMatchesStr("foo()"),
		`let temp_const$jscomp$0;
let temp_const$jscomp$1;
var temp$jscomp$2;
if ((temp_const$jscomp$0 = x) == null) {
  temp$jscomp$2 = void 0;
} else {
  var temp$jscomp$3;
  if ((temp_const$jscomp$1 = temp_const$jscomp$0.y.z) == null) {
    temp$jscomp$3 = void 0;
  } else {
    var temp_const$jscomp$5 = temp_const$jscomp$1;
    var temp_const$jscomp$4 = temp_const$jscomp$5.q;
    temp$jscomp$3 = temp_const$jscomp$4.call(temp_const$jscomp$5, foo());
  }
  temp$jscomp$2 = temp$jscomp$3;
}
a = temp$jscomp$2;
`)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: ExpressionDecomposerTest.exposeAnOptionalChain()
func (c *ExpressionDecomposerTest) TestExposeAnOptionalChain() {
	c.helperExposeExpression(
		"a = foo(arg1, opt?.chain())",
		exprMatchesStr("opt?.chain()"),
		`var temp_const$jscomp$1 = foo;
var temp_const$jscomp$0 = arg1;
a = temp_const$jscomp$1(temp_const$jscomp$0, opt?.chain());
`)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: ExpressionDecomposerTest.exposePartOfAnOptionalChain()
func (c *ExpressionDecomposerTest) TestExposePartOfAnOptionalChain() {
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE,
		// The non-optional part is fine to move.
		"nonOptional.part?.optional.chain?.continues()",
		exprMatchesStr("nonOptional.part"))

	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE,
		"nonOptional.part?.optional.chain?.continues()",
		// Cannot just move part of an optional chain.
		exprMatchesStr("nonOptional.part?.optional"))

	c.helperExposeExpression(
		"nonOptional.part?.optional.chain?.continues()",
		exprMatchesStr("nonOptional.part?.optional"),
		"let temp_const$jscomp$0;\nlet temp_const$jscomp$1;\nif ((temp_const$jscomp$0 = nonOptional.part) == null) {\n  void 0;\n} else {\n// `temp_const$jscomp$0.optional` is the expression we were trying to expose, and it\n// is now movable.\n  (temp_const$jscomp$1 = temp_const$jscomp$0.optional.chain) == null\n      ? void 0 : temp_const$jscomp$1.continues();\n}\n")
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: ExpressionDecomposerTest.canExposeMethodCallee()
func (c *ExpressionDecomposerTest) TestCanExposeMethodCallee() {
	// TODO(b/161802885): This should probably return DECOMPOSABLE, because it is not safe to
	// replace `foo.bar` with a temporary variable containing its value.
	c.helperCanExposeExpression(ExpressionDecomposer_DecompositionType_MOVABLE, "foo.bar()", exprMatchesStr("foo.bar"))

	// TODO(b/161802885): This should probably return DECOMPOSABLE, because it is not safe to
	// replace `foo?.bar` with a temporary variable containing its value.
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "(foo?.bar)()", exprMatchesStr("foo?.bar"))
}

// Java: ExpressionDecomposerTest.testExposeExpression7()
func (c *ExpressionDecomposerTest) TestExposeExpression7() {
	c.helperExposeExpression(
		"x = goo() && foo()",
		exprMatchesStr("foo()"),
		"var temp$jscomp$0; if (temp$jscomp$0 = goo()) temp$jscomp$0 = foo(); x = temp$jscomp$0;")
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: ExpressionDecomposerTest.exposeExpressionNullishCoalesce()
func (c *ExpressionDecomposerTest) TestExposeExpressionNullishCoalesce() {
	c.helperExposeExpression(
		"x = goo() ?? foo()",
		exprMatchesStr("foo()"),
		`var temp$jscomp$1;var temp$jscomp$0;
if((temp$jscomp$1 = goo()) != null) temp$jscomp$0 = temp$jscomp$1;
else temp$jscomp$0=foo(); x = temp$jscomp$0;
`)
}

// Java: ExpressionDecomposerTest.testExposeExpression8()
func (c *ExpressionDecomposerTest) TestExposeExpression8() {
	c.helperExposeExpression(
		"var x = 1 + (goo() && foo())",
		exprMatchesStr("foo()"),
		`var temp$jscomp$0;
if (temp$jscomp$0 = goo()) temp$jscomp$0 = foo();
var x = 1 + temp$jscomp$0;
`)
}

// Java: ExpressionDecomposerTest.testExposeExpression9()
func (c *ExpressionDecomposerTest) TestExposeExpression9() {
	c.helperExposeExpression(
		"const x = 1 + (goo() && foo())",
		exprMatchesStr("foo()"),
		`var temp$jscomp$0;
if (temp$jscomp$0 = goo()) temp$jscomp$0 = foo();
const x = 1 + temp$jscomp$0;
`)
}

// Java: ExpressionDecomposerTest.testExposeExpression10()
func (c *ExpressionDecomposerTest) TestExposeExpression10() {
	c.helperExposeExpression(
		"let x = 1 + (goo() && foo())",
		exprMatchesStr("foo()"),
		`var temp$jscomp$0;
if (temp$jscomp$0 = goo()) temp$jscomp$0 = foo();
let x = 1 + temp$jscomp$0;
`)
}

// Java: ExpressionDecomposerTest.testExposeExpression11()
func (c *ExpressionDecomposerTest) TestExposeExpression11() {
	c.helperExposeExpression(
		"if(goo() && foo());",
		exprMatchesStr("foo()"),
		`var temp$jscomp$0;
if (temp$jscomp$0 = goo()) temp$jscomp$0 = foo();
if(temp$jscomp$0);
`)
}

// Java: ExpressionDecomposerTest.testExposeExpression12()
func (c *ExpressionDecomposerTest) TestExposeExpression12() {
	c.helperExposeExpression(
		"switch(goo() && foo()){}",
		exprMatchesStr("foo()"),
		`var temp$jscomp$0;
if (temp$jscomp$0 = goo()) temp$jscomp$0 = foo();
switch(temp$jscomp$0){}
`)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: ExpressionDecomposerTest.exposeExpressionNullishCoalesceSwitch()
func (c *ExpressionDecomposerTest) TestExposeExpressionNullishCoalesceSwitch() {
	c.helperExposeExpression(
		"switch(goo() ?? foo()){}",
		exprMatchesStr("foo()"),
		`var temp$jscomp$1;var temp$jscomp$0;
if((temp$jscomp$1 = goo()) != null) temp$jscomp$0 = temp$jscomp$1;
else temp$jscomp$0 = foo(); switch(temp$jscomp$0){}
`)
}

// Java: ExpressionDecomposerTest.testExposeExpression13()
func (c *ExpressionDecomposerTest) TestExposeExpression13() {
	c.helperExposeExpression(
		"switch(1 + goo() + foo()){}",
		exprMatchesStr("foo()"),
		"var temp_const$jscomp$0 = 1 + goo(); switch(temp_const$jscomp$0 + foo()){}")
}

// Java: ExpressionDecomposerTest.testExposeExpression_inVanillaForInitializer_simpleExpression()
func (c *ExpressionDecomposerTest) TestExposeExpression_inVanillaForInitializer_simpleExpression() {
	c.helperExposeExpression(
		"for (x = goo() + foo();;) {}",
		exprMatchesStr("foo()"),
		`var temp_const$jscomp$0 = goo();
for (x = temp_const$jscomp$0 + foo();;) {}
`)
}

// Java: ExpressionDecomposerTest.testExposeExpression_inVanillaForInitializer_usingLabel()
func (c *ExpressionDecomposerTest) TestExposeExpression_inVanillaForInitializer_usingLabel() {
	c.helperExposeExpression(
		"LABEL: for (x = goo() + foo();;) {}",
		exprMatchesStr("foo()"),
		`var temp_const$jscomp$0 = goo();
LABEL: for (x = temp_const$jscomp$0 + foo();;) {}
`)
}

// Java: ExpressionDecomposerTest.testExposeExpression_inVanillaForInitializer_singleDeclaration_withLetOrConst()
func (c *ExpressionDecomposerTest) TestExposeExpression_inVanillaForInitializer_singleDeclaration_withLetOrConst() {
	for _, declarationKeyword := range []string{"let", "const"} {
		c.helperExposeExpression(
			"for ("+declarationKeyword+" x = goo() + foo();;) {}",
			exprMatchesStr("foo()"),
			strings.ReplaceAll(`var temp_const$jscomp$0 = goo();
for (DECLATION_KEYWORD x = temp_const$jscomp$0 + foo();;) {}
`, "DECLATION_KEYWORD", declarationKeyword))
	}
}

// Java: ExpressionDecomposerTest.testExposeExpression_inVanillaForInitializer_firstDeclaration_withLetOrConst()
func (c *ExpressionDecomposerTest) TestExposeExpression_inVanillaForInitializer_firstDeclaration_withLetOrConst() {
	for _, declarationKeyword := range []string{"let", "const"} {
		c.helperExposeExpression(
			"for ("+declarationKeyword+" x = goo() + foo(), y = 5;;) {}",
			exprMatchesStr("foo()"),
			strings.ReplaceAll(`var temp_const$jscomp$0 = goo();
for (DECLATION_KEYWORD x = temp_const$jscomp$0 + foo(), y = 5;;) {}
`, "DECLATION_KEYWORD", declarationKeyword))
	}
}

// Java: ExpressionDecomposerTest.testExposeExpression14()
func (c *ExpressionDecomposerTest) TestExposeExpression14() {
	c.helperExposeExpression(
		"function f(){ return goo() && foo();}",
		exprMatchesStr("foo()"),
		`function f() {
  var temp$jscomp$0; if (temp$jscomp$0 = goo()) temp$jscomp$0 = foo();
  return temp$jscomp$0;
}
`)
}

// Java: ExpressionDecomposerTest.testExposeExpression15()
func (c *ExpressionDecomposerTest) TestExposeExpression15() {
	// TODO(johnlenz): We really want a constant marking pass.
	// The value "goo" should be constant, but it isn't known to be so.
	c.helperExposeExpression(
		"if (goo(1, goo(2), (1 ? foo() : 0)));",
		exprMatchesStr("foo()"),
		`var temp_const$jscomp$1 = goo;
var temp_const$jscomp$0 = goo(2);
var temp$jscomp$2;
if (1) temp$jscomp$2 = foo(); else temp$jscomp$2 = 0;
if (temp_const$jscomp$1(1, temp_const$jscomp$0, temp$jscomp$2));
`)
}

// Java: ExpressionDecomposerTest.testExposeExpression16()
func (c *ExpressionDecomposerTest) TestExposeExpression16() {
	c.helperExposeExpression(
		"throw bar() && foo();",
		exprMatchesStr("foo()"),
		"var temp$jscomp$0; if (temp$jscomp$0 = bar()) temp$jscomp$0=foo(); throw temp$jscomp$0;")
}

// Java: ExpressionDecomposerTest.testExposeExpression17()
func (c *ExpressionDecomposerTest) TestExposeExpression17() {
	c.helperExposeExpression(
		"x.foo(y())",
		exprMatchesStr("y()"),
		`var temp_const$jscomp$1 = x;
var temp_const$jscomp$0 = temp_const$jscomp$1.foo;
temp_const$jscomp$0.call(temp_const$jscomp$1, y());
`)
}

// Java: ExpressionDecomposerTest.testExposeFreeCall()
func (c *ExpressionDecomposerTest) TestExposeFreeCall() {
	c.helperExposeExpression(
		"(0,x.foo)(y())",
		exprMatchesStr("y()"),
		`var temp_const$jscomp$0 = x.foo;
temp_const$jscomp$0(y());
`)
}

// Java: ExpressionDecomposerTest.testExposeTemplateLiteralFreeCall()
func (c *ExpressionDecomposerTest) TestExposeTemplateLiteralFreeCall() {
	c.helperExposeExpression(
		"foo`${x()}${y()}`",
		exprMatchesStr("y()"),
		"var temp_const$jscomp$1 = foo;\nvar temp_const$jscomp$0 = x();\ntemp_const$jscomp$1`${temp_const$jscomp$0}${y()}`;\n")
}

// Java: ExpressionDecomposerTest.testCanExposeTaggedTemplateLiteralInterpolation()
func (c *ExpressionDecomposerTest) TestCanExposeTaggedTemplateLiteralInterpolation() {
	c.helperCanExposeExpression(ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "x`${y()}`", exprMatchesStr("y()"))
	// TODO(b/251958225): Implement decomposition for this case.
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE, "x.foo`${y()}`", exprMatchesStr("y()"))
}

// Java: ExpressionDecomposerTest.testExposeExpression18()
func (c *ExpressionDecomposerTest) TestExposeExpression18() {
	c.helperExposeExpression(
		`const {a, b, c} = condition ?
  y() :
  {a: 0, b: 0, c: 1};
`,
		exprMatchesStr("y()"),
		`var temp$jscomp$0;
if (condition) {
  temp$jscomp$0 = y();
} else {
  temp$jscomp$0 = {a: 0, b: 0, c: 1};
}
const {a, b, c} = temp$jscomp$0;
`)
}

// Java: ExpressionDecomposerTest.testMoveClass1()
func (c *ExpressionDecomposerTest) TestMoveClass1() {
	c.shouldTestTypes = false
	// types don't come out quite the same before and after decomposition
	// TODO(bradfordcsmith): See TODO in helperMoveExpression()
	c.helperMoveExpression(
		"alert(class X {});",
		func(compiler AbstractCompiler) func(*rhino.Node) *rhino.Node {
			return expressionDecomposerTest_findClass
		},
		"var result$jscomp$0 = class X {}; alert(result$jscomp$0);")
}

// Java: ExpressionDecomposerTest.testMoveClass2()
func (c *ExpressionDecomposerTest) TestMoveClass2() {
	c.shouldTestTypes = false
	// types don't come out quite the same before and after decomposition
	// TODO(bradfordcsmith): See TODO in helperMoveExpression()
	c.helperMoveExpression(
		"console.log(1, 2, class X {});",
		func(compiler AbstractCompiler) func(*rhino.Node) *rhino.Node {
			return expressionDecomposerTest_findClass
		},
		"var result$jscomp$0 = class X {}; console.log(1, 2, result$jscomp$0);")
}

// Java: ExpressionDecomposerTest.testMoveYieldExpression1()
func (c *ExpressionDecomposerTest) TestMoveYieldExpression1() {
	c.helperMoveExpression(
		"function *f() { return { a: yield 1, c: foo(yield 2, yield 3) }; }",
		exprMatchesStr("yield 1"),
		`function *f() {
  var result$jscomp$0 = yield 1;
  return { a: result$jscomp$0, c: foo(yield 2, yield 3) };
}
`)

	c.helperMoveExpression(
		"function *f() { return { a: 0, c: foo(yield 2, yield 3) }; }",
		exprMatchesStr("yield 2"),
		`function *f() {
  var result$jscomp$0 = yield 2;
  return { a: 0, c: foo(result$jscomp$0, yield 3) };
}
`)

	c.helperMoveExpression(
		"function *f() { return { a: 0, c: foo(1, yield 3) }; }",
		exprMatchesStr("yield 3"),
		`function *f() {
  var result$jscomp$0 = yield 3;
  return { a: 0, c: foo(1, result$jscomp$0) };
}
`)
}

// Java: ExpressionDecomposerTest.testMoveYieldExpression2()
func (c *ExpressionDecomposerTest) TestMoveYieldExpression2() {
	c.helperMoveExpression(
		"function *f() { return (yield 1) || (yield 2); }",
		exprMatchesStr("yield 1"),
		`function *f() {
  var result$jscomp$0 = yield 1;
  return result$jscomp$0 || (yield 2);
}
`)
}

// Java: ExpressionDecomposerTest.testMoveYieldExpression3()
func (c *ExpressionDecomposerTest) TestMoveYieldExpression3() {
	c.helperMoveExpression(
		"function *f() { return x.y(yield 1); }",
		exprMatchesStr("yield 1"),
		`function *f() {
  var result$jscomp$0 = yield 1;
  return x.y(result$jscomp$0);
}
`)
}

// Java: ExpressionDecomposerTest.testExposeYieldExpression1()
func (c *ExpressionDecomposerTest) TestExposeYieldExpression1() {
	c.helperExposeExpression(
		"function *f(x) { return x || (yield 2); }",
		exprMatchesStr("yield 2"),
		`function *f(x) {
  var temp$jscomp$0;
  if (temp$jscomp$0=x); else temp$jscomp$0 = yield 2;
  return temp$jscomp$0
}
`)
}

// Java: ExpressionDecomposerTest.testExposeYieldExpression2()
func (c *ExpressionDecomposerTest) TestExposeYieldExpression2() {
	c.helperExposeExpression(
		"function *f() { return x.y(yield 1); }",
		exprMatchesStr("yield 1"),
		`function *f() {
  var temp_const$jscomp$1 = x;
  var temp_const$jscomp$0 = temp_const$jscomp$1.y;
  return temp_const$jscomp$0.call(temp_const$jscomp$1, yield 1);
}
`)
}

// Java: ExpressionDecomposerTest.testExposeYieldExpression3()
func (c *ExpressionDecomposerTest) TestExposeYieldExpression3() {
	c.helperExposeExpression(
		"function *f() { return g.call(yield 1); }",
		exprMatchesStr("yield 1"),
		`function *f() {
  var temp_const$jscomp$1 = g;
  var temp_const$jscomp$0 = temp_const$jscomp$1.call;
  return temp_const$jscomp$0.call(temp_const$jscomp$1, yield 1);
}
`)
}

// Java: ExpressionDecomposerTest.testExposeYieldExpression4()
func (c *ExpressionDecomposerTest) TestExposeYieldExpression4() {
	c.helperExposeExpression(
		"function *f() { return g.apply([yield 1, yield 2]); }",
		exprMatchesStr("yield 1"),
		`function *f() {
  var temp_const$jscomp$1 = g;
  var temp_const$jscomp$0 = temp_const$jscomp$1.apply;
  return temp_const$jscomp$0.call(temp_const$jscomp$1, [yield 1, yield 2]);
}
`)
}

// Java: ExpressionDecomposerTest.testExposePlusEquals1()
func (c *ExpressionDecomposerTest) TestExposePlusEquals1() {
	c.helperExposeExpression(
		"var x = 0; x += foo() + 1",
		exprMatchesStr("foo()"),
		"var x = 0; var temp_const$jscomp$0 = x; x = temp_const$jscomp$0 + (foo() + 1);")

	c.helperExposeExpression(
		"var x = 0; y = (x += foo()) + x",
		exprMatchesStr("foo()"),
		"var x = 0; var temp_const$jscomp$0 = x; y = (x = temp_const$jscomp$0 + foo()) + x")
}

// Java: ExpressionDecomposerTest.testExposePlusEquals2()
func (c *ExpressionDecomposerTest) TestExposePlusEquals2() {
	c.helperExposeExpression(
		"var x = {}; x.a += foo() + 1",
		exprMatchesStr("foo()"),
		`var x = {}; var temp_const$jscomp$0 = x;
var temp_const$jscomp$1 = temp_const$jscomp$0.a;
temp_const$jscomp$0.a = temp_const$jscomp$1 + (foo() + 1);
`)

	c.helperExposeExpression(
		"var x = {}; y = (x.a += foo()) + x.a",
		exprMatchesStr("foo()"),
		`var x = {}; var temp_const$jscomp$0 = x;
var temp_const$jscomp$1 = temp_const$jscomp$0.a;
y = (temp_const$jscomp$0.a = temp_const$jscomp$1 + foo()) + x.a
`)
}

// Java: ExpressionDecomposerTest.testExposePlusEquals3()
func (c *ExpressionDecomposerTest) TestExposePlusEquals3() {
	c.helperExposeExpression(
		"/** @const */ var XX = {}; XX.a += foo() + 1",
		exprMatchesStr("foo()"),
		`var XX = {};
var temp_const$jscomp$0 = XX.a;
XX.a = temp_const$jscomp$0 + (foo() + 1);
`)

	c.helperExposeExpression(
		"var XX = {}; y = (XX.a += foo()) + XX.a",
		exprMatchesStr("foo()"),
		`var XX = {};
var temp_const$jscomp$0 = XX.a;
y = (XX.a = temp_const$jscomp$0 + foo()) + XX.a
`)
}

// Java: ExpressionDecomposerTest.testExposePlusEquals4()
func (c *ExpressionDecomposerTest) TestExposePlusEquals4() {
	c.helperExposeExpression(
		"var x = {}; goo().a += foo() + 1",
		exprMatchesStr("foo()"),
		`var x = {};
var temp_const$jscomp$0 = goo();
var temp_const$jscomp$1 = temp_const$jscomp$0.a;
temp_const$jscomp$0.a = temp_const$jscomp$1 + (foo() + 1);
`)

	c.helperExposeExpression(
		"var x = {}; y = (goo().a += foo()) + goo().a",
		exprMatchesStr("foo()"),
		`var x = {};
var temp_const$jscomp$0 = goo();
var temp_const$jscomp$1 = temp_const$jscomp$0.a;
y = (temp_const$jscomp$0.a = temp_const$jscomp$1 + foo()) + goo().a
`)
}

// Java: ExpressionDecomposerTest.testExposePlusEquals5()
func (c *ExpressionDecomposerTest) TestExposePlusEquals5() {
	c.helperExposeExpression(
		"var x = {}; goo().a.b += foo() + 1",
		exprMatchesStr("foo()"),
		`var x = {};
var temp_const$jscomp$0 = goo().a;
var temp_const$jscomp$1 = temp_const$jscomp$0.b;
temp_const$jscomp$0.b = temp_const$jscomp$1 + (foo() + 1);
`)

	c.helperExposeExpression(
		"var x = {}; y = (goo().a.b += foo()) + goo().a",
		exprMatchesStr("foo()"),
		`var x = {};
var temp_const$jscomp$0 = goo().a;
var temp_const$jscomp$1 = temp_const$jscomp$0.b;
y = (temp_const$jscomp$0.b = temp_const$jscomp$1 + foo()) + goo().a
`)
}

// Dynamic receiver and key on LHS of assignment-op.
// Java: ExpressionDecomposerTest.testExposePlusEquals6()
func (c *ExpressionDecomposerTest) TestExposePlusEquals6() {
	c.helperExposeExpression(
		"var obj = {}; var nextKey = function() {}; obj[nextKey()] += foo() + 1",
		exprMatchesStr("foo()"),
		`var obj = {};
var nextKey = function() {};
var temp_const$jscomp$0 = obj;
var temp_const$jscomp$1 = nextKey();
var temp_const$jscomp$2 = temp_const$jscomp$0[temp_const$jscomp$1];
temp_const$jscomp$0[temp_const$jscomp$1] = temp_const$jscomp$2 + (foo() + 1);
`)

	c.helperExposeExpression(
		"var obj = {}; var nextKey = function() {}; y = (obj[nextKey()] += foo()) + obj[nextKey()]",
		exprMatchesStr("foo()"),
		`var obj = {};
var nextKey = function() {};
var temp_const$jscomp$0 = obj;
var temp_const$jscomp$1 = nextKey();
var temp_const$jscomp$2 = temp_const$jscomp$0[temp_const$jscomp$1];
y = (temp_const$jscomp$0[temp_const$jscomp$1] = temp_const$jscomp$2 + foo()) + obj[nextKey()]
`)
}

// Java: ExpressionDecomposerTest.testExposePlusEquals_getElemCallReceiverAndKey()
func (c *ExpressionDecomposerTest) TestExposePlusEquals_getElemCallReceiverAndKey() {
	c.helperExposeExpression(
		"var x = {}; goo()[goo()] += foo() + 1",
		exprMatchesStr("foo()"),
		`var x = {};
var temp_const$jscomp$0 = goo();
var temp_const$jscomp$1 = goo();
var temp_const$jscomp$2 = temp_const$jscomp$0[temp_const$jscomp$1];
temp_const$jscomp$0[temp_const$jscomp$1] = temp_const$jscomp$2 + (foo() + 1);
`)
}

// Java: ExpressionDecomposerTest.testExposeLogicalAssignment1()
func (c *ExpressionDecomposerTest) TestExposeLogicalAssignment1() {
	// Part of the work here is being done by Normalize, which converts all
	// instances of logical assignment operators into an larger expression
	// that separates the logical operation from the assignment.
	c.helperExposeExpression(
		"let x = 0; x ||= foo() + 1",
		exprMatchesStr("foo()"),
		`let x = 0;
if (x) {
} else {
   x = foo() + 1;
}
`)

	c.helperExposeExpression(
		"let x = 0; x &&= foo() + 1",
		exprMatchesStr("foo()"),
		`let x = 0;
if (x) {
   x = foo() + 1;
}
`)

	c.helperExposeExpression(
		"let x = 0; x ??= foo() + 1",
		exprMatchesStr("foo()"),
		`let x = 0;
var temp$jscomp$1;
if ((temp$jscomp$1 = x) != null) {
   temp$jscomp$1;
} else {
   x = foo() + 1;
}
`)
}

// Java: ExpressionDecomposerTest.testExposeLogicalAssignment2()
func (c *ExpressionDecomposerTest) TestExposeLogicalAssignment2() {
	// Part of the work here is being done by Normalize, which converts all
	// instances of logical assignment operators into an larger expression
	// that separates the logical operation from the assignment.
	c.helperExposeExpression(
		"let x = {}; x.a ||= foo() + 1",
		exprMatchesStr("foo()"),
		`let x = {};
let $jscomp$logical$assign$tmpm1146332801$0;
if (($jscomp$logical$assign$tmpm1146332801$0 = x).a) {
} else {
   var temp_const$jscomp$1 = $jscomp$logical$assign$tmpm1146332801$0;
   temp_const$jscomp$1.a = foo() + 1;
}
`)
	c.helperExposeExpression(
		"let x = {}; x[a] &&= foo() + 1",
		exprMatchesStr("foo()"),
		`let x = {};
let $jscomp$logical$assign$tmpm1146332801$0;
let $jscomp$logical$assign$tmpindexm1146332801$0;
if (($jscomp$logical$assign$tmpm1146332801$0 = x)
    [$jscomp$logical$assign$tmpindexm1146332801$0 = a]) {
    var temp_const$jscomp$2 = $jscomp$logical$assign$tmpm1146332801$0;
    var temp_const$jscomp$1 = $jscomp$logical$assign$tmpindexm1146332801$0;
    temp_const$jscomp$2[temp_const$jscomp$1] = foo() + 1;
}
`)
}

// Java: ExpressionDecomposerTest.testExposeObjectLit1()
func (c *ExpressionDecomposerTest) TestExposeObjectLit1() {
	// Validate that getter and setters methods are seen as side-effect
	// free and that values can move past them.  We don't need to be
	// concerned with exposing the getter or setter here but the
	// decomposer does not have a method of exposing properties, only variables.
	c.helperMoveExpression(
		"var x = {get a() {}, b: foo()};",
		exprMatchesStr("foo()"),
		"var result$jscomp$0=foo();var x = {get a() {}, b: result$jscomp$0};")

	c.helperMoveExpression(
		"var x = {set a(p) {}, b: foo()};",
		exprMatchesStr("foo()"),
		"var result$jscomp$0=foo();var x = {set a(p) {}, b: result$jscomp$0};")
}

// Java: ExpressionDecomposerTest.testMoveSpread_siblingOfCall_outOfArrayLiteral_usesTempArray()
func (c *ExpressionDecomposerTest) TestMoveSpread_siblingOfCall_outOfArrayLiteral_usesTempArray() {
	c.shouldTestTypes = false
	// types don't come out quite the same before and after decomposition
	// TODO(bradfordcsmith): See TODO in helperMoveExpression()
	c.helperExposeExpression(
		"[...x, foo()];",
		exprMatchesStr("foo()"),
		`var temp_const$jscomp$0 = [...x];
[...temp_const$jscomp$0, foo()];
`)
}

// Java: ExpressionDecomposerTest.testMoveSpread_siblingOfCall_outOfObjectLiteral_usesTempObject()
func (c *ExpressionDecomposerTest) TestMoveSpread_siblingOfCall_outOfObjectLiteral_usesTempObject() {
	c.shouldTestTypes = false
	// types don't come out quite the same before and after decomposition
	// TODO(bradfordcsmith): See TODO in helperMoveExpression()
	c.helperExposeExpression(
		"({...x, y: foo()});",
		exprMatchesStr("foo()"),
		`var temp_const$jscomp$0 = {...x};
({...temp_const$jscomp$0, y: foo()});
`)
}

// Java: ExpressionDecomposerTest.testMoveSpread_siblingOfCall_outOfFunctionCall_usesTempArray()
func (c *ExpressionDecomposerTest) TestMoveSpread_siblingOfCall_outOfFunctionCall_usesTempArray() {
	c.shouldTestTypes = false
	// types don't come out quite the same before and after decomposition
	// TODO(bradfordcsmith): See TODO in helperMoveExpression()
	c.helperExposeExpression(
		`function f() { }
f(...x, foo());
`,
		exprMatchesStr("foo()"),
		`function f() { }
var temp_const$jscomp$1 = f;
var temp_const$jscomp$0 = [...x];
temp_const$jscomp$1(...temp_const$jscomp$0, foo());
`)
}

// Java: ExpressionDecomposerTest.testMoveSpreadParent_siblingOfCall_outOfFunctionCall_usesNoTempArray()
func (c *ExpressionDecomposerTest) TestMoveSpreadParent_siblingOfCall_outOfFunctionCall_usesNoTempArray() {
	c.helperExposeExpression(
		`function f() { }
f([...x], foo());
`,
		exprMatchesStr("foo()"),
		`function f() { }
var temp_const$jscomp$1 = f;
var temp_const$jscomp$0 = [...x];
temp_const$jscomp$1(temp_const$jscomp$0, foo());
`)
}

// Java: ExpressionDecomposerTest.testMoveSpreadParent_siblingOfCall_outOfFunctionCall_usesNoTempObject()
func (c *ExpressionDecomposerTest) TestMoveSpreadParent_siblingOfCall_outOfFunctionCall_usesNoTempObject() {
	c.helperExposeExpression(
		`function f() { }
f({...x}, foo());
`,
		exprMatchesStr("foo()"),
		`function f() { }
var temp_const$jscomp$1 = f;
var temp_const$jscomp$0 = {...x};
temp_const$jscomp$1(temp_const$jscomp$0, foo());
`)
}

// Java: ExpressionDecomposerTest.testExposeExpressionInTemplateLibSub()
func (c *ExpressionDecomposerTest) TestExposeExpressionInTemplateLibSub() {
	c.helperExposeExpression(
		"` ${ foo() }  ${ goo() } `;",
		exprMatchesStr("goo()"),
		"var temp_const$jscomp$0 = foo(); ` ${ temp_const$jscomp$0 }  ${ goo() } `;")
}

// Java: ExpressionDecomposerTest.testExposeSubExpressionInTemplateLibSub()
func (c *ExpressionDecomposerTest) TestExposeSubExpressionInTemplateLibSub() {
	c.helperExposeExpression(
		"` ${ foo() + goo() } `;",
		exprMatchesStr("goo()"),
		"var temp_const$jscomp$0 = foo(); ` ${ temp_const$jscomp$0 + goo() } `;")
}

// Java: ExpressionDecomposerTest.testMoveExpressionInTemplateLibSub()
func (c *ExpressionDecomposerTest) TestMoveExpressionInTemplateLibSub() {
	c.helperMoveExpression(
		"` ${ foo() }  ${ goo() } `;",
		exprMatchesStr("foo()"),
		"var result$jscomp$0 = foo(); ` ${ result$jscomp$0 }  ${ goo() } `;")
}

// Java: ExpressionDecomposerTest.testExposeExpression_computedProp_withPureKey()
func (c *ExpressionDecomposerTest) TestExposeExpression_computedProp_withPureKey() {
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE,
		`({
  ['a' + 'b']: foo(),
});
`,
		exprMatchesStr("foo()"))
}

// Java: ExpressionDecomposerTest.testExposeObjectLitValue_computedProp_withImpureKey()
func (c *ExpressionDecomposerTest) TestExposeObjectLitValue_computedProp_withImpureKey() {
	c.helperExposeExpression(
		`({
  [goo()]: foo(),
});
`,
		exprMatchesStr("foo()"),
		`var temp_const$jscomp$0 = goo();
({
  [temp_const$jscomp$0]: foo(),
});
`)
}

// Java: ExpressionDecomposerTest.testExposeObjectLitValue_computedProp_asEarlierSibling_withImpureKeyAndValue()
func (c *ExpressionDecomposerTest) TestExposeObjectLitValue_computedProp_asEarlierSibling_withImpureKeyAndValue() {
	c.helperExposeExpression(
		`({
  [goo()]: qux(),
  bar: foo(),
});
`,
		exprMatchesStr("foo()"),
		`var temp_const$jscomp$1 = goo();
var temp_const$jscomp$0 = qux();
({
  [temp_const$jscomp$1]: temp_const$jscomp$0,
  bar: foo(),
});
`)
}

// Java: ExpressionDecomposerTest.testExposeObjectLitValue_memberFunctions_asEarlierSiblings_arePure()
func (c *ExpressionDecomposerTest) TestExposeObjectLitValue_memberFunctions_asEarlierSiblings_arePure() {
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE,
		`({
  a() { },
  get b() { },
  set b(v) { },

  bar: foo(),
});
`,
		exprMatchesStr("foo()"))
}

// Java: ExpressionDecomposerTest.testMoveSuperCall()
func (c *ExpressionDecomposerTest) TestMoveSuperCall() {
	c.helperMoveExpression(
		"class A { constructor() { super(foo()) } }",
		exprMatchesStr("foo()"),
		"class A{constructor(){var result$jscomp$0=foo();super(result$jscomp$0)}}")
}

// Java: ExpressionDecomposerTest.testMoveSuperCall_noSideEffects()
func (c *ExpressionDecomposerTest) TestMoveSuperCall_noSideEffects() {
	// String() is being used since it's known to not have side-effects.
	c.helperMoveExpression(
		"class A { constructor() { super(String()) } }",
		exprMatchesStr("String()"),
		"class A{constructor(){var result$jscomp$0=String();super(result$jscomp$0)}}")
}

// Java: ExpressionDecomposerTest.testExposeSuperCall()
func (c *ExpressionDecomposerTest) TestExposeSuperCall() {
	c.helperExposeExpression(
		"class A { constructor() { super(goo(), foo()) } }",
		exprMatchesStr("foo()"),
		`class A{ constructor(){
   var temp_const$jscomp$0=goo();
   super(temp_const$jscomp$0, foo())
}}
`)
}

// Java: ExpressionDecomposerTest.testExposeSuperCall_noSideEffects()
func (c *ExpressionDecomposerTest) TestExposeSuperCall_noSideEffects() {
	// String() is being used since it's known to not have side-effects.
	c.helperExposeExpression(
		"class A { constructor() { super(goo(), String()) } }",
		exprMatchesStr("String()"),
		`class A{ constructor(){
   var temp_const$jscomp$0=goo();
   super(temp_const$jscomp$0, String())
}}
`)
}

// Java: ExpressionDecomposerTest.testCannotDecomposeSuperMethodCall()
func (c *ExpressionDecomposerTest) TestCannotDecomposeSuperMethodCall() {
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE,
		"class A extends B { fn() { super.method(foo()) } }",
		exprMatchesStr("foo()"))
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE,
		"class A extends B { fn() { super['method'](foo()) } }",
		exprMatchesStr("foo()"))
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: ExpressionDecomposerTest.canMovePastFnDotCall()
func (c *ExpressionDecomposerTest) TestCanMovePastFnDotCall() {
	c.knownConstants.Add("fn")
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_MOVABLE, "fn.call(foo());", exprMatchesStr("foo()"))
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "unknownIfFn.call(foo());", exprMatchesStr("foo()"))
	c.helperMoveExpression(
		"fn.call(foo());",
		exprMatchesStr("foo()"),
		"var result$jscomp$0 = foo(); fn.call(result$jscomp$0);")
}

// Java: ExpressionDecomposerTest.testDecomposeEs6Features()
func (c *ExpressionDecomposerTest) TestDecomposeEs6Features() {
	c.knownConstants.Add("tag")
	c.knownConstants.Add("sink")
	// Tagged template literals
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "sink(tag`x`, foo())", exprMatchesStr("foo()"))

	// Await expressions
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE,
		"async function f() { sink(await 0, foo()); }",
		exprMatchesStr("foo()"))

	// Dynamic imports
	c.helperCanExposeExpression(
		ExpressionDecomposer_DecompositionType_DECOMPOSABLE, "sink(import('m'), foo())", exprMatchesStr("foo()"))
}

// Java: ExpressionDecomposerTest.helperCanExposeExpression(DecompositionType, String, Function<AbstractCompiler, Function<Node, Node>>)
func (c *ExpressionDecomposerTest) helperCanExposeExpression(
	expectedResult ExpressionDecomposer_DecompositionType,
	code string,
	nodeFinderFn expressionDecomposerTest_nodeFinder) {
	compiler := c.getCompiler()
	decomposer := c.createDecomposer(compiler)
	tree := expressionDecomposerTest_parse(compiler, code)
	AssertThat(tree).IsNotNull()

	externsRoot := expressionDecomposerTest_parse(compiler, "function goo() {} function foo() {}")
	AssertThat(externsRoot).IsNotNull()

	expresionNode := nodeFinderFn(compiler)(tree)

	result := decomposer.CanExposeExpression(expresionNode)
	AssertThat(result).IsEqualTo(expectedResult)
}

// Provides a toString() method that contains the source code a node
// represents where possible, or some explanatory text when not possible.
//
// Java: ExpressionDecomposerTest.NodeToSource
type expressionDecomposerTest_NodeToSource struct {
	compiler AbstractCompiler
	node     *rhino.Node
}

// Java: ExpressionDecomposerTest.NodeToSource(AbstractCompiler, Node)
func newExpressionDecomposerTest_NodeToSource(compiler AbstractCompiler, node *rhino.Node) *expressionDecomposerTest_NodeToSource {
	return &expressionDecomposerTest_NodeToSource{compiler: compiler, node: node}
}

// Java: ExpressionDecomposerTest.NodeToSource.toString()
func (n *expressionDecomposerTest_NodeToSource) String() string {
	if n.node.IsTemplateLitString() {
		// A string part of a template literal cannot be printed as code on its own.
		return fmt.Sprintf("[Template literal string: '%s']", n.node.GetRawString())
	} else if n.node.IsTemplateLitSub() {
		// The template literal substitution node cannot itself be turned into source code,
		// but we can do that for the expression inside of it.
		return fmt.Sprintf(
			"[Template literal substitution: '%s']", n.compiler.ToSourceOfNode(n.node.GetOnlyChild()))
	} else {
		return n.compiler.ToSourceOfNode(n.node)
	}
}

// Java: ExpressionDecomposerTest.nodeToSource(AbstractCompiler, Node)
func expressionDecomposerTest_nodeToSource(compiler AbstractCompiler, node *rhino.Node) string {
	return newExpressionDecomposerTest_NodeToSource(compiler, node).String()
}

// Java: ExpressionDecomposerTest.exprMatchesStr(String)
func exprMatchesStr(exprString string) expressionDecomposerTest_nodeFinder {
	return func(compiler AbstractCompiler) func(root *rhino.Node) *rhino.Node {
		return func(root *rhino.Node) *rhino.Node {
			// When matching trim off the trailing newline added by the compiler's pretty-print
			// option.
			isAMatch := func(node *rhino.Node) bool {
				return exprString == jutil.StringTrim(expressionDecomposerTest_nodeToSource(compiler, node))
			}
			containsAMatch := func(node *rhino.Node) bool {
				return strings.Contains(expressionDecomposerTest_nodeToSource(compiler, node), exprString)
			}
			matchingNode := NodeUtil_FindPreorder(root, isAMatch, containsAMatch)
			AssertWithMessage(
				"Expected node `%s` was not found in `%s`", exprString, compiler.ToSourceOfNode(root)).
				That(matchingNode).
				IsNotNull()
			return matchingNode
		}
	}
}

// Java: ExpressionDecomposerTest.helperExposeExpressionThenTypeCheck(String, Function<Node, Node>)
func (c *ExpressionDecomposerTest) helperExposeExpressionThenTypeCheck(
	code string, nodeFinder func(*rhino.Node) *rhino.Node) *rhino.Node {
	compiler := c.getCompiler()
	tree := expressionDecomposerTest_parse(compiler, code)
	AssertThat(tree).IsNotNull()

	decomposer := c.createDecomposer(compiler)
	decomposer.SetTempNamePrefix("temp")
	decomposer.SetResultNamePrefix("result")

	expr := nodeFinder(tree)

	decomposer.MaybeExposeExpression(expr)
	c.processForTypecheck(compiler, tree)

	return tree
}

// Java: ExpressionDecomposerTest.helperExposeExpression(String, Function<AbstractCompiler, Function<Node, Node>>, String)
func (c *ExpressionDecomposerTest) helperExposeExpression(
	code string, compilerToNodeFinder expressionDecomposerTest_nodeFinder, expectedResult string) {
	compiler := c.getCompiler()

	expectedRoot := expressionDecomposerTest_parse(compiler, expectedResult)
	tree := expressionDecomposerTest_parse(compiler, code)
	AssertThat(tree).IsNotNull()

	if c.shouldTestTypes {
		c.processForTypecheck(compiler, tree)
	}

	decomposer := c.createDecomposer(compiler)
	decomposer.SetTempNamePrefix("temp")
	decomposer.SetResultNamePrefix("result")

	nodeFinder := compilerToNodeFinder(compiler)
	expr := nodeFinder(tree)
	AssertWithMessage("Expected node was not found.").That(expr).IsNotNull()

	result := decomposer.CanExposeExpression(expr)
	AssertThat(result).IsEqualTo(ExpressionDecomposer_DecompositionType_DECOMPOSABLE)

	decomposer.MaybeExposeExpression(expr)
	c.validateSourceInfo(compiler, tree)
	AssertNode(tree).UsingSerializer(compiler.ToSourceOfNode).IsEqualTo(expectedRoot)

	if c.shouldTestTypes {
		decomposeThenTypeCheck := c.helperExposeExpressionThenTypeCheck(code, nodeFinder)
		c.checkTypeStringsEqualAsTree(decomposeThenTypeCheck, tree)
	}
}

// Java: ExpressionDecomposerTest.createDecomposer(Compiler)
func (c *ExpressionDecomposerTest) createDecomposer(compiler *Compiler) *ExpressionDecomposer {
	return compiler.CreateExpressionDecomposer(
		compiler.GetUniqueNameIdSupplier(), jutil.NewLinkedHashSet(c.knownConstants.ToSlice()...), c.newScope())
}

// Java: ExpressionDecomposerTest.helperMoveExpression(String, Function<AbstractCompiler, Function<Node, Node>>, String)
func (c *ExpressionDecomposerTest) helperMoveExpression(
	code string, compilerToNodeFinder expressionDecomposerTest_nodeFinder, expectedResult string) {
	compiler := c.getCompiler()
	nodeFinder := compilerToNodeFinder(compiler)

	expectedRoot := expressionDecomposerTest_parse(compiler, expectedResult)
	tree := expressionDecomposerTest_parse(compiler, code)
	originalTree := tree.CloneTree()
	AssertThat(tree).IsNotNull()

	if c.shouldTestTypes {
		c.processForTypecheck(compiler, tree)
	}

	decomposer := c.createDecomposer(compiler)
	decomposer.SetTempNamePrefix("temp")
	decomposer.SetResultNamePrefix("result")

	expr := nodeFinder(tree)
	AssertWithMessage("Expected node was not found.").That(expr).IsNotNull()

	decomposer.MoveExpression(expr)
	c.validateSourceInfo(compiler, tree)
	AssertNode(tree).UsingSerializer(compiler.ToSourceOfNode).IsEqualTo(expectedRoot)

	if c.shouldTestTypes {
		// find a basis for comparison:
		originalExpr := nodeFinder(originalTree)

		decomposer.MoveExpression(originalExpr)
		c.processForTypecheck(compiler, originalTree)

		// TODO(bradfordcsmith): Don't assume type check + decompose gives the same results as
		// decompose + type check.
		// There are legitimate cases where the types will be different from one order to another,
		// but not actually wrong.
		c.checkTypeStringsEqualAsTree(originalTree, tree)
	}
}

// Java: ExpressionDecomposerTest.checkTypeStringsEqualAsTree(Node, Node)
func (c *ExpressionDecomposerTest) checkTypeStringsEqualAsTree(rootExpected *rhino.Node, rootActual *rhino.Node) {
	expectedType := rootExpected.GetJSType()
	actualType := rootActual.GetJSType()

	if expectedType == nil || actualType == nil {
		AssertWithMessage("Expected %s but got %s", rootExpected, rootActual).
			That(actualType).
			IsEqualTo(expectedType)
	} else if expectedType.IsUnknownType() && actualType.IsUnknownType() {
		// continue
	} else {
		// we can't compare actual equality because the types are from different runs of the
		// type inference, so we just compare the strings.
		AssertWithMessage("Expected %s but got %s", rootExpected, rootActual).
			That(actualType.ToAnnotationString(rhino.JSType_Nullability_EXPLICIT)).
			IsEqualTo(expectedType.ToAnnotationString(rhino.JSType_Nullability_EXPLICIT))
	}

	child1 := rootExpected.GetFirstChild()
	child2 := rootActual.GetFirstChild()
	for child1 != nil {
		c.checkTypeStringsEqualAsTree(child1, child2)
		child1 = child1.GetNext()
		child2 = child2.GetNext()
	}
}

// Java: ExpressionDecomposerTest.getCompiler()
func (c *ExpressionDecomposerTest) getCompiler() *Compiler {
	compiler := NewCompiler()
	options := NewCompilerOptions()
	// If the specific test case requested an output language level,
	// use it. Otherwise, keep the default.
	if c.languageOut != nil {
		options.SetLanguageOut(*c.languageOut)
	}
	options.SetCodingConvention(NewGoogleCodingConvention())
	options.SetPrettyPrint(true)
	// Don't prefix the compiler output with `"use strict";`.
	// It's noise for these tests, and it interferes with tests that
	// want to use compiler.toSource() to string match expressions.
	options.SetEmitUseStrict(false)
	options.SetLanguageIn(CompilerOptions_LanguageMode_UNSUPPORTED)
	compiler.InitOptions(options)
	return compiler
}

// Java: ExpressionDecomposerTest.processForTypecheck(AbstractCompiler, Node)
func (c *ExpressionDecomposerTest) processForTypecheck(compiler AbstractCompiler, jsRoot *rhino.Node) {
	root := rhino.IR_Root(rhino.IR_Root(), rhino.IR_Root(jsRoot))
	registry := compiler.GetTypeRegistry()
	NewTypeCheck(compiler, NewSemanticReverseAbstractInterpreter(registry), registry).
		ProcessForTesting(root.GetFirstChild(), root.GetSecondChild())
	compiler.SetTypeCheckingHasRun(true)
}

// Java: ExpressionDecomposerTest.findClass(Node)
func expressionDecomposerTest_findClass(n *rhino.Node) *rhino.Node {
	if n.IsClass() {
		return n
	}
	for child := n.GetFirstChild(); child != nil; child = child.GetNext() {
		maybeClass := expressionDecomposerTest_findClass(child)
		if maybeClass != nil {
			return maybeClass
		}
	}
	return nil
}

// Java: ExpressionDecomposerTest.validateSourceInfo(Compiler, Node)
func (c *ExpressionDecomposerTest) validateSourceInfo(compiler *Compiler, subtree *rhino.Node) {
	NewSourceInfoCheck(compiler).SetCheckSubTree(subtree)
	// Source information problems are reported as compiler errors.
	AssertThat(compiler.GetErrors()).IsEmpty()
}

// Java: ExpressionDecomposerTest.parse(Compiler, String)
func expressionDecomposerTest_parse(compiler *Compiler, js string) *rhino.Node {
	n := compiler.ParseTestCode(js)
	NodeTraversal_Traverse(compiler, n, NewNormalize_NormalizeStatements(compiler, false, nil))
	AssertThat(compiler.GetErrors()).IsEmpty()
	return n
}

// Java: ExpressionDecomposerTest.newScope()
func (c *ExpressionDecomposerTest) newScope() *Scope {
	return Scope_CreateGlobalScope(rhino.NewNode(rhino.Token_ROOT))
}
