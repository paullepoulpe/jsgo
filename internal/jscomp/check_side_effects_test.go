// Port of com/google/javascript/jscomp/CheckSideEffectsTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/CheckSideEffectsTest.java)

/*
 * Copyright 2006 The Closure Compiler Authors.
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
)

// CheckSideEffectsTest tests CheckSideEffects.
type CheckSideEffectsTest struct {
	CompilerTestCaseBase

	e *DiagnosticType
}

// Java: CheckSideEffectsTest()
func NewCheckSideEffectsTest() *CheckSideEffectsTest {
	c := &CheckSideEffectsTest{}
	c.InitCompilerTestCaseBase(c)
	c.e = CheckSideEffects_USELESS_CODE_ERROR
	return c
}

func TestCheckSideEffectsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCheckSideEffectsTest() })
}

// Java: CheckSideEffectsTest.setUp()
func (c *CheckSideEffectsTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableParseTypeInfo()
	c.AllowExternsChanges()
}

// Java: CheckSideEffectsTest.getProcessor(Compiler)
func (c *CheckSideEffectsTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewCheckSideEffects(compiler, true, true)
}

// Java: CheckSideEffectsTest.testUselessCode()
func (c *CheckSideEffectsTest) TestUselessCode() {
	c.TestSame("function f(x) { if(x) return; }")
	c.TestWarning("function f(x) { if(x); }", c.e)
	c.TestSame("var f = x=>x")
	c.TestWarning("var f = (x)=>{ if(x); }", c.e)
	c.TestSame("if(x) x = y;")
	c.Test("if(x) x == bar();", "if(x) JSCOMPILER_PRESERVE(x == bar());", c.Warning(c.e))
	c.TestSame("x = 3;")
	c.Test("x == 3;", "JSCOMPILER_PRESERVE(x == 3);", c.Warning(c.e))
	c.TestSame("var x = 'test'")
	c.Test(
		"var x = 'test'\n'Breakstr'",
		"var x = 'test'\nJSCOMPILER_PRESERVE('Breakstr')",
		c.Warning(c.e))
	c.TestSame("")
	c.TestSame("foo();;;;bar();;;;")
	c.TestSame("var a, b; a = 5, b = 6")
	c.Test("var a, b; a = 5, b == 6", "var a, b; a = 5, JSCOMPILER_PRESERVE(b == 6)", c.Warning(c.e))
	c.Test("var a, b; a = (5, 6)", "var a, b; a = (JSCOMPILER_PRESERVE(5), 6)", c.Warning(c.e))
	c.Test("var a, b; a ||= (5, 6)", "var a, b; a ||= (JSCOMPILER_PRESERVE(5), 6)", c.Warning(c.e))
	c.Test(
		"var a, b; a = ('marker', 6)",
		"var a, b; a = (JSCOMPILER_PRESERVE('marker'), 6)",
		c.Warning(c.e))
	c.Test(
		"var a, b; a ||= ('marker', 6)",
		"var a, b; a ||= (JSCOMPILER_PRESERVE('marker'), 6)",
		c.Warning(c.e))
	c.Test(
		"var a, b; a = (bar(), 6, 7)",
		"var a, b; a = (bar(), JSCOMPILER_PRESERVE(6), 7)",
		c.Warning(c.e))
	c.Test(
		"var a, b; a = (bar(), bar(), 7, 8)",
		"var a, b; a = (bar(), bar(), JSCOMPILER_PRESERVE(7), 8)",
		c.Warning(c.e))
	c.Test(
		"var a, b; a ||= (bar(), bar(), 7, 8)",
		"var a, b; a ||= (bar(), bar(), JSCOMPILER_PRESERVE(7), 8)",
		c.Warning(c.e))
	c.TestSame("var a, b; a = (b = 7, 6)")
	c.TestSame("var a, b; a ||= (b ||= 7, 6)")
	c.TestSame(
		`function x(){}
function f(a, b){}
f(1,(x(), 2));
`)
	c.Test(
		`function x(){}
function f(a, b){}
f(1,(2, 3));
`,
		`function x(){}
function f(a, b){}
f(1,(JSCOMPILER_PRESERVE(2), 3));
`,
		c.Warning(c.e))
}

// Java: CheckSideEffectsTest.testUselessCodeTemplateStirngs()
func (c *CheckSideEffectsTest) TestUselessCodeTemplateStirngs() {
	c.Test(
		"var x = `TemplateA`\n'TestB'",
		"var x = `TemplateA`\nJSCOMPILER_PRESERVE('TestB')",
		c.Warning(c.e))
	c.Test("`LoneTemplate`", "JSCOMPILER_PRESERVE(`LoneTemplate`)", c.Warning(c.e))
	c.Test(
		"var name = 'Bad';\n"+
			"`${name}Template`;\n",
		"var name = 'Bad';\n"+
			"JSCOMPILER_PRESERVE(`${name}Template`)\n",
		c.Warning(c.e))
	c.Test(
		"var name = 'BadTail';\n"+
			"`Template${name}`;\n",
		"var name = 'BadTail';\n"+
			"JSCOMPILER_PRESERVE(`Template${name}`)\n",
		c.Warning(c.e))
	c.TestSame(
		"var name = 'Good';\n" +
			"var templateString = `${name}Template`;\n")
	c.TestSame("var templateString = `Template`;")
	c.TestSame("tagged`Template`;")
	c.TestSame("tagged`${name}Template`;")
}

// Java: CheckSideEffectsTest.testUselessCodeDestructuring()
func (c *CheckSideEffectsTest) TestUselessCodeDestructuring() {
	c.TestSame(
		`var obj = {
  itm1: 1,
  itm2: 2
}
var { itm1: de_item1, itm2: de_item2 } = obj;
`)
	c.TestSame(
		`var obj = {
  itm1: 1,
  itm2: 2
}
var { itm1, itm2 } = obj;
`)
	c.TestSame(
		`var arr = ['item1', 'item2', 'item3'];
var [ itm1 = 1, itm2 = 2 ] = arr;
`)
	c.TestSame(
		`var arr = ['item1', 'item2', 'item3'];
var [ itm1 = 1, itm2 = 2 ] = badArr;
`)
	c.TestSame(
		`var arr = ['item1', 'item2', 'item3'];
function f(){}
var [ itm1 = f(), itm2 = 2 ] = badArr;
`)
}

// Java: CheckSideEffectsTest.testUselessCodeParameters()
func (c *CheckSideEffectsTest) TestUselessCodeParameters() {
	c.TestSame("function c(a, b = 1) {}; c(1);")
	c.TestSame("function c(a, b = f()) {}; c(1);")
	c.TestSame("function c(a, {b, c}) {}; c(1);")
	c.TestSame("function c(a, {b, c}) {}; c(1, {b: 2, c: 3});")
}

// Java: CheckSideEffectsTest.testUselessCodeArrowFunctions()
func (c *CheckSideEffectsTest) TestUselessCodeArrowFunctions() {
	c.TestWarning("var f = s => {key:s}", c.e)
	c.TestWarning("var f = s => {key:s + 1}", c.e)
	c.TestWarning("var f = s => {s}", c.e)
	c.TestSame("var f = s => {s ||= 1}")
}

// just here to make sure import.meta doesn't break anything
// Java: CheckSideEffectsTest.testImportMeta()
func (c *CheckSideEffectsTest) TestImportMeta() {
	c.TestSame("var x = import.meta")
}

// Java: CheckSideEffectsTest.testUselessCodeInFor()
func (c *CheckSideEffectsTest) TestUselessCodeInFor() {
	c.TestSame("for(var x = 0; x < 100; x++) { foo(x) }")
	c.TestSame("for(; true; ) { bar() }")
	c.TestSame("for(foo(); true; foo()) { bar() }")
	c.Test(
		"for(i < 5; true; foo()) { bar() }",
		"for(JSCOMPILER_PRESERVE(i < 5); true; foo()) { bar() }",
		c.Warning(c.e))
	c.Test(
		"for(var i = foo(); i++; i < 5) { bar() }",
		"for(var i = foo(); i++; JSCOMPILER_PRESERVE(i < 5)) { bar() }",
		c.Warning(c.e))
	c.Test(
		"for(foo(); true; (1, bar())) { bar() }",
		"for(foo(); true; (JSCOMPILER_PRESERVE(1), bar())) { bar() }",
		c.Warning(c.e))
	c.TestSame("for(foo in bar) { foo() }")
	c.TestSame("for (i = 0; el = el.previousSibling; i++) {}")
	c.TestSame("for (i = 0; el = el.previousSibling; i++);")
}

// Java: CheckSideEffectsTest.testUslessCodeInClassStaticBlock()
func (c *CheckSideEffectsTest) TestUslessCodeInClassStaticBlock() {
	c.TestSame(
		`class C {
  static {
    function f(x) { if(x) return; }
  }
}
`)
	c.TestWarning(
		`class C {
  static {
    function f(x) { if(x); }
  }
}
`,
		c.e)
	c.TestSame(
		`class C {
  static {
    var f = x=>x;
  }
}
`)
	c.TestWarning(
		`class C {
  static {
    var f = (x)=>{ if(x); }
  }
}
`,
		c.e)
	c.TestSame(
		`class C {
  static {
    x = 3;
  }
}
`)
	c.Test(
		`class C {
  static {
    x == 3;
  }
}
`,
		`class C {
  static {
    JSCOMPILER_PRESERVE(x == 3);
  }
}
`,
		c.Warning(c.e))
	c.TestSame(
		`class C {
  static {
     foo();;;;bar();;;;
  }
}
`)
}

// Java: CheckSideEffectsTest.testTypeAnnotations()
func (c *CheckSideEffectsTest) TestTypeAnnotations() {
	c.Test("x;", "JSCOMPILER_PRESERVE(x);", c.Warning(c.e))
	c.Test("a.b.c.d;", "JSCOMPILER_PRESERVE(a.b.c.d);", c.Warning(c.e))
	c.TestSame("/** @type {Number} */ a.b.c.d;")
	c.TestSame("if (true) { /** @type {Number} */ a.b.c.d; }")
	c.Test(
		"function A() { this.foo; }",
		"function A() { JSCOMPILER_PRESERVE(this.foo); }",
		c.Warning(c.e))
	c.TestSame("function A() { /** @type {Number} */ this.foo; }")
}

// Java: CheckSideEffectsTest.testJSDocComments()
func (c *CheckSideEffectsTest) TestJSDocComments() {
	c.TestSame("function A() { /** This is a JsDoc comment */ this.foo; }")
	c.Test(
		"function A() { /* This is a normal comment */ this.foo; }",
		"function A() { /* This is a normal comment */ JSCOMPILER_PRESERVE(this.foo); }",
		c.Warning(c.e))
}

// Java: CheckSideEffectsTest.testIndirectCallExpressions()
func (c *CheckSideEffectsTest) TestIndirectCallExpressions() {
	// indirect call to remove context from eval.
	c.TestSame("(0, eval)('alert');")
	// indirect call to remove this context from module nested functions.
	c.TestSame("(0, modPrefix.foo)('alert');")
	c.TestSame("(0, modPrefix['foo'])('alert');")
	// tagged template literals are a form of call
	c.TestSame("(0, modPrefix.foo)`alert ${x}`;")
	c.TestSame("(0, modPrefix['foo'])`alert ${x}`;")
	// comma not in a call expression.
	c.Test("x = (0, foo) + 1;", "x = (JSCOMPILER_PRESERVE(0), foo) + 1;", c.Warning(c.e))
	// Plain function, no indirection needed to remove this
	c.TestSame("(0, otherFn)(1);")
	// Other expression.
	c.Test("(0, otherFn())(1);", "otherFn()(1);")
}

// Java: CheckSideEffectsTest.testIssue504()
func (c *CheckSideEffectsTest) TestIssue504() {
	// See https://blickly.github.io/closure-compiler-issues/#504, though the `void` idiom in the
	// original issue is still relevant externally (see ESLint's ignoreVoid option for the
	// no-floating-promises rule), so it's no longer an error.
	c.Test(
		"+f();",
		"JSCOMPILER_PRESERVE(+f());",
		c.Warning(c.e).
			WithMessage("Suspicious code. The result of the 'pos' operator is not being used."))
}

// Java: CheckSideEffectsTest.testExternFunctions()
func (c *CheckSideEffectsTest) TestExternFunctions() {
	externs := `/** @return {boolean}
  * @nosideeffects */
function noSideEffectsExtern(){}
/** @return {boolean}
  * @nosideeffects */
var noSideEffectsExtern2 = function(){};
/** @return {boolean} */ function hasSideEffectsExtern(){}
/** @return {boolean} */ var hasSideEffectsExtern2 = function(){}
`
	c.TestSame(c.Externs(externs), c.Srcs("alert(noSideEffectsExtern());"))
	c.Test(
		c.Externs(externs),
		c.Srcs("noSideEffectsExtern();"),
		c.Expected("JSCOMPILER_PRESERVE(noSideEffectsExtern());"),
		c.Warning(c.e).
			WithMessage(
				"Suspicious code. The result of the extern function call 'noSideEffectsExtern' is not being used."))
	c.Test(
		c.Externs(externs),
		c.Srcs("noSideEffectsExtern2();"),
		c.Expected("JSCOMPILER_PRESERVE(noSideEffectsExtern2());"),
		c.Warning(c.e).
			WithMessage(
				"Suspicious code. The result of the extern function call 'noSideEffectsExtern2' is not being used."))
	c.TestSame(c.Externs(externs), c.Srcs("hasSideEffectsExtern()"))
	c.TestSame(c.Externs(externs), c.Srcs("hasSideEffectsExtern2()"))
	// Methods redefined in inner scopes should not trigger a warning
	c.TestSame(
		c.Externs(externs),
		c.Srcs("(function() { function noSideEffectsExtern() {}; noSideEffectsExtern(); })()"))
}

// Java: CheckSideEffectsTest.testExternPropertyFunctions()
func (c *CheckSideEffectsTest) TestExternPropertyFunctions() {
	externs := `/** @const */ var foo = {};
/** @return {boolean}
  * @nosideeffects */
foo.noSideEffectsExtern = function(){}
`
	c.TestSame(c.Externs(externs), c.Srcs("alert(foo.noSideEffectsExtern());"))
	c.Test(
		c.Externs(externs),
		c.Srcs("foo.noSideEffectsExtern();"),
		c.Expected("JSCOMPILER_PRESERVE(foo.noSideEffectsExtern());"),
		c.Warning(c.e).
			WithMessage(
				"Suspicious code. The result of the extern function call 'foo.noSideEffectsExtern' is not being used."))
	// Methods redefined in inner scopes should not trigger a warning
	c.TestSame(
		c.Externs(externs),
		c.Srcs(
			`(function() {
  var foo = {};
  foo.noSideEffectsExtern = function() {};
  noSideEffectsExtern();
 })()
`))
}

// Java: CheckSideEffectsTest.testIgnoresRuntimeLibRequireStatementsOnly()
func (c *CheckSideEffectsTest) TestIgnoresRuntimeLibRequireStatementsOnly() {
	c.TestSame(c.Srcs(SourceFile_FromCode(AbstractCompiler_RUNTIME_LIB_DIR, "'require base'")))
	c.TestWarning(
		c.Srcs(SourceFile_FromCode(AbstractCompiler_RUNTIME_LIB_DIR, "'other string lit'")),
		c.Warning(CheckSideEffects_USELESS_CODE_ERROR))
	c.TestWarning(
		c.Srcs(
			SourceFile_FromCode(
				AbstractCompiler_RUNTIME_LIB_DIR, "function f() { 'require base' }")),
		c.Warning(CheckSideEffects_USELESS_CODE_ERROR))
	c.TestWarning(
		c.Srcs(SourceFile_FromCode("other/file.js", "'require base'")),
		c.Warning(CheckSideEffects_USELESS_CODE_ERROR))
}

// Java: CheckSideEffectsTest.testToStringWithRadixIsUselessCode()
func (c *CheckSideEffectsTest) TestToStringWithRadixIsUselessCode() {
	c.TestWarning("({}).toString(16);", c.e)
}

// Java: CheckSideEffectsTest.testOptionalChainToStringWithRadixIsUselessCode()
func (c *CheckSideEffectsTest) TestOptionalChainToStringWithRadixIsUselessCode() {
	c.TestWarning("({})?.toString(16);", c.e)
}
