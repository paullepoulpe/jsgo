// Port of com/google/javascript/jscomp/ClosureOptimizePrimitivesTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/ClosureOptimizePrimitivesTest.java)

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

	"github.com/paullepoulpe/jsgo/internal/parser"
)

// ClosureOptimizePrimitivesTest holds tests for ClosureOptimizePrimitives.
type ClosureOptimizePrimitivesTest struct {
	CompilerTestCaseBase

	canUseEs6Syntax bool
}

func NewClosureOptimizePrimitivesTest() *ClosureOptimizePrimitivesTest {
	c := &ClosureOptimizePrimitivesTest{}
	c.InitCompilerTestCaseBase(c)
	c.canUseEs6Syntax = true
	return c
}

func TestClosureOptimizePrimitivesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewClosureOptimizePrimitivesTest() })
}

// Java: ClosureOptimizePrimitivesTest.getProcessor(Compiler)
func (c *ClosureOptimizePrimitivesTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewClosureOptimizePrimitives(compiler, c.canUseEs6Syntax)
}

// Java: ClosureOptimizePrimitivesTest.testObjectCreateOddParams()
func (c *ClosureOptimizePrimitivesTest) TestObjectCreateOddParams() {
	c.TestSame("goog.object.create('a',1,2);")
}

// Java: ClosureOptimizePrimitivesTest.testObjectCreate1()
func (c *ClosureOptimizePrimitivesTest) TestObjectCreate1() {
	c.Test("var a = goog.object.create()", "var a = {}")
}

// Java: ClosureOptimizePrimitivesTest.testObjectCreate2()
func (c *ClosureOptimizePrimitivesTest) TestObjectCreate2() {
	c.Test(
		`var a = module$contents$goog$object_create(
'b',module$contents$goog$object_create('c','d'))
`,
		"var a = {'b':{'c':'d'}};")
}

// Java: ClosureOptimizePrimitivesTest.testObjectCreate3()
func (c *ClosureOptimizePrimitivesTest) TestObjectCreate3() {
	c.Test("var a = goog.object.create(1,2)", "var a = {1:2}")
}

// Java: ClosureOptimizePrimitivesTest.testObjectCreate4()
func (c *ClosureOptimizePrimitivesTest) TestObjectCreate4() {
	c.Test(
		"alert(goog.object.create(1,2).toString())", //
		"alert({1:2}.toString())")
}

// Java: ClosureOptimizePrimitivesTest.testObjectCreate5()
func (c *ClosureOptimizePrimitivesTest) TestObjectCreate5() {
	c.Test("goog.object.create('a',2).toString()", "({'a':2}).toString()")
}

// Java: ClosureOptimizePrimitivesTest.testObjectCreateGoogProvide()
func (c *ClosureOptimizePrimitivesTest) TestObjectCreateGoogProvide() {
	// TODO(user): Delete this once goog.object is a goog.module.
	c.Test("var a = goog$object$create('b',goog$object$create('c','d'))", "var a = {'b':{'c':'d'}};")
}

// Java: ClosureOptimizePrimitivesTest.testObjectCreateRewritten()
func (c *ClosureOptimizePrimitivesTest) TestObjectCreateRewritten() {
	c.EnableRewriteClosureCode()
	c.Test("var a = goog.object.create('a', 1, 'b', 2);", "var a = {'a': 1, 'b': 2};")
}

// Java: ClosureOptimizePrimitivesTest.testObjectCreateRewrittenRemoveParens()
func (c *ClosureOptimizePrimitivesTest) TestObjectCreateRewrittenRemoveParens() {
	c.EnableRewriteClosureCode()
	c.Test(
		"var a = goog.object.create(('a'), 1, ('b'), 2);", //
		"var a = {'a': 1, 'b': 2};")
}

// Java: ClosureOptimizePrimitivesTest.testObjectCreateNonConstKey1()
func (c *ClosureOptimizePrimitivesTest) TestObjectCreateNonConstKey1() {
	c.Test(
		"var a = goog.object.create('a', 1, 2, 3, foo, bar);", //
		"var a = {'a': 1, 2: 3, [foo]: bar};")

	// assertFS(getLastCompiler().getAllowableFeatures()).has(Feature.COMPUTED_PROPERTIES)
	AssertThat(c.GetLastCompiler().GetAllowableFeatures().Has(parser.FeatureSet_Feature_COMPUTED_PROPERTIES)).IsTrue()
}

// Java: ClosureOptimizePrimitivesTest.testObjectCreateNonConstKey2()
func (c *ClosureOptimizePrimitivesTest) TestObjectCreateNonConstKey2() {
	c.Test("var a = goog.object.create('a' + 'b', 0);", "var a = {['a' + 'b']: 0};")
}

// Java: ClosureOptimizePrimitivesTest.testObjectCreateNonConstKey3()
func (c *ClosureOptimizePrimitivesTest) TestObjectCreateNonConstKey3() {
	c.Test(
		`var a = module$contents$goog$object_create(
i++,module$contents$goog$object_create(foo(), 'd'))
`,
		"var a = {[i++]: {[foo()]: 'd'}};")
}

// Java: ClosureOptimizePrimitivesTest.testObjectCreateNonConstKey4()
func (c *ClosureOptimizePrimitivesTest) TestObjectCreateNonConstKey4() {
	c.Test(
		"alert(goog.object.create(a = 1, 2).toString())", //
		"alert({[a = 1]: 2}.toString())")
}

// Java: ClosureOptimizePrimitivesTest.testObjectCreateNonConstKey5()
func (c *ClosureOptimizePrimitivesTest) TestObjectCreateNonConstKey5() {
	c.Test(
		"goog.object.create(function foo() {}, 2).toString()", //
		"({[function foo() {}]: 2}).toString()")
}

// Java: ClosureOptimizePrimitivesTest.testObjectCreateNonConstKeyGoogProvide()
func (c *ClosureOptimizePrimitivesTest) TestObjectCreateNonConstKeyGoogProvide() {
	// TODO(user): Delete this once goog.object is a goog.module.
	c.Test(
		"var a = goog$object$create(i++,goog$object$create(foo(), 'd'))",
		"var a = {[i++]: {[foo()]: 'd'}};")
}

// Java: ClosureOptimizePrimitivesTest.testObjectCreateNonConstKeyNotEs6()
func (c *ClosureOptimizePrimitivesTest) TestObjectCreateNonConstKeyNotEs6() {
	c.canUseEs6Syntax = false
	c.TestSame("var a = goog.object.create(foo, bar)")
}

// Java: ClosureOptimizePrimitivesTest.testObjectCreateSetZeroArg()
func (c *ClosureOptimizePrimitivesTest) TestObjectCreateSetZeroArg() {
	c.Test("var a = goog.object.createSet()", "var a = {}")
}

// Java: ClosureOptimizePrimitivesTest.testObjectCreateSetTwoNumericalArgs()
func (c *ClosureOptimizePrimitivesTest) TestObjectCreateSetTwoNumericalArgs() {
	c.Test("var a = goog.object.createSet(1,2)", "var a = {1:true,2:true}")
}

// Java: ClosureOptimizePrimitivesTest.testObjectCreateSetOneNumericalArg()
func (c *ClosureOptimizePrimitivesTest) TestObjectCreateSetOneNumericalArg() {
	c.Test(
		"alert(goog.object.createSet(1).toString())", //
		"alert({1:true}.toString())")
}

// Java: ClosureOptimizePrimitivesTest.testObjectCreateSetOneStringArg()
func (c *ClosureOptimizePrimitivesTest) TestObjectCreateSetOneStringArg() {
	c.Test("goog.object.createSet('a').toString()", "({'a':true}).toString()")
}

// Java: ClosureOptimizePrimitivesTest.testObjectCreateSet_duplicate()
func (c *ClosureOptimizePrimitivesTest) TestObjectCreateSet_duplicate() {
	c.TestWarning("goog.object.createSet('a', 'a')", ClosureOptimizePrimitives_DUPLICATE_SET_MEMBER)
	c.TestWarning("goog.object.createSet(4, 4)", ClosureOptimizePrimitives_DUPLICATE_SET_MEMBER)
	c.TestWarning("goog.object.createSet(4, '4')", ClosureOptimizePrimitives_DUPLICATE_SET_MEMBER)
}

// Java: ClosureOptimizePrimitivesTest.testObjectCreateSetNonConstKey1()
func (c *ClosureOptimizePrimitivesTest) TestObjectCreateSetNonConstKey1() {
	c.Test(
		"var a = goog.object.createSet(foo, bar);", //
		"var a = {[foo]: true, [bar]: true};")
}

// Java: ClosureOptimizePrimitivesTest.testObjectCreateSetNonConstKey2()
func (c *ClosureOptimizePrimitivesTest) TestObjectCreateSetNonConstKey2() {
	c.Test(
		"alert(module$contents$goog$object_createSet(a = 1, 2).toString())",
		"alert({[a = 1]: true, 2: true}.toString())")
}

// Java: ClosureOptimizePrimitivesTest.testObjectCreateSetNonConstKeyGoogProvide()
func (c *ClosureOptimizePrimitivesTest) TestObjectCreateSetNonConstKeyGoogProvide() {
	// TODO(user): Delete this once goog.object is a goog.module.
	c.Test(
		"alert(goog$object$createSet(a = 1, 2).toString())",
		"alert({[a = 1]: true, 2: true}.toString())")
}

// Java: ClosureOptimizePrimitivesTest.testObjectCreateSetNonConstSingleKey()
func (c *ClosureOptimizePrimitivesTest) TestObjectCreateSetNonConstSingleKey() {
	// this triggers the 'first argument may be an array' case and we back off
	c.TestSame("goog.object.createSet(() => {}).toString()")
}

// Java: ClosureOptimizePrimitivesTest.testObjectCreateSetNonConstKeyNotEs6()
func (c *ClosureOptimizePrimitivesTest) TestObjectCreateSetNonConstKeyNotEs6() {
	c.canUseEs6Syntax = false
	c.TestSame("var a = goog.object.createSet(foo, bar);")
}

// Java: ClosureOptimizePrimitivesTest.testObjectCreateSetSingleNonLiteralArg()
func (c *ClosureOptimizePrimitivesTest) TestObjectCreateSetSingleNonLiteralArg() {
	c.TestSame("const arr = [1, 2, 3]; const s = goog.object.createSet(arr);")
	c.TestSame("const num = 1; const s = goog.object.createSet(num);")
	c.TestSame("const s = goog.object.createSet(undefined || [1, 2, 3]);")
}

// Java: ClosureOptimizePrimitivesTest.testObjectCreateSetSingleLiteralArg()
func (c *ClosureOptimizePrimitivesTest) TestObjectCreateSetSingleLiteralArg() {
	c.Test("const s = goog.object.createSet(3);", "const s = {3: true};")
	c.Test("const s = goog.object.createSet('a');", "const s = {'a': true};")
}

// Java: ClosureOptimizePrimitivesTest.testPropertyReflectionAdvanced()
func (c *ClosureOptimizePrimitivesTest) TestPropertyReflectionAdvanced() {
	c.Test("goog.reflect.objectProperty('push', [])", "JSCompiler_renameProperty('push', [])")
	c.TestSame("JSCompiler_renameProperty('push', [])")
}

// Java: ClosureOptimizePrimitivesTest.testEs6ArrowCompatibility()
func (c *ClosureOptimizePrimitivesTest) TestEs6ArrowCompatibility() {
	// Arrow
	c.Test("var f = () => goog.object.create(1, 2);", "var f = () => ({1: 2});")
}

// Java: ClosureOptimizePrimitivesTest.testEs6ClassCompatibility()
func (c *ClosureOptimizePrimitivesTest) TestEs6ClassCompatibility() {
	c.Test(
		`class C {
  constructor() {
    this.x = goog.object.create(1, 2);
  }
}
`,
		`class C {
  constructor() {
    this.x = {1: 2};
  }
}
`)
}

// Java: ClosureOptimizePrimitivesTest.testEs6MemberFunctionDefCompatibility()
func (c *ClosureOptimizePrimitivesTest) TestEs6MemberFunctionDefCompatibility() {
	c.Test(
		`var obj = {
  method() {
    return goog.object.create('a', 2);
  }
}
`,
		`var obj = {
  method() {
    return {'a': 2};
  }
}
`)
}

// Java: ClosureOptimizePrimitivesTest.testEs6ComputedPropCompatibility()
func (c *ClosureOptimizePrimitivesTest) TestEs6ComputedPropCompatibility() {
	c.Test(
		`var obj = {
  [goog.object.create(1, 2)]: 42
}
`,
		`var obj = {
  [{1: 2}]: 42
}
`)
}

// Java: ClosureOptimizePrimitivesTest.testEs6TemplateLitCompatibility()
func (c *ClosureOptimizePrimitivesTest) TestEs6TemplateLitCompatibility() {
	c.Test(
		"function tag(strings) {\n  return goog.object.create('a', 2);\n}\ntag`template`\n",
		"function tag(strings) {\n  return {'a': 2};\n}\ntag`template`\n")
}

// Java: ClosureOptimizePrimitivesTest.testEs6DestructuringCompatibility()
func (c *ClosureOptimizePrimitivesTest) TestEs6DestructuringCompatibility() {
	c.Test("var {a: x} = goog.object.create('a', 2);", "var {a: x} = {'a': 2};")
}

// Java: ClosureOptimizePrimitivesTest.testEs6AsyncCompatibility()
func (c *ClosureOptimizePrimitivesTest) TestEs6AsyncCompatibility() {
	c.Test(
		`async function foo() {
   return await goog.object.create('a', 2);
}
`,
		`async function foo() {
   return await {'a': 2};
}
`)
}
