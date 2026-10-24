// Port of com/google/javascript/jscomp/RemoveUnusedCodeClassPropertiesTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/RemoveUnusedCodeClassPropertiesTest.java)

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
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

type RemoveUnusedCodeClassPropertiesTest struct {
	CompilerTestCaseBase
	keepLocals  bool
	keepGlobals bool
}

// Java: RemoveUnusedCodeClassPropertiesTest()
func NewRemoveUnusedCodeClassPropertiesTest() *RemoveUnusedCodeClassPropertiesTest {
	c := &RemoveUnusedCodeClassPropertiesTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, RemoveUnusedCodeClassPropertiesTest_EXTERNS)
	c.keepLocals = true
	c.keepGlobals = true
	return c
}

func TestRemoveUnusedCodeClassPropertiesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewRemoveUnusedCodeClassPropertiesTest() })
}

const RemoveUnusedCodeClassPropertiesTest_EXTERNS = "/**\n" +
	" * @constructor\n" +
	" * @param {*=} opt_value\n" +
	" * @return {!Object}\n" +
	" */\n" +
	"function Object(opt_value) {}\n" +
	"/**\n" +
	" * @constructor\n" +
	" * @param {...*} var_args\n" +
	" */\n" +
	"function Function(var_args) {}\n" +
	"/**\n" +
	" * @constructor\n" +
	" * @param {*=} arg\n" +
	" * @return {string}\n" +
	" */\n" +
	"function String(arg) {}\n" +
	"/**\n" +
	" * @record\n" +
	" * @template VALUE\n" +
	" */\n" +
	"/**\n" +
	" * @template T\n" +
	" * @constructor\n" +
	" * @param {...*} var_args\n" +
	" * @return {!Array<?>}\n" +
	" */\n" +
	"function Array(var_args) {}\n" +
	"var window;\n" +
	"function alert(a) {}\n" +
	"function use(x) {}\n" +
	"var EXT = {};\n" +
	"EXT.ext;\n" +
	"var externVar;\n" +
	"function externFunction() {}\n" +
	"/** @type {Function} */\n" +
	"Object.defineProperties = function() {};\n" +
	"/** @type {Function} */\n" +
	"Object.prototype.constructor = function() {};\n" +
	"// NOTE: The following are needed to prevent NTI inexistent property warnings.\n" +
	"var $jscomp = {};\n" +
	"$jscomp.global = {}\n" +
	"/** @type {?} */\n" +
	"$jscomp.global.Object\n" +
	"function JSCompiler_renameProperty(p) {}\n" +
	"var goog = {};\n" +
	"goog.reflect = {};\n" +
	"goog.reflect.object = function(a, b) {};\n"

// Java: RemoveUnusedCodeClassPropertiesTest.getProcessor(Compiler)
func (c *RemoveUnusedCodeClassPropertiesTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewRemoveUnusedCode_Builder(compiler).
		RemoveLocalVars(!c.keepLocals).
		RemoveGlobals(!c.keepGlobals).
		RemoveUnusedPrototypeProperties(true).
		RemoveUnusedThisProperties(true).
		RemoveUnusedObjectDefinePropertiesDefinitions(true).
		Build()
}

// Java: RemoveUnusedCodeClassPropertiesTest.setUp()
func (c *RemoveUnusedCodeClassPropertiesTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	// Allow testing of features that aren't fully supported for output yet.
	c.EnableNormalize()
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.DisableCompareJsDoc()
	c.EnableGatherExternProperties()
	c.DisableTypeCheck()
}

// Java: RemoveUnusedCodeClassPropertiesTest.testSimple1()
func (c *RemoveUnusedCodeClassPropertiesTest) TestSimple1() {
	// A property defined on "this" can be removed
	c.Test("this.a = 2", "")
	c.Test("let x = (this.a = 2)", "let x = 2")
	c.TestSame("this.a = 2; let x = this.a;")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testSimple2()
func (c *RemoveUnusedCodeClassPropertiesTest) TestSimple2() {
	// A property defined on "this" can be removed, even when defined
	// as part of an expression
	c.Test("this.a = 2, alert(1);", "alert(1);")
	c.Test("const x = (this.a = 2, alert(1));", "const x = alert(1);")
	c.Test("const x = (alert(1), this.a = 2);", "const x = (alert(1), 2);")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testSimple3()
func (c *RemoveUnusedCodeClassPropertiesTest) TestSimple3() {
	// A property defined on an object other than "this" can not be removed.
	c.TestSame("var y = {}; y.a = 2")
	// and prevents the removal of the definition on 'this'.
	c.TestSame("var y = {}; y.a = 2; this.a = 2")
	// Some use of the property "a" prevents the removal.
	c.TestSame("var x; var y = {}; y.a = 2; this.a = 1; alert(x.a)")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testObjLit()
func (c *RemoveUnusedCodeClassPropertiesTest) TestObjLit() {
	// A property defined on an object other than "this" can not be removed.
	c.TestSame("({a:2})")
	// and prevent the removal of the definition on 'this'.
	c.TestSame("({a:0}); this.a = 1;")
	// ... even if it's quoted
	c.TestSame("({'a':0}); this.a = 1;")
	// Some use of the property "a" prevents the removal.
	c.TestSame("var x = ({a:0}); this.a = 1; alert(x.a)")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testExtern()
func (c *RemoveUnusedCodeClassPropertiesTest) TestExtern() {
	// A property defined in the externs is can not be removed.
	c.TestSame("this.ext = 2")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testExport()
func (c *RemoveUnusedCodeClassPropertiesTest) TestExport() {
	// An exported property can not be removed.
	c.TestSame("this.ext = 2; window['export'] = this.ext;")
	c.TestSame("function f() { this.ext = 2; } window['export'] = this.ext;")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testAssignOp1()
func (c *RemoveUnusedCodeClassPropertiesTest) TestAssignOp1() {
	// Properties defined using a compound assignment can be removed if the
	// result of the assignment expression is not immediately used.
	c.Test("this.x += 2", "")
	c.TestSame("const x = (this.x += 2)")
	c.TestSame("this.x += 2; const x = this.x;")
	// But, of course, a later use prevents its removal.
	c.TestSame("this.x += 2; let x = {}; x.x;")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testAssignOp2()
func (c *RemoveUnusedCodeClassPropertiesTest) TestAssignOp2() {
	// Properties defined using a compound assignment can be removed if the
	// result of the assignment expression is not immediately used.
	c.Test("this.a += 2, alert(1)", "alert(1)")
	c.Test("const x = (this.a += 2, alert(1))", "const x = alert(1)")
	c.TestSame("const x = (alert(1), this.a += 2)")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testInc1()
func (c *RemoveUnusedCodeClassPropertiesTest) TestInc1() {
	// Increments and Decrements are handled similarly to compound assignments
	// but need a placeholder value when replaced.
	c.Test("this.x++", "")
	c.TestSame("let x = (this.x++)")
	c.TestSame("this.x++; let x = this.x;")

	c.Test("--this.x", "")
	c.TestSame("let x = (--this.x)")
	c.TestSame("--this.x; let x = this.x;")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testInc2()
func (c *RemoveUnusedCodeClassPropertiesTest) TestInc2() {
	// Increments and Decrements are handled similarly to compound assignments
	// but need a placeholder value when replaced.
	c.Test("this.a++, alert()", "alert()")
	c.Test("let x = (this.a++, alert())", "let x = alert()")
	c.TestSame("let x = (alert(), this.a++)")

	c.Test("--this.a, alert()", "alert()")
	c.Test("let x = (--this.a, alert())", "let x = alert()")
	c.TestSame("let x = (alert(), --this.a)")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testDestructuringRest()
func (c *RemoveUnusedCodeClassPropertiesTest) TestDestructuringRest() {
	c.TestSame(
		"function Foo() {}\n" +
			"Foo.a = function() {};\n" +
			"({ ...Foo.a.b } = 0);\n")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testExprResult()
func (c *RemoveUnusedCodeClassPropertiesTest) TestExprResult() {
	c.Test("this.x", "")
	c.Test("externFunction().prototype.x", "externFunction()")
	// It doesn't make much sense to use optional chaining in these cases, but if you do,
	// it shouldn't prevent unused property removal
	c.Test("this?.x", "")
	c.Test("externFunction()?.prototype.x", "externFunction()")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testJSCompiler_renameProperty()
func (c *RemoveUnusedCodeClassPropertiesTest) TestJSCompiler_renameProperty() {
	// JSCompiler_renameProperty introduces a use of the property
	c.TestSame("var x; this.a = 2; x[JSCompiler_renameProperty('a')]")
	c.TestSame("this.a = 2; JSCompiler_renameProperty('a')")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testForIn()
func (c *RemoveUnusedCodeClassPropertiesTest) TestForIn() {
	// This is the basic assumption that this pass makes:
	// it can remove properties even when the object is used in a FOR-IN loop
	c.Test(
		"let x = {}; this.y = 1;for (var a in x) { alert(x[a]) }",
		"let x = {};            for (var a in x) { alert(x[a]) }")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testObjectKeys()
func (c *RemoveUnusedCodeClassPropertiesTest) TestObjectKeys() {
	// This is the basic assumption that this pass makes:
	// it can remove properties even when the object are referenced
	c.Test(
		"this.y = 1;alert(Object.keys(this))", // preserve format
		"           alert(Object.keys(this))")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testObjectReflection1()
func (c *RemoveUnusedCodeClassPropertiesTest) TestObjectReflection1() {
	// Verify reflection prevents removal.
	c.TestSame(
		"/** @constructor */ // preserve newlines\n" +
			"function A() { this.foo = 1; }\n" +
			"use(goog.reflect.object(A, {foo: 'foo'}));\n")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testObjectReflection2()
func (c *RemoveUnusedCodeClassPropertiesTest) TestObjectReflection2() {
	// Any object literal definition prevents removal.
	// Type based removal would allow this to be removed.
	c.TestSame(
		"/** @constructor */ // preserve newlines\n" +
			"function A() {this.foo = 1;}\n" +
			"use({foo: 'foo'});\n")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testIssue730()
func (c *RemoveUnusedCodeClassPropertiesTest) TestIssue730() {
	// Partial removal of properties can causes problems if the object is
	// sealed.
	c.TestSame(
		"function A() {this.foo = 0;}\n" +
			"function B() {this.a = new A();}\n" +
			"B.prototype.dostuff = function() { this.a.foo++; alert('hi'); }\n" +
			"new B().dostuff();\n")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testPrototypeProps1()
func (c *RemoveUnusedCodeClassPropertiesTest) TestPrototypeProps1() {
	c.Test(
		"function A() {this.foo = 1;}\n"+
			"A.prototype.foo = 0;\n"+
			"A.prototype.method = function() {this.foo++};\n"+
			"new A().method()\n",
		"function A() {             }\n"+
			"\n"+
			"A.prototype.method = function() {          };\n"+
			"new A().method()\n")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testPrototypeProps2()
func (c *RemoveUnusedCodeClassPropertiesTest) TestPrototypeProps2() {
	// don't remove properties that are exported by convention
	c.TestSame(
		"function A() {this._foo = 1;}\n" +
			"A.prototype._foo = 0;\n" +
			"A.prototype.method = function() {this._foo++};\n" +
			"new A().method()\n")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testConstructorProperty1()
func (c *RemoveUnusedCodeClassPropertiesTest) TestConstructorProperty1() {
	c.EnableTypeCheck()

	c.Test(
		"/** @constructor */ function C() {} C.prop = 1;",
		"/** @constructor */ function C() {}            ")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testConstructorProperty2()
func (c *RemoveUnusedCodeClassPropertiesTest) TestConstructorProperty2() {
	c.EnableTypeCheck()

	c.TestSame(
		"/** @constructor */ function C() {}\n" +
			"C.prop = 1;\n" +
			"function foo(a) { alert(a.prop) };\n" +
			"foo(C)\n")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testES6StaticProperty()
func (c *RemoveUnusedCodeClassPropertiesTest) TestES6StaticProperty() {
	c.Test(
		"class C { static prop() {} }", // preserve newline
		"class C {                  }")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testES6StaticProperty2()
func (c *RemoveUnusedCodeClassPropertiesTest) TestES6StaticProperty2() {
	c.Test("class C {} C.prop = 1;", "class C {}")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testObjectDefineProperties1()
func (c *RemoveUnusedCodeClassPropertiesTest) TestObjectDefineProperties1() {
	c.EnableTypeCheck()

	c.TestSame(
		"/** @constructor */ function C() {}\n" +
			"Object.defineProperties(C, {prop:{value:1}});\n" +
			"function foo(a) { alert(a.prop) };\n" +
			"foo(C)\n")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testObjectDefineProperties2()
func (c *RemoveUnusedCodeClassPropertiesTest) TestObjectDefineProperties2() {
	c.EnableTypeCheck()

	c.Test(
		"/** @constructor */ function C() {}\n"+
			"Object.defineProperties(C, {prop:{value:1}});\n",
		"/** @constructor */ function C() {}\n"+
			"Object.defineProperties(C, {});\n")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testObjectDefineProperties3()
func (c *RemoveUnusedCodeClassPropertiesTest) TestObjectDefineProperties3() {
	c.EnableTypeCheck()

	c.Test(
		"/** @constructor */ function C() {}\n"+
			"Object.defineProperties(C,\n"+
			"  {prop:{\n"+
			"    get:function(){},\n"+
			"    set:function(a){},\n"+
			"}});\n",
		"/** @constructor */ function C() {}\n"+
			"Object.defineProperties(C, {});\n")
}

// side-effect in definition retains property definition, but doesn't count as a reference
// Java: RemoveUnusedCodeClassPropertiesTest.testObjectDefineProperties4()
func (c *RemoveUnusedCodeClassPropertiesTest) TestObjectDefineProperties4() {
	c.EnableTypeCheck()

	c.Test(
		"/** @constructor */ function C() { this.prop = 3; }\n"+
			"Object.defineProperties(C, {prop:alert('')});\n",
		"/** @constructor */ function C() {                }\n"+
			"Object.defineProperties(C, {prop:alert('')});\n")
}

// quoted properties retains property
// Java: RemoveUnusedCodeClassPropertiesTest.testObjectDefineProperties5()
func (c *RemoveUnusedCodeClassPropertiesTest) TestObjectDefineProperties5() {
	c.EnableTypeCheck()

	c.TestSame(
		"/** @constructor */ function C() {}\n" +
			"Object.defineProperties(C, {'prop': {value: 1}});\n")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testObjectDefineProperties6()
func (c *RemoveUnusedCodeClassPropertiesTest) TestObjectDefineProperties6() {
	c.EnableTypeCheck()

	// an unknown destination object doesn't prevent removal.
	c.Test(
		"Object.defineProperties(externVar(), {prop:{value:1}});",
		"Object.defineProperties(externVar(), {              });")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testObjectDefineProperties7()
func (c *RemoveUnusedCodeClassPropertiesTest) TestObjectDefineProperties7() {
	c.EnableTypeCheck()

	c.Test(
		"/** @constructor */ function C() {}\n"+
			"Object.defineProperties(C, {prop:{get:function () {return new C}}});\n",
		"/** @constructor */ function C() {}\n"+
			"Object.defineProperties(C, {});\n")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testObjectDefineProperties8()
func (c *RemoveUnusedCodeClassPropertiesTest) TestObjectDefineProperties8() {
	c.EnableTypeCheck()

	c.Test(
		"/** @constructor */ function C() {}\n"+
			"Object.defineProperties(C, {prop:{set:function (a) {return alert(a)}}});\n",
		"/** @constructor */ function C() {}\n"+
			"Object.defineProperties(C, {});\n")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testObjectDefinePropertiesQuotesPreventRemoval()
func (c *RemoveUnusedCodeClassPropertiesTest) TestObjectDefinePropertiesQuotesPreventRemoval() {
	c.EnableTypeCheck()

	c.TestSame(
		"/** @constructor */ function C() { this.prop = 1; }\n" +
			"Object.defineProperties(C, {'prop':{set:function (a) {return alert(a.prop)}}});\n")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testObjectDefineProperties_usedSetter_notRemoved()
func (c *RemoveUnusedCodeClassPropertiesTest) TestObjectDefineProperties_usedSetter_notRemoved() {
	c.EnableTypeCheck()

	c.TestSame(
		"/** @constructor */ function C() {}\n" +
			"Object.defineProperties(C, {prop:{set:function (a) {alert(2)}}});\n" +
			"C.prop = 2;\n")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testPrototypeMethodDef_notConsideredSetterUse()
func (c *RemoveUnusedCodeClassPropertiesTest) TestPrototypeMethodDef_notConsideredSetterUse() {
	c.EnableTypeCheck()

	c.Test(
		"/** @constructor */ function C() {}\n"+
			"Object.defineProperties(C, {prop:{set:function (a) {alert(2)}}});\n"+
			"/** @constructor */ function D () {}\n"+
			"D.prototype.prop = function() {};\n",
		"/** @constructor */ function C() {}\n"+
			"Object.defineProperties(C, {});\n"+
			"/** @constructor */ function D () {}\n")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testEs6GettersWithoutTranspilation()
func (c *RemoveUnusedCodeClassPropertiesTest) TestEs6GettersWithoutTranspilation() {
	c.Test(
		"class C { get value() { return 0; } }", // preserve newline
		"class C {                           }")
	c.TestSame("class C { get value() { return 0; } } const x = (new C()).value")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testES6ClassComputedProperty()
func (c *RemoveUnusedCodeClassPropertiesTest) TestES6ClassComputedProperty() {
	c.TestSame("class C { ['test' + 3]() { return 0; } }")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testEs6SettersWithoutTranspilation()
func (c *RemoveUnusedCodeClassPropertiesTest) TestEs6SettersWithoutTranspilation() {
	c.Test(
		"class C { set value(val) { this.internalVal = val; } }", // preserve newline
		"class C {                                            }")

	c.Test(
		"class C { set value(val) { this.internalVal = val; } } (new C()).value = 3;",
		"class C { set value(val) {                         } } (new C()).value = 3;")
	c.TestSame(
		"class C {\n" +
			"  set value(val) {\n" +
			"    this.internalVal = val;\n" +
			"  }\n" +
			"  get value() {\n" +
			"    return this.internalVal;\n" +
			"  }\n" +
			"}\n" +
			"const y = new C();\n" +
			"y.value = 3;\n" +
			"const x = y.value;\n")
}

// All object literal fields are not removed, but the following
// tests assert that the pass does not fail.
// Java: RemoveUnusedCodeClassPropertiesTest.testEs6EnhancedObjLiteralsComputedValuesNotRemoved()
func (c *RemoveUnusedCodeClassPropertiesTest) TestEs6EnhancedObjLiteralsComputedValuesNotRemoved() {
	c.TestSame(
		"function getCar(make, model, value) {\n" +
			"  return {\n" +
			"    ['make' + make] : true\n" +
			"  };\n" +
			"}\n")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testEs6EnhancedObjLiteralsMethodShortHandNotRemoved()
func (c *RemoveUnusedCodeClassPropertiesTest) TestEs6EnhancedObjLiteralsMethodShortHandNotRemoved() {
	c.TestSame(
		"function getCar(make, model, value) {\n" +
			"  return {\n" +
			"    getModel() {\n" +
			"      return model;\n" +
			"    }\n" +
			"  };\n" +
			"}\n")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testEs6EnhancedObjLiteralsPropertyShorthand()
func (c *RemoveUnusedCodeClassPropertiesTest) TestEs6EnhancedObjLiteralsPropertyShorthand() {
	c.TestSame("function getCar(make, model, value) { return {model}; }")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testTranspiledEs6GettersRemoval()
func (c *RemoveUnusedCodeClassPropertiesTest) TestTranspiledEs6GettersRemoval() {
	c.EnableTypeCheck()
	c.Test(
		// This is the output of ES6->ES5 class getter converter.
		// See Es6TranspilationIntegrationTest.testEs5GettersAndSettersClasses test method.
		"/** @constructor @struct */\n"+
			"var C = function() {};\n"+
			"/** @type {?} */\n"+
			"C.prototype.value = 0;\n"+
			"$jscomp.global.Object.defineProperties(C.prototype, {\n"+
			"  value: {\n"+
			"    configurable: true,\n"+
			"    enumerable: true,\n"+
			"    /** @this {C} */\n"+
			"    get: function() {\n"+
			"      return 0;\n"+
			"    }\n"+
			"  }\n"+
			"});\n",
		"/** @constructor @struct */var C=function(){};\n"+
			"$jscomp.global.Object.defineProperties(C.prototype, {});\n")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testTranspiledEs6SettersRemoval()
func (c *RemoveUnusedCodeClassPropertiesTest) TestTranspiledEs6SettersRemoval() {
	c.EnableTypeCheck()
	c.Test(
		// This is the output of ES6->ES5 class setter converter.
		// See Es6TranspilationIntegrationTest.testEs5GettersAndSettersClasses test method.
		"/** @constructor @struct */\n"+
			"var C = function() {};\n"+
			"/** @type {?} */\n"+
			"C.prototype.value;\n"+
			"/** @type {?} */\n"+
			"C.prototype.internalVal;\n"+
			"$jscomp.global.Object.defineProperties(C.prototype, {\n"+
			"  value: {\n"+
			"    configurable: true,\n"+
			"    enumerable: true,\n"+
			"    /** @this {C} */\n"+
			"    set: function(val) {\n"+
			"      this.internalVal = val;\n"+
			"    }\n"+
			"  }\n"+
			"});\n",
		"/** @constructor @struct */var C=function(){};\n"+
			"$jscomp.global.Object.defineProperties(C.prototype, {});\n")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testEs6ArrowFunction()
func (c *RemoveUnusedCodeClassPropertiesTest) TestEs6ArrowFunction() {
	c.Test(
		"const arrow = () => this.a = 1;", // preserve newline
		"const arrow = () =>          1;")
	c.TestSame("const arrow = () => ({a: 2})")
	c.TestSame("var y = {}; const arrow = () => {y.a = 2; this.a = 2;}")
	c.Test(
		"function A() {\n"+
			"  this.foo = 1;\n"+
			"}\n"+
			"A.prototype.foo = 0;\n"+
			"A.prototype.getIncr = function() {\n"+
			"  return () => { this.foo++; };\n"+
			"};\n"+
			"new A().getIncr()\n",
		"function A() {\n"+
			"\n"+
			"}\n"+
			"\n"+
			"A.prototype.getIncr = function() {\n"+
			"  return () => {             };\n"+
			"};\n"+
			"new A().getIncr()\n")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testEs6Generator()
func (c *RemoveUnusedCodeClassPropertiesTest) TestEs6Generator() {
	c.Test(
		"function* gen() { yield this.a = 1; }", // preserve newline
		"function* gen() { yield          1; }")
	c.TestSame("function* gen() { yield this.a = 1; yield this.a; }")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testEs6Destructuring()
func (c *RemoveUnusedCodeClassPropertiesTest) TestEs6Destructuring() {
	// Test normal destructuring removal
	c.Test(
		"[this.x, this.y] = [1, 2]", // preserve newline
		"[              ] = [1, 2]")

	// Test normal destructuring, assignment prevent removal
	c.Test(
		"[this.x, this.y] = [1, 2] // preserve newline\n"+
			"var p = this.x;\n",
		"[this.x        ] = [1, 2] // preserve newline\n"+
			"var p = this.x;\n")

	// Test rest destructuring, `this` property
	c.Test(
		"[this.x, ...this.z] = [1, 2, 3]", // preserve newline
		"[      , ...this.z] = [1, 2, 3]")

	// Test rest destructuring with normal variable
	c.Test(
		"let z; [this.x, ...z] = [1, 2]", // preserve newline
		"let z; [      , ...z] = [1, 2]")

	// Test rest destructuring, assignment prevent removal
	c.Test(
		"[this.x, ...this.y] = [1, 2]; // preserve newline\n"+
			"var p = this.y;\n",
		"[      , ...this.y] = [1, 2]; // preserve newline\n"+
			"var p = this.y;\n")

	// Test destructuring rhs prevent removal
	c.TestSame(
		"let a;\n" +
			"this.x = 1; // preserve newline\n" +
			"this.y = 2;\n" +
			"[...a] = [this.x, this.y];\n")

	// Test nested destructuring
	c.Test(
		"let z; [this.x, [this.y, ...z]] = [1, [2]]", // preserve newline
		"let z; [      , [      , ...z]] = [1, [2]]")

	// Test normal object destructuring full removal
	c.Test("({a: this.x, b: this.y} = {a: 1, b: 2})", "({} = {a: 1, b: 2})")

	// Test normal object destructuring partial removal
	c.Test("let y; ({a: this.x, b: y} = {a: 1, b: 2})", "let y; ({           b: y} = {a: 1, b: 2})")

	// Test obj destructuring prevent removal
	c.Test(
		"({a: this.x, b: this.y} = {a: 1, b: 2});\n"+
			"var p = this.x;\n",
		"({a: this.x} = {a: 1, b: 2});\n"+
			"var p = this.x;\n")

	// Test obj destructuring with old style class
	c.TestSame(
		"/** @constructor */ function C () {\n" +
			"  this.a = 1;\n" +
			"}\n" +
			"let x;\n" +
			"({a: x} = new C());\n")

	// Test obj destructuring with new style class
	c.TestSame(
		"class C {\n" +
			"  constructor() {\n" +
			"     this.a = 1;\n" +
			"  }\n" +
			"}\n" +
			"let x;\n" +
			"({a: x} = new C());\n")

	// Test let destructuring
	c.TestSame(
		"class C {\n" +
			"  constructor() {\n" +
			"     this.a = 1;\n" +
			"  }\n" +
			"}\n" +
			"let {a: x} = new C();\n")

	// Test obj created at a different location and later used in destructuring
	c.TestSame(
		"class C {\n" +
			"  constructor() {\n" +
			"     this.a = 1;\n" +
			"  }\n" +
			"}\n" +
			"var obj = new C()\n" +
			"let x;\n" +
			"({a: x} = obj);\n")

	// Test obj destructuring with default value
	c.TestSame(
		"class C {\n" +
			"  constructor() {\n" +
			"     this.a = 1;\n" +
			"  }\n" +
			"}\n" +
			"let a;\n" +
			"({a = 2} = new C());\n")

	// Test obj nested destructuring
	c.TestSame(
		"class C {\n" +
			"  constructor() {\n" +
			"     this.a = 1;\n" +
			"  }\n" +
			"}\n" +
			"var obj = new C()\n" +
			"let a;\n" +
			"({x: {a}} = {x: obj});\n")

	// Computed Property string expression doesn't prevent removal.
	c.Test(
		"({['a']:0}); this.a = 1;", // preserve newline
		"({['a']:0});            ")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testDestrucuturing_assginmentToProperty_consideredUse()
func (c *RemoveUnusedCodeClassPropertiesTest) TestDestrucuturing_assginmentToProperty_consideredUse() {
	c.TestSame(
		"class Foo {\n" +
			"  constructor() {\n" +
			"    this.x = 0;\n" +
			"  }\n" +
			"}\n" +
			"\n" +
			"({a: new Foo().x} = {a: 0});\n")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testEs6DefaultParameter()
func (c *RemoveUnusedCodeClassPropertiesTest) TestEs6DefaultParameter() {
	c.Test(
		"function foo(x, y = this.a = 1) {}", // preserve newline
		"function foo(x, y =          1) {}")
	c.TestSame("this.a = 1; function foo(x, y = this.a) {}")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testEs8AsyncFunction()
func (c *RemoveUnusedCodeClassPropertiesTest) TestEs8AsyncFunction() {
	c.Test(
		"async function foo(promise) { // preserve newlines\n"+
			"   this.x = 1;\n"+
			"   return await promise;\n"+
			"}\n",
		"async function foo(promise) { // preserve newlines\n"+
			"\n"+
			"   return await promise;\n"+
			"}\n")

	c.TestSame(
		"async function foo() {\n" +
			"   this.x = 1;\n" +
			"   return await this.x;\n" +
			"}\n")

	c.TestSame(
		"this.x = 1;\n" +
			"async function foo() {\n" +
			"   return await this.x;\n" +
			"}\n")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testField()
func (c *RemoveUnusedCodeClassPropertiesTest) TestField() {
	c.Test(
		"class C {\n"+
			"  x = 1;\n"+
			"  y;\n"+
			"  z = 'hi';\n"+
			"  static x = 1;\n"+
			"  static y;\n"+
			"  static z = 'hi';\n"+
			"}\n",
		"class C {\n"+
			"\n"+
			"}\n")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testStaticFieldReferencingClassName()
func (c *RemoveUnusedCodeClassPropertiesTest) TestStaticFieldReferencingClassName() {
	c.keepGlobals = false
	c.keepLocals = false

	c.Test(
		"const C = class {\n"+
			"  constructor() {\n"+
			"  }\n"+
			"  static x = new C();\n"+
			"  static method() {\n"+
			"    return C.x;\n"+
			"  }\n"+
			"}\n",
		"const C = class {\n"+
			"  constructor() {\n"+
			"  }\n"+
			"  static x = new C();\n"+
			"}\n")

	c.Test(
		"const C = class {\n"+
			"  constructor() {\n"+
			"  }\n"+
			"  static x = new C();\n"+
			"  static y = alert();\n"+
			"  static method() {\n"+
			"    return C.x;\n"+
			"  }\n"+
			"}\n",
		"const C = class {\n"+
			"  constructor() {\n"+
			"  }\n"+
			"  static x = new C();\n"+
			"  static y = alert();\n"+
			"}\n")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testComputedField()
func (c *RemoveUnusedCodeClassPropertiesTest) TestComputedField() {
	c.TestSame(
		"class C {\n" +
			"  ['x'] = 1;\n" +
			"  'y';\n" +
			"  1 = 'hi';\n" +
			"  static ['x'] = 1;\n" +
			"  static 'y';\n" +
			"  static 1 = 'hi';\n" +
			"}\n")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testMixedField()
func (c *RemoveUnusedCodeClassPropertiesTest) TestMixedField() {
	// Computed properties cannot be removed, so only non-computed properties are removed
	c.Test(
		"class C {\n"+
			"  x = 1;\n"+
			"  y;\n"+
			"  z = 'hi';\n"+
			"  static x = 1;\n"+
			"  static y;\n"+
			"  static z = 'hi';\n"+
			"  ['x'] = 1;\n"+
			"  'y';\n"+
			"  1 = 'hi';\n"+
			"  static ['x'] = 1;\n"+
			"  static 'y';\n"+
			"  static 1 = 'hi';\n"+
			"}\n",
		"class C {\n"+
			"  ['x'] = 1;\n"+
			"  'y';\n"+
			"  1 = 'hi';\n"+
			"  static ['x'] = 1;\n"+
			"  static 'y';\n"+
			"  static 1 = 'hi';\n"+
			"}\n")

	c.TestSame(
		"class C {\n" +
			"  [alert()] = 5;\n" +
			"}\n")

	c.TestSame(
		"class C {\n" +
			"  static x = alert();\n" +
			"}\n")

	c.TestSame(
		"class C {\n" +
			"  x = alert();\n" +
			"}\n")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testPureOrBreakMyCode()
func (c *RemoveUnusedCodeClassPropertiesTest) TestPureOrBreakMyCode() {
	c.Test(
		"class C {\n"+
			"  constructor() {\n"+
			"   /** @const */this.used = /** @pureOrBreakMyCode */(alert());\n"+
			"   /** @const */this.unused = /** @pureOrBreakMyCode */(alert());\n"+
			"  }\n"+
			"};\n"+
			"/** @const */C.used = /** @pureOrBreakMyCode */(alert());\n"+
			"/** @const */C.unused = /** @pureOrBreakMyCode */(alert());\n"+
			"function foo() {\n"+
			"  return C.used;\n"+
			"}\n"+
			"foo();\n"+
			"function bar() {\n"+
			"  return new C().used;\n"+
			"}\n"+
			"bar();\n",
		"class C {\n"+
			"  constructor() {\n"+
			"    /** @const */this.used = /** @pureOrBreakMyCode */(alert());\n"+
			"  }\n"+
			"};\n"+
			"/** @const */C.used = /** @pureOrBreakMyCode */(alert());\n"+
			"function foo() {\n"+
			"  return C.used;\n"+
			"}\n"+
			"foo();\n"+
			"function bar() {\n"+
			"  return new C().used;\n"+
			"}\n"+
			"bar();\n")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testUnusedCodeInStaticInitializationBlock()
func (c *RemoveUnusedCodeClassPropertiesTest) TestUnusedCodeInStaticInitializationBlock() {
	c.keepLocals = false

	c.Test(
		"class C {\n"+
			"  static {\n"+
			"    let x = 1;\n"+
			"  }\n"+
			"}\n",
		"class C {\n"+
			"  static {\n"+
			"  }\n"+
			"}\n")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testClassReferencedInFieldInitializer()
func (c *RemoveUnusedCodeClassPropertiesTest) TestClassReferencedInFieldInitializer() {
	c.keepLocals = false
	c.keepGlobals = false

	c.TestSame(
		"class One {\n" +
			"  field = 1;\n" +
			"}\n" +
			"class Two {\n" +
			"  field = new One();\n" +
			"}\n" +
			"const two = new Two();\n" +
			"alert(two);\n" +
			"alert(two.field);\n" +
			"alert(two.field.field);\n")
}

// Java: RemoveUnusedCodeClassPropertiesTest.testClosureInFieldInitializerDoubleTraversalCrash()
func (c *RemoveUnusedCodeClassPropertiesTest) TestClosureInFieldInitializerDoubleTraversalCrash() {
	// TODO(b/571953140): traverses field initializers twice, crashing on removal
	c.keepLocals = false

	e := AssertThrows[jutil.IRuntimeException](func() {
		c.TestSame(
			"class C {\n" +
				"  fn = () => {\n" +
				"    function unused() {}\n" +
				"    return 1;\n" +
				"  };\n" +
				"}\n" +
				"alert(new C().fn());\n")
	})
	// Throwables.getRootCause(e)
	var root jutil.Throwable = e
	for root.GetCause() != nil && root.GetCause() != root {
		root = root.GetCause()
	}
	_, isIllegalState := root.(jutil.IIllegalStateException)
	AssertThat(isIllegalState).IsTrue()
	AssertThat(root.GetMessage()).Contains("Has no parent")
}
