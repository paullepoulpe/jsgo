// Port of com/google/javascript/jscomp/lint/CheckJSDocStyleTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/lint/CheckJSDocStyleTest.java)

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

	"github.com/paullepoulpe/jsgo/internal/parsing"
)

// CheckJSDocStyleTest is the test case for CheckJSDocStyle.
type CheckJSDocStyleTest struct {
	CompilerTestCaseBase

	codingConvention CodingConvention
}

// Java: CheckJSDocStyleTest()
func NewCheckJSDocStyleTest() *CheckJSDocStyleTest {
	c := &CheckJSDocStyleTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, "/** @fileoverview\n * @externs\n */")
	return c
}

func TestCheckJSDocStyleTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCheckJSDocStyleTest() })
}

// Java: CheckJSDocStyleTest.setUp()
func (c *CheckJSDocStyleTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.codingConvention = NewGoogleCodingConvention()
}

// Java: CheckJSDocStyleTest.getProcessor(Compiler)
func (c *CheckJSDocStyleTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewCheckJSDocStyle(compiler)
}

// Java: CheckJSDocStyleTest.getOptions()
func (c *CheckJSDocStyleTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetParseJsDocDocumentation(parsing.Config_JsDocParsing_INCLUDE_DESCRIPTIONS_NO_WHITESPACE)
	options.SetWarningLevel(CheckJSDocStyle_ALL_DIAGNOSTICS, CheckLevel_WARNING)
	return options
}

// Java: CheckJSDocStyleTest.getCodingConvention()
func (c *CheckJSDocStyleTest) GetCodingConvention() CodingConvention {
	return c.codingConvention
}

// Java: CheckJSDocStyleTest.testValidSuppress_onDeclaration()
func (c *CheckJSDocStyleTest) TestValidSuppress_onDeclaration() {
	c.TestSame("/** @const */ var global = this;")
	c.TestSame("/** @const */ goog.global = this;")
}

// Java: CheckJSDocStyleTest.testValidSuppress_withES6Modules01()
func (c *CheckJSDocStyleTest) TestValidSuppress_withES6Modules01() {
	c.TestSame("export /** @suppress {missingRequire} */ var x = new y.Z();")
}

// Java: CheckJSDocStyleTest.testValidSuppress_withES6Modules03()
func (c *CheckJSDocStyleTest) TestValidSuppress_withES6Modules03() {
	c.TestSame("export /** @const @suppress {duplicate} */ var google = {};")
}

// Java: CheckJSDocStyleTest.testExtraneousClassAnnotations()
func (c *CheckJSDocStyleTest) TestExtraneousClassAnnotations() {
	c.TestWarning(
		`/**
 * @constructor
 */
var X = class {};
`,
		CheckJSDocStyle_CLASS_DISALLOWED_JSDOC)
	c.TestWarning(
		`/**
 * @constructor
 */
class X {};
`,
		CheckJSDocStyle_CLASS_DISALLOWED_JSDOC)
	// TODO(tbreisacher): Warn for @extends too. We need to distinguish between cases like this
	// which are totally redundant...
	c.TestSame(
		`/**
 * @extends {Y}
 */
class X extends Y {};
`)
	// ... and ones like this which are not.
	c.TestSame(
		`/**
 * @extends {Y<number>}
 */
class X extends Y {};
`)
	c.TestSame(
		`/**
 * @implements {Z}
 */
class X extends Y {};
`)
	c.TestSame(
		`/**
 * @interface
 * @extends {Y}
 */
class X extends Y {};
`)
	c.TestSame(
		`/**
 * @record
 * @extends {Y}
 */
class X extends Y {};
`)
}

// Java: CheckJSDocStyleTest.testInvalidExtraneousClassAnnotations_withES6Modules()
func (c *CheckJSDocStyleTest) TestInvalidExtraneousClassAnnotations_withES6Modules() {
	c.TestWarning(
		`export
/**
 * @constructor
 */
var X = class {};
`,
		CheckJSDocStyle_CLASS_DISALLOWED_JSDOC)
}

// Java: CheckJSDocStyleTest.testValidExtraneousClassAnnotations_withES6Modules()
func (c *CheckJSDocStyleTest) TestValidExtraneousClassAnnotations_withES6Modules() {
	c.TestSame("export /** @extends {Y} */ class X extends Y {};")
}

// Java: CheckJSDocStyleTest.testNestedArrowFunctions()
func (c *CheckJSDocStyleTest) TestNestedArrowFunctions() {
	c.TestSame(
		`/**
 * @param {Object} a
 * @return {function(Object): boolean}
 */
var haskellStyleEquals = a => b => a == b;
`)
}

// Java: CheckJSDocStyleTest.testNestedArrowFunctions_withES6Modules()
func (c *CheckJSDocStyleTest) TestNestedArrowFunctions_withES6Modules() {
	c.TestSame(
		`export
/**
 * @param {Object} a
 * @return {function(Object): boolean}
 */
var haskellStyleEquals = a => b => a == b;
`)
}

// Java: CheckJSDocStyleTest.testGetterSetterMissingJsDoc()
func (c *CheckJSDocStyleTest) TestGetterSetterMissingJsDoc() {
	c.TestWarning("class Foo { get twentyone() { return 21; } }", CheckJSDocStyle_MISSING_JSDOC)
	c.TestWarning("class Foo { set someString(s) { this.someString_ = s; } }", CheckJSDocStyle_MISSING_JSDOC)
	c.TestSame("class Foo { /** @return {number} */ get twentyone() { return 21; } }")
	c.TestSame("class Foo { /** @param {string} s */ set someString(s) { this.someString_ = s; } }")
}

// Java: CheckJSDocStyleTest.testTypeAnnotationOnGetterSetter()
func (c *CheckJSDocStyleTest) TestTypeAnnotationOnGetterSetter() {
	c.TestWarning(
		"class Foo { /** @type {number} */ get twentyone() { return 21; } }",
		CheckJSDocStyle_INCORRECT_ANNOTATION_ON_GETTER_SETTER)
	c.TestWarning(
		"class Foo { /** @type {string} s */ set someString(s) { this.someString_ = s; } }",
		CheckJSDocStyle_INCORRECT_ANNOTATION_ON_GETTER_SETTER)
	c.TestNoWarning("class Foo { set someString( /** string */ s) { this.someString_ = s; } }")
}

// Java: CheckJSDocStyleTest.testGetterSetter_withES6Modules()
func (c *CheckJSDocStyleTest) TestGetterSetter_withES6Modules() {
	c.TestSame("export class Foo { /** @return {number} */ get twentyone() { return 21; } }")
}

// Java: CheckJSDocStyleTest.testMissingJsDoc()
func (c *CheckJSDocStyleTest) TestMissingJsDoc() {
	c.TestWarning("function f() {}", CheckJSDocStyle_MISSING_JSDOC)
	c.TestWarning("var f = function() {}", CheckJSDocStyle_MISSING_JSDOC)
	c.TestWarning("let f = function() {}", CheckJSDocStyle_MISSING_JSDOC)
	c.TestWarning("const f = function() {}", CheckJSDocStyle_MISSING_JSDOC)
	c.TestWarning("foo.bar = function() {}", CheckJSDocStyle_MISSING_JSDOC)
	c.TestWarning("Foo.prototype.bar = function() {}", CheckJSDocStyle_MISSING_JSDOC)
	c.TestWarning("class Foo { bar() {} }", CheckJSDocStyle_MISSING_JSDOC)
	c.TestWarning("class Foo { constructor(x) {} }", CheckJSDocStyle_MISSING_JSDOC)
	c.TestWarning("var Foo = class { bar() {} };", CheckJSDocStyle_MISSING_JSDOC)
	c.TestWarning("if (COMPILED) { var f = function() {}; }", CheckJSDocStyle_MISSING_JSDOC)
	c.TestWarning("var f = async function() {};", CheckJSDocStyle_MISSING_JSDOC)
	c.TestWarning("async function f() {};", CheckJSDocStyle_MISSING_JSDOC)
	c.TestWarning("Polymer({ method() {} });", CheckJSDocStyle_MISSING_JSDOC)
	c.TestWarning("Polymer({ method: function() {} });", CheckJSDocStyle_MISSING_JSDOC)
	c.TestSame("/** @return {string} */ function f() {}")
	c.TestSame("/** @return {string} */ var f = function() {}")
	c.TestSame("/** @return {string} */ let f = function() {}")
	c.TestSame("/** @return {string} */ const f = function() {}")
	c.TestSame("/** @return {string} */ foo.bar = function() {}")
	c.TestSame("/** @return {string} */ Foo.prototype.bar = function() {}")
	c.TestSame("class Foo { /** @return {string} */ bar() {} }")
	c.TestSame("class Foo { constructor(/** string */ x) {} }")
	c.TestSame("var Foo = class { /** @return {string} */ bar() {} };")
	c.TestSame("/** @param {string} s */ var f = async function(s) {};")
	c.TestSame("/** @param {string} s */ async function f(s) {};")
	c.TestSame("Polymer({ /** @return {null} */ method() {} });")
	c.TestSame("Polymer({ /** @return {null} */ method: function() {} });")
}

// Java: CheckJSDocStyleTest.testMissingJsDoc_withES6Modules01()
func (c *CheckJSDocStyleTest) TestMissingJsDoc_withES6Modules01() {
	c.TestWarning("export function f() {}", CheckJSDocStyle_MISSING_JSDOC)
}

// Java: CheckJSDocStyleTest.testMissingJsDoc_withES6Modules02()
func (c *CheckJSDocStyleTest) TestMissingJsDoc_withES6Modules02() {
	c.TestWarning("export var f = function() {}", CheckJSDocStyle_MISSING_JSDOC)
}

// Java: CheckJSDocStyleTest.testMissingJsDoc_withES6Modules03()
func (c *CheckJSDocStyleTest) TestMissingJsDoc_withES6Modules03() {
	c.TestWarning("export let f = function() {}", CheckJSDocStyle_MISSING_JSDOC)
}

// Java: CheckJSDocStyleTest.testMissingJsDoc_withES6Modules04()
func (c *CheckJSDocStyleTest) TestMissingJsDoc_withES6Modules04() {
	c.TestWarning("export const f = function() {}", CheckJSDocStyle_MISSING_JSDOC)
}

// Java: CheckJSDocStyleTest.testMissingJsDoc_withES6Modules09()
func (c *CheckJSDocStyleTest) TestMissingJsDoc_withES6Modules09() {
	c.TestWarning("export var f = async function() {};", CheckJSDocStyle_MISSING_JSDOC)
}

// Java: CheckJSDocStyleTest.testMissingJsDoc_noWarningIfInlineJsDocIsPresent()
func (c *CheckJSDocStyleTest) TestMissingJsDoc_noWarningIfInlineJsDocIsPresent() {
	c.TestSame("function /** string */ f() {}")
	c.TestSame("function f(/** string */ x) {}")
	c.TestSame("var f = function(/** string */ x) {}")
	c.TestSame("let f = function(/** string */ x) {}")
	c.TestSame("const f = function(/** string */ x) {}")
	c.TestSame("foo.bar = function(/** string */ x) {}")
	c.TestSame("Foo.prototype.bar = function(/** string */ x) {}")
	c.TestSame("class Foo { bar(/** string */ x) {} }")
	c.TestSame("var Foo = class { bar(/** string */ x) {} };")
}

// Java: CheckJSDocStyleTest.testMissingJsDoc_noWarningIfInlineJsDocIsPresent_withES6Modules()
func (c *CheckJSDocStyleTest) TestMissingJsDoc_noWarningIfInlineJsDocIsPresent_withES6Modules() {
	c.TestSame("export function /** string */ f() {}")
}

// Java: CheckJSDocStyleTest.testMissingJsDoc_noWarningIfNotTopLevel()
func (c *CheckJSDocStyleTest) TestMissingJsDoc_noWarningIfNotTopLevel() {
	c.TestSame(c.InIIFE("function f() {}"))
	c.TestSame(c.InIIFE("var f = function() {}"))
	c.TestSame(c.InIIFE("let f = function() {}"))
	c.TestSame(c.InIIFE("const f = function() {}"))
	c.TestSame(c.InIIFE("foo.bar = function() {}"))
	c.TestSame(c.InIIFE("class Foo { bar() {} }"))
	c.TestSame(c.InIIFE("var Foo = class { bar() {} };"))
	c.TestSame("myArray.forEach(function(elem) { alert(elem); });")
	c.TestSame(
		`Polymer({
  is: 'example-elem',
  /** @return {null} */
  someMethod: function() {},
});
`)
	c.TestSame(
		`Polymer({
  is: 'example-elem',
  /** @return {null} */
  someMethod() {},
});
`)
}

// Java: CheckJSDocStyleTest.testMissingJsDoc_noWarningIfNotTopLevelAndNoParams()
func (c *CheckJSDocStyleTest) TestMissingJsDoc_noWarningIfNotTopLevelAndNoParams() {
	c.TestSame(
		`describe('a karma test', function() {
  /** @ngInject */
  var helperFunction = function($compile, $rootScope) {};
})
`)
}

// Java: CheckJSDocStyleTest.testMissingJsDoc_noWarning_wizConstructorAndDeps()
func (c *CheckJSDocStyleTest) TestMissingJsDoc_noWarning_wizConstructorAndDeps() {
	// Exempt Wiz controller constructor and deps() method because Wiz automatically adds JSDoc
	// NOTE(lharker@): right now this does not warn because of b/124061048: the behavior is correct
	// but for the wrong reason.
	c.TestSame(
		`goog.module('a.b.MyController');
class MyController extends SomeParentController {
  static deps() { return {model: 0}; }
  constructor({model}) {}
}
registerController(MY_CONTROLLER, MyController);
`)
}

// Java: CheckJSDocStyleTest.testMissingJsDoc_noWarningOnTestFunctions()
func (c *CheckJSDocStyleTest) TestMissingJsDoc_noWarningOnTestFunctions() {
	c.TestSame("function testSomeFunctionality() {}")
	c.TestSame("var testSomeFunctionality = function() {};")
	c.TestSame("let testSomeFunctionality = function() {};")
	c.TestSame("window.testSomeFunctionality = function() {};")
	c.TestSame("const testSomeFunctionality = function() {};")
	c.TestSame("function setUp() {}")
	c.TestSame("function tearDown() {}")
	c.TestSame("var setUp = function() {};")
	c.TestSame("var tearDown = function() {};")
}

// Java: CheckJSDocStyleTest.testMissingJsDoc_noWarningOnTestMethods()
func (c *CheckJSDocStyleTest) TestMissingJsDoc_noWarningOnTestMethods() {
	c.TestSame("class MyClass { testSomeFunctionality() {} }")
	c.TestSame("goog.module('mod'); class MyClass { testSomeFunctionality() {} }")
	c.TestSame("a.b.c = class { testSomeFunctionality() {} }")
	c.TestSame("class MyClass { setUp() {} }")
	c.TestSame("class MyClass { tearDown() {} }")
}

// Java: CheckJSDocStyleTest.testMissingJsDoc_noWarningOnTestFunctions_withES6Modules()
func (c *CheckJSDocStyleTest) TestMissingJsDoc_noWarningOnTestFunctions_withES6Modules() {
	c.TestSame("export function testSomeFunctionality() {}")
}

// Java: CheckJSDocStyleTest.testMissingJsDoc_noWarningOnEmptyConstructor()
func (c *CheckJSDocStyleTest) TestMissingJsDoc_noWarningOnEmptyConstructor() {
	c.TestSame("class Foo { constructor() {} }")
}

// Java: CheckJSDocStyleTest.testMissingJsDoc_noWarningOnEmptyConstructor_withES6Modules()
func (c *CheckJSDocStyleTest) TestMissingJsDoc_noWarningOnEmptyConstructor_withES6Modules() {
	c.TestSame("export class Foo { constructor() {} }")
}

// Java: CheckJSDocStyleTest.testMissingJsDoc_googModule()
func (c *CheckJSDocStyleTest) TestMissingJsDoc_googModule() {
	c.TestWarning("goog.module('a.b.c'); function f() {}", CheckJSDocStyle_MISSING_JSDOC)
	c.TestWarning("goog.module('a.b.c'); var f = function() {};", CheckJSDocStyle_MISSING_JSDOC)
	// TODO(b/124061048): these should also warn for missing JSDoc
	c.TestSame("goog.module('a.b.c'); class Foo { constructor(x) {} }")
	c.TestSame("goog.module('a.b.c'); class Foo { someMethod() {} }")
}

// Java: CheckJSDocStyleTest.testMissingJsDoc_ES6Module01()
func (c *CheckJSDocStyleTest) TestMissingJsDoc_ES6Module01() {
	c.TestWarning("export default abc; function f() {}", CheckJSDocStyle_MISSING_JSDOC)
}

// Java: CheckJSDocStyleTest.testMissingJsDoc_ES6Module02()
func (c *CheckJSDocStyleTest) TestMissingJsDoc_ES6Module02() {
	c.TestWarning("export default abc; var f = function() {};", CheckJSDocStyle_MISSING_JSDOC)
}

// Java: CheckJSDocStyleTest.testMissingJsDoc_ES6Module03()
func (c *CheckJSDocStyleTest) TestMissingJsDoc_ES6Module03() {
	c.TestWarning("export function f() {};", CheckJSDocStyle_MISSING_JSDOC)
}

// Java: CheckJSDocStyleTest.testMissingJsDoc_ES6Module04()
func (c *CheckJSDocStyleTest) TestMissingJsDoc_ES6Module04() {
	c.TestWarning("export default function () {}", CheckJSDocStyle_MISSING_JSDOC)
}

// Java: CheckJSDocStyleTest.testMissingJsDoc_ES6Module05()
func (c *CheckJSDocStyleTest) TestMissingJsDoc_ES6Module05() {
	c.TestWarning("export default (foo) => { alert(foo); }", CheckJSDocStyle_MISSING_JSDOC)
}

// Java: CheckJSDocStyleTest.testMissingJsDoc_googModule_noWarning()
func (c *CheckJSDocStyleTest) TestMissingJsDoc_googModule_noWarning() {
	c.TestSame("goog.module('a.b.c'); /** @type {function()} */ function f() {}")
	c.TestSame("goog.module('a.b.c'); /** @type {function()} */ var f = function() {};")
	// No param constructors do not require JSDoc
	c.TestSame("goog.module('a.b.c'); class Foo { constructor() {} }")
}

// Java: CheckJSDocStyleTest.testMissingJsDoc_ES6Module_noWarning01()
func (c *CheckJSDocStyleTest) TestMissingJsDoc_ES6Module_noWarning01() {
	c.TestSame("export default abc; /** @type {function()} */ function f() {}")
}

// Java: CheckJSDocStyleTest.testMissingJsDoc_ES6Module_noWarning02()
func (c *CheckJSDocStyleTest) TestMissingJsDoc_ES6Module_noWarning02() {
	c.TestSame("export default abc; /** @type {function()} */ var f = function() {};")
}

// Java: CheckJSDocStyleTest.inIIFE(String)
// (static in Java)
func (c *CheckJSDocStyleTest) InIIFE(js string) string {
	return "(function() {\n" + js + "\n})()"
}

// Java: CheckJSDocStyleTest.testMissingParam_noWarning()
func (c *CheckJSDocStyleTest) TestMissingParam_noWarning() {
	c.TestSame(
		`/**
 * @param {string} x
 * @param {string} y
 */
function f(x, y) {}
`)
	c.TestSame(
		`/**
 * @param {string=} x
 */
function f(x = 1) {}
`)
	c.TestSame(
		`/**
 * @param {number=} x
 * @param {number=} y
 * @param {number=} z
 */
function f(x = 1, y = 2, z = 3) {}
`)
	c.TestSame(
		`/**
 * @param {...string} args
 */
function f(...args) {}
`)
	c.TestSame(
		`(function() {
  myArray.forEach(function(elem) { alert(elem); });
})();
`)
	c.TestSame(
		`(function() {
  myArray.forEach(elem => alert(elem));
})();
`)
	c.TestSame("/** @type {function(number)} */ function f(x) {}")
	c.TestSame("function f(/** string */ inlineArg) {}")
	c.TestSame("/** @export */ function f(/** string */ inlineArg) {}")
	c.TestSame("class Foo { constructor(/** string */ inlineArg) {} }")
	c.TestSame("class Foo { method(/** string */ inlineArg) {} }")
	c.TestSame("/** @export */ class Foo { constructor(/** string */ inlineArg) {} }")
	c.TestSame("class Foo { /** @export */ method(/** string */ inlineArg) {} }")
}

// Java: CheckJSDocStyleTest.testMissingParam_noWarning_withES6Modules()
func (c *CheckJSDocStyleTest) TestMissingParam_noWarning_withES6Modules() {
	c.TestSame("export class Foo { /** @export */ method(/** string */ inlineArg) {} }")
}

// Java: CheckJSDocStyleTest.testMissingParam()
func (c *CheckJSDocStyleTest) TestMissingParam() {
	c.TestWarning(
		`/**
 * @param {string} x
// No @param for y.
 */
function f(x, y) {}
`,
		CheckJSDocStyle_WRONG_NUMBER_OF_PARAMS)
	c.TestWarning(
		`/**
 * @param {string} x
 */
function f(x = 1) {}
`,
		CheckJSDocStyle_OPTIONAL_PARAM_NOT_MARKED_OPTIONAL)
	c.TestWarning(
		`/**
 * @param {string} x
// No @param for y.
 */
function f(x, y = 1) {}
`,
		CheckJSDocStyle_WRONG_NUMBER_OF_PARAMS)
	c.TestWarning("function f(/** string */ x, y) {}", CheckJSDocStyle_MISSING_PARAMETER_JSDOC)
	c.TestWarning("function f(x, /** string */ y) {}", CheckJSDocStyle_MISSING_PARAMETER_JSDOC)
	c.TestWarning("function /** string */ f(x) {}", CheckJSDocStyle_MISSING_PARAMETER_JSDOC)
	c.TestWarning(c.InIIFE("function f(/** string */ x, y) {}"), CheckJSDocStyle_MISSING_PARAMETER_JSDOC)
	c.TestWarning(c.InIIFE("function f(x, /** string */ y) {}"), CheckJSDocStyle_MISSING_PARAMETER_JSDOC)
	c.TestWarning(c.InIIFE("function /** string */ f(x) {}"), CheckJSDocStyle_MISSING_PARAMETER_JSDOC)
}

// Java: CheckJSDocStyleTest.testMissingParam_withES6Modules01()
func (c *CheckJSDocStyleTest) TestMissingParam_withES6Modules01() {
	c.TestWarning(
		`export
/**
 * @param {string} x
// No @param for y.
 */
function f(x, y) {}
`,
		CheckJSDocStyle_WRONG_NUMBER_OF_PARAMS)
}

// Java: CheckJSDocStyleTest.testMissingParam_withES6Modules02()
func (c *CheckJSDocStyleTest) TestMissingParam_withES6Modules02() {
	c.TestWarning(
		"export /** @param {string} x */ function f(x = 1) {}", CheckJSDocStyle_OPTIONAL_PARAM_NOT_MARKED_OPTIONAL)
}

// Java: CheckJSDocStyleTest.testMissingParam_withES6Modules03()
func (c *CheckJSDocStyleTest) TestMissingParam_withES6Modules03() {
	c.TestWarning("export function f(/** string */ x, y) {}", CheckJSDocStyle_MISSING_PARAMETER_JSDOC)
}

// Java: CheckJSDocStyleTest.testMissingParamWithDestructuringPattern()
func (c *CheckJSDocStyleTest) TestMissingParamWithDestructuringPattern() {
	c.TestWarning(
		`/**
 * @param {string} namedParam
 * @return {void}
 */
function f(namedParam, {destructuring:pattern}) {
}
`,
		CheckJSDocStyle_WRONG_NUMBER_OF_PARAMS)
	c.TestWarning(
		`/**
 * @param {string} namedParam
 * @return {void}
 */
function f({destructuring:pattern}, namedParam) {
}
`,
		CheckJSDocStyle_WRONG_NUMBER_OF_PARAMS)
	c.TestWarning(
		`/**
 * @param {string} namedParam
 * @return {void}
 */
function f(namedParam, [pattern]) {
}
`,
		CheckJSDocStyle_WRONG_NUMBER_OF_PARAMS)
	c.TestWarning(
		`/**
 * @param {string} namedParam
 * @return {void}
 */
function f([pattern], namedParam) {
}
`,
		CheckJSDocStyle_WRONG_NUMBER_OF_PARAMS)
	c.TestWarning(
		`/**
 * @param {{
 *   a: (string|undefined),
 *   b: (number|undefined),
 *   c: (boolean|undefined)
 * }} obj
 */
function create({a = 'hello', b = 8, c = false} = {}) {}
`,
		CheckJSDocStyle_OPTIONAL_PARAM_NOT_MARKED_OPTIONAL)
	// Same as above except there's an '=' to indicate that it's optional.
	c.TestSame(
		`/**
 * @param {{
 *   a: (string|undefined),
 *   b: (number|undefined),
 *   c: (boolean|undefined)
 * }=} obj
 */
function create({a = 'hello', b = 8, c = false} = {}) {}
`)
}

// Java: CheckJSDocStyleTest.testMissingParam_defaultValue()
func (c *CheckJSDocStyleTest) TestMissingParam_defaultValue() {
	c.TestWarning(
		`/**
 * @param {string} x
// No @param for y.
 */
function f(x, y = 0) {}
`,
		CheckJSDocStyle_WRONG_NUMBER_OF_PARAMS)
	c.TestWarning(
		`/**
 * @param {string} x
 * @param {number} y
 */
function f(x, y = 0) {}
`,
		CheckJSDocStyle_OPTIONAL_PARAM_NOT_MARKED_OPTIONAL)
	c.TestNoWarning(
		`/**
 * @param {string} x
 * @param {number=} y
 */
function f(x, y = 0) {}
`)
	c.TestWarning("function f(/** string */ x, y = 0) {}", CheckJSDocStyle_MISSING_PARAMETER_JSDOC)
	c.TestWarning(
		"function f(/** string */ x, /** number */ y = 0) {}", CheckJSDocStyle_OPTIONAL_PARAM_NOT_MARKED_OPTIONAL)
	c.TestNoWarning("function f(/** string */ x, /** number= */ y = 0) {}")
}

// Java: CheckJSDocStyleTest.testMissingParam_rest()
func (c *CheckJSDocStyleTest) TestMissingParam_rest() {
	c.TestWarning(
		`/**
 * @param {string} x
// No @param for y.
 */
function f(x, ...y) {}
`,
		CheckJSDocStyle_WRONG_NUMBER_OF_PARAMS)
	c.TestNoWarning(
		`/**
 * @param {string} x
 * @param {...number} y
 */
function f(x, ...y) {}
`)
	c.TestWarning("function f(/** string */ x, ...y) {}", CheckJSDocStyle_MISSING_PARAMETER_JSDOC)
	c.TestNoWarning("function f(/** string */ x, /** ...number */ ...y) {}")
}

// Java: CheckJSDocStyleTest.testInvalidMissingParamWithDestructuringPattern_withES6Modules01()
func (c *CheckJSDocStyleTest) TestInvalidMissingParamWithDestructuringPattern_withES6Modules01() {
	c.TestWarning(
		`export
/**
 * @param {string} namedParam
 * @return {void}
 */
function f(namedParam, {destructuring:pattern}) {
}
`,
		CheckJSDocStyle_WRONG_NUMBER_OF_PARAMS)
}

// Java: CheckJSDocStyleTest.testInvalidMissingParamWithDestructuringPattern_withES6Modules02()
func (c *CheckJSDocStyleTest) TestInvalidMissingParamWithDestructuringPattern_withES6Modules02() {
	c.TestWarning(
		`export
/**
 * @param {{
 *   a: (string|undefined),
 *   b: (number|undefined),
 *   c: (boolean|undefined)
 * }} obj
 */
function create({a = 'hello', b = 8, c = false} = {}) {}
`,
		CheckJSDocStyle_OPTIONAL_PARAM_NOT_MARKED_OPTIONAL)
}

// Java: CheckJSDocStyleTest.testValidMissingParamWithDestructuringPattern_withES6Modules()
func (c *CheckJSDocStyleTest) TestValidMissingParamWithDestructuringPattern_withES6Modules() {
	c.TestSame(
		`export
/**
 * @param {{
 *   a: (string|undefined),
 *   b: (number|undefined),
 *   c: (boolean|undefined)
 * }=} obj
 */
function create({a = 'hello', b = 8, c = false} = {}) {}
`)
}

// Java: CheckJSDocStyleTest.testMissingParamWithDestructuringPatternWithDefault()
func (c *CheckJSDocStyleTest) TestMissingParamWithDestructuringPatternWithDefault() {
	c.TestWarning(
		`/**
 * @param {string} namedParam
 * @return {void}
 */
function f(namedParam, {destructuring:pattern} = defaultValue) {
}
`,
		CheckJSDocStyle_WRONG_NUMBER_OF_PARAMS)
	c.TestWarning(
		`/**
 * @param {string} namedParam
 * @return {void}
 */
function f(namedParam, [pattern] = defaultValue) {
}
`,
		CheckJSDocStyle_WRONG_NUMBER_OF_PARAMS)
}

// Java: CheckJSDocStyleTest.testMissingParamWithDestructuringPatternWithDefault_withES6Modules()
func (c *CheckJSDocStyleTest) TestMissingParamWithDestructuringPatternWithDefault_withES6Modules() {
	c.TestWarning(
		`export
/**
 * @param {string} namedParam
 * @return {void}
 */
function f(namedParam, {destructuring:pattern} = defaultValue) {
}
`,
		CheckJSDocStyle_WRONG_NUMBER_OF_PARAMS)
}

// Java: CheckJSDocStyleTest.testParamWithNoTypeInfo()
func (c *CheckJSDocStyleTest) TestParamWithNoTypeInfo() {
	c.TestSame(
		`/**
 * @param x A param with no type information.
 */
function f(x) { }
`)
}

// Java: CheckJSDocStyleTest.testParamWithNoTypeInfo_optional()
func (c *CheckJSDocStyleTest) TestParamWithNoTypeInfo_optional() {
	c.TestWarning(
		`/**
 * @param x A param with no type information.
 */
function f(x = undefined) { }
`,
		CheckJSDocStyle_OPTIONAL_PARAM_NOT_MARKED_OPTIONAL)
}

// Java: CheckJSDocStyleTest.testParamWithNoTypeInfo_withES6Modules()
func (c *CheckJSDocStyleTest) TestParamWithNoTypeInfo_withES6Modules() {
	c.TestSame(
		`export
/**
 * @param x A param with no type information.
 */
function f(x) { }
`)
}

// Java: CheckJSDocStyleTest.testMissingPrivate_noWarningWithClosureConvention()
func (c *CheckJSDocStyleTest) TestMissingPrivate_noWarningWithClosureConvention() {
	c.codingConvention = NewClosureCodingConvention()
	c.TestSame(
		`/**
 * @return {number}
 * @private
 */
X.prototype.foo = function() { return 0; }
`)
}

// Java: CheckJSDocStyleTest.testMissingPrivate()
func (c *CheckJSDocStyleTest) TestMissingPrivate() {
	c.TestSame(
		`/**
 * @return {number}
 * @private
 */
X.prototype.foo_ = function() { return 0; }
`)
	c.TestSame(
		`/**
 * @type {number}
 * @private
 */
X.prototype.foo_ = 0;
`)
	c.TestSame(
		`/** @type {number} */
X.prototype['@some_special_property'] = 0;
`)
}

// Java: CheckJSDocStyleTest.testNoPrivateWarningsWithSuppressions()
func (c *CheckJSDocStyleTest) TestNoPrivateWarningsWithSuppressions() {
	c.TestNoWarning(
		`goog.module('mod');
class Foo {
  constructor() {
    /** @private {number} */
    this.n_;
    /** @private {number} */
    this.m_;
  }
  setUp() {
    /** @suppress {checkTypes} */
    this.n_ = ' not a number ';
    this.m_ = 1;
  }
  testSomething() {
    alert(this.n_ + this.m_);
  }
}
`)
}

// Java: CheckJSDocStyleTest.testMissingPrivate_dontWarnOnObjectLiteral()
func (c *CheckJSDocStyleTest) TestMissingPrivate_dontWarnOnObjectLiteral() {
	c.TestSame(
		`var obj = {
  /** @return {number} */
  foo_() { return 0; }
}
`)
}

// Java: CheckJSDocStyleTest.testMissingPrivate_dontWarnOnObjectLiteral_withES6Modules()
func (c *CheckJSDocStyleTest) TestMissingPrivate_dontWarnOnObjectLiteral_withES6Modules() {
	c.TestSame("export var obj = { /** @return {number} */ foo_() { return 0; } }")
}

// Java: CheckJSDocStyleTest.testOptionalArgs()
func (c *CheckJSDocStyleTest) TestOptionalArgs() {
	c.TestSame(
		`/**
 * @param {number=} n
 */
function f(n) {}
`)
	c.TestSame(
		`/**
 * @param {number} opt_n
 */
function f(opt_n) {}
`,
		CheckJSDocStyle_OPTIONAL_PARAM_NOT_MARKED_OPTIONAL)
	c.TestSame(
		`/**
 * @param {number=} opt_n
 */
function f(opt_n) {}
`)
}

// Java: CheckJSDocStyleTest.testValidOptionalArgs_withES6Modules()
func (c *CheckJSDocStyleTest) TestValidOptionalArgs_withES6Modules() {
	c.TestSame("export /** @param {number=} n */ function f(n) {}")
}

// Java: CheckJSDocStyleTest.testInvalidOptionalArgs_withES6Modules()
func (c *CheckJSDocStyleTest) TestInvalidOptionalArgs_withES6Modules() {
	c.TestSame(
		"export /** @param {number} opt_n */ function f(opt_n) {}",
		CheckJSDocStyle_OPTIONAL_PARAM_NOT_MARKED_OPTIONAL)
}

// Java: CheckJSDocStyleTest.testParamsOutOfOrder()
func (c *CheckJSDocStyleTest) TestParamsOutOfOrder() {
	c.TestWarning(
		`/**
 * @param {?} second
 * @param {?} first
 */
function f(first, second) {}
`,
		CheckJSDocStyle_INCORRECT_PARAM_NAME)
}

// Java: CheckJSDocStyleTest.testParamsOutOfOrder_withES6Modules()
func (c *CheckJSDocStyleTest) TestParamsOutOfOrder_withES6Modules() {
	c.TestWarning(
		`export
/**
 * @param {?} second
 * @param {?} first
 */
function f(first, second) {}
`,
		CheckJSDocStyle_INCORRECT_PARAM_NAME)
}

// Java: CheckJSDocStyleTest.testMixedStyles()
func (c *CheckJSDocStyleTest) TestMixedStyles() {
	c.TestWarning(
		`/**
 * @param {?} first
 * @param {string} second
 */
function f(first, /** string */ second) {}
`,
		CheckJSDocStyle_MIXED_PARAM_JSDOC_STYLES)
}

// Java: CheckJSDocStyleTest.testMixedStyles_withES6Modules()
func (c *CheckJSDocStyleTest) TestMixedStyles_withES6Modules() {
	c.TestWarning(
		`export
/**
 * @param {?} first
 * @param {string} second
 */
function f(first, /** string */ second) {}
`,
		CheckJSDocStyle_MIXED_PARAM_JSDOC_STYLES)
}

// Java: CheckJSDocStyleTest.testDestructuring()
func (c *CheckJSDocStyleTest) TestDestructuring() {
	c.TestSame(
		`/**
 * @param {{x: number, y: number}} point
 */
function getDistanceFromZero({x, y}) {}
`)
	c.TestSame("function getDistanceFromZero(/** {x: number, y: number} */ {x, y}) {}")
}

// Java: CheckJSDocStyleTest.testDestructuring_withES6Modules()
func (c *CheckJSDocStyleTest) TestDestructuring_withES6Modules() {
	c.TestSame("export function getDistanceFromZero(/** {x: number, y: number} */ {x, y}) {}")
}

// Java: CheckJSDocStyleTest.testMissingReturn_functionStatement_noWarning()
func (c *CheckJSDocStyleTest) TestMissingReturn_functionStatement_noWarning() {
	c.TestSame("/** @param {number} x */ function f(x) {}")
	c.TestSame("/** @constructor */ function f() {}")
	c.TestSame("/** @param {number} x */ function f(x) { function bar() { return x; } }")
	c.TestSame("/** @param {number} x */ function f(x) { return; }")
	c.TestSame("/** @param {number} x\n * @return {number} */ function f(x) { return x; }")
	c.TestSame("/** @param {number} x */ function /** number */ f(x) { return x; }")
}

// Java: CheckJSDocStyleTest.testMissingReturn_functionStatement_noWarning_withES6Modules()
func (c *CheckJSDocStyleTest) TestMissingReturn_functionStatement_noWarning_withES6Modules() {
	c.TestSame("export /** @param {number} x */ function f(x) {}")
}

// Java: CheckJSDocStyleTest.testMissingReturn_assign_noWarning()
func (c *CheckJSDocStyleTest) TestMissingReturn_assign_noWarning() {
	c.TestSame("/** @param {number} x */ f = function(x) {}")
	c.TestSame("/** @constructor */ f = function() {}")
	c.TestSame("/** @param {number} x */ f = function(x) { function bar() { return x; } }")
	c.TestSame("/** @param {number} x */ f = function(x) { return; }")
	c.TestSame("/** @param {number} x\n * @return {number} */ f = function(x) { return x; }")
}

// Java: CheckJSDocStyleTest.testMissingParamOrReturn_warnOnOverrideMethodsAndFields()
func (c *CheckJSDocStyleTest) TestMissingParamOrReturn_warnOnOverrideMethodsAndFields() {
	c.TestWarning(
		`/**
 * @override
 * @param {string} x
 */
 Foo.Bar = function(x) { return x; }
`, // function assigned to a field
		CheckJSDocStyle_MISSING_RETURN_JSDOC)
	c.TestWarning(
		`/**
 * @override
 * @param {string} x
 */
 function f(x) { return x; }
`,
		CheckJSDocStyle_MISSING_RETURN_JSDOC)
	c.TestWarning(
		`/**
 * @override
 * @return {string}
 */
 Foo.Bar = function(x) { return x; }
`,
		CheckJSDocStyle_MISSING_PARAMETER_JSDOC)
	c.TestWarning(
		`/**
 * @override
 * @param {string} x
 */
Foo.bar = function(x, y) {}
`,
		CheckJSDocStyle_WRONG_NUMBER_OF_PARAMS)
	// also test `@inheritDoc` annotations
	c.TestWarning(
		`/**
 * @inheritDoc
 * @return {string} x
 */
function f(x) { return x; }
`,
		CheckJSDocStyle_MISSING_PARAMETER_JSDOC)
	c.TestWarning(
		`/**
 * @inheritDoc
 * @return {string} x
 */
Foo.Bar = function(x) { return x; }
`, // assigned to a field
		CheckJSDocStyle_MISSING_PARAMETER_JSDOC)
	// inline param type
	c.TestNoWarning(
		`/**
 * @override
 * @return {string}
 */
var f = function(/** string */ x) { return x; }
`)
}

// Java: CheckJSDocStyleTest.testMissingReturn_var_noWarning()
func (c *CheckJSDocStyleTest) TestMissingReturn_var_noWarning() {
	c.TestSame("/** @param {number} x */ var f = function(x) {}")
	c.TestSame("/** @constructor */ var f = function() {}")
	c.TestSame("/** @param {number} x */ var f = function(x) { function bar() { return x; } }")
	c.TestSame("/** @param {number} x */ var f = function(x) { return; }")
	c.TestSame("/** @param {number} x\n * @return {number} */ var f = function(x) { return x; }")
	c.TestSame("/** @const {function(number): number} */ var f = function(x) { return x; }")
}

// Java: CheckJSDocStyleTest.testMissingReturn_constructor_noWarning()
func (c *CheckJSDocStyleTest) TestMissingReturn_constructor_noWarning() {
	c.TestSame("/** @constructor */ var C = function() { return null; }")
}

// Java: CheckJSDocStyleTest.testMissingReturn_class_constructor_noWarning()
func (c *CheckJSDocStyleTest) TestMissingReturn_class_constructor_noWarning() {
	c.TestSame("class C { /** @param {Array} x */ constructor(x) { return x; } }")
}

// Java: CheckJSDocStyleTest.testMissingReturn_var_noWarning_withES6Modules()
func (c *CheckJSDocStyleTest) TestMissingReturn_var_noWarning_withES6Modules() {
	c.TestSame("export /** @param {number} x */ var f = function(x) {}")
}

// Java: CheckJSDocStyleTest.testMissingReturn_functionStatement()
func (c *CheckJSDocStyleTest) TestMissingReturn_functionStatement() {
	c.TestWarning("/** @param {number} x */ function f(x) { return x; }", CheckJSDocStyle_MISSING_RETURN_JSDOC)
	c.TestWarning(
		`/** @param {number} x */
function f(x) {
  /** @param {number} x */
  function bar(x) {
    return x;
  }
}
`,
		CheckJSDocStyle_MISSING_RETURN_JSDOC)
	c.TestWarning(
		"/** @param {number} x */ function f(x) { if (true) { return x; } }", CheckJSDocStyle_MISSING_RETURN_JSDOC)
	c.TestWarning(
		"/** @param {number} x @constructor */ function f(x) { return x; }", CheckJSDocStyle_MISSING_RETURN_JSDOC)
}

// Java: CheckJSDocStyleTest.testMissingReturn_functionStatement_withES6Modules()
func (c *CheckJSDocStyleTest) TestMissingReturn_functionStatement_withES6Modules() {
	c.TestWarning(
		"export /** @param {number} x */ function f(x) { return x; }", CheckJSDocStyle_MISSING_RETURN_JSDOC)
}

// Java: CheckJSDocStyleTest.testMissingReturn_assign()
func (c *CheckJSDocStyleTest) TestMissingReturn_assign() {
	c.TestWarning("/** @param {number} x */ f = function(x) { return x; }", CheckJSDocStyle_MISSING_RETURN_JSDOC)
	c.TestWarning(
		`/** @param {number} x */
function f(x) {
  /** @param {number} x */
  bar = function(x) {
    return x;
  }
}
`,
		CheckJSDocStyle_MISSING_RETURN_JSDOC)
	c.TestWarning(
		"/** @param {number} x */ f = function(x) { if (true) { return x; } }",
		CheckJSDocStyle_MISSING_RETURN_JSDOC)
	c.TestWarning(
		"/** @param {number} x @constructor */ f = function(x) { return x; }",
		CheckJSDocStyle_MISSING_RETURN_JSDOC)
}

// Java: CheckJSDocStyleTest.testMissingReturn_assign_withES6Modules()
func (c *CheckJSDocStyleTest) TestMissingReturn_assign_withES6Modules() {
	c.TestWarning(
		`/** @param {number} x */
export
function f(x) {
  /** @param {number} x */
  bar = function(x) {
    return x;
  }
}
`,
		CheckJSDocStyle_MISSING_RETURN_JSDOC)
}

// Java: CheckJSDocStyleTest.testMissingReturn_var()
func (c *CheckJSDocStyleTest) TestMissingReturn_var() {
	c.TestWarning("/** @param {number} x */ var f = function(x) { return x; }", CheckJSDocStyle_MISSING_RETURN_JSDOC)
	c.TestWarning(
		`/** @param {number} x */
function f(x) {
  /** @param {number} x */
  var bar = function(x) {
    return x;
  }
}
`,
		CheckJSDocStyle_MISSING_RETURN_JSDOC)
	c.TestWarning(
		"/** @param {number} x */ var f = function(x) { if (true) { return x; } }",
		CheckJSDocStyle_MISSING_RETURN_JSDOC)
	c.TestWarning(
		"/** @param {number} x @constructor */ var f = function(x) { return x; }",
		CheckJSDocStyle_MISSING_RETURN_JSDOC)
}

// Java: CheckJSDocStyleTest.testMissingReturn_var_withES6Modules()
func (c *CheckJSDocStyleTest) TestMissingReturn_var_withES6Modules() {
	c.TestWarning(
		"export /** @param {number} x */ var f = function(x) { return x; }", CheckJSDocStyle_MISSING_RETURN_JSDOC)
}

// Java: CheckJSDocStyleTest.testExternsAnnotation()
func (c *CheckJSDocStyleTest) TestExternsAnnotation() {
	c.Test(c.Externs("function Example() {}"), c.Srcs(""), c.Warning(CheckJSDocStyle_EXTERNS_FILES_SHOULD_BE_ANNOTATED))
	c.TestSame(
		c.Externs(
			"/** @fileoverview Some super cool externs.\n * @externs\n */ function Example() {}"),
		c.Srcs(""))
	c.TestSame(
		c.Externs(
			`/** @fileoverview Some super cool externs.
 * @externs
 */
/** @constructor */ function Example() {}
/** @param {number} x */ function example2(x) {}
`),
		c.Srcs(""))
	c.Test(
		c.Srcs(
			"/** @fileoverview Some externs.\n * @externs\n */ /** @const */ var example;",
			"/** @fileoverview Some more.\n * @externs\n */ /** @const */ var example2;"),
		c.Expected(""))
}

// Java: CheckJSDocStyleTest.testInvalidExternsAnnotation_withES6Modules()
func (c *CheckJSDocStyleTest) TestInvalidExternsAnnotation_withES6Modules() {
	c.Test(
		c.Externs("export function Example() {}"),
		c.Srcs(""),
		c.Warning(CheckJSDocStyle_EXTERNS_FILES_SHOULD_BE_ANNOTATED))
}

// Java: CheckJSDocStyleTest.testValidExternsAnnotation_withES6Modules()
func (c *CheckJSDocStyleTest) TestValidExternsAnnotation_withES6Modules() {
	c.TestSame(
		c.Externs(
			`export /** @fileoverview Some super cool externs.
 * @externs
 */
function Example() {}
`),
		c.Srcs(""))
}

// Java: CheckJSDocStyleTest.testValidLicenseCommentWithoutExterns()
func (c *CheckJSDocStyleTest) TestValidLicenseCommentWithoutExterns() {
	c.TestSame(
		c.Srcs(
			`/**
 * @license
 * Copyright 2024 Google LLC
 */

function Example() {}
`))
}

// Java: CheckJSDocStyleTest.testValidLicenseCommentAndExterns_separateBlocks()
func (c *CheckJSDocStyleTest) TestValidLicenseCommentAndExterns_separateBlocks() {
	c.TestSame(
		c.Externs(
			`/**
 * @license
 * Copyright 2024 Google LLC
 */
/**
 * @fileoverview Some super cool externs.
 * @externs
 */

function Example() {}
`),
		c.Srcs(""))
}

// Java: CheckJSDocStyleTest.testValidLicenseCommentAndExterns_sameBlocks()
func (c *CheckJSDocStyleTest) TestValidLicenseCommentAndExterns_sameBlocks() {
	c.TestSame(
		c.Externs(
			`/**
 * @fileoverview Some super cool externs.
 * @externs
 * @license
 * Copyright 2024 Google LLC
 */

function Example() {}
`),
		c.Srcs(""))
}

// Java: CheckJSDocStyleTest.testInvalidLicenseComment_containsExterns()
func (c *CheckJSDocStyleTest) TestInvalidLicenseComment_containsExterns() {
	c.Test(
		c.Srcs(
			`/**
 * @license
 * Copyright 2024 Google LLC
 * @externs - oh no, this tag is treated as part of the @license!
 */

function Example() {}
`),
		c.Warning(CheckJSDocStyle_LICENSE_CONTAINS_AT_EXTERNS))
}

// Java: CheckJSDocStyleTest.testAtSignCodeDetectedWhenPresent()
func (c *CheckJSDocStyleTest) TestAtSignCodeDetectedWhenPresent() {
	c.TestWarning(
		"/** blah blah {@code blah blah} blah blah */ function f() {}",
		CheckJSDocStyle_PREFER_BACKTICKS_TO_AT_SIGN_CODE)
}
