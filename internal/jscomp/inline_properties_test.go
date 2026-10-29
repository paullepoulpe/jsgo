// Port of com/google/javascript/jscomp/InlinePropertiesTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/InlinePropertiesTest.java)

/*
 * Copyright 2012 The Closure Compiler Authors.
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

type InlinePropertiesTest struct {
	CompilerTestCaseBase

	runSmartNameRemoval bool
}

var InlinePropertiesTest_EXTERNS = CompilerTestCase_MINIMAL_EXTERNS +
	"Function.prototype.call=function(){};\n" +
	"Function.prototype.inherits=function(){};\n" +
	"prop.toString;\n" +
	"var google = { gears: { factory: {}, workerPool: {} } };\n" +
	"/** @type {?} */ var externUnknownVar;\n" +
	"/** @type {!Function} */ var externFn;\n"

// Java: InlinePropertiesTest()
func NewInlinePropertiesTest() *InlinePropertiesTest {
	c := &InlinePropertiesTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, InlinePropertiesTest_EXTERNS)
	c.runSmartNameRemoval = false
	return c
}

func TestInlinePropertiesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewInlinePropertiesTest() })
}

// Java: InlinePropertiesTest.getOptions()
func (c *InlinePropertiesTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	// Ignore a few type warnings: we intentionally trigger these warnings
	// to make sure that the pass still operates correctly with bad code.
	ignored := NewDiagnosticGroup(TypeCheck_INEXISTENT_PROPERTY, TypeValidator_TYPE_MISMATCH_WARNING)
	options.SetWarningLevel(ignored, CheckLevel_OFF)
	return options
}

// Java: InlinePropertiesTest.getProcessor(Compiler)
func (c *InlinePropertiesTest) GetProcessor(compiler *Compiler) CompilerPass {
	var pass CompilerPass = NewInlineProperties(compiler)
	if c.runSmartNameRemoval {
		removalPass := NewRemoveUnusedCode(compiler, RemoveUnusedCodeOptions{
			RemoveLocalVars:                               true,
			RemoveGlobals:                                 true,
			PreserveFunctionExpressionNames:               true,
			RemoveUnusedPrototypeProperties:               true,
			RemoveUnusedThisProperties:                    true,
			RemoveUnusedObjectDefinePropertiesDefinitions: true,
		})
		return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
			removalPass.Process(externs, root)
			pass.Process(externs, root)
		})
	}
	return pass
}

// Java: InlinePropertiesTest.setUp()
func (c *InlinePropertiesTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableTypeCheck()
	c.EnableNormalize()
	c.EnableClosurePass()
	c.EnableRewriteClosureProvides()
	c.EnableGatherExternProperties()
	c.runSmartNameRemoval = false
	c.ReplaceTypesWithColors()
	c.DisableCompareJsDoc()
}

// Java: InlinePropertiesTest.testConstInstanceProp1()
func (c *InlinePropertiesTest) TestConstInstanceProp1() {
	// Replace a reference to known constant property.
	c.Test(
		"/** @constructor */\n"+
			"function C() {\n"+
			"  this.foo = 1;\n"+
			"}\n"+
			"new C().foo;\n",
		"/** @constructor */\n"+
			"function C() {\n"+
			"  this.foo = 1;\n"+
			"}\n"+
			"new C(), 1;\n")

	c.Test(
		"/** @constructor */\n"+
			"function C() {\n"+
			"  {\n"+
			"    this.foo = 1;\n"+
			"  }\n"+
			"}\n"+
			"new C().foo;\n",
		"/** @constructor */\n"+
			"function C() {\n"+
			"  {\n"+
			"    this.foo = 1;\n"+
			"  }\n"+
			"}\n"+
			"new C(), 1;\n")
}

// Java: InlinePropertiesTest.testConstInstanceProp2()
func (c *InlinePropertiesTest) TestConstInstanceProp2() {
	// Replace a constant reference
	c.Test(
		"/** @constructor */\n"+
			"function C() {\n"+
			"  this.foo = 1;\n"+
			"}\n"+
			"var x = new C();\n"+
			"x.foo;\n",
		"/** @constructor */\n"+
			"function C() {\n"+
			"  this.foo = 1\n"+
			"}\n"+
			"var x = new C();\n"+
			"1;\n")
}

// Java: InlinePropertiesTest.testConstInstanceProp3()
func (c *InlinePropertiesTest) TestConstInstanceProp3() {
	// Replace a constant reference
	c.Test(
		"/** @constructor */\n"+
			"function C() {\n"+
			"  this.foo = 1;\n"+
			"}\n"+
			"/** @type {C} */\n"+
			"var x = new C();\n"+
			"x.foo;\n",
		"/** @constructor */\n"+
			"function C() {\n"+
			"  this.foo = 1\n"+
			"}\n"+
			"/** @type {C} */\n"+
			"var x = new C();\n"+
			"1;\n")
}

// Java: InlinePropertiesTest.testConstInstanceProp4()
func (c *InlinePropertiesTest) TestConstInstanceProp4() {
	// This pass replies on DisambiguateProperties to distinguish like named
	// properties so it doesn't handle this case.
	c.TestSame(
		"/** @constructor */\n" +
			"function C() {\n" +
			"  this.foo = 1;\n" +
			"}\n" +
			"/** @constructor */\n" +
			"function B() {\n" +
			"  this.foo = 1;\n" +
			"}\n" +
			"new C().foo;\n")
}

// Java: InlinePropertiesTest.testConstInstanceProp5()
func (c *InlinePropertiesTest) TestConstInstanceProp5() {
	c.Test(
		"/** @constructor */\n"+
			"function Foo() {\n"+
			"  /** @type {?number} */\n"+
			"  this.a = 1;\n"+
			"  /** @type {number} */\n"+
			"  this.b = 2;\n"+
			"}\n"+
			"var x = (new Foo).b;\n",
		"/** @constructor */\n"+
			"function Foo() {\n"+
			"  /** @type {?number} */\n"+
			"  this.a = 1;\n"+
			"  /** @type {number} */\n"+
			"  this.b = 2;\n"+
			"}\n"+
			"var x = (new Foo, 2);\n")
}

// Java: InlinePropertiesTest.testConstClassProps1()
func (c *InlinePropertiesTest) TestConstClassProps1() {
	// Inline constant class properties,
	c.Test(
		"/** @constructor */\n"+
			"function C() {\n"+
			"}\n"+
			"C.bar = 2;\n"+
			"C.foo = 1;\n"+
			"var z = C.foo;\n",
		"/** @constructor */\n"+
			"function C() {\n"+
			"}\n"+
			"C.bar = 2;\n"+
			"C.foo = 1;\n"+
			"var z = 1;\n")
}

// Java: InlinePropertiesTest.testConstClassProps2()
func (c *InlinePropertiesTest) TestConstClassProps2() {
	// Don't confuse, class properties with instance properties
	c.TestSame(
		"/** @constructor */\n" +
			"function C() {\n" +
			"  this.foo = 1;\n" +
			"}\n" +
			"var z = C.foo;\n")
}

// Java: InlinePropertiesTest.testConstClassProps3()
func (c *InlinePropertiesTest) TestConstClassProps3() {
	// Don't confuse, class properties with prototype properties
	c.TestSame(
		"/** @constructor */\n" +
			"function C() {}\n" +
			"C.prototype.foo = 1;\n" +
			"var z = C.foo;\n")
}

// Java: InlinePropertiesTest.testConstClassProps4()
func (c *InlinePropertiesTest) TestConstClassProps4() {
	// Don't confuse unique constructors with similiar function types
	c.TestSame(
		"/** @constructor */\n" +
			"function C() {}\n" +
			"/** @constructor @extends {C} */\n" +
			"function D() {}\n" +
			"/** @type {function(new:C): undefined} */\n" +
			"var x = D;\n" +
			"/** @type {number} */ x.foo = 1;\n" +
			"var z = C.foo;\n")
}

// Java: InlinePropertiesTest.testConstClassProps5()
func (c *InlinePropertiesTest) TestConstClassProps5() {
	// Don't confuse subtype constructors properties
	c.TestSame(
		"/** @constructor */\n" +
			"function C() {}\n" +
			"/** @constructor @extends {C} */\n" +
			"function D() {}\n" +
			"D.foo = 1;\n" +
			"var z = C.foo;\n")
}

// Java: InlinePropertiesTest.testConstClassProps6()
func (c *InlinePropertiesTest) TestConstClassProps6() {
	// Don't inline to unknowns
	c.TestSame(
		"/** @constructor */\n" +
			"function C() {}\n" +
			"C.foo = 1;\n" +
			"var z = externUnknownVar.foo;\n")
}

// Java: InlinePropertiesTest.testConstClassProps7()
func (c *InlinePropertiesTest) TestConstClassProps7() {
	// Don't inline to Function prop
	c.TestSame(
		"/** @constructor */\n" +
			"function C() {}\n" +
			"C.foo = 1;\n" +
			"var z = externFn.foo;\n")
}

// Java: InlinePropertiesTest.testNonConstClassProp1()
func (c *InlinePropertiesTest) TestNonConstClassProp1() {
	c.TestSame(
		"/** @constructor */\n" +
			"function C() {\n" +
			"}\n" +
			"C.foo = 1;\n" +
			"alert(C.foo);\n" +
			"delete C.foo;\n")
}

// Java: InlinePropertiesTest.testNonConstClassProp2()
func (c *InlinePropertiesTest) TestNonConstClassProp2() {
	c.TestSame(
		"/** @constructor */\n" +
			"function C() {\n" +
			"}\n" +
			"C.foo = 1;\n" +
			"alert(C.foo);\n" +
			"C.foo = 2;\n")
}

// Java: InlinePropertiesTest.testNonConstClassProp3()
func (c *InlinePropertiesTest) TestNonConstClassProp3() {
	c.TestSame(
		"/** @constructor */\n" +
			"function C() {\n" +
			"}\n" +
			"C.foo = 1;\n" +
			"function f(a) {\n" +
			" a.foo = 2;\n" +
			"}\n" +
			"alert(C.foo);\n" +
			"f(C);\n")
}

// Java: InlinePropertiesTest.testNonConstInstanceProp1()
func (c *InlinePropertiesTest) TestNonConstInstanceProp1() {
	c.TestSame(
		"/** @constructor */\n" +
			"function C() {\n" +
			"  this.foo = 1;\n" +
			"}\n" +
			"var x = new C();\n" +
			"alert(x.foo);\n" +
			"delete x.foo;\n")
}

// Java: InlinePropertiesTest.testNonConstInstanceProp2()
func (c *InlinePropertiesTest) TestNonConstInstanceProp2() {
	c.TestSame(
		"/** @constructor */\n" +
			"function C() {\n" +
			"  this.foo = 1;\n" +
			"}\n" +
			"var x = new C();\n" +
			"alert(x.foo);\n" +
			"x.foo = 2;\n")
}

// Java: InlinePropertiesTest.testNonConstructorInstanceProp1()
func (c *InlinePropertiesTest) TestNonConstructorInstanceProp1() {
	c.TestSame(
		"function C() {\n" +
			"  this.foo = 1;\n" +
			"  return this;\n" +
			"}\n" +
			"C().foo;\n")
}

// Java: InlinePropertiesTest.testConditionalInstanceProp1()
func (c *InlinePropertiesTest) TestConditionalInstanceProp1() {
	c.TestSame(
		"/** @constructor */\n" +
			"function C() {\n" +
			"  if (false) this.foo = 1;\n" +
			"}\n" +
			"new C().foo;\n")
}

// Java: InlinePropertiesTest.testConstPrototypeProp1()
func (c *InlinePropertiesTest) TestConstPrototypeProp1() {
	c.Test(
		"/** @constructor */\n"+
			"function C() {}\n"+
			"C.prototype.foo = 1;\n"+
			"new C().foo;\n",
		"/** @constructor */\n"+
			"function C() {}\n"+
			"C.prototype.foo = 1;\n"+
			"new C(), 1;\n")
}

// Java: InlinePropertiesTest.testConstPrototypeProp2()
func (c *InlinePropertiesTest) TestConstPrototypeProp2() {
	c.Test(
		"/** @constructor */\n"+
			"function C() {}\n"+
			"C.prototype.foo = 1;\n"+
			"var x = new C();\n"+
			"x.foo;\n",
		"/** @constructor */\n"+
			"function C() {}\n"+
			"C.prototype.foo = 1;\n"+
			"var x = new C();\n"+
			"1;\n")
}

// Java: InlinePropertiesTest.testConstPrototypePropInGlobalBlockScope()
func (c *InlinePropertiesTest) TestConstPrototypePropInGlobalBlockScope() {
	c.Test(
		"/** @constructor */\n"+
			"function C() {}\n"+
			"{\n"+
			"  C.prototype.foo = 1;\n"+
			"}\n"+
			"var x = new C();\n"+
			"x.foo;\n",
		"/** @constructor */\n"+
			"function C() {}\n"+
			"{\n"+
			"  C.prototype.foo = 1;\n"+
			"}\n"+
			"var x = new C();\n"+
			"1;\n")
}

// Java: InlinePropertiesTest.testGlobalThisNotInlined()
func (c *InlinePropertiesTest) TestGlobalThisNotInlined() {
	c.TestSame(
		"this.foo = 1;\n" +
			"/** @constructor */\n" +
			"function C() {\n" +
			"  foo;\n" +
			"}\n")
}

// Java: InlinePropertiesTest.testConstPrototypePropFromSuper()
func (c *InlinePropertiesTest) TestConstPrototypePropFromSuper() {
	c.Test(
		"/** @constructor */\n"+
			"function C() {}\n"+
			"C.prototype.foo = 1;\n"+
			"/** @constructor @extends {C} */\n"+
			"function D() {}\n"+
			"(new D).foo;\n",
		"/** @constructor */\n"+
			"function C() {}\n"+
			"C.prototype.foo = 1;\n"+
			"/** @constructor @extends {C} */\n"+
			"function D() {}\n"+
			"new D, 1;\n")
}

// Java: InlinePropertiesTest.testTypedPropInlining()
func (c *InlinePropertiesTest) TestTypedPropInlining() {
	c.Test(
		"/** @constructor */\n"+
			"function C() {}\n"+
			"C.prototype.foo = 1;\n"+
			"function f(/** !C */ x) { return x.foo; }\n"+
			"f(new C);\n",
		"/** @constructor */\n"+
			"function C() {}\n"+
			"C.prototype.foo = 1;\n"+
			"function f(/** !C */ x) { return 1; }\n"+
			"f(new C);\n")
}

// Java: InlinePropertiesTest.testTypeMismatchNoPropInlining()
func (c *InlinePropertiesTest) TestTypeMismatchNoPropInlining() {
	c.TestSame(
		"/** @constructor */\n" +
			"function C() {}\n" +
			"C.prototype.foo = 1;\n" +
			"function f(/** !C */ x) { return x.foo; }\n" +
			"f([]);\n")
}

// Java: InlinePropertiesTest.testTypeMismatchForPrototypeNoPropInlining()
func (c *InlinePropertiesTest) TestTypeMismatchForPrototypeNoPropInlining() {
	c.TestSame(
		"/** @constructor */\n" +
			"function C() {}\n" +
			"C.prototype.foo = 1;\n" +
			"function f(/** number */ x) {}\n" +
			"f(C.prototype);\n" +
			"new C().foo;\n")
}

// Java: InlinePropertiesTest.testStructuralInterfacesNoPropInlining()
func (c *InlinePropertiesTest) TestStructuralInterfacesNoPropInlining() {
	c.TestSame(
		"/** @record */ function I() {}\n" +
			"/** @type {number|undefined} */ I.prototype.foo;\n" +
			"\n" +
			"/** @constructor @implements {I} */\n" +
			"function C() {}\n" +
			"/** @override */\n" +
			"C.prototype.foo = 1;\n" +
			"\n" +
			"function f(/** !I */ x) { return x.foo; }\n" +
			"f([]);\n")
}

// Java: InlinePropertiesTest.testStructuralInterfacesNoPropInlining2()
func (c *InlinePropertiesTest) TestStructuralInterfacesNoPropInlining2() {
	c.runSmartNameRemoval = true

	c.Test(
		"/** @record */\n"+
			"function I() {\n"+
			"  /** @type {number} */ this.foo;\n"+
			"}\n"+
			"\n"+
			"/** @constructor @implements {I} */\n"+
			"function C() { /** @type {number} */ this.foo = 1; }\n"+
			"\n"+
			"function f(/** ? */ x) { return x.foo; }\n"+
			"f(new C());\n",
		"/** @constructor @implements {I} */\n"+
			"function C() { /** @type {number} */ this.foo = 1; }\n"+
			"\n"+
			"function f(/** ? */ x) { return x.foo; }\n"+
			"f(new C());\n")
}

// Java: InlinePropertiesTest.testConstInstanceProp_es6Class()
func (c *InlinePropertiesTest) TestConstInstanceProp_es6Class() {
	// Replace a reference to known constant property.
	c.Test(
		"class C {\n"+
			"  constructor() {\n"+
			"    this.foo = 1;\n"+
			"  }\n"+
			"}\n"+
			"new C().foo;\n",
		"class C {\n"+
			"  constructor() {\n"+
			"    this.foo = 1;\n"+
			"  }\n"+
			"}\n"+
			"new C(), 1;\n")
}

// Java: InlinePropertiesTest.testMultipleConstInstanceProp_es6Class()
func (c *InlinePropertiesTest) TestMultipleConstInstanceProp_es6Class() {
	c.Test(
		"class Foo {\n"+
			"  constructor() {\n"+
			"    /** @type {?number} */\n"+
			"    this.a = 1;\n"+
			"    /** @type {number} */\n"+
			"    this.b = 2;\n"+
			"  }\n"+
			"}\n"+
			"var x = (new Foo).b;\n",
		"class Foo {\n"+
			"  constructor() {\n"+
			"    /** @type {?number} */\n"+
			"    this.a = 1;\n"+
			"    /** @type {number} */\n"+
			"    this.b = 2;\n"+
			"  }\n"+
			"}\n"+
			"var x = (new Foo, 2);\n")
}

// Java: InlinePropertiesTest.testConstInstancePropInArrowFunction_es6Class()
func (c *InlinePropertiesTest) TestConstInstancePropInArrowFunction_es6Class() {
	// Don't replace a reference to known constant property defined in an arrow function.
	c.TestSame(
		"/** @unrestricted */ // make this not a struct, so we can define this.foo\n" +
			"class C {\n" +
			"  constructor() {\n" +
			"    (() => {\n" +
			"      this.foo = 1;\n" +
			"    })();\n" +
			"  }\n" +
			"}\n" +
			"new C().foo;\n")
}

// Java: InlinePropertiesTest.testConstInstancePropInForOf_es6Class()
func (c *InlinePropertiesTest) TestConstInstancePropInForOf_es6Class() {
	c.TestSame(
		`/** @unrestricted */
class C {
  constructor() {
    for (const x of [1, 2]) {
      this.foo = 1;
    }
  }
}
new C().foo;
`)
}

// Java: InlinePropertiesTest.testConstClassProps_es6Class()
func (c *InlinePropertiesTest) TestConstClassProps_es6Class() {
	// Inline constant class properties,
	c.Test(
		"class C {}\n"+
			"C.bar = 2;\n"+
			"C.foo = 1;\n"+
			"var z = C.foo;\n",
		"class C {}\n"+
			"C.bar = 2;\n"+
			"C.foo = 1;\n"+
			"var z = 1;\n")
}

// Java: InlinePropertiesTest.testConstClassPropsInheritedProp_es6Class()
func (c *InlinePropertiesTest) TestConstClassPropsInheritedProp_es6Class() {
	c.Test(
		"class C {}\n"+
			"class D extends C {}\n"+
			"C.foo = 1;\n"+
			"var z = D.foo;\n",
		"class C {}\n"+
			"class D extends C {}\n"+
			"C.foo = 1;\n"+
			"var z = 1;\n")
}

// Java: InlinePropertiesTest.testConstClassPropsInheritedPropChain_es6Class()
func (c *InlinePropertiesTest) TestConstClassPropsInheritedPropChain_es6Class() {
	c.Test(
		"class C {}\n"+
			"class D extends C {}\n"+
			"class E extends D {}\n"+
			"class F extends E {}\n"+
			"C.foo = 1;\n"+
			"var z = F.foo;\n",
		"class C {}\n"+
			"class D extends C {}\n"+
			"class E extends D {}\n"+
			"class F extends E {}\n"+
			"C.foo = 1;\n"+
			"var z = 1;\n")
}

// Java: InlinePropertiesTest.testConstClassPropsNonInheritedProp_es6Class()
func (c *InlinePropertiesTest) TestConstClassPropsNonInheritedProp_es6Class() {
	// Test that we don't accidentally treat the superclass as having a subclass prop
	c.TestSame(
		"class C {}\n" +
			"class D extends C {}\n" +
			"D.foo = 1;\n" +
			"var z = C.foo;\n")
}

// Java: InlinePropertiesTest.testNonConstClassProp_es6ClassWithStaticMethod()
func (c *InlinePropertiesTest) TestNonConstClassProp_es6ClassWithStaticMethod() {
	c.TestSame(
		"class C { static foo() {} }\n" +
			"alert(C.foo);\n" +
			"C.foo = 1;\n")
}

// Java: InlinePropertiesTest.testConstPrototypeProp_es6Class()
func (c *InlinePropertiesTest) TestConstPrototypeProp_es6Class() {
	c.Test(
		"class C {}\n"+
			"C.prototype.foo = 1;\n"+
			"new C().foo;\n",
		"class C {}\n"+
			"C.prototype.foo = 1;\n"+
			"new C(), 1;\n")
}

// Java: InlinePropertiesTest.testNonConstPrototypePropFromMemberFn()
func (c *InlinePropertiesTest) TestNonConstPrototypePropFromMemberFn() {
	c.TestSame(
		"class C {\n" +
			"  foo() {}\n" +
			"}\n" +
			"C.prototype.foo = 4;\n" +
			"(new C()).foo\n")
}

// Java: InlinePropertiesTest.testObjectPatternStringKeyDoesntInvalidateProp()
func (c *InlinePropertiesTest) TestObjectPatternStringKeyDoesntInvalidateProp() {
	c.Test(
		"/** @constructor */\n"+
			"function C() {\n"+
			"  this.foo = 3;\n"+
			"}\n"+
			"(new C()).foo\n"+
			"const {foo} = new C();\n",
		"/** @constructor */\n"+
			"function C() {\n"+
			"  this.foo = 3;\n"+
			"}\n"+
			"new C(), 3;\n"+
			"const {foo} = new C();\n")
}

// Java: InlinePropertiesTest.testNoInlineOnRecordType()
func (c *InlinePropertiesTest) TestNoInlineOnRecordType() {
	c.TestSame(
		"/** @record */\n" +
			"class C {}\n" +
			"C.bar = 2;\n" +
			"C.foo = 1;\n" +
			"var z = C.foo;\n")
}

// Java: InlinePropertiesTest.testNoInlineOnInterfaceType()
func (c *InlinePropertiesTest) TestNoInlineOnInterfaceType() {
	c.TestSame(
		"/** @interface */\n" +
			"class C {}\n" +
			"C.bar = 2;\n" +
			"C.foo = 1;\n" +
			"var z = C.foo;\n")
}

// Java: InlinePropertiesTest.testClassField()
func (c *InlinePropertiesTest) TestClassField() {
	c.Test(
		"class C {\n"+
			"  a = 1;\n"+
			"  b;\n"+
			"}\n"+
			"(new C()).a;\n"+
			"(new C()).b\n",
		"class C {\n"+
			"  a = 1;\n"+
			"  b;\n"+
			"}\n"+
			"new C(), 1;\n"+
			"(new C()).b\n")
}

// Java: InlinePropertiesTest.testClassFieldWithInheritance()
func (c *InlinePropertiesTest) TestClassFieldWithInheritance() {
	c.Test(
		"class C {\n"+
			"  a = 1;\n"+
			"  b;\n"+
			"}\n"+
			"class D extends C {};\n"+
			"(new D()).a;\n"+
			"(new D()).b\n",
		"class C {\n"+
			"  a = 1;\n"+
			"  b;\n"+
			"}\n"+
			"class D extends C {};\n"+
			"new D(), 1;\n"+
			"(new D()).b;\n")
}

// Java: InlinePropertiesTest.testClassField_static()
func (c *InlinePropertiesTest) TestClassField_static() {
	c.Test(
		"class C {\n"+
			"  static a = 1;\n"+
			"  static b;\n"+
			"}\n"+
			"C.a;\n"+
			"C.b\n",
		"class C {\n"+
			"  static a = 1;\n"+
			"  static b;\n"+
			"}\n"+
			"1;\n"+
			"C.b\n")
}

// Java: InlinePropertiesTest.testClassComputedField()
func (c *InlinePropertiesTest) TestClassComputedField() {
	c.TestSame(
		"/** @dict */\n" +
			"class C {\n" +
			"  ['a'] = 1;\n" +
			"  ['b'];\n" +
			"}\n" +
			"(new C())['a'];\n" +
			"(new C())['b']\n")
}

// Java: InlinePropertiesTest.testClassComputedField_static()
func (c *InlinePropertiesTest) TestClassComputedField_static() {
	c.TestSame(
		"/** @dict */\n" +
			"class C {\n" +
			"  static ['a'] = 1;\n" +
			"  static ['b'];\n" +
			"}\n" +
			"C['a']\n" +
			"C['b']\n")
}

// Java: InlinePropertiesTest.testUninitializedSubclassFieldPreventsInlining()
func (c *InlinePropertiesTest) TestUninitializedSubclassFieldPreventsInlining() {
	c.TestSame(
		`class Guard {
  constructor() { this.escaped = true; }
}
class RawGuard extends Guard {
  escaped;
}
(new RawGuard()).escaped;
`)
}

// Java: InlinePropertiesTest.testInitializedSubclassFieldPreventsInlining()
func (c *InlinePropertiesTest) TestInitializedSubclassFieldPreventsInlining() {
	c.TestSame(
		`class Guard {
  constructor() { this.escaped = true; }
}
class RawGuard extends Guard {
  escaped = false;
}
(new RawGuard()).escaped;
`)
}

// Java: InlinePropertiesTest.testFieldInitializerReadsLaterField()
func (c *InlinePropertiesTest) TestFieldInitializerReadsLaterField() {
	// TODO(b/571953100): inlines property read in field initializer before field is initialized
	c.Test(
		`class C {
  y = this.x;
  x = 1;
}
use((new C()).y);
`,
		`class C {
  y = 1;
  x = 1;
}
use((new C()).y);
`)
}

// Java: InlinePropertiesTest.testFieldInitializerAssigningThisPropInsideOuterConstructor()
func (c *InlinePropertiesTest) TestFieldInitializerAssigningThisPropInsideOuterConstructor() {
	// TODO(b/571953100): inlines this.x in inner class field initializer from outer constructor
	c.Test(
		`class Outer {
  constructor() {
    /** @unrestricted */
    class Inner {
      f = (use(this.x), this.x = 1);
    }
    new Inner();
  }
}
new Outer();
`,
		`class Outer {
  constructor() {
    class Inner {
      f = (use(1), this.x = 1);
    }
    new Inner();
  }
}
new Outer();
`)
}

// Java: InlinePropertiesTest.testPropertyNamedByRenameFunction()
func (c *InlinePropertiesTest) TestPropertyNamedByRenameFunction() {
	// `JSCompiler_renameProperty('foo', C)` names `foo` for reflective access, so the unknown
	// function may write `C.prototype.foo` and `foo` must not be inlined.
	c.TestSame(
		`/** @constructor */
function C() {}
C.prototype.foo = 1;
externUnknownVar(C.prototype, JSCompiler_renameProperty('foo', C), 2);
new C().foo;
`)
}

// Java: InlinePropertiesTest.testPropertyWrittenThroughRenamedName()
func (c *InlinePropertiesTest) TestPropertyWrittenThroughRenamedName() {
	c.TestSame(
		`/** @constructor */
function C() {}
C.prototype.foo = 1;
C.prototype[JSCompiler_renameProperty('foo', C)] = 2;
new C().foo;
`)
}
