// Port of com/google/javascript/jscomp/CrossChunkCodeMotionTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/CrossChunkCodeMotionTest.java)

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
)

type CrossChunkCodeMotionTest struct {
	CompilerTestCaseBase
	parentChunkCanSeeSymbolsDeclaredInChildren bool
}

// Java: CrossChunkCodeMotionTest()
func NewCrossChunkCodeMotionTest() *CrossChunkCodeMotionTest {
	c := &CrossChunkCodeMotionTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, CrossChunkCodeMotionTest_EXTERNS)
	c.parentChunkCanSeeSymbolsDeclaredInChildren = false
	return c
}

func TestCrossChunkCodeMotionTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCrossChunkCodeMotionTest() })
}

const CrossChunkCodeMotionTest_EXTERNS = "alert"

// Java: CrossChunkCodeMotionTest.setUp()
func (c *CrossChunkCodeMotionTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()

	c.parentChunkCanSeeSymbolsDeclaredInChildren = false
}

// Java: CrossChunkCodeMotionTest.getProcessor(Compiler)
func (c *CrossChunkCodeMotionTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewCrossChunkCodeMotion(
		compiler, compiler.GetChunkGraph(), c.parentChunkCanSeeSymbolsDeclaredInChildren)
}

// Java: CrossChunkCodeMotionTest.testFunctionMovement1()
func (c *CrossChunkCodeMotionTest) TestFunctionMovement1() {
	// This tests lots of things:
	// 1) f1 is declared in m1, and used in m2. Move it to m2
	// 2) f2 is declared in m1, and used in m3 twice. Move it to m3
	// 3) f3 is declared in m1, and used in m2+m3. It stays put
	// 4) g declared in m1 and never used. It stays put
	// 5) h declared in m2 and never used. It stays put
	// 6) f4 declared in m1 and used in m2 as var. It moves to m2

	chunks :=
		JSChunkGraphBuilder_ForStar().
			// m1
			AddChunk(
				"function f1(a) { alert(a); }\n" +
					"function f2(a) { alert(a); }\n" +
					"function f3(a) { alert(a); }\n" +
					"function f4() { alert(1); }\n" +
					"function g() { alert('ciao'); }\n").
			// m2
			AddChunk("f1('hi'); f3('bye'); var a = f4; function h(a) { alert('h:' + a); }").
			// m3
			AddChunk("f2('hi'); f2('hi'); f3('bye');").
			Build()

	c.Test(
		c.Srcs(chunks),
		c.Expected(
			// m1
			"function f3(a) { alert(a); } function g() { alert('ciao'); }",
			// m2
			"function f1(a) { alert(a); }\n"+
				"function f4() { alert(1); }\n"+
				"f1('hi'); f3('bye'); var a = f4;\n"+
				"function h(a) { alert('h:' + a); }\n",
			// m3
			"function f2(a) { alert(a); } f2('hi'); f2('hi'); f3('bye');"))
}

// Java: CrossChunkCodeMotionTest.testFunctionMovement2()
func (c *CrossChunkCodeMotionTest) TestFunctionMovement2() {
	// having f declared as a local variable should block the migration to m2
	chunks :=
		JSChunkGraphBuilder_ForStar().
			// m1
			AddChunk("function f(a) { alert(a); } function g() {var f = 1; f++}").
			// m2
			AddChunk("f(1);").
			Build()

	c.Test(
		c.Srcs(chunks),
		c.Expected(
			// m1
			"function g() {var f = 1; f++}",
			// m2
			"function f(a) { alert(a); } f(1);"))
}

// Java: CrossChunkCodeMotionTest.testFunctionMovement3()
func (c *CrossChunkCodeMotionTest) TestFunctionMovement3() {
	// having f declared as a arg should block the migration to m2
	chunks :=
		JSChunkGraphBuilder_ForStar().
			// m1
			AddChunk("function f(a) { alert(a); } function g(f) {f++}").
			// m2
			AddChunk("f(1);").
			Build()

	c.Test(
		c.Srcs(chunks),
		c.Expected(
			// m1
			"function g(f) {f++}",
			// m2
			"function f(a) { alert(a); } f(1);"))
}

// Java: CrossChunkCodeMotionTest.testFunctionMovement4()
func (c *CrossChunkCodeMotionTest) TestFunctionMovement4() {
	// Try out moving a function which returns a closure
	chunks :=
		JSChunkGraphBuilder_ForStar().
			// m1
			AddChunk("function f(){return function(a){}}").
			// m2
			AddChunk("var a = f();").
			Build()

	c.Test(
		c.Srcs(chunks),
		c.Expected(
			// m1
			"",
			// m2
			"function f(){return function(a){}} var a = f();"))
}

// Java: CrossChunkCodeMotionTest.testFunctionMovement5()
func (c *CrossChunkCodeMotionTest) TestFunctionMovement5() {
	// Try moving a recursive function [using factorials for kicks]
	chunks :=
		JSChunkGraphBuilder_ForStar().
			// m1
			AddChunk("function f(n){return (n<1)?1:f(n-1)}").
			// m2
			AddChunk("var a = f(4);").
			Build()

	c.Test(
		c.Srcs(chunks),
		c.Expected(
			// m1
			"",
			// m2
			"function f(n){return (n<1)?1:f(n-1)} var a = f(4);"))
}

// Java: CrossChunkCodeMotionTest.testFunctionMovement5b()
func (c *CrossChunkCodeMotionTest) TestFunctionMovement5b() {
	// Try moving a recursive function declared differently.
	chunks :=
		JSChunkGraphBuilder_ForStar().
			// m1
			AddChunk("var f = function(n){return (n<1)?1:f(n-1)};").
			// m2
			AddChunk("var a = f(4);").
			Build()

	c.Test(
		c.Srcs(chunks),
		c.Expected(
			// m1
			"",
			// m2
			"var f = function(n){return (n<1)?1:f(n-1)}; var a = f(4);"))
}

// Java: CrossChunkCodeMotionTest.testFunctionMovement5c()
func (c *CrossChunkCodeMotionTest) TestFunctionMovement5c() {
	// Try moving a recursive function declared differently, in a nested block scope.
	chunks :=
		JSChunkGraphBuilder_ForStar().
			// m1
			AddChunk("var f = function(n){if(true){if(true){return (n<1)?1:f(n-1)}}};").
			// m2
			AddChunk("var a = f(4);").
			Build()

	c.Test(
		c.Srcs(chunks),
		c.Expected(
			// m1
			"",
			// m2
			"var f = function(n){if(true){if(true){return (n<1)?1:f(n-1)}}};\n"+
				"var a = f(4);\n"))
}

// Java: CrossChunkCodeMotionTest.testFunctionMovement6()
func (c *CrossChunkCodeMotionTest) TestFunctionMovement6() {
	// Try out moving to the common ancestor
	chunks :=
		JSChunkGraphBuilder_ForChain().
			// m1
			AddChunk("function f(){return 1}").
			// m2
			AddChunk("var a = f();").
			// m3
			AddChunk("var b = f();").
			Build()

	c.Test(
		c.Srcs(chunks),
		c.Expected(
			// m1
			"",
			// m2
			"function f(){return 1} var a = f();",
			// m3
			"var b = f();"))
}

// Java: CrossChunkCodeMotionTest.testFunctionMovement7()
func (c *CrossChunkCodeMotionTest) TestFunctionMovement7() {
	// Try out moving to the common ancestor with deeper ancestry chain
	chunks :=
		JSChunkGraphBuilder_ForUnordered().
			// m1
			AddChunk("function f(){return 1}").
			// m2
			AddChunk("").
			// m3
			AddChunk("var a = f();").
			// m4
			AddChunk("var b = f();").
			// m5
			AddChunk("var c = f();").
			Build()

	chunks[1].AddDependency(chunks[0])
	chunks[2].AddDependency(chunks[1])
	chunks[3].AddDependency(chunks[1])
	chunks[4].AddDependency(chunks[1])

	c.Test(
		c.Srcs(chunks),
		c.Expected(
			// m1
			"",
			// m2
			"function f(){return 1}",
			// m3
			"var a = f();",
			// m4
			"var b = f();",
			// m5
			"var c = f();"))
}

// Java: CrossChunkCodeMotionTest.testFunctionMovement8()
func (c *CrossChunkCodeMotionTest) TestFunctionMovement8() {
	// Check what happens with named functions
	chunks :=
		JSChunkGraphBuilder_ForChain().
			// m1
			AddChunk("var v = function f(){return 1}").
			// m2
			AddChunk("v();").
			Build()

	c.Test(
		c.Srcs(chunks),
		c.Expected(
			// m1
			"",
			// m2
			"var v = function f(){return 1}; v();"))
}

// Java: CrossChunkCodeMotionTest.testFunctionNonMovement1()
func (c *CrossChunkCodeMotionTest) TestFunctionNonMovement1() {
	// This tests lots of things:
	// 1) we can't move it if it is a class with non-const attributes accessed
	// 2) if it's in an if statement, we can't move it
	// 3) if it's in an while statement, we can't move it [with some extra
	// block elements]
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk(
					"function f(){};f.prototype.bar=new f;\n" +
						"if(a)function f2(){}\n" +
						"{{while(a)function f3(){}}}\n").
				AddChunk("var a = new f();f2();f3();").
				Build()))
}

// Java: CrossChunkCodeMotionTest.testFunctionNonMovement2()
func (c *CrossChunkCodeMotionTest) TestFunctionNonMovement2() {
	// A generic case where 2 chunks depend on the first one. But it's the
	// common ancestor, so we can't move.
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk("function f(){return 1}").
				AddChunk("var a = f();").
				AddChunk("var b = f();").
				Build()))
}

// Java: CrossChunkCodeMotionTest.testEs6ClassMovement1()
func (c *CrossChunkCodeMotionTest) TestEs6ClassMovement1() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk("class f { bar() {} }").
				AddChunk("var a = new f();").
				Build()),
		c.Expected("", "class f { bar() {} } var a = new f();"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk("var f = class { bar() {} };").
				AddChunk("var a = new f();").
				Build()),
		c.Expected("", "var f = class { bar() {} }; var a = new f();"))
}

// Java: CrossChunkCodeMotionTest.testClassMovement1()
func (c *CrossChunkCodeMotionTest) TestClassMovement1() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk("function f(){} f.prototype.bar=function (){};").
				AddChunk("var a = new f();").
				Build()),
		c.Expected("", "function f(){} f.prototype.bar=function (){}; var a = new f();"))
}

// Java: CrossChunkCodeMotionTest.testEs6ClassMovement_instanceof()
func (c *CrossChunkCodeMotionTest) TestEs6ClassMovement_instanceof() {
	c.parentChunkCanSeeSymbolsDeclaredInChildren = true
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk("class f { bar(){} } 1 instanceof f;").
				AddChunk("var a = new f();").
				Build()),
		c.Expected(
			"'function' == typeof f && 1 instanceof f;", "class f { bar(){} } var a = new f();"))
}

// Java: CrossChunkCodeMotionTest.testClassMovement_instanceof()
func (c *CrossChunkCodeMotionTest) TestClassMovement_instanceof() {
	c.parentChunkCanSeeSymbolsDeclaredInChildren = true
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk("function f(){} f.prototype.bar=function (){}; 1 instanceof f;").
				AddChunk("var a = new f();").
				Build()),
		c.Expected(
			"'function' == typeof f && 1 instanceof f;",
			"function f(){} f.prototype.bar=function (){}; var a = new f();"))
}

// Java: CrossChunkCodeMotionTest.testEs6ClassMovement_instanceofTurnedOff()
func (c *CrossChunkCodeMotionTest) TestEs6ClassMovement_instanceofTurnedOff() {
	c.parentChunkCanSeeSymbolsDeclaredInChildren = false
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk("class f { bar(){} } 1 instanceof f;").
				AddChunk("var a = new f();").
				Build()))
}

// Java: CrossChunkCodeMotionTest.testClassMovement_instanceofTurnedOff()
func (c *CrossChunkCodeMotionTest) TestClassMovement_instanceofTurnedOff() {
	c.parentChunkCanSeeSymbolsDeclaredInChildren = false
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk("function f(){} f.prototype.bar=function (){}; 1 instanceof f;").
				AddChunk("var a = new f();").
				Build()))
}

// Java: CrossChunkCodeMotionTest.testEs6ClassMovement_instanceof2()
func (c *CrossChunkCodeMotionTest) TestEs6ClassMovement_instanceof2() {
	c.parentChunkCanSeeSymbolsDeclaredInChildren = true
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk("class f { bar(){} } (true && 1 instanceof f);").
				AddChunk("var a = new f();").
				Build()),
		c.Expected(
			"(true && ('function' == typeof f && 1 instanceof f));",
			"class f { bar(){} } var a = new f();"))
}

// Java: CrossChunkCodeMotionTest.testClassMovement_instanceof2()
func (c *CrossChunkCodeMotionTest) TestClassMovement_instanceof2() {
	c.parentChunkCanSeeSymbolsDeclaredInChildren = true
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk("function f(){} f.prototype.bar=function (){}; (true && 1 instanceof f);").
				AddChunk("var a = new f();").
				Build()),
		c.Expected(
			"(true && ('function' == typeof f && 1 instanceof f));",
			"function f(){} f.prototype.bar=function (){}; var a = new f();"))
}

// Java: CrossChunkCodeMotionTest.testClassMovement_alreadyGuardedInstanceof()
func (c *CrossChunkCodeMotionTest) TestClassMovement_alreadyGuardedInstanceof() {
	c.parentChunkCanSeeSymbolsDeclaredInChildren = true
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk(
					"function f(){} f.prototype.bar=function (){};\n"+
						"(true && ('undefined' != typeof f && 1 instanceof f));\n").
				AddChunk("var a = new f();").
				Build()),
		c.Expected(
			"(true && ('undefined' != typeof f && 1 instanceof f));",
			"function f(){} f.prototype.bar=function (){}; var a = new f();"))
}

// Java: CrossChunkCodeMotionTest.testClassMovement_alreadyGuardedInstanceof_functionGuard()
func (c *CrossChunkCodeMotionTest) TestClassMovement_alreadyGuardedInstanceof_functionGuard() {
	c.parentChunkCanSeeSymbolsDeclaredInChildren = true
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk(
					"function f(){} f.prototype.bar=function (){};\n"+
						"(true && ('function' == typeof f && 1 instanceof f));\n").
				AddChunk("var a = new f();").
				Build()),
		c.Expected(
			"(true && ('function' == typeof f && 1 instanceof f));",
			"function f(){} f.prototype.bar=function (){}; var a = new f();"))
}

// Java: CrossChunkCodeMotionTest.testClassMovement_instanceof3()
func (c *CrossChunkCodeMotionTest) TestClassMovement_instanceof3() {
	c.parentChunkCanSeeSymbolsDeclaredInChildren = true
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk("function f(){} f.prototype.bar=function (){}; f instanceof 1").
				AddChunk("var a = new f();").
				Build()))
}

// Java: CrossChunkCodeMotionTest.testClassMovement_instanceof_sideEffectsLhs()
func (c *CrossChunkCodeMotionTest) TestClassMovement_instanceof_sideEffectsLhs() {
	c.parentChunkCanSeeSymbolsDeclaredInChildren = true
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk(
					"class SafeHtml {} var out = raw; (out = escapeHtml(raw)) instanceof" +
						" SafeHtml;").
				AddChunk("new SafeHtml();").
				Build()))
}

// Java: CrossChunkCodeMotionTest.testClassMovement_instanceof_noRewriteRequired()
func (c *CrossChunkCodeMotionTest) TestClassMovement_instanceof_noRewriteRequired() {
	c.parentChunkCanSeeSymbolsDeclaredInChildren = true
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk("function f(){} f.prototype.bar=function (){}; 1 instanceof f; new f;").
				AddChunk("var a = new f();").
				Build()))
}

// Java: CrossChunkCodeMotionTest.testClassMovement_instanceof_noRewriteRequired2()
func (c *CrossChunkCodeMotionTest) TestClassMovement_instanceof_noRewriteRequired2() {
	c.parentChunkCanSeeSymbolsDeclaredInChildren = true
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("function f(){} f.prototype.bar=function (){}; new f;").
				AddChunk("1 instanceof f;").
				AddChunk("var a = new f();").
				Build()))
}

// Java: CrossChunkCodeMotionTest.testEs6ClassMovement2()
func (c *CrossChunkCodeMotionTest) TestEs6ClassMovement2() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				// m1
				AddChunk("class f {} f.prototype.bar=3; f.prototype.baz=5;").
				// m2
				AddChunk("f.prototype.baq = 7;").
				// m3
				AddChunk("f.prototype.baz = 9;").
				// m4
				AddChunk("var a = new f();").
				Build()),
		c.Expected(
			// m1
			"",
			// m2
			"",
			// m3
			"",
			// m4
			"class f {}\n"+
				"f.prototype.bar = 3;\n"+
				"f.prototype.baz = 5;\n"+
				"f.prototype.baq = 7;\n"+
				"f.prototype.baz = 9;\n"+
				"var a = new f();\n"))
}

// Java: CrossChunkCodeMotionTest.testClassMovement2()
func (c *CrossChunkCodeMotionTest) TestClassMovement2() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				// m1
				AddChunk("function f(){} f.prototype.bar=3; f.prototype.baz=5;").
				// m2
				AddChunk("f.prototype.baq = 7;").
				// m3
				AddChunk("f.prototype.baz = 9;").
				// m4
				AddChunk("var a = new f();").
				Build()),
		c.Expected(
			// m1
			"",
			// m2
			"",
			// m3
			"",
			// m4
			"function f(){}\n"+
				"f.prototype.bar=3;\n"+
				"f.prototype.baz=5;\n"+
				"f.prototype.baq = 7;\n"+
				"f.prototype.baz = 9;\n"+
				"var a = new f();\n"))
}

// Java: CrossChunkCodeMotionTest.testClassMovement3()
func (c *CrossChunkCodeMotionTest) TestClassMovement3() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				// m1
				AddChunk("var f = function() {}; f.prototype.bar=3; f.prototype.baz=5;").
				// m2
				AddChunk("f = 7;").
				// m3
				AddChunk("f = 9;").
				// m4
				AddChunk("f = 11;").
				Build()),
		c.Expected(
			// m1
			"",
			// m2
			"",
			// m3
			"",
			// m4
			"var f = function() {};\n"+
				"f.prototype.bar=3;\n"+
				"f.prototype.baz=5;\n"+
				"f = 7;\n"+
				"f = 9;\n"+
				"f = 11;\n"))
}

// Java: CrossChunkCodeMotionTest.testClassMovement4()
func (c *CrossChunkCodeMotionTest) TestClassMovement4() {
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk("function f(){} f.prototype.bar=3; f.prototype.baz=5;").
				AddChunk("f.prototype.baq = 7;").
				AddChunk("var a = new f();").
				Build()))
}

// Java: CrossChunkCodeMotionTest.testClassMovement5()
func (c *CrossChunkCodeMotionTest) TestClassMovement5() {
	chunks :=
		JSChunkGraphBuilder_ForUnordered().
			// m1
			AddChunk("function f(){} f.prototype.bar=3; f.prototype.baz=5;").
			// m2
			AddChunk("").
			// m3
			AddChunk("f.prototype.baq = 7;").
			// m4
			AddChunk("var a = new f();").
			Build()

	chunks[1].AddDependency(chunks[0])
	chunks[2].AddDependency(chunks[1])
	chunks[3].AddDependency(chunks[1])

	// m4 +
	c.Test(
		c.Srcs(chunks),
		c.Expected(
			// m1
			"",
			// m2
			"function f(){} f.prototype.bar=3; f.prototype.baz=5;",
			// m3
			"f.prototype.baq = 7;",
			// m4 +
			"var a = new f();"))
}

// Java: CrossChunkCodeMotionTest.testClassMovement6()
func (c *CrossChunkCodeMotionTest) TestClassMovement6() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				// m1
				AddChunk("function Foo(){} function Bar(){} goog.inherits(Bar, Foo); new Foo();").
				// m2
				AddChunk("new Bar();").
				Build()),
		c.Expected(
			// m1
			"function Foo(){} new Foo();",
			// m2
			"function Bar(){} goog.inherits(Bar, Foo); new Bar();"))
}

// Java: CrossChunkCodeMotionTest.testEs6ClassMovement6()
func (c *CrossChunkCodeMotionTest) TestEs6ClassMovement6() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				// m1
				AddChunk("class Foo{} class Bar extends Foo {} new Foo();").
				// m2
				AddChunk("new Bar();").
				Build()),
		c.Expected(
			// m1
			"class Foo {} new Foo();",
			// m2
			"class Bar extends Foo {} new Bar();"))
}

// Java: CrossChunkCodeMotionTest.testClassMovement7()
func (c *CrossChunkCodeMotionTest) TestClassMovement7() {
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				// m1
				AddChunk("function Foo(){} function Bar(){} goog.inherits(Bar, Foo); new Bar();").
				// m2
				AddChunk("new Foo();").
				Build()))
}

// Java: CrossChunkCodeMotionTest.testEs6ClassMovement7()
func (c *CrossChunkCodeMotionTest) TestEs6ClassMovement7() {
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().

				// m1
				AddChunk("class Foo {} class Bar extends Foo {} new Bar();").

				// m2
				AddChunk("new Foo();").
				Build()))
}

// Java: CrossChunkCodeMotionTest.testClassMovement8()
func (c *CrossChunkCodeMotionTest) TestClassMovement8() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				// m1
				AddChunk(
					"function Foo(){}\n"+
						"Object.defineProperties(Foo.prototype, {a: {get:function(){return 0;}}});\n").

				// m2
				AddChunk("new Foo();").
				Build()),
		c.Expected(
			"", // m1
			"function Foo(){}\n"+
				"Object.defineProperties(Foo.prototype, {a: {get:function(){return 0;}}});\n"+
				"new Foo();\n", // m2
		))
}

// Java: CrossChunkCodeMotionTest.testEs6ClassMovement8()
func (c *CrossChunkCodeMotionTest) TestEs6ClassMovement8() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("class Foo { get test() { return 0; }}").
				AddChunk("new Foo();").
				Build()),
		c.Expected(
			"", // m1
			"class Foo { get test() { return 0; }} new Foo();", // m2
		))
}

// Java: CrossChunkCodeMotionTest.testClassMovement_classFields()
func (c *CrossChunkCodeMotionTest) TestClassMovement_classFields() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("class Foo { a; static b = 2; ['c'] = 3; static 'd' = 'hi'; 1 = 2;}").
				AddChunk("new Foo();").
				Build()),
		c.Expected(
			"", "class Foo { a; static b = 2; ['c'] = 3; static 'd' = 'hi'; 1 = 2; } new Foo();"))
}

// Java: CrossChunkCodeMotionTest.testClassMovement_classStaticBlock1()
func (c *CrossChunkCodeMotionTest) TestClassMovement_classStaticBlock1() {
	// TODO(bradfordcsmith):Ideally the class would move
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("class Foo { static { } }").
				AddChunk("new Foo();").
				Build()),
		c.Expected("class Foo { static { } }", "new Foo();"))
}

// Java: CrossChunkCodeMotionTest.testClassMovement_classStaticBlock2()
func (c *CrossChunkCodeMotionTest) TestClassMovement_classStaticBlock2() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("var x = 1;").
				AddChunk("class Foo { static { x; } } new Foo();").
				Build()),
		c.Expected("", "var x = 1; class Foo { static { x; } } new Foo();"))
}

// Java: CrossChunkCodeMotionTest.testClassMovement_classStaticBlock3()
func (c *CrossChunkCodeMotionTest) TestClassMovement_classStaticBlock3() {
	// TODO(bradfordcsmith):Ideally the class and var would move to m3
	chunks :=
		JSChunkGraphBuilder_ForChain().
			// m1
			AddChunk("const x = 1; var y = 2;").
			// m2
			AddChunk("class Foo { static { y = 3; } }").
			// m3
			AddChunk("new Foo();").
			Build()

	c.Test(
		c.Srcs(chunks),
		c.Expected(
			// m1
			"const x = 1;",
			// m2
			"var y =2; class Foo { static { y = 3; } } ",
			// m3
			"new Foo();"))
}

// Java: CrossChunkCodeMotionTest.testClassMovement_classStaticBlock4()
func (c *CrossChunkCodeMotionTest) TestClassMovement_classStaticBlock4() {
	chunks :=
		JSChunkGraphBuilder_ForChain().
			// m1
			AddChunk("var x =1;").
			// m2
			AddChunk(
				"class Foo {\n" +
					"  static {\n" +
					"    x = 2;\n" +
					"  }\n" +
					"}\n" +
					"use(x);\n").
			// m3
			AddChunk("new Foo();").
			Build()

	c.Test(
		c.Srcs(chunks),
		c.Expected(
			// m1
			"",
			// m2
			"var x =1;\n"+
				"class Foo {\n"+
				"  static {\n"+
				"    x = 2\n"+
				"  }\n"+
				"}\n"+
				"use(x);\n",
			// m3
			"new Foo();"))
}

// Java: CrossChunkCodeMotionTest.testClassMovement_mixins()
func (c *CrossChunkCodeMotionTest) TestClassMovement_mixins() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"class ValueType {}\n"+
						"class Foo {}\n"+
						"ValueType.mixin(Foo, ValueType, 5, goog.reflect.objectProperty('foo', Foo))\n").
				AddChunk("new Foo();").
				Build()),
		c.Expected(
			"", // m1
			"class ValueType {}\n"+
				"class Foo {}\n"+
				"ValueType.mixin(Foo, ValueType, 5, goog.reflect.objectProperty('foo', Foo))\n"+
				"new Foo();\n", // m2
		))
}

// Java: CrossChunkCodeMotionTest.testPureOrBreakMyCodedStaticClassFieldIsMovable()
func (c *CrossChunkCodeMotionTest) TestPureOrBreakMyCodedStaticClassFieldIsMovable() {
	//
	//
	//
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"function wrapIt(val) { return [val]; }\n"+
						"class ImmovableClass { static x = wrapIt(123); }\n"+
						"class MovableClass { static y = /** @pureOrBreakMyCode */ wrapIt(123); }\n").
				AddChunk("new ImmovableClass(); new MovableClass();").
				Build()),
		c.Expected(
			"function wrapIt(val) { return [val] }\n"+
				"class ImmovableClass { static x = wrapIt(123);}\n",
			"class MovableClass { static y = /** @pureOrBreakMyCode */ wrapIt(123); }\n"+
				"new ImmovableClass(); new MovableClass();\n"))
}

// Java: CrossChunkCodeMotionTest.testStubMethodMovement()
func (c *CrossChunkCodeMotionTest) TestStubMethodMovement() {
	// The method stub can move, but the unstub definition cannot, because
	// CrossChunkCodeMotion doesn't know where individual methods are used.
	// CrossChunkMethodMotion is responsible for putting the unstub definitions
	// in the right places.
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				// m0
				AddChunk("function Foo(){} Foo.prototype.bar = JSCompiler_stubMethod(x);").
				// m1
				AddChunk("Foo.prototype.bar = JSCompiler_unstubMethod(x);").
				// m2
				AddChunk("new Foo();").
				Build()),
		c.Expected(
			// m0
			"",
			// m1
			"function Foo(){} Foo.prototype.bar = JSCompiler_stubMethod(x);\n"+
				"Foo.prototype.bar = JSCompiler_unstubMethod(x);\n",
			// m2
			"new Foo();"))
}

// Java: CrossChunkCodeMotionTest.testNoMoveSideEffectProperty()
func (c *CrossChunkCodeMotionTest) TestNoMoveSideEffectProperty() {
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				// m1
				AddChunk("function Foo(){}  Foo.prototype.bar = createSomething();").
				// m2
				AddChunk("new Foo();").
				Build()))
}

// Java: CrossChunkCodeMotionTest.testNoMoveSideEffectDefineProperties()
func (c *CrossChunkCodeMotionTest) TestNoMoveSideEffectDefineProperties() {
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				// m1
				AddChunk(
					"function Foo(){}\n" +
						"Object.defineProperties(Foo.prototype, {a: {get: createSomething()}})\n").
				// m2
				AddChunk("new Foo();").
				Build()))
}

// Java: CrossChunkCodeMotionTest.testNoMoveSideEffectDefinePropertiesComputed()
func (c *CrossChunkCodeMotionTest) TestNoMoveSideEffectDefinePropertiesComputed() {
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"function Foo(){}\n" +
						"Object.defineProperties(Foo.prototype,{[test()]:{get: function() {return 10;}}})\n").
				AddChunk("new Foo();").
				Build()))
}

// Java: CrossChunkCodeMotionTest.testAssignMovement()
func (c *CrossChunkCodeMotionTest) TestAssignMovement() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				// m1
				AddChunk("var f = 3; f = 5;").
				// m2
				AddChunk("var h = f;").
				Build()),
		c.Expected(
			// m1
			"",
			// m2
			"var f = 3; f = 5; var h = f;"))

	// don't move nested assigns
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("var f = 3; var g = f = 5;").
				AddChunk("var h = f;").
				Build()))
}

// Java: CrossChunkCodeMotionTest.testNoClassMovement2()
func (c *CrossChunkCodeMotionTest) TestNoClassMovement2() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				// m1
				AddChunk("var f = {}; f.h = 5;").
				// m2
				AddChunk("var h = f;").
				Build()),
		c.Expected(
			// m1
			"",
			// m2
			"var f = {}; f.h = 5; var h = f;"))

	// don't move nested getprop assigns
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("var f = {}; var g = f.h = 5;").
				AddChunk("var h = f;").
				Build()))
}

// Java: CrossChunkCodeMotionTest.testLiteralMovement1()
func (c *CrossChunkCodeMotionTest) TestLiteralMovement1() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				// m1
				AddChunk("var f = {'hi': 'mom', 'bye': function() {}};").
				// m2
				AddChunk("var h = f;").
				Build()),
		c.Expected(
			// m1
			"",
			// m2
			"var f = {'hi': 'mom', 'bye': function() {}}; var h = f;"))
}

// Java: CrossChunkCodeMotionTest.testLiteralMovement2()
func (c *CrossChunkCodeMotionTest) TestLiteralMovement2() {
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("var f = {'hi': 'mom', 'bye': shared};").
				AddChunk("var h = f;").
				Build()))
}

// Java: CrossChunkCodeMotionTest.testLiteralMovement3()
func (c *CrossChunkCodeMotionTest) TestLiteralMovement3() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				// m1
				AddChunk("var f = ['hi', function() {}];").
				// m2
				AddChunk("var h = f;").
				Build()),
		c.Expected(
			// m1
			"",
			// m2
			"var f = ['hi', function() {}]; var h = f;"))
}

// Java: CrossChunkCodeMotionTest.testLiteralMovement4()
func (c *CrossChunkCodeMotionTest) TestLiteralMovement4() {
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("var f = ['hi', shared];").
				AddChunk("var h = f;").
				Build()))
}

// Java: CrossChunkCodeMotionTest.testStringTemplateLiteralMovement1()
func (c *CrossChunkCodeMotionTest) TestStringTemplateLiteralMovement1() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				// m1
				AddChunk("var s = 'world'; var f = `hi ${s}`;").
				// m2
				AddChunk("var h = f;").
				Build()),
		c.Expected(
			// m1
			"",
			// m2
			"var s = 'world'; var f = `hi ${s}`; var h = f;"))
}

// Java: CrossChunkCodeMotionTest.testStringTemplateLiteralMovement2()
func (c *CrossChunkCodeMotionTest) TestStringTemplateLiteralMovement2() {
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("var f = `hi ${shared()}`;").
				AddChunk("var h = f;").
				Build()))
}

// Java: CrossChunkCodeMotionTest.testVarMovement1()
func (c *CrossChunkCodeMotionTest) TestVarMovement1() {
	// test moving a variable
	chunks :=
		JSChunkGraphBuilder_ForStar().
			// m1
			AddChunk("var a = 0;").
			// m2
			AddChunk("var x = a;").
			Build()

	c.Test(
		c.Srcs(chunks),
		c.Expected(
			// m1
			"",
			// m2
			"var a = 0; var x = a;"))
}

// Java: CrossChunkCodeMotionTest.testLetConstMovement()
func (c *CrossChunkCodeMotionTest) TestLetConstMovement() {
	// test moving a variable
	chunks :=
		JSChunkGraphBuilder_ForStar().
			// m1
			AddChunk("const a = 0;").
			// m2
			AddChunk("let x = a;").
			Build()

	c.Test(
		c.Srcs(chunks),
		c.Expected(
			// m1
			"",
			// m2
			"const a = 0; let x = a;"))
}

// Java: CrossChunkCodeMotionTest.testVarMovement2()
func (c *CrossChunkCodeMotionTest) TestVarMovement2() {
	// Test moving 1 variable out of the block
	chunks :=
		JSChunkGraphBuilder_ForStar().
			// m1
			AddChunk("var a = 0; var b = 1; var c = 2;").
			// m2
			AddChunk("var x = b;").
			Build()

	c.Test(
		c.Srcs(chunks),
		c.Expected(
			// m1
			"var a = 0; var c = 2;",
			// m2
			"var b = 1; var x = b;"))
}

// Java: CrossChunkCodeMotionTest.testLetConstMovement2()
func (c *CrossChunkCodeMotionTest) TestLetConstMovement2() {
	// Test moving 1 variable out of the block
	chunks :=
		JSChunkGraphBuilder_ForStar().
			// m1
			AddChunk("const a = 0; const b = 1; const c = 2;").
			// m2
			AddChunk("let x = b;").
			Build()

	c.Test(
		c.Srcs(chunks),
		c.Expected(
			// m1
			"const a = 0; const c = 2;",
			// m2
			"const b = 1; let x = b;"))
}

// Java: CrossChunkCodeMotionTest.testVarMovement3()
func (c *CrossChunkCodeMotionTest) TestVarMovement3() {
	// Test moving all variables out of the block
	chunks :=
		JSChunkGraphBuilder_ForStar().
			// m1
			AddChunk("var a = 0; var b = 1;").
			// m2
			AddChunk("var x = a + b;").
			Build()

	c.Test(
		c.Srcs(chunks),
		c.Expected(
			// m1
			"",
			// m2
			"var a = 0; var b = 1; var x = a + b;"))
}

// Java: CrossChunkCodeMotionTest.testLetConstMovement3()
func (c *CrossChunkCodeMotionTest) TestLetConstMovement3() {
	// Test moving all variables out of the block
	chunks :=
		JSChunkGraphBuilder_ForStar().
			// m1
			AddChunk("const a = 0; const b = 1;").
			// m2
			AddChunk("let x = a + b;").
			Build()

	c.Test(
		c.Srcs(chunks),
		c.Expected(
			// m1
			"",
			// m2
			"const a = 0; const b = 1; let x = a + b;"))
}

// Java: CrossChunkCodeMotionTest.testVarMovement4()
func (c *CrossChunkCodeMotionTest) TestVarMovement4() {
	// Test moving a function
	chunks :=
		JSChunkGraphBuilder_ForStar().
			// m1
			AddChunk("var a = function(){alert(1)};").
			// m2
			AddChunk("var x = a;").
			Build()

	c.Test(
		c.Srcs(chunks),
		c.Expected(
			// m1
			"",
			// m2
			"var a = function(){alert(1)}; var x = a;"))
}

// Java: CrossChunkCodeMotionTest.testLetConstMovement4()
func (c *CrossChunkCodeMotionTest) TestLetConstMovement4() {
	// Test moving a function
	chunks :=
		JSChunkGraphBuilder_ForStar().
			// m1
			AddChunk("const a = function(){alert(1)};").
			// m2
			AddChunk("let x = a;").
			Build()

	c.Test(
		c.Srcs(chunks),
		c.Expected(
			// m1
			"",
			// m2
			"const a = function(){alert(1)}; let x = a;"))
}

// Java: CrossChunkCodeMotionTest.testVarMovement5()
func (c *CrossChunkCodeMotionTest) TestVarMovement5() {
	// Don't move a function outside of scope
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk("var a = alert;").
				AddChunk("var x = a;").
				Build()))
}

// Java: CrossChunkCodeMotionTest.testLetConstMovement5()
func (c *CrossChunkCodeMotionTest) TestLetConstMovement5() {
	// Don't move a function outside of scope
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk("const a = alert;").
				AddChunk("let x = a;").
				Build()))
}

// Java: CrossChunkCodeMotionTest.testVarMovement6()
func (c *CrossChunkCodeMotionTest) TestVarMovement6() {
	// Test moving a var with no assigned value
	chunks :=
		JSChunkGraphBuilder_ForStar().
			// m1
			AddChunk("var a;").
			// m2
			AddChunk("var x = a;").
			Build()

	c.Test(
		c.Srcs(chunks),
		c.Expected(
			// m1
			"",
			// m2
			"var a; var x = a;"))
}

// Java: CrossChunkCodeMotionTest.testLetMovement6()
func (c *CrossChunkCodeMotionTest) TestLetMovement6() {
	// Test moving a let with no assigned value
	chunks :=
		JSChunkGraphBuilder_ForStar().
			// m1
			AddChunk("let a;").
			// m2
			AddChunk("let x = a;").
			Build()

	c.Test(
		c.Srcs(chunks),
		c.Expected(
			// m1
			"",
			// m2
			"let a; let x = a;"))
}

// Java: CrossChunkCodeMotionTest.testVarMovement7()
func (c *CrossChunkCodeMotionTest) TestVarMovement7() {
	// Don't move a variable higher in the dependency tree
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk("function f() {g();} f();").
				AddChunk("function g(){};").
				Build()))
}

// Java: CrossChunkCodeMotionTest.testVarMovement8()
func (c *CrossChunkCodeMotionTest) TestVarMovement8() {
	chunks :=
		JSChunkGraphBuilder_ForBush().
			// m1
			AddChunk("var a = 0;").
			// m2 -> m1
			AddChunk("").
			// m3 -> m2
			AddChunk("var x = a;").
			// m4 -> m2
			AddChunk("var y = a;").
			Build()

	c.Test(
		c.Srcs(chunks),
		c.Expected(
			// m1
			"",
			// m2
			"var a = 0;",
			// m3
			"var x = a;",
			// m4
			"var y = a;"))
}

// Java: CrossChunkCodeMotionTest.testLetConstMovement8()
func (c *CrossChunkCodeMotionTest) TestLetConstMovement8() {
	chunks :=
		JSChunkGraphBuilder_ForBush().
			// m1
			AddChunk("const a = 0;").
			// m2 -> m1
			AddChunk("").
			// m3 -> m2
			AddChunk("let x = a;").
			// m4 -> m2
			AddChunk("let y = a;").
			Build()

	c.Test(
		c.Srcs(chunks),
		c.Expected(
			// m1
			"",
			// m2
			"const a = 0;",
			// m3
			"let x = a;",
			// m4
			"let y = a;"))
}

// Java: CrossChunkCodeMotionTest.testVarMovement9()
func (c *CrossChunkCodeMotionTest) TestVarMovement9() {
	chunks :=
		JSChunkGraphBuilder_ForTree().
			// m1
			AddChunk("var a = 0; var b = 1; var c = 3;").
			// m2 -> m1
			AddChunk("").
			// m3 -> m1
			AddChunk("").
			// m4 -> m2
			AddChunk("a;").
			// m5 -> m2
			AddChunk("a;c;").
			// m6 -> m3
			AddChunk("b;").
			// m7 -> m4
			AddChunk("b;c;").
			Build()

	c.Test(
		c.Srcs(chunks),
		c.Expected(
			// m1
			"var c = 3;",
			// m2
			"var a = 0;",
			// m3
			"var b = 1;",
			// m4
			"a;",
			// m5
			"a;c;",
			// m6
			"b;",
			// m7
			"b;c;"))
}

// Java: CrossChunkCodeMotionTest.testConstMovement9()
func (c *CrossChunkCodeMotionTest) TestConstMovement9() {
	chunks :=
		JSChunkGraphBuilder_ForTree().
			// m1
			AddChunk("const a = 0; const b = 1; const c = 3;").
			// m2 -> m1
			AddChunk("").
			// m3 -> m1
			AddChunk("").
			// m4 -> m2
			AddChunk("a;").
			// m5 -> m2
			AddChunk("a;c;").
			// m6 -> m3
			AddChunk("b;").
			// m7 -> m4
			AddChunk("b;c;").
			Build()

	c.Test(
		c.Srcs(chunks),
		c.Expected(
			// m1
			"const c = 3;",
			// m2
			"const a = 0;",
			// m3
			"const b = 1;",
			// m4
			"a;",
			// m5
			"a;c;",
			// m6
			"b;",
			// m7
			"b;c;"))
}

// Java: CrossChunkCodeMotionTest.testClinit1()
func (c *CrossChunkCodeMotionTest) TestClinit1() {
	chunks :=
		JSChunkGraphBuilder_ForChain().
			// m1
			AddChunk("function Foo$clinit() { Foo$clinit = function() {}; }").
			// m2
			AddChunk("Foo$clinit();").
			Build()
	c.Test(
		c.Srcs(chunks),
		c.Expected(
			// m1
			"",
			// m2
			"function Foo$clinit() { Foo$clinit = function() {}; } Foo$clinit();"))
}

// Java: CrossChunkCodeMotionTest.testClinit2()
func (c *CrossChunkCodeMotionTest) TestClinit2() {
	chunks :=
		JSChunkGraphBuilder_ForChain().
			// m1
			AddChunk("var Foo$clinit = function() { Foo$clinit = function() {}; };").
			// m2
			AddChunk("Foo$clinit();").
			Build()
	c.Test(
		c.Srcs(chunks),
		c.Expected(
			// m1
			"",
			// m2
			"var Foo$clinit = function() { Foo$clinit = function() {}; }; Foo$clinit();"))
}

// Java: CrossChunkCodeMotionTest.testClone1()
func (c *CrossChunkCodeMotionTest) TestClone1() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				// m1
				AddChunk("function f(){} f.prototype.clone = function() { return new f };").
				// m2
				AddChunk("var a = (new f).clone();").
				Build()),
		c.Expected(
			// m1
			"",
			// m2
			"function f(){}\n"+
				"f.prototype.clone = function() { return new f() };\n"+
				"var a = (new f).clone();\n"))
}

// Java: CrossChunkCodeMotionTest.testClone2()
func (c *CrossChunkCodeMotionTest) TestClone2() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"function f(){}\n"+
						"f.prototype.cloneFun = function() {\n"+
						"  return function() {new f}\n"+
						"};\n").
				AddChunk("var a = (new f).cloneFun();").
				Build()),
		c.Expected(
			// m1
			"",
			"function f(){}\n"+
				"f.prototype.cloneFun = function() {\n"+
				"  return function() {new f}\n"+
				"};\n"+
				"var a = (new f).cloneFun();\n"))
}

// Java: CrossChunkCodeMotionTest.testBug4118005()
func (c *CrossChunkCodeMotionTest) TestBug4118005() {
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"var m = 1;\n" +
						"(function () {\n" +
						" var x = 1;\n" +
						" m = function() { return x };\n" +
						"})();\n").
				AddChunk("m();").
				Build()))
}

// Java: CrossChunkCodeMotionTest.testEmptyChunk()
func (c *CrossChunkCodeMotionTest) TestEmptyChunk() {
	// When the dest chunk is empty, it might try to move the code to the
	// one of the chunks that the empty chunk depends on. In some cases
	// this might ended up to be the same chunk as the definition of the code.
	// When that happens, CrossChunkCodeMotion might report a code change
	// while nothing is moved. This should not be a problem if we know all
	// chunks are non-empty.
	m1 := NewJSChunk("m1")
	m1.AddSourceFile(SourceFile_FromCode("m1", "function x() {}"))

	empty := NewJSChunk("empty")
	empty.AddDependency(m1)

	m2 := NewJSChunk("m2")
	m2.AddSourceFile(SourceFile_FromCode("m2", "x()"))
	m2.AddDependency(empty)

	m3 := NewJSChunk("m3")
	m3.AddSourceFile(SourceFile_FromCode("m3", "x()"))
	m3.AddDependency(empty)

	c.Test(c.Srcs([]*JSChunk{m1, empty, m2, m3}), c.Expected("", "function x() {}", "x()", "x()"))
}

// Java: CrossChunkCodeMotionTest.testAbstractMethod()
func (c *CrossChunkCodeMotionTest) TestAbstractMethod() {
	// m2 -> m1
	// m3 -> m1
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				// m1
				AddChunk(
					"var abstractMethod = function () {};\n"+
						"function F(){} F.prototype.bar=abstractMethod;\n"+
						"function G(){} G.prototype.bar=abstractMethod;\n").
				// m2 -> m1
				AddChunk("var f = new F();").
				// m3 -> m1
				AddChunk("var g = new G();").
				Build()),
		c.Expected(
			"var abstractMethod = function () {};",
			"function F(){} F.prototype.bar=abstractMethod; var f = new F();",
			"function G(){} G.prototype.bar=abstractMethod; var g = new G();"))
}

// Java: CrossChunkCodeMotionTest.testMovableUseBeforeDeclaration()
func (c *CrossChunkCodeMotionTest) TestMovableUseBeforeDeclaration() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				// m0
				AddChunk("function f() { g(); } function g() {}").
				// m1
				AddChunk("f();").
				Build()),
		c.Expected(
			// m0
			"",
			// m1
			"function g() {} function f() { g(); } f();"))
}

// Java: CrossChunkCodeMotionTest.testImmovableUseBeforeDeclaration()
func (c *CrossChunkCodeMotionTest) TestImmovableUseBeforeDeclaration() {
	// must recognize this as a reference to the following declaration
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"g(); // must recognize this as a reference to the following declaration\n" +
						"function g() {}\n").
				AddChunk("g();").
				Build()))

	// f() cannot move, so neither can g()
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				// m1
				AddChunk(
					"function f() { g(); }\n" +
						"function g() {}\n" +
						"// f() cannot move, so neither can g()\n" +
						"f();\n").
				// m2
				AddChunk("g();").
				Build()))
}

// Java: CrossChunkCodeMotionTest.testSplitDeclaration()
func (c *CrossChunkCodeMotionTest) TestSplitDeclaration() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				// m0
				AddChunk(
					"function a() { b(); }\n"+
						"function b() {}\n"+
						"function c() {}\n"+
						"a.prototype.x = function() { c(); };\n").
				// m1
				AddChunk("a();").
				Build()),
		c.Expected(
			// m0
			"",
			// m1
			"function b() {}\n"+
				"function c() {}\n"+
				"function a() { b(); }\n"+
				"a.prototype.x = function() { c(); };\n"+
				"a();\n"))
}

// Java: CrossChunkCodeMotionTest.testOutOfOrderAfterSplitDeclaration()
func (c *CrossChunkCodeMotionTest) TestOutOfOrderAfterSplitDeclaration() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				// m0
				AddChunk(
					"function a() { c(); }\n"+
						"function b() {}\n"+
						"a.prototype.x = function() { c(); };\n"+
						"function c() {}\n").
				// m1
				AddChunk("a();").
				Build()),
		c.Expected(
			// m0
			"function b() {}",
			// m1
			"function c() {}\n"+
				"function a() { c(); }\n"+
				"a.prototype.x = function() { c(); };\n"+
				"a();\n"))
}

// Java: CrossChunkCodeMotionTest.testOutOfOrderWithInterveningReferrer()
func (c *CrossChunkCodeMotionTest) TestOutOfOrderWithInterveningReferrer() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				// m0
				AddChunk("function a() { c(); } function b() { a(); } function c() {}").
				// m1
				AddChunk("b();").
				Build()),
		c.Expected(
			// m0
			"",
			// m1
			"function c() {} function a() { c(); } function b() { a(); } b();"))
}

// Java: CrossChunkCodeMotionTest.testOutOfOrderWithDifferentReferrers()
func (c *CrossChunkCodeMotionTest) TestOutOfOrderWithDifferentReferrers() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				// m0
				AddChunk("function a() { b(); } function b() {}").
				// m1
				AddChunk("b();").
				// m2
				AddChunk("a();").
				Build()),
		c.Expected(
			// m0
			"",
			// m1
			"function b() { } b();",
			"function a() { b(); } a();"))
}

// Java: CrossChunkCodeMotionTest.testCircularWithDifferentReferrers()
func (c *CrossChunkCodeMotionTest) TestCircularWithDifferentReferrers() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				// m0
				AddChunk("function a() { b(); } function b() { a(); }").
				// m1
				AddChunk("b();").
				// m2
				AddChunk("a();").
				Build()),
		c.Expected(
			// m0
			"",
			// m1
			"function a() { b(); } function b() { a(); } b();",
			"a();"))
}

// Java: CrossChunkCodeMotionTest.testSmallestCoveringDependencyDoesNotDependOnDeclarationChunk()
func (c *CrossChunkCodeMotionTest) TestSmallestCoveringDependencyDoesNotDependOnDeclarationChunk() {
	//       m0
	//      /  \
	//    m1   m2  // declaration in m1
	//    |    /|  // smallest common dep is m2
	//    m3_ | |  // best place for declaration is m3
	//   / | X  |
	//  /  |/ \ /
	// m4  m5  m6  // references in m5 and m6
	m :=
		JSChunkGraphBuilder_ForUnordered().
			AddChunk("").                // m0
			AddChunk("function f() {}"). // m1
			AddChunk("").                // m2
			AddChunk("").                // m3
			AddChunk("").                // m4
			AddChunk("f();").            // m5
			AddChunk("f();").            // m6
			Build()

	m[1].AddDependency(m[0])
	m[2].AddDependency(m[0])
	m[3].AddDependency(m[1])
	m[4].AddDependency(m[3])
	m[5].AddDependency(m[2])
	m[5].AddDependency(m[3])
	m[6].AddDependency(m[2])
	m[6].AddDependency(m[3])

	c.Test(
		c.Srcs(m),
		c.Expected(
			// m0
			"",
			// m1
			"",
			// m2
			"",
			// m3
			"function f() {}",
			// m4
			"",
			// m5
			"f();",
			// m6
			"f();"))
}

// Java: CrossChunkCodeMotionTest.testEarlyReferencesPinLateDeclarations()
func (c *CrossChunkCodeMotionTest) TestEarlyReferencesPinLateDeclarations() {
	// globalC.x == 2 - not safe to move declaration
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("function C() {} C.prototype.x = 1; var globalC = new C();").
				// globalC.x == 2 - not safe to move declaration
				AddChunk("C.prototype.x = 2; globalC.x;").
				AddChunk("new C().x;").
				Build()))
}

// Java: CrossChunkCodeMotionTest.testMovedInstanceofIsHandledCorrectly()
func (c *CrossChunkCodeMotionTest) TestMovedInstanceofIsHandledCorrectly() {
	c.parentChunkCanSeeSymbolsDeclaredInChildren = true
	// no need to guard instanceof
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				// m0
				AddChunk(
					"function C() {}\n"+
						"function X() {}\n"+
						"X.prototype.a = function(x) { return x instanceof C; }\n").
				// m1
				AddChunk("new C();").
				// m2
				AddChunk("new X();").
				Build()),
		c.Expected(
			// m0
			"",
			// m1
			"function C() {} new C();",
			// m2
			"function X() {}\n"+
				"// no need to guard instanceof\n"+
				"X.prototype.a = function(x) { return x instanceof C; }\n"+
				"new X();\n"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				// m0
				AddChunk(
					"function C() {}\n"+
						"function X() {}\n"+
						"X.prototype.a = function(x) { return x instanceof C; }\n").
				// m1
				AddChunk("new X();").
				// m2
				AddChunk("new C();").
				Build()),
		c.Expected(
			// m0
			"",
			// m1
			"function X() {}\n"+
				"X.prototype.a = function(x) { return 'function' == typeof C && x instanceof C; }\n"+
				"new X();\n",
			// m2
			"function C() {} new C();"))
}

// Java: CrossChunkCodeMotionTest.testValueNotWellDefined()
func (c *CrossChunkCodeMotionTest) TestValueNotWellDefined() {
	// TODO(bradfordcsmith): This prevents us from guaranteeing that all moves are made in a single
	//     pass. Code movement in the first pass may cause some variables to become "well defined"
	//     that weren't before, unblocking movement of some statements.
	// B is blocked from moving because A is used before it is defined.
	// See ReferenceCollection#isWellDefined and CrossChunkReferenceCollector#canMoveValue
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				// m0
				AddChunk("function f() { return A; } var A = 1; var B = A;").
				// m1
				AddChunk("f(); function f2() { return B; }").
				Build()),
		c.Expected(
			// m0
			"var A = 1; var B = A;",
			// m1
			"function f() { return A; } f(); function f2() { return B; }"))
}

// Java: CrossChunkCodeMotionTest.testDestructuringDeclarationsNotMovable()
func (c *CrossChunkCodeMotionTest) TestDestructuringDeclarationsNotMovable() {
	c.TestSame(
		c.Srcs(JSChunkGraphBuilder_ForChain().AddChunk("const [a] = [];").AddChunk("a;").Build()))
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("const {a} = { a: 1 };").
				AddChunk("a;").
				Build()))
}

// Java: CrossChunkCodeMotionTest.testDestructuringAssignmentsAreReferences()
func (c *CrossChunkCodeMotionTest) TestDestructuringAssignmentsAreReferences() {
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("let a = 1; [a] = [5];").
				AddChunk("a;").
				Build()))
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("let a = 1; ({x: a} = {x: 5});").
				AddChunk("a;").
				Build()))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("let a = 1;").
				AddChunk("[a] = [5];").
				AddChunk("a;").
				Build()),
		c.Expected("", "let a = 1; [a] = [5];", "a;"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("let a = 1;").
				AddChunk("({x: a} = {x: 5});").
				AddChunk("a;").
				Build()),
		c.Expected("", "let a = 1; ({x: a} = {x: 5});", "a;"))
}

// Java: CrossChunkCodeMotionTest.testDefaultParamValuesAreReferences()
func (c *CrossChunkCodeMotionTest) TestDefaultParamValuesAreReferences() {
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("let a = 1; function f(x = a) {} f();").
				AddChunk("a;").
				Build()))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("let a = 1; function f(x = a) {}").
				AddChunk("f();").
				Build()),
		c.Expected("", "let a = 1; function f(x = a) {} f();"))
}

// Java: CrossChunkCodeMotionTest.testSpreadCountsAsAReference()
func (c *CrossChunkCodeMotionTest) TestSpreadCountsAsAReference() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("let a = [];").
				AddChunk("function f(...args) {} f(...a);").
				AddChunk("a;").
				Build()),
		c.Expected("", "let a = []; function f(...args) {} f(...a);", "a;"))
}

// Java: CrossChunkCodeMotionTest.testObjectLiteralMethods()
func (c *CrossChunkCodeMotionTest) TestObjectLiteralMethods() {
	// Object literal methods, getters, and setters are movable and references within them are
	// handled correctly.
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("const a = 1;").
				AddChunk("const o = { foo() {return a;}, get x() {}, set x(v) {a = v;} };").
				AddChunk("o;").
				AddChunk("a;").
				Build()),
		c.Expected(
			"",
			"",
			"const a = 1; const o = { foo() {return a;}, get x() {}, set x(v) {a = v;} }; o;",
			"a;"))
}

// Java: CrossChunkCodeMotionTest.testComputedProperties()
func (c *CrossChunkCodeMotionTest) TestComputedProperties() {
	// Computed properties are movable if the key expression is a literal.
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("let a = { ['something']: 1};").
				AddChunk("a;").
				Build()),
		c.Expected("", "let a = { ['something']: 1}; a;"))

	// Computed properties are movable if the key is a well defined variable
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("const x = 1; let a = { [x]: 1};").
				AddChunk("a;").
				Build()),
		c.Expected("", "const x = 1; let a = { [x]: 1}; a;"))
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("const x = 1; let a = class { [x]() {} };").
				AddChunk("a;").
				Build()),
		c.Expected("", "const x = 1; let a = class { [x]() {} }; a;"))

	// Computed properties are not movable if the key is an unknown variable
	c.TestSame(
		c.Srcs(JSChunkGraphBuilder_ForChain().AddChunk("let a = { [x]: 1};").AddChunk("a;").Build()))
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("let a = class { [x]() {} };").
				AddChunk("a;").
				Build()))

	// Computed properties are not movable if the key is a
	c.TestSame(
		c.Srcs(JSChunkGraphBuilder_ForChain().AddChunk("let a = { [x]: 1};").AddChunk("a;").Build()))

	// references in computed properties are honored
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("let a = 1;").
				AddChunk("let b = { [a + 1]: 2 };").
				AddChunk("a;").
				Build()),
		c.Expected("", "let a = 1; let b = { [a + 1]: 2 };", "a;"))

	//  Well defined variables are moveable
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				// m1
				AddChunk("var a = function() {}; a.prototype.blah = function() {};").
				// m2
				AddChunk("a;").
				Build()),
		c.Expected(
			// m1
			"",
			// m2
			`var a = function() {};
a.prototype.blah = function() {};
a;
`))
}

// Java: CrossChunkCodeMotionTest.testComputedPropertiesUsingSymbol()
func (c *CrossChunkCodeMotionTest) TestComputedPropertiesUsingSymbol() {
	// Computed properties are movable if the key is a `Symbol.<any_property_on_symbol>`
	// `Symbol` is a built-in JavaScript object that is guaranteed to be unique. Therefore, it
	// is movable. Any properties on `Symbol` are also movable (e.g `Symbol.iterator`).
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("let a = { [Symbol.iterator]: 1};").
				AddChunk("a;").
				Build()),
		c.Expected(
			"", // everything has been moved out of this chunk
			"let a = { [Symbol.iterator]: 1}; a;"))

	// `Symbol` is declared in m1, and used in m2. It should be moved to m2 in expected output.
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				// m1
				AddChunk("let a = class { [Symbol.anyPropertyOnSymbolIsMovable]() {} };").
				// m2
				AddChunk("a;").
				Build()),
		c.Expected(
			// m1
			"",
			// m2
			"let a = class { [Symbol.anyPropertyOnSymbolIsMovable]() {} }; a;"))

	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				// m1
				AddChunk("let prop; let a = class { [prop]() {} };").
				// m2
				AddChunk("a;").
				Build()))

	// Now turn on transpilation and rerun the same above 3 unit tests.
	c.EnableTranspile()

	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				// m1
				AddChunk("let a = { [Symbol.iterator]: 1};").
				// m2
				AddChunk("a;").
				Build()),
		c.Expected(
			// m1
			// Ideally, we should be moving this code out of m1, but this would require us to move
			// an `ASSIGN` node which we don't want to do.
			`var $jscomp$compprop98447280$0 = {};
var a = ($jscomp$compprop98447280$0[Symbol.iterator] = 1, $jscomp$compprop98447280$0);
`,
			// m2
			"a;"))

	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				// m1
				AddChunk("let a = class { [Symbol.anyPropertyOnSymbolIsMovable]() {} };").
				// m2
				AddChunk("a;").
				Build()),
		c.Expected(
			// m1
			"",
			// m2
			`/** @constructor */ var a = function() {};
a.prototype[Symbol.anyPropertyOnSymbolIsMovable] = function() {};
a;
`))

	// Non-Symbol computed properties are currently not movable.
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				// m1
				AddChunk("let prop; let a = class { [prop]() {} };").
				// m2
				AddChunk("a;").
				Build()),
		c.Expected(
			// m1
			`var prop;
var $jscomp$compField$98447280$0 = prop;
/** @constructor */ var a = function() {};
a.prototype[$jscomp$compField$98447280$0] = function() {};
`,
			// m2
			"a;"))
}

// Java: CrossChunkCodeMotionTest.testComputedPropertiesUsingSymbolUnmovable()
func (c *CrossChunkCodeMotionTest) TestComputedPropertiesUsingSymbolUnmovable() {
	// Computed properties are not movable if the key is not explicitly called `Symbol.`
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("let a = { [Symbol123.iterator]: 1};").
				AddChunk("a;").
				Build()))

	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				// m1
				AddChunk("let a = class { [Symbol123.anyPropertyOnSymbolIsMovable]() {} };").
				// m2
				AddChunk("a;").
				Build()))
}

// Java: CrossChunkCodeMotionTest.testClassComputedFieldUnmovable()
func (c *CrossChunkCodeMotionTest) TestClassComputedFieldUnmovable() {
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("class C { ['x'] = 1; }").
				AddChunk("a;").
				Build()))
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("class C { static ['x'] = 1; }").
				AddChunk("a;").
				Build()))
}

// Java: CrossChunkCodeMotionTest.testPureOrBreakMyCode()
func (c *CrossChunkCodeMotionTest) TestPureOrBreakMyCode() {
	// should move a function if its top level invocation was annotated with @pureOrBreakMyCode
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"function wrapIt(val) { return [val]; }\n"+
						"const value = /** @pureOrBreakMyCode */ wrapIt(123);\n").
				AddChunk("const unwrappedValue = value[0];").
				Build()),
		c.Expected(
			"",
			"function wrapIt(val) { return [val]; }\n"+
				"const value = /** @pureOrBreakMyCode */ wrapIt(123);\n"+
				"const unwrappedValue = value[0];\n"))

	// should move a class if its static initializer containing a function call was annotated with
	// @pureOrBreakMyCode
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"function wrapIt(val) { return [val]; }\n"+
						"class ImmovableClass {}\n"+
						"ImmovableClass.prop = wrapIt(123);\n"+
						"class MovableClass {}\n"+
						"MovableClass.prop = /** @pureOrBreakMyCode */ wrapIt(123);\n").
				AddChunk("new ImmovableClass(); new MovableClass();").
				Build()),
		c.Expected(
			"function wrapIt(val) { return [val] }\n"+
				"class ImmovableClass {}\n"+
				"ImmovableClass.prop = wrapIt(123);\n",
			"class MovableClass {}\n"+
				"MovableClass.prop = /** @pureOrBreakMyCode */ wrapIt(123);\n"+
				"new ImmovableClass(); new MovableClass();\n"))

	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"class LowerCasePipe {}\n"+
						"/** @nocollapse */ LowerCasePipe.ɵpipe = /** @pureOrBreakMyCode*/ i0.ɵɵdefinePipe({ name: \"lowercase\", type: LowerCasePipe, pure: true });\n").
				AddChunk("new LowerCasePipe();").
				Build()),
		c.Expected(
			"",
			"class LowerCasePipe {}\n"+
				"/** @nocollapse */ LowerCasePipe.ɵpipe = /** @pureOrBreakMyCode*/ i0.ɵɵdefinePipe({ name: \"lowercase\", type: LowerCasePipe, pure: true });\n"+
				"new LowerCasePipe();\n"))
}
