// Port of com/google/javascript/jscomp/ConvertToDottedPropertiesTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/ConvertToDottedPropertiesTest.java)

/*
 * Copyright 2007 The Closure Compiler Authors.
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

import "testing"

// ConvertToDottedPropertiesTest is the tests for ConvertToDottedProperties.
type ConvertToDottedPropertiesTest struct {
	CompilerTestCaseBase
}

// Java: ConvertToDottedPropertiesTest()
func NewConvertToDottedPropertiesTest() *ConvertToDottedPropertiesTest {
	c := &ConvertToDottedPropertiesTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestConvertToDottedPropertiesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewConvertToDottedPropertiesTest() })
}

// Java: ConvertToDottedPropertiesTest.getProcessor(Compiler)
func (c *ConvertToDottedPropertiesTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewConvertToDottedProperties(compiler)
}

// Java: ConvertToDottedPropertiesTest.testConvert()
func (c *ConvertToDottedPropertiesTest) TestConvert() {
	c.Test("a['p']", "a.p")
	c.Test("a['_p_']", "a._p_")
	c.Test("a['_']", "a._")
	c.Test("a['$']", "a.$")
	c.Test("a.b.c['p']", "a.b.c.p")
	c.Test("a.b['c'].p", "a.b.c.p")
	c.Test("a['p']();", "a.p();")
	c.Test("a()['p']", "a().p")
	// ASCII in Unicode is safe.
	c.Test("a['AA']", "a.AA")
}

// Java: ConvertToDottedPropertiesTest.testDoNotConvert()
func (c *ConvertToDottedPropertiesTest) TestDoNotConvert() {
	c.TestSame("a[0]")
	c.TestSame("a['']")
	c.TestSame("a[' ']")
	c.TestSame("a[',']")
	c.TestSame("a[';']")
	c.TestSame("a[':']")
	c.TestSame("a['.']")
	c.TestSame("a['0']")
	c.TestSame("a['p ']")
	c.TestSame("a['p' + '']")
	c.TestSame("a[p]")
	c.TestSame("a[P]")
	c.TestSame("a[$]")
	c.TestSame("a[p()]")
	c.TestSame("a['default']")
	// Ignorable control characters are ok in Java identifiers, but not in JS.
	c.TestSame("a['A\x04']")
	// upper case lower half of o from phonetic extensions set.
	// valid in Safari, not in Firefox, IE.
	c.TestSame("a['\u1d17A']")
	// Latin capital N with tilde - nice if we handled it, but for now let's
	// only allow simple Latin (aka ASCII) to be converted.
	c.TestSame("a['\u00d1StuffAfter']")
}

// Java: ConvertToDottedPropertiesTest.testAlreadyDotted()
func (c *ConvertToDottedPropertiesTest) TestAlreadyDotted() {
	c.TestSame("a.b")
	c.TestSame("var a = {b: 0};")
}

// Java: ConvertToDottedPropertiesTest.testQuotedProps()
func (c *ConvertToDottedPropertiesTest) TestQuotedProps() {
	c.TestSame("({'':0})")
	c.TestSame("({'1.0':0})")
	c.TestSame("({'\u1d17A':0})")
	c.TestSame("({'a\x04b':0})")
}

// Java: ConvertToDottedPropertiesTest.test5746867()
func (c *ConvertToDottedPropertiesTest) Test5746867() {
	c.TestSame("var a = { '$\\\\' : 5 };")
	c.TestSame("var a = { 'x\\\\u0041$\\\\' : 5 };")
}

// Java: ConvertToDottedPropertiesTest.testOptionalChaining()
func (c *ConvertToDottedPropertiesTest) TestOptionalChaining() {
	c.Test("data?.['name']", "data?.name")
	c.Test("data?.['name']?.['first']", "data?.name?.first")
	c.Test("data['name']?.['first']", "data.name?.first")
	c.TestSame("a?.[0]")
	c.TestSame("a?.['']")
	c.TestSame("a?.[' ']")
	c.TestSame("a?.[',']")
	c.TestSame("a?.[';']")
	c.TestSame("a?.[':']")
	c.TestSame("a?.['.']")
	c.TestSame("a?.['0']")
	c.TestSame("a?.['p ']")
	c.TestSame("a?.['p' + '']")
	c.TestSame("a?.[p]")
	c.TestSame("a?.[P]")
	c.TestSame("a?.[$]")
	c.TestSame("a?.[p()]")
	c.TestSame("a?.['default']")
}

// Java: ConvertToDottedPropertiesTest.testComputedPropertyOrField()
func (c *ConvertToDottedPropertiesTest) TestComputedPropertyOrField() {
	c.Test("const test1 = {['prop1']:87};", "const test1 = {prop1:87};")
	c.Test(
		"const test1 = {['prop1']:87,['prop2']:bg,['prop3']:'hfd'};",
		"const test1 = {prop1:87,prop2:bg,prop3:'hfd'};")
	c.Test(
		"o = {['x']: async function(x) { return await x + 1; }};",
		"o = {x:async function (x) { return await x + 1; }};")
	c.Test("o = {['x']: function*(x) {}};", "o = {x: function*(x) {}};")
	c.Test(
		"o = {['x']: async function*(x) { return await x + 1; }};",
		"o = {x:async function*(x) { return await x + 1; }};")
	c.Test("class C {'x' = 0;  ['y'] = 1;}", "class C { x= 0;y= 1;}")
	c.Test("class C {'m'() {} }", "class C {m() {}}")
	c.Test("const o = {'b'() {}, ['c']() {}};", "const o = {b: function() {}, c:function(){}};")
	c.Test("o = {['x']: () => this};", "o = {x: () => this};")
	c.Test("const o = {get ['d']() {}};", "const o = {get d() {}};")
	c.Test("const o = { set ['e'](x) {}};", "const o = { set e(x) {}};")
	c.Test(
		"class C {'m'() {}  ['n']() {} 'x' = 0;  ['y'] = 1;}",
		"class C {m() {}  n() {} x= 0;y= 1;}")
	c.Test("const o = { get ['d']() {},  set ['e'](x) {}};", "const o = {get d() {},  set e(x){}};")
	c.Test(
		"const o = {['a']: 1,'b'() {}, ['c']() {},  get ['d']() {},  set ['e'](x) {}};",
		"const o = {a: 1,b: function() {}, c: function() {},  get d() {},  set e(x) {}};")
	// test static keyword
	c.Test(
		"class C {\n"+
			"'m'(){}\n"+
			"['n'](){}\n"+
			"static 'x' = 0;\n"+
			"static ['y'] = 1;}\n",
		"class C {\n"+
			"m(){}\n"+
			"n(){}\n"+
			"static x = 0;\n"+
			"static y= 1;}\n")
	c.Test(
		"window[\"MyClass\"] = class {\n"+
			"static [\"Register\"](){}\n"+
			"};\n",
		"window.MyClass = class {\n"+
			"static Register(){}\n"+
			"};\n")
	c.Test(
		"class C {\n"+
			"'method'(){}\n"+
			"async ['method1'](){}\n"+
			"*['method2'](){}\n"+
			"static ['smethod'](){}\n"+
			"static async ['smethod1'](){}\n"+
			"static *['smethod2'](){}}\n",
		"class C {\n"+
			"method(){}\n"+
			"async method1(){}\n"+
			"*method2(){}\n"+
			"static smethod(){}\n"+
			"static async smethod1(){}\n"+
			"static *smethod2(){}}\n")
	c.TestSame("const o = {[fn()]: 0}")
	c.TestSame("const test1 = {[0]:87};")
	c.TestSame("const test1 = {['default']:87};")
	c.TestSame("class C { ['constructor']() {} }")
	c.TestSame("class C { ['constructor'] = 0 }")
	c.TestSame("const o = {['__proto__']: v};")
	c.Test("class C { ['__proto__']() {} }", "class C { __proto__() {} }")
	c.Test("class C { ['__proto__'] = 0; }", "class C { __proto__ = 0; }")
}

// Java: ConvertToDottedPropertiesTest.testComputedPropertyWithDefaultValue()
func (c *ConvertToDottedPropertiesTest) TestComputedPropertyWithDefaultValue() {
	c.Test("const {['o']: o = 0} = {};", "const {o:o = 0} = {};")
}

// Java: ConvertToDottedPropertiesTest.testContinueOptionalChaining()
func (c *ConvertToDottedPropertiesTest) TestContinueOptionalChaining() {
	c.Test("const opt1 = window?.a?.['b'];", "const opt1 = window?.a?.b;")
	c.Test("const opt2 = window?.a['b'];", "const opt2 = window?.a.b;")
	c.Test(
		"const chain =\n"+
			"window['a'].x.y.b.x.y['c'].x.y?.d.x.y['e'].x.y\n"+
			"['f-f'].x.y?.['g-g'].x.y?.['h'].x.y['i'].x.y;\n",
		"const chain = window.a.x.y.b.x.y.c.x.y?.d.x.y.e.x.y\n"+
			"['f-f'].x.y?.['g-g'].x.y?.h.x.y.i.x.y;\n")
}

// Java: ConvertToDottedPropertiesTest.testClassPrivateProp_notPrivateProp_computedProp()
func (c *ConvertToDottedPropertiesTest) TestClassPrivateProp_notPrivateProp_computedProp() {
	c.TestSame("const x = {['#notAPrivateProp']: 1}")
}

// Java: ConvertToDottedPropertiesTest.testClassPrivateProp_notPrivateProp_computedFieldDef()
func (c *ConvertToDottedPropertiesTest) TestClassPrivateProp_notPrivateProp_computedFieldDef() {
	c.TestSame("class C {['#notAPrivateProp'] = 1}")
}

// Java: ConvertToDottedPropertiesTest.testClassPrivateProp_notPrivateProp_getterDef()
func (c *ConvertToDottedPropertiesTest) TestClassPrivateProp_notPrivateProp_getterDef() {
	c.TestSame("const x = {get '#notAPrivateProp'() {}}")
	c.TestSame("class C {get '#notAPrivateProp'() {}}")
}

// Java: ConvertToDottedPropertiesTest.testClassPrivateProp_notPrivateProp_setterDef()
func (c *ConvertToDottedPropertiesTest) TestClassPrivateProp_notPrivateProp_setterDef() {
	c.TestSame("const x = {set '#notAPrivateProp'(x) {}}")
	c.TestSame("class C {set '#notAPrivateProp'(x) {}}")
}

// Java: ConvertToDottedPropertiesTest.testClassPrivateProp_notPrivateProp_stringKey()
func (c *ConvertToDottedPropertiesTest) TestClassPrivateProp_notPrivateProp_stringKey() {
	c.TestSame("const x = {'#notAPrivateProp': 1}")
}

// Java: ConvertToDottedPropertiesTest.testClassPrivateProp_notPrivateProp_optchainGetelem()
func (c *ConvertToDottedPropertiesTest) TestClassPrivateProp_notPrivateProp_optchainGetelem() {
	c.TestSame("x?.['#notAPrivateProp']")
}

// Java: ConvertToDottedPropertiesTest.testClassPrivateProp_notPrivateProp_getelem()
func (c *ConvertToDottedPropertiesTest) TestClassPrivateProp_notPrivateProp_getelem() {
	c.TestSame("x['#notAPrivateProp']")
}

// Java: ConvertToDottedPropertiesTest.testUninitializedComputedField()
func (c *ConvertToDottedPropertiesTest) TestUninitializedComputedField() {
	c.Test("class C { ['x']; }", "class C { x; }")
}

// Java: ConvertToDottedPropertiesTest.testUninitializedQuotedField()
func (c *ConvertToDottedPropertiesTest) TestUninitializedQuotedField() {
	c.Test("class C { 'x'; }", "class C { x; }")
}

// Java: ConvertToDottedPropertiesTest.testUninitializedStaticComputedField()
func (c *ConvertToDottedPropertiesTest) TestUninitializedStaticComputedField() {
	c.Test("class C { static ['x']; }", "class C { static x; }")
}

// Java: ConvertToDottedPropertiesTest.testUninitializedStaticQuotedField()
func (c *ConvertToDottedPropertiesTest) TestUninitializedStaticQuotedField() {
	c.Test("class C { static 'x'; }", "class C { static x; }")
}

// Java: ConvertToDottedPropertiesTest.testUninitializedComputedField_notConverted()
func (c *ConvertToDottedPropertiesTest) TestUninitializedComputedField_notConverted() {
	c.TestSame("class C { ['constructor']; }")
	c.TestSame("class C { ['default']; }")
	c.TestSame("class C { ['#notAPrivateProp']; }")
}
