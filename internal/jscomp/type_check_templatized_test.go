// Port of com/google/javascript/jscomp/TypeCheckTemplatizedTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/TypeCheckTemplatizedTest.java)

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
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"testing"
)

// TypeCheckTemplatizedTest tests the type checker.
type TypeCheckTemplatizedTest struct {
	TypeCheckTestCaseBase
}

// Java: TypeCheckTemplatizedTest()
func NewTypeCheckTemplatizedTest() *TypeCheckTemplatizedTest {
	c := &TypeCheckTemplatizedTest{}
	c.InitTypeCheckTestCaseBase(c)
	return c
}

func TestTypeCheckTemplatizedTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewTypeCheckTemplatizedTest() })
}

// Java: TypeCheckTemplatizedTest.testTemplatizedArray1()
func (c *TypeCheckTemplatizedTest) TestTemplatizedArray1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {!Array<number>} a
* @return {string}
*/ var f = function(a) { return a[0]; };
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testTemplatizedArray2()
func (c *TypeCheckTemplatizedTest) TestTemplatizedArray2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {!Array<!Array<number>>} a
* @return {number}
*/ var f = function(a) { return a[0]; };
`).
		AddDiagnostic(`inconsistent return type
found   : Array<number>
required: number
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testTemplatizedArray3()
func (c *TypeCheckTemplatizedTest) TestTemplatizedArray3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {!Array<number>} a
* @return {number}
*/ var f = function(a) { a[1] = 0; return a[0]; };
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testTemplatizedArray4()
func (c *TypeCheckTemplatizedTest) TestTemplatizedArray4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {!Array<number>} a
*/ var f = function(a) { a[0] = 'a'; };
`).
		AddDiagnostic(`assignment
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testTemplatizedArray5()
func (c *TypeCheckTemplatizedTest) TestTemplatizedArray5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {!Array<*>} a
*/ var f = function(a) { a[0] = 'a'; };
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testTemplatizedArray6()
func (c *TypeCheckTemplatizedTest) TestTemplatizedArray6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {!Array<*>} a
* @return {string}
*/ var f = function(a) { return a[0]; };
`).
		AddDiagnostic(`inconsistent return type
found   : *
required: string
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testTemplatizedArray7()
func (c *TypeCheckTemplatizedTest) TestTemplatizedArray7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {?Array<number>} a
* @return {string}
*/ var f = function(a) { return a[0]; };
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testTemplatizedObject1()
func (c *TypeCheckTemplatizedTest) TestTemplatizedObject1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {!Object<number>} a
* @return {string}
*/ var f = function(a) { return a[0]; };
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testTemplatizedObjectOnWindow()
func (c *TypeCheckTemplatizedTest) TestTemplatizedObjectOnWindow() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns("/** @constructor */ window.Object = Object;").
		AddSource(`/** @param {!window.Object<number>} a
 *  @return {string}
 */ var f = function(a) { return a[0]; };
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testTemplatizedObjectOnWindow2()
func (c *TypeCheckTemplatizedTest) TestTemplatizedObjectOnWindow2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns("/** @const */ window.Object = Object;").
		AddSource(`/** @param {!window.Object<number>} a
 *  @return {string}
 */ var f = function(a) { return a[0]; };
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testTemplatizedObject2()
func (c *TypeCheckTemplatizedTest) TestTemplatizedObject2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {!Object<string,number>} a
* @return {string}
*/ var f = function(a) { return a['x']; };
`).
		AddDiagnostic(`inconsistent return type
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testTemplatizedObject3()
func (c *TypeCheckTemplatizedTest) TestTemplatizedObject3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {!Object<number,string>} a
* @return {string}
*/ var f = function(a) { return a['x']; };
`).
		AddDiagnostic(`restricted index type
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testTemplatizedObject4()
func (c *TypeCheckTemplatizedTest) TestTemplatizedObject4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @enum {string} */ var E = {A: 'a', B: 'b'};
/** @param {!Object<E,string>} a
* @return {string}
*/ var f = function(a) { return a['x']; };
`).
		AddDiagnostic(`restricted index type
found   : string
required: E<string>
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testTemplatizedObject5()
func (c *TypeCheckTemplatizedTest) TestTemplatizedObject5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function F() {
  /** @type {Object<number, string>} */ this.numbers = {};
}
(new F()).numbers['ten'] = '10';
`).
		AddDiagnostic(`restricted index type
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testTemplatized1()
func (c *TypeCheckTemplatizedTest) TestTemplatized1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {!Array<string>} */
var arr1 = [];
/** @type {!Array<number>} */
var arr2 = [];
arr1 = arr2;
`).
		AddDiagnostic(`assignment
found   : Array<number>
required: Array<string>
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testTemplatized2()
func (c *TypeCheckTemplatizedTest) TestTemplatized2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {!Array<string>} */
var arr1 = /** @type {!Array<number>} */([]);
`).
		AddDiagnostic(`initializing variable
found   : Array<number>
required: Array<string>
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testTemplatized3()
func (c *TypeCheckTemplatizedTest) TestTemplatized3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {Array<string>} */
var arr1 = /** @type {!Array<number>} */([]);
`).
		AddDiagnostic(`initializing variable
found   : Array<number>
required: (Array<string>|null)
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testTemplatized4()
func (c *TypeCheckTemplatizedTest) TestTemplatized4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {Array<string>} */
var arr1 = [];
/** @type {Array<number>} */
var arr2 = arr1;
`).
		AddDiagnostic(`initializing variable
found   : (Array<string>|null)
required: (Array<number>|null)
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testTemplatized5()
func (c *TypeCheckTemplatizedTest) TestTemplatized5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {Object<T>} obj
 * @return {boolean|undefined}
 * @template T
 */
var some = function(obj) {
  for (var key in obj) if (obj[key]) return true;
};
/** @return {!Array} */ function f() { return []; }
/** @return {!Array<string>} */ function g() { return []; }
some(f());
some(g());
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testTemplatized6()
func (c *TypeCheckTemplatizedTest) TestTemplatized6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ function I(){}
/** @param {T} a
 * @return {T}
 * @template T
*/
I.prototype.method;

/** @constructor 
 * @implements {I}
 */ function C(){}
/** @override*/ C.prototype.method = function(a) {}

/** @type {null} */ var some = new C().method('str');
`).
		AddDiagnostic(`initializing variable
found   : string
required: null
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testTemplatized7()
func (c *TypeCheckTemplatizedTest) TestTemplatized7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface
 *  @template Q
 */ function I(){}
/** @param {T} a
 * @return {T|Q}
 * @template T
*/
I.prototype.method;
/** @constructor 
 * @implements {I<number>}
 */ function C(){}
/** @override*/ C.prototype.method = function(a) {}
/** @type {null} */ var some = new C().method('str');
`).
		AddDiagnostic(`initializing variable
found   : (number|string)
required: null
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testTemplatized8()
func (c *TypeCheckTemplatizedTest) TestTemplatized8() {
	c.T().Skip("@Ignore")
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface
 *  @template Q
 */ function I(){}
/** @param {T} a
 * @return {T|Q}
 * @template T
*/
I.prototype.method;
/** @constructor 
 *  @implements {I<R>}
 *  @template R
 */ function C(){}
/** @override*/ C.prototype.method = function(a) {}
/** @type {C<number>} var x = new C();
/** @type {null} */ var some = x.method('str');
`).
		AddDiagnostic(`initializing variable
found   : (number|string)
required: null
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testTemplatized9()
func (c *TypeCheckTemplatizedTest) TestTemplatized9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface
 *  @template Q
 */ function I(){}
/** @param {T} a
 * @return {T|Q}
 * @template T
*/
I.prototype.method;
/** @constructor 
 *  @param {R} a
 *  @implements {I<R>}
 *  @template R
 */ function C(a){}
/** @override*/ C.prototype.method = function(a) {}
/** @type {null} */ var some = new C(1).method('str');
`).
		AddDiagnostic(`initializing variable
found   : (number|string)
required: null
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testTemplatized10()
func (c *TypeCheckTemplatizedTest) TestTemplatized10() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @template T
 */
function Parent() {};

/** @param {T} x */
Parent.prototype.method = function(x) {};

/**
 * @constructor
 * @extends {Parent<string>}
 */
function Child() {};
Child.prototype = new Parent();

(new Child()).method(123);
`).
		AddDiagnostic(`actual parameter 1 of Parent.prototype.method does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testTemplatized11()
func (c *TypeCheckTemplatizedTest) TestTemplatized11() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @template T
 * @constructor
 */
function C() {}

/**
 * @param {T|K} a
 * @return {T}
 * @template K
 */
C.prototype.method = function(a) {};

// method returns "?"
/** @type {void} */ var x = new C().method(1);
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testTemplatizedTypeSubtypes2()
func (c *TypeCheckTemplatizedTest) TestTemplatizedTypeSubtypes2() {
	var arrayOfNumber rhino.JSType = c.CreateTemplatizedType(c.GetNativeArrayType(), c.GetNativeNumberType())
	var arrayOfString rhino.JSType = c.CreateTemplatizedType(c.GetNativeArrayType(), c.GetNativeStringType())
	TypeSubject_AssertType(arrayOfString).IsNotSubtypeOf(c.CreateUnionType(arrayOfNumber, c.GetNativeNullVoidType()))
}

// Java: TypeCheckTemplatizedTest.testTemplatizedStructuralMatch1()
func (c *TypeCheckTemplatizedTest) TestTemplatizedStructuralMatch1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @record @template T */
function WithPropT() {}
/** @type {T} */ WithPropT.prototype.prop;
function f(/** !WithPropT<number> */ x){}
/** @constructor */ function Foo() {}
/** @type {number} */ Foo.prototype.prop;
f(new Foo);
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testTemplatizedStructuralMatch2()
func (c *TypeCheckTemplatizedTest) TestTemplatizedStructuralMatch2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @record @template T */
function WithPropT() {}
/** @type {T} */ WithPropT.prototype.prop
function f(/** !WithPropT<number> */ x){};
/** @constructor @template U */ function Foo() {}
/** @type {number} */ Foo.prototype.prop
f(new Foo)
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testTemplatizedStructuralMatch3()
func (c *TypeCheckTemplatizedTest) TestTemplatizedStructuralMatch3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @record @template T */
function WithPropT() {}
/** @type {T} */ WithPropT.prototype.prop
function f(/** !WithPropT<string> */ x){};
/** @constructor @template U */ function Foo() {}
/** @type {U} */ Foo.prototype.prop
f(new Foo)
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testTemplatizedStructuralMismatch1()
func (c *TypeCheckTemplatizedTest) TestTemplatizedStructuralMismatch1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @record @template T */
function WithPropT() {}
/** @type {T} */ WithPropT.prototype.prop
function f(/** !WithPropT<number> */ x){};
/** @constructor */ function Foo() {}
/** @type {string} */ Foo.prototype.prop = 'str'
f(new Foo)
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : Foo
required: WithPropT<number>
missing : []
mismatch: [prop]
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testTemplatizedStructuralMismatch2()
func (c *TypeCheckTemplatizedTest) TestTemplatizedStructuralMismatch2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @record @template T */
function WithPropT() {}
/** @type {T} */ WithPropT.prototype.prop
function f(/** !WithPropT<number> */ x){};
/** @constructor @template U */ function Foo() {}
/** @type {string} */ Foo.prototype.prop = 'str'
f(new Foo)
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : Foo<?>
required: WithPropT<number>
missing : []
mismatch: [prop]
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testTemplatizedStructuralMismatch3()
func (c *TypeCheckTemplatizedTest) TestTemplatizedStructuralMismatch3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @record @template T */
function WithPropT() {}
/** @type {T} */ WithPropT.prototype.prop
function f(/** !WithPropT<number> */ x){};
/**
 * @constructor
 * @template U
 * @param {U} x
 */
function Foo(x) {
  /** @type {U} */ this.prop = x
}
f(new Foo('str'))
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : Foo<string>
required: WithPropT<number>
missing : []
mismatch: [prop]
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testTemplatizedStructuralMismatch4()
func (c *TypeCheckTemplatizedTest) TestTemplatizedStructuralMismatch4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @record @template T */
function WithProp() {}
/** @type {T} */ WithProp.prototype.prop;
/** @constructor */
function Foo() {
  /** @type {number} */ this.prop = 4;
}
/**
 * @template U
 * @param {!WithProp<U>} x
 * @param {U} y
 */
function f(x, y){};
f(new Foo, 'str')
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : Foo
required: WithProp<string>
missing : []
mismatch: [prop]
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testTemplatizedStructuralMismatchNotFound()
func (c *TypeCheckTemplatizedTest) TestTemplatizedStructuralMismatchNotFound() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @record @template T */
function WithProp() {}
/** @type {T} */ WithProp.prototype.prop;
/** @constructor */
function Foo() {
  /** @type {number} */ this.prop = 4;
}
/** @constructor */
function Bar() {
  /** @type {string} */ this.prop = 'str';
}
/**
 * @template U
 * @param {!WithProp<U>} x
 * @param {!WithProp<U>} y
 */
function f(x, y){};
f(new Foo, new Bar)
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testRecursiveTemplatizedType_siblingGenericReferences()
func (c *TypeCheckTemplatizedTest) TestRecursiveTemplatizedType_siblingGenericReferences() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @template T, U */
class Foo {
  /** @return {!Foo<T|U, symbol>} */
  x() {
    return new Foo();
  }
}
/** @type {!Foo<number, string>} */
const a = new Foo();
/** @type {null} expected to error */
const b = a.x;
`).
		AddDiagnostic(`initializing variable
found   : function(this:Foo): Foo<(number|string),symbol>
required: null
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testRecursiveTemplatizedType_siblingGenericReferences_complex()
func (c *TypeCheckTemplatizedTest) TestRecursiveTemplatizedType_siblingGenericReferences_complex() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @template T, U, V */
class Foo {
  /** @return {!Foo<!Foo<U|V|symbol, null, !Array<!Foo<T|U>>>, T, V>} */
  x() {
    return new Foo();
  }
}
/** @type {!Foo<number, string, undefined>} */
const a = new Foo();
/** @type {null} expected to error */
const b = a.x;
`).
		AddDiagnostic(`initializing variable
found   : function(this:Foo): Foo<Foo<(string|symbol|undefined),null,Array<Foo<(number|string),?,?>>>,number,undefined>
required: null
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testb326131100_recursiveStackOverflow()
func (c *TypeCheckTemplatizedTest) Testb326131100_recursiveStackOverflow() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @template T */
class Foo {
  /** @return {!Foo<!Foo<?T>>} */
  asNullable() {
    return new Foo();
  }
}

/** @type {!Foo<number>} */
const a = new Foo();
/** @type {null} expected to error */
const b = a.asNullable();
`).
		AddDiagnostic(`initializing variable
found   : Foo<Foo<(null|number)>>
required: null
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testRecursiveTemplatizedType_unionWithMultipleTypes()
func (c *TypeCheckTemplatizedTest) TestRecursiveTemplatizedType_unionWithMultipleTypes() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @template T, U */
class Foo {
  /** @return {!Foo<!Foo<T|U, symbol>, *>} */
  x() {
    return new Foo();
  }
}
/** @type {!Foo<number, string>} */
const a = new Foo();
/** @type {null} expected to error */
const b = a.x;
`).
		AddDiagnostic(`initializing variable
found   : function(this:Foo): Foo<Foo<(number|string),symbol>,*>
required: null
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testRecursiveTemplatizedType_deeplyNestedTemplatizedTypes()
func (c *TypeCheckTemplatizedTest) TestRecursiveTemplatizedType_deeplyNestedTemplatizedTypes() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @template T */
class Foo {
  /** @return {!Foo<!Foo<!Array<(!Foo<T> | !Array<?Foo<T>> | undefined )>>>} */
  asNestedArray() {
    return new Foo();
  }
}
/** @type {!Foo<number>} */
const a = new Foo();
/** @type {null} expected to error */
const b = a.asNestedArray();
`).
		AddDiagnostic(`initializing variable
found   : Foo<Foo<Array<(Array<(Foo<number>|null)>|Foo<number>|undefined)>>>
required: null
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testSiblingAndInheritedGenerics()
func (c *TypeCheckTemplatizedTest) TestSiblingAndInheritedGenerics() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @template T, U */
class Parent {
  /** @return {!Parent<U, T>} */
  x() {}
}
/**
 * @template V, X
 * @extends {Parent<X, V>}
 */
class Child extends Parent {
  /** @return {!Child<X, V>} */
  y() {}
}
/** @type {!Child<number, null>} */
const a = new Child();
/** @type {null} expected to error */
const b = a.x();
/** @type {null} expected to error */
const c = a.y();
`).
		AddDiagnostic(`initializing variable
found   : Parent<number,null>
required: null
`).
		AddDiagnostic(`initializing variable
found   : Child<null,number>
required: null
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testSiblingAndInheritedGenerics_withUnions()
func (c *TypeCheckTemplatizedTest) TestSiblingAndInheritedGenerics_withUnions() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @template T, U */
class Parent {
  /** @return {!Parent<T|U, string>} */
  x() {}
}
/**
 * @template V, X
 * @extends {Parent<V|X, X|symbol>}
 */
class Child extends Parent {
  /** @return {!Child<!Array<X>, !Parent<V, V>>} */
  y() {}
}
/** @type {!Child<number, null>} */
const a = new Child();
/** @type {null} expected to error */
const b = a.x();
/** @type {null} expected to error */
const c = a.y();
`).
		AddDiagnostic(`initializing variable
found   : Parent<(null|number|symbol),string>
required: null
`).
		AddDiagnostic(`initializing variable
found   : Child<Array<null>,Parent<number,number>>
required: null
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testTemplateTypeOnEnclosingClass_usedInLet_es5()
func (c *TypeCheckTemplatizedTest) TestTemplateTypeOnEnclosingClass_usedInLet_es5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @template T
 * @constructor
 */
function GenericClass() {}
GenericClass.prototype.foo = function() {
  let anonymousFunction = function () {
    let /** ?T */ tRef = null;
  };
};
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testTemplateTypeOnEnclosingClass_usedInLet()
func (c *TypeCheckTemplatizedTest) TestTemplateTypeOnEnclosingClass_usedInLet() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @template T */
class GenericClass {
  /** @template V */
  foo(/** V */ param) {
    let anonymousFunction = function () {
      let /** ?T */ tRef = null;
      let another = function() {
        function andYetAnother() {
          let /** ?T */ tRefNested = null;
        }
      };
    };
  }
}
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testTemplateTypeOnEnclosingClass_usedInNestedFunction()
func (c *TypeCheckTemplatizedTest) TestTemplateTypeOnEnclosingClass_usedInNestedFunction() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @template T */
class GenericClass {
  foo() {
    let anonymousFunction = function () {
      /** @param {T} p */
      let another = function(p) {
      };
    };
  }
}
`).
		Run()
}

// Java: TypeCheckTemplatizedTest.testWellKnownSymbol_instantiatedInTemplateType()
func (c *TypeCheckTemplatizedTest) TestWellKnownSymbol_instantiatedInTemplateType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @template T */
class Foo {
  /** @return {!Array<T>} */
  x() {
    return [];
  }
}
/** @type {!Foo<typeof Symbol.iterator>} */
const a = new Foo();
/** @type {null} expected to error */
const b = a.x();
`).
		AddDiagnostic(`initializing variable
found   : Array<Symbol.iterator>
required: null
`).
		Run()
}
