// Port of com/google/javascript/jscomp/ijs/ConvertToTypedInterfaceTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/ijs/ConvertToTypedInterfaceTest.java)

/*
 * Copyright 2015 The Closure Compiler Authors.
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

// Unit tests for ConvertToTypedInterface.
type ConvertToTypedInterfaceTest struct {
	CompilerTestCaseBase
}

func NewConvertToTypedInterfaceTest() *ConvertToTypedInterfaceTest {
	c := &ConvertToTypedInterfaceTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestConvertToTypedInterfaceTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewConvertToTypedInterfaceTest() })
}

// Java: ConvertToTypedInterfaceTest.setUp()
func (c *ConvertToTypedInterfaceTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.AllowExternsChanges()
}

// Java: ConvertToTypedInterfaceTest.getProcessor(Compiler)
func (c *ConvertToTypedInterfaceTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewConvertToTypedInterface(compiler)
}

// Java: ConvertToTypedInterfaceTest.getOptions()
func (c *ConvertToTypedInterfaceTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetWarningLevel(DiagnosticGroups_MODULE_LOAD, CheckLevel_OFF)
	return options
}

// Java: ConvertToTypedInterfaceTest.testInferAnnotatedTypeFromTypeInference()
func (c *ConvertToTypedInterfaceTest) TestInferAnnotatedTypeFromTypeInference() {
	c.Test("/** @const */ var x = 5;", "/** @const {number} */ var x;")

	c.Test(
		"/** @constructor */ function Foo() { /** @const */ this.x = 5; }",
		"/** @constructor */ function Foo() {} \n /** @const {number} */ Foo.prototype.x;")
}

// Java: ConvertToTypedInterfaceTest.testExternsDefinitionsRespected()
func (c *ConvertToTypedInterfaceTest) TestExternsDefinitionsRespected() {
	c.Test(c.Externs("/** @type {number} */ var x;"), c.Srcs("x = 7;"), c.Expected(""))
}

// Java: ConvertToTypedInterfaceTest.testUnannotatedDeclaration()
func (c *ConvertToTypedInterfaceTest) TestUnannotatedDeclaration() {
	c.Test("var x;", "/** @const {UnusableType} */ var x;")
}

// Java: ConvertToTypedInterfaceTest.testEmptyClass()
func (c *ConvertToTypedInterfaceTest) TestEmptyClass() {
	c.Test("class x {;}", "class x {}")
}

// Java: ConvertToTypedInterfaceTest.testComputedFieldNoRHS()
func (c *ConvertToTypedInterfaceTest) TestComputedFieldNoRHS() {
	c.Test(
		`const s = Symbol();
class Foo {
  [s]
}
`,
		`const s = Symbol();
class Foo {
}
`)
	c.Test(
		`class Foo {
  [Symbol.iterator]
}
`,
		`class Foo {
  [Symbol.iterator]
}
`)
}

// Java: ConvertToTypedInterfaceTest.testDoubleAssignmentField()
func (c *ConvertToTypedInterfaceTest) TestDoubleAssignmentField() {
	c.Test(
		`let Snackbar = Snackbar1 = class Snackbar {
  x = 5
}
`,
		"/** @const @type {UnusableType} */ var Snackbar;")
}

// Java: ConvertToTypedInterfaceTest.testDoubleAssignmentStaticField()
func (c *ConvertToTypedInterfaceTest) TestDoubleAssignmentStaticField() {
	c.Test(
		`let A = B = class C {
  static y = true
}
`,
		"/** @const @type {UnusableType} */ var A;")
}

// Java: ConvertToTypedInterfaceTest.testSuperClassFields()
func (c *ConvertToTypedInterfaceTest) TestSuperClassFields() {
	c.Test(
		`class First{
  /** @const {number} */ fromFirst;
  constructor() {
    this.fromFirst = 1;
  };
};
class Second extends First {
   /** @const {number} */ fromSecond;
  constructor() {
    super();
    /** @override */this.fromFirst = 7;
    this.fromSecond = 5;
  };
};
`,
		`class First {
  /** @const {number} */ fromFirst;
  constructor() {
  }
}
class Second extends First {
   /** @const {number} */ fromSecond;
  constructor() {
  }
}
 /** @override */ Second.prototype.fromFirst;
`)
}

// Java: ConvertToTypedInterfaceTest.testThisAssignment()
func (c *ConvertToTypedInterfaceTest) TestThisAssignment() {
	// When a public field is created and assigned with a constructor, we only care about the
	// declaration.
	c.Test(
		`class Foo {
  /** @const {number} */ unique;
  constructor() {
    this.unique = 5;
   }
}
`,
		`class Foo {
  /** @const {number}*/ unique;
  constructor() {
  }
}
`)
	c.Test(
		`class Foo {
  /** @const {number} */ fieldOne;
  /** @const {boolean} */ fieldTwo;
  /** @const {string} */ fieldThree;
  constructor() {
    this.fieldOne = 1;
    this.fieldTwo = true
    this.fieldThree = 'three'
   }
}
`,
		`class Foo {
  /** @const {number} */ fieldOne;
  /** @const {boolean} */ fieldTwo;
  /** @const {string} */ fieldThree;
  constructor() {
  }
}
`)

	c.Test(
		`class Foo {
  /** @const {number} */ fieldOne;
  /** @const {boolean} */ fieldTwo;
  /** @const {string} */ fieldThree;
  constructor(x,y,z) {
    this.fieldOne = x;
    this.fieldTwo = y
    this.fieldThree = z
   }
}
`,
		`class Foo {
  /** @const {number} */ fieldOne;
  /** @const {boolean} */ fieldTwo;
  /** @const {string} */ fieldThree;
  constructor(x,y,z) {
  }
}
`)
}

// Java: ConvertToTypedInterfaceTest.testPrototypeDeclared()
func (c *ConvertToTypedInterfaceTest) TestPrototypeDeclared() {
	// When the prototype is declared, we only care about the initial declaration for public fields.
	c.Test(
		`class Foo {
  /** @const {number} */sameField = 10;
}
Foo.prototype.sameField = 8;
`,
		`class Foo {
  /** @const {number} */ sameField;
}
`)
}

// Java: ConvertToTypedInterfaceTest.testThisAndPrototype()
func (c *ConvertToTypedInterfaceTest) TestThisAndPrototype() {
	c.Test(
		`class Foo {
  /** @const @type {number} */sameField = 10;
  constructor() {
    this.sameField = 5;
   }
}
Foo.prototype.sameField = 8;
`,
		`class Foo {
  /** @const @type {number} */ sameField;
  constructor() {
  }
}
`)
}

// Java: ConvertToTypedInterfaceTest.testComputedFieldDef()
func (c *ConvertToTypedInterfaceTest) TestComputedFieldDef() {
	// Computed Field Def's are unanalyzable to the compiler so they will be dropped.
	c.Test(
		"const PREFIX = 'prefix';\n/** @unrestricted*/ class Foo {\n  [`${PREFIX}Field`] = 'prefixed field';\n}\n",
		` /** @const @type {string} */
var PREFIX;
/** @unrestricted*/ class Foo {
}
`)
	c.Test(
		`/** @unrestricted*/ class Foo {
  [Math.random()] = 'gone';
}
`,
		`/** @unrestricted*/ class Foo {
}
`)
}

// Java: ConvertToTypedInterfaceTest.testJSDocMemberFieldDef()
func (c *ConvertToTypedInterfaceTest) TestJSDocMemberFieldDef() {
	// Add type if not present
	c.Test(
		`class Foo {
  /** @const  */ myField = 5;
}
`,
		`class Foo {
  /** @const {number} */ myField;
}
`)
	// If no JSdoc present, set to UnusableType
	c.Test(
		`class Foo {
  myField = 5;
}
`,
		`class Foo {
  /** @const @type {UnusableType} */ myField;
}
`)
}

// Java: ConvertToTypedInterfaceTest.testStaticMemberFieldDef()
func (c *ConvertToTypedInterfaceTest) TestStaticMemberFieldDef() {
	// check that static remains in the declaration even when the assignment is removed
	c.Test(
		`class Foo {
  /** @const @type {number}*/ static myField = 5;
}
`,
		`class Foo {
  /** @const @type {number}*/ static myField;
}
`)
}

// Java: ConvertToTypedInterfaceTest.testUpdatingStaticField()
func (c *ConvertToTypedInterfaceTest) TestUpdatingStaticField() {
	// Updating the static field should not be included in i.js
	c.Test(
		`class Foo {
  /** @const {number}*/static myField = 1;
}
Foo.myField = 2;
`,
		`class Foo {
  /** @const {number}*/static myField;
}
`)
}

// Java: ConvertToTypedInterfaceTest.testThisInFieldInitializer()
func (c *ConvertToTypedInterfaceTest) TestThisInFieldInitializer() {
	c.Test(
		`class Base {
  /** @const @type {string}*/ baseField = 'base field';
  /** @const @type {string}*/anotherBaseField = this.baseField;
}
`,
		`class Base {
  /** @const @type {string}*/ baseField;
  /** @const @type {string}*/ anotherBaseField;
}
`)
	c.Test(
		`class Base {
  /** @const @type {string}*/ static fieldOne = 'I am static';
  /** @const @type {string}*/ static fieldTwo = this.fieldOne;
}
`,
		`class Base {
  /** @const @type {string}*/ static fieldOne;
  /** @const @type {string}*/ static fieldTwo;
}
`)

	c.Test(
		`class Base {
  /** @const @type {string}*/ static fieldStatic = 'I am static';
  /** @const @type {string}*/ instanceField = this.fieldStatic;
}
`,
		`class Base {
  /** @const @type {string}*/ static fieldStatic;
  /** @const @type {string}*/ instanceField;
}
`)
}

// Java: ConvertToTypedInterfaceTest.testSameFieldNameInDifferentClass()
func (c *ConvertToTypedInterfaceTest) TestSameFieldNameInDifferentClass() {
	c.Test(
		`class Base {
  /** @const @type {string}*/ sameField;
}
class Foo {
  constructor() {
    /** @const {number}*/ this.sameField;
  }
}
`,
		`class Base {
   /** @const @type {string}*/ sameField;
}
class Foo {
  constructor(){
  }
}
/** @const {number}*/ Foo.prototype.sameField
`)
}

// Java: ConvertToTypedInterfaceTest.testSameFieldDeclared()
func (c *ConvertToTypedInterfaceTest) TestSameFieldDeclared() {
	// If the one of the duplicated fields contains type annotation, we will keep that one.
	c.Test(
		`class Base {
  sameField;
  /** @const @type {string}*/ sameField = 'first';
}
`,
		`class Base {
  /** @const @type {string}*/ sameField
}
`)

	c.Test(
		`class Base {
  /** @const @type {string}*/ sameField = 'first';
}
class Foo {
  /** @const {number}*/ sameField = 3;
}
`,
		`class Base {
  /** @const @type {string}*/ sameField
}
class Foo {
 /** @const @type {number}*/ sameField
}
`)
}

// Java: ConvertToTypedInterfaceTest.testSimpleConstJsdocPropagation()
func (c *ConvertToTypedInterfaceTest) TestSimpleConstJsdocPropagation() {
	c.Test("/** @const */ var x = 5;", "/** @const {number} */ var x;")
	c.Test("/** @const */ var x = 5n;", "/** @const {bigint} */ var x;")
	c.Test("/** @const */ var x = true;", "/** @const {boolean} */ var x;")
	c.Test("/** @const */ var x = 'str';", "/** @const {string} */ var x;")
	c.Test("/** @const */ var x = `str`;", "/** @const {string} */ var x;")
	c.Test("/** @const */ var x = null;", "/** @const {null} */ var x;")
	c.Test("/** @const */ var x = void 0;", "/** @const {void} */ var x;")
	c.Test("/** @const */ var x = /a/;", "/** @const {!RegExp} */ var x;")

	c.Test(
		"/** @constructor */ function Foo() { /** @const */ this.x = 5; }",
		"/** @constructor */ function Foo() {} \n /** @const {number} */ Foo.prototype.x;")

	c.Test(
		"/** @const */ var x = cond ? true : 5;",
		"/** @const {UnusableType} */ var x;",
		c.Warning(ConvertToTypedInterface_CONSTANT_WITHOUT_EXPLICIT_TYPE))
}

// Java: ConvertToTypedInterfaceTest.testConstKeywordWithAnnotatedType()
func (c *ConvertToTypedInterfaceTest) TestConstKeywordWithAnnotatedType() {
	c.Test("/** @type {number} */ const x = 5;", "/** @const {number} */ var x;")
	c.Test("/** @type {!Foo} */ const f = new Foo;", "/** @const {!Foo} */ var f;")
}

// Java: ConvertToTypedInterfaceTest.testConstKeywordJsdocPropagation()
func (c *ConvertToTypedInterfaceTest) TestConstKeywordJsdocPropagation() {
	c.Test("const x = 5;", "/** @const {number} */ var x;")

	c.Test(
		"const x = 5, y = 'str', z = /abc/;",
		`/** @const {number} */ var x;
/** @const {string} */ var y;
/** @const {!RegExp} */ var z;
`)

	c.Test(
		"const x = cond ? true : 5;",
		"/** @const {UnusableType} */ var x;",
		c.Warning(ConvertToTypedInterface_CONSTANT_WITHOUT_EXPLICIT_TYPE))
}

// Java: ConvertToTypedInterfaceTest.testPropagateConstCast()
func (c *ConvertToTypedInterfaceTest) TestPropagateConstCast() {
	c.Test("const x = /** @type {!Array<number>} */ ([]);", "/** @const {!Array<number>} */ var x;")

	c.Test(
		"let /** (string|number) */ x = /** @type {number} */ (5);",
		"/** @type {(string|number)} */ var x;")
}

// Java: ConvertToTypedInterfaceTest.testSplitMultiDeclarations()
func (c *ConvertToTypedInterfaceTest) TestSplitMultiDeclarations() {
	c.Test(
		"var /** number */ x = 4, /** string */ y = 'str';",
		"/** @type {number} */ var x; /** @type {string} */ var y;")

	c.Test(
		"var /** number */ x, /** string */ y;",
		"/** @type {number} */ var x; /** @type {string} */ var y;")

	c.Test(
		"let /** number */ x = 4, /** string */ y = 'str';",
		"/** @type {number} */ var x; /** @type {string} */ var y;")

	c.Test(
		"let /** number */ x, /** string */ y;",
		"/** @type {number} */ let x; /** @type {string} */ let y;")
}

// Java: ConvertToTypedInterfaceTest.testThisPropertiesInConstructors()
func (c *ConvertToTypedInterfaceTest) TestThisPropertiesInConstructors() {
	c.Test(
		"/** @constructor */ function Foo() { /** @const {number} */ this.x; }",
		`/** @constructor */ function Foo() {}
/** @const {number} */ Foo.prototype.x
`)

	c.Test(
		"/** @constructor */ function Foo() { this.x = undefined; }",
		`/** @constructor */ function Foo() {}
/** @const {UnusableType} */ Foo.prototype.x;
`)

	c.Test(
		"/** @constructor */ function Foo() { /** @type {?number} */ this.x = null; this.x = 5; }",
		`/** @constructor */ function Foo() {}
/** @type {?number} */ Foo.prototype.x;
`)

	c.Test(
		"/** @constructor */ function Foo() { /** @const */ this.x = cond ? true : 5; }",
		`/** @constructor */ function Foo() {} /** @const {UnusableType} */ Foo.prototype.x;
`,
		c.Warning(ConvertToTypedInterface_CONSTANT_WITHOUT_EXPLICIT_TYPE))
}

// Java: ConvertToTypedInterfaceTest.testNonThisPropertiesInConstructors()
func (c *ConvertToTypedInterfaceTest) TestNonThisPropertiesInConstructors() {
	c.Test(
		"/** @constructor */ function Foo() { const obj = {}; obj.name = () => 5; alert(obj); }",
		"/** @constructor */ function Foo() {}")
}

// Java: ConvertToTypedInterfaceTest.testThisPropertiesInConstructorsAndPrototype()
func (c *ConvertToTypedInterfaceTest) TestThisPropertiesInConstructorsAndPrototype() {
	c.Test(
		`/** @constructor */ function Foo() { this.x = null; }
/** @type {?number} */ Foo.prototype.x = 5;
`,
		`/** @constructor */ function Foo() {}
/** @type {?number} */ Foo.prototype.x;
`)

	c.Test(
		`/** @constructor */ function Foo() { this.x = null; }
/** @type {?number} */ Foo.prototype.x;
`,
		`/** @constructor */ function Foo() {}
/** @type {?number} */ Foo.prototype.x;
`)

	c.Test(
		`/** @constructor */ function Foo() { this.x = null; }
Foo.prototype.x = 5;
`,
		`/** @constructor */ function Foo() {}
/** @const {UnusableType} */ Foo.prototype.x;
`)
}

// Java: ConvertToTypedInterfaceTest.testConstJsdocPropagationForGlobalNames()
func (c *ConvertToTypedInterfaceTest) TestConstJsdocPropagationForGlobalNames() {
	c.Test(
		"/** @type {!Array<string>} */ var x = []; /** @const */ var y = x;",
		"/** @type {!Array<string>} */ var x; /** @const */ var y = x;")

	c.Test(
		"/** @type {Object} */ var o = {}; /** @type {number} */ o.p = 5; /** @const */ var y = o;",
		"/** @type {Object} */ var o; /** @type {number} */ o.p; /** @const */ var y = o;")
}

// Java: ConvertToTypedInterfaceTest.testConstJsdocPropagationForConstructorNames()
func (c *ConvertToTypedInterfaceTest) TestConstJsdocPropagationForConstructorNames() {
	c.Test(
		`/** @constructor */
function Foo(/** number */ x) {
  /** @const */ this.x = x;
}
`,
		`/** @constructor */ function Foo(/** number */ x) {}
/** @const {number} */ Foo.prototype.x;
`)

	c.Test(
		`/** @constructor @param {!Array<string>} arr */
function Foo(arr) {
  /** @const */ this.arr = arr;
}
`,
		`/** @constructor @param {!Array<string>} arr */ function Foo(arr) {}
/** @const {!Array<string>} */ Foo.prototype.arr;
`)

	c.Test(
		`class Foo {
  constructor(/** number */ x) {
    /** @const */ this.x = x;
  }
}
`,
		`class Foo {
  constructor(/** number */ x) {}
}
/** @const {number} */ Foo.prototype.x;
`)

	c.Test(
		`class Foo {
  /** @param {number} x */
  constructor(x) {
    /** @const */ this.x = x;
  }
}
`,
		`class Foo {
  /** @param {number} x */
  constructor(x) {}
}
/** @const {number} */ Foo.prototype.x;
`)
}

// Java: ConvertToTypedInterfaceTest.testClassMethodsConflictWithOtherAssignment()
func (c *ConvertToTypedInterfaceTest) TestClassMethodsConflictWithOtherAssignment() {
	c.Test(
		`class Foo {
  /** @return {number} */
  method() {}
}
Foo.prototype.method = wrap(Foo.prototype.method);
`,
		`class Foo {
  /** @return {number} */
  method() {}
}
`)

	c.Test(
		`class Foo {
  constructor() {
    this.method = wrap(Foo.prototype.method);
  }
  /** @return {number} */
  method() {}
}
`,
		`class Foo {
  constructor() {}
  /** @return {number} */
  method() {}
}
`)
}

// Java: ConvertToTypedInterfaceTest.testMultipleSameNamedThisProperties()
func (c *ConvertToTypedInterfaceTest) TestMultipleSameNamedThisProperties() {
	c.Test(
		`class Foo {
  constructor(/** number */ x) {
    /** @const */ this.x = x;
  }
}
/** @template T */
class Bar {
  constructor(/** T */ x) {
    /** @const */ this.x = x;
  }
}
`,
		`class Foo {
  constructor(/** number */ x) {}
}
/** @const {number} */ Foo.prototype.x;
/** @template T */
class Bar {
  constructor(/** T */ x) {}
}
/** @const {T} */ Bar.prototype.x;
`)
}

// Java: ConvertToTypedInterfaceTest.testGlobalThis()
func (c *ConvertToTypedInterfaceTest) TestGlobalThis() {
	c.TestSame("/** @const */ this.globalNamespace = {};")
	c.Test(
		"/** @const */ this.globalNamespace = this.globalNamespace || {};",
		"/** @const */ this.globalNamespace = {};")
}

// Java: ConvertToTypedInterfaceTest.testGoogAddSingletonGetter()
func (c *ConvertToTypedInterfaceTest) TestGoogAddSingletonGetter() {
	c.TestSame("class Foo {}  goog.addSingletonGetter(Foo);")
}

// Java: ConvertToTypedInterfaceTest.testLegacyGoogModule()
func (c *ConvertToTypedInterfaceTest) TestLegacyGoogModule() {
	c.TestSame(
		`goog.module('a.b.c');
goog.module.declareLegacyNamespace();

exports = class {};
`)
}

// Java: ConvertToTypedInterfaceTest.testExternsAlias()
func (c *ConvertToTypedInterfaceTest) TestExternsAlias() {
	c.TestSame("const winAlias = window;")
	c.TestSame("const winAlias = window; const locationAlias = winAlias.location;")
}

// Java: ConvertToTypedInterfaceTest.testConstructorAlias1()
func (c *ConvertToTypedInterfaceTest) TestConstructorAlias1() {
	c.TestSame(
		`/** @constructor */
function Foo() {}
/** @const */ var FooAlias = Foo;
`)
}

// Java: ConvertToTypedInterfaceTest.testConstructorAlias2()
func (c *ConvertToTypedInterfaceTest) TestConstructorAlias2() {
	c.TestSame(
		`goog.module('a.b.c');

/** @constructor */
function Foo() {}
/** @const */ var FooAlias = Foo;
`)
}

// Java: ConvertToTypedInterfaceTest.testConstructorAlias3()
func (c *ConvertToTypedInterfaceTest) TestConstructorAlias3() {
	c.TestSame(
		`class Foo {}
/** @const */ var FooAlias = Foo;
`)
}

// Java: ConvertToTypedInterfaceTest.testConstructorAlias4()
func (c *ConvertToTypedInterfaceTest) TestConstructorAlias4() {
	c.TestSame(
		`goog.module('a.b.c');

class Foo {}
/** @constructor */ var FooAlias = Foo;
`)
}

// Java: ConvertToTypedInterfaceTest.testConstructorAlias5()
func (c *ConvertToTypedInterfaceTest) TestConstructorAlias5() {
	c.TestSame(
		`/** @constructor */
function Foo() {}
/** @constructor */ var FooAlias = Foo;
`)
}

// Java: ConvertToTypedInterfaceTest.testConstructorAlias6()
func (c *ConvertToTypedInterfaceTest) TestConstructorAlias6() {
	c.TestSame(
		`goog.provide('a.b.c.Foo');
goog.provide('FooAlias');

/** @constructor */
a.b.c.Foo = function() {};

/** @const */ var FooAlias = a.b.c.Foo;
`)
}

// Java: ConvertToTypedInterfaceTest.testConstructorAlias7()
func (c *ConvertToTypedInterfaceTest) TestConstructorAlias7() {
	c.TestSame(
		`class Foo {}
const FooAlias = Foo;
`)
}

// Java: ConvertToTypedInterfaceTest.testConstructorAlias8()
func (c *ConvertToTypedInterfaceTest) TestConstructorAlias8() {
	c.TestSame(
		`goog.module('a.b.c');

class Foo {}
const FooAlias = Foo;
`)
}

// Java: ConvertToTypedInterfaceTest.testRequireAlias1()
func (c *ConvertToTypedInterfaceTest) TestRequireAlias1() {
	c.TestSame(
		`goog.provide('FooAlias');

goog.require('a.b.c.Foo');

/** @const */ var FooAlias = a.b.c.Foo;
`)
}

// Java: ConvertToTypedInterfaceTest.testRequireAlias2()
func (c *ConvertToTypedInterfaceTest) TestRequireAlias2() {
	c.TestSame(
		`goog.provide('FooAlias');
goog.provide('BarAlias');

goog.require('a.b.c');

/** @const */ var FooAlias = a.b.c.Foo;
/** @const */ var BarAlias = a.b.c.Bar;
`)
}

// Java: ConvertToTypedInterfaceTest.testRequireAlias3()
func (c *ConvertToTypedInterfaceTest) TestRequireAlias3() {
	c.TestSame(
		`goog.module('FooAlias');

const Foo = goog.require('a.b.c.Foo');

exports = Foo;
`)
}

// Java: ConvertToTypedInterfaceTest.testRequireAlias4()
func (c *ConvertToTypedInterfaceTest) TestRequireAlias4() {
	c.TestSame(
		`goog.module('FooAlias');

const {Foo} = goog.require('a.b.c');

exports = Foo;
`)
}

// Java: ConvertToTypedInterfaceTest.testRequireAlias5()
func (c *ConvertToTypedInterfaceTest) TestRequireAlias5() {
	c.Test(
		`goog.module('FooAlias');

const {Foo} = goog.require('a.b.c');
/** @const {number} */
Foo = 1;

exports = Foo;
`,
		`goog.module('FooAlias');

/** @const @type {number} */ var Foo;

exports = Foo;
`)
}

// Java: ConvertToTypedInterfaceTest.testRequireTypeAlias1()
func (c *ConvertToTypedInterfaceTest) TestRequireTypeAlias1() {
	c.TestSame(
		`goog.provide('FooAlias');

goog.requireType('a.b.c.Foo');

/** @const */ var FooAlias = a.b.c.Foo;
`)
}

// Java: ConvertToTypedInterfaceTest.testRequireTypeAlias2()
func (c *ConvertToTypedInterfaceTest) TestRequireTypeAlias2() {
	c.TestSame(
		`goog.provide('FooAlias');
goog.provide('BarAlias');

goog.requireType('a.b.c');

/** @const */ var FooAlias = a.b.c.Foo;
/** @const */ var BarAlias = a.b.c.Bar;
`)
}

// Java: ConvertToTypedInterfaceTest.testRequireTypeAlias3()
func (c *ConvertToTypedInterfaceTest) TestRequireTypeAlias3() {
	c.TestSame(
		`goog.module('FooAlias');

const Foo = goog.requireType('a.b.c.Foo');

exports = Foo;
`)
}

// Java: ConvertToTypedInterfaceTest.testRequireTypeAlias4()
func (c *ConvertToTypedInterfaceTest) TestRequireTypeAlias4() {
	c.TestSame(
		`goog.module('FooAlias');

const {Foo} = goog.requireType('a.b.c');

exports = Foo;
`)
}

// Java: ConvertToTypedInterfaceTest.testRequireTypeAlias5()
func (c *ConvertToTypedInterfaceTest) TestRequireTypeAlias5() {
	c.Test(
		`goog.module('FooAlias');

const {Foo, Bar} = goog.requireType('a.b.c');
/** @const {string} */
Foo = 'foo';

exports = Foo;
`,
		`goog.module('FooAlias');

const {Bar} = goog.requireType('a.b.c');
/** @const @type {string} */ var Foo;

exports = Foo;
`)
}

// Java: ConvertToTypedInterfaceTest.testDestructuredAlias1()
func (c *ConvertToTypedInterfaceTest) TestDestructuredAlias1() {
	c.TestSame("const {Foo} = a.b.c; const {Bar} = x; exports = Foo;")
}

// Java: ConvertToTypedInterfaceTest.testDestructuredAlias2()
func (c *ConvertToTypedInterfaceTest) TestDestructuredAlias2() {
	c.TestSame("const {Foo, Bar} = a.b.c;  exports = Foo;")
}

// Java: ConvertToTypedInterfaceTest.testDestructuredAlias3()
func (c *ConvertToTypedInterfaceTest) TestDestructuredAlias3() {
	c.Test(
		`const {Foo} = a.b.c, bar = 1;
exports = Foo;
`,
		`const {Foo} = a.b.c;

/** @const @type {number} */ var bar;

exports = Foo;
`)
}

// Java: ConvertToTypedInterfaceTest.testDestructuredAlias4()
func (c *ConvertToTypedInterfaceTest) TestDestructuredAlias4() {
	c.Test(
		"const {Foo} = a.b.c, Baz = p.q.Baz; exports = Foo;",
		"const {Foo} = a.b.c; const Baz = p.q.Baz; exports = Foo;")
}

// Java: ConvertToTypedInterfaceTest.testDuplicateDeclarationWAliasRemoved1()
func (c *ConvertToTypedInterfaceTest) TestDuplicateDeclarationWAliasRemoved1() {
	c.Test(
		`const {Foo, Bar} = x.y;
/** @const {number} */
Bar = 1;

exports = Foo;
`,
		`const {Foo} = x.y

/** @const @type {number} */ var Bar;

exports = Foo;
`)
}

// Java: ConvertToTypedInterfaceTest.testDuplicateDeclarationWAliasRemoved2()
func (c *ConvertToTypedInterfaceTest) TestDuplicateDeclarationWAliasRemoved2() {
	c.Test(
		`/** @const {string} */
Foo = 'hello';

const {Foo, Bar} = x.y;

/** @const {number} */
Bar = 1;

exports = Foo;
`,
		`/** @const @type {string} */ var Foo;

/** @const @type {number} */ var Bar;

exports = Foo;
`)
}

// Java: ConvertToTypedInterfaceTest.testDuplicateDeclarationWAliasRemoved3()
func (c *ConvertToTypedInterfaceTest) TestDuplicateDeclarationWAliasRemoved3() {
	c.Test(
		`/** @const {string} */
Foo = 'hello';

const {Foo, Bar, Baz} = x.y;

/** @const {number} */
Baz = 1;

exports = Foo;
`,
		`/** @const @type {string} */ var Foo;

const {Bar} = x.y;

/** @const @type {number} */ var Baz;

exports = Foo;
`)
}

// Java: ConvertToTypedInterfaceTest.testBreakDownDestructuringForConst()
func (c *ConvertToTypedInterfaceTest) TestBreakDownDestructuringForConst() {
	c.Test(
		`goog.module('m');
const {x,y}=fn();
exports.x=x;
exports.y=y;
`,
		`goog.module("m");
/** @const @type {UnusableType} */
var x;
/** @const @type {UnusableType} */
var y;
exports.x=x;
exports.y = y;
`,
		c.Warning(ConvertToTypedInterface_CONSTANT_WITHOUT_EXPLICIT_TYPE))

	c.Test(
		`goog.module('m');
const {x=0} = fn();
exports.x = x;
`,
		`goog.module("m");
/** @const @type {UnusableType} */
var x;
exports.x = x;
`,
		c.Warning(ConvertToTypedInterface_CONSTANT_WITHOUT_EXPLICIT_TYPE))

	c.Test(
		`goog.module('m');
const {x:{y}} = fn();
exports.x = x;
`,
		`goog.module("m");
/** @const @type {UnusableType} */
var y;
exports.x = x;
`,
		c.Warning(ConvertToTypedInterface_CONSTANT_WITHOUT_EXPLICIT_TYPE))

	c.Test(
		`goog.module('m');
const {[x()]:x} = fn();
exports.x = x;
`,
		`goog.module("m");
/** @const @type {UnusableType} */
var x;
exports.x = x;
`,
		c.Warning(ConvertToTypedInterface_CONSTANT_WITHOUT_EXPLICIT_TYPE))
}

// Java: ConvertToTypedInterfaceTest.testBreakDownDestructuringWarning()
func (c *ConvertToTypedInterfaceTest) TestBreakDownDestructuringWarning() {
	c.TestWarning(
		`goog.module('m');
const x = ({y: obj.y} = fn());
exports.x = x;
`,
		ConvertToTypedInterface_CONSTANT_WITHOUT_EXPLICIT_TYPE)
}

// Java: ConvertToTypedInterfaceTest.testDuplicateDeclarationWAliasNotRemoved()
func (c *ConvertToTypedInterfaceTest) TestDuplicateDeclarationWAliasNotRemoved() {
	c.Test(
		`const {Foo, Bar} = x.y;

Bar = z;

exports = Foo;
`,
		`const {Foo, Bar} = x.y;

exports = Foo;
`)
}

// Java: ConvertToTypedInterfaceTest.testBreakDownDestructuringForLet()
func (c *ConvertToTypedInterfaceTest) TestBreakDownDestructuringForLet() {
	c.Test(
		`let {Foo, Bar} = a.b.c;
exports = Foo;
`,
		`/** @const @type {UnusableType} */
var Foo
/** @const @type {UnusableType} */
var Bar;
exports = Foo;
`,
		c.Warning(ConvertToTypedInterface_CONSTANT_WITHOUT_EXPLICIT_TYPE))
}

// Java: ConvertToTypedInterfaceTest.testAtConstAnnotationAlias()
func (c *ConvertToTypedInterfaceTest) TestAtConstAnnotationAlias() {
	c.Test(
		"/** @const */ var x = a.b.c; var y = x;",
		"/** @const */ var x = a.b.c; /** @const @type {UnusableType} */ var y;")
}

// Java: ConvertToTypedInterfaceTest.testConstPropagationPrivateProperties1()
func (c *ConvertToTypedInterfaceTest) TestConstPropagationPrivateProperties1() {
	c.Test(
		`/** @constructor */
function Foo() {
  /** @const @private */ this.x = someComplicatedExpression();
}
`,
		`/** @constructor */ function Foo() {}
/** @const @private {UnusableType} */ Foo.prototype.x;
`)
}

// Java: ConvertToTypedInterfaceTest.testConstPropagationPrivateProperties2()
func (c *ConvertToTypedInterfaceTest) TestConstPropagationPrivateProperties2() {
	c.Test(
		`goog.provide('a.b.c');

/** @private @const */
a.b.c.helper_ = someComplicatedExpression();
`,
		"goog.provide('a.b.c');   /** @private @const {UnusableType} */ a.b.c.helper_;")
}

// Java: ConvertToTypedInterfaceTest.testPrivateEs6ClassPropertyTypeOmitted_forNonTypeScript_doesNotOmit()
func (c *ConvertToTypedInterfaceTest) TestPrivateEs6ClassPropertyTypeOmitted_forNonTypeScript_doesNotOmit() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				"bar.js",
				`class Foo {
  constructor() {
    /** @private {SomeComplicatedType} */
    this.x;
  }
}
`)),
		c.Expected(
			`class Foo {
  constructor() {}
}
/** @private @type {SomeComplicatedType} */
Foo.prototype.x;
`))
}

// Java: ConvertToTypedInterfaceTest.testPrivateEs6ClassPropertyTypeOmitted_forTypeScriptFiles()
func (c *ConvertToTypedInterfaceTest) TestPrivateEs6ClassPropertyTypeOmitted_forTypeScriptFiles() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				"foo.closure.js",
				`class Foo {
  constructor() {
    /** @private {SomeComplicatedType} */
    this.x;
  }
}
`)),
		c.Expected(
			`class Foo {
  constructor() {}
}
/** @private @type {UnusableType} */
Foo.prototype.x;
`))
}

// Java: ConvertToTypedInterfaceTest.testNonPrivateVisibility_kept()
func (c *ConvertToTypedInterfaceTest) TestNonPrivateVisibility_kept() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				"foo.closure.js",
				`class Foo {
  constructor() {
    /** @protected {SomeComplicatedType} */
    this.x;
    /** @package {SomeComplicatedType} */
    this.y;
    /** @public {SomeComplicatedType} */
    this.z;
  }
}
`)),
		c.Expected(
			`class Foo {
  constructor() {}
}
/** @public @type {SomeComplicatedType} */
Foo.prototype.z;
/** @package @type {SomeComplicatedType} */
Foo.prototype.y;
/** @protected @type {SomeComplicatedType} */
Foo.prototype.x;
`))
}

// Java: ConvertToTypedInterfaceTest.testPrivateEs6ClassPropertyWithQmarkType_doesNotBecomeUnusableType()
func (c *ConvertToTypedInterfaceTest) TestPrivateEs6ClassPropertyWithQmarkType_doesNotBecomeUnusableType() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				"foo.closure.js",
				`class Foo {
  constructor() {
    /** @private {?} */
    this.x;
  }
}
`)),
		c.Expected(
			`class Foo {
  constructor() {}
}
/** @private @type {?} */
Foo.prototype.x;
`))
}

// Java: ConvertToTypedInterfaceTest.testPrivateEs5ClassPropertyTypeOmitted()
func (c *ConvertToTypedInterfaceTest) TestPrivateEs5ClassPropertyTypeOmitted() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				"foo.closure.js",
				`/** @constructor */
function Foo() {
  /** @const @private {SomeComplicatedType} */ this.x;
}
`)),
		c.Expected(
			`/** @constructor */ function Foo() {}
/** @const @private {UnusableType} */ Foo.prototype.x;
`))
}

// Java: ConvertToTypedInterfaceTest.testPrivateEs6PrototypeMethodSignatureOmitted()
func (c *ConvertToTypedInterfaceTest) TestPrivateEs6PrototypeMethodSignatureOmitted() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				"foo.closure.js",
				`class Foo {
  /**
   * @param {SomeType} a
   * @param {SomeType=} b
   * @param {...SomeType} c
   * @return {OtherType}
   * @private
   */
  x(a, b) {}
}
`)),
		c.Expected(
			`class Foo {
  /**
   * @param {IjsNoneType} a
   * @param {IjsNoneType=} b
   * @param {...IjsNoneType} c
   * @return {UnusableType}
   * @private
   */
  x(a, b) {}
}
`))
}

// Java: ConvertToTypedInterfaceTest.testPrivateEs5PrototypeMethodSignatureOmitted()
func (c *ConvertToTypedInterfaceTest) TestPrivateEs5PrototypeMethodSignatureOmitted() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				"foo.closure.js",
				`/** @constructor */
function Foo() {}
/**
 * @param {SomeType} a
 * @param {SomeType=} b
 * @return {OtherType}
 * @private
 */
Foo.prototype.x = function(a, b) {};
`)),
		c.Expected(
			`/** @constructor */ function Foo() {}
/**
 * @private
 * @param {IjsNoneType} a
 * @param {IjsNoneType=} b
 * @return {UnusableType}
 */
Foo.prototype.x = function(a, b) {};
`))
}

// Java: ConvertToTypedInterfaceTest.testDontEmitTypeNameOnPrivate_nonProperty()
func (c *ConvertToTypedInterfaceTest) TestDontEmitTypeNameOnPrivate_nonProperty() {
	c.Test(
		c.Srcs(
			SourceFile_FromCode(
				"foo.closure.js",
				`goog.module('a');

/** @private @typedef {string} */
let Typedef;
/** @private @enum {string} */
const Enum = {A: 0};
/** @private {boolean} */
const unusedPrivateLocal = true;
`)),
		c.Expected(
			`goog.module("a");

/** @private @typedef {string} */
let Typedef;
/** @private @enum {string} */
const Enum = {A: 0};
/** @private @const @type {boolean} */
var unusedPrivateLocal;
`))
}

// Java: ConvertToTypedInterfaceTest.testOverrideAnnotationCountsAsDeclaration()
func (c *ConvertToTypedInterfaceTest) TestOverrideAnnotationCountsAsDeclaration() {
	c.TestSame(
		`goog.provide('x.y.z.Bar');

goog.require('a.b.c.Foo');

/** @constructor @extends {a.b.c.Foo} */
x.y.z.Bar = function() {};

/** @override */
x.y.z.Bar.prototype.method = function(a, b, c) {};
`)

	c.TestSame(
		`goog.module('x.y.z');

const {Foo} = goog.require('a.b.c');

/** @constructor @extends {Foo} */
const Bar = class {
   /** @override */
   method(a, b, c) {}
}
exports.Bar = Bar;
`)
}

// Java: ConvertToTypedInterfaceTest.testConstJsdocPropagationForNames_optional()
func (c *ConvertToTypedInterfaceTest) TestConstJsdocPropagationForNames_optional() {
	c.Test(
		`/** @constructor */
function Foo(/** number= */ opt_x) {
  /** @const */ this.x = opt_x;
}
`,
		`/** @constructor */ function Foo(/** number= */ opt_x) {}
/** @const {number|undefined} */ Foo.prototype.x;
`)
}

// Java: ConvertToTypedInterfaceTest.testNotConfusedByOutOfOrderDeclarations()
func (c *ConvertToTypedInterfaceTest) TestNotConfusedByOutOfOrderDeclarations() {
	c.Test(
		`/** @constructor */
function Foo(/** boolean= */ opt_tag) {
  if (opt_tag) {
    Foo.tag = opt_tag;
  }
}
/** @type {boolean} */ Foo.tag = true;
`,
		`/** @constructor */
function Foo(/** boolean= */ opt_tag) {}
/** @type {boolean} */ Foo.tag;
`)
}

// Java: ConvertToTypedInterfaceTest.testConstJsdocPropagationForNames_rest()
func (c *ConvertToTypedInterfaceTest) TestConstJsdocPropagationForNames_rest() {
	c.Test(
		`/**
 * @constructor
 * @param {...number} nums
 */
function Foo(...nums) {
  /** @const */ this.nums = nums;
}
`,
		`/**
 * @constructor
 * @param {...number} nums
 */
function Foo(...nums) {}
/** @const {!Array<number>} */ Foo.prototype.nums;
`)
}

// Java: ConvertToTypedInterfaceTest.testOptionalRestParamFunction()
func (c *ConvertToTypedInterfaceTest) TestOptionalRestParamFunction() {
	c.Test(
		`/**
 * @param {?Object} o
 * @param {string=} str
 * @param {number=} num
 * @param {...!Object} rest
 */
function foo(o, str = '', num = 5, ...rest) {}
`,
		`/**
 * @param {?Object} o
 * @param {string=} str
 * @param {number=} num
 * @param {...!Object} rest
 */
function foo(o, str=void 0, num=void 0, ...rest) {}
`)
}

// Java: ConvertToTypedInterfaceTest.testDefaultValuesRemain()
func (c *ConvertToTypedInterfaceTest) TestDefaultValuesRemain() {
	c.Test("function f(x = 0) {}", "function f(x = void 0) {}")
	c.Test("function f(x = window.foobar()) {}", "function f(x = void 0) {}")
}

// Java: ConvertToTypedInterfaceTest.testGoogCallerLocationDefaultValueIsPreserved()
func (c *ConvertToTypedInterfaceTest) TestGoogCallerLocationDefaultValueIsPreserved() {
	c.TestSame("function f(a = goog.callerLocation()) {}")
	c.Test(
		"function f(a = 1, b = goog.callerLocation()) {}",
		"function f(a = void 0, b = goog.callerLocation()) {}")
}

// Java: ConvertToTypedInterfaceTest.testConstJsdocPropagationForNames_defaultValue()
func (c *ConvertToTypedInterfaceTest) TestConstJsdocPropagationForNames_defaultValue() {
	c.Test(
		`/**
 * @constructor
 * @param {string=} str
 */
function Foo(str = '') {
  /** @const */ this.s = str;
}
`,
		`/**
 * @constructor
 * @param {string=} str
 */
function Foo(str = void 0) {}
/** @const {string} */ Foo.prototype.s;
`)

	c.Test(
		`class Foo {
  /** @param {string=} str */
  constructor(str = '') {
    /** @const */ this.s = str;
  }
}
`,
		`class Foo {
  /** @param {string=} str */
  constructor(str = void 0) {}
}
/** @const {string} */ Foo.prototype.s;
`)
}

// Java: ConvertToTypedInterfaceTest.testConstWithDeclaredTypes()
func (c *ConvertToTypedInterfaceTest) TestConstWithDeclaredTypes() {
	c.Test("/** @const @type {number} */ var n = compute();", "/** @const @type {number} */ var n;")
	c.Test("/** @const {number} */ var n = compute();", "/** @const @type {number} */ var n;")
	c.Test("/** @const @return {void} */ var f = compute();", "/** @const @return {void} */ var f;")
	c.Test("/** @const @this {Array} */ var f = compute();", "/** @const @this {Array} x */ var f;")

	c.Test(
		"/** @const @param {number} x */ var f = compute();",
		"/** @const @param {number} x */ var f;")

	c.Test(
		"/** @const @constructor x */ var Foo = createConstructor();",
		"/** @const @constructor x */ var Foo;")
}

// Java: ConvertToTypedInterfaceTest.testRemoveUselessStatements()
func (c *ConvertToTypedInterfaceTest) TestRemoveUselessStatements() {
	c.Test("34", "")
	c.Test("'str'", "")
	c.Test("({x:4})", "")
	c.Test("debugger;", "")
	c.Test("throw 'error';", "")
	c.Test("label: debugger;", "")
}

// Java: ConvertToTypedInterfaceTest.testRemoveUnnecessaryBodies()
func (c *ConvertToTypedInterfaceTest) TestRemoveUnnecessaryBodies() {
	c.Test("function f(x,y) { /** @type {number} */ z = x + y; return z; }", "function f(x,y) {}")

	c.Test(
		"/** @return {number} */ function f(/** number */ x, /** number */ y) { return x + y; }",
		"/** @return {number} */ function f(/** number */ x, /** number */ y) {}")

	c.Test(
		"class Foo { method(/** string */ s) { return s.split(','); } }",
		"class Foo { method(/** string */ s) {} }")
}

// Java: ConvertToTypedInterfaceTest.testRemoveEmptyMembers()
func (c *ConvertToTypedInterfaceTest) TestRemoveEmptyMembers() {
	c.Test(
		"class Foo { ;; method(/** string */ s) {};; }",
		"class Foo { method(/** string */ s) {} }")
}

// Java: ConvertToTypedInterfaceTest.testEs6Modules()
func (c *ConvertToTypedInterfaceTest) TestEs6Modules() {
	c.TestSame("export default class {}")

	c.TestSame("import Foo from '/foo';")

	c.TestSame("export class Foo {}")

	c.TestSame("import {Foo} from '/foo';")

	c.TestSame(
		`import {Baz} from '/baz';

export /** @constructor */ function Foo() {}
/** @type {!Baz} */ Foo.prototype.baz
`)

	c.TestSame(
		`import {Bar, Baz} from '/a/b/c';

class Foo extends Bar {
  /** @return {!Baz} */ getBaz() {}
}

export {Foo};
`)

	c.Test(
		`export class Foo {
  /** @return {number} */ getTime() { return Date.now(); }
}
export default /** @return {number} */ () => 6
const BLAH = 'foobar';
export {BLAH};
`,
		`export class Foo {
  /** @return {number} */ getTime() {}
}
export default /** @return {number} */ () => {}
/** @const {string} */ var BLAH;
export {BLAH};
`)
}

// Java: ConvertToTypedInterfaceTest.testEs6ModulesExportedNameDeclarations()
func (c *ConvertToTypedInterfaceTest) TestEs6ModulesExportedNameDeclarations() {
	c.TestSame("/** @type {number} */ export let x;")
	c.TestSame("/** @type {number} */ export var x;")

	c.Test("/** @type {number} */ export var x = 5;", "/** @type {number} */ export var x;")
	c.Test("/** @type {number} */ export const X = 5;", "/** @const {number} */ export var X;")
	// TODO(blickly): Ideally, we would leave let declarations alone
	c.Test("/** @type {number} */ export let x = 5;", "/** @type {number} */ export var x;")

	c.Test(
		"/** @type {number} */ export const X = 5, Y = z;",
		"/** @const {number} */ export var X; /** @const {number} */ export var Y;")
}

// Java: ConvertToTypedInterfaceTest.testEs6ModulesDeclareModuleId()
func (c *ConvertToTypedInterfaceTest) TestEs6ModulesDeclareModuleId() {
	c.TestSame(
		`goog.declareModuleId('foo');
/** @type {number} */ export var x;
`)
}

// Java: ConvertToTypedInterfaceTest.testGoogModules()
func (c *ConvertToTypedInterfaceTest) TestGoogModules() {
	c.TestSame(
		`goog.module('x.y.z');

/** @constructor */ function Foo() {}

exports = Foo;
`)

	c.TestSame(
		`goog.module('x.y.z');

const Baz = goog.require('a.b.c');

/** @constructor */ function Foo() {}
/** @type {!Baz} */ Foo.prototype.baz

exports = Foo;
`)

	c.TestSame(
		`goog.module('x.y.z');

const {Bar, Baz} = goog.require('a.b.c');

/** @constructor */ function Foo() {}
/** @type {!Baz} */ Foo.prototype.baz

exports = Foo;
`)

	c.TestSame(
		c.Srcs(
			`goog.module('a.b.c');
/** @constructor */ function Foo() {}
Foo.prototype.display = function() {};
exports = Foo;
`,
			`goog.module('x.y.z');
/** @constructor */ function Foo() {}
Foo.prototype.display = function() {};
exports = Foo;
`))

	c.TestSame(
		c.Srcs(
			`/** @constructor */ function Foo() {}
Foo.prototype.display = function() {};
`,
			`goog.module('x.y.z');
/** @constructor */ function Foo() {}
Foo.prototype.display = function() {};
exports = Foo;
`))

	c.Test(
		c.Srcs(
			`goog.module('a.b.c');
/** @constructor */ function Foo() {
  /** @type {number} */ this.x = 5;
}
exports = Foo;
`,
			`goog.module('x.y.z');
/** @constructor */ function Foo() {
  /** @type {number} */ this.x = 99;
}
exports = Foo;
`),
		c.Expected(
			`goog.module('a.b.c');
/** @constructor */ function Foo() {}
/** @type {number} */ Foo.prototype.x;
exports = Foo;
`,
			`goog.module('x.y.z');
/** @constructor */ function Foo() {}
/** @type {number} */ Foo.prototype.x;
exports = Foo;
`))
}

// Java: ConvertToTypedInterfaceTest.testGoogModulesWithTypedefExports()
func (c *ConvertToTypedInterfaceTest) TestGoogModulesWithTypedefExports() {
	c.TestSame(
		`goog.module('x.y.z');

/** @typedef {number} */
exports.Foo;
`)
}

// Java: ConvertToTypedInterfaceTest.testGoogModulesWithUndefinedExports()
func (c *ConvertToTypedInterfaceTest) TestGoogModulesWithUndefinedExports() {
	c.TestWarning(
		`goog.module('x.y.z');

const Baz = goog.require('x.y.z.Baz');
const Foobar = goog.require('f.b.Foobar');

exports = (new Baz).getFoobar();
`,
		ConvertToTypedInterface_CONSTANT_WITHOUT_EXPLICIT_TYPE)

	c.TestWarning(
		`goog.module('x.y.z');

const Baz = goog.require('x.y.z.Baz');
const Foobar = goog.require('f.b.Foobar');

exports.foobar = (new Baz).getFoobar();
`,
		ConvertToTypedInterface_CONSTANT_WITHOUT_EXPLICIT_TYPE)
}

// Java: ConvertToTypedInterfaceTest.testGoogModulesWithAnnotatedUninferrableExports()
func (c *ConvertToTypedInterfaceTest) TestGoogModulesWithAnnotatedUninferrableExports() {
	c.Test(
		`goog.module('x.y.z');

const Baz = goog.require('x.y.z.Baz');
const Foobar = goog.require('f.b.Foobar');

/** @type {Foobar} */
exports.foobar = (new Baz).getFoobar();
`,
		`goog.module('x.y.z');

const Baz = goog.require('x.y.z.Baz');
const Foobar = goog.require('f.b.Foobar');

/** @type {Foobar} */
exports.foobar;
`)

	c.Test(
		`goog.module('x.y.z');

const Baz = goog.require('x.y.z.Baz');
const Foobar = goog.require('f.b.Foobar');

/** @type {Foobar} */
exports = (new Baz).getFoobar();
`,
		`goog.module('x.y.z');

const Baz = goog.require('x.y.z.Baz');
const Foobar = goog.require('f.b.Foobar');

/** @type {Foobar} */
exports = /** @type {?} */ (0);
`)
}

// Java: ConvertToTypedInterfaceTest.testCrossFileModifications()
func (c *ConvertToTypedInterfaceTest) TestCrossFileModifications() {
	c.Test(
		`goog.module('a.b.c');
othermodule.modify.something = othermodule.modify.something + 1;
`,
		"goog.module('a.b.c');")

	c.Test(
		`goog.module('a.b.c');
class Foo {}
Foo.something = othermodule.modify.something + 1;
exports = Foo;
`,
		`goog.module('a.b.c');
class Foo {}
/** @const {UnusableType} */ Foo.something
exports = Foo;
`)

	c.Test(
		`goog.provide('a.b.c');
otherfile.modify.something = otherfile.modify.something + 1;
`,
		"goog.provide('a.b.c');")

	c.Test(
		`goog.provide('a.b.c');
a.b.c.something = otherfile.modify.something + 1;
`,
		`goog.provide('a.b.c');
/** @const {UnusableType} */ a.b.c.something;
`)
}

// Java: ConvertToTypedInterfaceTest.testRemoveCalls()
func (c *ConvertToTypedInterfaceTest) TestRemoveCalls() {
	c.Test("alert('hello'); window.clearTimeout();", "")

	c.TestSame("goog.provide('a.b.c');")

	c.TestSame("goog.provide('a.b.c'); goog.require('x.y.z');")
}

// Java: ConvertToTypedInterfaceTest.testEnums()
func (c *ConvertToTypedInterfaceTest) TestEnums() {
	c.Test(
		"/** @const @enum {number} */ var E = { A: 1, B: 2, C: 3};",
		"/** @const @enum {number} */ var E = { A: 0, B: 0, C: 0};")

	c.Test(
		"/** @enum {number} */ var E = { A: foo(), B: bar(), C: baz()};",
		"/** @enum {number} */ var E = { A: 0, B: 0, C: 0};")

	// No corrosion into 0 for string enums. This allows conformance to emit fewer false positives.
	// Values longer than 10 gets truncated.
	c.Test(
		"/** @enum {string} */ var E = { A: 'hello', B: 'world', C: 'some_string_longer_than_10'};",
		"/** @enum {string} */ var E = { A: 'hello', B: 'world', C: 'some_str..'};")

	c.Test(
		"/** @enum {Object} */ var E = { A: {b: 'c'}, D: {e: 'f'} };",
		"/** @enum {Object} */ var E = { A: 0, D: 0};")

	c.Test(
		"var x = 7; /** @enum {number} */ var E = { A: x };",
		"/** @const {UnusableType} */ var x; /** @enum {number} */ var E = { A: 0 };")
}

// Java: ConvertToTypedInterfaceTest.testEnumInsideNamespace()
func (c *ConvertToTypedInterfaceTest) TestEnumInsideNamespace() {
	c.Test(
		"const ns = { /** @enum {number} */ ENUM: { A: 1, B: 2, C: 3} };",
		"const ns = { /** @enum {number} */ ENUM: { A: 0, B: 0, C: 0} };")
}

// Java: ConvertToTypedInterfaceTest.testTryCatch()
func (c *ConvertToTypedInterfaceTest) TestTryCatch() {
	c.Test(
		"try { /** @type {number} */ var n = foo(); } catch (e) { console.log(e); }",
		"/** @type {number} */ var n;")

	c.Test(
		`try {
  /** @type {number} */ var start = Date.now();
  doStuff();
} finally {
  /** @type {number} */ var end = Date.now();
}
`,
		"/** @type {number} */ var start; /** @type {number} */ var end;")
}

// Java: ConvertToTypedInterfaceTest.testTemplatedClass()
func (c *ConvertToTypedInterfaceTest) TestTemplatedClass() {
	c.Test(
		`/** @template T */
class Foo {
  /** @param {T} x */
  constructor(x) { /** @const */ this.x = x;}
}
`,
		`/** @template T */
class Foo {
  /** @param {T} x */
  constructor(x) {}
}
/** @const {T} */ Foo.prototype.x;
`)
}

// Java: ConvertToTypedInterfaceTest.testConstructorBodyWithThisDeclaration()
func (c *ConvertToTypedInterfaceTest) TestConstructorBodyWithThisDeclaration() {
	c.Test(
		"/** @constructor */ function Foo() { /** @type {number} */ this.num = 5;}",
		"/** @constructor */ function Foo() {} /** @type {number} */ Foo.prototype.num;")

	c.Test(
		"/** @constructor */ function Foo(b) { if (b) { /** @type {number} */ this.num = 5; } }",
		"/** @constructor */ function Foo(b) {} /** @type {number} */ Foo.prototype.num;")

	c.Test(
		"/** @constructor */ let Foo = function() { /** @type {number} */ this.num = 5;}",
		"/** @constructor */ let Foo = function() {}; /** @type {number} */ Foo.prototype.num;")

	c.Test(
		"class Foo { constructor() { /** @type {number} */ this.num = 5;} }",
		"class Foo { constructor() {} } /** @type {number} */ Foo.prototype.num;")
}

// Java: ConvertToTypedInterfaceTest.testConstructorBodyWithoutThisDeclaration()
func (c *ConvertToTypedInterfaceTest) TestConstructorBodyWithoutThisDeclaration() {
	c.Test(
		"/** @constructor */ function Foo(o) { o.num = 8; var x = 'str'; }",
		"/** @constructor */ function Foo(o) {}")
}

// Java: ConvertToTypedInterfaceTest.testIIFE()
func (c *ConvertToTypedInterfaceTest) TestIIFE() {
	c.Test("(function(){ /** @type {number} */ var n = 99; })();", "")
}

// Java: ConvertToTypedInterfaceTest.testConstants()
func (c *ConvertToTypedInterfaceTest) TestConstants() {
	c.Test("/** @const {number} */ var x = 5;", "/** @const {number} */ var x;")
}

// Java: ConvertToTypedInterfaceTest.testDefines()
func (c *ConvertToTypedInterfaceTest) TestDefines() {
	// NOTE: This is another pattern that is only allowed in externs.
	c.Test(
		"/** @define {number} */ var x = 5;", //
		"/** @define {number} */ var x;")

	c.Test(
		"/** @define {number} */ goog.define('goog.BLAH', 5);",
		"/** @define {number} */ goog.define('goog.BLAH', 0);")

	c.Test(
		"/** @define {string} */ const BLAH = goog.define('goog.BLAH', 'blah');",
		"/** @define {string} */ const BLAH = goog.define('goog.BLAH', '');")

	c.Test(
		"/** @define {boolean} */ goog.BLECH = goog.define('goog.BLAH', true);",
		"/** @define {boolean} */ goog.BLECH = goog.define('goog.BLAH', false);")

	c.Test(
		"/** @define {number|boolean} */ const X = goog.define('goog.XYZ', true);",
		"/** @define {number|boolean} */ const X = goog.define('goog.XYZ', 0);")
}

// Java: ConvertToTypedInterfaceTest.testNestedBlocks()
func (c *ConvertToTypedInterfaceTest) TestNestedBlocks() {
	c.Test("{ const x = foobar(); }", "")

	c.Test("{ /** @const */ let x = foobar(); }", "")

	c.Test("{ /** @const */ let x = foobar(); x = foobaz(); }", "")

	c.TestWarning(
		"{ /** @const */ var x = foobar(); }",
		ConvertToTypedInterface_CONSTANT_WITHOUT_EXPLICIT_TYPE)
}

// Java: ConvertToTypedInterfaceTest.testGoogProvidedTopLevelSymbol()
func (c *ConvertToTypedInterfaceTest) TestGoogProvidedTopLevelSymbol() {
	c.TestSame("goog.provide('Foo');  /** @constructor */ Foo = function() {};")
}

// Java: ConvertToTypedInterfaceTest.testIfs()
func (c *ConvertToTypedInterfaceTest) TestIfs() {
	c.Test("if (true) { var /** number */ x = 5; }", "/** @type {number} */ var x;")

	c.Test(
		"if (true) { var /** number */ x = 5; } else { var /** string */ y = 'str'; }",
		"/** @type {number} */ var x; /** @type {string} */ var  y;")

	c.Test("if (true) { if (false) { var /** number */ x = 5; } }", "/** @type {number} */ var x;")

	c.Test(
		"if (true) {} else { if (false) {} else { var /** number */ x = 5; } }",
		"/** @type {number} */ var x;")
}

// Java: ConvertToTypedInterfaceTest.testLoops()
func (c *ConvertToTypedInterfaceTest) TestLoops() {
	c.Test("while (true) { foo(); break; }", "")

	c.Test(
		"for (var i = 0; i < 10; i++) { var field = 88; }",
		"/** @const {UnusableType} */ var i; /** @const {UnusableType} */ var field;")

	c.Test(
		"for (var i = 0, arraySize = getSize(); i < arraySize; i++) { foo(arr[i]); }",
		"/** @const {UnusableType} */ var i; /** @const {UnusableType} */ var arraySize;")

	c.Test("while (i++ < 10) { var /** number */ field = i; }", "/** @type {number} */ var field;")

	c.Test(
		"label: while (i++ < 10) { var /** number */ field = i; }",
		"/** @type {number} */ var field;")

	c.Test(
		"do { var /** number */ field = i; } while (i++ < 10);",
		"/** @type {number} */ var field;")

	c.Test(
		"for (var /** number */ i = 0; i < 10; i++) { var /** number */ field = i; }",
		"/** @type {number} */ var i; /** @type {number} */ var field;")

	c.Test(
		"for (i = 0; i < 10; i++) { var /** number */ field = i; }",
		"/** @type {number} */ var field;")

	c.Test(
		"for (var i = 0; i < 10; i++) { var /** number */ field = i; }",
		"/** @const {UnusableType} */ var i; /** @type {number} */ var field;")
}

// Java: ConvertToTypedInterfaceTest.testSymbols()
func (c *ConvertToTypedInterfaceTest) TestSymbols() {
	c.TestSame("const sym = Symbol();")

	c.TestSame("/** @const */ var sym = Symbol();")

	c.Test("const sym = Symbol(computeDescription());", "const sym = Symbol();")

	c.Test(
		"/** @type {symbol} */ var sym = Symbol.for(computeDescription());",
		"/** @type {symbol} */ var sym;")
}

// Java: ConvertToTypedInterfaceTest.testNamespaces()
func (c *ConvertToTypedInterfaceTest) TestNamespaces() {
	c.TestSame("/** @const */ var ns = {}; /** @return {number} */ ns.fun = function(x,y,z) {}")

	c.TestSame("/** @const */ var ns = {}; ns.fun = function(x,y,z) {}")

	c.Test(
		"/** @const */ var ns = ns || {}; ns.fun = function(x,y,z) {}",
		"/** @const */ var ns = {}; ns.fun = function(x,y,z) {}")
}

// Java: ConvertToTypedInterfaceTest.testNonemptyNamespaces()
func (c *ConvertToTypedInterfaceTest) TestNonemptyNamespaces() {
	c.TestSame("/** @const */ var ns = {fun: function(x,y,z) {}}")

	c.Test(
		"/** @const */ var ns = {/** @type {number} */ n: 5};",
		"/** @const */ var ns = {/** @type {number} */ n: 0};")

	// NOTE: This pattern typechecks when found in externs, but not for code.
	// Since the goal of this pass is intended to be used as externs, this is acceptable.
	c.Test(
		"/** @const */ var ns = {/** @type {string} */ s: 'str'};",
		"/** @const */ var ns = {/** @type {string} */ s: 0};")

	c.Test(
		"/** @const */ var ns = {/** @const */ s: 'blahblahblah'};",
		"/** @const */ var ns = {/** @const {string} */ s: 0};")

	c.Test(
		"/** @const */ var ns = {untyped: foo()};",
		"/** @const */ var ns = {/** @const {UnusableType} */ untyped: 0};")
}

// Java: ConvertToTypedInterfaceTest.testConstKeywordNamespaces()
func (c *ConvertToTypedInterfaceTest) TestConstKeywordNamespaces() {
	c.TestSame("const ns = {}; /** @return {number} */ ns.fun = function(x,y,z) {}")
	c.TestSame("const ns = { /** @return {number} */ fun : goog.abstractMethod };")
	c.TestSame("const ns = {fun: function(x,y,z) {}}")
	c.TestSame("const ns = { /** @return {number} */ fun(x,y,z) {}}")
}

// Java: ConvertToTypedInterfaceTest.testRemoveIgnoredProperties()
func (c *ConvertToTypedInterfaceTest) TestRemoveIgnoredProperties() {
	c.Test(
		"/** @const */ var ns = {}; /** @return {number} */ ns['fun'] = function(x,y,z) {}",
		"/** @const */ var ns = {};")

	c.Test(
		"/** @constructor */ function Foo() {} Foo.prototype['fun'] = function(x,y,z) {}",
		"/** @constructor */ function Foo() {}")

	c.Test(
		"/** @constructor */ function Foo() {} /** @type {str} */ Foo['prototype'].method;",
		"/** @constructor */ function Foo() {}")
}

// Java: ConvertToTypedInterfaceTest.testRemoveRepeatedProperties()
func (c *ConvertToTypedInterfaceTest) TestRemoveRepeatedProperties() {
	c.Test(
		"/** @const */ var ns = {}; /** @type {number} */ ns.x = 5; ns.x = 7;",
		"/** @const */ var ns = {}; /** @type {number} */ ns.x;")

	c.Test(
		"/** @const */ var ns = {}; ns.x = 5; ns.x = 7;",
		"/** @const */ var ns = {}; /** @const {UnusableType} */ ns.x;")

	c.Test(
		"const ns = {}; /** @type {number} */ ns.x = 5; ns.x = 7;",
		"const ns = {}; /** @type {number} */ ns.x;")
}

// Java: ConvertToTypedInterfaceTest.testRemoveRepeatedDeclarations()
func (c *ConvertToTypedInterfaceTest) TestRemoveRepeatedDeclarations() {
	c.Test("/** @type {number} */ var x = 4; var x = 7;", "/** @type {number} */ var x;")

	c.Test("/** @type {number} */ var x = 4; x = 7;", "/** @type {number} */ var x;")

	c.Test("var x = 4; var x = 7;", "/** @const {UnusableType} */ var x;")

	c.Test("var x = 4; x = 7;", "/** @const {UnusableType} */ var x;")
}

// Java: ConvertToTypedInterfaceTest.testArrowFunctions()
func (c *ConvertToTypedInterfaceTest) TestArrowFunctions() {
	c.TestSame("/** @return {void} */ const f = () => {}")

	c.Test("/** @return {number} */ const f = () => 5", "/** @return {number} */ const f = () => {}")

	c.Test(
		"/** @return {string} */ const f = () => { return 'str' }",
		"/** @return {string} */ const f = () => {}")
}

// Java: ConvertToTypedInterfaceTest.testDontRemoveGoogModuleContents()
func (c *ConvertToTypedInterfaceTest) TestDontRemoveGoogModuleContents() {
	c.TestWarning(
		"goog.module('x.y.z'); var C = goog.require('a.b.C'); exports = new C;",
		ConvertToTypedInterface_CONSTANT_WITHOUT_EXPLICIT_TYPE)

	c.TestSame("goog.module('x.y.z.Foo'); exports = class {};")

	c.TestSame("goog.module('x.y.z'); exports.Foo = class {};")

	c.TestSame("goog.module('x.y.z.Foo'); class Foo {} exports = Foo;")

	c.TestSame("goog.module('x.y.z'); class Foo {} exports.Foo = Foo;")

	c.TestSame("goog.module('x.y.z'); class Foo {} exports = {Foo};")

	c.TestSame(
		`goog.module('x.y.z');
const C = goog.require('a.b.C');
class Foo extends C {}
exports = Foo;
`)
}

// Java: ConvertToTypedInterfaceTest.testDontPreserveUnknownTypeDeclarations()
func (c *ConvertToTypedInterfaceTest) TestDontPreserveUnknownTypeDeclarations() {
	c.TestSame("goog.forwardDeclare('MyType'); /** @type {MyType} */ var x;")

	c.Test(
		"goog.addDependency('zzz.js', ['MyType'], []); /** @type {MyType} */ var x;",
		"/** @type {MyType} */ var x;")

	// This is OK, because short-import goog.forwardDeclares don't declare a type.
	c.TestSame("goog.module('x.y.z'); var C = goog.forwardDeclare('a.b.C'); /** @type {C} */ var c;")
}

// Java: ConvertToTypedInterfaceTest.testAliasOfRequirePreserved()
func (c *ConvertToTypedInterfaceTest) TestAliasOfRequirePreserved() {
	c.TestSame(
		`goog.provide('a.b.c');

goog.require('ns.Foo');

/** @const */
a.b.c.FooAlias = ns.Foo;
`)

	c.TestSame(
		`goog.provide('a.b.c');

goog.require('ns.Foo');

/** @constructor */
a.b.c.FooAlias = ns.Foo;
`)

	c.TestSame(
		`goog.module('mymod');

const {Foo} = goog.require('ns.Foo');

/** @const */
var FooAlias = Foo;

/** @param {!FooAlias} f */
exports = function (f) {};
`)

	c.TestSame(
		`goog.module('mymod');

var Foo = goog.require('ns.Foo');

/** @constructor */
var FooAlias = Foo;

/** @param {!FooAlias} f */
exports = function (f) {};
`)
}

// Java: ConvertToTypedInterfaceTest.testAliasOfNonRequiredName()
func (c *ConvertToTypedInterfaceTest) TestAliasOfNonRequiredName() {
	c.TestSame(
		`goog.provide('a.b.c');

/** @const */
a.b.c.FooAlias = ns.Foo;
`)

	c.TestWarning(
		`goog.provide('a.b.c');

/** @constructor */
a.b.c.Bar = function() {
  /** @const */
  this.FooAlias = ns.Foo;
};
`,
		ConvertToTypedInterface_CONSTANT_WITHOUT_EXPLICIT_TYPE)

	c.TestWarning(
		`goog.module('a.b.c');

class FooAlias {
  constructor() {
    /** @const */
    this.FooAlias = window.Foo;
  }
};
`,
		ConvertToTypedInterface_CONSTANT_WITHOUT_EXPLICIT_TYPE)
}

// Java: ConvertToTypedInterfaceTest.testDuplicateClassMethods()
func (c *ConvertToTypedInterfaceTest) TestDuplicateClassMethods() {
	c.Test(
		`/** @constructor */ var Foo = function() {};
Foo.prototype.method = function() {};

Foo = class {
  method() {}
};
`,
		`/** @constructor */ var Foo = function() {};
Foo.prototype.method = function() {};
`)
}

// Java: ConvertToTypedInterfaceTest.testGoogScopeLeftoversAreRemoved()
func (c *ConvertToTypedInterfaceTest) TestGoogScopeLeftoversAreRemoved() {
	c.Test(
		`goog.provide('a.b.c.d.e.f.g');

/** @const */ var $jscomp = $jscomp || {};
/** @const */ $jscomp.scope = {};

$jscomp.scope.strayVariable = function() {};

a.b.c.d.e.f.g.Foo = class {};
`,
		`goog.provide('a.b.c.d.e.f.g');

a.b.c.d.e.f.g.Foo = class {};
`)

	c.Test(
		`goog.provide('a.b.c.d.e.f.g');

/** @const */ var $jscomp = $jscomp || {};
/** @const */ $jscomp.scope = {};

/** @constructor */
$jscomp.scope.strayCtor = function() { this.x = 5; };

a.b.c.d.e.f.g.Foo = class {};
`,
		`goog.provide('a.b.c.d.e.f.g');

a.b.c.d.e.f.g.Foo = class {};
`,
		c.Warning(ConvertToTypedInterface_GOOG_SCOPE_HIDDEN_TYPE))

	c.Test(
		`goog.provide('a.b.c.d.e.f.g');

/** @const */ var $jscomp = $jscomp || {};
/** @const */ $jscomp.scope = {};

$jscomp.scope.strayClass = class {
  constructor() { this.x = 5 };
  method() {};
};

a.b.c.d.e.f.g.Foo = class {};
`,
		`goog.provide('a.b.c.d.e.f.g');

a.b.c.d.e.f.g.Foo = class {};
`,
		c.Warning(ConvertToTypedInterface_GOOG_SCOPE_HIDDEN_TYPE))

	c.Test(
		`/** @const */ var $jscomp = $jscomp || {};
/** @const */ $jscomp.scope = {};

$jscomp.scope.strayClass = class {
  constructor() {
    this.Foo = class {};
  };
};
`,
		"",
		c.Warning(ConvertToTypedInterface_GOOG_SCOPE_HIDDEN_TYPE))
}

// Java: ConvertToTypedInterfaceTest.testDestructuringDoesntCrash()
func (c *ConvertToTypedInterfaceTest) TestDestructuringDoesntCrash() {
	c.Test(
		`goog.module('a.b.c');

const Enum = goog.require('Enum');
const Foo = goog.require('Foo');

const {A, B} = Enum;

/** @type {Foo} */
exports.foo = use(A, B);
`,
		`goog.module('a.b.c');

const Enum = goog.require('Enum');
const Foo = goog.require('Foo');

const {A, B} = Enum;

/** @type {Foo} */
exports.foo;
`)
}

// Java: ConvertToTypedInterfaceTest.testSameNamedStaticAndNonstaticMethodsDontCrash()
func (c *ConvertToTypedInterfaceTest) TestSameNamedStaticAndNonstaticMethodsDontCrash() {
	c.TestSame(
		`const Foo = class {
  static method() {}
  method() {}
}
`)
}

// Java: ConvertToTypedInterfaceTest.testRedeclarationOfClassMethodDoesntCrash()
func (c *ConvertToTypedInterfaceTest) TestRedeclarationOfClassMethodDoesntCrash() {
	c.Test(
		`class Foo {
  constructor() {
    /** @private */
    this.handleEvent_ = this.handleEvent_.bind(this);
  }
  /** @private @param {Event} e */
  handleEvent_(e) {}
}
`,
		`class Foo {
  constructor() {}
  /** @private @param {Event} e */
  handleEvent_(e) {}
}
`)

	c.Test(
		`class Foo {
  constructor() {
    /** @param {Event} e */
    this.handleEvent_ = function (e) {};
  }
  handleEvent_(e) {}
}
`,
		`class Foo {
  constructor() {
    /** @param {Event} e */
    this.handleEvent_ = function (e) {};
  }
}
`)
}

// Java: ConvertToTypedInterfaceTest.testGoogGlobalTyped()
func (c *ConvertToTypedInterfaceTest) TestGoogGlobalTyped() {
	c.TestSame("/** @const */ var goog = {}; /** @const */ goog.global = this;")
}

// Java: ConvertToTypedInterfaceTest.testAnonymousClassDoesntCrash()
func (c *ConvertToTypedInterfaceTest) TestAnonymousClassDoesntCrash() {
	c.Test("fooFactory(class { x = 5; });", "")

	c.Test("fooFactory(class { x; constructor() { this.x = 5} });", "")

	c.Test(
		"let Foo = fooFactory(class { x; constructor() { this.x = 5;} });", //
		"/** @const {UnusableType} */ var Foo;")

	c.Test(
		"let Foo = fooFactory(class { constructor() {} });",
		"/** @const {UnusableType} */ var Foo;")

	c.Test(
		`let Foo = fooFactory(class {
  constructor() {
    /** @type {number} */
    this.n = 5;
  }
});
`,
		"/** @const {UnusableType} */ var Foo;")

	c.Test(
		`/** @type {function(new:Int)} */
let Foo = fooFactory(class {
  constructor() {
    /** @type {number} */
    this.n = 5;
  }
});
`,
		"/** @type {function(new:Int)} */ var Foo;")
}

// Java: ConvertToTypedInterfaceTest.testPolymerCallsDeclareType()
func (c *ConvertToTypedInterfaceTest) TestPolymerCallsDeclareType() {
	c.TestSame("Polymer({is: 'exampleApp'});")
	c.TestSame("const ExampleApp = Polymer({is: 'exampleApp'});")
}

// Java: ConvertToTypedInterfaceTest.testWellKnownSymbolComputedPropsPreserved()
func (c *ConvertToTypedInterfaceTest) TestWellKnownSymbolComputedPropsPreserved() {
	c.Test(
		`class Foo {
  /** @return {?} */
  [Symbol.iterator]() { return []; }
  static [Symbol.dispose]() { return; }
}
`,
		`class Foo {
  /** @return {?} */
  [Symbol.iterator]() {}
  static [Symbol.dispose]() {}
}
`)
}

// Java: ConvertToTypedInterfaceTest.testSymbolProperties_nonClass()
func (c *ConvertToTypedInterfaceTest) TestSymbolProperties_nonClass() {
	c.TestSame(
		"/** @constructor */ function Foo() {} Foo.prototype[Symbol.iterator] = function(x,y,z)" +
			" {}")

	c.TestSame("/** @constructor */ function Foo() {} /** @type {str} */ Foo[Symbol.iterator];")

	c.Test("weirdDynamicThing()[Symbol.iterator] = function() {};", "")
}

// Java: ConvertToTypedInterfaceTest.testComputedPropertyInferenceDoesntCrash()
func (c *ConvertToTypedInterfaceTest) TestComputedPropertyInferenceDoesntCrash() {
	c.Test("const SomeMap = { [foo()]: 5 };", "const SomeMap = {};")

	c.Test(
		"const SomeBagOfMethods = { /** @return {number} */ method() { return 5; } };",
		"const SomeBagOfMethods = { /** @return {number} */ method() {} };")

	c.Test(
		"const SomeBagOfMethods = { /** @return {number} */ get x() { return 5; } };",
		"const SomeBagOfMethods = { /** @return {number} */ get x() {} };")

	c.Test(
		"const RandomStuff = { [foo()]: 4, method() {}, 9.4: 'bar', set y(x) {} };",
		"const RandomStuff = { method() {}, /** @const {UnusableType} */ '9.4': 0, set y(x) {} };")
}

// Java: ConvertToTypedInterfaceTest.testDescAnnotationCountsAsTyped()
func (c *ConvertToTypedInterfaceTest) TestDescAnnotationCountsAsTyped() {
	c.Test(
		`/** @desc Some description */
var MSG_DESCRIPTION = goog.getMsg('Text');
`,
		`/** @const {string} @desc Some description */
var MSG_DESCRIPTION;
`)

	c.Test(
		`goog.module('a.b.c');

/** @desc Some description */
exports.MSG_DESCRIPTION = goog.getMsg('Text');
`,
		`goog.module('a.b.c');

/** @const {string} @desc Some description */
exports.MSG_DESCRIPTION;
`)
}

// Java: ConvertToTypedInterfaceTest.testCommonJsModules()
func (c *ConvertToTypedInterfaceTest) TestCommonJsModules() {
	c.TestSame("const foo = require('./foo.js');")

	c.TestSame(
		`const {Baz} = require('./baz.js');

exports.Foo = class {
  /** @return {!Baz} */ getBaz() {}
};
`)

	c.Test(
		`module.exports = class Foo {
  /** @return {number} */ get42() { return 40 + 2; }
};
`,
		`module.exports = class Foo {
  /** @return {number} */ get42() {}
};
`)

	c.TestSame(
		`const {Bar, Baz} = require('./a/b/c.js');

class Foo extends Bar {
  /** @return {!Baz} */ getBaz() {}
}

exports = {Foo};
module.exports = {Foo};
`)
}

// Java: ConvertToTypedInterfaceTest.testEmptyFile()
func (c *ConvertToTypedInterfaceTest) TestEmptyFile() {
	c.Test(c.Srcs("const x = 42;", ""), c.Expected("/** @const {number} */ var x;", ""))
}

// Java: ConvertToTypedInterfaceTest.testEs5RecordWithCtorThisDefinitions()
func (c *ConvertToTypedInterfaceTest) TestEs5RecordWithCtorThisDefinitions() {
	c.Test(
		"/** @record */ function Foo() {  /** @type {number} */ this.x = 42;}",
		"/** @record */ function Foo() {}  /** @type {number} */ Foo.prototype.x;")
}

// Java: ConvertToTypedInterfaceTest.testEs5InterfaceWithCtorThisDefinitions()
func (c *ConvertToTypedInterfaceTest) TestEs5InterfaceWithCtorThisDefinitions() {
	c.Test(
		"/** @interface */ function Foo() {  /** @type {number} */ this.x = 42;}",
		"/** @interface */ function Foo() {}  /** @type {number} */ Foo.prototype.x;")
}

// Java: ConvertToTypedInterfaceTest.testPropertyAssignment()
func (c *ConvertToTypedInterfaceTest) TestPropertyAssignment() {
	// Test for b/123413988 which used to crash the ijs generator
	// These tests also exposes a bug tracked in b/124946590
	// TODO(b/124946590): Modify these tests when fixing the bug.
	c.Test(
		`goog.module('fooo');
var Foo;
Foo = class {
    constructor() {
        this.boundOnMouseMove = null;
    }
};
`,
		`goog.module('fooo');
/** @const @type {UnusableType} */ var Foo;
`)

	c.Test(
		`goog.module('fooo');
var Foo;
let Bar = Foo = class {
    constructor() {
        /** @type {null} */
        this.boundOnMouseMove = null;
    }
};
`,
		`goog.module('fooo');
/** @const @type {UnusableType} */ var Foo;
/** @const @type {UnusableType} */ var Bar
`)
}

// Java: ConvertToTypedInterfaceTest.testPolymerBehavior()
func (c *ConvertToTypedInterfaceTest) TestPolymerBehavior() {
	c.Test(
		`/** @polymerBehavior */
export const MyBehavior = {
  properties: {
    foo: String,
    /** @type {string|number} */
    bar: Number,
    /** @type {string|number} */
    baz: {
      type: String,
      value: function() {
        return "foo";
      },
      reflectToAttribute: true,
      observer: "bazChanged",
      readOnly: true
    }
  },
  observers: [
    "abc(foo)",
  ],
  /** @return {boolean} */
  abc() { return true; },
  /** @return {boolean} */
  xyz: function() { return false; }
};
/** @polymerBehavior */
export const MyBehaviorAlias = MyBehavior;
/** @polymerBehavior */
export const MyBehaviorArray = [MyBehavior];
/** @polymerBehavior */
exports.MyBehaviorAlias = MyBehavior;
/** @polymerBehavior */
exports.MyBehaviorArray = [MyBehavior];
/** @polymerBehavior */
export let InvalidBehavior1 = "foo";
/** @polymerBehavior */
exports.InvalidBehavior2 = function() { return "foo" };
export const NotABehavior = {
  properties: {
    foo: String
  }
};
`,
		`/** @polymerBehavior */
export const MyBehavior = {
// The "properties" configuration matters.
  properties: {
    foo: String,
// @type annotations on the properties matter.
    /** @type {string|number} */
    bar: Number,
    /** @type {string|number} */
    baz: {
// If the property definition is an object, only the "type" and "readOnly"
// sub-properties matter.
      type: String,
      readOnly: true
    }
  },
// The "observers" configuration doesn't matter.
  /** @const @type {UnusableType} */
  observers: 0,
// Methods matter, but only their signatures.
  /** @return {boolean} */
  abc() {},
  /** @return {boolean} */
  xyz: function() {}
};
// Behaviors can also be aliased or combined into arrays, and the RHS values matter.
/** @polymerBehavior */
export const MyBehaviorAlias = MyBehavior;
/** @polymerBehavior */
export const MyBehaviorArray = [MyBehavior];
/** @polymerBehavior */
exports.MyBehaviorAlias = MyBehavior;
/** @polymerBehavior */
exports.MyBehaviorArray = [MyBehavior];
// Not valid behavior types, can be simplified.
/**
 * @const
 * @polymerBehavior
 * @type {UnusableType}
 */
export var InvalidBehavior1
/** @polymerBehavior */
exports.InvalidBehavior2 = function() {};
// There's no @polymerBehavior annotation here, so don't preserve "properties".
export const NotABehavior = {
  /** @const @type {UnusableType} */
  properties: 0
};
`)
}

// Java: ConvertToTypedInterfaceTest.testGoogModuleGet()
func (c *ConvertToTypedInterfaceTest) TestGoogModuleGet() {
	c.Test(
		`goog.provide('a.b.c');
goog.provide('a.b.c.d');
goog.require('x.y.z');

/** @const */
a.b.c = {};
/** @const */
a.b.c.d = goog.module.get('x.y.z').d;
`,
		`goog.provide('a.b.c');
goog.provide('a.b.c.d');
goog.require('x.y.z');

/** @const */
a.b.c = {};
/** @const */
a.b.c.d = goog.module.get('x.y.z').d;
`)
}

// Java: ConvertToTypedInterfaceTest.testRemovesClosureUnawareCodeIfPresent()
func (c *ConvertToTypedInterfaceTest) TestRemovesClosureUnawareCodeIfPresent() {
	c.Test(
		`/** @fileoverview @closureUnaware */
goog.module('a.b.c');
/** @closureUnaware */
(function() {})();
exports.foo = 10;
`,
		`/**@fileoverview */
goog.module('a.b.c');

/** @const @type {number} */
exports.foo;
`)
}
