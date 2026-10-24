// Port of com/google/javascript/jscomp/TypedScopeCreatorTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/TypedScopeCreatorTest.java)

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
	"slices"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// TypedScopeCreatorTest holds the tests for TypedScopeCreator and
// TypeInference. Admittedly, the name is a bit of a misnomer.
type TypedScopeCreatorTest struct {
	CompilerTestCaseBase

	registry                 *rhino.JSTypeRegistry
	globalScope              *TypedScope
	lastLocalScope           *TypedScope
	lastFunctionScope        *TypedScope
	processClosurePrimitives bool

	// Maps a label name to information about the labeled statement.
	//
	// This map is recreated each time parseAndRunTypeInference() is
	// executed. TODO(bradfordcsmith): This map and LabeledStatement are also
	// in TypeInferenceTest. It would be good to unify them.
	labeledStatementMap *jutil.LinkedHashMap[string, *TypedScopeCreatorTest_LabeledStatement]
}

// TypedScopeCreatorTest_LabeledStatement stores information about a labeled
// statement and allows making assertions on it.
type TypedScopeCreatorTest_LabeledStatement struct {
	statementNode  rhino.Node
	enclosingScope *TypedScope
}

// Java: TypedScopeCreatorTest.LabeledStatement(Node, TypedScope)
func NewTypedScopeCreatorTest_LabeledStatement(statementNode rhino.Node, enclosingScope *TypedScope) *TypedScopeCreatorTest_LabeledStatement {
	s := &TypedScopeCreatorTest_LabeledStatement{}
	s.statementNode = rhino.CheckNotNull(statementNode)
	s.enclosingScope = jutil.CheckNotNull(enclosingScope)
	return s
}

// Java: TypedScopeCreatorTest()
func NewTypedScopeCreatorTest() *TypedScopeCreatorTest {
	c := &TypedScopeCreatorTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, CompilerTestCase_DEFAULT_EXTERNS)
	c.processClosurePrimitives = false
	return c
}

func TestTypedScopeCreatorTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewTypedScopeCreatorTest() })
}

// Java: TypedScopeCreatorTest.setUp()
func (c *TypedScopeCreatorTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableTypeInfoValidation()
}

// TypedScopeCreatorTest_ScopeFinder records the global scope, the last
// block and function scopes and every labeled statement.
type TypedScopeCreatorTest_ScopeFinder struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	// Go-only: language-adapter — the enclosing TypedScopeCreatorTest
	// instance of this inner class.
	outer *TypedScopeCreatorTest
}

// Java: TypedScopeCreatorTest.ScopeFinder.visit(NodeTraversal, Node, Node)
func (f *TypedScopeCreatorTest_ScopeFinder) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	scope := t.GetTypedScope()
	if scope.IsGlobal() {
		f.outer.globalScope = scope
	} else if scope.IsBlockScope() {
		// TODO(bradfordcsmith): use labels to find scopes instead of lastLocalScope
		f.outer.lastLocalScope = scope
	} else if scope.IsFunctionScope() {
		f.outer.lastFunctionScope = scope
	}
	if !parent.IsNil() && parent.IsLabel() && !n.IsLabelName() {
		// First child of a LABEL is a LABEL_NAME, n is the second child.
		labelNameNode := rhino.CheckNotNullMsg1(n.GetPrevious(), "%s", n)
		jutil.CheckStateMsg1(labelNameNode.IsLabelName(), "%s", labelNameNode)
		labelName := labelNameNode.GetString()
		AssertWithMessage("Duplicate label name: %s", labelName).
			That(f.outer.labeledStatementMap).
			DoesNotContainKey(labelName)
		f.outer.labeledStatementMap.Put(labelName, NewTypedScopeCreatorTest_LabeledStatement(n, scope))
	}
}

// Java: TypedScopeCreatorTest.getLabeledStatement(String)
func (c *TypedScopeCreatorTest) GetLabeledStatement(label string) *TypedScopeCreatorTest_LabeledStatement {
	AssertWithMessage("No statement found for label: %s", label).
		That(c.labeledStatementMap).
		ContainsKey(label)
	return c.labeledStatementMap.Get(label)
}

// Java: TypedScopeCreatorTest.getProcessor(Compiler)
func (c *TypedScopeCreatorTest) GetProcessor(compiler *Compiler) CompilerPass {
	c.registry = compiler.GetTypeRegistry()
	// Create a fresh statement map for each test case.
	c.labeledStatementMap = jutil.NewLinkedHashMap[string, *TypedScopeCreatorTest_LabeledStatement]()
	return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
		NewGatherModuleMetadata(compiler, false, ModuleLoader_ResolutionMode_BROWSER).Process(externs, root)
		NewModuleMapCreator(compiler, compiler.GetModuleMetadataMap()).Process(externs, root)
		NewInferConsts(compiler).Process(externs, root)

		scopeCreator := NewTypedScopeCreator(compiler)
		NewTypeInferencePass(compiler, compiler.GetReverseAbstractInterpreter(), scopeCreator).
			InferAllScopes(root.GetParent())
		NodeTraversal_Builder_().
			SetCompiler(compiler).
			SetCallback(&TypedScopeCreatorTest_ScopeFinder{outer: c}).
			SetScopeCreator(scopeCreator).
			TraverseRoots(externs, root)
	})
}

// Java: TypedScopeCreatorTest.getOptions()
func (c *TypedScopeCreatorTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetClosurePass(c.processClosurePrimitives)
	return options
}

// Java: TypedScopeCreatorTest.testVarDeclarationWithJSDocForObjPatWithOneVariable()
func (c *TypedScopeCreatorTest) TestVarDeclarationWithJSDocForObjPatWithOneVariable() {
	// Ignore JSDoc on a destructuring declaration, and just infer the type.
	// CheckJSDoc will issue a warning for the @type annotation.
	c.TestSame("/** @type {string} */ var {a} = {a: 1};")
	aVar := jutil.CheckNotNull(c.globalScope.GetVar("a"))
	TypeSubject_AssertType(aVar.GetType()).ToStringIsEqualTo("number")
	AssertThat(aVar.IsTypeInferred()).IsTrue()
}

// Java: TypedScopeCreatorTest.testVarDeclarationWithJSDocForObjPatWithMultipleVariables()
func (c *TypedScopeCreatorTest) TestVarDeclarationWithJSDocForObjPatWithMultipleVariables() {
	// Ignore JSDoc on a destructuring declaration
	// CheckJSDoc will issue a warning for the @type annotation.
	c.TestSame("/** @type {string} */ var {a, b} = {a: 1};")
	aVar := jutil.CheckNotNull(c.globalScope.GetVar("a"))
	TypeSubject_AssertType(aVar.GetType()).ToStringIsEqualTo("number")
	AssertThat(aVar.IsTypeInferred()).IsTrue()

	bVar := jutil.CheckNotNull(c.globalScope.GetVar("b"))
	TypeSubject_AssertType(bVar.GetType()).ToStringIsEqualTo("?")
	AssertThat(bVar.IsTypeInferred()).IsTrue()
}

// Java: TypedScopeCreatorTest.testVarDeclarationObjPatShorthandProp()
func (c *TypedScopeCreatorTest) TestVarDeclarationObjPatShorthandProp() {
	c.TestSame("var {/** number */ a} = {a: 1};")

	aVar := jutil.CheckNotNull(c.globalScope.GetVar("a"))
	TypeSubject_AssertType(aVar.GetType()).ToStringIsEqualTo("number")
	AssertThat(aVar.IsTypeInferred()).IsFalse()
}

// Java: TypedScopeCreatorTest.testVarDeclarationObjPatShorthandPropWithDefault()
func (c *TypedScopeCreatorTest) TestVarDeclarationObjPatShorthandPropWithDefault() {
	c.TestSame("var {/** number */ a = 2} = {a: 1};")

	aVar := jutil.CheckNotNull(c.globalScope.GetVar("a"))
	TypeSubject_AssertType(aVar.GetType()).ToStringIsEqualTo("number")
	AssertThat(aVar.IsTypeInferred()).IsFalse()
}

// Java: TypedScopeCreatorTest.testVarDeclarationObjPatNormalProp()
func (c *TypedScopeCreatorTest) TestVarDeclarationObjPatNormalProp() {
	c.TestSame("var {a: /** number */ a} = {a: 1};")

	aVar := jutil.CheckNotNull(c.globalScope.GetVar("a"))
	TypeSubject_AssertType(aVar.GetType()).ToStringIsEqualTo("number")
	AssertThat(aVar.IsTypeInferred()).IsFalse()
}

// Java: TypedScopeCreatorTest.testVarDeclarationObjPatNormalPropWithDefault()
func (c *TypedScopeCreatorTest) TestVarDeclarationObjPatNormalPropWithDefault() {
	c.TestSame("var {a: /** number */ a = 2} = {a: 1};")

	aVar := jutil.CheckNotNull(c.globalScope.GetVar("a"))
	TypeSubject_AssertType(aVar.GetType()).ToStringIsEqualTo("number")
	AssertThat(aVar.IsTypeInferred()).IsFalse()
}

// Java: TypedScopeCreatorTest.testVarDeclarationObjPatComputedProp()
func (c *TypedScopeCreatorTest) TestVarDeclarationObjPatComputedProp() {
	c.TestSame("var {['a']: /** number */ a} = {a: 1};")

	aVar := jutil.CheckNotNull(c.globalScope.GetVar("a"))
	TypeSubject_AssertType(aVar.GetType()).ToStringIsEqualTo("number")
	AssertThat(aVar.IsTypeInferred()).IsFalse()
}

// Java: TypedScopeCreatorTest.testVarDeclarationObjPatComputedPropWithDefault()
func (c *TypedScopeCreatorTest) TestVarDeclarationObjPatComputedPropWithDefault() {
	c.TestSame("var {['a']: /** number */ a = 2} = {a: 1};")

	aVar := jutil.CheckNotNull(c.globalScope.GetVar("a"))
	TypeSubject_AssertType(aVar.GetType()).ToStringIsEqualTo("number")
	AssertThat(aVar.IsTypeInferred()).IsFalse()
}

// Java: TypedScopeCreatorTest.testVarDeclarationArrayPat()
func (c *TypedScopeCreatorTest) TestVarDeclarationArrayPat() {
	c.TestSame("var [ /** number */ a ] = [1];")

	aVar := jutil.CheckNotNull(c.globalScope.GetVar("a"))
	TypeSubject_AssertType(aVar.GetType()).ToStringIsEqualTo("number")
	AssertThat(aVar.IsTypeInferred()).IsFalse()
}

// Java: TypedScopeCreatorTest.testVarDeclarationArrayPatWithDefault()
func (c *TypedScopeCreatorTest) TestVarDeclarationArrayPatWithDefault() {
	c.TestSame("var [ /** number */ a = 2 ] = [1];")

	aVar := jutil.CheckNotNull(c.globalScope.GetVar("a"))
	TypeSubject_AssertType(aVar.GetType()).ToStringIsEqualTo("number")
	AssertThat(aVar.IsTypeInferred()).IsFalse()
}

// Java: TypedScopeCreatorTest.testVarDeclarationArrayPatRest()
func (c *TypedScopeCreatorTest) TestVarDeclarationArrayPatRest() {
	// TODO(bradfordcsmith): Add a TypeCheck test case to ensure rest values are always Arrays
	c.TestSame("var [ ... /** !Array<number> */ a ] = [1];")

	aVar := jutil.CheckNotNull(c.globalScope.GetVar("a"))
	TypeSubject_AssertType(aVar.GetType()).ToStringIsEqualTo("Array<number>")
	AssertThat(aVar.IsTypeInferred()).IsFalse()
}

// Java: TypedScopeCreatorTest.testConstDeclarationObjectPatternInfersType_forAliasedConstructor()
func (c *TypedScopeCreatorTest) TestConstDeclarationObjectPatternInfersType_forAliasedConstructor() {
	c.TestSame(
		"const ns = {};\n" +
			"/** @constructor */\n" +
			"ns.Foo = function() {}\n" +
			"\n" +
			"const {Foo} = ns;\n" +
			"const /** !Foo */ fooInstance = new Foo();\n")

	fooVar := jutil.CheckNotNull(c.globalScope.GetVar("Foo"))
	TypeSubject_AssertType(fooVar.GetType()).ToStringIsEqualTo("(typeof ns.Foo)")
	AssertThat(fooVar.IsTypeInferred()).IsFalse()

	fooInstanceVar := jutil.CheckNotNull(c.globalScope.GetVar("fooInstance"))
	TypeSubject_AssertType(fooInstanceVar.GetType()).ToStringIsEqualTo("ns.Foo")
	AssertThat(fooInstanceVar.IsTypeInferred()).IsFalse()
}

// Java: TypedScopeCreatorTest.testConstDeclarationObjectPatternInfersType_forAliasedTypedef()
func (c *TypedScopeCreatorTest) TestConstDeclarationObjectPatternInfersType_forAliasedTypedef() {
	c.TestSame(
		"const ns = {};\n" +
			"/** @typedef {string} */ ns.Foo;\n" +
			"const {Foo} = ns;\n" +
			"let /** !Foo */ f = 'bar';\n")

	fooVar := jutil.CheckNotNull(c.globalScope.GetVar("Foo"))
	TypeSubject_AssertType(fooVar.GetType()).ToStringIsEqualTo("None")

	fooType := c.registry.GetGlobalType("Foo")
	TypeSubject_AssertType(fooType).IsEqualTo(c.GetNativeType(rhino.JSTypeNative_STRING_TYPE))

	fooInstanceVar := jutil.CheckNotNull(c.globalScope.GetVar("f"))
	TypeSubject_AssertType(fooInstanceVar.GetType()).ToStringIsEqualTo("string")
	AssertThat(fooInstanceVar.IsTypeInferred()).IsFalse()
}

// Java: TypedScopeCreatorTest.testConstDeclarationObjectPatternInfersTypeAsDeclared()
func (c *TypedScopeCreatorTest) TestConstDeclarationObjectPatternInfersTypeAsDeclared() {
	c.TestSame(
		"const /** {a: number} */ obj = {a: 3}; // preserve newline\n" +
			"const {a} = obj;\n")

	// we treat this as declaring a type on `a`
	aVar := jutil.CheckNotNull(c.globalScope.GetVar("a"))
	TypeSubject_AssertType(aVar.GetType()).ToStringIsEqualTo("number")
	AssertThat(aVar.IsTypeInferred()).IsFalse()
}

// Java: TypedScopeCreatorTest.testConstDeclarationObjectPatternInfersTypeGivenComputedProperty()
func (c *TypedScopeCreatorTest) TestConstDeclarationObjectPatternInfersTypeGivenComputedProperty() {
	c.TestSame(
		"const /** !IObject<string, number> */ obj = {a: 3}; // preserve newline\n" +
			"const {['foobar']: a} = obj;\n")

	aVar := jutil.CheckNotNull(c.globalScope.GetVar("a"))
	TypeSubject_AssertType(aVar.GetType()).ToStringIsEqualTo("number")
	AssertThat(aVar.IsTypeInferred()).IsTrue()
}

// Java: TypedScopeCreatorTest.testConstDeclarationObjectPatternInfersTypeGivenUnknownComputedProperty()
func (c *TypedScopeCreatorTest) TestConstDeclarationObjectPatternInfersTypeGivenUnknownComputedProperty() {
	c.TestSame(
		"var obj = {}; // preserve newline\n" +
			"const {['foobar']: a} = obj;\n")

	aVar := jutil.CheckNotNull(c.globalScope.GetVar("a"))
	TypeSubject_AssertType(aVar.GetType()).ToStringIsEqualTo("?")
	AssertThat(aVar.IsTypeInferred()).IsTrue()
}

// Java: TypedScopeCreatorTest.testConstDeclarationArrayPatternInfersType()
func (c *TypedScopeCreatorTest) TestConstDeclarationArrayPatternInfersType() {
	c.TestSame(
		"const /** !Iterable<number> */ arr = [1, 2, 3]; // preserve newline\n" +
			"const [a] = arr;\n")

	aVar := jutil.CheckNotNull(c.globalScope.GetVar("a"))
	TypeSubject_AssertType(aVar.GetType()).ToStringIsEqualTo("number")
	AssertThat(aVar.IsTypeInferred()).IsTrue()
}

// Java: TypedScopeCreatorTest.testConstDeclarationWithOrInRhs()
func (c *TypedScopeCreatorTest) TestConstDeclarationWithOrInRhs() {
	// needed because there is a special case for `var goog = goog || {};` that was crashing when
	// given a destructuring lhs.
	c.TestSame("let obj; const {a} = obj || {};")

	aVar := jutil.CheckNotNull(c.globalScope.GetVar("a"))
	TypeSubject_AssertType(aVar.GetType()).ToStringIsEqualTo("?")
	AssertThat(aVar.IsTypeInferred()).IsTrue()
}

// Java: TypedScopeCreatorTest.testConstDestructuringDeclarationWithUnknownTypeInExtendsClause()
func (c *TypedScopeCreatorTest) TestConstDestructuringDeclarationWithUnknownTypeInExtendsClause() {
	c.DisableTypeInfoValidation()
	c.TestWarning(
		c.Externs("var someUnknownExtern;"),
		c.Srcs("const {Parent} = someUnknownExtern; class Child extends Parent {}"),
		c.Warning(RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR))

	parentVar := jutil.CheckNotNull(c.globalScope.GetVar("Parent"))
	TypeSubject_AssertType(parentVar.GetType()).IsUnknown()

	childVar := jutil.CheckNotNull(c.globalScope.GetVar("Child"))
	childType := childVar.GetType().ToMaybeFunctionType()
	superclassCtor := childType.GetSuperClassConstructor()
	TypeSubject_AssertType(superclassCtor).IsNull()
}

// Java: TypedScopeCreatorTest.testConstDestructuringDeclarationWithUnknownPropertyInExtendsClause()
func (c *TypedScopeCreatorTest) TestConstDestructuringDeclarationWithUnknownPropertyInExtendsClause() {
	c.DisableTypeInfoValidation()
	c.TestWarning(
		c.Externs("var /** !Object */ someUnknownExtern;"),
		c.Srcs("const {Parent} = someUnknownExtern; class Child extends Parent {}"),
		c.Warning(RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR))

	parentVar := jutil.CheckNotNull(c.globalScope.GetVar("Parent"))
	TypeSubject_AssertType(parentVar.GetType()).IsUnknown()

	childVar := jutil.CheckNotNull(c.globalScope.GetVar("Child"))
	childType := childVar.GetType().ToMaybeFunctionType()
	superclassCtor := childType.GetSuperClassConstructor()
	TypeSubject_AssertType(superclassCtor).IsNull()
}

// TODO(bradfordcsmith): Add Object rest test case.
// Java: TypedScopeCreatorTest.testVarDeclarationNestedPatterns()
func (c *TypedScopeCreatorTest) TestVarDeclarationNestedPatterns() {
	c.TestSame(
		"var [\n" +
			"    {a: /** number */ a1},\n" +
			"    {a: /** string */ a2},\n" +
			"    ...{/** number */ length}\n" +
			"  ] = [{a: 1}, {a: '2'}, 1, 2, 3];\n")

	a1Var := jutil.CheckNotNull(c.globalScope.GetVar("a1"))
	TypeSubject_AssertType(a1Var.GetType()).ToStringIsEqualTo("number")
	AssertThat(a1Var.IsTypeInferred()).IsFalse()

	a2Var := jutil.CheckNotNull(c.globalScope.GetVar("a2"))
	TypeSubject_AssertType(a2Var.GetType()).ToStringIsEqualTo("string")
	AssertThat(a2Var.IsTypeInferred()).IsFalse()

	lengthVar := jutil.CheckNotNull(c.globalScope.GetVar("length"))
	TypeSubject_AssertType(lengthVar.GetType()).ToStringIsEqualTo("number")
	AssertThat(lengthVar.IsTypeInferred()).IsFalse()
}

// The following testAssign* tests check that we never treat qualified names in destructuring
// patterns as declared. CheckJSDoc will warn on those cases, so TypedScopeCreator just ignores
// them. The only way to 'declare' a qualified name is:
//
//	/** @type {number} */ a.b.c = rhs;
//
// Java: TypedScopeCreatorTest.testAssignWithJSDocForObjPatWithOneVariable()
func (c *TypedScopeCreatorTest) TestAssignWithJSDocForObjPatWithOneVariable() {
	// Ignore the JSDoc on the assignment
	c.TestSame("const ns = {}; (/** @type {number} */ {a: ns.a} = {a: 1});")

	aVar := c.globalScope.GetVar("ns.a")
	TypedVarSubject_AssertThat(aVar).IsNull()
}

// Java: TypedScopeCreatorTest.testAssignWithJSDocForObjPatWithMultipleVariables()
func (c *TypedScopeCreatorTest) TestAssignWithJSDocForObjPatWithMultipleVariables() {
	// Ignore the JSDoc on the assignment
	c.TestSame("const ns = {}; (/** @type {number} */ {a: ns.a, b: ns.b} = {a: 1});")

	aVar := c.globalScope.GetVar("ns.a")
	TypedVarSubject_AssertThat(aVar).IsNull()
}

// Java: TypedScopeCreatorTest.testAssignObjPatNormalProp()
func (c *TypedScopeCreatorTest) TestAssignObjPatNormalProp() {
	// CheckJSDoc will warn on the inline type annotation here, typechecking just ignores it.
	c.TestSame("const ns = {}; ({a: /** number */ ns.a} = {a: 1});")

	aVar := c.globalScope.GetVar("ns.a")
	TypedVarSubject_AssertThat(aVar).IsNull()
}

// Java: TypedScopeCreatorTest.testAssignObjPatComputedProp()
func (c *TypedScopeCreatorTest) TestAssignObjPatComputedProp() {
	// CheckJSDoc will warn on the inline type annotation here, typechecking just ignores it.
	c.TestSame("const ns = {}; ({['a']: /** number */ ns.a} = {a: 1});")

	aVar := c.globalScope.GetVar("ns.a")
	TypedVarSubject_AssertThat(aVar).IsNull()
}

// Java: TypedScopeCreatorTest.testAssignArrayPatWithJSDocOnAssign()
func (c *TypedScopeCreatorTest) TestAssignArrayPatWithJSDocOnAssign() {
	c.TestSame("const ns = {}; /** @type {number} */ [ ns.a ] = [1];")

	aVar := c.globalScope.GetVar("ns.a")
	TypedVarSubject_AssertThat(aVar).IsNull()
}

// Java: TypedScopeCreatorTest.testAssignArrayPatWithQualifiedName()
func (c *TypedScopeCreatorTest) TestAssignArrayPatWithQualifiedName() {
	// CheckJSDoc will warn on the inline type annotation here, typechecking just ignores it.
	c.TestSame("const ns = {}; [ /** number */ ns.a ] = [1];")

	aVar := c.globalScope.GetVar("ns.a")
	TypedVarSubject_AssertThat(aVar).IsNull()
}

// Java: TypedScopeCreatorTest.testAssignArrayPatWithQualifiedNameAndDefaultValue()
func (c *TypedScopeCreatorTest) TestAssignArrayPatWithQualifiedNameAndDefaultValue() {
	// CheckJSDoc will warn on the inline type annotation here, typechecking just ignores it.
	c.TestSame("const ns = {}; [ /** number */ ns.a = 1 ] = [];")

	aVar := c.globalScope.GetVar("ns.a")
	TypedVarSubject_AssertThat(aVar).IsNull()
}

// Java: TypedScopeCreatorTest.testForOfWithObjectPatVarDeclarationWithShorthand()
func (c *TypedScopeCreatorTest) TestForOfWithObjectPatVarDeclarationWithShorthand() {
	c.TestSame("for (var {/** number */ a} of {}) {}")

	aVar := jutil.CheckNotNull(c.globalScope.GetVar("a"))
	TypeSubject_AssertType(aVar.GetType()).ToStringIsEqualTo("number")
	AssertThat(aVar.IsTypeInferred()).IsFalse()
}

// Java: TypedScopeCreatorTest.testForOfWithArrayPatVarDeclaration()
func (c *TypedScopeCreatorTest) TestForOfWithArrayPatVarDeclaration() {
	c.TestSame("for (var [/** number */ a] of []) {}")

	aVar := jutil.CheckNotNull(c.globalScope.GetVar("a"))
	TypeSubject_AssertType(aVar.GetType()).ToStringIsEqualTo("number")
	AssertThat(aVar.IsTypeInferred()).IsFalse()
}

// Java: TypedScopeCreatorTest.testCastOnLhsDoesntDeclareProperty()
func (c *TypedScopeCreatorTest) TestCastOnLhsDoesntDeclareProperty() {
	c.TestSame("const ns = {}; /** @type {null} */ (ns.a) = null;")

	TypedVarSubject_AssertThat(c.globalScope.GetVar("ns.a")).IsNull()
}

// Java: TypedScopeCreatorTest.testClassInstancePropertyInArrowFunction()
func (c *TypedScopeCreatorTest) TestClassInstancePropertyInArrowFunction() {
	c.TestSame(
		"class C {\n" +
			"  constructor() {\n" +
			"    (() => {\n" +
			"      /** @type {number} */\n" +
			"      this.size = 0;\n" +
			"    })();\n" +
			"    CTOR: 0;\n" +
			"  }\n" +
			"}\n")

	// Declare `this.size` in the enclosing constructor's scope.
	ScopeSubject_AssertScope(c.GetLabeledStatement("CTOR").enclosingScope.GetParent()).
		Declares("this.size").
		Directly()
}

// Java: TypedScopeCreatorTest.testObjectLiteralShadowingFnExpressionName()
func (c *TypedScopeCreatorTest) TestObjectLiteralShadowingFnExpressionName() {
	c.TestSame(
		"var obj = {\n" +
			"  key: function f() {\n" +
			"    var f = {x: 0};\n" +
			"    FN: 0;\n" +
			"  }\n" +
			"};\n")

	// Test that we declare 'f.x' on the function block scope, not the enclosing function scope,
	// which also contains `f`.
	ScopeSubject_AssertScope(c.GetLabeledStatement("FN").enclosingScope).Declares("f.x").Directly()
}

// Java: TypedScopeCreatorTest.testForOfWithObjectPatConstDeclarationShadowedInLoop()
func (c *TypedScopeCreatorTest) TestForOfWithObjectPatConstDeclarationShadowedInLoop() {
	c.TestSame(
		"for (const {/** number */ a} of {}) {\n" +
			"  const /** string */ a = 'foo';\n" +
			"  IN_LOOP: a;\n" +
			"}\n")

	ScopeSubject_AssertScope(c.globalScope).DoesNotDeclare("a")
	loopBlockScope := c.GetLabeledStatement("IN_LOOP").enclosingScope
	loopInitializerScope := loopBlockScope.GetParent()
	AssertThat(loopInitializerScope).IsNotEqualTo(c.globalScope)

	aVarloopBlock := loopBlockScope.GetVar("a")
	TypeSubject_AssertType(aVarloopBlock.GetType()).ToStringIsEqualTo("string")
	AssertThat(aVarloopBlock.IsTypeInferred()).IsFalse()

	aVarLoopInit := loopInitializerScope.GetVar("a")
	TypeSubject_AssertType(aVarLoopInit.GetType()).ToStringIsEqualTo("number")
	AssertThat(aVarLoopInit.IsTypeInferred()).IsFalse()
}

// Java: TypedScopeCreatorTest.testDeclarativelyUnboundVarsWithoutTypes()
func (c *TypedScopeCreatorTest) TestDeclarativelyUnboundVarsWithoutTypes() {
	c.TestSame(
		"var uninitializedVar;\n" +
			"let uninitializedLet;\n" +
			"\n" +
			"/** @type {?} */ var uninitializedVarWithType;\n" +
			"/** @type {?} */ let uninitializedLetWithType;\n" +
			"\n" +
			"/** @typedef {number} */ var typedefVar;\n" +
			"/** @typedef {number} */ let typedefLet;\n" +
			"\n" +
			"var initializedVar = 1;\n" +
			"let initializedLet = 1;\n" +
			"\n" +
			"/** @type {number} */ var initializedVarWithType = 1;\n" +
			"/** @type {number} */ let initializedLetWithType = 1;\n" +
			"\n" +
			"/** @const */ var CONST_VAR = 1;\n" +
			"/** @const */ let CONST_LET = 1;\n" +
			"const CONST = 1;\n" +
			"\n" +
			"/** @const {number} */ var   CONST_VAR_WITH_TYPE = 1;\n" +
			"/** @const {number} */ let   CONST_LET_WITH_TYPE = 1;\n" +
			"/** @type  {number} */ const CONST_WITH_TYPE     = 1;\n")

	// After TypeInference, there should be no more variables in this set.
	AssertThat(slices.AppendSeq([]*TypedVar{}, c.globalScope.GetDeclarativelyUnboundVarsWithoutTypes())).IsEmpty()

	TypedVarSubject_AssertThat(c.globalScope.GetVar("uninitializedVar")).HasJSTypeThat().IsUnknown()
	TypedVarSubject_AssertThat(c.globalScope.GetVar("uninitializedLet")).HasJSTypeThat().IsUnknown()
}

// Java: TypedScopeCreatorTest.testRestParameters()
func (c *TypedScopeCreatorTest) TestRestParameters() {
	c.TestSame(
		"/** // preserve newlines\n" +
			" * @param {string} str\n" +
			" * @param {...number} nums\n" +
			" */\n" +
			"function doSomething(str, ...nums) {\n" +
			"  FUNCTION_BODY: 0;\n" +
			"}\n")
	functionBodyScope := c.GetLabeledStatement("FUNCTION_BODY").enclosingScope
	ScopeSubject_AssertScope(functionBodyScope).
		Declares("str").
		OnClosestContainerScope().
		WithTypeThat().
		ToStringIsEqualTo("string")
	ScopeSubject_AssertScope(functionBodyScope).
		Declares("nums").
		OnClosestContainerScope().
		WithTypeThat().
		ToStringIsEqualTo("Array<number>")
}

// Java: TypedScopeCreatorTest.testRestObjectPatternParameters()
func (c *TypedScopeCreatorTest) TestRestObjectPatternParameters() {
	c.TestSame(
		c.Externs("/** @type {number} */ Array.prototype.length"),
		c.Srcs(
			"/** // preserve newlines\n"+
				" * @param {...string} strs\n"+
				" */\n"+
				"function doSomething(...{length}) {\n"+
				"  FUNCTION_BODY: 0;\n"+
				"}\n"))
	functionBodyScope := c.GetLabeledStatement("FUNCTION_BODY").enclosingScope
	ScopeSubject_AssertScope(functionBodyScope).
		Declares("length").
		OnClosestContainerScope().
		WithTypeThat().
		ToStringIsEqualTo("number")
}

// Java: TypedScopeCreatorTest.testDefaultParameterFullJSDoc()
func (c *TypedScopeCreatorTest) TestDefaultParameterFullJSDoc() {
	c.TestSame(
		"/** // preserve newlines\n" +
			" * @param {string=} str\n" +
			" */\n" +
			"function doSomething(str = '') {\n" +
			"  FUNCTION_BODY: 0;\n" +
			"}\n")
	functionBodyScope := c.GetLabeledStatement("FUNCTION_BODY").enclosingScope
	ScopeSubject_AssertScope(functionBodyScope).
		Declares("str").
		OnClosestContainerScope().
		WithTypeThat().
		// TODO(b/117162687): this should just be `string`
		ToStringIsEqualTo("(string|undefined)")
}

// Java: TypedScopeCreatorTest.testDefaultParameterInlineJSDoc()
func (c *TypedScopeCreatorTest) TestDefaultParameterInlineJSDoc() {
	c.TestSame(
		" // preserve newlines\n" +
			"function doSomething(/** string= */ str = '') {\n" +
			"  FUNCTION_BODY: 0;\n" +
			"}\n")
	functionBodyScope := c.GetLabeledStatement("FUNCTION_BODY").enclosingScope
	ScopeSubject_AssertScope(functionBodyScope).
		Declares("str").
		OnClosestContainerScope().
		WithTypeThat().
		// TODO(b/117162687): this should just be `string`
		ToStringIsEqualTo("(string|undefined)")
}

// Java: TypedScopeCreatorTest.testDefaultParameterNoJSDoc()
func (c *TypedScopeCreatorTest) TestDefaultParameterNoJSDoc() {
	c.TestSame(
		" // preserve newlines\n" +
			"function doSomething(str = '') {\n" +
			"  FUNCTION_BODY: 0;\n" +
			"}\n")
	functionBodyScope := c.GetLabeledStatement("FUNCTION_BODY").enclosingScope
	ScopeSubject_AssertScope(functionBodyScope).
		Declares("str").
		OnClosestContainerScope().
		WithTypeThat().
		ToStringIsEqualTo("?")
}

// Java: TypedScopeCreatorTest.testDefaultParameterNoJSDocInferredToBeOptional()
func (c *TypedScopeCreatorTest) TestDefaultParameterNoJSDocInferredToBeOptional() {
	c.TestSame("function f(a = 3) {}")

	ScopeSubject_AssertScope(c.globalScope).
		Declares("f").
		WithTypeThat().
		ToStringIsEqualTo("function(?=): undefined")
}

// Java: TypedScopeCreatorTest.testDefaultParameterConflictingJSDoc()
func (c *TypedScopeCreatorTest) TestDefaultParameterConflictingJSDoc() {
	c.TestSame(
		"/** // preserve newlines\n" +
			" * @param {number=} str\n" +
			" */\n" +
			"function doSomething(/** string= */ str = '') {\n" +
			"  FUNCTION_BODY: 0;\n" +
			"}\n")
	functionBodyScope := c.GetLabeledStatement("FUNCTION_BODY").enclosingScope
	ScopeSubject_AssertScope(functionBodyScope).
		Declares("str").
		OnClosestContainerScope().
		WithTypeThat().
		// TODO(b/111523967): Should report an error when header and inline JSDoc types
		//     conflict.
		ToStringIsEqualTo("(number|undefined)")
}

// Java: TypedScopeCreatorTest.testDefaultParameterSetToUndefinedIsDeclaredAsPossiblyUndefined()
func (c *TypedScopeCreatorTest) TestDefaultParameterSetToUndefinedIsDeclaredAsPossiblyUndefined() {
	c.TestSame(
		"/** @param {number=} num */\n" +
			"function doSomething(num = undefined) {\n" +
			"  FUNCTION_BODY: 0;\n" +
			"}\n")
	functionBodyScope := c.GetLabeledStatement("FUNCTION_BODY").enclosingScope
	ScopeSubject_AssertScope(functionBodyScope).
		Declares("num").
		OnClosestContainerScope().
		WithTypeThat().
		ToStringIsEqualTo("(number|undefined)")
}

// Java: TypedScopeCreatorTest.testDefaultParameterFullJSDoc_setToUndefined()
func (c *TypedScopeCreatorTest) TestDefaultParameterFullJSDoc_setToUndefined() {
	c.TestSame(
		"/** // preserve newlines\n" +
			" * @param {(string|undefined)=} str // NOTE: we just drop the |undefined here\n" +
			" */\n" +
			"function doSomething(str = '') {\n" +
			"  FUNCTION_BODY: 0;\n" +
			"}\n")
	functionBodyScope := c.GetLabeledStatement("FUNCTION_BODY").enclosingScope
	ScopeSubject_AssertScope(functionBodyScope).
		Declares("str").
		OnClosestContainerScope().
		WithTypeThat().
		ToStringIsEqualTo("(string|undefined)")
}

// Java: TypedScopeCreatorTest.testDefaultParameterInferredNotUndefinedInCallback()
func (c *TypedScopeCreatorTest) TestDefaultParameterInferredNotUndefinedInCallback() {
	c.TestSame(
		"function takesCallback(/** function(string=): ? */ cb) {}\n" +
			"\n" +
			"takesCallback((str = '') => {})\n")

	strVar := jutil.CheckNotNull(c.lastFunctionScope.GetVar("str"))
	TypeSubject_AssertType(strVar.GetType()).IsString()
	AssertThat(strVar.IsTypeInferred()).IsTrue()

	strType := c.FindNameType("str", c.globalScope)
	TypeSubject_AssertType(strType).IsString() // the actual parameter is also typed as not-undefined
}

// Java: TypedScopeCreatorTest.testDefaultParameterInferredNotUndefinedInCallbackButLaterSetToUndefined()
func (c *TypedScopeCreatorTest) TestDefaultParameterInferredNotUndefinedInCallbackButLaterSetToUndefined() {
	c.TestSame(
		"function takesCallback(/** function(string=): ? */ cb) {}\n" +
			"\n" +
			"takesCallback((str = '') => {\n" +
			"  str = undefined;\n" +
			"})\n")

	strVar := jutil.CheckNotNull(c.lastFunctionScope.GetVar("str"))
	TypeSubject_AssertType(strVar.GetType()).ToStringIsEqualTo("(string|undefined)")
	AssertThat(strVar.IsTypeInferred()).IsTrue()
}

// Java: TypedScopeCreatorTest.testDefaultDestructuringParameterFullJSDoc()
func (c *TypedScopeCreatorTest) TestDefaultDestructuringParameterFullJSDoc() {
	c.TestSame(
		"/** @param {{str: (string|undefined)}=} data */\n" +
			"function f({str = ''} = {}) {\n" +
			"  FUNCTION_BODY: 0;\n" +
			"}\n")
	functionBodyScope := c.GetLabeledStatement("FUNCTION_BODY").enclosingScope
	ScopeSubject_AssertScope(functionBodyScope).
		Declares("str").
		OnClosestContainerScope().
		WithTypeThat().
		ToStringIsEqualTo("(string|undefined)")
}

// Java: TypedScopeCreatorTest.testDestructuringParameterWithNoJSDoc()
func (c *TypedScopeCreatorTest) TestDestructuringParameterWithNoJSDoc() {
	c.TestSame("function f([x, y], {z}) {}")

	fVar := jutil.CheckNotNull(c.globalScope.GetVar("f"))
	TypeSubject_AssertType(fVar.GetType()).ToStringIsEqualTo("function(?, ?): undefined")
	AssertThat(fVar.IsTypeInferred()).IsFalse()
}

// Java: TypedScopeCreatorTest.testArrayPatternParameterWithFullJSDoc()
func (c *TypedScopeCreatorTest) TestArrayPatternParameterWithFullJSDoc() {
	c.TestSame(
		"/**\n" +
			" * @param {string} x\n" +
			" * @param {!Iterable<number>} arr\n" +
			" */\n" +
			"function f(x, [y]) {}\n")

	fVar := jutil.CheckNotNull(c.globalScope.GetVar("f"))
	TypeSubject_AssertType(fVar.GetType()).
		ToStringIsEqualTo("function(string, Iterable<number,?,?>): undefined")
	AssertThat(fVar.IsTypeInferred()).IsFalse()

	xVar := jutil.CheckNotNull(c.lastFunctionScope.GetVar("x"))
	TypeSubject_AssertType(xVar.GetType()).ToStringIsEqualTo("string")
	AssertThat(xVar.IsTypeInferred()).IsFalse()

	yVar := jutil.CheckNotNull(c.lastFunctionScope.GetVar("y"))
	TypeSubject_AssertType(yVar.GetType()).ToStringIsEqualTo("number")
	AssertThat(yVar.IsTypeInferred()).IsFalse()
}

// Java: TypedScopeCreatorTest.testArrayPatternParameterWithRestWithFullJSDoc()
func (c *TypedScopeCreatorTest) TestArrayPatternParameterWithRestWithFullJSDoc() {
	c.TestSame("/** @param {!Iterable<number>} arr */ function f([x, ...y]) {}")

	fVar := jutil.CheckNotNull(c.globalScope.GetVar("f"))
	TypeSubject_AssertType(fVar.GetType()).ToStringIsEqualTo("function(Iterable<number,?,?>): undefined")
	AssertThat(fVar.IsTypeInferred()).IsFalse()

	xVar := jutil.CheckNotNull(c.lastFunctionScope.GetVar("x"))
	TypeSubject_AssertType(xVar.GetType()).ToStringIsEqualTo("number")
	AssertThat(xVar.IsTypeInferred()).IsFalse()

	yVar := jutil.CheckNotNull(c.lastFunctionScope.GetVar("y"))
	TypeSubject_AssertType(yVar.GetType()).ToStringIsEqualTo("Array<number>")
	AssertThat(yVar.IsTypeInferred()).IsFalse()
}

// Java: TypedScopeCreatorTest.testArrayPatternParameterWithDefaultName()
func (c *TypedScopeCreatorTest) TestArrayPatternParameterWithDefaultName() {
	c.TestSame(
		"const /** !Iterable<number> */ iter = [0];\n" +
			"\n" +
			"/** @param {!Iterable<number>=} p */\n" +
			"function f([x] = iter) {\n" +
			"}\n")

	// JSType on the actual node
	xType := c.FindNameType("x", c.lastFunctionScope)
	TypeSubject_AssertType(xType).IsNumber()

	// JSType in the scope
	xVar := jutil.CheckNotNull(c.lastFunctionScope.GetVar("x"))
	TypeSubject_AssertType(xVar.GetType()).ToStringIsEqualTo("number")
	AssertThat(xVar.IsTypeInferred()).IsTrue()

	// JSType on the array pattern
	arrayPatternType := c.FindTokenType(rhino.Token_ARRAY_PATTERN, c.globalScope)
	TypeSubject_AssertType(arrayPatternType).ToStringIsEqualTo("Iterable<number,?,?>")
}

// Java: TypedScopeCreatorTest.testArrayDestructuringPatternParameterWithDefaultArray()
func (c *TypedScopeCreatorTest) TestArrayDestructuringPatternParameterWithDefaultArray() {
	c.TestSame(
		"/** @param {!Iterable<number>=} p */\n" +
			"function f([x] = [0]) {}\n")

	xType := c.FindNameType("x", c.lastFunctionScope)
	// This is unknown because we infer `[0]` to have type `!Array<?>`
	TypeSubject_AssertType(xType).IsUnknown()
}

// Java: TypedScopeCreatorTest.testObjectPatternParameterWithFullJSDoc()
func (c *TypedScopeCreatorTest) TestObjectPatternParameterWithFullJSDoc() {
	c.TestSame("/** @param {{a: string, b: number}} arr */ function f({a, b}) {}")

	fVar := jutil.CheckNotNull(c.globalScope.GetVar("f"))
	TypeSubject_AssertType(fVar.GetType()).
		ToStringIsEqualTo("function({\n  a: string,\n  b: number\n}): undefined")
	AssertThat(fVar.IsTypeInferred()).IsFalse()

	aVar := jutil.CheckNotNull(c.lastFunctionScope.GetVar("a"))
	TypeSubject_AssertType(aVar.GetType()).ToStringIsEqualTo("string")
	AssertThat(aVar.IsTypeInferred()).IsFalse()

	bVar := jutil.CheckNotNull(c.lastFunctionScope.GetVar("b"))
	TypeSubject_AssertType(bVar.GetType()).ToStringIsEqualTo("number")
	AssertThat(bVar.IsTypeInferred()).IsFalse()
}

// Java: TypedScopeCreatorTest.testObjectPatternParameterWithUnknownPropertyWithFullJSDoc()
func (c *TypedScopeCreatorTest) TestObjectPatternParameterWithUnknownPropertyWithFullJSDoc() {
	c.TestSame("/** @param {{a: string}} arr */ function f({a, b}) {}")

	fVar := jutil.CheckNotNull(c.globalScope.GetVar("f"))
	TypeSubject_AssertType(fVar.GetType()).ToStringIsEqualTo("function({a: string}): undefined")
	AssertThat(fVar.IsTypeInferred()).IsFalse()

	aVar := jutil.CheckNotNull(c.lastFunctionScope.GetVar("a"))
	TypeSubject_AssertType(aVar.GetType()).ToStringIsEqualTo("string")
	AssertThat(aVar.IsTypeInferred()).IsFalse()

	bVar := jutil.CheckNotNull(c.lastFunctionScope.GetVar("b"))
	TypeSubject_AssertType(bVar.GetType()).ToStringIsEqualTo("?")
	AssertThat(bVar.IsTypeInferred()).IsTrue()
}

// Note: these tests expect the parameters to have inferred, not declared, types because of the
// order in which scope creation & named type resolution happens.
//
// <p>See also b/118710352
// Java: TypedScopeCreatorTest.testObjectPatternParameter_withUnresolvedForwardDeclaredType_isInferred()
func (c *TypedScopeCreatorTest) TestObjectPatternParameter_withUnresolvedForwardDeclaredType_isInferred() {
	// Note that this only emits one UNRECOGNIZED_TYPE_ERROR for SomeUnknownName. We don't emit
	// an error for SomeUnknownName#a, as that would be redundant.
	c.Test(
		c.Srcs(
			"goog.forwardDeclare('SomeUnknownName');\n"+
				"/** @param {!SomeUnknownName} obj */ function f({a}) {}\n"),
		c.Warning(RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR))

	aVar := jutil.CheckNotNull(c.lastFunctionScope.GetVar("a"))
	TypeSubject_AssertType(aVar.GetType()).ToStringIsEqualTo("?")
	AssertThat(aVar.IsTypeInferred()).IsTrue()
}

// Java: TypedScopeCreatorTest.testObjectPatternParameterWithFullJSDoc_withMissingProperty_isInferred()
func (c *TypedScopeCreatorTest) TestObjectPatternParameterWithFullJSDoc_withMissingProperty_isInferred() {
	c.TestSame(
		"class SomeName {}\n" +
			"/** @param {!SomeName} obj */ function f({a}) {}\n")

	aVar := jutil.CheckNotNull(c.lastFunctionScope.GetVar("a"))
	TypeSubject_AssertType(aVar.GetType()).ToStringIsEqualTo("?")
	// we still consider this a declared type
	AssertThat(aVar.IsTypeInferred()).IsTrue()
}

// Java: TypedScopeCreatorTest.testNestedObjectPatternParameter_withNamedType()
func (c *TypedScopeCreatorTest) TestNestedObjectPatternParameter_withNamedType() {
	c.TestSame(
		"/** @param {!SomeName} obj */ function f({data: {a}}) { A: a; }\n" +
			"class SomeOtherName {\n" +
			"  constructor() {\n" +
			"    /** @type {number} */\n" +
			"    this.a;\n" +
			"  }\n" +
			"}\n" +
			"class SomeName {\n" +
			"  constructor() {\n" +
			"    /** @type {!SomeOtherName} */\n" +
			"    this.data;\n" +
			"  }\n" +
			"}\n")

	// This variable becomes inferred, even though SomeName is declared before it is referenced,
	// just because the SomeName in the @param is still unresolved at type scope creation time.
	aVar := jutil.CheckNotNull(c.GetLabeledStatement("A").enclosingScope.GetVar("a"))
	TypeSubject_AssertType(aVar.GetType()).ToStringIsEqualTo("number")
	AssertThat(aVar.IsTypeInferred()).IsTrue()
}

// Java: TypedScopeCreatorTest.testTypecheckingTemplatizedTypesForObjectParameter()
func (c *TypedScopeCreatorTest) TestTypecheckingTemplatizedTypesForObjectParameter() {
	c.TestSame(
		"/**\n" +
			" * @constructor\n" +
			" * @template T\n" +
			" */\n" +
			"function TemplatizedClass() {\n" +
			"  /** @const {T} */\n" +
			"  this.data;\n" +
			"}\n" +
			"/** @param {!TemplatizedClass<!SomeName>} obj */\n" +
			"function f({data}) { A: data.a; }\n" +
			"\n" +
			"class SomeName {\n" +
			"  constructor() {\n" +
			"    /** @type {number} */\n" +
			"    this.a;\n" +
			"  }\n" +
			"}\n")

	aStatement := c.GetLabeledStatement("A")
	dataVar := jutil.CheckNotNull(aStatement.enclosingScope.GetVar("data"))

	TypeSubject_AssertType(dataVar.GetType()).ToStringIsEqualTo("SomeName")
	TypeSubject_AssertType(aStatement.statementNode.GetOnlyChild().GetJSType()).ToStringIsEqualTo("number")
}

// Java: TypedScopeCreatorTest.testTypecheckingAliasOfTypedefType_forObjectParameter()
func (c *TypedScopeCreatorTest) TestTypecheckingAliasOfTypedefType_forObjectParameter() {
	c.TestSame(
		"class Foo {\n" +
			"  constructor() {\n" +
			"    /** @type {number} */\n" +
			"    this.a;\n" +
			"  }\n" +
			"}\n" +
			"/** @typedef {!Foo} */\n" +
			"var TypedefOfFoo;\n" +
			"\n" +
			"/** @param {!ns.TypedefOfFoo} obj */\n" +
			"function f({a}) { A: a; }\n" +
			"const ns = {};\n" +
			"/** @const */\n" +
			"ns.TypedefOfFoo = TypedefOfFoo;\n")

	aVar := jutil.CheckNotNull(c.GetLabeledStatement("A").enclosingScope.GetVar("a"))
	TypeSubject_AssertType(aVar.GetType()).ToStringIsEqualTo("number")
	AssertThat(aVar.IsTypeInferred()).IsTrue()
}

// Java: TypedScopeCreatorTest.testArrayPatternParameter_withNamedType()
func (c *TypedScopeCreatorTest) TestArrayPatternParameter_withNamedType() {
	c.TestSame(
		"/** @param {!ASubclass} obj */ function f([a]) { A: a; }\n" +
			"/**\n" +
			" * @interface\n" +
			" * @extends {Iterable<T>}\n" +
			" * @template T\n" +
			" */\n" +
			"function TemplatizedClass() {\n" +
			"}\n" +
			"/** @interface @extends {TemplatizedClass<number>} */\n" +
			"class ASubclass {}\n")

	// Verify we get the correct type for 'a' after name resolution is complete
	aVar := jutil.CheckNotNull(c.GetLabeledStatement("A").enclosingScope.GetVar("a"))
	TypeSubject_AssertType(aVar.GetType()).ToStringIsEqualTo("number")
	AssertThat(aVar.IsTypeInferred()).IsTrue()
}

// Java: TypedScopeCreatorTest.testObjectPatternParameter_withNamedType_andRest()
func (c *TypedScopeCreatorTest) TestObjectPatternParameter_withNamedType_andRest() {
	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_ECMASCRIPT_2018)
	c.TestSame(
		"class SomeUnknownName {\n" +
			"  constructor() {\n" +
			"    /** @type {number} */\n" +
			"    this.data;\n" +
			"  }\n" +
			"}\n" +
			"/** @param {!SomeUnknownName} obj */ function f({...rest}) { REST: rest; }\n")

	aVar := jutil.CheckNotNull(c.GetLabeledStatement("REST").enclosingScope.GetVar("rest"))
	TypeSubject_AssertType(aVar.GetType()).ToStringIsEqualTo("Object")
	AssertThat(aVar.IsTypeInferred()).IsFalse()
}

// Java: TypedScopeCreatorTest.testMixedParametersWithFullJSDoc()
func (c *TypedScopeCreatorTest) TestMixedParametersWithFullJSDoc() {
	c.TestSame(
		"/**\n" +
			" * @param {string} x\n" +
			" * @param {{a: !Iterable<number>}} obj\n" +
			" * @param {!Iterable<{z: null}>} arr\n" +
			" */\n" +
			"function f(x, {a: [y]}, [{z}]) {}\n")

	fVar := jutil.CheckNotNull(c.globalScope.GetVar("f"))
	TypeSubject_AssertType(fVar.GetType()).
		ToStringIsEqualTo(
			"function(string, {a: Iterable<number,?,?>}, Iterable<{z: null},?,?>): undefined")
	AssertThat(fVar.IsTypeInferred()).IsFalse()

	xVar := jutil.CheckNotNull(c.lastFunctionScope.GetVar("x"))
	TypeSubject_AssertType(xVar.GetType()).ToStringIsEqualTo("string")
	AssertThat(xVar.IsTypeInferred()).IsFalse()

	yVar := jutil.CheckNotNull(c.lastFunctionScope.GetVar("y"))
	TypeSubject_AssertType(yVar.GetType()).ToStringIsEqualTo("number")
	AssertThat(yVar.IsTypeInferred()).IsFalse()

	zVar := jutil.CheckNotNull(c.lastFunctionScope.GetVar("z"))
	TypeSubject_AssertType(zVar.GetType()).ToStringIsEqualTo("null")
	AssertThat(zVar.IsTypeInferred()).IsFalse()
}

// Java: TypedScopeCreatorTest.testObjectPatternParameterWithComputedPropertyWithFullJSDoc()
func (c *TypedScopeCreatorTest) TestObjectPatternParameterWithComputedPropertyWithFullJSDoc() {
	c.TestSame(
		"/**\n" +
			" * @param {string} x\n" +
			" * @param {!Object<string, number>} arr\n" +
			" */\n" +
			"function f(x, {['foobar' + 3]: a}) {}\n")

	xVar := jutil.CheckNotNull(c.lastFunctionScope.GetVar("x"))
	TypeSubject_AssertType(xVar.GetType()).ToStringIsEqualTo("string")
	AssertThat(xVar.IsTypeInferred()).IsFalse()

	aVar := jutil.CheckNotNull(c.lastFunctionScope.GetVar("a"))
	TypeSubject_AssertType(aVar.GetType()).ToStringIsEqualTo("number")
	AssertThat(aVar.IsTypeInferred()).IsFalse()
}

// Java: TypedScopeCreatorTest.testDestructuringParametersInIifeInfersType_withNameArguments()
func (c *TypedScopeCreatorTest) TestDestructuringParametersInIifeInfersType_withNameArguments() {
	c.TestSame(
		"const /** {x:  number} */ data = {x: 3};\n" +
			"const /** !Iterable<string> */ strings = ['foo', 'bar'];\n" +
			"\n" +
			"(function ({x}, [y]) {})(data, strings);\n")

	xVar := jutil.CheckNotNull(c.lastFunctionScope.GetVar("x"))
	TypeSubject_AssertType(xVar.GetType()).ToStringIsEqualTo("number")
	AssertThat(xVar.IsTypeInferred()).IsTrue()

	yVar := jutil.CheckNotNull(c.lastFunctionScope.GetVar("y"))
	TypeSubject_AssertType(yVar.GetType()).ToStringIsEqualTo("string")
	AssertThat(yVar.IsTypeInferred()).IsTrue()
}

// Java: TypedScopeCreatorTest.testDestructuringParametersInIifeInfersType_withLiteralArguments()
func (c *TypedScopeCreatorTest) TestDestructuringParametersInIifeInfersType_withLiteralArguments() {
	c.TestSame("(function ({x}) {})({x: 3});")

	xVar := jutil.CheckNotNull(c.lastFunctionScope.GetVar("x"))
	TypeSubject_AssertType(xVar.GetType()).ToStringIsEqualTo("number")
	AssertThat(xVar.IsTypeInferred()).IsTrue()
}

// Java: TypedScopeCreatorTest.testDestructuringParametersInIifeInfersType_withLiteralArgumentsAndDefaultValue()
func (c *TypedScopeCreatorTest) TestDestructuringParametersInIifeInfersType_withLiteralArgumentsAndDefaultValue() {
	c.TestSame("(function ({x = 'bar' + 'baz'}) {})({x: true ? 3 : undefined});")

	xVar := jutil.CheckNotNull(c.lastFunctionScope.GetVar("x"))
	TypeSubject_AssertType(xVar.GetType()).ToStringIsEqualTo("(number|string)")
	AssertThat(xVar.IsTypeInferred()).IsTrue()
}

// Java: TypedScopeCreatorTest.testDestructuringParametersInCallbackInfersType()
func (c *TypedScopeCreatorTest) TestDestructuringParametersInCallbackInfersType() {
	c.TestSame(
		"function f(/** function({x: number}, !Iterable<string>) */ callback) {}\n" +
			"\n" +
			"f(function ({x}, [y]) {});\n")

	xVar := jutil.CheckNotNull(c.lastFunctionScope.GetVar("x"))
	TypeSubject_AssertType(xVar.GetType()).ToStringIsEqualTo("number")
	AssertThat(xVar.IsTypeInferred()).IsTrue()

	yVar := jutil.CheckNotNull(c.lastFunctionScope.GetVar("y"))
	TypeSubject_AssertType(yVar.GetType()).ToStringIsEqualTo("string")
	AssertThat(yVar.IsTypeInferred()).IsTrue()
}

// Java: TypedScopeCreatorTest.testOutOfOrderJSDocForDestructuringParameter()
func (c *TypedScopeCreatorTest) TestOutOfOrderJSDocForDestructuringParameter() {
	// Even though regular JSDoc parameters can be out of order, there's currently no way for
	// putting arbitrary orders on destructuring parameters.
	c.TestWarning(
		"/**\n"+
			" * @param {!Iterable<number>} arr\n"+
			" * @param {string} x\n"+
			" */\n"+
			"function f(x, [y, z]) {}\n",
		FunctionTypeBuilder_INEXISTENT_PARAM)

	fVar := jutil.CheckNotNull(c.globalScope.GetVar("f"))
	// TODO(b/77597706): it would make more sense for this to be function(string, ?) instead
	TypeSubject_AssertType(fVar.GetType()).ToStringIsEqualTo("function(string, string): undefined")
	AssertThat(fVar.IsTypeInferred()).IsFalse()
}

// Java: TypedScopeCreatorTest.testObjectPatternParameterWithInlineJSDoc()
func (c *TypedScopeCreatorTest) TestObjectPatternParameterWithInlineJSDoc() {
	c.TestSame("function f({/** number */ x}) {}")

	// TODO(b/112651122): infer that f takes {x: number}
	fVar := jutil.CheckNotNull(c.globalScope.GetVar("f"))
	TypeSubject_AssertType(fVar.GetType()).ToStringIsEqualTo("function(?): undefined")
	AssertThat(fVar.IsTypeInferred()).IsFalse()

	xVar := jutil.CheckNotNull(c.lastFunctionScope.GetVar("x"))
	TypeSubject_AssertType(xVar.GetType()).ToStringIsEqualTo("number")
	AssertThat(xVar.IsTypeInferred()).IsFalse()
}

// Java: TypedScopeCreatorTest.testArrayPatternParameterWithInlineJSDoc()
func (c *TypedScopeCreatorTest) TestArrayPatternParameterWithInlineJSDoc() {
	c.TestSame("function f([/** number */ x]) {}")

	// TODO(b/112651122): either forbid this case or infer that f takes an !Iterable<number>
	fVar := jutil.CheckNotNull(c.globalScope.GetVar("f"))
	TypeSubject_AssertType(fVar.GetType()).ToStringIsEqualTo("function(?): undefined")
	AssertThat(fVar.IsTypeInferred()).IsFalse()

	xVar := jutil.CheckNotNull(c.lastFunctionScope.GetVar("x"))
	TypeSubject_AssertType(xVar.GetType()).ToStringIsEqualTo("number")
	AssertThat(xVar.IsTypeInferred()).IsFalse()
}

// Java: TypedScopeCreatorTest.testArrayPatternParametersWithDifferingInlineJSDoc()
func (c *TypedScopeCreatorTest) TestArrayPatternParametersWithDifferingInlineJSDoc() {
	c.TestSame("function f([/** number */ x, /** string */ y]) {}")

	// TODO(b/112651122): forbid this case, as there's not a good way to type the function without
	// having tuple types.
	fVar := jutil.CheckNotNull(c.globalScope.GetVar("f"))
	TypeSubject_AssertType(fVar.GetType()).ToStringIsEqualTo("function(?): undefined")
	AssertThat(fVar.IsTypeInferred()).IsFalse()

	xVar := jutil.CheckNotNull(c.lastFunctionScope.GetVar("x"))
	TypeSubject_AssertType(xVar.GetType()).ToStringIsEqualTo("number")
	AssertThat(xVar.IsTypeInferred()).IsFalse()

	yVar := jutil.CheckNotNull(c.lastFunctionScope.GetVar("y"))
	TypeSubject_AssertType(yVar.GetType()).ToStringIsEqualTo("string")
	AssertThat(yVar.IsTypeInferred()).IsFalse()
}

// Java: TypedScopeCreatorTest.testObjectPatternCanAliasUnionTypeProperty()
func (c *TypedScopeCreatorTest) TestObjectPatternCanAliasUnionTypeProperty() {
	// tests that we can get properties off the union type `{objA|objB}` even though it's not
	// represented as an ObjectType. (this used to crash in TypedScopeCreator)
	c.TestSame(
		"const ns = {};\n" +
			"\n" +
			"/** @typedef {{propA: number}} */\n" +
			"let objA;\n" +
			"/** @typedef {{propB: string}} */\n" +
			"let objB;\n" +
			"\n" +
			"/** @type {objA|objB} */\n" +
			"ns.ctor;\n" +
			"\n" +
			"const {propA, propB} = ns.ctor;\n")

	propAVar := jutil.CheckNotNull(c.globalScope.GetVar("propA"))
	TypeSubject_AssertType(propAVar.GetType()).IsEqualTo(c.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE))
	AssertThat(propAVar.IsTypeInferred()).IsFalse()

	propBVar := jutil.CheckNotNull(c.globalScope.GetVar("propB"))
	TypeSubject_AssertType(propBVar.GetType()).IsEqualTo(c.GetNativeType(rhino.JSTypeNative_STRING_TYPE))
	AssertThat(propBVar.IsTypeInferred()).IsFalse()
}

// Java: TypedScopeCreatorTest.testConstNestedObjectPatternInArrayPattern()
func (c *TypedScopeCreatorTest) TestConstNestedObjectPatternInArrayPattern() {
	// regression test for a TypedScopeCreator crash
	c.DisableTypeInfoValidation()
	c.TestSame(
		"const /** !Array<{b: string}> */ a = [{b: 'bbb'}];\n" +
			"const [{b}] = a;\n")

	b := jutil.CheckNotNull(c.globalScope.GetVar("b"))
	TypeSubject_AssertType(b.GetType()).IsString()
	AssertThat(b.IsTypeInferred()).IsTrue() // b is inferred but not treated as declared
}

// Java: TypedScopeCreatorTest.testConstNestedObjectPatternWithComputedPropertyIsUnknown()
func (c *TypedScopeCreatorTest) TestConstNestedObjectPatternWithComputedPropertyIsUnknown() {
	// regression test for a TypedScopeCreator crash
	c.DisableTypeInfoValidation() // fails on a['foo'].b = function() {};
	c.TestSame(
		"const a = {};\n" +
			"/** @const */ a['foo'] = {};\n" +
			"/** @const */ a['foo'].b = function() {};\n" +
			"\n" +
			"const {['foo']: {b}} = a;\n")

	b := jutil.CheckNotNull(c.globalScope.GetVar("b"))
	TypeSubject_AssertType(b.GetType()).IsUnknown()
	AssertThat(b.IsTypeInferred()).IsTrue()
}

// Java: TypedScopeCreatorTest.testStubProperty()
func (c *TypedScopeCreatorTest) TestStubProperty() {
	c.TestSame("function Foo() {}; Foo.bar;")
	foo := typedScopeCreatorTest_cast[rhino.ObjectType](c.globalScope.GetVar("Foo").GetType())
	AssertThat(foo.HasProperty("bar")).IsFalse()
	TypeSubject_AssertType(foo.GetPropertyType("bar")).IsEqualTo(c.registry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE))
}

// Java: TypedScopeCreatorTest.testConstructorProperty()
func (c *TypedScopeCreatorTest) TestConstructorProperty() {
	c.TestSame("var foo = {}; /** @constructor */ foo.Bar = function() {};")
	foo := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("foo", c.globalScope))
	AssertThat(foo.HasProperty("Bar")).IsTrue()
	AssertThat(foo.IsPropertyTypeInferred("Bar")).IsFalse()

	fooBar := foo.GetPropertyType("Bar")
	AssertThat(fooBar.String()).IsEqualTo("(typeof foo.Bar)")
}

// Java: TypedScopeCreatorTest.testPrototypePropertyMethodWithoutAnnotation()
func (c *TypedScopeCreatorTest) TestPrototypePropertyMethodWithoutAnnotation() {
	// TODO(b/154044898): re-enable type info validation
	c.DisableTypeInfoValidation()
	c.TestSame(
		"var Foo = function Foo() {};\n" +
			"var proto = Foo.prototype = {\n" +
			"   bar: function(a, b){}\n" +
			"};\n" +
			"proto.baz = function(c) {};\n" +
			"(function() { proto.baz = function() {}; })();\n")
	foo := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("Foo", c.globalScope))
	AssertThat(foo.HasProperty("prototype")).IsTrue()

	fooProto := typedScopeCreatorTest_cast[rhino.ObjectType](foo.GetPropertyType("prototype"))
	AssertThat(fooProto.HasProperty("bar")).IsTrue()
	AssertThat(fooProto.GetPropertyType("bar").String()).IsEqualTo("function(?, ?): undefined")

	AssertThat(fooProto.HasProperty("baz")).IsTrue()
	AssertThat(fooProto.GetPropertyType("baz").String()).IsEqualTo("function(?): undefined")
}

// Java: TypedScopeCreatorTest.testEnumProperty()
func (c *TypedScopeCreatorTest) TestEnumProperty() {
	c.TestSame("var foo = {}; /** @enum */ foo.Bar = {XXX: 'xxx'};")
	foo := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("foo", c.globalScope))
	AssertThat(foo.HasProperty("Bar")).IsTrue()
	AssertThat(foo.IsPropertyTypeInferred("Bar")).IsFalse()
	AssertThat(foo.IsPropertyTypeDeclared("Bar")).IsTrue()

	fooBar := foo.GetPropertyType("Bar")
	AssertThat(fooBar.String()).IsEqualTo("enum{foo.Bar}")
}

// Java: TypedScopeCreatorTest.testInferredProperty1()
func (c *TypedScopeCreatorTest) TestInferredProperty1() {
	c.TestSame("var foo = {}; foo.Bar = 3;")
	foo := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("foo", c.globalScope))
	AssertWithMessage(foo.String()).That(foo.HasProperty("Bar")).IsTrue()
	AssertThat(foo.GetPropertyType("Bar").String()).IsEqualTo("number")
	AssertThat(foo.IsPropertyTypeInferred("Bar")).IsTrue()
}

// Java: TypedScopeCreatorTest.testInferredProperty1a()
func (c *TypedScopeCreatorTest) TestInferredProperty1a() {
	c.TestSame("var foo = {}; /** @type {number} */ foo.Bar = 3;")
	foo := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("foo", c.globalScope))
	AssertWithMessage(foo.String()).That(foo.HasProperty("Bar")).IsTrue()
	AssertThat(foo.GetPropertyType("Bar").String()).IsEqualTo("number")
	AssertThat(foo.IsPropertyTypeInferred("Bar")).IsFalse()
}

// Java: TypedScopeCreatorTest.testInferredProperty2()
func (c *TypedScopeCreatorTest) TestInferredProperty2() {
	c.TestSame("var foo = { Bar: 3 };")
	foo := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("foo", c.globalScope))
	AssertWithMessage(foo.String()).That(foo.HasProperty("Bar")).IsTrue()
	AssertThat(foo.GetPropertyType("Bar").String()).IsEqualTo("number")
	AssertThat(foo.IsPropertyTypeInferred("Bar")).IsTrue()
}

// Java: TypedScopeCreatorTest.testInferredProperty2b()
func (c *TypedScopeCreatorTest) TestInferredProperty2b() {
	c.TestSame("var foo = { /** @type {number} */ Bar: 3 };")
	foo := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("foo", c.globalScope))
	AssertWithMessage(foo.String()).That(foo.HasProperty("Bar")).IsTrue()
	AssertThat(foo.GetPropertyType("Bar").String()).IsEqualTo("number")
	AssertThat(foo.IsPropertyTypeInferred("Bar")).IsFalse()
}

// Java: TypedScopeCreatorTest.testInferredProperty2c()
func (c *TypedScopeCreatorTest) TestInferredProperty2c() {
	c.TestSame("var foo = { /** @return {number} */ Bar: 3 };")
	foo := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("foo", c.globalScope))
	AssertWithMessage(foo.String()).That(foo.HasProperty("Bar")).IsTrue()
	AssertThat(foo.GetPropertyType("Bar").String()).IsEqualTo("function(): number")
	AssertThat(foo.IsPropertyTypeInferred("Bar")).IsFalse()
}

// Java: TypedScopeCreatorTest.testInferredProperty3()
func (c *TypedScopeCreatorTest) TestInferredProperty3() {
	c.TestSame("var foo = { /** @type {number} */ get Bar() { return 3 } };")
	foo := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("foo", c.globalScope))
	AssertWithMessage(foo.String()).That(foo.HasProperty("Bar")).IsTrue()
	AssertThat(foo.GetPropertyType("Bar").String()).IsEqualTo("?")
	AssertThat(foo.IsPropertyTypeInferred("Bar")).IsTrue()
}

// Java: TypedScopeCreatorTest.testInferredProperty4()
func (c *TypedScopeCreatorTest) TestInferredProperty4() {
	c.TestSame("var foo = { /** @type {number} */ set Bar(a) {} };")
	foo := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("foo", c.globalScope))
	AssertWithMessage(foo.String()).That(foo.HasProperty("Bar")).IsTrue()
	AssertThat(foo.GetPropertyType("Bar").String()).IsEqualTo("?")
	AssertThat(foo.IsPropertyTypeInferred("Bar")).IsTrue()
}

// Java: TypedScopeCreatorTest.testInferredProperty5()
func (c *TypedScopeCreatorTest) TestInferredProperty5() {
	c.TestSame("var foo = { /** @return {number} */ get Bar() { return 3 } };")
	foo := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("foo", c.globalScope))
	AssertWithMessage(foo.String()).That(foo.HasProperty("Bar")).IsTrue()
	AssertThat(foo.GetPropertyType("Bar").String()).IsEqualTo("number")
	AssertThat(foo.IsPropertyTypeInferred("Bar")).IsFalse()
}

// Java: TypedScopeCreatorTest.testInferredProperty6()
func (c *TypedScopeCreatorTest) TestInferredProperty6() {
	c.TestSame("var foo = { /** @param {number} a */ set Bar(a) {} };")
	foo := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("foo", c.globalScope))
	AssertWithMessage(foo.String()).That(foo.HasProperty("Bar")).IsTrue()
	AssertThat(foo.GetPropertyType("Bar").String()).IsEqualTo("number")
	AssertThat(foo.IsPropertyTypeInferred("Bar")).IsFalse()
}

// Java: TypedScopeCreatorTest.testPrototypeInit()
func (c *TypedScopeCreatorTest) TestPrototypeInit() {
	c.TestSame(
		"/** @constructor */ var Foo = function() {};\n" +
			"Foo.prototype = {bar: 1}; var foo = new Foo();\n")
	foo := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("foo", c.globalScope))
	AssertThat(foo.HasProperty("bar")).IsTrue()
	AssertThat(foo.GetPropertyType("bar").String()).IsEqualTo("number")
	AssertThat(foo.IsPropertyTypeInferred("bar")).IsTrue()
}

// Java: TypedScopeCreatorTest.testBogusPrototypeInit()
func (c *TypedScopeCreatorTest) TestBogusPrototypeInit() {
	// This used to cause a compiler crash.
	c.TestSame(
		"goog.F = {}; /** @const */ goog.F.prototype = {};\n" +
			"/** @constructor */ goog.F = function() {};\n")
}

// Java: TypedScopeCreatorTest.testInferredPrototypeProperty1()
func (c *TypedScopeCreatorTest) TestInferredPrototypeProperty1() {
	c.TestSame(
		"/** @constructor */ var Foo = function() {};\n" +
			"Foo.prototype.bar = 1; var x = new Foo();\n")

	x := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("x", c.globalScope))
	AssertThat(x.HasProperty("bar")).IsTrue()
	AssertThat(x.GetPropertyType("bar").String()).IsEqualTo("number")
	AssertThat(x.IsPropertyTypeInferred("bar")).IsTrue()
}

// Java: TypedScopeCreatorTest.testInferredPrototypeProperty2()
func (c *TypedScopeCreatorTest) TestInferredPrototypeProperty2() {
	c.TestSame(
		"/** @constructor */ var Foo = function() {};\n" +
			"Foo.prototype = {bar: 1}; var x = new Foo();\n")

	x := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("x", c.globalScope))
	AssertThat(x.HasProperty("bar")).IsTrue()
	AssertThat(x.GetPropertyType("bar").String()).IsEqualTo("number")
	AssertThat(x.IsPropertyTypeInferred("bar")).IsTrue()
}

// Java: TypedScopeCreatorTest.testEnum()
func (c *TypedScopeCreatorTest) TestEnum() {
	c.TestSame("/** @enum */ var Foo = {BAR: 1}; var f = Foo;")
	f := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("f", c.globalScope))
	AssertThat(f.HasProperty("BAR")).IsTrue()
	TypeSubject_AssertType(f.GetPropertyType("BAR")).ToStringIsEqualTo("Foo<number>")
	IsInstanceOf[*rhino.EnumType](AssertThat(f))
}

// Java: TypedScopeCreatorTest.testLetEnum()
func (c *TypedScopeCreatorTest) TestLetEnum() {
	c.TestSame("/** @enum */ let Foo = {BAR: 1}; let f = Foo;")
	f := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("f", c.globalScope))
	AssertThat(f.HasProperty("BAR")).IsTrue()
	TypeSubject_AssertType(f.GetPropertyType("BAR")).ToStringIsEqualTo("Foo<number>")
	IsInstanceOf[*rhino.EnumType](AssertThat(f))
}

// Java: TypedScopeCreatorTest.testConstEnum()
func (c *TypedScopeCreatorTest) TestConstEnum() {
	c.TestSame("/** @enum */ const Foo = {BAR: 1}; const f = Foo;")
	f := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("f", c.globalScope))
	AssertThat(f.HasProperty("BAR")).IsTrue()
	TypeSubject_AssertType(f.GetPropertyType("BAR")).ToStringIsEqualTo("Foo<number>")
	IsInstanceOf[*rhino.EnumType](AssertThat(f))
}

// Java: TypedScopeCreatorTest.testEnumElement()
func (c *TypedScopeCreatorTest) TestEnumElement() {
	c.TestSame("/** @enum */ var Foo = {BAR: 1}; var f = Foo;")
	bar := c.globalScope.GetVar("Foo.BAR")
	TypedVarSubject_AssertThat(bar).IsNotNull()
	TypeSubject_AssertType(bar.GetType()).ToStringIsEqualTo("Foo<number>")
}

// Java: TypedScopeCreatorTest.testLetEnumElement()
func (c *TypedScopeCreatorTest) TestLetEnumElement() {
	c.TestSame("/** @enum */ let Foo = {BAR: 1}; let f = Foo;")
	bar := c.globalScope.GetVar("Foo.BAR")
	TypedVarSubject_AssertThat(bar).IsNotNull()
	TypeSubject_AssertType(bar.GetType()).ToStringIsEqualTo("Foo<number>")
}

// Java: TypedScopeCreatorTest.testConstEnumElement()
func (c *TypedScopeCreatorTest) TestConstEnumElement() {
	c.TestSame("/** @enum */ const Foo = {BAR: 1}; const f = Foo;")
	bar := c.globalScope.GetVar("Foo.BAR")
	TypedVarSubject_AssertThat(bar).IsNotNull()
	TypeSubject_AssertType(bar.GetType()).ToStringIsEqualTo("Foo<number>")
}

// Java: TypedScopeCreatorTest.testNamespacedEnum()
func (c *TypedScopeCreatorTest) TestNamespacedEnum() {
	c.TestSame(
		"goog.ui = {};\n" +
			"/** @constructor */goog.ui.Zippy = function() {};\n" +
			"/** @enum{string} */goog.ui.Zippy.EventType = { TOGGLE: 'toggle' };\n" +
			"var x = goog.ui.Zippy.EventType;\n" +
			"var y = goog.ui.Zippy.EventType.TOGGLE;\n")

	x := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("x", c.globalScope))
	AssertThat(x.IsEnumType()).IsTrue()
	AssertThat(x.HasProperty("TOGGLE")).IsTrue()
	AssertThat(x.GetReferenceName()).IsEqualTo("enum{goog.ui.Zippy.EventType}")

	y := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("y", c.globalScope))
	AssertThat(y.IsSubtypeOf(c.GetNativeType(rhino.JSTypeNative_STRING_TYPE))).IsTrue()
	AssertThat(y.IsEnumElementType()).IsTrue()
	AssertThat(y.GetReferenceName()).IsEqualTo("goog.ui.Zippy.EventType")
}

// Java: TypedScopeCreatorTest.testGlobalTypedefs1()
func (c *TypedScopeCreatorTest) TestGlobalTypedefs1() {
	c.TestSame(
		" // preserve newlines\n" +
			"/** @typedef {number} */\n" +
			"var VarTypedef;\n" +
			"{\n" +
			"  /** @typedef {number} */\n" +
			"  var VarTypedefInBlock; // still global despite enclosing block\n" +
			"}\n" +
			"/** @typedef {number} */\n" +
			"let LetTypedef;\n" +
			"// TODO(bradfordcsmith): We should probably disallow @typedef on const.\n" +
			"/** @typedef {number} */\n" +
			"const ConstTypedef = undefined;\n")
	TypeSubject_AssertType(c.registry.GetGlobalType("VarTypedef")).IsNumber()
	TypeSubject_AssertType(c.registry.GetGlobalType("VarTypedefInBlock")).IsNumber()
	TypeSubject_AssertType(c.registry.GetGlobalType("LetTypedef")).IsNumber()
	TypeSubject_AssertType(c.registry.GetGlobalType("ConstTypedef")).IsNumber()
}

// Java: TypedScopeCreatorTest.testGlobalTypedefs2()
func (c *TypedScopeCreatorTest) TestGlobalTypedefs2() {
	c.TestSame(
		" // preserve newlines\n" +
			"{\n" +
			"  /** @typedef {number} */\n" +
			"  var VarTypedefInBlock; // still global despite enclosing block\n" +
			"}\n")
}

// Java: TypedScopeCreatorTest.testLocalTypedefs1()
func (c *TypedScopeCreatorTest) TestLocalTypedefs1() {
	c.TestSame(
		" // preserve newlines\n" +
			"function f() {\n" +
			"  /** @typedef {number} */\n" +
			"  var VarTypedefInFunc;\n" +
			"  /** @typedef {number} */\n" +
			"  let LetTypedefInFunc;\n" +
			"  /** @typedef {number} */\n" +
			"// TODO(bradfordcsmith): We should probably disallow @typedef on const.\n" +
			"  const ConstTypedefInFunc = undefined;\n" +
			"}\n")
	fnRoot := c.lastLocalScope
	TypeSubject_AssertType(c.registry.GetType(fnRoot, "VarTypedefInFunc")).IsNumber()
	TypeSubject_AssertType(c.registry.GetType(fnRoot, "LetTypedefInFunc")).IsNumber()
	TypeSubject_AssertType(c.registry.GetType(fnRoot, "ConstTypedefInFunc")).IsNumber()
}

// Java: TypedScopeCreatorTest.testLocalTypedefs2()
func (c *TypedScopeCreatorTest) TestLocalTypedefs2() {
	c.TestSame(
		" // preserve newlines\n" +
			"{\n" +
			"  /** @typedef {number} */\n" +
			"  let LetTypedefInBlock;\n" +
			"  /** @typedef {number} */\n" +
			"  const ConstTypedefInBlock = undefined;\n" +
			"}\n")
	fnRoot := c.lastLocalScope

	TypeSubject_AssertType(c.registry.GetType(fnRoot, "LetTypedefInBlock")).IsNumber()
	TypeSubject_AssertType(c.registry.GetType(fnRoot, "ConstTypedefInBlock")).IsNumber()
}

// Java: TypedScopeCreatorTest.testEnumAlias()
func (c *TypedScopeCreatorTest) TestEnumAlias() {
	c.TestSame(
		"/** @enum */ var Foo = {BAR: 1};\n" +
			"/** @enum */ var FooAlias = Foo; var f = FooAlias;\n")

	AssertThat(c.registry.GetType(nil, "FooAlias").String()).IsEqualTo("Foo<number>")
	TypeSubject_AssertType(c.registry.GetType(nil, "Foo")).IsEqualTo(c.registry.GetType(nil, "FooAlias"))

	f := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("f", c.globalScope))
	AssertThat(f.HasProperty("BAR")).IsTrue()
	AssertThat(f.GetPropertyType("BAR").String()).IsEqualTo("Foo<number>")
	IsInstanceOf[*rhino.EnumType](AssertThat(f))
}

// Java: TypedScopeCreatorTest.testNamespacesEnumAlias()
func (c *TypedScopeCreatorTest) TestNamespacesEnumAlias() {
	c.TestSame(
		"/** @enum */ goog.Foo = {BAR: 1};\n" +
			"/** @enum */ goog.FooAlias = goog.Foo;\n")

	AssertThat(c.registry.GetType(nil, "goog.FooAlias").String()).IsEqualTo("goog.Foo<number>")
	TypeSubject_AssertType(c.registry.GetType(nil, "goog.FooAlias")).
		IsEqualTo(c.registry.GetType(nil, "goog.Foo"))
}

// Java: TypedScopeCreatorTest.testCollectedFunctionStub()
func (c *TypedScopeCreatorTest) TestCollectedFunctionStub() {
	c.TestSame(
		"/** @constructor */ function f() {\n" +
			"  /** @return {number} */ this.foo;\n" +
			"}\n" +
			"var x = new f();\n")
	x := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("x", c.globalScope))
	AssertThat(x.String()).IsEqualTo("f")
	AssertThat(x.HasProperty("foo")).IsTrue()
	AssertThat(x.GetPropertyType("foo").String()).IsEqualTo("function(this:f): number")
	AssertThat(x.IsPropertyTypeInferred("foo")).IsFalse()
}

// Java: TypedScopeCreatorTest.testCollectedFunctionStubLocal()
func (c *TypedScopeCreatorTest) TestCollectedFunctionStubLocal() {
	c.TestSame(
		"(function() {\n" +
			"/** @constructor */ function f() {\n" +
			"  /** @return {number} */ this.foo;\n" +
			"}\n" +
			"var x = new f();\n" +
			"});\n")
	x := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("x", c.lastLocalScope))
	AssertThat(x.String()).IsEqualTo("f")
	AssertThat(x.HasProperty("foo")).IsTrue()
	AssertThat(x.GetPropertyType("foo").String()).IsEqualTo("function(this:f): number")
	AssertThat(x.IsPropertyTypeInferred("foo")).IsFalse()
}

// Java: TypedScopeCreatorTest.testNamespacedFunctionStub()
func (c *TypedScopeCreatorTest) TestNamespacedFunctionStub() {
	c.TestSame("/** @param {number} x */ goog.foo;")

	goog := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("goog", c.globalScope))
	AssertThat(goog.HasProperty("foo")).IsTrue()
	AssertThat(goog.GetPropertyType("foo").String()).IsEqualTo("function(number): ?")
	AssertThat(goog.IsPropertyTypeDeclared("foo")).IsTrue()

	TypeSubject_AssertType(goog.GetPropertyType("foo")).IsEqualTo(c.globalScope.GetVar("goog.foo").GetType())
}

// Java: TypedScopeCreatorTest.testNamespacedFunctionStubLocal()
func (c *TypedScopeCreatorTest) TestNamespacedFunctionStubLocal() {
	c.TestSame(
		"(function() {\n" +
			"var goog = {};\n" +
			"/** @param {number} x */ goog.foo;\n" +
			"});\n")

	goog := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("goog", c.lastLocalScope))
	AssertThat(goog.HasProperty("foo")).IsTrue()
	AssertThat(goog.GetPropertyType("foo").String()).IsEqualTo("function(number): ?")
	AssertThat(goog.IsPropertyTypeDeclared("foo")).IsTrue()

	TypeSubject_AssertType(goog.GetPropertyType("foo")).IsEqualTo(c.lastLocalScope.GetVar("goog.foo").GetType())
}

// Java: TypedScopeCreatorTest.testCollectedCtorProperty1()
func (c *TypedScopeCreatorTest) TestCollectedCtorProperty1() {
	c.TestSame(
		"/** @constructor */ function f() {\n" +
			"  /** @type {number} */ this.foo = 3;\n" +
			"}\n" +
			"var x = new f();\n")
	x := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("x", c.globalScope))
	AssertThat(x.String()).IsEqualTo("f")
	AssertThat(x.HasProperty("foo")).IsTrue()
	AssertThat(x.GetPropertyType("foo").String()).IsEqualTo("number")
	AssertThat(x.IsPropertyTypeInferred("foo")).IsFalse()
	AssertThat(x.IsPropertyTypeDeclared("foo")).IsTrue()
}

// Java: TypedScopeCreatorTest.testCollectedCtorProperty2()
func (c *TypedScopeCreatorTest) TestCollectedCtorProperty2() {
	c.TestSame(
		"/** @constructor */ function f() {\n" +
			"  /** @const {number} */ this.foo = 3;\n" +
			"}\n" +
			"var x = new f();\n")
	x := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("x", c.globalScope))
	AssertThat(x.String()).IsEqualTo("f")
	AssertThat(x.HasProperty("foo")).IsTrue()
	AssertThat(x.GetPropertyType("foo").String()).IsEqualTo("number")
	AssertThat(x.IsPropertyTypeInferred("foo")).IsFalse()
	AssertThat(x.IsPropertyTypeDeclared("foo")).IsTrue()
}

// Java: TypedScopeCreatorTest.testCollectedCtorProperty3()
func (c *TypedScopeCreatorTest) TestCollectedCtorProperty3() {
	c.TestSame(
		"/** @constructor */ function f() {\n" +
			"  /** @const */ this.foo = 3;\n" +
			"}\n" +
			"var x = new f();\n")
	x := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("x", c.globalScope))
	AssertThat(x.String()).IsEqualTo("f")
	AssertThat(x.HasProperty("foo")).IsTrue()
	AssertThat(x.GetPropertyType("foo").String()).IsEqualTo("number")
	AssertThat(x.IsPropertyTypeInferred("foo")).IsFalse()
	AssertThat(x.IsPropertyTypeDeclared("foo")).IsTrue()
}

// Java: TypedScopeCreatorTest.testCollectedCtorProperty5()
func (c *TypedScopeCreatorTest) TestCollectedCtorProperty5() {
	c.TestSame(
		"/** @constructor */ function f() {\n" +
			"  /** @const */ this.foo = 'abc' + 'def';\n" +
			"}\n" +
			"var x = new f();\n")
	x := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("x", c.globalScope))
	AssertThat(x.String()).IsEqualTo("f")
	AssertThat(x.HasProperty("foo")).IsTrue()
	AssertThat(x.GetPropertyType("foo").String()).IsEqualTo("string")
	AssertThat(x.IsPropertyTypeInferred("foo")).IsFalse()
	AssertThat(x.IsPropertyTypeDeclared("foo")).IsTrue()
}

// Java: TypedScopeCreatorTest.testCollectedCtorBigIntProperty()
func (c *TypedScopeCreatorTest) TestCollectedCtorBigIntProperty() {
	c.TestSame(
		"/** @constructor */ function f() {\n" +
			"  /** @const */ this.foo = 1n + 2n;\n" +
			"}\n" +
			"var x = new f();\n")
	x := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("x", c.globalScope))
	AssertThat(x.String()).IsEqualTo("f")
	AssertThat(x.HasProperty("foo")).IsTrue()
	AssertThat(x.GetPropertyType("foo").String()).IsEqualTo("bigint")
	AssertThat(x.IsPropertyTypeInferred("foo")).IsFalse()
	AssertThat(x.IsPropertyTypeDeclared("foo")).IsTrue()
}

// Java: TypedScopeCreatorTest.testCollectedCtorProperty9()
func (c *TypedScopeCreatorTest) TestCollectedCtorProperty9() {
	c.TestSame(
		"/** @constructor */ function f() {}\n" +
			"f.prototype.init_f = function() {\n" +
			"  /** @const */ this.FOO = 'abc';\n" +
			"};\n" +
			"var x = new f();\n")
	x := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("x", c.globalScope))
	AssertThat(x.String()).IsEqualTo("f")
	AssertThat(x.HasProperty("FOO")).IsTrue()
	AssertThat(x.GetPropertyType("FOO").String()).IsEqualTo("string")
	AssertThat(x.IsPropertyTypeInferred("FOO")).IsFalse()
	AssertThat(x.IsPropertyTypeDeclared("FOO")).IsTrue()
}

// Java: TypedScopeCreatorTest.testCollectedCtorProperty10()
func (c *TypedScopeCreatorTest) TestCollectedCtorProperty10() {
	c.TestSame(
		"/** @constructor */ function f() {}\n" +
			"f.prototype.init_f = function() {\n" +
			"  /** @const */ this.foo = new String();\n" +
			"};\n" +
			"var x = new f();\n")
	x := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("x", c.globalScope))
	AssertThat(x.String()).IsEqualTo("f")
	AssertThat(x.HasProperty("foo")).IsTrue()
	AssertThat(x.GetPropertyType("foo").String()).IsEqualTo("String")
	AssertThat(x.IsPropertyTypeInferred("foo")).IsFalse()
	AssertThat(x.IsPropertyTypeDeclared("foo")).IsTrue()
}

// Java: TypedScopeCreatorTest.testCollectedCtorProperty11()
func (c *TypedScopeCreatorTest) TestCollectedCtorProperty11() {
	c.TestSame(
		"/** @constructor */ function f() {}\n" +
			"f.prototype.init_f = function() {\n" +
			"  /** @const */ this.foo = [];\n" +
			"};\n" +
			"var x = new f();\n")
	x := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("x", c.globalScope))
	AssertThat(x.String()).IsEqualTo("f")
	AssertThat(x.HasProperty("foo")).IsTrue()
	AssertThat(x.GetPropertyType("foo").String()).IsEqualTo("Array")
	AssertThat(x.IsPropertyTypeInferred("foo")).IsFalse()
	AssertThat(x.IsPropertyTypeDeclared("foo")).IsTrue()
}

// Java: TypedScopeCreatorTest.testCollectedCtorProperty12()
func (c *TypedScopeCreatorTest) TestCollectedCtorProperty12() {
	c.TestSame(
		c.Externs("/** @type {?} */ let unknown;"),
		c.Srcs(
			"/** @constructor */ function f() {}\n"+
				"f.prototype.init_f = function() {\n"+
				"  /** @const */ this.foo = !!unknown;\n"+
				"};\n"+
				"var x = new f();\n"))
	x := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("x", c.globalScope))
	AssertThat(x.String()).IsEqualTo("f")
	AssertThat(x.HasProperty("foo")).IsTrue()
	AssertThat(x.GetPropertyType("foo").String()).IsEqualTo("boolean")
	AssertThat(x.IsPropertyTypeInferred("foo")).IsFalse()
	AssertThat(x.IsPropertyTypeDeclared("foo")).IsTrue()
}

// Java: TypedScopeCreatorTest.testCollectedCtorProperty13()
func (c *TypedScopeCreatorTest) TestCollectedCtorProperty13() {
	c.TestSame(
		c.Externs("/** @type {?} */ let unknown;"),
		c.Srcs(
			"/** @constructor */ function f() {}\n"+
				"f.prototype.init_f = function() {\n"+
				"  /** @const */ this.foo = +unknown;\n"+
				"};\n"+
				"var x = new f();\n"))
	x := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("x", c.globalScope))
	AssertThat(x.String()).IsEqualTo("f")
	AssertThat(x.HasProperty("foo")).IsTrue()
	AssertThat(x.GetPropertyType("foo").String()).IsEqualTo("number")
	AssertThat(x.IsPropertyTypeInferred("foo")).IsFalse()
	AssertThat(x.IsPropertyTypeDeclared("foo")).IsTrue()
}

// Java: TypedScopeCreatorTest.testCollectedCtorProperty14()
func (c *TypedScopeCreatorTest) TestCollectedCtorProperty14() {
	c.TestSame(
		c.Externs("/** @type {?} */ let unknown;"),
		c.Srcs(
			"/** @constructor */ function f() {}\n"+
				"f.prototype.init_f = function() {\n"+
				"  /** @const */ this.foo = unknown + '';\n"+
				"};\n"+
				"var x = new f();\n"))
	x := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("x", c.globalScope))
	AssertThat(x.String()).IsEqualTo("f")
	AssertThat(x.HasProperty("foo")).IsTrue()
	AssertThat(x.GetPropertyType("foo").String()).IsEqualTo("string")
	AssertThat(x.IsPropertyTypeInferred("foo")).IsFalse()
	AssertThat(x.IsPropertyTypeDeclared("foo")).IsTrue()
}

// Java: TypedScopeCreatorTest.testCollectedCtorProperty15()
func (c *TypedScopeCreatorTest) TestCollectedCtorProperty15() {
	c.TestSame(
		"/**\n" +
			" * @constructor\n" +
			" * @param {string} a\n" +
			" */\n" +
			" function f(a) {\n" +
			"  /** @const */ this.foo = a;\n" +
			"};\n" +
			"var x = new f();\n")
	x := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("x", c.globalScope))
	AssertThat(x.String()).IsEqualTo("f")
	AssertThat(x.HasProperty("foo")).IsTrue()
	AssertThat(x.GetPropertyType("foo").String()).IsEqualTo("string")
	AssertThat(x.IsPropertyTypeInferred("foo")).IsFalse()
	AssertThat(x.IsPropertyTypeDeclared("foo")).IsTrue()
}

// Java: TypedScopeCreatorTest.testPropertyOnUnknownSuperClass1()
func (c *TypedScopeCreatorTest) TestPropertyOnUnknownSuperClass1() {
	c.TestWarning(
		"var goog = this.foo();\n"+
			"/** @constructor \n"+
			" * @extends {goog.Unknown} */\n"+
			"function Foo() {}\n"+
			"Foo.prototype.bar = 1;\n"+
			"var x = new Foo();\n",
		RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR)
	x := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("x", c.globalScope))
	AssertThat(x.String()).IsEqualTo("Foo")
	AssertThat(x.GetImplicitPrototype().HasOwnProperty("bar")).IsTrue()
	AssertThat(x.GetPropertyType("bar").String()).IsEqualTo("?")
	AssertThat(x.IsPropertyTypeInferred("bar")).IsTrue()
}

// Java: TypedScopeCreatorTest.testPropertyOnUnknownSuperClass2()
func (c *TypedScopeCreatorTest) TestPropertyOnUnknownSuperClass2() {
	c.TestWarning(
		"var goog = this.foo();\n"+
			"/** @constructor \n"+
			" * @extends {goog.Unknown} */\n"+
			"function Foo() {}\n"+
			"Foo.prototype = {bar: 1};\n"+
			"var x = new Foo();\n",
		RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR)
	x := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("x", c.globalScope))
	AssertThat(x.String()).IsEqualTo("Foo")
	AssertThat(x.GetImplicitPrototype().String()).IsEqualTo("Foo.prototype")
	AssertThat(x.GetImplicitPrototype().HasOwnProperty("bar")).IsTrue()
	AssertThat(x.GetPropertyType("bar").String()).IsEqualTo("?")
	AssertThat(x.IsPropertyTypeInferred("bar")).IsTrue()
}

// Java: TypedScopeCreatorTest.testSubBeforeSuper1()
func (c *TypedScopeCreatorTest) TestSubBeforeSuper1() {
	c.TestSame(
		"/** @interface\n" +
			" * @extends {MidI} */\n" +
			"function LowI() {}\n" +
			"/** @interface\n" +
			" * @extends {HighI} */\n" +
			"function MidI() {}\n" +
			"/** @interface */\n" +
			"function HighI() {}\n")
}

// Java: TypedScopeCreatorTest.testSubBeforeSuper2()
func (c *TypedScopeCreatorTest) TestSubBeforeSuper2() {
	c.TestSame(
		"/** @constructor\n" +
			" * @extends {MidI} */\n" +
			"function LowI() {}\n" +
			"/** @constructor\n" +
			" * @extends {HighI} */\n" +
			"function MidI() {}\n" +
			"/** @constructor */\n" +
			"function HighI() {}\n")
}

// Java: TypedScopeCreatorTest.testMethodBeforeFunction1()
func (c *TypedScopeCreatorTest) TestMethodBeforeFunction1() {
	// Adding a method to a hoisted constructor before its definition doesn't typecheck correctly.
	// Technically this code is valid JS but the TypedScopeCreator doesn't model function hoisting.
	// This is because hoisting causes a lot of new NamedTypes to appear, and NamedTypes are
	// generally buggy.
	c.TestSame(
		"var y = Window.prototype;\n" +
			"Window.prototype.alert = function(message) {};\n" +
			"/** @constructor */ function Window() {}\n" +
			"var window = new Window(); \n" +
			"var x = window;\n")
	x := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("x", c.globalScope))
	AssertThat(x.String()).IsEqualTo("Window")
	AssertThat(x.GetImplicitPrototype().HasOwnProperty("alert")).IsTrue()
	AssertThat(x.GetPropertyType("alert").String()).IsEqualTo("function(?): undefined")
	AssertThat(x.IsPropertyTypeDeclared("alert")).IsTrue()

	y := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("y", c.globalScope))
	AssertThat(y.GetPropertyType("alert").String()).IsEqualTo("function(?): undefined")
}

// Java: TypedScopeCreatorTest.testMethodBeforeFunction2()
func (c *TypedScopeCreatorTest) TestMethodBeforeFunction2() {
	c.TestSame(
		"var y = Window.prototype;\n" +
			"Window.prototype = {alert: function(message) {}};\n" +
			"/** @constructor */ function Window() {}\n" +
			"var window = new Window(); \n" +
			"var x = window;\n")
	x := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("x", c.globalScope))
	AssertThat(x.String()).IsEqualTo("Window")
	AssertThat(x.GetImplicitPrototype().HasOwnProperty("alert")).IsTrue()
	AssertThat(x.GetPropertyType("alert").String()).IsEqualTo("function(?): undefined")
	AssertThat(x.IsPropertyTypeDeclared("alert")).IsFalse()

	y := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("y", c.globalScope))
	TypeSubject_AssertType(y).WithTypeOfProp("alert").IsUnknown()
}

// Java: TypedScopeCreatorTest.testAddMethodsPrototypeTwoWays()
func (c *TypedScopeCreatorTest) TestAddMethodsPrototypeTwoWays() {
	c.TestSame(
		"/** @constructor */function A() {}\n" +
			"A.prototype = {m1: 5, m2: true};\n" +
			"A.prototype.m3 = 'third property!';\n" +
			"var x = new A();\n")

	instanceType := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("x", c.globalScope))
	AssertThat(instanceType.GetPropertiesCount()).
		IsEqualTo(c.GetNativeObjectType(rhino.JSTypeNative_OBJECT_TYPE).GetPropertiesCount() + 3)
	TypeSubject_AssertType(instanceType.GetPropertyType("m1")).IsEqualTo(c.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE))
	TypeSubject_AssertType(instanceType.GetPropertyType("m2")).IsEqualTo(c.GetNativeType(rhino.JSTypeNative_BOOLEAN_TYPE))
	TypeSubject_AssertType(instanceType.GetPropertyType("m3")).IsEqualTo(c.GetNativeType(rhino.JSTypeNative_STRING_TYPE))

	// Verify the prototype chain.
	// This is a special case where we want the anonymous object to
	// become a prototype.
	AssertThat(instanceType.HasOwnProperty("m1")).IsFalse()
	AssertThat(instanceType.HasOwnProperty("m2")).IsFalse()
	AssertThat(instanceType.HasOwnProperty("m3")).IsFalse()

	proto1 := instanceType.GetImplicitPrototype()
	AssertThat(proto1.HasOwnProperty("m1")).IsTrue()
	AssertThat(proto1.HasOwnProperty("m2")).IsTrue()
	AssertThat(proto1.HasOwnProperty("m3")).IsTrue()

	proto2 := proto1.GetImplicitPrototype()
	AssertThat(proto2.HasProperty("m1")).IsFalse()
	AssertThat(proto2.HasProperty("m2")).IsFalse()
	AssertThat(proto2.HasProperty("m3")).IsFalse()

	// Note: it's not necessary for correctness reasons to declare m1 and m2,  but it's trickier to
	// avoid it without extra string maniuplation
	ScopeSubject_AssertScope(c.globalScope).Declares("A.prototype.m1")
	ScopeSubject_AssertScope(c.globalScope).Declares("A.prototype.m2")
	ScopeSubject_AssertScope(c.globalScope).DoesNotDeclare("A.prototype.m3")
}

// Java: TypedScopeCreatorTest.testInferredVar()
func (c *TypedScopeCreatorTest) TestInferredVar() {
	c.TestSame("var x = 3; x = 'x'; x = true;")

	x := c.globalScope.GetVar("x")
	TypeSubject_AssertType(x.GetType()).ToStringIsEqualTo("(boolean|number|string)")
	AssertThat(x.IsTypeInferred()).IsTrue()
}

// Java: TypedScopeCreatorTest.testInferredLet()
func (c *TypedScopeCreatorTest) TestInferredLet() {
	c.TestSame("let x = 3; x = 'x'; x = true;")

	x := c.globalScope.GetVar("x")
	TypeSubject_AssertType(x.GetType()).ToStringIsEqualTo("(boolean|number|string)")
	AssertThat(x.IsTypeInferred()).IsTrue()
}

// Java: TypedScopeCreatorTest.testInferredConst()
func (c *TypedScopeCreatorTest) TestInferredConst() {
	c.TestSame("const x = 3;")

	x := c.globalScope.GetVar("x")
	TypeSubject_AssertType(x.GetType()).IsNumber()
	AssertThat(x.IsConst()).IsTrue()
	// Although we did infer the type, we'll consider it effectively declared because the variable
	// was declared to be constant. This is consistent with the way we handle the @const annotation
	// on var declarations.
	AssertThat(x.IsTypeInferred()).IsFalse()
}

// Java: TypedScopeCreatorTest.testInferredAnnotatedConst()
func (c *TypedScopeCreatorTest) TestInferredAnnotatedConst() {
	c.TestSame("/** @const */ var x = 3;")

	x := c.globalScope.GetVar("x")
	TypeSubject_AssertType(x.GetType()).IsNumber()
	AssertThat(x.IsConst()).IsFalse()
	AssertThat(x.IsTypeInferred()).IsFalse()
}

// Java: TypedScopeCreatorTest.testDeclaredVar()
func (c *TypedScopeCreatorTest) TestDeclaredVar() {
	c.TestSame("/** @type {?number} */ var x = 3; var y = x;")

	ScopeSubject_AssertScope(c.globalScope).Declares("x").Directly()
	x := c.globalScope.GetVar("x")
	TypeSubject_AssertType(x.GetType()).ToStringIsEqualTo("(null|number)")
	AssertThat(x.IsTypeInferred()).IsFalse()

	y := c.FindNameType("y", c.globalScope)
	AssertThat(y.String()).IsEqualTo("(null|number)")
}

// Java: TypedScopeCreatorTest.testDeclaredLet()
func (c *TypedScopeCreatorTest) TestDeclaredLet() {
	c.TestSame("/** @type {?number} */ let x = 3; let y = x;")

	ScopeSubject_AssertScope(c.globalScope).Declares("x").Directly()
	x := c.globalScope.GetVar("x")
	TypeSubject_AssertType(x.GetType()).ToStringIsEqualTo("(null|number)")
	AssertThat(x.IsTypeInferred()).IsFalse()

	y := c.FindNameType("y", c.globalScope)
	AssertThat(y.String()).IsEqualTo("(null|number)")
}

// Java: TypedScopeCreatorTest.testDeclaredConst()
func (c *TypedScopeCreatorTest) TestDeclaredConst() {
	c.TestSame("/** @type {?number} */ const x = 3; const y = x;")

	ScopeSubject_AssertScope(c.globalScope).Declares("x").Directly()
	x := c.globalScope.GetVar("x")
	TypeSubject_AssertType(x.GetType()).ToStringIsEqualTo("(null|number)")
	AssertThat(x.IsTypeInferred()).IsFalse()
	AssertThat(x.IsConst()).IsTrue()

	y := c.FindNameType("y", c.globalScope)
	AssertThat(y.String()).IsEqualTo("(null|number)")
}

// Java: TypedScopeCreatorTest.testStructuralInterfaceMatchingOnInterface1()
func (c *TypedScopeCreatorTest) TestStructuralInterfaceMatchingOnInterface1() {
	c.TestSame(
		"/** @record */ var I = function() {};\n" +
			"/** @type {number} */ I.prototype.bar;\n" +
			"I.prototype.baz = function(){};\n")

	i := c.globalScope.GetVar("I")
	AssertThat(i.GetType().String()).IsEqualTo("(typeof I)")
	AssertThat(i.GetType().IsInterface()).IsTrue()
	AssertThat(i.GetType().IsFunctionType()).IsTrue()
	AssertThat(i.GetType().ToMaybeFunctionType().IsStructuralInterface()).IsTrue()
}

// Java: TypedScopeCreatorTest.testStructuralInterfaceMatchingOnInterface2()
func (c *TypedScopeCreatorTest) TestStructuralInterfaceMatchingOnInterface2() {
	c.TestSame(
		"/** @interface */ var I = function() {};\n" +
			"/** @type {number} */ I.prototype.bar;\n" +
			"I.prototype.baz = function(){};\n")

	i := c.globalScope.GetVar("I")
	AssertThat(i.GetType().String()).IsEqualTo("(typeof I)")
	AssertThat(i.GetType().IsInterface()).IsTrue()
	AssertThat(i.GetType().IsFunctionType()).IsTrue()
	AssertThat(i.GetType().ToMaybeFunctionType().IsStructuralInterface()).IsFalse()
}

// Java: TypedScopeCreatorTest.testStructuralInterfaceMatchingOnInterface3()
func (c *TypedScopeCreatorTest) TestStructuralInterfaceMatchingOnInterface3() {
	c.TestSame(
		"/** @interface */ var I = function() {};\n" +
			"/** @type {number} */ I.prototype.bar;\n" +
			"/** @record */ I.prototype.baz = function() {};\n")

	iPrototype := c.globalScope.GetVar("I.prototype")
	baz := iPrototype.GetType().FindPropertyType("baz")
	AssertThat(baz.IsInterface()).IsTrue()
	AssertThat(baz.IsFunctionType()).IsTrue()
	AssertThat(baz.ToMaybeFunctionType().IsStructuralInterface()).IsTrue()
}

// Java: TypedScopeCreatorTest.testStructuralInterfaceMatchingOnInterface4()
func (c *TypedScopeCreatorTest) TestStructuralInterfaceMatchingOnInterface4() {
	c.TestSame(
		"/** @interface */ var I = function() {};\n" +
			"/** @type {number} */ I.prototype.bar;\n" +
			"/** @interface */ I.prototype.baz = function() {};\n")

	iPrototype := c.globalScope.GetVar("I.prototype")
	baz := iPrototype.GetType().FindPropertyType("baz")
	AssertThat(baz.IsInterface()).IsTrue()
	AssertThat(baz.IsFunctionType()).IsTrue()
	AssertThat(baz.ToMaybeFunctionType().IsStructuralInterface()).IsFalse()
}

// Java: TypedScopeCreatorTest.testStructuralInterfaceMatchingOnInterface5()
func (c *TypedScopeCreatorTest) TestStructuralInterfaceMatchingOnInterface5() {
	c.TestSame(
		"/** @constructor */ var C = function() {};\n" +
			"/** @type {number} */ C.prototype.bar;\n" +
			"/** @record */ C.prototype.baz = function() {};\n" +
			"var c = new C(); var cbaz = c.baz;\n")

	cBaz := c.globalScope.GetVar("cbaz")
	AssertThat(cBaz.GetType().IsFunctionType()).IsTrue()
	AssertThat(cBaz.GetType().ToMaybeFunctionType().IsStructuralInterface()).IsTrue()
}

// Java: TypedScopeCreatorTest.testStructuralInterfaceMatchingOnInterface6()
func (c *TypedScopeCreatorTest) TestStructuralInterfaceMatchingOnInterface6() {
	c.TestSame(
		"/** @constructor */ var C = function() {};\n" +
			"/** @type {number} */ C.prototype.bar;\n" +
			"/** @interface */ C.prototype.baz = function() {};\n" +
			"var c = new C(); var cbaz = c.baz;\n")

	cBaz := c.globalScope.GetVar("cbaz")
	AssertThat(cBaz.GetType().IsFunctionType()).IsTrue()
	AssertThat(cBaz.GetType().ToMaybeFunctionType().IsStructuralInterface()).IsFalse()
}

// Java: TypedScopeCreatorTest.testStructuralInterfaceMatchingOnInterface_symbolProps()
func (c *TypedScopeCreatorTest) TestStructuralInterfaceMatchingOnInterface_symbolProps() {
	c.TestSame(
		c.Externs(
			NewTestExternsBuilder().AddIterable().Build(),
			"/** @const {symbol} */\n"+
				"Symbol.number;\n"+
				"/** @const {symbol} */\n"+
				"Symbol.string;\n"),
		c.Srcs(
			"/** @record */\n"+
				"class I {\n"+
				"  constructor() {\n"+
				"    /** @type {string} */\n"+
				"    this[Symbol.string];\n"+
				"  }\n"+
				"  /** @type {number} */\n"+
				"  [Symbol.number];\n"+
				"}\n"))

	iCtor := c.globalScope.GetVar("I").GetType().AssertFunctionType()
	TypeSubject_AssertType(iCtor).ToStringIsEqualTo("(typeof I)")
	AssertThat(iCtor.IsStructuralInterface()).IsTrue()
	stringSym :=
		c.globalScope.GetVar("Symbol.string").GetType().ToMaybeKnownSymbolType()
	numberSym :=
		c.globalScope.GetVar("Symbol.number").GetType().ToMaybeKnownSymbolType()

	TypeSubject_AssertType(iCtor.GetInstanceType()).HasProperty(stringSym)
	TypeSubject_AssertType(iCtor.GetPrototype()).HasProperty(numberSym)
}

// Java: TypedScopeCreatorTest.testPropertiesOnInterface()
func (c *TypedScopeCreatorTest) TestPropertiesOnInterface() {
	c.TestSame(
		"/** @interface */ var I = function() {};\n" +
			"/** @type {number} */ I.prototype.bar;\n" +
			"I.prototype.baz = function(){};\n")

	i := c.globalScope.GetVar("I")
	AssertThat(i.GetType().String()).IsEqualTo("(typeof I)")
	AssertThat(i.GetType().IsInterface()).IsTrue()

	iPrototype := typedScopeCreatorTest_cast[rhino.ObjectType]((typedScopeCreatorTest_cast[rhino.ObjectType](i.GetType())).GetPropertyType("prototype"))
	AssertThat(iPrototype.String()).IsEqualTo("I.prototype")
	AssertThat(iPrototype.IsFunctionPrototypeType()).IsTrue()

	AssertThat(iPrototype.GetPropertyType("bar").String()).IsEqualTo("number")
	AssertThat(iPrototype.GetPropertyType("baz").String()).
		IsEqualTo("function(this:I): undefined")

	TypeSubject_AssertType(c.globalScope.GetVar("I.prototype").GetType()).IsEqualTo(iPrototype)
}

// Java: TypedScopeCreatorTest.testPropertiesOnInterface2()
func (c *TypedScopeCreatorTest) TestPropertiesOnInterface2() {
	c.TestSame(
		"/** @interface */ var I = function() {};\n" +
			"I.prototype = {baz: function(){}};\n" +
			"/** @type {number} */ I.prototype.bar;\n")

	i := c.globalScope.GetVar("I")
	AssertThat(i.GetType().String()).IsEqualTo("(typeof I)")
	AssertThat(i.GetType().IsInterface()).IsTrue()

	iPrototype := typedScopeCreatorTest_cast[rhino.ObjectType]((typedScopeCreatorTest_cast[rhino.ObjectType](i.GetType())).GetPropertyType("prototype"))
	AssertThat(iPrototype.String()).IsEqualTo("I.prototype")
	AssertThat(iPrototype.IsFunctionPrototypeType()).IsTrue()

	AssertThat(iPrototype.GetPropertyType("bar").String()).IsEqualTo("number")

	AssertThat(iPrototype.GetPropertyType("baz").String()).
		IsEqualTo("function(this:I): undefined")

	AssertThat(c.globalScope.GetVar("I.prototype").GetType()).IsEqualTo(iPrototype)
}

// TODO(johnlenz): A syntax for stubs using object literals?
// Java: TypedScopeCreatorTest.testStubsInExterns()
func (c *TypedScopeCreatorTest) TestStubsInExterns() {
	c.TestSame(
		c.Externs(
			"/** @constructor */ function Extern() {}\n"+
				"Extern.prototype.bar;\n"+
				"var e = new Extern(); e.baz;\n"),
		c.Srcs(
			"/** @constructor */ function Foo() {}\n"+
				"Foo.prototype.bar;\n"+
				"var f = new Foo(); f.baz;\n"))

	e := typedScopeCreatorTest_cast[rhino.ObjectType](c.globalScope.GetVar("e").GetType())
	AssertThat(e.GetPropertyType("bar").String()).IsEqualTo("?")
	AssertThat(e.GetPropertyType("baz").String()).IsEqualTo("?")

	f := typedScopeCreatorTest_cast[rhino.ObjectType](c.globalScope.GetVar("f").GetType())
	AssertThat(f.GetPropertyType("bar").String()).IsEqualTo("?")
	AssertThat(f.HasProperty("baz")).IsFalse()
}

// Java: TypedScopeCreatorTest.testStubsInExterns2()
func (c *TypedScopeCreatorTest) TestStubsInExterns2() {
	c.TestSame(
		c.Externs(
			"/** @constructor */ function Extern() {}\n"+
				"/** @type {Extern} */ var myExtern;\n"+
				"/** @type {number} */ myExtern.foo;\n"),
		c.Srcs(""))

	e := c.globalScope.GetVar("myExtern").GetType()
	AssertThat(e.String()).IsEqualTo("(Extern|null)")

	externType := typedScopeCreatorTest_cast[rhino.ObjectType](e.RestrictByNotNullOrUndefined())
	AssertWithMessage(c.globalScope.GetRootNode().ToStringTree()).
		That(externType.HasOwnProperty("foo")).
		IsTrue()
	AssertThat(externType.IsPropertyTypeDeclared("foo")).IsTrue()
	AssertThat(externType.GetPropertyType("foo").String()).IsEqualTo("number")
	AssertThat(externType.IsPropertyInExterns("foo")).IsTrue()
}

// Java: TypedScopeCreatorTest.testStubsInExterns3()
func (c *TypedScopeCreatorTest) TestStubsInExterns3() {
	c.TestSame(
		c.Externs(
			"/** @type {number} */ myExtern.foo;\n"+
				"/** @type {Extern} */ var myExtern;\n"+
				"/** @constructor */ function Extern() {}\n"),
		c.Srcs(""))

	e := c.globalScope.GetVar("myExtern").GetType()
	AssertThat(e.String()).IsEqualTo("(Extern|null)")

	externType := typedScopeCreatorTest_cast[rhino.ObjectType](e.RestrictByNotNullOrUndefined())
	AssertWithMessage(c.globalScope.GetRootNode().ToStringTree()).
		That(externType.HasOwnProperty("foo")).
		IsTrue()
	AssertThat(externType.IsPropertyTypeDeclared("foo")).IsTrue()
	AssertThat(externType.GetPropertyType("foo").String()).IsEqualTo("number")
	AssertThat(externType.IsPropertyInExterns("foo")).IsTrue()
}

// Java: TypedScopeCreatorTest.testStubsInExterns4()
func (c *TypedScopeCreatorTest) TestStubsInExterns4() {
	c.TestSame(
		c.Externs(
			"Extern.prototype.foo;\n"+
				"/** @constructor */ function Extern() {}\n"),
		c.Srcs(""))

	e := c.globalScope.GetVar("Extern").GetType()
	AssertThat(e.String()).IsEqualTo("(typeof Extern)")

	externProto := (typedScopeCreatorTest_cast[rhino.FunctionType](e)).GetPrototype()
	AssertWithMessage(c.globalScope.GetRootNode().ToStringTree()).
		That(externProto.HasOwnProperty("foo")).
		IsTrue()
	AssertThat(externProto.IsPropertyTypeInferred("foo")).IsTrue()
	AssertThat(externProto.GetPropertyType("foo").String()).IsEqualTo("?")
	AssertThat(externProto.IsPropertyInExterns("foo")).IsTrue()
}

// Java: TypedScopeCreatorTest.testPropertyInExterns1()
func (c *TypedScopeCreatorTest) TestPropertyInExterns1() {
	// Declaring a property on a non-native extern type (e.g. 'Extern') declares it as a property
	// on the instance type, but only in externs.
	c.TestSame(
		c.Externs(
			"/** @constructor */ function Extern() {}\n"+
				"/** @type {Extern} */ var extern;\n"+
				"/** @return {number} */ extern.one;\n"),
		c.Srcs(
			"/** @constructor */ function Normal() {}\n"+
				"/** @type {Normal} */ var normal;\n"+
				"/** @return {number} */ normal.one;\n"+
				"var result = new Extern().one();\n"))

	e := c.globalScope.GetVar("Extern").GetType()
	externInstance := (typedScopeCreatorTest_cast[rhino.FunctionType](e)).GetInstanceType()
	AssertThat(externInstance.HasOwnProperty("one")).IsTrue()
	AssertThat(externInstance.IsPropertyTypeDeclared("one")).IsTrue()
	AssertThat(externInstance.GetPropertyType("one").String()).IsEqualTo("function(): number")
	AssertThat(c.globalScope.GetVar("result").GetType().String()).IsEqualTo("number")

	n := c.globalScope.GetVar("Normal").GetType()
	normalInstance := (typedScopeCreatorTest_cast[rhino.FunctionType](n)).GetInstanceType()
	AssertThat(normalInstance.HasOwnProperty("one")).IsFalse()
}

// Java: TypedScopeCreatorTest.testPropertyInExterns2()
func (c *TypedScopeCreatorTest) TestPropertyInExterns2() {
	// Native extern types (such as Object) do not get stray properties declared, since this would
	// cause problems with bad externs (such as `/** @type {!Object} */ var api = {}; api.foo;`,
	// where we don't want to declare that all Objects have a "foo" property.  Nevertheless, the
	// specific qualified name (i.e. extern.one, in the example below) is still declared on the
	// global scope, so referring to the "one" property specifically on "extern" is still checked
	// as one would expect.
	c.TestSame(
		c.Externs(
			"/** @type {Object} */ var extern;\n"+
				"/** @return {number} */ extern.one;\n"),
		c.Srcs(
			"/** @type {Object} */ var normal;\n"+
				"/** @return {number} */ normal.one;\n"+
				"var result = extern.one();\n"))

	e := c.globalScope.GetVar("extern").GetType()
	AssertThat(e.Dereference().HasOwnProperty("one")).IsFalse()
	AssertThat(c.globalScope.GetVar("result").GetType().String()).IsEqualTo("number")

	normal := c.globalScope.GetVar("normal").GetType()
	AssertThat(normal.Dereference().HasOwnProperty("one")).IsFalse()
}

// Java: TypedScopeCreatorTest.testPropertyInExterns3()
func (c *TypedScopeCreatorTest) TestPropertyInExterns3() {
	c.TestSame(
		c.Externs(
			"/** @constructor\n"+
				" * @param {*=} x @return {!Object} */\n"+
				"function Object(x) {}\n"+
				"/** @type {number} */ Object.one;\n"),
		c.Srcs(""))

	obj := c.globalScope.GetVar("Object").GetType().Dereference()
	AssertThat(obj.HasOwnProperty("one")).IsTrue()
	AssertThat(obj.GetPropertyType("one").String()).IsEqualTo("number")
}

// Java: TypedScopeCreatorTest.testTypedStubsInExterns()
func (c *TypedScopeCreatorTest) TestTypedStubsInExterns() {
	c.TestSame(
		c.Externs(
			"/** @constructor\n"+
				" * @param {*} var_args */\n"+
				"function Function(var_args) {}\n"+
				"/** @type {!Function} */ Function.prototype.apply;\n"),
		c.Srcs("var f = new Function();"))

	f := typedScopeCreatorTest_cast[rhino.ObjectType](c.globalScope.GetVar("f").GetType())

	// The type of apply() on a function instance is resolved dynamically,
	// since apply varies with the type of the function it's called on.
	AssertThat(f.GetPropertyType("apply").String()).IsEqualTo("function(?=, (Object|null)=): ?")

	// The type of apply() on the function prototype just takes what it was
	// declared with.
	func_ := typedScopeCreatorTest_cast[rhino.FunctionType](c.globalScope.GetVar("Function").GetType())
	AssertThat(func_.GetPrototype().GetPropertyType("apply").String()).IsEqualTo("Function")
}

// Java: TypedScopeCreatorTest.testTypesInExterns()
func (c *TypedScopeCreatorTest) TestTypesInExterns() {
	c.TestSame(c.Externs(CompilerTypeTestCase_DEFAULT_EXTERNS), c.Srcs(""))

	v := c.globalScope.GetVar("Object")
	obj := typedScopeCreatorTest_cast[rhino.FunctionType](v.GetType())
	AssertThat(obj.String()).IsEqualTo("(typeof Object)")
	AssertThat(v.GetNode()).IsNotNull()
	AssertThat(v.GetInput()).IsNotNull()
}

// Java: TypedScopeCreatorTest.testSymbolsDeclaredInExterns()
func (c *TypedScopeCreatorTest) TestSymbolsDeclaredInExterns() {
	c.TestSame(
		c.Externs(
			CompilerTypeTestCase_DEFAULT_EXTERNS,
			"/** @const {symbol} */\n"+
				"Symbol.fooBar;\n"+
				"/** @type {number} */\n"+
				"Symbol.nonSymbolStatic;\n"),
		c.Srcs(""))

	fooBar := c.globalScope.GetVar("Symbol.fooBar")
	nonSymbolStatic := c.globalScope.GetVar("Symbol.nonSymbolStatic")

	TypeSubject_AssertType(fooBar.GetType()).ToStringIsEqualTo("Symbol.fooBar")
	TypeSubject_AssertType(nonSymbolStatic.GetType()).IsNumber()
}

// Java: TypedScopeCreatorTest.testSymbolIterator_includedByDefault()
func (c *TypedScopeCreatorTest) TestSymbolIterator_includedByDefault() {
	c.TestSame(c.Externs(""), c.Srcs(""))

	iterator := c.globalScope.GetVar("Symbol.iterator")

	TypeSubject_AssertType(iterator.GetType()).ToStringIsEqualTo("Symbol.iterator")
	AssertThat(iterator.GetType().IsKnownSymbolValueType()).IsTrue()
}

// Java: TypedScopeCreatorTest.testNewSymbolsInSrcCode()
func (c *TypedScopeCreatorTest) TestNewSymbolsInSrcCode() {
	c.TestSame(c.Srcs("/** @type {symbol} */ const s = Symbol();"))

	iterator := c.globalScope.GetVar("s")
	AssertThat(iterator.GetType().String()).IsEqualTo("symbol")
}

// Java: TypedScopeCreatorTest.testPropertyDeclarationOnInstanceType()
func (c *TypedScopeCreatorTest) TestPropertyDeclarationOnInstanceType() {
	c.TestSame(
		"/** @type {!Object} */ var a = {};\n" +
			"/** @type {number} */ a.name = 0;\n")

	AssertThat(c.globalScope.GetVar("a.name").GetType().String()).IsEqualTo("number")

	a := typedScopeCreatorTest_cast[rhino.ObjectType](c.globalScope.GetVar("a").GetType())
	AssertThat(a.HasProperty("name")).IsFalse()
	AssertThat(c.GetNativeObjectType(rhino.JSTypeNative_OBJECT_TYPE).HasProperty("name")).IsFalse()
}

// Java: TypedScopeCreatorTest.testPropertyDeclarationOnRecordType()
func (c *TypedScopeCreatorTest) TestPropertyDeclarationOnRecordType() {
	c.TestSame(
		"/** @type {{foo: number}} */ var a = {foo: 3};\n" +
			"/** @type {number} */ a.name = 0;\n")

	AssertThat(c.globalScope.GetVar("a.name").GetType().String()).IsEqualTo("number")

	a := typedScopeCreatorTest_cast[rhino.ObjectType](c.globalScope.GetVar("a").GetType())
	AssertThat(a.String()).IsEqualTo("{foo: number}")
	AssertThat(a.HasProperty("name")).IsFalse()
}

// Java: TypedScopeCreatorTest.testGlobalThis1()
func (c *TypedScopeCreatorTest) TestGlobalThis1() {
	c.TestSame(
		"/** @constructor */ function Window() {}\n" +
			"Window.prototype.alert = function() {};\n" +
			"var x = this;\n")

	x := typedScopeCreatorTest_cast[rhino.ObjectType](c.globalScope.GetVar("x").GetType())
	windowCtor := typedScopeCreatorTest_cast[rhino.FunctionType](c.globalScope.GetVar("Window").GetType())
	AssertThat(x.String()).IsEqualTo("global this")
	AssertThat(x.IsSubtypeOf(windowCtor.GetInstanceType())).IsTrue()
	AssertThat(x.Equals(windowCtor.GetInstanceType())).IsFalse()
	AssertThat(x.HasProperty("alert")).IsTrue()
}

// Java: TypedScopeCreatorTest.testGlobalThis2()
func (c *TypedScopeCreatorTest) TestGlobalThis2() {
	c.TestSame(
		"/** @constructor */ function Window() {}\n" +
			"Window.prototype = {alert: function() {}};\n" +
			"var x = this;\n")

	x := typedScopeCreatorTest_cast[rhino.ObjectType](c.globalScope.GetVar("x").GetType())
	windowCtor := typedScopeCreatorTest_cast[rhino.FunctionType](c.globalScope.GetVar("Window").GetType())
	AssertThat(x.String()).IsEqualTo("global this")
	AssertThat(x.IsSubtypeOf(windowCtor.GetInstanceType())).IsTrue()
	AssertThat(x.Equals(windowCtor.GetInstanceType())).IsFalse()
	AssertThat(x.HasProperty("alert")).IsTrue()
}

// Java: TypedScopeCreatorTest.testGlobalThis_includesProvidedNamesAsProperties()
func (c *TypedScopeCreatorTest) TestGlobalThis_includesProvidedNamesAsProperties() {
	c.processClosurePrimitives = true
	c.TestSame(
		c.Srcs(
			"goog.provide('my.alert');\n" +
				"/** @param {string} msg */\n" +
				"my.alert = function(msg) {};\n" +
				"var x = this;\n"))

	myType := c.FindNameType("my", c.globalScope).ToObjectType()
	globalThis := c.FindNameType("x", c.globalScope).ToObjectType()
	TypeSubject_AssertType(globalThis).ToStringIsEqualTo("global this")
	TypeSubject_AssertType(globalThis).WithTypeOfProp("my").IsSameInstanceAs(myType)
}

// Java: TypedScopeCreatorTest.testObjectLiteralCast()
func (c *TypedScopeCreatorTest) TestObjectLiteralCast() {
	// Verify that "goog.reflect.object" does not modify the types on
	// "A.B"
	c.TestSame(
		"var A = {};\n" +
			"/** @constructor */ A.B = function() {}\n" +
			"A.B.prototype.isEnabled = true;\n" +
			"goog.reflect.object(A.B, {isEnabled: 3})\n" +
			"var x = (new A.B()).isEnabled;\n")

	// Verify that "$jscomp.reflectObject" does not modify the types on
	// "A.B"
	c.TestSame(
		"var A = {};\n" +
			"var $jscomp = {};\n" +
			"/** @constructor */ A.B = function() {}\n" +
			"A.B.prototype.isEnabled = true;\n" +
			"$jscomp.reflectObject(A.B, {isEnabled: 3})\n" +
			"var x = (new A.B()).isEnabled;\n")

	AssertThat(c.FindTokenType(rhino.Token_OBJECTLIT, c.globalScope).String()).IsEqualTo("A.B")
	AssertThat(c.FindNameType("x", c.globalScope).String()).IsEqualTo("boolean")
}

// Java: TypedScopeCreatorTest.testBadObjectLiteralCast1()
func (c *TypedScopeCreatorTest) TestBadObjectLiteralCast1() {
	c.TestWarning(
		"var A = {}; /** @constructor */ A.B = function() {}\n"+
			"goog.reflect.object(A.B, 1)\n",
		ClosureCodingConvention_OBJECTLIT_EXPECTED)
}

// Java: TypedScopeCreatorTest.testBadObjectLiteralCast2()
func (c *TypedScopeCreatorTest) TestBadObjectLiteralCast2() {
	c.TestWarning("var A; goog.reflect.object(A.B, {})", TypedScopeCreator_CONSTRUCTOR_EXPECTED)
}

// Java: TypedScopeCreatorTest.testConstructorNode()
func (c *TypedScopeCreatorTest) TestConstructorNode() {
	c.TestSame("/** @constructor */ goog.Foo = function() {};")

	ctor := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("goog.Foo", c.globalScope))
	AssertThat(ctor).IsNotNull()
	AssertThat(ctor.IsConstructor()).IsTrue()
	AssertThat(ctor.String()).IsEqualTo("(typeof goog.Foo)")
}

// Java: TypedScopeCreatorTest.testClassDeclaration()
func (c *TypedScopeCreatorTest) TestClassDeclaration() {
	c.TestSame("class Foo {}")
	foo := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Foo", c.globalScope))
	AssertThat(foo.IsConstructor()).IsTrue()
	ScopeSubject_AssertScope(c.globalScope).Declares("Foo").WithTypeThat().IsEqualTo(foo)
}

// Java: TypedScopeCreatorTest.testClassDeclarationWithoutConstructor()
func (c *TypedScopeCreatorTest) TestClassDeclarationWithoutConstructor() {
	c.TestSame("class Foo {}")

	fooClass := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Foo", c.globalScope))
	fooProto := fooClass.GetPrototype()

	// Test class typing.
	AssertThat(fooClass.IsConstructor()).IsTrue()

	// Test constructor property.
	AssertThat(fooProto.HasOwnProperty("constructor")).IsTrue()
	AssertNode(fooProto.GetOwnPropertyDefSite("constructor")).IsNull()
	TypeSubject_AssertType(fooProto).WithTypeOfProp("constructor").ToStringIsEqualTo("(typeof Foo)")
}

// Java: TypedScopeCreatorTest.testClassDeclarationWithConstructor()
func (c *TypedScopeCreatorTest) TestClassDeclarationWithConstructor() {
	c.TestSame(
		"class Foo {\n" +
			"  /** @param {number} arg */\n" +
			"  constructor(arg) {\n" +
			"    CTOR_BODY:;\n" +
			"  }\n" +
			"}\n")

	fooClass := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Foo", c.globalScope))
	fooProto := fooClass.GetPrototype()
	params := fooClass.GetParameters()
	ctorDef := c.GetLabeledStatement("CTOR_BODY").statementNode.GetAncestor(3)

	// Test class typing.
	AssertThat(fooClass.IsConstructor()).IsTrue()
	AssertThat(params).HasSize(1)
	TypeSubject_AssertType(params[0].GetJSType()).IsNumber()

	// Test constructor property.
	AssertThat(fooProto.HasOwnProperty("constructor")).IsTrue()
	AssertNode(fooProto.GetOwnPropertyDefSite("constructor")).IsSameInstanceAs(ctorDef)
	TypeSubject_AssertType(fooProto).WithTypeOfProp("constructor").ToStringIsEqualTo("(typeof Foo)")
}

// Java: TypedScopeCreatorTest.testInterfaceClassDeclarationWithConstructor()
func (c *TypedScopeCreatorTest) TestInterfaceClassDeclarationWithConstructor() {
	c.TestSame(
		"/** @interface */\n" +
			"class Foo {\n" +
			"  constructor(arg) {}\n" +
			"}\n")

	fooClass := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Foo", c.globalScope))
	fooProto := fooClass.GetPrototype()

	// Test class typing.
	AssertThat(fooClass.IsInterface()).IsTrue()

	// Test constructor property.
	AssertThat(fooProto.HasOwnProperty("constructor")).IsFalse()
}

// Java: TypedScopeCreatorTest.testClassDeclarationWithImplementsForwardRef()
func (c *TypedScopeCreatorTest) TestClassDeclarationWithImplementsForwardRef() {
	c.TestSame(
		"/** @implements {Bar} */\n" +
			"class Foo {}\n" +
			"/** @interface */\n" +
			"class Bar {}\n" +
			"class SubFoo extends Foo {}\n")

	fooCtorType := c.globalScope.GetVar("Foo").GetType().ToMaybeFunctionType()
	fooType := fooCtorType.GetInstanceType()

	AssertThat(fooCtorType.LoosenTypecheckingDueToForwardReferencedSupertype()).IsFalse()
	AssertThat(fooType.LoosenTypecheckingDueToForwardReferencedSupertype()).IsTrue()

	subFooCtorType := c.globalScope.GetVar("SubFoo").GetType().ToMaybeFunctionType()
	subFooType := subFooCtorType.GetInstanceType()

	AssertThat(subFooCtorType.LoosenTypecheckingDueToForwardReferencedSupertype()).IsTrue()
	AssertThat(subFooType.LoosenTypecheckingDueToForwardReferencedSupertype()).IsTrue()
}

// Java: TypedScopeCreatorTest.testClassDeclarationWithExtends()
func (c *TypedScopeCreatorTest) TestClassDeclarationWithExtends() {
	c.TestSame(
		"class Bar {}\n" +
			"class Foo extends Bar {}\n")
	bar := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Bar", c.globalScope))
	foo := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Foo", c.globalScope))
	TypeSubject_AssertType(foo.GetInstanceType()).IsSubtypeOf(bar.GetInstanceType())
	TypeSubject_AssertType(foo.GetImplicitPrototype()).IsEqualTo(bar)
	ScopeSubject_AssertScope(c.globalScope).Declares("Bar").WithTypeThat().IsEqualTo(bar)
	ScopeSubject_AssertScope(c.globalScope).Declares("Foo").WithTypeThat().IsEqualTo(foo)

	AssertThat(foo.GetInstanceType().LoosenTypecheckingDueToForwardReferencedSupertype()).IsFalse()
	AssertThat(foo.LoosenTypecheckingDueToForwardReferencedSupertype()).IsFalse()
	AssertThat(foo.IsAmbiguousConstructor()).IsFalse()
}

// Java: TypedScopeCreatorTest.testClassDeclarationWithExtends_andCompatibleJsdoc()
func (c *TypedScopeCreatorTest) TestClassDeclarationWithExtends_andCompatibleJsdoc() {
	c.TestSame(
		"/** @template T */\n" +
			"class Bar {}\n" +
			"/** @extends {Bar<string>} */\n" +
			"class Foo extends Bar {}\n")
	bar := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Bar", c.globalScope))
	foo := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Foo", c.globalScope))
	TypeSubject_AssertType(foo.GetInstanceType()).IsSubtypeOf(bar.GetInstanceType())
	ScopeSubject_AssertScope(c.globalScope).Declares("Bar").WithTypeThat().IsEqualTo(bar)
	ScopeSubject_AssertScope(c.globalScope).Declares("Foo").WithTypeThat().IsEqualTo(foo)

	AssertThat(foo.GetInstanceType().LoosenTypecheckingDueToForwardReferencedSupertype()).IsFalse()
	AssertThat(foo.LoosenTypecheckingDueToForwardReferencedSupertype()).IsFalse()
	AssertThat(foo.IsAmbiguousConstructor()).IsFalse()
}

// Java: TypedScopeCreatorTest.testClassDeclarationWithExtends_andIncompatibleJsdoc()
func (c *TypedScopeCreatorTest) TestClassDeclarationWithExtends_andIncompatibleJsdoc() {
	c.TestSame(
		"class Bar {}\n" +
			"class Qux {}\n" +
			"/** @extends {Bar} */\n" +
			"class Foo extends Qux {}\n")
	bar := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Bar", c.globalScope))
	qux := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Qux", c.globalScope))
	foo := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Foo", c.globalScope))
	TypeSubject_AssertType(foo.GetInstanceType()).IsSubtypeOf(bar.GetInstanceType())
	ScopeSubject_AssertScope(c.globalScope).Declares("Bar").WithTypeThat().IsEqualTo(bar)
	ScopeSubject_AssertScope(c.globalScope).Declares("Qux").WithTypeThat().IsEqualTo(qux)
	ScopeSubject_AssertScope(c.globalScope).Declares("Foo").WithTypeThat().IsEqualTo(foo)

	AssertThat(foo.GetInstanceType().LoosenTypecheckingDueToForwardReferencedSupertype()).IsFalse()
	AssertThat(foo.LoosenTypecheckingDueToForwardReferencedSupertype()).IsFalse()
	AssertThat(foo.IsAmbiguousConstructor()).IsTrue()
}

// Java: TypedScopeCreatorTest.testAbstractClassDeclarationWithExtends_nameDeclAndAssign()
func (c *TypedScopeCreatorTest) TestAbstractClassDeclarationWithExtends_nameDeclAndAssign() {
	c.TestSame(
		"var Bar_1;\n" +
			"/** @abstract */\n" +
			"let Bar = Bar_1 = class Bar {}\n" +
			"class Foo extends Bar {}\n")

	bar1 := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Bar_1", c.globalScope))
	bar := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Bar", c.globalScope))
	foo := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Foo", c.globalScope))
	AssertThat(bar.IsAbstract()).IsTrue()
	AssertThat(bar).IsEqualTo(bar1)
	TypeSubject_AssertType(foo.GetInstanceType()).IsSubtypeOf(bar.GetInstanceType())
	TypeSubject_AssertType(foo.GetImplicitPrototype()).IsEqualTo(bar)
	ScopeSubject_AssertScope(c.globalScope).Declares("Bar").WithTypeThat().IsEqualTo(bar)
	ScopeSubject_AssertScope(c.globalScope).Declares("Foo").WithTypeThat().IsEqualTo(foo)

	AssertThat(foo.GetInstanceType().LoosenTypecheckingDueToForwardReferencedSupertype()).IsFalse()
	AssertThat(foo.LoosenTypecheckingDueToForwardReferencedSupertype()).IsFalse()
	AssertThat(foo.IsAmbiguousConstructor()).IsFalse()
}

// Java: TypedScopeCreatorTest.testClassDeclarationWithExtendsFromNamespace()
func (c *TypedScopeCreatorTest) TestClassDeclarationWithExtendsFromNamespace() {
	c.TestSame(
		"const ns = {};\n" +
			"ns.Bar = class {};\n" +
			"const nsAliased = ns;\n" +
			"class Foo extends nsAliased.Bar {}\n")
	bar := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("ns.Bar", c.globalScope))
	foo := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Foo", c.globalScope))

	TypeSubject_AssertType(foo.GetInstanceType()).IsSubtypeOf(bar.GetInstanceType())
	TypeSubject_AssertType(foo.GetImplicitPrototype()).IsEqualTo(bar)

	AssertThat(foo.GetInstanceType().LoosenTypecheckingDueToForwardReferencedSupertype()).IsFalse()
	AssertThat(foo.LoosenTypecheckingDueToForwardReferencedSupertype()).IsFalse()
	AssertThat(foo.IsAmbiguousConstructor()).IsFalse()
}

// Java: TypedScopeCreatorTest.testClassDeclarationWithExtendsFromNamespaceAndJSDoc()
func (c *TypedScopeCreatorTest) TestClassDeclarationWithExtendsFromNamespaceAndJSDoc() {
	c.TestSame(
		"const ns = {};\n" +
			"ns.Bar = class {};\n" +
			"const nsAliased = ns;\n" +
			"\n" +
			"/** @extends {nsAliased.Bar} */\n" +
			"class Foo extends nsAliased.Bar {}\n")
	bar := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("ns.Bar", c.globalScope))
	foo := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Foo", c.globalScope))

	TypeSubject_AssertType(foo.GetInstanceType()).IsSubtypeOf(bar.GetInstanceType())
	TypeSubject_AssertType(foo.GetImplicitPrototype()).IsEqualTo(bar)

	AssertThat(foo.GetInstanceType().LoosenTypecheckingDueToForwardReferencedSupertype()).IsFalse()
	AssertThat(foo.LoosenTypecheckingDueToForwardReferencedSupertype()).IsFalse()
	AssertThat(foo.IsAmbiguousConstructor()).IsFalse()
}

// Java: TypedScopeCreatorTest.testClassDeclarationWithExtendsFromNamespaceAndJSDoc_incompatible()
func (c *TypedScopeCreatorTest) TestClassDeclarationWithExtendsFromNamespaceAndJSDoc_incompatible() {
	c.TestSame(
		"const ns = {};\n" +
			"/** @template T */\n" +
			"ns.Bar = class {};\n" +
			"ns.Qux = class {};\n" +
			"const nsAliased = ns;\n" +
			"\n" +
			"/** @extends {nsAliased.Bar<number>} */\n" +
			"class Foo extends nsAliased.Bar {}\n")
	bar := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("ns.Bar", c.globalScope))
	foo := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Foo", c.globalScope))

	TypeSubject_AssertType(foo.GetInstanceType()).IsSubtypeOf(bar.GetInstanceType())
	TypeSubject_AssertType(foo.GetImplicitPrototype()).IsEqualTo(bar)

	AssertThat(foo.GetInstanceType().LoosenTypecheckingDueToForwardReferencedSupertype()).IsFalse()
	AssertThat(foo.LoosenTypecheckingDueToForwardReferencedSupertype()).IsFalse()
	AssertThat(foo.IsAmbiguousConstructor()).IsFalse()
}

// Java: TypedScopeCreatorTest.testClassDeclarationWithNestedExtends()
func (c *TypedScopeCreatorTest) TestClassDeclarationWithNestedExtends() {
	c.TestSame(
		"class Bar {}\n" +
			"class Foo extends class extends class extends Bar {} {} {}\n")
	bar := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Bar", c.globalScope))
	foo := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Foo", c.globalScope))
	TypeSubject_AssertType(foo.GetInstanceType()).IsSubtypeOf(bar.GetInstanceType())
}

// Java: TypedScopeCreatorTest.testClassDeclarationWithExtends_googModuleGet()
func (c *TypedScopeCreatorTest) TestClassDeclarationWithExtends_googModuleGet() {
	c.TestSame(
		c.Srcs(
			"goog.module('a.Bar');\n"+
				"class Bar {}\n"+
				"BAR: Bar;\n"+
				"exports = Bar;\n",
			"class Foo extends goog.module.get('a.Bar') {}"))
	moduleScope := c.GetLabeledStatement("BAR").enclosingScope
	foo := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Foo", c.globalScope))
	bar := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Bar", moduleScope))
	TypeSubject_AssertType(foo.GetInstanceType()).IsSubtypeOf(bar.GetInstanceType())
	TypeSubject_AssertType(foo.GetImplicitPrototype()).IsEqualTo(bar)
	ScopeSubject_AssertScope(moduleScope).Declares("Bar").WithTypeThat().IsEqualTo(bar)
	ScopeSubject_AssertScope(c.globalScope).Declares("Foo").WithTypeThat().IsEqualTo(foo)

	AssertThat(foo.GetInstanceType().LoosenTypecheckingDueToForwardReferencedSupertype()).IsFalse()
	AssertThat(foo.LoosenTypecheckingDueToForwardReferencedSupertype()).IsFalse()
	AssertThat(foo.IsAmbiguousConstructor()).IsFalse()
}

// Java: TypedScopeCreatorTest.testClassDeclarationWithExtends_googModuleGetProperty()
func (c *TypedScopeCreatorTest) TestClassDeclarationWithExtends_googModuleGetProperty() {
	c.TestSame(
		c.Srcs(
			"goog.module('a.b');\n"+
				"class Bar {}\n"+
				"BAR: Bar;\n"+
				"exports = {Bar};\n",
			"class Foo extends goog.module.get('a.b').Bar {}"))
	moduleScope := c.GetLabeledStatement("BAR").enclosingScope
	foo := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Foo", c.globalScope))
	bar := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Bar", moduleScope))
	TypeSubject_AssertType(foo.GetInstanceType()).IsSubtypeOf(bar.GetInstanceType())
	TypeSubject_AssertType(foo.GetImplicitPrototype()).IsEqualTo(bar)
	ScopeSubject_AssertScope(moduleScope).Declares("Bar").WithTypeThat().IsEqualTo(bar)
	ScopeSubject_AssertScope(c.globalScope).Declares("Foo").WithTypeThat().IsEqualTo(foo)

	AssertThat(foo.GetInstanceType().LoosenTypecheckingDueToForwardReferencedSupertype()).IsFalse()
	AssertThat(foo.LoosenTypecheckingDueToForwardReferencedSupertype()).IsFalse()
	AssertThat(foo.IsAmbiguousConstructor()).IsFalse()
}

// Java: TypedScopeCreatorTest.testClassDeclarationWithExtends_googModuleGetNestedProperty()
func (c *TypedScopeCreatorTest) TestClassDeclarationWithExtends_googModuleGetNestedProperty() {
	c.TestSame(
		c.Srcs(
			"goog.module('a.b');\n"+
				"class Bar {}\n"+
				"Bar.NestedClass = class {};\n"+
				"NESTED_CLASS: Bar.NestedClass;\n"+
				"exports = {Bar};\n",
			"class Foo extends goog.module.get('a.b').Bar.NestedClass {}"))
	moduleScope := c.GetLabeledStatement("NESTED_CLASS").enclosingScope
	foo := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Foo", c.globalScope))
	nestedClass :=
		typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Bar", moduleScope).FindPropertyType("NestedClass"))
	TypeSubject_AssertType(foo.GetInstanceType()).IsSubtypeOf(nestedClass.GetInstanceType())
	TypeSubject_AssertType(foo.GetImplicitPrototype()).IsEqualTo(nestedClass)

	AssertThat(foo.GetInstanceType().LoosenTypecheckingDueToForwardReferencedSupertype()).IsFalse()
	AssertThat(foo.LoosenTypecheckingDueToForwardReferencedSupertype()).IsFalse()
	AssertThat(foo.IsAmbiguousConstructor()).IsFalse()
}

// Java: TypedScopeCreatorTest.testClassDeclarationWithExtends_missingGoogModuleGet()
func (c *TypedScopeCreatorTest) TestClassDeclarationWithExtends_missingGoogModuleGet() {
	c.TestSame(c.Srcs("class Foo extends goog.module.get('not.a.real.module') {}"))
	foo := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Foo", c.globalScope))
	ScopeSubject_AssertScope(c.globalScope).Declares("Foo").WithTypeThat().IsEqualTo(foo)

	TypeSubject_AssertType(foo.GetPrototype().GetImplicitPrototype()).IsEqualTo(c.GetNativeType(rhino.JSTypeNative_OBJECT_TYPE))
}

// Java: TypedScopeCreatorTest.testClassDeclarationWithInheritedConstructor()
func (c *TypedScopeCreatorTest) TestClassDeclarationWithInheritedConstructor() {
	c.TestSame(
		"class Bar {\n" +
			"  constructor(/** string */ arg) {}\n" +
			"}\n" +
			"class Foo extends Bar {}\n")

	fooClass := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Foo", c.globalScope))
	fooProto := fooClass.GetPrototype()
	params := fooClass.GetParameters()

	// Test class typing.
	AssertThat(fooClass.IsConstructor()).IsTrue()
	AssertThat(params).HasSize(1)
	TypeSubject_AssertType(params[0].GetJSType()).IsString()

	// Test constructor property.
	AssertThat(fooProto.HasOwnProperty("constructor")).IsTrue()
	AssertNode(fooProto.GetOwnPropertyDefSite("constructor")).IsNull()
	TypeSubject_AssertType(fooProto).WithTypeOfProp("constructor").ToStringIsEqualTo("(typeof Foo)")
}

// Java: TypedScopeCreatorTest.testClassDeclarationWithOverriddenConstructor()
func (c *TypedScopeCreatorTest) TestClassDeclarationWithOverriddenConstructor() {
	c.TestSame(
		"class Bar {\n" +
			"  constructor(/** string */ arg) {}\n" +
			"}\n" +
			"class Foo extends Bar {\n" +
			"  constructor(/** number */ arg) { CTOR_BODY:super(''); }\n" +
			"  static method() {}\n" +
			"}\n")

	barCtor := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Bar", c.globalScope))
	barObject := barCtor.GetInstanceType()
	barConstructorProperty := barObject.GetPropertyType("constructor")

	fooCtor := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Foo", c.globalScope))
	fooObject := fooCtor.GetInstanceType()
	fooProto := fooCtor.GetPrototype()
	fooConstructorProperty := fooObject.GetPropertyType("constructor")

	superInvocation := c.GetLabeledStatement("CTOR_BODY").statementNode.GetFirstChild()

	superRef := superInvocation.GetFirstChild()
	AssertNode(superRef).HasToken(rhino.Token_SUPER).HasJSTypeThat().IsEqualTo(barCtor)

	TypeSubject_AssertType(fooCtor).WithTypeOfProp("method").IsNotUnknown()
	TypeSubject_AssertType(fooCtor).WithTypeOfProp("method").IsNotEmpty()

	params := fooCtor.GetParameters()
	AssertThat(params).HasSize(1)
	TypeSubject_AssertType(params[0].GetJSType()).IsNumber()
	TypeSubject_AssertType(barConstructorProperty).ToStringIsEqualTo("(typeof Bar)")

	fooCtorDef := NodeUtil_GetEnclosingFunction(superInvocation)
	TypeSubject_AssertType(fooConstructorProperty).ToStringIsEqualTo("(typeof Foo)")
	AssertNode(fooProto.GetOwnPropertyDefSite("constructor")).IsSameInstanceAs(fooCtorDef)

	TypeSubject_AssertType(fooConstructorProperty).IsSubtypeOf(barConstructorProperty)
	TypeSubject_AssertType(fooConstructorProperty).
		WithTypeOfProp("method").
		IsFunctionTypeThat().
		HasTypeOfThisThat().
		ToStringIsEqualTo("(typeof Foo)")
}

// Java: TypedScopeCreatorTest.testClassDeclarationWithNestedExtendsAndInheritedConstructor()
func (c *TypedScopeCreatorTest) TestClassDeclarationWithNestedExtendsAndInheritedConstructor() {
	c.TestSame(
		"class Bar {\n" +
			"  constructor(/** string */ arg) {}\n" +
			"}\n" +
			"class Foo extends class extends class extends Bar {} {} {}\n")
	foo := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Foo", c.globalScope))
	params := foo.GetParameters()
	AssertThat(params).HasSize(1)
	TypeSubject_AssertType(params[0].GetJSType()).IsString()
}

// Java: TypedScopeCreatorTest.testClassDeclarationWithMethod()
func (c *TypedScopeCreatorTest) TestClassDeclarationWithMethod() {
	c.TestSame(
		"class Foo {\n" +
			"  /** @param {string} arg */\n" +
			"  method(arg) {\n" +
			"    METHOD:;\n" +
			"    var /** number */ foo;\n" +
			"  }\n" +
			"}\n")
	methodBlockScope := c.GetLabeledStatement("METHOD").enclosingScope
	methodScope := methodBlockScope.GetParentScope().(*TypedScope)
	ScopeSubject_AssertScope(methodBlockScope).Declares("foo").Directly().WithTypeThat().IsNumber()
	ScopeSubject_AssertScope(methodScope).Declares("arg").Directly().WithTypeThat().IsString()

	method := typedScopeCreatorTest_cast[rhino.FunctionType](methodScope.GetRootNode().GetJSType())
	TypeSubject_AssertType(method).ToStringIsEqualTo("function(this:Foo, string): undefined")
	foo := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Foo", c.globalScope))
	TypeSubject_AssertType(foo.GetInstanceType()).WithTypeOfProp("method").IsEqualTo(method)
}

// Java: TypedScopeCreatorTest.testEs5ClassWithWellKnownSymbolTypedMethod()
func (c *TypedScopeCreatorTest) TestEs5ClassWithWellKnownSymbolTypedMethod() {
	c.TestSame(
		c.Externs(
			NewTestExternsBuilder().AddIterable().Build(),
			"/** @const {symbol}  */\n"+
				"Symbol.dispose;\n"),
		c.Srcs(
			"/** @constructor */\n"+
				"function Foo() {}\n"+
				"/**\n"+
				" * @param {string} x\n"+
				" * @param {number} y\n"+
				" */\n"+
				"Foo.prototype[Symbol.iterator] = function(x, y) {}\n"+
				"/** @return {symbol} */\n"+
				"Foo.prototype[Symbol.dispose] = function() { return Symbol.dispose; }\n"))

	foo := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Foo", c.globalScope))
	symbolIteratorType :=
		typedScopeCreatorTest_cast[*rhino.KnownSymbolType](c.FindNameType("Symbol.iterator", c.globalScope))
	symbolDisposeType :=
		typedScopeCreatorTest_cast[*rhino.KnownSymbolType](c.FindNameType("Symbol.dispose", c.globalScope))
	TypeSubject_AssertType(foo.GetInstanceType()).HasProperty(symbolIteratorType)
	TypeSubject_AssertType(foo.GetInstanceType()).HasProperty(symbolDisposeType)

	TypeSubject_AssertType(foo.GetInstanceType()).
		WithTypeOfPropSymbol(symbolIteratorType).
		ToStringIsEqualTo("function(this:Foo, string, number): undefined")
	TypeSubject_AssertType(foo.GetInstanceType()).
		WithTypeOfPropSymbol(symbolDisposeType).
		ToStringIsEqualTo("function(this:Foo): symbol")
}

// Java: TypedScopeCreatorTest.testClassWithWellKnownSymbolTypedMethod()
func (c *TypedScopeCreatorTest) TestClassWithWellKnownSymbolTypedMethod() {
	c.TestSame(
		c.Externs(
			NewTestExternsBuilder().AddIterable().Build(),
			"/** @const {symbol}  */\n"+
				"Symbol.dispose;\n"),
		c.Srcs(
			"class Foo {\n"+
				"  /**\n"+
				"   * @param {string} x\n"+
				"   * @param {number} y\n"+
				"   */\n"+
				"  [Symbol.iterator](x, y) {}\n"+
				"  /** @return {symbol} */\n"+
				"  [Symbol.dispose]() { return Symbol.dispose; }\n"+
				"}\n"))

	foo := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Foo", c.globalScope))
	symbolIteratorType :=
		typedScopeCreatorTest_cast[*rhino.KnownSymbolType](c.FindNameType("Symbol.iterator", c.globalScope))
	symbolDisposeType :=
		typedScopeCreatorTest_cast[*rhino.KnownSymbolType](c.FindNameType("Symbol.dispose", c.globalScope))
	TypeSubject_AssertType(foo.GetInstanceType()).HasProperty(symbolIteratorType)
	TypeSubject_AssertType(foo.GetInstanceType()).HasProperty(symbolDisposeType)

	TypeSubject_AssertType(foo.GetInstanceType()).
		WithTypeOfPropSymbol(symbolIteratorType).
		ToStringIsEqualTo("function(this:Foo, string, number): undefined")
	TypeSubject_AssertType(foo.GetInstanceType()).
		WithTypeOfPropSymbol(symbolDisposeType).
		ToStringIsEqualTo("function(this:Foo): symbol")
}

// Java: TypedScopeCreatorTest.testClassWithWellKnownSymbolTypedMethod_abstract()
func (c *TypedScopeCreatorTest) TestClassWithWellKnownSymbolTypedMethod_abstract() {
	c.TestSame(
		c.Externs(
			NewTestExternsBuilder().AddIterable().Build(),
			"/** @const {symbol}  */\n"+
				"Symbol.dispose;\n"),
		c.Srcs(
			"/** @abstract */\n"+
				"class Foo {\n"+
				"  /**\n"+
				"   * @param {string} x\n"+
				"   * @param {number} y\n"+
				"   * @abstract\n"+
				"   */\n"+
				"  [Symbol.iterator](x, y) {}\n"+
				"}\n"))

	foo := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Foo", c.globalScope))
	symbolIteratorType :=
		typedScopeCreatorTest_cast[*rhino.KnownSymbolType](c.FindNameType("Symbol.iterator", c.globalScope))
	TypeSubject_AssertType(foo.GetInstanceType()).HasProperty(symbolIteratorType)
	TypeSubject_AssertType(foo.GetInstanceType().GetPropertyTypeByKey(rhino.NewProperty_SymbolKey(symbolIteratorType))).
		IsFunctionTypeThat().
		IsAbstract()
}

// Java: TypedScopeCreatorTest.testClassDeclarationWithDucplicatePrototypeMethodAfterward()
func (c *TypedScopeCreatorTest) TestClassDeclarationWithDucplicatePrototypeMethodAfterward() {
	// Given
	c.TestSame(
		"class Foo {\n" +
			"  /** @param {string} arg */\n" +
			"  method(arg) { }\n" +
			"}\n" +
			"\n" +
			"// This declaration should be ignored in preference of the one in the CLASS body.\n" +
			"// Constructs like this are sometimes valid (e.g. in mod files) and have error reporting\n" +
			"// separate from `TypedScopeCreator`.\n" +
			"/** @param {number} arg */\n" +
			"Foo.prototype.method = function(arg) { }\n")

	// Then
	foo := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Foo", c.globalScope))
	TypeSubject_AssertType(foo.GetInstanceType()).
		WithTypeOfProp("method").
		ToStringIsEqualTo("function(this:Foo, string): undefined")
}

// Java: TypedScopeCreatorTest.testClassDeclarationWithDucplicateStaticMethodAfterward()
func (c *TypedScopeCreatorTest) TestClassDeclarationWithDucplicateStaticMethodAfterward() {
	// Given
	c.TestSame(
		"class Foo {\n" +
			"  /** @param {string} arg */\n" +
			"  static method(arg) { }\n" +
			"}\n" +
			"\n" +
			"// This declaration should be ignored in preference of the one in the CLASS body.\n" +
			"// Constructs like this are sometimes valid (e.g. in mod files) and have error reporting\n" +
			"// separate from `TypedScopeCreator`.\n" +
			"/** @param {number} arg */\n" +
			"Foo.method = function(arg) { }\n")

	// Then
	foo := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Foo", c.globalScope))
	TypeSubject_AssertType(foo).
		WithTypeOfProp("method").
		ToStringIsEqualTo("function(this:(typeof Foo), string): undefined")
}

// Java: TypedScopeCreatorTest.testClassDeclarationWithConflictingGetterSetter()
func (c *TypedScopeCreatorTest) TestClassDeclarationWithConflictingGetterSetter() {
	// Given
	c.TestSame(
		"class Foo {\n" +
			"// This is not valid, but will be warned about in TypeCheck later\n" +
			"  /** @return {string} */\n" +
			"  get field() { }\n" +
			"  /** @param {number} arg */\n" +
			"  set field(arg) {}\n" +
			"}\n")

	// Then
	foo := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Foo", c.globalScope))
	TypeSubject_AssertType(foo.GetPrototype()).WithTypeOfProp("field").IsString()
}

// Java: TypedScopeCreatorTest.testClassDeclarationWithGetterWithDuplicateAfterward()
func (c *TypedScopeCreatorTest) TestClassDeclarationWithGetterWithDuplicateAfterward() {
	// Given
	c.TestSame(
		"class Foo {\n" +
			"  /** @return {string} */\n" +
			"  get field() { }\n" +
			"}\n" +
			"\n" +
			"// This declaration should be ignored in preference of the one in the CLASS body.\n" +
			"// Constructs like this are sometimes valid (e.g. in mod files) and have error reporting\n" +
			"// separate from `TypedScopeCreator`.\n" +
			"/** @type {string} arg */\n" +
			"Foo.method = function(arg) { }\n")

	// Then
	foo := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Foo", c.globalScope))
	TypeSubject_AssertType(foo.GetPrototype()).WithTypeOfProp("field").IsString()
}

// Java: TypedScopeCreatorTest.testClassDeclarationWithMethodAndInlineParamDocs()
func (c *TypedScopeCreatorTest) TestClassDeclarationWithMethodAndInlineParamDocs() {
	c.TestSame(
		"class Foo {\n" +
			"  method(/** string */ arg) {\n" +
			"    METHOD:;\n" +
			"    var /** number */ foo;\n" +
			"  }\n" +
			"}\n")
	methodBlockScope := c.GetLabeledStatement("METHOD").enclosingScope
	methodScope := methodBlockScope.GetParentScope().(*TypedScope)
	ScopeSubject_AssertScope(methodBlockScope).Declares("foo").Directly().WithTypeThat().IsNumber()
	ScopeSubject_AssertScope(methodScope).Declares("arg").Directly().WithTypeThat().IsString()

	method := typedScopeCreatorTest_cast[rhino.FunctionType](methodScope.GetRootNode().GetJSType())
	TypeSubject_AssertType(method).ToStringIsEqualTo("function(this:Foo, string): undefined")
	foo := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Foo", c.globalScope))
	TypeSubject_AssertType(foo.GetInstanceType()).WithTypeOfProp("method").IsEqualTo(method)
}

// Java: TypedScopeCreatorTest.testClassDeclarationWithOverriddenMethod()
func (c *TypedScopeCreatorTest) TestClassDeclarationWithOverriddenMethod() {
	c.TestSame(
		"class Bar {\n" +
			"  method(/** string */ arg) {}\n" +
			"}\n" +
			"class Foo extends Bar {\n" +
			"  method(arg) {\n" +
			"    METHOD:;\n" +
			"  }\n" +
			"}\n")
	bar := (typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Bar", c.globalScope))).GetInstanceType()
	foo := (typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Foo", c.globalScope))).GetInstanceType()

	methodBlockScope := c.GetLabeledStatement("METHOD").enclosingScope
	methodScope := methodBlockScope.GetParentScope().(*TypedScope)
	ScopeSubject_AssertScope(methodScope).Declares("arg").Directly().WithTypeThat().IsString()
	ScopeSubject_AssertScope(methodScope).Declares("super").WithTypeThat().IsEqualTo(bar)

	method := typedScopeCreatorTest_cast[rhino.FunctionType](methodScope.GetRootNode().GetJSType())
	TypeSubject_AssertType(method).ToStringIsEqualTo("function(this:Foo, string): undefined")
	TypeSubject_AssertType(foo).WithTypeOfProp("method").IsEqualTo(method)
}

// Java: TypedScopeCreatorTest.testClassDeclarationWithStaticMethod()
func (c *TypedScopeCreatorTest) TestClassDeclarationWithStaticMethod() {
	c.TestSame(
		"class Foo {\n" +
			"  static method(/** string */ arg) {\n" +
			"    METHOD:;\n" +
			"    var /** number */ foo;\n" +
			"  }\n" +
			"}\n")
	methodBlockScope := c.GetLabeledStatement("METHOD").enclosingScope
	methodScope := methodBlockScope.GetParentScope().(*TypedScope)
	ScopeSubject_AssertScope(methodBlockScope).Declares("foo").Directly().WithTypeThat().IsNumber()
	ScopeSubject_AssertScope(methodScope).Declares("arg").Directly().WithTypeThat().IsString()

	method := typedScopeCreatorTest_cast[rhino.FunctionType](methodScope.GetRootNode().GetJSType())
	TypeSubject_AssertType(method).ToStringIsEqualTo("function(this:(typeof Foo), string): undefined")
	foo := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Foo", c.globalScope))
	TypeSubject_AssertType(foo).WithTypeOfProp("method").IsEqualTo(method)
}

// Java: TypedScopeCreatorTest.testClassDeclarationWithOverriddenStaticMethod()
func (c *TypedScopeCreatorTest) TestClassDeclarationWithOverriddenStaticMethod() {
	c.TestSame(
		"class Foo {\n" +
			"  /** @param {number} arg */\n" +
			"  static method(arg) {}\n" +
			"}\n" +
			"class Bar extends Foo {\n" +
			"  static method(arg) {\n" +
			"    METHOD: super.method(arg);\n" +
			"  }\n" +
			"}\n")
	fooCtor := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Foo", c.globalScope))
	superMethodCall := c.GetLabeledStatement("METHOD")

	methodOverrideBlockScope := superMethodCall.enclosingScope
	methodOverrideScope := methodOverrideBlockScope.GetParentScope().(*TypedScope)
	ScopeSubject_AssertScope(methodOverrideScope).Declares("arg").Directly().WithTypeThat().IsNumber()
	ScopeSubject_AssertScope(methodOverrideScope).Declares("super").WithTypeThat().IsEqualTo(fooCtor)

	methodOverride := typedScopeCreatorTest_cast[rhino.FunctionType](methodOverrideScope.GetRootNode().GetJSType())
	TypeSubject_AssertType(methodOverride).ToStringIsEqualTo("function(this:(typeof Bar), number): undefined")
	TypeSubject_AssertType(fooCtor).
		WithTypeOfProp("method").
		ToStringIsEqualTo("function(this:(typeof Foo), number): undefined")

	superRef := superMethodCall.statementNode.GetFirstChild().GetFirstFirstChild()
	AssertNode(superRef).HasToken(rhino.Token_SUPER).HasJSTypeThat().IsEqualTo(fooCtor)
}

// Java: TypedScopeCreatorTest.testClassDeclarationWithAsyncMethod()
func (c *TypedScopeCreatorTest) TestClassDeclarationWithAsyncMethod() {
	c.TestSame(
		"class Foo {\n" +
			"  /** @param {string} arg */\n" +
			"  async method(arg) {\n" +
			"    METHOD:;\n" +
			"  }\n" +
			"}\n")
	methodBlockScope := c.GetLabeledStatement("METHOD").enclosingScope
	methodScope := methodBlockScope.GetParentScope().(*TypedScope)
	ScopeSubject_AssertScope(methodScope).Declares("arg").Directly().WithTypeThat().IsString()

	method := typedScopeCreatorTest_cast[rhino.FunctionType](methodScope.GetRootNode().GetJSType())
	TypeSubject_AssertType(method).ToStringIsEqualTo("function(this:Foo, string): Promise<undefined>")
	foo := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Foo", c.globalScope))
	TypeSubject_AssertType(foo.GetInstanceType()).WithTypeOfProp("method").IsEqualTo(method)
}

// Java: TypedScopeCreatorTest.testClassDeclarationWithGeneratorMethod()
func (c *TypedScopeCreatorTest) TestClassDeclarationWithGeneratorMethod() {
	c.TestSame(
		"class Foo {\n" +
			"  /** @param {string} arg */\n" +
			"  * method(arg) {\n" +
			"    METHOD:;\n" +
			"  }\n" +
			"}\n")
	methodBlockScope := c.GetLabeledStatement("METHOD").enclosingScope
	methodScope := methodBlockScope.GetParentScope().(*TypedScope)
	ScopeSubject_AssertScope(methodScope).Declares("arg").Directly().WithTypeThat().IsString()

	method := typedScopeCreatorTest_cast[rhino.FunctionType](methodScope.GetRootNode().GetJSType())
	TypeSubject_AssertType(method).ToStringIsEqualTo("function(this:Foo, string): Generator<?,?,?>")
	foo := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Foo", c.globalScope))
	TypeSubject_AssertType(foo.GetInstanceType()).WithTypeOfProp("method").IsEqualTo(method)
}

// Java: TypedScopeCreatorTest.testClassDeclarationWithComputedPropertyMethod()
func (c *TypedScopeCreatorTest) TestClassDeclarationWithComputedPropertyMethod() {
	c.TestSame(
		"class Foo {\n" +
			"  /** @param {string} arg */\n" +
			"  ['method'](arg) {\n" +
			"    METHOD:;\n" +
			"    var /** number */ foo;\n" +
			"  }\n" +
			"}\n")
	methodBlockScope := c.GetLabeledStatement("METHOD").enclosingScope
	methodScope := methodBlockScope.GetParentScope().(*TypedScope)
	ScopeSubject_AssertScope(methodBlockScope).Declares("foo").Directly().WithTypeThat().IsNumber()
	ScopeSubject_AssertScope(methodScope).Declares("arg").Directly().WithTypeThat().IsString()

	method := typedScopeCreatorTest_cast[rhino.FunctionType](methodScope.GetRootNode().GetJSType())
	TypeSubject_AssertType(method).ToStringIsEqualTo("function(this:Foo, string): undefined")
}

// Java: TypedScopeCreatorTest.testClassDeclarationWithField()
func (c *TypedScopeCreatorTest) TestClassDeclarationWithField() {
	c.TestSame(
		"class Foo {\n" +
			"  /** @type {string} */\n" +
			"  a = 'hi';\n" +
			"  /** @type {symbol|undefined} */\n" +
			"  b;\n" +
			"  c = 0;\n" +
			"  /** @type {string|null} */\n" +
			"  d = null;\n" +
			"}\n")

	fooCtor := c.globalScope.GetVar("Foo").GetType().AssertFunctionType()
	fooInstance := fooCtor.GetInstanceType()

	AssertThat(fooInstance.HasOwnDeclaredProperty("a")).IsTrue()
	TypeSubject_AssertType(fooInstance.GetPropertyType("a")).IsString()

	AssertThat(fooInstance.HasOwnDeclaredProperty("b")).IsTrue()
	TypeSubject_AssertType(fooInstance.GetPropertyType("b")).ToStringIsEqualTo("(symbol|undefined)")

	AssertThat(fooInstance.HasOwnDeclaredProperty("c")).IsTrue()
	TypeSubject_AssertType(fooInstance.GetPropertyType("c")).IsNumber()

	AssertThat(fooInstance.HasOwnDeclaredProperty("d")).IsTrue()
	TypeSubject_AssertType(fooInstance.GetPropertyType("d")).ToStringIsEqualTo("(null|string)")
}

// Java: TypedScopeCreatorTest.testClassDeclarationWithFieldThisAndSuper()
func (c *TypedScopeCreatorTest) TestClassDeclarationWithFieldThisAndSuper() {
	c.TestSame(
		"class Foo {\n" +
			"  /** @type {number} */\n" +
			"  a = 2;\n" +
			"  /** @type {number} */\n" +
			"  b = this.a;\n" +
			"}\n" +
			"class Bar extends Foo {\n" +
			"  /** @type {number} */\n" +
			"  c = this.a + 1; // Must use `this` and not `super` to access field of parent class\n" +
			"}\n")

	fooCtor := c.globalScope.GetVar("Foo").GetType().AssertFunctionType()
	fooInstance := fooCtor.GetInstanceType()
	barCtor := c.globalScope.GetVar("Bar").GetType().AssertFunctionType()
	barInstance := barCtor.GetInstanceType()

	AssertThat(fooInstance.HasOwnDeclaredProperty("a")).IsTrue()
	TypeSubject_AssertType(fooInstance.GetPropertyType("a")).IsNumber()
	AssertThat(barInstance.HasOwnDeclaredProperty("a")).IsFalse()

	AssertThat(fooInstance.HasOwnDeclaredProperty("b")).IsTrue()
	TypeSubject_AssertType(fooInstance.GetPropertyType("b")).IsNumber()
	AssertThat(barInstance.HasOwnDeclaredProperty("b")).IsFalse()

	AssertThat(fooInstance.HasOwnDeclaredProperty("c")).IsFalse()
	AssertThat(barInstance.HasOwnDeclaredProperty("c")).IsTrue()
	TypeSubject_AssertType(barInstance.GetPropertyType("c")).IsNumber()
}

// When there is no JSDoc type declaration for a field, and the RHS is uninferrable or unknown,
// the field's type is set to the ALL type ('*')
// Java: TypedScopeCreatorTest.testPublicFieldWithUnknownTypeCast()
func (c *TypedScopeCreatorTest) TestPublicFieldWithUnknownTypeCast() {
	c.TestSame(
		"class Bar {\n" +
			"  x = /** @type {?} */ ({});\n" +
			"}\n")
	barCtor := c.globalScope.GetVar("Bar").GetType().AssertFunctionType()
	barInstance := barCtor.GetInstanceType()
	AssertThat(barInstance.HasOwnDeclaredProperty("x")).IsTrue()
	TypeSubject_AssertType(barInstance.GetPropertyType("x")).ToStringIsEqualTo("*")

	c.TestSame("class Foo { x; }")
	fooCtor := c.globalScope.GetVar("Foo").GetType().AssertFunctionType()
	fooInstance := fooCtor.GetInstanceType()
	AssertThat(fooInstance.HasOwnDeclaredProperty("x")).IsTrue()
	TypeSubject_AssertType(fooInstance.GetPropertyType("x")).ToStringIsEqualTo("*")
}

// Java: TypedScopeCreatorTest.testStaticFieldWithUnknownTypeCast()
func (c *TypedScopeCreatorTest) TestStaticFieldWithUnknownTypeCast() {
	c.TestSame("class Drum { static y = /** @type {?} */ (5); }")
	drumCtor := c.globalScope.GetVar("Drum").GetType().AssertFunctionType()
	AssertThat(drumCtor.HasOwnDeclaredProperty("y")).IsTrue()
	TypeSubject_AssertType(drumCtor.GetPropertyType("y")).ToStringIsEqualTo("*")
}

// Java: TypedScopeCreatorTest.testClassDeclarationWithStaticField()
func (c *TypedScopeCreatorTest) TestClassDeclarationWithStaticField() {
	c.TestSame(
		"class Foo {\n" +
			"  /** @type {string} */\n" +
			"  static a = 'hi';\n" +
			"  /** @type {symbol|undefined} */\n" +
			"  static b;\n" +
			"  static c = 0;\n" +
			"  static d;\n" +
			"  /** @type {string|null} */\n" +
			"  static e = null;\n" +
			"}\n")

	fooCtor := c.globalScope.GetVar("Foo").GetType().AssertFunctionType()

	AssertThat(fooCtor.HasOwnDeclaredProperty("a")).IsTrue()
	TypeSubject_AssertType(fooCtor.GetPropertyType("a")).IsString()

	AssertThat(fooCtor.HasOwnDeclaredProperty("b")).IsTrue()
	TypeSubject_AssertType(fooCtor.GetPropertyType("b")).ToStringIsEqualTo("(symbol|undefined)")

	AssertThat(fooCtor.HasOwnDeclaredProperty("c")).IsTrue()
	TypeSubject_AssertType(fooCtor.GetPropertyType("c")).IsNumber()

	AssertThat(fooCtor.HasOwnDeclaredProperty("d")).IsTrue()
	TypeSubject_AssertType(fooCtor.GetPropertyType("d")).ToStringIsEqualTo("*")

	AssertThat(fooCtor.HasOwnDeclaredProperty("e")).IsTrue()
	TypeSubject_AssertType(fooCtor.GetPropertyType("e")).ToStringIsEqualTo("(null|string)")
}

// Java: TypedScopeCreatorTest.testClassDeclarationWithStaticFieldThisAndSuper()
func (c *TypedScopeCreatorTest) TestClassDeclarationWithStaticFieldThisAndSuper() {
	c.TestSame(
		"class Foo {\n" +
			"  /** @type {number} */\n" +
			"  static a = 2;\n" +
			"  /** @type {number} */\n" +
			"  static b = this.a;\n" +
			"}\n" +
			"class Bar extends Foo {\n" +
			"  /** @type {number} */\n" +
			"  static c = this.a + 1; // Must use `this` and not `super` to access field of parent\n" +
			"// class\n" +
			"}\n")

	fooCtor := c.globalScope.GetVar("Foo").GetType().AssertFunctionType()
	barCtor := c.globalScope.GetVar("Bar").GetType().AssertFunctionType()

	AssertThat(fooCtor.HasOwnDeclaredProperty("a")).IsTrue()
	TypeSubject_AssertType(fooCtor.GetPropertyType("a")).IsNumber()
	AssertThat(barCtor.HasOwnDeclaredProperty("a")).IsFalse()

	AssertThat(fooCtor.HasOwnDeclaredProperty("b")).IsTrue()
	TypeSubject_AssertType(fooCtor.GetPropertyType("b")).IsNumber()
	AssertThat(barCtor.HasOwnDeclaredProperty("b")).IsFalse()

	AssertThat(fooCtor.HasOwnDeclaredProperty("c")).IsFalse()
	AssertThat(barCtor.HasOwnDeclaredProperty("c")).IsTrue()
	TypeSubject_AssertType(barCtor.GetPropertyType("c")).IsNumber()
}

// Java: TypedScopeCreatorTest.testTypeOfThisInObjectLitInPublicField()
func (c *TypedScopeCreatorTest) TestTypeOfThisInObjectLitInPublicField() {
	// Make sure the `this` on the RHS of a class field declaration has the appropriate type, even
	// when it appears within an object literal
	c.TestSame(
		"class Foo {\n" +
			"  x = 3;\n" +
			"  bar = {\n" +
			"    b: this.x\n" +
			"  };\n" +
			"}\n")

	fooCtor := c.globalScope.GetVar("Foo").GetType().AssertFunctionType()
	fooInstance := fooCtor.GetInstanceType()
	getProp :=
		fooInstance.
			GetPropertyNode("bar").
			GetOnlyChild(). // OBJECTLIT
			GetOnlyChild(). // STRING_KEY
			GetOnlyChild()  // GETPROP
	AssertNode(getProp).HasJSTypeThat().IsNumber()
	thisNode := getProp.GetOnlyChild()
	AssertNode(thisNode).HasJSTypeThat().IsEqualTo(fooInstance)

	c.TestSame(
		"class Foo {\n" +
			"  static baz = {\n" +
			"    e: this\n" +
			"  };\n" +
			"}\n")

	fooCtor = c.globalScope.GetVar("Foo").GetType().AssertFunctionType()
	thisNode =
		fooCtor.
			GetPropertyNode("baz").
			GetOnlyChild(). // OBJECTLIT
			GetOnlyChild(). // STRING_KEY
			GetOnlyChild()  // THIS
	AssertNode(thisNode).HasJSTypeThat().IsEqualTo(fooCtor)
}

// Java: TypedScopeCreatorTest.testGetterInObjectLitInPublicField()
func (c *TypedScopeCreatorTest) TestGetterInObjectLitInPublicField() {
	// Previously a bug caused getters & setters in object literals on the RHS of a field
	// declaration to be treated as if they belonged to the class.
	//
	// Make sure that `this` in an object literal getter or setter has the unknown type.
	c.TestSame(
		"class Foo {\n" +
			"  bar = {\n" +
			"    get baz() {\n" +
			"      X: this;\n" +
			"    }\n" +
			"  };\n" +
			"}\n")
	thisNode := c.GetLabeledStatement("X").statementNode.GetOnlyChild()
	// We don't create unique types for object literals, so even though we know
	// the getter will (probably) only be called on this object, we just
	// use unknown for the type of `this`.
	AssertNode(thisNode).HasJSTypeThat().IsUnknown()

	c.TestSame(
		"class Foo {\n" +
			"  static bar = {\n" +
			"    /** @this {string} */\n" +
			"    get baz() { Y: this; }\n" +
			"  };\n" +
			"}\n")

	thisNode = c.GetLabeledStatement("Y").statementNode.GetOnlyChild()
	AssertNode(thisNode).HasJSTypeThat().IsString()
}

// Java: TypedScopeCreatorTest.testMemberFunctionInObjectLitInPublicField()
func (c *TypedScopeCreatorTest) TestMemberFunctionInObjectLitInPublicField() {
	// Previously a bug caused methods in object literals on the RHS of a field declaration to be
	// treated as if they were methods on the class.
	//
	// Make sure that `this` in an object literal method has the unknown type.
	c.TestSame(
		"class Foo {\n" +
			"  bar = {\n" +
			"    /** @this {number} */\n" +
			"    member() {\n" +
			"      X: this;\n" +
			"    }\n" +
			"  };\n" +
			"}\n")

	thisNode := c.GetLabeledStatement("X").statementNode.GetOnlyChild()
	AssertNode(thisNode).HasJSTypeThat().IsNumber()
}

// Java: TypedScopeCreatorTest.testClassDeclarationWithDeprecatedField()
func (c *TypedScopeCreatorTest) TestClassDeclarationWithDeprecatedField() {
	c.TestSame(
		"class Foo {\n" +
			"  /** @deprecated There is a madness to this method */\n" +
			"  a = 'hi';\n" +
			"}\n")
}

// Java: TypedScopeCreatorTest.testClassDeclarationWithDeprecatedComputedField()
func (c *TypedScopeCreatorTest) TestClassDeclarationWithDeprecatedComputedField() {
	c.TestSame(
		"class Foo {\n" +
			"  /** @deprecated There is a madness to this method */\n" +
			"  ['a'] = 'hi';\n" +
			"}\n")
}

// Java: TypedScopeCreatorTest.testClassExpressionAssignment()
func (c *TypedScopeCreatorTest) TestClassExpressionAssignment() {
	c.TestSame("var Foo = class Bar {}")
	foo := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Foo", c.globalScope))
	AssertThat(foo.IsConstructor()).IsTrue()
	bar := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Bar", c.globalScope))
	AssertThat(bar).IsEqualTo(foo)
}

// Java: TypedScopeCreatorTest.testClassExpressionBleedingNameScope()
func (c *TypedScopeCreatorTest) TestClassExpressionBleedingNameScope() {
	c.TestSame(
		"var Foo = class Bar {\n" +
			"  constructor() {\n" +
			"    CTOR:;\n" +
			"  }\n" +
			"};\n")
	foo := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Foo", c.globalScope))
	ctorBlockScope := c.GetLabeledStatement("CTOR").enclosingScope
	ctorScope := ctorBlockScope.GetParentScope().(*TypedScope)
	classScope := ctorScope.GetParentScope().(*TypedScope)
	ScopeSubject_AssertScope(c.globalScope).Declares("Foo").WithTypeThat().IsEqualTo(foo)
	ScopeSubject_AssertScope(classScope).Declares("Bar").Directly().WithTypeThat().IsEqualTo(foo)
	ScopeSubject_AssertScope(c.globalScope).DoesNotDeclare("Bar")
}

// Java: TypedScopeCreatorTest.testClassExpression_nestedAssignmentToName()
func (c *TypedScopeCreatorTest) TestClassExpression_nestedAssignmentToName() {
	c.TestSame(
		c.Externs("function somethingElse() {}"),
		// test a pattern that tsickle generates for decorated classes
		c.Srcs(
			"var Foo_1;\n"+
				"var Foo_2;\n"+
				"var Foo = Foo_1 = Foo_2 = class Bar {};\n"+
				"Foo = /** @type {typeof Foo} */ (somethingElse());\n"+
				"/** @type {!Foo} */ var f;\n"))

	foo := c.FindNameType("Foo", c.globalScope)
	bar := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Bar", c.globalScope))
	TypedVarSubject_AssertThat(c.globalScope.GetVar("Foo")).HasJSTypeThat().IsEqualTo(bar)
	TypeSubject_AssertType(foo).IsEqualTo(bar)
	AssertThat(c.FindNameType("f", c.globalScope)).IsEqualTo(bar.GetInstanceType())
}

// Java: TypedScopeCreatorTest.testClassExpression_nestedAssignmentToConstQName()
func (c *TypedScopeCreatorTest) TestClassExpression_nestedAssignmentToConstQName() {
	c.TestSame(
		c.Srcs(
			"var Foo_1;\n" +
				"const ns = {};\n" +
				"/** @const */\n" +
				"ns.Foo = Foo_1 = class Bar {};\n" +
				"/** @type {!ns.Foo} */ var f;\n"))

	bar := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Bar", c.globalScope))
	AssertThat(c.FindNameType("f", c.globalScope)).IsEqualTo(bar.GetInstanceType())
}

// Java: TypedScopeCreatorTest.testClassExpression_nestedAssignmentToNonConstQName()
func (c *TypedScopeCreatorTest) TestClassExpression_nestedAssignmentToNonConstQName() {
	c.TestSame(
		c.Externs("function somethingElse() {}"),
		// test a pattern that tsickle generates for decorated classes
		c.Srcs(
			"var Foo_1;\n"+
				"const ns = {};\n"+
				"ns.Foo = Foo_1 = class Bar {};\n"+
				"ns.Foo = Foo_1 = /** @type {typeof ns.Foo} */ (somethingElse());\n"+
				"/** @type {!ns.Foo} */ var f;\n"))

	bar := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Bar", c.globalScope))
	AssertThat(c.FindNameType("f", c.globalScope)).IsEqualTo(bar.GetInstanceType())
}

// Java: TypedScopeCreatorTest.testClassExpressionWithMethod()
func (c *TypedScopeCreatorTest) TestClassExpressionWithMethod() {
	c.TestSame(
		"var Foo = class {\n" +
			"  /** @param {string} arg */\n" +
			"  method(arg) {\n" +
			"    METHOD:;\n" +
			"  }\n" +
			"};\n")
	methodBlockScope := c.GetLabeledStatement("METHOD").enclosingScope
	methodScope := methodBlockScope.GetParentScope().(*TypedScope)
	ScopeSubject_AssertScope(methodScope).Declares("arg").Directly().WithTypeThat().IsString()

	foo := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Foo", c.globalScope))
	TypeSubject_AssertType(foo.GetInstanceType()).
		WithTypeOfProp("method").
		ToStringIsEqualTo("function(this:Foo, string): undefined")
}

// Java: TypedScopeCreatorTest.testClassExpressionWithStaticMethod()
func (c *TypedScopeCreatorTest) TestClassExpressionWithStaticMethod() {
	c.TestSame(
		"var Foo = class {\n" +
			"  static method(/** string */ arg) {\n" +
			"    METHOD:;\n" +
			"    var /** number */ foo;\n" +
			"  }\n" +
			"}\n")
	methodBlockScope := c.GetLabeledStatement("METHOD").enclosingScope
	methodScope := methodBlockScope.GetParentScope().(*TypedScope)
	ScopeSubject_AssertScope(methodBlockScope).Declares("foo").Directly().WithTypeThat().IsNumber()
	ScopeSubject_AssertScope(methodScope).Declares("arg").Directly().WithTypeThat().IsString()

	method := typedScopeCreatorTest_cast[rhino.FunctionType](methodScope.GetRootNode().GetJSType())
	TypeSubject_AssertType(method).ToStringIsEqualTo("function(this:(typeof Foo), string): undefined")
	foo := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Foo", c.globalScope))
	TypeSubject_AssertType(foo).WithTypeOfProp("method").IsEqualTo(method)
}

// Java: TypedScopeCreatorTest.testClassExpressionWithStaticClassAssignedLater()
func (c *TypedScopeCreatorTest) TestClassExpressionWithStaticClassAssignedLater() {
	c.TestSame(
		"var Foo = class {};\n" +
			"// even though there's no JSDoc, this should be treated as a type declaration.\n" +
			"Foo.Something = class {};\n")
	ScopeSubject_AssertScope(c.globalScope).
		Declares("Foo").
		Directly().
		WithTypeThat().
		ToStringIsEqualTo("(typeof Foo)")
	ScopeSubject_AssertScope(c.globalScope).
		Declares("Foo.Something").
		Directly().
		WithTypeThat().
		ToStringIsEqualTo("(typeof Foo.Something)")
}

// Java: TypedScopeCreatorTest.testClassExpressionInCallback()
func (c *TypedScopeCreatorTest) TestClassExpressionInCallback() {
	c.TestSame(
		"function use(arg) {}\n" +
			"use(class Bar {\n" +
			"  constructor() {\n" +
			"    CTOR:;\n" +
			"  }\n" +
			"});\n")
	ctorBlockScope := c.GetLabeledStatement("CTOR").enclosingScope
	ctorScope := ctorBlockScope.GetParentScope().(*TypedScope)
	classScope := ctorScope.GetParentScope().(*TypedScope)
	ScopeSubject_AssertScope(classScope).
		Declares("Bar").
		Directly().
		WithTypeThat().
		// TODO(sdh): Print a better name (https://github.com/google/closure-compiler/issues/2982)
		ToStringIsEqualTo("(typeof <anonymous@testcode:2>)")
	ScopeSubject_AssertScope(c.globalScope).DoesNotDeclare("Bar")
}

// Java: TypedScopeCreatorTest.testClassStaticBlockAndVariableDeclarations()
func (c *TypedScopeCreatorTest) TestClassStaticBlockAndVariableDeclarations() {
	c.TestSame(
		"class Foo {\n" +
			"  static {\n" +
			"    STATIC:;\n" +
			"    let /** string */ x;\n" +
			"    const /** number */ y = 5;\n" +
			"    var /** boolean */ z;\n" +
			"  }\n" +
			"};\n")
	staticScope := c.GetLabeledStatement("STATIC").enclosingScope
	classScope := staticScope.GetParentScope().(*TypedScope)
	ScopeSubject_AssertScope(c.globalScope).
		Declares("Foo").
		Directly().
		WithTypeThat().
		ToStringIsEqualTo("(typeof Foo)")
	ScopeSubject_AssertScope(staticScope).Declares("x").Directly().WithTypeThat().IsString()
	ScopeSubject_AssertScope(classScope).DoesNotDeclare("x")
	ScopeSubject_AssertScope(c.globalScope).DoesNotDeclare("x")
	ScopeSubject_AssertScope(staticScope).Declares("y").Directly().WithTypeThat().IsNumber()
	ScopeSubject_AssertScope(classScope).DoesNotDeclare("y")
	ScopeSubject_AssertScope(c.globalScope).DoesNotDeclare("y")
	ScopeSubject_AssertScope(staticScope).Declares("z").Directly().WithTypeThat().IsBoolean()
	ScopeSubject_AssertScope(classScope).DoesNotDeclare("z")
	ScopeSubject_AssertScope(c.globalScope).DoesNotDeclare("z")
}

// Java: TypedScopeCreatorTest.testMultipleClassStaticBlocks()
func (c *TypedScopeCreatorTest) TestMultipleClassStaticBlocks() {
	c.TestSame(
		"class Foo {\n" +
			"  static {\n" +
			"    STATIC1:;\n" +
			"    var /** number */ x;\n" +
			"  }\n" +
			"  static {\n" +
			"    STATIC2:;\n" +
			"    var /** string */ y;\n" +
			"  }\n" +
			"};\n")
	staticScope1 := c.GetLabeledStatement("STATIC1").enclosingScope
	staticScope2 := c.GetLabeledStatement("STATIC2").enclosingScope
	classScope := staticScope1.GetParentScope().(*TypedScope)
	ScopeSubject_AssertScope(c.globalScope).
		Declares("Foo").
		Directly().
		WithTypeThat().
		ToStringIsEqualTo("(typeof Foo)")
	AssertThat(staticScope2.GetParentScope().(*TypedScope)).IsEqualTo(classScope)
	ScopeSubject_AssertScope(staticScope1).Declares("x").Directly().WithTypeThat().IsNumber()
	ScopeSubject_AssertScope(staticScope2).Declares("y").Directly().WithTypeThat().IsString()
	ScopeSubject_AssertScope(staticScope1).DoesNotDeclare("y")
	ScopeSubject_AssertScope(staticScope2).DoesNotDeclare("x")
	ScopeSubject_AssertScope(classScope).DoesNotDeclare("x")
	ScopeSubject_AssertScope(classScope).DoesNotDeclare("y")
	ScopeSubject_AssertScope(c.globalScope).DoesNotDeclare("x")
	ScopeSubject_AssertScope(c.globalScope).DoesNotDeclare("y")
}

// Java: TypedScopeCreatorTest.testClassStaticBlockWithIf()
func (c *TypedScopeCreatorTest) TestClassStaticBlockWithIf() {
	c.TestSame(
		"class Foo {\n" +
			"  static {\n" +
			"    STATIC:;\n" +
			"    if(true) {\n" +
			"      IF:;\n" +
			"      let /** boolean */ x;\n" +
			"      var /** string */ y;\n" +
			"      const /** number */ z = 1;\n" +
			"    }\n" +
			"  }\n" +
			"};\n")
	ifScope := c.GetLabeledStatement("IF").enclosingScope
	staticScope := c.GetLabeledStatement("STATIC").enclosingScope
	classScope := staticScope.GetParentScope().(*TypedScope)
	ScopeSubject_AssertScope(c.globalScope).
		Declares("Foo").
		Directly().
		WithTypeThat().
		ToStringIsEqualTo("(typeof Foo)")
	ScopeSubject_AssertScope(ifScope).Declares("x").Directly().WithTypeThat().IsBoolean()
	ScopeSubject_AssertScope(staticScope).DoesNotDeclare("x")
	ScopeSubject_AssertScope(classScope).DoesNotDeclare("x")
	ScopeSubject_AssertScope(c.globalScope).DoesNotDeclare("x")
	ScopeSubject_AssertScope(ifScope).Declares("y").OnSomeParent().WithTypeThat().IsString()
	ScopeSubject_AssertScope(staticScope).Declares("y").Directly().WithTypeThat().IsString()
	ScopeSubject_AssertScope(classScope).DoesNotDeclare("y")
	ScopeSubject_AssertScope(c.globalScope).DoesNotDeclare("y")
	ScopeSubject_AssertScope(ifScope).Declares("z").Directly().WithTypeThat().IsNumber()
	ScopeSubject_AssertScope(staticScope).DoesNotDeclare("z")
	ScopeSubject_AssertScope(classScope).DoesNotDeclare("z")
	ScopeSubject_AssertScope(c.globalScope).DoesNotDeclare("z")
}

// Java: TypedScopeCreatorTest.testClassStaticBlockWithFunction()
func (c *TypedScopeCreatorTest) TestClassStaticBlockWithFunction() {
	c.TestSame(
		"class Foo {\n" +
			"  static {\n" +
			"    STATIC:;\n" +
			"    function bar() {\n" +
			"      FUNCTION:;\n" +
			"    }\n" +
			"  }\n" +
			"};\n")
	functionScope := c.GetLabeledStatement("FUNCTION").enclosingScope
	staticScope := c.GetLabeledStatement("STATIC").enclosingScope
	classScope := staticScope.GetParentScope().(*TypedScope)
	AssertThat(functionScope.GetParentScope().(*TypedScope).GetParentScope().(*TypedScope)).IsEqualTo(staticScope)
	ScopeSubject_AssertScope(c.globalScope).
		Declares("Foo").
		Directly().
		WithTypeThat().
		ToStringIsEqualTo("(typeof Foo)")
	ScopeSubject_AssertScope(staticScope).
		Declares("bar").
		Directly().
		WithTypeThat().
		ToStringIsEqualTo("function(): undefined")
	ScopeSubject_AssertScope(classScope).DoesNotDeclare("bar")
	ScopeSubject_AssertScope(c.globalScope).DoesNotDeclare("bar")
}

// Java: TypedScopeCreatorTest.testClassStaticBlockWithThisLHS()
func (c *TypedScopeCreatorTest) TestClassStaticBlockWithThisLHS() {
	c.TestSame(
		"class Foo {\n" +
			"  static {\n" +
			"    STATIC:;\n" +
			"    /** @type {string} */\n" +
			"    this.x = 'hello';\n" +
			"  }\n" +
			"};\n")
	fooCtor := c.globalScope.GetVar("Foo").GetType().AssertFunctionType()
	staticScope := c.GetLabeledStatement("STATIC").enclosingScope
	classScope := staticScope.GetParentScope().(*TypedScope)
	AssertThat(fooCtor.HasProperty("x")).IsTrue()
	TypeSubject_AssertType(fooCtor.GetPropertyType("x")).IsString()
	AssertThat(fooCtor.IsPropertyTypeDeclared("x")).IsTrue()
	AssertThat(fooCtor.IsPropertyTypeInferred("x")).IsFalse()
	ScopeSubject_AssertScope(staticScope).DoesNotDeclare("x")
	ScopeSubject_AssertScope(classScope).DoesNotDeclare("x")
}

// Java: TypedScopeCreatorTest.testClassStaticBlockWithThisRHS()
func (c *TypedScopeCreatorTest) TestClassStaticBlockWithThisRHS() {
	c.TestSame(
		"class Foo {\n" +
			"  /** @type {number} */\n" +
			"  static y = 1;\n" +
			"  static {\n" +
			"    STATIC:;\n" +
			"    var x = this.y;\n" +
			"  }\n" +
			"};\n")
	fooCtor := c.globalScope.GetVar("Foo").GetType().AssertFunctionType()
	staticScope := c.GetLabeledStatement("STATIC").enclosingScope
	classScope := staticScope.GetParentScope().(*TypedScope)
	AssertThat(fooCtor.HasOwnDeclaredProperty("y")).IsTrue()
	TypeSubject_AssertType(fooCtor.GetPropertyType("y")).IsNumber()
	AssertThat(fooCtor.HasProperty("x")).IsFalse()
	ScopeSubject_AssertScope(staticScope).Declares("x").Directly().WithTypeThat().IsNumber()
	ScopeSubject_AssertScope(classScope).DoesNotDeclare("x")
}

// Java: TypedScopeCreatorTest.testSuperInEmptyClassStaticBlock()
func (c *TypedScopeCreatorTest) TestSuperInEmptyClassStaticBlock() {
	c.TestSame(
		"class Foo {}\n" +
			"class Bar extends Foo {\n" +
			"  static {\n" +
			"    STATIC:;\n" +
			"  }\n" +
			"};\n")
	fooCtor := c.globalScope.GetVar("Foo").GetType().AssertFunctionType()
	staticScope := c.GetLabeledStatement("STATIC").enclosingScope
	classBarScope := staticScope.GetParentScope().(*TypedScope)
	ScopeSubject_AssertScope(staticScope).Declares("super").Directly().WithTypeThat().IsEqualTo(fooCtor)
	ScopeSubject_AssertScope(classBarScope).DoesNotDeclare("super")
}

// Java: TypedScopeCreatorTest.testForLoopIntegration()
func (c *TypedScopeCreatorTest) TestForLoopIntegration() {
	c.TestSame("var y = 3; for (var x = true; x; y = x) {}")

	y := c.globalScope.GetVar("y")
	AssertThat(y.IsTypeInferred()).IsTrue()
	AssertThat(y.GetType().String()).IsEqualTo("(boolean|number)")
}

// Java: TypedScopeCreatorTest.testConstructorAlias()
func (c *TypedScopeCreatorTest) TestConstructorAlias() {
	c.TestSame(
		"/** @constructor */ var Foo = function() {};\n" +
			"/** @constructor */ var FooAlias = Foo;\n")
	AssertThat(c.registry.GetType(nil, "FooAlias").String()).IsEqualTo("Foo")
	TypeSubject_AssertType(c.registry.GetType(nil, "FooAlias")).IsEqualTo(c.registry.GetType(nil, "Foo"))
}

// Java: TypedScopeCreatorTest.testNamespacedConstructorAlias()
func (c *TypedScopeCreatorTest) TestNamespacedConstructorAlias() {
	c.TestSame(
		"/** @constructor */ goog.Foo = function() {};\n" +
			"/** @constructor */ goog.FooAlias = goog.Foo;\n")
	AssertThat(c.registry.GetType(nil, "goog.FooAlias").String()).IsEqualTo("goog.Foo")
	TypeSubject_AssertType(c.registry.GetType(nil, "goog.FooAlias")).
		IsEqualTo(c.registry.GetType(nil, "goog.Foo"))
}

// Java: TypedScopeCreatorTest.testConstructorAliasWithPrototypeMethodHasCorrectThisType()
func (c *TypedScopeCreatorTest) TestConstructorAliasWithPrototypeMethodHasCorrectThisType() {
	c.TestSame(
		"class Foo {}\n" +
			"const FooAlias = Foo;\n" +
			"FooAlias.prototype.fn = function() {};\n")

	fnAliasNameNode :=
		c.FindQualifiedNameNode("FooAlias.prototype.fn", c.globalScope.GetRootNode())

	fooCtorType := c.FindNameType("Foo", c.globalScope)
	fooType := fooCtorType.ToMaybeFunctionType().GetInstanceType()

	TypeSubject_AssertType(fnAliasNameNode.GetJSType()).
		IsFunctionTypeThat().
		HasTypeOfThisThat().
		IsEqualTo(fooType)
}

// Java: TypedScopeCreatorTest.testConstructorAliasPrototypePropIsInScope()
func (c *TypedScopeCreatorTest) TestConstructorAliasPrototypePropIsInScope() {
	c.TestSame(
		"class Foo {}\n" +
			"const FooAlias = Foo;\n" +
			"FooAlias.prototype.fn = function() {};\n")

	fooAliasPrototype := typedVarOf(c.globalScope.GetSlot("FooAlias.prototype"))
	TypedVarSubject_AssertThat(fooAliasPrototype).IsNotNull()
	TypeSubject_AssertType(fooAliasPrototype.GetType()).
		IsEqualTo(typedVarOf(c.globalScope.GetSlot("Foo.prototype")).GetType())
}

// Java: TypedScopeCreatorTest.testNoConstructorAlias_namespaceIsShadowedInParameter()
func (c *TypedScopeCreatorTest) TestNoConstructorAlias_namespaceIsShadowedInParameter() {
	c.TestWarning(
		"const a = {};\n"+
			"a.B = class {};\n"+
			"function f(a) {\n"+
			"  const B = a.B; // This does /not/ create a constructor alias as `a` is shadowed.\n"+
			"  var /** !B */ b;\n"+
			"}\n",
		RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR)
}

// Java: TypedScopeCreatorTest.testNoTypeResolution_namespaceShadowedAfterReference()
func (c *TypedScopeCreatorTest) TestNoTypeResolution_namespaceShadowedAfterReference() {
	c.TestWarning(
		"const a = {};\n"+
			"a.B = class {};\n"+
			"function f() {\n"+
			"  var /** !a.B */ b; // Treat 'a' as coming from 'let a;' below.\n"+
			"  let a;\n"+
			"}\n",
		RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR)
}

// Java: TypedScopeCreatorTest.testNoTypeResolution_namespaceShadowedOutsideScopeAfterReference()
func (c *TypedScopeCreatorTest) TestNoTypeResolution_namespaceShadowedOutsideScopeAfterReference() {
	c.TestWarning(
		"const a = {};\n"+
			"a.B = class {};\n"+
			"function f() {\n"+
			"  { var /** !a.B */ b; } // Treat 'a' as coming from 'let a;' below.\n"+
			"  let a;\n"+
			"}\n",
		RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR)
}

// Java: TypedScopeCreatorTest.testNoTypeResolution_namespaceShadowedAfterTemplateScopeReference()
func (c *TypedScopeCreatorTest) TestNoTypeResolution_namespaceShadowedAfterTemplateScopeReference() {
	c.TestWarning(
		"const a = {};\n"+
			"a.B = class {};\n"+
			"function f() {\n"+
			"  /**\n"+
			"   * @param {(!a.B|T)} b // Treat 'a' as coming from 'let a;' below.\n"+
			"   * @template T\n"+
			"   */\n"+
			"  function f(b) {}\n"+
			"  let a;\n"+
			"}\n",
		RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR)
}

// Java: TypedScopeCreatorTest.testTemplateType1()
func (c *TypedScopeCreatorTest) TestTemplateType1() {
	c.TestSame(
		"/**\n" +
			" * @param {function(this:T, ...)} fn\n" +
			" * @param {T} thisObj\n" +
			" * @template T\n" +
			" */\n" +
			"function bind(fn, thisObj) {}\n" +
			"/** @constructor */\n" +
			"function Foo() {}\n" +
			"/** @return {number} */\n" +
			"Foo.prototype.baz = function() {};\n" +
			"bind(function() { var g = this; var f = this.baz(); }, new Foo());\n")
	AssertThat(c.FindNameType("g", c.lastLocalScope).String()).IsEqualTo("Foo")
	AssertThat(c.FindNameType("f", c.lastLocalScope).String()).IsEqualTo("number")
}

// Java: TypedScopeCreatorTest.testTemplateType2()
func (c *TypedScopeCreatorTest) TestTemplateType2() {
	c.TestSame(
		"/**\n" +
			" * @param {T} x\n" +
			" * @return {T}\n" +
			" * @template T\n" +
			" */\n" +
			"function f(x) {\n" +
			"  return x;\n" +
			"}\n" +
			"/** @type {string} */\n" +
			"var val = 'hi';\n" +
			"var result = f(val);\n")
	AssertThat(c.FindNameType("result", c.globalScope).String()).IsEqualTo("string")
}

// Java: TypedScopeCreatorTest.testTemplateType2a()
func (c *TypedScopeCreatorTest) TestTemplateType2a() {
	c.TestSame(
		"/**\n" +
			" * @param {T} x\n" +
			" * @return {T|undefined}\n" +
			" * @template T\n" +
			" */\n" +
			"function f(x) {\n" +
			"  return x;\n" +
			"}\n" +
			"/** @type {string} */\n" +
			"var val = 'hi';\n" +
			"var result = f(val);\n")
	AssertThat(c.FindNameType("result", c.globalScope).String()).IsEqualTo("(string|undefined)")
}

// Java: TypedScopeCreatorTest.testTemplateType2b()
func (c *TypedScopeCreatorTest) TestTemplateType2b() {
	c.TestSame(
		"/**\n" +
			" * @param {T} x\n" +
			" * @return {T}\n" +
			" * @template T\n" +
			" */\n" +
			"function f(x) {\n" +
			"  return x;\n" +
			"}\n" +
			"/** @type {string|undefined} */\n" +
			"var val = 'hi';\n" +
			"var result = f(val);\n")
	AssertThat(c.FindNameType("result", c.globalScope).String()).IsEqualTo("(string|undefined)")
}

// Java: TypedScopeCreatorTest.testTemplateType3()
func (c *TypedScopeCreatorTest) TestTemplateType3() {
	c.TestSame(
		"/**\n" +
			" * @param {T} x\n" +
			" * @return {T}\n" +
			" * @template T\n" +
			" */\n" +
			"function f(x) {\n" +
			"  return x;\n" +
			"}\n" +
			"/** @type {string} */\n" +
			"var val1 = 'hi';\n" +
			"var result1 = f(val1);\n" +
			"/** @type {number} */\n" +
			"var val2 = 0;\n" +
			"var result2 = f(val2);\n")

	AssertThat(c.FindNameType("result1", c.globalScope).String()).IsEqualTo("string")
	AssertThat(c.FindNameType("result2", c.globalScope).String()).IsEqualTo("number")
}

// Java: TypedScopeCreatorTest.testTemplateType4()
func (c *TypedScopeCreatorTest) TestTemplateType4() {
	c.TestSame(
		"/**\n" +
			" * @param {T} x\n" +
			" * @return {T}\n" +
			" * @template T\n" +
			" */\n" +
			"function f(x) {\n" +
			"  return x;\n" +
			"}\n" +
			"/** @type {!Array<string>} */\n" +
			"var arr = [];\n" +
			"(function() {var result = f(arr);})();\n")

	resultType := c.FindNameType("result", c.lastLocalScope)
	AssertThat(resultType.String()).IsEqualTo("Array<string>")
}

// Java: TypedScopeCreatorTest.testTemplateType4a()
func (c *TypedScopeCreatorTest) TestTemplateType4a() {
	c.TestSame(
		"/**\n" +
			" * @param {function():T} x\n" +
			" * @return {T}\n" +
			" * @template T\n" +
			" */\n" +
			"function f(x) {\n" +
			"  return x;\n" +
			"}\n" +
			"/** @return {string} */\n" +
			"var g = function(){return 'hi'};\n" +
			"(function() {var result = f(g);})();\n")

	resultType := c.FindNameType("result", c.lastLocalScope)
	AssertThat(resultType.String()).IsEqualTo("string")
}

// Java: TypedScopeCreatorTest.testTemplateType4b()
func (c *TypedScopeCreatorTest) TestTemplateType4b() {
	c.TestSame(
		"/**\n" +
			" * @param {function(T):void} x\n" +
			" * @return {T}\n" +
			" * @template T\n" +
			" */\n" +
			"function f(x) {\n" +
			"  return x;\n" +
			"}\n" +
			"/** @param {string} x */\n" +
			"var g = function(x){};\n" +
			"(function() {var result = f(g);})();\n")

	resultType := c.FindNameType("result", c.lastLocalScope)
	AssertThat(resultType.String()).IsEqualTo("string")
}

// Java: TypedScopeCreatorTest.testTemplateType5()
func (c *TypedScopeCreatorTest) TestTemplateType5() {
	c.TestSame(
		"/**\n" +
			" * @param {Array<T>} arr\n" +
			" * @return {!Array<T>}\n" +
			" * @template T\n" +
			" */\n" +
			"function f(arr) {\n" +
			"  return arr;\n" +
			"}\n" +
			"/** @type {Array<string>} */\n" +
			"var arr = [];\n" +
			"var result = f(arr);\n")

	AssertThat(c.FindNameTypeStr("result", c.globalScope)).IsEqualTo("Array<string>")
}

// Java: TypedScopeCreatorTest.testTemplateType6()
func (c *TypedScopeCreatorTest) TestTemplateType6() {
	c.TestSame(
		"/**\n" +
			" * @param {Array<T>|string|undefined} arr\n" +
			" * @return {!Array<T>}\n" +
			" * @template T\n" +
			" */\n" +
			"function f(arr) {\n" +
			"  return arr;\n" +
			"}\n" +
			"/** @type {Array<string>} */\n" +
			"var arr = [];\n" +
			"var result = f(arr);\n")

	AssertThat(c.FindNameTypeStr("result", c.globalScope)).IsEqualTo("Array<string>")
}

// Java: TypedScopeCreatorTest.testTemplateType7()
func (c *TypedScopeCreatorTest) TestTemplateType7() {
	c.TestSame(
		"goog.array = {};\n" +
			"/**\n" +
			" * @param {Array<T>} arr\n" +
			" * @param {function(this:S, !T, number, !Array<!T>):boolean} f\n" +
			" * @param {!S=} opt_obj\n" +
			" * @return {!Array<T>}\n" +
			" * @template T,S\n" +
			" */\n" +
			"goog.array.filter = function(arr, f, opt_obj) {\n" +
			"  var res = [];\n" +
			"  for (var i = 0; i < arr.length; i++) {\n" +
			"     const val = arr[i];\n" +
			"     if (f.call(opt_obj, val, i, arr)) {\n" +
			"        res.push(val);\n" +
			"     }\n" +
			"  }\n" +
			"  return res;\n" +
			"}\n" +
			"/** @constructor */\n" +
			"function Foo() {}\n" +
			"/** @type {Array<string>} */\n" +
			"var arr = [];\n" +
			"var result = goog.array.filter(arr,\n" +
			"  function(a,b,c) {var self=this;}, new Foo());\n")

	AssertThat(c.FindNameType("self", c.lastFunctionScope).String()).IsEqualTo("Foo")
	AssertThat(c.FindNameType("a", c.lastFunctionScope).String()).IsEqualTo("string")
	AssertThat(c.FindNameType("b", c.lastFunctionScope).String()).IsEqualTo("number")
	AssertThat(c.FindNameType("c", c.lastFunctionScope).String()).IsEqualTo("Array<string>")
	AssertThat(c.FindNameType("result", c.globalScope).String()).IsEqualTo("Array<string>")
}

// Java: TypedScopeCreatorTest.testTemplateType7b()
func (c *TypedScopeCreatorTest) TestTemplateType7b() {
	c.TestSame(
		"goog.array = {};\n" +
			"/**\n" +
			" * @param {Array<T>} arr\n" +
			" * @param {function(this:S, !T, number, !Array<T>):boolean} f\n" +
			" * @param {!S=} opt_obj\n" +
			" * @return {!Array<T>}\n" +
			" * @template T,S\n" +
			" */\n" +
			"goog.array.filter = function(arr, f, opt_obj) {\n" +
			"  var res = [];\n" +
			"  for (var i = 0; i < arr.length; i++) {\n" +
			"     const val = arr[i];\n" +
			"     if (f.call(opt_obj, val, i, arr)) {\n" +
			"        res.push(val);\n" +
			"     }\n" +
			"  }\n" +
			"  return res;\n" +
			"}\n" +
			"/** @constructor */\n" +
			"function Foo() {}\n" +
			"/** @type {Array<string>} */\n" +
			"var arr = [];\n" +
			"var result = goog.array.filter(arr,\n" +
			"  function(a,b,c) {var self=this;}, new Foo());\n")

	AssertThat(c.FindNameType("self", c.lastFunctionScope).String()).IsEqualTo("Foo")
	AssertThat(c.FindNameType("a", c.lastFunctionScope).String()).IsEqualTo("string")
	AssertThat(c.FindNameType("b", c.lastFunctionScope).String()).IsEqualTo("number")
	AssertThat(c.FindNameType("c", c.lastFunctionScope).String()).IsEqualTo("Array<string>")
	AssertThat(c.FindNameType("result", c.globalScope).String()).IsEqualTo("Array<string>")
}

// Java: TypedScopeCreatorTest.testTemplateType7c()
func (c *TypedScopeCreatorTest) TestTemplateType7c() {
	c.TestSame(
		"goog.array = {};\n" +
			"/**\n" +
			" * @param {Array<T>} arr\n" +
			" * @param {function(this:S, T, number, Array<T>):boolean} f\n" +
			" * @param {!S=} opt_obj\n" +
			" * @return {!Array<T>}\n" +
			" * @template T,S\n" +
			" */\n" +
			"goog.array.filter = function(arr, f, opt_obj) {\n" +
			"  var res = [];\n" +
			"  for (var i = 0; i < arr.length; i++) {\n" +
			"     const val = arr[i];\n" +
			"     if (f.call(opt_obj, val, i, arr)) {\n" +
			"        res.push(val);\n" +
			"     }\n" +
			"  }\n" +
			"  return res;\n" +
			"}\n" +
			"/** @constructor */\n" +
			"function Foo() {}\n" +
			"/** @type {Array<string>} */\n" +
			"var arr = [];\n" +
			"var result = goog.array.filter(arr,\n" +
			"  function(a,b,c) {var self=this;}, new Foo());\n")

	AssertThat(c.FindNameType("self", c.lastFunctionScope).String()).IsEqualTo("Foo")
	AssertThat(c.FindNameType("a", c.lastFunctionScope).String()).IsEqualTo("string")
	AssertThat(c.FindNameType("b", c.lastFunctionScope).String()).IsEqualTo("number")
	AssertThat(c.FindNameType("c", c.lastFunctionScope).String()).IsEqualTo("(Array<string>|null)")
	AssertThat(c.FindNameType("result", c.globalScope).String()).IsEqualTo("Array<string>")
}

// Java: TypedScopeCreatorTest.testTemplateType8()
func (c *TypedScopeCreatorTest) IgnoredTestTemplateType8() {
	// TODO(johnlenz): somehow allow templated typedefs
	c.TestSame(
		"/** @constructor */ NodeList = function() {};\n" +
			"/** @constructor */ Arguments = function() {};\n" +
			"var goog = {};\n" +
			"goog.array = {};\n" +
			"/**\n" +
			" * @typedef {Array<T>|NodeList|Arguments|{length: number}}\n" +
			" * @template T\n" +
			" */\n" +
			"goog.array.ArrayLike;\n" +
			"/**\n" +
			" * @param {function(this:T, ...)} fn\n" +
			" * @param {T} thisObj\n" +
			" * @template T\n" +
			" */\n" +
			"function bind(fn, thisObj) {}\n" +
			"/** @constructor */\n" +
			"function Foo() {}\n" +
			"/** @return {number} */\n" +
			"Foo.prototype.baz = function() {};\n" +
			"bind(function() { var g = this; var f = this.baz(); }, new Foo());\n")
	AssertThat(c.FindNameType("g", c.lastLocalScope).String()).IsEqualTo("T")
	AssertThat(c.FindNameType("g", c.lastLocalScope).Equals(c.registry.GetType(nil, "Foo"))).IsTrue()
	AssertThat(c.FindNameType("f", c.lastLocalScope).String()).IsEqualTo("number")
}

// Java: TypedScopeCreatorTest.testTemplateType9()
func (c *TypedScopeCreatorTest) TestTemplateType9() {
	c.TestSame(
		"/** @constructor */\n" +
			"function Foo() {}\n" +
			"/**\n" +
			" * @this {T}\n" +
			" * @return {T}\n" +
			" * @template T\n" +
			" */\n" +
			"Foo.prototype.method = function() {};\n" +
			"/**\n" +
			" * @constructor\n" +
			" * @extends {Foo}\n" +
			" */\n" +
			"function Bar() {}\n" +
			"\n" +
			"var g = new Bar().method();\n")
	AssertThat(c.FindNameType("g", c.globalScope).String()).IsEqualTo("Bar")
}

// Java: TypedScopeCreatorTest.testTemplateType10()
func (c *TypedScopeCreatorTest) TestTemplateType10() {
	// NOTE: we would like the type within the function to remain "Foo"
	// we can handle this by support template type like "T extends Foo"
	// to provide a "minimum" type for "Foo" within the function body.
	c.TestSame(
		"/** @constructor */\n" +
			"function Foo() {}\n" +
			"\n" +
			"/**\n" +
			" * @this {T}\n" +
			" * @return {T} fn\n" +
			" * @template T\n" +
			" */\n" +
			"Foo.prototype.method = function() {var g = this;};\n")
	AssertThat(c.FindNameType("g", c.lastLocalScope).String()).IsEqualTo("T")
}

// Java: TypedScopeCreatorTest.testTemplateType_shadowingNativeType()
func (c *TypedScopeCreatorTest) TestTemplateType_shadowingNativeType() {
	// Verify that some special handling of externs for types built-in to JSTypeRegistry, like
	// Array, does not incorrectly run for local types shadowing native types.
	c.TestSame(
		"function local() {\n" +
			"  /** @template U @constructor @param {U} e*/\n" +
			"  function Array(e) { this.e = e; }\n" +
			"  /** @return {U} */ Array.prototype.first = function() { return this.e; };\n" +
			"  /** @type {!Array<string>} */ let arr = new Array();\n" +
			"\n" +
			"  const first = arr.first();\n" +
			"  LOCAL: local;\n" +
			"}\n")

	TypeSubject_AssertType(c.FindNameType("arr", c.lastLocalScope)).ToStringIsEqualTo("Array<string>")
	TypeSubject_AssertType(c.FindNameType("first", c.lastLocalScope)).IsString()

	TypeSubject_AssertType(c.FindNameType("Array", c.lastLocalScope)).
		IsNotEqualTo(c.registry.GetNativeType(rhino.JSTypeNative_ARRAY_FUNCTION_TYPE))
}

// Java: TypedScopeCreatorTest.testTemplatedThis_inClassInstanceMethod_isInferredToBe_receiverType()
func (c *TypedScopeCreatorTest) TestTemplatedThis_inClassInstanceMethod_isInferredToBe_receiverType() {
	c.TestSame(
		"class Foo {\n" +
			"  /**\n" +
			"   * @template THIS\n" +
			"   * @this {THIS}\n" +
			"   * @return {THIS}\n" +
			"   */\n" +
			"  clone() { return this; }\n" +
			"}\n" +
			"\n" +
			"var result = new Foo().clone();\n")

	TypeSubject_AssertType(c.FindNameType("result", c.globalScope)).ToStringIsEqualTo("Foo")
}

// Java: TypedScopeCreatorTest.testTemplatedThis_inClassInstanceMethod_invokedOnSuper_isInferredToBe_subtype()
func (c *TypedScopeCreatorTest) TestTemplatedThis_inClassInstanceMethod_invokedOnSuper_isInferredToBe_subtype() {
	c.TestSame(
		"class Foo {\n" +
			"  /**\n" +
			"   * @template THIS\n" +
			"   * @this {THIS}\n" +
			"   * @return {THIS}\n" +
			"   */\n" +
			"  clone() { return this; }\n" +
			"}\n" +
			"\n" +
			"class SubFoo extends Foo {\n" +
			"  other() {\n" +
			"    LABEL: super.clone();\n" +
			"  }\n" +
			"}\n")

	superCloneCall := c.GetLabeledStatement("LABEL").statementNode.GetOnlyChild()
	TypeSubject_AssertType(superCloneCall.GetJSType()).ToStringIsEqualTo("SubFoo")
}

// Java: TypedScopeCreatorTest.testTemplateType11()
func (c *TypedScopeCreatorTest) TestTemplateType11() {
	c.TestSame(
		"/**\n" +
			" * @this {T}\n" +
			" * @return {T} fn\n" +
			" * @template T\n" +
			" */\n" +
			"var method = function() {};\n" +
			"/**\n" +
			" * @constructor\n" +
			" */\n" +
			"function Bar() {}\n" +
			"\n" +
			"var g = method().call(new Bar());\n")
	// NOTE: we would like this to be "Bar"
	AssertThat(c.FindNameType("g", c.globalScope).String()).IsEqualTo("?")
}

// Java: TypedScopeCreatorTest.testTemplateType12()
func (c *TypedScopeCreatorTest) TestTemplateType12() {
	c.TestSame(
		"/** @constructor */\n" +
			"function Foo() {}\n" +
			"\n" +
			"/**\n" +
			" * @this {Array<T>|{length:number}}\n" +
			" * @return {T} fn\n" +
			" * @template T\n" +
			" */\n" +
			"Foo.prototype.method = function() {var g = this;};\n")
	AssertThat(c.FindNameType("g", c.lastLocalScope).String()).
		IsEqualTo("(Array<T>|{length: number})")
}

// Java: TypedScopeCreatorTest.testTemplateType13()
func (c *TypedScopeCreatorTest) IgnoredTestTemplateType13() {
	// TODO(johnlenz): allow template types in @type function expressions
	c.TestSame(
		"/**\n" +
			" * @type {function(T):T}\n" +
			" * @template T\n" +
			" */\n" +
			"var f;\n" +
			"/** @type {string} */\n" +
			"var val = 'hi';\n" +
			"var result = f(val);\n")
	AssertThat(c.FindNameType("result", c.globalScope).String()).IsEqualTo("(string|undefined)")
}

// Java: TypedScopeCreatorTest.testClassTemplateType1()
func (c *TypedScopeCreatorTest) TestClassTemplateType1() {
	// Verify that template types used in method signature are resolved.
	c.TestSame(
		"/**\n" +
			" * @constructor\n" +
			" * @template T\n" +
			" */\n" +
			"function C() {};\n" +
			"\n" +
			"/** @return {T} */\n" +
			"C.prototype.method = function() {}\n" +
			"\n" +
			"/** @type {C<string>} */ var x = new C();\n" +
			"var result = x.method();\n")
	AssertThat(c.FindNameType("result", c.globalScope).String()).IsEqualTo("string")
}

// Java: TypedScopeCreatorTest.testClassTemplateType2()
func (c *TypedScopeCreatorTest) TestClassTemplateType2() {
	// Verify that template types used in method signature on namespaced
	// objects are resolved.
	c.TestSame(
		"/** @const */ var ns = {};\n" +
			"/**\n" +
			" * @constructor\n" +
			" * @template T\n" +
			" */\n" +
			"ns.C = function() {};\n" +
			"\n" +
			"/** @return {T} */\n" +
			"ns.C.prototype.method = function() {}\n" +
			"\n" +
			"/** @type {ns.C<string>} */ var x = new ns.C();\n" +
			"var result = x.method();\n")
	AssertThat(c.FindNameType("result", c.globalScope).String()).IsEqualTo("string")
}

// Java: TypedScopeCreatorTest.testClassTemplateType3()
func (c *TypedScopeCreatorTest) TestClassTemplateType3() {
	// Verify that template types used for instance properties are recognized.
	c.TestSame(
		"/**\n" +
			" * @constructor\n" +
			" * @template T\n" +
			" */\n" +
			"function C() {\n" +
			"  /** @type {T} */\n" +
			"  this.foo;\n" +
			"};\n" +
			"\n" +
			"/** @type {C<string>} */ var x = new C();\n" +
			"var result = x.foo;\n")
	AssertThat(c.FindNameType("result", c.globalScope).String()).IsEqualTo("string")
}

// Java: TypedScopeCreatorTest.testClassTemplateType4_instanceProperty()
func (c *TypedScopeCreatorTest) TestClassTemplateType4_instanceProperty() {
	// Verify that template types used for instance properties are recognized.
	c.TestSame(
		"/** @const */ var ns = {};\n" +
			"/**\n" +
			" * @constructor\n" +
			" * @template T\n" +
			" */\n" +
			"ns.C = function() {\n" +
			"  /** @type {T} */\n" +
			"  this.foo;\n" +
			"};\n" +
			"\n" +
			"/** @type {ns.C<string>} */ var x = new ns.C();\n" +
			"var result = x.foo;\n")
	AssertThat(c.FindNameType("result", c.globalScope).String()).IsEqualTo("string")
}

// Java: TypedScopeCreatorTest.testClassTemplateType4_instancePropertyWithInheritance()
func (c *TypedScopeCreatorTest) TestClassTemplateType4_instancePropertyWithInheritance() {
	// Verify that template types used for instance properties are recognized
	// when shadowing a template type on the superclass with the same name.
	c.TestSame(
		"/**\n" +
			" * @template T\n" +
			" */\n" +
			"class Parent {}\n" +
			"\n" +
			"/**\n" +
			" * @template T\n" +
			" * @extends {Parent<number>}\n" +
			" */\n" +
			"class Child extends Parent {\n" +
			"  constructor() {\n" +
			"    super();\n" +
			"    /** @type {T} */\n" +
			"    this.foo;\n" +
			"  }\n" +
			"};\n" +
			"\n" +
			"/** @type {!Child<string>} */ var x = new Child();\n" +
			"var result = x.foo;\n")
	AssertThat(c.FindNameType("result", c.globalScope).String()).IsEqualTo("string")
}

// Java: TypedScopeCreatorTest.testClassTemplateType5()
func (c *TypedScopeCreatorTest) TestClassTemplateType5() {
	// Verify that template types used for prototype properties in stub
	// declarations are recognized.
	c.TestSame(
		"/**\n" +
			" * @constructor\n" +
			" * @template T\n" +
			" */\n" +
			"function C() {\n" +
			"};\n" +
			"\n" +
			"/** @type {T} */\n" +
			"C.prototype.foo;\n" +
			"\n" +
			"/** @type {C<string>} */ var x = new C();\n" +
			"var result = x.foo;\n")
	AssertThat(c.FindNameType("result", c.globalScope).String()).IsEqualTo("string")
}

// Java: TypedScopeCreatorTest.testClassTemplateType6()
func (c *TypedScopeCreatorTest) TestClassTemplateType6() {
	// Verify that template types used for prototype properties in assignment
	// expressions are recognized.
	c.TestSame(
		"/**\n" +
			" * @constructor\n" +
			" * @template T\n" +
			" */\n" +
			"function C() {\n" +
			"};\n" +
			"\n" +
			"/** @type {T} */\n" +
			"C.prototype.foo = 1;\n" +
			"\n" +
			"/** @type {C<string>} */ var x = new C();\n" +
			"var result = x.foo;\n")
	AssertThat(c.FindNameType("result", c.globalScope).String()).IsEqualTo("string")
}

// Java: TypedScopeCreatorTest.testClassTemplateType7()
func (c *TypedScopeCreatorTest) TestClassTemplateType7() {
	// Verify that template types used in prototype methods are recognized.
	c.TestSame(
		"/**\n" +
			" * @constructor\n" +
			" * @template T\n" +
			" */\n" +
			"function C() {};\n" +
			"\n" +
			"C.prototype.method = function() {\n" +
			"  /** @type {T} */ var local;\n" +
			"}\n")
	AssertThat(c.FindNameType("local", c.lastLocalScope).String()).IsEqualTo("T")
}

// Java: TypedScopeCreatorTest.testClassTemplateType8()
func (c *TypedScopeCreatorTest) TestClassTemplateType8() {
	// Verify that template types used in casts are recognized.
	c.TestSame(
		c.Externs("/** @type {?} */ let unknown;"),
		c.Srcs(
			"/**\n"+
				" * @constructor\n"+
				" * @template T\n"+
				" */\n"+
				"function C() {};\n"+
				"\n"+
				"C.prototype.method = function() {\n"+
				"  var local = /** @type {T} */ (unknown);\n"+
				"}\n"))
	AssertThat(c.FindNameType("local", c.lastLocalScope).String()).IsEqualTo("T")
}

// Java: TypedScopeCreatorTest.testClassTemplateInheritance1()
func (c *TypedScopeCreatorTest) TestClassTemplateInheritance1() {
	// Verify that template type inheritance works for prototype properties.
	c.TestSame(
		"/**\n" +
			" * @constructor\n" +
			" * @template T\n" +
			" */\n" +
			"function C() {};\n" +
			"\n" +
			"/** @type {T} */\n" +
			"C.prototype.foo = 1;\n" +
			"\n" +
			"/**\n" +
			" * @constructor\n" +
			" * @template T, U\n" +
			" * @extends {C<U>}\n" +
			" */\n" +
			"function D() {};\n" +
			"\n" +
			"/** @type {T} */\n" +
			"D.prototype.bar;\n" +
			"\n" +
			"/** @type {D<string, number>} */ var x = new D();\n" +
			"var result1 = x.foo;\n" +
			"var result2 = x.bar;\n")
	AssertThat(c.FindNameType("result1", c.globalScope).String()).IsEqualTo("number")
	AssertThat(c.FindNameType("result2", c.globalScope).String()).IsEqualTo("string")
}

// Java: TypedScopeCreatorTest.testClassTemplateInheritance2()
func (c *TypedScopeCreatorTest) TestClassTemplateInheritance2() {
	// Verify that template type inheritance works for properties and methods.
	c.TestSame(
		"/**\n" +
			" * @constructor\n" +
			" * @template T\n" +
			" */\n" +
			"function C() {};\n" +
			"\n" +
			"/** @return {T} */\n" +
			"C.prototype.method1 = function() {}\n" +
			"\n" +
			"/**\n" +
			" * @constructor\n" +
			" * @template T, U\n" +
			" * @extends {C<U>}\n" +
			" */\n" +
			"function D() {};\n" +
			"\n" +
			"/** @return {T} */\n" +
			"D.prototype.method2 = function() {}\n" +
			"\n" +
			"/** @type {D<boolean, string>} */ var x = new D();\n" +
			"var result1 = x.method1();\n" +
			"var result2 = x.method2();\n")
	AssertThat(c.FindNameType("result1", c.globalScope).String()).IsEqualTo("string")
	AssertThat(c.FindNameType("result2", c.globalScope).String()).IsEqualTo("boolean")
}

// Java: TypedScopeCreatorTest.testClassTemplateInheritance3()
func (c *TypedScopeCreatorTest) TestClassTemplateInheritance3() {
	// Verify that template type inheritance works when the superclass template
	// types are not specified.
	c.TestSame(
		"/**\n" +
			" * @constructor\n" +
			" * @template T\n" +
			" */\n" +
			"function C() {\n" +
			"  /** @type {T} */\n" +
			"  this.foo;\n" +
			"};\n" +
			"\n" +
			"/**\n" +
			" * @constructor\n" +
			" * @template T\n" +
			" * @extends {C}\n" +
			" */\n" +
			"function D() {\n" +
			"  /** @type {T} */\n" +
			"  this.bar;\n" +
			"};\n" +
			"\n" +
			"/** @type {D<boolean>} */ var x = new D();\n" +
			"var result1 = x.foo;\n" +
			"var result2 = x.bar;\n")
	AssertThat(c.FindNameType("result1", c.globalScope).String()).IsEqualTo("?")
	// TODO(nicksantos): There's a bug where the template name T clashes between
	// D and C.
	// assertEquals("boolean", findNameType("result2", globalScope).toString());
}

// Java: TypedScopeCreatorTest.testClassTemplateInheritance4()
func (c *TypedScopeCreatorTest) TestClassTemplateInheritance4() {
	// Verify that overriding methods works with template type inheritance.
	c.TestSame(
		"/**\n" +
			" * @constructor\n" +
			" * @template T\n" +
			" */\n" +
			"function C() {};\n" +
			"\n" +
			"/** @return {T} */\n" +
			"C.prototype.method = function() {}\n" +
			"\n" +
			"/**\n" +
			" * @constructor\n" +
			" * @template T, U\n" +
			" * @extends {C<U>}\n" +
			" */\n" +
			"function D() {};\n" +
			"\n" +
			"/** @override */\n" +
			"D.prototype.method = function() {}\n" +
			"\n" +
			"/** @type {D<boolean, string>} */ var x = new D();\n" +
			"var result = x.method();\n")
	AssertThat(c.FindNameType("result", c.globalScope).String()).IsEqualTo("string")
}

// Java: TypedScopeCreatorTest.testClassTemplateInheritance5()
func (c *TypedScopeCreatorTest) TestClassTemplateInheritance5() {
	// Verify that overriding methods works with template type inheritance.
	c.TestSame(
		"/**\n" +
			" * @constructor\n" +
			" * @template T\n" +
			" */\n" +
			"function C() {};\n" +
			"\n" +
			"/** @return {T} */\n" +
			"C.prototype.method1 = function() {}\n" +
			"\n" +
			"/**\n" +
			" * @constructor\n" +
			" * @template T, U\n" +
			" * @extends {C<U>}\n" +
			" */\n" +
			"function D() {};\n" +
			"\n" +
			"/** @return {T} */\n" +
			"D.prototype.method2 = function() {}\n" +
			"\n" +
			"/** @type {D<string, boolean>} */ var x = new D();\n" +
			"/** @type {C<boolean>} */ var y = x;\n" +
			"/** @type {C} */ var z = y;\n" +
			"var result1 = x.method2();\n" +
			"var result2 = y.method1();\n" +
			"var result3 = z.method1();\n")
	AssertThat(c.FindNameType("result1", c.globalScope).String()).IsEqualTo("string")
	AssertThat(c.FindNameType("result2", c.globalScope).String()).IsEqualTo("boolean")
	AssertThat(c.FindNameType("result3", c.globalScope).String()).IsEqualTo("T")
}

// Java: TypedScopeCreatorTest.testClosureParameterTypesWithoutJSDoc()
func (c *TypedScopeCreatorTest) TestClosureParameterTypesWithoutJSDoc() {
	c.TestSame(
		"/**\n" +
			" * @param {function(!Object)} bar\n" +
			" */\n" +
			"function foo(bar) {}\n" +
			"foo(function(baz) { var f = baz; })\n")
	AssertThat(c.FindNameType("f", c.lastLocalScope).String()).IsEqualTo("Object")
}

// Java: TypedScopeCreatorTest.testDuplicateExternProperty1()
func (c *TypedScopeCreatorTest) TestDuplicateExternProperty1() {
	c.TestSame(
		"/** @constructor */ function Foo() {}\n" +
			"Foo.prototype.bar;\n" +
			"/** @type {number} */ Foo.prototype.bar; var x = (new Foo).bar;\n")
	AssertThat(c.FindNameType("x", c.globalScope).String()).IsEqualTo("number")
}

// Java: TypedScopeCreatorTest.testDuplicateExternProperty2()
func (c *TypedScopeCreatorTest) TestDuplicateExternProperty2() {
	c.TestSame(
		"/** @constructor */ function Foo() {}\n" +
			"/** @type {number} */ Foo.prototype.bar;\n" +
			"Foo.prototype.bar; var x = (new Foo).bar;\n")
	AssertThat(c.FindNameType("x", c.globalScope).String()).IsEqualTo("number")
}

// Java: TypedScopeCreatorTest.testAbstractMethodOnPrototype()
func (c *TypedScopeCreatorTest) TestAbstractMethodOnPrototype() {
	c.TestSame(
		"/** @type {!Function} */ var abstractMethod;\n" +
			"/** @constructor */ function Foo() {}\n" +
			"/** @param {number} x */ Foo.prototype.bar = abstractMethod;\n")
	AssertThat(c.FindNameType("abstractMethod", c.globalScope).String()).IsEqualTo("Function")

	ctor := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Foo", c.globalScope))
	instance := ctor.GetInstanceType()
	AssertThat(instance.String()).IsEqualTo("Foo")

	proto := instance.GetImplicitPrototype()
	AssertThat(proto.String()).IsEqualTo("Foo.prototype")

	AssertThat(proto.GetPropertyType("bar").String()).IsEqualTo("function(this:Foo, number): ?")
}

// Java: TypedScopeCreatorTest.testAbstractMethod_cannotBeAssignedToName()
func (c *TypedScopeCreatorTest) TestAbstractMethod_cannotBeAssignedToName() {
	c.TestSame(
		"/** @type {!Function} */ var abstractMethod;\n" +
			// this results in 'y' being typed as the type of 'abstractMethod'
			"/** @param {number} x */ var y = abstractMethod;")
	AssertThat(c.FindNameType("y", c.globalScope).String()).IsEqualTo("Function")
	AssertThat(c.globalScope.GetVar("y").GetType().String()).IsEqualTo("Function")
}

// Java: TypedScopeCreatorTest.testAbstractMethod_inPrototypeObjectLiteral()
func (c *TypedScopeCreatorTest) TestAbstractMethod_inPrototypeObjectLiteral() {
	c.TestSame(
		"/** @type {!Function} */ var abstractMethod;\n" +
			"/** @constructor */ function Foo() {}\n" +
			"Foo.prototype = {/** @param {number} x */ bar: abstractMethod};\n")
	AssertThat(c.FindNameType("abstractMethod", c.globalScope).String()).IsEqualTo("Function")

	ctor := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("Foo", c.globalScope))
	instance := ctor.GetInstanceType()
	AssertThat(instance.String()).IsEqualTo("Foo")

	proto := instance.GetImplicitPrototype()
	AssertThat(proto.String()).IsEqualTo("Foo.prototype")

	TypeSubject_AssertType(proto.GetPropertyType("bar")).ToStringIsEqualTo("function(this:Foo, number): ?")
}

// Java: TypedScopeCreatorTest.testReturnTypeInference1()
func (c *TypedScopeCreatorTest) TestReturnTypeInference1() {
	c.TestSame("function f() {}")
	AssertThat(c.FindNameType("f", c.globalScope).String()).IsEqualTo("function(): undefined")
}

// Java: TypedScopeCreatorTest.testReturnTypeInference2()
func (c *TypedScopeCreatorTest) TestReturnTypeInference2() {
	c.TestSame("/** @return {?} */ function f() {}")
	AssertThat(c.FindNameType("f", c.globalScope).String()).IsEqualTo("function(): ?")
}

// Java: TypedScopeCreatorTest.testReturnTypeInference3()
func (c *TypedScopeCreatorTest) TestReturnTypeInference3() {
	c.TestSame("function f() {x: return 3;}")
	AssertThat(c.FindNameType("f", c.globalScope).String()).IsEqualTo("function(): ?")
}

// Java: TypedScopeCreatorTest.testReturnTypeInference4()
func (c *TypedScopeCreatorTest) TestReturnTypeInference4() {
	c.TestSame("function f() { throw 'error'; }")
	AssertThat(c.FindNameType("f", c.globalScope).String()).IsEqualTo("function(): ?")
}

// Java: TypedScopeCreatorTest.testReturnTypeInference5()
func (c *TypedScopeCreatorTest) TestReturnTypeInference5() {
	c.TestSame("function f() { if (true) { return 1; } }")
	AssertThat(c.FindNameType("f", c.globalScope).String()).IsEqualTo("function(): ?")
}

// Java: TypedScopeCreatorTest.testLiteralTypesInferred()
func (c *TypedScopeCreatorTest) TestLiteralTypesInferred() {
	c.SetAcceptedLanguage(CompilerOptions_LanguageMode_UNSUPPORTED)
	c.TestSame("null + true + false + 0 + '' + 0n + {}")
	AssertThat(c.FindTokenType(rhino.Token_NULL, c.globalScope).String()).IsEqualTo("null")
	AssertThat(c.FindTokenType(rhino.Token_TRUE, c.globalScope).String()).IsEqualTo("boolean")
	AssertThat(c.FindTokenType(rhino.Token_FALSE, c.globalScope).String()).IsEqualTo("boolean")
	AssertThat(c.FindTokenType(rhino.Token_NUMBER, c.globalScope).String()).IsEqualTo("number")
	AssertThat(c.FindTokenType(rhino.Token_STRINGLIT, c.globalScope).String()).IsEqualTo("string")
	AssertThat(c.FindTokenType(rhino.Token_BIGINT, c.globalScope).String()).IsEqualTo("bigint")
	AssertThat(c.FindTokenType(rhino.Token_OBJECTLIT, c.globalScope).String()).IsEqualTo("{}")
}

// Java: TypedScopeCreatorTest.testGlobalQualifiedNameInLocalScope()
func (c *TypedScopeCreatorTest) TestGlobalQualifiedNameInLocalScope() {
	// TODO(b/154044898): re-enable type info validation
	c.DisableTypeInfoValidation()
	c.TestSame(
		"var ns = {};\n" +
			"(function() {\n" +
			"    /** @param {number} x */ ns.foo = function(x) {}; })();\n" +
			"(function() { ns.foo(3); })();\n")
	TypedVarSubject_AssertThat(c.globalScope.GetVar("ns.foo")).IsNotNull()
	AssertThat(c.globalScope.GetVar("ns.foo").GetType().String()).
		IsEqualTo("function(number): undefined")
}

// Java: TypedScopeCreatorTest.testDeclaredObjectLitProperty1()
func (c *TypedScopeCreatorTest) TestDeclaredObjectLitProperty1() {
	c.TestSame("var x = {/** @type {number} */ y: 3};")
	xType := rhino.ObjectType_Cast(c.globalScope.GetVar("x").GetType())
	AssertThat(xType.GetPropertyType("y").String()).IsEqualTo("number")
	AssertThat(xType.String()).IsEqualTo("{y: number}")
}

// Java: TypedScopeCreatorTest.testDeclaredObjectLitProperty2()
func (c *TypedScopeCreatorTest) TestDeclaredObjectLitProperty2() {
	c.TestSame("var x = {/** @param {number} z */ y: function(z){}};")
	xType := rhino.ObjectType_Cast(c.globalScope.GetVar("x").GetType())
	AssertThat(xType.GetPropertyType("y").String()).IsEqualTo("function(number): undefined")
	AssertThat(xType.String()).IsEqualTo("{y: function(number): undefined}")
}

// Java: TypedScopeCreatorTest.testDeclaredObjectLitProperty3()
func (c *TypedScopeCreatorTest) TestDeclaredObjectLitProperty3() {
	c.TestSame(
		"function f() {\n" +
			"  var x = {/** @return {number} */ y: function(z){ return 3; }};\n" +
			"}\n")
	xType := rhino.ObjectType_Cast(c.lastLocalScope.GetVar("x").GetType())
	AssertThat(xType.GetPropertyType("y").String()).IsEqualTo("function(?): number")
	AssertThat(xType.String()).IsEqualTo("{y: function(?): number}")
}

// Java: TypedScopeCreatorTest.testDeclaredObjectLitProperty4()
func (c *TypedScopeCreatorTest) TestDeclaredObjectLitProperty4() {
	c.TestSame("var x = {y: 5, /** @type {number} */ z: 3};")
	xType := rhino.ObjectType_Cast(c.globalScope.GetVar("x").GetType())
	AssertThat(xType.GetPropertyType("y").String()).IsEqualTo("number")
	AssertThat(xType.IsPropertyTypeDeclared("y")).IsFalse()
	AssertThat(xType.IsPropertyTypeDeclared("z")).IsTrue()
	AssertThat(xType.String()).IsEqualTo("{\n  y: number,\n  z: number\n}")
}

// Java: TypedScopeCreatorTest.testDeclaredObjectLitProperty5()
func (c *TypedScopeCreatorTest) TestDeclaredObjectLitProperty5() {
	// TODO(b/154044898): re-enable type info validation
	c.DisableTypeInfoValidation()
	c.TestSame(
		"var x = {/** @type {number} */ prop: 3};\n" +
			"function f() { var y = x.prop; }\n")
	yType := c.lastLocalScope.GetVar("y").GetType()
	AssertThat(yType.String()).IsEqualTo("number")
}

// Java: TypedScopeCreatorTest.testDeclaredObjectLitProperty6()
func (c *TypedScopeCreatorTest) TestDeclaredObjectLitProperty6() {
	c.TestSame("var x = {/** This is JsDoc */ prop: function(){}};")
	prop := c.globalScope.GetVar("x.prop")
	propType := prop.GetType()
	AssertThat(propType.String()).IsEqualTo("function(): undefined")
	AssertThat(prop.IsTypeInferred()).IsFalse()
	AssertThat(rhino.ObjectType_Cast(c.globalScope.GetVar("x").GetType()).IsPropertyTypeInferred("prop")).
		IsFalse()
}

// Java: TypedScopeCreatorTest.testInferredObjectLitProperty1()
func (c *TypedScopeCreatorTest) TestInferredObjectLitProperty1() {
	c.TestSame("var x = {prop: 3};")
	prop := c.globalScope.GetVar("x.prop")
	propType := prop.GetType()
	AssertThat(propType.String()).IsEqualTo("number")
	AssertThat(prop.IsTypeInferred()).IsTrue()
	AssertThat(rhino.ObjectType_Cast(c.globalScope.GetVar("x").GetType()).IsPropertyTypeInferred("prop")).
		IsTrue()
}

// Java: TypedScopeCreatorTest.testInferredObjectLitProperty2()
func (c *TypedScopeCreatorTest) TestInferredObjectLitProperty2() {
	c.TestSame("var x = {prop: function(){}};")
	prop := c.globalScope.GetVar("x.prop")
	propType := prop.GetType()
	AssertThat(propType.String()).IsEqualTo("function(): undefined")
	AssertThat(prop.IsTypeInferred()).IsTrue()
	AssertThat(rhino.ObjectType_Cast(c.globalScope.GetVar("x").GetType()).IsPropertyTypeInferred("prop")).
		IsTrue()
}

// Java: TypedScopeCreatorTest.testDeclaredConstType1()
func (c *TypedScopeCreatorTest) TestDeclaredConstType1() {
	c.TestSame(
		"/** @const */ var x = 3;\n" +
			"function f() { var y = x; }\n")
	yType := c.lastLocalScope.GetVar("y").GetType()
	AssertThat(yType.String()).IsEqualTo("number")
}

// Java: TypedScopeCreatorTest.testDeclaredConstType_symbolCall()
func (c *TypedScopeCreatorTest) TestDeclaredConstType_symbolCall() {
	c.TestSame("/** @const */ var x = Symbol('test');  function f() { var y = x; }")
	yType := c.lastLocalScope.GetVar("y").GetType()
	AssertThat(yType).IsEqualTo(c.registry.GetNativeType(rhino.JSTypeNative_SYMBOL_TYPE))
}

// Java: TypedScopeCreatorTest.testDeclaredConstType_shadowingGlobalSymbol_notConfusedForNativeSymbol()
func (c *TypedScopeCreatorTest) TestDeclaredConstType_shadowingGlobalSymbol_notConfusedForNativeSymbol() {
	c.TestSame(
		"(function() {\n" +
			"function Symbol() {}\n" +
			"/** @const */ var x = Symbol();\n" +
			"function f() { var y = x; Y: y; }\n" +
			"})()\n")
	yType := c.GetLabeledStatement("Y").statementNode.GetOnlyChild().GetJSType()
	AssertThat(yType).IsEqualTo(c.registry.GetNativeType(rhino.JSTypeNative_VOID_TYPE))
}

// Java: TypedScopeCreatorTest.testDeclaredConstType2()
func (c *TypedScopeCreatorTest) TestDeclaredConstType2() {
	c.TestSame(
		"/** @const */ var x = {};\n" +
			"function f() { var y = x; }\n")
	yType := c.lastLocalScope.GetVar("y").GetType()
	AssertThat(yType.String()).IsEqualTo("{}")
}

// Java: TypedScopeCreatorTest.testDeclaredConstType3()
func (c *TypedScopeCreatorTest) TestDeclaredConstType3() {
	c.TestSame(
		"/** @const */ var x = {};\n" +
			"/** @const */ x.z = 'hi';\n" +
			"function f() { var y = x.z; }\n")
	yType := c.lastLocalScope.GetVar("y").GetType()
	AssertThat(yType.String()).IsEqualTo("string")
}

// Java: TypedScopeCreatorTest.testDeclaredConstType4()
func (c *TypedScopeCreatorTest) TestDeclaredConstType4() {
	c.TestSame(
		"/** @constructor */ function Foo() {}\n" +
			"/** @const */ Foo.prototype.z = 'hi';\n" +
			"function f() { var y = (new Foo()).z; }\n")
	yType := c.lastLocalScope.GetVar("y").GetType()
	AssertThat(yType.String()).IsEqualTo("string")

	fooType := (typedScopeCreatorTest_cast[rhino.FunctionType](c.globalScope.GetVar("Foo").GetType())).GetInstanceType()
	AssertThat(fooType.IsPropertyTypeDeclared("z")).IsTrue()
}

// Java: TypedScopeCreatorTest.testDeclaredConstType5a()
func (c *TypedScopeCreatorTest) TestDeclaredConstType5a() {
	c.TestSame(
		c.Externs(""),
		c.Srcs(
			"/** @const */ var goog = goog || {};\n"+
				"function f() { var y = goog; }\n"))
	yType := c.lastLocalScope.GetVar("y").GetType()
	AssertThat(yType.String()).IsEqualTo("{}")
}

// Java: TypedScopeCreatorTest.testDeclaredConstType6()
func (c *TypedScopeCreatorTest) TestDeclaredConstType6() {
	c.TestSame(
		"/**\n" +
			" * @param {{y:string}} a\n" +
			" * @constructor\n" +
			"*/\n" +
			"var C = function(a) { /** @const */ this.x = a.y;};\n" +
			"var instance = new C({y:'str'})\n")
	instance := typedScopeCreatorTest_cast[rhino.ObjectType](c.FindNameType("instance", c.globalScope))
	AssertThat(instance.String()).IsEqualTo("C")
	AssertThat(instance.HasProperty("x")).IsTrue()
	AssertThat(instance.GetPropertyType("x").String()).IsEqualTo("string")
	AssertThat(instance.IsPropertyTypeInferred("x")).IsFalse()
}

// Java: TypedScopeCreatorTest.testBadCtorInit1()
func (c *TypedScopeCreatorTest) TestBadCtorInit1() {
	c.TestWarning("/** @constructor */ var f;", TypedScopeCreator_CTOR_INITIALIZER)
}

// Java: TypedScopeCreatorTest.testBadCtorInit2()
func (c *TypedScopeCreatorTest) TestBadCtorInit2() {
	c.TestWarning("var x = {}; /** @constructor */ x.f;", TypedScopeCreator_CTOR_INITIALIZER)
}

// Java: TypedScopeCreatorTest.testBadIfaceInit1()
func (c *TypedScopeCreatorTest) TestBadIfaceInit1() {
	c.TestWarning("/** @interface */ var f;", TypedScopeCreator_IFACE_INITIALIZER)
}

// Java: TypedScopeCreatorTest.testBadIfaceInit2()
func (c *TypedScopeCreatorTest) TestBadIfaceInit2() {
	c.TestWarning("var x = {}; /** @interface */ x.f;", TypedScopeCreator_IFACE_INITIALIZER)
}

// Java: TypedScopeCreatorTest.testDeclaredCatchExpression1()
func (c *TypedScopeCreatorTest) TestDeclaredCatchExpression1() {
	// TODO(b/154044898): re-enable type info validation
	c.DisableTypeInfoValidation()
	c.TestSame("try {} catch (e) {}")
	TypedVarSubject_AssertThat(c.lastLocalScope.GetVar("e")).HasJSTypeThat().IsUnknown()
}

// Java: TypedScopeCreatorTest.testDeclaredCatchExpression2()
func (c *TypedScopeCreatorTest) TestDeclaredCatchExpression2() {
	c.TestSame("try {} catch (/** @type {string} */ e) {}")
	AssertThat(c.lastLocalScope.GetVar("e").GetType().String()).IsEqualTo("string")
}

// Java: TypedScopeCreatorTest.testDestructuringCatch()
func (c *TypedScopeCreatorTest) TestDestructuringCatch() {
	c.TestSame("try {} catch ({/** string */ message, /** number */ errno}) {}")
	TypeSubject_AssertType(c.lastLocalScope.GetVar("message").GetType()).ToStringIsEqualTo("string")
	TypeSubject_AssertType(c.lastLocalScope.GetVar("errno").GetType()).ToStringIsEqualTo("number")
}

// Java: TypedScopeCreatorTest.testDuplicateCatchVariableNames()
func (c *TypedScopeCreatorTest) TestDuplicateCatchVariableNames() {
	// TODO(b/154044898): re-enable type info validation
	c.DisableTypeInfoValidation()
	c.TestSame(
		"try {} // preserve newlines\n" +
			"catch (err) {\n" +
			"  FIRST_CATCH: err;\n" +
			"}\n" +
			"try {}\n" +
			"catch (err) {\n" +
			"  SECOND_CATCH: err;\n" +
			"}\n")
	firstCatchScope := c.GetLabeledStatement("FIRST_CATCH").enclosingScope
	ScopeSubject_AssertScope(firstCatchScope).Declares("err").Directly()
	firstErrVar := firstCatchScope.GetVar("err")

	secondCatchScope := c.GetLabeledStatement("SECOND_CATCH").enclosingScope
	ScopeSubject_AssertScope(firstCatchScope).Declares("err").Directly()
	AssertThat(firstCatchScope).IsNotSameInstanceAs(secondCatchScope)

	secondErrVar := secondCatchScope.GetVar("err")
	TypedVarSubject_AssertThat(firstErrVar).IsNotSameInstanceAs(secondErrVar)
}

// Java: TypedScopeCreatorTest.testGenerator1()
func (c *TypedScopeCreatorTest) TestGenerator1() {
	c.TestSame("function *gen() { yield 1; } var g = gen();")
	AssertThat(c.FindNameType("gen", c.globalScope).String()).
		IsEqualTo("function(): Generator<?,?,?>")
	AssertThat(c.FindNameType("g", c.globalScope).String()).IsEqualTo("Generator<?,?,?>")
}

// Java: TypedScopeCreatorTest.testGenerator2()
func (c *TypedScopeCreatorTest) TestGenerator2() {
	c.TestSame("var gen = function *() { yield 1; }; var g = gen();")
	AssertThat(c.FindNameType("gen", c.globalScope).String()).
		IsEqualTo("function(): Generator<?,?,?>")
	AssertThat(c.FindNameType("g", c.globalScope).String()).IsEqualTo("Generator<?,?,?>")
}

// Just check that this doesn't cause a StackOverflowError.
// Java: TypedScopeCreatorTest.testArgumentsStackOverflow()
func (c *TypedScopeCreatorTest) TestArgumentsStackOverflow() {
	// TODO(b/154044898): re-enable type info validation
	c.DisableTypeInfoValidation()
	js :=
		"/**\n" +
			" * @fileoverview\n" +
			" * @suppress {es5Strict}\n" +
			" */\n" +
			"\n" +
			"function getStackTrace() {\n" +
			"  try {\n" +
			"    throw 'error';\n" +
			"  } catch (e) {\n" +
			"    return 0;\n" +
			"  }\n" +
			"\n" +
			"  if (typeof (arguments.caller) != 'undefined') {}\n" +
			"\n" +
			"  return '';\n" +
			"}\n"
	c.TestSame(js)
}

// Java: TypedScopeCreatorTest.testSpreadInIifeBlocksTypeInference()
func (c *TypedScopeCreatorTest) TestSpreadInIifeBlocksTypeInference() {
	c.TestSame(
		"const /** !Array<string> */ arr = [];\n" +
			"(function (x = 4, y = 5) {})(...arr, undefined);\n")

	xVar := jutil.CheckNotNull(c.lastFunctionScope.GetVar("x"))
	TypeSubject_AssertType(xVar.GetType()).ToStringIsEqualTo("?")

	yVar := jutil.CheckNotNull(c.lastFunctionScope.GetVar("y"))
	TypeSubject_AssertType(yVar.GetType()).ToStringIsEqualTo("?")
}

// Java: TypedScopeCreatorTest.testAliasTypeFromClass()
func (c *TypedScopeCreatorTest) TestAliasTypeFromClass() {
	c.TestSame(
		"class Foo {};\n" +
			"/** @enum {number} */ Foo.E = {A: 1};\n" +
			"const F = Foo;\n" +
			"const E = F.E;\n")
	AssertThat(c.FindNameType("E", c.globalScope).String()).IsEqualTo("enum{Foo.E}")
}

// Java: TypedScopeCreatorTest.testAliasTypeFromNamespace()
func (c *TypedScopeCreatorTest) TestAliasTypeFromNamespace() {
	c.TestSame(
		"const Foo = {};\n" +
			"/** @enum {number} */ Foo.E = {A: 1};\n" +
			"const F = Foo;\n" +
			"const E = F.E;\n")
	AssertThat(c.FindNameType("E", c.globalScope).String()).IsEqualTo("enum{Foo.E}")
}

// Java: TypedScopeCreatorTest.testAliasTypeFromClassPrototype()
func (c *TypedScopeCreatorTest) TestAliasTypeFromClassPrototype() {
	// This is very weird but falls out of how we module types.
	c.TestSame(
		"class Foo {}\n" +
			"/** @enum {number} */ Foo.prototype.E = {A: 1};\n" +
			"const F = new Foo();\n" +
			"const E = F.E;\n")
	AssertThat(c.FindNameType("E", c.globalScope).String()).IsEqualTo("enum{Foo.prototype.E}")
}

// Java: TypedScopeCreatorTest.testAliasTypedef()
func (c *TypedScopeCreatorTest) TestAliasTypedef() {
	c.TestSame(
		"/** @typedef {number} */ let Foo;\n" +
			"const E = Foo;\n" +
			"/** @type {E} */ let x;\n")
	AssertThat(c.FindNameType("x", c.globalScope).String()).IsEqualTo("number")
}

// Java: TypedScopeCreatorTest.testAliasTypedefOntoNamespace()
func (c *TypedScopeCreatorTest) TestAliasTypedefOntoNamespace() {
	c.TestSame(
		"/** @typedef {number} */ let Foo;\n" +
			"const E = {};\n" +
			"/** @const */\n" +
			"E.F = Foo;\n" +
			"/** @type {E.F} */ let x;\n")
	TypeSubject_AssertType(c.FindNameType("x", c.globalScope)).IsNumber()
	eType := c.FindNameType("E", c.globalScope).ToMaybeObjectType()
	AssertThat(eType).IsNotNull()
	AssertThat(eType.IsPropertyTypeDeclared("F")).IsTrue()
}

// Java: TypedScopeCreatorTest.testAliasTypedefFromNamespace()
func (c *TypedScopeCreatorTest) TestAliasTypedefFromNamespace() {
	c.TestSame(
		"const Foo = {};\n" +
			"/** @typedef {number} */ Foo.E;\n" +
			"const E = Foo.E;\n" +
			"/** @type {E} */ let x;\n")
	AssertThat(c.FindNameType("x", c.globalScope).String()).IsEqualTo("number")
}

// Java: TypedScopeCreatorTest.testAliasTypedefFromNamespaceAlias()
func (c *TypedScopeCreatorTest) TestAliasTypedefFromNamespaceAlias() {
	c.TestSame(
		"const Foo = {};\n" +
			"/** @typedef {number} */ Foo.E;\n" +
			"const F = Foo;\n" +
			"const E = F.E;\n" +
			"/** @type {E} */ let x;\n")
	AssertThat(c.FindNameType("x", c.globalScope).String()).IsEqualTo("number")
}

// Java: TypedScopeCreatorTest.testTypedefName_hasNoType()
func (c *TypedScopeCreatorTest) TestTypedefName_hasNoType() {
	c.TestSame("/** @typedef {number} */ let Foo;")

	foo := typedVarOf(c.globalScope.GetSlot("Foo"))
	TypedVarSubject_AssertThat(foo).IsNotInferred()
	TypedVarSubject_AssertThat(foo).HasJSTypeThat().IsNoType()
}

// Java: TypedScopeCreatorTest.testTypedefQualifiedName_hasNoType()
func (c *TypedScopeCreatorTest) TestTypedefQualifiedName_hasNoType() {
	c.TestSame("const ns = {}; /** @typedef {number} */ ns.Foo;")

	foo := typedVarOf(c.globalScope.GetSlot("ns.Foo"))
	TypedVarSubject_AssertThat(foo).IsNotInferred()
	TypedVarSubject_AssertThat(foo).HasJSTypeThat().IsNoType() // Note: we would like this to be 'void' eventually.
}

// The {@link CheckJSDoc} pass already warns on these misplaced @typedefs
// Java: TypedScopeCreatorTest.testTypedef_notAllowedOnThisProp()
func (c *TypedScopeCreatorTest) TestTypedef_notAllowedOnThisProp() {
	c.TestSame("class C { constructor() { /** @typedef {string} */ this.name; } }")

	AssertThat(c.registry.GetType(c.globalScope, "this.name")).IsNull()
}

// Java: TypedScopeCreatorTest.testTypedef_notAllowedOnPrototypeProp()
func (c *TypedScopeCreatorTest) TestTypedef_notAllowedOnPrototypeProp() {
	c.TestSame("/** @constructor */ function C() {} /** @typedef {string} */ C.prototype.name;")

	AssertThat(c.registry.GetType(c.globalScope, "C.prototype.name")).IsNull()
}

// Java: TypedScopeCreatorTest.testTypedefName_usingLetWithLiteralRhs_hasNoTypeIfReassigned()
func (c *TypedScopeCreatorTest) TestTypedefName_usingLetWithLiteralRhs_hasNoTypeIfReassigned() {
	// This will cause a type error during typechecking.
	c.TestSame("/** @typedef {number} */ let Foo = 'a string'; Foo = 3;")

	foo := typedVarOf(c.globalScope.GetSlot("Foo"))
	TypedVarSubject_AssertThat(foo).IsNotInferred()
	TypedVarSubject_AssertThat(foo).HasJSTypeThat().IsNoType()
}

// Java: TypedScopeCreatorTest.testTypedefNameVar_usingConstWithLiteralRhs_hasLiteralType()
func (c *TypedScopeCreatorTest) TestTypedefNameVar_usingConstWithLiteralRhs_hasLiteralType() {
	c.TestSame("/** @enum */ let Enum = {A: 0}; /** @typedef {number} */ const Foo = Enum;")

	foo := typedVarOf(c.globalScope.GetSlot("Foo"))
	TypedVarSubject_AssertThat(foo).IsNotInferred()
	// Currently a typedef can be a const alias of a type.
	TypedVarSubject_AssertThat(foo).HasJSTypeThat().ToStringIsEqualTo("enum{Enum}")

	TypeSubject_AssertType(c.registry.GetGlobalType("Foo")).IsNumber()
}

// Java: TypedScopeCreatorTest.testTypedefNameVar_usingConstWithObjectLitRhs_canHaveAdditionalProperties()
func (c *TypedScopeCreatorTest) TestTypedefNameVar_usingConstWithObjectLitRhs_canHaveAdditionalProperties() {
	c.TestSame("/** @typedef {number} */ const Foo = {}; Foo.Builder = () => 0;")

	foo := typedVarOf(c.globalScope.GetSlot("Foo"))
	TypedVarSubject_AssertThat(foo).IsNotInferred()
	TypedVarSubject_AssertThat(foo).HasJSTypeThat().IsObjectTypeWithProperty("Builder")

	TypeSubject_AssertType(c.registry.GetGlobalType("Foo")).IsNumber()
}

// Java: TypedScopeCreatorTest.testTypedefQualifiedName_withLiteralRhs_hasNoType()
func (c *TypedScopeCreatorTest) TestTypedefQualifiedName_withLiteralRhs_hasNoType() {
	c.TestSame("const ns = {}; /** @typedef {number} */ ns.Foo = 'a string';")

	foo := typedVarOf(c.globalScope.GetSlot("ns.Foo"))
	TypedVarSubject_AssertThat(foo).IsNotNull()
	TypedVarSubject_AssertThat(foo).IsNotInferred()
	TypedVarSubject_AssertThat(foo).HasJSTypeThat().IsNoType()
}

// Java: TypedScopeCreatorTest.testTypedefQualifiedName_withObjectLiteralRhs_hasNoType()
func (c *TypedScopeCreatorTest) TestTypedefQualifiedName_withObjectLiteralRhs_hasNoType() {
	c.TestSame("const ns = {}; /** @typedef {number} */ ns.Foo = {}; ns.Foo.Builder = () => 0;")

	foo := typedVarOf(c.globalScope.GetSlot("ns.Foo"))
	TypedVarSubject_AssertThat(foo).IsNotInferred()
	TypedVarSubject_AssertThat(foo).HasJSTypeThat().IsNoType() // Note: eventually we may want to make this {}.

	TypeSubject_AssertType(c.registry.GetGlobalType("ns.Foo")).IsNumber()
}

// Java: TypedScopeCreatorTest.testObjectPatternInParametersWithDefaultGetsCorrectType()
func (c *TypedScopeCreatorTest) TestObjectPatternInParametersWithDefaultGetsCorrectType() {
	c.TestSame("/** @param {{a: string}=} obj */ function f({a = 'foo'} = {}) {}")

	patternType := c.FindTokenType(rhino.Token_OBJECT_PATTERN, c.globalScope)
	// the pattern has the type after the default value is evaluated, not ({a:string}|undefined)
	TypeSubject_AssertType(patternType).ToStringIsEqualTo("({a: string}|{})")

	aNameType := c.FindNameType("a", c.globalScope)
	TypeSubject_AssertType(aNameType).IsString()
}

// Java: TypedScopeCreatorTest.testGoogModuleMissingNamespaceDoesntCrash()
func (c *TypedScopeCreatorTest) TestGoogModuleMissingNamespaceDoesntCrash() {
	c.GetOptions().SetContinueAfterErrors(true)
	c.TestError("goog.module(); const x = 0;", ClosureRewriteModule_INVALID_MODULE_ID_ARG)
	c.TestError("goog.module(0); const x = 0;", ClosureRewriteModule_INVALID_MODULE_ID_ARG)
}

// Java: TypedScopeCreatorTest.testScriptAndGoogModuleWithShadowedVar()
func (c *TypedScopeCreatorTest) TestScriptAndGoogModuleWithShadowedVar() {
	c.TestSame(
		c.Srcs(
			"var x = 0; GLOBAL_X: x;", //
			"goog.module('a'); var x = 'str'; LOCAL_X: x;"))
	globalScope := c.GetLabeledStatement("GLOBAL_X").enclosingScope
	localScope := c.GetLabeledStatement("LOCAL_X").enclosingScope

	globalX := typedVarOf(globalScope.GetSlot("x"))
	localX := typedVarOf(localScope.GetSlot("x"))

	TypedVarSubject_AssertThat(globalX).IsNotEqualTo(localX)
	TypedVarSubject_AssertThat(globalX).HasJSTypeThat().IsNumber()
	TypedVarSubject_AssertThat(localX).HasJSTypeThat().IsString()
}

// Java: TypedScopeCreatorTest.testTwoGoogModulesWithSameNamedVar()
func (c *TypedScopeCreatorTest) TestTwoGoogModulesWithSameNamedVar() {
	c.TestSame(
		c.Srcs(
			"goog.module('a'); var x = 'str'; MOD_A_X: x;", //
			"goog.module('b'); var x = 0; MOD_B_X: x"))

	ScopeSubject_AssertScope(c.globalScope).DoesNotDeclare("x")

	modAX := typedVarOf(c.GetLabeledStatement("MOD_A_X").enclosingScope.GetSlot("x"))
	modBX := typedVarOf(c.GetLabeledStatement("MOD_B_X").enclosingScope.GetSlot("x"))

	TypedVarSubject_AssertThat(modAX).IsNotEqualTo(modBX)
	TypeSubject_AssertType(modAX.GetType()).IsString()
	TypeSubject_AssertType(modBX.GetType()).IsNumber()
}

// Java: TypedScopeCreatorTest.testTwoGoogModuleWithSameNamedFunctionDeclaration()
func (c *TypedScopeCreatorTest) TestTwoGoogModuleWithSameNamedFunctionDeclaration() {
	c.TestSame(
		c.Srcs(
			"goog.module('a'); function foo() {} MOD_A_FOO: foo;", //
			"goog.module('b'); function foo() {} MOD_B_FOO: foo;"))

	ScopeSubject_AssertScope(c.globalScope).DoesNotDeclare("foo")

	modAX := typedVarOf(c.GetLabeledStatement("MOD_A_FOO").enclosingScope.GetSlot("foo"))
	modBX := typedVarOf(c.GetLabeledStatement("MOD_B_FOO").enclosingScope.GetSlot("foo"))

	TypedVarSubject_AssertThat(modAX).IsNotEqualTo(modBX)
}

// Java: TypedScopeCreatorTest.testGoogModule_requiringMissingNamespace()
func (c *TypedScopeCreatorTest) TestGoogModule_requiringMissingNamespace() {
	c.TestSame("goog.module('a'); const b = goog.require('b'); B: b;")

	AssertNode(c.GetLabeledStatement("B").statementNode.GetOnlyChild()).HasJSTypeThat().IsUnknown()
}

// Java: TypedScopeCreatorTest.testGoogModule_requiringMissingNamespace_emptyDestructuring()
func (c *TypedScopeCreatorTest) TestGoogModule_requiringMissingNamespace_emptyDestructuring() {
	c.TestSame("goog.module('a'); const {} = goog.require('b');")
}

// Java: TypedScopeCreatorTest.testGoogModule_moduleBodyScopeIsStoredOnTheInput()
func (c *TypedScopeCreatorTest) TestGoogModule_moduleBodyScopeIsStoredOnTheInput() {
	c.TestSame("goog.module('a'); EXPORTS: exports;")

	moduleScope := c.GetLabeledStatement("EXPORTS").enclosingScope
	moduleBodyNode := moduleScope.GetRootNode()
	AssertNode(moduleBodyNode).HasToken(rhino.Token_MODULE_BODY)
	lastCompiler := c.GetLastCompiler()
	inputId := NodeUtil_GetInputId(moduleBodyNode)
	AssertThat(inputId).IsNotNull()
	compilerInput := lastCompiler.GetInput(inputId)
	AssertThat(compilerInput).IsNotNull()
	typedScopeStoredOnTheCompilerInput := compilerInput.GetTypedScope()
	AssertThat(typedScopeStoredOnTheCompilerInput).IsNotNull()
	AssertThat(typedScopeStoredOnTheCompilerInput).IsSameInstanceAs(moduleScope)
}

// Java: TypedScopeCreatorTest.testGoogModule_declaresExportsVariableImplicitly()
func (c *TypedScopeCreatorTest) TestGoogModule_declaresExportsVariableImplicitly() {
	c.TestSame("goog.module('a'); EXPORTS: exports;")

	exports := typedVarOf(c.GetLabeledStatement("EXPORTS").enclosingScope.GetSlot("exports"))
	TypedVarSubject_AssertThat(exports).IsNotInferred()
	TypedVarSubject_AssertThat(exports).HasJSTypeThat().IsLiteralObject()
}

// Java: TypedScopeCreatorTest.testGoogModule_exportsDeclaresTypeInModuleScope()
func (c *TypedScopeCreatorTest) TestGoogModule_exportsDeclaresTypeInModuleScope() {
	c.TestSame("goog.module('mod.a'); exports = class {}; MOD_A: 0;")

	TypeSubject_AssertType(c.registry.GetGlobalType("exports")).IsNull()
	TypeSubject_AssertType(c.registry.GetType(c.GetLabeledStatement("MOD_A").enclosingScope, "exports")).
		IsNotNull()
}

// Java: TypedScopeCreatorTest.testGoogModule_moduleObjectTypeMatchesImportExportAndModuleBody()
func (c *TypedScopeCreatorTest) TestGoogModule_moduleObjectTypeMatchesImportExportAndModuleBody() {
	c.TestSame(
		c.Srcs(
			"goog.module('mod.a');\n"+
				"exports = class {};\n"+
				"EXPORTED_A: exports;\n",
			"goog.module('mod.b');\n"+
				"const imported_a = goog.require('mod.a');\n"+
				"IMPORTED_A: imported_a;\n"))

	aExportsNode := c.GetLabeledStatement("EXPORTED_A").statementNode.GetOnlyChild()
	exportedAJSType := aExportsNode.GetJSType()
	AssertThat(exportedAJSType).IsNotNull()

	aModuleNode := NodeUtil_GetEnclosingType(aExportsNode, rhino.Token_MODULE_BODY)
	AssertNode(aModuleNode).HasJSTypeThat().IsEqualTo(exportedAJSType)

	importedANode := c.GetLabeledStatement("IMPORTED_A").statementNode.GetOnlyChild()
	AssertNode(importedANode).HasJSTypeThat().IsEqualTo(exportedAJSType)
}

// Java: TypedScopeCreatorTest.testGoogModule_namedExportDeclaresTypeInModuleScope()
func (c *TypedScopeCreatorTest) TestGoogModule_namedExportDeclaresTypeInModuleScope() {
	c.TestSame("goog.module('mod.a'); exports.B = class {}; MOD_A: 0;")

	TypeSubject_AssertType(c.registry.GetGlobalType("exports.B")).IsNull()
	TypeSubject_AssertType(c.registry.GetType(c.GetLabeledStatement("MOD_A").enclosingScope, "exports.B")).
		IsNotNull()
}

// Java: TypedScopeCreatorTest.testGoogRequireDefaultExport_getsInferredType()
func (c *TypedScopeCreatorTest) TestGoogRequireDefaultExport_getsInferredType() {
	// Do some indirection when assigning `exports = f()` so  that we don't find the type of
	// `exports` until flow-sensitive inference runs. This is to verify that required variables
	// types are tightened after flow-sensitive type inference.
	c.TestSame(
		c.Srcs(
			"goog.module('a'); /** @return {number} */ function f() {} exports = f();",
			"goog.module('b'); const x = goog.require('a'); X: x;"))

	xNode := c.GetLabeledStatement("X").statementNode.GetOnlyChild()
	AssertNode(xNode).HasJSTypeThat().IsNumber()

	xVar := typedVarOf(c.GetLabeledStatement("X").enclosingScope.GetSlot("x"))
	TypeSubject_AssertType(xVar.GetType()).IsNumber()
	TypedVarSubject_AssertThat(xVar).IsInferred()
}

// Java: TypedScopeCreatorTest.testGoogRequire_defaultExportWithExplicitType()
func (c *TypedScopeCreatorTest) TestGoogRequire_defaultExportWithExplicitType() {
	c.TestSame(
		c.Srcs(
			"goog.module('a');\n"+
				"/** @return {number} */ function f() {}\n"+
				"/** @type {number} */ exports = f();\n",
			"goog.module('b'); const x = goog.require('a'); X: x;"))

	xNode := c.GetLabeledStatement("X").statementNode.GetOnlyChild()
	AssertNode(xNode).HasJSTypeThat().IsNumber()

	xVar := typedVarOf(c.GetLabeledStatement("X").enclosingScope.GetSlot("x"))
	TypedVarSubject_AssertThat(xVar).HasJSTypeThat().IsNumber()
	TypedVarSubject_AssertThat(xVar).IsNotInferred()
}

// Java: TypedScopeCreatorTest.testGoogRequire_propertyOnObjectLiteralWithArrayInitializer()
func (c *TypedScopeCreatorTest) TestGoogRequire_propertyOnObjectLiteralWithArrayInitializer() {
	// Reproduce a bug seen when goog.requiring symbols from
	// https://github.com/google/closure-library/blob/master/closure/goog/i18n/datetimesymbols.js
	c.TestSame(
		c.Srcs(
			"/** @fileoverview @typeSummary */\n"+
				"goog.provide('goog.i18n.DateTimeSymbols');\n"+
				"goog.provide('goog.i18n.DateTimeSymbolsType');\n"+
				"/** @type {!goog.i18n.DateTimeSymbolsType} */\n"+
				"goog.i18n.DateTimeSymbols;\n"+
				"\n"+
				"/** @typedef {{MONTHS: !Array<string>}} */\n"+
				"goog.i18n.DateTimeSymbolsType;\n",
			"goog.module('b');\n"+
				"const {MONTHS} = goog.require('goog.i18n.DateTimeSymbols');\n"+
				"function i() { MONTHS: MONTHS; }\n"))

	monthsNode := c.GetLabeledStatement("MONTHS").statementNode.GetOnlyChild()
	AssertNode(monthsNode).HasJSTypeThat().ToStringIsEqualTo("Array<string>")

	monthsVar := typedVarOf(c.GetLabeledStatement("MONTHS").enclosingScope.GetSlot("MONTHS"))
	TypedVarSubject_AssertThat(monthsVar).HasJSTypeThat().ToStringIsEqualTo("Array<string>")
	TypedVarSubject_AssertThat(monthsVar).IsNotInferred()
}

// Java: TypedScopeCreatorTest.testGoogForwardDeclare_defaultExportWithExplicitType()
func (c *TypedScopeCreatorTest) TestGoogForwardDeclare_defaultExportWithExplicitType() {
	c.TestSame(
		c.Srcs(
			"goog.module('b'); const x = goog.forwardDeclare('a'); function inner() { X: x; }",
			"goog.module('a');\n"+
				"/** @return {number} */ function f() {}\n"+
				"/** @type {number} */ exports = f();\n"))

	xNode := c.GetLabeledStatement("X").statementNode.GetOnlyChild()
	AssertNode(xNode).HasJSTypeThat().IsNumber()

	xVar := typedVarOf(c.GetLabeledStatement("X").enclosingScope.GetSlot("x"))
	TypedVarSubject_AssertThat(xVar).HasJSTypeThat().IsNumber()
	TypedVarSubject_AssertThat(xVar).IsNotInferred()
}

// Java: TypedScopeCreatorTest.testGoogRequireType_namedExportOfClass()
func (c *TypedScopeCreatorTest) TestGoogRequireType_namedExportOfClass() {
	c.processClosurePrimitives = true
	c.TestSame(
		c.Srcs(
			"goog.module('b');\n"+
				"const x = goog.requireType('a.x');\n"+
				"function inner() { var /** !x.Y */ y; Y: y; }\n",
			"goog.module('a.x');\n"+
				"exports.Y = class {};\n"))

	xNode := c.GetLabeledStatement("Y").statementNode.GetOnlyChild()
	AssertNode(xNode).HasJSTypeThat().ToStringIsEqualTo("a.x.Y")
}

// Java: TypedScopeCreatorTest.testGoogForwardDeclare_providedClass()
func (c *TypedScopeCreatorTest) TestGoogForwardDeclare_providedClass() {
	c.processClosurePrimitives = true
	c.TestSame(
		c.Srcs(
			"goog.module('b');\n"+
				"const X = goog.forwardDeclare('a.X');\n"+
				"function inner() { var /** !X */ x; X: x; }\n",
			"goog.provide('a.X');\n"+
				"a.X = class {};\n"))

	xNode := c.GetLabeledStatement("X").statementNode.GetOnlyChild()
	AssertNode(xNode).HasJSTypeThat().ToStringIsEqualTo("a.X")
}

// Java: TypedScopeCreatorTest.testGoogRequireType_providedClass()
func (c *TypedScopeCreatorTest) TestGoogRequireType_providedClass() {
	c.processClosurePrimitives = true
	c.TestSame(
		c.Srcs(
			"goog.module('b');\n"+
				"const X = goog.requireType('a.X');\n"+
				"function inner() { var /** !X */ x; X: x; }\n",
			"goog.provide('a.X');\n"+
				"a.X = class {};\n"))

	xNode := c.GetLabeledStatement("X").statementNode.GetOnlyChild()
	AssertNode(xNode).HasJSTypeThat().ToStringIsEqualTo("a.X")
}

// Java: TypedScopeCreatorTest.testGoogRequireType_classOnProvidedNamespace()
func (c *TypedScopeCreatorTest) TestGoogRequireType_classOnProvidedNamespace() {
	c.processClosurePrimitives = true
	c.TestSame(
		c.Srcs(
			"goog.module('b');\n"+
				"const x = goog.requireType('a.x');\n"+
				"function inner() { var /** !x.Y */ y; Y: y; }\n",
			"goog.provide('a.x');\n"+
				"a.x.Y = class {};\n"))

	xNode := c.GetLabeledStatement("Y").statementNode.GetOnlyChild()
	AssertNode(xNode).HasJSTypeThat().ToStringIsEqualTo("a.x.Y")
}

// Java: TypedScopeCreatorTest.testGoogRequireType_destructuring_classOnProvidedNamespace()
func (c *TypedScopeCreatorTest) TestGoogRequireType_destructuring_classOnProvidedNamespace() {
	c.processClosurePrimitives = true
	c.TestSame(
		c.Srcs(
			"goog.module('b');\n"+
				"const {Y} = goog.requireType('a.x');\n"+
				"function inner() { var /** !Y */ y; Y: y; }\n",
			"goog.provide('a.x');\n"+
				"a.x.Y = class {};\n"))

	xNode := c.GetLabeledStatement("Y").statementNode.GetOnlyChild()
	AssertNode(xNode).HasJSTypeThat().ToStringIsEqualTo("a.x.Y")
}

// Java: TypedScopeCreatorTest.testGoogRequire_defaultExportOfConstructor()
func (c *TypedScopeCreatorTest) TestGoogRequire_defaultExportOfConstructor() {
	c.TestSame(
		c.Srcs(
			"goog.module('a.Foo'); exports = class Foo {};",
			"goog.module('b'); const Foo = goog.require('a.Foo'); var /** !Foo */ foo; FOO: foo"))

	xNode := c.GetLabeledStatement("FOO").statementNode.GetOnlyChild()
	AssertNode(xNode).HasJSTypeThat().ToStringIsEqualTo("a.Foo")
}

// Java: TypedScopeCreatorTest.testGoogModuleRequire_defaultExportAliasingLocalConstructor()
func (c *TypedScopeCreatorTest) TestGoogModuleRequire_defaultExportAliasingLocalConstructor() {
	c.TestSame(
		c.Srcs(
			"goog.module('a.Foo'); class LocalFoo {} exports = LocalFoo;",
			"goog.module('b'); const Foo = goog.require('a.Foo'); var /** !Foo */ foo; FOO: foo"))

	fooNode := c.GetLabeledStatement("FOO").statementNode.GetOnlyChild()
	AssertNode(fooNode).HasJSTypeThat().ToStringIsEqualTo("LocalFoo")

	fooCtorVar := c.GetLabeledStatement("FOO").enclosingScope.GetVar("Foo")
	TypedVarSubject_AssertThat(fooCtorVar).HasJSTypeThat().IsFunctionTypeThat().IsConstructorFor("LocalFoo")
	TypedVarSubject_AssertThat(fooCtorVar).IsNotInferred()
}

// Java: TypedScopeCreatorTest.testGoogModuleRequire_defaultExportOfLocalTypedef()
func (c *TypedScopeCreatorTest) TestGoogModuleRequire_defaultExportOfLocalTypedef() {
	c.TestSame(
		c.Srcs(
			"goog.module('a.Foo'); /** @typedef {number} */ var numType; exports = numType;",
			"goog.module('b'); const Foo = goog.require('a.Foo'); var /** !Foo */ x; X: x"))

	xNode := c.GetLabeledStatement("X").statementNode.GetOnlyChild()
	AssertNode(xNode).HasJSTypeThat().IsNumber()
}

// Java: TypedScopeCreatorTest.testGoogModuleRequire_defaultExportOfEnum()
func (c *TypedScopeCreatorTest) TestGoogModuleRequire_defaultExportOfEnum() {
	c.TestSame(
		c.Srcs(
			"goog.module('a.Foo'); exports = class Foo {};",
			"goog.module('b'); const Foo = goog.require('a.Foo'); var /** !Foo */ foo; FOO: foo"))

	xNode := c.GetLabeledStatement("FOO").statementNode.GetOnlyChild()
	AssertNode(xNode).HasJSTypeThat().ToStringIsEqualTo("a.Foo")
}

// Java: TypedScopeCreatorTest.testGoogModuleRequireAndRequireType_defaultExportOfClasses()
func (c *TypedScopeCreatorTest) TestGoogModuleRequireAndRequireType_defaultExportOfClasses() {
	c.TestSame(
		c.Srcs(
			"goog.module('a.Foo');\n"+
				"const Bar = goog.requireType('b.Bar');\n"+
				"var /** !Bar */ b;\n"+
				"B: b;\n"+
				"exports = class {};\n",
			"goog.module('b.Bar');\n"+
				"const Foo = goog.require('a.Foo');\n"+
				"const /** !Foo */ f = new Foo();\n"+
				"F: f;\n"+
				"exports = class {};\n"))

	fNode := c.GetLabeledStatement("F").statementNode.GetOnlyChild()
	AssertNode(fNode).HasJSTypeThat().ToStringIsEqualTo("a.Foo")

	bNode := c.GetLabeledStatement("B").statementNode.GetOnlyChild()
	AssertNode(bNode).HasJSTypeThat().ToStringIsEqualTo("b.Bar")
}

// Java: TypedScopeCreatorTest.testGoogModuleRequireAndRequireType_namedExportOfClasses()
func (c *TypedScopeCreatorTest) TestGoogModuleRequireAndRequireType_namedExportOfClasses() {
	c.TestSame(
		c.Srcs(
			"goog.module('a');\n"+
				"const {Bar} = goog.requireType('b');\n"+
				"var /** !Bar */ b;\n"+
				"B: b;\n"+
				"exports.Foo = class {};\n",
			"goog.module('b');\n"+
				"const {Foo} = goog.require('a');\n"+
				"const /** !Foo */ f = new Foo();\n"+
				"F: f;\n"+
				"exports.Bar = class {};\n"))

	fNode := c.GetLabeledStatement("F").statementNode.GetOnlyChild()
	AssertNode(fNode).HasJSTypeThat().ToStringIsEqualTo("a.Foo")

	bNode := c.GetLabeledStatement("B").statementNode.GetOnlyChild()
	AssertNode(bNode).HasJSTypeThat().ToStringIsEqualTo("b.Bar")
}

// Java: TypedScopeCreatorTest.testGoogModuleRequireAndRequireType_typedef()
func (c *TypedScopeCreatorTest) TestGoogModuleRequireAndRequireType_typedef() {
	c.TestSame(
		c.Srcs(
			"goog.module('a.Foo');\n"+
				"const Bar = goog.requireType('b.Bar');\n"+
				"var /** !Bar */ x;\n"+
				"X: x;\n",
			"goog.module('b.Bar');\n"+
				"/** @typedef {number} */\n"+
				"let numType;\n"+
				"exports = numType;\n"))

	fNode := c.GetLabeledStatement("X").statementNode.GetOnlyChild()
	AssertNode(fNode).HasJSTypeThat().IsNumber()

	barDeclaration := c.GetLabeledStatement("X").enclosingScope.GetVar("Bar").GetNameNode()
	TypeSubject_AssertType(barDeclaration.GetTypedefTypeProp()).IsNumber()
}

// Java: TypedScopeCreatorTest.testGoogModuleRequireAndRequireType_invalidDestructuringImport()
func (c *TypedScopeCreatorTest) TestGoogModuleRequireAndRequireType_invalidDestructuringImport() {
	c.TestError(
		c.Srcs(
			"goog.module('a.Foo');\n"+
				"const {Bar} = goog.requireType('b.Bar');\n",
			"goog.module('b.Bar');"),
		c.Error(ModuleMapCreator_DOES_NOT_HAVE_EXPORT_WITH_DETAILS))
}

// Java: TypedScopeCreatorTest.testReferenceGoogModulesByType()
func (c *TypedScopeCreatorTest) TestReferenceGoogModulesByType() {
	c.TestSame(
		c.Srcs(
			"goog.module('mod.A');\n"+
				"class A {}\n"+
				"exports = A;\n",
			"var /** !mod.A */ a;"))

	TypedVarSubject_AssertThat(c.globalScope.GetVar("a")).HasJSTypeThat().ToStringIsEqualTo("A")
}

// Java: TypedScopeCreatorTest.testCannotReferenceGoogDeclareModuleIdByType()
func (c *TypedScopeCreatorTest) TestCannotReferenceGoogDeclareModuleIdByType() {
	c.TestWarning(
		c.Srcs(
			"goog.declareModuleId('mod');\n"+
				"export class A {}\n",
			"var /** !mod.A */ a;"),
		c.Warning(RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR))
}

// Java: TypedScopeCreatorTest.testReferenceGoogModuleNestedClassByType()
func (c *TypedScopeCreatorTest) TestReferenceGoogModuleNestedClassByType() {
	c.TestSame(
		c.Srcs(
			"goog.module('mod.A');\n"+
				"class A {}\n"+
				"A.B = class {};\n"+
				"A.B.C = class {};\n"+
				"exports = A;\n",
			"var /** !mod.A.B.C */ a;"))

	TypedVarSubject_AssertThat(c.globalScope.GetVar("a")).HasJSTypeThat().ToStringIsEqualTo("A.B.C")
}

// Java: TypedScopeCreatorTest.testReferenceGoogModuleNestedClassByType_exportDefaultTypedef()
func (c *TypedScopeCreatorTest) TestReferenceGoogModuleNestedClassByType_exportDefaultTypedef() {
	// Verify that the unusal resolution of @typedefs via Node annotation works correctly.
	c.TestSame(
		c.Srcs(
			"goog.module('mod.A');\n"+
				"/** @typedef {string} */\n"+
				"let MyTypedef;\n"+
				"exports = MyTypedef;\n",
			"var /** !mod.A */ a;"))

	TypedVarSubject_AssertThat(c.globalScope.GetVar("a")).HasJSTypeThat().IsString()
}

// Java: TypedScopeCreatorTest.testReferenceGoogModuleByType_shadowingLocalCausingResolutionFailure()
func (c *TypedScopeCreatorTest) TestReferenceGoogModuleByType_shadowingLocalCausingResolutionFailure() {
	c.TestWarning(
		c.Srcs(
			"goog.module('mod.A');\n"+
				"class A {}\n"+
				"A.B = class {};\n"+
				"A.B.C = class {};\n"+
				"exports = A;\n",
			"goog.module('mod.B');\n"+
				"const mod = {};\n"+
				"var /** !mod.A.B.C */ a;\n"), // Resolves relative to the local `mod`.
		c.Warning(RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR))
}

// Java: TypedScopeCreatorTest.testReferenceGoogModuleByType_googProvideExtensionCausesResolutionFailure()
func (c *TypedScopeCreatorTest) TestReferenceGoogModuleByType_googProvideExtensionCausesResolutionFailure() {
	c.processClosurePrimitives = true
	c.TestWarning(
		c.Srcs(
			"goog.module('mod.A');\n"+
				"class A {}\n"+
				"A.B = class {};\n"+
				"A.B.C = class {};\n"+
				"exports = A;\n",
			"goog.provide('mod.A.B');\n"+
				"var /** !mod.A.B.C */ a;\n"), // Assumes this refers to the goog.provide
		c.Warning(RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR))
}

// Java: TypedScopeCreatorTest.testGoogRequireModuleWithoutExports()
func (c *TypedScopeCreatorTest) TestGoogRequireModuleWithoutExports() {
	c.TestSame(c.Srcs("goog.module('a');", "goog.module('b'); const a = goog.require('a');"))
}

// Java: TypedScopeCreatorTest.testGoogRequire_exportObjectLiteralWithLiteralValue()
func (c *TypedScopeCreatorTest) TestGoogRequire_exportObjectLiteralWithLiteralValue() {
	c.TestSame(
		c.Srcs(
			"goog.module('a'); exports = {x: 3};",
			"goog.module('b'); const x = goog.require('a'); X: x;", // can't destructure this
		))

	xNode := c.GetLabeledStatement("X").statementNode.GetOnlyChild()
	AssertNode(xNode).HasJSTypeThat().ToStringIsEqualTo("{x: number}")
}

// Java: TypedScopeCreatorTest.testGoogRequire_exportObjectLiteralWithLiteralFunctions()
func (c *TypedScopeCreatorTest) TestGoogRequire_exportObjectLiteralWithLiteralFunctions() {
	c.TestSame(
		c.Srcs(
			"goog.module('my.mod');\n"+
				"function memoize(fn) { return fn; }\n"+
				"exports = {\n"+
				"  /** @return {string} */\n"+
				"  bar: memoize(() => 'test')\n"+
				"};\n",
			"goog.module('other.mod');\n"+
				"const mod = goog.require('my.mod');\n"+
				"MOD: mod;\n"))

	moduleObject := c.GetLabeledStatement("MOD").statementNode.GetOnlyChild()
	AssertNode(moduleObject).
		HasJSTypeThat().
		WithTypeOfProp("bar").
		IsFunctionTypeThat().
		HasReturnTypeThat().
		IsString()
}

// Java: TypedScopeCreatorTest.testModuleTypedef_isNonnullableWithinModule()
func (c *TypedScopeCreatorTest) TestModuleTypedef_isNonnullableWithinModule() {
	c.TestSame(
		"goog.module('a.b'); var /** strType */ earlyRef;\n" +
			"/** @typedef {string} */ let strType;\n" +
			"var /** strType */ normalRef; MODULE: 0;\n")

	moduleScope := c.GetLabeledStatement("MODULE").enclosingScope
	TypedVarSubject_AssertThat(typedVarOf(moduleScope.GetSlot("earlyRef"))).HasJSTypeThat().IsString()
	TypedVarSubject_AssertThat(typedVarOf(moduleScope.GetSlot("normalRef"))).HasJSTypeThat().IsString()
}

// Java: TypedScopeCreatorTest.testModuleTypedef_isNonnullableWhenRequired()
func (c *TypedScopeCreatorTest) TestModuleTypedef_isNonnullableWhenRequired() {

	c.TestSame(
		c.Srcs(
			"goog.module('root');\n"+
				"/** @typedef {string} */ let strType;\n"+
				"exports.strType = strType;\n",
			"goog.module('a.b'); const {strType} = goog.require('root');\n"+
				"/** @param {strType} earlyRef */ function f(earlyRef) { EARLY_REF: earlyRef; }\n"+
				"var /** strType */ normalRef; NORMAL_REF: normalRef;\n"))

	AssertNode(c.GetLabeledStatement("EARLY_REF").statementNode.GetOnlyChild()).
		HasJSTypeThat().
		IsString()
	AssertNode(c.GetLabeledStatement("NORMAL_REF").statementNode.GetOnlyChild()).
		HasJSTypeThat().
		IsString()
}

// Java: TypedScopeCreatorTest.testGoogRequire_namedExportImportedAsNamespace()
func (c *TypedScopeCreatorTest) TestGoogRequire_namedExportImportedAsNamespace() {
	c.TestSame(
		c.Srcs(
			"goog.module('a'); /** @type {number} */ exports.x = 0;",
			"goog.module('b'); const a = goog.require('a'); X: a.x;"))

	xNode := c.GetLabeledStatement("X").statementNode.GetOnlyChild()
	AssertNode(xNode).HasJSTypeThat().IsNumber()

	xVar := typedVarOf(c.GetLabeledStatement("X").enclosingScope.GetSlot("a"))
	TypedVarSubject_AssertThat(xVar).HasJSTypeThat().IsObjectTypeWithProperty("x")
	TypedVarSubject_AssertThat(xVar).IsNotInferred()
}

// Java: TypedScopeCreatorTest.testGoogRequire_namedExportImportedAsNamespace_andModded()
func (c *TypedScopeCreatorTest) TestGoogRequire_namedExportImportedAsNamespace_andModded() {
	c.TestSame(
		c.Srcs(
			"goog.module('a'); class C { m() {} } exports = {C}; C: C;",
			"goog.module('b'); const a = goog.require('a'); a.C.prototype.n = function() {};"))

	cCtor := c.GetLabeledStatement("C").statementNode.GetOnlyChild()
	AssertNode(cCtor).HasJSTypeThat().IsFunctionTypeThat().IsConstructorFor("C")

	cCtorType := cCtor.GetJSType().AssertFunctionType()
	cPrototype := cCtorType.GetPrototype()
	TypeSubject_AssertType(cPrototype).HasDeclaredProperty("n")
	TypeSubject_AssertType(cPrototype.GetPropertyType("n")).
		IsFunctionTypeThat().
		HasTypeOfThisThat().
		IsEqualTo(cCtorType.GetInstanceType())
}

// Java: TypedScopeCreatorTest.testGoogRequire_destructuringInferredNamedExport()
func (c *TypedScopeCreatorTest) TestGoogRequire_destructuringInferredNamedExport() {
	c.TestSame(
		c.Srcs(
			"goog.module('a'); /** @return {number} */ function f() {} exports.x = f();",
			"goog.module('b'); const {x} = goog.require('a'); X: x;"))

	xNode := c.GetLabeledStatement("X").statementNode.GetOnlyChild()
	AssertNode(xNode).HasJSTypeThat().IsNumber()
}

// Java: TypedScopeCreatorTest.testGoogRequire_destructuringTypedNamedExport()
func (c *TypedScopeCreatorTest) TestGoogRequire_destructuringTypedNamedExport() {
	c.TestSame(
		c.Srcs(
			"goog.module('a'); /** @type {number} */ exports.x = 3;",
			"goog.module('b'); const {x} = goog.require('a'); X: x;"))

	xNode := c.GetLabeledStatement("X").statementNode.GetOnlyChild()
	AssertNode(xNode).HasJSTypeThat().IsNumber()

	xVar := typedVarOf(c.GetLabeledStatement("X").enclosingScope.GetSlot("x"))
	TypedVarSubject_AssertThat(xVar).HasJSTypeThat().IsNumber()
	TypedVarSubject_AssertThat(xVar).IsNotInferred()
}

// Java: TypedScopeCreatorTest.testGoogRequire_destructuringNonShorthand()
func (c *TypedScopeCreatorTest) TestGoogRequire_destructuringNonShorthand() {
	c.TestSame(
		c.Srcs(
			"goog.module('a'); /** @type {number} */ exports.x = 3;",
			"goog.module('b'); const {x: y} = goog.require('a'); Y: y;"))

	yNode := c.GetLabeledStatement("Y").statementNode.GetOnlyChild()
	AssertNode(yNode).HasJSTypeThat().IsNumber()
}

// Java: TypedScopeCreatorTest.testGoogRequire_exportObjectLiteral()
func (c *TypedScopeCreatorTest) TestGoogRequire_exportObjectLiteral() {
	c.TestSame(
		c.Srcs(
			"goog.module('a'); const y = 3; exports = {x: y};",
			"goog.module('b'); const {x} = goog.require('a'); X: x;"))

	xNode := c.GetLabeledStatement("X").statementNode.GetOnlyChild()
	AssertNode(xNode).HasJSTypeThat().IsNumber()
}

// Java: TypedScopeCreatorTest.testGoogRequire_namedExportDeclaredAsConstructor()
func (c *TypedScopeCreatorTest) TestGoogRequire_namedExportDeclaredAsConstructor() {
	c.TestSame(
		c.Srcs(
			"goog.module('a'); /** @constructor */ exports.X = function() {};",
			"goog.module('b'); const {X} = goog.require('a'); var /** !X */ x; X: x;"))

	xNode := c.GetLabeledStatement("X").statementNode.GetOnlyChild()
	AssertNode(xNode).HasJSTypeThat().ToStringIsEqualTo("a.X")
}

// Java: TypedScopeCreatorTest.testGoogRequire_namedExportDeclaredAsEnum()
func (c *TypedScopeCreatorTest) TestGoogRequire_namedExportDeclaredAsEnum() {
	c.TestSame(
		c.Srcs(
			"goog.module('a'); /** @enum {number} */ exports.Enum = {A: 0};",
			"goog.module('b'); const {Enum} = goog.require('a'); var /** !Enum */ x; X: x;"))

	xNode := c.GetLabeledStatement("X").statementNode.GetOnlyChild()
	AssertNode(xNode).HasJSTypeThat().ToStringIsEqualTo("a.Enum<number>")
}

// Java: TypedScopeCreatorTest.testGoogRequire_namedExportOfLocalEnum()
func (c *TypedScopeCreatorTest) TestGoogRequire_namedExportOfLocalEnum() {
	c.TestSame(
		c.Srcs(
			"goog.module('a'); /** @enum {number} */ const Enum = {A: 0}; exports.Enum = Enum;",
			"goog.module('b'); const {Enum} = goog.require('a'); var /** !Enum */ x; X: x"))

	xNode := c.GetLabeledStatement("X").statementNode.GetOnlyChild()
	AssertNode(xNode).HasJSTypeThat().ToStringIsEqualTo("Enum<number>")
}

// Java: TypedScopeCreatorTest.testGoogRequire_namedExportOfLocalEnum_inObjectLit()
func (c *TypedScopeCreatorTest) TestGoogRequire_namedExportOfLocalEnum_inObjectLit() {
	c.TestSame(
		c.Srcs(
			"goog.module('a'); /** @enum {number} */ const Enum = {A: 0}; exports = {Enum}",
			"goog.module('b'); const {Enum} = goog.require('a'); var /** !Enum */ x; X: x"))

	xNode := c.GetLabeledStatement("X").statementNode.GetOnlyChild()
	AssertNode(xNode).HasJSTypeThat().ToStringIsEqualTo("Enum<number>")
}

// Java: TypedScopeCreatorTest.testGoogRequire_namedExportDeclaredAsTypedef()
func (c *TypedScopeCreatorTest) TestGoogRequire_namedExportDeclaredAsTypedef() {
	c.TestSame(
		c.Srcs(
			"goog.module('a'); /** @typedef {number} */ exports.numType;",
			"goog.module('b'); const {numType} = goog.require('a'); var /** numType */ x; X: x;"))

	xNode := c.GetLabeledStatement("X").statementNode.GetOnlyChild()
	AssertNode(xNode).HasJSTypeThat().IsNumber()
}

// Java: TypedScopeCreatorTest.testGoogRequire_namedExportOfLocalTypedef()
func (c *TypedScopeCreatorTest) TestGoogRequire_namedExportOfLocalTypedef() {
	c.TestSame(
		c.Srcs(
			"goog.module('a'); /** @typedef {number} */ let numType; exports.numType = numType;",
			"goog.module('b'); const {numType} = goog.require('a'); var /** numType */ x; X: x;"))

	xNode := c.GetLabeledStatement("X").statementNode.GetOnlyChild()
	AssertNode(xNode).HasJSTypeThat().IsNumber()
}

// Java: TypedScopeCreatorTest.testGoogRequire_destructuringMissingExportDoesntCrash()
func (c *TypedScopeCreatorTest) TestGoogRequire_destructuringMissingExportDoesntCrash() {
	c.TestError(
		c.Srcs("goog.module('a');", "goog.module('b'); const {x} = goog.require('a'); X: x;"),
		c.Error(ModuleMapCreator_DOES_NOT_HAVE_EXPORT_WITH_DETAILS))

	xNode := c.GetLabeledStatement("X").statementNode.GetOnlyChild()
	AssertNode(xNode).HasJSTypeThat().IsUnknown() // Note: we warn for this case elsewhere.
}

// Java: TypedScopeCreatorTest.testRequire_requireRegularModuleInLoadModule()
func (c *TypedScopeCreatorTest) TestRequire_requireRegularModuleInLoadModule() {
	c.TestSame(
		c.Srcs(
			"goog.module('A'); exports = class {};",
			"goog.loadModule(function(exports) {\n"+
				"  goog.module('b');\n"+
				"  const A = goog.require('A');\n"+
				"  /** @type {!A} */ var a;\n"+
				"  A: a;\n"+
				"  return exports;\n"+
				"});\n"))

	aNode := c.GetLabeledStatement("A").statementNode.GetOnlyChild()
	AssertNode(aNode).HasJSTypeThat().ToStringIsEqualTo("A")
}

// Java: TypedScopeCreatorTest.testGoogRequire_requireLoadModuleInRegularModule()
func (c *TypedScopeCreatorTest) TestGoogRequire_requireLoadModuleInRegularModule() {
	c.TestSame(
		c.Srcs(
			"goog.loadModule(function(exports) {\n"+
				"  goog.module('A');\n"+
				"  /** @constructor*/\n"+
				"  exports = function() {};\n"+
				" return exports;\n"+
				"});\n",
			"goog.module('b'); const A = goog.require('A'); /** @type {!A} */ var a; A: a;"))

	aNode := c.GetLabeledStatement("A").statementNode.GetOnlyChild()
	AssertNode(aNode).HasJSTypeThat().ToStringIsEqualTo("A")
}

// Java: TypedScopeCreatorTest.testGoogRequire_requireLoadModuleInRegularModule_destructuring()
func (c *TypedScopeCreatorTest) TestGoogRequire_requireLoadModuleInRegularModule_destructuring() {
	c.TestSame(
		c.Srcs(
			"goog.loadModule(function(exports) {\n"+
				"  goog.module('mod');\n"+
				"  /** @constructor*/\n"+
				"  exports.A = function() {};\n"+
				"  return exports;\n"+
				"});\n",
			"goog.module('b'); const {A} = goog.require('mod'); /** @type {!A} */ var a; A: a;"))

	aNode := c.GetLabeledStatement("A").statementNode.GetOnlyChild()
	AssertNode(aNode).HasJSTypeThat().ToStringIsEqualTo("mod.A")
}

// Java: TypedScopeCreatorTest.testRequire_requireLoadModuleInRegularModule_inferredDefaultExport()
func (c *TypedScopeCreatorTest) TestRequire_requireLoadModuleInRegularModule_inferredDefaultExport() {
	c.TestSame(
		c.Srcs(
			"goog.loadModule(function(exports) {\n"+
				"  goog.module('a');\n"+
				"  /** @return {number} */ function f() {}\n"+
				"  exports = f();\n"+
				"  return exports;\n"+
				"});\n",
			"goog.module('b'); const x = goog.require('a'); X: x;"))

	// This is unknown because we visit the CFG for the module before visiting the CFG for the
	// function block. Users can work around this by explicitly typing their exports.
	xNode := c.GetLabeledStatement("X").statementNode.GetOnlyChild()
	AssertNode(xNode).HasJSTypeThat().IsUnknown()
}

// Java: TypedScopeCreatorTest.testRequire_requireLoadModuleInRegularModule_inferredNamedExport()
func (c *TypedScopeCreatorTest) TestRequire_requireLoadModuleInRegularModule_inferredNamedExport() {
	c.TestSame(
		c.Srcs(
			"goog.loadModule(function(exports) {\n"+
				"  goog.module('a');\n"+
				"  /** @return {number} */ function f() {}\n"+
				"  exports.x = f();\n"+
				"  return exports;\n"+
				"});\n",
			"goog.module('b'); const {x} = goog.require('a'); X: x;"))

	// This is unknown because we visit the CFG for the module before visiting the CFG for the
	// function block. Users can work around this by explicitly typing their exports.
	xNode := c.GetLabeledStatement("X").statementNode.GetOnlyChild()
	AssertNode(xNode).HasJSTypeThat().IsUnknown()
}

// Java: TypedScopeCreatorTest.testGoogModuleGet_nameDeclarationCreatesTypeAlias()
func (c *TypedScopeCreatorTest) TestGoogModuleGet_nameDeclarationCreatesTypeAlias() {
	c.TestSame(
		c.Srcs(
			"(function() {\n"+
				"  const B = goog.module.get('a.B');\n"+
				"  var /** !B */ b;\n"+
				"  B: b;\n"+
				"})();\n",
			"goog.module('a.B');\n"+
				"class B {}\n"+
				"exports = B;\n"))

	TypeSubject_AssertType(c.GetLabeledStatement("B").statementNode.GetOnlyChild().GetJSType()).
		ToStringIsEqualTo("B")
}

// Java: TypedScopeCreatorTest.testGoogModuleGet_destructruringCreatesTypeAlias()
func (c *TypedScopeCreatorTest) TestGoogModuleGet_destructruringCreatesTypeAlias() {
	c.TestSame(
		c.Srcs(
			"(function() {\n"+
				"  const {B} = goog.module.get('a');\n"+
				"  var /** !B */ b;\n"+
				"  B: b;\n"+
				"})();\n",
			"goog.module('a');\n"+
				"/** @typedef {number} */ let B;\n"+
				"exports.B = B;\n"))

	TypeSubject_AssertType(c.GetLabeledStatement("B").statementNode.GetOnlyChild().GetJSType()).IsNumber()
}

// Java: TypedScopeCreatorTest.testGoogModuleGet_getpropDoesNotCreateTypeAlias()
func (c *TypedScopeCreatorTest) TestGoogModuleGet_getpropDoesNotCreateTypeAlias() {
	// This test covers a pattern that works when the compiler rewrites goog.module.get before
	// typechecking, but is not supported in native module typechecking. It's a rare pattern and
	// would be more work to support in TypedScopeCreator.
	c.TestWarning(
		c.Srcs(
			"(function() {\n"+
				"  const B = goog.module.get('a').B;\n"+
				"  var /** !B */ b;\n"+
				"  B: b;\n"+
				"})();\n",
			"goog.module('a');\n"+
				"/** @typedef {number} */ let B;\n"+
				"exports.B = B;\n"),
		c.Warning(RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR))

	TypeSubject_AssertType(c.GetLabeledStatement("B").statementNode.GetOnlyChild().GetJSType()).IsUnknown()
}

// Java: TypedScopeCreatorTest.testRequire_requireLoadModule_defaultExportOfLocalClass()
func (c *TypedScopeCreatorTest) TestRequire_requireLoadModule_defaultExportOfLocalClass() {
	c.TestSame(
		c.Srcs(
			"goog.loadModule(function(exports) {\n"+
				"  goog.module('a');\n"+
				"  class X {}\n"+
				"  exports = X;\n"+
				"  return exports;\n"+
				"});\n",
			"goog.module('b');\n"+
				"const Y = goog.require('a');\n"+
				"var /** !Y */ x;\n"))
}

// Java: TypedScopeCreatorTest.testGoogProvide_singleNameWithConstJSDocIsDeclared()
func (c *TypedScopeCreatorTest) TestGoogProvide_singleNameWithConstJSDocIsDeclared() {
	c.processClosurePrimitives = true
	c.TestSame(c.Srcs("goog.provide('foo'); /** @const */ foo = 3;"))

	fooVar := c.globalScope.GetVar("foo")
	TypedVarSubject_AssertThat(fooVar).HasJSTypeThat().IsNumber()
	TypedVarSubject_AssertThat(fooVar).IsNotInferred()
}

// Java: TypedScopeCreatorTest.testGoogProvide_singleNameWithAtTypeJSDocIsDeclared()
func (c *TypedScopeCreatorTest) TestGoogProvide_singleNameWithAtTypeJSDocIsDeclared() {
	c.processClosurePrimitives = true
	c.TestSame(
		c.Srcs(
			"goog.provide('foo');\n" +
				"class Bar {}\n" +
				"/** @type {!Bar} */ foo = new Bar();\n"))

	fooVar := c.globalScope.GetVar("foo")
	TypedVarSubject_AssertThat(fooVar).HasJSTypeThat().ToStringIsEqualTo("Bar")
	TypedVarSubject_AssertThat(fooVar).IsNotInferred()
}

// Java: TypedScopeCreatorTest.testGoogProvide_singleNameAssignedInIfBlock()
func (c *TypedScopeCreatorTest) TestGoogProvide_singleNameAssignedInIfBlock() {
	c.processClosurePrimitives = true
	c.TestSame(c.Srcs("goog.provide('foo'); if (!foo) { foo = {}; };"))

	fooVar := c.globalScope.GetVar("foo")
	TypedVarSubject_AssertThat(fooVar).HasJSTypeThat().ToStringIsEqualTo("{}")
	TypedVarSubject_AssertThat(fooVar).IsInferred()
}

// Java: TypedScopeCreatorTest.testGoogProvide_ignoresBlockScopedShadow()
func (c *TypedScopeCreatorTest) TestGoogProvide_ignoresBlockScopedShadow() {
	c.processClosurePrimitives = true
	c.TestSame(c.Srcs("goog.provide('foo'); if (true) { const foo = 3; }"))

	fooVar := c.globalScope.GetVar("foo")
	TypedVarSubject_AssertThat(fooVar).HasJSTypeThat().ToStringIsEqualTo("{}")
	TypedVarSubject_AssertThat(fooVar).IsNotInferred()
}

// Java: TypedScopeCreatorTest.testGoogRequire_ofProvide_usedInEs6ExtendsClause()
func (c *TypedScopeCreatorTest) TestGoogRequire_ofProvide_usedInEs6ExtendsClause() {
	c.TestSame(
		c.Srcs(
			"goog.provide('a.b')\n"+
				"a.b.Foo = class {};\n",
			"goog.module('c');\n"+
				"const b = goog.require('a.b');\n"+
				"class C extends b.Foo {}\n"))

	cType := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("C", c.lastLocalScope))
	aFooType := typedScopeCreatorTest_cast[rhino.FunctionType](c.FindNameType("a.b.Foo", c.globalScope))

	TypeSubject_AssertType(cType.GetInstanceType()).IsSubtypeOf(aFooType.GetInstanceType())
	TypeSubject_AssertType(cType.GetImplicitPrototype()).IsEqualTo(aFooType)
}

// Java: TypedScopeCreatorTest.testGoogProvide_singleStubNameWithAtTypeJSDoc_ignoresJSDoc()
func (c *TypedScopeCreatorTest) TestGoogProvide_singleStubNameWithAtTypeJSDoc_ignoresJSDoc() {
	c.processClosurePrimitives = true
	// NOTE: this causes a CheckJSDoc warning.
	c.TestSame(c.Srcs("goog.provide('foo'); class Bar {} /** @type {!Bar} */ foo;"))

	fooVar := c.globalScope.GetVar("foo")
	TypedVarSubject_AssertThat(fooVar).HasJSTypeThat().IsLiteralObject()
	TypedVarSubject_AssertThat(fooVar).IsNotInferred()
}

// Java: TypedScopeCreatorTest.testGoogProvide_singleNameWithoutJSDocIsInferred()
func (c *TypedScopeCreatorTest) TestGoogProvide_singleNameWithoutJSDocIsInferred() {
	c.processClosurePrimitives = true
	c.TestSame(c.Srcs("goog.provide('foo'); foo = 3;"))

	fooVar := c.globalScope.GetVar("foo")
	TypedVarSubject_AssertThat(fooVar).HasJSTypeThat().IsNumber()
	TypedVarSubject_AssertThat(fooVar).IsInferred()
}

// Java: TypedScopeCreatorTest.testGoogProvide_singleNameWithoutDefinition()
func (c *TypedScopeCreatorTest) TestGoogProvide_singleNameWithoutDefinition() {
	c.processClosurePrimitives = true
	c.TestSame(c.Srcs("goog.provide('foo'); foo.Bar = class {};"))

	fooVar := c.globalScope.GetVar("foo")
	TypedVarSubject_AssertThat(fooVar).HasJSTypeThat().IsObjectTypeWithProperty("Bar")
	TypedVarSubject_AssertThat(fooVar).IsNotInferred()
}

// Java: TypedScopeCreatorTest.testGoogProvide_getPropWithoutJSDoc()
func (c *TypedScopeCreatorTest) TestGoogProvide_getPropWithoutJSDoc() {
	c.processClosurePrimitives = true
	c.TestSame(c.Srcs("goog.provide('foo.bar'); foo.bar = 3;"))

	fooVar := c.globalScope.GetVar("foo")
	TypedVarSubject_AssertThat(fooVar).HasJSTypeThat().WithTypeOfProp("bar").IsNumber()
	AssertThat(fooVar.GetType().ToMaybeObjectType().IsPropertyTypeDeclared("bar")).IsFalse()
}

// Java: TypedScopeCreatorTest.testGoogProvide_getPropTypedef_createsBothTypedefAndObjectType()
func (c *TypedScopeCreatorTest) TestGoogProvide_getPropTypedef_createsBothTypedefAndObjectType() {
	c.processClosurePrimitives = true
	c.TestSame(
		c.Srcs(
			"goog.provide('foo.bar');\n" +
				"goog.provide('foo.bar.baz');\n" +
				"/** @typedef {number} */ foo.bar;\n" +
				"/** @const */ foo.bar.baz = '';\n"))

	fooVar := c.globalScope.GetVar("foo")
	TypedVarSubject_AssertThat(fooVar).HasJSTypeThat().WithTypeOfProp("bar").IsNoType()
	TypedVarSubject_AssertThat(fooVar).HasJSTypeThat().HasDeclaredProperty("bar")

	TypeSubject_AssertType(c.registry.GetGlobalType("foo.bar")).IsNumber()
	fooBarVar := c.globalScope.GetVar("foo.bar")
	TypedVarSubject_AssertThat(fooBarVar).HasJSTypeThat().IsNoType() // Note: this is questionable.

	fooBarBaz := c.globalScope.GetVar("foo.bar.baz")
	TypedVarSubject_AssertThat(fooBarBaz).HasJSTypeThat().IsString()
}

// Java: TypedScopeCreatorTest.testGoogProvideGetPropWithTypeDeclaration()
func (c *TypedScopeCreatorTest) TestGoogProvideGetPropWithTypeDeclaration() {
	c.processClosurePrimitives = true
	c.TestSame(
		c.Srcs(
			"goog.provide('foo.bar');\n" +
				"/** @type {number} */ foo.bar = 3;\n" +
				"FOOBAR: foo.bar;\n"))

	fooVar := c.globalScope.GetVar("foo")
	TypedVarSubject_AssertThat(fooVar).HasJSTypeThat().WithTypeOfProp("bar").IsNumber()
	TypedVarSubject_AssertThat(fooVar).HasJSTypeThat().HasDeclaredProperty("bar")

	TypeSubject_AssertType(c.GetLabeledStatement("FOOBAR").statementNode.GetOnlyChild().GetJSType()).IsNumber()
}

// Java: TypedScopeCreatorTest.testGoogProvide_namespaceWithMultipleChildren()
func (c *TypedScopeCreatorTest) TestGoogProvide_namespaceWithMultipleChildren() {
	c.processClosurePrimitives = true
	c.TestSame(
		c.Srcs(
			"goog.provide('foo.bar');\n" +
				"goog.provide('foo.baz');\n" +
				"/** @const */\n" +
				"foo.bar = 3;\n" +
				"/** @const */\n" +
				"foo.baz = 'str';\n"))

	fooVar := c.globalScope.GetVar("foo")
	TypeSubject_AssertType(fooVar.GetType().FindPropertyType("bar")).IsNumber()
	TypedVarSubject_AssertThat(fooVar).HasJSTypeThat().HasDeclaredProperty("bar")
	TypeSubject_AssertType(fooVar.GetType().FindPropertyType("baz")).IsString()
	TypedVarSubject_AssertThat(fooVar).HasJSTypeThat().HasDeclaredProperty("baz")
}

// Java: TypedScopeCreatorTest.testProvide_longNamespace_declaresProperties()
func (c *TypedScopeCreatorTest) TestProvide_longNamespace_declaresProperties() {
	c.processClosurePrimitives = true
	c.TestSame(
		c.Srcs(
			"goog.provide('a.b.c.d.e.f');\n" +
				"\n" +
				"a.b.c.d.e.f.g = function() {};\n"))

	abcdeVar := c.globalScope.GetVar("a.b.c.d.e")
	TypedVarSubject_AssertThat(abcdeVar).HasJSTypeThat().HasDeclaredProperty("f")
}

// Java: TypedScopeCreatorTest.testGoogProvide_parentNamespaceIsInferred()
func (c *TypedScopeCreatorTest) TestGoogProvide_parentNamespaceIsInferred() {
	c.processClosurePrimitives = true
	c.TestSame(c.Srcs("goog.provide('a.b.c'); a.b = {};"))

	TypedVarSubject_AssertThat(c.globalScope.GetVar("a")).HasJSTypeThat().WithTypeOfProp("b").IsLiteralObject()
	TypedVarSubject_AssertThat(c.globalScope.GetVar("a.b")).IsNull() // Only declared qnames get TypedVars.
	TypedVarSubject_AssertThat(c.globalScope.GetVar("a.b.c")).IsNotInferred()
}

// Java: TypedScopeCreatorTest.testGoogProvide_extendsExternsNamespace()
func (c *TypedScopeCreatorTest) TestGoogProvide_extendsExternsNamespace() {
	c.processClosurePrimitives = true
	c.TestSame(
		c.Externs(
			// Go-only: language-adapter — Java passes the varargs as an array.
			NewTestExternsBuilder().AddClosureExterns().Build(),
			"/** @externs */ /** @const */ var google = {}; /** @type {string} */ google.boogle;"),
		c.Srcs(
			"goog.provide('google.ly');\n"+
				"\n"+
				"/** @return {number} */\n"+
				"google.ly = function() { return 0; }\n"),
		c.Warning(TypeValidator_DUP_VAR_DECLARATION_TYPE_MISMATCH))

	fooVar := c.globalScope.GetVar("google")
	TypedVarSubject_AssertThat(fooVar).HasJSTypeThat().HasDeclaredProperty("ly")
	TypedVarSubject_AssertThat(fooVar).HasJSTypeThat().HasDeclaredProperty("boogle")
}

// Java: TypedScopeCreatorTest.testGoogProvide_overwritingExternsFunction()
func (c *TypedScopeCreatorTest) TestGoogProvide_overwritingExternsFunction() {
	c.processClosurePrimitives = true
	c.TestSame(
		c.Externs(
			NewTestExternsBuilder().
				AddClosureExterns().
				AddExtra("/** @externs */ function eval(code) {}").
				Build()),
		c.Srcs(
			"goog.provide('eval.subNs');\n"+
				"eval.subNs = class {};\n"+
				"new eval.subNs();\n"),
		c.Warning(TypeValidator_DUP_VAR_DECLARATION_TYPE_MISMATCH))
}

// Java: TypedScopeCreatorTest.testLegacyGoogLoadModule_accessibleWithGoogRequire_exportingConstructor()
func (c *TypedScopeCreatorTest) TestLegacyGoogLoadModule_accessibleWithGoogRequire_exportingConstructor() {
	c.TestSame(
		c.Srcs(
			"goog.loadModule(function (exports) {\n"+
				"  goog.module('mod.A');\n"+
				"  goog.module.declareLegacyNamespace();\n"+
				"\n"+
				"  exports = class A {};\n"+
				"});\n",
			"goog.require('mod.A');\n"+
				"const /** !mod.A */ a = new mod.A();\n"+
				"A: a;\n"))

	TypeSubject_AssertType(c.GetLabeledStatement("A").statementNode.GetOnlyChild().GetJSType()).
		ToStringIsEqualTo("mod.A")
}

// Java: TypedScopeCreatorTest.testLegacyGoogModule_declaresParentNamespacesGlobally()
func (c *TypedScopeCreatorTest) TestLegacyGoogModule_declaresParentNamespacesGlobally() {
	c.processClosurePrimitives = true
	c.TestSame(c.Srcs("goog.module('a.b.c'); goog.module.declareLegacyNamespace(); MOD: 0;"))

	ScopeSubject_AssertScope(c.globalScope).Declares("a").Directly()
	ScopeSubject_AssertScope(c.globalScope).Declares("a.b").Directly()
	ScopeSubject_AssertScope(c.globalScope).Declares("a.b.c").Directly()
	ScopeSubject_AssertScope(c.GetLabeledStatement("MOD").enclosingScope).Declares("a").OnSomeParent()
}

// Java: TypedScopeCreatorTest.testLegacyGoogModule_isDeclaredPropertyOnParentNamespace()
func (c *TypedScopeCreatorTest) TestLegacyGoogModule_isDeclaredPropertyOnParentNamespace() {
	c.processClosurePrimitives = true
	c.TestSame(c.Srcs("goog.module('a.b.c'); goog.module.declareLegacyNamespace();"))

	TypedVarSubject_AssertThat(c.globalScope.GetVar("a.b")).HasJSTypeThat().HasDeclaredProperty("c")
}

// Java: TypedScopeCreatorTest.testLegacyGoogModule_addsJSDocToPropertyOnParentNamespace()
func (c *TypedScopeCreatorTest) TestLegacyGoogModule_addsJSDocToPropertyOnParentNamespace() {
	c.processClosurePrimitives = true
	c.TestSame(
		c.Srcs(
			"goog.module('a.b.c');\n" +
				"goog.module.declareLegacyNamespace();\n" +
				"/** @package */ exports = 0;\n"))

	TypedVarSubject_AssertThat(c.globalScope.GetVar("a.b")).HasJSTypeThat().HasDeclaredProperty("c")
	abType := c.globalScope.GetVar("a.b").GetType().ToMaybeObjectType()
	AssertThat(abType.GetPropertyJSDocInfo("c")).IsNotNull()
	AssertThat(abType.GetPropertyJSDocInfo("c").GetVisibility()).
		IsEqualTo(rhino.JSDocInfo_Visibility_PACKAGE)
}

// Java: TypedScopeCreatorTest.testGoogRequire_multipleRequiresInModule()
func (c *TypedScopeCreatorTest) TestGoogRequire_multipleRequiresInModule() {
	c.TestSame(
		c.Srcs(
			"goog.module('mod.A');\n"+
				"/** @interface */\n"+
				"exports = class {};\n",
			"goog.module('mod.B');\n"+
				"/** @interface */\n"+
				"exports = class {};\n",
			"goog.module('mod.C');\n"+
				"const A = goog.require('mod.A');\n"+
				"const B = goog.require('mod.B');\n"+
				"/** @implements {A} @implements {B} */\n"+
				"class C {};\n"+
				"MOD_C: 0;\n"))

	modCScope := c.GetLabeledStatement("MOD_C").enclosingScope
	aCtor := modCScope.GetVar("A").GetType().ToMaybeFunctionType()
	bCtor := modCScope.GetVar("B").GetType().ToMaybeFunctionType()
	cCtor := modCScope.GetVar("C").GetType().ToMaybeFunctionType()

	AssertThat(cCtor.GetAllImplementedInterfaces()).
		ContainsExactly(aCtor.GetInstanceType(), bCtor.GetInstanceType())
}

// Java: TypedScopeCreatorTest.testLegacyGoogModule_accessibleWithGoogRequire_exportingConstructor()
func (c *TypedScopeCreatorTest) TestLegacyGoogModule_accessibleWithGoogRequire_exportingConstructor() {
	c.TestSame(
		c.Srcs(
			"goog.module('mod.A');\n"+
				"goog.module.declareLegacyNamespace();\n"+
				"\n"+
				"exports = class A {};\n",
			"goog.require('mod.A');\n"+
				"const /** !mod.A */ a = new mod.A();\n"+
				"A: a;\n"))

	TypeSubject_AssertType(c.GetLabeledStatement("A").statementNode.GetOnlyChild().GetJSType()).
		ToStringIsEqualTo("mod.A")
}

// Java: TypedScopeCreatorTest.testLegacyGoogModule_accessibleWithGoogRequire_exportingLocalTypedef()
func (c *TypedScopeCreatorTest) TestLegacyGoogModule_accessibleWithGoogRequire_exportingLocalTypedef() {
	c.TestSame(
		c.Srcs(
			"goog.module('mod.A');\n"+
				"goog.module.declareLegacyNamespace();\n"+
				"\n"+
				"/** @typedef {number} */\n"+
				"let numType;\n"+
				"\n"+
				"exports = numType;\n",
			"goog.require('mod.A');\n"+
				"var /** !mod.A */ x;\n"+
				"X: x\n"))

	TypeSubject_AssertType(c.GetLabeledStatement("X").statementNode.GetOnlyChild().GetJSType()).IsNumber()
}

// Java: TypedScopeCreatorTest.testLegacyGoogModule_accessibleWithGoogRequire_exportingNamespace()
func (c *TypedScopeCreatorTest) TestLegacyGoogModule_accessibleWithGoogRequire_exportingNamespace() {
	c.TestSame(
		c.Srcs(
			"goog.module('mod');\n"+
				"goog.module.declareLegacyNamespace();\n"+
				"\n"+
				"exports.A = class A {};\n",
			"goog.require('mod');\n"+
				"const /** !mod.A */ a = new mod.A();\n"))
}

// Java: TypedScopeCreatorTest.testLegacyGoogModule_withNamedExport_extendedByGoogProvide()
func (c *TypedScopeCreatorTest) TestLegacyGoogModule_withNamedExport_extendedByGoogProvide() {
	c.TestSame(
		c.Srcs(
			"goog.module('mod');\n"+
				"goog.module.declareLegacyNamespace();\n"+
				"\n"+
				"exports.A = class A {};\n",
			"goog.provide('mod.B'); // This is bad style, but probably people do it.\n"+
				"mod.B = class B {};\n"))

	modType := c.globalScope.GetVar("mod").GetType()
	TypeSubject_AssertType(modType).WithTypeOfProp("A").ToStringIsEqualTo("(typeof mod.A)")
	TypeSubject_AssertType(modType).WithTypeOfProp("B").ToStringIsEqualTo("(typeof mod.B)")
}

// Java: TypedScopeCreatorTest.testLegacyGoogModule_canBeAliasedInGlobalScopeThenUsed()
func (c *TypedScopeCreatorTest) TestLegacyGoogModule_canBeAliasedInGlobalScopeThenUsed() {
	c.processClosurePrimitives = true
	c.TestSame(
		c.Srcs(
			"goog.module('mod.A');\n"+
				"goog.module.declareLegacyNamespace();\n"+
				"\n"+
				"class AInternal {}\n"+
				"exports = AInternal;\n",
			"goog.require('mod.A');\n"+
				"const A = mod.A;\n"+
				"const /** !A */ a = new A();\n"+
				"A: a;\n"))

	TypeSubject_AssertType(c.GetLabeledStatement("A").statementNode.GetOnlyChild().GetJSType()).
		ToStringIsEqualTo("AInternal")
}

// Java: TypedScopeCreatorTest.testLegacyGoogLoadModule_canBeAliasedInGlobalScopeThenUsed()
func (c *TypedScopeCreatorTest) TestLegacyGoogLoadModule_canBeAliasedInGlobalScopeThenUsed() {
	c.processClosurePrimitives = true
	c.TestSame(
		c.Srcs(
			"goog.loadModule(function(exports) {\n"+
				"  goog.module('mod.A');\n"+
				"  goog.module.declareLegacyNamespace();\n"+
				"\n"+
				"  exports = class A {};\n"+
				"  return exports;\n"+
				"});\n",
			"goog.require('mod.A');\n"+
				"const A = mod.A;\n"+
				"const /** !A */ a = new A();\n"+
				"A: a;\n"))

	TypeSubject_AssertType(c.GetLabeledStatement("A").statementNode.GetOnlyChild().GetJSType()).
		ToStringIsEqualTo("mod.A")
}

// Java: TypedScopeCreatorTest.testGoogModule_requiringGoogProvideShadowingGlobal()
func (c *TypedScopeCreatorTest) TestGoogModule_requiringGoogProvideShadowingGlobal() {
	c.processClosurePrimitives = true

	c.TestSame(
		c.Externs(
			NewTestExternsBuilder().AddClosureExterns().Build(),
			"/** @interface */ function EventTarget() {}"),
		c.Srcs(
			"goog.provide('goog.events.EventTarget');\n"+
				"/** @constructor */\n"+
				"goog.events.EventTarget = function() {};\n",
			"goog.module('mod.A');\n"+
				"const EventTarget = goog.require('goog.events.EventTarget');\n"+
				"LOCAL_EVENT_TARGET: EventTarget;\n"))

	TypeSubject_AssertType(c.GetLabeledStatement("LOCAL_EVENT_TARGET").statementNode.GetOnlyChild().GetJSType()).
		IsFunctionTypeThat().
		IsConstructorFor("goog.events.EventTarget")
}

// Java: TypedScopeCreatorTest.testGoogModule_requireShadowsGlobal_usedInExtendsClause()
func (c *TypedScopeCreatorTest) TestGoogModule_requireShadowsGlobal_usedInExtendsClause() {
	c.processClosurePrimitives = true

	c.TestSame(
		c.Externs(
			NewTestExternsBuilder().AddClosureExterns().Build(),
			"/** @interface */ function EventTarget() {}"),
		c.Srcs(
			"goog.provide('goog.events.EventTarget');\n"+
				"/** @constructor */\n"+
				"goog.events.EventTarget = function() {};\n",
			"goog.module('mod.A');\n"+
				"const EventTarget = goog.require('goog.events.EventTarget');\n"+
				"class Proxy extends EventTarget {}\n"+
				"PROXY: Proxy;\n"))

	proxyType :=
		c.GetLabeledStatement("PROXY").statementNode.GetOnlyChild().GetJSType().ToMaybeFunctionType()
	TypeSubject_AssertType(proxyType.GetSuperClassConstructor()).
		IsFunctionTypeThat().
		IsConstructorFor("goog.events.EventTarget")
}

// Java: TypedScopeCreatorTest.testLegacyGoogModule_withDefaultClassExport_createsGlobalType()
func (c *TypedScopeCreatorTest) TestLegacyGoogModule_withDefaultClassExport_createsGlobalType() {
	c.processClosurePrimitives = true
	c.TestSame(
		c.Srcs(
			"goog.module('mod.MyClass');\n" +
				"goog.module.declareLegacyNamespace();\n" +
				"class MyClass {}\n" +
				"exports = MyClass;\n"))

	TypeSubject_AssertType(c.registry.GetGlobalType("mod.MyClass")).ToStringIsEqualTo("MyClass")
}

// Java: TypedScopeCreatorTest.testLegacyGoogModule_withDefaultEnumExport_createsGlobalType()
func (c *TypedScopeCreatorTest) TestLegacyGoogModule_withDefaultEnumExport_createsGlobalType() {
	c.processClosurePrimitives = true
	c.TestSame(
		c.Srcs(
			"goog.module('mod.MyEnum');\n" +
				"goog.module.declareLegacyNamespace();\n" +
				"/** @enum {string} */\n" +
				"const MyEnum = {A: 'a', B: 'b'};\n" +
				"exports = MyEnum;\n"))

	TypeSubject_AssertType(c.registry.GetGlobalType("mod.MyEnum")).ToStringIsEqualTo("MyEnum<string>")
}

// Java: TypedScopeCreatorTest.testLegacyGoogModule_withDefaultTypedefExport_createsGlobalType()
func (c *TypedScopeCreatorTest) TestLegacyGoogModule_withDefaultTypedefExport_createsGlobalType() {
	c.processClosurePrimitives = true
	c.TestSame(
		c.Srcs(
			"goog.module('mod.MyTypedef');\n" +
				"goog.module.declareLegacyNamespace();\n" +
				"/** @typedef {number} */\n" +
				"let MyTypedef;\n" +
				"exports = MyTypedef;\n"))

	TypeSubject_AssertType(c.registry.GetGlobalType("mod.MyTypedef")).IsNumber()
}

// Java: TypedScopeCreatorTest.testLegacyGoogModule_withDefaultExport_extendedByGoogProvide()
func (c *TypedScopeCreatorTest) TestLegacyGoogModule_withDefaultExport_extendedByGoogProvide() {
	c.TestSame(
		c.Srcs(
			"goog.module('mod');\n"+
				"goog.module.declareLegacyNamespace();\n"+
				"\n"+
				"/** @return {number} */\n"+
				"exports = function() { return 0; };\n",
			"goog.provide('mod.B'); // This is bad style, but probably people do it.\n"+
				"mod.B = class B {};\n"))

	modType := c.globalScope.GetVar("mod").GetType()
	TypeSubject_AssertType(modType).IsFunctionTypeThat().HasReturnTypeThat().IsNumber()
	TypeSubject_AssertType(modType).WithTypeOfProp("B").ToStringIsEqualTo("(typeof mod.B)")
}

// Java: TypedScopeCreatorTest.testLegacyGoogModule_accessibleWithGoogRequire_exportingTypedef()
func (c *TypedScopeCreatorTest) TestLegacyGoogModule_accessibleWithGoogRequire_exportingTypedef() {
	c.TestSame(
		c.Srcs(
			"goog.module('mod');\n"+
				"goog.module.declareLegacyNamespace();\n"+
				"\n"+
				"/** @typedef {number} */\n"+
				"exports.numType;\n",
			"goog.require('mod');\n"+
				"var /** !mod.numType */ a;\n"))

	modType := c.globalScope.GetVar("a").GetType()
	TypeSubject_AssertType(modType).IsNumber()
}

// Java: TypedScopeCreatorTest.testGoogModuleRequiringGoogProvide_class()
func (c *TypedScopeCreatorTest) TestGoogModuleRequiringGoogProvide_class() {
	c.TestSame(
		c.Srcs(
			"goog.provide('a.b.Foo')\n"+
				"a.b.Foo = class {};\n",
			"goog.module('c');\n"+
				"const Foo = goog.require('a.b.Foo');\n"+
				"var /** !Foo */ x;\n"+
				"X: x;\n"))

	TypeSubject_AssertType(c.GetLabeledStatement("X").statementNode.GetOnlyChild().GetJSType()).
		ToStringIsEqualTo("a.b.Foo")
}

// Java: TypedScopeCreatorTest.testTypedef_namedExportInObjectLit()
func (c *TypedScopeCreatorTest) TestTypedef_namedExportInObjectLit() {
	c.TestSame(
		c.Srcs(
			"goog.module('a');\n"+
				"/** @typedef {number} */\n"+
				"let Foo;\n"+
				"exports = {Foo};\n",
			"goog.module('b');\n"+
				"const {Foo} = goog.require('a');\n"+
				"var /** !Foo */ x;\n"))
}

// Java: TypedScopeCreatorTest.testGoogModuleRequiringGoogProvide_classWithDestructuring()
func (c *TypedScopeCreatorTest) TestGoogModuleRequiringGoogProvide_classWithDestructuring() {
	c.TestSame(
		c.Srcs(
			"goog.provide('a.b')\n"+
				"a.b.Foo = class {};\n",
			"goog.module('c');\n"+
				"const {Foo} = goog.require('a.b');\n"+
				"var /** !Foo */ x;\n"+
				"X: x;\n"))

	TypeSubject_AssertType(c.GetLabeledStatement("X").statementNode.GetOnlyChild().GetJSType()).
		ToStringIsEqualTo("a.b.Foo")
}

// Java: TypedScopeCreatorTest.testEsModule_importStarMissingModule()
func (c *TypedScopeCreatorTest) TestEsModule_importStarMissingModule() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	// Make sure this does not crash the typechecker.
	c.TestError("import * as x from './invalid_path;'; X: x; export {x};", ModuleLoader_LOAD_WARNING)

	ScopeSubject_AssertScope(c.GetLabeledStatement("X").enclosingScope).Declares("x").WithTypeThat().IsUnknown()
	ScopeSubject_AssertScope(c.GetLabeledStatement("X").enclosingScope).
		Declares(Export_NAMESPACE).
		WithTypeThat().
		HasInferredProperty("x")
}

// Java: TypedScopeCreatorTest.testEsModule_importSpecsMissingModule()
func (c *TypedScopeCreatorTest) TestEsModule_importSpecsMissingModule() {
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
	// Make sure this does not crash the typechecker.
	c.TestError("import {x} from './invalid_path;'; X: x; export {x};", ModuleLoader_LOAD_WARNING)

	ScopeSubject_AssertScope(c.GetLabeledStatement("X").enclosingScope).Declares("x").WithTypeThat().IsUnknown()
	ScopeSubject_AssertScope(c.GetLabeledStatement("X").enclosingScope).
		Declares(Export_NAMESPACE).
		WithTypeThat().
		HasInferredProperty("x")
}

// Java: TypedScopeCreatorTest.testEsModule_importNameFromScript()
func (c *TypedScopeCreatorTest) TestEsModule_importNameFromScript() {
	// Make sure this does not crash the typechecker. Importing a script for the side effects is
	// fine, but you can't import a name.
	c.TestSame(
		c.Srcs(
			"const x = 'oops, you did not export me.';",
			"import {x} from './testcode0';\n"+
				"X: x;\n"))

	ScopeSubject_AssertScope(c.GetLabeledStatement("X").enclosingScope).Declares("x").WithTypeThat().IsUnknown()
}

// Java: TypedScopeCreatorTest.testEsModule_exportNameDeclaration()
func (c *TypedScopeCreatorTest) TestEsModule_exportNameDeclaration() {
	c.TestSame("/** @type {string|number} */ export let strOrNum = 0; MOD: 0;")

	strOrNum := c.GetLabeledStatement("MOD").enclosingScope.GetVar("strOrNum")
	TypedVarSubject_AssertThat(strOrNum).IsNotNull()
	TypedVarSubject_AssertThat(strOrNum).IsNotInferred()
	TypedVarSubject_AssertThat(strOrNum).HasJSTypeThat().ToStringIsEqualTo("(number|string)")
}

// Java: TypedScopeCreatorTest.testEsModule_importSpecs_declaredExport()
func (c *TypedScopeCreatorTest) TestEsModule_importSpecs_declaredExport() {
	c.TestSame(
		c.Srcs(
			"let /** number */ num;\n"+
				"export {num as y};\n",
			"import {y as x} from './testcode0';\n"+
				"X: x;\n"))

	AssertNode(c.GetLabeledStatement("X").statementNode.GetOnlyChild()).HasJSTypeThat().IsNumber()
	TypedVarSubject_AssertThat(c.GetLabeledStatement("X").enclosingScope.GetVar("x")).IsNotInferred()
}

// Java: TypedScopeCreatorTest.testEsModule_importSpecs_inferredExport_nameDeclaration()
func (c *TypedScopeCreatorTest) TestEsModule_importSpecs_inferredExport_nameDeclaration() {
	c.TestSame(
		c.Srcs(
			"/** @return {number} */ const f = () => 0;\n"+
				"export const y = f();\n",
			"import {y as x} from './testcode0';\n"+
				"X: x;\n"))

	AssertNode(c.GetLabeledStatement("X").statementNode.GetOnlyChild()).HasJSTypeThat().IsNumber()
	TypedVarSubject_AssertThat(c.GetLabeledStatement("X").enclosingScope.GetVar("x")).IsInferred()
}

// Java: TypedScopeCreatorTest.testEsModule_importSpecs_inferredExport_exportSpecs()
func (c *TypedScopeCreatorTest) TestEsModule_importSpecs_inferredExport_exportSpecs() {
	c.TestSame(
		c.Srcs(
			"/** @return {number} */ const f = () => 0;\n"+
				"const num = f();\n"+
				"export {num as y};\n",
			"import {y as x} from './testcode0';\n"+
				"X: x;\n"))

	AssertNode(c.GetLabeledStatement("X").statementNode.GetOnlyChild()).HasJSTypeThat().IsNumber()
	TypedVarSubject_AssertThat(c.GetLabeledStatement("X").enclosingScope.GetVar("x")).IsInferred()
}

// Java: TypedScopeCreatorTest.testEsModule_importSpecs_typedef()
func (c *TypedScopeCreatorTest) TestEsModule_importSpecs_typedef() {
	c.TestSame(
		c.Srcs(
			"/** @typedef {number} */ let numType;\n"+
				"export {numType};\n",
			"import {numType} from './testcode0';\n"+
				"var /** !numType */ x;\n"+
				"X: x;\n"))

	AssertNode(c.GetLabeledStatement("X").statementNode.GetOnlyChild()).HasJSTypeThat().IsNumber()
}

// Java: TypedScopeCreatorTest.testEsModule_importSpecs_fromGoogColonProvide()
func (c *TypedScopeCreatorTest) TestEsModule_importSpecs_fromGoogColonProvide() {
	c.TestSame(
		c.Srcs(
			"goog.provide('foo.bar');\n"+
				"/** @type {number} */\n"+
				"foo.bar.x = 0;\n",
			"import {x as xImported} from 'goog:foo.bar';\n"+
				"X: xImported;\n"))

	AssertNode(c.GetLabeledStatement("X").statementNode.GetOnlyChild()).HasJSTypeThat().IsNumber()
}

// Java: TypedScopeCreatorTest.testEsModule_importSpecs_fromGoogColonModule()
func (c *TypedScopeCreatorTest) TestEsModule_importSpecs_fromGoogColonModule() {
	c.TestSame(
		c.Srcs(
			"goog.module('foo.bar');\n"+
				"/** @type {number} */\n"+
				"exports.x = 0;\n",
			"import {x as xImported} from 'goog:foo.bar';\n"+
				"X: xImported;\n"))

	AssertNode(c.GetLabeledStatement("X").statementNode.GetOnlyChild()).HasJSTypeThat().IsNumber()
}

// Java: TypedScopeCreatorTest.testEsModule_moduleObjectTypeMatchesImportExportAndModuleBody()
func (c *TypedScopeCreatorTest) TestEsModule_moduleObjectTypeMatchesImportExportAndModuleBody() {
	c.TestSame(
		c.Srcs(
			SourceFile_FromCode(
				"a.js",
				"export class A {};\n"+
					"MOD_A: 0; // there's no way to reference the module object within the module\n"),
			SourceFile_FromCode(
				"b.js",
				"import * as a from './a.js';\n"+
					"IMPORTED_A: a;\n")))

	aModuleBody :=
		NodeUtil_GetEnclosingType(c.GetLabeledStatement("MOD_A").statementNode, rhino.Token_MODULE_BODY)
	moduleAObjectType := aModuleBody.GetJSType()
	AssertThat(moduleAObjectType).IsNotNull()

	importedANode := c.GetLabeledStatement("IMPORTED_A").statementNode.GetOnlyChild()
	AssertNode(importedANode).HasJSTypeThat().IsEqualTo(moduleAObjectType)
}

// Java: TypedScopeCreatorTest.testEsModule_exportDefault()
func (c *TypedScopeCreatorTest) TestEsModule_exportDefault() {
	c.TestSame(c.Srcs("export default 0;", "import x from './testcode0'; X: x;"))

	AssertNode(c.GetLabeledStatement("X").statementNode.GetOnlyChild()).HasJSTypeThat().IsNumber()
}

// Java: TypedScopeCreatorTest.testEsModule_exportDefaultImportedWithSpecs()
func (c *TypedScopeCreatorTest) TestEsModule_exportDefaultImportedWithSpecs() {
	c.TestSame(c.Srcs("export default 0;", "import {default as x} from './testcode0'; X: x;"))

	AssertNode(c.GetLabeledStatement("X").statementNode.GetOnlyChild()).HasJSTypeThat().IsNumber()
}

// Java: TypedScopeCreatorTest.testEsModule_importDefault_exportedUsingDefaultNamedKey()
func (c *TypedScopeCreatorTest) TestEsModule_importDefault_exportedUsingDefaultNamedKey() {
	c.TestSame(c.Srcs("const x = 0; export {x as default};", "import x from './testcode0'; X: x;"))

	AssertNode(c.GetLabeledStatement("X").statementNode.GetOnlyChild()).HasJSTypeThat().IsNumber()
}

// Java: TypedScopeCreatorTest.testEsModule_exportAndImportDefaultClassDeclaration()
func (c *TypedScopeCreatorTest) TestEsModule_exportAndImportDefaultClassDeclaration() {
	c.TestSame(
		c.Srcs(
			"export default class Button {}; const /** !Button */ b = new Button(); B1: b;",
			"import Button from './testcode0'; const /** !Button */ b = new Button(); B2: b;"))

	AssertNode(c.GetLabeledStatement("B1").statementNode.GetOnlyChild()).
		HasJSTypeThat().
		ToStringIsEqualTo("Button")
	AssertNode(c.GetLabeledStatement("B2").statementNode.GetOnlyChild()).
		HasJSTypeThat().
		ToStringIsEqualTo("Button")
}

// Java: TypedScopeCreatorTest.testEsModule_importAndImportDefault_nameLocalToModule()
func (c *TypedScopeCreatorTest) TestEsModule_importAndImportDefault_nameLocalToModule() {
	c.TestSame(
		c.Srcs(
			"class Button {} export default Button; const /** !Button */ b = new Button(); B1: b;",
			"import Button from './testcode0'; const /** !Button */ b = new Button(); B2: b;"))

	AssertNode(c.GetLabeledStatement("B1").statementNode.GetOnlyChild()).
		HasJSTypeThat().
		ToStringIsEqualTo("Button")
	AssertNode(c.GetLabeledStatement("B2").statementNode.GetOnlyChild()).
		HasJSTypeThat().
		ToStringIsEqualTo("Button")
}

// Java: TypedScopeCreatorTest.testEsModule_importDefault_fromGoogColonModule()
func (c *TypedScopeCreatorTest) TestEsModule_importDefault_fromGoogColonModule() {
	c.TestSame(
		c.Srcs(
			"goog.module('foo.bar');\n"+
				"/** @type {number} */\n"+
				"exports.x = 0;\n",
			"import bar from 'goog:foo.bar';\n"+
				"X: bar.x;\n"))

	AssertNode(c.GetLabeledStatement("X").statementNode.GetOnlyChild()).HasJSTypeThat().IsNumber()
}

// Java: TypedScopeCreatorTest.testEsModule_exportNameFrom()
func (c *TypedScopeCreatorTest) TestEsModule_exportNameFrom() {
	c.TestSame(
		c.Srcs(
			"export const x = 0; export class Button {}",
			"export {x, Button} from './testcode0';",
			"import {x as y, Button} from './testcode0';\n"+
				"Y: y;\n"+
				"const /** !Button */ b = new Button();\n"))

	AssertNode(c.GetLabeledStatement("Y").statementNode.GetOnlyChild()).HasJSTypeThat().IsNumber()
}

// Java: TypedScopeCreatorTest.testEsModule_exportStarFrom()
func (c *TypedScopeCreatorTest) TestEsModule_exportStarFrom() {
	c.TestSame(
		c.Srcs(
			"export const x = 0; export class Button {}",
			"export * from './testcode0';",
			"import {x as y, Button} from './testcode0';\n"+
				"Y: y;\n"+
				"const /** !Button */ b = new Button();\n"))

	AssertNode(c.GetLabeledStatement("Y").statementNode.GetOnlyChild()).HasJSTypeThat().IsNumber()
}

// Java: TypedScopeCreatorTest.testEsModule_importStar_declaredExport()
func (c *TypedScopeCreatorTest) TestEsModule_importStar_declaredExport() {
	c.TestSame(
		c.Srcs(
			"export let /** number */ x;",
			"import * as ns from './testcode0';\n"+
				"NS: ns;\n"+
				"X: ns.x;\n"))

	nsType := c.GetLabeledStatement("NS").statementNode.GetOnlyChild().GetJSType()
	TypeSubject_AssertType(nsType).HasDeclaredProperty("x")
	AssertNode(c.GetLabeledStatement("X").statementNode.GetOnlyChild()).HasJSTypeThat().IsNumber()
}

// Java: TypedScopeCreatorTest.testEsModule_importStar_inferredExport()
func (c *TypedScopeCreatorTest) TestEsModule_importStar_inferredExport() {
	c.TestSame(
		c.Srcs(
			"/** @return {number} */ const f = () => 0;\n"+
				"export const x = f();\n",
			"import * as ns from './testcode0';\n"+
				"NS: ns;\n"+
				"X: ns.x;\n"))

	nsType := c.GetLabeledStatement("NS").statementNode.GetOnlyChild().GetJSType()
	TypeSubject_AssertType(nsType).HasInferredProperty("x")
	AssertNode(c.GetLabeledStatement("X").statementNode.GetOnlyChild()).HasJSTypeThat().IsNumber()
}

// Java: TypedScopeCreatorTest.testEsModule_importStar_typedefInExportSpec()
func (c *TypedScopeCreatorTest) TestEsModule_importStar_typedefInExportSpec() {
	c.TestSame(
		c.Srcs(
			"/** @typedef {number} */ let numType; export {numType};",
			"import * as ns from './testcode0'; var /** !ns.numType */ x; X: x;"))

	AssertNode(c.GetLabeledStatement("X").statementNode.GetOnlyChild()).HasJSTypeThat().IsNumber()
}

// Java: TypedScopeCreatorTest.testEsModule_importStar_typedefExportDeclaration()
func (c *TypedScopeCreatorTest) TestEsModule_importStar_typedefExportDeclaration() {
	c.TestSame(
		c.Srcs(
			"/** @typedef {number} */ export let numType;",
			"import * as ns from './testcode0'; var /** !ns.numType */ x; X: x;"))

	AssertNode(c.GetLabeledStatement("X").statementNode.GetOnlyChild()).HasJSTypeThat().IsNumber()
}

// Java: TypedScopeCreatorTest.testEsModule_importStar_typedefReexported()
func (c *TypedScopeCreatorTest) TestEsModule_importStar_typedefReexported() {
	c.TestSame(
		c.Srcs(
			"/** @typedef {number} */ export let numType;",
			"import * as mod from './testcode0'; export {mod};",
			"import {mod} from './testcode1'; var /** !mod.numType */ x; X: x;"))

	AssertNode(c.GetLabeledStatement("X").statementNode.GetOnlyChild()).HasJSTypeThat().IsNumber()
}

// Java: TypedScopeCreatorTest.testEsModule_importStar_typedefReexportedThenImportStarred()
func (c *TypedScopeCreatorTest) TestEsModule_importStar_typedefReexportedThenImportStarred() {
	c.TestSame(
		c.Srcs(
			"/** @typedef {number} */ export let numType;",
			"import * as mod0 from './testcode0'; export {mod0};",
			"import * as mod1 from './testcode1'; var /** !mod1.mod0.numType */ x; X: x;"))

	AssertNode(c.GetLabeledStatement("X").statementNode.GetOnlyChild()).HasJSTypeThat().IsNumber()
}

// Java: TypedScopeCreatorTest.testGoogRequire_esModuleId_namespaceRequire()
func (c *TypedScopeCreatorTest) TestGoogRequire_esModuleId_namespaceRequire() {
	c.TestSame(
		c.Srcs(
			"goog.declareModuleId('a'); export const x = 0;",
			"goog.module('b'); const a = goog.require('a'); X: a.x;"))

	AssertNode(c.GetLabeledStatement("X").statementNode.GetOnlyChild()).HasJSTypeThat().IsNumber()
}

// Java: TypedScopeCreatorTest.testGoogRequire_esModuleId_destructuringRequire()
func (c *TypedScopeCreatorTest) TestGoogRequire_esModuleId_destructuringRequire() {
	c.TestSame(
		c.Srcs(
			"goog.declareModuleId('a'); export const x = 0;",
			"goog.module('b'); const {x} = goog.require('a'); X: x;"))

	AssertNode(c.GetLabeledStatement("X").statementNode.GetOnlyChild()).HasJSTypeThat().IsNumber()
}

// Java: TypedScopeCreatorTest.testGoogRequire_insideEsModule_namedExport()
func (c *TypedScopeCreatorTest) TestGoogRequire_insideEsModule_namedExport() {
	c.TestSame(
		c.Srcs(
			"goog.module('b'); exports.x = 0;",
			"const {x} = goog.require('b'); X: x; export {x};"))

	AssertNode(c.GetLabeledStatement("X").statementNode.GetOnlyChild()).HasJSTypeThat().IsNumber()
}

// Java: TypedScopeCreatorTest.testGoogRequire_insideEsModule_class()
func (c *TypedScopeCreatorTest) TestGoogRequire_insideEsModule_class() {
	c.TestSame(
		c.Srcs(
			"goog.module('b');\n"+
				"exports.X = class {};\n",
			"const {X} = goog.require('b');\n"+
				"var /** !X */ x;\n"+
				"X: x;\n"+
				"export {x};\n"))

	AssertNode(c.GetLabeledStatement("X").statementNode.GetOnlyChild()).
		HasJSTypeThat().
		ToStringIsEqualTo("b.X")
}

// Java: TypedScopeCreatorTest.testEsModuleImportCycle_namedExports()
func (c *TypedScopeCreatorTest) TestEsModuleImportCycle_namedExports() {
	// Note: we cannot reference circular imports directly in the module body but we can reference
	// them by type and inside functions.
	c.TestSame(
		c.Srcs(
			"import {Bar} from './testcode1';\n"+
				"var /** !Bar */ b;\n"+
				"BAR: b;\n"+
				"export class Foo {}\n",
			"import {Foo} from './testcode0';\n"+
				"var /** !Foo */ f;\n"+
				"FOO: f;\n"+
				"export class Bar {}\n"))

	AssertNode(c.GetLabeledStatement("FOO").statementNode.GetOnlyChild()).
		HasJSTypeThat().
		ToStringIsEqualTo("Foo")
	AssertNode(c.GetLabeledStatement("BAR").statementNode.GetOnlyChild()).
		HasJSTypeThat().
		ToStringIsEqualTo("Bar")
}

// Java: TypedScopeCreatorTest.testEsModuleImportCycle_aliasingNamedExports()
func (c *TypedScopeCreatorTest) TestEsModuleImportCycle_aliasingNamedExports() {
	c.TestSame(
		c.Srcs(
			"import {Bar as BarLocal} from './testcode1';\n"+
				"var /** !BarLocal */ b;\n"+
				"BAR: b;\n"+
				"export class Foo {}\n",
			"import {Foo as FooLocal} from './testcode0';\n"+
				"var /** !FooLocal */ f;\n"+
				"FOO: f;\n"+
				"export class Bar {}\n"))

	AssertNode(c.GetLabeledStatement("FOO").statementNode.GetOnlyChild()).
		HasJSTypeThat().
		ToStringIsEqualTo("Foo")
	AssertNode(c.GetLabeledStatement("BAR").statementNode.GetOnlyChild()).
		HasJSTypeThat().
		ToStringIsEqualTo("Bar")
}

// Java: TypedScopeCreatorTest.testEsModuleImportCycle_importStarPlusNamedExport()
func (c *TypedScopeCreatorTest) TestEsModuleImportCycle_importStarPlusNamedExport() {
	c.TestSame(
		c.Srcs(
			"import * as mod from './testcode1';\n"+
				"var /** !mod.Bar */ b;\n"+
				"BAR: b;\n"+
				"export class Foo {}\n",
			"import {Foo} from './testcode0';\n"+
				"var /** !Foo */ f;\n"+
				"FOO: f;\n"+
				"export class Bar {}\n"))

	AssertNode(c.GetLabeledStatement("FOO").statementNode.GetOnlyChild()).
		HasJSTypeThat().
		ToStringIsEqualTo("Foo")
	AssertNode(c.GetLabeledStatement("BAR").statementNode.GetOnlyChild()).
		HasJSTypeThat().
		ToStringIsEqualTo("Bar")
}

// Java: TypedScopeCreatorTest.testEsModuleImportStar_reexportLateReference()
func (c *TypedScopeCreatorTest) TestEsModuleImportStar_reexportLateReference() {
	c.TestSame(
		c.Srcs(
			"import * as mod from './testcode1';\n"+
				"function f() {\n"+
				"  BAR1: new mod.Bar();\n"+
				"}\n"+
				"export {mod};\n",
			"import {mod} from './testcode0';\n"+
				"var /** !mod.Bar */ b;\n"+
				"BAR: b;\n"+
				"function f() {\n"+
				"  BAR2: new mod.Bar();\n"+
				"}\n"+
				"export class Bar {}\n"))

	AssertNode(c.GetLabeledStatement("BAR").statementNode.GetOnlyChild()).
		HasJSTypeThat().
		ToStringIsEqualTo("Bar")
	AssertNode(c.GetLabeledStatement("BAR1").statementNode.GetOnlyChild()).
		HasJSTypeThat().
		ToStringIsEqualTo("Bar")
	AssertNode(c.GetLabeledStatement("BAR2").statementNode.GetOnlyChild()).
		HasJSTypeThat().
		ToStringIsEqualTo("Bar")
}

// Java: TypedScopeCreatorTest.testEsModuleImportCycle_importStar()
func (c *TypedScopeCreatorTest) TestEsModuleImportCycle_importStar() {
	c.TestSame(
		c.Srcs(
			"import * as mod from './testcode1';\n"+
				"var /** !mod.Bar */ b;\n"+
				"BAR: b;\n"+
				"export class Foo {}\n",
			"import * as mod from './testcode0';\n"+
				"var /** !mod.Foo */ f;\n"+
				"FOO: f;\n"+
				"export class Bar {}\n"))

	AssertNode(c.GetLabeledStatement("FOO").statementNode.GetOnlyChild()).
		HasJSTypeThat().
		ToStringIsEqualTo("Foo")
	AssertNode(c.GetLabeledStatement("BAR").statementNode.GetOnlyChild()).
		HasJSTypeThat().
		ToStringIsEqualTo("Bar")
}

// Java: TypedScopeCreatorTest.testInferredLegacyNsExport()
func (c *TypedScopeCreatorTest) TestInferredLegacyNsExport() {
	c.processClosurePrimitives = true
	c.TestSame(
		c.Srcs(
			"goog.module('a.b.c')\n" +
				"goog.module.declareLegacyNamespace();\n" +
				"MOD: 0;\n" +
				"\n" +
				"/** @return {number} */\n" +
				"function getNumber() { return 0; }\n" +
				"exports = getNumber();\n"))

	ab := c.globalScope.GetVar("a.b")
	TypedVarSubject_AssertThat(ab).IsNotInferred()
	TypedVarSubject_AssertThat(ab).HasJSTypeThat().WithTypeOfProp("c").IsNumber()

	abc := c.globalScope.GetVar("a.b.c")
	TypedVarSubject_AssertThat(abc).HasJSTypeThat().IsNumber()

	moduleBody := c.GetLabeledStatement("MOD").enclosingScope.GetRootNode()
	jutil.CheckStateMsg1(moduleBody.IsModuleBody(), "%s", moduleBody)
	AssertNode(moduleBody).HasJSTypeThat().IsNumber()
}

// Java: TypedScopeCreatorTest.testInferredLegacyNsExport_bundled()
func (c *TypedScopeCreatorTest) TestInferredLegacyNsExport_bundled() {
	c.processClosurePrimitives = true
	c.TestSame(
		c.Srcs(
			"goog.loadModule(function(exports) {\n" +
				"goog.module('a.b.c')\n" +
				"goog.module.declareLegacyNamespace();\n" +
				"MOD: 0;\n" +
				"\n" +
				"/** @return {number} */\n" +
				"function getNumber() { return 0; }\n" +
				"exports = getNumber();\n" +
				"return exports;\n" +
				"});\n"))

	ab := c.globalScope.GetVar("a.b")
	TypedVarSubject_AssertThat(ab).IsNotInferred()
	TypedVarSubject_AssertThat(ab).HasJSTypeThat().WithTypeOfProp("c").IsNumber()

	abc := c.globalScope.GetVar("a.b.c")
	TypedVarSubject_AssertThat(abc).HasJSTypeThat().IsNumber()

	fnBlock := c.GetLabeledStatement("MOD").enclosingScope.GetRootNode()
	exportsParam := fnBlock.GetPrevious().GetFirstChild()
	jutil.CheckState(exportsParam.GetString() == "exports")
	AssertNode(exportsParam).HasJSTypeThat().IsNumber()
}

// Java: TypedScopeCreatorTest.testGoogProvide_typedef_requiredInModule()
func (c *TypedScopeCreatorTest) TestGoogProvide_typedef_requiredInModule() {
	c.processClosurePrimitives = true
	c.TestSame(
		c.Srcs(
			"goog.provide('a.b.c');\n"+
				"/** @typedef {string} */\n"+
				"a.b.c;\n",
			"goog.module('mod');\n"+
				"goog.module.declareLegacyNamespace();\n"+
				"const c = goog.require('a.b.c');\n"+
				"/** @param {!c} d */\n"+
				"function f(d) {\n"+
				"  D: d;\n"+
				"}\n"))

	TypeSubject_AssertType(c.registry.GetGlobalType("a.b.c")).IsString()
	TypeSubject_AssertType(c.GetLabeledStatement("D").statementNode.GetOnlyChild().GetJSType()).IsString()
}

// Java: TypedScopeCreatorTest.testGoogModuleGetInBlockScopeDoesntCrash()
func (c *TypedScopeCreatorTest) TestGoogModuleGetInBlockScopeDoesntCrash() {
	c.TestSame(
		c.Srcs(
			"goog.module('a');\n"+
				"exports = 5;\n",
			"goog.require('a');\n"+
				"function f() {\n"+
				"  if (true) {\n"+
				"    D: var a = goog.module.get('a');\n"+
				"  }\n"+
				"}\n"))

	TypeSubject_AssertType(c.GetLabeledStatement("D").statementNode.GetOnlyChild().GetJSType()).IsNumber()
}

// Java: TypedScopeCreatorTest.testLegacyGoogModule_withExportedTypedefOnProperty()
func (c *TypedScopeCreatorTest) TestLegacyGoogModule_withExportedTypedefOnProperty() {
	c.processClosurePrimitives = true
	c.TestSame(
		c.Srcs(
			"goog.module('mod.a');\n"+
				"goog.module.declareLegacyNamespace();\n"+
				"\n"+
				"/** @typedef {number} */\n"+
				"exports.numType;\n",
			"goog.require('mod.a');\n"+
				"var /** !mod.a.numType */ a;\n"))

	modType := c.globalScope.GetVar("a").GetType()
	TypeSubject_AssertType(modType).IsNumber()
}

// Java: TypedScopeCreatorTest.testEsModule_exportFunctionDecl_createsLocalDeclaration()
func (c *TypedScopeCreatorTest) TestEsModule_exportFunctionDecl_createsLocalDeclaration() {
	c.TestSame("export function foo() {}; MOD: 0;")

	ScopeSubject_AssertScope(c.GetLabeledStatement("MOD").enclosingScope).
		Declares("foo").
		WithTypeThat().
		IsFunctionTypeThat().
		HasReturnTypeThat().
		IsVoid()
}

// Java: TypedScopeCreatorTest.testEsModule_exportDefaultFunctionDecl_createsLocalDeclaration()
func (c *TypedScopeCreatorTest) TestEsModule_exportDefaultFunctionDecl_createsLocalDeclaration() {
	c.TestSame("export default function foo() {}; MOD: 0;")

	ScopeSubject_AssertScope(c.GetLabeledStatement("MOD").enclosingScope).
		Declares("foo").
		WithTypeThat().
		IsFunctionTypeThat().
		HasReturnTypeThat().
		IsVoid()
}

// Java: TypedScopeCreatorTest.testConstantAliasOfForwardReferencedType_name()
func (c *TypedScopeCreatorTest) TestConstantAliasOfForwardReferencedType_name() {
	c.TestSame(
		"/** @const {!Type} */ let declaredVar;\n" +
			"const constantAlias = declaredVar;\n" +
			"class Type {}\n")

	declaredVar := c.globalScope.GetVar("declaredVar")
	constantAlias := c.globalScope.GetVar("constantAlias")

	TypedVarSubject_AssertThat(declaredVar).IsNotInferred()
	TypedVarSubject_AssertThat(constantAlias).IsNotInferred()
}

// Java: TypedScopeCreatorTest.testConstantAliasOfForwardReferencedType_qname()
func (c *TypedScopeCreatorTest) TestConstantAliasOfForwardReferencedType_qname() {
	c.TestSame(
		"/** @const {!Type} */ let declaredVar;\n" +
			"const a = {};\n" +
			"/** @const */\n" +
			"a.constantAlias = declaredVar;\n" +
			"class Type {}\n")

	declaredVar := c.globalScope.GetVar("declaredVar")
	aObject := c.globalScope.GetVar("a")
	constantAlias := c.globalScope.GetVar("a.constantAlias")

	TypedVarSubject_AssertThat(declaredVar).IsNotInferred()
	TypedVarSubject_AssertThat(aObject).HasJSTypeThat().HasDeclaredProperty("constantAlias")
	TypedVarSubject_AssertThat(constantAlias).IsNotInferred()
}

// Java: TypedScopeCreatorTest.testConstantAliasOfProperty()
func (c *TypedScopeCreatorTest) TestConstantAliasOfProperty() {
	c.TestSame(
		"/** @typedef {{name: string}} */\n" +
			"let Type;\n" +
			"class C {\n" +
			"  /** @param {!Type} obj */\n" +
			"  constructor(obj) {\n" +
			"    /** @private @const */\n" +
			"    this.name_ = obj.name;\n" +
			"  }\n" +
			"  method() {\n" +
			"    NAME: this.name_;\n" +
			"  }\n" +
			"}\n")

	thisDotName := c.GetLabeledStatement("NAME").statementNode.GetOnlyChild()
	AssertNode(thisDotName).HasJSTypeThat().IsString()

	cType := thisDotName.GetFirstChild().GetJSType().ToObjectType()
	TypeSubject_AssertType(cType).HasDeclaredProperty("name_")
	TypeSubject_AssertType(cType).WithTypeOfProp("name_").IsString()
}

// Java: TypedScopeCreatorTest.testConstantAliasOfForwardReferencedTypeProperty()
func (c *TypedScopeCreatorTest) TestConstantAliasOfForwardReferencedTypeProperty() {
	c.TestSame(
		"class C {\n" +
			"  /** @param {!Type} obj */\n" +
			"  constructor(obj) {\n" +
			"    /** @private @const */\n" +
			"    this.name_ = obj.name;\n" +
			"  }\n" +
			"  method() {\n" +
			"    NAME: this.name_;\n" +
			"  }\n" +
			"}\n" +
			"/** @typedef {{name: string}} */\n" +
			"let Type;\n")

	thisDotName := c.GetLabeledStatement("NAME").statementNode.GetOnlyChild()
	AssertNode(thisDotName).HasJSTypeThat().IsString()

	cType := thisDotName.GetFirstChild().GetJSType().ToObjectType()
	TypeSubject_AssertType(cType).HasInferredProperty("name_")
	TypeSubject_AssertType(cType).WithTypeOfProp("name_").IsString()
}

// Java: TypedScopeCreatorTest.testReportBadTypeAnnotationInTemplateParameter()
func (c *TypedScopeCreatorTest) TestReportBadTypeAnnotationInTemplateParameter() {
	c.TestWarning("var /** !Array<!MissingType> */ x;", RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR)
}

// Java: TypedScopeCreatorTest.testReportBadTypeAnnotation_invalidEnumTemplateParametersAreParsed()
func (c *TypedScopeCreatorTest) TestReportBadTypeAnnotation_invalidEnumTemplateParametersAreParsed() {
	// Enums cannot actually be templatized, but we still parse any types a user attempts to write.
	c.TestWarning(
		"/** @enum */ const Colors = {RED: 0}; var /** !Colors<!MissingType> */ x;",
		RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR)
}

// Java: TypedScopeCreatorTest.testReportBadTypeAnnotationInExtraTemplateParameter()
func (c *TypedScopeCreatorTest) TestReportBadTypeAnnotationInExtraTemplateParameter() {
	c.Test(
		c.Srcs("class C {} var /** !C<!MissingType> */ x;"),
		c.Warning(RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR),
		c.Warning(RhinoErrorReporter_TOO_MANY_TEMPLATE_PARAMS))
	c.Test(
		c.Srcs("var /** !Array<string, !MissingType> */ x;"),
		c.Warning(RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR),
		c.Warning(RhinoErrorReporter_TOO_MANY_TEMPLATE_PARAMS))
}

// Java: TypedScopeCreatorTest.testReportBadTypeAnnotationInForwardReferenceTemplateParameter()
func (c *TypedScopeCreatorTest) TestReportBadTypeAnnotationInForwardReferenceTemplateParameter() {
	c.TestWarning(
		"var /** !C<!MissingType> */ x; /** @template T */ class C {}",
		RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR)
}

// Java: TypedScopeCreatorTest.testReportBadTypeAnnotationInExtraForwardReferenceTemplateParameter()
func (c *TypedScopeCreatorTest) TestReportBadTypeAnnotationInExtraForwardReferenceTemplateParameter() {
	c.TestWarning(
		"var /** !C<!MissingType> */ x; class C {}", RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR)
	c.TestWarning(
		"var /** !C<string, !MissingType> */ x; /** @template T */ class C {}",
		RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR)
}

// Java: TypedScopeCreatorTest.testClassExtendsForwardReference_markedForLooserTypechecking()
func (c *TypedScopeCreatorTest) TestClassExtendsForwardReference_markedForLooserTypechecking() {
	c.TestSame(
		"/** @return {function(new: Object): ?} */\n" +
			"function mixin() {}\n" +
			"/** @extends {Parent} */\n" +
			"class Middle extends mixin() {}\n" +
			"class Child extends Middle {}\n" +
			"class Parent {}\n")

	parentCtor := typedScopeCreatorTest_cast[rhino.FunctionType](c.globalScope.GetVar("Parent").GetType())
	middleCtor := typedScopeCreatorTest_cast[rhino.FunctionType](c.globalScope.GetVar("Middle").GetType())
	childCtor := typedScopeCreatorTest_cast[rhino.FunctionType](c.globalScope.GetVar("Child").GetType())

	AssertThat(parentCtor.LoosenTypecheckingDueToForwardReferencedSupertype()).IsFalse()
	AssertThat(middleCtor.LoosenTypecheckingDueToForwardReferencedSupertype()).IsTrue()
	AssertThat(childCtor.LoosenTypecheckingDueToForwardReferencedSupertype()).IsTrue()

	AssertThat(parentCtor.GetPrototype().LoosenTypecheckingDueToForwardReferencedSupertype()).
		IsFalse()
	AssertThat(middleCtor.GetPrototype().LoosenTypecheckingDueToForwardReferencedSupertype()).
		IsTrue()
	AssertThat(childCtor.GetPrototype().LoosenTypecheckingDueToForwardReferencedSupertype()).
		IsTrue()

	AssertThat(parentCtor.GetInstanceType().LoosenTypecheckingDueToForwardReferencedSupertype()).
		IsFalse()
	AssertThat(middleCtor.GetInstanceType().LoosenTypecheckingDueToForwardReferencedSupertype()).
		IsTrue()
	AssertThat(childCtor.GetInstanceType().LoosenTypecheckingDueToForwardReferencedSupertype()).
		IsTrue()
}

// Java: TypedScopeCreatorTest.testClassExtendsForwardReference_forwardReferenceMarkedForLooserTypechecking()
func (c *TypedScopeCreatorTest) TestClassExtendsForwardReference_forwardReferenceMarkedForLooserTypechecking() {
	c.TestSame(
		"/** @return {function(new: Object): ?} */\n" +
			"function mixin() {}\n" +
			"/** @extends {Parent} */\n" +
			"class Middle extends mixin() {}\n" +
			"class Parent {}\n")
	var namedMiddleType *rhino.NamedType
	func() {
		closer := c.registry.GetResolver().OpenForDefinition()
		defer closer.Close()
		namedMiddleType = c.registry.CreateNamedType(c.globalScope, "Middle", nil, -1, -1)

		AssertThat(namedMiddleType.IsResolved()).IsFalse()
		AssertThat(namedMiddleType.LoosenTypecheckingDueToForwardReferencedSupertype()).IsFalse()
	}() // types are resolved here

	AssertThat(namedMiddleType.LoosenTypecheckingDueToForwardReferencedSupertype()).IsTrue()
}

// Java: TypedScopeCreatorTest.testClassImplementsForwardReference_markedForLooserTypechecking()
func (c *TypedScopeCreatorTest) TestClassImplementsForwardReference_markedForLooserTypechecking() {
	c.TestSame(
		"/** @implements {Parent} */\n" +
			"class Middle {}\n" +
			"\n" +
			"/** @interface @extends {Middle} */\n" +
			"class ExtendingChild {}\n" +
			"/** @implements {Middle} */\n" +
			"class ImplementingChild {}\n" +
			"\n" +
			"/** @interface */\n" +
			"class Parent {}\n")

	parentCtor := typedScopeCreatorTest_cast[rhino.FunctionType](c.globalScope.GetVar("Parent").GetType())
	middleCtor := typedScopeCreatorTest_cast[rhino.FunctionType](c.globalScope.GetVar("Middle").GetType())
	implementingChildCtor :=
		typedScopeCreatorTest_cast[rhino.FunctionType](c.globalScope.GetVar("ImplementingChild").GetType())
	extendingChildCtor := typedScopeCreatorTest_cast[rhino.FunctionType](c.globalScope.GetVar("ExtendingChild").GetType())

	AssertThat(parentCtor.LoosenTypecheckingDueToForwardReferencedSupertype()).IsFalse()
	// Static inheritance is not modeled for @implements so no need to loosen typechecking
	AssertThat(middleCtor.LoosenTypecheckingDueToForwardReferencedSupertype()).IsFalse()
	AssertThat(implementingChildCtor.LoosenTypecheckingDueToForwardReferencedSupertype()).IsFalse()
	AssertThat(extendingChildCtor.LoosenTypecheckingDueToForwardReferencedSupertype()).IsFalse()

	AssertThat(parentCtor.GetPrototype().LoosenTypecheckingDueToForwardReferencedSupertype()).
		IsFalse()
	AssertThat(middleCtor.GetPrototype().LoosenTypecheckingDueToForwardReferencedSupertype()).
		IsFalse()
	AssertThat(
		implementingChildCtor.
			GetPrototype().
			LoosenTypecheckingDueToForwardReferencedSupertype()).
		IsFalse()
	AssertThat(
		extendingChildCtor.GetPrototype().LoosenTypecheckingDueToForwardReferencedSupertype()).
		IsFalse()

	AssertThat(parentCtor.GetInstanceType().LoosenTypecheckingDueToForwardReferencedSupertype()).
		IsFalse()
	AssertThat(middleCtor.GetInstanceType().LoosenTypecheckingDueToForwardReferencedSupertype()).
		IsTrue()
	AssertThat(
		implementingChildCtor.
			GetInstanceType().
			LoosenTypecheckingDueToForwardReferencedSupertype()).
		IsTrue()
	AssertThat(
		extendingChildCtor.
			GetInstanceType().
			LoosenTypecheckingDueToForwardReferencedSupertype()).
		IsTrue()
}

// Java: TypedScopeCreatorTest.testInterfaceExtendsForwardReference_markedForLooserTypechecking()
func (c *TypedScopeCreatorTest) TestInterfaceExtendsForwardReference_markedForLooserTypechecking() {
	c.TestSame(
		"/** @interface @extends {Parent} */\n" +
			"class Middle {}\n" +
			"\n" +
			"/** @interface @extends {Middle} */\n" +
			"class ExtendingChild {}\n" +
			"/** @implements {Middle} */\n" +
			"class ImplementingChild {}\n" +
			"\n" +
			"/** @interface */\n" +
			"class Parent {}\n")

	parentCtor := typedScopeCreatorTest_cast[rhino.FunctionType](c.globalScope.GetVar("Parent").GetType())
	middleCtor := typedScopeCreatorTest_cast[rhino.FunctionType](c.globalScope.GetVar("Middle").GetType())
	implementingChildCtor :=
		typedScopeCreatorTest_cast[rhino.FunctionType](c.globalScope.GetVar("ImplementingChild").GetType())
	extendingChildCtor := typedScopeCreatorTest_cast[rhino.FunctionType](c.globalScope.GetVar("ExtendingChild").GetType())

	AssertThat(parentCtor.LoosenTypecheckingDueToForwardReferencedSupertype()).IsFalse()
	AssertThat(middleCtor.LoosenTypecheckingDueToForwardReferencedSupertype()).IsFalse()
	AssertThat(implementingChildCtor.LoosenTypecheckingDueToForwardReferencedSupertype()).IsFalse()
	AssertThat(extendingChildCtor.LoosenTypecheckingDueToForwardReferencedSupertype()).IsFalse()

	AssertThat(parentCtor.GetPrototype().LoosenTypecheckingDueToForwardReferencedSupertype()).
		IsFalse()
	AssertThat(middleCtor.GetPrototype().LoosenTypecheckingDueToForwardReferencedSupertype()).
		IsFalse()
	AssertThat(
		implementingChildCtor.
			GetPrototype().
			LoosenTypecheckingDueToForwardReferencedSupertype()).
		IsFalse()
	AssertThat(
		extendingChildCtor.GetPrototype().LoosenTypecheckingDueToForwardReferencedSupertype()).
		IsFalse()

	AssertThat(parentCtor.GetInstanceType().LoosenTypecheckingDueToForwardReferencedSupertype()).
		IsFalse()
	AssertThat(middleCtor.GetInstanceType().LoosenTypecheckingDueToForwardReferencedSupertype()).
		IsTrue()
	AssertThat(
		implementingChildCtor.
			GetInstanceType().
			LoosenTypecheckingDueToForwardReferencedSupertype()).
		IsTrue()
	AssertThat(
		extendingChildCtor.
			GetInstanceType().
			LoosenTypecheckingDueToForwardReferencedSupertype()).
		IsTrue()
}

// Java: TypedScopeCreatorTest.testGoogModule_exportsPropertyIsObjectLit_whosePropertiesNotInferredConst()
func (c *TypedScopeCreatorTest) TestGoogModule_exportsPropertyIsObjectLit_whosePropertiesNotInferredConst() {
	c.TestSame(
		c.Srcs(
			"goog.module('a.b.c');\n" +
				"function modifyThings() {\n" +
				"  exports.things.thing1 = 'something';\n" +
				"}\n" +
				"MOD_SCOPE: null;\n" +
				"exports.things = {\n" +
				"  thing1: null,\n" +
				"};\n"))

	moduleScope := c.GetLabeledStatement("MOD_SCOPE").enclosingScope

	exportsType := c.FindNameType("exports", moduleScope).ToObjectType()
	TypeSubject_AssertType(exportsType).HasDeclaredProperty("things")

	// Regression test for a bug where 'thing1' used to be a declared property, and assigning
	//   exports.things.thing1 = 'something'
	// caused a type mismatch error.
	exportsThingsType := exportsType.GetPropertyType("things").ToObjectType()
	TypeSubject_AssertType(exportsThingsType).HasInferredProperty("thing1")
	TypeSubject_AssertType(exportsThingsType).
		WithTypeOfProp("thing1").
		IsEqualTo(
			c.registry.CreateUnionType(
				c.registry.GetNativeType(rhino.JSTypeNative_STRING_TYPE), c.registry.GetNativeType(rhino.JSTypeNative_NULL_TYPE)))
}

// Java: TypedScopeCreatorTest.testLegacyModuleExternsConflict_defaultExport()
func (c *TypedScopeCreatorTest) TestLegacyModuleExternsConflict_defaultExport() {
	c.TestSame(
		c.Externs(
			"/** @const */\n"+
				"var a = {};\n"+
				"a.Foo = class {};\n"),
		c.Srcs(
			NewTestExternsBuilder().AddClosureExterns().Build(),
			"goog.module('a.Foo');\n"+
				"goog.module.declareLegacyNamespace();\n"+
				"class Foo {}\n"+
				"exports = Foo;\n",
			"goog.module('a.Bar');\n"+
				"const Foo = goog.require('a.Foo');\n"+
				"REQUIRED_FOO: Foo;\n"),
		// There are two warnings, one for assigning to 'a' and one for 'a.Foo'.
		c.Warning(TypeValidator_DUP_VAR_DECLARATION_TYPE_MISMATCH),
		c.Warning(TypeValidator_DUP_VAR_DECLARATION_TYPE_MISMATCH))

	aVar := c.globalScope.GetVar("a")
	AssertNode(aVar.GetNameNode()).IsFromExterns()

	aDotFoo := c.globalScope.GetVar("a.Foo")
	AssertNode(aDotFoo.GetNameNode()).IsFromExterns()

	TypedVarSubject_AssertThat(aVar).HasJSTypeThat().WithTypeOfProp("Foo").IsEqualTo(aDotFoo.GetType())

	TypeSubject_AssertType(c.registry.GetGlobalType("a.Foo")).
		IsEqualTo(aDotFoo.GetType().ToMaybeFunctionType().GetInstanceType())
}

// Java: TypedScopeCreatorTest.testLegacyModuleExternsConflict_namedExport()
func (c *TypedScopeCreatorTest) TestLegacyModuleExternsConflict_namedExport() {
	c.TestWarning(
		c.Externs(
			CompilerTestCase_DEFAULT_EXTERNS,
			"/** @const */\n"+
				"var a = {};\n"+
				"a.Foo = class {};\n"),
		c.Srcs(
			"goog.module('a');\n"+
				"goog.module.declareLegacyNamespace();\n"+
				"class Foo {}\n"+
				"exports.Foo = Foo;\n"),
		c.Warning(TypeValidator_DUP_VAR_DECLARATION_TYPE_MISMATCH))

	aVar := c.globalScope.GetVar("a")
	AssertNode(aVar.GetNameNode()).IsFromExterns()
	aDotFoo := c.globalScope.GetVar("a.Foo")
	AssertNode(aDotFoo.GetNameNode()).IsFromExterns()

	TypedVarSubject_AssertThat(aVar).HasJSTypeThat().WithTypeOfProp("Foo").IsEqualTo(aDotFoo.GetType())
}

// Java: TypedScopeCreatorTest.testRequireType_inheritanceChainWithIdenticalClassAndInterfaceName()
func (c *TypedScopeCreatorTest) TestRequireType_inheritanceChainWithIdenticalClassAndInterfaceName() {
	c.TestSame(
		c.Srcs(
			"goog.module('a.Foo');\n"+
				"const BFoo = goog.requireType('b.Foo');\n"+
				"var /** !BFoo */ b;\n"+
				"\n"+
				"/** @interface */\n"+
				"class Foo {}\n"+
				"exports = Foo;\n",
			"goog.module('b.Foo');\n"+
				"const AFoo = goog.require('a.Foo');\n"+
				"/** @implements {AFoo} */\n"+
				"class Foo {}\n"+
				"exports = Foo;\n"))
}

// Java: TypedScopeCreatorTest.testUnorderedGoogProvideAndModule_doesntCrash()
func (c *TypedScopeCreatorTest) TestUnorderedGoogProvideAndModule_doesntCrash() {
	c.TestSame(
		c.Srcs(
			"goog.provide('ns.boo');",
			"goog.module('test');\n"+
				"const boo = goog.require('ns.boo');\n",
			"ns = ns || {};\n"+
				"ns.boo = ns.boo || {};\n"))
}

// Java: TypedScopeCreatorTest.testLegacyGoogModule_redeclaredAsVar_reportsError()
func (c *TypedScopeCreatorTest) TestLegacyGoogModule_redeclaredAsVar_reportsError() {
	// This is a bad pattern but should not crash the compiler.
	c.Test(
		c.Srcs(
			"goog.module('GlobalName'); goog.module.declareLegacyNamespace();",
			"var GlobalName = class {};"),
		c.Warning(TypeValidator_DUP_VAR_DECLARATION_TYPE_MISMATCH))
}

// Java: TypedScopeCreatorTest.testGoogModuleIdAttachment_underGoogModule()
func (c *TypedScopeCreatorTest) TestGoogModuleIdAttachment_underGoogModule() {
	c.TestSame(
		c.Srcs(
			"goog.module('my.module');\n" +
				"\n" +
				"class A { }; CLASS: A;\n" +
				"/** @interface */ class B { }; IFACE: B;\n" +
				"function C() { }; FUNC: C;\n" +
				"/** @enum {number} */ const D = {}; ENUM: D;\n" +
				"CAST: /** @type {function()} */ (function() { });\n" +
				"\n" +
				"function scope() {\n" +
				"  class A { }; SCOPED_CLASS: A;\n" +
				"  /** @interface */ class B { }; SCOPED_IFACE: B;\n" +
				"  function C() { }; SCOPED_FUNC: C;\n" +
				"  /** @enum {number} */ const D = {}; SCOPED_ENUM: D;\n" +
				"  SCOPED_CAST: /** @type {function()} */ (function() { });\n" +
				"}\n"))

	c.AssertAttachedGoogModuleId("CLASS", jutil.StringPtr("my.module"))
	c.AssertAttachedGoogModuleId("IFACE", jutil.StringPtr("my.module"))
	c.AssertAttachedGoogModuleId("FUNC", jutil.StringPtr("my.module"))
	c.AssertAttachedGoogModuleId("ENUM", jutil.StringPtr("my.module"))
	c.AssertAttachedGoogModuleId("CAST", nil)

	c.AssertAttachedGoogModuleId("SCOPED_CLASS", jutil.StringPtr("my.module"))
	c.AssertAttachedGoogModuleId("SCOPED_IFACE", jutil.StringPtr("my.module"))
	c.AssertAttachedGoogModuleId("SCOPED_FUNC", jutil.StringPtr("my.module"))
	c.AssertAttachedGoogModuleId("SCOPED_ENUM", jutil.StringPtr("my.module"))
	c.AssertAttachedGoogModuleId("SCOPED_CAST", nil)
}

// Java: TypedScopeCreatorTest.testGoogModuleIdAttachment_underGoogLoadModule()
func (c *TypedScopeCreatorTest) TestGoogModuleIdAttachment_underGoogLoadModule() {
	c.TestSame(
		c.Srcs(
			"goog.loadModule(function(exports) {\n" +
				"  goog.module('my.module.1');\n" +
				"\n" +
				"  class A { }; CLASS1: A;\n" +
				"  /** @interface */ class B { }; IFACE1: B;\n" +
				"  function C() { }; FUNC1: C;\n" +
				"  /** @enum {number} */ const D = {}; ENUM1: D;\n" +
				"  CAST1: /** @type {function()} */ (function() { });\n" +
				"\n" +
				"  return exports;\n" +
				"})\n" +
				"\n" +
				"goog.loadModule(function(exports) {\n" +
				"  goog.module('my.module.2');\n" +
				"\n" +
				"  class A { }; CLASS2: A;\n" +
				"  /** @interface */ class B { }; IFACE2: B;\n" +
				"  function C() { }; FUNC2: C;\n" +
				"  /** @enum {number} */ const D = {}; ENUM2: D;\n" +
				"  CAST2: /** @type {function()} */ (function() { });\n" +
				"\n" +
				"  return exports;\n" +
				"})\n"))

	c.AssertAttachedGoogModuleId("CLASS1", jutil.StringPtr("my.module.1"))
	c.AssertAttachedGoogModuleId("IFACE1", jutil.StringPtr("my.module.1"))
	c.AssertAttachedGoogModuleId("FUNC1", jutil.StringPtr("my.module.1"))
	c.AssertAttachedGoogModuleId("ENUM1", jutil.StringPtr("my.module.1"))
	c.AssertAttachedGoogModuleId("CAST1", nil)

	c.AssertAttachedGoogModuleId("CLASS2", jutil.StringPtr("my.module.2"))
	c.AssertAttachedGoogModuleId("IFACE2", jutil.StringPtr("my.module.2"))
	c.AssertAttachedGoogModuleId("FUNC2", jutil.StringPtr("my.module.2"))
	c.AssertAttachedGoogModuleId("ENUM2", jutil.StringPtr("my.module.2"))
	c.AssertAttachedGoogModuleId("CAST2", nil)
}

// Java: TypedScopeCreatorTest.testGoogModuleIdAttachment_underGoogProvide()
func (c *TypedScopeCreatorTest) TestGoogModuleIdAttachment_underGoogProvide() {
	c.TestSame(
		c.Srcs(
			"goog.provide('my.module');\n" +
				"\n" +
				"class A { }; CLASS: A;\n" +
				"/** @interface */ class B { }; IFACE: B;\n" +
				"function C() { }; FUNC: C;\n" +
				"/** @enum {number} */ const D = {}; ENUM: D;\n" +
				"CAST: /** @type {function()} */ (function() { });\n" +
				"\n" +
				"function scope() {\n" +
				"  class A { }; SCOPED_CLASS: A;\n" +
				"  /** @interface */ class B { }; SCOPED_IFACE: B;\n" +
				"  function C() { }; SCOPED_FUNC: C;\n" +
				"  /** @enum {number} */ const D = {}; SCOPED_ENUM: D;\n" +
				"  SCOPED_CAST: /** @type {function()} */ (function() { });\n" +
				"}\n"))

	c.AssertAttachedGoogModuleId("CLASS", nil)
	c.AssertAttachedGoogModuleId("IFACE", nil)
	c.AssertAttachedGoogModuleId("FUNC", nil)
	c.AssertAttachedGoogModuleId("ENUM", nil)
	c.AssertAttachedGoogModuleId("CAST", nil)

	c.AssertAttachedGoogModuleId("SCOPED_CLASS", nil)
	c.AssertAttachedGoogModuleId("SCOPED_IFACE", nil)
	c.AssertAttachedGoogModuleId("SCOPED_FUNC", nil)
	c.AssertAttachedGoogModuleId("SCOPED_ENUM", nil)
	c.AssertAttachedGoogModuleId("SCOPED_CAST", nil)
}

// Java: TypedScopeCreatorTest.testGoogModuleIdAttachment_underEsModule()
func (c *TypedScopeCreatorTest) TestGoogModuleIdAttachment_underEsModule() {
	c.TestSame(
		c.Srcs(
			"export {};\n" +
				"\n" +
				"class A { }; CLASS: A;\n" +
				"/** @interface */ class B { }; IFACE: B;\n" +
				"function C() { }; FUNC: C;\n" +
				"/** @enum {number} */ const D = {}; ENUM: D;\n" +
				"CAST: /** @type {function()} */ (function() { });\n" +
				"\n" +
				"function scope() {\n" +
				"  class A { }; SCOPED_CLASS: A;\n" +
				"  /** @interface */ class B { }; SCOPED_IFACE: B;\n" +
				"  function C() { }; SCOPED_FUNC: C;\n" +
				"  /** @enum {number} */ const D = {}; SCOPED_ENUM: D;\n" +
				"  SCOPED_CAST: /** @type {function()} */ (function() { });\n" +
				"}\n"))

	c.AssertAttachedGoogModuleId("CLASS", nil)
	c.AssertAttachedGoogModuleId("IFACE", nil)
	c.AssertAttachedGoogModuleId("FUNC", nil)
	c.AssertAttachedGoogModuleId("ENUM", nil)
	c.AssertAttachedGoogModuleId("CAST", nil)

	c.AssertAttachedGoogModuleId("SCOPED_CLASS", nil)
	c.AssertAttachedGoogModuleId("SCOPED_IFACE", nil)
	c.AssertAttachedGoogModuleId("SCOPED_FUNC", nil)
	c.AssertAttachedGoogModuleId("SCOPED_ENUM", nil)
	c.AssertAttachedGoogModuleId("SCOPED_CAST", nil)
}

// Java: TypedScopeCreatorTest.testGoogModuleIdAttachment_underGoogDeclareModuleId()
func (c *TypedScopeCreatorTest) TestGoogModuleIdAttachment_underGoogDeclareModuleId() {
	c.TestSame(
		c.Srcs(
			"export {};\n" +
				"goog.declareModuleId('my.module');\n" +
				"\n" +
				"class A { }; CLASS: A;\n" +
				"/** @interface */ class B { }; IFACE: B;\n" +
				"function C() { }; FUNC: C;\n" +
				"/** @enum {number} */ const D = {}; ENUM: D;\n" +
				"CAST: /** @type {function()} */ (function() { });\n" +
				"\n" +
				"function scope() {\n" +
				"  class A { }; SCOPED_CLASS: A;\n" +
				"  /** @interface */ class B { }; SCOPED_IFACE: B;\n" +
				"  function C() { }; SCOPED_FUNC: C;\n" +
				"  /** @enum {number} */ const D = {}; SCOPED_ENUM: D;\n" +
				"  SCOPED_CAST: /** @type {function()} */ (function() { });\n" +
				"}\n"))

	c.AssertAttachedGoogModuleId("CLASS", jutil.StringPtr("my.module"))
	c.AssertAttachedGoogModuleId("IFACE", jutil.StringPtr("my.module"))
	c.AssertAttachedGoogModuleId("FUNC", jutil.StringPtr("my.module"))
	c.AssertAttachedGoogModuleId("ENUM", jutil.StringPtr("my.module"))
	c.AssertAttachedGoogModuleId("CAST", nil)

	c.AssertAttachedGoogModuleId("SCOPED_CLASS", jutil.StringPtr("my.module"))
	c.AssertAttachedGoogModuleId("SCOPED_IFACE", jutil.StringPtr("my.module"))
	c.AssertAttachedGoogModuleId("SCOPED_FUNC", jutil.StringPtr("my.module"))
	c.AssertAttachedGoogModuleId("SCOPED_ENUM", jutil.StringPtr("my.module"))
	c.AssertAttachedGoogModuleId("SCOPED_CAST", nil)
}

// Java: TypedScopeCreatorTest.assertAttachedGoogModuleId(String, String)
func (c *TypedScopeCreatorTest) AssertAttachedGoogModuleId(label string, expectedId *string) {
	type_ := c.GetLabeledStatement(label).statementNode.GetOnlyChild().GetJSType()

	var actualId *string
	if type_.IsFunctionType() {
		actualId = type_.ToMaybeFunctionType().GetGoogModuleId()
	} else if type_.IsEnumType() {
		actualId = type_.ToMaybeEnumType().GetGoogModuleId()
	} else {
		panic(jutil.NewAssertionError(jutil.ToJavaString(type_)))
	}

	AssertThat(actualId).IsEqualTo(expectedId)
}

// Java: TypedScopeCreatorTest.testObjectLitWithDeclaredSymbolProperty()
func (c *TypedScopeCreatorTest) TestObjectLitWithDeclaredSymbolProperty() {
	c.TestSame(
		c.Externs(
			NewTestExternsBuilder().AddIterable().Build(),
			"/** @const {symbol}  */\n"+
				"Symbol.foo;\n"),
		c.Srcs(
			"const foo = {\n"+
				"  /** @const */\n"+
				"  [Symbol.foo]: 0\n"+
				"};\n"))

	foo := c.FindNameType("foo", c.globalScope).AssertObjectType()
	symbolFooType := typedScopeCreatorTest_cast[*rhino.KnownSymbolType](c.FindNameType("Symbol.foo", c.globalScope))
	TypeSubject_AssertType(foo).WithTypeOfPropSymbol(symbolFooType).IsNumber()
}

// Java: TypedScopeCreatorTest.testForwardReferenceInUnion_extendedInterface_regressionTest()
func (c *TypedScopeCreatorTest) TestForwardReferenceInUnion_extendedInterface_regressionTest() {
	// Regression test for b/497046252#comment11.
	c.TestSame(
		c.Srcs(
			`/** @record */
class AtRecord {
  /** @return {!AtRecord} */
  getAtRecord() {}
}

/** @type {!Clazz} */
let clazzInstance;

/** @type {!AtRecord|!Clazz} */
let union;

/** @interface @extends {Parent} */
class Child {}

/** @implements {Child} */
class Clazz {
  /** @type {string} */
  x = 'hi';
}

/** @interface */
class Parent {}
`))
	unionType := c.FindNameType("union", c.globalScope)
	TypeSubject_AssertType(unionType).ToStringIsEqualTo("(AtRecord|Clazz)")
}

// Java: TypedScopeCreatorTest.testMemoization()
func (c *TypedScopeCreatorTest) TestMemoization() {
	root1 := TypedScopeCreatorTest_CreateEmptyRoot()
	root2 := TypedScopeCreatorTest_CreateEmptyRoot()
	compiler := NewCompiler()
	compiler.InitOptions(NewCompilerOptions())
	creator := NewTypedScopeCreator(compiler)
	scopeA := creator.CreateScope(root1, nil).(*TypedScope)
	AssertThat(creator.CreateScope(root1, nil)).IsSameInstanceAs(scopeA)
	AssertThat(creator.CreateScope(root2, nil)).IsNotSameInstanceAs(scopeA)
}

// Java: TypedScopeCreatorTest.testMemoizationPreconditionCheck()
func (c *TypedScopeCreatorTest) TestMemoizationPreconditionCheck() {
	compiler := NewCompiler()
	compiler.InitOptions(NewCompilerOptions())
	root := TypedScopeCreatorTest_CreateEmptyRoot()
	creator := NewTypedScopeCreator(compiler)
	scopeA := creator.CreateScope(root, nil).(*TypedScope)

	jutil.Try(func() {
		creator.CreateScope(root, scopeA)
		AssertWithMessage("Expected an IllegalStateException").Fail()
	}, func(expected *jutil.IllegalStateException) {
	})
}

// Java: TypedScopeCreatorTest.testFunctionPrototypeAssignmentType()
func (c *TypedScopeCreatorTest) TestFunctionPrototypeAssignmentType() {
	c.TestSame(
		`/**
 * @param {!Function} ctor
 */
function markImplementor(ctor) {
  ctor.prototype.$implements__Foo = true;
}
`)

	ctorPrototypeNode := c.FindQualifiedNameNode("ctor.prototype", c.globalScope.GetRootNode())
	ctorPrototypeType := jutil.CheckNotNull(ctorPrototypeNode.GetJSType())
	nativeFnInstancePrototype :=
		c.GetNativeObjectType(rhino.JSTypeNative_FUNCTION_INSTANCE_PROTOTYPE)

	AssertThat(ctorPrototypeType.ToObjectType().GetReferenceName()).IsEqualTo("?.prototype")
	AssertThat(ctorPrototypeType).IsSameInstanceAs(nativeFnInstancePrototype)
}

// Java: TypedScopeCreatorTest.testFunctionPrototypeAssignmentPreserved_b253690550()
func (c *TypedScopeCreatorTest) TestFunctionPrototypeAssignmentPreserved_b253690550() {
	c.TestSame(
		`/**
 * @param {!Function} childCtor
 * @param {!Function} parentCtor
 */
function inherits(childCtor, parentCtor) {
  // see b/253690550 - this used to clobber Function.prototype. This is a regression test
  // to make sure we preserve the native JSTypeRegistry prototype.
  childCtor.prototype = Object.create(parentCtor.prototype);
}

/**
 * @param {!Function} ctor
 */
function markImplementor(ctor) {
  ctor.prototype.$implements__Foo = true;
}
`)

	ctorPrototypeNode := c.FindQualifiedNameNode("ctor.prototype", c.globalScope.GetRootNode())
	ctorPrototypeType := jutil.CheckNotNull(ctorPrototypeNode.GetJSType())
	nativeFnInstancePrototype :=
		c.GetNativeObjectType(rhino.JSTypeNative_FUNCTION_INSTANCE_PROTOTYPE)

	AssertThat(ctorPrototypeType.ToObjectType().GetReferenceName()).IsEqualTo("?.prototype")
	AssertThat(ctorPrototypeType).IsSameInstanceAs(nativeFnInstancePrototype)
}

// Java: TypedScopeCreatorTest.createEmptyRoot()
func TypedScopeCreatorTest_CreateEmptyRoot() rhino.Node {
	script := rhino.IR_Script()
	root := rhino.IR_Root(rhino.IR_Root(), rhino.IR_Root(script))
	script.SetInputId(rhino.NewInputId("input"))
	return root
}

// Java: TypedScopeCreatorTest.findNameType(String, TypedScope)
func (c *TypedScopeCreatorTest) FindNameType(name string, scope *TypedScope) rhino.JSType {
	return c.FindTypeOnMatchedNode(func(n rhino.Node) bool { return n.MatchesQualifiedName(name) }, scope)
}

// Java: TypedScopeCreatorTest.findNameTypeStr(String, TypedScope)
func (c *TypedScopeCreatorTest) FindNameTypeStr(name string, scope *TypedScope) string {
	return c.FindNameType(name, scope).String()
}

// Java: TypedScopeCreatorTest.findTokenType(Token, TypedScope)
func (c *TypedScopeCreatorTest) FindTokenType(type_ rhino.Token, scope *TypedScope) rhino.JSType {
	return c.FindTypeOnMatchedNode(func(n rhino.Node) bool { return type_ == n.GetToken() }, scope)
}

// Java: TypedScopeCreatorTest.findTypeOnMatchedNode(Predicate, TypedScope)
func (c *TypedScopeCreatorTest) FindTypeOnMatchedNode(matcher func(rhino.Node) bool, scope *TypedScope) rhino.JSType {
	root := scope.GetRootNode()
	// Go-only: language-adapter — the ArrayDeque used as a stack
	// (push/pop at the head) is a slice whose top is its end.
	queue := []rhino.Node{root}
	for len(queue) != 0 {
		current := queue[len(queue)-1]
		queue = queue[:len(queue)-1]
		if matcher(current) && current.GetJSType() != nil {
			return current.GetJSType()
		}

		for child := current.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
			queue = append(queue, child)
		}
	}
	return nil
}

// Java: TypedScopeCreatorTest.getNativeType(JSTypeNative)
func (c *TypedScopeCreatorTest) GetNativeType(type_ rhino.JSTypeNative) rhino.JSType {
	return c.registry.GetNativeType(type_)
}

// Java: TypedScopeCreatorTest.getNativeObjectType(JSTypeNative)
func (c *TypedScopeCreatorTest) GetNativeObjectType(type_ rhino.JSTypeNative) rhino.ObjectType {
	return c.registry.GetNativeType(type_).(rhino.ObjectType)
}

// typedVarOf is the Java TypedScope.getSlot/getOwnSlot result type: the Go
// methods return a rhino.StaticSlot.
//
// Go-only: language-adapter — a nil slot stays a nil *TypedVar.
func typedVarOf(slot rhino.StaticSlot) *TypedVar {
	v, _ := slot.(*TypedVar)
	return v
}

// typedScopeCreatorTest_cast is the Java reference cast `(T) v`: a null
// reference stays null, any other value of the wrong type throws a
// ClassCastException.
//
// Go-only: language-adapter — a Go type assertion panics on a nil
// interface.
func typedScopeCreatorTest_cast[T any](v any) T {
	var null T
	if jutil.IsNil(v) {
		return null
	}
	return v.(T)
}
