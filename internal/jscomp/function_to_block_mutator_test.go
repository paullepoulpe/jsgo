// Port of com/google/javascript/jscomp/FunctionToBlockMutatorTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/FunctionToBlockMutatorTest.java)

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

// FunctionToBlockMutatorTest is the unit test of FunctionToBlockMutator.
type FunctionToBlockMutatorTest struct {
	// Go-only: test-helper — the Go test running the current test method.
	t *testing.T

	needsDefaultResult bool
	isCallInLoop       bool
}

// Java: FunctionToBlockMutatorTest()
func NewFunctionToBlockMutatorTest() *FunctionToBlockMutatorTest {
	return &FunctionToBlockMutatorTest{}
}

func TestFunctionToBlockMutatorTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewFunctionToBlockMutatorTest() })
}

// JUnitSetT gives the instance the Go test running its test method.
func (c *FunctionToBlockMutatorTest) JUnitSetT(t *testing.T) { c.t = t }

// JUnitBefore runs the @Before method setUp().
func (c *FunctionToBlockMutatorTest) JUnitBefore() { c.SetUp() }

// JUnitAfter: the class has no @After method.
func (c *FunctionToBlockMutatorTest) JUnitAfter() {}

// Java: FunctionToBlockMutatorTest.setUp()
func (c *FunctionToBlockMutatorTest) SetUp() {
	c.needsDefaultResult = false
	c.isCallInLoop = false
}

// Java: FunctionToBlockMutatorTest.testMutateNoReturnWithoutResultAssignment()
func (c *FunctionToBlockMutatorTest) TestMutateNoReturnWithoutResultAssignment() {
	c.HelperMutate("function foo(){}; foo();", "{}", "foo")
}

// Java: FunctionToBlockMutatorTest.testMutateNoReturnWithResultAssignment()
func (c *FunctionToBlockMutatorTest) TestMutateNoReturnWithResultAssignment() {
	c.needsDefaultResult = true
	c.HelperMutate("function foo(){}; var result = foo();", "{result = void 0}", "foo")
}

// Java: FunctionToBlockMutatorTest.testMutateNoValueReturnWithoutResultAssignment()
func (c *FunctionToBlockMutatorTest) TestMutateNoValueReturnWithoutResultAssignment() {
	c.HelperMutateWithResultName("function foo(){return;}; foo();", "{}", "foo", nil)
}

// Java: FunctionToBlockMutatorTest.testMutateNoValueReturnWithResultAssignment()
func (c *FunctionToBlockMutatorTest) TestMutateNoValueReturnWithResultAssignment() {
	c.HelperMutate("function foo(){return;}; var result = foo();", "{result = void 0}", "foo")
}

// Java: FunctionToBlockMutatorTest.testMutateValueReturnWithoutResultAssignment()
func (c *FunctionToBlockMutatorTest) TestMutateValueReturnWithoutResultAssignment() {
	c.HelperMutateWithResultName("function foo(){return true;}; foo();", "{true;}", "foo", nil)
}

// Java: FunctionToBlockMutatorTest.testMutateValueReturnWithResultAssignment()
func (c *FunctionToBlockMutatorTest) TestMutateValueReturnWithResultAssignment() {
	c.needsDefaultResult = true
	c.HelperMutateWithResultName("function foo(){return true;}; var x=foo();", "{x=true}", "foo", jutil.StringPtr("x"))
}

// Java: FunctionToBlockMutatorTest.testMutateWithMultipleReturns()
func (c *FunctionToBlockMutatorTest) TestMutateWithMultipleReturns() {
	c.needsDefaultResult = true
	c.HelperMutate(
		"function foo(){ if (0) {return 0} else {return 1} }; var result=foo();",
		`{
  JSCompiler_inline_label_foo_0: {
    if (0) {
      result = 0;
      break JSCompiler_inline_label_foo_0
    } else {
      result = 1;
      break JSCompiler_inline_label_foo_0
    }
    result=void 0
  }
}
`,
		"foo")
}

// Java: FunctionToBlockMutatorTest.testMutateWithParameters1()
func (c *FunctionToBlockMutatorTest) TestMutateWithParameters1() {
	// Simple call with useless parameter
	c.HelperMutateWithResultName("function foo(a){return true;}; foo(x);", "{true}", "foo", nil)
}

// Java: FunctionToBlockMutatorTest.testMutateWithParameters2()
func (c *FunctionToBlockMutatorTest) TestMutateWithParameters2() {
	// Simple call with parameter
	c.HelperMutateWithResultName("function foo(a){return x;}; foo(x);", "{x}", "foo", nil)
}

// Java: FunctionToBlockMutatorTest.testMutateWithParameters3()
func (c *FunctionToBlockMutatorTest) TestMutateWithParameters3() {
	// Parameter has side-effects.
	c.HelperMutateWithResultName(
		"function foo(a){return a;}; function x() { foo(x++); }", //
		"{ x++ }",
		"foo",
		nil)
}

// Java: FunctionToBlockMutatorTest.testMutate8()
func (c *FunctionToBlockMutatorTest) TestMutate8() {
	// Parameter has side-effects.
	c.HelperMutateWithResultName(
		"function foo(a){return a+a;}; foo(x++);",
		"{var a$jscomp$inline_0 = x++; a$jscomp$inline_0 + a$jscomp$inline_0;}",
		"foo",
		nil)
}

// Java: FunctionToBlockMutatorTest.testMutate8_withConstVar()
func (c *FunctionToBlockMutatorTest) TestMutate8_withConstVar() {
	// Parameter has side-effects.
	c.HelperMutateWithResultName(
		"function foo(/** @const */ a){return a+a;}; foo(x++);",
		"{var a$jscomp$inline_0 = x++; a$jscomp$inline_0 + a$jscomp$inline_0;}",
		"foo",
		nil)
}

// Java: FunctionToBlockMutatorTest.testMutateInitializeUninitializedVars1()
func (c *FunctionToBlockMutatorTest) TestMutateInitializeUninitializedVars1() {
	c.isCallInLoop = true
	c.HelperMutateWithResultName(
		"function foo(a){var b;return a;}; foo(1);",
		"{var b$jscomp$inline_1 = void 0; 1;}",
		"foo",
		nil)
}

// Java: FunctionToBlockMutatorTest.testMutateInitializeUninitializedVars2()
func (c *FunctionToBlockMutatorTest) TestMutateInitializeUninitializedVars2() {
	c.HelperMutateWithResultName(
		"function foo(a) {var b; for(b in c)return a;}; foo(1);",
		`{
  JSCompiler_inline_label_foo_2:
  {
    var b$jscomp$inline_1;
    for (b$jscomp$inline_1 in c) {
      1;
      break JSCompiler_inline_label_foo_2;
    }
  }
}
`,
		"foo",
		nil)
}

// Java: FunctionToBlockMutatorTest.testMutateInitializeUninitializedVarsNestedFunction()
func (c *FunctionToBlockMutatorTest) TestMutateInitializeUninitializedVarsNestedFunction() {
	c.isCallInLoop = true
	c.HelperMutateWithResultName(
		`function foo(a){
  var b = function(){ var c; };
  return a;
};
foo(1);
`,
		`{
  var b$jscomp$inline_1 = function(){ var c$jscomp$inline_2; };
  1;
}
`,
		"foo",
		nil)
}

// Java: FunctionToBlockMutatorTest.testMutateInitializeUninitializedLets1()
func (c *FunctionToBlockMutatorTest) TestMutateInitializeUninitializedLets1() {
	c.isCallInLoop = true
	c.HelperMutateWithResultName(
		"function foo(a){let b;return a;}; foo(1);",
		"{let b$jscomp$inline_1 = void 0; 1;}",
		"foo",
		nil)
}

// Java: FunctionToBlockMutatorTest.testMutateInitializeUninitializedLets2()
func (c *FunctionToBlockMutatorTest) TestMutateInitializeUninitializedLets2() {
	c.HelperMutateWithResultName(
		"function foo(a) {for(let b in c)return a;}; foo(1);",
		`{
  JSCompiler_inline_label_foo_2:
  {
    for (let b$jscomp$inline_1 in c) {
      1;
      break JSCompiler_inline_label_foo_2;
    }
  }
}
`,
		"foo",
		nil)
}

// Java: FunctionToBlockMutatorTest.testMutateInitializeUninitializedVarInLoop()
func (c *FunctionToBlockMutatorTest) TestMutateInitializeUninitializedVarInLoop() {
	c.isCallInLoop = true
	// Demonstrates that `b` should be initialized to `undefined` because
	// it is used in a loop.
	c.HelperMutateWithResultName(
		`function foo(a) {
  for (var i = 0; i < 10; i++) {
    var b;
    return a;
  }
}
foo(1);
`,
		`{
  JSCompiler_inline_label_foo_3:
  {
    var b$jscomp$inline_2 = void 0;
    var i$jscomp$inline_1 = 0;
    for (; i$jscomp$inline_1 < 10; i$jscomp$inline_1++) {
      1;
      break JSCompiler_inline_label_foo_3;
    }
  }
}
`,
		"foo",
		nil)
}

// Java: FunctionToBlockMutatorTest.testMutateInitializedVarInLoop()
func (c *FunctionToBlockMutatorTest) TestMutateInitializedVarInLoop() {
	c.isCallInLoop = true
	// Demonstrates that b's initializer should not be hoisted out of the loop or
	// lost when "foo" is inlined into the body of a loop
	c.HelperMutateWithResultName(
		`function foo(a) {
  for (var i = 0; i < 10; i++) {
    var b = 1;
    return a;
  }
}
foo(1);
`,
		`{
  JSCompiler_inline_label_foo_3:
  {
    var b$jscomp$inline_2 = void 0;
    var i$jscomp$inline_1 = 0;
    for (; i$jscomp$inline_1 < 10; i$jscomp$inline_1++) {
      b$jscomp$inline_2 = 1;
      1;
      break JSCompiler_inline_label_foo_3;
    }
  }
}
`,
		"foo",
		nil)
}

// Java: FunctionToBlockMutatorTest.testMutateCallInLoopVars1()
func (c *FunctionToBlockMutatorTest) TestMutateCallInLoopVars1() {
	src := `function foo(a) {
  var B = bar();
  a;
};
foo(1);
`
	// baseline: outside a loop, the constant remains constant.
	c.isCallInLoop = false
	c.HelperMutateWithResultName(src, "{var B$jscomp$inline_1 = bar(); 1;}", "foo", nil)
	// ... in a loop, the constant-ness is removed.
	// TODO(johnlenz): update this test to look for the const annotation.
	c.isCallInLoop = true
	c.HelperMutateWithResultName(src, "{var B$jscomp$inline_1 = bar(); 1;}", "foo", nil)
}

// Java: FunctionToBlockMutatorTest.testMutateFunctionDefinition()
func (c *FunctionToBlockMutatorTest) TestMutateFunctionDefinition() {
	// Function declarations are rewritten as function expressions.
	c.HelperMutateWithResultName(
		"function foo(a){function g(){}}; foo(1);",
		"{var g$jscomp$inline_1 = function() {};}",
		"foo",
		nil)
}

// Java: FunctionToBlockMutatorTest.testMutateFunctionDefinitionHoisting()
func (c *FunctionToBlockMutatorTest) TestMutateFunctionDefinitionHoisting() {
	c.HelperMutateWithResultName(
		`function foo(a){
  var b = g(a);
  function g(c){ return c; }
  var c = i();
  function h(){}
  function i(){}
}
foo(1);
`,
		`{
  var g$jscomp$inline_1 = function(c$jscomp$inline_6) {return c$jscomp$inline_6};
  var h$jscomp$inline_2 = function(){};
  var i$jscomp$inline_3 = function(){};
  var b$jscomp$inline_4 = g$jscomp$inline_1(1);
  var c$jscomp$inline_5 = i$jscomp$inline_3();
}
`,
		"foo",
		nil)
}

// Java: FunctionToBlockMutatorTest.validateSourceInfo(Compiler, Node)
func (c *FunctionToBlockMutatorTest) ValidateSourceInfo(compiler *Compiler, subtree rhino.Node) {
	NewSourceInfoCheck(compiler).SetCheckSubTree(subtree)
	AssertThat(compiler.GetErrors()).IsEmpty()
}

// Java: FunctionToBlockMutatorTest.helperMutate(String, String, String)
func (c *FunctionToBlockMutatorTest) HelperMutate(code string, expectedResult string, fnName string) {
	c.HelperMutateWithResultName(code, expectedResult, fnName, jutil.StringPtr("result"))
}

// Java: FunctionToBlockMutatorTest.helperMutate(String, String, String, @Nullable String)
func (c *FunctionToBlockMutatorTest) HelperMutateWithResultName(
	code string, expectedResult string, fnName string, resultName *string) {
	compiler := NewCompiler()
	compiler.InitCompilerOptionsIfTesting()
	mutator :=
		NewFunctionToBlockMutator(compiler, compiler.GetUniqueNameIdSupplier())

	compiler.Init(
		[]*SourceFile{},
		[]*SourceFile{SourceFile_FromCode("[testcode]", code)},
		compiler.GetOptions())
	compiler.Parse()
	script := compiler.GetRoot().GetSecondChild().GetFirstChild()

	externs := compiler.GetExternsRoot()
	js := compiler.GetJsRoot()

	Normalize_CreateNormalizeForOptimizations(compiler).Process(externs, js)
	GatherGetterAndSetterProperties_Update(compiler, externs, js)
	NewPureFunctionIdentifier_Driver(compiler).Process(externs, js)

	fnNode := FunctionToBlockMutatorTest_FindFunction(script, fnName)

	expectedRoot := compiler.ParseTestCode(expectedResult)
	jutil.CheckState(compiler.GetErrorCount() == 0)
	expected := expectedRoot.GetFirstChild()

	// inline tester
	tester :=
		func(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
			result :=
				mutator.Mutate(fnName, fnNode, n, resultName, c.needsDefaultResult, c.isCallInLoop)
			c.ValidateSourceInfo(compiler, result)
			AssertNode(result).UsingSerializer(compiler.ToSourceOfNode).IsEqualTo(expected)

			n.ReplaceWith(result)
			NewValidityCheck_VerifyConstants(compiler, false).Process(externs, js)
			Normalize_Builder_(compiler).AssertOnChange(true).Build().Process(externs, js)
			result.ReplaceWith(n)

			return true
		}

	test := NewFunctionToBlockMutatorTest_TestCallback(fnName, tester)
	NodeTraversal_Traverse(compiler, script, test)
}

// FunctionToBlockMutatorTest_Method is the functional interface Method.
type FunctionToBlockMutatorTest_Method func(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool

// FunctionToBlockMutatorTest_TestCallback is the static class TestCallback.
type FunctionToBlockMutatorTest_TestCallback struct {
	callname string
	method   FunctionToBlockMutatorTest_Method
	complete bool
}

// Java: FunctionToBlockMutatorTest.TestCallback(String, Method)
func NewFunctionToBlockMutatorTest_TestCallback(callname string, method FunctionToBlockMutatorTest_Method) *FunctionToBlockMutatorTest_TestCallback {
	return &FunctionToBlockMutatorTest_TestCallback{callname: callname, method: method}
}

// Java: FunctionToBlockMutatorTest.TestCallback.shouldTraverse(NodeTraversal, Node, Node)
func (cb *FunctionToBlockMutatorTest_TestCallback) ShouldTraverse(nodeTraversal *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	return !cb.complete
}

// Java: FunctionToBlockMutatorTest.TestCallback.visit(NodeTraversal, Node, Node)
func (cb *FunctionToBlockMutatorTest_TestCallback) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsCall() {
		first := n.GetFirstChild()
		if first.IsName() && first.GetString() == cb.callname {
			cb.complete = cb.method(t, n, parent)
		}
	}

	if parent.IsNil() {
		AssertThat(cb.complete).IsTrue()
	}
}

// Java: FunctionToBlockMutatorTest.findFunction(Node, String)
func FunctionToBlockMutatorTest_FindFunction(n rhino.Node, name string) rhino.Node {
	if n.IsFunction() {
		if n.GetFirstChild().GetString() == name {
			return n
		}
	}

	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		result := FunctionToBlockMutatorTest_FindFunction(c, name)
		if !result.IsNil() {
			return result
		}
	}

	return rhino.Node{}
}
