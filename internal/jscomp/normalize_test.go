// Port of com/google/javascript/jscomp/NormalizeTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/NormalizeTest.java)

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

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

const NormalizeTest_EXTERNS = "var window; var Arguments;"

// NormalizeTest is the unit test of Normalize.
type NormalizeTest struct {
	CompilerTestCaseBase
}

// Java: NormalizeTest()
func NewNormalizeTest() *NormalizeTest {
	c := &NormalizeTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, NormalizeTest_EXTERNS)
	return c
}

func TestNormalizeTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewNormalizeTest() })
}

// Java: NormalizeTest.getProcessor(Compiler)
func (c *NormalizeTest) GetProcessor(compiler *Compiler) CompilerPass {
	return Normalize_CreateNormalizeForOptimizations(compiler)
}

// Java: NormalizeTest.getOptions()
func (c *NormalizeTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetWarningLevel(DiagnosticGroups_MODULE_LOAD, CheckLevel_OFF)
	return options
}

// SetUp runs CompilerTestCase's @Before method, then NormalizeTest's
// (JUnit runs superclass @Before methods first).
func (c *NormalizeTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.CustomSetUp()
}

// Java: NormalizeTest.customSetUp()
func (c *NormalizeTest) CustomSetUp() {
	// Validate that Normalize copies colors onto any nodes it synthesizes
	c.EnableTypeInfoValidation()
	c.EnableTypeCheck()
	c.ReplaceTypesWithColors()
}

// Java: NormalizeTest.testMultipleForOfLoopsWithSameNameInductionVariable()
func (c *NormalizeTest) TestMultipleForOfLoopsWithSameNameInductionVariable() {
	srcs :=
		c.Srcs(
			`function* inorder1(t) {
    for (var x of []) {
      yield x;
    }
    for (var x of []) {
      yield x;
    }
}
`)
	expected :=
		c.Expected(
			`function* inorder1(t) {
    var x;
    for (x of []) {
      yield x;
    }
    for (x of []) {
      yield x;
    }
}
`)
	c.Test(srcs, expected) //
}

// Java: NormalizeTest.testConstAnnotationPropagation()
func (c *NormalizeTest) TestConstAnnotationPropagation() {
	c.Test(
		"const x = 3; var a,     b; var y = x + 2;", //
		"const x = 3; var a; var b; var y = x + 2;")
	root := c.GetLastCompiler().GetRoot()
	scriptNode :=
		root.GetLastChild(). // ROOT of input sources
					GetLastChild()
	// `const x = 3;`
	constNode := scriptNode.GetFirstChild()
	// `x`
	constName := constNode.GetOnlyChild()
	AssertThat(constName.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME)).IsTrue()
	// `var y = x + 2;`
	lastStatement := scriptNode.GetLastChild()
	// `y = x + 2`
	yVar := lastStatement.GetOnlyChild()
	secondNameNodeOfX :=
		yVar.GetFirstChild(). // `x + 2`
					GetFirstChild()
	AssertThat(secondNameNodeOfX.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME)).IsTrue()
}

// Java: NormalizeTest.testRestConstAnnotationPropagation()
func (c *NormalizeTest) TestRestConstAnnotationPropagation() {
	c.TestSame(
		`const {...x} = {a: 3};
var y = x;
`)
	root := c.GetLastCompiler().GetRoot()
	scriptNode :=
		root.GetLastChild(). // ROOT of input sources
					GetLastChild()
	// `const {...x} = {a: 3};`
	constNode := scriptNode.GetFirstChild()
	// `{...x}`
	objectPattern :=
		constNode.
			GetFirstChild(). // DESTRUCTURING_LHS
			GetFirstChild()
	// `x`
	constName :=
		objectPattern.
			GetFirstChild(). // OBJECT_REST
			GetFirstChild()
	AssertThat(constName.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME)).IsTrue()
	// `var y = x`
	lastStatement := scriptNode.GetLastChild()
	// `y = x`
	yVar := lastStatement.GetOnlyChild()
	secondNameNodeOfX := yVar.GetFirstChild()
	AssertThat(secondNameNodeOfX.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME)).IsTrue()
}

// Java: NormalizeTest.testRestConstAnnotationPropagation_onlyConstVars()
func (c *NormalizeTest) TestRestConstAnnotationPropagation_onlyConstVars() {
	c.TestSame(
		`const obj = {a: 3, b: 'string', c: null};
const {...rest} = obj;
const y = rest;
`)
	root := c.GetLastCompiler().GetRoot()
	scriptNode :=
		root.GetLastChild(). // ROOT of input sources
					GetLastChild()
	firstStatement := scriptNode.GetFirstChild()
	// `obj` from `const obj = {a: 3, b: 'string', c: null}`
	objName := firstStatement.GetFirstChild()
	AssertThat(objName.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME)).IsTrue()
	// `const {...rest} = obj`
	constNode := firstStatement.GetNext()
	// `{...rest} = obj
	destructuringNode := constNode.GetFirstChild()
	// `obj` from previous comment
	secondObjName := destructuringNode.GetLastChild()
	AssertThat(secondObjName.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME)).IsTrue()
	// `{...rest}
	objectPattern := destructuringNode.GetFirstChild()
	// `rest`
	constName :=
		objectPattern.
			GetFirstChild(). // OBJECT_REST
			GetFirstChild()
	AssertThat(constName.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME)).IsTrue()
	// `const y = rest`
	secondConstNode := constNode.GetNext()
	// `y = rest`
	yVar := secondConstNode.GetOnlyChild()
	AssertThat(yVar.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME)).IsTrue()
	secondConstName := yVar.GetFirstChild()
	AssertThat(secondConstName.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME)).IsTrue()
}

// Java: NormalizeTest.testConstRHSPropagation()
func (c *NormalizeTest) TestConstRHSPropagation() {
	c.TestSame("const obj = function inner() {inner();};")
	root := c.GetLastCompiler().GetRoot()
	scriptNode :=
		root.GetLastChild(). // ROOT of input sources
					GetLastChild()
	firstStatement := scriptNode.GetFirstChild()
	// `obj` from `const obj = function inner() {inner();};`
	objName := firstStatement.GetFirstChild()
	AssertThat(objName.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME)).IsTrue()
	functionName := objName.GetOnlyChild().GetFirstChild()
	AssertThat(functionName.IsName()).IsTrue()
	AssertThat(functionName.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME)).IsTrue()
}

// Java: NormalizeTest.testConstRHSPropagation2()
func (c *NormalizeTest) TestConstRHSPropagation2() {
	c.TestSame("/** @const */ var obj = function inner() {inner();};")
	root := c.GetLastCompiler().GetRoot()
	scriptNode :=
		root.GetLastChild(). // ROOT of input sources
					GetLastChild()
	firstStatement := scriptNode.GetFirstChild()
	// `obj` from `var obj = function inner() {inner();};`
	objName := firstStatement.GetFirstChild()
	AssertThat(objName.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME)).IsTrue()
	functionName := objName.GetOnlyChild().GetFirstChild()
	AssertThat(functionName.IsName()).IsTrue()
	AssertThat(functionName.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME)).IsTrue()
}

// Java: NormalizeTest.testNullishCoalesce()
func (c *NormalizeTest) TestNullishCoalesce() {
	c.Test("var a = x ?? y, b = foo()", "var a = x ?? y; var b = foo()")
	c.Test(
		`let x = a ?? b;
{ let x = a ?? b; }
{ let x = a ?? b; }
{ let x = a ?? b; }
{ let x = a ?? b; }
`,
		`let x = a ?? b;
{ let x$jscomp$1 = a ?? b; }
{ let x$jscomp$2 = a ?? b; }
{ let x$jscomp$3 = a ?? b; }
{ let x$jscomp$4 = a ?? b; }
`)
}

// Java: NormalizeTest.testSplitVar()
func (c *NormalizeTest) TestSplitVar() {
	c.TestSame("var a")
	c.Test("var a, b", "var a; var b")
	c.Test("var a, b, c", "var a; var b; var c")
	c.TestSame("var a = 0 ")
	c.Test("var a = 0 , b = foo()", "var a = 0; var b = foo()")
	c.Test("var a = 0, b = 1, c = 2", "var a = 0; var b = 1; var c = 2")
	c.Test(
		"var a = foo(1), b = foo(2), c = foo(3)", "var a = foo(1); var b = foo(2); var c = foo(3)")
	c.Test(
		"try{var b = foo(1), c = foo(2);} finally { foo(3) }",
		"try{var b = foo(1); var c = foo(2)} finally { foo(3); }")
	c.Test(
		"try{var b = foo(1),c = foo(2);} finally {}",
		"try{var b = foo(1); var c = foo(2)} finally {}")
	c.Test(
		"try{foo(0);} finally { var b = foo(1), c = foo(2); }",
		"try{foo(0);} finally {var b = foo(1); var c = foo(2)}")
	c.Test(
		"switch(a) {default: var b = foo(1), c = foo(2); break;}",
		"switch(a) {default: var b = foo(1); var c = foo(2); break;}")
	c.Test("do var a = foo(1), b; while(false);", "do{var a = foo(1); var b} while(false);")
	c.Test("a:var a,b,c;", "a:{ var a;var b; var c; }")
	c.Test("if (true) a:var a,b;", "if (true)a:{ var a; var b; }")
}

// Java: NormalizeTest.testSplitVar_forLoop()
func (c *NormalizeTest) TestSplitVar_forLoop() {
	// Verify vars extracted from FOR nodes are split.
	c.Test(
		"for(var a = 0, b = foo(1), c = 1; c < b; c++) foo(2)",
		"var a = 0; var b = foo(1); var c = 1; for(; c < b; c++) foo(2)")
	// Verify split vars properly introduce blocks when needed.
	c.Test("for(;;) var b = foo(1), c = foo(2);", "for(;;){var b = foo(1); var c = foo(2)}")
	c.Test("for(;;){var b = foo(1), c = foo(2);}", "for(;;){var b = foo(1); var c = foo(2)}")
	c.Test("a:for(var a,b,c;;);", "var a;var b; var c;a:for(;;);")
}

// Java: NormalizeTest.testSplitVar_forLoopCrash()
func (c *NormalizeTest) TestSplitVar_forLoopCrash() {
	// Verify b/174247914
	c.Test(
		`for (let j;;);
var i;
for(var i=0;;);
`,
		`for (let j;;);
var i;
i = 0;
for(;;);
`)
}

// Java: NormalizeTest.testSplitLet()
func (c *NormalizeTest) TestSplitLet() {
	c.TestSame("let a")
	c.Test("let a, b", "let a; let b")
	c.Test("let a, b, c", "let a; let b; let c")
	c.TestSame("let a = 0 ")
	c.Test("let a = 0 , b = foo()", "let a = 0; let b = foo()")
	c.Test("let a = 0, b = 1, c = 2", "let a = 0; let b = 1; let c = 2")
	c.Test(
		"let a = foo(1), b = foo(2), c = foo(3)", "let a = foo(1); let b = foo(2); let c = foo(3)")
	c.TestSame("for (let a = 0, b = 1;;) {}")
}

// Java: NormalizeTest.testLetManyBlocks()
func (c *NormalizeTest) TestLetManyBlocks() {
	c.Test(
		`let a = 'outer';
{ let a = 'inner1'; }
{ let a = 'inner2'; }
{ let a = 'inner3'; }
{ let a = 'inner4'; }
`,
		`let a = 'outer';
{ let a$jscomp$1 = 'inner1'; }
{ let a$jscomp$2 = 'inner2'; }
{ let a$jscomp$3 = 'inner3'; }
{ let a$jscomp$4 = 'inner4'; }
`)
}

// Java: NormalizeTest.testLetOutsideAndInsideForLoop()
func (c *NormalizeTest) TestLetOutsideAndInsideForLoop() {
	c.Test(
		`let a = 'outer';
for (let a = 'inner';;) {
  break;
}
alert(a);
`,
		`let a = 'outer';
for (let a$jscomp$1 = 'inner';;) {
  break;
}
alert(a);
`)
}

// Java: NormalizeTest.testLetOutsideAndInsideBlock()
func (c *NormalizeTest) TestLetOutsideAndInsideBlock() {
	c.Test(
		`let a = 'outer';
{
  let a = 'inner';
}
alert(a);
`,
		`let a = 'outer';
{
  let a$jscomp$1 = 'inner';
}
alert(a);
`)
}

// Java: NormalizeTest.testLetOutsideAndInsideFn()
func (c *NormalizeTest) TestLetOutsideAndInsideFn() {
	c.Test(
		`let a = 'outer';
function f() {
  let a = 'inner';
}
alert(a);
`,
		`let a = 'outer';
function f() {
  let a$jscomp$1 = 'inner';
}
alert(a);
`)
}

// Java: NormalizeTest.testRemoveEmptiesFromClass()
func (c *NormalizeTest) TestRemoveEmptiesFromClass() {
	c.Test(
		`class Foo {
  m1() {};
  m2() {};
}
`,
		`class Foo {
  m1() {}
  m2() {}
}
`)
}

// Java: NormalizeTest.testClassField()
func (c *NormalizeTest) TestClassField() {
	c.Test(
		`/** @unrestricted */
class Foo {
  f1;
  ['f2'] = 1;
  static f3;
  static 'f4' = 'hi';
}
`,
		`class Foo {
  f1
  ['f2'] = 1
  static f3
  static 'f4' = 'hi'
}
`)
}

// Java: NormalizeTest.testClassStaticBlock()
func (c *NormalizeTest) TestClassStaticBlock() {
	c.Test(
		`var x;
class Foo {
  static {
    var x;
    let y;
    this.x;
  }
  static {
    var x;
    let y;
  }
}
class Bar {
  static {
    var x;
    let y;
    this.x;
  }
}
`,
		`var x;
class Foo {
  static {
    var x$jscomp$1;
    let y;
    this.x;
  }
  static {
    var x$jscomp$2;
    let y$jscomp$1;
  }
}
class Bar {
  static {
    var x$jscomp$3;
    let y$jscomp$2;
    this.x;
  }
}
`)
}

// Java: NormalizeTest.testClassStaticBlock_innerFunctionHoisted()
func (c *NormalizeTest) TestClassStaticBlock_innerFunctionHoisted() {
	c.Test(
		`var x;
class Foo {
  static {
    this.x;
    function f1() {}
  }
  static {
    let y;
    function f2() {}
  }
}
class Bar {
  static {
    var z;
    function f3() {}
  }
}
`,
		`var x;
class Foo {
  static {
    function f1() {}
    this.x;
  }
  static {
    function f2() {}
    let y;
  }
}
class Bar {
  static {
    function f3() {}
    var z;
  }
}
`)
}

// Java: NormalizeTest.testClassInForLoop()
func (c *NormalizeTest) TestClassInForLoop() {
	c.TestSame("for (class a {};;) { break; }")
}

// Java: NormalizeTest.testFunctionInForLoop()
func (c *NormalizeTest) TestFunctionInForLoop() {
	c.TestSame("for (function a() {};;) { break; }")
}

// Java: NormalizeTest.testLetInGlobalHoistScope()
func (c *NormalizeTest) TestLetInGlobalHoistScope() {
	c.TestSame(
		`if (true) {
  let x = 1; alert(x);
}
`)
	c.Test(
		`if (true) {
  let x = 1; alert(x);
} else {
  let x = 1; alert(x);
}
`,
		`if (true) {
  let x = 1; alert(x);
} else {
  let x$jscomp$1 = 1; alert(x$jscomp$1);
}
`)
}

// Java: NormalizeTest.testConstInGlobalHoistScope()
func (c *NormalizeTest) TestConstInGlobalHoistScope() {
	c.TestSame(
		`if (true) {
  const x = 1; alert(x);
}
`)
	c.Test(
		`if (true) {
  const x = 1; alert(x);
} else {
  const x = 1; alert(x);
}
`,
		`if (true) {
  const x = 1; alert(x);
} else {
  const x$jscomp$1 = 1; alert(x$jscomp$1);
}
`)
}

// Java: NormalizeTest.testVarReferencedInHoistedFunction()
func (c *NormalizeTest) TestVarReferencedInHoistedFunction() {
	c.Test(
		`var f1 = function() {
  var x;
};

(function () {
  {
    var x = 0;
  }
  function f2() {
    alert(x);
  }
  f2();
})();
`,
		`var f1 = function() {
  var x;
};

(function () {
  function f2() {
    alert(x$jscomp$1);
  }
  {
    var x$jscomp$1 = 0;
  }
  f2();
})();
`)
}

// Java: NormalizeTest.testFunctionDeclInBlockScope()
func (c *NormalizeTest) TestFunctionDeclInBlockScope() {
	c.TestSame("var x; { function g() {} }")
	c.Test("var g; { function g() {} }", "var g; { function g$jscomp$1() {} }")
	c.TestSameInFunction("var x; { function g() {} }")
	c.TestInFunction("var g; { function g() {} }", "var g; { function g$jscomp$1() {} }")
}

// Java: NormalizeTest.testAssignShorthand()
func (c *NormalizeTest) TestAssignShorthand() {
	c.Test("x |= 1;", "x = x | 1;")
	c.Test("x ^= 1;", "x = x ^ 1;")
	c.Test("x &= 1;", "x = x & 1;")
	c.Test("x <<= 1;", "x = x << 1;")
	c.Test("x >>= 1;", "x = x >> 1;")
	c.Test("x >>>= 1;", "x = x >>> 1;")
	c.Test("x += 1;", "x = x + 1;")
	c.Test("x -= 1;", "x = x - 1;")
	c.Test("x *= 1;", "x = x * 1;")
	c.Test("x /= 1;", "x = x / 1;")
	c.Test("x %= 1;", "x = x % 1;")
	c.Test("/** @suppress {const} */ x += 1;", "x = x + 1;")
}

// Java: NormalizeTest.testAssignShorthandDontNormalizeWhenLHSNotName()
func (c *NormalizeTest) TestAssignShorthandDontNormalizeWhenLHSNotName() {
	c.TestSame("obj.x += 1;")
}

// Java: NormalizeTest.testLogicalAssignShorthand()
func (c *NormalizeTest) TestLogicalAssignShorthand() {
	c.Test("x ||= 1;", "x || (x = 1);")
	c.Test("x &&= 1;", "x && (x = 1);")
	c.Test("x ??= 1;", "x ?? (x = 1);")
}

// Java: NormalizeTest.testLogicalAssignPropertyReferenceShorthand()
func (c *NormalizeTest) TestLogicalAssignPropertyReferenceShorthand() {
	c.Test(
		c.Srcs("a.x ||= b"),
		c.Expected(
			`let $jscomp$logical$assign$tmpm1146332801$0;
($jscomp$logical$assign$tmpm1146332801$0 = a).x
   ||
($jscomp$logical$assign$tmpm1146332801$0.x = b);
`))
	c.Test(
		c.Srcs("a.foo &&= null"),
		c.Expected(
			`let $jscomp$logical$assign$tmpm1146332801$0;
($jscomp$logical$assign$tmpm1146332801$0 = a).foo
   &&
($jscomp$logical$assign$tmpm1146332801$0.foo = null);
`))
	c.Test(
		c.Srcs(
			`foo().x = null;
foo().x ??= y
`),
		c.Expected(
			`foo().x = null;
let $jscomp$logical$assign$tmpm1146332801$0;
($jscomp$logical$assign$tmpm1146332801$0 = foo()).x
   ??
($jscomp$logical$assign$tmpm1146332801$0.x = y);
`))
}

// Java: NormalizeTest.testLogicalAssignPropertyReferenceElementShorthand()
func (c *NormalizeTest) TestLogicalAssignPropertyReferenceElementShorthand() {
	c.Test(
		c.Srcs("a[x] ||= b"),
		c.Expected(
			`let $jscomp$logical$assign$tmpm1146332801$0;
let $jscomp$logical$assign$tmpindexm1146332801$0;
($jscomp$logical$assign$tmpm1146332801$0 = a)
[$jscomp$logical$assign$tmpindexm1146332801$0 = x]
   ||
($jscomp$logical$assign$tmpm1146332801$0
[$jscomp$logical$assign$tmpindexm1146332801$0] = b);
`))
	c.Test(
		c.Srcs("a[x + 5 + 's'] &&= b"),
		c.Expected(
			`let $jscomp$logical$assign$tmpm1146332801$0;
let $jscomp$logical$assign$tmpindexm1146332801$0;
($jscomp$logical$assign$tmpm1146332801$0 = a)
[$jscomp$logical$assign$tmpindexm1146332801$0 = (x + 5 + 's')]
   &&
($jscomp$logical$assign$tmpm1146332801$0
[$jscomp$logical$assign$tmpindexm1146332801$0] = b);
`))
	c.Test(
		c.Srcs("foo[x] ??= bar[y]"),
		c.Expected(
			`let $jscomp$logical$assign$tmpm1146332801$0;
let $jscomp$logical$assign$tmpindexm1146332801$0;
($jscomp$logical$assign$tmpm1146332801$0 = foo)
[$jscomp$logical$assign$tmpindexm1146332801$0 = x]
   ??
($jscomp$logical$assign$tmpm1146332801$0
[$jscomp$logical$assign$tmpindexm1146332801$0] = bar[y]);
`))
}

// Java: NormalizeTest.testLogicalAssignmentNestedName()
func (c *NormalizeTest) TestLogicalAssignmentNestedName() {
	c.Test(
		c.Srcs("a ||= (b &&= (c ??= d));"), //
		c.Expected("a || (a = (b && (b = (c ?? (c = d)))));"))
}

// Java: NormalizeTest.logicalAssignmentNestedPropertyReference()
func (c *NormalizeTest) LogicalAssignmentNestedPropertyReference() {
	c.Test(
		`const foo = {}, bar = {};
foo.x ||= (foo.y &&= (bar.z ??= 'something'));
`,
		`const foo = {}; const bar = {};
let $jscomp$logical$assign$tmpm1146332801$0;
let $jscomp$logical$assign$tmpm1146332801$1;
let $jscomp$logical$assign$tmpm1146332801$2;
($jscomp$logical$assign$tmpm1146332801$2 = foo).x
   ||
($jscomp$logical$assign$tmpm1146332801$2.x
 = ($jscomp$logical$assign$tmpm1146332801$1 = foo).y
   &&
($jscomp$logical$assign$tmpm1146332801$1.y
 = ($jscomp$logical$assign$tmpm1146332801$0 = bar).z
   ?? ($jscomp$logical$assign$tmpm1146332801$0.z = 'something')));
`)
}

// Java: NormalizeTest.testLogicalAssignmentInConciseArrowFunction()
func (c *NormalizeTest) TestLogicalAssignmentInConciseArrowFunction() {
	c.Test(
		c.Srcs("const f = (a, b) => a.x ||= b;"),
		c.Expected(
			"const f = (a, b) => {\n"+
				"  let $jscomp$logical$assign$tmpm1146332801$0;\n"+
				"  return ($jscomp$logical$assign$tmpm1146332801$0 = a).x\n"+
				"      || ($jscomp$logical$assign$tmpm1146332801$0.x = b);\n"+
				"};\n"))

	c.Test(
		c.Srcs("const f = (map, k, v) => (map[k] ??= v);"),
		c.Expected(
			"const f = (map, k, v) => {\n"+
				"  let $jscomp$logical$assign$tmpm1146332801$0;\n"+
				"  let $jscomp$logical$assign$tmpindexm1146332801$0;\n"+
				"  return ($jscomp$logical$assign$tmpm1146332801$0 = map)\n"+
				"      [$jscomp$logical$assign$tmpindexm1146332801$0 = k]\n"+
				"      ??\n"+
				"      ($jscomp$logical$assign$tmpm1146332801$0\n"+
				"      [$jscomp$logical$assign$tmpindexm1146332801$0] = v);\n"+
				"};\n"))
}

// Java: NormalizeTest.testDuplicateVarInExterns()
func (c *NormalizeTest) TestDuplicateVarInExterns() {
	c.Test(
		c.Externs("var extern;"),
		c.Srcs("/** @suppress {duplicate} */ var extern = 3;"),
		c.Expected("var extern = 3;"))
}

// Java: NormalizeTest.testUnhandled()
func (c *NormalizeTest) TestUnhandled() {
	c.TestSame("var x = y = 1")
}

// Java: NormalizeTest.testFor()
func (c *NormalizeTest) TestFor() {
	// Verify assignments are extracted from the FOR init node.
	c.Test("for(a = 0; a < 2 ; a++) foo();", "a = 0; for(; a < 2 ; a++) foo()")
	// Verify vars are extracted from the FOR init node.
	c.Test("for(var a = 0; c < b ; c++) foo()", "var a = 0; for(; c < b ; c++) foo()")
	// Verify vars are extracted from the FOR init before the label node.
	c.Test("a:for(var a = 0; c < b ; c++) foo()", "var a = 0; a:for(; c < b ; c++) foo()")
	// Verify vars are extracted from the FOR init before the labels node.
	c.Test("a:b:for(var a = 0; c < b ; c++) foo()", "var a = 0; a:b:for(; c < b ; c++) foo()")
	// Verify block are properly introduced for ifs.
	c.Test("if(x) for(var a = 0; c < b ; c++) foo()", "if(x){var a = 0; for(; c < b ; c++) foo()}")
	// Any other expression.
	c.Test("for(init(); a < 2 ; a++) foo();", "init(); for(; a < 2 ; a++) foo()")
	// Verify destructuring var declarations are extracted.
	c.Test(
		"for (var [a, b] = [1, 2]; a < 2; a = b++) foo();",
		"var a; var b; [a, b] = [1, 2]; for (; a < 2; a = b++) foo();")
}

// Java: NormalizeTest.testForIn1()
func (c *NormalizeTest) TestForIn1() {
	c.IgnoreWarnings(DiagnosticGroups_MISSING_PROPERTIES)
	// Verify nothing happens with simple for-in
	c.TestSame("for(a in b) foo();")
	// Verify vars are extracted from the FOR-IN node.
	c.Test("for(var a in b) foo()", "var a; for(a in b) foo()")
	// Verify vars are extracted from the FOR init before the label node.
	c.Test("a:for(var a in b) foo()", "var a; a:for(a in b) foo()")
	// Verify vars are extracted from the FOR init before the labels node.
	c.Test("a:b:for(var a in b) foo()", "var a; a:b:for(a in b) foo()")
	// Verify block are properly introduced for ifs.
	c.Test("if (x) for(var a in b) foo()", "if (x) { var a; for(a in b) foo() }")
	// Verify names in destructuring declarations are individually declared.
	c.Test(
		c.Externs(NewTestExternsBuilder().AddIterable().AddString().Build()),
		c.Srcs("for (var [a, b] in c) foo();"),
		c.Expected("var a; var b; for ([a, b] in c) foo();"))
	c.Test("for (var {a, b} in c) foo();", "var a; var b; for ({a: a, b: b} in c) foo();")
}

// Java: NormalizeTest.testForIn2()
func (c *NormalizeTest) TestForIn2() {
	c.SetExpectParseWarningsInThisTest()
	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT5)
	// Verify vars are extracted from the FOR-IN node.
	c.Test("for(var a = foo() in b) foo()", "var a = foo(); for(a in b) foo()")
}

// Java: NormalizeTest.testForOf()
func (c *NormalizeTest) TestForOf() {
	c.IgnoreWarnings(DiagnosticGroups_GLOBALLY_MISSING_PROPERTIES)
	// Verify nothing happens with simple for-of
	c.TestSame("for (a of b) foo();")
	// Verify vars are extracted from the FOR-OF node.
	c.Test("for (var a of b) foo()", "var a; for (a of b) foo()")
	// Verify vars are extracted from the FOR init before the label node.
	c.Test("a:for (var a of b) foo()", "var a; a: for (a of b) foo()")
	// Verify vars are extracted from the FOR init before the labels node.
	c.Test("a: b: for (var a of b) foo()", "var a; a: b: for (a of b) foo()")
	// Verify block are properly introduced for ifs.
	c.Test("if (x) for (var a of b) foo()", "if (x) { var a; for (a of b) foo() }")
	// Verify names in destructuring declarations are individually declared.
	c.Test("for (var [a, b] of c) foo();", "var a; var b; for ([a, b] of c) foo();")
	c.Test("for (var {a, b} of c) foo();", "var a; var b; for ({a: a, b: b} of c) foo();")
}

// Java: NormalizeTest.testForAwaitOf()
func (c *NormalizeTest) TestForAwaitOf() {
	c.IgnoreWarnings(DiagnosticGroups_GLOBALLY_MISSING_PROPERTIES)
	// Verify nothing happens with simple for-await-of
	c.TestSame("async () => { for await (a of b) foo(); }")
	// Verify vars are extracted from the FOR-AWAIT-OF node.
	c.Test(
		"async () => { for await (var a of b) foo() }",
		"async () => { var a; for await (a of b) foo() }")
	// Verify vars are extracted from the FOR init before the label node.
	c.Test(
		"async () => { a:for await (var a of b) foo() }",
		"async () => { var a; a: for await (a of b) foo() }")
	// Verify vars are extracted from the FOR init before the labels node.
	c.Test(
		"async () => { a: b: for await (var a of b) foo() }",
		"async () => { var a; a: b: for await (a of b) foo() }")
	// Verify block are properly introduced for ifs.
	c.Test(
		"async () => { if (x) for await (var a of b) foo() }",
		"async () => { if (x) { var a; for await (a of b) foo() } }")
	// Verify names in destructuring declarations are individually declared.
	c.Test(
		"async () => { for await (var [a, b] of c) foo(); }",
		"async () => { var a; var b; for await ([a, b] of c) foo(); }")
	c.Test(
		"async () => { for await (var {a, b} of c) foo(); }",
		"async () => { var a; var b; for await ({a: a, b: b} of c) foo(); }")
}

// Java: NormalizeTest.testWhile()
func (c *NormalizeTest) TestWhile() {
	// Verify while loops are converted to FOR loops.
	c.Test("while(c < b) foo()", "for(; c < b;) foo()")
}

// Java: NormalizeTest.testMoveFunctions1()
func (c *NormalizeTest) TestMoveFunctions1() {
	c.Test(
		"function f() { if (x) return; foo(); function foo() {} }",
		"function f() {function foo() {} if (x) return; foo(); }")
	c.Test(
		`function f() {
  function foo() {}
  if (x) return;
  foo();
  function bar() {}
}
`,
		`function f() {
  function foo() {}
  function bar() {}
  if (x) return;
  foo();
}
`)
}

// Java: NormalizeTest.testMoveFunctions2()
func (c *NormalizeTest) TestMoveFunctions2() {
	c.TestSame("function f() { function foo() {} }")
	c.TestSame("function f() { f(); {function bar() {}}}")
	c.TestSame("function f() { f(); if (true) {function bar() {}}}")
}

// Java: NormalizeTest.inFunction(String)
func NormalizeTest_InFunction(code string) string {
	return "(function(){" + code + "})"
}

// Java: NormalizeTest.c.TestSameInFunction(String)
func (c *NormalizeTest) TestSameInFunction(code string) {
	c.TestSame(NormalizeTest_InFunction(code))
}

// Java: NormalizeTest.c.TestInFunction(String, String)
func (c *NormalizeTest) TestInFunction(code string, expected string) {
	c.Test(NormalizeTest_InFunction(code), NormalizeTest_InFunction(expected))
}

// Java: NormalizeTest.testNormalizeFunctionDeclarations()
func (c *NormalizeTest) TestNormalizeFunctionDeclarations() {
	c.TestSame("function f() {}")
	c.TestSame("function f() {}; var g = f;")
	c.TestSame("var f = function () {}")
	c.Test("var f = function f() {}", "var f = function f$jscomp$1() {}")
	c.TestSame("var f = function g() {}")
	c.TestSame("var f = function(){}; var g = f;")
	c.TestSame("{function g() {}}")
	c.TestSame("if (function g() {}) {}")
	c.TestSame("if (true) {function g() {}}")
	c.TestSame("if (true) {} else {function g() {}}")
	c.TestSame("switch (function g() {}) {}")
	c.TestSame("switch (1) { case 1: function g() {}}")
	c.TestSame("if (true) {function g() {} function h() {}}")
	c.TestSameInFunction("function f() {}")
	c.TestSameInFunction("f(); {function g() {}}")
	c.TestSameInFunction("f(); if (true) {function g() {}}")
	c.TestSameInFunction("if (true) {} else {function g() {}}")
}

// Java: NormalizeTest.testMakeLocalNamesUnique()
func (c *NormalizeTest) TestMakeLocalNamesUnique() {
	c.IgnoreWarnings(DiagnosticGroups_GLOBALLY_MISSING_PROPERTIES)
	c.DisableCompareJsDoc()
	// Verify global names are untouched.
	c.TestSame("var a;")
	// Verify global names are untouched.
	c.TestSame("a;")
	// Local names are made unique.
	c.Test("var a;function foo(a){var b;a}", "var a;function foo(a$jscomp$1){var b;a$jscomp$1}")
	c.Test(
		"var a;function foo(){var b;a}function boo(){var b;a}",
		"var a;function foo(){var b;a}function boo(){var b$jscomp$1;a}")
	c.Test(
		"function foo(a){var b} function boo(a){var b}",
		"function foo(a){var b} function boo(a$jscomp$1){var b$jscomp$1}")
	// Verify function expressions are renamed.
	c.Test(
		"var a = function foo(){foo()};var b = function foo(){foo()};",
		"var a = function foo(){foo()};var b = function foo$jscomp$1(){foo$jscomp$1()};")
	// Verify catch exceptions names are made unique
	c.TestSame("try { } catch(e) {e;}")
	c.Test(
		"try { } catch(e) {e;}; try { } catch(e) {e;}",
		"try { } catch(e) {e;}; try { } catch(e$jscomp$1) {e$jscomp$1;}")
	c.Test(
		"try { } catch(e) {e; try { } catch(e) {e;}};",
		"try { } catch(e) {e; try { } catch(e$jscomp$1) {e$jscomp$1;} }; ")
	// Verify the 1st global redefinition of extern definition is not removed.
	c.TestSame("/** @suppress {duplicate} */ var window;")
	// Verify the 2nd global redefinition of extern definition is removed.
	c.Test(
		"/** @suppress {duplicate} */ var window; /** @suppress {duplicate} */ var window;",
		"/** @suppress {duplicate} */ var window;")
	// Verify local masking extern made unique.
	c.Test("function f() {var window}", "function f() {var window$jscomp$1}")
	// Verify import * as <alias> is renamed.
	c.Test(
		c.Srcs("let a = 5;", "import * as a from './a.js'; const TAU = 2 * a.PI;"),
		c.Expected(
			"let a = 5;", "import * as a$jscomp$1 from './a.js'; const TAU = 2 * a$jscomp$1.PI"))
	// Verify exported and imported names are untouched.
	c.Test(
		c.Srcs("var a;", "let a; export {a as a};"),
		c.Expected("var a;", "let a$jscomp$1; export {a$jscomp$1 as a};"))
	// Verify same local names in 2 different files get different new names
	c.Test(
		c.Srcs(
			"function foo() {var a; a;}",
			"function bar() {let a; let a$jscomp$1; a + a$jscomp$1;}"),
		c.Expected(
			"function foo() {var a; a;}",
			`function bar() {let a$jscomp$1; let a$jscomp$1$jscomp$1; a$jscomp$1 +
 a$jscomp$1$jscomp$1;}
`))
	c.Test(
		c.Srcs("var a;", "import {a as a} from './foo.js'; let b = a;"),
		c.Expected("var a;", "import {a as a$jscomp$1} from './foo.js'; let b = a$jscomp$1;"))
	c.Test(
		`function one() {
  var index = 1;
  return index;
}
function two() {
  var index = 2;
  var id = (index += 3).toString();
}
`,
		`function one() {
  var index = 1;
  return index;
}
function two() {
  var index$jscomp$1 = 2;
  var id = (index$jscomp$1 = index$jscomp$1 + 3).toString();
}
`)
}

// Java: NormalizeTest.testMakeParamNamesUnique()
func (c *NormalizeTest) TestMakeParamNamesUnique() {
	c.IgnoreWarnings(DiagnosticGroups_GLOBALLY_MISSING_PROPERTIES)
	c.Test(
		`function f(x) { x; }
function g(x) { x; }`,
		`function f(x) { x; }
function g(x$jscomp$1) { x$jscomp$1; }`)
	c.Test(
		`function f(x) { x; }
function g(...x) { x; }`,
		`function f(x) { x; }
function g(...x$jscomp$1) { x$jscomp$1; }`)
	c.Test(
		`function f(x) { x; }
function g({x: x}) { x; }`,
		`function f(x) { x; }
function g({x: x$jscomp$1}) { x$jscomp$1; }`)
	c.Test(
		`function f(x) { x; }
function g({x}) { x; }`,
		`function f(x) { x; }
function g({x: x$jscomp$1}) { x$jscomp$1; }`)
	c.Test(
		`function f(x) { x; }
function g({y: {x}}) { x; }`,
		`function f(x) { x; }
function g({y: {x: x$jscomp$1}}) { x$jscomp$1; }`)
}

// Java: NormalizeTest.testNormalize_createsUniqueNamesInFunctionBody_avoidsShadowingWithDefaultParams()
func (c *NormalizeTest) TestNormalize_createsUniqueNamesInFunctionBody_avoidsShadowingWithDefaultParams() {
	c.Test(
		`function x() {}
var y = 1;
function f(z=x, w=y) {
  let x = y;
  var y = 3;
  return w;
}
`,
		`function x() {}
var y = 1;
function f(z=x, w=y) {
  let x$jscomp$1 = y$jscomp$1;
  var y$jscomp$1 = 3;
  return w;
}
`)
}

// Java: NormalizeTest.testNoRenameParamNames()
func (c *NormalizeTest) TestNoRenameParamNames() {
	c.IgnoreWarnings(DiagnosticGroups_GLOBALLY_MISSING_PROPERTIES)
	c.TestSame("function f(x) { x; }")
	c.TestSame("function f(...x) { x; }")
	c.TestSame("function f({x: x}) { x; }")
	c.Test("function f({x}) { x; }", "function f({x: x}) { x; }")
	c.Test("function f({y: {x}}) { x; }", "function f({y: {x: x}}) { x; }")
}

// Java: NormalizeTest.testRemoveDuplicateVarDeclarations_destructuringArrayDeclaration()
func (c *NormalizeTest) TestRemoveDuplicateVarDeclarations_destructuringArrayDeclaration() {
	c.Test(
		"function f() { var a = 3; var [a] = [4]; }", //
		"function f() { var a = 3; [a]= [4]; }")
}

// Java: NormalizeTest.testRemoveDuplicateVarDeclarations_destructuringObjectDeclaration()
func (c *NormalizeTest) TestRemoveDuplicateVarDeclarations_destructuringObjectDeclaration() {
	c.Test(
		"function f() { var a = 3; var {a} = {a: 4}; }", //
		"function f() { var a = 3; ({a}= {a: 4});}")
}

// Java: NormalizeTest.testRemoveDuplicateVarDeclarations1()
func (c *NormalizeTest) TestRemoveDuplicateVarDeclarations1() {
	c.Test("function f() { var a; var a }", "function f() { var a; }")
	c.Test("function f() { var a = 1; var a = 2 }", "function f() { var a = 1; a = 2 }")
	// second declaration not assigned to any rhs, empty reference removed
	c.Test("function f() { var a = 1; var a; }", "function f() { var a = 1;}")
	// second declaration is a destructuring declaration; converted to assignment
	c.Test("function f() { var a = 1; var [a] = [2]; }", "function f() { var a = 1; [a]= [2]; }")
	// this should be an error in the parser but isn't
	c.Test("function f() { var a = 1; const a =2; }", "function f() { var a = 1; const a = 2 }")
	c.Test("var a = 1; function f(){ var a = 2 }", "var a = 1; function f(){ var a$jscomp$1 = 2 }")
	c.Test(
		"function f() { var a = 1; label1:var a = 2 }",
		"function f() { var a = 1; label1:{a = 2}}")
	c.Test("function f() { var a = 1; label1:var a }", "function f() { var a = 1; label1:{} }")
	c.Test("function f() { var a = 1; for(var a in b); }", "function f() { var a = 1; for(a in b);}")
}

// Java: NormalizeTest.testRemoveDuplicateVarDeclarations2()
func (c *NormalizeTest) TestRemoveDuplicateVarDeclarations2() {
	c.Test(
		"var e = 1; function f(){ try {} catch (e) {} var e = 2 }",
		"var e = 1; function f(){ try {} catch (e$jscomp$2) {} var e$jscomp$1 = 2 }")
}

// Java: NormalizeTest.testRemoveDuplicateVarDeclarations3()
func (c *NormalizeTest) TestRemoveDuplicateVarDeclarations3() {
	c.IgnoreWarnings(
		TypeCheck_FUNCTION_MASKS_VARIABLE,
		TypeValidator_TYPE_MISMATCH_WARNING,
		TypeValidator_DUP_VAR_DECLARATION)
	c.Test("var f = 1; function f(){}", "f = 1; function f(){}")
	c.Test("var f; function f(){}", "function f(){}")
	c.Test("function f(){} var f = 1;", "function f(){} f = 1;")
	c.Test("function f(){} var f;", "function f(){}")
	// TODO(johnlenz): Do we need to handle this differently for "third_party"
	// mode? Remove the previous function definitions?
	c.TestSame("function f(){} function f(){}")
}

// It's important that we not remove this var completely. See
// http://blickly.github.io/closure-compiler-issues/#290
// Java: NormalizeTest.testRemoveDuplicateVarDeclarations4()
func (c *NormalizeTest) TestRemoveDuplicateVarDeclarations4() {
	c.DisableCompareJsDoc()
	c.TestSame("if (!Arguments) { /** @suppress {duplicate} */ var Arguments = {}; }")
}

// If there are multiple duplicates, it's okay to remove all but the first.
// Java: NormalizeTest.testRemoveDuplicateVarDeclarations5()
func (c *NormalizeTest) TestRemoveDuplicateVarDeclarations5() {
	c.Test("var Arguments = {}; var Arguments = {};", "var Arguments = {}; Arguments = {};")
}

// Java: NormalizeTest.testRemoveVarDeclarationDuplicatesParam1()
func (c *NormalizeTest) TestRemoveVarDeclarationDuplicatesParam1() {
	c.Test(
		"function f(x) { alert(x); var x = 0; alert(x); }",
		"function f(x) { alert(x);     x = 0; alert(x); }")
}

// Java: NormalizeTest.testRemoveVarDeclarationDuplicatesParam2()
func (c *NormalizeTest) TestRemoveVarDeclarationDuplicatesParam2() {
	c.Test(
		"function f(x) { alert(x); var x; alert(x); }",
		"function f(x) { alert(x);        alert(x); }")
}

// Java: NormalizeTest.testRenamingConstants()
func (c *NormalizeTest) TestRenamingConstants() {
	c.TestSame("var ACONST = 4; var b = ACONST;")
	c.Test("var a, ACONST = 4;var b = ACONST;", "var a; var ACONST = 4; var b = ACONST;")
	c.TestSame("var ACONST; ACONST = 4; var b = ACONST;")
	c.TestSame("var ACONST = new Foo(); var b = ACONST;")
	c.TestSame("/** @const */ var aa; aa = 1;")
}

// Java: NormalizeTest.testSkipRenamingExterns()
func (c *NormalizeTest) TestSkipRenamingExterns() {
	c.Test(
		c.Externs("var EXTERN; var ext; ext.FOO;"),
		c.Srcs("var b = EXTERN; var c = ext.FOO"),
		c.Expected("var b = EXTERN; var c = ext.FOO"))
}

// Java: NormalizeTest.testIssue166e()
func (c *NormalizeTest) TestIssue166e() {
	c.Test(
		"var e = 2; try { throw 1 } catch(e) {}",
		"var e = 2; try { throw 1 } catch(e$jscomp$1) {}")
}

// Java: NormalizeTest.testIssue166f()
func (c *NormalizeTest) TestIssue166f() {
	c.Test(
		`function a() {
  var e = 2;
  try { throw 1 } catch(e) {}
}
`,
		`function a() {
  var e = 2;
  try { throw 1 } catch(e$jscomp$1) {}
}
`)
}

// Java: NormalizeTest.testIssue166g()
func (c *NormalizeTest) TestIssue166g() {
	c.Test(
		`function a() {
  try { throw 1 } catch(e) {}
  var e = 2;
}
`,
		`function a() {
  try { throw 1 } catch(e$jscomp$1) {}
  var e = 2;
}
`)
}

// Java: NormalizeTest.testLetsInSeparateBlocks()
func (c *NormalizeTest) TestLetsInSeparateBlocks() {
	c.Test(
		`if (x) {
  let e;
  alert(e);
}
if (y) {
  let e;
  alert(e);
}
`,
		`if (x) {
  let e;
  alert(e);
}
if (y) {
  let e$jscomp$1;
  alert(e$jscomp$1);
}
`)
}

// Java: NormalizeTest.testCatchesInSeparateBlocks()
func (c *NormalizeTest) TestCatchesInSeparateBlocks() {
	c.Test(
		`if (x) {
  try {
    throw 1;
  } catch (e) {
    alert(e);
  }
}
if (y) {
  try {
    throw 2;
  } catch (e) {
    alert(e);
  }
}
`,
		`if (x) {
  try {
    throw 1;
  } catch (e) {
    alert(e);
  }
}
if (y) {
  try {
    throw 2;
  } catch (e$jscomp$1) {
    alert(e$jscomp$1);
  }
}
`)
}

// Java: NormalizeTest.testDeclInCatchBlock()
func (c *NormalizeTest) TestDeclInCatchBlock() {
	c.Test(
		`var x;
try {
} catch (e) {
  let x;
}
`,
		`var x;
try {
} catch (e) {
  let x$jscomp$1
}
`)
}

// Java: NormalizeTest.testIssue()
func (c *NormalizeTest) TestIssue() {
	c.AllowExternsChanges()
	c.Test(c.Externs("var a,b,c; var a,b"), c.Srcs("a(), b()"), c.Expected("a(), b()"))
}

// Java: NormalizeTest.testIsConstant()
func (c *NormalizeTest) TestIsConstant() {
	c.TestSame("var CONST = 3; var b = CONST;")
	n := c.GetLastCompiler().GetRoot()
	constantNodes := c.FindNodesWithProperty(n, NormalizeTest_IsConstantName)
	AssertThat(constantNodes).HasSize(2)
	for hasProp := range constantNodes.All() {
		AssertThat(hasProp.GetString()).IsEqualTo("CONST")
	}
}

// Test that the name nodes of function expressions are unconditionally marked as const even if
// the LHS declaration is a var
// Java: NormalizeTest.testRHSFunctionExpressionNameNodeIsConstant()
func (c *NormalizeTest) TestRHSFunctionExpressionNameNodeIsConstant() {
	c.TestSame("var someNonConstVar = function foo() {};")
	n := c.GetLastCompiler().GetRoot()
	constantNodes := c.FindNodesWithProperty(n, NormalizeTest_IsConstantName)
	AssertThat(constantNodes).HasSize(1)
	AssertThat(constantNodes.ToSlice()[0].GetString()).IsEqualTo("foo")
}

// Test that the name nodes of function expressions are unconditionally marked as const if the LHS
// is a const
// Java: NormalizeTest.testRHSFunctionExpressionNameNodeIsConstant1()
func (c *NormalizeTest) TestRHSFunctionExpressionNameNodeIsConstant1() {
	c.TestSame("const someConstVar = function foo() {};")
	n := c.GetLastCompiler().GetRoot()
	constantNodes := c.FindNodesWithProperty(n, NormalizeTest_IsConstantName)
	AssertThat(constantNodes).HasSize(2) // {someConstVar, foo}
	itr := constantNodes.ToSlice()
	AssertThat(itr[0].GetString()).IsEqualTo("foo")
	AssertThat(itr[1].GetString()).IsEqualTo("someConstVar")
}

// Test that the name nodes of function expressions are unconditionally marked as const if the LHS
// is a const
// Java: NormalizeTest.testRHSFunctionExpressionNameNodeIsConstant2()
func (c *NormalizeTest) TestRHSFunctionExpressionNameNodeIsConstant2() {
	c.TestSame("const someConstVar = function foo() { foo(); };")
	n := c.GetLastCompiler().GetRoot()
	constantNodes := c.FindNodesWithProperty(n, NormalizeTest_IsConstantName)
	AssertThat(constantNodes).HasSize(3) // {someConstVar, foo, foo}
	itr := constantNodes.ToSlice()
	AssertThat(itr[0].GetString()).IsEqualTo("foo")
	AssertThat(itr[1].GetString()).IsEqualTo("foo")
	AssertThat(itr[2].GetString()).IsEqualTo("someConstVar")
}

// Test that the name nodes of function expressions are unconditionally marked as const
// irrespective of whether or not the LHS is a const
// Java: NormalizeTest.testRHSFunctionExpressionNameNodeIsConstant3()
func (c *NormalizeTest) TestRHSFunctionExpressionNameNodeIsConstant3() {
	c.TestSame(
		`const someConstVar = function foo() { foo(); };
// this call to something undefined but having the same name foo is not marked const
foo();
`)
	n := c.GetLastCompiler().GetRoot()
	constantNodes := c.FindNodesWithProperty(n, NormalizeTest_IsConstantName)
	AssertThat(constantNodes).HasSize(3) // {someConstVar, foo, foo}
	itr := constantNodes.ToSlice()
	AssertThat(itr[0].GetString()).IsEqualTo("foo")
	AssertThat(itr[1].GetString()).IsEqualTo("foo")
	AssertThat(itr[2].GetString()).IsEqualTo("someConstVar")
}

// Test that the name nodes of function expressions are unconditionally marked as const
// irrespective of whether or not the LHS is a const
// Java: NormalizeTest.testRHSFunctionExpressionNameNodeIsConstant4()
func (c *NormalizeTest) TestRHSFunctionExpressionNameNodeIsConstant4() {
	c.TestSame(
		`let someNonConstVar = function foo() { foo(); };
// this call to something undefined but having the same name foo is not marked const
foo();
`)
	n := c.GetLastCompiler().GetRoot()
	constantNodes := c.FindNodesWithProperty(n, NormalizeTest_IsConstantName)
	AssertThat(constantNodes).HasSize(2) // {foo, foo}
	itr := constantNodes.ToSlice()
	AssertThat(itr[0].GetString()).IsEqualTo("foo")
	AssertThat(itr[1].GetString()).IsEqualTo("foo")
}

// Test that the name nodes of function expressions are unconditionally marked as const
// irrespective of whether or not the LHS is a const
// Java: NormalizeTest.testFunctionExpressionNameNodeIsConstant()
func (c *NormalizeTest) TestFunctionExpressionNameNodeIsConstant() {
	c.TestSame("use(function foo(i) { foo(i-1);});")
	n := c.GetLastCompiler().GetRoot()
	constantNodes := c.FindNodesWithProperty(n, NormalizeTest_IsConstantName)
	AssertThat(constantNodes).HasSize(2) // {foo, foo}
	itr := constantNodes.ToSlice()
	AssertThat(itr[0].GetString()).IsEqualTo("foo")
	AssertThat(itr[1].GetString()).IsEqualTo("foo")
}

// Test that the name nodes of function expressions are unconditionally marked as const
// irrespective of whether or not the LHS is a const
// Java: NormalizeTest.testFunctionExpressionNameNodeIsConstant2()
func (c *NormalizeTest) TestFunctionExpressionNameNodeIsConstant2() {
	c.TestSame(
		`use(function foo(i) { foo(i-1);});
// this call to something undefined but having the same name foo is not marked const
foo();
`)
	n := c.GetLastCompiler().GetRoot()
	constantNodes := c.FindNodesWithProperty(n, NormalizeTest_IsConstantName)
	AssertThat(constantNodes).HasSize(2) // {foo, foo}
	itr := constantNodes.ToSlice()
	AssertThat(itr[0].GetString()).IsEqualTo("foo")
	AssertThat(itr[1].GetString()).IsEqualTo("foo")
}

// Java: NormalizeTest.testIsConstantByDestructuring()
func (c *NormalizeTest) TestIsConstantByDestructuring() {
	c.Test(
		"const {CONST} = {CONST:3}; let b = CONST;",
		"const {CONST: CONST} = {CONST:3}; let b = CONST;")
	n := c.GetLastCompiler().GetRoot()
	constantNodes := c.FindNodesWithProperty(n, NormalizeTest_IsConstantName)
	AssertThat(constantNodes).HasSize(2)
	for hasProp := range constantNodes.All() {
		AssertThat(hasProp.GetString()).IsEqualTo("CONST")
	}
}

// Java: NormalizeTest.testIsConstantByDestructuringWithDefault()
func (c *NormalizeTest) TestIsConstantByDestructuringWithDefault() {
	c.IgnoreWarnings(DiagnosticGroups_MISSING_PROPERTIES)
	c.Test("const {CONST = 3} = {}; var b = CONST;", "const {CONST: CONST = 3} = {}; var b = CONST;")
	n := c.GetLastCompiler().GetRoot()
	constantNodes := c.FindNodesWithProperty(n, NormalizeTest_IsConstantName)
	AssertThat(constantNodes).HasSize(2)
	for hasProp := range constantNodes.All() {
		AssertThat(hasProp.GetString()).IsEqualTo("CONST")
	}
}

// Java: NormalizeTest.testPropertyIsConstantIfMatchesConstantName()
func (c *NormalizeTest) TestPropertyIsConstantIfMatchesConstantName() {
	// verify that the /** @const */ 'other' doesn't accidentally cause the string key in
	// {'other: 4'} to be marked const
	c.TestSame("var a = {other: 4}; /** @const */ var other = 5;")
	n := c.GetLastCompiler().GetRoot()
	constantNodes := c.FindNodesWithProperty(n, NormalizeTest_IsConstantName)
	AssertThat(constantNodes).HasSize(1)
	for hasProp := range constantNodes.All() {
		AssertThat(hasProp.GetString()).IsEqualTo("other")
	}
}

// Java: NormalizeTest.testShadowFunctionName()
func (c *NormalizeTest) TestShadowFunctionName() {
	c.IgnoreWarnings(DiagnosticGroups_GLOBALLY_MISSING_PROPERTIES)
	c.Test(
		`function f() {
  var f = 'test';
  console.log(f);
}
`,
		`function f() {
  var f$jscomp$1 = 'test';
  console.log(f$jscomp$1);
}
`)
}

// Java: NormalizeTest.isConstantName(Node)
func NormalizeTest_IsConstantName(n rhino.Node) bool {
	return n.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME)
}

// Java: NormalizeTest.findNodesWithProperty(Node, Predicate<Node>)
func (c *NormalizeTest) FindNodesWithProperty(root rhino.Node, prop func(rhino.Node) bool) *jutil.LinkedHashSet[rhino.Node] {
	set := jutil.NewLinkedHashSet[rhino.Node]()

	NodeTraversal_Builder_().
		SetCompiler(c.GetLastCompiler()).
		SetCallbackFunc(
			func(t *NodeTraversal, node rhino.Node, parent rhino.Node) {
				if prop(node) {
					set.Add(node)
				}
			}).
		Traverse(root)
	return set
}

// normalizeTest_renamingConstantPropertiesTester is the anonymous
// CompilerTestCase subclass of testRenamingConstantProperties.
type normalizeTest_renamingConstantPropertiesTester struct {
	CompilerTestCaseBase
}

func newNormalizeTest_renamingConstantPropertiesTester() *normalizeTest_renamingConstantPropertiesTester {
	c := &normalizeTest_renamingConstantPropertiesTester{}
	c.InitCompilerTestCaseBase(c)
	return c
}

// Java: NormalizeTest.testRenamingConstantProperties().new CompilerTestCase() {...}.getProcessor(Compiler)
func (c *normalizeTest_renamingConstantPropertiesTester) GetProcessor(compiler *Compiler) CompilerPass {
	optimizer := NewPhaseOptimizer(compiler, nil)
	optimizer.AddOneTimePass(
		CompilerTestCase_MakePassFactory("es6NormalizeClasses", func(c AbstractCompiler) CompilerPass { return NewEs6NormalizeClasses(c) }))
	optimizer.AddOneTimePass(
		CompilerTestCase_MakePassFactory(
			"inlineAndCollapseProperties",
			func(comp AbstractCompiler) CompilerPass {
				return InlineAndCollapseProperties_Builder_(compiler).
					SetPropertyCollapseLevel(CompilerOptions_PropertyCollapseLevel_ALL).
					SetChunkOutputType(CompilerOptions_ChunkOutputType_GLOBAL_NAMESPACE).
					SetHaveModulesBeenRewritten(false).
					SetModuleResolutionMode(ModuleLoader_ResolutionMode_BROWSER).
					Build()
			}))
	return optimizer
}

// Java: NormalizeTest.testRenamingConstantProperties()
func (c *NormalizeTest) TestRenamingConstantProperties() {
	// In order to detect that foo.BAR is a constant, we need collapse
	// properties to run first so that we can tell if the initial value is
	// non-null and immutable. The Normalize pass doesn't modify the code
	// in these examples, it just infers const-ness of some variables, so
	// we call enableNormalize to make the Normalize.VerifyConstants pass run.

	// TODO(johnlenz): fix this so it is just another test case.
	tester := newNormalizeTest_renamingConstantPropertiesTester()
	tester.JUnitSetT(c.T())

	tester.SetUp()
	tester.EnableNormalize()
	tester.Test(
		"var a={}; a.ACONST = 4;var b = 1; b = a.ACONST;",
		"var a$ACONST = 4; var b = 1; b = a$ACONST;")
	tester.Test(
		"var a={b:{}}; a.b.ACONST = 4;var b = 1; b = a.b.ACONST;",
		"var a$b$ACONST = 4;var b = 1; b = a$b$ACONST;")
	tester.Test("var a = {FOO: 1};var b = 1; b = a.FOO;", "var a$FOO = 1; var b = 1; b = a$FOO;")
	tester.TestSame(
		tester.Externs("var EXTERN; var ext; ext.FOO;"), tester.Srcs("var b = EXTERN; var c = ext.FOO"))
	tester.Test(
		"var a={}; a.ACONST = 4; var b = 1; b = a.ACONST;",
		"var a$ACONST = 4; var b = 1; b = a$ACONST;")
	tester.Test(
		"var a = {}; function foo() { var d = a.CONST; }; (function(){a.CONST=4})();",
		"var a$CONST;function foo(){var d = a$CONST;}; (function(){a$CONST = 4})();")
	tester.Test(
		"var a = {}; a.ACONST = new Foo(); var b = 1; b = a.ACONST;",
		"var a$ACONST = new Foo(); var b = 1; b = a$ACONST;")
	tester.TearDown()
}

// Java: NormalizeTest.testFunctionBlock1()
func (c *NormalizeTest) TestFunctionBlock1() {
	c.Test("() => 1;", "() => { return 1; }")
}

// Java: NormalizeTest.testFunctionBlock2()
func (c *NormalizeTest) TestFunctionBlock2() {
	c.Test("var args = 1; var foo = () => args;", "var args = 1; var foo = () => { return args; }")
}

// Java: NormalizeTest.testBlocklessArrowFunction_withinArgs_getBlocks()
func (c *NormalizeTest) TestBlocklessArrowFunction_withinArgs_getBlocks() {
	// disable type checking to prevent `property map never defined` errors.
	c.DisableTypeCheck()
	c.DisableTypeInfoValidation()
	c.Test(
		"function sortAndConcatParams(params) { // arrow fn body missing block\n  return [...params].map(((k) => `k`));}\n",
		"function sortAndConcatParams(params) { // gets block {}\n  return [...params].map(((k) => { return `k`; }));}\n")
}

// Java: NormalizeTest.testArrowFunctionInFunction()
func (c *NormalizeTest) TestArrowFunctionInFunction() {
	c.Test(
		`function foo() {
  var x = () => 1;
  return x();
}
`,
		`function foo() {
  var x = () => { return 1; };
  return x();
}
`)
}

// Java: NormalizeTest.testES6ShorthandPropertySyntax01()
func (c *NormalizeTest) TestES6ShorthandPropertySyntax01() {
	c.Test("obj = {x, y};", "obj = {x: x, y: y}")
}

// Java: NormalizeTest.testES6ShorthandPropertySyntax02()
func (c *NormalizeTest) TestES6ShorthandPropertySyntax02() {
	c.Test("var foo = {x, y};", "var foo = {x: x, y: y}")
}

// Java: NormalizeTest.testES6ShorthandPropertySyntax03()
func (c *NormalizeTest) TestES6ShorthandPropertySyntax03() {
	c.Test(
		`function foo(a, b, c) {
  return {
    a,
    b,
    c
  };
}
`,
		`function foo(a, b, c) {
  return {
    a: a,
    b: b,
    c: c
  };
}
`)
}

// Java: NormalizeTest.testES6ShorthandPropertySyntax04()
func (c *NormalizeTest) TestES6ShorthandPropertySyntax04() {
	c.Test("var foo = {x};", "var foo = {x: x}")
}

// Java: NormalizeTest.testES6ShorthandPropertySyntax05()
func (c *NormalizeTest) TestES6ShorthandPropertySyntax05() {
	c.IgnoreWarnings(DiagnosticGroups_GLOBALLY_MISSING_PROPERTIES)
	c.Test("var {a = 5} = obj;", "var a; ({a = 5} = obj);")
}

// Java: NormalizeTest.testES6ShorthandPropertySyntax06()
func (c *NormalizeTest) TestES6ShorthandPropertySyntax06() {
	c.IgnoreWarnings(DiagnosticGroups_GLOBALLY_MISSING_PROPERTIES)
	c.Test("var {a = 5, b = 3} = obj;", "var a; var b; ({a = 5, b = 3} = obj);")
}

// Java: NormalizeTest.testES6ShorthandPropertySyntax07()
func (c *NormalizeTest) TestES6ShorthandPropertySyntax07() {
	c.IgnoreWarnings(DiagnosticGroups_GLOBALLY_MISSING_PROPERTIES)
	c.Test("var {a: a = 5, b = 3} = obj;", "var a; var b; ({a: a = 5, b: b = 3} = obj);")
}

// Java: NormalizeTest.testES6ShorthandPropertySyntax08()
func (c *NormalizeTest) TestES6ShorthandPropertySyntax08() {
	c.IgnoreWarnings(DiagnosticGroups_GLOBALLY_MISSING_PROPERTIES)
	c.Test("var {a, b} = obj;", "var a; var b; ({a, b} = obj);")
}

// Java: NormalizeTest.testES6ShorthandPropertySyntax09()
func (c *NormalizeTest) TestES6ShorthandPropertySyntax09() {
	c.IgnoreWarnings(DiagnosticGroups_GLOBALLY_MISSING_PROPERTIES)
	c.Test("({a = 5} = obj);", "({a: a = 5} = obj);")
}

// Java: NormalizeTest.testES6ShorthandPropertySyntax10()
func (c *NormalizeTest) TestES6ShorthandPropertySyntax10() {
	c.TestSame("function f(a = 5) {}")
}

// Java: NormalizeTest.testES6ShorthandPropertySyntax11()
func (c *NormalizeTest) TestES6ShorthandPropertySyntax11() {
	c.TestSame("[a = 5] = obj;")
}

// Java: NormalizeTest.testES6ShorthandPropertySyntax12()
func (c *NormalizeTest) TestES6ShorthandPropertySyntax12() {
	c.IgnoreWarnings(DiagnosticGroups_GLOBALLY_MISSING_PROPERTIES)
	c.TestSame("({a: a = 5} = obj)")
}

// Java: NormalizeTest.testES6ShorthandPropertySyntax13()
func (c *NormalizeTest) TestES6ShorthandPropertySyntax13() {
	c.TestSame("({['a']: a = 5} = obj);")
}

// Java: NormalizeTest.testRewriteExportSpecShorthand1()
func (c *NormalizeTest) TestRewriteExportSpecShorthand1() {
	c.Test("var a; export {a};", "var a; export {a as a};")
}

// Java: NormalizeTest.testRewriteExportSpecShorthand2()
func (c *NormalizeTest) TestRewriteExportSpecShorthand2() {
	c.Test(
		"let a, b, d; export {a, b as c, d};",
		"let a; let b; let d;export {a as a, b as c, d as d};")
}

// Java: NormalizeTest.testSplitExportDeclarationWithVar()
func (c *NormalizeTest) TestSplitExportDeclarationWithVar() {
	c.Test("export var a;", "var a; export {a as a};")
	c.Test("export var a = 4;", "var a = 4; export {a as a};")
	c.Test(
		"export var a, b;",
		`var a;
var b;
export {a as a, b as b};
`)
}

// Java: NormalizeTest.testSplitExportDeclarationWithShorthandProperty()
func (c *NormalizeTest) TestSplitExportDeclarationWithShorthandProperty() {
	c.Test("export var a = {b};", "var a = {b: b}; export {a as a};")
}

// Java: NormalizeTest.testSplitExportDeclarationWithDestructuring()
func (c *NormalizeTest) TestSplitExportDeclarationWithDestructuring() {
	c.IgnoreWarnings(DiagnosticGroups_MISSING_PROPERTIES)
	c.Test("export var {} = {};", "({} = {}); export {};")
	c.Test(
		`let obj = {a: 3, b: 2};
export var {a, b: d, e: f = 2} = obj;
`,
		`let obj = {a: 3, b: 2};
var a; var d; var f;
({a: a, b: d, e: f = 2} = obj);
export {a as a, d as d, f as f};
`)
}

// Java: NormalizeTest.testSplitExportDeclarationWithLet()
func (c *NormalizeTest) TestSplitExportDeclarationWithLet() {
	c.Test("export let a;", "let a; export {a as a};")
}

// Java: NormalizeTest.testSplitExportDeclarationWithConst()
func (c *NormalizeTest) TestSplitExportDeclarationWithConst() {
	c.Test("export const a = 17;", "const a = 17; export {a as a};")
}

// Java: NormalizeTest.testSplitExportDeclarationOfFunction()
func (c *NormalizeTest) TestSplitExportDeclarationOfFunction() {
	c.Test(
		"export function bar() {};",
		`function bar() {}
export {bar as bar};
`)
	// Don't need to split declarations in default exports since they are either unnamed, or the
	// name is declared in the module scope only.
	c.TestSame("export default function() {};")
	c.TestSame("export default function foo() {};")
}

// Java: NormalizeTest.testSplitExportDeclarationOfClass()
func (c *NormalizeTest) TestSplitExportDeclarationOfClass() {
	c.Test(
		"export class Foo {};",
		`class Foo {}
export {Foo as Foo};
`)
	c.TestSame("export default class Bar {}")
	c.TestSame("export default class {}")
}

// Java: NormalizeTest.testLogicalAssignInLabeledLoop()
func (c *NormalizeTest) TestLogicalAssignInLabeledLoop() {
	// TODO(b/571641723): Property logical assignment in labeled loop crashes AstValidator.
	e := AssertThrows[jutil.IRuntimeException](func() { c.TestSame("label: while (o.p ??= 1) {}") })
	// Throwables.getRootCause(e)
	var root jutil.Throwable = e
	for root.GetCause() != nil && root.GetCause() != root {
		root = root.GetCause()
	}
	_, isIllegalState := root.(jutil.IIllegalStateException)
	AssertThat(isIllegalState).IsTrue()
	AssertThat(root.GetMessage()).Contains("Expected 2 children, but was 3")
}
