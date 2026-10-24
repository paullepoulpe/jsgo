// Port of com/google/javascript/jscomp/InlineObjectLiteralsTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/InlineObjectLiteralsTest.java)

/*
 * Copyright 2011 The Closure Compiler Authors.
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
)

type InlineObjectLiteralsTest struct {
	CompilerTestCaseBase
}

// Java: InlineObjectLiteralsTest()
func NewInlineObjectLiteralsTest() *InlineObjectLiteralsTest {
	c := &InlineObjectLiteralsTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestInlineObjectLiteralsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewInlineObjectLiteralsTest() })
}

// Java: InlineObjectLiteralsTest.setUp()
func (c *InlineObjectLiteralsTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableNormalize()
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
}

// Java: InlineObjectLiteralsTest.getNumRepetitions()
func (c *InlineObjectLiteralsTest) GetNumRepetitions() int {
	// TODO(b/33104006): remove this override.
	return 2
}

// Java: InlineObjectLiteralsTest.getProcessor(Compiler)
func (c *InlineObjectLiteralsTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewInlineObjectLiterals(compiler, compiler.GetUniqueNameIdSupplier())
}

// Test object literal -> variable inlining

// Java: InlineObjectLiteralsTest.testObject0()
func (c *InlineObjectLiteralsTest) TestObject0() {
	// Don't mess with global variables, that is the job of CollapseProperties.
	c.TestSame("var a = {x:1}; f(a.x);")
}

// Java: InlineObjectLiteralsTest.testObject1()
func (c *InlineObjectLiteralsTest) TestObject1() {
	c.TestLocal(
		"var a = {x:x(), y:y()}; f(a.x, a.y);",
		"var JSCompiler_object_inline_x_0=x();\n"+
			"var JSCompiler_object_inline_y_1=y();\n"+
			"f(JSCompiler_object_inline_x_0, JSCompiler_object_inline_y_1);\n")
}

// Java: InlineObjectLiteralsTest.testObject1a()
func (c *InlineObjectLiteralsTest) TestObject1a() {
	c.TestLocal(
		"var a; a = {x:x, y:y}; f(a.x, a.y);",
		"var JSCompiler_object_inline_x_0;\n"+
			"var JSCompiler_object_inline_y_1;\n"+
			"(JSCompiler_object_inline_x_0=x,\n"+
			"JSCompiler_object_inline_y_1=y, true);\n"+
			"f(JSCompiler_object_inline_x_0, JSCompiler_object_inline_y_1);\n")
}

// Java: InlineObjectLiteralsTest.testObject2()
func (c *InlineObjectLiteralsTest) TestObject2() {
	c.TestLocal(
		"var a = {y:y}; a.x = z; f(a.x, a.y);",
		"var JSCompiler_object_inline_y_0 = y;\n"+
			"var JSCompiler_object_inline_x_1;\n"+
			"JSCompiler_object_inline_x_1=z;\n"+
			"f(JSCompiler_object_inline_x_1, JSCompiler_object_inline_y_0);\n")
}

// Java: InlineObjectLiteralsTest.testObject3()
func (c *InlineObjectLiteralsTest) TestObject3() {
	// Inlining the 'y' would cause the 'this' to be different in the
	// target function.
	c.TestSameLocal("var a = {y:y,x:x}; a.y(); f(a.x);")
	c.TestSameLocal("var a; a = {y:y,x:x}; a.y(); f(a.x);")
	c.TestSameLocal("var a = {y:y,x:x}; a.y?.(); f(a.x);")
	c.TestSameLocal("var a = {y:y,x:x}; a?.y(); f(a.x);")
	c.TestSameLocal("var a = {y:y,x:x}; a?.y?.(); f(a.x);")
	c.TestSameLocal("var a = {y:y,x:x}; a.y`hi`; f(a.x);")
}

// Java: InlineObjectLiteralsTest.testObject4()
func (c *InlineObjectLiteralsTest) TestObject4() {
	// Object literal is escaped.
	c.TestSameLocal("var a = {y:y}; a.x = z; f(a.x, a.y); g(a);")
	c.TestSameLocal("var a; a = {y:y}; a.x = z; f(a.x, a.y); g(a);")
}

// Java: InlineObjectLiteralsTest.testObject5()
func (c *InlineObjectLiteralsTest) TestObject5() {
	c.TestLocal(
		"var a = {x:x, y:y}; var b = {a:a}; f(b.a.x, b.a.y);",
		"var a = {x:x, y:y};\n"+
			"var JSCompiler_object_inline_a_0=a;\n"+
			"f(JSCompiler_object_inline_a_0.x, JSCompiler_object_inline_a_0.y);\n")
}

// https://github.com/google/closure-compiler/issues/3658
// Java: InlineObjectLiteralsTest.testNoCrashInliningObjectLiteral_conditionalAssignmentToObject_es2017()
func (c *InlineObjectLiteralsTest) TestNoCrashInliningObjectLiteral_conditionalAssignmentToObject_es2017() {
	CompilationLevel_ADVANCED_OPTIMIZATIONS.SetOptionsForCompilationLevel(c.GetOptions())

	src := "{let a;if(Math.random()){a={x:''};console.log(a.x)}};"
	expected :=
		"{\n" +
			"  var JSCompiler_object_inline_x_0;\n" +
			"  if (Math.random()) {\n" +
			"    JSCompiler_object_inline_x_0 = '', true;\n" +
			"    console.log(JSCompiler_object_inline_x_0);\n" +
			"  }\n" +
			"};\n"

	c.Test(src, expected)
}

// Java: InlineObjectLiteralsTest.testNoInlining_propertyAccessOnVar_unconditionalAssignmentToObject_es2017()
func (c *InlineObjectLiteralsTest) TestNoInlining_propertyAccessOnVar_unconditionalAssignmentToObject_es2017() {
	CompilationLevel_ADVANCED_OPTIMIZATIONS.SetOptionsForCompilationLevel(c.GetOptions())

	// `var` declaration corresponds to global scope: ineligible for inlining object literal
	src := "{var a; a={x:\"\"};console.log(a.x);}"
	c.TestSame(src)
}

// https://github.com/google/closure-compiler/issues/3658
// Java: InlineObjectLiteralsTest.testNoCrashInliningObjectLiteral_unconditionalAssignmentToObject_es2017()
func (c *InlineObjectLiteralsTest) TestNoCrashInliningObjectLiteral_unconditionalAssignmentToObject_es2017() {
	CompilationLevel_ADVANCED_OPTIMIZATIONS.SetOptionsForCompilationLevel(c.GetOptions())

	src := "{let a; a={x:\"\"};console.log(a.x);}"
	expected :=
		"{\n" +
			"  var JSCompiler_object_inline_x_0;\n" +
			"  JSCompiler_object_inline_x_0 = \"\", true;\n" +
			"  console.log(JSCompiler_object_inline_x_0);\n" +
			"}\n"
	c.Test(src, expected)
}

// https://github.com/google/closure-compiler/issues/3658
// Java: InlineObjectLiteralsTest.testNoInliningObjectLiteral_conditionalAssignmentToReturnValue_es2017()
func (c *InlineObjectLiteralsTest) TestNoInliningObjectLiteral_conditionalAssignmentToReturnValue_es2017() {
	CompilationLevel_ADVANCED_OPTIMIZATIONS.SetOptionsForCompilationLevel(c.GetOptions())

	src :=
		"const b = () => ({ x: '' })\n" +
			"function main() {\n" +
			"  let a;\n" +
			"  if (Math.random()) {\n" +
			"// regardless of scope, only direct assignments to object literals get inlined.\n" +
			"    a = b();\n" +
			"    alert(a.x);\n" +
			"  }\n" +
			"}\n" +
			"main();\n"
	c.TestSameLocal(src)
}

// https://github.com/google/closure-compiler/issues/3658
// Java: InlineObjectLiteralsTest.testNoInliningObjectLiteral_unconditionalAssignmentToReturnValue_es2017()
func (c *InlineObjectLiteralsTest) TestNoInliningObjectLiteral_unconditionalAssignmentToReturnValue_es2017() {
	CompilationLevel_ADVANCED_OPTIMIZATIONS.SetOptionsForCompilationLevel(c.GetOptions())

	src :=
		"const b = () => ({ x: '' })\n" +
			"function main() {\n" +
			"  let a;\n" +
			"// regardless of scope, only direct assignments to object literals get inlined.\n" +
			"  a = b();\n" +
			"  alert(a.x);\n" +
			"}\n" +
			"main();\n"
	c.TestSameLocal(src)
}

// Java: InlineObjectLiteralsTest.testObject6()
func (c *InlineObjectLiteralsTest) TestObject6() {
	c.TestLocal(
		"for (var i = 0; i < 5; i++) { var a = {i:i,x:x}; f(a.i, a.x); }",
		"for (var i = 0; i < 5; i++) {\n"+
			"  var JSCompiler_object_inline_i_0=i;\n"+
			"  var JSCompiler_object_inline_x_1=x;\n"+
			"  f(JSCompiler_object_inline_i_0,JSCompiler_object_inline_x_1)\n"+
			"}\n")
	c.TestLocal(
		"if (c) { var a = {i:i,x:x}; f(a.i, a.x); }",
		"if (c) {\n"+
			"  var JSCompiler_object_inline_i_0=i;\n"+
			"  var JSCompiler_object_inline_x_1=x;\n"+
			"  f(JSCompiler_object_inline_i_0,JSCompiler_object_inline_x_1)\n"+
			"}\n")
}

// Java: InlineObjectLiteralsTest.testObject6_let()
func (c *InlineObjectLiteralsTest) TestObject6_let() {
	c.TestLocal(
		"for (let i = 0; i < 5; i++) { let a = {i:i,x:x}; f(a.i, a.x); }",
		"for (let i = 0; i < 5; i++) {\n"+
			"  var JSCompiler_object_inline_i_0=i;\n"+
			"  var JSCompiler_object_inline_x_1=x;\n"+
			"  f(JSCompiler_object_inline_i_0,JSCompiler_object_inline_x_1)\n"+
			"}\n")
	c.TestLocal(
		"if (c) { var a = {i:i,x:x}; f(a.i, a.x); }",
		"if (c) {\n"+
			"  var JSCompiler_object_inline_i_0=i;\n"+
			"  var JSCompiler_object_inline_x_1=x;\n"+
			"  f(JSCompiler_object_inline_i_0,JSCompiler_object_inline_x_1)\n"+
			"}\n")
}

// Java: InlineObjectLiteralsTest.testObjectSwitch()
func (c *InlineObjectLiteralsTest) TestObjectSwitch() {
	c.TestLocal(
		"let x;\n"+
			"switch (x) {\n"+
			"  case 'h':\n"+
			"    let a;\n"+
			"    if (Math.random() > 0.5) {\n"+
			"      a = {i:i,x:x};\n"+
			"      f(a.i, a.x);\n"+
			"      break;\n"+
			"    }\n"+
			"}\n",
		"let x;\n"+
			"var JSCompiler_object_inline_i_0;\n"+
			"var JSCompiler_object_inline_x_1;\n"+
			" switch (x) {\n"+
			"   case 'h':\n"+
			"    if (Math.random() > .5) {\n"+
			"      JSCompiler_object_inline_i_0 = i, JSCompiler_object_inline_x_1 = x, true;\n"+
			"      f(JSCompiler_object_inline_i_0,JSCompiler_object_inline_x_1);\n"+
			"      break;\n"+
			"    }\n"+
			"}\n")
	c.TestLocal(
		"let x;\n"+
			" switch (x) {\n"+
			"   case 'h':\n"+
			"   case 'i':\n"+
			"     let a;\n"+
			"     if (true) {\n"+
			"       a = {i:i,x:x};\n"+
			"       f(a.i, a.x);\n"+
			"     }\n"+
			"     break;\n"+
			"}\n",
		"let x;\n"+
			"var JSCompiler_object_inline_i_0;\n"+
			"var JSCompiler_object_inline_x_1;\n"+
			" switch (x) {\n"+
			" case 'h':\n"+
			" case 'i':\n"+
			"  if (true) {\n"+
			"    JSCompiler_object_inline_i_0 = i, JSCompiler_object_inline_x_1 = x, true;\n"+
			"    f(JSCompiler_object_inline_i_0,JSCompiler_object_inline_x_1);\n"+
			"  }\n"+
			"break;\n"+
			"}\n")
	c.TestLocal(
		"let x;\n"+
			"switch (x) {\n"+
			"  case 'h':\n"+
			"    let a;\n"+
			"    if (Math.random() > 0.5) {\n"+
			"      a = {i:i,x:x};\n"+
			"      f(a.i, a.x);\n"+
			"      break;\n"+
			"    }\n"+
			"}\n",
		"let x;\n"+
			"var JSCompiler_object_inline_i_0;\n"+
			"var JSCompiler_object_inline_x_1;\n"+
			" switch (x) {\n"+
			" case 'h':\n"+
			"  if (Math.random() > .5) {\n"+
			"    JSCompiler_object_inline_i_0 = i, JSCompiler_object_inline_x_1 = x, true;\n"+
			"    f(JSCompiler_object_inline_i_0,JSCompiler_object_inline_x_1);\n"+
			"    break;\n"+
			"  }\n"+
			"}\n")
	c.TestLocal(
		"let x;\n"+
			"let z;\n"+
			"switch (x) {\n"+
			"  case 'h':\n"+
			"    if (Math.random() > 0.5) {\n"+
			"      z = {i:i,x:x};\n"+
			"      f(z.i, z.x);\n"+
			"    }\n"+
			"  case 'i':\n"+
			"    ({x:x} = f(z.i, z.x));\n"+
			"    break;\n"+
			"}\n",
		"var JSCompiler_object_inline_i_0;\n"+
			"var JSCompiler_object_inline_x_1;\n"+
			"let x;\n"+
			" switch (x) {\n"+
			" case 'h':\n"+
			"  if (Math.random() > .5) {\n"+
			"    JSCompiler_object_inline_i_0 = i, JSCompiler_object_inline_x_1 = x, true;\n"+
			"    f(JSCompiler_object_inline_i_0,JSCompiler_object_inline_x_1);\n"+
			"  }\n"+
			"  case 'i':\n"+
			"    ({x:x} = f(JSCompiler_object_inline_i_0, JSCompiler_object_inline_x_1));\n"+
			"    break;\n"+
			"}\n")
}

// Java: InlineObjectLiteralsTest.testObject7()
func (c *InlineObjectLiteralsTest) TestObject7() {
	c.TestLocal(
		"var a = {x:x, y:f()}; g(a.x);",
		"var JSCompiler_object_inline_x_0=x;\n"+
			"var JSCompiler_object_inline_y_1=f();\n"+
			"g(JSCompiler_object_inline_x_0)\n")
}

// Java: InlineObjectLiteralsTest.testObject8()
func (c *InlineObjectLiteralsTest) TestObject8() {
	c.TestSameLocal("var a = {x:x,y:y}; var b = {x:y}; f((c?a:b).x);")

	c.TestLocal(
		"var a; if(c) { a={x:x, y:y}; } else { a={x:y}; } f(a.x);",
		"var JSCompiler_object_inline_x_0;\n"+
			"var JSCompiler_object_inline_y_1;\n"+
			"if(c) JSCompiler_object_inline_x_0=x,\n"+
			"      JSCompiler_object_inline_y_1=y,\n"+
			"      true;\n"+
			"else JSCompiler_object_inline_x_0=y,\n"+
			"     JSCompiler_object_inline_y_1=void 0,\n"+
			"     true;\n"+
			"f(JSCompiler_object_inline_x_0)\n")
	c.TestLocal(
		"var a = {x:x,y:y}; var b = {x:y}; c ? f(a.x) : f(b.x);",
		"var JSCompiler_object_inline_x_0 = x;\n"+
			"var JSCompiler_object_inline_y_1 = y;\n"+
			"var JSCompiler_object_inline_x_2 = y;\n"+
			"c ? f(JSCompiler_object_inline_x_0):f(JSCompiler_object_inline_x_2)\n")
}

// Java: InlineObjectLiteralsTest.testObject9()
func (c *InlineObjectLiteralsTest) TestObject9() {
	// There is a call, so no inlining
	c.TestSameLocal(
		"function f(a,b) {\n" +
			"  var x = {a:a,b:b}; x.a(); return x.b;\n" +
			"}\n")

	c.TestLocal(
		"function f(a,b) {\n"+
			"  var x = {a:a,b:b}; g(x.a); x = {a:a,b:2}; return x.b;\n"+
			"}\n",
		"function f(a,b) {\n"+
			"  var JSCompiler_object_inline_a_0 = a;\n"+
			"  var JSCompiler_object_inline_b_1 = b;\n"+
			"  g(JSCompiler_object_inline_a_0);\n"+
			"  JSCompiler_object_inline_a_0 = a,\n"+
			"  JSCompiler_object_inline_b_1=2,\n"+
			"  true;\n"+
			"  return JSCompiler_object_inline_b_1\n"+
			"}\n")

	c.TestLocal(
		"function f(a,b) {\n"+
			"  var x = {a:a,b:b}; g(x.a); x.b = x.c = 2; return x.b;\n"+
			"}\n",
		"function f(a,b) {\n"+
			"  var JSCompiler_object_inline_a_0=a;\n"+
			"  var JSCompiler_object_inline_b_1=b;\n"+
			"  var JSCompiler_object_inline_c_2;\n"+
			"  g(JSCompiler_object_inline_a_0);\n"+
			"  JSCompiler_object_inline_b_1=JSCompiler_object_inline_c_2=2;\n"+
			"  return JSCompiler_object_inline_b_1\n"+
			"}\n")
}

// Java: InlineObjectLiteralsTest.testObject10()
func (c *InlineObjectLiteralsTest) TestObject10() {
	c.TestLocal(
		"var x; var b = f(); x = {a:a, b:b}; if(x.a) g(x.b);",
		"var JSCompiler_object_inline_a_0;\n"+
			"var JSCompiler_object_inline_b_1;\n"+
			"var b = f();\n"+
			"JSCompiler_object_inline_a_0=a,JSCompiler_object_inline_b_1=b,true;\n"+
			"if(JSCompiler_object_inline_a_0) g(JSCompiler_object_inline_b_1)\n")
	c.TestSameLocal("var x = {}; var b = f(); x = {a:a, b:b}; if(x.a) g(x.b) + x.c")
	c.TestLocal(
		"var x; var b = f(); x = {a:a, b:b}; x.c = c; if(x.a) g(x.b) + x.c",
		"var JSCompiler_object_inline_a_0;\n"+
			"var JSCompiler_object_inline_b_1;\n"+
			"var JSCompiler_object_inline_c_2;\n"+
			"var b = f();\n"+
			"JSCompiler_object_inline_a_0 = a,JSCompiler_object_inline_b_1 = b,\n"+
			"  JSCompiler_object_inline_c_2=void 0,true;\n"+
			"JSCompiler_object_inline_c_2 = c;\n"+
			"if (JSCompiler_object_inline_a_0)\n"+
			"  g(JSCompiler_object_inline_b_1) + JSCompiler_object_inline_c_2;\n")
	c.TestLocal(
		"var x = {a:a}; if (b) x={b:b}; f(x.a||x.b);",
		"var JSCompiler_object_inline_a_0 = a;\n"+
			"var JSCompiler_object_inline_b_1;\n"+
			"if(b) JSCompiler_object_inline_b_1 = b,\n"+
			"      JSCompiler_object_inline_a_0 = void 0,\n"+
			"      true;\n"+
			"f(JSCompiler_object_inline_a_0 || JSCompiler_object_inline_b_1)\n")
	c.TestLocal(
		"var x; var y = 5; x = {a:a, b:b, c:c}; if (b) x={b:b}; f(x.a||x.b);",
		"var JSCompiler_object_inline_a_0;\n"+
			"var JSCompiler_object_inline_b_1;\n"+
			"var JSCompiler_object_inline_c_2;\n"+
			"var y=5;\n"+
			"JSCompiler_object_inline_a_0=a,\n"+
			"JSCompiler_object_inline_b_1=b,\n"+
			"JSCompiler_object_inline_c_2=c,\n"+
			"true;\n"+
			"if (b) JSCompiler_object_inline_b_1=b,\n"+
			"       JSCompiler_object_inline_a_0=void 0,\n"+
			"       JSCompiler_object_inline_c_2=void 0,\n"+
			"       true;\n"+
			"f(JSCompiler_object_inline_a_0||JSCompiler_object_inline_b_1)\n")
}

// Java: InlineObjectLiteralsTest.testObject11()
func (c *InlineObjectLiteralsTest) TestObject11() {
	c.TestSameLocal("var x = {a:b}; (x = {a:a}).c = 5; f(x.a);")
	c.TestSameLocal("var x = {a:a}; f(x[a]); g(x[a]);")
}

// Java: InlineObjectLiteralsTest.testObject12()
func (c *InlineObjectLiteralsTest) TestObject12() {
	c.TestSameLocal("var a; a = {x:1, y:2}; f(a.x, a.y2);")
}

// Java: InlineObjectLiteralsTest.testObject13()
func (c *InlineObjectLiteralsTest) TestObject13() {
	c.TestSameLocal("var x = {a:1, b:2}; x = {a:3, b:x.a};")
}

// Java: InlineObjectLiteralsTest.testObject14()
func (c *InlineObjectLiteralsTest) TestObject14() {
	c.TestSameLocal("var x = {a:1}; if ('a' in x) { f(); }")
	c.TestSameLocal("var x = {a:1}; for (var y in x) { f(y); }")
}

// Java: InlineObjectLiteralsTest.testObject15()
func (c *InlineObjectLiteralsTest) TestObject15() {
	c.TestSameLocal("x = x || {}; f(x.a);")
}

// Java: InlineObjectLiteralsTest.testObject16()
func (c *InlineObjectLiteralsTest) TestObject16() {
	c.TestLocal(
		"function f(e) { bar(); x = {a: foo()}; var x; print(x.a); }",
		"function f(e) {\n"+
			"  var JSCompiler_object_inline_a_0;\n"+
			"  bar();\n"+
			"  JSCompiler_object_inline_a_0 = foo(), true;\n"+
			"  print(JSCompiler_object_inline_a_0);\n"+
			"}\n")
}

// Java: InlineObjectLiteralsTest.testObject17()
func (c *InlineObjectLiteralsTest) TestObject17() {
	// Note: Some day, with careful analysis, these two uses could be
	// disambiguated, and the second assignment could be inlined.
	c.TestSameLocal(
		"var a = {a: function(){}};\n" +
			"a.a();\n" +
			"a = {a1: 100};\n" +
			"print(a.a1);\n")
}

// Java: InlineObjectLiteralsTest.testObject18()
func (c *InlineObjectLiteralsTest) TestObject18() {
	c.TestSameLocal("var a,b; b=a={x:x, y:y}; f(b.x);")
}

// Java: InlineObjectLiteralsTest.testObject19()
func (c *InlineObjectLiteralsTest) TestObject19() {
	c.TestSameLocal("var a,b; if(c) { b=a={x:x, y:y}; } else { b=a={x:y}; } f(b.x);")
}

// Java: InlineObjectLiteralsTest.testObject20()
func (c *InlineObjectLiteralsTest) TestObject20() {
	c.TestSameLocal("var a,b; if(c) { b=a={x:x, y:y}; } else { b=a={x:y}; } f(a.x);")
}

// Java: InlineObjectLiteralsTest.testObject21()
func (c *InlineObjectLiteralsTest) TestObject21() {
	c.TestSameLocal("var a,b; b=a={x:x, y:y};")
	c.TestSameLocal(
		"var a,b; if(c) { b=a={x:x, y:y}; }\n" +
			"else { b=a={x:y}; } f(a.x); f(b.x)\n")
	c.TestSameLocal(
		"var a, b; if(c) { if (a={x:x, y:y}) f(); }\n" +
			"else { b=a={x:y}; } f(a.x);\n")
	c.TestSameLocal("var a,b; b = (a = {x:x, y:x});")
	c.TestSameLocal("var a,b; a = {x:x, y:x}; b = a")
	c.TestSameLocal("var a,b; a = {x:x, y:x}; b = x || a")
	c.TestSameLocal("var a,b; a = {x:x, y:x}; b = y && a")
	c.TestSameLocal("var a,b; a = {x:x, y:x}; b = y ? a : a")
	c.TestSameLocal("var a,b; a = {x:x, y:x}; b = y , a")
	c.TestSameLocal("b = x || (a = {x:1, y:2});")
}

// Java: InlineObjectLiteralsTest.testObject22()
func (c *InlineObjectLiteralsTest) TestObject22() {
	c.TestLocal(
		"while(1) { var a = {y:1}; if (b) a.x = 2; f(a.y, a.x);}",
		"for(;1;){\n"+
			" var JSCompiler_object_inline_y_0=1;\n"+
			" var JSCompiler_object_inline_x_1;\n"+
			" if(b) JSCompiler_object_inline_x_1=2;\n"+
			" f(JSCompiler_object_inline_y_0,JSCompiler_object_inline_x_1)\n"+
			"}\n")

	c.TestSameLocal("var a; while (1) { f(a.x, a.y); a = {x:1, y:1};}")
}

// Java: InlineObjectLiteralsTest.testObject23()
func (c *InlineObjectLiteralsTest) TestObject23() {
	c.TestLocal(
		"function f() {\n"+
			"  var templateData = {\n"+
			"    linkIds: {\n"+
			"      CHROME: 'cl',\n"+
			"      DISMISS: 'd'\n"+
			"    }\n"+
			"  };\n"+
			"  var html = templateData.linkIds.CHROME \n"+
			"       + \":\" + templateData.linkIds.DISMISS;\n"+
			"}\n",
		"function f(){\n"+
			"var JSCompiler_object_inline_CHROME_1='cl';\n"+
			"var JSCompiler_object_inline_DISMISS_2='d';\n"+
			"var html=JSCompiler_object_inline_CHROME_1 +\n"+
			" ':' +JSCompiler_object_inline_DISMISS_2}\n")
}

// Java: InlineObjectLiteralsTest.testObject24()
func (c *InlineObjectLiteralsTest) TestObject24() {
	c.TestLocal(
		"function f() {\n"+
			"  var linkIds = {\n"+
			"      CHROME: 1,\n"+
			"  };\n"+
			"  var g = function () {var o = {a: linkIds};}\n"+
			"}\n",
		"function f(){var linkIds={CHROME:1};\n"+
			"var g=function(){var JSCompiler_object_inline_a_0=linkIds}}\n")
}

// Java: InlineObjectLiteralsTest.testObject25()
func (c *InlineObjectLiteralsTest) TestObject25() {
	c.TestLocal(
		"var a = {x:f(), y:g()}; a = {y:g(), x:f()}; f(a.x, a.y);",
		"var JSCompiler_object_inline_x_0=f();\n"+
			"var JSCompiler_object_inline_y_1=g();\n"+
			"JSCompiler_object_inline_y_1=g(),\n"+
			"  JSCompiler_object_inline_x_0=f(),\n"+
			"  true;\n"+
			"f(JSCompiler_object_inline_x_0,JSCompiler_object_inline_y_1)\n")
}

// Java: InlineObjectLiteralsTest.testObject26()
func (c *InlineObjectLiteralsTest) TestObject26() {
	c.TestLocal(
		"var a = {}; a.b = function() {}; new a.b.c",
		"var JSCompiler_object_inline_b_0;\n"+
			"JSCompiler_object_inline_b_0=function(){};\n"+
			"new JSCompiler_object_inline_b_0.c\n")
}

// Java: InlineObjectLiteralsTest.testInlineObjectWithLet()
func (c *InlineObjectLiteralsTest) TestInlineObjectWithLet() {
	c.TestLocal(
		"let a = {x:x(), y:y()}; f(a.x, a.y);",
		"var JSCompiler_object_inline_x_0=x();\n"+
			"var JSCompiler_object_inline_y_1=y();\n"+
			"f(JSCompiler_object_inline_x_0, JSCompiler_object_inline_y_1);\n")
}

// Java: InlineObjectLiteralsTest.testInlineObjectWithConst()
func (c *InlineObjectLiteralsTest) TestInlineObjectWithConst() {
	c.TestLocal(
		"const a = {x:x(), y:y()}; f(a.x, a.y);",
		"var JSCompiler_object_inline_x_0=x();\n"+
			"var JSCompiler_object_inline_y_1=y();\n"+
			"f(JSCompiler_object_inline_x_0, JSCompiler_object_inline_y_1);\n")
}

// Java: InlineObjectLiteralsTest.testDontInlineLetInForLoopInit()
func (c *InlineObjectLiteralsTest) TestDontInlineLetInForLoopInit() {
	// Handling this case should be possible, but we don't currently have that logic in place.
	c.TestSameLocal("var i; for(let a = {x:x(), y:y()}; i < 0; i++) { f(a.x, a.y); }")
}

// Java: InlineObjectLiteralsTest.testDontInlineConstInForLoopInit()
func (c *InlineObjectLiteralsTest) TestDontInlineConstInForLoopInit() {
	c.TestSameLocal("var i; for(const a = {x:x(), y:y()}; i < 0; i++) { f(a.x, a.y); }")
}

// Java: InlineObjectLiteralsTest.testDoInlineObjectWithVarInForLoop()
func (c *InlineObjectLiteralsTest) TestDoInlineObjectWithVarInForLoop() {
	// The Normalize pass actually moves "var a = {x: x(), y: y()}" out of the for loop initializer
	// before this pass runs, which is why this works.
	c.TestLocal(
		"var i; for(var a = {x:x(), y:y()}; i < 0; i++) { f(a.x, a.y); }",
		"var i;\n"+
			"var JSCompiler_object_inline_x_0=x();\n"+
			"var JSCompiler_object_inline_y_1=y();\n"+
			"for(; i < 0; i++) {\n"+
			"f(JSCompiler_object_inline_x_0, JSCompiler_object_inline_y_1);\n"+
			"}\n")
}

// Java: InlineObjectLiteralsTest.testBug545a()
func (c *InlineObjectLiteralsTest) TestBug545a() {
	c.TestLocal("var a = {}", "")
}

// Java: InlineObjectLiteralsTest.testBug545b()
func (c *InlineObjectLiteralsTest) TestBug545b() {
	c.TestLocal("var a; a = {}", "true")
}

// Java: InlineObjectLiteralsTest.testIssue724()
func (c *InlineObjectLiteralsTest) TestIssue724() {
	c.TestSameLocal(
		"var getType; getType = {};\n" +
			"return functionToCheck &&\n" +
			"   getType.toString.apply(functionToCheck) ===\n" +
			"   '[object Function]';\n")
}

// Java: InlineObjectLiteralsTest.testNoInlineDeletedProperties()
func (c *InlineObjectLiteralsTest) TestNoInlineDeletedProperties() {
	c.TestSameLocal(
		"var foo = {bar:1};\n" +
			"delete foo.bar;\n" +
			"return foo.bar;\n")
}

// Java: InlineObjectLiteralsTest.testProto()
func (c *InlineObjectLiteralsTest) TestProto() {
	c.TestSameLocal(
		"var protoObject = {\n" +
			"  f: function() {\n" +
			"    return 1;\n" +
			"  }\n" +
			"};\n" +
			"var object = {\n" +
			"  __proto__: protoObject,\n" +
			"};\n" +
			"g(object.f);\n")

	c.TestSame(
		c.Externs(
			"var protoObject = {\n"+
				"  f: function() {\n"+
				"    return 1;\n"+
				"  }\n"+
				"};\n"+
				"var object = {\n"+
				"  __proto__: protoObject,\n"+
				"  g: false\n"+
				"};\n"+
				"g(object.g);\n"),
		c.Srcs(
			"var protoObject = {\n"+
				"  f: function() {\n"+
				"    return 1;\n"+
				"  }\n"+
				"};\n"+
				"var JSCompiler_object_inline___proto___0=protoObject;\n"+
				"var JSCompiler_object_inline_g_1=false;\n"+
				"g(JSCompiler_object_inline_g_1)\n"))
}

// Java: InlineObjectLiteralsTest.testSuper()
func (c *InlineObjectLiteralsTest) TestSuper() {
	c.TestSameLocal(
		"var superObject = {\n" +
			"  f() {\n" +
			"    return 1;\n" +
			"  }\n" +
			"};\n" +
			"var object = {\n" +
			"  __proto__: superObject,\n" +
			"  f() {\n" +
			"    return super.f();\n" +
			"  }\n" +
			"}\n" +
			"g(object.f());\n")
}

// Java: InlineObjectLiteralsTest.testShorthandFunctions()
func (c *InlineObjectLiteralsTest) TestShorthandFunctions() {
	c.TestSameLocal(
		"var object = {\n" +
			"  items: [],\n" +
			"  add(item) {\n" +
			"    this.items.push(item);\n" +
			"  },\n" +
			"};\n" +
			"object.add(1);\n")

	c.TestSameLocal(
		"var object = {\n" +
			"  one() {\n" +
			"    return 1\n" +
			"  },\n" +
			"};\n" +
			"object.one();\n")
}

// Java: InlineObjectLiteralsTest.testShorthandAssignments()
func (c *InlineObjectLiteralsTest) TestShorthandAssignments() {
	c.TestLocal(
		"var object = {\n"+
			"  x,\n"+
			"  y\n"+
			"};\n"+
			"f(object.x, object.y);\n",
		"var JSCompiler_object_inline_x_0=x;\n"+
			"var JSCompiler_object_inline_y_1=y;\n"+
			"f(JSCompiler_object_inline_x_0,JSCompiler_object_inline_y_1)\n")

	c.TestLocal(
		"var object = {\n"+
			"  x,\n"+
			"};\n"+
			"object.y = y\n"+
			"f(object.x, object.y);\n",
		"var JSCompiler_object_inline_x_0=x;\n"+
			"var JSCompiler_object_inline_y_1;\n"+
			"var JSCompiler_object_inline_y_1=y;\n"+
			"f(JSCompiler_object_inline_x_0,JSCompiler_object_inline_y_1)\n")
}

// Java: InlineObjectLiteralsTest.testComputedPropertyName()
func (c *InlineObjectLiteralsTest) TestComputedPropertyName() {
	c.TestSameLocal(
		"function addBar(name) {\n" +
			"  return name + 'Bar'\n" +
			"}\n" +
			"var object = {\n" +
			"  [addBar(\"foo\")]: 1\n" +
			"};\n")

	c.TestSameLocal(
		"var sym = Symbol('key');\n" +
			"var object = {\n" +
			"  [sym]: 1,\n" +
			"  x: true\n" +
			"}\n" +
			"use(object[sym]);\n")
}

// Java: InlineObjectLiteralsTest.testQuotedKeyThatIsNotRead()
func (c *InlineObjectLiteralsTest) TestQuotedKeyThatIsNotRead() {
	c.TestLocal("var obj = {'a.b.c': 'd'};", "var JSCompiler_object_inline_string_key_0 = 'd';")
	c.TestLocal(
		"var obj = {'@': 5, '!': 4, 'foo': 3};",
		"var JSCompiler_object_inline_string_key_0 = 5;\n"+
			"var JSCompiler_object_inline_string_key_1 = 4;\n"+
			"var JSCompiler_object_inline_foo_2 = 3;\n")

	c.TestSameLocal("var obj = {}; obj['@'] = 3;")
}

// Java: InlineObjectLiteralsTest.testQuotedKeyThatIsRead()
func (c *InlineObjectLiteralsTest) TestQuotedKeyThatIsRead() {
	c.TestSameLocal("var obj = {'a.b.c': 'd'}; use(obj['a.b.c']);")

	c.TestSameLocal("var obj = {}; obj['a'] = 3; use(obj['a']);")
}

// Java: InlineObjectLiteralsTest.testObjectSpread_readingFromSpreadAssignment()
func (c *InlineObjectLiteralsTest) TestObjectSpread_readingFromSpreadAssignment() {
	c.TestSameLocal(
		"var obj = {...foo};\n" +
			"use(obj.bar);\n")
}

// Java: InlineObjectLiteralsTest.testObjectSpread_readingFromOverwrittenProp()
func (c *InlineObjectLiteralsTest) TestObjectSpread_readingFromOverwrittenProp() {
	c.TestSameLocal(
		"var foo = {prop: 7};\n" +
			"var obj = {prop: 6, ...foo};\n" +
			"use(obj.prop);\n")
}

const InlineObjectLiteralsTest_LOCAL_PREFIX = "function local(){"
const InlineObjectLiteralsTest_LOCAL_POSTFIX = "}"

// Java: InlineObjectLiteralsTest.testLocal(String, String)
func (c *InlineObjectLiteralsTest) TestLocal(code string, result string) {
	c.Test(InlineObjectLiteralsTest_LOCAL_PREFIX+code+InlineObjectLiteralsTest_LOCAL_POSTFIX, InlineObjectLiteralsTest_LOCAL_PREFIX+result+InlineObjectLiteralsTest_LOCAL_POSTFIX)
}

// Java: InlineObjectLiteralsTest.testSameLocal(String)
func (c *InlineObjectLiteralsTest) TestSameLocal(code string) {
	c.TestSame(InlineObjectLiteralsTest_LOCAL_PREFIX + code + InlineObjectLiteralsTest_LOCAL_POSTFIX)
}

// Java: InlineObjectLiteralsTest.testObject1_optChainContinues()
func (c *InlineObjectLiteralsTest) TestObject1_optChainContinues() {
	// TODO(b/571220152): inlining an opt-chain segment start orphans the continuation
	e :=
		AssertThrows[jutil.IRuntimeException](
			func() {
				c.TestSame("function local(){var a = {x:x(), y:y()}; f(a?.x.toString(), a?.y);}")
			})
	// Throwables.getRootCause(e)
	var root jutil.Throwable = e
	for root.GetCause() != nil {
		root = root.GetCause()
	}
	IsInstanceOf[*jutil.IllegalStateException](AssertThat(root))
	AssertThat(root.GetMessage()).
		Contains("Start of optional chain node OPTCHAIN_GETPROP is not marked as the start")
}

// Java: InlineObjectLiteralsTest.testClassField_perIterationLetCapture()
func (c *InlineObjectLiteralsTest) TestClassField_perIterationLetCapture() {
	// TODO(b/572003863): let/const object literal in a loop body is split into function-scoped vars
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2022)
	c.TestLocal(
		`for (let i = 0; i < 2; i++) {
  let o = {a: i};
  arr.push(class { x = o.a; });
}
`,
		`for (let i = 0; i < 2; i++) {
  var JSCompiler_object_inline_a_0 = i;
  arr.push(class {
    x = JSCompiler_object_inline_a_0;
  });
}
`)
}

// Java: InlineObjectLiteralsTest.testClassField_perIterationLetCapture_arrowControl()
func (c *InlineObjectLiteralsTest) TestClassField_perIterationLetCapture_arrowControl() {
	// TODO(b/572003863): let/const object literal in a loop body is split into function-scoped vars
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2022)
	c.TestLocal(
		`for (let i = 0; i < 2; i++) {
  let o = {a: i};
  arr.push(() => o.a);
}
`,
		`for (let i = 0; i < 2; i++) {
  var JSCompiler_object_inline_a_0 = i;
  arr.push(() => {
    return JSCompiler_object_inline_a_0;
  });
}
`)
}
