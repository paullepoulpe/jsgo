// Port of com/google/javascript/jscomp/FunctionTypeBuilderTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/FunctionTypeBuilderTest.java)

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
	"github.com/paullepoulpe/jsgo/internal/rhino/rhinotesting"
)

// FunctionTypeBuilderTest holds the unit tests for FunctionTypeBuilder.
type FunctionTypeBuilderTest struct {
	CompilerTestCaseBase
}

// Java: FunctionTypeBuilderTest()
func NewFunctionTypeBuilderTest() *FunctionTypeBuilderTest {
	c := &FunctionTypeBuilderTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestFunctionTypeBuilderTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewFunctionTypeBuilderTest() })
}

// Java: FunctionTypeBuilderTest.setUp()
func (c *FunctionTypeBuilderTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableParseTypeInfo()
	c.EnableTypeCheck()
}

// functionTypeBuilderTest_noopPass is the anonymous CompilerPass of
// getProcessor.
type functionTypeBuilderTest_noopPass struct {
	_ byte // Go-only: language-adapter — distinct identity
}

func (p *functionTypeBuilderTest_noopPass) Process(externs rhino.Node, js rhino.Node) {}

// Java: FunctionTypeBuilderTest.getProcessor(Compiler)
func (c *FunctionTypeBuilderTest) GetProcessor(compiler *Compiler) CompilerPass {
	// By turning on type checking, the FunctionTypeBuilder will be invoked.
	return &functionTypeBuilderTest_noopPass{}
}

// Java: FunctionTypeBuilderTest.testValidBuiltInTypeRedefinition()
func (c *FunctionTypeBuilderTest) TestValidBuiltInTypeRedefinition() {
	c.TestSame(c.Externs(rhinotesting.BaseJSTypeTestCase_ALL_NATIVE_EXTERN_TYPES), c.Srcs(""))
}

// Java: FunctionTypeBuilderTest.testBuiltInTypeDifferentReturnType()
func (c *FunctionTypeBuilderTest) TestBuiltInTypeDifferentReturnType() {
	c.Test(
		c.Externs(
			"/**\n"+
				" * @constructor\n"+
				" * @param {*} opt_str\n"+
				" * @return {number}\n"+
				" */\n"+
				"function String(opt_str) {}\n"),
		c.Srcs(""),
		c.Warning(FunctionTypeBuilder_TYPE_REDEFINITION).
			WithMessage(
				"attempted re-definition of type String\n"+
					"found   : (typeof String)\n"+
					"expected: (typeof String)\n"))
}

// Java: FunctionTypeBuilderTest.testBuiltInTypeDifferentNumParams()
func (c *FunctionTypeBuilderTest) TestBuiltInTypeDifferentNumParams() {
	c.Test(
		c.Externs(
			"/**\n"+
				" * @constructor\n"+
				" * @return {string}\n"+
				" */\n"+
				"function String() {}\n"),
		c.Srcs(""),
		c.Warning(FunctionTypeBuilder_TYPE_REDEFINITION).
			WithMessage(
				"attempted re-definition of type String\n"+
					"found   : (typeof String)\n"+
					"expected: (typeof String)\n"))
}

// Java: FunctionTypeBuilderTest.testBuiltInTypeDifferentNumParams2()
func (c *FunctionTypeBuilderTest) TestBuiltInTypeDifferentNumParams2() {
	c.Test(
		c.Externs(
			"/**\n"+
				" * @constructor\n"+
				" * @return {string}\n"+
				" */\n"+
				"function String(opt_str, opt_nothing) {}\n"),
		c.Srcs(""),
		c.Warning(FunctionTypeBuilder_TYPE_REDEFINITION).
			WithMessage(
				"attempted re-definition of type String\n"+
					"found   : (typeof String)\n"+
					"expected: (typeof String)\n"))
}

// Java: FunctionTypeBuilderTest.testBuiltInTypeDifferentParamType()
func (c *FunctionTypeBuilderTest) TestBuiltInTypeDifferentParamType() {
	c.Test(
		c.Externs(
			"/**\n"+
				" * @constructor\n"+
				" * @return {string}\n"+
				" */\n"+
				"function String(opt_str) {}\n"),
		c.Srcs(""),
		c.Warning(FunctionTypeBuilder_TYPE_REDEFINITION).
			WithMessage(
				"attempted re-definition of type String\n"+
					"found   : (typeof String)\n"+
					"expected: (typeof String)\n"))
}

// Java: FunctionTypeBuilderTest.testBadFunctionTypeDefinition()
func (c *FunctionTypeBuilderTest) TestBadFunctionTypeDefinition() {
	c.Test(
		c.Externs("/** @constructor */function Function(opt_str) {}\n"),
		c.Srcs(""),
		c.Warning(FunctionTypeBuilder_TYPE_REDEFINITION).
			WithMessage(
				"attempted re-definition of type Function\n"+
					"found   : (typeof Function)\n"+
					"expected: (typeof Function)\n"))
}

// Java: FunctionTypeBuilderTest.testInlineJsDoc()
func (c *FunctionTypeBuilderTest) TestInlineJsDoc() {
	c.Test(
		c.Externs("/** @return {number} */ function f(/** string */ x) { return x; }"),
		c.Srcs(""),
		c.Warning(TypeValidator_TYPE_MISMATCH_WARNING).
			WithMessage(
				"inconsistent return type\n"+
					"found   : string\n"+
					"required: number\n"))
}

// Java: FunctionTypeBuilderTest.testInlineJsDoc2()
func (c *FunctionTypeBuilderTest) TestInlineJsDoc2() {
	c.Test(
		c.Externs(
			"/** @return {T} \n"+
				" @template T */\n"+
				"function f(/** T */ x) { return x; }\n"+
				"/** @type {string} */ var x = f(1);\n"),
		c.Srcs(""),
		c.Warning(TypeValidator_TYPE_MISMATCH_WARNING).
			WithMessage(
				"initializing variable\n"+
					"found   : number\n"+
					"required: string\n"))
}
