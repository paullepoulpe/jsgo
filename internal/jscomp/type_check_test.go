// Port of com/google/javascript/jscomp/TypeCheckTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/TypeCheckTest.java)

/*
 * Copyright 2006 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/rhinotesting"
	"strings"
	"testing"
)

// TypeCheckTest tests the type checker.
type TypeCheckTest struct {
	TypeCheckTestCaseBase
}

// Java: TypeCheckTest()
func NewTypeCheckTest() *TypeCheckTest {
	c := &TypeCheckTest{}
	c.InitTypeCheckTestCaseBase(c)
	return c
}

func TestTypeCheckTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewTypeCheckTest() })
}

var TypeCheckTest_BOUNDED_GENERICS_USE_MSG = "Bounded generic semantics are currently still in development"

var TypeCheckTest_SUGGESTION_CLASS = `/** @constructor
 */
function Suggest() {}
Suggest.prototype.a = 1;
Suggest.prototype.veryPossible = 1;
Suggest.prototype.veryPossible2 = 1;
`

var TypeCheckTest_ILLEGAL_PROPERTY_CREATION_MESSAGE = `Cannot add a property to a struct instance after it is constructed. (If you already declared the property, make sure to give it a type.)
`

var TypeCheckTest_PRIMITIVE_ASSERT_DEFS = `/**
 * @param {T} p
 * @return {U}
 * @template T
 * @template U :=
       mapunion(T, (E) => cond(sub(E, union('null', 'undefined')), none(), E)) =:
 * @closurePrimitive {asserts.truthy}
 */
function assertTruthy(p) { return p; }
/**
 * @param {*} p
 * @return {string}
 * @closurePrimitive {asserts.matchesReturn}
 */
function assertString(p) { return /** @type {string} */ (p); }
`

var TypeCheckTest_EXTERNS_WITH_IARRAYLIKE_DECLS = `/**
 * @constructor @implements IArrayLike<number>
 */
function Int8Array(length, opt_byteOffset, opt_length) {}
/** @type {number} */
Int8Array.prototype.length;
/**
* @constructor
* @extends {Int8Array}
*/
function Int8Array2(len) {};
/**
 * @interface
 * @extends {IArrayLike<number>}
 */
function IArrayLike2(){}

/**
 * @constructor
 * @implements {IArrayLike2}
 */
function Int8Array3(len) {};
/** @type {number} */
Int8Array3.prototype.length;
/**
 * @interface
 * @extends {IArrayLike<VALUE3>}
 * @template VALUE3
 */
function IArrayLike3(){}
/**
 * @constructor
 * @implements {IArrayLike3<number>}
 */
function Int8Array4(length) {};
/** @type {number} */
Int8Array4.prototype.length;
/**
 * @interface
 * @extends {IArrayLike<VALUE2>}
 * @template VALUE2
 */
function IArrayLike4(){}
/**
 * @interface
 * @extends {IArrayLike4<boolean>}
 */
function IArrayLike5(){}
/**
 * @constructor
 * @implements {IArrayLike5}
 */
function BooleanArray5(length) {};
/** @type {number} */
BooleanArray5.prototype.length;
`

var TypeCheckTest_EXTERNS_WITH_IOBJECT_DECLS = `/**
 * @constructor
 * @implements IObject<(string|number), number>
 */
function Object2() {}
/**
 * @constructor @struct
 * @implements IObject<number, number>
 */
function Object3() {}
`

var TypeCheckTest_MIXIN_DEFINITIONS = `/** @constructor */
function MyElement() {
  /** @type {string} */
  this.elemprop = 'asdf';
}
/** @record */
function Toggle() {}
/**
 * @param {string} x
 * @return {string}
 */
Toggle.prototype.foobar = function(x) {};
/**
 * @template T
 * @param {function(new:T)} superclass
 * @suppress {checkTypes} // TODO(b/74120976): fix bug and remove suppression
 */
function addToggle(superclass) {
  /**
   * @constructor
   * @extends {superclass}
   * @implements {Toggle}
   */
  function Clazz() {
    superclass.apply(this, arguments);
  }
  Clazz.prototype = Object.create(superclass.prototype);
  /** @override */
  Clazz.prototype.foobar = function(x) { return 'foobar ' + x; };
  return Clazz;
}
`

// Java: TypeCheckTest.testInitialTypingScope()
func (c *TypeCheckTest) TestInitialTypingScope() {
	var s *TypedScope = NewTypedScopeCreatorWithCodingConvention(c.compiler, CodingConventions_GetDefault()).CreateInitialScope(rhino.NewNodeWith2Children(rhino.Token_ROOT, rhino.NewNode(rhino.Token_ROOT), rhino.NewNode(rhino.Token_ROOT)))
	c.AssertTypeEquals(c.GetNativeArrayConstructorType(), s.GetVar("Array").GetType())
	c.AssertTypeEquals(c.GetNativeBooleanObjectConstructorType(), s.GetVar("Boolean").GetType())
	c.AssertTypeEquals(c.GetNativeDateConstructorType(), s.GetVar("Date").GetType())
	c.AssertTypeEquals(c.GetNativeNumberObjectConstructorType(), s.GetVar("Number").GetType())
	c.AssertTypeEquals(c.GetNativeObjectConstructorType(), s.GetVar("Object").GetType())
	c.AssertTypeEquals(c.GetNativeRegexpConstructorType(), s.GetVar("RegExp").GetType())
	c.AssertTypeEquals(c.GetNativeStringObjectConstructorType(), s.GetVar("String").GetType())
}

// Java: TypeCheckTest.testPrivateType()
func (c *TypeCheckTest) TestPrivateType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @private {number} */ var x = false;").
		AddDiagnostic(`initializing variable
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckTest.testTypeCheck1()
func (c *TypeCheckTest) TestTypeCheck1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@return {void}*/function foo(){ if (foo()) return; }").
		Run()
}

// Java: TypeCheckTest.testTypeCheck2()
func (c *TypeCheckTest) TestTypeCheck2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@return {void}*/function foo(){ var x=foo(); x--; }").
		AddDiagnostic(`increment/decrement
found   : undefined
required: number
`).
		Run()
}

// Java: TypeCheckTest.testTypeCheck4()
func (c *TypeCheckTest) TestTypeCheck4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@return {void}*/function foo(){ !foo(); }").
		Run()
}

// Java: TypeCheckTest.testTypeCheck5()
func (c *TypeCheckTest) TestTypeCheck5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@return {void}*/function foo(){ var a = +foo(); }").
		Run()
}

// Java: TypeCheckTest.testTypeCheck6()
func (c *TypeCheckTest) TestTypeCheck6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {void} */
function foo(){
  /** @type {undefined|number} */
  var a;
  if (a == foo())
    return;
}
`).
		Run()
}

// Java: TypeCheckTest.testTypeCheck8()
func (c *TypeCheckTest) TestTypeCheck8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {void} */
function foo() {
  do {} while (foo());
}
`).
		Run()
}

// Java: TypeCheckTest.testTypeCheck9()
func (c *TypeCheckTest) TestTypeCheck9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {void} */
function foo() {
  while (foo());
}
`).
		Run()
}

// Java: TypeCheckTest.testTypeCheck10()
func (c *TypeCheckTest) TestTypeCheck10() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {void} */
function foo() {
  for (;foo(););
}
`).
		Run()
}

// Java: TypeCheckTest.testTypeCheck11()
func (c *TypeCheckTest) TestTypeCheck11() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {!Number} */
var a;
/** @type {!String} */
var b;
a = b;
`).
		AddDiagnostic(`assignment
found   : String
required: Number
`).
		Run()
}

// Java: TypeCheckTest.testTypeCheck12()
func (c *TypeCheckTest) TestTypeCheck12() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@return {!Object}*/function foo(){var a = 3^foo();}").
		AddDiagnostic(`bad right operand to bitwise operator
found   : Object
required: (boolean|null|number|string|undefined)
`).
		Run()
}

// Java: TypeCheckTest.testTypeCheck13()
func (c *TypeCheckTest) TestTypeCheck13() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@type {!Number|!String}*/var i; i=/xx/;").
		AddDiagnostic(`assignment
found   : RegExp
required: (Number|String)
`).
		Run()
}

// Java: TypeCheckTest.testTypeCheck14()
func (c *TypeCheckTest) TestTypeCheck14() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@param {?} opt_a*/function foo(opt_a){}").
		Run()
}

// Java: TypeCheckTest.testTypeCheck15()
func (c *TypeCheckTest) TestTypeCheck15() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@type {Number|null} */var x;x=null;x=10;").
		AddDiagnostic(`assignment
found   : number
required: (Number|null)
`).
		Run()
}

// Java: TypeCheckTest.testTypeCheck16()
func (c *TypeCheckTest) TestTypeCheck16() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@type {Number|null} */var x='';").
		AddDiagnostic(`initializing variable
found   : string
required: (Number|null)
`).
		Run()
}

// Java: TypeCheckTest.testTypeCheck17()
func (c *TypeCheckTest) TestTypeCheck17() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @return {Number}
 * @param {Number} opt_foo
 */
function a(opt_foo){
  return /** @type {Number} */ (opt_foo);
}
`).
		Run()
}

// Java: TypeCheckTest.testTypeCheck18()
func (c *TypeCheckTest) TestTypeCheck18() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**@return {RegExp}
*/
 function a(){return new RegExp();}`).
		Run()
}

// Java: TypeCheckTest.testTypeCheck19()
func (c *TypeCheckTest) TestTypeCheck19() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**@return {Array}
*/
 function a(){return new Array();}`).
		Run()
}

// Java: TypeCheckTest.testTypeCheck20()
func (c *TypeCheckTest) TestTypeCheck20() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**@return {Date}
*/
 function a(){return new Date();}`).
		Run()
}

// Java: TypeCheckTest.testTypeCheckBasicDowncast()
func (c *TypeCheckTest) TestTypeCheckBasicDowncast() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */function foo() {}
/** @type {Object} */ var bar = new foo();
`).
		Run()
}

// Java: TypeCheckTest.testTypeCheckNoDowncastToNumber()
func (c *TypeCheckTest) TestTypeCheckNoDowncastToNumber() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */function foo() {}
/** @type {!Number} */ var bar = new foo();
`).
		AddDiagnostic(`initializing variable
found   : foo
required: Number
`).
		Run()
}

// Java: TypeCheckTest.testTypeCheck21()
func (c *TypeCheckTest) TestTypeCheck21() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {Array<String>} */var foo;").
		Run()
}

// Java: TypeCheckTest.testTypeCheck22()
func (c *TypeCheckTest) TestTypeCheck22() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {Element|Object} p */
function foo(p){}
/** @constructor */function Element(){}
/** @type {Element|Object} */var v;
foo(v);
`).
		Run()
}

// Java: TypeCheckTest.testTypeCheck23()
func (c *TypeCheckTest) TestTypeCheck23() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {(Object|Null)} */var foo; foo = null;").
		Run()
}

// Java: TypeCheckTest.testTypeCheck24()
func (c *TypeCheckTest) TestTypeCheck24() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */function MyType(){}
/** @type {(MyType|Null)} */var foo; foo = null;
`).
		Run()
}

// Java: TypeCheckTest.testTypeCheck25()
func (c *TypeCheckTest) TestTypeCheck25() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function foo(/** {a: number} */ obj) {};
foo({b: 'abc'});
`).
		AddDiagnostic(`actual parameter 1 of foo does not match formal parameter
found   : {
  a: (number|undefined),
  b: string
}
required: {a: number}
missing : []
mismatch: [a]
`).
		Run()
}

// Java: TypeCheckTest.testTypeCheck26()
func (c *TypeCheckTest) TestTypeCheck26() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function foo(/** {a: number} */ obj) {};
foo({a: 'abc'});
`).
		AddDiagnostic(`actual parameter 1 of foo does not match formal parameter
found   : {a: (number|string)}
required: {a: number}
missing : []
mismatch: [a]
`).
		Run()
}

// Java: TypeCheckTest.testTypeCheck27()
func (c *TypeCheckTest) TestTypeCheck27() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function foo(/** {a: number} */ obj) {};
foo({a: 123});
`).
		Run()
}

// Java: TypeCheckTest.testTypeCheck28()
func (c *TypeCheckTest) TestTypeCheck28() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function foo(/** ? */ obj) {};
foo({a: 123});
`).
		Run()
}

// Java: TypeCheckTest.testPropertyOnATypedefGetsResolved()
func (c *TypeCheckTest) TestPropertyOnATypedefGetsResolved() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`const a = {};
/** @typedef {{prop: string}} */
  a.b;
/** @typedef {string} */
  a.b.c;
const abAlias = a.b;
var /** !abAlias.whatIsThis */ y = 'foo';
`).
		AddDiagnosticType(RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR).
		Run()
}

// Java: TypeCheckTest.testPropertyOnTypedefAliasesGetsRecognized()
func (c *TypeCheckTest) TestPropertyOnTypedefAliasesGetsRecognized() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`const a = {};
/** @typedef {{prop: string}} */
  a.b;
/** @typedef {string} */
  a.b.c;
const aAlias = a;
var /** aAlias.b.c */ x;
`).
		AddDiagnosticType(RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR).
		Run()
}

// Java: TypeCheckTest.testPropertyOnTypedefAliasesGetsRecognized2()
func (c *TypeCheckTest) TestPropertyOnTypedefAliasesGetsRecognized2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`const a = {};
/** @typedef {{prop: string}} */
  a.b;
/** @typedef {string} */
  a.b.c;
const aAlias = a;
var /** aAlias.b */ x = {prop: 'foo'};
`).
		Run()
}

// Java: TypeCheckTest.testTypeCheckInlineReturns()
func (c *TypeCheckTest) TestTypeCheckInlineReturns() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function /** string */ foo(x) { return x; }
var /** number */ a = foo('abc');
`).
		AddDiagnostic(`initializing variable
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testTypeCheckDefaultExterns()
func (c *TypeCheckTest) TestTypeCheckDefaultExterns() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddArray().Build()).
		AddSource(`/** @param {string} x */
function f(x) {}
f([].length);
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testTypeCheckCustomExterns()
func (c *TypeCheckTest) TestTypeCheckCustomExterns() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns("/** @type {boolean} */ Array.prototype.oogabooga;").
		IncludeDefaultExterns().
		AddSource(`/** @param {string} x */ function f(x) {}
f([].oogabooga);
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : boolean
required: string
`).
		Run()
}

// Java: TypeCheckTest.testTypeCheckCustomExterns2()
func (c *TypeCheckTest) TestTypeCheckCustomExterns2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns("/** @enum {string} */ var Enum = {FOO: 1, BAR: 1};").
		IncludeDefaultExterns().
		AddSource("/** @param {Enum} x */ function f(x) {} f(Enum.FOO); f(true);").
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : boolean
required: Enum<string>
`).
		Run()
}

// Java: TypeCheckTest.testDontCrashOnRecursiveTemplateReference()
func (c *TypeCheckTest) TestDontCrashOnRecursiveTemplateReference() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor @implements {Iterable<VALUE>} @template VALUE */
function Foo() {
  /** @type {!Map<VALUE, VALUE>} */ this.map = new Map;
}
/** @override @return {!Iterator<VALUE>} */
Foo.prototype[Symbol.iterator] = function() {};
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testUnionOfFunctionAndType()
func (c *TypeCheckTest) TestUnionOfFunctionAndType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {null|(function(Number):void)} */ var a;
/** @type {(function(Number):void)|null} */ var b = null; a = b;
`).
		Run()
}

// Java: TypeCheckTest.testOptionalParameterComparedToUndefined()
func (c *TypeCheckTest) TestOptionalParameterComparedToUndefined() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param  {Number} opt_a */function foo(opt_a)
{if (opt_a==undefined) var b = 3;}
`).
		Run()
}

// Java: TypeCheckTest.testOptionalAllType()
func (c *TypeCheckTest) TestOptionalAllType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {*} opt_x */function f(opt_x) { return opt_x }
/** @type {*} */var y;
f(y);
`).
		Run()
}

// Java: TypeCheckTest.testUnaryPlusWithAllType()
func (c *TypeCheckTest) TestUnaryPlusWithAllType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {*} */var x; +x;").
		Run()
}

// Java: TypeCheckTest.testUnaryPlusWithNoType()
func (c *TypeCheckTest) TestUnaryPlusWithNoType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {!null} */var x; +x;").
		Run()
}

// Java: TypeCheckTest.testUnaryPlusWithUnknownType()
func (c *TypeCheckTest) TestUnaryPlusWithUnknownType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {?} */var x; +x;").
		Run()
}

// Java: TypeCheckTest.testOptionalUnknownNamedType()
func (c *TypeCheckTest) TestOptionalUnknownNamedType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {!T} opt_x
@return {undefined} */
function f(opt_x) { return opt_x; }
/** @constructor */var T = function() {};
`).
		AddDiagnostic(`inconsistent return type
found   : (T|undefined)
required: undefined
`).
		Run()
}

// Java: TypeCheckTest.testOptionalArgFunctionParam()
func (c *TypeCheckTest) TestOptionalArgFunctionParam() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {function(number=)} a */
function f(a) {a()};
`).
		Run()
}

// Java: TypeCheckTest.testOptionalArgFunctionParam2()
func (c *TypeCheckTest) TestOptionalArgFunctionParam2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {function(number=)} a */
function f(a) {a(3)};
`).
		Run()
}

// Java: TypeCheckTest.testOptionalArgFunctionParam3()
func (c *TypeCheckTest) TestOptionalArgFunctionParam3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {function(number=)} a */
function f(a) {a(undefined)};
`).
		Run()
}

// Java: TypeCheckTest.testOptionalArgFunctionParam4()
func (c *TypeCheckTest) TestOptionalArgFunctionParam4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {function(number=)} a */
function f(a) {a(3,4)};
`).
		AddDiagnostic(`Function a: called with 2 argument(s). Function requires at least 0 argument(s) and no more than 1 argument(s).
`).
		Run()
}

// Java: TypeCheckTest.testOptionalArgFunctionParamError()
func (c *TypeCheckTest) TestOptionalArgFunctionParamError() {
	expectedWarning := "Bad type annotation. variable length argument must be last." + parsing.JsDocInfoParser_BAD_TYPE_WIKI_LINK
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {function(...number, number=)} a */
function f(a) {};
`).
		AddDiagnostic(expectedWarning).
		Run()
}

// Java: TypeCheckTest.testOptionalNullableArgFunctionParam()
func (c *TypeCheckTest) TestOptionalNullableArgFunctionParam() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {function(?number=)} a */
function f(a) {a()};
`).
		Run()
}

// Java: TypeCheckTest.testOptionalNullableArgFunctionParam2()
func (c *TypeCheckTest) TestOptionalNullableArgFunctionParam2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {function(?number=)} a */
function f(a) {a(null)};
`).
		Run()
}

// Java: TypeCheckTest.testOptionalNullableArgFunctionParam3()
func (c *TypeCheckTest) TestOptionalNullableArgFunctionParam3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {function(?number=)} a */
function f(a) {a(3)};
`).
		Run()
}

// Java: TypeCheckTest.testOptionalArgFunctionReturn()
func (c *TypeCheckTest) TestOptionalArgFunctionReturn() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {function(number=)} */
function f() { return function(opt_x) { }; };
f()()
`).
		Run()
}

// Java: TypeCheckTest.testOptionalArgFunctionReturn2()
func (c *TypeCheckTest) TestOptionalArgFunctionReturn2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {function(Object=)} */
function f() { return function(opt_x) { }; };
f()({})
`).
		Run()
}

// Java: TypeCheckTest.testBooleanType()
func (c *TypeCheckTest) TestBooleanType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@type {boolean} */var x = 1 < 2;").
		Run()
}

// Java: TypeCheckTest.testBooleanReduction1()
func (c *TypeCheckTest) TestBooleanReduction1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@type {string} */var x; x = null || \"a\";").
		Run()
}

// Java: TypeCheckTest.testBooleanReduction1NullishCoalesce()
func (c *TypeCheckTest) TestBooleanReduction1NullishCoalesce() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@type {string} */var x; x = null ?? \"a\";").
		Run()
}

// Java: TypeCheckTest.testLogicalAssignment()
func (c *TypeCheckTest) TestLogicalAssignment() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@type {?} */var a; /** @type {?} */ var b; a||=b;").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@type {?} */var a; /** @type {?} */ var b; a&&=b;").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@type {?} */var a; /** @type {?} */ var b; a??=b;").
		Run()
}

// Java: TypeCheckTest.testAssignOrToNonNumeric()
func (c *TypeCheckTest) TestAssignOrToNonNumeric() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {boolean} x */
 function assignOr(x) {
   x ||= 'a';
 };
`).
		AddDiagnostic(`assignment
found   : string
required: boolean
`).
		Run()
}

// Java: TypeCheckTest.testAssignOrMayOrMayNotAssign()
func (c *TypeCheckTest) TestAssignOrMayOrMayNotAssign() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {string} x */
 function assignOr(x) {
   x ||= 5;
 };
`).
		AddDiagnostic(`assignment
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testAssignAndCheckRHSValid()
func (c *TypeCheckTest) TestAssignAndCheckRHSValid() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {string|undefined} */
var a; a &&= 0;
`).
		AddDiagnostic(`assignment
found   : number
required: (string|undefined)
`).
		Run()
}

// Java: TypeCheckTest.testAssignAndRHSNotExecuted()
func (c *TypeCheckTest) TestAssignAndRHSNotExecuted() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`let /** null */
n = null;
n &&= 'str';
`).
		AddDiagnostic(`assignment
found   : string
required: null
`).
		Run()
}

// Java: TypeCheckTest.testAssignCoalesceNoAssign()
func (c *TypeCheckTest) TestAssignCoalesceNoAssign() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {string} x */
 function assignCoalesce(x) {
   x ??= 5;
 };
`).
		AddDiagnostic(`assignment
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testAssignCoalesceCheckRHSValid()
func (c *TypeCheckTest) TestAssignCoalesceCheckRHSValid() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {string|undefined} */ var a; a ??= 0;").
		AddDiagnostic(`assignment
found   : number
required: (string|undefined)
`).
		Run()
}

// Java: TypeCheckTest.optChain()
func (c *TypeCheckTest) Test_optChain() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@type {?} */var x = a?.b;").
		Run()
}

// Java: TypeCheckTest.testBooleanReduction2()
func (c *TypeCheckTest) TestBooleanReduction2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {string} s
 * @return {string}
 */
(function(s) {
  return ((s == 'a') && s) || 'b';
})
`).
		Run()
}

// Java: TypeCheckTest.testBooleanReduction3()
func (c *TypeCheckTest) TestBooleanReduction3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {string} s
 * @return {string?}
 */
(function(s) {
  return s && null && 3;
})
`).
		Run()
}

// Java: TypeCheckTest.testBooleanReduction4()
func (c *TypeCheckTest) TestBooleanReduction4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {?Object} x
 * @return {?Object}
 */
(function(x) {
  return null || x || null;
})
`).
		Run()
}

// Java: TypeCheckTest.testBooleanReduction4NullishCoalesce()
func (c *TypeCheckTest) TestBooleanReduction4NullishCoalesce() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {?Object} x
 @return {!Object} */
(function(x) { return null ?? x ?? {} ; })
`).
		Run()
}

// Java: TypeCheckTest.testBooleanReduction5()
func (c *TypeCheckTest) TestBooleanReduction5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
* @param {Array|string} x
* @return {string?}
*/
var f = function(x) {
if (!x || typeof x == 'string') {
return x;
}
return null;
};
`).
		Run()
}

// Java: TypeCheckTest.testBooleanReduction6()
func (c *TypeCheckTest) TestBooleanReduction6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
* @param {Array|string|null} x
* @return {string?}
*/
var f = function(x) {
if (!(x && typeof x != 'string')) {
return x;
}
return null;
};
`).
		Run()
}

// Java: TypeCheckTest.testBooleanReduction7()
func (c *TypeCheckTest) TestBooleanReduction7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */var T = function() {};
/**
* @param {Array|T} x
* @return {null}
*/
var f = function(x) {
if (!x) {
return x;
}
return null;
};
`).
		Run()
}

// Java: TypeCheckTest.testNullishCoalesceWithAndExpressions()
func (c *TypeCheckTest) TestNullishCoalesceWithAndExpressions() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {!Array} s
    @return {string} */
(function(s) { return (s && undefined) ?? (s && "a"); })
`).
		Run()
}

// Java: TypeCheckTest.testNullishCoalesceShortCutsOnZero()
func (c *TypeCheckTest) TestNullishCoalesceShortCutsOnZero() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@type {number} */var x; x = 0 ?? \"a\";").
		Run()
}

// Java: TypeCheckTest.testNullishCoalesceUnionsOperatorTypesWithNullRemoved()
func (c *TypeCheckTest) TestNullishCoalesceUnionsOperatorTypesWithNullRemoved() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
* @param {string|null} x
* @return {string|number}
*/
var f = function(x) {
return x ?? 3;
};
`).
		Run()
}

// Java: TypeCheckTest.testNullishCoalesceShortCutsOnFalse()
func (c *TypeCheckTest) TestNullishCoalesceShortCutsOnFalse() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@type {boolean} */var x; x = false ?? undefined ?? [];").
		Run()
}

// Java: TypeCheckTest.testNullishCoalesceChaining()
func (c *TypeCheckTest) TestNullishCoalesceChaining() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@type {number} */var x; x = null ?? undefined ?? 3;").
		Run()
}

// Java: TypeCheckTest.testNullishCoalesceOptionalNullable()
func (c *TypeCheckTest) TestNullishCoalesceOptionalNullable() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {?Object=} optionalNullableObject
 * @return {!Object}
 */
function f(optionalNullableObject = undefined) {
 return optionalNullableObject ?? {}; // always returns !Object type
};
`).
		Run()
}

// Java: TypeCheckTest.testNullishCoalesceUseNonNullFunction()
func (c *TypeCheckTest) TestNullishCoalesceUseNonNullFunction() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {!Object} x
 * @return {!Object}
 */
 function useNonNull(x) {
   return x;
 };
/** @type {?Object} */ var x;
x ?? useNonNull(x);
`).
		AddDiagnostic(`actual parameter 1 of useNonNull does not match formal parameter
found   : null
required: Object
`).
		Run()
}

// Java: TypeCheckTest.testNullishCoalesceLeftTypeNotDefined()
func (c *TypeCheckTest) TestNullishCoalesceLeftTypeNotDefined() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {number} */ var a = x ?? 1;").
		Run()
}

// Java: TypeCheckTest.testNullAnd()
func (c *TypeCheckTest) TestNullAnd() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {null} */var x;
/** @type {number} */var r = x && x;
`).
		AddDiagnostic(`initializing variable
found   : null
required: number
`).
		Run()
}

// Java: TypeCheckTest.testNullOr()
func (c *TypeCheckTest) TestNullOr() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {null} */var x;
/** @type {number} */var r = x || x;
`).
		AddDiagnostic(`initializing variable
found   : null
required: number
`).
		Run()
}

// Java: TypeCheckTest.testBooleanPreservation1()
func (c *TypeCheckTest) TestBooleanPreservation1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**@type {string} */var x = "a";
x = ((x == "a") && x) || x == "b";
`).
		AddDiagnostic(`assignment
found   : (boolean|string)
required: string
`).
		Run()
}

// Java: TypeCheckTest.testBooleanPreservation2()
func (c *TypeCheckTest) TestBooleanPreservation2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**@type {string} */var x = "a";
x = (x == "a") || x;
`).
		AddDiagnostic(`assignment
found   : (boolean|string)
required: string
`).
		Run()
}

// Java: TypeCheckTest.testBooleanPreservation3()
func (c *TypeCheckTest) TestBooleanPreservation3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {Function?} x
 @return {boolean?} */
function f(x) { return x && x == "a"; }
`).
		AddDiagnostic(`condition always evaluates to false
left : Function
right: string
`).
		Run()
}

// Java: TypeCheckTest.testBooleanPreservation4()
func (c *TypeCheckTest) TestBooleanPreservation4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {Function?|boolean} x
 @return {boolean} */
function f(x) { return x && x == "a"; }
`).
		AddDiagnostic(`inconsistent return type
found   : (boolean|null)
required: boolean
`).
		Run()
}

// Java: TypeCheckTest.testWellKnownSymbolAccess1()
func (c *TypeCheckTest) TestWellKnownSymbolAccess1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {Array<string>} x
 */
function f(x) {
  const iter = x[Symbol.iterator]();
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testWellKnownSymbolAccess2()
func (c *TypeCheckTest) TestWellKnownSymbolAccess2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {IObject<string, number>} x
 */
function f(x) {
  const iter = x[Symbol.iterator]();
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testSymbolComparison1()
func (c *TypeCheckTest) TestSymbolComparison1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {symbol} x
 * @param {symbol} y
*/
function f(x, y) { return x === y; }
`).
		Run()
}

// Java: TypeCheckTest.testSymbolComparison2()
func (c *TypeCheckTest) TestSymbolComparison2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {symbol} x
 * @param {symbol} y
*/
function f(x, y) { return x == y; }
`).
		Run()
}

// Java: TypeCheckTest.testSymbolComparison3()
func (c *TypeCheckTest) TestSymbolComparison3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {symbol} x // primitive
 * @param {!Symbol} y // object
*/
function f(x, y) { return x == y; }
`).
		Run()
}

// Java: TypeCheckTest.testSymbolComparison4()
func (c *TypeCheckTest) TestSymbolComparison4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {symbol} x // primitive
 * @param {!Symbol} y // object
*/
function f(x, y) { return x === y; }
`).
		AddDiagnostic(`condition always evaluates to false
left : symbol
right: Symbol
`).
		Run()
}

// Java: TypeCheckTest.testSymbolComparison5()
func (c *TypeCheckTest) TestSymbolComparison5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {symbol} x // primitive
 * @param {(!Symbol|null)} y // object
*/
function f(x, y) { return x == y; }
`).
		Run()
}

// Java: TypeCheckTest.testSymbolComparison6()
func (c *TypeCheckTest) TestSymbolComparison6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {symbol} x // primitive
 * @param {(!Symbol|null)} y // object
*/
function f(x, y) { return x === y; }
`).
		AddDiagnostic(`condition always evaluates to false
left : symbol
right: (Symbol|null)
`).
		Run()
}

// Java: TypeCheckTest.testSymbolComparison7()
func (c *TypeCheckTest) TestSymbolComparison7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {symbol} x
 * @param {*} y
*/
function f(x, y) { return x == y; }
`).
		Run()
}

// Java: TypeCheckTest.testSymbolComparison8()
func (c *TypeCheckTest) TestSymbolComparison8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {symbol} x
 * @param {*} y
*/
function f(x, y) { return x === y; }
`).
		Run()
}

// Java: TypeCheckTest.testSymbolComparison9()
func (c *TypeCheckTest) TestSymbolComparison9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum {symbol} */ var E = {A:Symbol()};
/**
 * @param {symbol} x
 * @param {E} y
*/
function f(x, y) { return x == y; }
`).
		Run()
}

// Java: TypeCheckTest.testSymbolComparison10()
func (c *TypeCheckTest) TestSymbolComparison10() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum {symbol} */ var E = {A:Symbol()};
/**
 * @param {symbol} x
 * @param {E} y
*/
function f(x, y) { return x === y; }
`).
		Run()
}

// Java: TypeCheckTest.testSymbolComparison11()
func (c *TypeCheckTest) TestSymbolComparison11() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum {!Symbol} */ var E = {A:/** @type {!Symbol} */ (Object(Symbol()))};
/**
 * @param {symbol} x
 * @param {E} y
*/
function f(x, y) { return x == y; }
`).
		Run()
}

// Java: TypeCheckTest.testSymbolComparison12()
func (c *TypeCheckTest) TestSymbolComparison12() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum {!Symbol} */ var E = {A:/** @type {!Symbol} */ (Object(Symbol()))};
/**
 * @param {symbol} x
 * @param {E} y
*/
function f(x, y) { return x === y; }
`).
		AddDiagnostic(`condition always evaluates to false
left : symbol
right: E<Symbol>
`).
		Run()
}

// Java: TypeCheckTest.testSymbolComparison13()
func (c *TypeCheckTest) TestSymbolComparison13() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {symbol} x
 * @param {!Object} y
*/
function f(x, y) { return x == y; }
`).
		Run()
}

// Java: TypeCheckTest.testSymbolComparison14()
func (c *TypeCheckTest) TestSymbolComparison14() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {symbol} x
 * @param {!Object} y
*/
function f(x, y) { return x === y; }
`).
		AddDiagnostic(`condition always evaluates to false
left : symbol
right: Object
`).
		Run()
}

// Java: TypeCheckTest.testTypeOfReduction1()
func (c *TypeCheckTest) TestTypeOfReduction1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {string|number} x
 * @return {string}
 */
function f(x) {
  return typeof x == 'number' ? String(x) : x;
}
`).
		Run()
}

// Java: TypeCheckTest.testTypeOfReduction2()
func (c *TypeCheckTest) TestTypeOfReduction2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {string|number} x
 * @return {string}
 */
function f(x) {
  return typeof x != 'string' ? String(x) : x;
}
`).
		Run()
}

// Java: TypeCheckTest.testWellKnownSymbolTypeofNarrowing_removingOtherTypeFromUnion()
func (c *TypeCheckTest) TestWellKnownSymbolTypeofNarrowing_removingOtherTypeFromUnion() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @type {symbol} */ Symbol.gbigintInternal;
/** @interface */ function GbigintType() {}
/** @const {symbol} */ GbigintType.prototype.__doNotManuallySet;
/** @typedef {!GbigintType|(typeof Symbol.gbigintInternal)} */
let gbigint2;
`).
		AddSource(`/** @return {string|!gbigint2} */
function foo() {
  return 'lalala';
}

/** @return {string|!gbigint2} */
function bar() {
  const x = foo();
  if (typeof x === 'string') {
    return 'hi!';
  }
  return x;
}
`).
		Run()
}

// Java: TypeCheckTest.testWellKnownSymbolTypeofNarrowing_narrowingToSymbol()
func (c *TypeCheckTest) TestWellKnownSymbolTypeofNarrowing_narrowingToSymbol() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @type {symbol} */ Symbol.gbigintInternal;
/** @interface */ function GbigintType() {}
/** @const {symbol} */ GbigintType.prototype.__doNotManuallySet;
/** @typedef {!GbigintType|(typeof Symbol.gbigintInternal)} */
let gbigint2;
`).
		AddSource(`/** @return {string|!gbigint2} */
function foo() {
  return 'lalala';
}

function bar() {
  const x = foo();
  if (typeof x === 'symbol') {
    /** @type {!gbigint2} */
    const test = x;
  }
}
`).
		Run()
}

// Java: TypeCheckTest.testTypeOfReduction3()
func (c *TypeCheckTest) TestTypeOfReduction3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {number|null} x
 * @return {number}
 */
function f(x) {
  return typeof x == 'object' ? 1 : x;
}
`).
		Run()
}

// Java: TypeCheckTest.testTypeOfReduction4()
func (c *TypeCheckTest) TestTypeOfReduction4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {Object|undefined} x
 * @return {Object}
 */
function f(x) {
  return typeof x == 'undefined' ? {} : x;
}
`).
		Run()
}

// Java: TypeCheckTest.testTypeOfReduction5()
func (c *TypeCheckTest) TestTypeOfReduction5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum {string} */ var E = {A: 'a', B: 'b'};
/** @param {!E|number} x
 @return {string} */
function f(x) { return typeof x != 'number' ? x : 'a'; }
`).
		Run()
}

// Java: TypeCheckTest.testTypeOfReduction6()
func (c *TypeCheckTest) TestTypeOfReduction6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddString().Build()).
		AddSource(`/**
 * @param {number|string} x
 * @return {string}
 */
function f(x) {
  return typeof x == 'string' && x.length == 3 ? x : 'a';
}
`).
		Run()
}

// Java: TypeCheckTest.testTypeOfReduction7()
func (c *TypeCheckTest) TestTypeOfReduction7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {string} */
var f = function(x) {
  return typeof x == 'number' ? x : 'a';
}
`).
		AddDiagnostic(`inconsistent return type
found   : (number|string)
required: string
`).
		Run()
}

// Java: TypeCheckTest.testTypeOfReduction11()
func (c *TypeCheckTest) TestTypeOfReduction11() {
	c.TestClosureTypes(`/** @param {Array|string} x
@return {Array} */
function f(x) {
return goog.isObject(x) ? x : [];
}
`, nil)
}

// Java: TypeCheckTest.testTypeOfReduction12()
func (c *TypeCheckTest) TestTypeOfReduction12() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum {string} */ var E = {A: 'a', B: 'b'};
/** @param {E|Array} x
 @return {Array} */
function f(x) { return typeof x == 'object' ? x : []; }
`).
		Run()
}

// Java: TypeCheckTest.testTypeOfReduction13()
func (c *TypeCheckTest) TestTypeOfReduction13() {
	c.TestClosureTypes(`/** @enum {string} */ var E = {A: 'a', B: 'b'};
/** @param {E|Array} x
@return {Array} */
function f(x) { return goog.isObject(x) ? x : []; }
`, nil)
}

// Java: TypeCheckTest.testTypeOfReduction15()
func (c *TypeCheckTest) TestTypeOfReduction15() {
	c.TestClosureTypes(`function f(x) {
  return typeof arguments[0] == 'string' ? arguments[0] : 0;
}
`, nil)
}

// Java: TypeCheckTest.testTypeOfReduction16()
func (c *TypeCheckTest) TestTypeOfReduction16() {
	c.TestClosureTypes(`/** @interface */ function I() {}
/**
 * @param {*} x
 * @return {I}
 */
function f(x) {
  if(goog.isObject(x)) {
    return /** @type {I} */ (x);
  }
  return null;
}
`, nil)
}

// Java: TypeCheckTest.testQualifiedNameReduction1()
func (c *TypeCheckTest) TestQualifiedNameReduction1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var x = {}; /** @type {string?} */ x.a = 'a';
/** @return {string} */ var f = function() {
return x.a ? x.a : 'a'; }
`).
		Run()
}

// Java: TypeCheckTest.testQualifiedNameReduction2()
func (c *TypeCheckTest) TestQualifiedNameReduction2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {string?} a
@constructor */ var T = function(a) {this.a = a};
/** @return {string} */ T.prototype.f = function() {
return this.a ? this.a : 'a'; }
`).
		Run()
}

// Java: TypeCheckTest.testQualifiedNameReduction3()
func (c *TypeCheckTest) TestQualifiedNameReduction3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {string|Array} a
@constructor */ var T = function(a) {this.a = a};
/** @return {string} */ T.prototype.f = function() {
return typeof this.a == 'string' ? this.a : 'a'; }
`).
		Run()
}

// Java: TypeCheckTest.testQualifiedNameReduction5a()
func (c *TypeCheckTest) TestQualifiedNameReduction5a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var x = {/** @type {string} */ a:'b' };
/** @return {string} */ var f = function() {
return x.a; }
`).
		Run()
}

// Java: TypeCheckTest.testQualifiedNameReduction5b()
func (c *TypeCheckTest) TestQualifiedNameReduction5b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var x = {/** @type {number} */ a:12 };
/** @return {string} */
var f = function() {
  return x.a;
}
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testQualifiedNameReduction5c()
func (c *TypeCheckTest) TestQualifiedNameReduction5c() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {string} */ var f = function() {
var x = {/** @type {number} */ a:0 };
return (x.a) ? (x.a) : 'a'; }
`).
		AddDiagnostic(`inconsistent return type
found   : (number|string)
required: string
`).
		Run()
}

// Java: TypeCheckTest.testQualifiedNameReduction6()
func (c *TypeCheckTest) TestQualifiedNameReduction6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {string} */ var f = function() {
var x = {/** @return {string?} */ get a() {return 'a'}};
return x.a ? x.a : 'a'; }
`).
		Run()
}

// Java: TypeCheckTest.testQualifiedNameReduction7()
func (c *TypeCheckTest) TestQualifiedNameReduction7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {string} */ var f = function() {
var x = {/** @return {number} */ get a() {return 12}};
return x.a; }
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testQualifiedNameReduction7a()
func (c *TypeCheckTest) TestQualifiedNameReduction7a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {string} */ var f = function() {
var x = {get a() {return 12}};
return x.a; }
`).
		Run()
}

// Java: TypeCheckTest.testQualifiedNameReduction8()
func (c *TypeCheckTest) TestQualifiedNameReduction8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {string} */ var f = function() {
var x = {get a() {return 'a'}};
return x.a ? x.a : 'a'; }
`).
		Run()
}

// Java: TypeCheckTest.testQualifiedNameReduction9()
func (c *TypeCheckTest) TestQualifiedNameReduction9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {string} */ var f = function() {
var x = { /** @param {string} b */ set a(b) {}};
return x.a ? x.a : 'a'; }
`).
		Run()
}

// Java: TypeCheckTest.testQualifiedNameReduction10()
func (c *TypeCheckTest) TestQualifiedNameReduction10() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {string} */ var f = function() {
var x = { /** @param {number} b */ set a(b) {}};
return x.a ? x.a : 'a'; }
`).
		AddDiagnostic(`inconsistent return type
found   : (number|string)
required: string
`).
		Run()
}

// Java: TypeCheckTest.testUnknownsDontOverrideDeclaredTypesInLocalScope1()
func (c *TypeCheckTest) TestUnknownsDontOverrideDeclaredTypesInLocalScope1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ var C = function() {
  /** @type {string} */ this.a = 'str'};
/** @param {?} a
 @return {number} */
C.prototype.f = function(a) {
  this.a = a;
  return this.a;
}
`).
		AddDiagnostic(`inconsistent return type
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testUnknownsDontOverrideDeclaredTypesInLocalScope2()
func (c *TypeCheckTest) TestUnknownsDontOverrideDeclaredTypesInLocalScope2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ var C = function() {
  /** @type {string} */ this.a = 'str';
};
/** @type {C} */ var x = new C();
/** @param {?} a
 @return {number} */
C.prototype.f = function(a) {
  x.a = a;
  return x.a;
}
`).
		AddDiagnostic(`inconsistent return type
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testObjLitDef1a()
func (c *TypeCheckTest) TestObjLitDef1a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var x = {/** @type {number} */ a:12 };
x.a = 'a';
`).
		AddDiagnostic(`assignment to property a of x
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testObjLitDef1b()
func (c *TypeCheckTest) TestObjLitDef1b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(){
var x = {/** @type {number} */ a:12 };
x.a = 'a';
};
f();
`).
		AddDiagnostic(`assignment to property a of x
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testObjLitDef2a()
func (c *TypeCheckTest) TestObjLitDef2a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var x = {/** @param {number} b */ set a(b){} };
x.a = 'a';
`).
		AddDiagnostic(`assignment to property a of x
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testObjLitDef2b()
func (c *TypeCheckTest) TestObjLitDef2b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(){
var x = {/** @param {number} b */ set a(b){} };
x.a = 'a';
};
f();
`).
		AddDiagnostic(`assignment to property a of x
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testObjLitDef3a()
func (c *TypeCheckTest) TestObjLitDef3a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {string} */ var y;
var x = {/** @return {number} */ get a(){} };
y = x.a;
`).
		AddDiagnostic(`assignment
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testObjLitDef3b()
func (c *TypeCheckTest) TestObjLitDef3b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {string} */ var y;
function f(){
var x = {/** @return {number} */ get a(){} };
y = x.a;
};
f();
`).
		AddDiagnostic(`assignment
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testObjLitDef4()
func (c *TypeCheckTest) TestObjLitDef4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var x = {
  /** @return {number} */ a:12
};
`).
		AddDiagnostic(`assignment to property a of {a: function(): number}
found   : number
required: function(): number
`).
		Run()
}

// Java: TypeCheckTest.testMismatchInObjectKeyShouldShowBothTypes()
func (c *TypeCheckTest) TestMismatchInObjectKeyShouldShowBothTypes() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().IncludeDefaultExterns().
		AddSourceWithFilename("a.js", `goog.module('my.a');

/** @enum {string} */
const MyEnum = {
  KEY1: 'VALUE1',
  KEY2: 'VALUE2',
};

exports = {MyEnum};
`).
		AddSourceWithFilename("b.js", `goog.module('my.b');

/** @enum {string} */
const MyEnum = {
  KEY1: 'VALUE1',
  KEY2: 'VALUE2',
};

exports = {MyEnum};
`).
		AddSourceWithFilename("c.js", `goog.module('my.c');

const a = goog.require('my.a');
const b = goog.require('my.b');

/**
 * @typedef {!Object<!a.MyEnum, *>}
 */
let MyObjectTypedef;

/** @type {MyObjectTypedef} */
const myObject = {};

myObject[b.MyEnum.KEY1] = 'some value';
`).
		AddDiagnostic(`restricted index type
found   : string (enum definition: b.js:4:15)
required: string (enum definition: a.js:4:15)
`).
		Run()
}

// Java: TypeCheckTest.testObjLitDef5()
func (c *TypeCheckTest) TestObjLitDef5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var x = {};
/** @return {number} */ x.a = 12;
`).
		AddDiagnostic(`assignment to property a of x
found   : number
required: function(): number
`).
		Run()
}

// Java: TypeCheckTest.testObjLitDef6()
func (c *TypeCheckTest) TestObjLitDef6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var lit = /** @struct */ { 'x': 1 };").
		AddDiagnostic("Illegal key, the object literal is a struct").
		Run()
}

// Java: TypeCheckTest.testObjLitDef7()
func (c *TypeCheckTest) TestObjLitDef7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var lit = /** @dict */ { x: 1 };").
		AddDiagnostic("Illegal key, the object literal is a dict").
		Run()
}

// Java: TypeCheckTest.testInstanceOfReduction1()
func (c *TypeCheckTest) TestInstanceOfReduction1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ var T = function() {};
/** @param {T|string} x
@return {T} */
var f = function(x) {
if (x instanceof T) { return x; } else { return new T(); }
};
`).
		Run()
}

// Java: TypeCheckTest.testInstanceOfReduction2()
func (c *TypeCheckTest) TestInstanceOfReduction2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ var T = function() {};
/** @param {!T|string} x
@return {string} */
var f = function(x) {
if (x instanceof T) { return ''; } else { return x; }
};
`).
		Run()
}

// Java: TypeCheckTest.testUndeclaredGlobalProperty1()
func (c *TypeCheckTest) TestUndeclaredGlobalProperty1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const */ var x = {}; x.y = null;
function f(a) { x.y = a; }
/** @param {string} a */ function g(a) { }
function h() { g(x.y); }
`).
		Run()
}

// Java: TypeCheckTest.testUndeclaredGlobalProperty2()
func (c *TypeCheckTest) TestUndeclaredGlobalProperty2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const */ var x = {}; x.y = null;
function f() { x.y = 3; }
/** @param {string} a */ function g(a) { }
function h() { g(x.y); }
`).
		AddDiagnostic(`actual parameter 1 of g does not match formal parameter
found   : (null|number)
required: string
`).
		Run()
}

// Java: TypeCheckTest.testLocallyInferredGlobalProperty1()
func (c *TypeCheckTest) TestLocallyInferredGlobalProperty1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function F() {}
/** @type {number} */ F.prototype.z;
/** @const */ var x = {}; /** @type {F} */ x.y;
function f() { x.y.z = 'abc'; }
/** @param {number} x */ function g(x) {}
function h() { g(x.y.z); }
`).
		AddDiagnostic(`assignment to property z of F
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testPropertyInferredPropagation()
func (c *TypeCheckTest) TestPropertyInferredPropagation() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @return {Object} */
function f() { return {}; }
function g() { var x = f(); if (x.p) x.a = 'a'; else x.a = 'b'; }
function h() { var x = f(); x.a = false; }
`).
		Run()
}

// Java: TypeCheckTest.testPropertyInference1()
func (c *TypeCheckTest) TestPropertyInference1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function F() { this.x_ = true; }
/** @return {string} */
F.prototype.bar = function() { if (this.x_) return this.x_; };
`).
		AddDiagnostic(`inconsistent return type
found   : boolean
required: string
`).
		Run()
}

// Java: TypeCheckTest.testPropertyInference2()
func (c *TypeCheckTest) TestPropertyInference2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function F() { this.x_ = true; }
F.prototype.baz = function() { this.x_ = null; };
/** @return {string} */
F.prototype.bar = function() { if (this.x_) return this.x_; };
`).
		AddDiagnostic(`inconsistent return type
found   : boolean
required: string
`).
		Run()
}

// Java: TypeCheckTest.testPropertyInference3()
func (c *TypeCheckTest) TestPropertyInference3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function F() { this.x_ = true; }
F.prototype.baz = function() { this.x_ = 3; };
/** @return {string} */
F.prototype.bar = function() { if (this.x_) return this.x_; };
`).
		AddDiagnostic(`inconsistent return type
found   : (boolean|number)
required: string
`).
		Run()
}

// Java: TypeCheckTest.testPropertyInference4()
func (c *TypeCheckTest) TestPropertyInference4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function F() { }
F.prototype.x_ = 3;
/** @return {string} */
F.prototype.bar = function() { if (this.x_) return this.x_; };
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testPropertyInference5()
func (c *TypeCheckTest) TestPropertyInference5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function F() { }
F.prototype.baz = function() { this.x_ = 3; };
/** @return {string} */
F.prototype.bar = function() { if (this.x_) return this.x_; };
`).
		Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		Run()
}

// Java: TypeCheckTest.testPropertyInference6()
func (c *TypeCheckTest) TestPropertyInference6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function F() { }
(new F).x_ = 3;
/** @return {string} */
F.prototype.bar = function() { return this.x_; };
`).
		AddDiagnostic("Property x_ never defined on F").
		AddDiagnostic("Property x_ never defined on F").
		Run()
}

// Java: TypeCheckTest.testPropertyInference7()
func (c *TypeCheckTest) TestPropertyInference7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function F() { this.x_ = true; }
(new F).x_ = 3;
/** @return {string} */
F.prototype.bar = function() { return this.x_; };
`).
		AddDiagnostic(`inconsistent return type
found   : boolean
required: string
`).
		Run()
}

// Java: TypeCheckTest.testPropertyInference8()
func (c *TypeCheckTest) TestPropertyInference8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function F() {
  /** @type {string} */ this.x_ = 'x';
}
(new F).x_ = 3;
/** @return {string} */
F.prototype.bar = function() { return this.x_; };
`).
		AddDiagnostic(`assignment to property x_ of F
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testPropertyInference9()
func (c *TypeCheckTest) TestPropertyInference9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function A() {}
/** @return {function(): ?} */ function f() {
  return function() {};
}
var g = f();
/** @type {number} */ g.prototype.bar_ = null;
`).
		AddDiagnostic(`assignment
found   : null
required: number
`).
		Run()
}

// Java: TypeCheckTest.testPropertyInference10()
func (c *TypeCheckTest) TestPropertyInference10() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function A() {}
/** @return {function(): ?} */ function f() {
  return function() {};
}
var g = f();
/** @type {number} */ g.prototype.bar_ = 1;
var h = f();
/** @type {string} */ h.prototype.bar_ = 1;
`).
		AddDiagnostic(`assignment
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testNoPersistentTypeInferenceForObjectProperties()
func (c *TypeCheckTest) TestNoPersistentTypeInferenceForObjectProperties() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @param {Object} o
@param {string} x */
function s1(o,x) { o.x = x; }
/** @param {Object} o
@return {string} */
function g1(o) { return typeof o.x == 'undefined' ? '' : o.x; }
/** @param {Object} o
@param {number} x */
function s2(o,x) { o.x = x; }
/** @param {Object} o
@return {number} */
function g2(o) { return typeof o.x == 'undefined' ? 0 : o.x; }
`).
		Run()
}

// Java: TypeCheckTest.testNoPersistentTypeInferenceForFunctionProperties()
func (c *TypeCheckTest) TestNoPersistentTypeInferenceForFunctionProperties() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @param {Function} o
@param {string} x */
function s1(o,x) { o.x = x; }
/** @param {Function} o
@return {string} */
function g1(o) { return typeof o.x == 'undefined' ? '' : o.x; }
/** @param {Function} o
@param {number} x */
function s2(o,x) { o.x = x; }
/** @param {Function} o
@return {number} */
function g2(o) { return typeof o.x == 'undefined' ? 0 : o.x; }
`).
		Run()
}

// Java: TypeCheckTest.testStrictPropertiesOnFunctions2()
func (c *TypeCheckTest) TestStrictPropertiesOnFunctions2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {Function} o
 @param {string} x */
function s1(o,x) { o.x = x; }
`).
		AddDiagnostic("Property x never defined on Function").
		Run()
}

// Java: TypeCheckTest.testStrictPropertiesOnEnum1()
func (c *TypeCheckTest) TestStrictPropertiesOnEnum1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddString().Build()).
		AddSource(`/** @enum {string} */ var E = {S:'s'};
/** @param {E} e
 @return {string}  */
function s1(e) { return e.slice(1); }
`).
		Run()
}

// Java: TypeCheckTest.testStrictPropertiesOnEnum2()
func (c *TypeCheckTest) TestStrictPropertiesOnEnum2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddString().Build()).
		AddSource(`/** @enum {?string} */ var E = {S:'s'};
/** @param {E} e
 @return {string}  */
function s1(e) { return e.slice(1); }
`).
		Run()
}

// Java: TypeCheckTest.testStrictPropertiesOnEnum3()
func (c *TypeCheckTest) TestStrictPropertiesOnEnum3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddString().Build()).
		AddSource(`/** @enum {string} */ var E = {S:'s'};
/** @param {?E} e
 @return {string}  */
function s1(e) { return e.slice(1); }
`).
		Run()
}

// Java: TypeCheckTest.testStrictPropertiesOnEnum4()
func (c *TypeCheckTest) TestStrictPropertiesOnEnum4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddString().Build()).
		AddSource(`/** @enum {?string} */ var E = {S:'s'};
/** @param {?E} e
 @return {string}  */
function s1(e) { return e.slice(1); }
`).
		Run()
}

// Java: TypeCheckTest.testDestructuring()
func (c *TypeCheckTest) TestDestructuring() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {{id:(number|undefined)}=} o */
function f(o = {}) {
   const {id} = o;
}
`).
		Run()
}

// Java: TypeCheckTest.testObjectPropertyTypeInferredInLocalScope1a()
func (c *TypeCheckTest) TestObjectPropertyTypeInferredInLocalScope1a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @param {!Object} o
 * @return {string}
 */
function f(o) {
  o.x = 1;
  return o.x;
}
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testObjectPropertyTypeInferredInLocalScope1b()
func (c *TypeCheckTest) TestObjectPropertyTypeInferredInLocalScope1b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {!Object} o
 * @return {string}
 */
function f(o) {
  o.x = 1;
}
`).
		AddDiagnostic("Property x never defined on Object").
		Run()
}

// Java: TypeCheckTest.testObjectPropertyTypeInferredInLocalScope2a()
func (c *TypeCheckTest) TestObjectPropertyTypeInferredInLocalScope2a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/**
  @param {!Object} o
  @param {number?} x
  @return {string}
 */
function f(o, x) {
  o.x = 'a';
  if (x) {o.x = x;}
  return o.x;
}
`).
		AddDiagnostic(`inconsistent return type
found   : (number|string)
required: string
`).
		Run()
}

// Java: TypeCheckTest.testObjectPropertyTypeInferredInLocalScope2b()
func (c *TypeCheckTest) TestObjectPropertyTypeInferredInLocalScope2b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
  @param {!Object} o
  @param {number?} x
 */
function f(o, x) {
  o.x = 'a';
}
`).
		AddDiagnostic("Property x never defined on Object").
		Run()
}

// Java: TypeCheckTest.testObjectPropertyTypeInferredInLocalScope2c()
func (c *TypeCheckTest) TestObjectPropertyTypeInferredInLocalScope2c() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
  @param {!Object} o
  @param {number?} x
 */
function f(o, x) {
  if (x) { o.x = x; }
}
`).
		AddDiagnostic("Property x never defined on Object").
		Run()
}

// Java: TypeCheckTest.testObjectPropertyTypeInferredInLocalScope3()
func (c *TypeCheckTest) TestObjectPropertyTypeInferredInLocalScope3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {!Object} o
 * @param {number?} x
 * @return {string}
 */
function f(o, x) {
  if (x) {
    o.x = x;
  } else {
    o.x = 'a';
  }
  return o.x;
}
`).
		AddDiagnostic(`inconsistent return type
found   : (number|string)
required: string
`).
		Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		Run()
}

// Java: TypeCheckTest.testMissingPropertyOfUnion1()
func (c *TypeCheckTest) TestMissingPropertyOfUnion1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
  @param {number|string} obj
  @return {?}
 */
function f(obj, x) {
  return obj.foo;
}
`).
		AddDiagnostic("Property foo never defined on (Number|String)").
		Run()
}

// Java: TypeCheckTest.testMissingPropertyOfUnion2()
func (c *TypeCheckTest) TestMissingPropertyOfUnion2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
  @param {*=} opt_o
  @return {string|null}
 */
function f(opt_o) {
  if (opt_o) {
    if (typeof opt_o !== 'string') {
      return opt_o.toString();
    }
  }
  return null;
}
`).
		AddDiagnostic("Property toString never defined on *" + TypeCheck_POSSIBLE_INEXISTENT_PROPERTY_EXPLANATION).
		Run()
}

// Java: TypeCheckTest.testMismatchingOverridingInferredPropertyBeforeDeclaredProperty1()
func (c *TypeCheckTest) TestMismatchingOverridingInferredPropertyBeforeDeclaredProperty1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */var T = function() { this.x = ''; };
/** @type {number} */ T.prototype.x = 0;
`).
		AddDiagnostic(`assignment to property x of T
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testMismatchingOverridingInferredPropertyBeforeDeclaredProperty2()
func (c *TypeCheckTest) TestMismatchingOverridingInferredPropertyBeforeDeclaredProperty2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */var T = function() { this.x = ''; };
/** @type {number} */ T.prototype.x;
`).
		AddDiagnostic(`assignment to property x of T
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testMismatchingOverridingInferredPropertyBeforeDeclaredProperty3()
func (c *TypeCheckTest) TestMismatchingOverridingInferredPropertyBeforeDeclaredProperty3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @type {Object} */ var n = {};
/** @constructor */ n.T = function() { this.x = ''; };
/** @type {number} */ n.T.prototype.x = 0;
`).
		AddDiagnostic(`assignment to property x of n.T
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testMismatchingOverridingInferredPropertyBeforeDeclaredProperty4()
func (c *TypeCheckTest) TestMismatchingOverridingInferredPropertyBeforeDeclaredProperty4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var n = {};
/** @constructor */ n.T = function() { this.x = ''; };
/** @type {number} */ n.T.prototype.x = 0;
`).
		AddDiagnostic(`assignment to property x of n.T
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testPropertyUsedBeforeDefinition1()
func (c *TypeCheckTest) TestPropertyUsedBeforeDefinition1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ var T = function() {};
/** @return {string} */
T.prototype.f = function() { return this.g(); };
/** @return {number} */ T.prototype.g = function() { return 1; };
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testPropertyUsedBeforeDefinition2()
func (c *TypeCheckTest) TestPropertyUsedBeforeDefinition2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var n = {};
/** @constructor */ n.T = function() {};
/** @return {string} */
n.T.prototype.f = function() { return this.g(); };
/** @return {number} */ n.T.prototype.g = function() { return 1; };
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testAdd1()
func (c *TypeCheckTest) TestAdd1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@return {void}*/function foo(){var a = 'abc'+foo();}").
		Run()
}

// Java: TypeCheckTest.testAdd2()
func (c *TypeCheckTest) TestAdd2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@return {void}*/function foo(){var a = foo()+4;}").
		Run()
}

// Java: TypeCheckTest.testAdd3()
func (c *TypeCheckTest) TestAdd3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {string} */ var a = 'a';
/** @type {string} */ var b = 'b';
/** @type {string} */ var c = a + b;
`).
		Run()
}

// Java: TypeCheckTest.testAdd4()
func (c *TypeCheckTest) TestAdd4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {number} */ var a = 5;
/** @type {string} */ var b = 'b';
/** @type {string} */ var c = a + b;
`).
		Run()
}

// Java: TypeCheckTest.testAdd5()
func (c *TypeCheckTest) TestAdd5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {string} */ var a = 'a';
/** @type {number} */ var b = 5;
/** @type {string} */ var c = a + b;
`).
		Run()
}

// Java: TypeCheckTest.testAdd6()
func (c *TypeCheckTest) TestAdd6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {number} */ var a = 5;
/** @type {number} */ var b = 5;
/** @type {number} */ var c = a + b;
`).
		Run()
}

// Java: TypeCheckTest.testAdd7()
func (c *TypeCheckTest) TestAdd7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {number} */ var a = 5;
/** @type {string} */ var b = 'b';
/** @type {number} */ var c = a + b;
`).
		AddDiagnostic(`initializing variable
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testAdd8()
func (c *TypeCheckTest) TestAdd8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {string} */ var a = 'a';
/** @type {number} */ var b = 5;
/** @type {number} */ var c = a + b;
`).
		AddDiagnostic(`initializing variable
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testAdd9()
func (c *TypeCheckTest) TestAdd9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {number} */ var a = 5;
/** @type {number} */ var b = 5;
/** @type {string} */ var c = a + b;
`).
		AddDiagnostic(`initializing variable
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testAdd10()
func (c *TypeCheckTest) TestAdd10() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(c.SuppressMissingProperty("e", "f") + `/** @type {number} */ var a = 5;
/** @type {string} */ var c = a + d.e.f;
`).
		Run()
}

// Java: TypeCheckTest.testAdd11()
func (c *TypeCheckTest) TestAdd11() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(c.SuppressMissingProperty("e", "f") + `/** @type {number} */ var a = 5;
/** @type {number} */ var c = a + d.e.f;
`).
		Run()
}

// Java: TypeCheckTest.testAdd12()
func (c *TypeCheckTest) TestAdd12() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {(number|string)} */ function a() { return 5; }
/** @type {number} */ var b = 5;
/** @type {boolean} */ var c = a() + b;
`).
		AddDiagnostic(`initializing variable
found   : (number|string)
required: boolean
`).
		Run()
}

// Java: TypeCheckTest.testAdd13()
func (c *TypeCheckTest) TestAdd13() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {number} */ var a = 5;
/** @return {(number|string)} */ function b() { return 5; }
/** @type {boolean} */ var c = a + b();
`).
		AddDiagnostic(`initializing variable
found   : (number|string)
required: boolean
`).
		Run()
}

// Java: TypeCheckTest.testAdd14()
func (c *TypeCheckTest) TestAdd14() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {(null|string)} */ var a = unknown;
/** @type {number} */ var b = 5;
/** @type {boolean} */ var c = a + b;
`).
		AddDiagnostic(`initializing variable
found   : (number|string)
required: boolean
`).
		Run()
}

// Java: TypeCheckTest.testAdd15()
func (c *TypeCheckTest) TestAdd15() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {number} */ var a = 5;
/** @return {(number|string)} */ function b() { return 5; }
/** @type {boolean} */ var c = a + b();
`).
		AddDiagnostic(`initializing variable
found   : (number|string)
required: boolean
`).
		Run()
}

// Java: TypeCheckTest.testAdd16()
func (c *TypeCheckTest) TestAdd16() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {(undefined|string)} */ var a = unknown;
/** @type {number} */ var b = 5;
/** @type {boolean} */ var c = a + b;
`).
		AddDiagnostic(`initializing variable
found   : (number|string)
required: boolean
`).
		Run()
}

// Java: TypeCheckTest.testAdd17()
func (c *TypeCheckTest) TestAdd17() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {number} */ var a = 5;
/** @type {(undefined|string)} */ var b = unknown;
/** @type {boolean} */ var c = a + b;
`).
		AddDiagnostic(`initializing variable
found   : (number|string)
required: boolean
`).
		Run()
}

// Java: TypeCheckTest.testAdd18()
func (c *TypeCheckTest) TestAdd18() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f() {};
/** @type {string} */ var a = 'a';
/** @type {number} */ var c = a + f();
`).
		AddDiagnostic(`initializing variable
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testAdd19()
func (c *TypeCheckTest) TestAdd19() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {number} opt_x
@param {number} opt_y
@return {number} */ function f(opt_x, opt_y) {
return opt_x + opt_y;}
`).
		Run()
}

// Java: TypeCheckTest.testAdd20()
func (c *TypeCheckTest) TestAdd20() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {!Number} opt_x
@param {!Number} opt_y
@return {number} */ function f(opt_x, opt_y) {
return opt_x + opt_y;}
`).
		Run()
}

// Java: TypeCheckTest.testAdd21()
func (c *TypeCheckTest) TestAdd21() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {Number|Boolean} opt_x
@param {number|boolean} opt_y
@return {number} */ function f(opt_x, opt_y) {
return opt_x + opt_y;}
`).
		Run()
}

// Java: TypeCheckTest.testNumericComparison1()
func (c *TypeCheckTest) TestNumericComparison1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@param {number} a*/ function f(a) {return a < 3;}").
		Run()
}

// Java: TypeCheckTest.testNumericComparison2()
func (c *TypeCheckTest) TestNumericComparison2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@param {!Object} a*/ function f(a) {return a < 3;}").
		AddDiagnostic(`left side of numeric comparison
found   : Object
required: (bigint|number)
`).
		Run()
}

// Java: TypeCheckTest.testNumericComparison3()
func (c *TypeCheckTest) TestNumericComparison3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_PRIMITIVE_OPERATORS).
		AddSource("/**@param {string} a*/ function f(a) {return a < 3;}").
		Run()
}

// Java: TypeCheckTest.testNumericComparison4()
func (c *TypeCheckTest) TestNumericComparison4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_PRIMITIVE_OPERATORS).
		AddSource(`/**@param {(number|undefined)} a*/
function f(a) {return a < 3;}
`).
		Run()
}

// Java: TypeCheckTest.testNumericComparison5()
func (c *TypeCheckTest) TestNumericComparison5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@param {*} a*/ function f(a) {return a < 3;}").
		AddDiagnostic(`left side of numeric comparison
found   : *
required: (bigint|number)
`).
		Run()
}

// Java: TypeCheckTest.testNumericComparison6()
func (c *TypeCheckTest) TestNumericComparison6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@return {void} */ function foo() { if (3 >= foo()) return; }").
		AddDiagnostic(`right side of numeric comparison
found   : undefined
required: (bigint|number)
`).
		Run()
}

// Java: TypeCheckTest.testStringComparison1()
func (c *TypeCheckTest) TestStringComparison1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@param {string} a*/ function f(a) {return a < 'x';}").
		Run()
}

// Java: TypeCheckTest.testStringComparison2()
func (c *TypeCheckTest) TestStringComparison2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_PRIMITIVE_OPERATORS).
		AddSource("/**@param {Object} a*/ function f(a) {return a < 'x';}").
		Run()
}

// Java: TypeCheckTest.testStringComparison3()
func (c *TypeCheckTest) TestStringComparison3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_PRIMITIVE_OPERATORS).
		AddSource("/**@param {number} a*/ function f(a) {return a < 'x';}").
		Run()
}

// Java: TypeCheckTest.testStringComparison4()
func (c *TypeCheckTest) TestStringComparison4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**@param {string|undefined} a*/
function f(a) {return a < 'x';}
`).
		Run()
}

// Java: TypeCheckTest.testStringComparison5()
func (c *TypeCheckTest) TestStringComparison5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**@param {*} a*/
function f(a) {return a < 'x';}
`).
		Run()
}

// Java: TypeCheckTest.testStringComparison6()
func (c *TypeCheckTest) TestStringComparison6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_PRIMITIVE_OPERATORS).
		AddSource(`/**@return {void} */
function foo() { if ('a' >= foo()) return; }
`).
		AddDiagnostic(`right side of comparison
found   : undefined
required: string
`).
		Run()
}

// Java: TypeCheckTest.testValueOfComparison1()
func (c *TypeCheckTest) TestValueOfComparison1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddObject().Build()).
		AddSource(`/** @constructor */function O() {};
/**@override*/O.prototype.valueOf = function() { return 1; };
/**@param {!O} a
@param {!O} b*/ function f(a,b) { return a < b; }
`).
		Run()
}

// Java: TypeCheckTest.testValueOfComparison2()
func (c *TypeCheckTest) TestValueOfComparison2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_PRIMITIVE_OPERATORS).
		AddExterns(NewTestExternsBuilder().AddObject().Build()).
		AddSource(`/** @constructor */function O() {};
/**@override*/O.prototype.valueOf = function() { return 1; };
/**
 * @param {!O} a
 * @param {number} b
 */
function f(a,b) { return a < b; }
`).
		Run()
}

// Java: TypeCheckTest.testValueOfComparison3()
func (c *TypeCheckTest) TestValueOfComparison3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_PRIMITIVE_OPERATORS).
		AddExterns(NewTestExternsBuilder().AddObject().Build()).
		AddSource(`/** @constructor */function O() {};
/**@override*/O.prototype.toString = function() { return 'o'; };
/**
 * @param {!O} a
 * @param {string} b
 */
function f(a,b) { return a < b; }
`).
		Run()
}

// Java: TypeCheckTest.testGenericRelationalExpression()
func (c *TypeCheckTest) TestGenericRelationalExpression() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**@param {*} a
@param {*} b*/
function f(a,b) {return a < b;}
`).
		Run()
}

// Java: TypeCheckTest.testInstanceof1()
func (c *TypeCheckTest) TestInstanceof1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function foo(){
if (bar instanceof 3)return;}
`).
		AddDiagnostic(`instanceof requires an object
found   : number
required: Object
`).
		Run()
}

// Java: TypeCheckTest.testInstanceof2()
func (c *TypeCheckTest) TestInstanceof2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**@return {void}*/function foo(){
if (foo() instanceof Object)return;}
`).
		AddDiagnostic(`deterministic instanceof yields false
found   : undefined
required: NoObject
`).
		Run()
}

// Java: TypeCheckTest.testInstanceof3()
func (c *TypeCheckTest) TestInstanceof3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**@return {*} */function foo(){
if (foo() instanceof Object)return;}
`).
		Run()
}

// Java: TypeCheckTest.testInstanceof4()
func (c *TypeCheckTest) TestInstanceof4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**@return {(Object|number)} */function foo(){
if (foo() instanceof Object)return 3;}
`).
		Run()
}

// Java: TypeCheckTest.testInstanceof5()
func (c *TypeCheckTest) TestInstanceof5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {?} */ function foo(){
if (foo() instanceof Object)return;}
`).
		Run()
}

// Java: TypeCheckTest.testInstanceof6()
func (c *TypeCheckTest) TestInstanceof6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**@return {(Array|number)} */function foo(){
if (foo() instanceof Object)return 3;}
`).
		Run()
}

// Java: TypeCheckTest.testInstanceOfReduction3()
func (c *TypeCheckTest) TestInstanceOfReduction3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** 
 * @param {Object} x 
 * @param {Function} y 
 * @return {boolean} 
 */
var f = function(x, y) {
  return x instanceof y;
};
`).
		Run()
}

// Java: TypeCheckTest.testScoping1()
func (c *TypeCheckTest) TestScoping1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**@param {string} a*/function foo(a){
  /**@param {Array|string} a*/function bar(a){
    if (a instanceof Array)return;
  }
}
`).
		Run()
}

// Java: TypeCheckTest.testScoping2()
func (c *TypeCheckTest) TestScoping2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {number} */ var a;
function Foo() {
  /** @type {string} */ var a;
}
`).
		Run()
}

// Java: TypeCheckTest.testScoping3()
func (c *TypeCheckTest) TestScoping3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {Number} */ var b;
/** @type {!String} */ var b;
`).
		AddDiagnostic(`variable b redefined with type String, original definition at testcode0:1 with type (Number|null)
`).
		Run()
}

// Java: TypeCheckTest.testScoping4()
func (c *TypeCheckTest) TestScoping4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type{Number}*/var b; if (true) /** @type{!String} */var b;").
		AddDiagnostic(`variable b redefined with type String, original definition at testcode0:1 with type (Number|null)
`).
		Run()
}

// Java: TypeCheckTest.testScoping5()
func (c *TypeCheckTest) TestScoping5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("if (true) var b; var b;").
		Run()
}

// Java: TypeCheckTest.testScoping6()
func (c *TypeCheckTest) TestScoping6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("if (true) var b; if (true) var b;").
		Run()
}

// Java: TypeCheckTest.testScoping7()
func (c *TypeCheckTest) TestScoping7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */function A() {
  /** @type {!A} */this.a = null;
}
`).
		AddDiagnostic(`assignment to property a of A
found   : null
required: A
`).
		Run()
}

// Java: TypeCheckTest.testScoping8()
func (c *TypeCheckTest) TestScoping8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */function A() {}
/** @constructor */function B() {
  /** @type {!A} */this.a = null;
}
`).
		AddDiagnostic(`assignment to property a of B
found   : null
required: A
`).
		Run()
}

// Java: TypeCheckTest.testScoping9()
func (c *TypeCheckTest) TestScoping9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */function B() {
  /** @type {!A} */this.a = null;
}
/** @constructor */function A() {}
`).
		AddDiagnostic(`assignment to property a of B
found   : null
required: A
`).
		Run()
}

// Java: TypeCheckTest.testScoping10()
func (c *TypeCheckTest) TestScoping10() {
	var p *TypeCheckTestCase_TypeCheckResult = c.ParseAndTypeCheckWithScope("var a = function b(){};")
	ScopeSubject_AssertScope(p.scope).Declares("a")
	ScopeSubject_AssertScope(p.scope).DoesNotDeclare("b")
	AssertThat(p.scope.GetVar("a").GetType().
		String()).IsEqualTo("function(): undefined")
}

// Java: TypeCheckTest.testDontDropPropertiesInUnion1()
func (c *TypeCheckTest) TestDontDropPropertiesInUnion1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {{a: number}|{a:number, b:string}} x */
function f(x) {
  var /** null */ n = x.b;
}
`).
		AddDiagnosticType(TypeCheck_STRICT_INEXISTENT_UNION_PROPERTY).
		AddDiagnosticType(TypeValidator_TYPE_MISMATCH_WARNING).
		Run()
}

// Java: TypeCheckTest.testDontDropPropertiesInUnion2()
func (c *TypeCheckTest) TestDontDropPropertiesInUnion2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {{a:number, b:string}|{a: number}} x */
function f(x) {
  var /** null */ n = x.b;
}
`).
		AddDiagnosticType(TypeCheck_STRICT_INEXISTENT_UNION_PROPERTY).
		AddDiagnosticType(TypeValidator_TYPE_MISMATCH_WARNING).
		Run()
}

// Java: TypeCheckTest.testDontDropPropertiesInUnion3()
func (c *TypeCheckTest) TestDontDropPropertiesInUnion3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {{a: number}|{a:number, b:string}} x */
function f(x) {}
/** @param {{a: number}} x */
function g(x) { return x.b; }
`).
		AddDiagnosticType(TypeCheck_INEXISTENT_PROPERTY).
		Run()
}

// Java: TypeCheckTest.testDontDropPropertiesInUnion4()
func (c *TypeCheckTest) TestDontDropPropertiesInUnion4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {{a: number}|{a:number, b:string}} x */
function f(x) {}
/** @param {{c: number}} x */
function g(x) { return x.b; }
`).
		AddDiagnostic("Property b never defined on x").
		Run()
}

// Java: TypeCheckTest.testDontDropPropertiesInUnion5()
func (c *TypeCheckTest) TestDontDropPropertiesInUnion5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {{a: number}|{a: number, b: string}} x */
function f(x) {}
f({a: 123});
`).
		Run()
}

// Java: TypeCheckTest.testDontDropPropertiesInUnion6()
func (c *TypeCheckTest) TestDontDropPropertiesInUnion6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {{a: number}|{a: number, b: string}} x */
function f(x) {
  var /** null */ n = x;
}
`).
		AddDiagnostic(`initializing variable
found   : ({
  a: number,
  b: string
}|{a: number})
required: null
`).
		Run()
}

// Java: TypeCheckTest.testDontDropPropertiesInUnion7()
func (c *TypeCheckTest) TestDontDropPropertiesInUnion7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {{a: number}|{a:number, b:string}} x */
function f(x) {}
/** @param {{c: number}|{c:number, d:string}} x */
function g(x) { return x.b; }
`).
		AddDiagnosticType(TypeCheck_INEXISTENT_PROPERTY).
		Run()
}

// Java: TypeCheckTest.testScoping11()
func (c *TypeCheckTest) TestScoping11() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @return {number} */var a = function b(){ return b };").
		AddDiagnostic(`inconsistent return type
found   : function(): number
required: number
`).
		Run()
}

// Java: TypeCheckTest.testScoping12()
func (c *TypeCheckTest) TestScoping12() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function F() {}
/** @type {number} */ F.prototype.bar = 3;
/** @param {!F} f */ function g(f) {
  /** @return {string} */
  function h() {
    return f.bar;
  }
}
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testFunctionArguments1()
func (c *TypeCheckTest) TestFunctionArguments1() {
	c.TestFunctionType(`/** @param {number} a
@return {string} */
function f(a) {}
`, "function(number): string")
}

// Java: TypeCheckTest.testFunctionArguments2()
func (c *TypeCheckTest) TestFunctionArguments2() {
	c.TestFunctionType(`/** @param {number} opt_a
@return {string} */
function f(opt_a) {}
`, "function(number=): string")
}

// Java: TypeCheckTest.testFunctionArguments3()
func (c *TypeCheckTest) TestFunctionArguments3() {
	c.TestFunctionType(`/** @param {number} b
@return {string} */
function f(a,b) {}
`, "function(?, number): string")
}

// Java: TypeCheckTest.testFunctionArguments4()
func (c *TypeCheckTest) TestFunctionArguments4() {
	c.TestFunctionType(`/** @param {number} opt_a
@return {string} */
function f(a,opt_a) {}
`, "function(?, number=): string")
}

// Java: TypeCheckTest.testFunctionArguments5()
func (c *TypeCheckTest) TestFunctionArguments5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function a(opt_a,a) {}").
		AddDiagnostic("optional arguments must be at the end").
		Run()
}

// Java: TypeCheckTest.testFunctionArguments6()
func (c *TypeCheckTest) TestFunctionArguments6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function a(var_args,a) {}").
		AddDiagnostic("variable length argument must be last").
		Run()
}

// Java: TypeCheckTest.testFunctionArguments7()
func (c *TypeCheckTest) TestFunctionArguments7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {number} opt_a
@return {string} */
function a(a,opt_a,var_args) {}
`).
		Run()
}

// Java: TypeCheckTest.testFunctionArguments8()
func (c *TypeCheckTest) TestFunctionArguments8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function a(a,opt_a,var_args,b) {}").
		AddDiagnostic("variable length argument must be last").
		Run()
}

// Java: TypeCheckTest.testFunctionArguments9()
func (c *TypeCheckTest) TestFunctionArguments9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function a(a,opt_a,var_args,b,c) {}").
		AddDiagnostic("variable length argument must be last").
		Run()
}

// Java: TypeCheckTest.testFunctionArguments10()
func (c *TypeCheckTest) TestFunctionArguments10() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function a(a,opt_a,b,c) {}").
		AddDiagnostic("optional arguments must be at the end").
		Run()
}

// Java: TypeCheckTest.testFunctionArguments11()
func (c *TypeCheckTest) TestFunctionArguments11() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function a(a,opt_a,b,c,var_args,d) {}").
		AddDiagnostic("optional arguments must be at the end").
		Run()
}

// Java: TypeCheckTest.testFunctionArguments12()
func (c *TypeCheckTest) TestFunctionArguments12() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @param {String} foo  */function bar(baz){}").
		AddDiagnostic("parameter foo does not appear in bar's parameter list").
		Run()
}

// Java: TypeCheckTest.testFunctionArguments13()
func (c *TypeCheckTest) TestFunctionArguments13() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {boolean} */ function u() { return true; }
/** @param {boolean} b
@return {?boolean} */
function f(b) { if (u()) { b = null; } return b; }
`).
		AddDiagnostic(`assignment
found   : null
required: boolean
`).
		Run()
}

// Java: TypeCheckTest.testFunctionArguments14()
func (c *TypeCheckTest) TestFunctionArguments14() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {string} x
 * @param {number} opt_y
 * @param {boolean} var_args
 */ function f(x, opt_y, var_args) {}
f('3'); f('3', 2); f('3', 2, true); f('3', 2, true, false);
`).
		Run()
}

// Java: TypeCheckTest.testFunctionArguments15()
func (c *TypeCheckTest) TestFunctionArguments15() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {?function(*)} f */
function g(f) { f(1, 2); }
`).
		AddDiagnostic(`Function f: called with 2 argument(s). Function requires at least 1 argument(s) and no more than 1 argument(s).
`).
		Run()
}

// Java: TypeCheckTest.testFunctionArguments16()
func (c *TypeCheckTest) TestFunctionArguments16() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {...number} var_args */
function g(var_args) {}
g(1, true);
`).
		AddDiagnostic(`actual parameter 2 of g does not match formal parameter
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckTest.testUndefinedPassedForVarArgs()
func (c *TypeCheckTest) TestUndefinedPassedForVarArgs() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {...number} var_args */
function g(var_args) {}
g(undefined, 1);
`).
		AddDiagnostic(`actual parameter 1 of g does not match formal parameter
found   : undefined
required: number
`).
		Run()
}

// Java: TypeCheckTest.testFunctionArguments17()
func (c *TypeCheckTest) TestFunctionArguments17() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {booool|string} x */
function f(x) { g(x) }
/** @param {number} x */
function g(x) {}
`).
		AddDiagnostic("Bad type annotation. Unknown type booool").
		Run()
}

// Java: TypeCheckTest.testFunctionArguments18()
func (c *TypeCheckTest) TestFunctionArguments18() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(x) {}
f(/** @param {number} y */ (function() {}));
`).
		AddDiagnostic("parameter y does not appear in <anonymous>'s parameter list").
		Run()
}

// Java: TypeCheckTest.testVarArgParameterWithTemplateType()
func (c *TypeCheckTest) TestVarArgParameterWithTemplateType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddArray().Build()).
		AddSource(`/**
 * @template T
 * @param {...T} var_args
 * @return {T}
 */
function firstOf(var_args) { return arguments[0]; }
/** @type {null} */ var a = firstOf('hi', 1);
`).
		AddDiagnostic(`initializing variable
found   : (number|string)
required: null
`).
		Run()
}

// Java: TypeCheckTest.testRestParameters()
func (c *TypeCheckTest) TestRestParameters() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddArray().Build()).
		AddSource(`/**
 * @param {...number} x
 */
function f(...x) {
  var /** string */ s = x[0];
}
`).
		AddDiagnostic(`initializing variable
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testRestParameterWithTemplateType()
func (c *TypeCheckTest) TestRestParameterWithTemplateType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddArray().Build()).
		AddSource(`/**
 * @template T
 * @param {...T} rest
 * @return {T}
 */
function firstOf(...rest) {
  return rest[0];
}
/** @type {null} */ var a = firstOf('hi', 1);
`).
		AddDiagnostic(`initializing variable
found   : (number|string)
required: null
`).
		Run()
}

// Java: TypeCheckTest.testRestParametersWithGenericsNoWarning()
func (c *TypeCheckTest) TestRestParametersWithGenericsNoWarning() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddObject().AddArray().
		Build()).
		AddSource(`/**
 * @constructor
 * @template T
 */
function Foo() {}
/**
 * @template T
 * @param {...!Foo<T>} x
 */
function f(...x) {
  return 123;
}
`).
		Run()
}

// Java: TypeCheckTest.testPrintFunctionName1()
func (c *TypeCheckTest) TestPrintFunctionName1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = {}; goog.run = function(f) {};
goog.run();
`).
		AddDiagnostic(`Function goog.run: called with 0 argument(s). Function requires at least 1 argument(s) and no more than 1 argument(s).
`).
		Run()
}

// Java: TypeCheckTest.testPrintFunctionName2()
func (c *TypeCheckTest) TestPrintFunctionName2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ var Foo = function() {};
Foo.prototype.run = function(f) {};
(new Foo).run();
`).
		AddDiagnostic(`Function Foo.prototype.run: called with 0 argument(s). Function requires at least 1 argument(s) and no more than 1 argument(s).
`).
		Run()
}

// Java: TypeCheckTest.testFunctionInference1()
func (c *TypeCheckTest) TestFunctionInference1() {
	c.TestFunctionType("function f(a) {}", "function(?): undefined")
}

// Java: TypeCheckTest.testFunctionInference2()
func (c *TypeCheckTest) TestFunctionInference2() {
	c.TestFunctionType("function f(a,b) {}", "function(?, ?): undefined")
}

// Java: TypeCheckTest.testFunctionInference3()
func (c *TypeCheckTest) TestFunctionInference3() {
	c.TestFunctionType("function f(var_args) {}", "function(...?): undefined")
}

// Java: TypeCheckTest.testFunctionInference4()
func (c *TypeCheckTest) TestFunctionInference4() {
	c.TestFunctionType("function f(a,b,c,var_args) {}", "function(?, ?, ?, ...?): undefined")
}

// Java: TypeCheckTest.testFunctionInference5()
func (c *TypeCheckTest) TestFunctionInference5() {
	c.TestFunctionType(`/** @this {Date}
@return {string} */function f(a) {}
`, "function(this:Date, ?): string")
}

// Java: TypeCheckTest.testFunctionInference6()
func (c *TypeCheckTest) TestFunctionInference6() {
	c.TestFunctionType(`/** @this {Date}
@return {string} */function f(opt_a) {}
`, "function(this:Date, ?=): string")
}

// Java: TypeCheckTest.testFunctionInference7()
func (c *TypeCheckTest) TestFunctionInference7() {
	c.TestFunctionType("/** @this {Date} */function f(a,b,c,var_args) {}", "function(this:Date, ?, ?, ?, ...?): undefined")
}

// Java: TypeCheckTest.testFunctionInference8()
func (c *TypeCheckTest) TestFunctionInference8() {
	c.TestFunctionType("function f() {}", "function(): undefined")
}

// Java: TypeCheckTest.testFunctionInference9()
func (c *TypeCheckTest) TestFunctionInference9() {
	c.TestFunctionType("var f = function() {};", "function(): undefined")
}

// Java: TypeCheckTest.testFunctionInference10()
func (c *TypeCheckTest) TestFunctionInference10() {
	c.TestFunctionType(`/** @this {Date}
@param {boolean} b
@return {string} */var f = function(a,b) {};
`, "function(this:Date, ?, boolean): string")
}

// Java: TypeCheckTest.testFunctionInference11()
func (c *TypeCheckTest) TestFunctionInference11() {
	c.TestFunctionTypeWithFunctionType(`var goog = {};
/** @return {number}*/goog.f = function(){};
`, "goog.f", "function(): number")
}

// Java: TypeCheckTest.testFunctionInference12()
func (c *TypeCheckTest) TestFunctionInference12() {
	c.TestFunctionTypeWithFunctionType(`var goog = {};
goog.f = function(){};
`, "goog.f", "function(): undefined")
}

// Java: TypeCheckTest.testFunctionInference13()
func (c *TypeCheckTest) TestFunctionInference13() {
	c.TestFunctionTypeWithFunctionType(`var goog = {};
/** @constructor */ goog.Foo = function(){};
/** @param {!goog.Foo} f */function eatFoo(f){};
`, "eatFoo", "function(goog.Foo): undefined")
}

// Java: TypeCheckTest.testFunctionInference14()
func (c *TypeCheckTest) TestFunctionInference14() {
	c.TestFunctionTypeWithFunctionType(`var goog = {};
/** @constructor */ goog.Foo = function(){};
/** @return {!goog.Foo} */function eatFoo(){ return new goog.Foo; };
`, "eatFoo", "function(): goog.Foo")
}

// Java: TypeCheckTest.testFunctionInference15()
func (c *TypeCheckTest) TestFunctionInference15() {
	c.TestFunctionTypeWithFunctionType(`/** @constructor */ function f() {};
f.prototype.foo = function(){};
`, "f.prototype.foo", "function(this:f): undefined")
}

// Java: TypeCheckTest.testFunctionInference16()
func (c *TypeCheckTest) TestFunctionInference16() {
	c.TestFunctionTypeWithFunctionType(`/** @constructor */ function f() {};
f.prototype.foo = function(){};
`, "(new f).foo", "function(this:f): undefined")
}

// Java: TypeCheckTest.testFunctionInference17()
func (c *TypeCheckTest) TestFunctionInference17() {
	c.TestFunctionTypeWithFunctionType(`/** @constructor */ function f() {}
function abstractMethod() {}
/** @param {number} x */ f.prototype.foo = abstractMethod;
`, "(new f).foo", "function(this:f, number): ?")
}

// Java: TypeCheckTest.testFunctionInference18()
func (c *TypeCheckTest) TestFunctionInference18() {
	c.TestFunctionTypeWithFunctionType(`var goog = {};
/** @this {Date} */ goog.eatWithDate;
`, "goog.eatWithDate", "function(this:Date): ?")
}

// Java: TypeCheckTest.testFunctionInference19()
func (c *TypeCheckTest) TestFunctionInference19() {
	c.TestFunctionTypeWithFunctionType("/** @param {string} x */ var f;", "f", "function(string): ?")
}

// Java: TypeCheckTest.testFunctionInference20()
func (c *TypeCheckTest) TestFunctionInference20() {
	c.TestFunctionTypeWithFunctionType("/** @this {Date} */ var f;", "f", "function(this:Date): ?")
}

// Java: TypeCheckTest.testFunctionInference21a()
func (c *TypeCheckTest) TestFunctionInference21a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var f = function() { throw 'x' };
/** @return {boolean} */ var g = f;
`).
		Run()
}

// Java: TypeCheckTest.testFunctionInference21b()
func (c *TypeCheckTest) TestFunctionInference21b() {
	c.TestFunctionTypeWithFunctionType("var f = function() { throw 'x' };", "f", "function(): ?")
}

// Java: TypeCheckTest.testFunctionInference22()
func (c *TypeCheckTest) TestFunctionInference22() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {!Function} */ var f = function() { g(this); };
/** @param {boolean} x */ var g = function(x) {};
`).
		Run()
}

// Java: TypeCheckTest.testFunctionInference23a()
func (c *TypeCheckTest) TestFunctionInference23a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @type {!Function} */ var f = function() {
  /** @type {number} */ this.prop = 3;
};
/**
 * @param {Object} x
 * @return {string}
 */ var g = function(x) { return x.prop; };
`).
		Run()
}

// Java: TypeCheckTest.testFunctionInference23b()
func (c *TypeCheckTest) TestFunctionInference23b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {!Function} */ var f = function() {
  /** @type {number} */ this.prop = 3;
};
/**
 * @param {Object} x
 * @return {string}
 */ var g = function(x) { return x.prop; };
`).
		AddDiagnostic("Property prop never defined on Object").
		Run()
}

// Java: TypeCheckTest.testFunctionInference24()
func (c *TypeCheckTest) TestFunctionInference24() {
	c.TestFunctionType("var f = (/** number */ n, /** string= */ s) => null;", "function(number, string=): ?")
}

// Java: TypeCheckTest.testFunctionInference25()
func (c *TypeCheckTest) TestFunctionInference25() {
	c.TestFunctionType("var f = (/** number */ n, /** ...string */ s) => null;", "function(number, ...string): ?")
}

// Java: TypeCheckTest.testInnerFunction1()
func (c *TypeCheckTest) TestInnerFunction1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f() {
  /** @type {number} */ var x = 3;
  function g() { x = null; }
  return x;
}
`).
		AddDiagnostic(`assignment
found   : null
required: number
`).
		Run()
}

// Java: TypeCheckTest.testInnerFunction2()
func (c *TypeCheckTest) TestInnerFunction2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {number} */
function f() {
  var x = null;
  function g() { x = 3; }
  g();
  return x;
}
`).
		AddDiagnostic(`inconsistent return type
found   : (null|number)
required: number
`).
		Run()
}

// Java: TypeCheckTest.testInnerFunction3()
func (c *TypeCheckTest) TestInnerFunction3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var x = null;
/** @return {number} */
function f() {
  x = 3;
  /** @return {number} */
  function g() { x = true; return x; }
  return x;
}
`).
		AddDiagnostic(`inconsistent return type
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckTest.testInnerFunction4()
func (c *TypeCheckTest) TestInnerFunction4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var x = null;
/** @return {number} */
function f() {
  x = '3';
  /** @return {number} */
  function g() { x = 3; return x; }
  return x;
}
`).
		AddDiagnostic(`inconsistent return type
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testInnerFunction5()
func (c *TypeCheckTest) TestInnerFunction5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {number} */
function f() {
  var x = 3;
  /** @return {number} */
  function g() { var x = 3;x = true; return x; }
  return x;
}
`).
		AddDiagnostic(`inconsistent return type
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckTest.testInnerFunction6NullishCoalesce()
func (c *TypeCheckTest) TestInnerFunction6NullishCoalesce() {
	c.TestClosureTypes(`function f() {
 var x = null ?? function() {};
 function g() { if (goog.isFunction(x)) { x(1); } }
 g();
}
`, jutil.StringPtr(`Function x: called with 1 argument(s). Function requires at least 0 argument(s) and no more than 0 argument(s).
`))
}

// Java: TypeCheckTest.testInnerFunction7NullishCoalesce()
func (c *TypeCheckTest) TestInnerFunction7NullishCoalesce() {
	c.TestClosureTypes(`function f() {
 /** @type {function()} */
 var x = null ?? function() {};
 function g() { if (goog.isFunction(x)) { x(1); } }
 g();
}
`, jutil.StringPtr(`Function x: called with 1 argument(s). Function requires at least 0 argument(s) and no more than 0 argument(s).
`))
}

// Java: TypeCheckTest.testInnerFunction8()
func (c *TypeCheckTest) TestInnerFunction8() {
	c.TestClosureTypes(`function f() {
  function x() {};
  function g() { if (goog.isFunction(x)) { x(1); } }
  g();
}
`, jutil.StringPtr(`Function x: called with 1 argument(s). Function requires at least 0 argument(s) and no more than 0 argument(s).
`))
}

// Java: TypeCheckTest.testInnerFunction9()
func (c *TypeCheckTest) TestInnerFunction9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f() {
  var x = 3;
  function g() { x = null; };
  function h() { return x == null; }
  return h();
}
`).
		Run()
}

// Java: TypeCheckTest.testInnerFunction10()
func (c *TypeCheckTest) TestInnerFunction10() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f() {
  /** @type {?number} */ var x = null;
  /** @return {string} */
  function g() {
    if (!x) {
      x = 1;
    }
    return x;
  }
}
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testInnerFunction11()
func (c *TypeCheckTest) TestInnerFunction11() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f() {
  /** @type {?number} */ var x = null;
  /** @return {number} */
  function g() {
    x = 1;
    h();
    return x;
  }
  function h() {
    x = null;
  }
}
`).
		Run()
}

// Java: TypeCheckTest.testMethodInference1()
func (c *TypeCheckTest) TestMethodInference1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function F() {}
/** @return {number} */ F.prototype.foo = function() { return 3; };
/** @constructor
 * @extends {F} */
function G() {}
/** @override */ G.prototype.foo = function() { return true; };
`).
		AddDiagnostic(`inconsistent return type
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckTest.testMethodInference2()
func (c *TypeCheckTest) TestMethodInference2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = {};
/** @constructor */ goog.F = function() {};
/** @return {number} */ goog.F.prototype.foo = function() { return 3; };
/** @constructor
 * @extends {goog.F} */
goog.G = function() {};
/** @override */ goog.G.prototype.foo = function() { return true; };
`).
		AddDiagnostic(`inconsistent return type
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckTest.testMethodInference3()
func (c *TypeCheckTest) TestMethodInference3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function F() {}
/** @param {boolean} x
 * @return {number} */
F.prototype.foo = function(x) { return 3; };
/** @constructor
 * @extends {F} */
function G() {}
/** @override */
G.prototype.foo = function(x) { return x; };
`).
		AddDiagnostic(`inconsistent return type
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckTest.testMethodInference4()
func (c *TypeCheckTest) TestMethodInference4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function F() {}
/** @param {boolean} x
 * @return {number} */
F.prototype.foo = function(x) { return 3; };
/** @constructor
 * @extends {F} */
function G() {}
/** @override */
G.prototype.foo = function(y) { return y; };
`).
		AddDiagnostic(`inconsistent return type
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckTest.testMethodInference5()
func (c *TypeCheckTest) TestMethodInference5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function F() {}
/** @param {number} x
 * @return {string} */
F.prototype.foo = function(x) { return 'x'; };
/** @constructor
 * @extends {F} */
function G() {}
/** @type {number} */ G.prototype.num = 3;
/** @override */
G.prototype.foo = function(y) { return this.num + y; };
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testMethodInference6()
func (c *TypeCheckTest) TestMethodInference6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function F() {}
/** @param {number} x */ F.prototype.foo = function(x) { };
/** @constructor
 * @extends {F} */
function G() {}
/** @override */ G.prototype.foo = function() { };
(new G()).foo(1);
`).
		Run()
}

// Java: TypeCheckTest.testMethodInference7()
func (c *TypeCheckTest) TestMethodInference7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function F() {}
F.prototype.foo = function() { };
/** @constructor
 * @extends {F} */
function G() {}
/** @override */ G.prototype.foo = function(x, y) { };
`).
		AddDiagnostic(`mismatch of the foo property type and the type of the property it overrides from superclass F
original: function(this:F): undefined
override: function(this:G, ?, ?): undefined
`).
		Run()
}

// Java: TypeCheckTest.testMethodInference8()
func (c *TypeCheckTest) TestMethodInference8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function F() {}
F.prototype.foo = function() { };
/** @constructor
 * @extends {F} */
function G() {}
/** @override */
G.prototype.foo = function(opt_b, var_args) { };
(new G()).foo(1, 2, 3);
`).
		Run()
}

// Java: TypeCheckTest.testMethodInference9()
func (c *TypeCheckTest) TestMethodInference9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function F() {}
F.prototype.foo = function() { };
/** @constructor
 * @extends {F} */
function G() {}
/** @override */
G.prototype.foo = function(var_args, opt_b) { };
`).
		AddDiagnostic("variable length argument must be last").
		Run()
}

// Java: TypeCheckTest.testStaticMethodDeclaration1()
func (c *TypeCheckTest) TestStaticMethodDeclaration1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function F() { F.foo(true); }
/** @param {number} x */ F.foo = function(x) {};
`).
		AddDiagnostic(`actual parameter 1 of F.foo does not match formal parameter
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckTest.testStaticMethodDeclaration2()
func (c *TypeCheckTest) TestStaticMethodDeclaration2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = goog || {}; function f() { goog.foo(true); }
/** @param {number} x */ goog.foo = function(x) {};
`).
		AddDiagnostic(`actual parameter 1 of goog.foo does not match formal parameter
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckTest.testStaticMethodDeclaration3()
func (c *TypeCheckTest) TestStaticMethodDeclaration3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = goog || {}; function f() { goog.foo(true); }
goog.foo = function() {};
`).
		AddDiagnostic(`Function goog.foo: called with 1 argument(s). Function requires at least 0 argument(s) and no more than 0 argument(s).
`).
		Run()
}

// Java: TypeCheckTest.testDuplicateStaticMethodDecl1()
func (c *TypeCheckTest) TestDuplicateStaticMethodDecl1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = goog || {};
/** @param {number} x */ goog.foo = function(x) {};
/** @param {number} x */ goog.foo = function(x) {};
`).
		AddDiagnostic("variable goog.foo redefined, original definition at testcode0:2").
		Run()
}

// Java: TypeCheckTest.testDuplicateStaticMethodDecl2()
func (c *TypeCheckTest) TestDuplicateStaticMethodDecl2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = goog || {};
/** @param {number} x */ goog.foo = function(x) {};
/** @param {number} x
 * @suppress {duplicate} */
goog.foo = function(x) {};
`).
		Run()
}

// Java: TypeCheckTest.testDuplicateStaticMethodDecl3()
func (c *TypeCheckTest) TestDuplicateStaticMethodDecl3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = goog || {};
goog.foo = function(x) {};
goog.foo = function(x) {};
`).
		Run()
}

// Java: TypeCheckTest.testDuplicateStaticMethodDecl4()
func (c *TypeCheckTest) TestDuplicateStaticMethodDecl4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = goog || {};
/** @type {Function} */ goog.foo = function(x) {};
goog.foo = function(x) {};
`).
		Run()
}

// Java: TypeCheckTest.testDuplicateStaticMethodDecl5()
func (c *TypeCheckTest) TestDuplicateStaticMethodDecl5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = goog || {};
goog.foo = function(x) {};
/** @return {undefined} */ goog.foo = function(x) {};
`).
		AddDiagnostic("variable goog.foo redefined, original definition at testcode0:2").
		Run()
}

// Java: TypeCheckTest.testDuplicateStaticMethodDecl6()
func (c *TypeCheckTest) TestDuplicateStaticMethodDecl6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = goog || {};
goog.foo = function(x) {};
/**
 * @suppress {duplicate}
 * @return {undefined}
 */
goog.foo =
   /** @type {!Function} */ (function(x) {});
`).
		Run()
}

// Java: TypeCheckTest.testDuplicateStaticMethodDecl7()
func (c *TypeCheckTest) TestDuplicateStaticMethodDecl7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {string} */ var foo;
var z = function foo() {};
`).
		AddDiagnosticType(TypeCheck_FUNCTION_MASKS_VARIABLE).
		Run()
}

// Java: TypeCheckTest.testDuplicateStaticMethodDecl8()
func (c *TypeCheckTest) TestDuplicateStaticMethodDecl8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @fileoverview @suppress {duplicate} */
/** @type {string} */ var foo;
var z = function foo() {};
`).
		Run()
}

// Java: TypeCheckTest.testDuplicateStaticPropertyDecl1()
func (c *TypeCheckTest) TestDuplicateStaticPropertyDecl1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = goog || {};
/** @type {Foo} */ goog.foo;
/** @type {Foo} */ goog.foo;
/** @constructor */ function Foo() {}
`).
		Run()
}

// Java: TypeCheckTest.testDuplicateStaticPropertyDecl2()
func (c *TypeCheckTest) TestDuplicateStaticPropertyDecl2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = goog || {};
/** @type {Foo} */ goog.foo;
/** @type {Foo}
 * @suppress {duplicate} */ goog.foo;
/** @constructor */ function Foo() {}
`).
		Run()
}

// Java: TypeCheckTest.testDuplicateStaticPropertyDecl3()
func (c *TypeCheckTest) TestDuplicateStaticPropertyDecl3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = goog || {};
/** @type {!Foo} */ goog.foo;
/** @type {string} */ goog.foo;
/** @constructor */ function Foo() {}
`).
		AddDiagnostic(`variable goog.foo redefined with type string, original definition at testcode0:2 with type Foo
`).
		Run()
}

// Java: TypeCheckTest.testDuplicateStaticPropertyDecl4()
func (c *TypeCheckTest) TestDuplicateStaticPropertyDecl4() {
	c.TestClosureTypesMultipleWarnings(`/** @type {!Foo} */ goog.foo;
/** @type {string} */ goog.foo = 'x';
/** @constructor */ function Foo() {}
`, []string{`assignment to property foo of goog
found   : string
required: Foo`, "variable goog.foo redefined with type string, original definition at testcode:1 with type Foo"})
}

// Java: TypeCheckTest.testDuplicateStaticPropertyDecl5()
func (c *TypeCheckTest) TestDuplicateStaticPropertyDecl5() {
	c.TestClosureTypesMultipleWarnings(`var goog = goog || {};
/** @type {!Foo} */ goog.foo;
/** @type {string}
 * @suppress {duplicate} */ goog.foo = 'x';
/** @constructor */ function Foo() {}
`, []string{`assignment to property foo of goog
found   : string
required: Foo`})
}

// Java: TypeCheckTest.testDuplicateStaticPropertyDecl6()
func (c *TypeCheckTest) TestDuplicateStaticPropertyDecl6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = goog || {};
/** @type {string} */ goog.foo = 'y';
/** @type {string}
 * @suppress {duplicate} */ goog.foo = 'x';
`).
		Run()
}

// Java: TypeCheckTest.testDuplicateStaticPropertyDecl7()
func (c *TypeCheckTest) TestDuplicateStaticPropertyDecl7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = goog || {};
/** @param {string} x */ goog.foo;
/** @type {function(string)} */ goog.foo;
`).
		Run()
}

// Java: TypeCheckTest.testDuplicateStaticPropertyDecl8()
func (c *TypeCheckTest) TestDuplicateStaticPropertyDecl8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = goog || {};
/** @return {EventCopy} */ goog.foo;
/** @constructor */ function EventCopy() {}
/** @return {EventCopy} */ goog.foo;
`).
		Run()
}

// Java: TypeCheckTest.testDuplicateStaticPropertyDecl9()
func (c *TypeCheckTest) TestDuplicateStaticPropertyDecl9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = goog || {};
/** @return {EventCopy} */ goog.foo;
/** @return {EventCopy} */ goog.foo;
/** @constructor */ function EventCopy() {}
`).
		Run()
}

// Java: TypeCheckTest.testDuplicateStaticPropertyDec20()
func (c *TypeCheckTest) TestDuplicateStaticPropertyDec20() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @fileoverview
 * @suppress {duplicate}
 */
var goog = goog || {};
/** @type {string} */ goog.foo = 'y';
/** @type {string} */ goog.foo = 'x';
`).
		Run()
}

// Java: TypeCheckTest.testDuplicateLocalVarDecl()
func (c *TypeCheckTest) TestDuplicateLocalVarDecl() {
	c.TestClosureTypesMultipleWarnings(`/** @param {number} x */
function f(x) { /** @type {string} */ var x = ''; }
`, []string{"variable x redefined with type string, original definition at testcode:2 with type number", `initializing variable
found   : string
required: number`})
}

// Java: TypeCheckTest.testDuplicateLocalVarDeclSuppressed()
func (c *TypeCheckTest) TestDuplicateLocalVarDeclSuppressed() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @suppress {duplicate} */
function f(x) {
  var x = x;
  x.y = true;
}
`).
		Run()
}

// Java: TypeCheckTest.testShadowBleedingFunctionName()
func (c *TypeCheckTest) TestShadowBleedingFunctionName() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var f = function x() {
  var x;
  var /** undefined */ y = x;
};
`).
		Run()
}

// Java: TypeCheckTest.testDuplicateInstanceMethod1()
func (c *TypeCheckTest) TestDuplicateInstanceMethod1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function F() {}
F.prototype.bar = function() {};
F.prototype.bar = function() {};
`).
		Run()
}

// Java: TypeCheckTest.testDuplicateInstanceMethod2()
func (c *TypeCheckTest) TestDuplicateInstanceMethod2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function F() {}
/** jsdoc */ F.prototype.bar = function() {};
/** jsdoc */ F.prototype.bar = function() {};
`).
		Run()
}

// Java: TypeCheckTest.testDuplicateInstanceMethod3()
func (c *TypeCheckTest) TestDuplicateInstanceMethod3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function F() {}
F.prototype.bar = function() {};
/** jsdoc */ F.prototype.bar = function() {};
`).
		Run()
}

// Java: TypeCheckTest.testDuplicateInstanceMethod4()
func (c *TypeCheckTest) TestDuplicateInstanceMethod4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function F() {}
/** jsdoc */ F.prototype.bar = function() {};
F.prototype.bar = function() {};
`).
		Run()
}

// Java: TypeCheckTest.testDuplicateInstanceMethod5()
func (c *TypeCheckTest) TestDuplicateInstanceMethod5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function F() {}
/** jsdoc
 * @return {number} */ F.prototype.bar = function() {
  return 3;
};
/** jsdoc
 * @suppress {duplicate} */
F.prototype.bar = function() { return ''; };
`).
		AddDiagnostic(`inconsistent return type
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testDuplicateInstanceMethod6()
func (c *TypeCheckTest) TestDuplicateInstanceMethod6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function F() {}
/** jsdoc
 * @return {number} */ F.prototype.bar = function() {
  return 3;
};
/** jsdoc
 * @return {string}
 * @suppress {duplicate} */
F.prototype.bar = function() { return ''; };
`).
		AddDiagnostic(`assignment to property bar of F.prototype
found   : function(this:F): string
required: function(this:F): number
`).
		Run()
}

// Java: TypeCheckTest.testStubFunctionDeclaration1()
func (c *TypeCheckTest) TestStubFunctionDeclaration1() {
	c.TestFunctionTypeWithFunctionType(`/** @constructor */ function f() {};
/** @param {number} x
 * @param {string} y
 * @return {number} */ f.prototype.foo;
`, "(new f).foo", "function(this:f, number, string): number")
}

// Java: TypeCheckTest.testStubFunctionDeclaration2()
func (c *TypeCheckTest) TestStubFunctionDeclaration2() {
	c.TestExternFunctionType(`/** @constructor */ function f() {};
/** @constructor
 * @extends {f} */ f.subclass;
`, "f.subclass", "function(new:f.subclass): ?")
}

// Java: TypeCheckTest.testStubFunctionDeclaration_static_onEs5Class()
func (c *TypeCheckTest) TestStubFunctionDeclaration_static_onEs5Class() {
	c.TestFunctionTypeWithFunctionType(`/** @constructor */ function f() {};
/** @return {undefined} */ f.foo;
`, "f.foo", "function(): undefined")
}

// Java: TypeCheckTest.testStubFunctionDeclaration4()
func (c *TypeCheckTest) TestStubFunctionDeclaration4() {
	c.TestFunctionTypeWithFunctionType(`/** @constructor */ function f() {
  /** @return {number} */ this.foo;
}
`, "(new f).foo", "function(this:f): number")
}

// Java: TypeCheckTest.testStubFunctionDeclaration5()
func (c *TypeCheckTest) TestStubFunctionDeclaration5() {
	c.TestFunctionTypeWithFunctionType(`/** @constructor */ function f() {
  /** @type {Function} */ this.foo;
}
`, "(new f).foo", c.CreateNullableType(c.GetNativeFunctionType()).String())
}

// Java: TypeCheckTest.testStubFunctionDeclaration6()
func (c *TypeCheckTest) TestStubFunctionDeclaration6() {
	c.TestFunctionTypeWithFunctionType(`/** @constructor */ function f() {}
/** @type {Function} */ f.prototype.foo;
`, "(new f).foo", c.CreateNullableType(c.GetNativeFunctionType()).String())
}

// Java: TypeCheckTest.testStubFunctionDeclaration7()
func (c *TypeCheckTest) TestStubFunctionDeclaration7() {
	c.TestFunctionTypeWithFunctionType(`/** @constructor */ function f() {}
/** @type {Function} */ f.prototype.foo = function() {};
`, "(new f).foo", c.CreateNullableType(c.GetNativeFunctionType()).String())
}

// Java: TypeCheckTest.testStubFunctionDeclaration8()
func (c *TypeCheckTest) TestStubFunctionDeclaration8() {
	c.TestFunctionTypeWithFunctionType("/** @type {Function} */ var f = function() {}; ", "f", c.CreateNullableType(c.GetNativeFunctionType()).String())
}

// Java: TypeCheckTest.testStubFunctionDeclaration9()
func (c *TypeCheckTest) TestStubFunctionDeclaration9() {
	c.TestFunctionTypeWithFunctionType("/** @type {function():number} */ var f; ", "f", "function(): number")
}

// Java: TypeCheckTest.testStubFunctionDeclaration10()
func (c *TypeCheckTest) TestStubFunctionDeclaration10() {
	c.TestFunctionTypeWithFunctionType("/** @type {function(number):number} */ var f = function(x) {};", "f", "function(number): number")
}

// Java: TypeCheckTest.testStubMethodDeclarationDoesntBlockTypechecking_1()
func (c *TypeCheckTest) TestStubMethodDeclarationDoesntBlockTypechecking_1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
function Foo() {}
/** @return {number} */
Foo.prototype.method = function() {};
/**
 * @constructor
 * @implements {Foo}
 */
function Bar() {}
Bar.prototype.method;
var /** null */ n = (new Bar).method();
`).
		AddDiagnostic(`initializing variable
found   : number
required: null
`).
		Run()
}

// Java: TypeCheckTest.testStubMethodDeclarationDoesntBlockTypechecking_2()
func (c *TypeCheckTest) TestStubMethodDeclarationDoesntBlockTypechecking_2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */
function Foo() {}
/** @return {number} */
Foo.prototype.method = function() {};
/**
 * @constructor
 * @extends {Foo}
 */
function Bar() {}
Bar.prototype.method;
var /** null */ n = (new Bar).method();
`).
		AddDiagnostic(`initializing variable
found   : number
required: null
`).
		Run()
}

// Java: TypeCheckTest.testStubMethodDeclarationDoesntBlockTypechecking_3()
func (c *TypeCheckTest) TestStubMethodDeclarationDoesntBlockTypechecking_3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
var Foo = function() {};
/** @type {number} */
Foo.prototype.num;
/**
 * @constructor
 * @implements {Foo}
 */
var Bar = function() {};
/** @type {?} */
Bar.prototype.num;
var /** string */ x = (new Bar).num;
`).
		Run()
}

// Java: TypeCheckTest.testNestedFunctionInference1()
func (c *TypeCheckTest) TestNestedFunctionInference1() {
	nestedAssignOfFooAndBar := `/** @constructor */
function f() {};
f.prototype.foo = f.prototype.bar = function(){};
`
	c.TestFunctionTypeWithFunctionType(nestedAssignOfFooAndBar, "(new f).bar", "function(this:f): undefined")
}

// Java: TypeCheckTest.testFunctionType(String, String)
func (c *TypeCheckTest) TestFunctionType(functionDef string, functionType string) {
	c.TestFunctionTypeWithFunctionType(functionDef, "f", functionType)
}

// Java: TypeCheckTest.testFunctionType(String, String, String)
func (c *TypeCheckTest) TestFunctionTypeWithFunctionType(functionDef string, functionName string, functionType string) {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(functionDef + "/** @type {number} */var a=" + functionName + ";").
		AddDiagnostic(`initializing variable
` + ("found   : " + functionType + `
`) + "required: number").
		Run()
}

// Java: TypeCheckTest.testExternFunctionType(String, String, String)
func (c *TypeCheckTest) TestExternFunctionType(functionDef string, functionName string, functionType string) {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(functionDef).
		AddSource("/** @type {number} */var a=" + functionName + ";").
		AddDiagnostic(`initializing variable
` + ("found   : " + functionType + `
`) + "required: number").
		Run()
}

// Java: TypeCheckTest.testTypeRedefinition()
func (c *TypeCheckTest) TestTypeRedefinition() {
	c.TestClosureTypesMultipleWarnings(`a={};/**@enum {string}*/ a.A = {ZOR:'b'};
/** @constructor */ a.A = function() {}
`, []string{"variable a.A redefined with type (typeof a.A), original definition at testcode:1 with type enum{a.A}", `assignment to property A of a
found   : (typeof a.A)
required: enum{a.A}`})
}

// Java: TypeCheckTest.testIn1()
func (c *TypeCheckTest) TestIn1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("'foo' in Object").
		Run()
}

// Java: TypeCheckTest.testIn2()
func (c *TypeCheckTest) TestIn2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("3 in Object").
		Run()
}

// Java: TypeCheckTest.testIn3()
func (c *TypeCheckTest) TestIn3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("undefined in Object").
		Run()
}

// Java: TypeCheckTest.testIn4()
func (c *TypeCheckTest) TestIn4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("Date in Object").
		AddDiagnostic(`left side of 'in'
found   : function(new:Date, ?=, ?=, ?=, ?=, ?=, ?=, ?=): string
required: (string|symbol)
`).
		Run()
}

// Java: TypeCheckTest.testIn5()
func (c *TypeCheckTest) TestIn5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("'x' in null").
		AddDiagnostic(`'in' requires an object
found   : null
required: Object
`).
		Run()
}

// Java: TypeCheckTest.testIn6()
func (c *TypeCheckTest) TestIn6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {number} x */
function g(x) {}
g(1 in {});
`).
		AddDiagnostic(`actual parameter 1 of g does not match formal parameter
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckTest.testIn7()
func (c *TypeCheckTest) TestIn7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {number} x
 * @return {number}
 */
function g(x) { return 5; }
function f() {
  var x = {};
  x.foo = '3';
  return g(x.foo) in {};
}
`).
		AddDiagnostic(`actual parameter 1 of g does not match formal parameter
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testInWithThis_narrowsPropertyType()
func (c *TypeCheckTest) TestInWithThis_narrowsPropertyType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */
function Foo() {}
Foo.prototype.method = function() {
  if ('x' in this) {
    return this.x; // this access is allowed
  }
  return this.y; // this access causes a warning
}
`).
		AddDiagnostic("Property y never defined on Foo").
		Run()
}

// Java: TypeCheckTest.testInWithWellKnownSymbol()
func (c *TypeCheckTest) TestInWithWellKnownSymbol() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("Symbol.iterator in Object").
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testInWithUniqueSymbol()
func (c *TypeCheckTest) TestInWithUniqueSymbol() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("Symbol('foo') in Object").
		Run()
}

// Java: TypeCheckTest.testInWithSymbol()
func (c *TypeCheckTest) TestInWithSymbol() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function f(/** symbol */ s) { return s in {}; }").
		Run()
}

// Java: TypeCheckTest.testForIn1()
func (c *TypeCheckTest) TestForIn1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {boolean} x */ function f(x) {}
for (var k in {}) {
  f(k);
}
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : string
required: boolean
`).
		Run()
}

// Java: TypeCheckTest.testForIn2()
func (c *TypeCheckTest) TestForIn2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {boolean} x */ function f(x) {}
/** @enum {string} */ var E = {FOO: 'bar'};
/** @type {Object<E, string>} */ var obj = {};
var k = null;
for (k in obj) {
  f(k);
}
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : E<string>
required: boolean
`).
		Run()
}

// Java: TypeCheckTest.testForIn3()
func (c *TypeCheckTest) TestForIn3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {boolean} x */ function f(x) {}
/** @type {Object<number>} */ var obj = {};
for (var k in obj) {
  f(obj[k]);
}
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : number
required: boolean
`).
		Run()
}

// Java: TypeCheckTest.testForIn4()
func (c *TypeCheckTest) TestForIn4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {boolean} x */ function f(x) {}
/** @enum {string} */ var E = {FOO: 'bar'};
/** @type {Object<E, Array>} */ var obj = {};
for (var k in obj) {
  f(obj[k]);
}
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : (Array|null)
required: boolean
`).
		Run()
}

// Java: TypeCheckTest.testForIn5()
func (c *TypeCheckTest) TestForIn5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddObject().Build()).
		AddSource(`/** @param {boolean} x */ function f(x) {}
/** @constructor */ var E = function(){};
/** @override */ E.prototype.toString = function() { return ''; };
/** @type {Object<!E, number>} */ var obj = {};
for (var k in obj) {
  f(k);
}
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : string
required: boolean
`).
		Run()
}

// Java: TypeCheckTest.testForOf1()
func (c *TypeCheckTest) TestForOf1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {!Iterable<number>} */
var it = [1, 2, 3];
for (let x of it) {
  alert(x);
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testForOf2()
func (c *TypeCheckTest) TestForOf2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {boolean} x */ function f(x) {}
/** @type {!Iterable<number>} */ var it = [1, 2, 3];
for (let x of it) { f(x); }
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : number
required: boolean
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testForOf3()
func (c *TypeCheckTest) TestForOf3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {?Iterable<number>} */ var it = [1, 2, 3];
for (let x of it) {}
`).
		AddDiagnostic(`Can only iterate over a (non-null) Iterable type
found   : (Iterable<number,?,?>|null)
required: Iterable
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testForOf4()
func (c *TypeCheckTest) TestForOf4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** !IteratorLike<number> */ it) {
  for (let x of it) {}
}
`).
		AddDiagnostic(`Can only iterate over a (non-null) Iterable type
found   : IteratorLike<number,?,?>
required: Iterable
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testForOf5()
func (c *TypeCheckTest) TestForOf5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** string */ ch, /** string */ str) {
  for (ch of str) {}
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testForOf6()
func (c *TypeCheckTest) TestForOf6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** !Array<number> */ a) {
  for (let elem of a) {}
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testForOf7()
func (c *TypeCheckTest) TestForOf7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("for (let elem of ['array', 'literal']) {}").
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testForOf8()
func (c *TypeCheckTest) TestForOf8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** null */ x) {}
(function(/** !Array<string>|!Iterable<number> */ it) {
  for (let elem of it) {
    f(elem);
  }
})(['']);
`).
		IncludeDefaultExterns().
		AddDiagnosticType(DiagnosticType_Error("JSC_TYPE_MISMATCH", `actual parameter 1 of f does not match formal parameter
found   : (number|string)
required: null at testcode0 line 4 : 6
`)).
		Run()
}

// Java: TypeCheckTest.testUnionIndexAccessStringElem()
func (c *TypeCheckTest) TestUnionIndexAccessStringElem() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const {!ReadonlyArray<string>|!Array<string>} */ const a = [];
/** @const {!null} */ const b = a[0];
`).
		IncludeDefaultExterns().
		AddDiagnosticType(DiagnosticType_Error("JSC_TYPE_MISMATCH", `initializing variable
found   : string
required: None at testcode0 line 1 : 114
`)).
		Run()
}

// Java: TypeCheckTest.testUnionIndexAccessOverlappingElem()
func (c *TypeCheckTest) TestUnionIndexAccessOverlappingElem() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const {!ReadonlyArray<string|undefined>|!Array<string>} */ const a = [];
/** @const {!undefined} */ const b = a[0];
`).
		IncludeDefaultExterns().
		AddDiagnosticType(DiagnosticType_Error("JSC_TYPE_MISMATCH", `initializing variable
found   : (string|undefined)
required: None at testcode0 line 1 : 114
`)).
		Run()
}

// Java: TypeCheckTest.testUnionIndexAccessDisjointElem()
func (c *TypeCheckTest) TestUnionIndexAccessDisjointElem() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const {!ReadonlyArray<string|number>|!Array<boolean>} */ const a = [];
/** @const {!undefined} */ const b = a[0];
`).
		IncludeDefaultExterns().
		AddDiagnosticType(DiagnosticType_Error("JSC_TYPE_MISMATCH", `initializing variable
found   : (boolean|number|string)
required: None at testcode0 line 1 : 114
`)).
		Run()
}

// Java: TypeCheckTest.testReadonlyArrayUnionIndexAccessStringElem()
func (c *TypeCheckTest) TestReadonlyArrayUnionIndexAccessStringElem() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const {!ReadonlyArray<string>|!ReadonlyArray<string>} */ const a = [];
/** @const {!null} */ const b = a[0];
`).
		IncludeDefaultExterns().
		AddDiagnosticType(DiagnosticType_Error("JSC_TYPE_MISMATCH", `initializing variable
found   : string
required: None at testcode0 line 1 : 114
`)).
		Run()
}

// Java: TypeCheckTest.testReadonlyArrayUnionIndexAccessOverlappingElem()
func (c *TypeCheckTest) TestReadonlyArrayUnionIndexAccessOverlappingElem() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const {!ReadonlyArray<string|undefined>|!ReadonlyArray<string>} */ const a = [];
/** @const {!undefined} */ const b = a[0];
`).
		IncludeDefaultExterns().
		AddDiagnosticType(DiagnosticType_Error("JSC_TYPE_MISMATCH", `initializing variable
found   : (string|undefined)
required: None at testcode0 line 1 : 114
`)).
		Run()
}

// Java: TypeCheckTest.testReadonlyArrayUnionIndexAccessDisjointElem()
func (c *TypeCheckTest) TestReadonlyArrayUnionIndexAccessDisjointElem() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const {!ReadonlyArray<string|number>|!ReadonlyArray<boolean>} */ const a = [];
/** @const {!undefined} */ const b = a[0];
`).
		IncludeDefaultExterns().
		AddDiagnosticType(DiagnosticType_Error("JSC_TYPE_MISMATCH", `initializing variable
found   : (boolean|number|string)
required: None at testcode0 line 1 : 114
`)).
		Run()
}

// Java: TypeCheckTest.testUnionIndexAccessIncompatibleBecomesUnknown()
func (c *TypeCheckTest) TestUnionIndexAccessIncompatibleBecomesUnknown() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const {!ReadonlyArray<string>|number} */ const a = [];
/** @const {!undefined} */ const b = a[0];
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testComparison1()
func (c *TypeCheckTest) TestComparison1() {
	c.T().Skip("@Ignore")
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**@type null */var a;
/**@type !Date */var b;
if (a==b) {}
`).
		AddDiagnostic(`condition always evaluates to false
left : null
right: Date
`).
		Run()
}

// Java: TypeCheckTest.testComparison2()
func (c *TypeCheckTest) TestComparison2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**@type {number}*/var a;
/**@type {!Date} */var b;
if (a!==b) {}
`).
		AddDiagnostic(`condition always evaluates to true
left : number
right: Date
`).
		Run()
}

// Java: TypeCheckTest.testComparison3()
func (c *TypeCheckTest) TestComparison3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {(Object|undefined)} */var a;
var b = a == null
`).
		Run()
}

// Java: TypeCheckTest.testComparison4()
func (c *TypeCheckTest) TestComparison4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {(!Object|undefined)} */var a;
/** @type {!Object} */var b;
var c = a == b
`).
		Run()
}

// Java: TypeCheckTest.testComparison5()
func (c *TypeCheckTest) TestComparison5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {null} */var a;
/** @type {null} */var b;
a == b
`).
		AddDiagnostic(`condition always evaluates to true
left : null
right: null
`).
		Run()
}

// Java: TypeCheckTest.testComparison6()
func (c *TypeCheckTest) TestComparison6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {null} */var a;
/** @type {null} */var b;
a != b
`).
		AddDiagnostic(`condition always evaluates to false
left : null
right: null
`).
		Run()
}

// Java: TypeCheckTest.testComparison7()
func (c *TypeCheckTest) TestComparison7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var a;
var b;
a == b
`).
		AddDiagnostic(`condition always evaluates to true
left : undefined
right: undefined
`).
		Run()
}

// Java: TypeCheckTest.testComparison8()
func (c *TypeCheckTest) TestComparison8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {Array<string>} */ var a = [];
a[0] == null || a[1] == undefined
`).
		Run()
}

// Java: TypeCheckTest.testComparison9()
func (c *TypeCheckTest) TestComparison9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {Array<undefined>} */ var a = [];
a[0] == null
`).
		AddDiagnostic(`condition always evaluates to true
left : undefined
right: null
`).
		Run()
}

// Java: TypeCheckTest.testComparison10()
func (c *TypeCheckTest) TestComparison10() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {Array<undefined>} */ var a = [];
a[0] === null
`).
		Run()
}

// Java: TypeCheckTest.testComparison11()
func (c *TypeCheckTest) TestComparison11() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("(function(){}) == 'x'").
		AddDiagnostic(`condition always evaluates to false
left : function(): undefined
right: string
`).
		Run()
}

// Java: TypeCheckTest.testComparison12()
func (c *TypeCheckTest) TestComparison12() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("(function(){}) == 3").
		AddDiagnostic(`condition always evaluates to false
left : function(): undefined
right: number
`).
		Run()
}

// Java: TypeCheckTest.testComparison13()
func (c *TypeCheckTest) TestComparison13() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("(function(){}) == false").
		AddDiagnostic(`condition always evaluates to false
left : function(): undefined
right: boolean
`).
		Run()
}

// Java: TypeCheckTest.testComparison14()
func (c *TypeCheckTest) TestComparison14() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {function((Array|string), Object): number} */
function f(x, y) { return x === y; }
`).
		AddDiagnostic(`inconsistent return type
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckTest.testComparison15()
func (c *TypeCheckTest) TestComparison15() {
	c.TestClosureTypes(`/** @constructor */ function F() {}
/**
 * @param {number} x
 * @constructor
 * @extends {F}
 */
function G(x) {}
goog.inherits(G, F);
/**
 * @param {number} x
 * @constructor
 * @extends {G}
 */
function H(x) {}
goog.inherits(H, G);
/** @param {G} x */
function f(x) { return x.constructor === H; }
`, nil)
}

// Java: TypeCheckTest.testComparisonWithTemplatizedSubclass()
func (c *TypeCheckTest) TestComparisonWithTemplatizedSubclass() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @template T @constructor */
function Base() {}
/**
 * @template T
 * @constructor
 * @extends {Base<T>}
 */
function Sub() {}
/** @type {!Base<number>} */ var base = new Base();
/** @type {!Sub<number>} */ var sub = new Sub();
if (base === sub) {}
`).
		Run()
}

// Java: TypeCheckTest.testDeleteOperator1()
func (c *TypeCheckTest) TestDeleteOperator1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var x = {};
/** @return {string} */ function f() { return delete x['a']; }
`).
		AddDiagnostic(`inconsistent return type
found   : boolean
required: string
`).
		Run()
}

// Java: TypeCheckTest.testDeleteOperator2()
func (c *TypeCheckTest) TestDeleteOperator2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var obj = {};
/**
 * @param {string} x
 * @return {Object} */ function f(x) { return obj; }
/** @param {?number} x */ function g(x) {
  if (x) { delete f(x)['a']; }
}
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testEnumStaticMethod1()
func (c *TypeCheckTest) TestEnumStaticMethod1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum */ var Foo = {AAA: 1};
/** @param {number} x */ Foo.method = function(x) {};
Foo.method(true);
`).
		AddDiagnostic(`actual parameter 1 of Foo.method does not match formal parameter
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckTest.testEnumStaticMethod2()
func (c *TypeCheckTest) TestEnumStaticMethod2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum */ var Foo = {AAA: 1};
/** @param {number} x */ Foo.method = function(x) {};
function f() { Foo.method(true); }
`).
		AddDiagnostic(`actual parameter 1 of Foo.method does not match formal parameter
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckTest.testEnum1()
func (c *TypeCheckTest) TestEnum1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**@enum*/var a={BB:1,CC:2};
/**@type {a}*/var d;d=a.BB;
`).
		Run()
}

// Java: TypeCheckTest.testEnum3()
func (c *TypeCheckTest) TestEnum3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@enum*/var a={BB:1,BB:2}").
		AddDiagnostic("variable a.BB redefined, original definition at testcode0:1").
		Run()
}

// Java: TypeCheckTest.testEnum4()
func (c *TypeCheckTest) TestEnum4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@enum*/var a={BB:'string'}").
		AddDiagnostic(`assignment to property BB of enum{a}
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testEnum5()
func (c *TypeCheckTest) TestEnum5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@enum {?String}*/var a={BB:'string'}").
		AddDiagnostic(`assignment to property BB of enum{a}
found   : string
required: (String|null)
`).
		Run()
}

// Java: TypeCheckTest.testEnum6()
func (c *TypeCheckTest) TestEnum6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**@enum*/var a={BB:1,CC:2};
/**@type {!Array}*/var d;d=a.BB;
`).
		AddDiagnostic(`assignment
found   : a<number>
required: Array`).
		Run()
}

// Java: TypeCheckTest.testEnum7()
func (c *TypeCheckTest) TestEnum7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum */var a={AA:1,BB:2,CC:3};
/** @type {a} */var b=a.D;
`).
		AddDiagnostic("element D does not exist on this enum").
		Run()
}

// Java: TypeCheckTest.testEnum8()
func (c *TypeCheckTest) TestEnum8() {
	c.TestClosureTypesMultipleWarnings("/** @enum */var a=8;", []string{"enum initializer must be an object literal or an enum", `initializing variable
found   : number
required: enum{a}`})
}

// Java: TypeCheckTest.testEnum9()
func (c *TypeCheckTest) TestEnum9() {
	c.TestClosureTypesMultipleWarnings("/** @enum */ goog.a=8;", []string{`assignment to property a of goog
found   : number
required: enum{goog.a}`, "enum initializer must be an object literal or an enum"})
}

// Java: TypeCheckTest.testEnum10()
func (c *TypeCheckTest) TestEnum10() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum {number} */
goog.K = { A : 3 };
`).
		Run()
}

// Java: TypeCheckTest.testEnum11()
func (c *TypeCheckTest) TestEnum11() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum {number} */
goog.K = { 502 : 3 };
`).
		Run()
}

// Java: TypeCheckTest.testEnum12()
func (c *TypeCheckTest) TestEnum12() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum {number} */ var a = {};
/** @enum */ var b = a;
`).
		Run()
}

// Java: TypeCheckTest.testEnum13a()
func (c *TypeCheckTest) TestEnum13a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum {number} */ var a = {};
/** @enum {string} */ var b = a;
`).
		AddDiagnostic(`incompatible enum element types
found   : number
required: string`).
		Run()
}

// Java: TypeCheckTest.testEnum13b()
func (c *TypeCheckTest) TestEnum13b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum {number} */ var a = {};
/** @const */ var ns = {};
/** @enum {string} */ ns.b = a;
`).
		AddDiagnostic(`incompatible enum element types
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testEnum14()
func (c *TypeCheckTest) TestEnum14() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum {number} */ var a = {FOO:5};
/** @enum */ var b = a;
var c = b.FOO;
`).
		Run()
}

// Java: TypeCheckTest.testEnum15()
func (c *TypeCheckTest) TestEnum15() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum {number} */ var a = {FOO:5};
/** @enum */ var b = a;
var c = b.BAR;
`).
		AddDiagnostic("element BAR does not exist on this enum").
		Run()
}

// Java: TypeCheckTest.testEnum16()
func (c *TypeCheckTest) TestEnum16() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = {};
/**@enum*/goog .a={BB:1,BB:2}
`).
		AddDiagnostic("variable goog.a.BB redefined, original definition at testcode0:2").
		Run()
}

// Java: TypeCheckTest.testEnum17()
func (c *TypeCheckTest) TestEnum17() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = {};
/**@enum*/goog.a={BB:'string'}
`).
		AddDiagnostic(`assignment to property BB of enum{goog.a}
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testEnum18()
func (c *TypeCheckTest) TestEnum18() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**@enum*/ var E = {A: 1, B: 2};
/** @param {!E} x
@return {number} */
var f = function(x) { return x; };
`).
		Run()
}

// Java: TypeCheckTest.testEnum19()
func (c *TypeCheckTest) TestEnum19() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**@enum*/ var E = {A: 1, B: 2};
/** @param {number} x
@return {!E} */
var f = function(x) { return x; };
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: E<number>
`).
		Run()
}

// Java: TypeCheckTest.testEnum20()
func (c *TypeCheckTest) TestEnum20() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@enum*/ var E = {A: 1, B: 2}; var x = []; x[E.A] = 0;").
		Run()
}

// Java: TypeCheckTest.testEnum21()
func (c *TypeCheckTest) TestEnum21() {
	var n *rhino.Node = c.ParseAndTypeCheck(`/** @enum {string} */ var E = {A : 'a', B : 'b'};
/** @param {!E} x
@return {!E} */ function f(x) { return x; }
`)
	var nodeX *rhino.Node = n.GetLastChild().GetLastChild().
		GetLastChild().
		GetLastChild()
	var typeE rhino.JSType = nodeX.GetJSType()
	AssertThat(typeE.IsObject()).IsFalse()
	AssertThat(typeE.IsNullable()).IsFalse()
}

// Java: TypeCheckTest.testEnum22()
func (c *TypeCheckTest) TestEnum22() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**@enum*/ var E = {A: 1, B: 2};
/** @param {E} x
* @return {number} */ function f(x) {return x}
`).
		Run()
}

// Java: TypeCheckTest.testEnum23()
func (c *TypeCheckTest) TestEnum23() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**@enum*/ var E = {A: 1, B: 2};
/** @param {E} x
* @return {string} */ function f(x) {return x}
`).
		AddDiagnostic(`inconsistent return type
found   : E<number>
required: string
`).
		Run()
}

// Java: TypeCheckTest.testEnum24()
func (c *TypeCheckTest) TestEnum24() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**@enum {?Object} */ var E = {A: {}};
/** @param {E} x
* @return {!Object} */ function f(x) {return x}
`).
		AddDiagnostic(`inconsistent return type
found   : E<(Object|null)>
required: Object
`).
		Run()
}

// Java: TypeCheckTest.testEnum25()
func (c *TypeCheckTest) TestEnum25() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**@enum {!Object} */ var E = {A: {}};
/** @param {E} x
* @return {!Object} */ function f(x) {return x}
`).
		Run()
}

// Java: TypeCheckTest.testEnum26()
func (c *TypeCheckTest) TestEnum26() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var a = {}; /**@enum*/ a.B = {A: 1, B: 2};
/** @param {a.B} x
* @return {number} */ function f(x) {return x}
`).
		Run()
}

// Java: TypeCheckTest.testEnum27()
func (c *TypeCheckTest) TestEnum27() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum */ var A = {B: 1, C: 2};
function f(x) { return A == x; }
`).
		Run()
}

// Java: TypeCheckTest.testEnum28()
func (c *TypeCheckTest) TestEnum28() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum */ var A = {B: 1, C: 2};
function f(x) { return A.B == x; }
`).
		Run()
}

// Java: TypeCheckTest.testEnum29()
func (c *TypeCheckTest) TestEnum29() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum */ var A = {B: 1, C: 2};
/** @return {number} */ function f() { return A; }
`).
		AddDiagnostic(`inconsistent return type
found   : enum{A}
required: number
`).
		Run()
}

// Java: TypeCheckTest.testEnum30()
func (c *TypeCheckTest) TestEnum30() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum */ var A = {B: 1, C: 2};
/** @return {number} */ function f() { return A.B; }
`).
		Run()
}

// Java: TypeCheckTest.testEnum31()
func (c *TypeCheckTest) TestEnum31() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum */ var A = {B: 1, C: 2};
/** @return {A} */ function f() { return A; }
`).
		AddDiagnostic(`inconsistent return type
found   : enum{A}
required: A<number>
`).
		Run()
}

// Java: TypeCheckTest.testEnum32()
func (c *TypeCheckTest) TestEnum32() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum */ var A = {B: 1, C: 2};
/** @return {A} */ function f() { return A.B; }
`).
		Run()
}

// Java: TypeCheckTest.testEnum34()
func (c *TypeCheckTest) TestEnum34() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum */ var A = {B: 1, C: 2};
/** @param {number} x */ function f(x) { return x == A.B; }
`).
		Run()
}

// Java: TypeCheckTest.testEnum35()
func (c *TypeCheckTest) TestEnum35() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var a = a || {}; /** @enum */ a.b = {C: 1, D: 2};
/** @return {a.b} */ function f() { return a.b.C; }
`).
		Run()
}

// Java: TypeCheckTest.testEnum36()
func (c *TypeCheckTest) TestEnum36() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var a = a || {}; /** @enum */ a.b = {C: 1, D: 2};
/** @return {!a.b} */ function f() { return 1; }
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: a.b<number>
`).
		Run()
}

// Java: TypeCheckTest.testEnum37()
func (c *TypeCheckTest) TestEnum37() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = goog || {};
/** @enum {number} */ goog.a = {};
/** @enum */ var b = goog.a;
`).
		Run()
}

// Java: TypeCheckTest.testEnum38()
func (c *TypeCheckTest) TestEnum38() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum {MyEnum} */ var MyEnum = {};
/** @param {MyEnum} x */ function f(x) {}
`).
		AddDiagnostic("Cycle detected in inheritance chain of type MyEnum").
		Run()
}

// Java: TypeCheckTest.testEnum39()
func (c *TypeCheckTest) TestEnum39() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum {Number} */ var MyEnum = {FOO: new Number(1)};
/** @param {MyEnum} x
 * @return {number} */
function f(x) { return x == MyEnum.FOO && MyEnum.FOO == x; }
`).
		AddDiagnostic(`inconsistent return type
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckTest.testEnum40()
func (c *TypeCheckTest) TestEnum40() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum {Number} */ var MyEnum = {FOO: new Number(1)};
/** @param {number} x
 * @return {number} */
function f(x) { return x == MyEnum.FOO && MyEnum.FOO == x; }
`).
		AddDiagnostic(`inconsistent return type
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckTest.testEnum41()
func (c *TypeCheckTest) TestEnum41() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum {number} */ var MyEnum = {/** @const */ FOO: 1};
/** @return {string} */
function f() { return MyEnum.FOO; }
`).
		AddDiagnostic(`inconsistent return type
found   : MyEnum<number>
required: string
`).
		Run()
}

// Java: TypeCheckTest.testEnum42a()
func (c *TypeCheckTest) TestEnum42a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {number} x */ function f(x) {}
/** @enum {Object} */ var MyEnum = {FOO: {a: 1, b: 2}};
f(MyEnum.FOO.a);
`).
		AddDiagnostic("Property a never defined on MyEnum<Object>").
		Run()
}

// Java: TypeCheckTest.testEnum42b()
func (c *TypeCheckTest) TestEnum42b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {number} x */ function f(x) {}
/** @enum {!Object} */ var MyEnum = {FOO: {a: 1, b: 2}};
f(MyEnum.FOO.a);
`).
		AddDiagnostic("Property a never defined on MyEnum<Object>").
		Run()
}

// Java: TypeCheckTest.testEnum42c()
func (c *TypeCheckTest) TestEnum42c() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @param {number} x */ function f(x) {}
/** @enum {Object} */ var MyEnum = {FOO: {a: 1, b: 2}};
f(MyEnum.FOO.a);
`).
		Run()
}

// Java: TypeCheckTest.testEnum43()
func (c *TypeCheckTest) TestEnum43() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {number} x */ function f(x) {}
/** @enum {{a:number, b:number}} */ var MyEnum = {FOO: {a: 1, b: 2}};
f(MyEnum.FOO.a);
`).
		Run()
}

// Java: TypeCheckTest.testEnumDefinedInObjectLiteral()
func (c *TypeCheckTest) TestEnumDefinedInObjectLiteral() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var ns = {
  /** @enum {number} */
  Enum: {A: 1, B: 2},
};
/** @param {!ns.Enum} arg */
function f(arg) {}
`).
		Run()
}

// Java: TypeCheckTest.testAliasedEnum1()
func (c *TypeCheckTest) TestAliasedEnum1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum */ var YourEnum = {FOO: 3};
/** @enum */ var MyEnum = YourEnum;
/** @param {MyEnum} x */ function f(x) {} f(MyEnum.FOO);
`).
		Run()
}

// Java: TypeCheckTest.testAliasedEnum2()
func (c *TypeCheckTest) TestAliasedEnum2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum */ var YourEnum = {FOO: 3};
/** @enum */ var MyEnum = YourEnum;
/** @param {YourEnum} x */ function f(x) {} f(MyEnum.FOO);
`).
		Run()
}

// Java: TypeCheckTest.testAliasedEnum3()
func (c *TypeCheckTest) TestAliasedEnum3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum */ var YourEnum = {FOO: 3};
/** @enum */ var MyEnum = YourEnum;
/** @param {MyEnum} x */ function f(x) {} f(YourEnum.FOO);
`).
		Run()
}

// Java: TypeCheckTest.testAliasedEnum4()
func (c *TypeCheckTest) TestAliasedEnum4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum */ var YourEnum = {FOO: 3};
/** @enum */ var MyEnum = YourEnum;
/** @param {YourEnum} x */ function f(x) {} f(YourEnum.FOO);
`).
		Run()
}

// Java: TypeCheckTest.testAliasedEnum5()
func (c *TypeCheckTest) TestAliasedEnum5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum */ var YourEnum = {FOO: 3};
/** @enum */ var MyEnum = YourEnum;
/** @param {string} x */ function f(x) {} f(MyEnum.FOO);
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : YourEnum<number>
required: string
`).
		Run()
}

// Java: TypeCheckTest.testAliasedEnum_rhsIsStubDeclaration()
func (c *TypeCheckTest) TestAliasedEnum_rhsIsStubDeclaration() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`let YourEnum;
/** @enum */ const MyEnum = YourEnum;
`).
		AddDiagnostic(`initializing variable
found   : undefined
required: enum{MyEnum}
`).
		Run()
}

// Java: TypeCheckTest.testAliasedEnum_rhsIsNonEnum()
func (c *TypeCheckTest) TestAliasedEnum_rhsIsNonEnum() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`let YourEnum = 0;
/** @enum */ const MyEnum = YourEnum;
`).
		AddDiagnostic(`initializing variable
found   : number
required: enum{MyEnum}
`).
		Run()
}

// Java: TypeCheckTest.testConstAliasedEnum()
func (c *TypeCheckTest) TestConstAliasedEnum() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum */ var YourEnum = {FOO: 3};
const MyEnum = YourEnum;
/** @param {MyEnum} x */ function f(x) {} f(MyEnum.FOO);
`).
		Run()
}

// Java: TypeCheckTest.testConstAliasedEnum_constJSDoc()
func (c *TypeCheckTest) TestConstAliasedEnum_constJSDoc() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum */ var YourEnum = {FOO: 3};
/** @const */ var MyEnum = YourEnum;
/** @param {MyEnum} x */ function f(x) {} f(MyEnum.FOO);
`).
		Run()
}

// Java: TypeCheckTest.testConstAliasedEnum_qnameConstJSDoc()
func (c *TypeCheckTest) TestConstAliasedEnum_qnameConstJSDoc() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum */ var YourEnum = {FOO: 3};
const ns = {};
/** @const */
ns.MyEnum = YourEnum;
/** @param {ns.MyEnum} x */ function f(x) {} f(ns.MyEnum.FOO);
`).
		Run()
}

// Java: TypeCheckTest.testConstAliasedEnum_inObjectLit()
func (c *TypeCheckTest) TestConstAliasedEnum_inObjectLit() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum */ var YourEnum = {FOO: 3};
const ns = {/** @const */ MyEnum: YourEnum}
/** @param {ns.MyEnum} x */ function f(x) {} f(ns.MyEnum.FOO);
`).
		Run()
}

// Java: TypeCheckTest.testConstAliasedEnum_nestedInObjectLit()
func (c *TypeCheckTest) TestConstAliasedEnum_nestedInObjectLit() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum */ var YourEnum = {FOO: 3};
const ns = {/** @const */ x: {/** @const */ MyEnum: YourEnum}}
/** @param {ns.x.MyEnum} x */ function f(x) {} f(ns.x.MyEnum.FOO);
`).
		Run()
}

// Java: TypeCheckTest.testBackwardsEnumUse1()
func (c *TypeCheckTest) TestBackwardsEnumUse1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {string} */ function f() { return MyEnum.FOO; }
/** @enum {string} */ var MyEnum = {FOO: 'x'};
`).
		Run()
}

// Java: TypeCheckTest.testBackwardsEnumUse2()
func (c *TypeCheckTest) TestBackwardsEnumUse2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {number} */ function f() { return MyEnum.FOO; }
/** @enum {string} */ var MyEnum = {FOO: 'x'};
`).
		AddDiagnostic(`inconsistent return type
found   : MyEnum<string>
required: number
`).
		Run()
}

// Java: TypeCheckTest.testBackwardsEnumUse3()
func (c *TypeCheckTest) TestBackwardsEnumUse3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {string} */ function f() { return MyEnum.FOO; }
/** @enum {string} */ var YourEnum = {FOO: 'x'};
/** @enum {string} */ var MyEnum = YourEnum;
`).
		Run()
}

// Java: TypeCheckTest.testBackwardsEnumUse4()
func (c *TypeCheckTest) TestBackwardsEnumUse4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {number} */ function f() { return MyEnum.FOO; }
/** @enum {string} */ var YourEnum = {FOO: 'x'};
/** @enum {string} */ var MyEnum = YourEnum;
`).
		AddDiagnostic(`inconsistent return type
found   : YourEnum<string>
required: number
`).
		Run()
}

// Java: TypeCheckTest.testBackwardsEnumUse5()
func (c *TypeCheckTest) TestBackwardsEnumUse5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {string} */ function f() { return MyEnum.BAR; }
/** @enum {string} */ var YourEnum = {FOO: 'x'};
/** @enum {string} */ var MyEnum = YourEnum;
`).
		AddDiagnostic("element BAR does not exist on this enum").
		Run()
}

// Java: TypeCheckTest.testBackwardsTypedefUse2()
func (c *TypeCheckTest) TestBackwardsTypedefUse2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @this {MyTypedef} */ function f() {}
/** @typedef {!(Date|Array)} */ var MyTypedef;
`).
		Run()
}

// Java: TypeCheckTest.testBackwardsTypedefUse4()
func (c *TypeCheckTest) TestBackwardsTypedefUse4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {MyTypedef} */ function f() { return null; }
/** @typedef {string} */ var MyTypedef;
`).
		AddDiagnostic(`inconsistent return type
found   : null
required: string
`).
		Run()
}

// Java: TypeCheckTest.testBackwardsTypedefUse6()
func (c *TypeCheckTest) TestBackwardsTypedefUse6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {goog.MyTypedef} */ function f() { return null; }
var goog = {};
/** @typedef {string} */ goog.MyTypedef;
`).
		AddDiagnostic(`inconsistent return type
found   : null
required: string
`).
		Run()
}

// Java: TypeCheckTest.testBackwardsTypedefUse7()
func (c *TypeCheckTest) TestBackwardsTypedefUse7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {goog.MyTypedef} */ function f() { return null; }
var goog = {};
/** @typedef {Object} */ goog.MyTypedef;
`).
		Run()
}

// Java: TypeCheckTest.testBackwardsTypedefUse8()
func (c *TypeCheckTest) TestBackwardsTypedefUse8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {!Array} x */ function g(x) {}
/** @this {goog.MyTypedef} */ function f() { g(this); }
var goog = {};
/** @typedef {(Array|null|undefined)} */ goog.MyTypedef;
`).
		Run()
}

// Java: TypeCheckTest.testBackwardsTypedefUse9()
func (c *TypeCheckTest) TestBackwardsTypedefUse9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {!Array} x */ function g(x) {}
/** @this {goog.MyTypedef} */ function f() { g(this); }
var goog = {};
/** @typedef {(RegExp|null|undefined)} */ goog.MyTypedef;
`).
		AddDiagnostic(`actual parameter 1 of g does not match formal parameter
found   : RegExp
required: Array
`).
		Run()
}

// Java: TypeCheckTest.testBackwardsTypedefUse10()
func (c *TypeCheckTest) TestBackwardsTypedefUse10() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {goog.MyEnum} x */ function g(x) {}
var goog = {};
/** @enum {goog.MyTypedef} */ goog.MyEnum = {FOO: 1};
/** @typedef {number} */ goog.MyTypedef;
g(1);
`).
		AddDiagnostic(`actual parameter 1 of g does not match formal parameter
found   : number
required: goog.MyEnum<number>
`).
		Run()
}

// Java: TypeCheckTest.testBackwardsConstructor1()
func (c *TypeCheckTest) TestBackwardsConstructor1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f() { (new Foo(true)); }
/**
 * @constructor
 * @param {number} x */
var Foo = function(x) {};
`).
		AddDiagnostic(`actual parameter 1 of Foo does not match formal parameter
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckTest.testBackwardsConstructor2()
func (c *TypeCheckTest) TestBackwardsConstructor2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f() { (new Foo(true)); }
/**
 * @constructor
 * @param {number} x */
var YourFoo = function(x) {};
/**
 * @constructor
 * @param {number} x */
var Foo = YourFoo;
`).
		AddDiagnostic(`actual parameter 1 of Foo does not match formal parameter
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckTest.testMinimalConstructorAnnotation()
func (c *TypeCheckTest) TestMinimalConstructorAnnotation() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @constructor */function Foo(){}").
		Run()
}

// Java: TypeCheckTest.testGoodExtends1()
func (c *TypeCheckTest) TestGoodExtends1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */function base() {}
/** @constructor
 * @extends {base} */function derived() {}
`).
		Run()
}

// Java: TypeCheckTest.testGoodExtends2()
func (c *TypeCheckTest) TestGoodExtends2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor
 * @extends base */function derived() {}
/** @constructor */function base() {}
`).
		Run()
}

// Java: TypeCheckTest.testGoodExtends3()
func (c *TypeCheckTest) TestGoodExtends3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor
 * @extends {Object} */function base() {}
/** @constructor
 * @extends {base} */function derived() {}
`).
		Run()
}

// Java: TypeCheckTest.testGoodExtends4()
func (c *TypeCheckTest) TestGoodExtends4() {
	var n *rhino.Node = c.ParseAndTypeCheck(`var goog = {};
/** @constructor */goog.Base = function(){};
/** @constructor
  * @extends {goog.Base} */goog.Derived = function(){};
`)
	var subTypeName *rhino.Node = n.GetLastChild().GetLastChild().
		GetFirstChild()
	AssertThat(subTypeName.GetQualifiedName()).IsEqualTo("goog.Derived")
	var subCtorType rhino.FunctionType = subTypeName.GetNext().GetJSType().(rhino.FunctionType)
	AssertThat(subCtorType.GetInstanceType().String()).IsEqualTo("goog.Derived")
	var superType rhino.JSType = subCtorType.GetPrototype().GetImplicitPrototype()
	AssertThat(superType.String()).IsEqualTo("goog.Base")
}

// Java: TypeCheckTest.testGoodExtends5()
func (c *TypeCheckTest) TestGoodExtends5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */function base() {}
/** @extends {base}
 * @constructor */function derived() {}
`).
		Run()
}

// Java: TypeCheckTest.testGoodExtends6()
func (c *TypeCheckTest) TestGoodExtends6() {
	c.TestFunctionTypeWithFunctionType(CompilerTypeTestCase_CLOSURE_DEFS+`/** @constructor */function base() {}
/** @return {number} */
  base.prototype.foo = function() { return 1; };
/** @extends {base}
 * @constructor */function derived() {}
goog.inherits(derived, base);
`, "derived.superClass_.foo", "function(this:base): number")
}

// Java: TypeCheckTest.testGoodExtends7()
func (c *TypeCheckTest) TestGoodExtends7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor
 * @extends {Base} */ function Sub() {}
/** @return {number} */ function f() { return (new Sub()).foo; }
/** @constructor */ function Base() {}
/** @type {boolean} */ Base.prototype.foo = true;
`).
		AddDiagnostic(`inconsistent return type
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckTest.testGoodExtends8()
func (c *TypeCheckTest) TestGoodExtends8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Super() {}
Super.prototype.foo = function() {};
/** @constructor
 * @extends {Super} */ function Sub() {}
Sub.prototype = new Super();
/** @override */ Sub.prototype.foo = function() {};
`).
		Run()
}

// Java: TypeCheckTest.testGoodExtends9()
func (c *TypeCheckTest) TestGoodExtends9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Super() {}
/** @constructor
 * @extends {Super} */ function Sub() {}
Sub.prototype = new Super();
/** @return {Super} */ function foo() { return new Sub(); }
`).
		Run()
}

// Java: TypeCheckTest.testGoodExtends10()
func (c *TypeCheckTest) TestGoodExtends10() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Super() {}
/** @param {boolean} x */ Super.prototype.foo = function(x) {};
/** @constructor
 * @extends {Super} */ function Sub() {}
Sub.prototype = new Super();
(new Sub()).foo(0);
`).
		AddDiagnostic(`actual parameter 1 of Super.prototype.foo does not match formal parameter
found   : number
required: boolean
`).
		Run()
}

// Java: TypeCheckTest.testGoodExtends11()
func (c *TypeCheckTest) TestGoodExtends11() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor
 * @extends {Super} */ function Sub() {}
/** @constructor
 * @extends {Sub} */ function Sub2() {}
/** @constructor */ function Super() {}
/** @param {Super} x */ function foo(x) {}
foo(new Sub2());
`).
		Run()
}

// Java: TypeCheckTest.testGoodExtends12()
func (c *TypeCheckTest) TestGoodExtends12() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor
 * @extends {B}  */ function C() {}
/** @constructor
 * @extends {D}  */ function E() {}
/** @constructor
 * @extends {C}  */ function D() {}
/** @constructor
 * @extends {A} */ function B() {}
/** @constructor */ function A() {}
/** @param {number} x */ function f(x) {} f(new E());
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : E
required: number
`).
		Run()
}

// Java: TypeCheckTest.testGoodExtends13()
func (c *TypeCheckTest) TestGoodExtends13() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(CompilerTypeTestCase_CLOSURE_DEFS + `/** @param {Function} f */ function g(f) {
  /** @constructor */ function NewType() {};
  goog.inherits(NewType, f);
  (new NewType());
}
`).
		Run()
}

// Java: TypeCheckTest.testGoodExtends14()
func (c *TypeCheckTest) TestGoodExtends14() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(CompilerTypeTestCase_CLOSURE_DEFS + `/** @constructor */ function OldType() {}
/** @param {?function(new:OldType)} f */ function g(f) {
  /**
    * @constructor
    * @extends {OldType}
    */
  function NewType() {};
  goog.inherits(NewType, f);
  NewType.prototype.method = function() {
    NewType.superClass_.foo.call(this);
  };
}
`).
		AddDiagnostic("Property foo never defined on OldType.prototype").
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testGoodExtends_withAliasOfSuperclass()
func (c *TypeCheckTest) TestGoodExtends_withAliasOfSuperclass() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().IncludeDefaultExterns().
		AddSource(CompilerTypeTestCase_CLOSURE_DEFS + `/** @constructor */ const OldType = function () {};
const OldTypeAlias = OldType;

/**
  * @constructor
  * @extends {OldTypeAlias}
  */
function NewType() {};
// Verify that we recognize the inheritance even when goog.inherits references
// OldTypeAlias, not OldType.
goog.inherits(NewType, OldTypeAlias);
NewType.prototype.method = function() {
  NewType.superClass_.foo.call(this);
};
`).
		AddDiagnostic("Property foo never defined on OldType.prototype").
		Run()
}

// Java: TypeCheckTest.testBadExtends_withAliasOfSuperclass()
func (c *TypeCheckTest) TestBadExtends_withAliasOfSuperclass() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(CompilerTypeTestCase_CLOSURE_DEFS + `const ns = {};
/** @constructor */ ns.OldType = function () {};
const nsAlias = ns;

/**
  * @constructor
  * // no @extends here, NewType should not have a goog.inherits
  */
function NewType() {};
// Verify that we recognize the inheritance even when goog.inherits references
// nsAlias.OldType, not ns.OldType.
goog.inherits(NewType, ns.OldType);
`).
		AddDiagnostic("Missing @extends tag on type NewType").
		Run()
}

// Java: TypeCheckTest.testBadExtends_withNamespacedAliasOfSuperclass()
func (c *TypeCheckTest) TestBadExtends_withNamespacedAliasOfSuperclass() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(CompilerTypeTestCase_CLOSURE_DEFS + `const ns = {};
/** @constructor */ ns.OldType = function () {};
const nsAlias = ns;

/**
  * @constructor
  * // no @extends here, NewType should not have a goog.inhertis
  */
function NewType() {};
// Verify that we recognize the inheritance even when goog.inherits references
// nsAlias.OldType, not ns.OldType.
goog.inherits(NewType, nsAlias.OldType);
`).
		AddDiagnostic("Missing @extends tag on type NewType").
		Run()
}

// Java: TypeCheckTest.testBadExtends_withUnionType()
func (c *TypeCheckTest) TestBadExtends_withUnionType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
class Foo {}
/** @interface */
class Bar {}
/** @typedef {!Foo|!Bar} */
let Baz;
/**
 * @interface
 * @extends {Baz}
 */
class Blah {}
`).
		AddDiagnostic("Blah @extends non-object type (Bar|Foo)").
		Run()
}

// Java: TypeCheckTest.testGoodExtends_withNamespacedAliasOfSuperclass()
func (c *TypeCheckTest) TestGoodExtends_withNamespacedAliasOfSuperclass() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(CompilerTypeTestCase_CLOSURE_DEFS + `const ns = {};
/** @constructor */ ns.OldType = function () {};
const nsAlias = ns;

/**
  * @constructor
  * @extends {nsAlias.OldType}
  */
function NewType() {};
// Verify that we recognize the inheritance even when goog.inherits references
// nsAlias.OldType, not ns.OldType.
goog.inherits(NewType, nsAlias.OldType);
NewType.prototype.method = function() {
  NewType.superClass_.foo.call(this);
};
`).
		IncludeDefaultExterns().
		AddDiagnostic("Property foo never defined on ns.OldType.prototype").
		Run()
}

// Java: TypeCheckTest.testGoodExtends15()
func (c *TypeCheckTest) TestGoodExtends15() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(CompilerTypeTestCase_CLOSURE_DEFS + `/** @param {Function} f */ function g(f) {
  /** @constructor */ function NewType() {};
  goog.inherits(f, NewType);
  (new NewType());
}
`).
		Run()
}

// Java: TypeCheckTest.testGoodExtends16()
func (c *TypeCheckTest) TestGoodExtends16() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(CompilerTypeTestCase_CLOSURE_DEFS + `/** @constructor
 * @template T */
function C() {}
/** @constructor
 * @extends {C<string>} */
function D() {};
goog.inherits(D, C);
(new D())
`).
		Run()
}

// Java: TypeCheckTest.testGoodExtends17()
func (c *TypeCheckTest) TestGoodExtends17() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(CompilerTypeTestCase_CLOSURE_DEFS + `/** @constructor */
function C() {}

/** @interface
 * @template T */
function D() {}
/** @param {T} t */
D.prototype.method;

/** @constructor
 * @template T
 * @extends {C}
 * @implements {D<T>} */
function E() {};
goog.inherits(E, C);
/** @override */
E.prototype.method = function(t) {};

var e = /** @type {E<string>} */ (new E());
e.method(3);
`).
		AddDiagnostic(`actual parameter 1 of E.prototype.method does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testGoodExtends18()
func (c *TypeCheckTest) TestGoodExtends18() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
var MyInterface = function() {};
MyInterface.prototype = {
  /** @return {number} */
  method: function() {}
}
/** @extends {MyInterface}
 * @interface */
var MyOtherInterface = function() {};
MyOtherInterface.prototype = {
  /** @return {number}
   * @override */
  method: function() {}
}
`).
		Run()
}

// Java: TypeCheckTest.testGoodExtends19()
func (c *TypeCheckTest) TestGoodExtends19() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */
var MyType = function() {};
MyType.prototype = {
  /** @return {number} */
  method: function() {}
}
/** @constructor
 * @extends {MyType}
 */
var MyOtherType = function() {};
MyOtherType.prototype = {
  /** @return {number}
   * @override */
  method: function() {}
}
`).
		Run()
}

// Java: TypeCheckTest.testBadExtends1()
func (c *TypeCheckTest) TestBadExtends1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */function base() {}
/** @constructor
 * @extends {not_base} */function derived() {}
`).
		AddDiagnostic("Bad type annotation. Unknown type not_base").
		Run()
}

// Java: TypeCheckTest.testBadExtends2()
func (c *TypeCheckTest) TestBadExtends2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */function base() {
/** @type {!Number}*/
this.baseMember = new Number(4);
}
/** @constructor
  * @extends {base} */function derived() {}
/** @param {!String} x*/
function foo(x){ }
/** @type {!derived}*/var y;
foo(y.baseMember);
`).
		AddDiagnostic(`actual parameter 1 of foo does not match formal parameter
found   : Number
required: String
`).
		Run()
}

// Java: TypeCheckTest.testBadExtends3()
func (c *TypeCheckTest) TestBadExtends3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @extends {Object} */function base() {}").
		AddDiagnostic("@extends used without @constructor or @interface for base").
		Run()
}

// Java: TypeCheckTest.testBadExtends4()
func (c *TypeCheckTest) TestBadExtends4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor
 * @extends {bad} */ function Sub() {}
/** @constructor
 * @extends {Sub} */ function Sub2() {}
/** @param {Sub} x */ function foo(x) {}
foo(new Sub2());
`).
		AddDiagnostic("Bad type annotation. Unknown type bad").
		Run()
}

// Java: TypeCheckTest.testBadExtends5()
func (c *TypeCheckTest) TestBadExtends5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
var MyInterface = function() {};
MyInterface.prototype = {
  /** @return {number} */
  method: function() {}
}
/** @extends {MyInterface}
 * @interface */
var MyOtherInterface = function() {};
MyOtherInterface.prototype = {
  /** @return {string}
 @override */
  method: function() {}
}
`).
		AddDiagnostic(`mismatch of the method property on type MyOtherInterface and the type of the property it overrides from interface MyInterface
original: function(this:MyInterface): number
override: function(this:MyOtherInterface): string
`).
		Run()
}

// Java: TypeCheckTest.testBadExtends6()
func (c *TypeCheckTest) TestBadExtends6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */
var MyType = function() {};
MyType.prototype = {
  /** @return {number} */
  method: function() {}
}
/** @constructor
 * @extends {MyType}
 */
var MyOtherType = function() {};
MyOtherType.prototype = {
  /** @return {string}
   * @override */
  method: function() { return ''; }
}
`).
		AddDiagnostic(`mismatch of the method property type and the type of the property it overrides from superclass MyType
original: function(this:MyType): number
override: function(this:MyOtherType): string
`).
		Run()
}

// Java: TypeCheckTest.testLateExtends()
func (c *TypeCheckTest) TestLateExtends() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(CompilerTypeTestCase_CLOSURE_DEFS + `/** @constructor */ function Foo() {}
Foo.prototype.foo = function() {};
/** @constructor */function Bar() {}
goog.inherits(Foo, Bar);
`).
		AddDiagnostic("Missing @extends tag on type Foo").
		Run()
}

// Java: TypeCheckTest.testSuperclassMatch()
func (c *TypeCheckTest) TestSuperclassMatch() {
	c.compiler.GetOptions().SetCodingConvention(NewGoogleCodingConvention())
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ var Foo = function() {};
/** @constructor
 * @extends Foo */ var Bar = function() {};
Bar.inherits = function(x){};
Bar.inherits(Foo);
`).
		Run()
}

// Java: TypeCheckTest.testSuperclassMatchWithMixin()
func (c *TypeCheckTest) TestSuperclassMatchWithMixin() {
	c.compiler.GetOptions().SetCodingConvention(NewGoogleCodingConvention())
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ var Foo = function() {};
/** @constructor */ var Baz = function() {};
/** @constructor
 * @extends Foo */ var Bar = function() {};
Bar.inherits = function(x){};
Bar.mixin = function(y){};
Bar.inherits(Foo);
Bar.mixin(Baz);
`).
		Run()
}

// Java: TypeCheckTest.testSuperClassDefinedAfterSubClass1()
func (c *TypeCheckTest) TestSuperClassDefinedAfterSubClass1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor
 * @extends {Base} */ function A() {}
/** @constructor
 * @extends {Base} */ function B() {}
/** @constructor */ function Base() {}
/** @param {A|B} x
 * @return {B|A} */
function foo(x) { return x; }
`).
		Run()
}

// Java: TypeCheckTest.testSuperClassDefinedAfterSubClass2()
func (c *TypeCheckTest) TestSuperClassDefinedAfterSubClass2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor
 * @extends {Base} */ function A() {}
/** @constructor
 * @extends {Base} */ function B() {}
/** @param {A|B} x
 * @return {B|A} */
function foo(x) { return x; }
/** @constructor */ function Base() {}
`).
		Run()
}

// Java: TypeCheckTest.testGoodSuperCall()
func (c *TypeCheckTest) TestGoodSuperCall() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class A {
  /**
   * @param {string} a
   */
  constructor(a) {
    this.a = a;
  }
}
class B extends A {
  constructor() {
    super('b');
  }
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testBadSuperCall()
func (c *TypeCheckTest) TestBadSuperCall() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class A {
  /**
   * @param {string} a
   */
  constructor(a) {
    this.a = a;
  }
}
class B extends A {
  constructor() {
    super(5);
  }
}
`).
		AddDiagnostic(`actual parameter 1 of super does not match formal parameter
found   : number
required: string
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testDirectPrototypeAssignment1()
func (c *TypeCheckTest) TestDirectPrototypeAssignment1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Base() {}
Base.prototype.foo = 3;
/** @constructor
 * @extends {Base} */ function A() {}
A.prototype = new Base();
/** @return {string} */ function foo() { return (new A).foo; }
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testDirectPrototypeAssignment2()
func (c *TypeCheckTest) TestDirectPrototypeAssignment2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Base() {}
/** @constructor
 * @extends {Base} */ function A() {}
A.prototype = new Base();
A.prototype.foo = 3;
/** @return {string} */ function foo() { return (new Base).foo; }
`).
		AddDiagnosticType(TypeCheck_STRICT_INEXISTENT_PROPERTY).
		Run()
}

// Java: TypeCheckTest.testDirectPrototypeAssignment3()
func (c *TypeCheckTest) TestDirectPrototypeAssignment3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ var MainWidgetCreator = function() {};
/** @param {Function} ctor */
function createMainWidget(ctor) {
  /** @constructor */ function tempCtor() {};
  tempCtor.prototype = ctor.prototype;
  MainWidgetCreator.superClass_ = ctor.prototype;
  MainWidgetCreator.prototype = new tempCtor();
}
`).
		Run()
}

// Java: TypeCheckTest.testGoodImplements1()
func (c *TypeCheckTest) TestGoodImplements1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */function Disposable() {}
/** @implements {Disposable}
 * @constructor */function f() {}
`).
		Run()
}

// Java: TypeCheckTest.testGoodImplements2()
func (c *TypeCheckTest) TestGoodImplements2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */function Base1() {}
/** @interface */function Base2() {}
/** @constructor
 * @implements {Base1}
 * @implements {Base2}
 */ function derived() {}
`).
		Run()
}

// Java: TypeCheckTest.testGoodImplements3()
func (c *TypeCheckTest) TestGoodImplements3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */function Disposable() {}
/** @constructor
 * @implements {Disposable} */function f() {}
`).
		Run()
}

// Java: TypeCheckTest.testGoodImplements4()
func (c *TypeCheckTest) TestGoodImplements4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = {};
/** @type {!Function} */
goog.abstractMethod = function() {};
/** @interface */
goog.Disposable = function() {};
goog.Disposable.prototype.dispose = goog.abstractMethod;
/** @implements {goog.Disposable}
 * @constructor */
goog.SubDisposable = function() {};
/** @inheritDoc */
goog.SubDisposable.prototype.dispose = function() {};
`).
		Run()
}

// Java: TypeCheckTest.testGoodImplements5()
func (c *TypeCheckTest) TestGoodImplements5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
goog.Disposable = function() {};
/** @type {Function} */
goog.Disposable.prototype.dispose = function() {};
/** @implements {goog.Disposable}
 * @constructor */
goog.SubDisposable = function() {};
/** @param {number} key
 * @override */
goog.SubDisposable.prototype.dispose = function(key) {};
`).
		Run()
}

// Java: TypeCheckTest.testGoodImplements6()
func (c *TypeCheckTest) TestGoodImplements6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var myNullFunction = function() {};
/** @interface */
goog.Disposable = function() {};
/** @return {number} */
goog.Disposable.prototype.dispose = myNullFunction;
/** @implements {goog.Disposable}
 * @constructor */
goog.SubDisposable = function() {};
/** @return {number}
 * @override */
goog.SubDisposable.prototype.dispose = function() { return 0; };
`).
		Run()
}

// Java: TypeCheckTest.testGoodImplements7()
func (c *TypeCheckTest) TestGoodImplements7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var myNullFunction = function() {};
/** @interface */
goog.Disposable = function() {};
/** @return {number} */
goog.Disposable.prototype.dispose = function() {};
/** @implements {goog.Disposable}
 * @constructor */
goog.SubDisposable = function() {};
/** @return {number}
 * @override */
goog.SubDisposable.prototype.dispose = function() { return 0; };
`).
		Run()
}

// Java: TypeCheckTest.testGoodImplements8()
func (c *TypeCheckTest) TestGoodImplements8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
MyInterface = function() {};
MyInterface.prototype = {
  /** @return {number} */
  method: function() {}
}
/** @implements {MyInterface}
 * @constructor */
MyClass = function() {};
MyClass.prototype = {
  /** @return {number}
   * @override */
  method: function() { return 0; }
}
`).
		Run()
}

// Java: TypeCheckTest.testBadImplements1()
func (c *TypeCheckTest) TestBadImplements1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */function Base1() {}
/** @interface */function Base2() {}
/** @constructor
 * @implements {nonExistent}
 * @implements {Base2}
 */ function derived() {}
`).
		AddDiagnostic("Bad type annotation. Unknown type nonExistent").
		Run()
}

// Java: TypeCheckTest.testBadImplements2()
func (c *TypeCheckTest) TestBadImplements2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */function Disposable() {}
/** @implements {Disposable}
 */function f() {}
`).
		AddDiagnostic("@implements used without @constructor for f").
		Run()
}

// Java: TypeCheckTest.testBadImplements3()
func (c *TypeCheckTest) TestBadImplements3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = {};
/** @type {!Function} */ goog.abstractMethod = function(){};
/** @interface */ var Disposable = function() {};
Disposable.prototype.method = goog.abstractMethod;
/** @implements {Disposable}
 * @constructor */function f() {}
`).
		AddDiagnostic("property method on interface Disposable is not implemented by type f").
		Run()
}

// Java: TypeCheckTest.testBadImplements4()
func (c *TypeCheckTest) TestBadImplements4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */function Disposable() {}
/** @implements {Disposable}
 * @interface */function f() {}
`).
		AddDiagnostic(`f cannot implement this type; an interface can only extend, but not implement interfaces
`).
		Run()
}

// Java: TypeCheckTest.testBadImplements5()
func (c *TypeCheckTest) TestBadImplements5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */function Disposable() {}
/** @type {number} */ Disposable.prototype.bar = function() {};
`).
		AddDiagnostic(`assignment to property bar of Disposable.prototype
found   : function(): undefined
required: number
`).
		Run()
}

// Java: TypeCheckTest.testBadImplements6()
func (c *TypeCheckTest) TestBadImplements6() {
	c.TestClosureTypesMultipleWarnings(`/** @interface */function Disposable() {}
/** @type {function()} */ Disposable.prototype.bar = 3;
`, []string{`assignment to property bar of Disposable.prototype
found   : number
required: function(): ?`, "interface members can only be empty property declarations, empty functions, or goog.abstractMethod"})
}

// Java: TypeCheckTest.testBadImplements7()
func (c *TypeCheckTest) TestBadImplements7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
MyInterface = function() {};
MyInterface.prototype = {
  /** @return {number} */
  method: function() {}
}
/** @implements {MyInterface}
 * @constructor */
MyClass = function() {};
MyClass.prototype = {
  /** @return {string}
   * @override */
  method: function() { return ''; }
}
`).
		AddDiagnostic(`mismatch of the method property on type MyClass and the type of the property it overrides from interface MyInterface
original: function(): number
override: function(): string
`).
		Run()
}

// Java: TypeCheckTest.testBadImplements8()
func (c *TypeCheckTest) TestBadImplements8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
MyInterface = function() {};
MyInterface.prototype = {
  /** @return {number} */
  method: function() {}
}
/** @implements {MyInterface}
 * @constructor */
MyClass = function() {};
MyClass.prototype = {
  /** @return {number} */
  method: function() { return 0; }
}
`).
		AddDiagnostic(`property method already defined on interface MyInterface; use @override to override it
`).
		Run()
}

// Java: TypeCheckTest.testProtoDoesNotRequireOverrideFromInterface()
func (c *TypeCheckTest) TestProtoDoesNotRequireOverrideFromInterface() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().IncludeDefaultExterns().
		AddExterns("/** @type {Object} */ Object.prototype.__proto__;").
		AddSource(`/** @interface */
var MyInterface = function() {};
/** @constructor
 @implements {MyInterface} */
var MySuper = function() {};
/** @constructor
 @extends {MySuper} */
var MyClass = function() {};
MyClass.prototype = {
  __proto__: MySuper.prototype
}
`).
		Run()
}

// Java: TypeCheckTest.testConstructorClassTemplate()
func (c *TypeCheckTest) TestConstructorClassTemplate() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @template S,T
 */
function A() {}
`).
		Run()
}

// Java: TypeCheckTest.testGenericBoundExplicitUnknown()
func (c *TypeCheckTest) TestGenericBoundExplicitUnknown() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_BOUNDED_GENERICS).
		AddSource(`/**
 * @param {T} x
 * @template {?} T
 */
function f(x) {}
`).
		AddDiagnostic("Illegal upper bound '?' on template type parameter T").
		DiagnosticsAreErrors().
		Run()
}

// Java: TypeCheckTest.testGenericBoundArgAppError()
func (c *TypeCheckTest) TestGenericBoundArgAppError() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {T} x
 * @template {number} T
 */
function f(x) {}
var a = f('a');
`).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testGenericBoundArgApp()
func (c *TypeCheckTest) TestGenericBoundArgApp() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {T} x
 * @template {number} T
 */
function f(x) {}
var a = f(3);
`).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		Run()
}

// Java: TypeCheckTest.testGenericBoundReturnError()
func (c *TypeCheckTest) TestGenericBoundReturnError() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @return {T}
 * @template {number} T
 */
function f(x) { return 'a'; }
var a = f(0);
`).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		AddDiagnostic(`inconsistent return type
found   : string
required: T extends number
`).
		Run()
}

// Java: TypeCheckTest.testGenericBoundArgAppNullError()
func (c *TypeCheckTest) TestGenericBoundArgAppNullError() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {T} x
 * @template {number} T
 */
function f(x) { return x; }
var a = f(null);
`).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : null
required: number
`).
		Run()
}

// Java: TypeCheckTest.testGenericBoundArgAppNullable()
func (c *TypeCheckTest) TestGenericBoundArgAppNullable() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {T} x
 * @template {?number} T
 */
function f(x) { return x; }
var a = f(null);
`).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		Run()
}

// Java: TypeCheckTest.testGenericBoundArgInnerCallAppError()
func (c *TypeCheckTest) TestGenericBoundArgInnerCallAppError() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {string} x
 */
function stringID(x) { return x; }
/**
 * @param {T} x
 * @template {number|boolean} T
 */
function foo(x) { return stringID(x); }
`).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		AddDiagnostic(`actual parameter 1 of stringID does not match formal parameter
found   : T extends (boolean|number)
required: string
`).
		Run()
}

// Java: TypeCheckTest.testGenericBoundArgInnerCallApp()
func (c *TypeCheckTest) TestGenericBoundArgInnerCallApp() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {number} x
 */
function numID(x) { return x; }
/**
 * @param {T} x
 * @template {number} T
 */
function foo(x) { return numID(x); }
`).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		Run()
}

// Java: TypeCheckTest.testGenericBoundArgInnerCallAppSubtypeError()
func (c *TypeCheckTest) TestGenericBoundArgInnerCallAppSubtypeError() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {number} x
 */
function numID(x) { return x; }
/**
 * @param {T} x
 * @template {number|boolean|string} T
 */
function foo(x) { return numID(x); }
`).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		AddDiagnostic(`actual parameter 1 of numID does not match formal parameter
found   : T extends (boolean|number|string)
required: number
`).
		Run()
}

// Java: TypeCheckTest.testGenericBoundArgInnerAssignSubtype()
func (c *TypeCheckTest) TestGenericBoundArgInnerAssignSubtype() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {T} x
 * @param {number|string|boolean} y
 * @template {number|string} T
 */
function foo(x,y) { y=x; }
`).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		Run()
}

// Java: TypeCheckTest.testGenericBoundArgInnerAssignBoundedInvariantError()
func (c *TypeCheckTest) TestGenericBoundArgInnerAssignBoundedInvariantError() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {number} x
 * @param {T} y
 * @template {number|string} T
 */
function foo(x,y) { y=x; }
`).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		AddDiagnostic(`assignment
found   : number
required: T extends (number|string)
`).
		Run()
}

// Java: TypeCheckTest.testGenericBoundArgInnerAssignSubtypeError()
func (c *TypeCheckTest) TestGenericBoundArgInnerAssignSubtypeError() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {number} x
 * @param {T} y
 * @template {number|string} T
 */
function foo(x,y) { x=y; }
`).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		AddDiagnostic(`assignment
found   : T extends (number|string)
required: number
`).
		Run()
}

// Java: TypeCheckTest.testDoubleGenericBoundArgInnerAssignSubtype()
func (c *TypeCheckTest) TestDoubleGenericBoundArgInnerAssignSubtype() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {T} x
 * @param {T} y
 * @template {number|string} T
 */
function foo(x,y) { x=y; }
`).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		Run()
}

// Java: TypeCheckTest.testBoundedGenericForwardDeclaredParameterError()
func (c *TypeCheckTest) TestBoundedGenericForwardDeclaredParameterError() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @template {number} T
 */
class Foo {}
var /** !Foo<Str> */ a;
/** @typedef {string} */
let Str;
`).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		AddDiagnostic(`Bounded generic type error. string assigned to template type T is not a subtype of bound number
`).
		Run()
}

// Java: TypeCheckTest.testBoundedGenericParametrizedTypeVarError()
func (c *TypeCheckTest) TestBoundedGenericParametrizedTypeVarError() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @template {number|boolean} T
 */
function F() {}
var /** F<string> */ a;
`).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		AddDiagnostic(`Bounded generic type error. string assigned to template type T is not a subtype of bound (boolean|number)
`).
		Run()
}

// Java: TypeCheckTest.testBoundedGenericParametrizedTypeVar()
func (c *TypeCheckTest) TestBoundedGenericParametrizedTypeVar() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @param {T} x
 * @template {number|boolean} T
 */
function F(x) {}
var /** F<number> */ a;
`).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		Run()
}

// Java: TypeCheckTest.testDoubleGenericBoundArgInnerAssignSubtypeError()
func (c *TypeCheckTest) TestDoubleGenericBoundArgInnerAssignSubtypeError() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {T} x
 * @param {U} y
 * @template {number|string} T
 * @template {number|string} U
 */
function foo(x,y) { x=y; }
`).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		AddDiagnostic(`assignment
found   : U extends (number|string)
required: T extends (number|string)
`).
		Run()
}

// Java: TypeCheckTest.testGenericBoundBivariantTemplatized()
func (c *TypeCheckTest) TestGenericBoundBivariantTemplatized() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {!Array<T>} x
 * @param {!Array<number|boolean>} y
 * @template {number} T
 */
function foo(x,y) {
  var /** !Array<T> */ z = y;
  y = x;
}
`).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		Run()
}

// Java: TypeCheckTest.testConstructGenericBoundGenericBound()
func (c *TypeCheckTest) TestConstructGenericBoundGenericBound() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {T} x
 * @param {U} y
 * @template {boolean|string} T
 * @template {T} U
 */
function foo(x,y) { x=y; }
`).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		Run()
}

// Java: TypeCheckTest.testConstructGenericBoundGenericBoundError()
func (c *TypeCheckTest) TestConstructGenericBoundGenericBoundError() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {T} x
 * @param {U} y
 * @template {boolean|string} T
 * @template {T} U
 */
function foo(x,y) {
  var /** U */ z = x;
  x = y;
}
`).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		AddDiagnostic(`initializing variable
found   : T extends (boolean|string)
required: U extends T extends (boolean|string)
`).
		Run()
}

// Java: TypeCheckTest.testDoubleDeclareTemplateNameInFunctions()
func (c *TypeCheckTest) TestDoubleDeclareTemplateNameInFunctions() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {T} x
 * @template {boolean} T
 */
function foo(x) { return x; }
/**
 * @param {T} x
 * @template {string} T
 */
function bar(x) { return x; }
foo(true);
bar('Hi');
`).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		Run()
}

// Java: TypeCheckTest.testDoubleDeclareTemplateNameInFunctionsError()
func (c *TypeCheckTest) TestDoubleDeclareTemplateNameInFunctionsError() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {T} x
 * @template {boolean} T
 */
function foo(x) { return x; }
/**
 * @param {T} x
 * @template {string} T
 */
function bar(x) { return x; }
foo('Hi');
bar(true);
`).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		AddDiagnostic(`actual parameter 1 of foo does not match formal parameter
found   : string
required: boolean
`).
		AddDiagnostic(`actual parameter 1 of bar does not match formal parameter
found   : boolean
required: string
`).
		Run()
}

// Java: TypeCheckTest.testShadowTemplateNameInFunctionAndClass()
func (c *TypeCheckTest) TestShadowTemplateNameInFunctionAndClass() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {T} x
 * @template {boolean} T
 */
function foo(x) { return x; }
class Foo {
  /**
   * @param {T} x
   * @template {string} T
   */
  foo(x) { return x; }
}
var F = new Foo();
F.foo('Hi');
foo(true);
`).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		Run()
}

// Java: TypeCheckTest.testShadowTemplateNameInFunctionAndClassError()
func (c *TypeCheckTest) TestShadowTemplateNameInFunctionAndClassError() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {T} x
 * @template {boolean} T
 */
function foo(x) { return x; }
class Foo {
  /**
   * @param {T} x
   * @template {string} T
   */
  foo(x) { return x; }
}
var F = new Foo();
F.foo(true);
foo('Hi');
`).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		AddDiagnostic(`actual parameter 1 of Foo.prototype.foo does not match formal parameter
found   : boolean
required: string
`).
		AddDiagnostic(`actual parameter 1 of foo does not match formal parameter
found   : string
required: boolean
`).
		Run()
}

// Java: TypeCheckTest.testTemplateTypeBounds_passingBoundedTemplateType_toTemplatedFunction()
func (c *TypeCheckTest) TestTemplateTypeBounds_passingBoundedTemplateType_toTemplatedFunction() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @constructor */\nfunction Foo() { }\n\n/**\n * @template {!Foo} X\n * @param {X} x\n * @return {X}\n */\nfunction clone(x) {\n// The focus of this test is that `X` (already a template variable) is bound to `Y` at\n// this callsite. We confirm that by ensuring an `X` is returned from `cloneInternal`.\n  return cloneInternal(x);\n}\n\n/**\n * @template {!Foo} Y\n * @param {Y} x\n * @return {Y}\n */\nfunction cloneInternal(x) {\n  return x;\n}\n").
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		Run()
}

// Java: TypeCheckTest.testTemplateTypeBounds_onFreeFunctions_areAlwaysSpecialized()
func (c *TypeCheckTest) TestTemplateTypeBounds_onFreeFunctions_areAlwaysSpecialized() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */
function Foo() { }

/**
 * @template {!Foo} X
 * @param {X} x
 * @return {X}
 */
function identity(x) {
  return x;
}

var /** function(!Foo): !Foo */ y = identity;
`).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		Run()
}

// Java: TypeCheckTest.testFunctionTemplateShadowClassTemplate()
func (c *TypeCheckTest) TestFunctionTemplateShadowClassTemplate() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @template T
 */
class Foo {
  /**
   * @param {T} x
   * @template T
   */
  foo(x) { return x; }
}
const /** !Foo<string> */ f = new Foo();
f.foo(3);
`).
		AddDiagnostic(`actual parameter 1 of Foo.prototype.foo does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testFunctionTemplateShadowClassTemplateBounded()
func (c *TypeCheckTest) TestFunctionTemplateShadowClassTemplateBounded() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @template {string} T
 */
class Foo {
  /**
   * @param {T} x
   * @template {number} T
   */
  foo(x) { return x; }
}
const /** !Foo<string> */ f = new Foo();
f.foo(3);
`).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		AddDiagnostic(`actual parameter 1 of Foo.prototype.foo does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testCyclicGenericBoundGenericError()
func (c *TypeCheckTest) TestCyclicGenericBoundGenericError() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_BOUNDED_GENERICS).
		AddSource(`/**
 * @template {S} T
 * @template {T} U
 * @template {U} S
 */
class Foo { }
`).
		AddDiagnostic("Cycle detected in inheritance chain of type S").
		AddDiagnostic("Cycle detected in inheritance chain of type T").
		AddDiagnostic("Cycle detected in inheritance chain of type U").
		Run()
}

// Java: TypeCheckTest.testUnitCyclicGenericBoundGenericError()
func (c *TypeCheckTest) TestUnitCyclicGenericBoundGenericError() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_BOUNDED_GENERICS).
		AddSource(`/**
 * @template {T} T
 */
class Foo { }
`).
		AddDiagnostic("Cycle detected in inheritance chain of type T").
		Run()
}

// Java: TypeCheckTest.testSecondUnitCyclicGenericBoundGenericError()
func (c *TypeCheckTest) TestSecondUnitCyclicGenericBoundGenericError() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_BOUNDED_GENERICS).
		AddSource(`/**
 * @template {T} T
 * @template {T} U
 */
class Foo { }
`).
		AddDiagnostic("Cycle detected in inheritance chain of type T").
		AddDiagnostic("Cycle detected in inheritance chain of type U").
		Run()
}

// Java: TypeCheckTest.testConstructCyclicGenericBoundGenericBoundGraph()
func (c *TypeCheckTest) TestConstructCyclicGenericBoundGenericBoundGraph() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @template V, E
 */
class Vertex {}
/**
 * @template V, E
 */
class Edge {}
/**
 * @template {Vertex<V, E>} V
 * @template {Edge<V, E>} E
 */
class Graph {}
`).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		Run()
}

// Java: TypeCheckTest.testBoundedGenericParametrizedTypeReturnError()
func (c *TypeCheckTest) TestBoundedGenericParametrizedTypeReturnError() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @template T
 */
function C(x) {}
/**
 * @template {number|string} U
 * @param {C<boolean>} x
 * @return {C<U>}
 */
function f(x) { return x; }
`).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		AddDiagnostic(`inconsistent return type
found   : (C<boolean>|null)
required: (C<U extends (number|string)>|null)
`).
		Run()
}

// Java: TypeCheckTest.testBoundedGenericParametrizedTypeError()
func (c *TypeCheckTest) TestBoundedGenericParametrizedTypeError() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @param {T} x
 * @template T
 */
function C(x) {}
/**
 * @template {number|string} U
 * @param {C<U>} x
 */
function f(x) {}
var /** C<boolean> */ c;
f(c);
`).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : (C<boolean>|null)
required: (C<(number|string)>|null)
`).
		Run()
}

// Java: TypeCheckTest.testPartialNarrowingBoundedGenericParametrizedTypeError()
func (c *TypeCheckTest) TestPartialNarrowingBoundedGenericParametrizedTypeError() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @param {T} x
 * @template T
 */
function C(x) {}
/**
 * @template {number|string} U
 * @param {C<U>} x
 * @param {U} y
 */
function f(x,y) {}
var /** C<string> */ c;
f(c,false);
`).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : (C<string>|null)
required: (C<(number|string)>|null)
`).
		AddDiagnostic(`actual parameter 2 of f does not match formal parameter
found   : boolean
required: (number|string)
`).
		Run()
}

// Java: TypeCheckTest.testPartialNarrowingBoundedGenericParametrizedType()
func (c *TypeCheckTest) TestPartialNarrowingBoundedGenericParametrizedType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @param {T} x
 * @template T
 */
function C(x) {}
/**
 * @template {number|string} U
 * @param {C<U>} x
 * @param {U} y
 */
function f(x,y) {}
var /** C<number|string> */ c;
f(c,0);
`).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		Run()
}

// Java: TypeCheckTest.testUnmappedBoundedGenericTypeError()
func (c *TypeCheckTest) TestUnmappedBoundedGenericTypeError() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @template {number|string} T
 * @template {number|null} S
 */
class Foo {
  /**
   * @param {T} x
   */
  bar(x) { }
  /**
   * @param {S} x
   */
  baz(x) { }
}
var /** Foo<number> */ foo;
foo.baz(3);
foo.baz(null);
foo.baz('3');
`).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		AddDiagnostic(`actual parameter 1 of Foo.prototype.baz does not match formal parameter
found   : string
required: S extends (null|number)
`).
		Run()
}

// Java: TypeCheckTest.testFunctionBodyBoundedGenericError()
func (c *TypeCheckTest) TestFunctionBodyBoundedGenericError() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class C {}
/**
 * @template {C} T
 * @param {T} t
 */
function f(t) { t = new C(); }
`).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		AddDiagnostic(`assignment
found   : C
required: T extends (C|null)
`).
		Run()
}

// Java: TypeCheckTest.testFunctionBodyBoundedPropertyGenericError()
func (c *TypeCheckTest) TestFunctionBodyBoundedPropertyGenericError() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @template {number|string} T
 */
class Foo {
  constructor() {
    /** @type {T} */
    this.x;
  }
  m() {
  this.x = 0;
  }
}
function f(t) { t = new C(); }
`).
		AddDiagnostic(TypeCheckTest_BOUNDED_GENERICS_USE_MSG).
		AddDiagnostic(`assignment to property x of Foo
found   : number
required: T extends (number|string)
`).
		Run()
}

// Java: TypeCheckTest.testInterfaceExtends()
func (c *TypeCheckTest) TestInterfaceExtends() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */function A() {}
/** @interface 
 * @extends {A} */function B() {}
/** @constructor
 * @implements {B}
 */ function derived() {}
`).
		Run()
}

// Java: TypeCheckTest.testDontCrashOnDupPropDefinition()
func (c *TypeCheckTest) TestDontCrashOnDupPropDefinition() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const */
var ns = {};
/** @interface */
ns.I = function() {};
/** @interface */
ns.A = function() {};
/**
 * @constructor
 * @implements {ns.I}
 */
ns.A = function() {};
`).
		AddDiagnostic("variable ns.A redefined, original definition at testcode0:6").
		Run()
}

// Java: TypeCheckTest.testBadInterfaceExtends1()
func (c *TypeCheckTest) TestBadInterfaceExtends1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @interface
 * @extends {nonExistent}
 */
function A() {}
`).
		AddDiagnostic("Bad type annotation. Unknown type nonExistent").
		Run()
}

// Java: TypeCheckTest.testBadInterfaceExtendsNonExistentInterfaces()
func (c *TypeCheckTest) TestBadInterfaceExtendsNonExistentInterfaces() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface 
 * @extends {nonExistent1} 
 * @extends {nonExistent2} 
 */function A() {}
`).
		AddDiagnostic("Bad type annotation. Unknown type nonExistent1").
		AddDiagnostic("Bad type annotation. Unknown type nonExistent2").
		Run()
}

// Java: TypeCheckTest.testBadInterfaceExtends2()
func (c *TypeCheckTest) TestBadInterfaceExtends2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */function A() {}
/** @interface
 * @extends {A} */function B() {}
`).
		AddDiagnostic("B cannot extend this type; interfaces can only extend interfaces").
		Run()
}

// Java: TypeCheckTest.testBadInterfaceExtends3()
func (c *TypeCheckTest) TestBadInterfaceExtends3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */function A() {}
/** @constructor
 * @extends {A} */function B() {}
`).
		AddDiagnostic("B cannot extend this type; constructors can only extend constructors").
		Run()
}

// Java: TypeCheckTest.testBadInterfaceExtends4()
func (c *TypeCheckTest) TestBadInterfaceExtends4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */function A() {}
/** @constructor */function B() {}
B.prototype = A;
`).
		Run()
}

// Java: TypeCheckTest.testBadInterfaceExtends5()
func (c *TypeCheckTest) TestBadInterfaceExtends5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */function A() {}
/** @interface */function B() {}
B.prototype = A;
`).
		Run()
}

// Java: TypeCheckTest.testBadImplementsAConstructor()
func (c *TypeCheckTest) TestBadImplementsAConstructor() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */function A() {}
/** @constructor
 * @implements {A} */function B() {}
`).
		AddDiagnostic("can only implement interfaces").
		Run()
}

// Java: TypeCheckTest.testBadImplementsAConstructorWithSubclass()
func (c *TypeCheckTest) TestBadImplementsAConstructorWithSubclass() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */
function A() {}
/** @implements {A} */
class B {}
class C extends B {}
`).
		AddDiagnostic("can only implement interfaces").
		Run()
}

// Java: TypeCheckTest.testBadImplementsNonInterfaceType()
func (c *TypeCheckTest) TestBadImplementsNonInterfaceType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor
 * @implements {Boolean} */function B() {}
`).
		AddDiagnostic("can only implement interfaces").
		Run()
}

// Java: TypeCheckTest.testBadImplementsNonObjectType()
func (c *TypeCheckTest) TestBadImplementsNonObjectType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor
 * @implements {string} */function S() {}
`).
		AddDiagnostic("can only implement interfaces").
		Run()
}

// Java: TypeCheckTest.testBadImplementsDuplicateInterface1()
func (c *TypeCheckTest) TestBadImplementsDuplicateInterface1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface 
 * @template T
 */
function Foo() {}
/** @constructor 
 * @implements {Foo<?>}
 * @implements {Foo}
 */
function A() {}
`).
		AddDiagnostic(`Cannot @implement the same interface more than once
Repeated interface: Foo
`).
		Run()
}

// Java: TypeCheckTest.testBadImplementsDuplicateInterface2()
func (c *TypeCheckTest) TestBadImplementsDuplicateInterface2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface 
 * @template T
 */
function Foo() {}
/** @constructor 
 * @implements {Foo<string>}
 * @implements {Foo<number>}
 */
function A() {}
`).
		AddDiagnostic(`Cannot @implement the same interface more than once
Repeated interface: Foo
`).
		Run()
}

// Java: TypeCheckTest.testInterfaceAssignment1()
func (c *TypeCheckTest) TestInterfaceAssignment1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */var I = function() {};
/** @constructor
@implements {I} */var T = function() {};
var t = new T();
/** @type {!I} */var i = t;
`).
		Run()
}

// Java: TypeCheckTest.testInterfaceAssignment2()
func (c *TypeCheckTest) TestInterfaceAssignment2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */var I = function() {};
/** @constructor */var T = function() {};
var t = new T();
/** @type {!I} */var i = t;
`).
		AddDiagnostic(`initializing variable
found   : T
required: I
`).
		Run()
}

// Java: TypeCheckTest.testInterfaceAssignment3()
func (c *TypeCheckTest) TestInterfaceAssignment3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */var I = function() {};
/** @constructor
@implements {I} */var T = function() {};
var t = new T();
/** @type {I|number} */var i = t;
`).
		Run()
}

// Java: TypeCheckTest.testInterfaceAssignment4()
func (c *TypeCheckTest) TestInterfaceAssignment4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */var I1 = function() {};
/** @interface */var I2 = function() {};
/** @constructor
@implements {I1} */var T = function() {};
var t = new T();
/** @type {I1|I2} */var i = t;
`).
		Run()
}

// Java: TypeCheckTest.testInterfaceAssignment5()
func (c *TypeCheckTest) TestInterfaceAssignment5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */var I1 = function() {};
/** @interface */var I2 = function() {};
/** @constructor
 * @implements {I1}
 * @implements {I2}*/
var T = function() {};
var t = new T();
/** @type {I1} */var i1 = t;
/** @type {I2} */var i2 = t;
`).
		Run()
}

// Java: TypeCheckTest.testInterfaceAssignment6()
func (c *TypeCheckTest) TestInterfaceAssignment6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */var I1 = function() {};
/** @interface */var I2 = function() {};
/** @constructor
@implements {I1} */var T = function() {};
/** @type {!I1} */var i1 = new T();
/** @type {!I2} */var i2 = i1;
`).
		AddDiagnostic(`initializing variable
found   : I1
required: I2
`).
		Run()
}

// Java: TypeCheckTest.testInterfaceAssignment7()
func (c *TypeCheckTest) TestInterfaceAssignment7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */var I1 = function() {};
/** @interface
@extends {I1}*/var I2 = function() {};
/** @constructor
@implements {I2}*/var T = function() {};
var t = new T();
/** @type {I1} */var i1 = t;
/** @type {I2} */var i2 = t;
i1 = i2;
`).
		Run()
}

// Java: TypeCheckTest.testInterfaceAssignment8()
func (c *TypeCheckTest) TestInterfaceAssignment8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @interface */var I = function() {};
/** @type {I} */var i;
/** @type {Object} */var o = i;
new Object().prototype = i.prototype;
`).
		Run()
}

// Java: TypeCheckTest.testInterfaceAssignment9()
func (c *TypeCheckTest) TestInterfaceAssignment9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */var I = function() {};
/** @return {I?} */function f() { return null; }
/** @type {!I} */var i = f();
`).
		AddDiagnostic(`initializing variable
found   : (I|null)
required: I
`).
		Run()
}

// Java: TypeCheckTest.testInterfaceAssignment10()
func (c *TypeCheckTest) TestInterfaceAssignment10() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */var I1 = function() {};
/** @interface */var I2 = function() {};
/** @constructor
@implements {I2} */var T = function() {};
/** @return {!I1|!I2} */function f() { return new T(); }
/** @type {!I1} */var i1 = f();
`).
		AddDiagnostic(`initializing variable
found   : (I1|I2)
required: I1
`).
		Run()
}

// Java: TypeCheckTest.testInterfaceAssignment11()
func (c *TypeCheckTest) TestInterfaceAssignment11() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */var I1 = function() {};
/** @interface */var I2 = function() {};
/** @constructor */var T = function() {};
/** @return {!I1|!I2|!T} */function f() { return new T(); }
/** @type {!I1} */var i1 = f();
`).
		AddDiagnostic(`initializing variable
found   : (I1|I2|T)
required: I1
`).
		Run()
}

// Java: TypeCheckTest.testInterfaceAssignment12()
func (c *TypeCheckTest) TestInterfaceAssignment12() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */var I = function() {};
/** @constructor
@implements{I}*/var T1 = function() {};
/** @constructor
@extends {T1}*/var T2 = function() {};
/** @return {I} */function f() { return new T2(); }
`).
		Run()
}

// Java: TypeCheckTest.testInterfaceAssignment13()
func (c *TypeCheckTest) TestInterfaceAssignment13() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */var I = function() {};
/** @constructor
 * @implements {I}*/var T = function() {};
/** @constructor */function Super() {};
/** @return {I} */Super.prototype.foo = function() { return new T(); };
/** @constructor
 * @extends {Super} */function Sub() {}
/** @override
 * @return {T} */Sub.prototype.foo = function() { return new T(); };
`).
		Run()
}

// Java: TypeCheckTest.testGetprop1()
func (c *TypeCheckTest) TestGetprop1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @return {void}*/function foo(){foo().bar;}").
		AddDiagnostic(`No properties on this expression
found   : undefined
required: Object
`).
		Run()
}

// Java: TypeCheckTest.testGetprop2()
func (c *TypeCheckTest) TestGetprop2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var x = null; x.alert();").
		AddDiagnostic(`No properties on this expression
found   : null
required: Object
`).
		Run()
}

// Java: TypeCheckTest.testGetprop3()
func (c *TypeCheckTest) TestGetprop3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */
function Foo() { /** @type {?Object} */ this.x = null; }
Foo.prototype.initX = function() { this.x = {foo: 1}; };
Foo.prototype.bar = function() {
  if (this.x == null) { this.initX(); alert(this.x.foo); }
};
`).
		Run()
}

// Java: TypeCheckTest.testGetprop4()
func (c *TypeCheckTest) TestGetprop4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var x = null; x.prop = 3;").
		AddDiagnostic(`No properties on this expression
found   : null
required: Object
`).
		Run()
}

// Java: TypeCheckTest.testGetrop_interfaceWithoutTypeDeclaration()
func (c *TypeCheckTest) TestGetrop_interfaceWithoutTypeDeclaration() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */var I = function() {};
// Note that we didn't declare the type but we still expect JsCompiler to recognize the
// property.
I.prototype.foo;
var v = /** @type {I} */ (null);
v.foo = 5;
`).
		Run()
}

// Java: TypeCheckTest.testGetrop_interfaceEs6WithoutTypeDeclaration()
func (c *TypeCheckTest) TestGetrop_interfaceEs6WithoutTypeDeclaration() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
class I {
  constructor() {
    this.foo;
  }
}
var v = /** @type {I} */ (null);
v.foo = 5;
`).
		AddDiagnostic("Property foo never defined on I").
		AddDiagnostic(TypeCheckTest_ILLEGAL_PROPERTY_CREATION_MESSAGE).
		Run()
}

// Java: TypeCheckTest.testGetrop_interfaceEs6WithTypeDeclaration()
func (c *TypeCheckTest) TestGetrop_interfaceEs6WithTypeDeclaration() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
class I {
  constructor() {
    /** @type {number} */ this.foo;
  }
}
var v = /** @type {I} */ (null);
v.foo = 5;
`).
		Run()
}

// Java: TypeCheckTest.testSetprop1()
func (c *TypeCheckTest) TestSetprop1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @struct
 */
function Foo() { this.x = 123; }
`).
		Run()
}

// Java: TypeCheckTest.testSetprop2()
func (c *TypeCheckTest) TestSetprop2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @struct
 */
function Foo() {}
(new Foo()).x = 123;
`).
		AddDiagnostic(TypeCheckTest_ILLEGAL_PROPERTY_CREATION_MESSAGE).
		Run()
}

// Java: TypeCheckTest.testSetprop3()
func (c *TypeCheckTest) TestSetprop3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @struct
 */
function Foo() {}
(function() { (new Foo()).x = 123; })();
`).
		AddDiagnostic(TypeCheckTest_ILLEGAL_PROPERTY_CREATION_MESSAGE).
		Run()
}

// Java: TypeCheckTest.testSetprop4()
func (c *TypeCheckTest) TestSetprop4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @struct
 */
function Foo() { this.x = 123; }
(new Foo()).x = "asdf";
`).
		Run()
}

// Java: TypeCheckTest.testSetprop5()
func (c *TypeCheckTest) TestSetprop5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @struct
 */
function Foo() {}
(true ? new Foo() : {}).x = 123;
`).
		AddDiagnosticType(TypeCheck_ILLEGAL_PROPERTY_CREATION_ON_UNION_TYPE).
		Run()
}

// Java: TypeCheckTest.testSetprop6()
func (c *TypeCheckTest) TestSetprop6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @struct
 */
function Foo() {}
/**
 * @constructor
 * @param{Foo} f
 */
function Bar(f) { f.x = 123; }
`).
		AddDiagnostic(TypeCheckTest_ILLEGAL_PROPERTY_CREATION_MESSAGE).
		Run()
}

// Java: TypeCheckTest.testSetprop7()
func (c *TypeCheckTest) TestSetprop7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @struct
 */
function Foo() {
  var t = this;
  t.x = 123;
}
`).
		AddDiagnostic(TypeCheckTest_ILLEGAL_PROPERTY_CREATION_MESSAGE).
		Run()
}

// Java: TypeCheckTest.testSetprop8()
func (c *TypeCheckTest) TestSetprop8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @struct
 */
function Foo() {}
(new Foo()).x--;
`).
		AddDiagnostic(TypeCheckTest_ILLEGAL_PROPERTY_CREATION_MESSAGE).
		AddDiagnostic("Property x never defined on Foo").
		Run()
}

// Java: TypeCheckTest.testSetprop9()
func (c *TypeCheckTest) TestSetprop9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @struct
 */
function Foo() {}
(new Foo()).x += 123;
`).
		AddDiagnostic(TypeCheckTest_ILLEGAL_PROPERTY_CREATION_MESSAGE).
		AddDiagnostic("Property x never defined on Foo").
		Run()
}

// Java: TypeCheckTest.testSetprop10()
func (c *TypeCheckTest) TestSetprop10() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** 
 * @constructor 
 * @struct 
 */ 
function Square(side) { 
  this.side = side; 
} 
Square.prototype = /** @struct */ {
  area: function() { return this.side * this.side; }
};
Square.prototype.id = function(x) { return x; };
`).
		Run()
}

// Java: TypeCheckTest.testSetprop11()
func (c *TypeCheckTest) TestSetprop11() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @struct
 */
function Foo() {}
/** @constructor */
function Bar() {}
Bar.prototype = new Foo();
Bar.prototype.someprop = 123;
`).
		Run()
}

// Java: TypeCheckTest.testSetprop12()
func (c *TypeCheckTest) TestSetprop12() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @struct
 */
function Foo() {}
Foo.someprop = 123;
`).
		Run()
}

// Java: TypeCheckTest.testSetprop13()
func (c *TypeCheckTest) TestSetprop13() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @struct
 */
function Parent() {}
/**
 * @constructor
 * @extends {Parent}
 */
function Kid() {}
Kid.prototype.foo = 123;
var x = (new Kid()).foo;
`).
		Run()
}

// Java: TypeCheckTest.testSetprop14()
func (c *TypeCheckTest) TestSetprop14() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @struct
 */
function Top() {}
/**
 * @constructor
 * @extends {Top}
 */
function Mid() {}
/** blah blah */
Mid.prototype.foo = function() { return 1; };
/**
 * @constructor
 * @struct
 * @extends {Mid}
 */
function Bottom() {}
/** @override */
Bottom.prototype.foo = function() { return 3; };
`).
		Run()
}

// Java: TypeCheckTest.testSetprop15()
func (c *TypeCheckTest) TestSetprop15() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
function Peelable() {};
/** @return {undefined} */
Peelable.prototype.peel;
/**
 * @constructor
 * @struct
 */
function Fruit() {};
/**
 * @constructor
 * @extends {Fruit}
 * @implements {Peelable}
 */
function Banana() { };
function f() {};
/** @override */
Banana.prototype.peel = f;
`).
		Run()
}

// Java: TypeCheckTest.testGetpropDict1()
func (c *TypeCheckTest) TestGetpropDict1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @dict
 */
function Dict1(){ this['prop'] = 123; }
/** @param{Dict1} x */
function takesDict(x) { return x.prop; }
`).
		AddDiagnostic("Cannot do '.' access on a dict").
		Run()
}

// Java: TypeCheckTest.testGetpropDict2()
func (c *TypeCheckTest) TestGetpropDict2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @dict
 */
function Dict1(){ this['prop'] = 123; }
/**
 * @constructor
 * @extends {Dict1}
 */
function Dict1kid(){ this['prop'] = 123; }
/** @param{Dict1kid} x */
function takesDict(x) { return x.prop; }
`).
		AddDiagnostic("Cannot do '.' access on a dict").
		Run()
}

// Java: TypeCheckTest.testGetpropDict3()
func (c *TypeCheckTest) TestGetpropDict3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @dict
 */
function Dict1() { this['prop'] = 123; }
/** @constructor */
function NonDict() { this.prop = 321; }
/** @param{(NonDict|Dict1)} x */
function takesDict(x) { return x.prop; }
`).
		AddDiagnostic("Cannot do '.' access on a dict").
		Run()
}

// Java: TypeCheckTest.testGetpropDict4()
func (c *TypeCheckTest) TestGetpropDict4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @dict
 */
function Dict1() { this['prop'] = 123; }
/**
 * @constructor
 * @struct
 */
function Struct1() { this.prop = 123; }
/** @param{(Struct1|Dict1)} x */
function takesNothing(x) { return x.prop; }
`).
		AddDiagnostic("Cannot do '.' access on a dict").
		Run()
}

// Java: TypeCheckTest.testGetpropDict5()
func (c *TypeCheckTest) TestGetpropDict5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @dict
 */
function Dict1(){ this.prop = 123; }
`).
		AddDiagnostic("Cannot do '.' access on a dict").
		Run()
}

// Java: TypeCheckTest.testGetpropDict6()
func (c *TypeCheckTest) TestGetpropDict6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @dict
 */
function Foo() {}
function Bar() {}
Bar.prototype = new Foo();
Bar.prototype.someprop = 123;
`).
		Run()
}

// Java: TypeCheckTest.testGetpropDict7()
func (c *TypeCheckTest) TestGetpropDict7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("(/** @dict */ {'x': 123}).x = 321;").
		AddDiagnostic("Cannot do '.' access on a dict").
		Run()
}

// Java: TypeCheckTest.testDictEs5ClassCannotExtendStruct()
func (c *TypeCheckTest) TestDictEs5ClassCannotExtendStruct() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @struct @constructor*/
function StructParent() {}
/** @constructor @dict @extends {StructParent} */
function DictChild() {}
DictChild.prototype['prop'] = function() {};
`).
		AddDiagnostic("@dict class DictChild cannot extend @struct class StructParent").
		Run()
}

// Java: TypeCheckTest.testDictEs6ClassCannotExtendStruct()
func (c *TypeCheckTest) TestDictEs6ClassCannotExtendStruct() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class StructParent {}
/** @dict */
class DictChild extends StructParent {
  ['prop']() {}
}
`).
		AddDiagnostic("@dict class DictChild cannot extend @struct class StructParent").
		Run()
}

// Java: TypeCheckTest.testStructEs6ClassCannotExtendDict()
func (c *TypeCheckTest) TestStructEs6ClassCannotExtendDict() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @dict */
class DictParent {}
/** @struct */
class StructChild extends DictParent {
  ['prop']() {}
}
`).
		AddDiagnostic("@struct class StructChild cannot extend @dict class DictParent").
		AddDiagnostic("Cannot do '[]' access on a struct").
		Run()
}

// Java: TypeCheckTest.testStructEs6ClassCanExtendDict_transitive()
func (c *TypeCheckTest) TestStructEs6ClassCanExtendDict_transitive() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @dict */
class DictParent {}
/** @unrestricted */
class UnannotatedMiddleClass extends DictParent {}
/** @struct */
class StructChild extends UnannotatedMiddleClass {
  ['prop']() {}
}
`).
		AddDiagnostic("Cannot do '[]' access on a struct").
		Run()
}

// Java: TypeCheckTest.testGetelemStruct1()
func (c *TypeCheckTest) TestGetelemStruct1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @struct
 */
function Struct1(){ this.prop = 123; }
/** @param{Struct1} x */
function takesStruct(x) {
  var z = x;
  return z['prop'];
}
`).
		AddDiagnostic("Cannot do '[]' access on a struct").
		Run()
}

// Java: TypeCheckTest.testGetelemStruct2()
func (c *TypeCheckTest) TestGetelemStruct2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @struct
 */
function Struct1(){ this.prop = 123; }
/**
 * @constructor
 * @extends {Struct1}
 */
function Struct1kid(){ this.prop = 123; }
/** @param{Struct1kid} x */
function takesStruct2(x) { return x['prop']; }
`).
		AddDiagnostic("Cannot do '[]' access on a struct").
		Run()
}

// Java: TypeCheckTest.testGetelemStruct3()
func (c *TypeCheckTest) TestGetelemStruct3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @struct
 */
function Struct1(){ this.prop = 123; }
/**
 * @constructor
 * @extends {Struct1}
 */
function Struct1kid(){ this.prop = 123; }
var x = (new Struct1kid())['prop'];
`).
		AddDiagnostic("Cannot do '[]' access on a struct").
		Run()
}

// Java: TypeCheckTest.testGetelemStruct4()
func (c *TypeCheckTest) TestGetelemStruct4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @struct
 */
function Struct1() { this.prop = 123; }
/** @constructor */
function NonStruct() { this.prop = 321; }
/** @param{(NonStruct|Struct1)} x */
function takesStruct(x) { return x['prop']; }
`).
		AddDiagnostic("Cannot do '[]' access on a struct").
		Run()
}

// Java: TypeCheckTest.testGetelemStruct5()
func (c *TypeCheckTest) TestGetelemStruct5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @struct
 */
function Struct1() { this.prop = 123; }
/**
 * @constructor
 * @dict
 */
function Dict1() { this['prop'] = 123; }
/** @param{(Struct1|Dict1)} x */
function takesNothing(x) { return x['prop']; }
`).
		AddDiagnostic("Cannot do '[]' access on a struct").
		Run()
}

// Java: TypeCheckTest.testGetelemStruct6()
func (c *TypeCheckTest) TestGetelemStruct6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ function Foo(){}
/**
 * @constructor
 * @struct
 * @implements {Foo}
 */
function Bar(){ this.x = 123; }
var z = /** @type {Foo} */(new Bar())['x'];
`).
		Run()
}

// Java: TypeCheckTest.testGetelemStruct7()
func (c *TypeCheckTest) TestGetelemStruct7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @struct
 */
function Foo() {}
/** @constructor */
function Bar() {}
Bar.prototype = new Foo();
Bar.prototype['someprop'] = 123;
`).
		Run()
}

// Java: TypeCheckTest.testGetelemStruct_noErrorForSettingWellKnownSymbol()
func (c *TypeCheckTest) TestGetelemStruct_noErrorForSettingWellKnownSymbol() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @struct
 */
function Foo() {}
Foo.prototype[Symbol.iterator] = 123;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testGetelemStruct_noErrorForGettingWellKnownSymbol()
func (c *TypeCheckTest) TestGetelemStruct_noErrorForGettingWellKnownSymbol() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @struct
 */
function Foo() {}
/** @param {!Foo} foo */
function getIterator(foo) { return foo[Symbol.iterator](); }
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testInOnStruct()
func (c *TypeCheckTest) TestInOnStruct() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @struct
 */
function Foo() {}
if ('prop' in (new Foo())) {}
`).
		AddDiagnostic("Cannot use the IN operator with structs").
		Run()
}

// Java: TypeCheckTest.testForinOnStruct()
func (c *TypeCheckTest) TestForinOnStruct() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @struct
 */
function Foo() {}
for (var prop in (new Foo())) {}
`).
		AddDiagnostic("Cannot use the IN operator with structs").
		Run()
}

// Java: TypeCheckTest.testIArrayLikeAccess1()
func (c *TypeCheckTest) TestIArrayLikeAccess1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {!IArrayLike<T>} x
 * @return {T}
 * @template T
*/
function f(x) { return x[0]; }
function g(/** !Array<string> */ x) {
  var /** null */ y = f(x);
}
`).
		AddDiagnostic(`initializing variable
found   : string
required: null
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testIArrayLikeAccess2()
func (c *TypeCheckTest) TestIArrayLikeAccess2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {!IArrayLike<T>} x
 * @return {T}
 * @template T
*/
function f(x) { return x[0]; }
function g(/** !IArrayLike<string> */ x) {
  var /** null */ y = f(x);
}
`).
		AddDiagnostic(`initializing variable
found   : string
required: null
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testIteratorAccess1()
func (c *TypeCheckTest) TestIteratorAccess1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {!IteratorLike<T>} x
 * @return {T}
 * @template T
*/
function f(x) { return x[0]; }
function g(/** !Generator<string> */ x, /** !Iterator<number> */ it) {
  var /** null */ y = f(x);
  var /** null */ z = f(it);
}
`).
		AddDiagnostic(`initializing variable
found   : string
required: null
`).
		AddDiagnostic(`initializing variable
found   : number
required: null
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testIterableAccess1()
func (c *TypeCheckTest) TestIterableAccess1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {!Iterable<T>} x
 * @return {T}
 * @template T
*/
function f(x) { return x[0]; }
function g(/** !Generator<string> */ x) {
  var /** null */ y = f(x);
}
`).
		AddDiagnostic(`initializing variable
found   : string
required: null
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testIteratorIterableAccess1()
func (c *TypeCheckTest) TestIteratorIterableAccess1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {!IteratorIterable<T>} x
 * @return {T}
 * @template T
*/
function f(x) { return x[0]; }
function g(/** !Generator<string> */ x) {
  var /** null */ y = f(x);
}
`).
		AddDiagnostic(`initializing variable
found   : string
required: null
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testIteratorLikeAccess1()
func (c *TypeCheckTest) TestIteratorLikeAccess1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {!IteratorLike<T>} x
 * @return {T}
 * @template T
 */
function f(x) { return x[0]; }
function g(/** !Iterator<string> */ x) {
  var /** null */ y = f(x);
}
`).
		AddDiagnostic(`initializing variable
found   : string
required: null
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testIteratorFrom_acceptsString()
func (c *TypeCheckTest) TestIteratorFrom_acceptsString() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var /** !Iterator<string> */ it = Iterator.from('abc');
var /** !Iterator<number> */ it2 = Iterator.from([1, 2]);
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testIteratorFrom_rejectsNumber()
func (c *TypeCheckTest) TestIteratorFrom_rejectsNumber() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`Iterator.from(123);
`).
		AddDiagnostic(`actual parameter 1 of Iterator.from does not match formal parameter
found   : number
required: (Iterable<?,?,?>|IteratorLike<?,?,?>|string)
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testIteratorHelperMethods_callbacksAcceptNonBooleanAndNonVoid()
func (c *TypeCheckTest) TestIteratorHelperMethods_callbacksAcceptNonBooleanAndNonVoid() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {!Iterator<?string>} it
 * @param {!Array<number>} arr
 */
function f(it, arr) {
  var /** !Iterator<?string> */ filtered = it.filter((x) => x);
  var /** boolean */ someres = it.some((s) => s ? s.length : 0);
  var /** boolean */ everyres = it.every((s) => s);
  var /** ?string|undefined */ found = it.find((s) => s);
  it.forEach((x) => arr.push(1));
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testIteratorHelperMethods_chainingAndReturnTypes()
func (c *TypeCheckTest) TestIteratorHelperMethods_chainingAndReturnTypes() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {!Iterator<number>} it */
function f(it) {
  var /** !Iterator<string> */ mapped = it.map((x) => String(x));
  var /** !Iterator<number> */ taken = it.take(5);
  var /** !Iterator<number> */ dropped = it.drop(2);
  var /** !Iterator<string> */ flatMapped = it.flatMap((x) => [String(x)]);
  var /** number */ reduced = it.reduce((acc, x) => acc + x, 0);
  var /** !Array<number> */ arr = it.toArray();
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testIteratorFlatMap_rejectsStringReturn()
func (c *TypeCheckTest) TestIteratorFlatMap_rejectsStringReturn() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {!Iterator<number>} it */
function f(it) {
  it.flatMap((x) => 'hello');
}
`).
		AddDiagnostic(`inconsistent return type
found   : string
required: (Iterable<?,?,?>|Iterator<?,?,?>|IteratorLike<?,?,?>)
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testTestExternsBuilder_iteratorHelpers()
func (c *TypeCheckTest) TestTestExternsBuilder_iteratorHelpers() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddIterable().Build()).
		AddSource(`/** @param {!Iterator<number>} it */
function f(it) {
  var /** !Iterator<number> */ filtered = it.filter((x) => x);
  var /** boolean */ s = it.some((x) => x);
  it.forEach((x) => 42);
  var /** !Iterator<string> */ fromStr = Iterator.from('abc');
}
`).
		Run()
}

// Java: TypeCheckTest.testIteratorLike_optionalReturnAndThrow()
func (c *TypeCheckTest) TestIteratorLike_optionalReturnAndThrow() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`// Object with only next() satisfies IteratorLike record
var /** !IteratorLike<number> */ it1 = {
  next: function() { return {done: true, value: 1}; }
};
// Object with next(), return(), and throw() satisfies IteratorLike record
var /** !IteratorLike<number> */ it2 = {
  next: function() { return {done: true, value: 1}; },
  return: function(opt_val) { return {done: true, value: 1}; },
  throw: function(opt_err) { return {done: true, value: 1}; }
};
// Calling return() on IteratorLike is permitted without cast
function close(/** !IteratorLike<number> */ it) {
  if (it.return) {
    var /** !IIterableResult<number> */ res = it.return();
  }
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testIterator_returnMethod()
func (c *TypeCheckTest) TestIterator_returnMethod() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function closeIterator(/** !Iterator<number> */ it) {
  if (it.return) {
    var /** !IIterableResult<number> */ res = it.return();
    var /** !IIterableResult<number> */ res2 = it.return(42);
  }
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testIterator_throwMethod()
func (c *TypeCheckTest) TestIterator_throwMethod() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function throwIterator(/** !Iterator<number> */ it) {
  if (it.throw) {
    var /** !IIterableResult<number> */ res = it.throw();
    var /** !IIterableResult<number> */ res2 = it.throw(new Error('err'));
  }
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testArrayAccess1()
func (c *TypeCheckTest) TestArrayAccess1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var a = []; var b = a['hi'];").
		AddDiagnostic(`restricted index type
found   : string
required: number
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testArrayAccess2()
func (c *TypeCheckTest) TestArrayAccess2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var a = []; var b = a[[1,2]];").
		AddDiagnostic(`restricted index type
found   : Array<?>
required: number
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testArrayAccess3()
func (c *TypeCheckTest) TestArrayAccess3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var bar = [];
/** @return {void} */function baz(){};
var foo = bar[baz()];
`).
		AddDiagnostic(`restricted index type
found   : undefined
required: number
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testArrayAccess4()
func (c *TypeCheckTest) TestArrayAccess4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@return {!Array}*/function foo(){};var bar = foo()[foo()];").
		AddDiagnostic(`restricted index type
found   : Array
required: number
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testArrayAccess6()
func (c *TypeCheckTest) TestArrayAccess6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var bar = null[1];").
		AddDiagnostic(`only arrays or objects can be accessed
found   : null
required: Object
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testArrayAccess7()
func (c *TypeCheckTest) TestArrayAccess7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var bar = void 0; bar[0];").
		AddDiagnostic(`only arrays or objects can be accessed
found   : undefined
required: Object
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testArrayAccess8()
func (c *TypeCheckTest) TestArrayAccess8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var bar = void 0; bar[0]; bar[1];").
		AddDiagnostic(`only arrays or objects can be accessed
found   : undefined
required: Object
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testArrayAccess9()
func (c *TypeCheckTest) TestArrayAccess9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {?Array} */ function f() { return []; }
f()[{}]
`).
		AddDiagnostic(`restricted index type
found   : {}
required: number
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testPropAccess()
func (c *TypeCheckTest) TestPropAccess() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {*} x */var f = function(x) {
var o = String(x);
if (typeof o['a'] != 'undefined') { return o['a']; }
return null;
};
`).
		Run()
}

// Java: TypeCheckTest.testPropAccess2()
func (c *TypeCheckTest) TestPropAccess2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var bar = void 0; bar.baz;").
		AddDiagnostic(`No properties on this expression
found   : undefined
required: Object
`).
		Run()
}

// Java: TypeCheckTest.testPropAccess3()
func (c *TypeCheckTest) TestPropAccess3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var bar = void 0; bar.baz; bar.bax;").
		AddDiagnostic(`No properties on this expression
found   : undefined
required: Object
`).
		Run()
}

// Java: TypeCheckTest.testPropAccess4()
func (c *TypeCheckTest) TestPropAccess4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @param {*} x */ function f(x) { return x['hi']; }").
		Run()
}

// Java: TypeCheckTest.testSwitchCase1()
func (c *TypeCheckTest) TestSwitchCase1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**@type {number}*/var a;
/**@type {string}*/var b;
switch(a){case b:;}
`).
		AddDiagnostic(`case expression doesn't match switch
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testSwitchCase2()
func (c *TypeCheckTest) TestSwitchCase2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var a = null; switch (typeof a) { case 'foo': }").
		Run()
}

// Java: TypeCheckTest.testVar1()
func (c *TypeCheckTest) TestVar1() {
	var p *TypeCheckTestCase_TypeCheckResult = c.ParseAndTypeCheckWithScope("/** @type {(string|null)} */var a = null")
	c.AssertTypeEquals(c.CreateUnionType(c.GetNativeStringType(), c.GetNativeNullType()), p.scope.GetVar("a").GetType())
}

// Java: TypeCheckTest.testVar2()
func (c *TypeCheckTest) TestVar2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {Function} */ var a = function(){}").
		Run()
}

// Java: TypeCheckTest.testVar3()
func (c *TypeCheckTest) TestVar3() {
	var p *TypeCheckTestCase_TypeCheckResult = c.ParseAndTypeCheckWithScope("var a = 3;")
	c.AssertTypeEquals(c.GetNativeNumberType(), p.scope.GetVar("a").GetType())
}

// Java: TypeCheckTest.testVar4()
func (c *TypeCheckTest) TestVar4() {
	var p *TypeCheckTestCase_TypeCheckResult = c.ParseAndTypeCheckWithScope("var a = 3; a = 'string';")
	c.AssertTypeEquals(c.CreateUnionType(c.GetNativeStringType(), c.GetNativeNumberType()), p.scope.GetVar("a").GetType())
}

// Java: TypeCheckTest.testVar5()
func (c *TypeCheckTest) TestVar5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = {};
/** @type {string} */goog.foo = 'hello';
/** @type {number} */var a = goog.foo;
`).
		AddDiagnostic(`initializing variable
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testVar6()
func (c *TypeCheckTest) TestVar6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f() {
  return function() {
    /** @type {!Date} */
    var a = 7;
  };
}
`).
		AddDiagnostic(`initializing variable
found   : number
required: Date
`).
		Run()
}

// Java: TypeCheckTest.testVar7()
func (c *TypeCheckTest) TestVar7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {number} */var a, b;").
		AddDiagnostic("declaration of multiple variables with shared type information").
		Run()
}

// Java: TypeCheckTest.testVar8()
func (c *TypeCheckTest) TestVar8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var a, b;").
		Run()
}

// Java: TypeCheckTest.testVar9()
func (c *TypeCheckTest) TestVar9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @enum */var a;").
		AddDiagnostic("enum initializer must be an object literal or an enum").
		Run()
}

// Java: TypeCheckTest.testVar10()
func (c *TypeCheckTest) TestVar10() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {!Number} */var foo = 'abc';").
		AddDiagnostic(`initializing variable
found   : string
required: Number
`).
		Run()
}

// Java: TypeCheckTest.testVar11()
func (c *TypeCheckTest) TestVar11() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var /** @type {!Date} */foo = 'abc';").
		AddDiagnostic(`initializing variable
found   : string
required: Date
`).
		Run()
}

// Java: TypeCheckTest.testVar12()
func (c *TypeCheckTest) TestVar12() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var /** @type {!Date} */foo = 'abc',
/** @type {!RegExp} */bar = 5;
`).
		AddDiagnostic(`initializing variable
found   : string
required: Date
`).
		AddDiagnostic(`initializing variable
found   : number
required: RegExp
`).
		Run()
}

// Java: TypeCheckTest.testVar13()
func (c *TypeCheckTest) TestVar13() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var /** @type {number} */a,a;").
		Run()
}

// Java: TypeCheckTest.testVar14()
func (c *TypeCheckTest) TestVar14() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @return {number} */ function f() { var x; return x; }").
		AddDiagnostic(`inconsistent return type
found   : undefined
required: number
`).
		Run()
}

// Java: TypeCheckTest.testVar15()
func (c *TypeCheckTest) TestVar15() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {number} */
function f() { var x = x || {}; return x; }
`).
		AddDiagnostic(`inconsistent return type
found   : {}
required: number
`).
		Run()
}

// Java: TypeCheckTest.testVar15NullishCoalesce()
func (c *TypeCheckTest) TestVar15NullishCoalesce() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {number} */
function f() { var x = x ?? {}; return x; }
`).
		AddDiagnostic(`inconsistent return type
found   : {}
required: number
`).
		Run()
}

// Java: TypeCheckTest.testAssign1()
func (c *TypeCheckTest) TestAssign1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = {};
/** @type {number} */goog.foo = 'hello';
`).
		AddDiagnostic(`assignment to property foo of goog
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testAssign2()
func (c *TypeCheckTest) TestAssign2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = {};
/** @type {number}  */goog.foo = 3;
goog.foo = 'hello';
`).
		AddDiagnostic(`assignment to property foo of goog
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testAssign3()
func (c *TypeCheckTest) TestAssign3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = {};
/** @type {number}  */goog.foo = 3;
goog.foo = 4;
`).
		Run()
}

// Java: TypeCheckTest.testAssign4()
func (c *TypeCheckTest) TestAssign4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = {};
goog.foo = 3;
goog.foo = 'hello';
`).
		Run()
}

// Java: TypeCheckTest.testAssignInference()
func (c *TypeCheckTest) TestAssignInference() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {Array} x
 * @return {number}
 */
function f(x) {
  var y = null;
  y = x[0];
  if (y == null) { return 4; } else { return 6; }
}
`).
		Run()
}

// Java: TypeCheckTest.testAssignReadonlyArrayValueFails()
func (c *TypeCheckTest) TestAssignReadonlyArrayValueFails() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().IncludeDefaultExterns().
		AddSource(`const foo = /** @type {!ReadonlyArray<number>} */ ([5]);
foo[0] = 3;
`).
		DiagnosticsAreErrors().
		AddDiagnosticType(TypeCheck_PROPERTY_ASSIGNMENT_TO_READONLY_VALUE).
		Run()
}

// Java: TypeCheckTest.testAssignReadonlyArrayValueFailsWithoutTypeParameter()
func (c *TypeCheckTest) TestAssignReadonlyArrayValueFailsWithoutTypeParameter() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().IncludeDefaultExterns().
		AddSource(`const foo = /** @type {!ReadonlyArray} */ ([5]);
foo[0] = 3;
`).
		DiagnosticsAreErrors().
		AddDiagnosticType(TypeCheck_PROPERTY_ASSIGNMENT_TO_READONLY_VALUE).
		Run()
}

// Java: TypeCheckTest.testAssignReadonlyArrayLengthFails()
func (c *TypeCheckTest) TestAssignReadonlyArrayLengthFails() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().IncludeDefaultExterns().
		AddSource(`const foo = /** @type {!ReadonlyArray<number>} */ ([5]);
foo.length = 0;
`).
		DiagnosticsAreErrors().
		AddDiagnosticType(TypeCheck_PROPERTY_ASSIGNMENT_TO_READONLY_VALUE).
		Run()
}

// Java: TypeCheckTest.testOr1()
func (c *TypeCheckTest) TestOr1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {number}  */var a;
/** @type {number}  */var b;
a + b || undefined;
`).
		Run()
}

// Java: TypeCheckTest.testNullishCoalesceNumber()
func (c *TypeCheckTest) TestNullishCoalesceNumber() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {number}  */var a;
/** @type {number}  */var b;
a + b ?? undefined;
`).
		Run()
}

// Java: TypeCheckTest.testOr2()
func (c *TypeCheckTest) TestOr2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {number}  */var a;
/** @type {number}  */var b;
/** @type {number}  */var c = a + b || undefined;
`).
		AddDiagnostic(`initializing variable
found   : (number|undefined)
required: number
`).
		Run()
}

// Java: TypeCheckTest.testNullishCoalesceNumberVar()
func (c *TypeCheckTest) TestNullishCoalesceNumberVar() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {number}  */var a;
/** @type {number}  */var c = a ?? undefined;
`).
		Run()
}

// Java: TypeCheckTest.testOr3()
func (c *TypeCheckTest) TestOr3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {(number|undefined)} */var a;
/** @type {number}  */var c = a || 3;
`).
		Run()
}

// Java: TypeCheckTest.testNullishCoalesceNumberUndefined()
func (c *TypeCheckTest) TestNullishCoalesceNumberUndefined() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {(number|undefined)} */var a;
/** @type {number}  */var c = a ?? 3;
`).
		Run()
}

// Java: TypeCheckTest.testOr4()
func (c *TypeCheckTest) TestOr4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@type {number} */var x;x=null || \"a\";").
		AddDiagnostic(`assignment
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testOr5()
func (c *TypeCheckTest) TestOr5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@type {number} */var x;x=undefined || \"a\";").
		AddDiagnostic(`assignment
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testNullishCoalesceAssignment2()
func (c *TypeCheckTest) TestNullishCoalesceAssignment2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@type {number} */var x;x=undefined ?? \"a\";").
		AddDiagnostic(`assignment
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testOr6()
func (c *TypeCheckTest) TestOr6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {!Array=} opt_x */
function removeDuplicates(opt_x) {
  var x = opt_x || [];
  var /** undefined */ y = x;
}
`).
		AddDiagnostic(`initializing variable
found   : Array
required: undefined
`).
		Run()
}

// Java: TypeCheckTest.testNullishCoaleceAssignment3()
func (c *TypeCheckTest) TestNullishCoaleceAssignment3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {!Array=} opt_x */
function removeDuplicates(opt_x) {
  var x = opt_x ?? [];
  var /** undefined */ y = x;
}
`).
		AddDiagnostic(`initializing variable
found   : Array
required: undefined
`).
		Run()
}

// Java: TypeCheckTest.testAnd1()
func (c *TypeCheckTest) TestAnd1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {number}  */var a;
/** @type {number}  */var b;
a + b && undefined;
`).
		Run()
}

// Java: TypeCheckTest.testAnd2()
func (c *TypeCheckTest) TestAnd2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {number}  */var a;
/** @type {number}  */var b;
/** @type {number}  */var c = a + b && undefined;
`).
		AddDiagnostic(`initializing variable
found   : (number|undefined)
required: number
`).
		Run()
}

// Java: TypeCheckTest.testAnd3()
func (c *TypeCheckTest) TestAnd3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {(!Array|undefined)} */var a;
/** @type {number}  */var c = a && undefined;
`).
		AddDiagnostic(`initializing variable
found   : undefined
required: number
`).
		Run()
}

// Java: TypeCheckTest.testAnd4()
func (c *TypeCheckTest) TestAnd4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {number} x */function f(x){};
/** @type {null}  */var x; /** @type {number?} */var y;
if (x && y) { f(y) }
`).
		Run()
}

// Java: TypeCheckTest.testAnd5()
func (c *TypeCheckTest) TestAnd5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {number} x
@param {string} y*/function f(x,y){};
/** @type {number?} */var x; /** @type {string?} */var y;
if (x && y) { f(x, y) }
`).
		Run()
}

// Java: TypeCheckTest.testAnd6()
func (c *TypeCheckTest) TestAnd6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {number} x */function f(x){};
/** @type {number|undefined} */var x;
if (x && f(x)) { f(x) }
`).
		Run()
}

// Java: TypeCheckTest.testAnd7()
func (c *TypeCheckTest) TestAnd7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {null} */var x; if (x && x) {}").
		Run()
}

// Java: TypeCheckTest.testAnd8()
func (c *TypeCheckTest) TestAnd8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** (null | number | string) */ x) {
  (x && (typeof x === 'number')) && takesNum(x);
}
function takesNum(/** number */ n) {}
`).
		Run()
}

// Java: TypeCheckTest.testAnd9()
func (c *TypeCheckTest) TestAnd9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** (number|string|null) */ x) {
  if (x && typeof x === 'number') {
    takesNum(x);
  }
}
function takesNum(/** number */ x) {}
`).
		Run()
}

// Java: TypeCheckTest.testAnd10()
func (c *TypeCheckTest) TestAnd10() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** (null | number | string) */ x) {
  (x && (typeof x === 'string')) && takesNum(x);
}
function takesNum(/** number */ n) {}
`).
		AddDiagnostic(`actual parameter 1 of takesNum does not match formal parameter
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testHook()
func (c *TypeCheckTest) TestHook() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@return {void}*/function foo(){ var x=foo()?a:b; }").
		Run()
}

// Java: TypeCheckTest.testHookRestrictsType1()
func (c *TypeCheckTest) TestHookRestrictsType1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {(string|null)} */
function f() { return null;}
/** @type {(string|null)} */ var a = f();
/** @type {string} */
var b = a ? a : 'default';
`).
		Run()
}

// Java: TypeCheckTest.testHookRestrictsType2()
func (c *TypeCheckTest) TestHookRestrictsType2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {String} */
var a = null;
/** @type {null} */
var b = a ? null : a;
`).
		Run()
}

// Java: TypeCheckTest.testHookRestrictsType3()
func (c *TypeCheckTest) TestHookRestrictsType3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {String} */
var a;
/** @type {null} */
var b = (!a) ? a : null;
`).
		Run()
}

// Java: TypeCheckTest.testHookRestrictsType4()
func (c *TypeCheckTest) TestHookRestrictsType4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {(boolean|undefined)} */
var a;
/** @type {boolean} */
var b = a != null ? a : true;
`).
		Run()
}

// Java: TypeCheckTest.testHookRestrictsType5()
func (c *TypeCheckTest) TestHookRestrictsType5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {(boolean|undefined)} */
var a;
/** @type {(undefined)} */
var b = a == null ? a : undefined;
`).
		Run()
}

// Java: TypeCheckTest.testHookRestrictsType6()
func (c *TypeCheckTest) TestHookRestrictsType6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {(number|null|undefined)} */
var a;
/** @type {number} */
var b = a == null ? 5 : a;
`).
		Run()
}

// Java: TypeCheckTest.testHookRestrictsType7()
func (c *TypeCheckTest) TestHookRestrictsType7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {(number|null|undefined)} */
var a;
/** @type {number} */
var b = a == undefined ? 5 : a;
`).
		Run()
}

// Java: TypeCheckTest.testWhileRestrictsType1()
func (c *TypeCheckTest) TestWhileRestrictsType1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {null} x */ function g(x) {}
/** @param {number?} x */
function f(x) {
while (x) {
if (g(x)) { x = 1; }
x = x-1;
}
}
`).
		AddDiagnostic(`actual parameter 1 of g does not match formal parameter
found   : number
required: null
`).
		Run()
}

// Java: TypeCheckTest.testWhileRestrictsType2()
func (c *TypeCheckTest) TestWhileRestrictsType2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {number?} x
@return {number}*/
function f(x) {
/** @type {number} */var y = 0;
while (x) {
y = x;
x = x-1;
}
return y;
}
`).
		Run()
}

// Java: TypeCheckTest.testHigherOrderFunctions1()
func (c *TypeCheckTest) TestHigherOrderFunctions1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {function(number)} */var f;
f(true);
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckTest.testHigherOrderFunctions2()
func (c *TypeCheckTest) TestHigherOrderFunctions2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {function():!Date} */var f;
/** @type {boolean} */var a = f();
`).
		AddDiagnostic(`initializing variable
found   : Date
required: boolean
`).
		Run()
}

// Java: TypeCheckTest.testHigherOrderFunctions3()
func (c *TypeCheckTest) TestHigherOrderFunctions3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {function(this:Array):Date} */var f; new f").
		AddDiagnostic("cannot instantiate non-constructor, found type: function(this:Array): (Date|null)").
		Run()
}

// Java: TypeCheckTest.testHigherOrderFunctions4()
func (c *TypeCheckTest) TestHigherOrderFunctions4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {function(this:Array, ...number):Date} */var f; new f").
		AddDiagnostic(`cannot instantiate non-constructor, found type: function(this:Array, ...number): (Date|null)
`).
		Run()
}

// Java: TypeCheckTest.testHigherOrderFunctions5()
func (c *TypeCheckTest) TestHigherOrderFunctions5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {number} x */ function g(x) {}
/** @type {function(new:Array, ...number):Date} */ var f;
g(new f());
`).
		AddDiagnostic(`actual parameter 1 of g does not match formal parameter
found   : Array
required: number
`).
		Run()
}

// Java: TypeCheckTest.testConstructorAlias1()
func (c *TypeCheckTest) TestConstructorAlias1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ var Foo = function() {};
/** @type {number} */ Foo.prototype.bar = 3;
/** @constructor */ var FooAlias = Foo;
/** @return {string} */ function foo() {
  return (new FooAlias()).bar; }
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testConstructorAlias2()
func (c *TypeCheckTest) TestConstructorAlias2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ var Foo = function() {};
/** @constructor */ var FooAlias = Foo;
/** @type {number} */ FooAlias.prototype.bar = 3;
/** @return {string} */ function foo() {
  return (new Foo()).bar; }
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testConstructorAlias3()
func (c *TypeCheckTest) TestConstructorAlias3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ var Foo = function() {};
/** @type {number} */ Foo.prototype.bar = 3;
/** @constructor */ var FooAlias = Foo;
/** @return {string} */ function foo() {
  return (new FooAlias()).bar; }
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testConstructorAlias4()
func (c *TypeCheckTest) TestConstructorAlias4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ var Foo = function() {};
var FooAlias = Foo;
/** @type {number} */ FooAlias.prototype.bar = 3;
/** @return {string} */ function foo() {
  return (new Foo()).bar; }
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testConstructorAlias5()
func (c *TypeCheckTest) TestConstructorAlias5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ var Foo = function() {};
/** @constructor */ var FooAlias = Foo;
/** @return {FooAlias} */ function foo() {
  return new Foo(); }
`).
		Run()
}

// Java: TypeCheckTest.testConstructorAlias6()
func (c *TypeCheckTest) TestConstructorAlias6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ var Foo = function() {};
/** @constructor */ var FooAlias = Foo;
/** @return {Foo} */ function foo() {
  return new FooAlias(); }
`).
		Run()
}

// Java: TypeCheckTest.testConstructorAlias7()
func (c *TypeCheckTest) TestConstructorAlias7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = {};
/** @constructor */ goog.Foo = function() {};
/** @constructor */ goog.FooAlias = goog.Foo;
/** @return {number} */ function foo() {
  return new goog.FooAlias(); }
`).
		AddDiagnostic(`inconsistent return type
found   : goog.Foo
required: number
`).
		Run()
}

// Java: TypeCheckTest.testConstructorAlias8()
func (c *TypeCheckTest) TestConstructorAlias8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = {};
/**
 * @param {number} x
 * @constructor */
goog.Foo = function(x) {};
/**
 * @param {number} x
 * @constructor */
goog.FooAlias = goog.Foo;
/** @return {number} */ function foo() {
  return new goog.FooAlias(1); }
`).
		AddDiagnostic(`inconsistent return type
found   : goog.Foo
required: number
`).
		Run()
}

// Java: TypeCheckTest.testConstructorAlias9()
func (c *TypeCheckTest) TestConstructorAlias9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = {};
/**
 * @param {number} x
 * @constructor */
goog.Foo = function(x) {};
/** @constructor */ goog.FooAlias = goog.Foo;
/** @return {number} */ function foo() {
  return new goog.FooAlias(1); }
`).
		AddDiagnostic(`inconsistent return type
found   : goog.Foo
required: number
`).
		Run()
}

// Java: TypeCheckTest.testConstructorAlias10()
func (c *TypeCheckTest) TestConstructorAlias10() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {number} x
 * @constructor */
var Foo = function(x) {};
/** @constructor */ var FooAlias = Foo;
/** @return {number} */ function foo() {
  return new FooAlias(1); }
`).
		AddDiagnostic(`inconsistent return type
found   : Foo
required: number
`).
		Run()
}

// Java: TypeCheckTest.testConstructorAlias11()
func (c *TypeCheckTest) TestConstructorAlias11() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {number} x
 * @constructor */
var Foo = function(x) {};
/** @const */ var FooAlias = Foo;
/** @const */ var FooAlias2 = FooAlias;
/** @return {FooAlias2} */ function foo() {
  return 1; }
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: (Foo|null)
`).
		Run()
}

// Java: TypeCheckTest.testConstructorAliasWithBadAnnotation1()
func (c *TypeCheckTest) TestConstructorAliasWithBadAnnotation1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
/** @record */ var Bar = Foo;
`).
		AddDiagnostic("Annotation @record on Bar incompatible with aliased type.").
		Run()
}

// Java: TypeCheckTest.testConstructorAliasWithBadAnnotation2()
func (c *TypeCheckTest) TestConstructorAliasWithBadAnnotation2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
/** @interface */ var Bar = Foo;
`).
		AddDiagnostic("Annotation @interface on Bar incompatible with aliased type.").
		Run()
}

// Java: TypeCheckTest.testConstructorAliasWithBadAnnotation3()
func (c *TypeCheckTest) TestConstructorAliasWithBadAnnotation3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ function Foo() {}
/** @record */ var Bar = Foo;
`).
		AddDiagnostic("Annotation @record on Bar incompatible with aliased type.").
		Run()
}

// Java: TypeCheckTest.testConstructorAliasWithBadAnnotation4()
func (c *TypeCheckTest) TestConstructorAliasWithBadAnnotation4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ function Foo() {}
/** @constructor */ var Bar = Foo;
`).
		AddDiagnostic("Annotation @constructor on Bar incompatible with aliased type.").
		Run()
}

// Java: TypeCheckTest.testConstAliasedTypeCastInferredCorrectly1()
func (c *TypeCheckTest) TestConstAliasedTypeCastInferredCorrectly1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */
function Foo() {}
/** @return {number} */
Foo.prototype.foo = function() {};

/** @const */ var FooAlias = Foo;

var x = /** @type {!FooAlias} */ ({});
var /** null */ n = x.foo();
`).
		AddDiagnostic(`initializing variable
found   : number
required: null
`).
		Run()
}

// Java: TypeCheckTest.testConstAliasedTypeCastInferredCorrectly2()
func (c *TypeCheckTest) TestConstAliasedTypeCastInferredCorrectly2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */
function Foo() {}
/** @return {number} */
Foo.prototype.foo = function() {};

var ns = {};
/** @const */ ns.FooAlias = Foo;

var x = /** @type {!ns.FooAlias} */ ({});
var /** null */ n = x.foo();
`).
		AddDiagnostic(`initializing variable
found   : number
required: null
`).
		Run()
}

// Java: TypeCheckTest.testConstructorAliasedTypeCastInferredCorrectly()
func (c *TypeCheckTest) TestConstructorAliasedTypeCastInferredCorrectly() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */
function Foo() {}
/** @return {number} */
Foo.prototype.foo = function() {};

/** @constructor */ var FooAlias = Foo;

var x = /** @type {!FooAlias} */ ({});
var /** null */ n = x.foo();
`).
		AddDiagnostic(`initializing variable
found   : number
required: null
`).
		Run()
}

// Java: TypeCheckTest.testAliasedTypedef1()
func (c *TypeCheckTest) TestAliasedTypedef1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @typedef {string} */ var original;
/** @const */ var alias = original;
/** @type {alias} */ var x;
`).
		Run()
}

// Java: TypeCheckTest.testAliasedTypedef2()
func (c *TypeCheckTest) TestAliasedTypedef2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const */ var ns = {};
/** @typedef {string} */ var original;
/** @const */ ns.alias = original;
/** @type {ns.alias} */ var x;
`).
		Run()
}

// Java: TypeCheckTest.testAliasedTypedef3()
func (c *TypeCheckTest) TestAliasedTypedef3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @typedef {string} */ var original;
/** @const */ var alias = original;
/** @return {number} */ var f = function(/** alias */ x) { return x; }
`).
		AddDiagnostic(`inconsistent return type
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testAliasedTypedef4()
func (c *TypeCheckTest) TestAliasedTypedef4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const */ var ns = {};
/** @typedef {string} */ var original;
/** @const */ ns.alias = original;
/** @return {number} */ var f = function(/** ns.alias */ x) { return x; }
`).
		AddDiagnostic(`inconsistent return type
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testAliasedNonTypedef()
func (c *TypeCheckTest) TestAliasedNonTypedef() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {string} */ var notTypeName;
/** @const */ var alias = notTypeName;
/** @type {alias} */ var x;
`).
		AddDiagnostic(`Bad type annotation. Unknown type alias
It's possible that 'alias' refers to a value, not a type.
`).
		Run()
}

// Java: TypeCheckTest.testConstStringKeyDoesntCrash()
func (c *TypeCheckTest) TestConstStringKeyDoesntCrash() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */
function Foo() {}
var ns = {
  /** @const */ FooAlias: Foo
};
`).
		Run()
}

// Java: TypeCheckTest.testClosure7()
func (c *TypeCheckTest) TestClosure7() {
	c.TestClosureTypes(`/** @type {string|null|undefined} */ var a = foo();
/** @type {number} */
var b = goog.asserts.assert(a);
`, jutil.StringPtr(`initializing variable
found   : string
required: number
`))
}

// Java: TypeCheckTest.testPrimitiveAssertTruthy_removesNullAndUndefinedFromString()
func (c *TypeCheckTest) TestPrimitiveAssertTruthy_removesNullAndUndefinedFromString() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(TypeCheckTest_PRIMITIVE_ASSERT_DEFS + `function f(/** ?string|undefined */ str) {
  assertTruthy(str);
  const /** number */ n = str;
}
`).
		AddDiagnostic(`initializing variable
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testPrimitiveAssertString_narrowsAllTypeToString()
func (c *TypeCheckTest) TestPrimitiveAssertString_narrowsAllTypeToString() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(TypeCheckTest_PRIMITIVE_ASSERT_DEFS + `function f(/** * */ str) {
  assertString(str);
  const /** number */ n = str;
}
`).
		AddDiagnostic(`initializing variable
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testReturn1()
func (c *TypeCheckTest) TestReturn1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@return {void}*/function foo(){ return 3; }").
		AddDiagnostic(`inconsistent return type
found   : number
required: undefined
`).
		Run()
}

// Java: TypeCheckTest.testReturn2()
func (c *TypeCheckTest) TestReturn2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@return {!Number}*/function foo(){ return; }").
		AddDiagnostic(`inconsistent return type
found   : undefined
required: Number
`).
		Run()
}

// Java: TypeCheckTest.testReturn3()
func (c *TypeCheckTest) TestReturn3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@return {!Number}*/function foo(){ return 'abc'; }").
		AddDiagnostic(`inconsistent return type
found   : string
required: Number
`).
		Run()
}

// Java: TypeCheckTest.testReturn4()
func (c *TypeCheckTest) TestReturn4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**@return {!Number}
*/
 function a(){return new Array();}`).
		AddDiagnostic(`inconsistent return type
found   : Array<?>
required: Number
`).
		Run()
}

// Java: TypeCheckTest.testReturn5()
func (c *TypeCheckTest) TestReturn5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {number} n
 * @constructor
 */
function fn(n){ return }
`).
		Run()
}

// Java: TypeCheckTest.testReturn6()
func (c *TypeCheckTest) TestReturn6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {number} opt_a
@return {string} */
function a(opt_a) { return opt_a }
`).
		AddDiagnostic(`inconsistent return type
found   : (number|undefined)
required: string
`).
		Run()
}

// Java: TypeCheckTest.testReturn7()
func (c *TypeCheckTest) TestReturn7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */var A = function() {};
/** @constructor */var B = function() {};
/** @return {!B} */A.f = function() { return 1; };
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: B
`).
		Run()
}

// Java: TypeCheckTest.testReturn8()
func (c *TypeCheckTest) TestReturn8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */var A = function() {};
/** @constructor */var B = function() {};
/** @return {!B} */A.prototype.f = function() { return 1; };
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: B
`).
		Run()
}

// Java: TypeCheckTest.testInferredReturn1()
func (c *TypeCheckTest) TestInferredReturn1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f() {} /** @param {number} x */ function g(x) {}
g(f());
`).
		AddDiagnostic(`actual parameter 1 of g does not match formal parameter
found   : undefined
required: number
`).
		Run()
}

// Java: TypeCheckTest.testInferredReturn2()
func (c *TypeCheckTest) TestInferredReturn2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
Foo.prototype.bar = function() {};
/** @param {number} x */ function g(x) {}
g((new Foo()).bar());
`).
		AddDiagnostic(`actual parameter 1 of g does not match formal parameter
found   : undefined
required: number
`).
		Run()
}

// Java: TypeCheckTest.testInferredReturn3()
func (c *TypeCheckTest) TestInferredReturn3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
Foo.prototype.bar = function() {};
/** @constructor
 * @extends {Foo} */ function SubFoo() {}
/** @return {number}
 * @override  */
SubFoo.prototype.bar = function() { return 3; };
`).
		AddDiagnostic(`mismatch of the bar property type and the type of the property it overrides from superclass Foo
original: function(this:Foo): undefined
override: function(this:SubFoo): number
`).
		Run()
}

// Java: TypeCheckTest.testInferredReturn4()
func (c *TypeCheckTest) TestInferredReturn4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var x = function() {};
x = /** @type {function(): number} */ (function() { return 3; });
`).
		AddDiagnostic(`assignment
found   : function(): number
required: function(): undefined
`).
		Run()
}

// Java: TypeCheckTest.testInferredReturn5()
func (c *TypeCheckTest) TestInferredReturn5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {string} */
function f() {
  var x = function() {};
  x = /** @type {function(): number} */ (function() { return 3; });
  return x();
}
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testInferredReturn6()
func (c *TypeCheckTest) TestInferredReturn6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {string} */
function f() {
  var x = function() {};
  if (f())
    x = /** @type {function(): number} */
        (function() { return 3; });
  return x();
}
`).
		AddDiagnostic(`inconsistent return type
found   : (number|undefined)
required: string
`).
		Run()
}

// Java: TypeCheckTest.testInferredReturn7()
func (c *TypeCheckTest) TestInferredReturn7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
/** @param {number} x */ Foo.prototype.bar = function(x) {};
Foo.prototype.bar = function(x) { return 3; };
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: undefined
`).
		Run()
}

// Java: TypeCheckTest.testInferredReturn8()
func (c *TypeCheckTest) TestInferredReturn8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
/** @param {number} x */ Foo.prototype.bar = function(x) {};
/** @constructor
 * @extends {Foo} */ function SubFoo() {}
/** @override @param {number} x */ SubFoo.prototype.bar =
    function(x) { return 3; }
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: undefined
`).
		Run()
}

// Java: TypeCheckTest.testInferredReturn9()
func (c *TypeCheckTest) TestInferredReturn9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {function():string} x */
function f(x) {}
f(/** asdf */ function() { return 123; });
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testInferredParam1()
func (c *TypeCheckTest) TestInferredParam1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
/** @param {number} x */ Foo.prototype.bar = function(x) {};
/** @param {string} x */ function f(x) {}
Foo.prototype.bar = function(y) { f(y); };
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testInferredParam2()
func (c *TypeCheckTest) TestInferredParam2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {string} x */ function f(x) {}
/** @constructor */ function Foo() {}
/** @param {number} x */ Foo.prototype.bar = function(x) {};
/** @constructor
 * @extends {Foo} */ function SubFoo() {}
/** @override @return {void} */ SubFoo.prototype.bar =
    function(x) { f(x); }
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testInferredParam3()
func (c *TypeCheckTest) TestInferredParam3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {string} x */ function f(x) {}
/** @constructor */ function Foo() {}
/** @param {number=} x */ Foo.prototype.bar = function(x) {};
/** @constructor
 * @extends {Foo} */ function SubFoo() {}
/** @override @return {void} */ SubFoo.prototype.bar =
    function(x) { f(x); }; (new SubFoo()).bar();
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : (number|undefined)
required: string
`).
		Run()
}

// Java: TypeCheckTest.testInferredParam4()
func (c *TypeCheckTest) TestInferredParam4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {string} x */ function f(x) {}

/** @constructor */
function Foo() {}
/** @param {...number} x */
Foo.prototype.bar = function(x) {};

/**
 * @constructor
 * @extends {Foo}
 */
function SubFoo() {}
/**
 * @override
 * @return {void}
 */
SubFoo.prototype.bar = function(x) { f(x); };
(new SubFoo()).bar();
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testInferredParam5()
func (c *TypeCheckTest) TestInferredParam5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {string} x */ function f(x) {}
/** @constructor */ function Foo() {}
/** @param {...number} x */ Foo.prototype.bar = function(x) {};
/** @constructor
 * @extends {Foo} */ function SubFoo() {}
/** @override @param {number=} x
 * @param {...number} y  */
SubFoo.prototype.bar =
    function(x, y) { f(x); }; (new SubFoo()).bar();
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : (number|undefined)
required: string
`).
		Run()
}

// Java: TypeCheckTest.testInferredParam6()
func (c *TypeCheckTest) TestInferredParam6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {string} x */ function f(x) {}
/** @constructor */ function Foo() {}
/** @param {number=} x */ Foo.prototype.bar = function(x) {};
/** @constructor
 * @extends {Foo} */ function SubFoo() {}
/** @override @param {number=} x
 * @param {number=} y */
SubFoo.prototype.bar =
    function(x, y) { f(y); };
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : (number|undefined)
required: string
`).
		Run()
}

// Java: TypeCheckTest.testInferredParam7()
func (c *TypeCheckTest) TestInferredParam7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {string} x */ function f(x) {}
/** @type {function(number=,number=)} */
var bar = function(x, y) { f(y); };
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : (number|undefined)
required: string
`).
		Run()
}

// Java: TypeCheckTest.testInferredParam8()
func (c *TypeCheckTest) TestInferredParam8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {function(string)} callback */
function foo(callback) {}
foo(/** random JSDoc */ function(x) { var /** number */ n = x; });
`).
		AddDiagnostic(`initializing variable
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testFunctionLiteralParamWithInlineJSDocNotInferred()
func (c *TypeCheckTest) TestFunctionLiteralParamWithInlineJSDocNotInferred() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {function(string)} x */
function f(x) {}
f(function(/** number */ x) {});
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : function(number): undefined
required: function(string): ?
`).
		Run()
}

// Java: TypeCheckTest.testFunctionLiteralParamWithInlineJSDocNotInferredWithTwoParams()
func (c *TypeCheckTest) TestFunctionLiteralParamWithInlineJSDocNotInferredWithTwoParams() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @param {function(string, number)} x */\nfunction f(x) {}\nf((/** string */ str, num) => {\n// TODO(b/123583824): this should be a type mismatch warning.\n//    found   : number\n//    expected: string\n//  but the JSDoc on `str` blocks inference of `num`.\n  const /** string */ newStr = num;\n});\n").
		Run()
}

// Java: TypeCheckTest.testJSDocOnNoParensArrowFnParameterIsIgnored()
func (c *TypeCheckTest) TestJSDocOnNoParensArrowFnParameterIsIgnored() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @param {function(string)} callback */\nfunction foo(callback) {}\n// The `/** number */` JSDoc is attached to the entire arrow  function, not the\n// parameter `x`\nfoo(/** number */ x =>  { var /** number */ y = x; });\n").
		AddDiagnostic(`initializing variable
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testOverriddenParams1()
func (c *TypeCheckTest) TestOverriddenParams1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
/** @param {...?} var_args */
Foo.prototype.bar = function(var_args) {};
/**
 * @constructor
 * @extends {Foo}
 */ function SubFoo() {}
/**
 * @param {number} x
 * @override
 */
SubFoo.prototype.bar = function(x) {};
`).
		Run()
}

// Java: TypeCheckTest.testOverriddenParams2()
func (c *TypeCheckTest) TestOverriddenParams2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
/** @type {function(...?)} */
Foo.prototype.bar = function(var_args) {};
/**
 * @constructor
 * @extends {Foo}
 */ function SubFoo() {}
/**
 * @type {function(number)}
 * @override
 */
SubFoo.prototype.bar = function(x) {};
`).
		Run()
}

// Java: TypeCheckTest.testOverriddenParams3()
func (c *TypeCheckTest) TestOverriddenParams3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
/** @param {...number} var_args */
Foo.prototype.bar = function(var_args) { };
/**
 * @constructor
 * @extends {Foo}
 */ function SubFoo() {}
/**
 * @param {number} x
 * @override
 */
SubFoo.prototype.bar = function(x) {};
`).
		AddDiagnostic(`mismatch of the bar property type and the type of the property it overrides from superclass Foo
original: function(this:Foo, ...number): undefined
override: function(this:SubFoo, number): undefined
`).
		Run()
}

// Java: TypeCheckTest.testOverriddenParams4()
func (c *TypeCheckTest) TestOverriddenParams4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
/** @type {function(...number)} */
Foo.prototype.bar = function(var_args) {};
/**
 * @constructor
 * @extends {Foo}
 */ function SubFoo() {}
/**
 * @type {function(number)}
 * @override
 */
SubFoo.prototype.bar = function(x) {};
`).
		AddDiagnostic(`mismatch of the bar property type and the type of the property it overrides from superclass Foo
original: function(...number): ?
override: function(number): ?
`).
		Run()
}

// Java: TypeCheckTest.testOverriddenParams5()
func (c *TypeCheckTest) TestOverriddenParams5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
/** @param {number} x */
Foo.prototype.bar = function(x) { };
/**
 * @constructor
 * @extends {Foo}
 */ function SubFoo() {}
/**
 * @override
 */
SubFoo.prototype.bar = function() {};
(new SubFoo()).bar();
`).
		Run()
}

// Java: TypeCheckTest.testOverriddenParams6()
func (c *TypeCheckTest) TestOverriddenParams6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
/** @param {number} x */
Foo.prototype.bar = function(x) { };
/**
 * @constructor
 * @extends {Foo}
 */ function SubFoo() {}
/**
 * @override
 */
SubFoo.prototype.bar = function() {};
(new SubFoo()).bar(true);
`).
		AddDiagnostic(`actual parameter 1 of SubFoo.prototype.bar does not match formal parameter
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckTest.testOverriddenParams7()
func (c *TypeCheckTest) TestOverriddenParams7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ function Foo() {}
/** @param {number} x */
Foo.prototype.bar = function(x) { };
/**
 * @interface
 * @extends {Foo}
 */ function SubFoo() {}
/**
 * @override
 */
SubFoo.prototype.bar = function() {};
var subFoo = /** @type {SubFoo} */ (null);
subFoo.bar(true);
`).
		AddDiagnostic(`actual parameter 1 of SubFoo.prototype.bar does not match formal parameter
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckTest.testOverriddenParams8()
func (c *TypeCheckTest) TestOverriddenParams8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor
 * @template T */ function Foo() {}
/** @param {T} x */
Foo.prototype.bar = function(x) { };
/**
 * @constructor
 * @extends {Foo<string>}
 */ function SubFoo() {}
/**
 * @param {number} x
 * @override
 */
SubFoo.prototype.bar = function(x) {};
`).
		AddDiagnostic(`mismatch of the bar property type and the type of the property it overrides from superclass Foo
original: function(this:Foo, string): undefined
override: function(this:SubFoo, number): undefined
`).
		Run()
}

// Java: TypeCheckTest.testOverriddenReturn1()
func (c *TypeCheckTest) TestOverriddenReturn1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
/** @return {Object} */ Foo.prototype.bar =
    function() { return {}; };
/** @constructor
 * @extends {Foo} */ function SubFoo() {}
/** @return {SubFoo}
 * @override */ SubFoo.prototype.bar =
    function() { return new Foo(); }
`).
		AddDiagnostic(`inconsistent return type
found   : Foo
required: (SubFoo|null)
`).
		Run()
}

// Java: TypeCheckTest.testOverriddenReturn2()
func (c *TypeCheckTest) TestOverriddenReturn2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
/** @return {SubFoo} */ Foo.prototype.bar =
    function() { return new SubFoo(); };
/** @constructor
 * @extends {Foo} */ function SubFoo() {}
/** @return {Foo} x
 * @override */ SubFoo.prototype.bar =
    function() { return new SubFoo(); }
`).
		AddDiagnostic(`mismatch of the bar property type and the type of the property it overrides from superclass Foo
original: function(this:Foo): (SubFoo|null)
override: function(this:SubFoo): (Foo|null)
`).
		Run()
}

// Java: TypeCheckTest.testOverriddenReturn3()
func (c *TypeCheckTest) TestOverriddenReturn3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor
 * @template T */ function Foo() {}
/** @return {T} */ Foo.prototype.bar =
    function() { return null; };
/** @constructor
 * @extends {Foo<string>} */ function SubFoo() {}
/** @override */ SubFoo.prototype.bar =
    function() { return 3; }
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testOverriddenReturn4()
func (c *TypeCheckTest) TestOverriddenReturn4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor
 * @template T */ function Foo() {}
/** @return {T} */ Foo.prototype.bar =
    function() { return null; };
/** @constructor
 * @extends {Foo<string>} */ function SubFoo() {}
/** @return {number}
 * @override */ SubFoo.prototype.bar =
    function() { return 3; }
`).
		AddDiagnostic(`mismatch of the bar property type and the type of the property it overrides from superclass Foo
original: function(this:Foo): string
override: function(this:SubFoo): number
`).
		Run()
}

// Java: TypeCheckTest.testThis1()
func (c *TypeCheckTest) TestThis1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = {};
/** @constructor */goog.A = function(){};
/** @return {number} */
goog.A.prototype.n = function() { return this };
`).
		AddDiagnostic(`inconsistent return type
found   : goog.A
required: number
`).
		Run()
}

// Java: TypeCheckTest.testOverriddenProperty1()
func (c *TypeCheckTest) TestOverriddenProperty1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
/** @type {Object} */
Foo.prototype.bar = {};
/**
 * @constructor
 * @extends {Foo}
 */ function SubFoo() {}
/**
 * @type {Array}
 * @override
 */
SubFoo.prototype.bar = [];
`).
		Run()
}

// Java: TypeCheckTest.testOverriddenProperty2()
func (c *TypeCheckTest) TestOverriddenProperty2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {
  /** @type {Object} */
  this.bar = {};
}
/**
 * @constructor
 * @extends {Foo}
 */ function SubFoo() {}
/**
 * @type {Array}
 * @override
 */
SubFoo.prototype.bar = [];
`).
		Run()
}

// Java: TypeCheckTest.testOverriddenProperty3()
func (c *TypeCheckTest) TestOverriddenProperty3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {
}
/** @type {string} */ Foo.prototype.data;
/**
 * @constructor
 * @extends {Foo}
 */ function SubFoo() {}
/** @type {string|Object}
 * @override */
SubFoo.prototype.data = null;
`).
		AddDiagnostic(`mismatch of the data property type and the type of the property it overrides from superclass Foo
original: string
override: (Object|null|string)
`).
		Run()
}

// Java: TypeCheckTest.testOverriddenProperty4()
func (c *TypeCheckTest) TestOverriddenProperty4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
Foo.prototype.bar = null;
/**
 * @constructor
 * @extends {Foo}
 */ function SubFoo() {}
SubFoo.prototype.bar = 3;
`).
		Run()
}

// Java: TypeCheckTest.testOverriddenProperty5()
func (c *TypeCheckTest) TestOverriddenProperty5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
Foo.prototype.bar = null;
/**
 * @constructor
 * @extends {Foo}
 */ function SubFoo() {}
/** @override */ SubFoo.prototype.bar = 3;
`).
		Run()
}

// Java: TypeCheckTest.testOverriddenProperty6()
func (c *TypeCheckTest) TestOverriddenProperty6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
/** @type {?number} */ Foo.prototype.bar = null;
/**
 * @constructor
 * @extends {Foo}
 */ function SubFoo() {}
SubFoo.prototype.bar = 3;
`).
		Run()
}

// Java: TypeCheckTest.testOverriddenPropertyWithUnknown()
func (c *TypeCheckTest) TestOverriddenPropertyWithUnknown() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
/** @type {?number} */ Foo.prototype.bar = null;

/**
 * @constructor
 * @extends {Foo}
 */
function SubFoo() {}
/** @type {?} */
SubFoo.prototype.bar = 'not a number';
`).
		AddDiagnostic("property bar already defined on superclass Foo; use @override to override it").
		Run()
}

// Java: TypeCheckTest.testOverriddenPropertyWithUnknownInterfaceAncestor()
func (c *TypeCheckTest) TestOverriddenPropertyWithUnknownInterfaceAncestor() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface @extends {Unknown}  */ function Foo() {}
/** @type {string} */ Foo.prototype.data;
/** @constructor @implements {Foo} */ function SubFoo() {}
/** @type {string|Object} */
SubFoo.prototype.data = null;
`).
		AddDiagnostic("Bad type annotation. Unknown type Unknown").
		AddDiagnostic(`mismatch of the data property on type SubFoo and the type of the property it overrides from interface Foo
original: string
override: (Object|null|string)
`).
		Run()
}

// Java: TypeCheckTest.testThis2()
func (c *TypeCheckTest) TestThis2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = {};
/** @constructor */goog.A = function(){
  this.foo = null;
};
/** @return {number} */
goog.A.prototype.n = function() { return this.foo };
`).
		AddDiagnostic(`inconsistent return type
found   : null
required: number
`).
		Run()
}

// Java: TypeCheckTest.testThis3()
func (c *TypeCheckTest) TestThis3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = {};
/** @constructor */goog.A = function(){
  this.foo = null;
  this.foo = 5;
};
`).
		Run()
}

// Java: TypeCheckTest.testThis4()
func (c *TypeCheckTest) TestThis4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = {};
/** @constructor */goog.A = function(){
  /** @type {string?} */this.foo = null;
};
/** @return {number} */goog.A.prototype.n = function() {
  return this.foo };
`).
		AddDiagnostic(`inconsistent return type
found   : (null|string)
required: number
`).
		Run()
}

// Java: TypeCheckTest.testThis5()
func (c *TypeCheckTest) TestThis5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @this {Date}
@return {number}*/function h() { return this }
`).
		AddDiagnostic(`inconsistent return type
found   : Date
required: number
`).
		Run()
}

// Java: TypeCheckTest.testThis6()
func (c *TypeCheckTest) TestThis6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = {};
/** @constructor
@return {!Date} */
goog.A = function(){ return this };
`).
		AddDiagnostic(`inconsistent return type
found   : goog.A
required: Date
`).
		Run()
}

// Java: TypeCheckTest.testThis7()
func (c *TypeCheckTest) TestThis7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */function A(){};
/** @return {number} */A.prototype.n = function() { return this };
`).
		AddDiagnostic(`inconsistent return type
found   : A
required: number
`).
		Run()
}

// Java: TypeCheckTest.testThis8()
func (c *TypeCheckTest) TestThis8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */function A(){
  /** @type {string?} */this.foo = null;
};
/** @return {number} */A.prototype.n = function() {
  return this.foo };
`).
		AddDiagnostic(`inconsistent return type
found   : (null|string)
required: number
`).
		Run()
}

// Java: TypeCheckTest.testTypeOfThis_inStaticMethod_onEs5Ctor_isUnknown()
func (c *TypeCheckTest) TestTypeOfThis_inStaticMethod_onEs5Ctor_isUnknown() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */function A(){};
A.foo = 3;
/** @return {string} */ A.bar = function() { return this.foo; };
`).
		Run()
}

// Java: TypeCheckTest.testThis10()
func (c *TypeCheckTest) TestThis10() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */function A(){};
A.prototype.foo = 3;
/** @this {A}
@return {string} */
A.bar = function() { return this.foo; };
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testThis11()
func (c *TypeCheckTest) TestThis11() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {number} x */ function f(x) {}
/** @constructor */ function Ctor() {
  /** @this {Date} */
  this.method = function() {
    f(this);
  };
}
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : Date
required: number
`).
		Run()
}

// Java: TypeCheckTest.testThis12()
func (c *TypeCheckTest) TestThis12() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {number} x */ function f(x) {}
/** @constructor */ function Ctor() {}
Ctor.prototype['method'] = function() {
  f(this);
}
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : Ctor
required: number
`).
		Run()
}

// Java: TypeCheckTest.testThis13()
func (c *TypeCheckTest) TestThis13() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {number} x */ function f(x) {}
/** @constructor */ function Ctor() {}
Ctor.prototype = {
  method: function() {
    f(this);
  }
};
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : Ctor
required: number
`).
		Run()
}

// Java: TypeCheckTest.testThis14()
func (c *TypeCheckTest) TestThis14() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {number} x */ function f(x) {}
f(this.Object);
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : (typeof Object)
required: number
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testThisType_whichTypesSupport_freeCallsOfFunction()
func (c *TypeCheckTest) TestThisType_whichTypesSupport_freeCallsOfFunction() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {function(this:*)} */ function f() {}; f();").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {function(this:?)} */ function f() {}; f();").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {function(this:undefined)} */ function f() {}; f();").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {function(this:Object)} */ function f() {}; f();").
		Run()
}

// Java: TypeCheckTest.testThisTypeOfFunction2()
func (c *TypeCheckTest) TestThisTypeOfFunction2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function F() {}
/** @type {function(this:F)} */ function f() {}
f();
`).
		AddDiagnostic("\"function(this:F): ?\" must be called with a \"this\" type").
		Run()
}

// Java: TypeCheckTest.testThisTypeOfFunction3()
func (c *TypeCheckTest) TestThisTypeOfFunction3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function F() {}
F.prototype.bar = function() {};
var f = (new F()).bar; f();
`).
		AddDiagnostic("\"function(this:F): undefined\" must be called with a \"this\" type").
		Run()
}

// Java: TypeCheckTest.testThisTypeOfFunction4()
func (c *TypeCheckTest) TestThisTypeOfFunction4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddArray().AddObject().
		Build()).
		AddSource(`/** @constructor */ function F() {}
F.prototype.moveTo = function(x, y) {};
F.prototype.lineTo = function(x, y) {};
function demo() {
  var path = new F();
  var points = [[1,1], [2,2]];
  for (var i = 0; i < points.length; i++) {
    (i == 0 ? path.moveTo : path.lineTo)(points[i][0], points[i][1]);
  }
}
`).
		AddDiagnostic("\"function(this:F, ?, ?): undefined\" must be called with a \"this\" type").
		Run()
}

// Java: TypeCheckTest.testThisTypeOfFunction5()
func (c *TypeCheckTest) TestThisTypeOfFunction5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {function(this:number)} */
function f() {
  var /** number */ n = this;
}
`).
		Run()
}

// Java: TypeCheckTest.testGlobalThis1()
func (c *TypeCheckTest) TestGlobalThis1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Window() {}
/** @param {string} msg */
Window.prototype.alert = function(msg) {};
this.alert(3);
`).
		AddDiagnostic(`actual parameter 1 of Window.prototype.alert does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testGlobalThis2a()
func (c *TypeCheckTest) TestGlobalThis2a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Bindow() {}
/** @param {string} msg */
Bindow.prototype.alert = function(msg) {};
this.alert = 3;
(new Bindow()).alert(this.alert)
`).
		AddDiagnostic("Property alert never defined on global this").
		AddDiagnostic("Property alert never defined on global this").
		Run()
}

// Java: TypeCheckTest.testGlobalThis2b()
func (c *TypeCheckTest) TestGlobalThis2b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @constructor */ function Bindow() {}
/** @param {string} msg */
Bindow.prototype.alert = function(msg) {};
this.alert = 3;
(new Bindow()).alert(this.alert)
`).
		Run()
}

// Java: TypeCheckTest.testGlobalThis2c()
func (c *TypeCheckTest) TestGlobalThis2c() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @constructor */ function Bindow() {}
/** @param {string} msg */
Bindow.prototype.alert = function(msg) {};
/** @return {number} */ this.alert = function() { return 3; };
(new Bindow()).alert(this.alert())
`).
		AddDiagnostic(`actual parameter 1 of Bindow.prototype.alert does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testGlobalThis3()
func (c *TypeCheckTest) TestGlobalThis3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {string} msg */
function alert(msg) {};
this.alert(3);
`).
		AddDiagnostic(`actual parameter 1 of global this.alert does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testGlobalThis4()
func (c *TypeCheckTest) TestGlobalThis4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {string} msg */
var alert = function(msg) {};
this.alert(3);
`).
		AddDiagnostic(`actual parameter 1 of global this.alert does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testGlobalThis5()
func (c *TypeCheckTest) TestGlobalThis5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f() {
  /** @param {string} msg */
  var alert = function(msg) {};
}
this.alert(3);
`).
		AddDiagnostic("Property alert never defined on global this").
		Run()
}

// Java: TypeCheckTest.testGlobalThis6()
func (c *TypeCheckTest) TestGlobalThis6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {string} msg */
var alert = function(msg) {};
var x = 3;
x = 'msg';
this.alert(this.x);
`).
		Run()
}

// Java: TypeCheckTest.testGlobalThis7()
func (c *TypeCheckTest) TestGlobalThis7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Window() {}
/** @param {Window} msg */
var foo = function(msg) {};
foo(this);
`).
		Run()
}

// Java: TypeCheckTest.testGlobalThis8()
func (c *TypeCheckTest) TestGlobalThis8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Window() {}
/** @param {number} msg */
var foo = function(msg) {};
foo(this);
`).
		AddDiagnostic(`actual parameter 1 of foo does not match formal parameter
found   : global this
required: number
`).
		Run()
}

// Java: TypeCheckTest.testGlobalThis9()
func (c *TypeCheckTest) TestGlobalThis9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`// Window is not marked as a constructor, so the
// inheritance doesn't happen.
function Window() {}
Window.prototype.alert = function() {};
this.alert();
`).
		AddDiagnostic("Property alert never defined on global this").
		Run()
}

// Java: TypeCheckTest.testGlobalThisDoesNotIncludeVarsDeclaredWithConst()
func (c *TypeCheckTest) TestGlobalThisDoesNotIncludeVarsDeclaredWithConst() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {string} msg */
const alert = function(msg) {};
this.alert('boo');
`).
		AddDiagnostic("Property alert never defined on global this").
		Run()
}

// Java: TypeCheckTest.testGlobalThisDoesNotIncludeVarsDeclaredWithLet()
func (c *TypeCheckTest) TestGlobalThisDoesNotIncludeVarsDeclaredWithLet() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {string} msg */
let alert = function(msg) {};
this.alert('boo');
`).
		AddDiagnostic("Property alert never defined on global this").
		Run()
}

// Java: TypeCheckTest.testControlFlowRestrictsType1()
func (c *TypeCheckTest) TestControlFlowRestrictsType1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {String?} */ function f() { return null; }
/** @type {String?} */ var a = f();
/** @type {String} */ var b = new String('foo');
/** @type {null} */ var c = null;
if (a) {
  b = a;
} else {
  c = a;
}
`).
		Run()
}

// Java: TypeCheckTest.testControlFlowRestrictsType2()
func (c *TypeCheckTest) TestControlFlowRestrictsType2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {(string|null)} */ function f() { return null; }
/** @type {(string|null)} */ var a = f();
/** @type {string} */ var b = 'foo';
/** @type {null} */ var c = null;
if (a) {
  b = a;
} else {
  c = a;
}
`).
		AddDiagnostic(`assignment
found   : (null|string)
required: null
`).
		Run()
}

// Java: TypeCheckTest.testControlFlowRestrictsType3()
func (c *TypeCheckTest) TestControlFlowRestrictsType3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {(string|void)} */
var a;
/** @type {string} */
var b = 'foo';
if (a) {
  b = a;
}
`).
		Run()
}

// Java: TypeCheckTest.testControlFlowRestrictsType4()
func (c *TypeCheckTest) TestControlFlowRestrictsType4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {string} a */ function f(a){}
/** @type {(string|undefined)} */ var a;
a && f(a);
`).
		Run()
}

// Java: TypeCheckTest.testControlFlowRestrictsType5()
func (c *TypeCheckTest) TestControlFlowRestrictsType5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {undefined} a */ function f(a){}
/** @type {(!Array|undefined)} */ var a;
a || f(a);
`).
		Run()
}

// Java: TypeCheckTest.testControlFlowRestrictsType6()
func (c *TypeCheckTest) TestControlFlowRestrictsType6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {undefined} x */ function f(x) {}
/** @type {(string|undefined)} */ var a;
a && f(a);
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : string
required: undefined
`).
		Run()
}

// Java: TypeCheckTest.testControlFlowRestrictsType7()
func (c *TypeCheckTest) TestControlFlowRestrictsType7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {undefined} x */ function f(x) {}
/** @type {(string|undefined)} */ var a;
a && f(a);
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : string
required: undefined
`).
		Run()
}

// Java: TypeCheckTest.testControlFlowRestrictsType8()
func (c *TypeCheckTest) TestControlFlowRestrictsType8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {undefined} a */ function f(a){}
/** @type {(!Array|undefined)} */ var a;
if (a || f(a)) {}
`).
		Run()
}

// Java: TypeCheckTest.testControlFlowRestrictsType9()
func (c *TypeCheckTest) TestControlFlowRestrictsType9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {number?} x
 * @return {number}*/
var f = function(x) {
if (!x || x == 1) { return 1; } else { return x; }
};
`).
		Run()
}

// Java: TypeCheckTest.testControlFlowRestrictsType10()
func (c *TypeCheckTest) TestControlFlowRestrictsType10() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {number} x */ function f(x) {}
function g() {
  var y = null;
  for (var i = 0; i < 10; i++) {
    f(y);
    if (y != null) {
      // y is None the first time it goes through this branch
    } else {
      y = {};
    }
  }
};
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : (null|{})
required: number
`).
		Run()
}

// Java: TypeCheckTest.testControlFlowRestrictsType11()
func (c *TypeCheckTest) TestControlFlowRestrictsType11() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {boolean} x */ function f(x) {}
function g() {
  var y = null;
  if (y != null) {
    for (var i = 0; i < 10; i++) {
      f(y);
    }
  }
};
`).
		AddDiagnostic(`condition always evaluates to false
left : null
right: null
`).
		Run()
}

// Java: TypeCheckTest.testSwitchCase_primitiveDoesNotAutobox()
func (c *TypeCheckTest) TestSwitchCase_primitiveDoesNotAutobox() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {!String} */
var a = new String('foo');
switch (a) { case 'A': }
`).
		AddDiagnostic(`case expression doesn't match switch
found   : string
required: String
`).
		Run()
}

// Java: TypeCheckTest.testSwitchCase_unknownSwitchExprMatchesAnyCase()
func (c *TypeCheckTest) TestSwitchCase_unknownSwitchExprMatchesAnyCase() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var a = unknown;
switch (a) { case 'A':break; case null:break; }
`).
		Run()
}

// Java: TypeCheckTest.testSwitchCase_doesNotAutoboxStringToMatchNullableUnion()
func (c *TypeCheckTest) TestSwitchCase_doesNotAutoboxStringToMatchNullableUnion() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {?String} */
var a = unknown;
switch (a) { case 'A':break; case null:break; }
`).
		AddDiagnostic(`case expression doesn't match switch
found   : string
required: (String|null)
`).
		Run()
}

// Java: TypeCheckTest.testSwitchCase_doesNotAutoboxNumberToMatchNullableUnion()
func (c *TypeCheckTest) TestSwitchCase_doesNotAutoboxNumberToMatchNullableUnion() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {?Number} */
var a = unknown;
switch (a) { case 5:break; case null:break; }
`).
		AddDiagnostic(`case expression doesn't match switch
found   : number
required: (Number|null)
`).
		Run()
}

// Java: TypeCheckTest.testSwitchCase7()
func (c *TypeCheckTest) TestSwitchCase7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {number} x
 * @return {number}
 */
function g(x) { return 5; }
function f() {
  var x = {};
  x.foo = '3';
  switch (3) { case g(x.foo): return 3; }
}
`).
		AddDiagnostic(`actual parameter 1 of g does not match formal parameter
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testSwitchCase8()
func (c *TypeCheckTest) TestSwitchCase8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {number} x
 * @return {number}
 */
function g(x) { return 5; }
function f() {
  var x = {};
  x.foo = '3';
  switch (g(x.foo)) { case 3: return 3; }
}
`).
		AddDiagnostic(`actual parameter 1 of g does not match formal parameter
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testSwitchCase_allowsStructuralMatching()
func (c *TypeCheckTest) TestSwitchCase_allowsStructuralMatching() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @record */
class R {
  constructor() {
    /** @type {string} */
    this.str;
  }
}
/** @record */
class S {
  constructor() {
    /** @type {string} */
    this.str;
    /** @type {number} */
    this.num;
  }
}
 /**
 * @param {!R} r
 * @param {!S} s
 */
function f(r, s) {
  switch (r) {
    case s:
      return true;
  }
}
`).
		Run()
}

// Java: TypeCheckTest.testImplicitCast1()
func (c *TypeCheckTest) TestImplicitCast1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @constructor */ function Element() {};
/** @type {string}
  * @implicitCast */
Element.prototype.innerHTML;
`).
		AddSource("(new Element).innerHTML = new Array();").
		Run()
}

// Java: TypeCheckTest.testImplicitCast2()
func (c *TypeCheckTest) TestImplicitCast2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @constructor */ function Element() {};
/**
 * @type {string}
 * @implicitCast
 */
Element.prototype.innerHTML;
`).
		AddSource(`/** @constructor */ function C(e) {
  /** @type {Element} */ this.el = e;
}
C.prototype.method = function() {
  this.el.innerHTML = new Array();
};
`).
		Run()
}

// Java: TypeCheckTest.testImplicitCast3()
func (c *TypeCheckTest) TestImplicitCast3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @constructor */ function Element() {};
/**
 * @type {string}
 * @implicitCast
 */
Element.prototype.innerHTML;
`).
		AddSource(`/** @param {?Element} element
 * @param {string|number} text
 */
function f(element, text) {
  element.innerHTML = text;
}
`).
		Run()
}

// Java: TypeCheckTest.testImplicitCastSubclassAccess()
func (c *TypeCheckTest) TestImplicitCastSubclassAccess() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @constructor */ function Element() {};
/** @type {string}
  * @implicitCast */
Element.prototype.innerHTML;
/** @constructor
 * @extends Element */
function DIVElement() {};
`).
		AddSource("(new DIVElement).innerHTML = new Array();").
		Run()
}

// Java: TypeCheckTest.testImplicitCastNotInExterns()
func (c *TypeCheckTest) TestImplicitCastNotInExterns() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Element() {};
/**
 * @type {string}
 * @implicitCast
 */
Element.prototype.innerHTML;
(new Element).innerHTML = new Array();
`).
		AddDiagnostic(`assignment to property innerHTML of Element
found   : Array<?>
required: string
`).
		Run()
}

// Java: TypeCheckTest.testNumberNode()
func (c *TypeCheckTest) TestNumberNode() {
	var n *rhino.Node = rhino.IR_Number(0)
	c.TypeCheck(rhino.IR_ExprResult(n))
	c.AssertTypeEquals(c.GetNativeNumberType(), n.GetJSType())
}

// Java: TypeCheckTest.testStringNode()
func (c *TypeCheckTest) TestStringNode() {
	var n *rhino.Node = rhino.IR_String("hello")
	c.TypeCheck(rhino.IR_ExprResult(n))
	c.AssertTypeEquals(c.GetNativeStringType(), n.GetJSType())
}

// Java: TypeCheckTest.testBooleanNodeTrue()
func (c *TypeCheckTest) TestBooleanNodeTrue() {
	var trueNode *rhino.Node = rhino.IR_TrueNode()
	c.TypeCheck(rhino.IR_ExprResult(trueNode))
	c.AssertTypeEquals(c.GetNativeBooleanType(), trueNode.GetJSType())
}

// Java: TypeCheckTest.testBooleanNodeFalse()
func (c *TypeCheckTest) TestBooleanNodeFalse() {
	var falseNode *rhino.Node = rhino.IR_FalseNode()
	c.TypeCheck(rhino.IR_ExprResult(falseNode))
	c.AssertTypeEquals(c.GetNativeBooleanType(), falseNode.GetJSType())
}

// Java: TypeCheckTest.testUndefinedNode()
func (c *TypeCheckTest) TestUndefinedNode() {
	var p *rhino.Node = rhino.NewNode(rhino.Token_ADD)
	var n *rhino.Node = rhino.Node_NewStringWithToken(rhino.Token_NAME, "undefined")
	p.AddChildToBack(n)
	p.AddChildToBack(rhino.Node_NewNumber(5))
	c.TypeCheck(rhino.IR_ExprResult(p))
	c.AssertTypeEquals(c.GetNativeVoidType(), n.GetJSType())
}

// Java: TypeCheckTest.testNumberAutoboxing()
func (c *TypeCheckTest) TestNumberAutoboxing() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {Number} */var a = 4;").
		AddDiagnostic(`initializing variable
found   : number
required: (Number|null)
`).
		Run()
}

// Java: TypeCheckTest.testNumberUnboxing()
func (c *TypeCheckTest) TestNumberUnboxing() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {number} */var a = new Number(4);").
		AddDiagnostic(`initializing variable
found   : Number
required: number
`).
		Run()
}

// Java: TypeCheckTest.testStringAutoboxing()
func (c *TypeCheckTest) TestStringAutoboxing() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {String} */var a = 'hello';").
		AddDiagnostic(`initializing variable
found   : string
required: (String|null)
`).
		Run()
}

// Java: TypeCheckTest.testStringUnboxing()
func (c *TypeCheckTest) TestStringUnboxing() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {string} */var a = new String('hello');").
		AddDiagnostic(`initializing variable
found   : String
required: string
`).
		Run()
}

// Java: TypeCheckTest.testBooleanAutoboxing()
func (c *TypeCheckTest) TestBooleanAutoboxing() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {Boolean} */var a = true;").
		AddDiagnostic(`initializing variable
found   : boolean
required: (Boolean|null)
`).
		Run()
}

// Java: TypeCheckTest.testBooleanUnboxing()
func (c *TypeCheckTest) TestBooleanUnboxing() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {boolean} */var a = new Boolean(false);").
		AddDiagnostic(`initializing variable
found   : Boolean
required: boolean
`).
		Run()
}

// Java: TypeCheckTest.testIIFE1()
func (c *TypeCheckTest) TestIIFE1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var namespace = {};
/** @type {number} */ namespace.prop = 3;
(function(ns) {
  ns.prop = true;
})(namespace);
`).
		AddDiagnostic(`assignment to property prop of ns
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckTest.testIIFE2()
func (c *TypeCheckTest) TestIIFE2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
(function(ctor) {
  /** @type {boolean} */ ctor.prop = true;
})(Foo);
/** @return {number} */ function f() { return Foo.prop; }
`).
		AddDiagnostic(`inconsistent return type
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckTest.testIIFE3()
func (c *TypeCheckTest) TestIIFE3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
(function(ctor) {
  /** @type {boolean} */ ctor.prop = true;
})(Foo);
/** @param {number} x */ function f(x) {}
f(Foo.prop);
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckTest.testIIFE4()
func (c *TypeCheckTest) TestIIFE4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const */ var namespace = {};
(function(ns) {
  /**
   * @constructor
   * @param {number} x
   */
   ns.Ctor = function(x) {};
})(namespace);
new namespace.Ctor(true);
`).
		AddDiagnostic(`actual parameter 1 of namespace.Ctor does not match formal parameter
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckTest.testIIFE5()
func (c *TypeCheckTest) TestIIFE5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const */ var namespace = {};
(function(ns) {
  /**
   * @constructor
   */
   ns.Ctor = function() {};
   /** @type {boolean} */ ns.Ctor.prototype.bar = true;
})(namespace);
/** @param {namespace.Ctor} x
  * @return {number} */ function f(x) { return x.bar; }
`).
		AddDiagnostic(`inconsistent return type
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckTest.testNotIIFE1()
func (c *TypeCheckTest) TestNotIIFE1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {number} x */ function f(x) {}
/** @param {...?} x */ function g(x) {}
g(function(y) { f(y); }, true);
`).
		Run()
}

// Java: TypeCheckTest.testEnums()
func (c *TypeCheckTest) TestEnums() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var outer = function() {
  /** @enum {number} */
  var Level = {
    NONE: 0,
  };
  /** @type {!Level} */
  var l = Level.NONE;
}
`).
		Run()
}

// Java: TypeCheckTest.testStrictInterfaceCheck()
func (c *TypeCheckTest) TestStrictInterfaceCheck() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
function EventTarget() {}
/** @constructor
 * @implements {EventTarget} */
function Node() {}
/** @type {number} */ Node.prototype.index;
/** @param {EventTarget} x
 * @return {string} */
function foo(x) { return x.index; }
`).
		AddDiagnostic("Property index never defined on EventTarget").
		Run()
}

// Java: TypeCheckTest.testTemplateSubtyping_0()
func (c *TypeCheckTest) TestTemplateSubtyping_0() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @implements {IFoo<number>} */
class Foo { }

/**
 * @template T
 * @interface
 */
class IFoo { }

const /** !IFoo<number> */ x = new Foo();
`).
		Run()
}

// Java: TypeCheckTest.testTemplateSubtyping_1()
func (c *TypeCheckTest) TestTemplateSubtyping_1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @template T
 * @interface
 */
class IFoo { }

/** @implements {IFoo<number>} */
class FooA { }

/** @implements {IFoo<string>} */
class FooB extends FooA { }

const /** !IFoo<number> */ x = new FooB();
const /** !IFoo<string> */ y = new FooB();
`).
		Run()
}

// Java: TypeCheckTest.testTypedefBeforeUse()
func (c *TypeCheckTest) TestTypedefBeforeUse() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @typedef {Object<string, number>} */
var map;
/** @param {(map|function())} isResult */
var f = function(isResult) {
    while (true)
        isResult['t'];
};
`).
		Run()
}

// Java: TypeCheckTest.testScopedConstructors1()
func (c *TypeCheckTest) TestScopedConstructors1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function foo1() {
  /** @constructor */ function Bar() {
    /** @type {number} */ this.x = 3;
  }
}
function foo2() {
  /** @constructor */ function Bar() {
    /** @type {string} */ this.x = 'y';
  }
  /** * @param {Bar} b
   * @return {number}
   */
  function baz(b) { return b.x; }
}
`).
		AddDiagnostic(`inconsistent return type
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testScopedConstructors2()
func (c *TypeCheckTest) TestScopedConstructors2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {Function} f */
function foo1(f) {
  /** @param {Function} g */
  f.prototype.bar = function(g) {};
}
`).
		Run()
}

// Java: TypeCheckTest.testQualifiedNameInference1()
func (c *TypeCheckTest) TestQualifiedNameInference1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
/** @type {number?} */ Foo.prototype.bar = null;
/** @type {number?} */ Foo.prototype.baz = null;
/** @param {Foo} foo */
function f(foo) {
  while (true) {
    if (!foo.baz) break;
    foo.bar = null;
  }
  // Tests a bug where this condition always evaluated to true.
  return foo.bar == null;
}
`).
		Run()
}

// Java: TypeCheckTest.testQualifiedNameInference2()
func (c *TypeCheckTest) TestQualifiedNameInference2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var x = {};
x.y = c;
function f(a, b) {
  if (a) {
    if (b)
      x.y = 2;
    else
      x.y = 1;
  }
  return x.y == null;
}
`).
		Run()
}

// Java: TypeCheckTest.testQualifiedNameInference3()
func (c *TypeCheckTest) TestQualifiedNameInference3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var x = {};
x.y = c;
function f(a, b) {
  if (a) {
    if (b)
      x.y = 2;
    else
      x.y = 1;
  }
  return x.y == null;
} function g() { x.y = null; }
`).
		Run()
}

// Java: TypeCheckTest.testQualifiedNameInference4()
func (c *TypeCheckTest) TestQualifiedNameInference4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {string} x */ function f(x) {}
/**
 * @param {?string} x
 * @constructor
 */
function Foo(x) { this.x_ = x; }
Foo.prototype.bar = function() {
  if (this.x_) { f(this.x_); }
};
`).
		Run()
}

// Java: TypeCheckTest.testQualifiedNameInference5()
func (c *TypeCheckTest) TestQualifiedNameInference5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`var ns = {};
(function() {
    /** @param {number} x */ ns.foo = function(x) {}; })();
(function() { ns.foo(true); })();
`).
		AddDiagnostic(`actual parameter 1 of ns.foo does not match formal parameter
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckTest.testQualifiedNameInference6()
func (c *TypeCheckTest) TestQualifiedNameInference6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const */ var ns = {};
/** @param {number} x */ ns.foo = function(x) {};
(function() {
    ns.foo = function(x) {};
    ns.foo(true);
})();
`).
		AddDiagnostic(`actual parameter 1 of ns.foo does not match formal parameter
found   : boolean
required: number`).
		Run()
}

// Java: TypeCheckTest.testQualifiedNameInference7()
func (c *TypeCheckTest) TestQualifiedNameInference7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`var ns = {};
(function() {
  /** @constructor
  * @param {number} x */
  ns.Foo = function(x) {};
  /** @param {ns.Foo} x */ function f(x) {}
  f(new ns.Foo(true));
})();
`).
		AddDiagnostic(`actual parameter 1 of ns.Foo does not match formal parameter
found   : boolean
required: number`).
		Run()
}

// Java: TypeCheckTest.testQualifiedNameInference8()
func (c *TypeCheckTest) TestQualifiedNameInference8() {
	c.TestClosureTypesMultipleWarnings(`var ns = {};
(function() {
  /** @constructor
  * @param {number} x */
  ns.Foo = function(x) {};
})();
/** @param {ns.Foo} x */ function f(x) {}
f(new ns.Foo(true));
`, []string{"Property Foo never defined on ns", `actual parameter 1 of ns.Foo does not match formal parameter
found   : boolean
required: number`})
}

// Java: TypeCheckTest.testQualifiedNameInference9()
func (c *TypeCheckTest) TestQualifiedNameInference9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var ns = {};
ns.ns2 = {};
(function() {
  /** @constructor
  * @param {number} x */
  ns.ns2.Foo = function(x) {};
  /** @param {ns.ns2.Foo} x */ function f(x) {}
  f(new ns.ns2.Foo(true));
})();
`).
		AddDiagnostic(`actual parameter 1 of ns.ns2.Foo does not match formal parameter
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckTest.testQualifiedNameInference10()
func (c *TypeCheckTest) TestQualifiedNameInference10() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var ns = {};
ns.ns2 = {};
(function() {
  /** @interface */
  ns.ns2.Foo = function() {};
  /** @constructor
  * @implements {ns.ns2.Foo} */
  function F() {}
  (new F());
})();
`).
		Run()
}

// Java: TypeCheckTest.testQualifiedNameInference11()
func (c *TypeCheckTest) TestQualifiedNameInference11() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @constructor */ function Foo() {}
function f() {
  var x = new Foo();
  x.onload = function() {
    x.onload = null;
  };
}
`).
		Run()
}

// Java: TypeCheckTest.testQualifiedNameInference12()
func (c *TypeCheckTest) TestQualifiedNameInference12() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {function(this:Object)} x */ function f(x) {}
/** @constructor */ function Foo() {
  /** @type {number} */ this.bar = 3;
  f(function() { this.bar = true; });
}
`).
		Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		Run()
}

// Java: TypeCheckTest.testQualifiedNameInference13()
func (c *TypeCheckTest) TestQualifiedNameInference13() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @constructor */ function Foo() {}
function f(z) {
  var x = new Foo();
  if (z) {
    x.onload = function() {};
  } else {
    x.onload = null;
  };
}
`).
		Run()
}

// Java: TypeCheckTest.testQualifiedNameInference14()
func (c *TypeCheckTest) TestQualifiedNameInference14() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @constructor */ function Foo() {}
function f(z) {
  var x = new Foo();
  {
    x.onload = function() {};
  }
  {
    x.onload = null;
  };
}
`).
		AddDiagnostic(`assignment
found   : null
required: function(): undefined
`).
		Run()
}

// Java: TypeCheckTest.testScopeQualifiedNamesOnThis()
func (c *TypeCheckTest) TestScopeQualifiedNamesOnThis() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {
  /** @type {!Bar} */
  this.baz = new Bar();
}
Foo.prototype.foo = function() {
  this.baz.bar();
};
/** @constructor */ function Bar() {
  /** @type {!Foo} */
  this.baz = new Foo();
}
Bar.prototype.bar = function() {
  this.baz.foo();
};
`).
		Run()
}

// Java: TypeCheckTest.testSheqRefinedScope()
func (c *TypeCheckTest) TestSheqRefinedScope() {
	var n *rhino.Node = c.ParseAndTypeCheck(`/** @constructor */function A() {}
/** @constructor 
 @extends A */ function B() {}
/** @return {number} */
B.prototype.p = function() { return 1; }
/** @param {A} a
 @param {B} b */
function f(a, b) {
  b.p();
  if (a === b) {
    b.p();
  }
}
`)
	var nodeC *rhino.Node = n.GetLastChild().GetLastChild().
		GetLastChild().
		GetLastChild().
		GetLastChild().
		GetLastChild()
	var typeC rhino.JSType = nodeC.GetJSType()
	AssertThat(typeC.IsNumber()).IsTrue()
	var nodeB *rhino.Node = nodeC.GetFirstFirstChild()
	var typeB rhino.JSType = nodeB.GetJSType()
	AssertThat(typeB.String()).IsEqualTo("B")
}

// Java: TypeCheckTest.testAssignToUntypedVariable()
func (c *TypeCheckTest) TestAssignToUntypedVariable() {
	var n *rhino.Node = c.ParseAndTypeCheck("var z; z = 1;")
	var assign *rhino.Node = n.GetLastChild().GetFirstChild()
	var node *rhino.Node = assign.GetFirstChild()
	AssertThat(node.GetJSType().IsUnknownType()).IsFalse()
	AssertThat(node.GetJSType().String()).IsEqualTo("number")
}

// Java: TypeCheckTest.testAssignToUntypedProperty()
func (c *TypeCheckTest) TestAssignToUntypedProperty() {
	var n *rhino.Node = c.ParseAndTypeCheck(`/** @constructor */ function Foo() {}
Foo.prototype.a = 1;
(new Foo).a;
`)
	var node *rhino.Node = n.GetLastChild().GetFirstChild()
	AssertThat(node.GetJSType().IsUnknownType()).IsFalse()
	AssertThat(node.GetJSType().IsNumber()).IsTrue()
}

// Java: TypeCheckTest.testNew1()
func (c *TypeCheckTest) TestNew1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("new 4").
		AddDiagnostic("cannot instantiate non-constructor, found type: number").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("new 4").
		AddDiagnostic("cannot instantiate non-constructor, found type: number").
		Run()
}

// Java: TypeCheckTest.testNew2()
func (c *TypeCheckTest) TestNew2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var Math = {}; new Math()").
		AddDiagnostic("cannot instantiate non-constructor, found type: {}").
		Run()
}

// Java: TypeCheckTest.testNew3()
func (c *TypeCheckTest) TestNew3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("new Date()").
		Run()
}

// Java: TypeCheckTest.testNew4()
func (c *TypeCheckTest) TestNew4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @constructor */function A(){}; new A();").
		Run()
}

// Java: TypeCheckTest.testNew5()
func (c *TypeCheckTest) TestNew5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function A(){}; new A();").
		AddDiagnostic("cannot instantiate non-constructor, found type: function(): undefined").
		Run()
}

// Java: TypeCheckTest.testNew6()
func (c *TypeCheckTest) TestNew6() {
	var p *TypeCheckTestCase_TypeCheckResult = c.ParseAndTypeCheckWithScope(`/** @constructor */function A(){};
var a = new A();
`)
	var aType rhino.JSType = p.scope.GetVar("a").GetType()
	IsInstanceOf[rhino.ObjectType](AssertThat(aType))
	var aObjectType rhino.ObjectType = aType.(rhino.ObjectType)
	AssertThat(aObjectType.GetConstructor().GetReferenceName()).IsEqualTo("A")
}

// Java: TypeCheckTest.testNew7()
func (c *TypeCheckTest) TestNew7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {Function} opt_constructor */
function foo(opt_constructor) {
if (opt_constructor) { new opt_constructor; }
}
`).
		Run()
}

// Java: TypeCheckTest.testNew8()
func (c *TypeCheckTest) TestNew8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {Function} opt_constructor */
function foo(opt_constructor) {
new opt_constructor;
}
`).
		Run()
}

// Java: TypeCheckTest.testNew9()
func (c *TypeCheckTest) TestNew9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {Function} opt_constructor */
function foo(opt_constructor) {
new (opt_constructor || Array);
}
`).
		Run()
}

// Java: TypeCheckTest.testNew10()
func (c *TypeCheckTest) TestNew10() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = {};
/** @param {Function} opt_constructor */
goog.Foo = function(opt_constructor) {
new (opt_constructor || Array);
}
`).
		Run()
}

// Java: TypeCheckTest.testNew11()
func (c *TypeCheckTest) TestNew11() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {Function} c1 */
function f(c1) {
  var c2 = function(){};
  c1.prototype = new c2;
}
`).
		AddDiagnostic("cannot instantiate non-constructor, found type: function(): undefined").
		Run()
}

// Java: TypeCheckTest.testNew12()
func (c *TypeCheckTest) TestNew12() {
	var p *TypeCheckTestCase_TypeCheckResult = c.ParseAndTypeCheckWithScope("var a = new Array();")
	var a *TypedVar = p.scope.GetVar("a")
	c.AssertTypeEquals(c.GetNativeArrayType(), a.GetType())
}

// Java: TypeCheckTest.testNew13()
func (c *TypeCheckTest) TestNew13() {
	var p *TypeCheckTestCase_TypeCheckResult = c.ParseAndTypeCheckWithScope(`/** @constructor */function FooBar(){};
var a = new FooBar();
`)
	var a *TypedVar = p.scope.GetVar("a")
	IsInstanceOf[rhino.ObjectType](AssertThat(a.GetType()))
	AssertThat(a.GetType().String()).IsEqualTo("FooBar")
}

// Java: TypeCheckTest.testNew14()
func (c *TypeCheckTest) TestNew14() {
	var p *TypeCheckTestCase_TypeCheckResult = c.ParseAndTypeCheckWithScope(`/** @constructor */var FooBar = function(){};
var a = new FooBar();
`)
	var a *TypedVar = p.scope.GetVar("a")
	IsInstanceOf[rhino.ObjectType](AssertThat(a.GetType()))
	AssertThat(a.GetType().String()).IsEqualTo("FooBar")
}

// Java: TypeCheckTest.testNew15()
func (c *TypeCheckTest) TestNew15() {
	var p *TypeCheckTestCase_TypeCheckResult = c.ParseAndTypeCheckWithScope(`var goog = {};
/** @constructor */goog.A = function(){};
var a = new goog.A();
`)
	var a *TypedVar = p.scope.GetVar("a")
	IsInstanceOf[rhino.ObjectType](AssertThat(a.GetType()))
	AssertThat(a.GetType().String()).IsEqualTo("goog.A")
}

// Java: TypeCheckTest.testNew16()
func (c *TypeCheckTest) TestNew16() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** * @param {string} x
 * @constructor
 */
function Foo(x) {}
function g() { new Foo(1); }
`).
		AddDiagnostic(`actual parameter 1 of Foo does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testNew17()
func (c *TypeCheckTest) TestNew17() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var goog = {}; goog.x = 3; new goog.x").
		AddDiagnostic("cannot instantiate non-constructor, found type: number").
		Run()
}

// Java: TypeCheckTest.testNew18()
func (c *TypeCheckTest) TestNew18() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = {};
/** @constructor */ goog.F = function() {};
/** @constructor */ goog.G = goog.F;
`).
		Run()
}

// Java: TypeCheckTest.testNew19()
func (c *TypeCheckTest) TestNew19() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @constructor @abstract */ var Foo = function() {}; var foo = new Foo();").
		AddDiagnosticType(TypeCheck_INSTANTIATE_ABSTRACT_CLASS).
		Run()
}

// Java: TypeCheckTest.testNew20()
func (c *TypeCheckTest) TestNew20() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor @abstract */
function Bar() {};
/** @return {function(new:Bar)} */
function foo() {}
var Foo = foo();
var f = new Foo;
`).
		Run()
}

// Java: TypeCheckTest.testName1()
func (c *TypeCheckTest) TestName1() {
	c.AssertTypeEquals(c.GetNativeVoidType(), c.TestNameNode("undefined"))
}

// Java: TypeCheckTest.testName2()
func (c *TypeCheckTest) TestName2() {
	c.AssertTypeEquals(c.GetNativeObjectConstructorType(), c.TestNameNode("Object"))
}

// Java: TypeCheckTest.testName3()
func (c *TypeCheckTest) TestName3() {
	c.AssertTypeEquals(c.GetNativeArrayConstructorType(), c.TestNameNode("Array"))
}

// Java: TypeCheckTest.testName4()
func (c *TypeCheckTest) TestName4() {
	c.AssertTypeEquals(c.GetNativeDateConstructorType(), c.TestNameNode("Date"))
}

// Java: TypeCheckTest.testName5()
func (c *TypeCheckTest) TestName5() {
	c.AssertTypeEquals(c.GetNativeRegexpConstructorType(), c.TestNameNode("RegExp"))
}

// Java: TypeCheckTest.testNameNode(String)
func (c *TypeCheckTest) TestNameNode(name string) rhino.JSType {
	var node *rhino.Node = rhino.Node_NewStringWithToken(rhino.Token_NAME, name)
	var parent *rhino.Node = rhino.NewNodeWithChild(rhino.Token_SCRIPT, node)
	parent.SetInputId(rhino.NewInputId("code"))
	var externs *rhino.Node = rhino.NewNode(rhino.Token_SCRIPT)
	externs.SetInputId(rhino.NewInputId("externs"))
	var root *rhino.Node = rhino.IR_Root(rhino.IR_Root(externs), rhino.IR_Root(parent))
	c.MakeTypeCheck().ProcessForTesting(root.GetFirstChild(), root.GetSecondChild())
	return node.GetJSType()
}

// Java: TypeCheckTest.testBitOperation1()
func (c *TypeCheckTest) TestBitOperation1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@return {void}*/function foo(){ ~foo(); }").
		AddDiagnostic("operator ~ cannot be applied to undefined").
		Run()
}

// Java: TypeCheckTest.testBitOperation2()
func (c *TypeCheckTest) TestBitOperation2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@return {void}*/function foo(){var a = foo()<<3;}").
		AddDiagnostic("operator << cannot be applied to undefined").
		Run()
}

// Java: TypeCheckTest.testBitOperation3()
func (c *TypeCheckTest) TestBitOperation3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@return {void}*/function foo(){var a = 3<<foo();}").
		AddDiagnostic("operator << cannot be applied to undefined").
		Run()
}

// Java: TypeCheckTest.testBitOperation4()
func (c *TypeCheckTest) TestBitOperation4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@return {void}*/function foo(){var a = foo()>>>3;}").
		AddDiagnostic("operator >>> cannot be applied to undefined").
		Run()
}

// Java: TypeCheckTest.testBitOperation5()
func (c *TypeCheckTest) TestBitOperation5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@return {void}*/function foo(){var a = 3>>>foo();}").
		AddDiagnostic("operator >>> cannot be applied to undefined").
		Run()
}

// Java: TypeCheckTest.testBitOperation6()
func (c *TypeCheckTest) TestBitOperation6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@return {!Object}*/function foo(){var a = foo()&3;}").
		AddDiagnostic(`bad left operand to bitwise operator
found   : Object
required: (boolean|null|number|string|undefined)
`).
		Run()
}

// Java: TypeCheckTest.testBitOperation7()
func (c *TypeCheckTest) TestBitOperation7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_PRIMITIVE_OPERATORS).
		AddSource("var x = null; x |= undefined; x &= 3; x ^= '3'; x |= true;").
		Run()
}

// Java: TypeCheckTest.testBitOperation8()
func (c *TypeCheckTest) TestBitOperation8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_PRIMITIVE_OPERATORS).
		AddSource("var x = void 0; x |= new Number(3);").
		Run()
}

// Java: TypeCheckTest.testBitOperation9()
func (c *TypeCheckTest) TestBitOperation9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_PRIMITIVE_OPERATORS).
		AddSource("var x = void 0; x |= {};").
		AddDiagnostic(`bad right operand to bitwise operator
found   : {}
required: (boolean|null|number|string|undefined)
`).
		Run()
}

// Java: TypeCheckTest.testCall1()
func (c *TypeCheckTest) TestCall1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("3();").
		AddDiagnostic("number expressions are not callable").
		Run()
}

// Java: TypeCheckTest.testCall2()
func (c *TypeCheckTest) TestCall2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @param {!Number} foo*/function bar(foo){ bar('abc'); }").
		AddDiagnostic(`actual parameter 1 of bar does not match formal parameter
found   : string
required: Number
`).
		Run()
}

// Java: TypeCheckTest.testCall3()
func (c *TypeCheckTest) TestCall3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {Function|undefined} */var opt_f;
/** @type {some.unknown.type} */var f1;
var f2 = opt_f || f1;
f2();
`).
		AddDiagnostic("Bad type annotation. Unknown type some.unknown.type").
		Run()
}

// Java: TypeCheckTest.testCall3NullishCoalesce()
func (c *TypeCheckTest) TestCall3NullishCoalesce() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {Function|undefined} */var opt_f;
/** @type {some.unknown.type} */var f1;
var f2 = opt_f ?? f1;
f2();
`).
		AddDiagnostic("Bad type annotation. Unknown type some.unknown.type").
		Run()
}

// Java: TypeCheckTest.testCall4()
func (c *TypeCheckTest) TestCall4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@param {!RegExp} a*/var foo = function bar(a){ bar('abc'); }").
		AddDiagnostic(`actual parameter 1 of bar does not match formal parameter
found   : string
required: RegExp
`).
		Run()
}

// Java: TypeCheckTest.testCall5()
func (c *TypeCheckTest) TestCall5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@param {!RegExp} a*/var foo = function bar(a){ foo('abc'); }").
		AddDiagnostic(`actual parameter 1 of foo does not match formal parameter
found   : string
required: RegExp
`).
		Run()
}

// Java: TypeCheckTest.testCall6()
func (c *TypeCheckTest) TestCall6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {!Number} foo*/function bar(foo){}
bar('abc');
`).
		AddDiagnostic(`actual parameter 1 of bar does not match formal parameter
found   : string
required: Number
`).
		Run()
}

// Java: TypeCheckTest.testCall7()
func (c *TypeCheckTest) TestCall7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {!RegExp} a*/var foo = function bar(a){};
foo('abc');
`).
		AddDiagnostic(`actual parameter 1 of foo does not match formal parameter
found   : string
required: RegExp
`).
		Run()
}

// Java: TypeCheckTest.testCall8()
func (c *TypeCheckTest) TestCall8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {Function|number} */var f;
f();
`).
		AddDiagnostic("(Function|number) expressions are not callable").
		Run()
}

// Java: TypeCheckTest.testCall9()
func (c *TypeCheckTest) TestCall9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = {};
/** @constructor */ goog.Foo = function() {};
/** @param {!goog.Foo} a */ var bar = function(a){};
bar('abc');
`).
		AddDiagnostic(`actual parameter 1 of bar does not match formal parameter
found   : string
required: goog.Foo
`).
		Run()
}

// Java: TypeCheckTest.testCall10()
func (c *TypeCheckTest) TestCall10() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {Function} */var f;f();").
		Run()
}

// Java: TypeCheckTest.testCall11()
func (c *TypeCheckTest) TestCall11() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var f = new Function(); f();").
		Run()
}

// Java: TypeCheckTest.testCall12()
func (c *TypeCheckTest) TestCall12() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {*} x
 * @return {number}
 */
function f(x, y) {
  return x && x.foo();
}
`).
		AddDiagnostic(`inconsistent return type
found   : *
required: number
`).
		AddDiagnostic("Property foo never defined on *" + TypeCheck_POSSIBLE_INEXISTENT_PROPERTY_EXPLANATION).
		Run()
}

// Java: TypeCheckTest.testCall13()
func (c *TypeCheckTest) TestCall13() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var x;
function useX() { var /** string */ str = x(); }
function setX() { x = /** @return {number} */ () => 3; }
`).
		AddDiagnostic(`initializing variable
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testFunctionCall1()
func (c *TypeCheckTest) TestFunctionCall1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {number} x */ var foo = function(x) {};
foo.call(null, 3);
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testFunctionCall2()
func (c *TypeCheckTest) TestFunctionCall2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {number} x */ var foo = function(x) {};
foo.call(null, 'bar');
`).
		AddDiagnostic(`actual parameter 2 of foo.call does not match formal parameter
found   : string
required: number
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testFunctionCall3()
func (c *TypeCheckTest) TestFunctionCall3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {number} x
 * @constructor */
var Foo = function(x) { this.bar.call(null, x); };
/** @type {function(number)} */ Foo.prototype.bar;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testFunctionCall4()
func (c *TypeCheckTest) TestFunctionCall4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {string} x
 * @constructor */
var Foo = function(x) { this.bar.call(null, x); };
/** @type {function(number)} */ Foo.prototype.bar;
`).
		AddDiagnostic(`actual parameter 2 of this.bar.call does not match formal parameter
found   : string
required: number
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testFunctionCall5()
func (c *TypeCheckTest) TestFunctionCall5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {Function} handler
 * @constructor */
var Foo = function(handler) { handler.call(this, x); };
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testFunctionCall6()
func (c *TypeCheckTest) TestFunctionCall6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {Function} handler
 * @constructor */
var Foo = function(handler) { handler.apply(this, x); };
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testFunctionCall7()
func (c *TypeCheckTest) TestFunctionCall7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {Function} handler
 * @param {Object} opt_context */
var Foo = function(handler, opt_context) {
  handler.call(opt_context, x);
};
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testFunctionCall8()
func (c *TypeCheckTest) TestFunctionCall8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {Function} handler
 * @param {Object} opt_context */
var Foo = function(handler, opt_context) {
  handler.apply(opt_context, x);
};
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testFunctionCall9()
func (c *TypeCheckTest) TestFunctionCall9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor
 * @template T
 **/ function Foo() {}
/** @param {T} x */ Foo.prototype.bar = function(x) {}
var foo = /** @type {Foo<string>} */ (new Foo());
foo.bar(3);
`).
		AddDiagnostic(`actual parameter 1 of Foo.prototype.bar does not match formal parameter
found   : number
required: string
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testFunctionBind1()
func (c *TypeCheckTest) TestFunctionBind1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {function(string, number): boolean} */
function f(x, y) { return true; }
f.bind(null, 3);
`).
		AddDiagnostic(`actual parameter 2 of f.bind does not match formal parameter
found   : number
required: string
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testFunctionBind2()
func (c *TypeCheckTest) TestFunctionBind2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {function(number): boolean} */
function f(x) { return true; }
f(f.bind(null, 3)());
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : boolean
required: number
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testFunctionBind3()
func (c *TypeCheckTest) TestFunctionBind3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {function(number, string): boolean} */
function f(x, y) { return true; }
f.bind(null, 3)(true);
`).
		AddDiagnostic(`actual parameter 1 of function does not match formal parameter
found   : boolean
required: string
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testFunctionBind4()
func (c *TypeCheckTest) TestFunctionBind4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {...number} x */
function f(x) {}
f.bind(null, 3, 3, 3)(true);
`).
		AddDiagnostic(`actual parameter 1 of function does not match formal parameter
found   : boolean
required: number
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testFunctionBind5()
func (c *TypeCheckTest) TestFunctionBind5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {...number} x */
function f(x) {}
f.bind(null, true)(3, 3, 3);
`).
		AddDiagnostic(`actual parameter 2 of f.bind does not match formal parameter
found   : boolean
required: number
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testFunctionBind6()
func (c *TypeCheckTest) TestFunctionBind6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */
function MyType() {
  /** @type {number} */
  this.x = 0;
  var f = function() {
    this.x = 'str';
  }.bind(this);
}
`).
		AddDiagnostic(`assignment to property x of MyType
found   : string
required: number
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testFunctionBind7()
func (c *TypeCheckTest) TestFunctionBind7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */
function MyType() {
  /** @type {number} */
  this.x = 0;
}
var m = new MyType;
(function f() {this.x = 'str';}).bind(m);
`).
		AddDiagnostic(`assignment to property x of MyType
found   : string
required: number
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testFunctionBind8()
func (c *TypeCheckTest) TestFunctionBind8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */
function MyType() {}

/** @constructor */
function AnotherType() {}
AnotherType.prototype.foo = function() {};

/** @type {?} */
var m = new MyType;
(function f() {this.foo();}).bind(m);
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testFunctionBind9()
func (c *TypeCheckTest) TestFunctionBind9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */
function MyType() {}

/** @constructor */
function AnotherType() {}
AnotherType.prototype.foo = function() {};

var m = new MyType;
(function f() {this.foo();}).bind(m);
`).
		AddDiagnosticType(TypeCheck_INEXISTENT_PROPERTY).
		Run()
}

// Java: TypeCheckTest.testGoogBind1()
func (c *TypeCheckTest) TestGoogBind1() {
	c.TestClosureTypes(`goog.bind = function(var_args) {};
/** @type {function(number): boolean} */
function f(x, y) { return true; }
f(goog.bind(f, null, 'x')());
`, jutil.StringPtr(`actual parameter 1 of f does not match formal parameter
found   : boolean
required: number
`))
}

// Java: TypeCheckTest.testGoogBind2()
func (c *TypeCheckTest) TestGoogBind2() {
	c.TestClosureTypes(`goog.bind = function(var_args) {};
/** @type {function(boolean): boolean} */
function f(x, y) { return true; }
f(goog.bind(f, null, 'x')());
`, nil)
}

// Java: TypeCheckTest.testCast2()
func (c *TypeCheckTest) TestCast2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */function base() {}
/** @constructor
 @extends {base} */function derived() {}
/** @type {base} */ var baz = new derived();
`).
		Run()
}

// Java: TypeCheckTest.testCast3()
func (c *TypeCheckTest) TestCast3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */function base() {}
/** @constructor @extends {base} */function derived() {}
/** @type {!derived} */ var baz = new base();
`).
		AddDiagnostic(`initializing variable
found   : base
required: derived
`).
		Run()
}

// Java: TypeCheckTest.testCast3a()
func (c *TypeCheckTest) TestCast3a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */function Base() {}
/** @constructor @extends {Base} */function Derived() {}
var baseInstance = new Base();
/** @type {!Derived} */ var baz = baseInstance;
`).
		AddDiagnostic(`initializing variable
found   : Base
required: Derived
`).
		Run()
}

// Java: TypeCheckTest.testCast4()
func (c *TypeCheckTest) TestCast4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */function base() {}
/** @constructor
 * @extends {base} */function derived() {}
/** @type {!derived} */ var baz =
/** @type {!derived} */(new base());
`).
		Run()
}

// Java: TypeCheckTest.testCast4Types()
func (c *TypeCheckTest) TestCast4Types() {
	var root *rhino.Node = c.ParseAndTypeCheck(`/** @constructor */function base() {}
/** @constructor
 * @extends {base} */function derived() {}
/** @type {!derived} */ var baz =
/** @type {!derived} */(new base());
`)
	var castedExprNode *rhino.Node = root.GetLastChild().GetFirstFirstChild().
		GetFirstChild()
	AssertThat(castedExprNode.GetJSType().String()).IsEqualTo("derived")
	AssertThat(castedExprNode.GetJSTypeBeforeCast().String()).IsEqualTo("base")
}

// Java: TypeCheckTest.testCast5()
func (c *TypeCheckTest) TestCast5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */function foo() {}
/** @constructor */function bar() {}
var baz = /** @type {!foo} */(new bar);
`).
		AddDiagnostic(`invalid cast - must be a subtype or supertype
from: bar
to  : foo
`).
		Run()
}

// Java: TypeCheckTest.testCast5a()
func (c *TypeCheckTest) TestCast5a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */function foo() {}
/** @constructor */function bar() {}
var barInstance = new bar;
var baz = /** @type {!foo} */(barInstance);
`).
		AddDiagnostic(`invalid cast - must be a subtype or supertype
from: bar
to  : foo
`).
		Run()
}

// Java: TypeCheckTest.testCast6()
func (c *TypeCheckTest) TestCast6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */function foo() {}
/** @constructor 
 @extends foo */function bar() {}
var baz = /** @type {!bar} */(new bar);
var baz = /** @type {!foo} */(new foo);
var baz = /** @type {bar} */(new bar);
var baz = /** @type {foo} */(new foo);
var baz = /** @type {!foo} */(new bar);
var baz = /** @type {!bar} */(new foo);
var baz = /** @type {foo} */(new bar);
var baz = /** @type {bar} */(new foo);
`).
		Run()
}

// Java: TypeCheckTest.testCast7()
func (c *TypeCheckTest) TestCast7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var x = /** @type {foo} */ (new Object());").
		AddDiagnostic("Bad type annotation. Unknown type foo").
		Run()
}

// Java: TypeCheckTest.testCast8()
func (c *TypeCheckTest) TestCast8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function f() { return /** @type {foo} */ (new Object()); }").
		AddDiagnostic("Bad type annotation. Unknown type foo").
		Run()
}

// Java: TypeCheckTest.testCast9()
func (c *TypeCheckTest) TestCast9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var foo = {};
function f() { return /** @type {foo} */ (new Object()); }
`).
		AddDiagnostic(`Bad type annotation. Unknown type foo
It's possible that 'foo' refers to a value, not a type.
`).
		Run()
}

// Java: TypeCheckTest.testCast10()
func (c *TypeCheckTest) TestCast10() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var foo = function() {};
function f() { return /** @type {foo} */ (new Object()); }
`).
		AddDiagnostic(`Bad type annotation. Unknown type foo
It's possible that 'foo' refers to a value, not a type.
`).
		Run()
}

// Java: TypeCheckTest.testCast11()
func (c *TypeCheckTest) TestCast11() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = {}; goog.foo = {};
function f() { return /** @type {goog.foo} */ (new Object()); }
`).
		AddDiagnostic(`Bad type annotation. Unknown type goog.foo
It's possible that 'goog.foo' refers to a value, not a type.
`).
		Run()
}

// Java: TypeCheckTest.testCast12()
func (c *TypeCheckTest) TestCast12() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = {}; goog.foo = function() {};
function f() { return /** @type {goog.foo} */ (new Object()); }
`).
		AddDiagnostic(`Bad type annotation. Unknown type goog.foo
It's possible that 'goog.foo' refers to a value, not a type.
`).
		Run()
}

// Java: TypeCheckTest.testCast13()
func (c *TypeCheckTest) TestCast13() {
	c.TestClosureTypes(`goog.forwardDeclare('goog.foo');
goog.foo = function() {};
function f() { return /** @type {goog.foo} */ (new Object()); }
`, nil)
}

// Java: TypeCheckTest.testCast14()
func (c *TypeCheckTest) TestCast14() {
	c.TestClosureTypes(`goog.forwardDeclare('goog.bar');
function f() { return /** @type {goog.bar} */ (new Object()); }
`, nil)
}

// Java: TypeCheckTest.testCast15()
func (c *TypeCheckTest) TestCast15() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`for (var i = 0; i < 10; i++) {
var x = /** @type {Object|number} */ ({foo: 3});
/** @param {number} x */ function f(x) {}
f(x.foo);
f([].foo);
}
`).
		AddDiagnostic("Property foo never defined on Array<?>").
		Run()
}

// Java: TypeCheckTest.testCast15b()
func (c *TypeCheckTest) TestCast15b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`for (var i = 0; i < 10; i++) {
var x = /** @type {{foo:number}}|number} */ ({foo: 3});
/** @param {number} x */ function f(x) {}
f(x.foo);
f([].foo);
}
`).
		AddDiagnostic("Property foo never defined on Array<?>").
		Run()
}

// Java: TypeCheckTest.testCast16()
func (c *TypeCheckTest) TestCast16() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`for (var i = 0; i < 10; i++) {
var x = /** @type {Object|number} */ (
  {/** @type {string} */ foo: 3});
}
`).
		AddDiagnostic(`assignment to property foo of {foo: string}
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testCast17a()
func (c *TypeCheckTest) TestCast17a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
/** @type {Foo} */ var x = /** @type {Foo} */ (y)
`).
		Run()
}

// Java: TypeCheckTest.testCast17b()
func (c *TypeCheckTest) TestCast17b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
/** @type {Foo} */ var x = /** @type {Foo} */ ({})
`).
		Run()
}

// Java: TypeCheckTest.testCast19()
func (c *TypeCheckTest) TestCast19() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var x = 'string';
/** @type {number} */
var y = /** @type {number} */(x);
`).
		AddDiagnostic(`invalid cast - must be a subtype or supertype
from: string
to  : number
`).
		Run()
}

// Java: TypeCheckTest.testCast20()
func (c *TypeCheckTest) TestCast20() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum {boolean|null} */
var X = {
  AA: true,
  BB: false,
  CC: null
};
var y = /** @type {X} */(true);
`).
		Run()
}

// Java: TypeCheckTest.testCast21()
func (c *TypeCheckTest) TestCast21() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum {boolean|null} */
var X = {
  AA: true,
  BB: false,
  CC: null
};
var value = true;
var y = /** @type {X} */(value);
`).
		Run()
}

// Java: TypeCheckTest.testCast22()
func (c *TypeCheckTest) TestCast22() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var x = null;
var y = /** @type {number} */(x);
`).
		AddDiagnostic(`invalid cast - must be a subtype or supertype
from: null
to  : number
`).
		Run()
}

// Java: TypeCheckTest.testCast23()
func (c *TypeCheckTest) TestCast23() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var x = null;
var y = /** @type {Number} */(x);
`).
		Run()
}

// Java: TypeCheckTest.testCast24()
func (c *TypeCheckTest) TestCast24() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var x = undefined;
var y = /** @type {number} */(x);
`).
		AddDiagnostic(`invalid cast - must be a subtype or supertype
from: undefined
to  : number
`).
		Run()
}

// Java: TypeCheckTest.testCast25()
func (c *TypeCheckTest) TestCast25() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var x = undefined;
var y = /** @type {number|undefined} */(x);
`).
		Run()
}

// Java: TypeCheckTest.testCast26()
func (c *TypeCheckTest) TestCast26() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function fn(dir) {
  var node = dir ? 1 : 2;
  fn(/** @type {number} */ (node));
}
`).
		Run()
}

// Java: TypeCheckTest.testCast27()
func (c *TypeCheckTest) TestCast27() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ function I() {}
/** @constructor */ function C() {}
var x = new C();
var y = /** @type {I} */(x);
`).
		Run()
}

// Java: TypeCheckTest.testCast27a()
func (c *TypeCheckTest) TestCast27a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ function I() {}
/** @constructor */ function C() {}
/** @type {C} */ var x ;
var y = /** @type {I} */(x);
`).
		Run()
}

// Java: TypeCheckTest.testCast28()
func (c *TypeCheckTest) TestCast28() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ function I() {}
/** @constructor */ function C() {}
/** @type {!I} */ var x;
var y = /** @type {C} */(x);
`).
		Run()
}

// Java: TypeCheckTest.testCast28a()
func (c *TypeCheckTest) TestCast28a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ function I() {}
/** @constructor */ function C() {}
/** @type {I} */ var x;
var y = /** @type {C} */(x);
`).
		Run()
}

// Java: TypeCheckTest.testCast29a()
func (c *TypeCheckTest) TestCast29a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function C() {}
var x = new C();
var y = /** @type {{remoteJids: Array, sessionId: string}} */(x);
`).
		Run()
}

// Java: TypeCheckTest.testCast29b()
func (c *TypeCheckTest) TestCast29b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function C() {}
/** @type {C} */ var x;
var y = /** @type {{prop1: Array, prop2: string}} */(x);
`).
		Run()
}

// Java: TypeCheckTest.testCast29c()
func (c *TypeCheckTest) TestCast29c() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function C() {}
/** @type {{remoteJids: Array, sessionId: string}} */ var x ;
var y = /** @type {C} */(x);
`).
		Run()
}

// Java: TypeCheckTest.testCast30()
func (c *TypeCheckTest) TestCast30() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function C() {}
/** @type {function():string} */ var x ;
var y = /** @type {function():?} */(x);
`).
		Run()
}

// Java: TypeCheckTest.testCast31()
func (c *TypeCheckTest) TestCast31() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function C() {}
/** @type {function(*)} */ var x ;
var y = /** @type {function(string)} */(x);
`).
		Run()
}

// Java: TypeCheckTest.testCast32()
func (c *TypeCheckTest) TestCast32() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function C() {}
/** @type {Object} */ var x ;
var y = /** @type {null|{length:number}} */(x);
`).
		Run()
}

// Java: TypeCheckTest.testCast33a()
func (c *TypeCheckTest) TestCast33a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function C() {}
/** @type {null|undefined} */ var x ;
var y = /** @type {string?|undefined} */(x);
`).
		Run()
}

// Java: TypeCheckTest.testCast33b()
func (c *TypeCheckTest) TestCast33b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function C() {}
/** @type {null|undefined} */ var x ;
var y = /** @type {string|undefined} */(x);
`).
		Run()
}

// Java: TypeCheckTest.testCast33c()
func (c *TypeCheckTest) TestCast33c() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function C() {}
/** @type {null|undefined} */ var x ;
var y = /** @type {string?} */(x);
`).
		Run()
}

// Java: TypeCheckTest.testCast33d()
func (c *TypeCheckTest) TestCast33d() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function C() {}
/** @type {null|undefined} */ var x ;
var y = /** @type {null} */(x);
`).
		Run()
}

// Java: TypeCheckTest.testCast34a()
func (c *TypeCheckTest) TestCast34a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function C() {}
/** @type {Object} */ var x ;
var y = /** @type {Function} */(x);
`).
		Run()
}

// Java: TypeCheckTest.testCast34b()
func (c *TypeCheckTest) TestCast34b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function C() {}
/** @type {Function} */ var x ;
var y = /** @type {Object} */(x);
`).
		Run()
}

// Java: TypeCheckTest.testCastToNameRequiringPropertyResolution()
func (c *TypeCheckTest) TestCastToNameRequiringPropertyResolution() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns("var unknownVar;").
		AddSource("const foo = {bar: {}};\nconst bar = foo.bar;\nbar.Class = class {\n  /** @return {number} */\n  id() { return 0; }\n};\n// Because `foo.bar.Class = ...` was never directly assigned, the type 'foo.bar.Class'\n// is not in the JSTypeRegistry. It's resolved through NamedType#resolveViaProperty.\n// The same thing would have occurred if we assigned 'foo.bar.Class = ...' then\n// referred to '!bar.Class' in the JSDoc.\n// Verify that type inference correctly infers the id property's type.\nconst /** null */ n = /** @type {!foo.bar.Class} */ (unknownVar).id;\n").
		AddDiagnostic(`initializing variable
found   : function(this:bar.Class): number
required: null
`).
		Run()
}

// Java: TypeCheckTest.testNestedCasts()
func (c *TypeCheckTest) TestNestedCasts() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */var T = function() {};
/** @constructor */var V = function() {};
/**
* @param {boolean} b
* @return {T|V}
*/
function f(b) { return b ? new T() : new V(); }
/**
* @param {boolean} b
* @return {boolean|undefined}
*/
function g(b) { return b ? true : undefined; }
/** @return {T} */
function h() {
return /** @type {T} */ (f(/** @type {boolean} */ (g(true))));
}
`).
		Run()
}

// Java: TypeCheckTest.testNativeCast1()
func (c *TypeCheckTest) TestNativeCast1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {number} x */ function f(x) {}
f(String(true));
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testNativeCast2()
func (c *TypeCheckTest) TestNativeCast2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {string} x */ function f(x) {}
f(Number(true));
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testNativeCast3()
func (c *TypeCheckTest) TestNativeCast3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {number} x */ function f(x) {}
f(Boolean(''));
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckTest.testNativeCast4()
func (c *TypeCheckTest) TestNativeCast4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {number} x */ function f(x) {}
f(Array(1));
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : Array
required: number
`).
		Run()
}

// Java: TypeCheckTest.testBadConstructorCall()
func (c *TypeCheckTest) TestBadConstructorCall() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
Foo();
`).
		AddDiagnostic("Constructor (typeof Foo) should be called with the \"new\" keyword").
		Run()
}

// Java: TypeCheckTest.testTypeof()
func (c *TypeCheckTest) TestTypeof() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/**@return {void}*/function foo(){ var a = typeof foo(); }").
		Run()
}

// Java: TypeCheckTest.testTypeof2()
func (c *TypeCheckTest) TestTypeof2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function f(){ if (typeof 123 == 'numbr') return 321; }").
		AddDiagnostic("unknown type: numbr").
		Run()
}

// Java: TypeCheckTest.testTypeof3()
func (c *TypeCheckTest) TestTypeof3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f() {
  return (
      typeof 123 == 'number' ||
      typeof 123 == 'string' ||
      typeof 123 == 'boolean' ||
      typeof 123 == 'undefined' ||
      typeof 123 == 'function' ||
      typeof 123 == 'object' ||
      typeof 123 == 'symbol' ||
      typeof 123 == 'unknown'); }
`).
		Run()
}

// Java: TypeCheckTest.testConstDecl1()
func (c *TypeCheckTest) TestConstDecl1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {?number} x
* @return {boolean} */
function f(x) {
  if (x) { /** @const */ var y = x; return y } return true;
}
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: boolean
`).
		Run()
}

// Java: TypeCheckTest.testConstDecl2()
func (c *TypeCheckTest) TestConstDecl2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {?number} x */
function f(x) {
  if (x) {
    /** @const */ var y = x;
    /** @return {boolean} */ function g() { return y; }
  }
}
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: boolean
`).
		Run()
}

// Java: TypeCheckTest.testConstructorType1()
func (c *TypeCheckTest) TestConstructorType1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**@constructor*/function Foo(){}
/**@type{!Foo}*/var f = new Date();
`).
		AddDiagnostic(`initializing variable
found   : Date
required: Foo
`).
		Run()
}

// Java: TypeCheckTest.testConstructorType2()
func (c *TypeCheckTest) TestConstructorType2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**@constructor*/function Foo(){
/**@type{Number}*/this.bar = new Number(5);
}
/**@type{Foo}*/var f = new Foo();
/**@type{Number}*/var n = f.bar;
`).
		Run()
}

// Java: TypeCheckTest.testConstructorType3()
func (c *TypeCheckTest) TestConstructorType3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**@type{Foo}*/var f = new Foo();
/**@type{Number}*/var n = f.bar;
/**@constructor*/function Foo(){
/**@type{Number}*/this.bar = new Number(5);
}
`).
		Run()
}

// Java: TypeCheckTest.testConstructorType4()
func (c *TypeCheckTest) TestConstructorType4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**@constructor*/function Foo(){
/**@type{!Number}*/this.bar = new Number(5);
}
/**@type{!Foo}*/var f = new Foo();
/**@type{!String}*/var n = f.bar;
`).
		AddDiagnostic(`initializing variable
found   : Number
required: String
`).
		Run()
}

// Java: TypeCheckTest.testConstructorType5()
func (c *TypeCheckTest) TestConstructorType5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**@constructor*/function Foo(){}
if (Foo){}
`).
		Run()
}

// Java: TypeCheckTest.testConstructorType6()
func (c *TypeCheckTest) TestConstructorType6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */
function bar() {}
function _foo() {
 /** @param {bar} x */
  function f(x) {}
}
`).
		Run()
}

// Java: TypeCheckTest.testConstructorType7()
func (c *TypeCheckTest) TestConstructorType7() {
	var p *TypeCheckTestCase_TypeCheckResult = c.ParseAndTypeCheckWithScope("/** @constructor */function A(){};")
	var type_ rhino.JSType = p.scope.GetVar("A").GetType()
	IsInstanceOf[rhino.FunctionType](AssertThat(type_))
	var fType rhino.FunctionType = type_.(rhino.FunctionType)
	AssertThat(fType.GetReferenceName()).IsEqualTo("A")
}

// Java: TypeCheckTest.testConstructorType8()
func (c *TypeCheckTest) TestConstructorType8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var ns = {};
ns.create = function() { return function() {}; };
/** @constructor */ ns.Foo = ns.create();
ns.Foo.prototype = {x: 0, y: 0};
/**
 * @param {ns.Foo} foo
 * @return {string}
 */
function f(foo) {
  return foo.x;
}
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testConstructorType9()
func (c *TypeCheckTest) TestConstructorType9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var ns = {};
ns.create = function() { return function() {}; };
ns.extend = function(x) { return x; };
/** @constructor */ ns.Foo = ns.create();
ns.Foo.prototype = ns.extend({x: 0, y: 0});
/**
 * @param {ns.Foo} foo
 * @return {string}
 */
function f(foo) {
  return foo.x;
}
`).
		Run()
}

// Java: TypeCheckTest.testConstructorType10()
func (c *TypeCheckTest) TestConstructorType10() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */
function NonStr() {}
/**
 * @constructor
 * @struct
 * @extends{NonStr}
 */
function NonStrKid() {}
`).
		Run()
}

// Java: TypeCheckTest.testConstructorType11()
func (c *TypeCheckTest) TestConstructorType11() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */
function NonDict() {}
/**
 * @constructor
 * @dict
 * @extends{NonDict}
 */
function NonDictKid() {}
`).
		Run()
}

// Java: TypeCheckTest.testConstructorType12()
func (c *TypeCheckTest) TestConstructorType12() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @struct
 */
function Bar() {}
Bar.prototype = {};
`).
		Run()
}

// Java: TypeCheckTest.testBadStruct()
func (c *TypeCheckTest) TestBadStruct() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @struct */function Struct1() {}").
		AddDiagnostic("@struct used without @constructor for Struct1").
		Run()
}

// Java: TypeCheckTest.testBadDict()
func (c *TypeCheckTest) TestBadDict() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @dict */function Dict1() {}").
		AddDiagnostic("@dict used without @constructor for Dict1").
		Run()
}

// Java: TypeCheckTest.testAnonymousPrototype1()
func (c *TypeCheckTest) TestAnonymousPrototype1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var ns = {};
/** @constructor */ ns.Foo = function() {
  this.bar(3, 5);
};
ns.Foo.prototype = {
  bar: function(x) {}
};
`).
		AddDiagnostic(`Function ns.Foo.prototype.bar: called with 2 argument(s). Function requires at least 1 argument(s) and no more than 1 argument(s).
`).
		Run()
}

// Java: TypeCheckTest.testAnonymousPrototype2()
func (c *TypeCheckTest) TestAnonymousPrototype2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ var Foo = function() {};
Foo.prototype = {
  foo: function(x) {}
};
/**
 * @constructor
 * @implements {Foo}
 */ var Bar = function() {};
`).
		AddDiagnostic("property foo on interface Foo is not implemented by type Bar").
		Run()
}

// Java: TypeCheckTest.testAnonymousType1()
func (c *TypeCheckTest) TestAnonymousType1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f() { return {}; }
/** @constructor */
f().bar = function() {};
`).
		Run()
}

// Java: TypeCheckTest.testAnonymousType2()
func (c *TypeCheckTest) TestAnonymousType2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f() { return {}; }
/** @interface */
f().bar = function() {};
`).
		Run()
}

// Java: TypeCheckTest.testAnonymousType3()
func (c *TypeCheckTest) TestAnonymousType3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f() { return {}; }
/** @enum */
f().bar = {FOO: 1};
`).
		Run()
}

// Java: TypeCheckTest.testBang1()
func (c *TypeCheckTest) TestBang1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {Object} x
* @return {!Object} */
function f(x) { return x; }
`).
		AddDiagnostic(`inconsistent return type
found   : (Object|null)
required: Object
`).
		Run()
}

// Java: TypeCheckTest.testBang2()
func (c *TypeCheckTest) TestBang2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {Object} x
* @return {!Object} */
function f(x) { return x ? x : new Object(); }
`).
		Run()
}

// Java: TypeCheckTest.testBang3()
func (c *TypeCheckTest) TestBang3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {Object} x
* @return {!Object} */
function f(x) { return /** @type {!Object} */ (x); }
`).
		Run()
}

// Java: TypeCheckTest.testBang4()
func (c *TypeCheckTest) TestBang4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**@param {Object} x
@param {Object} y
@return {boolean}*/
function f(x, y) {
if (typeof x != 'undefined') { return x == y; }
else { return x != y; }
}
`).
		Run()
}

// Java: TypeCheckTest.testBang5()
func (c *TypeCheckTest) TestBang5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**@param {Object} x
@param {Object} y
@return {boolean}*/
function f(x, y) { return !!x && x == y; }
`).
		Run()
}

// Java: TypeCheckTest.testBang6()
func (c *TypeCheckTest) TestBang6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {Object?} x
* @return {Object} */
function f(x) { return x; }
`).
		Run()
}

// Java: TypeCheckTest.testBang7()
func (c *TypeCheckTest) TestBang7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**@param {(Object|string|null)} x
* @return {(Object|string)}*/function f(x) { return x; }
`).
		Run()
}

// Java: TypeCheckTest.testDefinePropertyOnNullableObject1()
func (c *TypeCheckTest) TestDefinePropertyOnNullableObject1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @type {Object} */ var n = {};
/** @type {number} */ n.x = 1;
/** @return {boolean} */ function f() { return n.x; }
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: boolean
`).
		Run()
}

// Java: TypeCheckTest.testDefinePropertyOnNullableObject1a()
func (c *TypeCheckTest) TestDefinePropertyOnNullableObject1a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const */ var n = {};
/** @type {number} */ n.x = 1;
/** @return {boolean} */function f() { return n.x; }
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: boolean
`).
		Run()
}

// Java: TypeCheckTest.testDefinePropertyOnObject()
func (c *TypeCheckTest) TestDefinePropertyOnObject() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @type {!Object} */ var n = {};
/** @type {number} */ n.x = 1;
/** @return {boolean} */function f() { return n.x; }
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: boolean
`).
		Run()
}

// Java: TypeCheckTest.testDefinePropertyOnNullableObject2()
func (c *TypeCheckTest) TestDefinePropertyOnNullableObject2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @constructor */ var T = function() {};
/** @param {T} t
@return {boolean} */function f(t) {
t.x = 1; return t.x; }
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: boolean
`).
		Run()
}

// Java: TypeCheckTest.testDefinePropertyOnNullableObject2b()
func (c *TypeCheckTest) TestDefinePropertyOnNullableObject2b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ var T = function() {};
/** @param {T} t */function f(t) { t.x = 1; }
`).
		AddDiagnostic("Property x never defined on T").
		Run()
}

// Java: TypeCheckTest.testUnknownConstructorInstanceType1()
func (c *TypeCheckTest) TestUnknownConstructorInstanceType1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @return {Array} */ function g(f) { return new f(); }").
		Run()
}

// Java: TypeCheckTest.testUnknownConstructorInstanceType2()
func (c *TypeCheckTest) TestUnknownConstructorInstanceType2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function g(f) { return /** @type {Array} */(new f()); }").
		Run()
}

// Java: TypeCheckTest.testUnknownConstructorInstanceType3()
func (c *TypeCheckTest) TestUnknownConstructorInstanceType3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function g(f) { var x = new f(); x.a = 1; return x; }").
		Run()
}

// Java: TypeCheckTest.testUnknownPrototypeChain1()
func (c *TypeCheckTest) TestUnknownPrototypeChain1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
* @param {Object} co
 * @return {Object}
 */
function inst(co) {
 /** @constructor */
 var c = function() {};
 c.prototype = co.prototype;
 return new c;
}
`).
		AddDiagnostic("Property prototype never defined on Object").
		Run()
}

// Java: TypeCheckTest.testUnknownPrototypeChain2()
func (c *TypeCheckTest) TestUnknownPrototypeChain2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {Function} co
 * @return {Object}
 */
function inst(co) {
 /** @constructor */
 var c = function() {};
 c.prototype = co.prototype;
 return new c;
}
`).
		Run()
}

// Java: TypeCheckTest.testNamespacedConstructor()
func (c *TypeCheckTest) TestNamespacedConstructor() {
	var root *rhino.Node = c.ParseAndTypeCheck(`var goog = {};
/** @constructor */ goog.MyClass = function() {};
/** @return {!goog.MyClass} */
function foo() { return new goog.MyClass(); }
`)
	var typeOfFoo rhino.JSType = root.GetLastChild().GetJSType()
	rhinotesting.IsInstanceOf[rhino.FunctionType](TypeSubject_AssertType(typeOfFoo).Subject)
	var retType rhino.JSType = (typeOfFoo.(rhino.FunctionType)).GetReturnType()
	rhinotesting.IsInstanceOf[rhino.ObjectType](TypeSubject_AssertType(retType).Subject)
	AssertThat((retType.(rhino.ObjectType)).GetReferenceName()).IsEqualTo("goog.MyClass")
}

// Java: TypeCheckTest.testComplexNamespace()
func (c *TypeCheckTest) TestComplexNamespace() {
	js := `var goog = {};
goog.foo = {};
goog.foo.bar = 5;
`
	var p *TypeCheckTestCase_TypeCheckResult = c.ParseAndTypeCheckWithScope(js)
	var googScopeType rhino.JSType = p.scope.GetVar("goog").GetType()
	IsInstanceOf[rhino.ObjectType](AssertThat(googScopeType))
	AssertWithMessage("foo property not present on goog type").That(googScopeType.HasProperty("foo")).
		IsTrue()
	AssertWithMessage("bar property present on goog type").That(googScopeType.HasProperty("bar")).
		IsFalse()
	var varNode *rhino.Node = p.root.GetFirstChild()
	NodeSubject_AssertNode(varNode).HasToken(rhino.Token_VAR)
	var googNodeType rhino.JSType = varNode.GetFirstChild().GetJSType()
	IsInstanceOf[rhino.ObjectType](AssertThat(googNodeType))
	AssertThat(googScopeType).IsSameInstanceAs(googNodeType)
	var getpropFoo1 *rhino.Node = varNode.GetNext().GetFirstFirstChild()
	NodeSubject_AssertNode(getpropFoo1).HasToken(rhino.Token_GETPROP)
	AssertThat(getpropFoo1.GetFirstChild().GetString()).IsEqualTo("goog")
	var googGetpropFoo1Type rhino.JSType = getpropFoo1.GetFirstChild().GetJSType()
	IsInstanceOf[rhino.ObjectType](AssertThat(googGetpropFoo1Type))
	AssertThat(googGetpropFoo1Type).IsSameInstanceAs(googScopeType)
	var googFooType rhino.JSType = (googScopeType.(rhino.ObjectType)).GetPropertyType("foo")
	IsInstanceOf[rhino.ObjectType](AssertThat(googFooType))
	var getpropFoo2 *rhino.Node = varNode.GetNext().GetNext().
		GetFirstFirstChild().
		GetFirstChild()
	NodeSubject_AssertNode(getpropFoo2).HasToken(rhino.Token_GETPROP)
	AssertThat(getpropFoo2.GetFirstChild().GetString()).IsEqualTo("goog")
	var googGetpropFoo2Type rhino.JSType = getpropFoo2.GetFirstChild().GetJSType()
	IsInstanceOf[rhino.ObjectType](AssertThat(googGetpropFoo2Type))
	AssertThat(googGetpropFoo2Type).IsSameInstanceAs(googScopeType)
	var googFooGetprop2Type rhino.JSType = getpropFoo2.GetJSType()
	IsInstanceOf[rhino.ObjectType](AssertWithMessage("goog.foo incorrectly annotated in goog.foo.bar selection").That(googFooGetprop2Type))
	var googFooGetprop2ObjectType rhino.ObjectType = googFooGetprop2Type.(rhino.ObjectType)
	AssertWithMessage("foo property present on goog.foo type").That(googFooGetprop2ObjectType.HasProperty("foo")).
		IsFalse()
	AssertWithMessage("bar property not present on goog.foo type").That(googFooGetprop2ObjectType.HasProperty("bar")).
		IsTrue()
	c.AssertTypeEqualsMsg("bar property on goog.foo type incorrectly inferred", c.GetNativeNumberType(), googFooGetprop2ObjectType.GetPropertyType("bar"))
}

// Java: TypeCheckTest.testAddingMethodsUsingPrototypeIdiomSimpleNamespace()
func (c *TypeCheckTest) TestAddingMethodsUsingPrototypeIdiomSimpleNamespace() {
	var js1Node *rhino.Node = c.ParseAndTypeCheckWithExterns(CompilerTypeTestCase_DEFAULT_EXTERNS, `/** @constructor */function A() {}
A.prototype.m1 = 5
`)
	var instanceType rhino.ObjectType = TypeCheckTestCase_GetInstanceType(js1Node)
	c.AssertHasXMorePropertiesThanNativeObject(instanceType, 1)
	c.CheckObjectType(instanceType, "m1", c.GetNativeNumberType())
}

// Java: TypeCheckTest.testAddingMethodsUsingPrototypeIdiomComplexNamespace1()
func (c *TypeCheckTest) TestAddingMethodsUsingPrototypeIdiomComplexNamespace1() {
	var p *TypeCheckTestCase_TypeCheckResult = c.ParseAndTypeCheckWithScopeWithExterns(CompilerTypeTestCase_DEFAULT_EXTERNS, `var goog = {};
goog.A = /** @constructor */function() {};
/** @type {number} */goog.A.prototype.m1 = 5
`)
	c.TestAddingMethodsUsingPrototypeIdiomComplexNamespace(p)
}

// Java: TypeCheckTest.testAddingMethodsUsingPrototypeIdiomComplexNamespace2()
func (c *TypeCheckTest) TestAddingMethodsUsingPrototypeIdiomComplexNamespace2() {
	var p *TypeCheckTestCase_TypeCheckResult = c.ParseAndTypeCheckWithScopeWithExterns(CompilerTypeTestCase_DEFAULT_EXTERNS, `var goog = {};
/** @constructor */goog.A = function() {};
/** @type {number} */goog.A.prototype.m1 = 5
`)
	c.TestAddingMethodsUsingPrototypeIdiomComplexNamespace(p)
}

// Java: TypeCheckTest.testAddingMethodsUsingPrototypeIdiomComplexNamespace(TypeCheckResult)
func (c *TypeCheckTest) TestAddingMethodsUsingPrototypeIdiomComplexNamespace(p *TypeCheckTestCase_TypeCheckResult) {
	var goog rhino.ObjectType = p.scope.GetVar("goog").GetType().(rhino.ObjectType)
	c.AssertHasXMorePropertiesThanNativeObject(goog, 1)
	var googA rhino.JSType = goog.GetPropertyType("A")
	AssertThat(googA).IsNotNull()
	IsInstanceOf[rhino.FunctionType](AssertThat(googA))
	var googAFunction rhino.FunctionType = googA.(rhino.FunctionType)
	var classA rhino.ObjectType = googAFunction.GetInstanceType()
	c.AssertHasXMorePropertiesThanNativeObject(classA, 1)
	c.CheckObjectType(classA, "m1", c.GetNativeNumberType())
}

// Java: TypeCheckTest.testAddingMethodsPrototypeIdiomAndObjectLiteralSimpleNamespace()
func (c *TypeCheckTest) TestAddingMethodsPrototypeIdiomAndObjectLiteralSimpleNamespace() {
	var js1Node *rhino.Node = c.ParseAndTypeCheckWithExterns(CompilerTypeTestCase_DEFAULT_EXTERNS, `/** @constructor */function A() {}
A.prototype = {m1: 5, m2: true}
`)
	var instanceType rhino.ObjectType = TypeCheckTestCase_GetInstanceType(js1Node)
	c.AssertHasXMorePropertiesThanNativeObject(instanceType, 2)
	c.CheckObjectType(instanceType, "m1", c.GetNativeNumberType())
	c.CheckObjectType(instanceType, "m2", c.GetNativeBooleanType())
}

// Java: TypeCheckTest.testDontAddMethodsIfNoConstructor()
func (c *TypeCheckTest) TestDontAddMethodsIfNoConstructor() {
	var js1Node *rhino.Node = c.ParseAndTypeCheck(`function A() {}
A.prototype = {m1: 5, m2: true}
`)
	var functionAType rhino.JSType = js1Node.GetFirstChild().GetJSType()
	AssertThat(functionAType.String()).IsEqualTo("function(): undefined")
	c.AssertTypeEquals(c.GetNativeUnknownType(), c.GetNativeFunctionType().GetPropertyType("m1"))
	c.AssertTypeEquals(c.GetNativeUnknownType(), c.GetNativeFunctionType().GetPropertyType("m2"))
}

// Java: TypeCheckTest.testFunctionAssignement()
func (c *TypeCheckTest) TestFunctionAssignement() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
* @param {string} ph0
* @param {string} ph1
* @return {string}
*/
function MSG_CALENDAR_ACCESS_ERROR(ph0, ph1) {return ''}
/** @type {Function} */
var MSG_CALENDAR_ADD_ERROR = MSG_CALENDAR_ACCESS_ERROR;
`).
		Run()
}

// Java: TypeCheckTest.testAddMethodsPrototypeTwoWays()
func (c *TypeCheckTest) TestAddMethodsPrototypeTwoWays() {
	var js1Node *rhino.Node = c.ParseAndTypeCheckWithExterns(CompilerTypeTestCase_DEFAULT_EXTERNS, `/** @constructor */function A() {}
A.prototype = {m1: 5, m2: true};
A.prototype.m3 = 'third property!';
`)
	var instanceType rhino.ObjectType = TypeCheckTestCase_GetInstanceType(js1Node)
	AssertThat(instanceType.String()).IsEqualTo("A")
	c.AssertHasXMorePropertiesThanNativeObject(instanceType, 3)
	c.CheckObjectType(instanceType, "m1", c.GetNativeNumberType())
	c.CheckObjectType(instanceType, "m2", c.GetNativeBooleanType())
	c.CheckObjectType(instanceType, "m3", c.GetNativeStringType())
}

// Java: TypeCheckTest.testPrototypePropertyTypes()
func (c *TypeCheckTest) TestPrototypePropertyTypes() {
	var js1Node *rhino.Node = c.ParseAndTypeCheckWithExterns(CompilerTypeTestCase_DEFAULT_EXTERNS, `/** @constructor */function A() {
  /** @type {string} */ this.m1;
  /** @type {Object?} */ this.m2 = {};
  /** @type {boolean} */ this.m3;
}
/** @type {string} */ A.prototype.m4;
/** @type {number} */ A.prototype.m5 = 0;
/** @type {boolean} */ A.prototype.m6;
`)
	var instanceType rhino.ObjectType = TypeCheckTestCase_GetInstanceType(js1Node)
	c.AssertHasXMorePropertiesThanNativeObject(instanceType, 6)
	c.CheckObjectType(instanceType, "m1", c.GetNativeStringType())
	c.CheckObjectType(instanceType, "m2", c.CreateUnionType(c.GetNativeObjectType(), c.GetNativeNullType()))
	c.CheckObjectType(instanceType, "m3", c.GetNativeBooleanType())
	c.CheckObjectType(instanceType, "m4", c.GetNativeStringType())
	c.CheckObjectType(instanceType, "m5", c.GetNativeNumberType())
	c.CheckObjectType(instanceType, "m6", c.GetNativeBooleanType())
}

// Java: TypeCheckTest.testValueTypeBuiltInPrototypePropertyType()
func (c *TypeCheckTest) TestValueTypeBuiltInPrototypePropertyType() {
	var node *rhino.Node = c.ParseAndTypeCheckWithExterns(NewTestExternsBuilder().AddString().Build(), "\"x\".charAt(0)")
	c.AssertTypeEquals(c.GetNativeStringType(), node.GetFirstFirstChild().GetJSType())
}

// Java: TypeCheckTest.testDeclareBuiltInConstructor()
func (c *TypeCheckTest) TestDeclareBuiltInConstructor() {
	var node *rhino.Node = c.ParseAndTypeCheckWithExterns(NewTestExternsBuilder().AddString().Build(), `/** @constructor */ var String = function(opt_str) {};
(new String("x")).charAt(0)
`)
	c.AssertTypeEquals(c.GetNativeStringType(), node.GetLastChild().GetFirstChild().
		GetJSType())
}

// Java: TypeCheckTest.testExtendBuiltInType1()
func (c *TypeCheckTest) TestExtendBuiltInType1() {
	externs := `/** @constructor */ var String = function(opt_str) {};
/**
* @param {number} start
* @param {number} opt_length
* @return {string}
*/
String.prototype.substr = function(start, opt_length) {};
`
	var n1 *rhino.Node = c.ParseAndTypeCheck(externs + "(new String(\"x\")).substr(0,1);")
	c.AssertTypeEquals(c.GetNativeStringType(), n1.GetLastChild().GetFirstChild().
		GetJSType())
}

// Java: TypeCheckTest.testExtendBuiltInType2()
func (c *TypeCheckTest) TestExtendBuiltInType2() {
	externs := `/** @constructor */ var String = function(opt_str) {};
/**
* @param {number} start
* @param {number} opt_length
* @return {string}
*/
String.prototype.substr = function(start, opt_length) {};
`
	var n2 *rhino.Node = c.ParseAndTypeCheck(externs + "\"x\".substr(0,1);")
	c.AssertTypeEquals(c.GetNativeStringType(), n2.GetLastChild().GetFirstChild().
		GetJSType())
}

// Java: TypeCheckTest.testExtendFunction1()
func (c *TypeCheckTest) TestExtendFunction1() {
	var n *rhino.Node = c.ParseAndTypeCheck(`/**@return {number}*/Function.prototype.f = function() { return 1; };
(new Function()).f();
`)
	var type_ rhino.JSType = n.GetLastChild().GetLastChild().
		GetJSType()
	c.AssertTypeEquals(c.GetNativeNumberType(), type_)
}

// Java: TypeCheckTest.testExtendFunction2()
func (c *TypeCheckTest) TestExtendFunction2() {
	var n *rhino.Node = c.ParseAndTypeCheck(`/**@return {number}*/Function.prototype.f = function() { return 1; };
(function() {}).f();
`)
	var type_ rhino.JSType = n.GetLastChild().GetLastChild().
		GetJSType()
	c.AssertTypeEquals(c.GetNativeNumberType(), type_)
}

// Java: TypeCheckTest.testClassExtendPrimitive()
func (c *TypeCheckTest) TestClassExtendPrimitive() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @extends {number} */ class C extends number {}").
		AddDiagnostic("C @extends non-object type number").
		Run()
}

// Java: TypeCheckTest.testInheritanceCheck1()
func (c *TypeCheckTest) TestInheritanceCheck1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */function Super() {};
/** @constructor
 @extends {Super} */function Sub() {};
Sub.prototype.foo = function() {};
`).
		Run()
}

// Java: TypeCheckTest.testInheritanceCheck2()
func (c *TypeCheckTest) TestInheritanceCheck2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */function Super() {};
/** @constructor
 @extends {Super} */function Sub() {};
/** @override */Sub.prototype.foo = function() {};
`).
		AddDiagnostic("property foo not defined on any superclass of Sub").
		Run()
}

// Java: TypeCheckTest.testInheritanceCheck3()
func (c *TypeCheckTest) TestInheritanceCheck3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */function Super() {};
Super.prototype.foo = function() {};
/** @constructor
 @extends {Super} */function Sub() {};
Sub.prototype.foo = function() {};
`).
		AddDiagnostic("property foo already defined on superclass Super; use @override to override it").
		Run()
}

// Java: TypeCheckTest.testInheritanceCheck4()
func (c *TypeCheckTest) TestInheritanceCheck4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */function Super() {};
Super.prototype.foo = function() {};
/** @constructor
 @extends {Super} */function Sub() {};
/** @override */Sub.prototype.foo = function() {};
`).
		Run()
}

// Java: TypeCheckTest.testInheritanceCheck5()
func (c *TypeCheckTest) TestInheritanceCheck5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */function Root() {};
Root.prototype.foo = function() {};
/** @constructor
 @extends {Root} */function Super() {};
/** @constructor
 @extends {Super} */function Sub() {};
Sub.prototype.foo = function() {};
`).
		AddDiagnostic("property foo already defined on superclass Root; use @override to override it").
		Run()
}

// Java: TypeCheckTest.testInheritanceCheck6()
func (c *TypeCheckTest) TestInheritanceCheck6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */function Root() {};
Root.prototype.foo = function() {};
/** @constructor
 @extends {Root} */function Super() {};
/** @constructor
 @extends {Super} */function Sub() {};
/** @override */Sub.prototype.foo = function() {};
`).
		Run()
}

// Java: TypeCheckTest.testInheritanceCheck7()
func (c *TypeCheckTest) TestInheritanceCheck7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = {};
/** @constructor */goog.Super = function() {};
goog.Super.prototype.foo = 3;
/** @constructor
 @extends {goog.Super} */goog.Sub = function() {};
goog.Sub.prototype.foo = 5;
`).
		Run()
}

// Java: TypeCheckTest.testInheritanceCheck8()
func (c *TypeCheckTest) TestInheritanceCheck8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = {};
/** @constructor */goog.Super = function() {};
goog.Super.prototype.foo = 3;
/** @constructor
 @extends {goog.Super} */goog.Sub = function() {};
/** @override */goog.Sub.prototype.foo = 5;
`).
		Run()
}

// Java: TypeCheckTest.testInheritanceCheck9_1()
func (c *TypeCheckTest) TestInheritanceCheck9_1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */function Super() {};
Super.prototype.foo = function() { return 3; };
/** @constructor
 @extends {Super} */function Sub() {};
/** @override
 @return {number} */Sub.prototype.foo =
function() { return 1; };
`).
		Run()
}

// Java: TypeCheckTest.testInheritanceCheck9_2()
func (c *TypeCheckTest) TestInheritanceCheck9_2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */function Super() {};
/** @return {number} */
Super.prototype.foo = function() { return 1; };
/** @constructor
 @extends {Super} */function Sub() {};
/** @override */Sub.prototype.foo =
function() {};
`).
		Run()
}

// Java: TypeCheckTest.testInheritanceCheck9_3()
func (c *TypeCheckTest) TestInheritanceCheck9_3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */function Super() {};
/** @return {number} */
Super.prototype.foo = function() { return 1; };
/** @constructor
 @extends {Super} */function Sub() {};
/** @override
 @return {string} */Sub.prototype.foo =
function() { return "some string" };
`).
		AddDiagnostic(`mismatch of the foo property type and the type of the property it overrides from superclass Super
original: function(this:Super): number
override: function(this:Sub): string
`).
		Run()
}

// Java: TypeCheckTest.testInheritanceCheck10_1()
func (c *TypeCheckTest) TestInheritanceCheck10_1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */function Root() {};
Root.prototype.foo = function() { return 3; };
/** @constructor
 @extends {Root} */function Super() {};
/** @constructor
 @extends {Super} */function Sub() {};
/** @override
 @return {number} */Sub.prototype.foo =
function() { return 1; };
`).
		Run()
}

// Java: TypeCheckTest.testInheritanceCheck10_2()
func (c *TypeCheckTest) TestInheritanceCheck10_2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */function Root() {};
/** @return {number} */
Root.prototype.foo = function() { return 1; };
/** @constructor
 @extends {Root} */function Super() {};
/** @constructor
 @extends {Super} */function Sub() {};
/** @override */Sub.prototype.foo =
function() {};
`).
		Run()
}

// Java: TypeCheckTest.testInheritanceCheck10_3()
func (c *TypeCheckTest) TestInheritanceCheck10_3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */function Root() {};
/** @return {number} */
Root.prototype.foo = function() { return 1; };
/** @constructor
 @extends {Root} */function Super() {};
/** @constructor
 @extends {Super} */function Sub() {};
/** @override
 @return {string} */Sub.prototype.foo =
function() { return "some string" };
`).
		AddDiagnostic(`mismatch of the foo property type and the type of the property it overrides from superclass Root
original: function(this:Root): number
override: function(this:Sub): string
`).
		Run()
}

// Java: TypeCheckTest.testInterfaceInheritanceCheck11()
func (c *TypeCheckTest) TestInterfaceInheritanceCheck11() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */function Super() {};
/** @param {number} bar */Super.prototype.foo = function(bar) {};
/** @constructor
 @extends {Super} */function Sub() {};
/** @override
  @param {string} bar */Sub.prototype.foo =
function(bar) {};
`).
		AddDiagnostic(`mismatch of the foo property type and the type of the property it overrides from superclass Super
original: function(this:Super, number): undefined
override: function(this:Sub, string): undefined
`).
		Run()
}

// Java: TypeCheckTest.testInheritanceCheck12()
func (c *TypeCheckTest) TestInheritanceCheck12() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = {};
/** @constructor */goog.Super = function() {};
goog.Super.prototype.foo = 3;
/** @constructor
 @extends {goog.Super} */goog.Sub = function() {};
/** @override */goog.Sub.prototype.foo = "some string";
`).
		Run()
}

// Java: TypeCheckTest.testInheritanceCheck13()
func (c *TypeCheckTest) TestInheritanceCheck13() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = {};
/** @constructor
 @extends {goog.Missing} */function Sub() {};
/** @override */Sub.prototype.foo = function() {};
`).
		AddDiagnostic(`Bad type annotation. Unknown type goog.Missing
It's possible that 'goog.Missing' refers to a value, not a type.
`).
		Run()
}

// Java: TypeCheckTest.testInheritanceCheck14()
func (c *TypeCheckTest) TestInheritanceCheck14() {
	c.TestClosureTypes(`/** @constructor
 @extends {goog.Missing} */
goog.Super = function() {};
/** @constructor
 @extends {goog.Super} */function Sub() {};
/** @override */ Sub.prototype.foo = function() {};
`, jutil.StringPtr(`Bad type annotation. Unknown type goog.Missing
It's possible that 'goog.Missing' refers to a value, not a type.
`))
}

// Java: TypeCheckTest.testInheritanceCheck15()
func (c *TypeCheckTest) TestInheritanceCheck15() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */function Super() {};
/** @param {number} bar */Super.prototype.foo;
/** @constructor
 @extends {Super} */function Sub() {};
/** @override
  @param {number} bar */Sub.prototype.foo =
function(bar) {};
`).
		Run()
}

// Java: TypeCheckTest.testInheritanceCheck16()
func (c *TypeCheckTest) TestInheritanceCheck16() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = {};
/** @constructor */goog.Super = function() {};
/** @type {number} */ goog.Super.prototype.foo = 3;
/** @constructor
 @extends {goog.Super} */goog.Sub = function() {};
/** @type {number} */ goog.Sub.prototype.foo = 5;
`).
		AddDiagnostic(`property foo already defined on superclass goog.Super; use @override to override it
`).
		Run()
}

// Java: TypeCheckTest.testInheritanceCheck17()
func (c *TypeCheckTest) TestInheritanceCheck17() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = {};
/** @constructor */goog.Super = function() {};
/** @param {number} x */ goog.Super.prototype.foo = function(x) {};
/** @constructor
 @extends {goog.Super} */goog.Sub = function() {};
/** @override @param {string} x */ goog.Sub.prototype.foo = function(x) {};
`).
		AddDiagnostic(`mismatch of the foo property type and the type of the property it overrides from superclass goog.Super
original: function(this:goog.Super, number): undefined
override: function(this:goog.Sub, string): undefined
`).
		Run()
}

// Java: TypeCheckTest.testInterfacePropertyOverride1()
func (c *TypeCheckTest) TestInterfacePropertyOverride1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */function Super() {};
Super.prototype.foo = function() {};
/** @interface
 @extends {Super} */function Sub() {};
Sub.prototype.foo = function() {};
`).
		Run()
}

// Java: TypeCheckTest.testInterfacePropertyOverride2()
func (c *TypeCheckTest) TestInterfacePropertyOverride2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */function Root() {};
Root.prototype.foo = function() {};
/** @interface
 @extends {Root} */function Super() {};
/** @interface
 @extends {Super} */function Sub() {};
Sub.prototype.foo = function() {};
`).
		Run()
}

// Java: TypeCheckTest.testInterfacePropertyBadOverrideFails()
func (c *TypeCheckTest) TestInterfacePropertyBadOverrideFails() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */function Super() {};
/** @type {number} */
Super.prototype.foo;
/** @interface @extends {Super} */function Sub() {};
/** @type {string} */
Sub.prototype.foo;
`).
		AddDiagnostic(`mismatch of the foo property on type Sub and the type of the property it overrides from interface Super
original: number
override: string
`).
		Run()
}

// Java: TypeCheckTest.testInterfaceInheritanceCheck1()
func (c *TypeCheckTest) TestInterfaceInheritanceCheck1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */function Super() {};
Super.prototype.foo = function() {};
/** @constructor
 @implements {Super} */function Sub() {};
Sub.prototype.foo = function() {};
`).
		AddDiagnostic("property foo already defined on interface Super; use @override to override it").
		Run()
}

// Java: TypeCheckTest.testInterfaceInheritanceCheck2()
func (c *TypeCheckTest) TestInterfaceInheritanceCheck2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */function Super() {};
Super.prototype.foo = function() {};
/** @constructor
 @implements {Super} */function Sub() {};
/** @override */Sub.prototype.foo = function() {};
`).
		Run()
}

// Java: TypeCheckTest.testInterfaceInheritanceCheck3()
func (c *TypeCheckTest) TestInterfaceInheritanceCheck3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */function Root() {};
/** @return {number} */Root.prototype.foo = function() {};
/** @interface
 @extends {Root} */function Super() {};
/** @constructor
 @implements {Super} */function Sub() {};
/** @return {number} */Sub.prototype.foo = function() { return 1;};
`).
		AddDiagnostic("property foo already defined on interface Root; use @override to override it").
		Run()
}

// Java: TypeCheckTest.testInterfaceInheritanceCheck4()
func (c *TypeCheckTest) TestInterfaceInheritanceCheck4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */function Root() {};
/** @return {number} */Root.prototype.foo = function() {};
/** @interface
 @extends {Root} */function Super() {};
/** @constructor
 @implements {Super} */function Sub() {};
/** @override
 * @return {number} */Sub.prototype.foo =
function() { return 1;};
`).
		Run()
}

// Java: TypeCheckTest.testInterfaceInheritanceCheck5()
func (c *TypeCheckTest) TestInterfaceInheritanceCheck5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */function Super() {};/** @return {string} */Super.prototype.foo = function() {};/** @constructor
 @implements {Super} */function Sub() {};/** @override
 @return {number} */Sub.prototype.foo = function() { return 1; };
`).
		AddDiagnostic(`mismatch of the foo property on type Sub and the type of the property it overrides from interface Super
original: function(this:Super): string
override: function(this:Sub): number
`).
		Run()
}

// Java: TypeCheckTest.testInterfaceInheritanceCheck6()
func (c *TypeCheckTest) TestInterfaceInheritanceCheck6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */function Root() {};/** @return {string} */Root.prototype.foo = function() {};/** @interface
 @extends {Root} */function Super() {};/** @constructor
 @implements {Super} */function Sub() {};/** @override
 @return {number} */Sub.prototype.foo = function() { return 1; };
`).
		AddDiagnostic(`mismatch of the foo property on type Sub and the type of the property it overrides from interface Root
original: function(this:Root): string
override: function(this:Sub): number
`).
		Run()
}

// Java: TypeCheckTest.testInterfaceInheritanceCheck7()
func (c *TypeCheckTest) TestInterfaceInheritanceCheck7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */function Super() {};
/** @param {number} bar */Super.prototype.foo = function(bar) {};
/** @constructor
 @implements {Super} */function Sub() {};
/** @override
  @param {string} bar */Sub.prototype.foo =
function(bar) {};
`).
		AddDiagnostic(`mismatch of the foo property on type Sub and the type of the property it overrides from interface Super
original: function(this:Super, number): undefined
override: function(this:Sub, string): undefined
`).
		Run()
}

// Java: TypeCheckTest.testInterfaceInheritanceCheck8()
func (c *TypeCheckTest) TestInterfaceInheritanceCheck8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor
 @implements {Super} */function Sub() {};
/** @override */Sub.prototype.foo = function() {};
`).
		AddDiagnostic("Bad type annotation. Unknown type Super").
		AddDiagnostic("property foo not defined on any superclass of Sub").
		Run()
}

// Java: TypeCheckTest.testInterfaceInheritanceCheck9()
func (c *TypeCheckTest) TestInterfaceInheritanceCheck9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ function I() {}
/** @return {number} */ I.prototype.bar = function() {};
/** @constructor */ function F() {}
/** @return {number} */ F.prototype.bar = function() {return 3; };
/** @return {number} */ F.prototype.foo = function() {return 3; };
/** @constructor
 * @extends {F}
 * @implements {I} */
function G() {}
/** @return {string} */ function f() { return new G().bar(); }
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testInterfaceInheritanceCheck10()
func (c *TypeCheckTest) TestInterfaceInheritanceCheck10() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ function I() {}
/** @return {number} */ I.prototype.bar = function() {};
/** @constructor */ function F() {}
/** @return {number} */ F.prototype.foo = function() {return 3; };
/** @constructor
 * @extends {F}
 * @implements {I} */
function G() {}
/** @return {number}
 * @override */
G.prototype.bar = G.prototype.foo;
/** @return {string} */ function f() { return new G().bar(); }
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testInterfaceInheritanceCheck12()
func (c *TypeCheckTest) TestInterfaceInheritanceCheck12() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ function I() {};
/** @type {string} */ I.prototype.foobar;
/** 
 * @constructor 
 * @implements {I} */
function C() {
/** 
 * @type {number} */ this.foobar = 2;};
/** @type {I} */ 
 var test = new C(); alert(test.foobar);
`).
		AddDiagnostic(`mismatch of the foobar property on type C and the type of the property it overrides from interface I
original: string
override: number
`).
		Run()
}

// Java: TypeCheckTest.testInterfaceInheritanceCheck13()
func (c *TypeCheckTest) TestInterfaceInheritanceCheck13() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function abstractMethod() {};
/** @interface */var base = function() {};
/** @extends {base} 
 @interface */ var Int = function() {}
/** @type {{bar : !Function}} */ var x; 
/** @type {!Function} */ base.prototype.bar = abstractMethod; 
/** @type {Int} */ var foo;
foo.bar();
`).
		Run()
}

// Java: TypeCheckTest.testInterfaceInheritanceCheck14()
func (c *TypeCheckTest) TestInterfaceInheritanceCheck14() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface
 @template T */function A() {};
/** @return {T} */A.prototype.foo = function() {};
/** @interface
 @template U
 @extends {A<U>} */function B() {};
/** @return {U} */B.prototype.bar = function() {};
/** @constructor
 @implements {B<string>} */function C() {};
/** @return {string}
 @override */C.prototype.foo = function() {};
/** @return {string}
 @override */C.prototype.bar = function() {};
`).
		Run()
}

// Java: TypeCheckTest.testInterfaceInheritanceCheck15()
func (c *TypeCheckTest) TestInterfaceInheritanceCheck15() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface
 @template T */function A() {};
/** @return {T} */A.prototype.foo = function() {};
/** @interface
 @template U
 @extends {A<U>} */function B() {};
/** @return {U} */B.prototype.bar = function() {};
/** @constructor
 @template V
 @implements {B<V>}
 */function C() {};
/** @return {V}
 @override */C.prototype.foo = function() {};
/** @return {V}
 @override */C.prototype.bar = function() {};
`).
		Run()
}

// Java: TypeCheckTest.testInterfaceInheritanceCheck16()
func (c *TypeCheckTest) TestInterfaceInheritanceCheck16() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface
 @template T */function A() {};
/** @return {T} */A.prototype.foo = function() {};
/** @return {T} */A.prototype.bar = function() {};
/** @constructor
 @implements {A<string>} */function B() {};
/** @override */B.prototype.foo = function() { return 'string'};
/** @override */B.prototype.bar = function() { return 3 };
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testInterfacePropertyNotImplemented()
func (c *TypeCheckTest) TestInterfacePropertyNotImplemented() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */function Int() {};
Int.prototype.foo = function() {};
/** @constructor
 @implements {Int} */function Foo() {};
`).
		AddDiagnostic("property foo on interface Int is not implemented by type Foo").
		Run()
}

// Java: TypeCheckTest.testInterfacePropertyNotImplemented2()
func (c *TypeCheckTest) TestInterfacePropertyNotImplemented2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */function Int() {};
Int.prototype.foo = function() {};
/** @interface
 @extends {Int} */function Int2() {};
/** @constructor
 @implements {Int2} */function Foo() {};
`).
		AddDiagnostic("property foo on interface Int is not implemented by type Foo").
		Run()
}

// Java: TypeCheckTest.testInterfacePropertyNotImplemented3()
func (c *TypeCheckTest) TestInterfacePropertyNotImplemented3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface  @template T */ function Int() {};
/** @return {T} */ Int.prototype.foo = function() {};

/** @constructor @implements {Int<string>} */ function Foo() {};
/** @return {number}  @override */ Foo.prototype.foo = function() {};
`).
		AddDiagnostic(`mismatch of the foo property on type Foo and the type of the property it overrides from interface Int
original: function(this:Int): string
override: function(this:Foo): number
`).
		Run()
}

// Java: TypeCheckTest.testStubConstructorImplementingInterface()
func (c *TypeCheckTest) TestStubConstructorImplementingInterface() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @interface */ function Int() {}
Int.prototype.foo = function() {};
/** @constructor
 @implements {Int} */ var Foo;
`).
		AddSource("").
		Run()
}

// Java: TypeCheckTest.testObjectLiteral()
func (c *TypeCheckTest) TestObjectLiteral() {
	var n *rhino.Node = c.ParseAndTypeCheck("var a = {m1: 7, m2: 'hello'}")
	var nameNode *rhino.Node = n.GetFirstFirstChild()
	var objectNode *rhino.Node = nameNode.GetFirstChild()
	NodeSubject_AssertNode(nameNode).HasToken(rhino.Token_NAME)
	NodeSubject_AssertNode(objectNode).HasToken(rhino.Token_OBJECTLIT)
	var objectType rhino.ObjectType = objectNode.GetJSType().(rhino.ObjectType)
	c.AssertTypeEquals(c.GetNativeNumberType(), objectType.GetPropertyType("m1"))
	c.AssertTypeEquals(c.GetNativeStringType(), objectType.GetPropertyType("m2"))
	c.AssertTypeEquals(objectType, nameNode.GetJSType())
}

// Java: TypeCheckTest.testObjectLiteralDeclaration1()
func (c *TypeCheckTest) TestObjectLiteralDeclaration1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var x = {
  /** @type {boolean} */ abc: true,
  /** @type {number} */ 'def': 0,
  /** @type {string} */ 3: 'fgh'
};
`).
		Run()
}

// Java: TypeCheckTest.testObjectLiteralDeclaration2()
func (c *TypeCheckTest) TestObjectLiteralDeclaration2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var x = {
  /** @type {boolean} */ abc: true
};
x.abc = 0;
`).
		AddDiagnostic(`assignment to property abc of x
found   : number
required: boolean
`).
		Run()
}

// Java: TypeCheckTest.testObjectLiteralDeclaration3()
func (c *TypeCheckTest) TestObjectLiteralDeclaration3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {{foo: !Function}} x */ function f(x) {}
f({foo: function() {}});
`).
		Run()
}

// Java: TypeCheckTest.testObjectLiteralDeclaration4()
func (c *TypeCheckTest) TestObjectLiteralDeclaration4() {
	c.TestClosureTypes(`var x = {
  /** @param {boolean} x */ abc: function(x) {}
};
/**
 * @param {string} x
 * @suppress {duplicate}
 */ x.abc = function(x) {};
`, jutil.StringPtr(`assignment to property abc of x
found   : function(string): undefined
required: function(boolean): undefined
`))
}

// Java: TypeCheckTest.testObjectLiteralDeclaration5()
func (c *TypeCheckTest) TestObjectLiteralDeclaration5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var x = {
  /** @param {boolean} x */ abc: function(x) {}
};
/**
 * @param {boolean} x
 * @suppress {duplicate}
 */ x.abc = function(x) {};
`).
		Run()
}

// Java: TypeCheckTest.testObjectLiteralDeclaration6()
func (c *TypeCheckTest) TestObjectLiteralDeclaration6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var x = {};
/**
 * @param {boolean} x
 * @suppress {duplicate}
 */ x.abc = function(x) {};
x = {
  /**
   * @param {boolean} x
   * @suppress {duplicate}
   */
  abc: function(x) {}
};
`).
		Run()
}

// Java: TypeCheckTest.testObjectLiteralDeclaration7()
func (c *TypeCheckTest) TestObjectLiteralDeclaration7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var x = {};
/**
 * @type {function(boolean): undefined}
 */ x.abc = function(x) {};
x = {
  /**
   * @param {boolean} x
   * @suppress {duplicate}
   */
  abc: function(x) {}
};
`).
		Run()
}

// Java: TypeCheckTest.testCallDateConstructorAsFunction()
func (c *TypeCheckTest) TestCallDateConstructorAsFunction() {
	var n *rhino.Node = c.ParseAndTypeCheck("Date()")
	c.AssertTypeEquals(c.GetNativeStringType(), n.GetFirstFirstChild().GetJSType())
}

// Java: TypeCheckTest.testCallErrorConstructorAsFunction()
func (c *TypeCheckTest) TestCallErrorConstructorAsFunction() {
	externs := `/** @constructor
    @param {string} message
    @return {!Error} */
function Error(message) {}
`
	var n *rhino.Node = c.ParseAndTypeCheckWithExterns(externs, "Error('x')")
	var call *rhino.Node = n.GetFirstFirstChild()
	AssertThat(call.IsCall()).IsTrue()
	c.AssertTypeEquals(call.GetFirstChild().GetJSType().
		ToMaybeFunctionType().
		GetInstanceType(), call.GetJSType())
}

// Java: TypeCheckTest.testCallArrayConstructorAsFunction()
func (c *TypeCheckTest) TestCallArrayConstructorAsFunction() {
	var n *rhino.Node = c.ParseAndTypeCheck("Array()")
	c.AssertTypeEquals(c.GetNativeArrayType(), n.GetFirstFirstChild().GetJSType())
}

// Java: TypeCheckTest.testPropertyTypeOfUnionType()
func (c *TypeCheckTest) TestPropertyTypeOfUnionType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var a = {};
/** @constructor */ a.N = function() {};
a.N.prototype.p = 1;
/** @constructor */ a.S = function() {};
a.S.prototype.p = 'a';
/** @param {!a.N|!a.S} x
@return {string} */
var f = function(x) { return x.p; };
`).
		AddDiagnostic(`inconsistent return type
found   : (number|string)
required: string
`).
		Run()
}

// Java: TypeCheckTest.testGetPropertyTypeOfUnionType_withMatchingTemplates()
func (c *TypeCheckTest) TestGetPropertyTypeOfUnionType_withMatchingTemplates() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface @template T */ function Foo() {};
/** @type {T} */
Foo.prototype.p;
/** @interface @template U */ function Bar() {};
/** @type {U} */
Bar.prototype.p;

/**
 * @param {!Foo<number>|!Bar<number>} x
 * @return {string}
 */
var f = function(x) { return x.p; };
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testGetPropertyTypeOfUnionType_withDifferingTemplates()
func (c *TypeCheckTest) TestGetPropertyTypeOfUnionType_withDifferingTemplates() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface @template T */ function Foo() {};
/** @type {T} */
Foo.prototype.p;
/** @interface @template U */ function Bar() {};
/** @type {U} */
Bar.prototype.p;

/**
 * @param {!Foo<number>|!Bar<string>} x
 * @return {string}
 */
var f = function(x) { return x.p; };
`).
		AddDiagnostic(`inconsistent return type
found   : (number|string)
required: string
`).
		Run()
}

// Java: TypeCheckTest.testGetPropertyTypeOfUnionType_withMembersThatExtendATemplatizedType()
func (c *TypeCheckTest) TestGetPropertyTypeOfUnionType_withMembersThatExtendATemplatizedType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface @template T */ function Foo() {};
/** @type {T} */
Foo.prototype.p;

/** @interface @extends {Foo<number>} */ function Bar() {};
/** @interface @extends {Foo<number>} */ function Baz() {}

/**
 * @param {!Bar|!Baz} x
 * @return {string}
 */
var f = function(x) { return x.p; };
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testInvalidAssignToPropertyTypeOfUnionType_withMatchingTemplates_doesntWarn()
func (c *TypeCheckTest) TestInvalidAssignToPropertyTypeOfUnionType_withMatchingTemplates_doesntWarn() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface @template T */ function Foo() {};
/** @type {T} */
Foo.prototype.p;
/** @interface @template U */ function Bar() {};
/** @type {U} */
Bar.prototype.p;

/**
 * @param {!Foo<number>|!Bar<number>} x
 */
var f = function(x) { x.p = 'not a number'; };
`).
		Run()
}

// Java: TypeCheckTest.testWarningOnInterfacePrototype()
func (c *TypeCheckTest) TestWarningOnInterfacePrototype() {
	c.T().Skip("@Ignore")
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ u.T = function() {};
/** @return {number} */ u.T.prototype = function() { };
`).
		AddDiagnostic("e of its definition").
		Run()
}

// Java: TypeCheckTest.testBadPropertyOnInterface1()
func (c *TypeCheckTest) TestBadPropertyOnInterface1() {
	c.T().Skip("@Ignore")
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ u.T = function() {};
/** @return {number} */ u.T.f = function() { return 1;};
`).
		AddDiagnostic("cannot reference an interface outside of its definition").
		Run()
}

// Java: TypeCheckTest.testBadPropertyOnInterface2()
func (c *TypeCheckTest) TestBadPropertyOnInterface2() {
	c.T().Skip("@Ignore")
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ function T() {};
/** @return {number} */ T.f = function() { return 1;};
`).
		AddDiagnostic("cannot reference an interface outside of its definition").
		Run()
}

// Java: TypeCheckTest.testBadPropertyOnInterface3()
func (c *TypeCheckTest) TestBadPropertyOnInterface3() {
	c.T().Skip("@Ignore")
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @interface */ u.T = function() {}; u.T.x").
		AddDiagnostic("cannot reference an interface outside of its definition").
		Run()
}

// Java: TypeCheckTest.testBadPropertyOnInterface4()
func (c *TypeCheckTest) TestBadPropertyOnInterface4() {
	c.T().Skip("@Ignore")
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @interface */ function T() {}; T.x;").
		AddDiagnostic("cannot reference an interface outside of its definition").
		Run()
}

// Java: TypeCheckTest.testAnnotatedPropertyOnInterface1()
func (c *TypeCheckTest) TestAnnotatedPropertyOnInterface1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ u.T = function() {};
/** @return {number} */ u.T.prototype.f = function() {};
`).
		Run()
}

// Java: TypeCheckTest.testAnnotatedPropertyOnInterface2()
func (c *TypeCheckTest) TestAnnotatedPropertyOnInterface2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ u.T = function() {};
/** @return {number} */ u.T.prototype.f = function() { };
`).
		Run()
}

// Java: TypeCheckTest.testAnnotatedPropertyOnInterface3()
func (c *TypeCheckTest) TestAnnotatedPropertyOnInterface3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ function T() {};
/** @return {number} */ T.prototype.f = function() { };
`).
		Run()
}

// Java: TypeCheckTest.testAnnotatedPropertyOnInterface4()
func (c *TypeCheckTest) TestAnnotatedPropertyOnInterface4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(CompilerTypeTestCase_CLOSURE_DEFS + `/** @interface */ function T() {};
/** @return {number} */ T.prototype.f = goog.abstractMethod;
`).
		Run()
}

// Java: TypeCheckTest.testWarnUnannotatedPropertyOnInterface1()
func (c *TypeCheckTest) TestWarnUnannotatedPropertyOnInterface1() {
	c.T().Skip("@Ignore")
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @interface */ u.T = function() {}; u.T.prototype.x;").
		AddDiagnostic("interface property x is not annotated").
		Run()
}

// Java: TypeCheckTest.testWarnUnannotatedPropertyOnInterface2()
func (c *TypeCheckTest) TestWarnUnannotatedPropertyOnInterface2() {
	c.T().Skip("@Ignore")
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @interface */ function T() {}; T.prototype.x;").
		AddDiagnostic("interface property x is not annotated").
		Run()
}

// Java: TypeCheckTest.testWarnUnannotatedPropertyOnInterface5()
func (c *TypeCheckTest) TestWarnUnannotatedPropertyOnInterface5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ u.T = function() {};
u.T.prototype.x = function() {};
`).
		Run()
}

// Java: TypeCheckTest.testWarnUnannotatedPropertyOnInterface6()
func (c *TypeCheckTest) TestWarnUnannotatedPropertyOnInterface6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ function T() {};
T.prototype.x = function() {};
`).
		Run()
}

// Java: TypeCheckTest.testWarnDataPropertyOnInterface1()
func (c *TypeCheckTest) TestWarnDataPropertyOnInterface1() {
	c.T().Skip("@Ignore")
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ u.T = function() {};
/** @type {number} */u.T.prototype.x;
`).
		AddDiagnostic("interface members can only be plain functions").
		Run()
}

// Java: TypeCheckTest.testDataPropertyOnInterface1()
func (c *TypeCheckTest) TestDataPropertyOnInterface1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ function T() {};
/** @type {number} */T.prototype.x;
`).
		Run()
}

// Java: TypeCheckTest.testDataPropertyOnInterface2()
func (c *TypeCheckTest) TestDataPropertyOnInterface2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ function T() {};
/** @type {number} */T.prototype.x;
/** @constructor 
 * @implements {T} 
 */
function C() {}
/** @override */
C.prototype.x = 'foo';
`).
		AddDiagnostic(`mismatch of the x property on type C and the type of the property it overrides from interface T
original: number
override: string
`).
		Run()
}

// Java: TypeCheckTest.testDataPropertyOnInterface3()
func (c *TypeCheckTest) TestDataPropertyOnInterface3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ function T() {};
/** @type {number} */T.prototype.x;
/** @constructor 
 * @implements {T} 
 */
function C() {}
/** @override */
C.prototype.x = 'foo';
`).
		AddDiagnostic(`mismatch of the x property on type C and the type of the property it overrides from interface T
original: number
override: string
`).
		Run()
}

// Java: TypeCheckTest.testDataPropertyOnInterface4()
func (c *TypeCheckTest) TestDataPropertyOnInterface4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ function T() {};
/** @type {number} */T.prototype.x;
/** @constructor 
 * @implements {T} 
 */
function C() { /** @type {string} */ 
 this.x = 'foo'; }
`).
		AddDiagnostic(`mismatch of the x property on type C and the type of the property it overrides from interface T
original: number
override: string
`).
		Run()
}

// Java: TypeCheckTest.testWarnDataPropertyOnInterface3()
func (c *TypeCheckTest) TestWarnDataPropertyOnInterface3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ u.T = function() {};
/** @type {number} */u.T.prototype.x = 1;
`).
		AddDiagnostic(`interface members can only be empty property declarations, empty functions, or goog.abstractMethod
`).
		Run()
}

// Java: TypeCheckTest.testWarnDataPropertyOnInterface4()
func (c *TypeCheckTest) TestWarnDataPropertyOnInterface4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ function T() {};
/** @type {number} */T.prototype.x = 1;
`).
		AddDiagnostic(`interface members can only be empty property declarations, empty functions, or goog.abstractMethod
`).
		Run()
}

// Java: TypeCheckTest.testErrorMismatchingPropertyOnInterface1()
func (c *TypeCheckTest) TestErrorMismatchingPropertyOnInterface1() {
	c.T().Skip("@Ignore")
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ u.T = function() {};
/** @param {Number} foo */u.T.prototype.x =
/** @param {String} foo */function(foo) {};
`).
		AddDiagnostic(`found   : 
required: 
`).
		Run()
}

// Java: TypeCheckTest.testErrorMismatchingPropertyOnInterface2()
func (c *TypeCheckTest) TestErrorMismatchingPropertyOnInterface2() {
	c.T().Skip("@Ignore")
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ function T() {};
/** @return {number} */T.prototype.x =
/** @return {string} */function() {};
`).
		AddDiagnostic(`found   : 
required: 
`).
		Run()
}

// Java: TypeCheckTest.testErrorMismatchingPropertyOnInterface3()
func (c *TypeCheckTest) TestErrorMismatchingPropertyOnInterface3() {
	c.T().Skip("@Ignore")
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ u.T = function() {};
/** @param {Number} foo */u.T.prototype.x =
function(foo, bar) {};
`).
		AddDiagnostic(`found   : 
required: 
`).
		Run()
}

// Java: TypeCheckTest.testErrorMismatchingPropertyOnInterface4()
func (c *TypeCheckTest) TestErrorMismatchingPropertyOnInterface4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ u.T = function() {};
/** @param {Number} foo */u.T.prototype.x =
function() {};
`).
		AddDiagnostic("parameter foo does not appear in u.T.prototype.x's parameter list").
		Run()
}

// Java: TypeCheckTest.testErrorMismatchingPropertyOnInterface5()
func (c *TypeCheckTest) TestErrorMismatchingPropertyOnInterface5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ function T() {};
/** @type {number} */T.prototype.x = function() { };
`).
		AddDiagnostic(`assignment to property x of T.prototype
found   : function(): undefined
required: number`).
		Run()
}

// Java: TypeCheckTest.testErrorMismatchingPropertyOnInterface6()
func (c *TypeCheckTest) TestErrorMismatchingPropertyOnInterface6() {
	c.TestClosureTypesMultipleWarnings(`/** @interface */ function T() {};
/** @return {number} */T.prototype.x = 1
`, []string{`assignment to property x of T.prototype
found   : number
required: function(this:T): number`, "interface members can only be empty property declarations, empty functions, or goog.abstractMethod"})
}

// Java: TypeCheckTest.testInterfaceNonEmptyFunction()
func (c *TypeCheckTest) TestInterfaceNonEmptyFunction() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ function T() {};
T.prototype.x = function() { return 'foo'; }
`).
		AddDiagnostic("interface member functions must have an empty body").
		Run()
}

// Java: TypeCheckTest.testDoubleNestedInterface()
func (c *TypeCheckTest) TestDoubleNestedInterface() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ var I1 = function() {};
/** @interface */ I1.I2 = function() {};
/** @interface */ I1.I2.I3 = function() {};
`).
		Run()
}

// Java: TypeCheckTest.testStaticDataPropertyOnNestedInterface()
func (c *TypeCheckTest) TestStaticDataPropertyOnNestedInterface() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ var I1 = function() {};
/** @interface */ I1.I2 = function() {};
/** @type {number} */ I1.I2.x = 1;
`).
		Run()
}

// Java: TypeCheckTest.testInterfaceInstantiation()
func (c *TypeCheckTest) TestInterfaceInstantiation() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @interface */var f = function(){}; new f").
		AddDiagnostic("cannot instantiate non-constructor, found type: (typeof f)").
		Run()
}

// Java: TypeCheckTest.testPrototypeLoop()
func (c *TypeCheckTest) TestPrototypeLoop() {
	c.TestClosureTypesMultipleWarnings(c.SuppressMissingProperty("foo")+`/** @constructor
 * @extends {T} */var T = function() {};
alert((new T).foo);
`, []string{"Cycle detected in inheritance chain of type T", "Could not resolve type in @extends tag of T"})
}

// Java: TypeCheckTest.testImplementsLoop()
func (c *TypeCheckTest) TestImplementsLoop() {
	c.TestClosureTypesMultipleWarnings(strings.ReplaceAll(`/** @constructor
 * @implements {T} */var T = function() {};
SUPPRESSION
alert((new T).foo);
`, "SUPPRESSION", c.SuppressMissingPropertyFor("T", "foo")), []string{"Cycle detected in inheritance chain of type T"})
}

// Java: TypeCheckTest.testImplementsExtendsLoop()
func (c *TypeCheckTest) TestImplementsExtendsLoop() {
	c.TestClosureTypesMultipleWarnings(`/**
 * @constructor
 * @implements {F}
 */
var G = function() {};
/**
 * @constructor
 * @extends {G}
 */
var F = function() {};
alert((new F).foo);
`, []string{"Cycle detected in inheritance chain of type F", "Property foo never defined on F"})
}

// Java: TypeCheckTest.testInterfaceExtendsLoop()
func (c *TypeCheckTest) TestInterfaceExtendsLoop() {
	c.T().Skip("@Ignore")
	c.TestClosureTypesMultipleWarnings(strings.ReplaceAll(`/** @interface
 * @extends {F} */var G = function() {};
/** @interface
 * @extends {G} */var F = function() {};
/** @constructor
 * @implements {F} */var H = function() {};
SUPPRESSION
alert((new H).foo);
`, "SUPPRESSION", c.SuppressMissingPropertyFor("H", "foo")), []string{"extends loop involving F, loop: F -> G -> F", "extends loop involving G, loop: G -> F -> G"})
}

// Java: TypeCheckTest.testInterfaceExtendsLoop2()
func (c *TypeCheckTest) TestInterfaceExtendsLoop2() {
	c.TestClosureTypesMultipleWarnings(strings.ReplaceAll(`/** @record
 * @extends {F} */var G = function() {};
/** @record
 * @extends {G} */var F = function() {};
/** @constructor
 * @implements {F} */var H = function() {};
SUPPRESSION
alert((new H).foo);
`, "SUPPRESSION", c.SuppressMissingPropertyFor("H", "foo")), []string{
		"Cycle detected in inheritance chain of type F",
		"Could not resolve type in @extends tag of G"})
}

// Java: TypeCheckTest.testInheritPropFromMultipleInterfaces1()
func (c *TypeCheckTest) TestInheritPropFromMultipleInterfaces1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
function High1() {}
/** @type {number|string} */
High1.prototype.prop;
/** @interface */
function High2() {}
/** @type {number} */
High2.prototype.prop;
/**
 * @interface
 * @extends {High1}
 * @extends {High2}
 */
function Low() {}
function f(/** !Low */ x) { var /** null */ n = x.prop; }
`).
		AddDiagnostic(`initializing variable
found   : (number|string)
required: null
`).
		Run()
}

// Java: TypeCheckTest.testInheritPropFromMultipleInterfaces2()
func (c *TypeCheckTest) TestInheritPropFromMultipleInterfaces2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
function High1() {}
/** @type {number} */
High1.prototype.prop;
/** @interface */
function High2() {}
/** @type {number|string} */
High2.prototype.prop;
/**
 * @interface
 * @extends {High1}
 * @extends {High2}
 */
function Low() {}
function f(/** !Low */ x) { var /** null */ n = x.prop; }
`).
		AddDiagnostic(`initializing variable
found   : number
required: null
`).
		Run()
}

// Java: TypeCheckTest.testInheritPropFromMultipleInterfaces3()
func (c *TypeCheckTest) TestInheritPropFromMultipleInterfaces3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @interface
 * @template T1
 */
function MyCollection() {}
/**
 * @interface
 * @template T2
 * @extends {MyCollection<T2>}
 */
function MySet() {}
/**
 * @interface
 * @template T3,T4
 */
function MyMapEntry() {}
/**
 * @interface
 * @template T5,T6
 */
function MyMultimap() {}
/** @return {MyCollection<MyMapEntry<T5, T6>>} */
MyMultimap.prototype.entries = function() {};
/**
 * @interface
 * @template T7,T8
 * @extends {MyMultimap<T7, T8>}
 */
function MySetMultimap() {}
/** @return {MySet<MyMapEntry<T7, T8>>} */
MySetMultimap.prototype.entries = function() {};
/**
 * @interface
 * @template T9,T10
 * @extends {MyMultimap<T9, T10>}
 */
function MyFilteredMultimap() {}
/**
 * @interface
 * @template T11,T12
 * @extends {MyFilteredMultimap<T11, T12>}
 * @extends {MySetMultimap<T11, T12>}
 */
function MyFilteredSetMultimap() {}
`).
		Run()
}

// Java: TypeCheckTest.testInheritSameGenericInterfaceFromDifferentPaths()
func (c *TypeCheckTest) TestInheritSameGenericInterfaceFromDifferentPaths() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const */ var ns = {};
/**
 * @constructor
 * @template T1
 */
ns.Foo = function() {};
/**
 * @interface
 * @template T2
 */
ns.High = function() {};
/** @type {!ns.Foo<T2>} */
ns.High.prototype.myprop;
/**
 * @interface
 * @template T3
 * @extends {ns.High<T3>}
 */
ns.Med1 = function() {};
/**
 * @interface
 * @template T4
 * @extends {ns.High<T4>}
 */
ns.Med2 = function() {};
/**
 * @interface
 * @template T5
 * @extends {ns.Med1<T5>}
 * @extends {ns.Med2<T5>}
 */
ns.Low = function() {};
`).
		Run()
}

// Java: TypeCheckTest.testConversionFromInterfaceToRecursiveConstructor()
func (c *TypeCheckTest) TestConversionFromInterfaceToRecursiveConstructor() {
	c.TestClosureTypesMultipleWarnings(c.SuppressMissingProperty("foo")+`/** @interface */ var OtherType = function() {}
/** @implements {MyType}
 * @constructor */
var MyType = function() {}
/** @type {MyType} */
var x = /** @type {!OtherType} */ (new Object());
`, []string{"Cycle detected in inheritance chain of type MyType", `initializing variable
found   : OtherType
required: (MyType|null)`})
}

// Java: TypeCheckTest.testDirectPrototypeAssign()
func (c *TypeCheckTest) TestDirectPrototypeAssign() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
/** @constructor */ function Bar() {}
/** @type {Array} */ Bar.prototype = new Foo()
`).
		Run()
}

// Java: TypeCheckTest.testResolutionViaRegistry1()
func (c *TypeCheckTest) TestResolutionViaRegistry1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ u.T = function() {};
/** @type {(number|string)} */ u.T.prototype.a;
/**
* @param {u.T} t
* @return {string}
*/
var f = function(t) { return t.a; };
`).
		AddDiagnostic(`inconsistent return type
found   : (number|string)
required: string
`).
		Run()
}

// Java: TypeCheckTest.testResolutionViaRegistry2()
func (c *TypeCheckTest) TestResolutionViaRegistry2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ u.T = function() {
  this.a = 0; };
/**
* @param {u.T} t
* @return {string}
*/
var f = function(t) { return t.a; };
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testResolutionViaRegistry3()
func (c *TypeCheckTest) TestResolutionViaRegistry3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ u.T = function() {};
/** @type {(number|string)} */ u.T.prototype.a = 0;
/**
* @param {u.T} t
* @return {string}
*/
var f = function(t) { return t.a; };
`).
		AddDiagnostic(`inconsistent return type
found   : (number|string)
required: string
`).
		Run()
}

// Java: TypeCheckTest.testResolutionViaRegistry4()
func (c *TypeCheckTest) TestResolutionViaRegistry4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ u.A = function() {};
/**
* @constructor
* @extends {u.A}
*/
u.A.A = function() {};
/**
* @constructor
* @extends {u.A}
*/
u.A.B = function() {};
var ab = new u.A.B();
/** @type {!u.A} */ var a = ab;
/** @type {!u.A.A} */ var aa = ab;
`).
		AddDiagnostic(`initializing variable
found   : u.A.B
required: u.A.A
`).
		Run()
}

// Java: TypeCheckTest.testResolutionViaRegistry5()
func (c *TypeCheckTest) TestResolutionViaRegistry5() {
	var n *rhino.Node = c.ParseAndTypeCheck("/** @constructor */ u.T = function() {}; u.T")
	var type_ rhino.JSType = n.GetLastChild().GetLastChild().
		GetJSType()
	AssertThat(type_.IsUnknownType()).IsFalse()
	IsInstanceOf[rhino.FunctionType](AssertThat(type_))
	AssertThat((type_.(rhino.FunctionType)).GetInstanceType().GetReferenceName()).IsEqualTo("u.T")
}

// Java: TypeCheckTest.testGatherProperyWithoutAnnotation1()
func (c *TypeCheckTest) TestGatherProperyWithoutAnnotation1() {
	var n *rhino.Node = c.ParseAndTypeCheck(`/** @constructor */ var T = function() {};
/** @type {!T} */var t; t.x; t;
`)
	var type_ rhino.JSType = n.GetLastChild().GetLastChild().
		GetJSType()
	AssertThat(type_.IsUnknownType()).IsFalse()
	IsInstanceOf[rhino.ObjectType](AssertThat(type_))
	var objectType rhino.ObjectType = type_.(rhino.ObjectType)
	AssertThat(objectType.HasProperty("x")).IsFalse()
}

// Java: TypeCheckTest.testGatherProperyWithoutAnnotation2()
func (c *TypeCheckTest) TestGatherProperyWithoutAnnotation2() {
	var ns *TypeCheckTestCase_TypeCheckResult = c.ParseAndTypeCheckWithScope("/** @type {!Object} */var t; t.x; t;")
	var n *rhino.Node = ns.root
	var type_ rhino.JSType = n.GetLastChild().GetLastChild().
		GetJSType()
	AssertThat(type_.IsUnknownType()).IsFalse()
	c.AssertTypeEquals(type_, c.GetNativeObjectType())
	IsInstanceOf[rhino.ObjectType](AssertThat(type_))
	var objectType rhino.ObjectType = type_.(rhino.ObjectType)
	AssertThat(objectType.HasProperty("x")).IsFalse()
}

// Java: TypeCheckTest.testFunctionMasksVariableBug()
func (c *TypeCheckTest) TestFunctionMasksVariableBug() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var x = 4; var f = function x(b) { return b ? 1 : x(true); };").
		AddDiagnostic("function x masks variable (IE bug)").
		Run()
}

// Java: TypeCheckTest.testDfa1()
func (c *TypeCheckTest) TestDfa1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var x = null;
 x = 1;
 /** @type {number} */ var y = x;`).
		Run()
}

// Java: TypeCheckTest.testDfa2()
func (c *TypeCheckTest) TestDfa2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function u() {}
/** @return {number} */ function f() {
var x = 'todo';
if (u()) { x = 1; } else { x = 2; } return x;
}
`).
		Run()
}

// Java: TypeCheckTest.testDfa3()
func (c *TypeCheckTest) TestDfa3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function u() {}
/** @return {number} */ function f() {
/** @type {number|string} */ var x = 'todo';
if (u()) { x = 1; } else { x = 2; } return x;
}
`).
		Run()
}

// Java: TypeCheckTest.testDfa4()
func (c *TypeCheckTest) TestDfa4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {Date?} d */ function f(d) {
if (!d) { return; }
/** @type {!Date} */ var e = d;
}
`).
		Run()
}

// Java: TypeCheckTest.testDfa5()
func (c *TypeCheckTest) TestDfa5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {string?} */ function u() {return 'a';}
/** @param {string?} x
@return {string} */ function f(x) {
while (!x) { x = u(); }
return x;
}
`).
		Run()
}

// Java: TypeCheckTest.testDfa6()
func (c *TypeCheckTest) TestDfa6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {Object?} */ function u() {return {};}
/** @param {Object?} x */ function f(x) {
while (x) { x = u(); if (!x) { x = u(); } }
}
`).
		Run()
}

// Java: TypeCheckTest.testDfa7()
func (c *TypeCheckTest) TestDfa7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ var T = function() {};
/** @type {Date?} */ T.prototype.x = null;
/** @param {!T} t */ function f(t) {
if (!t.x) { return; }
/** @type {!Date} */ var e = t.x;
}
`).
		Run()
}

// Java: TypeCheckTest.testDfa8()
func (c *TypeCheckTest) TestDfa8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ var T = function() {};
/** @type {number|string} */ T.prototype.x = '';
function u() {}
/** @param {!T} t
@return {number} */ function f(t) {
if (u()) { t.x = 1; } else { t.x = 2; } return t.x;
}
`).
		Run()
}

// Java: TypeCheckTest.testDfa9()
func (c *TypeCheckTest) TestDfa9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f() {
/** @type {string?} */var x;
x = null;
if (x == null) { return 0; } else { return 1; } }
`).
		AddDiagnostic(`condition always evaluates to true
left : null
right: null
`).
		Run()
}

// Java: TypeCheckTest.testDfa10()
func (c *TypeCheckTest) TestDfa10() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {null} x */ function g(x) {}
/** @param {string?} x */function f(x) {
if (!x) { x = ''; }
if (g(x)) { return 0; } else { return 1; } }
`).
		AddDiagnostic(`actual parameter 1 of g does not match formal parameter
found   : string
required: null
`).
		Run()
}

// Java: TypeCheckTest.testDfa11()
func (c *TypeCheckTest) TestDfa11() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {string} opt_x
@return {string} */
function f(opt_x) { if (!opt_x) {
throw new Error('x cannot be empty'); } return opt_x; }
`).
		Run()
}

// Java: TypeCheckTest.testDfa12()
func (c *TypeCheckTest) TestDfa12() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {string} x
 * @constructor
 */
var Bar = function(x) {};
/** @param {string} x */ function g(x) { return true; }
/** @param {string|number} opt_x */
function f(opt_x) {
  if (opt_x) { new Bar(g(opt_x) && 'x'); }
}
`).
		AddDiagnostic(`actual parameter 1 of g does not match formal parameter
found   : (number|string)
required: string
`).
		Run()
}

// Java: TypeCheckTest.testDfa13()
func (c *TypeCheckTest) TestDfa13() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {string} x
 * @param {number} y
 * @param {number} z
 */
function g(x, y, z) {}
function f() {
  var x = 'a'; g(x, x = 3, x);
}
`).
		Run()
}

// Java: TypeCheckTest.testTypeInferenceWithCast1()
func (c *TypeCheckTest) TestTypeInferenceWithCast1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**@return {(number|null|undefined)}*/function u(x) {return null;}
/**@param {number?} x
@return {number?}*/function f(x) {return x;}
/**@return {number?}*/function g(x) {
var y = /**@type {number?}*/(u(x)); return f(y);}
`).
		Run()
}

// Java: TypeCheckTest.testTypeInferenceWithCast2()
func (c *TypeCheckTest) TestTypeInferenceWithCast2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**@return {(number|null|undefined)}*/function u(x) {return null;}
/**@param {number?} x
@return {number?}*/function f(x) {return x;}
/**@return {number?}*/function g(x) {
var y; y = /**@type {number?}*/(u(x)); return f(y);}
`).
		Run()
}

// Java: TypeCheckTest.testTypeInferenceWithCast3()
func (c *TypeCheckTest) TestTypeInferenceWithCast3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**@return {(number|null|undefined)}*/function u(x) {return 1;}
/**@return {number}*/function g(x) {
return /**@type {number}*/(u(x));}
`).
		Run()
}

// Java: TypeCheckTest.testTypeInferenceWithCast4()
func (c *TypeCheckTest) TestTypeInferenceWithCast4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**@return {(number|null|undefined)}*/function u(x) {return 1;}
/**@return {number}*/function g(x) {
return /**@type {number}*/(u(x)) && 1;}
`).
		Run()
}

// Java: TypeCheckTest.testTypeInferenceWithCast5()
func (c *TypeCheckTest) TestTypeInferenceWithCast5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {number} x */ function foo(x) {}
/** @param {{length:*}} y */ function bar(y) {
  /** @type {string} */ y.length;
  foo(y.length);
}
`).
		AddDiagnostic(`actual parameter 1 of foo does not match formal parameter
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testTypeInferenceWithClosure1()
func (c *TypeCheckTest) TestTypeInferenceWithClosure1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {boolean} */
function f() {
  /** @type {?string} */ var x = null;
  function g() { x = 'y'; } g();
  return x == null;
}
`).
		Run()
}

// Java: TypeCheckTest.testTypeInferenceWithClosure2()
func (c *TypeCheckTest) TestTypeInferenceWithClosure2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {boolean} */
function f() {
  /** @type {?string} */ var x = null;
  function g() { x = 'y'; } g();
  return x === 3;
}
`).
		AddDiagnostic(`condition always evaluates to false
left : (null|string)
right: number
`).
		Run()
}

// Java: TypeCheckTest.testTypeInferenceWithNoEntry1()
func (c *TypeCheckTest) TestTypeInferenceWithNoEntry1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {number} x */ function f(x) {}
/** @constructor */ function Foo() {}
Foo.prototype.init = function() {
  /** @type {?{baz: number}} */ this.bar = {baz: 3};
};
/**
 * @extends {Foo}
 * @constructor
 */
function SubFoo() {}
/** Method */
SubFoo.prototype.method = function() {
  for (var i = 0; i < 10; i++) {
    f(this.bar);
    f(this.bar.baz);
  }
};
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : (null|{baz: number})
required: number
`).
		Run()
}

// Java: TypeCheckTest.testTypeInferenceWithNoEntry2()
func (c *TypeCheckTest) TestTypeInferenceWithNoEntry2() {
	c.TestClosureTypes(`/** @param {number} x */ function f(x) {}
/** @param {!Object} x */ function g(x) {}
/** @constructor */ function Foo() {}
Foo.prototype.init = function() {
  /** @type {?{baz: number}} */ this.bar = {baz: 3};
};
/**
 * @extends {Foo}
 * @constructor
 */
function SubFoo() {}
/** Method */
SubFoo.prototype.method = function() {
  for (var i = 0; i < 10; i++) {
    f(this.bar);
    goog.asserts.assert(this.bar);
    g(this.bar);
  }
};
`, jutil.StringPtr(`actual parameter 1 of f does not match formal parameter
found   : (null|{baz: number})
required: number
`))
}

// Java: TypeCheckTest.testForwardPropertyReference()
func (c *TypeCheckTest) TestForwardPropertyReference() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ var Foo = function() { this.init(); };
/** @return {string} */
Foo.prototype.getString = function() {
  return this.number_;
};
Foo.prototype.init = function() {
  /** @type {number} */
  this.number_ = 3;
};
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testNoForwardTypeDeclaration()
func (c *TypeCheckTest) TestNoForwardTypeDeclaration() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @param {MyType} x */ function f(x) {}").
		AddDiagnostic("Bad type annotation. Unknown type MyType").
		Run()
}

// Java: TypeCheckTest.testNoForwardTypeDeclarationAndNoBraces()
func (c *TypeCheckTest) TestNoForwardTypeDeclarationAndNoBraces() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @return The result. */ function f() {}").
		AddDiagnosticType(RhinoErrorReporter_JSDOC_MISSING_TYPE_WARNING).
		Run()
}

// Java: TypeCheckTest.testForwardTypeDeclaration1()
func (c *TypeCheckTest) TestForwardTypeDeclaration1() {
	f := `goog.forwardDeclare('MyType');
/** @param {!MyType} x */ function f(x) { }
`
	c.TestClosureTypes(f, nil)
	c.TestClosureTypes(f+"f(3);", jutil.StringPtr(`actual parameter 1 of f does not match formal parameter
found   : number
required: MyType
`))
}

// Java: TypeCheckTest.testForwardTypeDeclaration2()
func (c *TypeCheckTest) TestForwardTypeDeclaration2() {
	f := `goog.forwardDeclare('MyType');
/** @param {MyType} x */ function f(x) { }
`
	c.TestClosureTypes(f, nil)
	c.TestClosureTypes(f+"f(3);", jutil.StringPtr(`actual parameter 1 of f does not match formal parameter
found   : number
required: (MyType|null)
`))
}

// Java: TypeCheckTest.testForwardTypeDeclaration3()
func (c *TypeCheckTest) TestForwardTypeDeclaration3() {
	c.TestClosureTypes(`goog.forwardDeclare('MyType');
/** @param {MyType} x */ function f(x) { return x; }
/** @constructor */ var MyType = function() {};
f(3);
`, jutil.StringPtr(`actual parameter 1 of f does not match formal parameter
found   : number
required: (MyType|null)
`))
}

// Java: TypeCheckTest.testForwardTypeDeclaration4()
func (c *TypeCheckTest) TestForwardTypeDeclaration4() {
	c.TestClosureTypes(`goog.forwardDeclare('MyType');
/** @param {MyType} x */ function f(x) { return x; }
/** @constructor */ var MyType = function() {};
f(new MyType());
`, nil)
}

// Java: TypeCheckTest.testForwardTypeDeclaration5()
func (c *TypeCheckTest) TestForwardTypeDeclaration5() {
	c.TestClosureTypes(`goog.forwardDeclare('MyType');
/**
 * @constructor
 * @extends {MyType}
 */ var YourType = function() {};
/** @override */ YourType.prototype.method = function() {};
`, jutil.StringPtr("Could not resolve type in @extends tag of YourType"))
}

// Java: TypeCheckTest.testForwardTypeDeclaration6()
func (c *TypeCheckTest) TestForwardTypeDeclaration6() {
	c.TestClosureTypesMultipleWarnings(`goog.forwardDeclare('MyType');
/**
 * @constructor
 * @implements {MyType}
 */ var YourType = function() {};
/** @override */ YourType.prototype.method = function() {};
`, []string{"Could not resolve type in @implements tag of YourType", "property method not defined on any superclass of YourType"})
}

// Java: TypeCheckTest.testForwardTypeDeclaration7()
func (c *TypeCheckTest) TestForwardTypeDeclaration7() {
	c.TestClosureTypes(`goog.forwardDeclare('MyType');
/** @param {MyType=} x */
function f(x) { return x == undefined; }
`, nil)
}

// Java: TypeCheckTest.testForwardTypeDeclaration8()
func (c *TypeCheckTest) TestForwardTypeDeclaration8() {
	c.TestClosureTypes(`goog.forwardDeclare('MyType');
/** @param {MyType} x */
function f(x) { return x.name == undefined; }
`, nil)
}

// Java: TypeCheckTest.testForwardTypeDeclaration9()
func (c *TypeCheckTest) TestForwardTypeDeclaration9() {
	c.TestClosureTypes(`goog.forwardDeclare('MyType');
/** @param {MyType} x */
function f(x) { x.name = 'Bob'; }
`, nil)
}

// Java: TypeCheckTest.testForwardTypeDeclaration10()
func (c *TypeCheckTest) TestForwardTypeDeclaration10() {
	f := `goog.forwardDeclare('MyType');
/** @param {MyType|number} x */ function f(x) { }
`
	c.TestClosureTypes(f, nil)
	c.TestClosureTypes(f+"f(3);", nil)
	c.TestClosureTypes(f+"f('3');", jutil.StringPtr(`actual parameter 1 of f does not match formal parameter
found   : string
required: (MyType|null|number)
`))
}

// Java: TypeCheckTest.testForwardTypeDeclaration12()
func (c *TypeCheckTest) TestForwardTypeDeclaration12() {
	c.TestClosureTypes(`goog.forwardDeclare('MyType');
/**
 * @param {!Function} ctor
 * @return {MyType}
 */
function f(ctor) { return new ctor(); }
`, nil)
}

// Java: TypeCheckTest.testForwardTypeDeclaration13()
func (c *TypeCheckTest) TestForwardTypeDeclaration13() {
	c.TestClosureTypes(`goog.forwardDeclare('MyType');
/**
 * @param {!Function} ctor
 * @return {MyType}
 */
function f(ctor) { return (new ctor()).impossibleProp; }
`, jutil.StringPtr("Property impossibleProp never defined on ?"+TypeCheck_POSSIBLE_INEXISTENT_PROPERTY_EXPLANATION))
}

// Java: TypeCheckTest.testDuplicateTypeDef()
func (c *TypeCheckTest) TestDuplicateTypeDef() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = {};
/** @constructor */ goog.Bar = function() {};
/** @typedef {number} */ goog.Bar;
`).
		AddDiagnostic(`variable goog.Bar redefined with type None, original definition at testcode0:2 with type (typeof goog.Bar)
`).
		Run()
}

// Java: TypeCheckTest.testTypeDef1()
func (c *TypeCheckTest) TestTypeDef1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = {};
/** @typedef {number} */ goog.Bar;
/** @param {goog.Bar} x */ function f(x) {}
f(3);
`).
		Run()
}

// Java: TypeCheckTest.testTypeDef2()
func (c *TypeCheckTest) TestTypeDef2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = {};
/** @typedef {number} */ goog.Bar;
/** @param {goog.Bar} x */ function f(x) {}
f('3');
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testTypeDef3()
func (c *TypeCheckTest) TestTypeDef3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = {};
/** @typedef {number} */ var Bar;
/** @param {Bar} x */ function f(x) {}
f('3');
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testTypeDef4()
func (c *TypeCheckTest) TestTypeDef4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function A() {}
/** @constructor */ function B() {}
/** @typedef {(A|B)} */ var AB;
/** @param {AB} x */ function f(x) {}
f(new A()); f(new B()); f(1);
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : number
required: (A|B|null)
`).
		Run()
}

// Java: TypeCheckTest.testTypeDef5()
func (c *TypeCheckTest) TestTypeDef5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {AB} x */ function f(x) {}
/** @constructor */ function A() {}
/** @constructor */ function B() {}
/** @typedef {(A|B)} */ var AB;
f(new A()); f(new B()); f(1);
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : number
required: (A|B|null)
`).
		Run()
}

// Java: TypeCheckTest.testCircularTypeDef()
func (c *TypeCheckTest) TestCircularTypeDef() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var goog = {};
/** @typedef {number|Array<goog.Bar>} */ goog.Bar;
/** @param {goog.Bar} x */ function f(x) {}
f(3); f([3]); f([[3]]);
`).
		Run()
}

// Java: TypeCheckTest.testPrototypePropertyReference()
func (c *TypeCheckTest) TestPrototypePropertyReference() {
	var p *TypeCheckTestCase_TypeCheckResult = c.ParseAndTypeCheckWithScopeWithExterns(CompilerTypeTestCase_DEFAULT_EXTERNS, `/** @constructor */
function Foo() {}
/** @param {number} a */
Foo.prototype.bar = function(a){};
/** @param {Foo} f */
function baz(f) {
  Foo.prototype.bar.call(f, 3);
}
`)
	AssertThat(c.compiler.GetErrorCount()).IsEqualTo(0)
	AssertThat(c.compiler.GetWarningCount()).IsEqualTo(0)
	IsInstanceOf[rhino.FunctionType](AssertThat(p.scope.GetVar("Foo").GetType()))
	var fooType rhino.FunctionType = p.scope.GetVar("Foo").GetType().(rhino.FunctionType)
	AssertThat(fooType.GetPrototype().GetPropertyType("bar").
		String()).IsEqualTo("function(this:Foo, number): undefined")
}

// Java: TypeCheckTest.testMissingProperty1a()
func (c *TypeCheckTest) TestMissingProperty1a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
Foo.prototype.bar = function() { return this.a; };
Foo.prototype.baz = function() { this.a = 3; };
`).
		AddDiagnostic("Property a never defined on Foo").
		AddDiagnostic("Property a never defined on Foo").
		Run()
}

// Java: TypeCheckTest.testMissingProperty1b()
func (c *TypeCheckTest) TestMissingProperty1b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @constructor */ function Foo() {}
Foo.prototype.bar = function() { return this.a; };
Foo.prototype.baz = function() { this.a = 3; };
`).
		Run()
}

// Java: TypeCheckTest.testMissingProperty2a()
func (c *TypeCheckTest) TestMissingProperty2a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @constructor */ function Foo() {}
Foo.prototype.bar = function() { return this.a; };
Foo.prototype.baz = function() { this.b = 3; };
`).
		AddDiagnostic("Property a never defined on Foo").
		Run()
}

// Java: TypeCheckTest.testMissingProperty2b()
func (c *TypeCheckTest) TestMissingProperty2b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
Foo.prototype.baz = function() { this.b = 3; };
`).
		AddDiagnostic("Property b never defined on Foo").
		Run()
}

// Java: TypeCheckTest.testMissingProperty3a()
func (c *TypeCheckTest) TestMissingProperty3a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
Foo.prototype.bar = function() { return this.a; };
(new Foo).a = 3;
`).
		AddDiagnostic("Property a never defined on Foo").
		AddDiagnostic("Property a never defined on Foo").
		Run()
}

// Java: TypeCheckTest.testMissingProperty3b()
func (c *TypeCheckTest) TestMissingProperty3b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @constructor */ function Foo() {}
Foo.prototype.bar = function() { return this.a; };
(new Foo).a = 3;
`).
		Run()
}

// Java: TypeCheckTest.testMissingProperty4a()
func (c *TypeCheckTest) TestMissingProperty4a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
Foo.prototype.bar = function() { return this.a; };
(new Foo).b = 3;
`).
		AddDiagnostic("Property a never defined on Foo").
		Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		Run()
}

// Java: TypeCheckTest.testMissingProperty4b()
func (c *TypeCheckTest) TestMissingProperty4b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
(new Foo).b = 3;
`).
		AddDiagnostic("Property b never defined on Foo").
		Run()
}

// Java: TypeCheckTest.testMissingProperty5()
func (c *TypeCheckTest) TestMissingProperty5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
Foo.prototype.bar = function() { return this.a; };
/** @constructor */ function Bar() { this.a = 3; };
`).
		AddDiagnostic("Property a never defined on Foo").
		Run()
}

// Java: TypeCheckTest.testMissingProperty6a()
func (c *TypeCheckTest) TestMissingProperty6a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
Foo.prototype.bar = function() { return this.a; };
/** @constructor
 * @extends {Foo} */
function Bar() { this.a = 3; };
`).
		AddDiagnostic("Property a never defined on Foo").
		Run()
}

// Java: TypeCheckTest.testMissingProperty6b()
func (c *TypeCheckTest) TestMissingProperty6b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @constructor */ function Foo() {}
Foo.prototype.bar = function() { return this.a; };
/** @constructor
 * @extends {Foo} */
function Bar() { this.a = 3; };
`).
		Run()
}

// Java: TypeCheckTest.testMissingProperty7()
func (c *TypeCheckTest) TestMissingProperty7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {Object} obj */
function foo(obj) { return obj.impossible; }
`).
		AddDiagnostic("Property impossible never defined on Object" + TypeCheck_POSSIBLE_INEXISTENT_PROPERTY_EXPLANATION).
		Run()
}

// Java: TypeCheckTest.testMissingProperty8()
func (c *TypeCheckTest) TestMissingProperty8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {Object} obj */
function foo(obj) { return typeof obj.impossible; }
`).
		Run()
}

// Java: TypeCheckTest.testMissingProperty9()
func (c *TypeCheckTest) TestMissingProperty9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @param {Object} obj */
function foo(obj) { if (obj.impossible) { return true; } }
`).
		Run()
}

// Java: TypeCheckTest.testMissingProperty10()
func (c *TypeCheckTest) TestMissingProperty10() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @param {Object} obj */
function foo(obj) { while (obj.impossible) { return true; } }
`).
		Run()
}

// Java: TypeCheckTest.testMissingProperty11()
func (c *TypeCheckTest) TestMissingProperty11() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @param {Object} obj */
function foo(obj) { for (;obj.impossible;) { return true; } }
`).
		Run()
}

// Java: TypeCheckTest.testMissingProperty12()
func (c *TypeCheckTest) TestMissingProperty12() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @param {Object} obj */
function foo(obj) { do { } while (obj.impossible); }
`).
		Run()
}

// Java: TypeCheckTest.testMissingProperty15()
func (c *TypeCheckTest) TestMissingProperty15() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @param {Object} x */
function f(x) { if (x.foo) { x.foo(); } }
`).
		Run()
}

// Java: TypeCheckTest.testMissingProperty16()
func (c *TypeCheckTest) TestMissingProperty16() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @param {Object} x */
function f(x) { x.foo(); if (x.foo) {} }
`).
		AddDiagnostic("Property foo never defined on Object" + TypeCheck_POSSIBLE_INEXISTENT_PROPERTY_EXPLANATION).
		Run()
}

// Java: TypeCheckTest.testMissingProperty17()
func (c *TypeCheckTest) TestMissingProperty17() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {Object} x */
function f(x) { if (typeof x.foo == 'function') { x.foo(); } }
`).
		Run()
}

// Java: TypeCheckTest.testMissingProperty18()
func (c *TypeCheckTest) TestMissingProperty18() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @param {Object} x */
function f(x) { if (x.foo instanceof Function) { x.foo(); } }
`).
		Run()
}

// Java: TypeCheckTest.testMissingProperty19()
func (c *TypeCheckTest) TestMissingProperty19() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @param {Object} x */
function f(x) { if (x.bar) { if (x.foo) {} } else { x.foo(); } }
`).
		AddDiagnostic("Property foo never defined on Object" + TypeCheck_POSSIBLE_INEXISTENT_PROPERTY_EXPLANATION).
		Run()
}

// Java: TypeCheckTest.testMissingProperty20()
func (c *TypeCheckTest) TestMissingProperty20() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @param {Object} x */
function f(x) { if (x.foo) { } else { x.foo(); } }
`).
		AddDiagnostic("Property foo never defined on Object" + TypeCheck_POSSIBLE_INEXISTENT_PROPERTY_EXPLANATION).
		Run()
}

// Java: TypeCheckTest.testMissingProperty21()
func (c *TypeCheckTest) TestMissingProperty21() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @param {Object} x */
function f(x) { x.foo && x.foo(); }
`).
		Run()
}

// Java: TypeCheckTest.testMissingProperty22()
func (c *TypeCheckTest) TestMissingProperty22() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @param {Object} x
 * @return {boolean} */
function f(x) { return x.foo ? x.foo() : true; }
`).
		Run()
}

// Java: TypeCheckTest.testMissingProperty23()
func (c *TypeCheckTest) TestMissingProperty23() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function f(x) { x.impossible(); }").
		AddDiagnostic("Property impossible never defined on x" + TypeCheck_POSSIBLE_INEXISTENT_PROPERTY_EXPLANATION).
		Run()
}

// Java: TypeCheckTest.testMissingProperty24()
func (c *TypeCheckTest) TestMissingProperty24() {
	c.TestClosureTypes(`goog.forwardDeclare('MissingType');
/** @param {MissingType} x */
function f(x) { x.impossible(); }
`, nil)
}

// Java: TypeCheckTest.testMissingProperty25()
func (c *TypeCheckTest) TestMissingProperty25() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ var Foo = function() {};
Foo.prototype.bar = function() {};
/** @constructor */ var FooAlias = Foo;
(new FooAlias()).bar();
`).
		Run()
}

// Java: TypeCheckTest.testMissingProperty26()
func (c *TypeCheckTest) TestMissingProperty26() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ var Foo = function() {};
/** @constructor */ var FooAlias = Foo;
FooAlias.prototype.bar = function() {};
(new Foo()).bar();
`).
		Run()
}

// Java: TypeCheckTest.testMissingProperty27()
func (c *TypeCheckTest) TestMissingProperty27() {
	c.TestClosureTypes(`goog.forwardDeclare('MissingType');
/** @param {?MissingType} x */
function f(x) {
  for (var parent = x; parent; parent = parent.getParent()) {}
}
`, nil)
}

// Java: TypeCheckTest.testMissingProperty28a()
func (c *TypeCheckTest) TestMissingProperty28a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(obj) {
  /** @type {*} */ obj.foo;
  return obj.foo;
}
`).
		Run()
}

// Java: TypeCheckTest.testMissingProperty28b()
func (c *TypeCheckTest) TestMissingProperty28b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(obj) {
  /** @type {*} */ obj.foo;
  return obj.foox;
}
`).
		AddDiagnostic("Property foox never defined on obj" + TypeCheck_POSSIBLE_INEXISTENT_PROPERTY_EXPLANATION).
		Run()
}

// Java: TypeCheckTest.testMissingProperty29()
func (c *TypeCheckTest) TestMissingProperty29() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @constructor */ var Foo;
Foo.prototype.opera;
Foo.prototype.opera.postError;
`).
		AddSource("").
		Run()
}

// Java: TypeCheckTest.testMissingProperty30a()
func (c *TypeCheckTest) TestMissingProperty30a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {*} */
function f() {
 return {};
}
f().a = 3;
/** @param {Object} y */ function g(y) { return y.a; }
`).
		AddDiagnostic("Property a never defined on *").
		AddDiagnostic("Property a never defined on Object").
		Run()
}

// Java: TypeCheckTest.testMissingProperty30b()
func (c *TypeCheckTest) TestMissingProperty30b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @return {*} */
function f() {
 return {};
}
f().a = 3;
/** @param {Object} y */ function g(y) { return y.a; }
`).
		Run()
}

// Java: TypeCheckTest.testMissingProperty31a()
func (c *TypeCheckTest) TestMissingProperty31a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {Array|number} */
function f() {
 return [];
}
f().a = 3;
`).
		AddDiagnostic("Property a never defined on (Array|Number)").
		Run()
}

// Java: TypeCheckTest.testMissingProperty31b()
func (c *TypeCheckTest) TestMissingProperty31b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @return {Array|number} */
function f() {
 return [];
}
f().a = 3;
/** @param {Array} y */ function g(y) { return y.a; }
`).
		Run()
}

// Java: TypeCheckTest.testMissingProperty32()
func (c *TypeCheckTest) TestMissingProperty32() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @return {Array|number} */
function f() {
 return [];
}
f().a = 3;
/** @param {Date} y */ function g(y) { return y.a; }
`).
		AddDiagnostic("Property a never defined on Date").
		Run()
}

// Java: TypeCheckTest.testMissingProperty33()
func (c *TypeCheckTest) TestMissingProperty33() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @param {Object} x */
function f(x) { !x.foo || x.foo(); }
`).
		Run()
}

// Java: TypeCheckTest.testMissingProperty34()
func (c *TypeCheckTest) TestMissingProperty34() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @fileoverview
 * @suppress {missingProperties} */
/** @constructor */ function Foo() {}
Foo.prototype.bar = function() { return this.a; };
Foo.prototype.baz = function() { this.b = 3; };
`).
		Run()
}

// Java: TypeCheckTest.testMissingProperty35a()
func (c *TypeCheckTest) TestMissingProperty35a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
/** @constructor */ function Bar() {}
/** @constructor */ function Baz() {}
/** @param {Foo|Bar} x */ function f(x) { x.specialProp = 1; }
/** @param {Bar|Baz} x */ function g(x) { return x.specialProp; }
`).
		AddDiagnostic("Property specialProp never defined on (Foo|Bar)").
		AddDiagnostic("Property specialProp never defined on (Bar|Baz)").
		Run()
}

// Java: TypeCheckTest.testMissingProperty35b()
func (c *TypeCheckTest) TestMissingProperty35b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
/** @constructor */ function Bar() {}
/** @constructor */ function Baz() {}
/** @param {Foo|Bar} x */ function f(x) { x.specialProp = 1; }
/** @param {Bar|Baz} x */ function g(x) { return x.specialProp; }
`).
		Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		Run()
}

// Java: TypeCheckTest.testMissingProperty36a()
func (c *TypeCheckTest) TestMissingProperty36a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
Foo.prototype.baz = 0;
/** @constructor
 * @extends {Foo} */ function SubFoo() {}
SubFoo.prototype.bar = 0;
/** @param {{bar: number}} x */ function f(x) { return x.baz; }
`).
		AddDiagnostic("Property baz never defined on x").
		Run()
}

// Java: TypeCheckTest.testMissingProperty36b()
func (c *TypeCheckTest) TestMissingProperty36b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
Foo.prototype.baz = 0;
/** @constructor
 * @extends {Foo} */ function SubFoo() {}
SubFoo.prototype.bar = 0;
/** @param {{bar: number}} x */ function f(x) { return x.baz; }
`).
		Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		Run()
}

// Java: TypeCheckTest.testMissingProperty37a()
func (c *TypeCheckTest) TestMissingProperty37a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {{isVisible: boolean}} x */
function f(x){
  x.isVisible = false;
}
/** @constructor */
function Foo() {}
/**
 * @constructor
 * @extends {Foo}
 */
function SubFoo() {}
/** @type {boolean} */
SubFoo.prototype.isVisible = true;
/**
 * @param {Foo} x
 * @return {boolean}
 */
function g(x) { return x.isVisible; }
`).
		AddDiagnostic("Property isVisible never defined on Foo").
		Run()
}

// Java: TypeCheckTest.testMissingProperty37b()
func (c *TypeCheckTest) TestMissingProperty37b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {{isVisible: boolean}} x */ function f(x){
  x.isVisible = false;
}
/** @constructor */ function Foo() {}
/**
 * @constructor
 * @extends {Foo}
 */ function SubFoo() {}
/** @type {boolean} */ SubFoo.prototype.isVisible = true;
/**
 * @param {Foo} x
 * @return {boolean}
 */
function g(x) { return x.isVisible; }
`).
		Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		Run()
}

// Java: TypeCheckTest.testMissingProperty38()
func (c *TypeCheckTest) TestMissingProperty38() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
/** @constructor */ function Bar() {}
/** @return {Foo|Bar} */ function f() { return new Foo(); }
f().missing;
`).
		AddDiagnostic("Property missing never defined on (Foo|Bar)").
		Run()
}

// Java: TypeCheckTest.testMissingProperty39a()
func (c *TypeCheckTest) TestMissingProperty39a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddExterns(NewTestExternsBuilder().AddString().Build()).
		AddSource(`/** @return {string|number} */ function f() { return 3; }
f().length;
`).
		Run()
}

// Java: TypeCheckTest.testMissingProperty39b()
func (c *TypeCheckTest) TestMissingProperty39b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddString().Build()).
		AddSource(`/** @return {string|number} */ function f() { return 3; }
f().length;
`).
		AddDiagnostic("Property length not defined on all member types of (String|Number)").
		Run()
}

// Java: TypeCheckTest.testMissingProperty40a()
func (c *TypeCheckTest) TestMissingProperty40a() {
	c.TestClosureTypes(`goog.forwardDeclare('MissingType');
/** @param {MissingType} x */
function f(x) { x.impossible(); }
`, nil)
}

// Java: TypeCheckTest.testMissingProperty40b()
func (c *TypeCheckTest) TestMissingProperty40b() {
	c.TestClosureTypes(`goog.forwardDeclare('MissingType');
/** @param {(Array|MissingType)} x */
function f(x) { x.impossible(); }
`, jutil.StringPtr("Property impossible not defined on all member types of x"))
}

// Java: TypeCheckTest.testMissingProperty41a()
func (c *TypeCheckTest) TestMissingProperty41a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {(Array|Date)} x */
function f(x) { if (x.impossible) x.impossible(); }
`).
		AddDiagnostic("Property impossible never defined on (Array|Date)").
		AddDiagnostic("Property impossible never defined on (Array|Date)").
		Run()
}

// Java: TypeCheckTest.testMissingProperty41b()
func (c *TypeCheckTest) TestMissingProperty41b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @param {(Array|Date)} x */
function f(x) { if (x.impossible) x.impossible(); }
`).
		Run()
}

// Java: TypeCheckTest.testMissingProperty42()
func (c *TypeCheckTest) TestMissingProperty42() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {Object} x */
function f(x) {
  if (typeof x.impossible == 'undefined') throw Error();
  return x.impossible;
}
`).
		Run()
}

// Java: TypeCheckTest.testMissingProperty43()
func (c *TypeCheckTest) TestMissingProperty43() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`function f(x) {
 return /** @type {number} */ (x.impossible) && 1;
}
`).
		Run()
}

// Java: TypeCheckTest.testMissingProperty_notReportedInPropertyAbsenceCheck()
func (c *TypeCheckTest) TestMissingProperty_notReportedInPropertyAbsenceCheck() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`function f(/** !Object */ x) {
  if (x.y == null) throw new Error();
}
`).
		Run()
}

// Java: TypeCheckTest.optChainGetPropAllowLoosePropertyAccess()
func (c *TypeCheckTest) Test_optChainGetPropAllowLoosePropertyAccess() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @param {?Object} x */
function f(x) {
  x.y?.z;
}
`).
		Run()
}

// Java: TypeCheckTest.normalGetPropNotAllowLoosePropertyAccess()
func (c *TypeCheckTest) Test_normalGetPropNotAllowLoosePropertyAccess() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @param {?Object} x */
function f(x) {
  x.y.z;
}
`).
		AddDiagnostic("Property y never defined on Object" + TypeCheck_POSSIBLE_INEXISTENT_PROPERTY_EXPLANATION).
		Run()
}

// Java: TypeCheckTest.optChainGetElemAllowLoosePropertyAccess()
func (c *TypeCheckTest) Test_optChainGetElemAllowLoosePropertyAccess() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @param {?Object} x */
function f(x) {
  x.y?.[z];
}
`).
		Run()
}

// Java: TypeCheckTest.normalGetElemNotAllowLoosePropertyAccess()
func (c *TypeCheckTest) Test_normalGetElemNotAllowLoosePropertyAccess() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @param {?Object} x */
function f(x) {
  x.y[z];
}
`).
		AddDiagnostic("Property y never defined on Object" + TypeCheck_POSSIBLE_INEXISTENT_PROPERTY_EXPLANATION).
		Run()
}

// Java: TypeCheckTest.optChainCallAllowLoosePropertyAccess()
func (c *TypeCheckTest) Test_optChainCallAllowLoosePropertyAccess() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @param {?Object} x */
function f(x) {
  x.y?.();
}
`).
		Run()
}

// Java: TypeCheckTest.normalCallNotAllowLoosePropertyAccess()
func (c *TypeCheckTest) Test_normalCallNotAllowLoosePropertyAccess() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @param {?Object} x */
function f(x) {
  x.y();
}
`).
		AddDiagnostic("Property y never defined on Object" + TypeCheck_POSSIBLE_INEXISTENT_PROPERTY_EXPLANATION).
		Run()
}

// Java: TypeCheckTest.getNotFirstChildOfOptChainCallNotAllowLoosePropertyAccess()
func (c *TypeCheckTest) Test_getNotFirstChildOfOptChainCallNotAllowLoosePropertyAccess() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @param {?Object} x */
function f(x) {
  return false;
}
f?.(x.y)
`).
		AddDiagnostic("Property y never defined on x" + TypeCheck_POSSIBLE_INEXISTENT_PROPERTY_EXPLANATION).
		Run()
}

// Java: TypeCheckTest.getNotFirstChildOfOptionalGetElemNotAllowLoosePropertyAccess()
func (c *TypeCheckTest) Test_getNotFirstChildOfOptionalGetElemNotAllowLoosePropertyAccess() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @param {?Object} x */
function f(x) {
  x?.[y.z];
}
`).
		AddDiagnostic("Property z never defined on y" + TypeCheck_POSSIBLE_INEXISTENT_PROPERTY_EXPLANATION).
		Run()
}

// Java: TypeCheckTest.testOptChainGetPropProvidesThisForMethodCall()
func (c *TypeCheckTest) TestOptChainGetPropProvidesThisForMethodCall() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("class A {\n  foo() {}\n}\n/** @param {?A} a */\nfunction f(a) {\n// TypeCheck should not complain that foo() is getting called without a correctly typed\n// `this` value.\n  a?.foo();\n}\n").
		Run()
}

// Java: TypeCheckTest.testReflectObject1()
func (c *TypeCheckTest) TestReflectObject1() {
	c.TestClosureTypes(`goog.reflect = {};
goog.reflect.object = function(x, y){};
/** @constructor */ function A() {}
goog.reflect.object(A, {x: 3});
`, nil)
}

// Java: TypeCheckTest.testReflectObject2()
func (c *TypeCheckTest) TestReflectObject2() {
	c.TestClosureTypes(`goog.reflect = {};
goog.reflect.object = function(x, y){};
/** @param {string} x */ function f(x) {}
/** @constructor */ function A() {}
goog.reflect.object(A, {x: f(1 + 1)});
`, jutil.StringPtr(`actual parameter 1 of f does not match formal parameter
found   : number
required: string
`))
}

// Java: TypeCheckTest.testLends1()
func (c *TypeCheckTest) TestLends1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function extend(x, y) {}
/** @constructor */ function Foo() {}
extend(Foo, /** @lends */ ({bar: 1}));
`).
		AddDiagnostic("Bad type annotation. missing object name in @lends tag." + parsing.JsDocInfoParser_BAD_TYPE_WIKI_LINK).
		Run()
}

// Java: TypeCheckTest.testLends2()
func (c *TypeCheckTest) TestLends2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function extend(x, y) {}
/** @constructor */ function Foo() {}
extend(Foo, /** @lends {Foob} */ ({bar: 1}));
`).
		AddDiagnostic("Variable Foob not declared before @lends annotation.").
		Run()
}

// Java: TypeCheckTest.testLends3()
func (c *TypeCheckTest) TestLends3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function extend(x, y) {}
/** @constructor */ function Foo() {}
extend(Foo, {bar: 1});
alert(Foo.bar);
`).
		AddDiagnostic("Property bar never defined on Foo").
		Run()
}

// Java: TypeCheckTest.testLends4()
func (c *TypeCheckTest) TestLends4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function extend(x, y) {}
/** @constructor */ function Foo() {}
extend(Foo, /** @lends {Foo} */ ({bar: 1}));
alert(Foo.bar);
`).
		Run()
}

// Java: TypeCheckTest.testLends5()
func (c *TypeCheckTest) TestLends5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function extend(x, y) {}
/** @constructor */ function Foo() {}
extend(Foo, {bar: 1});
alert((new Foo()).bar);
`).
		AddDiagnostic("Property bar never defined on Foo").
		Run()
}

// Java: TypeCheckTest.testLends6()
func (c *TypeCheckTest) TestLends6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function extend(x, y) {}
/** @constructor */ function Foo() {}
extend(Foo, /** @lends {Foo.prototype} */ ({bar: 1}));
alert((new Foo()).bar);
`).
		Run()
}

// Java: TypeCheckTest.testLends7()
func (c *TypeCheckTest) TestLends7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function extend(x, y) {}
/** @constructor */ function Foo() {}
extend(Foo, /** @lends {Foo.prototype|Foo} */ ({bar: 1}));
`).
		AddDiagnostic("Bad type annotation. expected closing }" + parsing.JsDocInfoParser_BAD_TYPE_WIKI_LINK).
		Run()
}

// Java: TypeCheckTest.testLends8()
func (c *TypeCheckTest) TestLends8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function extend(x, y) {}
/** @type {number} */ var Foo = 3;
extend(Foo, /** @lends {Foo} */ ({bar: 1}));
`).
		AddDiagnostic("May only lend properties to object types. Foo has type number.").
		Run()
}

// Java: TypeCheckTest.testLends9()
func (c *TypeCheckTest) TestLends9() {
	c.TestClosureTypesMultipleWarnings(`function extend(x, y) {}
/** @constructor */ function Foo() {}
extend(Foo, /** @lends {!Foo} */ ({bar: 1}));
`, []string{"Bad type annotation. expected closing }" + parsing.JsDocInfoParser_BAD_TYPE_WIKI_LINK, "Bad type annotation. missing object name in @lends tag." + parsing.JsDocInfoParser_BAD_TYPE_WIKI_LINK})
}

// Java: TypeCheckTest.testDeclaredNativeTypeEquality()
func (c *TypeCheckTest) TestDeclaredNativeTypeEquality() {
	var n *rhino.Node = c.ParseAndTypeCheck("/** @constructor */ function Object() {};")
	c.AssertTypeEquals(c.registry.GetNativeType(rhino.JSTypeNative_OBJECT_FUNCTION_TYPE), n.GetFirstChild().GetJSType())
}

// Java: TypeCheckTest.testUndefinedVar()
func (c *TypeCheckTest) TestUndefinedVar() {
	var n *rhino.Node = c.ParseAndTypeCheck("var undefined;")
	c.AssertTypeEquals(c.registry.GetNativeType(rhino.JSTypeNative_VOID_TYPE), n.GetFirstFirstChild().GetJSType())
}

// Java: TypeCheckTest.testFlowScopeBug1()
func (c *TypeCheckTest) TestFlowScopeBug1() {
	var n *rhino.Node = c.ParseAndTypeCheck(`/** @param {number} a
* @param {number} b */
function f(a, b) {
/** @type {number} */
var i = 0;
for (; (i + a) < b; ++i) {}}
`)
	c.AssertTypeEquals(c.registry.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE), n.GetFirstChild().GetLastChild().
		GetLastChild().
		GetFirstChild().
		GetNext().
		GetFirstChild().
		GetJSType())
}

// Java: TypeCheckTest.testFlowScopeBug2()
func (c *TypeCheckTest) TestFlowScopeBug2() {
	var n *rhino.Node = c.ParseAndTypeCheck(`/** @constructor */ function Foo() {};
Foo.prototype.hi = false;
function foo(a, b) {
  /** @type {Array} */
  var arr;
  /** @type {number} */
  var iter;
  for (iter = 0; iter < arr.length; ++ iter) {
    /** @type {Foo} */
    var afoo = arr[iter];
    afoo;
  }
}
`)
	c.AssertTypeEquals(c.registry.CreateNullableType(c.registry.GetGlobalType("Foo")), n.GetLastChild().GetLastChild().
		GetLastChild().
		GetLastChild().
		GetLastChild().
		GetLastChild().
		GetJSType())
}

// Java: TypeCheckTest.testAddSingletonGetter()
func (c *TypeCheckTest) TestAddSingletonGetter() {
	var n *rhino.Node = c.ParseAndTypeCheck(`/** @constructor */ function Foo() {};
goog.addSingletonGetter(Foo);
`)
	var o rhino.ObjectType = n.GetFirstChild().GetJSType().(rhino.ObjectType)
	AssertThat(o.GetPropertyType("getInstance").String()).IsEqualTo("function(): Foo")
	AssertThat(o.GetPropertyType("instance_").String()).IsEqualTo("Foo")
}

// Java: TypeCheckTest.testTypeCheckStandaloneAST()
func (c *TypeCheckTest) TestTypeCheckStandaloneAST() {
	var externs *rhino.Node = rhino.IR_Root()
	var firstScript *rhino.Node = c.compiler.ParseTestCode("function Foo() { }")
	c.TypeCheck(firstScript)
	var root *rhino.Node = rhino.IR_Root(externs, rhino.IR_Root(firstScript.Detach()))
	var scopeCreator *TypedScopeCreator = NewTypedScopeCreator(c.compiler)
	var topScope *TypedScope = scopeCreator.CreateScope(root, nil).(*TypedScope)
	var secondScript *rhino.Node = c.compiler.ParseTestCode("new Foo")
	firstScript.ReplaceWith(secondScript)
	NewTypeCheckWithScope(c.compiler, NewSemanticReverseAbstractInterpreter(c.registry), c.registry, topScope, scopeCreator).Process(externs, secondScript.GetParent())
	AssertThat(c.compiler.GetWarningCount()).IsEqualTo(1)
	AssertThat(c.compiler.GetWarnings()[0].
		Description()).IsEqualTo("cannot instantiate non-constructor, found type: function(): undefined")
}

// Java: TypeCheckTest.testUpdateParameterTypeOnClosure()
func (c *TypeCheckTest) TestUpdateParameterTypeOnClosure() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/**
* @constructor
* @param {*=} opt_value
* @return {!Object}
*/
function Object(opt_value) {}
/**
* @constructor
* @param {...*} var_args
*/
function Function(var_args) {}
/**
* @type {Function}
*/
Object.prototype.constructor = function() {};
`).
		AddSource(`/**
* @param {function(): boolean} fn
*/
function f(fn) {}
f(function() { });
`).
		Run()
}

// Java: TypeCheckTest.testTemplatedThisType1()
func (c *TypeCheckTest) TestTemplatedThisType1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */
function Foo() {}
/**
 * @this {T}
 * @return {T}
 * @template T
 */
Foo.prototype.method = function() {};
/**
 * @constructor
 * @extends {Foo}
 */
function Bar() {}
var g = new Bar().method();
/**
 * @param {number} a
 */
function compute(a) {};
compute(g);
`).
		AddDiagnostic(`actual parameter 1 of compute does not match formal parameter
found   : Bar
required: number
`).
		Run()
}

// Java: TypeCheckTest.testTemplatedThisType2()
func (c *TypeCheckTest) TestTemplatedThisType2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @this {Array<T>|{length:number}}
 * @return {T}
 * @template T
 */
Array.prototype.method = function() {};
(function(){
  Array.prototype.method.call(arguments);
})();
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testTemplateType1()
func (c *TypeCheckTest) TestTemplateType1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
* @param {T} x
* @param {T} y
* @param {function(this:T, ...)} z
* @template T
*/
function f(x, y, z) {}
f(this, this, function() { this });
`).
		Run()
}

// Java: TypeCheckTest.testTemplateType2()
func (c *TypeCheckTest) TestTemplateType2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
* @param {T} x
* @param {function(this:T, ...)} y
* @template T
*/
function f(x, y) {}
f(0, function() {});
`).
		Run()
}

// Java: TypeCheckTest.testTemplateType3()
func (c *TypeCheckTest) TestTemplateType3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {T} v
 * @param {function(T)} f
 * @template T
 */
function call(v, f) { f.call(null, v); }
/** @type {string} */ var s;
call(3, function(x) {
 x = true;
 s = x;
});
`).
		AddDiagnostic(`assignment
found   : boolean
required: string
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testTemplateType4()
func (c *TypeCheckTest) TestTemplateType4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {...T} p
 * @return {T}
 * @template T
 */
function fn(p) { return p; }
/** @type {!Object} */ var x;
x = fn(3, null);
`).
		AddDiagnostic(`assignment
found   : (null|number)
required: Object
`).
		Run()
}

// Java: TypeCheckTest.testTemplateType5()
func (c *TypeCheckTest) TestTemplateType5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`const CGI_PARAM_RETRY_COUNT = 'rc';

/**
 * @param {...T} p
 * @return {T}
 * @template T
 */
function fn(p) { return p; }
/** @type {!Object} */ var x;

/** @return {void} */
function aScope() {
  x = fn(CGI_PARAM_RETRY_COUNT, 1);
}
`).
		AddDiagnostic(`assignment
found   : (number|string)
required: Object
`).
		Run()
}

// Java: TypeCheckTest.testTemplateType6()
func (c *TypeCheckTest) TestTemplateType6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {Array<T>} arr
 * @param {?function(T)} f
 * @return {T}
 * @template T
 */
function fn(arr, f) { return arr[0]; }
/** @param {Array<number>} arr */ function g(arr) {
  /** @type {!Object} */ var x = fn.call(null, arr, null);
}
`).
		AddDiagnostic(`initializing variable
found   : number
required: Object
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testTemplateType7()
func (c *TypeCheckTest) TestTemplateType7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddArray().Build()).
		AddSource(`/** @type {!Array<string>} */
var query = [];
query.push(1);
`).
		AddDiagnostic(`actual parameter 1 of Array.prototype.push does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testTemplateType8()
func (c *TypeCheckTest) TestTemplateType8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor
 * @template S,T
 */
function Bar() {}
/**
 * @param {Bar<T>} bar
 * @return {T}
 * @template T
 */
function fn(bar) {}
/** @param {Bar<number>} bar */ function g(bar) {
  /** @type {!Object} */ var x = fn(bar);
}
`).
		AddDiagnostic(`initializing variable
found   : number
required: Object
`).
		Run()
}

// Java: TypeCheckTest.testTemplateType9()
func (c *TypeCheckTest) TestTemplateType9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface
 * @template S,T
 */
function Bar() {}
/**
 * @param {Bar<T>} bar
 * @return {T}
 * @template T
 */
function fn(bar) {}
/** @param {Bar<number>} bar */ function g(bar) {
  /** @type {!Object} */ var x = fn(bar);
}
`).
		AddDiagnostic(`initializing variable
found   : number
required: Object
`).
		Run()
}

// Java: TypeCheckTest.testTemplateType10()
func (c *TypeCheckTest) TestTemplateType10() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor
 * @template T
 */
function Bar() {}

/** @type {!Bar<?>} */ var x;
/** @type {!Bar<number>} */ var y;
y = x;
`).
		Run()
}

// Java: TypeCheckTest.testTemplateType11()
func (c *TypeCheckTest) TestTemplateType11() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor
 * @template T
 */
function Foo() {}

/** @constructor
 * @extends {Foo<string>}
 */
function A() {}

/** @constructor
 * @extends {Foo<number>}
 */
function B() {}

/** @type {!Foo<string>} */ var a = new A();
/** @type {!Foo<string>} */ var b = new B();
`).
		AddDiagnostic(`initializing variable
found   : B
required: Foo<string>
`).
		Run()
}

// Java: TypeCheckTest.testTemplateType12()
func (c *TypeCheckTest) TestTemplateType12() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface
 * @template T
 */
function Foo() {}

/** @constructor
 * @implements {Foo<string>}
 */
function A() {}

/** @constructor
 * @implements {Foo<number>}
 */
function B() {}

/** @type {!Foo<string>} */ var a = new A();
/** @type {!Foo<string>} */ var b = new B();
`).
		AddDiagnostic(`initializing variable
found   : B
required: Foo<string>
`).
		Run()
}

// Java: TypeCheckTest.testTemplateType13()
func (c *TypeCheckTest) TestTemplateType13() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor
 * @template T
 */
function Foo() {}

/** @constructor
 * @template T
 * @extends {Foo<T>}
 */
function A() {}

var a1 = new A();
var a2 = /** @type {!A<string>} */ (new A());
var a3 = /** @type {!A<number>} */ (new A());
/** @type {!Foo<string>} */ var f1 = a1;
/** @type {!Foo<string>} */ var f2 = a2;
/** @type {!Foo<string>} */ var f3 = a3;
`).
		AddDiagnostic(`initializing variable
found   : A<number>
required: Foo<string>
`).
		Run()
}

// Java: TypeCheckTest.testTemplateType14()
func (c *TypeCheckTest) TestTemplateType14() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface
 * @template T
 */
function Foo() {}

/** @constructor
 * @template T
 * @implements {Foo<T>}
 */
function A() {}

var a1 = new A();
var a2 = /** @type {!A<string>} */ (new A());
var a3 = /** @type {!A<number>} */ (new A());
/** @type {!Foo<string>} */ var f1 = a1;
/** @type {!Foo<string>} */ var f2 = a2;
/** @type {!Foo<string>} */ var f3 = a3;
`).
		AddDiagnostic(`initializing variable
found   : A<number>
required: Foo<string>
`).
		Run()
}

// Java: TypeCheckTest.testTemplateType15()
func (c *TypeCheckTest) TestTemplateType15() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {{foo:T}} p
 * @return {T}
 * @template T
 */
function fn(p) { return p.foo; }
/** @type {!Object} */ var x;
x = fn({foo:3});
`).
		AddDiagnostic(`assignment
found   : number
required: Object
`).
		Run()
}

// Java: TypeCheckTest.testTemplateType16()
func (c *TypeCheckTest) TestTemplateType16() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function C() {
  /** @type {number} */ this.foo = 1
}
/**
 * @param {{foo:T}} p
 * @return {T}
 * @template T
 */
function fn(p) { return p.foo; }
/** @type {!Object} */ var x;
x = fn(new C());
`).
		AddDiagnostic(`assignment
found   : number
required: Object
`).
		Run()
}

// Java: TypeCheckTest.testTemplateType17()
func (c *TypeCheckTest) TestTemplateType17() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function C() {}
C.prototype.foo = 1;
/**
 * @param {{foo:T}} p
 * @return {T}
 * @template T
 */
function fn(p) { return p.foo; }
/** @type {!Object} */ var x;
x = fn(new C());
`).
		AddDiagnostic(`assignment
found   : number
required: Object
`).
		Run()
}

// Java: TypeCheckTest.testTemplateType18()
func (c *TypeCheckTest) TestTemplateType18() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function C() {}
C.prototype.foo = 1;
/**
 * @param {{foo:T}} p
 * @return {T}
 * @template T
 */
function fn(p) { return p.foo; }
/** @type {!Object} */ var x;
x = fn({});
`).
		Run()
}

// Java: TypeCheckTest.testTemplateType19()
func (c *TypeCheckTest) TestTemplateType19() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {T} t
 * @param {U} u
 * @return {{t:T, u:U}} 
 * @template T,U
 */
function fn(t, u) { return {t:t, u:u}; }
/** @type {null} */ var x = fn(1, 'str');
`).
		AddDiagnostic(`initializing variable
found   : {
  t: number,
  u: string
}
required: null
`).
		Run()
}

// Java: TypeCheckTest.testTemplateType20()
func (c *TypeCheckTest) TestTemplateType20() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function C() {
  /** @type {void} */ this.x;
}
/**
* @param {T} x
* @param {function(this:T, ...)} y
* @template T
*/
function f(x, y) {}
f(new C, /** @param {number} a */ function(a) {this.x = a;});
`).
		AddDiagnostic(`assignment to property x of C
found   : number
required: undefined
`).
		Run()
}

// Java: TypeCheckTest.testTemplateType21()
func (c *TypeCheckTest) TestTemplateType21() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface @template T */ function A() {}
/** @constructor @implements {A<Foo>} */
function Foo() {}
/** @constructor @implements {A<Bar>} */
function Bar() {}
/** @type {!Foo} */
var x = new Bar();
`).
		AddDiagnostic(`initializing variable
found   : Bar
required: Foo
`).
		Run()
}

// Java: TypeCheckTest.testTemplateType22()
func (c *TypeCheckTest) TestTemplateType22() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface @template T */ function A() {}
/** @interface @template T */ function B() {}
/** @constructor @implements {A<Foo>} */
function Foo() {}
/** @constructor @implements {B<Foo>} */
function Bar() {}
/** @constructor @implements {B<Foo>} */
function Qux() {}
/** @type {!Qux} */
var x = new Bar();
`).
		AddDiagnostic(`initializing variable
found   : Bar
required: Qux
`).
		Run()
}

// Java: TypeCheckTest.testTemplateType23()
func (c *TypeCheckTest) TestTemplateType23() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface @template T */ function A() {}
/** @constructor @implements {A<Foo>} */
function Foo() {}
/** @type {!Foo} */
var x = new Foo();
`).
		Run()
}

// Java: TypeCheckTest.testTemplateType24()
func (c *TypeCheckTest) TestTemplateType24() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @template T
 * @param {T} x
 */
function Foo(x) {
  /** @type {T} */
  this.p = x;
}
/** @return {Foo<Foo<T>>} */
Foo.prototype.m = function() {
  return null;
};
/** @return {T} */
Foo.prototype.get = function() {
  return this.p;
};
var /** null */ n = new Foo(new Object).m().get();
`).
		AddDiagnostic(`initializing variable
found   : (Foo<Object>|null)
required: null
`).
		Run()
}

// Java: TypeCheckTest.testTemplateType25()
func (c *TypeCheckTest) TestTemplateType25() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @template T
 * @param {T} x
 */
function Foo(x) {
  /** @type {T} */
  this.p = x;
}
/** @return {!Foo<!Foo<T>>} */
Foo.prototype.m = function() {
  return new Foo(new Foo(new Object));
};
/** @return {T} */
Foo.prototype.get = function() {
  return this.p;
};
var /** null */ n = new Foo(new Object).m().get();
`).
		AddDiagnostic(`initializing variable
found   : Foo<Object>
required: null
`).
		Run()
}

// Java: TypeCheckTest.testTemplateType26()
func (c *TypeCheckTest) TestTemplateType26() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {T} bar
 * @constructor
 * @template T
 */
function Bar(bar) {
  /** @type {T} */
  this.bar = bar;
}
/** @return {T} */
Bar.prototype.getBar = function() {
  return this.bar;
};
/**
 * @param {T} foo
 * @constructor
 * @template T
 * @extends {Bar<!Array<T>>}
 */
function Foo(foo) {
  /** @type {T} */
  this.foo = foo;
}
var /** null */ n = new Foo(new Object).getBar();
`).
		AddDiagnostic(`initializing variable
found   : Array<Object>
required: null
`).
		Run()
}

// Java: TypeCheckTest.testTemplateTypeCollidesWithParameter()
func (c *TypeCheckTest) TestTemplateTypeCollidesWithParameter() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {T} T
 * @template T
 */
function f(T) {}
`).
		AddDiagnostic(`variable T redefined with type undefined, original definition at testcode0:5 with type T
`).
		Run()
}

// Java: TypeCheckTest.testTemplateTypeForwardReference()
func (c *TypeCheckTest) TestTemplateTypeForwardReference() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {!Foo<string>} x */
function f(x) {}

/**
 * @template T
 * @constructor
 */
function Foo() {}

/** @param {!Foo<number>} x */
function g(x) {
  f(x);
}
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : Foo<number>
required: Foo<string>
`).
		Run()
}

// Java: TypeCheckTest.testTemplateTypeForwardReference_declared()
func (c *TypeCheckTest) TestTemplateTypeForwardReference_declared() {
	c.compiler.ForwardDeclareType("Foo")
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {!Foo<string>} x */
function f(x) {}

/**
 * @template T
 * @constructor
 */
function Foo() {}

/** @param {!Foo<number>} x */
function g(x) {
  f(x);
}
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : Foo<number>
required: Foo<string>
`).
		Run()
}

// Java: TypeCheckTest.testTemplateTypeForwardReferenceFunctionWithExtra()
func (c *TypeCheckTest) TestTemplateTypeForwardReferenceFunctionWithExtra() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {!Foo<string, boolean>} x */
function f(x) {}

/**
 * @constructor
 * @template T
 */
function Foo() {}

/** @param {!Foo<number>} x */
function g(x) {
  f(x);
}
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : Foo<number>
required: Foo<string>
`).
		Run()
}

// Java: TypeCheckTest.testTemplateTypeForwardReferenceVar()
func (c *TypeCheckTest) TestTemplateTypeForwardReferenceVar() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {!Foo<string>} x */
function f(x) {}

/**
 * @template T
 * @constructor
 */
var Foo = function() {}

/** @param {!Foo<number>} x */
function g(x) {
  f(x);
}
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : Foo<number>
required: Foo<string>
`).
		Run()
}

// Java: TypeCheckTest.testTemplateTypeForwardReference_declaredMissing()
func (c *TypeCheckTest) TestTemplateTypeForwardReference_declaredMissing() {
	c.compiler.ForwardDeclareType("Foo")
	c.compiler.ForwardDeclareType("DoesNotExist")
	c.NewTestLegacy().AddSource(`/** @param {!Foo<DoesNotExist>} x */
function f(x) {}
`).
		Run()
}

// Java: TypeCheckTest.testTemplateTypeForwardReference_extends()
func (c *TypeCheckTest) TestTemplateTypeForwardReference_extends() {
	c.compiler.ForwardDeclareType("Bar")
	c.compiler.ForwardDeclareType("Baz")
	c.NewTestLegacy().AddSource(`/** @constructor @extends {Bar<Baz>} */
function Foo() {}
/** @constructor */
function Bar() {}
`).
		Run()
}

// Java: TypeCheckTest.testSubtypeNotTemplated1()
func (c *TypeCheckTest) TestSubtypeNotTemplated1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface @template T */ function A() {}
/** @constructor @implements {A<U>} @template U */ function Foo() {}
function f(/** (!Object|!Foo<string>) */ x) {
  var /** null */ n = x;
}
`).
		AddDiagnostic(`initializing variable
found   : Object
required: null
`).
		Run()
}

// Java: TypeCheckTest.testSubtypeNotTemplated2()
func (c *TypeCheckTest) TestSubtypeNotTemplated2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface @template T */ function A() {}
/** @constructor @implements {A<U>} @template U */ function Foo() {}
function f(/** (!Object|!Foo) */ x) {
  var /** null */ n = x;
}
`).
		AddDiagnostic(`initializing variable
found   : Object
required: null
`).
		Run()
}

// Java: TypeCheckTest.testTemplateTypeWithUnresolvedType()
func (c *TypeCheckTest) TestTemplateTypeWithUnresolvedType() {
	c.TestClosureTypes(`goog.forwardDeclare('Color');
/** @interface @template T */ function C() {}
/** @return {!Color} */ C.prototype.method;
/** @constructor @implements {C} */ function D() {}
/** @override */ D.prototype.method = function() {};
`, nil)
}

// Java: TypeCheckTest.testTemplateTypeWithTypeDef1a()
func (c *TypeCheckTest) TestTemplateTypeWithTypeDef1a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @template T
 * @param {T} x
 */
function Generic(x) {}

/** @constructor */
function Foo() {}

/** @typedef {!Foo} */
var Bar;

/** @type {Generic<!Foo>} */ var x;
/** @type {Generic<!Bar>} */ var y;

x = y;
/** @type {null} */ var z1 = y;
`).
		AddDiagnostic(`initializing variable
found   : (Generic<Foo>|null)
required: null
`).
		Run()
}

// Java: TypeCheckTest.testTemplateTypeWithTypeDef1b()
func (c *TypeCheckTest) TestTemplateTypeWithTypeDef1b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @template T
 * @param {T} x
 */
function Generic(x) {}

/** @constructor */
function Foo() {}

/** @typedef {!Foo} */
var Bar;

/** @type {Generic<!Foo>} */ var x;
/** @type {Generic<!Bar>} */ var y;

y = x;
/** @type {null} */ var z1 = x;
`).
		AddDiagnostic(`initializing variable
found   : (Generic<Foo>|null)
required: null
`).
		Run()
}

// Java: TypeCheckTest.testTemplateTypeWithTypeDef2a()
func (c *TypeCheckTest) TestTemplateTypeWithTypeDef2a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @template T
 * @param {T} x
 */
function Generic(x) {}

/** @constructor */
function Foo() {}

/** @typedef {!Foo} */
var Bar;

function f(/** Generic<!Bar> */ x) {}
/** @type {Generic<!Foo>} */ var x;
f(x);
`).
		Run()
}

// Java: TypeCheckTest.testTemplateTypeWithTypeDef2b()
func (c *TypeCheckTest) TestTemplateTypeWithTypeDef2b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @template T
 * @param {T} x
 */
function Generic(x) {}

/** @constructor */
function Foo() {}

/** @typedef {!Foo} */
var Bar;

function f(/** Generic<!Bar> */ x) {}
/** @type {Generic<!Bar>} */ var x;
f(x);
`).
		Run()
}

// Java: TypeCheckTest.testTemplateTypeWithTypeDef2c()
func (c *TypeCheckTest) TestTemplateTypeWithTypeDef2c() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @template T
 * @param {T} x
 */
function Generic(x) {}

/** @constructor */
function Foo() {}

/** @typedef {!Foo} */
var Bar;

function f(/** Generic<!Foo> */ x) {}
/** @type {Generic<!Foo>} */ var x;
f(x);
`).
		Run()
}

// Java: TypeCheckTest.testTemplateTypeWithTypeDef2d()
func (c *TypeCheckTest) TestTemplateTypeWithTypeDef2d() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @template T
 * @param {T} x
 */
function Generic(x) {}

/** @constructor */
function Foo() {}

/** @typedef {!Foo} */
var Bar;

function f(/** Generic<!Foo> */ x) {}
/** @type {Generic<!Bar>} */ var x;
f(x);
`).
		Run()
}

// Java: TypeCheckTest.testTemplatedFunctionInUnion1()
func (c *TypeCheckTest) TestTemplatedFunctionInUnion1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
* @param {T} x
* @param {function(this:T, ...)|{fn:Function}} z
* @template T
*/
function f(x, z) {}
f([], function() { /** @type {string} */ var x = this });
`).
		AddDiagnostic(`initializing variable
found   : Array<?>
required: string
`).
		Run()
}

// Java: TypeCheckTest.testTemplateTypeRecursion1()
func (c *TypeCheckTest) TestTemplateTypeRecursion1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @typedef {{a: D2}} */
var D1;

/** @typedef {{b: D1}} */
var D2;

fn(x);


/**
 * @param {!D1} s
 * @template T
 */
var fn = function(s) {};
`).
		Run()
}

// Java: TypeCheckTest.testTemplateTypeRecursion2()
func (c *TypeCheckTest) TestTemplateTypeRecursion2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @typedef {{a: D2}} */
var D1;

/** @typedef {{b: D1}} */
var D2;

/** @type {D1} */ var x;
fn(x);


/**
 * @param {!D1} s
 * @template T
 */
var fn = function(s) {};
`).
		Run()
}

// Java: TypeCheckTest.testTemplateTypeRecursion3()
func (c *TypeCheckTest) TestTemplateTypeRecursion3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @typedef {{a: function(D2)}} */
var D1;

/** @typedef {{b: D1}} */
var D2;

/** @type {D1} */ var x;
fn(x);


/**
 * @param {!D1} s
 * @template T
 */
var fn = function(s) {};
`).
		Run()
}

// Java: TypeCheckTest.testFunctionLiteralUndefinedThisArgument()
func (c *TypeCheckTest) TestFunctionLiteralUndefinedThisArgument() {
	c.T().Skip("@Ignore")
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {function(this:T, ...)?} fn
 * @param {?T} opt_obj
 * @template T
 */
function baz(fn, opt_obj) {}
baz(function() { this; });
`).
		AddDiagnostic("Function literal argument refers to undefined this argument").
		Run()
}

// Java: TypeCheckTest.testFunctionLiteralDefinedThisArgument()
func (c *TypeCheckTest) TestFunctionLiteralDefinedThisArgument() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {function(this:T, ...)?} fn
 * @param {?T} opt_obj
 * @template T
 */
function baz(fn, opt_obj) {}
baz(function() { this; }, {});
`).
		Run()
}

// Java: TypeCheckTest.testFunctionLiteralDefinedThisArgument2()
func (c *TypeCheckTest) TestFunctionLiteralDefinedThisArgument2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddArray().Build()).
		AddSource(`/** @param {string} x */ function f(x) {}
/**
 * @param {?function(this:T, ...)} fn
 * @param {T=} opt_obj
 * @template T
 */
function baz(fn, opt_obj) {}
function g() { baz(function() { f(this.length); }, []); }
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testFunctionLiteralUnreadNullThisArgument()
func (c *TypeCheckTest) TestFunctionLiteralUnreadNullThisArgument() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {function(this:T, ...)?} fn
 * @param {?T} opt_obj
 * @template T
 */
function baz(fn, opt_obj) {}
baz(function() {}, null);
`).
		Run()
}

// Java: TypeCheckTest.testUnionTemplateThisType()
func (c *TypeCheckTest) TestUnionTemplateThisType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function F() {}
/** @return {F|Array} */ function g() { return []; }
/** @param {F} x */ function h(x) { }
/**
* @param {T} x
* @param {function(this:T, ...)} y
* @template T
*/
function f(x, y) {}
f(g(), function() { h(this); });
`).
		AddDiagnostic(`actual parameter 1 of h does not match formal parameter
found   : (Array|F|null)
required: (F|null)
`).
		Run()
}

// Java: TypeCheckTest.testRecordType1()
func (c *TypeCheckTest) TestRecordType1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {{prop: number}} x */
function f(x) {}
f({});
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : {prop: (number|undefined)}
required: {prop: number}
missing : []
mismatch: [prop]
`).
		Run()
}

// Java: TypeCheckTest.testRecordType2()
func (c *TypeCheckTest) TestRecordType2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {{prop: (number|undefined)}} x */
function f(x) {}
f({});
`).
		Run()
}

// Java: TypeCheckTest.testRecordType3()
func (c *TypeCheckTest) TestRecordType3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {{prop: number}} x */
function f(x) {}
f({prop: 'x'});
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : {prop: (number|string)}
required: {prop: number}
missing : []
mismatch: [prop]
`).
		Run()
}

// Java: TypeCheckTest.testRecordType4()
func (c *TypeCheckTest) TestRecordType4() {
	c.TestClosureTypesMultipleWarnings(`/** @param {{prop: (number|undefined)}} x */
function f(x) {}
/** @param {{prop: (string|undefined)}} x */
function g(x) {}
var x = {};
f(x);
g(x);
`, []string{`actual parameter 1 of f does not match formal parameter
found   : {prop: (number|string|undefined)}
required: {prop: (number|undefined)}
missing : []
mismatch: [prop]`, `actual parameter 1 of g does not match formal parameter
found   : {prop: (number|string|undefined)}
required: {prop: (string|undefined)}
missing : []
mismatch: [prop]`})
}

// Java: TypeCheckTest.testRecordType5()
func (c *TypeCheckTest) TestRecordType5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {{prop: (number|undefined)}} x */
function f(x) {}
/** @param {{otherProp: (string|undefined)}} x */
function g(x) {}
var x = {}; f(x); g(x);
`).
		Run()
}

// Java: TypeCheckTest.testRecordType6()
func (c *TypeCheckTest) TestRecordType6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {{prop: (number|undefined)}} x */
function f() { return {}; }
`).
		Run()
}

// Java: TypeCheckTest.testRecordType7()
func (c *TypeCheckTest) TestRecordType7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {{prop: (number|undefined)}} x */
function f() { var x = {}; g(x); return x; }
/** @param {number} x */
function g(x) {}
`).
		AddDiagnostic(`actual parameter 1 of g does not match formal parameter
found   : {prop: (number|undefined)}
required: number
`).
		Run()
}

// Java: TypeCheckTest.testRecordType8()
func (c *TypeCheckTest) TestRecordType8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {{prop: (number|string)}} x */
function f() { var x = {prop: 3}; g(x.prop); return x; }
/** @param {string} x */
function g(x) {}
`).
		AddDiagnostic(`actual parameter 1 of g does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testDuplicateRecordFields1()
func (c *TypeCheckTest) TestDuplicateRecordFields1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {{x:string, x:number}} a
 */
function f(a) {};
`).
		AddDiagnostic("Bad type annotation. Duplicate record field x." + parsing.JsDocInfoParser_BAD_TYPE_WIKI_LINK).
		Run()
}

// Java: TypeCheckTest.testDuplicateRecordFields2()
func (c *TypeCheckTest) TestDuplicateRecordFields2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {{name:string,number:x,number:y}} a
 */
function f(a) {};
`).
		AddDiagnostic("Bad type annotation. Unknown type x").
		AddDiagnostic("Bad type annotation. Duplicate record field number." + parsing.JsDocInfoParser_BAD_TYPE_WIKI_LINK).
		Run()
}

// Java: TypeCheckTest.testMultipleExtendsInterface1()
func (c *TypeCheckTest) TestMultipleExtendsInterface1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ function base1() {}
/** @interface */ function base2() {}
/** @interface
* @extends {base1}
* @extends {base2}
*/
function derived() {}
`).
		Run()
}

// Java: TypeCheckTest.testMultipleExtendsInterface2()
func (c *TypeCheckTest) TestMultipleExtendsInterface2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */function Int0() {};
/** @interface */function Int1() {};
Int0.prototype.foo = function() {};
/** @interface
 * @extends {Int0}
 * @extends {Int1} */
function Int2() {};
/** @constructor
 * @implements {Int2} */function Foo() {};
`).
		AddDiagnostic("property foo on interface Int0 is not implemented by type Foo").
		Run()
}

// Java: TypeCheckTest.testMultipleExtendsInterface3()
func (c *TypeCheckTest) TestMultipleExtendsInterface3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */function Int0() {};
/** @interface */function Int1() {};
Int1.prototype.foo = function() {};
/** @interface
 * @extends {Int0}
 * @extends {Int1} */
function Int2() {};
/** @constructor
 * @implements {Int2} */function Foo() {};
`).
		AddDiagnostic("property foo on interface Int1 is not implemented by type Foo").
		Run()
}

// Java: TypeCheckTest.testMultipleExtendsInterface4()
func (c *TypeCheckTest) TestMultipleExtendsInterface4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */function Int0() {};
/** @interface */function Int1() {};
/** @interface
 * @extends {Int0}
 * @extends {Int1}
 * @extends {number} */
function Int2() {};
/** @constructor
 * @implements {Int2} */function Foo() {};
`).
		AddDiagnostic("Int2 @extends non-object type number").
		Run()
}

// Java: TypeCheckTest.testMultipleExtendsInterface5()
func (c *TypeCheckTest) TestMultipleExtendsInterface5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */function Int0() {};
/** @constructor */function Int1() {};
/** @return {string} x */
/** @interface
 * @extends {Int0}
 * @extends {Int1} */
function Int2() {};
`).
		AddDiagnostic("Int2 cannot extend this type; interfaces can only extend interfaces").
		Run()
}

// Java: TypeCheckTest.testMultipleExtendsInterface6()
func (c *TypeCheckTest) TestMultipleExtendsInterface6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */function Super1() {};
/** @interface */function Super2() {};
/** @param {number} bar */Super2.prototype.foo = function(bar) {};
/** @interface @extends {Super1} @extends {Super2} */function Sub() {};
/** @override @param {string} bar */Sub.prototype.foo =
function(bar) {};
`).
		AddDiagnostic(`mismatch of the foo property on type Sub and the type of the property it overrides from interface Super2
original: function(this:Super2, number): undefined
override: function(this:Sub, string): undefined
`).
		Run()
}

// Java: TypeCheckTest.testMultipleExtendsInterfaceAssignment()
func (c *TypeCheckTest) TestMultipleExtendsInterfaceAssignment() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */var I1 = function() {};
/** @interface */ var I2 = function() {}
/** @interface
@extends {I1}
@extends {I2}*/
var I3 = function() {};
/** @constructor
@implements {I3}*/var T = function() {};
var t = new T();
/** @type {I1} */var i1 = t;
/** @type {I2} */var i2 = t;
/** @type {I3} */var i3 = t;
i1 = i3;
i2 = i3;
`).
		Run()
}

// Java: TypeCheckTest.testMultipleExtendsInterfaceParamPass()
func (c *TypeCheckTest) TestMultipleExtendsInterfaceParamPass() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
var I1 = function() {};
/** @interface */
var I2 = function() {}
/** @interface @extends {I1} @extends {I2} */
var I3 = function() {};
/** @constructor @implements {I3} */
var T = function() {};
var t = new T();
/**
 * @param {I1} x
 * @param {I2} y
 * @param {I3} z
 */
function foo(x,y,z){};
foo(t,t,t)
`).
		Run()
}

// Java: TypeCheckTest.testBadMultipleExtendsClass()
func (c *TypeCheckTest) TestBadMultipleExtendsClass() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function base1() {}
/** @constructor */ function base2() {}
/** @constructor
* @extends {base1}
* @extends {base2}
*/
function derived() {}
`).
		AddDiagnostic("Bad type annotation. type annotation incompatible with other annotations." + parsing.JsDocInfoParser_BAD_TYPE_WIKI_LINK).
		Run()
}

// Java: TypeCheckTest.testInterfaceExtendsResolution()
func (c *TypeCheckTest) TestInterfaceExtendsResolution() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface 
 @extends {A} */ function B() {};
/** @constructor 
 @implements {B} */ function C() {};
/** @interface */ function A() {};
`).
		Run()
}

// Java: TypeCheckTest.testPropertyCanBeDefinedInObject()
func (c *TypeCheckTest) TestPropertyCanBeDefinedInObject() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @interface */ function I() {};
I.prototype.bar = function() {};
/** @type {Object} */ var foo;
foo.bar();
`).
		Run()
}

// Java: TypeCheckTest.testExtendedInterfacePropertiesCompatibility1()
func (c *TypeCheckTest) TestExtendedInterfacePropertiesCompatibility1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */function Int0() {};
/** @interface */function Int1() {};
/** @type {number} */
Int0.prototype.foo;
/** @type {string} */
Int1.prototype.foo;
/** @interface
 * @extends {Int0}
 * @extends {Int1} */
function Int2() {};
`).
		AddDiagnostic(`Interface Int2 has a property foo with incompatible types in its super interfaces Int0 and Int1
`).
		Run()
}

// Java: TypeCheckTest.testExtendedInterfacePropertiesCompatibility2()
func (c *TypeCheckTest) TestExtendedInterfacePropertiesCompatibility2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */function Int0() {};
/** @interface */function Int1() {};
/** @interface */function Int2() {};
/** @type {number} */
Int0.prototype.foo;
/** @type {string} */
Int1.prototype.foo;
/** @type {Object} */
Int2.prototype.foo;
/** @interface
 * @extends {Int0}
 * @extends {Int1}
 * @extends {Int2}*/
function Int3() {};
`).
		AddDiagnostic(`Interface Int3 has a property foo with incompatible types in its super interfaces Int0 and Int1
`).
		AddDiagnostic(`Interface Int3 has a property foo with incompatible types in its super interfaces Int1 and Int2
`).
		Run()
}

// Java: TypeCheckTest.testExtendedInterfacePropertiesCompatibility3()
func (c *TypeCheckTest) TestExtendedInterfacePropertiesCompatibility3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */function Int0() {};
/** @interface */function Int1() {};
/** @type {number} */
Int0.prototype.foo;
/** @type {string} */
Int1.prototype.foo;
/** @interface
 * @extends {Int1} */ function Int2() {};
/** @interface
 * @extends {Int0}
 * @extends {Int2} */
function Int3() {};
`).
		AddDiagnostic(`Interface Int3 has a property foo with incompatible types in its super interfaces Int0 and Int1
`).
		Run()
}

// Java: TypeCheckTest.testExtendedInterfacePropertiesCompatibility4()
func (c *TypeCheckTest) TestExtendedInterfacePropertiesCompatibility4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */function Int0() {};
/** @interface
 * @extends {Int0} */ function Int1() {};
/** @type {number} */
Int0.prototype.foo;
/** @interface */function Int2() {};
/** @interface
 * @extends {Int2} */ function Int3() {};
/** @type {string} */
Int2.prototype.foo;
/** @interface
 * @extends {Int1}
 * @extends {Int3} */
function Int4() {};
`).
		AddDiagnostic(`Interface Int4 has a property foo with incompatible types in its super interfaces Int0 and Int2
`).
		Run()
}

// Java: TypeCheckTest.testExtendedInterfacePropertiesCompatibility5()
func (c *TypeCheckTest) TestExtendedInterfacePropertiesCompatibility5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */function Int0() {};
/** @interface */function Int1() {};
/** @type {number} */
Int0.prototype.foo;
/** @type {string} */
Int1.prototype.foo;
/** @interface
 * @extends {Int1} */ function Int2() {};
/** @interface
 * @extends {Int0}
 * @extends {Int2} */
function Int3() {};
/** @interface */function Int4() {};
/** @type {number} */
Int4.prototype.foo;
/** @interface
 * @extends {Int3}
 * @extends {Int4} */
function Int5() {};
`).
		AddDiagnostic(`Interface Int3 has a property foo with incompatible types in its super interfaces Int0 and Int1
`).
		AddDiagnostic(`Interface Int5 has a property foo with incompatible types in its super interfaces Int1 and Int4
`).
		Run()
}

// Java: TypeCheckTest.testExtendedInterfacePropertiesCompatibility6()
func (c *TypeCheckTest) TestExtendedInterfacePropertiesCompatibility6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */function Int0() {};
/** @interface */function Int1() {};
/** @type {number} */
Int0.prototype.foo;
/** @type {string} */
Int1.prototype.foo;
/** @interface
 * @extends {Int1} */ function Int2() {};
/** @interface
 * @extends {Int0}
 * @extends {Int2} */
function Int3() {};
/** @interface */function Int4() {};
/** @type {string} */
Int4.prototype.foo;
/** @interface
 * @extends {Int3}
 * @extends {Int4} */
function Int5() {};
`).
		AddDiagnostic(`Interface Int3 has a property foo with incompatible types in its super interfaces Int0 and Int1
`).
		Run()
}

// Java: TypeCheckTest.testExtendedInterfacePropertiesCompatibility7()
func (c *TypeCheckTest) TestExtendedInterfacePropertiesCompatibility7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */function Int0() {};
/** @interface */function Int1() {};
/** @type {number} */
Int0.prototype.foo;
/** @type {string} */
Int1.prototype.foo;
/** @interface
 * @extends {Int1} */ function Int2() {};
/** @interface
 * @extends {Int0}
 * @extends {Int2} */
function Int3() {};
/** @interface */function Int4() {};
/** @type {Object} */
Int4.prototype.foo;
/** @interface
 * @extends {Int3}
 * @extends {Int4} */
function Int5() {};
`).
		AddDiagnostic(`Interface Int3 has a property foo with incompatible types in its super interfaces Int0 and Int1
`).
		AddDiagnostic(`Interface Int5 has a property foo with incompatible types in its super interfaces Int1 and Int4
`).
		Run()
}

// Java: TypeCheckTest.testExtendedInterfacePropertiesCompatibility8()
func (c *TypeCheckTest) TestExtendedInterfacePropertiesCompatibility8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */function Int0() {};
/** @interface */function Int1() {};
/** @type {number} */
Int0.prototype.foo;
/** @type {string} */
Int1.prototype.bar;
/** @interface
 * @extends {Int1} */ function Int2() {};
/** @interface
 * @extends {Int0}
 * @extends {Int2} */
function Int3() {};
/** @interface */function Int4() {};
/** @type {Object} */
Int4.prototype.foo;
/** @type {Null} */
Int4.prototype.bar;
/** @interface
 * @extends {Int3}
 * @extends {Int4} */
function Int5() {};
`).
		AddDiagnostic(`Interface Int5 has a property bar with incompatible types in its super interfaces Int1 and Int4
`).
		AddDiagnostic(`Interface Int5 has a property foo with incompatible types in its super interfaces Int0 and Int4
`).
		Run()
}

// Java: TypeCheckTest.testExtendedInterfacePropertiesCompatibility9()
func (c *TypeCheckTest) TestExtendedInterfacePropertiesCompatibility9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface
 * @template T */function Int0() {};
/** @interface
 * @template T */function Int1() {};
/** @type {T} */
Int0.prototype.foo;
/** @type {T} */
Int1.prototype.foo;
/** @interface
 * @extends {Int0<number>}
 * @extends {Int1<string>} */
function Int2() {};
`).
		AddDiagnostic(`Interface Int2 has a property foo with incompatible types in its super interfaces Int0<number> and Int1<string>
`).
		Run()
}

// Java: TypeCheckTest.testExtendedInterfacePropertiesCompatibilityNoError()
func (c *TypeCheckTest) TestExtendedInterfacePropertiesCompatibilityNoError() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */function Int0() {};
/** @interface */function Int1() {};
/** @param {number} x */
Int0.prototype.foo;
/** @param {number} x */
Int1.prototype.foo;
/** @interface
 * @extends {Int0}
 * @extends {Int1} */
function Int2() {};
`).
		Run()
}

// Java: TypeCheckTest.testImplementedInterfacePropertiesShouldFailOnConflict()
func (c *TypeCheckTest) TestImplementedInterfacePropertiesShouldFailOnConflict() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */function Int0() {};
/** @interface */function Int1() {};
/** @type {number} */
Int0.prototype.foo;
/** @type {string} */
Int1.prototype.foo;
/** @constructor @implements {Int0} @implements {Int1} */
function Foo() {};
Foo.prototype.foo;
`).
		AddDiagnostic(`mismatch of the foo property on type Foo and the type of the property it overrides from interface Int1
original: string
override: number
`).
		Run()
}

// Java: TypeCheckTest.testGenerics1a()
func (c *TypeCheckTest) TestGenerics1a() {
	fnDecl := `/** 
 * @param {T} x 
 * @param {function(T):T} y 
 * @template T
 */ 
function f(x,y) { return y(x); }
`
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(fnDecl + `/** @type {string} */
var out;
/** @type {string} */
var result = f('hi', function(x) {
  out = x;
  return x;
});
`).
		Run()
}

// Java: TypeCheckTest.testGenerics1b()
func (c *TypeCheckTest) TestGenerics1b() {
	fnDecl := `/**
 * @param {T} x
 * @param {function(T):T} y
 * @template T
 */
function f(x,y) { return y(x); }
`
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(fnDecl + `/** @type {string} */
var out;
var result = f(0, function(x){ out = x; return x; });
`).
		AddDiagnostic(`assignment
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testFilter0()
func (c *TypeCheckTest) TestFilter0() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {T} arr
 * @return {T}
 * @template T
 */
var filter = function(arr){};
/** @type {!Array<string>} */
var arr;
/** @type {!Array<string>} */
var result = filter(arr);
`).
		Run()
}

// Java: TypeCheckTest.testFilter1()
func (c *TypeCheckTest) TestFilter1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {!Array<T>} arr
 * @return {!Array<T>}
 * @template T
 */
var filter = function(arr){};
/** @type {!Array<string>} */
var arr;
/** @type {!Array<string>} */
var result = filter(arr);
`).
		Run()
}

// Java: TypeCheckTest.testFilter2()
func (c *TypeCheckTest) TestFilter2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {!Array<T>} arr
 * @return {!Array<T>}
 * @template T
 */
var filter = function(arr){};
/** @type {!Array<string>} */
var arr;
/** @type {!Array<number>} */
var result = filter(arr);
`).
		AddDiagnostic(`initializing variable
found   : Array<string>
required: Array<number>
`).
		Run()
}

// Java: TypeCheckTest.testFilter3()
func (c *TypeCheckTest) TestFilter3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {Array<T>} arr
 * @return {Array<T>}
 * @template T
 */
var filter = function(arr){};
/** @type {Array<string>} */
var arr;
/** @type {Array<number>} */
var result = filter(arr);
`).
		AddDiagnostic(`initializing variable
found   : (Array<string>|null)
required: (Array<number>|null)
`).
		Run()
}

// Java: TypeCheckTest.testBackwardsInferenceGoogArrayFilter1()
func (c *TypeCheckTest) TestBackwardsInferenceGoogArrayFilter1() {
	c.TestClosureTypes(`/** @type {Array<string>} */
var arr;
/** @type {!Array<number>} */
var result = goog.array.filter(
   arr,
   function(item,index,src) {return false;});
`, jutil.StringPtr(`initializing variable
found   : Array<string>
required: Array<number>
`))
}

// Java: TypeCheckTest.testBackwardsInferenceGoogArrayFilter2()
func (c *TypeCheckTest) TestBackwardsInferenceGoogArrayFilter2() {
	c.TestClosureTypes(`/** @type {number} */
var out;
/** @type {Array<string>} */
var arr;
var out4 = goog.array.filter(
   arr,
   function(item,index,src) {out = item; return false});
`, jutil.StringPtr(`assignment
found   : string
required: number
`))
}

// Java: TypeCheckTest.testBackwardsInferenceGoogArrayFilter3()
func (c *TypeCheckTest) TestBackwardsInferenceGoogArrayFilter3() {
	c.TestClosureTypes(`/** @type {string} */
var out;
/** @type {Array<string>} */ var arr;
var result = goog.array.filter(
   arr,
   function(item,index,src) {out = index;});
`, jutil.StringPtr(`assignment
found   : number
required: string
`))
}

// Java: TypeCheckTest.testBackwardsInferenceGoogArrayFilter4()
func (c *TypeCheckTest) TestBackwardsInferenceGoogArrayFilter4() {
	c.TestClosureTypes(`/** @type {string} */
var out;
/** @type {Array<string>} */ var arr;
var out4 = goog.array.filter(
   arr,
   function(item,index,srcArr) {out = srcArr;});
`, jutil.StringPtr(`assignment
found   : (Array|null|{length: number})
required: string
`))
}

// Java: TypeCheckTest.testCatchExpression1()
func (c *TypeCheckTest) TestCatchExpression1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function fn() {
  /** @type {number} */
  var out = 0;
  try {
    foo();
  } catch (/** @type {string} */ e) {
    out = e;
  }
}
`).
		AddDiagnostic(`assignment
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testCatchExpression2()
func (c *TypeCheckTest) TestCatchExpression2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function fn() {
  /** @type {number} */
  var out = 0;
  /** @type {string} */
  var e;
  try {
    foo();
  } catch (e) {
    out = e;
  }
}
`).
		Run()
}

// Java: TypeCheckTest.testIssue1058()
func (c *TypeCheckTest) TestIssue1058() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
  * @constructor
  * @template CLASS
  */
var Class = function() {};

/**
  * @param {function(CLASS):CLASS} a
  * @template T
  */
Class.prototype.foo = function(a) {
  return 'string';
};

/** @param {number} a
  * @return {string} */
var a = function(a) { return '' };

new Class().foo(a);
`).
		Run()
}

// Java: TypeCheckTest.testDeterminacyIssue()
func (c *TypeCheckTest) TestDeterminacyIssue() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`(function() {
    /** @constructor */
    var ImageProxy = function() {};
    /** @constructor */
    var FeedReader = function() {};
    /** @type {ImageProxy} */
    FeedReader.x = new ImageProxy();
})();
`).
		Run()
}

// Java: TypeCheckTest.testUnknownTypeReport()
func (c *TypeCheckTest) TestUnknownTypeReport() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function id(x) { return x; }").
		EnableReportUnknownTypes().
		AddDiagnostic("could not determine the type of this expression").
		Run()
}

// Java: TypeCheckTest.testUnknownTypeReport_allowsUnknownIfStatement()
func (c *TypeCheckTest) TestUnknownTypeReport_allowsUnknownIfStatement() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function id(x) { x; }").
		EnableReportUnknownTypes().
		Run()
}

// Java: TypeCheckTest.testUnknownForIn()
func (c *TypeCheckTest) TestUnknownForIn() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var x = {'a':1}; var y;
for(
y
in x) {}
`).
		EnableReportUnknownTypes().
		Run()
}

// Java: TypeCheckTest.testUnknownTypeDisabledByDefault()
func (c *TypeCheckTest) TestUnknownTypeDisabledByDefault() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function id(x) { return x; }").
		Run()
}

// Java: TypeCheckTest.testNonexistentPropertyAccessOnStruct()
func (c *TypeCheckTest) TestNonexistentPropertyAccessOnStruct() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @struct
 */
var A = function() {};
/** @param {A} a */
function foo(a) {
  if (a.bar) { a.bar(); }
}
`).
		AddDiagnostic("Property bar never defined on A").
		Run()
}

// Java: TypeCheckTest.testNonexistentPropertyAccessOnStructOrObject()
func (c *TypeCheckTest) TestNonexistentPropertyAccessOnStructOrObject() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/**
 * @constructor
 * @struct
 */
var A = function() {};
/** @param {A|Object} a */
function foo(a) {
  if (a.bar) { a.bar(); }
}
`).
		Run()
}

// Java: TypeCheckTest.testNonexistentPropertyAccessOnExternStruct()
func (c *TypeCheckTest) TestNonexistentPropertyAccessOnExternStruct() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/**
 * @constructor
 * @struct
 */
var A = function() {};
`).
		AddSource(`/** @param {A} a */
function foo(a) {
  if (a.bar) { a.bar(); }
}
`).
		AddDiagnostic("Property bar never defined on A").
		Run()
}

// Java: TypeCheckTest.testNonexistentPropertyAccessStructSubtype()
func (c *TypeCheckTest) TestNonexistentPropertyAccessStructSubtype() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @struct
 */
var A = function() {};
/**
 * @constructor
 * @struct
 * @extends {A}
 */
var B = function() { this.bar = function(){}; };
/** @param {A} a */
function foo(a) {
  if (a.bar) { a.bar(); }
}
`).
		AddDiagnostic("Property bar never defined on A").
		Run()
}

// Java: TypeCheckTest.testNonexistentPropertyAccessStructInterfaceSubtype()
func (c *TypeCheckTest) TestNonexistentPropertyAccessStructInterfaceSubtype() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @interface
 * @struct
 */
var A = function() {};

/**
 * @interface
 * @struct
 * @extends {A}
 */
var B = function() {};
/** @return {void} */ B.prototype.bar = function(){};

/** @param {A} a */
function foo(a) {
  if (a.bar) { a.bar(); }
}
`).
		AddDiagnostic("Property bar never defined on A").
		Run()
}

// Java: TypeCheckTest.testNonexistentPropertyAccessStructRecordSubtype()
func (c *TypeCheckTest) TestNonexistentPropertyAccessStructRecordSubtype() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/**
 * @record
 * @struct
 */
var A = function() {};

/**
 * @record
 * @struct
 * @extends {A}
 */
var B = function() {};
/** @return {void} */ B.prototype.bar = function(){};

/** @param {A} a */
function foo(a) {
  if (a.bar) {
    a.bar();
  }
}
`).
		AddDiagnostic("Property bar never defined on A").
		Run()
}

// Java: TypeCheckTest.testNonexistentPropertyAccessStructSubtype2()
func (c *TypeCheckTest) TestNonexistentPropertyAccessStructSubtype2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @struct
 */
function Foo() {
  this.x = 123;
}
var objlit = /** @struct */ { y: 234 };
Foo.prototype = objlit;
var n = objlit.x;
`).
		AddDiagnostic("Property x never defined on Foo.prototype").
		Run()
}

// Java: TypeCheckTest.testIssue1024a()
func (c *TypeCheckTest) TestIssue1024a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`/** @param {Object} a */
function f(a) {
  a.prototype = '__proto'
}
/** @param {Object} b
 *  @return {!Object}
 */
function g(b) {
  return b.prototype
}
`).
		Run()
}

// Java: TypeCheckTest.testIssue1024b()
func (c *TypeCheckTest) TestIssue1024b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {Object} a */
function f(a) {
  a.prototype = {foo:3};
}
/** @param {Object} b
 */
function g(b) {
  b.prototype = function(){};
}
`).
		Run()
}

// Java: TypeCheckTest.testModuleReferenceNotAllowed()
func (c *TypeCheckTest) TestModuleReferenceNotAllowed() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @param {./Foo} z */ function f(z) {}").
		AddDiagnostic("Bad type annotation. Unknown type ./Foo").
		Run()
}

// Java: TypeCheckTest.testCheckObjectKeysBadKey1()
func (c *TypeCheckTest) TestCheckObjectKeysBadKey1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {!Object<!Object, number>} */ var k;").
		AddDiagnosticType(TypeCheck_NON_STRINGIFIABLE_OBJECT_KEY).
		Run()
}

// Java: TypeCheckTest.testCheckObjectKeysBadKey2()
func (c *TypeCheckTest) TestCheckObjectKeysBadKey2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {!Object<function(), number>} */ var k;").
		AddDiagnosticType(TypeCheck_NON_STRINGIFIABLE_OBJECT_KEY).
		Run()
}

// Java: TypeCheckTest.testCheckObjectKeysBadKey3()
func (c *TypeCheckTest) TestCheckObjectKeysBadKey3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {!Object<!Array<!Object>, number>} */ var k;").
		AddDiagnosticType(TypeCheck_NON_STRINGIFIABLE_OBJECT_KEY).
		Run()
}

// Java: TypeCheckTest.testCheckObjectKeysBadKey4()
func (c *TypeCheckTest) TestCheckObjectKeysBadKey4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {!Object<*, number>} */ var k;").
		AddDiagnosticType(TypeCheck_NON_STRINGIFIABLE_OBJECT_KEY).
		Run()
}

// Java: TypeCheckTest.testCheckObjectKeysBadKey5()
func (c *TypeCheckTest) TestCheckObjectKeysBadKey5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {(string|Object<Object, number>)} */ var k;").
		AddDiagnosticType(TypeCheck_NON_STRINGIFIABLE_OBJECT_KEY).
		Run()
}

// Java: TypeCheckTest.testCheckObjectKeysBadKey6()
func (c *TypeCheckTest) TestCheckObjectKeysBadKey6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {!Object<number, !Object<Object, number>>} */ var k;").
		AddDiagnosticType(TypeCheck_NON_STRINGIFIABLE_OBJECT_KEY).
		Run()
}

// Java: TypeCheckTest.testCheckObjectKeysBadKey7()
func (c *TypeCheckTest) TestCheckObjectKeysBadKey7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */
var MyClass = function() {};
/** @type {!Object<MyClass, number>} */
var k;
`).
		AddDiagnosticType(TypeCheck_NON_STRINGIFIABLE_OBJECT_KEY).
		Run()
}

// Java: TypeCheckTest.testCheckObjectKeysBadKey8()
func (c *TypeCheckTest) TestCheckObjectKeysBadKey8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum{!Object} */
var Enum = {};
/** @type {!Object<Enum, number>} */
var k;
`).
		AddDiagnosticType(TypeCheck_NON_STRINGIFIABLE_OBJECT_KEY).
		Run()
}

// Java: TypeCheckTest.testCheckObjectKeysBadKey9()
func (c *TypeCheckTest) TestCheckObjectKeysBadKey9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {function(!Object<!Object, number>)} */ var k;").
		AddDiagnosticType(TypeCheck_NON_STRINGIFIABLE_OBJECT_KEY).
		Run()
}

// Java: TypeCheckTest.testCheckObjectKeysBadKey10()
func (c *TypeCheckTest) TestCheckObjectKeysBadKey10() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {function(): !Object<!Object, number>} */ var k;").
		AddDiagnosticType(TypeCheck_NON_STRINGIFIABLE_OBJECT_KEY).
		Run()
}

// Java: TypeCheckTest.testCheckObjectKeysBadKey11()
func (c *TypeCheckTest) TestCheckObjectKeysBadKey11() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */
function X() {}
/** @constructor @extends {X} */
function X2() {}
/** @enum {!X} */
var XE = {A:new X};
/** @type {Object<(!XE|!X2), string>} */
var Y = {};
`).
		AddDiagnosticType(TypeCheck_NON_STRINGIFIABLE_OBJECT_KEY).
		Run()
}

// Java: TypeCheckTest.testCheckObjectKeysVariousTags1()
func (c *TypeCheckTest) TestCheckObjectKeysVariousTags1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {!Object<!Object, number>} */ var k;").
		AddDiagnosticType(TypeCheck_NON_STRINGIFIABLE_OBJECT_KEY).
		Run()
}

// Java: TypeCheckTest.testCheckObjectKeysVariousTags2()
func (c *TypeCheckTest) TestCheckObjectKeysVariousTags2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @param {!Object<!Object, number>} a */ var f = function(a) {};").
		AddDiagnosticType(TypeCheck_NON_STRINGIFIABLE_OBJECT_KEY).
		Run()
}

// Java: TypeCheckTest.testCheckObjectKeysVariousTags3()
func (c *TypeCheckTest) TestCheckObjectKeysVariousTags3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @return {!Object<!Object, number>} */ var f = function() {return {}};").
		AddDiagnosticType(TypeCheck_NON_STRINGIFIABLE_OBJECT_KEY).
		Run()
}

// Java: TypeCheckTest.testCheckObjectKeysVariousTags4()
func (c *TypeCheckTest) TestCheckObjectKeysVariousTags4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @typedef {!Object<!Object, number>} */ var MyType;").
		AddDiagnosticType(TypeCheck_NON_STRINGIFIABLE_OBJECT_KEY).
		Run()
}

// Java: TypeCheckTest.testCheckObjectKeysGoodKey1()
func (c *TypeCheckTest) TestCheckObjectKeysGoodKey1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {!Object<number, number>} */ var k;").
		Run()
}

// Java: TypeCheckTest.testCheckObjectKeysGoodKey2()
func (c *TypeCheckTest) TestCheckObjectKeysGoodKey2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {!Object<string, number>} */ var k;").
		Run()
}

// Java: TypeCheckTest.testCheckObjectKeysGoodKey3()
func (c *TypeCheckTest) TestCheckObjectKeysGoodKey3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {!Object<boolean, number>} */ var k;").
		Run()
}

// Java: TypeCheckTest.testCheckObjectKeysGoodKey4()
func (c *TypeCheckTest) TestCheckObjectKeysGoodKey4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {!Object<null, number>} */ var k;").
		Run()
}

// Java: TypeCheckTest.testCheckObjectKeysGoodKey5()
func (c *TypeCheckTest) TestCheckObjectKeysGoodKey5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {!Object<undefined, number>} */ var k;").
		Run()
}

// Java: TypeCheckTest.testCheckObjectKeysGoodKey6()
func (c *TypeCheckTest) TestCheckObjectKeysGoodKey6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {!Object<!Date, number>} */ var k;").
		Run()
}

// Java: TypeCheckTest.testCheckObjectKeysGoodKey7()
func (c *TypeCheckTest) TestCheckObjectKeysGoodKey7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {!Object<!RegExp, number>} */ var k;").
		Run()
}

// Java: TypeCheckTest.testCheckObjectKeysGoodKey8()
func (c *TypeCheckTest) TestCheckObjectKeysGoodKey8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {!Object<!Array, number>} */ var k;").
		Run()
}

// Java: TypeCheckTest.testCheckObjectKeysGoodKey9()
func (c *TypeCheckTest) TestCheckObjectKeysGoodKey9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {!Object<!Array<number>, number>} */ var k;").
		Run()
}

// Java: TypeCheckTest.testCheckObjectKeysGoodKey10()
func (c *TypeCheckTest) TestCheckObjectKeysGoodKey10() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {!Object<?, number>} */ var k;").
		Run()
}

// Java: TypeCheckTest.testCheckObjectKeysGoodKey11()
func (c *TypeCheckTest) TestCheckObjectKeysGoodKey11() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {!Object<(string|number), number>} */ var k").
		Run()
}

// Java: TypeCheckTest.testCheckObjectKeysGoodKey12()
func (c *TypeCheckTest) TestCheckObjectKeysGoodKey12() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {!Object<Object>} */ var k;").
		Run()
}

// Java: TypeCheckTest.testCheckObjectKeysGoodKey13()
func (c *TypeCheckTest) TestCheckObjectKeysGoodKey13() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
var MyInterface = function() {};
/** @type {!Object<!MyInterface, number>} */
var k;
`).
		Run()
}

// Java: TypeCheckTest.testCheckObjectKeysGoodKey14()
func (c *TypeCheckTest) TestCheckObjectKeysGoodKey14() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @typedef {{a: number}} */ var MyRecord;
/** @type {!Object<MyRecord, number>} */ var k;
`).
		Run()
}

// Java: TypeCheckTest.testCheckObjectKeysGoodKey15()
func (c *TypeCheckTest) TestCheckObjectKeysGoodKey15() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum{number} */
var Enum = {};
/** @type {!Object<Enum, number>} */
var k;
`).
		Run()
}

// Java: TypeCheckTest.testCheckObjectKeysClassWithToString()
func (c *TypeCheckTest) TestCheckObjectKeysClassWithToString() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddObject().Build()).
		AddSource(`/** @constructor */
var MyClass = function() {};
/** @override*/
MyClass.prototype.toString = function() { return ''; };
/** @type {!Object<!MyClass, number>} */
var k;
`).
		Run()
}

// Java: TypeCheckTest.testCheckObjectKeysClassInheritsToString()
func (c *TypeCheckTest) TestCheckObjectKeysClassInheritsToString() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddObject().Build()).
		AddSource(`/** @constructor */
var Parent = function() {};
/** @override */
Parent.prototype.toString = function() { return ''; };
/** @constructor @extends {Parent} */
var Child = function() {};
/** @type {!Object<!Child, number>} */
var k;
`).
		Run()
}

// Java: TypeCheckTest.testCheckObjectKeysForEnumUsingClassWithToString()
func (c *TypeCheckTest) TestCheckObjectKeysForEnumUsingClassWithToString() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddObject().Build()).
		AddSource(`/** @constructor */
var MyClass = function() {};
/** @override*/
MyClass.prototype.toString = function() { return ''; };
/** @enum{!MyClass} */
var Enum = {};
/** @type {!Object<Enum, number>} */
var k;
`).
		Run()
}

// Java: TypeCheckTest.testBadSuperclassInheritance1()
func (c *TypeCheckTest) TestBadSuperclassInheritance1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */
function Foo() {}
/** @type {number} */
Foo.prototype.myprop = 2;

/** @constructor @extends {Foo} */
function Bar() {}
/** @type {number} */
Bar.prototype.myprop = 1;
`).
		AddDiagnosticType(TypeCheck_HIDDEN_SUPERCLASS_PROPERTY).
		Run()
}

// Java: TypeCheckTest.testBadSuperclassInheritance2()
func (c *TypeCheckTest) TestBadSuperclassInheritance2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */
function Foo() {}
/** @type {number} */
Foo.prototype.myprop = 2;

/** @constructor @extends {Foo} */
function Bar() {}
/** @override @type {string} */
Bar.prototype.myprop = 'qwer';
`).
		AddDiagnosticType(TypeValidator_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH).
		Run()
}

// Java: TypeCheckTest.testBadSuperclassInheritance3()
func (c *TypeCheckTest) TestBadSuperclassInheritance3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */
function Foo() {}
/** @type {number} */
Foo.prototype.myprop = 2;

/** @constructor @extends {Foo} */
function Bar() {}
/** @override @type {string} */
Bar.prototype.myprop;
`).
		Run()
}

// Java: TypeCheckTest.testCheckObjectKeysWithNamedType()
func (c *TypeCheckTest) TestCheckObjectKeysWithNamedType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {!Object<!PseudoId, number>} */
var k;
/** @typedef {number|string} */
var PseudoId;
`).
		Run()
}

// Java: TypeCheckTest.testCheckObjectKeyRecursiveType()
func (c *TypeCheckTest) TestCheckObjectKeyRecursiveType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @typedef {!Object<string, !Predicate>} */ var Schema;
/** @typedef {function(*): boolean|!Schema} */ var Predicate;
/** @type {!Schema} */ var k;
`).
		Run()
}

// Java: TypeCheckTest.testDontOverrideNativeScalarTypes()
func (c *TypeCheckTest) TestDontOverrideNativeScalarTypes() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`string = 123;
var /** string */ s = 123;
`).
		AddDiagnostic(`initializing variable
found   : number
required: string
`).
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var string = goog.require('goog.string');
var /** string */ s = 123;
`).
		AddDiagnostic("Property require never defined on goog" + TypeCheck_POSSIBLE_INEXISTENT_PROPERTY_EXPLANATION).
		AddDiagnostic(`initializing variable
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testTemplateMap1()
func (c *TypeCheckTest) TestTemplateMap1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckTest_EXTERNS_WITH_IARRAYLIKE_DECLS).
		AddSource(`function f() {
  /** @type {Int8Array} */
  var x = new Int8Array(10);
  /** @type {IArrayLike<string>} */
  var y;
  y = x;
}
`).
		AddDiagnostic(`assignment
found   : (Int8Array|null)
required: (IArrayLike<string>|null)
`).
		Run()
}

// Java: TypeCheckTest.testTemplateMap2()
func (c *TypeCheckTest) TestTemplateMap2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckTest_EXTERNS_WITH_IARRAYLIKE_DECLS).
		AddSource(`function f() {
  /** @type {Int8Array} */
  var x = new Int8Array(10);

  /** @type {IObject<number, string>} */
  var z;
  z = x;
}
`).
		AddDiagnostic(`assignment
found   : (Int8Array|null)
required: (IObject<number,string>|null)
`).
		Run()
}

// Java: TypeCheckTest.testTemplateMap3()
func (c *TypeCheckTest) TestTemplateMap3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckTest_EXTERNS_WITH_IARRAYLIKE_DECLS).
		AddSource(`function f() {
  var x = new Int8Array(10);

  /** @type {IArrayLike<string>} */
  var y;
  y = x;
}
`).
		AddDiagnostic(`assignment
found   : Int8Array
required: (IArrayLike<string>|null)
`).
		Run()
}

// Java: TypeCheckTest.testTemplateMap4()
func (c *TypeCheckTest) TestTemplateMap4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckTest_EXTERNS_WITH_IARRAYLIKE_DECLS).
		AddSource(`function f() {
  var x = new Int8Array(10);

  /** @type {IObject<number, string>} */
  var z;
  z = x;
}
`).
		AddDiagnostic(`assignment
found   : Int8Array
required: (IObject<number,string>|null)
`).
		Run()
}

// Java: TypeCheckTest.testTemplateMap5()
func (c *TypeCheckTest) TestTemplateMap5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckTest_EXTERNS_WITH_IARRAYLIKE_DECLS).
		AddSource(`function f() {
  var x = new Int8Array(10);
  /** @type {IArrayLike<number>} */
  var y;
  y = x;
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testTemplateMap6()
func (c *TypeCheckTest) TestTemplateMap6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckTest_EXTERNS_WITH_IARRAYLIKE_DECLS).
		AddSource(`function f() {
  var x = new Int8Array(10);
  /** @type {IObject<number, number>} */
  var z;
  z = x;
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testArrayImplementsIArrayLike()
func (c *TypeCheckTest) TestArrayImplementsIArrayLike() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {!Array<number>} */ var arr = [];
var /** null */ n = arr[0];
`).
		AddDiagnostic(`initializing variable
found   : number
required: null
`).
		Run()
}

// Java: TypeCheckTest.testIArrayLike1()
func (c *TypeCheckTest) TestIArrayLike1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckTest_EXTERNS_WITH_IARRAYLIKE_DECLS).
		AddSource(`var arr = new Int8Array(7);
// no warning
arr[0] = 1;
arr[1] = 2;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testIArrayLike2()
func (c *TypeCheckTest) TestIArrayLike2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckTest_EXTERNS_WITH_IARRAYLIKE_DECLS).
		AddSource(`var arr = new Int8Array(7);
// have warnings
arr[3] = false;
`).
		AddDiagnostic(`assignment
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckTest.testIArrayLike3()
func (c *TypeCheckTest) TestIArrayLike3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckTest_EXTERNS_WITH_IARRAYLIKE_DECLS).
		AddSource(`var arr = new Int8Array2(10);
// have warnings
arr[3] = false;
`).
		AddDiagnostic(`assignment
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckTest.testIArrayLike4()
func (c *TypeCheckTest) TestIArrayLike4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckTest_EXTERNS_WITH_IARRAYLIKE_DECLS).
		AddSource(`var arr = new Int8Array2(10);
// have warnings
arr[3] = false;
`).
		AddDiagnostic(`assignment
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckTest.testIArrayLike5()
func (c *TypeCheckTest) TestIArrayLike5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckTest_EXTERNS_WITH_IARRAYLIKE_DECLS).
		AddSource(`var arr = new Int8Array3(10);
// have warnings
arr[3] = false;
`).
		AddDiagnostic(`assignment
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckTest.testIArrayLike6()
func (c *TypeCheckTest) TestIArrayLike6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckTest_EXTERNS_WITH_IARRAYLIKE_DECLS).
		AddSource(`var arr = new Int8Array4(10);
// have warnings
arr[3] = false;
`).
		AddDiagnostic(`assignment
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckTest.testIArrayLike7()
func (c *TypeCheckTest) TestIArrayLike7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckTest_EXTERNS_WITH_IARRAYLIKE_DECLS).
		AddSource(`var arr5 = new BooleanArray5(10);
arr5[2] = true;
arr5[3] = "";
`).
		AddDiagnostic(`assignment
found   : string
required: boolean
`).
		Run()
}

// Java: TypeCheckTest.testIArrayLike8()
func (c *TypeCheckTest) TestIArrayLike8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckTest_EXTERNS_WITH_IARRAYLIKE_DECLS).
		AddSource(`var arr2 = new Int8Array(10);
arr2[true] = 1;
`).
		AddDiagnostic(`restricted index type
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckTest.testIArrayLike9()
func (c *TypeCheckTest) TestIArrayLike9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckTest_EXTERNS_WITH_IARRAYLIKE_DECLS).
		AddSource(`var arr2 = new Int8Array2(10);
arr2[true] = 1;
`).
		AddDiagnostic(`restricted index type
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckTest.testIArrayLike10()
func (c *TypeCheckTest) TestIArrayLike10() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckTest_EXTERNS_WITH_IARRAYLIKE_DECLS).
		AddSource(`var arr2 = new Int8Array3(10);
arr2[true] = 1;
`).
		AddDiagnostic(`restricted index type
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckTest.testIArrayLike11()
func (c *TypeCheckTest) TestIArrayLike11() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckTest_EXTERNS_WITH_IARRAYLIKE_DECLS).
		AddSource(`var arr2 = new Int8Array4(10);
arr2[true] = 1;
`).
		AddDiagnostic(`restricted index type
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckTest.testIArrayLike12()
func (c *TypeCheckTest) TestIArrayLike12() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckTest_EXTERNS_WITH_IARRAYLIKE_DECLS).
		AddSource(`var arr2 = new BooleanArray5(10);
arr2['prop'] = true;
`).
		AddDiagnostic(`restricted index type
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testIArrayLike13()
func (c *TypeCheckTest) TestIArrayLike13() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckTest_EXTERNS_WITH_IARRAYLIKE_DECLS).
		AddSource(`var numOrStr = null ? 0 : 'prop';
var arr2 = new BooleanArray5(10);
arr2[numOrStr] = true;
`).
		AddDiagnostic(`restricted index type
found   : (number|string)
required: number
`).
		Run()
}

// Java: TypeCheckTest.testIArrayLikeStructuralMatch1()
func (c *TypeCheckTest) TestIArrayLikeStructuralMatch1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** !IArrayLike */ x){};
/** @constructor */
function Foo() {}
/** @type {number} */ Foo.prototype.length
f(new Foo)
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testIArrayLikeStructuralMatch2()
func (c *TypeCheckTest) TestIArrayLikeStructuralMatch2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** !IArrayLike */ x){};
/** @constructor */
function Foo() {
  /** @type {number} */ this.length = 5;
}
f(new Foo)
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testIArrayLikeStructuralMatch3()
func (c *TypeCheckTest) TestIArrayLikeStructuralMatch3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** !IArrayLike */ x){};
f({length: 5})
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testIArrayLikeStructuralMatch4()
func (c *TypeCheckTest) TestIArrayLikeStructuralMatch4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** !IArrayLike */ x){};
/** @const */ var ns = {};
/** @type {number} */ ns.length
f(ns)
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testIArrayLikeStructuralMatch5()
func (c *TypeCheckTest) TestIArrayLikeStructuralMatch5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** !IArrayLike */ x){};
var ns = function() {};
/** @type {number} */ ns.length
f(ns)
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testIArrayLikeStructuralMatch6()
func (c *TypeCheckTest) TestIArrayLikeStructuralMatch6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** !IArrayLike<string> */ x){};
/** @constructor */
function Foo() {}
/** @type {number} */ Foo.prototype.length
f(new Foo)
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testIObject1()
func (c *TypeCheckTest) TestIObject1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckTest_EXTERNS_WITH_IOBJECT_DECLS).
		AddSource(`var arr2 = new Object2();
arr2[0] = 1;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testIObject2()
func (c *TypeCheckTest) TestIObject2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckTest_EXTERNS_WITH_IOBJECT_DECLS).
		AddSource(`var arr2 = new Object2();
arr2['str'] = 1;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testIObject3()
func (c *TypeCheckTest) TestIObject3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckTest_EXTERNS_WITH_IOBJECT_DECLS).
		AddSource(`var arr2 = new Object2();
arr2[true] = 1;
`).
		AddDiagnostic(`restricted index type
found   : boolean
required: (number|string)
`).
		Run()
}

// Java: TypeCheckTest.testIObject4()
func (c *TypeCheckTest) TestIObject4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckTest_EXTERNS_WITH_IOBJECT_DECLS).
		AddSource(`var arr2 = new Object2();
arr2[function(){}] = 1;
`).
		AddDiagnostic(`restricted index type
found   : function(): undefined
required: (number|string)
`).
		Run()
}

// Java: TypeCheckTest.testIObject5()
func (c *TypeCheckTest) TestIObject5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckTest_EXTERNS_WITH_IOBJECT_DECLS).
		AddSource(`var arr2 = new Object2();
arr2[{}] = 1;
`).
		AddDiagnostic(`restricted index type
found   : {}
required: (number|string)
`).
		Run()
}

// Java: TypeCheckTest.testIObject6()
func (c *TypeCheckTest) TestIObject6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckTest_EXTERNS_WITH_IOBJECT_DECLS).
		AddSource(`var arr2 = new Object2();
arr2[undefined] = 1;
`).
		AddDiagnostic(`restricted index type
found   : undefined
required: (number|string)
`).
		Run()
}

// Java: TypeCheckTest.testIObject7()
func (c *TypeCheckTest) TestIObject7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckTest_EXTERNS_WITH_IOBJECT_DECLS).
		AddSource(`var arr2 = new Object2();
arr2[null] = 1;
`).
		AddDiagnostic(`restricted index type
found   : null
required: (number|string)
`).
		Run()
}

// Java: TypeCheckTest.testIObject8()
func (c *TypeCheckTest) TestIObject8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckTest_EXTERNS_WITH_IOBJECT_DECLS).
		AddSource(`var arr = new Object2();
/** @type {boolean} */
var x = arr[3];
`).
		AddDiagnostic(`initializing variable
found   : number
required: boolean
`).
		Run()
}

// Java: TypeCheckTest.testIObject9()
func (c *TypeCheckTest) TestIObject9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckTest_EXTERNS_WITH_IOBJECT_DECLS).
		AddSource(`var arr = new Object2();
/** @type {(number|string)} */
var x = arr[3];
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testIObject10()
func (c *TypeCheckTest) TestIObject10() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckTest_EXTERNS_WITH_IOBJECT_DECLS).
		AddSource(`var arr = new Object3();
/** @type {number} */
var x = arr[3];
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testIObject11()
func (c *TypeCheckTest) TestIObject11() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckTest_EXTERNS_WITH_IOBJECT_DECLS).
		AddSource(`var arr = new Object3();
/** @type {boolean} */
var x = arr[3];
`).
		AddDiagnostic(`initializing variable
found   : number
required: boolean
`).
		Run()
}

// Java: TypeCheckTest.testIObject12()
func (c *TypeCheckTest) TestIObject12() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckTest_EXTERNS_WITH_IOBJECT_DECLS).
		AddSource(`var arr = new Object3();
/** @type {string} */
var x = arr[3];
`).
		AddDiagnostic(`initializing variable
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testIObject13()
func (c *TypeCheckTest) TestIObject13() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckTest_EXTERNS_WITH_IOBJECT_DECLS).
		AddSource(`var arr = new Object3();
arr[3] = false;
`).
		AddDiagnostic(`assignment
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckTest.testIObject14()
func (c *TypeCheckTest) TestIObject14() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckTest_EXTERNS_WITH_IOBJECT_DECLS).
		AddSource(`var arr = new Object3();
arr[3] = 'value';
`).
		AddDiagnostic(`assignment
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testIObject15()
func (c *TypeCheckTest) TestIObject15() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckTest_EXTERNS_WITH_IOBJECT_DECLS).
		AddSource(`function f(/** !Object<string, string> */ x) {}
var /** !IObject<string, string> */ y;
f(y);
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testIObject16()
func (c *TypeCheckTest) TestIObject16() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckTest_EXTERNS_WITH_IOBJECT_DECLS).
		AddSource(`function f(/** !Object<string, string> */ x) {}
var /** !IObject<string, number> */ y;
f(y);
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : IObject<string,number>
required: Object<string,string>
`).
		Run()
}

// Java: TypeCheckTest.testEmptySubtypeRecord1()
func (c *TypeCheckTest) TestEmptySubtypeRecord1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @record */ function I1() {};
/** @type {number|undefined} */ I1.prototype.prop;

/** @record @extends {I1} */ function I2() {}
/** @record @extends {I2} */ function I3() {}
/** @type {string} */ I3.prototype.prop2;
/** @constructor @implements {I3} */ function C() {
    /** @type {number} */ this.prop = 1;
    /** @type {string} */ this.prop2 = 'str';
}

/** @param {I3} a */ function fn(a) {}
`).
		Run()
}

// Java: TypeCheckTest.testEmptySubtypeRecord2()
func (c *TypeCheckTest) TestEmptySubtypeRecord2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {!SubInterface} */ var value;
/** @record */ var SuperInterface = function () {};
/** @record @extends {SuperInterface} */ var SubInterface = function() {};
`).
		Run()
}

// Java: TypeCheckTest.testEmptySubtypeRecord3()
func (c *TypeCheckTest) TestEmptySubtypeRecord3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @record */ var SuperInterface = function () {};
/** @record @extends {SuperInterface} */ var SubInterface = function() {};
/** @type {!SubInterface} */ var value;
`).
		Run()
}

// Java: TypeCheckTest.testEmptySubtypeInterface1()
func (c *TypeCheckTest) TestEmptySubtypeInterface1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {!SubInterface} */ var value;
/** @interface */ var SuperInterface = function () {};
/** @interface @extends {SuperInterface} */ var SubInterface = function() {};
`).
		Run()
}

// Java: TypeCheckTest.testRecordWithOptionalProperty()
func (c *TypeCheckTest) TestRecordWithOptionalProperty() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**  @constructor */ function Foo() {};
Foo.prototype.str = 'foo';

var /** {str: string, opt_num: (undefined|number)} */ x = new Foo;
`).
		Run()
}

// Java: TypeCheckTest.testRecordWithUnknownProperty()
func (c *TypeCheckTest) TestRecordWithUnknownProperty() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**  @constructor */ function Foo() {};
Foo.prototype.str = 'foo';

var /** {str: string, unknown: ?} */ x = new Foo;
`).
		AddDiagnostic(`initializing variable
found   : Foo
required: {
  str: string,
  unknown: ?
}
missing : [unknown]
mismatch: []
`).
		Run()
}

// Java: TypeCheckTest.testRecordWithOptionalUnknownProperty()
func (c *TypeCheckTest) TestRecordWithOptionalUnknownProperty() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**  @constructor */ function Foo() {};
Foo.prototype.str = 'foo';

var /** {str: string, opt_unknown: (?|undefined)} */ x = new Foo;
`).
		Run()
}

// Java: TypeCheckTest.testRecordWithTopProperty()
func (c *TypeCheckTest) TestRecordWithTopProperty() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**  @constructor */ function Foo() {};
Foo.prototype.str = 'foo';

var /** {str: string, top: *} */ x = new Foo;
`).
		AddDiagnostic(`initializing variable
found   : Foo
required: {
  str: string,
  top: *
}
missing : [top]
mismatch: []
`).
		Run()
}

// Java: TypeCheckTest.testOptionalUnknownIsAssignableToUnknown()
func (c *TypeCheckTest) TestOptionalUnknownIsAssignableToUnknown() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** (undefined|?) */ opt_unknown) {
  var /** ? */ unknown = opt_unknown;
}
`).
		Run()
}

// Java: TypeCheckTest.testOptimizePropertyMap1()
func (c *TypeCheckTest) TestOptimizePropertyMap1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {!Function} */
function f() {
  var x = function() {};
  /** @type {number} */
  x.prop = 123;
  return x;
}
function g(/** !Function */ x) {
  var /** null */ n = x.prop;
}
`).
		AddDiagnostic("Property prop never defined on Function").
		Run()
}

// Java: TypeCheckTest.testOptimizePropertyMap1b()
func (c *TypeCheckTest) TestOptimizePropertyMap1b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {!Function} */
function f() {
  var x = function() {};
  /** @type {number} */
  x.prop = 123;
  return x;
}
function g(/** !Function */ x) {
  var /** null */ n = x.prop;
}
`).
		Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		Run()
}

// Java: TypeCheckTest.testOptimizePropertyMap2()
func (c *TypeCheckTest) TestOptimizePropertyMap2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @typedef {{a:number}} */
var Foo;
function f(/** !Foo */ x) {
  var y = x;
  /** @type {number} */
  y.b = 123;
}
function g(/** !Foo */ x) {
  var /** null */ n = x.b;
}
`).
		AddDiagnostic("Property b never defined on x").
		Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		Run()
}

// Java: TypeCheckTest.testOptimizePropertyMap2b()
func (c *TypeCheckTest) TestOptimizePropertyMap2b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @typedef {{a:number}} */
var Foo;
function f(/** !Foo */ x) {
  var y = x;
  /** @type {number} */
  y.b = 123;
}
`).
		AddDiagnostic("Property b never defined on y").
		Run()
}

// Java: TypeCheckTest.testOptimizePropertyMap3a()
func (c *TypeCheckTest) TestOptimizePropertyMap3a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @record */
function Foo() {}
/** @type {number} */
Foo.prototype.a;
function f(/** !Foo */ x) {
  var y = x;
  /** @type {number} */
  y.b = 123;
}
function g(/** !Foo */ x) {
  var /** null */ n = x.b;
}
`).
		AddDiagnostic("Property b never defined on Foo").
		AddDiagnostic("Property b never defined on Foo").
		Run()
}

// Java: TypeCheckTest.testOptimizePropertyMap3b()
func (c *TypeCheckTest) TestOptimizePropertyMap3b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @record */
function Foo() {}
/** @type {number} */
Foo.prototype.a;
function f(/** !Foo */ x) {
  var y = x;
  /** @type {number} */
  y.b = 123;
}
function g(/** !Foo */ x) {
  var /** null */ n = x.b;
}
`).
		Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		Run()
}

// Java: TypeCheckTest.testOptimizePropertyMap4()
func (c *TypeCheckTest) TestOptimizePropertyMap4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(x) {
  var y = { a: 1, b: 2 };
}
function g(x) {
  return x.b + 1;
}
`).
		Run()
}

// Java: TypeCheckTest.testOptimizePropertyMap5()
func (c *TypeCheckTest) TestOptimizePropertyMap5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(x) {
  var y = { a: 1, b: 2 };
}
function g() {
  var x = { c: 123 };
  return x.a + 1;
}
`).
		AddDiagnostic("Property a never defined on x").
		Run()
}

// Java: TypeCheckTest.testOptimizePropertyMap6()
func (c *TypeCheckTest) TestOptimizePropertyMap6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** {a:number} */ x) {
  var y = x;
  /** @type {number} */
  y.b = 123;
}
function g(/** {c:number} */ x) {
  var /** null */ n = x.b;
}
`).
		AddDiagnostic("Property b never defined on x").
		Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		Run()
}

// Java: TypeCheckTest.testOptimizePropertyMap7()
func (c *TypeCheckTest) TestOptimizePropertyMap7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f() {
  var x = {a:1};
  x.b = 2;
}
function g() {
  var y = {a:1};
  return y.b + 1;
}
`).
		AddDiagnostic("Property b never defined on y").
		Run()
}

// Java: TypeCheckTest.testOptimizePropertyMap8()
func (c *TypeCheckTest) TestOptimizePropertyMap8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** {a:number, b:number} */ x) {}
function g(/** {c:number} */ x) {
  var /** null */ n = x.b;
}
`).
		AddDiagnostic("Property b never defined on x").
		Run()
}

// Java: TypeCheckTest.testOptimizePropertyMap9()
func (c *TypeCheckTest) TestOptimizePropertyMap9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */
function Foo() {
  this.a = 123;
}
function f(/** {a: number, c: string} */ x) {
  x.b = 234;
}
function g(/** !Foo */ x) {
  return x.b + 5;
}
`).
		AddDiagnostic("Property b never defined on Foo").
		Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		Run()
}

// Java: TypeCheckTest.testDuplicateVariableDefinition1()
func (c *TypeCheckTest) TestDuplicateVariableDefinition1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @record */
function A() {}
/** @type {number} */
A.prototype.prop;
/** @record */
function B() {}
/** @type {number} */
B.prototype.prop;
/** @constructor */
function C() {}
/** @type {number} */
C.prototype.prop;
/** @return {(A|B|C)} */
function fun () {}
/** @return {(B|A|C)} */
function fun () {}
`).
		AddDiagnostic("variable fun redefined, original definition at testcode0:14").
		Run()
}

// Java: TypeCheckTest.testDuplicateVariableDefinition3()
func (c *TypeCheckTest) TestDuplicateVariableDefinition3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var ns = {};
/** @type {{x:number}} */ ns.x;
/** @type {{x:number}} */ ns.x;
`).
		Run()
}

// Java: TypeCheckTest.testDuplicateVariableDefinition3_1()
func (c *TypeCheckTest) TestDuplicateVariableDefinition3_1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var ns = {};
/** @type {{x:number}} */ ns.x;
/** @type {{x:string}} */ ns.x;
`).
		AddDiagnostic(`variable ns.x redefined with type {x: string}, original definition at testcode0:2 with type {x: number}
`).
		Run()
}

// Java: TypeCheckTest.testDuplicateVariableDefinition3_2()
func (c *TypeCheckTest) TestDuplicateVariableDefinition3_2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var ns = {};
/** @type {{x:number}} */ ns.x;
/** @type {{x:number, y:boolean}} */ ns.x;
`).
		AddDiagnostic(`variable ns.x redefined with type {
  x: number,
  y: boolean
}, original definition at testcode0:2 with type {x: number}
`).
		Run()
}

// Java: TypeCheckTest.testDuplicateVariableDefinition4()
func (c *TypeCheckTest) TestDuplicateVariableDefinition4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var ns = {};
/** @record */ function rec3(){}
/** @record */ function rec4(){}
/** @type {!rec3} */ ns.x;
/** @type {!rec4} */ ns.x;
`).
		AddDiagnostic(`variable ns.x redefined with type rec4, original definition at testcode0:4 with type rec3
`).
		Run()
}

// Java: TypeCheckTest.testDuplicateVariableDefinition5()
func (c *TypeCheckTest) TestDuplicateVariableDefinition5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var ns = {};
/** @record */ function rec3(){}
/** @record */ function rec4(){}
/** @type {number} */ rec4.prototype.prop;
/** @type {!rec3} */ ns.x;
/** @type {!rec4} */ ns.x;
`).
		AddDiagnostic(`variable ns.x redefined with type rec4, original definition at testcode0:5 with type rec3
`).
		Run()
}

// Java: TypeCheckTest.testDuplicateVariableDefinition6()
func (c *TypeCheckTest) TestDuplicateVariableDefinition6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var ns = {};
/** @record */ function rec3(){}
/** @type {number} */ rec3.prototype.prop;
/** @record */ function rec4(){}
/** @type {!rec3} */ ns.x;
/** @type {!rec4} */ ns.x;
`).
		AddDiagnostic(`variable ns.x redefined with type rec4, original definition at testcode0:5 with type rec3
`).
		Run()
}

// Java: TypeCheckTest.testDuplicateVariableDefinition7()
func (c *TypeCheckTest) TestDuplicateVariableDefinition7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @typedef {{prop:TD2}} */
  var TD1;

  /** @typedef {{prop:TD1}} */
  var TD2;

  var /** TD1 */ td1;
  var /** TD2 */ td2;

  td1 = td2;
`).
		Run()
}

// Java: TypeCheckTest.testDuplicateVariableDefinition8()
func (c *TypeCheckTest) TestDuplicateVariableDefinition8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var ns = {}
/** @record */ function rec(){}
/** @type {number} */ rec.prototype.prop;

/** @type {!rec} */ ns.x;
/** @type {{prop:number}} */ ns.x;

/** @type {{prop:number}} */ ns.y;
/** @type {!rec} */ ns.y;
`).
		AddDiagnostic(`variable ns.x redefined with type {prop: number}, original definition at testcode0:5 with type rec
`).
		AddDiagnostic(`variable ns.y redefined with type rec, original definition at testcode0:8 with type {prop: number}
`).
		Run()
}

// Java: TypeCheckTest.testDuplicateVariableDefinition8_2()
func (c *TypeCheckTest) TestDuplicateVariableDefinition8_2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var ns = {}
/** @record */ function rec(){}
/** @type {number} */ rec.prototype.prop;

/** @type {!rec} */ ns.x;
/** @type {{prop:string}} */ ns.x;

/** @type {{prop:number}} */ ns.y;
/** @type {!rec} */ ns.y;
`).
		AddDiagnostic(`variable ns.x redefined with type {prop: string}, original definition at testcode0:5 with type rec
`).
		AddDiagnostic(`variable ns.y redefined with type rec, original definition at testcode0:8 with type {prop: number}
`).
		Run()
}

// Java: TypeCheckTest.testDuplicateVariableDefinition8_3()
func (c *TypeCheckTest) TestDuplicateVariableDefinition8_3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var ns = {}
/** @record */ function rec(){}
/** @type {string} */ rec.prototype.prop;

/** @type {!rec} */ ns.x;
/** @type {{prop:string}} */ ns.x;

/** @type {{prop:number}} */ ns.y;
/** @type {!rec} */ ns.y;
`).
		AddDiagnostic(`variable ns.x redefined with type {prop: string}, original definition at testcode0:5 with type rec
`).
		AddDiagnostic(`variable ns.y redefined with type rec, original definition at testcode0:8 with type {prop: number}
`).
		Run()
}

// Java: TypeCheckTest.testDuplicateVariableDefinition8_4()
func (c *TypeCheckTest) TestDuplicateVariableDefinition8_4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @record @template T */ function I() {}
/** @type {T} */ I.prototype.prop;
var ns = {}
/** @record */ function rec(){}
/** @type {I} */ rec.prototype.prop;

/** @type {!rec} */ ns.x;
/** @type {{prop:I}} */ ns.x;
`).
		AddDiagnostic(`variable ns.x redefined with type {prop: (I|null)}, original definition at testcode0:7 with type rec
`).
		Run()
}

// Java: TypeCheckTest.testDuplicateVariableDefinition8_5()
func (c *TypeCheckTest) TestDuplicateVariableDefinition8_5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @record @template T */ function I() {}
/** @type {T} */ I.prototype.prop;
var ns = {}
/** @record */ function rec(){}
/** @type {I<number>} */ rec.prototype.prop;

/** @type {!rec} */ ns.x;
/** @type {{prop:I<number>}} */ ns.x;
`).
		AddDiagnostic(`variable ns.x redefined with type {prop: (I<number>|null)}, original definition at testcode0:7 with type rec
`).
		Run()
}

// Java: TypeCheckTest.testDuplicateVariableDefinition8_6()
func (c *TypeCheckTest) TestDuplicateVariableDefinition8_6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @record @template T */ function I() {}
/** @type {T} */ I.prototype.prop;
var ns = {}
/** @record */ function rec(){}
/** @type {I<number>} */ rec.prototype.prop;

/** @type {!rec} */ ns.x;
/** @type {{prop:I<string>}} */ ns.x;
`).
		AddDiagnostic(`variable ns.x redefined with type {prop: (I<string>|null)}, original definition at testcode0:7 with type rec
`).
		Run()
}

// Java: TypeCheckTest.testDuplicateVariableDefinition8_7()
func (c *TypeCheckTest) TestDuplicateVariableDefinition8_7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @record @template T */
function rec(){}
/** @type {T} */ rec.prototype.value;

/** @type {rec<string>} */ ns.x;
/** @type {{value: string}} */ ns.x;
`).
		AddDiagnostic(`variable ns.x redefined with type {value: string}, original definition at testcode0:5 with type (null|rec<string>)
`).
		Run()
}

// Java: TypeCheckTest.testModuloNullUndefThatWorkedWithoutSpecialSubtypingRules1()
func (c *TypeCheckTest) TestModuloNullUndefThatWorkedWithoutSpecialSubtypingRules1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */
function Foo() {}
function f(/** function(?Foo, !Foo) */ x) {
  return /** @type {function(!Foo, ?Foo)} */ (x);
}
`).
		Run()
}

// Java: TypeCheckTest.testModuloNullUndefThatWorkedWithoutSpecialSubtypingRules2()
func (c *TypeCheckTest) TestModuloNullUndefThatWorkedWithoutSpecialSubtypingRules2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */
function Foo() {}
function f(/** !Array<!Foo> */ to, /** !Array<?Foo> */ from) {
  to = from;
}
`).
		Run()
}

// Java: TypeCheckTest.testModuloNullUndefThatWorkedWithoutSpecialSubtypingRules3()
func (c *TypeCheckTest) TestModuloNullUndefThatWorkedWithoutSpecialSubtypingRules3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** ?Object */ x) {
  return {} instanceof x;
}
`).
		Run()
}

// Java: TypeCheckTest.testModuloNullUndefThatWorkedWithoutSpecialSubtypingRules4()
func (c *TypeCheckTest) TestModuloNullUndefThatWorkedWithoutSpecialSubtypingRules4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** ?Function */ x) {
  return x();
}
`).
		Run()
}

// Java: TypeCheckTest.testEs5ClassExtendingEs6Class()
func (c *TypeCheckTest) TestEs5ClassExtendingEs6Class() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {}
/** @constructor @extends {Foo} */ var Bar = function() {};
`).
		AddDiagnostic("ES5 class Bar cannot extend ES6 class Foo").
		Run()
}

// Java: TypeCheckTest.testEs5ClassExtendingEs6Class_noWarning()
func (c *TypeCheckTest) TestEs5ClassExtendingEs6Class_noWarning() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class A {}
/** @constructor @extends {A} */
const B = createSubclass(A);
`).
		Run()
}

// Java: TypeCheckTest.testNonNullTemplatedThis()
func (c *TypeCheckTest) TestNonNullTemplatedThis() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */
function C() {}

/**
  @return {THIS}
  @this {THIS}
  @template THIS
*/
C.prototype.method = function() {};

/** @return {C|null} */
function f() {
  return x;
};

/** @type {string} */ var s = f().method();
`).
		AddDiagnostic(`initializing variable
found   : C
required: string
`).
		Run()
}

// Java: TypeCheckTest.testSetPrototypeToNewInstance()
func (c *TypeCheckTest) TestSetPrototypeToNewInstance() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */
function C() {}
C.prototype = new C;
`).
		Run()
}

// Java: TypeCheckTest.testNoResolvedTypeAndGreatestSubtypeInference()
func (c *TypeCheckTest) TestNoResolvedTypeAndGreatestSubtypeInference() {
	c.TestClosureTypes(`goog.forwardDeclare('Foo');
/**
 * @param {boolean} pred
 * @param {?Foo} x
 */
function f(pred, x) {
  var y;
  if (pred) {
    y = null;
  } else {
    y = x;
  }
  var /** number */ z = y;
}
`, jutil.StringPtr(`initializing variable
found   : (Foo|null)
required: number
`))
}

// Java: TypeCheckTest.testNoResolvedTypeDoesntCauseInfiniteLoop()
func (c *TypeCheckTest) TestNoResolvedTypeDoesntCauseInfiniteLoop() {
	c.TestClosureTypes(`goog.forwardDeclare('Foo');
goog.forwardDeclare('Bar');

/** @interface */
var Baz = function() {};
/** @return {!Bar} */
Baz.prototype.getBar = function() {};
/** @constructor */
var Qux = function() {
  /** @type {?Foo} */
  this.jobRuntimeTracker_ = null;
};
/** @param {!Baz} job */
Qux.prototype.runRenderJobs_ = function(job) {
  for (var i = 0; i < 10; i++) {
    if (this.jobRuntimeTracker_) {
      goog.asserts.assert(job.getBar, '');
    }
  }
};
`, nil)
}

// Java: TypeCheckTest.testb38182645()
func (c *TypeCheckTest) Testb38182645() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @interface
 * @template VALUE
 */
function MyI() {}


/**
 * @constructor
 * @implements {MyI<K|V>}
 * @template K, V
 */
function MyMap() {}


/**
 * @param {!MyMap<string, T>} map
 * @return {T}
 * @template T
 */
function getValueFromNameAndMap(map) {
  return /** @type {?} */ (123);
}
var m = /** @type {!MyMap<string,number>} */ (new MyMap());
var /** null */ n = getValueFromNameAndMap(m);
`).
		AddDiagnostic(`initializing variable
found   : number
required: null
`).
		Run()
}

// Java: TypeCheckTest.testLocalType1()
func (c *TypeCheckTest) TestLocalType1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function C(){ /** @const */ this.a = true;}
function f() {
  // shadow
  /** @constructor */ function C(){ /** @const */ this.a = 1;}
}
`).
		Run()
}

// Java: TypeCheckTest.testLocalType2()
func (c *TypeCheckTest) TestLocalType2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function C(){ /** @const */ this.a = true;}
function f() {
  // shadow
  /** @constructor */ function C(){ /** @const */ this.a = 1;}
  /** @type {null} */ var x = new C().a;
}
`).
		AddDiagnostic(`initializing variable
found   : number
required: null
`).
		Run()
}

// Java: TypeCheckTest.testForwardDecl1()
func (c *TypeCheckTest) TestForwardDecl1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const */ var ns = {};
/** @constructor */ ns.Outer = function C() {};
/** @return {!ns.Outer.Inner} */
ns.Outer.prototype.method = function() {
  return new ns.Outer.Inner();
};
/** @constructor */ ns.Outer.Inner = function () {};
`).
		Run()
}

// Java: TypeCheckTest.testForwardDecl2()
func (c *TypeCheckTest) TestForwardDecl2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const */ var ns1 = {};
/** @const */ ns1.ns2 = {};
/** @constructor */ ns1.ns2.Outer = function C() {};
/** @return {!ns1.ns2.Outer.Inner} */
ns1.ns2.Outer.prototype.method = function() {
  return new ns1.ns2.Outer.Inner();
};
/** @constructor */ ns1.ns2.Outer.Inner = function () {};
`).
		Run()
}

// Java: TypeCheckTest.testMissingPropertiesWarningOnObject1()
func (c *TypeCheckTest) TestMissingPropertiesWarningOnObject1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */
function Foo() {
  this.prop = 123;
}
/** @param {!Object} x */
function f(x) {
  return x.prop;
}
`).
		AddDiagnostic("Property prop never defined on Object").
		Run()
}

// Java: TypeCheckTest.testStrictNumericOperators1()
func (c *TypeCheckTest) TestStrictNumericOperators1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var x = 'asdf' - 1;").
		AddDiagnostic(`left operand
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testStrictNumericOperators2()
func (c *TypeCheckTest) TestStrictNumericOperators2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var x = 1 - 'asdf';").
		AddDiagnostic(`right operand
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testStrictNumericOperators3()
func (c *TypeCheckTest) TestStrictNumericOperators3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var x = 'asdf'; x++;").
		AddDiagnostic(`increment/decrement
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testStrictNumericOperators4a()
func (c *TypeCheckTest) TestStrictNumericOperators4a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var x = -'asdf';").
		AddDiagnostic(`sign operator
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testStrictNumericOperators4b()
func (c *TypeCheckTest) TestStrictNumericOperators4b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var x = +'asdf';").
		Run()
}

// Java: TypeCheckTest.testStrictNumericOperators5()
func (c *TypeCheckTest) TestStrictNumericOperators5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var x = 1 < 'asdf';").
		AddDiagnostic(`right side of numeric comparison
found   : string
required: (bigint|number)
`).
		Run()
}

// Java: TypeCheckTest.testStrictNumericOperators6()
func (c *TypeCheckTest) TestStrictNumericOperators6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var x = 'asdf'; x *= 2;").
		AddDiagnostic(`left operand
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testStrictNumericOperators7()
func (c *TypeCheckTest) TestStrictNumericOperators7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var x = ~ 'asdf';").
		AddDiagnostic(`bitwise NOT
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testStrictNumericOperators8()
func (c *TypeCheckTest) TestStrictNumericOperators8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var x = 'asdf' | 1;").
		AddDiagnostic(`bad left operand to bitwise operator
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testStrictNumericOperators9()
func (c *TypeCheckTest) TestStrictNumericOperators9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var x = 'asdf' << 1;").
		AddDiagnostic(`operator <<
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testStrictNumericOperators10()
func (c *TypeCheckTest) TestStrictNumericOperators10() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var x = 1 >>> 'asdf';").
		AddDiagnostic(`operator >>>
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testClassFieldDictError()
func (c *TypeCheckTest) TestClassFieldDictError() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @dict */ class C { x=2; }").
		AddDiagnostic("Illegal key, the class is a dict").
		Run()
}

// Java: TypeCheckTest.testClassFieldStaticDictError()
func (c *TypeCheckTest) TestClassFieldStaticDictError() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @dict */ class C { static x=2; }").
		AddDiagnostic("Illegal key, the class is a dict").
		Run()
}

// Java: TypeCheckTest.testClassFieldUnrestricted()
func (c *TypeCheckTest) TestClassFieldUnrestricted() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @unrestricted */ class C { x = 2; }").
		Run()
}

// Java: TypeCheckTest.testClassFieldStaticUnrestricted()
func (c *TypeCheckTest) TestClassFieldStaticUnrestricted() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @unrestricted */ class C { static x = 2; }").
		Run()
}

// Java: TypeCheckTest.testClassFieldTypeError1()
func (c *TypeCheckTest) TestClassFieldTypeError1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class C {
  /** @type {string} */
  x = 2;
}
`).
		AddDiagnostic(`assignment to property x of C
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testClassFieldStaticTypeError1()
func (c *TypeCheckTest) TestClassFieldStaticTypeError1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class C {
  /** @type {string} */
  static x = 2;
}
`).
		AddDiagnostic(`assignment to property x of C
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testClassFieldTypeError2()
func (c *TypeCheckTest) TestClassFieldTypeError2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class C {
  /** @type {string} */
  x = '';
  constructor() {
    /** @type {number} */
    this.x = 1;
  }
}
`).
		AddDiagnostic(`assignment to property x of C
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testClassFieldTypeError3()
func (c *TypeCheckTest) TestClassFieldTypeError3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`const obj = {};
obj.C = class {
  /** @type {string} */
  x = 2;
}
`).
		AddDiagnostic(`assignment to property x of obj.C
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testClassFieldStaticTypeError3()
func (c *TypeCheckTest) TestClassFieldStaticTypeError3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`const obj = {};
obj.C = class {
  /** @type {string} */
  static x = 2;
}
`).
		AddDiagnostic(`assignment to property x of obj.C
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testClassInitializerTypeInitializer()
func (c *TypeCheckTest) TestClassInitializerTypeInitializer() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {number|undefined} y */
function f(y) {
class C { /** @type {string} */ x = y ?? 0; }
}
`).
		AddDiagnostic(`assignment to property x of C
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testClassInitializerTypeInference()
func (c *TypeCheckTest) TestClassInitializerTypeInference() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {string} s
 * @return {string}
 */
const stringIdentity = (s) => s;
class C { x = stringIdentity(0); }
`).
		AddDiagnostic(`actual parameter 1 of stringIdentity does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testClassComputedField()
func (c *TypeCheckTest) TestClassComputedField() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @dict */ class C { [x]=2; }").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @dict */ class C { 'x' = 2; }").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @dict */ class C { 1 = 2; }").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {number} x */ function takesNum(x) {}
/** @dict */
class C {
  /** @type {string} @suppress {checkTypes} */
  [takesNum('string')] = 2;
}
`).
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @unrestricted */ class C { [x] = 2; }").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @unrestricted */ class C { 'x' = 2; }").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @unrestricted */ class C { 1 = 2; }").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @unrestricted */
class C {
  /** @type {string} */
  [x] = 2;
}
`).
		Run()
}

// Java: TypeCheckTest.testClassComputedFieldStatic()
func (c *TypeCheckTest) TestClassComputedFieldStatic() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @dict */ class C { static [x]=2; }").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @dict */ class C { static 'x' = 2; }").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @dict */ class C { static 1 = 2; }").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @dict */
class C {
  /** @type {string}*/
  static [x] = 2;
}
`).
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @unrestricted */ class C { static [x]=2; }").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @unrestricted */ class C { static 'x' = 2; }").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @unrestricted */ class C { static 1 = 2; }").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @unrestricted */
class C {
  /** @type {string} */
  static [x] = 2;
}
`).
		Run()
}

// Java: TypeCheckTest.testClassComputedFieldNoInitializer()
func (c *TypeCheckTest) TestClassComputedFieldNoInitializer() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @dict */ class C { [x]; }").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @dict */ class C { 'x' }").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @dict */ class C { 1 }").
		Run()
}

// Java: TypeCheckTest.testClassComputedFieldNoInitializerStatic()
func (c *TypeCheckTest) TestClassComputedFieldNoInitializerStatic() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @dict */ class C { static [x]; }").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @dict */ class C { static 'x' }").
		Run()
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @dict */ class C { static 1 }").
		Run()
}

// Java: TypeCheckTest.testClassComputedFieldError()
func (c *TypeCheckTest) TestClassComputedFieldError() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("class C { [x] = 2; }").
		AddDiagnostic("Cannot do '[]' access on a struct").
		Run()
}

// Java: TypeCheckTest.testClassComputedFieldErrorStatic()
func (c *TypeCheckTest) TestClassComputedFieldErrorStatic() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("class C { static [x] = 2; }").
		AddDiagnostic("Cannot do '[]' access on a struct").
		Run()
}

// Java: TypeCheckTest.testStrictComparison1()
func (c *TypeCheckTest) TestStrictComparison1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var x = true < 'asdf';").
		AddDiagnostic(`expected matching types in comparison
found   : string
required: boolean
`).
		Run()
}

// Java: TypeCheckTest.testStrictComparison2()
func (c *TypeCheckTest) TestStrictComparison2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** (number|string) */ x, /** string */ y) {
  return x < y;
}
`).
		Run()
}

// Java: TypeCheckTest.testComparisonWithUnknownAndSymbol()
func (c *TypeCheckTest) TestComparisonWithUnknownAndSymbol() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @type {symbol} */ var x; /** @type {?} */ var y; x < y").
		AddDiagnostic(`left side of comparison
found   : symbol
required: (bigint|number|string)
`).
		Run()
}

// Java: TypeCheckTest.testComparisonInStrictModeNoSpuriousWarning()
func (c *TypeCheckTest) TestComparisonInStrictModeNoSpuriousWarning() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(x) {
  var y = 'asdf';
  var z = y < x;
}
`).
		Run()
}

// Java: TypeCheckTest.testComparisonTreatingUnknownAsNumber()
func (c *TypeCheckTest) TestComparisonTreatingUnknownAsNumber() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_PRIMITIVE_OPERATORS).
		AddSource(`/** @constructor */
function Foo() {}
function f(/** ? */ x) {
  return (new Foo) < x;
}
`).
		AddDiagnostic(`left side of comparison
found   : Foo
required: (bigint|number|string)
`).
		Run()
}

// Java: TypeCheckTest.testComparisonTreatingUnknownAsNumber_leftTypeUnknown()
func (c *TypeCheckTest) TestComparisonTreatingUnknownAsNumber_leftTypeUnknown() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_PRIMITIVE_OPERATORS).
		AddSource(`/** @constructor */
function Bar() {}
function f(/** ? */ x) {
  return x < (new Bar);
}
`).
		AddDiagnostic(`right side of comparison
found   : Bar
required: (bigint|number|string)
`).
		Run()
}

// Java: TypeCheckTest.testEnumOfSymbol1()
func (c *TypeCheckTest) TestEnumOfSymbol1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum {symbol} */
var ES = {A: Symbol('a'), B: Symbol('b')};

/** @type {!Object<ES, number>} */
var o = {};

o[ES.A] = 1;
`).
		Run()
}

// Java: TypeCheckTest.testEnumOfSymbol2()
func (c *TypeCheckTest) TestEnumOfSymbol2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum {symbol} */
var ES = {A: Symbol('a'), B: Symbol('b')};

/** @type {!Object<number, number>} */
var o = {};

o[ES.A] = 1;
`).
		AddDiagnostic(`restricted index type
found   : ES<symbol>
required: number
`).
		Run()
}

// Java: TypeCheckTest.testEnumOfSymbol4()
func (c *TypeCheckTest) TestEnumOfSymbol4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum {symbol} */
var ES = {A: Symbol('a'), B: Symbol('b')};

/** @const */
var o = {};

o[ES.A] = 1;
`).
		Run()
}

// Java: TypeCheckTest.testSymbol1()
func (c *TypeCheckTest) TestSymbol1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const */
var o = {};

if (o[Symbol.iterator]) { /** ok */ };
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testSymbol2()
func (c *TypeCheckTest) TestSymbol2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const */
var o = new Symbol();
`).
		AddDiagnostic("cannot instantiate non-constructor, found type: (typeof Symbol)").
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testUnknownExtends1()
func (c *TypeCheckTest) TestUnknownExtends1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @template T
 * @param {function(new: T)} clazz
 */
function f(clazz) {
  /**
   * @constructor
   * @extends {clazz}
   */
  function Foo() {}
}
`).
		Run()
}

// Java: TypeCheckTest.testUnknownExtends2()
func (c *TypeCheckTest) TestUnknownExtends2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** function(new: ?) */ clazz) {
  /**
   * @constructor
   * @extends {clazz}
   */
  function Foo() {}
}
`).
		AddDiagnostic("Could not resolve type in @extends tag of Foo").
		Run()
}

// Java: TypeCheckTest.testMixinApplication1()
func (c *TypeCheckTest) TestMixinApplication1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(TypeCheckTest_MIXIN_DEFINITIONS + `/**
 * @constructor
 * @extends {MyElement}
 * @implements {Toggle}
 */
var MyElementWithToggle = addToggle(MyElement);
(new MyElementWithToggle).foobar(123)
`).
		AddDiagnostic(`actual parameter 1 of MyElementWithToggle.prototype.foobar does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testMixinApplication1_inTypeSummaryWithVar()
func (c *TypeCheckTest) TestMixinApplication1_inTypeSummaryWithVar() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @typeSummary */
` + TypeCheckTest_MIXIN_DEFINITIONS + `/**
 * @constructor
 * @extends {MyElement}
 * @implements {Toggle}
 */
var MyElementWithToggle;
`).
		AddSource(`class SubToggle extends MyElementWithToggle {}
(new SubToggle).foobar(123);
`).
		AddDiagnostic(`actual parameter 1 of MyElementWithToggle.prototype.foobar does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testMixinApplication1_inTypeSummaryWithLet()
func (c *TypeCheckTest) TestMixinApplication1_inTypeSummaryWithLet() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @typeSummary */
` + TypeCheckTest_MIXIN_DEFINITIONS + `/**
 * @constructor
 * @extends {MyElement}
 * @implements {Toggle}
 */
let MyElementWithToggle;
`).
		AddSource(`class SubToggle extends MyElementWithToggle {}
(new SubToggle).foobar(123);
`).
		AddDiagnostic(`actual parameter 1 of MyElementWithToggle.prototype.foobar does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testMixinApplication2()
func (c *TypeCheckTest) TestMixinApplication2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(TypeCheckTest_MIXIN_DEFINITIONS + `/**
 * @constructor
 * @extends {MyElement}
 * @implements {Toggle}
 */
var MyElementWithToggle = addToggle(MyElement);
(new MyElementWithToggle).elemprop = 123;
`).
		AddDiagnostic(`assignment to property elemprop of MyElementWithToggle
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testMixinApplication3()
func (c *TypeCheckTest) TestMixinApplication3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(TypeCheckTest_MIXIN_DEFINITIONS + `var MyElementWithToggle = addToggle(MyElement);
/** @type {MyElementWithToggle} */ var x = 123;
`).
		AddDiagnostic(`Bad type annotation. Unknown type MyElementWithToggle
It's possible that 'MyElementWithToggle' refers to a value, not a type.
`).
		Run()
}

// Java: TypeCheckTest.testReassignedSymbolAccessedInClosureIsFlowInsensitive()
func (c *TypeCheckTest) TestReassignedSymbolAccessedInClosureIsFlowInsensitive() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddString().Build()).
		AddSource(`function f(){
  for (var x = {}; ; x = {y: x.y}) {
    x.y = 'str';
    x.y = x.y.length
    var g = (function(x) {
      return function() {
        if (x.y) -x.y;
      };
    })(x);
  }
}
`).
		AddDiagnostic(`sign operator
found   : (number|string)
required: number
`).
		Run()
}

// Java: TypeCheckTest.testMethodOverriddenDirectlyOnThis()
func (c *TypeCheckTest) TestMethodOverriddenDirectlyOnThis() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */
function Foo() {
  this.bar = function() { return 'str'; };
}
/** @return {number} */
Foo.prototype.bar = function() {};
`).
		AddDiagnostic(`inconsistent return type
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTest.testSuperclassDefinedInBlockUsingVar()
func (c *TypeCheckTest) TestSuperclassDefinedInBlockUsingVar() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`{
  /** @constructor */
  var Base = function() {};
  /** @param {number} x */
  Base.prototype.baz = function(x) {};
}
/** @constructor @extends {Base} */
var Foo = function() {};
/** @override */
Foo.prototype.baz = function(x) {
  var /** string */ y = x;
};
`).
		AddDiagnostic(`initializing variable
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testSuperclassDefinedInBlockOnNamespace()
func (c *TypeCheckTest) TestSuperclassDefinedInBlockOnNamespace() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const */
var ns = {};
{
  /** @constructor */
  ns.Base = function() {};
  /** @param {number} x */
  ns.Base.prototype.baz = function(x) {};
}
/** @constructor @extends {ns.Base} */
ns.Foo = function() {};
/** @override */
ns.Foo.prototype.baz = function(x) {
  var /** string */ y = x;
};
`).
		AddDiagnostic(`initializing variable
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testPropertyReferenceOnShadowingParameter()
func (c *TypeCheckTest) TestPropertyReferenceOnShadowingParameter() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const */
var ns = {};
ns.fn = function(/** number */ a){};
class Namespace {
  fn(/** string */ a){}
}
function test(/** !Namespace */ ns) { // The parameter 'ns' shadows the global 'ns'.
// Verify that ns.fn resolves to Namespace.prototype.fn, not the global ns.fn.
  ns.fn(0);
}
`).
		AddDiagnostic(`actual parameter 1 of Namespace.prototype.fn does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testPropertyReferenceOnShadowingParameter_whenGlobalVariableIsUndeclared()
func (c *TypeCheckTest) TestPropertyReferenceOnShadowingParameter_whenGlobalVariableIsUndeclared() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().IncludeDefaultExterns().
		AddExterns("ns.fn = function(/** number */ a){};").
		AddSource(`class Namespace {
  fn(/** string */ a){}
}
function test(/** !Namespace */ ns) { // The parameter 'ns' shadows the global 'ns'.
// Verify that ns.fn resolves to Namespace.prototype.fn, not the global ns.fn.
  ns.fn(0);
}
`).
		AddDiagnostic(`actual parameter 1 of Namespace.prototype.fn does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testPropertyReferenceOnShadowingVariable()
func (c *TypeCheckTest) TestPropertyReferenceOnShadowingVariable() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const */
var ns = {};
ns.fn = function(/** number */ a){};
class Namespace {
  fn(/** string */ a){}
}
function test() {
// The local 'ns' shadows the global 'ns'.
  const /** !Namespace */ ns = /** @type {!Namespace} */ (unknown);
// Verify that ns.fn resolves to Namespace.prototype.fn, not the global ns.fn.
  ns.fn(0);
}
`).
		AddDiagnostic(`actual parameter 1 of Namespace.prototype.fn does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testPropertyReferenceOnShadowingVariable_inferredLocalSlotOverridesGlobal()
func (c *TypeCheckTest) TestPropertyReferenceOnShadowingVariable_inferredLocalSlotOverridesGlobal() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const */
var ns = {};
ns.fn = function(/** number */ a){};
/** @unrestricted */
class Namespace {}
function test(/** !Namespace */ ns) { // The parameter 'ns' shadows the global 'ns'.
  ns.fn = function(/** string */ a) {};
// Verify that ns.fn resolves to the locally assigned ns.fn, not the global ns.fn.
  ns.fn(0);
}
`).
		AddDiagnostic("Property fn never defined on Namespace").
		AddDiagnostic(`actual parameter 1 of ns.fn does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testReassigningPropertyOnClassTypeInFunction_localSlotOverridesPrototype()
func (c *TypeCheckTest) TestReassigningPropertyOnClassTypeInFunction_localSlotOverridesPrototype() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @unrestricted */\nclass Namespace {\n  fn(/** string */ a) {}\n}\nfunction test(/** !Namespace */ ns) {\n// Assign ns.fn to a narrow type than Namespace.prototype.fn.\n  ns.fn = function(/** string|number */ a) {};\n// Verify that ns.fn resolves to the narrower local type, which can accept `number`.\n  ns.fn(0);\n}\n").
		Run()
}

// Java: TypeCheckTest.testPropertyReferenceOnShadowingThisProperty()
func (c *TypeCheckTest) TestPropertyReferenceOnShadowingThisProperty() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class C {
  constructor() {
    /** @type {?number} */
    this.size = null;
  }
  method() {
// TODO(b/132283774): Redeclaring 'this.size' should cause an error.
    /** @type {number} */
    this.size = unknown;
    const /** number */ num = this.size;
  }
}
`).
		Run()
}

// Java: TypeCheckTest.testThisReferenceTighteningInInnerScope()
func (c *TypeCheckTest) TestThisReferenceTighteningInInnerScope() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("class C {\n  constructor() {\n    /** @type {?number} */\n    this.size = 0;\n// Verify that the inferred type of `this.size`, which is `number`, propagates into\n// block scopes.\n    { const /** number */ n = this.size; }\n  }\n}\n").
		Run()
}

// Java: TypeCheckTest.testShadowedForwardReferenceHoisted()
func (c *TypeCheckTest) TestShadowedForwardReferenceHoisted() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ var C = function() { /** @type {string} */ this.prop = 's'; };
var fn = function() {
  /** @type {C} */ var x = f();
  /** @type {number} */ var n1 = x.prop;
  /** @type {number} */ var n2 = new C().prop;
  /** @constructor */ function C() { /** @type {number} */ this.prop = 1; };
  /** @return {C} */ function fn() { return new C(); };

}
`).
		Run()
}

// Java: TypeCheckTest.testTypeDeclarationsShadowOneAnotherWithFunctionScoping()
func (c *TypeCheckTest) TestTypeDeclarationsShadowOneAnotherWithFunctionScoping() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */
var C = function() { }; // (1)

var fn = function() {
  /** @type {?C} */ var innerC; // (2)

  /** @constructor */
  var C = function() { }; // (3)

  innerC = new C(); // (4)
}
`).
		Run()
}

// Java: TypeCheckTest.testTypeDeclarationsShadowOneAnotherWithModuleScoping_withTemplate()
func (c *TypeCheckTest) TestTypeDeclarationsShadowOneAnotherWithModuleScoping_withTemplate() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().IncludeDefaultExterns().
		AddExterns("class C {}").
		AddSource(`/** @type {!C} */ let innerC; // (2)

// NB: without the @template this test would pass due to b/110741413.
/** @template T */
class C {}; // (3)

innerC = new C(); // (4)
export {}
`).
		Run()
}

// Java: TypeCheckTest.testTypeDeclarationsShadowOneAnotherWithModuleScoping_withSuperclass()
func (c *TypeCheckTest) TestTypeDeclarationsShadowOneAnotherWithModuleScoping_withSuperclass() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().IncludeDefaultExterns().
		AddExterns("class C {}").
		AddSource(`/** @type {!C} */ let innerC; // (2)

class Parent {}
class C extends Parent {}; // (3)

/** @type {!Parent} */ const p = innerC; // (4)
export {}
`).
		Run()
}

// Java: TypeCheckTest.testTypeDeclarationsShadowOneAnotherWithFunctionScopingConsideringHoisting()
func (c *TypeCheckTest) TestTypeDeclarationsShadowOneAnotherWithFunctionScopingConsideringHoisting() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */
var C = function() { }; // (1)

var fn = function() {
  /** @constructor */
  var C = function() { }; // (2)

// This function will be hoisted above, and therefore type-analysed before, (2).
  /** @return {!C} */ // (3)
  function hoisted() { return new C(); }; // (4)
}
`).
		Run()
}

// Java: TypeCheckTest.testRefinedTypeInNestedShortCircuitingAndOr()
func (c *TypeCheckTest) TestRefinedTypeInNestedShortCircuitingAndOr() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function B() {}
/** @constructor @extends {B} */ function C() {}
/** @return {string} */ C.prototype.foo = function() {};
/** @return {boolean} */ C.prototype.bar = function() {};
/** @constructor @extends {B} */ function D() {}
/** @return {number} */ D.prototype.foo = function() {};
/** @param {!B} arg
    @return {null} */
function f(arg) {
  if ((arg instanceof C && arg.bar()) || arg instanceof D) {
    return arg.foo();
  }
  return null;
}
`).
		AddDiagnostic(`inconsistent return type
found   : (number|string)
required: null
`).
		Run()
}

// Java: TypeCheckTest.testTypeofType_notInScope()
func (c *TypeCheckTest) TestTypeofType_notInScope() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var /** typeof ns */ x;").
		AddDiagnostic("Missing type for `typeof` value. The value must be declared and const.").
		Run()
}

// Java: TypeCheckTest.testTypeofType_namespace()
func (c *TypeCheckTest) TestTypeofType_namespace() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const */ var ns = {};
var /** typeof ns */ x = ns;
`).
		Run()
}

// Java: TypeCheckTest.testTypeofType_namespaceMismatch()
func (c *TypeCheckTest) TestTypeofType_namespaceMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const */ var ns = {};
var /** typeof ns */ x = {};
`).
		AddDiagnostic(`initializing variable
found   : {}
required: {}
`).
		Run()
}

// Java: TypeCheckTest.testTypeofType_namespaceForwardReferenceMismatch()
func (c *TypeCheckTest) TestTypeofType_namespaceForwardReferenceMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var /** typeof ns */ x = {};
/** @const */ var ns = {};
`).
		AddDiagnostic(`initializing variable
found   : {}
required: {}
`).
		Run()
}

// Java: TypeCheckTest.testTypeofType_constructor()
func (c *TypeCheckTest) TestTypeofType_constructor() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddArray().Build()).
		AddSource(`/** @constructor */ function Foo() {}
var /** !Array<typeof Foo> */ x = [];
x.push(Foo);
`).
		Run()
}

// Java: TypeCheckTest.testTypeofType_constructorMismatch1()
func (c *TypeCheckTest) TestTypeofType_constructorMismatch1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddArray().Build()).
		AddSource(`/** @constructor */ function Foo() {}
var /** !Array<(typeof Foo)> */ x = [];
x.push(new Foo());
`).
		AddDiagnostic(`actual parameter 1 of Array.prototype.push does not match formal parameter
found   : Foo
required: (typeof Foo)
`).
		Run()
}

// Java: TypeCheckTest.testTypeofType_constructorMismatch2()
func (c *TypeCheckTest) TestTypeofType_constructorMismatch2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddArray().Build()).
		AddSource(`/** @constructor */ function Foo() {}
var /** !Array<!Foo> */ x = [];
var /** typeof Foo */ y = Foo;
x.push(y);
`).
		AddDiagnostic(`actual parameter 1 of Array.prototype.push does not match formal parameter
found   : (typeof Foo)
required: Foo
`).
		Run()
}

// Java: TypeCheckTest.testTypeofType_enum()
func (c *TypeCheckTest) TestTypeofType_enum() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum */ var Foo = {A: 1}
var /** typeof Foo */ x = Foo;
`).
		Run()
}

// Java: TypeCheckTest.testTypeofType_enumFromCall()
func (c *TypeCheckTest) TestTypeofType_enumFromCall() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum */ var Foo = {A: 1}
/** @return {*} */ function getFoo() { return Foo; }
var x = /** @type {typeof Foo} */ (getFoo());
`).
		Run()
}

// Java: TypeCheckTest.testTypeofType_enumMismatch1()
func (c *TypeCheckTest) TestTypeofType_enumMismatch1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum */ var Foo = {A: 1}
var /** typeof Foo */ x = Foo.A;
`).
		AddDiagnostic(`initializing variable
found   : Foo<number>
required: enum{Foo}
`).
		Run()
}

// Java: TypeCheckTest.testTypeofType_enumMismatch2()
func (c *TypeCheckTest) TestTypeofType_enumMismatch2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum */ var Foo = {A: 1}
/** @enum */ var Bar = {A: 1}
var /** typeof Foo */ x = Bar;
`).
		AddDiagnostic(`initializing variable
found   : enum{Bar}
required: enum{Foo}
`).
		Run()
}

// Java: TypeCheckTest.testTypeofType_castNamespaceIncludesPropertiesFromTypeofType()
func (c *TypeCheckTest) TestTypeofType_castNamespaceIncludesPropertiesFromTypeofType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const */ var ns1 = {};
/** @type {string} */ ns1.foo;
/** @const */ var ns2 = {};
/** @type {number} */ ns2.bar;

/** @const {typeof ns2} */ var ns = /** @type {?} */ (ns1);
/** @type {null} */ var x = ns.bar;
`).
		AddDiagnostic(`initializing variable
found   : number
required: null
`).
		Run()
}

// Java: TypeCheckTest.testTypeofType_castNamespaceDoesNotIncludeOwnProperties()
func (c *TypeCheckTest) TestTypeofType_castNamespaceDoesNotIncludeOwnProperties() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const */ var ns1 = {};
/** @type {string} */ ns1.foo;
/** @const */ var ns2 = {};
/** @type {number} */ ns2.bar;

/** @const {typeof ns2} */ var ns = /** @type {?} */ (ns1);
/** @type {null} */ var x = ns.foo;
`).
		AddDiagnostic("Property foo never defined on ns").
		Run()
}

// Java: TypeCheckTest.testTypeofType_namespaceTypeIsAnAliasNotACopy()
func (c *TypeCheckTest) TestTypeofType_namespaceTypeIsAnAliasNotACopy() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const */ var ns1 = {};
/** @type {string} */ ns1.foo;

/** @const {typeof ns1} */ var ns = /** @type {?} */ (x);

/** @type {string} */ ns1.bar;

/** @type {null} */ var x = ns.bar;
`).
		AddDiagnostic(`initializing variable
found   : string
required: null
`).
		Run()
}

// Java: TypeCheckTest.testTypeofType_namespacedTypeNameResolves()
func (c *TypeCheckTest) TestTypeofType_namespacedTypeNameResolves() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const */ var ns1 = {};
/** @constructor */ ns1.Foo = function() {};
/** @const {typeof ns1} */ var ns = /** @type {?} */ (x);
/** @type {!ns.Foo} */ var x = null;
`).
		AddDiagnostic(`initializing variable
found   : null
required: ns1.Foo
`).
		Run()
}

// Java: TypeCheckTest.testTypeofType_unknownType()
func (c *TypeCheckTest) TestTypeofType_unknownType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var /** ? */ x;
/** @type {typeof x} */ var y;
`).
		AddDiagnostic("Missing type for `typeof` value. The value must be declared and const.").
		Run()
}

// Java: TypeCheckTest.testTypeofType_namespacedTypeOnIndirectAccessNameResolves()
func (c *TypeCheckTest) TestTypeofType_namespacedTypeOnIndirectAccessNameResolves() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @const */ var ns1 = {};\n/** @constructor */ ns1.Foo = function() {};\nconst ns2 = ns1;\n// Verify this works although we have never seen any assignments to `ns2.Foo`\n/** @const {typeof ns2.Foo} */ var Foo2 = /** @type {?} */ (x);\n/** @type {null} */ var x = new Foo2();\n").
		AddDiagnostic(`initializing variable
found   : ns1.Foo
required: null
`).
		Run()
}

// Java: TypeCheckTest.testTypeofType_namespacedTypeMissing()
func (c *TypeCheckTest) TestTypeofType_namespacedTypeMissing() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const */ var ns = {};
/** @const {typeof ns.Foo} */ var Foo = /** @type {?} */ (x);
`).
		AddDiagnostic("Missing type for `typeof` value. The value must be declared and const.").
		Run()
}

// Java: TypeCheckTest.testTypeofType_withGlobalDeclaredVariableWithHoistedFunction()
func (c *TypeCheckTest) TestTypeofType_withGlobalDeclaredVariableWithHoistedFunction() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {string|number} */
var x = 1;
function g(/** typeof x */ a) {}
x = 'str';
g(null);
`).
		AddDiagnostic(`actual parameter 1 of g does not match formal parameter
found   : null
required: (number|string)
`).
		Run()
}

// Java: TypeCheckTest.testTypeofType_typeofForwardReferencedShadowingLocal()
func (c *TypeCheckTest) TestTypeofType_typeofForwardReferencedShadowingLocal() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const {string} */
var x = 'x';
function f() {
  function g(/** typeof x */ a) {}
  /** @const {number} */
  var x = 1;
  g(null);
}
`).
		AddDiagnostic(`actual parameter 1 of g does not match formal parameter
found   : null
required: number
`).
		Run()
}

// Java: TypeCheckTest.testTypeofType_withGlobalDeclaredVariableWithFunction()
func (c *TypeCheckTest) TestTypeofType_withGlobalDeclaredVariableWithFunction() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {string|number} */
var x = 1;
let g = function(/** typeof x */ a) {}
x = 'str';
g(null);
`).
		AddDiagnostic(`actual parameter 1 of g does not match formal parameter
found   : null
required: (number|string)
`).
		Run()
}

// Java: TypeCheckTest.testTypeofType_withGlobalInferredVariableWithFunctionExpression()
func (c *TypeCheckTest) TestTypeofType_withGlobalInferredVariableWithFunctionExpression() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var x = 1;
var g = function (/** typeof x */ a) {}
x = 'str';
g(null);
`).
		AddDiagnostic("Missing type for `typeof` value. The value must be declared and const.").
		Run()
}

// Java: TypeCheckTest.testTypeofType_withLocalInferredVariableInHoistedFunction()
func (c *TypeCheckTest) TestTypeofType_withLocalInferredVariableInHoistedFunction() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f() {
  var x = 1;
  function g(/** typeof x */ a) {}
  x = 'str';
  g(null);
}
`).
		AddDiagnostic("Missing type for `typeof` value. The value must be declared and const.").
		Run()
}

// Java: TypeCheckTest.testTypeofType_withLocalDeclaredVariableWithHoistedFunction()
func (c *TypeCheckTest) TestTypeofType_withLocalDeclaredVariableWithHoistedFunction() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f() {
  /** @type {string|number} */
  var x = 1;
  function g(/** typeof x */ a) {}
  x = 'str';
  g(null);
}
`).
		AddDiagnostic(`actual parameter 1 of g does not match formal parameter
found   : null
required: (number|string)
`).
		Run()
}

// Java: TypeCheckTest.testTypeofType_withLocalInferredVariableWithFunctionExpression()
func (c *TypeCheckTest) TestTypeofType_withLocalInferredVariableWithFunctionExpression() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f() {
  var x = 1;
  var g = function (/** typeof x */ a) {}
  x = 'str';
  g(null);
}
`).
		AddDiagnostic("Missing type for `typeof` value. The value must be declared and const.").
		Run()
}

// Java: TypeCheckTest.testTypeofType_withLocalDeclaredVariableWithFunctionExpression()
func (c *TypeCheckTest) TestTypeofType_withLocalDeclaredVariableWithFunctionExpression() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f() {
  /** @type {string|number} */
  var x = 1;
  var g = function (/** typeof x */ a) {}
  x = 'str';
  g(null);
}
`).
		AddDiagnostic(`actual parameter 1 of g does not match formal parameter
found   : null
required: (number|string)
`).
		Run()
}

// Java: TypeCheckTest.testTypeofType_withLocalInferredVariable()
func (c *TypeCheckTest) TestTypeofType_withLocalInferredVariable() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f() {
  var x = 1;
  /** @type {typeof x} */
  var g = null
  x = 'str';
}
`).
		AddDiagnostic("Missing type for `typeof` value. The value must be declared and const.").
		Run()
}

// Java: TypeCheckTest.testTypeofType_withLocalDeclaredVariable1()
func (c *TypeCheckTest) TestTypeofType_withLocalDeclaredVariable1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f() {
  /** @type {string|number} */
  var x = 1;
  /** @type {typeof x} */
  var g = null
  x = 'str';
}
`).
		AddDiagnostic(`initializing variable
found   : null
required: (number|string)
`).
		Run()
}

// Java: TypeCheckTest.testTypeofType_withLocalDeclaredVariableAfterReassignment()
func (c *TypeCheckTest) TestTypeofType_withLocalDeclaredVariableAfterReassignment() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f() {
  /** @type {string|number} */
  var x = 1;
  x = 'str';
  /** @type {typeof x} */
  var g = null
}
`).
		AddDiagnostic(`initializing variable
found   : null
required: (number|string)
`).
		Run()
}

// Java: TypeCheckTest.testTypeofType_withLocalDeclaredVariableAfterTightening()
func (c *TypeCheckTest) TestTypeofType_withLocalDeclaredVariableAfterTightening() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f() {
  /** @type {string|number} */
  var x = 'str';
  if (typeof x == 'string') {
    /** @type {typeof x} */
    let g = null
  }
}
`).
		AddDiagnostic(`initializing variable
found   : null
required: (number|string)
`).
		Run()
}

// Java: TypeCheckTest.testTypeNameAliasOnAliasedNamespace()
func (c *TypeCheckTest) TestTypeNameAliasOnAliasedNamespace() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {};
/** @enum {number} */ Foo.E = {A: 1};
const F = Foo;
const E = F.E;
/** @type {E} */ let e = undefined;
`).
		AddDiagnostic(`initializing variable
found   : undefined
required: Foo.E<number>
`).
		Run()
}

// Java: TypeCheckTest.testAssignOp()
func (c *TypeCheckTest) TestAssignOp() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function fn(someUnknown) {
  var x = someUnknown;
  x *= 2; // infer 'x' to now be 'number'
  var /** null */ y = x;
}
`).
		AddDiagnostic(`initializing variable
found   : number
required: null
`).
		Run()
}

// Java: TypeCheckTest.testAnnotatedObjectLiteralInBlocklessArrow1()
func (c *TypeCheckTest) TestAnnotatedObjectLiteralInBlocklessArrow1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** {g: function(number): undefined} */ x) {}
() => f({/** @param {string} x */ g(x) {}});
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : {g: function(string): undefined}
required: {g: function(number): undefined}
missing : []
mismatch: [g]
`).
		Run()
}

// Java: TypeCheckTest.testAnnotatedObjectLiteralInBlocklessArrow2()
func (c *TypeCheckTest) TestAnnotatedObjectLiteralInBlocklessArrow2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** {g: function(): string} */ x) {}
() => f({/** @constructor */ g: function() {}});
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : {g: (typeof <anonymous@testcode0:2>)}
required: {g: function(): string}
missing : []
mismatch: [g]
`).
		Run()
}

// Java: TypeCheckTest.testAnnotatedObjectLiteralInBlocklessArrow3()
func (c *TypeCheckTest) TestAnnotatedObjectLiteralInBlocklessArrow3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function G() {}
function f(/** {g: function(): string} */ x) {}
() => f({/** @constructor */ g: G});
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : {g: (typeof G)}
required: {g: function(): string}
missing : []
mismatch: [g]
`).
		Run()
}

// Java: TypeCheckTest.testNativePromiseTypeWithExterns()
func (c *TypeCheckTest) TestNativePromiseTypeWithExterns() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var p = new Promise(function(resolve, reject) { resolve(3); });
p.then(result => {}, error => {});
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testNativeThenableType()
func (c *TypeCheckTest) TestNativeThenableType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var /** Thenable */ t = null;").
		AddDiagnostic(`initializing variable
found   : null
required: {then: ?}
`).
		Run()
}

// Java: TypeCheckTest.testFunctionOnObjLitWithAngleBracketKey()
func (c *TypeCheckTest) TestFunctionOnObjLitWithAngleBracketKey() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const {!Object<function()>} */
var x = {'<': function() {}};
`).
		Run()
}

// Java: TypeCheckTest.testFunctionOnObjLitWithAngleBracketKeyChecked()
func (c *TypeCheckTest) TestFunctionOnObjLitWithAngleBracketKeyChecked() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const {!Object<function(): number>} */
var x = {'<': function() { return 'str'; }};
`).
		Run()
}

// Java: TypeCheckTest.testFunctionOnObjLitWithAngleBracketKeyAndJsdoc()
func (c *TypeCheckTest) TestFunctionOnObjLitWithAngleBracketKeyAndJsdoc() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const */
var x = {
  /** @param {number} arg */
  '<': function(arg) {},
};
`).
		Run()
}

// Java: TypeCheckTest.testClassOnObjLitWithAngleBracketKey()
func (c *TypeCheckTest) TestClassOnObjLitWithAngleBracketKey() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const */
var x = {
  /** @constructor */
  '<': function() {},
};
`).
		Run()
}

// Java: TypeCheckTest.testQuotedPropertyWithDotInType()
func (c *TypeCheckTest) TestQuotedPropertyWithDotInType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const */
var x = {
  /** @constructor */
  'y.A': function() {},
};
var /** x.y.A */ a;
`).
		AddDiagnostic(`Bad type annotation. Unknown type x.y.A
It's possible that 'x.y.A' refers to a value, not a type.
`).
		Run()
}

// Java: TypeCheckTest.testQuotedPropertyWithDotInCode()
func (c *TypeCheckTest) TestQuotedPropertyWithDotInCode() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const */
var x = {
  /** @constructor */
  'y.A': function() {},
};
var a = new x.y.A();
`).
		AddDiagnostic("Property y never defined on x").
		AddDiagnostic("Property A never defined on x.y" + TypeCheck_POSSIBLE_INEXISTENT_PROPERTY_EXPLANATION).
		Run()
}

// Java: TypeCheckTest.testClassTemplateParamsInMethodBody()
func (c *TypeCheckTest) TestClassTemplateParamsInMethodBody() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor @template T */
function Foo() {}
Foo.prototype.foo = function() {
  var /** T */ x;
};
`).
		Run()
}

// Java: TypeCheckTest.testClassTtlParamsInMethodBody()
func (c *TypeCheckTest) TestClassTtlParamsInMethodBody() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor @template T := 'number' =: */
function Foo() {}
Foo.prototype.foo = function() {
  var /** T */ x;
};
`).
		AddDiagnostic("Template type transformation T not allowed on classes or interfaces").
		AddDiagnostic("Bad type annotation. Unknown type T").
		Run()
}

// Java: TypeCheckTest.testClassTtlParamsInMethodSignature()
func (c *TypeCheckTest) TestClassTtlParamsInMethodSignature() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor @template T := 'number' =: */
function Foo() {}
/** @return {T} */
Foo.prototype.foo = function() {};
`).
		AddDiagnostic("Template type transformation T not allowed on classes or interfaces").
		AddDiagnostic("Bad type annotation. Unknown type T").
		Run()
}

// Java: TypeCheckTest.testFunctionTemplateParamsInBody()
func (c *TypeCheckTest) TestFunctionTemplateParamsInBody() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @template T */
function f(/** !Array<T> */ arr) {
  var /** T */ first = arr[0];
}
`).
		Run()
}

// Java: TypeCheckTest.testFunctionTtlParamsInBody()
func (c *TypeCheckTest) TestFunctionTtlParamsInBody() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @template T := 'number' =: */
function f(/** !Array<T> */ arr) {
  var /** T */ first = arr[0];
}
`).
		Run()
}

// Java: TypeCheckTest.testFunctionTemplateParamsInBodyMismatch()
func (c *TypeCheckTest) TestFunctionTemplateParamsInBodyMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @template T, V */
function f(/** !Array<T> */ ts, /** !Array<V> */ vs) {
  var /** T */ t = vs[0];
}
`).
		Run()
}

// Java: TypeCheckTest.testClassAndMethodTemplateParamsInMethodBody()
func (c *TypeCheckTest) TestClassAndMethodTemplateParamsInMethodBody() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor @template T */
function Foo() {}
/** @template U */
Foo.prototype.foo = function() {
  var /** T */ x;
  var /** U */ y;
};
`).
		Run()
}

// Java: TypeCheckTest.testNestedFunctionTemplates()
func (c *TypeCheckTest) TestNestedFunctionTemplates() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor @template A, B */
function Foo(/** A */ a, /** B */ b) {}
/** @template T */
function f(/** T */ t) {
  /** @template S */
  function g(/** S */ s) {
    var /** !Foo<T, S> */ foo = new Foo(t, s);
  }
}
`).
		Run()
}

// Java: TypeCheckTest.testNestedFunctionTemplatesMismatch()
func (c *TypeCheckTest) TestNestedFunctionTemplatesMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor @template A, B */
function Foo(/** A */ a, /** B */ b) {}
/** @template T */
function f(/** T */ t) {
  /** @template S */
  function g(/** S */ s) {
    var /** !Foo<T, S> */ foo = new Foo(s, t);
  }
}
`).
		Run()
}

// Java: TypeCheckTest.testTypeCheckingDoesntCrashOnDebuggerStatement()
func (c *TypeCheckTest) TestTypeCheckingDoesntCrashOnDebuggerStatement() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var x = 1; debugger; x = 2;").
		Run()
}

// Java: TypeCheckTest.testCastOnLhsOfAssignBlocksBadAssignmentWarning()
func (c *TypeCheckTest) TestCastOnLhsOfAssignBlocksBadAssignmentWarning() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var /** number */ x = 1; (/** @type {?} */ (x)) = 'foobar';").
		Run()
}

// Java: TypeCheckTest.testCastOnLhsWithAdditionalJSDoc()
func (c *TypeCheckTest) TestCastOnLhsWithAdditionalJSDoc() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {
  takesString(/** string */ stringParameter) {}
}
/** @param {*} all */
function f(all) {
  /** some jsdoc */
  (/** @type {!Foo} */ (all)).takesString(0);
}
`).
		Run()
}

// Java: TypeCheckTest.testInvalidComparisonsInStrictOperatorsMode()
func (c *TypeCheckTest) TestInvalidComparisonsInStrictOperatorsMode() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("function f(/** (void|string) */ x, /** void */ y) { return x < y; }").
		AddDiagnostic(`right side of comparison
found   : undefined
required: string
`).
		Run()
}

// Java: TypeCheckTest.testStrayTypedefOnRecordInstanceDoesNotModifyRecord()
func (c *TypeCheckTest) TestStrayTypedefOnRecordInstanceDoesNotModifyRecord() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @record */
function Component() {}

/** @const {!Component} */
var MyComponent = {};
/** @typedef {string} */
MyComponent.MyString;

/** @const {!Component} */
var OtherComponent = {};
`).
		Run()
}

// Java: TypeCheckTest.testStrayPropertyOnRecordInstanceDoesNotModifyRecord()
func (c *TypeCheckTest) TestStrayPropertyOnRecordInstanceDoesNotModifyRecord() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @record */
function Component() {}

/** @const {!Component} */
var MyComponent = {};
/** @const {string} */
MyComponent.NAME = 'MyComponent'; // strict inexistent property is usually suppressed

/** @const {!Component} */
var OtherComponent = {};
`).
		AddDiagnosticType(TypeCheck_STRICT_INEXISTENT_PROPERTY).
		Run()
}

// Java: TypeCheckTest.testStrayTypedefOnRecordTypedefDoesNotModifyRecord()
func (c *TypeCheckTest) TestStrayTypedefOnRecordTypedefDoesNotModifyRecord() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @typedef {{foo: number}} */
var Component;

/** @const {!Component} */
var MyComponent = {foo: 1};
/** @typedef {string} */
MyComponent.MyString;

/** @const {!Component} */
var OtherComponent = {foo: 2};
`).
		Run()
}

// Java: TypeCheckTest.testCanUseConstToAssignValueToTypedef()
func (c *TypeCheckTest) TestCanUseConstToAssignValueToTypedef() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @typedef {!Object<string, number>} */
const Type = {};
`).
		Run()
}

// Java: TypeCheckTest.testLetTypedefWithAssignedValue_reassigned()
func (c *TypeCheckTest) TestLetTypedefWithAssignedValue_reassigned() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @typedef {!Object<string, number>} */
let Type = {};
Type = {x: 3};
`).
		AddDiagnostic(`initializing variable
found   : {}
required: None
`).
		AddDiagnostic(`assignment
found   : {x: number}
required: None
`).
		Run()
}

// Java: TypeCheckTest.testLetTypedefWithAssignedValue_notReassigned()
func (c *TypeCheckTest) TestLetTypedefWithAssignedValue_notReassigned() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @typedef {!Object<string, number>} */
let Type = {};
`).
		Run()
}

// Java: TypeCheckTest.testStrayPropertyOnRecordTypedefDoesNotModifyRecord()
func (c *TypeCheckTest) TestStrayPropertyOnRecordTypedefDoesNotModifyRecord() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @typedef {{foo: number}} */
var Component;

/** @const {!Component} */
var MyComponent = {foo: 1 };
/** @const {string} */
MyComponent.NAME = 'MyComponent';

/** @const {!Component} */
var OtherComponent = {foo: 2};
`).
		AddDiagnosticType(TypeCheck_STRICT_INEXISTENT_PROPERTY).
		Run()
}

// Java: TypeCheckTest.testCheckRecursiveTypedefSubclassOfNominalClass()
func (c *TypeCheckTest) TestCheckRecursiveTypedefSubclassOfNominalClass() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @typedef {{self: !Bar}} */
var Foo;
/** @typedef {!Foo} */
var Bar;

/** @type {!Foo} */
var foo;

class Baz {
  constructor() {
    /** @type {!Baz} */ this.self = this;
  }
}

const /** !Baz */ x = foo;
`).
		AddDiagnostic(`initializing variable
found   : {self: {...}}
required: Baz
`).
		Run()
}

// Java: TypeCheckTest.testCheckNominalClassSubclassOfRecursiveTypedef()
func (c *TypeCheckTest) TestCheckNominalClassSubclassOfRecursiveTypedef() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @typedef {{self: !Bar}} */
var Foo;
/** @typedef {!Foo} */
var Bar;

class Baz {
  constructor() {
    /** @type {!Baz} */ this.self = this;
  }
}

const /** !Foo */ x = new Baz();
`).
		Run()
}

// Java: TypeCheckTest.testNullableToStringSubclassOfRecordLiteral()
func (c *TypeCheckTest) TestNullableToStringSubclassOfRecordLiteral() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddString().AddObject().
		Build()).
		AddSource("/** @struct @constructor */\nfunction MyClass() {}\n\n/**\n * @override\n * @return {?string}\n * @suppress {checkTypes} // J2CL allows this, so we need to test it.\n */\nMyClass.prototype.toString = function() {};\n\n// To do a strucutural match, this type needs at least one additional property.\n// There is no empty `{}` type.\n/** @return {number} */\nMyClass.prototype.x = function() {};\n\nvar /** {x: !Function} */ instance = new MyClass();\n").
		Run()
}

// Java: TypeCheckTest.testForwardReferencedRecordTypeUnionedWithConcreteType()
func (c *TypeCheckTest) TestForwardReferencedRecordTypeUnionedWithConcreteType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Concrete {}

/**
 * @param {!ForwardReferenced|!Concrete} arg
 */
var bar = (arg) => {};

/**
 * @typedef {{baz}}
 */
var ForwardReferenced;

/**
 * @param {!Concrete} arg
 */
var foo = (arg) => bar(arg);
`).
		Run()
}

// Java: TypeCheckTest.testUnionCollapse_recordWithOnlyOptionalPropertiesDoesNotSupercedeArrayOrFunction()
func (c *TypeCheckTest) TestUnionCollapse_recordWithOnlyOptionalPropertiesDoesNotSupercedeArrayOrFunction() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @typedef {{x: (string|undefined)}} */
var Options;
/** @constructor */
function Foo() {}
/** @typedef {function(new: Foo)|!Array<function(new: Foo)>} */
var FooCtor;
/** @type {FooCtor|Options} */
var x;
/** @type {null} */
var y = x;
`).
		AddDiagnostic(`initializing variable
found   : (Array<function(new:Foo): ?>|function(new:Foo): ?|{x: (string|undefined)})
required: null
`).
		Run()
}

// Java: TypeCheckTest.testCyclicUnionTypedefs()
func (c *TypeCheckTest) TestCyclicUnionTypedefs() {
	// Go-only: test-helper — Java asserts that running this test throws StackOverflowError. A
	// goroutine stack overflow is a fatal runtime error in Go, which cannot be recovered, so the
	// assertion cannot run.
	c.T().Skip("asserts a java.lang.StackOverflowError, which a Go program cannot recover from")
	// TODO(b/112964849): This case should not throw anything.
	builder := TypeCheckTestCase_TypeTestBuilder_NewTest().
		AddSource(`/** @typedef {Foo|string} */
var Bar;
/** @typedef {Bar|number} */
var Foo;
var /** Foo */ foo;
var /** null */ x = foo;
`).
		AddDiagnostic(`initializing variable
found   : (number|string)
required: null
`)
	AssertThrows[jutil.Throwable](builder.Run)
}

// Java: TypeCheckTest.testAtRecordOnVarObjectLiteral_warns()
func (c *TypeCheckTest) TestAtRecordOnVarObjectLiteral_warns() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @record */ var X = {};").
		AddDiagnostic(`initializing variable
found   : {}
required: function(this:X): ?
`).
		Run()
}

// Java: TypeCheckTest.testAtRecordOnConstObjectLiteral_warns()
func (c *TypeCheckTest) TestAtRecordOnConstObjectLiteral_warns() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @record */ const X = {};").
		AddDiagnostic(`initializing variable
found   : {}
required: function(this:X): ?
`).
		Run()
}

// Java: TypeCheckTest.testAtConstructorOnConstObjectLiteral_warns()
func (c *TypeCheckTest) TestAtConstructorOnConstObjectLiteral_warns() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("/** @constructor */ const X = {};").
		AddDiagnostic(`initializing variable
found   : {}
required: function(new:X): ?
`).
		Run()
}

// Java: TypeCheckTest.testDeclarationAnnotatedEnum_warnsIfAssignedNonEnumRhs()
func (c *TypeCheckTest) TestDeclarationAnnotatedEnum_warnsIfAssignedNonEnumRhs() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("const Y = {A: 0, B: 1}; /** @enum {number} */ const X = Y;").
		AddDiagnostic(`initializing variable
found   : {
  A: number,
  B: number
}
required: enum{X}
`).
		Run()
}

// Java: TypeCheckTest.testEs6ExtendCannotUseGoogInherits()
func (c *TypeCheckTest) TestEs6ExtendCannotUseGoogInherits() {
	c.TestClosureTypes(`class Super {}
/** @extends {Super} */
class Sub {}
goog.inherits(Sub, Super);
`, jutil.StringPtr("Do not use goog.inherits with ES6 classes. Use the ES6 `extends` keyword to inherit instead.\n"))
}

// Java: TypeCheckTest.testParameterShadowsNamespace()
func (c *TypeCheckTest) TestParameterShadowsNamespace() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor @param {!service.Service} service */
var Controller = function(service) {
  /** @private {!service.Service} */
  this.service = service;
};
/** @const */
var service = {};
/** @constructor */
service.Service = function() {};
`).
		AddDiagnostic(`Bad type annotation. Unknown type service.Service
It's possible that a local variable called 'service' is shadowing the intended global namespace.
`).
		Run()
}

// Java: TypeCheckTest.testValueUsedAsType()
func (c *TypeCheckTest) TestValueUsedAsType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {?} */ var Foo = 'Foo';
/** @type {!Foo} */ var foo = 'foo';
`).
		AddDiagnostic(`Bad type annotation. Unknown type Foo
It's possible that 'Foo' refers to a value, not a type.
`).
		Run()
}

// Java: TypeCheckTest.testBangOperatorOnForwardReferencedType()
func (c *TypeCheckTest) TestBangOperatorOnForwardReferencedType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @typedef {?number} */
var Foo;
/** @return {!Foo} */
function f() {}
/** @const {!Foo} */
var x = f();
`).
		Run()
}

// Java: TypeCheckTest.testBangOperatorOnForwardReferencedType_mismatch()
func (c *TypeCheckTest) TestBangOperatorOnForwardReferencedType_mismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {!Foo} */
var x;
/** @typedef {?number} */
var Foo;
/** @type {?Foo} */
var y;
x = y;
`).
		AddDiagnostic(`assignment
found   : (null|number)
required: number
`).
		Run()
}

// Java: TypeCheckTest.testBangOperatorOnIndirectlyForwardReferencedType()
func (c *TypeCheckTest) TestBangOperatorOnIndirectlyForwardReferencedType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {!b.Foo} */
function f() {}
/** @const {!b.Foo} */
var x = f();

const a = {};
/** @typedef {?number} */
a.Foo;
const b = a;
`).
		Run()
}

// Java: TypeCheckTest.testBangOperatorOnIndirectlyForwardReferencedType_mismatch()
func (c *TypeCheckTest) TestBangOperatorOnIndirectlyForwardReferencedType_mismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {?b.Foo} */
function f() {}
/** @const {!b.Foo} */
var x = f();

const a = {};
/** @typedef {?number} */
a.Foo;
const b = a;
`).
		AddDiagnostic(`initializing variable
found   : (null|number)
required: number
`).
		Run()
}

// Java: TypeCheckTest.testBangOperatorOnTypedefForShadowedNamespace()
func (c *TypeCheckTest) TestBangOperatorOnTypedefForShadowedNamespace() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @typedef {!service.Service} */
var serviceService;
/** @constructor @param {!service.Service} service */
var Controller = function(service) {
  /** @private {!serviceService} */
  this.service = service;
};
/** @const */
var service = {};
/** @constructor */
service.Service = function() {};
`).
		Run()
}

// Java: TypeCheckTest.testNullCheckAvoidsInexistentPropertyError()
func (c *TypeCheckTest) TestNullCheckAvoidsInexistentPropertyError() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`function use(x) {}
/** @constructor */ function Foo() {}
var /** !Foo */ foo = new Foo();
if (foo.bar != null) use(foo);
`).
		Run()
}

// Java: TypeCheckTest.testTripleEqualsNonNullCheckGivesInexistentPropertyError()
func (c *TypeCheckTest) TestTripleEqualsNonNullCheckGivesInexistentPropertyError() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		AddSource(`function use(x) {}
/** @constructor */ function Foo() {}
var /** !Foo */ foo = new Foo();
if (foo.bar !== null) use(foo);
`).
		AddDiagnostic("Property bar never defined on Foo").
		Run()
}

// Java: TypeCheckTest.testGenericConstructorCrash()
func (c *TypeCheckTest) TestGenericConstructorCrash() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @template T */
class Base {
 /**
  * @param {(function(function(new: T, ...?)): T)=} instantiate
  * @return {!Array<T>}
  */
 delegates(instantiate = undefined) {
   return [instantiate];
 };
}
class Bar {}
class Foo {}

/** @type {Base<Foo|Bar>} */
const c = new Base();

c.delegates(ctor => new ctor(42));
`).
		Run()
}

// Java: TypeCheckTest.testBug515079328_genericInTypedef()
func (c *TypeCheckTest) TestBug515079328_genericInTypedef() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().IncludeDefaultExterns().
		AddSourceWithFilename("internal.js", `goog.module('internal');
/**
 * @record
 * @template T
 */
function TypeGuard() {}
exports.TypeGuard = TypeGuard;
`).
		AddSourceWithFilename("guards.js", `goog.module('guards');
const rt = goog.requireType('internal');
/** @typedef {!rt.TypeGuard} */
exports.TypeGuard;
`).
		AddSourceWithFilename("test.js", `goog.module('test');
const {TypeGuard} = goog.require('guards');
const guards = goog.require('guards');
/** @type {!TypeGuard<string>} */
const destructuringType = 1;

/** @type {!guards.TypeGuard<null>} */
const nsImportType = 1;
`).
		AddDiagnostic(`initializing variable
found   : number
required: TypeGuard<string>
`).
		AddDiagnostic(`initializing variable
found   : number
required: TypeGuard<null>
`).
		Run()
}

// Java: TypeCheckTest.testLegacyGoogModuleExportTypecheckedAgainstGlobal_simpleModuleId()
func (c *TypeCheckTest) TestLegacyGoogModuleExportTypecheckedAgainstGlobal_simpleModuleId() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @const @suppress {duplicate} */
var globalNs = {};
`).
		AddSource(`goog.module('globalNs');
goog.module.declareLegacyNamespace();
exports = 0;
`).
		AddDiagnostic(`legacy goog.module export
found   : number
required: {}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testLegacyGoogModuleExportTypecheckedAgainstGlobal_simpleModuleIdWithProperties()
func (c *TypeCheckTest) TestLegacyGoogModuleExportTypecheckedAgainstGlobal_simpleModuleIdWithProperties() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @const @suppress {duplicate} */
var globalNs = {};
`).
		AddSource(`goog.module('globalNs');
goog.module.declareLegacyNamespace();
exports.prop = 0;
`).
		AddDiagnostic(`legacy goog.module export
found   : {prop: number}
required: {}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckTest.testLegacyGoogModuleExportTypecheckedAgainstGlobal_dottedModuleId()
func (c *TypeCheckTest) TestLegacyGoogModuleExportTypecheckedAgainstGlobal_dottedModuleId() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @const @suppress {duplicate} */
var globalNs = {};
/** @suppress {duplicate} */
globalNs.Ctor = class {};
`).
		IncludeDefaultExterns().
		AddSource(`goog.module('globalNs.Ctor');
goog.module.declareLegacyNamespace();
class Ctor {}
exports = Ctor;
`).
		AddDiagnostic(`assignment to property Ctor of globalNs
found   : (typeof Ctor)
required: (typeof globalNs.Ctor)
`).
		Run()
}

// Java: TypeCheckTest.testDynamicImportSpecifier()
func (c *TypeCheckTest) TestDynamicImportSpecifier() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var foo = undefined; import(foo);").
		AddDiagnostic(`dynamic import specifier
found   : undefined
required: string
`).
		Run()
}

// Java: TypeCheckTest.testDynamicImport1()
func (c *TypeCheckTest) TestDynamicImport1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroup_ForType(ModuleLoader_INVALID_MODULE_PATH)).
		AddSource("/** @type {number} */ var foo = import('foo.js');").
		AddDiagnostic(`initializing variable
found   : Promise<?>
required: number
`).
		Run()
}

// Java: TypeCheckTest.testDynamicImport2()
func (c *TypeCheckTest) TestDynamicImport2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroup_ForType(ModuleLoader_INVALID_MODULE_PATH)).
		AddSource("/** @type {Promise} */ var foo2 = import('foo.js');").
		Run()
}

// Java: TypeCheckTest.testDynamicImport3()
func (c *TypeCheckTest) TestDynamicImport3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroup_ForType(ModuleLoader_INVALID_MODULE_PATH)).
		AddSource("/** @type {Promise<{default: number}>} */ var foo = import('foo.js');").
		Run()
}

// Java: TypeCheckTest.testConflictingGetterSetterType()
func (c *TypeCheckTest) TestConflictingGetterSetterType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class C {
  /** @return {string} */
  get value() { }

  /** @param {number} v */
  set value(v) { }
}
`).
		AddDiagnosticType(TypeCheck_CONFLICTING_GETTER_SETTER_TYPE).
		Run()
}

// Java: TypeCheckTest.testConflictingGetterSetterTypeSuppressed()
func (c *TypeCheckTest) TestConflictingGetterSetterTypeSuppressed() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @suppress {checkTypes} */
class C {
  /** @return {string} */
  get value() { }

  /** @param {number} v */
  set value(v) { }
}
`).
		Run()
}

// Java: TypeCheckTest.testSuppressCheckTypesOnStatements()
func (c *TypeCheckTest) TestSuppressCheckTypesOnStatements() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {string} s */
function takesString(s) {}

takesString(0); // verify this is an error
/** @suppress {checkTypes} */
takesString(null);
/** @suppress {checkTypes} */
(0, takesString(undefined));
`).
		AddDiagnostic(`actual parameter 1 of takesString does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTest.testSymbolIteratorMethod()
func (c *TypeCheckTest) TestSymbolIteratorMethod() {
	c.T().Skip("@Ignore")
	TypeCheckTestCase_TypeTestBuilder_NewTest().IncludeDefaultExterns().
		AddSource("const /** number */ num = [][Symbol.iterator];").
		AddDiagnosticType(TypeValidator_TYPE_MISMATCH_WARNING).
		Run()
}

// Java: TypeCheckTest.testGoogRequireDynamic_missingSources()
func (c *TypeCheckTest) TestGoogRequireDynamic_missingSources() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().Suppress(DiagnosticGroups_MISSING_SOURCES_WARNINGS).
		AddSource(`class Test {
 /**
  * @return {!Object}
  *
  * @suppress {checkTypes}
  */
testGoogRequireDynamicStubbedAndWithLoadedModule() {
   goog.setImportHandlerInternalDoNotCallOrElse(() => Promise.resolve(null));
   goog.setUncompiledChunkIdHandlerInternalDoNotCallOrElse(s => s);

   goog.loadModule('goog.module("a.loaded.module"); exports.foo = 12;');

   return goog.requireDynamic('a.loaded.module')
       .then(({foo}) => assertEquals(foo, 12));
 } }
`).
		Run()
}

// Java: TypeCheckTest.testClassMultipleExtends_fromClosureJs()
func (c *TypeCheckTest) TestClassMultipleExtends_fromClosureJs() {
	c.NewTestLegacy().AddSourceWithFilename("testcode0.closure.js", `/** @fileoverview @suppress {checkTypes} */
class Foo {}
class Bar {}
/**
 * @extends {Foo}
 * @extends {Bar}
 */
class FooBar {}
`).
		Run()
	var fooBar rhino.FunctionType = c.compiler.GetTopScope().GetVar("FooBar").
		GetType().
		AssertFunctionType()
	AssertThat(fooBar.IsAmbiguousConstructor()).IsTrue()
}

// Java: TypeCheckTest.testReassignClassPrototype()
func (c *TypeCheckTest) TestReassignClassPrototype() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {}
Foo.prototype = {};
`).
		AddDiagnosticType(TypeInference_REASSIGN_CLASS_PROTOTYPE).
		DiagnosticsAreErrors().
		Run()
}

// Java: TypeCheckTest.testReassignClassPrototypeObjectType()
func (c *TypeCheckTest) TestReassignClassPrototypeObjectType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`class Foo {}
const x = /** @type {!Object} */ ({});
Foo.prototype = x;
`).
		AddDiagnosticType(TypeInference_REASSIGN_CLASS_PROTOTYPE).
		DiagnosticsAreErrors().
		Run()
}

// Java: TypeCheckTest.testClosureTypes(String, String)
func (c *TypeCheckTest) TestClosureTypes(js string, description *string) {
	var descriptions []string
	if description != nil {
		descriptions = []string{jutil.StringTrim(*description)}
	}
	c.TestClosureTypesMultipleWarnings(js, descriptions)
}

// Java: TypeCheckTest.testClosureTypesMultipleWarnings(String, List<String>)
func (c *TypeCheckTest) TestClosureTypesMultipleWarnings(js string, descriptions []string) {
	c.compiler.InitOptions(c.compiler.GetOptions())
	var jsRoot *rhino.Node = rhino.IR_Root(c.compiler.ParseTestCode(js))
	var externs *rhino.Node = rhino.IR_Root(c.compiler.ParseTestCode(NewTestExternsBuilder().AddString().AddClosureExterns().
		AddExtra(CompilerTypeTestCase_CLOSURE_DEFS).
		Build()))
	rhino.IR_Root(externs, jsRoot)
	AssertWithMessage("parsing error: " + typeCheckTests_join(c.compiler.GetErrors(), ", ")).That(c.compiler.GetErrorCount()).
		IsEqualTo(0)
	NewGatherModuleMetadata(c.compiler, false, ModuleLoader_ResolutionMode_BROWSER).Process(externs, jsRoot)
	NewProcessClosurePrimitives(c.compiler).Process(externs, jsRoot)
	NewTypeCheck(c.compiler, NewClosureReverseAbstractInterpreter(c.registry).Append(NewSemanticReverseAbstractInterpreter(c.registry)).GetFirst(), c.registry).ProcessForTesting(externs, jsRoot)
	AssertWithMessage("unexpected error(s) : " + typeCheckTests_join(c.compiler.GetErrors(), ", ")).That(c.compiler.GetErrorCount()).
		IsEqualTo(0)
	if descriptions == nil {
		AssertWithMessage("unexpected warning(s) : " + typeCheckTests_join(c.compiler.GetWarnings(), ", ")).That(c.compiler.GetWarningCount()).
			IsEqualTo(0)
	} else {
		AssertWithMessage("unexpected warning(s) : " + typeCheckTests_join(c.compiler.GetWarnings(), ", ")).That(c.compiler.GetWarningCount()).
			IsEqualTo(len(descriptions))
		actualWarningDescriptions := jutil.NewJavaHashSet[string]()
		for i := 0; i < len(descriptions); i++ {
			actualWarningDescriptions.Add(c.compiler.GetWarnings()[i].Description())
		}
		AssertThat(actualWarningDescriptions).IsEqualTo(jutil.NewJavaHashSetFromSlice(descriptions))
	}
}
