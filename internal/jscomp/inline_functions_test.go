// Port of com/google/javascript/jscomp/InlineFunctionsTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/InlineFunctionsTest.java)

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
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// InlineFunctionsTest: inline function tests.
type InlineFunctionsTest struct {
	CompilerTestCaseBase

	// Go-only: language-adapter — the receiver of maybeEnableInferConsts(),
	// which one test's anonymous subclass overrides.
	self interface{ MaybeEnableInferConsts() }

	inliningReach        CompilerOptions_Reach
	assumeStrictThis     bool
	assumeMinimumCapture bool
	maxSizeAfterInlining int

	uniqueIdSupplier *InlineFunctionsTest_UniqueIdSupplier
}

// Java: InlineFunctionsTest.EXTERNS
const InlineFunctionsTest_EXTERNS = "/** @nosideeffects */ function nochg(){}\nfunction chg(){}\n"

// Java: InlineFunctionsTest()
func NewInlineFunctionsTest() *InlineFunctionsTest {
	c := &InlineFunctionsTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, InlineFunctionsTest_EXTERNS)
	c.self = c
	c.uniqueIdSupplier = NewInlineFunctionsTest_UniqueIdSupplier()
	return c
}

func TestInlineFunctionsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewInlineFunctionsTest() })
}

// MaybeEnableInferConsts is overridden by one test method that needs to
// disable this.
//
// Java: InlineFunctionsTest.maybeEnableInferConsts()
func (c *InlineFunctionsTest) MaybeEnableInferConsts() {
	c.EnableInferConsts()
}

// Java: InlineFunctionsTest.setUp()
func (c *InlineFunctionsTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableDebugLogging(true)
	c.self.MaybeEnableInferConsts()
	c.EnableNormalize()
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.EnableComputeSideEffects()
	c.inliningReach = CompilerOptions_Reach_ALL
	c.assumeStrictThis = false
	c.assumeMinimumCapture = false
	c.maxSizeAfterInlining = CompilerOptions_UNLIMITED_FUN_SIZE_AFTER_INLINING
	c.uniqueIdSupplier = NewInlineFunctionsTest_UniqueIdSupplier()
}

// Java: InlineFunctionsTest.getProcessor(Compiler)
func (c *InlineFunctionsTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewInlineFunctions(
		compiler,
		c.uniqueIdSupplier.Get,
		c.inliningReach,
		c.assumeStrictThis,
		c.assumeMinimumCapture,
		c.maxSizeAfterInlining)
}

// InlineFunctionsTest_UniqueIdSupplier is the Supplier<String> of unique
// ids handed to InlineFunctions.
type InlineFunctionsTest_UniqueIdSupplier struct {
	nextId int
}

// Java: InlineFunctionsTest.UniqueIdSupplier()
func NewInlineFunctionsTest_UniqueIdSupplier() *InlineFunctionsTest_UniqueIdSupplier {
	return &InlineFunctionsTest_UniqueIdSupplier{}
}

// Java: InlineFunctionsTest.UniqueIdSupplier.get()
func (u *InlineFunctionsTest_UniqueIdSupplier) Get() string {
	id := u.nextId
	u.nextId++
	return "v" + strconv.Itoa(id) // prefix with "v" to distinguish from Compiler's unique id supplier.
}

// Returns the number of times the pass should be run before results are verified.
// Java: InlineFunctionsTest.getNumRepetitions()
func (c *InlineFunctionsTest) GetNumRepetitions() int {
	// Some inlining can only be done in multiple passes.
	return 3
}

// Java: InlineFunctionsTest.testNoInline()
func (c *InlineFunctionsTest) TestNoInline() {
	c.TestSame("/** @noinline */ function foo(){} foo();foo();foo();")
	c.TestSame("/** @noinline */ var foo = function(){}; foo();foo();foo();")
	c.TestSame("/** @noinline */ var foo = function(){}; foo?.();foo?.();foo?.();")
}

// Java: InlineFunctionsTest.testRequireInlining()
func (c *InlineFunctionsTest) TestRequireInlining() {
	c.Test(
		`/** @requireInlining */ function foo() {console.log('hi');console.log('bye');}
foo();foo();foo();
`,
		`{console.log('hi');console.log('bye');}
{console.log('hi');console.log('bye');}
{console.log('hi');console.log('bye');}
`)
}

// Java: InlineFunctionsTest.testEncourageInlining()
func (c *InlineFunctionsTest) TestEncourageInlining() {
	c.Test(
		`/** @encourageInlining */ function foo() {console.log('hi');console.log('bye');}
foo();foo();foo();
`,
		`{console.log('hi');console.log('bye');}
{console.log('hi');console.log('bye');}
{console.log('hi');console.log('bye');}
`)
}

// Java: InlineFunctionsTest.testEncourageInliningAliased()
func (c *InlineFunctionsTest) TestEncourageInliningAliased() {
	// Aliasing @requireInlining functions is allowed when they are
	// const. The idea is that that declaration will be inlined and then
	// a later InlineFunctions will inline as required.
	c.TestSame(
		`/** @encourageInlining */ function foo() {console.log('hi');console.log('bye');}
const bar = foo;
bar();bar();bar();
`)
}

// Java: InlineFunctionsTest.testEncourageInliningNested()
func (c *InlineFunctionsTest) TestEncourageInliningNested() {
	c.Test(
		`/** @encourageInlining */ function foo() {console.log('hi');console.log('bye');}
/** @encourageInlining */ function bar() {foo();}
bar();bar();bar();
`,
		`{{console.log('hi');console.log('bye');}}
{{console.log('hi');console.log('bye');}}
{{console.log('hi');console.log('bye');}}
`)
}

// Java: InlineFunctionsTest.testEncourageInliningDoesNotWarnWhenNotPossible()
func (c *InlineFunctionsTest) TestEncourageInliningDoesNotWarnWhenNotPossible() {
	c.TestSame(
		`/** @encourageInlining */ function foo(q) {console.log(q);console.log('bye');}
console.log(/** @type {?} */ (foo).name);
`)
}

// Java: InlineFunctionsTest.testRequireInliningAliased()
func (c *InlineFunctionsTest) TestRequireInliningAliased() {
	// Aliasing @requireInlining functions is allowed when they are
	// const. The idea is that that declaration will be inlined and then
	// a later InlineFunctions will inline as required.
	c.TestSame(
		`/** @requireInlining */ function foo() {console.log('hi');console.log('bye');}
const bar = foo;
bar();bar();bar();
`)
}

// Java: InlineFunctionsTest.testRequireInliningNested()
func (c *InlineFunctionsTest) TestRequireInliningNested() {
	c.Test(
		`/** @requireInlining */ function foo() {console.log('hi');console.log('bye');}
/** @requireInlining */ function bar() {foo();}
bar();bar();bar();
`,
		`{{console.log('hi');console.log('bye');}}
{{console.log('hi');console.log('bye');}}
{{console.log('hi');console.log('bye');}}
`)
}

// Java: InlineFunctionsTest.testRequireInliningCascades()
func (c *InlineFunctionsTest) TestRequireInliningCascades() {
	c.DisableCompareAsTree()
	c.Test(
		"const a = class {constructor() {} setA(value) {this.a = value;} setB(value) {this.b = value;}}; /** @encourageInlining */ function fromRecord(rec) { return new a().setA(rec.a).setB(rec.b); }; console.log(fromRecord({a: 1, b: fromRecord({a: 2, b: fromRecord({a: 4, b: 5}) }) }) );",
		`const a=class{constructor(){}setA(value){this.a=value}setB(value$jscomp$1){this.b=value$jscomp$1}};var JSCompiler_temp_const$jscomp$v1=console;var JSCompiler_temp_const$jscomp$v0=JSCompiler_temp_const$jscomp$v1.log;var JSCompiler_inline_result$jscomp$v2;{var rec$jscomp$inline_0={a:4,b:5};JSCompiler_inline_result$jscomp$v2=(new a).setA(rec$jscomp$inline_0.a).setB(rec$jscomp$inline_0.b)}var JSCompiler_inline_result$jscomp$v3;
{var rec$jscomp$inline_1={a:2,b:JSCompiler_inline_result$jscomp$v2};JSCompiler_inline_result$jscomp$v3=(new a).setA(rec$jscomp$inline_1.a).setB(rec$jscomp$inline_1.b)}var JSCompiler_inline_result$jscomp$v4;{var rec$jscomp$inline_2={a:1,b:JSCompiler_inline_result$jscomp$v3};JSCompiler_inline_result$jscomp$v4=(new a).setA(rec$jscomp$inline_2.a).setB(rec$jscomp$inline_2.b)}JSCompiler_temp_const$jscomp$v0.call(JSCompiler_temp_const$jscomp$v1,JSCompiler_inline_result$jscomp$v4)`)
}

// Java: InlineFunctionsTest.testNullableFunctionsNotInlined()
func (c *InlineFunctionsTest) TestNullableFunctionsNotInlined() {
	c.TestSame("var foo = null; if(x) { foo = function(){}; }; foo();foo();foo();")
	c.TestSame("var foo = null; if(x) { foo = function(){}; }; foo?.();foo?.();foo?.();")
	// Even when non-null, we don't inline if there are multiple definitions of a function (need
	// Dataflow Analysis for that).
	c.TestSame(
		`var foo = function() {return 1;}; if(x) { foo = function(){}; };
 foo?.();foo?.();foo?.();
`)
}

// Java: InlineFunctionsTest.testInlinedAtRegularAndOptChainCallsBoth()
func (c *InlineFunctionsTest) TestInlinedAtRegularAndOptChainCallsBoth() {
	// both regular and optional references are found and inlined
	c.Test("function foo() {return 1;} foo(); foo?.();", "1;1;")
}

// Java: InlineFunctionsTest.testQualifiedNamesNotInlined()
func (c *InlineFunctionsTest) TestQualifiedNamesNotInlined() {
	c.TestSame(
		`var a = { b :  function() {return 0;} };
var c = a.b();
`)
	c.TestSame(
		`var a = { b :  function() {return 0;} };
var c = a.b?.();
`)
	c.TestSame(
		`var a = { b :  function() {return 0;} };
var c = a?.b();
`)
}

// Java: InlineFunctionsTest.testInlineEmptyFunction1()
func (c *InlineFunctionsTest) TestInlineEmptyFunction1() {
	// Empty function, no params.
	c.Test(
		`function foo(){}
foo();
`,
		"void 0;")
}

// Java: InlineFunctionsTest.testInlineEmptyFunction1_optChain()
func (c *InlineFunctionsTest) TestInlineEmptyFunction1_optChain() {
	// Empty function, no params.
	c.Test("function foo(){} foo?.();", "void 0;")
}

// Java: InlineFunctionsTest.testInline_optChain_1()
func (c *InlineFunctionsTest) TestInline_optChain_1() {
	// Empty function, no params.
	c.Test("(function(){}).call?.(this)", "void 0;")
	c.Test("(function(){})?.call(this)", "void 0;")
}

// Java: InlineFunctionsTest.testInlineEmptyFunction2()
func (c *InlineFunctionsTest) TestInlineEmptyFunction2() {
	// Empty function, params with no side-effects.
	c.Test("function foo(){}\n foo(1, new Date, function(){});", "void 0;")
}

// Java: InlineFunctionsTest.testInlineEmptyFunction2_optChain()
func (c *InlineFunctionsTest) TestInlineEmptyFunction2_optChain() {
	c.Test("function foo(){}\n foo?.(1, new Date, function(){});", "void 0;")
}

// Java: InlineFunctionsTest.testInlineEmptyFunction3()
func (c *InlineFunctionsTest) TestInlineEmptyFunction3() {
	// Empty function, multiple references.
	c.Test("function foo(){}\n foo();foo();foo();", "void 0;void 0;void 0")
}

// Java: InlineFunctionsTest.testInlineEmptyFunction4()
func (c *InlineFunctionsTest) TestInlineEmptyFunction4() {
	// Empty function, params with side-effects forces block inlining.
	c.Test("function foo(){}\n foo(x());", "{var JSCompiler_inline_anon_param_v0 = x();}")
}

// Java: InlineFunctionsTest.testInlineEmptyFunction4_optChainCall()
func (c *InlineFunctionsTest) TestInlineEmptyFunction4_optChainCall() {
	// Empty function, params with side-effects forces block inlining.
	c.Test(
		"function foo(){}\n foo?.(x());",
		`var JSCompiler_inline_result$jscomp$v0;
{
  var JSCompiler_inline_anon_param_v1 = x();
  JSCompiler_inline_result$jscomp$v0 = void 0;
}
JSCompiler_inline_result$jscomp$v0;
`)
}

// Java: InlineFunctionsTest.testInlineEmptyFunction6()
func (c *InlineFunctionsTest) TestInlineEmptyFunction6() {
	c.Test("if (window) { f(); function f() {} }", "if (window) { void 0; }")
}

// Java: InlineFunctionsTest.testInlineEmptyFunction6_optChain()
func (c *InlineFunctionsTest) TestInlineEmptyFunction6_optChain() {
	c.Test("if (window) { f?.(); function f() {} }", "if (window) { void 0; }")
}

// Java: InlineFunctionsTest.testInlineFunctions1()
func (c *InlineFunctionsTest) TestInlineFunctions1() {
	// As simple a test as we can get.
	c.Test(
		`function foo(){ return 4 }
foo();
`,
		"4")
}

// Java: InlineFunctionsTest.testInlineFunctions1_optChain()
func (c *InlineFunctionsTest) TestInlineFunctions1_optChain() {
	// As simple a test as we can get.
	c.Test(
		`function foo(){ return 4 }
foo?.();
`,
		"4")
}

// Java: InlineFunctionsTest.testInlineFunctions1_arrowFn()
func (c *InlineFunctionsTest) TestInlineFunctions1_arrowFn() {
	c.Test(
		`const foo = () => 4;
foo();
`,
		"4")
}

// Java: InlineFunctionsTest.testInlineFunctions1_arrowFn_block()
func (c *InlineFunctionsTest) TestInlineFunctions1_arrowFn_block() {
	c.Test(
		`const foo = () => { return 4; };
foo();
`,
		"4")
}

// Java: InlineFunctionsTest.testInlineFunctions1_inArrowFn()
func (c *InlineFunctionsTest) TestInlineFunctions1_inArrowFn() {
	c.Test(
		`const foo = () => 4;
(() => foo())();
`,
		"4")
}

// Java: InlineFunctionsTest.testInlineFunctions1_inArrowFn_block()
func (c *InlineFunctionsTest) TestInlineFunctions1_inArrowFn_block() {
	c.Test(
		`const foo = () => 4;
(() => {
  let a = foo();
  [].forEach(() => a += 1);
  return a;
})();
`,
		// Why is the IIFE not inlined here?
		// This happens whether or not the function references
		// a local.
		`(() => {
  let a = 4;
  [].forEach(() => {
    return a = a + 1;
  });
  return a;
})();
`)
}

// Java: InlineFunctionsTest.testInlineFunctions1_arrowFnArg()
func (c *InlineFunctionsTest) TestInlineFunctions1_arrowFnArg() {
	c.Test(
		`const bar = (f) => f();
bar(() => 4);
`,
		"4")
}

// Java: InlineFunctionsTest.testInlineFunctions2()
func (c *InlineFunctionsTest) TestInlineFunctions2() {
	// inline simple constants
	// NOTE: CD is not inlined.
	c.Test(
		`var t;var AB=function(){return 4};
function BC(){return 6;}
CD=function(x){return x + 5};x=CD(3);y=AB();z=BC();
`,
		"var t;CD=function(x){return x+5};x=CD(3);y=4;z=6")
}

// calls present in the args list of an optional chaining call get inlined.
// Java: InlineFunctionsTest.testInline_simpleConstantsUnderOptionalChaining()
func (c *InlineFunctionsTest) TestInline_simpleConstantsUnderOptionalChaining() {
	// inline simple constants
	c.Test(
		`var AB=function(){return 4};
function BC(x){ return x;}
y=BC?.(AB());
`,
		"y = 4;")
}

// Java: InlineFunctionsTest.testInlineFunctions3()
func (c *InlineFunctionsTest) TestInlineFunctions3() {
	// inline simple constants
	c.Test(
		`var t;var AB=function(){return 4};
function BC(){return 6;}
var CD=function(x){return x + 5};x=CD(3);y=AB();z=BC();
`,
		"var t;x=3+5;y=4;z=6")
}

// Java: InlineFunctionsTest.testInlineFunctions3_optChain()
func (c *InlineFunctionsTest) TestInlineFunctions3_optChain() {
	// inline simple constants
	c.Test(
		`var t;var AB=function(){return 4};
function BC(){return 6;}
var CD=function(x){return x + 5};x=CD?.(3);y=AB?.();z=BC?.();
`,
		"var t;x=3+5;y=4;z=6")
}

// Java: InlineFunctionsTest.testInlineFunctions4()
func (c *InlineFunctionsTest) TestInlineFunctions4() {
	// don't inline if there are multiple definitions (need DFA for that).
	c.Test(
		`var t; var AB = function() { return 4 };
function BC() { return 6; }
CD = 0;
CD = function(x) { return x + 5 }; x = CD(3); y = AB(); z = BC();
`,
		"var t;CD=0;CD=function(x){return x+5};x=CD(3);y=4;z=6")
}

// Java: InlineFunctionsTest.testInlineFunctions4_optChain()
func (c *InlineFunctionsTest) TestInlineFunctions4_optChain() {
	// don't inline if there are multiple definitions (need DFA for that).
	c.Test(
		`var t; var AB = function() { return 4 };
function BC() { return 6; }
CD = 0;
CD = function(x) { return x + 5 }; x = CD?.(3); y = AB?.(); z = BC?.();
`,
		"var t;CD=0;CD=function(x){return x+5};x=CD?.(3);y=4;z=6")
}

// Java: InlineFunctionsTest.testInlineFunctions5()
func (c *InlineFunctionsTest) TestInlineFunctions5() {
	// inline additions
	c.Test(
		`var FOO_FN=function(x,y) { return "de" + x + "nu" + y };
var a = FOO_FN("ez", "ts")
`,
		"var a=\"de\"+\"ez\"+\"nu\"+\"ts\"")
}

// Java: InlineFunctionsTest.testInlineFunctions5_optChain()
func (c *InlineFunctionsTest) TestInlineFunctions5_optChain() {
	// inline additions
	c.Test(
		`var FOO_FN=function(x,y) { return "de" + x + "nu" + y };
var a = FOO_FN?.("ez", "ts")
`,
		"var a=\"de\"+\"ez\"+\"nu\"+\"ts\"")
}

// Java: InlineFunctionsTest.testInlineFunctions6()
func (c *InlineFunctionsTest) TestInlineFunctions6() {
	// more complex inlines
	c.Test(
		`function BAR_FN(x, y, z) { return z(nochg(x + y)) }
alert(BAR_FN(1, 2, baz))
`,
		"alert(baz(nochg(1+2)))")
}

// Java: InlineFunctionsTest.testInlineFunctions6_optChain()
func (c *InlineFunctionsTest) TestInlineFunctions6_optChain() {
	// more complex inlines
	c.Test(
		`function BAR_FN(x, y, z) { return z?.(nochg?.(x + y)) }
alert(BAR_FN?.(1, 2, baz))
`,
		"alert(baz?.(nochg?.(1+2)))")
}

// Java: InlineFunctionsTest.testInlineFunctions7()
func (c *InlineFunctionsTest) TestInlineFunctions7() {
	// inlines appearing multiple times
	c.Test(
		`function FN(x,y,z){return x+x+y}
var b=FN(1,2,3)
`,
		"var b=1+1+2")
}

// Java: InlineFunctionsTest.testInlineFunctions7_optChain()
func (c *InlineFunctionsTest) TestInlineFunctions7_optChain() {
	// inlines appearing multiple times
	c.Test(
		`function FN(x,y,z){return x+x+y}
var b=FN?.(1,2,3)
`,
		"var b=1+1+2")
}

// Java: InlineFunctionsTest.testInlineFunctions8()
func (c *InlineFunctionsTest) TestInlineFunctions8() {
	// check correct parenthesization
	c.Test(
		`function MUL(x,y){return x*y}function ADD(x,y){return x+y}
var a=1+MUL(2,3);var b=2*ADD(3,4)
`,
		"var a=1+2*3;var b=2*(3+4)")
}

// Java: InlineFunctionsTest.testInlineFunctions9()
func (c *InlineFunctionsTest) TestInlineFunctions9() {
	// don't inline if the input parameter is modified.
	c.Test(
		`function INC(x){return x++}
var y=INC(i)
`,
		`var y;{var x$jscomp$inline_0=i;
y=x$jscomp$inline_0++}
`)
}

// Java: InlineFunctionsTest.testInlineFunctions9_optChain()
func (c *InlineFunctionsTest) TestInlineFunctions9_optChain() {
	// don't inline if the input parameter is modified.
	c.Test(
		"function INC(x){return x++} var y=INC?.(i)",
		"var y;{var x$jscomp$inline_0=i; y=x$jscomp$inline_0++}")
}

// Java: InlineFunctionsTest.testInlineFunctions10()
func (c *InlineFunctionsTest) TestInlineFunctions10() {
	c.Test(
		`function INC(x){return x++}
var y = INC(i);
y = INC(i);
`,
		`var y;
{
  var x$jscomp$inline_0 = i;
  y = x$jscomp$inline_0++;
}
{
  var x$jscomp$inline_1 = i;
  y = x$jscomp$inline_1++;
}
`)
}

// Java: InlineFunctionsTest.testInlineFunctions11()
func (c *InlineFunctionsTest) TestInlineFunctions11() {
	c.Test(
		`function f(x){return x}
var y=f(i)
`,
		"var y=i")
}

// Java: InlineFunctionsTest.testInlineFunctions13()
func (c *InlineFunctionsTest) TestInlineFunctions13() {
	// inline as block if the input parameter has side-effects.
	c.Test(
		`function f(x){return x}
var y=f(i++)
`,
		// TODO: b/407603216 - Inline using DIRECT mode when return expressions don't have
		// references that can be side-effected. The output code can be `var y=i++`
		`var y = i++;
`)
}

// Java: InlineFunctionsTest.testInlineFunctions13a()
func (c *InlineFunctionsTest) TestInlineFunctions13a() {
	// inline as block if the input parameter has side-effects.
	c.Test(
		`function f(x){return random() || x}
var y=f(i++)
`,
		"var y;{var x$jscomp$inline_0=i++;y=random() || x$jscomp$inline_0}")
}

// Java: InlineFunctionsTest.testInlineFunctions14()
func (c *InlineFunctionsTest) TestInlineFunctions14() {
	// don't remove functions that are referenced on other ways
	c.Test(
		`function FOO(x){return x}
var BAR=function(y){return y};
b=FOO;a(BAR);x=FOO(1);y=BAR(2)
`,
		`function FOO(x){return x}
var BAR=function(y){return y};
b=FOO;a(BAR);x=1;y=2
`)
}

// Java: InlineFunctionsTest.testInlineFunctions15a()
func (c *InlineFunctionsTest) TestInlineFunctions15a() {
	// closure factories: do inline into global scope.
	c.Test(
		`function foo(){return function(a){return a+1}}
var b=function(){return c};
var d=b()+foo()
`,
		"var d=c+function(a){return a+1}")
}

// Java: InlineFunctionsTest.testInlineFunctions15b()
func (c *InlineFunctionsTest) TestInlineFunctions15b() {
	c.assumeMinimumCapture = false
	// closure factories: don't inline closure with locals into global scope.
	c.Test(
		`function foo(){var x; return function(a){ return a+1; };}
var b = function() { return c; };
var d = b() + foo()
`,
		"function foo() { var x; return function(a){return a+1}; } var d = c+foo();")
	c.assumeMinimumCapture = true
	c.Test(
		`function foo() { var x; return function(a) { return a + 1; }; }
var b = function() { return c; };
var d = b() + foo()
`,
		`var JSCompiler_inline_result$jscomp$v0;
{
  var x$jscomp$inline_0;
  JSCompiler_inline_result$jscomp$v0 = function(a$jscomp$inline_1) {
    return a$jscomp$inline_1 + 1;
  };
}
var d = c+JSCompiler_inline_result$jscomp$v0;
`)
}

// Java: InlineFunctionsTest.testInlineFunctions15c()
func (c *InlineFunctionsTest) TestInlineFunctions15c() {
	c.assumeMinimumCapture = false
	// closure factories: don't inline into non-global scope.
	c.Test(
		`function foo(){return function(a){return a+1}}
var b=function(){return c};
function _x(){ var d=b()+foo() }
`,
		`function foo(){return function(a){return a+1}}
function _x(){ var d=c+foo() }
`)
	c.assumeMinimumCapture = true
	// closure factories: don't inline into non-global scope.
	c.Test(
		`function foo(){return function(a){return a+1}}
var b=function(){return c};
function _x(){ var d=b()+foo() }
`,
		"function _x(){var d=c+function(a){return a+1}}")
}

// Java: InlineFunctionsTest.testInlineFunctions15d()
func (c *InlineFunctionsTest) TestInlineFunctions15d() {
	c.assumeMinimumCapture = false
	// closure factories: don't inline functions with vars.
	c.Test(
		`function foo() { var x; return function(a) { return a+1; }; }
var b = function() { return c; };
function _x(){ var d = b()+foo() }
`,
		`function foo() {
  var x;
  return function(a){ return a+1; };
}
function _x() { var d = c + foo(); }
`)
	c.assumeMinimumCapture = true
	// closure factories: inline functions with vars.
	c.Test(
		`function foo(){var x; return function(a){return a+1}}
var b = function(){return c};
function _x(){ var d=b()+foo() }
`,
		`function _x() {
  var JSCompiler_inline_result$jscomp$v0;
  {
    var x$jscomp$inline_0;
    JSCompiler_inline_result$jscomp$v0 = function(a$jscomp$inline_1) {
      return a$jscomp$inline_1+1
    };
  }
  var d = c + JSCompiler_inline_result$jscomp$v0;
}
`)
}

// Java: InlineFunctionsTest.testInlineFunctions16a()
func (c *InlineFunctionsTest) TestInlineFunctions16a() {
	c.assumeMinimumCapture = false
	c.TestSame("function foo(b){return window.bar(function(){c(b)})} var d=foo(e)")
	c.assumeMinimumCapture = true
	c.Test(
		`function foo(b){return window.bar(function(){c(b)})}
var d=foo(e)
`,
		`var d;{var b$jscomp$inline_0=e;
d=window.bar(function(){c(b$jscomp$inline_0)})}
`)
}

// Java: InlineFunctionsTest.testInlineFunctions16b()
func (c *InlineFunctionsTest) TestInlineFunctions16b() {
	c.Test(
		`function foo(){return window.bar(function(){c()})}
var d=foo(e)
`,
		"var d=window.bar(function(){c()})")
}

// Java: InlineFunctionsTest.testInlineFunctions17()
func (c *InlineFunctionsTest) TestInlineFunctions17() {
	// don't inline recursive functions
	c.TestSame("function foo(x){return x*x+foo(3)}var bar=foo(4)")
}

// Java: InlineFunctionsTest.testInlineFunctions19()
func (c *InlineFunctionsTest) TestInlineFunctions19() {
	// TRICKY ... test nested inlines
	// with block inlining possible
	c.Test(
		`function foo(a, b){return a+b}
function bar(d){return c}
var d=foo(bar(1),e)
`,
		"var d=c+e;")
}

// Java: InlineFunctionsTest.testInlineFunctions19_optChain()
func (c *InlineFunctionsTest) TestInlineFunctions19_optChain() {
	// TRICKY ... test nested inlines
	// with block inlining possible
	c.Test(
		`function foo(a, b){return a+b}
function bar(d){return c}
var d=foo?.(bar?.(1),e)
`,
		"var d=c+e;")
}

// Java: InlineFunctionsTest.testInlineFunctions21()
func (c *InlineFunctionsTest) TestInlineFunctions21() {
	// with block inlining possible
	c.Test(
		`function foo(a, b){return a+b}
function bar(d){return c}
var d=bar(foo(1,e))
`,
		"var d=c")
}

// Java: InlineFunctionsTest.testInlineFunctions22()
func (c *InlineFunctionsTest) TestInlineFunctions22() {
	// Another tricky case ... test nested compiler inlines
	c.Test(
		`function plex(a){if(a) return 0;else return 1;}
function foo(a, b){return bar(a+b)}
function bar(d){return plex(d)}
var d=foo(1,2)
`,
		`var d;
{
  JSCompiler_inline_label_plex_v0:{
    if(1+2) {
      d = 0;
      break JSCompiler_inline_label_plex_v0;
    } else {
      d = 1;
      break JSCompiler_inline_label_plex_v0;
    }
    d=void 0;
  }
}
`)
}

// Java: InlineFunctionsTest.testInlineFunctions23()
func (c *InlineFunctionsTest) TestInlineFunctions23() {
	// Test both orderings again
	c.Test(
		`function complex(a){if(a) return 0;else return 1;}
function bar(d){return complex(d)}
function foo(a, b){return bar(a+b)}
var d=foo(1,2)
`,
		`var d;
{
  JSCompiler_inline_label_complex_v0:{
    if (1+2) {
      d=0;
      break JSCompiler_inline_label_complex_v0;
    } else {
      d = 1;
      break JSCompiler_inline_label_complex_v0;
    }
    d=void 0
  }
}
`)
}

// Java: InlineFunctionsTest.testInlineFunctions24()
func (c *InlineFunctionsTest) TestInlineFunctions24() {
	// Don't inline functions with 'arguments' or 'this'
	c.TestSame("function foo(x){return this}foo(1)")
}

// Java: InlineFunctionsTest.testInlineFunctions26()
func (c *InlineFunctionsTest) TestInlineFunctions26() {
	// Don't inline external functions
	c.TestSame("function _foo(x){return x}_foo(1)")
}

// Java: InlineFunctionsTest.testInlineFunctions27()
func (c *InlineFunctionsTest) TestInlineFunctions27() {
	c.Test(
		`var window = {}; function foo(){window.bar++; return 3;}
var x = {y: 1, z: foo(2)};
`,
		`var window={};
var JSCompiler_inline_result$jscomp$v0;
{
  window.bar++;
  JSCompiler_inline_result$jscomp$v0 = 3;
}
var x = {y: 1, z: JSCompiler_inline_result$jscomp$v0};
`)
}

// Java: InlineFunctionsTest.testInlineFunctions28()
func (c *InlineFunctionsTest) TestInlineFunctions28() {
	c.Test(
		`var window = {}; function foo(){window.bar++; return 3;}
var x = {y: alert(), z: foo(2)};
`,
		`var window = {};
var JSCompiler_temp_const$jscomp$v0 = alert();
var JSCompiler_inline_result$jscomp$v1;
{
 window.bar++;
 JSCompiler_inline_result$jscomp$v1 = 3;}
var x = {
  y: JSCompiler_temp_const$jscomp$v0,
  z: JSCompiler_inline_result$jscomp$v1
};
`)
}

// Java: InlineFunctionsTest.testInlineFunctions29()
func (c *InlineFunctionsTest) TestInlineFunctions29() {
	c.Test(
		`var window = {}; function foo(){window.bar++; return 3;}
var x = {a: alert(), b: alert2(), c: foo(2)};
`,
		`var window = {};
var JSCompiler_temp_const$jscomp$v1 = alert();
var JSCompiler_temp_const$jscomp$v0 = alert2();
var JSCompiler_inline_result$jscomp$v2;
{
 window.bar++;
 JSCompiler_inline_result$jscomp$v2 = 3;}
var x = {
  a: JSCompiler_temp_const$jscomp$v1,
  b: JSCompiler_temp_const$jscomp$v0,
  c: JSCompiler_inline_result$jscomp$v2
};
`)
}

// Java: InlineFunctionsTest.testInlineFunctions30()
func (c *InlineFunctionsTest) TestInlineFunctions30() {
	// As simple a test as we can get.
	c.TestSame("function foo(){ return eval() } foo();")
}

// Java: InlineFunctionsTest.testInlineFunctions31()
func (c *InlineFunctionsTest) TestInlineFunctions31() {
	// Don't introduce a duplicate label in the same scope
	c.Test("function foo(){ lab:{4;} } lab:{foo();}", "lab:{{JSCompiler_inline_label_0:{4}}}")
}

// Java: InlineFunctionsTest.testInlineFunctions32()
func (c *InlineFunctionsTest) TestInlineFunctions32() {
	c.Test("let f = function() { return 5; }; f();", "5;")
}

// Java: InlineFunctionsTest.testInlineFunctions33()
func (c *InlineFunctionsTest) TestInlineFunctions33() {
	c.Test("const f = function() { return 5; }; f();", "5;")
}

// Java: InlineFunctionsTest.testInlineFunctions34()
func (c *InlineFunctionsTest) TestInlineFunctions34() {
	c.Test(
		`class X {}
(function(e) {
  for (var el = f(e.target); el != null; el = el.parent) {}
  function f(x) { return x instanceof X ? x : null; }
})({target:{}});
`,
		`class X {}
{
  var e$jscomp$inline_1 = {target:{}};
  var el$jscomp$inline_2;
  {
    var x$jscomp$inline_3 = e$jscomp$inline_1.target;
    el$jscomp$inline_2 = x$jscomp$inline_3 instanceof X ? x$jscomp$inline_3 : null;
  }
  for(;el$jscomp$inline_2 != null; el$jscomp$inline_2 = el$jscomp$inline_2.parent);
}
`)
}

// Java: InlineFunctionsTest.testInlineFunctions34_optChain()
func (c *InlineFunctionsTest) TestInlineFunctions34_optChain() {
	c.Test(
		`class X {}
(function(e) {
  for (var el = f?.(e.target); el != null; el = el.parent) {}
  function f(x) { return x instanceof X ? x : null; }
})({target:{}});
`,
		`class X {}
{
  var e$jscomp$inline_1 = {target:{}};
  var el$jscomp$inline_2;
  {
    var x$jscomp$inline_3 = e$jscomp$inline_1.target;
    el$jscomp$inline_2 = x$jscomp$inline_3 instanceof X ? x$jscomp$inline_3 : null;
  }
  for(;el$jscomp$inline_2 != null; el$jscomp$inline_2 = el$jscomp$inline_2.parent);
}
`)
}

// Same as above, except the for loop uses a let instead of a var. See b/73373371.
// Java: InlineFunctionsTest.testInlineFunctions35()
func (c *InlineFunctionsTest) TestInlineFunctions35() {
	c.Test(
		`class X {}
(function(e) {
  for (let el = f(e.target); el != null; el = el.parent) {}
  function f(x) { return x instanceof X ? x : null; }
})({target:{}});
`,
		`class X {}
{
  var e$jscomp$inline_1 = {target: {}};
  var JSCompiler_inline_result$jscomp$inline_2;
  {
    var x$jscomp$inline_3 = e$jscomp$inline_1.target;
    JSCompiler_inline_result$jscomp$inline_2 =
        x$jscomp$inline_3 instanceof X ? x$jscomp$inline_3 : null
  }
  for (let el$jscomp$inline_4 = JSCompiler_inline_result$jscomp$inline_2;
       el$jscomp$inline_4 != null;
       el$jscomp$inline_4 = el$jscomp$inline_4.parent) {}
}
`)
}

// Java: InlineFunctionsTest.testInlineFunctions35_optChain()
func (c *InlineFunctionsTest) TestInlineFunctions35_optChain() {
	c.Test(
		"class X {}\n"+
			"(function(e) {\n"+
			"// The for loop uses a `let` instead of a var.\n"+
			"  for (let el = f?.(e.target); el != null; el = el.parent) {}\n"+
			"  function f(x) { return x instanceof X ? x : null; }\n"+
			"})({target:{}});\n",
		`class X {}
{
  var e$jscomp$inline_1 = {target: {}};
  var JSCompiler_inline_result$jscomp$inline_2;
  {
    var x$jscomp$inline_3 = e$jscomp$inline_1.target;
    JSCompiler_inline_result$jscomp$inline_2 =
        x$jscomp$inline_3 instanceof X ? x$jscomp$inline_3 : null
  }
  for (let el$jscomp$inline_4 = JSCompiler_inline_result$jscomp$inline_2;
       el$jscomp$inline_4 != null;
       el$jscomp$inline_4 = el$jscomp$inline_4.parent) {}
}
`)
}

// Java: InlineFunctionsTest.testDontInlineFunctionsWithDirectArgumentsReferences()
func (c *InlineFunctionsTest) TestDontInlineFunctionsWithDirectArgumentsReferences() {
	c.TestSame("function foo() { return arguments[0]; } foo(1);")
}

// Java: InlineFunctionsTest.testInlineSwitchStatementWithReturns()
func (c *InlineFunctionsTest) TestInlineSwitchStatementWithReturns() {
	c.Test(
		`function foo(x) {
  switch (x) {
    case 1: return 2;
    case 2: return 3;
    default: return 4;
  }
}
function bar() {
  const f = foo(x);
  console.log(f);
}
bar();
`,
		`{
  var JSCompiler_inline_result$jscomp$v1;
  {
    JSCompiler_inline_label_foo_v2: {
      switch(x) {
        case 1:
          JSCompiler_inline_result$jscomp$v1 = 2;
          break JSCompiler_inline_label_foo_v2;
        case 2:
          JSCompiler_inline_result$jscomp$v1 = 3;
          break JSCompiler_inline_label_foo_v2;
        default:
          JSCompiler_inline_result$jscomp$v1 = 4;
          break JSCompiler_inline_label_foo_v2;
      }
      JSCompiler_inline_result$jscomp$v1 = void 0;
    }
  }
  const f$jscomp$inline_0 = JSCompiler_inline_result$jscomp$v1;
  console.log(f$jscomp$inline_0);
}
`)
}

// Java: InlineFunctionsTest.testDontInlineFunctionsWithArgumentsReferencesInArrowFunction()
func (c *InlineFunctionsTest) TestDontInlineFunctionsWithArgumentsReferencesInArrowFunction() {
	c.TestSame("function foo() { return () => arguments[0]; } foo(1);")
}

// Java: InlineFunctionsTest.testDoInlineFunctionsWithArgumentsReferencesInInnerVanillaFunction()
func (c *InlineFunctionsTest) TestDoInlineFunctionsWithArgumentsReferencesInInnerVanillaFunction() {
	c.Test(
		"function foo() { return function() { return arguments[0]; } } foo(1);",
		"(function() { return arguments[0]; })")
}

// Java: InlineFunctionsTest.testMixedModeInlining1()
func (c *InlineFunctionsTest) TestMixedModeInlining1() {
	// Base line tests, direct inlining
	c.Test(
		`function foo(){return 1}
foo();
`,
		"1;")
}

// Java: InlineFunctionsTest.testMixedModeInlining2()
func (c *InlineFunctionsTest) TestMixedModeInlining2() {
	// Base line tests, block inlining. Block inlining is needed by possible-side-effect parameter.
	c.Test("function foo(){return 1} foo(x());", "{var JSCompiler_inline_anon_param_v0=x();1}")
}

// Java: InlineFunctionsTest.testMixedModeInlining2_optChain()
func (c *InlineFunctionsTest) TestMixedModeInlining2_optChain() {
	// Base line tests, block inlining. Block inlining is needed by possible-side-effect parameter.
	c.Test(
		"function foo(){return 1} foo?.(x());",
		`var JSCompiler_inline_result$jscomp$v0;
{
  var JSCompiler_inline_anon_param_v1 = x();
  JSCompiler_inline_result$jscomp$v0 = 1;
}
JSCompiler_inline_result$jscomp$v0;
`)
}

// Java: InlineFunctionsTest.testMixedModeInlining3()
func (c *InlineFunctionsTest) TestMixedModeInlining3() {
	// Inline using both modes.
	c.Test(
		"function foo(){return 1} foo();foo(x());",
		"1;{var JSCompiler_inline_anon_param_v0=x();1}")
}

// Java: InlineFunctionsTest.testMixedModeInlining3_optChain()
func (c *InlineFunctionsTest) TestMixedModeInlining3_optChain() {
	// Inline using both modes.
	c.Test(
		"function foo(){return 1} foo?.();foo?.(x?.());",
		`1;
var JSCompiler_inline_result$jscomp$v0;
{
  var JSCompiler_inline_anon_param_v1 = x?.();
  JSCompiler_inline_result$jscomp$v0 = 1;
}
JSCompiler_inline_result$jscomp$v0;
`)
}

// Java: InlineFunctionsTest.testMixedModeInlining4()
func (c *InlineFunctionsTest) TestMixedModeInlining4() {
	// Inline using both modes. Alternating. Second call of each type has
	// side-effect-less parameter, this is thrown away.
	c.Test(
		"function foo(){return 1} foo(); foo(x()); foo(1);foo(1,x());",
		`1;
{
  var JSCompiler_inline_anon_param_v0 = x();
  1;
}
1;
{
  var JSCompiler_inline_anon_param_v3 = 1;
  var JSCompiler_inline_anon_param_v4 = x();
  1;
}
`)
}

// Java: InlineFunctionsTest.testMixedModeInliningCosting1()
func (c *InlineFunctionsTest) TestMixedModeInliningCosting1() {
	// Inline using both modes. Costing estimates.
	// Base line.
	c.Test(
		`function foo(a,b){return a+b+a+b+4+5+6+7+8+9+1+2+3+4+5}
foo(1,2);
foo(2,3)
`,
		`1+2+1+2+4+5+6+7+8+9+1+2+3+4+5;
2+3+2+3+4+5+6+7+8+9+1+2+3+4+5
`)
}

// Java: InlineFunctionsTest.testMixedModeInliningCosting2()
func (c *InlineFunctionsTest) TestMixedModeInliningCosting2() {
	// Don't inline here because the function definition can not be eliminated.
	// TODO(johnlenz): Should we add constant removing to the unit test?
	c.TestSame(
		`function foo(a,b){return a+b+a+b+4+5+6+7+8+9+1+2+3+4+5}
foo(1,2);
foo(2,3,x())
`)
}

// Java: InlineFunctionsTest.testMixedModeInliningCosting3()
func (c *InlineFunctionsTest) TestMixedModeInliningCosting3() {
	// Do inline here because the function definition can be eliminated.
	c.Test(
		`function foo(a,b){return a+b+a+b+4+5+6+7+8+9+1+2+3+10}
foo(1,2);
foo(2,3,x());
`,
		`1+2+1+2+4+5+6+7+8+9+1+2+3+10; // inlining of the foo(1,2); call
{
var a$jscomp$inline_0 = 2; // temp for arg 2
var b$jscomp$inline_1 = 3; // temp for arg 3
  var JSCompiler_inline_anon_param_v0 = x();
  a$jscomp$inline_0 + b$jscomp$inline_1 + a$jscomp$inline_0 + b$jscomp$inline_1 + 4 + 5 + 6 + 7 + 8 + 9 + 1 + 2 + 3 + 10;
}
`)
}

// Java: InlineFunctionsTest.testMixedModeInliningCosting4()
func (c *InlineFunctionsTest) TestMixedModeInliningCosting4() {
	// Threshold test.
	c.TestSame(
		`function foo(a,b){return a+b+a+b+4+5+6+7+8+9+1+2+3+4+101}
foo(1,2);
foo(2,3,x())
`)
}

// See b/72513540
// Java: InlineFunctionsTest.testDestructuringAssignInFunction()
func (c *InlineFunctionsTest) TestDestructuringAssignInFunction() {
	c.Test(
		`function foo(a) {
  [a] = [1];
}
foo(2);
`,
		`{
  var a$jscomp$inline_0 = 2;
  [a$jscomp$inline_0] = [1];
}
`)
}

// Java: InlineFunctionsTest.testDestructuringAssignInFunction_withArrayPattern_doesNotCrash()
func (c *InlineFunctionsTest) TestDestructuringAssignInFunction_withArrayPattern_doesNotCrash() {
	c.Test(
		`var Di = I(() => {
  function zv() {
    JSCOMPILER_PRESERVE(e), [f] = CN();
  }
  function CN() {
    return [1];
  }
});
`,
		"var Di = I(() => {});")
}

// Java: InlineFunctionsTest.testDestructuringAssignInFunction_withObjectPattern_doesNotCrash()
func (c *InlineFunctionsTest) TestDestructuringAssignInFunction_withObjectPattern_doesNotCrash() {
	c.Test(
		`var Di = I(() => {
  function zv() {
    JSCOMPILER_PRESERVE(e), ({f: g} = CN());
  }
  function CN() {
    return {f: 1};
  }
});
`,
		"var Di = I(() => {});")
}

// Java: InlineFunctionsTest.testNoInlineIfParametersModified1()
func (c *InlineFunctionsTest) TestNoInlineIfParametersModified1() {
	// Assignment
	c.Test(
		"function f(x){return x=1}f(undefined)",
		`{var x$jscomp$inline_0=undefined;
x$jscomp$inline_0=1}
`)
}

// Java: InlineFunctionsTest.testNoInlineIfParametersModified2()
func (c *InlineFunctionsTest) TestNoInlineIfParametersModified2() {
	c.Test(
		"function f(x){return (x)=1;}f(2)",
		`{var x$jscomp$inline_0=2;
x$jscomp$inline_0=1}
`)
}

// Java: InlineFunctionsTest.testNoInlineIfParametersModified3()
func (c *InlineFunctionsTest) TestNoInlineIfParametersModified3() {
	// Assignment variant.
	c.Test(
		"function f(x){return x*=2}f(2)",
		`{var x$jscomp$inline_0=2;
x$jscomp$inline_0*=2}
`)
}

// Java: InlineFunctionsTest.testNoInlineIfParametersModified4()
func (c *InlineFunctionsTest) TestNoInlineIfParametersModified4() {
	// Assignment in if.
	c.Test(
		"function f(x){return x?(x=2):0}f(2)",
		`{var x$jscomp$inline_0=2;
x$jscomp$inline_0?(
x$jscomp$inline_0=2):0}
`)
}

// Java: InlineFunctionsTest.testNoInlineIfParametersModified5()
func (c *InlineFunctionsTest) TestNoInlineIfParametersModified5() {
	// Assignment in if, multiple params
	c.Test(
		"function f(x,y){return x?(y=2):0}f(2,undefined)",
		`{var x$jscomp$inline_0 = 2; var y$jscomp$inline_1=undefined;
x$jscomp$inline_0 ? y$jscomp$inline_1 = 2 : 0;}
`)
}

// Java: InlineFunctionsTest.testNoInlineIfParametersModified6()
func (c *InlineFunctionsTest) TestNoInlineIfParametersModified6() {
	c.Test(
		"function f(x,y){return x?(y=2):0}f(2)",
		`{var x$jscomp$inline_0 = 2; var y$jscomp$inline_1=void 0;
x$jscomp$inline_0 ? y$jscomp$inline_1 = 2 : 0;}
`)
}

// Java: InlineFunctionsTest.testNoInlineIfParametersModified7()
func (c *InlineFunctionsTest) TestNoInlineIfParametersModified7() {
	// Increment
	c.Test(
		"function f(a){return++a<++a}f(1)",
		`{var a$jscomp$inline_0=1;
++a$jscomp$inline_0<
++a$jscomp$inline_0}
`)
}

// Java: InlineFunctionsTest.testInlineIfParametersModified8()
func (c *InlineFunctionsTest) TestInlineIfParametersModified8() {
	// OK, object parameter modified.
	c.Test("function f(a){return a.x=2}f(o)", "o.x=2")
}

// Java: InlineFunctionsTest.testInlineIfParametersModified9()
func (c *InlineFunctionsTest) TestInlineIfParametersModified9() {
	// OK, array parameter modified.
	c.Test("function f(a){return a[2]=2}f(o)", "o[2]=2")
}

// Java: InlineFunctionsTest.testInlineNeverPartialSubtitution1()
func (c *InlineFunctionsTest) TestInlineNeverPartialSubtitution1() {
	c.Test("function f(z){return x.y.z;}f(1)", "x.y.z")
}

// Java: InlineFunctionsTest.testInlineNeverPartialSubtitution2()
func (c *InlineFunctionsTest) TestInlineNeverPartialSubtitution2() {
	c.Test("function f(z){return x.y[z];}f(a)", "x.y[a]")
}

// Java: InlineFunctionsTest.testInlineNeverMutateConstants()
func (c *InlineFunctionsTest) TestInlineNeverMutateConstants() {
	c.Test(
		"function f(x){return x=1}f(undefined)",
		`{var x$jscomp$inline_0=undefined;
x$jscomp$inline_0=1}
`)
}

// Java: InlineFunctionsTest.testInlineNeverOverrideNewValues()
func (c *InlineFunctionsTest) TestInlineNeverOverrideNewValues() {
	c.Test(
		"function f(a){return++a<++a}f(1)",
		`{var a$jscomp$inline_0=1;
++a$jscomp$inline_0<++a$jscomp$inline_0}
`)
}

// Java: InlineFunctionsTest.testInlineMutableArgsReferencedOnce()
func (c *InlineFunctionsTest) TestInlineMutableArgsReferencedOnce() {
	c.Test("function foo(x){return x;}foo([])", "[]")
}

// Java: InlineFunctionsTest.testInlineMutableArgsReferencedOnce2()
func (c *InlineFunctionsTest) TestInlineMutableArgsReferencedOnce2() {
	c.assumeMinimumCapture = true
	// Don't inline a mutable value that will be reused.
	c.Test(
		"function foo(x){return function(){ return x; }} repeat(foo([]))",
		`var JSCompiler_inline_result$jscomp$v0;
{
var x$jscomp$inline_0=[];
JSCompiler_inline_result$jscomp$v0=function(){return x$jscomp$inline_0};
}
repeat(JSCompiler_inline_result$jscomp$v0)
`)
}

// Java: InlineFunctionsTest.testInlineMutableArgsReferencedOnce3()
func (c *InlineFunctionsTest) TestInlineMutableArgsReferencedOnce3() {
	c.assumeMinimumCapture = true
	// Don't inline a mutable value that will be reused.
	c.Test(
		`function f(a) {
  for(var i=0; i<0; i++) {
    g(a);
  }
}
f([]);
`,
		`{
  var a$jscomp$inline_0 = [];
  var i$jscomp$inline_1 = 0;
  for(; i$jscomp$inline_1 < 0; i$jscomp$inline_1++) {
    g(a$jscomp$inline_0)
  }
}
`)
}

// Java: InlineFunctionsTest.testInlineBlockMutableArgs1()
func (c *InlineFunctionsTest) TestInlineBlockMutableArgs1() {
	c.Test(
		"function foo(x){x+x}foo([])",
		`{var x$jscomp$inline_0=[];
x$jscomp$inline_0+x$jscomp$inline_0}
`)
}

// Java: InlineFunctionsTest.testInlineBlockMutableArgs2()
func (c *InlineFunctionsTest) TestInlineBlockMutableArgs2() {
	c.Test(
		"function foo(x){x+x}foo(new Date)",
		`{var x$jscomp$inline_0=new Date;
x$jscomp$inline_0+x$jscomp$inline_0}
`)
}

// Java: InlineFunctionsTest.testInlineBlockMutableArgs3()
func (c *InlineFunctionsTest) TestInlineBlockMutableArgs3() {
	c.Test(
		"function foo(x){x+x}foo(true&&new Date)",
		`{var x$jscomp$inline_0=true&&new Date;
x$jscomp$inline_0+x$jscomp$inline_0}
`)
}

// Java: InlineFunctionsTest.testInlineBlockMutableArgs4()
func (c *InlineFunctionsTest) TestInlineBlockMutableArgs4() {
	c.Test(
		"function foo(x){x+x}foo({})",
		`{var x$jscomp$inline_0={};
x$jscomp$inline_0+x$jscomp$inline_0}
`)
}

// Java: InlineFunctionsTest.testShadowVariables1()
func (c *InlineFunctionsTest) TestShadowVariables1() {
	// The Normalize pass now guarantees that that globals are never shadowed
	// by locals.
	// "foo" is inlined here as its parameter "a" doesn't conflict.
	// "bar" is assigned a new name.
	c.Test(
		`var a = 0;
function foo(a) { return 3+a; }
function bar() { var a = foo(4); }
bar();
`,
		"var a=0; {var a$jscomp$inline_0 = 3+4; }")
}

// Java: InlineFunctionsTest.testShadowVariables2()
func (c *InlineFunctionsTest) TestShadowVariables2() {
	// "foo" is inlined here as its parameter "a" doesn't conflict.
	// "bar" is inlined as its uses global "a", and does introduce any new
	// globals.
	c.Test(
		`var a=0;
function foo(a){return 3+a}
function bar(){a=foo(4)}
bar()
`,
		`var a=0;
{a=3+4}
`)
}

// Java: InlineFunctionsTest.testShadowVariables3()
func (c *InlineFunctionsTest) TestShadowVariables3() {
	// "foo" is inlined into exported "_bar", aliasing foo's "a".
	c.Test(
		`var a = 0;
function foo() { var a = 2; return 3+a; }
function _bar() { a = foo(); }
`,
		`var a=0;
function _bar() {
  {
     var a$jscomp$inline_0 = 2;
     a = 3 + a$jscomp$inline_0;
  }
}
`)
}

// Java: InlineFunctionsTest.testShadowVariables4()
func (c *InlineFunctionsTest) TestShadowVariables4() {
	// "foo" is inlined.
	// block access to global "a".
	c.Test(
		`var a=0;
function foo(){return 3+a}
function _bar(a){a=foo(4)+a}
`,
		`var a=0;function _bar(a$jscomp$1){
a$jscomp$1=
3+a+a$jscomp$1}
`)
}

// Java: InlineFunctionsTest.testShadowVariables6()
func (c *InlineFunctionsTest) TestShadowVariables6() {
	c.Test(
		`var a=0;
function foo() { var a = 4; return 3+a; }
function _bar(a) { a = foo(4); }
`,
		`var a=0;
function _bar(a$jscomp$2) {
  {
    var a$jscomp$inline_0 = 4;
    a$jscomp$2 = 3 + a$jscomp$inline_0;
  }
}
`)
}

// Java: InlineFunctionsTest.testShadowVariables7()
func (c *InlineFunctionsTest) TestShadowVariables7() {
	c.assumeMinimumCapture = false
	c.Test(
		`var a=3;
function foo(){return a}
(function(){var a=5;(function(){foo()})()})()
`,
		`var a=3;
{var a$jscomp$inline_0=5;{a}}
`)
	c.assumeMinimumCapture = true
	c.Test(
		`var a=3;
function foo(){return a}
(function(){var a=5;(function(){foo()})()})()
`,
		`var a=3;
{var a$jscomp$inline_0=5;{a}}
`)
}

// Java: InlineFunctionsTest.testShadowVariables8()
func (c *InlineFunctionsTest) TestShadowVariables8() {
	// this should be inlined
	c.Test(
		`var a=0;
function foo(){return 3}
function _bar(){var a=foo()}
`,
		`var a=0;
function _bar(){var a=3}
`)
}

// Java: InlineFunctionsTest.testShadowVariables9()
func (c *InlineFunctionsTest) TestShadowVariables9() {
	// this should be inlined too [even if the global is not declared]
	c.Test(
		`function foo(){return 3}
function _bar(){var a=foo()}
`,
		"function _bar(){var a=3}")
}

// Java: InlineFunctionsTest.testShadowVariables10()
func (c *InlineFunctionsTest) TestShadowVariables10() {
	// callee var must be renamed.
	c.Test(
		`var a;function foo(){return a}
function _bar(){var a=foo()}
`,
		"var a;function _bar(){var a$jscomp$1=a}")
}

// Java: InlineFunctionsTest.testShadowVariables11()
func (c *InlineFunctionsTest) TestShadowVariables11() {
	// The call has a local variable
	// which collides with the function being inlined
	c.Test(
		`var a=0;var b=1;
function foo(){return a+a}
function _bar(){var a=foo();alert(a)}
`,
		`var a=0;var b=1;
function _bar(){var a$jscomp$1=a+a;
alert(a$jscomp$1)}
`)
}

// Java: InlineFunctionsTest.testShadowVariables12()
func (c *InlineFunctionsTest) TestShadowVariables12() {
	// 2 globals colliding
	c.Test(
		`var a=0;var b=1;
function foo(){return a+b}
function _bar(){var a=foo(),b;alert(a)}
`,
		`var a=0;var b=1;
function _bar(){var a$jscomp$1=a+b,
b$jscomp$1;
alert(a$jscomp$1)}
`)
}

// Java: InlineFunctionsTest.testShadowVariables13()
func (c *InlineFunctionsTest) TestShadowVariables13() {
	// The only change is to remove the collision
	c.Test(
		`var a=0;var b=1;
function foo(){return a+a}
function _bar(){var c=foo();alert(c)}
`,
		`var a=0;var b=1;
function _bar(){var c=a+a;alert(c)}
`)
}

// Java: InlineFunctionsTest.testShadowVariables14()
func (c *InlineFunctionsTest) TestShadowVariables14() {
	// There is a collision even though it is not read.
	c.Test(
		`var a=0;var b=1;
function foo(){return a+b}
function _bar(){var c=foo(),b;alert(c)}
`,
		`var a=0;var b=1;
function _bar(){var c=a+b,
b$jscomp$1;alert(c)}
`)
}

// Java: InlineFunctionsTest.testShadowVariables15()
func (c *InlineFunctionsTest) TestShadowVariables15() {
	// Both parent and child reference a global
	c.Test(
		`var a=0;var b=1;
function foo(){return a+a}
function _bar(){var c=foo();alert(c+a)}
`,
		`var a=0;var b=1;
function _bar(){var c=a+a;alert(c+a)}
`)
}

// Java: InlineFunctionsTest.testShadowVariables16()
func (c *InlineFunctionsTest) TestShadowVariables16() {
	c.assumeMinimumCapture = false
	// Inline functions defined as a child of the CALL node.
	c.Test(
		`var a=3;
function foo(){return a}
(function(){var a=5;(function(){foo()})()})()
`,
		`var a=3;
{var a$jscomp$inline_0=5;{a}}
`)
	c.assumeMinimumCapture = true
	// Inline functions defined as a child of the CALL node.
	c.Test(
		`var a=3;
function foo(){return a}
(function(){var a=5;(function(){foo()})()})()
`,
		`var a=3;
{var a$jscomp$inline_0=5;{a}}
`)
}

// Java: InlineFunctionsTest.testShadowVariables17()
func (c *InlineFunctionsTest) TestShadowVariables17() {
	c.Test(
		`var a=0;
function bar(){return a+a}
function foo(){return bar()}
function _goo(){var a=2;var x=foo();}
`,
		`var a=0;
function _goo(){var a$jscomp$1=2;var x=a+a}
`)
}

// Java: InlineFunctionsTest.testShadowVariables18()
func (c *InlineFunctionsTest) TestShadowVariables18() {
	c.Test(
		`var a=0;
function bar() { return a + a; }
function foo() { var a=3; return bar(); }
function _goo() { var a=2; var x = foo(); }
`,
		`var a=0;
function _goo(){
  var a$jscomp$2 = 2;
  var x;
  {
    var a$jscomp$inline_0 = 3;
    x = a + a;
  }
}
`)
}

// Java: InlineFunctionsTest.testShadowVariables19()
func (c *InlineFunctionsTest) TestShadowVariables19() {
	c.Test(
		`let a = 0;
function bar() { return a + a; }
function foo() { let a = 3; return bar(); }
function _goo() { let a = 2; let x = foo(); }
`,
		`let a = 0;
function _goo() {
  let a$jscomp$2 = 2;
  var JSCompiler_inline_result$jscomp$v0;
  {
     let a$jscomp$inline_0 = 3;
     JSCompiler_inline_result$jscomp$v0 = a + a;
  }
  let x = JSCompiler_inline_result$jscomp$v0;
}
`)
}

// Java: InlineFunctionsTest.testShadowVariables20()
func (c *InlineFunctionsTest) TestShadowVariables20() {
	c.Test(
		`const a = 0;
function bar() { return a + a; }
function foo() { const a = 3; return bar(); }
function _goo() { const a = 2; const x = foo(); }
`,
		`const a = 0;
function _goo() {
  const a$jscomp$2 = 2;
  var JSCompiler_inline_result$jscomp$v0;
  {
    const a$jscomp$inline_0 = 3;
    JSCompiler_inline_result$jscomp$v0 = a + a;
  }
  const x = JSCompiler_inline_result$jscomp$v0;
}
`)
}

// Java: InlineFunctionsTest.testCostBasedInlining1()
func (c *InlineFunctionsTest) TestCostBasedInlining1() {
	c.TestSame(
		`function foo(a){return a}
foo=new Function("return 1");
foo(1)
`)
}

// Java: InlineFunctionsTest.testCostBasedInlining2()
func (c *InlineFunctionsTest) TestCostBasedInlining2() {
	// Baseline complexity tests.
	// Single call, function not removed.
	c.Test(
		`function foo(a){return a}
var b=foo;
function _t1(){return foo(1)}
`,
		`function foo(a){return a}
var b=foo;
function _t1(){return 1}
`)
}

// Java: InlineFunctionsTest.testCostBasedInlining3()
func (c *InlineFunctionsTest) TestCostBasedInlining3() {
	// Two calls, function not removed.
	c.Test(
		`function foo(a,b){return a+b}
var b=foo;
function _t1(){return foo(1,2)}
function _t2(){return foo(2,3)}
`,
		`function foo(a,b){return a+b}
var b=foo;
function _t1(){return 1+2}
function _t2(){return 2+3}
`)
}

// Java: InlineFunctionsTest.testCostBasedInlining4()
func (c *InlineFunctionsTest) TestCostBasedInlining4() {
	// Two calls, function not removed.
	// Here there isn't enough savings to justify inlining.
	c.TestSame(
		`function foo(a,b){return a+b+a+b}
var b=foo;
function _t1(){return foo(1,2)}
function _t2(){return foo(2,3)}
`)
}

// Java: InlineFunctionsTest.testCostBasedInlining5()
func (c *InlineFunctionsTest) TestCostBasedInlining5() {
	// Here there is enough savings to justify inlining.
	c.Test(
		`function foo(a,b){return a+b+a+b}
function _t1(){return foo(1,2)}
function _t2(){return foo(2,3)}
`,
		`function _t1(){return 1+2+1+2}
function _t2(){return 2+3+2+3}
`)
}

// Java: InlineFunctionsTest.testCostBasedInlining6()
func (c *InlineFunctionsTest) TestCostBasedInlining6() {
	// Here we have a threshold test.
	// Do inline here:
	c.Test(
		`function foo(a,b){return a+b+a+b+a+b+a+b+4+5+6+7+8+9+1+2+3+4+5}
function _t1(){return foo(1,2)}
function _t2(){return foo(2,3)}
`,
		`function _t1(){return 1+2+1+2+1+2+1+2+4+5+6+7+8+9+1+2+3+4+5}
function _t2(){return 2+3+2+3+2+3+2+3+4+5+6+7+8+9+1+2+3+4+5}
`)
}

// Java: InlineFunctionsTest.testCostBasedInlining7()
func (c *InlineFunctionsTest) TestCostBasedInlining7() {
	// Don't inline here (not enough savings):
	c.TestSame(
		`function foo(a,b){
    return a+b+a+b+a+b+a+b+4+5+6+7+8+9+1+2+3+4+5+6}
function _t1(){return foo(1,2)}
function _t2(){return foo(2,3)}
`)
}

// Java: InlineFunctionsTest.testCostBasedInlining9()
func (c *InlineFunctionsTest) TestCostBasedInlining9() {
	// Here both direct and block inlining is used.  The call to f as a
	// parameter is inlined directly, which the call to f with f as a parameter
	// is inlined using block inlining.
	c.Test(
		`function f(a){return chg() + a + a;}
var a = f(f(1));
`,
		`var a;
{var a$jscomp$inline_0=chg()+1+1;
a=chg()+a$jscomp$inline_0+a$jscomp$inline_0}
`)
}

// Java: InlineFunctionsTest.testCostBasedInlining11()
func (c *InlineFunctionsTest) TestCostBasedInlining11() {
	// With block inlining
	c.Test(
		`function f(a){return chg() + a + a;}
var a = f(f(1))
`,
		`var a;
{var a$jscomp$inline_0=chg()+1+1;
a=chg()+a$jscomp$inline_0+a$jscomp$inline_0}
`)
}

// Java: InlineFunctionsTest.testCostBasedInlining12()
func (c *InlineFunctionsTest) TestCostBasedInlining12() {
	c.Test(
		`function f(a){return 1 + a + a;}
var a = f(1) + f(2);
`,
		"var a=1+1+1+(1+2+2)")
}

// Java: InlineFunctionsTest.testCostBasedInlineForSimpleFunction()
func (c *InlineFunctionsTest) TestCostBasedInlineForSimpleFunction() {
	calls := 100
	src := "function f(a){return a;}\n"
	for i := 0; i < calls; i++ {
		src += "f(chg());\n"
	}
	expected := ""
	for i := 0; i < calls; i++ {
		expected += "chg();\n"
	}
	c.Test(src, expected)
}

// Java: InlineFunctionsTest.testCostBasedInliningComplex1()
func (c *InlineFunctionsTest) TestCostBasedInliningComplex1() {
	c.TestSame(
		`function foo(a){a()}
foo=new Function("return 1");
foo(1)
`)
}

// Java: InlineFunctionsTest.testCostBasedInliningComplex2()
func (c *InlineFunctionsTest) TestCostBasedInliningComplex2() {
	// Baseline complexity tests.
	// Single call, function not removed.
	c.Test(
		`function foo(a){a()}
var b=foo;
function _t1(){foo(x)}
`,
		`function foo(a){a()}
var b=foo;
function _t1(){{x()}}
`)
}

// Java: InlineFunctionsTest.testCostBasedInliningComplex3()
func (c *InlineFunctionsTest) TestCostBasedInliningComplex3() {
	// Two calls, function not removed.
	c.Test(
		`function foo(a,b){a+b}
var b=foo;
function _t1(){foo(1,2)}
function _t2(){foo(2,3)}
`,
		`function foo(a,b){a+b}
var b=foo;
function _t1(){{1+2}}
function _t2(){{2+3}}
`)
}

// Java: InlineFunctionsTest.testCostBasedInliningComplex4()
func (c *InlineFunctionsTest) TestCostBasedInliningComplex4() {
	// Two calls, function not removed.
	// Here there isn't enough savings to justify inlining.
	c.TestSame(
		`function foo(a,b){a+b+a+b}
var b=foo;
function _t1(){foo(1,2)}
function _t2(){foo(2,3)}
`)
}

// Java: InlineFunctionsTest.testCostBasedInliningComplex5()
func (c *InlineFunctionsTest) TestCostBasedInliningComplex5() {
	// Here there is enough savings to justify inlining.
	c.Test(
		`function foo(a,b){a+b+a+b}
function _t1(){foo(1,2)}
function _t2(){foo(2,3)}
`,
		`function _t1(){{1+2+1+2}}
function _t2(){{2+3+2+3}}
`)
}

// Java: InlineFunctionsTest.testCostBasedInliningComplex6()
func (c *InlineFunctionsTest) TestCostBasedInliningComplex6() {
	// Here we have a threshold test.
	// Do inline here:
	c.Test(
		`function foo(a,b){a+b+a+b+a+b+a+b+4+5+6+7+8+9+1}
function _t1(){foo(1,2)}
function _t2(){foo(2,3)}
`,
		`function _t1(){{1+2+1+2+1+2+1+2+4+5+6+7+8+9+1}}
function _t2(){{2+3+2+3+2+3+2+3+4+5+6+7+8+9+1}}
`)
}

// Java: InlineFunctionsTest.testCostBasedInliningComplex7()
func (c *InlineFunctionsTest) TestCostBasedInliningComplex7() {
	// Don't inline here (not enough savings):
	c.TestSame(
		`function foo(a,b){a+b+a+b+a+b+a+b+4+5+6+7+8+9+1+2}
function _t1(){foo(1,2)}
function _t2(){foo(2,3)}
`)
}

// Java: InlineFunctionsTest.testCostBasedInliningComplex8()
func (c *InlineFunctionsTest) TestCostBasedInliningComplex8() {
	// Verify multiple references in the same statement.
	c.TestSame(
		`function _f(a){1+a+a}
a=_f(1)+_f(1)
`)
}

// Java: InlineFunctionsTest.testCostBasedInliningComplex9()
func (c *InlineFunctionsTest) TestCostBasedInliningComplex9() {
	c.Test(
		`function f(a){1 + a + a;}
f(1);f(2);
`,
		`{1+1+1}
{1+2+2}
`)
}

// Java: InlineFunctionsTest.testDoubleInlining2()
func (c *InlineFunctionsTest) TestDoubleInlining2() {
	c.Test(
		`var foo = function(a) { return getWindow(a); };
var bar = function(b) { return b; };
foo(bar(x));
`,
		"getWindow(x)")
}

// Java: InlineFunctionsTest.testNoInlineOfNonGlobalFunction1()
func (c *InlineFunctionsTest) TestNoInlineOfNonGlobalFunction1() {
	c.Test(
		`var g;function _f(){function g(){return 0}}
function _h(){return g()}
`,
		`var g;function _f(){}
function _h(){return g()}
`)
}

// Java: InlineFunctionsTest.testNoInlineOfNonGlobalFunction2()
func (c *InlineFunctionsTest) TestNoInlineOfNonGlobalFunction2() {
	c.Test(
		`var g;
function _f() { var g = function() { return 0; }; }
function _h() { return g(); }
`,
		`var g;
function _f(){}
function _h(){ return g(); }
`)
}

// Java: InlineFunctionsTest.testNoInlineOfNonGlobalFunction3()
func (c *InlineFunctionsTest) TestNoInlineOfNonGlobalFunction3() {
	c.Test(
		`var g;function _f(){var g=function(){return 0}}
function _h(){return g()}
`,
		`var g;function _f(){}
function _h(){return g()}
`)
}

// Java: InlineFunctionsTest.testNoInlineOfNonGlobalFunction4()
func (c *InlineFunctionsTest) TestNoInlineOfNonGlobalFunction4() {
	c.Test(
		`var g;function _f(){function g(){return 0}}
function _h(){return g()}
`,
		`var g;function _f(){}
function _h(){return g()}
`)
}

// Java: InlineFunctionsTest.testNoInlineMaskedFunction()
func (c *InlineFunctionsTest) TestNoInlineMaskedFunction() {
	// Normalization makes this test of marginal value.
	// The unreferenced function is removed.
	c.Test(
		`var g=function(){return 0};
function _f(g){return g()}
`,
		"function _f(g$jscomp$1){return g$jscomp$1()}")
}

// Java: InlineFunctionsTest.testNoInlineNonFunction()
func (c *InlineFunctionsTest) TestNoInlineNonFunction() {
	c.TestSame("var g=3;function _f(){return g()}")
}

// Java: InlineFunctionsTest.testInlineCall()
func (c *InlineFunctionsTest) TestInlineCall() {
	c.Test("function f(g) { return g.h(); } f('x');", "\"x\".h()")
}

// Java: InlineFunctionsTest.testInlineFunctionWithArgsMismatch1()
func (c *InlineFunctionsTest) TestInlineFunctionWithArgsMismatch1() {
	c.Test("function f(g) { return g; } f();", "void 0")
}

// Java: InlineFunctionsTest.testInlineFunctionWithArgsMismatch2()
func (c *InlineFunctionsTest) TestInlineFunctionWithArgsMismatch2() {
	c.Test("function f() { return 0; } f(1);", "0")
}

// Java: InlineFunctionsTest.testInlineFunctionWithArgsMismatch3()
func (c *InlineFunctionsTest) TestInlineFunctionWithArgsMismatch3() {
	c.Test("function f(one, two, three) { return one + two + three; } f(1);", "1+void 0+void 0")
}

// Java: InlineFunctionsTest.testInlineFunctionWithArgsMismatch4()
func (c *InlineFunctionsTest) TestInlineFunctionWithArgsMismatch4() {
	c.Test(
		`function f(one, two, three) { return one + two + three; }
f(1,2,3,4,5);
`,
		"1+2+3")
}

// Java: InlineFunctionsTest.testComplexInlineNoResultNoParamCall1()
func (c *InlineFunctionsTest) TestComplexInlineNoResultNoParamCall1() {
	c.Test("function f(){a()}f()", "{a()}")
}

// Java: InlineFunctionsTest.testComplexInlineNoResultNoParamCall2()
func (c *InlineFunctionsTest) TestComplexInlineNoResultNoParamCall2() {
	c.Test(
		"function f() { if (true) { return; } else; } f();",
		"{JSCompiler_inline_label_f_v0:{ if(true)break JSCompiler_inline_label_f_v0;else; } }")
}

// Java: InlineFunctionsTest.testComplexInlineNoResultNoParamCall3()
func (c *InlineFunctionsTest) TestComplexInlineNoResultNoParamCall3() {
	// We now allow vars in the global space.
	//   Don't inline into vars into global scope.
	//   testSame("function f(){a();b();var z=1+1}f()");
	// But do inline into functions
	c.Test(
		"function f(){a();b();var z=1+1}function _foo(){f()}",
		"function _foo(){{a();b();var z$jscomp$inline_0=1+1}}")
}

// Java: InlineFunctionsTest.testComplexInlineNoResultWithParamCall1()
func (c *InlineFunctionsTest) TestComplexInlineNoResultWithParamCall1() {
	c.Test("function f(x){a(x)}f(1)", "{a(1)}")
}

// Java: InlineFunctionsTest.testComplexInlineNoResultWithParamCall2()
func (c *InlineFunctionsTest) TestComplexInlineNoResultWithParamCall2() {
	c.Test("function f(x,y){a(x)}var b=1;f(1,b)", "var b=1;{a(1)}")
}

// Java: InlineFunctionsTest.testComplexInlineNoResultWithParamCall3()
func (c *InlineFunctionsTest) TestComplexInlineNoResultWithParamCall3() {
	c.Test("function f(x,y){if (x) y(); return true;}var b=1;f(1,b)", "var b=1;{if(1)b();true}")
}

// Java: InlineFunctionsTest.testComplexInline1()
func (c *InlineFunctionsTest) TestComplexInline1() {
	c.Test(
		"function f(){if (true){return;}else;} z=f();",
		`{
  JSCompiler_inline_label_f_v0: {
    if (true) {
      z = void 0;
      break JSCompiler_inline_label_f_v0;
    }else;
    z = void 0;
  }
}
`)
}

// Java: InlineFunctionsTest.testComplexInline2()
func (c *InlineFunctionsTest) TestComplexInline2() {
	c.Test(
		"function f(){if (true){return;}else return;} z=f();",
		`{
  JSCompiler_inline_label_f_v0: {
    if(true) {
      z = void 0;
      break JSCompiler_inline_label_f_v0;
    } else {
      z=void 0;
      break JSCompiler_inline_label_f_v0;
    }
    z=void 0;
  }
}
`)
}

// Java: InlineFunctionsTest.testComplexInline3()
func (c *InlineFunctionsTest) TestComplexInline3() {
	c.Test(
		"function f(){if (true){return 1;}else return 0;} z=f();",
		`{
JSCompiler_inline_label_f_v0:{
if(true){z=1;
break JSCompiler_inline_label_f_v0;
} else {
z=0;
break JSCompiler_inline_label_f_v0;}
z = void 0
}
}
`)
}

// Java: InlineFunctionsTest.testComplexInline4()
func (c *InlineFunctionsTest) TestComplexInline4() {
	c.Test("function f(x){a(x)} z = f(1)", "{a(1);z=void 0}")
}

// Java: InlineFunctionsTest.testComplexInline5()
func (c *InlineFunctionsTest) TestComplexInline5() {
	c.Test("function f(x,y){a(x)}var b=1;z=f(1,b)", "var b=1;{a(1);z=void 0}")
}

// Java: InlineFunctionsTest.testComplexInline6()
func (c *InlineFunctionsTest) TestComplexInline6() {
	c.Test("function f(x,y){if (x) y(); return true;}var b=1;z=f(1,b)", "var b=1;{if(1)b();z=true}")
}

// Java: InlineFunctionsTest.testComplexInline7()
func (c *InlineFunctionsTest) TestComplexInline7() {
	c.Test(
		"function f(x,y){if (x) return y(); else return true;} var b=1;z=f(1,b)",
		`var b=1;
{
  JSCompiler_inline_label_f_v0: {
    if(1) {
      z = b();
      break JSCompiler_inline_label_f_v0;
    } else {
      z = true;
      break JSCompiler_inline_label_f_v0
    }
    z = void 0;
  }
}
`)
}

// Java: InlineFunctionsTest.testComplexInline8()
func (c *InlineFunctionsTest) TestComplexInline8() {
	c.Test("function f(x){a(x)}var z=f(1)", "var z;{a(1);z=void 0}")
}

// Java: InlineFunctionsTest.testInlineIntoLoopWithUninitializedVars1()
func (c *InlineFunctionsTest) TestInlineIntoLoopWithUninitializedVars1() {
	c.Test(
		`function f(){var x;if (true){x=7} return x;}
for(;;) {f();}
`,
		`for(;;) {
  {
    var x$jscomp$inline_0 = void 0;
    if (true) { x$jscomp$inline_0 = 7; }
    x$jscomp$inline_0;
  }
}
`)
}

// Java: InlineFunctionsTest.testInlineIntoLoopWithUninitializedVars2()
func (c *InlineFunctionsTest) TestInlineIntoLoopWithUninitializedVars2() {
	c.Test(
		`function f(){x = x||1; var x; return x;}
for(;;) {f();}
`,
		`for(;;) {
  {
    var x$jscomp$inline_0 = void 0;
    x$jscomp$inline_0 = x$jscomp$inline_0 || 1;
    x$jscomp$inline_0;
  }
}
`)
}

// Java: InlineFunctionsTest.testComplexInlineVars1()
func (c *InlineFunctionsTest) TestComplexInlineVars1() {
	c.Test(
		"function f() { if (true) { return; } else; } var z = f();",
		`var z;
{
  JSCompiler_inline_label_f_v0:{
    if(true) {
      z = void 0;
      break JSCompiler_inline_label_f_v0;
    } else;
    z = void 0;
  }
}
`)
}

// Java: InlineFunctionsTest.testComplexInlineVars2()
func (c *InlineFunctionsTest) TestComplexInlineVars2() {
	c.Test(
		"function f(){if (true){return;}else return;}var z=f();",
		`var z;
{
  JSCompiler_inline_label_f_v0:{
    if (true) {
      z = void 0;
      break JSCompiler_inline_label_f_v0;
    } else {
      z = void 0;
      break JSCompiler_inline_label_f_v0;
    }
    z=void 0;
  }
}
`)
}

// Java: InlineFunctionsTest.testComplexInlineVars3()
func (c *InlineFunctionsTest) TestComplexInlineVars3() {
	c.Test(
		"function f(){if (true){return 1;}else return 0;}var z=f();",
		`var z;
{
  JSCompiler_inline_label_f_v0:{
    if (true) {
      z = 1;
      break JSCompiler_inline_label_f_v0;
    } else {
      z = 0;
      break JSCompiler_inline_label_f_v0;
    }
    z = void 0;
  }
}
`)
}

// Java: InlineFunctionsTest.testComplexInlineVars4()
func (c *InlineFunctionsTest) TestComplexInlineVars4() {
	c.Test("function f(x){a(x)}var z = f(1)", "var z;{a(1);z=void 0}")
}

// Java: InlineFunctionsTest.testComplexInlineVars5()
func (c *InlineFunctionsTest) TestComplexInlineVars5() {
	c.Test("function f(x,y){a(x)}var b=1;var z=f(1,b)", "var b=1;var z;{a(1);z=void 0}")
}

// Java: InlineFunctionsTest.testComplexInlineVars6()
func (c *InlineFunctionsTest) TestComplexInlineVars6() {
	c.Test(
		"function f(x,y){if (x) y(); return true;}var b=1;var z=f(1,b)",
		"var b=1;var z;{if(1)b();z=true}")
}

// Java: InlineFunctionsTest.testComplexInlineVars7()
func (c *InlineFunctionsTest) TestComplexInlineVars7() {
	c.Test(
		`function f(x,y){if (x) return y(); else return true;}
var b=1;var z=f(1,b)
`,
		`var b=1;
var z;
{
  JSCompiler_inline_label_f_v0:{
    if (1) {
      z = b();
      break JSCompiler_inline_label_f_v0;
    } else {
      z = true;
      break JSCompiler_inline_label_f_v0;
    }
    z = void 0;
  }
}
`)
}

// Java: InlineFunctionsTest.testComplexInlineVars8()
func (c *InlineFunctionsTest) TestComplexInlineVars8() {
	c.Test("function f(x){a(x)}var x;var z=f(1)", "var x;var z;{a(1);z=void 0}")
}

// Java: InlineFunctionsTest.testComplexInlineVars9()
func (c *InlineFunctionsTest) TestComplexInlineVars9() {
	c.Test("function f(x){a(x)}var x;var z=f(1);var y", "var x;var z;{a(1);z=void 0}var y")
}

// Java: InlineFunctionsTest.testComplexInlineVars10()
func (c *InlineFunctionsTest) TestComplexInlineVars10() {
	c.Test(
		"function f(x){a(x)}var x=blah();var z=f(1);var y=blah();",
		"var x=blah();var z;{a(1);z=void 0}var y=blah()")
}

// Java: InlineFunctionsTest.testComplexInlineVars11()
func (c *InlineFunctionsTest) TestComplexInlineVars11() {
	c.Test(
		"function f(x){a(x)}var x=blah();var z=f(1);var y;",
		"var x=blah();var z;{a(1);z=void 0}var y")
}

// Java: InlineFunctionsTest.testComplexInlineVars12()
func (c *InlineFunctionsTest) TestComplexInlineVars12() {
	c.Test(
		"function f(x){a(x)}var x;var z=f(1);var y=blah();",
		"var x;var z;{a(1);z=void 0}var y=blah()")
}

// Java: InlineFunctionsTest.testComplexInlineInExpressionss1()
func (c *InlineFunctionsTest) TestComplexInlineInExpressionss1() {
	c.Test("function f(){a()}var z=f()", "var z;{a();z=void 0}")
}

// Java: InlineFunctionsTest.testComplexInlineInExpressionss2()
func (c *InlineFunctionsTest) TestComplexInlineInExpressionss2() {
	c.Test(
		"function f(){a()}c=z=f()",
		`var JSCompiler_inline_result$jscomp$v0;
{a();JSCompiler_inline_result$jscomp$v0=void 0;}
c=z=JSCompiler_inline_result$jscomp$v0
`)
}

// Java: InlineFunctionsTest.testComplexInlineInExpressionss3()
func (c *InlineFunctionsTest) TestComplexInlineInExpressionss3() {
	c.Test(
		"function f(){a()}c=z=f()",
		`var JSCompiler_inline_result$jscomp$v0;
{a();JSCompiler_inline_result$jscomp$v0=void 0;}
c=z=JSCompiler_inline_result$jscomp$v0
`)
}

// Java: InlineFunctionsTest.testComplexInlineInExpressionss4()
func (c *InlineFunctionsTest) TestComplexInlineInExpressionss4() {
	c.Test(
		"function f(){a()}if(z=f());",
		`var JSCompiler_inline_result$jscomp$v0;
{a();JSCompiler_inline_result$jscomp$v0=void 0;}
if(z=JSCompiler_inline_result$jscomp$v0);
`)
}

// Java: InlineFunctionsTest.testComplexInlineInExpressionss5()
func (c *InlineFunctionsTest) TestComplexInlineInExpressionss5() {
	c.Test(
		"function f(){a()}if(z.y=f());",
		`var JSCompiler_temp_const$jscomp$v0=z;
var JSCompiler_inline_result$jscomp$v1;
{a();JSCompiler_inline_result$jscomp$v1=void 0;}
if(JSCompiler_temp_const$jscomp$v0.y=JSCompiler_inline_result$jscomp$v1);
`)
}

// Java: InlineFunctionsTest.testComplexNoInline1()
func (c *InlineFunctionsTest) TestComplexNoInline1() {
	c.TestSame("function f(){a()}while(z=f())continue")
}

// Java: InlineFunctionsTest.testComplexNoInline2()
func (c *InlineFunctionsTest) TestComplexNoInline2() {
	c.TestSame("function f(){a()}do;while(z=f())")
}

// Java: InlineFunctionsTest.testComplexSample()
func (c *InlineFunctionsTest) TestComplexSample() {
	c.Test(
		`var foo = function(stylesString, opt_element) {
  var styleSheet = null;
  if (goog$userAgent$IE)
    styleSheet = 0;
  else
    var head = 0;

  goo$zoo(styleSheet, stylesString);
  return styleSheet;
};

var goo$zoo = function(element, stylesString) {
  if (goog$userAgent$IE)
    element.cssText = stylesString;
  else {
    var propToSet = 'innerText';
    element[propToSet] = stylesString;
  }
};
(function(){foo(a,b);})();
`,
		`{
  {
    var styleSheet$jscomp$inline_2 = null;
    if (goog$userAgent$IE)
      styleSheet$jscomp$inline_2 = 0;
    else
      var head$jscomp$inline_3 = 0;
    {
       var element$jscomp$inline_4 = styleSheet$jscomp$inline_2;
       var stylesString$jscomp$inline_5=a;
       if (goog$userAgent$IE)
         element$jscomp$inline_4.cssText = stylesString$jscomp$inline_5;
       else {
         var propToSet$jscomp$inline_6 = 'innerText';
         element$jscomp$inline_4[propToSet$jscomp$inline_6] =
             stylesString$jscomp$inline_5
       }
    }
    styleSheet$jscomp$inline_2;
  }
}
`)
}

// Java: InlineFunctionsTest.testComplexSampleNoInline()
func (c *InlineFunctionsTest) TestComplexSampleNoInline() {
	c.TestSame(
		`foo = function(stylesString, opt_element) {
  var styleSheet = null;
  if(goog$userAgent$IE)
    styleSheet=0;
  else
    var head=0;

  goo$zoo(styleSheet, stylesString);
  return styleSheet
};

goo$zoo = function(element, stylesString) {
  if (goog$userAgent$IE)
    element.cssText = stylesString;
  else{
    var propToSet = goog$userAgent$WEBKIT? 'innerText' : 'innerHTML';
    element[propToSet] = stylesString
  }
}
`)
}

// Test redefinition of parameter name.
// Java: InlineFunctionsTest.testComplexNoVarSub()
func (c *InlineFunctionsTest) TestComplexNoVarSub() {
	c.Test(
		`function foo(x){
var x;
y=x
}
foo(1)
`,
		"{y=1}")
}

// Java: InlineFunctionsTest.testComplexFunctionWithFunctionDefinition1()
func (c *InlineFunctionsTest) TestComplexFunctionWithFunctionDefinition1() {
	c.Test("function f(){call(function(){return})}f()", "{call(function(){return})}")
}

// Java: InlineFunctionsTest.testComplexFunctionWithFunctionDefinition2()
func (c *InlineFunctionsTest) TestComplexFunctionWithFunctionDefinition2() {
	c.assumeMinimumCapture = false
	// Don't inline if local names might be captured.
	c.TestSame(
		`function f(a) {
  call(function(){return});
}
f();
`)
	c.assumeMinimumCapture = true
	c.Test(
		`(function() {
  var f = function(a) { call(function(){return a}); };
  f();
})();
`,
		`{
  {
    var a$jscomp$inline_2 = void 0;
    call(function() { return a$jscomp$inline_2; });
  }
}
`)
}

// Java: InlineFunctionsTest.testComplexFunctionWithFunctionDefinition2a()
func (c *InlineFunctionsTest) TestComplexFunctionWithFunctionDefinition2a() {
	c.assumeMinimumCapture = false
	// Don't inline if local names might be captured.
	c.TestSame(
		`(function() {
  var f = function(a) { call(function(){return a}); };
  f();
})()
`)
	c.assumeMinimumCapture = true
	c.Test(
		`(function() {
  var f = function(a) { call(function(){return a}); };
  f();
})()
`,
		`{
  {
    var a$jscomp$inline_2 = void 0;
    call(function() { return a$jscomp$inline_2; });
  }
}
`)
}

// Java: InlineFunctionsTest.testComplexFunctionWithFunctionDefinition3()
func (c *InlineFunctionsTest) TestComplexFunctionWithFunctionDefinition3() {
	c.assumeMinimumCapture = false
	// Don't inline if local names might need to be captured.
	c.TestSame("function f(){var a; call(function(){return a})}f()")
	c.assumeMinimumCapture = true
	c.Test(
		"function f(){var a; call(function(){return a})}f()",
		"{var a$jscomp$inline_0;call(function(){return a$jscomp$inline_0})}")
}

// Java: InlineFunctionsTest.testDecomposePlusEquals()
func (c *InlineFunctionsTest) TestDecomposePlusEquals() {
	c.Test(
		"function f(){a=1;return 1} var x = 1; x += f()",
		`var x = 1;
var JSCompiler_temp_const$jscomp$v0 = x;
var JSCompiler_inline_result$jscomp$v1;
{a=1;
 JSCompiler_inline_result$jscomp$v1=1}
x = JSCompiler_temp_const$jscomp$v0 + JSCompiler_inline_result$jscomp$v1;
`)
}

// Java: InlineFunctionsTest.testDecomposeFunctionExpressionInCall()
func (c *InlineFunctionsTest) TestDecomposeFunctionExpressionInCall() {
	c.Test(
		`(function(map){descriptions_=map})(
function(){
var ret={};
ret[ONE]='a';
ret[TWO]='b';
return ret
}()
);
`,
		`var JSCompiler_inline_result$jscomp$v0;
{
var ret$jscomp$inline_0={};
ret$jscomp$inline_0[ONE]='a';
ret$jscomp$inline_0[TWO]='b';
JSCompiler_inline_result$jscomp$v0 = ret$jscomp$inline_0;
}
{
descriptions_ = JSCompiler_inline_result$jscomp$v0;
}
`)
}

// Java: InlineFunctionsTest.testInlineConstructor1()
func (c *InlineFunctionsTest) TestInlineConstructor1() {
	c.Test("function f() {} function _g() {f.call(this)}", "function _g() {void 0}")
}

// Java: InlineFunctionsTest.testInlineConstructor2()
func (c *InlineFunctionsTest) TestInlineConstructor2() {
	c.Test(
		"function f() {} f.prototype.a = 0; function _g() {f.call(this)}",
		"function f() {} f.prototype.a = 0; function _g() {void 0}")
}

// Java: InlineFunctionsTest.testInlineConstructor3()
func (c *InlineFunctionsTest) TestInlineConstructor3() {
	c.Test(
		`function f() {x.call(this)} f.prototype.a = 0;
function _g() {f.call(this)}
`,
		`function f() {x.call(this)} f.prototype.a = 0;
function _g() {{x.call(this)}}
`)
}

// Java: InlineFunctionsTest.testInlineConstructor4()
func (c *InlineFunctionsTest) TestInlineConstructor4() {
	c.Test(
		`function f() {x.call(this)} f.prototype.a = 0;
function _g() {var t = f.call(this)}
`,
		`function f() {x.call(this)} f.prototype.a = 0;
    function _g() {var t; {x.call(this); t = void 0}}
`)
}

// Java: InlineFunctionsTest.testFunctionExpressionInlining1()
func (c *InlineFunctionsTest) TestFunctionExpressionInlining1() {
	c.Test("(function(){})()", "void 0")
}

// Java: InlineFunctionsTest.testFunctionExpressionInlining2()
func (c *InlineFunctionsTest) TestFunctionExpressionInlining2() {
	c.Test("(function(){foo()})()", "{foo()}")
}

// Java: InlineFunctionsTest.testFunctionExpressionInlining3()
func (c *InlineFunctionsTest) TestFunctionExpressionInlining3() {
	c.Test("var a = (function(){return foo()})()", "var a = foo()")
}

// Java: InlineFunctionsTest.testFunctionExpressionInlining4()
func (c *InlineFunctionsTest) TestFunctionExpressionInlining4() {
	c.Test("var a; a = 1 + (function(){return foo()})()", "var a; a = 1 + foo()")
}

// Java: InlineFunctionsTest.testFunctionExpressionCallInlining1()
func (c *InlineFunctionsTest) TestFunctionExpressionCallInlining1() {
	c.Test("(function(){}).call(this)", "void 0")
}

// Java: InlineFunctionsTest.testFunctionExpressionCallInlining2()
func (c *InlineFunctionsTest) TestFunctionExpressionCallInlining2() {
	c.Test("(function(){foo(this)}).call(this)", "{foo(this)}")
}

// Java: InlineFunctionsTest.testFunctionExpressionCallInlining3()
func (c *InlineFunctionsTest) TestFunctionExpressionCallInlining3() {
	c.Test("var a = (function(){return foo(this)}).call(this)", "var a = foo(this)")
}

// Java: InlineFunctionsTest.testFunctionExpressionCallInlining4()
func (c *InlineFunctionsTest) TestFunctionExpressionCallInlining4() {
	c.Test("var a; a = 1 + (function(){return foo(this)}).call(this)", "var a; a = 1 + foo(this)")
}

// Java: InlineFunctionsTest.testFunctionExpressionCallInlining5()
func (c *InlineFunctionsTest) TestFunctionExpressionCallInlining5() {
	c.Test("a:(function(){return foo()})()", "a:foo()")
}

// Java: InlineFunctionsTest.testFunctionExpressionCallInlining6()
func (c *InlineFunctionsTest) TestFunctionExpressionCallInlining6() {
	c.Test("a:(function(){return foo()}).call(this)", "a:foo()")
}

// Java: InlineFunctionsTest.testFunctionExpressionCallInlining7()
func (c *InlineFunctionsTest) TestFunctionExpressionCallInlining7() {
	c.Test("a:(function(){})()", "a:void 0")
}

// Java: InlineFunctionsTest.testFunctionExpressionCallInlining8()
func (c *InlineFunctionsTest) TestFunctionExpressionCallInlining8() {
	c.Test("a:(function(){}).call(this)", "a:void 0")
}

// Java: InlineFunctionsTest.testFunctionExpressionCallInlining9()
func (c *InlineFunctionsTest) TestFunctionExpressionCallInlining9() {
	// ... with unused recursive name.
	c.Test("(function foo(){})()", "void 0")
}

// Java: InlineFunctionsTest.testFunctionExpressionCallInlining10()
func (c *InlineFunctionsTest) TestFunctionExpressionCallInlining10() {
	// ... with unused recursive name.
	c.Test("(function foo(){}).call(this)", "void 0")
}

// Java: InlineFunctionsTest.testFunctionExpressionCallInlining11a()
func (c *InlineFunctionsTest) TestFunctionExpressionCallInlining11a() {
	// Inline functions that return inner functions.
	c.Test("((function(){return function(){foo()}})())();", "{foo()}")
}

// Java: InlineFunctionsTest.testFunctionExpressionCallInlining11b()
func (c *InlineFunctionsTest) TestFunctionExpressionCallInlining11b() {
	c.assumeMinimumCapture = false
	// Can't inline functions that return inner functions and have local names.
	c.TestSame("((function(){var a; return function(){foo()}})())();")
	c.assumeMinimumCapture = true
	c.Test(
		"((function(){var a; return function(){foo()}})())();",
		`var JSCompiler_inline_result$jscomp$v0;
{var a$jscomp$inline_0;
JSCompiler_inline_result$jscomp$v0=function(){foo()};}
JSCompiler_inline_result$jscomp$v0()
`)
}

// Java: InlineFunctionsTest.testFunctionExpressionCallInlining11c()
func (c *InlineFunctionsTest) TestFunctionExpressionCallInlining11c() {
	// TODO(johnlenz): Can inline, not temps needed.
	c.assumeMinimumCapture = false
	c.TestSame(
		`function _x() {
  ((function(){return function(){foo()}})())();
}
`)
	c.assumeMinimumCapture = true
	c.Test(
		`function _x() {
  ((function(){return function(){foo()}})())();
}
`,
		`function _x() {
  {foo()}
}
`)
}

// Java: InlineFunctionsTest.testFunctionExpressionCallInlining11d()
func (c *InlineFunctionsTest) TestFunctionExpressionCallInlining11d() {
	// TODO(johnlenz): Can inline into a function containing eval, if
	// no names are introduced.
	c.assumeMinimumCapture = false
	c.TestSame(
		`function _x() {
  eval();
  ((function(){return function(){foo()}})())();
}
`)
	c.assumeMinimumCapture = true
	c.Test(
		`function _x() {
  eval();
  ((function(){return function(){foo()}})())();
}
`,
		`function _x() {
  eval();
  {foo()}
}
`)
}

// Java: InlineFunctionsTest.testFunctionExpressionCallInlining11e()
func (c *InlineFunctionsTest) TestFunctionExpressionCallInlining11e() {
	// No, don't inline into a function containing eval,
	// if temps are introduced.
	c.assumeMinimumCapture = false
	c.TestSame(
		`function _x() {
  eval();
  ((function(a){return function(){foo()}})())();
}
`)
	c.assumeMinimumCapture = true
	c.Test(
		`function _x() {
  eval();
  ((function(a){return function(){foo()}})())();
}
`,
		`function _x() {
  eval();
  {foo();}
}
`)
}

// Java: InlineFunctionsTest.testFunctionExpressionCallInlining12()
func (c *InlineFunctionsTest) TestFunctionExpressionCallInlining12() {
	// Can't inline functions that recurse.
	c.TestSame("(function foo(){foo()})()")
}

// Java: InlineFunctionsTest.testFunctionExpressionOmega()
func (c *InlineFunctionsTest) TestFunctionExpressionOmega() {
	// ... with unused recursive name.
	c.Test(
		"(function (f){f(f)})(function(f){f(f)})",
		`{
  var f$jscomp$inline_0 = function(f$jscomp$1) {
    f$jscomp$1(f$jscomp$1);
  };
  {
    {
       f$jscomp$inline_0(f$jscomp$inline_0);
    }
  }
}
`)
}

// Java: InlineFunctionsTest.testLocalFunctionInlining1()
func (c *InlineFunctionsTest) TestLocalFunctionInlining1() {
	c.Test("function _f(){ function g() {} g() }", "function _f(){ void 0 }")
}

// Java: InlineFunctionsTest.testLocalFunctionInlining2()
func (c *InlineFunctionsTest) TestLocalFunctionInlining2() {
	c.Test("function _f(){ function g() {foo(); bar();} g() }", "function _f(){ {foo(); bar();} }")
}

// Java: InlineFunctionsTest.testLocalFunctionInlining3()
func (c *InlineFunctionsTest) TestLocalFunctionInlining3() {
	c.Test("function _f(){ function g() {foo(); bar();} g() }", "function _f(){ {foo(); bar();} }")
}

// Java: InlineFunctionsTest.testLocalFunctionInlining4()
func (c *InlineFunctionsTest) TestLocalFunctionInlining4() {
	c.Test("function _f(){ function g() {return 1} return g() }", "function _f(){ return 1 }")
}

// Java: InlineFunctionsTest.testLocalFunctionInlining4_optChain()
func (c *InlineFunctionsTest) TestLocalFunctionInlining4_optChain() {
	c.Test("function _f(){ function g() {return 1} return g?.() }", "function _f(){ return 1 }")
}

// Java: InlineFunctionsTest.testLocalFunctionInlining5()
func (c *InlineFunctionsTest) TestLocalFunctionInlining5() {
	c.TestSame("function _f(){ function g() {this;} g() }")
}

// Java: InlineFunctionsTest.testLocalFunctionInlining6()
func (c *InlineFunctionsTest) TestLocalFunctionInlining6() {
	c.TestSame("function _f(){ function g() {this;} return g; }")
}

// Java: InlineFunctionsTest.testLocalFunctionInliningOnly1()
func (c *InlineFunctionsTest) TestLocalFunctionInliningOnly1() {
	c.inliningReach = CompilerOptions_Reach_ALL
	c.Test("function f(){} f()", "void 0;")
	c.inliningReach = CompilerOptions_Reach_LOCAL_ONLY
	c.TestSame("function f(){} f()")
}

// Java: InlineFunctionsTest.testLocalFunctionInliningOnly2()
func (c *InlineFunctionsTest) TestLocalFunctionInliningOnly2() {
	c.inliningReach = CompilerOptions_Reach_LOCAL_ONLY
	c.TestSame("function f(){} f()")
	c.Test(
		"function f(){ function g() {return 1} return g() }; f();",
		"function f(){ return 1 }; f();")
}

// Java: InlineFunctionsTest.testLocalFunctionInliningOnly3()
func (c *InlineFunctionsTest) TestLocalFunctionInliningOnly3() {
	c.inliningReach = CompilerOptions_Reach_LOCAL_ONLY
	c.TestSame("function f(){} f()")
	c.Test("(function(){ function g() {return 1} return g() })();", "(function(){ return 1 })();")
}

// Java: InlineFunctionsTest.testLocalFunctionInliningOnly4()
func (c *InlineFunctionsTest) TestLocalFunctionInliningOnly4() {
	c.inliningReach = CompilerOptions_Reach_LOCAL_ONLY
	c.TestSame("function f(){} f()")
	c.Test("(function(){ return (function() {return 1})() })();", "(function(){ return 1 })();")
}

// Java: InlineFunctionsTest.testInlineWithThis1()
func (c *InlineFunctionsTest) TestInlineWithThis1() {
	c.assumeStrictThis = false
	// If no "this" is provided it might need to be coerced to the global
	// "this".
	c.TestSame("function f(){} f.call();")
	c.TestSame("function f(){this} f.call();")
	c.assumeStrictThis = true
	// In strict mode, "this" is never coerced so we can use the provided value.
	c.Test("function f(){} f.call();", "{}")
	c.Test("function f(){this} f.call();", "{void 0;}")
}

// Java: InlineFunctionsTest.testInlineWithThis2()
func (c *InlineFunctionsTest) TestInlineWithThis2() {
	// "this" can always be replaced with "this"
	c.assumeStrictThis = false
	c.Test("function f(){} f.call(this);", "void 0")
	c.assumeStrictThis = true
	c.Test("function f(){} f.call(this);", "void 0")
}

// Java: InlineFunctionsTest.testInlineWithThis3()
func (c *InlineFunctionsTest) TestInlineWithThis3() {
	c.assumeStrictThis = false
	// If no "this" is provided it might need to be coerced to the global
	// "this".
	c.TestSame("function f(){} f.call([]);")
	c.assumeStrictThis = true
	// In strict mode, "this" is never coerced so we can use the provided value.
	c.Test("function f(){} f.call([]);", "{}")
}

// Java: InlineFunctionsTest.testInlineWithThis4()
func (c *InlineFunctionsTest) TestInlineWithThis4() {
	c.assumeStrictThis = false
	// If no "this" is provided it might need to be coerced to the global
	// "this".
	c.TestSame("function f(){} f.call(new g);")
	c.assumeStrictThis = true
	// In strict mode, "this" is never coerced so we can use the provided value.
	c.Test("function f(){} f.call(new g);", "{var JSCompiler_inline_this_v0=new g}")
}

// Java: InlineFunctionsTest.testInlineWithThis5()
func (c *InlineFunctionsTest) TestInlineWithThis5() {
	c.assumeStrictThis = false
	// If no "this" is provided it might need to be coerced to the global
	// "this".
	c.TestSame("function f(){} f.call(g());")
	c.assumeStrictThis = true
	// In strict mode, "this" is never coerced so we can use the provided value.
	c.Test("function f(){} f.call(g());", "{var JSCompiler_inline_this_v0=g()}")
}

// Java: InlineFunctionsTest.testInlineWithThis6()
func (c *InlineFunctionsTest) TestInlineWithThis6() {
	c.assumeStrictThis = false
	// If no "this" is provided it might need to be coerced to the global
	// "this".
	c.TestSame("function f(){this} f.call(new g);")
	c.assumeStrictThis = true
	// In strict mode, "this" is never coerced so we can use the provided value.
	c.Test(
		"function f(){this} f.call(new g);",
		"{var JSCompiler_inline_this_v0=new g;JSCompiler_inline_this_v0}")
}

// Java: InlineFunctionsTest.testInlineWithThis7()
func (c *InlineFunctionsTest) TestInlineWithThis7() {
	c.assumeStrictThis = true
	// In strict mode, "this" is never coerced so we can use the provided value.
	c.Test(
		"function f(a){a=1;this} f.call();",
		"{var a$jscomp$inline_0=void 0; a$jscomp$inline_0=1; void 0;}")
	c.uniqueIdSupplier = NewInlineFunctionsTest_UniqueIdSupplier()
	c.Test(
		"function f(a){a=1;this} f.call(x, x);",
		`{var JSCompiler_inline_this_v0 = x; var a$jscomp$inline_0=x; a$jscomp$inline_0=1;
    JSCompiler_inline_this_v0;}
`)
}

// http://en.wikipedia.org/wiki/Fixed_point_combinator#Y_combinator
// Java: InlineFunctionsTest.testFunctionExpressionYCombinator()
func (c *InlineFunctionsTest) TestFunctionExpressionYCombinator() {
	c.assumeMinimumCapture = false
	c.TestSame(
		`var factorial = ((function(M) {
      return ((function(f) {
                 return M(function(arg) {
                            return (f(f))(arg);
                            })
               })
              (function(f) {
                 return M(function(arg) {
                            return (f(f))(arg);
                           })
                 }));
     })
    (function(f) {
       return function(n) {
        if (n === 0)
          return 1;
        else
          return n * f(n - 1);
       };
     }));

factorial(5)
`)
	c.assumeMinimumCapture = true
	c.Test(
		`var factorial = ((function(M) {
      return ((function(f) {
                 return M(function(arg) {
                            return (f(f))(arg);
                            })
               })
              (function(f) {
                 return M(function(arg) {
                            return (f(f))(arg);
                           })
                 }));
     })
    (function(f) {
       return function(n) {
        if (n === 0)
          return 1;
        else
          return n * f(n - 1);
       };
     }));

factorial(5)
`,
		`var factorial;
{
  var M$jscomp$inline_2 = function(f$jscomp$2) {
    return function(n) {
      if (n === 0) {
        return 1;
      } else {
        return n * f$jscomp$2(n - 1);
      }
    };
  };
  {
    var f$jscomp$inline_7 = function(f$jscomp$inline_5) {
      return M$jscomp$inline_2(function(arg$jscomp$inline_6) {
        return f$jscomp$inline_5(f$jscomp$inline_5)(arg$jscomp$inline_6);
      });
    };
    factorial = M$jscomp$inline_2(function(arg$jscomp$inline_8) {
      return f$jscomp$inline_7(f$jscomp$inline_7)(arg$jscomp$inline_8);
    });
  }
}
factorial(5);
`)
}

// Java: InlineFunctionsTest.testRenamePropertyFunction()
func (c *InlineFunctionsTest) TestRenamePropertyFunction() {
	c.TestSame(
		`function JSCompiler_renameProperty(x) {return x}
JSCompiler_renameProperty('foo')
`)
}

// Java: InlineFunctionsTest.testReplacePropertyFunction()
func (c *InlineFunctionsTest) TestReplacePropertyFunction() {
	// baseline: an alias doesn't prevents declaration removal, but not
	// inlining.
	c.Test(
		`function f(x) {return x}
foo(window, f); f(1)
`,
		`function f(x) {return x}
foo(window, f); 1
`)
}

// Java: InlineFunctionsTest.testInlineWithClosureContainingThis()
func (c *InlineFunctionsTest) TestInlineWithClosureContainingThis() {
	c.Test("(function (){return f(function(){return this})})();", "f(function(){return this})")
}

// Java: InlineFunctionsTest.testIssue5159924a()
func (c *InlineFunctionsTest) TestIssue5159924a() {
	c.Test(
		`function f() { if (x()) return y(); }
while(1) { var m = f() || z() }
`,
		`for(;1;) {
  var JSCompiler_inline_result$jscomp$v0;
  {
    JSCompiler_inline_label_f_v1: {
      if(x()) {
        JSCompiler_inline_result$jscomp$v0 = y();
        break JSCompiler_inline_label_f_v1;
      }
      JSCompiler_inline_result$jscomp$v0 = void 0;
    }
  }
  var m = JSCompiler_inline_result$jscomp$v0 || z()
}
`)
}

// Java: InlineFunctionsTest.testIssue5159924b()
func (c *InlineFunctionsTest) TestIssue5159924b() {
	c.Test(
		`function f() { if (x()) return y(); }
while(1) { var m = f(); }
`,
		`for(;1;) {
  var m;
  {
    JSCompiler_inline_label_f_v0: {
      if(x()) {
        m = y();
        break JSCompiler_inline_label_f_v0;
      }
      m = void 0;
    }
  }
}
`)
}

// Java: InlineFunctionsTest.testInlineObject()
func (c *InlineFunctionsTest) TestInlineObject() {
	c.DisableCompareAsTree()
	c.EnableComputeSideEffects()
	c.inliningReach = CompilerOptions_Reach_LOCAL_ONLY
	c.assumeStrictThis = true
	c.assumeMinimumCapture = true
	// TODO(johnlenz): normalize the AST so an AST comparison can be done.
	// As is, the expected AST does not match the actual correct result:
	// The AST matches "g.a()" with a FREE_CALL annotation, but this as
	// expected string would fail as it won't be mark as a free call.
	// "(0,g.a)()" matches the output, but not the resulting AST.
	c.Test("function inner(){function f(){return g.a}(f())()}", "function inner(){(0,g.a)()}")
}

// Java: InlineFunctionsTest.testReproduceBug401582520()
func (c *InlineFunctionsTest) TestReproduceBug401582520() {
	// InlineFunctions is inlining reading a property before the side-effectful calls that
	// initializes it occurs.
	c.Test(
		`  var module$exports$foo$Foo$$0clinit = function() {
    module$exports$foo$Foo$$0clinit = function() {
    };
    module$exports$foo$Foo$options_ = ["a", "b", "c", "d"];
  };

          var module$exports$foo$Foo$create = function(o$jscomp$1) {
            module$exports$foo$Foo$$0clinit();
            return new module$exports$foo$Foo(o$jscomp$1);
          };

          var module$exports$foo$Foo = class {
            constructor(o) {
              this.o_ = o;
            }
          };

          var module$exports$foo$Foo$options_;

  console.log(function(JSCompiler_StaticMethods_isKnown$self) {
    return module$exports$foo$Foo$options_.includes(JSCompiler_StaticMethods_isKnown$self.o_);
  }(function(o$jscomp$2) {
    return module$exports$foo$Foo$create(o$jscomp$2);
  }("abc")));
`,
		`var module$exports$foo$Foo$$0clinit = function() {
  module$exports$foo$Foo$$0clinit = function() {
  };
  module$exports$foo$Foo$options_ = ["a", "b", "c", "d"];
};
var module$exports$foo$Foo = class {
  constructor(o) {
    this.o_ = o;
  }
};
var module$exports$foo$Foo$options_;
var JSCompiler_temp_const$jscomp$v1 = console;
var JSCompiler_temp_const$jscomp$v0 = JSCompiler_temp_const$jscomp$v1.log;
var JSCompiler_inline_result$jscomp$v2;
{
  var JSCompiler_StaticMethods_isKnown$self$jscomp$inline_0;
  {
    module$exports$foo$Foo$$0clinit();
    JSCompiler_StaticMethods_isKnown$self$jscomp$inline_0 = new module$exports$foo$Foo("abc");
  }
  JSCompiler_inline_result$jscomp$v2 =
      module$exports$foo$Foo$options_.includes(
          JSCompiler_StaticMethods_isKnown$self$jscomp$inline_0.o_);
}
JSCompiler_temp_const$jscomp$v0.call(
    JSCompiler_temp_const$jscomp$v1,
    JSCompiler_inline_result$jscomp$v2);
`)
}

// Java: InlineFunctionsTest.testReproduceBug401582520_simplified()
func (c *InlineFunctionsTest) TestReproduceBug401582520_simplified() {
	// InlineFunctions is inlining reading a property before the side-effectful calls that
	// initializes it occurs.
	c.Test(
		`var options_;

/** @noinline */
var initAndReturnA = function() {
  options_ = ["a", "b", "c", "d"];
  return "a";
};

console.log(function(x) {
  return options_.includes(x);
}(initAndReturnA()));
`,
		`var options_;

/** @noinline */
var initAndReturnA = function() {
  options_ = ["a", "b", "c", "d"];
  return "a";
};

var JSCompiler_temp_const$jscomp$v1 = console;
var JSCompiler_temp_const$jscomp$v0 = JSCompiler_temp_const$jscomp$v1.log;
var JSCompiler_inline_result$jscomp$v2;
{
  var x$jscomp$inline_0 = initAndReturnA();
  JSCompiler_inline_result$jscomp$v2 = options_.includes(x$jscomp$inline_0);
}
JSCompiler_temp_const$jscomp$v0.call(JSCompiler_temp_const$jscomp$v1, JSCompiler_inline_result$jscomp$v2);
`)
}

// Java: InlineFunctionsTest.testBug4944818()
func (c *InlineFunctionsTest) TestBug4944818() {
	c.Test(
		`var getDomServices_ = function(self) {
  if (!self.domServices_) {
    self.domServices_ = goog$component$DomServices.get(
        self.appContext_);
  }

  return self.domServices_;
};

var getOwnerWin_ = function(self) {
  return getDomServices_(self).getDomHelper().getWindow();
};

HangoutStarter.prototype.launchHangout = function() {
  var self = a.b;
  var myUrl = new goog.Uri(getOwnerWin_(self).location.href);
};
`,
		`HangoutStarter.prototype.launchHangout = function() {
  var self$jscomp$2 = a.b;
  var JSCompiler_temp_const$jscomp$v0 = goog.Uri;
  var JSCompiler_inline_result$jscomp$v1;
  {
  var self$jscomp$inline_0 = self$jscomp$2;
  if (!self$jscomp$inline_0.domServices_) {
    self$jscomp$inline_0.domServices_ = goog$component$DomServices.get(
        self$jscomp$inline_0.appContext_);
  }
  JSCompiler_inline_result$jscomp$v1=self$jscomp$inline_0.domServices_;
  }
  var myUrl = new JSCompiler_temp_const$jscomp$v0(
      JSCompiler_inline_result$jscomp$v1.getDomHelper().
          getWindow().location.href)
}
`)
}

// http://blickly.github.io/closure-compiler-issues/#423
// Java: InlineFunctionsTest.testIssue423()
func (c *InlineFunctionsTest) TestIssue423() {
	c.assumeMinimumCapture = false
	c.Test(
		`(function($) {
  $.fn.multicheck = function(options) {
    initialize.call(this, options);
  };

  function initialize(options) {
    options.checkboxes = $(this).siblings(':checkbox');
    preload_check_all.call(this);
  }

  function preload_check_all() {
    $(this).data('checkboxes');
  }
})(jQuery)
`,
		`(function($){
  $.fn.multicheck=function(options$jscomp$1) {
    {
      options$jscomp$1.checkboxes=$(this).siblings(':checkbox');
      {
        $(this).data('checkboxes');
      }
    }
  }
})(jQuery)
`)
}

// http://blickly.github.io/closure-compiler-issues/#423
// Java: InlineFunctionsTest.testIssue423_minCap()
func (c *InlineFunctionsTest) TestIssue423_minCap() {
	c.assumeMinimumCapture = true
	c.Test(
		`(function($) {
  $.fn.multicheck = function(options) {
    initialize.call(this, options);
  };

  function initialize(options) {
    options.checkboxes = $(this).siblings(':checkbox');
    preload_check_all.call(this);
  }

  function preload_check_all() {
    $(this).data('checkboxes');
  }
})(jQuery)
`,
		`{
  var $$jscomp$inline_0=jQuery;
  $$jscomp$inline_0.fn.multicheck=function(options$jscomp$inline_4) {
    {
      options$jscomp$inline_4.checkboxes =
          $$jscomp$inline_0(this).siblings(':checkbox');
      {
         $$jscomp$inline_0(this).data('checkboxes')
      }
    }
  }
}
`)
}

// http://blickly.github.io/closure-compiler-issues/#728
// Java: InlineFunctionsTest.testIssue728()
func (c *InlineFunctionsTest) TestIssue728() {
	f := "var f = function() { return false; };"
	calls := jutil.NewStringBuilder()
	folded := jutil.NewStringBuilder()
	for i := 0; i < 30; i++ {
		calls.Append("if (!f()) alert('x');")
		folded.Append("if (!false) alert('x');")
	}
	c.Test(f+calls.String(), folded.String())
}

// Java: InlineFunctionsTest.testAnonymous1()
func (c *InlineFunctionsTest) TestAnonymous1() {
	c.assumeMinimumCapture = false
	c.Test(
		"(function(){var a=10;(function(){var b=a;a++;alert(b)})()})();",
		`{
  var a$jscomp$inline_1 = 10;
  {
    var b$jscomp$inline_2 = a$jscomp$inline_1;
    a$jscomp$inline_1++;
    alert(b$jscomp$inline_2);
  }
}
`)
	c.assumeMinimumCapture = true
	c.Test(
		"(function(){var a=10;(function(){var b=a;a++;alert(b)})()})();",
		`{
  var a$jscomp$inline_1 = 10;
  {
    var b$jscomp$inline_3 = a$jscomp$inline_1;
    a$jscomp$inline_1++;
    alert(b$jscomp$inline_3);
  }
}
`)
}

// Java: InlineFunctionsTest.testAnonymous2()
func (c *InlineFunctionsTest) TestAnonymous2() {
	c.TestSame(
		`(function(){
  eval();
  (function(){
    var b=a;
    a++;
    alert(b)
  })()
})();
`)
}

// Java: InlineFunctionsTest.testAnonymous3()
func (c *InlineFunctionsTest) TestAnonymous3() {
	// Introducing a new value into is tricky
	c.assumeMinimumCapture = false
	c.TestSame("(function(){var a=10;(function(){arguments;})()})();")
	c.assumeMinimumCapture = true
	c.Test(
		"(function(){var a=10;(function(){arguments;})()})();",
		"{var a$jscomp$inline_0=10;(function(){arguments;})();}")
	c.Test("(function(){(function(){arguments;})()})();", "{(function(){arguments;})()}")
}

// Java: InlineFunctionsTest.testLoopWithFunctionWithFunction()
func (c *InlineFunctionsTest) TestLoopWithFunctionWithFunction() {
	c.assumeMinimumCapture = true
	c.Test(
		`function _testLocalVariableInLoop_() {
  var result = 0;
  function foo() {
    var arr = [1, 2, 3, 4, 5];
    for (var i = 0, l = arr.length; i < l; i++) {
      var j = arr[i];
      // Don't inline this function; the correct behavior depends on captured values.
      (function() {
        var k = j;
        setTimeout(function() { result += k; }, 5 * i);
      })();
    }
  }
  foo();
}
`,
		`function _testLocalVariableInLoop_() {
  var result = 0;
  {
    var arr$jscomp$inline_0 = [1,2,3,4,5];
    var i$jscomp$inline_1 = 0;
    var l$jscomp$inline_2 = arr$jscomp$inline_0.length;
    for(; i$jscomp$inline_1 < l$jscomp$inline_2; i$jscomp$inline_1++) {
      var j$jscomp$inline_3 = arr$jscomp$inline_0[i$jscomp$inline_1];
      (function(){
        var k$jscomp$inline_4 = j$jscomp$inline_3;
        setTimeout(function() {
          result += k$jscomp$inline_4;
        }, 5 * i$jscomp$inline_1)
      })()
    }
  }
}
`)
}

// Java: InlineFunctionsTest.testMethodWithFunctionWithFunction()
func (c *InlineFunctionsTest) TestMethodWithFunctionWithFunction() {
	c.assumeMinimumCapture = true
	c.Test(
		`function _testLocalVariable_() {
  var result = 0;
  function foo() {
      var j = [i];
      (function(j) {
        setTimeout(function() { result += j; }, 5 * i);
      })(j);
      j = null;
  }
  foo();
}
`,
		`function _testLocalVariable_(){
  var result = 0;
  {
    var j$jscomp$inline_1 = [i];
    {
      var j$jscomp$inline_3 = j$jscomp$inline_1; // this temp is needed.
      setTimeout(function(){ result += j$jscomp$inline_3; }, 5*i);
    }
    j$jscomp$inline_1 = null; // because this value can be modified later.
  }
}
`)
}

// Inline a single reference function into deeper chunks
// Java: InlineFunctionsTest.testCrossChunkInlining1()
func (c *InlineFunctionsTest) TestCrossChunkInlining1() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				// m1
				AddChunk("function foo(){return f(1)+g(2)+h(3);}").
				// m2
				AddChunk("foo()").
				Build()),
		c.Expected(
			// m1
			"",
			// m2
			"f(1)+g(2)+h(3);"))
}

// Inline a single reference function into shallow chunks, only if it
// is cheaper than the call itself.
// Java: InlineFunctionsTest.testCrossChunkInlining2()
func (c *InlineFunctionsTest) TestCrossChunkInlining2() {
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("foo()").
				AddChunk("function foo(){return f(1)+g(2)+h(3);}").
				Build()))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				// m1
				AddChunk("foo()").
				// m2
				AddChunk("function foo(){return f();}").
				Build()),
		c.Expected( // m1
			"f();",
			// m2
			""))
}

// Inline a multi-reference functions into shallow chunks, only if it
// is cheaper than the call itself.
// Java: InlineFunctionsTest.testCrossChunkInlining3()
func (c *InlineFunctionsTest) TestCrossChunkInlining3() {
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("foo()").
				AddChunk("function foo(){return f(1)+g(2)+h(3);}").
				AddChunk("foo()").
				Build()))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("foo()").
				AddChunk("function foo(){return f();}").
				AddChunk("foo()").
				Build()),
		c.Expected(
			// m1
			"f();",
			// m2
			"",
			// m3
			"f();"))
}

// Java: InlineFunctionsTest.test6671158()
func (c *InlineFunctionsTest) Test6671158() {
	tester := newInlineFunctionsTest_test6671158Tester()
	tester.JUnitSetT(c.T())

	tester.SetUp()
	tester.Test(
		`function f() {return g()}
function Y(a){a.loader_()}
function _Z(){}
function _X() { new _Z(a,b, Y(singleton), f()) }
`,
		`function _Z(){}
function _X(){
  var JSCompiler_temp_const$jscomp$v2=_Z;
  var JSCompiler_temp_const$jscomp$v1=a;
  var JSCompiler_temp_const$jscomp$v0=b;
  var JSCompiler_inline_result$jscomp$v3;
  {
    singleton.loader_();
    JSCompiler_inline_result$jscomp$v3=void 0;
  }
  new JSCompiler_temp_const$jscomp$v2(
    JSCompiler_temp_const$jscomp$v1,
    JSCompiler_temp_const$jscomp$v0,
    JSCompiler_inline_result$jscomp$v3,
    g())}
`)
	tester.TearDown()
}

// Java: InlineFunctionsTest.test6671158b()
func (c *InlineFunctionsTest) Test6671158b() {
	c.Test(
		`function f() {return g()}
function Y(a){a.loader_()}
function _Z(){}
function _X() { new _Z(a,b, Y(singleton), f()) }
`,
		`function _Z(){}
function _X(){
  var JSCompiler_temp_const$jscomp$v1=a;
  var JSCompiler_temp_const$jscomp$v0=b;
  var JSCompiler_inline_result$jscomp$v2;
  {
    singleton.loader_();
    JSCompiler_inline_result$jscomp$v2=void 0;
  }
  new _Z(
    JSCompiler_temp_const$jscomp$v1,
    JSCompiler_temp_const$jscomp$v0,
    JSCompiler_inline_result$jscomp$v2,
    g())}
`)
}

// Java: InlineFunctionsTest.test8609285a()
func (c *InlineFunctionsTest) Test8609285a() {
	c.Test(
		"function f(x){ for(x in y){} } f()",
		"{var x$jscomp$inline_0=void 0;for(x$jscomp$inline_0 in y);}")
}

// Java: InlineFunctionsTest.test8609285b()
func (c *InlineFunctionsTest) Test8609285b() {
	c.Test(
		"function f(x){ for(var x in y){} } f()",
		"{var x$jscomp$inline_0=void 0;for(x$jscomp$inline_0 in y);}")
}

// http://blickly.github.io/closure-compiler-issues/#1101
// Java: InlineFunctionsTest.testIssue1101()
func (c *InlineFunctionsTest) TestIssue1101() {
	c.Test(
		`var x = (function (saved) {
    return modifyObjProp(obj) + saved;
  })(obj["prop"]);
`,
		`var x;
{
  var saved$jscomp$inline_0=obj["prop"];x=modifyObjProp(obj)+
     saved$jscomp$inline_0
}
`)
}

// Java: InlineFunctionsTest.testMaxFunSizeAfterInlining()
func (c *InlineFunctionsTest) TestMaxFunSizeAfterInlining() {
	c.maxSizeAfterInlining = 1
	c.Test( // Always inline single-statement functions
		`function g() { return 123; }
function f() { g(); }
`,
		"function f() { 123; }")
	c.maxSizeAfterInlining = 10
	c.Test( // Always inline at the top level
		`function g() { 123; return 123; }
g();
`,
		"{ 123; 123; }")
	c.maxSizeAfterInlining = 1
	c.TestSame( // g is too big to be inlined
		`function g() { 123; return 123; }
g();
`)
	c.maxSizeAfterInlining = 20
	c.Test(
		`function g() { 123; return 123; }
function f() {
  g();
}
`,
		"")
	// g's size ends up exceeding the max size because all inlining decisions
	// were made in the same inlining round.
	c.maxSizeAfterInlining = 25
	c.Test(
		`function f1() { 1; return 1; }
function f2() { 2; return 2; }
function f3() { 3; return 3; }
function f4() { 4; return 4; }
function g() {
  f1(); f2(); f3(); f4();
}
g(); g(); g();
`,
		"function g() { {1; 1;} {2; 2;} {3; 3;} {4; 4;} }\ng(); g(); g();")
	c.maxSizeAfterInlining = CompilerOptions_UNLIMITED_FUN_SIZE_AFTER_INLINING
}

// Java: InlineFunctionsTest.testArrowFunctionUndefined()
func (c *InlineFunctionsTest) TestArrowFunctionUndefined() {
	c.Test("var empty = () => {}; empty();", "void 0;")
}

// Java: InlineFunctionsTest.testFunctionProperty()
func (c *InlineFunctionsTest) TestFunctionProperty() {
	c.TestSame(
		`function Foo() {
  this.bar = 0;
}
var a = new Foo();
a.bar;
`)
}

// Java: InlineFunctionsTest.testFunctionRestParam()
func (c *InlineFunctionsTest) TestFunctionRestParam() {
	c.Test(
		`var f = function(...args) { return args[0]; }
f(8);
`,
		"[8][0]")
	c.Test(
		`var f = function(...args) { return args.length; }
f();
`,
		"[].length;")
	c.Test(
		`var f = function(...args) { return args.length; }
f(3, 4);
`,
		"[3, 4].length;")
}

// Java: InlineFunctionsTest.testArrowFunctionRestParam()
func (c *InlineFunctionsTest) TestArrowFunctionRestParam() {
	c.Test(
		`function foo() {
  var f = (...args) => args[0];
  return f(8);
}
foo();
`,
		"[8][0]")
	c.Test(
		`function foo() {
  var f = (x, ...args) => x + args[0];
  return f(1, 8);
}
foo();
`,
		`{
  var JSCompiler_inline_result$jscomp$inline_2;
  {
    var x$jscomp$inline_3 = 1;
    var args$jscomp$inline_4 = [8];
    JSCompiler_inline_result$jscomp$inline_2 = x$jscomp$inline_3 + args$jscomp$inline_4[0];
  }
  JSCompiler_inline_result$jscomp$inline_2
}
`)
	c.Test(
		`function foo() {
  var f = (...args) => args[0].bar;
  return f({bar: 8});
}
foo();
`,
		"[{bar: 8}][0].bar")
}

// Java: InlineFunctionsTest.testRestObjectPattern()
func (c *InlineFunctionsTest) TestRestObjectPattern() {
	c.TestSame(
		`function countArgs(...{length}) {
  return length;
}
countArgs(1, 1, 1, 1, 1);
`)
	c.TestSame(
		`function countArgs(x, ...{length}) {
  return length;
}
countArgs(1, 1, 1, 1, 1);
`)
	c.TestSame(
		`function countArgs(x, ...{length: length}) {
  return length;
}
countArgs(1, 1, 1, 1, 1);
`)
	c.TestSame(
		`function f(...{'a': x}) {
  return x;
}
f(null,null,null,3,null);
`)
	c.TestSame(
		`function f(...{3: x}) {
  return x;
}
f(null,null,null,3,null);
`)
	c.TestSame(
		`function f(...{x: y}) {
  return y;
}
f(null,null,null,3,null);
`)
	c.TestSame(
		`function f(...{p: x, 3:y}) {
  return y;
}
f(null, null, null, 3, null);
`)
}

// Java: InlineFunctionsTest.testObjectPatternParam()
func (c *InlineFunctionsTest) TestObjectPatternParam() {
	c.TestSame(
		`function foo({x}) {
  return x+1;
}
foo({x:5});
`)
	c.TestSame(
		`function foo({x:y}) {
  return y+1;
}
foo({x:5});
`)
	c.TestSame(
		`function foo({'x':y}) {
  return y+1;
}
foo({x:5});
`)
	c.TestSame(
		`function foo({x}, {y}) {
  return x+y;
}
foo({x:5}, {y:6});
`)
	c.TestSame(
		`function foo({x}, {y}) {
  return x+1;
}
foo({x:5}, {y:6});
`)
	c.TestSame(
		`function foo({x: {y}}) {
  return y + 1;
}
foo(obj);
`)
	c.TestSame(
		`function foo({a: b, x: {y}}) {
  return y + 1;
}
foo(obj);
`)
	c.TestSame(
		`function f({x}) {
  return x;
}
class Foo {constructor() {this.x = 0;}}
f(new Foo());
`)
	c.TestSame(
		`function f({x}) {
  alert(x);
  alert(x)
}
class Foo {constructor() {this.x = 0;}}
f(new Foo());
`)
	c.TestSame(
		`function f({x}, {y}) {
  alert(x);
  alert(x);
  return y;
}
class Foo {constructor() {this.x = 0;}}
f(new Foo(), {y:6});
`)
	c.TestSame(
		`function f({x, y}) {
  return x + y;
}
f(obj);
`)
	c.TestSame(
		`function f({x, y}, {z}) {
  alert(z);
  return x + y;
}
f(obj, new Foo());
`)
	c.TestSame(
		`function f({x, y}) {
  return x + y;
}
f(getArg());
`)
	c.TestSame(
		`function f({a, b, c}) {
  return b + c;
}
f(x);
`)
	c.TestSame(
		`function f({3:x}) {
  return x;
}
f({3:1});
`)
	c.TestSame(
		`function f({x:{3:y}}) {
  return y;
}
f({x:{3:1}});
`)
	c.TestSame(
		`function f({p: x, 3: y}) {
  return x;
}
f({p:1, 3:2});
`)
	c.TestSame(
		`function f({prop1, prop2}) {
  return prop1;
}
f({prop1:5, prop2:6});
`)
	c.TestSame(
		`function f({'foo bar':x}) {
  return x;
}
f({'foo bar': 1});
`)
	c.TestSame(
		`function f({'foo_bar':x}) {
  return x;
}
f({'foo_bar': 1});
`)
	c.TestSame(
		`function f({123: x}) {
 return x;
}
f({123: 1});
`)
	c.TestSame(
		`function f({'1foo': x}) {
  return x;
}
f({'1foo': 1});
`)
	c.TestSame(
		`function f({foo1 : x}) {
  return x;
}
f({foo1 : 1});
`)
	c.TestSame(
		`function f({'foo1': x}) {
  return x;
}
f({'foo1': 1});
`)
	c.TestSame(
		`function f({$foo : x}) {
  return x;
}
f({$foo : 1});
`)
	c.TestSame(
		`function f({_foo : x}) {
  return x;
}
f({_foo : 1});
`)
}

// Java: InlineFunctionsTest.testDefaultObjectPatternParam()
func (c *InlineFunctionsTest) TestDefaultObjectPatternParam() {
	c.TestSame(
		`function foo({x} = {x:5}) {
  return x+1;
}
foo();
`)
	c.TestSame(
		`function foo({x} = {x:5}, {y} = {y:3}) {
  return x+y;
}
foo();
`)
	c.TestSame(
		`let defaultObj = {x: 5};
function foo({x} = defaultObj) {
  return x;
}
foo();
`)
	c.TestSame(
		`function f({a, b, c} = {a:1, b:2, c:3}) {
 return b+c;
}
f();
`)
	c.TestSame(
		`function f({p:x, 3:y} = {p:1, 3:2}) {
 return x+y;
}
f();
`)
	c.TestSame(
		`function foo(a, { b = '', c = '' } = {}) {
  return a;
}
foo();
`)
}

// Java: InlineFunctionsTest.testObjectPatternParamWithMultipleDestructuredNames()
func (c *InlineFunctionsTest) TestObjectPatternParamWithMultipleDestructuredNames() {
	c.TestSame(
		`function f({x, y}) {
  return x + y;
}
alert(f({x: sideEffects1(), y: sideEffects2()}));
`)
}

// Java: InlineFunctionsTest.testArrayPatternParam()
func (c *InlineFunctionsTest) TestArrayPatternParam() {
	c.TestSame("function f([x]) { return x; } alert(f([3]));")
}

// Java: InlineFunctionsTest.testObjectPatternParamWithDefaults()
func (c *InlineFunctionsTest) TestObjectPatternParamWithDefaults() {
	c.TestSame(
		`function foo({x = 3}) {
  return x;
}
alert(foo({}));
`)
	c.TestSame(
		`function foo({x = 3}) {
  return x;
}
alert(foo({x: 3}));
`)
	c.TestSame(
		`function foo({x = 3}, y) {
  return y;
}
alert(foo({x: 3}, 4));
`)
	c.TestSame(
		`function foo({x = 3} = {}) {
  return x;
}
alert(foo());
`)
	c.TestSame(
		`function foo(x, ...{length = 5}) {
  return x * length;
}
alert(foo(3));
`)
}

// Java: InlineFunctionsTest.testDefaultParam()
func (c *InlineFunctionsTest) TestDefaultParam() {
	c.TestSame(
		`function foo(a, b = 1) {
  return a + b;
}
foo(1);
`)
	c.TestSame(
		`function foo(a, b = 1) {
  return a + b;
}
foo(1, 2);
`)
	c.TestSame(
		`function foo(a = 1, b = 2) {
  return a + b;
}
foo(3, 4);
`)
	c.TestSame(
		`function foo(a, b = {foo: 5}) {
  return a + b.foo;
}
foo(3, {foo: 9});
`)
	c.TestSame(
		`function foo(a, b = {'foo': 5}) {
  return a + b['foo'];
}
foo(3, {'foo': 9});
`)
	c.TestSame(
		`function foo(a, b = {foo: 5, bar: 6}) {
  return a + b.foo + b.bar;
}
foo(3, {foo: 1, bar: 2});
`)
	c.TestSame(
		`function foo(a, b = {foo: 5}) {
  return a + b.foo;
}
foo(3);
`)
	c.TestSame(
		`function foo(a, b = {foo: 5, bar: 6}) {
  return a + b.foo + b.bar;
}
foo(3, {foo: 1});
`)
	c.TestSame(
		`function foo(a, b = [1, 2]) {
  return a + b[1];
}
foo(3, [7, 8]);
`)
	c.TestSame(
		`function foo(a, b = []) {
  return a + b[1];
}
foo(3, [7, 8]);
`)
}

// Java: InlineFunctionsTest.testDefaultParam_argIsUndefined()
func (c *InlineFunctionsTest) TestDefaultParam_argIsUndefined() {
	c.TestSame(
		`function foo(a, b = 1) {
  return a + b;
}
foo(1, undefined);
`)
}

// Java: InlineFunctionsTest.testDefaultParam_argIsUnknown()
func (c *InlineFunctionsTest) TestDefaultParam_argIsUnknown() {
	c.TestSame(
		`function foo(a, b = 1) {
  return a + b;
}
foo(1, x);
`)
}

// Java: InlineFunctionsTest.testDefaultParam_withAssign()
func (c *InlineFunctionsTest) TestDefaultParam_withAssign() {
	c.TestSame(
		`function foo(x = undefined) {
  if (!x) {
    x = 2;
  }
  return x;
}
foo(4);
`)
}

// Java: InlineFunctionsTest.testNestedDefaultParam()
func (c *InlineFunctionsTest) TestNestedDefaultParam() {
	c.TestSame(
		`function foo(a = b = 1) {
  return a;
}
foo();
`)
	c.TestSame(
		`function foo(c = {a:(b = 1)}) {
  return c;
}
foo();
`)
}

// Java: InlineFunctionsTest.testSpreadCall_withKnownArray_isNotInlined_andFunctionIsPreserved()
func (c *InlineFunctionsTest) TestSpreadCall_withKnownArray_isNotInlined_andFunctionIsPreserved() {
	c.TestSame(
		`function foo(x, y) {
  return x + y;
}
var args = [0, 1];
foo(...args);
`)
}

// Java: InlineFunctionsTest.testSpreadCall_withKnownArray_andLiteralArg_isNotInlined_andFunctionIsPreserved()
func (c *InlineFunctionsTest) TestSpreadCall_withKnownArray_andLiteralArg_isNotInlined_andFunctionIsPreserved() {
	c.TestSame(
		`function foo(x, y, z) {
  return x + y + z;
}
var args = [0, 1];
foo(2, ...args);
`)
}

// Java: InlineFunctionsTest.testSpreadCall_withArrayLiteral_isNotInlined_andFunctionIsPreserved()
func (c *InlineFunctionsTest) TestSpreadCall_withArrayLiteral_isNotInlined_andFunctionIsPreserved() {
	c.TestSame(
		`function foo(x, y) {
  return x + y;
}
foo(...[0, 1]);
`)
}

// Java: InlineFunctionsTest.testSpreadCall_withArrayLiteral_isNotInlined_andFunctionIsPreserved_optChain()
func (c *InlineFunctionsTest) TestSpreadCall_withArrayLiteral_isNotInlined_andFunctionIsPreserved_optChain() {
	c.TestSame(
		`function foo(x, y) {
  return x + y;
}
foo?.(...[0, 1]);
`)
}

// Java: InlineFunctionsTest.testSpreadCall_withMultipleSpreads_isNotInlined_andFunctionIsPreserved()
func (c *InlineFunctionsTest) TestSpreadCall_withMultipleSpreads_isNotInlined_andFunctionIsPreserved() {
	c.TestSame(
		`function foo(x, y) {
  return x + y;
}
var args = [0];
foo(...args, ...[1]);
`)
}

// Java: InlineFunctionsTest.testSpreadCall_withArrayLiteral_forRestParam_isNotInlined_andFunctionIsPreserved()
func (c *InlineFunctionsTest) TestSpreadCall_withArrayLiteral_forRestParam_isNotInlined_andFunctionIsPreserved() {
	c.TestSame(
		`function foo(...args) {
  return args.length;
}
foo(...[0,1]);
`)
}

// Java: InlineFunctionsTest.testSpreadCall_withKnownArray_intoIife_isNotInlined()
func (c *InlineFunctionsTest) TestSpreadCall_withKnownArray_intoIife_isNotInlined() {
	c.TestSame(
		`var args = [0, 1];
(function foo(x, y) { return x + y; })(...args);
`)
}

// Java: InlineFunctionsTest.testSpreadCall_withKnownArray_andLiteralArg_intoIife_isNotInlined()
func (c *InlineFunctionsTest) TestSpreadCall_withKnownArray_andLiteralArg_intoIife_isNotInlined() {
	c.TestSame(
		`var args = [0, 1];
(function foo(x, y, z) {
  return x + y + z;
})(2, ...args);
`)
}

// Java: InlineFunctionsTest.testSpreadCall_withArrayLiteral_intoIife_isNotInlined()
func (c *InlineFunctionsTest) TestSpreadCall_withArrayLiteral_intoIife_isNotInlined() {
	c.TestSame("(function (x, y) {  return x + y; })(...[0, 1]);")
}

// Java: InlineFunctionsTest.testSpreadCall_withSpreadAsIndirectDescendent_call_isInlined()
func (c *InlineFunctionsTest) TestSpreadCall_withSpreadAsIndirectDescendent_call_isInlined() {
	c.Test(
		`function foo(x) {
  return x;
}

foo(bar(...arg));
`,
		`bar(...arg);
`)
}

// Java: InlineFunctionsTest.testSpreadCall_withSpreadAsIndirectDescendent_arrrayLit_isInlined()
func (c *InlineFunctionsTest) TestSpreadCall_withSpreadAsIndirectDescendent_arrrayLit_isInlined() {
	c.Test(
		`function foo(x) {
  return x;
}

foo([...arg]);
`,
		`[...arg];
`)
}

// Java: InlineFunctionsTest.testSpreadCall_withBlockInlinableFunction_asSpreadArg_isInlined()
func (c *InlineFunctionsTest) TestSpreadCall_withBlockInlinableFunction_asSpreadArg_isInlined() {
	c.Test(
		`function foo(x) {
  qux(); // Something with a side-effect.
  return x;
}

bar(...foo(arg));
`,
		`var JSCompiler_temp_const$jscomp$v0 = bar;
var JSCompiler_inline_result$jscomp$v1;
{
  var x$jscomp$inline_0 = arg;
  qux();
  JSCompiler_inline_result$jscomp$v1 = x$jscomp$inline_0;
}

JSCompiler_temp_const$jscomp$v0(...JSCompiler_inline_result$jscomp$v1);
`)
}

// Java: InlineFunctionsTest.testSpreadCall_withBlockInlinableFunction_withPrecedingSpreadArg_isInlined()
func (c *InlineFunctionsTest) TestSpreadCall_withBlockInlinableFunction_withPrecedingSpreadArg_isInlined() {
	c.Test(
		`function foo(x) {
  qux(); // Something with a side-effect.
  return x;
}

bar(...zap(), foo(arg));
`,
		`var JSCompiler_temp_const$jscomp$v1 = bar;
var JSCompiler_temp_const$jscomp$v0 = [...zap()];
var JSCompiler_inline_result$jscomp$v2;
{
  var x$jscomp$inline_0 = arg;
  qux();
  JSCompiler_inline_result$jscomp$v2 = x$jscomp$inline_0;
}

JSCompiler_temp_const$jscomp$v1(
  ...JSCompiler_temp_const$jscomp$v0,
  JSCompiler_inline_result$jscomp$v2);
`)
}

// Java: InlineFunctionsTest.testSpreadCall_withBlockInlinableFunction_withPrecedingSpreadArg_isInlined_optChain()
func (c *InlineFunctionsTest) TestSpreadCall_withBlockInlinableFunction_withPrecedingSpreadArg_isInlined_optChain() {
	c.Test(
		`function foo(x) {
  qux?.(); // Something with a side-effect.
  return x;
}

bar?.(...zap?.(), foo?.(arg));
`,
		"let JSCompiler_temp_const$jscomp$v0;\n"+
			"// Optional call to `bar` gets decomposed to if- else.\n"+
			"if ((JSCompiler_temp_const$jscomp$v0 = bar) == null) {\n"+
			"  void 0;\n"+
			"} else {\n"+
			"  var JSCompiler_temp_const$jscomp$v2 = [...zap?.()];\n"+
			"  var JSCompiler_inline_result$jscomp$v3;\n"+
			"  {\n"+
			"    qux?.();\n"+
			"    JSCompiler_inline_result$jscomp$v3 = arg;\n"+
			"  }\n"+
			"  JSCompiler_temp_const$jscomp$v0(...JSCompiler_temp_const$jscomp$v2,\n"+
			" JSCompiler_inline_result$jscomp$v3);\n"+
			"}\n")
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: InlineFunctionsTest.methodCallIsDecomposed()
func (c *InlineFunctionsTest) TestMethodCallIsDecomposed() {
	c.Test(
		`function foo(x) {
  qux(); // Something with a side-effect.
  return x;
}

bar.method(foo(arg));
`,
		`var JSCompiler_temp_const$jscomp$v1 = bar;
var JSCompiler_temp_const$jscomp$v0 = JSCompiler_temp_const$jscomp$v1.method;
var JSCompiler_inline_result$jscomp$v2;
{
  var x$jscomp$inline_0 = arg;
  qux();
  JSCompiler_inline_result$jscomp$v2 = x$jscomp$inline_0;
 }
JSCompiler_temp_const$jscomp$v0.call(JSCompiler_temp_const$jscomp$v1,
 JSCompiler_inline_result$jscomp$v2);
`)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: InlineFunctionsTest.methodCallOnThisIsDecomposed()
func (c *InlineFunctionsTest) TestMethodCallOnThisIsDecomposed() {
	c.Test(
		`function foo(x) {
  qux(); // Something with a side-effect.
  return x;
}

class C {
  method() {
    this.privateMethod(foo(arg));
  }
}
`,
		`class C {
  method() {
    var JSCompiler_temp_const$jscomp$v0 = this.privateMethod
    var JSCompiler_inline_result$jscomp$v1;
    {
      var x$jscomp$inline_0 = arg;
      qux();
      JSCompiler_inline_result$jscomp$v1 = x$jscomp$inline_0;
    }
    JSCompiler_temp_const$jscomp$v0.call(
        this,
        JSCompiler_inline_result$jscomp$v1);
  }
}
`)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: InlineFunctionsTest.methodCallOnSuperIsDecomposed()
func (c *InlineFunctionsTest) TestMethodCallOnSuperIsDecomposed() {
	c.Test(
		`function foo(x) {
  qux(); // Something with a side-effect.
  return x;
}

class B {
  privateMethod(y) {}
}

class C extends B {
  method() {
    super.privateMethod(foo(arg));
  }
}
`,
		`class B {
  privateMethod(y) {}
}

class C extends B {
  method() {
    var JSCompiler_temp_const$jscomp$v0 = super.privateMethod;
    var JSCompiler_inline_result$jscomp$v1
    {
      var x$jscomp$inline_0 = arg;
      qux();
      JSCompiler_inline_result$jscomp$v1 = x$jscomp$inline_0;
    }
    JSCompiler_temp_const$jscomp$v0.call(
        this,
        JSCompiler_inline_result$jscomp$v1);
  }
}
`)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: InlineFunctionsTest.methodCallWithNestedCallIsInlined()
func (c *InlineFunctionsTest) TestMethodCallWithNestedCallIsInlined() {
	c.assumeStrictThis = true
	c.Test(
		`function foo(x) {
  qux(); // Something with a side-effect.
  return x;
}

function bar(x) {
  return x;
}
bar.call(null, foo(arg));
`,
		`var JSCompiler_inline_result$jscomp$v0;
{
  var x$jscomp$inline_0 = arg;
  qux();
  JSCompiler_inline_result$jscomp$v0 = x$jscomp$inline_0;
}
{
  JSCompiler_inline_result$jscomp$v0;
}
`)
}

// Java: InlineFunctionsTest.testGeneratorFunction()
func (c *InlineFunctionsTest) TestGeneratorFunction() {
	c.TestSame(
		`function* foo() {}
var bar = foo();
`)
	c.TestSame(
		`function* foo() {
  yield 'X';
  return 'Y';
}
var bar = foo();
`)
	c.TestSame(
		`function* foo() {
  yield 'X';
}
foo();
`)
	c.TestSame(
		`function* foo() {
  return 1;
}
foo();
`)
}

// Java: InlineFunctionsTest.testAsyncFunction()
func (c *InlineFunctionsTest) TestAsyncFunction() {
	c.TestSame(
		`async function foo() {}
foo().then(result => { alert(result); } );
`)
	c.TestSame(
		`async function foo() {
  return 'Y';
}
foo().then(result => { alert(result); } );
`)
	// We could possibly inline here since the return statement already contains a promise.
	c.TestSame(
		`async function foo() {
 return new Promise((resolve, reject) => { resolve('Success'); } );
}
foo().then(result => { alert(result); } );
`)
}

// Java: InlineFunctionsTest.testFunctionReferencingLetInNonGlobalBlock()
func (c *InlineFunctionsTest) TestFunctionReferencingLetInNonGlobalBlock() {
	c.TestSame(
		`if (true) {
  let value = 1;
  var g = function() {
    return value + 1;
  }
}
alert(g(10));
`)
}

// Java: InlineFunctionsTest.testNotInliningFunctionWithSuper()
func (c *InlineFunctionsTest) TestNotInliningFunctionWithSuper() {
	// Super field accessor arrow functions like this one are used for transpilation of some
	// features such as async functions and async generators. If inlined, it becomes a syntax error
	// as the inner function has no super.
	c.TestSame(
		`class A { m(){ return 1 } };
class B extends A {
  m() {
    const super$m = () => super.m;
    const jscomp$this = this;
    return function*() {
      yield super$m().call(jscomp$this);
    };
  }
}
`)
}

// Java: InlineFunctionsTest.testClassField()
func (c *InlineFunctionsTest) TestClassField() {
	c.Test(
		`class C {
  constructor() {
   /** @type {someType} */
    this.field = function f(){ function g() {} g() };
  }
}
C.staticField = function f(){ function g() {} g() };
`,
		`class C {
  constructor() {
    this.field = function f(){ void 0 };
  }
}
C.staticField = function f$jscomp$1() { void 0};
`)
	c.Test(
		`class C {
  field = function f(){ function g() {} g() };
  static staticField = function f(){ function g() {} g() };
}
`,
		`class C {
  field = function f(){ void 0 };
  static staticField = function f$jscomp$1() { void 0};
}
`)
	c.Test(
		`class C {
  constructor() {
   /** @return {number} */
    this.field = function() {return 1 };
  }
}
C.staticField = function() {return 1 };
console.log(new C().field());
console.log(C.staticField());
`,
		`class C {
  constructor() {
    this.field = function() {return 1 };
  }
}
C.staticField = function() {return 1 };
console.log(new C().field());
console.log(C.staticField());
`)
	c.TestSame(
		`class C {
  field = function() {return 1 };
  static staticField = function() {return 1 };
}
console.log(new C().field());
console.log(C.staticField());
`)
	c.TestSame(
		`function f({x}) {
  return x;
}
class Foo { x = 0; }
f(new Foo());
`)
}

// Java: InlineFunctionsTest.testClassField_doesNotInlineValueThatChanges()
func (c *InlineFunctionsTest) TestClassField_doesNotInlineValueThatChanges() {
	c.TestSame(
		"let nextVal = 0;\n" +
			"function getNextVal() {\n" +
			"  nextVal++;\n" +
			"  return nextVal;\n" +
			"}\n" +
			"\n" +
			"// Class expression.\n" +
			"const Uid = class {\n" +
			"  // Should not inline `nextVal` because it changes.\n" +
			"  val = getNextVal();\n" +
			"};\n")
	c.TestSame(
		"let nextVal = 0;\n" +
			"function getNextVal() {\n" +
			"  nextVal++;\n" +
			"  return nextVal;\n" +
			"}\n" +
			"\n" +
			"// Class declaration.\n" +
			"class Uid {\n" +
			"  // Should not inline `nextVal` because it changes.\n" +
			"  val = getNextVal();\n" +
			"}\n")
}

// Java: InlineFunctionsTest.testVariableUsedAsArgumentAndReassignedInFollowingArgument()
func (c *InlineFunctionsTest) TestVariableUsedAsArgumentAndReassignedInFollowingArgument() {
	// See https://github.com/google/closure-compiler/issues/4115.
	c.Test(
		`function f(value, begin, end) {
 return value.slice(begin, end);
}
window.g = function(a, b, c) {
  return f(a, b + 1, b = c);
};
`,
		`window.g = function(a, b, c) {
 var JSCompiler_inline_result$jscomp$v0;
{
  var value$jscomp$inline_0 = a;
  var begin$jscomp$inline_1 = b + 1;
  var end$jscomp$inline_2 = b = c;
  JSCompiler_inline_result$jscomp$v0 =
   value$jscomp$inline_0.slice(begin$jscomp$inline_1, end$jscomp$inline_2);
}
return JSCompiler_inline_result$jscomp$v0;
};
`)
}

// Java: InlineFunctionsTest.testVariableUsedAsArgumentAndModifiedInFollowingFunctionArgumentCall()
func (c *InlineFunctionsTest) TestVariableUsedAsArgumentAndModifiedInFollowingFunctionArgumentCall() {
	c.Test(
		`let text = 'abc'
let index = 1

// This function is necessary. Inlining the slice call will make temps for both args
function slice(begin, end) {
  return text.slice(begin, end);
}

function getEndIndex() {
  if (index)
    // This line is the only relevant line in this function - the index is being modified.
    // Everything else is artifical complexity to prevent the function from being inlined
    return ++index;

  return getEndIndex();
}

window.bug = slice(index - 1, getEndIndex());
`,
		`let text = 'abc';
let index = 1;
function getEndIndex() {
  if (index) {
    return ++index;
  }
  return getEndIndex();
}
var JSCompiler_temp_const$jscomp$v0 = window;
var JSCompiler_inline_result$jscomp$v1;
{
  var begin$jscomp$inline_0 = index - 1;
  var end$jscomp$inline_1 = getEndIndex();
  JSCompiler_inline_result$jscomp$v1 = text.slice(begin$jscomp$inline_0, end$jscomp$inline_1);
}
JSCompiler_temp_const$jscomp$v0.bug = JSCompiler_inline_result$jscomp$v1;
`)
}

// Java: InlineFunctionsTest.testVariableUsedAsArgumentAndReassignedInFollowingArgument2()
func (c *InlineFunctionsTest) TestVariableUsedAsArgumentAndReassignedInFollowingArgument2() {
	// See https://github.com/google/closure-compiler/issues/4115.
	c.Test(
		`function f(value, begin, end) {
 return value.slice(begin, end);
}
window.g = function(a, b, c) {
  return f(a, b.x + 1, b.x = c);
};
`,
		`window.g = function(a, b, c) {
    var JSCompiler_inline_result$jscomp$v0;
{
  var value$jscomp$inline_0 = a;
  var begin$jscomp$inline_1 = b.x + 1;
  var end$jscomp$inline_2 = b.x = c;
  JSCompiler_inline_result$jscomp$v0 =
      value$jscomp$inline_0.slice(begin$jscomp$inline_1, end$jscomp$inline_2);
}
return JSCompiler_inline_result$jscomp$v0;
};
`)
}

// Java: InlineFunctionsTest.testDifferentPropertiesAsArguments()
func (c *InlineFunctionsTest) TestDifferentPropertiesAsArguments() {
	// See https://github.com/google/closure-compiler/issues/4115.
	c.Test(
		`function f(value, begin, end) {
 return value.slice(begin, end);
}
window.g = function(a, b, c) {
  return f(a, b.x + 1, b.y = c);
};
`,
		`window.g = function(a, b, c) {
  var JSCompiler_inline_result$jscomp$v0;
  {
    var value$jscomp$inline_0 = a;
    var begin$jscomp$inline_1 = b.x + 1;
    var end$jscomp$inline_2 = b.y = c;
    JSCompiler_inline_result$jscomp$v0 = value$jscomp$inline_0.slice(begin$jscomp$inline_1, end$jscomp$inline_2);
  }
  return JSCompiler_inline_result$jscomp$v0;
};
`)
}

// Java: InlineFunctionsTest.testVariableUsedAsArgumentAndReassignedInPreviousArgument()
func (c *InlineFunctionsTest) TestVariableUsedAsArgumentAndReassignedInPreviousArgument() {
	c.Test(
		`function f(value, begin, end) {
 return value.slice(begin, end);
}
window.g = function(a, b, c) {
  return f(a, b = c, b + 1);
};
`,
		"window.g = function(a, b, c) {\n"+
			"  var JSCompiler_inline_result$jscomp$v0;\n"+
			"  {\n"+
			"    // Evaluates all args until `b=c` in order\n"+
			"    var value$jscomp$inline_0 = a;\n"+
			"    var begin$jscomp$inline_1 = b = c;\n"+
			"    JSCompiler_inline_result$jscomp$v0 = value$jscomp$inline_0.slice(begin$jscomp$inline_1, b + 1);\n"+
			"  }\n"+
			"  return JSCompiler_inline_result$jscomp$v0;\n"+
			"};\n")
}

// Java: InlineFunctionsTest.testVariableUsedAsArgumentMultipleTimes()
func (c *InlineFunctionsTest) TestVariableUsedAsArgumentMultipleTimes() {
	c.Test(
		`function f(value, begin, end) {
 return value.slice(begin, end);
}
window.g = function(a, b, c) {
// No temp needed as none of the arguments mutate state.
  return f(a, b, b + 1);
};
`,
		`window.g = function(a, b, c) {
  return a.slice(b, b + 1);
};
`)
}

// Java: InlineFunctionsTest.testVariableUsedAsArgumentMultipleTimes2()
func (c *InlineFunctionsTest) TestVariableUsedAsArgumentMultipleTimes2() {
	c.Test(
		`function f(value, begin, end) {
 return value.slice(begin, end);
}
window.g = function(a, b, c) {
  return f(a, ++b, b++);
};
`,
		`window.g = function(a, b, c) {
  var JSCompiler_inline_result$jscomp$v0;
  {
    var value$jscomp$inline_0 = a;
    var begin$jscomp$inline_1 = ++b;
    var end$jscomp$inline_2 = b++;
    JSCompiler_inline_result$jscomp$v0 = value$jscomp$inline_0.slice(begin$jscomp$inline_1, end$jscomp$inline_2);
  }
  return JSCompiler_inline_result$jscomp$v0;
};
`)
}

// Java: InlineFunctionsTest.testVariableUsedAsArgumentMultipleTimes3()
func (c *InlineFunctionsTest) TestVariableUsedAsArgumentMultipleTimes3() {
	c.Test(
		`function add(a,b) {
 return a+b;
}
window.g = function(value) {
  return add(++value, value++);
};
`,
		`window.g = function(value) {
  var JSCompiler_inline_result$jscomp$v0;
  {
    var a$jscomp$inline_0 = ++value;
    var b$jscomp$inline_1 = value++;
    JSCompiler_inline_result$jscomp$v0 = a$jscomp$inline_0 + b$jscomp$inline_1;
  }
  return JSCompiler_inline_result$jscomp$v0;
};
`)
}

// Java: InlineFunctionsTest.testVariableUsedAsArgumentMultipleTimes4()
func (c *InlineFunctionsTest) TestVariableUsedAsArgumentMultipleTimes4() {
	c.Test(
		`function add(a,b) {
 return a+b;
}
window.g = function(value) {
  return add(++value, ++value);
};
`,
		`window.g = function(value) {
  var JSCompiler_inline_result$jscomp$v0;
  {
    var a$jscomp$inline_0 = ++value;
    var b$jscomp$inline_1 = ++value;
    JSCompiler_inline_result$jscomp$v0 = a$jscomp$inline_0 + b$jscomp$inline_1;
  }
  return JSCompiler_inline_result$jscomp$v0;
};
`)
}

// Java: InlineFunctionsTest.testVariableUsedAsArgumentMultipleTimes5()
func (c *InlineFunctionsTest) TestVariableUsedAsArgumentMultipleTimes5() {
	c.Test(
		`function add(value, a,b) {
 return value+a+b;
}
window.g = function(value) {
  return add(value, ++value, ++value);
};
`,
		`window.g = function(value$jscomp$1) {
  var JSCompiler_inline_result$jscomp$v0;
  {
    var value$jscomp$inline_0 = value$jscomp$1;
    var a$jscomp$inline_1 = ++value$jscomp$1;
    var b$jscomp$inline_2 = ++value$jscomp$1;
    JSCompiler_inline_result$jscomp$v0 = value$jscomp$inline_0 + a$jscomp$inline_1 + b$jscomp$inline_2;
  }
  return JSCompiler_inline_result$jscomp$v0;
};
`)
}

// Java: InlineFunctionsTest.testVariableUsedAsArgumentMultipleTimes6()
func (c *InlineFunctionsTest) TestVariableUsedAsArgumentMultipleTimes6() {
	c.Test(
		`value = 10;
function add(a,b) {
 return (value = 0)+a+b;
}
window.g = function() {
  return add(++value, ++value);
};
`,
		`value = 10;
window.g = function() {
  var JSCompiler_inline_result$jscomp$v0;
  {
    var a$jscomp$inline_0 = ++value;
    var b$jscomp$inline_1 = ++value;
    JSCompiler_inline_result$jscomp$v0 = (value = 0 ) + a$jscomp$inline_0 + b$jscomp$inline_1;
  }
  return JSCompiler_inline_result$jscomp$v0;
};
`)
}

// inlineFunctionsTest_test6671158Tester is the anonymous subclass of
// InlineFunctionsTest in test6671158, which does not enable InferConsts.
type inlineFunctionsTest_test6671158Tester struct {
	*InlineFunctionsTest
}

// Java: InlineFunctionsTest.test6671158().new InlineFunctionsTest() {...}
func newInlineFunctionsTest_test6671158Tester() *inlineFunctionsTest_test6671158Tester {
	t := &inlineFunctionsTest_test6671158Tester{NewInlineFunctionsTest()}
	t.self = t
	return t
}

// Java: InlineFunctionsTest.test6671158().new InlineFunctionsTest() {...}.maybeEnableInferConsts()
func (t *inlineFunctionsTest_test6671158Tester) MaybeEnableInferConsts() {}
