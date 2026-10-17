// Port of com/google/javascript/jscomp/CrossChunkMethodMotionTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/CrossChunkMethodMotionTest.java)

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

type CrossChunkMethodMotionTest struct {
	CompilerTestCaseBase
	canMoveExterns bool
	noStubs        bool
}

// Java: CrossChunkMethodMotionTest()
func NewCrossChunkMethodMotionTest() *CrossChunkMethodMotionTest {
	c := &CrossChunkMethodMotionTest{}
	c.InitCompilerTestCaseBase(c)
	c.canMoveExterns = false
	c.noStubs = false
	return c
}

func TestCrossChunkMethodMotionTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCrossChunkMethodMotionTest() })
}

const CrossChunkMethodMotionTest_STUB_DECLARATIONS = CrossChunkMethodMotion_STUB_DECLARATIONS

// Java: CrossChunkMethodMotionTest.getProcessor(Compiler)
func (c *CrossChunkMethodMotionTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewCrossChunkMethodMotion(compiler, NewIdGenerator(), c.canMoveExterns, c.noStubs)
}

// Java: CrossChunkMethodMotionTest.getOptions()
func (c *CrossChunkMethodMotionTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	// pretty printing makes it much easier to read the failure messages.
	options.SetPrettyPrint(true)
	return options
}

// JUnitBefore runs the @Before methods: CompilerTestCase.setUp() (the
// superclass one) first, then customSetUp().
func (c *CrossChunkMethodMotionTest) JUnitBefore() {
	c.CompilerTestCaseBase.JUnitBefore()
	c.CustomSetUp()
}

// Java: CrossChunkMethodMotionTest.customSetUp()
func (c *CrossChunkMethodMotionTest) CustomSetUp() {
	c.canMoveExterns = false
	c.noStubs = false
	c.EnableNormalize()
}

// Java: CrossChunkMethodMotionTest.moveMethodAssignedToPrototype()
func (c *CrossChunkMethodMotionTest) TestMoveMethodAssignedToPrototype() {
	c.TestSame(
		// bar property is defined in externs, so it cannot be moved
		c.Externs("IFoo.prototype.bar;"),
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"function Foo() {}\n"+
						"Foo.prototype.bar = function() {};\n").
				AddChunk("(new Foo).bar()").
				Build()))

	c.canMoveExterns = true
	c.Test(
		c.Externs("IFoo.prototype.bar;"),
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"function Foo() {}\n"+
						"Foo.prototype.bar = function() {};\n").
				AddChunk("(new Foo).bar()").
				Build()),
		c.Expected(
			CrossChunkMethodMotionTest_STUB_DECLARATIONS+
				"function Foo() {}\n"+
				"Foo.prototype.bar = JSCompiler_stubMethod(0);\n",
			// Chunk 2
			"Foo.prototype.bar = JSCompiler_unstubMethod(0, function() {});\n"+
				"(new Foo).bar()\n"))
}

// Java: CrossChunkMethodMotionTest.moveMethodDefinedInPrototypeLiteralWithStubs()
func (c *CrossChunkMethodMotionTest) TestMoveMethodDefinedInPrototypeLiteralWithStubs() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"function Foo() {}\n"+
						"Foo.prototype = { method: function() {} };\n").
				AddChunk("(new Foo).method()").
				Build()),
		c.Expected(
			CrossChunkMethodMotionTest_STUB_DECLARATIONS+
				"function Foo() {}\n"+
				"Foo.prototype = { method: JSCompiler_stubMethod(0) };\n",
			// Chunk 2
			"Foo.prototype.method =\n"+
				"    JSCompiler_unstubMethod(0, function() {});\n"+
				"(new Foo).method()\n"))
}

// Java: CrossChunkMethodMotionTest.moveMethodDefinedInPrototypeLiteralWithoutStubs()
func (c *CrossChunkMethodMotionTest) TestMoveMethodDefinedInPrototypeLiteralWithoutStubs() {
	c.noStubs = true
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"function Foo() {}\n"+
						"Foo.prototype = { method: function() {} };\n").
				AddChunk("(new Foo).method()").
				Build()),
		c.Expected(
			"function Foo() {}\n"+
				"Foo.prototype = {};\n",
			// Chunk 2
			"Foo.prototype.method = function() {};\n"+
				"(new Foo).method()\n"))
}

// Java: CrossChunkMethodMotionTest.moveMethodDefinedInPrototypeLiteralUsingShorthandSyntaxWithStub()
func (c *CrossChunkMethodMotionTest) TestMoveMethodDefinedInPrototypeLiteralUsingShorthandSyntaxWithStub() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"function Foo() {}\n"+
						"Foo.prototype = { method() {} };\n").
				AddChunk("(new Foo).method()").
				Build()),
		c.Expected(
			CrossChunkMethodMotionTest_STUB_DECLARATIONS+
				"function Foo() {}\n"+
				"Foo.prototype = { method: JSCompiler_stubMethod(0) };\n",
			// Chunk 2
			"Foo.prototype.method =\n"+
				"    JSCompiler_unstubMethod(0, function() {});\n"+
				"(new Foo).method()\n"))
}

// Java: CrossChunkMethodMotionTest.moveMethodDefinedInPrototypeLiteralUsingShorthandSyntaxWithoutStub()
func (c *CrossChunkMethodMotionTest) TestMoveMethodDefinedInPrototypeLiteralUsingShorthandSyntaxWithoutStub() {
	c.noStubs = true
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"function Foo() {}\n"+
						"Foo.prototype = { method() {} };\n").
				AddChunk("(new Foo).method()").
				Build()),
		c.Expected(
			"function Foo() {}\n"+
				"Foo.prototype = {};\n",
			// Chunk 2
			"Foo.prototype.method = function() {};\n"+
				"(new Foo).method()\n"))
}

// Java: CrossChunkMethodMotionTest.doNotMoveMethodDefinedInPrototypeLiteralContainingSuper()
func (c *CrossChunkMethodMotionTest) TestDoNotMoveMethodDefinedInPrototypeLiteralContainingSuper() {
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"function Foo() {}\n" +
						"Foo.prototype = { method() { return super.toString(); } };\n").
				AddChunk("(new Foo).method()").
				Build()))
}

// Java: CrossChunkMethodMotionTest.doNotMoveMethodDefinedInPrototypeLiteralAsComputedProp()
func (c *CrossChunkMethodMotionTest) TestDoNotMoveMethodDefinedInPrototypeLiteralAsComputedProp() {
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"function Foo() {}\n" +
						"Foo.prototype = { [1]:  {} };\n").
				AddChunk("(new Foo)[1]()").
				Build()))
}

// Java: CrossChunkMethodMotionTest.moveClassMethod()
func (c *CrossChunkMethodMotionTest) TestMoveClassMethod() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("class Foo { method() {} }").
				AddChunk("(new Foo).method()").
				Build()),
		c.Expected(
			CrossChunkMethodMotionTest_STUB_DECLARATIONS+
				"class Foo {}\n"+
				"Foo.prototype.method = JSCompiler_stubMethod(0);\n",
			// Chunk 2
			"Foo.prototype.method = JSCompiler_unstubMethod(0, function() {});\n"+
				"(new Foo).method();\n"))

	// Same as above, but reference to the method is via an optional chain
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("class Foo { method() {} }").
				AddChunk("(new Foo)?.method()").
				Build()),
		c.Expected(
			CrossChunkMethodMotionTest_STUB_DECLARATIONS+
				"class Foo {}\n"+
				"Foo.prototype.method = JSCompiler_stubMethod(0);\n",
			// Chunk 2
			"Foo.prototype.method = JSCompiler_unstubMethod(0, function() {});\n"+
				"(new Foo)?.method();\n"))
}

// Java: CrossChunkMethodMotionTest.doNotMoveClassMethodContainingSuper()
func (c *CrossChunkMethodMotionTest) TestDoNotMoveClassMethodContainingSuper() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"class Bar { method() {} }\n"+
						"class Foo extends Bar { method2() { super.method(); } }\n").
				AddChunk("(new Foo).method2()").
				Build()),
		c.Expected(
			CrossChunkMethodMotionTest_STUB_DECLARATIONS+
				"class Bar {}\n"+
				"Bar.prototype.method = JSCompiler_stubMethod(0);\n"+
				"class Foo extends Bar { method2() { super.method(); } }\n"+
				"\n",
			"Bar.prototype.method = JSCompiler_unstubMethod(0, function() {});\n"+
				"(new Foo).method2();\n"))
}

// Java: CrossChunkMethodMotionTest.doNotMoveClassMethodContainingSuperInAnArrow()
func (c *CrossChunkMethodMotionTest) TestDoNotMoveClassMethodContainingSuperInAnArrow() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"class Bar { method() {} }\n"+
						"class Foo extends Bar { method2() { return () => super.method(); } }\n").
				AddChunk("(new Foo).method2()").
				Build()),
		c.Expected(
			CrossChunkMethodMotionTest_STUB_DECLARATIONS+
				"class Bar {}\n"+
				"Bar.prototype.method = JSCompiler_stubMethod(0);\n"+
				"class Foo extends Bar { method2() { return () => super.method(); } }\n"+
				"\n",
			"Bar.prototype.method = JSCompiler_unstubMethod(0, function() {});\n"+
				"(new Foo).method2();\n"))
}

// Java: CrossChunkMethodMotionTest.moveClassMethodContainingObjLitContainingSuper()
func (c *CrossChunkMethodMotionTest) TestMoveClassMethodContainingObjLitContainingSuper() {
	// Don't be fooled by `super` that isn't referring to the method's `super`.
	// This `super` isn't really a reference within `method()`
	// It refers to Object.prototype.toString.
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"class Foo {\n"+
						"  method() {\n"+
						"    return {\n"+
						"      objLitMethod() {\n"+
						"// This `super` isn't really a reference within `method()`\n"+
						"// It refers to Object.prototype.toString.\n"+
						"        super.toString;\n"+
						"      }\n"+
						"    };\n"+
						"  }\n"+
						"}\n").
				AddChunk("(new Foo).method();").
				Build()),
		c.Expected(
			CrossChunkMethodMotionTest_STUB_DECLARATIONS+
				"class Foo {\n"+
				"}\n"+
				"Foo.prototype.method = JSCompiler_stubMethod(0);\n",
			"Foo.prototype.method = JSCompiler_unstubMethod(0, function() {\n"+
				"  return {\n"+
				"    objLitMethod() {\n"+
				"      super.toString;\n"+
				"    }\n"+
				"  };\n"+
				"});\n"+
				"(new Foo).method();\n"))
}

// Java: CrossChunkMethodMotionTest.doNotMoveClassMethodContainingSuperDefaultParam()
func (c *CrossChunkMethodMotionTest) TestDoNotMoveClassMethodContainingSuperDefaultParam() {
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"class Bar { defaultValue() { return 1; } }\n" +
						"class Foo extends Bar { method(x = super.defaultValue()) { return x; } }\n").
				AddChunk("(new Foo).method2()").
				Build()))
}

// Java: CrossChunkMethodMotionTest.doNotMoveClassConstructor()
func (c *CrossChunkMethodMotionTest) TestDoNotMoveClassConstructor() {
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"class Foo { constructor() { } }\n").
				AddChunk("(new Foo).constructor").
				Build()))
}

// Java: CrossChunkMethodMotionTest.doNotMoveClassComputedPropertyMethod()
func (c *CrossChunkMethodMotionTest) TestDoNotMoveClassComputedPropertyMethod() {
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("const methodName = 'method';").
				AddChunk("class Foo { [methodName]() {} }").
				AddChunk("(new Foo)[methodName]()").
				Build()))
}

// Java: CrossChunkMethodMotionTest.moveClassMethodForConstDefinition()
func (c *CrossChunkMethodMotionTest) TestMoveClassMethodForConstDefinition() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("const Foo = class FooInternal { method() {} }").
				AddChunk("(new Foo).method()").
				Build()),
		c.Expected(
			CrossChunkMethodMotionTest_STUB_DECLARATIONS+
				"const Foo = class FooInternal {}\n"+
				"Foo.prototype.method = JSCompiler_stubMethod(0);\n",
			// Chunk 2
			"Foo.prototype.method = JSCompiler_unstubMethod(0, function() {});\n"+
				"(new Foo).method();\n"))
}

// Java: CrossChunkMethodMotionTest.doNotMoveFunctionCall_thatIsSideEffected()
func (c *CrossChunkMethodMotionTest) TestDoNotMoveFunctionCall_thatIsSideEffected() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	chunks :=
		JSChunkGraphBuilder_ForChain().
			// m1
			AddChunk(
				"var a = 0;\n" +
					"function f1(a) { return a + 1 }\n" +
					"var b = f1(1);\n" +
					"a += 1;\n").
			// m2
			AddChunk("var c = b").
			Build()

	c.TestSame(c.Srcs(chunks))
}

// Java: CrossChunkMethodMotionTest.doNotMoveClassMethodWithLocalClassNameReference()
func (c *CrossChunkMethodMotionTest) TestDoNotMoveClassMethodWithLocalClassNameReference() {
	// We could probably rewrite the internal reference, but it is unlikely that the added
	// complexity of doing so would be worthwhile.
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("const Foo = class FooInternal { method() { FooInternal; } }").
				AddChunk("(new Foo).method()").
				Build()))
}

// Java: CrossChunkMethodMotionTest.doNotMoveGetterDefinedInPrototypeLiteral()
func (c *CrossChunkMethodMotionTest) TestDoNotMoveGetterDefinedInPrototypeLiteral() {
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"function Foo() {}\n" +
						"Foo.prototype = { get method() {} };\n").
				AddChunk("(new Foo).method()").
				Build()))
}

// Java: CrossChunkMethodMotionTest.doNotMoveClassGetter()
func (c *CrossChunkMethodMotionTest) TestDoNotMoveClassGetter() {
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("class Foo { get method() {} }").
				AddChunk("(new Foo).method()").
				Build()))
}

// Java: CrossChunkMethodMotionTest.movePrototypeMethodWithoutStub()
func (c *CrossChunkMethodMotionTest) TestMovePrototypeMethodWithoutStub() {
	c.TestSame(
		c.Externs("IFoo.prototype.bar;"),
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"function Foo() {}\n"+
						"Foo.prototype.bar = function() {};\n").
				AddChunk("(new Foo).bar()").
				Build()))

	c.canMoveExterns = true
	c.noStubs = true
	c.Test(
		c.Externs("IFoo.prototype.bar;"),
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"function Foo() {}\n"+
						"Foo.prototype.bar = function() {};\n").
				AddChunk("(new Foo).bar()").
				Build()),
		c.Expected(
			"function Foo() {}",
			// Chunk 2
			"Foo.prototype.bar = function() {};\n"+
				"(new Foo).bar()\n"))
}

// Java: CrossChunkMethodMotionTest.movePrototypeMethodImplementingInterfaceWithoutStub()
func (c *CrossChunkMethodMotionTest) TestMovePrototypeMethodImplementingInterfaceWithoutStub() {
	c.DisableCompareJsDoc() // multistage compilation erases the @implements
	c.TestSame(
		c.Externs(
			"/** @interface */\n"+
				"class IFoo {\n"+
				"  ifooMethod() {}\n"+
				"}\n"),
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"/**\n"+
						" * @constructor\n"+
						" * @implements {IFoo}\n"+
						" */\n"+
						"function Foo() {}\n"+
						"Foo.prototype.ifooMethod = function() {};\n").
				AddChunk("(new Foo).ifooMethod()").
				Build()))

	c.canMoveExterns = true
	c.noStubs = true
	c.Test(
		c.Externs(
			"/** @interface */\n"+
				"class IFoo {\n"+
				"  ifooMethod() {}\n"+
				"}\n"),
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"/**\n"+
						" * @constructor\n"+
						" * @implements {IFoo}\n"+
						" */\n"+
						"function Foo() {}\n"+
						"Foo.prototype.ifooMethod = function() {};\n").
				AddChunk("(new Foo).ifooMethod()").
				Build()),
		c.Expected(
			"/**\n"+
				" * @constructor\n"+
				" */\n"+
				"function Foo() {}\n",
			// Chunk 2
			"Foo.prototype.ifooMethod = function() {};\n"+
				"(new Foo).ifooMethod()\n"))
}

// Java: CrossChunkMethodMotionTest.moveClassMethodWithoutStub()
func (c *CrossChunkMethodMotionTest) TestMoveClassMethodWithoutStub() {
	c.TestSame(
		c.Externs("IFoo.prototype.bar;"),
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("class Foo { bar() {} }").
				AddChunk("(new Foo).bar()").
				Build()))

	c.canMoveExterns = true
	c.noStubs = true
	c.Test(
		c.Externs("IFoo.prototype.bar;"),
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("class Foo { bar() {} }").
				AddChunk("(new Foo).bar()").
				Build()),
		c.Expected(
			"class Foo {}",
			// Chunk 2
			"Foo.prototype.bar = function() {};\n"+
				"(new Foo).bar()\n"))
}

// Java: CrossChunkMethodMotionTest.moveClassMethodImplementingExternsInterfaceWithoutStub()
func (c *CrossChunkMethodMotionTest) TestMoveClassMethodImplementingExternsInterfaceWithoutStub() {
	c.DisableCompareJsDoc() // multistage compilation deletes the @implements
	c.TestSame(
		c.Externs(
			"/** @interface */\n"+
				"class IFoo {\n"+
				"  ifooMethod() {}\n"+
				"}\n"),
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"/** @implements {IFoo} */\n"+
						"class Foo { ifooMethod() {} }\n").
				AddChunk("(new Foo).ifooMethod()").
				Build()))

	c.canMoveExterns = true
	c.noStubs = true
	c.Test(
		c.Externs(
			"/** @interface */\n"+
				"class IFoo {\n"+
				"  ifooMethod() {}\n"+
				"}\n"),
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"/** @implements {IFoo} */\n"+
						"class Foo { ifooMethod() {} }\n").
				AddChunk("(new Foo).ifooMethod()").
				Build()),
		c.Expected(
			"class Foo {}",
			// Chunk 2
			"Foo.prototype.ifooMethod = function() {};\n"+
				"(new Foo).ifooMethod()\n"))
}

// Java: CrossChunkMethodMotionTest.doNotMovePrototypeMethodIfAliasedAndNoStubs()
func (c *CrossChunkMethodMotionTest) TestDoNotMovePrototypeMethodIfAliasedAndNoStubs() {
	// don't move if noStubs enabled and there's a reference to the method to be moved
	c.noStubs = true
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"function Foo() {}\n" +
						"Foo.prototype.m = function() {};\n" +
						"Foo.prototype.m2 = Foo.prototype.m;\n").
				AddChunk("(new Foo).m()").
				Build()))

	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"function Foo() {}\n" +
						"Foo.prototype.m = function() {};\n" +
						"Foo.prototype.m2 = Foo.prototype.m;\n").
				AddChunk("(new Foo).m(), (new Foo).m2()").
				Build()))

	c.noStubs = false

	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"function Foo() {}\n"+
						"Foo.prototype.m = function() {};\n"+
						"Foo.prototype.m2 = Foo.prototype.m;\n").
				AddChunk("(new Foo).m()").
				Build()),
		c.Expected(
			CrossChunkMethodMotionTest_STUB_DECLARATIONS+
				"function Foo() {}\n"+
				"Foo.prototype.m = JSCompiler_stubMethod(0);\n"+
				"Foo.prototype.m2 = Foo.prototype.m;\n",
			// Chunk 2
			"Foo.prototype.m = JSCompiler_unstubMethod(0, function() {});\n"+
				"(new Foo).m()\n"))

	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"function Foo() {}\n"+
						"Foo.prototype.m = function() {};\n"+
						"Foo.prototype.m2 = Foo.prototype.m;\n").
				AddChunk("(new Foo).m(), (new Foo).m2()").
				Build()),
		c.Expected(
			CrossChunkMethodMotionTest_STUB_DECLARATIONS+
				"function Foo() {}\n"+
				"Foo.prototype.m = JSCompiler_stubMethod(0);\n"+
				"Foo.prototype.m2 = Foo.prototype.m;\n",
			// Chunk 2
			"Foo.prototype.m = JSCompiler_unstubMethod(0, function() {});\n"+
				"(new Foo).m(), (new Foo).m2()\n"))
}

// Java: CrossChunkMethodMotionTest.doNotMoveClassMethodIfAliasedAndNoStubs()
func (c *CrossChunkMethodMotionTest) TestDoNotMoveClassMethodIfAliasedAndNoStubs() {
	// don't move if noStubs enabled and there's a reference to the method to be moved
	c.noStubs = true
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"class Foo { m() {} }\n" +
						"Foo.prototype.m2 = Foo.prototype.m;\n").
				AddChunk("(new Foo).m()").
				Build()))

	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"class Foo { m() {} }\n" +
						"Foo.prototype.m2 = Foo.prototype.m;\n").
				AddChunk("(new Foo).m(), (new Foo).m2()").
				Build()))

	c.noStubs = false

	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"class Foo { m() {} }\n"+
						"Foo.prototype.m2 = Foo.prototype.m;\n").
				AddChunk("(new Foo).m()").
				Build()),
		c.Expected(
			CrossChunkMethodMotionTest_STUB_DECLARATIONS+
				"class Foo {}\n"+
				"Foo.prototype.m = JSCompiler_stubMethod(0);\n"+
				"Foo.prototype.m2 = Foo.prototype.m;\n",
			// Chunk 2
			"Foo.prototype.m = JSCompiler_unstubMethod(0, function() {});\n"+
				"(new Foo).m()\n"))

	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"class Foo { m() {} }\n"+
						"Foo.prototype.m2 = Foo.prototype.m;\n").
				AddChunk("(new Foo).m(), (new Foo).m2()").
				Build()),
		c.Expected(
			CrossChunkMethodMotionTest_STUB_DECLARATIONS+
				"class Foo {}\n"+
				"Foo.prototype.m = JSCompiler_stubMethod(0);\n"+
				"Foo.prototype.m2 = Foo.prototype.m;\n",
			// Chunk 2
			"Foo.prototype.m = JSCompiler_unstubMethod(0, function() {});\n"+
				"(new Foo).m(), (new Foo).m2()\n"))
}

// Java: CrossChunkMethodMotionTest.doNotMovePrototypeMethodRedeclaredInSiblingChunk()
func (c *CrossChunkMethodMotionTest) TestDoNotMovePrototypeMethodRedeclaredInSiblingChunk() {
	// don't move if it can be overwritten when a sibling of the first referencing chunk is loaded.
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk(
					"function Foo() {}\n" +
						"Foo.prototype.method = function() {};\n").
				AddChunk("Foo.prototype.method = function() {};").
				AddChunk("(new Foo).method()").
				Build()))
}

// Java: CrossChunkMethodMotionTest.doNotMoveClassMethodRedeclaredInSiblingChunk()
func (c *CrossChunkMethodMotionTest) TestDoNotMoveClassMethodRedeclaredInSiblingChunk() {
	// don't move if it can be overwritten when a sibling of the first referencing chunk is loaded.
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk("class Foo { method() {} }").
				AddChunk("Foo.prototype.method = function() {};").
				AddChunk("(new Foo).method()").
				Build()))
}

// Java: CrossChunkMethodMotionTest.doNotMovePrototypeMethodRedeclaredInDependentChunk()
func (c *CrossChunkMethodMotionTest) TestDoNotMovePrototypeMethodRedeclaredInDependentChunk() {
	// don't move if it can be overwritten by a chunk depending on the first referencing chunk.
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"function Foo() {}\n" +
						"Foo.prototype.method = function() {};\n").
				AddChunk("(new Foo).method()").
				AddChunk("Foo.prototype.method = function() {};").
				Build()))
}

// Java: CrossChunkMethodMotionTest.doNotMoveClassMethodRedeclaredInDependentChunk()
func (c *CrossChunkMethodMotionTest) TestDoNotMoveClassMethodRedeclaredInDependentChunk() {
	// don't move if it can be overwritten by a chunk depending on the first referencing chunk.
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("class Foo { method() {} }").
				AddChunk("(new Foo).method()").
				AddChunk("Foo.prototype.method = function() {};").
				Build()))
}

// Java: CrossChunkMethodMotionTest.doNotMovePrototypeMethodRedeclaredBeforeFirstReferencingChunk()
func (c *CrossChunkMethodMotionTest) TestDoNotMovePrototypeMethodRedeclaredBeforeFirstReferencingChunk() {
	// Note: it is reasonable to move the method in this case,
	// but it is difficult enough to prove that we don't.
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"function Foo() {}\n" +
						"Foo.prototype.method = function() {};\n").
				AddChunk("Foo.prototype.method = function() {};").
				AddChunk("(new Foo).method()").
				Build()))
}

// Java: CrossChunkMethodMotionTest.doNotMoveClassMethodRedeclaredBeforeFirstReferencingChunk()
func (c *CrossChunkMethodMotionTest) TestDoNotMoveClassMethodRedeclaredBeforeFirstReferencingChunk() {
	// Note: it is reasonable to move the method in this case,
	// but it is difficult enough to prove that we don't.
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("class Foo { method() {} }").
				AddChunk("Foo.prototype.method = function() {};").
				AddChunk("(new Foo).method()").
				Build()))
}

// Java: CrossChunkMethodMotionTest.movePrototypeRecursiveMethod()
func (c *CrossChunkMethodMotionTest) TestMovePrototypeRecursiveMethod() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"function Foo() {}\n"+
						"Foo.prototype.baz = function() { this.baz(); };\n").
				AddChunk("(new Foo).baz()").
				Build()),
		c.Expected(
			CrossChunkMethodMotionTest_STUB_DECLARATIONS+
				"function Foo() {}\n"+
				"Foo.prototype.baz = JSCompiler_stubMethod(0);\n",
			// Chunk 2
			"Foo.prototype.baz = JSCompiler_unstubMethod(0, function() { this.baz(); });\n"+
				"(new Foo).baz()\n"))
}

// Java: CrossChunkMethodMotionTest.moveInstanceRecursiveMethod()
func (c *CrossChunkMethodMotionTest) TestMoveInstanceRecursiveMethod() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("class Foo { baz() { this.baz(); } }").
				AddChunk("(new Foo).baz()").
				Build()),
		c.Expected(
			CrossChunkMethodMotionTest_STUB_DECLARATIONS+
				"class Foo {}\n"+
				"Foo.prototype.baz = JSCompiler_stubMethod(0);\n",
			// Chunk 2
			"Foo.prototype.baz = JSCompiler_unstubMethod(0, function() { this.baz(); });\n"+
				"(new Foo).baz()\n"))
}

// Java: CrossChunkMethodMotionTest.doNotMoveNonLiteralFunction()
func (c *CrossChunkMethodMotionTest) TestDoNotMoveNonLiteralFunction() {
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"function Foo() {}\n" +
						"Foo.prototype.baz = shared;\n").
				AddChunk("(new Foo).baz()").
				Build()))

	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"class Foo {}\n" +
						"Foo.prototype.baz = shared;\n").
				AddChunk("(new Foo).baz()").
				Build()))
}

// Java: CrossChunkMethodMotionTest.movePrototypeDeclarationsInTheRightOrder()
func (c *CrossChunkMethodMotionTest) TestMovePrototypeDeclarationsInTheRightOrder() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"function Foo() {}\n"+
						"Foo.prototype.baz = function() { return 1; };\n"+
						"Foo.prototype.baz = function() { return 2; };\n").
				AddChunk("(new Foo).baz()").
				Build()),
		c.Expected(
			CrossChunkMethodMotionTest_STUB_DECLARATIONS+
				"function Foo() {}\n"+
				"Foo.prototype.baz = JSCompiler_stubMethod(1);\n"+
				"Foo.prototype.baz = JSCompiler_stubMethod(0);\n",
			// Chunk 2
			"Foo.prototype.baz = JSCompiler_unstubMethod(1, function() { return 1; });\n"+
				"Foo.prototype.baz = JSCompiler_unstubMethod(0, function() { return 2; });\n"+
				"(new Foo).baz()\n"))
}

// Java: CrossChunkMethodMotionTest.moveClassMethodAndReclarationInTheRightOrder()
func (c *CrossChunkMethodMotionTest) TestMoveClassMethodAndReclarationInTheRightOrder() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"class Foo { baz() { return 1; } }\n"+
						"Foo.prototype.baz = function() { return 2; };\n").
				AddChunk("(new Foo).baz()").
				Build()),
		c.Expected(
			CrossChunkMethodMotionTest_STUB_DECLARATIONS+
				"class Foo {}\n"+
				"Foo.prototype.baz = JSCompiler_stubMethod(1);\n"+
				"Foo.prototype.baz = JSCompiler_stubMethod(0);\n",
			// Chunk 2
			"Foo.prototype.baz =\n"+
				"JSCompiler_unstubMethod(1, function() { return 1; });\n"+
				"Foo.prototype.baz =\n"+
				"JSCompiler_unstubMethod(0, function() { return 2; });\n"+
				"(new Foo).baz()\n"))
}

// Java: CrossChunkMethodMotionTest.movePrototypeMethodsForDifferentClassesInTheRightOrder()
func (c *CrossChunkMethodMotionTest) TestMovePrototypeMethodsForDifferentClassesInTheRightOrder() {
	m :=
		JSChunkGraphBuilder_ForUnordered().
			AddChunk(
				"function Foo() {}\n" +
					"Foo.prototype.baz = function() { return 1; };\n" +
					"function Goo() {}\n" +
					"Goo.prototype.baz = function() { return 2; };\n").

			// Chunk 2, depends on 1
			AddChunk("").
			// Chunk 3, depends on 2
			AddChunk("(new Foo).baz()").
			// Chunk 4, depends on 3
			AddChunk("").
			// Chunk 5, depends on 3
			AddChunk("(new Goo).baz()").
			Build()

	m[1].AddDependency(m[0])
	m[2].AddDependency(m[1])
	m[3].AddDependency(m[2])
	m[4].AddDependency(m[2])

	c.Test(
		c.Srcs(m),
		c.Expected(
			CrossChunkMethodMotionTest_STUB_DECLARATIONS+
				"function Foo() {}\n"+
				"Foo.prototype.baz = JSCompiler_stubMethod(1);\n"+
				"function Goo() {}\n"+
				"Goo.prototype.baz = JSCompiler_stubMethod(0);\n",
			// Chunk 2
			"",
			// Chunk 3
			"Foo.prototype.baz = JSCompiler_unstubMethod(1, function() { return 1; });\n"+
				"Goo.prototype.baz = JSCompiler_unstubMethod(0, function() { return 2; });\n"+
				"(new Foo).baz()\n",
			// Chunk 4
			"",
			// Chunk 5
			"(new Goo).baz()"))
}

// Java: CrossChunkMethodMotionTest.moveClassMethodsForDifferentClassesInTheRightOrder()
func (c *CrossChunkMethodMotionTest) TestMoveClassMethodsForDifferentClassesInTheRightOrder() {
	m :=
		JSChunkGraphBuilder_ForUnordered().
			AddChunk(
				"class Foo { baz() { return 1; } }\n" +
					"class Goo { baz() { return 2; } }\n").
			// Chunk 2, depends on 1
			AddChunk("").
			// Chunk 3, depends on 2
			AddChunk("(new Foo).baz()").
			// Chunk 4, depends on 3
			AddChunk("").
			// Chunk 5, depends on 3
			AddChunk("(new Goo).baz()").
			Build()

	m[1].AddDependency(m[0])
	m[2].AddDependency(m[1])
	m[3].AddDependency(m[2])
	m[4].AddDependency(m[2])

	c.Test(
		c.Srcs(m),
		c.Expected(
			CrossChunkMethodMotionTest_STUB_DECLARATIONS+
				"class Foo {}\n"+
				"Foo.prototype.baz = JSCompiler_stubMethod(1);\n"+
				"class Goo {}\n"+
				"Goo.prototype.baz = JSCompiler_stubMethod(0);\n",
			// Chunk 2
			"",
			// Chunk 3
			"Foo.prototype.baz = JSCompiler_unstubMethod(1, function() { return 1; });\n"+
				"Goo.prototype.baz = JSCompiler_unstubMethod(0, function() { return 2; });\n"+
				"(new Foo).baz()\n",
			// Chunk 4
			"",
			// Chunk 5
			"(new Goo).baz()"))
}

// Java: CrossChunkMethodMotionTest.doNotMovePrototypeMethodUsedInMultiplepDependentChunks()
func (c *CrossChunkMethodMotionTest) TestDoNotMovePrototypeMethodUsedInMultiplepDependentChunks() {
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk(
					"function Foo() {}\n" +
						"Foo.prototype.baz = function() {};\n").
				AddChunk("(new Foo).baz()").
				AddChunk("(new Foo).baz()").
				Build()))
}

// Java: CrossChunkMethodMotionTest.doNotMoveClassMethodUsedInMultiplepDependentChunks()
func (c *CrossChunkMethodMotionTest) TestDoNotMoveClassMethodUsedInMultiplepDependentChunks() {
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk("class Foo { baz() {} }").
				AddChunk("(new Foo).baz()").
				AddChunk("(new Foo).baz()").
				Build()))
}

// Java: CrossChunkMethodMotionTest.movePrototypeMethodToDeepestCommonDependencyOfReferencingChunks()
func (c *CrossChunkMethodMotionTest) TestMovePrototypeMethodToDeepestCommonDependencyOfReferencingChunks() {
	chunks :=
		JSChunkGraphBuilder_ForUnordered().
			AddChunk(
				"function Foo() {}\n" +
					"Foo.prototype.baz = function() {};\n").
			// Chunk 2
			// a blank chunk in the middle
			AddChunk("").
			// Chunk 3
			AddChunk("(new Foo).baz() , 1").
			// Chunk 4
			AddChunk("(new Foo).baz() , 2").
			Build()

	chunks[1].AddDependency(chunks[0])
	chunks[2].AddDependency(chunks[1])
	chunks[3].AddDependency(chunks[1])
	c.Test(
		c.Srcs(chunks),
		c.Expected(
			CrossChunkMethodMotionTest_STUB_DECLARATIONS+
				"function Foo() {}\n"+
				"Foo.prototype.baz = JSCompiler_stubMethod(0);\n",
			// Chunk 2
			"Foo.prototype.baz = JSCompiler_unstubMethod(0, function() {});",
			// Chunk 3
			"(new Foo).baz() , 1",
			// Chunk 4
			"(new Foo).baz() , 2"))
}

// Java: CrossChunkMethodMotionTest.moveClassMethodToDeepestCommonDependencyOfReferencingChunks()
func (c *CrossChunkMethodMotionTest) TestMoveClassMethodToDeepestCommonDependencyOfReferencingChunks() {
	chunks :=
		JSChunkGraphBuilder_ForUnordered().
			AddChunk("class Foo { baz() {} }").
			// Chunk 2
			// a blank chunk in the middle
			AddChunk("").
			// Chunk 3
			AddChunk("(new Foo).baz() , 1").
			// Chunk 4
			AddChunk("(new Foo).baz() , 2").
			Build()

	chunks[1].AddDependency(chunks[0])
	chunks[2].AddDependency(chunks[1])
	chunks[3].AddDependency(chunks[1])

	c.Test(
		c.Srcs(chunks),
		c.Expected(
			CrossChunkMethodMotionTest_STUB_DECLARATIONS+
				"class Foo {}\n"+
				"Foo.prototype.baz = JSCompiler_stubMethod(0);\n",
			// Chunk 2
			"Foo.prototype.baz = JSCompiler_unstubMethod(0, function() {});",
			// Chunk 3
			"(new Foo).baz() , 1",
			// Chunk 4
			"(new Foo).baz() , 2"))
}

// Java: CrossChunkMethodMotionTest.movePrototypeMethodThatRefersToAnotherOnTheSameClass()
func (c *CrossChunkMethodMotionTest) TestMovePrototypeMethodThatRefersToAnotherOnTheSameClass() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"function Foo() {}\n"+
						"Foo.prototype.baz = function() {};\n").
				// Chunk 2
				AddChunk("Foo.prototype.callBaz = function() { this.baz(); }").
				// Chunk 3
				AddChunk("(new Foo).callBaz()").
				Build()),
		c.Expected(
			CrossChunkMethodMotionTest_STUB_DECLARATIONS+
				"function Foo() {}\n"+
				"Foo.prototype.baz = JSCompiler_stubMethod(0);\n",
			// Chunk 2
			"Foo.prototype.callBaz = JSCompiler_stubMethod(1);",
			// Chunk 3
			"Foo.prototype.callBaz =\n"+
				"  JSCompiler_unstubMethod(1, function() { this.baz(); });\n"+
				"Foo.prototype.baz = JSCompiler_unstubMethod(0, function() {});\n"+
				"(new Foo).callBaz()\n"))
}

// Java: CrossChunkMethodMotionTest.movePrototypeMethodThatRefersToAnClassMethodOnTheSameClass()
func (c *CrossChunkMethodMotionTest) TestMovePrototypeMethodThatRefersToAnClassMethodOnTheSameClass() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("class Foo { baz() {} }").
				// Chunk 2
				AddChunk("Foo.prototype.callBaz = function() { this.baz(); }").
				// Chunk 3
				AddChunk("(new Foo).callBaz()").
				Build()),
		c.Expected(
			CrossChunkMethodMotionTest_STUB_DECLARATIONS+
				"class Foo {}\n"+
				"Foo.prototype.baz = JSCompiler_stubMethod(0);\n",
			// Chunk 2
			"Foo.prototype.callBaz = JSCompiler_stubMethod(1);",
			// Chunk 3
			"Foo.prototype.callBaz = JSCompiler_unstubMethod(1, function() { this.baz(); });\n"+
				"Foo.prototype.baz = JSCompiler_unstubMethod(0, function() {});\n"+
				"(new Foo).callBaz()\n"))
}

// Java: CrossChunkMethodMotionTest.doNotMovePrototypeMethodDefinitionThatFollowsFirstUse()
func (c *CrossChunkMethodMotionTest) TestDoNotMovePrototypeMethodDefinitionThatFollowsFirstUse() {
	// if the programmer screws up the chunk order, we don't try to correct
	// the mistake.
	// call before definition
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"function Foo() {}\n"+
						"Foo.prototype.baz = function() {};\n").
				// Chunk 2
				// call before definition
				AddChunk("(new Foo).callBaz()").
				// Chunk 3
				AddChunk("Foo.prototype.callBaz = function() { this.baz(); }").
				Build()),
		c.Expected(
			CrossChunkMethodMotionTest_STUB_DECLARATIONS+
				"function Foo() {}\n"+
				"Foo.prototype.baz = JSCompiler_stubMethod(0);\n",
			// Chunk 2
			"(new Foo).callBaz()",
			// Chunk 3
			"Foo.prototype.baz = JSCompiler_unstubMethod(0, function() {});\n"+
				"Foo.prototype.callBaz = function() { this.baz(); };\n"))
}

// Java: CrossChunkMethodMotionTest.movePrototypeMethodPastUsageInAGlobalFunction()
func (c *CrossChunkMethodMotionTest) TestMovePrototypeMethodPastUsageInAGlobalFunction() {
	// usage here doesn't really happen until x() is called, so
	// it's OK to move the definition of baz().
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"function Foo() {}\n"+
						"Foo.prototype.baz = function() {};\n"+
						"// usage here doesn't really happen until x() is called, so\n"+
						"// it's OK to move the definition of baz().\n"+
						"function x() { return (new Foo).baz(); }\n").
				// Chunk 2
				AddChunk("x();").
				Build()),
		c.Expected(
			CrossChunkMethodMotionTest_STUB_DECLARATIONS+
				"function Foo() {}\n"+
				"Foo.prototype.baz = JSCompiler_stubMethod(0);\n"+
				"function x() { return (new Foo).baz(); }\n",
			// Chunk 2
			"Foo.prototype.baz = JSCompiler_unstubMethod(0, function() {});\n"+
				"x();\n"))
}

// Java: CrossChunkMethodMotionTest.moveClassMethodPastUsageInAGlobalFunction()
func (c *CrossChunkMethodMotionTest) TestMoveClassMethodPastUsageInAGlobalFunction() {
	// usage here doesn't really happen until x() is called, so
	// it's OK to move the definition of baz().
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"class Foo {\n"+
						"  baz() {}\n"+
						"}\n"+
						"// usage here doesn't really happen until x() is called, so\n"+
						"// it's OK to move the definition of baz().\n"+
						"function x() { return (new Foo).baz(); }\n").
				// Chunk 2
				AddChunk("x();").
				Build()),
		c.Expected(
			CrossChunkMethodMotionTest_STUB_DECLARATIONS+
				"class Foo {}\n"+
				"Foo.prototype.baz = JSCompiler_stubMethod(0);\n"+
				"function x() { return (new Foo).baz(); }\n",
			// Chunk 2
			"Foo.prototype.baz = JSCompiler_unstubMethod(0, function() {});\n"+
				"x();\n"))
}

// Read of closure variable disables method motions.
// Java: CrossChunkMethodMotionTest.doNotMovePrototypeMethodThatUsesLocalClosureVariable()
func (c *CrossChunkMethodMotionTest) TestDoNotMovePrototypeMethodThatUsesLocalClosureVariable() {
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"function Foo() {}\n" +
						"(function() {\n" +
						"  var x = 'x';\n" +
						"  Foo.prototype.baz = function() {x};\n" +
						"})();\n").
				AddChunk("var y = new Foo(); y.baz();").
				Build()))
}

// Java: CrossChunkMethodMotionTest.doNotMoveClassMethodThatUsesLocalClosureVariable()
func (c *CrossChunkMethodMotionTest) TestDoNotMoveClassMethodThatUsesLocalClosureVariable() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"const Foo = (function() {\n" +
						"  var x = 'x';\n" +
						"  return class Foo { baz() { return x; } };\n" +
						"})();\n").
				AddChunk("var y = new Foo(); y.baz();").
				Build()))
}

// Java: CrossChunkMethodMotionTest.movePrototypeMethodThatDefinesOtherMethodsOnSameGlobalClass()
func (c *CrossChunkMethodMotionTest) TestMovePrototypeMethodThatDefinesOtherMethodsOnSameGlobalClass() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"function Foo() {}\n"+
						"Foo.prototype.b1 = function() {\n"+
						"  var x = 1;\n"+
						"  Foo.prototype.b2 = function() {\n"+
						"    Foo.prototype.b3 = function() {\n"+
						"      x;\n"+
						"    }\n"+
						"  }\n"+
						"};\n").
				// Chunk 2
				AddChunk("var y = new Foo(); y.b1();").
				// Chunk 3
				AddChunk("y = new Foo(); z.b2();").
				// Chunk 4
				AddChunk("y = new Foo(); z.b3();").
				Build()),
		c.Expected(
			CrossChunkMethodMotionTest_STUB_DECLARATIONS+
				"function Foo() {}\n"+
				"Foo.prototype.b1 = JSCompiler_stubMethod(0);\n",
			// Chunk 2
			"Foo.prototype.b1 = JSCompiler_unstubMethod(0, function() {\n"+
				"  var x = 1;\n"+
				"  Foo.prototype.b2 = function() {\n"+
				"    Foo.prototype.b3 = function() {\n"+
				"      x;\n"+
				"    }\n"+
				"  }\n"+
				"});\n"+
				"var y = new Foo(); y.b1();\n",
			// Chunk 3
			"y = new Foo(); z.b2();",
			// Chunk 4
			"y = new Foo(); z.b3();"))
}

// Java: CrossChunkMethodMotionTest.moveClassMethodThatDefinesOtherMethodsOnSameGlobalClass()
func (c *CrossChunkMethodMotionTest) TestMoveClassMethodThatDefinesOtherMethodsOnSameGlobalClass() {
	// b2 cannot be extracted, because it contains a reference to x
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"class Foo {\n"+
						"  b1() {\n"+
						"    var x = 1;\n"+
						"// b2 cannot be extracted, because it contains a reference to x\n"+
						"    Foo.prototype.b2 = function() {\n"+
						"      Foo.prototype.b3 = function() {\n"+
						"        x;\n"+
						"      }\n"+
						"    }\n"+
						"  };\n"+
						"}\n").
				AddChunk("var y = new Foo(); y.b1();").
				AddChunk("y = new Foo(); z.b2();").
				AddChunk("y = new Foo(); z.b3();").
				Build()),
		c.Expected(
			CrossChunkMethodMotionTest_STUB_DECLARATIONS+
				"class Foo {}\n"+
				"Foo.prototype.b1 = JSCompiler_stubMethod(0);\n",
			// Chunk 2
			"Foo.prototype.b1 = JSCompiler_unstubMethod(0, function() {\n"+
				"  var x = 1;\n"+
				"  Foo.prototype.b2 = function() {\n"+
				"    Foo.prototype.b3 = function() {\n"+
				"      x;\n"+
				"    }\n"+
				"  }\n"+
				"});\n"+
				"var y = new Foo(); y.b1();\n",
			// Chunk 3
			"y = new Foo(); z.b2();",
			// Chunk 4
			"y = new Foo(); z.b3();"))
}

// Java: CrossChunkMethodMotionTest.extractPrototypeMethodDefinedInAnotherMethodWhenNoClosureReferencePreventsIt()
func (c *CrossChunkMethodMotionTest) TestExtractPrototypeMethodDefinedInAnotherMethodWhenNoClosureReferencePreventsIt() {
	// definition of b2 can be extracted, because it doesn't refer to any
	// variables
	// defined by b1.
	// definition of b3 cannot be extracted, because it refers to x
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"function Foo() {}\n"+
						"Foo.prototype.b1 = function() {\n"+
						"// definition of b2 can be extracted, because it doesn't refer to any\n"+
						"// variables\n"+
						"// defined by b1.\n"+
						"  Foo.prototype.b2 = function() {\n"+
						"    var x = 1;\n"+
						"// definition of b3 cannot be extracted, because it refers to x\n"+
						"    Foo.prototype.b3 = function() {\n"+
						"      x;\n"+
						"    }\n"+
						"  }\n"+
						"};\n").
				AddChunk("var y = new Foo(); y.b1();").
				AddChunk("y = new Foo(); z.b2();").
				AddChunk("y = new Foo(); z.b3();").
				Build()),
		c.Expected(
			CrossChunkMethodMotionTest_STUB_DECLARATIONS+
				"function Foo() {}\n"+
				"Foo.prototype.b1 = JSCompiler_stubMethod(0);\n",
			// Chunk 2
			"Foo.prototype.b1 = JSCompiler_unstubMethod(0, function() {\n"+
				"  Foo.prototype.b2 = JSCompiler_stubMethod(1);\n"+
				"});\n"+
				"var y = new Foo(); y.b1();\n",
			// Chunk 3
			"Foo.prototype.b2 = JSCompiler_unstubMethod(1, function() {\n"+
				"  var x = 1;\n"+
				"  Foo.prototype.b3 = function() {\n"+
				"    x;\n"+
				"  }\n"+
				"});\n"+
				"y = new Foo(); z.b2();\n",
			// Chunk 4
			"y = new Foo(); z.b3();"))
}

// Java: CrossChunkMethodMotionTest.extractClassMethodDefinedInAnotherMethodWhenNoClosureReferencePreventsIt()
func (c *CrossChunkMethodMotionTest) TestExtractClassMethodDefinedInAnotherMethodWhenNoClosureReferencePreventsIt() {
	// definition of b2 can be extracted, because it doesn't refer to any
	// variables
	// defined by b1.
	// definition of b3 cannot be extracted, because it refers to x
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"class Foo {\n"+
						"  b1() {\n"+
						"// definition of b2 can be extracted, because it doesn't refer to any\n"+
						"// variables\n"+
						"// defined by b1.\n"+
						"    Foo.prototype.b2 = function() {\n"+
						"      var x = 1;\n"+
						"// definition of b3 cannot be extracted, because it refers to x\n"+
						"      Foo.prototype.b3 = function() {\n"+
						"        x;\n"+
						"      }\n"+
						"    }\n"+
						"  }\n"+
						"}\n").
				// Chunk 2
				AddChunk("var y = new Foo(); y.b1();").
				// Chunk 3
				AddChunk("y = new Foo(); z.b2();").
				// Chunk 4
				AddChunk("y = new Foo(); z.b3();").
				Build()),
		c.Expected(
			CrossChunkMethodMotionTest_STUB_DECLARATIONS+
				"class Foo {}\n"+
				"Foo.prototype.b1 = JSCompiler_stubMethod(0);\n",
			// Chunk 2
			"Foo.prototype.b1 =\n"+
				"    JSCompiler_unstubMethod(\n"+
				"        0,\n"+
				"        function() {\n"+
				"          Foo.prototype.b2 = JSCompiler_stubMethod(1);\n"+
				"        });\n"+
				"\n"+
				"\n"+
				"var y = new Foo(); y.b1();\n",
			// Chunk 3
			"Foo.prototype.b2 = JSCompiler_unstubMethod(1, function() {\n"+
				"  var x = 1;\n"+
				"  Foo.prototype.b3 = function() {\n"+
				"    x;\n"+
				"  }\n"+
				"});\n"+
				"y = new Foo(); z.b2();\n",
			// Chunk 4
			"y = new Foo(); z.b3();"))
}

// Read of global variable is fine.
// Java: CrossChunkMethodMotionTest.movePrototypeMethodThatReadsGlobalVar()
func (c *CrossChunkMethodMotionTest) TestMovePrototypeMethodThatReadsGlobalVar() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"function Foo() {}\n"+
						"var x = 'x';\n"+
						"Foo.prototype.baz = function(){x};\n").
				AddChunk("var y = new Foo(); y.baz();").
				Build()),
		c.Expected(
			CrossChunkMethodMotionTest_STUB_DECLARATIONS+
				"function Foo() {}\n"+
				"var x = 'x';\n"+
				"Foo.prototype.baz = JSCompiler_stubMethod(0);\n",
			// Chunk 2
			"Foo.prototype.baz = JSCompiler_unstubMethod(0, function(){x});\n"+
				"var y = new Foo(); y.baz();\n"))
}

// Read of global variable is fine.
// Java: CrossChunkMethodMotionTest.moveClassMethodThatReadsGlobalVar()
func (c *CrossChunkMethodMotionTest) TestMoveClassMethodThatReadsGlobalVar() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"class Foo {\n"+
						"  baz() { x; }\n"+
						"}\n"+
						"var x = 'x';\n").
				AddChunk("var y = new Foo(); y.baz();").
				Build()),
		c.Expected(
			CrossChunkMethodMotionTest_STUB_DECLARATIONS+
				"class Foo {}\n"+
				"Foo.prototype.baz = JSCompiler_stubMethod(0);\n"+
				"var x = 'x';\n"+
				"\n",
			// Chunk 2
			"Foo.prototype.baz = JSCompiler_unstubMethod(0, function(){x});\n"+
				"var y = new Foo(); y.baz();\n"))
}

// Read of a local is fine.
// Java: CrossChunkMethodMotionTest.movePrototypeMethodThatReferencesOnlyLocalVariables()
func (c *CrossChunkMethodMotionTest) TestMovePrototypeMethodThatReferencesOnlyLocalVariables() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"function Foo() {}\n"+
						"Foo.prototype.baz = function(){var x = 1;x};\n").
				AddChunk("var y = new Foo(); y.baz();").
				Build()),
		c.Expected(
			CrossChunkMethodMotionTest_STUB_DECLARATIONS+
				"function Foo() {}\n"+
				"Foo.prototype.baz = JSCompiler_stubMethod(0);\n",
			// Chunk 2
			"Foo.prototype.baz = JSCompiler_unstubMethod(\n"+
				"    0, function(){var x = 1; x});\n"+
				"var y = new Foo(); y.baz();\n"))
}

// Read of a local is fine.
// Java: CrossChunkMethodMotionTest.moveClassMethodThatReferencesOnlyLocalVariables()
func (c *CrossChunkMethodMotionTest) TestMoveClassMethodThatReferencesOnlyLocalVariables() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk("class Foo { baz() {var x = 1; x; } }").
				AddChunk("var y = new Foo(); y.baz();").
				Build()),
		c.Expected(
			CrossChunkMethodMotionTest_STUB_DECLARATIONS+
				"class Foo {}\n"+
				"Foo.prototype.baz = JSCompiler_stubMethod(0);\n",
			// Chunk 2
			"Foo.prototype.baz = JSCompiler_unstubMethod(\n"+
				"    0, function(){var x = 1; x});\n"+
				"var y = new Foo(); y.baz();\n"))
}

// An anonymous inner function reading a closure variable is fine.
// Java: CrossChunkMethodMotionTest.movePrototypeMethodContainingClosureOverLocalVariable()
func (c *CrossChunkMethodMotionTest) TestMovePrototypeMethodContainingClosureOverLocalVariable() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"function Foo() {}\n"+
						"Foo.prototype.baz = function() {\n"+
						"  var x = 1;\n"+
						"  return function(){x}\n"+
						"};\n").
				AddChunk("var y = new Foo(); y.baz();").
				Build()),
		c.Expected(
			CrossChunkMethodMotionTest_STUB_DECLARATIONS+
				"function Foo() {}\n"+
				"Foo.prototype.baz = JSCompiler_stubMethod(0);\n",
			// Chunk 2
			"Foo.prototype.baz = JSCompiler_unstubMethod(\n"+
				"    0, function(){var x = 1; return function(){x}});\n"+
				"var y = new Foo(); y.baz();\n"))
}

// Java: CrossChunkMethodMotionTest.moveClassMethodContainingClosureOverLocalVariable()
func (c *CrossChunkMethodMotionTest) TestMoveClassMethodContainingClosureOverLocalVariable() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"class Foo {\n"+
						"  baz() {\n"+
						"    var x = 1;\n"+
						"    return function(){x}\n"+
						"  }\n"+
						"}\n").
				AddChunk("var y = new Foo(); y.baz();").
				Build()),
		c.Expected(
			CrossChunkMethodMotionTest_STUB_DECLARATIONS+
				"class Foo {}\n"+
				"Foo.prototype.baz = JSCompiler_stubMethod(0);\n",
			// Chunk 2
			"Foo.prototype.baz = JSCompiler_unstubMethod(\n"+
				"    0, function(){var x = 1; return function(){x}});\n"+
				"var y = new Foo(); y.baz();\n"))
}

// Java: CrossChunkMethodMotionTest.staticBlockWithoutMethodReference()
func (c *CrossChunkMethodMotionTest) TestStaticBlockWithoutMethodReference() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"class Bar {\n"+
						"  method() {\n"+
						"  }\n"+
						"  static {\n"+
						"  }\n"+
						"}\n"+
						"class Foo extends Bar {\n"+
						"  method2() {\n"+
						"    return () => super.method();\n"+
						"  }\n"+
						"}\n").
				AddChunk("(new Foo).method2()").
				Build()),
		c.Expected(
			CrossChunkMethodMotionTest_STUB_DECLARATIONS+
				"class Bar {\n"+
				"  static {\n"+
				"  }\n"+
				"}\n"+
				"Bar.prototype.method = JSCompiler_stubMethod(0);\n"+
				"class Foo extends Bar {\n"+
				"  method2() {\n"+
				"    return () => super.method();\n"+
				"  }\n"+
				"}\n",
			// Chunk 2
			"Bar.prototype.method = JSCompiler_unstubMethod(0, function() {});\n"+
				"(new Foo).method2()\n"))
}

// Java: CrossChunkMethodMotionTest.referenceToMethodInOwnStaticBlock()
func (c *CrossChunkMethodMotionTest) TestReferenceToMethodInOwnStaticBlock() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"class Bar {\n" +
						"  method() {\n" +
						"  }\n" +
						"  static {\n" +
						"    this.prototype.method;\n" +
						"  }\n" +
						"}\n" +
						"class Foo extends Bar {\n" +
						"  method2() {\n" +
						"    return () => super.method();\n" +
						"  }\n" +
						"}\n").
				AddChunk("(new Foo).method2()").
				Build()))
}

// Java: CrossChunkMethodMotionTest.staticBlockReferenceToMethodInDifferentClassNoMovement()
func (c *CrossChunkMethodMotionTest) TestStaticBlockReferenceToMethodInDifferentClassNoMovement() {
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"class Bar {\n" +
						"  method() {}\n" +
						"}\n" +
						"class Foo extends Bar {\n" +
						"  static {\n" +
						"    (new Bar).method();\n" +
						"  }\n" +
						"  method2() {\n" +
						"    return () => super.method();\n" +
						"  }\n" +
						"}\n").
				AddChunk("(new Foo).method2()").
				Build()))
}

// Java: CrossChunkMethodMotionTest.staticBlockReferenceToMethodInDifferentClassWithMovement()
func (c *CrossChunkMethodMotionTest) TestStaticBlockReferenceToMethodInDifferentClassWithMovement() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"class Bar {\n"+
						"  method() {}\n"+
						"}\n").
				AddChunk(
					"class Foo extends Bar {\n"+
						"  static {\n"+
						"    (new Bar()).method();\n"+
						"  }\n"+
						"  method2() {\n"+
						"    return () => { return super.method(); };\n"+
						"  }\n"+
						"}\n").
				AddChunk("(new Foo).method2()").
				Build()),
		c.Expected(
			CrossChunkMethodMotionTest_STUB_DECLARATIONS+
				"class Bar {\n"+
				"}\n"+
				"Bar.prototype.method = JSCompiler_stubMethod(0);\n",
			// Chunk 2
			"Bar.prototype.method = JSCompiler_unstubMethod(0, function() {});\n"+
				"class Foo extends Bar {\n"+
				"  static {\n"+
				"    (new Bar()).method();\n"+
				"  }\n"+
				"  method2() {\n"+
				"    return () => { return super.method(); };\n"+
				"  }\n"+
				"}\n",
			// Chunk 3
			"(new Foo()).method2();"))
}

// Java: CrossChunkMethodMotionTest.testIssue600()
func (c *CrossChunkMethodMotionTest) TestIssue600() {
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"var jQuery1 = (function() {\n" +
						"  var jQuery2 = function() {};\n" +
						"  var theLoneliestNumber = 1;\n" +
						"  jQuery2.prototype = {\n" +
						"    size: function() {\n" +
						"      return theLoneliestNumber;\n" +
						"    }\n" +
						"  };\n" +
						"  return jQuery2;\n" +
						"})();\n").
				AddChunk(
					"(function() {\n" +
						"  var div = jQuery1('div');\n" +
						"  div.size();\n" +
						"})();\n").
				Build()))
}

// Java: CrossChunkMethodMotionTest.testIssue600b()
func (c *CrossChunkMethodMotionTest) TestIssue600b() {
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"var jQuery1 = (function() {\n" +
						"  var jQuery2 = function() {};\n" +
						"  jQuery2.prototype = {\n" +
						"    size: function() {\n" +
						"      return 1;\n" +
						"    }\n" +
						"  };\n" +
						"  return jQuery2;\n" +
						"})();\n").
				AddChunk(
					"(function() {\n" +
						"  var div = jQuery1('div');\n" +
						"  div.size();\n" +
						"})();\n").
				Build()))
}

// Java: CrossChunkMethodMotionTest.testIssue600c()
func (c *CrossChunkMethodMotionTest) TestIssue600c() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"var jQuery2 = function() {};\n"+
						"jQuery2.prototype = {\n"+
						"  size: function() {\n"+
						"    return 1;\n"+
						"  }\n"+
						"};\n").
				AddChunk(
					"(function() {\n"+
						"  var div = jQuery2('div');\n"+
						"  div.size();\n"+
						"})();\n").
				Build()),
		c.Expected(
			CrossChunkMethodMotionTest_STUB_DECLARATIONS+
				"var jQuery2 = function() {};\n"+
				"jQuery2.prototype = {\n"+
				"  size: JSCompiler_stubMethod(0)\n"+
				"};\n",
			// Chunk 2
			"jQuery2.prototype.size=\n"+
				"    JSCompiler_unstubMethod(0,function(){return 1});\n"+
				"(function() {\n"+
				"  var div = jQuery2('div');\n"+
				"  div.size();\n"+
				"})();\n"))
}

// Java: CrossChunkMethodMotionTest.testIssue600d()
func (c *CrossChunkMethodMotionTest) TestIssue600d() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"var jQuery2 = function() {};\n"+
						"(function() {\n"+
						"  jQuery2.prototype = {\n"+
						"    size: function() {\n"+
						"      return 1;\n"+
						"    }\n"+
						"  };\n"+
						"})();\n").
				AddChunk(
					"(function() {\n"+
						"  var div = jQuery2('div');\n"+
						"  div.size();\n"+
						"})();\n").
				Build()),
		c.Expected(
			CrossChunkMethodMotionTest_STUB_DECLARATIONS+
				"var jQuery2 = function() {};\n"+
				"(function() {\n"+
				"  jQuery2.prototype = {\n"+
				"    size: JSCompiler_stubMethod(0)\n"+
				"  };\n"+
				"})();\n",
			"jQuery2.prototype.size=\n"+
				"    JSCompiler_unstubMethod(0,function(){return 1});\n"+
				"(function() {\n"+
				"  var div = jQuery2('div');\n"+
				"  div.size();\n"+
				"})();\n"))
}

// Java: CrossChunkMethodMotionTest.testIssue600e()
func (c *CrossChunkMethodMotionTest) TestIssue600e() {
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"var jQuery2 = function() {};\n" +
						"(function() {\n" +
						"  var theLoneliestNumber = 1;\n" +
						"  jQuery2.prototype = {\n" +
						"    size: function() {\n" +
						"      return theLoneliestNumber;\n" +
						"    }\n" +
						"  };\n" +
						"})();\n").
				AddChunk(
					"(function() {\n" +
						"  var div = jQuery2('div');\n" +
						"  div.size();\n" +
						"})();\n").
				Build()))
}

// Java: CrossChunkMethodMotionTest.testPrototypeOfThisAssign()
func (c *CrossChunkMethodMotionTest) TestPrototypeOfThisAssign() {
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"/** @constructor */\n" +
						"function F() {}\n").
				AddChunk("this.prototype.foo = function() {};").
				AddChunk("(new F()).foo();").
				Build()))
}

// Java: CrossChunkMethodMotionTest.testDestructuring()
func (c *CrossChunkMethodMotionTest) TestDestructuring() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"/** @constructor */\n"+
						"function F() {}\n"+
						"F.prototype.foo = function() {};\n").
				AddChunk("const {foo} = new F();").
				Build()),
		c.Expected(
			CrossChunkMethodMotionTest_STUB_DECLARATIONS+
				"/** @constructor */\n"+
				"function F() {}\n"+
				"F.prototype.foo = JSCompiler_stubMethod(0);\n",
			"F.prototype.foo = JSCompiler_unstubMethod(0, function(){});\n"+
				"const {foo} = new F();\n"))
}

// Java: CrossChunkMethodMotionTest.testDestructuringWithQuotedProp()
func (c *CrossChunkMethodMotionTest) TestDestructuringWithQuotedProp() {
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"/** @constructor */\n" +
						"function F() {}\n" +
						"F.prototype.foo = function() {};\n").
				AddChunk("const {'foo': foo} = new F();").
				Build()))
}

// Java: CrossChunkMethodMotionTest.testDestructuringWithComputedProp()
func (c *CrossChunkMethodMotionTest) TestDestructuringWithComputedProp() {
	// See https://github.com/google/closure-compiler/issues/3145
	c.TestSame(
		c.Srcs(
			JSChunkGraphBuilder_ForChain().
				AddChunk(
					"/** @constructor */\n" +
						"function F() {}\n" +
						"F.prototype['foo'] = function() {};\n").
				AddChunk("const {['foo']: foo} = new F();").
				Build()))
}
