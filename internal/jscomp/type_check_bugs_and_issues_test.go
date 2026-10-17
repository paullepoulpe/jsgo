// Port of com/google/javascript/jscomp/TypeCheckBugsAndIssuesTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/TypeCheckBugsAndIssuesTest.java)

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

// TypeCheckBugsAndIssuesTest tests the type checker.
type TypeCheckBugsAndIssuesTest struct {
	typeCheckTests_noFixture
}

// Java: TypeCheckBugsAndIssuesTest()
func NewTypeCheckBugsAndIssuesTest() *TypeCheckBugsAndIssuesTest {
	c := &TypeCheckBugsAndIssuesTest{}
	return c
}

func TestTypeCheckBugsAndIssuesTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewTypeCheckBugsAndIssuesTest() })
}

// Java: TypeCheckBugsAndIssuesTest.testIssue61a()
func (c *TypeCheckBugsAndIssuesTest) TestIssue61a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var ns = {};
(function() {
  /** @param {string} b */
  ns.a = function(b) {};
})();
function d() {
  ns.a(123);
}
`).
		AddDiagnostic(`actual parameter 1 of ns.a does not match formal parameter
found   : number
required: string
`).
		Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testIssue61b()
func (c *TypeCheckBugsAndIssuesTest) TestIssue61b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const */ var ns = {};
(function() {
  /** @param {string} b */
  ns.a = function(b) {};
})();
ns.a(123);
`).
		AddDiagnostic(`actual parameter 1 of ns.a does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testIssue61c()
func (c *TypeCheckBugsAndIssuesTest) TestIssue61c() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var ns = {};
(function() {
  /** @param {string} b */
  ns.a = function(b) {};
})();
ns.a(123);
`).
		AddDiagnostic(`actual parameter 1 of ns.a does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testIssue86()
func (c *TypeCheckBugsAndIssuesTest) TestIssue86() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */ function I() {}
/** @return {number} */ I.prototype.get = function(){};
/** @constructor 
 * @implements {I} */ function F() {}
/** @override */ F.prototype.get = function() { return true; };
`).
		AddDiagnostic(`inconsistent return type
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testIssue124()
func (c *TypeCheckBugsAndIssuesTest) TestIssue124() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var t = null;
function test() {
  if (t != null) { t = null; }
  t = 1;
}
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testIssue124b()
func (c *TypeCheckBugsAndIssuesTest) TestIssue124b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var t = null;
function test() {
  if (t != null) { t = null; }
  t = undefined;
}
`).
		AddDiagnostic(`condition always evaluates to false
left : (null|undefined)
right: null
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testIssue259()
func (c *TypeCheckBugsAndIssuesTest) TestIssue259() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {number} x */ function f(x) {}
/** @constructor */
var Clock = function() {
  /** @constructor */
  this.Date = function() {};
  f(new this.Date());
};
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : this.Date
required: number
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testIssue301()
func (c *TypeCheckBugsAndIssuesTest) TestIssue301() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddString().AddArray().
		Build()).
		AddSource(`Array.indexOf = function() {};
var s = 'hello';
alert(s.toLowerCase.indexOf('1'));
`).
		AddDiagnostic("Property indexOf never defined on String.prototype.toLowerCase").
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testIssue368()
func (c *TypeCheckBugsAndIssuesTest) TestIssue368() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo(){}
/**
 * @param {number} one
 * @param {string} two
 */
Foo.prototype.add = function(one, two) {};
/**
 * @constructor
 * @extends {Foo}
 */
function Bar(){}
/** @override */
Bar.prototype.add = function(ignored) {};
(new Bar()).add(1, 2);
`).
		AddDiagnostic(`actual parameter 2 of Bar.prototype.add does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testIssue380()
func (c *TypeCheckBugsAndIssuesTest) TestIssue380() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddArray().AddObject().
		Build()).
		AddSource(`/** @type { function(string): {innerHTML: string} } */
document.getElementById;
var list = /** @type {!Array<string>} */ ['hello', 'you'];
list.push('?');
document.getElementById('node').innerHTML = list.toString();
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testIssue483()
func (c *TypeCheckBugsAndIssuesTest) TestIssue483() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddArray().Build()).
		AddSource(`/** @constructor */ function C() {
  /** @type {?Array} */ this.a = [];
}
C.prototype.f = function() {
  if (this.a.length > 0) {
    g(this.a);
  }
};
/** @param {number} a */ function g(a) {}
`).
		AddDiagnostic(`actual parameter 1 of g does not match formal parameter
found   : Array
required: number
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testIssue537a()
func (c *TypeCheckBugsAndIssuesTest) TestIssue537a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
Foo.prototype = {method: function() {}};
/**
 * @constructor
 * @extends {Foo}
 */
function Bar() {
  Foo.call(this);
  if (this.baz()) this.method(1);
}
Bar.prototype = {
  baz: function() {
    return true;
  }
};
Bar.prototype.__proto__ = Foo.prototype;
`).
		AddDiagnostic(`Function Foo.prototype.method: called with 1 argument(s). Function requires at least 0 argument(s) and no more than 0 argument(s).
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testIssue537b()
func (c *TypeCheckBugsAndIssuesTest) TestIssue537b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
Foo.prototype = {method: function() {}};
/**
 * @constructor
 * @extends {Foo}
 */
function Bar() {
  Foo.call(this);
  if (this.baz(1)) this.method();
}
Bar.prototype = {
  baz: function() {
    return true;
  }
};
Bar.prototype.__proto__ = Foo.prototype;
`).
		AddDiagnostic(`Function Bar.prototype.baz: called with 1 argument(s). Function requires at least 0 argument(s) and no more than 0 argument(s).
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testIssue537c()
func (c *TypeCheckBugsAndIssuesTest) TestIssue537c() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
/**
 * @constructor
 * @extends {Foo}
 */
function Bar() {
  Foo.call(this);
  if (this.baz2()) alert(1);
}
Bar.prototype = {
  baz: function() {
    return true;
  }
};
Bar.prototype.__proto__ = Foo.prototype;
`).
		AddDiagnostic("Property baz2 never defined on Bar").
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testIssue537d()
func (c *TypeCheckBugsAndIssuesTest) TestIssue537d() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
Foo.prototype = {
  /** @return {Bar} */ x: function() { new Bar(); },
  /** @return {Foo} */ y: function() { new Bar(); }
};
/**
 * @constructor
 * @extends {Foo}
 */
function Bar() {
  this.xy = 3;
}
/** @return {Bar} */ function f() { return new Bar(); }
/** @return {Foo} */ function g() { return new Bar(); }
Bar.prototype = {
  /** @override @return {Bar} */ x: function() { new Bar(); },
  /** @override @return {Foo} */ y: function() { new Bar(); }
};
Bar.prototype.__proto__ = Foo.prototype;
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testIssue586()
func (c *TypeCheckBugsAndIssuesTest) TestIssue586() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */
var MyClass = function() {};
/** @param {boolean} success */
MyClass.prototype.fn = function(success) {};
MyClass.prototype.test = function() {
  this.fn();
  this.fn = function() {};
};
`).
		AddDiagnostic(`Function MyClass.prototype.fn: called with 0 argument(s). Function requires at least 1 argument(s) and no more than 1 argument(s).
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testIssue635()
func (c *TypeCheckBugsAndIssuesTest) TestIssue635() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */
function F() {}
F.prototype.bar = function() { this.baz(); };
F.prototype.baz = function() {};
/** @constructor */
function G() {}
G.prototype.bar = F.prototype.bar;
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testIssue635b()
func (c *TypeCheckBugsAndIssuesTest) TestIssue635b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */
function F() {}
/** @constructor */
function G() {}
/** @type {function(new:G)} */ var x = F;
`).
		AddDiagnostic(`initializing variable
found   : (typeof F)
required: function(new:G): ?
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testIssue669()
func (c *TypeCheckBugsAndIssuesTest) TestIssue669() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {{prop1: (Object|undefined)}} */
function f(a) {
  var results;
  if (a) {
    results = {};
    results.prop1 = {a: 3};
  } else {
    results = {prop2: 3};
  }
  return results;
}
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testIssue688()
func (c *TypeCheckBugsAndIssuesTest) TestIssue688() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @const */ var SOME_DEFAULT =
    /** @type {TwoNumbers} */ ({first: 1, second: 2});
/**
* Class defining an interface with two numbers.
* @interface
*/
function TwoNumbers() {}
/** @type {number} */
TwoNumbers.prototype.first;
/** @type {number} */
TwoNumbers.prototype.second;
/** @return {number} */ function f() { return SOME_DEFAULT; }
`).
		AddDiagnostic(`inconsistent return type
found   : (TwoNumbers|null)
required: number
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testIssue700()
func (c *TypeCheckBugsAndIssuesTest) TestIssue700() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {{text: string}} opt_data
 * @return {string}
 */
function temp1(opt_data) {
  return opt_data.text;
}

/**
 * @param {{activity: (boolean|number|string|null|Object)}} opt_data
 * @return {string}
 */
function temp2(opt_data) {
  /** @suppress {checkTypes} */
  function __inner() {
    return temp1(opt_data.activity);
  }
  return __inner();
}

/**
 * @param {{n: number, text: string, b: boolean}} opt_data
 * @return {string}
 */
function temp3(opt_data) {
  return 'n: ' + opt_data.n + ', t: ' + opt_data.text + '.';
}

function callee() {
  var output = temp3({
    n: 0,
    text: 'a string',
    b: true
  })
  alert(output);
}

callee();
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testIssue725()
func (c *TypeCheckBugsAndIssuesTest) TestIssue725() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @typedef {{name: string}} */ var RecordType1;
/** @typedef {{name2222: string}} */ var RecordType2;
/** @param {RecordType1} rec */ function f(rec) {
  alert(rec.name2222);
}
`).
		AddDiagnostic("Property name2222 never defined on rec").
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testIssue726()
func (c *TypeCheckBugsAndIssuesTest) TestIssue726() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
/** @param {number} x */ Foo.prototype.bar = function(x) {};
/** @return {!Function} */
Foo.prototype.getDeferredBar = function() {
  var self = this;
  return function() {
    self.bar(true);
  };
};
`).
		AddDiagnostic(`actual parameter 1 of Foo.prototype.bar does not match formal parameter
found   : boolean
required: number
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testIssue765()
func (c *TypeCheckBugsAndIssuesTest) TestIssue765() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */
var AnotherType = function(parent) {
    /** @param {string} stringParameter Description... */
    this.doSomething = function(stringParameter) {};
};
/** @constructor */
var YetAnotherType = function() {
    this.field = new AnotherType(self);
    this.testfun=function(stringdata) {
        this.field.doSomething(null);
    };
};
`).
		AddDiagnostic(`actual parameter 1 of AnotherType.doSomething does not match formal parameter
found   : null
required: string
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testIssue783()
func (c *TypeCheckBugsAndIssuesTest) TestIssue783() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */
var Type = function() {
  /** @type {Type} */
  this.me_ = this;
};
Type.prototype.doIt = function() {
  var me = this.me_;
  for (var i = 0; i < me.unknownProp; i++) {}
};
`).
		AddDiagnostic("Property unknownProp never defined on Type").
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testIssue791()
func (c *TypeCheckBugsAndIssuesTest) TestIssue791() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {{func: function()}} obj */
function test1(obj) {}
var fnStruc1 = {};
fnStruc1.func = function() {};
test1(fnStruc1);
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testIssue810()
func (c *TypeCheckBugsAndIssuesTest) TestIssue810() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */
var Type = function() {
  this.prop = x;
};
Type.prototype.doIt = function(obj) {
  this.prop = obj.unknownProp;
};
`).
		AddDiagnostic("Property unknownProp never defined on obj" + TypeCheck_POSSIBLE_INEXISTENT_PROPERTY_EXPLANATION).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testIssue1002()
func (c *TypeCheckBugsAndIssuesTest) TestIssue1002() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @interface */
var I = function() {};
/** @constructor @implements {I} */
var A = function() {};
/** @constructor @implements {I} */
var B = function() {};
var f = function() {
  if (A === B) {
    new B();
  }
};
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testIssue1023()
func (c *TypeCheckBugsAndIssuesTest) TestIssue1023() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */
function F() {}
(function() {
  F.prototype = {
    /** @param {string} x */
    bar: function(x) {  }
  };
})();
(new F()).bar(true)
`).
		AddDiagnostic(`actual parameter 1 of F.prototype.bar does not match formal parameter
found   : boolean
required: string
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testIssue1047()
func (c *TypeCheckBugsAndIssuesTest) TestIssue1047() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 */
function C2() {}

/**
 * @constructor
 */
function C3(c2) {
  /**
   * @type {C2} 
   * @private
   */
  this.c2_;

  var x = this.c2_.prop;
}
`).
		AddDiagnostic("Property prop never defined on C2").
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testIssue1056()
func (c *TypeCheckBugsAndIssuesTest) TestIssue1056() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {Array} */ var x = null;
x.push('hi');
`).
		AddDiagnostic(`No properties on this expression
found   : null
required: Object
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testIssue1072()
func (c *TypeCheckBugsAndIssuesTest) TestIssue1072() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {string} x
 * @return {number}
 */
var f1 = function(x) {
  return 3;
};

/** Function */
var f2 = function(x) {
  if (!x) throw new Error()
  return /** @type {number} */ (f1('x'))
}

/**
 * @param {string} x
 */
var f3 = function(x) {};

f1(f3);
`).
		AddDiagnostic(`actual parameter 1 of f1 does not match formal parameter
found   : function(string): undefined
required: string
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testIssue1123()
func (c *TypeCheckBugsAndIssuesTest) TestIssue1123() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {function(number)} g */ function f(g) {}
f(function(a, b) {})
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : function(?, ?): undefined
required: function(number): ?
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testIssue1201()
func (c *TypeCheckBugsAndIssuesTest) TestIssue1201() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {function(this:void)} f */ function g(f) {}
/** @constructor */ function F() {}
/** desc */ F.prototype.bar = function() {};
g(new F().bar);
`).
		AddDiagnostic(`actual parameter 1 of g does not match formal parameter
found   : function(this:F): undefined
required: function(this:undefined): ?
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testIssue1201b()
func (c *TypeCheckBugsAndIssuesTest) TestIssue1201b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {function(this:void)} f */ function g(f) {}
/** @constructor */ function F() {}
/** desc */ F.prototype.bar = function() {};
var f = new F();
g(f.bar.bind(f));
`).
		IncludeDefaultExterns().
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testIssue1201c()
func (c *TypeCheckBugsAndIssuesTest) TestIssue1201c() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {function(this:void)} f */ function g(f) {}
g(function() { this.alert() })
`).
		AddDiagnostic(`No properties on this expression
found   : undefined
required: Object
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testIssue926a()
func (c *TypeCheckBugsAndIssuesTest) TestIssue926a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** x */ function error() {}
/**
 * @constructor
 * @param {string} error
 */
function C(error) {
 /** @const */ this.e = error;
}
/** @type {number} */ var x = (new C('x')).e;
`).
		AddDiagnostic(`initializing variable
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testIssue926b()
func (c *TypeCheckBugsAndIssuesTest) TestIssue926b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */
function A() {
 /** @constructor */
 function B() {}
 /** @type {!B} */ this.foo = new B();
 /** @type {!B} */ var C = new B();
}
/** @type {number} */ var x = (new A()).foo;
`).
		AddDiagnostic(`initializing variable
found   : B
required: number
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testBug592170()
func (c *TypeCheckBugsAndIssuesTest) TestBug592170() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {Function} opt_f ... */
function foo(opt_f) {
  /** @type {Function} */
  return opt_f || function() {};
}
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testNullishCoalesceTypeIsFirstOrSecondArgument()
func (c *TypeCheckBugsAndIssuesTest) TestNullishCoalesceTypeIsFirstOrSecondArgument() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {Function} opt_f ... */
function foo(opt_f) {
  /** @type {Function} */
  return opt_f ?? function() {};
}
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testBug901455a()
func (c *TypeCheckBugsAndIssuesTest) TestBug901455a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {(number|undefined)} */ function a() { return 3; }
var b = undefined === a()
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testBug901455b()
func (c *TypeCheckBugsAndIssuesTest) TestBug901455b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {(number|undefined)} */ function a() { return 3; }
var b = a() === undefined
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testBug908701()
func (c *TypeCheckBugsAndIssuesTest) TestBug908701() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddString().Build()).
		AddSource(`/** @type {String} */ var s = new String('foo');
var b = s.match(/a/) != null;
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testBug908625()
func (c *TypeCheckBugsAndIssuesTest) TestBug908625() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */function A(){}
/** @constructor
 * @extends A */function B(){}
/** @param {B} b

 @return {(A|undefined)} */function foo(b){return b}
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testBug911118a()
func (c *TypeCheckBugsAndIssuesTest) TestBug911118a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource("var a = function() {}; /** @type {!null} */ const b = a;").
		AddDiagnostic(`initializing variable
found   : function(): undefined
required: None
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testBug911118b()
func (c *TypeCheckBugsAndIssuesTest) TestBug911118b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`function nullFunction() {};
var foo = nullFunction;
foo = function() {};
foo();
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testBug909000()
func (c *TypeCheckBugsAndIssuesTest) TestBug909000() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */function A(){}
/** @param {!A} a
@return {boolean}*/
function y(a) { return a }
`).
		AddDiagnostic(`inconsistent return type
found   : A
required: boolean
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testBug930117()
func (c *TypeCheckBugsAndIssuesTest) TestBug930117() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {boolean} x */function f(x){}
f(null);
`).
		AddDiagnostic(`actual parameter 1 of f does not match formal parameter
found   : null
required: boolean
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testBug1484445()
func (c *TypeCheckBugsAndIssuesTest) TestBug1484445() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor */ function Foo() {}
/** @type {number?} */ Foo.prototype.bar = null;
/** @type {number?} */ Foo.prototype.baz = null;
/** @param {Foo} foo */
function f(foo) {
  while (true) {
    if (foo.bar == null && foo.baz == null) {
      foo.bar;
    }
  }
}
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testBug1859535()
func (c *TypeCheckBugsAndIssuesTest) TestBug1859535() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @param {Function} childCtor Child class.
 * @param {Function} parentCtor Parent class.
 */
var inherits = function(childCtor, parentCtor) {
  /** @constructor */
  function tempCtor() {};
  tempCtor.prototype = parentCtor.prototype;
  childCtor.superClass_ = parentCtor.prototype;
  childCtor.prototype = new tempCtor();
  /** @override */ childCtor.prototype.constructor = childCtor;
};
/**
 * @param {Function} constructor
 * @param {Object} var_args
 * @return {Object}
 */
var factory = function(constructor, var_args) {
  /** @constructor */
  var tempCtor = function() {};
  tempCtor.prototype = constructor.prototype;
  var obj = new tempCtor();
  constructor.apply(obj, arguments);
  return obj;
};
`).
		IncludeDefaultExterns().
		Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testBug1940591()
func (c *TypeCheckBugsAndIssuesTest) TestBug1940591() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {Object} */
var a = {};
/** @type {number} */
a.name = 0;
/**
 * @param {Function} x anything.
 */
a.g = function(x) { x.name = 'a'; }
`).
		Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testBug1942972()
func (c *TypeCheckBugsAndIssuesTest) TestBug1942972() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`var google = {
  gears: {
    factory: {},
    workerPool: {}
  }
};

google.gears = {factory: {}};
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testBug1943776()
func (c *TypeCheckBugsAndIssuesTest) TestBug1943776() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return  {{foo: Array}} */
function bar() {
  return {foo: []};
}
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testBug1987544()
func (c *TypeCheckBugsAndIssuesTest) TestBug1987544() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {string} x */ function foo(x) {}
var duration;
if (true && !(duration = 3)) {
 foo(duration);
}
`).
		AddDiagnostic(`actual parameter 1 of foo does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testBug1940769()
func (c *TypeCheckBugsAndIssuesTest) TestBug1940769() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {!Object} */
function proto(obj) { return obj.prototype; }
/**
 * @constructor
 * @extends {Map}
 */
function Map2() { Map.call(this); };
Map2.prototype = proto(Map);
`).
		AddExterns(NewTestExternsBuilder().AddMap().AddObject().
			Build()).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testBug2335992()
func (c *TypeCheckBugsAndIssuesTest) TestBug2335992() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @return {*} */ function f() { return 3; }
var x = f();
/** @type {string} */
x.y = 3;
`).
		AddDiagnostic(`assignment
found   : number
required: string
`).
		Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testBug2341812()
func (c *TypeCheckBugsAndIssuesTest) TestBug2341812() {
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
		Suppress(DiagnosticGroups_STRICT_MISSING_PROPERTIES).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testBug7701884()
func (c *TypeCheckBugsAndIssuesTest) TestBug7701884() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(NewTestExternsBuilder().AddArray().Build()).
		AddSource(`/**
 * @param {Array<T>} x
 * @param {function(T)} y
 * @template T
 */
var forEach = function(x, y) {
  for (var i = 0; i < x.length; i++) y(x[i]);
};
/** @param {number} x */
function f(x) {}
/** @param {?} x */
function h(x) {
  var top = null;
  forEach(x, function(z) { top = z; });
  if (top) f(top);
}
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testBug8017789()
func (c *TypeCheckBugsAndIssuesTest) TestBug8017789() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {(map|function())} isResult */
var f = function(isResult) {
    while (true)
        isResult['t'];
};
/** @typedef {Object<string, number>} */
var map;
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testBug12441160()
func (c *TypeCheckBugsAndIssuesTest) TestBug12441160() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @param {string} a */ 
function use(a) {};
/**
 * @param {function(this:THIS)} fn
 * @param {THIS} context 
 * @constructor
 * @template THIS
 */
var P = function(fn, context) {}

/** @constructor */
function C() { /** @type {number} */ this.a = 1; }

/** @return {P} */ 
C.prototype.method = function() {
   return new P(function() { use(this.a); }, this);
};

`).
		AddDiagnostic(`actual parameter 1 of use does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testBug13641083a()
func (c *TypeCheckBugsAndIssuesTest) TestBug13641083a() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @constructor @struct */ function C() {};
new C().bar;
`).
		AddDiagnosticType(TypeCheck_INEXISTENT_PROPERTY).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testBug13641083b()
func (c *TypeCheckBugsAndIssuesTest) TestBug13641083b() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/** @type {?} */ var C;
C.bar + 1;
`).
		AddDiagnosticType(TypeCheck_POSSIBLE_INEXISTENT_PROPERTY).
		Run()
}

// Java: TypeCheckBugsAndIssuesTest.testBug12722936()
func (c *TypeCheckBugsAndIssuesTest) TestBug12722936() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddSource(`/**
 * @constructor
 * @template T
 */
function X() {}
/** @constructor */ function C() {
  /** @type {!X<boolean>}*/
  this.a = new X();
  /** @type {null} */ var x = this.a;
};

`).
		AddDiagnostic(`initializing variable
found   : X<boolean>
required: null
`).
		Run()
}
