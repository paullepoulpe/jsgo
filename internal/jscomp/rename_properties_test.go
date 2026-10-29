// Port of com/google/javascript/jscomp/RenamePropertiesTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/RenamePropertiesTest.java)

/*
 * Copyright 2005 The Closure Compiler Authors.
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
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// RenamePropertiesTest tests RenameProperties.
//
// The Java @Test methods not named test* get a "Test" prefix in Go so that
// RunJUnit4 runs them.
//
// Java: public final class RenamePropertiesTest extends CompilerTestCase
type RenamePropertiesTest struct {
	CompilerTestCaseBase

	renameProperties    *RenameProperties
	generatePseudoNames bool
	prevUsedPropertyMap *VariableMap               // may be nil
	propertyNameFilter  func(name rhino.Node) bool // may be nil
}

// Java: RenamePropertiesTest.EXTERNS
const renamePropertiesTest_EXTERNS = `var window;
prop.toString;
var google = { gears: { factory: {}, workerPool: {} } };
`

// Java: RenamePropertiesTest()
func NewRenamePropertiesTest() *RenamePropertiesTest {
	c := &RenamePropertiesTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, renamePropertiesTest_EXTERNS)
	c.propertyNameFilter = func(name rhino.Node) bool { return true }
	return c
}

func TestRenamePropertiesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewRenamePropertiesTest() })
}

// Java: RenamePropertiesTest.setUp()
func (c *RenamePropertiesTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.generatePseudoNames = false
	c.prevUsedPropertyMap = nil
	c.EnableNormalize()
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.EnableGatherExternProperties()
}

// Java: RenamePropertiesTest.testPrototypeProperties()
func (c *RenamePropertiesTest) TestPrototypeProperties() {
	c.Test(
		`Bar.prototype.getA = function(){}; bar.getA();
Bar.prototype.getB = function(){};
Bar.prototype.getC = function(){};
`,
		`Bar.prototype.a = function(){}; bar.a();
Bar.prototype.b = function(){};
Bar.prototype.c = function(){}
`)
}

// Java: RenamePropertiesTest.testPrototypeProperties_optChain()
func (c *RenamePropertiesTest) TestPrototypeProperties_optChain() {
	c.Test(
		`Bar.prototype.getA = function(){}; bar?.getA();
Bar.prototype.getB = function(){};
Bar.prototype.getC = function(){};
`,
		`Bar.prototype.a = function(){}; bar?.a();
Bar.prototype.b = function(){};
Bar.prototype.c = function(){}
`)
}

// Java: RenamePropertiesTest.testPrototypePropertiesAsObjLitKeys1()
func (c *RenamePropertiesTest) TestPrototypePropertiesAsObjLitKeys1() {
	c.Test(
		"Bar.prototype = {2: function(){}, getA: function(){}}; bar[2]();",
		"Bar.prototype = {2: function(){}, a: function(){}}; bar[2]();")
}

// Java: RenamePropertiesTest.testPrototypePropertiesAsObjLitKeys2()
func (c *RenamePropertiesTest) TestPrototypePropertiesAsObjLitKeys2() {
	c.TestSame("Bar.prototype = {get 2(){}}; bar[2];")

	c.TestSame("Bar.prototype = {get 'a'(){}}; bar['a'];")

	c.Test("Bar.prototype = {get getA(){}}; bar.getA;", "Bar.prototype = {get a(){}}; bar.a;")
}

// Java: RenamePropertiesTest.testPrototypePropertiesAsObjLitKeys3()
func (c *RenamePropertiesTest) TestPrototypePropertiesAsObjLitKeys3() {
	c.TestSame("Bar.prototype = {set 2(x){}}; bar[2];")

	c.TestSame("Bar.prototype = {set 'a'(x){}}; bar['a'];")

	c.Test("Bar.prototype = {set getA(x){}}; bar.getA;", "Bar.prototype = {set a(x){}}; bar.a;")

	c.Test("Bar.prototype = {set getA(x){}}; bar?.getA;", "Bar.prototype = {set a(x){}}; bar?.a;")
}

// Java: RenamePropertiesTest.testMixedQuotedAndUnquotedObjLitKeys1()
func (c *RenamePropertiesTest) TestMixedQuotedAndUnquotedObjLitKeys1() {
	c.Test(
		"Bar = {getA: function(){}, 'getB': function(){}}; bar.getA();",
		"Bar = {a: function(){}, 'getB': function(){}}; bar.a();")

	c.Test(
		"Bar = {getA: function(){}, 'getB': function(){}}; bar?.getA();",
		"Bar = {a: function(){}, 'getB': function(){}}; bar?.a();")
}

// Java: RenamePropertiesTest.testMixedQuotedAndUnquotedObjLitKeys2()
func (c *RenamePropertiesTest) TestMixedQuotedAndUnquotedObjLitKeys2() {
	c.Test(
		"Bar = {getA: function(){}, 'getB': function(){}}; bar.getA();",
		"Bar = {a: function(){}, 'getB': function(){}}; bar.a();")
}

// Java: RenamePropertiesTest.testQuotedPrototypeProperty()
func (c *RenamePropertiesTest) TestQuotedPrototypeProperty() {
	c.TestSame("Bar.prototype['getA'] = function(){}; bar['getA']();")
	c.TestSame("Bar.prototype['getA'] = function(){}; bar?.['getA']();")
}

// Java: RenamePropertiesTest.testOverlappingOriginalAndGeneratedNames()
func (c *RenamePropertiesTest) TestOverlappingOriginalAndGeneratedNames() {
	c.Test(
		"Bar.prototype = {b: function(){}, a: function(){}}; bar.b();",
		"Bar.prototype = {a: function(){}, b: function(){}}; bar.a();")
}

// Java: RenamePropertiesTest.testRenamePropertiesWithLeadingUnderscores()
func (c *RenamePropertiesTest) TestRenamePropertiesWithLeadingUnderscores() {
	c.Test(
		"Bar.prototype = {_getA: function(){}, _b: 0}; bar._getA();",
		"Bar.prototype = {a: function(){}, b: 0}; bar.a();")

	c.Test(
		"Bar.prototype = {_getA: function(){}, _b: 0}; bar?._getA();",
		"Bar.prototype = {a: function(){}, b: 0}; bar?.a();")
}

// Java: RenamePropertiesTest.testPropertyAddedToObject()
func (c *RenamePropertiesTest) TestPropertyAddedToObject() {
	c.Test("var foo = {}; foo.prop = '';", "var foo = {}; foo.a = '';")
}

// Java: RenamePropertiesTest.testPropertyAddedToFunction()
func (c *RenamePropertiesTest) TestPropertyAddedToFunction() {
	c.Test("var foo = function(){}; foo.prop = '';", "var foo = function(){}; foo.a = '';")
}

// Java: RenamePropertiesTest.testPropertyOfObjectOfUnknownType()
func (c *RenamePropertiesTest) TestPropertyOfObjectOfUnknownType() {
	c.Test("var foo = x(); foo.prop = '';", "var foo = x(); foo.a = '';")
}

// Java: RenamePropertiesTest.testSetPropertyOfThis()
func (c *RenamePropertiesTest) TestSetPropertyOfThis() {
	c.Test("this.prop = 'bar'", "this.a = 'bar'")
}

// Java: RenamePropertiesTest.testReadPropertyOfThis()
func (c *RenamePropertiesTest) TestReadPropertyOfThis() {
	c.Test("f(this.prop);", "f(this.a);")

	c.Test("f(this?.prop);", "f(this?.a);")
}

// Java: RenamePropertiesTest.testObjectLiteralInLocalScope()
func (c *RenamePropertiesTest) TestObjectLiteralInLocalScope() {
	c.Test(
		"function x() { var foo = {prop1: 'bar', prop2: 'baz'}; }",
		"function x() { var foo = {a: 'bar', b: 'baz'}; }")
}

// Java: RenamePropertiesTest.testIncorrectAttemptToAccessQuotedProperty()
func (c *RenamePropertiesTest) TestIncorrectAttemptToAccessQuotedProperty() {
	// The correct way to call the quoted 'getFoo' method is: bar['getFoo']().
	c.Test(
		"Bar.prototype = {'B': 0, 'getFoo': function(){}}; bar.getFoo();",
		"Bar.prototype = {'B': 0, 'getFoo': function(){}}; bar.a();")
}

// Java: RenamePropertiesTest.testSetQuotedPropertyOfThis()
func (c *RenamePropertiesTest) TestSetQuotedPropertyOfThis() {
	c.TestSame("this['prop'] = 'bar';")
}

// Java: RenamePropertiesTest.testExternedPropertyName()
func (c *RenamePropertiesTest) TestExternedPropertyName() {
	c.Test(
		"Bar.prototype = {toString: function(){}, foo: 0}; bar.toString();",
		"Bar.prototype = {toString: function(){}, a: 0}; bar.toString();")
}

// Java: RenamePropertiesTest.testExternedPropertyNameDefinedByObjectLiteral()
func (c *RenamePropertiesTest) TestExternedPropertyNameDefinedByObjectLiteral() {
	c.TestSame("function x() { var foo = google.gears.factory; }")
}

// Java: RenamePropertiesTest.testAvoidingConflictsBetweenQuotedAndUnquotedPropertyNames()
func (c *RenamePropertiesTest) TestAvoidingConflictsBetweenQuotedAndUnquotedPropertyNames() {
	c.Test(
		"Bar.prototype.foo = function(){}; Bar.prototype['a'] = 0; bar.foo();",
		"Bar.prototype.b = function(){}; Bar.prototype['a'] = 0; bar.b();")
}

// Java: RenamePropertiesTest.testSamePropertyNameQuotedAndUnquoted()
func (c *RenamePropertiesTest) TestSamePropertyNameQuotedAndUnquoted() {
	c.Test(
		"Bar.prototype.prop = function(){}; y = {'prop': 0};",
		"Bar.prototype.a = function(){}; y = {'prop': 0};")
}

// Java: RenamePropertiesTest.testStaticAndInstanceMethodWithSameName()
func (c *RenamePropertiesTest) TestStaticAndInstanceMethodWithSameName() {
	c.Test(
		`Bar = function(){}; Bar.getA = function(){};
Bar.prototype.getA = function(){}; Bar.getA(); bar.getA();
`,
		`Bar = function(){}; Bar.a = function(){};
Bar.prototype.a = function(){}; Bar.a(); bar.a();
`)
}

// Java: RenamePropertiesTest.testRenamePropertiesFunctionCall1()
func (c *RenamePropertiesTest) TestRenamePropertiesFunctionCall1() {
	c.Test(
		"var foo = {myProp: 0}; f(foo[JSCompiler_renameProperty('myProp')]);",
		"var foo = {a: 0}; f(foo['a']);")
	// JSCompiler does not handle optional calls to `JSCompiler_renameProperty`. The
	// OPTCHAIN_GETELEM is handled consistently with regular GETELEM.
	c.Test(
		"var foo = {myProp: 0}; f(foo?.[JSCompiler_renameProperty('myProp')]);",
		"var foo = {a: 0}; f(foo?.['a']);")
}

// Java: RenamePropertiesTest.testRenamePropertiesFunctionCall2()
func (c *RenamePropertiesTest) TestRenamePropertiesFunctionCall2() {
	c.Test(
		`var foo = {myProp: 0};
f(JSCompiler_renameProperty('otherProp.myProp.someProp'));
foo.myProp = 1; foo.theirProp = 2; foo.yourProp = 3;
`,
		`var foo = {a: 0}; f('b.a.c');
foo.a = 1; foo.d = 2; foo.e = 3;
`)
}

// Java: RenamePropertiesTest.testRemoveRenameFunctionStubs1()
func (c *RenamePropertiesTest) TestRemoveRenameFunctionStubs1() {
	c.Test("function JSCompiler_renameProperty(x) { return x; }", "")
}

// Java: RenamePropertiesTest.testRemoveRenameFunctionStubs2()
func (c *RenamePropertiesTest) TestRemoveRenameFunctionStubs2() {
	c.Test(
		`function JSCompiler_renameProperty(x) { return x; }
var foo = {myProp: 0}; f(foo[JSCompiler_renameProperty('myProp')]);
`,
		"var foo = {a: 0}; f(foo['a']);")
}

// Java: RenamePropertiesTest.testGeneratePseudoNames()
func (c *RenamePropertiesTest) TestGeneratePseudoNames() {
	c.generatePseudoNames = true
	c.Test("var foo={}; foo.bar=1; foo['abc']=2", "var foo={}; foo.$bar$=1; foo['abc']=2")
}

// Java: RenamePropertiesTest.testChunks()
func (c *RenamePropertiesTest) TestChunks() {
	chunk1Js :=
		`function Bar(){} Bar.prototype.getA=function(x){};
var foo;foo.getA(foo);foo.doo=foo;foo.bloo=foo;
`

	chunk2Js :=
		`function Far(){} Far.prototype.getB=function(y){};
var too;too.getB(too);too.woo=too;too.bloo=too;
`

	chunk3Js :=
		`function Car(){} Car.prototype.getC=function(z){};
var noo;noo.getC(noo);noo.zoo=noo;noo.cloo=noo;
`

	chunk1 := NewJSChunk("m1")
	chunk1.AddSourceFile(SourceFile_FromCode("input1", chunk1Js))

	chunk2 := NewJSChunk("m2")
	chunk2.AddSourceFile(SourceFile_FromCode("input2", chunk2Js))

	chunk3 := NewJSChunk("m3")
	chunk3.AddSourceFile(SourceFile_FromCode("input3", chunk3Js))

	chunks := []*JSChunk{chunk1, chunk2, chunk3}
	compiler := c.CompileChunks("", chunks)

	result := compiler.GetResult()
	AssertThat(result.Success).IsTrue()

	AssertThat(compiler.ToSourceOfJSChunk(chunk1)).
		IsEqualTo(
			"function Bar(){}Bar.prototype.b=function(x){};var foo;foo.b(foo);foo.f=foo;foo.a=foo;")

	AssertThat(compiler.ToSourceOfJSChunk(chunk2)).
		IsEqualTo(
			"function Far(){}Far.prototype.c=function(y){};var too;too.c(too);too.g=too;too.a=too;")

	// Note that properties that occur most often globally get the earliest
	// names. The "getC" property, which doesn't occur until chunk 3, is
	// renamed to an earlier name in the alphabet than "woo", which appears
	// in chunk 2, because "getC" occurs more total times across all chunks.
	// Might be better to give early chunks the shortest names, but this is
	// how the pass currently works.
	AssertThat(compiler.ToSourceOfJSChunk(chunk3)).
		IsEqualTo(
			"function Car(){}Car.prototype.d=function(z){};var noo;noo.d(noo);noo.h=noo;noo.e=noo;")
}

// Java: RenamePropertiesTest.testPropertyAffinityOff()
func (c *RenamePropertiesTest) TestPropertyAffinityOff() {
	c.Test(
		`var foo={};foo.x=1;foo.y=2;foo.z=3;
function f1() { foo.z; foo.z; foo.z; foo.y}
function f2() {                      foo.x}
`,
		`var foo={};foo.b=1;foo.c=2;foo.a=3;
function f1() { foo.a; foo.a; foo.a; foo.c}
function f2() {                      foo.b}
`)

	c.Test(
		`var foo={};foo.x=1;foo.y=2;foo.z=3;
function f1() { foo.z; foo.z; foo.z; foo.y}
function f2() { foo.z; foo.z; foo.z; foo.x}
`,
		`var foo={};foo.b=1;foo.c=2;foo.a=3;
function f1() { foo.a; foo.a; foo.a; foo.c}
function f2() { foo.a; foo.a; foo.a; foo.b}
`)
}

// Java: RenamePropertiesTest.testPrototypePropertiesStable()
func (c *RenamePropertiesTest) TestPrototypePropertiesStable() {
	c.TestStableRenaming(
		`Bar.prototype.getA = function(){}; bar.getA();
Bar.prototype.getB = function(){};
`,
		`Bar.prototype.a = function(){}; bar.a();
Bar.prototype.b = function(){}
`,
		`Bar.prototype.get = function(){}; bar.get();
Bar.prototype.getA = function(){}; bar.getA();
Bar.prototype.getB = function(){};
`,
		`Bar.prototype.c = function(){}; bar.c();
Bar.prototype.a = function(){}; bar.a();
Bar.prototype.b = function(){}
`)
}

// Java: RenamePropertiesTest.testPrototypePropertiesAsObjLitKeysStable()
func (c *RenamePropertiesTest) TestPrototypePropertiesAsObjLitKeysStable() {
	c.TestStableRenaming(
		"Bar.prototype = {2: function(){}, getA: function(){}}; bar[2]();",
		"Bar.prototype = {2: function(){}, a: function(){}}; bar[2]();",
		"Bar.prototype = {getB: function(){},getA: function(){}}; bar.getB();",
		"Bar.prototype = {b: function(){},a: function(){}}; bar.b();")
}

// Java: RenamePropertiesTest.testMixedQuotedAndUnquotedObjLitKeysStable()
func (c *RenamePropertiesTest) TestMixedQuotedAndUnquotedObjLitKeysStable() {
	c.TestStableRenaming(
		"Bar = {getA: function(){}, 'getB': function(){}}; bar.getA();",
		"Bar = {a: function(){}, 'getB': function(){}}; bar.a();",
		`Bar = {get: function(){}, getA: function(){}, 'getB': function(){}};
bar.getA();bar.get();
`,
		`Bar = {b: function(){}, a: function(){}, 'getB': function(){}};
bar.a();bar.b();
`)
}

// Java: RenamePropertiesTest.testOverlappingOriginalAndGeneratedNamesStable()
func (c *RenamePropertiesTest) TestOverlappingOriginalAndGeneratedNamesStable() {
	c.TestStableRenaming(
		"Bar.prototype = {b: function(){}, a: function(){}}; bar.b();",
		"Bar.prototype = {a: function(){}, b: function(){}}; bar.a();",
		`Bar.prototype = {c: function(){}, b: function(){}, a: function(){}};
bar.b();
`,
		`Bar.prototype = {c: function(){}, a: function(){}, b: function(){}};
bar.a();
`)
}

// Java: RenamePropertiesTest.testStableWithTrickyExternsChanges()
func (c *RenamePropertiesTest) TestStableWithTrickyExternsChanges() {
	c.Test(
		"Bar.prototype = {b: function(){}, a: function(){}}; bar.b();",
		"Bar.prototype = {a: function(){}, b: function(){}}; bar.a();")
	c.prevUsedPropertyMap = c.renameProperties.GetPropertyMap()
	externs := renamePropertiesTest_EXTERNS + "prop.b;"
	c.Test(
		c.Externs(externs),
		c.Srcs(
			`Bar.prototype = {new_f: function(){}, b: function(){},
a: function(){}};bar.b();
`),
		c.Expected(
			`Bar.prototype = {c:function(){}, b:function(){}, a:function(){}};
bar.b();
`))
}

// Java: RenamePropertiesTest.testRenamePropertiesWithLeadingUnderscoresStable()
func (c *RenamePropertiesTest) TestRenamePropertiesWithLeadingUnderscoresStable() {
	c.TestStableRenaming(
		"Bar.prototype = {_getA: function(){}, _b: 0}; bar._getA();",
		"Bar.prototype = {a: function(){}, b: 0}; bar.a();",
		"Bar.prototype = {_getA: function(){}, _c: 1, _b: 0}; bar._getA();",
		"Bar.prototype = {a: function(){}, c: 1, b: 0}; bar.a();")
}

// Java: RenamePropertiesTest.testPropertyAddedToObjectStable()
func (c *RenamePropertiesTest) TestPropertyAddedToObjectStable() {
	c.TestStableRenaming(
		"var foo = {}; foo.prop = '';",
		"var foo = {}; foo.a = '';",
		"var foo = {}; foo.prop = ''; foo.a='';",
		"var foo = {}; foo.a = ''; foo.b='';")
}

// Java: RenamePropertiesTest.testAvoidingConflictsBetQuotedAndUnquotedPropertyNamesStable()
func (c *RenamePropertiesTest) TestAvoidingConflictsBetQuotedAndUnquotedPropertyNamesStable() {
	c.TestStableRenaming(
		"Bar.prototype.foo = function(){}; Bar.prototype['b'] = 0; bar.foo();",
		"Bar.prototype.a = function(){}; Bar.prototype['b'] = 0; bar.a();",
		"Bar.prototype.foo = function(){}; Bar.prototype['a'] = 0; bar.foo();",
		"Bar.prototype.b = function(){}; Bar.prototype['a'] = 0; bar.b();")
}

// Java: RenamePropertiesTest.testRenamePropertiesFunctionCallStable()
func (c *RenamePropertiesTest) TestRenamePropertiesFunctionCallStable() {
	c.TestStableRenaming(
		`var foo = {myProp: 0};
f(JSCompiler_renameProperty('otherProp.myProp.someProp'));
foo.myProp = 1; foo.theirProp = 2; foo.yourProp = 3;
`,
		`var foo = {a: 0}; f('b.a.c');
foo.a = 1; foo.d = 2; foo.e = 3;
`,
		`var bar = {newProp: 0}; var foo = {myProp: 0};
f(JSCompiler_renameProperty('otherProp.myProp.someProp'));
foo.myProp = 1; foo.theirProp = 2; foo.yourProp = 3;
`,
		`var bar = {f: 0}; var foo = {a: 0}; f('b.a.c');
foo.a = 1; foo.d = 2; foo.e = 3;
`)
}

// Java: RenamePropertiesTest.testStableRenaming(String, String, String, String)
func (c *RenamePropertiesTest) TestStableRenaming(input1 string, expected1 string, input2 string, expected2 string) {
	c.Test(input1, expected1)
	c.prevUsedPropertyMap = c.renameProperties.GetPropertyMap()
	c.Test(input2, expected2)
}

// Java: RenamePropertiesTest.testPrototypePropertyForArrowFunction()
func (c *RenamePropertiesTest) TestPrototypePropertyForArrowFunction() {
	c.Test(
		"Bar.prototype = {2: () => {}, getA: () => {}}; bar[2]();",
		"Bar.prototype = {2: () => {}, a:    () => {}}; bar[2]();")
}

// Java: RenamePropertiesTest.testArrayDestructuring()
func (c *RenamePropertiesTest) TestArrayDestructuring() {
	c.TestSame("var [first, second] = someArray")
}

// Java: RenamePropertiesTest.testDestructuredProperties()
func (c *RenamePropertiesTest) TestDestructuredProperties() {
	// using destructuring shorthand
	c.Test("var {   foo,   bar } = { foo: 1, bar: 2 }", "var { b:foo, a:bar } = {    b:1,    a:2 }")

	// without destructuring shorthand
	c.Test(
		"var { foo:foo, bar:bar } = { foo:1, bar:2 }",
		"var {   b:foo,   a:bar } = {   b:1,   a:2 }")

	c.Test(
		"var foo = { bar: 1, baz: 2 }; var foo1 = foo.bar; var foo2 = foo.baz; ",
		"var foo = {   a: 1,   b: 2 }; var foo1 = foo.a;   var foo2 = foo.b;")
}

// Java: RenamePropertiesTest.testNestedDestructuringProperties()
func (c *RenamePropertiesTest) TestNestedDestructuringProperties() {
	c.Test(
		"var {outer: {inner}} = {outer: {inner: 'value'}};",
		"var {b: {a: inner}} = {b: {a: 'value'}};")
}

// Java: RenamePropertiesTest.testComputedPropertyNamesInObjectLit()
func (c *RenamePropertiesTest) TestComputedPropertyNamesInObjectLit() {
	// TODO (simranarora) A restriction of this pass is that quoted and unquoted property
	// references cannot be mixed.
	c.Test(
		`var a = {
  ['val' + ++i]: i,
  ['val' + ++i]: i
};
a.val1;
`,
		`var a = {
  ['val' + ++i]: i, // don't rename here
  ['val' + ++i]: i
};
a.a;
`) // rename here
}

// Java: RenamePropertiesTest.testComputedMethodPropertyNamesInClass()
func (c *RenamePropertiesTest) TestComputedMethodPropertyNamesInClass() {
	// TODO (simranarora) A restriction of this pass is that quoted and unquoted property
	// references cannot be mixed.

	// Concatination for computed property
	c.Test(
		`class Bar {
  constructor(){}
  ['f'+'oo']() {
    return 1
  }
}
var bar = new Bar()
bar.foo();
`,
		`class Bar {
  constructor(){}
  ['f'+'oo']() { // don't rename here
    return 1
  }
}
var bar = new Bar()
bar.a();
`) // rename here

	// Without property concatination
	c.Test(
		`class Bar {
  constructor(){}
  ['foo']() {
    return 1
  }
}
var bar = new Bar()
bar.foo();
`,
		`class Bar {
  constructor(){}
  ['foo']() { // don't rename here
    return 1
  }
}
var bar = new Bar()
bar.a();
`) // rename here
}

// Java: RenamePropertiesTest.testClasses()
func (c *RenamePropertiesTest) TestClasses() {
	// Call class method inside class scope - due to the scoping rules of javascript, the "getA()"
	// call inside of getB() refers to a method getA() in the outer scope and not the getA() method
	// inside the Bar class
	c.Test(
		`function getA() {};
class Bar {
  constructor(){}
  getA() {
    return 1
  }
  getB(x) {
    getA();
  }
}
`,
		`function getA() {};
class Bar {
  constructor(){}
  a() {
    return 1
  }
  b(x) {
    getA();
  }
}
`)

	// Call class method inside class scope - due to the scoping rules of javascript,
	// the "this.getA()" call inside of getB() refers to a method getA() in the Bar class and
	// not the getA() method in the outer scope
	c.Test(
		`function getA() {};
class Bar {
  constructor(){}
  getA() {
    return 1
  }
  getB(x) {
    this.getA();
  }
}
`,
		`function getA() {};
class Bar {
  constructor(){}
  a() {
    return 1
  }
  b(x) {
    this.a();
  }
}
`)

	// Call class method outside class scope
	c.Test(
		`class Bar {
  constructor(){}
  getB(x) {}
}
var too;
var too = new Bar();
too.getB(too);
`,
		`class Bar {
  constructor(){}
  a(x) {}
}
var too;
var too = new Bar();
too.a(too);
`)
}

// Java: RenamePropertiesTest.testGetSetInClass()
func (c *RenamePropertiesTest) TestGetSetInClass() {
	c.Test(
		`class Bar {
  constructor(foo){
    this.foo = foo;
  }
  get foo() {
    return this.foo;
  }
  set foo(x) {
    this.foo = x;
  }
}
var barObj = new Bar();
barObj.foo();
barObj.foo(1);
`,
		`class Bar {
  constructor(foo){
    this.a = foo;
  }
  get a() {
    return this.a;
  }
  set a(x) {
    this.a = x;
  }
}
var barObj = new Bar();
barObj.a();
barObj.a(1);
`)
}

// Java: RenamePropertiesTest.testStaticMethodInClass()
func (c *RenamePropertiesTest) TestStaticMethodInClass() {
	c.Test(
		`class Bar {
  static double(n) {
    return n*2
  }
}
Bar.double(1);
`,
		`class Bar {
  static a(n) {
    return n*2
  }
}
Bar.a(1);
`)
}

// Java: RenamePropertiesTest.testClassFields()
func (c *RenamePropertiesTest) TestClassFields() {
	c.Test(
		`class Bar {
  field = 7;
}
var bar = new Bar();
bar.field;
`,
		`class Bar {
  a = 7;
}
var bar = new Bar();
bar.a;
`)
}

// Java: RenamePropertiesTest.testClassFieldsInExterns()
func (c *RenamePropertiesTest) TestClassFieldsInExterns() {
	externs :=
		`class Foo {
  helloWorld;
  get bar() {}
  set baz(v) {}
}
`
	c.TestSame(c.Externs(externs), c.Srcs("new Foo().helloWorld; new Foo().bar; new Foo().baz = 1;"))
}

// Java: RenamePropertiesTest.testClassFieldWithFunctionRHS()
func (c *RenamePropertiesTest) TestClassFieldWithFunctionRHS() {
	c.Test(
		`class Bar {
  superClass_ = function f(){};
}
var bar = new Bar();
bar.superClass_;
`,
		`class Bar {
  a = function f(){};
}
var bar = new Bar();
bar.a;
`)
}

// Java: RenamePropertiesTest.testStaticClassFields()
func (c *RenamePropertiesTest) TestStaticClassFields() {
	c.Test(
		`class Bar {
  static field = 1;
}
Bar.field;
`,
		`class Bar {
  static a = 1;
}
Bar.a;
`)
}

// Java: RenamePropertiesTest.testClassComputedFields()
func (c *RenamePropertiesTest) TestClassComputedFields() {
	c.Test(
		`class Bar {
  ['field'] = 1;
}
var bar = new Bar()
bar.field;
`,
		`class Bar {
  ['field'] = 1;
}
var bar = new Bar()
bar.a;
`)
}

// Java: RenamePropertiesTest.testStaticClassComputedFields()
func (c *RenamePropertiesTest) TestStaticClassComputedFields() {
	c.Test(
		`class Bar {
  static ['field'] = 1;
}
Bar.field;
`,
		`class Bar {
  static ['field'] = 1;
}
Bar.a;
`)
}

// Java: RenamePropertiesTest.testClassMixedFields()
func (c *RenamePropertiesTest) TestClassMixedFields() {
	c.Test(
		`class Bar {
  field = 1;
  ['field'] = 2;
  static 1 = 5;
}
var bar = new Bar()
bar.field;
`,
		`class Bar {
  a = 1;
  ['field'] = 2;
  static 1 = 5;
}
var bar = new Bar()
bar.a;
`)
}

// Java: RenamePropertiesTest.testObjectMethodProperty()
func (c *RenamePropertiesTest) TestObjectMethodProperty() {
	// ES5 version
	c.Test(
		`var foo = {
  bar: 1,
  myFunc: function myFunc() {
    return this.bar
  }
};
foo.myFunc();
`,
		`var foo = {
  a: 1,
  b: function myFunc() {
    return this.a
  }
};
foo.b();
`)

	// ES6 version
	c.Test(
		`var foo = {
  bar: 1,
  myFunc() {
    return this.bar
  }
};
foo.myFunc();
`,
		`var foo = {
  a: 1,
  b() {
    return this.a
  }
};
foo.b();
`)
}

// Java: RenamePropertiesTest.testPrototypePropertiesUsingFilterFunction()
func (c *RenamePropertiesTest) TestPrototypePropertiesUsingFilterFunction() {
	c.propertyNameFilter =
		func(node rhino.Node) bool {
			name := node.GetString()
			// Strings.nullToEmpty(name): getString() is never null.
			return strings.HasPrefix(name, "get")
		}

	c.Test(
		`Bar.prototype.getA = function(){}; bar.getA();
Bar.prototype.getB = function(){};
Bar.prototype.setC = function(){};
`,
		`Bar.prototype.a = function(){}; bar.a();
Bar.prototype.b = function(){};
Bar.prototype.setC = function(){}
`)
}

// Java: RenamePropertiesTest.testPrototypePropertiesUsingFilterOnFilename()
func (c *RenamePropertiesTest) TestPrototypePropertiesUsingFilterOnFilename() {
	c.propertyNameFilter =
		func(node rhino.Node) bool {
			name := ""
			if sourceFileName := node.GetSourceFileName(); sourceFileName != nil { // Strings.nullToEmpty
				name = *sourceFileName
			}
			return name == "foo.js"
		}

	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				"foo.js",
				`Foo.prototype.getA = function(){};
Foo.prototype.getB = function(){};
foo.getA();
foo.getB();
`),
			SourceFile_FromCode(
				"bar.js",
				`Bar.prototype.getA = function(){};
bar.getA();
`)),
		c.Expected(
			SourceFile_FromCode(
				"foo.js",
				`Foo.prototype.getA = function(){};
Foo.prototype.b = function(){};
foo.getA();
foo.b();
`),
			SourceFile_FromCode(
				"bar.js",
				`Bar.prototype.getA = function(){};
bar.getA();
`)))
}

// Java: RenamePropertiesTest.testPrototypeAndRenameFunctionsFilterOnFilename()
func (c *RenamePropertiesTest) TestPrototypeAndRenameFunctionsFilterOnFilename() {
	c.propertyNameFilter =
		func(node rhino.Node) bool {
			name := ""
			if sourceFileName := node.GetSourceFileName(); sourceFileName != nil { // Strings.nullToEmpty
				name = *sourceFileName
			}
			return name == "foo.js"
		}

	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				"foo.js",
				`var foo = {myProp: 0, myProp2: 1};
f(foo[JSCompiler_renameProperty('myProp')]);
f(foo[JSCompiler_renameProperty('myProp2')]);
`),
			SourceFile_FromCode(
				"bar.js",
				`var bar = {myProp: 0};
f(bar[JSCompiler_renameProperty('myProp')]);
`)),
		c.Expected(
			SourceFile_FromCode(
				"foo.js",
				`var foo = {myProp: 0, b: 1};
f(foo['myProp']);
f(foo['b']);
`),
			SourceFile_FromCode(
				"bar.js",
				`var bar = {myProp: 0};
f(bar['myProp']);
`)))
}

// Java: RenamePropertiesTest.testPrototypeAndRenameFunctionsFilterOnFilename2()
func (c *RenamePropertiesTest) TestPrototypeAndRenameFunctionsFilterOnFilename2() {
	c.propertyNameFilter =
		func(node rhino.Node) bool {
			name := ""
			if sourceFileName := node.GetSourceFileName(); sourceFileName != nil { // Strings.nullToEmpty
				name = *sourceFileName
			}
			return name == "foo.js"
		}

	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				"foo.js",
				`var foo = {myProp: 0, myProp2: 1};
f(JSCompiler_renameProperty('otherProp.myProp.someProp'));
f(JSCompiler_renameProperty('otherProp.myProp2.someProp'));
`),
			SourceFile_FromCode(
				"bar.js",
				`var bar = {myProp: 0};
f(JSCompiler_renameProperty('otherProp.myProp.someProp'));
`)),
		c.Expected(
			SourceFile_FromCode(
				"foo.js",
				`var foo = {myProp: 0, d: 1};
f('otherProp.myProp.someProp');
f('otherProp.d.someProp');
`),
			SourceFile_FromCode(
				"bar.js",
				`var bar = {myProp: 0};
f('otherProp.myProp.someProp');
`)))
}

// Java: RenamePropertiesTest.compileChunks(String, JSChunk[])
func (c *RenamePropertiesTest) CompileChunks(externs string, chunks []*JSChunk) *Compiler {
	externsInput := SourceFile_FromCode("externs", externs)

	options := NewCompilerOptions()
	options.SetEmitUseStrict(false)
	options.SetPropertyRenaming(PropertyRenamingPolicy_ALL_UNQUOTED)

	compiler := NewCompiler()
	_ = compiler.CompileChunks(
		[]*SourceFile{externsInput}, append([]*JSChunk(nil), chunks...), options)
	return compiler
}

// Java: RenamePropertiesTest.getProcessor(Compiler)
func (c *RenamePropertiesTest) GetProcessor(compiler *Compiler) CompilerPass {
	c.renameProperties = NewRenamePropertiesWithFilter(
		compiler,
		c.generatePseudoNames,
		c.prevUsedPropertyMap,
		ASCIICharSet{},
		ASCIICharSet{},
		NewDefaultNameGenerator(),
		c.propertyNameFilter)
	return c.renameProperties
}

// Java: RenamePropertiesTest.testAvoidingConflictsBetweenQuotedAndUnquotedClassFields()
func (c *RenamePropertiesTest) TestAvoidingConflictsBetweenQuotedAndUnquotedClassFields() {
	// TODO(b/571953710): COMPUTED_FIELD_DEF string keys not reserved
	c.Test(
		`class C {
  foo = 1;
  'a' = alert(2);
}
new C().foo;
`,
		`class C {
  a = 1;
  ["a"] = alert(2);
}
(new C()).a;
`)
}
