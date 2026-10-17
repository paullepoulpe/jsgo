// Port of com/google/javascript/jscomp/TypeCheckCovarianceTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/TypeCheckCovarianceTest.java)

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
	"testing"
)

// TypeCheckCovarianceTest tests the type checker.
type TypeCheckCovarianceTest struct {
	typeCheckTests_noFixture
}

// Java: TypeCheckCovarianceTest()
func NewTypeCheckCovarianceTest() *TypeCheckCovarianceTest {
	c := &TypeCheckCovarianceTest{}
	return c
}

func TestTypeCheckCovarianceTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewTypeCheckCovarianceTest() })
}

// Java: TypeCheckCovarianceTest.testIterableCovariant()
func (c *TypeCheckCovarianceTest) TestIterableCovariant() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** !Iterable<(number|string)>*/ x){};
function g(/** !Iterable<number> */ arr) {
    f(arr);
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckCovarianceTest.testLocalShadowOfIterableNotCovariant()
func (c *TypeCheckCovarianceTest) TestLocalShadowOfIterableNotCovariant() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @template T */
class Iterable {}
function f(/** !Iterable<(number|string)>*/ x) {};
function g(/** !Iterable<number> */ arr) {
    f(arr);
}
export {};
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : Iterable<number>
required: Iterable<(number|string)>
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckCovarianceTest.testIterableNotContravariant()
func (c *TypeCheckCovarianceTest) TestIterableNotContravariant() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** !Iterable<number>*/ x){};
function g(/** !Iterable<(number|string)> */ arr) {
    f(arr);
}
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : Iterable<(number|string),?,?>
required: Iterable<number,?,?>
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckCovarianceTest.testIterableCovariantWhenComparingToSubtype()
func (c *TypeCheckCovarianceTest) TestIterableCovariantWhenComparingToSubtype() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @constructor
 * @implements {Iterable<T>}
 * @template T
 */
function Set() {}
/** @override */
Set.prototype[Symbol.iterator] = function() {};
`).
		AddSource(`function f(/** !Iterable<(number|string)>*/ x){};
function g(/** !Set<number> */ arr) {
    f(arr);
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckCovarianceTest.testIteratorCovariant()
func (c *TypeCheckCovarianceTest) TestIteratorCovariant() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** !Iterator<(string|number)>*/ x){};
function g(/** !Iterator<number> */ arr) {
    f(arr);
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckCovarianceTest.testIteratorLikeCovariant()
func (c *TypeCheckCovarianceTest) TestIteratorLikeCovariant() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** !IteratorLike<(string|number)>*/ x){};
function g(/** !IteratorLike<number> */ arr) {
    f(arr);
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckCovarianceTest.testIteratorIterableCovariant()
func (c *TypeCheckCovarianceTest) TestIteratorIterableCovariant() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** !IteratorIterable<(string|number)>*/ x){};
function g(/** !IteratorIterable<number> */ arr) {
    f(arr);
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckCovarianceTest.testIIterableResultCovariant()
func (c *TypeCheckCovarianceTest) TestIIterableResultCovariant() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** !IIterableResult<(string|number)>*/ x){};
function g(/** !IIterableResult<number> */ arr) {
    f(arr);
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckCovarianceTest.testGeneratorCovariant()
func (c *TypeCheckCovarianceTest) TestGeneratorCovariant() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** !Generator<(string|number)>*/ x){};
function g(/** !Generator<number> */ arr) {
    f(arr);
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckCovarianceTest.testIterableImplementorInvariant()
func (c *TypeCheckCovarianceTest) TestIterableImplementorInvariant() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/**
 * @constructor
 * @implements {Iterable<T>}
 * @template T
 */
function Set() {}
/** @override */
Set.prototype[Symbol.iterator] = function() {};
`).
		IncludeDefaultExterns().
		AddSource(`function f(/** !Set<(string|number)>*/ x){};
function g(/** !Set<number> */ arr) {
    f(arr);
}
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : Set<number>
required: Set<(number|string)>
`).
		Run()
}

// Java: TypeCheckCovarianceTest.testIArrayLikeCovariant1()
func (c *TypeCheckCovarianceTest) TestIArrayLikeCovariant1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** !IArrayLike<(string|number)>*/ x){};
function g(/** !IArrayLike<number> */ arr) {
    f(arr);
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckCovarianceTest.testIArrayLikeCovariant2()
func (c *TypeCheckCovarianceTest) TestIArrayLikeCovariant2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** !IArrayLike<(string|number)>*/ x){};
function g(/** !Array<number> */ arr) {
    f(arr);
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckCovarianceTest.testIArrayLikeBivaraint()
func (c *TypeCheckCovarianceTest) TestIArrayLikeBivaraint() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** !IArrayLike<number>*/ x){};
function g(/** !IArrayLike<(string|number)> */ arr) {
    f(arr);
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckCovarianceTest.testCovarianceForRecordType1()
func (c *TypeCheckCovarianceTest) TestCovarianceForRecordType1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @constructor */
function C() {}
/** @constructor
  * @extends {C}
  */
function C2() {}
`).
		AddSource(`/** @type {{prop: C}} */
var r1;
/** @type {{prop: C2}} */
var r2;
r1 = r2;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckCovarianceTest.testCovarianceForRecordType2()
func (c *TypeCheckCovarianceTest) TestCovarianceForRecordType2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @constructor */
function C() {}
/** @constructor
  * @extends {C}
  */
function C2() {}
`).
		AddSource(`/** @type {{prop: C, prop2: C}} */
var r1;
/** @type {{prop: C2, prop2: C}} */
var r2;
r1 = r2;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckCovarianceTest.testCovarianceForRecordType3()
func (c *TypeCheckCovarianceTest) TestCovarianceForRecordType3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @constructor */
function C() {}
/** @constructor @extends {C} */
function C2() {}
`).
		AddSource(`/** @type {{prop: C}} */
var r1;
/** @type {{prop: C2, prop2: C}} */
var r2;
r1 = r2;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckCovarianceTest.testCovarianceForRecordType4()
func (c *TypeCheckCovarianceTest) TestCovarianceForRecordType4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @constructor */
function C() {}
/** @constructor @extends {C} */
function C2() {}
`).
		AddSource(`/** @type {{prop: C, prop2: C}} */
var r1;
/** @type {{prop: C2}} */
var r2;
r1 = r2;
`).
		AddDiagnostic(`assignment
found   : {prop: (C2|null)}
required: {
  prop: (C|null),
  prop2: (C|null)
}
missing : [prop2]
mismatch: []
`).
		Run()
}

// Java: TypeCheckCovarianceTest.testCovarianceForRecordType5()
func (c *TypeCheckCovarianceTest) TestCovarianceForRecordType5() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @constructor */
function C() {}
/** @constructor */
function C2() {}
`).
		AddSource(`/** @type {{prop: C}} */
var r1;
/** @type {{prop: C2}} */
var r2;
r1 = r2;
`).
		AddDiagnostic(`assignment
found   : {prop: (C2|null)}
required: {prop: (C|null)}
missing : []
mismatch: [prop]
`).
		Run()
}

// Java: TypeCheckCovarianceTest.testCovarianceForRecordType6()
func (c *TypeCheckCovarianceTest) TestCovarianceForRecordType6() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @constructor */
function C() {}
/** @constructor @extends {C} */
function C2() {}
`).
		AddSource(`/** @type {{prop: C2}} */
var r1;
/** @type {{prop: C}} */
var r2;
r1 = r2;
`).
		AddDiagnostic(`assignment
found   : {prop: (C|null)}
required: {prop: (C2|null)}
missing : []
mismatch: [prop]
`).
		Run()
}

// Java: TypeCheckCovarianceTest.testCovarianceForRecordType7()
func (c *TypeCheckCovarianceTest) TestCovarianceForRecordType7() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @constructor */
function C() {}
/** @constructor @extends {C} */
function C2() {}
`).
		AddSource(`/** @type {{prop: C2, prop2: C2}} */
var r1;
/** @type {{prop: C2, prop2: C}} */
var r2;
r1 = r2;
`).
		AddDiagnostic(`assignment
found   : {
  prop: (C2|null),
  prop2: (C|null)
}
required: {
  prop: (C2|null),
  prop2: (C2|null)
}
missing : []
mismatch: [prop2]
`).
		Run()
}

// Java: TypeCheckCovarianceTest.testCovarianceForRecordType8()
func (c *TypeCheckCovarianceTest) TestCovarianceForRecordType8() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @constructor */
function Foo(){}
/** @type {number} */
Foo.prototype.x = 5
/** @type {string} */
Foo.prototype.y = 'str'
`).
		AddSource(`/** @type {{x: number, y: string}} */
var r1 = {x: 1, y: 'value'};

/** @type {!Foo} */
var f = new Foo();
r1 = f;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckCovarianceTest.testCovarianceForRecordType9()
func (c *TypeCheckCovarianceTest) TestCovarianceForRecordType9() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @constructor */
function Foo(){}
/** @type {number} */
Foo.prototype.x1 = 5
/** @type {string} */
Foo.prototype.y = 'str'
`).
		AddSource(`/** @type {{x: number, y: string}} */
var r1 = {x: 1, y: 'value'};

/** @type {!Foo} */
var f = new Foo();
f = r1;
`).
		AddDiagnostic(`assignment
found   : {
  x: number,
  y: string
}
required: Foo
`).
		Run()
}

// Java: TypeCheckCovarianceTest.testCovarianceForRecordType10()
func (c *TypeCheckCovarianceTest) TestCovarianceForRecordType10() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @constructor */
function Foo() {}
/** @type {{x: !Foo}} */
Foo.prototype.x = {x: new Foo()};
`).
		AddSource(`/** @type {!Foo} */
var o = new Foo();

/** @type {{x: !Foo}} */
var r = {x : new Foo()};
r = o;
`).
		AddDiagnostic(`assignment
found   : Foo
required: {x: Foo}
missing : []
mismatch: [x]
`).
		Run()
}

// Java: TypeCheckCovarianceTest.testCovarianceForRecordType11()
func (c *TypeCheckCovarianceTest) TestCovarianceForRecordType11() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @interface */
function Foo() {}
/** @constructor @implements {Foo} */
function Bar1() {}
/** @return {number} */
Bar1.prototype.y = function(){return 1;};
/** @constructor @implements {Foo} */
function Bar() {}
/** @return {string} */
Bar.prototype.y = function(){return 'test';};
`).
		AddSource(`function fun(/** Foo */f) {
  f.y();
}
fun(new Bar1())
fun(new Bar());
`).
		AddDiagnosticType(TypeCheck_STRICT_INEXISTENT_PROPERTY).
		Run()
}

// Java: TypeCheckCovarianceTest.testCovarianceForRecordType12()
func (c *TypeCheckCovarianceTest) TestCovarianceForRecordType12() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @interface */
function Foo() {}
/** @constructor @implements {Foo} */
function Bar1() {}
/** @constructor @implements {Foo} */
function Bar() {}
/** @return {undefined} */
Bar.prototype.y = function(){};
`).
		AddSource(`/** @type{Foo} */
var f = new Bar1();
f.y();
`).
		AddDiagnosticType(TypeCheck_STRICT_INEXISTENT_PROPERTY).
		Run()
}

// Java: TypeCheckCovarianceTest.testCovarianceForRecordType13()
func (c *TypeCheckCovarianceTest) TestCovarianceForRecordType13() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @interface */
function I() {}
/** @constructor @implements {I} */
function C() {}
/** @return {undefined} */
C.prototype.y = function(){};
`).
		AddSource(`/** @type{{x: {obj: I}}} */
var ri;
ri.x.obj.y();
`).
		AddDiagnosticType(TypeCheck_STRICT_INEXISTENT_PROPERTY).
		Run()
}

// Java: TypeCheckCovarianceTest.testCovarianceForRecordType14a()
func (c *TypeCheckCovarianceTest) TestCovarianceForRecordType14a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @interface */
function I() {}
/** @constructor */
function C() {}
/** @return {undefined} */
C.prototype.y = function(){};
`).
		AddSource(`/** @type{({x: {obj: I}}|{x: {obj: C}})} */
var ri;
ri.x.obj.y();
`).
		IncludeDefaultExterns().
		Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		Run()
}

// Java: TypeCheckCovarianceTest.testCovarianceForRecordType14b()
func (c *TypeCheckCovarianceTest) TestCovarianceForRecordType14b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().IncludeDefaultExterns().
		AddExterns(`/** @interface */
function I() {}
/** @constructor */
function C() {}
/** @return {undefined} */
C.prototype.y = function(){};
`).
		AddSource(`/** @type{({x: {obj: I}}|{x: {obj: C}})} */
var ri;
ri.x.obj.y();
`).
		AddDiagnostic("Property y not defined on all member types of (I|C)").
		Run()
}

// Java: TypeCheckCovarianceTest.testCovarianceForRecordType15()
func (c *TypeCheckCovarianceTest) TestCovarianceForRecordType15() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @constructor */
function C() {}
/** @return {undefined} */
C.prototype.y1 = function(){};
/** @constructor */
function C1() {}
/** @return {undefined} */
C1.prototype.y = function(){};
`).
		AddSource(`/** @type{({x: {obj: C}}|{x: {obj: C1}})} */
var ri;
ri.x.obj.y1();
ri.x.obj.y();
`).
		IncludeDefaultExterns().
		Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		Run()
}

// Java: TypeCheckCovarianceTest.testCovarianceForRecordType16()
func (c *TypeCheckCovarianceTest) TestCovarianceForRecordType16() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @constructor */
function C() {}
/** @return {number} */
C.prototype.y = function(){return 1;};
/** @constructor */
function C1() {}
/** @return {string} */
C1.prototype.y = function(){return 'test';};
`).
		AddSource(`/** @type{({x: {obj: C}}|{x: {obj: C1}})} */
var ri;
ri.x.obj.y();
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckCovarianceTest.testCovarianceForRecordType17()
func (c *TypeCheckCovarianceTest) TestCovarianceForRecordType17() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @interface */
function Foo() {}
/** @constructor @implements {Foo} */
function Bar1() {}
Bar1.prototype.y = function(){return {};};
/** @constructor @implements {Foo} */
function Bar() {}
/** @return {number} */
Bar.prototype.y = function(){return 1;};
`).
		AddSource(`/** @type {Foo} */ var f;
f.y();
`).
		AddDiagnosticType(TypeCheck_STRICT_INEXISTENT_PROPERTY).
		Run()
}

// Java: TypeCheckCovarianceTest.testCovarianceForRecordType18()
func (c *TypeCheckCovarianceTest) TestCovarianceForRecordType18() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @constructor*/
function Bar1() {}
/** @type {{x: number}} */
Bar1.prototype.prop;
/** @constructor */
function Bar() {}
/** @type {{x: number, y: number}} */
Bar.prototype.prop;
`).
		AddSource(`/** @type {{x: number}} */ var f;
f.z;
`).
		IncludeDefaultExterns().
		Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		Run()
}

// Java: TypeCheckCovarianceTest.testCovarianceForRecordType19a()
func (c *TypeCheckCovarianceTest) TestCovarianceForRecordType19a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @constructor */
function Bar1() {}
/** @type {number} */
Bar1.prototype.prop;
/** @type {number} */
Bar1.prototype.prop1;
/** @constructor */
function Bar2() {}
/** @type {number} */
Bar2.prototype.prop;
`).
		AddSource(`/** @type {(Bar1|Bar2)} */ var b;
var x = b.prop1
`).
		IncludeDefaultExterns().
		Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		Run()
}

// Java: TypeCheckCovarianceTest.testCovarianceForRecordType19b()
func (c *TypeCheckCovarianceTest) TestCovarianceForRecordType19b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().IncludeDefaultExterns().
		AddExterns(`/** @constructor */
function Bar1() {}
/** @type {number} */
Bar1.prototype.prop;
/** @type {number} */
Bar1.prototype.prop1;
/** @constructor */
function Bar2() {}
/** @type {number} */
Bar2.prototype.prop;
`).
		AddSource(`/** @type {(Bar1|Bar2)} */ var b;
var x = b.prop1
`).
		AddDiagnostic("Property prop1 not defined on all member types of (Bar1|Bar2)").
		Run()
}

// Java: TypeCheckCovarianceTest.testCovarianceForRecordType20()
func (c *TypeCheckCovarianceTest) TestCovarianceForRecordType20() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().IncludeDefaultExterns().
		AddExterns(`/** @constructor */
function Bar1() {}
/** @type {number} */
Bar1.prototype.prop;
/** @type {number} */
Bar1.prototype.prop1;
/** @type {number} */
Bar1.prototype.prop2;
`).
		AddSource(`/** @type {{prop2:number}} */ var c;
/** @type {(Bar1|{prop:number, prop2: number})} */ var b;
// there should be no warning saying that
// prop2 is not defined on b;
var x = b.prop2
`).
		Run()
}

// Java: TypeCheckCovarianceTest.testCovarianceForRecordType20_2()
func (c *TypeCheckCovarianceTest) TestCovarianceForRecordType20_2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().IncludeDefaultExterns().
		AddSource(`/** @type {{prop2:number}} */ var c;
/** @type {({prop:number, prop1: number, prop2: number}|
{prop:number, prop2: number})} */ var b;
// there should be no warning saying that
// prop2 is not defined on b;
var x = b.prop2
`).
		Run()
}

// Java: TypeCheckCovarianceTest.testCovarianceForRecordType21()
func (c *TypeCheckCovarianceTest) TestCovarianceForRecordType21() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns("").
		AddSource(`/** @constructor */
function Bar1() {};
/** @type {number} */
Bar1.prototype.propName;
/** @type {number} */
Bar1.prototype.propName1;
/** @type {{prop2:number}} */ var c;
/** @type {(Bar1|{propName:number, propName1: number})} */ var b;
var x = b.prop2;
`).
		AddDiagnostic("Property prop2 never defined on b").
		Run()
}

// Java: TypeCheckCovarianceTest.testCovarianceForRecordType23()
func (c *TypeCheckCovarianceTest) TestCovarianceForRecordType23() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @constructor */
function A() {}
/** @constructor @extends{A} */
function B() {}

/** @constructor */
function C() {}
/** @type {B} */
C.prototype.prop2;
/** @type {number} */
C.prototype.prop3;

/** @constructor */
function D() {}
/** @type {number} */
D.prototype.prop;
/** @type {number} */
D.prototype.prop1;
/** @type {B} */
D.prototype.prop2;
`).
		AddSource(`/** @type {{prop2: A}} */ var record;
var xhr = new C();
if (true) { xhr = new D(); }
// there should be no warning saying that
// prop2 is not defined on b;
var x = xhr.prop2
`).
		Run()
}

// Java: TypeCheckCovarianceTest.testCovarianceForRecordType24()
func (c *TypeCheckCovarianceTest) TestCovarianceForRecordType24() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @constructor */
function C() {}

/** @type {!Function} */
C.prototype.abort = function() {};

/** @type{number} */
C.prototype.test2 = 1;
`).
		AddSource(`function f() {
  /** @type{{abort: !Function, count: number}} */
  var x;
}

function f2() {
  /** @type{(C|{abort: Function})} */
  var y;
  y.abort();
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckCovarianceTest.testCovarianceForRecordType25()
func (c *TypeCheckCovarianceTest) TestCovarianceForRecordType25() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @constructor */
function C() {}

/** @type {!Function} */
C.prototype.abort = function() {};

/** @type{number} */
C.prototype.test2 = 1;
`).
		AddSource(`function f() {
  /** @type{!Function} */ var f;
  var x = {abort: f, count: 1}
  return x;
}

function f2() {
  /** @type{(C|{test2: number})} */
  var y;
  y.abort();
}
`).
		AddDiagnosticType(TypeCheck_STRICT_INEXISTENT_UNION_PROPERTY).
		Run()
}

// Java: TypeCheckCovarianceTest.testCovarianceForRecordType26()
func (c *TypeCheckCovarianceTest) TestCovarianceForRecordType26() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @constructor */
function C() {}

C.prototype.abort = function() {};

/** @type{number} */
C.prototype.test2 = 1;
`).
		AddSource(`function f() {
  /** @type{{abort: !Function}} */
  var x;
}

function f2() {
  /** @type{(C|{test2: number})} */
  var y;
  /** @type {C} */ (y).abort();
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckCovarianceTest.testCovarianceForRecordType26AndAHalf()
func (c *TypeCheckCovarianceTest) TestCovarianceForRecordType26AndAHalf() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @constructor */
function C() {}

C.prototype.abort = function() {};

/** @type{number} */
C.prototype.test2 = 1;
var g = function /** !C */(){};
`).
		AddSource(`function f() {
  /** @type{{abort: !Function}} */
  var x;
}
function f2() {
  var y = g();
  y.abort();
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckCovarianceTest.testCovarianceForRecordType27()
func (c *TypeCheckCovarianceTest) TestCovarianceForRecordType27() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @constructor */
function C(){}
/** @constructor @extends {C} */
function C2() {}
`).
		AddSource(`/** @type {{prop2:C}} */ var c;
/** @type {({prop:number, prop1: number, prop2: C}|
{prop:number, prop1: number, prop2: number})} */ var b;
var x = b.prop2;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckCovarianceTest.testCovarianceForRecordType28()
func (c *TypeCheckCovarianceTest) TestCovarianceForRecordType28() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @constructor */
function XMLHttpRequest() {}
/**
 * @return {undefined}
 */
XMLHttpRequest.prototype.abort = function() {};

/** @constructor */
function XDomainRequest() {}

XDomainRequest.prototype.abort = function() {};
`).
		AddSource(`/**
 * @typedef {{abort: !Function, close: !Function}}
 */
var WritableStreamSink;
function sendCrossOrigin() {
  var xhr = new XMLHttpRequest;
  xhr = new XDomainRequest;
  return function() {
    xhr.abort();
  };
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckCovarianceTest.testCovarianceForRecordType29()
func (c *TypeCheckCovarianceTest) TestCovarianceForRecordType29() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @constructor */
function XMLHttpRequest() {}
/**
 * @type {!Function}
 */
XMLHttpRequest.prototype.abort = function() {};

/** @constructor */
function XDomainRequest() {}
/**
 * @type {!Function}
 */
XDomainRequest.prototype.abort = function() {};
`).
		AddSource(`/**
 * @typedef {{close: !Function, abort: !Function}}
 */
var WritableStreamSink;
function sendCrossOrigin() {
  var xhr = new XMLHttpRequest;
  xhr = new XDomainRequest;
  return function() {
    xhr.abort();
  };
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckCovarianceTest.testCovarianceForRecordType30()
func (c *TypeCheckCovarianceTest) TestCovarianceForRecordType30() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @constructor */
function A() {}
`).
		AddSource(`/**
 * @type {{prop1: (A)}}
 */
var r1;
/**
 * @type {{prop1: (A|undefined)}}
 */
var r2;
r1 = r2
`).
		AddDiagnostic(`assignment
found   : {prop1: (A|null|undefined)}
required: {prop1: (A|null)}
missing : []
mismatch: [prop1]
`).
		Run()
}

// Java: TypeCheckCovarianceTest.testCovarianceForRecordType31()
func (c *TypeCheckCovarianceTest) TestCovarianceForRecordType31() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(`/** @constructor */
function A() {}
`).
		AddSource(`/**
 * @type {{prop1: (A|null)}}
 */
var r1;
/**
 * @type {{prop1: (A|null|undefined)}}
 */
var r2;
r1 = r2
`).
		AddDiagnostic(`assignment
found   : {prop1: (A|null|undefined)}
required: {prop1: (A|null)}
missing : []
mismatch: [prop1]
`).
		Run()
}

// Java: TypeCheckCovarianceTest.testCovariantIThenable1()
func (c *TypeCheckCovarianceTest) TestCovariantIThenable1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {!IThenable<string|number>} */ var x;
function fn(/** !IThenable<string> */ a ) {
  x = a;
}
`).
		Run()
}

// Java: TypeCheckCovarianceTest.testCovariantIThenable2()
func (c *TypeCheckCovarianceTest) TestCovariantIThenable2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {!IThenable<string>} */ var x;
function fn(/** !IThenable<string|number> */ a ) {
  x = a;
}
`).
		AddDiagnostic(`assignment
found   : IThenable<(number|string)>
required: IThenable<string>
`).
		Run()
}

// Java: TypeCheckCovarianceTest.testCovariantIThenable3()
func (c *TypeCheckCovarianceTest) TestCovariantIThenable3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {!Promise<string|number>} */ var x;
function fn(/** !Promise<string> */ a ) {
  x = a;
}
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckCovarianceTest.testCovariantIThenable4()
func (c *TypeCheckCovarianceTest) TestCovariantIThenable4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {!Promise<string>} */ var x;
function fn(/** !Promise<string|number> */ a ) {
  x = a;
}
`).
		AddDiagnostic(`assignment
found   : Promise<(number|string)>
required: Promise<string>
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckCovarianceTest.testCovariantIThenableNonThenable1()
func (c *TypeCheckCovarianceTest) TestCovariantIThenableNonThenable1() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {!Array<string>} */ var x;
function fn(/** !IThenable<string> */ a ) {
  x = a;
}
`).
		AddDiagnostic(`assignment
found   : IThenable<string>
required: Array<string>
`).
		Run()
}

// Java: TypeCheckCovarianceTest.testCovariantIThenableNonThenable2()
func (c *TypeCheckCovarianceTest) TestCovariantIThenableNonThenable2() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {!IThenable<string>} */ var x;
function fn(/** !Array<string> */ a ) {
  x = a;
}
`).
		AddDiagnostic(`assignment
found   : Array<string>
required: IThenable<string>
`).
		Run()
}

// Java: TypeCheckCovarianceTest.testCovariantIThenableNonThenable3()
func (c *TypeCheckCovarianceTest) TestCovariantIThenableNonThenable3() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
  @constructor
  @template T
 */
function C() {}
/** @type {!C<string>} */ var x;
function fn(/** !IThenable<string> */ a ) {
  x = a;
}
`).
		AddDiagnostic(`assignment
found   : IThenable<string>
required: C<string>
`).
		Run()
}

// Java: TypeCheckCovarianceTest.testCovariantIThenableNonThenable4()
func (c *TypeCheckCovarianceTest) TestCovariantIThenableNonThenable4() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
  @constructor
  @template T
 */
function C() {}
/** @type {!IThenable<string>} */ var x;
function fn(/** !C<string> */ a ) {
  x = a;
}
`).
		AddDiagnostic(`assignment
found   : C<string>
required: IThenable<string>
`).
		Run()
}

// Java: TypeCheckCovarianceTest.testCovariantIThenableNonNativeSubclass()
func (c *TypeCheckCovarianceTest) TestCovariantIThenableNonNativeSubclass() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @implements {IThenable<T>}
 * @template T
 */
class CustomPromise {}
/** @type {!CustomPromise<(number|string)>} */ var x;
function fn(/** !CustomPromise<string> */ a ) {
  x = a;
}
`).
		Run()
}

// Java: TypeCheckCovarianceTest.testReadonlyArrayCovariant()
func (c *TypeCheckCovarianceTest) TestReadonlyArrayCovariant() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function f(/** !ReadonlyArray<(number|string)>*/ x) {};
function g(/** !ReadonlyArray<number> */ arr) {
    f(arr);
}
export {};
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckCovarianceTest.testReadonlyArrayCovariant_withArraySubtype()
func (c *TypeCheckCovarianceTest) TestReadonlyArrayCovariant_withArraySubtype() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddArray().AddExtra(`/**
 * @param {...T} args
 * @return {!Array<T>}
 * @this {!IArrayLike<T>}
 * @template T
 */
ReadonlyArray.prototype.foo = function(args) {};
/**
 * @param {...T} args
 * @return {!Array<T>}
 * @this {!IArrayLike<T>}
 * @override
 * @template T
 */
Array.prototype.foo = function(args) {};
`).
		Build()).
		AddSource(`/**
 * @param {!Array<number>} a
 * @return {!ReadonlyArray<*>}
 */
var good = function(a) { return a; };
/**
 * @param {!Array<*>} a
 * @return {!ReadonlyArray<number>}
 */
var bad = function(a) { return a; };
`).
		AddDiagnostic(`inconsistent return type
found   : Array<*>
required: ReadonlyArray<number>
missing : []
mismatch: [filter,forEach,values]
`).
		Run()
}

// Java: TypeCheckCovarianceTest.testReadonlyArrayCovariant_withCustomArray_bug()
func (c *TypeCheckCovarianceTest) TestReadonlyArrayCovariant_withCustomArray_bug() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddArray().AddExtra(`/**
 * @constructor
 * @extends {Array<T>}
 * @template T, U
 */
function MyArray() {}

/** @param {T} p */
ReadonlyArray.prototype.takes = function(p) {};
/** @override @param {T} p */
Array.prototype.takes = function(p) {};
/** @override @param {T} p */
MyArray.prototype.takes = function(p) {};

/** @return {T} */
ReadonlyArray.prototype.returns = function() {};
/** @override @return {T} */
Array.prototype.returns = function() {};
/** @override @return {U} */
MyArray.prototype.returns = function() {};
`).
		Build()).
		AddSource(`/**
 * @param {!MyArray<number, number>} a
 * @return {!ReadonlyArray<*>}
 */
var good = function(a) { return a; };
/**
 * @param {!MyArray<number, *>} a
 * @return {!ReadonlyArray<number>}
 */
var bad = function(a) { return a; };
`).
		AddDiagnostic(`inconsistent return type
found   : MyArray<number,number>
required: ReadonlyArray<*>
missing : []
mismatch: [includes,takes]
`).
		AddDiagnostic(`inconsistent return type
found   : MyArray<number,*>
required: ReadonlyArray<number>
missing : []
mismatch: [returns]
`).
		Run()
}

// Java: TypeCheckCovarianceTest.testIThenableSubtypes()
func (c *TypeCheckCovarianceTest) TestIThenableSubtypes() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddPromise().AddExtra(`/**
 * @record
 * @struct
 * @template TYPE
 * @extends {IThenable<TYPE>}
 */
function PromiseLike() {}
/**
 * Copy of Closure library's Thenable
 *
 * @interface
 * @extends {IThenable<TYPE>}
 * @template TYPE
 */
function GoogThenable() {}
`).
		Build()).
		AddSource(`/**
 * @param {!GoogThenable<undefined>} t
 * @return {!PromiseLike}
 */
function f(t) { return t; };
`).
		Run()
}
