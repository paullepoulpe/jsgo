// Port of com/google/javascript/jscomp/JSDocInfoPrinterTest.java @ closure-compiler v20261015

/*
 * Copyright 2014 The Closure Compiler Authors.
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

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// jsDocInfoPrinterTest holds the fixture of JSDocInfoPrinterTest.
type jsDocInfoPrinterTest struct {
	t                *testing.T
	builder          *rhino.JSDocInfo_Builder
	jsDocInfoPrinter *JSDocInfoPrinter
}

// Java: JSDocInfoPrinterTest.setUp()
func newJSDocInfoPrinterTest(t *testing.T) *jsDocInfoPrinterTest {
	return &jsDocInfoPrinterTest{
		t:                t,
		builder:          rhino.JSDocInfo_Builder_().ParseDocumentation(),
		jsDocInfoPrinter: NewJSDocInfoPrinterWithPrintDesc( /* useOriginalName= */ false /* printDesc= */, true),
	}
}

// assertPrints is assertThat(jsDocInfoPrinter.print(builder.buildAndReset())).isEqualTo(want).
func (x *jsDocInfoPrinterTest) assertPrints(want string) {
	x.t.Helper()
	info := x.builder.BuildAndReset()
	if got := x.jsDocInfoPrinter.Print(info); got != want {
		x.t.Errorf("print:\n got: %q\nwant: %q", got, want)
	}
}

// jsDocInfoPrinterTestType is new JSTypeExpression(JsDocInfoParser.parseTypeString(s), source).
func jsDocInfoPrinterTestType(s string, source string) *rhino.JSTypeExpression {
	return rhino.NewJSTypeExpression(parsing.JsDocInfoParser_ParseTypeString(s), source)
}

func jsDocInfoPrinterTestSet(names ...string) *jutil.LinkedHashSet[string] {
	return jutil.NewLinkedHashSet(names...)
}

func TestJSDocInfoPrinter_Basic(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordConstancy()
	x.assertPrints("/** @const */ ")
	x.builder.RecordConstructor()
	x.assertPrints("/** @constructor */ ")
	x.builder.RecordSuppressions(jsDocInfoPrinterTestSet("globalThis", "uselessCode"))
	x.assertPrints("/**\n * @suppress {globalThis,uselessCode}\n */\n")
}

func TestJSDocInfoPrinter_Suppressions(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordSuppressionsWithDescription(jsDocInfoPrinterTestSet("globalThis", "uselessCode"), "Common description.")
	x.assertPrints("/**\n * @suppress {globalThis,uselessCode} Common description.\n */\n")
}

func TestJSDocInfoPrinter_Suppressions_multipleLineDescription(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordSuppressionsWithDescription(
		jsDocInfoPrinterTestSet("globalThis", "uselessCode"),
		"Common description.\n More on another line.")
	x.assertPrints("/**\n * @suppress {globalThis,uselessCode} Common description.\n * More on another line.\n */\n")
}

func TestJSDocInfoPrinter_Suppressions_multiple(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordSuppressionsWithDescription(jsDocInfoPrinterTestSet("globalThis", "uselessCode"), "Common description.")
	x.builder.RecordSuppressions(jsDocInfoPrinterTestSet("const")) // has no description
	x.assertPrints("/**\n * @suppress {globalThis,uselessCode} Common description.\n * @suppress {const}\n */\n")
}

func TestJSDocInfoPrinter_Suppressions_multiple_printOrder(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordSuppressions(jsDocInfoPrinterTestSet("const")) // has no description
	x.builder.RecordSuppressionsWithDescription(jsDocInfoPrinterTestSet("uselessCode", "globalThis"), "Common description.")
	// @suppress printed in order in which it is recorded(parsed)
	// warnings inside a suppress printed in natural order for consistency
	x.assertPrints("/**\n * @suppress {const}\n * @suppress {globalThis,uselessCode} Common description.\n */\n")
}

func TestJSDocInfoPrinter_Final(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordFinality()
	x.assertPrints("/** @final */ ")
}

func TestJSDocInfoPrinter_DescTag(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordDescription("foo")
	x.assertPrints("/**\n * @desc foo\n */\n")
}

func TestJSDocInfoPrinter_MultilineDescTag(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordDescription("foo\nbar")
	x.assertPrints("/**\n * @desc foo\n * bar\n */\n")
}

func TestJSDocInfoPrinter_RecordTag(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordImplicitMatch()
	x.assertPrints("/** @record */ ")
}

func TestJSDocInfoPrinter_Template(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordTemplateTypeName("T")
	x.builder.RecordTemplateTypeName("U")
	x.assertPrints("/**\n * @template T\n * @template U\n */\n")
}

func TestJSDocInfoPrinter_TemplateBound_single(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordTemplateTypeNameWithBound("T", jsDocInfoPrinterTestType("!Array<number>", ""))
	x.assertPrints("/**\n * @template {!Array<number>} T\n */\n")
}

func TestJSDocInfoPrinter_TemplateBound_nullabilityIsPreserved(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordTemplateTypeNameWithBound("T", jsDocInfoPrinterTestType("?Array<number>", ""))
	x.assertPrints("/**\n * @template {?Array<number>} T\n */\n")
}

func TestJSDocInfoPrinter_TemplateBound_explicitlyOnUnknown_isOmitted(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordTemplateTypeNameWithBound("T", jsDocInfoPrinterTestType("?", ""))
	x.assertPrints("/**\n * @template T\n */\n")
}

func TestJSDocInfoPrinter_TemplatesBound_multiple(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordTemplateTypeNameWithBound("T", jsDocInfoPrinterTestType("!Array<number>", ""))
	x.builder.RecordTemplateTypeNameWithBound("U", jsDocInfoPrinterTestType("boolean", ""))
	x.assertPrints("/**\n * @template {!Array<number>} T\n * @template {boolean} U\n */\n")
}

func TestJSDocInfoPrinter_TemplatesBound_mixedWithUnbounded(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordTemplateTypeNameWithBound("T", jsDocInfoPrinterTestType("!Object", ""))
	x.builder.RecordTemplateTypeName("S")
	x.builder.RecordTemplateTypeName("R")
	x.builder.RecordTemplateTypeNameWithBound("U", jsDocInfoPrinterTestType("*", ""))
	x.builder.RecordTemplateTypeName("Q")
	x.assertPrints("/**\n * @template {!Object} T\n * @template S\n * @template R\n * @template {*} U\n * @template Q\n */\n")
}

func TestJSDocInfoPrinter_TypeTransformationLanguageTemplate(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordTypeTransformation("T", rhino.IR_String("Promise"))
	x.assertPrints("/**\n * @template T := \"Promise\" =:\n */\n")
}

func TestJSDocInfoPrinter_Param(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordParameter("foo", jsDocInfoPrinterTestType("number", "<testParam>"))
	x.builder.RecordParameter("bar", jsDocInfoPrinterTestType("string", "<testParam>"))
	x.assertPrints("/**\n * @param {number} foo\n * @param {string} bar\n */\n")

	x.builder.RecordParameter(
		"foo", rhino.NewJSTypeExpression(rhino.NewNodeWithChild(rhino.Token_EQUALS, rhino.IR_String("number")), "<testParam>"))
	x.assertPrints("/**\n * @param {number=} foo\n */\n")

	x.builder.RecordParameter(
		"foo", rhino.NewJSTypeExpression(rhino.NewNodeWithChild(rhino.Token_ITER_REST, rhino.IR_String("number")), "<testParam>"))
	x.assertPrints("/**\n * @param {...number} foo\n */\n")

	x.builder.RecordParameter(
		"foo", rhino.NewJSTypeExpression(rhino.NewNodeWithChild(rhino.Token_ITER_REST, rhino.IR_Empty()), "<testParam>"))
	x.assertPrints("/**\n * @param {...} foo\n */\n")

	x.builder.RecordParameter("foo", nil)
	x.assertPrints("/**\n * @param foo\n */\n")
}

func TestJSDocInfoPrinter_RecordTypes(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordType(jsDocInfoPrinterTestType("{foo: number}", "<testRecordTypes>"))
	x.assertPrints("/** @type {{foo:number}} */ ")

	x.builder.RecordType(jsDocInfoPrinterTestType("{foo}", "<testRecordTypes>"))
	x.assertPrints("/** @type {{foo}} */ ")

	x.builder.RecordType(jsDocInfoPrinterTestType("{foo, bar}", "<testRecordTypes>"))
	x.assertPrints("/** @type {{foo,bar}} */ ")

	x.builder.RecordType(jsDocInfoPrinterTestType("{foo: number, bar}", "<testRecordTypes>"))
	x.assertPrints("/** @type {{foo:number,bar}} */ ")

	x.builder.RecordType(jsDocInfoPrinterTestType("{foo, bar: number}", "<testRecordTypes>"))
	x.assertPrints("/** @type {{foo,bar:number}} */ ")
}

func TestJSDocInfoPrinter_Typeof(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordType(jsDocInfoPrinterTestType("typeof foo", "<testTypeof>"))
	x.assertPrints("/** @type {typeof foo} */ ")
}

func TestJSDocInfoPrinter_Types(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordReturnType(jsDocInfoPrinterTestType("number|string", "<testTypes>"))
	x.assertPrints("/**\n * @return {(number|string)}\n */\n")

	x.builder.RecordParameter(
		"foo", rhino.NewJSTypeExpression(rhino.NewNodeWithChild(rhino.Token_ITER_REST, rhino.IR_String("number")), "<testTypes>"))
	x.assertPrints("/**\n * @param {...number} foo\n */\n")
	x.builder.RecordTypedef(rhino.NewJSTypeExpression(rhino.NewNode(rhino.Token_QMARK), "<testTypes>"))
	x.assertPrints("/** @typedef {?} */ ")
	x.builder.RecordType(rhino.NewJSTypeExpression(rhino.NewNode(rhino.Token_VOID), "<testTypes>"))
	x.assertPrints("/** @type {void} */ ")

	// Object types
	x.builder.RecordEnumParameterType(jsDocInfoPrinterTestType("{foo:number,bar:string}", "<testTypes>"))
	x.assertPrints("/** @enum {{foo:number,bar:string}} */ ")

	x.builder.RecordEnumParameterType(jsDocInfoPrinterTestType("{foo:(number|string)}", "<testTypes>"))
	x.assertPrints("/** @enum {{foo:(number|string)}} */ ")

	// Nullable/non-nullable types.
	x.builder.RecordType(jsDocInfoPrinterTestType("?Object", "<testTypes>"))
	x.assertPrints("/** @type {?Object} */ ")
	x.builder.RecordType(jsDocInfoPrinterTestType("!Object", "<testTypes>"))
	x.assertPrints("/** @type {!Object} */ ")

	// Array types
	x.builder.RecordType(jsDocInfoPrinterTestType("!Array<(number|string)>", "<testTypes>"))
	x.assertPrints("/** @type {!Array<(number|string)>} */ ")
	x.builder.RecordType(jsDocInfoPrinterTestType("Array", "<testTypes>"))
	x.builder.RecordInlineType()
	x.assertPrints("/** Array */ ")

	// Other template types
	x.builder.RecordType(jsDocInfoPrinterTestType("!Set<number|string>", "<testTypes>"))
	x.assertPrints("/** @type {!Set<(number|string)>} */ ")
	x.builder.RecordType(jsDocInfoPrinterTestType("!Map<!Foo, !Bar<!Baz|string>>", "<testTypes>"))
	x.assertPrints("/** @type {!Map<!Foo,!Bar<(!Baz|string)>>} */ ")
	x.builder.RecordType(jsDocInfoPrinterTestType("Map", "<testTypes>"))
	x.builder.RecordInlineType()
	x.assertPrints("/** Map */ ")
}

func TestJSDocInfoPrinter_Inheritance(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordImplementedInterface(jsDocInfoPrinterTestType("Foo", "<testInheritance>"))
	x.assertPrints("/**\n * @implements {Foo}\n */\n")

	x.builder.RecordImplementedInterface(jsDocInfoPrinterTestType("!Foo", "<testInheritance>"))
	x.assertPrints("/**\n * @implements {Foo}\n */\n")

	x.builder.RecordBaseType(jsDocInfoPrinterTestType("Foo", "<testInheritance>"))
	x.assertPrints("/**\n * @extends {Foo}\n */\n")

	x.builder.RecordBaseType(jsDocInfoPrinterTestType("!Foo", "<testInheritance>"))
	x.builder.RecordImplementedInterface(jsDocInfoPrinterTestType("Bar", "<testInheritance>"))
	x.builder.RecordImplementedInterface(jsDocInfoPrinterTestType("Bar.Baz", "<testInheritance>"))
	x.assertPrints("/**\n * @extends {Foo}\n * @implements {Bar}\n * @implements {Bar.Baz}\n */\n")
}

func TestJSDocInfoPrinter_InterfaceInheritance(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordInterface()
	x.builder.RecordExtendedInterface(jsDocInfoPrinterTestType("Foo", "<testInterfaceInheritance>"))
	x.builder.RecordExtendedInterface(jsDocInfoPrinterTestType("Bar", "<testInterfaceInheritance>"))
	x.assertPrints("/**\n * @interface\n * @extends {Foo}\n * @extends {Bar}\n */\n")
}

func TestJSDocInfoPrinter_Functions(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	for _, tc := range []struct{ typ, want string }{
		{"function()", "/** @type {function()} */ "},
		{"function(foo,bar)", "/** @type {function(foo,bar)} */ "},
		{"function(foo):number", "/** @type {function(foo):number} */ "},
		{"function(new:goog,number)", "/** @type {function(new:goog,number)} */ "},
		{"function(this:number,...)", "/** @type {function(this:number,...)} */ "},
		{"function(...number)", "/** @type {function(...number)} */ "},
		{"function():void", "/** @type {function():void} */ "},
		{"function():number", "/** @type {function():number} */ "},
		{"function(string):number", "/** @type {function(string):number} */ "},
		{"function(this:foo):?", "/** @type {function(this:foo):?} */ "},
	} {
		x.builder.RecordType(jsDocInfoPrinterTestType(tc.typ, "<testFunctions>"))
		x.assertPrints(tc.want)
	}
}

func TestJSDocInfoPrinter_Defines(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordDefineType(jsDocInfoPrinterTestType("string", "<testDefines>"))
	x.assertPrints("/** @define {string} */ ")
}

func TestJSDocInfoPrinter_ConstDefines(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordDefineType(jsDocInfoPrinterTestType("string", "<testDefines>"))
	x.builder.RecordConstancy()
	x.assertPrints("/** @define {string} */ ")
}

func TestJSDocInfoPrinter_BlockDescription(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordBlockDescription("Description of the thing")
	x.assertPrints("/**\n * Description of the thing\n */\n")
}

func TestJSDocInfoPrinter_ParamDescriptions(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordParameter("foo", jsDocInfoPrinterTestType("number", "<testParam>"))
	x.builder.RecordParameter("bar", jsDocInfoPrinterTestType("string", "<testParam>"))
	// The parser will retain leading whitespace for descriptions.
	x.builder.RecordParameterDescription("foo", " A number for foo")
	x.builder.RecordParameterDescription("bar", " A multline\n     description for bar")
	x.assertPrints("/**\n * @param {number} foo A number for foo\n * @param {string} bar A multline\n *     description for bar\n */\n")
}

func TestJSDocInfoPrinter_ReturnDescription(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordReturnType(jsDocInfoPrinterTestType("boolean", "<testReturn>"))
	x.builder.RecordReturnDescription("The return value")
	x.assertPrints("/**\n * @return {boolean} The return value\n */\n")
}

func TestJSDocInfoPrinter_AllDescriptions(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordBlockDescription("Description of the thing")
	x.builder.RecordParameter("foo", jsDocInfoPrinterTestType("number", "<testParam>"))
	x.builder.RecordParameter("bar", jsDocInfoPrinterTestType("string", "<testParam>"))
	x.builder.RecordParameterDescription("foo", " A number for foo")
	x.builder.RecordParameterDescription("bar", " A multline\n     description for bar")
	x.builder.RecordReturnType(jsDocInfoPrinterTestType("boolean", "<testReturn>"))
	x.builder.RecordReturnDescription("The return value")
	x.assertPrints("/**\n * Description of the thing\n *\n * @param {number} foo A number for foo\n" +
		" * @param {string} bar A multline\n *     description for bar\n * @return {boolean} The return value\n */\n")
}

func TestJSDocInfoPrinter_Deprecated(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordDeprecated()
	x.builder.RecordDeprecationReason("See {@link otherClass} for more info.")
	x.builder.RecordType(jsDocInfoPrinterTestType("string", "<testDeprecated>"))
	x.assertPrints("/**\n * @type {string}\n * @deprecated See {@link otherClass} for more info.\n */\n")
}

func TestJSDocInfoPrinter_Deprecated_noReason(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordDeprecated()
	x.builder.RecordType(jsDocInfoPrinterTestType("string", "<testDeprecated>"))
	x.assertPrints("/**\n * @type {string}\n * @deprecated\n */\n")
}

// Tests that a @see is sufficient to populate a JSDocInfo.
func TestJSDocInfoPrinter_JSDocIsPopulated_withSeeReferenceAlone(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordReference("SomeClassName for more details")
	x.assertPrints("/**\n * @see SomeClassName for more details\n */\n")
}

// Tests that an @Author is sufficient to populate a JSDocInfo.
func TestJSDocInfoPrinter_JSDocIsPopulated_withAuthorAlone(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordAuthor("John Doe.")
	x.assertPrints("/**\n * @author John Doe.\n */\n")
}

func TestJSDocInfoPrinter_SeeReferencePrintedWithOtherAnnotations(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordType(jsDocInfoPrinterTestType("string", "<testSee>"))
	x.builder.RecordReference("SomeClassName for more details")
	x.assertPrints("/**\n * @see SomeClassName for more details\n * @type {string}\n */\n")
}

func TestJSDocInfoPrinter_AuthorReferencePrintedWithOtherAnnotations(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordType(jsDocInfoPrinterTestType("string", "<testAuthor>"))
	x.builder.RecordAuthor("John Doe.")
	x.assertPrints("/**\n * @author John Doe.\n * @type {string}\n */\n")
}

func TestJSDocInfoPrinter_Externs(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordExterns()
	x.assertPrints("/** @externs */ ")
}

func TestJSDocInfoPrinter_TypeSummary(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordTypeSummary()
	x.assertPrints("/** @typeSummary */ ")
}

func TestJSDocInfoPrinter_Export(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordExport()
	x.assertPrints("/** @export */ ")
}

func TestJSDocInfoPrinter_Abstract(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordAbstract()
	x.assertPrints("/** @abstract */ ")
}

func TestJSDocInfoPrinter_ImplicitCast(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordImplicitCast()
	x.assertPrints("/** @implicitCast */ ")
}

func TestJSDocInfoPrinter_ClosurePrimitive(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordClosurePrimitiveId("testPrimitive")
	x.assertPrints("/** @closurePrimitive {testPrimitive} */ ")
}

func TestJSDocInfoPrinter_NoCollapse(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordNoCollapse()
	x.assertPrints("/** @nocollapse */ ")
}

func TestJSDocInfoPrinter_NgInject(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordNgInject(true)
	x.assertPrints("/** @ngInject */ ")
}

func TestJSDocInfoPrinter_TsType(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordTsType("():string")
	x.assertPrints("/** @tsType ():string */ ")
}

func TestJSDocInfoPrinter_TsType_multipleTsTypes(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordTsType("():string")
	x.builder.RecordTsType("(x:string):number")
	x.assertPrints("/** @tsType ():string @tsType (x:string):number */ ")
}

func TestJSDocInfoPrinter_Polymer(t *testing.T) {
	x := newJSDocInfoPrinterTest(t)
	x.builder.RecordPolymer()
	x.builder.RecordPolymerBehavior()
	x.builder.RecordMixinFunction()
	x.builder.RecordMixinClass()
	x.builder.RecordCustomElement()
	x.assertPrints("/**\n * @polymer\n * @polymerBehavior\n * @mixinFunction\n * @mixinClass\n * @customElement\n */\n")
}
