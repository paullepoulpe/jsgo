// Port of com/google/javascript/jscomp/CoalesceVariableNamesTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/CoalesceVariableNamesTest.java)

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
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"testing"
)

type CoalesceVariableNamesTest struct {
	CompilerTestCaseBase
	usePseudoName bool
}

// Java: CoalesceVariableNamesTest()
func NewCoalesceVariableNamesTest() *CoalesceVariableNamesTest {
	c := &CoalesceVariableNamesTest{}
	c.InitCompilerTestCaseBase(c)
	c.usePseudoName = false
	return c
}

func TestCoalesceVariableNamesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCoalesceVariableNamesTest() })
}

// The spacing in this file is not exactly standard but it greatly helps
// picking out which variable names are merged.
// Java: CoalesceVariableNamesTest.setUp()
func (c *CoalesceVariableNamesTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.usePseudoName = false
}

// Java: CoalesceVariableNamesTest.getProcessor(Compiler)
func (c *CoalesceVariableNamesTest) GetProcessor(compiler *Compiler) CompilerPass {
	return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
		// enableNormalize would require output of CoalesceVariableNames to be normalized,
		// so we just manually normalize the input instead.
		normalize := Normalize_CreateNormalizeForOptimizations(compiler)
		normalize.Process(externs, root)
		NewCoalesceVariableNames(compiler, c.usePseudoName).Process(externs, root)
	})
}

// Java: CoalesceVariableNamesTest.testSimple()
func (c *CoalesceVariableNamesTest) TestSimple() {
	c.InFunctionWithExpected("var x; var y; x=1; x; y=1; y; return y", "var x;        x=1; x; x=1; x; return x")

	c.InFunctionWithExpected("var x,y; x=1; x; y=1; y", "var x  ; x=1; x; x=1; x")

	c.InFunction("var x; var y; x=1; y=2; y; x")

	c.InFunctionWithExpected("y=0; var x, y; y; x=0; x", "y=0; var y   ; y; y=0;y")

	c.InFunctionWithExpected("var x,y; x=1; y=x; y", "var x  ; x=1; x=x; x")

	c.InFunctionWithExpected("var x,y; x=1; y=x+1; y", "var x  ; x=1; x=x+1; x")

	c.InFunctionWithExpected("x=1; x; y=2; y; var x; var y", "x=1; x; x=2; x; var x")

	c.InFunctionWithExpected("var x=1; var y=x+1; return y", "var x=1;     x=x+1; return x")

	c.InFunction("var x=1; var y=0; x = x + 1; y")

	c.InFunctionWithExpected("var x=1; x+=1;     var y=0; y", "var x=1; x = x + 1;     x=0; x")

	c.InFunctionWithExpected(
		"var x=1; foo(bar(x+=1));     var y=0; y", "var x=1; foo(bar(x = x + 1));    x=0; x")

	c.InFunction("var y; var x=1; f(x = x + 1, y)")

	c.InFunction("var x; var y; y = y + 1, y, x = 1; x")
}

// Java: CoalesceVariableNamesTest.testCoaleseLetAndConst()
func (c *CoalesceVariableNamesTest) TestCoaleseLetAndConst() {
	c.InFunctionWithExpected(
		"let x; const y = 1; x = y + 1; return x",
		// `let` must become `var`, because we might be coalescing
		// variables declared in different blocks.
		// See testLetAndConstDifferentBlock()
		"var x;       x = 1; x = x + 1; return x")
}

// Java: CoalesceVariableNamesTest.testLetAndConstDifferentBlock()
func (c *CoalesceVariableNamesTest) TestLetAndConstDifferentBlock() {
	c.InFunctionWithExpected(
		"if(1) { const x = 0; x } else { let y = 0; y }",
		"if(1) {   var x = 0; x } else {     x = 0; x }")
}

// Java: CoalesceVariableNamesTest.testCoalesceLetRequiresInitWithinALoop()
func (c *CoalesceVariableNamesTest) TestCoalesceLetRequiresInitWithinALoop() {
	c.InFunctionWithExpected(
		"for (let i = 0; i < 3; ++i) {\n"+
			"  let something;\n"+
			"  if (i == 0) {\n"+
			"    const x = 'hi';\n"+
			"    alert(x);\n"+
			"    something = x + ' there';\n"+
			"  }\n"+
			"  alert(something);\n"+
			"}\n",
		"for (let i = 0; i < 3; ++i) {\n"+
			"// we must initialize `something` on each loop iteration\n"+
			"  var something = void 0;\n"+
			"  if (i == 0) {\n"+
			"    something = 'hi';\n"+
			"    alert(something); // always alerts 'hi'\n"+
			"    something = something + ' there';\n"+
			"  }\n"+
			"  alert(something);\n"+
			"}\n")
}

// Java: CoalesceVariableNamesTest.testCoalesceVariableNames_doesNotWronglyUndeclareOuterName()
func (c *CoalesceVariableNamesTest) TestCoalesceVariableNames_doesNotWronglyUndeclareOuterName() {
	c.DisableValidateAstChangeMarking()
	c.Test(
		"const [{localKey:localKey$jscomp$2} = {localKey:void 0}] =\n"+
			" entries.filter(a => { var b; [, b] = a; return b === c; });\n"+
			"alert(localKey$jscomp$2);\n",
		// The CoalesceVariableNames pass merges the names `a` and `b` in the arrow function, but in
		// doing so does not delete the const keyword which declared the outer name
		// `localKey$jscomp$2`
		"const [{localKey:localKey$jscomp$2} = {localKey:void 0}] =\n"+
			"entries.filter(a => { [, a] = a; return a === c; })\n"+
			"alert(localKey$jscomp$2);\n")
}

// Java: CoalesceVariableNamesTest.testCoalesceVariableNames_worksInsideExpressionResults()
func (c *CoalesceVariableNamesTest) TestCoalesceVariableNames_worksInsideExpressionResults() {
	c.InFunctionWithExpected(
		"a.b = function() {\n"+
			"  const [{localKey:localKey$jscomp$2} = {localKey:void 0}] =\n"+
			"    entries.filter(a => {\n"+
			"      var b;\n"+
			"      [, b] = a;\n"+
			"      return b === c;\n"+
			"    });\n"+
			"  alert(localKey$jscomp$2);\n"+
			"}\n",
		"a.b = function() {\n"+
			"  const [{localKey:localKey$jscomp$2} = {localKey:void 0}] =\n"+
			"    entries.filter(a => {\n"+
			"      [, a] = a;\n"+
			"      return a === c;\n"+
			"    })\n"+
			"  alert(localKey$jscomp$2);\n"+
			"}\n")
}

// Java: CoalesceVariableNamesTest.testMergeThreeVarNames()
func (c *CoalesceVariableNamesTest) TestMergeThreeVarNames() {
	c.InFunctionWithExpected("var x,y,z; x=1; x; y=1; y; z=1; z", "var x    ; x=1; x; x=1; x; x=1; x")
}

// Java: CoalesceVariableNamesTest.testDifferentBlock()
func (c *CoalesceVariableNamesTest) TestDifferentBlock() {
	c.InFunctionWithExpected(
		"if(1) { var x = 0; x } else { var y = 0; y }",
		"if(1) { var x = 0; x } else {     x = 0; x }")
}

// Java: CoalesceVariableNamesTest.testLoops()
func (c *CoalesceVariableNamesTest) TestLoops() {
	c.InFunction("var x; for ( ; 1; ) { x; x = 1; var y = 1; y }")

	c.InFunctionWithExpected(
		"var y = 1; y; for ( ; 1; ) { var x = 1; x }",
		"var y = 1; y; for ( ; 1; ) {     y = 1; y }")
}

// Java: CoalesceVariableNamesTest.testEscaped()
func (c *CoalesceVariableNamesTest) TestEscaped() {
	c.InFunction("function f() { x } var x = 1; x; var y = 0; y; f()")
}

// Java: CoalesceVariableNamesTest.testFor()
func (c *CoalesceVariableNamesTest) TestFor() {
	c.InFunctionWithExpected("var x = 1; x; for (;;) var y; y = 1; y", "var x = 1; x; for (;;)      ; x = 1; x")
}

// Java: CoalesceVariableNamesTest.testForIn()
func (c *CoalesceVariableNamesTest) TestForIn() {
	// We lose some precision here, unless we have "branched-backward-dataflow".
	c.InFunction("var x = 1; var k; x; var y; for ( y in k ) y")

	c.InFunctionWithExpected(
		"var x = 1,     k; x; y = 1; for (var y in k) { y }",
		"var x = 1; var k; x; x = 1; for (    x in k) { x }")

	c.InFunction("function f(param){ var foo; for([foo] in arr); param }")
}

// Java: CoalesceVariableNamesTest.testForLoopCoalesceWithFollowingCode()
func (c *CoalesceVariableNamesTest) TestForLoopCoalesceWithFollowingCode() {
	c.InFunctionWithExpected(
		"for (;;) { const a = 3; } const y = 1; y;", //
		"for (;;) { var   a = 3; }       a = 1; a;")
	c.InFunctionWithExpected(
		"for (let a = 3;;) { a; } const y = 1; y;", //
		"for (var a = 3;;) { a; }       a = 1; a;")
	c.InFunctionWithExpected(
		"for (const x in k) { x; } const y = 1; y;", //
		"for (var   x in k) { x; }       x = 1; x;")
	c.InFunctionWithExpected(
		"for (let x in k) { x; } const y = 1; y;", //
		"for (var x in k) { x; }       x = 1; x;")
	c.InFunctionWithExpected(
		"for (const x of k) { x; } const y = 1; y;", //
		"for (var   x of k) { x; }       x = 1; x;")
	c.InFunctionWithExpected(
		"for (let x of k) { x; } const y = 1; y;", //
		"for (var x of k) { x; }       x = 1; x;")
}

// Java: CoalesceVariableNamesTest.testForOf()
func (c *CoalesceVariableNamesTest) TestForOf() {
	// We lose some precision here, unless we have "branched-backward-dataflow".
	c.InFunction("var x = 1; var k; x; var y; for ( y of k ) y")

	c.InFunctionWithExpected(
		"var x = 1,     k; x; y = 1; for (var y of k) { y }",
		"var x = 1; var k; x; x = 1; for (    x of k) { x }")

	c.InFunction("function f(param){ var foo; for([foo] of arr); param }")
}

// Java: CoalesceVariableNamesTest.testLoopInductionVar()
func (c *CoalesceVariableNamesTest) TestLoopInductionVar() {
	c.InFunctionWithExpected(
		"for(var x = 0; x < 10; x++){}\n"+
			"for(var y = 0; y < 10; y++){}\n"+
			"for(var z = 0; z < 10; z++){}\n",
		"var x=0;\n"+
			"for(;x<10;x++);\n"+
			"x=0;\n"+
			"for(;x<10;x++);\n"+
			"x=0;\n"+
			"for(;x<10;x++) {}\n")

	c.InFunctionWithExpected(
		"for(var x = 0; x < 10; x++){z} for(var y = 0, z = 0; y < 10; y++){z}",
		"var x=0;for(;x<10;x++)z;x=0;var z=0;for(;x<10;x++)z")

	c.InFunctionWithExpected("var x = 1; x; for (var y; y=1; ) {y}", "var x = 1; x; for ( ; x=1; ) {x}")

	c.InFunctionWithExpected("var x = 1; x; y = 1; while(y) var y; y", "var x = 1; x; x = 1; for (;x;); x")

	// It's not the job of the coalesce variables pass to remove unused labels
	c.InFunctionWithExpected("var x = 1; x; f:var y; y=1", "var x = 1; x; f:{} x=1")
}

// Java: CoalesceVariableNamesTest.testSwitchCase()
func (c *CoalesceVariableNamesTest) TestSwitchCase() {
	c.InFunctionWithExpected(
		"var x = 1; switch(x) { case 1: var y; case 2: } y = 1; y",
		"var x = 1; switch(x) { case 1:        case 2: } x = 1; x")
}

// Java: CoalesceVariableNamesTest.testDuplicatedVar()
func (c *CoalesceVariableNamesTest) TestDuplicatedVar() {
	// Is there a shorter version without multiple declarations?
	c.InFunctionWithExpected(
		"z = 1; var x = 0; x; z; var y = 2,     z = 1; y; z;",
		"z = 1; var x = 0; x; z;     x = 2; var z = 1; x; z;")
}

// Java: CoalesceVariableNamesTest.testTryCatch()
func (c *CoalesceVariableNamesTest) TestTryCatch() {
	c.InFunction("try {} catch (e) { } var x = 4; x;")
	c.InFunction("var x = 4; x; try {} catch (e) { }")
}

// Java: CoalesceVariableNamesTest.testOptionalChainingUnsafeCoalesce()
func (c *CoalesceVariableNamesTest) TestOptionalChainingUnsafeCoalesce() {
	c.InFunction(
		`var x = 1;
try {
  cb?.();
  x = 2;
} catch (e) {
  var y = 3;
  use(y);
}
use(x);
`)
}

// Java: CoalesceVariableNamesTest.testDeadAssignment()
func (c *CoalesceVariableNamesTest) TestDeadAssignment() {
	c.InFunction("var x = 6; var y; y = 4 ; x")
	c.InFunction("var y = 3; y = y + 4; x")
	c.InFunction("var y = 3; y = y + 1; x")
	c.InFunction("y = 3; var x; var y = 1 ; x")
}

// Java: CoalesceVariableNamesTest.testParameter()
func (c *CoalesceVariableNamesTest) TestParameter() {
	c.Test("function FUNC(param) {var x = 0; x}", "function FUNC(param) {param = 0; param}")
}

// Java: CoalesceVariableNamesTest.testParameter2()
func (c *CoalesceVariableNamesTest) TestParameter2() {
	// Make sure two formal parameter name never merges.
	c.TestSame("function FUNC(x,y) {x = 0; x; y = 0; y}")
	c.TestSame("function FUNC(x,y,z) {x = 0; x; y = 0; z = 0; z}")
}

// Java: CoalesceVariableNamesTest.testParameter3()
func (c *CoalesceVariableNamesTest) TestParameter3() {
	// Make sure the formal parameter declaration is consider a def.
	c.TestSame("function FUNC(x) {var y; y = 0; x; y}")
}

// Java: CoalesceVariableNamesTest.testParameter4a()
func (c *CoalesceVariableNamesTest) TestParameter4a() {
	// Make sure that we do not merge two-arg functions because of the
	// IE sort bug (see comments in computeEscaped)
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT3)
	c.Test(
		"function FUNC(x, y) {var a,b; y; a=0; a; x; b=0; b}",
		"function FUNC(x, y) {var a;   y; a=0; a; x; a=0; a}")
}

// Java: CoalesceVariableNamesTest.testParameter4b()
func (c *CoalesceVariableNamesTest) TestParameter4b() {
	// Go ahead and merge, if language-out is ES5, because that means IE8 does
	// not need to be supported.
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.Test(
		"function FUNC(x, y) {var a,b; y; a=0; a; x; b=0; b}",
		"function FUNC(x, y) {         y; y=0; y; x; x=0; x}")
}

// Java: CoalesceVariableNamesTest.testParameter5()
func (c *CoalesceVariableNamesTest) TestParameter5() {
	// Merge parameters
	c.Test(
		"function FUNC(x, y, z) {var a,b; y; a=0; a; x; b=0; b}",
		"function FUNC(x, y, z) {         y; y=0; y; x; x=0; x}")
}

// Java: CoalesceVariableNamesTest.testLiveRangeChangeWithinCfgNode()
func (c *CoalesceVariableNamesTest) TestLiveRangeChangeWithinCfgNode() {
	c.InFunction("var x; var y; x = 1, y = 2, y, x")
	c.InFunction("var x; var y; x = 1,x; y")

	// We lose some precisions within the node itself.
	c.InFunction("var x; var y; y = 1, y, x = 1; x")

	c.InFunctionWithExpected("var x; var y; y = 1; y, x = 1; x", "var x; x = 1; x, x = 1; x")

	c.InFunction("var x; var y; y = 1, x = 1, x, y = y + 1, y")

	c.InFunction("var x; var y; y = 1, x = 1, x, y = y + 1, y")
}

// Java: CoalesceVariableNamesTest.testUsePseudoName_forWithUninitializedLetInside()
func (c *CoalesceVariableNamesTest) TestUsePseudoName_forWithUninitializedLetInside() {
	c.usePseudoName = false
	src :=
		"for( var i = 0;i < 10; i++) {\n" +
			"  let x;\n" +
			"  if(true) {x=3} else { x = undefined;}\n" +
			"  i;\n" +
			"}\n" +
			"const y = new Set();\n" +
			"i; y;\n"
	expected :=
		"var i = 0;\n" +
			"for(; i < 10; i++) {\n" +
			"  var x = void 0;\n" +
			"  if(true) {x=3} else { x = undefined;}\n" +
			"  i;\n" +
			"}\n" +
			"x = new Set();\n" +
			"i; x;\n"
	c.InFunctionWithExpected(src, expected)

	// make sure we generate the correct code with `usePseudoName = true`
	c.usePseudoName = true
	expectedPseudoNames :=
		"var i = 0;\n" +
			"for(; i < 10; i++) {\n" +
			"  var x_y = void 0; // preserves \"reset\" before every iteration\n" +
			"  if(true) {x_y=3} else { x_y = undefined;}\n" +
			"  i;\n" +
			"}\n" +
			"x_y = new Set();\n" +
			"i; x_y;\n"
	c.InFunctionWithExpected(src, expectedPseudoNames)
}

// Java: CoalesceVariableNamesTest.test_doesNotCoalesceOverlappingLiveRange()
func (c *CoalesceVariableNamesTest) Test_doesNotCoalesceOverlappingLiveRange() {
	c.usePseudoName = false
	src :=
		"for( var i = 0;i < 10; i++) {\n" +
			"  let x = i;\n" +
			"  let y;\n" +
			"  x;\n" +
			"  y;\n" +
			"}\n"
	expected :=
		"var i = 0;\n" +
			"for(; i < 10; i++) {\n" +
			"  let x = i;\n" +
			"  let y;\n" +
			"  x;\n" +
			"  y;\n" +
			"}\n"
	c.InFunctionWithExpected(src, expected)

	src =
		"for( var i = 0;i < 10; i++) {\n" +
			"  let x = i;\n" +
			"  let y;\n" +
			"  y;\n" +
			"  x;\n" +
			"}\n"
	expected =
		"var i = 0;\n" +
			"for(; i < 10; i++) {\n" +
			"  let x = i;\n" +
			"  let y;\n" +
			"  y;\n" +
			"  x;\n" +
			"}\n"
	c.InFunctionWithExpected(src, expected)

	c.usePseudoName = true // same for usePseudoName = true code path
	c.InFunctionWithExpected(src, expected)
}

// We could potentially coaleasce here but we miss.
// Java: CoalesceVariableNamesTest.test_doesNotCoalesceNonOverlappingLiveRange()
func (c *CoalesceVariableNamesTest) Test_doesNotCoalesceNonOverlappingLiveRange() {
	c.usePseudoName = false
	src :=
		"for( var i = 0;i < 10; i++) {\n" +
			"  let x;\n" +
			"  x;\n" +
			"  let y;\n" +
			"  y;\n" +
			"}\n"
	expected :=
		"var i = 0;\n" +
			"for(; i < 10; i++) {\n" +
			"  let x;\n" +
			"  x;\n" +
			"  let y;\n" +
			"  y;\n" +
			"}\n"
	c.InFunctionWithExpected(src, expected)
	c.usePseudoName = true // same for usePseudoName = true code path
	c.InFunctionWithExpected(src, expected)
}

// Java: CoalesceVariableNamesTest.test_doesNotUsePseudoName_coalesces_deletesUninitializedDeclInLoop()
func (c *CoalesceVariableNamesTest) Test_doesNotUsePseudoName_coalesces_deletesUninitializedDeclInLoop() {
	c.usePseudoName = false
	src :=
		"for( var i = 0;i < 10; i++) {\n" +
			"  const x = i;\n" +
			"  let y; // safe to delete when coalesced\n" +
			"  x;\n" +
			"  y = x;\n" +
			"}\n"
	expected :=
		"var i = 0;\n" +
			"for(; i < 10; i++) {\n" +
			"  var x = i;\n" +
			"  x;\n" +
			"  x = x;\n" +
			"}\n"
	c.InFunctionWithExpected(src, expected)
}

// Java: CoalesceVariableNamesTest.test_usePseudoName_coalesces_deletesUninitializedDeclInLoop()
func (c *CoalesceVariableNamesTest) Test_usePseudoName_coalesces_deletesUninitializedDeclInLoop() {
	c.usePseudoName = true
	src :=
		"for( var i = 0;i < 10; i++) {\n" +
			"  const x = i;\n" +
			"  let y; // safe to delete when coalesced\n" +
			"  x;\n" +
			"  y=x;\n" +
			"}\n"
	expected :=
		"var i = 0;\n" +
			"for(; i < 10; i++) {\n" +
			"  var x_y = i;\n" +
			"  x_y;\n" +
			"  x_y = x_y;\n" +
			"}\n"
	c.InFunctionWithExpected(src, expected)
}

// Java: CoalesceVariableNamesTest.test_doesNotUsePseudoName_coalesces_resetsUninitializedDeclInLoopWhenLive()
func (c *CoalesceVariableNamesTest) Test_doesNotUsePseudoName_coalesces_resetsUninitializedDeclInLoopWhenLive() {
	c.usePseudoName = false
	src := `var y;
while (someExtern) {
  let x;
  x;
  x = 1;
}
y = 2;
y;
`
	expected := `var y;
for (; someExtern;) {
  y = void 0;
  y;
  y = 1;
}
y = 2;
y;
`
	c.InFunctionWithExpected(src, expected)
}

// Java: CoalesceVariableNamesTest.test_usePseudoName_coalesces_resetsUninitializedDeclInLoopWhenLive()
func (c *CoalesceVariableNamesTest) Test_usePseudoName_coalesces_resetsUninitializedDeclInLoopWhenLive() {
	c.usePseudoName = true
	src := `var y;
while (someExtern) {
  let x;
  x;
  x = 1;
}
y = 2;
y;
`
	expected := `var x_y;
for (; someExtern;) {
  x_y = void 0;
  x_y;
  x_y = 1;
}
x_y = 2;
x_y;
`
	c.InFunctionWithExpected(src, expected)
}

// Java: CoalesceVariableNamesTest.test_doesNotUsePseudoName_coalesces_resetsUninitializedDeclInDoWhileWhenLive()
func (c *CoalesceVariableNamesTest) Test_doesNotUsePseudoName_coalesces_resetsUninitializedDeclInDoWhileWhenLive() {
	c.usePseudoName = false
	src := `var y;
do {
  let x;
  x;
  x = 1;
} while (someExtern);
y = 2;
y;
`
	expected := `var y;
do {
  y = void 0;
  y;
  y = 1;
} while (someExtern);
y = 2;
y;
`
	c.InFunctionWithExpected(src, expected)
}

// Java: CoalesceVariableNamesTest.testUsePseudoName_forWithUninitializedLetInside2()
func (c *CoalesceVariableNamesTest) TestUsePseudoName_forWithUninitializedLetInside2() {
	c.usePseudoName = true
	src :=
		"for( var i = 0;i < 10; i++) {\n" +
			"  let x;\n" +
			"  if(someExtern) {x=3}\n" +
			"  i; x;\n" +
			"}\n" +
			"const y = new Set();\n" +
			"i; y;\n"
	expected :=
		"var i = 0;\n" +
			"for(; i < 10; i++) {\n" +
			"  var x_y = void 0; // preserves \"reset\" before every iteration\n" +
			"  if(someExtern) {x_y=3}\n" +
			"  i; x_y;\n" +
			"}\n" +
			"x_y = new Set();\n" +
			"i; x_y;\n"
	c.InFunctionWithExpected(src, expected)
}

// Java: CoalesceVariableNamesTest.testUsePseudoName_forWithUninitializedLetInside_onlyTwoNodesInLiveness()
func (c *CoalesceVariableNamesTest) TestUsePseudoName_forWithUninitializedLetInside_onlyTwoNodesInLiveness() {
	c.usePseudoName = false
	c.InFunctionWithExpected(
		"for(;;) { let x; x;}  const y = new Set(); y;",
		"for(;;) { var x = void 0; x; } x = new Set(); x;")

	c.usePseudoName = true
	c.InFunctionWithExpected(
		"for(;;) { let x; x;}  const y = new Set(); y;",
		"for(;;) { var x_y = void 0; x_y; } x_y = new Set(); x_y;")
}

// Java: CoalesceVariableNamesTest.testLiveRangeChangeWithinCfgNode2()
func (c *CoalesceVariableNamesTest) TestLiveRangeChangeWithinCfgNode2() {
	c.InFunction("var x; var y; var a; var b; y = 1, a = 1, y, a, x = 1, b = 1; x; b")

	c.InFunctionWithExpected(
		"var x; var y; var a; var b; y = 1, a = 1, y, a, x = 1; x; b = 1; b",
		"var x; var y; var a;        y = 1, a = 1, y, a, x = 1; x; x = 1; x")

	c.InFunctionWithExpected(
		"var x; var y; var a; var b; y = 1, a = 1, y, x = 1; a; x; b = 1; b",
		"var x; var y; var a;        y = 1, a = 1, y, x = 1; a; x; x = 1; x")
}

// Java: CoalesceVariableNamesTest.testFunctionNameReuse()
func (c *CoalesceVariableNamesTest) TestFunctionNameReuse() {
	c.InFunction("function x() {}; x(); var y = 1; y")

	c.InFunction("function x() { } x(); var y = 1; y")

	c.InFunction("function x() { } x(); var y = 1; y")

	// Can't merge because of possible escape.
	c.InFunction("function x() {return x}; x(); var y = 1; y")

	c.InFunction("function x() {} var y = 1; y; x")

	c.InFunction("function x() { } var y = 1; y; x")

	c.InFunction("function x() { } var y = 1; y; x = 1; x")

	c.InFunction("function x() {} var y = 1; y; x = 1; x")
}

// Java: CoalesceVariableNamesTest.testCoalesceVarinStaticBlock()
func (c *CoalesceVariableNamesTest) TestCoalesceVarinStaticBlock() {
	c.Test(
		"function foo() {\n"+
			"   var x = 0;\n"+
			"   var y = 2;\n"+
			"   const C = class{\n"+
			"      static{\n"+
			"       print(x);\n"+
			"      }\n"+
			"   }\n"+
			"}\n",
		"function foo() {\n"+
			"   var x = 0;\n"+
			"   var y = 2;\n"+
			"   y = class{\n"+
			"     static{\n"+
			"      print(x);\n"+
			"     }\n"+
			"   }\n"+
			"}\n")

	c.Test(
		"function foo() {\n"+
			"   var x = 0;\n"+
			"   var y = 2;\n"+
			"   const C = class {\n"+
			"      function(){\n"+
			"       print(x);\n"+
			"      }\n"+
			"   }\n"+
			"}\n",
		"function foo() {\n"+
			"   var x = 0;\n"+
			"   var y = 2;\n"+
			"   y = class {\n"+
			"     function(){\n"+
			"      print(x);\n"+
			"     }\n"+
			"  }\n"+
			"}\n")
}

// Java: CoalesceVariableNamesTest.testBug65688660()
func (c *CoalesceVariableNamesTest) TestBug65688660() {
	c.Test(
		"function f(param) {\n"+
			"  if (true) {\n"+
			"    const b1 = [];\n"+
			"    for (const [key, value] of []) {}\n"+
			"  }\n"+
			"  if (true) {\n"+
			"    const b2 = [];\n"+
			"    for (const kv of []) {\n"+
			"      const key2 = kv.key;\n"+
			"    }\n"+
			"  }\n"+
			"}\n",
		"function f(param) {\n"+
			"  if (true) {\n"+
			"    param = [];\n"+
			"    for (const [key, value] of []) {}\n"+
			"  }\n"+
			"  if (true) {\n"+
			"    param = [];\n"+
			"    for (const kv of []) {\n"+
			"      param = kv.key;\n"+
			"    }\n"+
			"  }\n"+
			"}\n")
}

// Java: CoalesceVariableNamesTest.testBug1401831()
func (c *CoalesceVariableNamesTest) TestBug1401831() {
	// Verify that we don't wrongly merge "opt_a2" and "i" without considering
	// arguments[0] aliasing it.
	src :=
		"function f(opt_a2){\n" +
			"  var buffer;\n" +
			"  if(opt_a2){\n" +
			"    var i=0;\n" +
			"    for(;i<arguments.length;i++)\n" +
			"      buffer=buffer+arguments[i];\n" +
			"  }\n" +
			"  return buffer;\n" +
			"}\n"
	c.TestSame(src)
}

// Code inside a class is automatically in strict mode, so duplicated parameter names are not
// allowed.
// Java: CoalesceVariableNamesTest.testBug64898400()
func (c *CoalesceVariableNamesTest) TestBug64898400() {
	c.TestSame("class C { f(a, b, c) {} }")
	c.TestSame("class C { f(a, b=0, c=0) {} }")
}

// Java: CoalesceVariableNamesTest.testDontCoalesceClassDeclarationsWithConstDeclaration()
func (c *CoalesceVariableNamesTest) TestDontCoalesceClassDeclarationsWithConstDeclaration() {
	c.TestSame(
		"function f() { // preserve newline\n" +
			"  class A {}\n" +
			"  const b = {};\n" +
			"  return b;\n" +
			"}\n")
}

// Java: CoalesceVariableNamesTest.testDontCoalesceClassDeclarationsWithDestructuringDeclaration()
func (c *CoalesceVariableNamesTest) TestDontCoalesceClassDeclarationsWithDestructuringDeclaration() {
	// See https://github.com/google/closure-compiler/issues/3019 - this used to cause a syntax
	// error by coalescing `B` and `C` without converting `class B {}` to a non-block-scoped
	// declaration.
	c.TestSame(
		"function f(obj) {\n" +
			"  class B {}\n" +
			"  console.log(B);\n" +
			"  const {default: C} = obj;\n" +
			"  return {obj, C};\n" +
			"}\n")
}

// Java: CoalesceVariableNamesTest.testObjDestructuringConst()
func (c *CoalesceVariableNamesTest) TestObjDestructuringConst() {
	c.Test(
		"function f(obj) {\n"+
			"  {\n"+
			"    const {foo: foo} = obj;\n"+
			"    alert(foo);\n"+
			"  }\n"+
			"}\n",
		"function f(obj) {\n"+
			"  {\n"+
			"    ({foo: obj} = obj);\n"+
			"    alert(obj);\n"+
			"  }\n"+
			"}\n")
}

// Java: CoalesceVariableNamesTest.testObjDestructuringConstWithMultipleDeclarations()
func (c *CoalesceVariableNamesTest) TestObjDestructuringConstWithMultipleDeclarations() {
	c.Test(
		"function f(obj) {\n"+
			"  {\n"+
			"    const {foo: foo} = obj;\n"+
			"    alert(foo);\n"+
			"  }\n"+
			"  {\n"+
			"    const {bar: bar} = obj;\n"+
			"    alert(bar);\n"+
			"  }\n"+
			"}\n",
		"function f(obj) {\n"+
			"  {\n"+
			"    const {foo: foo} = obj;\n"+
			"    alert(foo);\n"+
			"  }\n"+
			"  {\n"+
			"    ({bar: obj} = obj);\n"+
			"    alert(obj);\n"+
			"  }\n"+
			"}\n")
}

// Java: CoalesceVariableNamesTest.testObjDestructuringConstWithMultipleLvaluesInDecl()
func (c *CoalesceVariableNamesTest) TestObjDestructuringConstWithMultipleLvaluesInDecl() {
	c.TestSame(
		"function f() {\n" +
			"  const obj = {};\n" +
			"  const {prop1: foo, prop2: bar} = obj;\n" +
			"  alert(foo);\n" +
			"}\n")
}

// Java: CoalesceVariableNamesTest.testObjDestructuringVar()
func (c *CoalesceVariableNamesTest) TestObjDestructuringVar() {
	c.TestSame(
		"function f(param) {\n" +
			"  var bar;\n" +
			"  ({prop1:param, prop2:bar} = param);\n" +
			"  alert(param);\n" +
			"}\n")
}

// Java: CoalesceVariableNamesTest.testObjDestructuringVarInAsyncFn()
func (c *CoalesceVariableNamesTest) TestObjDestructuringVarInAsyncFn() {
	c.TestSame(
		"async function f(param) {\n" +
			"  var bar;\n" +
			"  ({prop1:param, prop2:bar} = param);\n" +
			"  alert(param);\n" +
			"}\n")
}

// Java: CoalesceVariableNamesTest.testObjDestructuringVarInGeneratorFn()
func (c *CoalesceVariableNamesTest) TestObjDestructuringVarInGeneratorFn() {
	c.TestSame(
		"function *f(param) {\n" +
			"  var bar;\n" +
			"  ({prop1:param, prop2:bar} = param);\n" +
			"  alert(param);\n" +
			"}\n")
}

// Java: CoalesceVariableNamesTest.testObjDestructuringVarInAsyncGeneratorFn()
func (c *CoalesceVariableNamesTest) TestObjDestructuringVarInAsyncGeneratorFn() {
	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	c.TestSame(
		"async function *f(param) {\n" +
			"  var bar;\n" +
			"  ({prop1:param, prop2:bar} = param);\n" +
			"  alert(param);\n" +
			"}\n")
}

// Java: CoalesceVariableNamesTest.testLetWithSingleLValuesInForLoopCoalesced()
func (c *CoalesceVariableNamesTest) TestLetWithSingleLValuesInForLoopCoalesced() {
	c.Test(
		"function f(x) {\n"+
			"  for (let y = x + 1;;) {\n"+
			"    alert(y);\n"+
			"  }\n"+
			"}\n",
		"function f(x) {\n"+
			"  for (x = x + 1;;) {\n"+
			"    alert(x);\n"+
			"  }\n"+
			"}\n")
}

// Java: CoalesceVariableNamesTest.testLetWithMultipleLValuesInForLoopNotCoalesced()
func (c *CoalesceVariableNamesTest) TestLetWithMultipleLValuesInForLoopNotCoalesced() {
	c.TestSame(
		"function f(x) {\n" +
			"  for (let y = x + 1, z = 0;;) {\n" +
			"    alert(y + z);\n" +
			"  }\n" +
			"}\n")
}

// Java: CoalesceVariableNamesTest.testConstDestructuringDeclInForOf_dropsConst()
func (c *CoalesceVariableNamesTest) TestConstDestructuringDeclInForOf_dropsConst() {
	c.Test(
		"function f(param) { for (let [y] of []) {} }",
		"function f(param) { for ([param] of []) {} }")
}

// Java: CoalesceVariableNamesTest.testConstDestructuringInForOfCoalescedWithUseInBlock()
func (c *CoalesceVariableNamesTest) TestConstDestructuringInForOfCoalescedWithUseInBlock() {
	// TODO(b/121276933): coalesce `x` and `y`
	c.InFunction("var x = 1; for (let [y] of iter) { y }")
}

// Java: CoalesceVariableNamesTest.testReplaceRhsOfDestructuringDeclaration()
func (c *CoalesceVariableNamesTest) TestReplaceRhsOfDestructuringDeclaration() {
	c.InFunctionWithExpected(
		"let unused = 0; let arr = [1, 2, 3]; const [a, b, c] = arr; alert(a + b + c);",
		"var unused = 0; unused = [1, 2, 3]; const [a, b, c] = unused; alert(a + b + c);")
}

// Java: CoalesceVariableNamesTest.testReplaceRhsOfDestructuringDeclaration_withPseudoNames()
func (c *CoalesceVariableNamesTest) TestReplaceRhsOfDestructuringDeclaration_withPseudoNames() {
	c.usePseudoName = true
	c.InFunctionWithExpected(
		"let unused = 0; let arr = [1, 2, 3]; const [a, b, c] = arr; alert(a + b + c);",
		"var arr_unused = 0;\n"+
			"arr_unused = [1, 2, 3];\n"+
			"const [a, b, c] = arr_unused;\n"+
			"alert(a + b + c);\n")
}

// Java: CoalesceVariableNamesTest.testDestructuringDefaultValue()
func (c *CoalesceVariableNamesTest) TestDestructuringDefaultValue() {
	c.TestSame("function f(param) {  var a;  [a = param] = {};  param;  }")

	c.Test(
		"function f(param) {  var a;  [a = param] = {};  a;  }",
		"function f(param) {  [param = param] = {};  param;  }")
}

// Java: CoalesceVariableNamesTest.testSpread_ofArray_consideredRead()
func (c *CoalesceVariableNamesTest) TestSpread_ofArray_consideredRead() {
	c.TestSame(
		"function f() {\n" +
			"  var a = 9;\n" +
			"  var b = 6;\n" +
			"\n" +
			"  ([...a]); // Read `a`.\n" +
			"  return b;\n" +
			"}\n")

	c.Test(
		"function f() {\n"+
			"  var a = 9;\n"+
			"  read(a);\n"+
			"\n"+
			"  var b = 6;\n"+
			"  ([...b]); // Read `b`.\n"+
			"}\n",
		"function f() {\n"+
			"  var a = 9;\n"+
			"  read(a);\n"+
			"\n"+
			"  a = 6;\n"+
			"  ([...a]);\n"+
			"}\n")
}

// Java: CoalesceVariableNamesTest.testSpread_ofObject_consideredRead()
func (c *CoalesceVariableNamesTest) TestSpread_ofObject_consideredRead() {
	c.TestSame(
		"function f() {\n" +
			"  var a = 9;\n" +
			"  var b = 6;\n" +
			"\n" +
			"  ({...a}); // Read `a`.\n" +
			"  return b;\n" +
			"}\n")

	c.Test(
		"function f() {\n"+
			"  var a = 9;\n"+
			"  read(a);\n"+
			"\n"+
			"  var b = 6;\n"+
			"  ({...b}); // Read `b`.\n"+
			"}\n",
		"function f() {\n"+
			"  var a = 9;\n"+
			"  read(a);\n"+
			"\n"+
			"  a = 6;\n"+
			"  ({...a});\n"+
			"}\n")
}

// Java: CoalesceVariableNamesTest.testRest_fromArray_consideredWrite()
func (c *CoalesceVariableNamesTest) TestRest_fromArray_consideredWrite() {
	c.TestSame(
		"function f() {\n" +
			"  var a = 9;\n" +
			"  var b = 6;\n" +
			"\n" +
			"  ([...a] = itr); // Write `a`.\n" +
			"  return b;\n" +
			"}\n")

	c.Test(
		"function f() {\n"+
			"  var a = 9;\n"+
			"  read(a);\n"+
			"\n"+
			"  var b = 6;\n"+
			"  ([...b] = itr); // Write `b`.\n"+
			"}\n",
		"function f() {\n"+
			"  var a = 9;\n"+
			"  read(a);\n"+
			"\n"+
			"  a = 6;\n"+
			"  ([...a] = itr);\n"+
			"}\n")
}

// Java: CoalesceVariableNamesTest.testRest_fromObject_consideredWrite()
func (c *CoalesceVariableNamesTest) TestRest_fromObject_consideredWrite() {
	c.TestSame(
		"function f() {\n" +
			"  var a = 9;\n" +
			"  var b = 6;\n" +
			"\n" +
			"  ({...a} = obj); // Write `a`.\n" +
			"  return b;\n" +
			"}\n")

	c.Test(
		"function f() {\n"+
			"  var a = 9;\n"+
			"  read(a);\n"+
			"\n"+
			"  var b = 6;\n"+
			"  ({...b} = itr); // Write `b`.\n"+
			"}\n",
		"function f() {\n"+
			"  var a = 9;\n"+
			"  read(a);\n"+
			"\n"+
			"  a = 6;\n"+
			"  ({...a} = itr);\n"+
			"}\n")
}

// Java: CoalesceVariableNamesTest.testDestructuringEvaluationOrder()
func (c *CoalesceVariableNamesTest) TestDestructuringEvaluationOrder() {
	// Since the "a = 5" assignment is evaluated before "a = param" (which is
	// conditionally evaluated), we must not coalesce param and a.
	c.TestSame("function f(param) { var a; [a = param] = (a = 5, {});  a; }")
}

// We would normally coalesce 'key' with 'collidesWithKey', but in doing so we'd change the 'let'
// on line 2 to a 'var' which would cause the inner function to capture the wrong value of 'val'.
// Java: CoalesceVariableNamesTest.testCaptureLet()
func (c *CoalesceVariableNamesTest) TestCaptureLet() {
	c.TestSame(
		"function f(param) {\n" +
			"  for (let [key, val] of foo()) {\n" +
			"    param = (x) => { return val(x); };\n" +
			"  }\n" +
			"  let collidesWithKey = 5;\n" +
			"  return param(collidesWithKey);\n" +
			"}\n")
}

// Compare to the earlier case. Since the two-lvalue declaration `var [key, val]` gets normalized
// we still coalesce `key` with `collidesWithKey`.
// Java: CoalesceVariableNamesTest.testCaptureVar()
func (c *CoalesceVariableNamesTest) TestCaptureVar() {
	c.Test(
		"function f(param) {\n"+
			"  for (var [key, val] of foo()) {\n"+
			"    param = (x) => { return val(x); };\n"+
			"  }\n"+
			"  let collidesWithKey = 5;\n"+
			"  return param(collidesWithKey);\n"+
			"}\n",
		"function f(param) {\n"+
			"  var key;\n"+
			"  var val;\n"+
			"  for ([key, val] of foo()) {\n"+
			"    param = (x) => { return val(x); };\n"+
			"  }\n"+
			"  key = 5;\n"+
			"  return param(key);\n"+
			"}\n")
}

// Java: CoalesceVariableNamesTest.testDestructuring()
func (c *CoalesceVariableNamesTest) TestDestructuring() {
	c.TestSame(
		"function f() {\n" +
			"  const [x, y] = foo(5);\n" +
			"  let z = foo(x);\n" +
			"  return x;\n" +
			"}\n")
}

// Java: CoalesceVariableNamesTest.testDeterministic()
func (c *CoalesceVariableNamesTest) TestDeterministic() {
	// Make the variable interference graph a pentagon.
	//         a - b
	//        /     \
	//        e     c
	//         \   /
	//           d
	// The coloring partitioning would be:
	//  a = { a, c }
	//  b = { b, d }
	//  e = { e }
	c.InFunctionWithExpected(
		"var a;\n"+
			"var b;\n"+
			"var c;\n"+
			"var d;\n"+
			"var e;\n"+
			"a=1; b=1; a; b;\n"+
			"b=1; c=1; b; c;\n"+
			"c=1; d=1; c; d;\n"+
			"d=1; e=1; d; e;\n"+
			"e=1; a=1; e; a;\n",
		"var a;\n"+
			"var b;\n"+
			"var e;\n"+
			"a=1; b=1; a; b;\n"+
			"b=1; a=1; b; a;\n"+
			"a=1; b=1; a; b;\n"+
			"b=1; e=1; b; e;\n"+
			"e=1; a=1; e; a;\n")

	// If we favor "d" first by declaring "d" earlier,
	// the coloring partitioning would be:
	//  b = { b, e }
	//  d = { d, a }
	//  c = { c }
	c.InFunctionWithExpected(
		"var d,a,b,c,e;\n"+
			"  a=1; b=1; a; b;\n"+
			"  b=1; c=1; b; c;\n"+
			"  c=1; d=1; c; d;\n"+
			"  d=1; e=1; d; e;\n"+
			"  e=1; a=1; e; a;\n",
		"var d;var b;var c;\n"+
			"d=1;b=1;d;b;\n"+
			"b=1;c=1;b;c;\n"+
			"c=1;d=1;c;d;\n"+
			"d=1;b=1;d;b;\n"+
			"b=1;d=1;b;d\n")
}

// Sometimes live range can be cross even within a VAR declaration.
// Java: CoalesceVariableNamesTest.testVarLiveRangeCross()
func (c *CoalesceVariableNamesTest) TestVarLiveRangeCross() {
	c.InFunctionWithExpected("var a={}; var b=a.S(); b", "var a={};     a=a.S(); a")

	c.InFunctionWithExpected(
		"var a = {}; var b = a.S(),     c = b.SS(); b; c",
		"var a = {};     a = a.S(); var c = a.SS(); a; c")

	c.InFunctionWithExpected(
		"var a={}; var b=a.S(); var c=a.SS(); var d=a.SSS(); b; c; d",
		"var a={}; var b=a.S(); var c=a.SS();     a=a.SSS(); b; c; a")

	c.InFunctionWithExpected(
		"var a={}; var b=a.S(); var c=a.SS(); var d=a.SSS(); b; c; d",
		"var a={}; var b=a.S(); var c=a.SS();     a=a.SSS(); b; c; a")

	c.InFunction("var a={}; d=1; d; var b=a.S(); var c=a.SS(); var d=a.SSS(); b; c; d")
}

// Java: CoalesceVariableNamesTest.testBug1445366()
func (c *CoalesceVariableNamesTest) TestBug1445366() {
	// An assignment might not be complete if the RHS throws an exception.
	c.InFunction(
		"var iframe = getFrame();\n" +
			"try {\n" +
			"  var win = iframe.contentWindow;\n" +
			"} catch (e) {\n" +
			"} finally {\n" +
			"  if (win)\n" +
			"    this.setupWinUtil_();\n" +
			"  else\n" +
			"    this.load();\n" +
			"}\n")

	// Verify that we can still coalesce it if there are no handlers.
	c.InFunctionWithExpected(
		"var iframe = getFrame();\n"+
			"var win = iframe.contentWindow;\n"+
			"if (win)\n"+
			"  this.setupWinUtil_();\n"+
			"else\n"+
			"  this.load();\n",
		"var iframe = getFrame();\n"+
			"iframe = iframe.contentWindow;\n"+
			"if (iframe)\n"+
			"  this.setupWinUtil_();\n"+
			"else\n"+
			"  this.load();\n")
}

// Parameter 'e' is never used, but if we coalesce 'command' with 'e' then the 'if (command)'
// check will produce an incorrect result if none of the 'case' statements is executed.
// Java: CoalesceVariableNamesTest.testCannotReuseAnyParamsBug()
func (c *CoalesceVariableNamesTest) TestCannotReuseAnyParamsBug() {
	c.TestSame(
		"function handleKeyboardShortcut(e, key, isModifierPressed) {\n" +
			"  if (!isModifierPressed) {\n" +
			"    return false;\n" +
			"  }\n" +
			"  var command;\n" +
			"  switch (key) {\n" +
			"    case 'b': // Ctrl+B\n" +
			"      command = COMMAND.BOLD;\n" +
			"      break;\n" +
			"    case 'i': // Ctrl+I\n" +
			"      command = COMMAND.ITALIC;\n" +
			"      break;\n" +
			"    case 'u': // Ctrl+U\n" +
			"      command = COMMAND.UNDERLINE;\n" +
			"      break;\n" +
			"    case 's': // Ctrl+S\n" +
			"      return true;\n" +
			"  }\n" +
			"\n" +
			"  if (command) {\n" +
			"    this.fieldObject.execCommand(command);\n" +
			"    return true;\n" +
			"  }\n" +
			"\n" +
			"  return false;\n" +
			"};\n")
}

// Java: CoalesceVariableNamesTest.testCannotReuseAnyParamsBugWithDestructuring()
func (c *CoalesceVariableNamesTest) TestCannotReuseAnyParamsBugWithDestructuring() {
	c.TestSame(
		"function handleKeyboardShortcut({type: type}, key, isModifierPressed) {\n" +
			"  if (!isModifierPressed) {\n" +
			"    return false;\n" +
			"  }\n" +
			"  var command;\n" +
			"  switch (key) {\n" +
			"    case 'b': // Ctrl+B\n" +
			"      command = COMMAND.BOLD;\n" +
			"      break;\n" +
			"    case 'i': // Ctrl+I\n" +
			"      command = COMMAND.ITALIC;\n" +
			"      break;\n" +
			"    case 'u': // Ctrl+U\n" +
			"      command = COMMAND.UNDERLINE;\n" +
			"      break;\n" +
			"    case 's': // Ctrl+S\n" +
			"      return true;\n" +
			"  }\n" +
			"\n" +
			"  if (command) {\n" +
			"    this.fieldObject.execCommand(command);\n" +
			"    return true;\n" +
			"  }\n" +
			"\n" +
			"  return false;\n" +
			"};\n")
}

// Java: CoalesceVariableNamesTest.testForInWithAssignment()
func (c *CoalesceVariableNamesTest) TestForInWithAssignment() {
	c.InFunctionWithExpected(
		"function f(commands) {\n"+
			"  var k, v, ref;\n"+
			"  for (k in ref = commands) {\n"+
			"    v = ref[k];\n"+
			"    alert(k + ':' + v);\n"+
			"  }\n"+
			"}\n",
		"function f(commands){\n"+
			"var k;\n"+
			"var ref;\n"+
			"for(k in ref = commands) {\n"+
			"  commands=ref[k];\n"+
			"  alert(k+':'+commands)\n"+
			"}}\n")
}

// Java: CoalesceVariableNamesTest.testUsePseudoNames()
func (c *CoalesceVariableNamesTest) TestUsePseudoNames() {
	c.usePseudoName = true
	c.InFunctionWithExpected(
		"var x   = 0; print(x  ); var   y = 1; print(  y)",
		"var x_y = 0; print(x_y);     x_y = 1; print(x_y)")

	c.InFunctionWithExpected(
		"var x_y = 1; var x    = 0; print(x   ); var     y = 1; print(   y); print(x_y);",
		"var x_y = 1; var x_y$ = 0; print(x_y$);      x_y$ = 1; print(x_y$); print(x_y);")

	c.InFunctionWithExpected(
		"var x_y = 1;\n"+
			"function f() {\n"+
			"  var x    = 0;\n"+
			"  print(x  );\n"+
			"  var y = 1;\n"+
			"  print( y);\n"+
			"  print(x_y);\n"+
			"}\n",
		"function f(){\n"+
			"  var x_y$=0;\n"+
			"  print(x_y$);\n"+
			"  x_y$=1;\n"+
			"  print(x_y$);\n"+
			"  print(x_y)\n"+
			"}\n"+
			"var x_y=1\n")

	c.InFunctionWithExpected(
		"var x   = 0;\n"+
			"print(x  );\n"+
			"var   y = 1;\n"+
			"print(  y);\n"+
			"var closure_var;\n"+
			"function bar() {\n"+
			"  print(closure_var);\n"+
			"}\n",
		"function bar(){\n"+
			"  print(closure_var)\n"+
			"}\n"+
			"var x_y=0;\n"+
			"print(x_y);\n"+
			"x_y=1;\n"+
			"print(x_y);\n"+
			"var closure_var\n")
}

// Java: CoalesceVariableNamesTest.testUsePseudoNamesWithLets()
func (c *CoalesceVariableNamesTest) TestUsePseudoNamesWithLets() {
	c.usePseudoName = true
	c.InFunctionWithExpected(
		"var x_y = 1;\n"+
			"function f() {\n"+
			"  let x    = 0;\n"+
			"  print(x  );\n"+
			"  let y = 1;\n"+
			"  print( y);\n"+
			"  print(x_y);\n"+
			"}\n",
		"function f(){\n"+
			"  var x_y$=0;\n"+
			"  print(x_y$);\n"+
			"  x_y$ = 1;\n"+
			"  print(x_y$);\n"+
			"  print(x_y)\n"+
			"}\n"+
			"var x_y=1\n")
}

// Java: CoalesceVariableNamesTest.testMaxVars()
func (c *CoalesceVariableNamesTest) TestMaxVars() {
	code := ""
	for i := 0; i < LiveVariablesAnalysis_MAX_VARIABLES_TO_ANALYZE+1; i++ {
		code += fmt.Sprintf("var x%d = 0; print(x%d);", i, i)
	}
	c.InFunction(code)
}

// Testing Es6 features
// Java: CoalesceVariableNamesTest.testCoalesceInInnerBlock()
func (c *CoalesceVariableNamesTest) TestCoalesceInInnerBlock() {
	c.InFunctionWithExpected("{ var x = 1; var y = 2; y }", "{ var x = 1;     x = 2; x }")

	c.InFunctionWithExpected("var x = 1; var y = 2; y;", "var x = 1;     x = 2; x;")
}

// Java: CoalesceVariableNamesTest.testLetSimple()
func (c *CoalesceVariableNamesTest) TestLetSimple() {
	c.InFunctionWithExpected("let x = 0; x; let y = 5; y", "var x = 0; x;     x = 5; x")

	c.InFunctionWithExpected(
		"var x = 1; var y = 2; { let z = 3; y; }", "var x = 1;     x = 2; { let z = 3; x; }")

	// First let in a block - It is unsafe for { let x = 0; x; } let y = 1; to be coalesced as
	// { let x = 0; x; } x = 1; because x will be out of scope outside of the inner scope!
	c.InFunctionWithExpected("{ let x = 0; x; } let y = 5; y;", "{ var x = 0; x; }     x = 5; x;")

	// The following situation will never happen in practice because at this point, the code has
	// been normalized so no two variables will have the same name
	// --> var x = 1; var y = 2; { let x = 3; y }
}

// Java: CoalesceVariableNamesTest.testLetDifferentBlocks()
func (c *CoalesceVariableNamesTest) TestLetDifferentBlocks() {
	c.InFunctionWithExpected(
		"var x = 0; if (1) { let y = 1; x } else { let z = 1; x }",
		"var x = 0; if (1) { var y = 1; x } else {     y = 1; x }")

	c.InFunctionWithExpected(
		"var x = 0; if (1) { let y = 1; y } else { let z = 1 + x; z }",
		"var x = 0; if (1) {     x = 1; x } else {     x = 1 + x; x }")

	c.InFunctionWithExpected(
		"var x = 0; if (1) { let y = 1; y } else { let z = 1; z }; x",
		"var x = 0; if (1) { var y = 1; y } else {     y = 1; y }; x")

	c.InFunctionWithExpected(
		"if (1) { var x = 0; let y = 1; y + x} else { let z = 1; z } y;",
		"if (1) { var x = 0; let y = 1; y + x} else {     x = 1; x } y;")

	c.InFunctionWithExpected(
		"if(1) { var x = 0; x } else { var y = 0; y }",
		"if(1) { var x = 0; x } else {     x = 0; x }")

	c.InFunctionWithExpected(
		"if (a) {\n"+
			"   return a;\n"+
			" } else {\n"+
			"   let b = a;\n"+
			"   let c = 1;\n"+
			"   return c;\n"+
			" }\n"+
			" return a;\n",
		"if (a) {\n"+
			"    return a;\n"+
			"  } else {\n"+
			"    var b = a;\n"+
			"        b = 1;\n"+
			"    return b;\n"+
			"  }\n"+
			"  return a;\n")
}

// Java: CoalesceVariableNamesTest.testLetWhileLoops()
func (c *CoalesceVariableNamesTest) TestLetWhileLoops() {
	// Simple
	// It violates the temporal deadzone for let x = 1; while (1) { x = 2; x; let x = 0; x } to
	// be output.
	c.InFunction("let x = 1; for(;1;) { x; x = 2; let y = 0; y }")

	c.InFunction("let x = 1; for(;1;) { x = 2; x; let y = 0; y } x;")
}

// Java: CoalesceVariableNamesTest.testLetForLoops()
func (c *CoalesceVariableNamesTest) TestLetForLoops() {
	// TODO (simranarora) We should get rid of declaration hoisting from the normalize pass.
	// Right now, because of declaration hoisting, this following test reads the expected code as:
	// var x = 1; for ( ; x < 10; x++) { let y = 2; x + y } x = 3
	c.InFunctionWithExpected(
		"for (let x = 1; x < 10; x ++) { let y = 2; x + y; } let z = 3;",
		"for(var x=1;x<10;x++){let y=2;x+y}x=3")

	c.InFunctionWithExpected(
		"var w = 0; for (let x = 1; x < 10; x ++) { let y = 2; x + y; } var z = 3;",
		"var w = 0; for (    w = 1; w < 10; w ++) { let y = 2; w + y; }     w = 3;")

	// Closure capture of the let variable
	// Here z should not be coalesced because variables used in functions are considered escaped
	// and this pass does not touch any escaped variables
	c.InFunction("let x = 3; for (let z = 1; z < 10; z++) { use(() => {z}); }")

	c.InFunction("for (let x = 1; x < 10; x++) { use(() => {x}); } let z = 3;")

	// Multiple lets declared in loop head
	c.InFunction("for (let x = 1, y = 2, z = 3; (x + z) < 10; x ++) { x + z; }")

	// Here the variable y is redeclared because the variable z in the header of the for-loop has
	// not been declared before
	c.InFunction("let y = 2; for (let x = 1, z = 3; (x + z) < 10; x ++) { x + z; }")
}

// Java: CoalesceVariableNamesTest.testArrowFunctions()
func (c *CoalesceVariableNamesTest) TestArrowFunctions() {
	c.InFunctionWithExpected(
		"var x = 1; var y = () => { let z = 0; z }", "var x = 1;     x = () => { let z = 0; z }")

	c.InFunctionWithExpected(
		"var x = 1; var y = () => { let z = 0; z }; y();",
		"var x = 1;     x = () => { let z = 0; z }; x();")

	c.InFunctionWithExpected(
		"var x = 1; var y = () => { let z = 0; z }; x;",
		"var x = 1; var y = () => { let z = 0; z }; x;")

	c.InFunctionWithExpected(
		"var x = () => { let z = 0; let y = 1; y }", "var x = () => { var z = 0;     z = 1; z }")

	c.InFunctionWithExpected(
		"var x = 1; var y = 2; var f = () => x + 1",
		"var x = 1; var y = 2;     y = () => { return x + 1; }")

	// Coalesce with arrow function parameters
	c.InFunctionWithExpected("(x) => { var y = 1; y; }", "(x) => {     x = 1; x; }")

	c.InFunctionWithExpected("(x) => { let y = 1; y; }", "(x) => {     x = 1; x; }")
}

// Java: CoalesceVariableNamesTest.testCodeWithTwoFunctions()
func (c *CoalesceVariableNamesTest) TestCodeWithTwoFunctions() {
	// We only want to coalesce within a function, not across functions
	c.Test(
		"function FUNC1() {\n"+
			"  var x = 1;\n"+
			"  var y = 2;\n"+
			"          y;\n"+
			"}\n"+
			"function FUNC2() {\n"+
			"  var z = 3;\n"+
			"  var w = 4;\n"+
			"          w;\n"+
			"}\n",
		"function FUNC1() {\n"+
			"  var x = 1;\n"+
			"      x = 2;\n"+
			"          x;\n"+
			"}\n"+
			"function FUNC2() {\n"+
			"  var z = 3;\n"+
			"      z = 4;\n"+
			"          z;\n"+
			"}\n")

	// Two arrow functions
	c.Test(
		"() => { var x = 1; var y = 2; y; };\n"+
			"() => { var z = 3; var w = 4; w; };\n",
		"() => { var x = 1;     x = 2; x; };\n"+
			"() => { var z = 3;     z = 4; z; };\n")
}

// Java: CoalesceVariableNamesTest.testNestedFunctionCoalescing()
func (c *CoalesceVariableNamesTest) TestNestedFunctionCoalescing() {
	c.Test(
		"function FUNC1() {\n"+
			"  var x = 1;\n"+
			"  var y = 2;\n"+
			"          y;\n"+
			"  function FUNC2() {\n"+
			"    var z = 3;\n"+
			"    var w = 4;\n"+
			"            w;\n"+
			"  }\n"+
			"}\n",
		"function FUNC1() {\n"+
			"  function FUNC2() {\n"+
			"    var z = 3;\n"+
			"        z = 4;\n"+
			"            z\n"+
			"  }\n"+
			"  var x = 1;\n"+
			"  x = 2;\n"+
			"  x\n"+
			"}\n")
}

// Java: CoalesceVariableNamesTest.testInstanceFieldInitializerCaptureIsNotClobbered()
func (c *CoalesceVariableNamesTest) TestInstanceFieldInitializerCaptureIsNotClobbered() {
	c.TestSame(
		`function render(input) {
  let s = sanitize(input);
  class Box { html = s; }
  let t = input;
  log(t, input);
  return new Box();
}
`)
}

// Java: CoalesceVariableNamesTest.testControl_innerFunctionCaptureIsNotClobbered()
func (c *CoalesceVariableNamesTest) TestControl_innerFunctionCaptureIsNotClobbered() {
	// `s` is escaped (FUNCTION is a CFG root), so it is excluded from the interference graph
	// and `t` cannot be coalesced onto it. Nothing coalesces; `t` stays `t`.
	c.TestSame(
		`function render(input) {
  let s = sanitize(input);
  let make = function() { return {html: s}; };
  let t = input;
  log(t, input);
  return make();
}
`)
}

// Java: CoalesceVariableNamesTest.testComputedFieldInitializerCaptureIsNotClobbered()
func (c *CoalesceVariableNamesTest) TestComputedFieldInitializerCaptureIsNotClobbered() {
	c.TestSame(
		`function render(input) {
  let s = sanitize(input);
  class Box { ["html"] = s; }
  let t = input;
  log(t, input);
  return new Box();
}
`)
}

// Java: CoalesceVariableNamesTest.inFunction(String)
func (c *CoalesceVariableNamesTest) InFunction(src string) {
	c.TestSame("function FUNC(){" + src + "}")
}

// Java: CoalesceVariableNamesTest.inFunction(String, String)
func (c *CoalesceVariableNamesTest) InFunctionWithExpected(src string, expected string) {
	c.Test("function FUNC(){"+src+"}", "function FUNC(){"+expected+"}")
}
