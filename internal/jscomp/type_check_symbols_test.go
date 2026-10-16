// Port of com/google/javascript/jscomp/TypeCheckSymbolsTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/TypeCheckSymbolsTest.java)

package jscomp

import (
	"testing"
)

// TypeCheckSymbolsTest tests the type checker.
type TypeCheckSymbolsTest struct {
	typeCheckTests_noFixture
}

// Java: TypeCheckSymbolsTest()
func NewTypeCheckSymbolsTest() *TypeCheckSymbolsTest {
	c := &TypeCheckSymbolsTest{}
	return c
}

func TestTypeCheckSymbolsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewTypeCheckSymbolsTest() })
}

var TypeCheckSymbolsTest_SYMBOL_EXTERNS = NewTestExternsBuilder().AddIterable().AddArray().
	Build()

// Java: TypeCheckSymbolsTest.symbolMethod_class()
func (c *TypeCheckSymbolsTest) Test_symbolMethod_class() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckSymbolsTest_SYMBOL_EXTERNS).
		AddExterns(`/** @const {symbol} */
Symbol.foobar;
`).
		AddSource(`class Foo {
  /**
   * @param {string} x
   * @param {number} y
   */
  [Symbol.foobar](x, y) {}
}
new Foo()[Symbol.foobar](4, 5);
`).
		AddDiagnostic(`actual parameter 1 of function does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckSymbolsTest.symbolMethod_es5Constructor()
func (c *TypeCheckSymbolsTest) Test_symbolMethod_es5Constructor() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckSymbolsTest_SYMBOL_EXTERNS).
		AddExterns(`/** @const {symbol} */
Symbol.foobar;
`).
		AddSource(`/** @constructor */
function Foo() {}
/**
 * @param {string} x
 * @param {number} y
 */
Foo.prototype[Symbol.foobar] = function(x, y) {};
new Foo()[Symbol.foobar](4, 5);
`).
		AddDiagnostic(`actual parameter 1 of function does not match formal parameter
found   : number
required: string
`).
		Run()
}

// Java: TypeCheckSymbolsTest.symbolMethod_es5Constructor_invalidOverride()
func (c *TypeCheckSymbolsTest) Test_symbolMethod_es5Constructor_invalidOverride() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckSymbolsTest_SYMBOL_EXTERNS).
		AddExterns(`/** @const {symbol} */
Symbol.foobar;
`).
		AddSource(`/** @constructor */
function Parent() {}
/**
 * @param {string} x
 */
Parent.prototype[Symbol.foobar] = function(x) {}

/** @constructor @extends {Parent} */
function Child() {}
/**
 * @param {number} x
 * @override
 */
Child.prototype[Symbol.foobar] = function(x) {}
`).
		AddDiagnostic(`mismatch of the Symbol.foobar property type and the type of the property it overrides from superclass Parent
original: function(this:Parent, string): undefined
override: function(this:Child, number): undefined
`).
		Run()
}

// Java: TypeCheckSymbolsTest.symbolMethod_class_invalidOverride()
func (c *TypeCheckSymbolsTest) Test_symbolMethod_class_invalidOverride() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckSymbolsTest_SYMBOL_EXTERNS).
		AddExterns(`/** @const {symbol} */
Symbol.foobar;
`).
		AddSource(`class Parent {
  /**
   * @param {string} x
   */
  [Symbol.foobar](x) {}
}
class Child extends Parent {
  /**
   * @param {number} x
   * @override
   */
  [Symbol.foobar](x) {}
}
`).
		AddDiagnostic(`mismatch of the Symbol.foobar property type and the type of the property it overrides from superclass Parent
original: function(this:Parent, string): undefined
override: function(this:Child, number): undefined
`).
		Run()
}

// Java: TypeCheckSymbolsTest.symbolMethod_class_missingOverrideJsDoc()
func (c *TypeCheckSymbolsTest) Test_symbolMethod_class_missingOverrideJsDoc() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckSymbolsTest_SYMBOL_EXTERNS).
		AddExterns(`/** @const {symbol} */
Symbol.foobar;
`).
		AddSource(`class Parent {
  /** @param {string} x */
  [Symbol.foobar](x) {}
}
class Child extends Parent {
  [Symbol.foobar](x) {}
}
`).
		AddDiagnostic(`property Symbol.foobar already defined on superclass Parent; use @override to override it
`).
		Run()
}

// Java: TypeCheckSymbolsTest.symbolMethod_class_invalidImplements()
func (c *TypeCheckSymbolsTest) Test_symbolMethod_class_invalidImplements() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckSymbolsTest_SYMBOL_EXTERNS).
		AddExterns(`/** @const {symbol} */
Symbol.foobar;
`).
		AddSource(`/** @interface */
class Parent {
  /** @param {string} x */
  [Symbol.foobar](x) {}
}
/** @implements {Parent} */
class Child {
  /**
   * @param {number} x
   * @override
   */
  [Symbol.foobar](x) {}
}
`).
		AddDiagnostic(`mismatch of the Symbol.foobar property on type Child and the type of the property it overrides from interface Parent
original: function(this:Parent, string): undefined
override: function(this:Child, number): undefined
`).
		Run()
}

// Java: TypeCheckSymbolsTest.symbolMethod_class_implementsIterable_missingOverrideJSDoc()
func (c *TypeCheckSymbolsTest) Test_symbolMethod_class_implementsIterable_missingOverrideJSDoc() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckSymbolsTest_SYMBOL_EXTERNS).
		AddExterns(`/** @const {symbol} */
Symbol.foobar;
`).
		AddSource(`/** @implements {Iterable<string>} */
class I {
  /** @return {!Iterator<string>} */
  [Symbol.iterator]() {
    return ['foo', 'bar'][Symbol.iterator]();
  }
}
`).
		AddDiagnostic(`property Symbol.iterator already defined on interface Iterable; use @override to override it
`).
		Run()
}

// Java: TypeCheckSymbolsTest.symbolProp_assignmentTypeMismatch()
func (c *TypeCheckSymbolsTest) Test_symbolProp_assignmentTypeMismatch() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckSymbolsTest_SYMBOL_EXTERNS).
		AddExterns(`/** @const {symbol} */
Symbol.number;
`).
		AddSource(`class C {
  constructor() {
    /** @type {number} */
    this[Symbol.number] = 0;
  }
}
new C()[Symbol.number] = 'a string';
`).
		AddDiagnostic(`assignment
found   : string
required: number
`).
		Run()
}

// Java: TypeCheckSymbolsTest.symbolProp_onNamespace()
func (c *TypeCheckSymbolsTest) Test_symbolProp_onNamespace() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckSymbolsTest_SYMBOL_EXTERNS).
		AddExterns(`/** @const {symbol} */
Symbol.myProp;
`).
		AddSource(`const ns = {};
/** @type {string} */
ns[Symbol.myProp] = 'hello';
/** @type {null} */
const x = ns[Symbol.myProp];
`).
		AddDiagnostic(`initializing variable
found   : string
required: null
`).
		Run()
}

// Java: TypeCheckSymbolsTest.knownSymbolProp_onTypeWithIndexSignature()
func (c *TypeCheckSymbolsTest) Test_knownSymbolProp_onTypeWithIndexSignature() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckSymbolsTest_SYMBOL_EXTERNS).
		AddExterns(`/** @const {symbol} */
Symbol.myProp;
`).
		AddSource(`/** @type {!Object<symbol, number>} */
const o = {};
o[Symbol.myProp] = 'hello';
/** @type {null} */
const x = o[Symbol.myProp];
`).
		Run()
}

// Java: TypeCheckSymbolsTest.unknownSymbolPropAccess()
func (c *TypeCheckSymbolsTest) Test_unknownSymbolPropAccess() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckSymbolsTest_SYMBOL_EXTERNS).
		AddExterns(`/** @const {symbol} */
Symbol.good;
/** @const {symbol} */
Symbol.bad;
`).
		AddSource(`class C {
  constructor() {
    /** @type {number} */
    this[Symbol.good] = 0;
  }
}
use(new C()[Symbol.good]);
use(new C()[Symbol.bad]);
`).
		Run()
}

// Java: TypeCheckSymbolsTest.objLitWithComputedProp_checkedAgainstStructuralType()
func (c *TypeCheckSymbolsTest) Test_objLitWithComputedProp_checkedAgainstStructuralType() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckSymbolsTest_SYMBOL_EXTERNS).
		AddExterns(`/** @const {symbol} */
Symbol.foobar;
`).
		AddSource(`/** @record */
class Foo {
  constructor() {
    /** @type {string} */
    this[Symbol.foobar];
  }
}
/** @type {!Foo} */
const foo = {
  /** @const */
  [Symbol.foobar]: 123
};
`).
		AddDiagnostic(`initializing variable
found   : {[Symbol.foobar]: number}
required: Foo
`).
		Run()
}

// Java: TypeCheckSymbolsTest.objLitWithComputedProp_checkedAgainstStructuralType_forMissingRequiredProp()
func (c *TypeCheckSymbolsTest) Test_objLitWithComputedProp_checkedAgainstStructuralType_forMissingRequiredProp() {
	TypeCheckTestCase_TypeTestBuilder_NewTest().AddExterns(TypeCheckSymbolsTest_SYMBOL_EXTERNS).
		AddExterns(`/** @const {symbol} */
Symbol.foobar;
`).
		AddSource(`/** @record */
class Foo {
  constructor() {
    /** @type {string} */
    this[Symbol.foobar];
  }
}
/** @type {!Foo} */
const foo = {};
`).
		AddDiagnostic(`initializing variable
found   : {}
required: Foo
`).
		Run()
}
