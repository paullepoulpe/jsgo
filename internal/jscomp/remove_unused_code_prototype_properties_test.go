// Port of com/google/javascript/jscomp/RemoveUnusedCodePrototypePropertiesTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/RemoveUnusedCodePrototypePropertiesTest.java)

/*
 * Copyright 2017 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"testing"
)

type RemoveUnusedCodePrototypePropertiesTest struct {
	CompilerTestCaseBase
	keepLocals  bool
	keepGlobals bool
}

// Java: RemoveUnusedCodePrototypePropertiesTest()
func NewRemoveUnusedCodePrototypePropertiesTest() *RemoveUnusedCodePrototypePropertiesTest {
	c := &RemoveUnusedCodePrototypePropertiesTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, RemoveUnusedCodePrototypePropertiesTest_EXTERNS)
	c.keepLocals = true
	c.keepGlobals = false
	return c
}

func TestRemoveUnusedCodePrototypePropertiesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewRemoveUnusedCodePrototypePropertiesTest() })
}

var RemoveUnusedCodePrototypePropertiesTest_EXTERNS = CompilerTestCase_MINIMAL_EXTERNS +
	"var window;\n" +
	"var Math = {};\n" +
	"Math.random = function() {};\n" +
	"function alert(x) {}\n" +
	"function externFunction() {}\n" +
	"externFunction.prototype.externPropName;\n" +
	"var mExtern;\n" +
	"mExtern.bExtern;\n" +
	"mExtern['cExtern'];\n" +
	"\n" +
	"/** @const */\n" +
	"var goog = {};\n" +
	"goog.reflect.objectProperty = function(name) { };\n"

// Java: RemoveUnusedCodePrototypePropertiesTest.getProcessor(Compiler)
func (c *RemoveUnusedCodePrototypePropertiesTest) GetProcessor(compiler *Compiler) CompilerPass {
	return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
		NewRemoveUnusedCode(compiler, RemoveUnusedCodeOptions{
			RemoveLocalVars:                 !c.keepLocals,
			RemoveGlobals:                   !c.keepGlobals,
			RemoveUnusedPrototypeProperties: true,
		}).Process(externs, root)
	})
}

// Java: RemoveUnusedCodePrototypePropertiesTest.getOptions()
func (c *RemoveUnusedCodePrototypePropertiesTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetWarningLevel(DiagnosticGroups_MODULE_LOAD, CheckLevel_OFF)
	return options
}

// Java: RemoveUnusedCodePrototypePropertiesTest.setUp()
func (c *RemoveUnusedCodePrototypePropertiesTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	// Allow testing of features that aren't fully supported for output yet.
	c.EnableNormalize()
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.EnableGatherExternProperties()
	c.keepLocals = true
	c.keepGlobals = false
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testClassPropertiesNotRemoved()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestClassPropertiesNotRemoved() {
	c.keepGlobals = true
	// This whole test class runs with removeUnusedClassProperties disabled.
	c.TestSame("/** @constructor */ function C() {} C.unused = 3;")
	c.TestSame(
		"/** @constructor */ function C() {} Object.defineProperties(C, {unused: {value: 3}});")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testUnusedPrototypeFieldReference()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestUnusedPrototypeFieldReference() {
	c.Test(
		"function C() {} C.prototype.x; new C();", // x is not actually read
		"function C() {}                new C();")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testUnusedReferenceToFieldWithGetter()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestUnusedReferenceToFieldWithGetter() {
	// Reference to a field with a getter should not be removed unless we know it has no side
	// effects.
	// TODO(bradfordcsmith): Implement removal for the no-side-effect cases.
	c.TestSame("function C() {} C.prototype = { get x() {} }; new C().x")
	c.TestSame("function C() {} C.prototype = { get x() { alert('x'); } }; new C().x")
	c.TestSame("class C { get x() {} } new C().x;")
	c.TestSame("class C { get x() { alert('x'); } } new C().x")
	c.TestSame("let c = { get x() {} }; c.x;")
	c.TestSame("let c = { get x() { alert('x'); } }; c.x;")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testAnonymousPrototypePropertyRemoved()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestAnonymousPrototypePropertyRemoved() {
	c.Test("({}.prototype.x = 5, externFunction())", "externFunction();")
	c.Test("({}).prototype.x = 5;", "")
	c.Test("({}).prototype.x = externFunction();", "externFunction();")
	c.Test("externFunction({}.prototype.x = 5);", "externFunction(5);")
	c.Test("externFunction().prototype.x = 5;", "externFunction();")
	c.Test("externFunction().prototype.x = externFunction();", "externFunction(), externFunction();")

	// make sure an expression with '.prototype' is traversed when it should be and not when it
	// shouldn't.
	c.Test(
		"function C() {} externFunction(C).prototype.x = 5;", // preserve format
		"function C() {} externFunction(C);")
	c.Test("function C() {} ({ C: C }).prototype.x = 5;", "")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testAnonymousPrototypePropertyNoRemoveSideEffect1()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestAnonymousPrototypePropertyNoRemoveSideEffect1() {
	c.Test(
		"function A() { // preserve format\n"+
			"  externFunction('me');\n"+
			"  return function(){}\n"+
			"}\n"+
			"A().prototype.foo = function() {};\n",
		"function A() { // preserve format\n"+
			"  externFunction('me');\n"+
			"  return function(){}\n"+
			"}\n"+
			"A();\n")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testAnonymousPrototypePropertyNoRemoveSideEffect2()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestAnonymousPrototypePropertyNoRemoveSideEffect2() {
	c.Test(
		"function A() { externFunction('me'); return function(){}; } A().prototype.foo++;",
		"function A() { externFunction('me'); return function(){}; } A();")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testIncPrototype()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestIncPrototype() {
	c.Test("function A() {} A.prototype.x = 1; A.prototype.x++;", "")
	c.Test(
		"function A() {} A.prototype.x = 1; A.prototype.x++; new A();",
		"function A() {}                                     new A();")
	c.Test("externFunction().prototype.x++", "externFunction()")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testRenamePropertyFunctionTest()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestRenamePropertyFunctionTest() {
	c.Test(
		"function C() {} // preserve formatting\n"+
			"C.prototype.unreferenced = function() {};\n"+
			"C.prototype.renamed = function() {};\n"+
			"JSCompiler_renameProperty('renamed');\n"+
			"new C();\n",
		"function C() {} // preserve formatting\n"+
			"C.prototype.renamed = function() {};\n"+
			"JSCompiler_renameProperty('renamed');\n"+
			"new C();\n")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testNonPrototypePropertiesAreKept()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestNonPrototypePropertiesAreKept() {
	// foo cannot be removed because it is called
	// x cannot be removed because a property is set on it and we don't know where it comes from
	// x.a cannot be removed because we don't know where x comes from.
	// x.prototype.b *can* be removed because we consider it safe to remove prototype properties
	// that have no references.
	c.Test(
		"function foo(x) { x.a = 1; x.prototype.b = 2; }; foo({});",
		"function foo(x) { x.a = 1; }; foo({});")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testAnalyzePrototypeProperties()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestAnalyzePrototypeProperties() {
	// Basic removal for prototype properties
	c.Test(
		"function e(){}\n"+
			"e.prototype.a = function(){};\n"+
			"e.prototype.b = function(){};\n"+
			"var x = new e; x.a()\n",
		"function e(){}\n"+
			"e.prototype.a = function(){};\n"+
			"var x = new e; x.a()\n")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testObjectLiteralPrototype()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestObjectLiteralPrototype() {
	c.Test(
		"function e(){} e.prototype = {a: function(){}, b: function(){}}; var x = new e; x.a()",
		"function e(){} e.prototype = {a: function(){}                 }; var x = new e; x.a()")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testObjectLiteralPrototypeUnusedPropDefinitionWithSideEffects()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestObjectLiteralPrototypeUnusedPropDefinitionWithSideEffects() {
	c.Test(
		"function e(){} e.prototype = {a: alert('a'), b: function(){}}; new e;",
		"function e(){} e.prototype = {a: alert('a')                 }; new e;")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testPropertiesDefinedInExterns()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestPropertiesDefinedInExterns() {
	c.Test(
		"function e(){}\n"+
			"e.prototype.a = function(){};\n"+
			"e.prototype.bExtern = function(){};\n"+
			"var x = new e;x.a()\n",
		"function e(){}\n"+
			"e.prototype.a = function(){};\n"+
			"e.prototype.bExtern = function(){};\n"+
			"var x = new e; x.a()\n")
	c.TestSame(
		"function e(){}\n" +
			"e.prototype = {a: function(){}, bExtern: function(){}};\n" +
			"var x = new e; x.a()\n")

	c.TestSame(
		"class C {\n" +
			"  constructor() {}\n" +
			"  bExtern() {} // property name defined in externs.\n" +
			"}\n" +
			"new C();\n")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testAliasing1()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestAliasing1() {
	// Aliasing a property is not enough for it to count as used
	c.Test(
		"function e(){}\n"+
			"e.prototype.method1 = function(){};\n"+
			"e.prototype.method2 = function(){};\n"+
			"// aliases\n"+
			"e.prototype.alias1 = e.prototype.method1;\n"+
			"e.prototype.alias2 = e.prototype.method2;\n"+
			"var x = new e; x.method1()\n",
		"function e(){}\n"+
			"e.prototype.method1 = function(){};\n"+
			"var x = new e; x.method1()\n")

	// Using an alias should keep it
	c.Test(
		"function e(){}\n"+
			"e.prototype.method1 = function(){};\n"+
			"e.prototype.method2 = function(){};\n"+
			"// aliases\n"+
			"e.prototype.alias1 = e.prototype.method1;\n"+
			"e.prototype.alias2 = e.prototype.method2;\n"+
			"var x=new e; x.alias1()\n",
		"function e(){}\n"+
			"e.prototype.method1 = function(){};\n"+
			"e.prototype.alias1 = e.prototype.method1;\n"+
			"var x = new e; x.alias1()\n")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testAliasing2()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestAliasing2() {
	// Aliasing a property is not enough for it to count as used
	c.Test(
		"function e(){}\n"+
			"e.prototype.method1 = function(){};\n"+
			"// aliases\n"+
			"e.prototype.alias1 = e.prototype.method1;\n"+
			"(new e).method1()\n",
		"function e(){}\n"+
			"e.prototype.method1 = function(){};\n"+
			"(new e).method1()\n")

	// Using an alias should keep it
	c.TestSame(
		"function e(){}\n" +
			"e.prototype.method1 = function(){};\n" +
			"// aliases\n" +
			"e.prototype.alias1 = e.prototype.method1;\n" +
			"(new e).alias1()\n")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testAliasing3()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestAliasing3() {
	// Aliasing a property is not enough for it to count as used
	c.TestSame(
		"function e(){}\n" +
			"e.prototype.method1 = function(){};\n" +
			"e.prototype.method2 = function(){};\n" +
			"// aliases\n" +
			"e.prototype['alias1'] = e.prototype.method1;\n" +
			"e.prototype['alias2'] = e.prototype.method2;\n" +
			"new e;\n")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testAliasing4()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestAliasing4() {
	// Aliasing a property is not enough for it to count as used
	c.Test(
		"function e(){}\n"+
			"e.prototype['alias1'] = e.prototype.method1 = function(){};\n"+
			"e.prototype['alias2'] = e.prototype.method2 = function(){};\n"+
			"new e;\n",
		"function e(){}\n"+
			"e.prototype['alias1'] = function(){};\n"+
			"e.prototype['alias2'] = function(){};\n"+
			"new e;\n")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testAliasing5()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestAliasing5() {
	// An exported alias must preserved any referenced values in the
	// referenced function.
	c.TestSame(
		"function e(){}\n" +
			"e.prototype.method1 = function(){this.method2()};\n" +
			"e.prototype.method2 = function(){};\n" +
			"// aliases\n" +
			"e.prototype['alias1'] = e.prototype.method1;\n" +
			"new e;\n")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testAliasing6()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestAliasing6() {
	// An exported alias must preserved any referenced values in the
	// referenced function.
	c.Test(
		"function e(){}\n"+
			"e.prototype.method1 = function(){this.method2()};\n"+
			"e.prototype.method2 = function(){};\n"+
			"// aliases\n"+
			"window['alias1'] = e.prototype.method1;\n",
		"function e(){}\n"+
			"e.prototype.method1=function(){this.method2()};\n"+
			"e.prototype.method2=function(){};\n"+
			"window['alias1']=e.prototype.method1;\n")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testAliasing7()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestAliasing7() {
	// An exported alias must preserved any referenced values in the
	// referenced function.
	c.Test(
		"function e(){}\n"+
			"e.prototype['alias1'] = e.prototype.method1 = function(){this.method2()};\n"+
			"e.prototype.method2 = function(){};\n"+
			"new e;\n",
		"function e(){}\n"+
			"e.prototype['alias1'] = function(){this.method2()};\n"+
			"e.prototype.method2 = function(){};\n"+
			"new e;\n")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testExportedMethodsByNamingConvention()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestExportedMethodsByNamingConvention() {
	classAndItsMethodAliasedAsExtern :=
		"function Foo() {}\n" +
			"Foo.prototype.method = function() {};\n" +
			"// not removed\n" +
			"Foo.prototype.unused = function() {};\n" +
			"// removed\n" +
			"var _externInstance = new Foo();\n" +
			"Foo.prototype._externMethod = Foo.prototype.method // aliased here\n"

	compiled :=
		"function Foo(){}\n" +
			"Foo.prototype.method = function(){};\n" +
			"var _externInstance = new Foo;\n" +
			"Foo.prototype._externMethod = Foo.prototype.method\n"

	c.Test(classAndItsMethodAliasedAsExtern, compiled)
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testExportedMethodsByNamingConventionAlwaysExported()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestExportedMethodsByNamingConventionAlwaysExported() {
	classAndItsMethodAliasedAsExtern :=
		"function Foo() {}\n" +
			"Foo.prototype.method = function() {};\n" +
			"// not removed\n" +
			"Foo.prototype.unused = function() {};\n" +
			"// removed\n" +
			"var _externInstance = new Foo();\n" +
			"Foo.prototype._externMethod = Foo.prototype.method // aliased here\n"

	compiled :=
		"function Foo(){}\n" +
			"Foo.prototype.method = function(){};\n" +
			"var _externInstance = new Foo;\n" +
			"Foo.prototype._externMethod = Foo.prototype.method\n"

	c.Test(classAndItsMethodAliasedAsExtern, compiled)
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testExternMethodsFromExternsFile()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestExternMethodsFromExternsFile() {
	classAndItsMethodAliasedAsExtern :=
		"function Foo() {}\n" +
			"Foo.prototype.bar_ = function() {};\n" +
			"// not removed\n" +
			"Foo.prototype.unused = function() {};\n" +
			"// removed\n" +
			"var instance = new Foo;\n" +
			"Foo.prototype.externPropName = Foo.prototype.bar_ // aliased here\n"

	compiled :=
		"function Foo(){}\n" +
			"Foo.prototype.bar_ = function(){};\n" +
			"new Foo;\n" +
			"Foo.prototype.externPropName = Foo.prototype.bar_\n"

	c.Test(classAndItsMethodAliasedAsExtern, compiled)
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testPropertyReferenceGraph()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestPropertyReferenceGraph() {
	// test a prototype property graph that looks like so:
	// b -> a, c -> b, c -> a, d -> c, e -> a, e -> f
	constructor := "function Foo() {}"
	defA := "Foo.prototype.a = function() { Foo.superClass_.a.call(this); };"
	defB := "Foo.prototype.b = function() { this.a(); };"
	defC :=
		"Foo.prototype.c = function() {\n" +
			"Foo.superClass_.c.call(this); this.b(); this.a(); };\n"
	defD := "Foo.prototype.d = function() { this.c(); };"
	defE := "Foo.prototype.e = function() { this.a(); this.f(); };"
	defF := "Foo.prototype.f = function() { };"
	fullClassDef := constructor + defA + defB + defC + defD + defE + defF

	// ensure that all prototypes are compiled out if none are used
	c.Test(fullClassDef, "")

	// make sure that the right prototypes are called for each use
	callA := "(new Foo()).a();"
	callB := "(new Foo()).b();"
	callC := "(new Foo()).c();"
	callD := "(new Foo()).d();"
	callE := "(new Foo()).e();"
	callF := "(new Foo()).f();"
	c.Test(fullClassDef+callA, constructor+defA+callA)
	c.Test(fullClassDef+callB, constructor+defA+defB+callB)
	c.Test(fullClassDef+callC, constructor+defA+defB+defC+callC)
	c.Test(fullClassDef+callD, constructor+defA+defB+defC+defD+callD)
	c.Test(fullClassDef+callE, constructor+defA+defE+defF+callE)
	c.Test(fullClassDef+callF, constructor+defF+callF)

	c.Test(fullClassDef+callA+callC, constructor+defA+defB+defC+callA+callC)
	c.Test(fullClassDef+callB+callC, constructor+defA+defB+defC+callB+callC)
	c.Test(
		fullClassDef+callA+callB+callC,
		constructor+defA+defB+defC+callA+callB+callC)
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testPropertiesDefinedWithGetElem()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestPropertiesDefinedWithGetElem() {
	c.TestSame("function Foo() {} Foo.prototype['elem'] = function() {}; new Foo;")
	c.TestSame("function Foo() {} Foo.prototype[1 + 1] = function() {}; new Foo;")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testQuotedProperties()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestQuotedProperties() {
	// Basic removal for prototype replacement
	c.TestSame("function e(){} e.prototype = {'a': function(){}, 'b': function(){}}; new e;")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testNeverRemoveImplicitlyUsedProperties()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestNeverRemoveImplicitlyUsedProperties() {
	c.TestSame(
		"function Foo() {}\n" +
			"Foo.prototype.length = 3;\n" +
			"Foo.prototype.toString = function() { return 'Foo'; };\n" +
			"Foo.prototype.valueOf = function() { return 'Foo'; };\n" +
			"new Foo;\n")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testPropertyDefinedInBranch()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestPropertyDefinedInBranch() {
	c.Test("function Foo() {} if (true) Foo.prototype.baz = function() {};", "if (true);")
	c.Test("function Foo() {} while (true) Foo.prototype.baz = function() {};", "while (true);")
	c.Test("function Foo() {} for (;;) Foo.prototype.baz = function() {};", "for (;;);")
	c.Test(
		"function Foo() {} do Foo.prototype.baz = function() {}; while(true);", "do; while(true);")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testUsingAnonymousObjectsToDefeatRemoval()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestUsingAnonymousObjectsToDefeatRemoval() {
	c.Test("function Foo() {} Foo.prototype.baz = 3; new Foo;", "function Foo() {} new Foo;")
	c.TestSame("function Foo() {} Foo.prototype.baz = 3; new Foo; var x = {}; x.baz;")
	c.TestSame("function Foo() {} Foo.prototype.baz = 3; new Foo; var x = {baz: 5}; x;")
	// quoted properties still prevent removal
	c.TestSame("function Foo() {} Foo.prototype.baz = 3; new Foo; var x = {'baz': 5}; x;")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testGlobalFunctionsInGraph()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestGlobalFunctionsInGraph() {
	c.Test(
		"var x = function() { (new Foo).baz(); };\n"+
			"var y = function() { x(); };\n"+
			"function Foo() {}\n"+
			"Foo.prototype.baz = function() { y(); };\n",
		"")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testGlobalFunctionsInGraph2()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestGlobalFunctionsInGraph2() {
	c.Test(
		"var x = function() { (new Foo).baz(); };\n"+
			"var y = function() { x(); };\n"+
			"function Foo() { this.baz(); }\n"+
			"Foo.prototype.baz = function() { y(); };\n",
		"")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testGlobalFunctionsInGraph3()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestGlobalFunctionsInGraph3() {
	c.Test(
		"var x = function() { (new Foo).baz(); };\n"+
			"var y = function() { x(); };\n"+
			"function Foo() { this.baz(); }\n"+
			"Foo.prototype.baz = function() { x(); };\n",
		"")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testGlobalFunctionsInGraph4()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestGlobalFunctionsInGraph4() {
	c.Test(
		"var x = function() { (new Foo).baz(); };\n"+
			"var y = function() { x(); };\n"+
			"function Foo() { Foo.prototype.baz = function() { y(); }; }\n",
		"")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testGlobalFunctionsInGraph5()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestGlobalFunctionsInGraph5() {
	c.Test(
		"function Foo() {}\n"+
			"Foo.prototype.methodA = function() {};\n"+
			"function x() { (new Foo).methodA(); }\n"+
			"Foo.prototype.methodB = function() { x(); };\n",
		"")

	c.keepGlobals = true
	c.Test(
		"function Foo() {}\n"+
			"Foo.prototype.methodA = function() {};\n"+
			"function x() { (new Foo).methodA(); }\n"+
			"Foo.prototype.methodB = function() { x(); };\n",
		"function Foo() {}\n"+
			"Foo.prototype.methodA = function() {};\n"+
			"function x() { (new Foo).methodA(); }\n")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testGlobalFunctionsInGraph6()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestGlobalFunctionsInGraph6() {
	c.TestSame(
		"function Foo() {}\n" +
			"Foo.prototype.methodA = function() {};\n" +
			"function x() { (new Foo).methodA(); }\n" +
			"Foo.prototype.methodB = function() { x(); };\n" +
			"(new Foo).methodB();\n")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testGlobalFunctionsInGraph7()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestGlobalFunctionsInGraph7() {
	c.keepGlobals = true
	c.TestSame("function Foo() {} Foo.prototype.methodA = function() {}; this.methodA();")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testGlobalFunctionsInGraph8()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestGlobalFunctionsInGraph8() {
	c.Test(
		"let x = function() { (new Foo).baz(); };\n"+
			"const y = function() { x(); };\n"+
			"function Foo() { Foo.prototype.baz = function() { y(); }; }\n",
		"")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testGetterBaseline()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestGetterBaseline() {
	c.keepGlobals = true
	c.Test(
		"function Foo() {}\n"+
			"Foo.prototype = {\n"+
			"  methodA: function() {},\n"+
			"  methodB: function() { x(); }\n"+
			"};\n"+
			"function x() { (new Foo).methodA(); }\n",
		"function Foo() {}\n"+
			"Foo.prototype = {\n"+
			"  methodA: function() {}\n"+
			"};\n"+
			"function x() { (new Foo).methodA(); }\n")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testGetter1()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestGetter1() {
	c.Test(
		"function Foo() {}\n"+
			"Foo.prototype = {\n"+
			"  get methodA() {},\n"+
			"  get methodB() { x(); }\n"+
			"};\n"+
			"function x() { (new Foo).methodA; }\n"+
			"new Foo();\n",
		"function Foo() {}\n"+
			"// x() and all methods of Foo removed.\n"+
			"Foo.prototype = {};\n"+
			"new Foo();\n")

	c.keepGlobals = true
	c.Test(
		"function Foo() {}\n"+
			"Foo.prototype = {\n"+
			"  get methodA() {},\n"+
			"  get methodB() { x(); }\n"+
			"};\n"+
			"function x() { (new Foo).methodA; }\n",
		"function Foo() {}\n"+
			"Foo.prototype = {\n"+
			"  get methodA() {}\n"+
			"};\n"+
			"// x() keeps methodA alive\n"+
			"function x() { (new Foo).methodA; }\n")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testGetter2()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestGetter2() {
	c.keepGlobals = true
	c.Test(
		"function Foo() {}\n"+
			"Foo.prototype = {\n"+
			"  get methodA() {},\n"+
			"  set methodA(a) {},\n"+
			"  get methodB() { x(); },\n"+
			"  set methodB(a) { x(); }\n"+
			"};\n"+
			"function x() { (new Foo).methodA; }\n",
		"function Foo() {}\n"+
			"Foo.prototype = {\n"+
			"  get methodA() {},\n"+
			"  set methodA(a) {}\n"+
			"};\n"+
			"function x() { (new Foo).methodA; }\n")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testHook1()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestHook1() {
	c.Test(
		"/** @constructor */ function Foo() {}\n"+
			"Foo.prototype.method1 =\n"+
			"    Math.random()\n"+
			"        ? function() { this.method2(); }\n"+
			"        : function() { this.method3(); };\n"+
			"Foo.prototype.method2 = function() {};\n"+
			"Foo.prototype.method3 = function() {};\n",
		"")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testHook2()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestHook2() {
	c.TestSame(
		"/** @constructor */ function Foo() {}\n" +
			"Foo.prototype.method1 =\n" +
			"    Math.random()\n" +
			"        ? function() { this.method2(); }\n" +
			"        : function() { this.method3(); };\n" +
			"Foo.prototype.method2 = function() {};\n" +
			"Foo.prototype.method3 = function() {};\n" +
			"(new Foo()).method1();\n")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testDestructuringProperty()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestDestructuringProperty() {
	// Makes the cases below shorter because we don't have to add references
	// to globals to keep them around and just test prototype property removal.
	c.keepGlobals = true
	c.Test(
		"function Foo() {} Foo.prototype.a = function() {}; var {} = new Foo();",
		"function Foo() {} var {} = new Foo();")

	c.Test(
		"function Foo() {}\n"+
			"Foo.prototype.a = function() {};\n"+
			"Foo.prototype.b = function() {}\n"+
			"var {a} = new Foo();\n",
		"function Foo() {}\n"+
			"Foo.prototype.a = function() {};\n"+
			"var {a} = new Foo();\n")

	c.Test(
		"function Foo() {}\n"+
			"Foo.prototype.a = function() {};\n"+
			"Foo.prototype.b = function() {}\n"+
			"var {a:x} = new Foo();\n",
		"function Foo() {}\n"+
			"Foo.prototype.a = function() {};\n"+
			"var {a:x} = new Foo();\n")

	c.TestSame(
		"function Foo() {}\n" +
			"Foo.prototype.a = function() {};\n" +
			"Foo.prototype.b = function() {}\n" +
			"var {a, b} = new Foo();\n")

	c.TestSame(
		"function Foo() {}\n" +
			"Foo.prototype.a = function() {};\n" +
			"Foo.prototype.b = function() {}\n" +
			"var {a:x, b:y} = new Foo();\n")

	c.TestSame(
		"function Foo() {} // preserve newlines\n" +
			"Foo.prototype.a = function() {};\n" +
			"let x;\n" +
			"({a:x} = new Foo());\n")

	c.TestSame(
		"function Foo() {}\n" +
			"Foo.prototype.a = function() {};\n" +
			"function f({a:x}) { x; }; f(new Foo());\n")

	c.TestSame(
		"function Foo() {}\n" +
			"Foo.prototype.a = function() {};\n" +
			"var {a : x = 3} = new Foo();\n")

	c.Test(
		"function Foo() {}\n"+
			"Foo.prototype.a = function() {};\n"+
			"Foo.prototype.b = function() {}\n"+
			"var {a : a = 3} = new Foo();\n",
		"function Foo() {}\n"+
			"Foo.prototype.a = function() {};\n"+
			"var {a : a = 3} = new Foo();\n")

	c.TestSame(
		"function Foo() {}\n" +
			"Foo.prototype.a = function() {};\n" +
			"let { a : [b, c, d] } = new Foo();\n")

	c.TestSame(
		"function Foo() {}\n" +
			"Foo.prototype.a = function() {};\n" +
			"const { a : { b : { c : d = '' }}} = new Foo();\n")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testDestructuringRest()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestDestructuringRest() {
	// Makes the cases below shorter because we don't have to add references
	// to globals to keep them around and just test prototype property removal.
	c.keepGlobals = true

	c.TestSame(
		"function Foo() {}\n" +
			"Foo.prototype.a = function() {};\n" +
			"({ ...new Foo().a.b } = 0);\n")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testOptionalGetPropPreventsRemoval()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestOptionalGetPropPreventsRemoval() {
	c.Test(
		"class C {\n"+
			"  constructor() {\n"+
			"    this.x = 1;\n"+
			"  }\n"+
			"  optChainGetPropRef() {}\n"+
			"  optChainCallRef() {}\n"+
			"  unreferenced() {}\n"+
			"}\n"+
			"var c = new C;\n"+
			"c?.optChainGetPropRef()\n"+
			"c.optChainCallRef?.()\n"+
			"// no call to unreferenced()\n",
		"class C {\n"+
			"  constructor() {\n"+
			"    this.x = 1;\n"+
			"  }\n"+
			"  optChainGetPropRef() {} // kept\n"+
			"  optChainCallRef() {} // kept\n"+
			"// unreferenced() removed\n"+
			"}\n"+
			"var c = new C;\n"+
			"c?.optChainGetPropRef()\n"+
			"c.optChainCallRef?.()\n"+
			"// no call to unreferenced()\n")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testEs6Class()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestEs6Class() {
	c.TestSame(
		"class C {\n" +
			"  constructor() { // constructor is not removable\n" +
			"    this.x = 1;\n" +
			"  }\n" +
			"}\n" +
			"new C();\n")

	c.Test(
		"class C {\n"+
			"  constructor() {\n"+
			"    this.x = 1;\n"+
			"  }\n"+
			"  foo() {}\n"+
			"}\n"+
			"var c = new C\n",
		"class C {\n"+
			"  constructor() { // constructor is not removable\n"+
			"    this.x = 1;\n"+
			"  }\n"+
			"}\n"+
			"new C();\n")

	c.TestSame(
		"class C {\n" +
			"  constructor() {\n" +
			"    this.x = 1;\n" +
			"  }\n" +
			"  foo() {}\n" +
			"}\n" +
			"var c = new C\n" +
			"c.foo()\n")

	c.Test(
		"class C {\n"+
			"  constructor() {\n"+
			"    this.x = 1;\n"+
			"  }\n"+
			"  static foo() {}\n"+
			"}\n"+
			"new C;\n",
		"class C {\n"+
			"  constructor() { // constructor is not removable\n"+
			"    this.x = 1;\n"+
			"  }\n"+
			"// TODO(b/139319709): Remove this. static method removal is disabled.\n"+
			"  static foo() {}\n"+
			"}\n"+
			"new C();\n")

	c.Test(
		"class C {\n"+
			"  constructor() {\n"+
			"    this.x = 1;\n"+
			"  }\n"+
			"  get foo() {}\n"+
			"  set foo(val) {}\n"+
			"}\n"+
			"var c = new C\n",
		"class C { constructor() { this.x = 1; } } new C")

	c.TestSame(
		"class C {\n" +
			"  constructor() {\n" +
			"    this.x = 1;\n" +
			"  }\n" +
			"  get foo() {}\n" +
			"  set foo(val) {}\n" +
			"}\n" +
			"var c = new C;\n" +
			"c.foo = 3;\n")

	c.TestSame(
		"class C {\n" +
			"  constructor() {\n" +
			"    this.x = 1;\n" +
			"  }\n" +
			"  get foo() {}\n" +
			"  set foo(val) {}\n" +
			"}\n" +
			"var c = new C;\n" +
			"c.foo;\n")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testEs6Extends()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestEs6Extends() {
	c.TestSame(
		"class C {\n" +
			"  constructor() {\n" +
			"    this.x = 1;\n" +
			"  }\n" +
			"}\n" +
			"class D extends C {\n" +
			"  constructor() {}\n" +
			"}\n" +
			"new D();\n")

	c.TestSame(
		"class C {\n" +
			"  constructor() {\n" +
			"    this.x = 1;\n" +
			"  }\n" +
			"  foo() {}\n" +
			"}\n" +
			"class D extends C {\n" +
			"  constructor() {}\n" +
			"  foo() {\n" +
			"     return super.foo()\n" +
			"  }\n" +
			"}\n" +
			"var d = new D\n" +
			"d.foo()\n")

	c.Test(
		"class C {\n"+
			"  constructor() {\n"+
			"    this.x = 1;\n"+
			"  }\n"+
			"  foo() {}\n"+
			"}\n"+
			"class D extends C {\n"+
			"  constructor() {}\n"+
			"  foo() {\n"+
			"     return super.foo()\n"+
			"  }\n"+
			"}\n"+
			"var d = new D;\n",
		"class C {\n"+
			"  constructor() {\n"+
			"    this.x = 1;\n"+
			"  }\n"+
			"}\n"+
			"class D extends C {\n"+
			"  constructor() {}\n"+
			"}\n"+
			"new D;\n")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testAnonClasses()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestAnonClasses() {
	// Make sure class expression names are removed.
	c.keepLocals = false
	c.Test(
		"var C = class InnerC {\n"+
			"  constructor() {\n"+
			"    this.x = 1;\n"+
			"  }\n"+
			"  foo() {}\n"+
			"};\n"+
			"new C;\n",
		"var C = class { constructor() { this.x = 1; } }; new C;")

	c.TestSame(
		"var C = class {\n" +
			"  constructor() {\n" +
			"    this.x = 1;\n" +
			"  }\n" +
			"  foo() {}\n" +
			"}\n" +
			"var c = new C()\n" +
			"c.foo()\n")

	c.Test(
		"var C = class {}\n"+
			"C.D = class {\n"+
			"  constructor() {\n"+
			"    this.x = 1;\n"+
			"  }\n"+
			"  foo() {}\n"+
			"}\n"+
			"new C.D();\n",
		"var C = class {}\n"+
			"C.D = class{\n"+
			"  constructor() {\n"+
			"    this.x = 1;\n"+
			"  }\n"+
			"}\n"+
			"new C.D();\n")

	c.Test(
		"externFunction(class C { constructor() { } externPropName() { } })",
		"externFunction(class   { constructor() { } externPropName() { } })")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testBaseClassExpressionHasSideEffects()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestBaseClassExpressionHasSideEffects() {
	// Make sure names are removed from class expressions.
	c.keepLocals = false
	c.TestSame(
		"function getBaseClass() { return class {}; }\n" +
			"class C extends getBaseClass() {}\n")
	c.Test(
		"function getBaseClass() { return class {}; }\n"+
			"const C = class InnerC extends getBaseClass() {};\n",
		"function getBaseClass() { return class {}; }\n"+
			"(class extends getBaseClass() {})\n")
	c.Test(
		"function getBaseClass() { return class {}; }\n"+
			"let C;\n"+
			"C = class InnerC extends getBaseClass() {}\n",
		"function getBaseClass() { return class {}; }\n"+
			"(class extends getBaseClass() {})\n")
	c.Test(
		"function getBaseClass() { return class {}; }\n"+
			"externFunction(class InnerC extends getBaseClass() {})\n",
		"function getBaseClass() { return class {}; }\n"+
			"externFunction(class extends getBaseClass() {})\n")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testModules()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestModules() {
	c.TestSame("export default function(){}")
	c.TestSame("export class C {};")
	c.TestSame("class Bar {} export {Bar}")

	c.TestSame("import { square, diag } from '/lib';")
	c.TestSame("import * as lib from '/lib';")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testReflection_reflectProperty_pinsReflectedName()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestReflection_reflectProperty_pinsReflectedName() {
	c.TestSame(
		"/** @constructor */\n" +
			"function Foo() {}\n" +
			"Foo.prototype.handle = function(x, y) { alert(y); };\n" +
			"\n" +
			"goog.reflect.objectProperty('handle');\n" +
			"alert(new Foo());\n")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testReflection_reflectProperty_onlyPinsReflectedName()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestReflection_reflectProperty_onlyPinsReflectedName() {
	c.Test(
		"/** @constructor */\n"+
			"function Foo() {}\n"+
			"Foo.prototype.handle = function(x, y) { alert(y); };\n"+
			"\n"+
			"goog.reflect.objectProperty('not_handle');\n"+
			"alert(new Foo());\n",
		"/** @constructor */\n"+
			"function Foo() {}\n"+
			"\n"+
			"goog.reflect.objectProperty('not_handle');\n"+
			"alert(new Foo());\n")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testReflection_reflectProperty_onlyPinsReflectedName_whenNameMissing()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestReflection_reflectProperty_onlyPinsReflectedName_whenNameMissing() {
	c.Test(
		"/** @constructor */\n"+
			"function Foo() {}\n"+
			"Foo.prototype.handle = function(x, y) { alert(y); };\n"+
			"\n"+
			"goog.reflect.objectProperty();\n"+
			"alert(new Foo());\n",
		"/** @constructor */\n"+
			"function Foo() {}\n"+
			"\n"+
			"goog.reflect.objectProperty();\n"+
			"alert(new Foo());\n")
}

// Java: RemoveUnusedCodePrototypePropertiesTest.testPureOrBreakMyCode()
func (c *RemoveUnusedCodePrototypePropertiesTest) TestPureOrBreakMyCode() {
	c.Test(
		"/** @constructor */\n"+
			"function Foo() {}\n"+
			"Foo.prototype.used = /** @pureOrBreakMyCode */(alert());\n"+
			"Foo.prototype.unused = /** @pureOrBreakMyCode */(alert());\n"+
			"function foo() {\n"+
			"  return new Foo().used;\n"+
			"}\n"+
			"foo();\n",
		"/** @constructor */\n"+
			"function Foo() {}\n"+
			"Foo.prototype.used = /** @pureOrBreakMyCode */(alert());\n"+
			"function foo() {\n"+
			"  return new Foo().used;\n"+
			"}\n"+
			"foo();\n")
}
