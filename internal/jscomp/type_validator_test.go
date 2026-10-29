// Port of test/com/google/javascript/jscomp/TypeValidatorTest.java @ closure-compiler v20261021

/*
 * Copyright 2009 The Closure Compiler Authors.
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

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Type-checking tests that can use methods from CompilerTestCase
//
// Java: public final class TypeValidatorTest extends CompilerTestCase
type TypeValidatorTest struct {
	CompilerTestCaseBase
}

// Java: TypeValidatorTest()
func NewTypeValidatorTest() *TypeValidatorTest {
	c := &TypeValidatorTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestTypeValidatorTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewTypeValidatorTest() })
}

// Java: TypeValidatorTest.setUp()
func (c *TypeValidatorTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableTypeCheck()
}

// Java: TypeValidatorTest.getProcessor(Compiler)
func (c *TypeValidatorTest) GetProcessor(compiler *Compiler) CompilerPass {
	return CompilerPassFunc(func(externs rhino.Node, n rhino.Node) {
		// Do nothing: we're in it for the type-checking.
	})
}

// Java: TypeValidatorTest.testBasicMismatch()
func (c *TypeValidatorTest) TestBasicMismatch() {
	c.TestWarning("/** @param {number} x */ function f(x) {} f('a');", TypeValidator_TYPE_MISMATCH_WARNING)
	c.assertThatRecordedMismatches().
		ComparingElementsUsing(TypeValidatorTest_HAVE_SAME_TYPES).
		ContainsExactly(c.fromNatives(rhino.JSTypeNative_STRING_TYPE, rhino.JSTypeNative_NUMBER_TYPE))
}

// Java: TypeValidatorTest.testFunctionMismatch()
func (c *TypeValidatorTest) TestFunctionMismatch() {
	c.TestWarning(
		`/** 
 * @param {function(string): number} x 
 * @return {function(boolean): string} 
 */ function f(x) { return x; }
`,
		TypeValidator_TYPE_MISMATCH_WARNING)
	registry := c.GetLastCompiler().GetTypeRegistry()
	string_ := registry.GetNativeType(rhino.JSTypeNative_STRING_TYPE)
	bool_ := registry.GetNativeType(rhino.JSTypeNative_BOOLEAN_TYPE)
	number := registry.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)
	firstFunction := registry.CreateFunctionType(number, string_)
	secondFunction := registry.CreateFunctionType(string_, bool_)
	c.assertThatRecordedMismatches().
		ComparingElementsUsing(TypeValidatorTest_HAVE_SAME_TYPES).
		ContainsExactly(
			TypeMismatch_CreateForTesting(firstFunction, secondFunction),
			c.fromNatives(rhino.JSTypeNative_STRING_TYPE, rhino.JSTypeNative_BOOLEAN_TYPE),
			c.fromNatives(rhino.JSTypeNative_NUMBER_TYPE, rhino.JSTypeNative_STRING_TYPE))
}

// Java: TypeValidatorTest.testFunctionMismatch2()
func (c *TypeValidatorTest) TestFunctionMismatch2() {
	c.TestWarning(
		`/** 
 * @param {function(string): number} x 
 * @return {function(boolean): number} 
 */ function f(x) { return x; }
`,
		TypeValidator_TYPE_MISMATCH_WARNING)
	registry := c.GetLastCompiler().GetTypeRegistry()
	string_ := registry.GetNativeType(rhino.JSTypeNative_STRING_TYPE)
	bool_ := registry.GetNativeType(rhino.JSTypeNative_BOOLEAN_TYPE)
	number := registry.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)
	firstFunction := registry.CreateFunctionType(number, string_)
	secondFunction := registry.CreateFunctionType(number, bool_)
	c.assertThatRecordedMismatches().
		ComparingElementsUsing(TypeValidatorTest_HAVE_SAME_TYPES).
		ContainsExactly(
			TypeMismatch_CreateForTesting(firstFunction, secondFunction),
			c.fromNatives(rhino.JSTypeNative_STRING_TYPE, rhino.JSTypeNative_BOOLEAN_TYPE))
}

// Java: TypeValidatorTest.testFunctionMismatchMediumLengthTypes()
func (c *TypeValidatorTest) TestFunctionMismatchMediumLengthTypes() {
	c.Test(
		c.Externs(""),
		c.Srcs(
			`/**
 * @param {{a: string, b: string, c: string, d: string, e: string}} x
 */
function f(x) {}
var y = {a:'',b:'',c:'',d:'',e:0};
f(y);
`),
		c.Warning(TypeValidator_TYPE_MISMATCH_WARNING).
			WithMessage(
				`actual parameter 1 of f does not match formal parameter
found   : {
  a: string,
  b: string,
  c: string,
  d: string,
  e: (number|string)
}
required: {
  a: string,
  b: string,
  c: string,
  d: string,
  e: string
}
missing : []
mismatch: [e]
`))
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: TypeValidatorTest.missingPropertiesPrintedWithMismatch_defaultParam()
func (c *TypeValidatorTest) TestMissingPropertiesPrintedWithMismatch_defaultParam() {
	c.Test(
		c.Externs(""),
		c.Srcs(
			`/** @unrestricted */
class NumberFormatSymbolsMap {
  constructor() {}
}
/** @type {string} */
NumberFormatSymbolsMap.prototype.DEF_CURRENCY_CODE;

var /** !NumberFormatSymbolsMap */ numberSymbol;
/**
 * @record
 */
const Type = class {
  constructor() {
    /** @type {string} */ this.CURRENCY_PATTERN;
    /** @type {string} */ this.DEF_CURRENCY_CODE;
  }
};
/**
 * @param {!Type=} symbols
 * @constructor
 */
let NumberFormat = function(symbols) {};
new NumberFormat(numberSymbol);
`),
		c.Warning(TypeValidator_TYPE_MISMATCH_WARNING).WithMessageContaining("missing : [CURRENCY_PATTERN]"))
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: TypeValidatorTest.missingPropertiesPrintedWithMismatch_nullableParam()
func (c *TypeValidatorTest) TestMissingPropertiesPrintedWithMismatch_nullableParam() {
	c.Test(
		c.Externs(""),
		c.Srcs(
			`/** @unrestricted */
class NumberFormatSymbolsMap {
  constructor() {}
}
/** @type {string} */
NumberFormatSymbolsMap.prototype.DEF_CURRENCY_CODE;

var /** !NumberFormatSymbolsMap */ numberSymbol;
/**
 * @record
 */
const Type = class {
  constructor() {
    /** @type {string} */ this.CURRENCY_PATTERN;
    /** @type {string} */ this.DEF_CURRENCY_CODE;
  }
};
/**
 * @param {?Type} symbols
 * @constructor
 */
let NumberFormat = function(symbols) {};
new NumberFormat(numberSymbol);
`),
		c.Warning(TypeValidator_TYPE_MISMATCH_WARNING).WithMessageContaining("missing : [CURRENCY_PATTERN]"))
}

// Make sure the 'found' and 'required' strings are not identical when there is a mismatch. See
// https://code.google.com/p/closure-compiler/issues/detail?id=719.
// Java: TypeValidatorTest.testFunctionMismatchLongTypes()
func (c *TypeValidatorTest) TestFunctionMismatchLongTypes() {
	c.Test(
		c.Externs(""),
		c.Srcs(
			`/**
 * @param {{a: string, b: string, c: string, d: string, e: string,
 *          f: string, g: string, h: string, i: string, j: string, k: string}} x
 */
function f(x) {}
var y = {a:'',b:'',c:'',d:'',e:'',f:'',g:'',h:'',i:'',j:'',k:0};
f(y);
`),
		c.Warning(TypeValidator_TYPE_MISMATCH_WARNING).
			WithMessage(
				`actual parameter 1 of f does not match formal parameter
found   : {a: string, b: string, c: string, d: string, e: string, f: string, g: string, h: string, i: string, j: string, k: (number|string)}
required: {a: string, b: string, c: string, d: string, e: string, f: string, g: string, h: string, i: string, j: string, k: string}
missing : []
mismatch: [k]
`))
}

// Same as testFunctionMismatchLongTypes, but with one of the types being a typedef.
// Java: TypeValidatorTest.testFunctionMismatchTypedef()
func (c *TypeValidatorTest) TestFunctionMismatchTypedef() {
	c.Test(
		c.Externs(""),
		c.Srcs(
			`/**
 * @typedef {{a: string, b: string, c: string, d: string, e: string,
 *            f: string, g: string, h: string, i: string, j: string, k: string}} x
 */
var t;
/**
 * @param {t} x
 */
function f(x) {}
var y = {a:'',b:'',c:'',d:'',e:'',f:'',g:'',h:'',i:'',j:'',k:0};
f(y);
`),
		c.Warning(TypeValidator_TYPE_MISMATCH_WARNING).
			WithMessage(
				`actual parameter 1 of f does not match formal parameter
found   : {a: string, b: string, c: string, d: string, e: string, f: string, g: string, h: string, i: string, j: string, k: (number|string)}
required: {a: string, b: string, c: string, d: string, e: string, f: string, g: string, h: string, i: string, j: string, k: string}
missing : []
mismatch: [k]
`))
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: TypeValidatorTest.bug_testMismatchRecursively_throughFields()
func (c *TypeValidatorTest) TestBug_testMismatchRecursively_throughFields() {
	c.TestWarning(
		`class Foo {
  constructor() {
    /** @type {number} */ this.x;
  }
}

function f(/** {x: string} */ a) {
  const /** !Foo */ b = a;
}
`,
		TypeValidator_TYPE_MISMATCH_WARNING)
	// TODO(b/148169932): There should be another mismatch {found: string, required: number}.
	c.assertThatRecordedMismatches().HasSize(1)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: TypeValidatorTest.bug_testMismatchRecursively_throughTemplates()
func (c *TypeValidatorTest) TestBug_testMismatchRecursively_throughTemplates() {
	c.TestWarning(
		`/**
 * @interface
 * @template T
 */
class Foo { }

/** @implements {Foo<string>} */
class Bar { }

function f(/** !Bar */ a) {
  const /** !Foo<number> */ b = a;
}
`,
		TypeValidator_TYPE_MISMATCH_WARNING)
	// TODO(b/148169932): There should be an another mismatch {found: string, required: number}.
	c.assertThatRecordedMismatches().IsEmpty()
}

// Java: TypeValidatorTest.testNullUndefined()
func (c *TypeValidatorTest) TestNullUndefined() {
	c.TestWarning(
		`/** @param {string} x */ function f(x) {}
f(/** @type {string|null|undefined} */ ('a'));
`,
		TypeValidator_TYPE_MISMATCH_WARNING)
	c.assertThatRecordedMismatches().IsEmpty()
}

// Java: TypeValidatorTest.testSubclass()
func (c *TypeValidatorTest) TestSubclass() {
	c.TestWarning(
		`/** @constructor */
function Super() {}
/**
 * @constructor
 * @extends {Super}
 */
function Sub() {}
/** @param {Sub} x */ function f(x) {}
f(/** @type {Super} */ (new Sub));
`,
		TypeValidator_TYPE_MISMATCH_WARNING)
	c.assertThatRecordedMismatches().IsEmpty()
}

// Java: TypeValidatorTest.testUnionsMismatch()
func (c *TypeValidatorTest) TestUnionsMismatch() {
	c.TestWarning(
		`/** @param {number|string} x */
function f(x) {}
f(/** @type {boolean|string} */ ('a'));
`,
		TypeValidator_TYPE_MISMATCH_WARNING)
	c.assertThatRecordedMismatches().IsEmpty()
}

// Java: TypeValidatorTest.testModuloNullUndef1()
func (c *TypeValidatorTest) TestModuloNullUndef1() {
	c.TestSame(
		c.Srcs(
			SourceFile_FromCode(
				"foo.java.js",
				`function f(/** number */ to, /** (number|null) */ from) {
  to = from;
}
`)))
}

// Java: TypeValidatorTest.testModuloNullUndef2()
func (c *TypeValidatorTest) TestModuloNullUndef2() {
	c.TestSame(
		c.Srcs(
			SourceFile_FromCode(
				"foo.java.js",
				`function f(/** number */ to, /** (number|undefined) */ from) {
  to = from;
}
`)))
}

// Java: TypeValidatorTest.testModuloNullUndef3()
func (c *TypeValidatorTest) TestModuloNullUndef3() {
	c.TestSame(
		c.Srcs(
			SourceFile_FromCode(
				"foo.java.js",
				`/** @constructor */
function Foo() {}
function f(/** !Foo */ to, /** ?Foo */ from) {
  to = from;
}
`)))
}

// Java: TypeValidatorTest.testModuloNullUndef4()
func (c *TypeValidatorTest) TestModuloNullUndef4() {
	c.TestSame(
		c.Srcs(
			SourceFile_FromCode(
				"foo.java.js",
				`/** @constructor */
function Foo() {}
/** @constructor @extends {Foo} */
function Bar() {}
function f(/** !Foo */ to, /** ?Bar */ from) {
  to = from;
}
`)))
}

// Java: TypeValidatorTest.testModuloNullUndef5()
func (c *TypeValidatorTest) TestModuloNullUndef5() {
	c.TestSame(
		c.Srcs(
			SourceFile_FromCode(
				"foo.java.js",
				`function f(/** {a: number} */ to, /** {a: (null|number)} */ from) {
  to = from;
}
`)))
}

// Java: TypeValidatorTest.testModuloNullUndef6()
func (c *TypeValidatorTest) TestModuloNullUndef6() {
	c.TestSame(
		c.Srcs(
			SourceFile_FromCode(
				"foo.java.js",
				`function f(/** {a: number} */ to, /** ?{a: (null|number)} */ from) {
  to = from;
}
`)))
}

// Java: TypeValidatorTest.testModuloNullUndef7()
func (c *TypeValidatorTest) TestModuloNullUndef7() {
	c.TestSame(
		c.Srcs(
			SourceFile_FromCode(
				"foo.java.js",
				`/** @constructor */
function Foo() {}
function f(/** function():!Foo */ to, /** function():?Foo */ from) {
  to = from;
}
`)))
}

// Java: TypeValidatorTest.testModuloNullUndef8()
func (c *TypeValidatorTest) TestModuloNullUndef8() {
	c.TestSame(
		c.Srcs(
			SourceFile_FromCode(
				"foo.java.js",
				`/**
 * @constructor
 * @template T
 */
function Foo() {}
function f(/** !Foo<number> */ to, /** !Foo<(number|null)> */ from) {
  to = from;
}
`)))
}

// Java: TypeValidatorTest.testModuloNullUndef9()
func (c *TypeValidatorTest) TestModuloNullUndef9() {
	c.TestSame(
		c.Srcs(
			SourceFile_FromCode(
				"foo.java.js",
				`/** @interface */
function Foo() {}
/** @type {function(?number)} */
Foo.prototype.prop;
/** @constructor @implements {Foo} */
function Bar() {}
/** @type {function(number)} */
Bar.prototype.prop;
`)))
}

// Java: TypeValidatorTest.testModuloNullUndef10()
func (c *TypeValidatorTest) TestModuloNullUndef10() {
	c.TestSame(
		c.Srcs(
			SourceFile_FromCode(
				"foo.java.js",
				`/** @constructor */
function Bar() {}
/** @type {!number} */
Bar.prototype.prop;
function f(/** ?number*/ n) {
  (new Bar).prop = n;
}
`)))
}

// Java: TypeValidatorTest.testModuloNullUndef11()
func (c *TypeValidatorTest) TestModuloNullUndef11() {
	c.TestSame(
		c.Srcs(
			SourceFile_FromCode(
				"foo.java.js",
				`function f(/** number */ n) {}
f(/** @type {?number} */ (null));
`)))
}

// Java: TypeValidatorTest.testModuloNullUndef12()
func (c *TypeValidatorTest) TestModuloNullUndef12() {
	// Only warn for the file not ending in .java.js
	c.TestWarning(
		c.Srcs(
			SourceFile_FromCode(
				"foo.js",
				`function f(/** number */ to, /** (number|null) */ from) {
  to = from;
}
`),
			SourceFile_FromCode(
				"foo.java.js",
				`function g(/** number */ to, /** (number|null) */ from) {
  to = from;
}
`)),
		c.Warning(TypeValidator_TYPE_MISMATCH_WARNING))
}

// Java: TypeValidatorTest.testModuloNullUndef13()
func (c *TypeValidatorTest) TestModuloNullUndef13() {
	c.TestSame(c.Srcs(SourceFile_FromCode("foo.java.js", "var /** @type {{ a:number }} */ x = null;")))
}

// Java: TypeValidatorTest.testInheritanceModuloNullUndef1()
func (c *TypeValidatorTest) TestInheritanceModuloNullUndef1() {
	c.TestSame(
		c.Srcs(
			SourceFile_FromCode(
				"foo.java.js",
				`/** @constructor */
function Foo() {}
/** @return {string} */
Foo.prototype.toString = function() { return ''; };
/** @constructor @extends {Foo} */
function Bar() {}
/**
 * @override
 * @return {?string}
 */
Bar.prototype.toString = function() { return null; };
`)))
}

// Java: TypeValidatorTest.testInheritanceModuloNullUndef2()
func (c *TypeValidatorTest) TestInheritanceModuloNullUndef2() {
	c.TestSame(
		c.Srcs(
			SourceFile_FromCode(
				"foo.java.js",
				`/** @interface */
function Foo() {}
/** @return {string} */
Foo.prototype.toString = function() {};
/** @constructor @implements {Foo} */
function Bar() {}
/**
 * @override
 * @return {?string}
 */
Bar.prototype.toString = function() {};
`)))
}

// Java: TypeValidatorTest.testInheritanceModuloNullUndef3()
func (c *TypeValidatorTest) TestInheritanceModuloNullUndef3() {
	c.TestSame(
		c.Srcs(
			SourceFile_FromCode(
				"foo.java.js",
				`/** @interface */
function High1() {}
/** @type {number} */
High1.prototype.prop;
/** @interface */
function High2() {}
/** @type {?number} */
High2.prototype.prop;
/**
 * @interface
 * @extends {High1}
 * @extends {High2}
 */
function Low() {}
`)))
}

// Java: TypeValidatorTest.testOptionalProperties_dontNeedInitializer_inConstructor()
func (c *TypeValidatorTest) TestOptionalProperties_dontNeedInitializer_inConstructor() {
	c.TestSame(
		`/** @interface */
function Foo() {}
/** @type {boolean|undefined} */
Foo.prototype.prop;

/** @constructor @implements {Foo} */
function Bar() {
// No initializer here.
}
`)
}

// Java: TypeValidatorTest.testDuplicateSuppression()
func (c *TypeValidatorTest) TestDuplicateSuppression() {
	c.TestWarning(
		`/** @const */
var ns2 = {};
/** @type {number} */
ns2.x = 3;
/** @type {number} */
ns2.x = 3;
`,
		TypeValidator_DUP_VAR_DECLARATION)
	c.TestWarning(
		`/** @const */
var ns2 = {};
/** @type {number|string} */ // so that the second assignment is accepted.
ns2.x = 3;
/** @type {string} */
ns2.x = 'a';
`,
		TypeValidator_DUP_VAR_DECLARATION_TYPE_MISMATCH)
	// catch variables in different catch blocks are not duplicate declarations
	c.TestSame(
		`try { throw 1; } catch (/** @type {number} */ err) {}
try { throw 1; } catch (/** @type {number} */ err) {}
`)
	// duplicates suppressed on 1st declaration.
	c.TestSame(
		`/** @const */
var ns1 = {};
/** @type {number} @suppress {duplicate} */
ns1.x = 3;
/** @type {number} */
ns1.x = 3;
`)
	// duplicates suppressed on 2nd declaration.
	c.TestSame(
		`/** @const */
var ns1 = {};
/** @type {number} */
ns1.x = 3;
/** @type {number} @suppress {duplicate} */
ns1.x = 3;
`)
	// duplicates suppressed on file level.
	c.TestSame(
		`/** @fileoverview @suppress {duplicate} */
/** @const */
var ns1 = {};
/** @type {number} */
ns1.x = 3;
/** @type {number} */
ns1.x = 3;
`)
}

// Java: TypeValidatorTest.testDuplicateSuppression_class()
func (c *TypeValidatorTest) TestDuplicateSuppression_class() {
	c.TestWarning(
		`class X { constructor() {} }
function X() {}
`,
		TypeValidator_DUP_VAR_DECLARATION_TYPE_MISMATCH)
	c.TestSame(
		`/** @suppress {duplicate} */
class X { constructor() {} }
function X() {}
`)
}

// Java: TypeValidatorTest.testDuplicateSuppression_typeMismatch()
func (c *TypeValidatorTest) TestDuplicateSuppression_typeMismatch() {
	// duplicate diagnostic category includes type mismatches.
	c.TestSame(
		`/** @const */
var ns1 = {};
/** @type {number} */
ns1.x = 3;
/** @type {string} @suppress {duplicate} */
ns1.x = 3;
`)
	c.TestSame(
		`/** @const */
var ns1 = {};
/** @type {number} @suppress {duplicate} */
ns1.x = 3;
/** @type {string} */
ns1.x = 3;
`)
}

// Java: TypeValidatorTest.testDuplicateSuppression_stubs()
func (c *TypeValidatorTest) TestDuplicateSuppression_stubs() {
	// No duplicate warning because the first declaration is a stub declaration (property access)
	c.TestSame(
		`/** @const */
var ns0 = {};
/** @type {number} */
ns0.x;
/** @type {number} */
ns0.x;
`)
	// Type mismatch on stub.
	c.TestWarning(
		`/** @const */
var ns3 = {};
/** @type {number} */
ns3.x;
/** @type {string} */
ns3.x;
`,
		TypeValidator_DUP_VAR_DECLARATION_TYPE_MISMATCH)
}

// Java: TypeValidatorTest.testDuplicateSuppression_topLevelVariables()
func (c *TypeValidatorTest) TestDuplicateSuppression_topLevelVariables() {
	c.TestWarning(
		`/** @type {number} */
var w;
/** @type {number} */
var w;
`,
		TypeValidator_DUP_VAR_DECLARATION)
	c.TestWarning(
		`/** @type {number} */
var y = 3;
/** @type {string} */
var y = 3;
`,
		TypeValidator_DUP_VAR_DECLARATION_TYPE_MISMATCH)
	// @suppress on file level.
	c.TestSame(
		`/** @fileoverview  @suppress {duplicate} */
/** @type {number} */
var x;
/** @type {number} */
var x;
`)
	// @suppress on variable declaration.
	c.TestSame(
		`/** @type {number} */
var z;
/** @type {number} @suppress {duplicate} */
var z;
`)
}

// Java: TypeValidatorTest.testDuplicateSuppression_topLevelFunctions()
func (c *TypeValidatorTest) TestDuplicateSuppression_topLevelFunctions() {
	c.TestWarning(
		`/** @return {number} */
function f() {}
/** @return {number} */
function f() {}
`,
		TypeValidator_DUP_VAR_DECLARATION)
	c.TestWarning(
		`/** @return {number} */
function f() {}
/** @return {string} */
function f() {}
`,
		TypeValidator_DUP_VAR_DECLARATION_TYPE_MISMATCH)
	c.TestSame(
		`/** @return {number} */
function f() {}
/** @return {string} @suppress {duplicate} */
function f() {}
`)
	c.TestSame(
		`/** @return {number} */
function f() {}
/** @return {string} @suppress {duplicate} */
function f() {}
`)
}

// Java: TypeValidatorTest.fromNatives(JSTypeNative, JSTypeNative)
func (c *TypeValidatorTest) fromNatives(a rhino.JSTypeNative, b rhino.JSTypeNative) TypeMismatch {
	registry := c.GetLastCompiler().GetTypeRegistry()
	return TypeMismatch_CreateForTesting(registry.GetNativeType(a), registry.GetNativeType(b))
}

// Java: TypeValidatorTest.assertThatRecordedMismatches()
func (c *TypeValidatorTest) assertThatRecordedMismatches() *Subject {
	return AssertThat(c.GetLastCompiler().GetTypeMismatches())
}

// Java: TypeValidatorTest.HAVE_SAME_TYPES
var TypeValidatorTest_HAVE_SAME_TYPES = Correspondence_TransformingBoth(
	func(x TypeMismatch) []rhino.JSType { return []rhino.JSType{x.Found(), x.Required()} },
	func(x TypeMismatch) []rhino.JSType { return []rhino.JSType{x.Found(), x.Required()} },
	"has same types as")
