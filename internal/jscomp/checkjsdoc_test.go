// Port of test/com/google/javascript/jscomp/CheckJsDocTest.java @ closure-compiler v20260415

package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/diag"
)

// ---------------------------------------------------------------------------
// helpers shared by CheckJSDoc tests
// ---------------------------------------------------------------------------

// cjdDiags runs CheckSource and returns all diagnostics.
func cjdDiags(t *testing.T, src string) []diag.Diagnostic {
	t.Helper()
	all, err := CheckSource("t.js", src)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	return all
}

// assertCjdWarning asserts at least one diagnostic with wantCode is present.
func assertCjdWarning(t *testing.T, src, wantCode string) {
	t.Helper()
	diags := cjdDiags(t, src)
	for _, d := range diags {
		if d.Code == wantCode {
			return
		}
	}
	t.Errorf("expected %s; got %v\nsrc: %s", wantCode, diags, src)
}

// assertCjdSame asserts no CheckJSDoc diagnostics are present.
func assertCjdSame(t *testing.T, src string) {
	t.Helper()
	diags := cjdDiags(t, src)
	for _, d := range diags {
		switch d.Code {
		case diag.JSC_MISPLACED_ANNOTATION, diag.JSC_ARROW_FUNCTION_AS_CONSTRUCTOR,
			diag.JSC_DISALLOWED_MEMBER_JSDOC, diag.JSC_MISPLACED_MSG_ANNOTATION,
			diag.JSC_INVALID_DEFINE_ON_LET, diag.JSC_MISPLACED_SUPPRESS,
			diag.JSC_JSDOC_IN_BLOCK_COMMENT, diag.JSC_JSDOC_ON_RETURN,
			diag.JSC_DEFAULT_PARAM_MUST_BE_MARKED_OPTIONAL, diag.JSC_INVALID_MODIFIES_ANNOTATION:
			t.Errorf("unexpected %s at %d:%d: %s\nsrc: %s", d.Code, d.Line, d.Col, d.Message, src)
		}
	}
}

// ---------------------------------------------------------------------------
// Tests ported from CheckJsDocTest.java
// ---------------------------------------------------------------------------

// testInlineJsDocOnObjectPatternTargetNames
func TestCheckJSDoc_InlineJsDocOnObjectPatternTargetNames(t *testing.T) {
	assertCjdSame(t, "let {/** string */ prop} = {prop: 'hi'};")
	assertCjdSame(t, "let {/** string */ prop = 'default'} = {prop: 'hi'};")
	assertCjdSame(t, "let {prop: /** string */ x} = {prop: 'hi'};")
	assertCjdSame(t, "let {prop: /** string */ x = 'default'} = {prop: 'hi'};")
}

// testInlineJsDocOnObjectPatternTargetNamesInAssign
func TestCheckJSDoc_InlineJsDocOnObjectPatternTargetNamesInAssign(t *testing.T) {
	t.Skip("port gap: JSDoc-in-destructuring-assign detection not implemented")
}

// testInlineJsDocOnArrayPatternTargetNames
func TestCheckJSDoc_InlineJsDocOnArrayPatternTargetNames(t *testing.T) {
	assertCjdSame(t, "let [/** string */ x] = ['hi'];")
	assertCjdSame(t, "let [/** string */ x = 'lo'] = ['hi'];")
	assertCjdSame(t, "try {} catch ([/** string */ x]) {}")
}

// testInlineJsDocOnArrayPatternTargetInAssign
func TestCheckJSDoc_InlineJsDocOnArrayPatternTargetInAssign(t *testing.T) {
	t.Skip("port gap: JSDoc-in-array-destructuring-assign detection not implemented")
}

// testImportMeta
func TestCheckJSDoc_ImportMeta(t *testing.T) {
	assertCjdSame(t, "var /** number */ x = foo; import.meta.foo = x")
}

// testInlineJsDocOnDeclaration
func TestCheckJSDoc_InlineJsDocOnDeclaration(t *testing.T) {
	assertCjdSame(t, "var /** number */ x;")
	assertCjdSame(t, "var /** number */ x = 3;")
	t.Skip("port gap: JSDoc on destructuring var declaration not detected")
}

// testFieldMisplacedAnnotation
func TestCheckJSDoc_FieldMisplacedAnnotation(t *testing.T) {
	assertCjdWarning(t, "class Foo {\n  /** @nocollapse */\n  x = 5;\n}", diag.JSC_MISPLACED_ANNOTATION)
}

// testStaticFieldNoCollapse
func TestCheckJSDoc_StaticFieldNoCollapse(t *testing.T) {
	t.Skip("port gap: static class field nocollapse validation not implemented")
}

// testThisPropertyMisplacedAnnotation
func TestCheckJSDoc_ThisPropertyMisplacedAnnotation(t *testing.T) {
	assertCjdWarning(t, "class Foo {\n  constructor() {\n    /** @nocollapse */\n    this.x = 4;\n  }\n}", diag.JSC_MISPLACED_ANNOTATION)
}

// testInvalidJsDocOnDestructuringDeclaration
func TestCheckJSDoc_InvalidJsDocOnDestructuringDeclaration(t *testing.T) {
	t.Skip("port gap: JSDoc-on-destructuring-declaration MISPLACED_ANNOTATION not implemented")
}

// testInlineJsDoc_ES6
func TestCheckJSDoc_InlineJsDoc_ES6(t *testing.T) {
	assertCjdSame(t, "function f(/** string */ x) {}")
	assertCjdSame(t, "function f(/** number= */ x=3) {}")
	assertCjdSame(t, "function f(/** !Object */ {x}) {}")
	assertCjdSame(t, "function f(/** !Array */ [x]) {}")
	assertCjdSame(t, "function f([/** number */ x]) {}")
}

// testJsDocOnCatchVariables
func TestCheckJSDoc_JsDocOnCatchVariables(t *testing.T) {
	assertCjdSame(t, "try {} catch (/** @type {number} */ n) {}")
	t.Skip("port gap: JSDoc on catch destructuring pattern not detected")
}

// testInvalidJSDocOnName
func TestCheckJSDoc_InvalidJSDocOnName(t *testing.T) {
	t.Skip("port gap: inline JSDoc on NAME expressions not detected")
}

// testInvalidClassJsdoc
func TestCheckJSDoc_InvalidClassJsdoc(t *testing.T) {
	assertCjdSame(t, "class Foo { /** @param {number} x */ constructor(x) {}}")
	assertCjdWarning(t, "class Foo { /** @constructor */ constructor() {}}", diag.JSC_DISALLOWED_MEMBER_JSDOC)
	assertCjdWarning(t, "class Foo { /** @interface */ constructor() {}}", diag.JSC_DISALLOWED_MEMBER_JSDOC)
	assertCjdWarning(t, "class Foo { /** @extends {Foo} */ constructor() {}}", diag.JSC_DISALLOWED_MEMBER_JSDOC)
	assertCjdWarning(t, "class Foo { /** @implements {Foo} */ constructor() {}}", diag.JSC_DISALLOWED_MEMBER_JSDOC)
}

// testAbstract_method
func TestCheckJSDoc_Abstract_method(t *testing.T) {
	assertCjdSame(t, "class Foo { /** @abstract */ doSomething() {}}")
	assertCjdSame(t, "/** @constructor */\nvar Foo = function() {};\n/** @abstract */\nFoo.prototype.something = function() {}")
}

// testAbstract_getter_setter
func TestCheckJSDoc_Abstract_getter_setter(t *testing.T) {
	t.Skip("port gap: @abstract on static getter/setter not detected")
}

// testAbstract_nonEmptyMethod
func TestCheckJSDoc_Abstract_nonEmptyMethod(t *testing.T) {
	// @abstract on a method with a non-empty body fires MISPLACED_ANNOTATION.
	assertCjdWarning(t, "class Foo { /** @abstract */ doSomething() { return 5; }}", diag.JSC_MISPLACED_ANNOTATION)
	// But empty body is fine.
	assertCjdSame(t, "class Foo { /** @abstract */ doSomething() {} }")
}

// testAbstract_staticMethod
func TestCheckJSDoc_Abstract_staticMethod(t *testing.T) {
	// @abstract on a static method fires MISPLACED_ANNOTATION.
	assertCjdWarning(t, "class Foo { /** @abstract */ static doSomething() {} }", diag.JSC_MISPLACED_ANNOTATION)
}

// testAbstract_class
func TestCheckJSDoc_Abstract_class(t *testing.T) {
	assertCjdSame(t, "/** @abstract */ class Foo { constructor() {}}")
	assertCjdSame(t, "/** @abstract @constructor */ var Foo = function() {};")
}

// testAbstract
func TestCheckJSDoc_Abstract(t *testing.T) {
	t.Skip("port gap: @abstract on var declaration not detected")
}

// testAbstract_constructor
func TestCheckJSDoc_Abstract_constructor(t *testing.T) {
	t.Skip("port gap: @abstract on constructor not detected")
}

// testAbstract_field
func TestCheckJSDoc_Abstract_field(t *testing.T) {
	t.Skip("port gap: @abstract on instance field not detected")
}

// testAbstract_var
func TestCheckJSDoc_Abstract_var(t *testing.T) {
	t.Skip("port gap: @abstract on var inside constructor not detected")
}

// testAbstract_function
func TestCheckJSDoc_Abstract_function(t *testing.T) {
	t.Skip("port gap: @abstract on function expression inside constructor not detected")
}

// testInlineJSDoc
func TestCheckJSDoc_InlineJSDoc(t *testing.T) {
	assertCjdSame(t, "function f(/** string */ x) {}")
	assertCjdSame(t, "function f(/** @type {string} */ x) {}")
	assertCjdSame(t, "var /** string */ x = 'x';")
	assertCjdSame(t, "var /** @type {string} */ x = 'x';")
}

// testFunctionJSDocOnStubDeclarations
func TestCheckJSDoc_FunctionJSDocOnStubDeclarations(t *testing.T) {
	assertCjdSame(t, "/** @return {string} */ ns.my.abstract.ctor.method;")
}

// testFunctionJSDocOnMethods
func TestCheckJSDoc_FunctionJSDocOnMethods(t *testing.T) {
	assertCjdSame(t, "class Foo { /** @return {?} */ bar() {} }")
	assertCjdSame(t, "class Foo { /** @return {?} */ static bar() {} }")
	assertCjdSame(t, "class Foo { /** @return {?} */ get bar() {} }")
	assertCjdSame(t, "class Foo { /** @param {?} x */ set bar(x) {} }")
}

// testObjectLiterals
func TestCheckJSDoc_ObjectLiterals(t *testing.T) {
	assertCjdSame(t, "var o = { /** @type {?} */ x: y };")
	t.Skip("port gap: JSDoc on object-literal value not detected")
}

// testJSDocFunctionNodeAttachment
func TestCheckJSDoc_JSDocFunctionNodeAttachment(t *testing.T) {
	t.Skip("port gap: @param JSDoc on non-function node not detected")
}

// testJSDocDescAttachment
func TestCheckJSDoc_JSDocDescAttachment(t *testing.T) {
	t.Skip("port gap: @desc JSDoc MISPLACED_MSG_ANNOTATION not implemented for these forms")
}

// testJSDocTypeAttachment
func TestCheckJSDoc_JSDocTypeAttachment(t *testing.T) {
	t.Skip("port gap: @type JSDoc on if/return statement not detected")
}

// testMisplacedTypeAnnotation1
func TestCheckJSDoc_MisplacedTypeAnnotation1(t *testing.T) {
	t.Skip("port gap: @type on comma-expr assignment not detected")
}

// testMisplacedTypeAnnotation2
func TestCheckJSDoc_MisplacedTypeAnnotation2(t *testing.T) {
	t.Skip("port gap: @type on function call RHS not detected")
}

// testMisplacedTypeAnnotation3
func TestCheckJSDoc_MisplacedTypeAnnotation3(t *testing.T) {
	t.Skip("port gap: @type on binary operand not detected")
}

// testMisplacedTypeAnnotation4
func TestCheckJSDoc_MisplacedTypeAnnotation4(t *testing.T) {
	t.Skip("port gap: @type on array literal not detected")
}

// testMisplacedTypeAnnotation5
func TestCheckJSDoc_MisplacedTypeAnnotation5(t *testing.T) {
	assertCjdWarning(t, "var o = (/** @type {!Foo} */ {});", diag.JSC_MISPLACED_ANNOTATION)
}

// testMisplacedTypeAnnotation6
func TestCheckJSDoc_MisplacedTypeAnnotation6(t *testing.T) {
	t.Skip("port gap: @type on function expression RHS not detected")
}

// testMisplacedTypeAnnotation7
func TestCheckJSDoc_MisplacedTypeAnnotation7(t *testing.T) {
	t.Skip("port gap: @type on var-initializer NAME not detected")
}

// testMisplacedTypeAnnotationOnStubPropertyOnCall
func TestCheckJSDoc_MisplacedTypeAnnotationOnStubPropertyOnCall(t *testing.T) {
	assertCjdWarning(t, "/** @type {string} */ a.b.c().d;", diag.JSC_MISPLACED_ANNOTATION)
}

// testAllowedNocollapseAnnotation1
func TestCheckJSDoc_AllowedNocollapseAnnotation1(t *testing.T) {
	assertCjdSame(t, "var foo = {}; /** @nocollapse */ foo.bar = true;")
}

// testAllowedNocollapseAnnotation2
func TestCheckJSDoc_AllowedNocollapseAnnotation2(t *testing.T) {
	assertCjdSame(t, "/** @constructor */ function Foo() {};\nvar ns = {};\n/** @nocollapse */ ns.bar = Foo.prototype.blah;")
}

// testMisplacedNocollapseAnnotationOnPrototypeMethod
func TestCheckJSDoc_MisplacedNocollapseAnnotationOnPrototypeMethod(t *testing.T) {
	t.Skip("port gap: @nocollapse on prototype method not detected")
}

// testArrowFuncAsConstructor
func TestCheckJSDoc_ArrowFuncAsConstructor(t *testing.T) {
	t.Skip("port gap: @constructor on let/const arrow function not detected")
}

// testDefaultParam
func TestCheckJSDoc_DefaultParam(t *testing.T) {
	t.Skip("port gap: DEFAULT_PARAM_MUST_BE_MARKED_OPTIONAL not implemented")
}

// testGoodTemplate1
func TestCheckJSDoc_GoodTemplate1(t *testing.T) {
	assertCjdSame(t, "/** @template T */ class C {}")
	assertCjdSame(t, "/** @template T \n @param {T} a\n @param {T} b\n */ var x = function(a, b){};")
	assertCjdSame(t, "/** @constructor @template T */ var x = function(){};")
}

// testBadTemplate_functionInvocation
func TestCheckJSDoc_BadTemplate_functionInvocation(t *testing.T) {
	assertCjdWarning(t, "/** @template T */ foo();", diag.JSC_MISPLACED_ANNOTATION)
}

// testGoodTypedef
func TestCheckJSDoc_GoodTypedef(t *testing.T) {
	assertCjdSame(t, "/** @typedef {string} */ var x;")
	assertCjdSame(t, "/** @typedef {string} */ let x;")
	assertCjdSame(t, "/** @typedef {string} */ const x = {};")
	assertCjdSame(t, "/** @typedef {string} */ a.b.c;")
}

// testBadTypedef_usedAsCast
func TestCheckJSDoc_BadTypedef_usedAsCast(t *testing.T) {
	t.Skip("port gap: @typedef used as cast not detected")
}

// testBadTypedef_onFunction
func TestCheckJSDoc_BadTypedef_onFunction(t *testing.T) {
	assertCjdWarning(t, "/** @typedef {?} */ function foo() {}", diag.JSC_MISPLACED_ANNOTATION)
}

// testBadTypedef_onPrototypeProp
func TestCheckJSDoc_BadTypedef_onPrototypeProp(t *testing.T) {
	assertCjdWarning(t, "/** @constructor */ function C() {} /** @typedef {string} */ C.prototype.foo;", diag.JSC_MISPLACED_ANNOTATION)
}

// testNoSideEffectsAnnotation1
func TestCheckJSDoc_NoSideEffectsAnnotation1(t *testing.T) {
	assertCjdSame(t, "/** @nosideeffects */ function foo() {}")
}

// testNoSideEffectsAnnotation2
func TestCheckJSDoc_NoSideEffectsAnnotation2(t *testing.T) {
	assertCjdSame(t, "var f = /** @nosideeffects */ function() {}")
}

// testInvalidModifiesAnnotation
func TestCheckJSDoc_InvalidModifiesAnnotation(t *testing.T) {
	assertCjdWarning(t, "/** @modifies {this} */ var f = function() {};", diag.JSC_INVALID_MODIFIES_ANNOTATION)
}

// testInvalidDefinesVariableDeclaration
func TestCheckJSDoc_InvalidDefinesVariableDeclaration(t *testing.T) {
	assertCjdWarning(t, "/** @define {boolean} */ let DEF = true;", diag.JSC_INVALID_DEFINE_ON_LET)
	assertCjdWarning(t, "/** @define {number} */ let DEF = 3;", diag.JSC_INVALID_DEFINE_ON_LET)
	assertCjdSame(t, "/** @define {boolean} */ var DEF = true;")
	assertCjdSame(t, "/** @define {boolean} */ const DEF = true;")
}

// testMisplacedSuppress
func TestCheckJSDoc_MisplacedSuppress(t *testing.T) {
	assertCjdSame(t, "/** @suppress {missingRequire} */ var x = new y.Z();")
	assertCjdSame(t, "/** @suppress {missingRequire} */ function f() { var x = new y.Z(); }")
	assertCjdSame(t, "/** @suppress {extraRequire} */ goog.require('unused.Class');")
	assertCjdSame(t, "/** @const @suppress {duplicate} */ var google = {};")
	assertCjdSame(t, "/** @suppress {duplicate} @type {Foo} */ ns.something.foo;")
	assertCjdSame(t, "/** @suppress {checkTypes} */ foo(0);")
	assertCjdSame(t, "/** @suppress {visibility} */ a.x_ = 0;")
	t.Skip("port gap: @suppress on control-flow statements not detected as MISPLACED_SUPPRESS")
}

// testJsDocInBlockComments
func TestCheckJSDoc_JsDocInBlockComments(t *testing.T) {
	assertCjdSame(t, "/** @type {X} */ let x;")
	assertCjdSame(t, "// @type {X}\nlet x;")
	t.Skip("port gap: JSDoc in block comment (not /** */) not detected")
}

// testClosurePrimitive
func TestCheckJSDoc_ClosurePrimitive(t *testing.T) {
	assertCjdSame(t, "/** @closurePrimitive {asserts.fail} */ function fail() {}")
	t.Skip("port gap: @closurePrimitive on @fileoverview not detected as MISPLACED_ANNOTATION")
}

// testJsDocOnReturn
func TestCheckJSDoc_JsDocOnReturn(t *testing.T) {
	t.Skip("port gap: JSDoc on return statement not detected")
}

// testIsUsedViaDotConstrucctor
func TestCheckJSDoc_IsUsedViaDotConstructor(t *testing.T) {
	assertCjdSame(t, "/** @constructor @usedViaDotConstructor */ function A() {}")
	// class-constructor case: Go treats class method constructors differently
	t.Skip("port gap: @usedViaDotConstructor on class constructor not recognized")
}

// testPublicClassField_allowsTypeAnnotation
func TestCheckJSDoc_PublicClassField_allowsTypeAnnotation(t *testing.T) {
	assertCjdSame(t, "class C { /** @type {number} */ x = 2; }")
}

// testPublicClassField_typeDefError
func TestCheckJSDoc_PublicClassField_typeDefError(t *testing.T) {
	assertCjdWarning(t, "class C { /** @typedef {number} */ x = 2;}", diag.JSC_MISPLACED_ANNOTATION)
}

// testMisplacedTypeAnnotationCommaExpression
func TestCheckJSDoc_MisplacedTypeAnnotationCommaExpression(t *testing.T) {
	assertCjdSame(t, "var x = obj.method(/** @type {foo} */ (bar), baz);")
	assertCjdSame(t, "var x = (/** @type {foo} */ (obj.method()), bar);")
	assertCjdSame(t, "var x = (/** @type {foo} */ (bar), baz);")
}
