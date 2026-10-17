// Port of com/google/javascript/jscomp/AstValidatorTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/AstValidatorTest.java)

/*
 * Copyright 2011 The Closure Compiler Authors.
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
	"iter"
	"math/big"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// AstValidatorTest holds tests for AstValidator.
//
// Go-only: test-helper — Java test methods not named test* are
// Test_<javaName> so that RunJUnit4 finds them.
type AstValidatorTest struct {
	CompilerTestCaseBase

	lastCheckViolationMessages *jutil.ArrayList[string]
	typeInfoValidationMode     AstValidator_TypeInfoValidation
}

func NewAstValidatorTest() *AstValidatorTest {
	c := &AstValidatorTest{typeInfoValidationMode: AstValidator_TypeInfoValidation_JSTYPE}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestAstValidatorTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewAstValidatorTest() })
}

// Java: AstValidatorTest.getProcessor(Compiler)
func (c *AstValidatorTest) GetProcessor(compiler *Compiler) CompilerPass {
	return c.createValidator(compiler)
}

// Java: AstValidatorTest.getOptions()
func (c *AstValidatorTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetWarningLevel(DiagnosticGroups_MODULE_LOAD, CheckLevel_OFF)
	return options
}

// astValidatorTest_handler is the anonymous ViolationHandler of
// createValidator.
type astValidatorTest_handler struct {
	messages *jutil.ArrayList[string]
}

// Java: AstValidatorTest.createValidator(Compiler).new ViolationHandler() {...}.handleViolation(String, Node)
func (h *astValidatorTest_handler) HandleViolation(message string, n *rhino.Node) {
	h.messages.Add(message)
}

// Java: AstValidatorTest.createValidator(Compiler)
func (c *AstValidatorTest) createValidator(compiler *Compiler) *AstValidator {
	c.lastCheckViolationMessages = jutil.NewArrayList[string]()
	astValidator := NewAstValidatorWithHandler(
		compiler,
		&astValidatorTest_handler{messages: c.lastCheckViolationMessages},
		/* validateScriptFeatures= */ true,
		/* shouldValidateRequiredInlinings= */ true)
	astValidator.SetTypeValidationMode(c.typeInfoValidationMode)
	return astValidator
}

// Java: AstValidatorTest.setUp()
func (c *AstValidatorTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.DisableAstValidation()
	c.DisableNormalize()
	c.EnableTypeCheck()
}

// Java: AstValidatorTest.testIsParenthesizedProperty()
func (c *AstValidatorTest) TestIsParenthesizedProperty() {

	// Since we're building the AST by hand, there won't be any types on it.
	c.typeInfoValidationMode = AstValidator_TypeInfoValidation_NONE

	n := rhino.IR_String("a")
	n.SetIsParenthesized(true)
	c.setTestSourceLocationForTree(n)

	c.expectValid(n, astValidatorTest_Check_EXPRESSION)

	n.SetToken(rhino.Token_STRING_KEY) // A string key cannot be parenthesized
	// We have to put the STRING_KEY into an object and give it a child, so we have an
	// expression to validate that is valid other than the bad parenthesized property.
	objNode := rhino.IR_Objectlit(n)
	n.AddChildToFront(rhino.IR_Number(0))
	objNode.SrcrefTree(n)
	c.expectInvalid(objNode, astValidatorTest_Check_EXPRESSION, "non-expression is parenthesized")
}

// Java: AstValidatorTest.testValidGetPropAndGetElem()
func (c *AstValidatorTest) TestValidGetPropAndGetElem() {

	c.valid(
		`a.b;
a['b'];
`)
}

// Java: AstValidatorTest.testClass()
func (c *AstValidatorTest) TestClass() {

	c.valid(
		`class C {
  get m1() {return 1}
  set m1(a) {}
  m2(a) {}
}

/** @dict */
class D {
  ['m2']() {}
  [m2]() {}
}
`)

	c.typeInfoValidationMode = AstValidator_TypeInfoValidation_NONE // synthetic AST w/o types

	cls := rhino.NewNodeWith2Children(rhino.Token_CLASS, rhino.IR_Name("C"), rhino.IR_Empty())
	members := rhino.NewNode(rhino.Token_CLASS_MEMBERS)
	cls.AddChildToBack(members)
	c.setTestSourceLocationForTree(cls)
	c.expectValid(cls, astValidatorTest_Check_STATEMENT)
	method1 :=
		rhino.NewNodeWithChild(rhino.Token_MEMBER_FUNCTION_DEF, rhino.IR_Function(rhino.IR_Name(""), rhino.IR_ParamList(), rhino.IR_Block()))
	method1.SrcrefTree(members)
	members.AddChildToBack(method1)
	c.expectInvalid(cls, astValidatorTest_Check_STATEMENT)

	members.DetachChildren()

	// Invalid empty string
	method2 := rhino.Node_NewStringWithToken(rhino.Token_MEMBER_FUNCTION_DEF, "")
	method2.AddChildToBack(rhino.IR_Function(rhino.IR_Name(""), rhino.IR_ParamList(), rhino.IR_Block()))
	method2.SrcrefTree(members)
	members.AddChildToBack(method2)

	c.expectInvalid(cls, astValidatorTest_Check_STATEMENT)
}

// Java: AstValidatorTest.testClassField()
func (c *AstValidatorTest) TestClassField() {

	c.valid("class C {x}")
	c.valid("class C {x = 2}")
	c.valid("class C {x = 2;}")
	c.valid("class C {x; y;}")
	c.valid(
		`class C {
  x
  y
}
`)
}

// Java: AstValidatorTest.testClassFieldStatic()
func (c *AstValidatorTest) TestClassFieldStatic() {

	c.valid("class C {static x}")
	c.valid("class C {static x = 2}")
	c.valid("class C {static x = 2;}")
	c.valid("class C {static x; static y;}")
	c.valid(
		`class C {
  static x
  static y
}
`)
}

// Java: AstValidatorTest.testClassField_super()
func (c *AstValidatorTest) TestClassField_super() {

	c.valid(
		`class B { static x = 2}
class C extends B {
  static y = super.x;
}
`)
	c.valid(
		`class B { static x = 2}
class C extends B {
  static y = super.x;
}
`)
	c.valid(
		`class B { static x = 2}
class C extends B {
  static y = super.x + 2;
}
`)
	c.valid(
		`class B { static x = 2 }
class C extends B {
   static y = super.x + 2 + 4 + 6 - 8;
}
`)
	c.valid(
		`class B { static x = 2 }
/** @unrestricted */ class C extends B {
   static ['y'] = super.x;
}
`)
	c.valid(
		`class C { x=1;}
class D extends C {
  y = () => super.x;
}
`)
	c.invalid(
		`class C { x=1;}
class D extends C {
  y = function() { super.x; }
}
`)
}

// Java: AstValidatorTest.testClassComputedField()
func (c *AstValidatorTest) TestClassComputedField() {

	c.valid("/** @dict */ class C { [x]; }")
	c.valid("/** @dict */ class C { ['x']=2; }")
	c.valid("/** @dict */ class C { 'x'=2; }")
	c.valid("/** @dict */ class C { 1=2; }")
	c.valid(
		`/** @unrestricted */
class C {
  [x]=2
  static y = 4
}
`)
}

// Java: AstValidatorTest.testClassComputedFieldStatic()
func (c *AstValidatorTest) TestClassComputedFieldStatic() {

	c.valid("/** @dict */ class C { static [x]; }")
	c.valid("/** @dict */ class C { static ['x']=2; }")
	c.valid("/** @dict */ class C { static 'x'=2; }")
	c.valid("/** @dict */ class C { static 1=2; }")
	c.valid(
		`/** @unrestricted */
class C {
  static [x]=2
  static y = 4
}
`)
}

// Java: AstValidatorTest.testFeatureValidation_classField()
func (c *AstValidatorTest) TestFeatureValidation_classField() {

	c.testFeatureValidation(
		`class C {
  x=2;
}
`,
		parser.FeatureSet_Feature_PUBLIC_CLASS_FIELDS)
}

// Java: AstValidatorTest.testFeatureValidation_privateClassField()
func (c *AstValidatorTest) TestFeatureValidation_privateClassField() {
	c.testFeatureValidation(
		`class C {
  #x=2;
}
`,
		parser.FeatureSet_Feature_PRIVATE_ELEMENTS)
}

// Java: AstValidatorTest.testFeatureValidation_staticPrivateClassField()
func (c *AstValidatorTest) TestFeatureValidation_staticPrivateClassField() {
	c.testFeatureValidation(
		`class C {
  static #x=2;
}
`,
		parser.FeatureSet_Feature_PRIVATE_ELEMENTS)
}

// Java: AstValidatorTest.testFeatureValidation_privateGetter()
func (c *AstValidatorTest) TestFeatureValidation_privateGetter() {
	c.testFeatureValidation(
		`class C {
  get #x() { return 1; }
}
`,
		parser.FeatureSet_Feature_PRIVATE_ELEMENTS)
}

// Java: AstValidatorTest.testFeatureValidation_privateSetter()
func (c *AstValidatorTest) TestFeatureValidation_privateSetter() {
	c.testFeatureValidation(
		`class C {
  set #x(val) {}
}
`,
		parser.FeatureSet_Feature_PRIVATE_ELEMENTS)
}

// Java: AstValidatorTest.testFeatureValidation_privateMethod()
func (c *AstValidatorTest) TestFeatureValidation_privateMethod() {
	c.testFeatureValidation(
		`class C {
  #m() {}
}
`,
		parser.FeatureSet_Feature_PRIVATE_ELEMENTS)
}

// Java: AstValidatorTest.testFeatureValidation_privateBrandCheck()
func (c *AstValidatorTest) TestFeatureValidation_privateBrandCheck() {
	c.testFeatureValidation(
		`class C {
  #x;
  check(obj) {
    return #x in obj;
  }
}
`,
		parser.FeatureSet_Feature_PRIVATE_ELEMENTS)
}

// Java: AstValidatorTest.testFeatureValidation_staticPrivateMethod()
func (c *AstValidatorTest) TestFeatureValidation_staticPrivateMethod() {
	c.testFeatureValidation(
		`class C {
  static #m() {}
}
`,
		parser.FeatureSet_Feature_PRIVATE_ELEMENTS)
}

// Java: AstValidatorTest.testFeatureValidation_staticPrivateGetter()
func (c *AstValidatorTest) TestFeatureValidation_staticPrivateGetter() {
	c.testFeatureValidation(
		`class C {
  static get #x() { return 1; }
}
`,
		parser.FeatureSet_Feature_PRIVATE_ELEMENTS)
}

// Java: AstValidatorTest.testFeatureValidation_staticPrivateSetter()
func (c *AstValidatorTest) TestFeatureValidation_staticPrivateSetter() {
	c.testFeatureValidation(
		`class C {
  static set #x(val) {}
}
`,
		parser.FeatureSet_Feature_PRIVATE_ELEMENTS)
}

// Java: AstValidatorTest.testFeatureValidation_classComputedField()
func (c *AstValidatorTest) TestFeatureValidation_classComputedField() {

	c.testFeatureValidation(
		`/** @dict */
class C {
  [x]=2;
}
`,
		parser.FeatureSet_Feature_PUBLIC_CLASS_FIELDS)
}

// Java: AstValidatorTest.testClassStaticBlock()
func (c *AstValidatorTest) TestClassStaticBlock() {

	c.valid(
		`class C {
  static {
  }
}
`)
}

// Java: AstValidatorTest.testClassStaticBlock_this()
func (c *AstValidatorTest) TestClassStaticBlock_this() {

	// TODO (user): re-enable TypeInfoValidation and TypeCheck
	c.typeInfoValidationMode = AstValidator_TypeInfoValidation_NONE
	c.DisableTypeCheck()
	c.valid(
		`class C {
  static {
    this.prop = 4;
  }
}
`)
}

// Java: AstValidatorTest.testClassStaticBlock_superProperty()
func (c *AstValidatorTest) TestClassStaticBlock_superProperty() {

	// TODO (user): re-enable TypeInfoValidation and TypeCheck
	c.typeInfoValidationMode = AstValidator_TypeInfoValidation_NONE
	c.DisableTypeCheck()
	c.valid(
		`class D {
  static field;
}
class C extends D {
  static {
    if (Foo) { super.field = 'hello'; }
  }
}
`)
}

// Java: AstValidatorTest.testClassStaticBlock_invalidSuperConstructorCall()
func (c *AstValidatorTest) TestClassStaticBlock_invalidSuperConstructorCall() {

	// TODO (user): re-enable TypeInfoValidation and TypeCheck
	c.typeInfoValidationMode = AstValidator_TypeInfoValidation_NONE
	c.DisableTypeCheck()
	c.invalid(
		`class D {}
class C extends D {
  static {
    super();
  }
}
`)
}

// Java: AstValidatorTest.testFeatureValidation_classStaticBlock()
func (c *AstValidatorTest) TestFeatureValidation_classStaticBlock() {

	c.testFeatureValidation(
		`class C {
  static {
  }
}
`,
		parser.FeatureSet_Feature_CLASS_STATIC_BLOCK)
}

// Java: AstValidatorTest.testFor()
func (c *AstValidatorTest) TestFor() {

	c.valid("for(var a;;);")
	c.valid("for(let a;;);")
	c.valid("for(var a = 0;;);")
	c.valid("for(var a;b;c);")
	c.valid("for(var a;;c);")
}

// Java: AstValidatorTest.testForIn()
func (c *AstValidatorTest) TestForIn() {

	c.valid("for(var a in b);")
	c.valid("for(let a in b);")
	c.valid("for(const a in b);")
	c.valid("for(a in b);")
	c.valid("for(a in []);")
	c.valid("for(a in {});")

	// Test that initializers are banned (except for simple vars - see testQuestionableForIn)
	c.typeInfoValidationMode = AstValidator_TypeInfoValidation_NONE
	c.expectInvalid(
		rhino.NewNodeWith2Children(rhino.Token_FOR_IN, rhino.IR_ConstNode(rhino.IR_Name("a"), rhino.IR_Number(1)), rhino.IR_Name("b")),
		astValidatorTest_Check_STATEMENT)
	c.expectInvalid(
		rhino.NewNodeWith2Children(
			rhino.Token_FOR_IN,
			rhino.IR_Var(rhino.NewNodeWith2Children(rhino.Token_DESTRUCTURING_LHS, rhino.IR_ObjectPattern(), rhino.IR_Objectlit())),
			rhino.IR_Name("b")),
		astValidatorTest_Check_STATEMENT)
}

// Java: AstValidatorTest.testForOf()
func (c *AstValidatorTest) TestForOf() {

	c.valid("for(var a of b);")
	c.valid("for(let a of b);")
	c.valid("for(const a of b);")
	c.valid("for(a of b);")
	c.valid("for(a of []);")
	c.valid("for(a of /** @type {!Iterable<?>} */ ({}));")
	c.valid("for (const [] of b);")
	c.valid("for (const {} of b);")

	// Test that initializers are banned
	c.typeInfoValidationMode = AstValidator_TypeInfoValidation_NONE
	c.expectInvalid(
		rhino.NewNodeWith2Children(rhino.Token_FOR_OF, rhino.IR_VarWithValue(rhino.IR_Name("a"), rhino.IR_Number(1)), rhino.IR_Name("b")), astValidatorTest_Check_STATEMENT)
	c.expectInvalid(
		rhino.NewNodeWith2Children(rhino.Token_FOR_OF, rhino.IR_ConstNode(rhino.IR_Name("a"), rhino.IR_Number(1)), rhino.IR_Name("b")),
		astValidatorTest_Check_STATEMENT)
	c.expectInvalid(
		rhino.NewNodeWith2Children(
			rhino.Token_FOR_OF,
			rhino.IR_Var(rhino.NewNodeWith2Children(rhino.Token_DESTRUCTURING_LHS, rhino.IR_ObjectPattern(), rhino.IR_Objectlit())),
			rhino.IR_Name("b")),
		astValidatorTest_Check_STATEMENT)
}

// Java: AstValidatorTest.testForAwaitOf()
func (c *AstValidatorTest) TestForAwaitOf() {

	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	c.valid("async () => { for await(var a of b); }")
	c.valid("async () => { for await(let a of b); }")
	c.valid("async () => { for await(const a of b); }")
	c.valid("async () => { for await(a of b); }")
	c.valid("async () => { for await(a of []); }")
	c.valid("async () => { for await(a of /** @type {!Iterable<?>} */ ({})); }")

	// For tests, we permit top-level await and just verify that the initializer is banned.
	c.typeInfoValidationMode = AstValidator_TypeInfoValidation_NONE

	wrapInEs6Script :=
		func(n *rhino.Node) *rhino.Node {
			script := rhino.IR_ScriptOfNodes(n)
			script.PutBooleanProp(rhino.Node_ES6_MODULE, true)
			c.setTestSourceLocationForTree(script)
			return n
		}

	c.expectInvalid(
		wrapInEs6Script(
			rhino.NewNodeWith2Children(rhino.Token_FOR_AWAIT_OF, rhino.IR_VarWithValue(rhino.IR_Name("a"), rhino.IR_Number(1)), rhino.IR_Name("b"))),
		astValidatorTest_Check_STATEMENT)
	c.expectInvalid(
		wrapInEs6Script(
			rhino.NewNodeWith2Children(rhino.Token_FOR_AWAIT_OF, rhino.IR_ConstNode(rhino.IR_Name("a"), rhino.IR_Number(1)), rhino.IR_Name("b"))),
		astValidatorTest_Check_STATEMENT)
	c.expectInvalid(
		wrapInEs6Script(
			rhino.NewNodeWith2Children(
				rhino.Token_FOR_AWAIT_OF,
				rhino.IR_Var(rhino.NewNodeWith2Children(rhino.Token_DESTRUCTURING_LHS, rhino.IR_ObjectPattern(), rhino.IR_Objectlit())),
				rhino.IR_Name("b"))),
		astValidatorTest_Check_STATEMENT)
}

// Java: AstValidatorTest.testTopLevelForAwaitOf()
func (c *AstValidatorTest) TestTopLevelForAwaitOf() {

	// Since we're building the AST by hand, there won't be any types on it.
	c.typeInfoValidationMode = AstValidator_TypeInfoValidation_NONE

	c.SetLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT, CompilerOptions_LanguageMode_ECMASCRIPT5)

	forAwaitOfNode := rhino.NewNode(rhino.Token_FOR_AWAIT_OF)
	forAwaitOfNode.AddChildToBack(rhino.IR_Name("x"))
	forAwaitOfNode.AddChildToBack(rhino.IR_Name("y"))
	forAwaitOfNode.AddChildToBack(rhino.IR_Block())

	script := rhino.IR_ScriptOfNodes(forAwaitOfNode)
	c.setTestSourceLocationForTree(script)

	c.expectInvalid(
		forAwaitOfNode, astValidatorTest_Check_STATEMENT, "'await' expression is not within an async function")

	// Top-level await is only allowed in modules.
	script.PutBooleanProp(rhino.Node_ES6_MODULE, true)
	c.expectValid(forAwaitOfNode, astValidatorTest_Check_STATEMENT)
}

// Java: AstValidatorTest.testIncDec()
func (c *AstValidatorTest) TestIncDec() {

	c.valid("x++")
	c.valid("x--")
	c.valid("++x")
	c.valid("--x")
	c.valid("const a = {b: 0}; a.b++")
	c.valid("a['b']++")
	c.valid("/** @type {number} */ (x)++;")

	c.typeInfoValidationMode = AstValidator_TypeInfoValidation_NONE
	c.expectInvalid(rhino.NewNodeWith2Children(rhino.Token_INC, rhino.IR_Name("x"), rhino.IR_Name("x")), astValidatorTest_Check_EXPRESSION)
	c.expectInvalid(rhino.NewNodeWithChild(rhino.Token_INC, rhino.IR_ArrayPattern()), astValidatorTest_Check_EXPRESSION)
	c.expectInvalid(rhino.NewNodeWithChild(rhino.Token_INC, rhino.IR_ObjectPattern()), astValidatorTest_Check_EXPRESSION)
}

// Java: AstValidatorTest.testCompoundAssign()
func (c *AstValidatorTest) TestCompoundAssign() {

	c.valid("a += 1;")
	c.valid("a['b'] += 1;")
	c.valid("const a = {b: 0}; a.b += 1;")
	c.valid("const a = {b: '0'}; /** @type {?} */ (a.b) += 1;")

	c.typeInfoValidationMode = AstValidator_TypeInfoValidation_NONE
	c.expectInvalid(rhino.NewNodeWith2Children(rhino.Token_ASSIGN_ADD, rhino.IR_ArrayPattern(), rhino.IR_Number(0)), astValidatorTest_Check_EXPRESSION)
	c.expectInvalid(rhino.NewNodeWith2Children(rhino.Token_ASSIGN_ADD, rhino.IR_ObjectPattern(), rhino.IR_Number(0)), astValidatorTest_Check_EXPRESSION)
}

// Java: AstValidatorTest.testQuestionableForIn()
func (c *AstValidatorTest) TestQuestionableForIn() {

	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.SetExpectParseWarningsInThisTest()
	c.valid("for(var a = 1 in b);")
}

// Java: AstValidatorTest.testDebugger()
func (c *AstValidatorTest) TestDebugger() {

	c.valid("debugger;")
}

// Java: AstValidatorTest.testValidScript()
func (c *AstValidatorTest) TestValidScript() {

	// Since we're building the AST by hand, there won't be any types on it.
	c.typeInfoValidationMode = AstValidator_TypeInfoValidation_NONE

	n := rhino.NewNode(rhino.Token_SCRIPT)
	c.expectInvalid(n, astValidatorTest_Check_SCRIPT)
	n.SetInputId(rhino.NewInputId("something_input"))
	n.SetStaticSourceFile(rhino.NewSimpleSourceFile("something", rhino.StaticSourceFile_SourceKind_STRONG))
	c.expectValid(n, astValidatorTest_Check_SCRIPT)
	c.expectInvalid(n, astValidatorTest_Check_STATEMENT)
	c.expectInvalid(n, astValidatorTest_Check_EXPRESSION)
}

// Java: AstValidatorTest.testValidStatement1()
func (c *AstValidatorTest) TestValidStatement1() {

	// Since we're building the AST by hand, there won't be any types on it.
	c.typeInfoValidationMode = AstValidator_TypeInfoValidation_NONE

	n := rhino.NewNode(rhino.Token_RETURN)
	c.setTestSourceLocationForTree(n)
	c.expectInvalid(n, astValidatorTest_Check_EXPRESSION)
	c.expectValid(n, astValidatorTest_Check_STATEMENT)
	c.expectInvalid(n, astValidatorTest_Check_SCRIPT)
}

// Java: AstValidatorTest.testValidExpression1()
func (c *AstValidatorTest) TestValidExpression1() {

	// Since we're building the AST by hand, there won't be any types on it.
	c.typeInfoValidationMode = AstValidator_TypeInfoValidation_NONE

	n := rhino.NewNodeWithChild(rhino.Token_ARRAYLIT, rhino.NewNode(rhino.Token_EMPTY))
	c.setTestSourceLocationForTree(n)
	c.expectValid(n, astValidatorTest_Check_EXPRESSION)
	c.expectInvalid(n, astValidatorTest_Check_STATEMENT)
	c.expectInvalid(n, astValidatorTest_Check_SCRIPT)
}

// Java: AstValidatorTest.testValidExpression2()
func (c *AstValidatorTest) TestValidExpression2() {

	// Since we're building the AST by hand, there won't be any types on it.
	c.typeInfoValidationMode = AstValidator_TypeInfoValidation_NONE

	n := rhino.NewNodeWithChild(rhino.Token_NOT, rhino.NewNode(rhino.Token_TRUE))
	c.setTestSourceLocationForTree(n)
	c.expectValid(n, astValidatorTest_Check_EXPRESSION)
	c.expectInvalid(n, astValidatorTest_Check_STATEMENT)
	c.expectInvalid(n, astValidatorTest_Check_SCRIPT)
}

// Java: AstValidatorTest.testValidConst()
func (c *AstValidatorTest) TestValidConst() {

	c.valid("const x = r;")
	c.valid("const [x] = r;")
	c.valid("const {x} = {x: 1};")
	c.valid("const x = r, y = r;")
	c.valid("const {x} = {x: 1}, y = r;")
	c.valid("const x = r, {y} = {y: 1};")
}

// Java: AstValidatorTest.testInvalidConst()
func (c *AstValidatorTest) TestInvalidConst() {

	n := rhino.NewNode(rhino.Token_CONST)
	c.expectInvalid(n, astValidatorTest_Check_STATEMENT)

	n.AddChildToBack(rhino.IR_Name("x"))
	c.expectInvalid(n, astValidatorTest_Check_STATEMENT)

	n = rhino.NewNode(rhino.Token_CONST)
	n.AddChildToBack(rhino.NewNode(rhino.Token_DESTRUCTURING_LHS))
	n.GetFirstChild().AddChildToBack(rhino.NewNode(rhino.Token_OBJECT_PATTERN))

	c.expectInvalid(n, astValidatorTest_Check_STATEMENT)
}

// Java: AstValidatorTest.testInvalidConstLanguageLevel()
func (c *AstValidatorTest) TestInvalidConstLanguageLevel() {

	// Since we're building the AST by hand, there won't be any types on it.
	c.typeInfoValidationMode = AstValidator_TypeInfoValidation_NONE

	n := rhino.IR_ConstNode(rhino.IR_Name("x"), rhino.IR_Number(3))
	c.setTestSourceLocationForTree(n)

	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.expectInvalid(n, astValidatorTest_Check_STATEMENT)

	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_2015)
	c.expectValid(n, astValidatorTest_Check_STATEMENT)
}

// Java: AstValidatorTest.testInvalidLetLanguageLevel()
func (c *AstValidatorTest) TestInvalidLetLanguageLevel() {

	// Since we're building the AST by hand, there won't be any types on it.
	c.typeInfoValidationMode = AstValidator_TypeInfoValidation_NONE

	n := rhino.IR_LetWithValue(rhino.IR_Name("x"), rhino.IR_Number(3))
	c.setTestSourceLocationForTree(n)

	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.expectInvalid(n, astValidatorTest_Check_STATEMENT)

	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_2015)
	c.expectValid(n, astValidatorTest_Check_STATEMENT)
}

// Java: AstValidatorTest.testNewTargetIsValidExpression()
func (c *AstValidatorTest) TestNewTargetIsValidExpression() {

	// Since we're building the AST by hand, there won't be any types on it.
	c.typeInfoValidationMode = AstValidator_TypeInfoValidation_NONE

	n := rhino.NewNode(rhino.Token_NEW_TARGET)
	c.setTestSourceLocationForTree(n)
	c.expectValid(n, astValidatorTest_Check_EXPRESSION)
}

// Java: AstValidatorTest.testImportMetaIsValidExpression()
func (c *AstValidatorTest) TestImportMetaIsValidExpression() {

	// Since we're building the AST by hand, there won't be any types on it.
	c.typeInfoValidationMode = AstValidator_TypeInfoValidation_NONE

	n := rhino.NewNode(rhino.Token_IMPORT_META)
	c.setTestSourceLocationForTree(n)
	c.expectValid(n, astValidatorTest_Check_EXPRESSION)
}

// Java: AstValidatorTest.testCastOnLeftSideOfAssign()
func (c *AstValidatorTest) TestCastOnLeftSideOfAssign() {

	// Since we're building the AST by hand, there won't be any types on it.
	c.typeInfoValidationMode = AstValidator_TypeInfoValidation_NONE

	jsdoc := rhino.JSDocInfo_Builder_()
	jsdoc.RecordType(rhino.NewJSTypeExpression(rhino.IR_String("number"), "<AstValidatorTest>"))
	n :=
		rhino.IR_ExprResult(rhino.NewNodeWith2Children(rhino.Token_ASSIGN, rhino.IR_Cast(rhino.IR_Name("x"), jsdoc.Build()), rhino.IR_Number(0)))
	c.setTestSourceLocationForTree(n)
	c.expectValid(n, astValidatorTest_Check_STATEMENT)
}

// Java: AstValidatorTest.testInvalidEmptyStatement()
func (c *AstValidatorTest) TestInvalidEmptyStatement() {

	n := rhino.NewNodeWithChild(rhino.Token_EMPTY, rhino.NewNode(rhino.Token_TRUE))
	c.setTestSourceLocationForTree(n)
	c.expectInvalid(n, astValidatorTest_Check_STATEMENT)
	n.DetachChildren()
	c.expectValid(n, astValidatorTest_Check_STATEMENT)
}

// Java: AstValidatorTest.testInvalidNumberStatement()
func (c *AstValidatorTest) TestInvalidNumberStatement() {

	// Since we're building the AST by hand, there won't be any types on it.
	c.typeInfoValidationMode = AstValidator_TypeInfoValidation_NONE

	n := rhino.IR_Number(1)
	c.setTestSourceLocationForTree(n)
	c.expectInvalid(n, astValidatorTest_Check_STATEMENT)
	n = rhino.IR_ExprResult(n)
	c.setTestSourceLocationForTree(n)
	c.expectValid(n, astValidatorTest_Check_STATEMENT)
}

// Java: AstValidatorTest.testInvalidBigIntStatement()
func (c *AstValidatorTest) TestInvalidBigIntStatement() {

	// Since we're building the AST by hand, there won't be any types on it.
	c.typeInfoValidationMode = AstValidator_TypeInfoValidation_NONE

	n := rhino.IR_Bigint(big.NewInt(1))
	c.setTestSourceLocationForTree(n)
	c.expectInvalid(n, astValidatorTest_Check_STATEMENT)
	n = rhino.IR_ExprResult(n)
	c.setTestSourceLocationForTree(n)
	c.expectValid(n, astValidatorTest_Check_STATEMENT)
}

// Java: AstValidatorTest.testValidRestParameter()
func (c *AstValidatorTest) TestValidRestParameter() {

	c.SetLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_2015, CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.valid("function f(a,...rest){}")
	c.valid("function f(a,...[re,...st]){}")
}

// Java: AstValidatorTest.testDefaultParameter()
func (c *AstValidatorTest) TestDefaultParameter() {

	c.SetLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_2015, CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.valid("function f(a, b = 0){}")
}

// Java: AstValidatorTest.testAwaitExpression()
func (c *AstValidatorTest) TestAwaitExpression() {

	c.SetLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT, CompilerOptions_LanguageMode_ECMASCRIPT5)

	c.valid("async () => { await 1; }")
	c.valid("async function f() { await 1; }")
	c.valid("async function f() { return await 1; }")
	c.valid("async function f() { await await 1; }")
	c.valid("class C { async m() { await 1; } }")

	// Since we're building the AST by hand, there won't be any types on it.
	c.typeInfoValidationMode = AstValidator_TypeInfoValidation_NONE

	awaitNode := rhino.NewNode(rhino.Token_AWAIT)
	awaitNode.AddChildToBack(rhino.IR_Number(1))
	parentFunction :=
		rhino.IR_Function(rhino.IR_Name("foo"), rhino.IR_ParamList(), rhino.IR_BlockOfNode(rhino.IR_ReturnNodeWithExpr(awaitNode)))
	c.setTestSourceLocationForTree(awaitNode)
	parentFunction.SetIsAsyncFunction(true)
	c.expectValid(awaitNode, astValidatorTest_Check_EXPRESSION)
}

// Java: AstValidatorTest.testNoAwaitExpressionInDefaultParams()
func (c *AstValidatorTest) TestNoAwaitExpressionInDefaultParams() {

	// We're inserting our own Nodes below, and we won't be bothering to put valid type
	// information on them.
	c.typeInfoValidationMode = AstValidator_TypeInfoValidation_NONE
	scriptNode :=
		c.parseValidScript(
			"async function outer(){\n// `await` in a parameter default value is not allowed,\n// regardless of whether the function with the parameter is an\n// async function or enclosed in an async function.\n  async function inner(a = replaceWithAwait) {\n  }\n}\n")

	awaitNode := rhino.NewNode(rhino.Token_AWAIT)
	awaitNode.AddChildToBack(rhino.IR_Number(1))
	nodeToReplace :=
		astValidatorTest_findFirst(NodeUtil_PreOrderIterable(scriptNode), func(node *rhino.Node) bool { return node.IsName() && node.GetString() == ("replaceWithAwait") })
	nodeToReplace.ReplaceWith(awaitNode)
	c.expectInvalid(awaitNode, astValidatorTest_Check_EXPRESSION)
}

// Java: AstValidatorTest.testNoYieldExpressionInDefaultParams()
func (c *AstValidatorTest) TestNoYieldExpressionInDefaultParams() {

	// We're inserting our own Nodes below, and we won't be bothering to put valid type
	// information on them.
	c.typeInfoValidationMode = AstValidator_TypeInfoValidation_NONE
	scriptNode :=
		c.parseValidScript(
			"function *outer(){\n// `yield` in a parameter default value is not allowed,\n// regardless of whether the function with the parameter is a generator\n// or is enclosed in a generator function.\n  function *inner(a = replaceWithYield) {\n  }\n}\n")

	yieldNode := rhino.NewNode(rhino.Token_YIELD)
	yieldNode.AddChildToBack(rhino.IR_Number(1))
	nodeToReplace :=
		astValidatorTest_findFirst(NodeUtil_PreOrderIterable(scriptNode), func(node *rhino.Node) bool { return node.IsName() && node.GetString() == ("replaceWithYield") })
	nodeToReplace.ReplaceWith(yieldNode)
	c.expectInvalid(yieldNode, astValidatorTest_Check_EXPRESSION)
}

// Java: AstValidatorTest.testAwaitExpressionNonAsyncFunction()
func (c *AstValidatorTest) TestAwaitExpressionNonAsyncFunction() {

	c.SetLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT, CompilerOptions_LanguageMode_ECMASCRIPT5)
	awaitNode := rhino.NewNode(rhino.Token_AWAIT)
	awaitNode.AddChildToBack(rhino.IR_Number(1))
	parentFunction :=
		rhino.IR_Function(rhino.IR_Name("foo"), rhino.IR_ParamList(), rhino.IR_BlockOfNode(rhino.IR_ReturnNodeWithExpr(awaitNode)))
	parentFunction.SetIsAsyncFunction(false)
	c.expectInvalid(awaitNode, astValidatorTest_Check_EXPRESSION)
}

// Java: AstValidatorTest.testYieldExpressionNonGeneratorFunction()
func (c *AstValidatorTest) TestYieldExpressionNonGeneratorFunction() {

	yieldNode := rhino.NewNode(rhino.Token_YIELD)
	yieldNode.AddChildToBack(rhino.IR_Number(1))
	parentFunction :=
		rhino.IR_Function(rhino.IR_Name("foo"), rhino.IR_ParamList(), rhino.IR_BlockOfNode(rhino.IR_ReturnNodeWithExpr(yieldNode)))
	parentFunction.SetIsGeneratorFunction(false)
	c.expectInvalid(yieldNode, astValidatorTest_Check_EXPRESSION)
}

// Java: AstValidatorTest.testTopLevelAwait()
func (c *AstValidatorTest) TestTopLevelAwait() {

	// Since we're building the AST by hand, there won't be any types on it.
	c.typeInfoValidationMode = AstValidator_TypeInfoValidation_NONE

	c.SetLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT, CompilerOptions_LanguageMode_ECMASCRIPT5)
	awaitNode := rhino.NewNode(rhino.Token_AWAIT)
	awaitNode.AddChildToBack(rhino.IR_Number(1))
	script := rhino.IR_ScriptOfNodes(rhino.IR_ExprResult(awaitNode))
	c.setTestSourceLocationForTree(script)

	c.expectInvalid(
		awaitNode, astValidatorTest_Check_EXPRESSION, "'await' expression is not within an async function")

	// Top-level await is only allowed in modules.
	script.PutBooleanProp(rhino.Node_ES6_MODULE, true)
	c.expectValid(awaitNode, astValidatorTest_Check_EXPRESSION)
}

// Java: AstValidatorTest.testYieldExpressionNoFunction()
func (c *AstValidatorTest) TestYieldExpressionNoFunction() {

	// Since we're building the AST by hand, there won't be any types on it.
	c.typeInfoValidationMode = AstValidator_TypeInfoValidation_NONE

	c.SetLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT, CompilerOptions_LanguageMode_ECMASCRIPT5)
	n := rhino.NewNode(rhino.Token_YIELD)
	n.AddChildToBack(rhino.IR_Number(1))
	c.expectInvalid(n, astValidatorTest_Check_EXPRESSION)
}

// Java: AstValidatorTest.testInvalidArrayPattern0()
func (c *AstValidatorTest) TestInvalidArrayPattern0() {

	// Since we're building the AST by hand, there won't be any types on it.
	c.typeInfoValidationMode = AstValidator_TypeInfoValidation_NONE

	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_2015)

	// [...x = 1] = []
	n :=
		rhino.IR_Assign(
			rhino.NewNodeWithChild(
				rhino.Token_ARRAY_PATTERN,
				rhino.NewNodeWithChild(
					rhino.Token_ITER_REST, rhino.NewNodeWith2Children(rhino.Token_DEFAULT_VALUE, rhino.IR_Name("x"), rhino.IR_Arraylit()))),
			rhino.IR_Arraylit())
	c.expectInvalid(n, astValidatorTest_Check_EXPRESSION)
}

// Java: AstValidatorTest.testValidDestructuringAssignment0()
func (c *AstValidatorTest) TestValidDestructuringAssignment0() {

	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_2015)
	c.valid("var [x] = obj;")
	c.valid("var [x = 1] = obj;")
	c.valid("var [...y] = obj;")
	c.valid("var [x, ...y] = obj;")

	c.valid("[x] = [];")
	c.valid("[x = 1] = [];")
	c.valid("[x.y] = [];")
	c.valid("[x.y = 1] = [];")
	c.valid("[x['y']] = [];")
	c.valid("[x['y'] = 1] = [];")
	c.valid("[x().y] = [];")
	c.valid("[x().y = 1] = [];")
	c.valid("[x()['y']] = [];")
	c.valid("[x()['y'] = 1] = [];")

	c.valid("([...y] = obj);")
	c.valid("([x, ...y] = obj);")
	c.valid("([...this.x] = obj);")
	c.valid("([...this['x']] = obj);")
	c.valid("([...x.y] = obj);")
	c.valid("([...x['y']] = obj);")
	c.valid("([...x().y] = obj);")
	c.valid("([...x()['y']] = obj);")
}

// Java: AstValidatorTest.testValidDestructuringAssignment1()
func (c *AstValidatorTest) TestValidDestructuringAssignment1() {

	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_2015)
	c.valid("var {a:b} = {a: 1};")
	c.valid("({a:b} = {a: 1});")
	c.valid("({a:b.c} = {a: 1});")
	c.valid("({a:b().c} = {a: 1});")
	c.valid("({a:b['c']} = {a: 1});")
	c.valid("({a:b()['c']} = {a: 1});")
	c.valid("({a:b.c = 1} = {a: 1});")
	c.valid("({a:b().c = 1} = {a: 1});")
	c.valid("({a:b['c'] = 1} = {a: 1});")
	c.valid("({a:b()['c'] = 1} = {a: 1});")
}

// Java: AstValidatorTest.testValidDestructuringAssignment2()
func (c *AstValidatorTest) TestValidDestructuringAssignment2() {

	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_2015)
	c.valid("var {['a']:b} = obj;")
	c.valid("({['a']:b} = obj);")
	c.valid("({['a']:this.b} = obj);")
	c.valid("({['a']:b.c} = obj);")
	c.valid("({['a']:b.c = 1} = obj);")
	c.valid("({['a']:b().c} = obj);")
	c.valid("({['a']:b().c = 1} = obj);")
	c.valid("({['a']:b['c']} = obj);")
	c.valid("({['a']:b['c'] = 1} = obj);")
	c.valid("({['a']:b()['c']} = obj);")
	c.valid("({['a']:b()['c'] = 1} = obj);")
}

// Java: AstValidatorTest.testObjectRestAssignment()
func (c *AstValidatorTest) TestObjectRestAssignment() {

	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	c.valid("var {a, ...rest} = {a: 1, b: 2};")
	c.valid("({a:b, ...rest} = {a: 1, b: 2});")
	c.valid("({a:b.c, ...rest} = {a: 1, b: 2});")
	c.valid("({a:b().c, ...rest} = {a: 1, b: 2});")
	c.valid("({a:b['c'], ...rest} = {a: 1, b: 2});")
	c.valid("({a:b()['c'], ...rest} = {a: 1, b: 2});")
	c.valid("({a:b.c = 1, ...rest} = {a: 1, b: 2});")
	c.valid("({a:b().c = 1, ...rest} = {a: 1, b: 2});")
	c.valid("({a:b['c'] = 1, ...rest} = {a: 1, b: 2});")
	c.valid("({a:b()['c'] = 1, ...rest} = {a: 1, b: 2});")
}

// Java: AstValidatorTest.testInvalidObjectRestForLanguageLevel()
func (c *AstValidatorTest) TestInvalidObjectRestForLanguageLevel() {

	// Since we're building the AST by hand, there won't be any types on it.
	c.typeInfoValidationMode = AstValidator_TypeInfoValidation_NONE

	n := rhino.IR_Assign(rhino.IR_ObjectPattern(rhino.IR_ObjectRest(rhino.IR_Name("x"))), rhino.IR_Objectlit())
	c.setTestSourceLocationForTree(n)

	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_2015)
	c.expectInvalid(n, astValidatorTest_Check_EXPRESSION)

	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_2018)
	c.expectValid(n, astValidatorTest_Check_EXPRESSION)
}

// Java: AstValidatorTest.testInvalidArrayRestForLanguageLevel()
func (c *AstValidatorTest) TestInvalidArrayRestForLanguageLevel() {

	// Since we're building the AST by hand, there won't be any types on it.
	c.typeInfoValidationMode = AstValidator_TypeInfoValidation_NONE

	n := rhino.IR_Assign(rhino.IR_ArrayPattern(rhino.IR_IterRest(rhino.IR_Name("x"))), rhino.IR_Arraylit())
	c.setTestSourceLocationForTree(n)

	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT5)
	c.expectInvalid(n, astValidatorTest_Check_EXPRESSION)

	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_2015)
	c.expectValid(n, astValidatorTest_Check_EXPRESSION)
}

// Java: AstValidatorTest.testInvalidDestructuringDeclaration()
func (c *AstValidatorTest) TestInvalidDestructuringDeclaration() {

	// Since we're building the AST by hand, there won't be any types on it.
	c.typeInfoValidationMode = AstValidator_TypeInfoValidation_NONE

	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_2015)

	// missing DESTRUCTURING_LHS
	//     n := rhino.IR_Var(rhino.IR_ObjectPattern())
	//     c.expectInvalid(n, astValidatorTest_Check_STATEMENT)

	//     n = rhino.IR_Let(rhino.IR_ObjectPattern())
	//     c.expectInvalid(n, astValidatorTest_Check_STATEMENT)

	//     n = rhino.NewNodeWithChild(rhino.Token_CONST, rhino.IR_ObjectPattern())
	//     c.expectInvalid(n, astValidatorTest_Check_STATEMENT)

	// missing a right-hand side
	n := rhino.IR_Var(rhino.NewNodeWithChild(rhino.Token_DESTRUCTURING_LHS, rhino.IR_ObjectPattern()))
	c.expectInvalid(n, astValidatorTest_Check_STATEMENT)

	n = rhino.IR_Let(rhino.NewNodeWithChild(rhino.Token_DESTRUCTURING_LHS, rhino.IR_ObjectPattern()))
	c.expectInvalid(n, astValidatorTest_Check_STATEMENT)

	n = rhino.NewNodeWithChild(rhino.Token_CONST, rhino.NewNodeWithChild(rhino.Token_DESTRUCTURING_LHS, rhino.IR_ObjectPattern()))
	c.expectInvalid(n, astValidatorTest_Check_STATEMENT)
}

// Java: AstValidatorTest.testInvalidDestructuringAssignment()
func (c *AstValidatorTest) TestInvalidDestructuringAssignment() {

	// Since we're building the AST by hand, there won't be any types on it.
	c.typeInfoValidationMode = AstValidator_TypeInfoValidation_NONE

	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_2015)

	n :=
		rhino.IR_Assign(rhino.NewNodeWithChild(rhino.Token_OBJECT_PATTERN, rhino.NewNode(rhino.Token_ARRAY_PATTERN)), rhino.IR_Objectlit())
	c.expectInvalid(n, astValidatorTest_Check_EXPRESSION)

	n =
		rhino.IR_Assign(
			rhino.NewNodeWithChild(rhino.Token_ARRAY_PATTERN, rhino.IR_ComputedProp(rhino.IR_String("x"), rhino.IR_Number(1))),
			rhino.IR_Objectlit())
	c.expectInvalid(n, astValidatorTest_Check_EXPRESSION)

	stringkey := rhino.IR_StringKey("x")
	stringkey.AddChildToFront(rhino.IR_ComputedProp(rhino.IR_String("x"), rhino.IR_Number(1)))
	n = rhino.IR_Assign(rhino.NewNodeWithChild(rhino.Token_OBJECT_PATTERN, stringkey), rhino.IR_Objectlit())
	c.expectInvalid(n, astValidatorTest_Check_EXPRESSION)
}

// Java: AstValidatorTest.testGetter()
func (c *AstValidatorTest) TestGetter() {

	c.valid(
		`class C {
  get m1() {return 1}
}

/** @dict */
class D {
  get ['m2']() {return 2}
}
`)

	// Since we're modifying the AST by hand below, there won't be any types on it.
	c.typeInfoValidationMode = AstValidator_TypeInfoValidation_NONE

	cls := rhino.NewNodeWith2Children(rhino.Token_CLASS, rhino.IR_Name("C"), rhino.IR_Empty())
	members := rhino.NewNode(rhino.Token_CLASS_MEMBERS)
	cls.AddChildToBack(members)
	c.setTestSourceLocationForTree(cls)
	c.expectValid(cls, astValidatorTest_Check_STATEMENT)

	// Invalid getter with parameters
	getter1 := rhino.Node_NewStringWithToken(rhino.Token_GETTER_DEF, "prop")
	getter1.AddChildToBack(rhino.IR_Function(rhino.IR_Name(""), rhino.IR_ParamList(rhino.IR_Name("foo")), rhino.IR_Block()))
	getter1.SrcrefTree(members)
	members.AddChildToBack(getter1)
	c.expectInvalid(cls, astValidatorTest_Check_STATEMENT)

	members.DetachChildren()

	// Invalid getter with function name
	getter2 := rhino.Node_NewStringWithToken(rhino.Token_GETTER_DEF, "prop")
	getter2.AddChildToBack(rhino.IR_Function(rhino.IR_Name("foo"), rhino.IR_ParamList(), rhino.IR_Block()))
	getter2.SrcrefTree(members)
	members.AddChildToBack(getter2)
	c.expectInvalid(cls, astValidatorTest_Check_STATEMENT)

	members.DetachChildren()

	// Invalid computed property style getter with parameters
	getter3 :=
		rhino.NewNodeWith2Children(
			rhino.Token_COMPUTED_PROP,
			rhino.IR_String("prop"),
			rhino.IR_Function(rhino.IR_Name(""), rhino.IR_ParamList(rhino.IR_Name("foo")), rhino.IR_Block()))
	getter3.PutBooleanProp(rhino.Node_COMPUTED_PROP_GETTER, true)
	getter3.SrcrefTree(members)
	members.AddChildToBack(getter3)
	c.expectInvalid(cls, astValidatorTest_Check_STATEMENT)

	members.DetachChildren()

	// Invalid computed property style getter with function name
	getter4 :=
		rhino.NewNodeWith2Children(
			rhino.Token_COMPUTED_PROP,
			rhino.IR_String("prop"),
			rhino.IR_Function(rhino.IR_Name("foo"), rhino.IR_ParamList(), rhino.IR_Block()))
	getter4.PutBooleanProp(rhino.Node_COMPUTED_PROP_GETTER, true)
	getter4.SrcrefTree(members)
	members.AddChildToBack(getter4)
	c.expectInvalid(cls, astValidatorTest_Check_STATEMENT)
}

// Java: AstValidatorTest.testSuperInvalidAtScriptLevel()
func (c *AstValidatorTest) TestSuperInvalidAtScriptLevel() {

	scriptNode := rhino.IR_Script()
	superNode := rhino.IR_SuperNode()
	superStatementNode := rhino.IR_ExprResult(rhino.IR_Getprop(superNode, "prop"))
	scriptNode.AddChildToBack(superStatementNode)

	c.expectInvalid(superStatementNode, astValidatorTest_Check_STATEMENT)
}

// Java: AstValidatorTest.superInvalidWithOptChainCall()
func (c *AstValidatorTest) Test_superInvalidWithOptChainCall() {

	superNode := rhino.IR_SuperNode()
	optChainCallNode := rhino.IR_StartOptChainCall(superNode)

	c.expectInvalid(optChainCallNode, astValidatorTest_Check_STATEMENT)
}

// Java: AstValidatorTest.superInvalidWithOptionalGetProp()
func (c *AstValidatorTest) Test_superInvalidWithOptionalGetProp() {

	superNode := rhino.IR_SuperNode()
	optChainGetPropNode := rhino.IR_StartOptChainGetprop(superNode, "prop")

	c.expectInvalid(optChainGetPropNode, astValidatorTest_Check_STATEMENT)
}

// Java: AstValidatorTest.superInvalidWithOptionalGetElem()
func (c *AstValidatorTest) Test_superInvalidWithOptionalGetElem() {

	superNode := rhino.IR_SuperNode()
	exprNode := rhino.IR_Name("expr")
	optChainGetElemNode := rhino.IR_StartOptChainGetelem(superNode, exprNode)

	c.expectInvalid(optChainGetElemNode, astValidatorTest_Check_STATEMENT)
}

// Java: AstValidatorTest.optChainGetPropInvalidWithNoStartOfChain()
func (c *AstValidatorTest) Test_optChainGetPropInvalidWithNoStartOfChain() {

	innerGetProp := rhino.IR_ContinueOptChainGetprop(rhino.IR_Name("expr"), "prop1")
	outterGetProp := rhino.IR_ContinueOptChainGetprop(innerGetProp, "prop2")

	c.expectInvalid(outterGetProp, astValidatorTest_Check_STATEMENT)
}

// Java: AstValidatorTest.optChainGetElemInvalidWithNoStartOfChain()
func (c *AstValidatorTest) Test_optChainGetElemInvalidWithNoStartOfChain() {

	innerGetElem := rhino.IR_ContinueOptChainGetelem(rhino.IR_Name("expr"), rhino.IR_Name("prop1"))
	outterGetElem := rhino.IR_ContinueOptChainGetelem(innerGetElem, rhino.IR_Name("prop2"))

	c.expectInvalid(outterGetElem, astValidatorTest_Check_STATEMENT)
}

// Java: AstValidatorTest.optChainCallInvalidWithNoStartOfChain()
func (c *AstValidatorTest) Test_optChainCallInvalidWithNoStartOfChain() {

	innerCall := rhino.IR_ContinueOptChainCall(rhino.IR_Name("f"), rhino.IR_Name("arg1"))
	outterCall := rhino.IR_ContinueOptChainCall(innerCall, rhino.IR_Name("arg2"))

	c.expectInvalid(outterCall, astValidatorTest_Check_STATEMENT)
}

// Java: AstValidatorTest.testSuperInvalidInNonMemberFunction()
func (c *AstValidatorTest) TestSuperInvalidInNonMemberFunction() {

	c.invalid(
		`function nonMethod() {
  return super.toString();
}
`)
}

// Java: AstValidatorTest.testSuperPropReferenceIsValidInClassWithoutExtendsClause()
func (c *AstValidatorTest) TestSuperPropReferenceIsValidInClassWithoutExtendsClause() {

	c.valid(
		`class C {
  method() {
    super.prop;
    super['prop'];
  }
}
`)
}

// Java: AstValidatorTest.testSuperPropReferenceIsValidInObjectLiteralMethod()
func (c *AstValidatorTest) TestSuperPropReferenceIsValidInObjectLiteralMethod() {

	c.valid(
		`const myObj = {
  method() {
    super.prop;
    super['prop'];
  }
};
`)
}

// Java: AstValidatorTest.testSuperInvalidAsAnExpression()
func (c *AstValidatorTest) TestSuperInvalidAsAnExpression() {

	scriptNode :=
		c.parseValidScript(
			`class D {}
class C extends D {
  method() {
    return replaceWithSuper;
  }
}
`)

	nodeToReplace :=
		astValidatorTest_findFirst(NodeUtil_PreOrderIterable(scriptNode), func(node *rhino.Node) bool { return node.IsName() && node.GetString() == ("replaceWithSuper") })
	nodeToReplace.ReplaceWith(rhino.IR_SuperNode())

	c.expectInvalid(scriptNode, astValidatorTest_Check_SCRIPT)
}

// Java: AstValidatorTest.testSuperInvalidIfTypeInfoIsMissing()
func (c *AstValidatorTest) TestSuperInvalidIfTypeInfoIsMissing() {

	scriptNode :=
		c.parseValidScript(
			`class C {
  method() {
    return super.toString();
  }
}
`)

	superNode :=
		astValidatorTest_findFirst(NodeUtil_PreOrderIterable(scriptNode), (*rhino.Node).IsSuper)
	// Erase the type information to make the super node invalid
	superNode.SetJSType(nil)

	c.expectInvalid(scriptNode, astValidatorTest_Check_SCRIPT)
}

// Java: AstValidatorTest.testValidSuperConstructorCall()
func (c *AstValidatorTest) TestValidSuperConstructorCall() {

	c.valid(
		`class D {}
class C extends D {
  constructor() {
    super();
  }
}
`)
}

// Java: AstValidatorTest.testInvalidSuperConstructorCallInNonConstructor()
func (c *AstValidatorTest) TestInvalidSuperConstructorCallInNonConstructor() {

	c.invalid(
		`class D {}
class C extends D {
  nonConstructor() {
    super();
  }
}
`)
}

// Java: AstValidatorTest.testSuperConstructorCallInvalidInClassWithoutExtendsClause()
func (c *AstValidatorTest) TestSuperConstructorCallInvalidInClassWithoutExtendsClause() {

	c.invalid(
		`class C {
  constructor() {
    super();
  }
}
`)
}

// Java: AstValidatorTest.testSuperConstructorCallInvalidInObjectLiteralMethod()
func (c *AstValidatorTest) TestSuperConstructorCallInvalidInObjectLiteralMethod() {

	c.invalid(
		`const myObj = {
  constructor() {
    return super();
  }
}
`)
}

// Java: AstValidatorTest.testSetter()
func (c *AstValidatorTest) TestSetter() {

	c.valid(
		`class C {
  set m1(value) {}
}

/** @dict */
class D {
  set ['m2'](value) {}
}
`)

	// Since we're modifying the AST by hand below, there won't be types on some nodes that need
	// them.
	c.typeInfoValidationMode = AstValidator_TypeInfoValidation_NONE

	cls := rhino.NewNodeWith2Children(rhino.Token_CLASS, rhino.IR_Name("C"), rhino.IR_Empty())
	members := rhino.NewNode(rhino.Token_CLASS_MEMBERS)
	cls.AddChildToBack(members)
	c.setTestSourceLocationForTree(cls)
	c.expectValid(cls, astValidatorTest_Check_STATEMENT)

	// Invalid setter with no parameters
	setter1 := rhino.Node_NewStringWithToken(rhino.Token_SETTER_DEF, "prop")
	setter1.AddChildToBack(rhino.IR_Function(rhino.IR_Name(""), rhino.IR_ParamList(), rhino.IR_Block()))
	setter1.SrcrefTree(members)
	members.AddChildToBack(setter1)
	c.expectInvalid(cls, astValidatorTest_Check_STATEMENT)

	members.DetachChildren()

	// Invalid setter with function name
	setter2 := rhino.Node_NewStringWithToken(rhino.Token_SETTER_DEF, "prop")
	setter2.AddChildToBack(rhino.IR_Function(rhino.IR_Name("foo"), rhino.IR_ParamList(rhino.IR_Name("value")), rhino.IR_Block()))
	setter2.SrcrefTree(members)
	members.AddChildToBack(setter2)
	c.expectInvalid(cls, astValidatorTest_Check_STATEMENT)

	members.DetachChildren()

	// Invalid computed property style setter with parameters
	setter3 :=
		rhino.NewNodeWith2Children(
			rhino.Token_COMPUTED_PROP,
			rhino.IR_String("prop"),
			rhino.IR_Function(rhino.IR_Name(""), rhino.IR_ParamList(), rhino.IR_Block()))
	setter3.PutBooleanProp(rhino.Node_COMPUTED_PROP_SETTER, true)
	setter3.SrcrefTree(members)
	members.AddChildToBack(setter3)
	c.expectInvalid(cls, astValidatorTest_Check_STATEMENT)

	members.DetachChildren()

	// Invalid computed property style setter with function name
	setter4 :=
		rhino.NewNodeWith2Children(
			rhino.Token_COMPUTED_PROP,
			rhino.IR_String("prop"),
			rhino.IR_Function(rhino.IR_Name("foo"), rhino.IR_ParamList(rhino.IR_Name("value")), rhino.IR_Block()))
	setter4.PutBooleanProp(rhino.Node_COMPUTED_PROP_SETTER, true)
	setter4.SrcrefTree(members)
	members.AddChildToBack(setter4)
	c.expectInvalid(cls, astValidatorTest_Check_STATEMENT)
}

// Java: AstValidatorTest.testFeatureValidation_getter()
func (c *AstValidatorTest) TestFeatureValidation_getter() {

	c.testFeatureValidation("var obj = {get f() {}};", parser.FeatureSet_Feature_GETTER)
}

// Java: AstValidatorTest.testFeatureValidation_setter()
func (c *AstValidatorTest) TestFeatureValidation_setter() {

	c.testFeatureValidation("var obj = {set f(x) {}};", parser.FeatureSet_Feature_SETTER)
}

// Java: AstValidatorTest.testFeatureValidation_classGetterSetter()
func (c *AstValidatorTest) TestFeatureValidation_classGetterSetter() {

	c.testFeatureValidation("class C { get f() {} }", parser.FeatureSet_Feature_CLASS_GETTER_SETTER)
	c.testFeatureValidation("class C { set f(x) {} }", parser.FeatureSet_Feature_CLASS_GETTER_SETTER)
}

// Java: AstValidatorTest.testFeatureValidation_arrowFunctions()
func (c *AstValidatorTest) TestFeatureValidation_arrowFunctions() {

	c.testFeatureValidation("var arrow = () => 3", parser.FeatureSet_Feature_ARROW_FUNCTIONS)
	c.testFeatureValidation("var asyncArrow = async () => 3", parser.FeatureSet_Feature_ARROW_FUNCTIONS)
}

// Java: AstValidatorTest.testFeatureValidation_blockScopedFunctionDeclaration()
func (c *AstValidatorTest) TestFeatureValidation_blockScopedFunctionDeclaration() {

	c.testFeatureValidation("{ function f() {} }", parser.FeatureSet_Feature_BLOCK_SCOPED_FUNCTION_DECLARATION)
	c.testFeatureValidation(
		"function f() { if (true) { function g() {} } }",
		parser.FeatureSet_Feature_BLOCK_SCOPED_FUNCTION_DECLARATION)
	c.valid("function f() {}")
}

// Java: AstValidatorTest.testFeatureValidation_classes()
func (c *AstValidatorTest) TestFeatureValidation_classes() {

	c.testFeatureValidation("class C {}", parser.FeatureSet_Feature_CLASSES)
	c.testFeatureValidation("var C = class {}", parser.FeatureSet_Feature_CLASSES)
}

// Java: AstValidatorTest.testFeatureValidation_computedProperties()
func (c *AstValidatorTest) TestFeatureValidation_computedProperties() {

	c.testFeatureValidation("var obj = { ['foo' + 3]: 4};", parser.FeatureSet_Feature_COMPUTED_PROPERTIES)
	c.testFeatureValidation("var { ['foo' + 3]: x} = obj;", parser.FeatureSet_Feature_COMPUTED_PROPERTIES)
	c.testFeatureValidation("/** @dict */ class C { ['foobar']() {} }", parser.FeatureSet_Feature_COMPUTED_PROPERTIES)
}

// Java: AstValidatorTest.testFeatureValidation_defaultParameters()
func (c *AstValidatorTest) TestFeatureValidation_defaultParameters() {

	c.testFeatureValidation("function f(a = 1) {}", parser.FeatureSet_Feature_DEFAULT_PARAMETERS)
	c.testFeatureValidation("((a = 3) => a)", parser.FeatureSet_Feature_DEFAULT_PARAMETERS)
}

// Java: AstValidatorTest.testFeatureValidation_arrayDestructuring()
func (c *AstValidatorTest) TestFeatureValidation_arrayDestructuring() {

	c.testFeatureValidation("var x, [a, b] = arr;", parser.FeatureSet_Feature_ARRAY_DESTRUCTURING)
	c.testFeatureValidation("x = 0, [a, b] = obj;", parser.FeatureSet_Feature_ARRAY_DESTRUCTURING)
	c.testFeatureValidation("for ([a, b] of c) {}", parser.FeatureSet_Feature_ARRAY_DESTRUCTURING)
	c.testFeatureValidation("function f([a, b]) {}", parser.FeatureSet_Feature_ARRAY_DESTRUCTURING)
}

// Java: AstValidatorTest.testFeatureValidation_objectDestructuring()
func (c *AstValidatorTest) TestFeatureValidation_objectDestructuring() {

	c.testFeatureValidation("var x, {a, b} = {a: 1, b: 2};", parser.FeatureSet_Feature_OBJECT_DESTRUCTURING)
	c.testFeatureValidation("(x = 0, {a, b} = {a: 1, b: 2});", parser.FeatureSet_Feature_OBJECT_DESTRUCTURING)
	c.testFeatureValidation(
		`/** @type {!Array<{a: string, b:string}>} */
const c = [];
for ({a, b} of c) {}
`,
		parser.FeatureSet_Feature_OBJECT_DESTRUCTURING)
	c.testFeatureValidation(
		`/** @param {{a: string, b:string}} p1 */
function f({a, b}) {}
`,
		parser.FeatureSet_Feature_OBJECT_DESTRUCTURING)
}

// Java: AstValidatorTest.testFeatureValidation_extendedObjectLiterals()
func (c *AstValidatorTest) TestFeatureValidation_extendedObjectLiterals() {

	c.testFeatureValidation("var obj = { x };", parser.FeatureSet_Feature_SHORTHAND_OBJECT_PROPERTIES)
}

// Java: AstValidatorTest.testFeatureValidation_forOf()
func (c *AstValidatorTest) TestFeatureValidation_forOf() {

	c.testFeatureValidation("for (const a of b) {}", parser.FeatureSet_Feature_FOR_OF)
}

// Java: AstValidatorTest.testFeatureValidation_forAwaitOf()
func (c *AstValidatorTest) TestFeatureValidation_forAwaitOf() {

	c.testFeatureValidation("async () => { for await (const a of b) {} }", parser.FeatureSet_Feature_FOR_AWAIT_OF)
}

// Java: AstValidatorTest.testFeatureValidation_generatorFunctions()
func (c *AstValidatorTest) TestFeatureValidation_generatorFunctions() {

	c.testFeatureValidation("const f = function *() {}", parser.FeatureSet_Feature_GENERATORS)
	c.testFeatureValidation("function *f() {}", parser.FeatureSet_Feature_GENERATORS)
	c.testFeatureValidation("class C { *f() {} }", parser.FeatureSet_Feature_GENERATORS)
}

// Java: AstValidatorTest.testFeatureValidation_memberDeclarations()
func (c *AstValidatorTest) TestFeatureValidation_memberDeclarations() {

	c.testFeatureValidation("class C { f() {} }", parser.FeatureSet_Feature_MEMBER_DECLARATIONS)
	c.testFeatureValidation("var obj = { f() {} };", parser.FeatureSet_Feature_MEMBER_DECLARATIONS)
}

// Java: AstValidatorTest.testFeatureValidation_newTarget()
func (c *AstValidatorTest) TestFeatureValidation_newTarget() {

	c.testFeatureValidation("function f() { new.target }", parser.FeatureSet_Feature_NEW_TARGET)
}

// Java: AstValidatorTest.testFeatureValidation_restParameters()
func (c *AstValidatorTest) TestFeatureValidation_restParameters() {

	c.testFeatureValidation("function f(...rest) {}", parser.FeatureSet_Feature_REST_PARAMETERS)
}

// Java: AstValidatorTest.testFeatureValidation_spreadExpressions()
func (c *AstValidatorTest) TestFeatureValidation_spreadExpressions() {

	c.testFeatureValidation("f(...arr);", parser.FeatureSet_Feature_SPREAD_EXPRESSIONS)
	c.testFeatureValidation("var arr = [...something];", parser.FeatureSet_Feature_SPREAD_EXPRESSIONS)
	c.testFeatureValidation("var obj = {...something};", parser.FeatureSet_Feature_OBJECT_LITERALS_WITH_SPREAD)
}

// Java: AstValidatorTest.testFeatureValidation_super()
func (c *AstValidatorTest) TestFeatureValidation_super() {

	c.testFeatureValidation("const obj = { f() { super.toString(); } };", parser.FeatureSet_Feature_SUPER)
}

// Java: AstValidatorTest.testFeatureValidation_templateLiterals()
func (c *AstValidatorTest) TestFeatureValidation_templateLiterals() {

	c.testFeatureValidation("`foo ${3} bar `", parser.FeatureSet_Feature_TEMPLATE_LITERALS)
	c.testFeatureValidation("tag`foo ${3} bar`", parser.FeatureSet_Feature_TEMPLATE_LITERALS)
}

// Java: AstValidatorTest.testFeatureValidation_taggedTemplateLiteralWithInvalidEscapes()
func (c *AstValidatorTest) TestFeatureValidation_taggedTemplateLiteralWithInvalidEscapes() {

	// Tagged template literals are allowed to contain invalid escape sequences,
	// so these should not cause compiler errors or warnings.
	c.testFeatureValidation("foo`\\unicode`", parser.FeatureSet_Feature_TEMPLATE_LITERALS)
	c.testFeatureValidation("foo`\\xray ${3} \\u{42`", parser.FeatureSet_Feature_TEMPLATE_LITERALS)
}

// Java: AstValidatorTest.testFeatureValidation_modules()
func (c *AstValidatorTest) TestFeatureValidation_modules() {

	// Modules need to be set up better than we're doing here to avoid type check throwing an
	// exception
	c.DisableTypeCheck()
	c.typeInfoValidationMode = AstValidator_TypeInfoValidation_NONE
	c.testFeatureValidation("export {x};", parser.FeatureSet_Feature_MODULES)
	c.testFeatureValidation("import {x} from './foo.js';", parser.FeatureSet_Feature_MODULES)
}

// Java: AstValidatorTest.testFeatureValidation_exponentOp()
func (c *AstValidatorTest) TestFeatureValidation_exponentOp() {

	c.testFeatureValidation("2 ** 3", parser.FeatureSet_Feature_EXPONENT_OP)
	c.testFeatureValidation("x **= 3;", parser.FeatureSet_Feature_EXPONENT_OP)
}

// Java: AstValidatorTest.testFeatureValidation_nullishCoalesceOp()
func (c *AstValidatorTest) TestFeatureValidation_nullishCoalesceOp() {

	c.testFeatureValidation("x ?? y", parser.FeatureSet_Feature_NULL_COALESCE_OP)
	c.testFeatureValidation("x ?? y ?? z", parser.FeatureSet_Feature_NULL_COALESCE_OP)
}

// Java: AstValidatorTest.testFeatureValidation_logicalAssignmentOp()
func (c *AstValidatorTest) TestFeatureValidation_logicalAssignmentOp() {

	// TODO (user): re-enable TypeInfoValidation and TypeCheck
	c.typeInfoValidationMode = AstValidator_TypeInfoValidation_NONE
	c.DisableTypeCheck()
	c.testFeatureValidation("x ||= y", parser.FeatureSet_Feature_LOGICAL_ASSIGNMENT)
	c.testFeatureValidation("x &&= y", parser.FeatureSet_Feature_LOGICAL_ASSIGNMENT)
	c.testFeatureValidation("x ??= y", parser.FeatureSet_Feature_LOGICAL_ASSIGNMENT)
}

// Java: AstValidatorTest.testFeatureValidation_optChain()
func (c *AstValidatorTest) TestFeatureValidation_optChain() {

	c.testFeatureValidation("x?.y", parser.FeatureSet_Feature_OPTIONAL_CHAINING)
	c.testFeatureValidation("x?.()", parser.FeatureSet_Feature_OPTIONAL_CHAINING)
	c.testFeatureValidation("x?.[1]", parser.FeatureSet_Feature_OPTIONAL_CHAINING)
}

// Java: AstValidatorTest.testFeatureValidation_asyncFunctions()
func (c *AstValidatorTest) TestFeatureValidation_asyncFunctions() {

	c.testFeatureValidation("const f = async function() {}", parser.FeatureSet_Feature_ASYNC_FUNCTIONS)
	c.testFeatureValidation("async function f() {}", parser.FeatureSet_Feature_ASYNC_FUNCTIONS)
	c.testFeatureValidation("class C { async f() {} }", parser.FeatureSet_Feature_ASYNC_FUNCTIONS)
	c.testFeatureValidation("(async () => {})", parser.FeatureSet_Feature_ASYNC_FUNCTIONS)
}

// Java: AstValidatorTest.testFeatureValidation_asyncGeneratorFunctions()
func (c *AstValidatorTest) TestFeatureValidation_asyncGeneratorFunctions() {

	c.testFeatureValidation("const f = async function *() {}", parser.FeatureSet_Feature_ASYNC_GENERATORS)
	c.testFeatureValidation("async function *f() {}", parser.FeatureSet_Feature_ASYNC_GENERATORS)
	c.testFeatureValidation("class C { async *f() {} }", parser.FeatureSet_Feature_ASYNC_GENERATORS)
}

// Java: AstValidatorTest.testFeatureValidation_objectLiteralsWithSpread()
func (c *AstValidatorTest) TestFeatureValidation_objectLiteralsWithSpread() {

	c.testFeatureValidation("var obj = {...something};", parser.FeatureSet_Feature_OBJECT_LITERALS_WITH_SPREAD)
}

// Java: AstValidatorTest.testValidFeatureInScript()
func (c *AstValidatorTest) TestValidFeatureInScript() {

	// Since we're building the AST by hand, there won't be any types on it.
	c.typeInfoValidationMode = AstValidator_TypeInfoValidation_NONE

	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_2015)

	n := rhino.NewNode(rhino.Token_SCRIPT)
	n.SetInputId(rhino.NewInputId("something_input"))
	n.SetStaticSourceFile(rhino.NewSimpleSourceFile("something", rhino.StaticSourceFile_SourceKind_STRONG))
	c.expectValid(n, astValidatorTest_Check_SCRIPT)

	n.AddChildToFront(rhino.IR_LetWithValue(rhino.IR_Name("a"), rhino.IR_Number(3)))
	n.SrcrefTree(n)
	c.expectInvalid(n, astValidatorTest_Check_SCRIPT)

	n.PutProp(rhino.Node_FEATURE_SET, parser.FeatureSet_BARE_MINIMUM.With(parser.FeatureSet_Feature_LET_DECLARATIONS))
	c.expectValid(n, astValidatorTest_Check_SCRIPT)

	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT3) // resets compiler's allowable featureSet
	c.expectInvalid(n, astValidatorTest_Check_SCRIPT)
	// violation reported from {@code validateFeature} call of LET because
	// `!allowbleFeatures.Has(feature)`
	AssertThat(c.lastCheckViolationMessages).Contains("AST should not contain let declaration")
	// violation reported from {@code validateScript} call because script's `FEATURE_SET` is not
	// a subset of compiler's allowable featureSet.
	AssertThat(c.lastCheckViolationMessages).
		Contains("SCRIPT node contains these unallowable features:[let declaration]")
}

// Java: AstValidatorTest.testUnresolvedTypesAreBanned()
func (c *AstValidatorTest) TestUnresolvedTypesAreBanned() {
	compiler := c.CreateCompiler()
	compiler.InitOptions(c.GetOptions())
	registry := compiler.GetTypeRegistry()
	closer := registry.GetResolver().OpenForDefinition()

	var unresolvedFunction rhino.JSType = registry.CreateFunctionType(registry.GetNativeType(rhino.JSTypeNative_VOID_TYPE))
	AssertThat(unresolvedFunction.IsResolved()).IsFalse()

	foo := rhino.IR_Name("foo").SetJSType(unresolvedFunction)
	expr := rhino.IR_ExprResult(foo)
	c.setTestSourceLocationForTree(expr)

	c.expectInvalid(expr, astValidatorTest_Check_STATEMENT)

	closer.Close()

	c.expectValid(expr, astValidatorTest_Check_STATEMENT)
}

// Java: AstValidatorTest.testValidatesColorInfoOnExpression()
func (c *AstValidatorTest) TestValidatesColorInfoOnExpression() {

	c.typeInfoValidationMode = AstValidator_TypeInfoValidation_COLOR

	foo := rhino.IR_String("foo")
	expr := rhino.IR_ExprResult(foo)
	script := rhino.IR_ScriptOfNodes(expr)
	c.setTestSourceLocationForTree(script)

	c.expectInvalid(expr, astValidatorTest_Check_STATEMENT)

	compiler := c.CreateCompiler()

	registry := compiler.GetTypeRegistry()
	foo.SetJSType(registry.GetNativeType(rhino.JSTypeNative_STRING_TYPE))
	c.expectInvalid(expr, astValidatorTest_Check_STATEMENT)

	foo.SetJSType(nil)
	foo.SetColor(colors.StandardColors_STRING)

	c.expectValid(expr, astValidatorTest_Check_STATEMENT)
}

// Java: AstValidatorTest.testSwitchStatement()
func (c *AstValidatorTest) TestSwitchStatement() {

	// Since we're building the AST by hand, there won't be any types on it.
	c.typeInfoValidationMode = AstValidator_TypeInfoValidation_NONE

	switchStatement :=
		`function foo(x) {
  switch(x) {
    case 1:
      var y = 5;
    case 2:
      var y = 5;
  }
}
`

	c.valid(switchStatement)

	rootScriptNode := c.parseValidScript(switchStatement)
	blockNode := rhino.IR_Block()

	firstCaseNode :=
		astValidatorTest_findFirst(NodeUtil_PreOrderIterable(rootScriptNode), (*rhino.Node).IsCase)

	firstCaseNode.AddChildToFront(blockNode)
	// A CASE node is only allowed to have 2 children, 1 expression and 1 block.
	c.expectInvalid(rootScriptNode, astValidatorTest_Check_SCRIPT)
}

// Java: AstValidatorTest.testShadowContent_validatesSingularExpectedStructure()
func (c *AstValidatorTest) TestShadowContent_validatesSingularExpectedStructure() {

	// Since we're building the AST by hand, there won't be any types on it.
	c.typeInfoValidationMode = AstValidator_TypeInfoValidation_NONE

	shadowHost := rhino.IR_Name("f")
	f := rhino.IR_ExprResult(shadowHost)
	c.expectValid(f, astValidatorTest_Check_STATEMENT)

	shadowHost.SetClosureUnawareShadow(rhino.IR_Name("x"))
	c.expectInvalid(f, astValidatorTest_Check_STATEMENT, "Shadow reference node is not a ROOT node")

	shadowHost.SetClosureUnawareShadow(rhino.IR_Root())
	c.expectInvalid(f, astValidatorTest_Check_STATEMENT, "Shadow root node's child is not a script node")

	shadowHost.SetClosureUnawareShadow(
		rhino.IR_Root(rhino.IR_ScriptOfNodes(rhino.IR_ExprResult(rhino.IR_Name("x")), rhino.IR_ExprResult(rhino.IR_Name("y")))))
	c.expectInvalid(f, astValidatorTest_Check_STATEMENT, "Shadow SCRIPT node child has more than one child")

	shadowHost.SetClosureUnawareShadow(rhino.IR_Root(c.parseValidScript("function foo(){}").Detach()))
	c.expectInvalid(f, astValidatorTest_Check_STATEMENT, "Shadow SCRIPT node child is not an expr result node")

	shadowHost.SetClosureUnawareShadow(rhino.IR_Root(c.parseValidScript("(x++)").Detach()))
	c.expectInvalid(f, astValidatorTest_Check_STATEMENT, "Shadow node EXPR_RESULT child is not a call")

	shadowHost.SetClosureUnawareShadow(rhino.IR_Root(c.parseValidScript("x()").Detach()))
	c.expectInvalid(f, astValidatorTest_Check_STATEMENT, "Shadow node CALL child is not a function")

	shadowHost.SetClosureUnawareShadow(
		rhino.IR_Root(c.parseValidScript("sink(function(){})").Detach()))
	c.expectValid(f, astValidatorTest_Check_STATEMENT)
}

// Java: AstValidatorTest.checkRequiredInlinings_failsWithRemainingRequiredInlinings()
func (c *AstValidatorTest) Test_checkRequiredInlinings_failsWithRemainingRequiredInlinings() {

	// Write a reference to a function that required inlining. We need an enclosing statement
	// because the expression checker will look for a parent.
	myFunction := rhino.IR_Name("myFunction")
	rhino.IR_Root(rhino.IR_ScriptOfNodes(rhino.IR_ExprResult(myFunction)))

	// Add the requireInlining annotation.
	jsDocInfo := rhino.JSDocInfo_Builder_()
	jsDocInfo.RecordRequireInlining()
	myFunction.SetJSDocInfo(jsDocInfo.Build())

	// A surviving requireInlining annotation should cause a validation failure.
	c.expectInvalid(
		myFunction,
		astValidatorTest_Check_EXPRESSION,
		`Type information missing
/** @requireInlining */ myFunction`,
		"@requireInlining node failed to be inlined.",
		`Type information missing
/** @requireInlining */ myFunction`)
}

// Java: AstValidatorTest.valid(String)
func (c *AstValidatorTest) valid(code string) {
	c.TestSame(code)
	AssertThat(c.lastCheckViolationMessages).IsEmpty()
}

// parseValidScript parses and validates code for a script that should be
// valid. Returns the SCRIPT Node created by the parsing.
//
// Java: AstValidatorTest.parseValidScript(String)
func (c *AstValidatorTest) parseValidScript(scriptCode string) *rhino.Node {
	c.TestSame(scriptCode)
	AssertThat(c.lastCheckViolationMessages).IsEmpty()
	return c.GetLastCompiler().GetJsRoot().GetFirstChild()
}

// Java: AstValidatorTest.invalid(String)
func (c *AstValidatorTest) invalid(code string) {
	c.TestSame(code)
	AssertThat(c.lastCheckViolationMessages).IsNotEmpty()
}

// astValidatorTest_Check is AstValidatorTest.Check.
type astValidatorTest_Check int32

const (
	astValidatorTest_Check_SCRIPT astValidatorTest_Check = iota
	astValidatorTest_Check_STATEMENT
	astValidatorTest_Check_EXPRESSION
)

// doCheck performs validation check: n is the tree of nodes ot check;
// level the level at which the check starts. Returns the list of
// violations.
//
// Java: AstValidatorTest.doCheck(Node, Check)
func (c *AstValidatorTest) doCheck(n *rhino.Node, level astValidatorTest_Check) *jutil.ArrayList[string] {
	compiler := c.CreateCompiler()
	compiler.InitOptions(c.GetOptions())
	validator := c.createValidator(compiler)
	switch level {
	case astValidatorTest_Check_SCRIPT:
		validator.ValidateScript(n)
	case astValidatorTest_Check_STATEMENT:
		validator.ValidateStatement(n)
	case astValidatorTest_Check_EXPRESSION:
		validator.ValidateExpression(n)
	}
	return c.lastCheckViolationMessages
}

// expectInvalid is the two Java overloads expectInvalid(Node, Check) and
// expectInvalid(Node, Check, String...).
//
// Go-only: test-helper — javac binds a call without messages to the
// non-varargs overload (some violation), and one with messages to the
// varargs overload (exactly these violations).
//
// Java: AstValidatorTest.expectInvalid(Node, Check)
// Java: AstValidatorTest.expectInvalid(Node, Check, String...)
func (c *AstValidatorTest) expectInvalid(n *rhino.Node, level astValidatorTest_Check, validationMessages ...string) {
	if len(validationMessages) == 0 {
		AssertThat(c.doCheck(n, level)).IsNotEmpty()
		return
	}
	AssertThat(c.doCheck(n, level)).ContainsExactlyElementsIn(validationMessages)
}

// Java: AstValidatorTest.expectValid(Node, Check)
func (c *AstValidatorTest) expectValid(n *rhino.Node, level astValidatorTest_Check) {
	AssertThat(c.doCheck(n, level)).IsEmpty()
}

// setTestSourceLocationForTree adds source location information to a tree
// of nodes.
//
// All non-ROOT nodes are expected to have at minimum a property indicating
// their original source file. Make sure this is true for all of the nodes in
// this tree to avoid property validation violations.
//
// Java: AstValidatorTest.setTestSourceLocationForTree(Node)
func (c *AstValidatorTest) setTestSourceLocationForTree(n *rhino.Node) {
	jutil.CheckStateMsg(!n.IsRoot(), "ROOT nodes don't get source files")
	n.SetSourceFileForTesting("testcode")
	// make sure all child nodes get the same source file.
	n.SrcrefTree(n)
}

// testFeatureValidation tests that AstValidator checks for the given
// feature in the AST
//
// This will raise an error if a) the AST parsed from code lacks feature, or
// b) AstValidator does not validate feature's presence in the AST.
//
// Java: AstValidatorTest.testFeatureValidation(String, Feature)
func (c *AstValidatorTest) testFeatureValidation(code string, feature parser.FeatureSet_Feature) {
	c.valid(code)
	script := c.GetLastCompiler().GetJsRoot().GetFirstChild()

	// Remove `feature` from the SCRIPT node's feature set, checking that it
	// was originally present, and then validate that AstValidator errors
	// because it expects `feature` to be present.
	currentFeatures := NodeUtil_GetFeatureSetOfScript(script)
	AssertThat(currentFeatures.ContainsFeature(feature)).IsTrue()

	script.PutProp(rhino.Node_FEATURE_SET, currentFeatures.WithoutFeatures(feature))
	c.expectInvalid(script, astValidatorTest_Check_SCRIPT)
}

// astValidatorTest_findFirst is
// stream(nodes).filter(predicate).findFirst().get().
//
// Go-only: test-helper — Optional.get() throws NoSuchElementException when
// nothing matches.
func astValidatorTest_findFirst(nodes iter.Seq[*rhino.Node], predicate func(*rhino.Node) bool) *rhino.Node {
	for n := range nodes {
		if predicate(n) {
			return n
		}
	}
	panic(jutil.NewNoSuchElementExceptionMsg("No value present"))
}
