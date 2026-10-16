// Port of com/google/javascript/jscomp/CodePrinterEs6TypedTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/CodePrinterEs6TypedTest.java)

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

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CodePrinterEs6TypedTest is a round-trip test for TypeScript-style inline
// type syntax.
//
// Each expression is manually constructed and then code printed, as parsing
// type syntax is not supported.
type CodePrinterEs6TypedTest struct {
	CodePrinterTestBase
}

func NewCodePrinterEs6TypedTest() *CodePrinterEs6TypedTest {
	c := &CodePrinterEs6TypedTest{}
	c.InitCodePrinterTestBase()
	return c
}

func TestCodePrinterEs6TypedTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCodePrinterEs6TypedTest() })
}

// Java: CodePrinterEs6TypedTest.assertPrettyPrint(Node, String)
func (c *CodePrinterEs6TypedTest) AssertPrettyPrint(root *rhino.Node, js string) {
	parsed := NewCodePrinter_Builder(root).
		SetCompilerOptions(
			c.NewCompilerOptions(
				codePrinterTestBase_CompilerOptionBuilderFunc(func(options *CompilerOptions) {
					options.SetPrettyPrint(true)
					options.SetPreferSingleQuotes(true)
					options.SetGentsMode(true)
				}))).
		Build()
	parsed = jutil.StringTrim(parsed) // strip trailing line break.
	AssertThat(parsed).IsEqualTo(js)
}

// Java: CodePrinterEs6TypedTest.testVariableDeclaration()
func (c *CodePrinterEs6TypedTest) TestVariableDeclaration() {
	c.AssertPrettyPrint(CodePrinterEs6TypedTest_ConstructVarDeclarationWithType("any"), "var foo: any = 'hello';")
	c.AssertPrettyPrint(CodePrinterEs6TypedTest_ConstructVarDeclarationWithType("number"), "var foo: number = 'hello';")
	c.AssertPrettyPrint(CodePrinterEs6TypedTest_ConstructVarDeclarationWithType("boolean"), "var foo: boolean = 'hello';")
	c.AssertPrettyPrint(CodePrinterEs6TypedTest_ConstructVarDeclarationWithType("string"), "var foo: string = 'hello';")
	c.AssertPrettyPrint(CodePrinterEs6TypedTest_ConstructVarDeclarationWithType("void"), "var foo: void = 'hello';")
	c.AssertPrettyPrint(CodePrinterEs6TypedTest_ConstructVarDeclarationWithType("hello"), "var foo: hello = 'hello';")
}

// CodePrinterEs6TypedTest_ConstructVarDeclarationWithType constructs
// `var foo: <typeName> = "hello"`.
//
// Java: CodePrinterEs6TypedTest.constructVarDeclarationWithType(String)
func CodePrinterEs6TypedTest_ConstructVarDeclarationWithType(typeName string) *rhino.Node {
	lhs := rhino.IR_Name("foo")
	lhs.SetDeclaredTypeExpression(rhino.TypeDeclarationsIR_NamedType(typeName))
	return rhino.IR_VarWithValue(lhs, rhino.IR_String("hello"))
}

// Java: CodePrinterEs6TypedTest.testFunctionParamDeclaration()
func (c *CodePrinterEs6TypedTest) TestFunctionParamDeclaration() {
	stringType := rhino.TypeDeclarationsIR_NamedType("string")
	xParam := rhino.IR_Name("x")
	xParam.SetDeclaredTypeExpression(stringType)
	fnNode := rhino.IR_Function(rhino.IR_Name("foo"), rhino.IR_ParamList(xParam), rhino.IR_Block())
	c.AssertPrettyPrint(fnNode, "function foo(x: string) {\n}")
}

// Java: CodePrinterEs6TypedTest.testFunctionParamDeclaration_defaultValue()
func (c *CodePrinterEs6TypedTest) TestFunctionParamDeclaration_defaultValue() {
	stringType := rhino.TypeDeclarationsIR_NamedType("string")
	xParam := rhino.IR_Name("x")
	xParam.SetDeclaredTypeExpression(stringType)
	fnNode := rhino.IR_Function(
		rhino.IR_Name("foo"),
		rhino.IR_ParamList(rhino.NewNodeWith2Children(rhino.Token_DEFAULT_VALUE, xParam, rhino.IR_String("hello"))),
		rhino.IR_Block())
	c.AssertPrettyPrint(fnNode, "function foo(x: string = 'hello') {\n}")
}

// Java: CodePrinterEs6TypedTest.testFunctionParamDeclaration_arrow()
func (c *CodePrinterEs6TypedTest) TestFunctionParamDeclaration_arrow() {
	stringType := rhino.TypeDeclarationsIR_NamedType("string")
	param := rhino.IR_Name("x")
	param.SetDeclaredTypeExpression(stringType)
	fnNode := rhino.NewNodeWith3Children(
		rhino.Token_FUNCTION,
		rhino.IR_Name(""),
		rhino.IR_ParamList(param),
		rhino.IR_Add(rhino.IR_String("hello"), rhino.IR_Name("x")))
	fnNode.SetIsArrowFunction(true)
	fnStatement := rhino.IR_ExprResult(fnNode)
	c.AssertPrettyPrint(fnStatement, "(x: string) => 'hello' + x;")
}

// Java: CodePrinterEs6TypedTest.testFunctionReturn()
func (c *CodePrinterEs6TypedTest) TestFunctionReturn() {
	returnType := rhino.TypeDeclarationsIR_NamedType("string")
	fnNode := rhino.IR_Function(rhino.IR_Name("foo"), rhino.IR_ParamList(),
		rhino.IR_BlockOfNode(rhino.IR_ReturnNodeWithExpr(rhino.IR_String("hello"))))
	fnNode.SetDeclaredTypeExpression(returnType)
	c.AssertPrettyPrint(fnNode, "function foo(): string {\n  return 'hello';\n}")
}

// Java: CodePrinterEs6TypedTest.testFunctionReturn_arrow()
func (c *CodePrinterEs6TypedTest) TestFunctionReturn_arrow() {
	returnType := rhino.TypeDeclarationsIR_NamedType("string")
	fnNode := rhino.NewNodeWith3Children(rhino.Token_FUNCTION, rhino.IR_Name(""), rhino.IR_ParamList(), rhino.IR_String("hello"))
	fnNode.SetDeclaredTypeExpression(returnType)
	fnNode.SetIsArrowFunction(true)
	fnStatement := rhino.IR_ExprResult(fnNode)
	c.AssertPrettyPrint(fnStatement, "(): string => 'hello';")
}

// Java: CodePrinterEs6TypedTest.testCompositeType()
func (c *CodePrinterEs6TypedTest) TestCompositeType() {
	type_ := rhino.TypeDeclarationsIR_NamedType("mymod.ns.Type")
	lhs := rhino.IR_Name("foo")
	lhs.SetDeclaredTypeExpression(type_)
	var_ := rhino.IR_Var(lhs)
	c.AssertPrettyPrint(var_, "var foo: mymod.ns.Type;")
}

// Java: CodePrinterEs6TypedTest.testArrayType()
func (c *CodePrinterEs6TypedTest) TestArrayType() {
	type_ := rhino.NewNodeWithChild(rhino.Token_ARRAY_TYPE, rhino.TypeDeclarationsIR_NamedType("string"))
	lhs := rhino.IR_Name("foo")
	lhs.SetDeclaredTypeExpression(type_)
	var_ := rhino.IR_Var(lhs)
	c.AssertPrettyPrint(var_, "var foo: string[];")
}

// Java: CodePrinterEs6TypedTest.testArrayType_qualifiedType()
func (c *CodePrinterEs6TypedTest) TestArrayType_qualifiedType() {
	arrayType := rhino.TypeDeclarationsIR_ArrayType(rhino.TypeDeclarationsIR_NamedType("mymod.ns.Type"))
	lhs := rhino.IR_Name("foo")
	lhs.SetDeclaredTypeExpression(arrayType)
	var_ := rhino.IR_Var(lhs)
	c.AssertPrettyPrint(var_, "var foo: mymod.ns.Type[];")
}

// Java: CodePrinterEs6TypedTest.testParameterizedType()
func (c *CodePrinterEs6TypedTest) TestParameterizedType() {
	type_ := rhino.TypeDeclarationsIR_ParameterizedType(
		rhino.TypeDeclarationsIR_NamedType("my.parameterized.Type"),
		[]*rhino.Node{rhino.TypeDeclarationsIR_NamedType("ns.A"), rhino.TypeDeclarationsIR_NamedType("ns.B")})
	lhs := rhino.IR_Name("x")
	lhs.SetDeclaredTypeExpression(type_)
	var_ := rhino.IR_Var(lhs)
	c.AssertPrettyPrint(var_, "var x: my.parameterized.Type<ns.A, ns.B>;")
}
