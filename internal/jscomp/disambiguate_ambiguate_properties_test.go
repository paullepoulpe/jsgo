// Port of com/google/javascript/jscomp/disambiguate/AmbiguatePropertiesTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/disambiguate/AmbiguatePropertiesTest.java)

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
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// AmbiguatePropertiesTest is the unit test for AmbiguateProperties.
//
// Java: public final class AmbiguatePropertiesTest extends CompilerTestCase
type AmbiguatePropertiesTest struct {
	CompilerTestCaseBase

	lastPass *AmbiguateProperties
}

// Java: AmbiguatePropertiesTest.EXTERNS
var ambiguatePropertiesTest_EXTERNS = CompilerTestCase_MINIMAL_EXTERNS +
	`Function.prototype.call=function(){};
Function.prototype.inherits=function(){};
Object.defineProperties = function(typeRef, definitions) {};
Object.prototype.toString = function() {};
var google = { gears: { factory: {}, workerPool: {} } };
/** @return {?} */ function any() {};
/** @constructor */ function Window() {}
/** @const */ var window = {};
`

// Java: AmbiguatePropertiesTest()
func NewAmbiguatePropertiesTest() *AmbiguatePropertiesTest {
	c := &AmbiguatePropertiesTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, ambiguatePropertiesTest_EXTERNS)
	return c
}

func TestAmbiguatePropertiesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewAmbiguatePropertiesTest() })
}

// JUnitBefore runs the @Before methods: CompilerTestCase.setUp, then
// customSetUp.
func (c *AmbiguatePropertiesTest) JUnitBefore() {
	c.CompilerTestCaseBase.JUnitBefore()
	c.CustomSetUp()
}

// Java: AmbiguatePropertiesTest.customSetUp()
func (c *AmbiguatePropertiesTest) CustomSetUp() {
	c.EnableTypeCheck()
	c.ReplaceTypesWithColors()
	c.EnableNormalize()
	c.EnableGatherExternProperties()
	c.DisableCompareJsDoc() // removeTypes also deletes JSDocInfo
}

// Java: AmbiguatePropertiesTest.getProcessor(Compiler)
func (c *AmbiguatePropertiesTest) GetProcessor(compiler *Compiler) CompilerPass {
	return CompilerPassFunc(func(externs *rhino.Node, root *rhino.Node) {
		c.lastPass = AmbiguateProperties_MakePassForTesting(
			compiler,
			jutil.NewLinkedHashSet[uint16]('$'),
			jutil.NewLinkedHashSet[uint16]('$'),
			c.GetGatheredExternProperties())
		c.lastPass.Process(externs, root)
	})
}

// Java: AmbiguatePropertiesTest.testOneVar1()
func (c *AmbiguatePropertiesTest) TestOneVar1() {
	c.Test(
		"/** @constructor */ var Foo = function(){};Foo.prototype.b = 0;",
		"/** @constructor */ var Foo = function(){};Foo.prototype.a = 0;")
}

// Java: AmbiguatePropertiesTest.testOneVar2()
func (c *AmbiguatePropertiesTest) TestOneVar2() {
	c.Test(
		`/** @constructor */ var Foo = function(){};
Foo.prototype = {b: 0};
`,
		`/** @constructor */ var Foo = function(){};
Foo.prototype = {a: 0};
`)
}

// Java: AmbiguatePropertiesTest.testOneVar3()
func (c *AmbiguatePropertiesTest) TestOneVar3() {
	c.Test(
		`/** @constructor */ var Foo = function(){};
Foo.prototype = {get b() {return 0}};
`,
		`/** @constructor */ var Foo = function(){};
Foo.prototype = {get a() {return 0}};
`)
}

// Java: AmbiguatePropertiesTest.testOneVar4()
func (c *AmbiguatePropertiesTest) TestOneVar4() {
	c.Test(
		`/** @constructor */ var Foo = function(){};
Foo.prototype = {set b(a) {}};
`,
		`/** @constructor */ var Foo = function(){};
Foo.prototype = {set a(a) {}};
`)
}

// Java: AmbiguatePropertiesTest.testTwoVar1()
func (c *AmbiguatePropertiesTest) TestTwoVar1() {
	js :=
		`/** @constructor */ var Foo = function(){};
Foo.prototype.z=0;
Foo.prototype.z=0;
Foo.prototype.x=0;
`
	output :=
		`/** @constructor */ var Foo = function(){};
Foo.prototype.a=0;
Foo.prototype.a=0;
Foo.prototype.b=0;
`
	c.Test(js, output)
}

// Java: AmbiguatePropertiesTest.testTwoVar2()
func (c *AmbiguatePropertiesTest) TestTwoVar2() {
	js :=
		`/** @constructor */ var Foo = function(){};
Foo.prototype = {z:0, z:1, x:0};
`
	output :=
		`/** @constructor */ var Foo = function(){};
Foo.prototype = {a:0, a:1, b:0};
`
	c.Test(js, output)
}

// Java: AmbiguatePropertiesTest.testTwoIndependentVar()
func (c *AmbiguatePropertiesTest) TestTwoIndependentVar() {
	js :=
		`/** @constructor */ var Foo = function(){};
Foo.prototype.b = 0;
/** @constructor */ var Bar = function(){};
Bar.prototype.c = 0;
`
	output :=
		`/** @constructor */ var Foo = function(){};
Foo.prototype.a=0;
/** @constructor */ var Bar = function(){};
Bar.prototype.a=0;
`
	c.Test(js, output)
}

// Java: AmbiguatePropertiesTest.testTwoTypesTwoVar()
func (c *AmbiguatePropertiesTest) TestTwoTypesTwoVar() {
	js :=
		`/** @constructor */ var Foo = function(){};
Foo.prototype.r = 0;
Foo.prototype.g = 0;
/** @constructor */ var Bar = function(){};
Bar.prototype.c = 0;
Bar.prototype.r = 0;
`
	output :=
		`/** @constructor */ var Foo = function(){};
Foo.prototype.a=0;
Foo.prototype.b=0;
/** @constructor */ var Bar = function(){};
Bar.prototype.b=0;
Bar.prototype.a=0;
`
	c.Test(js, output)
}

// Java: AmbiguatePropertiesTest.testUnion_withUnrelatedPropertyAccess()
func (c *AmbiguatePropertiesTest) TestUnion_withUnrelatedPropertyAccess() {
	js :=
		`/** @constructor */ var Foo = function(){};
/** @constructor */ var Bar = function(){};
Foo.prototype.foodoo=0;
Bar.prototype.bardoo=0;
// variable exists that could be a Foo or a Bar
/** @type {Foo|Bar} */
var U = any();
// We don't actually access either foodoo or bardoo on it,
// though, so it's OK if they end up having the same name
U.joint
`
	output :=
		`/** @constructor */ var Foo = function(){};
/** @constructor */ var Bar = function(){};
Foo.prototype.a=0;
Bar.prototype.a=0;
/** @type {Foo|Bar} */
var U = any();
U.b
`
	c.Test(js, output)
}

// Java: AmbiguatePropertiesTest.testUnion_withRelatedPropertyAccess()
func (c *AmbiguatePropertiesTest) TestUnion_withRelatedPropertyAccess() {
	js :=
		`/** @constructor */ var Foo = function(){};
/** @constructor */ var Bar = function(){};
Foo.prototype.foodoo=0;
Bar.prototype.bardoo=0;
// variable exists that could be a Foo or a Bar
/** @type {Foo|Bar} */
var U = any();
// both foodoo and bardoo are accessed through that variable,
// so they must have different names.
U.foodoo;
U.bardoo
`
	output :=
		`/** @constructor */ var Foo = function(){};
/** @constructor */ var Bar = function(){};
Foo.prototype.b=0;
Bar.prototype.a=0;
/** @type {Foo|Bar} */
var U = any();
U.b;
U.a
`
	c.Test(js, output)
}

// Java: AmbiguatePropertiesTest.testUnions()
func (c *AmbiguatePropertiesTest) TestUnions() {
	js :=
		`/** @constructor */ var Foo = function(){};
/** @constructor */ var Bar = function(){};
/** @constructor */ var Baz = function(){};
/** @constructor */ var Bat = function(){};
Foo.prototype.lone1=0;
Bar.prototype.lone2=0;
Baz.prototype.lone3=0;
Bat.prototype.lone4=0;
/** @type {Foo|Bar} */
var U1 = any();
U1.j1;
U1.j2;
/** @type {Baz|Bar} */
var U2 = any();
U2.j3;
U2.j4;
/** @type {Baz|Bat} */
var U3 = any();
U3.j5;
U3.j6
`
	output :=
		`/** @constructor */ var Foo = function(){};
/** @constructor */ var Bar = function(){};
/** @constructor */ var Baz = function(){};
/** @constructor */ var Bat = function(){};
Foo.prototype.c=0;
Bar.prototype.e=0;
Baz.prototype.e=0;
Bat.prototype.c=0;
/** @type {Foo|Bar} */
var U1 = any();
U1.a;
U1.b;
/** @type {Baz|Bar} */
var U2 = any();
U2.c;
U2.d;
/** @type {Baz|Bat} */
var U3 = any();
U3.a;
U3.b
`
	c.Test(js, output)
}

// Java: AmbiguatePropertiesTest.testExtends()
func (c *AmbiguatePropertiesTest) TestExtends() {
	c.EnableClosurePass()

	js :=
		`/** @constructor */ var Foo = function(){};
Foo.prototype.x=0;
/** @constructor
 @extends Foo */ var Bar = function(){};
goog.inherits(Bar, Foo);
Bar.prototype.y=0;
Bar.prototype.z=0;
/** @constructor */ var Baz = function(){};
Baz.prototype.l=0;
Baz.prototype.m=0;
Baz.prototype.n=0;
(new Baz).m
`
	output :=
		`/** @constructor */ var Foo = function(){};
Foo.prototype.a=0;
/** @constructor
 @extends Foo */ var Bar = function(){};
goog.inherits(Bar, Foo);
Bar.prototype.b=0;
Bar.prototype.c=0;
/** @constructor */ var Baz = function(){};
Baz.prototype.b=0;
Baz.prototype.a=0;
Baz.prototype.c=0;
(new Baz).a
`
	c.Test(js, output)
}

// Java: AmbiguatePropertiesTest.testLotsOfVars()
func (c *AmbiguatePropertiesTest) TestLotsOfVars() {
	var js strings.Builder
	var output strings.Builder
	js.WriteString("/** @constructor */ var Foo = function(){};\n")
	js.WriteString("/** @constructor */ var Bar = function(){};\n")
	output.WriteString(js.String())

	vars := 10
	for i := 0; i < vars; i++ {
		js.WriteString("Foo.prototype.var" + strconv.Itoa(i) + " = 0;")
		js.WriteString("Bar.prototype.var" + strconv.Itoa(i+10000) + " = 0;")
		output.WriteString("Foo.prototype." + string(rune('a'+i)) + "=0;")
		output.WriteString("Bar.prototype." + string(rune('a'+i)) + "=0;")
	}
	c.Test(js.String(), output.String())
}

// Java: AmbiguatePropertiesTest.testLotsOfClasses()
func (c *AmbiguatePropertiesTest) TestLotsOfClasses() {
	var b strings.Builder
	classes := 10
	for i := 0; i < classes; i++ {
		c_ := "Foo" + strconv.Itoa(i)
		b.WriteString("/** @constructor */ var " + c_ + " = function(){};\n")
		b.WriteString(c_ + ".prototype.varness" + strconv.Itoa(i) + " = 0;")
	}
	js := b.String()
	c.Test(js, jutil.StringReplaceAll(js, "varness\\d+", "a"))
}

// Java: AmbiguatePropertiesTest.testFunctionType()
func (c *AmbiguatePropertiesTest) TestFunctionType() {
	js :=
		`/** @constructor */ function Foo(){};
/** @return {Bar} */
Foo.prototype.fun = function() { return new Bar(); };
/** @constructor */ function Bar(){};
Bar.prototype.bazz;
(new Foo).fun().bazz();
`
	output :=
		`/** @constructor */ function Foo(){};
/** @return {Bar} */
Foo.prototype.a = function() { return new Bar(); };
/** @constructor */ function Bar(){};
Bar.prototype.a;
(new Foo).a().a();
`
	c.Test(js, output)
}

// Java: AmbiguatePropertiesTest.testPrototypePropertiesAsObjLitKeys1()
func (c *AmbiguatePropertiesTest) TestPrototypePropertiesAsObjLitKeys1() {
	c.Test(
		`/** @constructor */ function Bar() {};
Bar.prototype = {2: function(){}, getA: function(){}};
`,
		`/** @constructor */ function Bar() {};
Bar.prototype = {2: function(){}, a: function(){}};
`)
}

// Java: AmbiguatePropertiesTest.testPrototypePropertiesAsObjLitKeys2()
func (c *AmbiguatePropertiesTest) TestPrototypePropertiesAsObjLitKeys2() {
	c.TestSame(
		`/** @constructor */ function Bar() {};
Bar.prototype = {2: function(){}, 'getA': function(){}};
`)
}

// Java: AmbiguatePropertiesTest.testQuotedPrototypeProperty()
func (c *AmbiguatePropertiesTest) TestQuotedPrototypeProperty() {
	c.TestSame(
		`/** @constructor */ function Bar() {};
Bar.prototype['getA'] = function(){};
var bar = new Bar();
bar['getA']();
`)
}

// Java: AmbiguatePropertiesTest.testObjectDefineProperties()
func (c *AmbiguatePropertiesTest) TestObjectDefineProperties() {
	js :=
		`/** @constructor */ var Bar = function(){};
Bar.prototype.bar = 0;
/** @struct @constructor */ var Foo = function() {
  this.bar_ = 'bar';
};
/** @type {?} */ Foo.prototype.bar;
Object.defineProperties(Foo.prototype, {
  bar: {
    configurable: true,
    enumerable: true,
    /** @this {Foo} */ get: function() { return this.bar_;},
    /** @this {Foo} */ set: function(value$jscomp$1) { this.bar_ = value$jscomp$1; }
  }
});
`

	result :=
		`/** @constructor */ var Bar = function(){};
Bar.prototype.a = 0;
/** @struct @constructor */ var Foo = function() {
  this.b = 'bar';
};
/** @type {?} */ Foo.prototype.a;
Object.defineProperties(Foo.prototype, {
  a: {
    configurable: true,
    enumerable: true,
    /** @this {Foo} */ get: function() { return this.b;},
    /** @this {Foo} */ set: function(value$jscomp$1) { this.b = value$jscomp$1; }
  }
});
`

	c.Test(js, result)
}

// Java: AmbiguatePropertiesTest.testObjectDefinePropertiesQuoted()
func (c *AmbiguatePropertiesTest) TestObjectDefinePropertiesQuoted() {
	js :=
		`/** @constructor */ var Bar = function(){};
Bar.prototype.bar = 0;
/** @struct @constructor */ var Foo = function() {
  this.bar_ = 'bar';
};
/** @type {?} */ Foo.prototype['bar'];
Object.defineProperties(Foo.prototype, {
  'a': {
    configurable: true,
    enumerable: true,
    /** @this {Foo} */ get: function() { return this.bar_;},
    /** @this {Foo} */ set: function(value$jscomp$1) { this.bar_ = value$jscomp$1; }
  }
});
`

	result :=
		`/** @constructor */ var Bar = function(){};
Bar.prototype.b = 0;
/** @struct @constructor */ var Foo = function() {
  this.b = 'bar';
};
/** @type {?} */ Foo.prototype['bar'];
Object.defineProperties(Foo.prototype, {
  'a': {
    configurable: true,
    enumerable: true,
    /** @this {Foo} */ get: function() { return this.b;},
    /** @this {Foo} */ set: function(value$jscomp$1) { this.b = value$jscomp$1; }
  }
});
`

	c.Test(js, result)
}

// Java: AmbiguatePropertiesTest.testObjectDefinePropertiesComputed()
func (c *AmbiguatePropertiesTest) TestObjectDefinePropertiesComputed() {
	js :=
		`/** @constructor */ var Bar = function(){};
Bar.prototype.bar = 0;
/** @struct @constructor */ var Foo = function() {
  this.bar_ = 'bar';
};
/** @type {?} */ Foo.prototype['bar'];
Object.defineProperties(Foo.prototype, {
  ['a']: {
    configurable: true,
    enumerable: true,
    /** @this {Foo} */ get: function() { return this.bar_;},
    /** @this {Foo} */ set: function(value$jscomp$1) { this.bar_ = value$jscomp$1; }
  }
});
`

	result :=
		`/** @constructor */ var Bar = function(){};
Bar.prototype.b = 0;
/** @struct @constructor */ var Foo = function() {
  this.b = 'bar';
};
/** @type {?} */ Foo.prototype['bar'];
Object.defineProperties(Foo.prototype, {
  ['a']: {
    configurable: true,
    enumerable: true,
    /** @this {Foo} */ get: function() { return this.b;},
    /** @this {Foo} */ set: function(value$jscomp$1) { this.b = value$jscomp$1; }
  }
});
`

	c.Test(js, result)
}

// Java: AmbiguatePropertiesTest.testObjectDefinePropertiesMemberFn()
func (c *AmbiguatePropertiesTest) TestObjectDefinePropertiesMemberFn() {
	// NOTE: this is very odd code, and people should not be writing it anyway. We really just
	// don't want a crash.
	js :=
		`/** @constructor */ var Bar = function(){};
Bar.prototype.bar = 0;
/** @struct @constructor */ var Foo = function() {
  this.bar_ = 'bar';
};
/** @type {?} */ Foo.prototype['bar'];
Object.defineProperties(Foo.prototype, {
  a() {}
});
`

	result :=
		`/** @constructor */ var Bar = function(){};
Bar.prototype.a = 0;
/** @struct @constructor */ var Foo = function() {
  this.b = 'bar';
};
/** @type {?} */ Foo.prototype['bar'];
Object.defineProperties(Foo.prototype, {
  a() {}
});
`

	c.Test(js, result)
}

// Java: AmbiguatePropertiesTest.testOverlappingOriginalAndGeneratedNames()
func (c *AmbiguatePropertiesTest) TestOverlappingOriginalAndGeneratedNames() {
	c.Test(
		`/** @constructor */ function Bar(){};
Bar.prototype.b = function(){};
Bar.prototype.a = function(){};
var bar = new Bar();
bar.b();
`,
		`/** @constructor */ function Bar(){};
Bar.prototype.a = function(){};
Bar.prototype.b = function(){};
var bar = new Bar();
bar.a();
`)
}

// Java: AmbiguatePropertiesTest.testPropertyAddedToObject()
func (c *AmbiguatePropertiesTest) TestPropertyAddedToObject() {
	c.TestSame("var foo = {}; foo.prop = '';")
}

// Java: AmbiguatePropertiesTest.testPropertyAddedToFunction()
func (c *AmbiguatePropertiesTest) TestPropertyAddedToFunction() {
	c.TestSame("var foo = function(){}; foo.prop = '';")
}

// Java: AmbiguatePropertiesTest.testPropertyAddedToFunctionIndirectly()
func (c *AmbiguatePropertiesTest) TestPropertyAddedToFunctionIndirectly() {
	c.TestSame(
		`var foo = function(){}; foo.prop = ''; foo.baz = '';
function f(/** function(): void */ fun) { fun.bar = ''; fun.baz = ''; }
`)
}

// Java: AmbiguatePropertiesTest.testConstructorTreatedAsSubtypeOfFunction()
func (c *AmbiguatePropertiesTest) TestConstructorTreatedAsSubtypeOfFunction() {
	js :=
		`Function.prototype.a = 1;
/** @constructor */
function F() {}
F.y = 2;
`

	output :=
		`Function.prototype.a = 1;
// F is a subtype of Function, so we can't ambiguate "F.b" to "F.a".
/** @constructor */
function F() {}
F.b = 2;
`
	c.Test(js, output)
}

// Java: AmbiguatePropertiesTest.testPropertyOfObjectOfUnknownType()
func (c *AmbiguatePropertiesTest) TestPropertyOfObjectOfUnknownType() {
	c.TestSame("var foo = x(); foo.prop = '';")
}

// Java: AmbiguatePropertiesTest.testPropertyOnParamOfUnknownType()
func (c *AmbiguatePropertiesTest) TestPropertyOnParamOfUnknownType() {
	c.TestSame(
		`/** @constructor */ function Foo(){};
Foo.prototype.prop = 0;
function go(aFoo){
  aFoo.prop = 1;
}
`)
}

// Java: AmbiguatePropertiesTest.testSetPropertyOfGlobalThis()
func (c *AmbiguatePropertiesTest) TestSetPropertyOfGlobalThis() {
	c.Test("this.prop = 'bar'", "this.a = 'bar'")
}

// Java: AmbiguatePropertiesTest.testReadPropertyOfGlobalThis()
func (c *AmbiguatePropertiesTest) TestReadPropertyOfGlobalThis() {
	c.TestSame(c.Externs(ambiguatePropertiesTest_EXTERNS+"Object.prototype.prop;"), c.Srcs("f(this.prop);"))
}

// Java: AmbiguatePropertiesTest.testSetQuotedPropertyOfThis()
func (c *AmbiguatePropertiesTest) TestSetQuotedPropertyOfThis() {
	c.TestSame("this['prop'] = 'bar';")
}

// Java: AmbiguatePropertiesTest.testExternedPropertyName()
func (c *AmbiguatePropertiesTest) TestExternedPropertyName() {
	c.Test(
		`/** @constructor */ var Bar = function(){};
/** @override */ Bar.prototype.toString = function(){};
Bar.prototype.func = function(){};
var bar = new Bar();
bar.toString();
`,
		`/** @constructor */ var Bar = function(){};
/** @override */ Bar.prototype.toString = function(){};
Bar.prototype.a = function(){};
var bar = new Bar();
bar.toString();
`)
}

// Java: AmbiguatePropertiesTest.testExternedPropertyNameDefinedByObjectLiteral()
func (c *AmbiguatePropertiesTest) TestExternedPropertyNameDefinedByObjectLiteral() {
	c.TestSame("/**@constructor*/function Bar(){};Bar.prototype.factory")
}

// Java: AmbiguatePropertiesTest.testStaticAndInstanceMethodWithSameName()
func (c *AmbiguatePropertiesTest) TestStaticAndInstanceMethodWithSameName() {
	c.Test(
		`/** @constructor */function Bar(){}; Bar.getA = function(){};
Bar.prototype.getA = function(){}; Bar.getA();
var bar = new Bar(); bar.getA();
`,
		`/** @constructor */function Bar(){}; Bar.a = function(){};
Bar.prototype.a = function(){}; Bar.a();
var bar = new Bar(); bar.a();
`)
}

// Java: AmbiguatePropertiesTest.testStaticAndInstanceProperties()
func (c *AmbiguatePropertiesTest) TestStaticAndInstanceProperties() {
	c.Test(
		`/** @constructor */function Bar(){};
Bar.getA = function(){};
Bar.prototype.getB = function(){};
`,
		`/** @constructor */function Bar(){}; Bar.a = function(){};
Bar.prototype.a = function(){};
`)
}

// Java: AmbiguatePropertiesTest.testStaticAndSubInstanceProperties()
func (c *AmbiguatePropertiesTest) TestStaticAndSubInstanceProperties() {
	c.EnableClosurePass()

	js :=
		`/** @constructor */ var Foo = function(){};
Foo.x=0;
/** @constructor
 @extends Foo */ var Bar = function(){};
goog.inherits(Bar, Foo);
Bar.y=0;
Bar.prototype.z=0;
`
	output :=
		`/** @constructor */ var Foo = function(){};
Foo.a=0;
/** @constructor
 @extends Foo */ var Bar = function(){};
goog.inherits(Bar, Foo);
Bar.a=0;
Bar.prototype.a=0;
`
	c.Test(js, output)
}

// Java: AmbiguatePropertiesTest.testStatic()
func (c *AmbiguatePropertiesTest) TestStatic() {
	js :=
		`/** @constructor */ var Foo = function() {};
Foo.x = 0;
/** @constructor */ var Bar = function() {};
Bar.y = 0;
`
	output :=
		`/** @constructor */ var Foo = function() {};
Foo.a = 0;
/** @constructor */ var Bar = function() {};
Bar.a = 0;
`
	c.Test(js, output)
}

// Java: AmbiguatePropertiesTest.testClassWithStaticsPassedToUnrelatedFunction()
func (c *AmbiguatePropertiesTest) TestClassWithStaticsPassedToUnrelatedFunction() {
	js :=
		`/** @constructor */ var Foo = function() {};
Foo.x = 0;
/** @param {!Function} x */ function f(x) { x.y = 1; x.z = 2; }
f(Foo)
`
	output :=
		`/** @constructor */ var Foo = function() {};
Foo.a = 0;
/** @param {!Function} x */ function f(x) { x.y = 1; x.z = 2; }
f(Foo)
`
	c.Test(js, output)
}

// Java: AmbiguatePropertiesTest.testClassWithStaticsPassedToRelatedFunction()
func (c *AmbiguatePropertiesTest) TestClassWithStaticsPassedToRelatedFunction() {
	js :=
		`/** @constructor */ var Foo = function() {};
Foo.x = 0;
/** @param {!Function} x */ function f(x) { x.y = 1; x.x = 2;}
f(Foo)
`
	c.TestSame(js)
}

// Java: AmbiguatePropertiesTest.testFunctionInstanceProtoype_isInvalidating()
func (c *AmbiguatePropertiesTest) TestFunctionInstanceProtoype_isInvalidating() {
	js :=
		`class Foo {
  bar() { }
}

function addPrototypeProp(/** !Function */ ctor) {
  ctor.prototype.kip = 0;
}
`
	output :=
		`class Foo {
  a() { }
}

function addPrototypeProp(/** !Function */ ctor) {
  ctor.prototype.kip = 0;
}
`
	c.Test(js, output)
}

// Java: AmbiguatePropertiesTest.testTypeMismatch()
func (c *AmbiguatePropertiesTest) TestTypeMismatch() {
	c.IgnoreWarnings(DiagnosticGroups_CHECK_TYPES)
	c.TestSame(
		`/** @constructor */var Foo = function(){};
/** @constructor */var Bar = function(){};
Foo.prototype.b = 0;
/** @type {Foo} */
var F = new Bar();
`)
}

// Java: AmbiguatePropertiesTest.testRenamingMap()
func (c *AmbiguatePropertiesTest) TestRenamingMap() {
	js :=
		`/** @constructor */ var Foo = function(){};
Foo.prototype.z=0;
Foo.prototype.z=0;
Foo.prototype.x=0;
Foo.prototype.y=0;
`
	output :=
		`/** @constructor */ var Foo = function(){};
Foo.prototype.a=0;
Foo.prototype.a=0;
Foo.prototype.b=0;
Foo.prototype.c=0;
`
	c.Test(js, output)

	answerMap := jutil.NewLinkedHashMap[string, string]()
	answerMap.Put("x", "b")
	answerMap.Put("y", "c")
	answerMap.Put("z", "a")
	AssertThat(c.lastPass.GetRenamingMap()).IsEqualTo(answerMap)
}

// Java: AmbiguatePropertiesTest.testInline()
func (c *AmbiguatePropertiesTest) TestInline() {
	js :=
		`/** @interface */ function Foo(){}
Foo.prototype.x = function(){};
/**
 * @constructor
 * @implements {Foo}
 */
function Bar(){}
Bar.prototype.y;
/** @override */
Bar.prototype.x = function() { return this.y; };
Bar.prototype.z = function() {};
/** @type {Foo} */ (new Bar).y;
`
	output :=
		`/** @interface */ function Foo(){}
Foo.prototype.b = function(){};
/**
 * @constructor
 * @implements {Foo}
 */
function Bar(){}
Bar.prototype.a;
/** @override */
Bar.prototype.b = function() { return this.a; };
Bar.prototype.c = function() {};
(new Bar).a;
`
	c.Test(js, output)
}

// Java: AmbiguatePropertiesTest.testImplementsAndExtends()
func (c *AmbiguatePropertiesTest) TestImplementsAndExtends() {
	js :=
		`/** @interface */ function Foo() {}
/** @constructor */ function Bar(){}
Bar.prototype.y = function() { return 3; };
/**
 * @constructor
 * @extends {Bar}
 * @implements {Foo}
 */
function SubBar(){ }
/** @param {Foo} x1 */ function f(x1) { x1.z = 3; }
/** @param {SubBar} x2 */ function g(x2) { x2.z = 3; }
`
	output :=
		`/** @interface */ function Foo(){}
/** @constructor */ function Bar(){}
Bar.prototype.b = function() { return 3; };
/**
 * @constructor
 * @extends {Bar}
 * @implements {Foo}
 */
function SubBar(){}
/** @param {Foo} x1 */ function f(x1) { x1.a = 3; }
/** @param {SubBar} x2 */ function g(x2) { x2.a = 3; }
`
	c.Test(js, output)
}

// Java: AmbiguatePropertiesTest.testImplementsAndExtends_respectsUndeclaredProperties()
func (c *AmbiguatePropertiesTest) TestImplementsAndExtends_respectsUndeclaredProperties() {
	js :=
		`/** @interface */ function A() {}
/**
 * @constructor
 */
function C1(){}
/**
 * @constructor
 * @extends {C1}
 * @implements {A}
 */
function C2(){}
/** @param {C1} x1 */ function f(x1) { x1.y = 3; }
/** @param {A} x2 */ function g(x2) { x2.z = 3; }
`
	output :=
		`/** @interface */ function A(){}
/**
 * @constructor
 */
function C1(){}
/**
 * @constructor
 * @extends {C1}
 * @implements {A}
 */
function C2(){}
/** @param {C1} x1 */ function f(x1) { x1.a = 3; }
/** @param {A} x2 */ function g(x2) { x2.b = 3; }
`
	c.Test(js, output)
}

// Java: AmbiguatePropertiesTest.testImplementsAndExtendsEs6Class_respectsUndeclaredProperties()
func (c *AmbiguatePropertiesTest) TestImplementsAndExtendsEs6Class_respectsUndeclaredProperties() {
	js :=
		`/** @interface */
class A {
  constructor() {
// Optional property; C2 does not need to declare it implemented.
    /** @type {number|undefined} */
    this.y;
  }
}
class C1 {
  constructor() {
    /** @type {number} */
    this.z;
  }
}
/** @implements {A} */
class C2 extends C1 {}
/** @param {A} x1 */ function f(x1) { x1.y = 3; }
/** @param {C1} x2 */ function g(x2) { x2.z = 3; }
`
	output :=
		`/** @interface */
class A {
  constructor() {
    /** @type {number|undefined} */
    this.a;
  }
}
class C1 {
  constructor() {
    /** @type {number} */
    this.b;
  }
}
/** @implements {A} */
class C2 extends C1 {}
/** @param {A} x1 */ function f(x1) { x1.a = 3; }
/** @param {C1} x2 */ function g(x2) { x2.b = 3; }
`
	c.Test(js, output)
}

// Java: AmbiguatePropertiesTest.testExtendsInterface()
func (c *AmbiguatePropertiesTest) TestExtendsInterface() {
	js :=
		`/** @interface */ function A() {}
/** @interface
 @extends {A} */ function B() {}
/** @param {A} x1 */ function f(x1) { x1.y = 3; }
/** @param {B} x2 */ function g(x2) { x2.z = 3; }
`
	output :=
		`/** @interface */ function A(){}
/** @interface
 @extends {A} */ function B(){}
/** @param {A} x1 */ function f(x1) { x1.a = 3; }
/** @param {B} x2 */ function g(x2) { x2.b = 3; }
`
	c.Test(js, output)
}

// Java: AmbiguatePropertiesTest.testInterfaceWithSubInterfaceAndDirectImplementors()
func (c *AmbiguatePropertiesTest) TestInterfaceWithSubInterfaceAndDirectImplementors() {
	c.IgnoreWarnings(DiagnosticGroups_CHECK_TYPES)
	c.Test(
		`/** @interface */ function Foo(){};
Foo.prototype.foo = function(){};
/** @interface @extends {Foo} */ function Bar(){};
/** @constructor @implements {Foo} */ function Baz(){};
Baz.prototype.baz = function(){};
`,
		`/** @interface */ function Foo(){};
Foo.prototype.b = function(){};
/** @interface @extends {Foo} */ function Bar(){};
/** @constructor @implements {Foo} */ function Baz(){};
Baz.prototype.a = function(){};
`)
}

// Java: AmbiguatePropertiesTest.testPredeclaredType()
func (c *AmbiguatePropertiesTest) TestPredeclaredType() {
	c.EnableClosurePass()
	c.EnableRewriteClosureProvides()

	js :=
		`goog.forwardDeclare('goog.Foo');
/** @constructor */
function A() {
  this.x = 3;
}
/** @param {goog.Foo} x */
function f(x) { x.y = 4; }
`
	result :=
		`/** @constructor */
function A() {
  this.a = 3;
}
/** @param {goog.Foo} x */
function f(x) { x.y = 4; }
`
	c.Test(
		c.Externs(ambiguatePropertiesTest_EXTERNS+NewTestExternsBuilder().AddClosureExterns().Build()),
		c.Srcs(js),
		c.Expected(result))
}

// Java: AmbiguatePropertiesTest.testBug14291280()
func (c *AmbiguatePropertiesTest) TestBug14291280() {
	js :=
		`/** @constructor
 @template T */

function C() {

  this.aa = 1;

}

/** @return {C.<T>} */

C.prototype.method = function() {

  return this;

}

/** @type {C.<string>} */

var x = new C().method();

x.bb = 2;

x.cc = 3;
`
	result :=
		`/** @constructor
 @template T */

function C() {
  this.b = 1;
}
/** @return {C.<T>} */

C.prototype.a = function() {
  return this;
};
/** @type {C.<string>} */

var x = (new C).a();
x.c = 2;
x.d = 3;
`
	c.Test(js, result)
}

// Java: AmbiguatePropertiesTest.testAmbiguateWithAnAlias()
func (c *AmbiguatePropertiesTest) TestAmbiguateWithAnAlias() {
	js :=
		`/** @constructor */ function Foo() { this.abc = 5; }

/** @const */ var alias = Foo;

/** @constructor @extends alias */

function Bar() {

  this.xyz = 7;

}
`
	result :=
		`/** @constructor */ function Foo() { this.a = 5; }

/** @const */ var alias = Foo;

/** @constructor @extends alias */

function Bar() {

  this.b = 7;

}
`
	c.Test(js, result)
}

// Java: AmbiguatePropertiesTest.testAmbiguateWithAliases()
func (c *AmbiguatePropertiesTest) TestAmbiguateWithAliases() {
	js :=
		`/** @constructor */ function Foo() { this.abc = 5; }

/** @const */ var alias = Foo;

/** @constructor @extends alias */

function Bar() {

  this.def = 7;

}

/** @constructor @extends alias */

function Baz() {

  this.xyz = 8;

}
`
	result :=
		`/** @constructor */ function Foo() { this.a = 5; }

/** @const */ var alias = Foo;

/** @constructor @extends alias */

function Bar() {

  this.b = 7;

}

/** @constructor @extends alias */

function Baz() {

  this.b = 8;

}
`
	c.Test(js, result)
}

// Java: AmbiguatePropertiesTest.testAmbiguateWithStructuralInterfaces()
func (c *AmbiguatePropertiesTest) TestAmbiguateWithStructuralInterfaces() {
	js :=
		`/** @record */
function Record() {}
/** @type {number|undefined} */
Record.prototype.recordProp;

function f(/** !Record */ a) { use(a.recordProp); }

/** @constructor */
function Type() {
  /** @const */
  this.classProp = 'a';
}
f(new Type)
`

	expected :=
		`/** @record */
function Record() {}
/** @type {number|undefined} */
Record.prototype.recordProp;

function f(/** !Record */ a) { use(a.recordProp); }

/** @constructor */
function Type() {
  /** @const */
  this.a = 'a';
}
f(new Type)
`

	c.Test(js, expected)
}

// Java: AmbiguatePropertiesTest.structuralTypesNotAmbiguated_forwardRefsNotInUnion()
func (c *AmbiguatePropertiesTest) TestStructuralTypesNotAmbiguated_forwardRefsNotInUnion() {
	// edge case where the early referenced to "Params" caused data.name => data.a
	c.Test(
		`class Foo {
  method() {}
}
/** @param {!Params=} data */
function f(data) {
  return data.name;
}
/** @typedef {{name: string}} */
let Params;
`,
		`class Foo {
  a() {}
}
/** @param {!Params=} data */
function f(data) {
  return data.name;
}
/** @typedef {{name: string}} */
let Params;
`)
}

// Java: AmbiguatePropertiesTest.testObjectSpreadDoesNotCrash()
func (c *AmbiguatePropertiesTest) TestObjectSpreadDoesNotCrash() {
	c.TestSame("var lit1 = {...a};")
}

// Java: AmbiguatePropertiesTest.testObjectRestDoesNotCrash()
func (c *AmbiguatePropertiesTest) TestObjectRestDoesNotCrash() {
	// only normalization changes, nothing else
	c.Test("var {...a} = b;", "var a; ({...a} = b);")
}

// Java: AmbiguatePropertiesTest.testUnrelatedObjectLiterals()
func (c *AmbiguatePropertiesTest) TestUnrelatedObjectLiterals() {
	c.TestSame(
		`/** @constructor */ function Foo() {}
/** @constructor */ function Bar() {}
var lit1 = {a: new Foo()};
var lit2 = {b: new Bar()};
`)
}

// Java: AmbiguatePropertiesTest.testObjectLitTwoVar()
func (c *AmbiguatePropertiesTest) TestObjectLitTwoVar() {
	js :=
		`/** @constructor */ var Foo = function(){};
Foo.prototype = {z:0, z:1, x:0};
`
	output :=
		`/** @constructor */ var Foo = function(){};
Foo.prototype = {a:0, a:1, b:0};
`
	c.Test(js, output)
}

// Java: AmbiguatePropertiesTest.testObjectLitTwoIndependentVar()
func (c *AmbiguatePropertiesTest) TestObjectLitTwoIndependentVar() {
	js :=
		`/** @constructor */ var Foo = function(){};
Foo.prototype = {b: 0};
/** @constructor */ var Bar = function(){};
Bar.prototype = {c: 0};
`
	output :=
		`/** @constructor */ var Foo = function(){};
Foo.prototype = {a: 0};
/** @constructor */ var Bar = function(){};
Bar.prototype = {a: 0};
`
	c.Test(js, output)
}

// Java: AmbiguatePropertiesTest.testObjectLitTwoTypesTwoVar()
func (c *AmbiguatePropertiesTest) TestObjectLitTwoTypesTwoVar() {
	js :=
		`/** @constructor */ var Foo = function(){};
Foo.prototype = {r: 0, g: 0};
/** @constructor */ var Bar = function(){};
Bar.prototype = {c: 0, r: 0};
`
	output :=
		`/** @constructor */ var Foo = function(){};
Foo.prototype = {a: 0, b: 0};
/** @constructor */ var Bar = function(){};
Bar.prototype = {b: 0, a: 0};
`
	c.Test(js, output)
}

// Java: AmbiguatePropertiesTest.testObjectLitUnion()
func (c *AmbiguatePropertiesTest) TestObjectLitUnion() {
	js :=
		`/** @constructor */ var Foo = function(){};
/** @constructor */ var Bar = function(){};
Foo.prototype = {foodoo: 0};
Bar.prototype = {bardoo: 0};
/** @type {Foo|Bar} */
var U = any();
U.joint;
`
	output :=
		`/** @constructor */ var Foo = function(){};
/** @constructor */ var Bar = function(){};
Foo.prototype = {a: 0};
Bar.prototype = {a: 0};
/** @type {Foo|Bar} */
var U = any();
U.b;
`
	c.Test(js, output)
}

// Java: AmbiguatePropertiesTest.testObjectLitUnions()
func (c *AmbiguatePropertiesTest) TestObjectLitUnions() {
	js :=
		`/** @constructor */ var Foo = function(){};
/** @constructor */ var Bar = function(){};
/** @constructor */ var Baz = function(){};
/** @constructor */ var Bat = function(){};
Foo.prototype = {lone1: 0};
Bar.prototype = {lone2: 0};
Baz.prototype = {lone3: 0};
Bat.prototype = {lone4: 0};
/** @type {Foo|Bar} */
var U1 = any();
U1.j1;
U1.j2;
/** @type {Baz|Bar} */
var U2 = any();
U2.j3;
U2.j4;
/** @type {Baz|Bat} */
var U3 = any();
U3.j5;
U3.j6
`
	output :=
		`/** @constructor */ var Foo = function(){};
/** @constructor */ var Bar = function(){};
/** @constructor */ var Baz = function(){};
/** @constructor */ var Bat = function(){};
Foo.prototype = {c: 0};
Bar.prototype = {e: 0};
Baz.prototype = {e: 0};
Bat.prototype = {c: 0};
/** @type {Foo|Bar} */
var U1 = any();
U1.a;
U1.b;
/** @type {Baz|Bar} */
var U2 = any();
U2.c;
U2.d;
/** @type {Baz|Bat} */
var U3 = any();
U3.a;
U3.b
`
	c.Test(js, output)
}

// Java: AmbiguatePropertiesTest.testObjectLitExtends()
func (c *AmbiguatePropertiesTest) TestObjectLitExtends() {
	c.EnableClosurePass()

	js :=
		`/** @constructor */ var Foo = function(){};
Foo.prototype = {x: 0};
/** @constructor
 @extends Foo */ var Bar = function(){};
goog.inherits(Bar, Foo);
Bar.prototype = {y: 0, z: 0};
/** @constructor */ var Baz = function(){};
Baz.prototype = {l: 0, m: 0, n: 0};
(new Baz).m
`
	output :=
		`/** @constructor */ var Foo = function(){};
Foo.prototype = {a: 0};
/** @constructor
 @extends Foo */ var Bar = function(){};
goog.inherits(Bar, Foo);
Bar.prototype = {b: 0, c: 0};
/** @constructor */ var Baz = function(){};
Baz.prototype = {b: 0, a: 0, c: 0};
(new Baz).a
`
	c.Test(js, output)
}

// Java: AmbiguatePropertiesTest.testObjectLitLotsOfClasses()
func (c *AmbiguatePropertiesTest) TestObjectLitLotsOfClasses() {
	var b strings.Builder
	classes := 10
	for i := 0; i < classes; i++ {
		c_ := "Foo" + strconv.Itoa(i)
		b.WriteString("/** @constructor */ var " + c_ + " = function(){};\n")
		b.WriteString(c_ + ".prototype = {varness" + strconv.Itoa(i) + ": 0};")
	}
	js := b.String()
	c.Test(js, jutil.StringReplaceAll(js, "varness\\d+", "a"))
}

// Java: AmbiguatePropertiesTest.testObjectLitFunctionType()
func (c *AmbiguatePropertiesTest) TestObjectLitFunctionType() {
	js :=
		`/** @constructor */ function Foo(){};
Foo.prototype = { /** @return {Bar} */ fun: function() { return new Bar(); } };
/** @constructor */ function Bar(){};
Bar.prototype.bazz;
(new Foo).fun().bazz();
`
	output :=
		`/** @constructor */ function Foo(){};
Foo.prototype = { /** @return {Bar} */ a: function() { return new Bar(); } };
/** @constructor */ function Bar(){};
Bar.prototype.a;
(new Foo).a().a();
`
	c.Test(js, output)
}

// Java: AmbiguatePropertiesTest.testObjectLitOverlappingOriginalAndGeneratedNames()
func (c *AmbiguatePropertiesTest) TestObjectLitOverlappingOriginalAndGeneratedNames() {
	c.Test(
		`/** @constructor */ function Bar(){};
Bar.prototype = {b: function(){}, a: function(){}};
var bar = new Bar();
bar.b();
`,
		`/** @constructor */ function Bar(){};
Bar.prototype = {a: function(){}, b: function(){}};
var bar = new Bar();
bar.a();
`)
}

// Java: AmbiguatePropertiesTest.testEnum()
func (c *AmbiguatePropertiesTest) TestEnum() {
	c.TestSame(
		`/** @enum {string} */ var Foo = {X: 'y'};
var x = Foo.X
`)
}

// Java: AmbiguatePropertiesTest.testUnannotatedConstructorsDontCrash()
func (c *AmbiguatePropertiesTest) TestUnannotatedConstructorsDontCrash() {
	c.TestSame(
		`function Foo() {}
Foo.prototype.a;
function Bar() {}
Bar.prototype.a;
`)
}

// Java: AmbiguatePropertiesTest.testGenericPrototypeObject()
func (c *AmbiguatePropertiesTest) TestGenericPrototypeObject() {
	js :=
		`/**
 * @constructor
 * @template T
 */
function Foo() {
  this.a = 1;
}
/** @constructor @extends {Foo<number>} */
function Bar() {}
/** @constructor */
function Baz() {
  this.b = 2;
}
`

	output :=
		`/**
 * @constructor
 * @template T
 */
function Foo() {
  this.a = 1;
}
/** @constructor @extends {Foo<number>} */
function Bar() {}
/** @constructor */
function Baz() {
  this.a = 2;
}
`

	c.Test(js, output)
}

// Java: AmbiguatePropertiesTest.testRelatedInstanceAndPrototypePropNotAmbiguated()
func (c *AmbiguatePropertiesTest) TestRelatedInstanceAndPrototypePropNotAmbiguated() {
	c.Test(
		`class Foo {
  constructor() {
    this.y = 0;
  }
  x() {}
}
`,
		`class Foo {
  constructor() {
    this.b = 0;
  }
  a() {}
}
`)
}

// Java: AmbiguatePropertiesTest.testRelatedInstanceAndSuperclassPrototypePropNotAmbiguated()
func (c *AmbiguatePropertiesTest) TestRelatedInstanceAndSuperclassPrototypePropNotAmbiguated() {
	c.Test(
		`class Foo {
  x() {}
}
class Bar extends Foo {
  constructor() {
    this.y = 0;
  }
}
`,
		`class Foo {
  a() {}
}
class Bar extends Foo {
  constructor() {
    this.b = 0;
  }
}
`)
}

// Java: AmbiguatePropertiesTest.testPropertiesWithTypesThatHaveBeenNarrowed()
func (c *AmbiguatePropertiesTest) TestPropertiesWithTypesThatHaveBeenNarrowed() {
	js :=
		`/** @constructor */
function Foo() {
  /** @type {?Object} */
  this.prop1 = {};
  /** @type {?Object} */
  this.prop2;
  this.headerObserver_;
}
Foo.prototype.m = function() {
  this.prop2 = {};
  return this.headerObserver_;
};
`

	output :=
		`/** @constructor */
function Foo() {
  /** @type {?Object} */
  this.d = {};
  /** @type {?Object} */
  this.b;
  this.a;
}
Foo.prototype.c = function() {
  this.b = {};
  return this.a;
};
`

	c.Test(js, output)
}

// Java: AmbiguatePropertiesTest.testDontRenamePrototypeWithoutExterns()
func (c *AmbiguatePropertiesTest) TestDontRenamePrototypeWithoutExterns() {
	js :=
		`/** @interface */
function Foo() {}
/** @return {!Foo} */
Foo.prototype.foo = function() {};
`

	output :=
		`/** @interface */
function Foo() {}
/** @return {!Foo} */
Foo.prototype.a = function() {};
`

	c.Test(c.Externs(""), c.Srcs(js), c.Expected(output))
}

// Java: AmbiguatePropertiesTest.testInvalidRenameFunction_doesNotCrash()
func (c *AmbiguatePropertiesTest) TestInvalidRenameFunction_doesNotCrash() {
	c.TestSame("const p = JSCompiler_renameProperty('foo', 0, 1)")
	c.TestSame("const p = JSCompiler_renameProperty(0)")
	c.TestSame("const p = JSCompiler_renameProperty('a.b')")
}

// Java: AmbiguatePropertiesTest.testJSCompiler_renameProperty_twoArgs_blocksAmbiguation()
func (c *AmbiguatePropertiesTest) TestJSCompiler_renameProperty_twoArgs_blocksAmbiguation() {
	c.TestSame(
		`/** @constructor */
function Foo() {}
Foo.prototype.bar = 3;
const barName = JSCompiler_renameProperty('bar', Foo);
`)
}

// Java: AmbiguatePropertiesTest.testJSCompiler_renameProperty_oneArg_blocksAmbiguation()
func (c *AmbiguatePropertiesTest) TestJSCompiler_renameProperty_oneArg_blocksAmbiguation() {
	c.TestSame(
		`/** @constructor */
function Foo() {}
Foo.prototype.bar = 3;
const barName = JSCompiler_renameProperty('bar');
`)
}

// Java: AmbiguatePropertiesTest.testSingleClass_withSingleMemberFn_ambiguated()
func (c *AmbiguatePropertiesTest) TestSingleClass_withSingleMemberFn_ambiguated() {
	c.Test("class Foo { methodFoo() {} }", "class Foo { a() {} }")
}

// Java: AmbiguatePropertiesTest.testSingleClass_withSingleMemberField_ambiguated()
func (c *AmbiguatePropertiesTest) TestSingleClass_withSingleMemberField_ambiguated() {
	c.Test("class Foo { field = 2; } ", "class Foo { a = 2; }")
}

// Java: AmbiguatePropertiesTest.testSingleClass_withSingleQuotedMemberField_notAmbiguated()
func (c *AmbiguatePropertiesTest) TestSingleClass_withSingleQuotedMemberField_notAmbiguated() {
	c.TestSame("/** @dict */ class Foo { 'field' = 2; }")
}

// Java: AmbiguatePropertiesTest.testSingleClass_withSingleStaticMemberField_ambiguated()
func (c *AmbiguatePropertiesTest) TestSingleClass_withSingleStaticMemberField_ambiguated() {
	c.Test(
		`class Foo {
  static field = 1;
  static {
    Foo.field = 2;
    this.field = 3;
  }
}
`,
		`class Foo {
  static a = 1;
  static {
    Foo.a = 2;
    this.a = 3;
  }
}
`)
}

// Java: AmbiguatePropertiesTest.testSingleClass_withSingleComputedMemberField_notAmbiguated()
func (c *AmbiguatePropertiesTest) TestSingleClass_withSingleComputedMemberField_notAmbiguated() {
	c.TestSame("/** @dict */ class Foo { ['field'] = 2; }")
}

// Java: AmbiguatePropertiesTest.testSingleClass_withSingleStaticComputedMemberField_notAmbiguated()
func (c *AmbiguatePropertiesTest) TestSingleClass_withSingleStaticComputedMemberField_notAmbiguated() {
	c.TestSame("/** @dict */ class Foo { static ['field'] = 2; }")
}

// Java: AmbiguatePropertiesTest.testSingleClass_withSingleNumberMemberField_notAmbiguated()
func (c *AmbiguatePropertiesTest) TestSingleClass_withSingleNumberMemberField_notAmbiguated() {
	c.TestSame("/** @dict */ class Foo { 1 = 2; }")
}

// Java: AmbiguatePropertiesTest.testSingleClass_withSingleStaticNumberMemberField_notAmbiguated()
func (c *AmbiguatePropertiesTest) TestSingleClass_withSingleStaticNumberMemberField_notAmbiguated() {
	c.TestSame("/** @dict */ class Foo { static 1 = 2; }")
}

// Java: AmbiguatePropertiesTest.testQuotedMemberFieldInClass_notAmbiguated()
func (c *AmbiguatePropertiesTest) TestQuotedMemberFieldInClass_notAmbiguated() {
	c.TestSame("/** @dict */ class Foo { 'fieldFoo' }")
}

// Java: AmbiguatePropertiesTest.testQuotedMemberFieldInClassReservesPropertyName()
func (c *AmbiguatePropertiesTest) TestQuotedMemberFieldInClassReservesPropertyName() {
	c.Test(
		"/** @unrestricted */ class Foo { 'a'; foo; }",
		"/** @unrestricted */ class Foo { 'a'; b; }")
}

// Java: AmbiguatePropertiesTest.testSingleClass_withTwoMemberFields_notAmbiguated()
func (c *AmbiguatePropertiesTest) TestSingleClass_withTwoMemberFields_notAmbiguated() {
	c.Test("class Foo { field1; field2 = 2;}", "class Foo { a; b = 2;}")
}

// Java: AmbiguatePropertiesTest.testSingleClass_withStaticAndPrototypeMemberFields_ambiguated()
func (c *AmbiguatePropertiesTest) TestSingleClass_withStaticAndPrototypeMemberFields_ambiguated() {
	c.Test("class Foo { static staticField; memberField; }", "class Foo { static a; a; }")
}

// Java: AmbiguatePropertiesTest.testTwoUnrelatedClasses_withMemberFields_ambiguated()
func (c *AmbiguatePropertiesTest) TestTwoUnrelatedClasses_withMemberFields_ambiguated() {
	c.Test("class Foo { fieldFoo; } class Bar { fieldBar; }", "class Foo { a; } class Bar { a; }")
}

// Java: AmbiguatePropertiesTest.testTwoUnrelatedClasses_withStaticMemberFields_ambiguated()
func (c *AmbiguatePropertiesTest) TestTwoUnrelatedClasses_withStaticMemberFields_ambiguated() {
	c.Test(
		`class Foo { static fieldFoo = 2; }
class Bar { static fieldBar = 2; }
`,
		`class Foo { static a = 2; }
class Bar { static a = 2; }
`)
}

// Java: AmbiguatePropertiesTest.testEs6SuperclassStaticField_notAmbiguated()
func (c *AmbiguatePropertiesTest) TestEs6SuperclassStaticField_notAmbiguated() {
	c.Test(
		`class Foo { static fieldFoo; }
class Bar extends Foo { static fieldBar; }
`,
		// Since someone could access Foo.methodFoo() through Bar, make sure the fields get
		// distinct names.
		`class Foo { static b; }
class Bar extends Foo { static a; }
`)
}

// Java: AmbiguatePropertiesTest.testSingleClass_withMultipleMemberFields_notAmbiguated()
func (c *AmbiguatePropertiesTest) TestSingleClass_withMultipleMemberFields_notAmbiguated() {
	c.Test(
		`/** @unrestricted */
class Foo {
  field = 1;
  ['hi'] = 4;
  static 1 = 2;
  hello = 5;
}
`,
		`/** @unrestricted */
class Foo {
  a = 1;
  ['hi'] = 4;
  static 1 = 2;
  b = 5;
}
`)
}

// Java: AmbiguatePropertiesTest.testQuotedMemberFnInClass_notAmbiguated()
func (c *AmbiguatePropertiesTest) TestQuotedMemberFnInClass_notAmbiguated() {
	c.TestSame("/** @dict */ class Foo { 'methodFoo'() {} }")
}

// Java: AmbiguatePropertiesTest.testQuotedMemberFnInClassReservesPropertyName()
func (c *AmbiguatePropertiesTest) TestQuotedMemberFnInClassReservesPropertyName() {
	c.Test(
		"/** @unrestricted */ class Foo { 'a'() {} foo() {} }",
		"/** @unrestricted */ class Foo { 'a'() {} b() {} }")
}

// Java: AmbiguatePropertiesTest.testSingleClass_withTwoMemberFns_notAmbiguated()
func (c *AmbiguatePropertiesTest) TestSingleClass_withTwoMemberFns_notAmbiguated() {
	c.Test("class Foo { method1() {} method2() {} }", "class Foo { a() {} b() {} }")
}

// Java: AmbiguatePropertiesTest.testSingleClass_withStaticAndPrototypeMemberFns_ambiguated()
func (c *AmbiguatePropertiesTest) TestSingleClass_withStaticAndPrototypeMemberFns_ambiguated() {
	c.Test("class Foo { static method1() {} method2() {} }", "class Foo { static a() {} a() {} }")
}

// Java: AmbiguatePropertiesTest.testTwoUnrelatedClasses_withMemberFns_ambiguated()
func (c *AmbiguatePropertiesTest) TestTwoUnrelatedClasses_withMemberFns_ambiguated() {
	c.Test(
		"class Foo { methodFoo() {} } class Bar { methodBar() {} }",
		"class Foo { a() {} } class Bar { a() {} }")
}

// Java: AmbiguatePropertiesTest.testTwoUnrelatedClasses_withStaticMemberFns_ambiguated()
func (c *AmbiguatePropertiesTest) TestTwoUnrelatedClasses_withStaticMemberFns_ambiguated() {
	c.Test(
		`class Foo { static methodFoo() {} }
class Bar { static methodBar() {} }
`,
		`class Foo { static a() {} }
class Bar { static a() {} }
`)
}

// Java: AmbiguatePropertiesTest.testEs6SuperclassStaticMethod_notAmbiguated()
func (c *AmbiguatePropertiesTest) TestEs6SuperclassStaticMethod_notAmbiguated() {
	c.Test(
		`class Foo { static methodFoo() {} }
class Bar extends Foo { static methodBar() {} }
`,
		// Since someone could access Foo.methodFoo() through Bar, make sure the methods get
		// distinct names.
		`class Foo { static b() {} }
class Bar extends Foo { static a() {} }
`)
}

// Java: AmbiguatePropertiesTest.testEs6SubclassChain_withStaticMethods_notAmbiguated()
func (c *AmbiguatePropertiesTest) TestEs6SubclassChain_withStaticMethods_notAmbiguated() {
	c.Test(
		`class Foo { static methodFoo() { alert('foo'); } }
class Bar extends Foo { static methodBar() {alert('bar'); } }
class Baz extends Bar { static methodBaz() {alert('baz'); } }
class Moo extends Baz { static methodMoo() { alert('moo'); } }
Moo.methodFoo();
`,
		// All four static methods must get distinct names, so that Moo.a resolves correctly to
		// Foo.a
		`class Foo { static a() {alert('foo'); } }
class Bar extends Foo { static b() {alert('bar'); } }
class Baz extends Bar { static c() {alert('baz'); } }
class Moo extends Baz { static d() { alert('moo'); } }
Moo.a();
`)
}

// Java: AmbiguatePropertiesTest.testEs5ClassWithExtendsChainStaticMethods_notAmbiguated()
func (c *AmbiguatePropertiesTest) TestEs5ClassWithExtendsChainStaticMethods_notAmbiguated() {
	c.Test(
		`/** @constructor */ function Foo () {}
Foo.methodFoo = function() { alert('foo'); };
class Bar extends Foo { static methodBar() {alert('bar'); } }
class Baz extends Bar { static methodBaz() {alert('baz'); } }
class Moo extends Baz { static methodMoo() { alert('moo'); } }
Moo.methodFoo();
`,
		`/** @constructor */ function Foo () {}
Foo.a = function() { alert('foo'); };
class Bar extends Foo { static b() {alert('bar'); } }
class Baz extends Bar { static c() {alert('baz'); } }
class Moo extends Baz { static d() { alert('moo'); } }
Moo.a();
`)
}

// Java: AmbiguatePropertiesTest.testEs5ClassWithEs5SubclassWtaticMethods_ambiguated()
func (c *AmbiguatePropertiesTest) TestEs5ClassWithEs5SubclassWtaticMethods_ambiguated() {
	c.Test(
		`/** @constructor */ function Foo () {}
Foo.methodFoo = function() { alert('foo'); };
/** @constructor @extends {Foo} */ function Bar() {}
Bar.methodBar = function() { alert('bar'); };
`,
		`/** @constructor */ function Foo () {}
Foo.a = function() { alert('foo'); };
/** @constructor @extends {Foo} */ function Bar() {}
Bar.a = function() { alert('bar'); };
`)
}

// Java: AmbiguatePropertiesTest.testClassWithSuperclassStaticMethodsCalledWithSuper_ambiguated()
func (c *AmbiguatePropertiesTest) TestClassWithSuperclassStaticMethodsCalledWithSuper_ambiguated() {
	c.Test(
		`class Foo { static methodFoo() {} }
class Bar extends Foo { static methodBar() { super.methodFoo(); } }
`,
		// Since someone could access Foo.methodFoo() through Bar, make sure the methods get
		// distinct names.
		`class Foo { static a() {} }
class Bar extends Foo { static b() { super.a(); } }
`)
}

// Java: AmbiguatePropertiesTest.testGetterInClass_ambiguated()
func (c *AmbiguatePropertiesTest) TestGetterInClass_ambiguated() {
	c.Test("class Foo { get prop() {} }", "class Foo { get a() {} }")
}

// Java: AmbiguatePropertiesTest.testQuotedGetterInClass_isNotAmbiguated()
func (c *AmbiguatePropertiesTest) TestQuotedGetterInClass_isNotAmbiguated() {
	c.TestSame("/** @dict */ class Foo { get 'prop'() {} }")
}

// Java: AmbiguatePropertiesTest.testQuotedGetterInClassReservesPropertyName()
func (c *AmbiguatePropertiesTest) TestQuotedGetterInClassReservesPropertyName() {
	c.Test(
		"/** @unrestricted */ class Foo { get 'a'() {} foo() {} }",
		"/** @unrestricted */ class Foo { get 'a'() {} b() {} }")
}

// Java: AmbiguatePropertiesTest.testSetterInClass_isAmbiguated()
func (c *AmbiguatePropertiesTest) TestSetterInClass_isAmbiguated() {
	c.Test("class Foo { set prop(x) {} }", "class Foo { set a(x) {} }")
}

// Java: AmbiguatePropertiesTest.testQuotedSetterInClass_notAmbiguated()
func (c *AmbiguatePropertiesTest) TestQuotedSetterInClass_notAmbiguated() {
	c.TestSame("/** @dict */ class Foo { set 'prop'(x) {} }")
}

// Java: AmbiguatePropertiesTest.testSameGetterAndSetterInClass_ambiguated()
func (c *AmbiguatePropertiesTest) TestSameGetterAndSetterInClass_ambiguated() {
	c.Test("class Foo { get prop() {} set prop(x) {} }", "class Foo { get a() {} set a(x) {} }")
}

// Java: AmbiguatePropertiesTest.testDistinctGetterAndSetterInClass_notAmbiguated()
func (c *AmbiguatePropertiesTest) TestDistinctGetterAndSetterInClass_notAmbiguated() {
	c.Test("class Foo { set propA(x) {} get propB() {} }", "class Foo { set a(x) {} get b() {} }")
}

// Java: AmbiguatePropertiesTest.testComputedMemberFunctionInClass_notAmbiguated()
func (c *AmbiguatePropertiesTest) TestComputedMemberFunctionInClass_notAmbiguated() {
	c.TestSame("/** @dict */ class Foo { ['method']() {}}")
}

// Java: AmbiguatePropertiesTest.testSimpleComputedMemberFnInClassReservesPropertyName()
func (c *AmbiguatePropertiesTest) TestSimpleComputedMemberFnInClassReservesPropertyName() {
	c.Test(
		"/** @unrestricted */ class Foo { ['a']() {} foo() {} }",
		"/** @unrestricted */ class Foo { ['a']() {} b() {} }")
}

// Java: AmbiguatePropertiesTest.testComplexComputedMemberFnInClassDoesntReservePropertyName()
func (c *AmbiguatePropertiesTest) TestComplexComputedMemberFnInClassDoesntReservePropertyName() {
	// we don't try to evaluate 'a' + '' to see that it's identical to 'a', and so end up renaming
	// 'foo' -> 'a'
	// The property name invalidation is already just a heuristic, so just handle the very simple
	// case of ['a']() {}
	c.Test(
		"/** @unrestricted */ class Foo { ['a' + '']() {} foo() {} }",
		"/** @unrestricted */ class Foo { ['a' + '']() {} a() {} }")
}

// Java: AmbiguatePropertiesTest.testEs6ClassConstructorMethod_notAmbiguated()
func (c *AmbiguatePropertiesTest) TestEs6ClassConstructorMethod_notAmbiguated() {
	c.TestSame("class Foo { constructor() {} }")
}

// Java: AmbiguatePropertiesTest.testOptchainGet_ambiguated()
func (c *AmbiguatePropertiesTest) TestOptchainGet_ambiguated() {
	c.Test("class Foo { foo() { this?.foo(); } };", "class Foo { a() { this?.a(); } };")
}

// Java: AmbiguatePropertiesTest.testAmbiguateEs6ClassMethodsDoesntCrashOnClassInACast()
func (c *AmbiguatePropertiesTest) TestAmbiguateEs6ClassMethodsDoesntCrashOnClassInACast() {
	// the cast causes the actual CLASS node to have the unknown type, so verify that the pass
	// can handle it not being a function type.
	c.TestSame(
		`const Foo = /** @type {?} */ (class {
  method() {}
});
class Bar {
  method() {}
}
`)
}

// Java: AmbiguatePropertiesTest.testObjectSetPrototypeOfIsIgnored()
func (c *AmbiguatePropertiesTest) TestObjectSetPrototypeOfIsIgnored() {
	c.Test(
		c.Externs("Object.setPrototypeOf = function(obj, proto) {}"),
		c.Srcs(
			`/** @constructor */
function Foo() {}
Foo.fooMethod = () => 3;
/** @constructor */
function Bar() {}
Bar.barMethod = () => 4;
Object.setPrototypeOf(Foo, Bar);
`),
		c.Expected(
			`/** @constructor */
function Foo() {}
Foo.a = () => { return 3; };
/** @constructor */
function Bar() {}
Bar.a = () => { return 4; };
Object.setPrototypeOf(Foo, Bar);
`))
	// now trying to reference Foo.barMethod will not work, and will call barMethod instead.
	// AmbiguateProperties currently ignores this case
}

// Java: AmbiguatePropertiesTest.testComputedPropertyInObjectLiteral_notAmbiguated()
func (c *AmbiguatePropertiesTest) TestComputedPropertyInObjectLiteral_notAmbiguated() {
	c.TestSame("const obj = {['a']: 3, b: 4}")
}

// Java: AmbiguatePropertiesTest.testQuotedPropertyInObjectLiteralReservesPropertyName()
func (c *AmbiguatePropertiesTest) TestQuotedPropertyInObjectLiteralReservesPropertyName() {
	c.Test(
		"const obj = {'a': 3}; class C { method() {}}",
		// `method` is renamed to `b`, not `a`, to avoid colliding with the computed prop.
		// This is just a heuristic; JSCompiler cannot statically determine all string
		// property names.
		"const obj = {'a': 3}; class C { b() {}}")
}

// Java: AmbiguatePropertiesTest.testQuotedMemberFnInObjectLiteralReservesPropertyName()
func (c *AmbiguatePropertiesTest) TestQuotedMemberFnInObjectLiteralReservesPropertyName() {
	c.Test(
		"const obj = {'a'() {}}; class C { method() {}}",
		"const obj = {'a'() {}}; class C { b() {}}")
}

// Java: AmbiguatePropertiesTest.testComputedPropertyInObjectLiteralReservesPropertyName()
func (c *AmbiguatePropertiesTest) TestComputedPropertyInObjectLiteralReservesPropertyName() {
	c.Test(
		"const obj = {['a']: 3}; class C { method() {}}",
		"const obj = {['a']: 3}; class C { b() {}}")
}

// Java: AmbiguatePropertiesTest.testMemberFnInObjectLiteralPreventsPropertyRenaming()
func (c *AmbiguatePropertiesTest) TestMemberFnInObjectLiteralPreventsPropertyRenaming() {
	// Don't rename the class member 'm' because the object literal type is invalidating, and
	// also has a property 'm'
	c.TestSame("const obj = {m() {}}; class C {m() {}}")
}

// Java: AmbiguatePropertiesTest.testSimplePropInObjectLiteralPreventsPropertyRenaminge()
func (c *AmbiguatePropertiesTest) TestSimplePropInObjectLiteralPreventsPropertyRenaminge() {
	// Don't rename the class member 'm' because the object literal type is invalidating, and
	// also has a property 'm'
	c.TestSame("const obj = {m: 0}; class C {m() {}}")
}

// Java: AmbiguatePropertiesTest.testObjectPatternDeclarationWithStringKey_ambiguated()
func (c *AmbiguatePropertiesTest) TestObjectPatternDeclarationWithStringKey_ambiguated() {
	c.Test(
		`class Foo {
  method() {}
}
const {method} = new Foo();
`,
		`class Foo {
  a() {}
}
const {a: method} = new Foo();
`)
}

// Java: AmbiguatePropertiesTest.testObjectPatternDeclarationWithStringKeWithDefault_ambiguated()
func (c *AmbiguatePropertiesTest) TestObjectPatternDeclarationWithStringKeWithDefault_ambiguated() {
	c.Test(
		`class Foo {
  method() {}
}
const {method = () => 3} = new Foo();
`,
		`class Foo {
  a() {}
}
const {
  a: method = () => { return 3; }
} = new Foo();
`)
}

// Java: AmbiguatePropertiesTest.testNestedObjectPatternWithStringKey_ambiguated()
func (c *AmbiguatePropertiesTest) TestNestedObjectPatternWithStringKey_ambiguated() {
	c.Test(
		`class Foo {
  method() {}
}
const {foo: {method}} = {foo: new Foo()};
`,
		`class Foo {
  a() {}
}
// note: we rename the 'method' access but not 'foo', because we don't try ambiguating
// properties on object literal types.
const {foo: {a: method}} = {foo: new Foo()};
`)
}

// Java: AmbiguatePropertiesTest.testObjectPatternParameterWithStringKey_ambiguated()
func (c *AmbiguatePropertiesTest) TestObjectPatternParameterWithStringKey_ambiguated() {
	c.Test(
		`class Foo {
  method() {}
}
/** @param {!Foo} foo */
function f({method}) {}
`,
		`class Foo {
  a() {}
}
/** @param {!Foo} foo */
function f({a: method}) {}
`)
}

// Java: AmbiguatePropertiesTest.testObjectPatternQuotedStringKey_notAmbiguated()
func (c *AmbiguatePropertiesTest) TestObjectPatternQuotedStringKey_notAmbiguated() {
	// this emits a warning for a computed property access on a struct, since property ambiguation
	// will break this code.
	c.IgnoreWarnings(DiagnosticGroups_CHECK_TYPES)
	c.Test(
		`class Foo {
  method() {}
}
const {'method': method} = new Foo();
`,
		`class Foo {
  a() {}
}
const {'method': method} = new Foo();
`)
}

// Java: AmbiguatePropertiesTest.testObjectPatternComputedProperty_notAmbiguated()
func (c *AmbiguatePropertiesTest) TestObjectPatternComputedProperty_notAmbiguated() {
	// this emits a warning for a computed property access on a struct, since property ambiguation
	// will break this code.
	c.IgnoreWarnings(DiagnosticGroups_CHECK_TYPES)
	c.Test(
		`class Foo {
  method() {}
}
const {['method']: method} = new Foo();
`,
		`class Foo {
  a() {}
}
const {['method']: method} = new Foo();
`)
}

// Java: AmbiguatePropertiesTest.testMixinPropertiesNotAmbiguated()
func (c *AmbiguatePropertiesTest) TestMixinPropertiesNotAmbiguated() {
	c.Test(
		`class A {
  constructor() {
    this.aProp = 'aProp';
  }
}

/**
 * @template T
 * @param {function(new: T)} baseType
 * @return {?}
 */
function mixinX(baseType) {
  return class extends baseType {
    constructor() {
      super();
      this.x = 'x';
    }
  };
}
/** @constructor @extends {A} */
const BSuper = mixinX(A);

class B extends BSuper {
  constructor() {
    super();
    this.bProp = 'bProp';
  }
}
`,
		`class A {
  constructor() {
    this.a = 'aProp';
  }
}

/**
 * @template T
 * @param {function(new: T)} baseType
 * @return {?}
 */
function mixinX(baseType) {
  return class extends baseType {
    constructor() {
      super();
      this.x = 'x';
    }
  };
}
/** @constructor @extends {A} */
const BSuper = mixinX(A);

class B extends BSuper {
  constructor() {
    super();
// Note that ambiguating 'bProp' => 'a' would be incorrect because BSuper extends A via
// the mixin, so 'bProp' would conflict with 'aProp'. JSCompiler doesn't actually track
// that B extends A. Instead it conservatively invalidates all properties on B because
// BSuper is from a mixin, not a class literal.
    this.bProp = 'bProp';
  }
}
`)
}

// Java: AmbiguatePropertiesTest.testMixinPrototypePropertiesNotAmbiguated()
func (c *AmbiguatePropertiesTest) TestMixinPrototypePropertiesNotAmbiguated() {
	c.Test(
		`class A {
  aMethod() {}
}

function mixinX(baseType) {
  return class extends baseType {
    x() {}
  };
}
/** @constructor @extends {A} */
const BSuper = mixinX(A);

class B extends BSuper {
  bMethod() {}
}
`,
		`class A {
  a() {}
}

function mixinX(baseType) {
  return class extends baseType {
    x() {}
  };
}
/** @constructor @extends {A} */
const BSuper = mixinX(A);

class B extends BSuper {
// Note that ambiguating 'bMethod' => 'a' would be incorrect because BSuper extends A
// via the mixin, so 'bMethod' would conflict with 'aMethod'.
  bMethod() {}
}
`)
}

// Java: AmbiguatePropertiesTest.testMixinConstructorPropertiesNotAmbiguated()
func (c *AmbiguatePropertiesTest) TestMixinConstructorPropertiesNotAmbiguated() {
	c.Test(
		`class A {
  static aMethod() {}
}

/**
 * @template T
 * @param {function(new: T)} baseType
 * @return {?}
 */
function mixinX(baseType) {
  return class extends baseType {
    static x() {}
  };
}
/** @constructor @extends {A} */
const BSuper = mixinX(A);

class B extends BSuper {
  static bMethod() {}
}
`,
		`class A {
  static a() {}
}

/**
 * @template T
 * @param {function(new: T)} baseType
 * @return {?}
 */
function mixinX(baseType) {
  return class extends baseType {
    static x() {}
  };
}
/** @constructor @extends {A} */
const BSuper = mixinX(A);

class B extends BSuper {
  static bMethod() {}
}
`)
}

// Java: AmbiguatePropertiesTest.testObjectUsedAsMapIsInvalidating()
func (c *AmbiguatePropertiesTest) TestObjectUsedAsMapIsInvalidating() {
	c.IgnoreWarnings(DiagnosticGroups_MISSING_PROPERTIES)
	c.Test(
		`/** @param {!Object<string>} obj$jscomp$6 */
function f(obj$jscomp$6) {
  return obj$jscomp$6.x;
}
class OtherClass {
  constructor() {
    this.y = 0;
  }
}
`,
		`/** @param {!Object<string>} obj$jscomp$6 */
function f(obj$jscomp$6) {
// Can't ambiguate this 'x' to 'a' since we don't track different Object types uniquely
// and OtherClass is also assignable to !Object<string>
  return obj$jscomp$6.x;
}
class OtherClass {
  constructor() {
    this.a = 0;
  }
}
`)
}

// Java: AmbiguatePropertiesTest.testAmbiguateMethodAddedToStringPrototype_accessedOffScalar()
func (c *AmbiguatePropertiesTest) TestAmbiguateMethodAddedToStringPrototype_accessedOffScalar() {
	c.Test(
		`/** @return {string} */
String.prototype.customMethod = function() { return ''; };
'FOOBAR'.customMethod();
class Bar {
  bar() {}
}
new Bar().bar();
`,
		`/** @return {string} */
String.prototype.a = function() { return ''; };
'FOOBAR'.a();
class Bar {
  a() {}
}
new Bar().a();
`)
}

// Java: AmbiguatePropertiesTest.testAmbiguateMethodAddedToStringPrototype()
func (c *AmbiguatePropertiesTest) TestAmbiguateMethodAddedToStringPrototype() {
	c.Test(
		`/** @return {string} */
String.prototype.customMethod = function() { return ''; };
new String('FOOBAR').customMethod();
class Bar {
  bar() {}
}
new Bar().bar();
`,
		`/** @return {string} */
String.prototype.a = function() { return ''; };
new String('FOOBAR').a();
class Bar {
  a() {}
}
new Bar().a();
`)
}

// Java: AmbiguatePropertiesTest.testEs2022StaticInitializationBlock_notAmbiguated()
func (c *AmbiguatePropertiesTest) TestEs2022StaticInitializationBlock_notAmbiguated() {
	c.TestSame(
		`class Foo {
  static {
    alert('Nothing to ambiguate');
  }
}
`)
}
