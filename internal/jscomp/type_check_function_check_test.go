// Port of com/google/javascript/jscomp/TypeCheckFunctionCheckTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/TypeCheckFunctionCheckTest.java)

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

// TypeCheckFunctionCheckTest tests the type checker.
type TypeCheckFunctionCheckTest struct {
	CompilerTestCaseBase
	convention CodingConvention
}

// Java: TypeCheckFunctionCheckTest()
func NewTypeCheckFunctionCheckTest() *TypeCheckFunctionCheckTest {
	c := &TypeCheckFunctionCheckTest{}
	c.InitCompilerTestCaseBase(c)
	c.convention = nil
	return c
}

func TestTypeCheckFunctionCheckTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewTypeCheckFunctionCheckTest() })
}

// Java: TypeCheckFunctionCheckTest.getProcessor(Compiler)
func (c *TypeCheckFunctionCheckTest) GetProcessor(compiler *Compiler) CompilerPass {
	return CompilerPassFunc(func(externs *rhino.Node, root *rhino.Node) {})
}

// Java: TypeCheckFunctionCheckTest.getCodingConvention()
func (c *TypeCheckFunctionCheckTest) GetCodingConvention() CodingConvention {
	return c.convention
}

// Java: TypeCheckFunctionCheckTest.setUp()
func (c *TypeCheckFunctionCheckTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.convention = NewGoogleCodingConvention()
	c.EnableParseTypeInfo()
	c.EnableTypeCheck()
}

// Java: TypeCheckFunctionCheckTest.testFunctionAritySimple()
func (c *TypeCheckFunctionCheckTest) TestFunctionAritySimple() {
	c.AssertOk("", "")
	c.AssertOk("a", "'a'")
	c.AssertOk("a,b", "10, 20")
}

// Java: TypeCheckFunctionCheckTest.testFunctionArityWithOptionalArgs()
func (c *TypeCheckFunctionCheckTest) TestFunctionArityWithOptionalArgs() {
	c.AssertOk("a,b,opt_c", "1,2")
	c.AssertOk("a,b,opt_c", "1,2,3")
	c.AssertOk("a,opt_b,opt_c", "1")
}

// Java: TypeCheckFunctionCheckTest.testFunctionArityWithVarArgs()
func (c *TypeCheckFunctionCheckTest) TestFunctionArityWithVarArgs() {
	c.AssertOk("var_args", "")
	c.AssertOk("var_args", "1,2")
	c.AssertOk("a,b,var_args", "1,2")
	c.AssertOk("a,b,var_args", "1,2,3")
	c.AssertOk("a,b,var_args", "1,2,3,4,5")
	c.AssertOk("a,opt_b,var_args", "1")
	c.AssertOk("a,opt_b,var_args", "1,2")
	c.AssertOk("a,opt_b,var_args", "1,2,3")
	c.AssertOk("a,opt_b,var_args", "1,2,3,4,5")
}

// Java: TypeCheckFunctionCheckTest.testWrongNumberOfArgs()
func (c *TypeCheckFunctionCheckTest) TestWrongNumberOfArgs() {
	c.AssertWarning("a,b,opt_c", "1", TypeCheck_WRONG_ARGUMENT_COUNT)
	c.AssertWarning("a,b,opt_c", "1,2,3,4", TypeCheck_WRONG_ARGUMENT_COUNT)
	c.AssertWarning("a,b", "1, 2, 3", TypeCheck_WRONG_ARGUMENT_COUNT)
	c.AssertWarning("", "1, 2, 3", TypeCheck_WRONG_ARGUMENT_COUNT)
	c.AssertWarning("a,b,c,d", "1, 2, 3", TypeCheck_WRONG_ARGUMENT_COUNT)
	c.AssertWarning("a,b,var_args", "1", TypeCheck_WRONG_ARGUMENT_COUNT)
	c.AssertWarning("a,b,opt_c,var_args", "1", TypeCheck_WRONG_ARGUMENT_COUNT)
}

// Java: TypeCheckFunctionCheckTest.testVarArgsLast()
func (c *TypeCheckFunctionCheckTest) TestVarArgsLast() {
	c.AssertWarning("a,b,var_args,c", "1,2,3,4", FunctionTypeBuilder_VAR_ARGS_MUST_BE_LAST)
}

// Java: TypeCheckFunctionCheckTest.testOptArgsLast()
func (c *TypeCheckFunctionCheckTest) TestOptArgsLast() {
	c.AssertWarning("a,b,opt_d,c", "1, 2, 3", FunctionTypeBuilder_OPTIONAL_ARG_AT_END)
	c.AssertWarning("a,b,opt_d,c", "1, 2", FunctionTypeBuilder_OPTIONAL_ARG_AT_END)
}

// Java: TypeCheckFunctionCheckTest.testFunctionsWithJsDoc1()
func (c *TypeCheckFunctionCheckTest) TestFunctionsWithJsDoc1() {
	c.TestSame("/** @param {*=} c */ function foo(a,b,c) {} foo(1,2);")
}

// Java: TypeCheckFunctionCheckTest.testFunctionsWithJsDoc2()
func (c *TypeCheckFunctionCheckTest) TestFunctionsWithJsDoc2() {
	c.TestSame("/** @param {*=} c */ function foo(a,b,c) {} foo(1,2,3);")
}

// Java: TypeCheckFunctionCheckTest.testFunctionsWithJsDoc3()
func (c *TypeCheckFunctionCheckTest) TestFunctionsWithJsDoc3() {
	c.TestSame(`/** @param {*=} c 
 * @param {*=} b */
function foo(a,b,c) {} foo(1);
`)
}

// Java: TypeCheckFunctionCheckTest.testFunctionsWithJsDoc4()
func (c *TypeCheckFunctionCheckTest) TestFunctionsWithJsDoc4() {
	c.TestSame("/** @param {...*} a */ var foo = function(a) {}; foo();")
}

// Java: TypeCheckFunctionCheckTest.testFunctionsWithJsDoc5()
func (c *TypeCheckFunctionCheckTest) TestFunctionsWithJsDoc5() {
	c.TestSame("/** @param {...*} a */ var foo = function(a) {}; foo(1,2);")
}

// Java: TypeCheckFunctionCheckTest.testFunctionsWithJsDoc6()
func (c *TypeCheckFunctionCheckTest) TestFunctionsWithJsDoc6() {
	c.TestWarning("/** @param {...*} b */ var foo = function(a, b) {}; foo();", TypeCheck_WRONG_ARGUMENT_COUNT)
}

// Java: TypeCheckFunctionCheckTest.testFunctionsWithJsDoc7()
func (c *TypeCheckFunctionCheckTest) TestFunctionsWithJsDoc7() {
	fooDfn := "/** @param {*} [b] */ var foo = function(b) {};"
	c.TestSame(fooDfn + "foo();")
	c.TestSame(fooDfn + "foo(1);")
	c.TestWarning(fooDfn+"foo(1, 2);", TypeCheck_WRONG_ARGUMENT_COUNT)
}

// Java: TypeCheckFunctionCheckTest.testFunctionWithDefaultCodingConvention()
func (c *TypeCheckFunctionCheckTest) TestFunctionWithDefaultCodingConvention() {
	c.convention = CodingConventions_GetDefault()
	c.TestWarning("var foo = function(x) {}; foo(1, 2);", TypeCheck_WRONG_ARGUMENT_COUNT)
	c.TestWarning("var foo = function(opt_x) {}; foo(1, 2);", TypeCheck_WRONG_ARGUMENT_COUNT)
	c.TestWarning("var foo = function(var_args) {}; foo(1, 2);", TypeCheck_WRONG_ARGUMENT_COUNT)
}

// Java: TypeCheckFunctionCheckTest.testMethodCalls()
func (c *TypeCheckFunctionCheckTest) TestMethodCalls() {
	METHOD_DEFS := `/** @constructor */
function Foo() {}
// Methods defined in a separate functions and then added via assignment
function twoArg(arg1, arg2) {};
Foo.prototype.prototypeMethod = twoArg;
Foo.staticMethod = twoArg;
// Constructor that specifies a return type
/**
 * @constructor
 * @return {Bar}
 */
function Bar() {}
`
	c.TestWarning(METHOD_DEFS+"var f = new Foo();f.prototypeMethod(1, 2, 3);", TypeCheck_WRONG_ARGUMENT_COUNT)
	c.TestWarning(METHOD_DEFS+"var f = new Foo();f.prototypeMethod(1);", TypeCheck_WRONG_ARGUMENT_COUNT)
	c.TestWarning(METHOD_DEFS+"Foo.staticMethod(1, 2, 3);", TypeCheck_WRONG_ARGUMENT_COUNT)
	c.TestWarning(METHOD_DEFS+"Foo.staticMethod(1);", TypeCheck_WRONG_ARGUMENT_COUNT)
	c.TestWarning(METHOD_DEFS+"Foo();", TypeCheck_CONSTRUCTOR_NOT_CALLABLE)
	c.TestSame(METHOD_DEFS + "Bar();")
	c.Test(c.Externs(METHOD_DEFS), c.Srcs("Foo();"), c.Warning(TypeCheck_CONSTRUCTOR_NOT_CALLABLE))
	c.TestSame(c.Externs(METHOD_DEFS), c.Srcs("Bar();"))
}

// Java: TypeCheckFunctionCheckTest.assertOk(String, String)
func (c *TypeCheckFunctionCheckTest) AssertOk(params string, arguments string) {
	c.TestSame("function foo(" + params + ") {} foo(" + arguments + ");")
}

// Java: TypeCheckFunctionCheckTest.assertWarning(String, String, DiagnosticType)
func (c *TypeCheckFunctionCheckTest) AssertWarning(params string, arguments string, type_ *DiagnosticType) {
	c.TestWarning("function foo("+params+") {} foo("+arguments+");", type_)
}
