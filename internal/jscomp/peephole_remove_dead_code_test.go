// Port of com/google/javascript/jscomp/PeepholeRemoveDeadCodeTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/PeepholeRemoveDeadCodeTest.java)

/*
 * Copyright 2004 The Closure Compiler Authors.
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

// PeepholeRemoveDeadCodeTest is the tests for PeepholeRemoveDeadCode in
// isolation. Tests for the interaction of multiple peephole passes are in
// PeepholeIntegrationTest.
type PeepholeRemoveDeadCodeTest struct {
	CompilerTestCaseBase
}

var peepholeRemoveDeadCodeTest_MATH = "/** @const */ var Math = {};\n" +
	"/** @nosideeffects */ Math.random = function(){};\n" +
	"/** @nosideeffects */ Math.sin = function(){};\n"

// Java: PeepholeRemoveDeadCodeTest()
func NewPeepholeRemoveDeadCodeTest() *PeepholeRemoveDeadCodeTest {
	c := &PeepholeRemoveDeadCodeTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, peepholeRemoveDeadCodeTest_MATH)
	return c
}

func TestPeepholeRemoveDeadCodeTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewPeepholeRemoveDeadCodeTest() })
}

// peepholeRemoveDeadCodeTest_processor is the anonymous CompilerPass of
// getProcessor.
type peepholeRemoveDeadCodeTest_processor struct {
	this     *PeepholeRemoveDeadCodeTest
	compiler *Compiler
}

// Java: PeepholeRemoveDeadCodeTest.getProcessor(Compiler).new CompilerPass() {...}.process(Node, Node)
func (p *peepholeRemoveDeadCodeTest_processor) Process(externs *rhino.Node, root *rhino.Node) {
	peepholePass :=
		NewPeepholeOptimizationsPass(p.compiler, p.this.GetName(), NewPeepholeRemoveDeadCode())
	peepholePass.Process(externs, root)
}

// Java: PeepholeRemoveDeadCodeTest.getProcessor(Compiler)
func (c *PeepholeRemoveDeadCodeTest) GetProcessor(compiler *Compiler) CompilerPass {
	return &peepholeRemoveDeadCodeTest_processor{this: c, compiler: compiler}
}

// Java: PeepholeRemoveDeadCodeTest.setUp()
func (c *PeepholeRemoveDeadCodeTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	// This pass doesn't need the AST to be normalized as it runs both before and after
	// `denormalize`. Since it runs PureFunctionsIdentifier pass which expects the AST to be
	// normalized, we're doing `enableNormalize` for these tests.
	c.EnableNormalize()
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.EnableComputeSideEffects()
}

// Java: PeepholeRemoveDeadCodeTest.getNumRepetitions()
func (c *PeepholeRemoveDeadCodeTest) GetNumRepetitions() int {
	return 2
}

// Java: PeepholeRemoveDeadCodeTest.foldSame(String)
func (c *PeepholeRemoveDeadCodeTest) FoldSame(js string) {
	c.TestSame(js)
}

// Java: PeepholeRemoveDeadCodeTest.fold(String, String)
func (c *PeepholeRemoveDeadCodeTest) Fold(js string, expected string) {
	c.Test(js, expected)
}

// Java: PeepholeRemoveDeadCodeTest.testRemoveNoOpLabelledStatement()
func (c *PeepholeRemoveDeadCodeTest) TestRemoveNoOpLabelledStatement() {
	c.Fold("a: break a;", "")
	c.Fold("a: { break a; }", "")
	c.Fold( //
		"a: { break a; console.log('unreachable'); }", //
		"")
	c.Fold( //
		"a: { break a; var x = 1; } x = 2;", //
		"var x; x = 2;")
	c.FoldSame("b: { var x = 1; } x = 2;")
	c.FoldSame("a: b: { var x = 1; } x = 2;")
}

// Java: PeepholeRemoveDeadCodeTest.testRemoveUselessLabelWithFollowingBreak()
func (c *PeepholeRemoveDeadCodeTest) TestRemoveUselessLabelWithFollowingBreak() {
	c.Fold("a:b: break b;", "")
	// Note: the break is only removed if the parent
	// is the break target.
	c.FoldSame("a:b: break a;")
}

// Java: PeepholeRemoveDeadCodeTest.testFoldBlock()
func (c *PeepholeRemoveDeadCodeTest) TestFoldBlock() {
	c.Fold("{{foo()}}", "foo()")
	c.Fold("{foo();{}}", "foo()")
	c.Fold("{{foo()}{}}", "foo()")
	c.Fold("{{foo()}{bar()}}", "foo();bar()")
	c.Fold("{if(false)foo(); {bar()}}", "bar()")
	c.Fold("{if(false)if(false)if(false)foo(); {bar()}}", "bar()")
	c.Fold("{'hi'}", "")
	c.Fold("{x==3}", "")
	c.Fold("{`hello ${foo}`}", "")
	c.Fold("{ (function(){x++}) }", "")
	c.FoldSame("function f(){return;}")
	c.Fold("function f(){return 3;}", "function f(){return 3}")
	c.FoldSame("function f(){if(x)return; x=3; return; }")
	c.Fold("{x=3;;;y=2;;;}", "x=3;y=2")
	// Cases to test for empty block.
	c.Fold("while(x()){x}", "while(x());")
	c.Fold("while(x()){x()}", "while(x())x()")
	c.Fold("for(x=0;x<100;x++){x}", "for(x=0;x<100;x++);")
	c.Fold("for(x in y){x}", "for(x in y);")
	c.Fold("for (x of y) {x}", "for(x of y);")
	c.FoldSame("for (let x = 1; x <10; x++ ) {}")
	c.FoldSame("for (var x = 1; x <10; x++ ) {}")
}

// Java: PeepholeRemoveDeadCodeTest.testFoldBlockWithDeclaration_notNormalized()
func (c *PeepholeRemoveDeadCodeTest) TestFoldBlockWithDeclaration_notNormalized() {
	c.DisableNormalize()
	c.DisableComputeSideEffects()
	c.FoldSame("{let x}")
	c.FoldSame("function f() {let x}")
	c.FoldSame("{const x = 1}")
	c.FoldSame("{x = 2; y = 4; let z;}")
	c.Fold("{'hi'; let x;}", "{let x}")
	c.Fold("{x = 4; {let y}}", "x = 4; {let y}")
	c.FoldSame("{class C {}} {class C {}}")
	c.Fold("{label: var x}", "label: var x")
	// `{label: let x}` is a syntax error
	c.FoldSame("{label: var x; let y;}")
}

// Java: PeepholeRemoveDeadCodeTest.testFoldBlockWithDeclaration_normalized()
func (c *PeepholeRemoveDeadCodeTest) TestFoldBlockWithDeclaration_normalized() {
	c.Fold("{let x}", "let x")
	c.FoldSame("function f() {let x}")
	c.Fold("{const x = 1}", "const x = 1;")
	c.Fold("{x = 2; y = 4; let z;}", "x = 2; y = 4; let z;")
	c.Fold("{'hi'; let x;}", "let x;")
	c.Fold("{x = 4; {let y}}", "x = 4; let y;")
	c.Fold("{class C {}} {class C {}}", "class C {} class C$jscomp$1 {}")
	c.Fold("{label: var x}", "label: var x")
	// `{label: let x}` is a syntax error
	c.Fold("{label: var x; let y;}", "label: var x; let y;")
}

// Try to remove spurious blocks with multiple children
// Java: PeepholeRemoveDeadCodeTest.testFoldBlocksWithManyChildren()
func (c *PeepholeRemoveDeadCodeTest) TestFoldBlocksWithManyChildren() {
	c.Fold("function f() { if (false) {} }", "function f(){}")
	c.Fold("function f() { { if (false) {} if (true) {} {} } }", "function f(){}")
	c.Fold(
		"{var x; var y; var z; class Foo { constructor() { var a; { var b; } } } }",
		"var x;var y;var z;class Foo { constructor() { var a;var b} }")
	c.Fold("{var x; var y; var z; { { var a; { var b; } } } }", "var x;var y;var z; var a;var b")
}

// Java: PeepholeRemoveDeadCodeTest.testIf()
func (c *PeepholeRemoveDeadCodeTest) TestIf() {
	c.Fold("if (1){ x=1; } else { x = 2;}", "x=1")
	c.Fold("if (false){ x = 1; } else { x = 2; }", "x=2")
	c.Fold("if (undefined){ x = 1; } else { x = 2; }", "x=2")
	c.Fold("if (null){ x = 1; } else { x = 2; }", "x=2")
	c.Fold("if (void 0){ x = 1; } else { x = 2; }", "x=2")
	c.Fold("if (void foo()){ x = 1; } else { x = 2; }", "foo();x=2")
	c.Fold("if (false){ x = 1; } else if (true) { x = 3; } else { x = 2; }", "x=3")
	c.Fold("if (x){ x = 1; } else if (false) { x = 3; }", "if(x)x=1")
}

// Java: PeepholeRemoveDeadCodeTest.testHook()
func (c *PeepholeRemoveDeadCodeTest) TestHook() {
	c.Fold("true ? a() : b()", "a()")
	c.Fold("false ? a() : b()", "b()")
	c.Fold("a() ? b() : true", "a() && b()")
	c.Fold("a() ? true : b()", "a() || b()")
	c.Fold("(a = true) ? b() : c()", "a = true, b()")
	c.Fold("(a = false) ? b() : c()", "a = false, c()")
	c.Fold("do {f()} while((a = true) ? b() : c())", "do {f()} while((a = true) , b())")
	c.Fold("do {f()} while((a = false) ? b() : c())", "do {f()} while((a = false) , c())")
	c.Fold("var x = (true) ? 1 : 0", "var x=1")
	c.Fold("var y = (true) ? ((false) ? 12 : (cond ? 1 : 2)) : 13", "var y=cond?1:2")
	c.FoldSame("var z=x?void 0:y()")
	c.FoldSame("z=x?void 0:y()")
	c.FoldSame("z*=x?void 0:y()")
	c.FoldSame("var z=x?y():void 0")
	c.FoldSame("(w?x:void 0).y=z")
	c.FoldSame("(w?x:void 0).y+=z")
	c.Fold("y = (x ? void 0 : void 0)", "y = void 0")
	c.Fold("y = (x ? f() : f())", "y = f()")
	c.Fold("(function(){}) ? function(){} : function(){}", "")
}

// Java: PeepholeRemoveDeadCodeTest.testConstantConditionWithSideEffect1()
func (c *PeepholeRemoveDeadCodeTest) TestConstantConditionWithSideEffect1() {
	c.Fold("if (b=true) x=1;", "b=true;x=1")
	c.Fold("if (b=/ab/) x=1;", "b=/ab/;x=1")
	c.Fold("if (b=/ab/){ x=1; } else { x=2; }", "b=/ab/;x=1")
	c.Fold("var b;b=/ab/;if(b)x=1;", "var b;b=/ab/;x=1")
	c.FoldSame("var b;b=f();if(b)x=1;")
	c.Fold("var b=/ab/;if(b)x=1;", "var b=/ab/;x=1")
	c.FoldSame("var b=f();if(b)x=1;")
	c.FoldSame("b=b++;if(b)x=b;")
	c.Fold("(b=0,b=1);if(b)x=b;", "b=0,b=1;if(b)x=b;")
	c.Fold("b=1;if(foo,b)x=b;", "b=1;x=b;")
	c.FoldSame("b=1;if(foo=1,b)x=b;")
}

// Java: PeepholeRemoveDeadCodeTest.testConstantConditionWithSideEffect2()
func (c *PeepholeRemoveDeadCodeTest) TestConstantConditionWithSideEffect2() {
	c.Fold("(b=true)?x=1:x=2;", "b=true,x=1")
	c.Fold("(b=false)?x=1:x=2;", "b=false,x=2")
	c.Fold("if (b=/ab/) x=1;", "b=/ab/;x=1")
	c.Fold("var b;b=/ab/;(b)?x=1:x=2;", "var b;b=/ab/;x=1")
	c.FoldSame("var b;b=f();(b)?x=1:x=2;")
	c.Fold("var b=/ab/;(b)?x=1:x=2;", "var b=/ab/;x=1")
	c.FoldSame("var b=f();(b)?x=1:x=2;")
}

// Java: PeepholeRemoveDeadCodeTest.testConstantConditionWithSideEffect_coalesce()
func (c *PeepholeRemoveDeadCodeTest) TestConstantConditionWithSideEffect_coalesce() {
	c.Fold("b = null; b ?? (x = 1)", "b = null; void 0 ?? (x = 1)")
	c.Fold("b = undefined; b ?? (x = 1)", "b = undefined; void 0 ?? (x = 1)")
	c.Fold("b = (fn(), null); b ?? (x = 1)", "b = (fn(), null); void 0 ?? (x = 1)")
	c.Fold("b = 34; b ?? (x = 1)", "b = 34; 0 ?? (x = 1)")
	c.Fold("b = 'test'; b ?? (x = 1)", "b = 'test'; 0 ?? (x = 1)")
	c.Fold("b = []; b ?? (x = 1)", "b = []; 0 ?? (x = 1)")
	c.Fold("b = (fn(), 0); b ?? (x = 1)", " b= (fn(), 0); 0 ?? (x = 1)")
	c.FoldSame("b = fn(); b ?? (x = 1)")
}

// Java: PeepholeRemoveDeadCodeTest.testVarLifting()
func (c *PeepholeRemoveDeadCodeTest) TestVarLifting() {
	c.Fold("if(true)var a", "var a")
	c.Fold("if(false)var a", "var a")
	// More var lifting tests in PeepholeIntegrationTests
}

// Java: PeepholeRemoveDeadCodeTest.testLetConstLifting()
func (c *PeepholeRemoveDeadCodeTest) TestLetConstLifting() {
	c.Fold("if(true) {const x = 1}", "const x = 1;")
	c.Fold("if(false) {const x = 1}", "")
	c.Fold("if(true) {let x}", "let x;")
	c.Fold("if(false) {let x}", "")
	c.Fold("if(false) {const x = 1;  function f() { return x; }}", "")
}

// Java: PeepholeRemoveDeadCodeTest.testLetConstLifting_removePartOfBlock()
func (c *PeepholeRemoveDeadCodeTest) TestLetConstLifting_removePartOfBlock() {
	c.Fold(
		"function f() {\n"+
			"  return 0;\n"+
			"  let x = 0;\n"+
			"  x++;\n"+
			"}\n",
		"function f() {\n"+
			"  let x;\n"+
			"  return 0;\n"+
			"}\n")
	c.Fold(
		"function f() {\n"+
			"  return 0;\n"+
			"  const [x, y, [[z]]] = 1;\n"+
			"}\n",
		"function f() {\n"+
			"  let x;\n"+
			"  let y;\n"+
			"  let z;\n"+
			"  return 0;\n"+
			"}\n")
	c.Fold(
		"function f() {\n"+
			"  return 0;\n"+
			"  const C = class Bar {};\n"+
			"}\n",
		"function f() {\n"+
			"  let C;\n"+
			"  return 0;\n"+
			"}\n")
}

// Java: PeepholeRemoveDeadCodeTest.testLetConstLifting_removePartOfBlock_withHoistedFunction()
func (c *PeepholeRemoveDeadCodeTest) TestLetConstLifting_removePartOfBlock_withHoistedFunction() {
	c.Fold(
		"function f() {\n"+
			"  return 0;\n"+
			"  // Everything after this is dead code, except for the function declaration, which is\n"+
			"  // hoisted. `foo` could in theory reference x, y, or C. Add a stub 'let' declaration for\n"+
			"  // each to avoid violating the invariant that all NAME nodes in the AST are declared.\n"+
			"  const x = 1;\n"+
			"  let y;\n"+
			"  class C {}\n"+
			"  function foo() {}\n"+
			"}\n",
		"function f() {\n"+
			"  function foo() {}\n"+
			"  let C;\n"+
			"  let y;\n"+
			"  let x;\n"+
			"  return 0;\n"+
			"}\n")
	c.Fold(
		"function f(param) {\n"+
			"  return 0;\n"+
			"  // everything after this is dead code.\n"+
			"  if (param) {\n"+
			"    function foo() {}\n"+
			"    const x = 1;\n"+
			"    let y;\n"+
			"    class C {}\n"+
			"  }\n"+
			"}\n",
		"function f(param) {\n"+
			"  return 0;\n"+
			"}\n")
	c.Fold(
		"function f(param) {\n"+
			"  if (param) {\n"+
			"    return 0;\n"+
			"    const x = 1;\n"+
			"  }\n"+
			"  return 1;\n"+
			"}\n",
		"function f(param) {\n"+
			"  if (param) {\n"+
			"    let x;\n"+
			"    return 0;\n"+
			"  }\n"+
			"  return 1;\n"+
			"}\n")
}

// Java: PeepholeRemoveDeadCodeTest.testFoldUselessFor()
func (c *PeepholeRemoveDeadCodeTest) TestFoldUselessFor() {
	c.Fold("for(;false;) { foo() }", "")
	c.Fold("for(;void 0;) { foo() }", "")
	c.Fold("for(;undefined;) { foo() }", "")
	c.Fold("for(;true;) foo() ", "for(;;) foo() ")
	c.FoldSame("for(;;) foo()")
	c.Fold("for(;false;) { var a = 0; }", "var a")
	c.Fold("for(;false;) { const a = 0; }", "")
	c.Fold("for(;false;) { let a = 0; }", "")
	// Make sure it plays nice with minimizing
	c.Fold("for(;false;) { foo(); continue }", "")
	c.Fold("l1:for(;false;) {  }", "")
}

// Java: PeepholeRemoveDeadCodeTest.testFoldUselessDo()
func (c *PeepholeRemoveDeadCodeTest) TestFoldUselessDo() {
	c.Fold("do { foo() } while(false);", "foo()")
	c.Fold("do { foo() } while(void 0);", "foo()")
	c.Fold("do { foo() } while(undefined);", "foo()")
	c.FoldSame("do { foo() } while(true);")
	c.Fold("do { var a = 0; } while(false);", "var a=0")
	c.Fold("do { var a = 0; } while(!{a:foo()});", "var a=0;foo()")
	// Can't fold with break or continues.
	c.FoldSame("do { foo(); continue; } while(0)")
	c.FoldSame("do { try { foo() } catch (e) { break; } } while (0);")
	c.FoldSame("do { foo(); break; } while(0)")
	c.Fold("do { for (;;) {foo(); continue;} } while(0)", "for (;;) {foo(); continue;}")
	c.FoldSame("l1: do { for (;;) { foo() } } while(0)")
	c.Fold("do { switch (1) { default: foo(); break} } while(0)", "foo();")
	c.Fold(
		"do { switch (1) { default: foo(); continue} } while(0)",
		"do { foo(); continue } while(0)")
	c.Fold(
		"l1: { do { x = 1; break l1; } while (0); x = 2; }", //
		"l1: { x = 1; break l1; }")
	c.Fold("do { x = 1; } while (x = 0);", "x = 1; x = 0;")
	c.Fold(
		"let x = 1; (function() { do { let x = 2; } while (x = 10, false); })();",
		"let x = 1; (function() { let x$jscomp$1 = 2; x = 10 })();")
}

// Java: PeepholeRemoveDeadCodeTest.testFoldEmptyDo()
func (c *PeepholeRemoveDeadCodeTest) TestFoldEmptyDo() {
	c.Fold("do { } while(true);", "for (;;);")
}

// Java: PeepholeRemoveDeadCodeTest.testMinimizeLoop_withConstantCondition_vanillaFor()
func (c *PeepholeRemoveDeadCodeTest) TestMinimizeLoop_withConstantCondition_vanillaFor() {
	c.Fold("for(;true;) foo()", "for(;;) foo()")
	c.Fold("for(;0;) foo()", "")
	c.Fold("for(;0.0;) foo()", "")
	c.Fold("for(;NaN;) foo()", "")
	c.Fold("for(;null;) foo()", "")
	c.Fold("for(;undefined;) foo()", "")
	c.Fold("for(;'';) foo()", "")
}

// Java: PeepholeRemoveDeadCodeTest.testMinimizeLoop_withConstantCondition_doWhile()
func (c *PeepholeRemoveDeadCodeTest) TestMinimizeLoop_withConstantCondition_doWhile() {
	c.Fold("do { foo(); } while (true)", "do { foo(); } while (true);")
	c.Fold("do { foo(); } while (0)", "foo();")
	c.Fold("do { foo(); } while (0.0)", "foo();")
	c.Fold("do { foo(); } while (NaN)", "foo();")
	c.Fold("do { foo(); } while (null)", "foo();")
	c.Fold("do { foo(); } while (undefined)", "foo();")
	c.Fold("do { foo(); } while ('')", "foo();")
}

// Java: PeepholeRemoveDeadCodeTest.testFoldConstantCommaExpressions()
func (c *PeepholeRemoveDeadCodeTest) TestFoldConstantCommaExpressions() {
	c.Fold("if (true, false) {foo()}", "")
	c.Fold("if (false, true) {foo()}", "foo()")
	c.Fold("true, foo()", "foo()")
	c.Fold("true, foo?.()", "foo?.()")
	c.Fold("(1 + 2 + ''), foo()", "foo()")
	c.Fold("(1 + 2 + ''), foo?.()", "foo?.()")
}

// Java: PeepholeRemoveDeadCodeTest.testRemoveUselessOps1()
func (c *PeepholeRemoveDeadCodeTest) TestRemoveUselessOps1() {
	c.FoldSame("(function () { f(); })();")
}

// TODO(b/538125186): Do not remove optional chain side effects in PeepholeRemoveDeadCode
// Java: PeepholeRemoveDeadCodeTest.testCallSideEffectsPreserved()
func (c *PeepholeRemoveDeadCodeTest) TestCallSideEffectsPreserved() {
	// Functions calls known to be free of side effects are removed.
	c.Fold("Math.random()", "")
	c.Fold("Math?.random()", "")
	c.Fold("Math.random(f() + g())", "f(),g();")
	c.FoldSame("Math?.random(f() + g());")
	c.Fold("Math.random(f(),g(),h())", "f(),g(),h();")
	c.FoldSame("Math?.random(f(),g(),h());")
	c.Fold("Symbol()", "")
	c.Fold("Symbol('desc')", "")
	c.Fold("Symbol(f())", "f();")

	// Calls to functions with unknown side-effects are preserved.
	c.FoldSame("f();")
	c.FoldSame("f?.();")
	c.FoldSame("(function () { f(); })();")
}

// Java: PeepholeRemoveDeadCodeTest.testRemoveUselessOps2()
func (c *PeepholeRemoveDeadCodeTest) TestRemoveUselessOps2() {
	// There are four place where expression results are discarded:
	//  - a top-level expression EXPR_RESULT
	//  - the LHS of a COMMA
	//  - the FOR init expression
	//  - the FOR increment expression
	// We know that this function has no side effects because of the
	// PureFunctionIdentifier.
	c.Fold("(function () {})();", "")
	// Uncalled function expressions are removed
	c.Fold("(function () {});", "")
	c.Fold("(function f() {});", "")
	c.Fold("(function* f() {})", "")
	// ... including any code they contain.
	c.Fold("(function () {foo();});", "")
	// Useless operators are removed.
	c.Fold("+f()", "f()")
	c.Fold("+f?.()", "f?.()")
	c.Fold("a=(+f(),g())", "a=(f(),g())")
	c.Fold("a=(+f?.(),g())", "a=(f?.(),g())")
	c.Fold("a=(true,g())", "a=g()")
	c.Fold("f(),true", "f()")
	c.Fold("f() + g()", "f(),g()")
	c.Fold("for(;;+f()){}", "for(;;f()){}")
	c.Fold("for(+f();;g()){}", "for(f();;g()){}")
	c.Fold("for(;;Math.random(f(),g(),h())){}", "for(;;f(),g(),h()){}")
	// The optimization cascades into conditional expressions:
	c.Fold("g() && +f()", "g() && f()")
	c.Fold("g() || +f()", "g() || f()")
	c.Fold("x ? g() : +f()", "x ? g() : f()")
	c.Fold("+x()", "x()")
	c.Fold("+x() * 2", "x()")
	c.Fold("-(+x() * 2)", "x()")
	c.Fold("2 -(+x() * 2)", "x()")
	c.Fold("x().foo", "x()")
	c.FoldSame("x().foo()")
	c.FoldSame("x++")
	c.FoldSame("++x")
	c.FoldSame("x--")
	c.FoldSame("--x")
	c.FoldSame("x = 2")
	c.FoldSame("x *= 2")
	// Sanity check, other expression are left alone.
	c.FoldSame("function f() {}")
	c.FoldSame("var x;")
}

// TODO(b/538166699): Do not remove switch condition with side effects in PeepholeRemoveDeadCode
// Java: PeepholeRemoveDeadCodeTest.testOptimizeSwitch()
func (c *PeepholeRemoveDeadCodeTest) TestOptimizeSwitch() {
	c.Fold("switch(a){}", "")
	c.Fold("switch(foo()){}", "foo()")
	c.Fold("switch(x++){}", "x++")
	c.Fold("switch(a=b){}", "a=b")
	c.Fold("switch(x++){default:foo();}", "x++; foo()")
	c.Fold("switch(a=b){default:foo();}", "a=b; foo()")
	c.Fold("switch(a){default:}", "")
	c.Fold("switch(a){default:break;}", "")
	c.Fold("switch(a){default:var b;break;}", "var b")
	c.Fold("switch(a){case 1: default:}", "")
	c.Fold("switch(a){default: case 1:}", "")
	c.Fold("switch(a){default: break; case 1:break;}", "")
	c.Fold("switch(a){default: var b; break; case 1: var c; break;}", "var c; var b;")
	c.Fold("var x=1; switch(x) { case 1: var y; }", "var y; var x=1;")

	// Can't remove cases if a default exists and is not the last case.
	c.FoldSame("function f() {switch(a){default: return; case 1: break;}}")
	c.FoldSame("function f() {switch(1){default: return; case 1: break;}}") // foldable
	c.FoldSame("function f() {switch(a){case 1: foo();}}")
	c.FoldSame("function f() {switch(a){case 3: case 2: case 1: foo();}}")

	c.Fold("function f() {switch(a){case 2: case 1: default: foo();}}", "function f() { foo(); }")
	c.Fold("switch(a){case 1: default:break; case 2: foo()}", "switch(a){case 2: foo()}")
	c.FoldSame("switch(a){case 1: goo(); default:break; case 2: foo()}")

	// TODO(johnlenz): merge the useless "case 2"
	c.FoldSame("switch(a){case 1: goo(); case 2:break; case 3: foo()}")

	// Can't remove unused code with a "var" in it.
	c.Fold("switch(1){case 2: var x=0;}", "var x;")
	c.Fold(
		"switch ('repeated') {\n"+
			"case 'repeated':\n"+
			"  foo();\n"+
			"  break;\n"+
			"case 'repeated':\n"+
			"  var x=0;\n"+
			"  break;\n"+
			"}\n",
		"var x; foo();")

	// Can't remove cases if something useful is done.
	c.FoldSame("switch(a){case 1: var c =2; break;}")
	c.FoldSame("function f() {switch(a){case 1: return;}}")
	c.FoldSame("x:switch(a){case 1: break x;}")

	c.Fold(
		"switch ('foo') {\n"+
			"case 'foo':\n"+
			"  foo();\n"+
			"  break;\n"+
			"case 'bar':\n"+
			"  bar();\n"+
			"  break;\n"+
			"}\n",
		"foo();")
	c.Fold(
		"switch ('noMatch') {\n"+
			"case 'foo':\n"+
			"  foo();\n"+
			"  break;\n"+
			"case 'bar':\n"+
			"  bar();\n"+
			"  break;\n"+
			"}\n",
		"")
	c.Fold(
		"switch ('fallThru') {\n"+
			"case 'fallThru':\n"+
			"  if (foo(123) > 0) {\n"+
			"    foobar(1);\n"+
			"    break;\n"+
			"  }\n"+
			"  foobar(2);\n"+
			"case 'bar':\n"+
			"  bar();\n"+
			"}\n",
		"switch ('fallThru') {\n"+
			"case 'fallThru':\n"+
			"  if (foo(123) > 0) {\n"+
			"    foobar(1);\n"+
			"    break;\n"+
			"  }\n"+
			"  foobar(2);\n"+
			"  bar();\n"+
			"}\n")
	c.Fold(
		"switch ('fallThru') {\n"+
			"case 'fallThru':\n"+
			"  foo();\n"+
			"case 'bar':\n"+
			"  bar();\n"+
			"}\n",
		"foo();\n"+
			"bar();\n")
	c.Fold(
		"switch ('hasDefaultCase') {\n"+
			"  case 'foo':\n"+
			"    foo();\n"+
			"    break;\n"+
			"  default:\n"+
			"    bar();\n"+
			"    break;\n"+
			"}\n",
		"bar();")
	c.Fold(
		"switch ('repeated') {\n"+
			"case 'repeated':\n"+
			"  foo();\n"+
			"  break;\n"+
			"case 'repeated':\n"+
			"  bar();\n"+
			"  break;\n"+
			"}\n",
		"foo();")
	c.FoldSame(
		"switch ('foo') {\n" +
			"case 'bar':\n" +
			"  bar();\n" +
			"  break;\n" +
			"case notConstant:\n" +
			"  foobar();\n" +
			"  break;\n" +
			"case 'foo':\n" +
			"  foo();\n" +
			"  break;\n" +
			"}\n")
	c.Fold(
		"switch (1) {\n"+
			"case 1:\n"+
			"  foo();\n"+
			"  break;\n"+
			"case 2:\n"+
			"  bar();\n"+
			"  break;\n"+
			"}\n",
		"foo();")
	c.Fold(
		"switch (1) {\n"+
			"case 1.1:\n"+
			"  foo();\n"+
			"  break;\n"+
			"case 2:\n"+
			"  bar();\n"+
			"  break;\n"+
			"}\n",
		"")
	c.Fold(
		"switch (0) {\n"+
			"case NaN:\n"+
			"  foobar();\n"+
			"  break;\n"+
			"case -0.0:\n"+
			"  foo();\n"+
			"  break;\n"+
			"case 2:\n"+
			"  bar();\n"+
			"  break;\n"+
			"}\n",
		"foo();")
	c.FoldSame(
		"switch ('\\v') {\n" +
			"case '\\u000B':\n" +
			"  foo();\n" +
			"}\n")
	c.Fold(
		"switch ('empty') {\n"+
			"case 'empty':\n"+
			"case 'foo':\n"+
			"  foo();\n"+
			"}\n",
		"foo()")

	c.Fold(
		"let x;\n"+
			"switch (use(x)) {\n"+
			"  default: {let y;}\n"+
			"}\n",
		"let x;\n"+
			"use(x);\n"+
			"let y;\n")

	c.Fold(
		"let x;\n"+
			"switch (use?.(x)) {\n"+
			"  default: {let y;}\n"+
			"}\n",
		"let x;\n"+
			"use?.(x);\n"+
			"let y;\n")

	c.Fold(
		"let x;\n"+
			"switch (use(x)) {\n"+
			"  default: let y;\n"+
			"}\n",
		"let x;\n"+
			"use(x);\n"+
			"let y;\n")
}

// Java: PeepholeRemoveDeadCodeTest.testOptimizeSwitchBug335145701()
func (c *PeepholeRemoveDeadCodeTest) TestOptimizeSwitchBug335145701() {
	c.FoldSame(
		"function foo() { alert('foo()'); }\n" +
			"switch (1) {\n" +
			"  case 1: break;\n" +
			"  case foo(): break;\n" +
			"}\n")
	c.FoldSame(
		"function foo() { alert('foo()'); }\n" +
			"switch (1) {\n" +
			"  case 0: break;\n" +
			"  case 1: break;\n" +
			"  case foo(): break;\n" +
			"}\n")
	c.FoldSame(
		"function foo() { alert('foo()'); }\n" +
			"switch (1) {\n" +
			"  case 0: alert('bar'); break;\n" +
			"  case 1: break;\n" +
			"  case foo(): break;\n" +
			"}\n")
	c.FoldSame(
		"function foo() { alert('foo()'); return 1; }\n" +
			"switch (1) {\n" +
			"  case 0: break;\n" +
			"  case foo(): break;\n" +
			"  case 2: break;\n" +
			"}\n")
	c.Fold(
		"function foo() { alert('foo()'); }\n"+
			"switch (1) {\n"+
			"  case foo(): break;\n"+
			"  case (0,1): break;\n"+
			"}\n",
		"function foo() { alert('foo()'); }\n"+
			"switch (1) {\n"+
			"  case foo(): break;\n"+
			"  case 1: break;\n"+
			"}\n")
	c.FoldSame(
		"function foo() { alert('foo()'); }\n" +
			"switch (x) {\n" +
			"  case 1: break;\n" +
			"  case foo(): break;\n" +
			"}\n")
	c.Fold(
		"// not valid to remove the useless case 1,\n"+
			"// it would cause the default to run and it has side-effects\n"+
			"switch (1) {\n"+
			"  case 1: break;\n"+
			"  default:\n"+
			"    bar();\n"+
			"    break;\n"+
			"}\n",
		"")
	c.FoldSame(
		"function foo() { alert('foo()'); }\n" +
			"switch (1) {\n" +
			"  case 0: alert('bar'); break;\n" +
			"  case 1: break;\n" +
			"  case foo(): break;\n" +
			"}\n")
	c.FoldSame(
		"    function foo() { alert('foo()'); }\n" +
			"    switch (bar()) {\n" +
			"      case 1: break;\n" +
			"      case foo(): break;\n" +
			"    }\n")
	c.Fold(
		"// is not valid to remove the first useless case 1,\n"+
			"// because it matches and the second should not run\n"+
			"switch (1) {\n"+
			"  case 1: break;\n"+
			"  case 1: bar(); break;\n"+
			"}\n",
		"")
}

// Java: PeepholeRemoveDeadCodeTest.testOptimizeSwitchBug11536863()
func (c *PeepholeRemoveDeadCodeTest) TestOptimizeSwitchBug11536863() {
	c.Fold(
		"outer: {\n"+
			"  switch (2) {\n"+
			"    case 2:\n"+
			"      f();\n"+
			"      break outer;\n"+
			"  }\n"+
			"}\n",
		"outer: {f(); break outer;}")
}

// `a[b]` could trigger a getter or setter, and have side effects. However, we always assume it
// does not (even though it's unsound) because the code size cost of assuming all GETELEM nodes
// have side effects is unacceptable.
// Java: PeepholeRemoveDeadCodeTest.testUnusedGetElemRemoved()
func (c *PeepholeRemoveDeadCodeTest) TestUnusedGetElemRemoved() {
	c.Fold("a[b]", "")
	c.Fold("a?.[b]", "")
}

// Java: PeepholeRemoveDeadCodeTest.testOptimizeSwitch2()
func (c *PeepholeRemoveDeadCodeTest) TestOptimizeSwitch2() {
	c.Fold(
		"outer: switch (2) {\n"+
			"  case 2:\n"+
			"    f();\n"+
			"    break outer;\n"+
			"}\n",
		"outer: {f(); break outer;}")
}

// Java: PeepholeRemoveDeadCodeTest.testOptimizeSwitch3()
func (c *PeepholeRemoveDeadCodeTest) TestOptimizeSwitch3() {
	c.Fold(
		"switch (1) {\n"+
			"  case 1:\n"+
			"  case 2:\n"+
			"  case 3: {\n"+
			"    break;\n"+
			"  }\n"+
			"  case 4:\n"+
			"  case 5:\n"+
			"  case 6:\n"+
			"  default:\n"+
			"    fail('Should not get here');\n"+
			"    break;\n"+
			"}\n",
		"")
}

// Java: PeepholeRemoveDeadCodeTest.testOptimizeSwitchWithLabellessBreak()
func (c *PeepholeRemoveDeadCodeTest) TestOptimizeSwitchWithLabellessBreak() {
	c.Fold(
		"function f() {\n"+
			"  switch('x') {\n"+
			"    case 'x': var x = 1; break;\n"+
			"    case 'y': break;\n"+
			"  }\n"+
			"}\n",
		"function f() { var x = 1; }")
	// TODO(moz): Convert this to an if statement for better optimization
	c.FoldSame(
		"function f() {\n" +
			"  switch(x) {\n" +
			"    case 'y': break;\n" +
			"    default: var x = 1;\n" +
			"  }\n" +
			"}\n")
	c.Fold(
		"var exit;\n"+
			"switch ('a') {\n"+
			"  case 'a':\n"+
			"    break;\n"+
			"  default:\n"+
			"    exit = 21;\n"+
			"    break;\n"+
			"}\n"+
			"switch(exit) {\n"+
			"  case 21: throw 'x';\n"+
			"  default : console.log('good');\n"+
			"}\n",
		"var exit;\n"+
			"switch(exit) {\n"+
			"  case 21: throw 'x';\n"+
			"  default : console.log('good');\n"+
			"}\n")
	c.Fold(
		"let x = 1;\n"+
			"switch('x') {\n"+
			"  case 'x': let x = 2; break;\n"+
			"}\n",
		"let x = 1;\n"+
			"let x$jscomp$1 = 2\n")
}

// Java: PeepholeRemoveDeadCodeTest.testOptimizeSwitchWithLabelledBreak()
func (c *PeepholeRemoveDeadCodeTest) TestOptimizeSwitchWithLabelledBreak() {
	c.Fold(
		"function f() {\n"+
			"  label:\n"+
			"  switch('x') {\n"+
			"    case 'x': break label;\n"+
			"    case 'y': throw f;\n"+
			"  }\n"+
			"}\n",
		"function f() { }")
	c.Fold(
		"function f() {\n"+
			"  label:\n"+
			"  switch('x') {\n"+
			"    case 'x': break label;\n"+
			"    default: throw f;\n"+
			"  }\n"+
			"}\n",
		"function f() { }")
}

// Java: PeepholeRemoveDeadCodeTest.testOptimizeSwitchWithReturn()
func (c *PeepholeRemoveDeadCodeTest) TestOptimizeSwitchWithReturn() {
	c.Fold(
		"function f() {\n"+
			"  switch('x') {\n"+
			"    case 'x': return 1;\n"+
			"    case 'y': return 2;\n"+
			"  }\n"+
			"}\n",
		"function f() { return 1; }")
	c.Fold(
		"function f() {\n"+
			"  let x = 1;\n"+
			"  switch('x') {\n"+
			"    case 'x': { let x = 2; } return 3;\n"+
			"    case 'y': return 4;\n"+
			"  }\n"+
			"}\n",
		"function f() {\n"+
			"  let x = 1;\n"+
			"  let x$jscomp$1 = 2;\n"+
			"  return 3;\n"+
			"}\n")
}

// Java: PeepholeRemoveDeadCodeTest.testOptimizeSwitchWithThrow()
func (c *PeepholeRemoveDeadCodeTest) TestOptimizeSwitchWithThrow() {
	c.Fold(
		"function f() {\n"+
			"  switch('x') {\n"+
			"    case 'x': throw f;\n"+
			"    case 'y': throw f;\n"+
			"  }\n"+
			"}\n",
		"function f() { throw f; }")
}

// Java: PeepholeRemoveDeadCodeTest.testOptimizeSwitchWithContinue()
func (c *PeepholeRemoveDeadCodeTest) TestOptimizeSwitchWithContinue() {
	c.Fold(
		"function f() {\n"+
			"  for (;;) {\n"+
			"    switch('x') {\n"+
			"      case 'x': continue;\n"+
			"      case 'y': continue;\n"+
			"    }\n"+
			"  }\n"+
			"}\n",
		"function f() { for (;;) { continue; } }")
}

// Java: PeepholeRemoveDeadCodeTest.testOptimizeSwitchWithDefaultCaseWithFallthru()
func (c *PeepholeRemoveDeadCodeTest) TestOptimizeSwitchWithDefaultCaseWithFallthru() {
	c.FoldSame(
		"function f() {\n" +
			"  switch(a) {\n" +
			"    case 'x':\n" +
			"    case foo():\n" +
			"    default: return 3\n" +
			"  }\n" +
			"}\n")
}

// GitHub issue #1722: https://github.com/google/closure-compiler/issues/1722
// Java: PeepholeRemoveDeadCodeTest.testOptimizeSwitchWithDefaultCase()
func (c *PeepholeRemoveDeadCodeTest) TestOptimizeSwitchWithDefaultCase() {
	c.Fold(
		"function f() {\n"+
			"  switch('x') {\n"+
			"    case 'x': return 1;\n"+
			"    case 'y': return 2;\n"+
			"    default: return 3\n"+
			" }\n"+
			"}\n",
		"function f() { return 1; }")
	c.Fold(
		"switch ('hasDefaultCase') {\n"+
			"  case 'foo':\n"+
			"    foo();\n"+
			"    break;\n"+
			"  default:\n"+
			"    bar();\n"+
			"    break;\n"+
			"}\n",
		"bar();")
	c.FoldSame("switch (x) { default: if (a) { break; } bar(); }")
	// Potentially foldable
	c.FoldSame(
		"switch (x) {\n" +
			"  case x:\n" +
			"    foo();\n" +
			"    break;\n" +
			"  default:\n" +
			"    if (a) { break; }\n" +
			"    bar();\n" +
			"}\n")
	c.Fold(
		"switch ('hasDefaultCase') {\n"+
			"  case 'foo':\n"+
			"    foo();\n"+
			"    break;\n"+
			"  default:\n"+
			"    if (true) { break; }\n"+
			"    bar();\n"+
			"}\n",
		"")
	c.Fold(
		"switch ('hasDefaultCase') {\n"+
			"  case 'foo':\n"+
			"    foo();\n"+
			"    break;\n"+
			"  default:\n"+
			"    if (a) { break; }\n"+
			"    bar();\n"+
			"}\n",
		"switch ('hasDefaultCase') { default: if (a) { break; } bar(); }")
	c.Fold(
		"l: switch ('hasDefaultCase') {\n"+
			"  case 'foo':\n"+
			"    foo();\n"+
			"    break;\n"+
			"  default:\n"+
			"    if (a) { break l; }\n"+
			"    bar();\n"+
			"    break;\n"+
			"}\n",
		"l:{ if (a) { break l; } bar(); }")
	c.Fold(
		"switch ('hasDefaultCase') {\n"+
			"  case 'foo':\n"+
			"    bar();\n"+
			"    break;\n"+
			"  default:\n"+
			"    foo();\n"+
			"    break;\n"+
			"}\n",
		"foo();")
	c.Fold("switch (a()) { default: bar(); break;}", "a(); bar();")
	c.Fold("switch (a?.()) { default: bar(); break;}", "a?.(); bar();")
	c.Fold("switch (a()) { default: break; bar();}", "a();")
	c.Fold(
		"loop:\n"+
			"for (;;) {\n"+
			"  switch (a()) {\n"+
			"    default:\n"+
			"      bar();\n"+
			"      break loop;\n"+
			"  }\n"+
			"}\n",
		"loop: for (;;) { a(); bar(); break loop; }")
}

// Java: PeepholeRemoveDeadCodeTest.testTreatSwitchAsExit()
func (c *PeepholeRemoveDeadCodeTest) TestTreatSwitchAsExit() {
	c.Fold(
		"a: {switch(x){ case 1: case 2: break a; default: foo(); break a;} bar(); }",
		"a: {switch(x){ case 1: case 2: break a; default: foo(); break a;} }")
}

// Java: PeepholeRemoveDeadCodeTest.testDontTreatSwitchAsExit()
func (c *PeepholeRemoveDeadCodeTest) TestDontTreatSwitchAsExit() {
	c.FoldSame("a: {switch(x){ case 1:break a; default: b: {foo(); break b;}} bar(); }")
	c.FoldSame("a: {switch(x){ case 1:break a; default: foo(); break;} bar(); }")
	c.FoldSame("a: {b: switch(x){ case 1:break a; default: foo(); break b;} bar(); }")
	c.FoldSame("a: {switch(x){ case 1: b: { foo(); break b; } default: break a; } bar(); }")
	c.FoldSame("a: {switch(x){ case 1: break a; } bar(); }")
	c.FoldSame("a: {switch(x){ case 1: if (y) { break; } default: break a;} bar(); }")
	c.FoldSame("a: {switch(x){ case 1: if (y) { break; } break a; default: break a;} bar(); }")
}

// Java: PeepholeRemoveDeadCodeTest.testTreatTryAsExit()
func (c *PeepholeRemoveDeadCodeTest) TestTreatTryAsExit() {
	c.Fold(
		"a: {try { foo(); break a; } catch (e) { foo(); break a; } bar(); }",
		"a: {try { foo(); break a; } catch (e) { foo(); break a; } }")
	c.Fold(
		"a: {try { foo(); } finally { foo(); break a; } bar(); }",
		"a: {try { foo(); } finally { foo(); break a; } }")
	c.Fold(
		"a: {try { foo(); break a; } finally { foo(); } bar(); }",
		"a: {try { foo(); break a; } finally { foo(); } }")
}

// Java: PeepholeRemoveDeadCodeTest.testDontTreatTryAsExit()
func (c *PeepholeRemoveDeadCodeTest) TestDontTreatTryAsExit() {
	c.FoldSame("a: {try { foo(); break a; } catch (e) { foo(); } bar(); }")
	c.FoldSame("a: {try { foo(); break a; } catch (e) { } bar(); }")
	c.FoldSame("a: {try { foo(); } catch (e) { foo(); break a; } bar(); }")
	c.FoldSame("a: {try { foo(); } finally { foo(); } bar(); }")
	c.FoldSame("a: {try { b: { foo(); break b; } } finally { foo(); } bar(); }")
	c.FoldSame("a: { b: try { foo(); break a; } finally { foo(); } bar(); }")
}

// Java: PeepholeRemoveDeadCodeTest.testRemoveNumber()
func (c *PeepholeRemoveDeadCodeTest) TestRemoveNumber() {
	c.Fold("3", "")
}

// Java: PeepholeRemoveDeadCodeTest.testRemoveVarGet1()
func (c *PeepholeRemoveDeadCodeTest) TestRemoveVarGet1() {
	c.Fold("a", "")
}

// Java: PeepholeRemoveDeadCodeTest.testRemoveVarGet2()
func (c *PeepholeRemoveDeadCodeTest) TestRemoveVarGet2() {
	c.Fold("var a = 1;a", "var a = 1")
}

// Java: PeepholeRemoveDeadCodeTest.testRemoveUnusedGetProp()
func (c *PeepholeRemoveDeadCodeTest) TestRemoveUnusedGetProp() {
	c.Fold("var a = {};a.b", "var a = {}")
}

// Java: PeepholeRemoveDeadCodeTest.testRemoveUnusedOptChainGetProp()
func (c *PeepholeRemoveDeadCodeTest) TestRemoveUnusedOptChainGetProp() {
	c.Fold("var a = {};a?.b", "var a = {}")
}

// Java: PeepholeRemoveDeadCodeTest.testRemoveUnusedGetProp2()
func (c *PeepholeRemoveDeadCodeTest) TestRemoveUnusedGetProp2() {
	c.Fold("var a = {};a.b=1;a.b", "var a = {};a.b=1")
}

// Java: PeepholeRemoveDeadCodeTest.testRemoveUnusedOptChainGetProp2()
func (c *PeepholeRemoveDeadCodeTest) TestRemoveUnusedOptChainGetProp2() {
	c.Fold("var a = {};a.b=1;a?.b", "var a = {};a.b=1")
}

// Java: PeepholeRemoveDeadCodeTest.testRemovePrototypeGet1()
func (c *PeepholeRemoveDeadCodeTest) TestRemovePrototypeGet1() {
	c.Fold("var a = {};a.prototype.b", "var a = {}")
}

// Java: PeepholeRemoveDeadCodeTest.testRemoveOptChainPrototypeGet1()
func (c *PeepholeRemoveDeadCodeTest) TestRemoveOptChainPrototypeGet1() {
	c.Fold("var a = {};a?.prototype.b", "var a = {}")
}

// Java: PeepholeRemoveDeadCodeTest.testRemovePrototypeGet2()
func (c *PeepholeRemoveDeadCodeTest) TestRemovePrototypeGet2() {
	c.Fold("var a = {};a.prototype.b = 1;a.prototype.b", "var a = {};a.prototype.b = 1")
}

// Java: PeepholeRemoveDeadCodeTest.testRemoveOptChainPrototypeGet2()
func (c *PeepholeRemoveDeadCodeTest) TestRemoveOptChainPrototypeGet2() {
	c.Fold("var a = {};a.prototype.b = 1;a?.prototype.b", "var a = {};a.prototype.b = 1")
}

// Java: PeepholeRemoveDeadCodeTest.testNotRemovePrototypeGet2()
func (c *PeepholeRemoveDeadCodeTest) TestNotRemovePrototypeGet2() {
	c.FoldSame("var a = {};a.prototype.b = 1; let x = a.prototype.b")
}

// Java: PeepholeRemoveDeadCodeTest.testNotRemoveOptChainPrototypeGet2()
func (c *PeepholeRemoveDeadCodeTest) TestNotRemoveOptChainPrototypeGet2() {
	c.FoldSame("var a = {};a.prototype.b = 1; let x = a?.prototype.b")
}

// Java: PeepholeRemoveDeadCodeTest.testRemoveAdd1()
func (c *PeepholeRemoveDeadCodeTest) TestRemoveAdd1() {
	c.Fold("1 + 2", "")
}

// Java: PeepholeRemoveDeadCodeTest.testNoRemoveVar1()
func (c *PeepholeRemoveDeadCodeTest) TestNoRemoveVar1() {
	c.FoldSame("var a = 1")
}

// Java: PeepholeRemoveDeadCodeTest.testNoRemoveVar2()
func (c *PeepholeRemoveDeadCodeTest) TestNoRemoveVar2() {
	c.FoldSame("var a = 1, b = 2")
}

// Java: PeepholeRemoveDeadCodeTest.testNoRemoveAssign1()
func (c *PeepholeRemoveDeadCodeTest) TestNoRemoveAssign1() {
	c.FoldSame("a = 1")
}

// Java: PeepholeRemoveDeadCodeTest.testNoRemoveAssign2()
func (c *PeepholeRemoveDeadCodeTest) TestNoRemoveAssign2() {
	c.FoldSame("a = b = 1")
}

// Java: PeepholeRemoveDeadCodeTest.testNoRemoveAssign3()
func (c *PeepholeRemoveDeadCodeTest) TestNoRemoveAssign3() {
	c.Fold("1 + (a = 2)", "a = 2")
}

// Java: PeepholeRemoveDeadCodeTest.testNoRemoveAssign4()
func (c *PeepholeRemoveDeadCodeTest) TestNoRemoveAssign4() {
	c.FoldSame("x.a = 1")
}

// Java: PeepholeRemoveDeadCodeTest.testNoRemoveAssign5()
func (c *PeepholeRemoveDeadCodeTest) TestNoRemoveAssign5() {
	c.FoldSame("x.a = x.b = 1")
}

// Java: PeepholeRemoveDeadCodeTest.testNoRemoveAssign6()
func (c *PeepholeRemoveDeadCodeTest) TestNoRemoveAssign6() {
	c.Fold("1 + (x.a = 2)", "x.a = 2")
}

// Java: PeepholeRemoveDeadCodeTest.testNoRemoveCall1()
func (c *PeepholeRemoveDeadCodeTest) TestNoRemoveCall1() {
	c.FoldSame("a()")
}

// Java: PeepholeRemoveDeadCodeTest.testNoRemoveOptChainCall1()
func (c *PeepholeRemoveDeadCodeTest) TestNoRemoveOptChainCall1() {
	c.FoldSame("a?.()")
}

// Java: PeepholeRemoveDeadCodeTest.testNoRemoveCall2()
func (c *PeepholeRemoveDeadCodeTest) TestNoRemoveCall2() {
	c.Fold("a()+b()", "a(),b()")
}

// Java: PeepholeRemoveDeadCodeTest.testNoRemoveOptChainCall2()
func (c *PeepholeRemoveDeadCodeTest) TestNoRemoveOptChainCall2() {
	c.Fold("a?.()+b?.()", "a?.(),b?.()")
}

// Java: PeepholeRemoveDeadCodeTest.testNoRemoveCall3()
func (c *PeepholeRemoveDeadCodeTest) TestNoRemoveCall3() {
	c.FoldSame("a() && b()")
}

// Java: PeepholeRemoveDeadCodeTest.testNoRemoveOptChainCall3()
func (c *PeepholeRemoveDeadCodeTest) TestNoRemoveOptChainCall3() {
	c.FoldSame("a?.() && b?.()")
}

// Java: PeepholeRemoveDeadCodeTest.testNoRemoveCall4()
func (c *PeepholeRemoveDeadCodeTest) TestNoRemoveCall4() {
	c.FoldSame("a() || b()")
}

// Java: PeepholeRemoveDeadCodeTest.testNoRemoveOptChainCall4()
func (c *PeepholeRemoveDeadCodeTest) TestNoRemoveOptChainCall4() {
	c.FoldSame("a?.() || b?.()")
}

// Java: PeepholeRemoveDeadCodeTest.testNoRemoveCall4NullishCoalesce()
func (c *PeepholeRemoveDeadCodeTest) TestNoRemoveCall4NullishCoalesce() {
	c.FoldSame("a() ?? b()")
}

// Java: PeepholeRemoveDeadCodeTest.testNoRemoveOptChainCall4NullishCoalesce()
func (c *PeepholeRemoveDeadCodeTest) TestNoRemoveOptChainCall4NullishCoalesce() {
	c.FoldSame("a?.() ?? b?.()")
}

// Java: PeepholeRemoveDeadCodeTest.testNoRemoveCall5NullishCoalesce()
func (c *PeepholeRemoveDeadCodeTest) TestNoRemoveCall5NullishCoalesce() {
	c.Fold("a() ?? 1", "a()")
}

// Java: PeepholeRemoveDeadCodeTest.testNoRemoveOptChainCall5NullishCoalesce()
func (c *PeepholeRemoveDeadCodeTest) TestNoRemoveOptChainCall5NullishCoalesce() {
	c.Fold("a?.() ?? 1", "a?.()")
}

// Java: PeepholeRemoveDeadCodeTest.testNoRemoveOptionalChainSideEffects()
func (c *PeepholeRemoveDeadCodeTest) TestNoRemoveOptionalChainSideEffects() {
	c.FoldSame("a?.[f()];")
	c.FoldSame("a?.[x++];")
	c.FoldSame("a?.[delete cache[k]];")
	c.FoldSame("a?.b[f()];")
	c.FoldSame("Math?.sin(mutate());")
}

// Java: PeepholeRemoveDeadCodeTest.testNoRemoveCall6NullishCoalesce()
func (c *PeepholeRemoveDeadCodeTest) TestNoRemoveCall6NullishCoalesce() {
	c.FoldSame("1 ?? a()")
}

// Java: PeepholeRemoveDeadCodeTest.testNoRemoveCall5()
func (c *PeepholeRemoveDeadCodeTest) TestNoRemoveCall5() {
	c.Fold("a() || 1", "a()")
}

// Java: PeepholeRemoveDeadCodeTest.testNoRemoveCall6()
func (c *PeepholeRemoveDeadCodeTest) TestNoRemoveCall6() {
	c.FoldSame("1 || a()")
}

// Java: PeepholeRemoveDeadCodeTest.testNoRemoveThrow1()
func (c *PeepholeRemoveDeadCodeTest) TestNoRemoveThrow1() {
	c.FoldSame("function f(){throw a()}")
}

// Java: PeepholeRemoveDeadCodeTest.testNoRemoveThrow2()
func (c *PeepholeRemoveDeadCodeTest) TestNoRemoveThrow2() {
	c.FoldSame("function f(){throw a}")
}

// Java: PeepholeRemoveDeadCodeTest.testNoRemoveThrow3()
func (c *PeepholeRemoveDeadCodeTest) TestNoRemoveThrow3() {
	c.FoldSame("function f(){throw 10}")
}

// Java: PeepholeRemoveDeadCodeTest.testRedundantIfRemoved()
func (c *PeepholeRemoveDeadCodeTest) TestRedundantIfRemoved() {
	c.Fold("if(x()) 1", "x()")
}

// Java: PeepholeRemoveDeadCodeTest.testRemoveInControlStructure3()
func (c *PeepholeRemoveDeadCodeTest) TestRemoveInControlStructure3() {
	c.Fold("for(1;2;3) 4", "for(;;);")
}

// Java: PeepholeRemoveDeadCodeTest.testHook1()
func (c *PeepholeRemoveDeadCodeTest) TestHook1() {
	c.Fold("1 ? 2 : 3", "")
}

// Java: PeepholeRemoveDeadCodeTest.testHook2()
func (c *PeepholeRemoveDeadCodeTest) TestHook2() {
	c.Fold("x ? a() : 3", "x && a()")
}

// Java: PeepholeRemoveDeadCodeTest.testHook3()
func (c *PeepholeRemoveDeadCodeTest) TestHook3() {
	c.Fold("x ? 2 : a()", "x || a()")
}

// Java: PeepholeRemoveDeadCodeTest.testHook4()
func (c *PeepholeRemoveDeadCodeTest) TestHook4() {
	c.FoldSame("x ? a() : b()")
}

// Java: PeepholeRemoveDeadCodeTest.testHook5()
func (c *PeepholeRemoveDeadCodeTest) TestHook5() {
	c.Fold("a() ? 1 : 2", "a()")
}

// Java: PeepholeRemoveDeadCodeTest.testHook6()
func (c *PeepholeRemoveDeadCodeTest) TestHook6() {
	c.Fold("a() ? b() : 2", "a() && b()")
}

// TODO(johnlenz): Consider adding a post optimization pass to
// convert OR into HOOK to save parentheses when the operator
// precedents would require them.
// Java: PeepholeRemoveDeadCodeTest.testHook7()
func (c *PeepholeRemoveDeadCodeTest) TestHook7() {
	c.Fold("a() ? 1 : b()", "a() || b()")
}

// Java: PeepholeRemoveDeadCodeTest.testHook8()
func (c *PeepholeRemoveDeadCodeTest) TestHook8() {
	c.FoldSame("a() ? b() : c()")
}

// Java: PeepholeRemoveDeadCodeTest.testHook9()
func (c *PeepholeRemoveDeadCodeTest) TestHook9() {
	c.Fold("true ? a() : (function f() {})()", "a()")
	c.Fold("false ? a() : (function f() {alert(x)})()", "(function f() {alert(x)})()")
}

// Java: PeepholeRemoveDeadCodeTest.testHook10()
func (c *PeepholeRemoveDeadCodeTest) TestHook10() {
	c.Fold("((function () {}), true) ? a() : b()", "a()")
	c.Fold("((function () {alert(x)})(), true) ? a() : b()", "(function(){alert(x)})(),a()")
}

// Java: PeepholeRemoveDeadCodeTest.testShortCircuit1()
func (c *PeepholeRemoveDeadCodeTest) TestShortCircuit1() {
	c.FoldSame("1 && a()")
}

// Java: PeepholeRemoveDeadCodeTest.testShortCircuit2NullishCoalesce()
func (c *PeepholeRemoveDeadCodeTest) TestShortCircuit2NullishCoalesce() {
	c.Fold("1 ?? a() ?? 2", "1 ?? a()")
}

// Java: PeepholeRemoveDeadCodeTest.testShortCircuit3NullishCoalesce()
func (c *PeepholeRemoveDeadCodeTest) TestShortCircuit3NullishCoalesce() {
	c.Fold("a() ?? 1 ?? 2", "a()")
}

// Java: PeepholeRemoveDeadCodeTest.testShortCircuit4NullishCoalesce()
func (c *PeepholeRemoveDeadCodeTest) TestShortCircuit4NullishCoalesce() {
	c.FoldSame("a() ?? 1 ?? b()")
}

// Java: PeepholeRemoveDeadCodeTest.testShortCircuit2()
func (c *PeepholeRemoveDeadCodeTest) TestShortCircuit2() {
	c.Fold("1 && a() && 2", "1 && a()")
}

// Java: PeepholeRemoveDeadCodeTest.testShortCircuit3()
func (c *PeepholeRemoveDeadCodeTest) TestShortCircuit3() {
	c.Fold("a() && 1 && 2", "a()")
}

// Java: PeepholeRemoveDeadCodeTest.testShortCircuit4()
func (c *PeepholeRemoveDeadCodeTest) TestShortCircuit4() {
	c.FoldSame("a() && 1 && b()")
}

// Java: PeepholeRemoveDeadCodeTest.testComplex1()
func (c *PeepholeRemoveDeadCodeTest) TestComplex1() {
	c.Fold("1 && a() + b() + c()", "1 && (a(), b(), c())")
}

// Java: PeepholeRemoveDeadCodeTest.testComplex2()
func (c *PeepholeRemoveDeadCodeTest) TestComplex2() {
	c.Fold("1 && (a() ? b() : 1)", "1 && (a() && b())")
}

// Java: PeepholeRemoveDeadCodeTest.testComplex3()
func (c *PeepholeRemoveDeadCodeTest) TestComplex3() {
	c.Fold("1 && (a() ? b() : 1 + c())", "1 && (a() ? b() : c())")
}

// Java: PeepholeRemoveDeadCodeTest.testComplex4()
func (c *PeepholeRemoveDeadCodeTest) TestComplex4() {
	c.Fold("1 && (a() ? 1 : 1 + c())", "1 && (a() || c())")
}

// Java: PeepholeRemoveDeadCodeTest.testComplex5()
func (c *PeepholeRemoveDeadCodeTest) TestComplex5() {
	// can't simplify LHS of short circuit statements with side effects
	c.FoldSame("(a() ? 1 : 1 + c()) && foo()")
}

// Java: PeepholeRemoveDeadCodeTest.testNoRemoveFunctionDeclaration1()
func (c *PeepholeRemoveDeadCodeTest) TestNoRemoveFunctionDeclaration1() {
	c.FoldSame("function foo(){}")
}

// Java: PeepholeRemoveDeadCodeTest.testNoRemoveFunctionDeclaration2()
func (c *PeepholeRemoveDeadCodeTest) TestNoRemoveFunctionDeclaration2() {
	c.FoldSame("var foo = function (){}")
}

// Java: PeepholeRemoveDeadCodeTest.testNoSimplifyFunctionArgs1()
func (c *PeepholeRemoveDeadCodeTest) TestNoSimplifyFunctionArgs1() {
	c.FoldSame("f(1 + 2, 3 + g())")
}

// Java: PeepholeRemoveDeadCodeTest.testNoSimplifyFunctionArgs2()
func (c *PeepholeRemoveDeadCodeTest) TestNoSimplifyFunctionArgs2() {
	c.FoldSame("1 && f(1 + 2, 3 + g())")
}

// Java: PeepholeRemoveDeadCodeTest.testNoSimplifyFunctionArgs3()
func (c *PeepholeRemoveDeadCodeTest) TestNoSimplifyFunctionArgs3() {
	c.FoldSame("1 && foo(a() ? b() : 1 + c())")
}

// Java: PeepholeRemoveDeadCodeTest.testNoRemoveInherits1()
func (c *PeepholeRemoveDeadCodeTest) TestNoRemoveInherits1() {
	c.FoldSame("var a = {}; this.b = {}; var goog = {}; goog.inherits(b, a)")
}

// Java: PeepholeRemoveDeadCodeTest.testNoRemoveInherits2()
func (c *PeepholeRemoveDeadCodeTest) TestNoRemoveInherits2() {
	c.Fold(
		"var a = {}; this.b = {}; var goog = {}; goog.inherits(b, a) + 1",
		"var a = {}; this.b = {}; var goog = {}; goog.inherits(b, a)")
}

// Java: PeepholeRemoveDeadCodeTest.testNoRemoveInherits3()
func (c *PeepholeRemoveDeadCodeTest) TestNoRemoveInherits3() {
	c.FoldSame("this.a = {}; var b = {}; b.inherits(a);")
}

// Java: PeepholeRemoveDeadCodeTest.testNoRemoveInherits4()
func (c *PeepholeRemoveDeadCodeTest) TestNoRemoveInherits4() {
	c.Fold("this.a = {}; var b = {}; b.inherits(a) + 1;", "this.a = {}; var b = {}; b.inherits(a)")
}

// Java: PeepholeRemoveDeadCodeTest.testRemoveFromLabel1()
func (c *PeepholeRemoveDeadCodeTest) TestRemoveFromLabel1() {
	c.Fold("LBL: void 0", "")
}

// Java: PeepholeRemoveDeadCodeTest.testRemoveFromLabel2()
func (c *PeepholeRemoveDeadCodeTest) TestRemoveFromLabel2() {
	c.Fold("LBL: foo() + 1 + bar()", "LBL: foo(),bar()")
}

// Java: PeepholeRemoveDeadCodeTest.testCall()
func (c *PeepholeRemoveDeadCodeTest) TestCall() {
	c.FoldSame("foo(0)")
	// We use a function with no side-effects, otherwise the entire invocation would be preserved.
	c.Fold("Math.sin(0);", "")
	c.Fold("1 + Math.sin(0);", "")
}

// Java: PeepholeRemoveDeadCodeTest.testCall_toString()
func (c *PeepholeRemoveDeadCodeTest) TestCall_toString() {
	c.Fold("(10).toString();", "")
	c.Fold("(10).toString(16);", "")
	c.Fold("x.toString(2);", "")
	c.Fold("foo().toString(16);", "foo();")
	c.Fold("(10)?.toString(16);", "")
	c.FoldSame("(10).toString(16, 2);")
}

// Java: PeepholeRemoveDeadCodeTest.testCall_containingSpread()
func (c *PeepholeRemoveDeadCodeTest) TestCall_containingSpread() {
	// We use a function with no side-effects, otherwise the entire invocation would be preserved.
	c.Fold("Math.sin(...c)", "([...c])")
	c.Fold("Math.sin(4, ...c, a)", "([...c])")
	c.Fold("Math.sin(foo(), ...c, bar())", "(foo(), [...c], bar())")
	c.Fold("Math.sin(...a, b, ...c)", "([...a], [...c])")
	c.Fold("Math.sin(...b, ...c)", "([...b], [...c])")
}

// TODO(b/538125186): Do not remove optional chain side effects in PeepholeRemoveDeadCode
// Java: PeepholeRemoveDeadCodeTest.testOptChainCall_containingSpread()
func (c *PeepholeRemoveDeadCodeTest) TestOptChainCall_containingSpread() {
	// We use a function with no side-effects, otherwise the entire invocation would be preserved.
	c.FoldSame("Math?.sin(...c);")
	c.FoldSame("Math?.sin(4, ...c, a);")
	c.FoldSame("Math?.sin(foo(), ...c, bar());")
	c.FoldSame("Math?.sin(...a, b, ...c);")
	c.FoldSame("Math?.sin(...b, ...c);")
}

// Java: PeepholeRemoveDeadCodeTest.testNew()
func (c *PeepholeRemoveDeadCodeTest) TestNew() {
	c.FoldSame("new foo(0)")
	// We use a function with no side-effects, otherwise the entire invocation would be preserved.
	c.Fold("new Date;", "")
	c.Fold("1 + new Date;", "")
	c.FoldSame("new (class { x = foo(); })()")
	c.Fold("new (class { x = 1; })()", "")
}

// Java: PeepholeRemoveDeadCodeTest.testNew_containingSpread()
func (c *PeepholeRemoveDeadCodeTest) TestNew_containingSpread() {
	// We use a function with no side-effects, otherwise the entire invocation would be preserved.
	c.Fold("new Date(...c)", "([...c])")
	c.Fold("new Date(4, ...c, a)", "([...c])")
	c.Fold("new Date(foo(), ...c, bar())", "(foo(), [...c], bar())")
	c.Fold("new Date(...a, b, ...c)", "([...a], [...c])")
	c.Fold("new Date(...b, ...c)", "([...b], [...c])")
}

// Java: PeepholeRemoveDeadCodeTest.testTaggedTemplateLit_simpleTemplate()
func (c *PeepholeRemoveDeadCodeTest) TestTaggedTemplateLit_simpleTemplate() {
	c.FoldSame("foo`Simple`")
	// We use a function with no side-effects, otherwise the entire invocation would be preserved.
	c.Fold("Math.sin`Simple`", "")
	c.Fold("1 + Math.sin`Simple`", "")
}

// Java: PeepholeRemoveDeadCodeTest.testTaggedTemplateLit_substitutingTemplate()
func (c *PeepholeRemoveDeadCodeTest) TestTaggedTemplateLit_substitutingTemplate() {
	c.FoldSame("foo`Complex ${butSafe}`")
	// We use a function with no side-effects, otherwise the entire invocation would be preserved.
	c.Fold("Math.sin`Complex ${butSafe}`", "")
	c.Fold("Math.sin`Complex ${andDangerous()}`", "andDangerous()")
}

// Java: PeepholeRemoveDeadCodeTest.testFoldAssign()
func (c *PeepholeRemoveDeadCodeTest) TestFoldAssign() {
	c.Fold("x=x", "")
	c.FoldSame("x=xy")
	c.FoldSame("x=x + 1")
	c.FoldSame("x.a=x.a")
	c.Fold("var y=(x=x)", "var y=x")
	c.Fold("y=1 + (x=x)", "y=1 + x")
}

// Java: PeepholeRemoveDeadCodeTest.testTryCatchFinally()
func (c *PeepholeRemoveDeadCodeTest) TestTryCatchFinally() {
	c.FoldSame("try {foo()} catch (e) {bar()}")
	c.FoldSame("try { try {foo()} catch (e) {bar()}} catch (x) {bar()}")
	c.Fold("try {var x = 1} finally {}", "var x = 1;")
	c.FoldSame("try {var x = 1} finally {x()}")
	c.Fold( //
		"function f() { return; try{ var x = 1; }finally{} }", "function f() { var x; return; }")
	c.Fold("try {} finally {x()}", "x()")
	c.Fold("try {} catch (e) { bar()} finally {x()}", "x()")
	c.Fold("try {} catch (e) { bar()}", "")
	c.Fold("try {} catch (e) { var a = 0; } finally {x()}", "var a; x()")
	c.Fold("try {} catch (e) {}", "")
	c.Fold("try {} finally {}", "")
	c.Fold("try {} catch (e) {} finally {}", "")
	c.Fold("L1:try {} catch (e) {} finally {}", "")
	c.Fold("L2:L1:try {} catch (e) {} finally {}", "")
}

// Java: PeepholeRemoveDeadCodeTest.testObjectLiteral()
func (c *PeepholeRemoveDeadCodeTest) TestObjectLiteral() {
	c.Fold("({})", "")
	c.Fold("({a:1})", "")
	c.Fold("({a:foo()})", "foo()")
	c.Fold("({'a':foo()})", "foo()")
	// Object-spread may tigger getters.
	c.FoldSame("({...a})")
	c.FoldSame("({...foo()})")
}

// Java: PeepholeRemoveDeadCodeTest.testArrayLiteral()
func (c *PeepholeRemoveDeadCodeTest) TestArrayLiteral() {
	c.Fold("([])", "")
	c.Fold("([1])", "")
	c.Fold("([a])", "")
	c.Fold("([foo()])", "foo()")
}

// Java: PeepholeRemoveDeadCodeTest.testArrayLiteral_containingSpread()
func (c *PeepholeRemoveDeadCodeTest) TestArrayLiteral_containingSpread() {
	c.FoldSame("([...c])")
	c.Fold("([4, ...c, a])", "([...c])")
	c.Fold("([foo(), ...c, bar()])", "(foo(), [...c], bar())")
	c.Fold("([...a, b, ...c])", "([...a], [...c])")
	c.FoldSame("([...b, ...c])")
	// It would also be fine if the spreads were split apart.
}

// Java: PeepholeRemoveDeadCodeTest.testAwait()
func (c *PeepholeRemoveDeadCodeTest) TestAwait() {
	c.FoldSame("async function f() { await something(); }")
	c.FoldSame("async function f() { await some.thing(); }")
}

// Java: PeepholeRemoveDeadCodeTest.testEmptyPatternInDeclarationRemoved()
func (c *PeepholeRemoveDeadCodeTest) TestEmptyPatternInDeclarationRemoved() {
	c.Fold("var [] = [];", "")
	c.Fold("let [] = [];", "")
	c.Fold("const [] = [];", "")
	c.Fold("var {} = [];", "")
	c.Fold("var [] = foo();", "foo()")
}

// Java: PeepholeRemoveDeadCodeTest.testEmptyArrayPatternInAssignRemoved()
func (c *PeepholeRemoveDeadCodeTest) TestEmptyArrayPatternInAssignRemoved() {
	c.Fold("({} = {});", "")
	c.Fold("({} = foo());", "foo()")
	c.Fold("[] = [];", "")
	c.Fold("[] = foo();", "foo()")
}

// Java: PeepholeRemoveDeadCodeTest.testEmptyPatternInParamsNotRemoved()
func (c *PeepholeRemoveDeadCodeTest) TestEmptyPatternInParamsNotRemoved() {
	c.FoldSame("function f([], a) {}")
	c.FoldSame("function f({}, a) {}")
}

// Java: PeepholeRemoveDeadCodeTest.testEmptyPatternInForOfLoopNotRemoved()
func (c *PeepholeRemoveDeadCodeTest) TestEmptyPatternInForOfLoopNotRemoved() {
	c.FoldSame("for (let [] of foo()) {}")
	c.FoldSame("for (const [] of foo()) {}")
	c.FoldSame("for ([] of foo()) {}")
	c.FoldSame("for ({} of foo()) {}")
}

// Java: PeepholeRemoveDeadCodeTest.testEmptySlotInArrayPatternRemoved()
func (c *PeepholeRemoveDeadCodeTest) TestEmptySlotInArrayPatternRemoved() {
	c.Fold("[,,] = foo();", "foo()")
	c.Fold("[a,b,,] = foo();", "[a,b] = foo();")
	c.Fold("[a,[],b,[],[]] = foo();", "[a,[],b] = foo();")
	c.Fold("[a,{},b,{},{}] = foo();", "[a,{},b] = foo();")
	c.Fold("function f([,,,]) {}", "function f([]) {}")
	c.FoldSame("[[], [], [], ...rest] = foo()")
}

// Java: PeepholeRemoveDeadCodeTest.testEmptySlotInArrayPatternWithDefaultValueMaybeRemoved()
func (c *PeepholeRemoveDeadCodeTest) TestEmptySlotInArrayPatternWithDefaultValueMaybeRemoved() {
	c.Fold("[a,[] = 0] = [];", "[a] = [];")
	c.FoldSame("[a,[] = foo()] = [];")
}

// Java: PeepholeRemoveDeadCodeTest.testEmptyKeyInObjectPatternRemoved()
func (c *PeepholeRemoveDeadCodeTest) TestEmptyKeyInObjectPatternRemoved() {
	c.Fold("const {f: {}} = {};", "")
	c.Fold("const {f: []} = {};", "")
	c.Fold("const {f: {}, g} = {};", "const {g} = {};")
	c.Fold("const {f: [], g} = {};", "const {g} = {};")
	c.FoldSame("const {[foo()]: {}} = {};")
}

// Java: PeepholeRemoveDeadCodeTest.testEmptyKeyInObjectPatternWithDefaultValueMaybeRemoved()
func (c *PeepholeRemoveDeadCodeTest) TestEmptyKeyInObjectPatternWithDefaultValueMaybeRemoved() {
	c.Fold("const {f: {} = 0} = {};", "")
	// In theory the following case could be reduced to `foo()`, but that gets more complicated to
	// implement for object patterns with multiple keys with side effects.
	// Instead the pass backs off for any default with a possible side effect
	c.FoldSame("const {f: {} = foo()} = {};")
}

// Java: PeepholeRemoveDeadCodeTest.testEmptyKeyInObjectPatternNotRemovedWithObjectRest()
func (c *PeepholeRemoveDeadCodeTest) TestEmptyKeyInObjectPatternNotRemovedWithObjectRest() {
	c.FoldSame("const {f: {}, ...g} = foo()")
	c.FoldSame("const {f: [], ...g} = foo()")
}

// Java: PeepholeRemoveDeadCodeTest.testUndefinedDefaultParameterRemoved()
func (c *PeepholeRemoveDeadCodeTest) TestUndefinedDefaultParameterRemoved() {
	c.Fold(
		"function f(x=undefined,y) {  }", //
		"function f(x,y)             {  }")
	c.Fold(
		"function f(x,y=undefined,z) {  }", //
		"function f(x,y          ,z) {  }")
	c.Fold(
		"function f(x=undefined,y=undefined,z=undefined) {  }", //
		"function f(x,          y,          z)           {  }")
}

// Java: PeepholeRemoveDeadCodeTest.testPureVoidDefaultParameterRemoved()
func (c *PeepholeRemoveDeadCodeTest) TestPureVoidDefaultParameterRemoved() {
	c.Fold(
		"function f(x = void 0) {  }", //
		"function f(x         ) {  }")
	c.Fold(
		"function f(x = void \"XD\") {  }", //
		"function f(x              ) {  }")
	c.Fold(
		"function f(x = void f()) {  }", //
		"function f(x)            {  }")
}

// Java: PeepholeRemoveDeadCodeTest.testNoDefaultParameterNotRemoved()
func (c *PeepholeRemoveDeadCodeTest) TestNoDefaultParameterNotRemoved() {
	c.FoldSame("function f(x,y) {  }")
	c.FoldSame("function f(x) {  }")
	c.FoldSame("function f() {  }")
}

// Java: PeepholeRemoveDeadCodeTest.testEffectfulDefaultParameterNotRemoved()
func (c *PeepholeRemoveDeadCodeTest) TestEffectfulDefaultParameterNotRemoved() {
	c.FoldSame("function f(x = void console.log(1)) {  }")
	c.FoldSame("function f(x = void f()) { alert(x); }")
}

// Java: PeepholeRemoveDeadCodeTest.testDestructuringUndefinedDefaultParameter()
func (c *PeepholeRemoveDeadCodeTest) TestDestructuringUndefinedDefaultParameter() {
	c.Fold(
		"function f({a=undefined,b=1,c}) {  }", //
		"function f({a          ,b=1,c}) {  }")
	c.Fold(
		"function f({a={},b=0}=undefined) {  }", //
		"function f({a={},b=0}) {  }")
	c.Fold(
		"function f({a=undefined,b=0}) {  }", //
		"function f({a,b=0}) {  }")
	c.Fold(
		" function f({a: {b = undefined}}) {  }", //
		" function f({a: {b}}) {  }")
	c.FoldSame("function f({a,b}) {  }")
	c.FoldSame("function f({a=0, b=1}) {  }")
	c.FoldSame("function f({a=0,b=0}={}) {  }")
	c.FoldSame("function f({a={},b=0}={}) {  }")
}

// Java: PeepholeRemoveDeadCodeTest.testUndefinedDefaultObjectPatterns()
func (c *PeepholeRemoveDeadCodeTest) TestUndefinedDefaultObjectPatterns() {
	c.Fold(
		"const {a = undefined} = obj;", //
		"const {a} = obj;")
	c.Fold(
		"const {a = void 0} = obj;", //
		"const {a} = obj;")
}

// Java: PeepholeRemoveDeadCodeTest.testBatchC_destructuringOptimization()
func (c *PeepholeRemoveDeadCodeTest) TestBatchC_destructuringOptimization() {
	// OPP-015: Pure undefined / void default values
	c.Fold("const {a = void 0, b = undefined} = obj;", "const {a, b} = obj;")
	c.Fold("const [a = void 0, b = undefined] = arr;", "const [a, b] = arr;")

	// Empty nested patterns in object patterns
	c.Fold("const {a: {}, b: {}} = obj;", "")
	c.Fold("const {a: {}, b} = obj;", "const {b} = obj;")
	c.Fold("const {a: [], b} = obj;", "const {b} = obj;")

	// Trailing empty slots and empty nested patterns in array patterns
	c.Fold("const [a, , ,] = arr;", "const [a] = arr;")
	c.Fold("const [a, [], {}] = arr;", "const [a] = arr;")
	c.Fold("let [a, [] = 0] = arr;", "let [a] = arr;")

	// Guard cases: Object rest property exclusion, effectful defaults, computed keys
	c.FoldSame("const {a: {}, ...rest} = obj;")
	c.FoldSame("const {a: [] = foo()} = obj;")
	c.FoldSame("const {[foo()]: {}} = obj;")
	c.FoldSame("function f({}, x) {}")
	c.FoldSame("for (const {} of iter) {}")
	c.FoldSame("const [, a] = arr;")
}

// Java: PeepholeRemoveDeadCodeTest.testDoNotRemoveGetterOnlyAccess()
func (c *PeepholeRemoveDeadCodeTest) TestDoNotRemoveGetterOnlyAccess() {
	c.FoldSame(
		"var a = {\n" +
			"  get property() {}\n" +
			"};\n" +
			"a.property;\n")
	c.FoldSame(
		"var a = {\n" +
			"  get property() {}\n" +
			"};\n" +
			"a?.property;\n")
	c.FoldSame(
		"var a = {};\n" +
			"Object.defineProperty(a, 'property', {\n" +
			"  get() {}\n" +
			"});\n" +
			"a.property;\n")
	c.FoldSame(
		"var a = {};\n" +
			"Object.defineProperty(a, 'property', {\n" +
			"  get() {}\n" +
			"});\n" +
			"a?.property;\n")
}

// Java: PeepholeRemoveDeadCodeTest.testDoNotRemoveNestedGetterOnlyAccess()
func (c *PeepholeRemoveDeadCodeTest) TestDoNotRemoveNestedGetterOnlyAccess() {
	c.FoldSame(
		"var a = {\n" +
			"  b: { get property() {} }\n" +
			"};\n" +
			"a.b.property;\n")
}

// Java: PeepholeRemoveDeadCodeTest.testRemoveAfterNestedGetterOnlyAccess()
func (c *PeepholeRemoveDeadCodeTest) TestRemoveAfterNestedGetterOnlyAccess() {
	c.Fold(
		"var a = {\n"+
			"  b: { get property() {} }\n"+
			"};\n"+
			"a.b.property.d.e;\n",
		"var a = {\n"+
			"  b: { get property() {} }\n"+
			"};\n"+
			"a.b.property;\n")
}

// Java: PeepholeRemoveDeadCodeTest.testFoldLabelledEmptyBlock()
func (c *PeepholeRemoveDeadCodeTest) TestFoldLabelledEmptyBlock() {
	c.Fold("a:{}", "")
	c.Fold("a:b:{}", "")
	c.Fold("a:b:c:{}", "")
}

// Java: PeepholeRemoveDeadCodeTest.testRetainSetterOnlyAccess()
func (c *PeepholeRemoveDeadCodeTest) TestRetainSetterOnlyAccess() {
	c.FoldSame(
		"var a = {\n" +
			"  set property(v) {}\n" +
			"};\n" +
			"a.property;\n")
	c.FoldSame(
		"var a = {\n" +
			"  set property(v) {}\n" +
			"};\n" +
			"a?.property;\n")
}

// Java: PeepholeRemoveDeadCodeTest.testDoNotRemoveGetterSetterAccess()
func (c *PeepholeRemoveDeadCodeTest) TestDoNotRemoveGetterSetterAccess() {
	c.FoldSame(
		"var a = {\n" +
			"  get property() {},\n" +
			"  set property(x) {}\n" +
			"};\n" +
			"a.property;\n")
}

// Java: PeepholeRemoveDeadCodeTest.testDoNotRemoveSetSetterToGetter()
func (c *PeepholeRemoveDeadCodeTest) TestDoNotRemoveSetSetterToGetter() {
	c.FoldSame(
		"var a = {\n" +
			"  get property() {},\n" +
			"  set property(x) {}\n" +
			"};\n" +
			"a.property = a.property;\n")
}

// Java: PeepholeRemoveDeadCodeTest.testDoNotRemoveAccessIfOtherPropertyIsGetter()
func (c *PeepholeRemoveDeadCodeTest) TestDoNotRemoveAccessIfOtherPropertyIsGetter() {
	c.FoldSame(
		"var a = {\n" +
			"  get property() {}\n" +
			"};\n" +
			"var b = {\n" +
			"  property: 0,\n" +
			"};\n" +
			"// This pass should be conservative and not remove this since it sees a getter for\n" +
			"// \"property\"\n" +
			"b.property;\n")
	c.FoldSame(
		"var a = {};\n" +
			"Object.defineProperty(a, 'property', {\n" +
			"  get() {}\n" +
			"});\n" +
			"var b = {\n" +
			"  property: 0,\n" +
			"};\n" +
			"b.property;\n")
}

// Java: PeepholeRemoveDeadCodeTest.testFunctionCallReferencesGetterIsNotRemoved()
func (c *PeepholeRemoveDeadCodeTest) TestFunctionCallReferencesGetterIsNotRemoved() {
	c.FoldSame(
		"var a = {\n" +
			"  get property() {}\n" +
			"};\n" +
			"function foo() { a.property; }\n" +
			"foo();\n")
}

// Java: PeepholeRemoveDeadCodeTest.testFunctionCallReferencesSetterIsNotRemoved()
func (c *PeepholeRemoveDeadCodeTest) TestFunctionCallReferencesSetterIsNotRemoved() {
	c.FoldSame(
		"var a = {\n" +
			"  set property(v) {}\n" +
			"};\n" +
			"function foo() { a.property = 0; }\n" +
			"foo();\n")
}

// Java: PeepholeRemoveDeadCodeTest.testClassField()
func (c *PeepholeRemoveDeadCodeTest) TestClassField() {
	c.Fold(
		"class C {\n"+
			"  f1 = (5,2);\n"+
			"}\n",
		"class C {\n"+
			"  f1 = 2;\n"+
			"}\n")
}

// Java: PeepholeRemoveDeadCodeTest.testThis()
func (c *PeepholeRemoveDeadCodeTest) TestThis() {
	c.Fold(
		"class C {\n"+
			"  constructor() {\n"+
			"    this.f1 = (5,2);\n"+
			"  }\n"+
			"}\n",
		"class C {\n"+
			"  constructor() {\n"+
			"    this.f1 = 2;\n"+
			"  }\n"+
			"}\n")
}

// Java: PeepholeRemoveDeadCodeTest.testClassStaticBlock()
func (c *PeepholeRemoveDeadCodeTest) TestClassStaticBlock() {
	c.Fold(
		"class C {\n"+
			"  static {\n"+
			"  }\n"+
			"}\n",
		"class C {\n"+
			"}\n")
	c.FoldSame(
		"class C {\n" +
			"  static {\n" +
			"    this.x = 0;\n" +
			"  }\n" +
			"}\n")
}

// Java: PeepholeRemoveDeadCodeTest.testRemoveUnreachableOptionalChainingCall()
func (c *PeepholeRemoveDeadCodeTest) TestRemoveUnreachableOptionalChainingCall() {
	c.Fold("(null)?.();", "")
	c.Fold("(void 0)?.();", "")
	c.Fold("(undefined)?.();", "")
	c.Fold("(void 0)?.(0)", "")
	c.Fold("(void 0)?.(function f() {})", "")
	c.Fold("(null)?.x;", "")
	c.Fold("(void 0)?.x;", "")
	c.Fold("(null)?.['x'];", "")
	c.Fold("(void 0)?.['x'];", "")
	c.Fold("(null)?.[x];", "")
	c.Fold("(void 0)?.[x];", "")
	// arguments with unknown side effects are also removed
	c.Fold("(void 0)?.(f(), g())", "")
	// void arguments with unknown side effects are preserved
	c.Fold("(void f())?.();", "f();")
	c.Fold("g((void f())?.());", "g(void f());")
	c.FoldSame("(f(), null)?.()")
	c.FoldSame("f?.()")
	c.Fold("a?.x;", "")
	c.Fold("a?.['x'];", "")
}

// Java: PeepholeRemoveDeadCodeTest.testRemoveUnusedVoid()
func (c *PeepholeRemoveDeadCodeTest) TestRemoveUnusedVoid() {
	// remove void at statement level
	c.Fold("void 0;", "")
	c.Fold("void foo();", "foo();")
	// preserve void when passed somewhere else
	c.FoldSame("use(void 0);")
	c.FoldSame("use(void foo());")
	c.FoldSame("use(() => void foo());")
	c.Fold("void use(() => void foo());", "use(() => void foo());")
}

// Java: PeepholeRemoveDeadCodeTest.testInFn(String, String)
func (c *PeepholeRemoveDeadCodeTest) TestInFn(js string, expected string) {
	pre := "function f() {"
	post := "}"
	c.Test(pre+js+post, pre+expected+post)
}

// Java: PeepholeRemoveDeadCodeTest.testInLoop(String, String)
func (c *PeepholeRemoveDeadCodeTest) TestInLoop(js string, expected string) {
	c.TestInLoopWithExpectedBeforeLoop(js, "", expected)
}

// Java: PeepholeRemoveDeadCodeTest.testInLoop(String, String, String)
func (c *PeepholeRemoveDeadCodeTest) TestInLoopWithExpectedBeforeLoop(js string, expectedBeforeLoop string, expected string) {
	pre := "for (;;) {"
	post := "}"
	c.Test(pre+js+post, expectedBeforeLoop+pre+expected+post)
}

// Java: PeepholeRemoveDeadCodeTest.testRemoveDeadStatements1()
func (c *PeepholeRemoveDeadCodeTest) TestRemoveDeadStatements1() {
	c.Test("throw 1; x;", "throw 1;")
	c.Test("throw 1; alert(1)", "throw 1;")
	c.Test("throw 1; var x = 1", "var x; throw 1;")
}

// Java: PeepholeRemoveDeadCodeTest.testRemoveDeadStatements2()
func (c *PeepholeRemoveDeadCodeTest) TestRemoveDeadStatements2() {
	c.TestInFn("return; x;", "return;")
	c.TestInFn("return; alert(1)", "return;")
	c.TestInFn("return; var x = 1", "var x; return;")
}

// Java: PeepholeRemoveDeadCodeTest.testRemoveDeadStatements3()
func (c *PeepholeRemoveDeadCodeTest) TestRemoveDeadStatements3() {
	c.TestInLoop("break; x;", "break;")
	c.TestInLoop("break; alert(1)", "break;")
	c.TestInLoopWithExpectedBeforeLoop("break; var x = 1", "var x;", "break;")
}

// Java: PeepholeRemoveDeadCodeTest.testRemoveDeadStatements4()
func (c *PeepholeRemoveDeadCodeTest) TestRemoveDeadStatements4() {
	c.TestInLoop("continue; x;", "continue;")
	c.TestInLoop("continue; alert(1)", "continue;")
	c.TestInLoopWithExpectedBeforeLoop("continue; var x = 1", "var x;", "continue;")
}

// Java: PeepholeRemoveDeadCodeTest.testRemovalRequiresRedeclaration()
func (c *PeepholeRemoveDeadCodeTest) TestRemovalRequiresRedeclaration() {
	c.Test( //
		"while(1) { break; var x = 1}", //
		"var x; for(;;) { break }")
	c.Test( //
		"while(1) { break; var x=1; var y=1 }", //
		"var y; var x; for(;;) { break }")
	c.Test( //
		"while(1) { break; var [x, [[[y]]]] = [];}", //
		"var y; var x; for(;;) { break }")
}

// Java: PeepholeRemoveDeadCodeTest.testRemovalRequiresRedeclaration_normalizeDisabled()
func (c *PeepholeRemoveDeadCodeTest) TestRemovalRequiresRedeclaration_normalizeDisabled() {
	c.DisableNormalize()
	c.DisableComputeSideEffects()
	c.Test( //
		"while(1) { break; var x = 1}", //
		"var x; while (1) { break; }")
	c.Test( //
		"while(1) { break; var x=1; var y=1 }", //
		"var y; var x; while (1) { break; }")
	c.Test( //
		"while(1) { break; var [x, [[[y]]]] = [];}", //
		"var y; var x; while (1) { break; }")
}

// Java: PeepholeRemoveDeadCodeTest.testRemoveDo()
func (c *PeepholeRemoveDeadCodeTest) TestRemoveDo() {
	c.Test("do { print(1); break } while(1)", "do { print(1); break } while(1)")
	c.Test("while(1) { break; do { print(1); break } while(1) }", "for (;;) { break; }")
}

// Java: PeepholeRemoveDeadCodeTest.testSwitchCase()
func (c *PeepholeRemoveDeadCodeTest) TestSwitchCase() {
	c.Test(
		"function f() { switch(x) { case 1: break; default: return 5; foo()}}",
		"function f() { switch(x) { case 1: break; default: return 5;}}")
	c.Test(
		"function f() { switch(x) { default: return; case 1: foo(); bar()}}",
		"function f() { switch(x) { default: return; case 1: foo(); bar()}}")
	c.Test(
		"function f() { switch(x) { default: return; case 1: return 5;bar()}}",
		"function f() { switch(x) { default: return; case 1: return 5;}}")
}

// Java: PeepholeRemoveDeadCodeTest.testBatchE_switchConstantDiscriminantFolding()
func (c *PeepholeRemoveDeadCodeTest) TestBatchE_switchConstantDiscriminantFolding() {
	// OPP-023: Switch Constant Discriminant Folding & Dead Case Pruning
	c.Fold(
		"switch ('foo') {\n"+
			"  case 'foo':\n"+
			"    foo();\n"+
			"    break;\n"+
			"  case 'bar':\n"+
			"    bar();\n"+
			"    break;\n"+
			"}\n",
		"foo();")

	c.Fold(
		"switch ('noMatch') {\n"+
			"  case 'foo':\n"+
			"    foo();\n"+
			"    break;\n"+
			"  case 'bar':\n"+
			"    bar();\n"+
			"    break;\n"+
			"}\n",
		"")

	c.Fold("switch(1){case 2: var x=0;}", "var x;")

	// Guard cases
	c.FoldSame("switch (x) { case 1: foo(); break; }")
	c.FoldSame("function f() { switch(a) { case 1: foo(); } }")
	c.FoldSame("function f() { switch(a) { default: return; case 1: break; } }")
}

// Java: PeepholeRemoveDeadCodeTest.testBatchE_pureExpressionStatementAndBlockElimination()
func (c *PeepholeRemoveDeadCodeTest) TestBatchE_pureExpressionStatementAndBlockElimination() {
	// OPP-024: Pure Expression Statement & Block Elimination
	c.Fold("{ 'hi' }", "")
	c.Fold("{ x == 3 }", "")
	c.Fold("{ `hello ${foo}` }", "")
	c.Fold("{ (function(){ x++; }) }", "")
	c.Fold("{{ foo() }}", "foo()")
	c.Fold("{ foo(); {} }", "foo()")
	c.Fold("{{ foo() } { bar() }}", "foo(); bar()")

	// Guard cases: side-effectful statements and declarations
	c.FoldSame("foo();")
	c.FoldSame("function f() { return 3; }")
	c.FoldSame("for (let x = 1; x < 10; x++) {}")
}

// Java: PeepholeRemoveDeadCodeTest.testUnusedClassWithUninitializedField()
func (c *PeepholeRemoveDeadCodeTest) TestUnusedClassWithUninitializedField() {
	c.Test("(class { x; });", "")
}

// Java: PeepholeRemoveDeadCodeTest.testUnusedClassWithInstanceFieldSideEffects()
func (c *PeepholeRemoveDeadCodeTest) TestUnusedClassWithInstanceFieldSideEffects() {
	c.Test("(class { x = alert(1); });", "")
	c.Test("(class { x = alert(1); y = alert(2); });", "")
}

// Java: PeepholeRemoveDeadCodeTest.testUnusedClassWithInstanceFieldSideEffects_thisAndSuper()
func (c *PeepholeRemoveDeadCodeTest) TestUnusedClassWithInstanceFieldSideEffects_thisAndSuper() {
	c.Test("(class { x = this.y(); });", "")
	c.Test("(class extends Object { x = super.toString(); });", "")
}

// Java: PeepholeRemoveDeadCodeTest.testUnusedClassWithComputedInstanceFieldSideEffects()
func (c *PeepholeRemoveDeadCodeTest) TestUnusedClassWithComputedInstanceFieldSideEffects() {
	c.Test("(class { [alert(1)] = alert(2); });", "alert(1);")
	c.Test("(class { ['x'] = alert(2); });", "")
}

// Java: PeepholeRemoveDeadCodeTest.testUnusedClassWithComputedInstanceFieldSideEffects_uninitialized()
func (c *PeepholeRemoveDeadCodeTest) TestUnusedClassWithComputedInstanceFieldSideEffects_uninitialized() {
	c.Test("(class { [alert(1)]; });", "alert(1);")
}

// Java: PeepholeRemoveDeadCodeTest.testUnusedClassWithStaticField_uninitialized()
func (c *PeepholeRemoveDeadCodeTest) TestUnusedClassWithStaticField_uninitialized() {
	c.Test("(class { static x; });", "")
	c.Test("(class { static [alert(1)]; });", "alert(1);")
}

// Java: PeepholeRemoveDeadCodeTest.testUnusedClassWithStaticFieldSideEffects()
func (c *PeepholeRemoveDeadCodeTest) TestUnusedClassWithStaticFieldSideEffects() {
	c.Test("(class { static x = alert(1); });", "alert(1);")
	c.Test("(class { static x = alert(1); static y = alert(2); });", "alert(1), alert(2);")
	c.Test("(class { static [alert(1)] = alert(2); });", "alert(1), alert(2);")
}

// Java: PeepholeRemoveDeadCodeTest.testUnusedClassWithStaticFieldSideEffects_thisAndSuper()
func (c *PeepholeRemoveDeadCodeTest) TestUnusedClassWithStaticFieldSideEffects_thisAndSuper() {
	// TODO(b/574118308): fix the `super.y()` and `this.y()` inlining.
	c.Test("(class extends Parent { static x = this.y(); } );", "this.y();")
	ex :=
		AssertThrows[jutil.Throwable](
			func() { c.TestSame("(class extends Parent { static x = super.y(); });") })
	AssertThat(ex.GetMessage()).
		Contains(
			"super property references are only allowed in methods, class static blocks and class" +
				" fields")
	// This is correct - we know `super.toString()` has no side-effects so is safe to remove.
	c.Test("(class extends Object { static x = super.toString(); });", "")
}
