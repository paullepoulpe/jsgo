// Port of com/google/javascript/jscomp/CodePrinterTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/CodePrinterTest.java)

/*
 * Copyright 2004 The Closure Compiler Authors.
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
	"fmt"
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CodePrinterTest is the JUnit 4 class CodePrinterTest.
type CodePrinterTest struct {
	CodePrinterTestBase
}

func NewCodePrinterTest() *CodePrinterTest {
	c := &CodePrinterTest{}
	c.InitCodePrinterTestBase()
	return c
}

func TestCodePrinterTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCodePrinterTest() })
}

// Java: CodePrinterTest.testBigInt()
func (c *CodePrinterTest) TestBigInt() {
	c.languageMode = CompilerOptions_LanguageMode_UNSUPPORTED
	c.AssertPrintSame("1n")
	c.AssertPrint("0b10n", "2n")
	c.AssertPrint("0o3n", "3n")
	c.AssertPrint("0x4n", "4n")
	c.AssertPrintSame("-5n")
	c.AssertPrint("-0b110n", "-6n")
	c.AssertPrint("-0o7n", "-7n")
	c.AssertPrint("-0x8n", "-8n")
	c.AssertPrintSame("1000000000n")
	c.AssertPrintSame("-10000000000n")
	c.AssertPrint(
		"7182582809266874004429817360811601465721521679690913024666851123813995272999n",
		"0xfe132a356ec5f9279946c8fdf29780abd0387a8277e811069d174660b385b27n")
	c.AssertPrint(
		"0x000000fe132a356ec5f9279946c8fdf29780abd0387a8277e811069d174660b385b27n",
		"0xfe132a356ec5f9279946c8fdf29780abd0387a8277e811069d174660b385b27n")
	c.AssertPrint(
		"-7182582809266874004429817360811601465721521679690913024666851123813995272999n",
		"-0xfe132a356ec5f9279946c8fdf29780abd0387a8277e811069d174660b385b27n")
	c.AssertPrint(
		"-0o00774114521526730576223631215443757451360052750070365011677201040647213506301316055447n",
		"-0xfe132a356ec5f9279946c8fdf29780abd0387a8277e811069d174660b385b27n")
	c.AssertPrint(
		"0o774114521526730576223631215443757451360052750070365011677201040647213506301316055447n",
		"0xfe132a356ec5f9279946c8fdf29780abd0387a8277e811069d174660b385b27n")
	c.AssertPrint(
		"0b111111100001001100101010001101010110111011000101111110010010011110011001010001101100100011111101111100101001011110000000101010111101000000111000011110101000001001110111111010000001000100000110100111010001011101000110011000001011001110000101101100100111n",
		"0xfe132a356ec5f9279946c8fdf29780abd0387a8277e811069d174660b385b27n")
}

// Java: CodePrinterTest.testTrailingCommaInArrayAndObjectWithPrettyPrint()
func (c *CodePrinterTest) TestTrailingCommaInArrayAndObjectWithPrettyPrint() {
	c.AssertPrettyPrintSame("({a:1, b:2,});\n")
	c.AssertPrettyPrintSame("[1, 2, 3,];\n")
	// An array starting with a hole is printed ideally but this is very rare.
	c.AssertPrettyPrintSame("[, ];\n")
}

// Java: CodePrinterTest.testTrailingCommaInArrayAndObjectWithoutPrettyPrint()
func (c *CodePrinterTest) TestTrailingCommaInArrayAndObjectWithoutPrettyPrint() {
	c.AssertPrint("({a:1, b:2,})", "({a:1,b:2})")
	c.AssertPrint("[1, 2, 3,]", "[1,2,3]")

	// When the last element is empty,  the trailing comma must be kept.
	c.AssertPrintSame("[,]")   // same as `[undefined]`
	c.AssertPrintSame("[a,,]") // same as `[a, undefined]`
}

// Java: CodePrinterTest.testNoTrailingCommaInEmptyArrayLiteral()
func (c *CodePrinterTest) TestNoTrailingCommaInEmptyArrayLiteral() {
	// In cases where we modify the AST we might wind up with an array literal that has no elements
	// yet still has a trailing comma. This is meant to test for that. We need to build the tree
	// manually because an array literal with no elements and a trailing comma has a different
	// meaning: it represents a single undefined element.
	arrLit := rhino.IR_Arraylit()
	arrLit.SetTrailingComma(true)
	c.AssertPrettyPrintNode("[]", arrLit)
}

// Java: CodePrinterTest.testNoTrailingCommaInEmptyObjectLiteral()
func (c *CodePrinterTest) TestNoTrailingCommaInEmptyObjectLiteral() {
	// In cases where we modify the AST we might wind up with an object literal that has no elements
	// yet still has a trailing comma. This is meant to test for that. We need to build the tree
	// manually because an object literal with no elements and a trailing comma is a syntax error.
	objLit := rhino.IR_Objectlit()
	objLit.SetTrailingComma(true)
	c.AssertPrettyPrintNode("{}", objLit)
}

// Java: CodePrinterTest.testNoTrailingCommaInEmptyParamList()
func (c *CodePrinterTest) TestNoTrailingCommaInEmptyParamList() {
	// In cases where we modify the AST we might wind up with a parameter list that has no elements
	// yet still has a trailing comma. This is meant to test for that. We need to build the tree
	// manually because a parameter list with no elements and a trailing comma is a syntax error.
	paramList := rhino.IR_ParamList()
	rhino.IR_Function(rhino.IR_Name("f"), paramList, rhino.IR_Block())
	paramList.SetTrailingComma(true)
	c.AssertPrettyPrintNode("()", paramList)
}

// Java: CodePrinterTest.testNoTrailingCommaInEmptyCall()
func (c *CodePrinterTest) TestNoTrailingCommaInEmptyCall() {
	// In cases where we modify the AST we might wind up with a call node that has no elements
	// yet still has a trailing comma. This is meant to test for that. We need to build the tree
	// manually because a call node with no elements and a trailing comma is a syntax error.
	call := rhino.IR_Call(rhino.IR_Name("f"))
	call.SetTrailingComma(true)
	c.AssertPrettyPrintNode("f()", call)
}

// Java: CodePrinterTest.testNoTrailingCommaInEmptyOptChainCall()
func (c *CodePrinterTest) TestNoTrailingCommaInEmptyOptChainCall() {
	// In cases where we modify the AST we might wind up with an optional chain call node that has
	// no elements yet still has a trailing comma. This is meant to test for that. We need to build
	// the tree manually because an optional chain call node with no elements and a trailing comma
	// is a syntax error.
	optChainCall := rhino.IR_StartOptChainCall(rhino.IR_Name("f"))
	optChainCall.SetTrailingComma(true)
	c.AssertPrettyPrintNode("f?.()", optChainCall)
}

// Java: CodePrinterTest.testNoTrailingCommaInEmptyNew()
func (c *CodePrinterTest) TestNoTrailingCommaInEmptyNew() {
	// In cases where we modify the AST we might wind up with a new node that has no elements
	// yet still has a trailing comma. This is meant to test for that. We need to build the tree
	// manually because a new node with no elements and a trailing comma is a syntax error.
	newNode := rhino.IR_NewNode(rhino.IR_Name("f"))
	newNode.SetTrailingComma(true)
	c.AssertPrettyPrintNode("new f()", newNode)
}

// Java: CodePrinterTest.testTrailingCommaInParameterListWithPrettyPrint()
func (c *CodePrinterTest) TestTrailingCommaInParameterListWithPrettyPrint() {
	c.languageMode = CompilerOptions_LanguageMode_UNSUPPORTED
	c.AssertPrettyPrintSame("function f(a, b,) {\n}\n")
	c.AssertPrettyPrintSame("f(1, 2,);\n")
	c.AssertPrettyPrintSame("f?.(1, 2,);\n")
	c.AssertPrettyPrintSame("let x = new Number(1,);\n")
}

// Java: CodePrinterTest.testTrailingCommaInParameterListWithoutPrettyPrint()
func (c *CodePrinterTest) TestTrailingCommaInParameterListWithoutPrettyPrint() {
	c.languageMode = CompilerOptions_LanguageMode_UNSUPPORTED
	c.AssertPrint("function f(a, b,) {}", "function f(a,b){}")
	c.AssertPrint("f(1, 2,);", "f(1,2)")
	c.AssertPrint("f?.(1, 2,);", "f?.(1,2)")
	c.AssertPrint("let x = new Number(1,);", "let x=new Number(1)")
}

// Java: CodePrinterTest.optChain()
func (c *CodePrinterTest) TestOptChain() {
	c.languageMode = CompilerOptions_LanguageMode_UNSUPPORTED
	c.AssertPrintSame("a.b?.c")
	c.AssertPrintSame("a.b?.[\"c\"]")
	c.AssertPrintSame("a.b?.()")
	c.AssertPrintSame("a?.b.c?.d")
	c.AssertPrintSame("(a?.b).c")
	c.AssertPrintSame("(a.b?.c.d).e")
	c.AssertPrintSame("(a?.[b])[c]")
	c.AssertPrintSame("(a.b?.())()")
}

// Java: CodePrinterTest.testUnescapedUnicodeLineSeparator_2018()
func (c *CodePrinterTest) TestUnescapedUnicodeLineSeparator_2018() {
	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT_2018
	c.AssertPrintSame("`\u2028`")

	c.AssertPrint("'\u2028'", "\"\\u2028\"")
	c.AssertPrint("\"\u2028\"", "\"\\u2028\"")

	c.outputCharset = jutil.StandardCharsets_UTF_8
	// printed as a unicode escape for ES_2018 output
	c.AssertPrint("'\u2028'", "\"\\u2028\"")
	c.AssertPrint("\"\u2028\"", "\"\\u2028\"")
}

// Java: CodePrinterTest.testUnescapedUnicodeLineSeparator_2019()
func (c *CodePrinterTest) TestUnescapedUnicodeLineSeparator_2019() {
	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT_2019
	c.AssertPrint("'\u2028'", "\"\\u2028\"")
	c.AssertPrint("\"\u2028\"", "\"\\u2028\"")

	c.outputCharset = jutil.StandardCharsets_UTF_8
	// left unescaped for ES_2019 out
	c.AssertPrint("'\u2028'", "\"\u2028\"")
	c.AssertPrint("\"\u2028\"", "\"\u2028\"")
}

// Java: CodePrinterTest.testUnescapedUnicodeParagraphSeparator_2018()
func (c *CodePrinterTest) TestUnescapedUnicodeParagraphSeparator_2018() {
	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT_2018
	c.AssertPrintSame("`\u2029`")

	c.AssertPrint("'\u2029'", "\"\\u2029\"")
	c.AssertPrint("\"\u2029\"", "\"\\u2029\"")

	c.outputCharset = jutil.StandardCharsets_UTF_8
	// printed as a unicode escape for ES_2018 output
	c.AssertPrint("'\u2029'", "\"\\u2029\"")
	c.AssertPrint("\"\u2029\"", "\"\\u2029\"")
}

// Java: CodePrinterTest.testUnescapedUnicodeParagraphSeparator_2019()
func (c *CodePrinterTest) TestUnescapedUnicodeParagraphSeparator_2019() {
	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT_2019

	c.AssertPrint("'\u2029'", "\"\\u2029\"")
	c.AssertPrint("\"\u2029\"", "\"\\u2029\"")

	c.outputCharset = jutil.StandardCharsets_UTF_8
	// left unescaped for ES_2019 out
	c.AssertPrint("'\u2029'", "\"\u2029\"")
	c.AssertPrint("\"\u2029\"", "\"\u2029\"")
}

// Java: CodePrinterTest.testOptionalCatchBlock()
func (c *CodePrinterTest) TestOptionalCatchBlock() {
	c.AssertPrintSame("try{}catch{}")
	c.AssertPrintSame("try{}catch{}finally{}")
}

// Java: CodePrinterTest.testExponentiationOperator()
func (c *CodePrinterTest) TestExponentiationOperator() {
	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT_2016
	c.AssertPrintSame("x**y")
	// Exponentiation is right associative
	c.AssertPrint("x**(y**z)", "x**y**z")
	c.AssertPrintSame("(x**y)**z")
	// parens are kept because ExponentiationExpression cannot expand to
	//	 UnaryExpression ** ExponentiationExpression
	c.AssertPrintSame("(-x)**y")
	// parens are kept because unary operators are higher precedence than '**'
	c.AssertPrintSame("-(x**y)")
	// parens are not needed for a unary operator on the right operand
	c.AssertPrint("x**(-y)", "x**-y")
	// NOTE: "-x**y" is a syntax error tested in ParserTest

	// ** has a higher precedence than /
	c.AssertPrint("x/(y**z)", "x/y**z")
	c.AssertPrintSame("(x/y)**z")
}

// Java: CodePrinterTest.testExponentiationAssignmentOperator()
func (c *CodePrinterTest) TestExponentiationAssignmentOperator() {
	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT_2016
	c.AssertPrintSame("x**=y")
}

// Java: CodePrinterTest.testNullishCoalesceOperator()
func (c *CodePrinterTest) TestNullishCoalesceOperator() {
	c.AssertPrintSame("x??y??z")
	// Nullish coalesce is left associative
	c.AssertPrintSame("x??(y??z)")
	c.AssertPrint("(x??y)??z", "x??y??z")
	// // parens are kept because logical AND and logical OR must be separated from '??'
	c.AssertPrintSame("(x&&y)??z")
	c.AssertPrintSame("(x??y)||z")
	c.AssertPrintSame("x??(y||z)")
	// NOTE: "x&&y??z" is a syntax error tested in ParserTest
}

// Java: CodePrinterTest.testNullishCoalesceOperator2()
func (c *CodePrinterTest) TestNullishCoalesceOperator2() {
	// | has higher precedence than ??
	c.AssertPrint("(a|b)??c", "a|b??c")
	c.AssertPrintSame("(a??b)|c")
	c.AssertPrintSame("a|(b??c)")
	c.AssertPrint("a??(b|c)", "a??b|c")
	// ?? has higher precedence than : ? (conditional)
	c.AssertPrint("(a??b)?(c??d):(e??f)", "a??b?c??d:e??f")
	c.AssertPrintSame("a??(b?c:d)")
	c.AssertPrintSame("(a?b:c)??d")
}

// Java: CodePrinterTest.testLogicalAssignmentOperator()
func (c *CodePrinterTest) TestLogicalAssignmentOperator() {
	c.AssertPrintSame("x||=y")
	c.AssertPrintSame("x&&=y")
	c.AssertPrintSame("x??=y")
}

// Java: CodePrinterTest.testObjectLiteralWithSpread()
func (c *CodePrinterTest) TestObjectLiteralWithSpread() {
	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT_NEXT
	c.AssertPrintSame("({...{}})")
	c.AssertPrintSame("({...x})")
	c.AssertPrintSame("({...x,a:1})")
	c.AssertPrintSame("({a:1,...x})")
	c.AssertPrintSame("({a:1,...x,b:1})")
	c.AssertPrintSame("({...x,...y})")
	c.AssertPrintSame("({...x,...f()})")
	c.AssertPrintSame("({...{...{}}})")
}

// Java: CodePrinterTest.testObjectLiteralWithComma()
func (c *CodePrinterTest) TestObjectLiteralWithComma() {
	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT_NEXT
	c.AssertPrintSame("({[(a,b)]:c})")
	c.AssertPrintSame("({a:(b,c)})")
	c.AssertPrintSame("({[(a,b)]:(c,d)})")
	c.AssertPrintSame("({[(a,b)]:c,[d]:(e,f)})")
}

// Java: CodePrinterTest.testPrint()
func (c *CodePrinterTest) TestPrint() {
	c.AssertPrint("10 + a + b", "10+a+b")
	c.AssertPrint("10 + (30*50)", "10+30*50")
	c.AssertPrint("with(x) { x + 3; }", "with(x)x+3")
	c.AssertPrint("\"aa'a\"", "\"aa'a\"")
	c.AssertPrint("\"aa\\\"a\"", "'aa\"a'")
	c.AssertPrint("function foo()\n{return 10;}", "function foo(){return 10}")
	c.AssertPrint("a instanceof b", "a instanceof b")
	c.AssertPrint("typeof(a)", "typeof a")
	c.AssertPrint(
		"var foo = x ? { a : 1 } : {a: 3, b:4, \"default\": 5, \"foo-bar\": 6}",
		"var foo=x?{a:1}:{a:3,b:4,\"default\":5,\"foo-bar\":6}")

	// Safari: needs ';' at the end of a throw statement
	c.AssertPrint("function foo(){throw 'error';}", "function foo(){throw\"error\";}")

	// The code printer does not eliminate unnecessary blocks.
	c.AssertPrint("var x = 10; { var y = 20; }", "var x=10;{var y=20}")

	c.AssertPrint("while (x-- > 0);", "while(x-- >0);")
	c.AssertPrint("x-- >> 1", "x-- >>1")

	c.AssertPrint("(function () {})(); ", "(function(){})()")

	// Associativity
	c.AssertPrint("var a,b,c,d;a || (b&& c) && (a || d)", "var a,b,c,d;a||b&&c&&(a||d)")
	c.AssertPrint(
		"var a,b,c; a || (b || c); a * (b * c); a | (b | c)",
		"var a,b,c;a||(b||c);a*(b*c);a|(b|c)")
	c.AssertPrint(
		"var a,b,c; a / b / c;a / (b / c); a - (b - c);", "var a,b,c;a/b/c;a/(b/c);a-(b-c)")

	// Nested assignments
	c.AssertPrint("var a,b; a = b = 3;", "var a,b;a=b=3")
	c.AssertPrint("var a,b,c,d; a = (b = c = (d = 3));", "var a,b,c,d;a=b=c=d=3")
	c.AssertPrint("var a,b,c; a += (b = c += 3);", "var a,b,c;a+=b=c+=3")
	c.AssertPrint("var a,b,c; a *= (b -= c);", "var a,b,c;a*=b-=c")

	// Precedence
	c.AssertPrint("a ? delete b[0] : 3", "a?delete b[0]:3")
	c.AssertPrint("(delete a[0])/10", "delete a[0]/10")

	// optional '()' for new

	// simple new
	c.AssertPrint("new A", "new A")
	c.AssertPrint("new A()", "new A")
	c.AssertPrint("new A('x')", "new A(\"x\")")

	// calling instance method directly after new
	c.AssertPrint("new A().a()", "(new A).a()")
	c.AssertPrint("(new A).a()", "(new A).a()")

	// this case should be fixed
	c.AssertPrint("new A('y').a()", "(new A(\"y\")).a()")

	// internal class
	c.AssertPrint("new A.B", "new A.B")
	c.AssertPrint("new A.B()", "new A.B")
	c.AssertPrint("new A.B('z')", "new A.B(\"z\")")

	// calling instance method directly after new internal class
	c.AssertPrint("(new A.B).a()", "(new A.B).a()")
	c.AssertPrint("new A.B().a()", "(new A.B).a()")
	// this case should be fixed
	c.AssertPrint("new A.B('w').a()", "(new A.B(\"w\")).a()")

	// calling new on the result of a call
	c.AssertPrintSame("new (a())")
	c.AssertPrint("new (a())()", "new (a())")
	c.AssertPrintSame("new (a.b())")
	c.AssertPrint("new (a.b())()", "new (a.b())")

	// Operators: make sure we don't convert binary + and unary + into ++
	c.AssertPrint("x + +y", "x+ +y")
	c.AssertPrint("x - (-y)", "x- -y")
	c.AssertPrint("x++ +y", "x++ +y")
	c.AssertPrint("x-- -y", "x-- -y")
	c.AssertPrint("x++ -y", "x++-y")

	// Label
	c.AssertPrint("foo:for(;;){break foo;}", "foo:for(;;)break foo")
	c.AssertPrint("foo:while(1){continue foo;}", "foo:while(1)continue foo")
	c.AssertPrintSame("foo:;")
	c.AssertPrint("foo: {}", "foo:;")

	// Object literals.
	c.AssertPrint("({})", "({})")
	c.AssertPrint("var x = {};", "var x={}")
	c.AssertPrint("({}).x", "({}).x")
	c.AssertPrint("({})['x']", "({})[\"x\"]")
	c.AssertPrint("({}) instanceof Object", "({})instanceof Object")
	c.AssertPrint("({}) || 1", "({})||1")
	c.AssertPrint("1 || ({})", "1||{}")
	c.AssertPrint("({}) ? 1 : 2", "({})?1:2")
	c.AssertPrint("0 ? ({}) : 2", "0?{}:2")
	c.AssertPrint("0 ? 1 : ({})", "0?1:{}")
	c.AssertPrint("typeof ({})", "typeof{}")
	c.AssertPrint("f({})", "f({})")

	// Anonymous function expressions.
	c.AssertPrint("(function(){})", "(function(){})")
	c.AssertPrint("(function(){})()", "(function(){})()")
	c.AssertPrint("(function(){})instanceof Object", "(function(){})instanceof Object")
	c.AssertPrint("(function(){}).bind().call()", "(function(){}).bind().call()")
	c.AssertPrint("var x = function() { };", "var x=function(){}")
	c.AssertPrint("var x = function() { }();", "var x=function(){}()")
	c.AssertPrint("(function() {}), 2", "(function(){}),2")

	// Name functions expression.
	c.AssertPrint("(function f(){})", "(function f(){})")

	// Function declaration.
	c.AssertPrint("function f(){}", "function f(){}")

	// Make sure we don't treat non-Latin character escapes as raw strings.
	c.AssertPrint("({ 'a': 4, '\\u0100': 4 })", "({\"a\":4,\"\\u0100\":4})")
	c.AssertPrint("({ a: 4, '\\u0100': 4 })", "({a:4,\"\\u0100\":4})")

	// Test if statement and for statements with single statements in body.
	c.AssertPrint("if (true) { alert();}", "if(true)alert()")
	c.AssertPrint("if (false) {} else {alert(\"a\");}", "if(false);else alert(\"a\")")
	c.AssertPrint("for(;;) { alert();};", "for(;;)alert()")

	c.AssertPrint("do { alert(); } while(true);", "do alert();while(true)")
	c.AssertPrint("myLabel: { alert();}", "myLabel:alert()")
	c.AssertPrint("myLabel: for(;;) continue myLabel;", "myLabel:for(;;)continue myLabel")

	// Test nested var statement
	c.AssertPrint("if (true) var x; x = 4;", "if(true)var x;x=4")

	// Non-latin identifier. Make sure we keep them escaped.
	c.AssertPrint("\\u00fb", "\\u00fb")
	c.AssertPrint("\\u00fa=1", "\\u00fa=1")
	c.AssertPrint("function \\u00f9(){}", "function \\u00f9(){}")
	c.AssertPrint("x.\\u00f8", "x.\\u00f8")
	c.AssertPrint("x.\\u00f8", "x.\\u00f8")
	c.AssertPrint("abc\\u4e00\\u4e01jkl", "abc\\u4e00\\u4e01jkl")

	// Test the right-associative unary operators for spurious parens
	c.AssertPrint("! ! true", "!!true")
	c.AssertPrint("!(!(true))", "!!true")
	c.AssertPrint("typeof(void(0))", "typeof void 0")
	c.AssertPrint("typeof(void(!0))", "typeof void!0")
	c.AssertPrint("+ - + + - + 3", "+-+ +-+3") // chained unary plus/minus
	c.AssertPrint("+(--x)", "+--x")
	c.AssertPrint("-(++x)", "-++x")

	// needs a space to prevent an ambiguous parse
	c.AssertPrint("-(--x)", "- --x")
	c.AssertPrint("!(~~5)", "!~~5")
	c.AssertPrint("~(a/b)", "~(a/b)")

	// Preserve parens to overcome greedy binding of NEW
	c.AssertPrint("new (foo.bar()).factory(baz)", "new (foo.bar().factory)(baz)")
	c.AssertPrint("new (bar()).factory(baz)", "new (bar().factory)(baz)")
	c.AssertPrint("new (new foobar(x)).factory(baz)", "new (new foobar(x)).factory(baz)")

	// Make sure that HOOK is right associative
	c.AssertPrint("a ? b : (c ? d : e)", "a?b:c?d:e")
	c.AssertPrint("a ? (b ? c : d) : e", "a?b?c:d:e")
	c.AssertPrint("(a ? b : c) ? d : e", "(a?b:c)?d:e")

	// Test nested ifs
	c.AssertPrint("if (x) if (y); else;", "if(x)if(y);else;")

	// Test comma.
	c.AssertPrint("a,b,c", "a,b,c")
	c.AssertPrint("(a,b),c", "a,b,c")
	c.AssertPrint("a,(b,c)", "a,b,c")
	c.AssertPrint("x=a,b,c", "x=a,b,c")
	c.AssertPrint("x=(a,b),c", "x=(a,b),c")
	c.AssertPrint("x=a,(b,c)", "x=a,b,c")
	c.AssertPrint("x=a,y=b,z=c", "x=a,y=b,z=c")
	c.AssertPrint("x=(a,y=b,z=c)", "x=(a,y=b,z=c)")
	c.AssertPrint("x=[a,b,c,d]", "x=[a,b,c,d]")
	c.AssertPrint("x=[(a,b,c),d]", "x=[(a,b,c),d]")
	c.AssertPrint("x=[(a,(b,c)),d]", "x=[(a,b,c),d]")
	c.AssertPrint("x=[a,(b,c,d)]", "x=[a,(b,c,d)]")
	c.AssertPrint("var x=(a,b)", "var x=(a,b)")
	c.AssertPrint("var x=a,b,c", "var x=a,b,c")
	c.AssertPrint("var x=(a,b),c", "var x=(a,b),c")
	c.AssertPrint("var x=a,b=(c,d)", "var x=a,b=(c,d)")
	c.AssertPrint("var x=(a,b)(c);", "var x=(a,b)(c)")
	c.AssertPrint("var x=(a,b)`c`;", "var x=(a,b)`c`")
	c.AssertPrint("foo(a,b,c,d)", "foo(a,b,c,d)")
	c.AssertPrint("foo((a,b,c),d)", "foo((a,b,c),d)")
	c.AssertPrint("foo((a,(b,c)),d)", "foo((a,b,c),d)")
	c.AssertPrint("f(a+b,(c,d,(e,f,g)))", "f(a+b,(c,d,e,f,g))")
	c.AssertPrint("({}) , 1 , 2", "({}),1,2")
	c.AssertPrint("({}) , {} , {}", "({}),{},{}")

	c.AssertPrintSame("var a=(b=c,d)")
	c.AssertPrintSame("var a=(b[c]=d,e)")
	c.AssertPrintSame("var a=(b[c]=d,e[f]=g,h)")

	c.AssertPrint("var a = /** @type {?} */ (b=c,d)", "var a=(b=c,d)")
	c.AssertPrint("var a = /** @type {?} */ (b[c]=d,e)", "var a=(b[c]=d,e)")
	c.AssertPrint("var a = /** @type {?} */ (b[c]=d,e[f]=g,h)", "var a=(b[c]=d,e[f]=g,h)")

	// EMPTY nodes
	c.AssertPrint("if (x){}", "if(x);")
	c.AssertPrint("if(x);", "if(x);")
	c.AssertPrint("if(x)if(y);", "if(x)if(y);")
	c.AssertPrint("if(x){if(y);}", "if(x)if(y);")
	c.AssertPrint("if(x){if(y){};;;}", "if(x)if(y);")
}

// Java: CodePrinterTest.testPrintNewVoid()
func (c *CodePrinterTest) TestPrintNewVoid() {
	// Odd looking but valid. This, of course, will cause a runtime exception but
	// should not cause a parse error as "new void 0" would.
	c.AssertPrintSame("new (void 0)")
}

// Java: CodePrinterTest.testPrintComma1()
func (c *CodePrinterTest) TestPrintComma1() {
	node := rhino.IR_VarWithValue(rhino.IR_Name("a"), rhino.IR_Comma(rhino.IR_Comma(rhino.IR_Name("b"), rhino.IR_Name("c")), rhino.IR_Name("d")))
	c.AssertPrintNode("var a=(b,c,d)", node)
}

// Java: CodePrinterTest.testPrintComma2()
func (c *CodePrinterTest) TestPrintComma2() {
	node := rhino.IR_VarWithValue(rhino.IR_Name("a"), rhino.IR_Comma(rhino.IR_Name("b"), rhino.IR_Comma(rhino.IR_Name("c"), rhino.IR_Name("d"))))
	c.AssertPrintNode("var a=(b,c,d)", node)
}

// Java: CodePrinterTest.testPrettyPrintJSDoc()
func (c *CodePrinterTest) TestPrettyPrintJSDoc() {
	c.AssertPrettyPrintSame(
		("/** @type {number} */ \n" +
			"var x;\n"))
}

// Java: CodePrinterTest.testPrintCast1()
func (c *CodePrinterTest) TestPrintCast1() {
	c.AssertPrint("var x = /** @type {number} */ (0);", "var x=0")
	c.AssertPrettyPrintSame("var x = /** @type {number} */ (0);\n")
}

// Java: CodePrinterTest.testPrintCast2()
func (c *CodePrinterTest) TestPrintCast2() {
	c.AssertPrint("var x = (2+3) * 4;", "var x=(2+3)*4")
	c.AssertPrint("var x = /** @type {number} */ (2+3) * 4;", "var x=(2+3)*4")
	c.AssertPrettyPrintSame("var x = (/** @type {number} */ (2 + 3)) * 4;\n")
}

// Java: CodePrinterTest.testPrintCast3()
func (c *CodePrinterTest) TestPrintCast3() {
	c.AssertPrint("var x = (2*3) + 4;", "var x=2*3+4")
	c.AssertPrint("var x = /** @type {number} */ (2*3) + 4;", "var x=2*3+4")
	c.AssertPrettyPrintSame("var x = /** @type {number} */ (2 * 3) + 4;\n")
}

// Java: CodePrinterTest.testLetConstInIf()
func (c *CodePrinterTest) TestLetConstInIf() {
	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT_NEXT
	c.AssertPrint("if (true) { let x; };", "if(true){let x}")
	c.AssertPrint("if (true) { const x = 0; };", "if(true){const x=0}")
}

// Java: CodePrinterTest.testPrintBlockScopedFunctions()
func (c *CodePrinterTest) TestPrintBlockScopedFunctions() {
	// Safari 3 needs a "{" around a single function
	c.AssertPrint("if (true) function foo(){return}", "if(true){function foo(){return}}")
	c.AssertPrint("if(x){;;function y(){};;}", "if(x){function y(){}}")
}

// Java: CodePrinterTest.testPrintArrayPatternVar()
func (c *CodePrinterTest) TestPrintArrayPatternVar() {
	c.AssertPrintSame("var []=[]")
	c.AssertPrintSame("var [a]=[1]")
	c.AssertPrintSame("var [a,b]=[1,2]")
	c.AssertPrintSame("var [a,...b]=[1,2]")
	c.AssertPrintSame("var [,b]=[1,2]")
	c.AssertPrintSame("var [,,,,,,g]=[1,2,3,4,5,6,7]")
	c.AssertPrintSame("var [a,,c]=[1,2,3]")
	c.AssertPrintSame("var [a,,,d]=[1,2,3,4]")
	c.AssertPrintSame("var [a,,c,,e]=[1,2,3,4,5]")
}

// Java: CodePrinterTest.testPrintArrayPatternLet()
func (c *CodePrinterTest) TestPrintArrayPatternLet() {
	c.AssertPrintSame("let []=[]")
	c.AssertPrintSame("let [a]=[1]")
	c.AssertPrintSame("let [a,b]=[1,2]")
	c.AssertPrintSame("let [a,...b]=[1,2]")
	c.AssertPrintSame("let [,b]=[1,2]")
	c.AssertPrintSame("let [,,,,,,g]=[1,2,3,4,5,6,7]")
	c.AssertPrintSame("let [a,,c]=[1,2,3]")
	c.AssertPrintSame("let [a,,,d]=[1,2,3,4]")
	c.AssertPrintSame("let [a,,c,,e]=[1,2,3,4,5]")
}

// Java: CodePrinterTest.testPrintArrayPatternConst()
func (c *CodePrinterTest) TestPrintArrayPatternConst() {
	c.AssertPrintSame("const []=[]")
	c.AssertPrintSame("const [a]=[1]")
	c.AssertPrintSame("const [a,b]=[1,2]")
	c.AssertPrintSame("const [a,...b]=[1,2]")
	c.AssertPrintSame("const [,b]=[1,2]")
	c.AssertPrintSame("const [,,,,,,g]=[1,2,3,4,5,6,7]")
	c.AssertPrintSame("const [a,,c]=[1,2,3]")
	c.AssertPrintSame("const [a,,,d]=[1,2,3,4]")
	c.AssertPrintSame("const [a,,c,,e]=[1,2,3,4,5]")
}

// Java: CodePrinterTest.testPrintArrayPatternAssign()
func (c *CodePrinterTest) TestPrintArrayPatternAssign() {
	c.AssertPrintSame("[]=[]")
	c.AssertPrintSame("[a]=[1]")
	c.AssertPrintSame("[a,b]=[1,2]")
	c.AssertPrintSame("[a,...b]=[1,2]")
	c.AssertPrintSame("[,b]=[1,2]")
	c.AssertPrintSame("[,,,,,,g]=[1,2,3,4,5,6,7]")
	c.AssertPrintSame("[a,,c]=[1,2,3]")
	c.AssertPrintSame("[a,,,d]=[1,2,3,4]")
	c.AssertPrintSame("[a,,c,,e]=[1,2,3,4,5]")
}

// Java: CodePrinterTest.testPrintArrayPatternWithInitializer()
func (c *CodePrinterTest) TestPrintArrayPatternWithInitializer() {
	c.AssertPrintSame("[x=1]=[]")
	c.AssertPrintSame("[a,,c=2,,e]=[1,2,3,4,5]")
	c.AssertPrintSame("[a=1,b=2,c=3]=foo()")
	c.AssertPrintSame("[a=(1,2),b]=foo()")
	c.AssertPrintSame("[a=[b=(1,2)]=bar(),c]=foo()")
}

// Java: CodePrinterTest.testPrintNestedArrayPattern()
func (c *CodePrinterTest) TestPrintNestedArrayPattern() {
	c.AssertPrintSame("var [a,[b,c],d]=[1,[2,3],4]")
	c.AssertPrintSame("var [[[[a]]]]=[[[[1]]]]")

	c.AssertPrintSame("[a,[b,c],d]=[1,[2,3],4]")
	c.AssertPrintSame("[[[[a]]]]=[[[[1]]]]")
}

// Java: CodePrinterTest.testPrettyPrintArrayPattern()
func (c *CodePrinterTest) TestPrettyPrintArrayPattern() {
	c.AssertPrettyPrint("let [a,b,c]=foo();", "let [a, b, c] = foo();\n")
}

// Java: CodePrinterTest.testPrintObjectPatternVar()
func (c *CodePrinterTest) TestPrintObjectPatternVar() {
	c.AssertPrintSame("var {a}=foo()")
	c.AssertPrintSame("var {a,b}=foo()")
	c.AssertPrintSame("var {a:a,b:b}=foo()")
}

// Java: CodePrinterTest.testPrintObjectPatternLet()
func (c *CodePrinterTest) TestPrintObjectPatternLet() {
	c.AssertPrintSame("let {a}=foo()")
	c.AssertPrintSame("let {a,b}=foo()")
	c.AssertPrintSame("let {a:a,b:b}=foo()")
}

// Java: CodePrinterTest.testPrintObjectPatternConst()
func (c *CodePrinterTest) TestPrintObjectPatternConst() {
	c.AssertPrintSame("const {a}=foo()")
	c.AssertPrintSame("const {a,b}=foo()")
	c.AssertPrintSame("const {a:a,b:b}=foo()")
}

// Java: CodePrinterTest.testPrintObjectPatternAssign()
func (c *CodePrinterTest) TestPrintObjectPatternAssign() {
	c.AssertPrintSame("({a}=foo())")
	c.AssertPrintSame("({a,b}=foo())")
	c.AssertPrintSame("({a:a,b:b}=foo())")
}

// Java: CodePrinterTest.testPrintNestedObjectPattern()
func (c *CodePrinterTest) TestPrintNestedObjectPattern() {
	c.AssertPrintSame("({a:{b,c}}=foo())")
	c.AssertPrintSame("({a:{b:{c:{d}}}}=foo())")
}

// Java: CodePrinterTest.testPrintObjectPatternInitializer()
func (c *CodePrinterTest) TestPrintObjectPatternInitializer() {
	c.AssertPrintSame("({a=1}=foo())")
	c.AssertPrintSame("({a:{b=2}}=foo())")
	c.AssertPrintSame("({a:b=2}=foo())")
	c.AssertPrintSame("({a,b:{c=2}}=foo())")
	c.AssertPrintSame("({a:{b=2},c}=foo())")
	c.AssertPrintSame("({a=(1,2),b}=foo())")
	c.AssertPrintSame("({a:b=(1,2),c}=foo())")
}

// Java: CodePrinterTest.testPrintObjectPatternWithRest()
func (c *CodePrinterTest) TestPrintObjectPatternWithRest() {
	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT_NEXT
	c.AssertPrintSame("const {a,...rest}=foo()")
	c.AssertPrintSame("var {a,...rest}=foo()")
	c.AssertPrintSame("let {a,...rest}=foo()")
	c.AssertPrintSame("({a,...rest}=foo())")
	c.AssertPrintSame("({a=2,...rest}=foo())")
	c.AssertPrintSame("({a:b=2,...rest}=foo())")
}

// Java: CodePrinterTest.testPrettyPrintObjectPattern()
func (c *CodePrinterTest) TestPrettyPrintObjectPattern() {
	c.AssertPrettyPrint("const {a,b,c}=foo();", "const {a, b, c} = foo();\n")
}

// Java: CodePrinterTest.testPrintMixedDestructuring()
func (c *CodePrinterTest) TestPrintMixedDestructuring() {
	c.AssertPrintSame("({a:[b,c]}=foo())")
	c.AssertPrintSame("[a,{b,c}]=foo()")
}

// Java: CodePrinterTest.testPrintDestructuringInParamList()
func (c *CodePrinterTest) TestPrintDestructuringInParamList() {
	c.AssertPrintSame("function f([a]){}")
	c.AssertPrintSame("function f([a,b]){}")
	c.AssertPrintSame("function f([a,b]=c()){}")
	c.AssertPrintSame("function f([a=(1,2),b=(3,4)]=c()){}")
	c.AssertPrintSame("function f({a}){}")
	c.AssertPrintSame("function f({a,b}){}")
	c.AssertPrintSame("function f({a,b}=c()){}")
	c.AssertPrintSame("function f([a,{b,c}]){}")
	c.AssertPrintSame("function f({a,b:[c,d]}){}")
}

// Java: CodePrinterTest.testPrintDestructuringInRestParam()
func (c *CodePrinterTest) TestPrintDestructuringInRestParam() {
	c.AssertPrintSame("function f(...[a,b]){}")
	c.AssertPrintSame("function f(...{length:num_params}){}")
}

// Java: CodePrinterTest.testDestructuringForInLoops()
func (c *CodePrinterTest) TestDestructuringForInLoops() {
	c.AssertPrintSame("for({a}in b)c")
	c.AssertPrintSame("for(var {a}in b)c")
	c.AssertPrintSame("for(let {a}in b)c")
	c.AssertPrintSame("for(const {a}in b)c")

	c.AssertPrintSame("for({a:b}in c)d")
	c.AssertPrintSame("for(var {a:b}in c)d")
	c.AssertPrintSame("for(let {a:b}in c)d")
	c.AssertPrintSame("for(const {a:b}in c)d")

	c.AssertPrintSame("for([a]in b)c")
	c.AssertPrintSame("for(var [a]in b)c")
	c.AssertPrintSame("for(let [a]in b)c")
	c.AssertPrintSame("for(const [a]in b)c")
}

// Java: CodePrinterTest.testDestructuringForOfLoops1()
func (c *CodePrinterTest) TestDestructuringForOfLoops1() {
	c.AssertPrintSame("for({a}of b)c")
	c.AssertPrintSame("for(var {a}of b)c")
	c.AssertPrintSame("for(let {a}of b)c")
	c.AssertPrintSame("for(const {a}of b)c")

	c.AssertPrintSame("for({a:b}of c)d")
	c.AssertPrintSame("for(var {a:b}of c)d")
	c.AssertPrintSame("for(let {a:b}of c)d")
	c.AssertPrintSame("for(const {a:b}of c)d")

	c.AssertPrintSame("for([a]of b)c")
	c.AssertPrintSame("for(var [a]of b)c")
	c.AssertPrintSame("for(let [a]of b)c")
	c.AssertPrintSame("for(const [a]of b)c")
}

// Java: CodePrinterTest.testDestructuringForOfLoops2()
func (c *CodePrinterTest) TestDestructuringForOfLoops2() {
	// The destructuring 'var' statement is a child of the for-of loop, but
	// not the first child.
	c.AssertPrintSame("for(a of b)var {x}=y")
}

// Java: CodePrinterTest.testBreakTrustedStrings()
func (c *CodePrinterTest) TestBreakTrustedStrings() {
	// Break scripts
	c.AssertPrint("'<script>'", "\"<script>\"")
	c.AssertPrint("'</script>'", "\"\\x3c/script>\"")
	c.AssertPrint("\"</script> </SCRIPT>\"", "\"\\x3c/script> \\x3c/SCRIPT>\"")

	c.AssertPrint("'-->'", "\"--\\x3e\"")
	c.AssertPrint("']]>'", "\"]]\\x3e\"")
	c.AssertPrint("' --></script>'", "\" --\\x3e\\x3c/script>\"")

	c.AssertPrint("/--> <\\/script>/g", "/--\\x3e <\\/script>/g")

	// Break HTML start comments. Certain versions of WebKit
	// begin an HTML comment when they see this.
	c.AssertPrint("'<!-- I am a string -->'", "\"\\x3c!-- I am a string --\\x3e\"")

	c.AssertPrint("'<=&>'", "\"<=&>\"")
}

// Java: CodePrinterTest.testBreakUntrustedStrings()
func (c *CodePrinterTest) TestBreakUntrustedStrings() {
	c.trustedStrings = false

	// Break scripts
	c.AssertPrint("'<script>'", "\"\\x3cscript\\x3e\"")
	c.AssertPrint("'</script>'", "\"\\x3c/script\\x3e\"")
	c.AssertPrint("\"</script> </SCRIPT>\"", "\"\\x3c/script\\x3e \\x3c/SCRIPT\\x3e\"")

	c.AssertPrint("'-->'", "\"--\\x3e\"")
	c.AssertPrint("']]>'", "\"]]\\x3e\"")
	c.AssertPrint("' --></script>'", "\" --\\x3e\\x3c/script\\x3e\"")

	c.AssertPrint("/--> <\\/script>/g", "/--\\x3e <\\/script>/g")

	// Break HTML start comments. Certain versions of WebKit
	// begin an HTML comment when they see this.
	c.AssertPrint("'<!-- I am a string -->'", "\"\\x3c!-- I am a string --\\x3e\"")

	c.AssertPrint("'<=&>'", "\"\\x3c\\x3d\\x26\\x3e\"")
	c.AssertPrint("/(?=x)/", "/(?=x)/")
}

// Java: CodePrinterTest.testBreakTemplateLiterals()
func (c *CodePrinterTest) TestBreakTemplateLiterals() {
	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT_NEXT

	// Break scripts in template literals
	c.AssertPrint("`</script>`", "`\\x3c/script>`")
	c.AssertPrint("`</SCRIPT>`", "`\\x3c/SCRIPT>`")
	c.AssertPrint("`</style>`", "`\\x3c/style>`")
	c.AssertPrint("`</STYLE>`", "`\\x3c/STYLE>`")

	// Break comments in template literals
	c.AssertPrint("`-->`", "`--\\x3e`")
	c.AssertPrint("`]]>`", "`]]\\x3e`")
	c.AssertPrint("` --></script>`", "` --\\x3e\\x3c/script>`")
	c.AssertPrint("`<!-- I am a string -->`", "`\\x3c!-- I am a string --\\x3e`")
}

// Java: CodePrinterTest.testBreakUntrustedTemplateLiterals()
func (c *CodePrinterTest) TestBreakUntrustedTemplateLiterals() {
	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT_NEXT
	c.trustedStrings = false

	// Break scripts in template literals
	c.AssertPrint("`</script>`", "`\\x3c/script\\x3e`")
	c.AssertPrint("`<!-- I am a string -->`", "`\\x3c!-- I am a string --\\x3e`")
	c.AssertPrint("`<=&>`", "`\\x3c\\x3d\\x26\\x3e`")
}

// Java: CodePrinterTest.testHtmlComments()
func (c *CodePrinterTest) TestHtmlComments() {
	c.AssertPrint("3< !(--x)", "3< !--x")
	c.AssertPrint("while (x-- > 0) {}", "while(x-- >0);")
}

// Java: CodePrinterTest.testPrintArray()
func (c *CodePrinterTest) TestPrintArray() {
	c.AssertPrint("[void 0, void 0]", "[void 0,void 0]")
	c.AssertPrint("[undefined, undefined]", "[undefined,undefined]")
	c.AssertPrint("[ , , , undefined]", "[,,,undefined]")
	c.AssertPrint("[ , , , 0]", "[,,,0]")
}

// Java: CodePrinterTest.testHook()
func (c *CodePrinterTest) TestHook() {
	c.AssertPrint("a ? b = 1 : c = 2", "a?b=1:c=2")
	c.AssertPrint("x = a ? b = 1 : c = 2", "x=a?b=1:c=2")
	c.AssertPrint("(x = a) ? b = 1 : c = 2", "(x=a)?b=1:c=2")

	c.AssertPrint("x, a ? b = 1 : c = 2", "x,a?b=1:c=2")
	c.AssertPrint("x, (a ? b = 1 : c = 2)", "x,a?b=1:c=2")
	c.AssertPrint("(x, a) ? b = 1 : c = 2", "(x,a)?b=1:c=2")

	c.AssertPrint("a ? (x, b) : c = 2", "a?(x,b):c=2")
	c.AssertPrint("a ? b = 1 : (x,c)", "a?b=1:(x,c)")

	c.AssertPrint("a ? b = 1 : c = 2 + x", "a?b=1:c=2+x")
	c.AssertPrint("(a ? b = 1 : c = 2) + x", "(a?b=1:c=2)+x")
	c.AssertPrint("a ? b = 1 : (c = 2) + x", "a?b=1:(c=2)+x")

	c.AssertPrint("a ? (b?1:2) : 3", "a?b?1:2:3")
}

// Java: CodePrinterTest.testForIn()
func (c *CodePrinterTest) TestForIn() {
	c.AssertPrintSame("for(a in b)c")
	c.AssertPrintSame("for(var a in b)c")
	c.AssertPrintSame("for(var a in b=c)d")
	c.AssertPrintSame("for(var a in b,c)d")
}

// Java: CodePrinterTest.testPrintInOperatorInForLoop()
func (c *CodePrinterTest) TestPrintInOperatorInForLoop() {
	// Check for in expression in for's init expression.
	// Check alone, with + (higher precedence), with ?: (lower precedence),
	// and with conditional.
	c.AssertPrint(
		"var a={}; for (var i = (\"length\" in a); i;) {}",
		"var a={};for(var i=(\"length\"in a);i;);")
	c.AssertPrint(
		"var a={}; for (var i = (\"length\" in a) ? 0 : 1; i;) {}",
		"var a={};for(var i=(\"length\"in a)?0:1;i;);")
	c.AssertPrint(
		"var a={}; for (var i = (\"length\" in a) + 1; i;) {}",
		"var a={};for(var i=(\"length\"in a)+1;i;);")
	c.AssertPrint(
		"var a={};for (var i = (\"length\" in a|| \"size\" in a);;);",
		"var a={};for(var i=(\"length\"in a)||(\"size\"in a);;);")
	c.AssertPrint(
		"var a={};for (var i = (a || a) || (\"size\" in a);;);",
		"var a={};for(var i=a||a||(\"size\"in a);;);")

	// Test works with unary operators and calls.
	c.AssertPrint(
		"var a={}; for (var i = -(\"length\" in a); i;) {}",
		"var a={};for(var i=-(\"length\"in a);i;);")
	c.AssertPrint(
		"var a={};function b_(p){ return p;};for(var i=1,j=b_(\"length\" in a);;) {}",
		"var a={};function b_(p){return p}for(var i=1,j=b_(\"length\"in a);;);")

	// Test we correctly handle an in operator in the test clause.
	c.AssertPrint("var a={}; for (;(\"length\" in a);) {}", "var a={};for(;\"length\"in a;);")

	// Test we correctly handle an in operator inside a comma.
	c.AssertPrintSame("for(x,(y in z);;)foo()")
	c.AssertPrintSame("for(var x,w=(y in z);;)foo()")

	// And in operator inside a hook.
	c.AssertPrintSame("for(a=c?0:(0 in d);;)foo()")

	// And inside an arrow function body
	c.AssertPrint(
		"var a={}; for(var i = () => (0 in a); i;) {}", "var a={};for(var i=()=>(0 in a);i;);")
	c.AssertPrint(
		"var a={}; for(var i = () => ({} in a); i;) {}", "var a={};for(var i=()=>({}in a);i;);")

	// And inside a destructuring declaration
	c.AssertPrint(
		"var a={}; for(var {noop} = (\"prop\" in a); noop;) {}",
		"var a={};for(var {noop}=(\"prop\"in a);noop;);")
}

// Java: CodePrinterTest.testForOf()
func (c *CodePrinterTest) TestForOf() {
	c.AssertPrintSame("for(a of b)c")
	c.AssertPrintSame("for(var a of b)c")
	c.AssertPrintSame("for(var a of b=c)d")
	c.AssertPrintSame("for(var a of(b,c))d")
}

// In pretty-print mode, make sure there is a space before and after the 'of' in a for/of loop.
// Java: CodePrinterTest.testForOfPretty()
func (c *CodePrinterTest) TestForOfPretty() {
	c.AssertPrettyPrintSame("for ([x, y] of b) {\n  c;\n}\n")
	c.AssertPrettyPrintSame("for (x of [[1, 2]]) {\n  c;\n}\n")
	c.AssertPrettyPrintSame("for ([x, y] of [[1, 2]]) {\n  c;\n}\n")
}

// Java: CodePrinterTest.testForAwaitOf()
func (c *CodePrinterTest) TestForAwaitOf() {
	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT_NEXT

	c.AssertPrintSame("async()=>{for await(a of b)c}")
	c.AssertPrintSame("async()=>{for await(var a of b)c}")
	c.AssertPrintSame("async()=>{for await(var a of b=c)d}")
	c.AssertPrintSame("async()=>{for await(var a of(b,c))d}")
}

// In pretty-print mode, make sure there is a space before and after the 'of' in a for/of loop.
// Java: CodePrinterTest.testForAwaitOfPretty()
func (c *CodePrinterTest) TestForAwaitOfPretty() {
	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT_NEXT

	c.AssertPrettyPrintSame("async() => {\n  for await ([x, y] of b) {\n    c;\n  }\n};\n")
	c.AssertPrettyPrintSame("async() => {\n  for await (x of [[1, 2]]) {\n    c;\n  }\n};\n")
	c.AssertPrettyPrintSame("async() => {\n  for await ([x, y] of [[1, 2]]) {\n    c;\n  }\n};\n")
}

// Java: CodePrinterTest.testLetFor()
func (c *CodePrinterTest) TestLetFor() {
	c.AssertPrintSame("for(let a=0;a<5;a++)b")
	c.AssertPrintSame("for(let a in b)c")
	c.AssertPrintSame("for(let a of b)c")
	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT_NEXT
	c.AssertPrintSame("async()=>{for await(let a of b)c}")
}

// Java: CodePrinterTest.testConstFor()
func (c *CodePrinterTest) TestConstFor() {
	c.AssertPrintSame("for(const a=5;b<a;b++)c")
	c.AssertPrintSame("for(const a in b)c")
	c.AssertPrintSame("for(const a of b)c")

	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT_NEXT
	c.AssertPrintSame("async()=>{for await(const a of b)c}")
}

// Java: CodePrinterTest.testLiteralProperty()
func (c *CodePrinterTest) TestLiteralProperty() {
	c.AssertPrint("(64).toString()", "(64).toString()")
}

// Make sure that the code generator doesn't associate an
// else clause with the wrong if clause.
// Java: CodePrinterTest.testAmbiguousElseClauses()
func (c *CodePrinterTest) TestAmbiguousElseClauses() {
	c.AssertPrintNode(
		"if(x)if(y);else;",
		rhino.NewNodeWith2Children(
			rhino.Token_IF,
			rhino.Node_NewStringWithToken(rhino.Token_NAME, "x"),
			rhino.NewNodeWithChild(
				rhino.Token_BLOCK,
				rhino.NewNodeWith3Children(
					rhino.Token_IF,
					rhino.Node_NewStringWithToken(rhino.Token_NAME, "y"),
					rhino.NewNode(rhino.Token_BLOCK),

					// ELSE clause for the inner if
					rhino.NewNode(rhino.Token_BLOCK)))))

	c.AssertPrintNode(
		"if(x){if(y);}else;",
		rhino.NewNodeWith3Children(
			rhino.Token_IF,
			rhino.Node_NewStringWithToken(rhino.Token_NAME, "x"),
			rhino.NewNodeWithChild(
				rhino.Token_BLOCK,
				rhino.NewNodeWith2Children(rhino.Token_IF, rhino.Node_NewStringWithToken(rhino.Token_NAME, "y"), rhino.NewNode(rhino.Token_BLOCK))),

			// ELSE clause for the outer if
			rhino.NewNode(rhino.Token_BLOCK)))

	c.AssertPrintNode(
		"if(x)if(y);else{if(z);}else;",
		rhino.NewNodeWith3Children(
			rhino.Token_IF,
			rhino.Node_NewStringWithToken(rhino.Token_NAME, "x"),
			rhino.NewNodeWithChild(
				rhino.Token_BLOCK,
				rhino.NewNodeWith3Children(
					rhino.Token_IF,
					rhino.Node_NewStringWithToken(rhino.Token_NAME, "y"),
					rhino.NewNode(rhino.Token_BLOCK),
					rhino.NewNodeWithChild(
						rhino.Token_BLOCK,
						rhino.NewNodeWith2Children(
							rhino.Token_IF, rhino.Node_NewStringWithToken(rhino.Token_NAME, "z"), rhino.NewNode(rhino.Token_BLOCK))))),

			// ELSE clause for the outermost if
			rhino.NewNode(rhino.Token_BLOCK)))
}

// Java: CodePrinterTest.testLineBreak()
func (c *CodePrinterTest) TestLineBreak() {
	// line break after function if in a statement context
	c.AssertLineBreak(
		("function a() {}\n" +
			"function b() {}\n"),
		("function a(){}\n" +
			"function b(){}\n"))

	// line break after ; after a function
	c.AssertLineBreak(
		("var a = {};\n" +
			"a.foo = function () {}\n" +
			"function b() {}\n"),
		("var a={};a.foo=function(){};\n" +
			"function b(){}\n"))

	// break after comma after a function
	c.AssertLineBreak(
		("var a = {\n" +
			"  b: function() {},\n" +
			"  c: function() {}\n" +
			"};\n" +
			"alert(a);\n"),
		("var a={b:function(){},\n" +
			"c:function(){}};\n" +
			"alert(a)"))
}

// Java: CodePrinterTest.assertLineBreak(String, String)
func (c *CodePrinterTest) AssertLineBreak(js string, expected string) {
	AssertThat(
		c.ParsePrint(
			js,
			c.NewCompilerOptions(
				codePrinterTestBase_CompilerOptionBuilderFunc(func(options *CompilerOptions) {
					options.SetPrettyPrint(false)
					options.SetLineBreak(true)
					options.SetLineLengthThreshold(
						CompilerOptions_DEFAULT_LINE_LENGTH_THRESHOLD)
				})))).
		IsEqualTo(expected)
}

// Java: CodePrinterTest.testPrettyPrinter()
func (c *CodePrinterTest) TestPrettyPrinter() {
	// Ensure that the pretty printer inserts line breaks at appropriate
	// places.
	c.AssertPrettyPrint("(function(){})();", "(function() {\n})();\n")
	c.AssertPrettyPrint("var a = (function() {});alert(a);", "var a = function() {\n};\nalert(a);\n")

	// Check we correctly handle putting brackets around all if clauses so
	// we can put breakpoints inside statements.
	c.AssertPrettyPrint(
		"if (1) {}",
		("if (1) {\n" +
			"}\n"))
	c.AssertPrettyPrint(
		"if (1) {alert(\"\");}",
		("if (1) {\n" +
			"  alert(\"\");\n" +
			"}\n"))
	c.AssertPrettyPrint(
		"if (1)alert(\"\");",
		("if (1) {\n" +
			"  alert(\"\");\n" +
			"}\n"))
	c.AssertPrettyPrint(
		"if (1) {alert();alert();}",
		("if (1) {\n" +
			"  alert();\n" +
			"  alert();\n" +
			"}\n"))

	// Don't add blocks if they weren't there already.
	c.AssertPrettyPrint("label: alert();", "label: alert();\n")

	// But if statements and loops get blocks automagically.
	c.AssertPrettyPrint(
		"if (1) alert();",
		("if (1) {\n" +
			"  alert();\n" +
			"}\n"))
	c.AssertPrettyPrint(
		"for (;;) alert();",
		("for (;;) {\n" +
			"  alert();\n" +
			"}\n"))

	c.AssertPrettyPrint(
		"while (1) alert();",
		("while (1) {\n" +
			"  alert();\n" +
			"}\n"))

	// Do we put else clauses in blocks?
	c.AssertPrettyPrint(
		"if (1) {} else {alert(a);}",
		("if (1) {\n" +
			"} else {\n" +
			"  alert(a);\n" +
			"}\n"))

	// Do we add blocks to else clauses?
	c.AssertPrettyPrint(
		"if (1) alert(a); else alert(b);",
		("if (1) {\n" +
			"  alert(a);\n" +
			"} else {\n" +
			"  alert(b);\n" +
			"}\n"))

	// Do we put for bodies in blocks?
	c.AssertPrettyPrint(
		"for(;;) { alert();}",
		("for (;;) {\n" +
			"  alert();\n" +
			"}\n"))
	c.AssertPrettyPrint(
		"for(;;) {}",
		("for (;;) {\n" +
			"}\n"))
	c.AssertPrettyPrint(
		"for(;;) { alert(); alert(); }",
		("for (;;) {\n" +
			"  alert();\n" +
			"  alert();\n" +
			"}\n"))
	c.AssertPrettyPrint(
		"for(var x=0;x<10;x++) { alert(); alert(); }",
		("for (var x = 0; x < 10; x++) {\n" +
			"  alert();\n" +
			"  alert();\n" +
			"}\n"))

	// How about do loops?
	c.AssertPrettyPrint(
		"do { alert(); } while(true);",
		("do {\n" +
			"  alert();\n" +
			"} while (true);\n"))

	// label?
	c.AssertPrettyPrint(
		"myLabel: { alert();}",
		("myLabel: {\n" +
			"  alert();\n" +
			"}\n"))
	c.AssertPrettyPrint("myLabel: {}", "myLabel: {\n}\n")
	c.AssertPrettyPrint("myLabel: ;", "myLabel: ;\n")

	// Don't move the label on a loop, because then break {label} and
	// continue {label} won't work.
	c.AssertPrettyPrint(
		"myLabel: for(;;) continue myLabel;",
		("myLabel: for (;;) {\n" +
			"  continue myLabel;\n" +
			"}\n"))

	c.AssertPrettyPrint("var a;", "var a;\n")
	c.AssertPrettyPrint("i--", "i--;\n")
	c.AssertPrettyPrint("i++", "i++;\n")

	// There must be a space before and after binary operators.
	c.AssertPrettyPrint("var foo = 3+5;", "var foo = 3 + 5;\n")

	// There should be spaces between the ternary operator
	c.AssertPrettyPrint("var foo = bar ? 3 : null;", "var foo = bar ? 3 : null;\n")

	// Ensure that string literals after return and throw have a space.
	c.AssertPrettyPrint(
		"function foo() { return \"foo\"; }", "function foo() {\n  return \"foo\";\n}\n")
	c.AssertPrettyPrint("throw \"foo\";", "throw \"foo\";\n")

	// Test that loops properly have spaces inserted.
	c.AssertPrettyPrint("do{ alert(); } while(true);", "do {\n  alert();\n} while (true);\n")
	c.AssertPrettyPrint("while(true) { alert(); }", "while (true) {\n  alert();\n}\n")
}

// Java: CodePrinterTest.testPrettyPrinter2()
func (c *CodePrinterTest) TestPrettyPrinter2() {
	c.AssertPrettyPrint(
		"if(true) f();",
		("if (true) {\n" +
			"  f();\n" +
			"}\n"))

	c.AssertPrettyPrint(
		"if (true) { f() } else { g() }",
		("if (true) {\n" +
			"  f();\n" +
			"} else {\n" +
			"  g();\n" +
			"}\n"))

	c.AssertPrettyPrint(
		"if(true) f(); for(;;) g();",
		("if (true) {\n" +
			"  f();\n" +
			"}\n" +
			"for (;;) {\n" +
			"  g();\n" +
			"}\n"))
}

// Java: CodePrinterTest.testPrettyPrinter3()
func (c *CodePrinterTest) TestPrettyPrinter3() {
	c.AssertPrettyPrint(
		"try {} catch(e) {}if (1) {alert();alert();}",
		("try {\n" +
			"} catch (e) {\n" +
			"}\n" +
			"if (1) {\n" +
			"  alert();\n" +
			"  alert();\n" +
			"}\n"))

	c.AssertPrettyPrint(
		"try {} finally {}if (1) {alert();alert();}",
		("try {\n" +
			"} finally {\n" +
			"}\n" +
			"if (1) {\n" +
			"  alert();\n" +
			"  alert();\n" +
			"}\n"))

	c.AssertPrettyPrint(
		"try {} catch(e) {} finally {} if (1) {alert();alert();}",
		("try {\n" +
			"} catch (e) {\n" +
			"} finally {\n" +
			"}\n" +
			"if (1) {\n" +
			"  alert();\n" +
			"  alert();\n" +
			"}\n"))
}

// Java: CodePrinterTest.testPrettyPrinter4()
func (c *CodePrinterTest) TestPrettyPrinter4() {
	c.AssertPrettyPrint(
		"function f() {}if (1) {alert();}",
		("function f() {\n" +
			"}\n" +
			"if (1) {\n" +
			"  alert();\n" +
			"}\n"))

	c.AssertPrettyPrint(
		"var f = function() {};if (1) {alert();}",
		("var f = function() {\n" +
			"};\n" +
			"if (1) {\n" +
			"  alert();\n" +
			"}\n"))

	c.AssertPrettyPrint(
		"(function() {})();if (1) {alert();}",
		("(function() {\n" +
			"})();\n" +
			"if (1) {\n" +
			"  alert();\n" +
			"}\n"))

	c.AssertPrettyPrint(
		"(function() {alert();alert();})();if (1) {alert();}",
		("(function() {\n" +
			"  alert();\n" +
			"  alert();\n" +
			"})();\n" +
			"if (1) {\n" +
			"  alert();\n" +
			"}\n"))
}

// Java: CodePrinterTest.testPrettyPrinter_arrow()
func (c *CodePrinterTest) TestPrettyPrinter_arrow() {
	c.AssertPrettyPrint("(a)=>123;", "a => 123;\n")
}

// Java: CodePrinterTest.testPrettyPrinter_defaultValue()
func (c *CodePrinterTest) TestPrettyPrinter_defaultValue() {
	c.AssertPrettyPrint("(a=1)=>123;", "(a = 1) => 123;\n")
	c.AssertPrettyPrint("[a=(1,2)]=[];", "[a = (1, 2)] = [];\n")
}

// For https://github.com/google/closure-compiler/issues/782
// Java: CodePrinterTest.testPrettyPrinter_spaceBeforeSingleQuote()
func (c *CodePrinterTest) TestPrettyPrinter_spaceBeforeSingleQuote() {
	c.AssertPrettyPrintWithOptionBuilder(
		"var f = function() { return 'hello'; };",
		("var f = function() {\n" +
			"  return 'hello';\n" +
			"};\n"),
		codePrinterTestBase_CompilerOptionBuilderFunc(func(options *CompilerOptions) {
			options.SetPreferSingleQuotes(true)
		}))
}

// For https://github.com/google/closure-compiler/issues/782
// Java: CodePrinterTest.testPrettyPrinter_spaceBeforeUnaryOperators()
func (c *CodePrinterTest) TestPrettyPrinter_spaceBeforeUnaryOperators() {
	c.AssertPrettyPrint(
		"var f = function() { return !b; };",
		("var f = function() {\n" +
			"  return !b;\n" +
			"};\n"))
	c.AssertPrettyPrint(
		"var f = function*(){yield -b}",
		("var f = function*() {\n" +
			"  yield -b;\n" +
			"};\n"))
	c.AssertPrettyPrint(
		"var f = function() { return +b; };",
		("var f = function() {\n" +
			"  return +b;\n" +
			"};\n"))
	c.AssertPrettyPrint(
		"var f = function() { throw ~b; };",
		("var f = function() {\n" +
			"  throw ~b;\n" +
			"};\n"))
	c.AssertPrettyPrint(
		"var f = function() { return ++b; };",
		("var f = function() {\n" +
			"  return ++b;\n" +
			"};\n"))
	c.AssertPrettyPrint(
		"var f = function() { return --b; };",
		("var f = function() {\n" +
			"  return --b;\n" +
			"};\n"))
}

// Java: CodePrinterTest.testPrettyPrinter_varLetConst()
func (c *CodePrinterTest) TestPrettyPrinter_varLetConst() {
	c.AssertPrettyPrint("var x=0;", "var x = 0;\n")
	c.AssertPrettyPrint("const x=0;", "const x = 0;\n")
	c.AssertPrettyPrint("let x=0;", "let x = 0;\n")
}

// Java: CodePrinterTest.testPrettyPrinter_number()
func (c *CodePrinterTest) TestPrettyPrinter_number() {
	c.AssertPrettyPrintSame("var x = 10;\n")
	c.AssertPrettyPrintSame("var x = 1.;\n")
	c.AssertPrettyPrint("var x = 0xFE;", "var x = 254;\n")
	c.AssertPrettyPrintSame("var x = 1" + strings.Repeat("0", 100) + ";\n") // a googol
	c.AssertPrettyPrintSame("f(10000);\n")
	c.AssertPrettyPrint("var x = -10000;\n", "var x = -10000;\n")
	c.AssertPrettyPrint("var x = y - -10000;\n", "var x = y - -10000;\n")
	c.AssertPrettyPrint("f(-10000);\n", "f(-10000);\n")
	c.AssertPrettyPrintSame("x < 2592000;\n")
	c.AssertPrettyPrintSame("x < 1000.000;\n")
	c.AssertPrettyPrintSame("x < 1000.912;\n")
	c.AssertPrettyPrintSame("var x = 1E20;\n")
	c.AssertPrettyPrintSame("var x = 1E1;\n")
	c.AssertPrettyPrintSame("var x = void 0;\n")
	c.AssertPrettyPrintSame("foo(-0);\n")
	c.AssertPrettyPrint("var x = 4-1000;", "var x = 4 - 1000;\n")
}

// Java: CodePrinterTest.testTypeAnnotations()
func (c *CodePrinterTest) TestTypeAnnotations() {
	c.AssertTypeAnnotations(
		"/** @constructor */ function Foo(){}", "/**\n * @constructor\n */\nfunction Foo() {\n}\n")
}

// Java: CodePrinterTest.testNonNullTypes()
func (c *CodePrinterTest) TestNonNullTypes() {
	c.AssertTypeAnnotations(
		("/** @constructor */\n" +
			"function Foo() {}\n" +
			"/** @return {!Foo} */\n" +
			"Foo.prototype.f = function() { return new Foo; };\n"),
		("/**\n" +
			" * @constructor\n" +
			" */\n" +
			"function Foo() {\n" +
			"}\n" +
			"/**\n" +
			" * @return {!Foo}\n" +
			" */\n" +
			"Foo.prototype.f = function() {\n" +
			"  return new Foo();\n" +
			"};\n"))
}

// Java: CodePrinterTest.testTypeAnnotationsTypeDef()
func (c *CodePrinterTest) TestTypeAnnotationsTypeDef() {
	// TODO(johnlenz): It would be nice if there were some way to preserve
	// typedefs but currently they are resolved into the basic types in the
	// type registry.
	c.AssertTypeAnnotations(
		("/** @const */ var goog = {};\n" +
			"/** @const */ goog.java = {};\n" +
			"/** @typedef {Array<number>} */ goog.java.Long;\n" +
			"/** @param {!goog.java.Long} a*/\n" +
			"function f(a){};\n"),
		("/** @const */ var goog = {};\n" +
			"/** @const */ goog.java = {};\n" +
			"goog.java.Long;\n" +
			"/**\n" +
			" * @param {!Array<number>} a\n" +
			" * @return {undefined}\n" +
			" */\n" +
			"function f(a) {\n" +
			"}\n"))
}

// Java: CodePrinterTest.testTypeAnnotationsAssign()
func (c *CodePrinterTest) TestTypeAnnotationsAssign() {
	c.AssertTypeAnnotations(
		"/** @constructor */ var Foo = function(){}",
		("/**\n" +
			" * @constructor\n" +
			" */\n" +
			"var Foo = function() {\n" +
			"};\n"))
}

// Java: CodePrinterTest.testTypeAnnotationsNamespace_varWithoutJSDoc()
func (c *CodePrinterTest) TestTypeAnnotationsNamespace_varWithoutJSDoc() {
	c.AssertTypeAnnotations(
		("var a = {};\n" +
			"/** @constructor */ a.Foo = function(){}\n"),
		("var a = {};\n" +
			"/**\n" +
			" * @constructor\n" +
			" */\n" +
			"a.Foo = function() {\n" +
			"};\n"))
}

// Java: CodePrinterTest.testTypeAnnotationsNamespace_varWithConstJSDoc()
func (c *CodePrinterTest) TestTypeAnnotationsNamespace_varWithConstJSDoc() {
	c.AssertTypeAnnotations(
		("/** @const */\n" +
			"var a = {};\n" +
			"/** @constructor */ a.Foo = function(){}\n"),
		("/** @const */ var a = {};\n" +
			"/**\n" +
			" * @constructor\n" +
			" */\n" +
			"a.Foo = function() {\n" +
			"};\n"))
}

// Java: CodePrinterTest.testTypeAnnotationsNamespace_constDeclarationWithoutJSDoc()
func (c *CodePrinterTest) TestTypeAnnotationsNamespace_constDeclarationWithoutJSDoc() {
	c.AssertTypeAnnotations(
		("const a = {};\n" +
			"/** @constructor */ a.Foo = function(){}\n"),
		("const a = {};\n" +
			"/**\n" +
			" * @constructor\n" +
			" */\n" +
			"a.Foo = function() {\n" +
			"};\n"))
}

// Java: CodePrinterTest.testTypeAnnotationsNamespace_constDeclarationWithJSDoc()
func (c *CodePrinterTest) TestTypeAnnotationsNamespace_constDeclarationWithJSDoc() {
	c.AssertTypeAnnotations(
		("/** @export */\n" +
			"const a = {};\n" +
			"/** @constructor */ a.Foo = function(){}\n"),
		("/** @export */ const a = {};\n" +
			"/**\n" +
			" * @constructor\n" +
			" */\n" +
			"a.Foo = function() {\n" +
			"};\n"))
}

// Java: CodePrinterTest.testTypeAnnotationsNamespace_qnameWithConstJSDoc()
func (c *CodePrinterTest) TestTypeAnnotationsNamespace_qnameWithConstJSDoc() {
	c.AssertTypeAnnotations(
		("/** @const */\n" +
			"var a = {};\n" +
			"/** @const */\n" +
			"a.b = {};\n" +
			"/** @constructor */ a.b.Foo = function(){}\n"),
		("/** @const */ var a = {};\n" +
			"/** @const */ a.b = {};\n" +
			"/**\n" +
			" * @constructor\n" +
			" */\n" +
			"a.b.Foo = function() {\n" +
			"};\n"))
}

// Java: CodePrinterTest.testTypeAnnotationsMemberSubclass()
func (c *CodePrinterTest) TestTypeAnnotationsMemberSubclass() {
	c.AssertTypeAnnotations(
		("/** @const */ var a = {};\n" +
			"/** @constructor */ a.Foo = function(){};\n" +
			"/** @constructor\n" +
			" @extends {a.Foo} */ a.Bar = function(){}\n"),
		("/** @const */ var a = {};\n" +
			"/**\n" +
			" * @constructor\n" +
			" */\n" +
			"a.Foo = function() {\n" +
			"};\n" +
			"/**\n" +
			" * @extends {a.Foo}\n" +
			" * @constructor\n" +
			" */\n" +
			"a.Bar = function() {\n" +
			"};\n"))
}

// Java: CodePrinterTest.testTypeAnnotationsInterface()
func (c *CodePrinterTest) TestTypeAnnotationsInterface() {
	c.AssertTypeAnnotations(
		("/** @const */ var a = {};\n" +
			"/** @interface */ a.Foo = function(){};\n" +
			"/** @interface\n" +
			" @extends {a.Foo} */ a.Bar = function(){}\n"),
		("/** @const */ var a = {};\n" +
			"/**\n" +
			" * @interface\n" +
			" */\n" +
			"a.Foo = function() {\n" +
			"};\n" +
			"/**\n" +
			" * @extends {a.Foo}\n" +
			" * @interface\n" +
			" */\n" +
			"a.Bar = function() {\n" +
			"};\n"))
}

// Java: CodePrinterTest.testTypeAnnotationsMultipleInterface()
func (c *CodePrinterTest) TestTypeAnnotationsMultipleInterface() {
	c.AssertTypeAnnotations(
		("/** @const */ var a = {};\n" +
			"/** @interface */ a.Foo1 = function(){};\n" +
			"/** @interface */ a.Foo2 = function(){};\n" +
			"/** @interface\n" +
			" @extends {a.Foo1}\n" +
			" @extends {a.Foo2} */\n" +
			"a.Bar = function(){}\n"),
		("/** @const */ var a = {};\n" +
			"/**\n" +
			" * @interface\n" +
			" */\n" +
			"a.Foo1 = function() {\n" +
			"};\n" +
			"/**\n" +
			" * @interface\n" +
			" */\n" +
			"a.Foo2 = function() {\n" +
			"};\n" +
			"/**\n" +
			" * @extends {a.Foo1}\n" +
			" * @extends {a.Foo2}\n" +
			" * @interface\n" +
			" */\n" +
			"a.Bar = function() {\n" +
			"};\n"))
}

// Java: CodePrinterTest.testTypeAnnotationsMember()
func (c *CodePrinterTest) TestTypeAnnotationsMember() {
	c.AssertTypeAnnotations(
		("var a = {};\n" +
			"/** @constructor */ a.Foo = function(){}\n" +
			"/** @param {string} foo\n" +
			"  * @return {number} */\n" +
			"a.Foo.prototype.foo = function(foo) { return 3; };\n" +
			"/** @type {!Array|undefined} */\n" +
			"a.Foo.prototype.bar = [];\n"),
		("var a = {};\n" +
			"/**\n" +
			" * @constructor\n" +
			" */\n" +
			"a.Foo = function() {\n" +
			"};\n" +
			"/**\n" +
			" * @param {string} foo\n" +
			" * @return {number}\n" +
			" */\n" +
			"a.Foo.prototype.foo = function(foo) {\n" +
			"  return 3;\n" +
			"};\n" +
			"/** @type {!Array<?>} */\n" +
			"a.Foo.prototype.bar = [];\n"))
}

// Java: CodePrinterTest.testTypeAnnotationsMemberStub()
func (c *CodePrinterTest) TestTypeAnnotationsMemberStub() {
	// TODO(blickly): Investigate why the method's type isn't preserved.
	c.AssertTypeAnnotations(
		("/** @interface */ function I(){};\n" +
			"/** @return {undefined} @param {number} x */ I.prototype.method;\n"),
		"/**\n * @interface\n */\nfunction I() {\n}\nI.prototype.method;\n")
}

// Java: CodePrinterTest.testTypeAnnotationsImplements()
func (c *CodePrinterTest) TestTypeAnnotationsImplements() {
	c.AssertTypeAnnotations(
		("/** @const */ var a = {};\n" +
			"/** @constructor */ a.Foo = function(){};\n" +
			"/** @interface */ a.I = function(){};\n" +
			"/** @record */ a.I2 = function(){};\n" +
			"/** @record @extends {a.I2} */ a.I3 = function(){};\n" +
			"/** @constructor\n" +
			" @extends {a.Foo}\n" +
			" * @implements {a.I}\n" +
			" @implements {a.I2}\n" +
			" */ a.Bar = function(){}\n"),
		("/** @const */ var a = {};\n" +
			"/**\n" +
			" * @constructor\n" +
			" */\n" +
			"a.Foo = function() {\n" +
			"};\n" +
			"/**\n" +
			" * @interface\n" +
			" */\n" +
			"a.I = function() {\n" +
			"};\n" +
			"/**\n" +
			" * @record\n" +
			" */\n" +
			"a.I2 = function() {\n" +
			"};\n" +
			"/**\n" +
			" * @extends {a.I2}\n" +
			" * @record\n" +
			" */\n" +
			"a.I3 = function() {\n" +
			"};\n" +
			"/**\n" +
			" * @extends {a.Foo}\n" +
			" * @implements {a.I}\n" +
			" * @implements {a.I2}\n" +
			" * @constructor\n" +
			" */\n" +
			"a.Bar = function() {\n" +
			"};\n"))
}

// Java: CodePrinterTest.testTypeAnnotationClassImplements()
func (c *CodePrinterTest) TestTypeAnnotationClassImplements() {
	c.AssertTypeAnnotations(
		("/** @interface */ class Foo {}\n" +
			"/** @implements {Foo} */ class Bar {}\n"),
		("/**\n" +
			" * @interface\n" +
			" */\n" +
			"class Foo {\n" +
			"}\n" +
			"/**\n" +
			" * @implements {Foo}\n" +
			" */\n" +
			"class Bar {\n" +
			"}\n"))
}

// Java: CodePrinterTest.testTypeAnnotationClassMember()
func (c *CodePrinterTest) TestTypeAnnotationClassMember() {
	c.AssertTypeAnnotations(
		("class Foo {\n" +
			"  /** @return {number} */ method(/** string */ arg) {}\n" +
			"}\n"),
		("class Foo {\n" +
			"  /**\n" +
			"   * @param {string} arg\n" +
			"   * @return {number}\n" +
			"   */\n" +
			"  method(arg) {\n" +
			"  }\n" +
			"}\n"))
}

// Java: CodePrinterTest.testTypeAnnotationClassConstructor()
func (c *CodePrinterTest) TestTypeAnnotationClassConstructor() {
	c.AssertTypeAnnotations(
		("/**\n" +
			" * @template T\n" +
			" */\n" +
			"class Foo {\n" +
			"  /** @param {T} arg */\n" +
			"  constructor(arg) {}\n" +
			"}\n"),
		("/**\n" +
			" * @template T\n" +
			" */\n" +
			"class Foo {\n" +
			"  /**\n" +
			"   * @param {T} arg\n" +
			"   */\n" +
			"  constructor(arg) {\n" +
			"  }\n" +
			"}\n"))
}

// Java: CodePrinterTest.testRestParameter()
func (c *CodePrinterTest) TestRestParameter() {
	c.AssertTypeAnnotations(
		("/** @param {...string} args */\n" +
			"function f(...args) {}\n"),
		("/**\n" +
			" * @param {...string} args\n" +
			" * @return {undefined}\n" +
			" */\n" +
			"function f(...args) {\n" +
			"}\n"))
}

// Java: CodePrinterTest.testDefaultParameter()
func (c *CodePrinterTest) TestDefaultParameter() {
	c.AssertTypeAnnotations(
		("/** @param {string=} msg */\n" +
			"function f(msg = 'hi') {}\n"),
		("/**\n" +
			" * @param {string=} msg\n" +
			" * @return {undefined}\n" +
			" */\n" +
			"function f(msg = \"hi\") {\n" +
			"}\n"))
}

// Java: CodePrinterTest.testObjectDestructuringParameter()
func (c *CodePrinterTest) TestObjectDestructuringParameter() {
	c.AssertTypeAnnotations(
		("/** @param {{a: number, b: number}} ignoredName */\n" +
			"function f({a, b}) {}\n"),
		// The `ignoredName` JSDoc name is ignored
		// The whitespace in the `{a, b}` output must match the input
		("/**\n" +
			" * @param {{a: number, b: number}} p0\n" +
			" * @return {undefined}\n" +
			" */\n" +
			"function f({a, b}) {\n" +
			"}\n"))
}

// Java: CodePrinterTest.testObjectDestructuringParameterWithDefault()
func (c *CodePrinterTest) TestObjectDestructuringParameterWithDefault() {
	c.AssertTypeAnnotations(
		("/** @param {{a: number, b: number}=} ignoredName */\n" +
			"function f({a, b} = {a: 1, b: 2}) {}\n"),
		// The `ignoredName` JSDoc name is ignored
		// The whitespace in the `{a, b}` output must match the input
		("/**\n" +
			" * @param {{a: number, b: number}=} p0\n" +
			" * @return {undefined}\n" +
			" */\n" +
			"function f({a, b} = {a:1, b:2}) {\n" +
			"}\n"))
}

// Java: CodePrinterTest.testArrayDestructuringParameter()
func (c *CodePrinterTest) TestArrayDestructuringParameter() {
	c.AssertTypeAnnotations(
		("/** @param {!Iterable<number>} ignoredName */\n" +
			"function f([a, b]) {}\n"),
		// The `ignoredName` JSDoc name is ignored
		// The whitespace in the `[a, b]` output must match the input
		("/**\n" +
			" * @param {!Iterable<number,?,?>} p0\n" +
			" * @return {undefined}\n" +
			" */\n" +
			"function f([a, b]) {\n" +
			"}\n"))
}

// Java: CodePrinterTest.testArrayDestructuringParameterWithDefault()
func (c *CodePrinterTest) TestArrayDestructuringParameterWithDefault() {
	c.AssertTypeAnnotations(
		("/** @param {!Iterable<number>=} ignoredName */\n" +
			"function f([a, b] = [1, 2]) {}\n"),
		// The `ignoredName` JSDoc name is ignored
		// The whitespace in the `[a, b]` output must match the input
		("/**\n" +
			" * @param {!Iterable<number,?,?>=} p0\n" +
			" * @return {undefined}\n" +
			" */\n" +
			"function f([a, b] = [1, 2]) {\n" +
			"}\n"))
}

// Java: CodePrinterTest.testU2UFunctionTypeAnnotation1()
func (c *CodePrinterTest) TestU2UFunctionTypeAnnotation1() {
	c.AssertTypeAnnotations(
		"/** @type {!Function} */ var x = function() {}",
		"/** @type {!Function} */\nvar x = function() {\n};\n")
}

// Java: CodePrinterTest.testU2UFunctionTypeAnnotation2()
func (c *CodePrinterTest) TestU2UFunctionTypeAnnotation2() {
	// TODO(johnlenz): we currently report the type of the RHS which is not
	// correct, we should export the type of the LHS.
	c.AssertTypeAnnotations(
		"/** @type {Function} */ var x = function() {}",
		"/** @type {!Function} */\nvar x = function() {\n};\n")
}

// Java: CodePrinterTest.testEmitUnknownParamTypesAsAllType()
func (c *CodePrinterTest) TestEmitUnknownParamTypesAsAllType() {
	// x is unused, so NTI infers that x can be omitted.
	c.AssertTypeAnnotations(
		"var a = function(x) {}",
		("/**\n" +
			" * @param {?} x\n" +
			" * @return {undefined}\n" +
			" */\n" +
			"var a = function(x) {\n" +
			"};\n"))
}

// Java: CodePrinterTest.testOptionalTypesAnnotation()
func (c *CodePrinterTest) TestOptionalTypesAnnotation() {
	c.AssertTypeAnnotations(
		"/** @param {string=} x */ var a = function(x) {}",
		("/**\n" +
			" * @param {string=} x\n" +
			" * @return {undefined}\n" +
			" */\n" +
			"var a = function(x) {\n" +
			"};\n"))
}

// Java: CodePrinterTest.testOptionalTypesAnnotation2()
func (c *CodePrinterTest) TestOptionalTypesAnnotation2() {
	c.AssertTypeAnnotations(
		"/** @param {undefined=} x */ var a = function(x) {}",
		("/**\n" +
			" * @param {undefined=} x\n" +
			" * @return {undefined}\n" +
			" */\n" +
			"var a = function(x) {\n" +
			"};\n"))
}

// Java: CodePrinterTest.testVariableArgumentsTypesAnnotation()
func (c *CodePrinterTest) TestVariableArgumentsTypesAnnotation() {
	c.AssertTypeAnnotations(
		"/** @param {...string} x */ var a = function(x) {}",
		("/**\n" +
			" * @param {...string} x\n" +
			" * @return {undefined}\n" +
			" */\n" +
			"var a = function(x) {\n" +
			"};\n"))
}

// Java: CodePrinterTest.testTempConstructor()
func (c *CodePrinterTest) TestTempConstructor() {
	c.AssertTypeAnnotations(
		("var x = function() {\n" +
			"  /** @constructor */ function t1() {}\n" +
			"  /** @constructor */ function t2() {}\n" +
			"  t1.prototype = t2.prototype\n" +
			"}\n"),
		("/**\n" +
			" * @return {undefined}\n" +
			" */\n" +
			"var x = function() {\n" +
			"  /**\n" +
			"   * @constructor\n" +
			"   */\n" +
			"  function t1() {\n" +
			"  }\n" +
			"  /**\n" +
			"   * @constructor\n" +
			"   */\n" +
			"  function t2() {\n" +
			"  }\n" +
			"  t1.prototype = t2.prototype;\n" +
			"};\n"))
}

// Java: CodePrinterTest.testEnumAnnotation1()
func (c *CodePrinterTest) TestEnumAnnotation1() {
	c.AssertTypeAnnotations(
		"/** @enum {string} */ const Enum = {FOO: 'x', BAR: 'y'};",
		"/** @enum {string} */\nconst Enum = {FOO:\"x\", BAR:\"y\"};\n")
}

// Java: CodePrinterTest.testEnumAnnotation2()
func (c *CodePrinterTest) TestEnumAnnotation2() {
	c.AssertTypeAnnotations(
		("/** @const */ var goog = goog || {};\n" +
			"/** @enum {string} */ goog.Enum = {FOO: 'x', BAR: 'y'};\n" +
			"/** @const */ goog.Enum2 = goog.x ? {} : goog.Enum;\n"),
		("/** @const */ var goog = goog || {};\n" +
			"/** @enum {string} */\n" +
			"goog.Enum = {FOO:\"x\", BAR:\"y\"};\n" +
			"/** @type {(!Object|{})} */\n" +
			"goog.Enum2 = goog.x ? {} : goog.Enum;\n"))
}

// Java: CodePrinterTest.testEnumAnnotation3()
func (c *CodePrinterTest) TestEnumAnnotation3() {
	c.AssertTypeAnnotations(
		"/** @enum {!Object} */ var Enum = {FOO: {}};",
		"/** @enum {!Object} */\nvar Enum = {FOO:{}};\n")
}

// Java: CodePrinterTest.testEnumAnnotation4()
func (c *CodePrinterTest) TestEnumAnnotation4() {
	c.AssertTypeAnnotations(
		("/** @enum {number} */ var E = {A:1, B:2};\n" +
			"function f(/** !E */ x) { return x; }\n"),
		("/** @enum {number} */\n" +
			"var E = {A:1, B:2};\n" +
			"/**\n" +
			" * @param {number} x\n" +
			" * @return {?}\n" +
			" */\n" +
			"function f(x) {\n" +
			"  return x;\n" +
			"}\n"))
}

// Java: CodePrinterTest.testClosureLibraryTypeAnnotationExamples()
func (c *CodePrinterTest) TestClosureLibraryTypeAnnotationExamples() {
	c.AssertTypeAnnotations(
		("/** @const */ var goog = goog || {};\n" +
			"/** @param {Object} obj */goog.removeUid = function(obj) {};\n" +
			"/** @param {Object} obj The object to remove the field from. */\n" +
			"goog.removeHashCode = goog.removeUid;\n"),
		("/** @const */ var goog = goog || {};\n" +
			"/**\n" +
			" * @param {(Object|null)} obj\n" +
			" * @return {undefined}\n" +
			" */\n" +
			"goog.removeUid = function(obj) {\n" +
			"};\n" +
			"/**\n" +
			" * @param {(Object|null)} p0\n" +
			" * @return {undefined}\n" +
			" */\n" +
			"goog.removeHashCode = goog.removeUid;\n"))
}

// Java: CodePrinterTest.testFunctionTypeAnnotation()
func (c *CodePrinterTest) TestFunctionTypeAnnotation() {
	c.AssertTypeAnnotations(
		"/**\n * @param {{foo:number}} arg\n */\nfunction f(arg) {}",
		"/**\n * @param {{foo: number}} arg\n * @return {undefined}\n */\nfunction f(arg) {\n}\n")
	c.AssertTypeAnnotations(
		"/**\n * @param {number} arg\n */\nfunction f(arg) {}",
		"/**\n * @param {number} arg\n * @return {undefined}\n */\nfunction f(arg) {\n}\n")
	c.AssertTypeAnnotations(
		"/**\n * @param {!Array<string>} arg\n */\nfunction f(arg) {}",
		"/**\n * @param {!Array<string>} arg\n * @return {undefined}\n */\nfunction f(arg) {\n}\n")
}

// Java: CodePrinterTest.testFunctionWithThisTypeAnnotation()
func (c *CodePrinterTest) TestFunctionWithThisTypeAnnotation() {
	c.AssertTypeAnnotations(
		"/**\n * @this {{foo:number}}\n */\nfunction foo() {}",
		"/**\n * @return {undefined}\n * @this {{foo: number}}\n */\nfunction foo() {\n}\n")
	c.AssertTypeAnnotations(
		"/**\n * @this {!Array<string>}\n */\nfunction foo() {}",
		"/**\n * @return {undefined}\n * @this {!Array<string>}\n */\nfunction foo() {\n}\n")
}

// Java: CodePrinterTest.testReturnWithTypeAnnotation()
func (c *CodePrinterTest) TestReturnWithTypeAnnotation() {
	c.preserveTypeAnnotations = true
	c.AssertPrettyPrint(
		"function f() { return (/** @return {number} */ function() { return 42; }); }",
		("function f() {\n" +
			"  return (/**\n" +
			"   * @return {number}\n" +
			"   */\n" +
			"  function() {\n" +
			"    return 42;\n" +
			"  });\n" +
			"}\n"))
}

// Java: CodePrinterTest.testDeprecatedAnnotationIncludesNewline()
func (c *CodePrinterTest) TestDeprecatedAnnotationIncludesNewline() {
	js := ("/**\n" +
		" * @type {number}\n" +
		" * @deprecated See {@link replacementClass} for more details.\n" +
		" */\n" +
		"var x;\n")

	c.AssertPrettyPrint(js, js)
}

// Java: CodePrinterTest.testNonJSDocCommentsPrinted_nonTrailing_blockComment()
func (c *CodePrinterTest) TestNonJSDocCommentsPrinted_nonTrailing_blockComment() {
	c.preserveNonJSDocComments = true
	c.AssertPrettyPrint(
		"/* testComment */ function Foo(){}", "/* testComment */ function Foo() {\n}\n")
}

// Java: CodePrinterTest.testNonJSDocCommentsPrinted_endOfFile_lineComment()
func (c *CodePrinterTest) TestNonJSDocCommentsPrinted_endOfFile_lineComment() {
	c.preserveNonJSDocComments = true
	c.AssertPrettyPrint(
		("function f1() {}\n" +
			"if (true) {\n" +
			"// first\n" +
			"f1();\n" +
			"}\n" +
			"// second\n"),
		("function f1() {\n" +
			"}\n" +
			"if (true) {\n" +
			"  // first\n" +
			"  f1();\n" +
			"}\n" +
			" // second\n"))
}

// Java: CodePrinterTest.testNonJSDocCommentsPrinted_endOfBlockComment()
func (c *CodePrinterTest) TestNonJSDocCommentsPrinted_endOfBlockComment() {
	c.preserveNonJSDocComments = true
	c.AssertPrettyPrint(
		("function f1() {}\n" +
			"if (true) {\n" +
			"// first\n" +
			"f1();\n" +
			"/* second */\n" +
			"}\n"),
		("function f1() {\n" +
			"}\n" +
			"if (true) {\n" +
			"  // first\n" +
			"  f1(); \n" +
			"  /* second */\n" +
			"}\n"))
}

// Java: CodePrinterTest.testNonJSDocCommentsPrinted_endOfBlock_manyMixedComments()
func (c *CodePrinterTest) TestNonJSDocCommentsPrinted_endOfBlock_manyMixedComments() {
	c.preserveNonJSDocComments = true
	c.AssertPrettyPrint(
		("function f1() {}\n" +
			"if (true) {\n" +
			"// first\n" +
			"f1();\n" +
			"// second\n" +
			"/* third */\n" +
			"// fourth\n" +
			"}\n"),
		("function f1() {\n" +
			"}\n" +
			"if (true) {\n" +
			"  // first\n" +
			"  f1(); \n" +
			"  // second\n" +
			"  /* third */\n" +
			"  // fourth\n" +
			"}\n"))
}

// Java: CodePrinterTest.testNonJSDocCommentsPrinted_lastTrailing()
func (c *CodePrinterTest) TestNonJSDocCommentsPrinted_lastTrailing() {
	c.preserveNonJSDocComments = true
	c.AssertPrettyPrint(
		("function f1() {}\n" +
			"if (true) {\n" +
			"// first\n" +
			"f1(); // second\n" +
			"}\n"),
		("function f1() {\n" +
			"}\n" +
			"if (true) {\n" +
			"  // first\n" +
			"  f1(); // second\n" +
			"}\n"))
}

// Java: CodePrinterTest.testNonJSDocCommentsPrinted_nonTrailing_lineComment()
func (c *CodePrinterTest) TestNonJSDocCommentsPrinted_nonTrailing_lineComment() {
	c.preserveNonJSDocComments = true
	c.AssertPrettyPrint("// testComment\nfunction Foo(){}", "// testComment\nfunction Foo() {\n}\n")
}

// Java: CodePrinterTest.testNonJSDocCommentsPrinted_betweenCode_sameLine()
func (c *CodePrinterTest) TestNonJSDocCommentsPrinted_betweenCode_sameLine() {
	c.preserveNonJSDocComments = true

	c.AssertPrettyPrint(
		"function /* testComment */ Foo(){}", "function/* testComment */ Foo() {\n}\n")
}

// Java: CodePrinterTest.testNonJSDocCommentsPrinted_betweenCode_differentLines()
func (c *CodePrinterTest) TestNonJSDocCommentsPrinted_betweenCode_differentLines() {
	c.preserveNonJSDocComments = true
	c.AssertPrettyPrint(
		"function /* testComment */\nFoo(){}", "function/* testComment */\nFoo() {\n}\n")
}

// Java: CodePrinterTest.testNonJSDocCommentsPrinted_nonTrailing_inlineComments()
func (c *CodePrinterTest) TestNonJSDocCommentsPrinted_nonTrailing_inlineComments() {
	c.preserveNonJSDocComments = true
	// tests inline comments in parameter lists are parsed and printed
	c.AssertPrettyPrint(
		"function Foo(/*first*/ x, /* second*/ y) {}",
		"function Foo(/*first*/ x, /* second*/ y) {\n}\n")
}

// Args on new line are condensed onto the same line by prettyPrint
// Java: CodePrinterTest.testArgs_noComments_newLines()
func (c *CodePrinterTest) TestArgs_noComments_newLines() {
	c.AssertPrettyPrint(" var rpcid = new RpcId(a,\n b, \nc);", "var rpcid = new RpcId(a, b, c);\n")
}

// Comments are printed when args on new line are condensed onto the same line by prettyPrint
// Java: CodePrinterTest.testNonJSDocCommentsPrinted_nonTrailing_inlineComments_newLines()
func (c *CodePrinterTest) TestNonJSDocCommentsPrinted_nonTrailing_inlineComments_newLines() {
	c.preserveNonJSDocComments = true
	c.AssertPrettyPrint(
		" var rpcid = new RpcId(a,\n /* comment1 */ b, \n/* comment1 */ c);",
		"var rpcid = new RpcId(a, /* comment1 */ b, /* comment1 */ c);\n")
}

// Java: CodePrinterTest.testNonJSDocCommentsPrinted_trailing_and_nonTrailing_inlineComments()
func (c *CodePrinterTest) TestNonJSDocCommentsPrinted_trailing_and_nonTrailing_inlineComments() {
	c.preserveNonJSDocComments = true
	// tests inline trailing comments in parameter lists are parsed and printed
	c.AssertPrettyPrint(
		"function Foo(x //first\n, /* second*/ y) {}",
		"function Foo(x //first\n, /* second*/ y) {\n}\n")
}

// Java: CodePrinterTest.testNonJSDocCommentsPrinted_trailing_inlineComments_paramList()
func (c *CodePrinterTest) TestNonJSDocCommentsPrinted_trailing_inlineComments_paramList() {
	c.preserveNonJSDocComments = true
	c.AssertPrettyPrint("function Foo(x) {}", "function Foo(x) {\n}\n")
	c.AssertPrettyPrint("function Foo(x /*first*/) {}", "function Foo(x /*first*/) {\n}\n")
	c.AssertPrettyPrint("function Foo(x //first\n) {}", "function Foo(x //first\n) {\n}\n")
}

// Java: CodePrinterTest.testClassExtendsLeftHandSideExpression()
func (c *CodePrinterTest) TestClassExtendsLeftHandSideExpression() {
	c.AssertPrettyPrint(
		"class A {} class B extends (0, A) {}", "class A {\n}\nclass B extends(0, A) {\n}\n")
}

// Same as above, but tests argList instead of Param list
// Java: CodePrinterTest.testNonJSDocCommentsPrinted_trailing_inlineComments_callArgList()
func (c *CodePrinterTest) TestNonJSDocCommentsPrinted_trailing_inlineComments_callArgList() {
	c.preserveNonJSDocComments = true
	c.AssertPrettyPrint("foo(x);", "foo(x);\n")
	c.AssertPrettyPrint("foo(x /*first*/);", "foo(x /*first*/);\n")
	c.AssertPrettyPrint("foo(x //first\n);", "foo(x //first\n);\n")
}

// Java: CodePrinterTest.assertPrettyPrintSame(String)
func (c *CodePrinterTest) AssertPrettyPrintSame(js string) {
	c.AssertPrettyPrint(js, js)
}

// Java: CodePrinterTest.assertPrettyPrint(String, String)
func (c *CodePrinterTest) AssertPrettyPrint(js string, expected string) {
	c.AssertPrettyPrintWithOptionBuilder(
		js,
		expected,
		codePrinterTestBase_CompilerOptionBuilderFunc(func(options *CompilerOptions) {
			/* no-op */
		}))
}

// Java: CodePrinterTest.assertPrettyPrint(String, String, CompilerOptionBuilder)
func (c *CodePrinterTest) AssertPrettyPrintWithOptionBuilder(js string, expected string, optionBuilder CodePrinterTestBase_CompilerOptionBuilder) {
	AssertThat(
		c.ParsePrint(
			js,
			c.NewCompilerOptions(
				codePrinterTestBase_CompilerOptionBuilderFunc(func(options *CompilerOptions) {
					options.SetPrettyPrint(true)
					options.SetPreserveTypeAnnotations(true)
					options.SetLineBreak(false)
					options.SetLineLengthThreshold(
						CompilerOptions_DEFAULT_LINE_LENGTH_THRESHOLD)
					optionBuilder.SetOptions(options)
				})))).
		IsEqualTo(expected)
}

// Java: CodePrinterTest.assertTypeAnnotations(String, String)
func (c *CodePrinterTest) AssertTypeAnnotations(js string, expected string) {
	actual := NewCodePrinter_Builder(c.ParseTypeChecked(js /* typeChecked= */, true)).
		SetCompilerOptions(
			c.NewCompilerOptions(
				codePrinterTestBase_CompilerOptionBuilderFunc(func(options *CompilerOptions) {
					options.SetPrettyPrint(true)
					options.SetLineBreak(false)
					options.SetLineLengthThreshold(
						CompilerOptions_DEFAULT_LINE_LENGTH_THRESHOLD)
				}))).
		SetOutputTypes(true).
		SetTypeRegistry(c.lastCompiler.GetTypeRegistry()).
		Build()

	AssertThat(actual).IsEqualTo(expected)
}

// Java: CodePrinterTest.testSubtraction()
func (c *CodePrinterTest) TestSubtraction() {
	compiler := NewCompiler()
	n := compiler.ParseTestCode("x - -4")
	AssertThat(compiler.GetErrorCount()).IsEqualTo(0)

	AssertThat(c.PrintNode(n)).IsEqualTo("x- -4")
}

// Java: CodePrinterTest.testFunctionWithCall()
func (c *CodePrinterTest) TestFunctionWithCall() {
	c.AssertPrint(
		("var user = new function() {\n" +
			"alert(\"foo\")}\n"),
		"var user=new function(){alert(\"foo\")}")
	c.AssertPrint(
		("var user = new function() {\n" +
			"this.name = \"foo\";\n" +
			"this.local = function(){alert(this.name)};}\n"),
		"var user=new function(){this.name=\"foo\";this.local=function(){alert(this.name)}}")
}

// Java: CodePrinterTest.testLineLength()
func (c *CodePrinterTest) TestLineLength() {
	// list
	c.AssertLineLength("var aba,bcb,cdc", "var aba,bcb,\ncdc")

	// operators, and two breaks
	c.AssertLineLength(
		"\"foo\"+\"bar,baz,bomb\"+\"whee\"+\";long-string\"\n+\"aaa\"",
		("\"foo\"+\"bar,baz,bomb\"+\n" +
			"\"whee\"+\";long-string\"+\n" +
			"\"aaa\""))

	// assignment
	c.AssertLineLength(
		"var abazaba=1234",
		("var abazaba=\n" +
			"1234"))

	// statements
	c.AssertLineLength(
		"var abab=1;var bab=2",
		("var abab=1;\n" +
			"var bab=2"))

	// don't break regexes
	c.AssertLineLength(
		"var a=/some[reg](ex),with.*we?rd|chars/i;var b=a",
		("var a=/some[reg](ex),with.*we?rd|chars/i;\n" +
			"var b=a"))

	// don't break strings
	c.AssertLineLength(
		"var a=\"foo,{bar};baz\";var b=a",
		("var a=\"foo,{bar};baz\";\n" +
			"var b=a"))

	// don't break before post inc/dec
	c.AssertLineLength(
		"var a=\"a\";a++;var b=\"bbb\";",
		("var a=\"a\";a++;\n" +
			"var b=\"bbb\""))
}

// Java: CodePrinterTest.assertLineLength(String, String)
func (c *CodePrinterTest) AssertLineLength(js string, expected string) {
	AssertThat(
		c.ParsePrint(
			js,
			c.NewCompilerOptions(
				codePrinterTestBase_CompilerOptionBuilderFunc(func(options *CompilerOptions) {
					options.SetPrettyPrint(false)
					options.SetLineBreak(true)
					options.SetLineLengthThreshold(10)
				})))).
		IsEqualTo(expected)
}

// Java: CodePrinterTest.testParsePrintParse()
func (c *CodePrinterTest) TestParsePrintParse() {
	c.TestReparse("3;")
	c.TestReparse("var a = b;")
	c.TestReparse("var x, y, z;")
	c.TestReparse("try { foo() } catch(e) { bar() }")
	c.TestReparse("try { foo() } catch(e) { bar() } finally { stuff() }")
	c.TestReparse("try { foo() } finally { stuff() }")
	c.TestReparse("throw 'me'")
	c.TestReparse("function foo(a) { return a + 4; }")
	c.TestReparse("function foo() { return; }")
	c.TestReparse("var a = function(a, b) { foo(); return a + b; }")
	c.TestReparse("b = [3, 4, 'paul', \"Buchhe it\",,5];")
	c.TestReparse("v = (5, 6, 7, 8)")
	c.TestReparse("d = 34.0; x = 0; y = .3; z = -22")
	c.TestReparse("d = -x; t = !x + ~y;")
	c.TestReparse(
		("'hi'; /* just a test */ stuff(a,b) \n" +
			" foo(); // and another \n" +
			" bar();\n"))
	c.TestReparse("a = b++ + ++c; a = b++-++c; a = - --b; a = - ++b;")
	c.TestReparse("a++; b= a++; b = ++a; b = a--; b = --a; a+=2; b-=5")
	c.TestReparse("a = (2 + 3) * 4;")
	c.TestReparse("a = 1 + (2 + 3) + 4;")
	c.TestReparse("x = a ? b : c; x = a ? (b,3,5) : (foo(),bar());")
	c.TestReparse(
		("a = b | c || d ^ e\n" +
			"&& f & !g != h << i <= j < k >>> l > m * n % !o\n"))
	c.TestReparse(
		("a == b; a != b; a === b; a == b == a;\n" +
			" (a == b) == a; a == (b == a);\n"))
	c.TestReparse("if (a > b) a = b; if (b < 3) a = 3; else c = 4;")
	c.TestReparse("if (a == b) { a++; } if (a == 0) { a++; } else { a --; }")
	c.TestReparse("for (var i in a) b += i;")
	c.TestReparse(
		("for (var i = 0; i < 10; i++){ b /= 2;\n" +
			" if (b == 2)break;else continue;}\n"))
	c.TestReparse("for (x = 0; x < 10; x++) a /= 2;")
	c.TestReparse("for (;;) a++;")
	c.TestReparse("while(true) { blah(); }while(true) blah();")
	c.TestReparse("do stuff(); while(a>b);")
	c.TestReparse("[0, null, , true, false, this];")
	c.TestReparse("s.replace(/absc/, 'X').replace(/ab/gi, 'Y');")
	c.TestReparse("new Foo; new Bar(a, b,c);")
	c.TestReparse("with(foo()) { x = z; y = t; } with(bar()) a = z;")
	c.TestReparse("delete foo['bar']; delete foo;")
	c.TestReparse("var x = { 'a':'paul', 1:'3', 2:(3,4) };")
	c.TestReparse(
		("switch(a) { case 2: case 3: stuff(); break;\n" +
			"case 4: morestuff(); break; default: done();}\n"))
	c.TestReparse("x = foo['bar'] + foo['my stuff'] + foo[bar] + f.stuff;")
	c.TestReparse("a.v = b.v; x['foo'] = y['zoo'];")
	c.TestReparse("'test' in x; 3 in x; a in x;")
	c.TestReparse("'foo\"bar' + \"foo'c\" + 'stuff\\n and \\\\more'")
	c.TestReparse("x.__proto__;")
}

// Java: CodePrinterTest.testReparse(String)
func (c *CodePrinterTest) TestReparse(code string) {
	compiler := NewCompiler()
	parseOnce := c.Parse(code)
	parseTwice := c.Parse(NewCodePrinter_Builder(parseOnce).Build())
	AssertNode(parseTwice).UsingSerializer(compiler.ToSourceOfNode).IsEqualTo(parseOnce)
}

// Java: CodePrinterTest.testDoLoopIECompatibility()
func (c *CodePrinterTest) TestDoLoopIECompatibility() {
	// Do loops within IFs cause syntax errors in IE6 and IE7.
	c.AssertPrint(
		"function f(){if(e1){do foo();while(e2)}else foo()}",
		"function f(){if(e1){do foo();while(e2)}else foo()}")

	c.AssertPrint(
		"function f(){if(e1)do foo();while(e2)else foo()}",
		"function f(){if(e1){do foo();while(e2)}else foo()}")

	c.AssertPrint("if(x){do{foo()}while(y)}else bar()", "if(x){do foo();while(y)}else bar()")

	c.AssertPrint("if(x)do{foo()}while(y);else bar()", "if(x){do foo();while(y)}else bar()")

	c.AssertPrint("if(x){do{foo()}while(y)}", "if(x){do foo();while(y)}")

	c.AssertPrint("if(x)do{foo()}while(y);", "if(x){do foo();while(y)}")

	c.AssertPrint("if(x)A:do{foo()}while(y);", "if(x){A:do foo();while(y)}")

	c.AssertPrint(
		"var i = 0;a: do{b: do{i++;break b;} while(0);} while(0);",
		"var i=0;a:do{b:do{i++;break b}while(0)}while(0)")
}

// Java: CodePrinterTest.testFunctionSafariCompatibility()
func (c *CodePrinterTest) TestFunctionSafariCompatibility() {
	// Functions within IFs cause syntax errors on Safari.
	c.AssertPrint(
		"function f(){if(e1){function goo(){return true}}else foo()}",
		"function f(){if(e1){function goo(){return true}}else foo()}")

	c.AssertPrint(
		"function f(){if(e1)function goo(){return true}else foo()}",
		"function f(){if(e1){function goo(){return true}}else foo()}")

	c.AssertPrint("if(e1){function goo(){return true}}", "if(e1){function goo(){return true}}")

	c.AssertPrint("if(e1)function goo(){return true}", "if(e1){function goo(){return true}}")
}

// Java: CodePrinterTest.testExponents()
func (c *CodePrinterTest) TestExponents() {
	c.AssertPrintNumberOfInt("1", 1)
	c.AssertPrintNumberOfInt("10", 10)
	c.AssertPrintNumberOfInt("100", 100)
	c.AssertPrintNumberOfInt("1E3", 1000)
	c.AssertPrintNumberOfInt("1E4", 10000)
	c.AssertPrintNumberOfInt("1E5", 100000)
	c.AssertPrintNumber("1E18", 1000000000000000000.0)
	c.AssertPrintNumber("1E5", 100000.0)
	c.AssertPrintNumber("100000.1", 100000.1)

	c.AssertPrintNumber("1E-6", 0.000001)
	c.AssertPrintNumber("0x38d7ea4c68001", 0x38d7ea4c68001p0)
	c.AssertPrintNumber("0x7fffffffffffffff", 0x7fffffffffffffffp0)

	c.AssertPrintNumber(".01", 0.01)
	c.AssertPrintNumber("1.01", 1.01)
}

// Java: CodePrinterTest.testFormatNumber()
func (c *CodePrinterTest) TestFormatNumber() {
	// Schubfach boundary values
	AssertThat(CodeConsumer_FormatNumber(2.0e-3)).IsEqualTo(".002")
	AssertThat(CodeConsumer_FormatNumber(1.0e23)).IsEqualTo("1E23")
	AssertThat(CodeConsumer_FormatNumber(jutil.DoubleParseDouble("9.007199254740993e15"))).
		IsEqualTo("9007199254740992")

	// Subnormals and MAX_VALUE
	AssertThat(CodeConsumer_FormatNumber(math.SmallestNonzeroFloat64)).IsEqualTo("4.9E-324")
	AssertThat(CodeConsumer_FormatNumber(math.MaxFloat64)).IsEqualTo("1.7976931348623157E308")

	// Powers of 10
	AssertThat(CodeConsumer_FormatNumber(100)).IsEqualTo("100")
	AssertThat(CodeConsumer_FormatNumber(1000)).IsEqualTo("1E3")
	AssertThat(CodeConsumer_FormatNumber(1e12)).IsEqualTo("1E12")
	AssertThat(CodeConsumer_FormatNumber(1e13)).IsEqualTo("1E13")

	// Hex-vs-decimal crossover integers around 1000000000000L
	AssertThat(CodeConsumer_FormatNumber(float64(int64(1000000000000)))).IsEqualTo("1E12")
	AssertThat(CodeConsumer_FormatNumber(float64(int64(1000000000001)))).IsEqualTo("0xe8d4a51001")

	// Leading-zero fractions
	AssertThat(CodeConsumer_FormatNumber(0.5)).IsEqualTo(".5")
	AssertThat(CodeConsumer_FormatNumber(0.002)).IsEqualTo(".002")
	AssertThat(CodeConsumer_FormatNumber(0.0001)).IsEqualTo("1E-4")

	// .0E stripping
	AssertThat(CodeConsumer_FormatNumber(1.0e21)).IsEqualTo("1E21")
}

// Java: CodePrinterTest.testBiggerThanMaxLongNumericLiterals()
func (c *CodePrinterTest) TestBiggerThanMaxLongNumericLiterals() {
	// Since ECMAScript implements IEEE 754 "round to nearest, ties to even",
	// any literal in the range [0x7ffffffffffffe00,0x8000000000000400] will
	// round to the same value, namely 2^63. The fact that we print this as
	// 2^63-1 doesn't matter, since it must be rounded back to 2^63 at runtime.
	// See:
	//   http://www.ecma-international.org/ecma-262/5.1/#sec-8.5
	c.AssertPrint("9223372036854775808", "0x7fffffffffffffff")
	c.AssertPrint("0x8000000000000000", "0x7fffffffffffffff")
	c.AssertPrint(
		"0b1000000000000000000000000000000000000000000000000000000000000000", "0x7fffffffffffffff")
	c.AssertPrint("0o1000000000000000000000", "0x7fffffffffffffff")
}

// Make sure to test as both a String and a Node, because
// negative numbers do not parse consistently from strings.
// Java: CodePrinterTest.assertPrintNumber(String, double)
func (c *CodePrinterTest) AssertPrintNumber(expected string, number float64) {
	c.AssertPrint(jutil.DoubleToString(number), expected)
	c.AssertPrintNode(expected, rhino.Node_NewNumber(number))
}

// Java: CodePrinterTest.assertPrintNumber(String, int)
func (c *CodePrinterTest) AssertPrintNumberOfInt(expected string, number int) {
	c.AssertPrint(strconv.Itoa(number), expected)
	c.AssertPrintNode(expected, rhino.Node_NewNumber(float64(number)))
}

// Java: CodePrinterTest.testDirectEval()
func (c *CodePrinterTest) TestDirectEval() {
	c.AssertPrint("eval('1');", "eval(\"1\")")
}

// Java: CodePrinterTest.testIndirectEval()
func (c *CodePrinterTest) TestIndirectEval() {
	n := c.Parse("eval('1');")
	c.AssertPrintNode("eval(\"1\")", n)
	n.GetFirstFirstChild().GetFirstChild().PutBooleanProp(rhino.Node_DIRECT_EVAL, false)
	c.AssertPrintNode("(0,eval)(\"1\")", n)
}

// Java: CodePrinterTest.freeCallTaggedTemplate()
func (c *CodePrinterTest) TestFreeCallTaggedTemplate() {
	n := c.Parse("a.b`xyz`")
	call := n.GetFirstFirstChild()
	AssertThat(call.IsTaggedTemplateLit()).IsTrue()
	call.PutBooleanProp(rhino.Node_FREE_CALL, true)
	c.AssertPrintNode("(0,a.b)`xyz`", n)
}

// Java: CodePrinterTest.freeCallOptChain()
func (c *CodePrinterTest) TestFreeCallOptChain() {
	n := c.Parse("(a?.b)()")
	call := n.GetFirstFirstChild()
	AssertThat(call.IsCall()).IsTrue()
	call.PutBooleanProp(rhino.Node_FREE_CALL, true)
	c.AssertPrintNode("(0,a?.b)()", n)
}

// Java: CodePrinterTest.freeCallOptChainOptChainCall()
func (c *CodePrinterTest) TestFreeCallOptChainOptChainCall() {
	n := c.Parse("(a?.b)?.()")
	call := n.GetFirstFirstChild()
	AssertThat(call.IsOptChainCall()).IsTrue()
	call.PutBooleanProp(rhino.Node_FREE_CALL, true)
	c.AssertPrintNode("(0,a?.b)?.()", n)
}

// Java: CodePrinterTest.optChainCalleeForNewRequiresParentheses()
func (c *CodePrinterTest) TestOptChainCalleeForNewRequiresParentheses() {
	c.AssertPrintSame("new (a?.b)")
}

// Java: CodePrinterTest.testFreeCall1()
func (c *CodePrinterTest) TestFreeCall1() {
	c.AssertPrint("foo(a);", "foo(a)")
	c.AssertPrint("x.foo(a);", "x.foo(a)")
}

// Java: CodePrinterTest.testFreeCall2()
func (c *CodePrinterTest) TestFreeCall2() {
	n := c.Parse("foo(a);")
	c.AssertPrintNode("foo(a)", n)
	call := n.GetFirstFirstChild()
	AssertThat(call.IsCall()).IsTrue()
	call.PutBooleanProp(rhino.Node_FREE_CALL, true)
	c.AssertPrintNode("foo(a)", n)
}

// Java: CodePrinterTest.testFreeCall3()
func (c *CodePrinterTest) TestFreeCall3() {
	n := c.Parse("x.foo(a);")
	c.AssertPrintNode("x.foo(a)", n)
	call := n.GetFirstFirstChild()
	AssertThat(call.IsCall()).IsTrue()
	call.PutBooleanProp(rhino.Node_FREE_CALL, true)
	c.AssertPrintNode("(0,x.foo)(a)", n)
}

// Java: CodePrinterTest.testPrintScript()
func (c *CodePrinterTest) TestPrintScript() {
	// Verify that SCRIPT nodes not marked as synthetic are printed as
	// blocks.
	ast := rhino.NewNodeWith2Children(
		rhino.Token_SCRIPT,
		rhino.NewNodeWithChild(rhino.Token_EXPR_RESULT, rhino.Node_NewString("f")),
		rhino.NewNodeWithChild(rhino.Token_EXPR_RESULT, rhino.Node_NewString("g")))
	result := NewCodePrinter_Builder(ast).SetPrettyPrint(true).Build()
	AssertThat(result).IsEqualTo("\"f\";\n\"g\";\n")
}

// Java: CodePrinterTest.testObjectLit()
func (c *CodePrinterTest) TestObjectLit() {
	c.AssertPrint("({x:1})", "({x:1})")
	c.AssertPrint("var x=({x:1})", "var x={x:1}")
	c.AssertPrint("var x={'x':1}", "var x={\"x\":1}")
	c.AssertPrint("var x={1:1}", "var x={1:1}")
	c.AssertPrint("({},42)+0", "({},42)+0")
}

// Java: CodePrinterTest.testObjectLit2()
func (c *CodePrinterTest) TestObjectLit2() {
	c.AssertPrint("var x={1:1}", "var x={1:1}")
	c.AssertPrint("var x={'1':1}", "var x={1:1}")
	c.AssertPrint("var x={'1.0':1}", "var x={\"1.0\":1}")
	c.AssertPrint("var x={1.5:1}", "var x={\"1.5\":1}")
}

// Java: CodePrinterTest.testObjectLit3()
func (c *CodePrinterTest) TestObjectLit3() {
	c.AssertPrint("var x={3E9:1}", "var x={3E9:1}")
	c.AssertPrint(
		"var x={'3000000000':1}", // More than 31 bits
		"var x={3E9:1}")
	c.AssertPrint("var x={'3000000001':1}", "var x={3000000001:1}")
	c.AssertPrint(
		"var x={'6000000001':1}", // More than 32 bits
		"var x={6000000001:1}")
	c.AssertPrint(
		"var x={\"12345678901234567\":1}", // More than 53 bits
		"var x={\"12345678901234567\":1}")
}

// Java: CodePrinterTest.testObjectLit4()
func (c *CodePrinterTest) TestObjectLit4() {
	// More than 128 bits.
	c.AssertPrint(
		"var x={\"123456789012345671234567890123456712345678901234567\":1}",
		"var x={\"123456789012345671234567890123456712345678901234567\":1}")
}

// Java: CodePrinterTest.testExtendedObjectLit()
func (c *CodePrinterTest) TestExtendedObjectLit() {
	c.AssertPrintSame("var a={b}")
	c.AssertPrintSame("var a={b,c}")
	c.AssertPrintSame("var a={b,c:d,e}")
	c.AssertPrintSame("var a={b,c(){},d,e:f}")
}

// Java: CodePrinterTest.testComputedProperties()
func (c *CodePrinterTest) TestComputedProperties() {
	c.AssertPrintSame("var a={[b]:c}")
	c.AssertPrintSame("var a={[b+3]:c}")

	c.AssertPrintSame("var a={[b](){}}")
	c.AssertPrintSame("var a={[b](){alert(foo)}}")
	c.AssertPrintSame("var a={*[b](){yield\"foo\"}}")
	c.AssertPrintSame("var a={[b]:()=>c}")

	c.AssertPrintSame("var a={get [b](){return null}}")
	c.AssertPrintSame("var a={set [b](val){window.b=val}}")
}

// Java: CodePrinterTest.testComputedPropertiesClassMethods()
func (c *CodePrinterTest) TestComputedPropertiesClassMethods() {
	c.AssertPrintSame("class C{[m](){}}")

	c.AssertPrintSame("class C{[\"foo\"+bar](){alert(1)}}")
}

// Java: CodePrinterTest.testGetter()
func (c *CodePrinterTest) TestGetter() {
	c.AssertPrint("var x = {}", "var x={}")
	c.AssertPrint("var x = {get a() {return 1}}", "var x={get a(){return 1}}")
	c.AssertPrint("var x = {get a() {}, get b(){}}", "var x={get a(){},get b(){}}")

	c.AssertPrint("var x = {get 'a'() {return 1}}", "var x={get \"a\"(){return 1}}")

	c.AssertPrint("var x = {get 1() {return 1}}", "var x={get 1(){return 1}}")

	c.AssertPrint("var x = {get \"()\"() {return 1}}", "var x={get \"()\"(){return 1}}")

	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT5
	c.AssertPrintSame("var x={get function(){return 1}}")
}

// Java: CodePrinterTest.testGetterInEs3()
func (c *CodePrinterTest) TestGetterInEs3() {
	// Getters and setters and not supported in ES3 but if someone sets the
	// the ES3 output mode on an AST containing them we still produce them.
	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT3

	getter := rhino.Node_NewStringWithToken(rhino.Token_GETTER_DEF, "f")
	getter.AddChildToBack(NodeUtil_EmptyFunction())
	c.AssertPrintNode("({get f(){}})", rhino.IR_ExprResult(rhino.IR_Objectlit(getter)))
}

// Java: CodePrinterTest.testSetter()
func (c *CodePrinterTest) TestSetter() {
	c.AssertPrint("var x = {}", "var x={}")
	c.AssertPrint("var x = {set a(y) {return 1}}", "var x={set a(y){return 1}}")

	c.AssertPrint("var x = {get 'a'() {return 1}}", "var x={get \"a\"(){return 1}}")

	c.AssertPrint("var x = {set 1(y) {return 1}}", "var x={set 1(y){return 1}}")

	c.AssertPrint("var x = {set \"(x)\"(y) {return 1}}", "var x={set \"(x)\"(y){return 1}}")

	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT5
	c.AssertPrintSame("var x={set function(x){}}")
}

// Java: CodePrinterTest.testSetterInEs3()
func (c *CodePrinterTest) TestSetterInEs3() {
	// Getters and setters and not supported in ES3 but if someone sets the
	// the ES3 output mode on an AST containing them we still produce them.
	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT3

	getter := rhino.Node_NewStringWithToken(rhino.Token_SETTER_DEF, "f")
	getter.AddChildToBack(rhino.IR_Function(rhino.IR_Name(""), rhino.IR_ParamList(rhino.IR_Name("a")), rhino.IR_Block()))
	c.AssertPrintNode("({set f(a){}})", rhino.IR_ExprResult(rhino.IR_Objectlit(getter)))
}

// Java: CodePrinterTest.testNegNoCollapse()
func (c *CodePrinterTest) TestNegNoCollapse() {
	c.AssertPrint("var x = - - 2;", "var x=- -2")
	c.AssertPrint("var x = - (2);", "var x=-2")
}

// Java: CodePrinterTest.defaultBuilder(Node)
func (c *CodePrinterTest) DefaultBuilder(jsRoot *rhino.Node) *CodePrinter_Builder {
	return NewCodePrinter_Builder(jsRoot).
		SetCompilerOptions(
			c.NewCompilerOptions(
				codePrinterTestBase_CompilerOptionBuilderFunc(func(options *CompilerOptions) {
					options.SetPrettyPrint(false)
					options.SetLineBreak(false)
					options.SetLineLengthThreshold(0)
				}))).
		SetOutputTypes(false).
		SetTypeRegistry(c.lastCompiler.GetTypeRegistry())
}

// Java: CodePrinterTest.testStrict()
func (c *CodePrinterTest) TestStrict() {
	result := c.DefaultBuilder(c.ParseTypeChecked("var x" /* typeChecked= */, true)).SetTagAsStrict(true).Build()
	AssertThat(result).IsEqualTo("'use strict';var x")
}

// Java: CodePrinterTest.testStrictPretty()
func (c *CodePrinterTest) TestStrictPretty() {
	result := c.DefaultBuilder(c.ParseTypeChecked("var x" /* typeChecked= */, true)).
		SetTagAsStrict(true).
		SetPrettyPrint(true).
		Build()
	AssertThat(result).IsEqualTo("'use strict';\nvar x;\n")
}

// Java: CodePrinterTest.testIjs()
func (c *CodePrinterTest) TestIjs() {
	result := c.DefaultBuilder(c.ParseTypeChecked("var x" /* typeChecked= */, true)).SetTagAsTypeSummary(true).Build()
	AssertThat(result).IsEqualTo("/** @fileoverview @typeSummary */\nvar x")
}

// Java: CodePrinterTest.testIjsWithProvideAlreadyProvided()
func (c *CodePrinterTest) TestIjsWithProvideAlreadyProvided() {
	c.AssertPrettyPrintSame(
		("/** @provideAlreadyProvided */ \n" +
			"goog.provide(\"a.b.c\");\n"))
}

// Java: CodePrinterTest.testArrayLiteral()
func (c *CodePrinterTest) TestArrayLiteral() {
	c.AssertPrint("var x = [,];", "var x=[,]")
	c.AssertPrint("var x = [,,];", "var x=[,,]")
	c.AssertPrint("var x = [,s,,];", "var x=[,s,,]")
	c.AssertPrint("var x = [,s];", "var x=[,s]")
	c.AssertPrint("var x = [s,];", "var x=[s]")
}

// Java: CodePrinterTest.testZero()
func (c *CodePrinterTest) TestZero() {
	c.AssertPrint("var x ='\\0';", "var x=\"\\x00\"")
	c.AssertPrint("var x ='\\x00';", "var x=\"\\x00\"")
	c.AssertPrint("var x ='\\u0000';", "var x=\"\\x00\"")
	c.AssertPrint("var x ='\\u00003';", "var x=\"\\x003\"")
}

// Java: CodePrinterTest.testOctalInString()
func (c *CodePrinterTest) TestOctalInString() {
	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT5
	c.AssertPrint("var x ='\\0';", "var x=\"\\x00\"")
	c.AssertPrint("var x ='\\07';", "var x=\"\\u0007\"")

	// Octal 12 = Hex 0A = \n
	c.AssertPrint("var x ='\\012';", "var x=\"\\n\"")

	// Octal 13 = Hex 0B = \v
	c.AssertPrint("var x ='\\013';", "var x=\"\\v\"")

	// Octal 34 = Hex 1C
	c.AssertPrint("var x ='\\034';", "var x=\"\\u001c\"")

	// 8 and 9 are not octal digits
	c.AssertPrint("var x ='\\08';", "var x=\"\\x008\"")
	c.AssertPrint("var x ='\\09';", "var x=\"\\x009\"")

	// Only the first two digits are part of the octal literal.
	c.AssertPrint("var x ='\\01234';", "var x=\"\\n34\"")
}

// Java: CodePrinterTest.testOctalInStringNoLeadingZero()
func (c *CodePrinterTest) TestOctalInStringNoLeadingZero() {
	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT5
	c.AssertPrint("var x ='\\7';", "var x=\"\\u0007\"")

	// Octal 12 = Hex 0A = \n
	c.AssertPrint("var x ='\\12';", "var x=\"\\n\"")

	// Octal 13 = Hex 0B = \v.
	c.AssertPrint("var x ='\\13';", "var x=\"\\v\"")

	// Octal 34 = Hex 1C
	c.AssertPrint("var x ='\\34';", "var x=\"\\u001c\"")

	// Octal 240 = Hex A0
	c.AssertPrint("var x ='\\240';", "var x=\"\\u00a0\"")

	// Only the first three digits are part of the octal literal.
	c.AssertPrint("var x ='\\2400';", "var x=\"\\u00a00\"")

	// Only the first two digits are part of the octal literal because '8'
	// is not an octal digit.
	// Octal 67 = Hex 37 = "7"
	c.AssertPrint("var x ='\\6789';", "var x=\"789\"")

	// 8 and 9 are not octal digits. '\' is ignored and the digit
	// is just a regular character.
	c.AssertPrint("var x ='\\8';", "var x=\"8\"")
	c.AssertPrint("var x ='\\9';", "var x=\"9\"")

	// Only the first three digits are part of the octal literal.
	// Octal 123 = Hex 53 = "S"
	c.AssertPrint("var x ='\\1234';", "var x=\"S4\"")
}

// Java: CodePrinterTest.testUnicode()
func (c *CodePrinterTest) TestUnicode() {
	c.AssertPrint("var x ='\\x0f';", "var x=\"\\u000f\"")
	c.AssertPrint("var x ='\\x68';", "var x=\"h\"")
	c.AssertPrint("var x ='\\x7f';", "var x=\"\\u007f\"")
}

// Separate from testNumericKeys() so we can set allowWarnings.
// Java: CodePrinterTest.testOctalNumericKey()
func (c *CodePrinterTest) TestOctalNumericKey() {
	c.allowWarnings = true
	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT5

	c.AssertPrint("var x = {010: 1};", "var x={8:1}")
}

// Java: CodePrinterTest.testNumericKeys()
func (c *CodePrinterTest) TestNumericKeys() {
	c.AssertPrint("var x = {'010': 1};", "var x={\"010\":1}")

	c.AssertPrint("var x = {0x10: 1};", "var x={16:1}")
	c.AssertPrint("var x = {'0x10': 1};", "var x={\"0x10\":1}")

	// I was surprised at this result too.
	c.AssertPrint("var x = {.2: 1};", "var x={\"0.2\":1}")
	c.AssertPrint("var x = {'.2': 1};", "var x={\".2\":1}")

	c.AssertPrint("var x = {0.2: 1};", "var x={\"0.2\":1}")
	c.AssertPrint("var x = {'0.2': 1};", "var x={\"0.2\":1}")
}

// Java: CodePrinterTest.testIssue582()
func (c *CodePrinterTest) TestIssue582() {
	c.AssertPrint("var x = -0.0;", "var x=-0")
}

// Java: CodePrinterTest.testIssue942()
func (c *CodePrinterTest) TestIssue942() {
	c.AssertPrint("var x = {0: 1};", "var x={0:1}")
}

// Java: CodePrinterTest.testIssue601()
func (c *CodePrinterTest) TestIssue601() {
	c.AssertPrint("'\\v' == 'v'", "\"\\v\"==\"v\"")
	c.AssertPrint("'\\u000B' == '\\v'", "\"\\v\"==\"\\v\"")
	c.AssertPrint("'\\x0B' == '\\v'", "\"\\v\"==\"\\v\"")
}

// Java: CodePrinterTest.testIssue620()
func (c *CodePrinterTest) TestIssue620() {
	c.AssertPrint("alert(/ / / / /);", "alert(/ // / /)")
	c.AssertPrint("alert(/ // / /);", "alert(/ // / /)")
}

// Java: CodePrinterTest.testIssue5746867()
func (c *CodePrinterTest) TestIssue5746867() {
	c.AssertPrint("var a = { '$\\\\' : 5 };", "var a={\"$\\\\\":5}")
}

// Java: CodePrinterTest.testCommaSpacing()
func (c *CodePrinterTest) TestCommaSpacing() {
	c.AssertPrint("var a = (b = 5, c = 5);", "var a=(b=5,c=5)")
	c.AssertPrettyPrint("var a = (b = 5, c = 5);", "var a = (b = 5, c = 5);\n")
}

// Java: CodePrinterTest.testManyCommas()
func (c *CodePrinterTest) TestManyCommas() {
	numCommas := 10000
	var numbers []string
	numbers = append(numbers, "0")
	numbers = append(numbers, "1")
	current := rhino.NewNodeWith2Children(rhino.Token_COMMA, rhino.Node_NewNumber(0), rhino.Node_NewNumber(1))
	for i := 2; i < numCommas; i++ {
		current = rhino.NewNodeWithChild(rhino.Token_COMMA, current)

		// 1000 is printed as 1E3, and screws up our test.
		num := i % 1000
		numbers = append(numbers, strconv.Itoa(num))
		current.AddChildToBack(rhino.Node_NewNumber(float64(num)))
	}

	expected := strings.Join(numbers, ",")
	actual := strings.ReplaceAll(c.PrintNode(current), "\n", "")
	AssertThat(actual).IsEqualTo(expected)
}

// Java: CodePrinterTest.testManyAdds()
func (c *CodePrinterTest) TestManyAdds() {
	numAdds := 10000
	var numbers []string
	numbers = append(numbers, "0")
	numbers = append(numbers, "1")
	current := rhino.NewNodeWith2Children(rhino.Token_ADD, rhino.Node_NewNumber(0), rhino.Node_NewNumber(1))
	for i := 2; i < numAdds; i++ {
		current = rhino.NewNodeWithChild(rhino.Token_ADD, current)

		// 1000 is printed as 1E3, and screws up our test.
		num := i % 1000
		numbers = append(numbers, strconv.Itoa(num))
		current.AddChildToBack(rhino.Node_NewNumber(float64(num)))
	}

	expected := strings.Join(numbers, "+")
	actual := strings.ReplaceAll(c.PrintNode(current), "\n", "")
	AssertThat(actual).IsEqualTo(expected)
}

// Java: CodePrinterTest.testMinusNegativeZero()
func (c *CodePrinterTest) TestMinusNegativeZero() {
	// Negative zero is weird, because we have to be able to distinguish
	// it from positive zero (there are some subtle differences in behavior).
	c.AssertPrint("x- -0", "x- -0")
}

// Java: CodePrinterTest.testStringEscapeSequences()
func (c *CodePrinterTest) TestStringEscapeSequences() {
	// From the SingleEscapeCharacter grammar production.
	c.AssertPrintSame("var x=\"\\b\"")
	c.AssertPrintSame("var x=\"\\f\"")
	c.AssertPrintSame("var x=\"\\n\"")
	c.AssertPrintSame("var x=\"\\r\"")
	c.AssertPrintSame("var x=\"\\t\"")
	c.AssertPrintSame("var x=\"\\v\"")
	c.AssertPrint("var x=\"\\\"\"", "var x='\"'")
	c.AssertPrint("var x=\"\\'\"", "var x=\"'\"")

	// From the LineTerminator grammar
	c.AssertPrint("var x=\"\\u000A\"", "var x=\"\\n\"")
	c.AssertPrint("var x=\"\\u000D\"", "var x=\"\\r\"")
	c.AssertPrintSame("var x=\"\\u2028\"")
	c.AssertPrintSame("var x=\"\\u2029\"")

	// Now with regular expressions.
	c.AssertPrintSame("var x=/\\b/")
	c.AssertPrintSame("var x=/\\f/")
	c.AssertPrintSame("var x=/\\n/")
	c.AssertPrintSame("var x=/\\r/")
	c.AssertPrintSame("var x=/\\t/")
	c.AssertPrintSame("var x=/\\v/")
	c.AssertPrintSame("var x=/\\u000A/")
	c.AssertPrintSame("var x=/\\u000D/")
	c.AssertPrintSame("var x=/\\u2028/")
	c.AssertPrintSame("var x=/\\u2029/")
}

// Java: CodePrinterTest.testRegexp_escape()
func (c *CodePrinterTest) TestRegexp_escape() {
	c.AssertPrintSame("/\\bword\\b/")
	c.AssertPrintSame("/Java\\BScript/")
	c.AssertPrintSame("/\\ca/")
	c.AssertPrintSame("/\\cb/")
	c.AssertPrintSame("/\\cc/")
	c.AssertPrintSame("/\\cA/")
	c.AssertPrintSame("/\\cB/")
	c.AssertPrintSame("/\\cC/")
	c.AssertPrintSame("/\\d/")
	c.AssertPrintSame("/\\D/")
	c.AssertPrintSame("/\\0/")
	c.AssertPrintSame("/\\\\/")
	c.AssertPrintSame("/(.)\\1/")
	c.AssertPrintSame("/\\k<group>/")
	c.AssertPrintSame("/(?<group>a)\\k<group>/")
	c.AssertPrintSame("/\\x0B/") // Don't print this as \v (as is done in strings)
}

// Java: CodePrinterTest.testRegexp_unnecessaryEscape()
func (c *CodePrinterTest) TestRegexp_unnecessaryEscape() {
	c.AssertPrint("/\\a/", "/a/")
	c.AssertPrint("/\\e/", "/e/")
	c.AssertPrint("/\\g/", "/g/")
	c.AssertPrint("/\\h/", "/h/")
	c.AssertPrint("/\\i/", "/i/")
	c.AssertPrint("/\\¡/", "/\\u00a1/")
}

// Java: CodePrinterTest.testRegexp_afterLt()
func (c *CodePrinterTest) TestRegexp_afterLt() {
	c.AssertPrint("x < /script/", "x< /script/")
	c.AssertPrint("x << /script/", "x<< /script/")
}

// Java: CodePrinterTest.testKeywordProperties1()
func (c *CodePrinterTest) TestKeywordProperties1() {
	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT5
	c.AssertPrintSame("x.foo=2")
	c.AssertPrintSame("x.function=2")

	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT3
	c.AssertPrintSame("x.foo=2")
}

// Java: CodePrinterTest.testKeywordProperties1a()
func (c *CodePrinterTest) TestKeywordProperties1a() {
	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT5
	nodes := c.Parse("x.function=2")
	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT3
	c.AssertPrintNode("x[\"function\"]=2", nodes)
}

// Java: CodePrinterTest.testKeywordProperties2()
func (c *CodePrinterTest) TestKeywordProperties2() {
	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT5
	c.AssertPrintSame("x={foo:2}")
	c.AssertPrintSame("x={function:2}")

	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT3
	c.AssertPrintSame("x={foo:2}")
}

// Java: CodePrinterTest.testKeywordProperties2a()
func (c *CodePrinterTest) TestKeywordProperties2a() {
	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT5
	nodes := c.Parse("x={function:2}")
	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT3
	c.AssertPrintNode("x={\"function\":2}", nodes)
}

// Java: CodePrinterTest.testIssue1062()
func (c *CodePrinterTest) TestIssue1062() {
	c.AssertPrintSame("3*(4%3*5)")
}

// Java: CodePrinterTest.testPreserveTypeAnnotations()
func (c *CodePrinterTest) TestPreserveTypeAnnotations() {
	c.preserveTypeAnnotations = true
	c.AssertPrintSame("/** @type {foo} */ var bar")
	c.AssertPrintSame("function/** void */ f(/** string */ s,/** number */ n){}")

	c.preserveTypeAnnotations = false
	c.AssertPrint("/** @type {foo} */ var bar;", "var bar")
}

// Java: CodePrinterTest.testPreserveTypeAnnotations2()
func (c *CodePrinterTest) TestPreserveTypeAnnotations2() {
	c.preserveTypeAnnotations = true

	c.AssertPrintSame("/** @const */ var ns={}")

	c.AssertPrintSame(
		("/**\n" +
			" * @const\n" +
			" * @suppress {const,duplicate}\n" +
			" */\n" +
			"var ns={}\n"))
}

// Java: CodePrinterTest.testDefaultParameters()
func (c *CodePrinterTest) TestDefaultParameters() {
	c.AssertPrintSame("function f(a=0){}")
	c.AssertPrintSame("function f(a,b=0){}")
	c.AssertPrintSame("function f(a=(1,2),b){}")
	c.AssertPrintSame("function f(a,b=(1,2)){}")
}

// Java: CodePrinterTest.testRestParameters()
func (c *CodePrinterTest) TestRestParameters() {
	c.AssertPrintSame("function f(...args){}")
	c.AssertPrintSame("function f(first,...rest){}")
}

// Java: CodePrinterTest.testDefaultParametersWithRestParameters()
func (c *CodePrinterTest) TestDefaultParametersWithRestParameters() {
	c.AssertPrintSame("function f(first=0,...args){}")
	c.AssertPrintSame("function f(first,second=0,...rest){}")
}

// Java: CodePrinterTest.testSpreadExpression()
func (c *CodePrinterTest) TestSpreadExpression() {
	c.AssertPrintSame("f(...args)")
	c.AssertPrintSame("f(...arrayOfArrays[0])")
	c.AssertPrintSame("f(...[1,2,3])")
	c.AssertPrintSame("f(...([1],[2]))")
}

// Java: CodePrinterTest.testClass()
func (c *CodePrinterTest) TestClass() {
	c.AssertPrintSame("class C{}")
	c.AssertPrintSame("(class C{})")
	c.AssertPrintSame("class C extends D{}")
	c.AssertPrintSame("class C{static member(){}}")
	c.AssertPrintSame("class C{member(){}get f(){}}")
	c.AssertPrintSame("var x=class C{}")
}

// Java: CodePrinterTest.testClassComputedProperties()
func (c *CodePrinterTest) TestClassComputedProperties() {

	c.AssertPrintSame("class C{[x](){}}")
	c.AssertPrintSame("class C{get [x](){}}")
	c.AssertPrintSame("class C{set [x](val){}}")

	c.AssertPrintSame("class C{static [x](){}}")
	c.AssertPrintSame("class C{static get [x](){}}")
	c.AssertPrintSame("class C{static set [x](val){}}")
}

// Java: CodePrinterTest.testClassPretty()
func (c *CodePrinterTest) TestClassPretty() {
	c.AssertPrettyPrint("class C{}", "class C {\n}\n")
	c.AssertPrettyPrint(
		"class C{member(){}get f(){}}",
		("class C {\n" +
			"  member() {\n" +
			"  }\n" +
			"  get f() {\n" +
			"  }\n" +
			"}\n"))
	c.AssertPrettyPrint("var x=class C{}", "var x = class C {\n};\n")
}

// Java: CodePrinterTest.testClassField()
func (c *CodePrinterTest) TestClassField() {
	c.AssertPrettyPrintSame(
		("class C {\n" +
			"  x;\n" +
			"}\n"))
	c.AssertPrettyPrintSame(
		("class C {\n" +
			"  x = 2;\n" +
			"}\n"))
	c.AssertPrettyPrintSame(
		("class C {\n" +
			"  x = 2;\n" +
			"  y = 3;\n" +
			"}\n"))
	c.AssertPrettyPrint(
		("class C {\n" +
			"  x = 2\n" +
			"  y = 3\n" +
			"}\n"),
		("class C {\n" +
			"  x = 2;\n" +
			"  y = 3;\n" +
			"}\n"))
	c.AssertPrettyPrint(
		"class C {x=2;y=3}",
		("class C {\n" +
			"  x = 2;\n" +
			"  y = 3;\n" +
			"}\n"))
}

// Java: CodePrinterTest.testClassFieldCheckState()
func (c *CodePrinterTest) TestClassFieldCheckState() {
	c.AssertPrettyPrintSame(
		("/** @interface */ \n" +
			"class C {\n" +
			"  x;\n" +
			"}\n"))
	c.AssertPrettyPrintSame(
		("/** @record */ \n" +
			"class C {\n" +
			"  x;\n" +
			"}\n"))
}

// Java: CodePrinterTest.testClassFieldStatic()
func (c *CodePrinterTest) TestClassFieldStatic() {
	c.AssertPrettyPrintSame(
		("class C {\n" +
			"  static x;\n" +
			"}\n"))
	c.AssertPrettyPrintSame(
		("class C {\n" +
			"  static x = 2;\n" +
			"}\n"))
	c.AssertPrettyPrintSame(
		("class C {\n" +
			"  static x = 2;\n" +
			"  static y = 3;\n" +
			"}\n"))
	c.AssertPrettyPrintSame(
		("/** @interface */ \n" +
			"class C {\n" +
			"  static x;\n" +
			"}\n"))
	c.AssertPrettyPrintSame(
		("/** @record */ \n" +
			"class C {\n" +
			"  static x;\n" +
			"}\n"))
}

// Java: CodePrinterTest.testComputedClassFieldLiteralStringNumber()
func (c *CodePrinterTest) TestComputedClassFieldLiteralStringNumber() {
	c.AssertPrettyPrint(
		"class C { 'str' = 2;}",
		("class C {\n" +
			"  [\"str\"] = 2;\n" +
			"}\n"))
	c.AssertPrettyPrint(
		"class C { 1 = 2;}",
		("class C {\n" +
			"  [1] = 2;\n" +
			"}\n"))
}

// Java: CodePrinterTest.testComputedClassField()
func (c *CodePrinterTest) TestComputedClassField() {
	c.AssertPrettyPrintSame(
		("class C {\n" +
			"  [x];\n" +
			"}\n"))
	c.AssertPrettyPrintSame(
		("class C {\n" +
			"  [x] = 2;\n" +
			"}\n"))
	c.AssertPrettyPrintSame(
		("class C {\n" +
			"  [x] = 2;\n" +
			"  y = 3;\n" +
			"}\n"))
	c.AssertPrettyPrintSame(
		("class C {\n" +
			"  [x] = 2;\n" +
			"  [y] = 3;\n" +
			"}\n"))
}

// Java: CodePrinterTest.testComputedClassFieldStatic()
func (c *CodePrinterTest) TestComputedClassFieldStatic() {
	c.AssertPrettyPrintSame(
		("class C {\n" +
			"  static [x];\n" +
			"}\n"))
	c.AssertPrettyPrintSame(
		("class C {\n" +
			"  static [x] = 2;\n" +
			"}\n"))
	c.AssertPrettyPrintSame(
		("class C {\n" +
			"  static [x] = 2;\n" +
			"  static y = 3;\n" +
			"}\n"))
	c.AssertPrettyPrintSame(
		("class C {\n" +
			"  static [x] = 2;\n" +
			"  static [y] = 3;\n" +
			"}\n"))
}

// Java: CodePrinterTest.testPrivateClassProperties_definition()
func (c *CodePrinterTest) TestPrivateClassProperties_definition() {
	c.AssertPrintSame("class C{#f;}")
	c.AssertPrintSame("class C{#m(){}}")
	c.AssertPrintSame("class C{*#g(){}}")
	c.AssertPrintSame("class C{get #g(){}}")
	c.AssertPrintSame("class C{set #s(x){}}")
	c.AssertPrintSame("class C{get #p(){}set #p(x){}}")
	c.AssertPrintSame("class C{async #a(){}}")
	c.AssertPrintSame("class C{async *#ag(){}}")

	c.AssertPrintSame("class C{static #sf;}")
	c.AssertPrintSame("class C{static #sm(){}}")
	c.AssertPrintSame("class C{static *#sg(){}}")
	c.AssertPrintSame("class C{static get #sg(){}}")
	c.AssertPrintSame("class C{static set #ss(x){}}")
	c.AssertPrintSame("class C{static get #sp(){}static set #sp(x){}}")
	c.AssertPrintSame("class C{static async #sa(){}}")
	c.AssertPrintSame("class C{static async *#sag(){}}")
}

// Java: CodePrinterTest.testPrivateClassProperties_usage()
func (c *CodePrinterTest) TestPrivateClassProperties_usage() {
	c.AssertPrintSame("class C{#f;#g=this.#f;}")
	c.AssertPrintSame("class C{#f;#g=this?.#f;}")
	c.AssertPrintSame("class C{#m(){this.#m()}}")
	c.AssertPrintSame("class C{#m(){this?.#m()}}")

	c.AssertPrintSame("class C{static #f;static #m(){const t=this;t.#f}}")
	c.AssertPrintSame("class C{static #f;static #m(){const t=this;t?.#f}}")
	c.AssertPrintSame("class C{static #f;static #m(){const t=this;#f in t}}")
}

// Java: CodePrinterTest.testSuper()
func (c *CodePrinterTest) TestSuper() {
	c.AssertPrintSame("class C extends foo(){}")
	c.AssertPrintSame("class C extends m.foo(){}")
	c.AssertPrintSame("class C extends D{member(){super.foo()}}")
}

// Java: CodePrinterTest.testNewTarget()
func (c *CodePrinterTest) TestNewTarget() {
	c.AssertPrintSame("function f(){new.target}")
	c.AssertPrint("function f() {\nnew\n.\ntarget;\n}", "function f(){new.target}")
}

// Java: CodePrinterTest.testImportMeta()
func (c *CodePrinterTest) TestImportMeta() {
	c.AssertPrintSame("import.meta")
	c.AssertPrintSame("import.meta.url")
	c.AssertPrintSame("console.log(import.meta.url)")
}

// Java: CodePrinterTest.testGeneratorYield()
func (c *CodePrinterTest) TestGeneratorYield() {
	c.AssertPrintSame("function*f(){yield 1}")
	c.AssertPrintSame("function*f(){yield}")
	c.AssertPrintSame("function*f(){yield 1?0:2}")
	c.AssertPrintSame("function*f(){yield 1,0}")
	c.AssertPrintSame("function*f(){1,yield 0}")
	c.AssertPrintSame("function*f(){yield(a=0)}")
	c.AssertPrintSame("function*f(){a=yield 0}")
	c.AssertPrintSame("function*f(){(yield 1)+(yield 1)}")
	// Parens required for evaluating arrow function expression i.e. `yield (() => expr)`
	c.AssertPrintSame("function*f(){yield(()=>({}))}")
}

// Java: CodePrinterTest.testGeneratorYieldPretty()
func (c *CodePrinterTest) TestGeneratorYieldPretty() {
	c.AssertPrettyPrint(
		"function *f() {yield 1}",
		("function* f() {\n" +
			"  yield 1;\n" +
			"}\n"))

	c.AssertPrettyPrint(
		"function *f() {yield}",
		("function* f() {\n" +
			"  yield;\n" +
			"}\n"))
}

// Java: CodePrinterTest.testMemberGeneratorYield1()
func (c *CodePrinterTest) TestMemberGeneratorYield1() {
	c.AssertPrintSame("class C{*member(){(yield 1)+(yield 1)}}")
	c.AssertPrintSame("class C{*[0](){(yield 1)+(yield 1)}}")
	c.AssertPrintSame("var obj={*member(){(yield 1)+(yield 1)}}")
	c.AssertPrintSame("var obj={*[0](){(yield 1)+(yield 1)}}")
	c.AssertPrintSame("var obj={[0]:function*(){(yield 1)+(yield 1)}}")
}

// Java: CodePrinterTest.testArrowFunction_zeroParams()
func (c *CodePrinterTest) TestArrowFunction_zeroParams() {
	c.AssertPrintSame("()=>1")
	c.AssertPrint("(()=>1)", "()=>1")
	c.AssertPrintSame("()=>{}")
	c.AssertPrint("(()=>a),b", "()=>a,b")
	c.AssertPrint("()=>(a=b)", "()=>a=b")
}

// Java: CodePrinterTest.testArrowFunction_oneParam()
func (c *CodePrinterTest) TestArrowFunction_oneParam() {
	c.AssertPrintSame("a=>b")
	c.AssertPrintSame("([a])=>b")
	c.AssertPrintSame("(...a)=>b")
	c.AssertPrintSame("(a=0)=>b")
	c.AssertPrintSame("(a=>b)(1)")
	c.AssertPrintSame("fn?.(a=>a)")
	c.AssertPrint("(a=>a)?.['length']", "(a=>a)?.[\"length\"]")
	c.AssertPrintSame("(a=>a)?.(1)")
	c.AssertPrintSame("(a=>a)?.length")
	c.AssertPrintSame("a=>a?.length")
	c.AssertPrintSame("var z={x:a=>1}")
	c.AssertPrintSame("[1,2].forEach(x=>y)")
}

// Java: CodePrinterTest.testArrowFunction_manyParams()
func (c *CodePrinterTest) TestArrowFunction_manyParams() {
	c.AssertPrintSame("(a,b)=>b")
}

// Java: CodePrinterTest.testArrowFunction_bodyEdgeCases()
func (c *CodePrinterTest) TestArrowFunction_bodyEdgeCases() {
	c.AssertPrintSame("()=>(a,b)")
	c.AssertPrintSame("()=>({a:1})")
	c.AssertPrintSame("()=>{return 1}")
}

// Java: CodePrinterTest.testAsyncFunction()
func (c *CodePrinterTest) TestAsyncFunction() {
	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT_NEXT
	c.AssertPrintSame("async function f(){}")
	c.AssertPrintSame("let f=async function f(){}")
	c.AssertPrintSame("let f=async function(){}")
	// implicit semicolon prevents async being treated as a keyword
	c.AssertPrint("async\nfunction f(){}", "async;function f(){}")
	c.AssertPrint("let f=async\nfunction f(){}", "let f=async;function f(){}")
}

// Java: CodePrinterTest.testAsyncGeneratorFunction()
func (c *CodePrinterTest) TestAsyncGeneratorFunction() {
	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT_NEXT
	c.AssertPrintSame("async function*f(){}")
	c.AssertPrintSame("let f=async function*f(){}")
	c.AssertPrintSame("let f=async function*(){}")
	// implicit semicolon prevents async being treated as a keyword
	c.AssertPrint("async\nfunction*f(){}", "async;function*f(){}")
	c.AssertPrint("let f=async\nfunction*f(){}", "let f=async;function*f(){}")
}

// Java: CodePrinterTest.testAsyncArrowFunction()
func (c *CodePrinterTest) TestAsyncArrowFunction() {
	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT_NEXT
	c.AssertPrintSame("async()=>1")
	c.AssertPrint("async (a) => 1", "async a=>1")

	// implicit semicolon prevents async being treated as a keyword
	c.AssertPrint("f=async\n()=>1", "f=async;()=>1")
}

// Java: CodePrinterTest.testAsyncMethod()
func (c *CodePrinterTest) TestAsyncMethod() {
	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT_NEXT
	c.AssertPrintSame("o={async m(){}}")
	c.AssertPrintSame("o={async[a+b](){}}")
	c.AssertPrintSame("o={[0]:async function(){}}") // (not technically a method)
	c.AssertPrintSame("class C{async m(){}}")
	c.AssertPrintSame("class C{async[a+b](){}}")
	c.AssertPrintSame("class C{static async m(){}}")
	c.AssertPrintSame("class C{static async[a+b](){}}")
}

// Java: CodePrinterTest.testAsyncGeneratorMethod()
func (c *CodePrinterTest) TestAsyncGeneratorMethod() {
	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT_NEXT
	c.AssertPrintSame("o={async *m(){}}")
	c.AssertPrintSame("o={async*[a+b](){}}")
	c.AssertPrintSame("o={[0]:async*function(){}}") // (not technically a method)
	c.AssertPrintSame("class C{async *m(){}}")
	c.AssertPrintSame("class C{async*[a+b](){}}")
	c.AssertPrintSame("class C{static async *m(){}}")
	c.AssertPrintSame("class C{static async*[a+b](){}}")
}

// Java: CodePrinterTest.testAwaitExpression()
func (c *CodePrinterTest) TestAwaitExpression() {
	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT_NEXT
	c.AssertPrintSame("async function f(promise){return await promise}")
	c.AssertPrintSame("pwait=async function(promise){return await promise}")
	c.AssertPrintSame("class C{async pwait(promise){await promise}}")
	c.AssertPrintSame("o={async pwait(promise){await promise}}")
	c.AssertPrintSame("pwait=async promise=>await promise")
}

// Regression test for b/235871063 - necessary parens dropped around awaited arrow function.
// Java: CodePrinterTest.testParansAroudAwaitArrowFunction()
func (c *CodePrinterTest) TestParansAroudAwaitArrowFunction() {
	// Parens required for evaluating arrow function expression i.e. `await (() => expr)`
	c.AssertPrint(
		"async function f(){return await (()=>new Promise((resolve)=>setTimeout(resolve,0)));}",
		"async function f(){return await (()=>new Promise(resolve=>setTimeout(resolve,0)))}")
	fmt.Println("--------------")
	// Parens not required for evaluating new
	c.AssertPrint(
		"async function f(){return await new Promise((resolve)=>setTimeout(resolve,0));}",
		"async function f(){return await new Promise(resolve=>setTimeout(resolve,0))}")
}

// Regression test for b/28633247 - necessary parens dropped around arrow functions.
// <p>Many of these cases use single param arrows because their PARAM_LIST parens should also be
// dropped, which can make this harder to parse for humans.
// Java: CodePrinterTest.testParensAroundArrow()
func (c *CodePrinterTest) TestParensAroundArrow() {
	// Parens required for non-assignment binary operator
	c.AssertPrint("x||((_)=>true)", "x||(_=>true)")
	// Parens required for unary operator
	c.AssertPrint("void((e)=>e*5)", "void(e=>e*5)")
	// Parens not required for comma operator
	c.AssertPrint("((_) => true), ((_) => false)", "_=>true,_=>false")
	// Parens not required for right side of assignment operator
	// NOTE: An arrow function on the left side would be a parse error.
	c.AssertPrint("x = ((_) => _ + 1)", "x=_=>_+1")
	// Parens required for template tag
	c.AssertPrint("((_)=>\"\")`template`", "(_=>\"\")`template`")
	// Parens required to reference a property
	c.AssertPrintSame("((a,b,c)=>a+b+c).length")
	c.AssertPrintSame("((a,b,c)=>a+b+c)[\"length\"]")
	// Parens not required when evaluating property name.
	// (It doesn't make much sense to do it, though.)
	c.AssertPrint("x[((_)=>0)]", "x[_=>0]")
	// Parens required to call the arrow function immediately
	c.AssertPrint("((x)=>x*5)(10)", "(x=>x*5)(10)")
	// Parens not required for function call arguments
	c.AssertPrint("x(((_) => true), ((_) => false))", "x(_=>true,_=>false)")
	// Parens required for first operand to a conditional, but not the rest.
	c.AssertPrint("((x)=>1)?a:b", "(x=>1)?a:b")
	c.AssertPrint("x?((x)=>0):((x)=>1)", "x?x=>0:x=>1")
	c.AssertPrint("new ((x)=>x)", "new (x=>x)")
	c.AssertPrintSame("new C(x=>x)")
}

// Java: CodePrinterTest.testParensAroundArrowFnInCast()
func (c *CodePrinterTest) TestParensAroundArrowFnInCast() {
	c.preserveTypeAnnotations = false
	c.AssertPrint("x(/** @type {?} */ (()=>{x}))", "x(()=>{x})")
	c.AssertPrint("x(/** @type {?} */ (()=>{x})())", "x((()=>{x})())")
	c.AssertPrint("x(/** @type {string} */ (/** @type {?} */ (()=>{x}))())", "x((()=>{x})())")

	c.preserveTypeAnnotations = true
	c.AssertPrintSame("x(/** @type {?} */ (()=>{x}))")
	c.AssertPrintSame("x(/** @type {?} */ (()=>{x})())")
	c.AssertPrintSame("x(/** @type {string} */ (/** @type {?} */ (()=>{x}))())")
}

// Java: CodePrinterTest.testParensAroundVariableDeclarator()
func (c *CodePrinterTest) TestParensAroundVariableDeclarator() {
	c.AssertPrintSame("var o=(test++,{one:1})")
	c.AssertPrintSame("({one}=(test++,{one:1}))")
	c.AssertPrintSame("[one]=(test++,[1])")

	c.AssertPrintSame("var {one}=(test++,{one:1})")
	c.AssertPrintSame("var [one]=(test++,[1])")
	c.AssertPrint(
		"var {one}=/** @type {{one: number}} */(test++,{one:1})", "var {one}=(test++,{one:1})")
	c.AssertPrint("var [one]=/** @type {!Array<number>} */(test++,[1])", "var [one]=(test++,[1])")
}

// Java: CodePrinterTest.testParensAroundArrowReturnValue()
func (c *CodePrinterTest) TestParensAroundArrowReturnValue() {
	c.AssertPrintSame("()=>({})")
	c.AssertPrintSame("()=>({a:1})")
	c.AssertPrintSame("()=>({a:1,b:2})")
	c.AssertPrint("()=>/** @type {Object} */({})", "()=>({})")
	c.AssertPrint("()=>/** @type {Object} */({a:1})", "()=>({a:1})")
	c.AssertPrint("()=>/** @type {Object} */({a:1,b:2})", "()=>({a:1,b:2})")
	c.AssertPrint("()=>/** @type {number} */(3)", "()=>3")
	c.AssertPrint("()=>/** @type {Object} */ ({}={})", "()=>({}={})")

	c.AssertPrintSame("()=>(1,2)")
	c.AssertPrintSame("()=>({},2)")
	c.AssertPrintSame("()=>(1,{})")
	c.AssertPrint("()=>/** @type {?} */(1,2)", "()=>(1,2)")
	c.AssertPrint("()=>/** @type {?} */({},2)", "()=>({},2)")
	c.AssertPrint("()=>/** @type {?} */(1,{})", "()=>(1,{})")

	// Test object literals more deeply nested
	c.AssertPrintSame("fn=()=>({})||3")
	c.AssertPrintSame("fn=()=>3||{}")
	c.AssertPrintSame("fn=()=>({}={})")
	c.AssertPrintSame("()=>function(){}") // don't need parentheses around a function
	c.AssertPrintSame("for(var i=()=>({});;);")

	c.preserveTypeAnnotations = true
	c.AssertPrintSame("()=>/** @type {Object} */ ({})")
}

// Java: CodePrinterTest.testPrettyArrowFunction()
func (c *CodePrinterTest) TestPrettyArrowFunction() {
	c.AssertPrettyPrint(
		"if (x) {var f = ()=>{alert(1); alert(2)}}",
		("if (x) {\n" +
			"  var f = () => {\n" +
			"    alert(1);\n" +
			"    alert(2);\n" +
			"  };\n" +
			"}\n"))
}

// Java: CodePrinterTest.testPrettyPrint_switch()
func (c *CodePrinterTest) TestPrettyPrint_switch() {
	c.AssertPrettyPrint(
		"switch(something){case 0:alert(0);break;case 1:alert(1);break}",
		("switch(something) {\n" +
			"  case 0:\n" +
			"    alert(0);\n" +
			"    break;\n" +
			"  case 1:\n" +
			"    alert(1);\n" +
			"    break;\n" +
			"}\n"))
}

// Java: CodePrinterTest.testBlocksInCaseArePreserved()
func (c *CodePrinterTest) TestBlocksInCaseArePreserved() {
	js := ("switch(something) {\n" +
		"  case 0:\n" +
		"    {\n" +
		"      const x = 1;\n" +
		"      break;\n" +
		"    }\n" +
		"  case 1:\n" +
		"    break;\n" +
		"  case 2:\n" +
		"    console.log(`case 2!`);\n" +
		"    {\n" +
		"      const x = 2;\n" +
		"      break;\n" +
		"    }\n" +
		"}\n")
	c.AssertPrettyPrint(js, js)
}

// Java: CodePrinterTest.testBlocksArePreserved()
func (c *CodePrinterTest) TestBlocksArePreserved() {
	js := ("console.log(0);\n" +
		"{\n" +
		"  let x = 1;\n" +
		"  console.log(x);\n" +
		"}\n" +
		"console.log(x);\n")
	c.AssertPrettyPrint(js, js)
}

// Java: CodePrinterTest.testBlocksNotPreserved()
func (c *CodePrinterTest) TestBlocksNotPreserved() {
	c.AssertPrint("if (x) {};", "if(x);")
	c.AssertPrint("while (x) {};", "while(x);")
}

// Java: CodePrinterTest.testEmptyClassStaticBlock()
func (c *CodePrinterTest) TestEmptyClassStaticBlock() {
	c.AssertPrettyPrintSame("class C {\n  static {\n  }\n}\n")
	c.AssertPrint("class C {\n  static {\n  }\n}\n", "class C{static{}}")
	c.AssertPrettyPrintSame("let a = class {\n  static {\n  }\n};\n")
}

// Java: CodePrinterTest.testClassStaticBlock()
func (c *CodePrinterTest) TestClassStaticBlock() {
	c.AssertPrettyPrintSame(
		("class C {\n" +
			"  static field1 = 1;\n" +
			"  static field2 = 2;\n" +
			"  static {\n" +
			"    let x = this.field1;\n" +
			"    let y = this.field2;\n" +
			"  }\n" +
			"}\n"))
	c.AssertPrettyPrintSame(
		("class C {\n" +
			"  static {\n" +
			"    this.field1 = 1;\n" +
			"    this.field2 = 2;\n" +
			"  }\n" +
			"}\n"))
	c.AssertPrettyPrintSame(
		("let a = class {\n" +
			"  static field1 = 1;\n" +
			"  static field2 = 2;\n" +
			"  static {\n" +
			"    let x = this.field1;\n" +
			"    let y = this.field2;\n" +
			"  }\n" +
			"};\n"))
	c.AssertPrettyPrintSame(
		("let a = class {\n" +
			"  static {\n" +
			"    this.field1 = 1;\n" +
			"    this.field2 = 2;\n" +
			"  }\n" +
			"};\n"))
}

// Java: CodePrinterTest.testMultipleClassStaticBlocks()
func (c *CodePrinterTest) TestMultipleClassStaticBlocks() {
	// empty
	c.AssertPrettyPrintSame("class C {\n  static {\n  }\n  static {\n  }\n}\n")
	c.AssertPrettyPrintSame("let a = class {\n  static {\n  }\n  static {\n  }\n};\n")
	// multiple fields
	c.AssertPrettyPrintSame(
		("class C {\n" +
			"  static field1 = 1;\n" +
			"  static field2 = 2;\n" +
			"  static {\n" +
			"    let x = this.field1;\n" +
			"  }\n" +
			"  static {\n" +
			"    let y = this.field2;\n" +
			"  }\n" +
			"}\n"))
	c.AssertPrettyPrintSame(
		("class C {\n" +
			"  static {\n" +
			"    this.field1 = 1;\n" +
			"  }\n" +
			"  static {\n" +
			"    this.field2 = 2;\n" +
			"  }\n" +
			"}\n"))
	c.AssertPrettyPrintSame(
		("let a = class {\n" +
			"  static field1 = 1;\n" +
			"  static field2 = 2;\n" +
			"  static {\n" +
			"    let x = this.field1;\n" +
			"  }\n" +
			"  static {\n" +
			"    let y = this.field2;\n" +
			"  }\n" +
			"};\n"))
	c.AssertPrettyPrintSame(
		("let a = class {\n" +
			"  static {\n" +
			"    this.field1 = 1;\n" +
			"  }\n" +
			"  static {\n" +
			"    this.field2 = 2;\n" +
			"  }\n" +
			"};\n"))
}

// Java: CodePrinterTest.testDeclarations()
func (c *CodePrinterTest) TestDeclarations() {
	c.AssertPrintSame("let x")
	c.AssertPrintSame("let x,y")
	c.AssertPrintSame("let x=1")
	c.AssertPrintSame("let x=1,y=2")
	c.AssertPrintSame("if(a){let x}")

	c.AssertPrintSame("const x=1")
	c.AssertPrintSame("const x=1,y=2")
	c.AssertPrintSame("if(a){const x=1}")

	c.AssertPrintSame("function f(){}")
	c.AssertPrintSame("if(a){function f(){}}")
	c.AssertPrintSame("if(a)(function(){})")

	c.AssertPrintSame("class f{}")
	c.AssertPrintSame("if(a){class f{}}")
	c.AssertPrintSame("if(a)(class{})")

	c.AssertPrintSame("class P extends(_Component=Component){}")
}

// Java: CodePrinterTest.testImports()
func (c *CodePrinterTest) TestImports() {
	c.diagnosticsToIgnore = []*DiagnosticGroup{DiagnosticGroups_MODULE_LOAD} // allow importing nonexistent modules
	c.AssertPrintSame("import x from\"./foo\"")
	c.AssertPrintSame("import\"./foo\"")
	c.AssertPrintSame("import x,{a as b}from\"./foo\"")
	c.AssertPrintSame("import{a as b,c as d}from\"./foo\"")
	c.AssertPrintSame("import x,{a}from\"./foo\"")
	c.AssertPrintSame("import{a,c}from\"./foo\"")
	c.AssertPrintSame("import x,*as f from\"./foo\"")
	c.AssertPrintSame("import*as f from\"./foo\"")
}

// Java: CodePrinterTest.testExports()
func (c *CodePrinterTest) TestExports() {
	// export declarations
	c.diagnosticsToIgnore = []*DiagnosticGroup{DiagnosticGroups_MODULE_LOAD} // allow importing nonexistent modules
	c.AssertPrintSame("export var x=1")
	c.AssertPrintSame("export var x;export var y")
	c.AssertPrintSame("export let x=1")
	c.AssertPrintSame("export const x=1")
	c.AssertPrintSame("export function f(){}")
	c.AssertPrintSame("export class f{}")
	c.AssertPrintSame("export class f{}export class b{}")

	// export all from
	c.AssertPrint("export * from './a.b.c'", "export*from\"./a.b.c\"")

	// from
	c.AssertPrintSame("export{a}from\"./a.b.c\"")
	c.AssertPrintSame("export{a as x}from\"./a.b.c\"")
	c.AssertPrintSame("export{a,b}from\"./a.b.c\"")
	c.AssertPrintSame("export{a as x,b as y}from\"./a.b.c\"")
	c.AssertPrintSame("export{a}")
	c.AssertPrintSame("export{a as x}")

	c.AssertPrintSame("export{a,b}")
	c.AssertPrintSame("export{a as x,b as y}")

	// export default
	c.AssertPrintSame("export default x")
	c.AssertPrintSame("export default 1")
	c.AssertPrintSame("export default class Foo{}export function f(){}")
	c.AssertPrintSame("export function f(){}export default class Foo{}")
}

// Java: CodePrinterTest.testExportAsyncFunction()
func (c *CodePrinterTest) TestExportAsyncFunction() {
	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT_2017
	c.AssertPrintSame("export async function f(){}")
}

// Java: CodePrinterTest.testTemplateLiteral()
func (c *CodePrinterTest) TestTemplateLiteral() {
	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT_NEXT
	// We need to use the raw string instead of the normalized string in template literals
	c.AssertPrintSame("`hello`")
	c.AssertPrintSame("`\\\\bhello`")
	c.AssertPrintSame("`hel\rlo`")
	c.AssertPrintSame("`hel\\rlo`")
	c.AssertPrintSame("`hel\r\nlo`")
	c.AssertPrintSame("`hel\\r\\nlo`")
	c.AssertPrint("`hello`\n'world'", "`hello`;\"world\"")
	c.AssertPrint("`hello`\n`world`", "`hello``world`")
	c.AssertPrint("var x=`TestA`\n`TemplateB`", "var x=`TestA``TemplateB`")
	c.AssertPrintSame("`hello``world`")

	c.AssertPrintSame("`hello${world}!`")
	c.AssertPrintSame("`hello${world} ${name}!`")

	c.AssertPrintSame("`hello${(function(){let x=3})()}`")
	c.AssertPrintSame("(function(){})()`${(function(){})()}`")
	c.AssertPrintSame("url`hello`")
	c.AssertPrintSame("url(`hello`)")
	c.AssertPrintSame("`\\u{2026}`")
	c.AssertPrintSame("`start\\u{2026}end`")
	c.AssertPrintSame("`\\u{1f42a}`")
	c.AssertPrintSame("`start\\u{1f42a}end`")
	c.AssertPrintSame("`\\u2026`")
	c.AssertPrintSame("`start\\u2026end`")
	c.AssertPrintSame("`\"`")
	c.AssertPrintSame("`'`")
	c.AssertPrintSame("`\\\"`")
	c.AssertPrintSame("`\\'`")
	c.AssertPrintSame("`\\``")

	c.AssertPrintSame("foo`\\unicode`")
	// b/114808380
	c.AssertPrintSame("String.raw`a\\ b`")

	// Nested substitutions.
	c.AssertPrintSame("`Hello ${x?`Alice`:`Bob`}?`")
	c.AssertPrintSame("`Hello ${x?`Alice ${y(`Kitten`)}`:`Bob`}?`")

	// Substitution without padding.
	c.AssertPrintSame("`Unbroken${x}string`")

	// Template strings terminate statements if needed.
	c.AssertPrintSame("let a;`a`")
}

// Java: CodePrinterTest.testMultiLineTemplateLiteral_preservesInteralNewAndBlankLines()
func (c *CodePrinterTest) TestMultiLineTemplateLiteral_preservesInteralNewAndBlankLines() {
	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT_NEXT
	c.AssertPrintSame(
		("var y=`hello // Line break (0 blank lines).\n" +
			"world\n" +
			" // Single blank line.\n" +
			"foo\n" +
			" // Multiple blank lines.\n" +
			"\n" +
			"\n" +
			"bar`\n"))

	c.AssertPrettyPrintSame(
		("var y = `hello // Line break (0 blank lines).\n" +
			"world\n" +
			" // Single blank line.\n" +
			"foo\n" +
			" // Multiple blank lines.\n" +
			"\n" +
			"\n" +
			"bar`;\n"))
}

// Java: CodePrinterTest.testMultiLineTemplateLiteral_doesNotPreserveNewLines_inSubstituions()
func (c *CodePrinterTest) TestMultiLineTemplateLiteral_doesNotPreserveNewLines_inSubstituions() {
	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT_NEXT
	c.AssertPrint(
		("var y=`Hello ${x\n" +
			"+\n" +
			"z\n" +
			"}`\n"),
		"var y=`Hello ${x+z}`")

	c.AssertPrettyPrint(
		("var y=`Hello ${x\n" +
			"+\n" +
			"z\n" +
			"}`\n"),
		"var y = `Hello ${x + z}`;\n")
}

// Java: CodePrinterTest.testMultiLineTemplateLiteral_notIndented_byPrettyPrint()
func (c *CodePrinterTest) TestMultiLineTemplateLiteral_notIndented_byPrettyPrint() {
	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT_NEXT

	// We intentionally put all the delimiter characters on the start of their own line to check
	// their indentation.
	c.AssertPrettyPrint(
		("function indentScope() {\n" +
			"  var y =\n" +
			"// Open backtick.\n" +
			"`hello\n" +
			"world\n" +
			"foo\n" +
			"${ // Open substituion.\n" +
			"bing\n" +
			"// Close substitution.\n" +
			"}\n" +
			"bar\n" +
			"`; // Close backtick.\n" +
			"}\n"),
		("function indentScope() {\n" +
			"  var y = `hello\n" +
			"world\n" +
			"foo\n" +
			"${bing}\n" +
			"bar\n" +
			"`;\n" +
			"}\n"))
}

// Java: CodePrinterTest.testMultiLineTemplateLiteral_brokenOntoLastLine_isNotCollapsed()
func (c *CodePrinterTest) TestMultiLineTemplateLiteral_brokenOntoLastLine_isNotCollapsed() {
	// related to b/117613188

	// Given
	// Configure these so that the printer would otherwise attempt to reuse an existing newline.
	codePrinterOptions := NewCompilerOptions()
	codePrinterOptions.SetLineLengthThreshold(30) // Must be big compared to the last line length.

	input := ("`hello\n" +
		"world\n" +
		"foo\n" +
		"bar`;\n")

	// When
	actual := NewCodePrinter_Builder(c.Parse(input)).
		SetCompilerOptions(codePrinterOptions).
		SetPrettyPrint(false).
		Build()

	// Then
	AssertThat(actual).
		IsEqualTo(
			("`hello\n" +
				"world\n" +
				"foo\n" +
				"bar`"))
}

// Java: CodePrinterTest.testEs6GoogModule()
func (c *CodePrinterTest) TestEs6GoogModule() {
	code := ("goog.module('foo.bar');\n" +
		"const STR = '3';\n" +
		"function fn() {\n" +
		"  alert(STR);\n" +
		"}\n" +
		"exports.fn = fn;\n")
	expectedCode := ("goog.module('foo.bar');\n" +
		"var module$exports$foo$bar = {};\n" +
		"const STR = '3';\n" +
		"function fn() {\n" +
		"  alert(STR);\n" +
		"}\n" +
		"exports.fn = fn;\n")

	compilerOptions := NewCompilerOptions()
	compilerOptions.SetClosurePass(true)
	compilerOptions.SetPreserveDetailedSourceInfo(true)
	compilerOptions.SetContinueAfterErrors(true)
	compiler := NewCompiler()
	compiler.DisableThreads()
	c.CheckWithOriginalName(code, expectedCode, compilerOptions)
}

// Java: CodePrinterTest.testEs6ArrowFunctionSetsOriginalNameForThis()
func (c *CodePrinterTest) TestEs6ArrowFunctionSetsOriginalNameForThis() {
	code := "(x)=>{this.foo[0](3);}"
	// TODO(tomnguyen): Avoid printing the `$jscomp$this$3556498$0 = this` line.
	// TODO(tomnguyen): `function(x) {` should print as an `=>` function.
	expectedCode := ("var $jscomp$this$3556498$0 = this;\n" +
		"(function(x) {\n" +
		"  this.foo[0](3);\n" +
		"});\n")
	compilerOptions := NewCompilerOptions()
	compilerOptions.SkipAllCompilerPasses()
	compilerOptions.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.CheckWithOriginalName(code, expectedCode, compilerOptions)
}

// Java: CodePrinterTest.testEs6ArrowFunctionSetsOriginalNameForArguments()
func (c *CodePrinterTest) TestEs6ArrowFunctionSetsOriginalNameForArguments() {
	// With original names in output set, the end result is not correct code, but the "this" is
	// not rewritten.
	code := "(x)=>{arguments[0]();}"
	expectedCode := ("var $jscomp$arguments$3556498$0 = arguments;\n" +
		"(function(x) {\n" +
		"  arguments[0]();\n" +
		"});\n")
	compilerOptions := NewCompilerOptions()
	compilerOptions.SkipAllCompilerPasses()
	compilerOptions.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.CheckWithOriginalName(code, expectedCode, compilerOptions)
}

// Java: CodePrinterTest.testEs6NewTargetBare()
func (c *CodePrinterTest) TestEs6NewTargetBare() {
	c.AssertPrintSame("class C{constructor(){new.target.prototype}}")
}

// Java: CodePrinterTest.testEs6NewTargetPrototype()
func (c *CodePrinterTest) TestEs6NewTargetPrototype() {
	c.AssertPrintSame(
		"class C{constructor(){var callable=Object.setPrototypeOf(obj,new.target.prototype)}}")
}

// Java: CodePrinterTest.testEs6NewTargetConditional()
func (c *CodePrinterTest) TestEs6NewTargetConditional() {
	c.AssertPrint(
		("function f() {\n" +
			"  if (!new.target) throw 'Must be called with new!';\n" +
			"}\n"),
		"function f(){if(!new.target)throw\"Must be called with new!\";}")
}

// Java: CodePrinterTest.testGoogScope()
func (c *CodePrinterTest) TestGoogScope() {
	// TODO(mknichel): Function declarations need to be rewritten to match the original source
	// instead of being assigned to a local variable with duplicate JS Doc.
	code := ("goog.provide('foo.bar');\n" +
		"goog.require('baz.qux.Quux');\n" +
		"goog.require('foo.ScopedType');\n" +
		"\n" +
		"goog.scope(function() {\n" +
		"var Quux = baz.qux.Quux;\n" +
		"var ScopedType = foo.ScopedType;\n" +
		"\n" +
		"var STR = '3';\n" +
		"/** @param {ScopedType} obj */\n" +
		"function fn(obj) {\n" +
		"  alert(STR);\n" +
		"  alert(Quux.someProperty);\n" +
		"}\n" +
		"}); // goog.scope\n")
	expectedCode := ("goog.provide('foo.bar');\n" +
		"goog.require('baz.qux.Quux');\n" +
		"goog.require('foo.ScopedType');\n" +
		"/**\n" +
		" * @param {ScopedType} obj\n" +
		" */\n" +
		"var $jscomp$scope$3556498$1$fn = /**\n" +
		" * @param {ScopedType} obj\n" +
		" */\n" +
		"function(obj) {\n" +
		"  alert(STR);\n" +
		"  alert(Quux.someProperty);\n" +
		"};\n" +
		"var $jscomp$scope$3556498$0$STR = '3';\n")

	compilerOptions := NewCompilerOptions()
	compilerOptions.SetChecksOnly(true)
	compilerOptions.SetClosurePass(true)
	compilerOptions.SetPreserveDetailedSourceInfo(true)
	compilerOptions.SetCheckTypes(true)
	compilerOptions.SetContinueAfterErrors(true)
	compilerOptions.SetPreserveClosurePrimitives(true)
	compiler := NewCompiler()
	compiler.DisableThreads()
	compiler.Compile(
		[]*SourceFile{}, // Externs
		[]*SourceFile{SourceFile_FromCode("test", code)},
		compilerOptions)
	node := compiler.GetRoot().GetLastChild().GetFirstChild()

	codePrinterOptions := NewCompilerOptions()
	codePrinterOptions.SetPreferSingleQuotes(true)
	codePrinterOptions.SetLineLengthThreshold(80)
	codePrinterOptions.SetPreserveTypeAnnotations(true)
	codePrinterOptions.SetUseOriginalNamesInOutput(true)
	AssertThat(
		NewCodePrinter_Builder(node).
			SetCompilerOptions(codePrinterOptions).
			SetPrettyPrint(true).
			SetLineBreak(true).
			Build()).
		IsEqualTo(expectedCode)
}

// Java: CodePrinterTest.testEscapeDollarInTemplateLiteralInOutput()
func (c *CodePrinterTest) TestEscapeDollarInTemplateLiteralInOutput() {
	compilerOptions := NewCompilerOptions()
	compilerOptions.SkipAllCompilerPasses()
	compilerOptions.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2015)

	c.CheckWithOriginalName(
		"let Foo; const x = `${Foo}`;", "let Foo;\nconst x = `${Foo}`;\n", compilerOptions)

	c.CheckWithOriginalName("const x = `\\${Foo}`;", "const x = `\\${Foo}`;\n", compilerOptions)

	c.CheckWithOriginalName(
		"let Foo; const x = `${Foo}\\${Foo}`;",
		"let Foo;\nconst x = `${Foo}\\${Foo}`;\n",
		compilerOptions)

	c.CheckWithOriginalName(
		"let Foo; const x = `\\${Foo}${Foo}`;",
		"let Foo;\nconst x = `\\${Foo}${Foo}`;\n",
		compilerOptions)
}

// Java: CodePrinterTest.testEscapeDollarInTemplateLiteralEs5Output()
func (c *CodePrinterTest) TestEscapeDollarInTemplateLiteralEs5Output() {
	compilerOptions := NewCompilerOptions()
	compilerOptions.SkipAllCompilerPasses()
	compilerOptions.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)

	c.CheckWithOriginalName(
		"let Foo; const x = `${Foo}`;", "var Foo;\nvar x = '' + Foo;\n", compilerOptions)

	c.CheckWithOriginalName("const x = `\\${Foo}`;", "var x = '${Foo}';\n", compilerOptions)

	c.CheckWithOriginalName(
		"let Foo; const x = `${Foo}\\${Foo}`;",
		"var Foo;\nvar x = Foo + '${Foo}';\n",
		compilerOptions)
	c.CheckWithOriginalName(
		"let Foo; const x = `\\${Foo}${Foo}`;",
		"var Foo;\nvar x = '${Foo}' + Foo;\n",
		compilerOptions)
}

// Java: CodePrinterTest.testDoNotEscapeDollarInRegex()
func (c *CodePrinterTest) TestDoNotEscapeDollarInRegex() {
	compilerOptions := NewCompilerOptions()
	compilerOptions.SkipAllCompilerPasses()
	compilerOptions.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.CheckWithOriginalName("var x = /\\$qux/;", "var x = /\\$qux/;\n", compilerOptions)
	c.CheckWithOriginalName("var x = /$qux/;", "var x = /$qux/;\n", compilerOptions)
}

// Java: CodePrinterTest.testDoNotEscapeDollarInStringLiteral()
func (c *CodePrinterTest) TestDoNotEscapeDollarInStringLiteral() {
	code := "var x = '\\$qux';"
	expectedCode := "var x = '$qux';\n"
	compilerOptions := NewCompilerOptions()
	compilerOptions.SkipAllCompilerPasses()
	compilerOptions.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.CheckWithOriginalName(code, expectedCode, compilerOptions)
	c.CheckWithOriginalName("var x = '\\$qux';", "var x = '$qux';\n", compilerOptions)
	c.CheckWithOriginalName("var x = '$qux';", "var x = '$qux';\n", compilerOptions)
}

// Java: CodePrinterTest.testPrettyPrinterIfElseIfAddedBlock()
func (c *CodePrinterTest) TestPrettyPrinterIfElseIfAddedBlock() {
	c.AssertPrettyPrintSame(
		("if (0) {\n" +
			"  0;\n" +
			"} else if (1) {\n" +
			"  if (2) {\n" +
			"    2;\n" +
			"  }\n" +
			"} else if (3) {\n" +
			"  3;\n" +
			"}\n"))

	c.AssertPrettyPrint(
		"if(0)if(1)1;else 2;else 3;",
		("if (0) {\n" +
			"  if (1) {\n" +
			"    1;\n" +
			"  } else {\n" +
			"    2;\n" +
			"  }\n" +
			"} else {\n" +
			"  3;\n" +
			"}\n"))
}

// Java: CodePrinterTest.testNonJSDocCommentsPrinted_getProp()
func (c *CodePrinterTest) TestNonJSDocCommentsPrinted_getProp() {
	c.preserveNonJSDocComments = true
	// TODO(b/228156705): Fix comment printing properly for GETPROP.
	c.AssertPrettyPrint("a.// testComment\nb", "// testComment\na.b;\n")
}

// Java: CodePrinterTest.testPrintsContentOfClosureUnawareShadow()
func (c *CodePrinterTest) TestPrintsContentOfClosureUnawareShadow() {
	shadowHost := rhino.IR_Name("SHADOW")
	astFactory := NewCompiler().CreateAstFactoryWithoutTypes()
	closureUnawareFn := astFactory.CreateEmptyFunction(AstFactory_Type_OfColor(colors.StandardColors_UNKNOWN))
	shadowHost.SetClosureUnawareShadow(
		rhino.IR_Root(rhino.IR_ScriptOfNodes(rhino.IR_ExprResult(rhino.IR_Call(rhino.IR_Name("PRESERVE"), closureUnawareFn)))))

	AssertThat(
		NewCodePrinter_Builder(shadowHost).SetCompilerOptions(NewCompilerOptions()).Build()).
		IsEqualTo("function(){}")
}

// Java: CodePrinterTest.checkWithOriginalName(String, String, CompilerOptions)
func (c *CodePrinterTest) CheckWithOriginalName(code string, expectedCode string, compilerOptions *CompilerOptions) {
	compilerOptions.SetCheckSymbols(true)
	compilerOptions.SetCheckTypes(true)
	compilerOptions.SetPreserveDetailedSourceInfo(true)
	compilerOptions.SetPreserveClosurePrimitives(true)
	compilerOptions.SetClosurePass(true)
	compiler := NewCompiler()
	compiler.DisableThreads()
	compiler.Compile(
		[]*SourceFile{}, // Externs
		[]*SourceFile{SourceFile_FromCode("test", code)},
		compilerOptions)
	node := compiler.GetRoot().GetLastChild().GetFirstChild()

	codePrinterOptions := NewCompilerOptions()
	codePrinterOptions.SetPreferSingleQuotes(true)
	codePrinterOptions.SetLineLengthThreshold(80)
	codePrinterOptions.SetUseOriginalNamesInOutput(true)
	AssertThat(
		NewCodePrinter_Builder(node).
			SetCompilerOptions(codePrinterOptions).
			SetPrettyPrint(true).
			SetLineBreak(true).
			Build()).
		IsEqualTo(expectedCode)
}

// Java: CodePrinterTest.testNullishCoalesceCastInOr()
func (c *CodePrinterTest) TestNullishCoalesceCastInOr() {
	// TODO(b/571641705): CodeGenerator omits required parentheses around ?? in CAST-wrapped logical expressions
	c.AssertPrintNode(
		"x??y||z",
		rhino.IR_Or(rhino.NewNodeWithChild(rhino.Token_CAST, rhino.IR_Coalesce(rhino.IR_Name("x"), rhino.IR_Name("y"))), rhino.IR_Name("z")))
}

// Java: CodePrinterTest.testNullishCoalesceExtends()
func (c *CodePrinterTest) TestNullishCoalesceExtends() {
	// TODO(b/571641705): CodeGenerator omits required parentheses around ?? in class extends clauses
	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT_NEXT
	AssertThat(c.PrintNode(c.Parse("class C extends (a ?? b) {}"))).IsEqualTo("class C extends a??b{}")
}

// Java: CodePrinterTest.testNullishCoalesceAssignExtends()
func (c *CodePrinterTest) TestNullishCoalesceAssignExtends() {
	// TODO(b/571641705): CodeGenerator omits required parentheses around ??= in class extends clauses
	c.languageMode = CompilerOptions_LanguageMode_ECMASCRIPT_NEXT
	AssertThat(c.PrintNode(c.Parse("class C extends (a ??= b) {}"))).IsEqualTo("class C extends a??=b{}")
}
