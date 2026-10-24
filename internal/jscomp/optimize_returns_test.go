// Port of com/google/javascript/jscomp/OptimizeReturnsTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/OptimizeReturnsTest.java)

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
)

type OptimizeReturnsTest struct {
	CompilerTestCaseBase
}

// Java: OptimizeReturnsTest()
func NewOptimizeReturnsTest() *OptimizeReturnsTest {
	c := &OptimizeReturnsTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, CompilerTestCase_DEFAULT_EXTERNS+OptimizeReturnsTest_EXTERNAL_SYMBOLS)
	return c
}

func TestOptimizeReturnsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewOptimizeReturnsTest() })
}

// Java: OptimizeReturnsTest.getProcessor(Compiler)
func (c *OptimizeReturnsTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewOptimizeReturns(compiler)
}

const OptimizeReturnsTest_EXTERNAL_SYMBOLS = "var extern;\n" +
	"extern.externalMethod\n"

// Java: OptimizeReturnsTest.setUp()
func (c *OptimizeReturnsTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableNormalize() // Required for `OptimizeCalls`.
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.EnableGatherExternProperties()
}

// Java: OptimizeReturnsTest.nullishCoalesceReturnRemoved()
func (c *OptimizeReturnsTest) TestNullishCoalesceReturnRemoved() {
	c.Test(
		"var f = (function() {return 1}) ?? (function() {return 2}); f();",
		"var f = function() { return; } ?? function() { return; }; f();")
}

// Java: OptimizeReturnsTest.testNoRewriteUsedResult1()
func (c *OptimizeReturnsTest) TestNoRewriteUsedResult1() {
	source :=
		"function a(){return 1}\n" +
			"var x = a()\n"
	c.TestSame(source)
}

// Java: OptimizeReturnsTest.testNoRewriteUsedResult2()
func (c *OptimizeReturnsTest) TestNoRewriteUsedResult2() {
	source :=
		"var a = function(){return 1}\n" +
			"a(); var b = a()\n"
	c.TestSame(source)
}

// Java: OptimizeReturnsTest.testNoRewriteUsedClassMethodResult1()
func (c *OptimizeReturnsTest) TestNoRewriteUsedClassMethodResult1() {
	source :=
		"class C { method() {return 1} }\n" +
			"var x = new C().method()\n"
	c.TestSame(source)
}

// Java: OptimizeReturnsTest.testNoRewriteUnsedObjectMethodResult1()
func (c *OptimizeReturnsTest) TestNoRewriteUnsedObjectMethodResult1() {
	source :=
		"var o = { method() {return 1} }\n" +
			"o.method()\n"
	c.TestSame(source)
}

// Java: OptimizeReturnsTest.testNoRewriteDestructuredArray1()
func (c *OptimizeReturnsTest) TestNoRewriteDestructuredArray1() {
	source :=
		"var x = function() { return 1; };\n" +
			"[x] = []\n" +
			"x()\n"
	c.TestSame(source)
}

// Java: OptimizeReturnsTest.testNoRewriteDestructuredArray2()
func (c *OptimizeReturnsTest) TestNoRewriteDestructuredArray2() {
	source :=
		"var x = function() { return 1; };\n" +
			"[x = function() {}] = []\n" +
			"x()\n"
	c.TestSame(source)
}

// Java: OptimizeReturnsTest.testNoRewriteDestructuredArray3()
func (c *OptimizeReturnsTest) TestNoRewriteDestructuredArray3() {
	source :=
		"class C { method() { return 1 }}\n" +
			"[x.method] = []\n" +
			"x()\n"
	c.TestSame(source)
}

// Java: OptimizeReturnsTest.testNoRewriteDestructuredArray4()
func (c *OptimizeReturnsTest) TestNoRewriteDestructuredArray4() {
	source :=
		"class C { method() { return 1 }}\n" +
			"[x.method] = []\n" +
			"x()\n"
	c.TestSame(source)
}

// Java: OptimizeReturnsTest.testNoRewriteDestructuredObject1()
func (c *OptimizeReturnsTest) TestNoRewriteDestructuredObject1() {
	source :=
		"var x = function() { return 1; };\n" +
			"({a:x} = {})\n" +
			"x()\n"
	c.TestSame(source)
}

// Java: OptimizeReturnsTest.testNoRewriteDestructuredObject2()
func (c *OptimizeReturnsTest) TestNoRewriteDestructuredObject2() {
	source :=
		"var x = function() { return 1; };\n" +
			"({a:x = function() {}} = {})\n" +
			"x()\n"
	c.TestSame(source)
}

// Java: OptimizeReturnsTest.testNoRewriteDestructuredObject3()
func (c *OptimizeReturnsTest) TestNoRewriteDestructuredObject3() {
	source :=
		"class C { method() { return 1 }}\n" +
			"({a:x.method} = {})\n" +
			"x()\n"
	c.TestSame(source)
}

// Java: OptimizeReturnsTest.testNoRewriteDestructuredObject4()
func (c *OptimizeReturnsTest) TestNoRewriteDestructuredObject4() {
	source :=
		"class C { method() { return 1 }}\n" +
			"({a:x.method = function() {}} = {})\n" +
			"x()\n"
	c.TestSame(source)
}

// Java: OptimizeReturnsTest.testNoRewriteTagged1()
func (c *OptimizeReturnsTest) TestNoRewriteTagged1() {
	// TODO(johnlenz): support this. Unused return can be removed.
	source :=
		"var f = function() { return 1; };\n" +
			"f`tagged`\n"
	c.TestSame(source)
}

// Java: OptimizeReturnsTest.testNoRewriteTagged2()
func (c *OptimizeReturnsTest) TestNoRewriteTagged2() {
	// Tagged use prevents optimizations
	source :=
		"var f = function() { return 1; };\n" +
			"var x = f`tagged`\n"
	c.TestSame(source)
}

// Java: OptimizeReturnsTest.testNoRewriteTagged3()
func (c *OptimizeReturnsTest) TestNoRewriteTagged3() {
	// Tagged use is not ignored.
	source :=
		"var f = function() { return 1; };\n" +
			"var x = f`tagged`\n" +
			"f()\n"
	c.TestSame(source)
}

// Java: OptimizeReturnsTest.testNoRewriteConstructorProp()
func (c *OptimizeReturnsTest) TestNoRewriteConstructorProp() {
	source :=
		"class C { constructor() { return 1 } }\n" +
			"x.constructor()\n"
	c.TestSame(source)
}

// Java: OptimizeReturnsTest.testRewriteUnusedResult1()
func (c *OptimizeReturnsTest) TestRewriteUnusedResult1() {
	source :=
		"function a(){return 1}\n" +
			"a()\n"
	expected :=
		"function a(){return}\n" +
			"a()\n"
	c.Test(source, expected)
}

// Java: OptimizeReturnsTest.testRewriteUnusedResult2()
func (c *OptimizeReturnsTest) TestRewriteUnusedResult2() {
	source :=
		"var a; a = function(){return 1}\n" +
			"a()\n"
	expected :=
		"var a; a = function(){return}\n" +
			"a()\n"
	c.Test(source, expected)
}

// Java: OptimizeReturnsTest.testRewriteUnusedResult3()
func (c *OptimizeReturnsTest) TestRewriteUnusedResult3() {
	source :=
		"var a = function(){return 1}\n" +
			"a()\n"
	expected :=
		"var a = function(){return}\n" +
			"a()\n"
	c.Test(source, expected)
}

// Java: OptimizeReturnsTest.testRewriteUnusedResult4a()
func (c *OptimizeReturnsTest) TestRewriteUnusedResult4a() {
	source :=
		"var a = function(){return a()}\n" +
			"a()\n"
	c.TestSame(source)
}

// Java: OptimizeReturnsTest.testRewriteUnusedResult4b()
func (c *OptimizeReturnsTest) TestRewriteUnusedResult4b() {
	source :=
		"var a = function b(){return b()}\n" +
			"a()\n"
	c.TestSame(source)
}

// Java: OptimizeReturnsTest.testNoRewriteWhenAliasedDuringAssignment()
func (c *OptimizeReturnsTest) TestNoRewriteWhenAliasedDuringAssignment() {
	source :=
		"var a, b;\n" +
			"a = b = function (){return 1}\n" +
			"use(a()); // result used\n" +
			"b()\n" // result unused
	c.TestSame(source)
}

// Java: OptimizeReturnsTest.testRewriteUnusedResult4c()
func (c *OptimizeReturnsTest) TestRewriteUnusedResult4c() {
	source :=
		"function a(){return a()}\n" +
			"a()\n"
	c.TestSame(source)
}

// Java: OptimizeReturnsTest.testRewriteUnusedResult5()
func (c *OptimizeReturnsTest) TestRewriteUnusedResult5() {
	source :=
		"function a(){}\n" +
			"a.prototype.foo = function(args) {return args};\n" +
			"var o = new a;\n" +
			"o.foo()\n"
	expected :=
		"function a(){}\n" +
			"a.prototype.foo = function(args) {args;return};\n" +
			"var o = new a;\n" +
			"o.foo()\n"
	c.Test(source, expected)
}

// Java: OptimizeReturnsTest.testRewriteUnusedResult6()
func (c *OptimizeReturnsTest) TestRewriteUnusedResult6() {
	source :=
		"function a(){return (g = 1)}\n" +
			"a()\n"
	expected :=
		"function a(){g = 1;return}\n" +
			"a()\n"
	c.Test(source, expected)
}

// Java: OptimizeReturnsTest.testRewriteUnusedResult7a()
func (c *OptimizeReturnsTest) TestRewriteUnusedResult7a() {
	source :=
		"function a() { return 1 }\n" +
			"function b() { return a() }\n" +
			"function c() { return b() }\n" +
			"c();\n"

	expected :=
		"function a() { return 1 }\n" +
			"function b() { return a() }\n" +
			"function c() { b(); return }\n" +
			"c();\n"
	c.Test(source, expected)
}

// Java: OptimizeReturnsTest.testRewriteUnusedResult7b()
func (c *OptimizeReturnsTest) TestRewriteUnusedResult7b() {
	source :=
		"c();\n" +
			"function c() { return b() }\n" +
			"function b() { return a() }\n" +
			"function a() { return 1 }\n"

	// Iteration 1.
	expected :=
		"c();\n" +
			"function c() { b(); return }\n" +
			"function b() { return a() }\n" +
			"function a() { return 1 }\n"
	c.Test(source, expected)

	// Iteration 2.
	source = expected
	expected =
		"c();\n" +
			"function c() { b(); return }\n" +
			"function b() { a(); return }\n" +
			"function a() { return 1 }\n"
	c.Test(source, expected)

	// Iteration 3.
	source = expected
	expected =
		"c();\n" +
			"function c() { b(); return }\n" +
			"function b() { a(); return }\n" +
			"function a() { return }\n"
	c.Test(source, expected)
}

// Java: OptimizeReturnsTest.testRewriteUnusedResult8()
func (c *OptimizeReturnsTest) TestRewriteUnusedResult8() {
	source :=
		"function a() { return c() }\n" +
			"function b() { return a() }\n" +
			"function c() { return b() }\n" +
			"c();\n"
	c.TestSame(source)
}

// Java: OptimizeReturnsTest.testRewriteUnusedResult9()
func (c *OptimizeReturnsTest) TestRewriteUnusedResult9() {
	// Proves that the deleted function scope is reported.
	source :=
		"function a(){return function() {};}\n" +
			"a()\n"
	expected :=
		"function a(){(function() {}); return}\n" +
			"a()\n"
	c.Test(source, expected)
}

// Java: OptimizeReturnsTest.testRewriteUsedResult10()
func (c *OptimizeReturnsTest) TestRewriteUsedResult10() {
	source :=
		"class C { method() {return 1} }\n" +
			"new C().method()\n"
	expected :=
		"class C { method() {return} }\n" +
			"new C().method()\n"
	c.Test(source, expected)
}

// Java: OptimizeReturnsTest.testRewriteUnusedTemplateLitResult()
func (c *OptimizeReturnsTest) TestRewriteUnusedTemplateLitResult() {
	// Proves that the deleted function scope is reported.
	source :=
		"function a(){ return `template`; }\n" +
			"a()\n"
	expected :=
		"function a(){ return; }\n" +
			"a()\n"
	c.Test(source, expected)
}

// Java: OptimizeReturnsTest.testRewriteUnusedAsyncResult1()
func (c *OptimizeReturnsTest) TestRewriteUnusedAsyncResult1() {
	// Async function returns can be dropped if no-one waits on the returned
	// promise.
	source :=
		"async function a(){return promise}\n" +
			"a()\n"
	expected :=
		"async function a(){promise; return}\n" +
			"a()\n"
	c.Test(source, expected)
}

// Java: OptimizeReturnsTest.testRewriteUnusedGeneratorResult1()
func (c *OptimizeReturnsTest) TestRewriteUnusedGeneratorResult1() {
	// Generator function returns can be dropped if no-one uses the returned
	// iterator.
	source :=
		"function *a(){return value}\n" +
			"a()\n"
	expected :=
		"function *a(){value; return}\n" +
			"a()\n"
	c.Test(source, expected)
}

// Java: OptimizeReturnsTest.testNoRewriteObjLit1()
func (c *OptimizeReturnsTest) TestNoRewriteObjLit1() {
	source :=
		"var a = {b:function(){return 1;}}\n" +
			"for(c in a) (a[c])();\n" +
			"a.b()\n"
	c.TestSame(source)
}

// Java: OptimizeReturnsTest.testNoRewriteObjLit2()
func (c *OptimizeReturnsTest) TestNoRewriteObjLit2() {
	source :=
		"var a = {b:function fn(){return 1;}}\n" +
			"for(c in a) (a[c])();\n" +
			"a.b()\n"
	c.TestSame(source)

	sourceOptChain :=
		"var a = {b:function fn(){return 1;}}\n" +
			"for(c in a) (a[c])?.();\n" +
			"a.b?.()\n"
	c.TestSame(sourceOptChain)
}

// Java: OptimizeReturnsTest.testNoRewriteArrLit()
func (c *OptimizeReturnsTest) TestNoRewriteArrLit() {
	source :=
		"var a = [function(){return 1;}]\n" +
			"(a[0])();\n"
	c.TestSame(source)

	sourceOptChain :=
		"var a = [function(){return 1;}]\n" +
			"(a[0])?.();\n"
	c.TestSame(sourceOptChain)
}

// Java: OptimizeReturnsTest.testPrototypeMethod1()
func (c *OptimizeReturnsTest) TestPrototypeMethod1() {
	source :=
		"function c(){}\n" +
			"c.prototype.a = function(){return 1}\n" +
			"var x = new c;\n" +
			"x.a()\n"
	result :=
		"function c(){}\n" +
			"c.prototype.a = function(){return}\n" +
			"var x = new c;\n" +
			"x.a()\n"
	c.Test(source, result)
}

// Java: OptimizeReturnsTest.testPrototypeMethod2()
func (c *OptimizeReturnsTest) TestPrototypeMethod2() {
	source :=
		"function c(){}\n" +
			"c.prototype.a = function(){return 1}\n" +
			"goog.reflect.object({a: 'v'})\n" +
			"var x = new c;\n" +
			"x.a()\n"
	c.TestSame(source)
}

// Java: OptimizeReturnsTest.testPrototypeMethod3()
func (c *OptimizeReturnsTest) TestPrototypeMethod3() {
	source :=
		"function c(){}\n" +
			"c.prototype.a = function(){return 1}\n" +
			"var x = new c;\n" +
			"for(var key in goog.reflect.object({a: 'v'})){ x[key](); }\n" +
			"x.a()\n"
	c.TestSame(source)
}

// Java: OptimizeReturnsTest.testPrototypeMethod4()
func (c *OptimizeReturnsTest) TestPrototypeMethod4() {
	source :=
		"function c(){}\n" +
			"c.prototype.a = function(){return 1}\n" +
			"var x = new c;\n" +
			"for(var key in goog.reflect.object({a: 'v'})){ x[key](); }\n"
	c.TestSame(source)
}

// Java: OptimizeReturnsTest.testCallOrApply()
func (c *OptimizeReturnsTest) TestCallOrApply() {
	// TODO(johnlenz): Add support for .apply
	c.Test("function a() {return 1}; a.call(new foo);", "function a() {return  }; a.call(new foo);")
	c.TestSame("function a() {return 1}; a.apply(new foo);")
}

// Java: OptimizeReturnsTest.testCallOrApply_optChain()
func (c *OptimizeReturnsTest) TestCallOrApply_optChain() {
	c.Test(
		"function a() {return 1}; a?.call(new foo);", //
		"function a() {return  }; a?.call(new foo);")
	c.TestSame("function a() {return 1}; a?.apply(new foo);")
}

// Java: OptimizeReturnsTest.testRewriteUseSiteRemoval()
func (c *OptimizeReturnsTest) TestRewriteUseSiteRemoval() {
	source :=
		"function a() { return {\"_id\" : 1} }\n" +
			"a();\n"
	expected :=
		"function a() { ({\"_id\" : 1}); return }\n" +
			"a();\n"
	c.Test(source, expected)
}

// Java: OptimizeReturnsTest.testUnknownDefinitionAllowRemoval()
func (c *OptimizeReturnsTest) TestUnknownDefinitionAllowRemoval() {
	// TODO(johnlenz): allow this to be optimized.
	c.TestSame(
		"let x = functionFactory();\n" +
			"x(1, 2);\n" +
			"x = function(a,b) { return b; }\n")
}

// Java: OptimizeReturnsTest.testReturnNotRemovedFromRecursiveNamedFunctionExpression()
func (c *OptimizeReturnsTest) TestReturnNotRemovedFromRecursiveNamedFunctionExpression() {
	c.TestSame(
		"let x = function innerName(n) {\n" +
			"  if (n < 1) {\n" +
			"    return 0\n" +
			"  } else {\n" +
			"    return innerName(n - 1) + n;\n" +
			"  }\n" +
			"\n" +
			"}\n" +
			"x(3);\n")
}

// optional chaining version of the above test
// Java: OptimizeReturnsTest.testReturnNotRemovedFromRecursiveNamedFunctionExpression_optChain()
func (c *OptimizeReturnsTest) TestReturnNotRemovedFromRecursiveNamedFunctionExpression_optChain() {
	c.TestSame(
		"let x = function innerName(n) {\n" +
			"  if (n < 1) {\n" +
			"    return 0\n" +
			"  } else {\n" +
			"    return innerName(n - 1) + n;\n" +
			"  }\n" +
			"\n" +
			"}\n" +
			"x?.(3);\n")
}

// Java: OptimizeReturnsTest.testRewriteUnusedResultWithSafeReference1()
func (c *OptimizeReturnsTest) TestRewriteUnusedResultWithSafeReference1() {
	source :=
		"function a(){return 1}\n" +
			"typeof a\n" +
			"a()\n"
	expected :=
		"function a(){return}\n" +
			"typeof a\n" +
			"a()\n"
	c.Test(source, expected)
}

// Java: OptimizeReturnsTest.testRewriteUnusedResultWithSafeReference2()
func (c *OptimizeReturnsTest) TestRewriteUnusedResultWithSafeReference2() {
	source :=
		"function a(){return 1}\n" +
			"x instanceof a\n" +
			"a instanceof x\n" +
			"a()\n"
	expected :=
		"function a(){return}\n" +
			"x instanceof a\n" +
			"a instanceof x\n" +
			"a()\n"
	c.Test(source, expected)
}

// Java: OptimizeReturnsTest.testRewriteUnusedResultWithSafeReference3()
func (c *OptimizeReturnsTest) TestRewriteUnusedResultWithSafeReference3() {
	source :=
		"function a(){return 1}\n" +
			"x in a\n" +
			"a in x\n" +
			"a()\n"
	expected :=
		"function a(){return}\n" +
			"x in a\n" +
			"a in x\n" +
			"a()\n"
	c.Test(source, expected)
}

// Java: OptimizeReturnsTest.testRewriteUnusedResultWithSafeReference4()
func (c *OptimizeReturnsTest) TestRewriteUnusedResultWithSafeReference4() {
	source :=
		"function a(){return 1}\n" +
			"a.x\n" +
			"a['x']\n" +
			"a()\n"
	expected :=
		"function a(){return}\n" +
			"a.x\n" +
			"a['x']\n" +
			"a()\n"
	c.Test(source, expected)
}

// Java: OptimizeReturnsTest.testRewriteUnusedResultWithSafeReference5()
func (c *OptimizeReturnsTest) TestRewriteUnusedResultWithSafeReference5() {
	source :=
		"function a(){return 1}\n" +
			"for (x in a) {}\n" +
			"for (x of a) {}\n" +
			"a()\n"
	expected :=
		"function a(){return}\n" +
			"for (x in a) {}\n" +
			"for (x of a) {}\n" +
			"a()\n"
	c.Test(source, expected)

	// optional versions
	sourceOptChainCall :=
		"function a(){return 1}\n" +
			"for (x in a) {}\n" +
			"for (x of a) {}\n" +
			"a?.()\n"
	expectedOptChainCall :=
		"function a(){return}\n" +
			"for (x in a) {}\n" +
			"for (x of a) {}\n" +
			"a?.()\n"
	c.Test(sourceOptChainCall, expectedOptChainCall)
}

// Java: OptimizeReturnsTest.testNoRewriteUnusedResultWithUnsafeReference1()
func (c *OptimizeReturnsTest) TestNoRewriteUnusedResultWithUnsafeReference1() {
	// call to 'a.x' escapes 'a' as 'this'
	source :=
		"function a(){return 1}\n" +
			"a.x()\n" +
			"a()\n"
	c.TestSame(source)

	// call to 'a?.x' escapes 'a' as 'this'
	sourceOptChainCall :=
		"function a(){return 1}\n" +
			"a?.x()\n" +
			"a()\n"
	c.TestSame(sourceOptChainCall)
}

// Java: OptimizeReturnsTest.testNoRewriteUnusedResultWithUnsafeReference2()
func (c *OptimizeReturnsTest) TestNoRewriteUnusedResultWithUnsafeReference2() {
	// call to 'a[x]' escapes 'a' as 'this'
	source :=
		"function a(){return 1}\n" +
			"a['x']()\n" +
			"a()\n"
	c.TestSame(source)

	// call to 'a?.[x]' escapes 'a' as 'this'
	sourceOptChainCall :=
		"function a(){return 1}\n" +
			"a?.['x']()\n" +
			"a()\n"
	c.TestSame(sourceOptChainCall)
}

// Java: OptimizeReturnsTest.testNoRewriteUnusedResultWithUnsafeReference4()
func (c *OptimizeReturnsTest) TestNoRewriteUnusedResultWithUnsafeReference4() {
	// call to 'a' is assigned an unknown value
	// TODO(johnlenz): optimize this.
	source :=
		"function a(){return 1}\n" +
			"for (a of x) {}\n" +
			"a()\n"
	c.TestSame(source)
}
