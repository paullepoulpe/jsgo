// Port of com/google/javascript/jscomp/AstFactoryTest.java @ closure-compiler v20261015

/*
 * Copyright 2018 The Closure Compiler Authors.
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
	"reflect"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/rhinotesting"
)

// The Java @Test methods not named test* get a "Test" prefix so that
// RunJUnit4 runs them.
//
// Java: public class AstFactoryTest
type AstFactoryTest struct {
	// Go-only: test-helper — the Go test running the current test method.
	t *testing.T

	compiler            *Compiler
	runtimeJsLibManager *RuntimeJsLibManager
}

func TestAstFactoryTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return &AstFactoryTest{} })
}

// JUnitSetT gives the instance the Go test running its test method.
func (c *AstFactoryTest) JUnitSetT(t *testing.T) { c.t = t }

// JUnitBefore runs the @Before method.
func (c *AstFactoryTest) JUnitBefore() { c.SetUp() }

// JUnitAfter: the class has no @After method.
func (c *AstFactoryTest) JUnitAfter() {}

// Java: AstFactoryTest.setUp()
func (c *AstFactoryTest) SetUp() {
	c.compiler = NewCompiler()
	c.runtimeJsLibManager =
		RuntimeJsLibManager_Create(
			RuntimeJsLibManager_RuntimeLibraryMode_RECORD_AND_VALIDATE_FIELDS,
			astFactoryTest_TestResourceProvider,
			c.compiler.GetChangeTracker(),
			astFactoryTest_alwaysThrowNodeSupplier)
}

// Java: AstFactoryTest.alwaysThrowNodeSupplier()
func astFactoryTest_alwaysThrowNodeSupplier() rhino.Node {
	panic(jutil.NewUnsupportedOperationException())
}

// Java: AstFactoryTest.TestResourceProvider.parse(String, String)
func astFactoryTest_TestResourceProvider(unused1 string, unused2 string) rhino.Node {
	panic(jutil.NewUnsupportedOperationException())
}

// Java: AstFactoryTest.getRegistry()
func (c *AstFactoryTest) getRegistry() *rhino.JSTypeRegistry {
	return c.compiler.GetTypeRegistry()
}

// Java: AstFactoryTest.getNativeType(JSTypeNative)
func (c *AstFactoryTest) getNativeType(nativeType rhino.JSTypeNative) rhino.JSType {
	return c.getRegistry().GetNativeType(nativeType)
}

// Java: AstFactoryTest.parseWithoutTypes(String)
func (c *AstFactoryTest) parseWithoutTypes(source string) rhino.Node {
	return c.parseWithoutTypesWithExterns("", source)
}

// Java: AstFactoryTest.parseWithoutTypes(String, String)
func (c *AstFactoryTest) parseWithoutTypesWithExterns(externs string, source string) rhino.Node {
	// parse the test code
	options := NewCompilerOptions()
	c.compiler.Init(
		[]*SourceFile{SourceFile_FromCode("externs", externs)},
		[]*SourceFile{SourceFile_FromCode("source", source)},
		options)
	c.compiler.ParseInputs()

	AssertWithMessage("parse error").That(c.compiler.GetErrors()).IsEmpty()
	return c.compiler.GetJsRoot()
}

// Java: AstFactoryTest.parseAndAddTypes(String)
func (c *AstFactoryTest) parseAndAddTypes(source string) rhino.Node {
	return c.parseAndAddTypesWithExterns("", source)
}

// Java: AstFactoryTest.parseAndAddTypes(String, String)
func (c *AstFactoryTest) parseAndAddTypesWithExterns(externs string, source string) rhino.Node {
	c.parseWithoutTypesWithExterns(externs, source)

	// then add types to the AST & do type checks
	// TODO(bradfordcsmith): Fail if there are type checking errors.
	typeCheck :=
		NewTypeCheck(
			c.compiler, c.compiler.GetReverseAbstractInterpreter(), c.compiler.GetTypeRegistry())
	typeCheck.ProcessForTesting(c.compiler.GetExternsRoot(), c.compiler.GetJsRoot())
	c.compiler.SetTypeCheckingHasRun(true)
	return c.compiler.GetJsRoot()
}

// Java: AstFactoryTest.parseAndAddColors(String)
func (c *AstFactoryTest) parseAndAddColors(source string) rhino.Node {
	return c.parseAndAddColorsWithExterns("", source)
}

// Java: AstFactoryTest.parseAndAddColors(String, String)
func (c *AstFactoryTest) parseAndAddColorsWithExterns(externs string, source string) rhino.Node {
	c.parseAndAddTypesWithExterns(externs, source)
	NewConvertTypesToColors(
		c.compiler, SerializationOptions_Builder_().SetIncludeDebugInfo(true).Build()).
		Process(c.compiler.GetExternsRoot(), c.compiler.GetJsRoot())
	return c.compiler.GetJsRoot()
}

// Java: AstFactoryTest.createTestAstFactory()
func (c *AstFactoryTest) createTestAstFactory() *AstFactory {
	return AstFactory_CreateFactoryWithTypes(
		c.compiler.GetLifeCycleStage(), c.getRegistry(), c.runtimeJsLibManager)
}

// Java: AstFactoryTest.createTestAstFactoryWithColors()
func (c *AstFactoryTest) createTestAstFactoryWithColors() *AstFactory {
	// the built-in color registry is available only if we've run parseAndAddColors()
	var registry *colors.ColorRegistry
	if c.compiler.HasOptimizationColors() {
		registry = c.compiler.GetColorRegistry()
	} else {
		registry = colors.ColorRegistry_Builder_().SetDefaultNativeColorsForTesting().Build()
	}
	return AstFactory_CreateFactoryWithColors(c.compiler.GetLifeCycleStage(), registry, c.runtimeJsLibManager)
}

// Java: AstFactoryTest.createTestAstFactoryWithoutTypes()
func (c *AstFactoryTest) createTestAstFactoryWithoutTypes() *AstFactory {
	return AstFactory_CreateFactoryWithoutTypes(c.compiler.GetLifeCycleStage(), c.runtimeJsLibManager)
}

// astFactoryTest_redeclarationHandler is getScope's lambda
// (Scope s, String name, Node n, CompilerInput input) -> {}.
type astFactoryTest_redeclarationHandler struct {
	_ byte // Go-only: language-adapter — distinct identity
}

func (*astFactoryTest_redeclarationHandler) OnRedeclaration(s *Scope, name string, n rhino.Node, input *CompilerInput) {
}

// Java: AstFactoryTest.getScope(Node)
func (c *AstFactoryTest) getScope(root rhino.Node) *Scope {
	// Normal passes use SyntacticScopeCreator, so that's what we use here.
	var redeclarationHandler SyntacticScopeCreator_RedeclarationHandler = &astFactoryTest_redeclarationHandler{}
	scopeCreator := NewSyntacticScopeCreatorWithRedeclarationHandler(c.compiler, redeclarationHandler)
	return scopeCreator.CreateScope(root, nil).(*Scope)
}

// Java: AstFactoryTest.testStringLiteral_jstypes()
func (c *AstFactoryTest) TestStringLiteral_jstypes() {
	astFactory := c.createTestAstFactory()
	stringLiteral := astFactory.CreateString("hello")
	AssertNode(stringLiteral).HasType(rhino.Token_STRINGLIT)
	AssertThat(stringLiteral.GetString()).IsEqualTo("hello")
	TypeSubject_AssertType(stringLiteral.GetJSType()).IsString()
}

// Java: AstFactoryTest.testStringLiteral_colors()
func (c *AstFactoryTest) TestStringLiteral_colors() {
	astFactory := c.createTestAstFactoryWithColors()
	stringLiteral := astFactory.CreateString("hello")
	AssertNode(stringLiteral).HasType(rhino.Token_STRINGLIT)
	AssertThat(stringLiteral.GetString()).IsEqualTo("hello")
	AssertNode(stringLiteral).HasColorThat().IsEqualTo(colors.StandardColors_STRING)
}

// Java: AstFactoryTest.testNumberLiteral_jstypes()
func (c *AstFactoryTest) TestNumberLiteral_jstypes() {
	astFactory := c.createTestAstFactory()
	numberLiteral := astFactory.CreateNumber(2112.0)
	AssertNode(numberLiteral).HasType(rhino.Token_NUMBER)
	AssertThat(numberLiteral.GetDouble()).IsEqualTo(2112.0)
	TypeSubject_AssertType(numberLiteral.GetJSType()).IsNumber()
}

// Java: AstFactoryTest.testNumberLiteral_colors()
func (c *AstFactoryTest) TestNumberLiteral_colors() {
	astFactory := c.createTestAstFactoryWithColors()
	numberLiteral := astFactory.CreateNumber(2112.0)
	AssertNode(numberLiteral).HasType(rhino.Token_NUMBER)
	AssertThat(numberLiteral.GetDouble()).IsEqualTo(2112.0)
	AssertNode(numberLiteral).HasColorThat().IsEqualTo(colors.StandardColors_NUMBER)
}

// Java: AstFactoryTest.testBooleanLiteral_jstypes()
func (c *AstFactoryTest) TestBooleanLiteral_jstypes() {
	astFactory := c.createTestAstFactory()
	trueNode := astFactory.CreateBoolean(true)
	AssertNode(trueNode).HasType(rhino.Token_TRUE)
	TypeSubject_AssertType(trueNode.GetJSType()).IsBoolean()
	falseNode := astFactory.CreateBoolean(false)
	AssertNode(falseNode).HasType(rhino.Token_FALSE)
	TypeSubject_AssertType(falseNode.GetJSType()).IsBoolean()
}

// Java: AstFactoryTest.testBooleanLiteral_colors()
func (c *AstFactoryTest) TestBooleanLiteral_colors() {
	astFactory := c.createTestAstFactoryWithColors()
	trueNode := astFactory.CreateBoolean(true)
	AssertNode(trueNode).HasType(rhino.Token_TRUE)
	AssertNode(trueNode).HasColorThat().IsEqualTo(colors.StandardColors_BOOLEAN)
}

// Java: AstFactoryTest.testVoidExpression_jstypes()
func (c *AstFactoryTest) TestVoidExpression_jstypes() {
	astFactory := c.createTestAstFactory()
	voidNode := astFactory.CreateVoid(astFactory.CreateNumber(0))
	AssertNode(voidNode).HasType(rhino.Token_VOID)
	AssertNode(voidNode).HasJSTypeThat().IsVoid()
}

// Java: AstFactoryTest.testVoidExpression_colors()
func (c *AstFactoryTest) TestVoidExpression_colors() {
	astFactory := c.createTestAstFactoryWithColors()
	voidNode := astFactory.CreateVoid(astFactory.CreateNumber(0))
	AssertNode(voidNode).HasType(rhino.Token_VOID)
	AssertNode(voidNode).HasColorThat().IsEqualTo(colors.StandardColors_NULL_OR_VOID)
}

// Java: AstFactoryTest.testNotExpression_jstypes()
func (c *AstFactoryTest) TestNotExpression_jstypes() {
	astFactory := c.createTestAstFactory()
	notNode := astFactory.CreateNot(astFactory.CreateNumber(0))
	AssertNode(notNode).HasType(rhino.Token_NOT)
	AssertNode(notNode).HasJSTypeThat().IsBoolean()
}

// Java: AstFactoryTest.testNotExpression_colors()
func (c *AstFactoryTest) TestNotExpression_colors() {
	astFactory := c.createTestAstFactoryWithColors()
	notNode := astFactory.CreateNot(astFactory.CreateNumber(0))
	AssertNode(notNode).HasType(rhino.Token_NOT)
	AssertNode(notNode).HasColorThat().IsEqualTo(colors.StandardColors_BOOLEAN)
}

// Java: AstFactoryTest.testCreateArgumentsReference_jstypes()
func (c *AstFactoryTest) TestCreateArgumentsReference_jstypes() {
	// Make sure the compiler's type registry includes the standard externs definition for
	// Arguments.
	c.parseAndAddTypesWithExterns(NewTestExternsBuilder().AddArguments().Build(), "")
	astFactory := c.createTestAstFactory()
	argumentsNode := astFactory.CreateArgumentsReference()
	AssertNode(argumentsNode).MatchesName("arguments")
	TypeSubject_AssertType(argumentsNode.GetJSType()).IsEqualTo(c.getRegistry().GetGlobalType("Arguments"))
}

// Java: AstFactoryTest.testCreateArgumentsReference_colors()
func (c *AstFactoryTest) TestCreateArgumentsReference_colors() {
	root := c.parseAndAddColorsWithExterns(
		NewTestExternsBuilder().AddArguments().Build(), "function f() { arguments; }")
	astFactory := c.createTestAstFactoryWithColors()
	block := NodeUtil_GetFunctionBody(root.GetFirstFirstChild())
	argumentsReferenceNode := block.GetFirstFirstChild()
	argumentsReferenceColor := argumentsReferenceNode.GetColor()
	argumentsNode := astFactory.CreateArgumentsReference()
	AssertNode(argumentsNode).MatchesName("arguments")
	AssertNode(argumentsNode).HasColorThat().IsEqualTo(argumentsReferenceColor)
}

// Java: AstFactoryTest.testCreateSingleVarNameDeclaration_jstypes()
func (c *AstFactoryTest) TestCreateSingleVarNameDeclaration_jstypes() {
	astFactory := c.createTestAstFactory()
	valueNode := astFactory.CreateBoolean(true)
	constNode := astFactory.CreateSingleVarNameDeclarationWithValue("myTrue", valueNode)
	AssertNode(constNode).IsVar().HasOneChildThat().IsName("myTrue")
	nameNode := constNode.GetOnlyChild()
	AssertNode(nameNode).HasOneChildThat().IsEqualTo(valueNode)
	AssertNode(nameNode).HasJSTypeThat().IsEqualTo(valueNode.GetJSType())
}

// Java: AstFactoryTest.testCreateSingleVarNameDeclaration_colors()
func (c *AstFactoryTest) TestCreateSingleVarNameDeclaration_colors() {
	astFactory := c.createTestAstFactoryWithColors()
	valueNode := astFactory.CreateBoolean(true)
	constNode := astFactory.CreateSingleVarNameDeclarationWithValue("myTrue", valueNode)
	AssertNode(constNode).IsVar().HasOneChildThat().IsName("myTrue")
	nameNode := constNode.GetOnlyChild()
	AssertNode(nameNode).HasOneChildThat().IsEqualTo(valueNode)
	AssertNode(nameNode).HasColorThat().IsEqualTo(valueNode.GetColor())
}

// Java: AstFactoryTest.testCreateSingleConstNameDeclaration_jstypes()
func (c *AstFactoryTest) TestCreateSingleConstNameDeclaration_jstypes() {
	astFactory := c.createTestAstFactory()
	valueNode := astFactory.CreateBoolean(true)
	constNode := astFactory.CreateSingleConstNameDeclaration("myTrue", valueNode)
	AssertNode(constNode).IsConst().HasOneChildThat().IsName("myTrue")
	nameNode := constNode.GetOnlyChild()
	AssertNode(nameNode).HasOneChildThat().IsEqualTo(valueNode)
	AssertNode(nameNode).HasJSTypeThat().IsEqualTo(valueNode.GetJSType())
}

// Java: AstFactoryTest.testCreateSingleConstNameDeclaration_colors()
func (c *AstFactoryTest) TestCreateSingleConstNameDeclaration_colors() {
	astFactory := c.createTestAstFactoryWithColors()
	valueNode := astFactory.CreateBoolean(true)
	constNode := astFactory.CreateSingleConstNameDeclaration("myTrue", valueNode)
	AssertNode(constNode).IsConst().HasOneChildThat().IsName("myTrue")
	nameNode := constNode.GetOnlyChild()
	AssertNode(nameNode).HasOneChildThat().IsEqualTo(valueNode)
	AssertNode(nameNode).HasColorThat().IsEqualTo(valueNode.GetColor())
}

// Java: AstFactoryTest.testCreateNameWithJSType()
func (c *AstFactoryTest) TestCreateNameWithJSType() {
	astFactory := c.createTestAstFactory()
	x := astFactory.CreateName("x", AstFactory_Type_OfJSTypeNative(rhino.JSTypeNative_STRING_TYPE))
	AssertNode(x).HasType(rhino.Token_NAME)
	AssertThat(x.GetString()).IsEqualTo("x")
	TypeSubject_AssertType(x.GetJSType()).IsString()
}

// Java: AstFactoryTest.testCreateNameWithNativeType()
func (c *AstFactoryTest) TestCreateNameWithNativeType() {
	astFactory := c.createTestAstFactory()
	x := astFactory.CreateName("x", AstFactory_Type_OfJSTypeNative(rhino.JSTypeNative_STRING_TYPE))
	AssertNode(x).HasType(rhino.Token_NAME)
	AssertThat(x.GetString()).IsEqualTo("x")
	TypeSubject_AssertType(x.GetJSType()).IsString()
}

// Java: AstFactoryTest.testCreateNameWithColor()
func (c *AstFactoryTest) TestCreateNameWithColor() {
	astFactory := c.createTestAstFactoryWithColors()
	x := astFactory.CreateName("x", AstFactory_Type_OfJSTypeNativeAndColor(rhino.JSTypeNative_STRING_TYPE, colors.StandardColors_STRING))
	AssertNode(x).HasType(rhino.Token_NAME)
	AssertThat(x.GetString()).IsEqualTo("x")
	AssertNode(x).HasColorThat().IsEqualTo(colors.StandardColors_STRING)
}

// Java: AstFactoryTest.testCreateNameFromScope_jstypes()
func (c *AstFactoryTest) TestCreateNameFromScope_jstypes() {
	astFactory := c.createTestAstFactory()
	root := c.parseAndAddTypes("/** @type {string} */ const X = 'hi';")
	scope := c.getScope(root)
	x := astFactory.CreateNameInScope(scope, "X")
	AssertNode(x).HasType(rhino.Token_NAME)
	AssertThat(x.GetString()).IsEqualTo("X")
	TypeSubject_AssertType(x.GetJSType()).IsString()
}

// Java: AstFactoryTest.testCreateNameFromScope_colors()
func (c *AstFactoryTest) TestCreateNameFromScope_colors() {
	astFactory := c.createTestAstFactoryWithColors()
	root := c.parseAndAddColors("/** @type {string} */ const X = 'hi';")
	scope := c.getScope(root)
	x := astFactory.CreateNameInScope(scope, "X")
	AssertNode(x).HasType(rhino.Token_NAME)
	AssertThat(x.GetString()).IsEqualTo("X")
	AssertNode(x).HasColorThat().IsEqualTo(colors.StandardColors_STRING)
}

// Java: AstFactoryTest.testCreateNameFromScope_crashesIfMissingVariable()
func (c *AstFactoryTest) TestCreateNameFromScope_crashesIfMissingVariable() {
	astFactory := c.createTestAstFactory()
	root := c.parseAndAddTypes("/** @type {string} */ const X = 'hi';")
	scope := c.getScope(root)
	AssertThrows[jutil.Throwable](func() { astFactory.CreateNameInScope(scope, "missing") })
}

// Java: AstFactoryTest.testCreateGetPropFromScope_defaultsToUnknownJSTypeWhenNull()
func (c *AstFactoryTest) TestCreateGetPropFromScope_defaultsToUnknownJSTypeWhenNull() {
	astFactory := c.createTestAstFactory()
	receiver := astFactory.CreateNameWithUnknownType("JQ")
	typeTemplate := rhino.IR_Name("JQ")
	jutil.CheckStateMsg(typeTemplate.GetJSType() == nil, "getJSType does not return null ")
	x := astFactory.CreateGetProp(receiver, "$", AstFactory_Type_(typeTemplate))
	AssertNode(x).HasType(rhino.Token_GETPROP)
	AssertNode(x).MatchesQualifiedName("JQ.$")
	AssertNode(x).HasJSTypeThat().IsUnknown()
}

// Java: AstFactoryTest.testCreateGetPropFromScope_defaultsToUnknownColorWhenNull()
func (c *AstFactoryTest) TestCreateGetPropFromScope_defaultsToUnknownColorWhenNull() {
	astFactory := c.createTestAstFactoryWithColors()
	receiver := astFactory.CreateNameWithUnknownType("JQ")
	typeTemplate := rhino.IR_Name("JQ")
	jutil.CheckStateMsg(typeTemplate.GetColor() == nil, "getColor does not return null ")
	x := astFactory.CreateGetProp(receiver, "$", AstFactory_Type_(typeTemplate))
	AssertNode(x).HasType(rhino.Token_GETPROP)
	AssertNode(x).MatchesQualifiedName("JQ.$")
	AssertNode(x).HasColorThat().IsEqualTo(colors.StandardColors_UNKNOWN)
}

// Java: AstFactoryTest.testCreateThisReference()
func (c *AstFactoryTest) TestCreateThisReference() {
	astFactory := c.createTestAstFactory()
	x := astFactory.CreateThis(AstFactory_Type_OfJSType(c.getNativeType(rhino.JSTypeNative_STRING_TYPE)))
	AssertNode(x).HasType(rhino.Token_THIS)
	TypeSubject_AssertType(x.GetJSType()).IsString()
}

// Java: AstFactoryTest.testCreateSuperReference()
func (c *AstFactoryTest) TestCreateSuperReference() {
	astFactory := c.createTestAstFactory()
	x := astFactory.CreateSuper(AstFactory_Type_OfJSType(c.getNativeType(rhino.JSTypeNative_STRING_TYPE)))
	AssertNode(x).HasType(rhino.Token_SUPER)
	TypeSubject_AssertType(x.GetJSType()).IsString()
}

// Java: AstFactoryTest.testCreateNewTargetReference()
func (c *AstFactoryTest) TestCreateNewTargetReference() {
	astFactory := c.createTestAstFactory()

	x := astFactory.CreateNewTarget(AstFactory_Type_OfJSType(c.getNativeType(rhino.JSTypeNative_STRING_TYPE)))
	AssertNode(x).HasType(rhino.Token_NEW_TARGET)
	TypeSubject_AssertType(x.GetJSType()).IsString()
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: AstFactoryTest.createThisForEs6ClassMember_jstypes()
func (c *AstFactoryTest) TestCreateThisForEs6ClassMember_jstypes() {
	astFactory := c.createTestAstFactory()
	root := c.parseAndAddTypes(
		`class C {
  method() {}
}
`)
	classNode := root.GetFirstChild(). // script
						GetFirstChild()
	// class node
	memberDef := classNode.
		GetLastChild(). // class members
		GetFirstChild()
	// member function def
	instanceType := classNode.GetJSTypeRequired().AssertFunctionType().GetInstanceType()
	thisAlias := astFactory.CreateThisForEs6ClassMember(memberDef)
	AssertNode(thisAlias).HasType(rhino.Token_THIS)
	AssertNode(thisAlias).HasJSTypeThat().IsEqualTo(instanceType)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: AstFactoryTest.createThisForEs6ClassMember_colors()
func (c *AstFactoryTest) TestCreateThisForEs6ClassMember_colors() {
	astFactory := c.createTestAstFactoryWithColors()
	root := c.parseAndAddColors(
		`class C {
  method() {}
}
`)
	classNode := root.GetFirstChild(). // script
						GetFirstChild()
	// class node
	memberDef := classNode.
		GetLastChild(). // class members
		GetFirstChild()
	// member function def
	instanceType := colors.Color_CreateUnion(classNode.GetColor().GetInstanceColors())
	thisAlias := astFactory.CreateThisForEs6ClassMember(memberDef)
	AssertNode(thisAlias).HasType(rhino.Token_THIS)
	AssertNode(thisAlias).HasColorThat().IsEqualTo(instanceType)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: AstFactoryTest.createThisForEs6ClassStaticMember_jstypes()
func (c *AstFactoryTest) TestCreateThisForEs6ClassStaticMember_jstypes() {
	astFactory := c.createTestAstFactory()
	root := c.parseAndAddTypes(
		`class C {
  static method() {}
}
`)
	classNode := root.GetFirstChild(). // script
						GetFirstChild()
	// class node
	memberDef := classNode.
		GetLastChild(). // class members
		GetFirstChild()
	// member function def
	thisAlias := astFactory.CreateThisForEs6ClassMember(memberDef)
	AssertNode(thisAlias).HasType(rhino.Token_THIS)
	AssertNode(thisAlias).HasJSTypeThat().IsEqualTo(classNode.GetJSType())
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: AstFactoryTest.createThisForEs6ClassStaticMember_colors()
func (c *AstFactoryTest) TestCreateThisForEs6ClassStaticMember_colors() {
	astFactory := c.createTestAstFactoryWithColors()
	root := c.parseAndAddColors(
		`class C {
  static method() {}
}
`)
	classNode := root.GetFirstChild(). // script
						GetFirstChild()
	// class node
	memberDef := classNode.
		GetLastChild(). // class members
		GetFirstChild()
	// member function def
	thisAlias := astFactory.CreateThisForEs6ClassMember(memberDef)
	AssertNode(thisAlias).HasType(rhino.Token_THIS)
	AssertNode(thisAlias).HasColorThat().IsEqualTo(classNode.GetColor())
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: AstFactoryTest.createThisAliasReferenceForFunction()
func (c *AstFactoryTest) TestCreateThisAliasReferenceForFunction() {
	astFactory := c.createTestAstFactory()
	root := c.parseAndAddTypes(
		`class C {
  method() {}
}
`)
	classNode := root.GetFirstChild(). // script
						GetFirstChild()
	// class node
	instanceType := classNode.GetJSTypeRequired().AssertFunctionType().GetInstanceType()
	thisAlias := astFactory.CreateThisAliasReferenceForEs6Class("thisAlias", classNode)
	AssertNode(thisAlias).HasType(rhino.Token_NAME)
	AssertThat(thisAlias.GetString()).IsEqualTo("thisAlias")
	AssertNode(thisAlias).HasJSTypeThat().IsEqualTo(instanceType)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: AstFactoryTest.createThisAliasReferenceForFunctionWithoutTypes()
func (c *AstFactoryTest) TestCreateThisAliasReferenceForFunctionWithoutTypes() {
	astFactory := c.createTestAstFactoryWithoutTypes()
	root := c.parseWithoutTypes(
		`class C {
  method() {}
}
`)
	classNode := root.GetFirstChild(). // script
						GetFirstChild()
	// class
	thisAlias := astFactory.CreateThisAliasReferenceForEs6Class("thisAlias", classNode)
	AssertNode(thisAlias).HasType(rhino.Token_NAME)
	AssertThat(thisAlias.GetString()).IsEqualTo("thisAlias")
	AssertThat(thisAlias.GetJSType()).IsNull()
}

// Java: AstFactoryTest.testCreateGetpropWithColorFromNode()
func (c *AstFactoryTest) TestCreateGetpropWithColorFromNode() {
	astFactory := c.createTestAstFactoryWithColors()
	receiver := astFactory.CreateNameWithUnknownType("x")
	typeTemplate := astFactory.CreateNumber(0)
	getProp := astFactory.CreateGetProp(receiver, "y", AstFactory_Type_(typeTemplate))
	AssertNode(getProp).HasToken(rhino.Token_GETPROP)
	AssertNode(getProp).HasStringThat().IsEqualTo("y")
	AssertNode(getProp).HasFirstChildThat().IsEqualTo(receiver)
	AssertNode(getProp).HasColorThat().IsEqualTo(colors.StandardColors_NUMBER)
}

// Java: AstFactoryTest.testCreateStartOptChainGetprop()
func (c *AstFactoryTest) TestCreateStartOptChainGetprop() {
	astFactory := c.createTestAstFactoryWithColors()
	receiver := astFactory.CreateNameWithUnknownType("x")
	typeTemplate := astFactory.CreateNumber(0)
	getProp := astFactory.CreateStartOptChainGetprop(receiver, "y", AstFactory_Type_(typeTemplate))
	AssertNode(getProp).HasToken(rhino.Token_OPTCHAIN_GETPROP)
	AssertNode(getProp).HasStringThat().IsEqualTo("y")
	AssertNode(getProp).HasFirstChildThat().IsEqualTo(receiver)
	AssertNode(getProp).HasColorThat().IsEqualTo(colors.StandardColors_NUMBER)
	AssertNode(getProp).IsOptionalChainStart()
}

// Java: AstFactoryTest.testCreateContinueOptChainGetprop()
func (c *AstFactoryTest) TestCreateContinueOptChainGetprop() {
	astFactory := c.createTestAstFactoryWithColors()
	receiver := astFactory.CreateNameWithUnknownType("x")
	typeTemplate := astFactory.CreateNumber(0)
	getProp := astFactory.CreateContinueOptChainGetprop(receiver, "y", AstFactory_Type_(typeTemplate))
	AssertNode(getProp).HasToken(rhino.Token_OPTCHAIN_GETPROP)
	AssertNode(getProp).HasStringThat().IsEqualTo("y")
	AssertNode(getProp).HasFirstChildThat().IsEqualTo(receiver)
	AssertNode(getProp).HasColorThat().IsEqualTo(colors.StandardColors_NUMBER)
	AssertNode(getProp).IsNotOptionalChainStart()
}

// Java: AstFactoryTest.testCreateStartOptChainGetelem()
func (c *AstFactoryTest) TestCreateStartOptChainGetelem() {
	astFactory := c.createTestAstFactoryWithColors()
	receiver := astFactory.CreateNameWithUnknownType("x")
	elem := astFactory.CreateNumber(0)
	typeTemplate := astFactory.CreateNumber(0)
	getElem := astFactory.CreateStartOptChainGetelem(receiver, elem, AstFactory_Type_(typeTemplate))
	AssertNode(getElem).HasToken(rhino.Token_OPTCHAIN_GETELEM)
	AssertNode(getElem).HasFirstChildThat().IsEqualTo(receiver)
	AssertNode(getElem).HasSecondChildThat().IsEqualTo(elem)
	AssertNode(getElem).HasColorThat().IsEqualTo(colors.StandardColors_NUMBER)
	AssertNode(getElem).IsOptionalChainStart()
}

// Java: AstFactoryTest.testCreateContinueOptChainGetelem()
func (c *AstFactoryTest) TestCreateContinueOptChainGetelem() {
	astFactory := c.createTestAstFactoryWithColors()
	receiver := astFactory.CreateNameWithUnknownType("x")
	elem := astFactory.CreateNumber(0)
	typeTemplate := astFactory.CreateNumber(0)
	getElem := astFactory.CreateContinueOptChainGetelem(receiver, elem, AstFactory_Type_(typeTemplate))
	AssertNode(getElem).HasToken(rhino.Token_OPTCHAIN_GETELEM)
	AssertNode(getElem).HasFirstChildThat().IsEqualTo(receiver)
	AssertNode(getElem).HasSecondChildThat().IsEqualTo(elem)
	AssertNode(getElem).HasColorThat().IsEqualTo(colors.StandardColors_NUMBER)
	AssertNode(getElem).IsNotOptionalChainStart()
}

// Java: AstFactoryTest.testCreateStartOptChainCall()
func (c *AstFactoryTest) TestCreateStartOptChainCall() {
	astFactory := c.createTestAstFactoryWithColors()
	receiver := astFactory.CreateNameWithUnknownType("x")
	typeTemplate := astFactory.CreateNumber(0)
	arg1 := astFactory.CreateNumber(1)
	arg2 := astFactory.CreateNumber(2)
	call := astFactory.CreateStartOptChainCall(receiver, AstFactory_Type_(typeTemplate), arg1, arg2)
	AssertNode(call).HasToken(rhino.Token_OPTCHAIN_CALL)
	AssertNode(call).HasFirstChildThat().IsEqualTo(receiver)
	AssertNode(call).HasSecondChildThat().IsEqualTo(arg1)
	AssertNode(call).HasLastChildThat().IsEqualTo(arg2)
	AssertNode(call).HasColorThat().IsEqualTo(colors.StandardColors_NUMBER)
	AssertNode(call).IsOptionalChainStart()
}

// Java: AstFactoryTest.testCreateContinueOptChainCall()
func (c *AstFactoryTest) TestCreateContinueOptChainCall() {
	astFactory := c.createTestAstFactoryWithColors()
	receiver := astFactory.CreateNameWithUnknownType("x")
	typeTemplate := astFactory.CreateNumber(0)
	arg1 := astFactory.CreateNumber(1)
	arg2 := astFactory.CreateNumber(2)
	call := astFactory.CreateContinueOptChainCall(receiver, AstFactory_Type_(typeTemplate), arg1, arg2)
	AssertNode(call).HasToken(rhino.Token_OPTCHAIN_CALL)
	AssertNode(call).HasFirstChildThat().IsEqualTo(receiver)
	AssertNode(call).HasSecondChildThat().IsEqualTo(arg1)
	AssertNode(call).HasLastChildThat().IsEqualTo(arg2)
	AssertNode(call).HasColorThat().IsEqualTo(colors.StandardColors_NUMBER)
	AssertNode(call).IsNotOptionalChainStart()
}

// Java: AstFactoryTest.testCreateStringKey_jstypes()
func (c *AstFactoryTest) TestCreateStringKey_jstypes() {
	astFactory := c.createTestAstFactory()
	numberNode := astFactory.CreateNumber(2112.0)
	stringKeyNode := astFactory.CreateStringKey("key", numberNode)
	AssertNode(stringKeyNode).HasType(rhino.Token_STRING_KEY)
	AssertThat(stringKeyNode.GetString()).IsEqualTo("key")
	AssertThat(astFactoryTest_childList(stringKeyNode)).ContainsExactly(numberNode)
	TypeSubject_AssertType(stringKeyNode.GetJSType()).IsNumber()
}

// Java: AstFactoryTest.testCreateStringKey_colors()
func (c *AstFactoryTest) TestCreateStringKey_colors() {
	astFactory := c.createTestAstFactoryWithColors()
	numberNode := astFactory.CreateNumber(2112.0)
	stringKeyNode := astFactory.CreateStringKey("key", numberNode)
	AssertNode(stringKeyNode).HasType(rhino.Token_STRING_KEY)
	AssertThat(stringKeyNode.GetString()).IsEqualTo("key")
	AssertThat(astFactoryTest_childList(stringKeyNode)).ContainsExactly(numberNode)
	AssertNode(stringKeyNode).HasColorThat().IsEqualTo(colors.StandardColors_NUMBER)
}

// Java: AstFactoryTest.testCreateComputedProperty_jstypes()
func (c *AstFactoryTest) TestCreateComputedProperty_jstypes() {
	astFactory := c.createTestAstFactory()
	numberNode := astFactory.CreateNumber(2112.0)
	stringLiteral := astFactory.CreateString("string literal key")
	computedPropertyNode := astFactory.CreateComputedProperty(stringLiteral, numberNode)
	AssertNode(computedPropertyNode).HasType(rhino.Token_COMPUTED_PROP)
	AssertThat(astFactoryTest_childList(computedPropertyNode)).
		ContainsExactly(stringLiteral, numberNode).
		InOrder()
	TypeSubject_AssertType(computedPropertyNode.GetJSType()).IsNumber()
}

// Java: AstFactoryTest.testCreateComputedProperty_colors()
func (c *AstFactoryTest) TestCreateComputedProperty_colors() {
	astFactory := c.createTestAstFactoryWithColors()
	numberNode := astFactory.CreateNumber(2112.0)
	stringLiteral := astFactory.CreateString("string literal key")
	computedPropertyNode := astFactory.CreateComputedProperty(stringLiteral, numberNode)
	AssertNode(computedPropertyNode).HasType(rhino.Token_COMPUTED_PROP)
	AssertThat(astFactoryTest_childList(computedPropertyNode)).
		ContainsExactly(stringLiteral, numberNode).
		InOrder()
	AssertNode(computedPropertyNode).HasColorThat().IsEqualTo(colors.StandardColors_NUMBER)
}

// Java: AstFactoryTest.testCreateGetterDef()
func (c *AstFactoryTest) TestCreateGetterDef() {
	astFactory := c.createTestAstFactory()
	valueNode := astFactory.CreateString("value")
	getterNode := astFactory.CreateGetterDef("name", valueNode)
	AssertNode(getterNode). //
				HasToken(rhino.Token_GETTER_DEF).
				HasStringThat().
				IsEqualTo("name")
	AssertNode(getterNode). //
				HasOneChildThat().
				IsFunction().
				HasXChildren(3)
	getterFunctionNode := getterNode.GetOnlyChild()
	// The first child is an empty string representing that the function has no name
	AssertNode(getterFunctionNode).HasFirstChildThat().IsName("")
	// the second child is the parameter list, which should be empty.
	AssertNode(getterFunctionNode). //
					HasSecondChildThat().
					IsParamList().
					HasNoChildren()
	getterFunctionBody := getterFunctionNode.GetLastChild()
	AssertNode(getterFunctionBody). //
					HasOneChildThat().
					HasToken(rhino.Token_RETURN).
					HasOneChildThat().
					IsEqualTo(valueNode)
}

// Java: AstFactoryTest.testCreateGetElem_jstypes()
func (c *AstFactoryTest) TestCreateGetElem_jstypes() {
	astFactory := c.createTestAstFactory()
	objectName := astFactory.CreateName("obj", AstFactory_Type_OfJSTypeNative(rhino.JSTypeNative_OBJECT_TYPE))
	stringLiteral := astFactory.CreateString("string literal key")
	getElemNode := astFactory.CreateGetElem(objectName, stringLiteral)
	AssertNode(getElemNode).HasType(rhino.Token_GETELEM)
	AssertThat(astFactoryTest_childList(getElemNode)).ContainsExactly(objectName, stringLiteral).InOrder()
	// TODO(bradfordcsmith): When receiver is an Array<T> or an Object<K, V>, use the template type
	// here.
	TypeSubject_AssertType(getElemNode.GetJSType()).IsUnknown()
}

// Java: AstFactoryTest.testCreateGetElem_colors()
func (c *AstFactoryTest) TestCreateGetElem_colors() {
	astFactory := c.createTestAstFactoryWithColors()
	objectName := astFactory.CreateName("obj", AstFactory_Type_OfJSTypeNativeAndColor(rhino.JSTypeNative_STRING_TYPE, colors.StandardColors_STRING))
	stringLiteral := astFactory.CreateString("string literal key")
	getElemNode := astFactory.CreateGetElem(objectName, stringLiteral)
	AssertNode(getElemNode).HasType(rhino.Token_GETELEM)
	AssertThat(astFactoryTest_childList(getElemNode)).ContainsExactly(objectName, stringLiteral).InOrder()
	AssertNode(getElemNode).HasColorThat().IsEqualTo(colors.StandardColors_UNKNOWN)
}

// Java: AstFactoryTest.testCreateComma_jstypes()
func (c *AstFactoryTest) TestCreateComma_jstypes() {
	astFactory := c.createTestAstFactory()
	stringNode := astFactory.CreateString("hi")
	numberNode := astFactory.CreateNumber(2112.0)
	commaNode := astFactory.CreateComma(stringNode, numberNode)
	AssertNode(commaNode).HasType(rhino.Token_COMMA)
	AssertThat(astFactoryTest_childList(commaNode)).ContainsExactly(stringNode, numberNode).InOrder()
	TypeSubject_AssertType(commaNode.GetJSType()).IsNumber()
}

// Java: AstFactoryTest.testCreateComma_colors()
func (c *AstFactoryTest) TestCreateComma_colors() {
	astFactory := c.createTestAstFactoryWithColors()
	stringNode := astFactory.CreateString("hi")
	numberNode := astFactory.CreateNumber(2112.0)
	commaNode := astFactory.CreateComma(stringNode, numberNode)
	AssertNode(commaNode).HasType(rhino.Token_COMMA)
	AssertThat(astFactoryTest_childList(commaNode)).ContainsExactly(stringNode, numberNode).InOrder()
	AssertNode(commaNode).HasColorThat().IsEqualTo(colors.StandardColors_NUMBER)
}

// Java: AstFactoryTest.testCreateCommas_jstypes()
func (c *AstFactoryTest) TestCreateCommas_jstypes() {
	astFactory := c.createTestAstFactory()
	stringNode := astFactory.CreateString("hi")
	numberNode := astFactory.CreateNumber(2112.0)
	trueNode := astFactory.CreateBoolean(true)
	falseNode := astFactory.CreateBoolean(false)
	// "hi", 2112, true, false
	stringNumberTrueFalse := astFactory.CreateCommas(stringNode, numberNode, trueNode, falseNode)
	// ("hi", 2112, true), false
	AssertNode(stringNumberTrueFalse).HasType(rhino.Token_COMMA)
	stringNumberTrue := stringNumberTrueFalse.GetFirstChild()
	AssertThat(astFactoryTest_childList(stringNumberTrueFalse)).
		ContainsExactly(stringNumberTrue, falseNode).
		InOrder()
	TypeSubject_AssertType(stringNumberTrueFalse.GetJSType()).IsBoolean()
	// ("hi", 2112), true
	stringNumber := stringNumberTrue.GetFirstChild()
	AssertThat(astFactoryTest_childList(stringNumberTrue)).ContainsExactly(stringNumber, trueNode).InOrder()
	TypeSubject_AssertType(stringNumberTrue.GetJSType()).IsBoolean()
	// "hi", 2112
	AssertThat(astFactoryTest_childList(stringNumber)).ContainsExactly(stringNode, numberNode)
	TypeSubject_AssertType(stringNumber.GetJSType()).IsNumber()
}

// Java: AstFactoryTest.testCreateIn_jstypes()
func (c *AstFactoryTest) TestCreateIn_jstypes() {
	astFactory := c.createTestAstFactory()
	prop := astFactory.CreateString("prop")
	obj := rhino.IR_Name("obj")
	// TODO(bradfordcsmith): This should have a type on it.
	n := astFactory.CreateIn(prop, obj)
	AssertNode(n).HasType(rhino.Token_IN)
	TypeSubject_AssertType(n.GetJSType()).IsBoolean()
	AssertThat(astFactoryTest_childList(n)).ContainsExactly(prop, obj).InOrder()
}

// Java: AstFactoryTest.testCreateIn_colors()
func (c *AstFactoryTest) TestCreateIn_colors() {
	astFactory := c.createTestAstFactoryWithColors()
	prop := astFactory.CreateString("prop")
	obj := rhino.IR_Name("obj")
	// TODO(bradfordcsmith): This should have a type on it.
	n := astFactory.CreateIn(prop, obj)
	AssertNode(n).HasType(rhino.Token_IN)
	AssertNode(n).HasColorThat().IsEqualTo(colors.StandardColors_BOOLEAN)
	AssertThat(astFactoryTest_childList(n)).ContainsExactly(prop, obj).InOrder()
}

// Java: AstFactoryTest.testCreateAnd_jstypes()
func (c *AstFactoryTest) TestCreateAnd_jstypes() {
	astFactory := c.createTestAstFactory()
	numberLiteral := astFactory.CreateNumber(2112.0)
	stringLiteral := astFactory.CreateString("hello")
	andNode := astFactory.CreateAnd(numberLiteral, stringLiteral)
	AssertNode(andNode).HasType(rhino.Token_AND)
	AssertThat(astFactoryTest_childList(andNode)).ContainsExactly(numberLiteral, stringLiteral).InOrder()
	TypeSubject_AssertType(andNode.GetJSType()).ToStringIsEqualTo("(number|string)")
}

// Java: AstFactoryTest.testCreateAnd_colors()
func (c *AstFactoryTest) TestCreateAnd_colors() {
	astFactory := c.createTestAstFactoryWithColors()
	numberLiteral := astFactory.CreateNumber(2112.0)
	stringLiteral := astFactory.CreateString("hello")
	andNode := astFactory.CreateAnd(numberLiteral, stringLiteral)
	AssertNode(andNode).HasType(rhino.Token_AND)
	AssertThat(astFactoryTest_childList(andNode)).ContainsExactly(numberLiteral, stringLiteral).InOrder()
	AssertNode(andNode).HasColorThat().HasAlternates(colors.StandardColors_STRING, colors.StandardColors_NUMBER)
}

// Java: AstFactoryTest.testCreateAndWithAlwaysFalsyLhs()
func (c *AstFactoryTest) TestCreateAndWithAlwaysFalsyLhs() {
	astFactory := c.createTestAstFactory()
	nullNode := astFactory.CreateNull()
	stringLiteral := astFactory.CreateString("hello")
	andNode := astFactory.CreateAnd(nullNode, stringLiteral)
	AssertNode(andNode).HasType(rhino.Token_AND)
	AssertThat(astFactoryTest_childList(andNode)).ContainsExactly(nullNode, stringLiteral).InOrder()
	// NULL_TYPE doesn't contain any truthy values, so its type is the only possibility
	// but AstFactory is simpler than the type inferencer and does not realize this.
	TypeSubject_AssertType(andNode.GetJSType()).ToStringIsEqualTo("(null|string)")
}

// Java: AstFactoryTest.testCreateAndWithAlwaysTruthyLhs()
func (c *AstFactoryTest) TestCreateAndWithAlwaysTruthyLhs() {
	astFactory := c.createTestAstFactory()
	nonNullObject := astFactory.CreateName("nonNullObject", AstFactory_Type_OfJSTypeNative(rhino.JSTypeNative_OBJECT_TYPE))
	stringLiteral := astFactory.CreateString("hello")
	andNode := astFactory.CreateAnd(nonNullObject, stringLiteral)
	AssertNode(andNode).HasType(rhino.Token_AND)
	AssertThat(astFactoryTest_childList(andNode)).ContainsExactly(nonNullObject, stringLiteral).InOrder()
	// OBJECT_TYPE doesn't contain any falsy values, so the RHS type is the only possibility
	// but AstFactory is simpler than the type inferencer and does not realize this.
	TypeSubject_AssertType(andNode.GetJSType()).ToStringIsEqualTo("(Object|string)")
}

// Java: AstFactoryTest.testCreateOr()
func (c *AstFactoryTest) TestCreateOr() {
	astFactory := c.createTestAstFactory()
	numberLiteral := astFactory.CreateNumber(2112.0)
	stringLiteral := astFactory.CreateString("hello")
	andNode := astFactory.CreateOr(numberLiteral, stringLiteral)
	AssertNode(andNode).HasType(rhino.Token_OR)
	AssertThat(astFactoryTest_childList(andNode)).ContainsExactly(numberLiteral, stringLiteral).InOrder()
	TypeSubject_AssertType(andNode.GetJSType()).ToStringIsEqualTo("(number|string)")
}

// Java: AstFactoryTest.testCreateOrWithAlwaysFalsyLhs()
func (c *AstFactoryTest) TestCreateOrWithAlwaysFalsyLhs() {
	astFactory := c.createTestAstFactory()
	nullLiteral := astFactory.CreateNull()
	stringLiteral := astFactory.CreateString("hello")
	andNode := astFactory.CreateOr(nullLiteral, stringLiteral)
	AssertNode(andNode).HasType(rhino.Token_OR)
	AssertThat(astFactoryTest_childList(andNode)).ContainsExactly(nullLiteral, stringLiteral).InOrder()
	// NULL_TYPE doesn't contain any truthy values, so the RHS type is the only possibility
	// but AstFactory is simpler than the type inferencer and does not realize this.
	TypeSubject_AssertType(andNode.GetJSType()).ToStringIsEqualTo("(null|string)")
}

// Java: AstFactoryTest.testCreateOrWithAlwaysTruthyLhs()
func (c *AstFactoryTest) TestCreateOrWithAlwaysTruthyLhs() {
	astFactory := c.createTestAstFactory()
	nonNullObject := astFactory.CreateName("nonNullObject", AstFactory_Type_OfJSTypeNative(rhino.JSTypeNative_OBJECT_TYPE))
	stringLiteral := astFactory.CreateString("hello")
	andNode := astFactory.CreateOr(nonNullObject, stringLiteral)
	AssertNode(andNode).HasType(rhino.Token_OR)
	AssertThat(astFactoryTest_childList(andNode)).ContainsExactly(nonNullObject, stringLiteral).InOrder()
	// OBJECT_TYPE doesn't contain any falsy values, so the RHS won't be evaluated
	// but AstFactory is simpler than the type inferencer and does not realize this.
	TypeSubject_AssertType(andNode.GetJSType()).ToStringIsEqualTo("(Object|string)")
}

// Java: AstFactoryTest.testCreateAdd_stringAndNumber_jstypes()
func (c *AstFactoryTest) TestCreateAdd_stringAndNumber_jstypes() {
	astFactory := c.createTestAstFactory()
	zero := astFactory.CreateNumber(0)
	str := astFactory.CreateString("x")
	add := astFactory.CreateAdd(zero, str)
	AssertNode(add).HasToken(rhino.Token_ADD)
	AssertThat(astFactoryTest_childList(add)).ContainsExactly(zero, str)
	AssertNode(add).HasJSTypeThat().IsEqualTo(c.getNativeType(rhino.JSTypeNative_BIGINT_NUMBER_STRING))
}

// Java: AstFactoryTest.testCreateAdd_stringAndNumber_colors()
func (c *AstFactoryTest) TestCreateAdd_stringAndNumber_colors() {
	astFactory := c.createTestAstFactoryWithColors()
	zero := astFactory.CreateNumber(0)
	str := astFactory.CreateString("x")
	add := astFactory.CreateAdd(zero, str)
	AssertNode(add).HasToken(rhino.Token_ADD)
	AssertThat(astFactoryTest_childList(add)).ContainsExactly(zero, str)
	AssertNode(add).
		HasColorThat().
		HasAlternates(colors.StandardColors_BIGINT, colors.StandardColors_NUMBER, colors.StandardColors_STRING)
}

// Java: AstFactoryTest.testCreateSub_jstypes()
func (c *AstFactoryTest) TestCreateSub_jstypes() {
	astFactory := c.createTestAstFactory()
	zero := astFactory.CreateNumber(0)
	one := astFactory.CreateNumber(1)
	sub := astFactory.CreateSub(zero, one)
	AssertNode(sub).HasToken(rhino.Token_SUB)
	AssertThat(astFactoryTest_childList(sub)).ContainsExactly(zero, one)
	AssertNode(sub).HasJSTypeThat().IsEqualTo(c.getNativeType(rhino.JSTypeNative_NUMBER_TYPE))
}

// Java: AstFactoryTest.testCreateSub_colors()
func (c *AstFactoryTest) TestCreateSub_colors() {
	astFactory := c.createTestAstFactoryWithColors()
	zero := astFactory.CreateNumber(0)
	one := astFactory.CreateNumber(1)
	sub := astFactory.CreateSub(zero, one)
	AssertNode(sub).HasToken(rhino.Token_SUB)
	AssertThat(astFactoryTest_childList(sub)).ContainsExactly(zero, one)
	AssertNode(sub).HasColorThat().IsEqualTo(colors.StandardColors_NUMBER)
}

// Java: AstFactoryTest.testCreateLessThan_jstypes()
func (c *AstFactoryTest) TestCreateLessThan_jstypes() {
	astFactory := c.createTestAstFactory()
	zero := astFactory.CreateNumber(0)
	one := astFactory.CreateNumber(1)
	lt := astFactory.CreateLessThan(zero, one)
	AssertNode(lt).HasToken(rhino.Token_LT)
	AssertThat(astFactoryTest_childList(lt)).ContainsExactly(zero, one)
	AssertNode(lt).HasJSTypeThat().IsEqualTo(c.getNativeType(rhino.JSTypeNative_BOOLEAN_TYPE))
}

// Java: AstFactoryTest.testCreateLessThan_colors()
func (c *AstFactoryTest) TestCreateLessThan_colors() {
	astFactory := c.createTestAstFactoryWithColors()
	zero := astFactory.CreateNumber(0)
	one := astFactory.CreateNumber(1)
	lt := astFactory.CreateLessThan(zero, one)
	AssertNode(lt).HasToken(rhino.Token_LT)
	AssertThat(astFactoryTest_childList(lt)).ContainsExactly(zero, one)
	AssertNode(lt).HasColorThat().IsEqualTo(colors.StandardColors_BOOLEAN)
}

// Java: AstFactoryTest.testCreateBitwiseAnd_jstypes()
func (c *AstFactoryTest) TestCreateBitwiseAnd_jstypes() {
	astFactory := c.createTestAstFactory()
	zero := astFactory.CreateNumber(0)
	one := astFactory.CreateNumber(1)
	bitAnd := astFactory.CreateBitwiseAnd(zero, one)
	AssertNode(bitAnd).HasToken(rhino.Token_BITAND)
	AssertThat(astFactoryTest_childList(bitAnd)).ContainsExactly(zero, one)
	AssertNode(bitAnd).HasJSTypeThat().IsEqualTo(c.getNativeType(rhino.JSTypeNative_NUMBER_TYPE))
}

// Java: AstFactoryTest.testCreateBitwiseAnd_colors()
func (c *AstFactoryTest) TestCreateBitwiseAnd_colors() {
	astFactory := c.createTestAstFactoryWithColors()
	zero := astFactory.CreateNumber(0)
	one := astFactory.CreateNumber(1)
	bitAnd := astFactory.CreateBitwiseAnd(zero, one)
	AssertNode(bitAnd).HasToken(rhino.Token_BITAND)
	AssertThat(astFactoryTest_childList(bitAnd)).ContainsExactly(zero, one)
	AssertNode(bitAnd).HasColorThat().IsEqualTo(colors.StandardColors_NUMBER)
}

// Java: AstFactoryTest.testCreateRightShift_jstypes()
func (c *AstFactoryTest) TestCreateRightShift_jstypes() {
	astFactory := c.createTestAstFactory()
	zero := astFactory.CreateNumber(0)
	one := astFactory.CreateNumber(1)
	rightShift := astFactory.CreateRightShift(zero, one)
	AssertNode(rightShift).HasToken(rhino.Token_RSH)
	AssertThat(astFactoryTest_childList(rightShift)).ContainsExactly(zero, one)
	AssertNode(rightShift).HasJSTypeThat().IsEqualTo(c.getNativeType(rhino.JSTypeNative_NUMBER_TYPE))
}

// Java: AstFactoryTest.testCreateRightShift_colors()
func (c *AstFactoryTest) TestCreateRightShift_colors() {
	astFactory := c.createTestAstFactoryWithColors()
	zero := astFactory.CreateNumber(0)
	one := astFactory.CreateNumber(1)
	rightShift := astFactory.CreateRightShift(zero, one)
	AssertNode(rightShift).HasToken(rhino.Token_RSH)
	AssertThat(astFactoryTest_childList(rightShift)).ContainsExactly(zero, one)
	AssertNode(rightShift).HasColorThat().IsEqualTo(colors.StandardColors_NUMBER)
}

// Java: AstFactoryTest.testCreateInc_prefix_jstypes()
func (c *AstFactoryTest) TestCreateInc_prefix_jstypes() {
	astFactory := c.createTestAstFactory()
	x := astFactory.CreateNameWithUnknownType("x")
	inc := astFactory.CreateInc(x /* isPost= */, false)
	AssertNode(inc).HasToken(rhino.Token_INC)
	AssertThat(inc.GetBooleanProp(rhino.Node_INCRDECR_PROP)).IsFalse()
	AssertThat(astFactoryTest_childList(inc)).ContainsExactly(x)
	AssertNode(inc).HasJSTypeThat().IsEqualTo(c.getNativeType(rhino.JSTypeNative_NUMBER_TYPE))
}

// Java: AstFactoryTest.testCreateInc_postfix_jstypes()
func (c *AstFactoryTest) TestCreateInc_postfix_jstypes() {
	astFactory := c.createTestAstFactory()
	x := astFactory.CreateNameWithUnknownType("x")
	inc := astFactory.CreateInc(x /* isPost= */, true)
	AssertNode(inc).HasToken(rhino.Token_INC)
	AssertThat(inc.GetBooleanProp(rhino.Node_INCRDECR_PROP)).IsTrue()
	AssertThat(astFactoryTest_childList(inc)).ContainsExactly(x)
	AssertNode(inc).HasJSTypeThat().IsEqualTo(c.getNativeType(rhino.JSTypeNative_NUMBER_TYPE))
}

// Java: AstFactoryTest.testCreateInc_prefix_colors()
func (c *AstFactoryTest) TestCreateInc_prefix_colors() {
	astFactory := c.createTestAstFactoryWithColors()
	x := astFactory.CreateNameWithUnknownType("x")
	inc := astFactory.CreateInc(x /* isPost= */, false)
	AssertNode(inc).HasToken(rhino.Token_INC)
	AssertThat(inc.GetBooleanProp(rhino.Node_INCRDECR_PROP)).IsFalse()
	AssertThat(astFactoryTest_childList(inc)).ContainsExactly(x)
	AssertNode(inc).HasColorThat().IsEqualTo(colors.StandardColors_NUMBER)
}

// Java: AstFactoryTest.testCreateCallWithTypeFromNode()
func (c *AstFactoryTest) TestCreateCallWithTypeFromNode() {
	astFactory := c.createTestAstFactory()
	root := c.parseAndAddTypes(
		`/**
 * @param {string} arg1
 * @param {number} arg2
 * @return {string}
 */
function foo() { return arg1; }
`)
	scope := c.getScope(root)
	// foo("hi", 2112)
	callee := astFactory.CreateNameInScope(scope, "foo")
	arg1 := astFactory.CreateString("hi")
	arg2 := astFactory.CreateNumber(2112.0)
	callNode := astFactory.CreateCall(callee, AstFactory_Type_(astFactory.CreateString("tmp")), arg1, arg2)
	AssertNode(callNode).HasType(rhino.Token_CALL)
	AssertThat(callNode.GetBooleanProp(rhino.Node_FREE_CALL)).IsTrue()
	AssertThat(astFactoryTest_childList(callNode)).ContainsExactly(callee, arg1, arg2).InOrder()
	TypeSubject_AssertType(callNode.GetJSType()).IsString()
}

// Java: AstFactoryTest.testCreateCallWithColorFromNode()
func (c *AstFactoryTest) TestCreateCallWithColorFromNode() {
	astFactory := c.createTestAstFactoryWithColors()
	root := c.parseAndAddColors(
		`/**
 * @param {string} arg1
 * @param {number} arg2
 * @return {string}
 */
function foo() { return arg1; }
`)
	scope := c.getScope(root)
	// foo("hi", 2112)
	callee := astFactory.CreateNameInScope(scope, "foo")
	arg1 := astFactory.CreateString("hi")
	arg2 := astFactory.CreateNumber(2112.0)
	callNode := astFactory.CreateCall(callee, AstFactory_Type_(astFactory.CreateString("tmp")), arg1, arg2)
	AssertNode(callNode).HasType(rhino.Token_CALL)
	AssertThat(callNode.GetBooleanProp(rhino.Node_FREE_CALL)).IsTrue()
	AssertThat(astFactoryTest_childList(callNode)).ContainsExactly(callee, arg1, arg2).InOrder()
	AssertNode(callNode).HasColorThat().IsEqualTo(colors.StandardColors_STRING)
}

// Java: AstFactoryTest.testCreateFreeCall()
func (c *AstFactoryTest) TestCreateFreeCall() {
	astFactory := c.createTestAstFactory()
	root := c.parseAndAddTypes(
		`/**
 * @param {string} arg1
 * @param {number} arg2
 * @return {string}
 */
function foo() { return arg1; }
`)
	scope := c.getScope(root)
	// foo("hi", 2112)
	callee := astFactory.CreateNameInScope(scope, "foo")
	arg1 := astFactory.CreateString("hi")
	arg2 := astFactory.CreateNumber(2112.0)
	callNode := astFactory.CreateCall(callee, AstFactory_Type_OfJSTypeNative(rhino.JSTypeNative_STRING_TYPE), arg1, arg2)
	AssertNode(callNode).HasType(rhino.Token_CALL)
	AssertThat(callNode.GetBooleanProp(rhino.Node_FREE_CALL)).IsTrue()
	AssertThat(astFactoryTest_childList(callNode)).ContainsExactly(callee, arg1, arg2).InOrder()
	TypeSubject_AssertType(callNode.GetJSType()).IsString()
}

// Java: AstFactoryTest.testCreateMethodCall_throws()
func (c *AstFactoryTest) TestCreateMethodCall_throws() {
	astFactory := c.createTestAstFactory()
	c.parseAndAddTypes(
		`class Foo {
  /**
   * @param {string} arg1
   * @param {number} arg2
   * @return {string}
   */
  method(arg1, arg2) { return arg1; }
}
const foo = new Foo();
`)
	scope := c.compiler.GetTranspilationNamespace()
	// createQName only accepts globally qualified qnames. foo.method is a prototype method access.
	// foo.method("hi", 2112)
	AssertThrows[jutil.Throwable](func() { astFactory.CreateQName(scope, "foo.method") })
}

// Java: AstFactoryTest.testCreateStaticMethodCallDotCallThrows()
func (c *AstFactoryTest) TestCreateStaticMethodCallDotCallThrows() {
	// NOTE: This method is testing both createCall() and createQName()
	astFactory := c.createTestAstFactory()
	c.parseAndAddTypes(
		`class Foo {
  /**
   * @param {string} arg1
   * @param {number} arg2
   * @return {string}
   */
  static method(arg1, arg2) { return arg1; }
}
`)
	scope := c.compiler.GetTranspilationNamespace()
	// createQName only accepts globally qualified qnames. While Foo.method is a global qualified
	// name, its '.call' property is not.
	// Foo.method.call(null, "hi", 2112)
	AssertThrows[jutil.Throwable](func() { astFactory.CreateQName(scope, "Foo.method.call") })
}

// Java: AstFactoryTest.testCreateQNameWithUnknownTypeFromString_jstype()
func (c *AstFactoryTest) TestCreateQNameWithUnknownTypeFromString_jstype() {
	astFactory := c.createTestAstFactory()
	objDotInnerDotStr := astFactory.CreateQNameWithUnknownType("obj.inner.str")
	AssertNode(objDotInnerDotStr).MatchesQualifiedName("obj.inner.str")
	objDotInner := objDotInnerDotStr.GetFirstChild()
	obj := objDotInner.GetFirstChild()
	AssertNode(obj).HasJSTypeThat().IsUnknown()
	AssertNode(objDotInner).HasJSTypeThat().IsUnknown()
	AssertNode(objDotInnerDotStr).HasJSTypeThat().IsUnknown()
}

// Java: AstFactoryTest.testCreateQNameWithUnknownTypeFromString_colors()
func (c *AstFactoryTest) TestCreateQNameWithUnknownTypeFromString_colors() {
	astFactory := c.createTestAstFactoryWithColors()
	objDotInnerDotStr := astFactory.CreateQNameWithUnknownType("obj.inner.str")
	AssertNode(objDotInnerDotStr).MatchesQualifiedName("obj.inner.str")
	objDotInner := objDotInnerDotStr.GetFirstChild()
	obj := objDotInner.GetFirstChild()
	AssertNode(obj).HasColorThat().IsEqualTo(colors.StandardColors_UNKNOWN)
	AssertNode(objDotInner).HasColorThat().IsEqualTo(colors.StandardColors_UNKNOWN)
	AssertNode(objDotInnerDotStr).HasColorThat().IsEqualTo(colors.StandardColors_UNKNOWN)
}

// Java: AstFactoryTest.testCreateQNameFromString()
func (c *AstFactoryTest) TestCreateQNameFromString() {
	astFactory := c.createTestAstFactory()
	c.parseAndAddTypes(
		`const obj = {
  inner: {
    str: 'hi',
  }
};
`)
	scope := c.compiler.GetTranspilationNamespace()
	objDotInnerDotStr := astFactory.CreateQName(scope, "obj.inner.str")
	AssertNode(objDotInnerDotStr).MatchesQualifiedName("obj.inner.str")
	objDotInner := objDotInnerDotStr.GetFirstChild()
	obj := objDotInner.GetFirstChild()
	AssertNode(obj).HasJSTypeThat().ToStringIsEqualTo("{inner: {str: string}}")
	AssertNode(objDotInner).HasJSTypeThat().ToStringIsEqualTo("{str: string}")
	AssertNode(objDotInnerDotStr).HasJSTypeThat().IsString()
}

// Java: AstFactoryTest.testCreateQNameFromBaseNamePlusStringIterable()
func (c *AstFactoryTest) TestCreateQNameFromBaseNamePlusStringIterable() {
	astFactory := c.createTestAstFactory()
	c.parseAndAddTypes(
		`const obj = {
  inner: {
    str: 'hi',
  }
};
`)
	scope := c.compiler.GetTranspilationNamespace()
	objDotInnerDotStr := astFactory.CreateQNameWithPropertyList(scope, "obj", []string{"inner", "str"})
	AssertNode(objDotInnerDotStr).MatchesQualifiedName("obj.inner.str")
	objDotInner := objDotInnerDotStr.GetFirstChild()
	obj := objDotInner.GetFirstChild()
	AssertNode(obj).HasJSTypeThat().ToStringIsEqualTo("{inner: {str: string}}")
	AssertNode(objDotInner).HasJSTypeThat().ToStringIsEqualTo("{str: string}")
	AssertNode(objDotInnerDotStr).HasJSTypeThat().IsString()
}

// Java: AstFactoryTest.testCreateNameMaintainsNormalization()
func (c *AstFactoryTest) TestCreateNameMaintainsNormalization() {
	astFactory := c.createTestAstFactory()
	root := c.parseAndAddTypes("const obj = {}")
	// Simulate normalization adding the IS_CONSTANT_NAME property ot `obj` NAME node
	objName := root.GetFirstChild(). // SCRIPT
						GetFirstChild(). // CONST
						GetFirstChild()
	// obj NAME
	AssertNode(objName).IsName("obj")
	objName.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
	newObjName := astFactory.CreateNameInScope(c.compiler.GetTranspilationNamespace(), "obj")
	// Assert that the newly created NAME node gets the IS_CONSTANT_NAME property it should in
	// order to be consistent with normalization.
	AssertThat(newObjName.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME)).IsTrue()
}

// Java: AstFactoryTest.testCreateQNameFromStringVarArgs()
func (c *AstFactoryTest) TestCreateQNameFromStringVarArgs() {
	astFactory := c.createTestAstFactory()
	c.parseAndAddTypes(
		`const obj = {
  inner: {
    str: 'hi',
  }
};
`)
	objDotInnerDotStr := astFactory.CreateQNameWithProperties(c.compiler.GetTranspilationNamespace(), "obj", "inner", "str")
	AssertNode(objDotInnerDotStr).MatchesQualifiedName("obj.inner.str")
	objDotInner := objDotInnerDotStr.GetFirstChild()
	obj := objDotInner.GetFirstChild()
	AssertNode(obj).HasJSTypeThat().ToStringIsEqualTo("{inner: {str: string}}")
	AssertNode(objDotInner).HasJSTypeThat().ToStringIsEqualTo("{str: string}")
	AssertNode(objDotInnerDotStr).HasJSTypeThat().IsString()
}

// Java: AstFactoryTest.testCreateQNameFromSimpleStringAndTypedScope()
func (c *AstFactoryTest) TestCreateQNameFromSimpleStringAndTypedScope() {
	astFactory := c.createTestAstFactory()
	scope := TypedScope_CreateGlobalScope(rhino.IR_Root())
	scope.Declare("x", rhino.IR_Name("x"), c.getNativeType(rhino.JSTypeNative_NUMBER_TYPE), nil, true)
	name := astFactory.CreateQNameUsingJSTypeInfo(scope, "x")
	AssertNode(name).HasStringThat().IsEqualTo("x")
	AssertNode(name).HasJSTypeThat().IsNumber()
}

// Java: AstFactoryTest.testCreateQNameFromDottedStringAndTypedScope()
func (c *AstFactoryTest) TestCreateQNameFromDottedStringAndTypedScope() {
	astFactory := c.createTestAstFactory()
	scope := TypedScope_CreateGlobalScope(rhino.IR_Root())
	// Declare a global "x" with the type "{y: number}".
	objectWithYProp := c.getRegistry().CreateAnonymousObjectType(nil)
	objectWithYProp.DefineDeclaredProperty("y", c.getNativeType(rhino.JSTypeNative_NUMBER_TYPE), rhino.NilNode)
	scope.Declare("x", rhino.IR_Name("x"), objectWithYProp, nil, true)
	name := astFactory.CreateQNameUsingJSTypeInfo(scope, "x.y")
	AssertNode(name).MatchesQualifiedName("x.y")
	AssertNode(name).HasJSTypeThat().IsNumber()
	AssertNode(name.GetFirstChild()).HasJSTypeThat().IsObjectTypeWithProperty("y")
}

// Java: AstFactoryTest.testCreateQNameFromStringAndTypedScope_crashesGivenMissingName()
func (c *AstFactoryTest) TestCreateQNameFromStringAndTypedScope_crashesGivenMissingName() {
	astFactory := c.createTestAstFactory()
	scope := TypedScope_CreateGlobalScope(rhino.IR_Root())
	AssertThrows[jutil.Throwable](func() { astFactory.CreateQNameUsingJSTypeInfo(scope, "x") })
}

// Java: AstFactoryTest.testCreateQNameFromStringAndTypedScope_crashesGivenLocalScope()
func (c *AstFactoryTest) TestCreateQNameFromStringAndTypedScope_crashesGivenLocalScope() {
	astFactory := c.createTestAstFactory()
	root := rhino.IR_Root()
	block := rhino.IR_Block()
	root.AddChildToFront(rhino.IR_ScriptOfNodes(block))
	globalScope := TypedScope_CreateGlobalScope(root)
	globalScope.Declare("x", rhino.IR_Name("x"), c.getNativeType(rhino.JSTypeNative_NUMBER_TYPE), nil, true)
	localScope := NewTypedScope(globalScope, block)
	unused := astFactory.CreateQNameUsingJSTypeInfo(globalScope, "x")
	_ = unused
	AssertThrows[*jutil.IllegalArgumentException](
		func() { astFactory.CreateQNameUsingJSTypeInfo(localScope, "x") })
}

// Java: AstFactoryTest.testCreateConstructorCall()
func (c *AstFactoryTest) TestCreateConstructorCall() {
	astFactory := c.createTestAstFactory()
	root := c.parseAndAddTypes(
		`class A {}
class B extends A {}
`)
	classBNode := root.GetFirstChild(). // script node
						GetSecondChild()
	classBExtendsNode := classBNode.GetSecondChild()
	classBType := classBNode.GetJSType().ToMaybeFunctionType()
	classBInstanceType := classBType.GetInstanceType()
	// simulate creating a call to super() intended to go in a constructor for B
	callee := rhino.IR_SuperNode().SetJSType(classBExtendsNode.GetJSType())
	arg1 := astFactory.CreateString("hi")
	arg2 := astFactory.CreateNumber(2112.0)
	callNode := astFactory.CreateConstructorCall(AstFactory_Type_(classBNode), callee, arg1, arg2)
	AssertNode(callNode).HasType(rhino.Token_CALL)
	AssertThat(callNode.GetBooleanProp(rhino.Node_FREE_CALL)).IsTrue()
	AssertThat(astFactoryTest_childList(callNode)).ContainsExactly(callee, arg1, arg2).InOrder()
	TypeSubject_AssertType(callNode.GetJSType()).IsEqualTo(classBInstanceType)
}

// Java: AstFactoryTest.testCreateEmptyFunction()
func (c *AstFactoryTest) TestCreateEmptyFunction() {
	astFactory := c.createTestAstFactory()
	// just a quick way to get a valid function type
	root := c.parseAndAddTypes("function foo() {}")
	functionType := root.GetFirstChild(). // script
						GetFirstChild(). // function
						GetJSType()
	emptyFunction := astFactory.CreateEmptyFunction(AstFactory_Type_OfJSType(functionType))
	AssertNode(emptyFunction).HasToken(rhino.Token_FUNCTION)
	TypeSubject_AssertType(emptyFunction.GetJSType()).IsEqualTo(functionType)
}

// Java: AstFactoryTest.testCreateFunction()
func (c *AstFactoryTest) TestCreateFunction() {
	astFactory := c.createTestAstFactory()
	// just a quick way to get a valid function type
	root := c.parseAndAddTypes("function foo() {}")
	functionType := root.GetFirstChild(). // script
						GetFirstChild(). // function
						GetJSType()
	paramList := rhino.IR_ParamList()
	body := rhino.IR_Block()
	functionNode := astFactory.CreateFunction("bar", paramList, body, AstFactory_Type_OfJSType(functionType))
	AssertNode(functionNode).HasToken(rhino.Token_FUNCTION)
	TypeSubject_AssertType(functionNode.GetJSType()).IsEqualTo(functionType)
	functionNameNode := functionNode.GetFirstChild()
	AssertNode(functionNameNode).IsName("bar")
	AssertThat(astFactoryTest_childList(functionNode)).
		ContainsExactly(functionNameNode, paramList, body).
		InOrder()
}

// Java: AstFactoryTest.testCreateMemberFunctionDef_jstypes()
func (c *AstFactoryTest) TestCreateMemberFunctionDef_jstypes() {
	astFactory := c.createTestAstFactory()
	// just a quick way to get a valid function type
	root := c.parseAndAddTypes("function foo() {}")
	functionType := root.GetFirstChild(). // script
						GetFirstChild(). // function
						GetJSType()
	paramList := rhino.IR_ParamList()
	body := rhino.IR_Block()
	functionNode := astFactory.CreateFunction("", paramList, body, AstFactory_Type_OfJSType(functionType))
	memberFunctionDef := astFactory.CreateMemberFunctionDef("bar", functionNode)
	AssertNode(memberFunctionDef).HasToken(rhino.Token_MEMBER_FUNCTION_DEF)
	AssertThat(memberFunctionDef.GetString()).IsEqualTo("bar")
	TypeSubject_AssertType(memberFunctionDef.GetJSType()).IsEqualTo(functionType)
}

// Java: AstFactoryTest.testCreateMemberFunctionDef_colors()
func (c *AstFactoryTest) TestCreateMemberFunctionDef_colors() {
	astFactory := c.createTestAstFactoryWithColors()
	// just a quick way to get a valid function type
	root := c.parseAndAddColors("function foo() {}")
	functionType := root.GetFirstChild(). // script
						GetFirstChild(). // function
						GetColor()
	paramList := rhino.IR_ParamList()
	body := rhino.IR_Block()
	functionNode := rhino.IR_Function(rhino.IR_Name(""), paramList, body).SetColor(functionType)
	memberFunctionDef := astFactory.CreateMemberFunctionDef("bar", functionNode)
	AssertNode(memberFunctionDef).HasToken(rhino.Token_MEMBER_FUNCTION_DEF)
	AssertThat(memberFunctionDef.GetString()).IsEqualTo("bar")
	AssertNode(memberFunctionDef).HasColorThat().IsEqualTo(functionType)
}

// Java: AstFactoryTest.testCreateZeroArgFunction()
func (c *AstFactoryTest) TestCreateZeroArgFunction() {
	astFactory := c.createTestAstFactory()
	// just a quick way to get a valid function type
	root := c.parseAndAddTypes("/** @return {number} */ function foo() {}")
	functionType := root.GetFirstChild(). // script
						GetFirstChild(). // function
						GetJSType()
	body := rhino.IR_Block()
	returnType := c.getNativeType(rhino.JSTypeNative_NUMBER_TYPE)
	functionNode := astFactory.CreateZeroArgFunction("bar", body, returnType)
	TypeSubject_AssertType(functionNode.GetJSType()).IsEqualTo(functionType)
	AssertNode(functionNode).HasToken(rhino.Token_FUNCTION)
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
//
// Java: AstFactoryTest.createZeroArgGeneratorFunction()
func (c *AstFactoryTest) TestCreateZeroArgGeneratorFunction() {
	astFactory := c.createTestAstFactory()
	// just a quick way to get a valid generator function type
	root := c.parseAndAddTypes("/** @return {number} */ function *foo() { return 1; }")
	functionType := root.GetFirstChild(). // script
						GetFirstChild(). // function
						GetJSType()
	body := rhino.IR_Block()
	returnType := c.getNativeType(rhino.JSTypeNative_NUMBER_TYPE)
	functionNode := astFactory.CreateZeroArgGeneratorFunction("bar", body, returnType)
	TypeSubject_AssertType(functionNode.GetJSType()).IsEqualTo(functionType)
	AssertNode(functionNode).HasToken(rhino.Token_FUNCTION)
}

// Java: AstFactoryTest.testCreateZeroArgFunctionForExpression()
func (c *AstFactoryTest) TestCreateZeroArgFunctionForExpression() {
	astFactory := c.createTestAstFactory()
	// quick way to get a function to contain the new arrow function and another arrow function
	// to compare types with
	root := c.parseAndAddTypes(
		`class C {
  /** @return {number} */
  foo() {
// TODO(b/118435472): compiler should be able to infer the return type
    /**
     * @return {number}
     */
    const orig = () => 1; // new arrow function exactly like this one
  }
}
`)
	existingArrowFunctionNode := root.GetFirstChild(). // script
								GetFirstChild(). // class
								GetLastChild().  // class members
								GetFirstChild(). // foo member function def
								GetFirstChild(). // foo function node
								GetLastChild().  // foo function body
								GetFirstChild(). // const
								GetOnlyChild().  // orig name node
								GetOnlyChild()
	expression := astFactory.CreateNumber(1)
	newArrowFunctionNode := astFactory.CreateZeroArgArrowFunctionForExpression(expression)
	AssertNode(newArrowFunctionNode).IsEqualTo(existingArrowFunctionNode)
	AssertNode(newArrowFunctionNode).
		HasJSTypeThat().
		IsEqualTo(existingArrowFunctionNode.GetJSTypeRequired())
}

// Java: AstFactoryTest.testCreateAssignFromNodes_jstypes()
func (c *AstFactoryTest) TestCreateAssignFromNodes_jstypes() {
	astFactory := c.createTestAstFactory()
	lhs := astFactory.CreateName("x", AstFactory_Type_OfJSTypeNative(rhino.JSTypeNative_STRING_TYPE))
	rhs := astFactory.CreateNumber(0)
	assign := astFactory.CreateAssign(lhs, rhs)
	AssertNode(assign).HasToken(rhino.Token_ASSIGN)
	AssertNode(assign).HasFirstChildThat().IsEqualTo(lhs)
	AssertNode(assign).HasSecondChildThat().IsEqualTo(rhs)
	AssertNode(assign).HasJSTypeThat().IsNumber()
	// take the rhs type, not lhs type
}

// Java: AstFactoryTest.testCreateAssignFromNodes_colors()
func (c *AstFactoryTest) TestCreateAssignFromNodes_colors() {
	astFactory := c.createTestAstFactoryWithColors()
	lhs := astFactory.CreateName("x", AstFactory_Type_OfColor(colors.StandardColors_STRING))
	rhs := astFactory.CreateNumber(0)
	assign := astFactory.CreateAssign(lhs, rhs)
	AssertNode(assign).HasToken(rhino.Token_ASSIGN)
	AssertNode(assign).HasFirstChildThat().IsEqualTo(lhs)
	AssertNode(assign).HasSecondChildThat().IsEqualTo(rhs)
	AssertNode(assign).HasColorThat().IsEqualTo(colors.StandardColors_NUMBER)
	// rhs type, not lhs type
}

// Java: AstFactoryTest.testCreateObjectLit_empty()
func (c *AstFactoryTest) TestCreateObjectLit_empty() {
	astFactory := c.createTestAstFactory()
	// just a quick way to get a valid object literal type
	root := c.parseAndAddTypes("({})")
	objectLitType := root.GetFirstChild(). // script
						GetFirstChild(). // expression result
						GetOnlyChild().  // object literal
						GetJSType()
	objectLit := astFactory.CreateObjectLit()
	TypeSubject_AssertType(objectLit.GetJSType()).ToStringIsEqualTo("{}")
	// isInstanceOf(objectLitType.getClass()): the same concrete JSType class.
	AssertThat(reflect.TypeOf(objectLit.GetJSType())).IsEqualTo(reflect.TypeOf(objectLitType))
	AssertNode(objectLit).HasToken(rhino.Token_OBJECTLIT)
	AssertNode(objectLit).HasChildren(false)
}

// Java: AstFactoryTest.testCreateObjectLit_empty_withColor()
func (c *AstFactoryTest) TestCreateObjectLit_empty_withColor() {
	astFactory := c.createTestAstFactoryWithColors()
	objColor := colors.Color_SingleBuilder().SetId(colors.ColorId_FromAscii("1")).Build()
	objectLit := astFactory.CreateObjectLitWithType(AstFactory_Type_OfJSTypeNativeAndColor(rhino.JSTypeNative_UNKNOWN_TYPE, objColor))
	AssertNode(objectLit).HasColorThat().IsEqualTo(objColor)
	AssertNode(objectLit).HasToken(rhino.Token_OBJECTLIT)
	AssertNode(objectLit).HasChildren(false)
}

// Java: AstFactoryTest.testCreateObjectLit_withElements()
func (c *AstFactoryTest) TestCreateObjectLit_withElements() {
	astFactory := c.createTestAstFactory()
	// just a quick way to get a valid object literal type
	root := c.parseAndAddTypes("({})")
	objectLitType := root.GetFirstChild(). // script
						GetFirstChild(). // expression result
						GetOnlyChild().  // object literal
						GetJSType()
	spread := rhino.IR_ObjectSpread(rhino.IR_Name("a"))
	stringKey := rhino.IR_StringKeyWithValue("b", rhino.IR_Number(0))
	objectLit := astFactory.CreateObjectLit(spread, stringKey)
	TypeSubject_AssertType(objectLit.GetJSType()).ToStringIsEqualTo("{}")
	// isInstanceOf(objectLitType.getClass()): the same concrete JSType class.
	AssertThat(reflect.TypeOf(objectLit.GetJSType())).IsEqualTo(reflect.TypeOf(objectLitType))
	AssertNode(objectLit).HasToken(rhino.Token_OBJECTLIT)
	AssertNode(objectLit.GetFirstChild()).IsSameInstanceAs(spread)
	AssertNode(objectLit.GetSecondChild()).IsSameInstanceAs(stringKey)
}

// Java: AstFactoryTest.testCreateObjectLit_withElementsAndType()
func (c *AstFactoryTest) TestCreateObjectLit_withElementsAndType() {
	astFactory := c.createTestAstFactory()
	// Try creating an object literal that implements Thenable and type it as one.
	thenableType := c.getRegistry().GetNativeType(rhino.JSTypeNative_THENABLE_TYPE)
	thenStringKey := rhino.IR_StringKeyWithValue(
		"then",
		astFactory.CreateZeroArgArrowFunctionForExpression(astFactory.CreateString("hi")))
	objectLit := astFactory.CreateObjectLitWithType(AstFactory_Type_OfJSType(thenableType), thenStringKey)
	AssertThat(objectLit.GetJSType()).IsEqualTo(thenableType)
	AssertNode(objectLit).HasToken(rhino.Token_OBJECTLIT)
	AssertNode(objectLit).HasOneChildThat().IsSameInstanceAs(thenStringKey)
}

// Java: AstFactoryTest.testCreateObjectLit_empty_colors()
func (c *AstFactoryTest) TestCreateObjectLit_empty_colors() {
	astFactory := c.createTestAstFactoryWithColors()
	objectLit := astFactory.CreateObjectLit()
	AssertNode(objectLit).HasColorThat().IsEqualTo(colors.StandardColors_TOP_OBJECT)
	AssertNode(objectLit).HasToken(rhino.Token_OBJECTLIT)
	AssertNode(objectLit).HasChildren(false)
}

// Java: AstFactoryTest.testCreateDelProp_jstypes()
func (c *AstFactoryTest) TestCreateDelProp_jstypes() {
	astFactory := c.createTestAstFactory()
	getprop := rhino.IR_Getprop(rhino.IR_Name("obj"), "prop")
	delprop := astFactory.CreateDelProp(getprop)
	AssertNode(delprop).HasToken(rhino.Token_DELPROP)
	TypeSubject_AssertType(delprop.GetJSType()).IsEqualTo(c.getNativeType(rhino.JSTypeNative_BOOLEAN_TYPE))
	AssertNode(delprop).HasChildren(true)
}

// Java: AstFactoryTest.testCreateDelProp_colors()
func (c *AstFactoryTest) TestCreateDelProp_colors() {
	astFactory := c.createTestAstFactoryWithColors()
	getprop := rhino.IR_Getprop(rhino.IR_Name("obj"), "prop")
	delprop := astFactory.CreateDelProp(getprop)
	AssertNode(delprop).HasToken(rhino.Token_DELPROP)
	AssertNode(delprop).HasColorThat().IsEqualTo(colors.StandardColors_BOOLEAN)
	AssertNode(delprop).HasChildren(true)
}

// Java: AstFactoryTest.testCreateSheq_jstypes()
func (c *AstFactoryTest) TestCreateSheq_jstypes() {
	astFactory := c.createTestAstFactory()
	left := rhino.IR_String("left")
	right := rhino.IR_Number(0)
	sheq := astFactory.CreateSheq(left, right)
	AssertNode(sheq).HasToken(rhino.Token_SHEQ)
	TypeSubject_AssertType(sheq.GetJSType()).IsEqualTo(c.getNativeType(rhino.JSTypeNative_BOOLEAN_TYPE))
}

// Java: AstFactoryTest.testCreateSheq_colors()
func (c *AstFactoryTest) TestCreateSheq_colors() {
	astFactory := c.createTestAstFactoryWithColors()
	left := rhino.IR_String("left")
	right := rhino.IR_Number(0)
	sheq := astFactory.CreateSheq(left, right)
	AssertNode(sheq).HasToken(rhino.Token_SHEQ)
	AssertNode(sheq).HasColorThat().IsEqualTo(colors.StandardColors_BOOLEAN)
}

// Java: AstFactoryTest.testCreateEq_jstypes()
func (c *AstFactoryTest) TestCreateEq_jstypes() {
	astFactory := c.createTestAstFactory()
	left := rhino.IR_String("left")
	right := rhino.IR_Number(0)
	sheq := astFactory.CreateEq(left, right)
	AssertNode(sheq).HasToken(rhino.Token_EQ)
	TypeSubject_AssertType(sheq.GetJSType()).IsEqualTo(c.getNativeType(rhino.JSTypeNative_BOOLEAN_TYPE))
}

// Java: AstFactoryTest.testCreateEq_colors()
func (c *AstFactoryTest) TestCreateEq_colors() {
	astFactory := c.createTestAstFactoryWithColors()
	left := rhino.IR_String("left")
	right := rhino.IR_Number(0)
	eq := astFactory.CreateEq(left, right)
	AssertNode(eq).HasToken(rhino.Token_EQ)
	AssertNode(eq).HasColorThat().IsEqualTo(colors.StandardColors_BOOLEAN)
}

// Java: AstFactoryTest.testCreateNe_jstypes()
func (c *AstFactoryTest) TestCreateNe_jstypes() {
	astFactory := c.createTestAstFactory()
	left := rhino.IR_String("left")
	right := rhino.IR_Number(0)
	ne := astFactory.CreateNe(left, right)
	AssertNode(ne).HasToken(rhino.Token_NE)
	TypeSubject_AssertType(ne.GetJSType()).IsEqualTo(c.getNativeType(rhino.JSTypeNative_BOOLEAN_TYPE))
}

// Java: AstFactoryTest.testCreateNe_colors()
func (c *AstFactoryTest) TestCreateNe_colors() {
	astFactory := c.createTestAstFactoryWithColors()
	left := rhino.IR_String("left")
	right := rhino.IR_Number(0)
	ne := astFactory.CreateNe(left, right)
	AssertNode(ne).HasToken(rhino.Token_NE)
	AssertNode(ne).HasColorThat().IsEqualTo(colors.StandardColors_BOOLEAN)
}

// Java: AstFactoryTest.testCreateHook_jstypes()
func (c *AstFactoryTest) TestCreateHook_jstypes() {
	astFactory := c.createTestAstFactory()
	stringType := c.getNativeType(rhino.JSTypeNative_STRING_TYPE)
	numberType := c.getNativeType(rhino.JSTypeNative_NUMBER_TYPE)
	condition := rhino.IR_FalseNode()
	left := rhino.IR_Name("left").SetJSType(stringType)
	right := rhino.IR_Number(0).SetJSType(numberType)
	hook := astFactory.CreateHook(condition, left, right)
	AssertNode(hook).HasToken(rhino.Token_HOOK)
	TypeSubject_AssertType(hook.GetJSType()).IsEqualTo(c.getRegistry().CreateUnionType(stringType, numberType))
}

// Java: AstFactoryTest.testCreateHook_colors()
func (c *AstFactoryTest) TestCreateHook_colors() {
	astFactory := c.createTestAstFactoryWithColors()
	condition := rhino.IR_FalseNode()
	left := astFactory.CreateString("left")
	right := astFactory.CreateNumber(0)
	hook := astFactory.CreateHook(condition, left, right)
	AssertNode(hook).HasToken(rhino.Token_HOOK)
	AssertNode(hook).HasColorThat().HasAlternates(colors.StandardColors_STRING, colors.StandardColors_NUMBER)
}

// Java: AstFactoryTest.testCreateArraylit_jstypes()
func (c *AstFactoryTest) TestCreateArraylit_jstypes() {
	// Given
	astFactory := c.createTestAstFactory()
	numberType := c.getNativeType(rhino.JSTypeNative_NUMBER_TYPE)
	first := rhino.IR_Number(0).SetJSType(numberType)
	second := rhino.IR_Number(1).SetJSType(numberType)
	third := rhino.IR_Number(2).SetJSType(numberType)
	expected := c.parseAndAddTypes("[0, 1, 2]").
		GetFirstChild(). // Script
		GetFirstChild(). // Expression
		GetFirstChild()
	// Array
	// When
	array := astFactory.CreateArraylit(first, second, third)
	// Then
	AssertNode(array).IsEquivalentTo(expected)
	TypeSubject_AssertType(array.GetJSType()).IsEqualTo(expected.GetJSType())
}

// Java: AstFactoryTest.testCreateArraylit_colors()
func (c *AstFactoryTest) TestCreateArraylit_colors() {
	// Given
	astFactory := c.createTestAstFactoryWithColors()
	first := astFactory.CreateNumber(0)
	second := astFactory.CreateNumber(1)
	third := astFactory.CreateNumber(2)
	expected := c.parseAndAddColors("[0, 1, 2]").
		GetFirstChild(). // Script
		GetFirstChild(). // Expression
		GetFirstChild()
	// Array
	// When
	array := astFactory.CreateArraylit(first, second, third)
	// Then
	AssertNode(array).IsEquivalentTo(expected)
	AssertThat(array.GetColor()).IsEqualTo(expected.GetColor())
}

// Java: AstFactoryTest.testCreateNewNode_jstypes()
func (c *AstFactoryTest) TestCreateNewNode_jstypes() {
	// Given
	astFactory := c.createTestAstFactory()
	numberType := c.getNativeType(rhino.JSTypeNative_NUMBER_TYPE)
	first := rhino.IR_Number(0).SetJSType(numberType)
	second := rhino.IR_Number(1).SetJSType(numberType)
	classNode := c.parseAndAddTypes(
		`class Example { constructor(arg0, arg1) {} }
new Example(0, 1);
`).
		GetFirstChild(). // Script
		GetFirstChild()
	// class
	expected := classNode.
		GetNext(). // ExpressionResult
		GetFirstChild()
	// NewExpression
	// When
	newExpr := astFactory.CreateNewNode(astFactory.CreateName("Example", AstFactory_Type_(classNode)), first, second)
	// Then
	AssertNode(newExpr).IsEquivalentTo(expected)
	TypeSubject_AssertType(newExpr.GetJSType()).IsEqualTo(expected.GetJSType())
}

// Java: AstFactoryTest.testCreateNewNode_colors()
func (c *AstFactoryTest) TestCreateNewNode_colors() {
	// Given
	astFactory := c.createTestAstFactoryWithColors()
	first := astFactory.CreateNumber(0)
	second := astFactory.CreateNumber(1)
	classNode := c.parseAndAddColors(
		`class Example { constructor(arg0, arg1) {} }
new Example(0, 1);
`).
		GetFirstChild(). // Script
		GetFirstChild()
	// class
	expected := classNode.
		GetNext(). // ExpressionResult
		GetFirstChild()
	// NewExpression
	// When
	newExpr := astFactory.CreateNewNode(
		astFactory.CreateName("Example", AstFactory_Type_OfColor(classNode.GetColor())), first, second)
	// Then
	AssertNode(newExpr).IsEquivalentTo(expected)
	AssertNode(newExpr).HasColorThat().IsEqualTo(expected.GetColor())
}

// Java: AstFactoryTest.testCreateRuntimeField_throwsIfFieldNotInjected()
func (c *AstFactoryTest) TestCreateRuntimeField_throwsIfFieldNotInjected() {
	// Given
	astFactory := c.createTestAstFactoryWithoutTypes()
	scope := rhinotesting.NewMapBasedScope(jutil.NewLinkedHashMap[string, rhino.JSType]())
	field := c.runtimeJsLibManager.GetJsLibField("$jscomp.global")
	// When
	AssertThrows[*jutil.IllegalStateException](func() { astFactory.CreateQNameOfField(scope, field) })
}

// Java: AstFactoryTest.testCreateRuntimeField_succeeds()
func (c *AstFactoryTest) TestCreateRuntimeField_succeeds() {
	// Given
	astFactory := c.createTestAstFactoryWithoutTypes()
	scope := rhinotesting.NewMapBasedScope(jutil.NewLinkedHashMap[string, rhino.JSType]())
	field := c.runtimeJsLibManager.GetJsLibField("$jscomp.global")
	c.runtimeJsLibManager.InjectLibForField("$jscomp.global")
	// When
	result := astFactory.CreateQNameOfField(scope, field)
	// Then
	AssertNode(result).MatchesQualifiedName("$jscomp.global")
}

// Java: AstFactoryTest.testCreateRuntimeField_inExternMode_succeeds_andUsesUnderscoredName()
func (c *AstFactoryTest) TestCreateRuntimeField_inExternMode_succeeds_andUsesUnderscoredName() {
	// Given
	externs := rhino.IR_Script()
	externs.SetStaticSourceFile(SourceFile_FromCodeWithKind("externs.js", "", rhino.StaticSourceFile_SourceKind_EXTERN))
	c.runtimeJsLibManager =
		RuntimeJsLibManager_Create(
			RuntimeJsLibManager_RuntimeLibraryMode_EXTERN_FIELD_NAMES,
			astFactoryTest_TestResourceProvider,
			c.compiler.GetChangeTracker(),
			func() rhino.Node { return externs })
	astFactory := c.createTestAstFactoryWithoutTypes()
	scope := rhinotesting.NewMapBasedScope(jutil.NewLinkedHashMap[string, rhino.JSType]())
	field := c.runtimeJsLibManager.GetJsLibField("$jscomp.global")
	c.runtimeJsLibManager.InjectLibForField("$jscomp.global")
	// When
	result := astFactory.CreateQNameOfField(scope, field)
	// Then
	AssertNode(result).MatchesName("$jscomp_global")
}

// Java: AstFactoryTest.testCreateJscompMakeIteratorCall_throwsIfJscompMakeIteratorNotInjected()
func (c *AstFactoryTest) TestCreateJscompMakeIteratorCall_throwsIfJscompMakeIteratorNotInjected() {
	// Given
	astFactory := c.createTestAstFactoryWithoutTypes()
	iterable := rhino.IR_Name("arr")
	scope := rhinotesting.NewMapBasedScope(jutil.NewLinkedHashMap[string, rhino.JSType]())
	// When
	AssertThrows[*jutil.IllegalStateException](
		func() { astFactory.CreateJSCompMakeIteratorCall(iterable, scope) })
}

// Java: AstFactoryTest.testCreateJscompMakeIteratorCall_succeeds()
func (c *AstFactoryTest) TestCreateJscompMakeIteratorCall_succeeds() {
	// Given
	astFactory := c.createTestAstFactoryWithoutTypes()
	iterable := rhino.IR_Name("arr")
	c.runtimeJsLibManager.InjectLibForField("$jscomp.makeIterator")
	// When
	result := astFactory.CreateJSCompMakeIteratorCall(iterable, rhinotesting.MapBasedScope_EmptyScope())
	// Then
	AssertNode(result).IsCall()
	AssertNode(result).HasFirstChildThat().MatchesQualifiedName("$jscomp.makeIterator")
	AssertNode(result).HasLastChildThat().IsEqualTo(iterable)
}

// Java: AstFactoryTest.testCreateJscompIteratorCloseCall_throwsIfJscompIteratorCloseNotInjected()
func (c *AstFactoryTest) TestCreateJscompIteratorCloseCall_throwsIfJscompIteratorCloseNotInjected() {
	// Given
	astFactory := c.createTestAstFactoryWithoutTypes()
	iter := rhino.IR_Name("iter")
	iterResult := rhino.IR_Name("res")
	scope := rhinotesting.NewMapBasedScope(jutil.NewLinkedHashMap[string, rhino.JSType]())

	// When
	AssertThrows[*jutil.IllegalStateException](
		func() { astFactory.CreateJscompIteratorCloseCall(iter, iterResult, scope) })
}

// Java: AstFactoryTest.testCreateJscompIteratorCloseCall_succeeds()
func (c *AstFactoryTest) TestCreateJscompIteratorCloseCall_succeeds() {
	// Given
	astFactory := c.createTestAstFactoryWithoutTypes()
	iter := rhino.IR_Name("iter")
	iterResult := rhino.IR_Name("res")
	c.runtimeJsLibManager.InjectLibForField("$jscomp.iteratorClose")

	// When
	result :=
		astFactory.CreateJscompIteratorCloseCall(iter, iterResult, rhinotesting.MapBasedScope_EmptyScope())

	// Then
	AssertNode(result).IsCall()
	AssertNode(result).HasFirstChildThat().MatchesQualifiedName("$jscomp.iteratorClose")
	AssertNode(result).HasSecondChildThat().IsEqualTo(iter)
	AssertNode(result).HasLastChildThat().IsEqualTo(iterResult)
}

// Java: AstFactoryTest.childList(Node)
func astFactoryTest_childList(parent rhino.Node) []rhino.Node {
	list := []rhino.Node{}
	for child := parent.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		list = append(list, child)
	}
	return list
}
