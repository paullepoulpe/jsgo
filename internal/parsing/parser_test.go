// Port of com/google/javascript/jscomp/parsing/ParserTest.java @ closure-compiler v20261015

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

package parsing

// ParserTest extends rhino.testing.BaseJSTypeTestCase only to compare the
// JSDoc types of a few nodes with JSTypes of its registry
// (assertNodeHasJSDocInfoWithJSType); parserTest embeds
// rhinotesting.BaseJSTypeTestCase.

import (
	"testing"
	"time"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/rhinotesting"
)

const parserTest_TRAILING_COMMA_MESSAGE = "Trailing comma is not legal in an ECMA-262 object initializer"

var parserTest_MISSING_GT_MESSAGE = "Bad type annotation. missing closing >" + JsDocInfoParser_BAD_TYPE_WIKI_LINK

var parserTest_UNNECESSARY_BRACES_MESSAGE = "Bad type annotation. braces are not required here" + JsDocInfoParser_BAD_TYPE_WIKI_LINK

const parserTest_NAME_NOT_RECOGNIZED_MESSAGE = "name not recognized due to syntax error."

const parserTest_UNLABELED_BREAK = "unlabelled break must be inside loop or switch"

const parserTest_UNEXPECTED_CONTINUE = "continue must be inside loop"

const parserTest_UNEXPECTED_RETURN = "return must be inside function"

const parserTest_UNEXPECTED_LABELLED_CONTINUE = "continue can only use labels of iteration statements"

const parserTest_UNEXPECTED_YIELD = "yield must be inside generator function"

const parserTest_UNEXPECTED_AWAIT = "await must be inside asynchronous function"
const parserTest_UNEXPECTED_FOR_AWAIT_OF = "'for-await-of' used in a non-async function context"

const parserTest_UNDEFINED_LABEL = "undefined label"

const parserTest_HTML_COMMENT_WARNING = "In some cases, '<!--' and '-->' are treated as a '//' for legacy reasons. " +
	"Removing this from your code is safe for all browsers currently in use."
const parserTest_INVALID_ASSIGNMENT_TARGET = "invalid assignment target"
const parserTest_INVALID_FOR_AWAIT_INIT = "for-await-of statement may not have initializer"
const parserTest_INVALID_FOR_AWAIT_MULT_DECLS = "for-await-of statement may not have more than one variable declaration"

const parserTest_SEMICOLON_EXPECTED = "Semi-colon expected"

const parserTest_INVALID_PRIVATE_ID = "Private identifiers may not be used in this context"
const parserTest_PRIVATE_FIELD_NOT_DEFINED = "Private fields must be declared in an enclosing class"
const parserTest_PRIVATE_METHOD_NOT_DEFINED = "Private methods must be declared in an enclosing class"
const parserTest_PRIVATE_FIELD_DELETED = "Private fields cannot be deleted"

// parserTest holds the fields of ParserTest on top of BaseJSTypeTestCase.
type parserTest struct {
	*rhinotesting.BaseJSTypeTestCase
	t *testing.T

	mode             Config_LanguageMode
	parsingMode      Config_JsDocParsing
	strictMode       Config_StrictMode
	isIdeMode        bool
	expectedFeatures *parser.FeatureSet
}

// newParserTest is setUp().
//
// Java: ParserTest.setUp()
func newParserTest(t *testing.T) *parserTest {
	pt := &parserTest{BaseJSTypeTestCase: rhinotesting.NewBaseJSTypeTestCase(t), t: t}
	pt.mode = Config_LanguageMode_UNSUPPORTED
	pt.parsingMode = Config_JsDocParsing_INCLUDE_DESCRIPTIONS_NO_WHITESPACE
	pt.strictMode = Config_StrictMode_STRICT
	pt.isIdeMode = false
	pt.expectedFeatures = parser.FeatureSet_BARE_MINIMUM
	return pt
}

const parserTest_STRING_CONTINUATIONS_WARNING = "String continuations are not recommended. " +
	"See https://google.github.io/styleguide/jsguide.html#features-strings-no-line-continuations"

// Java: ParserTest.testParseUnescapedLineSep()
func TestParseUnescapedLineSep(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("`\u2028`;")

	pt.expectFeatures(parser.FeatureSet_Feature_UNESCAPED_UNICODE_LINE_OR_PARAGRAPH_SEP)
	pt.parse("\"\u2028\";")
	pt.parse("'\u2028';")
}

// Java: ParserTest.testParseUnescapedParagraphSep()
func TestParseUnescapedParagraphSep(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("`\u2029`;")

	pt.expectFeatures(parser.FeatureSet_Feature_UNESCAPED_UNICODE_LINE_OR_PARAGRAPH_SEP)
	pt.parse("\"\u2029\";")
	pt.parse("'\u2029';")
}

// Java: ParserTest.testOptionalCatchBinding()
func TestOptionalCatchBinding(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_OPTIONAL_CATCH_BINDING)

	pt.parse("try {} catch {}")
	pt.parse("try {} catch {} finally {}")
}

// Java: ParserTest.testOptionalCatchBindingSourceInfo()
func TestOptionalCatchBindingSourceInfo(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_OPTIONAL_CATCH_BINDING)

	result := pt.parse("try {} catch     {}")
	catchNode := result.GetFirstFirstChild().GetNext().GetFirstChild()
	assertNode(t, catchNode).HasToken(rhino.Token_CATCH)
	emptyNode := catchNode.GetFirstChild()
	assertNode(t, emptyNode).HasToken(rhino.Token_EMPTY).HasLength(5) // The length matches the whitespace.
}

// Java: ParserTest.testExponentOperator()
func TestExponentOperator(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("-x**y", "Unary operator '-' requires parentheses before '**'")
	pt.parseError("+x**y", "Unary operator '+' requires parentheses before '**'")
	pt.parseError("!x**y", "Unary operator '!' requires parentheses before '**'")
	pt.parseError("~x**y", "Unary operator '~' requires parentheses before '**'")
	pt.parseError("typeof x**y", "Unary operator 'typeof' requires parentheses before '**'")
	pt.parseError("void x**y", "Unary operator 'void' requires parentheses before '**'")
	pt.parseError("delete x**y", "Unary operator 'delete' requires parentheses before '**'")
	pt.parseError("await x**y", "Unary operator 'await' requires parentheses before '**'")
	pt.parseError(
		"async function f() { await x ** y; }",
		"Unary operator 'await' requires parentheses before '**'")
	pt.parseError(
		"async () => { await x ** y; };",
		"Unary operator 'await' requires parentheses before '**'")
	pt.parseError(
		"({ async f() { await x ** y; } });",
		"Unary operator 'await' requires parentheses before '**'")
	pt.parseError(
		"class C { async m() { await x ** y; } }",
		"Unary operator 'await' requires parentheses before '**'")

	pt.expectFeatures(parser.FeatureSet_Feature_EXPONENT_OP)

	pt.parse("x**y")
	// Parentheses are required for disambiguation when a unary expression is desired as
	// the left operand.
	pt.parse("-(x**y)")
	pt.parse("(-x)**y")
	pt.parse("async function f() { (await x) ** y; }")
	// Parens are not required for unary operator on the right operand
	pt.parse("x**-y")
	pt.parse("async function f() { x ** await y; }")
	pt.parse("x/y**z")

	pt.parse("2 ** 3 > 3")

	pt.mode = Config_LanguageMode_ECMASCRIPT_2015
	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parseWarning("x**y", requiresLanguageModeMessage(parser.FeatureSet_Feature_EXPONENT_OP))
}

// Java: ParserTest.testExponentAssignmentOperator()
func TestExponentAssignmentOperator(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_EXPONENT_OP)
	pt.parse("x**=y;")

	pt.mode = Config_LanguageMode_ECMASCRIPT_2015
	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parseWarning("x**=y;", requiresLanguageModeMessage(parser.FeatureSet_Feature_EXPONENT_OP))
}

// Java: ParserTest.testFunction()
func TestFunction(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("var f = function(x,y,z) { return 0; }")
	pt.parse("function f(x,y,z) { return 0; }")

	pt.isIdeMode = true
	pt.parseError("function f(x y z) {}", "',' expected")
}

// Java: ParserTest.testFunctionTrailingComma()
func TestFunctionTrailingComma(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("var f = function(x,y,z,) {}")
	pt.parse("function f(x,y,z,) {}")
}

// Java: ParserTest.testFunctionTrailingCommaPreES8()
func TestFunctionTrailingCommaPreES8(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT_2016

	pt.parseError("var f = function(x,y,z,) {}", "Invalid trailing comma in formal parameter list")
	pt.parseError("function f(x,y,z,) {}", "Invalid trailing comma in formal parameter list")
}

// Java: ParserTest.testFunctionExtraTrailingComma()
func TestFunctionExtraTrailingComma(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("var f = function(x,y,z,,) {}", "')' expected")
	pt.parseError("function f(x,y,z,,) {}", "')' expected")
}

// Java: ParserTest.testCallTrailingComma()
func TestCallTrailingComma(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("f(x,y,z,);")
}

// Java: ParserTest.testCallTrailingCommaPreES8()
func TestCallTrailingCommaPreES8(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT_2016

	pt.parseError("f(x,y,z,);", "Invalid trailing comma in arguments list")
}

// Java: ParserTest.testCallExtraTrailingComma()
func TestCallExtraTrailingComma(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("f(x,y,z,,);", "')' expected")
}

// Java: ParserTest.testTrailingCommaArray()
func TestTrailingCommaArray(t *testing.T) {
	pt := newParserTest(t)

	arrayLit :=
		pt.parse("[1, 2, 3,]"). // SCRIPT
					GetOnlyChild(). // EXPR_RESULT
					GetOnlyChild()  // ARRAYLIT
	assertNode(t, arrayLit).HasType(rhino.Token_ARRAYLIT)
	assertNode(t, arrayLit).HasTrailingComma()
}

// Java: ParserTest.testTrailingCommaObject()
func TestTrailingCommaObject(t *testing.T) {
	pt := newParserTest(t)

	objectlit :=
		pt.parse("var obj = {a:1,b:2,};"). // SCRIPT
							GetOnlyChild(). // VAR
							GetOnlyChild(). // NAME
							GetOnlyChild()  // OBJECTLIT
	assertNode(t, objectlit).HasType(rhino.Token_OBJECTLIT)
	assertNode(t, objectlit).HasTrailingComma()
}

// Java: ParserTest.testTrailingCommaParamList()
func TestTrailingCommaParamList(t *testing.T) {
	pt := newParserTest(t)

	paramList :=
		pt.parse("function f(a, b,) {}"). // SCRIPT
							GetOnlyChild().  // FUNCTION
							GetSecondChild() // PARAM_LIST
	assertNode(t, paramList).HasType(rhino.Token_PARAM_LIST)
	assertNode(t, paramList).HasTrailingComma()
}

// Java: ParserTest.testTrailingCommaParamListArrow()
func TestTrailingCommaParamListArrow(t *testing.T) {
	pt := newParserTest(t)

	paramList :=
		pt.parse("x = (a, b,) => {};"). // SCRIPT
						GetOnlyChild().   // EXPR_RESULT
						GetOnlyChild().   // ASSIGN
						GetSecondChild(). // FUNCTION
						GetSecondChild()  // PARAM_LIST
	assertNode(t, paramList).HasType(rhino.Token_PARAM_LIST)
	assertNode(t, paramList).HasTrailingComma()

	paramList =
		pt.parse("x = (a,) => {};"). // SCRIPT
						GetOnlyChild().   // EXPR_RESULT
						GetOnlyChild().   // ASSIGN
						GetSecondChild(). // FUNCTION
						GetSecondChild()  // PARAM_LIST
	assertNode(t, paramList).HasType(rhino.Token_PARAM_LIST)
	assertNode(t, paramList).HasTrailingComma()

	paramList =
		pt.parse("x = async (a, b,) => {};"). // SCRIPT
							GetOnlyChild().   // EXPR_RESULT
							GetOnlyChild().   // ASSIGN
							GetSecondChild(). // FUNCTION
							GetSecondChild()  // PARAM_LIST
	assertNode(t, paramList).HasType(rhino.Token_PARAM_LIST)
	assertNode(t, paramList).HasTrailingComma()
}

// Java: ParserTest.testTrailingCommaCallNode()
func TestTrailingCommaCallNode(t *testing.T) {
	pt := newParserTest(t)

	call :=
		pt.parse("f(a, b,);"). // SCRIPT
					GetOnlyChild(). // EXPR_RESULT
					GetOnlyChild()  // CALL
	assertNode(t, call).HasType(rhino.Token_CALL)
	assertNode(t, call).HasTrailingComma()
}

// Java: ParserTest.testTrailingCommaNewNode()
func TestTrailingCommaNewNode(t *testing.T) {
	pt := newParserTest(t)

	call :=
		pt.parse("new f(a, b,);"). // SCRIPT
						GetOnlyChild(). // EXPR_RESULT
						GetOnlyChild()  // NEW
	assertNode(t, call).HasType(rhino.Token_NEW)
	assertNode(t, call).HasTrailingComma()
}

// Java: ParserTest.testTrailingCommaOptChainCallNode()
func TestTrailingCommaOptChainCallNode(t *testing.T) {
	pt := newParserTest(t)

	call :=
		pt.parse("f?.(a, b,);"). // SCRIPT
						GetOnlyChild(). // EXPR_RESULT
						GetOnlyChild()  // OPTCHAIN_CALL
	assertNode(t, call).HasType(rhino.Token_OPTCHAIN_CALL)
	assertNode(t, call).HasTrailingComma()
}

// Java: ParserTest.testArrayWithElisions()
func TestArrayWithElisions(t *testing.T) {
	pt := newParserTest(t)

	arrayLit :=
		pt.parse("[  , 1,   ,]"). // SCRIPT
						GetOnlyChild(). // EXPR_RESULT
						GetOnlyChild()  // ARRAYLIT

	assertNode(t, arrayLit).HasType(rhino.Token_ARRAYLIT)
	assertNode(t, arrayLit).HasXChildren(3)

	assertNode(t, arrayLit.GetFirstChild()).HasType(rhino.Token_EMPTY)
	assertNode(t, arrayLit.GetFirstChild()).HasCharno(3)
	assertNode(t, arrayLit.GetFirstChild()).HasLength(0)

	assertNode(t, arrayLit.GetSecondChild()).HasToken(rhino.Token_NUMBER)

	assertNode(t, arrayLit.GetChildAtIndex(2)).HasType(rhino.Token_EMPTY)
	assertNode(t, arrayLit.GetChildAtIndex(2)).HasLength(0)
	assertNode(t, arrayLit.GetChildAtIndex(2)).HasCharno(10)

	assertNode(t, arrayLit).HasTrailingComma()
}

// Java: ParserTest.testWhile()
func TestWhile(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("while(1) { break; }")
}

// Java: ParserTest.testNestedWhile()
func TestNestedWhile(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("while(1) { while(1) { break; } }")
}

// Java: ParserTest.testBreak()
func TestBreak(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("break;", parserTest_UNLABELED_BREAK)
}

// Java: ParserTest.testContinue()
func TestContinue(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("continue;", parserTest_UNEXPECTED_CONTINUE)
}

// Java: ParserTest.testBreakCrossFunction()
func TestBreakCrossFunction(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("while(1) { var f = function() { break; } }", parserTest_UNLABELED_BREAK)
}

// Java: ParserTest.testBreakCrossFunctionInFor()
func TestBreakCrossFunctionInFor(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("while(1) {for(var f = function () { break; };;) {}}", parserTest_UNLABELED_BREAK)
}

// Java: ParserTest.testBreakInForOf()
func TestBreakInForOf(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY
	pt.expectFeatures(parser.FeatureSet_Feature_FOR_OF)
	pt.parse(
		`for (var x of [1, 2, 3]) {
  if (x == 2) break;
}
`)
}

// Java: ParserTest.testContinueToSwitch()
func TestContinueToSwitch(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("switch(1) {case(1): continue; }", parserTest_UNEXPECTED_CONTINUE)
}

// Java: ParserTest.testContinueToSwitchWithNoCases()
func TestContinueToSwitchWithNoCases(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("switch(1){}")
}

// Java: ParserTest.testContinueToSwitchWithTwoCases()
func TestContinueToSwitchWithTwoCases(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("switch(1){case(1):break;case(2):continue;}", parserTest_UNEXPECTED_CONTINUE)
}

// Java: ParserTest.testContinueToSwitchWithDefault()
func TestContinueToSwitchWithDefault(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("switch(1){case(1):break;case(2):default:continue;}", parserTest_UNEXPECTED_CONTINUE)
}

// Java: ParserTest.testContinueToLabelSwitch()
func TestContinueToLabelSwitch(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("while(1) {a: switch(1) {case(1): continue a; }}", parserTest_UNEXPECTED_LABELLED_CONTINUE)
}

// Java: ParserTest.testContinueOutsideSwitch()
func TestContinueOutsideSwitch(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("b: while(1) { a: switch(1) { case(1): continue b; } }")
}

// Java: ParserTest.testContinueNotCrossFunction1()
func TestContinueNotCrossFunction1(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("a:switch(1){case(1):var f = function(){a:while(1){continue a;}}}")
}

// Java: ParserTest.testContinueNotCrossFunction2()
func TestContinueNotCrossFunction2(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError(
		"a:switch(1){case(1):var f = function(){while(1){continue a;}}}",
		parserTest_UNDEFINED_LABEL+" \"a\"")
}

// Java: ParserTest.testContinueInForOf()
func TestContinueInForOf(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY
	pt.expectFeatures(parser.FeatureSet_Feature_FOR_OF)
	pt.parse(
		`for (var x of [1, 2, 3]) {
  if (x == 2) continue;
}
`)
}

// Java: ParserTest.testVarSourceLocations()
func TestVarSourceLocations(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true

	n := pt.parse("var x, y = 1;")
	var_ := n.GetFirstChild()
	assertNode(t, var_).HasType(rhino.Token_VAR)

	x := var_.GetFirstChild()
	assertNode(t, x).HasType(rhino.Token_NAME)
	assertNode(t, x).HasCharno(jutil.StringLength("var "))

	y := x.GetNext()
	assertNode(t, y).HasType(rhino.Token_NAME)
	assertNode(t, y).HasCharno(jutil.StringLength("var x, "))
	assertNode(t, y).HasLength(jutil.StringLength("y = 1"))
}

// Java: ParserTest.testSourceLocationsNonAscii()
func TestSourceLocationsNonAscii(t *testing.T) {
	pt := newParserTest(t)

	n := pt.parse("'\uc548\ub155\uc138\uacc4!'")
	exprResult := n.GetFirstChild()
	string := exprResult.GetFirstChild()
	assertNode(t, string).HasType(rhino.Token_STRINGLIT)
	assertNode(t, string).HasLength(7) // 2 quotes, plus 5 characters
}

// Java: ParserTest.testNumericSeparatorOnDecimal()
func TestNumericSeparatorOnDecimal(t *testing.T) {
	pt := newParserTest(t)

	result := pt.parse("1_000_000;")
	pt.expectFeatures(parser.FeatureSet_Feature_NUMERIC_SEPARATOR)
	assertNode(t, result).HasType(rhino.Token_SCRIPT)
	exprResult := result.GetOnlyChild()
	assertNode(t, exprResult).HasType(rhino.Token_EXPR_RESULT)
	numberNode := exprResult.GetOnlyChild()
	assertNode(t, numberNode).IsNumber(1000000)
}

// Java: ParserTest.testNumericSeparatorOnBinary()
func TestNumericSeparatorOnBinary(t *testing.T) {
	pt := newParserTest(t)

	result := pt.parse("0b1_0000;")
	pt.expectFeatures(parser.FeatureSet_Feature_NUMERIC_SEPARATOR)
	assertNode(t, result).HasType(rhino.Token_SCRIPT)
	exprResult := result.GetOnlyChild()
	assertNode(t, exprResult).HasType(rhino.Token_EXPR_RESULT)
	numberNode := exprResult.GetOnlyChild()
	assertNode(t, numberNode).IsNumber(16)
}

// Java: ParserTest.testNumericSeparatorOnOctal()
func TestNumericSeparatorOnOctal(t *testing.T) {
	pt := newParserTest(t)

	result := pt.parse("0o01_00;")
	pt.expectFeatures(parser.FeatureSet_Feature_NUMERIC_SEPARATOR)
	assertNode(t, result).HasType(rhino.Token_SCRIPT)
	exprResult := result.GetOnlyChild()
	assertNode(t, exprResult).HasType(rhino.Token_EXPR_RESULT)
	numberNode := exprResult.GetOnlyChild()
	assertNode(t, numberNode).IsNumber(64)
}

// Java: ParserTest.testNumericSeparatorOnHex()
func TestNumericSeparatorOnHex(t *testing.T) {
	pt := newParserTest(t)

	result := pt.parse("0x01_01")
	pt.expectFeatures(parser.FeatureSet_Feature_NUMERIC_SEPARATOR)
	assertNode(t, result).HasType(rhino.Token_SCRIPT)
	exprResult := result.GetOnlyChild()
	assertNode(t, exprResult).HasType(rhino.Token_EXPR_RESULT)
	numberNode := exprResult.GetOnlyChild()
	assertNode(t, numberNode).IsNumber(257)
}

// Java: ParserTest.testNumericSeparatorOnBigInt()
func TestNumericSeparatorOnBigInt(t *testing.T) {
	pt := newParserTest(t)

	bigint :=
		pt.parse("1_000n"). // SCRIPT
					GetOnlyChild(). // EXPR_RESULT
					GetOnlyChild()  // BIGINT
	assertNode(t, bigint).HasType(rhino.Token_BIGINT)
	assertNode(t, bigint).IsBigInt(jutil.NewBigIntegerFromString("1000"))
}

// Java: ParserTest.testTrailingNumericSeparator()
func TestTrailingNumericSeparator(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("1000_", "Trailing numeric separator")
	pt.parseError("0b0001_", "Trailing numeric separator")
	pt.parseError("0o100_", "Trailing numeric separator")
	pt.parseError("0x0F_", "Trailing numeric separator")
	pt.parseError("1000_n", "Trailing numeric separator")
	pt.parseError("0b0001_n", "Trailing numeric separator")
	pt.parseError("0o100_n", "Trailing numeric separator")
	pt.parseError("0x0F_n", "Trailing numeric separator")
}

// Java: ParserTest.testNumericSeparatorWithDecimalPoint()
func TestNumericSeparatorWithDecimalPoint(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("3_.141", "Trailing numeric separator")
	pt.parseError("3._141", "Trailing numeric separator")
	pt.parseError("3._", "Trailing numeric separator")
	pt.parseError("0._141", "Trailing numeric separator")
	pt.parseError("0._", "Trailing numeric separator")
	pt.parseError("3._e1", "Trailing numeric separator")
	pt.parseError("3.1_e1", "Trailing numeric separator")
	pt.parseError("3.141_", "Trailing numeric separator")
	pt.parseError(".141_", "Trailing numeric separator")
	pt.parseError("_.1", parserTest_SEMICOLON_EXPECTED)
	pt.parseError("._141", "primary expression expected")

	result := pt.parse("3.1_41;")
	pt.expectFeatures(parser.FeatureSet_Feature_NUMERIC_SEPARATOR)
	assertNode(t, result).HasType(rhino.Token_SCRIPT)
	exprResult := result.GetOnlyChild()
	assertNode(t, exprResult).HasType(rhino.Token_EXPR_RESULT)
	numberNode := exprResult.GetOnlyChild()
	assertNode(t, numberNode).IsNumber(3.141)
}

// Java: ParserTest.testNumericSeparatorWarning()
func TestNumericSeparatorWarning(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT_2020
	pt.parseWarning("1_000", requiresLanguageModeMessage(parser.FeatureSet_Feature_NUMERIC_SEPARATOR))
}

// Java: ParserTest.testReturn()
func TestReturn(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("function foo() { return 1; }")
	pt.parseError("return;", parserTest_UNEXPECTED_RETURN)
	pt.parseError("return 1;", parserTest_UNEXPECTED_RETURN)
}

// Java: ParserTest.testThrow()
func TestThrow(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("throw Error();")
	pt.parse("throw new Error();")
	pt.parse("throw '';")
	pt.parseError("throw;", "semicolon/newline not allowed after 'throw'")
	pt.parseError("throw\nError();", "semicolon/newline not allowed after 'throw'")
}

// Java: ParserTest.testLabel1()
func TestLabel1(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("foo:bar")
}

// Java: ParserTest.testLabel2()
func TestLabel2(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("{foo:bar}")
}

// Java: ParserTest.testLabel3()
func TestLabel3(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("foo:bar:baz")
}

// Java: ParserTest.testDuplicateLabelWithoutBraces()
func TestDuplicateLabelWithoutBraces(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("foo:foo:bar", "Duplicate label \"foo\"")
}

// Java: ParserTest.testDuplicateLabelWithBraces()
func TestDuplicateLabelWithBraces(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("foo:{bar;foo:baz}", "Duplicate label \"foo\"")
}

// Java: ParserTest.testDuplicateLabelWithFor()
func TestDuplicateLabelWithFor(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("foo:for(;;){foo:bar}", "Duplicate label \"foo\"")
}

// Java: ParserTest.testNonDuplicateLabelSiblings()
func TestNonDuplicateLabelSiblings(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("foo:1;foo:2")
}

// Java: ParserTest.testNonDuplicateLabelCrossFunction()
func TestNonDuplicateLabelCrossFunction(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("foo:(function(){foo:2})")
}

// Java: ParserTest.testLabeledFunctionDeclaration()
func TestLabeledFunctionDeclaration(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError(
		"foo:function f() {}",
		"Lexical declarations are only allowed at top level or inside a block.")
}

// Java: ParserTest.testLabeledClassDeclaration()
func TestLabeledClassDeclaration(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError(
		"foo:class Foo {}",
		"Lexical declarations are only allowed at top level or inside a block.")
}

// Java: ParserTest.testLabeledLetDeclaration()
func TestLabeledLetDeclaration(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError(
		"foo: let x = 0;", "Lexical declarations are only allowed at top level or inside a block.")
}

// Java: ParserTest.testLabeledConstDeclaration()
func TestLabeledConstDeclaration(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError(
		"foo: const x = 0;",
		"Lexical declarations are only allowed at top level or inside a block.")
}

// Java: ParserTest.testMethodNamedStatic()
func TestMethodNamedStatic(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("class C { static(a, b) {} }")
}

// Java: ParserTest.testLinenoCharnoAssign1()
func TestLinenoCharnoAssign1(t *testing.T) {
	pt := newParserTest(t)

	assign := pt.parse("a = b").GetFirstFirstChild()

	assertNode(t, assign).HasType(rhino.Token_ASSIGN)
	assertNode(t, assign).HasLineno(1)
	assertNode(t, assign).HasCharno(0)
}

// Java: ParserTest.testLinenoCharnoAssign2()
func TestLinenoCharnoAssign2(t *testing.T) {
	pt := newParserTest(t)

	assign := pt.parse("\n a.g.h.k    =  45").GetFirstFirstChild()

	assertNode(t, assign).HasType(rhino.Token_ASSIGN)
	assertNode(t, assign).HasLineno(2)
	assertNode(t, assign).HasCharno(1)
}

// Java: ParserTest.testLinenoCharnoCall()
func TestLinenoCharnoCall(t *testing.T) {
	pt := newParserTest(t)

	call := pt.parse("\n foo(123);").GetFirstFirstChild()

	assertNode(t, call).HasType(rhino.Token_CALL)
	assertNode(t, call).HasLineno(2)
	assertNode(t, call).HasCharno(1)
}

// Java: ParserTest.testLinenoCharnoGetProp1()
func TestLinenoCharnoGetProp1(t *testing.T) {
	pt := newParserTest(t)

	getprop := pt.parse("\n foo.bar").GetFirstFirstChild()

	assertNode(t, getprop).HasType(rhino.Token_GETPROP)
	assertNode(t, getprop).HasLineno(2)
	assertNode(t, getprop).HasCharno(5)
	assertNode(t, getprop).HasStringThat().IsEqualTo("bar")
}

// Java: ParserTest.testLinenoCharnoGetProp2()
func TestLinenoCharnoGetProp2(t *testing.T) {
	pt := newParserTest(t)

	getprop := pt.parse(`
 foo.
bar`).GetFirstFirstChild()

	assertNode(t, getprop).HasType(rhino.Token_GETPROP)
	assertNode(t, getprop).HasLineno(3)
	assertNode(t, getprop).HasCharno(0)
	assertNode(t, getprop).HasStringThat().IsEqualTo("bar")
}

// Java: ParserTest.testLinenoCharnoGetelem1()
func TestLinenoCharnoGetelem1(t *testing.T) {
	pt := newParserTest(t)

	call := pt.parse("\n foo[123]").GetFirstFirstChild()

	assertNode(t, call).HasType(rhino.Token_GETELEM)
	assertNode(t, call).HasLineno(2)
	assertNode(t, call).HasCharno(1)
}

// Java: ParserTest.testLinenoCharnoGetelem2()
func TestLinenoCharnoGetelem2(t *testing.T) {
	pt := newParserTest(t)

	call := pt.parse("\n   \n foo()[123]").GetFirstFirstChild()

	assertNode(t, call).HasType(rhino.Token_GETELEM)
	assertNode(t, call).HasLineno(3)
	assertNode(t, call).HasCharno(1)
}

// Java: ParserTest.testLinenoCharnoGetelem3()
func TestLinenoCharnoGetelem3(t *testing.T) {
	pt := newParserTest(t)

	call := pt.parse("\n   \n (8 + kl)[123]").GetFirstFirstChild()

	assertNode(t, call).HasType(rhino.Token_GETELEM)
	assertNode(t, call).HasLineno(3)
	assertNode(t, call).HasCharno(1)
}

// Java: ParserTest.testLinenoCharnoForComparison()
func TestLinenoCharnoForComparison(t *testing.T) {
	pt := newParserTest(t)

	lt := pt.parse("for (; i < j;){}").GetFirstChild().GetSecondChild()

	assertNode(t, lt).HasType(rhino.Token_LT)
	assertNode(t, lt).HasLineno(1)
	assertNode(t, lt).HasCharno(7)
}

// Java: ParserTest.testLinenoCharnoHook()
func TestLinenoCharnoHook(t *testing.T) {
	pt := newParserTest(t)

	n := pt.parse("\n a ? 9 : 0").GetFirstFirstChild()

	assertNode(t, n).HasType(rhino.Token_HOOK)
	assertNode(t, n).HasLineno(2)
	assertNode(t, n).HasCharno(1)
}

// Java: ParserTest.testLinenoCharnoArrayLiteral()
func TestLinenoCharnoArrayLiteral(t *testing.T) {
	pt := newParserTest(t)

	n := pt.parse("\n  [8, 9]").GetFirstFirstChild()

	assertNode(t, n).HasType(rhino.Token_ARRAYLIT)
	assertNode(t, n).HasLineno(2)
	assertNode(t, n).HasCharno(2)

	n = n.GetFirstChild()

	assertNode(t, n).HasType(rhino.Token_NUMBER)
	assertNode(t, n).HasLineno(2)
	assertNode(t, n).HasCharno(3)

	n = n.GetNext()

	assertNode(t, n).HasType(rhino.Token_NUMBER)
	assertNode(t, n).HasLineno(2)
	assertNode(t, n).HasCharno(6)
}

// Java: ParserTest.testLinenoCharnoObjectLiteral()
func TestLinenoCharnoObjectLiteral(t *testing.T) {
	pt := newParserTest(t)

	n := pt.parse(`

 var a = {a:0
,b :1};`).GetFirstFirstChild().GetFirstChild()

	assertNode(t, n).HasType(rhino.Token_OBJECTLIT)
	assertNode(t, n).HasLineno(3)
	assertNode(t, n).HasCharno(9)

	key := n.GetFirstChild()

	assertNode(t, key).HasType(rhino.Token_STRING_KEY)
	assertNode(t, key).HasLineno(3)
	assertNode(t, key).HasCharno(10)

	value := key.GetFirstChild()

	assertNode(t, value).HasType(rhino.Token_NUMBER)
	assertNode(t, value).HasLineno(3)
	assertNode(t, value).HasCharno(12)

	key = key.GetNext()

	assertNode(t, key).HasType(rhino.Token_STRING_KEY)
	assertNode(t, key).HasLineno(4)
	assertNode(t, key).HasCharno(1)

	value = key.GetFirstChild()

	assertNode(t, value).HasType(rhino.Token_NUMBER)
	assertNode(t, value).HasLineno(4)
	assertNode(t, value).HasCharno(4)
}

// Java: ParserTest.testLinenoCharnoObjectLiteralMemberFunction()
func TestLinenoCharnoObjectLiteralMemberFunction(t *testing.T) {
	pt := newParserTest(t)

	n := pt.parse("var a = {\n fn() {} };").GetFirstFirstChild().GetFirstChild()

	assertNode(t, n).HasType(rhino.Token_OBJECTLIT)
	assertNode(t, n).HasLineno(1)
	assertNode(t, n).HasCharno(8)

	// fn() {}
	key := n.GetFirstChild()

	assertNode(t, key).HasType(rhino.Token_MEMBER_FUNCTION_DEF)
	assertNode(t, key).HasLineno(2)
	assertNode(t, key).HasCharno(1)
	assertNode(t, key).HasLength(2) // "fn"

	value := key.GetFirstChild()

	assertNode(t, value).HasType(rhino.Token_FUNCTION)
	assertNode(t, value).HasLineno(2)
	assertNode(t, value).HasCharno(1)
	assertNode(t, value).HasLength(7) // "fn() {}"
}

// Java: ParserTest.testLinenoCharnoEs6Class()
func TestLinenoCharnoEs6Class(t *testing.T) {
	pt := newParserTest(t)

	n := pt.parse(`class C {
  fn1() {}
  static fn2() {}
 };`).GetFirstChild()

	assertNode(t, n).HasType(rhino.Token_CLASS)
	assertNode(t, n).HasLineno(1)
	assertNode(t, n).HasCharno(0)

	members := nodeUtil_getClassMembers(n)
	assertNode(t, members).HasType(rhino.Token_CLASS_MEMBERS)

	// fn1 () {}
	memberFn := members.GetFirstChild()

	assertNode(t, memberFn).HasType(rhino.Token_MEMBER_FUNCTION_DEF)
	assertNode(t, memberFn).HasLineno(2)
	assertNode(t, memberFn).HasCharno(2)
	assertNode(t, memberFn).HasLength(3) // "fn"

	fn := memberFn.GetFirstChild()
	assertNode(t, fn).HasType((rhino.Token_FUNCTION))
	assertNode(t, fn).HasLineno(2)
	assertNode(t, fn).HasCharno(2)
	assertNode(t, fn).HasLength(8) // "fn1() {}"

	// static fn2() {}
	memberFn = memberFn.GetNext()

	assertNode(t, memberFn).HasType(rhino.Token_MEMBER_FUNCTION_DEF)
	assertNode(t, memberFn).HasLineno(3)
	assertNode(t, memberFn).HasCharno(9)
	assertNode(t, memberFn).HasLength(3) // "fn2"
	assertNode(t, memberFn).IsStatic()

	fn = memberFn.GetFirstChild()
	assertNode(t, fn).HasType((rhino.Token_FUNCTION))
	assertNode(t, fn).HasLineno(3)
	assertNode(t, fn).HasCharno(2)
	assertNode(t, fn).HasLength(15) // "static fn2() {}"
}

// Java: ParserTest.testLinenoCharnoAdd()
func TestLinenoCharnoAdd(t *testing.T) {
	pt := newParserTest(t)

	pt.testLinenoCharnoBinop("+")
}

// Java: ParserTest.testLinenoCharnoSub()
func TestLinenoCharnoSub(t *testing.T) {
	pt := newParserTest(t)

	pt.testLinenoCharnoBinop("-")
}

// Java: ParserTest.testLinenoCharnoMul()
func TestLinenoCharnoMul(t *testing.T) {
	pt := newParserTest(t)

	pt.testLinenoCharnoBinop("*")
}

// Java: ParserTest.testLinenoCharnoDiv()
func TestLinenoCharnoDiv(t *testing.T) {
	pt := newParserTest(t)

	pt.testLinenoCharnoBinop("/")
}

// Java: ParserTest.testLinenoCharnoMod()
func TestLinenoCharnoMod(t *testing.T) {
	pt := newParserTest(t)

	pt.testLinenoCharnoBinop("%")
}

// Java: ParserTest.testLinenoCharnoShift()
func TestLinenoCharnoShift(t *testing.T) {
	pt := newParserTest(t)

	pt.testLinenoCharnoBinop("<<")
}

// Java: ParserTest.testLinenoCharnoBinaryAnd()
func TestLinenoCharnoBinaryAnd(t *testing.T) {
	pt := newParserTest(t)

	pt.testLinenoCharnoBinop("&")
}

// Java: ParserTest.testLinenoCharnoAnd()
func TestLinenoCharnoAnd(t *testing.T) {
	pt := newParserTest(t)

	pt.testLinenoCharnoBinop("&&")
}

// Java: ParserTest.testLinenoCharnoBinaryOr()
func TestLinenoCharnoBinaryOr(t *testing.T) {
	pt := newParserTest(t)

	pt.testLinenoCharnoBinop("|")
}

// Java: ParserTest.testLinenoCharnoOr()
func TestLinenoCharnoOr(t *testing.T) {
	pt := newParserTest(t)

	pt.testLinenoCharnoBinop("||")
}

// Java: ParserTest.testLinenoCharnoLt()
func TestLinenoCharnoLt(t *testing.T) {
	pt := newParserTest(t)

	pt.testLinenoCharnoBinop("<")
}

// Java: ParserTest.testLinenoCharnoLe()
func TestLinenoCharnoLe(t *testing.T) {
	pt := newParserTest(t)

	pt.testLinenoCharnoBinop("<=")
}

// Java: ParserTest.testLinenoCharnoGt()
func TestLinenoCharnoGt(t *testing.T) {
	pt := newParserTest(t)

	pt.testLinenoCharnoBinop(">")
}

// Java: ParserTest.testLinenoCharnoGe()
func TestLinenoCharnoGe(t *testing.T) {
	pt := newParserTest(t)

	pt.testLinenoCharnoBinop(">=")
}

// Java: ParserTest.testLinenoCharnoNegativeNumber()
func TestLinenoCharnoNegativeNumber(t *testing.T) {
	pt := newParserTest(t)

	n := pt.parse("var x = -1000").GetFirstChild().GetFirstFirstChild()

	assertNode(t, n).HasLineno(1)
	assertNode(t, n).HasCharno(8)
	assertNode(t, n).HasLength(5)
}

// Java: ParserTest.testJSDocAttachment1()
func TestJSDocAttachment1(t *testing.T) {
	pt := newParserTest(t)

	varNode := pt.parse("/** @type {number} */var a;").GetFirstChild()

	// VAR
	assertNode(t, varNode).HasType(rhino.Token_VAR)
	pt.assertNodeHasJSDocInfoWithJSType(varNode, pt.NUMBER_TYPE)

	// VAR NAME
	varNameNode := varNode.GetFirstChild()
	assertNode(t, varNameNode).HasType(rhino.Token_NAME)
	assertThat(t, varNameNode.GetJSDocInfo()).IsNull()

	pt.strictMode = Config_StrictMode_SLOPPY
	pt.expectFeatures(parser.FeatureSet_Feature_LET_DECLARATIONS)

	letNode := pt.parse("/** @type {number} */let a;").GetFirstChild()

	// LET
	assertNode(t, letNode).HasType(rhino.Token_LET)
	pt.assertNodeHasJSDocInfoWithJSType(letNode, pt.NUMBER_TYPE)

	// LET NAME
	letNameNode := letNode.GetFirstChild()
	assertNode(t, letNameNode).HasType(rhino.Token_NAME)
	assertThat(t, letNameNode.GetJSDocInfo()).IsNull()

	pt.expectFeatures(parser.FeatureSet_Feature_CONST_DECLARATIONS)
	constNode := pt.parse("/** @type {number} */const a = 0;").GetFirstChild()

	// CONST
	assertNode(t, constNode).HasType(rhino.Token_CONST)
	pt.assertNodeHasJSDocInfoWithJSType(constNode, pt.NUMBER_TYPE)

	// CONST NAME
	constNameNode := constNode.GetFirstChild()
	assertNode(t, constNameNode).HasType(rhino.Token_NAME)
	assertThat(t, constNameNode.GetJSDocInfo()).IsNull()
}

// Java: ParserTest.testJSDocAttachment2()
func TestJSDocAttachment2(t *testing.T) {
	pt := newParserTest(t)

	varNode := pt.parse("/** @type {number} */var a,b;").GetFirstChild()

	// VAR
	assertNode(t, varNode).HasType(rhino.Token_VAR)
	pt.assertNodeHasJSDocInfoWithJSType(varNode, pt.NUMBER_TYPE)

	// First NAME
	nameNode1 := varNode.GetFirstChild()
	assertNode(t, nameNode1).HasType(rhino.Token_NAME)
	assertThat(t, nameNode1.GetJSDocInfo()).IsNull()

	// Second NAME
	nameNode2 := nameNode1.GetNext()
	assertNode(t, nameNode2).HasType(rhino.Token_NAME)
	assertThat(t, nameNode2.GetJSDocInfo()).IsNull()
}

// Java: ParserTest.testJSDocAttachment3()
func TestJSDocAttachment3(t *testing.T) {
	pt := newParserTest(t)

	assignNode := pt.parse("/** @type {number} */goog.FOO = 5;").GetFirstFirstChild()
	assertNode(t, assignNode).HasType(rhino.Token_ASSIGN)
	pt.assertNodeHasJSDocInfoWithJSType(assignNode, pt.NUMBER_TYPE)
}

// Java: ParserTest.testJSDocAttachment4()
func TestJSDocAttachment4(t *testing.T) {
	pt := newParserTest(t)

	varNode := pt.parse("var a, /** @define {number} */ b = 5;").GetFirstChild()

	// ASSIGN
	assertNode(t, varNode).HasType(rhino.Token_VAR)
	assertThat(t, varNode.GetJSDocInfo()).IsNull()

	// a
	a := varNode.GetFirstChild()
	assertThat(t, a.GetJSDocInfo()).IsNull()

	// b
	b := a.GetNext()
	info := b.GetJSDocInfo()
	assertThat(t, info).IsNotNull()
	assertThat(t, info.IsDefine()).IsTrue()
	pt.AssertTypeEqualsExpr(pt.NUMBER_TYPE, info.GetType())
}

// Java: ParserTest.testJSDocAttachment5()
func TestJSDocAttachment5(t *testing.T) {
	pt := newParserTest(t)

	varNode :=
		pt.parse("var /** @type {number} */a, /** @define {number} */b = 5;").GetFirstChild()

	// ASSIGN
	assertNode(t, varNode).HasType(rhino.Token_VAR)
	assertThat(t, varNode.GetJSDocInfo()).IsNull()

	// a
	a := varNode.GetFirstChild()
	assertThat(t, a.GetJSDocInfo()).IsNotNull()
	info := a.GetJSDocInfo()
	assertThat(t, info).IsNotNull()
	assertThat(t, info.IsDefine()).IsFalse()
	pt.AssertTypeEqualsExpr(pt.NUMBER_TYPE, info.GetType())

	// b
	b := a.GetNext()
	info = b.GetJSDocInfo()
	assertThat(t, info).IsNotNull()
	assertThat(t, info.IsDefine()).IsTrue()
	pt.AssertTypeEqualsExpr(pt.NUMBER_TYPE, info.GetType())
}

// Java: ParserTest.testJSDocAttachment6()
func TestJSDocAttachment6(t *testing.T) {
	pt := newParserTest(t)

	functionNode :=
		pt.parse(
			`var a = /** @param {number} index */5;
/** @return {boolean} */function f(index){}
`).
			GetSecondChild()

	assertNode(t, functionNode).HasType(rhino.Token_FUNCTION)
	info := functionNode.GetJSDocInfo()
	assertThat(t, info).IsNotNull()
	assertThat(t, info.HasParameter("index")).IsFalse()
	assertThat(t, info.HasReturnType()).IsTrue()
	pt.AssertTypeEqualsExpr(pt.BOOLEAN_TYPE, info.GetReturnType())
}

// Java: ParserTest.testJSDocAttachment7()
func TestJSDocAttachment7(t *testing.T) {
	pt := newParserTest(t)

	varNode := pt.parse("/** */var a;").GetFirstChild()

	// VAR
	assertNode(t, varNode).HasType(rhino.Token_VAR)

	// NAME
	nameNode := varNode.GetFirstChild()
	assertNode(t, nameNode).HasType(rhino.Token_NAME)
	assertThat(t, nameNode.GetJSDocInfo()).IsNull()
}

// Java: ParserTest.testJSDocAttachment8()
func TestJSDocAttachment8(t *testing.T) {
	pt := newParserTest(t)

	varNode := pt.parse("/** x */var a;").GetFirstChild()

	// VAR
	assertNode(t, varNode).HasType(rhino.Token_VAR)

	// NAME
	nameNode := varNode.GetFirstChild()
	assertNode(t, nameNode).HasType(rhino.Token_NAME)
	assertThat(t, nameNode.GetJSDocInfo()).IsNull()
}

// Java: ParserTest.testJSDocAttachment9()
func TestJSDocAttachment9(t *testing.T) {
	pt := newParserTest(t)

	varNode := pt.parse("/** \n x */var a;").GetFirstChild()

	// VAR
	assertNode(t, varNode).HasType(rhino.Token_VAR)

	// NAME
	nameNode := varNode.GetFirstChild()
	assertNode(t, nameNode).HasType(rhino.Token_NAME)
	assertThat(t, nameNode.GetJSDocInfo()).IsNull()
}

// Java: ParserTest.testJSDocAttachment10()
func TestJSDocAttachment10(t *testing.T) {
	pt := newParserTest(t)

	varNode := pt.parse("/** x\n */var a;").GetFirstChild()

	// VAR
	assertNode(t, varNode).HasType(rhino.Token_VAR)

	// NAME
	nameNode := varNode.GetFirstChild()
	assertNode(t, nameNode).HasType(rhino.Token_NAME)
	assertThat(t, nameNode.GetJSDocInfo()).IsNull()
}

// Java: ParserTest.testJSDocAttachment11()
func TestJSDocAttachment11(t *testing.T) {
	pt := newParserTest(t)

	varNode := pt.parse("/** @type {{x : number, 'y' : string, z}} */var a;").GetFirstChild()

	// VAR
	assertNode(t, varNode).HasType(rhino.Token_VAR)
	pt.assertNodeHasJSDocInfoWithJSType(
		varNode,
		pt.CreateRecordTypeBuilder().AddProperty("x", pt.NUMBER_TYPE, nil).AddProperty("y", pt.STRING_TYPE, nil).AddProperty("z", pt.UNKNOWN_TYPE, nil).Build())

	// NAME
	nameNode := varNode.GetFirstChild()
	assertNode(t, nameNode).HasType(rhino.Token_NAME)
	assertThat(t, nameNode.GetJSDocInfo()).IsNull()
}

// Java: ParserTest.testJSDocAttachment12()
func TestJSDocAttachment12(t *testing.T) {
	pt := newParserTest(t)

	varNode := pt.parse("var a = {/** @type {Object} */ b: c};").GetFirstChild()
	objectLitNode := varNode.GetFirstFirstChild()
	assertNode(t, objectLitNode).HasType(rhino.Token_OBJECTLIT)
	assertThat(t, objectLitNode.GetFirstChild().GetJSDocInfo()).IsNotNull()
}

// Java: ParserTest.testJSDocAttachment13()
func TestJSDocAttachment13(t *testing.T) {
	pt := newParserTest(t)

	varNode := pt.parse("/** foo */ var a;").GetFirstChild()
	assertThat(t, varNode.GetJSDocInfo()).IsNotNull()
}

// Java: ParserTest.testJSDocAttachment14()
func TestJSDocAttachment14(t *testing.T) {
	pt := newParserTest(t)

	varNode := pt.parse("/** */ var a;").GetFirstChild()
	assertThat(t, varNode.GetJSDocInfo()).IsNull()
}

// Java: ParserTest.testJSDocAttachment15()
func TestJSDocAttachment15(t *testing.T) {
	pt := newParserTest(t)

	varNode := pt.parse("/** \n * \n */ var a;").GetFirstChild()
	assertThat(t, varNode.GetJSDocInfo()).IsNull()
}

// Java: ParserTest.testJSDocAttachment16()
func TestJSDocAttachment16(t *testing.T) {
	pt := newParserTest(t)

	exprCall := pt.parse("/** @private */ x(); function f() {};").GetFirstChild()
	assertNode(t, exprCall).HasType(rhino.Token_EXPR_RESULT)
	assertThat(t, exprCall.GetNext().GetJSDocInfo()).IsNull()
	assertThat(t, exprCall.GetFirstChild().GetJSDocInfo()).IsNotNull()
}

// Java: ParserTest.testJSDocAttachmentForCastFnCall()
func TestJSDocAttachmentForCastFnCall(t *testing.T) {
	pt := newParserTest(t)

	fn :=
		pt.parse(
			`function f() {
  return /** @type {string} */ (g(1 /** @desc x */));
};
`).
			GetFirstChild()
	assertNode(t, fn).HasType(rhino.Token_FUNCTION)
	cast := fn.GetLastChild().GetFirstFirstChild()
	assertNode(t, cast).HasType(rhino.Token_CAST)
}

// Java: ParserTest.testJSDocAttachmentForCastName()
func TestJSDocAttachmentForCastName(t *testing.T) {
	pt := newParserTest(t)

	fn :=
		pt.parse(
			`function f() {
  var x = /** @type {string} */ (y);
};
`).
			GetFirstChild()
	assertNode(t, fn).HasType(rhino.Token_FUNCTION)
	cast := fn.GetLastChild().GetFirstFirstChild().GetFirstChild()
	assertNode(t, cast).HasType(rhino.Token_CAST)
}

// Java: ParserTest.testJSDocAttachmentForCastLhs()
func TestJSDocAttachmentForCastLhs(t *testing.T) {
	pt := newParserTest(t)

	expr := pt.parse("/** some jsdoc */ (/** @type {?} */ (a)).b = 0;").GetOnlyChild()
	lhs := expr.GetFirstFirstChild() // child is ASSIGN, grandchild is the GETPROP
	cast := lhs.GetFirstChild()
	assertNode(t, cast).HasToken(rhino.Token_CAST)
	assertThat(t, cast.GetJSDocInfo()).IsNotNull()
	// TODO(b/123955687): this should be true
	assertThat(t, cast.GetJSDocInfo().HasType()).IsFalse()
}

// Java: ParserTest.testJSDocAttachment19()
func TestJSDocAttachment19(t *testing.T) {
	pt := newParserTest(t)

	fn :=
		pt.parse(
			`function f() {
  /** @type {string} */
  return;
};
`).
			GetFirstChild()
	assertNode(t, fn).HasType(rhino.Token_FUNCTION)

	ret := fn.GetLastChild().GetFirstChild()
	assertNode(t, ret).HasType(rhino.Token_RETURN)
	assertThat(t, ret.GetJSDocInfo()).IsNotNull()
}

// Java: ParserTest.testJSDocAttachment20()
func TestJSDocAttachment20(t *testing.T) {
	pt := newParserTest(t)

	fn :=
		pt.parse(
			`function f() {
  /** @type {string} */
  if (true) return;
};
`).
			GetFirstChild()
	assertNode(t, fn).HasType(rhino.Token_FUNCTION)

	ret := fn.GetLastChild().GetFirstChild()
	assertNode(t, ret).HasType(rhino.Token_IF)
	assertThat(t, ret.GetJSDocInfo()).IsNotNull()
}

// Java: ParserTest.testJSDocAttachment21()
func TestJSDocAttachment21(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY
	pt.expectFeatures(parser.FeatureSet_Feature_CONST_DECLARATIONS)
	pt.parse("/** @param {string} x */ const f = function() {};")

	pt.expectFeatures(parser.FeatureSet_Feature_LET_DECLARATIONS)
	pt.parse("/** @param {string} x */ let f = function() {};")
}

// Java: ParserTest.testJSDocAttachment22()
func TestJSDocAttachment22(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY
	pt.expectFeatures(parser.FeatureSet_Feature_MODULES)

	n := pt.parse("/** @param {string} x */ export function f(x) {};")
	export := n.GetFirstFirstChild()

	assertNode(t, export).HasType(rhino.Token_EXPORT)
	assertThat(t, export.GetJSDocInfo()).IsNull()
	assertThat(t, export.GetFirstChild().GetJSDocInfo()).IsNotNull()
	assertThat(t, export.GetFirstChild().GetJSDocInfo().HasParameter("x")).IsTrue()
}

// Java: ParserTest.testInlineJSDocAttachmentToVar()
func TestInlineJSDocAttachmentToVar(t *testing.T) {
	pt := newParserTest(t)

	letNode := pt.parse("let /** string */ x = 'a';").GetFirstChild()
	assertNode(t, letNode).HasType(rhino.Token_LET)

	info := letNode.GetFirstChild().GetJSDocInfo()
	assertThat(t, info).IsNotNull()
	pt.AssertTypeEqualsExpr(pt.STRING_TYPE, info.GetType())
}

// Java: ParserTest.testNodeInsideParens_simpleString()
func TestNodeInsideParens_simpleString(t *testing.T) {
	pt := newParserTest(t)

	exprRes := pt.parse("('a')").GetFirstChild()
	str := exprRes.GetFirstChild()

	assertNode(t, exprRes).HasType(rhino.Token_EXPR_RESULT)
	assertNode(t, str).HasType(rhino.Token_STRINGLIT)
	assertThat(t, str.GetIsParenthesized()).IsTrue()
}

// Java: ParserTest.testNodesInsideParens()
func TestNodesInsideParens(t *testing.T) {
	pt := newParserTest(t)

	andExprRes := pt.parse("((x&&y) && z)").GetFirstChild()
	assertNode(t, andExprRes).HasType(rhino.Token_EXPR_RESULT)

	//  The expression `((x&&y) && z)` is marked as parenthesized
	andNode := andExprRes.GetFirstChild()
	assertNode(t, andNode).HasType(rhino.Token_AND)
	assertThat(t, andNode.GetIsParenthesized()).IsTrue()

	// The expression `(x&&y)` is marked as parenthesized
	// These parens are still recorded despite not affecting the structure of the parent expression.
	xAndy := andNode.GetFirstChild()
	assertNode(t, xAndy).HasType(rhino.Token_AND)
	assertThat(t, xAndy.GetIsParenthesized()).IsTrue()

	// inner nodes `x` and `y` of a `(x&&y)` aren't marked as parenthesized
	x := xAndy.GetFirstChild()
	y := xAndy.GetSecondChild()
	assertNode(t, x).HasType(rhino.Token_NAME)
	assertNode(t, y).HasType(rhino.Token_NAME)
	assertThat(t, x.GetIsParenthesized()).IsFalse()
	assertThat(t, y.GetIsParenthesized()).IsFalse()

	z := andNode.GetSecondChild()
	assertNode(t, z).HasType(rhino.Token_NAME)
	assertThat(t, z.GetIsParenthesized()).IsFalse()
}

// Java: ParserTest.testInlineNonJSDocCommentAttachmentToVar()
func TestInlineNonJSDocCommentAttachmentToVar(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	letNode := pt.parse("let /* blah */ x = 'a';").GetFirstChild()
	assertNode(t, letNode).HasType(rhino.Token_LET)

	nonJSDocComment := letNode.GetFirstChild().GetNonJSDocComment()
	assertThat(t, nonJSDocComment).IsNotNull()
	assertThat(t, nonJSDocComment.IsInline()).IsTrue()
	assertThat(t, nonJSDocComment.IsEndingAsLineComment()).IsFalse()
	assertThat(t, nonJSDocComment.GetCommentString()).Contains("/* blah */")
}

// Java: ParserTest.testNoInlineNonJSDocCommentAttachmentWithoutParsingMode()
func TestNoInlineNonJSDocCommentAttachmentWithoutParsingMode(t *testing.T) {
	pt := newParserTest(t)

	letNode := pt.parse("let /* blah */ x = 'a';").GetFirstChild()
	assertNode(t, letNode).HasType(rhino.Token_LET)

	nonJSDocComment := letNode.GetFirstChild().GetNonJSDocComment()
	assertThat(t, nonJSDocComment).IsNull()
}

// Java: ParserTest.testInlineJSDocAttachmentToObjPatNormalProp()
func TestInlineJSDocAttachmentToObjPatNormalProp(t *testing.T) {
	pt := newParserTest(t)

	letNode :=
		pt.parse("let { normalProp: /** string */ normalPropTarget } = {};").GetFirstChild()
	assertNode(t, letNode).HasType(rhino.Token_LET)

	destructuringLhs := letNode.GetFirstChild()
	objectPattern := destructuringLhs.GetFirstChild()

	normalProp := objectPattern.GetFirstChild()
	assertNode(t, normalProp).HasType(rhino.Token_STRING_KEY)
	normalPropTarget := normalProp.GetOnlyChild()
	pt.assertNodeHasJSDocInfoWithJSType(normalPropTarget, pt.STRING_TYPE)
}

// Java: ParserTest.testInlineNonJSDocCommentAttachmentToObjPatNormalProp()
func TestInlineNonJSDocCommentAttachmentToObjPatNormalProp(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	letNode := pt.parse("let { normalProp: /* blah */ normalPropTarget } = {};").GetFirstChild()
	assertNode(t, letNode).HasType(rhino.Token_LET)

	destructuringLhs := letNode.GetFirstChild()
	objectPattern := destructuringLhs.GetFirstChild()

	normalProp := objectPattern.GetFirstChild()
	assertNode(t, normalProp).HasType(rhino.Token_STRING_KEY)
	normalPropTarget := normalProp.GetOnlyChild()
	assertThat(t, normalPropTarget.GetNonJSDocCommentString()).Contains("/* blah */")
}

// Java: ParserTest.testNoInlineNonJSDocCommentAttachmentToObjWithoutParsingMode()
func TestNoInlineNonJSDocCommentAttachmentToObjWithoutParsingMode(t *testing.T) {
	pt := newParserTest(t)

	letNode := pt.parse("let { normalProp: /* blah */ normalPropTarget } = {};").GetFirstChild()
	assertNode(t, letNode).HasType(rhino.Token_LET)

	destructuringLhs := letNode.GetFirstChild()
	objectPattern := destructuringLhs.GetFirstChild()

	normalProp := objectPattern.GetFirstChild()
	assertNode(t, normalProp).HasType(rhino.Token_STRING_KEY)
	normalPropTarget := normalProp.GetOnlyChild()
	assertThat(t, normalPropTarget.GetNonJSDocComment()).IsNull()
}

// Java: ParserTest.testInlineJSDocAttachmentToObjPatNormalPropKey()
func TestInlineJSDocAttachmentToObjPatNormalPropKey(t *testing.T) {
	pt := newParserTest(t)

	letNode := pt.parse("let { /** string */ normalProp: normalProp } = {};").GetFirstChild()
	assertNode(t, letNode).HasType(rhino.Token_LET)

	destructuringLhs := letNode.GetFirstChild()
	objectPattern := destructuringLhs.GetFirstChild()

	normalProp := objectPattern.GetFirstChild()
	assertNode(t, normalProp).HasType(rhino.Token_STRING_KEY)
	// TODO(bradfordcsmith): Putting the inline jsdoc on the key should be an error,
	//     because it isn't clear what that should mean.
	pt.assertNodeHasNoJSDocInfo(normalProp)
}

// Java: ParserTest.testInlineNonJSDocCommentAttachmentToObjPatNormalPropKey()
func TestInlineNonJSDocCommentAttachmentToObjPatNormalPropKey(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	letNode := pt.parse("let { /* blah */ normalProp: normalProp } = {};").GetFirstChild()
	assertNode(t, letNode).HasType(rhino.Token_LET)

	destructuringLhs := letNode.GetFirstChild()
	objectPattern := destructuringLhs.GetFirstChild()

	normalPropKey := objectPattern.GetFirstChild()
	assertNode(t, normalPropKey).HasType(rhino.Token_STRING_KEY)
	assertThat(t, normalPropKey.GetNonJSDocCommentString()).Contains("/* blah */")
}

// Java: ParserTest.testInlineNonJSDocCommentAttachment_numberAsKey()
func TestInlineNonJSDocCommentAttachment_numberAsKey(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	script := pt.parse("var { /*comment */ 3: x} = foo();")
	var_ := script.GetFirstChild()
	destructuringLHS := var_.GetFirstChild()
	objLit := destructuringLHS.GetFirstChild()
	three := objLit.GetFirstChild()
	x := three.GetFirstChild()
	assertThat(t, three.GetNonJSDocComment()).IsNotNull()
	assertThat(t, three.GetNonJSDocComment().GetCommentString()).IsEqualTo("/*comment */")
	assertThat(t, three.GetNonJSDocComment().IsInline()).IsTrue()
	assertThat(t, x.GetNonJSDocComment()).IsNull()
}

// Java: ParserTest.testInlineNonJSDocCommentAttachment_quotedKey()
func TestInlineNonJSDocCommentAttachment_quotedKey(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	script := pt.parse("var { /*comment */ 'a': x} = foo();")
	var_ := script.GetFirstChild()
	destructuringLHS := var_.GetFirstChild()
	objLit := destructuringLHS.GetFirstChild()
	a := objLit.GetFirstChild()
	x := a.GetFirstChild()
	assertThat(t, a.GetNonJSDocComment()).IsNotNull()
	assertThat(t, a.GetNonJSDocComment().GetCommentString()).IsEqualTo("/*comment */")
	assertThat(t, a.GetNonJSDocComment().IsInline()).IsTrue()
	assertThat(t, x.GetNonJSDocComment()).IsNull()
}

// Java: ParserTest.testInlineJSDocAttachmentToObjPatShorthandProp()
func TestInlineJSDocAttachmentToObjPatShorthandProp(t *testing.T) {
	pt := newParserTest(t)

	letNode := pt.parse("let { /** string */ shorthandProp } = {};").GetFirstChild()
	assertNode(t, letNode).HasType(rhino.Token_LET)

	destructuringLhs := letNode.GetFirstChild()
	objectPattern := destructuringLhs.GetFirstChild()

	shorthandProp := objectPattern.GetFirstChild()
	assertNode(t, shorthandProp).HasType(rhino.Token_STRING_KEY)
	shorthandPropTarget := shorthandProp.GetOnlyChild()
	pt.assertNodeHasJSDocInfoWithJSType(shorthandPropTarget, pt.STRING_TYPE)
}

// Java: ParserTest.testInlineNonJSDocAttachmentToObjPatShorthandProp()
func TestInlineNonJSDocAttachmentToObjPatShorthandProp(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	letNode := pt.parse("let { /* blah */ shorthandProp } = {};").GetFirstChild()
	assertNode(t, letNode).HasType(rhino.Token_LET)

	destructuringLhs := letNode.GetFirstChild()
	objectPattern := destructuringLhs.GetFirstChild()

	shorthandProp := objectPattern.GetFirstChild()
	assertNode(t, shorthandProp).HasType(rhino.Token_STRING_KEY)
	assertThat(t, shorthandProp.GetNonJSDocCommentString()).Contains("/* blah */")
}

// Java: ParserTest.testInlineJSDocAttachmentToObjPatNormalPropWithDefault()
func TestInlineJSDocAttachmentToObjPatNormalPropWithDefault(t *testing.T) {
	pt := newParserTest(t)

	letNode :=
		pt.parse("let { normalPropWithDefault: /** string */ normalPropWithDefault = 'hi' } = {};").
			GetFirstChild()
	assertNode(t, letNode).HasType(rhino.Token_LET)

	destructuringLhs := letNode.GetFirstChild()
	objectPattern := destructuringLhs.GetFirstChild()

	normalPropWithDefault := objectPattern.GetFirstChild()
	assertNode(t, normalPropWithDefault).HasType(rhino.Token_STRING_KEY)
	normalPropDefaultValue := normalPropWithDefault.GetOnlyChild()
	assertNode(t, normalPropDefaultValue).HasType(rhino.Token_DEFAULT_VALUE)
	normalPropWithDefaultTarget := normalPropDefaultValue.GetFirstChild()
	pt.assertNodeHasJSDocInfoWithJSType(normalPropWithDefaultTarget, pt.STRING_TYPE)
}

// Java: ParserTest.testInlineNonJSDocCommentAttachmentToObjPatNormalPropWithDefault()
func TestInlineNonJSDocCommentAttachmentToObjPatNormalPropWithDefault(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	letNode :=
		pt.parse("let { normalPropWithDefault: /* blah */ normalPropWithDefault = 'hi' } = {};").
			GetFirstChild()
	assertNode(t, letNode).HasType(rhino.Token_LET)

	destructuringLhs := letNode.GetFirstChild()
	objectPattern := destructuringLhs.GetFirstChild()

	normalPropWithDefault := objectPattern.GetFirstChild()
	assertNode(t, normalPropWithDefault).HasType(rhino.Token_STRING_KEY)
	normalPropDefaultValue := normalPropWithDefault.GetOnlyChild()
	assertNode(t, normalPropDefaultValue).HasType(rhino.Token_DEFAULT_VALUE)
	normalPropWithDefaultTarget := normalPropDefaultValue.GetFirstChild()
	assertThat(t, normalPropWithDefaultTarget.GetNonJSDocCommentString()).Contains("/* blah */")
}

// Java: ParserTest.testInlineJSDocAttachmentToObjPatShorthandWithDefault()
func TestInlineJSDocAttachmentToObjPatShorthandWithDefault(t *testing.T) {
	pt := newParserTest(t)

	letNode :=
		pt.parse("let { /** string */ shorthandPropWithDefault = 'lo' } = {};").GetFirstChild()
	assertNode(t, letNode).HasType(rhino.Token_LET)

	destructuringLhs := letNode.GetFirstChild()
	objectPattern := destructuringLhs.GetFirstChild()

	shorthandPropWithDefault := objectPattern.GetFirstChild()
	assertNode(t, shorthandPropWithDefault).HasType(rhino.Token_STRING_KEY)
	shorthandPropDefaultValue := shorthandPropWithDefault.GetOnlyChild()
	assertNode(t, shorthandPropDefaultValue).HasType(rhino.Token_DEFAULT_VALUE)
	shorthandPropWithDefaultTarget := shorthandPropDefaultValue.GetFirstChild()
	pt.assertNodeHasJSDocInfoWithJSType(shorthandPropWithDefaultTarget, pt.STRING_TYPE)
}

// Java: ParserTest.testInlineNonJSDocCommentAttachmentToObjPatShorthandWithDefault()
func TestInlineNonJSDocCommentAttachmentToObjPatShorthandWithDefault(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	letNode :=
		pt.parse("let { /* blah */ shorthandPropWithDefault = 'lo' } = {};").GetFirstChild()
	assertNode(t, letNode).HasType(rhino.Token_LET)

	destructuringLhs := letNode.GetFirstChild()
	objectPattern := destructuringLhs.GetFirstChild()

	shorthandPropWithDefault := objectPattern.GetFirstChild()
	assertNode(t, shorthandPropWithDefault).HasType(rhino.Token_STRING_KEY)
	shorthandPropDefaultValue := shorthandPropWithDefault.GetOnlyChild()
	assertNode(t, shorthandPropDefaultValue).HasType(rhino.Token_DEFAULT_VALUE)
	shorthandPropWithDefaultTarget := shorthandPropDefaultValue.GetFirstChild()
	assertThat(t, shorthandPropWithDefaultTarget.GetNonJSDocCommentString()).Contains("/* blah */")
}

// Java: ParserTest.testObjLitKeyNonJSDocComment()
func TestObjLitKeyNonJSDocComment(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	exprResult := pt.parse(`a = {
// blah
1n: 3};`).GetFirstChild()
	assignNode := exprResult.GetFirstChild()
	objectLit := assignNode.GetSecondChild()
	stringKey := objectLit.GetFirstChild()
	assertThat(t, stringKey.GetNonJSDocCommentString()).Contains("// blah")
}

// Java: ParserTest.testLabelNonJSDocComment()
func TestLabelNonJSDocComment(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	label := pt.parse(`
// blah
label:
for(;;){ continue label; }`).GetFirstChild()
	assertNode(t, label).HasType(rhino.Token_LABEL)
	assertThat(t, label.GetNonJSDocCommentString()).Contains("// blah")
}

// Java: ParserTest.testFieldNonJSDocComment()
func TestFieldNonJSDocComment(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	field :=
		pt.parse(`class C{
// blah
 field }`).GetFirstChild().GetLastChild().GetFirstChild()
	assertNode(t, field).HasType(rhino.Token_MEMBER_FIELD_DEF)
	assertThat(t, field.GetNonJSDocCommentString()).Contains("// blah")
}

// Java: ParserTest.testPropertyNameAssignmentNonJSDocComment()
func TestPropertyNameAssignmentNonJSDocComment(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	varName :=
		pt.parse("let {key: \n// blah\nvarName} = someObject").
			GetFirstChild().
			GetFirstChild().
			GetFirstChild().
			GetFirstChild().
			GetFirstChild()
	assertNode(t, varName).HasType(rhino.Token_NAME)
	assertThat(t, varName.GetNonJSDocCommentString()).Contains("// blah")
}

// Java: ParserTest.testGetPropCallNonJSDocComment()
func TestGetPropCallNonJSDocComment(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	exprResult := pt.parse(`foo.
// blah
baz(3);`).GetFirstChild()
	call := exprResult.GetFirstChild()
	assertNode(t, call).HasType(rhino.Token_CALL)

	bazAccess := call.GetFirstChild()
	assertThat(t, bazAccess.GetNonJSDocCommentString()).Contains("// blah")
}

// Java: ParserTest.testGetPropOptionalCallNonJSDocComment()
func TestGetPropOptionalCallNonJSDocComment(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	exprResult := pt.parse(`foo?.
// blah
baz(3);`).GetFirstChild()
	callNode := exprResult.GetFirstChild()
	assertNode(t, callNode).HasType(rhino.Token_OPTCHAIN_CALL)

	bazAccess := callNode.GetFirstChild()
	assertThat(t, bazAccess.GetNonJSDocCommentString()).Contains("// blah")
}

// Java: ParserTest.testGetPropAssignmentNonJSDocComment()
func TestGetPropAssignmentNonJSDocComment(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	exprResult := pt.parse(`foo.
// blah
baz = 5;`).GetFirstChild()
	assign := exprResult.GetFirstChild()
	assertNode(t, assign).HasType(rhino.Token_ASSIGN)

	bazAccess := assign.GetFirstChild()
	assertThat(t, bazAccess.GetNonJSDocCommentString()).Contains("// blah")
}

// Java: ParserTest.testNonJSDocCommentsOnAdjacentNodes()
func TestNonJSDocCommentsOnAdjacentNodes(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	root :=
		pt.parse(
			`// comment before GETPROP
a(
// comment on GETPROP
)
.b();
// comment after GETPROP
c();
`)
	exprResultA := root.GetFirstChild()
	assertNode(t, exprResultA).HasType(rhino.Token_EXPR_RESULT)
	assertThat(t, exprResultA.GetNonJSDocCommentString()).IsEqualTo("// comment before GETPROP")

	nodeB := exprResultA.GetFirstFirstChild()
	assertNode(t, nodeB).HasType(rhino.Token_GETPROP)
	assertThat(t, nodeB.GetNonJSDocCommentString()).IsEqualTo("// comment on GETPROP")

	exprResultC := root.GetLastChild()
	assertNode(t, exprResultC).HasType(rhino.Token_EXPR_RESULT)
	assertThat(t, exprResultC.GetNonJSDocCommentString()).IsEqualTo("// comment after GETPROP")
}

// Java: ParserTest.testInlineJSDocAttachmentToObjPatComputedPropKey()
func TestInlineJSDocAttachmentToObjPatComputedPropKey(t *testing.T) {
	pt := newParserTest(t)

	letNode :=
		pt.parse("let { /** string */ ['computedProp']: computedProp } = {};").GetFirstChild()
	assertNode(t, letNode).HasType(rhino.Token_LET)

	destructuringLhs := letNode.GetFirstChild()
	objectPattern := destructuringLhs.GetFirstChild()

	computedProp := objectPattern.GetFirstChild()
	assertNode(t, computedProp).HasType(rhino.Token_COMPUTED_PROP)
	// TODO(bradfordcsmith): Putting inline JSDoc on the computed property key should be an error,
	//     since it's not clear what it should mean.
	pt.assertNodeHasNoJSDocInfo(computedProp)
}

// Java: ParserTest.testInlineJSDocAttachmentToObjPatComputedProp()
func TestInlineJSDocAttachmentToObjPatComputedProp(t *testing.T) {
	pt := newParserTest(t)

	letNode :=
		pt.parse("let { ['computedProp']: /** string */ computedProp } = {};").GetFirstChild()
	assertNode(t, letNode).HasType(rhino.Token_LET)

	destructuringLhs := letNode.GetFirstChild()
	objectPattern := destructuringLhs.GetFirstChild()

	computedProp := objectPattern.GetFirstChild()
	assertNode(t, computedProp).HasType(rhino.Token_COMPUTED_PROP)
	computedPropTarget := computedProp.GetSecondChild()
	pt.assertNodeHasJSDocInfoWithJSType(computedPropTarget, pt.STRING_TYPE)
}

// Java: ParserTest.testInlineJSDocAttachmentToObjPatComputedPropWithDefault()
func TestInlineJSDocAttachmentToObjPatComputedPropWithDefault(t *testing.T) {
	pt := newParserTest(t)

	letNode :=
		pt.parse("let { ['computedPropWithDefault']: /** string */ computedProp = 'go' } = {};").
			GetFirstChild()
	assertNode(t, letNode).HasType(rhino.Token_LET)

	destructuringLhs := letNode.GetFirstChild()
	objectPattern := destructuringLhs.GetFirstChild()

	computedPropWithDefault := objectPattern.GetFirstChild()
	assertNode(t, computedPropWithDefault).HasType(rhino.Token_COMPUTED_PROP)
	computedPropDefaultValue := computedPropWithDefault.GetSecondChild()
	assertNode(t, computedPropDefaultValue).HasType(rhino.Token_DEFAULT_VALUE)
	computedPropWithDefaultTarget := computedPropDefaultValue.GetFirstChild()
	pt.assertNodeHasJSDocInfoWithJSType(computedPropWithDefaultTarget, pt.STRING_TYPE)
}

// Java: ParserTest.testInlineJSDocAttachmentToObjPatNormalPropWithQualifiedName()
func TestInlineJSDocAttachmentToObjPatNormalPropWithQualifiedName(t *testing.T) {
	pt := newParserTest(t)

	exprResult :=
		pt.parse("({ normalProp: /** string */ ns.normalPropTarget } = {});").GetFirstChild()
	assignNode := exprResult.GetFirstChild()
	assertNode(t, assignNode).HasType(rhino.Token_ASSIGN)

	objectPattern := assignNode.GetFirstChild()

	normalProp := objectPattern.GetFirstChild()
	assertNode(t, normalProp).HasType(rhino.Token_STRING_KEY)
	nsNormalPropTarget := normalProp.GetOnlyChild()
	pt.assertNodeHasJSDocInfoWithJSType(nsNormalPropTarget, pt.STRING_TYPE)
}

// Java: ParserTest.testInlineJSDocAttachmentToObjPatNormalPropWithQualifiedNameWithDefault()
func TestInlineJSDocAttachmentToObjPatNormalPropWithQualifiedNameWithDefault(t *testing.T) {
	pt := newParserTest(t)

	exprResult :=
		pt.parse("({ normalProp: /** string */ ns.normalPropTarget = 'foo' } = {});").GetFirstChild()
	assignNode := exprResult.GetFirstChild()
	assertNode(t, assignNode).HasType(rhino.Token_ASSIGN)

	objectPattern := assignNode.GetFirstChild()

	normalProp := objectPattern.GetFirstChild()
	assertNode(t, normalProp).HasType(rhino.Token_STRING_KEY)
	defaultValue := normalProp.GetFirstChild()
	assertNode(t, defaultValue).HasType(rhino.Token_DEFAULT_VALUE)
	nsNormalPropTarget := defaultValue.GetFirstChild()
	pt.assertNodeHasJSDocInfoWithJSType(nsNormalPropTarget, pt.STRING_TYPE)
}

// Java: ParserTest.testInlineJSDocAttachmentToArrayPatElement()
func TestInlineJSDocAttachmentToArrayPatElement(t *testing.T) {
	pt := newParserTest(t)

	letNode := pt.parse("let [/** string */ x] = [];").GetFirstChild()
	assertNode(t, letNode).HasType(rhino.Token_LET)

	destructuringLhs := letNode.GetFirstChild()
	arrayPattern := destructuringLhs.GetFirstChild()
	xVarName := arrayPattern.GetFirstChild()
	pt.assertNodeHasJSDocInfoWithJSType(xVarName, pt.STRING_TYPE)
}

// Java: ParserTest.testInlineJSDocAttachmentToArrayPatElementWithDefault()
func TestInlineJSDocAttachmentToArrayPatElementWithDefault(t *testing.T) {
	pt := newParserTest(t)

	letNode := pt.parse("let [/** string */ x = 'hi'] = [];").GetFirstChild()
	assertNode(t, letNode).HasType(rhino.Token_LET)

	destructuringLhs := letNode.GetFirstChild()
	arrayPattern := destructuringLhs.GetFirstChild()
	defaultValue := arrayPattern.GetFirstChild()
	assertNode(t, defaultValue).HasType(rhino.Token_DEFAULT_VALUE)

	xVarName := defaultValue.GetFirstChild()
	pt.assertNodeHasJSDocInfoWithJSType(xVarName, pt.STRING_TYPE)
}

// Java: ParserTest.testInlineJSDocAttachmentToArrayPatElementQualifiedName()
func TestInlineJSDocAttachmentToArrayPatElementQualifiedName(t *testing.T) {
	pt := newParserTest(t)

	exprResult := pt.parse("[/** string */ x.y.z] = [];").GetFirstChild()
	assignNode := exprResult.GetFirstChild()
	assertNode(t, assignNode).HasType(rhino.Token_ASSIGN)

	arrayPattern := assignNode.GetFirstChild()
	assertNode(t, arrayPattern).HasType(rhino.Token_ARRAY_PATTERN)
	xYZName := arrayPattern.GetFirstChild()
	pt.assertNodeHasJSDocInfoWithJSType(xYZName, pt.STRING_TYPE)
}

// Java: ParserTest.testInlineJSDocAttachmentToArrayPatElementQualifiedNameWithDefault()
func TestInlineJSDocAttachmentToArrayPatElementQualifiedNameWithDefault(t *testing.T) {
	pt := newParserTest(t)

	exprResult := pt.parse("[/** string */ x.y.z = 'foo'] = [];").GetFirstChild()
	assignNode := exprResult.GetFirstChild()
	assertNode(t, assignNode).HasType(rhino.Token_ASSIGN)

	arrayPattern := assignNode.GetFirstChild()
	assertNode(t, arrayPattern).HasType(rhino.Token_ARRAY_PATTERN)
	defaultValue := arrayPattern.GetOnlyChild()
	assertNode(t, defaultValue).HasType(rhino.Token_DEFAULT_VALUE)
	xYZName := defaultValue.GetFirstChild()
	pt.assertNodeHasJSDocInfoWithJSType(xYZName, pt.STRING_TYPE)
}

// Java: ParserTest.testInlineJSDocAttachmentToArrayPatElementAfterElision()
func TestInlineJSDocAttachmentToArrayPatElementAfterElision(t *testing.T) {
	pt := newParserTest(t)

	letNode := pt.parse("let [, /** string */ x] = [];").GetFirstChild()
	assertNode(t, letNode).HasType(rhino.Token_LET)

	destructuringLhs := letNode.GetFirstChild()
	arrayPattern := destructuringLhs.GetFirstChild()
	empty := arrayPattern.GetFirstChild()
	assertNode(t, empty).HasToken(rhino.Token_EMPTY)
	assertNode(t, empty).HasCharno(5)
	assertNode(t, empty).HasLength(1)
	xVarName := arrayPattern.GetSecondChild()
	pt.assertNodeHasJSDocInfoWithJSType(xVarName, pt.STRING_TYPE)
}

// Java: ParserTest.testInlineJSDocAttachmentToObjLitNormalProp()
func TestInlineJSDocAttachmentToObjLitNormalProp(t *testing.T) {
	pt := newParserTest(t)

	letNode := pt.parse("let x = { normalProp: /** string */ normalPropTarget };").GetFirstChild()
	assertNode(t, letNode).HasType(rhino.Token_LET)

	xNode := letNode.GetFirstChild()
	objectLit := xNode.GetFirstChild()

	normalProp := objectLit.GetFirstChild()
	assertNode(t, normalProp).HasType(rhino.Token_STRING_KEY)
	normalPropTarget := normalProp.GetOnlyChild()
	// TODO(bradfordcsmith): Make sure CheckJsDoc considers this an error, because it doesn't
	//     make sense to have inline JSDoc on the value.
	pt.assertNodeHasJSDocInfoWithNoJSType(normalPropTarget)
}

// Java: ParserTest.testInlineJSDocAttachmentToObjLitNormalPropKey()
func TestInlineJSDocAttachmentToObjLitNormalPropKey(t *testing.T) {
	pt := newParserTest(t)

	letNode := pt.parse("let x = { /** string */ normalProp: normalProp };").GetFirstChild()
	assertNode(t, letNode).HasType(rhino.Token_LET)

	xNode := letNode.GetFirstChild()
	objectLit := xNode.GetFirstChild()

	normalProp := objectLit.GetFirstChild()
	assertNode(t, normalProp).HasType(rhino.Token_STRING_KEY)
	// TODO(bradfordcsmith): We should either disallow inline JSDoc here or correctly pull the type
	//     out of it.
	pt.assertNodeHasJSDocInfoWithNoJSType(normalProp)
}

// Java: ParserTest.testJSDocAttachmentToObjLitNormalPropKey()
func TestJSDocAttachmentToObjLitNormalPropKey(t *testing.T) {
	pt := newParserTest(t)

	letNode :=
		pt.parse("let x = { /** @type {string} */ normalProp: normalProp };").GetFirstChild()
	assertNode(t, letNode).HasType(rhino.Token_LET)

	xNode := letNode.GetFirstChild()
	objectLit := xNode.GetFirstChild()

	normalProp := objectLit.GetFirstChild()
	assertNode(t, normalProp).HasType(rhino.Token_STRING_KEY)
	pt.assertNodeHasJSDocInfoWithJSType(normalProp, pt.STRING_TYPE)
}

// Java: ParserTest.testInlineJSDocAttachmentToObjLitShorthandProp()
func TestInlineJSDocAttachmentToObjLitShorthandProp(t *testing.T) {
	pt := newParserTest(t)

	letNode := pt.parse("let x = { /** string */ shorthandProp };").GetFirstChild()
	assertNode(t, letNode).HasType(rhino.Token_LET)

	xNode := letNode.GetFirstChild()
	objectLit := xNode.GetFirstChild()

	shorthandPropKey := objectLit.GetFirstChild()
	assertNode(t, shorthandPropKey).HasType(rhino.Token_STRING_KEY)
	// TODO(bradfordcsmith): We should either disallow inline JSDoc here or correctly pull the type
	//     out of it.
	pt.assertNodeHasJSDocInfoWithNoJSType(shorthandPropKey)
	shorthandPropTarget := shorthandPropKey.GetOnlyChild()
	pt.assertNodeHasNoJSDocInfo(shorthandPropTarget)
}

// Java: ParserTest.testJSDocAttachmentToObjLitShorthandProp()
func TestJSDocAttachmentToObjLitShorthandProp(t *testing.T) {
	pt := newParserTest(t)

	letNode := pt.parse("let x = { /** @type {string} */ shorthandProp };").GetFirstChild()
	assertNode(t, letNode).HasType(rhino.Token_LET)

	xNode := letNode.GetFirstChild()
	objectLit := xNode.GetFirstChild()

	shorthandPropKey := objectLit.GetFirstChild()
	assertNode(t, shorthandPropKey).HasType(rhino.Token_STRING_KEY)
	pt.assertNodeHasJSDocInfoWithJSType(shorthandPropKey, pt.STRING_TYPE)
	shorthandPropTarget := shorthandPropKey.GetOnlyChild()
	pt.assertNodeHasNoJSDocInfo(shorthandPropTarget)
}

// Java: ParserTest.testInlineJSDocAttachment1()
func TestInlineJSDocAttachment1(t *testing.T) {
	pt := newParserTest(t)

	fn := pt.parse("function f(/** string */ x) {}").GetFirstChild()
	assertNode(t, fn).HasType(rhino.Token_FUNCTION)

	info := fn.GetSecondChild().GetFirstChild().GetJSDocInfo()
	assertThat(t, info).IsNotNull()
	pt.AssertTypeEqualsExpr(pt.STRING_TYPE, info.GetType())
}

// Java: ParserTest.testInline_BlockCommentAttachment()
func TestInline_BlockCommentAttachment(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	fn := pt.parse("function f(/* blah */ x) {}").GetFirstChild()
	assertNode(t, fn).HasType(rhino.Token_FUNCTION)
	xNode := fn.GetSecondChild().GetFirstChild()
	assertThat(t, xNode.GetNonJSDocCommentString()).Contains("/* blah */")
}

// Java: ParserTest.testInline_LineCommentAttachment()
func TestInline_LineCommentAttachment(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	fn := pt.parse("function f( // blah\n x) {}").GetFirstChild()
	assertNode(t, fn).HasType(rhino.Token_FUNCTION)
	xNode := fn.GetSecondChild().GetFirstChild()
	assertThat(t, xNode.GetNonJSDocCommentString()).Contains("// blah")
}

// Java: ParserTest.testInlineNonJSDocComments_TrailingAndNonTrailing_ParamList()
func TestInlineNonJSDocComments_TrailingAndNonTrailing_ParamList(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	fn := pt.parse("function f(x /* first */ , /* second */ y ) {}").GetFirstChild()
	assertNode(t, fn).HasType(rhino.Token_FUNCTION)

	paramListNode := fn.GetSecondChild()
	xNode := paramListNode.GetFirstChild()
	yNode := paramListNode.GetSecondChild()

	assertThat(t, xNode.GetTrailingNonJSDocCommentString()).Contains("/* first */")
	assertThat(t, yNode.GetNonJSDocCommentString()).IsEqualTo("/* second */")
}

// Java: ParserTest.testInlineNonJSDocTrailingComments_ParamList()
func TestInlineNonJSDocTrailingComments_ParamList(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	fn := pt.parse("function f(x /* first */ , y ) {}").GetFirstChild()
	assertNode(t, fn).HasType(rhino.Token_FUNCTION)

	paramListNode := fn.GetSecondChild()
	xNode := paramListNode.GetFirstChild()
	yNode := xNode.GetNext()

	assertThat(t, xNode.GetTrailingNonJSDocCommentString()).Contains("/* first */")
	assertThat(t, yNode.GetNonJSDocCommentString()).IsEmpty()
}

// Java: ParserTest.testInlineNonJSDocTrailingComments_formalParamList_SingleParam()
func TestInlineNonJSDocTrailingComments_formalParamList_SingleParam(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	fn := pt.parse("function f(x /* first */) {}").GetFirstChild()
	assertNode(t, fn).HasType(rhino.Token_FUNCTION)

	paramListNode := fn.GetSecondChild()
	xNode := paramListNode.GetFirstChild()
	assertNode(t, xNode).HasType(rhino.Token_NAME)
	assertThat(t, xNode.GetTrailingNonJSDocCommentString()).Contains("/* first */")
}

// Java: ParserTest.testInlineNonJSDocTrailingComments_ParamList_MultiLine()
func TestInlineNonJSDocTrailingComments_ParamList_MultiLine(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	fn :=
		pt.parse(
			`function f(
  x,// first
  y // second
){}
`).
			GetFirstChild()
	assertNode(t, fn).HasType(rhino.Token_FUNCTION)

	paramListNode := fn.GetSecondChild()
	xNode := paramListNode.GetFirstChild()
	yNode := paramListNode.GetSecondChild()

	assertThat(t, xNode.GetTrailingNonJSDocCommentString()).IsEqualTo("// first")

	assertThat(t, yNode.GetTrailingNonJSDocCommentString()).IsEqualTo("// second")
}

// Java: ParserTest.testInlineNonJSDocTrailingComments_ParamList_MultiLine_BlockComments()
func TestInlineNonJSDocTrailingComments_ParamList_MultiLine_BlockComments(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	fn :=
		pt.parse(
			`function f(
  x, /* first */
  y /* second */
) {}
`).
			GetFirstChild()
	assertNode(t, fn).HasType(rhino.Token_FUNCTION)

	paramListNode := fn.GetSecondChild()
	xNode := paramListNode.GetFirstChild()
	yNode := paramListNode.GetSecondChild()

	assertThat(t, xNode.GetTrailingNonJSDocCommentString()).IsEqualTo("/* first */")

	assertThat(t, yNode.GetTrailingNonJSDocCommentString()).IsEqualTo("/* second */")
}

// Java: ParserTest.testInlineNonJSDocTrailingComments_ParamList_MultiLine_SingleBlockComments()
func TestInlineNonJSDocTrailingComments_ParamList_MultiLine_SingleBlockComments(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	fn :=
		pt.parse(
			`function f(x, /* first */
y
) {}
`).
			GetFirstChild()
	assertNode(t, fn).HasType(rhino.Token_FUNCTION)

	paramListNode := fn.GetSecondChild()
	xNode := paramListNode.GetFirstChild()
	yNode := paramListNode.GetSecondChild()

	assertThat(t, xNode.GetTrailingNonJSDocCommentString()).IsEqualTo("/* first */")

	assertThat(t, yNode.GetNonJSDocComment()).IsNull()
}

// Java: ParserTest.testNonJSDocTrailingCommentOnConstant()
func TestNonJSDocTrailingCommentOnConstant(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	cnst := pt.parse("const A = 1; // comment").GetFirstChild()
	assertNode(t, cnst).HasType(rhino.Token_CONST)

	assertThat(t, cnst.GetTrailingNonJSDocCommentString()).IsEqualTo("// comment")
}

// Java: ParserTest.testNonJSDocTrailingCommentOnConstantNoWhitespace()
func TestNonJSDocTrailingCommentOnConstantNoWhitespace(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	cnst := pt.parse("const A = 1;// comment").GetFirstChild()
	assertNode(t, cnst).HasType(rhino.Token_CONST)

	assertThat(t, cnst.GetTrailingNonJSDocCommentString()).IsEqualTo("// comment")
}

// Java: ParserTest.testNonJSDocTrailingCommentOnConstantWithMoreWhitespace()
func TestNonJSDocTrailingCommentOnConstantWithMoreWhitespace(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	cnst := pt.parse("const A = 1;   // comment").GetFirstChild()
	assertNode(t, cnst).HasType(rhino.Token_CONST)

	assertThat(t, cnst.GetTrailingNonJSDocCommentString()).IsEqualTo("// comment")
}

// Java: ParserTest.testNonJSDocTrailingCommentOnConstantFollowedByConstant()
func TestNonJSDocTrailingCommentOnConstantFollowedByConstant(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	cnst :=
		pt.parse(
			`const A = 1; // comment

const B = 2;
`).
			GetFirstChild()
	assertNode(t, cnst).HasType(rhino.Token_CONST)

	assertThat(t, cnst.GetTrailingNonJSDocCommentString()).IsEqualTo("// comment")
}

// Java: ParserTest.testMultipleNonJSDocTrailingComments()
func TestMultipleNonJSDocTrailingComments(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	n :=
		pt.parse(
			`const A = 1; /* A1 */ /* A2 */ // A3

const B = 2;
`)
	a := n.GetFirstChild()
	assertNode(t, a).HasType(rhino.Token_CONST)
	b := n.GetLastChild()
	assertNode(t, a).HasType(rhino.Token_CONST)
	assertThat(t, b.GetNonJSDocCommentString()).IsEqualTo("/* A2 */// A3")
}

// Java: ParserTest.testNonJSDocTrailingCommentAfterFunction()
func TestNonJSDocTrailingCommentAfterFunction(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	n := pt.parse("function foo(){} // comment").GetFirstChild()
	assertNode(t, n).HasType(rhino.Token_FUNCTION)

	assertThat(t, n.GetTrailingNonJSDocCommentString()).IsEqualTo("// comment")
}

// Java: ParserTest.testNonJSDocTrailingCommentAfterFunctionCall()
func TestNonJSDocTrailingCommentAfterFunctionCall(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	n :=
		pt.parse(
			`function g(){ f(); // comment
f();}
`)
	exprRes := n.GetFirstChild().GetLastChild().GetFirstChild()
	assertNode(t, exprRes).HasType(rhino.Token_EXPR_RESULT)

	assertThat(t, exprRes.GetTrailingNonJSDocCommentString()).IsEqualTo("// comment")
}

// Java: ParserTest.testNonJSDocTrailingCommentAfterFunctionCallInBlock()
func TestNonJSDocTrailingCommentAfterFunctionCallInBlock(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	n :=
		pt.parse(
			`if (true) {
  f1(); // comment1 on f1()
  // comment2
  // comment3
}
`)
	exprRes := n.GetFirstChild().GetLastChild().GetFirstChild()
	assertNode(t, exprRes).HasType(rhino.Token_EXPR_RESULT)

	assertThat(t, exprRes.GetTrailingNonJSDocCommentString()).
		IsEqualTo(`// comment1 on f1()
// comment2
// comment3`)
}

// Java: ParserTest.testLastNonJSDocCommentInBlock()
func TestLastNonJSDocCommentInBlock(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	n :=
		pt.parse(
			`if (true) {
  f();
  /* comment */
}
`)
	exprRes := n.GetFirstChild().GetLastChild().GetFirstChild()
	assertNode(t, exprRes).HasType(rhino.Token_EXPR_RESULT)

	assertThat(t, exprRes.GetTrailingNonJSDocCommentString()).IsEqualTo("\n/* comment */")
}

// Java: ParserTest.testLastNonJSDocCommentInBlockWithBlankLines()
func TestLastNonJSDocCommentInBlockWithBlankLines(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	n :=
		pt.parse(
			`if (true) {
  f();


  /* comment */
}
`)
	exprRes := n.GetFirstChild().GetLastChild().GetFirstChild()
	assertNode(t, exprRes).HasType(rhino.Token_EXPR_RESULT)

	// TODO(b/242294987): This should keep the blank lines.
	assertThat(t, exprRes.GetTrailingNonJSDocCommentString()).IsEqualTo("\n/* comment */")
}

// Java: ParserTest.testInlineCommentInFunctionCallInBlock()
func TestInlineCommentInFunctionCallInBlock(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	n :=
		pt.parse(
			`if (true) {
  f(0, 1 /* comment */);}
`)
	exprRes := n.GetFirstChild().GetLastChild().GetFirstChild()
	assertNode(t, exprRes).HasType(rhino.Token_EXPR_RESULT)
	// TODO(b/242294987): This should not be an "end of block" comment (which we treat as trailing
	// comment on the last child), but a comment on the argument.
	assertThat(t, exprRes.GetTrailingNonJSDocCommentString()).IsEqualTo("\n/* comment */")
}

// Java: ParserTest.testNonJSDocBigCommentInbetween()
func TestNonJSDocBigCommentInbetween(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	n :=
		pt.parse(
			`let x = 0; // comment on x
//
// more comment
let y = 1;
`)

	fstLetDecl := n.GetFirstChild()
	sndLetDecl := n.GetLastChild()

	assertNode(t, fstLetDecl).HasType(rhino.Token_LET)
	assertThat(t, fstLetDecl.GetTrailingNonJSDocCommentString()).IsEqualTo("// comment on x")
	assertNode(t, sndLetDecl).HasType(rhino.Token_LET)
	assertThat(t, sndLetDecl.GetNonJSDocCommentString()).IsEqualTo("//\n// more comment")
}

// Java: ParserTest.testInlineNonJSDocCommentsOnSeparateLetDeclarations()
func TestInlineNonJSDocCommentsOnSeparateLetDeclarations(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS

	n := pt.parse("let a /* leading a */ = {c} /* trailing */; let /* leading b */ b  = {d};")
	letADecl := n.GetFirstChild()
	letBDecl := n.GetLastChild()

	assertNode(t, letADecl).HasType(rhino.Token_LET)
	assertNode(t, letBDecl).HasType(rhino.Token_LET)
	assertThat(t, letADecl.GetFirstFirstChild().GetNonJSDocCommentString()).
		Contains("/* leading a */")
	assertThat(t, letADecl.GetTrailingNonJSDocCommentString()).Contains("/* trailing */")
	assertThat(t, letBDecl.GetFirstChild().GetNonJSDocCommentString()).Contains("/* leading b */")
}

// Java: ParserTest.testMultipleInline_LineCommentsAttachment()
func TestMultipleInline_LineCommentsAttachment(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	fn :=
		pt.parse(
			`function f(
  // blah1
  x,
  // blah2
  y) {}
`).
			GetFirstChild()
	assertNode(t, fn).HasType(rhino.Token_FUNCTION)

	xNode := fn.GetSecondChild().GetFirstChild()
	assertThat(t, xNode.GetNonJSDocCommentString()).Contains("// blah1")

	yNode := fn.GetSecondChild().GetSecondChild()
	assertThat(t, yNode.GetNonJSDocCommentString()).IsEqualTo("// blah2")
}

// Java: ParserTest.testMultipleInline_MixedCommentsAttachment()
func TestMultipleInline_MixedCommentsAttachment(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	fn :=
		pt.parse(
			`function f(
  /* blah1 */ x,
  // blah2
  y
) {}
`).
			GetFirstChild()
	assertNode(t, fn).HasType(rhino.Token_FUNCTION)

	xNode := fn.GetSecondChild().GetFirstChild()
	assertThat(t, xNode.GetNonJSDocCommentString()).Contains("/* blah1 */")

	yNode := fn.GetSecondChild().GetSecondChild()
	assertThat(t, yNode.GetNonJSDocCommentString()).IsEqualTo("// blah2")
}

// Java: ParserTest.testMultipleInline_NonJSDocCommentsGetAttachedToSameNode()
func TestMultipleInline_NonJSDocCommentsGetAttachedToSameNode(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	fn := pt.parse(`function f(/* blah1 */
// blah
 x) {}`).GetFirstChild()
	assertNode(t, fn).HasType(rhino.Token_FUNCTION)

	xNode := fn.GetSecondChild().GetFirstChild()
	assertThat(t, xNode.GetNonJSDocCommentString()).Contains("/* blah1 */\n// blah")
}

// Java: ParserTest.testBoth_TrailingAndNonTrailing_NonJSDocCommentsGetAttachedToSameNode_MultiLine()
func TestBoth_TrailingAndNonTrailing_NonJSDocCommentsGetAttachedToSameNode_MultiLine(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	fn :=
		pt.parse(
			`function f(
     /* blah1 */
     x // blah
  ) {}
`).
			GetFirstChild()
	assertNode(t, fn).HasType(rhino.Token_FUNCTION)

	xNode := fn.GetSecondChild().GetFirstChild()
	assertThat(t, xNode.GetNonJSDocCommentString()).Contains("/* blah1 */")
	assertThat(t, xNode.GetTrailingNonJSDocCommentString()).Contains("// blah")
}

// Java: ParserTest.testEndOfFileNonJSDocComments_lineComment()
func TestEndOfFileNonJSDocComments_lineComment(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	scriptNode :=
		pt.parse(
			`function f1() {}
// first
f1();
// second
`)

	assertNode(t, scriptNode).HasType(rhino.Token_SCRIPT)

	exprNode := scriptNode.GetLastChild()
	assertNode(t, exprNode).HasType(rhino.Token_EXPR_RESULT)

	assertThat(t, scriptNode.GetTrailingNonJSDocCommentString()).IsEqualTo("// second")
}

// Java: ParserTest.testEndOfFileNonJSDocComments_blockComment()
func TestEndOfFileNonJSDocComments_blockComment(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	scriptNode :=
		pt.parse(
			`function f1() {}
// first
f1();
/* second */
`)

	assertNode(t, scriptNode).HasType(rhino.Token_SCRIPT)

	exprNode := scriptNode.GetLastChild()
	assertNode(t, exprNode).HasType(rhino.Token_EXPR_RESULT)

	assertThat(t, scriptNode.GetTrailingNonJSDocCommentString()).IsEqualTo("/* second */")
}

// Java: ParserTest.testEndOfFileNonJSDocComments_manyComments()
func TestEndOfFileNonJSDocComments_manyComments(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	scriptNode :=
		pt.parse(
			`function f1() {}
// first
f1();
// second
/* third */
// fourth
`)

	assertNode(t, scriptNode).HasType(rhino.Token_SCRIPT)

	exprNode := scriptNode.GetLastChild()
	assertNode(t, exprNode).HasType(rhino.Token_EXPR_RESULT)

	assertThat(t, scriptNode.GetTrailingNonJSDocCommentString()).
		IsEqualTo(`// second
/* third */
// fourth`)
}

// Java: ParserTest.testEndOfFileNonJSDocComments()
func TestEndOfFileNonJSDocComments(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	scriptNode :=
		pt.parse(
			`function f1() {}
if (true) {
// first
f1(); // second
}
// third
`)

	assertNode(t, scriptNode).HasType(rhino.Token_SCRIPT)

	exprNode := scriptNode.GetLastChild().GetLastChild().GetFirstChild()
	assertNode(t, exprNode).HasType(rhino.Token_EXPR_RESULT)
	assertThat(t, exprNode.GetNonJSDocCommentString()).IsEqualTo("// first")

	callNode := scriptNode.GetLastChild().GetLastChild().GetLastChild()
	assertThat(t, callNode.GetTrailingNonJSDocCommentString()).IsEqualTo("// second")

	assertThat(t, scriptNode.GetTrailingNonJSDocCommentString()).IsEqualTo("// third")
}

// Java: ParserTest.testBoth_TrailingAndNonTrailing_NonJSDocCommentsGetAttachedToSameNode_SingleLine()
func TestBoth_TrailingAndNonTrailing_NonJSDocCommentsGetAttachedToSameNode_SingleLine(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	fn :=
		pt.parse(
			`function f(/* blah1 */ x // blah
) {}
`).
			GetFirstChild()
	assertNode(t, fn).HasType(rhino.Token_FUNCTION)

	xNode := fn.GetSecondChild().GetFirstChild()
	assertThat(t, xNode.GetNonJSDocCommentString()).Contains("/* blah1 */")
	assertThat(t, xNode.GetTrailingNonJSDocCommentString()).Contains("/ blah")
}

// Java: ParserTest.testBothJSDocAndNonJSDocCommentsGetAttached()
func TestBothJSDocAndNonJSDocCommentsGetAttached(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	fn := pt.parse("function f(/** string */ // nonJSDoc\n x) {}").GetFirstChild()
	assertNode(t, fn).HasType(rhino.Token_FUNCTION)

	xNode := fn.GetSecondChild().GetFirstChild()
	info := xNode.GetJSDocInfo()
	assertThat(t, info).IsNotNull()
	pt.assertNodeHasJSDocInfoWithJSType(xNode, pt.STRING_TYPE)
	assertThat(t, xNode.GetNonJSDocCommentString()).Contains("// nonJSDoc")
}

// Java: ParserTest.testBothNonJSDocAndJSDocCommentsGetAttached()
func TestBothNonJSDocAndJSDocCommentsGetAttached(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	fn := pt.parse("function f(// nonJSDoc\n /** string */ x) {}").GetFirstChild()
	assertNode(t, fn).HasType(rhino.Token_FUNCTION)

	xNode := fn.GetSecondChild().GetFirstChild()
	info := xNode.GetJSDocInfo()
	assertThat(t, info).IsNotNull()
	pt.assertNodeHasJSDocInfoWithJSType(xNode, pt.STRING_TYPE)
	assertThat(t, xNode.GetNonJSDocCommentString()).Contains("// nonJSDoc")
}

// Java: ParserTest.testInlineTrailingNonJSDocComments_FunctionArgsAndBody()
func TestInlineTrailingNonJSDocComments_FunctionArgsAndBody(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	fn := pt.parse("function f(x /* first */ ) { /* second */ let y;}").GetFirstChild()
	assertNode(t, fn).HasType(rhino.Token_FUNCTION)

	xNode := fn.GetSecondChild().GetOnlyChild()
	yNode := fn.GetLastChild().GetFirstChild()

	assertThat(t, xNode.GetTrailingNonJSDocCommentString()).Contains("/* first */")
	assertThat(t, yNode.GetNonJSDocCommentString()).Contains("/* second */")
}

// Java: ParserTest.testInlineNonJSDocComments_FunctionCall()
func TestInlineNonJSDocComments_FunctionCall(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS

	exprRes := pt.parse("function f(x) {let y;}; f( /* first */  1)").GetLastChild()
	assertNode(t, exprRes).HasType(rhino.Token_EXPR_RESULT)

	call := exprRes.GetFirstChild()
	oneArgNode := call.GetSecondChild()

	assertThat(t, oneArgNode.GetNonJSDocCommentString()).Contains("/* first */")
}

// Java: ParserTest.testInlineTrailingNonJSDocComments_MultipleArgs()
func TestInlineTrailingNonJSDocComments_MultipleArgs(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS

	exprRes := pt.parse("function f(x, y) {}; f( 1 /* first */, 2 );").GetLastChild()
	assertNode(t, exprRes).HasType(rhino.Token_EXPR_RESULT)

	call := exprRes.GetFirstChild()
	oneArgNode := call.GetSecondChild()
	twoArgNode := oneArgNode.GetNext()

	assertThat(t, oneArgNode.GetTrailingNonJSDocCommentString()).Contains("/* first */")
	assertThat(t, twoArgNode.GetNonJSDocCommentString()).IsEmpty()
}

// Java: ParserTest.testInlineTrailingNonJSDocComments_SingleArgument()
func TestInlineTrailingNonJSDocComments_SingleArgument(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS

	exprRes := pt.parse("function f(x, y) {}; f( 1 /* first */);").GetLastChild()
	assertNode(t, exprRes).HasType(rhino.Token_EXPR_RESULT)

	call := exprRes.GetFirstChild()
	oneArgNode := call.GetSecondChild()

	assertThat(t, oneArgNode.GetTrailingNonJSDocCommentString()).Contains("/* first */")
}

// Java: ParserTest.testInlineJSDocAttachment2()
func TestInlineJSDocAttachment2(t *testing.T) {
	pt := newParserTest(t)

	fn := pt.parse("function f(/** ? */ x) {}").GetFirstChild()
	assertNode(t, fn).HasType(rhino.Token_FUNCTION)

	info := fn.GetSecondChild().GetFirstChild().GetJSDocInfo()
	assertThat(t, info).IsNotNull()
	pt.AssertTypeEqualsExpr(pt.UNKNOWN_TYPE, info.GetType())
}

// Java: ParserTest.testInlineJSDocAttachment3()
func TestInlineJSDocAttachment3(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("function f(/** @type {string} */ x) {}")
}

// Java: ParserTest.testInlineJSDocAttachment4()
func TestInlineJSDocAttachment4(t *testing.T) {
	pt := newParserTest(t)

	pt.parse(
		`function f(/**
 * @type {string}
 */ x) {}
`)
}

// Java: ParserTest.testInlineJSDocAttachment5()
func TestInlineJSDocAttachment5(t *testing.T) {
	pt := newParserTest(t)

	vardecl := pt.parse("var /** string */ x = 'asdf';").GetFirstChild()
	info := vardecl.GetFirstChild().GetJSDocInfo()
	assertThat(t, info).IsNotNull()
	assertThat(t, info.HasType()).IsTrue()
	pt.AssertTypeEqualsExpr(pt.STRING_TYPE, info.GetType())
}

// Java: ParserTest.testInlineJSDocAttachment6()
func TestInlineJSDocAttachment6(t *testing.T) {
	pt := newParserTest(t)

	fn := pt.parse("function f(/** {attr: number} */ x) {}").GetFirstChild()
	assertNode(t, fn).HasType(rhino.Token_FUNCTION)

	info := fn.GetSecondChild().GetFirstChild().GetJSDocInfo()
	assertThat(t, info).IsNotNull()
	pt.AssertTypeEqualsExpr(
		pt.CreateRecordTypeBuilder().AddProperty("attr", pt.NUMBER_TYPE, nil).Build(), info.GetType())
}

// Java: ParserTest.testInlineJSDocWithOptionalType()
func TestInlineJSDocWithOptionalType(t *testing.T) {
	pt := newParserTest(t)

	fn := pt.parse("function f(/** string= */ x) {}").GetFirstChild()
	assertNode(t, fn).HasType(rhino.Token_FUNCTION)

	info := fn.GetSecondChild().GetFirstChild().GetJSDocInfo()
	assertThat(t, info.GetType().IsOptionalArg()).IsTrue()
}

// Java: ParserTest.testInlineJSDocWithVarArgs()
func TestInlineJSDocWithVarArgs(t *testing.T) {
	pt := newParserTest(t)

	fn := pt.parse("function f(/** ...string */ x) {}").GetFirstChild()
	assertNode(t, fn).HasType(rhino.Token_FUNCTION)

	info := fn.GetSecondChild().GetFirstChild().GetJSDocInfo()
	assertThat(t, info.GetType().IsVarArgs()).IsTrue()
}

// Java: ParserTest.testInlineJSDocReturnType()
func TestInlineJSDocReturnType(t *testing.T) {
	pt := newParserTest(t)

	fn := pt.parse("function /** string */ f(x) {}").GetFirstChild()
	assertNode(t, fn).HasType(rhino.Token_FUNCTION)

	info := fn.GetFirstChild().GetJSDocInfo()
	assertThat(t, info.HasType()).IsTrue()
	pt.AssertTypeEqualsExpr(pt.STRING_TYPE, info.GetType())
}

// Java: ParserTest.testInlineJSDocReturnType_generator1()
func TestInlineJSDocReturnType_generator1(t *testing.T) {
	pt := newParserTest(t)

	fn := pt.parse("function * /** string */ f(x) {}").GetFirstChild()
	assertNode(t, fn).HasType(rhino.Token_FUNCTION)

	info := fn.GetFirstChild().GetJSDocInfo()
	assertThat(t, info.HasType()).IsTrue()
	pt.AssertTypeEqualsExpr(pt.STRING_TYPE, info.GetType())
}

// Java: ParserTest.testInlineJSDocReturnType_generator2()
func TestInlineJSDocReturnType_generator2(t *testing.T) {
	pt := newParserTest(t)

	fn := pt.parse("function /** string */ *f(x) {}").GetFirstChild()
	assertNode(t, fn).HasType(rhino.Token_FUNCTION)

	info := fn.GetFirstChild().GetJSDocInfo()
	assertThat(t, info.HasType()).IsTrue()
	pt.AssertTypeEqualsExpr(pt.STRING_TYPE, info.GetType())
}

// Java: ParserTest.testInlineJSDocReturnType_async()
func TestInlineJSDocReturnType_async(t *testing.T) {
	pt := newParserTest(t)

	fn := pt.parse("async function /** string */ f(x) {}").GetFirstChild()
	assertNode(t, fn).HasType(rhino.Token_FUNCTION)

	info := fn.GetFirstChild().GetJSDocInfo()
	assertThat(t, info.HasType()).IsTrue()
	pt.AssertTypeEqualsExpr(pt.STRING_TYPE, info.GetType())
}

// Java: ParserTest.testIncorrectJSDocDoesNotAlterJSParsing1()
func TestIncorrectJSDocDoesNotAlterJSParsing1(t *testing.T) {
	pt := newParserTest(t)

	pt.assertNodeEquality(
		pt.parse("var a = [1,2]"),
		pt.parseWarning("/** @type {Array<number} */var a = [1,2]", parserTest_MISSING_GT_MESSAGE))
}

// Java: ParserTest.testIncorrectJSDocDoesNotAlterJSParsing2()
func TestIncorrectJSDocDoesNotAlterJSParsing2(t *testing.T) {
	pt := newParserTest(t)

	pt.assertNodeEquality(
		pt.parse("var a = [1,2]"),
		pt.parseWarning("/** @type {Array.<number}*/var a = [1,2]", parserTest_MISSING_GT_MESSAGE))
}

// Java: ParserTest.testIncorrectJSDocDoesNotAlterJSParsing3()
func TestIncorrectJSDocDoesNotAlterJSParsing3(t *testing.T) {
	pt := newParserTest(t)

	pt.assertNodeEquality(
		pt.parse("C.prototype.say=function(nums) {alert(nums.join(','));};"),
		pt.parseWarning(
			`/** @param {Array.<number} nums */
C.prototype.say=function(nums) {alert(nums.join(','));};
`,
			parserTest_MISSING_GT_MESSAGE))
}

// Java: ParserTest.testIncorrectJSDocDoesNotAlterJSParsing4()
func TestIncorrectJSDocDoesNotAlterJSParsing4(t *testing.T) {
	pt := newParserTest(t)

	pt.assertNodeEquality(
		pt.parse("C.prototype.say=function(nums) {alert(nums.join(','));};"),
		pt.parse(
			`/** @return {boolean} */
C.prototype.say=function(nums) {alert(nums.join(','));};
`))
}

// Java: ParserTest.testIncorrectJSDocDoesNotAlterJSParsing5()
func TestIncorrectJSDocDoesNotAlterJSParsing5(t *testing.T) {
	pt := newParserTest(t)

	pt.assertNodeEquality(
		pt.parse("C.prototype.say=function(nums) {alert(nums.join(','));};"),
		pt.parse(
			`/** @param {boolean} this is some string*/
C.prototype.say=function(nums) {alert(nums.join(','));};
`))
}

// Java: ParserTest.testIncorrectJSDocDoesNotAlterJSParsing6()
func TestIncorrectJSDocDoesNotAlterJSParsing6(t *testing.T) {
	pt := newParserTest(t)

	pt.assertNodeEquality(
		pt.parse("C.prototype.say=function(nums) {alert(nums.join(','));};"),
		pt.parseWarning(
			`/** @param {bool!*%E$} */
C.prototype.say=function(nums) {alert(nums.join(','));};
`,
			"Bad type annotation. expected closing }"+JsDocInfoParser_BAD_TYPE_WIKI_LINK,
			"Bad type annotation. expecting a variable name in a @param tag."+
				JsDocInfoParser_BAD_TYPE_WIKI_LINK))
}

// Java: ParserTest.testIncorrectJSDocDoesNotAlterJSParsing7()
func TestIncorrectJSDocDoesNotAlterJSParsing7(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true

	pt.assertNodeEquality(
		pt.parse("C.prototype.say=function(nums) {alert(nums.join(','));};"),
		pt.parseWarning(
			`/** @see */
C.prototype.say=function(nums) {alert(nums.join(','));};
`,
			"@see tag missing description"))
}

// Java: ParserTest.testIncorrectJSDocDoesNotAlterJSParsing8()
func TestIncorrectJSDocDoesNotAlterJSParsing8(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true

	pt.assertNodeEquality(
		pt.parse("C.prototype.say=function(nums) {alert(nums.join(','));};"),
		pt.parseWarning(
			`/** @author */
C.prototype.say=function(nums) {alert(nums.join(','));};
`,
			"@author tag missing author"))
}

// Java: ParserTest.testIncorrectJSDocDoesNotAlterJSParsing9()
func TestIncorrectJSDocDoesNotAlterJSParsing9(t *testing.T) {
	pt := newParserTest(t)

	pt.assertNodeEquality(
		pt.parse("C.prototype.say=function(nums) {alert(nums.join(','));};"),
		pt.parseWarning(
			`/** @someillegaltag */
C.prototype.say=function(nums) {alert(nums.join(','));};
`,
			"illegal use of unknown JSDoc tag \"someillegaltag\"; ignoring it. Place another"+
				" character before the @ to stop JSCompiler from parsing it as an annotation."))
}

// Java: ParserTest.testMisplacedDescAnnotation_noWarning()
func TestMisplacedDescAnnotation_noWarning(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("/** @desc Foo. */ var MSG_BAR = goog.getMsg('hello');")
	pt.parse("/** @desc Foo. */ x.y.z.MSG_BAR = goog.getMsg('hello');")
	pt.parse("/** @desc Foo. */ MSG_BAR = goog.getMsg('hello');")
	pt.parse("var msgs = {/** @desc x */ MSG_X: goog.getMsg('x')}")
}

// Java: ParserTest.testUnescapedSlashInRegexpCharClass()
func TestUnescapedSlashInRegexpCharClass(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("var foo = /[/]/;")
	pt.parse("var foo = /[hi there/]/;")
	pt.parse("var foo = /[/yo dude]/;")
	pt.parse("var foo = /\\/[@#$/watashi/wa/suteevu/desu]/;")
}

// Java: ParserTest.testMalformedRegexp()
func TestMalformedRegexp(t *testing.T) {
	pt := newParserTest(t)

	// Simple repro case
	js := "var x = com\\"
	pt.parseError(js, "Invalid escape sequence")

	// The original repro case as reported.
	js =
		`(function() {
  var url="";
  switch(true)
  {
    case /a.com\/g|l.i/N/.test(url):
      return "";
    case /b.com\/T/.test(url):
      return "";
  }
}
)();
`
	pt.parseError(js, "primary expression expected")
}

// Java: ParserTest.testParse()
func TestParse(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY
	a := rhino.Node_NewStringWithToken(rhino.Token_NAME, "a")
	a.AddChildToFront(rhino.Node_NewStringWithToken(rhino.Token_NAME, "b"))
	testCases := []parserTest_ParserResult{
		parserTest_ParserResult{"3;", pt.createScript(rhino.NewNodeWithChild(rhino.Token_EXPR_RESULT, rhino.Node_NewNumber(3.0)))},
		parserTest_ParserResult{"var a = b;", pt.createScript(rhino.NewNodeWithChild(rhino.Token_VAR, a))},
	}

	for _, testCase := range testCases {
		pt.assertNodeEquality(testCase.node, pt.parse(testCase.code))
	}
}

// Java: ParserTest.testPostfixExpression()
func TestPostfixExpression(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("a++")
	pt.parse("a.b--")
	pt.parse("a[0]++")
	pt.parse("/** @type {number} */ (a)++;")

	pt.parseError("a()++", "Invalid postfix increment operand.")
	pt.parseError("(new C)--", "Invalid postfix decrement operand.")
	pt.parseError("this++", "Invalid postfix increment operand.")
	pt.parseError("(a--)++", "Invalid postfix increment operand.")
	pt.parseError("(+a)++", "Invalid postfix increment operand.")
	pt.parseError("[1,2]++", "Invalid postfix increment operand.")
	pt.parseError("'literal'++", "Invalid postfix increment operand.")
	pt.parseError("/** @type {number} */ (a())++;", "Invalid postfix increment operand.")
}

// Java: ParserTest.testUnaryExpression()
func TestUnaryExpression(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parse("delete a.b")
	pt.parse("delete a[0]")
	pt.parse("void f()")
	pt.parse("typeof new C")
	pt.parse("++a[0]")
	pt.parse("--a.b")
	pt.parse("+{a: 1}")
	pt.parse("-[1,2]")
	pt.parse("~'42'")
	pt.expectFeatures(parser.FeatureSet_Feature_SUPER)
	pt.parse("!super.a")
	pt.expectFeatures()

	pt.parseError("delete f()", "Invalid delete operand. Only properties can be deleted.")
	pt.parseError("++a++", "Invalid prefix increment operand.")
	pt.parseError("--{a: 1}", "Invalid prefix decrement operand.")
	pt.parseError("++this", "Invalid prefix increment operand.")
	pt.parseError("++(-a)", "Invalid prefix increment operand.")
	pt.parseError("++{a: 1}", "Invalid prefix increment operand.")
	pt.parseError("++'literal'", "Invalid prefix increment operand.")
	pt.parseError("++delete a.b", "Invalid prefix increment operand.")
}

// Java: ParserTest.testUnaryExpressionWithBigInt()
func TestUnaryExpressionWithBigInt(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("+1n", "Cannot convert a BigInt value to a number")
	pt.parseError("delete 6n", "Invalid delete operand. Only properties can be deleted.")
}

// Java: ParserTest.testAutomaticSemicolonInsertion()
func TestAutomaticSemicolonInsertion(t *testing.T) {
	pt := newParserTest(t)

	// var statements
	pt.assertNodeEquality(pt.parse("var x = 1\nvar y = 2"), pt.parse("var x = 1; var y = 2;"))
	pt.assertNodeEquality(pt.parse("var x = 1\n, y = 2"), pt.parse("var x = 1, y = 2;"))

	// assign statements
	pt.assertNodeEquality(pt.parse("x = 1\ny = 2"), pt.parse("x = 1; y = 2;"))

	pt.assertNodeEquality(
		pt.parse("x = 1\n;y = 2"), //
		pt.parse("x = 1; y = 2;"))
	pt.assertNodeEquality(
		pt.parse("if (true) 1\n; else {}"), //
		pt.parse("if (true) 1; else {}"))

	// if/else statements
	pt.assertNodeEquality(pt.parse("if (x)\n;else{}"), pt.parse("if (x) {} else {}"))
}

// Java: ParserTest.testAutomaticSemicolonInsertion_curly()
func TestAutomaticSemicolonInsertion_curly(t *testing.T) {
	pt := newParserTest(t)

	pt.assertNodeEquality(
		pt.parse("while (true) { 1 }"), //
		pt.parse("while (true) { 1; }"))
}

// Java: ParserTest.testAutomaticSemicolonInsertion_doWhile()
func TestAutomaticSemicolonInsertion_doWhile(t *testing.T) {
	pt := newParserTest(t)

	pt.assertNodeEquality(
		pt.parse("do {} while (true) 1;"), //
		pt.parse("do {} while (true); 1;"))
}

// Java: ParserTest.testAutomaticSemicolonInsertion_examplesFromSpec()
func TestAutomaticSemicolonInsertion_examplesFromSpec(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("{ 1 2 } 3", parserTest_SEMICOLON_EXPECTED)

	pt.assertNodeEquality(pt.parse("{ 1\n2 } 3"), pt.parse("{ 1; 2; } 3;"))

	pt.parseError("for (a; b\n)", "';' expected")

	pt.assertNodeEquality(
		pt.parse("function f() { return\na + b }"), pt.parse("function f() { return; a + b; }"))

	pt.assertNodeEquality(pt.parse("a = b\n++c"), pt.parse("a = b; ++c;"))

	pt.parseError("if (a > b)\nelse c = d", "primary expression expected")

	pt.assertNodeEquality(pt.parse("a = b + c\n(d + e).print()"), pt.parse("a = b + c(d + e).print()"))
}

// Java: ParserTest.testAutomaticSemicolonInsertion_restrictedRules()
func TestAutomaticSemicolonInsertion_restrictedRules(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("x\n++;", "primary expression expected")
	pt.assertNodeEquality(
		pt.parse("function f() { return\n1; }"), //
		pt.parse("function f() { return;1; }"))
	pt.assertNodeEquality(
		pt.parse("while (true) { continue\nlabel; }"), //
		pt.parse("while (true) { continue; label; }"))
	pt.assertNodeEquality(
		pt.parse("while (true) { break\nlabel; }"), //
		pt.parse("while (true) { break; label; }"))
	pt.parseError("throw\n1;", "semicolon/newline not allowed after 'throw'")
	pt.parseError("yield\nvalue;", "primary expression expected")
	pt.parseError("()\n=> 1;", "No newline allowed before '=>'")
}

// Java: ParserTest.testMethodInObjectLiteral()
func TestMethodInObjectLiteral(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_MEMBER_DECLARATIONS)
	pt.testMethodInObjectLiteral("var a = {b() {}};")
	pt.testMethodInObjectLiteral("var a = {b() { alert('b'); }};")

	// Static methods not allowed in object literals.
	pt.expectFeatures()
	pt.parseError(
		"var a = {static b() { alert('b'); }};", "Cannot use keyword in short object literal")
}

// Java: ParserTest.testBigIntAsObjectLiteralPropertyName()
func TestBigIntAsObjectLiteralPropertyName(t *testing.T) {
	pt := newParserTest(t)

	objectLit :=
		pt.parse("({1n() {}, 1n: 0})"). // SCRIPT
						GetOnlyChild(). // EXPR_RESULT
						GetOnlyChild()  // OBJECTLIT
	assertNode(t, objectLit).HasType(rhino.Token_OBJECTLIT)
	computedProp := objectLit.GetFirstChild()
	assertNode(t, computedProp).HasType(rhino.Token_COMPUTED_PROP)
	bigint := computedProp.GetFirstChild()
	assertNode(t, bigint).HasType(rhino.Token_BIGINT)
	stringKey := objectLit.GetLastChild()
	assertNode(t, stringKey).HasType(rhino.Token_STRING_KEY)
}

// Java: ParserTest.testShorthandObjectProperty()
func TestShorthandObjectProperty(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_SHORTHAND_OBJECT_PROPERTIES)
	pt.testShorthandObjectProperty("var a = {b};")
	pt.testShorthandObjectProperty("var a = {b, c};")
	pt.testShorthandObjectProperty("var a = {b, c: d, e};")
	pt.testShorthandObjectProperty("var a = {type};")
	pt.testShorthandObjectProperty("var a = {declare};")
	pt.testShorthandObjectProperty("var a = {namespace};")
	pt.testShorthandObjectProperty("var a = {module};")

	pt.expectFeatures()
	pt.parseError("var a = { '!@#$%' };", "':' expected")
	pt.parseError("var a = { 123 };", "':' expected")
	pt.parseError("var a = { let };", "Cannot use keyword in short object literal")
	pt.parseError("var a = { else };", "Cannot use keyword in short object literal")
}

// Java: ParserTest.testComputedPropertiesObjLit()
func TestComputedPropertiesObjLit(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_COMPUTED_PROPERTIES)

	// Method
	pt.testComputedProperty("var x = {  [prop + '_']() {} }")
	// NOTE: we treat string and number keys as if they were computed properties for method
	// shorthand, but not getters and setters.
	pt.testComputedProperty("var x = {  'abc'() {} }")
	pt.testComputedProperty("var x = {  123() {} }")

	// Getter
	pt.testComputedProperty("var x = {  get [prop + '_']() {} }")

	// Setter
	pt.testComputedProperty("var x = { set [prop + '_'](val) {} }")

	// Generator method
	pt.mode = Config_LanguageMode_ECMASCRIPT_2015
	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parse("var x = { *[prop + '_']() {} }")
	pt.parse("var x = { *'abc'() {} }")
	pt.parse("var x = { *123() {} }")

	pt.mode = Config_LanguageMode_ECMASCRIPT_2017
	pt.parse("var x = { async [prop + '_']() {} }")
	pt.parse("var x = { async 'abc'() {} }")
	pt.parse("var x = { async 123() {} }")
}

// Java: ParserTest.testComputedMethodClass()
func TestComputedMethodClass(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY
	pt.expectFeatures(parser.FeatureSet_Feature_CLASSES, parser.FeatureSet_Feature_COMPUTED_PROPERTIES)
	pt.parse("class X { [prop + '_']() {} }")
	// Note that we pretend string and number keys are computed property names, because
	// this makes it easier to treat class and object-literal cases consistently.
	pt.parse("class X { 'abc'() {} }")
	pt.parse("class X { 123() {} }")

	pt.parse("class X { static [prop + '_']() {} }")
	pt.parse("class X { static 'abc'() {} }")
	pt.parse("class X { static 123() {} }")

	pt.parse("class X { *[prop + '_']() {} }")
	pt.parse("class X { *'abc'() {} }")
	pt.parse("class X { *123() {} }")

	pt.parse("class X { async [prop + '_']() {} }")
	pt.parse("class X { async 'abc'() {} }")
	pt.parse("class X { async 123() {} }")
}

// Java: ParserTest.testBigIntComputedMethodClass()
func TestBigIntComputedMethodClass(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_CLASSES, parser.FeatureSet_Feature_COMPUTED_PROPERTIES)
	var classMembers *rhino.Node
	var computedProp *rhino.Node
	var bigint *rhino.Node
	bigintValue := jutil.NewBigIntegerFromString("123")
	classMembers =
		pt.parse("class X { 123n() {} }"). // SCRIPT
							GetOnlyChild(). // CLASS
							GetLastChild()  // CLASS_MEMBERS
	assertNode(t, classMembers).HasType(rhino.Token_CLASS_MEMBERS)
	computedProp = classMembers.GetOnlyChild()
	assertNode(t, computedProp).HasType(rhino.Token_COMPUTED_PROP)
	bigint = computedProp.GetFirstChild()
	assertNode(t, bigint).HasType(rhino.Token_BIGINT)
	assertNode(t, bigint).IsBigInt(bigintValue)
	classMembers =
		pt.parse("class X { static 123n() {} }"). // SCRIPT
								GetOnlyChild(). // CLASS
								GetLastChild()  // CLASS_MEMBERS
	assertNode(t, classMembers).HasType(rhino.Token_CLASS_MEMBERS)
	computedProp = classMembers.GetOnlyChild()
	assertNode(t, computedProp).HasType(rhino.Token_COMPUTED_PROP)
	assertNode(t, computedProp).IsStatic()
	bigint = computedProp.GetFirstChild()
	assertNode(t, bigint).HasType(rhino.Token_BIGINT)
	assertNode(t, bigint).IsBigInt(bigintValue)
	classMembers =
		pt.parse("class X { *123n() {} }"). // SCRIPT
							GetOnlyChild(). // CLASS
							GetLastChild()  // CLASS_MEMBERS
	assertNode(t, classMembers).HasType(rhino.Token_CLASS_MEMBERS)
	computedProp = classMembers.GetOnlyChild()
	assertNode(t, computedProp).HasType(rhino.Token_COMPUTED_PROP)
	bigint = computedProp.GetFirstChild()
	assertNode(t, bigint).HasType(rhino.Token_BIGINT)
	assertNode(t, bigint).IsBigInt(bigintValue)
	classMembers =
		pt.parse("class X { async 123n() {} }"). // SCRIPT
								GetOnlyChild(). // CLASS
								GetLastChild()  // CLASS_MEMBERS
	assertNode(t, classMembers).HasType(rhino.Token_CLASS_MEMBERS)
	computedProp = classMembers.GetOnlyChild()
	assertNode(t, computedProp).HasType(rhino.Token_COMPUTED_PROP)
	bigint = computedProp.GetFirstChild()
	assertNode(t, bigint).HasType(rhino.Token_BIGINT)
	assertNode(t, bigint).IsBigInt(bigintValue)
}

// Java: ParserTest.testComputedProperty()
func TestComputedProperty(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_COMPUTED_PROPERTIES)

	pt.testComputedProperty(
		`var prop = 'some complex expression';

var x = {
  [prop]: 'foo'
}
`)

	pt.testComputedProperty(
		`var prop = 'some complex expression';

var x = {
  [prop + '!']: 'foo'
}
`)

	pt.testComputedProperty(
		`var prop;

var x = {
  [prop = 'some expr']: 'foo'
}
`)

	pt.testComputedProperty(
		`var x = {
  [1 << 8]: 'foo'
}
`)

	js :=
		`var x = {
  [1 << 8]: 'foo',
  [1 << 7]: 'bar'
}
`
	pt.mode = Config_LanguageMode_ECMASCRIPT_2015
	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parse(js)
	pt.mode = Config_LanguageMode_ECMASCRIPT5
	warning := requiresLanguageModeMessage(parser.FeatureSet_Feature_COMPUTED_PROPERTIES)
	pt.parseWarning(js, warning, warning)
}

// Java: ParserTest.testTrailingCommaWarning1()
func TestTrailingCommaWarning1(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("var a = ['foo', 'bar'];")
}

// Java: ParserTest.testTrailingCommaWarning2()
func TestTrailingCommaWarning2(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("var a = ['foo',,'bar'];")
}

// Java: ParserTest.testTrailingCommaWarning3()
func TestTrailingCommaWarning3(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT3
	pt.expectFeatures(parser.FeatureSet_Feature_TRAILING_COMMA)
	pt.parseWarning("var a = ['foo', 'bar',];", parserTest_TRAILING_COMMA_MESSAGE)
	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parse("var a = ['foo', 'bar',];")
}

// Java: ParserTest.testTrailingCommaWarning4()
func TestTrailingCommaWarning4(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT3
	pt.expectFeatures(parser.FeatureSet_Feature_TRAILING_COMMA)
	pt.parseWarning("var a = [,];", parserTest_TRAILING_COMMA_MESSAGE)
	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parse("var a = [,];")
}

// Java: ParserTest.testTrailingCommaWarning5()
func TestTrailingCommaWarning5(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("var a = {'foo': 'bar'};")
}

// Java: ParserTest.testTrailingCommaWarning6()
func TestTrailingCommaWarning6(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT3
	pt.expectFeatures(parser.FeatureSet_Feature_TRAILING_COMMA)
	pt.parseWarning("var a = {'foo': 'bar',};", parserTest_TRAILING_COMMA_MESSAGE)
	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parse("var a = {'foo': 'bar',};")
}

// Java: ParserTest.testTrailingCommaWarning7()
func TestTrailingCommaWarning7(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("var a = {,};", "'}' expected")
}

// Java: ParserTest.testCatchClauseForbidden()
func TestCatchClauseForbidden(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("try { } catch (e if true) {}", "')' expected")
}

// Java: ParserTest.testConstForbidden()
func TestConstForbidden(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.expectFeatures(parser.FeatureSet_Feature_CONST_DECLARATIONS)
	pt.parseWarning("const x = 3;", requiresLanguageModeMessage(parser.FeatureSet_Feature_CONST_DECLARATIONS))
}

// Java: ParserTest.testAnonymousFunctionExpression()
func TestAnonymousFunctionExpression(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parseError("function () {}", "'identifier' expected")

	pt.mode = Config_LanguageMode_ECMASCRIPT_2015
	pt.parseError("function () {}", "'identifier' expected")

	pt.isIdeMode = true
	pt.parseError("function () {}", "'identifier' expected", "unnamed function statement")
}

// Java: ParserTest.testAnonymousFunctionExpressionInClosureUnawareCode()
func TestAnonymousFunctionExpressionInClosureUnawareCode(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT_2015
	pt.parseError(
		"/** @closureUnaware */ (function() { function () {} }).call(globalThis)",
		"'identifier' expected")
}

// Java: ParserTest.testArrayDestructuringVar()
func TestArrayDestructuringVar(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.strictMode = Config_StrictMode_SLOPPY
	pt.expectFeatures(parser.FeatureSet_Feature_ARRAY_DESTRUCTURING)
	pt.parseWarning("var [x,y] = foo();", requiresLanguageModeMessage(parser.FeatureSet_Feature_ARRAY_DESTRUCTURING))

	pt.mode = Config_LanguageMode_ECMASCRIPT_2015
	pt.parse("var [x,y] = foo();")
}

// Java: ParserTest.testLHSOfNonVanillaEqualsOperator()
func TestLHSOfNonVanillaEqualsOperator(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY
	pt.mode = Config_LanguageMode_ECMASCRIPT_2015

	// object pattern or array pattern on the lhs of vanilla equals passes
	pt.parse("for (let [i,j] = [2,0]; j < 2; [i,j]  =  [j, j+1]) {}")
	pt.parse("for (let [i,j] = [2,0]; j < 2; {i,j}  =  [j, j+1]) {}")

	// error with object pattern or array pattern on the lhs of +=
	pt.parseError(
		"for (let [i,j] = [2,0]; j < 2; [i,j]  +=  [j, j+1]) {}", "invalid assignment target")
	pt.parseError(
		"for (let [i,j] = [2,0]; j < 2; {i,j}  +=  [j, j+1]) {}", "invalid assignment target")

	pt.parse("let [i,j]  =  [2, 2];")
	pt.parseError("let [i,j]  +=  [2, 2];", "destructuring must have an initializer")
}

// Java: ParserTest.testArrayDestructuringVarInvalid()
func TestArrayDestructuringVarInvalid(t *testing.T) {
	pt := newParserTest(t)

	// arbitrary LHS assignment target not allowed
	pt.parseError(
		"var [x,y[15]] = foo();", "Only an identifier or destructuring pattern is allowed here.")
}

// Java: ParserTest.testArrayDestructuringAssign()
func TestArrayDestructuringAssign(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.strictMode = Config_StrictMode_SLOPPY
	pt.expectFeatures(parser.FeatureSet_Feature_ARRAY_DESTRUCTURING)
	pt.parseWarning("[x,y] = foo();", requiresLanguageModeMessage(parser.FeatureSet_Feature_ARRAY_DESTRUCTURING))

	pt.mode = Config_LanguageMode_ECMASCRIPT_2015
	pt.parse("[x,y] = foo();")
	// arbitrary LHS assignment target is allowed
	pt.parse("[x,y[15]] = foo();")
}

// Java: ParserTest.testArrayDestructuringInitializer()
func TestArrayDestructuringInitializer(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY
	pt.expectFeatures(parser.FeatureSet_Feature_ARRAY_DESTRUCTURING)
	pt.parse("var [x=1,y] = foo();")
	pt.parse("[x=1,y] = foo();")
	pt.parse("var [x,y=2] = foo();")
	pt.parse("[x,y=2] = foo();")

	pt.parse("var [[a] = ['b']] = [];")
	pt.parse("[[a] = ['b']] = [];")
	// arbitrary LHS target allowed in assignment, but not declaration
	pt.parse("[[a.x] = ['b']] = [];")
}

// Java: ParserTest.testArrayDestructuringInitializerInvalid()
func TestArrayDestructuringInitializerInvalid(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError(
		"var [[a.x] = ['b']] = [];",
		"Only an identifier or destructuring pattern is allowed here.")
}

// Java: ParserTest.testArrayDestructuringDeclarationRest()
func TestArrayDestructuringDeclarationRest(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY

	pt.expectFeatures(parser.FeatureSet_Feature_ARRAY_DESTRUCTURING, parser.FeatureSet_Feature_ARRAY_PATTERN_REST)
	pt.parse("var [first, ...rest] = foo();")
	pt.parse("let [first, ...rest] = foo();")
	pt.parse("const [first, ...rest] = foo();")

	// nested destructuring in regular parameters and rest parameters
	pt.parse("var [first, {a, b}, ...[re, st, ...{length}]] = foo();")

	pt.expectFeatures()
	pt.parseError(
		"var [first, ...more = 'default'] = foo();",
		"A default value cannot be specified after '...'")
	pt.parseError("var [first, ...more, last] = foo();", "']' expected")

	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.parseWarning(
		"var [first, ...rest] = foo();",
		requiresLanguageModeMessage(parser.FeatureSet_Feature_ARRAY_DESTRUCTURING),
		requiresLanguageModeMessage(parser.FeatureSet_Feature_ARRAY_PATTERN_REST))
}

// Java: ParserTest.testObjectDestructuringDeclarationRest()
func TestObjectDestructuringDeclarationRest(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY

	pt.expectFeatures(parser.FeatureSet_Feature_OBJECT_DESTRUCTURING, parser.FeatureSet_Feature_OBJECT_PATTERN_REST)
	pt.parse("var {first, ...rest} = foo();")
	pt.parse("let {first, ...rest} = foo();")
	pt.parse("const {first, ...rest} = foo();")

	pt.expectFeatures()
	pt.parseError(
		"var {first, ...more = 'default'} = foo();",
		"A default value cannot be specified after '...'")
	pt.parseError("var {first, ...more, last} = foo();", "'}' expected")

	pt.mode = Config_LanguageMode_ECMASCRIPT_2015
	pt.parseWarning(
		"var {first, ...rest} = foo();", requiresLanguageModeMessage(parser.FeatureSet_Feature_OBJECT_PATTERN_REST))
}

// Java: ParserTest.testArrayLiteralDeclarationSpread()
func TestArrayLiteralDeclarationSpread(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY

	pt.expectFeatures(parser.FeatureSet_Feature_SPREAD_EXPRESSIONS)
	pt.parse("var o = [first, ...spread];")

	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.parseWarning(
		"var o = [first, ...spread];", requiresLanguageModeMessage(parser.FeatureSet_Feature_SPREAD_EXPRESSIONS))
}

// Java: ParserTest.testObjectLiteralDeclarationSpread()
func TestObjectLiteralDeclarationSpread(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY

	pt.expectFeatures(parser.FeatureSet_Feature_OBJECT_LITERALS_WITH_SPREAD)
	pt.parse("var o = {first: 1, ...spread};")

	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.parseWarning(
		"var o = {first: 1, ...spread};",
		requiresLanguageModeMessage(parser.FeatureSet_Feature_OBJECT_LITERALS_WITH_SPREAD))

	pt.mode = Config_LanguageMode_ECMASCRIPT_2015
	pt.parseWarning(
		"var o = {first: 1, ...spread};",
		requiresLanguageModeMessage(parser.FeatureSet_Feature_OBJECT_LITERALS_WITH_SPREAD))
}

// Java: ParserTest.testArrayDestructuringAssignRest()
func TestArrayDestructuringAssignRest(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY
	pt.expectFeatures(parser.FeatureSet_Feature_ARRAY_DESTRUCTURING, parser.FeatureSet_Feature_ARRAY_PATTERN_REST)
	pt.parse("[first, ...rest] = foo();")
	// nested destructuring in regular parameters and rest parameters
	pt.parse("[first, {a, b}, ...[re, st, ...{length}]] = foo();")
	// arbitrary LHS assignment target is allowed
	pt.parse("[x, ...y[15]] = foo();")

	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.parseWarning(
		"var [first, ...rest] = foo();",
		requiresLanguageModeMessage(parser.FeatureSet_Feature_ARRAY_DESTRUCTURING),
		requiresLanguageModeMessage(parser.FeatureSet_Feature_ARRAY_PATTERN_REST))
}

// Java: ParserTest.testObjectDestructuringAssignRest()
func TestObjectDestructuringAssignRest(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY
	pt.expectFeatures(parser.FeatureSet_Feature_OBJECT_DESTRUCTURING, parser.FeatureSet_Feature_OBJECT_PATTERN_REST)
	pt.parse("const {first, ...rest} = foo();")

	pt.mode = Config_LanguageMode_ECMASCRIPT_2015
	pt.parseWarning(
		"var {first, ...rest} = foo();", requiresLanguageModeMessage(parser.FeatureSet_Feature_OBJECT_PATTERN_REST))
}

// Java: ParserTest.testArrayDestructuringAssignRestInvalid()
func TestArrayDestructuringAssignRestInvalid(t *testing.T) {
	pt := newParserTest(t)

	// arbitrary LHS assignment target not allowed
	pt.parseError(
		"var [x, ...y[15]] = foo();",
		"Only an identifier or destructuring pattern is allowed here.")

	pt.parseError(
		"[first, ...more = 'default'] = foo();", "A default value cannot be specified after '...'")
	pt.parseError("var [first, ...more, last] = foo();", "']' expected")
}

// Java: ParserTest.testArrayDestructuringFnDeclaration()
func TestArrayDestructuringFnDeclaration(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY
	pt.expectFeatures(parser.FeatureSet_Feature_ARRAY_DESTRUCTURING)
	pt.parse("function f([x, y]) { use(x); use(y); }")
	pt.parse("function f([x, [y, z]]) {}")
	pt.parse("function f([x, {y, foo: z}]) {}")
	pt.parse("function f([x, y] = [1, 2]) { use(x); use(y); }")
	pt.parse("function f([x1, x2]) {}")
}

// Java: ParserTest.testArrayDestructuringFnDeclarationInvalid()
func TestArrayDestructuringFnDeclarationInvalid(t *testing.T) {
	pt := newParserTest(t)

	// arbitrary LHS expression not allowed as a formal parameter
	pt.parseError(
		"function f([a[0], x]) {}", "Only an identifier or destructuring pattern is allowed here.")
	// restriction applies to sub-patterns
	pt.parseError(
		"function f([a, [x.foo]]) {}",
		"Only an identifier or destructuring pattern is allowed here.")
	pt.parseError(
		"function f([a, {foo: x.foo}]) {}",
		"Only an identifier or destructuring pattern is allowed here.")
}

// Java: ParserTest.testObjectDestructuringVar()
func TestObjectDestructuringVar(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY
	pt.expectFeatures(parser.FeatureSet_Feature_OBJECT_DESTRUCTURING)
	pt.parse("var {x, y} = foo();")
	pt.parse("var {x: x, y: y} = foo();")
	pt.parse("var {x: {y, z}} = foo();")
	pt.parse("var {x: {y: {z}}} = foo();")

	// Useless, but legal.
	pt.parse("var {} = foo();")
}

// Java: ParserTest.testObjectDestructuringVarInvalid()
func TestObjectDestructuringVarInvalid(t *testing.T) {
	pt := newParserTest(t)

	// Arbitrary LHS target not allowed in declaration
	pt.parseError("var {x.a, y} = foo();", "'}' expected")
	pt.parseError(
		"var {a: x.a, y} = foo();", "Only an identifier or destructuring pattern is allowed here.")
}

// Java: ParserTest.testObjectDestructuringVarWithInitializer()
func TestObjectDestructuringVarWithInitializer(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY
	pt.expectFeatures(parser.FeatureSet_Feature_OBJECT_DESTRUCTURING, parser.FeatureSet_Feature_DEFAULT_PARAMETERS)
	pt.parse("var {x = 1} = foo();")
	pt.parse("var {x: {y = 1}} = foo();")
	pt.parse("var {x: y = 1} = foo();")
	pt.parse("var {x: v1 = 5, y: v2 = 'str'} = foo();")
	pt.parse("var {k1: {k2 : x} = bar(), k3: y} = foo();")
}

// Java: ParserTest.testObjectDestructuringAssign()
func TestObjectDestructuringAssign(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parseError("({x, y}) = foo();", "invalid assignment target")
	pt.expectFeatures(parser.FeatureSet_Feature_OBJECT_DESTRUCTURING)
	pt.parse("({x, y} = foo());")
	pt.parse("({x: x, y: y} = foo());")
	pt.parse("({x: {y, z}} = foo());")
	pt.parse("({k1: {k2 : x} = bar(), k3: y} = foo());")

	// Useless, but legal.
	pt.parse("({} = foo());")
}

// Java: ParserTest.testObjectDestructuringAssignWithInitializer()
func TestObjectDestructuringAssignWithInitializer(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parseError("({x = 1}) = foo();", "invalid assignment target")
	pt.expectFeatures(parser.FeatureSet_Feature_OBJECT_DESTRUCTURING)
	pt.parse("({x = 1} = foo());")
	pt.parse("({x: {y = 1}} = foo());")
	pt.parse("({x: y = 1} = foo());")
	pt.parse("({x: v1 = 5, y: v2 = 'str'} = foo());")
	pt.parse("({k1: {k2 : x} = bar(), k3: y} = foo());")
}

// Java: ParserTest.testObjectDestructuringWithInitializerInvalid()
func TestObjectDestructuringWithInitializerInvalid(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("var {{x}} = foo();", "'}' expected")
	pt.parseError("({{x}}) = foo();", "'}' expected")
	pt.parseError("({{a} = {a: 'b'}}) = foo();", "'}' expected")
	pt.parseError("({{a : b} = {a: 'b'}}) = foo();", "'}' expected")
}

// Java: ParserTest.testObjectDestructuringFnDeclaration()
func TestObjectDestructuringFnDeclaration(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY
	pt.expectFeatures(parser.FeatureSet_Feature_OBJECT_DESTRUCTURING)
	pt.parse("function f({x, y}) { use(x); use(y); }")
	pt.parse("function f({w, x: {y, z}}) {}")
	pt.parse("function f({x, y} = {x:1, y:2}) {}")
	pt.parse("function f({x1, x2}) {}")
}

// Java: ParserTest.testObjectDestructuringFnDeclarationInvalid()
func TestObjectDestructuringFnDeclarationInvalid(t *testing.T) {
	pt := newParserTest(t)

	// arbitrary LHS expression not allowed as a formal parameter
	pt.parseError("function f({a[0], x}) {}", "'}' expected")
	pt.parseError(
		"function f({foo: a[0], x}) {}",
		"Only an identifier or destructuring pattern is allowed here.")
	// restriction applies to sub-patterns
	pt.parseError(
		"function f({a, foo: [x.foo]}) {}",
		"Only an identifier or destructuring pattern is allowed here.")
	pt.parseError(
		"function f({a, x: {foo: x.foo}}) {}",
		"Only an identifier or destructuring pattern is allowed here.")
}

// Java: ParserTest.testObjectDestructuringComputedProp()
func TestObjectDestructuringComputedProp(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY

	pt.parseError("var {[x]} = z;", "':' expected")

	pt.expectFeatures(parser.FeatureSet_Feature_OBJECT_DESTRUCTURING)
	pt.parse("var {[x]: y} = z;")
	pt.parse("var { [foo()] : [x,y,z] = bar() } = baz();")
}

// Java: ParserTest.testObjectDestructuringStringAndNumberKeys()
func TestObjectDestructuringStringAndNumberKeys(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY

	pt.parseError("var { 'hello world' } = foo();", "':' expected")
	pt.parseError("var { 4 } = foo();", "':' expected")
	pt.parseError("var { 'hello' = 'world' } = foo();", "':' expected")
	pt.parseError("var { 2 = 5 } = foo();", "':' expected")

	pt.expectFeatures(parser.FeatureSet_Feature_OBJECT_DESTRUCTURING)
	pt.parse("var {'s': x} = foo();")
	pt.parse("var {3: x} = foo();")
}

// Java: ParserTest.testObjectNumberKeysSpecial()
func TestObjectNumberKeysSpecial(t *testing.T) {
	pt := newParserTest(t)

	n := pt.parse("var a = {12345678901234567890: 2}")

	objectLit := n.GetFirstChild().GetFirstFirstChild()
	assertThat(t, objectLit.GetToken()).IsEqualTo(rhino.Token_OBJECTLIT)

	number := objectLit.GetFirstChild()
	assertThat(t, number.GetToken()).IsEqualTo(rhino.Token_STRING_KEY)
	assertThat(t, number.GetString()).IsEqualTo("12345678901234567000")
}

// Java: ParserTest.testObjectDestructuringKeywordKeys()
func TestObjectDestructuringKeywordKeys(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY
	pt.expectFeatures(parser.FeatureSet_Feature_OBJECT_DESTRUCTURING)
	pt.parse("var {if: x, else: y} = foo();")
	pt.parse("var {while: x=1, for: y} = foo();")
	pt.parse("var {type} = foo();")
	pt.parse("var {declare} = foo();")
	pt.parse("var {module} = foo();")
	pt.parse("var {namespace} = foo();")
}

// Java: ParserTest.testObjectDestructuringKeywordKeysInvalid()
func TestObjectDestructuringKeywordKeysInvalid(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("var {while} = foo();", "cannot use keyword 'while' here.")
	pt.parseError("var {implements} = foo();", "cannot use keyword 'implements' here.")
}

// Java: ParserTest.testObjectDestructuringComplexTarget()
func TestObjectDestructuringComplexTarget(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parseError(
		"var {foo: bar.x} = baz();",
		"Only an identifier or destructuring pattern is allowed here.")

	pt.parseError(
		"var {foo: bar[x]} = baz();",
		"Only an identifier or destructuring pattern is allowed here.")

	pt.expectFeatures(parser.FeatureSet_Feature_OBJECT_DESTRUCTURING)
	pt.parse("({foo: bar.x} = baz());")
	pt.parse("for ({foo: bar.x} in baz());")

	pt.parse("({foo: bar[x]} = baz());")
	pt.parse("for ({foo: bar[x]} in baz());")
}

// Java: ParserTest.testObjectDestructuringExtraParens()
func TestObjectDestructuringExtraParens(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY
	pt.expectFeatures(parser.FeatureSet_Feature_OBJECT_DESTRUCTURING)
	pt.parse("({x: y} = z);")
	pt.parse("({x: (y)} = z);")
	pt.parse("({x: ((y))} = z);")

	pt.expectFeatures(parser.FeatureSet_Feature_ARRAY_DESTRUCTURING)
	pt.parse("([x] = y);")
	pt.parse("[(x), y] = z;")
	pt.parse("[x, (y)] = z;")
}

// Java: ParserTest.testObjectDestructuringExtraParensInvalid()
func TestObjectDestructuringExtraParensInvalid(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("[x, ([y])] = z;", parserTest_INVALID_ASSIGNMENT_TARGET)
	pt.parseError("[x, (([y]))] = z;", parserTest_INVALID_ASSIGNMENT_TARGET)
}

// Java: ParserTest.testObjectLiteralCannotUseDestructuring()
func TestObjectLiteralCannotUseDestructuring(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parseError("var o = {x = 5}", "Default value cannot appear at top level of an object literal.")
}

// Java: ParserTest.testMixedDestructuring()
func TestMixedDestructuring(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY
	pt.expectFeatures(parser.FeatureSet_Feature_ARRAY_DESTRUCTURING, parser.FeatureSet_Feature_OBJECT_DESTRUCTURING)
	pt.parse("var {x: [y, z]} = foo();")
	pt.parse("var [x, {y, z}] = foo();")

	pt.parse("({x: [y, z]} = foo());")
	pt.parse("[x, {y, z}] = foo();")

	pt.parse("function f({x: [y, z]}) {}")
	pt.parse("function f([x, {y, z}]) {}")
}

// Java: ParserTest.testMixedDestructuringWithInitializer()
func TestMixedDestructuringWithInitializer(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY
	pt.expectFeatures(parser.FeatureSet_Feature_ARRAY_DESTRUCTURING, parser.FeatureSet_Feature_OBJECT_DESTRUCTURING)
	pt.parse("var {x: [y, z] = [1, 2]} = foo();")
	pt.parse("var [x, {y, z} = {y: 3, z: 4}] = foo();")

	pt.parse("({x: [y, z] = [1, 2]} = foo());")
	pt.parse("[x, {y, z} = {y: 3, z: 4}] = foo();")

	pt.parse("function f({x: [y, z] = [1, 2]}) {}")
	pt.parse("function f([x, {y, z} = {y: 3, z: 4}]) {}")
}

// Java: ParserTest.testDestructuringNoRHS()
func TestDestructuringNoRHS(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY

	pt.parseError("var {x: y};", "destructuring must have an initializer")
	pt.parseError("let {x: y};", "destructuring must have an initializer")
	pt.parseError("const {x: y};", "const variables must have an initializer")
	pt.parseError("var {x};", "destructuring must have an initializer")
	pt.parseError("let {x};", "destructuring must have an initializer")
	pt.parseError("const {x};", "const variables must have an initializer")
	pt.parseError("var [x, y];", "destructuring must have an initializer")
	pt.parseError("let [x, y];", "destructuring must have an initializer")
	pt.parseError("const [x, y];", "const variables must have an initializer")
}

// Java: ParserTest.testComprehensions()
func TestComprehensions(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY
	error := "unsupported language feature:" + " array/generator comprehensions"

	// array comprehensions
	pt.parseError("[for (x of y) z];", error)
	pt.parseError("[for ({x,y} of z) x+y];", error)
	pt.parseError("[for (x of y) if (x<10) z];", error)
	pt.parseError("[for (a = 5 of v) a];", "'identifier' expected")

	// generator comprehensions
	pt.parseError("(for (x of y) z);", error)
	pt.parseError("(for ({x,y} of z) x+y);", error)
	pt.parseError("(for (x of y) if (x<10) z);", error)
	pt.parseError("(for (a = 5 of v) a);", "'identifier' expected")
}

// Java: ParserTest.testLetForbidden1()
func TestLetForbidden1(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.expectFeatures(parser.FeatureSet_Feature_LET_DECLARATIONS)
	pt.parseWarning("let x = 3;", requiresLanguageModeMessage(parser.FeatureSet_Feature_LET_DECLARATIONS))
}

// Java: ParserTest.testLetForbidden2()
func TestLetForbidden2(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.expectFeatures(parser.FeatureSet_Feature_LET_DECLARATIONS)
	pt.parseWarning(
		"function f() { let x = 3; };", requiresLanguageModeMessage(parser.FeatureSet_Feature_LET_DECLARATIONS))
}

// Java: ParserTest.testBlockScopedFunctionDeclaration()
func TestBlockScopedFunctionDeclaration(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_BLOCK_SCOPED_FUNCTION_DECLARATION)
	pt.parse("{ function foo() {} }")
	pt.parse("if (true) { function foo() {} }")
	pt.parse("{ function* gen() {} }")
	pt.parse("if (true) function foo() {}")
	pt.parse("if (true) function foo() {} else {}")
	pt.parse("if (true) {} else function foo() {}")
	pt.parse("if (true) function foo() {} else function foo() {}")

	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.expectFeatures()
	// Function expressions and functions directly inside other functions do not trigger this
	pt.parse("function foo() {}")
	pt.parse("(function foo() {})")
	pt.parse("function foo() { function bar() {} }")
	pt.parse("{ var foo = function() {}; }")
	pt.parse("{ var foo = function bar() {}; }")
	pt.parse("{ (function() {})(); }")
	pt.parse("{ (function foo() {})(); }")

	pt.parseWarning(
		"{ function f() {} }",
		requiresLanguageModeMessage(parser.FeatureSet_Feature_BLOCK_SCOPED_FUNCTION_DECLARATION))
}

// Java: ParserTest.testLetForbidden3()
func TestLetForbidden3(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.parseError("function f() { var let = 3; }", "'identifier' expected")

	pt.mode = Config_LanguageMode_ECMASCRIPT_2015
	pt.parseError("function f() { var let = 3; }", "'identifier' expected")
}

// Java: ParserTest.testYieldForbidden()
func TestYieldForbidden(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("function f() { yield 3; }", "primary expression expected")
	pt.parseError("function f(x = yield 3) { return x }", "primary expression expected")
	pt.parseError(
		"function *f(x = yield 3) { return x }", "`yield` is illegal in parameter default value.")
	pt.parseError(
		"function *f() { return function*(x = yield 1) { return x; } }",
		"`yield` is illegal in parameter default value.")
}

// Java: ParserTest.testGenerator()
func TestGenerator(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_GENERATORS)
	pt.parse("var obj = { *f() { yield 3; } };")
	pt.parse("function* f() { yield 3; }")
	pt.parse("function f() { return function* g() {} }")

	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.parseWarning("function* f() { yield 3; }", requiresLanguageModeMessage(parser.FeatureSet_Feature_GENERATORS))
	pt.parseWarning(
		"var obj = { * f() { yield 3; } };",
		requiresLanguageModeMessage(parser.FeatureSet_Feature_GENERATORS),
		requiresLanguageModeMessage(parser.FeatureSet_Feature_MEMBER_DECLARATIONS))
}

// Java: ParserTest.testBracelessFunctionForbidden()
func TestBracelessFunctionForbidden(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("var sq = function(x) x * x;", "'{' expected")
}

// Java: ParserTest.testGeneratorsForbidden()
func TestGeneratorsForbidden(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("var i = (x for (x in obj));", "')' expected")
}

// Java: ParserTest.testGettersForbidden1()
func TestGettersForbidden1(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT3
	pt.expectFeatures(parser.FeatureSet_Feature_GETTER)
	pt.parseError("var x = {get foo() { return 3; }};", IRFactory_GETTER_ERROR_MESSAGE)
}

// Java: ParserTest.testGettersForbidden2()
func TestGettersForbidden2(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT3
	pt.parseError("var x = {get foo bar() { return 3; }};", "'(' expected")
}

// Java: ParserTest.testGettersForbidden3()
func TestGettersForbidden3(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT3
	pt.parseError("var x = {a getter:function b() { return 3; }};", "'}' expected")
}

// Java: ParserTest.testGettersForbidden4()
func TestGettersForbidden4(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT3
	pt.parseError("var x = {\"a\" getter:function b() { return 3; }};", "':' expected")
}

// Java: ParserTest.testGettersForbidden5()
func TestGettersForbidden5(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT3
	pt.expectFeatures(parser.FeatureSet_Feature_GETTER)
	pt.parseError("var x = {a: 2, get foo() { return 3; }};", IRFactory_GETTER_ERROR_MESSAGE)
}

// Java: ParserTest.testGettersForbidden6()
func TestGettersForbidden6(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT3
	pt.expectFeatures(parser.FeatureSet_Feature_GETTER)
	pt.parseError("var x = {get 'foo'() { return 3; }};", IRFactory_GETTER_ERROR_MESSAGE)
}

// Java: ParserTest.testSettersForbidden()
func TestSettersForbidden(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT3
	pt.expectFeatures(parser.FeatureSet_Feature_SETTER)
	pt.parseError("var x = {set foo(a) { y = 3; }};", IRFactory_SETTER_ERROR_MESSAGE)
}

// Java: ParserTest.testSettersForbidden2()
func TestSettersForbidden2(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT3
	// TODO(johnlenz): maybe just report the first error, when not in IDE mode?
	pt.parseError("var x = {a setter:function b() { return 3; }};", "'}' expected")
}

// Java: ParserTest.testFileOverviewJSDoc1()
func TestFileOverviewJSDoc1(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true

	n := pt.parse("/** @fileoverview Hi mom! */ function Foo() {}")
	assertNode(t, n.GetFirstChild()).HasType(rhino.Token_FUNCTION)
	assertThat(t, n.GetJSDocInfo()).IsNotNull()
	assertThat(t, n.GetFirstChild().GetJSDocInfo()).IsNull()
	assertThat(t, n.GetJSDocInfo().GetFileOverview()).IsEqualTo("Hi mom!")
}

// Java: ParserTest.testFileOverviewJSDoc_notOnTopOfFile()
func TestFileOverviewJSDoc_notOnTopOfFile(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true

	n := pt.parse("// some comment \n let x; /** @fileoverview Hi mom! */ class Foo {}")
	assertNode(t, n).HasType(rhino.Token_SCRIPT)
	assertThat(t, n.GetJSDocInfo()).IsNotNull()
	assertThat(t, n.GetJSDocInfo().GetFileOverview()).IsEqualTo("Hi mom!")

	letNode := n.GetFirstChild()
	assertNode(t, letNode).HasType(rhino.Token_LET)
	assertThat(t, letNode.GetJSDocInfo()).IsNull()

	classNode := n.GetSecondChild()
	assertNode(t, classNode).HasType(rhino.Token_CLASS)
	assertThat(t, classNode.GetJSDocInfo()).IsNull()
}

// Java: ParserTest.testFileOverviewJSDocDoesNotHoseParsing()
func TestFileOverviewJSDocDoesNotHoseParsing(t *testing.T) {
	pt := newParserTest(t)

	assertNode(t, pt.parse("/** @fileoverview Hi mom! \n */ function Foo() {}").GetFirstChild()).
		HasType(rhino.Token_FUNCTION)
	assertNode(t, pt.parse("/** @fileoverview Hi mom! \n * * * */ function Foo() {}").GetFirstChild()).
		HasType(rhino.Token_FUNCTION)
	assertNode(t, pt.parse("/** @fileoverview \n * x */ function Foo() {}").GetFirstChild()).
		HasType(rhino.Token_FUNCTION)
	assertNode(t, pt.parse("/** @fileoverview \n * x \n */ function Foo() {}").GetFirstChild()).
		HasType(rhino.Token_FUNCTION)
}

// Java: ParserTest.testFileOverviewJSDoc2()
func TestFileOverviewJSDoc2(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true

	n :=
		pt.parse(
			`/** @fileoverview Hi mom! */
 /** @constructor */ function Foo() {}
`)
	assertThat(t, n.GetJSDocInfo()).IsNotNull()
	assertThat(t, n.GetJSDocInfo().GetFileOverview()).IsEqualTo("Hi mom!")
	assertThat(t, n.GetFirstChild().GetJSDocInfo()).IsNotNull()
	assertThat(t, n.GetFirstChild().GetJSDocInfo().HasFileOverview()).IsFalse()
	assertThat(t, n.GetFirstChild().GetJSDocInfo().IsConstructor()).IsTrue()
}

// Java: ParserTest.testFileoverview_firstOneWins()
func TestFileoverview_firstOneWins(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true

	n :=
		pt.parse(
			`/** @fileoverview First */
/** @fileoverview Second */
`)

	assertThat(t, n.GetJSDocInfo()).IsNotNull()
	assertThat(t, n.GetJSDocInfo().GetFileOverview()).IsEqualTo("First")
}

// Java: ParserTest.testFileoverview_firstOneWins_implicitFileoverview()
func TestFileoverview_firstOneWins_implicitFileoverview(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true

	n :=
		pt.parse(
			`/** @typeSummary */
/** @fileoverview Second */
`)

	assertThat(t, n.GetJSDocInfo()).IsNotNull()
	assertThat(t, n.GetJSDocInfo().IsTypeSummary()).IsTrue()
	assertThat(t, n.GetJSDocInfo().GetFileOverview()).IsNull()
}

// Java: ParserTest.testFileoverview_firstOneWins_suppressionsAccumulate()
func TestFileoverview_firstOneWins_suppressionsAccumulate(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true

	n :=
		pt.parse(
			`/** @fileoverview @suppress {const} */
/** @fileoverview @suppress {checkTypes} */
`)

	assertThat(t, n.GetJSDocInfo()).IsNotNull()
	assertThat(t, n.GetJSDocInfo().GetSuppressions()).ContainsExactly("const", "checkTypes")
}

// Java: ParserTest.testFileoverview_firstOneWins_externsAccumulate()
func TestFileoverview_firstOneWins_externsAccumulate(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true

	n :=
		pt.parse(
			`/** @fileoverview First */
/** @externs */
`)

	assertThat(t, n.GetJSDocInfo()).IsNotNull()
	assertThat(t, n.GetJSDocInfo().GetFileOverview()).IsEqualTo("First")
	assertThat(t, n.GetJSDocInfo().IsExterns()).IsTrue()
}

// Java: ParserTest.testFileoverview_typeSummaryAccumulates_withExterns()
func TestFileoverview_typeSummaryAccumulates_withExterns(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true

	n :=
		pt.parse(
			`/** @fileoverview @typeSummary */
/** @externs */
;
`)

	assertThat(t, n.GetJSDocInfo()).IsNotNull()
	assertThat(t, n.GetJSDocInfo().IsExterns()).IsTrue()
	assertThat(t, n.GetJSDocInfo().IsTypeSummary()).IsTrue()
}

// Java: ParserTest.testImportantComment()
func TestImportantComment(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true

	n := pt.parse("/*! Hi mom! */ function Foo() {}")
	assertNode(t, n.GetFirstChild()).HasType(rhino.Token_FUNCTION)
	assertThat(t, n.GetJSDocInfo()).IsNotNull()
	assertThat(t, n.GetFirstChild().GetJSDocInfo()).IsNull()
	assertThat(t, n.GetJSDocInfo().GetLicense()).IsEqualTo(" Hi mom! ")
}

// Java: ParserTest.testBlockCommentParsed()
func TestBlockCommentParsed(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	n := pt.parse("/* Hi mom! */ \n function Foo() {}")
	assertNode(t, n.GetFirstChild()).HasType(rhino.Token_FUNCTION)
	assertThat(t, n.GetFirstChild().GetNonJSDocCommentString()).IsEqualTo("/* Hi mom! */")
	assertThat(t, n.GetFirstChild().GetNonJSDocComment().IsInline()).IsFalse()
}

// Java: ParserTest.testBlockCommentNotParsedWithoutParsingMode()
func TestBlockCommentNotParsedWithoutParsingMode(t *testing.T) {
	pt := newParserTest(t)

	n := pt.parse("/* Hi mom! */ \n function Foo() {}")
	assertNode(t, n.GetFirstChild()).HasType(rhino.Token_FUNCTION)
	assertThat(t, n.GetFirstChild().GetNonJSDocComment()).IsNull()
}

// Java: ParserTest.testLineCommentParsed()
func TestLineCommentParsed(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	n := pt.parse("// Hi mom! \n function Foo() {}")
	assertNode(t, n.GetFirstChild()).HasType(rhino.Token_FUNCTION)
	assertThat(t, n.GetFirstChild().GetNonJSDocCommentString()).IsEqualTo("// Hi mom!")
}

// Java: ParserTest.testManyIndividualCommentsGetAttached()
func TestManyIndividualCommentsGetAttached(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	n := pt.parse("// Hi Mom! \n function Foo() {} \n // Hi Dad! \n  function Bar() {}")
	assertNode(t, n.GetFirstChild()).HasType(rhino.Token_FUNCTION)
	assertThat(t, n.GetFirstChild().GetNonJSDocCommentString()).IsEqualTo("// Hi Mom!")
	assertNode(t, n.GetSecondChild()).HasType(rhino.Token_FUNCTION)
	assertThat(t, n.GetSecondChild().GetNonJSDocCommentString()).IsEqualTo("// Hi Dad!")
}

// Java: ParserTest.testIndividualCommentsAroundClasses()
func TestIndividualCommentsAroundClasses(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	n := pt.parse("// comment A \n class A{} // trailing a \n // comment B \n  class Bar{}")
	classA := n.GetFirstChild()
	classB := n.GetSecondChild()
	assertNode(t, classA).HasType(rhino.Token_CLASS)
	assertThat(t, classA.GetNonJSDocCommentString()).IsEqualTo("// comment A")
	assertThat(t, classA.GetTrailingNonJSDocCommentString()).IsEqualTo("// trailing a")
	assertNode(t, classB).HasType(rhino.Token_CLASS)
	assertThat(t, classB.GetNonJSDocCommentString()).IsEqualTo("// comment B")
}

// Java: ParserTest.testMultipleLinedCommentsAttachedToSameNode()
func TestMultipleLinedCommentsAttachedToSameNode(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	n := pt.parse("// Hi Mom! \n // And Dad! \n  function Foo() {} ")
	assertNode(t, n.GetFirstChild()).HasType(rhino.Token_FUNCTION)
	assertThat(t, n.GetFirstChild().GetNonJSDocCommentString()).IsEqualTo("// Hi Mom!\n// And Dad!")
}

// Java: ParserTest.testEmptyLinesBetweenCommentsIsPreserved()
func TestEmptyLinesBetweenCommentsIsPreserved(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	n := pt.parse("/* Hi Mom!*/ \n\n\n // And Dad! \n  function Foo() {} ")
	assertNode(t, n.GetFirstChild()).HasType(rhino.Token_FUNCTION)
	assertThat(t, n.GetFirstChild().GetNonJSDocCommentString()).
		IsEqualTo(`/* Hi Mom!*/


// And Dad!`)
}

// Java: ParserTest.testObjectLiteralDoc1()
func TestObjectLiteralDoc1(t *testing.T) {
	pt := newParserTest(t)

	n := pt.parse("var x = {/** @type {number} */ 1: 2};")

	objectLit := n.GetFirstFirstChild().GetFirstChild()
	assertNode(t, objectLit).HasType(rhino.Token_OBJECTLIT)

	number := objectLit.GetFirstChild()
	assertNode(t, number).HasType(rhino.Token_STRING_KEY)
	assertThat(t, number.GetJSDocInfo()).IsNotNull()
}

// Java: ParserTest.testDuplicatedParam()
func TestDuplicatedParam(t *testing.T) {
	pt := newParserTest(t)

	pt.parseWarning("function foo(x, x) {}", "Duplicate parameter name \"x\"")
	pt.parseWarning("function foo(x, n, x) {}", "Duplicate parameter name \"x\"")
	pt.parseWarning("function foo(n, x, x) {}", "Duplicate parameter name \"x\"")
	pt.parseWarning("function foo(x, {x}) {}", "Duplicate parameter name \"x\"")
	pt.parseWarning("function foo(x, {n: {x}}) {}", "Duplicate parameter name \"x\"")
	pt.parseWarning("function foo(x, [x]) {}", "Duplicate parameter name \"x\"")
	pt.parseWarning("function foo(x, ...x) {}", "Duplicate parameter name \"x\"")
	pt.parseWarning("function foo(...[x, x]) {}", "Duplicate parameter name \"x\"")
	pt.parseWarning("function foo(...[x,,x]) {}", "Duplicate parameter name \"x\"")
	pt.parseWarning("function foo(x, x = 1) {}", "Duplicate parameter name \"x\"")
	pt.parseWarning("function foo([x, x] = [1, 2]) {}", "Duplicate parameter name \"x\"")
	pt.parseWarning("function foo({x, x} = {x: 1}) {}", "Duplicate parameter name \"x\"")
	pt.parseWarning("function foo(x, {[n()]: x}) {}", "Duplicate parameter name \"x\"")
	pt.parseWarning("function foo(x = x) {}")
	pt.parseWarning("function foo({x: x}) {}")
	pt.parseWarning("function foo(foo) {}")
}

// Java: ParserTest.testLetAsIdentifier()
func TestLetAsIdentifier(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT3
	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parse("var let")

	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.parse("var let")

	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.strictMode = Config_StrictMode_STRICT
	pt.parseError("var let", "'identifier' expected")

	pt.mode = Config_LanguageMode_ECMASCRIPT_2015
	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parse("var let")

	pt.mode = Config_LanguageMode_ECMASCRIPT_2015
	pt.strictMode = Config_StrictMode_STRICT
	pt.parseError("var let", "'identifier' expected")
}

// Java: ParserTest.testLet()
func TestLet(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY
	pt.expectFeatures(parser.FeatureSet_Feature_LET_DECLARATIONS)

	pt.parse("let x;")
	pt.parse("let x = 1;")
	pt.parse("let x, y = 2;")
	pt.parse("let x = 1, y = 2;")
}

// Java: ParserTest.testConst()
func TestConst(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY

	pt.parseError("const x;", "const variables must have an initializer")
	pt.parseError("const x, y = 2;", "const variables must have an initializer")

	pt.expectFeatures(parser.FeatureSet_Feature_CONST_DECLARATIONS)

	pt.parse("const x = 1;")
	pt.parse("const x = 1, y = 2;")
}

// Java: ParserTest.testYield1()
func TestYield1(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT3
	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parse("var yield")

	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.parse("var yield")

	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.strictMode = Config_StrictMode_STRICT
	pt.parseError("var yield", "'identifier' expected")

	pt.mode = Config_LanguageMode_ECMASCRIPT_2015
	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parse("var yield")

	pt.mode = Config_LanguageMode_ECMASCRIPT_2015
	pt.strictMode = Config_StrictMode_STRICT
	pt.parseError("var yield", "'identifier' expected")
}

// Java: ParserTest.testYield2()
func TestYield2(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_GENERATORS)
	pt.parse("function * f() { yield; }")
	pt.parse("function * f() { yield /a/i; }")

	pt.expectFeatures()
	pt.parseError("function * f() { 1 + yield; }", "primary expression expected")
	pt.parseError("function * f() { 1 + yield 2; }", "primary expression expected")
	pt.parseError("function * f() { yield 1 + yield 2; }", "primary expression expected")
	pt.parseError("function * f() { yield(1) + yield(2); }", "primary expression expected")
	pt.expectFeatures(parser.FeatureSet_Feature_GENERATORS)
	pt.parse("function * f() { (yield 1) + (yield 2); }") // OK
	pt.parse("function * f() { yield * yield; }")         // OK  (yield * (yield))
	pt.expectFeatures()
	pt.parseError("function * f() { yield + yield; }", "primary expression expected")
	pt.expectFeatures(parser.FeatureSet_Feature_GENERATORS)
	pt.parse("function * f() { (yield) + (yield); }") // OK
	pt.parse("function * f() { return yield; }")      // OK
	pt.parse("function * f() { return yield 1; }")    // OK
	pt.parse(
		`function * f() {
  yield * // line break allowed here
      [1, 2, 3];
}
`)
	pt.expectFeatures()
	pt.parseError(
		`function * f() {
  yield // line break not allowed here
      *[1, 2, 3];
}
`,
		"'}' expected")
	pt.parseError("function * f() { yield *; }", "yield* requires an expression")
}

// Java: ParserTest.testYield3()
func TestYield3(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_GENERATORS)
	// TODO(johnlenz): validate "yield" parsing. Firefox rejects this
	// use of "yield".
	pt.parse("function * f() { yield , yield; }")
}

// Java: ParserTest.testStringLineContinuationWarningsByMode()
func TestStringLineContinuationWarningsByMode(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_STRING_CONTINUATION)
	pt.strictMode = Config_StrictMode_SLOPPY

	pt.mode = Config_LanguageMode_ECMASCRIPT3
	pt.parseWarning(
		"'one\\\ntwo';",
		requiresLanguageModeMessage(parser.FeatureSet_Feature_STRING_CONTINUATION),
		parserTest_STRING_CONTINUATIONS_WARNING)

	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.parseWarning("'one\\\ntwo';", parserTest_STRING_CONTINUATIONS_WARNING)

	pt.mode = Config_LanguageMode_ECMASCRIPT_2015
	pt.parseWarning("'one\\\ntwo';", parserTest_STRING_CONTINUATIONS_WARNING)
}

// Java: ParserTest.testStringLineContinuationNormalization()
func TestStringLineContinuationNormalization(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_STRING_CONTINUATION)
	pt.strictMode = Config_StrictMode_SLOPPY

	n := pt.parseWarning("'one\\\ntwo';", parserTest_STRING_CONTINUATIONS_WARNING)
	assertThat(t, n.GetFirstFirstChild().GetString()).IsEqualTo("onetwo")

	n = pt.parseWarning("'one\\\rtwo';", parserTest_STRING_CONTINUATIONS_WARNING)
	assertThat(t, n.GetFirstFirstChild().GetString()).IsEqualTo("onetwo")

	n = pt.parseWarning("'one\\\r\ntwo';", parserTest_STRING_CONTINUATIONS_WARNING)
	assertThat(t, n.GetFirstFirstChild().GetString()).IsEqualTo("onetwo")

	n = pt.parseWarning("'one \\\ntwo';", parserTest_STRING_CONTINUATIONS_WARNING)
	assertThat(t, n.GetFirstFirstChild().GetString()).IsEqualTo("one two")

	n = pt.parseWarning("'one\\\n two';", parserTest_STRING_CONTINUATIONS_WARNING)
	assertThat(t, n.GetFirstFirstChild().GetString()).IsEqualTo("one two")
}

// Java: ParserTest.testStringContinuationIssue3492()
func TestStringContinuationIssue3492(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_STRING_CONTINUATION)
	pt.strictMode = Config_StrictMode_SLOPPY

	// This test runs on this input code, which technically has STRING_CONTINUATIONS_WARNINGs at 2
	// places:
	// ```
	// function x() {
	//    a = "\ <--- STRING_CONTINUATIONS_WARNING
	//    \ \ <--- STRING_CONTINUATIONS_WARNING
	//    ";
	// };
	// ```
	// We deduplicate based on the key `sourceName:line_no:col_no:warning` being the same. Here,
	// even though the `STRING_CONTINUATIONS_WARNING` happens at 2 places, JSCompiler reports the
	// warnings within a string with a location pointing to the start of the string token. This
	// happens here `charno(token.location.start)`:
	// google3/third_party/java_src/jscomp/java/com/google/javascript/jscomp/parsing/IRFactory.java?rcl=719034735&l=3482
	// Hence, after deduplication, we only get one error.
	pt.parseWarning(
		`function x() {
        a = "\
        \ \
        ";
};
`,
		"Unnecessary escape: '\\ ' is equivalent to just ' '",
		parserTest_STRING_CONTINUATIONS_WARNING)
}

// Java: ParserTest.testStringLiteral()
func TestStringLiteral(t *testing.T) {
	pt := newParserTest(t)

	n := pt.parse("'foo'")
	stringNode := n.GetFirstFirstChild()
	assertNode(t, stringNode).HasType(rhino.Token_STRINGLIT)
	assertThat(t, stringNode.GetString()).IsEqualTo("foo")
}

// Java: ParserTest.testUseTemplateLiteral()
func TestUseTemplateLiteral(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_TEMPLATE_LITERALS)
	pt.testTemplateLiteral("f`hello world`;")
	pt.testTemplateLiteral("`hello ${name} ${world}`.length;")
}

// Java: ParserTest.testTemplateLiterals()
func TestTemplateLiterals(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_TEMPLATE_LITERALS)
	pt.testTemplateLiteral("``")
	pt.testTemplateLiteral("`\"`")
	pt.testTemplateLiteral("`\\``")
	pt.testTemplateLiteral("`hello world`;")
	pt.testTemplateLiteral("`hello\nworld`;")
	pt.testTemplateLiteral("`string containing \\`escaped\\` backticks`;")
	pt.testTemplateLiteral("{ `in block` }")
	pt.testTemplateLiteral("{ `in ${block}` }")
}

// Java: ParserTest.testEscapedTemplateLiteral()
func TestEscapedTemplateLiteral(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_TEMPLATE_LITERALS)
	pt.assertSimpleTemplateLiteral("${escaped}", "`\\${escaped}`")
}

// Java: ParserTest.testTemplateLiteralWithNulChar()
func TestTemplateLiteralWithNulChar(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_TEMPLATE_LITERALS)
	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parse("var test = `\nhello\\0`")
}

// Java: ParserTest.testTemplateLiteralWithNewline()
func TestTemplateLiteralWithNewline(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_TEMPLATE_LITERALS)
	pt.assertSimpleTemplateLiteral("hello\nworld", "`hello\nworld`")
	pt.assertSimpleTemplateLiteral("\n", "`\r`")
	pt.assertSimpleTemplateLiteral("\n", "`\r\n`")
	pt.assertSimpleTemplateLiteral("\\\n", "`\\\\\n`")
	pt.assertSimpleTemplateLiteral("\\\n", "`\\\\\r\n`")
	pt.assertSimpleTemplateLiteral("\r\n", "`\\r\\n`")       // template literals support explicit escapes
	pt.assertSimpleTemplateLiteral("\\r\\n", "`\\\\r\\\\n`") // note: no actual newlines here
}

// Java: ParserTest.testTemplateLiteralWithLineContinuation()
func TestTemplateLiteralWithLineContinuation(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY
	pt.expectFeatures(parser.FeatureSet_Feature_TEMPLATE_LITERALS)
	n := pt.parseWarning("`string \\\ncontinuation`", parserTest_STRING_CONTINUATIONS_WARNING)
	templateLiteral := n.GetFirstFirstChild()
	stringNode := templateLiteral.GetFirstChild()
	assertNode(t, stringNode).HasType(rhino.Token_TEMPLATELIT_STRING)
	assertThat(t, stringNode.GetCookedString()).IsEqualTo("string continuation")
}

// Java: ParserTest.testTemplateLiteralSubstitution()
func TestTemplateLiteralSubstitution(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY
	pt.expectFeatures(parser.FeatureSet_Feature_TEMPLATE_LITERALS)
	pt.parse("`hello ${name}`;")
	pt.parse("`hello ${name} ${world}`;")
	pt.parse("`hello ${name }`")

	pt.expectFeatures()
	pt.parseError("`hello ${name", "Expected '}' after expression in template literal")
	pt.parseError("`hello ${name tail}", "Expected '}' after expression in template literal")
}

// Java: ParserTest.testUnterminatedTemplateLiteral()
func TestUnterminatedTemplateLiteral(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parseError("`hello", "Unterminated template literal")
	pt.parseError("`hello\\`", "Unterminated template literal")
}

// Java: ParserTest.testTemplateLiteralOctalEscapes()
func TestTemplateLiteralOctalEscapes(t *testing.T) {
	pt := newParserTest(t)

	pt.assertSimpleTemplateLiteral("\u0000", "`\\0`")
	pt.assertSimpleTemplateLiteral("aaa\u0000aaa", "`aaa\\0aaa`")
}

// Java: ParserTest.testIncorrectEscapeSequenceInTemplateLiteral()
func TestIncorrectEscapeSequenceInTemplateLiteral(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("`hello\\x`", "Hex digit expected")

	pt.parseError("`hello\\1`", "Invalid escape sequence")
	pt.parseError("`hello\\2`", "Invalid escape sequence")
	pt.parseError("`hello\\3`", "Invalid escape sequence")
	pt.parseError("`hello\\4`", "Invalid escape sequence")
	pt.parseError("`hello\\5`", "Invalid escape sequence")
	pt.parseError("`hello\\6`", "Invalid escape sequence")
	pt.parseError("`hello\\7`", "Invalid escape sequence")
	pt.parseError("`hello\\8`", "Invalid escape sequence")
	pt.parseError("`hello\\9`", "Invalid escape sequence")
	pt.parseError("`hello\\00`", "Invalid escape sequence")
	pt.parseError("`hello\\01`", "Invalid escape sequence")
	pt.parseError("`hello\\02`", "Invalid escape sequence")
	pt.parseError("`hello\\03`", "Invalid escape sequence")
	pt.parseError("`hello\\04`", "Invalid escape sequence")
	pt.parseError("`hello\\05`", "Invalid escape sequence")
	pt.parseError("`hello\\06`", "Invalid escape sequence")
	pt.parseError("`hello\\07`", "Invalid escape sequence")
	pt.parseError("`hello\\08`", "Invalid escape sequence")
	pt.parseError("`hello\\09`", "Invalid escape sequence")

	// newline before invalid escape sequence
	pt.parseError("`\n\\1`", "Invalid escape sequence")
	pt.parseError("`\n\\1 ${0}`", "Invalid escape sequence")
}

// Java: ParserTest.testTemplateLiteralSubstitutionWithCast()
func TestTemplateLiteralSubstitutionWithCast(t *testing.T) {
	pt := newParserTest(t)

	root := pt.parse("`${ /** @type {?} */ (3)}`")
	exprResult := root.GetFirstChild()
	templateLiteral := exprResult.GetFirstChild()
	assertNode(t, templateLiteral).HasType(rhino.Token_TEMPLATELIT)

	substitution := templateLiteral.GetSecondChild()
	assertNode(t, substitution).HasType(rhino.Token_TEMPLATELIT_SUB)

	cast := substitution.GetFirstChild()
	assertNode(t, cast).HasType(rhino.Token_CAST)

	number := cast.GetFirstChild()
	assertNode(t, number).HasType(rhino.Token_NUMBER)
}

// Java: ParserTest.testExponentialLiterals()
func TestExponentialLiterals(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("0e0")
	pt.parse("0E0")
	pt.parse("0E1")
	pt.parse("1E0")
	pt.parse("1E-0")
	pt.parse("10E10")
	pt.parse("10E-10")
	pt.parse("1.0E1")
	pt.parseError("01E0", parserTest_SEMICOLON_EXPECTED)
	pt.parseError("0E", "Exponent part must contain at least one digit")
	pt.parseError("1E-", "Exponent part must contain at least one digit")
	pt.parseError("1E1.1", parserTest_SEMICOLON_EXPECTED)
}

// Java: ParserTest.testBigIntLiteralZero()
func TestBigIntLiteralZero(t *testing.T) {
	pt := newParserTest(t)

	bigint :=
		pt.parse("0n;"). // SCRIPT
					GetOnlyChild(). // EXPR_RESULT
					GetOnlyChild()  // BIGINT
	assertNode(t, bigint).HasType(rhino.Token_BIGINT)
	assertNode(t, bigint).HasLineno(1)
	assertNode(t, bigint).HasCharno(0)
	assertNode(t, bigint).HasLength(2)
	assertNode(t, bigint).IsBigInt(jutil.BigInteger_ZERO)
}

// Java: ParserTest.testBigIntLiteralPositive()
func TestBigIntLiteralPositive(t *testing.T) {
	pt := newParserTest(t)

	bigint :=
		pt.parse("1n;"). // SCRIPT
					GetOnlyChild(). // EXPR_RESULT
					GetOnlyChild()  // BIGINT
	assertNode(t, bigint).HasType(rhino.Token_BIGINT)
	assertNode(t, bigint).HasLineno(1)
	assertNode(t, bigint).HasCharno(0)
	assertNode(t, bigint).HasLength(2)
	assertNode(t, bigint).IsBigInt(jutil.BigInteger_ONE)
}

// Java: ParserTest.testBigIntLiteralNegative()
func TestBigIntLiteralNegative(t *testing.T) {
	pt := newParserTest(t)

	neg :=
		pt.parse("-1n;"). // SCRIPT
					GetOnlyChild(). // EXPR_RESULT
					GetOnlyChild()  // NEG

	assertNode(t, neg).HasType(rhino.Token_NEG)
	assertNode(t, neg).HasLineno(1)
	assertNode(t, neg).HasCharno(0)
	assertNode(t, neg).HasLength(3)

	bigint := neg.GetOnlyChild()
	assertNode(t, bigint).HasType(rhino.Token_BIGINT)
	assertNode(t, bigint).HasLineno(1)
	assertNode(t, bigint).HasCharno(1)
	assertNode(t, bigint).HasLength(2)
	assertNode(t, bigint).IsBigInt(jutil.BigInteger_ONE)
}

// Java: ParserTest.testBigIntLiteralBinary()
func TestBigIntLiteralBinary(t *testing.T) {
	pt := newParserTest(t)

	bigint :=
		pt.parse("0b10000n;"). // SCRIPT
					GetOnlyChild(). // EXPR_RESULT
					GetOnlyChild()  // BIGINT
	assertNode(t, bigint).HasType(rhino.Token_BIGINT)
	assertNode(t, bigint).HasLineno(1)
	assertNode(t, bigint).HasCharno(0)
	assertNode(t, bigint).HasLength(8)
	assertNode(t, bigint).IsBigInt(jutil.NewBigIntegerFromString("16"))
}

// Java: ParserTest.testBigIntLiteralOctal()
func TestBigIntLiteralOctal(t *testing.T) {
	pt := newParserTest(t)

	bigint :=
		pt.parse("0o100n;"). // SCRIPT
					GetOnlyChild(). // EXPR_RESULT
					GetOnlyChild()  // BIGINT
	assertNode(t, bigint).HasType(rhino.Token_BIGINT)
	assertNode(t, bigint).HasLineno(1)
	assertNode(t, bigint).HasCharno(0)
	assertNode(t, bigint).HasLength(6)
	assertNode(t, bigint).IsBigInt(jutil.NewBigIntegerFromString("64"))
}

// Java: ParserTest.testBigIntLiteralHex()
func TestBigIntLiteralHex(t *testing.T) {
	pt := newParserTest(t)

	bigint :=
		pt.parse("0xFn;"). // SCRIPT
					GetOnlyChild(). // EXPR_RESULT
					GetOnlyChild()  // BIGINT
	assertNode(t, bigint).HasType(rhino.Token_BIGINT)
	assertNode(t, bigint).HasLineno(1)
	assertNode(t, bigint).HasCharno(0)
	assertNode(t, bigint).HasLength(4)
	assertNode(t, bigint).IsBigInt(jutil.NewBigIntegerFromString("15"))
}

// Java: ParserTest.testBigIntInFunctionStatement()
func TestBigIntInFunctionStatement(t *testing.T) {
	pt := newParserTest(t)

	add :=
		pt.parse("function f(/** @type {bigint} */ x) { 0n + x }"). // SCRIPT
										GetOnlyChild(). // FUNCTION
										GetLastChild(). // BLOCK
										GetOnlyChild(). // EXPR_RESULT
										GetOnlyChild()  // ADD
	assertNode(t, add).HasToken(rhino.Token_ADD)
	assertNode(t, add.GetFirstChild()).IsBigInt(jutil.BigInteger_ZERO)
	assertNode(t, add.GetLastChild()).IsName("x")
}

// Java: ParserTest.testBigIntLiteralErrors()
func TestBigIntLiteralErrors(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("01n;", "SyntaxError: nonzero BigInt can't have leading zero")
	pt.parseError(".1n", "Semi-colon expected")
	pt.parseError("0.1n", "Semi-colon expected")
	pt.parseError("1e1n", "Semi-colon expected")
}

// Java: ParserTest.testBigIntLiteralWarning()
func TestBigIntLiteralWarning(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT_2019
	pt.parseWarning(
		"1n;",
		"This language feature is only supported for ECMASCRIPT_2020 mode or better: bigint")
}

// Java: ParserTest.testBigIntFeatureRecorded()
func TestBigIntFeatureRecorded(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("1n;")
	pt.expectFeatures(parser.FeatureSet_Feature_BIGINT)
}

// Java: ParserTest.testBigIntLiteralInCall()
func TestBigIntLiteralInCall(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("alert(1n)")
}

// Java: ParserTest.testBinaryLiterals()
func TestBinaryLiterals(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_BINARY_LITERALS)
	pt.mode = Config_LanguageMode_ECMASCRIPT3
	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parseWarning("0b0001;", requiresLanguageModeMessage(parser.FeatureSet_Feature_BINARY_LITERALS))
	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.parseWarning("0b0001;", requiresLanguageModeMessage(parser.FeatureSet_Feature_BINARY_LITERALS))
	pt.mode = Config_LanguageMode_ECMASCRIPT_2015
	pt.parse("0b0001;")
}

// Java: ParserTest.testOctalLiterals()
func TestOctalLiterals(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_OCTAL_LITERALS)
	pt.mode = Config_LanguageMode_ECMASCRIPT3
	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parseWarning("0o0001;", requiresLanguageModeMessage(parser.FeatureSet_Feature_OCTAL_LITERALS))
	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.parseWarning("0o0001;", requiresLanguageModeMessage(parser.FeatureSet_Feature_OCTAL_LITERALS))
	pt.mode = Config_LanguageMode_ECMASCRIPT_2015
	pt.parse("0o0001;")
}

// Java: ParserTest.testOctalEscapes()
func TestOctalEscapes(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT3
	pt.strictMode = Config_StrictMode_SLOPPY
	n := pt.parse("var x = 'This is a \u00a9 copyright symbol.'") // 251 is the octal value for "©"
	assertNode(t, n.GetFirstFirstChild()).IsName("x")
	assertNode(t, n.GetFirstFirstChild().GetFirstChild()).IsString("This is a \u00a9 copyright symbol.")
}

// Java: ParserTest.testOldStyleOctalLiterals()
func TestOldStyleOctalLiterals(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT3
	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parseWarning("0001;", "Octal integer literals are not supported in strict mode.")

	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.parseWarning("0001;", "Octal integer literals are not supported in strict mode.")

	pt.mode = Config_LanguageMode_ECMASCRIPT_2015
	pt.parseWarning("0001;", "Octal integer literals are not supported in strict mode.")
}

// Java: ParserTest.testOldStyleOctalLiterals_strictMode()
func TestOldStyleOctalLiterals_strictMode(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_STRICT

	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.parseError("0001;", "Octal integer literals are not supported in strict mode.")

	pt.mode = Config_LanguageMode_ECMASCRIPT_2015
	pt.parseError("0001;", "Octal integer literals are not supported in strict mode.")
}

// Java: ParserTest.testInvalidOctalLiterals()
func TestInvalidOctalLiterals(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT3
	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parseError("0o08;", "Invalid octal digit in octal literal.")

	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.parseError("0o08;", "Invalid octal digit in octal literal.")

	pt.mode = Config_LanguageMode_ECMASCRIPT_2015
	pt.parseError("0o08;", "Invalid octal digit in octal literal.")
}

// Java: ParserTest.testInvalidOldStyleOctalLiterals()
func TestInvalidOldStyleOctalLiterals(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT3
	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parseError("08;", "Invalid octal digit in octal literal.")
	pt.parseError("01238;", "Invalid octal digit in octal literal.")

	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.parseError("08;", "Invalid octal digit in octal literal.")
	pt.parseError("01238;", "Invalid octal digit in octal literal.")

	pt.mode = Config_LanguageMode_ECMASCRIPT_2015
	pt.parseError("08;", "Invalid octal digit in octal literal.")
	pt.parseError("01238;", "Invalid octal digit in octal literal.")
}

// Java: ParserTest.testGetter_ObjectLiteral_Es3()
func TestGetter_ObjectLiteral_Es3(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_GETTER)
	pt.mode = Config_LanguageMode_ECMASCRIPT3
	pt.strictMode = Config_StrictMode_SLOPPY

	pt.parseError("var x = {get 1(){}};", IRFactory_GETTER_ERROR_MESSAGE)
	pt.parseError("var x = {get 'a'(){}};", IRFactory_GETTER_ERROR_MESSAGE)
	pt.parseError("var x = {get a(){}};", IRFactory_GETTER_ERROR_MESSAGE)
	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.parse("var x = {get 1(){}};")
	pt.parse("var x = {get 'a'(){}};")
	pt.parse("var x = {get a(){}};")
}

// Java: ParserTest.testGetter_ObjectLiteral_Es5()
func TestGetter_ObjectLiteral_Es5(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_GETTER)
	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.strictMode = Config_StrictMode_SLOPPY

	pt.parse("var x = {get 1(){}};")
	pt.parse("var x = {get 'a'(){}};")
	pt.parse("var x = {get a(){}};")
}

// Java: ParserTest.testGetterInvalid_ObjectLiteral_EsNext()
func TestGetterInvalid_ObjectLiteral_EsNext(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures()
	pt.strictMode = Config_StrictMode_SLOPPY

	pt.parseError("var x = {get a(b){}};", "')' expected")
}

// Java: ParserTest.testGetter_Computed_ObjectLiteral_Es6()
func TestGetter_Computed_ObjectLiteral_Es6(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_GETTER, parser.FeatureSet_Feature_COMPUTED_PROPERTIES)
	pt.strictMode = Config_StrictMode_SLOPPY

	pt.parse("var x = {get [1](){}};")
	pt.parse("var x = {get ['a'](){}};")
	pt.parse("var x = {get [a](){}};")
}

// Java: ParserTest.testGetterInvalid_Computed_ObjectLiteral_EsNext()
func TestGetterInvalid_Computed_ObjectLiteral_EsNext(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures()
	pt.strictMode = Config_StrictMode_SLOPPY

	pt.parseError("var x = {get [a](b){}};", "')' expected")
}

// Java: ParserTest.testGetter_ClassSyntax()
func TestGetter_ClassSyntax(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_CLASSES, parser.FeatureSet_Feature_GETTER)
	pt.strictMode = Config_StrictMode_SLOPPY

	pt.parse("class Foo { get 1() {} };")
	pt.parse("class Foo { get 'a'() {} };")
	pt.parse("class Foo { get a() {} };")
}

// Java: ParserTest.testGetterInvalid_ClassSyntax_EsNext()
func TestGetterInvalid_ClassSyntax_EsNext(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures()
	pt.strictMode = Config_StrictMode_SLOPPY

	pt.parseError("class Foo { get a(b) {} };", "')' expected")
}

// Java: ParserTest.testGetter_Computed_ClassSyntax()
func TestGetter_Computed_ClassSyntax(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_CLASSES, parser.FeatureSet_Feature_GETTER, parser.FeatureSet_Feature_COMPUTED_PROPERTIES)
	pt.strictMode = Config_StrictMode_SLOPPY

	pt.parse("class Foo { get [1]() {} };")
	pt.parse("class Foo { get ['a']() {} };")
	pt.parse("class Foo { get [a]() {} };")
}

// Java: ParserTest.testGetterInvalid_Computed_ClassSyntax_EsNext()
func TestGetterInvalid_Computed_ClassSyntax_EsNext(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures()
	pt.strictMode = Config_StrictMode_SLOPPY

	pt.parseError("class Foo { get [a](b) {} };", "')' expected")
}

// Java: ParserTest.testSetter_ObjectLiteral_Es3()
func TestSetter_ObjectLiteral_Es3(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_SETTER)
	pt.mode = Config_LanguageMode_ECMASCRIPT3
	pt.strictMode = Config_StrictMode_SLOPPY

	pt.parseError("var x = {set 1(x){}};", IRFactory_SETTER_ERROR_MESSAGE)
	pt.parseError("var x = {set 'a'(x){}};", IRFactory_SETTER_ERROR_MESSAGE)
	pt.parseError("var x = {set a(x){}};", IRFactory_SETTER_ERROR_MESSAGE)
}

// Java: ParserTest.testSetter_ObjectLiteral_Es5()
func TestSetter_ObjectLiteral_Es5(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_SETTER)
	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.strictMode = Config_StrictMode_SLOPPY

	pt.parse("var x = {set 1(x){}};")
	pt.parse("var x = {set 'a'(x){}};")
	pt.parse("var x = {set a(x){}};")
}

// Java: ParserTest.testSetter_ObjectLiteral_Es6()
func TestSetter_ObjectLiteral_Es6(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_SETTER)
	pt.mode = Config_LanguageMode_ECMASCRIPT_2015
	pt.strictMode = Config_StrictMode_SLOPPY

	pt.parse("var x = {set 1(x){}};")
	pt.parse("var x = {set 'a'(x){}};")
	pt.parse("var x = {set a(x){}};")

	pt.parse("var x = {set setter(x = 5) {}};")
	pt.parse("var x = {set setter(x = a) {}};")
	pt.parse("var x = {set setter(x = a + 5) {}};")

	pt.parse("var x = {set setter([x, y, z]) {}};")
	pt.parse("var x = {set setter([x, y, ...z]) {}};")
	pt.parse("var x = {set setter([x, y, z] = [1, 2, 3]) {}};")
	pt.parse("var x = {set setter([x = 1, y = 2, z = 3]) {}};")

	pt.parse("var x = {set setter({x, y, z}) {}};")
	pt.parse("var x = {set setter({x, y, z} = {x: 1, y: 2, z: 3}) {}};")
	pt.parse("var x = {set setter({x = 1, y = 2, z = 3}) {}};")
}

// Java: ParserTest.testSetterInvalid_ObjectLiteral_EsNext()
func TestSetterInvalid_ObjectLiteral_EsNext(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures()
	pt.strictMode = Config_StrictMode_SLOPPY

	pt.parseError("var x = {set a() {}};", "Setter must have exactly 1 parameter, found 0")
	pt.parseError("var x = {set a(x, y) {}};", "Setter must have exactly 1 parameter, found 2")
	pt.parseError("var x = {set a(...x, y) {}};", "Setter must have exactly 1 parameter, found 2")
	pt.parseError("var x = {set a(...x) {}};", "Setter must not have a rest parameter")
}

// Java: ParserTest.testSetter_Computed_ObjectLiteral_Es6()
func TestSetter_Computed_ObjectLiteral_Es6(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_SETTER, parser.FeatureSet_Feature_COMPUTED_PROPERTIES)
	pt.mode = Config_LanguageMode_ECMASCRIPT_2015
	pt.strictMode = Config_StrictMode_SLOPPY

	pt.parse("var x = {set [setter](x = 5) {}};")
	pt.parse("var x = {set [setter](x = a) {}};")
	pt.parse("var x = {set [setter](x = a + 5) {}};")

	pt.parse("var x = {set [setter]([x, y, z]) {}};")
	pt.parse("var x = {set [setter]([x, y, ...z]) {}};")
	pt.parse("var x = {set [setter]([x, y, z] = [1, 2, 3]) {}};")
	pt.parse("var x = {set [setter]([x = 1, y = 2, z = 3]) {}};")

	pt.parse("var x = {set [setter]({x, y, z}) {}};")
	pt.parse("var x = {set [setter]({x, y, z} = {x: 1, y: 2, z: 3}) {}};")
	pt.parse("var x = {set [setter]({x = 1, y = 2, z = 3}) {}};")
}

// Java: ParserTest.testSetterInvalid_Computed_ObjectLiteral_EsNext()
func TestSetterInvalid_Computed_ObjectLiteral_EsNext(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures()
	pt.strictMode = Config_StrictMode_SLOPPY

	pt.parseError("var x = {set [setter]() {}};", "Setter must have exactly 1 parameter, found 0")
	pt.parseError("var x = {set [setter](x, y) {}};", "Setter must have exactly 1 parameter, found 2")
	pt.parseError(
		"var x = {set [setter](...x, y) {}};", "Setter must have exactly 1 parameter, found 2")
	pt.parseError("var x = {set [setter](...x) {}};", "Setter must not have a rest parameter")
}

// Java: ParserTest.testSetter_ClassSyntax()
func TestSetter_ClassSyntax(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_CLASSES, parser.FeatureSet_Feature_SETTER)

	pt.parse("class Foo { set setter(x = 5) {} };")
	pt.parse("class Foo { set setter(x = a) {} };")
	pt.parse("class Foo { set setter(x = a + 5) {} };")

	pt.parse("class Foo { set setter([x, y, z]) {} };")
	pt.parse("class Foo { set setter([x, y, ...z]) {}};")
	pt.parse("class Foo { set setter([x, y, z] = [1, 2, 3]) {} };")
	pt.parse("class Foo { set setter([x = 1, y = 2, z = 3]) {} };")

	pt.parse("class Foo { set setter({x, y, z}) {}};")
	pt.parse("class Foo { set setter({x, y, z} = {x: 1, y: 2, z: 3}) {} };")
	pt.parse("class Foo { set setter({x = 1, y = 2, z = 3}) {} };")
}

// Java: ParserTest.testSetterInvalid_ClassSyntax_EsNext()
func TestSetterInvalid_ClassSyntax_EsNext(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures()

	pt.parseError("class Foo { set setter() {} };", "Setter must have exactly 1 parameter, found 0")
	pt.parseError(
		"class Foo { set setter(x, y) {} };", "Setter must have exactly 1 parameter, found 2")
	pt.parseError(
		"class Foo { set setter(...x, y) {} };", "Setter must have exactly 1 parameter, found 2")
	pt.parseError("class Foo { set setter(...x) {} };", "Setter must not have a rest parameter")
}

// Java: ParserTest.testSetter_Computed_ClassSyntax()
func TestSetter_Computed_ClassSyntax(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_CLASSES, parser.FeatureSet_Feature_SETTER, parser.FeatureSet_Feature_COMPUTED_PROPERTIES)
	pt.mode = Config_LanguageMode_ECMASCRIPT_2015

	pt.parse("class Foo { set [setter](x = 5) {} };")
	pt.parse("class Foo { set [setter](x = a) {} };")
	pt.parse("class Foo { set [setter](x = a + 5) {} };")

	pt.parse("class Foo { set [setter]([x, y, z]) {} };")
	pt.parse("class Foo { set [setter]([x, y, ...z]) {}};")
	pt.parse("class Foo { set [setter]([x, y, z] = [1, 2, 3]) {} };")
	pt.parse("class Foo { set [setter]([x = 1, y = 2, z = 3]) {} };")

	pt.parse("class Foo { set [setter]({x, y, z}) {}};")
	pt.parse("class Foo { set [setter]({x, y, z} = {x: 1, y: 2, z: 3}) {} };")
	pt.parse("class Foo { set [setter]({x = 1, y = 2, z = 3}) {} };")
}

// Java: ParserTest.testSetterInvalid_Computed_ClassSyntax_EsNext()
func TestSetterInvalid_Computed_ClassSyntax_EsNext(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures()

	pt.parseError("class Foo { set [setter]() {} };", "Setter must have exactly 1 parameter, found 0")
	pt.parseError(
		"class Foo { set [setter](x, y) {} };", "Setter must have exactly 1 parameter, found 2")
	pt.parseError(
		"class Foo { set [setter](...x, y) {} };", "Setter must have exactly 1 parameter, found 2")
	pt.parseError("class Foo { set [setter](...x) {} };", "Setter must not have a rest parameter")
}

// Java: ParserTest.testLamestWarningEver()
func TestLamestWarningEver(t *testing.T) {
	pt := newParserTest(t)

	// This used to be a warning.
	pt.parse("var x = /** @type {undefined} */ (y);")
	pt.parse("var x = /** @type {void} */ (y);")
}

// Java: ParserTest.testUnfinishedComment()
func TestUnfinishedComment(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("/** this is a comment ", "unterminated comment")
}

// Java: ParserTest.testHtmlStartCommentAtStartOfLine()
func TestHtmlStartCommentAtStartOfLine(t *testing.T) {
	pt := newParserTest(t)

	pt.parseWarning("<!-- This text is ignored.\nalert(1)", parserTest_HTML_COMMENT_WARNING)
}

// Java: ParserTest.testHtmlStartComment()
func TestHtmlStartComment(t *testing.T) {
	pt := newParserTest(t)

	pt.parseWarning("alert(1) <!-- This text is ignored.\nalert(2)", parserTest_HTML_COMMENT_WARNING)
}

// Java: ParserTest.testHtmlEndCommentAtStartOfLine()
func TestHtmlEndCommentAtStartOfLine(t *testing.T) {
	pt := newParserTest(t)

	pt.parseWarning("alert(1)\n --> This text is ignored.", parserTest_HTML_COMMENT_WARNING)
}

// Java: ParserTest.testHtmlEndComment()
func TestHtmlEndComment(t *testing.T) {
	pt := newParserTest(t)

	pt.parse(`while (x --> 0) {
  alert(1)
}`)
}

// Java: ParserTest.testParseBlockDescription()
func TestParseBlockDescription(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true

	n := pt.parse("/** This is a variable. */ var x;")
	var_ := n.GetFirstChild()
	assertThat(t, var_.GetJSDocInfo()).IsNotNull()
	assertThat(t, var_.GetJSDocInfo().GetBlockDescription()).IsEqualTo("This is a variable.")
}

// Java: ParserTest.testUnnamedFunctionStatement()
func TestUnnamedFunctionStatement(t *testing.T) {
	pt := newParserTest(t)

	// Statements
	pt.parseError("function() {};", "'identifier' expected")
	pt.parseError("if (true) { function() {}; }", "'identifier' expected")
	pt.parse("function f() {};")
	// Expressions
	pt.parse("(function f() {});")
	pt.parse("(function () {});")
}

// Java: ParserTest.testReservedKeywords()
func TestReservedKeywords(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_ES3_KEYWORDS_AS_IDENTIFIERS)
	pt.mode = Config_LanguageMode_ECMASCRIPT3
	pt.strictMode = Config_StrictMode_SLOPPY

	pt.parseError("var boolean;", "identifier is a reserved word")
	pt.parseError("function boolean() {};", "identifier is a reserved word")
	pt.parseError("boolean = 1;", "identifier is a reserved word")

	pt.expectFeatures()
	pt.parseError("class = 1;", "'identifier' expected")
	pt.parseError("public = 2;", "primary expression expected")

	pt.mode = Config_LanguageMode_ECMASCRIPT5

	pt.expectFeatures(parser.FeatureSet_Feature_ES3_KEYWORDS_AS_IDENTIFIERS)
	pt.parse("var boolean;")
	pt.parse("function boolean() {};")
	pt.parse("boolean = 1;")

	pt.expectFeatures()
	pt.parseError("class = 1;", "'identifier' expected")
	pt.parseError("var import = 0;", "'identifier' expected")
	// TODO(johnlenz): reenable
	// parse("public = 2;");

	pt.mode = Config_LanguageMode_ECMASCRIPT5

	pt.expectFeatures(parser.FeatureSet_Feature_ES3_KEYWORDS_AS_IDENTIFIERS)
	pt.parse("var boolean;")
	pt.parse("function boolean() {};")
	pt.parse("boolean = 1;")

	pt.expectFeatures()
	pt.parseError("public = 2;", "primary expression expected")
	pt.parseError("class = 1;", "'identifier' expected")

	pt.mode = Config_LanguageMode_ECMASCRIPT_2015
	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parseError("const else = 1;", "'identifier' expected")
}

// Java: ParserTest.testTypeScriptKeywords()
func TestTypeScriptKeywords(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("type = 2;")
	pt.parse("var type = 3;")
	pt.parse("type\nx = 5")
	pt.parse("while (i--) { type = types[i]; }")

	pt.parse("declare = 2;")
	pt.parse("var declare = 3;")
	pt.parse("declare\nx = 5")
	pt.parse("while (i--) { declare = declares[i]; }")

	pt.parse("module = 2;")
	pt.parse("var module = 3;")
	pt.parse("module\nx = 5")
	pt.parse("while (i--) { module = module[i]; }")
}

// Java: ParserTest.testKeywordsAsProperties1()
func TestKeywordsAsProperties1(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_KEYWORDS_AS_PROPERTIES)
	pt.mode = Config_LanguageMode_ECMASCRIPT3
	pt.strictMode = Config_StrictMode_SLOPPY

	pt.parseWarning("var x = {function: 1};", IRFactory_INVALID_ES3_PROP_NAME)
	pt.parseWarning("x.function;", IRFactory_INVALID_ES3_PROP_NAME)
	pt.parseWarning("var x = {class: 1};", IRFactory_INVALID_ES3_PROP_NAME)
	pt.expectFeatures()
	pt.parse("var x = {'class': 1};")
	pt.expectFeatures(parser.FeatureSet_Feature_KEYWORDS_AS_PROPERTIES)
	pt.parseWarning("x.class;", IRFactory_INVALID_ES3_PROP_NAME)
	pt.expectFeatures()
	pt.parse("x['class'];")
	pt.parse("var x = {let: 1};") // 'let' is not reserved in ES3
	pt.parse("x.let;")
	pt.parse("var x = {yield: 1};") // 'yield' is not reserved in ES3
	pt.parse("x.yield;")
	pt.expectFeatures(parser.FeatureSet_Feature_KEYWORDS_AS_PROPERTIES)
	pt.parseWarning("x.prototype.catch = function() {};", IRFactory_INVALID_ES3_PROP_NAME)
	pt.parseWarning("x().catch();", IRFactory_INVALID_ES3_PROP_NAME)

	pt.mode = Config_LanguageMode_ECMASCRIPT5

	pt.parse("var x = {function: 1};")
	pt.parse("x.function;")
	pt.parse("var x = {get function(){} };")
	pt.parse("var x = {set function(a){} };")
	pt.parse("var x = {class: 1};")
	pt.parse("x.class;")
	pt.expectFeatures()
	pt.parse("var x = {let: 1};")
	pt.parse("x.let;")
	pt.parse("var x = {yield: 1};")
	pt.parse("x.yield;")
	pt.expectFeatures(parser.FeatureSet_Feature_KEYWORDS_AS_PROPERTIES)
	pt.parse("x.prototype.catch = function() {};")
	pt.parse("x().catch();")

	pt.mode = Config_LanguageMode_ECMASCRIPT5

	pt.parse("var x = {function: 1};")
	pt.parse("x.function;")
	pt.parse("var x = {get function(){} };")
	pt.parse("var x = {set function(a){} };")
	pt.parse("var x = {class: 1};")
	pt.parse("x.class;")
	pt.expectFeatures()
	pt.parse("var x = {let: 1};")
	pt.parse("x.let;")
	pt.parse("var x = {yield: 1};")
	pt.parse("x.yield;")
	pt.expectFeatures(parser.FeatureSet_Feature_KEYWORDS_AS_PROPERTIES)
	pt.parse("x.prototype.catch = function() {};")
	pt.parse("x().catch();")
}

// Java: ParserTest.testKeywordsAsProperties2()
func TestKeywordsAsProperties2(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("var x = {get 'function'(){} };")
	pt.parse("var x = {get 1(){} };")
	pt.parse("var x = {set 'function'(a){} };")
	pt.parse("var x = {set 1(a){} };")
}

// Java: ParserTest.testKeywordsAsProperties3()
func TestKeywordsAsProperties3(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("var x = {get 'function'(){} };")
	pt.parse("var x = {get 1(){} };")
	pt.parse("var x = {set 'function'(a){} };")
	pt.parse("var x = {set 1(a){} };")
}

// Java: ParserTest.testKeywordsAsPropertiesInExterns1()
func TestKeywordsAsPropertiesInExterns1(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT3
	pt.strictMode = Config_StrictMode_SLOPPY

	pt.parse(`/** @fileoverview
@externs
*/
 var x = {function: 1};`)
}

// Java: ParserTest.testKeywordsAsPropertiesInExterns2()
func TestKeywordsAsPropertiesInExterns2(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT3
	pt.strictMode = Config_StrictMode_SLOPPY

	pt.parse(`/** @fileoverview
@externs
*/
 var x = {}; x.function + 1;`)
}

// Java: ParserTest.testExternsLanguageMode()
func TestExternsLanguageMode(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT3
	pt.strictMode = Config_StrictMode_SLOPPY

	pt.parse(`/** @fileoverview
@externs
*/
 const x = 1;`)
}

// Java: ParserTest.testUnicodeInIdentifiers()
func TestUnicodeInIdentifiers(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("var \u00e0")
	pt.parse("var cos\u03b8")
	pt.parse("if(true){foo=\u03b1}")
	pt.parse("if(true){foo=\u0394}else bar()")
}

// Java: ParserTest.testUnicodeEscapeInIdentifiers()
func TestUnicodeEscapeInIdentifiers(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("var \\u00fb")
	pt.parse("var \\u00fbtest\\u00fb")
	pt.parse("Js\\u00C7ompiler")
	pt.parse("Js\\u0043ompiler")
	pt.parse("if(true){foo=\\u03b5}")
	pt.parse("if(true){foo=\\u03b5}else bar()")
}

// Java: ParserTest.testUnicodePointEscapeInIdentifiers()
func TestUnicodePointEscapeInIdentifiers(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("var \\u{0043}")
	pt.parse("var \\u{0043}test\\u{0043}")
	pt.parse("var \\u0043test\\u{0043}")
	pt.parse("var \\u{0043}test\\u0043")
	pt.parse("Js\\u{0043}ompiler")
	pt.parse("Js\\u{275}ompiler")
	pt.parse("var \\u0043;{43}")
}

// Java: ParserTest.testUnicodePointEscapeStringLiterals()
func TestUnicodePointEscapeStringLiterals(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("var i = '\\u0043ompiler'")
	pt.parse("var i = '\\u{43}ompiler'")
	pt.parse("var i = '\\u{1f42a}ompiler'")
	pt.parse("var i = '\\u{2603}ompiler'")
	pt.parse("var i = '\\u{1}ompiler'")
}

// Java: ParserTest.testUnicodePointEscapeTemplateLiterals()
func TestUnicodePointEscapeTemplateLiterals(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("var i = `\\u0043ompiler`")
	pt.parse("var i = `\\u{43}ompiler`")
	pt.parse("var i = `\\u{1f42a}ompiler`")
	pt.parse("var i = `\\u{2603}ompiler`")
	pt.parse("var i = `\\u{1}ompiler`")
}

// Java: ParserTest.testInvalidUnicodePointEscapeInIdentifiers()
func TestInvalidUnicodePointEscapeInIdentifiers(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("var \\u{defg", "Invalid escape sequence")
	pt.parseError("var \\u{03b5", "Invalid escape sequence")
	pt.parseError("var \\u43{43}", "Invalid escape sequence")
	pt.parseError("var \\u{defgRestOfIdentifier", "Invalid escape sequence")
	pt.parseError("var \\u03b5}", "primary expression expected")
	pt.parseError("var \\u{03b5}}}", "primary expression expected")
	pt.parseError("var \\u{03b5}{}", parserTest_SEMICOLON_EXPECTED)
	pt.parseError("var \\u0043{43}", parserTest_SEMICOLON_EXPECTED)
	pt.parseError("var \\u{DEFG}", "Invalid escape sequence")
	pt.parseError("Js\\u{}ompiler", "Invalid escape sequence")
	// Legal unicode but invalid in identifier
	pt.parseError("Js\\u{99}ompiler", "Invalid escape sequence")
	pt.parseError("Js\\u{10000}ompiler", "Invalid escape sequence")
	// Code points outside the BMP must not be truncated to their low 16 bits. U+10041 shares its
	// low bits with U+0041 ('A') and U+110041 is above the U+10FFFF maximum; both used to be
	// silently accepted.
	pt.parseError("Js\\u{10041}ompiler", "Invalid escape sequence")
	pt.parseError("Js\\u{110041}ompiler", "Invalid escape sequence")
}

// Java: ParserTest.testInvalidUnicodePointEscapeStringLiterals()
func TestInvalidUnicodePointEscapeStringLiterals(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("var i = '\\u{defg'", "Hex digit expected")
	pt.parseError("var i = '\\u{defgRestOfIdentifier'", "Hex digit expected")
	pt.parseError("var i = '\\u{DEFG}'", "Hex digit expected")
	pt.parseError("var i = 'Js\\u{}ompiler'", "Empty unicode escape")
	pt.parseError("var i = '\\u{345", "Hex digit expected")
	pt.parseError("var i = '\\u{110000}'", "Undefined Unicode code-point")
}

// Java: ParserTest.testInvalidUnicodePointEscapeTemplateLiterals()
func TestInvalidUnicodePointEscapeTemplateLiterals(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("var i = `\\u{defg`", "Hex digit expected")
	pt.parseError("var i = `\\u{defgRestOfIdentifier`", "Hex digit expected")
	pt.parseError("var i = `\\u{DEFG}`", "Hex digit expected")
	pt.parseError("var i = `Js\\u{}ompiler`", "Empty unicode escape")
	pt.parseError("var i = `\\u{345`", "Hex digit expected")
	pt.parseError("var i = `\\u{110000}`", "Undefined Unicode code-point")
}

// Java: ParserTest.testEs2018LiftIllegalEscapeSequenceRestrictionOnTaggedTemplates()
func TestEs2018LiftIllegalEscapeSequenceRestrictionOnTaggedTemplates(t *testing.T) {
	pt := newParserTest(t)

	// These should not generate errors, even though they contain illegal escape sequences.
	// https://github.com/tc39/proposal-template-literal-revision
	pt.parse("latex`\\unicode`")
	pt.parse("foo`\\xerxes`")
	pt.parse("bar`\\u{h}ere`")
	pt.parse("bar`\\u{43`")

	// tagged malformed template literal throws error
	pt.parseError("foo`\\unicode", "Unterminated template literal")
	// normal template literals still throw error
	pt.parseError("var bad = `\\unicode`;", "Hex digit expected")
}

// Java: ParserTest.testInvalidEscape()
func TestInvalidEscape(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("var \\x39abc", "Invalid escape sequence")
	pt.parseError("var abc\\t", "Invalid escape sequence")
}

// Java: ParserTest.testUnnecessaryEscape()
func TestUnnecessaryEscape(t *testing.T) {
	pt := newParserTest(t)

	pt.parseWarning("var str = '\\a'", "Unnecessary escape: '\\a' is equivalent to just 'a'")
	pt.parse("var str = '\\b'")
	pt.parseWarning("var str = '\\c'", "Unnecessary escape: '\\c' is equivalent to just 'c'")
	pt.parseWarning("var str = '\\d'", "Unnecessary escape: '\\d' is equivalent to just 'd'")
	pt.parseWarning("var str = '\\e'", "Unnecessary escape: '\\e' is equivalent to just 'e'")
	pt.parse("var str = '\\f'")
	pt.parse("var str = '\\/'")
	pt.parse("var str = '\\0'")
	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parseWarning("var str = '\\1'", "Unnecessary escape: '\\1' is equivalent to just '1'")
	pt.parseWarning("var str = '\\2'", "Unnecessary escape: '\\2' is equivalent to just '2'")
	pt.parseWarning("var str = '\\3'", "Unnecessary escape: '\\3' is equivalent to just '3'")
	pt.parseWarning("var str = '\\4'", "Unnecessary escape: '\\4' is equivalent to just '4'")
	pt.parseWarning("var str = '\\5'", "Unnecessary escape: '\\5' is equivalent to just '5'")
	pt.parseWarning("var str = '\\6'", "Unnecessary escape: '\\6' is equivalent to just '6'")
	pt.parseWarning("var str = '\\7'", "Unnecessary escape: '\\7' is equivalent to just '7'")
	pt.parseWarning("var str = '\\8'", "Unnecessary escape: '\\8' is equivalent to just '8'")
	pt.parseWarning("var str = '\\9'", "Unnecessary escape: '\\9' is equivalent to just '9'")
	pt.parseWarning("var str = '\\%'", "Unnecessary escape: '\\%' is equivalent to just '%'")

	pt.parseWarning("var str = '\\$'", "Unnecessary escape: '\\$' is equivalent to just '$'")
}

// Java: ParserTest.testUnnecessaryEscapeUntaggedTemplateLiterals()
func TestUnnecessaryEscapeUntaggedTemplateLiterals(t *testing.T) {
	pt := newParserTest(t)

	pt.parseWarning("var str = `\\a`", "Unnecessary escape: '\\a' is equivalent to just 'a'")
	pt.parse("var str = `\\b`")
	pt.parseWarning("var str = `\\c`", "Unnecessary escape: '\\c' is equivalent to just 'c'")
	pt.parseWarning("var str = `\\d`", "Unnecessary escape: '\\d' is equivalent to just 'd'")
	pt.parseWarning("var str = `\\e`", "Unnecessary escape: '\\e' is equivalent to just 'e'")
	pt.parse("var str = `\\f`")
	pt.parseWarning("var str = `\\/`", "Unnecessary escape: '\\/' is equivalent to just '/'")
	pt.parse("var str = `\\0`")
	pt.parseWarning("var str = `\\%`", "Unnecessary escape: '\\%' is equivalent to just '%'")

	// single and double quotes have no meaning in a template lit
	pt.parseWarning("var str = `\\\"`", "Unnecessary escape: '\\\"' is equivalent to just '\"'")
	pt.parseWarning("var str = `\\'`", "Unnecessary escape: \"\\'\" is equivalent to just \"'\"")

	// $ needs to be escaped to distinguish it from use of ${}
	pt.parse("var str = `\\$`")
	// ` needs to be escaped to avoid ending the template lit
	pt.parse("var str = `\\``")
}

// Java: ParserTest.testUnnecessaryEscapeTaggedTemplateLiterals()
func TestUnnecessaryEscapeTaggedTemplateLiterals(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_TEMPLATE_LITERALS)

	// Don't warn for unnecessary escapes in tagged template literals since they may access the
	// raw string value
	pt.parse("var str = String.raw`\\a`")
	pt.parse("var str = String.raw`\\b`")
	pt.parse("var str = String.raw`\\c`")
	pt.parse("var str = String.raw`\\d`")
	pt.parse("var str = String.raw`\\e`")
	pt.parse("var str = String.raw`\\f`")
	pt.parse("var str = String.raw`\\/`")
	pt.parse("var str = String.raw`\\0`")
	pt.parse("var str = String.raw`\\8`")
	pt.parse("var str = String.raw`\\9`")
	pt.parse("var str = String.raw`\\%`")

	pt.parse("var str = String.raw`\\$`")
	pt.parse("var str = String.raw`\\``")
}

// Java: ParserTest.testEOFInUnicodeEscape()
func TestEOFInUnicodeEscape(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("var \\u1", "Invalid escape sequence")
	pt.parseError("var \\u12", "Invalid escape sequence")
	pt.parseError("var \\u123", "Invalid escape sequence")
}

// Java: ParserTest.testEndOfIdentifierInUnicodeEscape()
func TestEndOfIdentifierInUnicodeEscape(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("var \\u1 = 1;", "Invalid escape sequence")
	pt.parseError("var \\u12 = 2;", "Invalid escape sequence")
	pt.parseError("var \\u123 = 3;", "Invalid escape sequence")
}

// Java: ParserTest.testInvalidUnicodeEscape()
func TestInvalidUnicodeEscape(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("var \\uDEFG", "Invalid escape sequence")
}

// Java: ParserTest.testUnicodeEscapeInvalidIdentifierStart()
func TestUnicodeEscapeInvalidIdentifierStart(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("var \\u0037yler", "Character '7' (U+0037) is not a valid identifier start char")
	pt.parseError("var \\u{37}yler", "Character '7' (U+0037) is not a valid identifier start char")
	pt.parseError("var \\u0020space", "Invalid escape sequence")
}

// Java: ParserTest.testUnicodeEscapeInvalidIdentifierChar()
func TestUnicodeEscapeInvalidIdentifierChar(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("var sp\\u0020ce", "Invalid escape sequence")
}

// Java: ParserTest.testKeywordAsIdentifier()
func TestKeywordAsIdentifier(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("var while;", "'identifier' expected")
	pt.parseError("var wh\\u0069le;", "'identifier' expected")
	pt.parseError("var v\\u0061r;", "'identifier' expected")
}

// Java: ParserTest.testUnicodeEscapeInKeywords()
func TestUnicodeEscapeInKeywords(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("v\\u0061r b = 1;", "'identifier' expected")
	pt.parseError("i\\u0066 (b) { }", "'identifier' expected")
	pt.parseError("f\\u0075nction foo() {}", "'identifier' expected")
	pt.parseError("l\\u0065t x = 1;", "'identifier' expected")
	pt.parseError("l\\u0065t = 1;", "'identifier' expected")
	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parse("l\\u0065t = 1;")
}

// Java: ParserTest.testUnicodeEscapeInPropertyNames()
func TestUnicodeEscapeInPropertyNames(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("obj.v\\u0061r")
	pt.parse("var x = { v\\u0061r: 1 };")
}

// Java: ParserTest.testGetPropFunctionName()
func TestGetPropFunctionName(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("function a.b() {}", "'(' expected")
	pt.parseError("var x = function a.b() {}", "'(' expected")
}

// Java: ParserTest.testIdeModePartialTree()
func TestIdeModePartialTree(t *testing.T) {
	pt := newParserTest(t)

	partialTree := pt.parseError("function Foo() {} f.", "'identifier' expected")
	assertThat(t, partialTree).IsNull()

	pt.isIdeMode = true
	partialTree = pt.parseError("function Foo() {} f.", "'identifier' expected")
	assertThat(t, partialTree).IsNotNull()
}

// Java: ParserTest.testForEach()
func TestForEach(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError(
		`function f(stamp, status) {
  for each ( var curTiming in this.timeLog.timings ) {
    if ( curTiming.callId == stamp ) {
      curTiming.flag = status;
      break;
    }
  }
};
`,
		"'(' expected")
}

// Java: ParserTest.testValidTypeAnnotation1()
func TestValidTypeAnnotation1(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("/** @type {string} */ var o = 'str';")
	pt.parse("var /** @type {string} */ o = 'str', /** @type {number} */ p = 0;")
	pt.parse("/** @type {function():string} */ function o() { return 'str'; }")
	pt.parse("var o = {}; /** @type {string} */ o.prop = 'str';")
	pt.parse("var o = {}; /** @type {string} */ o['prop'] = 'str';")
	pt.parse("var o = { /** @type {string} */ prop : 'str' };")
	pt.parse("var o = { /** @type {string} */ 'prop' : 'str' };")
	pt.parse("var o = { /** @type {string} */ 1 : 'str' };")
}

// Java: ParserTest.testValidTypeAnnotation2()
func TestValidTypeAnnotation2(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY
	pt.expectFeatures(parser.FeatureSet_Feature_GETTER)
	pt.parse("var o = { /** @type {string} */ get prop() { return 'str' }};")
	pt.expectFeatures(parser.FeatureSet_Feature_SETTER)
	pt.parse("var o = { /** @type {string} */ set prop(s) {}};")
}

// Java: ParserTest.testValidTypeAnnotation3()
func TestValidTypeAnnotation3(t *testing.T) {
	pt := newParserTest(t)

	// This one we don't currently support in the type checker but
	// we would like to.
	pt.parse("try {} catch (/** @type {Error} */ e) {}")
}

// Java: ParserTest.testValidTypeAnnotation4()
func TestValidTypeAnnotation4(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY
	pt.expectFeatures(parser.FeatureSet_Feature_MODULES)
	pt.parse("/** @type {number} */ export var x = 3;")
}

// Java: ParserTest.testTypeofJsdoc()
func TestTypeofJsdoc(t *testing.T) {
	pt := newParserTest(t)

	pt.assertNodeEquality(pt.parse("var b = 0;"), pt.parse("var /** typeof a */ b = 0;"))

	pt.assertNodeEquality(
		pt.parse("var b = 0;"),
		pt.parseWarning("var /** typeof {a} */ b = 0;", parserTest_UNNECESSARY_BRACES_MESSAGE))

	pt.assertNodeEquality(
		pt.parse("var b = 0;"),
		pt.parseWarning("var /** typeof <a> */ b = 0;", parserTest_NAME_NOT_RECOGNIZED_MESSAGE))
}

// Java: ParserTest.testParsingAssociativity()
func TestParsingAssociativity(t *testing.T) {
	pt := newParserTest(t)

	pt.assertNodeEquality(pt.parse("x * y * z"), pt.parse("(x * y) * z"))
	pt.assertNodeEquality(pt.parse("x + y + z"), pt.parse("(x + y) + z"))
	pt.assertNodeEquality(pt.parse("x | y | z"), pt.parse("(x | y) | z"))
	pt.assertNodeEquality(pt.parse("x & y & z"), pt.parse("(x & y) & z"))
	pt.assertNodeEquality(pt.parse("x ^ y ^ z"), pt.parse("(x ^ y) ^ z"))
	pt.assertNodeEquality(pt.parse("x || y || z"), pt.parse("(x || y) || z"))
	pt.assertNodeEquality(pt.parse("x && y && z"), pt.parse("(x && y) && z"))
}

// Java: ParserTest.testIssue1116()
func TestIssue1116(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("/**/")
}

// Java: ParserTest.testUnterminatedStringLiteral()
func TestUnterminatedStringLiteral(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("var unterm = 'forgot closing quote", "Unterminated string literal")

	pt.parseError(
		`var unterm = 'forgot closing quote
alert(unterm);
`,
		"Unterminated string literal")

	// test combo of a string continuation + useless escape warning + unterminated literal error
	// create a TestErrorReporter so that we can expect both a warning and an error
	js := "var unterm = ' \\\n \\a \n"

	testErrorReporter :=
		newTestErrorReporter().
			expectAllWarnings("Unnecessary escape: '\\a' is equivalent to just 'a'").
			expectAllErrors("Unterminated string literal")
	file := rhino.NewSimpleSourceFile("input", rhino.StaticSourceFile_SourceKind_STRONG)
	ParserRunner_Parse(file, js, pt.createConfig(), testErrorReporter)

	// verify we reported both the warning and error
	testErrorReporter.verifyHasEncounteredAllWarningsAndErrors(t)
}

// Java: ParserTest.testUnterminatedRegExp()
func TestUnterminatedRegExp(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("var unterm = /forgot trailing slash", "Expected '/' in regular expression literal")

	pt.parseError(
		`var unterm = /forgot trailing slash
alert(unterm);
`,
		"Expected '/' in regular expression literal")
}

// Java: ParserTest.testRegExp()
func TestRegExp(t *testing.T) {
	pt := newParserTest(t)

	pt.assertNodeEquality(pt.parse("/a/"), pt.script(pt.expr(pt.regex("a"))))
	pt.assertNodeEquality(pt.parse("/\\\\/"), pt.script(pt.expr(pt.regex("\\\\"))))
	pt.assertNodeEquality(pt.parse("/\\s/"), pt.script(pt.expr(pt.regex("\\s"))))
	pt.assertNodeEquality(pt.parse("/\\u000A/"), pt.script(pt.expr(pt.regex("\\u000A"))))
	pt.assertNodeEquality(pt.parse("/[\\]]/"), pt.script(pt.expr(pt.regex("[\\]]"))))
	pt.assertNodeEquality(pt.parse("/\\k<group>/"), pt.script(pt.expr(pt.regex("\\k<group>"))))
	pt.assertNodeEquality(
		pt.parse("/(?<group>a)\\k<group>/"), pt.script(pt.expr(pt.regex("(?<group>a)\\k<group>"))))
}

// Java: ParserTest.testRegExpError()
func TestRegExpError(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("/a\\/", "Expected '/' in regular expression literal")
	pt.parseError("/\\ca\\/", "Expected '/' in regular expression literal")
	pt.parseError("/\u0008.\\/", "Expected '/' in regular expression literal")
}

// Java: ParserTest.testRegExpUnicode()
func TestRegExpUnicode(t *testing.T) {
	pt := newParserTest(t)

	pt.assertNodeEquality(pt.parse("/\\u10fA/"), pt.script(pt.expr(pt.regex("\\u10fA"))))
	pt.assertNodeEquality(pt.parse("/\\u{10fA}/u"), pt.script(pt.expr(pt.regex("\\u{10fA}", "u"))))
	pt.assertNodeEquality(pt.parse("/\\u{1fA}/u"), pt.script(pt.expr(pt.regex("\\u{1fA}", "u"))))
	pt.assertNodeEquality(pt.parse("/\\u{10FFFF}/u"), pt.script(pt.expr(pt.regex("\\u{10FFFF}", "u"))))
}

// Java: ParserTest.testRegExpFlags()
func TestRegExpFlags(t *testing.T) {
	pt := newParserTest(t)

	// Various valid combinations.
	pt.parse("/a/")
	pt.parse("/a/i")
	pt.parse("/a/g")
	pt.parse("/a/m")
	pt.parse("/a/ig")
	pt.parse("/a/gm")
	pt.parse("/a/mgi")

	// Invalid combinations
	pt.parseError("/a/a", "Invalid RegExp flag 'a'")
	pt.parseError("/a/b", "Invalid RegExp flag 'b'")
	pt.parseError(
		"/a/abc", "Invalid RegExp flag 'a'", "Invalid RegExp flag 'b'", "Invalid RegExp flag 'c'")
}

// Java: ParserTest.testES6RegExpFlags()
func TestES6RegExpFlags(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_REGEXP_FLAG_Y)
	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parse("/a/y")
	pt.expectFeatures(parser.FeatureSet_Feature_REGEXP_FLAG_U)
	pt.parse("/a/u")

	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.expectFeatures(parser.FeatureSet_Feature_REGEXP_FLAG_Y)
	pt.parseWarning("/a/y", requiresLanguageModeMessage(parser.FeatureSet_Feature_REGEXP_FLAG_Y))
	pt.expectFeatures(parser.FeatureSet_Feature_REGEXP_FLAG_U)
	pt.parseWarning("/a/u", requiresLanguageModeMessage(parser.FeatureSet_Feature_REGEXP_FLAG_U))
	pt.parseWarning(
		"/a/yu",
		requiresLanguageModeMessage(parser.FeatureSet_Feature_REGEXP_FLAG_Y),
		requiresLanguageModeMessage(parser.FeatureSet_Feature_REGEXP_FLAG_U))
}

// Java: ParserTest.testES2018RegExpFlagS()
func TestES2018RegExpFlagS(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_REGEXP_FLAG_S)
	pt.parse("/a/s")

	pt.mode = Config_LanguageMode_ECMASCRIPT_2015
	pt.expectFeatures(parser.FeatureSet_Feature_REGEXP_FLAG_S)
	pt.parseWarning("/a/s", requiresLanguageModeMessage(parser.FeatureSet_Feature_REGEXP_FLAG_S))
	pt.parseWarning(
		"/a/us", // 'u' added in es6
		requiresLanguageModeMessage(parser.FeatureSet_Feature_REGEXP_FLAG_S))

	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.parseWarning(
		"/a/us", // 'u' added in es6
		requiresLanguageModeMessage(parser.FeatureSet_Feature_REGEXP_FLAG_U),
		requiresLanguageModeMessage(parser.FeatureSet_Feature_REGEXP_FLAG_S))
}

// Java: ParserTest.testES2022RegExpFlagD()
func TestES2022RegExpFlagD(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_REGEXP_FLAG_D)
	pt.parse("/a/d")

	pt.mode = Config_LanguageMode_ECMASCRIPT_2015
	pt.expectFeatures(parser.FeatureSet_Feature_REGEXP_FLAG_D)
	pt.parseWarning("/a/d", requiresLanguageModeMessage(parser.FeatureSet_Feature_REGEXP_FLAG_D))
	pt.parseWarning(
		"/a/ud", // 'u' added in es6
		requiresLanguageModeMessage(parser.FeatureSet_Feature_REGEXP_FLAG_D))

	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.parseWarning(
		"/a/ud", // 'u' added in es6
		requiresLanguageModeMessage(parser.FeatureSet_Feature_REGEXP_FLAG_U),
		requiresLanguageModeMessage(parser.FeatureSet_Feature_REGEXP_FLAG_D))
}

// Java: ParserTest.testDefaultParameters()
func TestDefaultParameters(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY
	pt.expectFeatures(parser.FeatureSet_Feature_DEFAULT_PARAMETERS)
	pt.parse("function f(a, b=0) {}")
	pt.parse("function f(a, b=0, c) {}")

	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.parseWarning("function f(a, b=0) {}", requiresLanguageModeMessage(parser.FeatureSet_Feature_DEFAULT_PARAMETERS))
}

// Java: ParserTest.testDefaultParameterInlineJSDoc()
func TestDefaultParameterInlineJSDoc(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_DEFAULT_PARAMETERS)
	functionNode := pt.parse("function f(/** number */ a = 0) {}").GetFirstChild()
	parameterList := functionNode.GetSecondChild()
	defaultValue := parameterList.GetFirstChild()
	assertNode(t, defaultValue).HasType(rhino.Token_DEFAULT_VALUE)

	aName := defaultValue.GetFirstChild()
	assertNode(t, aName).HasType(rhino.Token_NAME)
	pt.assertNodeHasJSDocInfoWithJSType(aName, pt.NUMBER_TYPE)
}

// Java: ParserTest.testDefaultParameterInlineNonJSDocComment()
func TestDefaultParameterInlineNonJSDocComment(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	pt.expectFeatures(parser.FeatureSet_Feature_DEFAULT_PARAMETERS)
	functionNode := pt.parse("function f(/* number */ a = 0) {}").GetFirstChild()
	parameterList := functionNode.GetSecondChild()
	defaultValue := parameterList.GetFirstChild()
	assertNode(t, defaultValue).HasType(rhino.Token_DEFAULT_VALUE)

	aName := defaultValue.GetFirstChild()
	assertNode(t, aName).HasType(rhino.Token_NAME)
	assertThat(t, aName.GetNonJSDocCommentString()).Contains("/* number */")
}

// Java: ParserTest.testRestParameters()
func TestRestParameters(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY

	pt.parseError("(...xs, x) => xs", "')' expected")
	pt.parseError(
		"function f(...a[0]) {}", "Only an identifier or destructuring pattern is allowed here.")
	pt.parseError("function f(...y, z) {}", "A rest parameter must be last in a parameter list.")

	pt.expectFeatures(parser.FeatureSet_Feature_REST_PARAMETERS)
	pt.parse("function f(...b) {}")
	pt.parse("(...xs) => xs")
	pt.parse("(x, ...xs) => xs")
	pt.parse("(x, y, ...xs) => xs")
}

// Java: ParserTest.testTrailingCommaAfterRestParameters()
func TestTrailingCommaAfterRestParameters(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY

	pt.parseError("function f(...a,) {}", "A trailing comma must not follow a rest parameter.")
	// http://b/184088793, trailing commas should be allowed in arrow function then
	// the proper error is reported.
	pt.parseError("((...xs,) => xs)", "')' expected")
}

// Java: ParserTest.testDestructuredRestParameters()
func TestDestructuredRestParameters(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parseError(
		"function f(...[a[0]]) {}", "Only an identifier or destructuring pattern is allowed here.")

	pt.expectFeatures(parser.FeatureSet_Feature_REST_PARAMETERS, parser.FeatureSet_Feature_ARRAY_DESTRUCTURING)
	pt.parse("(...[x]) => xs")
	pt.parse("(...[x, y]) => xs")
	pt.parse("(a, b, c, ...[x, y, z]) => x")
}

// Java: ParserTest.testRestParameters_ES5()
func TestRestParameters_ES5(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.strictMode = Config_StrictMode_SLOPPY
	pt.expectFeatures(parser.FeatureSet_Feature_REST_PARAMETERS)
	pt.parseWarning("function f(...b) {}", requiresLanguageModeMessage(parser.FeatureSet_Feature_REST_PARAMETERS))
}

// Java: ParserTest.testExpressionsThatLookLikeParameters1()
func TestExpressionsThatLookLikeParameters1(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parseError("();", "invalid parenthesized expression")
	pt.parseError("(...xs);", "invalid parenthesized expression")
	pt.parseError("(x, ...xs);", "A rest parameter must be in a parameter list.")
	pt.parseError("(a, b, c, ...xs);", "A rest parameter must be in a parameter list.")
}

// Java: ParserTest.testExpressionsThatLookLikeParameters2()
func TestExpressionsThatLookLikeParameters2(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parseError("!()", "invalid parenthesized expression")
	pt.parseError("().method", "invalid parenthesized expression")
	pt.parseError("() || a", "invalid parenthesized expression")
	pt.parseError("() && a", "invalid parenthesized expression")
	pt.parseError("x = ()", "invalid parenthesized expression")
}

// Java: ParserTest.testExpressionsThatLookLikeParameters3()
func TestExpressionsThatLookLikeParameters3(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parseError("!(...x)", "invalid parenthesized expression")
	pt.parseError("(...x).method", "invalid parenthesized expression")
	pt.parseError("(...x) || a", "invalid parenthesized expression")
	pt.parseError("(...x) && a", "invalid parenthesized expression")
	pt.parseError("x = (...x)", "invalid parenthesized expression")
}

// Java: ParserTest.testDefaultParametersWithRestParameters()
func TestDefaultParametersWithRestParameters(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parseError("function f(a=1, ...b=3) {}", "A default value cannot be specified after '...'")

	pt.expectFeatures(parser.FeatureSet_Feature_DEFAULT_PARAMETERS, parser.FeatureSet_Feature_REST_PARAMETERS)
	pt.parse("function f(a=0, ...b) {}")
	pt.parse("function f(a, b=0, ...c) {}")
	pt.parse("function f(a, b=0, c=1, ...d) {}")
}

// Java: ParserTest.testClass1()
func TestClass1(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_CLASSES)
	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parse("class C {}")

	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.parseWarning("class C {}", requiresLanguageModeMessage(parser.FeatureSet_Feature_CLASSES))

	pt.mode = Config_LanguageMode_ECMASCRIPT3
	pt.parseWarning("class C {}", requiresLanguageModeMessage(parser.FeatureSet_Feature_CLASSES))
}

// Java: ParserTest.testClass2()
func TestClass2(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_CLASSES)
	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parse("class C {}")

	pt.parse(
		`class C {
  member() {}
  get prop() {}
  set prop(a) {}
}
`)

	pt.parse(
		`class C {
  static member() {}
  static get prop() {}
  static set prop(a) {}
}
`)
}

// Java: ParserTest.testClass3()
func TestClass3(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_CLASSES)
	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parse(
		`class C {
  member() {};
  get prop() {};
  set prop(a) {};
}
`)

	pt.parse(
		`class C {
  static member() {};
  static get prop() {};
  static set prop(a) {};
}
`)
}

// Java: ParserTest.testClassKeywordsAsMethodNames()
func TestClassKeywordsAsMethodNames(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_CLASSES, parser.FeatureSet_Feature_KEYWORDS_AS_PROPERTIES)
	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parse(
		`class KeywordMethods {
  continue() {}
  throw() {}
  else() {}
}
`)
}

// Java: ParserTest.testClassReservedWordsAsMethodNames()
func TestClassReservedWordsAsMethodNames(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_CLASSES, parser.FeatureSet_Feature_KEYWORDS_AS_PROPERTIES)
	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parse(
		`class C {
  import() {};
  get break() {};
  set break(a) {};
}
`)

	pt.parse(
		`class C {
  static import() {};
  static get break() {};
  static set break(a) {};
}
`)
}

// Java: ParserTest.testClass_semicolonsInBodyAreIgnored()
func TestClass_semicolonsInBodyAreIgnored(t *testing.T) {
	pt := newParserTest(t)

	tree :=
		pt.parse(
			`class C {
  foo() {};;;;;;
}
`)

	members := tree.GetFirstChild().GetChildAtIndex(2)
	assertThat(t, members.IsClassMembers()).IsTrue()
	assertThat(t, members.GetChildCount()).IsEqualTo(1)
}

// Java: ParserTest.testClass_constructorMember_legalModifiers()
func TestClass_constructorMember_legalModifiers(t *testing.T) {
	pt := newParserTest(t)

	errorMsg := "Class constructor may not be getter, setter, async, or generator."

	// The expected default case.
	pt.parse("class A { constructor() { } }")

	// Modifiers are legal on a static "constructor" member.
	pt.parse("class A { static constructor() { } }")
	pt.parse("class A { static get constructor() { } }")
	pt.parse("class A { static set constructor(x) { } }")
	pt.parse("class A { static async constructor() { } }")
	pt.parse("class A { static *constructor() { } }")

	// Modifiers are illegal on an instance "constructor" member.
	pt.parseError("class A { get constructor() { } }", errorMsg)
	pt.parseError("class A { set constructor(x) { } }", errorMsg)
	pt.parseError("class A { async constructor() { } }", errorMsg)
	pt.parseError("class A { *constructor() { } }", errorMsg)

	// Modifiers are also illegal on a constructor declared using a string literal.
	// TODO(b/123769080): These should be parse errors, but this case can't be detected currently.
	pt.parse("class A { 'constructor'() { } }")
	pt.parse("class A { get 'constructor'() { } }")
	pt.parse("class A { set 'constructor'(x) { } }")
	pt.parse("class A { async 'constructor'() { } }")
	pt.parse("class A { *'constructor'() { } }")

	// Modifiers are legal on computed properties that happen to be named "constructor".
	pt.parse("class A { ['constructor']() { } }")
	pt.parse("class A { get ['constructor']() { } }")
	pt.parse("class A { set ['constructor'](x) { } }")
	pt.parse("class A { async ['constructor']() { } }")
	pt.parse("class A { *['constructor']() { } }")
}

// Java: ParserTest.testClass_constructorMember_atMostOne()
func TestClass_constructorMember_atMostOne(t *testing.T) {
	pt := newParserTest(t)

	errorMsg := "Class may have only one constructor."

	// The expected default cases.
	pt.parse("class A { }")
	pt.parse("class A { constructor() { } }")

	// Not more than one.
	pt.parseError("class A { constructor() { } constructor() { } }", errorMsg)
	pt.parseError(
		"class A { constructor() { } constructor() { } constructor() { } }", errorMsg, errorMsg)
	// TODO(b/123769080): These should be parse errors, but this case can't be detected currently.
	pt.parse("class A { constructor() { } 'constructor'() { } }")

	// Computed properties can't be the class constructor.
	pt.parse("class A { constructor() { } ['constructor']() { } }")

	// Statics can't be the class constructor.
	pt.parse("class A { constructor() { } static constructor() { } }")
}

// Java: ParserTest.testSingleFieldNoInitializerNoSemicolon()
func TestSingleFieldNoInitializerNoSemicolon(t *testing.T) {
	pt := newParserTest(t)

	n := pt.parse("class C{ field }").GetFirstChild().GetLastChild()
	assertNode(t, n.GetFirstChild()).HasType(rhino.Token_MEMBER_FIELD_DEF)
	assertNode(t, n.GetFirstChild()).HasStringThat().IsEqualTo("field")
}

// Java: ParserTest.testSingleFieldWithInitializerNoSemicolon()
func TestSingleFieldWithInitializerNoSemicolon(t *testing.T) {
	pt := newParserTest(t)

	// check for automatic semicolon insertion
	n := pt.parse("class C{field = 2}").GetFirstChild().GetLastChild().GetFirstChild()
	assertNode(t, n).HasType(rhino.Token_MEMBER_FIELD_DEF)
	assertNode(t, n.GetFirstChild()).IsEqualTo(rhino.IR_Number(2.0))
}

// Java: ParserTest.testSingleFieldNoInitializerWithSemicolon()
func TestSingleFieldNoInitializerWithSemicolon(t *testing.T) {
	pt := newParserTest(t)

	n := pt.parse("class C{ field; }").GetFirstChild().GetLastChild()
	assertNode(t, n.GetFirstChild()).HasStringThat().IsEqualTo("field")
}

// Java: ParserTest.testSingleFieldWithInitializerWithSemicolon()
func TestSingleFieldWithInitializerWithSemicolon(t *testing.T) {
	pt := newParserTest(t)

	n := pt.parse("class C{field = dog;}").GetFirstChild().GetLastChild().GetFirstChild()
	assertNode(t, n).HasType(rhino.Token_MEMBER_FIELD_DEF)
	assertNode(t, n.GetFirstChild()).IsEqualTo(rhino.IR_Name("dog"))
}

// Java: ParserTest.testSingleComputedFieldNoInitializerNoSemicolon()
func TestSingleComputedFieldNoInitializerNoSemicolon(t *testing.T) {
	pt := newParserTest(t)

	// check for automatic semicolon insertion
	n := pt.parse("class C{['field']}").GetFirstChild().GetLastChild().GetFirstChild()
	assertNode(t, n).HasType(rhino.Token_COMPUTED_FIELD_DEF)
	assertNode(t, n.GetFirstChild()).HasStringThat().IsEqualTo("field")
}

// Java: ParserTest.testSingleComputedFieldWithInitializerNoSemicolon()
func TestSingleComputedFieldWithInitializerNoSemicolon(t *testing.T) {
	pt := newParserTest(t)

	// check for automatic semicolon insertion
	n := pt.parse("class C{['field'] = 2}").GetFirstChild().GetLastChild().GetFirstChild()
	assertNode(t, n).HasType(rhino.Token_COMPUTED_FIELD_DEF)
	assertNode(t, n.GetFirstChild()).HasStringThat().IsEqualTo("field")
	assertNode(t, n.GetLastChild()).IsEqualTo(rhino.IR_Number(2.0))
}

// Java: ParserTest.testSingleComputedFieldNoInitializerWithSemicolon()
func TestSingleComputedFieldNoInitializerWithSemicolon(t *testing.T) {
	pt := newParserTest(t)

	n := pt.parse("class C{['field'];}").GetFirstChild().GetLastChild().GetFirstChild()
	assertNode(t, n).HasType(rhino.Token_COMPUTED_FIELD_DEF)
	assertNode(t, n.GetFirstChild()).HasStringThat().IsEqualTo("field")
}

// Java: ParserTest.testSingleComputedFieldWithInitializerWithSemicolon()
func TestSingleComputedFieldWithInitializerWithSemicolon(t *testing.T) {
	pt := newParserTest(t)

	n := pt.parse("class C{['field'] = dog;}").GetFirstChild().GetLastChild().GetFirstChild()
	assertNode(t, n).HasType(rhino.Token_COMPUTED_FIELD_DEF)
	assertNode(t, n.GetFirstChild()).HasStringThat().IsEqualTo("field")
	assertNode(t, n.GetLastChild()).IsEqualTo(rhino.IR_Name("dog"))
}

// Java: ParserTest.testMultipleFieldsNoInitializerNoLineBreakNoSemicolon()
func TestMultipleFieldsNoInitializerNoLineBreakNoSemicolon(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("class C{ a b c }", "Semi-colon expected")
}

// Java: ParserTest.testMultipleFieldsNoInitializerWithLineBreakNoSemicolon()
func TestMultipleFieldsNoInitializerWithLineBreakNoSemicolon(t *testing.T) {
	pt := newParserTest(t)

	n := pt.parse("class C{ a \n b \n c }").GetFirstChild().GetLastChild()
	assertNode(t, n.GetFirstChild()).HasType(rhino.Token_MEMBER_FIELD_DEF)
	assertNode(t, n.GetFirstChild()).HasStringThat().IsEqualTo("a")
	assertNode(t, n.GetSecondChild()).HasType(rhino.Token_MEMBER_FIELD_DEF)
	assertNode(t, n.GetSecondChild()).HasStringThat().IsEqualTo("b")
	assertNode(t, n.GetLastChild()).HasType(rhino.Token_MEMBER_FIELD_DEF)
	assertNode(t, n.GetLastChild()).HasStringThat().IsEqualTo("c")
}

// Java: ParserTest.testMultipleFieldsNoLineBreakNoSemicolon()
func TestMultipleFieldsNoLineBreakNoSemicolon(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("class C{ a = 2 b = 4 c = 5}", "Semi-colon expected")
}

// Java: ParserTest.testMultipleFieldsWithLineBreakNoSemicolon()
func TestMultipleFieldsWithLineBreakNoSemicolon(t *testing.T) {
	pt := newParserTest(t)

	// check for automatic semicolon insertion
	n := pt.parse("class C{ a = 2 \n b = 4 \n c = 5}").GetFirstChild().GetLastChild()
	assertNode(t, n.GetFirstChild()).HasType(rhino.Token_MEMBER_FIELD_DEF)
	assertNode(t, n.GetFirstChild()).HasStringThat().IsEqualTo("a")
	assertNode(t, n.GetFirstFirstChild()).IsEqualTo(rhino.IR_Number(2.0))
	assertNode(t, n.GetSecondChild()).HasType(rhino.Token_MEMBER_FIELD_DEF)
	assertNode(t, n.GetSecondChild()).HasStringThat().IsEqualTo("b")
	assertNode(t, n.GetSecondChild().GetFirstChild()).IsEqualTo(rhino.IR_Number(4.0))
	assertNode(t, n.GetLastChild()).HasType(rhino.Token_MEMBER_FIELD_DEF)
	assertNode(t, n.GetLastChild()).HasStringThat().IsEqualTo("c")
	assertNode(t, n.GetLastChild().GetFirstChild()).IsEqualTo(rhino.IR_Number(5.0))
}

// Java: ParserTest.testMultipleFieldsNoLineBreakWithSemicolon()
func TestMultipleFieldsNoLineBreakWithSemicolon(t *testing.T) {
	pt := newParserTest(t)

	n := pt.parse("class C{field = 2; hi = 3;}").GetFirstChild().GetLastChild()
	assertNode(t, n.GetFirstChild()).HasType(rhino.Token_MEMBER_FIELD_DEF)
	assertNode(t, n.GetFirstChild()).HasStringThat().IsEqualTo("field")
	assertNode(t, n.GetFirstFirstChild()).IsEqualTo(rhino.IR_Number(2.0))
	assertNode(t, n.GetSecondChild()).HasStringThat().IsEqualTo("hi")
	assertNode(t, n.GetSecondChild().GetFirstChild()).IsEqualTo(rhino.IR_Number(3.0))
}

// Java: ParserTest.testMultipleFieldsWithLineBreakWithSemicolon()
func TestMultipleFieldsWithLineBreakWithSemicolon(t *testing.T) {
	pt := newParserTest(t)

	n := pt.parse("class C{field = 2; \n hi = 3;}").GetFirstChild().GetLastChild()
	assertNode(t, n.GetFirstChild()).HasType(rhino.Token_MEMBER_FIELD_DEF)
	assertNode(t, n.GetFirstChild()).HasStringThat().IsEqualTo("field")
	assertNode(t, n.GetFirstFirstChild()).IsEqualTo(rhino.IR_Number(2.0))
	assertNode(t, n.GetSecondChild()).HasStringThat().IsEqualTo("hi")
	assertNode(t, n.GetSecondChild().GetFirstChild()).IsEqualTo(rhino.IR_Number(3.0))
}

// Java: ParserTest.testMultipleComputedFieldsNoLineBreakNoSemicolon()
func TestMultipleComputedFieldsNoLineBreakNoSemicolon(t *testing.T) {
	pt := newParserTest(t)

	// tree is parsed correctly even though it looks weird
	n := pt.parse("class C{ ['a'] = 2 ['b'] = 4 ['c'] = 5}").GetFirstChild().GetLastChild()
	assertNode(t, n.GetFirstChild()).HasType(rhino.Token_COMPUTED_FIELD_DEF)
}

// Java: ParserTest.testMultipleComputedFieldsWithLineBreakNoSemicolon()
func TestMultipleComputedFieldsWithLineBreakNoSemicolon(t *testing.T) {
	pt := newParserTest(t)

	// tree is parsed correctly even though it looks weird
	n := pt.parse("class C{ ['a'] = 2 \n ['b'] = 4 \n ['c'] = 5}").GetFirstChild().GetLastChild()
	assertNode(t, n.GetFirstChild()).HasType(rhino.Token_COMPUTED_FIELD_DEF)
}

// Java: ParserTest.testMultipleComputedFieldsNoLineBreakWithSemicolon()
func TestMultipleComputedFieldsNoLineBreakWithSemicolon(t *testing.T) {
	pt := newParserTest(t)

	n := pt.parse("class C{['field'] = 2; ['hi'] = 3;}").GetFirstChild().GetLastChild()
	assertNode(t, n.GetFirstChild()).HasType(rhino.Token_COMPUTED_FIELD_DEF)
	assertNode(t, n.GetFirstFirstChild()).HasStringThat().IsEqualTo("field")
	assertNode(t, n.GetFirstChild().GetLastChild()).IsEqualTo(rhino.IR_Number(2.0))
	assertNode(t, n.GetSecondChild()).HasType(rhino.Token_COMPUTED_FIELD_DEF)
	assertNode(t, n.GetSecondChild().GetFirstChild()).HasStringThat().IsEqualTo("hi")
	assertNode(t, n.GetSecondChild().GetLastChild()).IsEqualTo(rhino.IR_Number(3.0))
}

// Java: ParserTest.testMultipleComputedFieldsWithLineBreaksWithSemicolons()
func TestMultipleComputedFieldsWithLineBreaksWithSemicolons(t *testing.T) {
	pt := newParserTest(t)

	n := pt.parse("class C{['field'] = 2; \n ['hi'] = 3;}").GetFirstChild().GetLastChild()
	assertNode(t, n.GetFirstChild()).HasType(rhino.Token_COMPUTED_FIELD_DEF)
	assertNode(t, n.GetFirstFirstChild()).HasStringThat().IsEqualTo("field")
	assertNode(t, n.GetFirstChild().GetLastChild()).IsEqualTo(rhino.IR_Number(2.0))
	assertNode(t, n.GetSecondChild()).HasType(rhino.Token_COMPUTED_FIELD_DEF)
	assertNode(t, n.GetSecondChild().GetFirstChild()).HasStringThat().IsEqualTo("hi")
	assertNode(t, n.GetSecondChild().GetLastChild()).IsEqualTo(rhino.IR_Number(3.0))
}

// Java: ParserTest.testMultipleMixedFieldsNoLineBreakWithSemicolon()
func TestMultipleMixedFieldsNoLineBreakWithSemicolon(t *testing.T) {
	pt := newParserTest(t)

	n := pt.parse("class C{  b = 4; ['a'] = 2 }").GetFirstChild().GetLastChild()
	assertNode(t, n.GetFirstChild()).HasType(rhino.Token_MEMBER_FIELD_DEF)
	assertNode(t, n.GetSecondChild()).HasType(rhino.Token_COMPUTED_FIELD_DEF)
}

// Java: ParserTest.testMultipleMixedFieldsNoLineBreakNoSemicolon()
func TestMultipleMixedFieldsNoLineBreakNoSemicolon(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("class C{ ['a'] = 2 b = 4}", "Semi-colon expected")

	// tree is parsed correctly even though it looks weird
	n := pt.parse("class C{b = 4 ['a'] = 2 }").GetFirstChild().GetLastChild()
	assertNode(t, n.GetFirstChild()).HasType(rhino.Token_MEMBER_FIELD_DEF)
}

// Java: ParserTest.testMultipleMixedFieldsWithLineBreakNoSemicolon()
func TestMultipleMixedFieldsWithLineBreakNoSemicolon(t *testing.T) {
	pt := newParserTest(t)

	// tree is parsed correctly even though it looks weird
	n := pt.parse("class C{b = 4 \n ['a'] = 2 }").GetFirstChild().GetLastChild()
	assertNode(t, n.GetFirstChild()).HasType(rhino.Token_MEMBER_FIELD_DEF)
}

// Java: ParserTest.testComputedFieldStringLit()
func TestComputedFieldStringLit(t *testing.T) {
	pt := newParserTest(t)

	n := pt.parse("class C {'x' = 2;}").GetFirstChild().GetLastChild().GetFirstChild()
	assertNode(t, n).HasType(rhino.Token_COMPUTED_FIELD_DEF)
	assertNode(t, n.GetFirstChild()).HasType(rhino.Token_STRINGLIT)
	assertNode(t, n.GetFirstChild()).HasStringThat().IsEqualTo("x")
	assertNode(t, n.GetLastChild()).HasType(rhino.Token_NUMBER)
	assertNode(t, n.GetLastChild()).IsEqualTo(rhino.IR_Number(2.0))
}

// Java: ParserTest.testComputedFieldName()
func TestComputedFieldName(t *testing.T) {
	pt := newParserTest(t)

	n := pt.parse("class C {[a] = 2;}").GetFirstChild().GetLastChild().GetFirstChild()
	assertNode(t, n).HasType(rhino.Token_COMPUTED_FIELD_DEF)
	assertNode(t, n.GetFirstChild()).HasType(rhino.Token_NAME)
	assertNode(t, n.GetFirstChild()).HasStringThat().IsEqualTo("a")
	assertNode(t, n.GetLastChild()).HasType(rhino.Token_NUMBER)
	assertNode(t, n.GetLastChild()).IsEqualTo(rhino.IR_Number(2.0))
}

// Java: ParserTest.testComputedFieldNumber()
func TestComputedFieldNumber(t *testing.T) {
	pt := newParserTest(t)

	n := pt.parse("class C{1 = 2;}").GetFirstChild().GetLastChild().GetFirstChild()
	assertNode(t, n).HasType(rhino.Token_COMPUTED_FIELD_DEF)
	assertNode(t, n.GetFirstChild()).HasType(rhino.Token_NUMBER)
	assertNode(t, n.GetFirstChild()).IsEqualTo(rhino.IR_Number(1.0))
	assertNode(t, n.GetLastChild()).HasType(rhino.Token_NUMBER)
	assertNode(t, n.GetLastChild()).IsEqualTo(rhino.IR_Number(2.0))
}

// Java: ParserTest.testStaticClassFields()
func TestStaticClassFields(t *testing.T) {
	pt := newParserTest(t)

	n := pt.parse("class C{static field = 2;}").GetFirstChild().GetLastChild()
	assertNode(t, n.GetFirstChild()).IsStatic()

	n = pt.parse("class C{static field = 2; hi = 3;}").GetFirstChild().GetLastChild()
	assertNode(t, n.GetFirstChild()).IsStatic()

	n = pt.parse("class C{static field = 2; static hi = 3;}").GetFirstChild().GetLastChild()
	assertNode(t, n.GetFirstChild()).IsStatic()
	assertNode(t, n.GetSecondChild()).IsStatic()
}

// Java: ParserTest.testStaticClassComputedFields()
func TestStaticClassComputedFields(t *testing.T) {
	pt := newParserTest(t)

	n := pt.parse("class C{static ['field'] = 2;}").GetFirstChild().GetLastChild()
	assertNode(t, n.GetFirstChild()).IsStatic()

	n = pt.parse("class C{static ['field'] = 2; ['hi'] = 3;}").GetFirstChild().GetLastChild()
	assertNode(t, n.GetFirstChild()).IsStatic()

	n = pt.parse("class C{static ['field'] = 2; static ['hi'] = 3;}").GetFirstChild().GetLastChild()
	assertNode(t, n.GetFirstChild()).IsStatic()
	assertNode(t, n.GetSecondChild()).IsStatic()

	n = pt.parse("class C{static 'hi' = 2;}").GetFirstChild().GetLastChild()
	assertNode(t, n.GetFirstChild()).IsStatic()

	n = pt.parse("class C{static 1 = 2;}").GetFirstChild().GetLastChild()
	assertNode(t, n.GetFirstChild()).IsStatic()
}

// Java: ParserTest.testClassField_es2020()
func TestClassField_es2020(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT_2020
	pt.expectFeatures(parser.FeatureSet_Feature_PUBLIC_CLASS_FIELDS)
	pt.parseWarning("class C{field = 2;}", requiresLanguageModeMessage(parser.FeatureSet_Feature_PUBLIC_CLASS_FIELDS))
}

// Java: ParserTest.testClassComputedField_es2020()
func TestClassComputedField_es2020(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT_2020
	pt.expectFeatures(parser.FeatureSet_Feature_PUBLIC_CLASS_FIELDS)
	pt.parseWarning("class C{['a'] = 2;}", requiresLanguageModeMessage(parser.FeatureSet_Feature_PUBLIC_CLASS_FIELDS))
}

// Java: ParserTest.testClassComputedField_es2020_noWarningOrRecordingOfFeaturesInClosureUnawareCode()
func TestClassComputedField_es2020_noWarningOrRecordingOfFeaturesInClosureUnawareCode(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT_2020
	pt.expectFeatures(parser.FeatureSet_Feature_PUBLIC_CLASS_FIELDS)
	pt.parseWarning(
		`/** @fileoverview @closureUnaware */
/** @closureUnaware */ (
  function() { class C{ a = 2;} }).call(globalThis);
`,
		// This warning is expected - it is used to halt compilation during parsing when
		// closure-unaware code is present when it should not be.
		"@closureUnaware annotation is not allowed in this compilation",
		"@closureUnaware annotation is not allowed in this compilation")

	// Note that the PUBLIC_CLASS_FIELDS feature is not recorded here. This is OK as we currently do
	// not support transpiling closure-unaware code.
	pt.expectFeatures(parser.FeatureSet_Feature_CLASSES)
}

// Java: ParserTest.testClassComputedField_es2020_warnsForCodeOutsideClosureUnawareRange()
func TestClassComputedField_es2020_warnsForCodeOutsideClosureUnawareRange(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT_2020
	pt.expectFeatures(parser.FeatureSet_Feature_PUBLIC_CLASS_FIELDS)
	pt.parseWarning(
		`/** @fileoverview @closureUnaware */
/** @closureUnaware */ (
  function() {
    class C{ a = 2;}
}).call(globalThis);
class C{
  a = 2;
}
`,
		"@closureUnaware annotation is not allowed in this compilation",
		"This language feature is only supported for ECMASCRIPT_2022 mode or better: Public class fields",
		"@closureUnaware annotation is not allowed in this compilation")

	pt.expectFeatures(parser.FeatureSet_Feature_CLASSES, parser.FeatureSet_Feature_PUBLIC_CLASS_FIELDS)
}

// Java: ParserTest.testPrivateProperty_unstable()
func TestPrivateProperty_unstable(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_UNSTABLE
	pt.expectFeatures(parser.FeatureSet_Feature_PRIVATE_ELEMENTS)
	pt.parseWarning(
		"class C { #f = 2; }", requiresLanguageModeMessage(parser.FeatureSet_Feature_PRIVATE_ELEMENTS))
}

// Java: ParserTest.testPrivateProperty_singleClassMember()
func TestPrivateProperty_singleClassMember(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_PRIVATE_ELEMENTS)

	pt.parse("class C { #f; }")
	pt.parse("class C { #m() {} }")
	pt.parse("class C { *#g() {} }")
	pt.parse("class C { get #g() {} }")
	pt.parse("class C { set #s(x) {} }")
	pt.parse("class C { get #p() {} set #p(x) {} }")
	pt.parse("class C { async #a() {} }")
	pt.parse("class C { async *#ag() {} }")

	pt.parse("class C { static #sf; }")
	pt.parse("class C { static #sm() {} }")
	pt.parse("class C { static *#sg() {} }")
	pt.parse("class C { static get #sg() {} }")
	pt.parse("class C { static set #ss(x) {} }")
	pt.parse("class C { static get #sp() {} static set #sp(x) {} }")
	pt.parse("class C { static async #sa() {} }")
	pt.parse("class C { static async *#sag() {} }")
}

// Java: ParserTest.testPrivateProperty_duplicateIdentifierOnInstanceAndStatic()
func TestPrivateProperty_duplicateIdentifierOnInstanceAndStatic(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_PRIVATE_ELEMENTS)

	pt.parseError("class C { #x; static #x; }", "Identifier '#x' has already been declared")
	pt.parseError("class C { #x; static [#x]; }", parserTest_INVALID_PRIVATE_ID)
}

// Java: ParserTest.testPrivateProperty_definition_linenocharno()
func TestPrivateProperty_definition_linenocharno(t *testing.T) {
	pt := newParserTest(t)

	n :=
		pt.parse(
			`class C {
  #pf = 1;
  #pm() {}
}
`).
			GetFirstChild()

	members := nodeUtil_getClassMembers(n)

	privateField := members.GetFirstChild()
	assertThat(t, privateField.GetLineno()).IsEqualTo(2)
	assertThat(t, privateField.GetCharno()).IsEqualTo(2)
	assertThat(t, privateField.GetLength()).IsEqualTo(8) // Includes the assignment

	privateMethod := members.GetLastChild()
	assertThat(t, privateMethod.GetLineno()).IsEqualTo(3)
	assertThat(t, privateMethod.GetCharno()).IsEqualTo(2)
	assertThat(t, privateMethod.GetLength()).IsEqualTo(3) // Just the method name.
}

// Java: ParserTest.testPrivateProperty_multipleClassMembers()
func TestPrivateProperty_multipleClassMembers(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("class C { #f; #g; }")
	pt.parse("class C { #m() {} #n() {} }")
	pt.parse("class C { get #g() {} get #h() {} }")
	pt.parse("class C { set #s(x) {} set #t(x) {} }")
	pt.parse("class C { get #s() {} set #s(x) {} get #t() {} set #t(x) {} }")

	pt.parse("class C { static #sf; static #sg; }")
	pt.parse("class C { static #sm() {} static #sn() {} }")
	pt.parse("class C { static get #sg() {} static get #sh() {} }")
	pt.parse("class C { static set #ss(x) {} static set #st(x) {} }")
	pt.parse(
		`class C { static get #ss() {} static set #ss(x) {}
static get #st() {} static set #st(x) {} }
`)

	pt.parse("class C { #a; #b; c() {} d() {} get #e() {} set #e(x) {} set #f(x) {} }")
	pt.parse("class C { static #a; #b; static c() {} d() {} static get #e() {} set #f(x) {} }")
	pt.parse(
		`class C { static #a; static #b; static c() {} static d() {}
static get #e() {} static set #f(x) {} }
`)
}

// Java: ParserTest.testPrivateProperty_invalid_identifierAlreadyDeclared()
func TestPrivateProperty_invalid_identifierAlreadyDeclared(t *testing.T) {
	pt := newParserTest(t)

	expectedError := "Identifier '#p' has already been declared"

	pt.parseError("class C { #p; #p; }", expectedError)
	pt.parseError("class C { #p() {} #p() {} }", expectedError)
	pt.parseError("class C { get #p() {} get #p() {} }", expectedError)
	pt.parseError("class C { set #p(x) {} set #p(x) {} }", expectedError)

	pt.parseError("class C { static #p; static #p; }", expectedError)
	pt.parseError("class C { static #p() {} static #p() {} }", expectedError)
	pt.parseError("class C { static get #p() {} static get #p() {} }", expectedError)
	pt.parseError("class C { static set #p(x) {} static set #p(x) {} }", expectedError)

	pt.parseError("class C { #p; #p() {} }", expectedError)
	pt.parseError("class C { #p; get #p() {} }", expectedError)
	pt.parseError("class C { #p; set #p(x) {} }", expectedError)

	pt.parseError("class C { #p() {} #p; }", expectedError)
	pt.parseError("class C { #p() {} get #p() {} }", expectedError)
	pt.parseError("class C { #p() {} set #p(x) {} }", expectedError)

	pt.parseError("class C { get #p() {} #p; }", expectedError)
	pt.parseError("class C { get #p() {} #p() {} }", expectedError)
	pt.parse( /**/ "class C { get #p() {} set #p(x) {} }") // OK

	pt.parseError("class C { set #p(x) {} #p; }", expectedError)
	pt.parseError("class C { set #p(x) {} #p() {} }", expectedError)
	pt.parse( /**/ "class C { set #p(x) {} get #p() {} }") // OK

	pt.parseError("class C { #p; static #p; }", expectedError) // Cross-static duplicate field
	pt.parseError("class C { #p() {} static #p() {} }", expectedError)
	pt.parseError("class C { get #p() {} static get #p() {} }", expectedError)
	pt.parseError("class C { set #p(x) {} static set #p(x) {} }", expectedError)
}

// Java: ParserTest.testPrivateProperty_invalid_identifierAlreadyDeclared_crossStaticGetterSetter()
func TestPrivateProperty_invalid_identifierAlreadyDeclared_crossStaticGetterSetter(t *testing.T) {
	pt := newParserTest(t)

	expectedError := "Identifier '#p' has already been declared"

	// Repeating in all orders as the check is unique for each to a degree.
	pt.parseError("class C { get #p() {} static set #p(x) {} }", expectedError)
	pt.parseError("class C { set #p(x) {} static get #p() {} }", expectedError)
	pt.parseError("class C { static get #p() {} set #p(x) {} }", expectedError)
	pt.parseError("class C { static set #p(x) {} get #p() {} }", expectedError)
}

// Java: ParserTest.testPrivateProperty_classMembersReferencingOtherPrivateField()
func TestPrivateProperty_classMembersReferencingOtherPrivateField(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_PRIVATE_ELEMENTS)

	pt.parse("class C { #f = 1; f = this.#f; }")

	pt.parse("class C { static #sf = 1; static sf = this.#sf; }")
	pt.parse("class C { static #sf = 1; static sf = C.#sf; }")
	pt.parse("class C { static #sf; static { this.#sf = 1; } }")
	pt.parse("class C { static #sf; static { C.#sf = 1; } }")
}

// Java: ParserTest.testPrivateProperty_reference_linenocharno()
func TestPrivateProperty_reference_linenocharno(t *testing.T) {
	pt := newParserTest(t)

	n :=
		pt.parse(
			`class C {
  #pf1 = 1;
  #pf2 = this.#pf1;
}
`).
			GetFirstChild()

	members := nodeUtil_getClassMembers(n)

	field2 := members.GetLastChild()
	assertNode(t, field2).HasType(rhino.Token_MEMBER_FIELD_DEF)

	field2GetProp := field2.GetFirstChild()
	assertNode(t, field2GetProp).HasType(rhino.Token_GETPROP)
	assertNode(t, field2GetProp).HasStringThat().IsEqualTo("#pf1")
	assertThat(t, field2GetProp.GetLineno()).IsEqualTo(3)
	assertThat(t, field2GetProp.GetCharno()).IsEqualTo(14)
	assertThat(t, field2GetProp.GetLength()).IsEqualTo(4)
}

// Java: ParserTest.testPrivateProperty_classMethodsReferencingOtherPrivateProp()
func TestPrivateProperty_classMethodsReferencingOtherPrivateProp(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_PRIVATE_ELEMENTS)

	pt.parse("class C { #f = 1; method() { this.#f = 2; } }")
	pt.parse("class C { #f = 1; #g = 2; method() { this.#f = this.#g; } }")
	pt.parse("class C { #f = 1; method() { const t = this; t.#f = 2; } }")
	pt.parse("class C { #f = 1; method() { const f = () => { this.#f = 2; }; } }")
	pt.parse("class C { #f = 1; method() { const x = [this.#f]; } }")
	pt.parse("class C { #f = 1; method() { const x = {y: this.#f}; } }")
	pt.parse("class C { #pm() { this.#pm(); } }")

	pt.parse("class C { static #f = 1; static method() { this.#f = 2; } }")
	pt.parse("class C { static #f = 1; static method() { C.#f = 2; } }")
	pt.parse("class C { static #f = 1; static method() { const x = [this.#f]; } }")
	pt.parse("class C { static #f = 1; static method() { const x = {y: C.#f}; } }")
	pt.parse("class C { static #pm() { this.#pm(); } }")
	pt.parse("class C { static #pm() { C.#pm(); } }")
}

// Java: ParserTest.testPrivateProperty_nestedClasses()
func TestPrivateProperty_nestedClasses(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_PRIVATE_ELEMENTS)

	pt.parse(
		`class Outer {
  #of1;
  #of2;
  om() {
    this.#of1;
    this.#of2;
    class Inner {
      #if1;
      #if2;
      im() {
        this.#of1;
        this.#of2;
        this.#if1;
        this.#if2;
      }
    }
  }
}
`)
}

// Java: ParserTest.testPrivateProperty_nestedClasses_invalid_referenceInnerFieldFromOuter()
func TestPrivateProperty_nestedClasses_invalid_referenceInnerFieldFromOuter(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError(
		`class Outer {
  #of1;
  #of2;
  om() {
    this.#of1;
    this.#of2;
    this.#if1; // Invalid
    this.#if2; // Invalid
    class Inner {
      #if1;
      #if2;
      im() {
        this.#of1;
        this.#of2;
        this.#if1;
        this.#if2;
      }
    }
  }
}
`,
		parserTest_PRIVATE_FIELD_NOT_DEFINED,
		parserTest_PRIVATE_FIELD_NOT_DEFINED)
}

// Java: ParserTest.testPrivateProperty_classMethodsReferencingOtherPrivatePropViaOptionalChain()
func TestPrivateProperty_classMethodsReferencingOtherPrivatePropViaOptionalChain(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_PRIVATE_ELEMENTS)

	pt.parse("class C { #f = 1; method() { const t = this; t?.#f; } }")
	pt.parse("class C { #pm() { const t = this; t?.#pm(); } }")
}

// Java: ParserTest.testPrivateProperty_destructuredAssignmentDefaultValue()
func TestPrivateProperty_destructuredAssignmentDefaultValue(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("class C { #px = 1; method(x) { const { y = this.#px } = x; } }")
	pt.parse("class C { #px = 1; method(x) { const { y: z = this.#px } = x; } }")
}

// Java: ParserTest.testPrivateProperty_invalid_nonExistentPrivateProp()
func TestPrivateProperty_invalid_nonExistentPrivateProp(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("class C { #f = this.#missing; }", parserTest_PRIVATE_FIELD_NOT_DEFINED)
	pt.parseError("class C { method() { this.#missing = 1; } }", parserTest_PRIVATE_FIELD_NOT_DEFINED)
	pt.parseError(
		"class C { #f = 1; method() { this.#f = this.#missing; } }", parserTest_PRIVATE_FIELD_NOT_DEFINED)
	pt.parseError("class C { method() { const t = this; t.#missing; } }", parserTest_PRIVATE_FIELD_NOT_DEFINED)
	pt.parseError("class C { method() { const t = this; t?.#missing; } }", parserTest_PRIVATE_FIELD_NOT_DEFINED)
	pt.parseError(
		"class C { method() { const f = () => { this.#missing; }; } }", parserTest_PRIVATE_FIELD_NOT_DEFINED)
	pt.parseError(
		"class C { method() { class D { method() { this.#missing = 1; } } } }",
		parserTest_PRIVATE_FIELD_NOT_DEFINED)
	pt.parseError(
		"class C { method() { class D { method() { const x = this.#missing; } } } }",
		parserTest_PRIVATE_FIELD_NOT_DEFINED)

	pt.parseError("class C { method() { this.#missing(); } }", parserTest_PRIVATE_METHOD_NOT_DEFINED)

	pt.parseError("class C { static #f = this.#missing; }", parserTest_PRIVATE_FIELD_NOT_DEFINED)
	pt.parseError("class C { static #f = C.#missing; }", parserTest_PRIVATE_FIELD_NOT_DEFINED)

	pt.parseError("class C { static { this.#missing; } }", parserTest_PRIVATE_FIELD_NOT_DEFINED)
	pt.parseError("class C { static { C.#missing; } }", parserTest_PRIVATE_FIELD_NOT_DEFINED)

	pt.parseError("class C { static { this.#missing(); } }", parserTest_PRIVATE_METHOD_NOT_DEFINED)
	pt.parseError("class C { static { C.#missing(); } }", parserTest_PRIVATE_METHOD_NOT_DEFINED)
}

// Java: ParserTest.testPrivateProperty_invalid_subclassAccessingBasePrivateField()
func TestPrivateProperty_invalid_subclassAccessingBasePrivateField(t *testing.T) {
	pt := newParserTest(t)

	// We expect this to fail because private fields in a base class cannot be accessed by
	// subclasses.
	// The JS parser correctly rejects this with a parse error as #baseField is not declared in Sub.
	pt.parseError(
		`class Base {
  #baseField = 1;
}
class Sub extends Base {
  constructor() {
    super();
    this.#baseField = 2;
  }
}
`,
		parserTest_PRIVATE_FIELD_NOT_DEFINED)
}

// Java: ParserTest.testPrivateProperty_invalid_deletePrivateField()
func TestPrivateProperty_invalid_deletePrivateField(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("class C { #f = 1; method() { delete this.#f; } }", parserTest_PRIVATE_FIELD_DELETED)
	pt.parseError(
		"class C { #f = 1; method() { const t = this; delete t.#f; } }", parserTest_PRIVATE_FIELD_DELETED)
	pt.parseError(
		"class C { #f = 1; method() { const t = this; delete t?.#f; } }", parserTest_PRIVATE_FIELD_DELETED)
	pt.parseError(
		"class C { #f = 1; method() { const a = {b: this}; delete ((a.b).#f); } }",
		parserTest_PRIVATE_FIELD_DELETED)

	pt.parseError(
		"class C { static #f = 1; static method() { delete this.#f; } }", parserTest_PRIVATE_FIELD_DELETED)
	pt.parseError(
		"class C { static #f = 1; static method() { delete C.#f; } }", parserTest_PRIVATE_FIELD_DELETED)
	pt.parseError("class C { static #f = 1; static { delete this.#f; } }", parserTest_PRIVATE_FIELD_DELETED)
	pt.parseError("class C { static #f = 1; static { delete C.#f; } }", parserTest_PRIVATE_FIELD_DELETED)
}

// Java: ParserTest.testPrivateProperty_invalid_deleteUndeclaredPrivateField()
func TestPrivateProperty_invalid_deleteUndeclaredPrivateField(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError(
		"class C { method() { delete this.#f; } }",
		parserTest_PRIVATE_FIELD_DELETED,
		parserTest_PRIVATE_FIELD_NOT_DEFINED)
	pt.parseError(
		"class C { method() { const t = this; delete t.#f; } }",
		parserTest_PRIVATE_FIELD_DELETED,
		parserTest_PRIVATE_FIELD_NOT_DEFINED)
	pt.parseError(
		"class C { method() { const t = this; delete t?.#f; } }",
		parserTest_PRIVATE_FIELD_DELETED,
		parserTest_PRIVATE_FIELD_NOT_DEFINED)
	pt.parseError(
		"class C { method() { const a = {b: this}; delete ((a.b).#f); } }",
		parserTest_PRIVATE_FIELD_DELETED,
		parserTest_PRIVATE_FIELD_NOT_DEFINED)

	pt.parseError(
		"class C { static method() { delete this.#f; } }",
		parserTest_PRIVATE_FIELD_DELETED,
		parserTest_PRIVATE_FIELD_NOT_DEFINED)
	pt.parseError(
		"class C { static method() { delete C.#f; } }",
		parserTest_PRIVATE_FIELD_DELETED,
		parserTest_PRIVATE_FIELD_NOT_DEFINED)
	pt.parseError(
		"class C { static { delete this.#f; } }", parserTest_PRIVATE_FIELD_DELETED, parserTest_PRIVATE_FIELD_NOT_DEFINED)
	pt.parseError(
		"class C { static { delete C.#f; } }", parserTest_PRIVATE_FIELD_DELETED, parserTest_PRIVATE_FIELD_NOT_DEFINED)
}

// Java: ParserTest.testPrivateProperty_invalid_objectLiteralProperty()
func TestPrivateProperty_invalid_objectLiteralProperty(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("const x = { #pf: 1 }", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("const x = { #pm() {} }", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("const x = { get #pp() {} }", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("const x = { set #ps(x) {} }", parserTest_INVALID_PRIVATE_ID)

	pt.parseError("class C { method() { const x = { #pf: 1 }; } }", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("class C { method() { const x = { #pm() {} }; } }", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("class C { method() { const x = { get #pp() {} }; } }", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("class C { method() { const x = { set #ps(x) {} }; } }", parserTest_INVALID_PRIVATE_ID)
}

// Java: ParserTest.testPrivateProperty_invalid_referencePrivateOnSuper()
func TestPrivateProperty_invalid_referencePrivateOnSuper(t *testing.T) {
	pt := newParserTest(t)

	// Private field read and write on super
	pt.parseError("class C { #f = 1; m() { return super.#f; } }", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("class C { #f = 1; m(v) { super.#f = v; } }", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("class C { #f = 1; m() { super.#f += 1; } }", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("class C { #f = 1; m() { super.#f++; } }", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("class C { #f = 1; m() { ++super.#f; } }", parserTest_INVALID_PRIVATE_ID)

	// Private method call on super
	pt.parseError("class C { #m() {} test() { super.#m(); } }", parserTest_INVALID_PRIVATE_ID)

	// Private getter/setter on super
	pt.parseError("class C { get #g() {} m() { return super.#g; } }", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("class C { set #s(v) {} m(v) { super.#s = v; } }", parserTest_INVALID_PRIVATE_ID)

	// Optional chaining with super private property
	pt.parseError(
		"class C { #f = 1; m() { super?.#f; } }", "Optional chaining is forbidden in super?.")

	// Destructuring assignment to super private property
	pt.parseError("class C { #f = 1; m() { [super.#f] = [1]; } }", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("class C { #f = 1; m() { ({ x: super.#f } = { x: 1 }); } }", parserTest_INVALID_PRIVATE_ID)

	// Nested class referencing outer class private method/field on super
	pt.parseError(
		`class Base {
  #priv() {}
  #field = 1;
  static {
    class Sub extends Base {
      m() {
        return super.#priv();
      }
      getF() {
        return super.#field;
      }
    }
  }
}
`,
		parserTest_INVALID_PRIVATE_ID,
		parserTest_INVALID_PRIVATE_ID)
}

// Java: ParserTest.testPrivateProperty_invalid_destructuredAssignment()
func TestPrivateProperty_invalid_destructuredAssignment(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("const { #px } = x;", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("const { x: #px } = x;", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("const { x = #px } = x;", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("const { x: y = #px } = x;", parserTest_INVALID_PRIVATE_ID)

	pt.parseError("class C { method() { const { #px } = x; } }", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("class C { method() { const { x: #px } = x; } }", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("class C { method() { const { x = #px } = x; } }", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("class C { method() { const { x: y = #px } = x; } }", parserTest_INVALID_PRIVATE_ID)
}

// Java: ParserTest.testPrivateProperty_invalid_variableName()
func TestPrivateProperty_invalid_variableName(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("const #pv = 1;", parserTest_INVALID_PRIVATE_ID)

	pt.parseError("class C { method() { const #pv = 1; } }", parserTest_INVALID_PRIVATE_ID)
}

// Java: ParserTest.testPrivateProperty_invalid_functionName()
func TestPrivateProperty_invalid_functionName(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("function #pf() {}", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("function* #pf() {}", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("async function #pf() {}", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("async function* #pf() {}", parserTest_INVALID_PRIVATE_ID)

	pt.parseError("class C { method() { function #pf() {} } }", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("class C { method() { function* #pf() {} } }", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("class C { method() { async function #pf() {} } }", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("class C { method() { async function* #pf() {} } }", parserTest_INVALID_PRIVATE_ID)
}

// Java: ParserTest.testPrivateProperty_invalid_paramName()
func TestPrivateProperty_invalid_paramName(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("function f(#p) {}", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("class C { method(#p) {} }", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("class C { #p; method(#p) {} }", parserTest_INVALID_PRIVATE_ID)

	pt.parseError("class C { method() { function f(#p) {} } }", parserTest_INVALID_PRIVATE_ID)
}

// Java: ParserTest.testPrivateProperty_invalid_className()
func TestPrivateProperty_invalid_className(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("class #PC {}", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("class C extends #PSC {}", parserTest_INVALID_PRIVATE_ID)
}

// Java: ParserTest.testPrivateProperty_invalid_referencePrivateOutsideClass()
func TestPrivateProperty_invalid_referencePrivateOutsideClass(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("class C { #f = 1; } const c = new C(); c.#f;", parserTest_INVALID_PRIVATE_ID)
}

// Java: ParserTest.testPrivateProperty_invalid_referencePrivateOutsideClassViaOptionalChain()
func TestPrivateProperty_invalid_referencePrivateOutsideClassViaOptionalChain(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("class C { #f = 1; } const c = new C(); c?.#f;", parserTest_INVALID_PRIVATE_ID)
}

// Java: ParserTest.testPrivateProperty_invalid_referencePrivatePropFromObjectLiteral()
func TestPrivateProperty_invalid_referencePrivatePropFromObjectLiteral(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("const o = {}; o.#f = 1;", parserTest_INVALID_PRIVATE_ID)
}

// Java: ParserTest.testPrivateProperty_invalid_referencePrivatePropFromObjectLiteralViaOptionalChain()
func TestPrivateProperty_invalid_referencePrivatePropFromObjectLiteralViaOptionalChain(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("const o = {}; o?.#f;", parserTest_INVALID_PRIVATE_ID)
}

// Java: ParserTest.testPrivateProperty_invalid_lineNumber()
func TestPrivateProperty_invalid_lineNumber(t *testing.T) {
	pt := newParserTest(t)

	var errors []string
	var lines []int
	var cols []int
	reporter := &parserTest_recordingErrorReporter{
		error_: func(message string, sourceName string, line int, lineOffset int) {
			errors = append(errors, message)
			lines = append(lines, line)
			cols = append(cols, lineOffset)
		},
	}
	ParserRunner_Parse(
		rhino.NewSimpleSourceFile("testcode", rhino.StaticSourceFile_SourceKind_STRONG),
		"const x = #foo;\nconst y = {\n  #pf: 1,\n};\nconst z = #bar;",
		pt.createConfig(),
		reporter)
	assertThat(t, errors).
		ContainsExactly(parserTest_INVALID_PRIVATE_ID, parserTest_INVALID_PRIVATE_ID, parserTest_INVALID_PRIVATE_ID).
		InOrder()
	assertThat(t, lines).ContainsExactly(1, 3, 5).InOrder()
	assertThat(t, cols).ContainsExactly(10, 2, 10).InOrder()
}

// parserTest_recordingErrorReporter is the anonymous ErrorReporter of
// testPrivateProperty_invalid_lineNumber: warnings are ignored, errors are
// passed to error_.
//
// Go-only: language-adapter — anonymous class.
type parserTest_recordingErrorReporter struct {
	error_ func(message string, sourceName string, line int, lineOffset int)
}

// Java: ErrorReporter.warning(String, String, int, int)
func (r *parserTest_recordingErrorReporter) Warning(message string, sourceName string, line int, lineOffset int) {
}

// Java: ErrorReporter.error(String, String, int, int)
func (r *parserTest_recordingErrorReporter) Error(message string, sourceName string, line int, lineOffset int) {
	r.error_(message, sourceName, line, lineOffset)
}

// Java: ParserTest.testPrivateProperty_invalid_importName()
func TestPrivateProperty_invalid_importName(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("import #pi from './someModule'", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("import {#pi} from './someModule'", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("import {x as #pi} from './someModule'", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("import * as #pi from './someModule'", parserTest_INVALID_PRIVATE_ID)
}

// Java: ParserTest.testPrivateProperty_invalid_exportName()
func TestPrivateProperty_invalid_exportName(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("export const #px = 1", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("export var #px = 1", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("export function #pf() {}", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("export class #pc {}", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("export {#px}", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("export {x as #px}", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("export {#px as default}", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("export {#y as class}", parserTest_INVALID_PRIVATE_ID)

	pt.parseError("export {x as #px} from './someModule'", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("export {default as #pd} from './someModule'", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("export {#px as default} from './someModule'", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("export {#pc as class} from './someModule'", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("export {#px} from './someModule'", parserTest_INVALID_PRIVATE_ID)
}

// Java: ParserTest.testPrivateProperty_invalid_labeledStatement()
func TestPrivateProperty_invalid_labeledStatement(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("#pl: while (true) {}", parserTest_INVALID_PRIVATE_ID)
	pt.parseError("while (true) { break #pl; }", parserTest_INVALID_PRIVATE_ID, "undefined label \"#pl\"")
	pt.parseError("while (true) { continue #pl; }", parserTest_INVALID_PRIVATE_ID, "undefined label \"#pl\"")

	pt.parseError("class C { method() { #pl: while (true) {} } }", parserTest_INVALID_PRIVATE_ID)
	pt.parseError(
		"class C { method() { while (true) { break #pl; } } }",
		parserTest_INVALID_PRIVATE_ID,
		"undefined label \"#pl\"")
	pt.parseError(
		"class C { method() { while (true) { continue #pl; } } }",
		parserTest_INVALID_PRIVATE_ID,
		"undefined label \"#pl\"")
}

// Java: ParserTest.testPrivateProperty_inOperatorWithPrivateProp_valid()
func TestPrivateProperty_inOperatorWithPrivateProp_valid(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_PRIVATE_ELEMENTS)

	pt.parse("class C { #f = 1; static isC(x) { return #f in x; } }")
	pt.parse("class C { #f = this; m() { return #f in this.#f; } }")
}

// Java: ParserTest.testPrivateProperty_inOperatorWithPrivateProp_linenocharno()
func TestPrivateProperty_inOperatorWithPrivateProp_linenocharno(t *testing.T) {
	pt := newParserTest(t)

	n :=
		pt.parse(
			`class C {
  #f = 1;
  static isC(x) {
    return #f in x;
  }
}
`).
			GetFirstChild()

	members := nodeUtil_getClassMembers(n)

	staticMethod := members.GetLastChild()
	assertNode(t, staticMethod).HasType(rhino.Token_MEMBER_FUNCTION_DEF)
	methodBlock := staticMethod.GetFirstChild().GetLastChild()
	assertNode(t, methodBlock).HasType(rhino.Token_BLOCK)
	returnStatement := methodBlock.GetFirstChild()
	assertNode(t, returnStatement).HasType(rhino.Token_RETURN)
	inExpression := returnStatement.GetFirstChild()
	assertNode(t, inExpression).HasType(rhino.Token_IN)
	inExpressionLeft := inExpression.GetFirstChild()
	assertNode(t, inExpressionLeft).HasType(rhino.Token_NAME)
	assertNode(t, inExpressionLeft).HasStringThat().IsEqualTo("#f")
	assertThat(t, inExpressionLeft.GetLineno()).IsEqualTo(4)
	assertThat(t, inExpressionLeft.GetCharno()).IsEqualTo(11)
	assertThat(t, inExpressionLeft.GetLength()).IsEqualTo(2)
}

// Java: ParserTest.testPrivateProperty_inOperatorWithPrivateProp_invalid_nonExistentPrivateProp()
func TestPrivateProperty_inOperatorWithPrivateProp_invalid_nonExistentPrivateProp(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_PRIVATE_ELEMENTS)

	pt.parseError("class C { static isC(x) { return #missing in x; } }", parserTest_PRIVATE_FIELD_NOT_DEFINED)
}

// Java: ParserTest.testPrivateProperty_inOperatorWithPrivateProp_invalid_privatePropOnRhs()
func TestPrivateProperty_inOperatorWithPrivateProp_invalid_privatePropOnRhs(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_PRIVATE_ELEMENTS)

	pt.parseError("class C { #f = 1; static isC(x) { return #f in #f; } }", parserTest_INVALID_PRIVATE_ID)
}

// Java: ParserTest.testPrivateProperty_inOperatorWithPrivateProp_invalid_surroundingParenthesis()
func TestPrivateProperty_inOperatorWithPrivateProp_invalid_surroundingParenthesis(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_PRIVATE_ELEMENTS)

	pt.parseError("class C { #f = 1; static isC(x) { return (#f) in x; } }", parserTest_INVALID_PRIVATE_ID)
}

// Java: ParserTest.testPrivateProperty_inOperatorWithPrivateProp_invalid_notInClass()
func TestPrivateProperty_inOperatorWithPrivateProp_invalid_notInClass(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("const o = {}; if (#f in o) {}", parserTest_INVALID_PRIVATE_ID)
}

// Java: ParserTest.testPrivateProperty_stringKeyThatLooksLikePrivateProp()
func TestPrivateProperty_stringKeyThatLooksLikePrivateProp(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("const o = {'#notAPrivateProp': 1};")
}

// Java: ParserTest.testEmptyClassStaticBlock()
func TestEmptyClassStaticBlock(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("class C { static { } }")
	pt.parse("let a = class { static { } };")
}

// Java: ParserTest.testReturnInClassStaticBlock()
func TestReturnInClassStaticBlock(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("function f() {class C { static { return; } }}", "return must be inside function")
	pt.parseError("class C { static { return; } }", "return must be inside function")
	pt.parse("class C {static {function f() {return;}}}")
}

// Java: ParserTest.testContinueInClassStaticBlock()
func TestContinueInClassStaticBlock(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("class C {static {continue;}}", parserTest_UNEXPECTED_CONTINUE)
	pt.parseError("for (let i = 1; i < 2; i++) {class C {static {continue;}}}", parserTest_UNEXPECTED_CONTINUE)
	pt.parseError("while (true) {class C {static {continue;}}}", parserTest_UNEXPECTED_CONTINUE)
	pt.parse("class C {static {while (true) {continue;}}}")
	pt.parse("class C {static {for (let i = 1; i < 2; i++) {continue;}}}")
	pt.parseError(
		"x: for (let i = 1; i < 2; i++) {class C {static {continue x;}}}",
		parserTest_UNDEFINED_LABEL+" \"x\"")
	pt.parseError("x: while (true) {class C {static {continue x;}}}", parserTest_UNDEFINED_LABEL+" \"x\"")
	pt.parse("class C {static {x: while (true) {continue x;}}}")
	pt.parse("class C {static {x: for (let i = 1; i < 2; i++) {continue x;}}}")
	pt.parseError("class C {static {x: { while (true) {continue x;}}}}", parserTest_UNEXPECTED_LABELLED_CONTINUE)
	pt.parseError(
		"class C {static {x: {for (let i = 1; i < 2; i++) {continue x;}}}}",
		parserTest_UNEXPECTED_LABELLED_CONTINUE)
}

// Java: ParserTest.testBreakInClassStaticBlock()
func TestBreakInClassStaticBlock(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("class C {static {break;}}", parserTest_UNLABELED_BREAK)
	pt.parseError("for (let i = 1; i < 2; i++) {class C {static {break;}}}", parserTest_UNLABELED_BREAK)
	pt.parseError("while (true) {class C {static {break;}}}", parserTest_UNLABELED_BREAK)
	pt.parse("class C {static {while (true) {break;}}}")
	pt.parse("class C {static {for (let i = 1; i < 2; i++) {break;}}}")
	pt.parseError(
		"x: for (let i = 1; i < 2; i++) {class C {static {break x;}}}", parserTest_UNDEFINED_LABEL+" \"x\"")
	pt.parseError("x: while (true) {class C {static {break x;}}}", parserTest_UNDEFINED_LABEL+" \"x\"")
	pt.parse("class C {static {x: while (true) {break x;}}}")
	pt.parse("class C {static {x: for (let i = 1; i < 2; i++) {break x;}}}")
	pt.parse("class C {static {x: { while (true) {break x;}}}}")
	pt.parse("class C {static {x: {for (let i = 1; i < 2; i++) {break x;}}}}")
}

// Java: ParserTest.testYieldInClassStaticBlock()
func TestYieldInClassStaticBlock(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("class C {static {var ind; yield;}}", "primary expression expected")
	pt.parseError("function* f(ind) {class C{ static {yield ind; ind++;}}}", parserTest_UNEXPECTED_YIELD)
	pt.parse("class C{ static {function* f(ind) {yield ind; ind++;}}}")
}

// Java: ParserTest.testClassStaticBlock_this()
func TestClassStaticBlock_this(t *testing.T) {
	pt := newParserTest(t)

	// multiple fields
	pt.parse(
		`class C {
static field1 = 1; static field2 = 2; static field3 = 3;
static {
let x = this.field1; let y = this.field2; let z = this.field3;
}
}
`)
	pt.parse("class C { static { this.field1 = 1; this.field2 = 2; this.field3 = 3; } }")
	pt.parse(
		`let a = class {
static field1 = 1; static field2 = 2; static field3 = 3;
static {
let x = this.field1; let y = this.field2; let z = this.field3;
}
};
`)
	pt.parse("let a = class { static { this.field1 = 1; this.field2 = 2; this.field3 = 3; } };")
	// functions
	pt.parse(
		`class C {
static field1 = 1;
static {
function incr() { return ++A.field1; }
console.log(incr());
if(incr()) {
this.field2 = 2;
}
}
}
`)
	// try catch
	pt.parse(
		`class C {
static field1 = 1;
static {
try {
this.field1 = 2;
}
catch {
}
}
}
`)
}

// Java: ParserTest.testClassStaticBlock_inheritance()
func TestClassStaticBlock_inheritance(t *testing.T) {
	pt := newParserTest(t)

	// It is a syntax error to call super() in a class static initialization block.
	// Must get reported in CheckSuper.java.
	pt.parse("class Base {} class C extends Base { static { super(); } }")
	// allow accessing static properties of the base class
	pt.parse("class Base { static y; } class C extends Base { static { super.y; } }")
	// allow accessing non-static properties of the base class
	pt.parse("class Base { y; } class C extends Base { static { super.y; } }")
}

// Java: ParserTest.testClassExtendsLeftHandSideExpression()
func TestClassExtendsLeftHandSideExpression(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("class A {} class B extends (0, A) {}")
	pt.parseError("class A {} class B extends 0, A {}", "'{' expected")
}

// Java: ParserTest.testMultipleClassStaticBlocks()
func TestMultipleClassStaticBlocks(t *testing.T) {
	pt := newParserTest(t)

	// empty
	pt.parse("class C { static { } static { } }")
	pt.parse("let a = class { static { } static { } };")
	// multiple fields
	pt.parse(
		`class C {
static field1 = 1; static field2 = 2; static field3 = 3;
static {
let x = this.field1; let y = this.field2;
}
static {
let z = this.field3;
}
}
`)
	pt.parse("class C { static { this.field1 = 1; this.field2 = 2; } static { this.field3 = 3; } }")
	pt.parse(
		`let a = class {
static field1 = 1; static field2 = 2; static field3 = 3;
static {
let x = this.field1; let y = this.field2;
}
static {
let z = this.field3;
}
};
`)
	pt.parse(
		`let a = class {
static {
this.field1 = 1; this.field2 = 2;
}
static {
 this.field3 = 3;
}
};
`)
}

// Java: ParserTest.testClassStaticBlock_linenocharno()
func TestClassStaticBlock_linenocharno(t *testing.T) {
	pt := newParserTest(t)

	n := pt.parse(`class C {
 static {}
 }`).GetFirstChild()

	assertNode(t, n).HasType(rhino.Token_CLASS)
	assertNode(t, n).HasLineno(1)
	assertNode(t, n).HasCharno(0)

	members := nodeUtil_getClassMembers(n)
	assertNode(t, members).HasType(rhino.Token_CLASS_MEMBERS)

	staticBlock := members.GetFirstChild()

	assertNode(t, staticBlock).HasType(rhino.Token_BLOCK)
	assertNode(t, staticBlock).HasLineno(2)
	assertNode(t, staticBlock).HasCharno(8)
	assertNode(t, staticBlock).HasLength(2)
}

// Java: ParserTest.testClassStaticBlock_invalid()
func TestClassStaticBlock_invalid(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("class { {} }", "'identifier' expected")
	pt.parseError("class { static { static { } } }", "'identifier' expected")
	pt.parseError("var o = { static {} };", "Cannot use keyword in short object literal")
}

// Java: ParserTest.testClassStaticSuper()
func TestClassStaticSuper(t *testing.T) {
	pt := newParserTest(t)

	pt.parse(
		`class Bar {
  static double(n) {
    return n * 2
  }
}
class Baz extends Bar {
  // Used from a static field initializer.
  static val1 = super.double(6);

  static val2;
  static {
    // Used from within a static block.
    Baz.val2 = super.double(5);
  }
}
`)
}

// Java: ParserTest.testSuper1()
func TestSuper1(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_SUPER)
	pt.strictMode = Config_StrictMode_SLOPPY

	// TODO(johnlenz): super in global scope should be a syntax error
	pt.parse("super;")

	pt.parse("function f() {super;};")

	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.parseWarning("super;", requiresLanguageModeMessage(parser.FeatureSet_Feature_SUPER))

	pt.mode = Config_LanguageMode_ECMASCRIPT3
	pt.parseWarning("super;", requiresLanguageModeMessage(parser.FeatureSet_Feature_SUPER))
}

// Java: ParserTest.testNewTarget()
func TestNewTarget(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_NEW_TARGET)
	pt.strictMode = Config_StrictMode_SLOPPY

	pt.parseError("new.target;", "new.target must be inside a function")

	pt.parse("function f() { new.target; };")

	pt.mode = Config_LanguageMode_ECMASCRIPT3
	pt.parseWarning(
		"class C { f() { new.target; } }",
		requiresLanguageModeMessage(parser.FeatureSet_Feature_CLASSES),
		requiresLanguageModeMessage(parser.FeatureSet_Feature_MEMBER_DECLARATIONS),
		requiresLanguageModeMessage(parser.FeatureSet_Feature_NEW_TARGET))

	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.parseWarning(
		"class C { f() { new.target; } }",
		requiresLanguageModeMessage(parser.FeatureSet_Feature_CLASSES),
		requiresLanguageModeMessage(parser.FeatureSet_Feature_MEMBER_DECLARATIONS),
		requiresLanguageModeMessage(parser.FeatureSet_Feature_NEW_TARGET))

	pt.mode = Config_LanguageMode_ECMASCRIPT_2015
	pt.expectFeatures(parser.FeatureSet_Feature_CLASSES, parser.FeatureSet_Feature_MEMBER_DECLARATIONS, parser.FeatureSet_Feature_NEW_TARGET)
	pt.parse("class C { f() { new.target; } }")
}

// Java: ParserTest.testNewDotSomethingInvalid()
func TestNewDotSomethingInvalid(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY

	pt.parseError("function f(){new.something}", "'target' expected")
}

// Java: ParserTest.hookWithDecimalNotParsedAsOptionalChaining()
func TestHookWithDecimalNotParsedAsOptionalChaining(t *testing.T) {
	pt := newParserTest(t)

	n := pt.parse("a?.1:2").GetFirstFirstChild()

	assertNode(t, n).HasType(rhino.Token_HOOK)
	assertNode(t, n).HasLineno(1)
	assertNode(t, n).HasCharno(0)
	assertNode(t, n.GetFirstChild()).IsEqualTo(rhino.IR_Name("a"))
	assertNode(t, n.GetSecondChild()).IsEqualTo(rhino.IR_Number(0.1))
	assertNode(t, n.GetLastChild()).IsEqualTo(rhino.IR_Number(2.0))
}

// Java: ParserTest.optChainGetProp()
func TestOptChainGetProp(t *testing.T) {
	pt := newParserTest(t)

	n := pt.parse("a?.b").GetFirstFirstChild()

	assertNode(t, n).HasType(rhino.Token_OPTCHAIN_GETPROP)
	assertNode(t, n).IsOptionalChainStart()
	assertNode(t, n).HasLineno(1)
	assertNode(t, n).HasCharno(3)
	assertNode(t, n).HasLength(1)
	assertNode(t, n.GetFirstChild()).IsEqualTo(rhino.IR_Name("a"))
	assertNode(t, n).HasStringThat().IsEqualTo("b")
}

// Java: ParserTest.optChainGetPropWithKeyword()
func TestOptChainGetPropWithKeyword(t *testing.T) {
	pt := newParserTest(t)

	n := pt.parse("a?.finally").GetFirstFirstChild()

	assertNode(t, n).HasType(rhino.Token_OPTCHAIN_GETPROP)
	assertNode(t, n).IsOptionalChainStart()
	assertNode(t, n).HasLineno(1)
	assertNode(t, n).HasCharno(3)
	assertNode(t, n).HasLength(7)
	assertNode(t, n.GetFirstChild()).IsEqualTo(rhino.IR_Name("a"))
	assertNode(t, n).HasStringThat().IsEqualTo("finally")
}

// Java: ParserTest.optChainGetElem()
func TestOptChainGetElem(t *testing.T) {
	pt := newParserTest(t)

	n := pt.parse("a?.[1]").GetFirstFirstChild()

	assertNode(t, n).HasType(rhino.Token_OPTCHAIN_GETELEM)
	assertNode(t, n).IsOptionalChainStart()
	assertNode(t, n).HasLineno(1)
	assertNode(t, n).HasCharno(0)
	assertNode(t, n).HasLength(6)
	assertNode(t, n.GetFirstChild()).IsEqualTo(rhino.IR_Name("a"))
	assertNode(t, n.GetSecondChild()).IsEqualTo(rhino.IR_Number(1.0))
}

// Java: ParserTest.optChainCall()
func TestOptChainCall(t *testing.T) {
	pt := newParserTest(t)

	n := pt.parse("a?.()").GetFirstFirstChild()

	assertNode(t, n).HasType(rhino.Token_OPTCHAIN_CALL)
	assertNode(t, n).IsOptionalChainStart()
	assertNode(t, n).HasLineno(1)
	assertNode(t, n).HasCharno(0)
	assertNode(t, n).HasLength(5)
	assertNode(t, n.GetFirstChild()).IsEqualTo(rhino.IR_Name("a"))
}

// Java: ParserTest.optChainStartOfChain_innerChainIsArgOfACall()
func TestOptChainStartOfChain_innerChainIsArgOfACall(t *testing.T) {
	pt := newParserTest(t)

	optChainCall := pt.parse("a?.b?.(x?.y);").GetFirstFirstChild()
	assertThat(t, optChainCall.IsOptChainCall()).IsTrue()

	optChainGetProp := optChainCall.GetFirstChild() // `a?.b`
	optChainArg := optChainCall.GetLastChild()      // `x?.y`

	// Check that `a?.b` is the start of an opt chain
	assertThat(t, optChainGetProp.IsOptChainGetProp()).IsTrue()
	assertThat(t, optChainGetProp.IsOptionalChainStart()).IsTrue()

	// Check that `x?.y` is the start of an opt chain
	assertThat(t, optChainArg.IsOptChainGetProp()).IsTrue()
	assertThat(t, optChainArg.IsOptionalChainStart()).IsTrue()
}

// Java: ParserTest.optChainStartOfChain_optChainGetProp()
func TestOptChainStartOfChain_optChainGetProp(t *testing.T) {
	pt := newParserTest(t)

	outerGet := pt.parse("a?.b.c").GetFirstFirstChild()
	innerGet := outerGet.GetFirstChild()

	// `a?.b.c`
	assertNode(t, outerGet).HasType(rhino.Token_OPTCHAIN_GETPROP)
	assertNode(t, outerGet).IsNotOptionalChainStart()
	assertNode(t, outerGet).HasLineno(1)
	assertNode(t, outerGet).HasCharno(5)
	assertNode(t, outerGet).HasLength(1)

	// `a?.b`
	assertNode(t, innerGet).HasType(rhino.Token_OPTCHAIN_GETPROP)
	assertNode(t, innerGet).IsOptionalChainStart()
	assertNode(t, innerGet).HasLineno(1)
	assertNode(t, innerGet).HasCharno(3)
	assertNode(t, innerGet).HasLength(1)

	assertNode(t, outerGet).HasStringThat().IsEqualTo("c")
}

// Java: ParserTest.optChainStartOfChain_optChainGetElem()
func TestOptChainStartOfChain_optChainGetElem(t *testing.T) {
	pt := newParserTest(t)

	outerGet := pt.parse("a?.[b][c]").GetFirstFirstChild()
	innerGet := outerGet.GetFirstChild()

	// `a?.[b][c]`
	assertNode(t, outerGet).HasType(rhino.Token_OPTCHAIN_GETELEM)
	assertNode(t, outerGet).IsNotOptionalChainStart()
	assertNode(t, outerGet).HasLineno(1)
	assertNode(t, outerGet).HasCharno(0)
	assertNode(t, outerGet).HasLength(9)
	//
	// // `a?.[b]`
	assertNode(t, innerGet).HasType(rhino.Token_OPTCHAIN_GETELEM)
	assertNode(t, innerGet).IsOptionalChainStart()
	assertNode(t, innerGet).HasLineno(1)
	assertNode(t, innerGet).HasCharno(0)
	assertNode(t, innerGet).HasLength(6)

	assertNode(t, outerGet.GetSecondChild()).IsEqualTo(rhino.IR_Name("c"))
}

// Java: ParserTest.optChainStartOfChain_optChainCall()
func TestOptChainStartOfChain_optChainCall(t *testing.T) {
	pt := newParserTest(t)

	outerCall := pt.parse("a?.()(b)").GetFirstFirstChild()
	innerCall := outerCall.GetFirstChild()

	// `a?.()(b)`
	assertNode(t, outerCall).HasType(rhino.Token_OPTCHAIN_CALL)
	assertNode(t, outerCall).IsNotOptionalChainStart()
	assertNode(t, outerCall).HasLineno(1)
	assertNode(t, outerCall).HasCharno(0)
	assertNode(t, outerCall).HasLength(8)

	// `a?.()`
	assertNode(t, innerCall).HasType(rhino.Token_OPTCHAIN_CALL)
	assertNode(t, innerCall).IsOptionalChainStart()
	assertNode(t, innerCall).HasLineno(1)
	assertNode(t, innerCall).HasCharno(0)
	assertNode(t, innerCall).HasLength(5)

	assertNode(t, outerCall.GetSecondChild()).IsEqualTo(rhino.IR_Name("b"))
}

// Java: ParserTest.optChainParens_optChainGetProp()
func TestOptChainParens_optChainGetProp(t *testing.T) {
	pt := newParserTest(t)

	outerGet := pt.parse("(a?.b).c").GetFirstFirstChild()
	innerGet := outerGet.GetFirstChild()

	// `(a?.b).c`
	assertNode(t, outerGet).HasType(rhino.Token_GETPROP)
	assertNode(t, outerGet).HasLineno(1)
	assertNode(t, outerGet).HasCharno(7)
	assertNode(t, outerGet).HasLength(1)

	// `a?.b`
	assertNode(t, innerGet).HasType(rhino.Token_OPTCHAIN_GETPROP)
	assertNode(t, innerGet).IsOptionalChainStart()
	assertNode(t, innerGet).HasLineno(1)
	assertNode(t, innerGet).HasCharno(4)
	assertNode(t, innerGet).HasLength(1)

	assertNode(t, outerGet).HasStringThat().IsEqualTo("c")
}

// Java: ParserTest.callExpressionBeforeOptionalGetProp()
func TestCallExpressionBeforeOptionalGetProp(t *testing.T) {
	pt := newParserTest(t)

	get := pt.parse("a()?.b").GetFirstFirstChild()
	call := get.GetFirstChild()

	assertNode(t, get).HasType(rhino.Token_OPTCHAIN_GETPROP)
	assertNode(t, get).HasLineno(1)
	assertNode(t, get).HasCharno(5)
	assertNode(t, get).HasLength(1)

	assertNode(t, call).HasType(rhino.Token_CALL)
	assertNode(t, call).HasLineno(1)
	assertNode(t, call).HasCharno(0)
	assertNode(t, call).HasLength(3)
}

// Java: ParserTest.optChainChain()
func TestOptChainChain(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("a?.b?.c")
	pt.parse("a.b?.c")
	pt.parse("a?.b?.[1]")
	pt.parse("a?.b?.()")
	pt.parse("a?.[1]?.b")
	pt.parse("a?.[1]?.b()")
	pt.parse("a?.b?.c?.d")
	pt.parse("a?.b?.c?.[1]")
	pt.parse("a?.b?.c?.()")
	pt.parse("a?.(c)?.b")
	pt.parse("a().b?.c")
}

// Java: ParserTest.optChainAssignError()
func TestOptChainAssignError(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("a?.b = c", "invalid assignment target")
}

// Java: ParserTest.optChainConstructorError()
func TestOptChainConstructorError(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("new a?.()", "Optional chaining is forbidden in construction contexts.")
	pt.parseError("new a?.b()", "Optional chaining is forbidden in construction contexts.")
}

// Java: ParserTest.optChainTemplateLiteralError()
func TestOptChainTemplateLiteralError(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("a?.()?.`hello`", "template literal cannot be used within optional chaining")
	pt.parseError("a?.`hello`", "template literal cannot be used within optional chaining")
	pt.parseError("a?.b`hello`", "template literal cannot be used within optional chaining")
	// https://github.com/tc39/test262/blob/master/test/language/expressions/optional-chaining/early-errors-tail-position-template-string-esi.js
	// test to prevent automatic semicolon insertion rules
	pt.parseError("a?.b\n`hello`", "template literal cannot be used within optional chaining")
}

// Java: ParserTest.optChainMiscErrors()
func TestOptChainMiscErrors(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("super?.()", "Optional chaining is forbidden in super?.")
	pt.parseError("super?.foo", "Optional chaining is forbidden in super?.")
	pt.parseError("new?.target", "Optional chaining is forbidden in `new?.target` contexts.")
	pt.parseError("import?.('foo')", "Optional chaining is forbidden in import?.")
}

// Java: ParserTest.optChainDeleteValid()
func TestOptChainDeleteValid(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("delete a?.b")
}

// Java: ParserTest.optChainEs2019()
func TestOptChainEs2019(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_OPTIONAL_CHAINING)
	pt.mode = Config_LanguageMode_ECMASCRIPT_2019

	pt.parseWarning("a?.b", requiresLanguageModeMessage(parser.FeatureSet_Feature_OPTIONAL_CHAINING))
}

// Java: ParserTest.optChainSyntaxError()
func TestOptChainSyntaxError(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("a?.{}", "syntax error: { not allowed in optional chain")
	// optional chain cannot be applied on a BLOCK `{}`
	pt.parseError("{a:x}?.a", "primary expression expected")
	pt.parse("({a:x})?.a")
}

// Java: ParserTest.testArrow1()
func TestArrow1(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_ARROW_FUNCTIONS)
	pt.strictMode = Config_StrictMode_SLOPPY

	pt.parse("()=>1;")
	pt.parse("()=>{}")
	pt.parse("(a,b) => a + b;")
	pt.parse("a => b")
	pt.parse("a => { return b }")
	pt.parse("a => b")
	pt.parse("var x = (a => b);")

	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.parseWarning("a => b", requiresLanguageModeMessage(parser.FeatureSet_Feature_ARROW_FUNCTIONS))

	pt.mode = Config_LanguageMode_ECMASCRIPT3
	pt.parseWarning("a => b;", requiresLanguageModeMessage(parser.FeatureSet_Feature_ARROW_FUNCTIONS))
}

// Java: ParserTest.testArrowInvalid1()
func TestArrowInvalid1(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parseError("*()=>1;", "primary expression expected")
	pt.parseError("var f = x\n=>2", "No newline allowed before '=>'")
	pt.parseError("f = (x,y)\n=>2;", "No newline allowed before '=>'")
	pt.parseError("f( (x,y)\n=>2)", "No newline allowed before '=>'")
}

// Java: ParserTest.testAsyncFunction()
func TestAsyncFunction(t *testing.T) {
	pt := newParserTest(t)

	asyncFunctionExpressionSource := "f = async function() {};"
	asyncFunctionDeclarationSource := "async function f() {}"
	pt.expectFeatures(parser.FeatureSet_Feature_ASYNC_FUNCTIONS)

	for _, m := range Config_LanguageMode_values() {
		pt.mode = m
		if m == Config_LanguageMode_ECMASCRIPT3 {
			pt.strictMode = Config_StrictMode_SLOPPY
		} else {
			pt.strictMode = Config_StrictMode_STRICT
		}
		if m.FeatureSet().Has(parser.FeatureSet_Feature_ASYNC_FUNCTIONS) {
			pt.parse(asyncFunctionExpressionSource)
			pt.parse(asyncFunctionDeclarationSource)
		} else {
			pt.parseWarning(
				asyncFunctionExpressionSource, requiresLanguageModeMessage(parser.FeatureSet_Feature_ASYNC_FUNCTIONS))
			pt.parseWarning(
				asyncFunctionDeclarationSource, requiresLanguageModeMessage(parser.FeatureSet_Feature_ASYNC_FUNCTIONS))
		}
	}
}

// Java: ParserTest.testAsyncNamedFunction()
func TestAsyncNamedFunction(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT_2015
	pt.expectFeatures(
		parser.FeatureSet_Feature_CLASSES,
		parser.FeatureSet_Feature_MEMBER_DECLARATIONS,
		parser.FeatureSet_Feature_CONST_DECLARATIONS,
		parser.FeatureSet_Feature_LET_DECLARATIONS)
	pt.parse(
		`class C {
  async(x) { return x; }
}
const c = new C();
c.async(1);
let foo = async(5);
`)
}

// Java: ParserTest.testAsyncGeneratorFunction()
func TestAsyncGeneratorFunction(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_ASYNC_FUNCTIONS, parser.FeatureSet_Feature_GENERATORS, parser.FeatureSet_Feature_ASYNC_GENERATORS)
	pt.parse("async function *f(){}")
	pt.parse("f = async function *(){}")
	pt.parse("class C { async *foo(){} }")
}

// Java: ParserTest.testAsyncArrowFunction()
func TestAsyncArrowFunction(t *testing.T) {
	pt := newParserTest(t)

	pt.doAsyncArrowFunctionTest("f = async (x) => x + 1")
	pt.doAsyncArrowFunctionTest("f = async x => x + 1")
}

// Java: ParserTest.testAsyncArrowInvalid()
func TestAsyncArrowInvalid(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("f = not_async (x) => x + 1;", "'=>' unexpected")
}

// Java: ParserTest.testAsyncMethod()
func TestAsyncMethod(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_ASYNC_FUNCTIONS)
	pt.parse("o={async m(){}}")
	pt.parse("o={async [a+b](){}}")
	pt.parse("class C{async m(){}}")
	pt.parse("class C{static async m(){}}")
	pt.parse("class C{async [a+b](){}}")
	pt.parse("class C{static async [a+b](){}}")
}

// Java: ParserTest.testInvalidAsyncMethod()
func TestInvalidAsyncMethod(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY
	pt.expectFeatures(parser.FeatureSet_Feature_MEMBER_DECLARATIONS)
	// 'async' allowed as a name
	pt.parse("o={async(){}}")
	pt.parse("class C{async(){}}")
	pt.parse("class C{static async(){}}")

	pt.expectFeatures()
	pt.parse("o={async:false}")
	// newline after 'async' forces it to be the property name
	pt.parseError("o={async\nm(){}}", "'}' expected")
	pt.parseError("o={static async\nm(){}}", "Cannot use keyword in short object literal")

	pt.expectFeatures(parser.FeatureSet_Feature_PUBLIC_CLASS_FIELDS)
	pt.parse("class C{async};")              // class field
	pt.parse("class C{async\nm(){}}")        // class field
	pt.parse("class C{static async\nm(){}}") // class field
}

// Java: ParserTest.testAwait_inFunctions()
func TestAwait_inFunctions(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_ASYNC_FUNCTIONS)
	pt.parse("async function f(p) { await p }")
	pt.parse("async function f(p) { async function f2() { await p; } }")
	pt.parse("f = async function(p) { await p }")
	pt.parse("f = async(p) => await p")
	pt.parse("class C{ async m(p) { await p }}")
	pt.parse("class C{ static async m(p) { await p }}")
}

// Java: ParserTest.testAwait_inMethods()
func TestAwait_inMethods(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("class C { async method() { await 2; } }")
	pt.parse("const x = { async method() { await 2; } }")
}

// Java: ParserTest.testAwait_inObjectLiteralComputedPropertyName()
func TestAwait_inObjectLiteralComputedPropertyName(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("async function f() { return { [await 1]: 1 }; }")
	pt.parse("async function f() { return { [await 2]() {} }; }")
	pt.parse("async function f() { return { *[await 3]() {} }; }")
	pt.parse("async function f() { return { async [await 4]() {} }; }")
	pt.parse("async function f() { return { async *[await 5]() {} }; }")
	pt.parse("async function f() { return { get [await 6]() {} }; }")
	pt.parse("async function f() { return { set [await 7](x) {} }; }")

	pt.parse("const f = async () => ({ [await 8]: 1 });")
}

// Java: ParserTest.testAwait_inClassComputedPropertyName()
func TestAwait_inClassComputedPropertyName(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("async function f() { return class { [await 1] = 1 }; }")
	pt.parse("async function f() { return class { [await 2]() {} }; }")
	pt.parse("async function f() { return class { *[await 3]() {} }; }")
	pt.parse("async function f() { return class { async [await 4]() {} }; }")
	pt.parse("async function f() { return class { async *[await 5]() {} }; }")
	pt.parse("async function f() { return class { get [await 6]() {} }; }")
	pt.parse("async function f() { return class { set [await 7](x) {} }; }")

	pt.parse("const f = async () => class { [await 8] = 1 };")
}

// Java: ParserTest.testAwait_invalid_missingExpression()
func TestAwait_invalid_missingExpression(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("async function f() { await; }", "primary expression expected")
}

// Java: ParserTest.testAwait_invalid_nonAsyncFunction()
func TestAwait_invalid_nonAsyncFunction(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("function f() { return await 5; }", parserTest_UNEXPECTED_AWAIT)
	pt.parseError("async function f() { function f2() { return await 5; } }", parserTest_UNEXPECTED_AWAIT)
}

// Java: ParserTest.testAwait_invalid_nonAsyncMethod()
func TestAwait_invalid_nonAsyncMethod(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("class C { method() { await 1; } }", parserTest_UNEXPECTED_AWAIT)
	pt.parseError("const x = { method() { await 2; } }", parserTest_UNEXPECTED_AWAIT)
}

// Java: ParserTest.testAwait_invalid_defaultParameter()
func TestAwait_invalid_defaultParameter(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError(
		"async function f(x = await 15) { return x; }",
		"`await` is illegal in parameter default value.")
}

// Java: ParserTest.testAwait_invalid_inClassStaticBlock()
func TestAwait_invalid_inClassStaticBlock(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("class C { static { await 1; } }", parserTest_UNEXPECTED_AWAIT)
	pt.parseError("async function f() { class C { static { await 1; } } }", parserTest_UNEXPECTED_AWAIT)
	pt.parseError("async () => { class C { static { await 1; } } }", parserTest_UNEXPECTED_AWAIT)

	// Double check that await in this position nested within an async function is valid.
	pt.parse("class C { static { async function f() { await 1; } } }")
	pt.parse("class C { static { async () => { await 1; } } }")
}

// Java: ParserTest.testAwait_invalid_inClassFieldInitializer()
func TestAwait_invalid_inClassFieldInitializer(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("class C { f = await 1; }", parserTest_UNEXPECTED_AWAIT)
	pt.parseError("class C { static f = await 1; }", parserTest_UNEXPECTED_AWAIT)
	pt.parseError("async function f() { class C { static f = await 1; } }", parserTest_UNEXPECTED_AWAIT)
	pt.parseError("async () => { class C { static f = await 1; } }", parserTest_UNEXPECTED_AWAIT)

	// Double check that await in this position nested within an async function is valid.
	pt.parse("class C { static f = async function f() { await 1; } }")
	pt.parse("class C { static f = async () => { await 1; } }")
}

// Java: ParserTest.testAwait_inClassMethods()
func TestAwait_inClassMethods(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("class C { method() { await 1; } }", parserTest_UNEXPECTED_AWAIT)

	pt.parse("class C { async method() { await 2; } }")
}

// Java: ParserTest.testAwait_topLevelInEsModule()
func TestAwait_topLevelInEsModule(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_TOP_LEVEL_AWAIT)
	pt.parse("export const x = await 15;")
	pt.parse("export const x = 1; await 15;")
}

// Java: ParserTest.testAwait_topLevelInObjectLiteral()
func TestAwait_topLevelInObjectLiteral(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_TOP_LEVEL_AWAIT)
	pt.parse("export const x = 1; const obj = { key: await x };")
}

// Java: ParserTest.testAwait_topLevelInBlock()
func TestAwait_topLevelInBlock(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_TOP_LEVEL_AWAIT)
	pt.parse("export const x = 1; { await x; }")
}

// Java: ParserTest.testAwait_topLevelInLoops()
func TestAwait_topLevelInLoops(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_TOP_LEVEL_AWAIT)
	pt.parse("export const x = 1; for (;;) { await x; }")
	pt.parse("export const x = 1; while (true) { await x; }")
	pt.parse("export const x = 1; do { await x; } while (true);")
}

// Java: ParserTest.testAwait_topLevelInTryCatch()
func TestAwait_topLevelInTryCatch(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_TOP_LEVEL_AWAIT)
	pt.parse("export const x = 1; try { await x; } catch (e) {}")
	pt.parse("export const x = 1; try {} catch (e) { await x; }")
	pt.parse("export const x = 1; try {} finally { await x; }")
}

// Java: ParserTest.testAwait_topLevelInClassExtends()
func TestAwait_topLevelInClassExtends(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_TOP_LEVEL_AWAIT)
	pt.parse("export const x = class {}; class C extends (await x) {}")
}

// Java: ParserTest.testAwait_invalid_topLevelInScript()
func TestAwait_invalid_topLevelInScript(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("const x = await 15;", parserTest_UNEXPECTED_AWAIT)
	pt.parseError("await 15;", parserTest_UNEXPECTED_AWAIT)
}

// Java: ParserTest.testAwait_invalid_topLevelInClassExtends()
func TestAwait_invalid_topLevelInClassExtends(t *testing.T) {
	pt := newParserTest(t)

	// Await is not allowed in extends clause without parentheses
	// Use export to ensure module context where await is a keyword
	pt.parseError("export const x = 1; class C extends await x {}", "'{' expected")
}

// Java: ParserTest.testAwait_invalid_topLevelInGoogModule()
func TestAwait_invalid_topLevelInGoogModule(t *testing.T) {
	pt := newParserTest(t)

	// Await is not allowed in goog.module
	pt.parseError("goog.module('m'); await 1;", parserTest_UNEXPECTED_AWAIT)
}

// Java: ParserTest.testAwait_invalid_topLevelInGoogProvide()
func TestAwait_invalid_topLevelInGoogProvide(t *testing.T) {
	pt := newParserTest(t)

	// Await is not allowed in goog.provide
	pt.parseError("goog.provide('m'); await 1;", parserTest_UNEXPECTED_AWAIT)
}

// Java: ParserTest.testAwait_topLevelInObjectLiteralComputedPropertyName()
func TestAwait_topLevelInObjectLiteralComputedPropertyName(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_TOP_LEVEL_AWAIT)
	pt.parse("export const x = { [await 1]: 1 };")
	pt.parse("export const x = { [await 2]() {} };")
	pt.parse("export const x = { *[await 3]() {} };")
	pt.parse("export const x = { async [await 4]() {} };")
	pt.parse("export const x = { async *[await 5]() {} };")
	pt.parse("export const x = { get [await 6]() {} };")
	pt.parse("export const x = { set [await 7](x) {} };")

	pt.parse("export const x = 1; const y = { [await 8]: 8 };")
}

// Java: ParserTest.testAwait_topLevelAwaitInClassComputedPropertyName()
func TestAwait_topLevelAwaitInClassComputedPropertyName(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_TOP_LEVEL_AWAIT)
	pt.parse("export class C { [await 1] = 1 };")
	pt.parse("export class C { [await 2]() {} };")
	pt.parse("export class C { *[await 3]() {} };")
	pt.parse("export class C { async [await 4]() {} };")
	pt.parse("export class C { async *[await 5]() {} };")
	pt.parse("export class C { get [await 6]() {} };")
	pt.parse("export class C { set [await 7](x) {} };")

	pt.parse("export class C { [await 8] = 1 };")

	pt.parse("export const x = 1; class C { [await 9] = 1 };")
}

// Java: ParserTest.testAwait_invalid_topLevelAwaitInClassComputedPropertyNameInScript()
func TestAwait_invalid_topLevelAwaitInClassComputedPropertyNameInScript(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("class C { [await 1] = 1 };", parserTest_UNEXPECTED_AWAIT)
}

// Java: ParserTest.testForAwaitOf_inFunctions()
func TestForAwaitOf_inFunctions(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY

	pt.expectFeatures(parser.FeatureSet_Feature_FOR_AWAIT_OF)
	pt.parse("async () => { for await (a of b) c;}")
	pt.parse("async () => { for await (var a of b) c;}")
	pt.parse("async () => { for await (a.x of b) c;}")
	pt.parse("async () => { for await ([a1, a2, a3] of b) c;}")
	pt.parse("async () => { for await (const {x, y, z} of b) c;}")
	// default value inside a pattern isn't an initializer
	pt.parse("async () => { for await (const {x, y = 2, z} of b) c;}")
	pt.expectFeatures(parser.FeatureSet_Feature_FOR_AWAIT_OF, parser.FeatureSet_Feature_LET_DECLARATIONS)
	pt.parse("async () => { for await (let a of b) c;}")
	pt.expectFeatures(parser.FeatureSet_Feature_FOR_AWAIT_OF, parser.FeatureSet_Feature_CONST_DECLARATIONS)
	pt.parse("async () => { for await (const a of b) c;}")
}

// Java: ParserTest.testForAwaitOf_invalid_declWithInitializer()
func TestForAwaitOf_invalid_declWithInitializer(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY

	pt.parseError("async () => { for await (a = 1 of b) c;}", parserTest_INVALID_ASSIGNMENT_TARGET)
	pt.parseError("async () => { for await (var a = 1 of b) c;}", parserTest_INVALID_FOR_AWAIT_INIT)
	pt.parseError("async () => { for await (let a = 1 of b) c;}", parserTest_INVALID_FOR_AWAIT_INIT)
	pt.parseError("async () => { for await (const a = 1 of b) c;}", parserTest_INVALID_FOR_AWAIT_INIT)
	pt.parseError("async () => { for await (let {a} = {} of b) c;}", parserTest_INVALID_FOR_AWAIT_INIT)
}

// Java: ParserTest.testForAwaitOf_invalid_multipleDecls()
func TestForAwaitOf_invalid_multipleDecls(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY

	pt.parseError("async () => { for await (a, b of c) d;}", parserTest_INVALID_ASSIGNMENT_TARGET)

	pt.parseError("async () => { for await (var a, b of c) d;}", parserTest_INVALID_FOR_AWAIT_MULT_DECLS)
	pt.parseError("async () => { for await (let a, b of c) d;}", parserTest_INVALID_FOR_AWAIT_MULT_DECLS)
	pt.parseError("async () => { for await (const a, b of c) d;}", parserTest_INVALID_FOR_AWAIT_MULT_DECLS)
}

// Java: ParserTest.testForAwaitOf_topLevelInEsModule()
func TestForAwaitOf_topLevelInEsModule(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_TOP_LEVEL_AWAIT)
	pt.parse("export const x = 1; for await (let a of b) foo();")
}

// Java: ParserTest.testForAwaitOf_invalid_topLevelInScript()
func TestForAwaitOf_invalid_topLevelInScript(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("for await (let a of b) foo();", parserTest_UNEXPECTED_FOR_AWAIT_OF)
}

// Java: ParserTest.testFor_ES5()
func TestFor_ES5(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("for (var x; x != 10; x = next()) {}")
	pt.parse("for (var x; x != 10; x = next());")
	pt.parse("for (var x = 0; x != 10; x++) {}")
	pt.parse("for (var x = 0; x != 10; x++);")

	pt.parse("var x; for (x; x != 10; x = next()) {}")
	pt.parse("var x; for (x; x != 10; x = next());")

	pt.parseError("for (x in {};;) {}", "')' expected")
}

// Java: ParserTest.testFor_ES6()
func TestFor_ES6(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY

	pt.expectFeatures(parser.FeatureSet_Feature_LET_DECLARATIONS)
	pt.parse("for (let x; x != 10; x = next()) {}")
	pt.parse("for (let x; x != 10; x = next());")
	pt.parse("for (let x = 0; x != 10; x++) {}")
	pt.parse("for (let x = 0; x != 10; x++);")

	pt.expectFeatures(parser.FeatureSet_Feature_CONST_DECLARATIONS)
	pt.parse("for (const x = 0; x != 10; x++) {}")
	pt.parse("for (const x = 0; x != 10; x++);")
}

// Java: ParserTest.testForConstNoInitializer()
func TestForConstNoInitializer(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("for (const x; x != 10; x = next()) {}", "const variables must have an initializer")
	pt.parseError("for (const x; x != 10; x = next());", "const variables must have an initializer")
}

// Java: ParserTest.testForIn_ES6()
func TestForIn_ES6(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY

	pt.parse("for (a in b) c;")
	pt.parse("for (var a in b) c;")

	pt.expectFeatures(parser.FeatureSet_Feature_LET_DECLARATIONS)
	pt.parse("for (let a in b) c;")

	pt.expectFeatures(parser.FeatureSet_Feature_CONST_DECLARATIONS)
	pt.parse("for (const a in b) c;")

	pt.expectFeatures()
	pt.parseError("for (a,b in c) d;", parserTest_INVALID_ASSIGNMENT_TARGET)
	pt.parseError(
		"for (var a,b in c) d;",
		"for-in statement may not have more than one variable declaration")

	pt.parseError(
		"for (let a,b in c) d;",
		"for-in statement may not have more than one variable declaration")

	pt.parseError(
		"for (const a,b in c) d;",
		"for-in statement may not have more than one variable declaration")

	pt.parseError("for (a=1 in b) c;", parserTest_INVALID_ASSIGNMENT_TARGET)
	pt.parseError("for (let a=1 in b) c;", "for-in statement may not have initializer")
	pt.parseError("for (const a=1 in b) c;", "for-in statement may not have initializer")
	pt.parseError("for (var a=1 in b) c;", "for-in statement may not have initializer")
	pt.parseError("for (\"a\" in b) c;", parserTest_INVALID_ASSIGNMENT_TARGET)
}

// Java: ParserTest.testForIn_ES5()
func TestForIn_ES5(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.strictMode = Config_StrictMode_SLOPPY

	pt.parse("for (a in b) c;")
	pt.parse("for (var a in b) c;")

	pt.parseError("for (a=1 in b) c;", parserTest_INVALID_ASSIGNMENT_TARGET)
	pt.parseWarning("for (var a=1 in b) c;", "for-in statement should not have initializer")
}

// Java: ParserTest.testForInDestructuring()
func TestForInDestructuring(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY

	pt.expectFeatures(parser.FeatureSet_Feature_OBJECT_DESTRUCTURING)
	pt.parse("for ({a} in b) c;")
	pt.parse("for (var {a} in b) c;")
	pt.expectFeatures(parser.FeatureSet_Feature_OBJECT_DESTRUCTURING, parser.FeatureSet_Feature_LET_DECLARATIONS)
	pt.parse("for (let {a} in b) c;")
	pt.expectFeatures(parser.FeatureSet_Feature_OBJECT_DESTRUCTURING, parser.FeatureSet_Feature_CONST_DECLARATIONS)
	pt.parse("for (const {a} in b) c;")

	pt.expectFeatures(parser.FeatureSet_Feature_OBJECT_DESTRUCTURING)
	pt.parse("for ({a: b} in c) d;")
	pt.parse("for (var {a: b} in c) d;")
	pt.expectFeatures(parser.FeatureSet_Feature_OBJECT_DESTRUCTURING, parser.FeatureSet_Feature_LET_DECLARATIONS)
	pt.parse("for (let {a: b} in c) d;")
	pt.expectFeatures(parser.FeatureSet_Feature_OBJECT_DESTRUCTURING, parser.FeatureSet_Feature_CONST_DECLARATIONS)
	pt.parse("for (const {a: b} in c) d;")

	pt.expectFeatures(parser.FeatureSet_Feature_ARRAY_DESTRUCTURING)
	pt.parse("for ([a] in b) c;")
	pt.parse("for (var [a] in b) c;")
	pt.expectFeatures(parser.FeatureSet_Feature_ARRAY_DESTRUCTURING, parser.FeatureSet_Feature_LET_DECLARATIONS)
	pt.parse("for (let [a] in b) c;")
	pt.expectFeatures(parser.FeatureSet_Feature_ARRAY_DESTRUCTURING, parser.FeatureSet_Feature_CONST_DECLARATIONS)
	pt.parse("for (const [a] in b) c;")
}

// Java: ParserTest.testForInDestructuringInvalid()
func TestForInDestructuringInvalid(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY

	pt.parseError("for ({a: b} = foo() in c) d;", parserTest_INVALID_ASSIGNMENT_TARGET)
	pt.parseError("for (var {a: b} = foo() in c) d;", "for-in statement may not have initializer")
	pt.parseError("for (let {a: b} = foo() in c) d;", "for-in statement may not have initializer")
	pt.parseError("for (const {a: b} = foo() in c) d;", "for-in statement may not have initializer")

	pt.parseError("for ([a] = foo() in b) c;", parserTest_INVALID_ASSIGNMENT_TARGET)
	pt.parseError("for (var [a] = foo() in b) c;", "for-in statement may not have initializer")
	pt.parseError("for (let [a] = foo() in b) c;", "for-in statement may not have initializer")
	pt.parseError("for (const [a] = foo() in b) c;", "for-in statement may not have initializer")
}

// Java: ParserTest.testForOf1()
func TestForOf1(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY

	pt.expectFeatures(parser.FeatureSet_Feature_FOR_OF)
	pt.parse("for(a of b) c;")
	pt.parse("for(var a of b) c;")
	pt.expectFeatures(parser.FeatureSet_Feature_FOR_OF, parser.FeatureSet_Feature_LET_DECLARATIONS)
	pt.parse("for(let a of b) c;")
	pt.expectFeatures(parser.FeatureSet_Feature_FOR_OF, parser.FeatureSet_Feature_CONST_DECLARATIONS)
	pt.parse("for(const a of b) c;")
}

// Java: ParserTest.testForOf2()
func TestForOf2(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY

	pt.parseError("for(a=1 of b) c;", parserTest_INVALID_ASSIGNMENT_TARGET)
	pt.parseError("for(var a=1 of b) c;", "for-of statement may not have initializer")
	pt.parseError("for(let a=1 of b) c;", "for-of statement may not have initializer")
	pt.parseError("for(const a=1 of b) c;", "for-of statement may not have initializer")
}

// Java: ParserTest.testForOf3()
func TestForOf3(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY

	pt.parseError(
		"for(var a, b of c) d;",
		"for-of statement may not have more than one variable declaration")
	pt.parseError(
		"for(let a, b of c) d;",
		"for-of statement may not have more than one variable declaration")
	pt.parseError(
		"for(const a, b of c) d;",
		"for-of statement may not have more than one variable declaration")
}

// Java: ParserTest.testForOf4()
func TestForOf4(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY

	pt.parseError("for(a, b of c) d;", parserTest_INVALID_ASSIGNMENT_TARGET)
}

// Java: ParserTest.testDestructuringInForLoops()
func TestDestructuringInForLoops(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY

	// Destructuring forbids an initializer in for-in/for-of
	pt.parseError("for (var {x: y} = foo() in bar()) {}", "for-in statement may not have initializer")
	pt.parseError("for (let {x: y} = foo() in bar()) {}", "for-in statement may not have initializer")
	pt.parseError(
		"for (const {x: y} = foo() in bar()) {}", "for-in statement may not have initializer")

	pt.parseError("for (var {x: y} = foo() of bar()) {}", "for-of statement may not have initializer")
	pt.parseError("for (let {x: y} = foo() of bar()) {}", "for-of statement may not have initializer")
	pt.parseError(
		"for (const {x: y} = foo() of bar()) {}", "for-of statement may not have initializer")

	// but requires it in a vanilla for loop
	pt.parseError("for (var {x: y};;) {}", "destructuring must have an initializer")
	pt.parseError("for (let {x: y};;) {}", "destructuring must have an initializer")
	pt.parseError("for (const {x: y};;) {}", "const variables must have an initializer")
}

// Java: ParserTest.testInvalidDestructuring()
func TestInvalidDestructuring(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY

	// {x: 5} and {x: 'str'} are valid object literals but not valid patterns.
	pt.parseError("for ({x: 5} in foo()) {}", parserTest_INVALID_ASSIGNMENT_TARGET)
	pt.parseError("for ({x: 'str'} in foo()) {}", parserTest_INVALID_ASSIGNMENT_TARGET)
	pt.parseError("var {x: 5} = foo();", parserTest_INVALID_ASSIGNMENT_TARGET)
	pt.parseError("var {x: 'str'} = foo();", parserTest_INVALID_ASSIGNMENT_TARGET)
	pt.parseError("({x: 5} = foo());", parserTest_INVALID_ASSIGNMENT_TARGET)
	pt.parseError("({x: 'str'} = foo());", parserTest_INVALID_ASSIGNMENT_TARGET)

	// {method(){}} is a valid object literal but not a valid object pattern.
	pt.parseError("function f({method(){}}) {}", "'}' expected")
	pt.parseError("function f({method(){}} = foo()) {}", "'}' expected")
}

// Java: ParserTest.testForOfPatterns()
func TestForOfPatterns(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY

	pt.expectFeatures(parser.FeatureSet_Feature_FOR_OF, parser.FeatureSet_Feature_OBJECT_DESTRUCTURING)
	pt.parse("for({x} of b) c;")
	pt.parse("for({x: y} of b) c;")

	pt.expectFeatures(parser.FeatureSet_Feature_FOR_OF, parser.FeatureSet_Feature_ARRAY_DESTRUCTURING)
	pt.parse("for([x, y] of b) c;")
	pt.parse("for([x, ...y] of b) c;")

	pt.expectFeatures(parser.FeatureSet_Feature_FOR_OF, parser.FeatureSet_Feature_OBJECT_DESTRUCTURING, parser.FeatureSet_Feature_LET_DECLARATIONS)
	pt.parse("for(let {x} of b) c;")
	pt.parse("for(let {x: y} of b) c;")

	pt.expectFeatures(parser.FeatureSet_Feature_FOR_OF, parser.FeatureSet_Feature_ARRAY_DESTRUCTURING, parser.FeatureSet_Feature_LET_DECLARATIONS)
	pt.parse("for(let [x, y] of b) c;")
	pt.parse("for(let [x, ...y] of b) c;")

	pt.expectFeatures(parser.FeatureSet_Feature_FOR_OF, parser.FeatureSet_Feature_OBJECT_DESTRUCTURING, parser.FeatureSet_Feature_CONST_DECLARATIONS)
	pt.parse("for(const {x} of b) c;")
	pt.parse("for(const {x: y} of b) c;")

	pt.expectFeatures(parser.FeatureSet_Feature_FOR_OF, parser.FeatureSet_Feature_ARRAY_DESTRUCTURING, parser.FeatureSet_Feature_CONST_DECLARATIONS)
	pt.parse("for(const [x, y] of b) c;")
	pt.parse("for(const [x, ...y] of b) c;")

	pt.expectFeatures(parser.FeatureSet_Feature_FOR_OF, parser.FeatureSet_Feature_OBJECT_DESTRUCTURING)
	pt.parse("for(var {x} of b) c;")
	pt.parse("for(var {x: y} of b) c;")

	pt.expectFeatures(parser.FeatureSet_Feature_FOR_OF, parser.FeatureSet_Feature_ARRAY_DESTRUCTURING)
	pt.parse("for(var [x, y] of b) c;")
	pt.parse("for(var [x, ...y] of b) c;")
}

// Java: ParserTest.testForOfPatternsWithInitializer()
func TestForOfPatternsWithInitializer(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY

	pt.parseError("for({x}=a of b) c;", parserTest_INVALID_ASSIGNMENT_TARGET)
	pt.parseError("for({x: y}=a of b) c;", parserTest_INVALID_ASSIGNMENT_TARGET)
	pt.parseError("for([x, y]=a of b) c;", parserTest_INVALID_ASSIGNMENT_TARGET)
	pt.parseError("for([x, ...y]=a of b) c;", parserTest_INVALID_ASSIGNMENT_TARGET)

	pt.parseError("for(let {x}=a of b) c;", "for-of statement may not have initializer")
	pt.parseError("for(let {x: y}=a of b) c;", "for-of statement may not have initializer")
	pt.parseError("for(let [x, y]=a of b) c;", "for-of statement may not have initializer")
	pt.parseError("for(let [x, ...y]=a of b) c;", "for-of statement may not have initializer")

	pt.parseError("for(const {x}=a of b) c;", "for-of statement may not have initializer")
	pt.parseError("for(const {x: y}=a of b) c;", "for-of statement may not have initializer")
	pt.parseError("for(const [x, y]=a of b) c;", "for-of statement may not have initializer")
	pt.parseError("for(const [x, ...y]=a of b) c;", "for-of statement may not have initializer")
}

// Java: ParserTest.testImport()
func TestImport(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_MODULES)
	pt.strictMode = Config_StrictMode_SLOPPY

	pt.parse("import 'someModule'")
	pt.parse("import d from './someModule'")
	pt.parse("import {} from './someModule'")
	pt.parse("import {x, y} from './someModule'")
	pt.parse("import {x as x1, y as y1} from './someModule'")
	pt.parse("import {x as x1, y as y1, } from './someModule'")
	pt.parse("import {default as d, class as c} from './someModule'")
	pt.parse("import d, {x as x1, y as y1} from './someModule'")
	pt.parse("import * as sm from './someModule'")

	pt.expectFeatures()
	pt.parseError("import class from './someModule'", "cannot use keyword 'class' here.")
	pt.parseError("import * as class from './someModule'", "'identifier' expected")
	pt.parseError("import {a as class} from './someModule'", "'identifier' expected")
	pt.parseError("import {class} from './someModule'", "'as' expected")
}

// Java: ParserTest.testExport()
func TestExport(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY

	pt.expectFeatures(parser.FeatureSet_Feature_MODULES)
	pt.parse("export const x = 1")
	pt.parse("export var x = 1")
	pt.parse("export function f() {}")
	pt.parse("export class c {}")
	pt.parse("export {x, y}")
	pt.parse("export {x as x1}")
	pt.parse("export {x as x1, y as x2}")
	pt.parse("export {x as default, y as class}")

	pt.expectFeatures()
	pt.parseError("export {default as x}", "cannot use keyword 'default' here.")
	pt.parseError("export {package as x}", "cannot use keyword 'package' here.")
	pt.parseError("export {package}", "cannot use keyword 'package' here.")

	pt.expectFeatures(parser.FeatureSet_Feature_MODULES)
	pt.parse("export {x as x1, y as y1} from './someModule'")
	pt.parse("export {x as x1, y as y1, } from './someModule'")
	pt.parse("export {default as d} from './someModule'")
	pt.parse("export {d as default, c as class} from './someModule'")
	pt.parse("export {default as default, class as class} from './someModule'")
	pt.parse("export {class} from './someModule'")
	pt.parse("export * from './someModule'")

	pt.expectFeatures()
	pt.parseError("export * as s from './someModule';", "'from' expected")
}

// Java: ParserTest.testExportAsync()
func TestExportAsync(t *testing.T) {
	pt := newParserTest(t)

	pt.strictMode = Config_StrictMode_SLOPPY

	pt.expectFeatures(parser.FeatureSet_Feature_MODULES, parser.FeatureSet_Feature_ASYNC_FUNCTIONS)
	pt.parse("export async function f() {}")
}

// Java: ParserTest.testImportExportTypescriptKeyword()
func TestImportExportTypescriptKeyword(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("export { namespace };")
	pt.parse("import { namespace } from './input0.js';")
}

// Java: ParserTest.testGoogModule()
func TestGoogModule(t *testing.T) {
	pt := newParserTest(t)

	tree := pt.parse("goog.module('example');")
	assertNode(t, tree).HasType(rhino.Token_SCRIPT)
	assertThat(t, tree.GetStaticSourceFile()).IsNotNull()
	assertNode(t, tree.GetFirstChild()).HasType(rhino.Token_MODULE_BODY)
	assertThat(t, tree.GetFirstChild().GetStaticSourceFile()).IsNotNull()
}

// Java: ParserTest.testShebang()
func TestShebang(t *testing.T) {
	pt := newParserTest(t)

	pt.parse("#!/usr/bin/node\n var x = 1;")
	pt.parseError("var x = 1; \n #!/usr/bin/node", "Shebang comment must be at the start of the file")
}

// Java: ParserTest.testInvalidPoundUsage()
func TestInvalidPoundUsage(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("var x = 1; \n# Wrong-style comment", "Invalid usage of #")
}

// Java: ParserTest.testLookaheadGithubIssue699()
func TestLookaheadGithubIssue699(t *testing.T) {
	pt := newParserTest(t)

	start := time.Now()
	pt.parse(
		`[1,[1,[1,[1,[1,[1,
[1,[1,[1,[1,[1,[1,
[1,[1,[1,[1,[1,[1,
[1,[1,[1,[1,[1,[1,
[1,[1,[1,[1,[1,[1,
[1,[1,
[1]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]
`)

	stop := time.Now()

	assertWithMessage(t, "runtime").That(stop.Sub(start).Milliseconds()).IsLessThan(5000)
}

// Java: ParserTest.testInvalidHandling1()
func TestInvalidHandling1(t *testing.T) {
	pt := newParserTest(t)

	pt.parse(
		`/**
 * @fileoverview Definition.
 * @mods {ns.bar}
 * @modName mod
 *
 * @extends {ns.bar}
 * @author someone
 */

goog.provide('ns.foo');
`)
}

// Java: ParserTest.testUtf8()
func TestUtf8(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.strictMode = Config_StrictMode_SLOPPY
	n := pt.parse("\ufefffunction f() {}\n")
	fn := n.GetFirstChild()
	assertNode(t, fn).HasType(rhino.Token_FUNCTION)
}

// Java: ParserTest.testParseDeep1()
func TestParseDeep1(t *testing.T) {
	pt := newParserTest(t)

	code := "var x; x = \n"
	for i := 1; i < 15000; i++ {
		code += "  '" + jutil.ToJavaString(i) + "' +\n"
	}
	code += "'end';n"
	pt.parse(code)
}

// Java: ParserTest.testParseDeep2()
func TestParseDeep2(t *testing.T) {
	pt := newParserTest(t)

	code := "var x; x = \n"
	for i := 1; i < 15000; i++ {
		code += "  '" + jutil.ToJavaString(i) + "' +\n"
	}
	code += "'end'; /** a comment */\n"
	pt.parse(code)
}

// Java: ParserTest.testParseDeep3()
func TestParseDeep3(t *testing.T) {
	pt := newParserTest(t)

	code := "var x; x = \n"
	for i := 1; i < 15000; i++ {
		code += "  '" + jutil.ToJavaString(i) + "' +\n"
	}
	code += "  /** @type {string} */ (x);\n"
	pt.parse(code)
}

// Java: ParserTest.testParseDeep4()
func TestParseDeep4(t *testing.T) {
	pt := newParserTest(t)

	// Currently, we back off if there is any JSDoc in the tree of binary expressions
	code := "var x; x = \n"
	for i := 1; i < 15000; i++ {
		if i == 5 {
			code += "  /** @type {string} */ (x) +\n"
		}
		code += "  '" + jutil.ToJavaString(i) + "' +\n"
	}
	code += "'end';n"
	var e jutil.IRuntimeException
	jutil.Try(func() {
		pt.parse(code)
	}, func(caught jutil.IRuntimeException) { e = caught })
	if e == nil {
		// Go-only: the deep-recursion difference. Java's IRFactory recurses
		// once per binary operator here (the JSDoc disables the iterative
		// processBinaryExpressionHelper) and overflows the thread stack; the
		// StackOverflowError is rethrown by ParserRunner.parse as a
		// RuntimeException. A goroutine stack grows instead of overflowing
		// (and an overflow could not be caught), so the Go parse succeeds.
		t.Skip("deep-recursion difference: Java overflows its stack here, Go does not")
	}
	// expected exception
	assertThat(t, e).HasMessageThat().Contains("Exception parsing")
}

// Java: ParserTest.testParseInlineSourceMap()
func TestParseInlineSourceMap(t *testing.T) {
	pt := newParserTest(t)

	code :=
		`var X = (function () {
    function X(input) {
        this.y = input;
    }
    return X;
}());
console.log(new X(1));
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiZm9vLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiZm9vLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiJBQUFBO0lBR0UsV0FBWSxLQUFhO1FBQ3ZCLElBQUksQ0FBQyxDQUFDLEdBQUcsS0FBSyxDQUFDO0lBQ2pCLENBQUM7SUFDSCxRQUFDO0FBQUQsQ0FBQyxBQU5ELElBTUM7QUFFRCxPQUFPLENBQUMsR0FBRyxDQUFDLElBQUksQ0FBQyxDQUFDLENBQUMsQ0FBQyxDQUFDLENBQUMifQ==
`
	result := pt.doParse(code)
	assertThat(t, result.SourceMapURL).
		IsEqualTo(
			"data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiZm9vLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiZm9vLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiJBQUFBO0lBR0UsV0FBWSxLQUFhO1FBQ3ZCLElBQUksQ0FBQyxDQUFDLEdBQUcsS0FBSyxDQUFDO0lBQ2pCLENBQUM7SUFDSCxRQUFDO0FBQUQsQ0FBQyxBQU5ELElBTUM7QUFFRCxPQUFPLENBQUMsR0FBRyxDQUFDLElBQUksQ0FBQyxDQUFDLENBQUMsQ0FBQyxDQUFDLENBQUMifQ==")
}

// Java: ParserTest.testParseSourceMapRelativeURL()
func TestParseSourceMapRelativeURL(t *testing.T) {
	pt := newParserTest(t)

	code :=
		`var X = (function () {
    function X(input) {
        this.y = input;
    }
    return X;
}());
console.log(new X(1));
//# sourceMappingURL=somefile.js.map
`
	result := pt.doParse(code)
	assertThat(t, result.SourceMapURL).IsEqualTo("somefile.js.map")
}

// Java: ParserTest.testParseSourceMapAbsoluteURL()
func TestParseSourceMapAbsoluteURL(t *testing.T) {
	pt := newParserTest(t)

	code :=
		`console.log('asdf');
//# sourceMappingURL=/some/absolute/path/to/somefile.js.map
`
	result := pt.doParse(code)
	assertThat(t, result.SourceMapURL).IsEqualTo("/some/absolute/path/to/somefile.js.map")
}

// Java: ParserTest.testParseSourceMapAbsoluteURLHTTP()
func TestParseSourceMapAbsoluteURLHTTP(t *testing.T) {
	pt := newParserTest(t)

	code :=
		`console.log('asdf');
//# sourceMappingURL=http://google.com/some/absolute/path/to/somefile.js.map
`
	result := pt.doParse(code)
	assertThat(t, result.SourceMapURL).
		IsEqualTo("http://google.com/some/absolute/path/to/somefile.js.map")
}

// Java: ParserTest.testIncorrectAssignmentDoesntCrash()
func TestIncorrectAssignmentDoesntCrash(t *testing.T) {
	pt := newParserTest(t)

	// Check that error make sense in default "stop on error" mode.
	pt.parseError("[1 + 2] = 3;", "invalid assignment target")

	// Ensure that in IDE mode parser doesn't crash. It produces much more errors but it's
	// "ignore errors" mode so it's ok.
	pt.isIdeMode = true
	pt.parseError(
		"[1 + 2] = 3;",
		"invalid assignment target",
		"']' expected",
		"invalid assignment target",
		"Semi-colon expected",
		"Semi-colon expected",
		"primary expression expected",
		"invalid assignment target",
		"Semi-colon expected",
		"primary expression expected",
		"Semi-colon expected")
}

// Java: ParserTest.testDynamicImport()
func TestDynamicImport(t *testing.T) {
	pt := newParserTest(t)

	dynamicImportUses := []string{
		"import('foo')",
		"import('foo').then(function(a) { return a; })",
		"var moduleNamespace = import('foo')",
		"Promise.all([import('foo')]).then(function(a) { return a; })",
		"function foo() { foo(); import('foo'); } foo();",
	}
	pt.expectFeatures(parser.FeatureSet_Feature_DYNAMIC_IMPORT)

	for _, m := range Config_LanguageMode_values() {
		pt.mode = m
		if m == Config_LanguageMode_ECMASCRIPT3 {
			pt.strictMode = Config_StrictMode_SLOPPY
		} else {
			pt.strictMode = Config_StrictMode_STRICT
		}
		if m.FeatureSet().Has(parser.FeatureSet_Feature_DYNAMIC_IMPORT) {
			for _, importUseSource := range dynamicImportUses {
				pt.parse(importUseSource)
			}
		} else {
			for _, importUseSource := range dynamicImportUses {
				pt.parseWarning(importUseSource, requiresLanguageModeMessage(parser.FeatureSet_Feature_DYNAMIC_IMPORT))
			}
		}
	}

	pt.mode = Config_LanguageMode_ECMASCRIPT_2020
	pt.parseError("function foo() { import bar from './someModule'; }", "'(' expected")
}

// Java: ParserTest.testAwaitDynamicImport()
func TestAwaitDynamicImport(t *testing.T) {
	pt := newParserTest(t)

	awaitDynamicImportUses := []string{
		"(async function() { return await import('foo'); })()",
		"(async function() { await import('foo').then(function(a) { return a; }); })()",
		"(async function() { var moduleNamespace = await import('foo'); })()",
		`(async function() {
await Promise.all([import('foo')]).then(function(a) { return a; }); })()
`,
	}
	pt.expectFeatures(parser.FeatureSet_Feature_DYNAMIC_IMPORT, parser.FeatureSet_Feature_ASYNC_FUNCTIONS)

	for _, importUseSource := range awaitDynamicImportUses {
		pt.parse(importUseSource)
	}
}

// Java: ParserTest.testImportMeta()
func TestImportMeta(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_MODULES, parser.FeatureSet_Feature_IMPORT_META)

	tree := pt.parse("import.meta")
	assertNode(t, tree.GetFirstFirstChild()).IsEqualTo(rhino.IR_ExprResult(rhino.IR_ImportMeta()))
}

// Java: ParserTest.testImportMeta_es5()
func TestImportMeta_es5(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.expectFeatures(parser.FeatureSet_Feature_MODULES, parser.FeatureSet_Feature_IMPORT_META)

	pt.parseWarning(
		"import.meta",
		requiresLanguageModeMessage(parser.FeatureSet_Feature_MODULES),
		requiresLanguageModeMessage(parser.FeatureSet_Feature_IMPORT_META))
}

// Java: ParserTest.testImportMeta_es6()
func TestImportMeta_es6(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT_2015
	pt.expectFeatures(parser.FeatureSet_Feature_MODULES, parser.FeatureSet_Feature_IMPORT_META)

	pt.parseWarning("import.meta", requiresLanguageModeMessage(parser.FeatureSet_Feature_IMPORT_META))
}

// Java: ParserTest.testImportMeta_inExpression()
func TestImportMeta_inExpression(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_MODULES, parser.FeatureSet_Feature_IMPORT_META)

	propTree := pt.parse("import.meta.url")
	assertNode(t, propTree.GetFirstFirstChild()).
		IsEqualTo(rhino.IR_ExprResult(rhino.IR_Getprop(rhino.IR_ImportMeta(), "url")))

	callTree := pt.parse("f(import.meta.url)")
	assertNode(t, callTree.GetFirstFirstChild()).
		IsEqualTo(
			rhino.IR_ExprResult(pt.freeCall(rhino.IR_Call(rhino.IR_Name("f"), rhino.IR_Getprop(rhino.IR_ImportMeta(), "url")))))
}

// Java: ParserTest.testImportMeta_asDotProperty()
func TestImportMeta_asDotProperty(t *testing.T) {
	pt := newParserTest(t)

	tree := pt.parse("x.import.meta")
	assertNode(t, tree.GetFirstChild()).
		IsEqualTo(rhino.IR_ExprResult(rhino.IR_GetpropWithMoreProps(rhino.IR_Name("x"), "import", "meta")))
}

// Java: ParserTest.testNullishCoalesce()
func TestNullishCoalesce(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_NULL_COALESCE_OP)

	tree := pt.parse("x??y")
	assertNode(t, tree.GetFirstChild()).
		IsEqualTo(rhino.IR_ExprResult(rhino.IR_Coalesce(rhino.IR_Name("x"), rhino.IR_Name("y"))))
}

// Java: ParserTest.testNullishCoalesce_es2019()
func TestNullishCoalesce_es2019(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT_2019
	pt.expectFeatures(parser.FeatureSet_Feature_NULL_COALESCE_OP)

	pt.parseWarning("x??y", requiresLanguageModeMessage(parser.FeatureSet_Feature_NULL_COALESCE_OP))
}

// Java: ParserTest.testNullishCoalesce_withLogicalAND_shouldFail()
func TestNullishCoalesce_withLogicalAND_shouldFail(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("x&&y??z", "Logical OR and logical AND require parentheses when used with '??'")
}

// Java: ParserTest.testNullishCoalesce_withLogicalOR_shouldFail()
func TestNullishCoalesce_withLogicalOR_shouldFail(t *testing.T) {
	pt := newParserTest(t)

	pt.parseError("x??y||z", "Logical OR and logical AND require parentheses when used with '??'")
}

// Java: ParserTest.testNullishCoalesce_withLogicalANDinParens()
func TestNullishCoalesce_withLogicalANDinParens(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_NULL_COALESCE_OP)

	tree := pt.parse("(x&&y)??z")
	assertNode(t, tree.GetFirstChild()).
		IsEqualTo(rhino.IR_ExprResult(rhino.IR_Coalesce(rhino.IR_And(rhino.IR_Name("x"), rhino.IR_Name("y")), rhino.IR_Name("z"))))
}

// Java: ParserTest.testNullishCoalesce_chaining()
func TestNullishCoalesce_chaining(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_NULL_COALESCE_OP)

	tree := pt.parse("x??y??z")
	expr := tree.GetFirstChild()
	coalesce := expr.GetFirstChild()

	assertNode(t, expr).
		IsEqualTo(
			rhino.IR_ExprResult(rhino.IR_Coalesce(rhino.IR_Coalesce(rhino.IR_Name("x"), rhino.IR_Name("y")), rhino.IR_Name("z"))))
	assertNode(t, expr).HasLineno(1).HasCharno(0).HasLength(7)
	assertNode(t, coalesce).HasType(rhino.Token_COALESCE)
	assertNode(t, coalesce).HasLineno(1).HasCharno(0).HasLength(7)
	assertNode(t, coalesce.GetFirstChild()).HasType(rhino.Token_COALESCE)
	assertNode(t, coalesce.GetFirstChild()).HasLineno(1).HasCharno(0).HasLength(4)
}

// Java: ParserTest.testAssignOR()
func TestAssignOR(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_LOGICAL_ASSIGNMENT)
	tree := pt.parse("x||=y")
	assertNode(t, tree.GetFirstChild()).
		IsEqualTo(rhino.IR_ExprResult(rhino.IR_AssignOr(rhino.IR_Name("x"), rhino.IR_Name("y"))))
}

// Java: ParserTest.testAssignOr_es2020()
func TestAssignOr_es2020(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT_2020
	pt.expectFeatures(parser.FeatureSet_Feature_LOGICAL_ASSIGNMENT)

	pt.parseWarning("x||=y", requiresLanguageModeMessage(parser.FeatureSet_Feature_LOGICAL_ASSIGNMENT))
}

// Java: ParserTest.testAssignAnd()
func TestAssignAnd(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_LOGICAL_ASSIGNMENT)
	tree := pt.parse("x&&=y")
	assertNode(t, tree.GetFirstChild()).
		IsEqualTo(rhino.IR_ExprResult(rhino.IR_AssignAnd(rhino.IR_Name("x"), rhino.IR_Name("y"))))
}

// Java: ParserTest.testAssignAnd_es2020()
func TestAssignAnd_es2020(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT_2020
	pt.expectFeatures(parser.FeatureSet_Feature_LOGICAL_ASSIGNMENT)

	pt.parseWarning("x&&=y", requiresLanguageModeMessage(parser.FeatureSet_Feature_LOGICAL_ASSIGNMENT))
}

// Java: ParserTest.testAssignCoalesce()
func TestAssignCoalesce(t *testing.T) {
	pt := newParserTest(t)

	pt.expectFeatures(parser.FeatureSet_Feature_LOGICAL_ASSIGNMENT)
	tree := pt.parse("x??=y")
	assertNode(t, tree.GetFirstChild()).
		IsEqualTo(rhino.IR_ExprResult(rhino.IR_AssignCoalesce(rhino.IR_Name("x"), rhino.IR_Name("y"))))
}

// Java: ParserTest.testAssignCoalesce_es2020()
func TestAssignCoalesce_es2020(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT_2020
	pt.expectFeatures(parser.FeatureSet_Feature_LOGICAL_ASSIGNMENT)

	pt.parseWarning("x??=y", requiresLanguageModeMessage(parser.FeatureSet_Feature_LOGICAL_ASSIGNMENT))
}

// Java: ParserTest.testNoDuplicateComments_arrow_fn()
func TestNoDuplicateComments_arrow_fn(t *testing.T) {
	pt := newParserTest(t)

	pt.isIdeMode = true
	pt.parsingMode = Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	comments := pt.parseComments("const a = (/** number */ n) => {}")

	assertThat(t, comments).HasSize(1)
	assertThat(t, comments[0].Value).IsEqualTo("/** number */")
}

// Java: ParserTest.testIndirectCallName()
func TestIndirectCallName(t *testing.T) {
	pt := newParserTest(t)

	script := pt.parse("(0, foo)();")
	assertNode(t, script).HasToken(rhino.Token_SCRIPT)
	exprResult := script.GetFirstChild()
	assertNode(t, exprResult).HasToken(rhino.Token_EXPR_RESULT)

	// the `(0, foo)()` should have been converted to `foo()` with the
	// `FREE_CALL` property.
	call := exprResult.GetFirstChild()
	assertNode(t, call).IsFreeCall().HasFirstChildThat().IsName("foo")
}

// Java: ParserTest.testIndirectCallGetProp()
func TestIndirectCallGetProp(t *testing.T) {
	pt := newParserTest(t)

	script := pt.parse("(0, foo.bar)();")
	assertNode(t, script).HasToken(rhino.Token_SCRIPT)
	exprResult := script.GetFirstChild()
	assertNode(t, exprResult).HasToken(rhino.Token_EXPR_RESULT)

	// the `(0, foo)()` should have been converted to `foo.bar()` with the
	// `FREE_CALL` property, so it will actually get printed as `(0, foo.bar)()`.
	call := exprResult.GetFirstChild()
	assertNode(t, call).IsFreeCall().HasFirstChildThat().MatchesQualifiedName("foo.bar")
}

// Java: ParserTest.testFreeCall1()
func TestFreeCall1(t *testing.T) {
	pt := newParserTest(t)

	script := pt.parse("foo();")
	assertNode(t, script).HasToken(rhino.Token_SCRIPT)
	firstExpr := script.GetFirstChild()
	call := firstExpr.GetFirstChild()
	assertNode(t, call).HasToken(rhino.Token_CALL)

	assertNode(t, call).IsFreeCall()
}

// Java: ParserTest.testFreeCall2()
func TestFreeCall2(t *testing.T) {
	pt := newParserTest(t)

	script := pt.parse("x.foo();")
	assertNode(t, script).HasToken(rhino.Token_SCRIPT)
	firstExpr := script.GetFirstChild()
	call := firstExpr.GetFirstChild()
	assertNode(t, call).HasToken(rhino.Token_CALL)

	assertNode(t, call).IsNotFreeCall()
}

// Java: ParserTest.testTaggedTemplateFreeCall1()
func TestTaggedTemplateFreeCall1(t *testing.T) {
	pt := newParserTest(t)

	script := pt.parse("foo``;")
	assertNode(t, script).HasToken(rhino.Token_SCRIPT)
	firstExpr := script.GetFirstChild()
	call := firstExpr.GetFirstChild()
	assertNode(t, call).HasToken(rhino.Token_TAGGED_TEMPLATELIT)

	assertNode(t, call).IsFreeCall()
}

// Java: ParserTest.testTaggedTemplateFreeCall2()
func TestTaggedTemplateFreeCall2(t *testing.T) {
	pt := newParserTest(t)

	script := pt.parse("x.foo``;")
	assertNode(t, script).HasToken(rhino.Token_SCRIPT)
	firstExpr := script.GetFirstChild()
	call := firstExpr.GetFirstChild()
	assertNode(t, call).HasToken(rhino.Token_TAGGED_TEMPLATELIT)

	assertNode(t, call).IsNotFreeCall()
}

// Java: ParserTest.optionalFreeCall1()
func TestOptionalFreeCall1(t *testing.T) {
	pt := newParserTest(t)

	script := pt.parse("foo?.();")
	assertNode(t, script).HasToken(rhino.Token_SCRIPT)
	firstExpr := script.GetFirstChild()
	call := firstExpr.GetFirstChild()
	assertNode(t, call).HasToken(rhino.Token_OPTCHAIN_CALL)

	assertNode(t, call).IsFreeCall()
}

// Java: ParserTest.optChainFreeCall()
func TestOptChainFreeCall(t *testing.T) {
	pt := newParserTest(t)

	script := pt.parse("x?.foo();")
	assertNode(t, script).HasToken(rhino.Token_SCRIPT)
	firstExpr := script.GetFirstChild()
	call := firstExpr.GetFirstChild()
	assertNode(t, call).HasToken(rhino.Token_OPTCHAIN_CALL)

	assertNode(t, call).IsNotFreeCall()
}

// Java: ParserTest.testIsInClosureUnawareSubtreeProperty_isCorrectlySet()
func TestIsInClosureUnawareSubtreeProperty_isCorrectlySet(t *testing.T) {
	pt := newParserTest(t)

	script :=
		pt.parseWarning(
			`/**
 * @fileoverview
 * @closureUnaware
 */
goog.module('a.b');
/** @closureUnaware */
(function() {
  const x = 5;
}).call(globalThis);
let string1 = "ClosureAwareString";
/** @closureUnaware */
(function() {
  const y = 5;
}).call(globalThis);
`,
			"@closureUnaware annotation is not allowed in this compilation",
			"@closureUnaware annotation is not allowed in this compilation",
			"@closureUnaware annotation is not allowed in this compilation")

	moduleBody := script.GetFirstChild()

	firstUnawareFunctionExprResult := moduleBody.GetSecondChild()
	awareLet := firstUnawareFunctionExprResult.GetNext()
	secondUnawareFunctionExprResult := awareLet.GetNext()
	firstUnawareShadowHost :=
		firstUnawareFunctionExprResult.GetFirstFirstChild().GetFirstChild()
	secondUnawareShadowHost :=
		secondUnawareFunctionExprResult.GetFirstFirstChild().GetFirstChild()
	firstShadowRoot := firstUnawareShadowHost.GetClosureUnawareShadow()
	secondShadowRoot := secondUnawareShadowHost.GetClosureUnawareShadow()
	firstUnawareFunctionFunctionNode :=
		codeSubTree_findFirstNode(firstShadowRoot, (*rhino.Node).IsFunction)
	secondUnawareFunctionFunctionNode :=
		codeSubTree_findFirstNode(secondShadowRoot, (*rhino.Node).IsFunction)

	assertNode(t, firstUnawareShadowHost).IsName("$jscomp_wrap_closure_unaware_code")
	assertNode(t, secondUnawareShadowHost).IsName("$jscomp_wrap_closure_unaware_code")
	assertThat(t, firstUnawareShadowHost.GetClosureUnawareShadow()).IsNotNull()
	assertThat(t, firstUnawareShadowHost.GetIsInClosureUnawareSubtree()).IsFalse()
	assertThat(t, secondUnawareShadowHost.GetClosureUnawareShadow()).IsNotNull()
	assertThat(t, secondUnawareShadowHost.GetIsInClosureUnawareSubtree()).IsFalse()

	assertThat(t, firstUnawareFunctionFunctionNode.GetIsInClosureUnawareSubtree()).IsTrue()
	assertThat(t, secondUnawareFunctionFunctionNode.GetIsInClosureUnawareSubtree()).IsTrue()

	assertThat(t, awareLet.GetIsInClosureUnawareSubtree()).IsFalse()

	assertThat(t,
		rhino.Node_ValidateMemorySensitivePropertyGuarantees(
			firstUnawareFunctionFunctionNode,
			true,
			firstUnawareFunctionFunctionNode.GetFirstChild(),
			true)).
		IsTrue()
	assertThat(t,
		rhino.Node_ValidateMemorySensitivePropertyGuarantees(
			secondUnawareFunctionFunctionNode,
			true,
			secondUnawareFunctionFunctionNode.GetFirstChild(),
			true)).
		IsTrue()
	assertThat(t,
		rhino.Node_ValidateMemorySensitivePropertyGuarantees(
			firstUnawareFunctionFunctionNode, true, secondUnawareFunctionFunctionNode, true)).
		IsTrue()

	assertThat(t,
		rhino.Node_ValidateMemorySensitivePropertyGuarantees(
			firstUnawareFunctionFunctionNode, true, awareLet, false)).
		IsTrue()
	assertThat(t,
		rhino.Node_ValidateMemorySensitivePropertyGuarantees(
			secondUnawareFunctionFunctionNode, true, awareLet, false)).
		IsTrue()
	assertThat(t,
		rhino.Node_ValidateMemorySensitivePropertyGuarantees(
			firstUnawareFunctionFunctionNode, true, awareLet, false)).
		IsTrue()
}

// Java: ParserTest.closureUnawareAnnotation_specifyMinificationLevel_inFileOverview_succeeds()
func TestClosureUnawareAnnotation_specifyMinificationLevel_inFileOverview_succeeds(t *testing.T) {
	pt := newParserTest(t)

	errorReporter :=
		newTestErrorReporter().
			expectAllWarnings(
				"@closureUnaware annotation is not allowed in this compilation",
				"@closureUnaware annotation is not allowed in this compilation")
	pt.doParseWithReporter(
		`/** @fileoverview @closureUnaware {SIMPLE} */
/** @closureUnaware */ (
  function() {}).call(globalThis);
`,
		errorReporter)
}

// Java: ParserTest.closureUnawareAnnotation_specifyMinificationLevel_inPerFunctionDoc_errors()
func TestClosureUnawareAnnotation_specifyMinificationLevel_inPerFunctionDoc_errors(t *testing.T) {
	pt := newParserTest(t)

	errorReporter :=
		newTestErrorReporter().
			expectAllErrors("@closureUnaware mode can only be specified at the fileoverview level").
			expectAllWarnings(
				"@closureUnaware annotation is not allowed in this compilation",
				"@closureUnaware annotation is not allowed in this compilation")
	pt.doParseWithReporter(
		`/** @fileoverview @closureUnaware */
/** @closureUnaware {SIMPLE} */ (
  function() {}).call(globalThis);
`,
		errorReporter)
}

// Java: ParserTest.closureUnawareAnnotation_specifyMinificationLevel_inMultipleDoc_errors()
func TestClosureUnawareAnnotation_specifyMinificationLevel_inMultipleDoc_errors(t *testing.T) {
	pt := newParserTest(t)

	errorReporter :=
		newTestErrorReporter().
			expectAllErrors("@closureUnaware mode can only be specified at the fileoverview level").
			expectAllWarnings(
				"@closureUnaware annotation is not allowed in this compilation",
				"@closureUnaware annotation is not allowed in this compilation")
	pt.doParseWithReporter(
		`/** @fileoverview @closureUnaware {SIMPLE} */
/** @closureUnaware {SIMPLE} */ (
  function() {}).call(globalThis);
`,
		errorReporter)
}

// Java: ParserTest.closureUnawareAnnotation_usedOutsideFileWithClosureUnaware_errors()
func TestClosureUnawareAnnotation_usedOutsideFileWithClosureUnaware_errors(t *testing.T) {
	pt := newParserTest(t)

	errorReporter :=
		newTestErrorReporter().
			expectAllErrors(
				"invalid @closureUnaware node: must be in a file annotated with @closureUnaware").
			expectAllWarnings("@closureUnaware annotation is not allowed in this compilation")
	pt.doParseWithReporter(
		`/** @closureUnaware */ (
  function() {}).call(globalThis);
`,
		errorReporter)
}

// Java: ParserTest.closureUnawareAnnotation_onNonFunctionCall_errors()
func TestClosureUnawareAnnotation_onNonFunctionCall_errors(t *testing.T) {
	pt := newParserTest(t)

	pt.mode = Config_LanguageMode_ECMASCRIPT_2020
	pt.expectFeatures(parser.FeatureSet_Feature_CONST_DECLARATIONS)
	errorReporter :=
		newTestErrorReporter().
			expectAllErrors("@closureUnaware root must be function").
			expectAllWarnings(
				"@closureUnaware annotation is not allowed in this compilation",
				"@closureUnaware annotation is not allowed in this compilation")
	pt.doParseWithReporter(
		`/** @fileoverview @closureUnaware */
/** @closureUnaware */ const x = 0;
`,
		errorReporter)
}

// Java: ParserTest.testLinenoCharnoBinop(String)
func (pt *parserTest) testLinenoCharnoBinop(binop string) {
	t := pt.t
	t.Helper()
	op := pt.parse("var a = 89 " + binop + " 76;").GetFirstChild().GetFirstFirstChild()

	assertNode(t, op).HasLineno(1)
	assertNode(t, op).HasCharno(8)
}

// Java: ParserTest.assertNodeEquality(Node, Node)
func (pt *parserTest) assertNodeEquality(expected *rhino.Node, found *rhino.Node) {
	pt.t.Helper()
	assertNode(pt.t, found).IsEqualTo(expected)
}

// Java: ParserTest.createScript(Node)
func (pt *parserTest) createScript(n *rhino.Node) *rhino.Node {
	script := rhino.NewNode(rhino.Token_SCRIPT)
	script.AddChildToBack(n)
	return script
}

// Java: ParserTest.testMethodInObjectLiteral(String)
func (pt *parserTest) testMethodInObjectLiteral(js string) {
	pt.t.Helper()
	pt.mode = Config_LanguageMode_ECMASCRIPT_2015
	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parse(js)

	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.parseWarning(js, requiresLanguageModeMessage(parser.FeatureSet_Feature_MEMBER_DECLARATIONS))
}

// Java: ParserTest.testShorthandObjectProperty(String)
func (pt *parserTest) testShorthandObjectProperty(js string) {
	pt.t.Helper()
	pt.mode = Config_LanguageMode_ECMASCRIPT_2015
	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parse(js)

	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.parseWarning(js, requiresLanguageModeMessage(parser.FeatureSet_Feature_SHORTHAND_OBJECT_PROPERTIES))
}

// Java: ParserTest.testComputedProperty(String)
func (pt *parserTest) testComputedProperty(js string) {
	pt.t.Helper()
	pt.mode = Config_LanguageMode_ECMASCRIPT_2015
	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parse(js)

	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.parseWarning(js, requiresLanguageModeMessage(parser.FeatureSet_Feature_COMPUTED_PROPERTIES))
}

// Java: ParserTest.testTemplateLiteral(String)
func (pt *parserTest) testTemplateLiteral(s string) *rhino.Node {
	pt.t.Helper()
	pt.mode = Config_LanguageMode_ECMASCRIPT5
	pt.strictMode = Config_StrictMode_SLOPPY
	pt.parseWarning(s, requiresLanguageModeMessage(parser.FeatureSet_Feature_TEMPLATE_LITERALS))

	pt.mode = Config_LanguageMode_ECMASCRIPT_2015
	return pt.parse(s)
}

// Java: ParserTest.assertSimpleTemplateLiteral(String, String)
func (pt *parserTest) assertSimpleTemplateLiteral(expectedContents string, literal string) {
	t := pt.t
	t.Helper()
	node := pt.testTemplateLiteral(literal).GetFirstFirstChild()
	assertNode(t, node).HasType(rhino.Token_TEMPLATELIT)
	assertThat(t, node.GetChildCount()).IsEqualTo(1)
	assertNode(t, node.GetFirstChild()).HasType(rhino.Token_TEMPLATELIT_STRING)
	assertThat(t, node.GetFirstChild().GetCookedString()).IsEqualTo(expectedContents)
}

// Java: ParserTest.doAsyncArrowFunctionTest(String)
func (pt *parserTest) doAsyncArrowFunctionTest(arrowFunctionSource string) {
	pt.t.Helper()
	pt.expectFeatures(parser.FeatureSet_Feature_ASYNC_FUNCTIONS, parser.FeatureSet_Feature_ARROW_FUNCTIONS)

	for _, m := range Config_LanguageMode_values() {
		pt.mode = m
		if m == Config_LanguageMode_ECMASCRIPT3 {
			pt.strictMode = Config_StrictMode_SLOPPY
		} else {
			pt.strictMode = Config_StrictMode_STRICT
		}
		if m.FeatureSet().Has(parser.FeatureSet_Feature_ASYNC_FUNCTIONS) {
			pt.parse(arrowFunctionSource)
		} else if m.FeatureSet().Has(parser.FeatureSet_Feature_ARROW_FUNCTIONS) {
			pt.parseWarning(arrowFunctionSource, requiresLanguageModeMessage(parser.FeatureSet_Feature_ASYNC_FUNCTIONS))
		} else {
			pt.parseWarning(
				arrowFunctionSource,
				requiresLanguageModeMessage(parser.FeatureSet_Feature_ARROW_FUNCTIONS),
				requiresLanguageModeMessage(parser.FeatureSet_Feature_ASYNC_FUNCTIONS))
		}
	}
}

// Java: ParserTest.freeCall(Node)
func (pt *parserTest) freeCall(n *rhino.Node) *rhino.Node {
	n.PutBooleanProp(rhino.Node_FREE_CALL, true)
	return n
}

// Java: ParserTest.assertNodeHasJSDocInfoWithJSType(Node, JSType)
func (pt *parserTest) assertNodeHasJSDocInfoWithJSType(node *rhino.Node, jsType rhino.JSType) {
	t := pt.t
	t.Helper()
	info := node.GetJSDocInfo()
	assertWithMessage(t, jutil.LenientFormat("Node has no JSDocInfo: %s", node)).That(info).IsNotNull()
	pt.AssertTypeEqualsExpr(jsType, info.GetType())
}

// Java: ParserTest.assertNodeHasJSDocInfoWithNoJSType(Node)
func (pt *parserTest) assertNodeHasJSDocInfoWithNoJSType(node *rhino.Node) {
	t := pt.t
	t.Helper()
	info := node.GetJSDocInfo()
	assertWithMessage(t, jutil.LenientFormat("Node has no JSDocInfo: %s", node)).That(info).IsNotNull()
	type_ := info.GetType()
	assertWithMessage(t, "JSDoc unexpectedly has type").That(type_).IsNull()
}

// Java: ParserTest.assertNodeHasNoJSDocInfo(Node)
func (pt *parserTest) assertNodeHasNoJSDocInfo(node *rhino.Node) {
	t := pt.t
	t.Helper()
	info := node.GetJSDocInfo()
	assertWithMessage(t, jutil.LenientFormat("Node %s has unexpected JSDocInfo %s", node, info)).That(info).IsNull()
}

// Java: ParserTest.requiresLanguageModeMessage(Feature)
func requiresLanguageModeMessage(feature parser.FeatureSet_Feature) string {
	return IRFactory_LanguageFeatureWarningMessage(feature)
}

// Java: ParserTest.script(Node)
func (pt *parserTest) script(stmt *rhino.Node) *rhino.Node {
	return rhino.NewNodeWithChild(rhino.Token_SCRIPT, stmt)
}

// Java: ParserTest.expr(Node)
func (pt *parserTest) expr(n *rhino.Node) *rhino.Node {
	return rhino.NewNodeWithChild(rhino.Token_EXPR_RESULT, n)
}

// Java: ParserTest.regex(String) / ParserTest.regex(String, String)
func (pt *parserTest) regex(regex string, flag ...string) *rhino.Node {
	if len(flag) == 0 {
		return rhino.NewNodeWithChild(rhino.Token_REGEXP, rhino.Node_NewString(regex))
	}
	return rhino.NewNodeWith2Children(rhino.Token_REGEXP, rhino.Node_NewString(regex), rhino.Node_NewString(flag[0]))
}

// parseError verifies that the given code has the given parse errors.
// Returns, if in IDE mode, a partial tree.
//
// Java: ParserTest.parseError(String, String...)
func (pt *parserTest) parseError(source string, errors ...string) *rhino.Node {
	t := pt.t
	t.Helper()
	assertThat(t, errors).IsNotEmpty()

	testErrorReporter := newTestErrorReporter().expectAllErrors(errors...)
	result := ParserRunner_Parse(
		rhino.NewSimpleSourceFile("input", rhino.StaticSourceFile_SourceKind_STRONG),
		source,
		pt.createConfig(),
		testErrorReporter)
	script := result.Ast

	// check expected features if specified
	assertFS(t, result.Features).Contains(pt.expectedFeatures)

	// verifying that all errors were seen
	testErrorReporter.verifyHasEncounteredAllWarningsAndErrors(t)

	return script
}

// parseWarning verifies that the given code has the given parse warnings.
// Returns the parse tree.
//
// Java: ParserTest.parseWarning(String, String...)
func (pt *parserTest) parseWarning(string_ string, warnings ...string) *rhino.Node {
	pt.t.Helper()
	return pt.doParse(string_, warnings...).Ast
}

// Java: ParserTest.doParse(String, String...)
func (pt *parserTest) doParse(string_ string, warnings ...string) *ParserRunner_ParseResult {
	pt.t.Helper()
	testErrorReporter := newTestErrorReporter().expectAllWarnings(warnings...)
	return pt.doParseWithReporter(string_, testErrorReporter)
}

// Java: ParserTest.doParse(String, TestErrorReporter)
func (pt *parserTest) doParseWithReporter(string_ string, testErrorReporter *testErrorReporter) *ParserRunner_ParseResult {
	t := pt.t
	t.Helper()
	file := rhino.NewSimpleSourceFile("input", rhino.StaticSourceFile_SourceKind_STRONG)
	result := ParserRunner_Parse(file, string_, pt.createConfig(), testErrorReporter)

	// check expected features if specified
	assertFS(t, result.Features).Contains(pt.expectedFeatures)

	// verifying that all warnings were seen
	testErrorReporter.verifyHasEncounteredAllWarningsAndErrors(t)
	pt.assertSourceInfoPresent(result.Ast)
	return result
}

// Java: ParserTest.assertSourceInfoPresent(Node)
func (pt *parserTest) assertSourceInfoPresent(node *rhino.Node) {
	t := pt.t
	t.Helper()
	deque := []*rhino.Node{node}

	for len(deque) != 0 {
		node = deque[0]
		deque = deque[1:]

		assertWithMessage(t, jutil.LenientFormat("Source information must be present on %s", node)).
			That(node.GetLineno()).IsAtLeast(0)

		assertWithMessage(t, jutil.LenientFormat("Source length must be nonnegative on %s", node)).
			That(node.GetLength()).IsAtLeast(0)

		for child := node.GetFirstChild(); child != nil; child = child.GetNext() {
			deque = append(deque, child)
		}
	}
}

// parse verifies that the given code has no parse warnings or errors.
// Returns the parse tree.
//
// Java: ParserTest.parse(String)
func (pt *parserTest) parse(string_ string) *rhino.Node {
	pt.t.Helper()
	return pt.parseWarning(string_)
}

// parseComments returns all comments recorded by the parser in IDE mode.
//
// Assumes `isIdeMode` is true and an appropriate `parsingMode` is
// configured.
//
// Java: ParserTest.parseComments(String)
func (pt *parserTest) parseComments(string_ string) []*parser.Comment {
	pt.t.Helper()
	return pt.doParse(string_).Comments
}

// Java: ParserTest.createConfig()
func (pt *parserTest) createConfig() *Config {
	if pt.isIdeMode {
		return ParserRunner_CreateConfigWithModes(
			pt.mode, pt.parsingMode, Config_RunMode_KEEP_GOING, nil, true, pt.strictMode)
	} else {
		return ParserRunner_CreateConfig(pt.mode, nil, pt.strictMode)
	}
}

// expectFeatures sets expectedFeatures based on the list of features.
//
// Java: ParserTest.expectFeatures(Feature...)
func (pt *parserTest) expectFeatures(features ...parser.FeatureSet_Feature) {
	pt.expectedFeatures = parser.FeatureSet_BARE_MINIMUM.WithFeatures(features...)
}

// parserTest_ParserResult is ParserTest.ParserResult.
type parserTest_ParserResult struct {
	code string
	node *rhino.Node
}

// ---------------------------------------------------------------------------
// Truth subjects used by ParserTest beyond those of JsDocInfoParserTest.

// IsAtLeast is ComparableSubject.isAtLeast for an int.
func (s *subject) IsAtLeast(min int) {
	s.t.Helper()
	if v, ok := truthNorm(s.actual).(int64); !ok || v < int64(min) {
		s.fail("expected to be at least %d, was %v", min, s.actual)
	}
}

// HasMessageThat is ThrowableSubject.hasMessageThat().
func (s *subject) HasMessageThat() *subject {
	s.t.Helper()
	t, ok := s.actual.(jutil.Throwable)
	if !ok {
		s.fail("expected a Throwable, was %v", s.actual)
	}
	return &subject{t: s.t, actual: t.GetMessage(), msg: s.msg}
}

// IsLessThan is ComparableSubject.isLessThan for an integer.
func (s *subject) IsLessThan(max int64) {
	s.t.Helper()
	if v, ok := truthNorm(s.actual).(int64); !ok || v >= max {
		s.fail("expected to be less than %d, was %v", max, s.actual)
	}
}

// ---------------------------------------------------------------------------
// jscomp/testing/CodeSubTree and jscomp.NodeUtil helpers (jscomp is not
// ported yet).

// codeSubTree_findFirstNode is CodeSubTree.findFirstNode(Node, Predicate):
// the first match in a pre-order traversal (NodeUtil.visitPreOrder).
func codeSubTree_findFirstNode(rootNode *rhino.Node, predicate func(*rhino.Node) bool) *rhino.Node {
	var matches []*rhino.Node
	var visit func(n *rhino.Node)
	visit = func(n *rhino.Node) {
		if predicate(n) {
			matches = append(matches, n)
		}
		for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
			visit(c)
		}
	}
	visit(rootNode)
	jutil.CheckStateMsg(len(matches) != 0, "no nodes found")
	return matches[0]
}

// nodeUtil_getClassMembers is NodeUtil.getClassMembers(Node).
func nodeUtil_getClassMembers(n *rhino.Node) *rhino.Node {
	jutil.CheckArgument(n.IsClass())
	return n.GetLastChild()
}
