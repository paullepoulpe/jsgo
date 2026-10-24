// Port of com/google/javascript/jscomp/InlineVariablesTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/InlineVariablesTest.java)

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
	"strconv"
	"strings"
	"testing"
)

// InlineVariablesTest verifies that valid candidates for inlining are
// inlined, but that no dangerous inlining occurs.
type InlineVariablesTest struct {
	CompilerTestCaseBase

	inlineLocalsOnly bool
}

// Java: InlineVariablesTest()
func NewInlineVariablesTest() *InlineVariablesTest {
	c := &InlineVariablesTest{}
	c.InitCompilerTestCaseBase(c)
	c.inlineLocalsOnly = false
	return c
}

func TestInlineVariablesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewInlineVariablesTest() })
}

// Java: InlineVariablesTest.setUp()
func (c *InlineVariablesTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableNormalize()
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.DisableCompareJsDoc()
	// NOTE: We are not enabling var checks here, so it is OK to use undeclared variables in
	// these tests. They will be treated as if they were externs.
}

// Java: InlineVariablesTest.getProcessor(Compiler)
func (c *InlineVariablesTest) GetProcessor(compiler *Compiler) CompilerPass {
	mode := InlineVariables_Mode_ALL
	if c.inlineLocalsOnly {
		mode = InlineVariables_Mode_LOCALS_ONLY
	}
	return NewInlineVariables(compiler, mode)
}

// Java: InlineVariablesTest.tearDown()
func (c *InlineVariablesTest) TearDown() {
	c.inlineLocalsOnly = false
}

// Java: InlineVariablesTest.testArgAlias()
func (c *InlineVariablesTest) TestArgAlias() {
	// same scope
	c.Test(
		`function foo(arg) {
  const argAlias = arg;
  const argAliasAlias = argAlias
  use(argAliasAlias);
}
`,
		`function foo(arg) {
  use(arg);
}
`)
	// nested scope
	c.Test(
		`function foo(arg) {
  const argAlias = arg;
  {
    const argAliasAlias = argAlias
    use(argAliasAlias);
  }
}
`,
		`function foo(arg) {
  {
    use(arg);
  }
}
`)
}

// Java: InlineVariablesTest.testInlineExpressionUsingAlias()
func (c *InlineVariablesTest) TestInlineExpressionUsingAlias() {
	c.Test(
		"function foo(arg1, arg2) {\n"+
			"  const arg1Alias = arg1;\n"+
			"  const arg2Alias = arg2;\n"+
			"// Create an inner scope to force analysis of `product` before the alias\n"+
			"// variables.\n"+
			"  {\n"+
			"    const product = arg1Alias * arg2Alias;\n"+
			"// Use `product` twice so if it gets inlined, it will have to create new\n"+
			"// references to the operands.\n"+
			"// Either the expression should not be inlined, or some effort will have\n"+
			"// to be made to avoid creating new references to `arg1Alias` and `arg2Alias`\n"+
			"// before they are considered for inlining themselves.\n"+
			"    use(product);\n"+
			"    use(product);\n"+
			"  }\n"+
			"}\n",
		`function foo(arg1, arg2) {
  {
    const product = arg1 * arg2;
    use(product);
    use(product);
  }
}
`)
}

// Java: InlineVariablesTest.testPassDoesntProduceInvalidCode1()
func (c *InlineVariablesTest) TestPassDoesntProduceInvalidCode1() {
	c.TestSame(
		`function f(x = void 0) {
  var z;
  {
    const y = {};
    x && (y['x'] = x);
    z = y;
  }
  return z;
}
`)
}

// Java: InlineVariablesTest.testPassDoesntProduceInvalidCode2()
func (c *InlineVariablesTest) TestPassDoesntProduceInvalidCode2() {
	c.TestSame(
		`function f(x = void 0) {
  {
    var z;
    const y = {};
    x && (y['x'] = x);
    z = y;
  }
  return z;
}
`)
}

// Java: InlineVariablesTest.testPassDoesntProduceInvalidCode3()
func (c *InlineVariablesTest) TestPassDoesntProduceInvalidCode3() {
	c.Test(
		`function f(x = void 0) {
  var z;
  const y = {};
  x && (y['x'] = x);
  z = y;
  {
    return z;
  }
}
`,
		`function f(x = void 0) {

  const y = {};
  x && (y['x'] = x);
  y;
  {
    return y;
  }
}
`)
}

// Test respect for scopes and blocks
// Java: InlineVariablesTest.testInlineGlobal()
func (c *InlineVariablesTest) TestInlineGlobal() {
	c.Test("var x = 1; var z = x;", "var z = 1;")
}

// Java: InlineVariablesTest.testNoInlineAnnotation()
func (c *InlineVariablesTest) TestNoInlineAnnotation() {
	c.TestSame("/** @noinline */ var x = 1; use(x);")
	c.Test(
		"/** @noinline */ var x = 1; use(x); var z = x; use(z);", //
		"/** @noinline */ var x = 1; use(x);            use(x);")
}

// Java: InlineVariablesTest.testNoInlineExportedName()
func (c *InlineVariablesTest) TestNoInlineExportedName() {
	c.TestSame("var _x = 1; var z = _x;")
}

// Java: InlineVariablesTest.testNoInlineExportedName2()
func (c *InlineVariablesTest) TestNoInlineExportedName2() {
	c.TestSame(
		`var f = function() {};
var _x = f;
var y = function() { _x(); };
var _y = f;
`)
}

// Java: InlineVariablesTest.testDontTreatLocalVariablesAsExportedByConvention()
func (c *InlineVariablesTest) TestDontTreatLocalVariablesAsExportedByConvention() {
	// In the GoogleCodingConvention, globals starting with "_" are exported. (see the above two
	// tests). Verify that we don't accidentally apply the same convention to locals.
	c.Test("function f() { var _x = 1; var z = _x; }", "function f() { var z = 1; }")
}

// Java: InlineVariablesTest.testDoNotInlineIncrement()
func (c *InlineVariablesTest) TestDoNotInlineIncrement() {
	c.TestSame("var x = 1; x++;")
}

// Java: InlineVariablesTest.testDoNotInlineDecrement()
func (c *InlineVariablesTest) TestDoNotInlineDecrement() {
	c.TestSame("var x = 1; x--;")
}

// Java: InlineVariablesTest.testDoNotInlineIntoLhsOfAssign()
func (c *InlineVariablesTest) TestDoNotInlineIntoLhsOfAssign() {
	c.TestSame("var x = 1; x += 3;")
}

// Java: InlineVariablesTest.testInlineIntoRhsOfAssign()
func (c *InlineVariablesTest) TestInlineIntoRhsOfAssign() {
	c.Test("var x = 1; var y = x;", "var y = 1;")
}

// Java: InlineVariablesTest.testInlineInFunction1()
func (c *InlineVariablesTest) TestInlineInFunction1() {
	c.Test(
		"function baz() { var x = 1; var z = x; }", //
		"function baz() { var z = 1;            }")
}

// Java: InlineVariablesTest.testInlineInFunction2()
func (c *InlineVariablesTest) TestInlineInFunction2() {
	c.Test(
		`function baz() {
  var a = new obj();
  result = a;
}
`,
		`function baz() {
  result = new obj();
}
`)
}

// Java: InlineVariablesTest.testInlineInFunction3()
func (c *InlineVariablesTest) TestInlineInFunction3() {
	c.TestSame(
		`function baz() {
  var a = new obj();
  (function(){a;})();
  result = a;
}
`)
}

// Java: InlineVariablesTest.testInlineInFunction4()
func (c *InlineVariablesTest) TestInlineInFunction4() {
	c.TestSame(
		`function baz() {
  var a = new obj();
  foo.result = a;
}
`)
}

// Java: InlineVariablesTest.testInlineInFunction5()
func (c *InlineVariablesTest) TestInlineInFunction5() {
	c.TestSame(
		`function baz() {
var a = (foo = new obj());
foo.x();
result = a;
}
`)
}

// Java: InlineVariablesTest.testInlineInFunction6()
func (c *InlineVariablesTest) TestInlineInFunction6() {
	c.Test(
		"function baz() { { var x = 1; var z = x; } }", //
		"function baz() { { var z = 1;            } }")
}

// Java: InlineVariablesTest.testInlineInFunction7()
func (c *InlineVariablesTest) TestInlineInFunction7() {
	c.Test(
		"function baz() { var x = 1; { var z = x; } }", //
		"function baz() {            { var z = 1; } }")
}

// Java: InlineVariablesTest.testInlineIntoArrowFunction1()
func (c *InlineVariablesTest) TestInlineIntoArrowFunction1() {
	c.Test(
		"var x = 0; var f = () => x + 1;", //
		"           var f = () => 0 + 1;")
}

// Java: InlineVariablesTest.testInlineIntoArrowFunction2()
func (c *InlineVariablesTest) TestInlineIntoArrowFunction2() {
	c.Test(
		"var x = 0; var f = () => { return x + 1; }", //
		"           var f = () => { return 0 + 1; }")
}

// Java: InlineVariablesTest.testInlineAcrossChunks()
func (c *InlineVariablesTest) TestInlineAcrossChunks() {
	// TODO(kushal): Make decision about overlap with CrossChunkCodeMotion
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("var a = 2;").
				AddChunk("var b = a;").
				Build()),
		// Note that there's nothing ensuring "var a = 2;" will be loaded before "var b = a;", but
		// if it's not, this code would be broken uncompiled and we'd get an "a is not defined"
		// error, so the compiler inlines anyway.
		c.Expected("", "var b = 2;"))
}

// Java: InlineVariablesTest.testInliningAcrossSiblingChunk_assignmentChunkIsFirst()
func (c *InlineVariablesTest) TestInliningAcrossSiblingChunk_assignmentChunkIsFirst() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForTree().
				AddChunk("var a;"). // root chunk
				// Two child chunks that both depend on the root chunk. Neither is guaranteed
				// to be loaded, and if both load the order is not guaranteed.
				AddChunk("a = true;").
				AddChunk("function foo() { if (a) console.log('test'); } foo();").
				Build()),
		// The expected code has the correct behavior. The compiler doesn't know whether `a = true;`
		// will execute before the second child chunk or not.
		c.Expected(
			"var a;",
			// root
			"a = true;", // child 1
			"(function() { if (a) console.log('test'); })();"))
	// child 2
}

// Java: InlineVariablesTest.testInliningAcrossSiblingChunk_useChunkIsFirst()
func (c *InlineVariablesTest) TestInliningAcrossSiblingChunk_useChunkIsFirst() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForTree().
				AddChunk("var a;"). // root chunk
				// Two child chunks that both depend on the root chunk. Neither is guaranteed
				// to be loaded, and if both load the order is not guaranteed.
				AddChunk("function foo() { if (a) console.log('test'); } foo();").
				AddChunk("a = true;").
				Build()),
		// the expected code has the correct behavior. The compiler doesn't know whether `a = true;`
		// will execute before the first child chunk or not.
		c.Expected(
			"var a;", // root
			"(function() { if (a) console.log('test'); })();", // child 1
			"a = true;"))
	// child 2
}

// Java: InlineVariablesTest.testDoNotExitConditional1()
func (c *InlineVariablesTest) TestDoNotExitConditional1() {
	c.Test(
		"if (true) { var x = 1; } var z = x; use(z);", //
		"if (true) { var x = 1; }            use(x);")
}

// Java: InlineVariablesTest.testDoNotExitConditional2()
func (c *InlineVariablesTest) TestDoNotExitConditional2() {
	c.Test(
		"if (true) var x = 1; var z = x; use(z);", //
		"if (true) var x = 1;            use(x);")
}

// Java: InlineVariablesTest.testDoNotExitConditional3()
func (c *InlineVariablesTest) TestDoNotExitConditional3() {
	c.Test(
		"var x; if (true) x=1; var z = x; use(z);", //
		"var x; if (true) x=1;            use(x);")
}

// Java: InlineVariablesTest.testDoNotExitLoop()
func (c *InlineVariablesTest) TestDoNotExitLoop() {
	c.Test(
		"while (z) { var x = 3; } var y = x; use(y);", //
		"while (z) { var x = 3; }            use(x);")
}

// Java: InlineVariablesTest.testDoNotExitForLoop()
func (c *InlineVariablesTest) TestDoNotExitForLoop() {
	c.Test(
		"for (var i = 1; false; false) var z = i;", //
		"for (         ; false; false) var z = 1;")
	c.TestSame("for (; false; false) var i = 1; var z = i;")
	c.TestSame("for (var i in {}); var z = i;")
}

// Java: InlineVariablesTest.testConstInLoopAssignedOnceInLifetime()
func (c *InlineVariablesTest) TestConstInLoopAssignedOnceInLifetime() {
	c.Test(
		`for (; z; ) {
  const x = 1;
  use(x);
}
`,
		`for (; z; ) {
  use(1);
}
`)
	c.Test(
		`function foo(arg) {
  for (; z; ) {
    const alias = arg;
    use(alias);
    use2(alias);
  }
}
`,
		`function foo(arg) {
  for (; z; ) {
    use(arg);
    use2(arg);
  }
}
`)
}

// Java: InlineVariablesTest.testInlineSubscopeAlias()
func (c *InlineVariablesTest) TestInlineSubscopeAlias() {
	c.Test(
		`var x = function() {
  var self = this;
  return function() {
    var y = self;
    use(y);
  }
}
`,
		`var x = function() {
  var self = this;
  return function() {

    use(self);
  }
}
`)
	c.Test(
		`var x = function() {
  var y = [1];
  return function() {
    var z = y;
    use(z);
  };
}
`,
		`var x = function() {
  var y = [1];
  return function() {

    use(y);
  };
}
`)
	c.Test(
		`var x = function() {
  var y = 1;
  return function() {
    var z = y;
    use(z);
  };
}
`,
		`var x = function() {

  return function() {

    use(1);
  };
}
`)
}

// Java: InlineVariablesTest.testDoNotExitTry()
func (c *InlineVariablesTest) TestDoNotExitTry() {
	c.TestSame("try { var x = y; } catch (e) {} var z = y; ")
	c.TestSame("try { throw e; var x = 1; } catch (e) {} var z = x; ")
}

// Java: InlineVariablesTest.testDoNotEnterCatch()
func (c *InlineVariablesTest) TestDoNotEnterCatch() {
	c.TestSame("try { } catch (e) { var z = e; } ")
}

// Java: InlineVariablesTest.testDoNotEnterFinally()
func (c *InlineVariablesTest) TestDoNotEnterFinally() {
	c.Test(
		`try { throw e; var x = 1; } catch (e) {}
finally  { var z = x; use(z); }
`,
		`try { throw e; var x = 1; } catch (e) {}
finally  {            use(x); }
`)
}

// Java: InlineVariablesTest.testInsideIfConditional()
func (c *InlineVariablesTest) TestInsideIfConditional() {
	c.Test("var a = foo(); if (a) { alert(3); }", "if (foo()) { alert(3); }")
	c.Test("var a; a = foo(); if (a) { alert(3); }", "if (foo()) { alert(3); }")
}

// Java: InlineVariablesTest.testOnlyReadAtInitialization()
func (c *InlineVariablesTest) TestOnlyReadAtInitialization() {
	c.Test("var a; a = foo();", "foo();")
	c.Test("var a; if (a = foo()) { alert(3); }", "if (foo()) { alert(3); }")
	c.Test("var a; switch (a = foo()) {}", "switch(foo()) {}")
	c.Test(
		"var a; function f(){ return a = foo(); }", //
		"       function f(){ return     foo(); }")
	c.Test(
		"function f(){ var a; return a = foo(); }", //
		"function f(){        return     foo(); }")
	c.Test("var a; with (a = foo()) { alert(3); }", "with (foo()) { alert(3); }")
	c.Test("var a; b = (a = foo());", "b = foo();")
	c.Test(
		"var a; while(a = foo()) { alert(3); }", //
		"       while(    foo()) { alert(3); }")
	c.Test(
		"var a; for(;a = foo();) { alert(3); }", //
		"       for(;    foo();) { alert(3); }")
	c.Test(
		"var a; do {} while(a = foo()) { alert(3); }", //
		"       do {} while(    foo()) { alert(3); }")
}

// Java: InlineVariablesTest.testImmutableWithSingleReferenceAfterInitialzation()
func (c *InlineVariablesTest) TestImmutableWithSingleReferenceAfterInitialzation() {
	c.Test("var a; a = 1;", "1;")
	c.Test("var a; if (a = 1) { alert(3); }", "if (1) { alert(3); }")
	c.Test("var a; switch (a = 1) {}", "switch(1) {}")
	c.Test(
		"var a; function f(){ return a = 1; }", //
		"       function f(){ return 1; }")
	c.Test(
		"function f(){ var a; return a = 1; }", //
		"function f(){        return     1; }")
	c.Test("var a; with (a = 1) { alert(3); }", "with (1) { alert(3); }")
	c.Test("var a; b = (a = 1);", "b = 1;")
	c.Test(
		"var a; while(a = 1) { alert(3); }", //
		"       while(    1) { alert(3); }")
	c.Test(
		"var a; for(;a = 1;) { alert(3); }", //
		"       for(;    1;) { alert(3); }")
	c.Test(
		"var a; do {} while(a = 1) { alert(3); }", //
		"       do {} while(    1) { alert(3); }")
}

// Java: InlineVariablesTest.testSingleReferenceAfterInitialzation()
func (c *InlineVariablesTest) TestSingleReferenceAfterInitialzation() {
	c.Test("var a; a = foo();a;", "foo();")
	c.TestSame("var a; if (a = foo()) { alert(3); } a;")
	c.TestSame("var a; switch (a = foo()) {} a;")
	c.TestSame("var a; function f(){ return a = foo(); } a;")
	c.TestSame("function f(){ var a; return a = foo(); a;}")
	c.TestSame("var a; with (a = foo()) { alert(3); } a;")
	c.TestSame("var a; b = (a = foo()); a;")
	c.TestSame("var a; while(a = foo()) { alert(3); } a;")
	c.TestSame("var a; for(;a = foo();) { alert(3); } a;")
	c.TestSame("var a; do {} while(a = foo()) { alert(3); } a;")
}

// Java: InlineVariablesTest.testInsideIfBranch()
func (c *InlineVariablesTest) TestInsideIfBranch() {
	c.TestSame("var a = foo(); if (1) { alert(a); }")
}

// Java: InlineVariablesTest.testInsideAndConditional()
func (c *InlineVariablesTest) TestInsideAndConditional() {
	c.Test("var a = foo(); a && alert(3);", "foo() && alert(3);")
}

// Java: InlineVariablesTest.testInsideAndBranch()
func (c *InlineVariablesTest) TestInsideAndBranch() {
	c.TestSame("var a = foo(); 1 && alert(a);")
}

// Java: InlineVariablesTest.testInsideOrBranch()
func (c *InlineVariablesTest) TestInsideOrBranch() {
	c.TestSame("var a = foo(); 1 || alert(a);")
}

// Java: InlineVariablesTest.testInsideHookBranch()
func (c *InlineVariablesTest) TestInsideHookBranch() {
	c.TestSame("var a = foo(); 1 ? alert(a) : alert(3)")
}

// Java: InlineVariablesTest.testInsideHookConditional()
func (c *InlineVariablesTest) TestInsideHookConditional() {
	c.Test(
		"var a = foo(); a ? alert(1) : alert(3)", //
		"        foo()    ? alert(1) : alert(3)")
}

// Java: InlineVariablesTest.testInsideOrBranchInsideIfConditional()
func (c *InlineVariablesTest) TestInsideOrBranchInsideIfConditional() {
	c.TestSame("var a = foo(); if (x || a) {}")
}

// Java: InlineVariablesTest.testInsideOrBranchInsideIfConditionalWithConstant()
func (c *InlineVariablesTest) TestInsideOrBranchInsideIfConditionalWithConstant() {
	// We don't inline non-immutable constants into branches.
	c.TestSame("var a = [false]; if (x || a) {}")
}

// Java: InlineVariablesTest.testCrossFunctionsAsLeftLeaves()
func (c *InlineVariablesTest) TestCrossFunctionsAsLeftLeaves() {
	// Ensures getNext() understands how to walk past a function leaf
	c.Test(
		c.Srcs("var x = function() {};", "", "function cow() {} var z = x;"),
		c.Expected("", "", "function cow() {} var z = function() {};"))
	c.Test(
		c.Srcs("var x = function() {};", "", "var cow = function() {}; var z = x;"),
		c.Expected("", "", "var cow = function() {}; var z = function() {};"))
	c.Test(
		"var x = a; (function() { a++; })(); var z = x; use(z);",
		"var x = a; (function() { a++; })();            use(x);")
	c.Test(
		`var x = a;
function cow(){ a++; }
cow();
var z = x;
`,
		"var x = a;\n"+
			" // declaration removed\n"+
			"(function(){ a++; })();\n"+
			" // unused `z` removed\n")
	c.Test(
		"var x = a; cow(); var z = x; z; function cow() { a++; };",
		"var x = a; cow();            x; function cow() { a++; };")
}

// Test movement of constant values
// Java: InlineVariablesTest.testDoCrossFunction()
func (c *InlineVariablesTest) TestDoCrossFunction() {
	c.Test(
		"var x = 1; foo(); var z = x; use(z);", //
		"           foo();            use(1);")
}

// Java: InlineVariablesTest.testDoNotCrossReferencingFunction()
func (c *InlineVariablesTest) TestDoNotCrossReferencingFunction() {
	c.Test(
		`var f = function() { var z = x; };
var x = 1;
f();
var z = x;
use(z)
f();
`,
		`var f = function() { var z$jscomp$1 = x; };
var x = 1;
f();

use(x)
f();
`)
}

// Test tricky declarations and references
// Java: InlineVariablesTest.testChainedAssignment()
func (c *InlineVariablesTest) TestChainedAssignment() {
	c.Test("var a = 2, b = 2; var c = b;", "var a = 2; var c = 2;")
	c.Test("var a = 2, b = 2; var c = a;", "var b = 2; var c = 2;")
	c.Test("var a = b = 2; var f = 3; var c = a;", "var f = 3; var c = b = 2;")
	c.TestSame("var a = b = 2; var c = b;")
}

// Java: InlineVariablesTest.testForIn()
func (c *InlineVariablesTest) TestForIn() {
	c.TestSame("for (var i in j) { var c = i; }")
	c.TestSame("var i = 0; for (i in j) ;")
	c.TestSame("var i = 0; for (i in j) { var c = i; }")
	c.TestSame("i = 0; for (var i in j) { var c = i; }")
	c.TestSame("var j = {'key':'value'}; for (var i in j) {print(i)};")
}

// Test movement of values that have (may) side effects
// Java: InlineVariablesTest.testDoCrossNewVariables()
func (c *InlineVariablesTest) TestDoCrossNewVariables() {
	c.Test("var x = foo(); var z = x;", "var z = foo();")
}

// Java: InlineVariablesTest.testAliasAfterFunctionCall()
func (c *InlineVariablesTest) TestAliasAfterFunctionCall() {
	c.Test(
		"var x = foo(); bar(); var z = x; use(z);", //
		"var x = foo(); bar();            use(x);")
}

// Test movement of values that are complex but lack side effects
// Java: InlineVariablesTest.testDoNotCrossAssignment()
func (c *InlineVariablesTest) TestDoNotCrossAssignment() {
	c.Test(
		"var x = {}; var y = x.a; x.a = 1; var z = y; use(z);", //
		"var x = {}; var y = x.a; x.a = 1;            use(y);")
	c.TestSame("var a = this.id; foo(this.id = 3, a);")
}

// Java: InlineVariablesTest.testAliasAfterDelete()
func (c *InlineVariablesTest) TestAliasAfterDelete() {
	c.Test(
		"var x = {}; var y = x.a; delete x.a; var z = y; use(z);", //
		"var x = {}; var y = x.a; delete x.a;            use(y);")
}

// Java: InlineVariablesTest.testInlineAliasButNotSnapshot()
func (c *InlineVariablesTest) TestInlineAliasButNotSnapshot() {
	c.Test(
		"var b = 1; var a = b; b += 2; var c = a; use(c);", //
		"var b = 1; var a = b; b += 2;            use(a);")
}

// Java: InlineVariablesTest.testAliasAfterIncrement()
func (c *InlineVariablesTest) TestAliasAfterIncrement() {
	c.Test(
		"var a = b.c; b.c++; var d = a; use(d);", //
		"var a = b.c; b.c++;            use(a);")
}

// Java: InlineVariablesTest.testConstructorBeforeAlias()
func (c *InlineVariablesTest) TestConstructorBeforeAlias() {
	c.Test(
		"var a = b; new Foo(); var c = a; use(c);", //
		"var a = b; new Foo();            use(a);")
}

// Java: InlineVariablesTest.testDoNotInlineAliasAssignedBeforeAliasedVar()
func (c *InlineVariablesTest) TestDoNotInlineAliasAssignedBeforeAliasedVar() {
	// TODO(bradfordcsmith): This is consistent with past behavior, but it would be good to fix it.
	// Assumes we do not rely on undefined variables (not technically correct!)
	c.Test(
		"var a = b; var b = 3; alert(a)", //
		"           var b = 3; alert(b);")
	c.Test(
		`f();
function f() {
  var alias = original;
  alert(alias)
}
var original = 3;
`,
		`f();
function f() {

  alert(original)
}
var original = 3;
`)
	c.Test(
		`f();
var original = 3;
function f() {
  var alias = original;
  alert(alias)
}
`,
		`f();
var original = 3;
function f() {

  alert(original)
}
`)
}

// Java: InlineVariablesTest.testOverlappingInlines()
func (c *InlineVariablesTest) TestOverlappingInlines() {
	source := `a = function(el, x, opt_y) {
  var cur = bar(el);
  opt_y = x.y;
  x = x.x;
  var dx = x - cur.x;
  var dy = opt_y - cur.y;
  foo(
      el,
      el.offsetLeft + dx,
      el.offsetTop + dy);
};
`
	expected := `a = function(el, x, opt_y) {
  var cur = bar(el);
  opt_y = x.y;
  x = x.x;
  foo(
      el,
      el.offsetLeft + (x - cur.x),
      el.offsetTop + (opt_y - cur.y));
};
`
	c.Test(source, expected)
}

// Java: InlineVariablesTest.testOverlappingInlineFunctions()
func (c *InlineVariablesTest) TestOverlappingInlineFunctions() {
	source := `a = function() {
  var b = function(args) {var n;};
  var c = function(args) {};
  d(b,c);
};
`
	expected := `a = function() {
  d(function(args){var n;}, function(args){});
};
`
	c.Test(source, expected)
}

// Java: InlineVariablesTest.testInlineIntoLoops()
func (c *InlineVariablesTest) TestInlineIntoLoops() {
	c.Test(
		"var x = true; while (true) alert(   x);", //
		"              while (true) alert(true);")
	c.Test(
		"var x = true; while (true) for (var i in {}) alert(   x);",
		"              while (true) for (var i in {}) alert(true);")
	c.TestSame("var x = [true]; while (true) alert(x);")
}

// Java: InlineVariablesTest.testInlineIntoFunction()
func (c *InlineVariablesTest) TestInlineIntoFunction() {
	c.Test(
		"var x = false; var f = function() { alert(    x); };",
		"               var f = function() { alert(false); };")
	c.TestSame("var x = [false]; var f = function() { alert(x); };")
}

// Java: InlineVariablesTest.testNoInlineIntoNamedFunction()
func (c *InlineVariablesTest) TestNoInlineIntoNamedFunction() {
	c.TestSame("f(); var x = false; function f() { alert(x); };")
}

// Java: InlineVariablesTest.testInlineIntoNestedNonHoistedNamedFunctions()
func (c *InlineVariablesTest) TestInlineIntoNestedNonHoistedNamedFunctions() {
	c.Test(
		"f(); var x = false; if (false) function f() { alert(    x); };",
		"f();                if (false) function f() { alert(false); };")
}

// Java: InlineVariablesTest.testNoInlineIntoNestedNamedFunctions()
func (c *InlineVariablesTest) TestNoInlineIntoNestedNamedFunctions() {
	c.TestSame("f(); var x = false; function f() { if (false) { alert(x); } };")
}

// Java: InlineVariablesTest.testNoInlineMutatedVariable()
func (c *InlineVariablesTest) TestNoInlineMutatedVariable() {
	c.TestSame("var x = false; if (true) { var y = x; x = true; }")
}

// Java: InlineVariablesTest.testInlineImmutableMultipleTimes()
func (c *InlineVariablesTest) TestInlineImmutableMultipleTimes() {
	c.Test(
		"var x = null; var y =    x, z =    x;", //
		"              var y = null, z = null;")
	c.Test(
		"var x = 3; var y = x, z = x;", //
		"           var y = 3, z = 3;")
}

// Java: InlineVariablesTest.testInlineStringMultipleTimesAllStrings()
func (c *InlineVariablesTest) TestInlineStringMultipleTimesAllStrings() {
	c.Test(
		`var x = 'abcdefghijklmnopqrstuvwxyz';
var y = x, z = x;
`,
		`var y = 'abcdefghijklmnopqrstuvwxyz',
    z = 'abcdefghijklmnopqrstuvwxyz';
`)
}

// Java: InlineVariablesTest.testNoInlineBackwards()
func (c *InlineVariablesTest) TestNoInlineBackwards() {
	c.TestSame("var y = x; var x = null;")
}

// Java: InlineVariablesTest.testNoInlineOutOfBranch()
func (c *InlineVariablesTest) TestNoInlineOutOfBranch() {
	c.TestSame("if (true) var x = null; var y = x;")
}

// Java: InlineVariablesTest.testInterferingInlines()
func (c *InlineVariablesTest) TestInterferingInlines() {
	c.Test(
		"var a = 3; var f = function() { var x = a; alert(x); };",
		"           var f = function() {            alert(3); };")
}

// Java: InlineVariablesTest.testInlineIntoTryCatch()
func (c *InlineVariablesTest) TestInlineIntoTryCatch() {
	c.Test(
		`var a = true;
try { var b = a; }
catch (e) { var c = a + b; var d = true; }
finally { var f = a + b + c + d; }
`,
		`try { var b = true; }
catch (e) { var c = true + b; var d = true; }
finally { var f = true + b + c + d; }
`)
}

// Make sure that we still inline constants that are not provably
// written before they're read.
// Java: InlineVariablesTest.testInlineConstants()
func (c *InlineVariablesTest) TestInlineConstants() {
	c.Test(
		"function foo() { return  XXX; } var XXX = true;", //
		"function foo() { return true; }                ")
}

// Java: InlineVariablesTest.testInlineStringWhenWorthwhile()
func (c *InlineVariablesTest) TestInlineStringWhenWorthwhile() {
	c.Test("var x = 'a'; foo(x, x, x);", "foo('a', 'a', 'a');")
}

// Java: InlineVariablesTest.testInlineConstantAlias()
func (c *InlineVariablesTest) TestInlineConstantAlias() {
	c.Test(
		"var XXX = new Foo(); q(XXX); var YYY = XXX; bar(YYY)",
		"var XXX = new Foo(); q(XXX);                bar(XXX)")
}

// Java: InlineVariablesTest.testInlineConstantAliasWithAnnotation()
func (c *InlineVariablesTest) TestInlineConstantAliasWithAnnotation() {
	c.Test(
		"/** @const */ var xxx = new Foo(); q(xxx); var YYY = xxx; bar(YYY)",
		"/** @const */ var xxx = new Foo(); q(xxx);                bar(xxx)")
}

// Java: InlineVariablesTest.testInlineConstantAliasWithNonConstant()
func (c *InlineVariablesTest) TestInlineConstantAliasWithNonConstant() {
	c.Test(
		"var XXX = new Foo(); q(XXX); var y = XXX; bar(  y); baz(  y)",
		"var XXX = new Foo(); q(XXX);              bar(XXX); baz(XXX)")
}

// Java: InlineVariablesTest.testCascadingInlines()
func (c *InlineVariablesTest) TestCascadingInlines() {
	c.Test(
		`var XXX = 4;
function f() {
  var YYY = XXX;
  bar(YYY);
  baz(YYY);
}
`,
		`function f() {
  bar(4);
  baz(4);
}
`)
}

// Java: InlineVariablesTest.testNoInlineGetpropIntoCall()
func (c *InlineVariablesTest) TestNoInlineGetpropIntoCall() {
	c.Test("var a = b; a();", "b();")
	c.Test("var a = b.c; f(a);", "f(b.c);")
	c.TestSame("var a = b.c; a();")
}

// Java: InlineVariablesTest.testInlineFunctionDeclaration()
func (c *InlineVariablesTest) TestInlineFunctionDeclaration() {
	c.Test(
		"var f = function () {}; var a =              f;", //
		"                        var a = function () {};")
	c.Test(
		"var f = function () {}; foo(); var a =              f;", //
		"                        foo(); var a = function () {};")
	c.Test(
		"var f = function () {}; foo(             f);", //
		"                        foo(function () {});")
	c.Test(
		"var f = function () {}; function g() {var a = f; return a;}", //
		"var f = function () {}; function g() {           return f;}")
	c.TestSame("var f = function () {}; function g() {h(f);}")
}

// Java: InlineVariablesTest.test2388531()
func (c *InlineVariablesTest) Test2388531() {
	c.TestSame(
		`var f = function () {};
var g = function () {};
goog.inherits(f, g);
`)
	c.TestSame(
		`var f = function () {};
var g = function () {};
goog$inherits(f, g);
`)
}

// Java: InlineVariablesTest.testRecursiveFunction1()
func (c *InlineVariablesTest) TestRecursiveFunction1() {
	c.TestSame("var x = 0; (function x() { return x ? x() : 3; })();")
}

// Java: InlineVariablesTest.testRecursiveFunction2()
func (c *InlineVariablesTest) TestRecursiveFunction2() {
	c.TestSame("function y() { return y(); }")
}

// Java: InlineVariablesTest.testUnreferencedBleedingFunction()
func (c *InlineVariablesTest) TestUnreferencedBleedingFunction() {
	c.TestSame("var x = function y() {}")
}

// Java: InlineVariablesTest.testReferencedBleedingFunction()
func (c *InlineVariablesTest) TestReferencedBleedingFunction() {
	c.TestSame("var x = function y() { return y(); }")
}

// Java: InlineVariablesTest.testInlineAliases1()
func (c *InlineVariablesTest) TestInlineAliases1() {
	c.Test(
		"var x = this.foo(); this.bar(); var y = x; this.baz(y);",
		"var x = this.foo(); this.bar();            this.baz(x);")
}

// Java: InlineVariablesTest.testInlineAliases1b()
func (c *InlineVariablesTest) TestInlineAliases1b() {
	c.Test(
		"var x = this.foo(); this.bar(); var y; y = x; this.baz(y);",
		"var x = this.foo(); this.bar();            x; this.baz(x);")
}

// Java: InlineVariablesTest.testInlineAliases1c()
func (c *InlineVariablesTest) TestInlineAliases1c() {
	c.Test(
		"var x; x = this.foo(); this.bar(); var y = x; this.baz(y);",
		"var x; x = this.foo(); this.bar();            this.baz(x);")
}

// Java: InlineVariablesTest.testInlineAliases1d()
func (c *InlineVariablesTest) TestInlineAliases1d() {
	c.Test(
		"var x; x = this.foo(); this.bar(); var y; y = x; this.baz(y);",
		"var x; x = this.foo(); this.bar();            x; this.baz(x);")
}

// Java: InlineVariablesTest.testInlineAliases2()
func (c *InlineVariablesTest) TestInlineAliases2() {
	c.Test(
		`var x = this.foo();
this.bar();
function f() {
  var y = x;
  this.baz(y);
}
`,
		`var x = this.foo();
this.bar();
function f() {
  this.baz(x);
}
`)
}

// Java: InlineVariablesTest.testInlineAliases2b()
func (c *InlineVariablesTest) TestInlineAliases2b() {
	c.Test(
		`var x = this.foo();
this.bar();
function f() {
  var y;
  y = x;
  this.baz(y);
}
`,
		`var x = this.foo();
this.bar();
function f() {

  x;
  this.baz(x);
}
`)
}

// Java: InlineVariablesTest.testInlineAliases2c()
func (c *InlineVariablesTest) TestInlineAliases2c() {
	c.Test(
		`var x;
x = this.foo();
this.bar();
function f() {
  var y = x;
  this.baz(y);
}
`,
		`var x;
x = this.foo();
this.bar();
function f() {
  this.baz(x);
}
`)
}

// Java: InlineVariablesTest.testInlineAliases2d()
func (c *InlineVariablesTest) TestInlineAliases2d() {
	c.Test(
		`var x;
x = this.foo();
this.bar();
function f() {
  var y;
  y = x;
  this.baz(y);
}
`,
		`var x;
x = this.foo();
this.bar();
function f() {

  x;
  this.baz(x);
}
`)
}

// Java: InlineVariablesTest.testInlineAliasesInLoop()
func (c *InlineVariablesTest) TestInlineAliasesInLoop() {
	c.Test(
		`function f() {
  var x = extern();
  for (var i = 0; i < 5; i++) {
    (function() {
       var y = x; window.setTimeout(function() { extern(y); }, 0);
     })();
  }
}
`,
		`function f() {
  var x = extern();
  for (var i = 0; i < 5; i++) {
    (function() {
       window.setTimeout(function() { extern(x); }, 0);
     })();
  }
}
`)
}

// Java: InlineVariablesTest.testNoInlineAliasesInLoop()
func (c *InlineVariablesTest) TestNoInlineAliasesInLoop() {
	c.TestSame(
		`function f() {
  for (var i = 0; i < 5; i++) {
    var x = extern();
    (function() {
       var y = x; window.setTimeout(function() { extern(y); }, 0);
     })();
  }
}
`)
}

// Java: InlineVariablesTest.testNoInlineAliases1()
func (c *InlineVariablesTest) TestNoInlineAliases1() {
	c.TestSame("var x = this.foo(); this.bar(); var y = x; x = 3; this.baz(y);")
}

// Java: InlineVariablesTest.testNoInlineAliases1b()
func (c *InlineVariablesTest) TestNoInlineAliases1b() {
	c.TestSame("var x = this.foo(); this.bar(); var y; y = x; x = 3; this.baz(y);")
}

// Java: InlineVariablesTest.testNoInlineAliases2()
func (c *InlineVariablesTest) TestNoInlineAliases2() {
	c.TestSame("var x = this.foo(); this.bar(); var y = x; y = 3; this.baz(y); ")
}

// Java: InlineVariablesTest.testNoInlineAliases2b()
func (c *InlineVariablesTest) TestNoInlineAliases2b() {
	c.TestSame("var x = this.foo(); this.bar(); var y; y = x; y = 3; this.baz(y); ")
}

// Java: InlineVariablesTest.testNoInlineAliases3()
func (c *InlineVariablesTest) TestNoInlineAliases3() {
	c.TestSame(
		`var x = this.foo();
this.bar();
function f() {
  var y = x;
  g();
  this.baz(y);
}
function g() {
  x = 3;
}
`)
}

// Java: InlineVariablesTest.testNoInlineAliases3b()
func (c *InlineVariablesTest) TestNoInlineAliases3b() {
	c.TestSame(
		`var x = this.foo();
this.bar();

function f() {
var y;
y = x;
g();
this.baz(y);
}
function g() {
x = 3;
}
`)
}

// Java: InlineVariablesTest.testNoInlineAliases4()
func (c *InlineVariablesTest) TestNoInlineAliases4() {
	c.TestSame(
		`var x = this.foo();
this.bar();
function f() {
  var y = x;
  y = 3;
  this.baz(y);
}
`)
}

// Java: InlineVariablesTest.testNoInlineAliases4b()
func (c *InlineVariablesTest) TestNoInlineAliases4b() {
	c.TestSame(
		`var x = this.foo();
this.bar();
function f() {
  var y;
  y = x;
  y = 3;
  this.baz(y);
}
`)
}

// Java: InlineVariablesTest.testNoInlineAliases5()
func (c *InlineVariablesTest) TestNoInlineAliases5() {
	c.TestSame(
		`var x = this.foo();
this.bar();
var y = x;
this.bing();
this.baz(y);
x = 3;
`)
}

// Java: InlineVariablesTest.testNoInlineAliases5b()
func (c *InlineVariablesTest) TestNoInlineAliases5b() {
	c.TestSame(
		`var x = this.foo();
this.bar();
var y;
y = x;
this.bing();
this.baz(y);
x = 3;
`)
}

// Java: InlineVariablesTest.testNoInlineAliases6()
func (c *InlineVariablesTest) TestNoInlineAliases6() {
	c.TestSame(
		`var x = this.foo();
this.bar();
var y = x;
this.bing();
this.baz(y);
y = 3;
`)
}

// Java: InlineVariablesTest.testNoInlineAliases6b()
func (c *InlineVariablesTest) TestNoInlineAliases6b() {
	c.TestSame(
		`var x = this.foo();
this.bar();
var y;
y = x;
this.bing();
this.baz(y);
y = 3;
`)
}

// Java: InlineVariablesTest.testNoInlineAliases7()
func (c *InlineVariablesTest) TestNoInlineAliases7() {
	c.TestSame(
		`var x = this.foo();
this.bar();
function f() {
  var y = x;
  this.bing();
  this.baz(y);
  x = 3;
}
`)
}

// Java: InlineVariablesTest.testNoInlineAliases7b()
func (c *InlineVariablesTest) TestNoInlineAliases7b() {
	c.TestSame(
		`var x = this.foo();
this.bar();
function f() {
  var y;
  y = x;
  this.bing();
  this.baz(y);
  x = 3;
}
`)
}

// Java: InlineVariablesTest.testNoInlineAliases8()
func (c *InlineVariablesTest) TestNoInlineAliases8() {
	c.TestSame(
		`var x = this.foo();
this.bar();
function f() {
  var y = x;
  this.baz(y);
  y = 3;
}
`)
}

// Java: InlineVariablesTest.testNoInlineAliases8b()
func (c *InlineVariablesTest) TestNoInlineAliases8b() {
	c.TestSame(
		`var x = this.foo();
this.bar();
function f() {
  var y;
  y = x;
  this.baz(y);
  y = 3;
}
`)
}

// Java: InlineVariablesTest.testSideEffectOrder()
func (c *InlineVariablesTest) TestSideEffectOrder() {
	// z can not be changed by the call to y, so x can be inlined.
	EXTERNS := "var z; function f(){}"
	c.Test(
		c.Externs(EXTERNS), //
		c.Srcs("var x = f(y.a, y); z =          x"),
		c.Expected("               z = f(y.a, y);"))
	// z.b can be changed by the call to y, so x can not be inlined.
	c.TestSame(c.Externs(EXTERNS), c.Srcs("var x = f(y.a, y); z.b = x;"))
}

// Java: InlineVariablesTest.testInlineParameterAlias1()
func (c *InlineVariablesTest) TestInlineParameterAlias1() {
	c.Test(
		`function f(x) {
  var y = x;
  g();
  y;y;
}
`,
		`function f(x) {
  g();
  x;x;
}
`)
}

// Java: InlineVariablesTest.testInlineParameterAlias2()
func (c *InlineVariablesTest) TestInlineParameterAlias2() {
	c.Test(
		`function f(x) {
  var y; y = x;
  g();
  y;y;
}
`,
		`function f(x) {
  x;
  g();
  x;x;
}
`)
}

// Java: InlineVariablesTest.testInlineFunctionAlias1a()
func (c *InlineVariablesTest) TestInlineFunctionAlias1a() {
	c.Test(
		`function f(x) {}
var y = f;
g();
y();y();
`,
		`var y = function(x) {};
g();
y();y();
`)
}

// Java: InlineVariablesTest.testInlineFunctionAlias1b()
func (c *InlineVariablesTest) TestInlineFunctionAlias1b() {
	c.Test(
		`function f(x) {};
f;var y = f;
g();
y();y();
`,
		`function f(x) {};
f;g();
f();f();
`)
}

// Java: InlineVariablesTest.testInlineFunctionAlias2a()
func (c *InlineVariablesTest) TestInlineFunctionAlias2a() {
	c.Test(
		`function f(x) {}
var y; y = f;
g();
y();y();
`,
		`var y; y = function(x) {};
g();
y();y();
`)
}

// Java: InlineVariablesTest.testInlineFunctionAlias2b()
func (c *InlineVariablesTest) TestInlineFunctionAlias2b() {
	c.Test(
		`function f(x) {};
f; var y; y = f;
g();
y();y();
`,
		`function f(x) {};
f; f;
g();
f();f();
`)
}

// Java: InlineVariablesTest.testInlineSwitchVar()
func (c *InlineVariablesTest) TestInlineSwitchVar() {
	c.Test(
		"var x = y; switch (x) {}", //
		"           switch (y) {}")
}

// Java: InlineVariablesTest.testInlineSwitchLet()
func (c *InlineVariablesTest) TestInlineSwitchLet() {
	c.Test(
		"let x = y; switch (x) {}", //
		"           switch (y) {}")
}

// Successfully inlines 'values' and 'e'
// Java: InlineVariablesTest.testInlineIntoForLoop1()
func (c *InlineVariablesTest) TestInlineIntoForLoop1() {
	c.Test(
		`function calculate_hashCode() {
  var values = [1, 2, 3, 4, 5];
  var hashCode = 1;
  for (var $array = values, i = 0; i < $array.length; i++) {
    var e = $array[i];
    hashCode = 31 * hashCode + calculate_hashCode(e);
  }
  return hashCode;
}
`,
		`function calculate_hashCode() {
  var hashCode = 1;
  var $array = [1, 2, 3, 4, 5];
  var i = 0;
  for (; i < $array.length; i++) {
    hashCode = 31 * hashCode + calculate_hashCode($array[i]);
  }
  return hashCode;
}
`)
}

// Inlines 'e' but fails to inline 'values'
// TODO(tbreisacher): Investigate and see if we can improve this.
// Java: InlineVariablesTest.testInlineIntoForLoop2()
func (c *InlineVariablesTest) TestInlineIntoForLoop2() {
	c.Test(
		`function calculate_hashCode() {
  let values = [1, 2, 3, 4, 5];
  let hashCode = 1;
  for (let $array = values, i = 0; i < $array.length; i++) {
    let e = $array[i];
    hashCode = 31 * hashCode + calculate_hashCode(e);
  }
  return hashCode;
}
`,
		`function calculate_hashCode() {
  let values = [1, 2, 3, 4, 5];
  let hashCode = 1;
  for (let $array = values, i = 0; i < $array.length; i++) {
    hashCode = 31 * hashCode + calculate_hashCode($array[i]);
  }
  return hashCode;
}
`)
}

// This used to be inlined, but regressed when we switched to the ES6 scope creator.
// Java: InlineVariablesTest.testNoInlineCatchAliasVar1()
func (c *InlineVariablesTest) TestNoInlineCatchAliasVar1() {
	c.TestSame(
		`try {
} catch (e) {
  var y = e;
  g();
  y;y;
}
`)
}

// This used to be inlined, but regressed when we switched to the ES6 scope creator.
// Java: InlineVariablesTest.testNoInlineCatchAliasVar2()
func (c *InlineVariablesTest) TestNoInlineCatchAliasVar2() {
	c.TestSame(
		`try {
} catch (e) {
  var y; y = e;
  g();
  y;y;
}
`)
}

// Java: InlineVariablesTest.testInlineCatchAliasLet1()
func (c *InlineVariablesTest) TestInlineCatchAliasLet1() {
	c.Test(
		`try {
} catch (e) {
  let y = e;
  g();
  y;y;
}
`,
		`try {
} catch (e) {
  g();
  e;e;
}
`)
}

// Java: InlineVariablesTest.testInlineCatchAliasLet2()
func (c *InlineVariablesTest) TestInlineCatchAliasLet2() {
	c.Test(
		`try {
} catch (e) {
  let y; y = e;
  g();
  y;y;
}
`,
		`try {
} catch (e) {
  e;
  g();
  e;e;
}
`)
}

// Java: InlineVariablesTest.testInlineThis()
func (c *InlineVariablesTest) TestInlineThis() {
	c.Test(
		`/** @constructor */
function C() {}

C.prototype.m = function() {
  var self = this;
  if (true) {
    alert(self);
  }
};
`,
		`(/** @constructor */
function() {}).prototype.m = function() {
  if (true) {
    alert(this);
  }
};
`)
}

// Java: InlineVariablesTest.testVarInBlock1()
func (c *InlineVariablesTest) TestVarInBlock1() {
	c.Test(
		"function f(x) { if (true) {var y = x; y; y;} }", //
		"function f(x) { if (true) {           x; x;} }")
}

// Java: InlineVariablesTest.testVarInBlock2()
func (c *InlineVariablesTest) TestVarInBlock2() {
	c.Test(
		"function f(x) { switch (0) { case 0: { var y = x; y; y; } } }",
		"function f(x) { switch (0) { case 0: { x; x; } } }")
}

// Java: InlineVariablesTest.testLocalsOnly1()
func (c *InlineVariablesTest) TestLocalsOnly1() {
	c.inlineLocalsOnly = true
	c.Test(
		"var x=1; x; function f() {var x = 1; x;}", //
		"var x=1; x; function f() {           1;}")
}

// Java: InlineVariablesTest.testLocalsOnly2()
func (c *InlineVariablesTest) TestLocalsOnly2() {
	c.inlineLocalsOnly = true
	c.Test(
		`/** @const */
var X=1;
X;
function f() {
  /** @const */
  var X = 1; X;
}
`,
		`/** @const */
var X=1;
X;
function f() {
  1;
}
`)
}

// Java: InlineVariablesTest.testInlineUndefined1()
func (c *InlineVariablesTest) TestInlineUndefined1() {
	c.Test("var x; x;", "void 0;")
}

// Java: InlineVariablesTest.testInlineUndefined2()
func (c *InlineVariablesTest) TestInlineUndefined2() {
	c.TestSame("var x; x++;")
}

// Java: InlineVariablesTest.testInlineUndefined3()
func (c *InlineVariablesTest) TestInlineUndefined3() {
	c.TestSame("var x; var x;")
}

// Java: InlineVariablesTest.testInlineUndefined4()
func (c *InlineVariablesTest) TestInlineUndefined4() {
	c.Test("var x; x; x;", "void 0; void 0;")
}

// Java: InlineVariablesTest.testInlineUndefined5()
func (c *InlineVariablesTest) TestInlineUndefined5() {
	c.TestSame("var x; for(x in a) {}")
}

// Java: InlineVariablesTest.testIssue90()
func (c *InlineVariablesTest) TestIssue90() {
	c.Test(
		"var x;      x && alert(1)", //
		"       void 0 && alert(1)")
}

// Java: InlineVariablesTest.testRenamePropertyFunction()
func (c *InlineVariablesTest) TestRenamePropertyFunction() {
	c.TestSame(
		`var JSCompiler_renameProperty;
JSCompiler_renameProperty('foo')
`)
}

// Java: InlineVariablesTest.testThisAlias()
func (c *InlineVariablesTest) TestThisAlias() {
	c.Test(
		"function f() { var a = this;    a.y();    a.z(); }", //
		"function f() {               this.y(); this.z(); }")
}

// Java: InlineVariablesTest.testThisEscapedAlias()
func (c *InlineVariablesTest) TestThisEscapedAlias() {
	c.TestSame("function f() { var a = this; var g = function() { a.y(); }; a.z(); }")
}

// Java: InlineVariablesTest.testThisAlias_classMemberFieldDef()
func (c *InlineVariablesTest) TestThisAlias_classMemberFieldDef() {
	c.TestSame(
		`function f() {
  var self = this;
  use(self);
  class H {
    owner = self;
  }
}
`)
}

// Java: InlineVariablesTest.testThisAlias_classStaticMemberFieldDef()
func (c *InlineVariablesTest) TestThisAlias_classStaticMemberFieldDef() {
	c.TestSame(
		`function f() {
  var self = this;
  use(self);
  class H {
    static owner = self;
  }
}
`)
}

// Java: InlineVariablesTest.testThisAlias_classComputedFieldDefValue()
func (c *InlineVariablesTest) TestThisAlias_classComputedFieldDefValue() {
	c.TestSame(
		`function f() {
  var self = this;
  use(self);
  class H {
    ['owner'] = self;
  }
}
`)
}

// Java: InlineVariablesTest.testThisAlias_classComputedFieldDefKey()
func (c *InlineVariablesTest) TestThisAlias_classComputedFieldDefKey() {
	c.Test(
		`function f() {
  var self = this;
  use(self);
  class H {
    [self] = 1;
  }
}
`,
		`function f() {
  use(this);
  class H {
    [this] = 1;
  }
}
`)
}

// Java: InlineVariablesTest.testThisAlias_classStaticBlock()
func (c *InlineVariablesTest) TestThisAlias_classStaticBlock() {
	c.TestSame(
		`function f() {
  var self = this;
  use(self);
  class H {
    static {
      use(self);
    }
  }
}
`)
}

// Java: InlineVariablesTest.testThisAlias_globalClassMemberFieldDef()
func (c *InlineVariablesTest) TestThisAlias_globalClassMemberFieldDef() {
	c.TestSame(
		`var self = this;
use(self);
class H {
  owner = self;
}
`)
}

// Java: InlineVariablesTest.testInlineNamedFunction()
func (c *InlineVariablesTest) TestInlineNamedFunction() {
	c.Test("function f() {} f();", "(function(){})()")
}

// Java: InlineVariablesTest.testIssue378ModifiedArguments1()
func (c *InlineVariablesTest) TestIssue378ModifiedArguments1() {
	c.TestSame(
		`function g(callback) {
  var f = callback;
  arguments[0] = this;
  f.apply(this, arguments);
}
`)
}

// Java: InlineVariablesTest.testIssue378ModifiedArguments2()
func (c *InlineVariablesTest) TestIssue378ModifiedArguments2() {
	c.TestSame(
		`function g(callback) {
  /** @const */
  var f = callback;
  arguments[0] = this;
  f.apply(this, arguments);
}
`)
}

// Java: InlineVariablesTest.testIssue378EscapedArguments1()
func (c *InlineVariablesTest) TestIssue378EscapedArguments1() {
	c.TestSame(
		`function g(callback) {
  var f = callback;
  h(arguments,this);
  f.apply(this, arguments);
}
function h(a,b) {
  a[0] = b;
}
`)
}

// Java: InlineVariablesTest.testIssue378EscapedArguments2()
func (c *InlineVariablesTest) TestIssue378EscapedArguments2() {
	c.TestSame(
		`function g(callback) {
  /** @const */
  var f = callback;
  h(arguments,this);
  f.apply(this);
}
function h(a,b) {
  a[0] = b;
}
`)
}

// Java: InlineVariablesTest.testIssue378EscapedArguments3()
func (c *InlineVariablesTest) TestIssue378EscapedArguments3() {
	c.Test(
		`function g(callback) {
  var f = callback;
  f.apply(this, arguments);
}
`,
		`function g(callback) {

  callback.apply(this, arguments);
}
`)
}

// Java: InlineVariablesTest.testIssue378EscapedArguments4()
func (c *InlineVariablesTest) TestIssue378EscapedArguments4() {
	c.Test(
		`function g(callback) {
  var f = callback;
  h(arguments[0],this);
  f.apply(this, arguments);
}
function h(a,b) {
  a[0] = b;
}
`,
		`function g(callback) {

  h(arguments[0],this);
  callback.apply(this, arguments);
}
function h(a,b) {
  a[0] = b;
}
`)
}

// Java: InlineVariablesTest.testIssue378ArgumentsRead1()
func (c *InlineVariablesTest) TestIssue378ArgumentsRead1() {
	c.Test(
		`function g(callback) {
  var f = callback;
  var g = arguments[0];
  f.apply(this, arguments);
}
`,
		`function g(callback) {

  var g = arguments[0];
  callback.apply(this, arguments);
}
`)
}

// Java: InlineVariablesTest.testIssue378ArgumentsRead2()
func (c *InlineVariablesTest) TestIssue378ArgumentsRead2() {
	c.Test(
		`function g(callback) {
  var f = callback;
  h(arguments[0],this);
  f.apply(this, arguments[0]);
}
function h(a,b) {
  a[0] = b;
}
`,
		`function g(callback) {
  h(arguments[0],this);
  callback.apply(this, arguments[0]);
}
function h(a,b) {
  a[0] = b;
}
`)
}

// Java: InlineVariablesTest.testArgumentsModifiedInOuterFunction()
func (c *InlineVariablesTest) TestArgumentsModifiedInOuterFunction() {
	c.Test(
		`function g(callback) {
  var f = callback;
  arguments[0] = this;
  f.apply(this, arguments);
  function inner(callback) {
    var x = callback;
    x.apply(this);
  }
}
`,
		`function g(callback) {
  var f = callback;
  arguments[0] = this;
  f.apply(this, arguments);
  function inner(callback) {
    callback.apply(this);
  }
}
`)
}

// Java: InlineVariablesTest.testArgumentsModifiedInInnerFunction()
func (c *InlineVariablesTest) TestArgumentsModifiedInInnerFunction() {
	c.Test(
		`function g(callback) {
  var f = callback;
  f.apply(this, arguments);
  function inner(callback) {
    var x = callback;
    arguments[0] = this;
    x.apply(this);
  }
}
`,
		`function g(callback) {
  callback.apply(this, arguments);
  function inner(callback) {
    var x = callback;
    arguments[0] = this;
    x.apply(this);
  }
}
`)
}

// Java: InlineVariablesTest.testNoInlineRedeclaredExterns()
func (c *InlineVariablesTest) TestNoInlineRedeclaredExterns() {
	externs := "var test = 1;"
	code := "/** @suppress {duplicate} */ var test = 2;alert(test);"
	c.TestSame(c.Externs(externs), c.Srcs(code))
}

// Java: InlineVariablesTest.testBug6598844()
func (c *InlineVariablesTest) TestBug6598844() {
	c.TestSame(
		`function F() { this.a = 0; }
F.prototype.inc = function() { this.a++; return 10; };
F.prototype.bar = function() { var x = this.inc(); this.a += x; };
`)
}

// Java: InlineVariablesTest.testExternalIssue1053()
func (c *InlineVariablesTest) TestExternalIssue1053() {
	c.TestSame("var u; function f() { u = Random(); var x = u; f(); alert(x===u)}")
}

// Java: InlineVariablesTest.testHoistedFunction1()
func (c *InlineVariablesTest) TestHoistedFunction1() {
	c.Test("var x = 1; function f() { return x; }", "function f() { return 1; }")
}

// Java: InlineVariablesTest.testHoistedFunction2()
func (c *InlineVariablesTest) TestHoistedFunction2() {
	c.TestSame(
		`var impl_0;
b(a());
function a() { impl_0 = {}; }
function b() { window['f'] = impl_0; }
`)
}

// Java: InlineVariablesTest.testHoistedFunction3()
func (c *InlineVariablesTest) TestHoistedFunction3() {
	c.TestSame(
		`var impl_0;
b();
impl_0 = 1;
function b() { window['f'] = impl_0; }
`)
}

// Java: InlineVariablesTest.testHoistedFunction4()
func (c *InlineVariablesTest) TestHoistedFunction4() {
	c.Test(
		`var impl_0;
impl_0 = 1;
b();
function b() { window['f'] = impl_0; }
`,
		`1;
b();
function b() { window['f'] = 1; }
`)
}

// Java: InlineVariablesTest.testHoistedFunction5()
func (c *InlineVariablesTest) TestHoistedFunction5() {
	c.TestSame(
		`a();
var debug = 1;
function b() { return debug; }
function a() { return b(); }
`)
}

// Java: InlineVariablesTest.testHoistedFunction6()
func (c *InlineVariablesTest) TestHoistedFunction6() {
	c.Test(
		`var debug = 1;
a();
function b() { return debug; }
function a() { return b(); }
`,
		`a();
function b() { return 1; }
function a() { return b(); }
`)
}

// Java: InlineVariablesTest.testIssue354()
func (c *InlineVariablesTest) TestIssue354() {
	c.Test(
		`var enabled = true;
function Widget() {}
Widget.prototype = {
  frob: function() {
    search();
  }
};
function search() {
  if (enabled)
    alert(1);
  else
    alert(2);
}
window.foo = new Widget();
window.bar = search;
`,
		`function Widget() {}
Widget.prototype = {
  frob: function() {
    search();
  }
};
function search() {
  if (true)
    alert(1);
  else
    alert(2);
}
window.foo = new Widget();
window.bar = search;
`)
}

// Test respect for scopes and blocks
// Java: InlineVariablesTest.testIssue1177()
func (c *InlineVariablesTest) TestIssue1177() {
	c.TestSame("function x_64(){var x_7;for(;;);var x_68=x_7=x_7;}")
	c.TestSame("function x_64(){var x_7;for(;;);var x_68=x_7=x_7++;}")
	c.TestSame("function x_64(){var x_7;for(;;);var x_68=x_7=x_7*2;}")
}

// GitHub issue #1234: https://github.com/google/closure-compiler/issues/1234
// Java: InlineVariablesTest.testSwitchGithubIssue1234()
func (c *InlineVariablesTest) TestSwitchGithubIssue1234() {
	c.TestSame(
		`var x;
switch ('a') {
  case 'a':
    break;
  default:
    x = 1;
    break;
}
use(x);
`)
}

// Java: InlineVariablesTest.testLateGuardedAssign()
func (c *InlineVariablesTest) TestLateGuardedAssign() {
	c.Test(
		`let x;
function f() {
  x === void 0 && (x = 1);
  const y = x;
  use(y);
}
`,
		`let x;
function f() {
  x === void 0 && (x = 1);

  use(x);
}
`)
}

// Java: InlineVariablesTest.testLetConst()
func (c *InlineVariablesTest) TestLetConst() {
	c.Test(
		`function f(x) {
  if (true) {
    let y = x; y; y;
  }
}
`,
		`function f(x) {
  if (true) {
    x; x;
  }
}
`)
	c.Test(
		`function f(x) {
  if (true) {
    const y = x; y; y;
    }
  }
`,
		`function f(x) {
  if (true) {
    x; x;
  }
}
`)
	c.Test(
		`function f(x) {
  let y;
  {
    let y = x; y;
  }
}
`,
		`function f(x) {
  let y;
  {
    x;
  }
}
`)
	c.Test(
		`function f(x) {
  let y = x; y; const g = 2;
  {
    const g = 3; let y = g; y;
  }
}
`,
		`function f(x) {
  x;
  {3;}
}
`)
}

// Java: InlineVariablesTest.testGenerators()
func (c *InlineVariablesTest) TestGenerators() {
	c.Test(
		`function* f() {
  let x = 1;
  yield x;
}
`,
		`function* f() {
  yield 1;
}
`)
	c.Test(
		`function* f(x) {
  let y = x++
  yield y;
}
`,
		`function* f(x) {
  yield x++;
}
`)
}

// Java: InlineVariablesTest.testForOf()
func (c *InlineVariablesTest) TestForOf() {
	c.TestSame(" var i = 0; for(i of n) {}")
	c.TestSame("for( var i of n) { var x = i; }")
}

// Java: InlineVariablesTest.testTemplateStrings()
func (c *InlineVariablesTest) TestTemplateStrings() {
	c.Test(
		" var name = 'Foo'; `Hello ${name}`", //
		"                   `Hello ${'Foo'}`")
	c.Test(
		"var name = 'Foo'; var foo = name; `Hello ${foo}`", //
		"                                  `Hello ${'Foo'}`")
	c.Test(
		"var age = 3; `Age: ${age}`", //
		"             `Age: ${3}`")
}

// Java: InlineVariablesTest.testTaggedTemplateLiterals()
func (c *InlineVariablesTest) TestTaggedTemplateLiterals() {
	c.Test(
		"var name = 'Foo';\n"+
			"function myTag(strings, nameExp, numExp) {\n"+
			"  var modStr;\n"+
			"  if (numExp > 2) {\n"+
			"    modStr = nameExp + 'Bar'\n"+
			"  } else {\n"+
			"    modStr = nameExp + 'BarBar'\n"+
			"  }\n"+
			"}\n"+
			"var output = myTag`My name is ${name} ${3}`;\n",
		"var output = function(strings, nameExp, numExp) {\n"+
			"  var modStr;\n"+
			"  if (numExp > 2) {\n"+
			"    modStr = nameExp + 'Bar'\n"+
			"  } else {\n"+
			"    modStr = nameExp + 'BarBar'\n"+
			"  }\n"+
			"}`My name is ${'Foo'} ${3}`;\n")
	c.Test(
		"var name = 'Foo';\n"+
			"function myTag(strings, nameExp, numExp) {\n"+
			"  var modStr;\n"+
			"  if (numExp > 2) {\n"+
			"    modStr = nameExp + 'Bar'\n"+
			"  } else {\n"+
			"    modStr = nameExp + 'BarBar'\n"+
			"  }\n"+
			"}\n"+
			"var output = myTag`My name is ${name} ${3}`;\n"+
			"output = myTag`My name is ${name} ${2}`;\n",
		"function myTag(strings, nameExp, numExp) {\n"+
			"  var modStr;\n"+
			"  if (numExp > 2) {\n"+
			"    modStr = nameExp + 'Bar'\n"+
			"  } else {\n"+
			"    modStr = nameExp + 'BarBar'\n"+
			"  }\n"+
			"}\n"+
			"var output = myTag`My name is ${'Foo'} ${3}`;\n"+
			"output = myTag`My name is ${'Foo'} ${2}`;\n")
}

// Java: InlineVariablesTest.testDestructuring()
func (c *InlineVariablesTest) TestDestructuring() {
	c.Test(
		`var [a, b, c] = [1, 2, 3]
var x = a;
x; x;
`,
		`var a;
var b;
var c;
[a, b, c] = [1, 2, 3]
var x = a; // we fail to inline a into x because normalization adds stub declaration
// for a. But that's an acceptable compromise.
x; x;
`)
	c.TestSame("var x = 1; ({[0]: x} = {});")
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: InlineVariablesTest.dontInlineConditionalDefaultValueAssignment()
func (c *InlineVariablesTest) TestDontInlineConditionalDefaultValueAssignment() {
	c.TestSame("let x; const {y = (x = 0)} = obj(); use(x);")
}

// Java: InlineVariablesTest.testFunctionInlinedAcrossScript()
func (c *InlineVariablesTest) TestFunctionInlinedAcrossScript() {
	srcs := []string{
		"function f() {}", //
		"use(f);",
	}

	expected := []string{
		"", //
		"use(function() {});",
	}

	c.Test(c.Srcs(srcs), c.Expected(expected))
}

// Java: InlineVariablesTest.testFunctionVarInliningElidesFunctionName()
func (c *InlineVariablesTest) TestFunctionVarInliningElidesFunctionName() {
	// When there's only one reference to a function name (e.g. in the rhs of an assignment), we
	// should elide that function name.
	c.Test(
		"function notRecursive(x) { return 4 }; exports.y = notRecursive;",
		";exports.y = function(x) { return 4; }")
	// In this case, the function is recursive, so there's more than one reference, we should keep
	// the name.
	c.TestSame(
		`function factorial(x) { if (x == 1) return 1; return x + factorial(x - 1); };
exports.x = factorial;
`)
}

// Java: InlineVariablesTest.testNoInlinePastOptionalChainGetProp()
func (c *InlineVariablesTest) TestNoInlinePastOptionalChainGetProp() {
	c.Test(
		`function f() {
  var t = this.consume();
  var ready = this?.ready;
  return t + ready;
}
`,
		`function f() {
  var t = this.consume();
  return t + this?.ready;
}
`)
}

// Java: InlineVariablesTest.testNoInlinePastOptionalChainGetElem()
func (c *InlineVariablesTest) TestNoInlinePastOptionalChainGetElem() {
	c.Test(
		`function f() {
  var t = this.consume();
  var ready = this?.['ready'];
  return t + ready;
}
`,
		`function f() {
  var t = this.consume();
  return t + this?.['ready'];
}
`)
}

// Java: InlineVariablesTest.testNoInlinePastGetProp()
func (c *InlineVariablesTest) TestNoInlinePastGetProp() {
	c.Test(
		`function f() {
  var t = this.consume();
  var ready = this.ready;
  return t + ready;
}
`,
		`function f() {
  var t = this.consume();
  return t + this.ready;
}
`)
}

// Java: InlineVariablesTest.testNoInlinePastOptionalChainCall()
func (c *InlineVariablesTest) TestNoInlinePastOptionalChainCall() {
	c.TestSame(
		`function f() {
  var t = this.consume();
  var ready = this?.ready();
  return t + ready;
}
`)
}

// Java: InlineVariablesTest.testCrossScopeAliasChainWithInterveningReassignment()
func (c *InlineVariablesTest) TestCrossScopeAliasChainWithInterveningReassignment() {
	c.Test(
		`function f() {
  var C = compute();
  var D = C;
  {
    let A = D;
    C = attacker();
    sink(A); sink(A);
  }
}
`,
		`function f() {
  var C = compute();
  {
    let A = C;
    C = attacker();
    sink(A); sink(A);
  }
}
`)
}

// Java: InlineVariablesTest.testClassDeclaration_notInlined()
func (c *InlineVariablesTest) TestClassDeclaration_notInlined() {
	// When there's only one reference to a class name (e.g. in the rhs of an assignment), we
	// should elide that class name.
	c.TestSame("class MyClass { foo() {} } ns.y = MyClass;")
}

// Java: InlineVariablesTest.testClassExpressionInlined()
func (c *InlineVariablesTest) TestClassExpressionInlined() {
	c.Test("var C = class {}; use(C);", "use(class {});")
	c.Test("const C = class {}; use(C);", "use(class {});")
	c.Test("let C = class {}; use(C);", "use(class {});")
}

// Java: InlineVariablesTest.testModerateMotionStepLimit_doesNotInlinedAcrossTooManyNodes()
func (c *InlineVariablesTest) TestModerateMotionStepLimit_doesNotInlinedAcrossTooManyNodes() {
	var sb strings.Builder
	sb.WriteString("var C = class {};\n")
	// Generate enough side-effect-free class declarations to exceed MAX_MODERATE_MOTION_STEPS
	// (5000)
	for i := 0; i < 1000; i++ {
		sb.WriteString("class Extra" + strconv.Itoa(i) + " { foo() {} }\n")
	}
	sb.WriteString("use(C);\n")
	c.TestSame(sb.String())
}

// Java: InlineVariablesTest.testClassExpressionWithExtendsInlined()
func (c *InlineVariablesTest) TestClassExpressionWithExtendsInlined() {
	c.Test(
		"var C = class extends Base { foo() {} }; use(C);",
		"use(class extends Base { foo() {} });")
}

// Java: InlineVariablesTest.testClassWithSideEffectsNotMovedAcrossCalls()
func (c *InlineVariablesTest) TestClassWithSideEffectsNotMovedAcrossCalls() {
	c.TestSame(
		`var C = class { static { doSideEffects(); } };
doSomething();
use(C);
`)
}

// Java: InlineVariablesTest.testClassWithSideEffectingStaticField_notMovedAcrossCalls()
func (c *InlineVariablesTest) TestClassWithSideEffectingStaticField_notMovedAcrossCalls() {
	c.TestSame(
		`var C = class { static x = doSideEffects(); };
doSomething();
use(C);
`)
}

// Java: InlineVariablesTest.testClassWithPureStaticField_inlinedAcrossCalls()
func (c *InlineVariablesTest) TestClassWithPureStaticField_inlinedAcrossCalls() {
	c.TestSame(
		`var C = class { static x = 1; };
doSomething();
use(C);
`)
}

// Java: InlineVariablesTest.testClassWithImpureStaticField_notInlinedAcrossCalls()
func (c *InlineVariablesTest) TestClassWithImpureStaticField_notInlinedAcrossCalls() {
	c.TestSame(
		`var C = class { static x = compute(); };
doSomething();
use(C);
`)
}

// Java: InlineVariablesTest.testClassWithSideEffectingComputedProperty_notMovedAcrossCalls()
func (c *InlineVariablesTest) TestClassWithSideEffectingComputedProperty_notMovedAcrossCalls() {
	c.TestSame(
		`var C = class { [computeKey()]() {} };
doSomething();
use(C);
`)
}

// Java: InlineVariablesTest.testClassWithInstanceFieldSideEffects_canMoveBeforeInstantiation()
func (c *InlineVariablesTest) TestClassWithInstanceFieldSideEffects_canMoveBeforeInstantiation() {
	c.TestSame(
		`var C = class { x = compute(); };
doSomething();
use(C);
`)
}

// Java: InlineVariablesTest.testClassWithImpureExtendsNotMovedAcrossCalls()
func (c *InlineVariablesTest) TestClassWithImpureExtendsNotMovedAcrossCalls() {
	c.TestSame(
		`var C = class extends getBase() {};
doSomething();
use(C);
`)
}

// Java: InlineVariablesTest.testField_declaredConstAssignedInsideInstanceFieldInitializer()
func (c *InlineVariablesTest) TestField_declaredConstAssignedInsideInstanceFieldInitializer() {
	// TODO(b/572004916): @const assigned in field initializer is inlined before instantiation
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2022)
	c.Test(
		`/** @const */ var a;
class C {
  x = (a = 1);
}
alert(a);
new C();
`,
		`class C {
  x = 1;
}
alert(1);
new C();
`)
}

// Java: InlineVariablesTest.testField_declaredConstAssignedInsideInstanceFieldInitializer_local()
func (c *InlineVariablesTest) TestField_declaredConstAssignedInsideInstanceFieldInitializer_local() {
	// TODO(b/572004916): local @const assigned in field initializer is inlined before instantiation
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2022)
	c.Test(
		`function f() {
  /** @const */ var a;
  class C {
    x = (a = 1);
  }
  alert(a);
  return new C();
}
`,
		`function f() {
  class C {
    x = 1;
  }
  alert(1);
  return new C();
}
`)
}
