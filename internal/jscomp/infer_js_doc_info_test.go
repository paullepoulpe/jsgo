// Port of com/google/javascript/jscomp/InferJSDocInfoTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/InferJSDocInfoTest.java)

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
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Java: private static enum InferJSDocInfoTest.DeclarationKeyword
type inferJSDocInfoTest_DeclarationKeyword int

const (
	inferJSDocInfoTest_DeclarationKeyword_VAR inferJSDocInfoTest_DeclarationKeyword = iota
	inferJSDocInfoTest_DeclarationKeyword_LET
	inferJSDocInfoTest_DeclarationKeyword_CONST
)

// Java: InferJSDocInfoTest.DeclarationKeyword.toJs()
func (k inferJSDocInfoTest_DeclarationKeyword) toJs() string {
	switch k {
	case inferJSDocInfoTest_DeclarationKeyword_VAR:
		return "var"
	case inferJSDocInfoTest_DeclarationKeyword_LET:
		return "let"
	default:
		return "const"
	}
}

// InferJSDocInfoTest holds tests for InferJSDocInfo.
//
// Java: public final class InferJSDocInfoTest extends CompilerTestCase
type InferJSDocInfoTest struct {
	CompilerTestCaseBase
}

// Java: InferJSDocInfoTest()
func NewInferJSDocInfoTest() *InferJSDocInfoTest {
	c := &InferJSDocInfoTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestInferJSDocInfoTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewInferJSDocInfoTest() })
}

// Java: InferJSDocInfoTest.getOptions()
func (c *InferJSDocInfoTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetParseJsDocDocumentation(parsing.Config_JsDocParsing_INCLUDE_DESCRIPTIONS_NO_WHITESPACE)
	return options
}

// Java: InferJSDocInfoTest.getProcessor(Compiler)
func (c *InferJSDocInfoTest) GetProcessor(compiler *Compiler) CompilerPass {
	return CompilerPassFunc(func(externs *rhino.Node, root *rhino.Node) {
		scopeCreator := NewTypedScopeCreator(compiler)

		NewTypeInferencePass(compiler, compiler.GetReverseAbstractInterpreter(), scopeCreator).
			InferAllScopes(root.GetParent())

		NewInferJSDocInfo(compiler).Process(externs, root)
	})
}

// Java: InferJSDocInfoTest.testJSDocFromExternTypesIsPreserved()
func (c *InferJSDocInfoTest) TestJSDocFromExternTypesIsPreserved() {
	// Given
	c.TestSame(
		c.Externs(
			"/**\n"+
				" * I'm an Object.\n"+
				" * @param {*=} x\n"+
				" * @return {!Object}\n"+
				" * @constructor\n"+
				" */\n"+
				"function Object(x) {};\n"),
		c.Srcs("var x = new Object();"))

	xType := c.inferredTypeOfName("x")
	AssertThat(xType.String()).IsEqualTo("Object")

	// Then
	AssertThat(xType.GetJSDocInfo().GetBlockDescription()).IsEqualTo("I'm an Object.")
}

// Java: InferJSDocInfoTest.testJSDocFromInstanceNodesIsIgnored()
func (c *InferJSDocInfoTest) TestJSDocFromInstanceNodesIsIgnored() {
	// Given
	c.TestSame(
		c.Srcs(
			"/**\n" +
				" * I'm a user type.\n" +
				" * @constructor\n" +
				" */\n" +
				"function Foo() {};\n" +
				"\n" +
				"/** I'm a Foo instance */ // This should not be attached to a type.\n" +
				"var x = new Foo();\n"))

	xType := c.inferredTypeOfName("x")
	AssertThat(xType.String()).IsEqualTo("Foo")

	// Then
	AssertThat(xType.GetJSDocInfo().GetBlockDescription()).IsEqualTo("I'm a user type.")
}

// Java: InferJSDocInfoTest.testJSDocIsNotPropagatedToNativeFunctionType()
func (c *InferJSDocInfoTest) TestJSDocIsNotPropagatedToNativeFunctionType() {
	// Given
	c.TestSame(
		c.Srcs(
			"/**\n" +
				"* I'm some function.\n" +
				" * @type {!Function}\n" +
				" */\n" +
				"var x;\n"))

	xType := c.inferredTypeOfName("x")
	AssertThat(xType.String()).IsEqualTo("Function")

	// Then
	AssertThat(xType.GetJSDocInfo()).IsNull()
}

// Java: InferJSDocInfoTest.testJSDocFromNamedFunctionPropagatesToDefinedType()
func (c *InferJSDocInfoTest) TestJSDocFromNamedFunctionPropagatesToDefinedType() {
	// Given
	c.TestSame(
		c.Srcs(
			"/**\n" +
				" * I'm a user type.\n" +
				" * @constructor\n" +
				" */\n" +
				"function Foo() {};\n" +
				"\n" +
				"var x = new Foo();\n"))

	xType := c.inferredTypeOfName("x")
	AssertThat(xType.String()).IsEqualTo("Foo")

	// Then
	AssertThat(xType.GetJSDocInfo().GetBlockDescription()).IsEqualTo("I'm a user type.")
}

// Java: InferJSDocInfoTest.testJSDocFromNamedEs6ClassPropagatesToDefinedType()
func (c *InferJSDocInfoTest) TestJSDocFromNamedEs6ClassPropagatesToDefinedType() {
	// Given
	c.TestSame(
		c.Srcs(
			"/** I'm a user type. */\n" +
				"class Foo { };\n" +
				"\n" +
				"var x = new Foo();\n"))

	xType := c.inferredTypeOfName("x")
	AssertThat(xType.String()).IsEqualTo("Foo")

	// Then
	AssertThat(xType.GetJSDocInfo().GetBlockDescription()).IsEqualTo("I'm a user type.")
}

// Java: InferJSDocInfoTest.testJSDocFromGlobalAssignmentPropagatesToDefinedType()
func (c *InferJSDocInfoTest) TestJSDocFromGlobalAssignmentPropagatesToDefinedType() {
	// Given
	c.TestSame(
		c.Srcs(
			"/**\n" +
				" * I'm a user type.\n" +
				" * @constructor\n" +
				" */\n" +
				"Foo = function() {};\n" +
				"\n" +
				"var x = new Foo();\n"))

	xType := c.inferredTypeOfName("x")
	AssertThat(xType.String()).IsEqualTo("Foo")

	// Then
	AssertThat(xType.GetJSDocInfo().GetBlockDescription()).IsEqualTo("I'm a user type.")
}

// Java: InferJSDocInfoTest.testJSDocFromNamespacedNameAssignmentPropagatesToDefinedType()
func (c *InferJSDocInfoTest) TestJSDocFromNamespacedNameAssignmentPropagatesToDefinedType() {
	// Given
	c.TestSame(
		c.Srcs(
			"var namespace = {}\n" +
				"\n" +
				"/**\n" +
				" * I'm a user type.\n" +
				" * @constructor\n" +
				" */\n" +
				"namespace.Foo = function() {};\n" +
				"\n" +
				"var x = new namespace.Foo();\n"))

	xType := c.inferredTypeOfName("x")
	AssertThat(xType.String()).IsEqualTo("namespace.Foo")

	// Then
	AssertThat(xType.GetJSDocInfo().GetBlockDescription()).IsEqualTo("I'm a user type.")
}

// Java: InferJSDocInfoTest.testJSDocFromVarAssignmentPropagatesToDefinedType()
func (c *InferJSDocInfoTest) TestJSDocFromVarAssignmentPropagatesToDefinedType() {
	c.testJSDocFromVariableAssignmentPropagatesToDefinedType(inferJSDocInfoTest_DeclarationKeyword_VAR)
}

// Java: InferJSDocInfoTest.testJSDocFromLetAssignmentPropagatesToDefinedType()
func (c *InferJSDocInfoTest) TestJSDocFromLetAssignmentPropagatesToDefinedType() {
	c.testJSDocFromVariableAssignmentPropagatesToDefinedType(inferJSDocInfoTest_DeclarationKeyword_LET)
}

// Java: InferJSDocInfoTest.testJSDocFromConstAssignmentPropagatesToDefinedType()
func (c *InferJSDocInfoTest) TestJSDocFromConstAssignmentPropagatesToDefinedType() {
	c.testJSDocFromVariableAssignmentPropagatesToDefinedType(inferJSDocInfoTest_DeclarationKeyword_CONST)
}

// Java: InferJSDocInfoTest.testJSDocFromVariableAssignmentPropagatesToDefinedType(DeclarationKeyword keyword)
func (c *InferJSDocInfoTest) testJSDocFromVariableAssignmentPropagatesToDefinedType(keyword inferJSDocInfoTest_DeclarationKeyword) {
	// Given
	c.TestSame(
		c.Srcs(
			strings.ReplaceAll(
				"/**\n"+
					" * I'm a user type.\n"+
					" * @constructor\n"+
					" */\n"+
					"KEYWORD Foo = function() {};\n"+
					"\n"+
					"// `var x` is just a hook to access type \"Foo\".\n"+
					"var x = new Foo();\n",
				"KEYWORD", keyword.toJs())))

	xType := c.inferredTypeOfName("x")
	AssertThat(xType.String()).IsEqualTo("Foo")

	// Then
	AssertThat(xType.GetJSDocInfo().GetBlockDescription()).IsEqualTo("I'm a user type.")
}

// Java: InferJSDocInfoTest.testJSDocFromVariableNameAssignmentPropagatesToDefinedType()
func (c *InferJSDocInfoTest) TestJSDocFromVariableNameAssignmentPropagatesToDefinedType() {
	// Given
	c.TestSame(
		c.Srcs(
			"var /**\n" +
				" * I'm a user type.\n" +
				" * @constructor\n" +
				" */\n" +
				"Foo = function() {};\n" +
				"\n" +
				"var x = new Foo();\n"))

	xType := c.inferredTypeOfName("x")
	AssertThat(xType.String()).IsEqualTo("Foo")

	// Then
	AssertThat(xType.GetJSDocInfo().GetBlockDescription()).IsEqualTo("I'm a user type.")
}

// Java: InferJSDocInfoTest.testJSDocIsPropagatedToClasses_Es5()
func (c *InferJSDocInfoTest) TestJSDocIsPropagatedToClasses_Es5() {
	// Given
	c.TestSame(
		c.Srcs(
			"/**\n" +
				" * I'm a user class.\n" +
				" * @constructor\n" +
				" */\n" +
				"var Foo = function() {};\n" +
				"\n" +
				"var x = new Foo();\n"))

	xType := c.inferredTypeOfName("x")
	AssertThat(xType.String()).IsEqualTo("Foo")

	// Then
	AssertThat(xType.GetJSDocInfo().GetBlockDescription()).IsEqualTo("I'm a user class.")
}

// Java: InferJSDocInfoTest.testJSDocIsPropagatedToClasses_Es6()
func (c *InferJSDocInfoTest) TestJSDocIsPropagatedToClasses_Es6() {
	// Given
	c.TestSame(
		c.Srcs(
			"/** I'm a user class. */\n" +
				"var Foo = class { };\n" +
				"\n" +
				"var x = new Foo();\n"))

	xType := c.inferredTypeOfName("x")
	AssertThat(xType.String()).IsEqualTo("Foo")

	// Then
	AssertThat(xType.GetJSDocInfo().GetBlockDescription()).IsEqualTo("I'm a user class.")
}

// Java: InferJSDocInfoTest.testJSDocIsPropagatedToCtors()
func (c *InferJSDocInfoTest) TestJSDocIsPropagatedToCtors() {
	// Given
	c.TestSame(
		c.Srcs(
			"/**\n" +
				" * I'm a user class.\n" +
				" * @constructor\n" +
				" */\n" +
				"var Foo = function() {};\n" +
				"\n" +
				"var x = Foo;\n"))

	xType := c.inferredTypeOfName("x")
	AssertThat(xType.String()).IsEqualTo("(typeof Foo)")

	// Then
	AssertThat(xType.GetJSDocInfo().GetBlockDescription()).IsEqualTo("I'm a user class.")
}

// Java: InferJSDocInfoTest.testJSDocIsPropagatedToInterfaces()
func (c *InferJSDocInfoTest) TestJSDocIsPropagatedToInterfaces() {
	// Given
	c.TestSame(
		c.Srcs(
			"/**\n" +
				" * I'm a user interface.\n" +
				" * @interface\n" +
				" */\n" +
				"var Foo = function() {};\n" +
				"\n" +
				"var x = /** @type {!Foo} */ ({});\n"))

	xType := c.inferredTypeOfName("x")
	AssertThat(xType.String()).IsEqualTo("Foo")

	// Then
	AssertThat(xType.GetJSDocInfo().GetBlockDescription()).IsEqualTo("I'm a user interface.")
}

// Java: InferJSDocInfoTest.testJSDocIsPropagatedToRecords()
func (c *InferJSDocInfoTest) TestJSDocIsPropagatedToRecords() {
	// Given
	c.TestSame(
		c.Srcs(
			"/**\n" +
				" * I'm a user record.\n" +
				" * @record\n" +
				" */\n" +
				"var Foo = function() {};\n" +
				"\n" +
				"var x = /** @type {!Foo} */ ({});\n"))

	xType := c.inferredTypeOfName("x")
	AssertThat(xType.String()).IsEqualTo("Foo")

	// Then
	AssertThat(xType.GetJSDocInfo().GetBlockDescription()).IsEqualTo("I'm a user record.")
}

// Java: InferJSDocInfoTest.testJSDocIsPropagatedToEnums()
func (c *InferJSDocInfoTest) TestJSDocIsPropagatedToEnums() {
	// Given
	c.TestSame(
		c.Srcs(
			"/**\n" +
				" * I'm a user enum.\n" +
				" * @enum {number}\n" +
				" */\n" +
				"var Foo = {BAR: 0};\n" +
				"\n" +
				"var x = Foo;\n"))

	xType := c.inferredTypeOfName("x")
	AssertThat(xType.String()).IsEqualTo("enum{Foo}")

	// Then
	AssertThat(xType.GetJSDocInfo().GetBlockDescription()).IsEqualTo("I'm a user enum.")
}

// Java: InferJSDocInfoTest.testJSDocIsPropagatedToEnumElements()
func (c *InferJSDocInfoTest) TestJSDocIsPropagatedToEnumElements() {
	// Given
	c.TestSame(
		c.Srcs(
			"/**\n" +
				" * I'm a user enum.\n" +
				" * @enum {number}\n" +
				" */\n" +
				"var Foo = {BAR: 0};\n" +
				"\n" +
				"var x = Foo.BAR;\n"))

	xType := c.inferredTypeOfName("x")
	AssertThat(xType.String()).IsEqualTo("Foo<number>")

	// Then
	AssertThat(xType.GetJSDocInfo().GetBlockDescription()).IsEqualTo("I'm a user enum.")
}

// Java: InferJSDocInfoTest.testJSDocIsPropagatedToFunctionTypes()
func (c *InferJSDocInfoTest) TestJSDocIsPropagatedToFunctionTypes() {
	// Given
	c.TestSame(
		c.Srcs(
			"/**\n" +
				" * I'm a custom function.\n" +
				" * @param {*} a\n" +
				" * @return {string}\n" +
				" */\n" +
				"function test(a) {};\n" +
				"\n" +
				"var x = test;\n"))

	xType := c.inferredTypeOfName("x")
	AssertThat(xType.String()).IsEqualTo("function(*): string")

	// Then
	AssertThat(xType.GetJSDocInfo().GetBlockDescription()).IsEqualTo("I'm a custom function.")
}

// Java: InferJSDocInfoTest.testJSDocIsPropagatedToScopedTypes()
func (c *InferJSDocInfoTest) TestJSDocIsPropagatedToScopedTypes() {
	// Given
	c.TestSame(
		c.Srcs(
			"(() => {\n" +
				"  /**\n" +
				"   * I'm a scoped user class.\n" +
				"   * @constructor\n" +
				"   */\n" +
				"  var Foo = function() {};\n" +
				"\n" +
				"  var x = new Foo(); // Just a hook to access type \"Foo\".\n" +
				"})();\n"))

	xType := c.inferredTypeOfName("x")
	AssertThat(xType.String()).IsEqualTo("Foo")

	// Then
	AssertThat(xType.GetJSDocInfo().GetBlockDescription()).IsEqualTo("I'm a scoped user class.")
}

// Java: InferJSDocInfoTest.testJSDocIsPropagatedToFunctionTypesFromMethodAssigments()
func (c *InferJSDocInfoTest) TestJSDocIsPropagatedToFunctionTypesFromMethodAssigments() {
	// Given
	c.TestSame(
		c.Srcs(
			"/** @constructor */\n" +
				"function Foo() {};\n" +
				"\n" +
				"/**\n" +
				" * I'm a method.\n" +
				" * @return {number} a\n" +
				" */\n" +
				"Foo.prototype.method = function() { };\n" +
				"\n" +
				"var x = Foo.prototype.method;\n"))

	xType := c.inferredTypeOfName("x")
	AssertThat(xType.String()).IsEqualTo("function(this:Foo): number")

	// Then
	AssertThat(xType.GetJSDocInfo().GetBlockDescription()).IsEqualTo("I'm a method.")
}

// Java: InferJSDocInfoTest.testJSDocIsNotPropagatedToFunctionTypesFromConstAliases()
func (c *InferJSDocInfoTest) TestJSDocIsNotPropagatedToFunctionTypesFromConstAliases() {
	// Given
	c.TestSame(
		c.Srcs(
			"var foo = function() {};\n" +
				"\n" +
				"/**\n" +
				" * I'm an alias\n" +
				" * @const\n" +
				" */\n" +
				"var x = foo;\n"))

	xType := c.inferredTypeOfName("x")
	AssertThat(xType.String()).IsEqualTo("function(): undefined")

	// Then
	AssertThat(xType.GetJSDocInfo()).IsNull()
}

// Java: InferJSDocInfoTest.testJSDocIsNotPropagatedToFunctionTypesFromPropertyAliases()
func (c *InferJSDocInfoTest) TestJSDocIsNotPropagatedToFunctionTypesFromPropertyAliases() {
	// Given
	c.TestSame(
		c.Srcs(
			"/**\n" +
				" * @constructor\n" +
				" * @param {function(*): string} callback\n" +
				" */\n" +
				"function Foo(callback) {\n" +
				"\n" +
				"  /**\n" +
				"   * CB1\n" +
				"   * @private @const\n" +
				"   */\n" +
				"  this.callback1 = callback;\n" +
				"  /**\n" +
				"   * CB2\n" +
				"   * @private @const\n" +
				"   */\n" +
				"  this.callback2 = callback;\n" +
				"  /**\n" +
				"   * CB3\n" +
				"   * @private @const\n" +
				"   */\n" +
				"  this.callback3 = callback;\n" +
				"\n" +
				"};\n"))

	callbackType := c.inferredTypeOfName("callback")
	fooInstanceType := c.inferredTypeOfName("Foo").ToMaybeFunctionType().GetInstanceType()
	AssertThat(callbackType.String()).IsEqualTo("function(*): string")

	// Then
	AssertThat(callbackType.GetJSDocInfo()).IsNull()
	AssertThat(fooInstanceType.GetPropertyJSDocInfo("callback1").GetBlockDescription()).
		IsEqualTo("CB1")
	AssertThat(fooInstanceType.GetPropertyJSDocInfo("callback2").GetBlockDescription()).
		IsEqualTo("CB2")
	AssertThat(fooInstanceType.GetPropertyJSDocInfo("callback3").GetBlockDescription()).
		IsEqualTo("CB3")
}

// Java: InferJSDocInfoTest.testJSDocItPropagatedToFunctionTypesFromPropertyAliases_withUniqueTypeDeclarations()
func (c *InferJSDocInfoTest) TestJSDocItPropagatedToFunctionTypesFromPropertyAliases_withUniqueTypeDeclarations() {
	// Given
	c.TestSame(
		c.Srcs(
			"/**\n" +
				" * @constructor\n" +
				" * @param {function(*): string} callback\n" +
				" */\n" +
				"function Foo(callback) {\n" +
				"\n" +
				"  /**\n" +
				"   * CB1\n" +
				"   * @private @const {function(string): string}\n" +
				"   */\n" +
				"  this.callback = callback;\n" +
				"\n" +
				"};\n"))

	callbackType := c.inferredTypeOfName("callback")
	fooInstanceType := c.inferredTypeOfName("Foo").ToMaybeFunctionType().GetInstanceType()
	thisDotCallbackType :=
		fooInstanceType.GetPropertyType("callback").ToMaybeFunctionType()

	AssertThat(callbackType.String()).IsEqualTo("function(*): string")
	AssertThat(thisDotCallbackType.String()).IsEqualTo("function(string): string")

	// Then
	// While `this.callback = callback` doesn't attach any JSDocInfo to the parameter
	// type `function(*): string`, it does attach JSDocInfo to the narrower function type
	// `function(string): string`.
	AssertThat(callbackType.GetJSDocInfo()).IsNull()
	AssertThat(fooInstanceType.GetPropertyJSDocInfo("callback").GetBlockDescription()).
		IsEqualTo("CB1")
	AssertThat(thisDotCallbackType.GetJSDocInfo().GetBlockDescription()).IsEqualTo("CB1")
}

// Java: InferJSDocInfoTest.testJSDocIsPropagatedDistinctlyToMatchingStructuralTypes_ObjectLiteralTypes()
func (c *InferJSDocInfoTest) TestJSDocIsPropagatedDistinctlyToMatchingStructuralTypes_ObjectLiteralTypes() {
	// Given
	c.TestSame(
		c.Srcs(
			"/**\n" +
				" * I'm test0.\n" +
				" */\n" +
				"var test0 = {a: 4, b: 5};\n" +
				"\n" +
				"// The type of this object *looks* is the same, but it is different since the type\n" +
				"// may get more properties later. Therefore, it should get a different JSDoc.\n" +
				"/**\n" +
				" * I'm test1.\n" +
				" */\n" +
				"var test1 = {a: 4, b: 5};\n"))

	test0Type := c.inferredTypeOfName("test0")
	test1Type := c.inferredTypeOfName("test1")
	// For some reason `test0Type` and `test1Type` aren't equal. This is good, but unexpected.
	AssertThat(test1Type).IsNotSameInstanceAs(test0Type)

	// Then
	AssertThat(test0Type.GetJSDocInfo().GetBlockDescription()).IsEqualTo("I'm test0.")
	AssertThat(test1Type.GetJSDocInfo().GetBlockDescription()).IsEqualTo("I'm test1.")
}

// Java: InferJSDocInfoTest.testJSDocIsPropagatedDistinctlyToMatchingStructuralTypes_FunctionTypes()
func (c *InferJSDocInfoTest) TestJSDocIsPropagatedDistinctlyToMatchingStructuralTypes_FunctionTypes() {
	// Given
	c.TestSame(
		c.Srcs(
			"/**\n" +
				" * I'm test0.\n" +
				" * @param {*} a\n" +
				" * @return {string}\n" +
				" */\n" +
				"function test0(a) {};\n" +
				"\n" +
				"// The type of this function *looks* is the same, but it is different since the type\n" +
				"// may get more properties later. Therefore, it should get a different JSDoc.\n" +
				"/**\n" +
				" * I'm test1.\n" +
				" * @param {*} a\n" +
				" * @return {string}\n" +
				" */\n" +
				"function test1(a) {};\n"))

	test0Type := c.inferredTypeOfName("test0")
	test1Type := c.inferredTypeOfName("test1")
	// We really only care that they match, not about equality.
	// TODO(b/111070482): That fact that these are equal yet have different JSDocs is bad.
	AssertThat(test1Type).IsEqualTo(test0Type)
	AssertThat(test1Type).IsNotSameInstanceAs(test0Type)

	// Then
	AssertThat(test0Type.GetJSDocInfo().GetBlockDescription()).IsEqualTo("I'm test0.")
	AssertThat(test1Type.GetJSDocInfo().GetBlockDescription()).IsEqualTo("I'm test1.")
}

// Java: InferJSDocInfoTest.testJSDocIsNotOverriddenByStructuralTypeAssignments_ObjectLiteralTypes()
func (c *InferJSDocInfoTest) TestJSDocIsNotOverriddenByStructuralTypeAssignments_ObjectLiteralTypes() {
	// Given
	c.TestSame(
		c.Srcs(
			"/**\n" +
				" * I'm test0.\n" +
				" */\n" +
				"var test0 = {a: 4, b: 5};\n" +
				"\n" +
				"// The type of these objects should be identical.\n" +
				"/**\n" +
				" * I'm test1.\n" +
				"// We need this so inferrence understands thes two objects/types are stuck together.\n" +
				" * @const\n" +
				" */\n" +
				"var test1 = test0;\n"))

	test0Type := c.inferredTypeOfName("test0")
	test1Type := c.inferredTypeOfName("test1")
	AssertThat(test1Type).IsSameInstanceAs(test0Type)

	// Then
	AssertThat(test0Type.GetJSDocInfo().GetBlockDescription()).IsEqualTo("I'm test0.")
}

// Java: InferJSDocInfoTest.testJSDocIsNotOverriddenByStructuralTypeAssignments_FunctionTypes()
func (c *InferJSDocInfoTest) TestJSDocIsNotOverriddenByStructuralTypeAssignments_FunctionTypes() {
	// Given
	c.TestSame(
		c.Srcs(
			"/**\n" +
				" * I'm test0.\n" +
				" * @param {*} a\n" +
				" * @return {string}\n" +
				" */\n" +
				"function test0(a) {};\n" +
				"\n" +
				"// The type of these objects should be identical.\n" +
				"/**\n" +
				" * I'm test1.\n" +
				" * @param {string} a\n" +
				" * @return {*}\n" +
				"// We need this so inferrence understands thes two objects/types are stuck together.\n" +
				" * @const\n" +
				" */\n" +
				"var test1 = test0;\n"))

	test0Type := c.inferredTypeOfName("test0")
	test1Type := c.inferredTypeOfName("test1")
	AssertThat(test1Type).IsSameInstanceAs(test0Type)

	// Then
	AssertThat(test0Type.GetJSDocInfo().GetBlockDescription()).IsEqualTo("I'm test0.")
}

// Java: InferJSDocInfoTest.testJSDocIsPropagatedToFieldProperties_Es5()
func (c *InferJSDocInfoTest) TestJSDocIsPropagatedToFieldProperties_Es5() {
	// Given
	c.TestSame(
		c.Srcs(
			"/** @constructor */\n" +
				"function Foo() {\n" +
				"  /**\n" +
				"   * I'm a field.\n" +
				"   * @const\n" +
				"   */\n" +
				"   this.field = 5;\n" +
				"};\n" +
				"\n" +
				"var x = new Foo();\n"))

	xType := c.inferredTypeOfName("x").(rhino.ObjectType)
	AssertThat(xType.String()).IsEqualTo("Foo")

	// Then
	AssertThat(xType.GetPropertyJSDocInfo("field").GetBlockDescription()).IsEqualTo("I'm a field.")
}

// Java: InferJSDocInfoTest.testJSDocIsPropagatedToFieldProperties_Es6Class()
func (c *InferJSDocInfoTest) TestJSDocIsPropagatedToFieldProperties_Es6Class() {
	// Given
	c.TestSame(
		c.Srcs(
			"class Foo {\n" +
				"  constructor() {\n" +
				"    /**\n" +
				"     * I'm a field.\n" +
				"     * @const\n" +
				"     */\n" +
				"     this.field = 5;\n" +
				"  }\n" +
				"}\n" +
				"\n" +
				"var x = new Foo();\n"))

	xType := c.inferredTypeOfName("x").(rhino.ObjectType)
	AssertThat(xType.String()).IsEqualTo("Foo")

	// Then
	AssertThat(xType.GetPropertyJSDocInfo("field").GetBlockDescription()).IsEqualTo("I'm a field.")
}

// Java: InferJSDocInfoTest.testJSDocIsPropagatedToClassField()
func (c *InferJSDocInfoTest) TestJSDocIsPropagatedToClassField() {
	// Given
	c.TestSame(
		c.Srcs(
			"class Foo {\n" +
				"  /**\n" +
				"   * I'm a field.\n" +
				"   * @const\n" +
				"   */\n" +
				"   field = 5;\n" +
				"}\n" +
				"\n" +
				"var x = new Foo();\n"))

	xType := c.inferredTypeOfName("x").(rhino.ObjectType)
	AssertThat(xType.String()).IsEqualTo("Foo")

	// Then
	AssertThat(xType.GetPropertyJSDocInfo("field").GetBlockDescription()).IsEqualTo("I'm a field.")
}

// Java: InferJSDocInfoTest.testJsDocIsPropagatedToClassField_computed()
func (c *InferJSDocInfoTest) TestJsDocIsPropagatedToClassField_computed() {
	// Given
	c.TestSame(
		c.Externs(
			NewTestExternsBuilder().
				AddIterable().
				AddExtra("/** @const {symbol} */ Symbol.foobar").
				Build()),
		c.Srcs(
			"class Foo {\n"+
				"  /**\n"+
				"   * I'm a field.\n"+
				"   * @const\n"+
				"   */\n"+
				"   [Symbol.foobar] = 5;\n"+
				"}\n"+
				"\n"+
				"var x = new Foo();\n"))

	xType := c.inferredTypeOfName("x").(rhino.ObjectType)
	AssertThat(xType.String()).IsEqualTo("Foo")
	foobarSymbol := c.inferredTypeOfName("Symbol.foobar").ToMaybeKnownSymbolType()

	// Then
	AssertThat(
		xType.GetPropertyJSDocInfoByKey(rhino.NewProperty_SymbolKey(foobarSymbol)).GetBlockDescription()).
		IsEqualTo("I'm a field.")
}

// Java: InferJSDocInfoTest.testJSDocIsNotPropagatedFromSuppressionsOnFieldProperties_Es5()
func (c *InferJSDocInfoTest) TestJSDocIsNotPropagatedFromSuppressionsOnFieldProperties_Es5() {
	// Given
	c.TestSame(
		c.Srcs(
			"/** @constructor */\n" +
				"function Foo() {\n" +
				"  /**\n" +
				"   * I'm a field.\n" +
				"   * @const\n" +
				"   */\n" +
				"   this.field = 5;\n" +
				"};\n" +
				"\n" +
				"var x = new Foo(); // Just a hook to access type \"Foo\".\n" +
				"/** @suppress {checkTypes} */\n" +
				"x.field = 'test';\n"))

	xType := c.inferredTypeOfName("x").(rhino.ObjectType)
	AssertThat(xType.String()).IsEqualTo("Foo")

	// Then
	AssertThat(xType.GetPropertyJSDocInfo("field").GetBlockDescription()).IsEqualTo("I'm a field.")
}

// Java: InferJSDocInfoTest.testJSDocIsNotPropagatedFromSuppressionsOnEarlyRefToFieldProperties_Es5()
func (c *InferJSDocInfoTest) TestJSDocIsNotPropagatedFromSuppressionsOnEarlyRefToFieldProperties_Es5() {
	// Given
	c.TestSame(
		c.Srcs(
			"/** @param {!Foo} x */\n" +
				"function f(x) {\n" +
				"  /** @suppress {checkTypes} */\n" +
				"  x.field = 'test';\n" +
				"}\n" +
				"/** @constructor */\n" +
				"function Foo() {\n" +
				"  /**\n" +
				"   * I'm a field.\n" +
				"   * @const\n" +
				"   */\n" +
				"   this.field = 5;\n" +
				"};\n"))

	xType := c.inferredTypeOfName("x").(rhino.ObjectType)
	AssertThat(xType.String()).IsEqualTo("Foo")

	// Then
	AssertThat(xType.GetPropertyJSDocInfo("field").GetBlockDescription()).IsEqualTo("I'm a field.")
}

// Java: InferJSDocInfoTest.testJSDocIsPropagatedToGetterProperties_Es5()
func (c *InferJSDocInfoTest) TestJSDocIsPropagatedToGetterProperties_Es5() {
	// Given
	c.TestSame(
		c.Srcs(
			"/** @constructor */\n" +
				"function Foo() {};\n" +
				"\n" +
				"Foo.prototype = {\n" +
				"  /**\n" +
				"   * I'm a getter.\n" +
				"   * @return {number}\n" +
				"   */\n" +
				"  get getter() {}\n" +
				"};\n" +
				"\n" +
				"var x = new Foo();\n"))

	xType := c.inferredTypeOfName("x").(rhino.ObjectType)
	AssertThat(xType.String()).IsEqualTo("Foo")

	// Then
	AssertThat(xType.GetPropertyJSDocInfo("getter").GetBlockDescription()).
		IsEqualTo("I'm a getter.")
}

// Java: InferJSDocInfoTest.testJSDocIsPropagatedToGetterProperties_Es6Class()
func (c *InferJSDocInfoTest) TestJSDocIsPropagatedToGetterProperties_Es6Class() {
	// Given
	c.TestSame(
		c.Srcs(
			"class Foo {\n" +
				"  /**\n" +
				"   * I'm a getter.\n" +
				"   * @return {number}\n" +
				"   */\n" +
				"  get getter() {}\n" +
				"}\n" +
				"\n" +
				"var x = new Foo();\n"))

	xType := c.inferredTypeOfName("x").(rhino.ObjectType)
	AssertThat(xType.String()).IsEqualTo("Foo")

	// Then
	AssertThat(xType.GetPropertyJSDocInfo("getter").GetBlockDescription()).
		IsEqualTo("I'm a getter.")
}

// Java: InferJSDocInfoTest.testJSDocIsPropagatedToSetterProperties_Es5()
func (c *InferJSDocInfoTest) TestJSDocIsPropagatedToSetterProperties_Es5() {
	// Given
	c.TestSame(
		c.Srcs(
			"/** @constructor */\n" +
				"function Foo() {};\n" +
				"\n" +
				"Foo.prototype = {\n" +
				"  /**\n" +
				"   * I'm a setter.\n" +
				"   * @param {number} a\n" +
				"   */\n" +
				"  set setter(a) {}\n" +
				"};\n" +
				"\n" +
				"var x = new Foo();\n"))

	xType := c.inferredTypeOfName("x").(rhino.ObjectType)
	AssertThat(xType.String()).IsEqualTo("Foo")

	// Then
	AssertThat(xType.GetPropertyJSDocInfo("setter").GetBlockDescription()).
		IsEqualTo("I'm a setter.")
}

// Java: InferJSDocInfoTest.testJSDocIsPropagatedToSetterProperties_Es6Class()
func (c *InferJSDocInfoTest) TestJSDocIsPropagatedToSetterProperties_Es6Class() {
	// Given
	c.TestSame(
		c.Srcs(
			"class Foo {\n" +
				"  /**\n" +
				"   * I'm a setter.\n" +
				"   * @param {number} a\n" +
				"   */\n" +
				"  set setter(a) {}\n" +
				"}\n" +
				"\n" +
				"var x = new Foo();\n"))

	xType := c.inferredTypeOfName("x").(rhino.ObjectType)
	AssertThat(xType.String()).IsEqualTo("Foo")

	// Then
	AssertThat(xType.GetPropertyJSDocInfo("setter").GetBlockDescription()).
		IsEqualTo("I'm a setter.")
}

// Java: InferJSDocInfoTest.testJSDocIsPropagatedToMethodProperties_Es5()
func (c *InferJSDocInfoTest) TestJSDocIsPropagatedToMethodProperties_Es5() {
	// Given
	c.TestSame(
		c.Srcs(
			"/** @constructor */\n" +
				"function Foo() {};\n" +
				"\n" +
				"/**\n" +
				" * I'm a method.\n" +
				" * @return {number} a\n" +
				" */\n" +
				"Foo.prototype.method = function() { };\n" +
				"\n" +
				"var x = new Foo();\n"))

	xType := c.inferredTypeOfName("x").(rhino.ObjectType)
	AssertThat(xType.String()).IsEqualTo("Foo")

	// Then
	AssertThat(xType.GetPropertyJSDocInfo("method").GetBlockDescription()).
		IsEqualTo("I'm a method.")
}

// Java: InferJSDocInfoTest.testJSDocIsPropagatedToMethodProperties_Es6Class()
func (c *InferJSDocInfoTest) TestJSDocIsPropagatedToMethodProperties_Es6Class() {
	// Given
	c.TestSame(
		c.Srcs(
			"class Foo {\n" +
				"  /**\n" +
				"   * I'm a method.\n" +
				"   * @return {number} a\n" +
				"   */\n" +
				"  method() { }\n" +
				"}\n" +
				"\n" +
				"var x = new Foo();\n"))

	xType := c.inferredTypeOfName("x").(rhino.ObjectType)
	AssertThat(xType.String()).IsEqualTo("Foo")

	// Then
	AssertThat(xType.GetPropertyJSDocInfo("method").GetBlockDescription()).
		IsEqualTo("I'm a method.")
}

// Java: InferJSDocInfoTest.testJsDocIsPropagatedToMethodProperties_computed()
func (c *InferJSDocInfoTest) TestJsDocIsPropagatedToMethodProperties_computed() {
	// Given
	c.TestSame(
		c.Externs(
			NewTestExternsBuilder().
				AddIterable().
				AddExtra("/** @const {symbol} */ Symbol.foobar").
				Build()),
		c.Srcs(
			"class Foo {\n"+
				"  /**\n"+
				"   * I'm a method.\n"+
				"   * @return {number} a\n"+
				"   */\n"+
				"   [Symbol.foobar]() { }\n"+
				"}\n"+
				"\n"+
				"var x = new Foo();\n"))

	xType := c.inferredTypeOfName("x").(rhino.ObjectType)
	AssertThat(xType.String()).IsEqualTo("Foo")
	foobarSymbol := c.inferredTypeOfName("Symbol.foobar").ToMaybeKnownSymbolType()

	// Then
	AssertThat(
		xType.GetPropertyJSDocInfoByKey(rhino.NewProperty_SymbolKey(foobarSymbol)).GetBlockDescription()).
		IsEqualTo("I'm a method.")
}

// Java: InferJSDocInfoTest.testJSDocIsPropagatedToAbstractMethodProperties_Es5()
func (c *InferJSDocInfoTest) TestJSDocIsPropagatedToAbstractMethodProperties_Es5() {
	// Given
	c.TestSame(
		c.Srcs(
			"/** @constructor */\n" +
				"function Foo() {};\n" +
				"\n" +
				"/**\n" +
				" * I'm a method.\n" +
				" * @abstract\n" +
				" * @return {number} a\n" +
				" */\n" +
				"Foo.prototype.method = goog.abstractMethod;\n" +
				"\n" +
				"var x = new Foo();\n"))

	xType := c.inferredTypeOfName("x").(rhino.ObjectType)
	AssertThat(xType.String()).IsEqualTo("Foo")

	// Then
	AssertThat(xType.GetPropertyJSDocInfo("method").GetBlockDescription()).
		IsEqualTo("I'm a method.")
}

// Java: InferJSDocInfoTest.testJSDocIsPropagatedToAbstractMethodProperties_Es6Class()
func (c *InferJSDocInfoTest) TestJSDocIsPropagatedToAbstractMethodProperties_Es6Class() {
	// Given
	c.TestSame(
		c.Srcs(
			"class Foo {\n" +
				"  /**\n" +
				"   * I'm a method.\n" +
				"   * @return {number} a\n" +
				"   * @abstract\n" +
				"   */\n" +
				"  method() { }\n" +
				"}\n" +
				"\n" +
				"var x = new Foo();\n"))

	xType := c.inferredTypeOfName("x").(rhino.ObjectType)
	AssertThat(xType.String()).IsEqualTo("Foo")

	// Then
	AssertThat(xType.GetPropertyJSDocInfo("method").GetBlockDescription()).
		IsEqualTo("I'm a method.")
}

// Java: InferJSDocInfoTest.testJSDocIsPropagatedToStaticProperties_Es5()
func (c *InferJSDocInfoTest) TestJSDocIsPropagatedToStaticProperties_Es5() {
	// Given
	c.TestSame(
		c.Srcs(
			"/** @constructor */\n" +
				"function Foo() {};\n" +
				"\n" +
				"/**\n" +
				" * I'm a static.\n" +
				" * @return {number} a\n" +
				" */\n" +
				"Foo.static = function() { };\n" +
				"\n" +
				"var x = Foo;\n"))

	xType := c.inferredTypeOfName("x").(rhino.ObjectType)
	AssertThat(xType.String()).IsEqualTo("(typeof Foo)")

	// Then
	AssertThat(xType.GetPropertyJSDocInfo("static").GetBlockDescription()).
		IsEqualTo("I'm a static.")
}

// Java: InferJSDocInfoTest.testJSDocIsPropagatedToStaticProperties_Es6Class()
func (c *InferJSDocInfoTest) TestJSDocIsPropagatedToStaticProperties_Es6Class() {
	// Given
	c.TestSame(
		c.Srcs(
			"class Foo {\n" +
				"  /**\n" +
				"   * I'm a static.\n" +
				"   * @return {number} a\n" +
				"   */\n" +
				"   static static() { }\n" +
				"}\n" +
				"\n" +
				"var x = Foo;\n"))

	xType := c.inferredTypeOfName("x").(rhino.ObjectType)
	AssertThat(xType.String()).IsEqualTo("(typeof Foo)")

	// Then
	AssertThat(xType.GetPropertyJSDocInfo("static").GetBlockDescription()).
		IsEqualTo("I'm a static.")
}

// TODO(b/30710701): Constructor docs should be used in some way. This is probably similar to how
// access control is being differentiated between constructor invocation and constructors as
// namespaces. The decision for both of these cases should be made together.
// Java: InferJSDocInfoTest.testJSDocFromConstructorsIsIgnored_Es6Class()
func (c *InferJSDocInfoTest) TestJSDocFromConstructorsIsIgnored_Es6Class() {
	// Given
	c.TestSame(
		c.Srcs(
			"/** I'm a class. */\n" +
				"class Foo {\n" +
				"  /** I'm a constructor. */\n" +
				"  constructor() { }\n" +
				"}\n" +
				"\n" +
				"var x = new Foo();\n"))

	xType := c.inferredTypeOfName("x").(rhino.ObjectType)
	AssertThat(xType.String()).IsEqualTo("Foo")

	// Then
	AssertThat(xType.GetJSDocInfo().GetBlockDescription()).IsEqualTo("I'm a class.")
	AssertThat(xType.GetPropertyJSDocInfo("constructor")).IsNotNull()
}

// Java: InferJSDocInfoTest.testJSDocDoesNotPropagateFromStructuralTypesToClassProperties()
func (c *InferJSDocInfoTest) TestJSDocDoesNotPropagateFromStructuralTypesToClassProperties() {
	// Given
	c.TestSame(
		c.Srcs(
			"/** @constructor */\n" +
				"function Foo() {};\n" +
				"\n" +
				"/**\n" +
				" * I'm a free function.\n" +
				" * @return {number}\n" +
				" */\n" +
				"function free() { return 0; }\n" +
				"\n" +
				"/**\n" +
				" * I'm a method.\n" +
				" * @return {number} a\n" +
				" */\n" +
				"Foo.prototype.method = free;\n" +
				"\n" +
				"var x = new Foo();\n"))

	freeType := c.inferredTypeOfName("free")
	xType := c.inferredTypeOfName("x").(rhino.ObjectType)
	AssertThat(xType).IsNotSameInstanceAs(freeType)

	// Then
	AssertThat(freeType.GetJSDocInfo().GetBlockDescription()).IsEqualTo("I'm a free function.")
	AssertThat(xType.GetPropertyJSDocInfo("method").GetBlockDescription()).
		IsEqualTo("I'm a method.")
}

// Java: InferJSDocInfoTest.testJSDocDoesNotPropagateBackwardFromInstancesToTypes()
func (c *InferJSDocInfoTest) TestJSDocDoesNotPropagateBackwardFromInstancesToTypes() {
	// Given
	c.TestSame(
		"/** @constructor */\n" +
			"function Foo() {}\n" +
			"\n" +
			"var x = new Foo();\n" +
			"\n" +
			"/** @type {number} */\n" +
			"x.bar = 4;\n")

	xType := c.inferredTypeOfName("x").(rhino.ObjectType)
	AssertThat(xType.String()).IsEqualTo("Foo")

	// Then
	AssertThat(xType.HasProperty("bar")).IsFalse()
	AssertThat(xType.GetOwnPropertyJSDocInfo("bar")).IsNull()
}

// Java: InferJSDocInfoTest.testJSDocFromDuplicateDefinitionDoesNotOverrideJSDocFromOriginal()
func (c *InferJSDocInfoTest) TestJSDocFromDuplicateDefinitionDoesNotOverrideJSDocFromOriginal() {
	// Given
	c.TestSame(
		c.Srcs(
			"/**\n" +
				" * I'm a user type.\n" +
				" * @constructor\n" +
				" */\n" +
				"var Foo = function() {};\n" +
				"\n" +
				"/**\n" +
				" * I'm a different type.\n" +
				" * @constructor\n" +
				" */\n" +
				"Foo = function() {};\n" +
				"\n" +
				"var x = new Foo();\n"))

	xType := c.inferredTypeOfName("x")
	AssertThat(xType.String()).IsEqualTo("Foo")

	// Then
	AssertThat(xType.GetJSDocInfo().GetBlockDescription()).IsEqualTo("I'm a user type.")
}

// Java: InferJSDocInfoTest.testJSDocFromDuplicateDefinitionIsUsedIfThereWasNoOriginalJSDocFromOriginal()
func (c *InferJSDocInfoTest) TestJSDocFromDuplicateDefinitionIsUsedIfThereWasNoOriginalJSDocFromOriginal() {
	// Given
	c.TestSame(
		c.Srcs(
			"var Foo = function() {};\n" +
				"\n" +
				"/**\n" +
				" * I'm a different type.\n" +
				" * @constructor\n" +
				" */\n" +
				"Foo = function() {};\n" +
				"\n" +
				"var x = new Foo();\n"))

	xType := c.inferredTypeOfName("x")
	AssertThat(xType.String()).IsEqualTo("Foo")

	// Then
	AssertThat(xType.GetJSDocInfo().GetBlockDescription()).IsEqualTo("I'm a different type.")
}

// Java: InferJSDocInfoTest.testJSDocFromDuplicatePropertyDefinitionDoesNotOverrideJSDoc_es6Style()
func (c *InferJSDocInfoTest) TestJSDocFromDuplicatePropertyDefinitionDoesNotOverrideJSDoc_es6Style() {
	// Given
	c.TestSame(
		c.Srcs(
			"class Foo {\n" +
				"\n" +
				"  /**\n" +
				"   * I'm the first definition\n" +
				"   * @param {string} s\n" +
				"   */\n" +
				"  x(s) {}\n" +
				"\n" +
				"  /**\n" +
				"   * I'm the second definition\n" +
				"   * @param {string} s\n" +
				"   */\n" +
				"  x(s) {}\n" +
				"}\n" +
				"\n" +
				"var x = new Foo();\n"))

	xType := c.inferredTypeOfName("x")
	AssertThat(xType.String()).IsEqualTo("Foo")

	// Then
	AssertThat(xType.AssertObjectType().GetPropertyJSDocInfo("x").GetBlockDescription()).
		IsEqualTo("I'm the first definition")
}

// Java: InferJSDocInfoTest.testJSDocFromDuplicatePropertyDefinitionDoesNotOverrideJSDoc_es5Style()
func (c *InferJSDocInfoTest) TestJSDocFromDuplicatePropertyDefinitionDoesNotOverrideJSDoc_es5Style() {
	// Given
	c.TestSame(
		c.Srcs(
			"/**\n" +
				" * @constructor\n" +
				" */\n" +
				"var Foo = function() {};\n" +
				"\n" +
				"/**\n" +
				" * I'm the first definition\n" +
				" * @type {string}\n" +
				" */\n" +
				"Foo.prototype.x;\n" +
				"\n" +
				"/**\n" +
				" * I'm the second definition\n" +
				" * @type {string}\n" +
				" */\n" +
				"Foo.prototype.x;\n" +
				"\n" +
				"var x = new Foo();\n"))

	xType := c.inferredTypeOfName("x")
	AssertThat(xType.String()).IsEqualTo("Foo")

	// Then
	AssertThat(xType.AssertObjectType().GetPropertyJSDocInfo("x").GetBlockDescription()).
		IsEqualTo("I'm the first definition")
}

// Java: InferJSDocInfoTest.testJSDocIsPropagatedToTypeFromObjectLiteralPrototype()
func (c *InferJSDocInfoTest) TestJSDocIsPropagatedToTypeFromObjectLiteralPrototype() {
	c.TestSame(
		"/** @constructor */\n" +
			"function Foo() {}\n" +
			"\n" +
			"Foo.prototype = {\n" +
			"  /** Property a. */ a: function() {},\n" +
			"  /** Property b. */ get b() {},\n" +
			"  /** Property c. */ set c(x) {},\n" +
			"  /** Property d. */ d() {}\n" +
			"};\n")

	ctor := c.inferredTypeOfName("Foo").ToMaybeFunctionType()
	prototype := ctor.GetInstanceType().GetImplicitPrototype()

	AssertThat(prototype.GetOwnPropertyJSDocInfo("a").GetBlockDescription()).
		IsEqualTo("Property a.")
	AssertThat(prototype.GetOwnPropertyJSDocInfo("b").GetBlockDescription()).
		IsEqualTo("Property b.")
	AssertThat(prototype.GetOwnPropertyJSDocInfo("c").GetBlockDescription()).
		IsEqualTo("Property c.")
	AssertThat(prototype.GetOwnPropertyJSDocInfo("d").GetBlockDescription()).
		IsEqualTo("Property d.")
}

// Java: InferJSDocInfoTest.testClassWithUnknownTypeDoesNotCrash()
func (c *InferJSDocInfoTest) TestClassWithUnknownTypeDoesNotCrash() {
	// Code that looks like this and @suppresses {checkTypes} is generated by tsickle.
	c.IgnoreWarnings(TypeCheck_CONFLICTING_EXTENDED_TYPE)
	c.EnableTypeCheck()
	c.TestSame(
		"/** @param {?} base */\n" +
			"function mixin(base) {\n" +
			"  return (/** @type {?} */ ((class extends base {\n" +
			"    /** @param {...?} args */ constructor(...args) {}\n" +
			"  })));\n" +
			"}\n")
}

// inferredTypeOfName returns the inferred type of the reference name
// anywhere in the AST.
//
// Java: InferJSDocInfoTest.inferredTypeOfName(String)
func (c *InferJSDocInfoTest) inferredTypeOfName(name string) rhino.JSType {
	return c.inferredTypeOfLocalName(name, c.GetLastCompiler().GetRoot())
}

// inferredTypeOfLocalName returns the inferred type of the reference name
// under root AST node..
//
// Java: InferJSDocInfoTest.inferredTypeOfLocalName(String, Node)
func (c *InferJSDocInfoTest) inferredTypeOfLocalName(name string, root *rhino.Node) rhino.JSType {
	// ArrayDeque used as a stack (push / pop at the head).
	queue := []*rhino.Node{}
	queue = append(queue, root)

	for len(queue) != 0 {
		current := queue[len(queue)-1]
		queue = queue[:len(queue)-1]
		if current.MatchesQualifiedName(name) && current.GetJSType() != nil {
			return current.GetJSType()
		}

		for child := current.GetFirstChild(); child != nil; child = child.GetNext() {
			queue = append(queue, child)
		}
	}
	return nil
}
