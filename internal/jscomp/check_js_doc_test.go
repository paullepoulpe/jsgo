// Port of com/google/javascript/jscomp/CheckJsDocTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/CheckJsDocTest.java)

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

// CheckJsDocTest tests CheckJSDoc.
type CheckJsDocTest struct {
	CompilerTestCaseBase

	jsdocParsingMode parsing.Config_JsDocParsing
}

// Java: CheckJsDocTest()
func NewCheckJsDocTest() *CheckJsDocTest {
	c := &CheckJsDocTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestCheckJsDocTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCheckJsDocTest() })
}

// Java: CheckJsDocTest.setUp()
func (c *CheckJsDocTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.jsdocParsingMode = parsing.Config_JsDocParsing_INCLUDE_DESCRIPTIONS_WITH_WHITESPACE
}

// Java: CheckJsDocTest.getProcessor(Compiler)
func (c *CheckJsDocTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewCheckJSDoc(compiler)
}

// Java: CheckJsDocTest.getOptions()
func (c *CheckJsDocTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetWarningLevel(DiagnosticGroups_MISPLACED_SUPPRESS, CheckLevel_WARNING)
	options.SetParseJsDocDocumentation(c.jsdocParsingMode)
	options.SetPreserveDetailedSourceInfo(true)
	return options
}

// Java: CheckJsDocTest.testInlineJsDocOnObjectPatternTargetNames()
func (c *CheckJsDocTest) TestInlineJsDocOnObjectPatternTargetNames() {
	c.TestSame("let {/** string */ prop} = {prop: 'hi'};")
	c.TestSame("let {/** string */ prop = 'default'} = {prop: 'hi'};")
	c.TestSame("let {prop: /** string */ x} = {prop: 'hi'};")
	c.TestSame("let {prop: /** string */ x = 'default'} = {prop: 'hi'};")
	// TODO(bradfordcsmith): These two cases should be misplaced annotations.
	c.TestSame("let {/** string */ prop: x} = {prop: 'hi'};")
	c.TestSame("let {/** string */ prop: x = 'default'} = {prop: 'hi'};")
}

// Java: CheckJsDocTest.testInlineJsDocOnObjectPatternTargetNamesInAssign()
func (c *CheckJsDocTest) TestInlineJsDocOnObjectPatternTargetNamesInAssign() {
	// We don't allow any inline type annotations in object pattern assigns.
	// Simple variable names must have their type declared upon the variable declaration, and
	// qualified names must be declared outside of a destructuring pattern.
	c.TestWarning("({/** string */ prop} = {prop: 'hi'});", CheckJSDoc_MISPLACED_ANNOTATION)
	c.TestWarning(
		"({['prop']: /** string */ prop} = {prop: 'hi'});", CheckJSDoc_MISPLACED_ANNOTATION)
	c.TestWarning("({prop: /** string */ ns.prop} = {prop: 'hi'});", CheckJSDoc_MISPLACED_ANNOTATION)
	c.TestWarning(
		"({prop: /** string */ ns.prop = 'default'} = {prop: 'hi'});",
		CheckJSDoc_MISPLACED_ANNOTATION)
	c.TestWarning(
		"({prop: /** string */ f().y['z']} = {prop: 'hi'});", CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testInlineJsDocOnArrayPatternTargetNames()
func (c *CheckJsDocTest) TestInlineJsDocOnArrayPatternTargetNames() {
	c.TestSame("let [/** string */ x] = ['hi'];")
	c.TestSame("let [/** string */ x = 'lo'] = ['hi'];")
	c.TestSame("try {} catch ([/** string */ x]) {}")
}

// Java: CheckJsDocTest.testInlineJsDocOnArrayPatternTargetInAssign()
func (c *CheckJsDocTest) TestInlineJsDocOnArrayPatternTargetInAssign() {
	// We don't allow any inline type annotations in array pattern assigns.
	// Simple variable names must have their type declared upon the variable declaration, and
	// qualified names must be declared outside of a destructuring pattern.
	c.TestWarning("[/** string */ x] = ['hi'];", CheckJSDoc_MISPLACED_ANNOTATION)
	c.TestWarning("[/** string */ x = 'foo'] = ['hi'];", CheckJSDoc_MISPLACED_ANNOTATION)
	c.TestWarning("[/** string */ x.y] = ['hi'];", CheckJSDoc_MISPLACED_ANNOTATION)
	c.TestWarning("[/** string */ x.y = 'lo'] = ['hi'];", CheckJSDoc_MISPLACED_ANNOTATION)
	c.TestWarning("[/** string */ f().y['z']] = [];", CheckJSDoc_MISPLACED_ANNOTATION)
}

// just here to make sure import.meta doesn't break anything
// Java: CheckJsDocTest.testImportMeta()
func (c *CheckJsDocTest) TestImportMeta() {
	c.TestSame("var /** number */ x = foo; import.meta.foo = x")
}

// Java: CheckJsDocTest.testInlineJsDocOnDeclaration()
func (c *CheckJsDocTest) TestInlineJsDocOnDeclaration() {
	c.TestSame("var /** number */ x;")
	c.TestSame("var /** number */ x = 3;")
	// See testInlineJsDocOnArrayPatternTargetNames() for cases where the inline JSDoc is attached
	// to the names within the pattern.
	c.TestWarning("var /** !Array<number> */ [x, y] = someArr;", CheckJSDoc_MISPLACED_ANNOTATION)
	c.TestWarning("var /** {x: number} */ {x} = someObj;", CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testFieldMisplacedAnnotation()
func (c *CheckJsDocTest) TestFieldMisplacedAnnotation() {
	c.TestWarning(
		`class Foo {
  /** @nocollapse */
  x = 5;
}
`,
		CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testStaticFieldNoCollapse()
func (c *CheckJsDocTest) TestStaticFieldNoCollapse() {
	c.TestSame(
		`class Bar {
  /** @nocollapse */
  static y = 1
}
`)
}

// Java: CheckJsDocTest.testThisPropertyMisplacedAnnotation()
func (c *CheckJsDocTest) TestThisPropertyMisplacedAnnotation() {
	c.TestWarning(
		`class Foo {
  constructor() {
    /** @nocollapse */
    this.x = 4;
  }
}
`,
		CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testThisInFunctionMisplacedAnnotation()
func (c *CheckJsDocTest) TestThisInFunctionMisplacedAnnotation() {
	c.TestWarning(
		`/** @constructor */
function Bar() {
 /** @nocollapse */
 this.x = 4;
}
`,
		CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testInvalidJsDocOnDestructuringDeclaration()
func (c *CheckJsDocTest) TestInvalidJsDocOnDestructuringDeclaration() {
	// Type annotations are not allowed on any declaration containing >=1 destructuring pattern.
	c.TestWarning("/** @type {number} */ const [a] = arr;", CheckJSDoc_MISPLACED_ANNOTATION)
	c.TestWarning("/** @type {number} */ const {a} = obj;", CheckJSDoc_MISPLACED_ANNOTATION)
	c.TestWarning("/** @type {number} */ const a = 1, [b] = arr;", CheckJSDoc_MISPLACED_ANNOTATION)
	c.TestWarning("/** @type {number} */ const a = 1, {b} = obj;", CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testInlineJsDoc_ES6()
func (c *CheckJsDocTest) TestInlineJsDoc_ES6() {
	c.TestSame("function f(/** string */ x) {}")
	c.TestSame("function f(/** number= */ x=3) {}")
	c.TestSame("function f(/** !Object */ {x}) {}")
	c.TestSame("function f(/** !Array */ [x]) {}")
	c.TestSame("function f([/** number */ x]) {}")
}

// Java: CheckJsDocTest.testValidInlineJsDoc_ES6_withES6Modules()
func (c *CheckJsDocTest) TestValidInlineJsDoc_ES6_withES6Modules() {
	c.TestSame("export function f(/** string */ x) {};")
}

// Java: CheckJsDocTest.testInlineJsDoc_ES6_withES6Modules()
func (c *CheckJsDocTest) TestInlineJsDoc_ES6_withES6Modules() {
	c.TestSame("export function f([/** number */ x]) {};")
}

// TODO(tbreisacher): These should be a CheckJSDoc_MISPLACED_ANNOTATION warning instead of silently failing.
// Java: CheckJsDocTest.testInlineJsDocInsideObjectParams()
func (c *CheckJsDocTest) TestInlineJsDocInsideObjectParams() {
	c.TestSame("function f({ prop: {/** string */ x} }) {}")
	c.TestSame("function f({ prop: {x: /** string */ y} }) {}")
	c.TestSame("function f({ /** number */ x }) {}")
	c.TestSame("function f({ prop: /** number */ x }) {}")
}

// Java: CheckJsDocTest.testInlineJsDocInsideObjectParams_withES6Modules()
func (c *CheckJsDocTest) TestInlineJsDocInsideObjectParams_withES6Modules() {
	c.TestSame("export function f({ prop: {/** string */ x} }) {};")
}

// Java: CheckJsDocTest.testJsDocOnCatchVariables()
func (c *CheckJsDocTest) TestJsDocOnCatchVariables() {
	c.TestSame("try {} catch (/** @type {number} */ n) {}")
	// See testInlineJsDocOnArrayPatternTargetNames for cases where the inline JsDoc is attached to
	// the names within the pattern.
	c.TestWarning(
		"try {} catch (/** @type {!Array<number>} */ [x, y]) {}", CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testInvalidJSDocOnName()
func (c *CheckJsDocTest) TestInvalidJSDocOnName() {
	c.TestWarning("var f = () => /** @type {number} */ n;", CheckJSDoc_MISPLACED_ANNOTATION)
	c.TestWarning("var f = /** @type {number} */ n;", CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testInvalidClassJsdoc()
func (c *CheckJsDocTest) TestInvalidClassJsdoc() {
	c.TestSame("class Foo { /** @param {number} x */ constructor(x) {}}")
	c.TestWarning("class Foo { /** @constructor */ constructor() {}}", CheckJSDoc_DISALLOWED_MEMBER_JSDOC)
	c.TestWarning("class Foo { /** @interface */ constructor() {}}", CheckJSDoc_DISALLOWED_MEMBER_JSDOC)
	c.TestWarning("class Foo { /** @extends {Foo} */ constructor() {}}", CheckJSDoc_DISALLOWED_MEMBER_JSDOC)
	c.TestWarning("class Foo { /** @implements {Foo} */ constructor() {}}", CheckJSDoc_DISALLOWED_MEMBER_JSDOC)
}

// Java: CheckJsDocTest.testValidClassJsdoc_withES6Modules()
func (c *CheckJsDocTest) TestValidClassJsdoc_withES6Modules() {
	c.TestSame("export class Foo { /** @param {number} x */ constructor(x) {}; };")
}

// Java: CheckJsDocTest.testInvalidClassJsdoc_withES6Modules()
func (c *CheckJsDocTest) TestInvalidClassJsdoc_withES6Modules() {
	c.TestWarning(
		"export class Foo { /** @constructor */ constructor() {}; };", CheckJSDoc_DISALLOWED_MEMBER_JSDOC)
}

// Java: CheckJsDocTest.testAbstract_method()
func (c *CheckJsDocTest) TestAbstract_method() {
	c.TestSame("class Foo { /** @abstract */ doSomething() {}}")
	c.TestSame(
		`/** @constructor */
var Foo = function() {};
/** @abstract */
Foo.prototype.something = function() {}
`)
	c.TestSame(
		`/** @constructor */
let Foo = function() {};
/** @abstract */
Foo.prototype.something = function() {}
`)
	c.TestSame(
		`/** @constructor */
const Foo = function() {};
/** @abstract */
Foo.prototype.something = function() {}
`)
	c.TestSame(
		`/** @constructor */
const Foo = function() {
  /** @abstract @return {string} */
  this.something;
};
`)
}

// Java: CheckJsDocTest.testAbstract_method_withES6Modules()
func (c *CheckJsDocTest) TestAbstract_method_withES6Modules() {
	c.TestSame("export class Foo { /** @abstract */ doSomething() {}; };")
}

// Java: CheckJsDocTest.testAbstract_getter_setter()
func (c *CheckJsDocTest) TestAbstract_getter_setter() {
	c.TestSame("class Foo { /** @abstract */ get foo() {}}")
	c.TestSame("class Foo { /** @abstract */ set foo(val) {}}")
	c.TestWarning("class Foo { /** @abstract */ static get foo() {}}", CheckJSDoc_MISPLACED_ANNOTATION)
	c.TestWarning("class Foo { /** @abstract */ static set foo(val) {}}", CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testValidAbstract_getter_setter_withES6Modules()
func (c *CheckJsDocTest) TestValidAbstract_getter_setter_withES6Modules() {
	c.TestSame("export class Foo { /** @abstract */ get foo() {}; };")
}

// Java: CheckJsDocTest.testInvalidAbstract_getter_setter_withES6Modules()
func (c *CheckJsDocTest) TestInvalidAbstract_getter_setter_withES6Modules() {
	c.TestWarning(
		"export class Foo { /** @abstract */ static get foo() {}; };", CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testAbstract_nonEmptyMethod()
func (c *CheckJsDocTest) TestAbstract_nonEmptyMethod() {
	c.TestWarning("class Foo { /** @abstract */ doSomething() { return 0; }}", CheckJSDoc_MISPLACED_ANNOTATION)
	c.TestWarning(
		`/** @constructor */
var Foo = function() {};
/** @abstract */
Foo.prototype.something = function() { return 0; }
`,
		CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testAbstract_nonEmptyMethod_withES6Modules()
func (c *CheckJsDocTest) TestAbstract_nonEmptyMethod_withES6Modules() {
	c.TestWarning(
		"export class Foo { /** @abstract */ doSomething() { return 0; }; };",
		CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testAbstract_staticMethod()
// (static in Java)
func (c *CheckJsDocTest) TestAbstract_staticMethod() {
	c.TestWarning("class Foo { /** @abstract */ static doSomething() {}}", CheckJSDoc_MISPLACED_ANNOTATION)
	c.TestWarning(
		`/** @constructor */
var Foo = function() {};
/** @abstract */
Foo.something = function() {}
`,
		CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testAbstract_staticMethod_withES6Modules()
// (static in Java)
func (c *CheckJsDocTest) TestAbstract_staticMethod_withES6Modules() {
	c.TestWarning(
		"export class Foo { /** @abstract */ static doSomething() {}; };", CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testAbstract_class()
func (c *CheckJsDocTest) TestAbstract_class() {
	c.TestSame("/** @abstract */ class Foo { constructor() {}}")
	c.TestSame("/** @abstract */ exports.Foo = class {}")
	c.TestSame("/** @abstract */ const Foo = class {}")
	c.TestSame("/** @abstract @constructor */ exports.Foo = function() {}")
	c.TestSame("/** @abstract @constructor */ var Foo = function() {};")
	c.TestSame("/** @abstract @constructor */ var Foo = function() { var x = 1; };")
}

// Java: CheckJsDocTest.testAbstract_class_withES6Modules()
func (c *CheckJsDocTest) TestAbstract_class_withES6Modules() {
	c.TestSame("export /** @abstract */ class Foo { constructor() {}; };")
}

// Java: CheckJsDocTest.testAbstract_class_withSymbolNamedMethod()
func (c *CheckJsDocTest) TestAbstract_class_withSymbolNamedMethod() {
	c.TestSame("/** @abstract */ class Foo { /** @abstract */ [Symbol.iterator]() {} }")
}

// Java: CheckJsDocTest.testAbstract()
func (c *CheckJsDocTest) TestAbstract() {
	c.TestWarning("/** @abstract */ var Foo;", CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testInvalidAbstract_withES6Modules()
func (c *CheckJsDocTest) TestInvalidAbstract_withES6Modules() {
	c.TestWarning("export /** @abstract */ var Foo;", CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testAbstract_constructor()
func (c *CheckJsDocTest) TestAbstract_constructor() {
	c.TestWarning("class Foo { /** @abstract */ constructor() {}}", CheckJSDoc_MISPLACED_ANNOTATION)
	// ES5 constructors are treated as class definitions and tested above.
	// This is valid if foo() returns an abstract class constructor
	c.TestSame(
		"/** @constructor */ var C = foo(); /** @abstract */ C.prototype.method = function() {};")
}

// Java: CheckJsDocTest.testInvalidAbstract_constructor_withES6Modules()
func (c *CheckJsDocTest) TestInvalidAbstract_constructor_withES6Modules() {
	c.TestWarning("export class Foo { /** @abstract */ constructor() {}; };", CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testValidAbstract_constructor_withES6Modules()
func (c *CheckJsDocTest) TestValidAbstract_constructor_withES6Modules() {
	c.TestSame(
		`export /** @constructor */ var C = foo();
/** @abstract */ C.prototype.method = function() {};
`)
}

// Java: CheckJsDocTest.testAbstract_field()
func (c *CheckJsDocTest) TestAbstract_field() {
	c.TestWarning("class Foo { constructor() { /** @abstract */ this.x = 1;}}", CheckJSDoc_MISPLACED_ANNOTATION)
	c.TestWarning(
		`/** @constructor */
var Foo = function() {
  /** @abstract */
  this.x = 1;
};
`,
		CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testAbstract_field_withES6Modules()
func (c *CheckJsDocTest) TestAbstract_field_withES6Modules() {
	c.TestWarning(
		"export class Foo { constructor() { /** @abstract */ this.x = 1; }; };",
		CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testAbstract_var()
func (c *CheckJsDocTest) TestAbstract_var() {
	c.TestWarning("class Foo { constructor() {/** @abstract */ var x = 1;}}", CheckJSDoc_MISPLACED_ANNOTATION)
	c.TestWarning(
		`/** @constructor */
var Foo = function() {
  /** @abstract */
  var x = 1;
};
`,
		CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testAbstract_var_withES6Modules()
func (c *CheckJsDocTest) TestAbstract_var_withES6Modules() {
	c.TestWarning(
		"export class Foo { constructor() { /** @abstract */ var x = 1; }; };",
		CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testAbstract_function()
func (c *CheckJsDocTest) TestAbstract_function() {
	c.TestWarning(
		"class Foo { constructor() {/** @abstract */ var x = function() {};}}",
		CheckJSDoc_MISPLACED_ANNOTATION)
	c.TestWarning(
		`/** @constructor */
var Foo = function() {
  /** @abstract */
  var x = function() {};
};
`,
		CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testAbstract_function_withES6Modules()
func (c *CheckJsDocTest) TestAbstract_function_withES6Modules() {
	c.TestWarning(
		"export class Foo { constructor() { /** @abstract */ var x = function() {}; }}",
		CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testInlineJSDoc()
func (c *CheckJsDocTest) TestInlineJSDoc() {
	c.TestSame("function f(/** string */ x) {}")
	c.TestSame("function f(/** @type {string} */ x) {}")
	c.TestSame("var /** string */ x = 'x';")
	c.TestSame("var /** @type {string} */ x = 'x';")
	c.TestSame("var /** string */ x, /** number */ y;")
	c.TestSame("var /** @type {string} */ x, /** @type {number} */ y;")
}

// Java: CheckJsDocTest.testInlineJSDoc_withES6Modules()
func (c *CheckJsDocTest) TestInlineJSDoc_withES6Modules() {
	c.TestSame("export function f(/** string */ x) {}")
}

// Java: CheckJsDocTest.testFunctionJSDocOnStubDeclarations()
func (c *CheckJsDocTest) TestFunctionJSDocOnStubDeclarations() {
	c.TestSame("/** @return {string} */ ns.my.abstract.ctor.method;")
	// should this only be allowed on well-known Symbols ?
	c.TestSame("/** @return {string} */ ns.my.abstract.ctor['method'];")
}

// Java: CheckJsDocTest.testFunctionJSDocOnMethods()
func (c *CheckJsDocTest) TestFunctionJSDocOnMethods() {
	c.TestSame("class Foo { /** @return {?} */ bar() {} }")
	c.TestSame("class Foo { /** @return {?} */ static bar() {} }")
	c.TestSame("class Foo { /** @return {?} */ get bar() {} }")
	c.TestSame("class Foo { /** @param {?} x */ set bar(x) {} }")
	c.TestSame("class Foo { /** @return {?} */ [bar]() {} }")
	c.TestSame("class Foo { /** @return {?} */ static [bar]() {} }")
	c.TestSame("class Foo { /** @return {?} */ get [bar]() {} }")
	c.TestSame("class Foo { /** @return {?} x */ set [bar](x) {} }")
}

// Java: CheckJsDocTest.testValidFunctionJSDocOnMethods_withES6Modules()
func (c *CheckJsDocTest) TestValidFunctionJSDocOnMethods_withES6Modules() {
	c.TestSame("export class Foo { /** @return {?} */ bar() {} }")
}

// Java: CheckJsDocTest.testObjectLiterals()
func (c *CheckJsDocTest) TestObjectLiterals() {
	c.TestSame("var o = { /** @type {?} */ x: y };")
	c.TestWarning("var o = { x: /** @type {?} */ y };", CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testValidObjectLiterals_withES6Modules()
func (c *CheckJsDocTest) TestValidObjectLiterals_withES6Modules() {
	c.TestSame("export var o = { /** @type {?} */ x: y };")
}

// Java: CheckJsDocTest.testInvalidObjectLiterals_withES6Modules()
func (c *CheckJsDocTest) TestInvalidObjectLiterals_withES6Modules() {
	c.TestWarning("export var o = { x: /** @type {?} */ y };", CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testMethodsOnObjectLiterals()
func (c *CheckJsDocTest) TestMethodsOnObjectLiterals() {
	c.TestSame("var x = { /** @return {?} */ foo() {} };")
	c.TestSame("var x = { /** @return {?} */ [foo]() {} };")
	c.TestSame("var x = { /** @return {?} */ foo: someFn };")
	c.TestSame("var x = { /** @return {?} */ [foo]: someFn };")
}

// Java: CheckJsDocTest.testMethodsOnObjectLiterals_withES6Modules()
func (c *CheckJsDocTest) TestMethodsOnObjectLiterals_withES6Modules() {
	c.TestSame("export var x = { /** @return {?} */ foo() {} };")
}

// Java: CheckJsDocTest.testJSDocFunctionNodeAttachment()
func (c *CheckJsDocTest) TestJSDocFunctionNodeAttachment() {
	c.TestWarning(
		`var a = /** @param {number} index */5;
/** @return boolean */function f(index){}
`,
		CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testJSDocFunctionNodeAttachment_withES6Modules()
func (c *CheckJsDocTest) TestJSDocFunctionNodeAttachment_withES6Modules() {
	c.TestWarning(
		`export var a = /** @param {number} index */ 5;
export /** @return boolean */ function f(index){}
`,
		CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testJSDocDescAttachment()
func (c *CheckJsDocTest) TestJSDocDescAttachment() {
	c.TestWarning(
		"function f() { return /** @type {string} */ (g(1 /** @desc x */)); };",
		CheckJSDoc_MISPLACED_MSG_ANNOTATION)
	c.TestWarning("/** @desc Foo. */ var bar = goog.getMsg('hello');", CheckJSDoc_MISPLACED_MSG_ANNOTATION)
	c.TestWarning("/** @desc Foo. */ x.y.z.bar = goog.getMsg('hello');", CheckJSDoc_MISPLACED_MSG_ANNOTATION)
	c.TestWarning("var msgs = {/** @desc x */ x: goog.getMsg('x')}", CheckJSDoc_MISPLACED_MSG_ANNOTATION)
	c.TestWarning("/** @desc Foo. */ bar = goog.getMsg('x');", CheckJSDoc_MISPLACED_MSG_ANNOTATION)
	c.TestWarning("/** @desc Foo. */ var {bar} = goog.getMsg('x');", CheckJSDoc_MISPLACED_MSG_ANNOTATION)
	c.TestWarning("/** @desc Foo. */ let {bar} = goog.getMsg('x');", CheckJSDoc_MISPLACED_MSG_ANNOTATION)
	c.TestWarning("/** @desc Foo. */ const {bar} = goog.getMsg('x');", CheckJSDoc_MISPLACED_MSG_ANNOTATION)
	c.TestWarning(
		"var bar;\n/** @desc Foo. */ ({bar} = goog.getMsg('x'));", CheckJSDoc_MISPLACED_MSG_ANNOTATION)
	// allow in TS code
	c.TestNoWarning(
		c.Srcs(
			SourceFile_FromCode(
				"foo.closure.js", "/** @desc Foo. */ var bar = goog.getMsg('hello');")))
}

// Java: CheckJsDocTest.testJSDocDescInExterns()
func (c *CheckJsDocTest) TestJSDocDescInExterns() {
	c.TestWarning("/** @desc Foo. */ x.y.z.MSG_bar;", CheckJSDoc_MISPLACED_MSG_ANNOTATION)
	c.TestSame(c.Externs("/** @desc Foo. */ x.y.z.MSG_bar;"), c.Srcs(""))
}

// Java: CheckJsDocTest.testJSDocTypeAttachment()
func (c *CheckJsDocTest) TestJSDocTypeAttachment() {
	c.TestWarning("function f() {  /** @type {string} */ if (true) return; };", CheckJSDoc_MISPLACED_ANNOTATION)
	c.TestWarning("function f() {  /** @type {string} */  return; };", CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testJSDocTypeAttachment_withES6Modules()
func (c *CheckJsDocTest) TestJSDocTypeAttachment_withES6Modules() {
	c.TestWarning(
		"export function f() { /** @type {string} */ if (true) return; };", CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testJSDocOnExports()
func (c *CheckJsDocTest) TestJSDocOnExports() {
	c.TestSame(
		`goog.module('foo');
/** @const {!Array<number>} */
exports = [];
`)
}

// Java: CheckJsDocTest.testMisplacedTypeAnnotation1()
func (c *CheckJsDocTest) TestMisplacedTypeAnnotation1() {
	// misuse with COMMA
	c.TestWarning(
		"var o = {}; /** @type {string} */ o.prop1 = 1, o.prop2 = 2;", CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testMisplacedTypeAnnotation2()
func (c *CheckJsDocTest) TestMisplacedTypeAnnotation2() {
	// missing parentheses for the cast.
	c.TestWarning("var o = /** @type {string} */ getValue();", CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testMisplacedTypeAnnotation3()
func (c *CheckJsDocTest) TestMisplacedTypeAnnotation3() {
	// missing parentheses for the cast.
	c.TestWarning("var o = 1 + /** @type {string} */ value;", CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testMisplacedTypeAnnotation4()
func (c *CheckJsDocTest) TestMisplacedTypeAnnotation4() {
	// missing parentheses for the cast.
	c.TestWarning("var o = /** @type {!Array.<string>} */ ['hello', 'you'];", CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testMisplacedTypeAnnotation5()
func (c *CheckJsDocTest) TestMisplacedTypeAnnotation5() {
	// missing parentheses for the cast.
	c.TestWarning("var o = (/** @type {!Foo} */ {});", CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testMisplacedTypeAnnotation6()
func (c *CheckJsDocTest) TestMisplacedTypeAnnotation6() {
	c.TestWarning(
		"var o = /** @type {function():string} */ function() {return 'str';}",
		CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testMisplacedTypeAnnotation7()
func (c *CheckJsDocTest) TestMisplacedTypeAnnotation7() {
	c.TestWarning("var x = /** @type {string} */ y;", CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testMisplacedTypeAnnotationOnStubPropertyOnCall()
func (c *CheckJsDocTest) TestMisplacedTypeAnnotationOnStubPropertyOnCall() {
	c.TestWarning("/** @type {string} */ a.b.c().d;", CheckJSDoc_MISPLACED_ANNOTATION)
	c.TestWarning("/** @typedef {string} */ a.b.c().d;", CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testMisplacedTypeAnnotation_withES6Modules()
func (c *CheckJsDocTest) TestMisplacedTypeAnnotation_withES6Modules() {
	c.TestWarning(
		"export var o = {}; /** @type {string} */ o.prop1 = 1, o.prop2 = 2;", CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testAllowedNocollapseAnnotation1()
func (c *CheckJsDocTest) TestAllowedNocollapseAnnotation1() {
	c.TestSame("var foo = {}; /** @nocollapse */ foo.bar = true;")
}

// Java: CheckJsDocTest.testAllowedNocollapseAnnotation2()
func (c *CheckJsDocTest) TestAllowedNocollapseAnnotation2() {
	c.TestSame(
		`/** @constructor */ function Foo() {};
var ns = {};
/** @nocollapse */ ns.bar = Foo.prototype.blah;
`)
}

// Java: CheckJsDocTest.testAllowedNoCollapseAnnotationOnEsClassMember()
func (c *CheckJsDocTest) TestAllowedNoCollapseAnnotationOnEsClassMember() {
	c.TestSame(
		`class Foo {
  /** @nocollapse */ static bar() {}
  /** @nocollapse */ static get bar() {}
}
`)
}

// Java: CheckJsDocTest.testAllowedNoCollapseAnnotation_onEsClassStaticPropAssignment()
func (c *CheckJsDocTest) TestAllowedNoCollapseAnnotation_onEsClassStaticPropAssignment() {
	c.TestSame(
		`class Foo {}
/** @nocollapse */ Foo.foo = true
`)
}

// Java: CheckJsDocTest.testAllowedNocollapseAnnotation_withES6Modules()
func (c *CheckJsDocTest) TestAllowedNocollapseAnnotation_withES6Modules() {
	c.TestSame("export var foo = {}; /** @nocollapse */ foo.bar = true;")
}

// Java: CheckJsDocTest.testMisplacedNocollapseAnnotationOnPrototypeMethod()
func (c *CheckJsDocTest) TestMisplacedNocollapseAnnotationOnPrototypeMethod() {
	c.TestWarning(
		`/** @constructor */
function Foo() {};
/** @nocollapse */
Foo.prototype.bar = function() {};
`,
		CheckJSDoc_MISPLACED_ANNOTATION)
	c.TestWarning(
		`class Foo {}
/** @nocollapse */
Foo.prototype.bar = function() {};
`,
		CheckJSDoc_MISPLACED_ANNOTATION)
	c.TestWarning(
		`class Foo {
  /** @nocollapse */ bar() {}
}
`,
		CheckJSDoc_MISPLACED_ANNOTATION)
	c.TestWarning(
		`class Foo {
  /** @nocollapse */ get bar() {}
}
`,
		CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testMisplacedNocollapseAnnotation_withES6Modules()
func (c *CheckJsDocTest) TestMisplacedNocollapseAnnotation_withES6Modules() {
	c.TestWarning(
		`export /** @constructor */ function foo() {};
/** @nocollapse */ foo.prototype.bar = function() {};
`,
		CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testNocollapseInExterns()
func (c *CheckJsDocTest) TestNocollapseInExterns() {
	c.Test(
		c.Externs("var foo = {}; /** @nocollapse */ foo.bar = true;"),
		c.Srcs("foo.bar;"),
		c.Warning(CheckJSDoc_MISPLACED_ANNOTATION))
}

// Java: CheckJsDocTest.testNocollapseInExterns_withES6Modules()
func (c *CheckJsDocTest) TestNocollapseInExterns_withES6Modules() {
	c.Test(
		c.Externs("export var foo = {}; /** @nocollapse */ foo.bar = true;"),
		c.Srcs("foo.bar;"),
		c.Warning(CheckJSDoc_MISPLACED_ANNOTATION))
}

// Java: CheckJsDocTest.testArrowFuncAsConstructor()
func (c *CheckJsDocTest) TestArrowFuncAsConstructor() {
	c.TestError("/** @constructor */ var a = ()=>{}; var b = a();", CheckJSDoc_ARROW_FUNCTION_AS_CONSTRUCTOR)
	c.TestError("var a = /** @constructor */ ()=>{}; var b = a();", CheckJSDoc_ARROW_FUNCTION_AS_CONSTRUCTOR)
	c.TestError("/** @constructor */ let a = ()=>{}; var b = a();", CheckJSDoc_ARROW_FUNCTION_AS_CONSTRUCTOR)
	c.TestError("/** @constructor */ const a = ()=>{}; var b = a();", CheckJSDoc_ARROW_FUNCTION_AS_CONSTRUCTOR)
	c.TestError("var a; /** @constructor */ a = ()=>{}; var b = a();", CheckJSDoc_ARROW_FUNCTION_AS_CONSTRUCTOR)
}

// Java: CheckJsDocTest.testArrowFuncAsConstructor_withES6Modules()
func (c *CheckJsDocTest) TestArrowFuncAsConstructor_withES6Modules() {
	c.TestError(
		"export /** @constructor */ var a = ()=>{}; var b = a();", CheckJSDoc_ARROW_FUNCTION_AS_CONSTRUCTOR)
}

// Java: CheckJsDocTest.testValidRestParameter()
func (c *CheckJsDocTest) TestValidRestParameter() {
	c.TestNoWarning(c.Srcs("function f(...args) {}"))
	c.TestNoWarning(c.Srcs("/** @param {...string} args */ function f(...args) {}"))
	c.TestNoWarning(c.Srcs("function f(...[x, y, ...args]) {}"))
	c.TestNoWarning(c.Srcs("/** @param {...number} args */ function f(...[x, y, ...args]) {}"))
	c.TestNoWarning(c.Srcs("function f(.../** ...string */ args) {}"))
}

// Java: CheckJsDocTest.testRestParameter_missingVarArgsInParamJSDoc()
func (c *CheckJsDocTest) TestRestParameter_missingVarArgsInParamJSDoc() {
	c.Test(
		c.Srcs("/** @param {string} args */ function f(...args) {}"),
		c.Warning(CheckJSDoc_BAD_REST_PARAMETER_ANNOTATION))
}

// Java: CheckJsDocTest.testRestParameter_missingVarArgsInDestructuringParamJSDoc()
func (c *CheckJsDocTest) TestRestParameter_missingVarArgsInDestructuringParamJSDoc() {
	c.Test(
		c.Srcs("/** @param {number} rest */ function f(...[x, y, ...args]) {}"),
		c.Warning(CheckJSDoc_BAD_REST_PARAMETER_ANNOTATION))
	c.Test(
		c.Srcs(
			`/**
 * @param {number} p
 * @param {string} rest
 */
function f(p, ...[x, y, ...args]) {}
`),
		c.Warning(CheckJSDoc_BAD_REST_PARAMETER_ANNOTATION))
}

// Java: CheckJsDocTest.testRestParameter_missingVarArgsInInlineJSDoc()
func (c *CheckJsDocTest) TestRestParameter_missingVarArgsInInlineJSDoc() {
	c.Test(
		c.Srcs("function f(.../** string */ args) {}"),
		c.Warning(CheckJSDoc_BAD_REST_PARAMETER_ANNOTATION))
}

// Java: CheckJsDocTest.testDefaultParam()
func (c *CheckJsDocTest) TestDefaultParam() {
	c.TestError("function f(/** number */ x=0) {}", CheckJSDoc_DEFAULT_PARAM_MUST_BE_MARKED_OPTIONAL)
	c.TestSame("function f(/** number= */ x=0) {}")
}

// Java: CheckJsDocTest.testInvalidDefaultParam_withES6Modules()
func (c *CheckJsDocTest) TestInvalidDefaultParam_withES6Modules() {
	c.TestError("export function f(/** number */ x=0) {}", CheckJSDoc_DEFAULT_PARAM_MUST_BE_MARKED_OPTIONAL)
}

// Java: CheckJsDocTest.testValidDefaultParam_withES6Modules()
func (c *CheckJsDocTest) TestValidDefaultParam_withES6Modules() {
	c.TestSame("export function f(/** number= */ x=0) {}")
}

// Java: CheckJsDocTest.testBadTemplate(String)
func (c *CheckJsDocTest) TestBadTemplate(code string) {
	c.TestWarning(code, CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testGoodTemplate1()
func (c *CheckJsDocTest) TestGoodTemplate1() {
	c.TestSame("/** @template T */ class C {}")
	c.TestSame(
		`class C { /** @template T 
 @param {T} a
 @param {T} b 
 */
constructor(a,b){} }
`)
	c.TestSame("class C {/** @template T \n @param {T} a\n @param {T} b \n */ method(a,b){} }")
	c.TestSame("/** @template T \n @param {T} a\n @param {T} b\n */ var x = function(a, b){};")
	c.TestSame("/** @constructor @template T */ var x = function(){};")
	c.TestSame("/** @interface @template T */ var x = function(){};")
}

// Java: CheckJsDocTest.testGoodTemplate_withES6Modules()
func (c *CheckJsDocTest) TestGoodTemplate_withES6Modules() {
	c.TestSame(
		`export class C { /** @template T
 @param {T} a
 @param {T} b
 */
constructor(a,b){} }
`)
	c.TestSame("export /** @template T */ function f(/** T */ a) {}")
}

// Java: CheckJsDocTest.testGoodTemplate_FunctionDefinition()
func (c *CheckJsDocTest) TestGoodTemplate_FunctionDefinition() {
	c.TestSame("/** @template T */ function f() {}")
	c.TestSame("/** @template T */ var f = function() {};")
	c.TestSame("/** @template T */ Foo.prototype.f = function() {};")
	c.TestSame("/** @template T */ function f(/** T */ a) {}")
}

// Java: CheckJsDocTest.testBadTemplate_functionInvocation()
func (c *CheckJsDocTest) TestBadTemplate_functionInvocation() {
	c.TestBadTemplate("/** @template T */ foo();")
}

// Java: CheckJsDocTest.testGoodTypedef()
func (c *CheckJsDocTest) TestGoodTypedef() {
	c.TestSame("/** @typedef {string} */ var x;")
	c.TestSame("/** @typedef {string} */ let x;")
	c.TestSame("/** @typedef {string} */ const x = {};")
	c.TestSame("/** @typedef {string} */ a.b.c;")
	c.TestSame("/** @typedef {string} */ a.b.c = {};")
}

// Java: CheckJsDocTest.testBadTypedef_usedAsCast()
func (c *CheckJsDocTest) TestBadTypedef_usedAsCast() {
	c.TestWarning("const s = /** @typedef {?} */ (0);", CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testBadTypedef_onFunction()
func (c *CheckJsDocTest) TestBadTypedef_onFunction() {
	c.TestWarning("/** @typedef {?} */ function foo() {}", CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testBadTypedef_onInstanceProp()
func (c *CheckJsDocTest) TestBadTypedef_onInstanceProp() {
	c.TestWarning(
		"class C { constructor() { /** @typedef {string} */ this.foo = ''; }}",
		CheckJSDoc_MISPLACED_ANNOTATION)
	c.TestWarning(
		"class C { constructor() { this.foo = {}; /** @typedef {string} */ this.foo.bar = ''; }}",
		CheckJSDoc_MISPLACED_ANNOTATION)
	c.TestWarning(
		`class D {}
class C extends D {
  constructor() {
    super();
    /** @typedef {string} */
    super.foo = ''; }
}
`,
		CheckJSDoc_MISPLACED_ANNOTATION)
	c.TestWarning(
		"class C { constructor() { /** @typedef {string} */ this.foo; }}", CheckJSDoc_MISPLACED_ANNOTATION)
	c.TestWarning(
		"/** @constructor */ function C() { /** @typedef {string} */ this.foo; }",
		CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testBadTypedef_onPrototypeProp()
func (c *CheckJsDocTest) TestBadTypedef_onPrototypeProp() {
	c.TestWarning(
		"/** @constructor */ function C() {} /** @typedef {string} */ C.prototype.foo;",
		CheckJSDoc_MISPLACED_ANNOTATION)
	c.TestWarning(
		c.Externs("/** @constructor */ function C() {} /** @typedef {string} */ C.prototype.foo;"),
		c.Srcs(""),
		c.Warning(CheckJSDoc_MISPLACED_ANNOTATION))
	c.TestWarning(
		"/** @constructor */ function C() {} /** @typedef {string} */ C.prototype.foo = 'foobar';",
		CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testBadTypedef_withES6Modules()
func (c *CheckJsDocTest) TestBadTypedef_withES6Modules() {
	c.TestWarning(
		"export /** @typedef {{foo: string}} */ class C { constructor() { this.foo = ''; }}",
		CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testNoSideEffectsAnnotation1()
func (c *CheckJsDocTest) TestNoSideEffectsAnnotation1() {
	c.TestSame("/** @nosideeffects */ function foo() {}")
}

// Java: CheckJsDocTest.testNoSideEffectsAnnotation2()
func (c *CheckJsDocTest) TestNoSideEffectsAnnotation2() {
	c.TestSame("var f = /** @nosideeffects */ function() {}")
}

// Java: CheckJsDocTest.testNoSideEffectsAnnotation3()
func (c *CheckJsDocTest) TestNoSideEffectsAnnotation3() {
	c.TestSame("/** @nosideeffects */ var f = function() {}")
}

// Java: CheckJsDocTest.testNoSideEffectsAnnotation4()
func (c *CheckJsDocTest) TestNoSideEffectsAnnotation4() {
	c.TestSame(
		`var f = function() {};
/** @nosideeffects */ f.x = function() {}
`)
}

// Java: CheckJsDocTest.testNoSideEffectsAnnotation5()
func (c *CheckJsDocTest) TestNoSideEffectsAnnotation5() {
	c.TestSame(
		`var f = function() {};
f.x = /** @nosideeffects */ function() {}
`)
}

// Java: CheckJsDocTest.testNoSideEffectsAnnotations()
func (c *CheckJsDocTest) TestNoSideEffectsAnnotations() {
	c.TestSame("export /** @nosideeffects */ function foo() {}")
}

// Java: CheckJsDocTest.testInvalidModifiesAnnotation()
func (c *CheckJsDocTest) TestInvalidModifiesAnnotation() {
	c.TestError("/** @modifies {this} */ var f = function() {};", CheckJSDoc_INVALID_MODIFIES_ANNOTATION)
}

// Java: CheckJsDocTest.testInvalidModifiesAnnotation_withES6Modules()
func (c *CheckJsDocTest) TestInvalidModifiesAnnotation_withES6Modules() {
	c.TestError("export /** @modifies {this} */ var f = function() {};", CheckJSDoc_INVALID_MODIFIES_ANNOTATION)
}

// Java: CheckJsDocTest.testInvalidDefinesVariableDeclaration()
func (c *CheckJsDocTest) TestInvalidDefinesVariableDeclaration() {
	c.TestError("/** @define {boolean} */ let DEF = true;", CheckJSDoc_INVALID_DEFINE_ON_LET)
	c.TestError("/** @define {number} */ let DEF = 3;", CheckJSDoc_INVALID_DEFINE_ON_LET)
	c.TestSame("/** @define {boolean} */ var DEF = true;")
	c.TestSame("/** @define {boolean} */ const DEF = true;")
	c.TestError("/** @define {boolean} */ let DEF;", CheckJSDoc_INVALID_DEFINE_ON_LET)
	c.TestError("/** @define {boolean} */ let EXTERN_DEF;", CheckJSDoc_INVALID_DEFINE_ON_LET)
	c.TestSame("let a = {}; /** @define {boolean} */ a.B = false;")
}

// Java: CheckJsDocTest.testMisplacedSuppress()
func (c *CheckJsDocTest) TestMisplacedSuppress() {
	c.TestSame("/** @suppress {missingRequire} */ var x = new y.Z();")
	c.TestSame("/** @suppress {missingRequire} */ function f() { var x = new y.Z(); }")
	c.TestSame("/** @suppress {missingRequire} */ var f = function() { var x = new y.Z(); }")
	c.TestSame(
		`var obj = {
  /** @suppress {uselessCode} */
  f: function() {},
}
`)
	c.TestSame(
		`var obj = {
  /** @suppress {uselessCode} */
  f() {},
}
`)
	c.TestSame(
		`var obj = {
  /** @suppress {uselessCode} */
  ['h' + 6]() {},
}
`)
	c.TestSame(
		`class Example {
  /** @suppress {uselessCode} */
  f() {}
}
`)
	c.TestSame(
		`class Example {
  /** @suppress {uselessCode} */
  static f() {}
}
`)
	c.TestSame(
		`class Example {
  /** @suppress {uselessCode} */
  get f() {}
}
`)
	c.TestSame(
		`class Example {
  /**
   * @param {string} val
   * @suppress {uselessCode}
   */
  set f(val) {}
}
`)
	c.TestSame(
		`class Example {
  /** @suppress {uselessCode} */
  ['f' + 7]() {}
}
`)
	c.TestWarning(
		`var obj = {
  /** @suppress {uselessCode} */
  ['h' + 6]: 'hello',
}
`,
		CheckJSDoc_MISPLACED_SUPPRESS)
	c.TestSame("/** @suppress {extraRequire} */ goog.require('unused.Class');")
	c.TestSame("/** @const @suppress {duplicate} */ var google = {};")
	c.TestSame("/** @suppress {const} */ var google = {};")
	c.TestSame("/** @suppress {duplicate} @type {Foo} */ ns.something.foo;")
	c.TestSame("/** @suppress {with} */ with (o) { x; }")
	// Suppressions are not allowed:
	//  * on arbitrary sub-expressions within a statement
	//  * on control-flow structures like loops, if, switch, or on blocks
	c.TestWarning("const {/** @suppress {duplicate} */ Foo} = foo();", CheckJSDoc_MISPLACED_SUPPRESS)
	c.TestWarning("foo(/** @suppress {duplicate} */ ns.x = 7);", CheckJSDoc_MISPLACED_SUPPRESS)
	c.TestWarning("/** @suppress {duplicate} */ switch (0) { default: ; }", CheckJSDoc_MISPLACED_SUPPRESS)
	c.TestWarning("/** @suppress {duplicate} */ do {} while (true);", CheckJSDoc_MISPLACED_SUPPRESS)
	c.TestWarning("/** @suppress {duplicate} */ while (true) {}", CheckJSDoc_MISPLACED_SUPPRESS)
	c.TestWarning("/** @suppress {duplicate} */ for (; true; ) {}", CheckJSDoc_MISPLACED_SUPPRESS)
	c.TestWarning("/** @suppress {duplicate} */ { (0); }", CheckJSDoc_MISPLACED_SUPPRESS)
	c.TestWarning("/** @suppress {duplicate} */ if (true) {}", CheckJSDoc_MISPLACED_SUPPRESS)
	c.TestSame("/** @suppress {uselessCode} */ goog.require('unused.Class');")
	c.TestSame("/** @suppress {checkTypes} */ foo(0);")
	c.TestSame("/** @suppress {visibility} */ a.x_ = 0;")
	c.TestSame("/** @suppress {visibility} */ a.x_ += 0;")
	c.TestSame("/** @suppress {visibility} */ a.x_ *= 0;")
	c.TestSame("/** @suppress {visibility} */ a.x_ /= 0;")
}

// Java: CheckJsDocTest.testClassFieldSuppressed()
func (c *CheckJsDocTest) TestClassFieldSuppressed() {
	c.TestSame(
		`class Example {
  /** @suppress {uselessCode} */
  x = 2;
}
`)
	c.TestSame(
		`class Example {
  /** @suppress {uselessCode} */
  static x = 2
}
`)
	c.TestSame(
		`class Example {
  /** @suppress {uselessCode} */
  ['x'] = 2
}
`)
	c.TestSame(
		`class Example {
  /** @suppress {uselessCode} */
  static ['x'] = 2
}
`)
	c.TestSame(
		`class Example {
  /** @type {string}  */
   x = 2
}
`)
}

// Java: CheckJsDocTest.testImplicitCastOnlyAllowedInExterns()
func (c *CheckJsDocTest) TestImplicitCastOnlyAllowedInExterns() {
	c.TestSame(
		c.Externs(
			`/** @constructor */ function Element() {};
/**
 * @type {string}
 * @implicitCast
 */
Element.prototype.innerHTML;
`),
		c.Srcs(""))
	c.TestWarning(
		`/** @constructor */ function Element() {};
/**
 * @type {string}
 * @implicitCast
 */
Element.prototype.innerHTML;
`,
		TypeCheck_ILLEGAL_IMPLICIT_CAST)
}

// Java: CheckJsDocTest.testJsDocInBlockComments()
func (c *CheckJsDocTest) TestJsDocInBlockComments() {
	c.TestSame("/** @type {X} */ let x;")
	c.TestSame("// @type {X}\nlet x;")
	c.TestWarning("/* @type {X} */ let x;", CheckJSDoc_JSDOC_IN_BLOCK_COMMENT)
	// jsdoc tags contain letters only, no underscores etc.
	c.TestSame("/* @cc_on */ var x = 3;")
	// a jsdoc tag can't be immediately followed by a paren
	c.TestSame("/* @TODO(username) */ var x = 3;")
}

// Java: CheckJsDocTest.testClosurePrimitive()
func (c *CheckJsDocTest) TestClosurePrimitive() {
	c.TestSame("/** @closurePrimitive {asserts.fail} */ function fail() {}")
	c.TestSame("/** @closurePrimitive {asserts.fail} */ goog.asserts.fail = function() {};")
	c.TestSame("/** @closurePrimitive {asserts.fail} */ let fail = function() {}")
	c.TestWarning("/** @fileoverview @closurePrimitive {asserts.fail} */", CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testJsDocOnReturn()
func (c *CheckJsDocTest) TestJsDocOnReturn() {
	// JSDoc on a class defined in return statement is a warning.
	c.TestWarning(
		`/** @return {function(new:EventTarget)} */
function get() {
/** @implements {EventTarget} */
return class {};
}
`,
		CheckJSDoc_JSDOC_ON_RETURN)
	c.TestWarning("function get() {\n/** @enum {string} */\nreturn {A: 'a'};\n}", CheckJSDoc_JSDOC_ON_RETURN)
	c.TestWarning("function get() {\n/** @typedef {string} */\nreturn 'a';\n}", CheckJSDoc_MISPLACED_ANNOTATION)
	// No warning when returning a class annotated with JSDoc.
	c.TestSame(
		`/** @return {function(new:EventTarget)} */
function mixin() {
/** @implements {EventTarget} */
class MyEventTarget {}
return MyEventTarget;
}
`)
	// No warning for regular JSDoc.
	c.TestSame(
		`function f() {
/** Some value. */
return 5;
}
`)
}

// Java: CheckJsDocTest.testConstructorFieldError_withTypeAnnotation()
func (c *CheckJsDocTest) TestConstructorFieldError_withTypeAnnotation() {
	c.TestSame("class C { /** @constructor */ x = function() {}; }")
}

// Java: CheckJsDocTest.testPublicClassField_allowsTypeAnnotation()
func (c *CheckJsDocTest) TestPublicClassField_allowsTypeAnnotation() {
	c.TestSame("class C { /** @type {number} */ x = 2; }")
}

// Java: CheckJsDocTest.testPublicClassComputedField_allowsTypeAnnotation()
func (c *CheckJsDocTest) TestPublicClassComputedField_allowsTypeAnnotation() {
	c.TestSame("class C { /** @type {number} */ [x] = 2; }")
}

// Java: CheckJsDocTest.testPublicClassField_typeDefError()
func (c *CheckJsDocTest) TestPublicClassField_typeDefError() {
	c.TestWarning("class C { /** @typedef {number} */ x = 2;}", CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testPublicClassComputedField_typeDefError()
func (c *CheckJsDocTest) TestPublicClassComputedField_typeDefError() {
	c.TestWarning("class C { /** @typedef {number} */ [x] = 2;}", CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testMangleClosureModuleExportsContentsTypes()
func (c *CheckJsDocTest) TestMangleClosureModuleExportsContentsTypes() {
	// disable parsing anything besides types; otherwise this test case fails because the
	// "sourceComment"s from the actual/expected JSDocInfo do not match
	c.jsdocParsingMode = parsing.Config_JsDocParsing_TYPES_ONLY
	c.Test(
		"/** @type {!module$exports$foo$bar} */ let x;",
		"/** @type {!UnrecognizedType_module$exports$foo$bar} */ let x;")
	c.Test(
		c.Srcs(
			"goog.module('foo.bar'); exports = class {};",
			"/** @type {!module$exports$foo$bar} */ let x;"),
		c.Expected(
			"goog.module('foo.bar'); exports = class {};",
			"/** @type {!UnrecognizedType_module$exports$foo$bar} */ let x;"))
	c.Test(
		"/** @type {!module$exports$foo$bar.A.B} */ let x;",
		"/** @type {!UnrecognizedType_module$exports$foo$bar.A.B} */ let x;")
	c.Test(
		"/** @type {!module$contents$foo$bar_local} */ let x;",
		"/** @type {!UnrecognizedType_module$contents$foo$bar_local} */ let x;")
	c.Test(
		"/** @type {!Array<module$exports$foo$bar>} */ let x;",
		"/** @type {!Array<UnrecognizedType_module$exports$foo$bar>} */ let x;")
}

// Java: CheckJsDocTest.testNameDeclarationAndAssignForAbstractClass()
func (c *CheckJsDocTest) TestNameDeclarationAndAssignForAbstractClass() {
	c.TestSame(
		`/** @abstract */
let A = A_1 = class A {}
`)
}

// Java: CheckJsDocTest.testNameDeclarationAndAssignForTemplatedClass()
func (c *CheckJsDocTest) TestNameDeclarationAndAssignForTemplatedClass() {
	c.TestSame(
		`/** @template T */
let A = A_1 = class A {}
`)
}

// Java: CheckJsDocTest.testIsUsedViaDotConstrucctor()
func (c *CheckJsDocTest) TestIsUsedViaDotConstrucctor() {
	c.TestSame("/** @constructor @usedViaDotConstructor */ function A() {}")
	c.TestSame("class Foo {/** @usedViaDotConstructor */ constructor() {} }")
	c.TestWarning("/** @usedViaDotConstructor */ function A() {}", CheckJSDoc_MISPLACED_ANNOTATION)
}

// Java: CheckJsDocTest.testMisplacedTypeAnnotationCommaExpression()
func (c *CheckJsDocTest) TestMisplacedTypeAnnotationCommaExpression() {
	c.TestSame("var x = obj.method(/** @type {foo} */ (bar), baz);")
	c.TestSame("var x = (/** @type {foo} */ (obj.method()), bar);")
	c.TestSame("var x = (/** @type {foo} */ (bar), baz);")
}
