// Port of com/google/javascript/jscomp/TypeInferenceTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/TypeInferenceTest.java)

/*
 * Copyright 2008 The Closure Compiler Authors.
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
	"strconv"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/rhinotesting"
)

// TypeInferenceTest tests TypeInference.
//
// These unit tests don't test the full type inference that happens in a
// normal compilation. This is because a normal compilation creates many
// instances of the TypeInference class. More precisely, one instance is
// created for every scope that is a valid root of a control-flow graph (as
// defined by NodeUtil#isValidCfgRoot(Node) such as, for example, functions.
//
// These unit tests only ever create a single TypeInference instance. That
// means: these unit tests ignore any code nested within a function.
//
// To write tests that more fully model the typechecking in a full
// compilation, which visits every function body, use TypedScopeCreatorTest
// or TypeCheckTest.
type TypeInferenceTest struct {
	typeCheckTests_noFixture

	compiler        *Compiler
	registry        *rhino.JSTypeRegistry
	closer          *rhino.JSTypeResolver_Closer
	assumptions     *jutil.JavaHashMap[string, rhino.JSType]
	assumedThisType rhino.JSType
	returnScope     FlowScope

	// Maps a label name to information about the labeled statement.
	//
	// This map is recreated each time parseAndRunTypeInference() is executed.
	labeledStatementMap *jutil.JavaHashMap[string, *TypeInferenceTest_LabeledStatement]
}

// Java: TypeInferenceTest()
func NewTypeInferenceTest() *TypeInferenceTest {
	return &TypeInferenceTest{}
}

// JUnitBefore runs the @Before method setUp().
func (c *TypeInferenceTest) JUnitBefore() { c.SetUp() }

func TestTypeInferenceTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewTypeInferenceTest() })
}

var TypeInferenceTest_ASSERTION_FUNCTION_MAP = CodingConvention_AssertionFunctionLookup_Of(NewClosureCodingConvention().GetAssertionFunctions())

// TypeInferenceTest_LabeledStatement stores information about a labeled
// statement and allows making assertions on it.
type TypeInferenceTest_LabeledStatement struct {
	statementNode  rhino.Node
	enclosingScope *TypedScope
}

// Java: TypeInferenceTest.LabeledStatement(Node, TypedScope)
func NewTypeInferenceTest_LabeledStatement(statementNode rhino.Node, enclosingScope *TypedScope) *TypeInferenceTest_LabeledStatement {
	s := &TypeInferenceTest_LabeledStatement{}
	s.statementNode = rhino.CheckNotNull(statementNode)
	s.enclosingScope = jutil.CheckNotNull(enclosingScope)
	return s
}

// Java: TypeInferenceTest.setUp()
func (c *TypeInferenceTest) SetUp() {
	c.compiler = NewCompiler()
	options := NewCompilerOptions()
	options.SetClosurePass(true)
	options.SetLanguageIn(CompilerOptions_LanguageMode_UNSUPPORTED)
	c.compiler.InitOptions(options)
	c.registry = c.compiler.GetTypeRegistry()
	c.assumptions = jutil.NewJavaHashMap[string, rhino.JSType]()
	c.returnScope = nil
	c.closer = c.registry.GetResolver().OpenForDefinition()
}

// Java: TypeInferenceTest.assumingThisType(JSType)
func (c *TypeInferenceTest) AssumingThisType(type_ rhino.JSType) {
	c.assumedThisType = type_
}

// Java: TypeInferenceTest.assuming(String, JSType)
func (c *TypeInferenceTest) Assuming(name string, type_ rhino.JSType) {
	c.assumptions.Put(name, type_)
}

// AssumingNative declares a name with a given type in the parent scope of
// the test case code.
//
// Java: TypeInferenceTest.assuming(String, JSTypeNative)
func (c *TypeInferenceTest) AssumingNative(name string, type_ rhino.JSTypeNative) {
	c.Assuming(name, c.registry.GetNativeType(type_))
}

// InFunction runs an instance of TypeInference over the given code after
// wrapping it in a function.
//
// Does not visit any nested functions.
//
// Java: TypeInferenceTest.inFunction(String)
func (c *TypeInferenceTest) InFunction(js string) {
	// Parse the body of the function.
	thisBlock := ""
	if c.assumedThisType != nil {
		thisBlock = "/** @this {" + c.assumedThisType.String() + "} */"
	}
	c.ParseAndRunTypeInference("(" + thisBlock + " function() {" + js + "});")
}

// InScript runs an instance of TypeInference over the given code.
//
// Does not visit any nested functions.
//
// Java: TypeInferenceTest.inScript(String)
func (c *TypeInferenceTest) InScript(js string) {
	script := c.compiler.ParseTestCode(js)
	AssertWithMessage("%s", "parsing error: "+typeCheckTests_join(c.compiler.GetErrors(), ", ")).
		That(c.compiler.GetErrorCount()).
		IsEqualTo(0)
	root := rhino.IR_Root(rhino.IR_Root(), rhino.IR_Root(script))
	NewGatherModuleMetadata(c.compiler /* processCommonJsModules= */, false, ModuleLoader_ResolutionMode_BROWSER).
		Process(root.GetFirstChild(), root.GetSecondChild())
	NewModuleMapCreator(c.compiler, c.compiler.GetModuleMetadataMap()).
		Process(root.GetFirstChild(), root.GetSecondChild())

	c.ParseAndRunTypeInferenceOfNode(root, root)
}

// InGenerator runs an instance of TypeInference over the given code after
// wrapping it in a generator.
//
// Does not visit any nested functions.
//
// Java: TypeInferenceTest.inGenerator(String)
func (c *TypeInferenceTest) InGenerator(js string) {
	jutil.CheckState(c.assumedThisType == nil)
	c.ParseAndRunTypeInference("(function *() {" + js + "});")
}

// Java: TypeInferenceTest.withModules(ImmutableList<String>)
func (c *TypeInferenceTest) WithModules(js []string) {
	script := c.compiler.ParseTestCodeOfList(js)
	chunk := NewJSChunk("entry")
	for s := script.GetFirstChild(); !s.IsNil(); s = s.GetNext() {
		chunk.AddSourceFile(c.compiler.GetInput(s.GetInputId()).GetSourceFile())
	}
	c.compiler.InitChunks([]*SourceFile{}, []*JSChunk{chunk}, c.compiler.GetOptions())
	c.compiler.InitializeModuleLoader()
	AssertWithMessage("%s", "parsing error: "+typeCheckTests_join(c.compiler.GetErrors(), ", ")).
		That(c.compiler.GetErrorCount()).
		IsEqualTo(0)
	root := rhino.IR_Root(rhino.IR_Root(), rhino.IR_Root(script))
	NewGatherModuleMetadata(c.compiler /* processCommonJsModules= */, false, ModuleLoader_ResolutionMode_BROWSER).
		Process(root.GetFirstChild(), root.GetSecondChild())
	NewModuleMapCreator(c.compiler, c.compiler.GetModuleMetadataMap()).
		Process(root.GetFirstChild(), root.GetSecondChild())

	c.ParseAndRunTypeInferenceOfNode(root, root)
}

// Java: TypeInferenceTest.parseAndRunTypeInference(String)
func (c *TypeInferenceTest) ParseAndRunTypeInference(js string) {
	script := c.compiler.ParseTestCode(js)
	root := rhino.IR_Root(rhino.IR_Root(), rhino.IR_Root(script))
	AssertWithMessage("%s", "parsing error: "+typeCheckTests_join(c.compiler.GetErrors(), ", ")).
		That(c.compiler.GetErrorCount()).
		IsEqualTo(0)

	// SCRIPT -> EXPR_RESULT -> FUNCTION
	// `(function() { TEST CODE HERE });`
	function := script.GetFirstFirstChild()
	c.ParseAndRunTypeInferenceOfNode(root, function)
}

// typeInferenceTest_labelCollector is the anonymous AbstractScopedCallback
// of parseAndRunTypeInference(Node, Node): it creates the scopes and
// populates a map allowing us to look up labeled statements later.
type typeInferenceTest_labelCollector struct {
	NodeTraversal_AbstractScopedCallbackBase
	outer *TypeInferenceTest
}

// Java: TypeInferenceTest.parseAndRunTypeInference(Node, Node) (anonymous AbstractScopedCallback).enterScope(NodeTraversal)
func (cb *typeInferenceTest_labelCollector) EnterScope(t *NodeTraversal) {
	t.GetTypedScope()
}

// Java: TypeInferenceTest.parseAndRunTypeInference(Node, Node) (anonymous AbstractScopedCallback).visit(NodeTraversal, Node, Node)
func (cb *typeInferenceTest_labelCollector) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	scope := t.GetTypedScope()
	if !parent.IsNil() && parent.IsLabel() && !n.IsLabelName() {
		// First child of a LABEL is a LABEL_NAME, n is the second child.
		labelNameNode := rhino.CheckNotNullMsg1(n.GetPrevious(), "%s", n)
		jutil.CheckStateMsg1(labelNameNode.IsLabelName(), "%s", labelNameNode)
		labelName := labelNameNode.GetString()
		AssertWithMessage("Duplicate label name: %s", labelName).
			That(cb.outer.labeledStatementMap).
			DoesNotContainKey(labelName)
		cb.outer.labeledStatementMap.Put(labelName, NewTypeInferenceTest_LabeledStatement(n, scope))
	}
}

// Java: TypeInferenceTest.parseAndRunTypeInference(Node, Node)
func (c *TypeInferenceTest) ParseAndRunTypeInferenceOfNode(root rhino.Node, cfgRoot rhino.Node) {
	c.closer.Close()

	scopeCreator := NewTypedScopeCreator(c.compiler)
	var assumedScope *TypedScope
	func() {
		closer := c.registry.GetResolver().OpenForDefinition()
		defer closer.Close()
		// Create the scope with the assumptions.
		// Also populate a map allowing us to look up labeled statements later.
		c.labeledStatementMap = jutil.NewJavaHashMap[string, *TypeInferenceTest_LabeledStatement]()
		NodeTraversal_Builder_().
			SetCompiler(c.compiler).
			SetCallback(&typeInferenceTest_labelCollector{outer: c}).
			SetScopeCreator(scopeCreator).
			Traverse(root)
		assumedScope = scopeCreator.CreateScopeOfNode(cfgRoot)
		for name, type_ := range c.assumptions.All() {
			assumedScope.Declare(name, rhino.NilNode, type_, nil, false)
		}
		scopeCreator.ResolveWeakImportsPreResolution()
	}()
	scopeCreator.FinishAndFreeze()
	// Create the control graph.
	cfg :=
		ControlFlowAnalysis_Builder_().
			SetCompiler(c.compiler).
			SetCfgRoot(cfgRoot).
			SetTraverseFunctions(true).
			SetIncludeEdgeAnnotations(true).
			ComputeCfg()
	// Create a simple reverse abstract interpreter.
	rai := c.compiler.GetReverseAbstractInterpreter()
	// Do the type inference by data-flow analysis.
	dfa :=
		NewTypeInference(c.compiler, cfg, rai, assumedScope, scopeCreator, TypeInferenceTest_ASSERTION_FUNCTION_MAP)
	dfa.Analyze()
	// Get the scope of the implicit return.
	rtnState := cfg.GetImplicitReturn().GetAnnotation().(*DataFlowAnalysis_LinearFlowState[FlowScope])
	if cfgRoot.IsFunction() {
		// Reset the flow scope's syntactic scope to the function block, rather than the function
		// node
		// itself.  This allows pulling out local vars from the function by name to verify their
		// types.
		c.returnScope =
			rtnState.GetIn().WithSyntacticScope(scopeCreator.CreateScopeOfNode(cfgRoot.GetLastChild()))
	} else {
		c.returnScope = rtnState.GetIn()
	}

	c.closer = c.registry.GetResolver().OpenForDefinition()
}

// Java: TypeInferenceTest.getLabeledStatement(String)
func (c *TypeInferenceTest) GetLabeledStatement(label string) *TypeInferenceTest_LabeledStatement {
	AssertWithMessage("No statement found for label: %s", label).
		That(c.labeledStatementMap).
		ContainsKey(label)
	return c.labeledStatementMap.Get(label)
}

// AssertScopeEnclosing returns a ScopeSubject for the scope containing the
// labeled statement.
//
// Asserts that a statement with the given label existed in the code last
// passed to parseAndRunTypeInference().
//
// Java: TypeInferenceTest.assertScopeEnclosing(String)
func (c *TypeInferenceTest) AssertScopeEnclosing(label string) *ScopeSubject {
	return ScopeSubject_AssertScope(c.GetLabeledStatement(label).enclosingScope)
}

// AssertTypeOfExpression returns a TypeSubject for the JSType of the
// expression with the given label.
//
// Asserts that a statement with the given label existed in the code last
// passed to parseAndRunTypeInference(). Also asserts that the statement is an
// EXPR_RESULT whose expression has a non-null JSType.
//
// Java: TypeInferenceTest.assertTypeOfExpression(String)
func (c *TypeInferenceTest) AssertTypeOfExpression(label string) *rhinotesting.TypeSubject {
	statementNode := c.GetLabeledStatement(label).statementNode
	AssertWithMessage("Not an expression statement.").That(statementNode.IsExprResult()).IsTrue()
	jsType := statementNode.GetOnlyChild().GetJSType()
	AssertWithMessage("Expression type is null").That(jsType).IsNotNull()
	return TypeSubject_AssertType(jsType)
}

// Java: TypeInferenceTest.getType(String)
func (c *TypeInferenceTest) GetType(name string) rhino.JSType {
	AssertWithMessage("The return scope should not be null.").That(c.returnScope).IsNotNull()
	var_, _ := c.returnScope.GetSlot(name).(rhino.StaticTypedSlot)
	AssertWithMessage("The variable %s is missing from the scope.", name).That(var_).IsNotNull()
	return var_.GetType()
}

// GetParamNameNode returns the NAME node `name` from the PARAM_LIST of the
// top level of type inference.
//
// Java: TypeInferenceTest.getParamNameNode(String)
func (c *TypeInferenceTest) GetParamNameNode(name string) rhino.Node {
	staticScope := jutil.CheckNotNullMsg1(c.returnScope.GetDeclarationScope(), "%s", c.returnScope)
	slot := jutil.CheckNotNullMsg1(staticScope.GetSlot(name), "%s", staticScope)
	declaration := jutil.CheckNotNullMsg1(slot.GetDeclaration(), "%s", slot)
	node := rhino.CheckNotNullMsg1(declaration.GetNode(), "%s", declaration)

	NodeSubject_AssertNode(node).HasType(rhino.Token_NAME)
	var paramLists []rhino.Node
	for ancestor := range node.GetAncestors().All() {
		if ancestor.IsParamList() {
			paramLists = append(paramLists, ancestor)
		}
	}
	AssertThat(paramLists).IsNotEmpty()

	return node
}

// Java: TypeInferenceTest.verify(String, JSType)
func (c *TypeInferenceTest) Verify(name string, type_ rhino.JSType) {
	rhinotesting.NewTruth(assertionTB{}).
		AssertWithMessage("Mismatch for %s", name).
		AboutTypes().
		That(c.GetType(name)).
		IsEqualTo(type_)
}

// Java: TypeInferenceTest.verify(String, JSTypeNative)
func (c *TypeInferenceTest) VerifyNative(name string, type_ rhino.JSTypeNative) {
	c.Verify(name, c.registry.GetNativeType(type_))
}

// Java: TypeInferenceTest.verifyUnequal(String, JSType)
func (c *TypeInferenceTest) VerifyUnequal(name string, type_ rhino.JSType) {
	rhinotesting.NewTruth(assertionTB{}).
		AssertWithMessage("Mismatch for %s", name).
		AboutTypes().
		That(c.GetType(name)).
		IsNotEqualTo(type_)
}

// Java: TypeInferenceTest.verifyUnequal(String, JSTypeNative)
func (c *TypeInferenceTest) VerifyUnequalNative(name string, type_ rhino.JSTypeNative) {
	c.VerifyUnequal(name, c.registry.GetNativeType(type_))
}

// Java: TypeInferenceTest.verifySubtypeOf(String, JSType)
func (c *TypeInferenceTest) VerifySubtypeOf(name string, type_ rhino.JSType) {
	varType := c.GetType(name)
	AssertWithMessage("The variable %s is missing a type.", name).That(varType).IsNotNull()
	AssertWithMessage("The type %s of variable %s is not a subtype of %s.", varType, name, type_).
		That(varType.IsSubtypeOf(type_)).
		IsTrue()
}

// Java: TypeInferenceTest.verifySubtypeOf(String, JSTypeNative)
func (c *TypeInferenceTest) VerifySubtypeOfNative(name string, type_ rhino.JSTypeNative) {
	c.VerifySubtypeOf(name, c.registry.GetNativeType(type_))
}

// Java: TypeInferenceTest.createEnumType(String, JSTypeNative)
func (c *TypeInferenceTest) CreateEnumType(name string, elemType rhino.JSTypeNative) *rhino.EnumType {
	return c.CreateEnumTypeOfJSType(name, c.registry.GetNativeType(elemType))
}

// Java: TypeInferenceTest.createEnumType(String, JSType)
func (c *TypeInferenceTest) CreateEnumTypeOfJSType(name string, elemType rhino.JSType) *rhino.EnumType {
	return rhino.EnumType_Builder_(c.registry).SetName(jutil.StringPtr(name)).SetElementType(elemType).Build()
}

// Java: TypeInferenceTest.createUndefinableType(JSTypeNative)
func (c *TypeInferenceTest) CreateUndefinableType(type_ rhino.JSTypeNative) rhino.JSType {
	return c.registry.CreateUnionType(
		c.registry.GetNativeType(type_), c.registry.GetNativeType(rhino.JSTypeNative_VOID_TYPE))
}

// Java: TypeInferenceTest.createNullableType(JSTypeNative)
func (c *TypeInferenceTest) CreateNullableType(type_ rhino.JSTypeNative) rhino.JSType {
	return c.CreateNullableTypeOfJSType(c.registry.GetNativeType(type_))
}

// Java: TypeInferenceTest.createNullableType(JSType)
func (c *TypeInferenceTest) CreateNullableTypeOfJSType(type_ rhino.JSType) rhino.JSType {
	return c.registry.CreateNullableType(type_)
}

// Java: TypeInferenceTest.createUnionType(JSTypeNative, JSTypeNative)
func (c *TypeInferenceTest) CreateUnionType(type1 rhino.JSTypeNative, type2 rhino.JSTypeNative) rhino.JSType {
	return c.registry.CreateUnionType(c.registry.GetNativeType(type1), c.registry.GetNativeType(type2))
}

// CreateRecordType returns a record type with a field `fieldName` and
// JSType specified by `fieldType`.
//
// Java: TypeInferenceTest.createRecordType(String, JSType)
func (c *TypeInferenceTest) CreateRecordType(fieldName string, fieldType rhino.JSType) rhino.JSType {
	property := jutil.NewImmutableMapBuilder[string, rhino.JSType]().Put(fieldName, fieldType).BuildOrThrow()
	return c.registry.CreateRecordType(property)
}

// Java: TypeInferenceTest.createRecordType(String, JSTypeNative)
func (c *TypeInferenceTest) CreateRecordTypeNative(fieldName string, fieldType rhino.JSTypeNative) rhino.JSType {
	return c.CreateRecordType(fieldName, c.registry.GetNativeType(fieldType))
}

// Java: TypeInferenceTest.createMultiParamUnionType(JSTypeNative...)
func (c *TypeInferenceTest) CreateMultiParamUnionType(variants ...rhino.JSTypeNative) rhino.JSType {
	return c.registry.CreateUnionTypeOfNatives(variants...)
}

// Java: TypeInferenceTest.testAssumption()
func (c *TypeInferenceTest) TestAssumption() {
	c.AssumingNative("x", rhino.JSTypeNative_NUMBER_TYPE)
	c.InFunction("")
	c.VerifyNative("x", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testVar()
func (c *TypeInferenceTest) TestVar() {
	c.InFunction("var x = 1;")
	c.VerifyNative("x", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testMissingTypeAnnotationsInfersUnknownReturnInArrow()
func (c *TypeInferenceTest) TestMissingTypeAnnotationsInfersUnknownReturnInArrow() {
	c.InFunction("const getNum = () => 2; const a = getNum();")
	c.VerifyUnequalNative("a", rhino.JSTypeNative_NUMBER_TYPE)
	c.VerifyNative("a", rhino.JSTypeNative_UNKNOWN_TYPE)
}

// Java: TypeInferenceTest.testMissingTypeAnnotationsInfersUnknownReturn()
func (c *TypeInferenceTest) TestMissingTypeAnnotationsInfersUnknownReturn() {
	c.InFunction("const getNum = function() {return 2; }; const a = getNum();")
	c.VerifyUnequalNative("a", rhino.JSTypeNative_NUMBER_TYPE)
	c.VerifyNative("a", rhino.JSTypeNative_UNKNOWN_TYPE)
}

// Java: TypeInferenceTest.testEmptyVar()
func (c *TypeInferenceTest) TestEmptyVar() {
	c.InFunction("var x;")
	c.VerifyNative("x", rhino.JSTypeNative_VOID_TYPE)
}

// Java: TypeInferenceTest.testAssignment()
func (c *TypeInferenceTest) TestAssignment() {
	c.AssumingNative("x", rhino.JSTypeNative_OBJECT_TYPE)
	c.InFunction("x = 1;")
	c.VerifyNative("x", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testExprWithinCast()
func (c *TypeInferenceTest) TestExprWithinCast() {
	c.AssumingNative("x", rhino.JSTypeNative_OBJECT_TYPE)
	c.InFunction("/** @type {string} */ (x = 1);")
	c.VerifyNative("x", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testGetProp()
func (c *TypeInferenceTest) TestGetProp() {
	c.Assuming("x", c.CreateNullableType(rhino.JSTypeNative_OBJECT_TYPE))
	c.InFunction("x.y();")
	c.VerifyNative("x", rhino.JSTypeNative_OBJECT_TYPE)
}

// Java: TypeInferenceTest.testOptChainGetProp_nullObject()
func (c *TypeInferenceTest) TestOptChainGetProp_nullObject() {
	c.InFunction("let x = null; let a = x?.y;")
	c.VerifyNative("a", rhino.JSTypeNative_VOID_TYPE)
}

// Java: TypeInferenceTest.testOptChainGetElem_accessedByName()
func (c *TypeInferenceTest) TestOptChainGetElem_accessedByName() {
	c.InFunction("let x = { y : 5}; let a = x?.[y];")
	c.VerifyNative("a", rhino.JSTypeNative_UNKNOWN_TYPE)
}

// Java: TypeInferenceTest.testOptChainGetElem_accessedByString()
func (c *TypeInferenceTest) TestOptChainGetElem_accessedByString() {
	c.InFunction("let x = { y : 5}; let a = x?.['y'];")
	c.VerifyNative("a", rhino.JSTypeNative_UNKNOWN_TYPE)
}

// Java: TypeInferenceTest.testNormalGetElem_accessedByString()
func (c *TypeInferenceTest) TestNormalGetElem_accessedByString() {
	c.InFunction("let x = { y : 5}; let a = x['y'];")
	c.VerifyNative("a", rhino.JSTypeNative_UNKNOWN_TYPE)
}

// Java: TypeInferenceTest.testSimpleGetProp_missingPropAccessedOnRecordType()
func (c *TypeInferenceTest) TestSimpleGetProp_missingPropAccessedOnRecordType() {
	c.Assuming("x", c.CreateRecordTypeNative("y", rhino.JSTypeNative_STRING_TYPE))
	c.InFunction("a = x.z;")
	c.VerifyNative("a", rhino.JSTypeNative_UNKNOWN_TYPE)
}

// Java: TypeInferenceTest.testSimpleGetProp_missingPropAccessedOnStringType()
func (c *TypeInferenceTest) TestSimpleGetProp_missingPropAccessedOnStringType() {
	c.AssumingNative("x", rhino.JSTypeNative_STRING_TYPE)
	c.InFunction("a = x.z;")
	c.VerifyNative("a", rhino.JSTypeNative_UNKNOWN_TYPE)
}

// Java: TypeInferenceTest.testOptChainGetProp_nullObj()
func (c *TypeInferenceTest) TestOptChainGetProp_nullObj() {
	c.AssumingNative("x", rhino.JSTypeNative_NULL_TYPE)
	c.InFunction("let a = x?.y;")
	c.VerifyNative("a", rhino.JSTypeNative_VOID_TYPE)
}

// Java: TypeInferenceTest.testOptChain_accessingInexistentPropOnUnknownType()
func (c *TypeInferenceTest) TestOptChain_accessingInexistentPropOnUnknownType() {
	c.AssumingNative("x", rhino.JSTypeNative_UNKNOWN_TYPE)
	c.InFunction("a = x?.z")
	c.VerifyNative("a", rhino.JSTypeNative_UNKNOWN_TYPE)
}

// Java: TypeInferenceTest.testOptChainGetProp_stringProp()
func (c *TypeInferenceTest) TestOptChainGetProp_stringProp() {
	c.Assuming("x", c.CreateRecordTypeNative("y", rhino.JSTypeNative_STRING_TYPE))
	c.InFunction("let a = x?.y;")
	c.VerifyNative("a", rhino.JSTypeNative_STRING_TYPE)
}

// Java: TypeInferenceTest.testOptChainGetProp_nullableStringProp()
func (c *TypeInferenceTest) TestOptChainGetProp_nullableStringProp() {
	c.Assuming("x", c.CreateRecordType("y", c.CreateUnionType(rhino.JSTypeNative_NULL_TYPE, rhino.JSTypeNative_STRING_TYPE)))
	c.InFunction("let a = x?.y;")
	c.Verify("a", c.CreateUnionType(rhino.JSTypeNative_NULL_TYPE, rhino.JSTypeNative_STRING_TYPE))
}

// Java: TypeInferenceTest.testOptChainGetProp_functionProp()
func (c *TypeInferenceTest) TestOptChainGetProp_functionProp() {
	var funcType rhino.JSType = c.registry.CreateFunctionType(c.registry.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE))
	var lhsType rhino.JSType = c.CreateRecordType("y", funcType)
	c.Assuming("x", lhsType)
	c.InFunction("let a = x?.y;")
	c.Verify("x", lhsType)
	c.Verify("a", funcType)
}

// Java: TypeInferenceTest.testOptChainGetProp_functionProp_nullableObj()
func (c *TypeInferenceTest) TestOptChainGetProp_functionProp_nullableObj() {
	c.Assuming("x", c.CreateNullableTypeOfJSType(c.CreateRecordTypeNative("y", rhino.JSTypeNative_FUNCTION_TYPE)))
	c.InFunction("let a = x?.y;")
	c.Verify("a", c.CreateUnionType(rhino.JSTypeNative_VOID_TYPE, rhino.JSTypeNative_FUNCTION_TYPE))
}

// Java: TypeInferenceTest.testOptChainGetProp_stringProp_nullableObj()
func (c *TypeInferenceTest) TestOptChainGetProp_stringProp_nullableObj() {
	c.Assuming("x", c.CreateNullableTypeOfJSType(c.CreateRecordTypeNative("y", rhino.JSTypeNative_STRING_TYPE)))
	c.InFunction("let a = x?.y;")
	c.Verify("a", c.CreateUnionType(rhino.JSTypeNative_STRING_TYPE, rhino.JSTypeNative_VOID_TYPE))
}

// Java: TypeInferenceTest.testOptChainGetProp_numberProp()
func (c *TypeInferenceTest) TestOptChainGetProp_numberProp() {
	c.Assuming("x", c.CreateRecordTypeNative("y", rhino.JSTypeNative_NUMBER_TYPE))
	c.InFunction("let a = x?.y;")
	c.VerifyNative("a", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testOptChainGetProp_numberProp_nullableObj()
func (c *TypeInferenceTest) TestOptChainGetProp_numberProp_nullableObj() {
	c.Assuming("x", c.CreateNullableTypeOfJSType(c.CreateRecordTypeNative("y", rhino.JSTypeNative_NUMBER_TYPE)))
	c.InFunction("let a = x?.y;")
	c.Verify("a", c.CreateUnionType(rhino.JSTypeNative_VOID_TYPE, rhino.JSTypeNative_NUMBER_TYPE))
}

// Java: TypeInferenceTest.testOptChainGetProp_nullProp()
func (c *TypeInferenceTest) TestOptChainGetProp_nullProp() {
	c.Assuming("x", c.CreateRecordTypeNative("y", rhino.JSTypeNative_NULL_TYPE))
	c.InFunction("let a = x?.y;")
	c.VerifyNative("a", rhino.JSTypeNative_NULL_TYPE)
}

// Java: TypeInferenceTest.testOptChainGetProp_inexistentProp()
func (c *TypeInferenceTest) TestOptChainGetProp_inexistentProp() {
	c.Assuming("x", c.CreateRecordTypeNative("y", rhino.JSTypeNative_STRING_TYPE))
	c.InFunction("let a = x?.z;")
	c.VerifyNative("a", rhino.JSTypeNative_UNKNOWN_TYPE)
}

// Java: TypeInferenceTest.testSimpleOptChain_withUnsetType_trailingGetProp()
func (c *TypeInferenceTest) TestSimpleOptChain_withUnsetType_trailingGetProp() {
	c.AssumingNative("x", rhino.JSTypeNative_UNKNOWN_TYPE)
	c.InFunction("let a = (x?.y?.z).q;")
	c.VerifyNative("a", rhino.JSTypeNative_UNKNOWN_TYPE)
}

// Java: TypeInferenceTest.testGetProp_withOptionalChainObject_voidable()
func (c *TypeInferenceTest) TestGetProp_withOptionalChainObject_voidable() {
	var recordType rhino.JSType = c.CreateRecordType("y", c.CreateRecordTypeNative("z", rhino.JSTypeNative_STRING_TYPE))
	var lhsType rhino.JSType = c.registry.CreateUnionType(c.registry.GetNativeType(rhino.JSTypeNative_VOID_TYPE), recordType)
	c.Assuming("x", lhsType)
	c.InFunction("let a = (x?.y).z")
	c.Verify("x", recordType)
	c.VerifyNative("a", rhino.JSTypeNative_STRING_TYPE)
}

// Java: TypeInferenceTest.testOptChainGetProp_unconditionalAssignmentToObjInInnerNodes()
func (c *TypeInferenceTest) TestOptChainGetProp_unconditionalAssignmentToObjInInnerNodes() {
	c.Assuming("x", c.CreateRecordType("y", c.registry.CreateFunctionType(c.registry.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE))))
	c.InFunction("x?.y(x=5);")
	c.VerifyNative("x", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testOptChainGetProp_conditionalAssignmentToObjInInnerNodes()
func (c *TypeInferenceTest) TestOptChainGetProp_conditionalAssignmentToObjInInnerNodes() {
	var nullableRecordType rhino.JSType = c.CreateNullableTypeOfJSType(c.CreateRecordType("y", c.registry.CreateFunctionType(c.registry.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE))))
	c.Assuming("x", nullableRecordType)
	c.InFunction("x?.y(x=5);")
	c.Verify("x", c.registry.CreateUnionType(c.registry.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE), c.registry.GetNativeType(rhino.JSTypeNative_NULL_TYPE)))
}

// Java: TypeInferenceTest.testOptChainGetProp_typeAnnotationOnObjInInnerNodes_unconditionalChain()
func (c *TypeInferenceTest) TestOptChainGetProp_typeAnnotationOnObjInInnerNodes_unconditionalChain() {
	var recordType rhino.JSType = c.CreateRecordType("y", c.registry.CreateFunctionType(c.registry.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)))
	c.Assuming("x", recordType)
	c.InFunction("x?.y(/** @type {number} */ (x)); ")
	c.Verify("x", recordType)
}

// Java: TypeInferenceTest.testGetProp_typeAnnotationOnObjInInnerNodes_unconditionalChain()
func (c *TypeInferenceTest) TestGetProp_typeAnnotationOnObjInInnerNodes_unconditionalChain() {
	var recordType rhino.JSType = c.CreateRecordType("y", c.registry.CreateFunctionType(c.registry.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)))
	c.Assuming("x", recordType)
	c.InFunction("x.y(/** @type {number} */ (x)); ")
	c.Verify("x", recordType)
}

// Java: TypeInferenceTest.testOptChainGetProp_typeAnnotationOnObjInInnerNodes_conditionalChain()
func (c *TypeInferenceTest) TestOptChainGetProp_typeAnnotationOnObjInInnerNodes_conditionalChain() {
	var nullableRecordType rhino.JSType = c.CreateNullableTypeOfJSType(c.CreateRecordType("y", c.registry.CreateFunctionType(c.registry.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE))))
	c.Assuming("x", nullableRecordType)
	c.InFunction("x?.y(/** @type {number} */ (x)); ")
	c.Verify("x", nullableRecordType)
}

// Java: TypeInferenceTest.testOptChainGetProp_conditionalChangeToOuterVariable()
func (c *TypeInferenceTest) TestOptChainGetProp_conditionalChangeToOuterVariable() {
	var lhsNullableRecordType rhino.JSType = c.CreateNullableTypeOfJSType(c.CreateRecordType("b", c.registry.CreateFunctionType(c.registry.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE))))
	c.Assuming("a", lhsNullableRecordType)
	c.InFunction(`let x = 'x';
a?.b(x=5);
a;
x;
`)
	c.Verify("a", lhsNullableRecordType)
	c.Verify("x", c.CreateUnionType(rhino.JSTypeNative_NUMBER_TYPE, rhino.JSTypeNative_STRING_TYPE))
}

// Java: TypeInferenceTest.testOptChainGetProp_unconditionalChangeToOuterVariable_inNullableReceiver()
func (c *TypeInferenceTest) TestOptChainGetProp_unconditionalChangeToOuterVariable_inNullableReceiver() {
	var recordType rhino.JSType = c.CreateRecordType("b", c.registry.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE))
	var nullableRecordType rhino.JSType = c.CreateNullableTypeOfJSType(recordType)
	var lhsType rhino.JSType = c.registry.CreateFunctionType(nullableRecordType)
	c.Assuming("a", lhsType)
	c.InFunction(`let x = 'x';
let res = a(x = 'some')?.b
a;
x;
`)
	c.Verify("a", lhsType)
	c.VerifyNative("x", rhino.JSTypeNative_STRING_TYPE)
	c.Verify("res", c.CreateUnionType(rhino.JSTypeNative_VOID_TYPE, rhino.JSTypeNative_NUMBER_TYPE))
}

// Java: TypeInferenceTest.testRegularGetProp_unconditionalChangeToOuterVariable_inNullableReceiver()
func (c *TypeInferenceTest) TestRegularGetProp_unconditionalChangeToOuterVariable_inNullableReceiver() {
	var funcType rhino.JSType = c.registry.CreateFunctionType(c.registry.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE))
	var recordType rhino.JSType = c.CreateRecordType("b", funcType)
	var nullableRecordType rhino.JSType = c.CreateNullableTypeOfJSType(recordType)
	var lhsType rhino.JSType = c.registry.CreateFunctionType(nullableRecordType)
	c.Assuming("a", lhsType)
	c.InFunction(`let x = 'x';
let res = a(x = 1).b
a;
x;
`)
	c.Verify("a", lhsType)
	c.VerifyNative("x", rhino.JSTypeNative_NUMBER_TYPE)
	c.Verify("res", funcType)
}

// Java: TypeInferenceTest.testOptChainGetProp_unconditionalChangeToOuterVariable_inNonNullableReceiver()
func (c *TypeInferenceTest) TestOptChainGetProp_unconditionalChangeToOuterVariable_inNonNullableReceiver() {
	var recordType rhino.JSType = c.CreateRecordType("b", c.registry.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE))
	var lhsType rhino.JSType = c.registry.CreateFunctionType(recordType)
	c.Assuming("a", lhsType)
	c.InFunction(`let x = 'x';
let res = a(x = 'some')?.b
a;
x;
`)
	c.Verify("a", lhsType)
	c.VerifyNative("x", rhino.JSTypeNative_STRING_TYPE)
	c.VerifyNative("res", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testRegularGetProp_changeToOuterVariable_inNonNullableReceiver_andOptNodes()
func (c *TypeInferenceTest) TestRegularGetProp_changeToOuterVariable_inNonNullableReceiver_andOptNodes() {
	var funcType rhino.JSType = c.registry.CreateFunctionType(c.registry.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE))
	var recordType rhino.JSType = c.CreateRecordType("b", funcType)
	var lhsType rhino.JSType = c.registry.CreateFunctionType(recordType)
	c.Assuming("a", lhsType)
	c.InFunction(`let x = 'x';
let res = a(x = 1).b(x = 'x')
a;
x;
`)
	c.Verify("a", lhsType)
	c.VerifyNative("x", rhino.JSTypeNative_STRING_TYPE)
	c.VerifyNative("res", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testOptChainGetProp_changeToOuterVariable_inNonNullableReceiver_andOptNodes()
func (c *TypeInferenceTest) TestOptChainGetProp_changeToOuterVariable_inNonNullableReceiver_andOptNodes() {
	var funcType rhino.JSType = c.registry.CreateFunctionType(c.registry.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE))
	var recordType rhino.JSType = c.CreateRecordType("b", funcType)
	var lhsType rhino.JSType = c.registry.CreateFunctionType(recordType)
	c.Assuming("a", lhsType)
	c.InFunction(`let x = 'x';
let res = a(x = 1)?.b(x = 'x')
a;
x; res;
`)
	c.Verify("a", lhsType)
	c.VerifyNative("x", rhino.JSTypeNative_STRING_TYPE)
	c.VerifyNative("res", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testRegularGetProp_changeToOuterVariable_inNullableReceiver_andOptNodes()
func (c *TypeInferenceTest) TestRegularGetProp_changeToOuterVariable_inNullableReceiver_andOptNodes() {
	var funcType rhino.JSType = c.registry.CreateFunctionType(c.registry.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE))
	var recordType rhino.JSType = c.CreateRecordType("b", funcType)
	var nullableRecordType rhino.JSType = c.CreateNullableTypeOfJSType(recordType)
	var lhsType rhino.JSType = c.registry.CreateFunctionType(nullableRecordType)
	c.Assuming("a", lhsType)
	c.InFunction(`let x = 'x';
let res = a(x = 1).b(x = 'x')
a;
x;
`)
	c.Verify("a", lhsType)
	c.Verify("x", c.registry.CreateUnionTypeOfNatives(rhino.JSTypeNative_STRING_TYPE))
	c.VerifyNative("res", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testOptChainGetProp_changeToOuterVariable_inNullableReceiver_andOptNodes()
func (c *TypeInferenceTest) TestOptChainGetProp_changeToOuterVariable_inNullableReceiver_andOptNodes() {
	var funcType rhino.JSType = c.registry.CreateFunctionType(c.registry.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE))
	var recordType rhino.JSType = c.CreateRecordType("b", funcType)
	var nullableRecordType rhino.JSType = c.CreateNullableTypeOfJSType(recordType)
	var lhsType rhino.JSType = c.registry.CreateFunctionType(nullableRecordType)
	c.Assuming("a", lhsType)
	c.InFunction(`let x = 'x';
let res = a(x = 1)?.b(x = 'x')
a;
x;
`)
	c.Verify("a", lhsType)
	c.Verify("x", c.registry.CreateUnionTypeOfNatives(rhino.JSTypeNative_NUMBER_TYPE, rhino.JSTypeNative_STRING_TYPE))
	c.Verify("res", c.CreateUnionType(rhino.JSTypeNative_VOID_TYPE, rhino.JSTypeNative_NUMBER_TYPE))
}

// Java: TypeInferenceTest.testOptChainGetProp_unconditionalChangeToOuterVariable()
func (c *TypeInferenceTest) TestOptChainGetProp_unconditionalChangeToOuterVariable() {
	var lhsType rhino.JSType = c.CreateRecordType("b", c.registry.CreateFunctionType(c.registry.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)))
	c.Assuming("a", lhsType)
	c.InFunction(`let x = 'x';
a?.b(x=5);
a;
x;
`)
	c.Verify("a", lhsType)
	c.VerifyNative("x", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testRegularGetProp_unconditionalChangeToOuterVariable()
func (c *TypeInferenceTest) TestRegularGetProp_unconditionalChangeToOuterVariable() {
	var lhsType rhino.JSType = c.CreateRecordType("b", c.registry.CreateFunctionType(c.registry.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)))
	c.Assuming("a", lhsType)
	c.InFunction(`let x = 'x';
a.b(x=5);
a;
x;
`)
	c.Verify("a", lhsType)
	c.VerifyNative("x", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testOptChainGetProp_invalidChangeToOuterVariable()
func (c *TypeInferenceTest) TestOptChainGetProp_invalidChangeToOuterVariable() {
	var lhsType rhino.JSType = c.CreateRecordType("b", c.registry.CreateFunctionType(c.registry.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)))
	c.Assuming("a", lhsType)
	c.InFunction("let x = 'x';\na?.c(x=5); // even though prop `c` is inexistent, x=5 will run, and typeOf(x) will\n// change.\na;\nx;\n")
	c.Verify("a", lhsType)
	c.VerifyNative("x", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testOptChainGetProp_multipleChains()
func (c *TypeInferenceTest) TestOptChainGetProp_multipleChains() {
	var lhsType rhino.JSType = c.CreateRecordType("b", c.registry.CreateFunctionType(c.registry.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)))
	c.Assuming("a", lhsType)
	c.InFunction("let res = a?.b(x?.y.z);")
	c.VerifyNative("res", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testOptChainGetProp_multipleChains_nullableReceiver()
func (c *TypeInferenceTest) TestOptChainGetProp_multipleChains_nullableReceiver() {
	var lhsType rhino.JSType = c.CreateNullableTypeOfJSType(c.CreateRecordType("b", c.registry.CreateFunctionType(c.registry.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE))))
	c.Assuming("a", lhsType)
	c.InFunction("let res = a?.b(x?.y.z);")
	c.Verify("res", c.registry.CreateUnionTypeOfNatives(rhino.JSTypeNative_VOID_TYPE, rhino.JSTypeNative_NUMBER_TYPE))
}

// Java: TypeInferenceTest.testGetElemDereference()
func (c *TypeInferenceTest) TestGetElemDereference() {
	c.Assuming("x", c.CreateUndefinableType(rhino.JSTypeNative_OBJECT_TYPE))
	c.InFunction("x['z'] = 3;")
	c.VerifyNative("x", rhino.JSTypeNative_OBJECT_TYPE)
}

// Java: TypeInferenceTest.testGetElemDereference_knownSymbol()
func (c *TypeInferenceTest) TestGetElemDereference_knownSymbol() {
	var array rhino.ObjectType = c.GetNativeObjectType(rhino.JSTypeNative_ARRAY_TYPE)
	var sym *rhino.KnownSymbolType = rhino.NewKnownSymbolType(c.registry, "sym")
	array.DefineDeclaredPropertyByKey(rhino.NewProperty_SymbolKey(sym), c.GetNativeType(rhino.JSTypeNative_STRING_TYPE), rhino.NilNode)
	c.Assuming("o", array)
	c.Assuming("sym", sym)
	c.InFunction("const x = o[sym];")
	c.VerifyNative("x", rhino.JSTypeNative_STRING_TYPE)
}

// Java: TypeInferenceTest.testGetElemDereference_knownSymbol_nullable()
func (c *TypeInferenceTest) TestGetElemDereference_knownSymbol_nullable() {
	var array rhino.ObjectType = c.GetNativeObjectType(rhino.JSTypeNative_ARRAY_TYPE)
	var sym *rhino.KnownSymbolType = rhino.NewKnownSymbolType(c.registry, "sym")
	array.DefineDeclaredPropertyByKey(rhino.NewProperty_SymbolKey(sym), c.GetNativeType(rhino.JSTypeNative_STRING_TYPE), rhino.NilNode)
	c.Assuming("o", c.CreateNullableTypeOfJSType(array))
	c.Assuming("sym", sym)
	c.InFunction("const x = o[sym];")
	c.VerifyNative("x", rhino.JSTypeNative_STRING_TYPE)
}

// Java: TypeInferenceTest.testGetElemDereference_knownSymbol_union()
func (c *TypeInferenceTest) TestGetElemDereference_knownSymbol_union() {
	var array rhino.ObjectType = c.GetNativeObjectType(rhino.JSTypeNative_ARRAY_TYPE)
	var promise rhino.ObjectType = c.GetNativeObjectType(rhino.JSTypeNative_PROMISE_TYPE)
	var sym *rhino.KnownSymbolType = rhino.NewKnownSymbolType(c.registry, "sym")
	array.DefineDeclaredPropertyByKey(rhino.NewProperty_SymbolKey(sym), c.GetNativeType(rhino.JSTypeNative_STRING_TYPE), rhino.NilNode)
	promise.DefineDeclaredPropertyByKey(rhino.NewProperty_SymbolKey(sym), c.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE), rhino.NilNode)
	c.Assuming("o", c.registry.CreateUnionType(array, promise))
	c.Assuming("sym", sym)
	c.InFunction("const x = o[sym];")
	c.Verify("x", c.registry.CreateUnionType(c.GetNativeType(rhino.JSTypeNative_STRING_TYPE), c.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)))
}

// Java: TypeInferenceTest.testGetElemDereference_knownSymbol_union_notOnAllAlternates()
func (c *TypeInferenceTest) TestGetElemDereference_knownSymbol_union_notOnAllAlternates() {
	var array rhino.ObjectType = c.GetNativeObjectType(rhino.JSTypeNative_ARRAY_TYPE)
	var promise rhino.ObjectType = c.GetNativeObjectType(rhino.JSTypeNative_PROMISE_TYPE)
	var sym *rhino.KnownSymbolType = rhino.NewKnownSymbolType(c.registry, "sym")
	array.DefineDeclaredPropertyByKey(rhino.NewProperty_SymbolKey(sym), c.GetNativeType(rhino.JSTypeNative_STRING_TYPE), rhino.NilNode)
	c.Assuming("o", c.registry.CreateUnionType(array, promise))
	c.Assuming("sym", sym)
	c.InFunction("const x = o[sym];")
	c.VerifyNative("x", rhino.JSTypeNative_UNKNOWN_TYPE)
}

// Java: TypeInferenceTest.testGetElemDereference_knownSymbol_onAllType()
func (c *TypeInferenceTest) TestGetElemDereference_knownSymbol_onAllType() {
	var sym *rhino.KnownSymbolType = rhino.NewKnownSymbolType(c.registry, "sym")
	c.AssumingNative("o", rhino.JSTypeNative_ALL_TYPE)
	c.Assuming("sym", sym)
	c.InFunction("const x = o[sym];")
	c.VerifyNative("x", rhino.JSTypeNative_UNKNOWN_TYPE)
}

// Java: TypeInferenceTest.testGetElemDereference_knownSymbol_onUnknownType()
func (c *TypeInferenceTest) TestGetElemDereference_knownSymbol_onUnknownType() {
	var sym *rhino.KnownSymbolType = rhino.NewKnownSymbolType(c.registry, "sym")
	c.AssumingNative("o", rhino.JSTypeNative_UNKNOWN_TYPE)
	c.Assuming("sym", sym)
	c.InFunction("const x = o[sym];")
	c.VerifyNative("x", rhino.JSTypeNative_UNKNOWN_TYPE)
}

// Java: TypeInferenceTest.testIf1()
func (c *TypeInferenceTest) TestIf1() {
	c.Assuming("x", c.CreateNullableType(rhino.JSTypeNative_OBJECT_TYPE))
	c.InFunction("var y = {}; if (x) { y = x; }")
	c.VerifySubtypeOfNative("y", rhino.JSTypeNative_OBJECT_TYPE)
}

// Java: TypeInferenceTest.testIf1a()
func (c *TypeInferenceTest) TestIf1a() {
	c.Assuming("x", c.CreateNullableType(rhino.JSTypeNative_OBJECT_TYPE))
	c.InFunction("var y = {}; if (x != null) { y = x; }")
	c.VerifySubtypeOfNative("y", rhino.JSTypeNative_OBJECT_TYPE)
}

// Java: TypeInferenceTest.testNullishCoalesceNullableObject()
func (c *TypeInferenceTest) TestNullishCoalesceNullableObject() {
	c.Assuming("x", c.CreateNullableType(rhino.JSTypeNative_OBJECT_TYPE))
	c.InFunction("let z = x ?? {}")
	c.VerifyNative("z", rhino.JSTypeNative_OBJECT_TYPE)
}

// Java: TypeInferenceTest.testNullishCoalesceNullableUnion()
func (c *TypeInferenceTest) TestNullishCoalesceNullableUnion() {
	c.Assuming("x", c.CreateNullableTypeOfJSType(c.CreateUnionType(rhino.JSTypeNative_OBJECT_TYPE, rhino.JSTypeNative_STRING_TYPE)))
	c.InFunction("let z = x ?? {}")
	c.Verify("z", c.CreateUnionType(rhino.JSTypeNative_STRING_TYPE, rhino.JSTypeNative_OBJECT_TYPE))
}

// Java: TypeInferenceTest.testIf2()
func (c *TypeInferenceTest) TestIf2() {
	c.Assuming("x", c.CreateNullableType(rhino.JSTypeNative_OBJECT_TYPE))
	c.InFunction("var y = x; if (x) { y = x; } else { y = {}; }")
	c.VerifySubtypeOfNative("y", rhino.JSTypeNative_OBJECT_TYPE)
}

// Java: TypeInferenceTest.testIf3()
func (c *TypeInferenceTest) TestIf3() {
	c.Assuming("x", c.CreateNullableType(rhino.JSTypeNative_OBJECT_TYPE))
	c.InFunction("var y = 1; if (x) { y = x; }")
	c.Verify("y", c.CreateUnionType(rhino.JSTypeNative_OBJECT_TYPE, rhino.JSTypeNative_NUMBER_TYPE))
}

// Java: TypeInferenceTest.testPropertyInference1()
func (c *TypeInferenceTest) TestPropertyInference1() {
	var thisType rhino.ObjectType = c.registry.CreateAnonymousObjectType(nil)
	thisType.DefineDeclaredProperty("foo", c.CreateUndefinableType(rhino.JSTypeNative_STRING_TYPE), rhino.NilNode)
	c.AssumingThisType(thisType)
	c.InFunction("var y = 1; if (this.foo) { y = this.foo; }")
	c.Verify("y", c.CreateUnionType(rhino.JSTypeNative_NUMBER_TYPE, rhino.JSTypeNative_STRING_TYPE))
}

// Java: TypeInferenceTest.testPropertyInference2()
func (c *TypeInferenceTest) TestPropertyInference2() {
	var thisType rhino.ObjectType = c.registry.CreateAnonymousObjectType(nil)
	thisType.DefineDeclaredProperty("foo", c.CreateUndefinableType(rhino.JSTypeNative_STRING_TYPE), rhino.NilNode)
	c.AssumingThisType(thisType)
	c.InFunction("var y = 1; this.foo = 'x'; y = this.foo;")
	c.VerifyNative("y", rhino.JSTypeNative_STRING_TYPE)
}

// Java: TypeInferenceTest.testPropertyInference3()
func (c *TypeInferenceTest) TestPropertyInference3() {
	var thisType rhino.ObjectType = c.registry.CreateAnonymousObjectType(nil)
	thisType.DefineDeclaredProperty("foo", c.CreateUndefinableType(rhino.JSTypeNative_STRING_TYPE), rhino.NilNode)
	c.AssumingThisType(thisType)
	c.InFunction("var y = 1; this.foo = x; y = this.foo;")
	c.Verify("y", c.CreateUndefinableType(rhino.JSTypeNative_STRING_TYPE))
}

// Java: TypeInferenceTest.testAssert1()
func (c *TypeInferenceTest) TestAssert1() {
	var startType rhino.JSType = c.CreateNullableType(rhino.JSTypeNative_OBJECT_TYPE)
	c.Assuming("x", startType)
	c.InFunction("out1 = x; goog.asserts.assert(x); out2 = x;")
	c.Verify("out1", startType)
	c.VerifyNative("out2", rhino.JSTypeNative_OBJECT_TYPE)
}

// Java: TypeInferenceTest.testAssert1a()
func (c *TypeInferenceTest) TestAssert1a() {
	var startType rhino.JSType = c.CreateNullableType(rhino.JSTypeNative_OBJECT_TYPE)
	c.Assuming("x", startType)
	c.InFunction("out1 = x; goog.asserts.assert(x !== null); out2 = x;")
	c.Verify("out1", startType)
	c.VerifyNative("out2", rhino.JSTypeNative_OBJECT_TYPE)
}

// Java: TypeInferenceTest.testAssert2()
func (c *TypeInferenceTest) TestAssert2() {
	var startType rhino.JSType = c.CreateNullableType(rhino.JSTypeNative_OBJECT_TYPE)
	c.Assuming("x", startType)
	c.InFunction("goog.asserts.assert(1, x); out1 = x;")
	c.Verify("out1", startType)
}

// Java: TypeInferenceTest.testAssert3()
func (c *TypeInferenceTest) TestAssert3() {
	var startType rhino.JSType = c.CreateNullableType(rhino.JSTypeNative_OBJECT_TYPE)
	c.Assuming("x", startType)
	c.Assuming("y", startType)
	c.InFunction("out1 = x; goog.asserts.assert(x && y); out2 = x; out3 = y;")
	c.Verify("out1", startType)
	c.VerifyNative("out2", rhino.JSTypeNative_OBJECT_TYPE)
	c.VerifyNative("out3", rhino.JSTypeNative_OBJECT_TYPE)
}

// Java: TypeInferenceTest.testAssert4()
func (c *TypeInferenceTest) TestAssert4() {
	var startType rhino.JSType = c.CreateNullableType(rhino.JSTypeNative_OBJECT_TYPE)
	c.Assuming("x", startType)
	c.Assuming("y", startType)
	c.InFunction("out1 = x; goog.asserts.assert(x && !y); out2 = x; out3 = y;")
	c.Verify("out1", startType)
	c.VerifyNative("out2", rhino.JSTypeNative_OBJECT_TYPE)
	c.VerifyNative("out3", rhino.JSTypeNative_NULL_TYPE)
}

// Java: TypeInferenceTest.testAssert5()
func (c *TypeInferenceTest) TestAssert5() {
	var startType rhino.JSType = c.CreateNullableType(rhino.JSTypeNative_OBJECT_TYPE)
	c.Assuming("x", startType)
	c.Assuming("y", startType)
	c.InFunction("goog.asserts.assert(x || y); out1 = x; out2 = y;")
	c.Verify("out1", startType)
	c.Verify("out2", startType)
}

// Java: TypeInferenceTest.testAssert5NullishCoalesce()
func (c *TypeInferenceTest) TestAssert5NullishCoalesce() {
	var startType rhino.JSType = c.CreateNullableType(rhino.JSTypeNative_OBJECT_TYPE)
	c.Assuming("x", startType)
	c.Assuming("y", startType)
	c.InFunction("goog.asserts.assert(x ?? y); out1 = x; out2 = y;")
	c.Verify("out1", startType)
	c.Verify("out2", startType)
}

// Java: TypeInferenceTest.testAssert6()
func (c *TypeInferenceTest) TestAssert6() {
	var startType rhino.JSType = c.CreateNullableType(rhino.JSTypeNative_OBJECT_TYPE)
	c.Assuming("x", c.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE))
	c.Assuming("x.y", startType)
	c.InFunction("out1 = x.y; goog.asserts.assert(x.y); out2 = x.y;")
	c.Verify("out1", startType)
	c.VerifyNative("out2", rhino.JSTypeNative_OBJECT_TYPE)
}

// Java: TypeInferenceTest.testAssert7()
func (c *TypeInferenceTest) TestAssert7() {
	var startType rhino.JSType = c.CreateNullableType(rhino.JSTypeNative_OBJECT_TYPE)
	c.Assuming("x", startType)
	c.InFunction("out1 = x; out2 = goog.asserts.assert(x);")
	c.Verify("out1", startType)
	c.VerifyNative("out2", rhino.JSTypeNative_OBJECT_TYPE)
}

// Java: TypeInferenceTest.testAssert8()
func (c *TypeInferenceTest) TestAssert8() {
	var startType rhino.JSType = c.CreateNullableType(rhino.JSTypeNative_OBJECT_TYPE)
	c.Assuming("x", startType)
	c.InFunction("out1 = x; out2 = goog.asserts.assert(x != null);")
	c.Verify("out1", startType)
	c.VerifyNative("out2", rhino.JSTypeNative_BOOLEAN_TYPE)
}

// Java: TypeInferenceTest.testAssert9()
func (c *TypeInferenceTest) TestAssert9() {
	var startType rhino.JSType = c.CreateNullableType(rhino.JSTypeNative_NUMBER_TYPE)
	c.Assuming("x", startType)
	c.InFunction("out1 = x; out2 = goog.asserts.assert(y = x);")
	c.Verify("out1", startType)
	c.VerifyNative("out2", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testAssert11()
func (c *TypeInferenceTest) TestAssert11() {
	var startType rhino.JSType = c.CreateNullableType(rhino.JSTypeNative_OBJECT_TYPE)
	c.Assuming("x", startType)
	c.Assuming("y", startType)
	c.InFunction("var z = goog.asserts.assert(x || y);")
	c.Verify("x", startType)
	c.Verify("y", startType)
}

// Java: TypeInferenceTest.testPrimitiveAssertTruthy_narrowsNullableObjectToObject()
func (c *TypeInferenceTest) TestPrimitiveAssertTruthy_narrowsNullableObjectToObject() {
	var startType rhino.JSType = c.CreateNullableType(rhino.JSTypeNative_OBJECT_TYPE)
	c.IncludePrimitiveTruthyAssertionFunction("assertTruthy")
	c.Assuming("x", startType)
	c.InFunction("out1 = x; assertTruthy(x); out2 = x;")
	c.Verify("out1", startType)
	c.VerifyNative("out2", rhino.JSTypeNative_OBJECT_TYPE)
}

// Java: TypeInferenceTest.testPrimitiveAssertTruthy_narrowsNullableObjectInNeqNullToObject()
func (c *TypeInferenceTest) TestPrimitiveAssertTruthy_narrowsNullableObjectInNeqNullToObject() {
	var startType rhino.JSType = c.CreateNullableType(rhino.JSTypeNative_OBJECT_TYPE)
	c.IncludePrimitiveTruthyAssertionFunction("assertTruthy")
	c.Assuming("x", startType)
	c.InFunction("out1 = x; assertTruthy(x !== null); out2 = x;")
	c.Verify("out1", startType)
	c.VerifyNative("out2", rhino.JSTypeNative_OBJECT_TYPE)
}

// Java: TypeInferenceTest.testPrimitiveAssertTruthy_ignoresSecondArgumentEvenIfNullable()
func (c *TypeInferenceTest) TestPrimitiveAssertTruthy_ignoresSecondArgumentEvenIfNullable() {
	var startType rhino.JSType = c.CreateNullableType(rhino.JSTypeNative_OBJECT_TYPE)
	c.IncludePrimitiveTruthyAssertionFunction("assertTruthy")
	c.Assuming("x", startType)
	c.InFunction("assertTruthy(1, x); out1 = x;")
	c.Verify("out1", startType)
}

// Java: TypeInferenceTest.testAssertNumber_narrowsAllTypeToNumber()
func (c *TypeInferenceTest) TestAssertNumber_narrowsAllTypeToNumber() {
	var startType rhino.JSType = c.CreateNullableType(rhino.JSTypeNative_ALL_TYPE)
	c.IncludeGoogAssertionFn("assertNumber", c.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE))
	c.Assuming("x", startType)
	c.InFunction("out1 = x; goog.asserts.assertNumber(x); out2 = x;")
	c.Verify("out1", startType)
	c.VerifyNative("out2", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testAssertNumber_doesNotNarrowNamesInExpression()
func (c *TypeInferenceTest) TestAssertNumber_doesNotNarrowNamesInExpression() {
	var startType rhino.JSType = c.CreateNullableType(rhino.JSTypeNative_ALL_TYPE)
	c.IncludeGoogAssertionFn("assertNumber", c.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE))
	c.Assuming("x", startType)
	c.InFunction("goog.asserts.assertNumber(x + x); out1 = x;")
	c.Verify("out1", startType)
}

// Java: TypeInferenceTest.testAssertNumber_returnsNumberGivenExpression()
func (c *TypeInferenceTest) TestAssertNumber_returnsNumberGivenExpression() {
	var startType rhino.JSType = c.CreateNullableType(rhino.JSTypeNative_ALL_TYPE)
	c.IncludeGoogAssertionFn("assertNumber", c.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE))
	c.Assuming("x", startType)
	c.InFunction("out1 = x; out2 = goog.asserts.assertNumber(x + x);")
	c.Verify("out1", startType)
	c.VerifyNative("out2", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testPrimitiveAssertNumber_narrowsAllTypeToNumber()
func (c *TypeInferenceTest) TestPrimitiveAssertNumber_narrowsAllTypeToNumber() {
	var startType rhino.JSType = c.CreateNullableType(rhino.JSTypeNative_ALL_TYPE)
	c.IncludePrimitiveAssertionFn("assertNumber", c.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE))
	c.Assuming("x", startType)
	c.InFunction("out1 = x; assertNumber(x); out2 = x;")
	c.Verify("out1", startType)
	c.VerifyNative("out2", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testPrimitiveAssertNumber_doesNotNarrowNamesInExpression()
func (c *TypeInferenceTest) TestPrimitiveAssertNumber_doesNotNarrowNamesInExpression() {
	var startType rhino.JSType = c.CreateNullableType(rhino.JSTypeNative_ALL_TYPE)
	c.IncludePrimitiveAssertionFn("assertNumber", c.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE))
	c.Assuming("x", startType)
	c.InFunction("assertNumber(x + x); out1 = x;")
	c.Verify("out1", startType)
}

// Java: TypeInferenceTest.testPrimitiveAssertNumber_returnsNumberGivenExpression()
func (c *TypeInferenceTest) TestPrimitiveAssertNumber_returnsNumberGivenExpression() {
	var startType rhino.JSType = c.CreateNullableType(rhino.JSTypeNative_ALL_TYPE)
	c.IncludePrimitiveAssertionFn("assertNumber", c.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE))
	c.Assuming("x", startType)
	c.InFunction("out1 = x; out2 = assertNumber(x + x);")
	c.Verify("out1", startType)
	c.VerifyNative("out2", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testBigIntTypeAssignment()
func (c *TypeInferenceTest) TestBigIntTypeAssignment() {
	c.AssumingNative("x", rhino.JSTypeNative_UNKNOWN_TYPE)
	c.AssumingNative("y", rhino.JSTypeNative_UNKNOWN_TYPE)
	c.InFunction("x = 1n; y = BigInt(1);")
	c.VerifyNative("x", rhino.JSTypeNative_BIGINT_TYPE)
	c.VerifyNative("y", rhino.JSTypeNative_BIGINT_TYPE)
}

// Java: TypeInferenceTest.testAssertBigInt_narrowsAllTypeToBigInt()
func (c *TypeInferenceTest) TestAssertBigInt_narrowsAllTypeToBigInt() {
	var startType rhino.JSType = c.CreateNullableType(rhino.JSTypeNative_ALL_TYPE)
	c.IncludePrimitiveAssertionFn("assertBigInt", c.GetNativeType(rhino.JSTypeNative_BIGINT_TYPE))
	c.Assuming("x", startType)
	c.InFunction("out1 = x; assertBigInt(x); out2 = x;")
	c.Verify("out1", startType)
	c.VerifyNative("out2", rhino.JSTypeNative_BIGINT_TYPE)
}

// Java: TypeInferenceTest.testBigIntPresence()
func (c *TypeInferenceTest) TestBigIntPresence() {
	TypeInferenceTest_TestForAllBigInt(c.GetNativeType(rhino.JSTypeNative_BIGINT_TYPE))
	TypeInferenceTest_TestForAllBigInt(c.GetNativeType(rhino.JSTypeNative_BIGINT_OBJECT_TYPE))
	TypeInferenceTest_TestForNoBigInt(c.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE))
	TypeInferenceTest_TestForNoBigInt(c.GetNativeType(rhino.JSTypeNative_STRING_TYPE))
	TypeInferenceTest_TestForNoBigInt(c.GetNativeType(rhino.JSTypeNative_ALL_TYPE))
	TypeInferenceTest_TestForNoBigInt(c.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE))
	TypeInferenceTest_TestForNoBigInt(c.GetNativeType(rhino.JSTypeNative_NO_TYPE))
	TypeInferenceTest_TestForAllBigInt(c.CreateUnionType(rhino.JSTypeNative_BIGINT_TYPE, rhino.JSTypeNative_BIGINT_OBJECT_TYPE))
	TypeInferenceTest_TestForBigIntOrNumber(c.GetNativeType(rhino.JSTypeNative_BIGINT_NUMBER))
	TypeInferenceTest_TestForBigIntOrOther(c.CreateUnionType(rhino.JSTypeNative_BIGINT_TYPE, rhino.JSTypeNative_STRING_TYPE))
	TypeInferenceTest_TestForNoBigInt(c.GetNativeType(rhino.JSTypeNative_NUMBER_STRING))
	TypeInferenceTest_TestForBigIntOrNumber(c.CreateUnionType(rhino.JSTypeNative_NUMBER_OBJECT_TYPE, rhino.JSTypeNative_BIGINT_NUMBER))
	TypeInferenceTest_TestForBigIntOrNumber(c.CreateUnionType(rhino.JSTypeNative_BIGINT_OBJECT_TYPE, rhino.JSTypeNative_BIGINT_NUMBER))
	TypeInferenceTest_TestForBigIntOrNumber(c.registry.CreateUnionType(c.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE), c.CreateUnionType(rhino.JSTypeNative_BIGINT_TYPE, rhino.JSTypeNative_BIGINT_OBJECT_TYPE)))
	TypeInferenceTest_TestForBigIntOrNumber(c.registry.CreateUnionType(c.GetNativeType(rhino.JSTypeNative_BIGINT_TYPE), c.CreateUnionType(rhino.JSTypeNative_NUMBER_TYPE, rhino.JSTypeNative_NUMBER_OBJECT_TYPE)))
	TypeInferenceTest_TestForBigIntOrOther(c.registry.CreateUnionType(c.GetNativeType(rhino.JSTypeNative_BIGINT_TYPE), c.CreateUnionType(rhino.JSTypeNative_STRING_TYPE, rhino.JSTypeNative_STRING_OBJECT_TYPE)))
	TypeInferenceTest_TestForAllBigInt(c.registry.CreateUnionType(c.GetNativeType(rhino.JSTypeNative_BIGINT_OBJECT_TYPE), c.CreateEnumType("Enum", rhino.JSTypeNative_BIGINT_TYPE).GetElementsType()))
	TypeInferenceTest_TestForBigIntOrNumber(c.registry.CreateUnionType(c.GetNativeType(rhino.JSTypeNative_BIGINT_TYPE), c.CreateEnumType("Enum", rhino.JSTypeNative_NUMBER_TYPE).GetElementsType()))
	TypeInferenceTest_TestForBigIntOrOther(c.registry.CreateUnionType(c.GetNativeType(rhino.JSTypeNative_BIGINT_TYPE), c.CreateEnumType("Enum", rhino.JSTypeNative_STRING_TYPE).GetElementsType()))
	TypeInferenceTest_TestForAllBigInt(c.CreateEnumType("Enum", rhino.JSTypeNative_BIGINT_TYPE).GetElementsType())
	TypeInferenceTest_TestForNoBigInt(c.CreateEnumType("Enum", rhino.JSTypeNative_NUMBER_TYPE).GetElementsType())
	TypeInferenceTest_TestForAllBigInt(c.CreateEnumTypeOfJSType("Enum", c.CreateEnumType("Enum", rhino.JSTypeNative_BIGINT_TYPE).GetElementsType()).GetElementsType())
	TypeInferenceTest_TestForNoBigInt(c.CreateEnumTypeOfJSType("Enum", c.CreateEnumType("Enum", rhino.JSTypeNative_NUMBER_TYPE).GetElementsType()).GetElementsType())
	TypeInferenceTest_TestForAllBigInt(c.CreateEnumTypeOfJSType("Enum", c.CreateUnionType(rhino.JSTypeNative_BIGINT_TYPE, rhino.JSTypeNative_BIGINT_OBJECT_TYPE)).GetElementsType())
	TypeInferenceTest_TestForBigIntOrNumber(c.CreateEnumType("Enum", rhino.JSTypeNative_BIGINT_NUMBER).GetElementsType())
	TypeInferenceTest_TestForBigIntOrOther(c.CreateEnumTypeOfJSType("Enum", c.CreateUnionType(rhino.JSTypeNative_BIGINT_TYPE, rhino.JSTypeNative_STRING_TYPE)).GetElementsType())
}

// Java: TypeInferenceTest.testBigIntWithUnaryPlus()
func (c *TypeInferenceTest) TestBigIntWithUnaryPlus() {
	c.AssumingNative("x", rhino.JSTypeNative_BIGINT_TYPE)
	c.AssumingNative("y", rhino.JSTypeNative_BIGINT_OBJECT_TYPE)
	c.AssumingNative("z", rhino.JSTypeNative_BIGINT_NUMBER)
	c.InFunction("valueType = +x; objectType = +y; unionType = +z;")
	c.VerifyNative("valueType", rhino.JSTypeNative_NO_TYPE)
	c.VerifyNative("objectType", rhino.JSTypeNative_NO_TYPE)
	c.VerifyNative("unionType", rhino.JSTypeNative_NO_TYPE)
}

// Java: TypeInferenceTest.testBigIntEnumWithUnaryPlus()
func (c *TypeInferenceTest) TestBigIntEnumWithUnaryPlus() {
	var enumElementBigIntType *rhino.EnumElementType = c.CreateEnumType("MyEnum", rhino.JSTypeNative_BIGINT_TYPE).GetElementsType()
	var enumElementUnionType *rhino.EnumElementType = c.CreateEnumType("MyEnum", rhino.JSTypeNative_BIGINT_NUMBER).GetElementsType()
	c.Assuming("x", enumElementBigIntType)
	c.Assuming("y", c.registry.CreateUnionType(enumElementBigIntType, c.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)))
	c.Assuming("z", enumElementUnionType)
	c.InFunction("enumElementBigIntType = +x; unionEnumType = +y; enumElementUnionType = +z;")
	c.VerifyNative("enumElementBigIntType", rhino.JSTypeNative_NO_TYPE)
	c.VerifyNative("unionEnumType", rhino.JSTypeNative_NO_TYPE)
	c.VerifyNative("enumElementUnionType", rhino.JSTypeNative_NO_TYPE)
}

// Java: TypeInferenceTest.testBigIntWithLogicalNOT()
func (c *TypeInferenceTest) TestBigIntWithLogicalNOT() {
	c.AssumingNative("x", rhino.JSTypeNative_BIGINT_TYPE)
	c.AssumingNative("y", rhino.JSTypeNative_BIGINT_OBJECT_TYPE)
	c.AssumingNative("z", rhino.JSTypeNative_BIGINT_NUMBER)
	c.InFunction("valueType = !x; objectType = !y; unionType = !z;")
	c.VerifyNative("valueType", rhino.JSTypeNative_BOOLEAN_TYPE)
	c.VerifyNative("objectType", rhino.JSTypeNative_BOOLEAN_TYPE)
	c.VerifyNative("unionType", rhino.JSTypeNative_BOOLEAN_TYPE)
}

// Java: TypeInferenceTest.testBigIntWithTypeOfOperation()
func (c *TypeInferenceTest) TestBigIntWithTypeOfOperation() {
	c.AssumingNative("x", rhino.JSTypeNative_BIGINT_TYPE)
	c.AssumingNative("y", rhino.JSTypeNative_BIGINT_OBJECT_TYPE)
	c.Assuming("z", c.CreateUnionType(rhino.JSTypeNative_BIGINT_TYPE, rhino.JSTypeNative_BIGINT_OBJECT_TYPE))
	c.InFunction("valueType = typeof x; objectType = typeof y; unionType = typeof z;")
	c.VerifyNative("valueType", rhino.JSTypeNative_STRING_TYPE)
	c.VerifyNative("objectType", rhino.JSTypeNative_STRING_TYPE)
	c.VerifyNative("unionType", rhino.JSTypeNative_STRING_TYPE)
}

// Java: TypeInferenceTest.testBigIntWithDeleteOperation()
func (c *TypeInferenceTest) TestBigIntWithDeleteOperation() {
	c.AssumingNative("x", rhino.JSTypeNative_BIGINT_TYPE)
	c.AssumingNative("y", rhino.JSTypeNative_BIGINT_OBJECT_TYPE)
	c.Assuming("z", c.CreateUnionType(rhino.JSTypeNative_BIGINT_TYPE, rhino.JSTypeNative_BIGINT_OBJECT_TYPE))
	c.InFunction("valueType = delete x; objectType = delete y; unionType = delete z;")
	c.VerifyNative("valueType", rhino.JSTypeNative_BOOLEAN_TYPE)
	c.VerifyNative("objectType", rhino.JSTypeNative_BOOLEAN_TYPE)
	c.VerifyNative("unionType", rhino.JSTypeNative_BOOLEAN_TYPE)
}

// Java: TypeInferenceTest.testBigIntWithVoidOperation()
func (c *TypeInferenceTest) TestBigIntWithVoidOperation() {
	c.AssumingNative("x", rhino.JSTypeNative_BIGINT_TYPE)
	c.AssumingNative("y", rhino.JSTypeNative_BIGINT_OBJECT_TYPE)
	c.Assuming("z", c.CreateUnionType(rhino.JSTypeNative_BIGINT_TYPE, rhino.JSTypeNative_BIGINT_OBJECT_TYPE))
	c.InFunction("valueType = void x; objectType = void y; unionType = void z;")
	c.VerifyNative("valueType", rhino.JSTypeNative_VOID_TYPE)
	c.VerifyNative("objectType", rhino.JSTypeNative_VOID_TYPE)
	c.VerifyNative("unionType", rhino.JSTypeNative_VOID_TYPE)
}

// Java: TypeInferenceTest.testBigIntWithUnaryMinus()
func (c *TypeInferenceTest) TestBigIntWithUnaryMinus() {
	c.AssumingNative("x", rhino.JSTypeNative_BIGINT_TYPE)
	c.AssumingNative("y", rhino.JSTypeNative_BIGINT_OBJECT_TYPE)
	c.AssumingNative("u", rhino.JSTypeNative_UNKNOWN_TYPE)
	c.AssumingNative("a", rhino.JSTypeNative_ALL_TYPE)
	c.AssumingNative("z1", rhino.JSTypeNative_BIGINT_NUMBER)
	c.Assuming("z2", c.CreateUnionType(rhino.JSTypeNative_BIGINT_TYPE, rhino.JSTypeNative_STRING_TYPE))
	c.InFunction(`valueType = -x; objectType = -y; unknownType = -u; allType = -a; bigintNum = -z1;
 bigintOther = -z2;
`)
	c.VerifyNative("valueType", rhino.JSTypeNative_BIGINT_TYPE)
	c.VerifyNative("objectType", rhino.JSTypeNative_BIGINT_TYPE)
	c.VerifyNative("unknownType", rhino.JSTypeNative_NUMBER_TYPE)
	c.VerifyNative("allType", rhino.JSTypeNative_NUMBER_TYPE)
	c.VerifyNative("bigintNum", rhino.JSTypeNative_BIGINT_NUMBER)
	c.VerifyNative("bigintOther", rhino.JSTypeNative_BIGINT_NUMBER)
}

// Java: TypeInferenceTest.testBigIntWithBitwiseNOT()
func (c *TypeInferenceTest) TestBigIntWithBitwiseNOT() {
	c.AssumingNative("x", rhino.JSTypeNative_BIGINT_TYPE)
	c.AssumingNative("y", rhino.JSTypeNative_BIGINT_OBJECT_TYPE)
	c.AssumingNative("u", rhino.JSTypeNative_UNKNOWN_TYPE)
	c.AssumingNative("a", rhino.JSTypeNative_ALL_TYPE)
	c.AssumingNative("z1", rhino.JSTypeNative_BIGINT_NUMBER)
	c.Assuming("z2", c.CreateUnionType(rhino.JSTypeNative_BIGINT_TYPE, rhino.JSTypeNative_STRING_TYPE))
	c.InFunction(`valueType = ~x; objectType = ~y; unknownType = ~u; allType = ~a; bigintOrNumber = ~z1;
 bigintOrOther = ~z2;
`)
	c.VerifyNative("valueType", rhino.JSTypeNative_BIGINT_TYPE)
	c.VerifyNative("objectType", rhino.JSTypeNative_BIGINT_TYPE)
	c.VerifyNative("unknownType", rhino.JSTypeNative_NUMBER_TYPE)
	c.VerifyNative("allType", rhino.JSTypeNative_NUMBER_TYPE)
	c.VerifyNative("bigintOrNumber", rhino.JSTypeNative_BIGINT_NUMBER)
	c.VerifyNative("bigintOrOther", rhino.JSTypeNative_BIGINT_NUMBER)
}

// Java: TypeInferenceTest.testIncrementOnBigInt()
func (c *TypeInferenceTest) TestIncrementOnBigInt() {
	c.AssumingNative("x", rhino.JSTypeNative_BIGINT_TYPE)
	c.AssumingNative("y", rhino.JSTypeNative_BIGINT_OBJECT_TYPE)
	c.AssumingNative("z", rhino.JSTypeNative_BIGINT_NUMBER)
	c.Assuming("a", c.CreateUnionType(rhino.JSTypeNative_BIGINT_TYPE, rhino.JSTypeNative_STRING_TYPE))
	c.InFunction("valueType = x++; objectType = y++; bigintNumber = z++; bigintOther = a++;")
	c.VerifyNative("valueType", rhino.JSTypeNative_BIGINT_TYPE)
	c.VerifyNative("objectType", rhino.JSTypeNative_BIGINT_TYPE)
	c.VerifyNative("bigintNumber", rhino.JSTypeNative_BIGINT_NUMBER)
	c.VerifyNative("bigintOther", rhino.JSTypeNative_BIGINT_NUMBER)
}

// Java: TypeInferenceTest.testDecrementOnBigInt()
func (c *TypeInferenceTest) TestDecrementOnBigInt() {
	c.AssumingNative("x", rhino.JSTypeNative_BIGINT_TYPE)
	c.AssumingNative("y", rhino.JSTypeNative_BIGINT_OBJECT_TYPE)
	c.AssumingNative("z", rhino.JSTypeNative_BIGINT_NUMBER)
	c.Assuming("a", c.CreateUnionType(rhino.JSTypeNative_BIGINT_TYPE, rhino.JSTypeNative_STRING_TYPE))
	c.InFunction("valueType = x--; objectType = y--; bigintNumber = z++; bigintOther = a++;")
	c.VerifyNative("valueType", rhino.JSTypeNative_BIGINT_TYPE)
	c.VerifyNative("objectType", rhino.JSTypeNative_BIGINT_TYPE)
	c.VerifyNative("bigintNumber", rhino.JSTypeNative_BIGINT_NUMBER)
	c.VerifyNative("bigintOther", rhino.JSTypeNative_BIGINT_NUMBER)
}

// Java: TypeInferenceTest.testAdditionWithBigInt()
func (c *TypeInferenceTest) TestAdditionWithBigInt() {
	c.AssumingNative("b", rhino.JSTypeNative_BIGINT_TYPE)
	c.AssumingNative("B", rhino.JSTypeNative_BIGINT_OBJECT_TYPE)
	c.AssumingNative("n", rhino.JSTypeNative_NUMBER_TYPE)
	c.AssumingNative("bn", rhino.JSTypeNative_BIGINT_NUMBER)
	c.Assuming("bs", c.CreateUnionType(rhino.JSTypeNative_BIGINT_TYPE, rhino.JSTypeNative_STRING_TYPE))
	c.AssumingNative("s", rhino.JSTypeNative_STRING_TYPE)
	c.AssumingNative("u", rhino.JSTypeNative_UNKNOWN_TYPE)
	c.AssumingNative("ns", rhino.JSTypeNative_NUMBER_STRING)
	c.InFunction(`valueTypePlusSelf = b + b;
objectTypePlusSelf = B + B;
valuePlusObject = b + B;
bigintPlusNumber = b + n;
bigintNumberPlusSelf = bn + bn;
bigintStringConcat = b + s;
bigintNumberStringConcat = bn + s
bigintOtherStringConcat = bs + s
bigintStringConcatWithSelf = bs + bs
bigintPlusUnknown = b + u;
bigintPlusNumberString = b + ns;
`)
	c.VerifyNative("valueTypePlusSelf", rhino.JSTypeNative_BIGINT_TYPE)
	c.VerifyNative("objectTypePlusSelf", rhino.JSTypeNative_BIGINT_TYPE)
	c.VerifyNative("valuePlusObject", rhino.JSTypeNative_BIGINT_TYPE)
	c.VerifyNative("bigintPlusNumber", rhino.JSTypeNative_NO_TYPE)
	c.VerifyNative("bigintNumberPlusSelf", rhino.JSTypeNative_BIGINT_NUMBER)
	c.VerifyNative("bigintStringConcat", rhino.JSTypeNative_STRING_TYPE)
	c.VerifyNative("bigintNumberStringConcat", rhino.JSTypeNative_STRING_TYPE)
	c.VerifyNative("bigintOtherStringConcat", rhino.JSTypeNative_STRING_TYPE)
	c.VerifyNative("bigintStringConcatWithSelf", rhino.JSTypeNative_NO_TYPE)
	c.VerifyNative("bigintPlusUnknown", rhino.JSTypeNative_NO_TYPE)
	c.VerifyNative("bigintPlusNumberString", rhino.JSTypeNative_NO_TYPE)
}

// Java: TypeInferenceTest.testBigIntCompatibleBinaryOperator()
func (c *TypeInferenceTest) TestBigIntCompatibleBinaryOperator() {
	c.AssumingNative("b", rhino.JSTypeNative_BIGINT_TYPE)
	c.AssumingNative("B", rhino.JSTypeNative_BIGINT_OBJECT_TYPE)
	c.AssumingNative("n", rhino.JSTypeNative_NUMBER_TYPE)
	c.AssumingNative("bn", rhino.JSTypeNative_BIGINT_NUMBER)
	c.AssumingNative("s", rhino.JSTypeNative_STRING_TYPE)
	c.AssumingNative("u", rhino.JSTypeNative_UNKNOWN_TYPE)
	c.AssumingNative("ns", rhino.JSTypeNative_NUMBER_STRING)
	c.InFunction(`valueTypeWithSelf = b * b;
objectTypeWithSelf = B * B;
valueWithObject = b * B;
bigintWithNumber = b * n;
bigintNumberWithSelf = bn * bn;
bigintWithOther = b * s;
bigintWithUnknown = b * u;
bigintWithNumberString = b * ns;
`)
	c.VerifyNative("valueTypeWithSelf", rhino.JSTypeNative_BIGINT_TYPE)
	c.VerifyNative("objectTypeWithSelf", rhino.JSTypeNative_BIGINT_TYPE)
	c.VerifyNative("valueWithObject", rhino.JSTypeNative_BIGINT_TYPE)
	c.VerifyNative("bigintWithNumber", rhino.JSTypeNative_NO_TYPE)
	c.VerifyNative("bigintNumberWithSelf", rhino.JSTypeNative_BIGINT_NUMBER)
	c.VerifyNative("bigintWithOther", rhino.JSTypeNative_NO_TYPE)
	c.VerifyNative("bigintWithUnknown", rhino.JSTypeNative_NO_TYPE)
	c.VerifyNative("bigintWithNumberString", rhino.JSTypeNative_NO_TYPE)
}

// Java: TypeInferenceTest.testAssignOpWithBigInt()
func (c *TypeInferenceTest) TestAssignOpWithBigInt() {
	c.AssumingNative("b", rhino.JSTypeNative_BIGINT_TYPE)
	c.AssumingNative("n", rhino.JSTypeNative_NUMBER_TYPE)
	c.AssumingNative("s", rhino.JSTypeNative_STRING_TYPE)
	c.AssumingNative("u", rhino.JSTypeNative_UNKNOWN_TYPE)
	c.AssumingNative("bn", rhino.JSTypeNative_BIGINT_NUMBER)
	c.AssumingNative("bigintWithSelf", rhino.JSTypeNative_BIGINT_TYPE)
	c.AssumingNative("bigintWithNumber", rhino.JSTypeNative_BIGINT_TYPE)
	c.AssumingNative("bigintWithOther", rhino.JSTypeNative_BIGINT_TYPE)
	c.AssumingNative("bigintConcatString", rhino.JSTypeNative_BIGINT_TYPE)
	c.AssumingNative("stringConcatBigInt", rhino.JSTypeNative_STRING_TYPE)
	c.AssumingNative("bigintWithUnknown", rhino.JSTypeNative_BIGINT_TYPE)
	c.AssumingNative("bigintNumberWithSelf", rhino.JSTypeNative_BIGINT_NUMBER)
	c.AssumingNative("bigintNumberWithBigInt", rhino.JSTypeNative_BIGINT_NUMBER)
	c.AssumingNative("bigintNumberWithNumber", rhino.JSTypeNative_BIGINT_NUMBER)
	c.InFunction(`bigintWithSelf *= b;
bigintWithNumber *= n;
bigintWithOther *= s;
bigintConcatString += s
stringConcatBigInt += b
bigintWithUnknown *= u;
bigintNumberWithSelf *= bn;
bigintNumberWithBigInt *= b
bigintNumberWithNumber *= n
`)
	c.VerifyNative("bigintWithSelf", rhino.JSTypeNative_BIGINT_TYPE)
	c.VerifyNative("bigintWithNumber", rhino.JSTypeNative_NO_TYPE)
	c.VerifyNative("bigintWithOther", rhino.JSTypeNative_NO_TYPE)
	c.VerifyNative("bigintConcatString", rhino.JSTypeNative_STRING_TYPE)
	c.VerifyNative("stringConcatBigInt", rhino.JSTypeNative_STRING_TYPE)
	c.VerifyNative("bigintWithUnknown", rhino.JSTypeNative_NO_TYPE)
	c.VerifyNative("bigintNumberWithSelf", rhino.JSTypeNative_BIGINT_NUMBER)
	c.VerifyNative("bigintNumberWithBigInt", rhino.JSTypeNative_NO_TYPE)
	c.VerifyNative("bigintNumberWithNumber", rhino.JSTypeNative_NO_TYPE)
}

// Java: TypeInferenceTest.testUnsignedRightShiftWithBigInt()
func (c *TypeInferenceTest) TestUnsignedRightShiftWithBigInt() {
	c.AssumingNative("b", rhino.JSTypeNative_BIGINT_TYPE)
	c.AssumingNative("n", rhino.JSTypeNative_NUMBER_TYPE)
	c.AssumingNative("assignBigIntOnLeft", rhino.JSTypeNative_BIGINT_TYPE)
	c.AssumingNative("assignBigIntOnRight", rhino.JSTypeNative_NUMBER_TYPE)
	c.AssumingNative("assignBigIntOnBothSides", rhino.JSTypeNative_BIGINT_TYPE)
	c.InFunction(`bigintOnLeft = b >>> n;
bigintOnRight = n >>> b;
bigintOnBothSides = b >>> b;
assignBigIntOnLeft >>>= n
assignBigIntOnRight >>>= b
assignBigIntOnBothSides >>>= b
`)
	c.VerifyNative("bigintOnLeft", rhino.JSTypeNative_NO_TYPE)
	c.VerifyNative("bigintOnRight", rhino.JSTypeNative_NO_TYPE)
	c.VerifyNative("bigintOnBothSides", rhino.JSTypeNative_NO_TYPE)
	c.VerifyNative("assignBigIntOnLeft", rhino.JSTypeNative_NO_TYPE)
	c.VerifyNative("assignBigIntOnRight", rhino.JSTypeNative_NO_TYPE)
	c.VerifyNative("assignBigIntOnBothSides", rhino.JSTypeNative_NO_TYPE)
}

// Java: TypeInferenceTest.testBigIntComparison()
func (c *TypeInferenceTest) TestBigIntComparison() {
	c.AssumingNative("b", rhino.JSTypeNative_BIGINT_TYPE)
	c.AssumingNative("n", rhino.JSTypeNative_NUMBER_TYPE)
	c.InFunction("bigintOnly = b > b; bigintAndOther = b > n")
	c.VerifyNative("bigintOnly", rhino.JSTypeNative_BOOLEAN_TYPE)
	c.VerifyNative("bigintAndOther", rhino.JSTypeNative_BOOLEAN_TYPE)
}

// Java: TypeInferenceTest.testLogicalBinaryOperatorsWithBigInt()
func (c *TypeInferenceTest) TestLogicalBinaryOperatorsWithBigInt() {
	c.AssumingNative("b", rhino.JSTypeNative_BIGINT_TYPE)
	c.AssumingNative("B", rhino.JSTypeNative_BIGINT_OBJECT_TYPE)
	c.AssumingNative("n", rhino.JSTypeNative_NUMBER_TYPE)
	c.AssumingNative("bn", rhino.JSTypeNative_BIGINT_NUMBER)
	c.AssumingNative("s", rhino.JSTypeNative_STRING_TYPE)
	c.AssumingNative("u", rhino.JSTypeNative_UNKNOWN_TYPE)
	c.AssumingNative("ns", rhino.JSTypeNative_NUMBER_STRING)
	c.InFunction(`valueTypeWithSelf = b && b;
objectTypeWithSelf = B && B;
valueWithObject = b && B;
bigintWithNumber = b && n;
bigintNumberWithSelf = bn && bn;
bigintWithOther = b && s;
bigintWithUnknown = b && u;
bigintWithNumberString = b && ns;
`)
	c.VerifyNative("valueTypeWithSelf", rhino.JSTypeNative_BIGINT_TYPE)
	c.VerifyNative("objectTypeWithSelf", rhino.JSTypeNative_BIGINT_OBJECT_TYPE)
	c.Verify("valueWithObject", c.CreateUnionType(rhino.JSTypeNative_BIGINT_TYPE, rhino.JSTypeNative_BIGINT_OBJECT_TYPE))
	c.VerifyNative("bigintWithNumber", rhino.JSTypeNative_BIGINT_NUMBER)
	c.VerifyNative("bigintNumberWithSelf", rhino.JSTypeNative_BIGINT_NUMBER)
	c.Verify("bigintWithOther", c.CreateUnionType(rhino.JSTypeNative_BIGINT_TYPE, rhino.JSTypeNative_STRING_TYPE))
	c.Verify("bigintWithUnknown", c.CreateUnionType(rhino.JSTypeNative_BIGINT_TYPE, rhino.JSTypeNative_UNKNOWN_TYPE))
	c.VerifyNative("bigintWithNumberString", rhino.JSTypeNative_BIGINT_NUMBER_STRING)
}

// Java: TypeInferenceTest.testTernaryOperatorWithBigInt()
func (c *TypeInferenceTest) TestTernaryOperatorWithBigInt() {
	c.AssumingNative("b", rhino.JSTypeNative_BIGINT_TYPE)
	c.AssumingNative("B", rhino.JSTypeNative_BIGINT_OBJECT_TYPE)
	c.AssumingNative("n", rhino.JSTypeNative_NUMBER_TYPE)
	c.AssumingNative("bn", rhino.JSTypeNative_BIGINT_NUMBER)
	c.AssumingNative("s", rhino.JSTypeNative_STRING_TYPE)
	c.AssumingNative("u", rhino.JSTypeNative_UNKNOWN_TYPE)
	c.AssumingNative("v", rhino.JSTypeNative_VOID_TYPE)
	c.AssumingNative("ns", rhino.JSTypeNative_NUMBER_STRING)
	c.InFunction(`valueTypeWithSelf = v ? b : b;
objectTypeWithSelf = v ? B : B;
valueWithObject = v ? b : B;
bigintWithNumber = v ? b : n;
bigintNumberWithSelf = v ? bn : bn;
bigintWithOther = v ? b : s;
bigintWithUnknown = v ? b : u;
bigintWithNumberString = v ? b : ns;
`)
	c.VerifyNative("valueTypeWithSelf", rhino.JSTypeNative_BIGINT_TYPE)
	c.VerifyNative("objectTypeWithSelf", rhino.JSTypeNative_BIGINT_OBJECT_TYPE)
	c.Verify("valueWithObject", c.CreateUnionType(rhino.JSTypeNative_BIGINT_TYPE, rhino.JSTypeNative_BIGINT_OBJECT_TYPE))
	c.VerifyNative("bigintWithNumber", rhino.JSTypeNative_BIGINT_NUMBER)
	c.VerifyNative("bigintNumberWithSelf", rhino.JSTypeNative_BIGINT_NUMBER)
	c.Verify("bigintWithOther", c.CreateUnionType(rhino.JSTypeNative_BIGINT_TYPE, rhino.JSTypeNative_STRING_TYPE))
	c.Verify("bigintWithUnknown", c.CreateUnionType(rhino.JSTypeNative_BIGINT_TYPE, rhino.JSTypeNative_UNKNOWN_TYPE))
	c.VerifyNative("bigintWithNumberString", rhino.JSTypeNative_BIGINT_NUMBER_STRING)
}

// Java: TypeInferenceTest.testAssertBoolean_narrowsAllTypeToBoolean()
func (c *TypeInferenceTest) TestAssertBoolean_narrowsAllTypeToBoolean() {
	var startType rhino.JSType = c.CreateNullableType(rhino.JSTypeNative_ALL_TYPE)
	c.IncludeGoogAssertionFn("assertBoolean", c.GetNativeType(rhino.JSTypeNative_BOOLEAN_TYPE))
	c.Assuming("x", startType)
	c.InFunction("out1 = x; goog.asserts.assertBoolean(x); out2 = x;")
	c.Verify("out1", startType)
	c.VerifyNative("out2", rhino.JSTypeNative_BOOLEAN_TYPE)
}

// Java: TypeInferenceTest.testAssertString_narrowsAllTypeToString()
func (c *TypeInferenceTest) TestAssertString_narrowsAllTypeToString() {
	var startType rhino.JSType = c.CreateNullableType(rhino.JSTypeNative_ALL_TYPE)
	c.IncludeGoogAssertionFn("assertString", c.GetNativeType(rhino.JSTypeNative_STRING_TYPE))
	c.Assuming("x", startType)
	c.InFunction("out1 = x; goog.asserts.assertString(x); out2 = x;")
	c.Verify("out1", startType)
	c.VerifyNative("out2", rhino.JSTypeNative_STRING_TYPE)
}

// Java: TypeInferenceTest.testAssertFunction_narrowsAllTypeToFunction()
func (c *TypeInferenceTest) TestAssertFunction_narrowsAllTypeToFunction() {
	var startType rhino.JSType = c.CreateNullableType(rhino.JSTypeNative_ALL_TYPE)
	c.IncludeGoogAssertionFn("assertFunction", c.GetNativeType(rhino.JSTypeNative_FUNCTION_TYPE))
	c.Assuming("x", startType)
	c.InFunction("out1 = x; goog.asserts.assertFunction(x); out2 = x;")
	c.Verify("out1", startType)
	c.VerifySubtypeOfNative("out2", rhino.JSTypeNative_FUNCTION_TYPE)
}

// Java: TypeInferenceTest.testAssertElement_doesNotChangeElementType()
func (c *TypeInferenceTest) TestAssertElement_doesNotChangeElementType() {
	var elementType rhino.JSType = c.registry.CreateObjectType(jutil.StringPtr("Element"), c.registry.GetNativeObjectType(rhino.JSTypeNative_OBJECT_TYPE))
	c.IncludeGoogAssertionFn("assertElement", elementType)
	c.Assuming("x", elementType)
	c.InFunction("out1 = x; goog.asserts.assertElement(x); out2 = x;")
	c.Verify("out1", elementType)
	c.Verify("out2", elementType)
}

// Java: TypeInferenceTest.testAssertObject_narrowsNullableArrayToArray()
func (c *TypeInferenceTest) TestAssertObject_narrowsNullableArrayToArray() {
	var startType rhino.JSType = c.CreateNullableType(rhino.JSTypeNative_ARRAY_TYPE)
	c.IncludeGoogAssertionFn("assertObject", c.GetNativeType(rhino.JSTypeNative_OBJECT_TYPE))
	c.Assuming("x", startType)
	c.InFunction("out1 = x; goog.asserts.assertObject(x); out2 = x;")
	c.Verify("out1", startType)
	c.VerifyNative("out2", rhino.JSTypeNative_ARRAY_TYPE)
}

// Java: TypeInferenceTest.testAssertObject_narrowsNullableObjectToObject()
func (c *TypeInferenceTest) TestAssertObject_narrowsNullableObjectToObject() {
	var startType rhino.JSType = c.CreateNullableType(rhino.JSTypeNative_OBJECT_TYPE)
	c.IncludeGoogAssertionFn("assertObject", c.GetNativeType(rhino.JSTypeNative_OBJECT_TYPE))
	c.Assuming("x", startType)
	c.InFunction("out1 = x; goog.asserts.assertObject(x); out2 = x;")
	c.Verify("out1", startType)
	c.VerifyNative("out2", rhino.JSTypeNative_OBJECT_TYPE)
}

// Java: TypeInferenceTest.testAssertObject_narrowsQualifiedNameArgument()
func (c *TypeInferenceTest) TestAssertObject_narrowsQualifiedNameArgument() {
	var startType rhino.JSType = c.CreateNullableType(rhino.JSTypeNative_OBJECT_TYPE)
	c.IncludeGoogAssertionFn("assertObject", c.GetNativeType(rhino.JSTypeNative_OBJECT_TYPE))
	c.Assuming("x", c.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE))
	c.Assuming("x.y", startType)
	c.InFunction("out1 = x.y; goog.asserts.assertObject(x.y); out2 = x.y;")
	c.Verify("out1", startType)
	c.VerifyNative("out2", rhino.JSTypeNative_OBJECT_TYPE)
}

// Java: TypeInferenceTest.testAssertObject_inCastToArray_returnsArray()
func (c *TypeInferenceTest) TestAssertObject_inCastToArray_returnsArray() {
	var startType rhino.JSType = c.CreateNullableType(rhino.JSTypeNative_ALL_TYPE)
	c.IncludeGoogAssertionFn("assertObject", c.GetNativeType(rhino.JSTypeNative_OBJECT_TYPE))
	c.Assuming("x", startType)
	c.InFunction(`out1 = x;
out2 = /** @type {!Array} */ (goog.asserts.assertObject(x));
`)
	c.Verify("out1", startType)
	c.VerifyNative("out2", rhino.JSTypeNative_ARRAY_TYPE)
}

// Java: TypeInferenceTest.testAssertArray_narrowsNullableAllTypeToArray()
func (c *TypeInferenceTest) TestAssertArray_narrowsNullableAllTypeToArray() {
	var startType rhino.JSType = c.CreateNullableType(rhino.JSTypeNative_ALL_TYPE)
	c.IncludeGoogAssertionFn("assertArray", c.GetNativeType(rhino.JSTypeNative_ARRAY_TYPE))
	c.Assuming("x", startType)
	c.InFunction("out1 = x; goog.asserts.assertArray(x); out2 = x;")
	c.Verify("out1", startType)
	c.VerifySubtypeOfNative("out2", rhino.JSTypeNative_ARRAY_TYPE)
}

// Java: TypeInferenceTest.testAssertArray_narrowsObjectTypeToArray()
func (c *TypeInferenceTest) TestAssertArray_narrowsObjectTypeToArray() {
	var startType rhino.JSType = c.GetNativeType(rhino.JSTypeNative_OBJECT_TYPE)
	c.IncludeGoogAssertionFn("assertArray", c.GetNativeType(rhino.JSTypeNative_ARRAY_TYPE))
	c.Assuming("x", startType)
	c.InFunction("out1 = x; goog.asserts.assertArray(x); out2 = x;")
	c.Verify("out1", startType)
	c.VerifySubtypeOfNative("out2", rhino.JSTypeNative_ARRAY_TYPE)
}

// Java: TypeInferenceTest.testAssertInstanceof_invalidCall_setsArgToUnknownType()
func (c *TypeInferenceTest) TestAssertInstanceof_invalidCall_setsArgToUnknownType() {
	var startType rhino.JSType = c.CreateNullableType(rhino.JSTypeNative_ALL_TYPE)
	c.IncludeAssertInstanceof()
	c.Assuming("x", startType)
	c.InFunction("out1 = x; goog.asserts.assertInstanceof(x); out2 = x;")
	c.Verify("out1", startType)
	c.VerifyNative("out2", rhino.JSTypeNative_UNKNOWN_TYPE)
}

// Java: TypeInferenceTest.testAssertInstanceof_stringCtor_narrowsAllTypeToString()
func (c *TypeInferenceTest) TestAssertInstanceof_stringCtor_narrowsAllTypeToString() {
	var startType rhino.JSType = c.CreateNullableType(rhino.JSTypeNative_ALL_TYPE)
	c.IncludeAssertInstanceof()
	c.Assuming("x", startType)
	c.InFunction("out1 = x; goog.asserts.assertInstanceof(x, String); out2 = x;")
	c.Verify("out1", startType)
	c.VerifyNative("out2", rhino.JSTypeNative_STRING_OBJECT_TYPE)
}

// Java: TypeInferenceTest.testAssertInstanceof_unknownCtor_setsStringToUnknown()
func (c *TypeInferenceTest) TestAssertInstanceof_unknownCtor_setsStringToUnknown() {
	c.IncludeAssertInstanceof()
	c.AssumingNative("x", rhino.JSTypeNative_STRING_TYPE)
	c.AssumingNative("Foo", rhino.JSTypeNative_UNKNOWN_TYPE)
	c.InFunction("out1 = x; goog.asserts.assertInstanceof(x, Foo); out2 = x;")
	c.VerifyNative("out1", rhino.JSTypeNative_STRING_TYPE)
	c.VerifyNative("out2", rhino.JSTypeNative_UNKNOWN_TYPE)
}

// Java: TypeInferenceTest.testAssertInstanceof_stringCtor_narrowsUnknownToString()
func (c *TypeInferenceTest) TestAssertInstanceof_stringCtor_narrowsUnknownToString() {
	var startType rhino.JSType = c.registry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
	c.IncludeAssertInstanceof()
	c.Assuming("x", startType)
	c.InFunction("out1 = x; goog.asserts.assertInstanceof(x, String); out2 = x;")
	c.Verify("out1", startType)
	c.VerifyNative("out2", rhino.JSTypeNative_STRING_OBJECT_TYPE)
}

// Java: TypeInferenceTest.testAssertInstanceof_objectCtor_doesNotChangeStringType()
func (c *TypeInferenceTest) TestAssertInstanceof_objectCtor_doesNotChangeStringType() {
	var startType rhino.JSType = c.registry.GetNativeType(rhino.JSTypeNative_STRING_OBJECT_TYPE)
	c.IncludeAssertInstanceof()
	c.Assuming("x", startType)
	c.InFunction("out1 = x; goog.asserts.assertInstanceof(x, Object); out2 = x;")
	c.Verify("out1", startType)
	c.VerifyNative("out2", rhino.JSTypeNative_STRING_OBJECT_TYPE)
}

// Java: TypeInferenceTest.testAssertInstanceof_stringCtor_narrowsObjOrVoidToString()
func (c *TypeInferenceTest) TestAssertInstanceof_stringCtor_narrowsObjOrVoidToString() {
	var startType rhino.JSType = c.CreateUnionType(rhino.JSTypeNative_OBJECT_TYPE, rhino.JSTypeNative_VOID_TYPE)
	c.IncludeAssertInstanceof()
	c.Assuming("x", startType)
	c.InFunction("out1 = x; goog.asserts.assertInstanceof(x, String); var r = x;")
	c.Verify("out1", startType)
	c.VerifyNative("x", rhino.JSTypeNative_STRING_OBJECT_TYPE)
}

// Java: TypeInferenceTest.testAssertInstanceof_stringCtor_returnsStringFromObjOrVoid()
func (c *TypeInferenceTest) TestAssertInstanceof_stringCtor_returnsStringFromObjOrVoid() {
	var startType rhino.JSType = c.CreateUnionType(rhino.JSTypeNative_OBJECT_TYPE, rhino.JSTypeNative_VOID_TYPE)
	c.IncludeAssertInstanceof()
	c.Assuming("x", startType)
	c.InFunction("out1 = x; var y = goog.asserts.assertInstanceof(x, String);")
	c.Verify("out1", startType)
	c.VerifyNative("y", rhino.JSTypeNative_STRING_OBJECT_TYPE)
}

// Java: TypeInferenceTest.testTypeInferenceOccursInConstObjectProperties()
func (c *TypeInferenceTest) TestTypeInferenceOccursInConstObjectProperties() {
	c.InFunction(`/** @return {string} */
function foo() { return ''; }

const obj = {
   prop: foo(),
}
LABEL: obj.prop;
`)
	c.AssertTypeOfExpression("LABEL").ToStringIsEqualTo("string")
	// Java resolves isNotEqualTo(JSTypeNative) to Subject.isNotEqualTo(Object).
	c.AssertTypeOfExpression("LABEL").Subject.IsNotEqualTo(rhino.JSTypeNative_UNKNOWN_TYPE)
}

// Java: TypeInferenceTest.testReturn1()
func (c *TypeInferenceTest) TestReturn1() {
	c.Assuming("x", c.CreateNullableType(rhino.JSTypeNative_OBJECT_TYPE))
	c.InFunction(`if (x) { return x; }
x = {};
return x;`)
	c.VerifyNative("x", rhino.JSTypeNative_OBJECT_TYPE)
}

// Java: TypeInferenceTest.testReturn2()
func (c *TypeInferenceTest) TestReturn2() {
	c.Assuming("x", c.CreateNullableType(rhino.JSTypeNative_NUMBER_TYPE))
	c.InFunction(`if (!x) { x = 0; }
return x;`)
	c.VerifyNative("x", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testWhile1()
func (c *TypeInferenceTest) TestWhile1() {
	c.Assuming("x", c.CreateNullableType(rhino.JSTypeNative_NUMBER_TYPE))
	c.InFunction("while (!x) { if (x == null) { x = 0; } else { x = 1; } }")
	c.VerifyNative("x", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testWhile2()
func (c *TypeInferenceTest) TestWhile2() {
	c.Assuming("x", c.CreateNullableType(rhino.JSTypeNative_NUMBER_TYPE))
	c.InFunction("while (!x) { x = {}; }")
	c.VerifySubtypeOf("x", c.CreateUnionType(rhino.JSTypeNative_OBJECT_TYPE, rhino.JSTypeNative_NUMBER_TYPE))
}

// Java: TypeInferenceTest.testDo()
func (c *TypeInferenceTest) TestDo() {
	c.Assuming("x", c.CreateNullableType(rhino.JSTypeNative_OBJECT_TYPE))
	c.InFunction("do { x = 1; } while (!x);")
	c.VerifyNative("x", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testFor1()
func (c *TypeInferenceTest) TestFor1() {
	c.AssumingNative("y", rhino.JSTypeNative_NUMBER_TYPE)
	c.InFunction("var x = null; var i = null; for (i=y; !i; i=1) { x = 1; }")
	c.Verify("x", c.CreateNullableType(rhino.JSTypeNative_NUMBER_TYPE))
	c.VerifyNative("i", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testForInWithExistingVar()
func (c *TypeInferenceTest) TestForInWithExistingVar() {
	c.AssumingNative("y", rhino.JSTypeNative_OBJECT_TYPE)
	c.InFunction(`var x = null;
var i = null;
for (i in y) {
  I_INSIDE_LOOP: i;
  X_AT_LOOP_START: x;
  x = 1;
  X_AT_LOOP_END: x;
}
X_AFTER_LOOP: x;
I_AFTER_LOOP: i;
`)
	c.AssertScopeEnclosing("I_INSIDE_LOOP").Declares("i").
		OnClosestHoistScope()
	c.AssertScopeEnclosing("I_INSIDE_LOOP").Declares("x").
		OnClosestHoistScope()
	c.AssertTypeOfExpression("I_INSIDE_LOOP").ToStringIsEqualTo("string")
	c.AssertTypeOfExpression("I_AFTER_LOOP").ToStringIsEqualTo("(null|string)")
	c.AssertTypeOfExpression("X_AT_LOOP_START").ToStringIsEqualTo("(null|number)")
	c.AssertTypeOfExpression("X_AT_LOOP_END").ToStringIsEqualTo("number")
	c.AssertTypeOfExpression("X_AFTER_LOOP").ToStringIsEqualTo("(null|number)")
}

// Java: TypeInferenceTest.testForInWithRedeclaredVar()
func (c *TypeInferenceTest) TestForInWithRedeclaredVar() {
	c.AssumingNative("y", rhino.JSTypeNative_OBJECT_TYPE)
	c.InFunction(`var i = null;
for (var i in y) { // i redeclared here, but really the same variable
  I_INSIDE_LOOP: i;
}
I_AFTER_LOOP: i;
`)
	c.AssertScopeEnclosing("I_INSIDE_LOOP").Declares("i").
		OnClosestHoistScope()
	c.AssertTypeOfExpression("I_INSIDE_LOOP").ToStringIsEqualTo("string")
	c.AssertScopeEnclosing("I_AFTER_LOOP").Declares("i").
		Directly()
	c.AssertTypeOfExpression("I_AFTER_LOOP").ToStringIsEqualTo("(null|string)")
}

// Java: TypeInferenceTest.testForInWithLet()
func (c *TypeInferenceTest) TestForInWithLet() {
	c.AssumingNative("y", rhino.JSTypeNative_OBJECT_TYPE)
	c.InFunction(`FOR_IN_LOOP: for (let i in y) { // preserve newlines
  I_INSIDE_LOOP: i;
}
AFTER_LOOP: 1;
`)
	c.AssertScopeEnclosing("I_INSIDE_LOOP").Declares("i").
		OnScopeLabeled("FOR_IN_LOOP")
	c.AssertTypeOfExpression("I_INSIDE_LOOP").ToStringIsEqualTo("string")
	c.AssertScopeEnclosing("AFTER_LOOP").DoesNotDeclare("i")
}

// Java: TypeInferenceTest.testForInWithConst()
func (c *TypeInferenceTest) TestForInWithConst() {
	c.AssumingNative("y", rhino.JSTypeNative_OBJECT_TYPE)
	c.InFunction(`FOR_IN_LOOP: for (const i in y) { // preserve newlines
  I_INSIDE_LOOP: i;
}
AFTER_LOOP: 1;
`)
	c.AssertScopeEnclosing("I_INSIDE_LOOP").Declares("i").
		OnScopeLabeled("FOR_IN_LOOP")
	c.AssertTypeOfExpression("I_INSIDE_LOOP").ToStringIsEqualTo("string")
	c.AssertScopeEnclosing("AFTER_LOOP").DoesNotDeclare("i")
}

// Java: TypeInferenceTest.testFor4()
func (c *TypeInferenceTest) TestFor4() {
	c.Assuming("x", c.CreateNullableType(rhino.JSTypeNative_OBJECT_TYPE))
	c.InFunction(`var y = {};
if (x) { for (var i = 0; i < 10; i++) { break; } y = x; }
`)
	c.VerifySubtypeOfNative("y", rhino.JSTypeNative_OBJECT_TYPE)
}

// Java: TypeInferenceTest.testFor5()
func (c *TypeInferenceTest) TestFor5() {
	c.Assuming("y", c.Templatize(c.GetNativeObjectType(rhino.JSTypeNative_ARRAY_TYPE), []rhino.JSType{c.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)}))
	c.InFunction("var x = null; for (var i = 0; i < y.length; i++) { x = y[i]; }")
	c.Verify("x", c.CreateNullableType(rhino.JSTypeNative_NUMBER_TYPE))
	c.VerifyNative("i", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testFor6()
func (c *TypeInferenceTest) TestFor6() {
	c.Assuming("y", c.GetNativeObjectType(rhino.JSTypeNative_ARRAY_TYPE))
	c.InFunction(`var x = null;
for (var i = 0; i < y.length; i++) {
 if (y[i] == 'z') { x = y[i]; }
}
`)
	c.Verify("x", c.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE))
	c.VerifyNative("i", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testSwitch1()
func (c *TypeInferenceTest) TestSwitch1() {
	c.AssumingNative("x", rhino.JSTypeNative_NUMBER_TYPE)
	c.InFunction(`var y = null; switch(x) {
case 1: y = 1; break;
case 2: y = {};
case 3: y = {};
default: y = 0;}
`)
	c.VerifyNative("y", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testSwitch2()
func (c *TypeInferenceTest) TestSwitch2() {
	c.AssumingNative("x", rhino.JSTypeNative_ALL_TYPE)
	c.InFunction(`var y = null; switch (typeof x) {
case 'string':
  y = x;
  return;
default:
  y = 'a';
}
`)
	c.VerifyNative("y", rhino.JSTypeNative_STRING_TYPE)
}

// Java: TypeInferenceTest.testSwitch3()
func (c *TypeInferenceTest) TestSwitch3() {
	c.Assuming("x", c.CreateNullableTypeOfJSType(c.CreateUnionType(rhino.JSTypeNative_NUMBER_TYPE, rhino.JSTypeNative_STRING_TYPE)))
	c.InFunction(`var y; var z; switch (typeof x) {
case 'string':
  y = 1; z = null;
  return;
case 'number':
  y = x; z = null;
  return;
default:
  y = 1; z = x;
}
`)
	c.VerifyNative("y", rhino.JSTypeNative_NUMBER_TYPE)
	c.VerifyNative("z", rhino.JSTypeNative_NULL_TYPE)
}

// Java: TypeInferenceTest.testSwitch4()
func (c *TypeInferenceTest) TestSwitch4() {
	c.AssumingNative("x", rhino.JSTypeNative_ALL_TYPE)
	c.InFunction(`var y = null; switch (typeof x) {
case 'string':
case 'number':
  y = x;
  return;
default:
  y = 1;
}
`)
	c.Verify("y", c.CreateUnionType(rhino.JSTypeNative_NUMBER_TYPE, rhino.JSTypeNative_STRING_TYPE))
}

// Java: TypeInferenceTest.testCall1()
func (c *TypeInferenceTest) TestCall1() {
	c.Assuming("x", c.CreateNullableTypeOfJSType(c.registry.CreateFunctionType(c.registry.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE))))
	c.InFunction("var y = x();")
	c.VerifyNative("y", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testNew1()
func (c *TypeInferenceTest) TestNew1() {
	c.Assuming("x", c.CreateNullableTypeOfJSType(c.registry.GetNativeType(rhino.JSTypeNative_FUNCTION_TYPE)))
	c.InFunction("var y = new x();")
	c.VerifyNative("y", rhino.JSTypeNative_UNKNOWN_TYPE)
}

// Java: TypeInferenceTest.testNew2()
func (c *TypeInferenceTest) TestNew2() {
	c.InFunction(`/**
 * @constructor
 * @param {T} x
 * @template T
 */
function F(x) {}
var x = /** @type {!Array<number>} */ ([]);
var result = new F(x);
`)
	AssertThat(c.GetType("result").String()).IsEqualTo("F<Array<number>>")
}

// Java: TypeInferenceTest.testNew3()
func (c *TypeInferenceTest) TestNew3() {
	c.InFunction(`/**
 * @constructor
 * @param {Array<T>} x
 * @param {T} y
 * @param {S} z
 * @template T,S
 */
function F(x,y,z) {}
var x = /** @type {!Array<number>} */ ([]);
var y = /** @type {string} */ ('foo');
var z = /** @type {boolean} */ (true);
var result = new F(x,y,z);
`)
	AssertThat(c.GetType("result").String()).IsEqualTo("F<(number|string),boolean>")
}

// Java: TypeInferenceTest.testNew4()
func (c *TypeInferenceTest) TestNew4() {
	c.InFunction(`/**
 * @constructor
 * @param {!Array<T>} x
 * @param {T} y
 * @param {S} z
 * @param {U} m
 * @template T,S,U
 */
function F(x,y,z,m) {}
var /** !Array<number> */ x = [];
var y = 'foo';
var z = true;
var m = 9;
var result = new F(x,y,z,m);
`)
	AssertThat(c.GetType("result").String()).IsEqualTo("F<(number|string),boolean,number>")
}

// Java: TypeInferenceTest.testNew_onCtor_instantiatingTemplatizedType_withNoTemplateInformation()
func (c *TypeInferenceTest) TestNew_onCtor_instantiatingTemplatizedType_withNoTemplateInformation() {
	c.InFunction(`/**
 * @constructor
 * @template T
 */
function Foo() {}

var result = new Foo();
`)
	AssertThat(c.GetType("result").String()).IsEqualTo("Foo<?>")
}

// Java: TypeInferenceTest.testNew_onCtor_instantiatingTemplatizedType_specializedOnSecondaryTemplate()
func (c *TypeInferenceTest) TestNew_onCtor_instantiatingTemplatizedType_specializedOnSecondaryTemplate() {
	c.InFunction(`/**
 * @constructor
 * @template T
 */
function Foo() {}

/**
 * @template U
 * @param {function(new:Foo<U>)} ctor
 * @param {U} arg
 * @return {!Foo<U>}
 */
function create(ctor, arg) {
  return new ctor(arg);
}

var result = create(Foo, 0);
`)
	AssertThat(c.GetType("result").String()).IsEqualTo("Foo<number>")
}

// Java: TypeInferenceTest.testNewRest()
func (c *TypeInferenceTest) TestNewRest() {
	c.InFunction(`/**
 * @constructor
 * @param {Array<T>} x
 * @param {T} y
 * @param {...S} rest
 * @template T,S
 */
function F(x, y, ...rest) {}
var x = /** @type {!Array<number>} */ ([]);
var y = /** @type {string} */ ('foo');
var z = /** @type {boolean} */ (true);
var result = new F(x,y,z);
`)
	AssertThat(c.GetType("result").String()).IsEqualTo("F<(number|string),boolean>")
}

// Java: TypeInferenceTest.testParamNodeType_simpleName()
func (c *TypeInferenceTest) TestParamNodeType_simpleName() {
	c.ParseAndRunTypeInference("(/** @param {number} i */ function(i) {})")
	NodeSubject_AssertNode(c.GetParamNameNode("i")).HasJSTypeThat().
		IsNumber()
}

// Java: TypeInferenceTest.testParamNodeType_rest()
func (c *TypeInferenceTest) TestParamNodeType_rest() {
	c.ParseAndRunTypeInference("(/** @param {...number} i */ function(...i) {})")
	NodeSubject_AssertNode(c.GetParamNameNode("i")).HasJSTypeThat().
		ToStringIsEqualTo("Array<number>")
}

// Java: TypeInferenceTest.testParamNodeType_arrayDestructuring()
func (c *TypeInferenceTest) TestParamNodeType_arrayDestructuring() {
	c.ParseAndRunTypeInference("(/** @param {!Iterable<number>} i */ function([i]) {})")
	NodeSubject_AssertNode(c.GetParamNameNode("i")).HasJSTypeThat().
		IsNumber()
}

// Java: TypeInferenceTest.testParamNodeType_objectDestructuring()
func (c *TypeInferenceTest) TestParamNodeType_objectDestructuring() {
	c.ParseAndRunTypeInference("(/** @param {{a: number}} i */ function({a: i}) {})")
	NodeSubject_AssertNode(c.GetParamNameNode("i")).HasJSTypeThat().
		IsNumber()
}

// Java: TypeInferenceTest.testParamNodeType_simpleName_withDefault()
func (c *TypeInferenceTest) TestParamNodeType_simpleName_withDefault() {
	c.ParseAndRunTypeInference("(/** @param {number=} i */ function(i = 9) {})")
	NodeSubject_AssertNode(c.GetParamNameNode("i")).HasJSTypeThat().
		IsNumber()
}

// Java: TypeInferenceTest.testParamNodeType_arrayDestructuring_withDefault()
func (c *TypeInferenceTest) TestParamNodeType_arrayDestructuring_withDefault() {
	c.ParseAndRunTypeInference(`(/** @param {!Iterable<number>=} unused */
function([i] = /** @type ({!Array<number>} */ ([])) {})
`)
	NodeSubject_AssertNode(c.GetParamNameNode("i")).HasJSTypeThat().
		IsUnknown()
}

// Java: TypeInferenceTest.testParamNodeType_objectDestructuring_withDefault()
func (c *TypeInferenceTest) TestParamNodeType_objectDestructuring_withDefault() {
	c.ParseAndRunTypeInference("(/** @param {{a: number}=} i */ function({a: i} = {a: 9}) {})")
	NodeSubject_AssertNode(c.GetParamNameNode("i")).HasJSTypeThat().
		IsNumber()
}

// Java: TypeInferenceTest.testParamNodeType_arrayDestructuring_withDefault_nestedInPattern()
func (c *TypeInferenceTest) TestParamNodeType_arrayDestructuring_withDefault_nestedInPattern() {
	c.ParseAndRunTypeInference("(/** @param {!Iterable<number>} i */ function([i = 9]) {})")
	NodeSubject_AssertNode(c.GetParamNameNode("i")).HasJSTypeThat().
		IsNumber()
}

// Java: TypeInferenceTest.testParamNodeType_objectDestructuring_withDefault_nestedInPattern()
func (c *TypeInferenceTest) TestParamNodeType_objectDestructuring_withDefault_nestedInPattern() {
	c.ParseAndRunTypeInference("(/** @param {{a: number}} i */ function({a: i = 9}) {})")
	NodeSubject_AssertNode(c.GetParamNameNode("i")).HasJSTypeThat().
		IsNumber()
}

// Java: TypeInferenceTest.testInnerFunction1()
func (c *TypeInferenceTest) TestInnerFunction1() {
	c.InFunction("var x = 1; function f() { x = null; };")
	c.VerifyNative("x", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testInnerFunction2()
func (c *TypeInferenceTest) TestInnerFunction2() {
	c.InFunction("var x = 1; var f = function() { x = null; };")
	c.VerifyNative("x", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testFunctionDeclarationHasBlockScope()
func (c *TypeInferenceTest) TestFunctionDeclarationHasBlockScope() {
	c.InFunction(`BLOCK_SCOPE: {
  BEFORE_DEFINITION: f;
  function f() {}
  AFTER_DEFINITION: f;
}
AFTER_BLOCK: f;
`)
	c.AssertScopeEnclosing("BEFORE_DEFINITION").Declares("f").
		OnScopeLabeled("BLOCK_SCOPE")
	c.AssertTypeOfExpression("BEFORE_DEFINITION").ToStringIsEqualTo("function(): undefined")
	c.AssertTypeOfExpression("AFTER_DEFINITION").ToStringIsEqualTo("function(): undefined")
	c.AssertScopeEnclosing("AFTER_BLOCK").DoesNotDeclare("f")
}

// Java: TypeInferenceTest.testHook()
func (c *TypeInferenceTest) TestHook() {
	c.Assuming("x", c.CreateNullableType(rhino.JSTypeNative_OBJECT_TYPE))
	c.InFunction("var y = x ? x : {};")
	c.VerifySubtypeOfNative("y", rhino.JSTypeNative_OBJECT_TYPE)
}

// Java: TypeInferenceTest.testThrow()
func (c *TypeInferenceTest) TestThrow() {
	c.Assuming("x", c.CreateNullableType(rhino.JSTypeNative_NUMBER_TYPE))
	c.InFunction(`var y = 1;
if (x == null) { throw new Error('x is null') }
y = x;
`)
	c.VerifyNative("y", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testTry1()
func (c *TypeInferenceTest) TestTry1() {
	c.AssumingNative("x", rhino.JSTypeNative_NUMBER_TYPE)
	c.InFunction("var y = null; try { y = null; } finally { y = x; }")
	c.VerifyNative("y", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testTry2()
func (c *TypeInferenceTest) TestTry2() {
	c.AssumingNative("x", rhino.JSTypeNative_NUMBER_TYPE)
	c.InFunction(`var y = null;
try {  } catch (e) { y = null; } finally { y = x; }
`)
	c.VerifyNative("y", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testTry3()
func (c *TypeInferenceTest) TestTry3() {
	c.AssumingNative("x", rhino.JSTypeNative_NUMBER_TYPE)
	c.InFunction("var y = null; try { y = x; } catch (e) { }")
	c.VerifyNative("y", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testCatch1()
func (c *TypeInferenceTest) TestCatch1() {
	c.InFunction("var y = null; try { foo(); } catch (e) { y = e; }")
	c.VerifyNative("y", rhino.JSTypeNative_UNKNOWN_TYPE)
}

// Java: TypeInferenceTest.testCatch2()
func (c *TypeInferenceTest) TestCatch2() {
	c.InFunction("var y = null; var e = 3; try { foo(); } catch (e) { y = e; }")
	c.VerifyNative("y", rhino.JSTypeNative_UNKNOWN_TYPE)
}

// Java: TypeInferenceTest.testUnknownType1()
func (c *TypeInferenceTest) TestUnknownType1() {
	c.InFunction("var y = 3; y = x;")
	c.VerifyNative("y", rhino.JSTypeNative_UNKNOWN_TYPE)
}

// Java: TypeInferenceTest.testUnknownType2()
func (c *TypeInferenceTest) TestUnknownType2() {
	c.AssumingNative("x", rhino.JSTypeNative_ARRAY_TYPE)
	c.InFunction("var y = 5; y = x[0];")
	c.VerifyNative("y", rhino.JSTypeNative_UNKNOWN_TYPE)
}

// Java: TypeInferenceTest.testInfiniteLoop1()
func (c *TypeInferenceTest) TestInfiniteLoop1() {
	c.Assuming("x", c.CreateNullableType(rhino.JSTypeNative_OBJECT_TYPE))
	c.InFunction("x = {}; while(x != null) { x = {}; }")
}

// Java: TypeInferenceTest.testInfiniteLoop2()
func (c *TypeInferenceTest) TestInfiniteLoop2() {
	c.Assuming("x", c.CreateNullableType(rhino.JSTypeNative_OBJECT_TYPE))
	c.InFunction("x = {}; do { x = null; } while (x == null);")
}

// Java: TypeInferenceTest.testJoin1()
func (c *TypeInferenceTest) TestJoin1() {
	var unknownOrNull rhino.JSType = c.CreateUnionType(rhino.JSTypeNative_NULL_TYPE, rhino.JSTypeNative_UNKNOWN_TYPE)
	c.AssumingNative("x", rhino.JSTypeNative_BOOLEAN_TYPE)
	c.Assuming("unknownOrNull", unknownOrNull)
	c.InFunction("var y; if (x) y = unknownOrNull; else y = null;")
	c.Verify("y", unknownOrNull)
}

// Java: TypeInferenceTest.testJoin2()
func (c *TypeInferenceTest) TestJoin2() {
	var unknownOrNull rhino.JSType = c.CreateUnionType(rhino.JSTypeNative_NULL_TYPE, rhino.JSTypeNative_UNKNOWN_TYPE)
	c.AssumingNative("x", rhino.JSTypeNative_BOOLEAN_TYPE)
	c.Assuming("unknownOrNull", unknownOrNull)
	c.InFunction("var y; if (x) y = null; else y = unknownOrNull;")
	c.Verify("y", unknownOrNull)
}

// Java: TypeInferenceTest.testArrayLit()
func (c *TypeInferenceTest) TestArrayLit() {
	c.Assuming("x", c.CreateNullableType(rhino.JSTypeNative_OBJECT_TYPE))
	c.InFunction("var y = 3; if (x) { x = [y = x]; }")
	c.Verify("x", c.CreateUnionType(rhino.JSTypeNative_NULL_TYPE, rhino.JSTypeNative_ARRAY_TYPE))
	c.Verify("y", c.CreateUnionType(rhino.JSTypeNative_NUMBER_TYPE, rhino.JSTypeNative_OBJECT_TYPE))
}

// Java: TypeInferenceTest.testGetElem()
func (c *TypeInferenceTest) TestGetElem() {
	c.Assuming("x", c.CreateNullableType(rhino.JSTypeNative_OBJECT_TYPE))
	c.InFunction("var y = 3; if (x) { x = x[y = x]; }")
	c.VerifyNative("x", rhino.JSTypeNative_UNKNOWN_TYPE)
	c.Verify("y", c.CreateUnionType(rhino.JSTypeNative_NUMBER_TYPE, rhino.JSTypeNative_OBJECT_TYPE))
}

// Java: TypeInferenceTest.testEnumRAI1()
func (c *TypeInferenceTest) TestEnumRAI1() {
	var enumType rhino.JSType = c.CreateEnumType("MyEnum", rhino.JSTypeNative_ARRAY_TYPE).GetElementsType()
	c.Assuming("x", enumType)
	c.InFunction("var y = null; if (x) y = x;")
	c.Verify("y", c.CreateNullableTypeOfJSType(enumType))
}

// Java: TypeInferenceTest.testEnumRAI2()
func (c *TypeInferenceTest) TestEnumRAI2() {
	var enumType rhino.JSType = c.CreateEnumType("MyEnum", rhino.JSTypeNative_NUMBER_TYPE).GetElementsType()
	c.Assuming("x", enumType)
	c.InFunction("var y = null; if (typeof x == 'number') y = x;")
	c.Verify("y", c.CreateNullableTypeOfJSType(enumType))
}

// Java: TypeInferenceTest.testEnumRAI3()
func (c *TypeInferenceTest) TestEnumRAI3() {
	var enumType rhino.JSType = c.CreateEnumType("MyEnum", rhino.JSTypeNative_NUMBER_TYPE).GetElementsType()
	c.Assuming("x", enumType)
	c.InFunction("var y = null; if (x && typeof x == 'number') y = x;")
	c.Verify("y", c.CreateNullableTypeOfJSType(enumType))
}

// Java: TypeInferenceTest.testEnumRAI4()
func (c *TypeInferenceTest) TestEnumRAI4() {
	var enumType rhino.JSType = c.CreateEnumTypeOfJSType("MyEnum", c.CreateUnionType(rhino.JSTypeNative_STRING_TYPE, rhino.JSTypeNative_NUMBER_TYPE)).GetElementsType()
	c.Assuming("x", enumType)
	c.InFunction("var y = null; if (typeof x == 'number') y = x;")
	c.Verify("y", c.CreateNullableType(rhino.JSTypeNative_NUMBER_TYPE))
}

// Java: TypeInferenceTest.testShortCircuitingAnd()
func (c *TypeInferenceTest) TestShortCircuitingAnd() {
	c.AssumingNative("x", rhino.JSTypeNative_NUMBER_TYPE)
	c.InFunction("var y = null; if (x && (y = 3)) { }")
	c.Verify("y", c.CreateNullableType(rhino.JSTypeNative_NUMBER_TYPE))
}

// Java: TypeInferenceTest.testShortCircuitingAnd2()
func (c *TypeInferenceTest) TestShortCircuitingAnd2() {
	c.AssumingNative("x", rhino.JSTypeNative_NUMBER_TYPE)
	c.InFunction("var y = null; var z = 4; if (x && (y = 3)) { z = y; }")
	c.VerifyNative("z", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testShortCircuitingOr()
func (c *TypeInferenceTest) TestShortCircuitingOr() {
	c.AssumingNative("x", rhino.JSTypeNative_NUMBER_TYPE)
	c.InFunction("var y = null; if (x || (y = 3)) { }")
	c.Verify("y", c.CreateNullableType(rhino.JSTypeNative_NUMBER_TYPE))
}

// Java: TypeInferenceTest.testShortCircuitingNullishCoalesce()
func (c *TypeInferenceTest) TestShortCircuitingNullishCoalesce() {
	c.AssumingNative("x", rhino.JSTypeNative_NUMBER_TYPE)
	c.InFunction("var y = null; if (x ?? (y = 3)) { }")
	c.VerifyNative("y", rhino.JSTypeNative_NULL_TYPE)
}

// Java: TypeInferenceTest.testShortCircuitingNullishCoalesceIf()
func (c *TypeInferenceTest) TestShortCircuitingNullishCoalesceIf() {
	c.AssumingNative("x", rhino.JSTypeNative_NUMBER_TYPE)
	c.InFunction("var y = null; var z = 5; if (x ?? (y = 3)) { z = y }")
	c.VerifyNative("y", rhino.JSTypeNative_NULL_TYPE)
}

// Java: TypeInferenceTest.testShortCircuitingOr2()
func (c *TypeInferenceTest) TestShortCircuitingOr2() {
	c.AssumingNative("x", rhino.JSTypeNative_NUMBER_TYPE)
	c.InFunction("var y = null; var z = 4; if (x || (y = 3)) { z = y; }")
	c.Verify("z", c.CreateNullableType(rhino.JSTypeNative_NUMBER_TYPE))
}

// Java: TypeInferenceTest.testShortCircuitingNullishCoalseceNumber()
func (c *TypeInferenceTest) TestShortCircuitingNullishCoalseceNumber() {
	c.AssumingNative("x", rhino.JSTypeNative_NUMBER_TYPE)
	c.InFunction("var z = x ?? null")
	c.VerifyNative("z", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testNullishCoalesce()
func (c *TypeInferenceTest) TestNullishCoalesce() {
	c.Assuming("x", c.CreateNullableType(rhino.JSTypeNative_OBJECT_TYPE))
	c.InFunction("var y = 1; x ?? (y = x);")
	c.Verify("y", c.CreateNullableType(rhino.JSTypeNative_NUMBER_TYPE))
}

// Java: TypeInferenceTest.nullishCoalesceWithHook()
func (c *TypeInferenceTest) Test_nullishCoalesceWithHook() {
	c.Assuming("x", c.CreateNullableType(rhino.JSTypeNative_OBJECT_TYPE))
	c.InFunction("var z = (x) ?? (x ? 'hi' : false)")
	c.Verify("z", c.CreateMultiParamUnionType(rhino.JSTypeNative_OBJECT_TYPE, rhino.JSTypeNative_BOOLEAN_TYPE, rhino.JSTypeNative_STRING_TYPE))
}

// Java: TypeInferenceTest.nullishCoalesceRemoveNull()
func (c *TypeInferenceTest) Test_nullishCoalesceRemoveNull() {
	c.Assuming("x", c.CreateNullableType(rhino.JSTypeNative_NUMBER_TYPE))
	c.InFunction("x = x ?? 3")
	c.VerifyNative("x", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.nullishCoalesceZeroIsValid()
func (c *TypeInferenceTest) Test_nullishCoalesceZeroIsValid() {
	c.Assuming("x", c.CreateNullableType(rhino.JSTypeNative_NUMBER_TYPE))
	c.InFunction("var y = ''; if (x ?? (y = x)) { }")
	c.Verify("y", c.CreateNullableType(rhino.JSTypeNative_STRING_TYPE))
}

// Java: TypeInferenceTest.nullishCoalesceFalseIsValid()
func (c *TypeInferenceTest) Test_nullishCoalesceFalseIsValid() {
	c.Assuming("x", c.CreateNullableType(rhino.JSTypeNative_BOOLEAN_TYPE))
	c.InFunction("var y = ''; x ?? (y = x)")
	c.Verify("y", c.CreateNullableType(rhino.JSTypeNative_STRING_TYPE))
}

// Java: TypeInferenceTest.testAssignInCondition()
func (c *TypeInferenceTest) TestAssignInCondition() {
	c.Assuming("x", c.CreateNullableType(rhino.JSTypeNative_NUMBER_TYPE))
	c.InFunction("var y; if (!(y = x)) { y = 3; }")
	c.VerifyNative("y", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testInstanceOf1()
func (c *TypeInferenceTest) TestInstanceOf1() {
	c.AssumingNative("x", rhino.JSTypeNative_OBJECT_TYPE)
	c.InFunction("var y = null; if (x instanceof String) y = x;")
	c.Verify("y", c.CreateNullableType(rhino.JSTypeNative_STRING_OBJECT_TYPE))
}

// Java: TypeInferenceTest.testInstanceOf2()
func (c *TypeInferenceTest) TestInstanceOf2() {
	c.Assuming("x", c.CreateNullableType(rhino.JSTypeNative_OBJECT_TYPE))
	c.InFunction("var y = 1; if (x instanceof String) y = x;")
	c.Verify("y", c.CreateUnionType(rhino.JSTypeNative_STRING_OBJECT_TYPE, rhino.JSTypeNative_NUMBER_TYPE))
}

// Java: TypeInferenceTest.testInstanceOf3()
func (c *TypeInferenceTest) TestInstanceOf3() {
	c.Assuming("x", c.CreateUnionType(rhino.JSTypeNative_STRING_OBJECT_TYPE, rhino.JSTypeNative_NUMBER_OBJECT_TYPE))
	c.InFunction("var y = null; if (x instanceof String) y = x;")
	c.Verify("y", c.CreateNullableType(rhino.JSTypeNative_STRING_OBJECT_TYPE))
}

// Java: TypeInferenceTest.testInstanceOf4()
func (c *TypeInferenceTest) TestInstanceOf4() {
	c.Assuming("x", c.CreateUnionType(rhino.JSTypeNative_STRING_OBJECT_TYPE, rhino.JSTypeNative_NUMBER_OBJECT_TYPE))
	c.InFunction("var y = null; if (x instanceof String); else y = x;")
	c.Verify("y", c.CreateNullableType(rhino.JSTypeNative_NUMBER_OBJECT_TYPE))
}

// Java: TypeInferenceTest.testInstanceOf5()
func (c *TypeInferenceTest) TestInstanceOf5() {
	c.AssumingNative("x", rhino.JSTypeNative_OBJECT_TYPE)
	c.InFunction("var y = null; if (x instanceof String); else y = x;")
	c.Verify("y", c.CreateNullableType(rhino.JSTypeNative_OBJECT_TYPE))
}

// Java: TypeInferenceTest.testInstanceOf6()
func (c *TypeInferenceTest) TestInstanceOf6() {
	var startType rhino.JSType = c.registry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
	c.Assuming("x", startType)
	c.InFunction("out1 = x; if (x instanceof String) out2 = x;")
	c.Verify("out1", startType)
	c.VerifyNative("out2", rhino.JSTypeNative_STRING_OBJECT_TYPE)
}

// Java: TypeInferenceTest.testUnary()
func (c *TypeInferenceTest) TestUnary() {
	c.AssumingNative("x", rhino.JSTypeNative_NUMBER_TYPE)
	c.InFunction("var y = +x;")
	c.VerifyNative("y", rhino.JSTypeNative_NUMBER_TYPE)
	c.InFunction("var z = -x;")
	c.VerifyNative("z", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testAdd1()
func (c *TypeInferenceTest) TestAdd1() {
	c.AssumingNative("x", rhino.JSTypeNative_NUMBER_TYPE)
	c.InFunction("var y = x + 5;")
	c.VerifyNative("y", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testAdd2()
func (c *TypeInferenceTest) TestAdd2() {
	c.AssumingNative("x", rhino.JSTypeNative_NUMBER_TYPE)
	c.InFunction("var y = x + '5';")
	c.VerifyNative("y", rhino.JSTypeNative_STRING_TYPE)
}

// Java: TypeInferenceTest.testAdd3()
func (c *TypeInferenceTest) TestAdd3() {
	c.AssumingNative("x", rhino.JSTypeNative_NUMBER_TYPE)
	c.InFunction("var y = '5' + x;")
	c.VerifyNative("y", rhino.JSTypeNative_STRING_TYPE)
}

// Java: TypeInferenceTest.testAssignAdd()
func (c *TypeInferenceTest) TestAssignAdd() {
	c.AssumingNative("x", rhino.JSTypeNative_NUMBER_TYPE)
	c.InFunction("x += '5';")
	c.VerifyNative("x", rhino.JSTypeNative_STRING_TYPE)
}

// Java: TypeInferenceTest.testAssignAdd_lhsNodeInputType()
func (c *TypeInferenceTest) TestAssignAdd_lhsNodeInputType() {
	c.InFunction("L: { /** @type {*} */ let data; data += ''; }")
	block := c.GetLabeledStatement("L").statementNode
	// In L: { let data; data += ''; }, the block has two children: LET and EXPR_RESULT.
	exprResult := block.GetLastChild()
	assignAdd := exprResult.GetFirstChild()
	lhs := assignAdd.GetFirstChild()
	// Pre-fix: LHS node under ASSIGN_ADD gets overwritten with the result type (STRING_TYPE).
	// Post-fix: LHS node preserves its input type (UNKNOWN_TYPE / ALL_TYPE).
	NodeSubject_AssertNode(lhs).HasJSTypeThat().IsEqualTo(c.GetNativeType(rhino.JSTypeNative_ALL_TYPE))
}

// Java: TypeInferenceTest.testAssignOrToNumeric()
func (c *TypeInferenceTest) TestAssignOrToNumeric() {
	c.InFunction("var y = false; y ||= 20")
	c.Verify("y", c.CreateUnionType(rhino.JSTypeNative_NUMBER_TYPE, rhino.JSTypeNative_BOOLEAN_TYPE))
}

// Java: TypeInferenceTest.testComputedClassFieldsInControlFlow()
func (c *TypeInferenceTest) TestComputedClassFieldsInControlFlow() {
	c.InFunction(`let y;
class Foo {
  static [y = null] = (y = '');
  [y = false] = [y = null];
}
`)
	c.VerifyNative("y", rhino.JSTypeNative_STRING_TYPE)
}

// Java: TypeInferenceTest.testClassFieldsInControlFlow()
func (c *TypeInferenceTest) TestClassFieldsInControlFlow() {
	c.InFunction(`let y;
class Foo {
  static y = (y = '');
  z = [y = null];
}
`)
	c.VerifyNative("y", rhino.JSTypeNative_STRING_TYPE)
}

// Java: TypeInferenceTest.testThisTypeAfterMemberField()
func (c *TypeInferenceTest) TestThisTypeAfterMemberField() {
	var thisType rhino.JSType = c.CreateRecordTypeNative("x", rhino.JSTypeNative_STRING_TYPE)
	c.AssumingThisType(thisType)
	c.InFunction(`let thisDotX;
(class C {
  /** @type {number} */
  static x = 0;
}, thisDotX = this.x);
`)
	c.VerifyNative("thisDotX", rhino.JSTypeNative_STRING_TYPE)
}

// Java: TypeInferenceTest.testFunction()
func (c *TypeInferenceTest) TestFunction() {
	c.InScript(`let y;
function foo() {
  y = 'hi';
}
foo();
`)
	c.VerifyNative("y", rhino.JSTypeNative_VOID_TYPE)
}

// Java: TypeInferenceTest.testClassStaticBlock()
func (c *TypeInferenceTest) TestClassStaticBlock() {
	c.InScript(`let y;
class Foo {
  static {
    y = 'hi';
  }
}
`)
	c.VerifyNative("y", rhino.JSTypeNative_VOID_TYPE)
}

// Java: TypeInferenceTest.testSuper()
func (c *TypeInferenceTest) TestSuper() {
	c.InScript(`class Foo {
  static str;
}
class Bar extends Foo {
  static {
    super.str = 'hi';
  }
}
let x = Bar.str;
`)
	c.VerifyNative("x", rhino.JSTypeNative_ALL_TYPE)
}

// Java: TypeInferenceTest.testAssignOrNoAssign()
func (c *TypeInferenceTest) TestAssignOrNoAssign() {
	c.InFunction("var x; var y; y = false; x = (y || 'foo');")
	c.Verify("x", c.CreateUnionType(rhino.JSTypeNative_BOOLEAN_TYPE, rhino.JSTypeNative_STRING_TYPE))
	c.InFunction("var a; var b; var c; a = true; c = (a || (b = 'foo'));")
	c.Verify("c", c.CreateUnionType(rhino.JSTypeNative_BOOLEAN_TYPE, rhino.JSTypeNative_STRING_TYPE))
	c.VerifyNative("a", rhino.JSTypeNative_BOOLEAN_TYPE)
	c.Verify("b", c.CreateUnionType(rhino.JSTypeNative_VOID_TYPE, rhino.JSTypeNative_STRING_TYPE))
	c.InFunction("var y; y = true; y ||= 'foo';")
	c.Verify("y", c.CreateUnionType(rhino.JSTypeNative_BOOLEAN_TYPE, rhino.JSTypeNative_STRING_TYPE))
}

// Java: TypeInferenceTest.testAssignOrToBooleanEitherAbsoluteFalseOrTrue()
func (c *TypeInferenceTest) TestAssignOrToBooleanEitherAbsoluteFalseOrTrue() {
	c.AssumingNative("x", rhino.JSTypeNative_NULL_TYPE)
	c.InFunction("x ||= 'foo';")
	c.VerifyNative("x", rhino.JSTypeNative_STRING_TYPE)
	c.AssumingNative("y", rhino.JSTypeNative_OBJECT_TYPE)
	c.InFunction("y ||= 'foo';")
	c.VerifyNative("y", rhino.JSTypeNative_OBJECT_TYPE)
}

// Java: TypeInferenceTest.testAssignOrLHSFalsyRHSTruthy()
func (c *TypeInferenceTest) TestAssignOrLHSFalsyRHSTruthy() {
	c.AssumingNative("x", rhino.JSTypeNative_NULL_TYPE)
	c.AssumingNative("obj", rhino.JSTypeNative_OBJECT_TYPE)
	c.InFunction("x ||= obj;")
	c.VerifyNative("x", rhino.JSTypeNative_OBJECT_TYPE)
}

// Java: TypeInferenceTest.testAssignAndToNumeric()
func (c *TypeInferenceTest) TestAssignAndToNumeric() {
	c.InFunction("var y = true; y &&= 20")
	c.Verify("y", c.CreateUnionType(rhino.JSTypeNative_NUMBER_TYPE, rhino.JSTypeNative_BOOLEAN_TYPE))
}

// Java: TypeInferenceTest.testAssignAndNoAssign()
func (c *TypeInferenceTest) TestAssignAndNoAssign() {
	c.InFunction("var x = true && 'foo';")
	c.Verify("x", c.CreateUnionType(rhino.JSTypeNative_BOOLEAN_TYPE, rhino.JSTypeNative_STRING_TYPE))
	c.InFunction("var y = false; y &&= 'foo';")
	c.Verify("y", c.CreateUnionType(rhino.JSTypeNative_BOOLEAN_TYPE, rhino.JSTypeNative_STRING_TYPE))
}

// Java: TypeInferenceTest.testAssignAndToBooleanEitherAbsoluteFalseOrTrue()
func (c *TypeInferenceTest) TestAssignAndToBooleanEitherAbsoluteFalseOrTrue() {
	c.AssumingNative("x", rhino.JSTypeNative_NULL_TYPE)
	c.InFunction("x &&= 'foo';")
	c.VerifyNative("x", rhino.JSTypeNative_NULL_TYPE)
	c.AssumingNative("y", rhino.JSTypeNative_OBJECT_TYPE)
	c.InFunction("y &&= 'foo';")
	c.VerifyNative("y", rhino.JSTypeNative_STRING_TYPE)
}

// Java: TypeInferenceTest.testAssignAndLHSTruthyRHSFalsy()
func (c *TypeInferenceTest) TestAssignAndLHSTruthyRHSFalsy() {
	c.AssumingNative("x", rhino.JSTypeNative_OBJECT_TYPE)
	c.AssumingNative("null", rhino.JSTypeNative_NULL_TYPE)
	c.InFunction("x &&= null;")
	c.VerifyNative("x", rhino.JSTypeNative_NULL_TYPE)
}

// Java: TypeInferenceTest.testAssignCoalesceToNumeric()
func (c *TypeInferenceTest) TestAssignCoalesceToNumeric() {
	c.AssumingNative("x", rhino.JSTypeNative_NULL_TYPE)
	c.InFunction("x ??= 10;")
	c.VerifyNative("x", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testAssignCoalesceNoAssign()
func (c *TypeInferenceTest) TestAssignCoalesceNoAssign() {
	c.AssumingNative("x", rhino.JSTypeNative_STRING_TYPE)
	c.InFunction("x ??= 10;")
	c.VerifyNative("x", rhino.JSTypeNative_STRING_TYPE)
}

// Java: TypeInferenceTest.testAssignCoalesceRHSAssignmentScope()
func (c *TypeInferenceTest) TestAssignCoalesceRHSAssignmentScope() {
	c.InFunction("var y; var x = null; x ??= (y = 6)")
	c.Verify("y", c.CreateUnionType(rhino.JSTypeNative_VOID_TYPE, rhino.JSTypeNative_NUMBER_TYPE))
	c.VerifyNative("x", rhino.JSTypeNative_NUMBER_TYPE)
	c.InFunction("var y; var x = true; x ??= (y = 'a' + 6)")
	c.VerifyNative("y", rhino.JSTypeNative_VOID_TYPE)
	c.VerifyNative("x", rhino.JSTypeNative_BOOLEAN_TYPE)
}

// Java: TypeInferenceTest.testAssignCoalesceJoin()
func (c *TypeInferenceTest) TestAssignCoalesceJoin() {
	c.Assuming("x", c.CreateNullableType(rhino.JSTypeNative_NUMBER_TYPE))
	c.InFunction("var y = ''; x ??= (y = x, 1)")
	c.Verify("y", c.CreateNullableType(rhino.JSTypeNative_STRING_TYPE))
	c.VerifyNative("x", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testComparison()
func (c *TypeInferenceTest) TestComparison() {
	c.InFunction("var x = 'foo'; var y = (x = 3) < 4;")
	c.VerifyNative("x", rhino.JSTypeNative_NUMBER_TYPE)
	c.InFunction("var x = 'foo'; var y = (x = 3) > 4;")
	c.VerifyNative("x", rhino.JSTypeNative_NUMBER_TYPE)
	c.InFunction("var x = 'foo'; var y = (x = 3) <= 4;")
	c.VerifyNative("x", rhino.JSTypeNative_NUMBER_TYPE)
	c.InFunction("var x = 'foo'; var y = (x = 3) >= 4;")
	c.VerifyNative("x", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testComparisonWithBigInt()
func (c *TypeInferenceTest) TestComparisonWithBigInt() {
	c.InFunction("var x = 'foo'; var y = (x = 3n) < 4;")
	c.VerifyNative("x", rhino.JSTypeNative_BIGINT_TYPE)
	c.InFunction("var x = 'foo'; var y = (x = 3n) > 4;")
	c.VerifyNative("x", rhino.JSTypeNative_BIGINT_TYPE)
	c.InFunction("var x = 'foo'; var y = (x = 3n) <= 4;")
	c.VerifyNative("x", rhino.JSTypeNative_BIGINT_TYPE)
	c.InFunction("var x = 'foo'; var y = (x = 3n) >= 4;")
	c.VerifyNative("x", rhino.JSTypeNative_BIGINT_TYPE)
}

// Java: TypeInferenceTest.testThrownExpression()
func (c *TypeInferenceTest) TestThrownExpression() {
	c.InFunction(`var x = 'foo';
try { throw new Error(x = 3); } catch (ex) {}
`)
	c.VerifyNative("x", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testObjectLit()
func (c *TypeInferenceTest) TestObjectLit() {
	c.InFunction("var x = {}; var out = x.a;")
	c.VerifyNative("out", rhino.JSTypeNative_UNKNOWN_TYPE)
	c.InFunction("var x = {a:1}; var out = x.a;")
	c.VerifyNative("out", rhino.JSTypeNative_NUMBER_TYPE)
	c.InFunction("var x = {a:1}; var out = x.a; x.a = 'string'; var out2 = x.a;")
	c.VerifyNative("out", rhino.JSTypeNative_NUMBER_TYPE)
	c.VerifyNative("out2", rhino.JSTypeNative_STRING_TYPE)
	c.InFunction("var x = { get a() {return 1} }; var out = x.a;")
	c.VerifyNative("out", rhino.JSTypeNative_UNKNOWN_TYPE)
	c.InFunction(`var x = {
  /** @return {number} */ get a() {return 1}
};
var out = x.a;
`)
	c.VerifyNative("out", rhino.JSTypeNative_NUMBER_TYPE)
	c.InFunction("var x = { set a(b) {} }; var out = x.a;")
	c.VerifyNative("out", rhino.JSTypeNative_UNKNOWN_TYPE)
	c.InFunction(`var x = {
/** @param {number} b */ set a(b) {} };
var out = x.a;
`)
	c.VerifyNative("out", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testObjectSpread_isInferredToBeObject()
func (c *TypeInferenceTest) TestObjectSpread_isInferredToBeObject() {
	var recordType rhino.JSType = c.registry.CreateRecordType(jutil.NewImmutableMapBuilder[string, rhino.JSType]().Put("x", c.GetNativeType(rhino.JSTypeNative_STRING_TYPE)).Put("y", c.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)).BuildOrThrow())
	c.Assuming("obj", recordType)
	c.AssumingNative("before", rhino.JSTypeNative_BOOLEAN_TYPE)
	c.AssumingNative("after", rhino.JSTypeNative_NULL_TYPE)
	c.InFunction("let spread = {before, ...obj, after};")
	c.VerifyNative("spread", rhino.JSTypeNative_OBJECT_TYPE)
}

// Java: TypeInferenceTest.testCast1()
func (c *TypeInferenceTest) TestCast1() {
	c.InFunction("var x = /** @type {Object} */ (this);")
	c.Verify("x", c.CreateNullableType(rhino.JSTypeNative_OBJECT_TYPE))
}

// Java: TypeInferenceTest.testCast2()
func (c *TypeInferenceTest) TestCast2() {
	c.InFunction(`/** @return {boolean} */
Object.prototype.method = function() { return true; };
var x = /** @type {Object} */ (this).method;
`)
	c.Verify("x", c.registry.CreateFunctionTypeWithInstanceType(c.registry.GetNativeObjectType(rhino.JSTypeNative_OBJECT_TYPE), c.registry.GetNativeType(rhino.JSTypeNative_BOOLEAN_TYPE), []rhino.JSType{}))
}

// Java: TypeInferenceTest.testBackwardsInferenceCall()
func (c *TypeInferenceTest) TestBackwardsInferenceCall() {
	c.InFunction(`/** @param {{foo: (number|undefined)}} x */
function f(x) {}
var y = {};
f(y);
`)
	AssertThat(c.GetType("y").String()).IsEqualTo("{foo: (number|undefined)}")
}

// Java: TypeInferenceTest.testBackwardsInferenceCallRestParameter()
func (c *TypeInferenceTest) TestBackwardsInferenceCallRestParameter() {
	c.InFunction(`/** @param {...{foo: (number|undefined)}} rest */
function f(...rest) {}
var y = {};
f(y);
`)
	AssertThat(c.GetType("y").String()).IsEqualTo("{foo: (number|undefined)}")
}

// Java: TypeInferenceTest.testBackwardsInferenceNew()
func (c *TypeInferenceTest) TestBackwardsInferenceNew() {
	c.InFunction(`/**
 * @constructor
 * @param {{foo: (number|undefined)}} x
 */
function F(x) {}
var y = {};
new F(y);
`)
	AssertThat(c.GetType("y").String()).IsEqualTo("{foo: (number|undefined)}")
}

// Java: TypeInferenceTest.testNoThisInference()
func (c *TypeInferenceTest) TestNoThisInference() {
	var thisType rhino.JSType = c.CreateNullableType(rhino.JSTypeNative_OBJECT_TYPE)
	c.AssumingThisType(thisType)
	c.InFunction("var out = 3; if (goog.isNull(this)) out = this;")
	c.Verify("out", c.CreateUnionType(rhino.JSTypeNative_OBJECT_TYPE, rhino.JSTypeNative_NUMBER_TYPE))
}

// Java: TypeInferenceTest.testRecordInference()
func (c *TypeInferenceTest) TestRecordInference() {
	c.InFunction(`/** @param {{a: boolean}|{b: string}} x */
function f(x) {}
var out = {};
f(out);
`)
	AssertThat(c.GetType("out").String()).IsEqualTo(`{
  a: (boolean|undefined),
  b: (string|undefined)
}`)
}

// Java: TypeInferenceTest.testLotsOfBranchesGettingMerged()
func (c *TypeInferenceTest) TestLotsOfBranchesGettingMerged() {
	code := `var a = -1;
`
	code += `switch(foo()) { 
`
	for i := 0; i < 100; i++ {
		code += "case " + strconv.Itoa(i) + ": a = " + strconv.Itoa(i) + `; break; 
`
	}
	code += `default: a = undefined; break;
`
	code += `}
`
	c.InFunction(code)
	AssertThat(c.GetType("a").String()).IsEqualTo("(number|undefined)")
}

// Java: TypeInferenceTest.testIssue785()
func (c *TypeInferenceTest) TestIssue785() {
	c.InFunction(`/** @param {string|{prop: (string|undefined)}} x */
function f(x) {}
var out = {};
f(out);
`)
	AssertThat(c.GetType("out").String()).IsEqualTo("{prop: (string|undefined)}")
}

// Java: TypeInferenceTest.testFunctionTemplateType_specializedFunctionType_copiesColorIdCompnents()
func (c *TypeInferenceTest) TestFunctionTemplateType_specializedFunctionType_copiesColorIdCompnents() {
	c.InFunction(`/**
 * @template T
 * @param {T} a
 * @return {T}
 */
function f(a) {}
TEMPLATE: f;

SPECIALIZED: f(10);
`)
	var templateFn rhino.FunctionType = c.GetLabeledStatement("TEMPLATE").statementNode.GetFirstChild().GetJSType().
		ToMaybeFunctionType()
	var specializedFn rhino.FunctionType = c.GetLabeledStatement("SPECIALIZED").statementNode.GetFirstFirstChild().GetJSType().
		ToMaybeFunctionType()
	AssertThat(specializedFn).IsNotEqualTo(templateFn)
	AssertThat(specializedFn.GetReferenceName()).IsEqualTo(templateFn.GetReferenceName())
	AssertThat(specializedFn.GetSource()).IsEqualTo(templateFn.GetSource())
	AssertThat(specializedFn.GetGoogModuleId()).IsEqualTo(templateFn.GetGoogModuleId())
}

// Java: TypeInferenceTest.testFunctionTemplateType_literalParam()
func (c *TypeInferenceTest) TestFunctionTemplateType_literalParam() {
	c.InFunction(`/**
 * @template T
 * @param {T} a
 * @return {T}
 */
function f(a){}

var result = f(10);
`)
	c.VerifyNative("result", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testFunctionTemplateType_unionsPossibilities()
func (c *TypeInferenceTest) TestFunctionTemplateType_unionsPossibilities() {
	c.InFunction(`/**
 * @template T
 * @param {T} a
 * @param {T} b
 * @return {T}
 */
function f(a, b){}

var result = f(10, 'x');
`)
	c.Verify("result", c.registry.CreateUnionTypeOfNatives(rhino.JSTypeNative_NUMBER_TYPE, rhino.JSTypeNative_STRING_TYPE))
}

// Java: TypeInferenceTest.testFunctionTemplateType_willUseUnknown()
func (c *TypeInferenceTest) TestFunctionTemplateType_willUseUnknown() {
	c.InFunction(`/**
 * @template T
 * @param {T} a
 * @return {T}
 */
function f(a){}

var result = f(/** @type {?} */ ({}));
`)
	c.VerifyNative("result", rhino.JSTypeNative_UNKNOWN_TYPE)
}

// Java: TypeInferenceTest.testFunctionTemplateType_willUseUnknown_butPrefersTighterTypes()
func (c *TypeInferenceTest) TestFunctionTemplateType_willUseUnknown_butPrefersTighterTypes() {
	c.InFunction("/**\n * @template T\n * @param {T} a\n * @param {T} b\n * @param {T} c\n * @return {T}\n */\nfunction f(a, b, c){}\n\n// Make sure `?` is dispreferred before *and* after a known type.\nvar result = f('x', /** @type {?} */ ({}), 5);\n")
	c.Verify("result", c.registry.CreateUnionTypeOfNatives(rhino.JSTypeNative_NUMBER_TYPE, rhino.JSTypeNative_STRING_TYPE))
}

// Java: TypeInferenceTest.testFunctionTemplateType_recursesIntoFunctionParams()
func (c *TypeInferenceTest) TestFunctionTemplateType_recursesIntoFunctionParams() {
	c.InFunction(`/**
 * @template T
 * @param {function(T)} a
 * @return {T}
 */
function f(a){}

var result = f(function(/** number */ a) { });
`)
	c.VerifyNative("result", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testFunctionTemplateType_recursesIntoFunctionParams_throughUnknown()
func (c *TypeInferenceTest) TestFunctionTemplateType_recursesIntoFunctionParams_throughUnknown() {
	c.InFunction(`/**
 * @template T
 * @param {function(T)=} a
 * @return {T}
 */
function f(a){}

var result = f(/** @type {?} */ ({}));
`)
	c.VerifyNative("result", rhino.JSTypeNative_UNKNOWN_TYPE)
}

// Java: TypeInferenceTest.testFunctionTemplateType_unpacksUnions_fromParamType()
func (c *TypeInferenceTest) TestFunctionTemplateType_unpacksUnions_fromParamType() {
	c.InFunction(`/**
 * @template T
 * @param {!Iterable<T>|number} a
 * @return {T}
 */
function f(a){}

var result = f(/** @type {!Iterable<number>} */ ({}));
`)
	c.VerifyNative("result", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testFunctionTemplateType_unpacksUnions_fromArgType()
func (c *TypeInferenceTest) TestFunctionTemplateType_unpacksUnions_fromArgType() {
	c.InFunction(`/**
 * @template T
 * @param {!Iterable<T>} a
 * @return {T}
 */
function f(a){}

// The arg type is illegal, but the inference should still work.
var result = f(/** @type {!Iterable<number>|number} */ ({}));
`)
	c.VerifyNative("result", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testFunctionTemplateType_unpacksUnions_fromArgType_acrossSubtypes()
func (c *TypeInferenceTest) TestFunctionTemplateType_unpacksUnions_fromArgType_acrossSubtypes() {
	c.InFunction(`/**
 * @template T
 * @param {!Iterable<T>} a
 * @return {T}
 */
function f(a){}

var result = f(/** @type {!Array<number>|!Generator<string>} */ ({}));
`)
	c.Verify("result", c.registry.CreateUnionTypeOfNatives(rhino.JSTypeNative_NUMBER_TYPE, rhino.JSTypeNative_STRING_TYPE))
}

// Java: TypeInferenceTest.testTypeTransformationTypePredicate()
func (c *TypeInferenceTest) TestTypeTransformationTypePredicate() {
	c.InFunction(`/**
 * @return {R}
 * @template R := 'number' =:
 */
function f(a){}
var result = f(10);
`)
	c.VerifyNative("result", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testTypeTransformationConditional()
func (c *TypeInferenceTest) TestTypeTransformationConditional() {
	c.InFunction(`/**
 * @param {T} a
 * @param {N} b
 * @return {R}
 * @template T, N
 * @template R := cond( eq(T, N), 'string', 'boolean' ) =:
 */
function f(a, b){}
var result = f(1, 2);
var result2 = f(1, 'a');
`)
	c.VerifyNative("result", rhino.JSTypeNative_STRING_TYPE)
	c.VerifyNative("result2", rhino.JSTypeNative_BOOLEAN_TYPE)
}

// Java: TypeInferenceTest.testTypeTransformationNoneType()
func (c *TypeInferenceTest) TestTypeTransformationNoneType() {
	c.InFunction(`/**
 * @return {R}
 * @template R := none() =:
 */
function f(){}
var result = f(10);
`)
	c.VerifyNative("result", rhino.JSTypeNative_UNKNOWN_TYPE)
}

// Java: TypeInferenceTest.testTypeTransformationUnionType()
func (c *TypeInferenceTest) TestTypeTransformationUnionType() {
	c.InFunction(`/**
 * @param {S} a
 * @param {N} b
 * @return {R}
 * @template S, N
 * @template R := union(S, N) =:
 */
function f(a, b) {}
var result = f(1, 'a');
`)
	c.Verify("result", c.CreateUnionType(rhino.JSTypeNative_STRING_TYPE, rhino.JSTypeNative_NUMBER_TYPE))
}

// Java: TypeInferenceTest.testTypeTransformationMapunion()
func (c *TypeInferenceTest) TestTypeTransformationMapunion() {
	c.InFunction(`/**
 * @param {U} a
 * @return {R}
 * @template U
 * @template R :=
 * mapunion(U, (x) => cond(eq(x, 'string'), 'boolean', 'null'))
 * =:
 */
function f(a) {}
/** @type {string|number} */ var x;
var result = f(x);
`)
	c.Verify("result", c.CreateUnionType(rhino.JSTypeNative_BOOLEAN_TYPE, rhino.JSTypeNative_NULL_TYPE))
}

// Java: TypeInferenceTest.testTypeTransformationObjectUseCase()
func (c *TypeInferenceTest) TestTypeTransformationObjectUseCase() {
	c.InFunction(`/** 
 * @param {T} a
 * @return {R}
 * @template T 
 * @template R := 
 * mapunion(T, (x) => 
 *      cond(eq(x, 'string'), 'String',
 *      cond(eq(x, 'number'), 'Number',
 *      cond(eq(x, 'boolean'), 'Boolean',
 *      cond(eq(x, 'null'), 'Object', 
 *      cond(eq(x, 'undefined'), 'Object',
 *      x)))))) 
 * =:
 */
function Object(a) {}
/** @type {(string|number|boolean)} */
var o;
var r = Object(o);
`)
	c.Verify("r", c.CreateMultiParamUnionType(rhino.JSTypeNative_STRING_OBJECT_TYPE, rhino.JSTypeNative_NUMBER_OBJECT_TYPE, rhino.JSTypeNative_BOOLEAN_OBJECT_TYPE))
}

// Java: TypeInferenceTest.testTypeTransformationObjectUseCase2()
func (c *TypeInferenceTest) TestTypeTransformationObjectUseCase2() {
	c.InFunction(`/** 
 * @param {T} a
 * @return {R}
 * @template T 
 * @template R := 
 * mapunion(T, (x) => 
 *      cond(eq(x, 'string'), 'String',
 *      cond(eq(x, 'number'), 'Number',
 *      cond(eq(x, 'boolean'), 'Boolean',
 *      cond(eq(x, 'null'), 'Object', 
 *      cond(eq(x, 'undefined'), 'Object',
 *      x)))))) 
 * =:
 */
function fn(a) {}
/** @type {(string|null|undefined)} */
var o;
var r = fn(o);
`)
	c.VerifyNative("r", rhino.JSTypeNative_OBJECT_TYPE)
}

// Java: TypeInferenceTest.testTypeTransformationObjectUseCase3()
func (c *TypeInferenceTest) TestTypeTransformationObjectUseCase3() {
	c.InFunction(`/** 
 * @param {T} a
 * @return {R}
 * @template T 
 * @template R := 
 * mapunion(T, (x) => 
 *      cond(eq(x, 'string'), 'String',
 *      cond(eq(x, 'number'), 'Number',
 *      cond(eq(x, 'boolean'), 'Boolean',
 *      cond(eq(x, 'null'), 'Object', 
 *      cond(eq(x, 'undefined'), 'Object',
 *      x)))))) 
 * =:
 */
function fn(a) {}
/** @type {(Array|undefined)} */
var o;
var r = fn(o);
`)
	c.VerifyNative("r", rhino.JSTypeNative_OBJECT_TYPE)
}

// Java: TypeInferenceTest.testTypeTransformationTypeOfVarWithInstanceOfConstructor()
func (c *TypeInferenceTest) TestTypeTransformationTypeOfVarWithInstanceOfConstructor() {
	c.InFunction(`/** @constructor */
function Bar() {}
var b = new Bar();
/** 
 * @return {R}
 * @template R := typeOfVar('b') =:
 */
function f(){}
var r = f();
`)
	c.Verify("r", c.GetType("b"))
}

// Java: TypeInferenceTest.testTypeTransformationTypeOfVarWithConstructor()
func (c *TypeInferenceTest) TestTypeTransformationTypeOfVarWithConstructor() {
	c.InFunction(`/** @constructor */
function Bar() {}
/** 
 * @return {R}
 * @template R := typeOfVar('Bar') =:
 */
function f(){}
var r = f();
`)
	c.Verify("r", c.GetType("Bar"))
}

// Java: TypeInferenceTest.testTypeTransformationTypeOfVarWithTypedef()
func (c *TypeInferenceTest) TestTypeTransformationTypeOfVarWithTypedef() {
	c.InFunction(`/** @typedef {(string|number)} */
var NumberLike;
/** @type {!NumberLike} */
var x;
/**
 * @return {R}
 * @template R := typeOfVar('x') =:
 */
function f(){}
var r = f();
`)
	c.Verify("r", c.GetType("x"))
}

// Java: TypeInferenceTest.testTypeTransformationWithTypeFromConstructor()
func (c *TypeInferenceTest) TestTypeTransformationWithTypeFromConstructor() {
	c.InFunction(`/** @constructor */
function Bar(){}
var x = new Bar();
/** 
 * @return {R}
 * @template R := 'Bar' =:
 */
function f(){}
var r = f();
`)
	c.Verify("r", c.GetType("x"))
}

// Java: TypeInferenceTest.testTypeTransformationWithTypeFromTypedef()
func (c *TypeInferenceTest) TestTypeTransformationWithTypeFromTypedef() {
	c.InFunction(`/** @typedef {(string|number)} */
var NumberLike;
/** @type {!NumberLike} */
var x;
/**
 * @return {R}
 * @template R := 'NumberLike' =:
 */
function f(){}
var r = f();
`)
	c.Verify("r", c.CreateUnionType(rhino.JSTypeNative_STRING_TYPE, rhino.JSTypeNative_NUMBER_TYPE))
}

// Java: TypeInferenceTest.testTypeTransformationWithTypeFromNamespace()
func (c *TypeInferenceTest) TestTypeTransformationWithTypeFromNamespace() {
	c.InFunction(`var wiz
/** @constructor */
wiz.async.Response = function() {};
/**
 * @return {R}
 * @template R := typeOfVar('wiz.async.Response') =:
 */
function f(){}
var r = f();
`)
	c.Verify("r", c.GetType("wiz.async.Response"))
}

// Java: TypeInferenceTest.testTypeTransformationWithNativeTypeExpressionFunction()
func (c *TypeInferenceTest) TestTypeTransformationWithNativeTypeExpressionFunction() {
	c.InFunction(`/** @type {function(string, boolean)} */
var x;
/**
 * @return {R}
 * @template R := typeExpr('function(string, boolean)') =:
 */
function f(){}
var r = f();
`)
	c.Verify("r", c.GetType("x"))
}

// Java: TypeInferenceTest.testTypeTransformationWithNativeTypeExpressionFunctionReturn()
func (c *TypeInferenceTest) TestTypeTransformationWithNativeTypeExpressionFunctionReturn() {
	c.InFunction(`/** @type {function(): number} */
var x;
/**
 * @return {R}
 * @template R := typeExpr('function(): number') =:
 */
function f(){}
var r = f();
`)
	c.Verify("r", c.GetType("x"))
}

// Java: TypeInferenceTest.testTypeTransformationWithNativeTypeExpressionFunctionThis()
func (c *TypeInferenceTest) TestTypeTransformationWithNativeTypeExpressionFunctionThis() {
	c.InFunction(`/** @type {function(this:boolean, string)} */
var x;
/**
 * @return {R}
 * @template R := typeExpr('function(this:boolean, string)') =:
 */
function f(){}
var r = f();
`)
	c.Verify("r", c.GetType("x"))
}

// Java: TypeInferenceTest.testTypeTransformationWithNativeTypeExpressionFunctionVarargs()
func (c *TypeInferenceTest) TestTypeTransformationWithNativeTypeExpressionFunctionVarargs() {
	c.InFunction(`/** @type {function(string, ...number): number} */
var x;
/**
 * @return {R}
 * @template R := typeExpr('function(string, ...number): number') =:
 */
function f(){}
var r = f();
`)
	c.Verify("r", c.GetType("x"))
}

// Java: TypeInferenceTest.testTypeTransformationWithNativeTypeExpressionFunctionOptional()
func (c *TypeInferenceTest) TestTypeTransformationWithNativeTypeExpressionFunctionOptional() {
	c.InFunction(`/** @type {function(?string=, number=)} */
var x;
/**
 * @return {R}
 * @template R := typeExpr('function(?string=, number=)') =:
 */
function f(){}
var r = f();
`)
	c.Verify("r", c.GetType("x"))
}

// Java: TypeInferenceTest.testTypeTransformationRecordFromObject()
func (c *TypeInferenceTest) TestTypeTransformationRecordFromObject() {
	c.InFunction(`/** 
 * @param {T} a
 * @return {R}
 * @template T 
 * @template R := record(T) =:
 */
function f(a) {}
/** @type {{foo:?}} */
var e;
/** @type {?} */
var bar;
var r = f({foo:bar});
`)
	AssertThat(c.GetType("r").IsRecordType()).IsTrue()
	c.Verify("r", c.GetType("e"))
}

// Java: TypeInferenceTest.testTypeTransformationRecordFromObjectNested()
func (c *TypeInferenceTest) TestTypeTransformationRecordFromObjectNested() {
	c.InFunction(`/** 
 * @param {T} a
 * @return {R}
 * @template T 
 * @template R :=
 * maprecord(record(T), (k, v) => record({[k]:record(v)})) =:
 */
function f(a) {}
/** @type {{foo:!Object, bar:!Object}} */
var e;
var r = f({foo:{}, bar:{}});
`)
	AssertThat(c.GetType("r").IsRecordType()).IsTrue()
	c.Verify("r", c.GetType("e"))
}

// Java: TypeInferenceTest.testTypeTransformationRecordFromObjectWithTemplatizedType()
func (c *TypeInferenceTest) TestTypeTransformationRecordFromObjectWithTemplatizedType() {
	c.InFunction(`/** 
 * @param {T} a
 * @return {R}
 * @template T 
 * @template R := record(T) =:
 */
function f(a) {}
/** @type {{foo:!Array<number>}} */
var e;
/** @type {!Array<number>} */
var something;
var r = f({foo:something});
`)
	AssertThat(c.GetType("r").IsRecordType()).IsTrue()
	c.Verify("r", c.GetType("e"))
}

// Java: TypeInferenceTest.testTypeTransformationIsTemplatizedPartially()
func (c *TypeInferenceTest) TestTypeTransformationIsTemplatizedPartially() {
	c.InFunction(`/**
 * @constructor
 * @template T, U
 */
function Foo() {}
/**
 * @template T := cond(isTemplatized(type('Foo', 'number')), 'number', 'string') =:
 * @return {T}
 */
function f() { return 123; }
var x = f();
`)
	AssertThat(c.GetType("x").IsNumber()).IsTrue()
}

// Java: TypeInferenceTest.testAssertTypeofProp()
func (c *TypeInferenceTest) TestAssertTypeofProp() {
	c.Assuming("x", c.CreateNullableType(rhino.JSTypeNative_OBJECT_TYPE))
	c.InFunction(`goog.asserts.assert(typeof x.prop != 'undefined');
out = x.prop;
`)
	c.VerifyNative("out", rhino.JSTypeNative_CHECKED_UNKNOWN_TYPE)
}

// Java: TypeInferenceTest.testIsArray()
func (c *TypeInferenceTest) TestIsArray() {
	c.Assuming("x", c.CreateNullableType(rhino.JSTypeNative_OBJECT_TYPE))
	c.InFunction("goog.asserts.assert(Array.isArray(x));")
	c.VerifyNative("x", rhino.JSTypeNative_ARRAY_TYPE)
}

// Java: TypeInferenceTest.testNotIsArray()
func (c *TypeInferenceTest) TestNotIsArray() {
	c.Assuming("x", c.CreateUnionType(rhino.JSTypeNative_ARRAY_TYPE, rhino.JSTypeNative_NUMBER_TYPE))
	c.InFunction("goog.asserts.assert(!Array.isArray(x));")
	c.VerifyNative("x", rhino.JSTypeNative_NUMBER_TYPE)
}

// Java: TypeInferenceTest.testYield1()
func (c *TypeInferenceTest) TestYield1() {
	c.InGenerator("var x = yield 3;")
	c.Verify("x", c.registry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE))
}

// Java: TypeInferenceTest.testYield2()
func (c *TypeInferenceTest) TestYield2() {
	c.InGenerator(`var obj;
yield (obj = {a: 3, b: '4'});
var a = obj.a;
var b = obj.b;
`)
	c.Verify("a", c.registry.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE))
	c.Verify("b", c.registry.GetNativeType(rhino.JSTypeNative_STRING_TYPE))
}

// Java: TypeInferenceTest.testTemplateLiteral1()
func (c *TypeInferenceTest) TestTemplateLiteral1() {
	c.InFunction("var x = `foobar`; X: x;")
	c.AssertTypeOfExpression("X").IsString()
}

// Java: TypeInferenceTest.testSpreadExpression()
func (c *TypeInferenceTest) TestSpreadExpression() {
	c.InFunction(`let x = 1; // x is initially a number
let y = [...[x = 'hi', 'there']]; // reassign x a string in the spread
X: x;
`)
	c.AssertTypeOfExpression("X").ToStringIsEqualTo("string")
}

// Java: TypeInferenceTest.testTaggedTemplateLiteral1()
func (c *TypeInferenceTest) TestTaggedTemplateLiteral1() {
	c.Assuming("getNumber", c.registry.CreateFunctionType(c.registry.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)))
	c.InFunction("var num = getNumber``; NUM: num;")
	c.AssertTypeOfExpression("NUM").IsNumber()
}

// Java: TypeInferenceTest.testRestParamType()
func (c *TypeInferenceTest) TestRestParamType() {
	c.ParseAndRunTypeInference(`(
/** // preserve newlines
 * @param {...number} nums
 */
function(str, ...nums) {
  NUMS: nums;
  let n = null;
  N_START: n;
  if (nums.length > 0) {
    n = nums[0];
    N_IF_TRUE: n;
  } else {
    N_IF_FALSE: n;
  }
  N_FINAL: n;
}
);
`)
	c.AssertTypeOfExpression("N_START").ToStringIsEqualTo("null")
	c.AssertTypeOfExpression("N_IF_TRUE").ToStringIsEqualTo("number")
	c.AssertTypeOfExpression("N_IF_FALSE").ToStringIsEqualTo("null")
	c.AssertTypeOfExpression("N_FINAL").ToStringIsEqualTo("(null|number)")
}

// Java: TypeInferenceTest.testObjectDestructuringDeclarationInference()
func (c *TypeInferenceTest) TestObjectDestructuringDeclarationInference() {
	var recordType rhino.JSType = c.registry.CreateRecordType(jutil.NewImmutableMapBuilder[string, rhino.JSType]().Put("x", c.GetNativeType(rhino.JSTypeNative_STRING_TYPE)).Put("y", c.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)).BuildOrThrow())
	c.Assuming("obj", recordType)
	c.InFunction(`let {x, y} = obj;  // preserve newline
X: x;
Y: y;
`)
	c.AssertTypeOfExpression("X").ToStringIsEqualTo("string")
	c.AssertTypeOfExpression("Y").ToStringIsEqualTo("number")
	c.AssertScopeEnclosing("X").Declares("x").
		WithTypeThat().
		ToStringIsEqualTo("string")
}

// Java: TypeInferenceTest.testObjectDestructuringDeclarationInferenceWithDefaultValue()
func (c *TypeInferenceTest) TestObjectDestructuringDeclarationInferenceWithDefaultValue() {
	c.InFunction(`var /** {x: (?string|undefined)} */ obj;
let {x = 3} = obj;  // preserve newline
X: x;
`)
	c.AssertTypeOfExpression("X").ToStringIsEqualTo("(null|number|string)")
}

// Java: TypeInferenceTest.testObjectDestructuringDeclarationInferenceWithUnnecessaryDefaultValue()
func (c *TypeInferenceTest) TestObjectDestructuringDeclarationInferenceWithUnnecessaryDefaultValue() {
	c.InFunction(`var /** {x: string} */ obj;
let {x = 3} = obj;  // we ignore the default value's type
X: x;
`)
	c.AssertTypeOfExpression("X").ToStringIsEqualTo("(number|string)")
}

// Java: TypeInferenceTest.testObjectDestructuringDeclarationInference_unknownRhsAndKnownDefaultValue()
func (c *TypeInferenceTest) TestObjectDestructuringDeclarationInference_unknownRhsAndKnownDefaultValue() {
	c.InFunction(`var /** ? */ obj;
let {x = 3} = obj;  // preserve newline
X: x;
`)
	c.AssertTypeOfExpression("X").ToStringIsEqualTo("?")
}

// Java: TypeInferenceTest.testObjectDestructuringDeclarationInference_knownRhsAndUnknownDefaultValue()
func (c *TypeInferenceTest) TestObjectDestructuringDeclarationInference_knownRhsAndUnknownDefaultValue() {
	c.InFunction(`var /** {x: (string|undefined)} */ obj;
let {x = someUnknown} = obj;  // preserve newline
X: x;
`)
	c.AssertTypeOfExpression("X").ToStringIsEqualTo("?")
}

// Java: TypeInferenceTest.testObjectDestructuringDeclaration_defaultValueEvaluatedAfterComputedProperty()
func (c *TypeInferenceTest) TestObjectDestructuringDeclaration_defaultValueEvaluatedAfterComputedProperty() {
	c.InFunction(`var /** !Object<string, (number|undefined)> */ obj = {};
var a = 1;
const {[a = 'string']: b = a} = obj
A: a
B: b
`)
	c.AssertTypeOfExpression("A").ToStringIsEqualTo("string")
	c.AssertTypeOfExpression("B").ToStringIsEqualTo("(number|string)")
}

// Java: TypeInferenceTest.testObjectDestructuringDeclarationInferenceWithUnknownProperty()
func (c *TypeInferenceTest) TestObjectDestructuringDeclarationInferenceWithUnknownProperty() {
	var recordType rhino.JSType = c.registry.CreateRecordType(jutil.NewImmutableMapBuilder[string, rhino.JSType]().BuildOrThrow())
	c.Assuming("obj", recordType)
	c.InFunction(`let {x} = obj;  // preserve newline
X: x;
`)
	c.AssertTypeOfExpression("X").ToStringIsEqualTo("?")
}

// Java: TypeInferenceTest.testObjectDestructuringDoesInferenceWithinComputedProp()
func (c *TypeInferenceTest) TestObjectDestructuringDoesInferenceWithinComputedProp() {
	c.InFunction(`let y = 'foobar';  // preserve newline
let {[y = 3]: z} = {};
Y: y
Z: z
`)
	c.AssertTypeOfExpression("Y").ToStringIsEqualTo("number")
	c.AssertTypeOfExpression("Z").ToStringIsEqualTo("?")
}

// Java: TypeInferenceTest.testObjectDestructuringUsesIObjectTypeForComputedProp()
func (c *TypeInferenceTest) TestObjectDestructuringUsesIObjectTypeForComputedProp() {
	c.InFunction(`let /** !IObject<string, number> */ myObj = {['foo']: 3};  // preserve newline
let {[42]: x} = myObj;
X: x
`)
	c.AssertTypeOfExpression("X").ToStringIsEqualTo("number")
}

// Java: TypeInferenceTest.testObjectDestructuringDeclarationWithNestedPattern()
func (c *TypeInferenceTest) TestObjectDestructuringDeclarationWithNestedPattern() {
	c.InFunction(`let /** {a: {b: number}} */ obj = {a: {b: 3}};
let {a: {b: x}} = obj;
X: x
`)
	c.AssertTypeOfExpression("X").ToStringIsEqualTo("number")
}

// Java: TypeInferenceTest.testObjectDestructuringAssignmentToQualifiedName()
func (c *TypeInferenceTest) TestObjectDestructuringAssignmentToQualifiedName() {
	c.InFunction(`const ns = {};
({x: ns.x} = {x: 3});
X: ns.x;
`)
	c.AssertTypeOfExpression("X").ToStringIsEqualTo("number")
}

// Java: TypeInferenceTest.testObjectDestructuringDeclarationInForOf()
func (c *TypeInferenceTest) TestObjectDestructuringDeclarationInForOf() {
	c.InFunction(`const /** !Iterable<{x: number}> */ data = [{x: 3}];
for (let {x} of data) {
  X: x;
}
`)
	c.AssertTypeOfExpression("X").ToStringIsEqualTo("number")
}

// Java: TypeInferenceTest.testObjectDestructuringAssignInForOf()
func (c *TypeInferenceTest) TestObjectDestructuringAssignInForOf() {
	c.InFunction(`const /** !Iterable<{x: number}> */ data = [{x: 3}];
var x;
for ({x} of data) {
  X: x;
}
`)
	c.AssertTypeOfExpression("X").ToStringIsEqualTo("number")
}

// Java: TypeInferenceTest.testObjectDestructuringParameterWithDefaults()
func (c *TypeInferenceTest) TestObjectDestructuringParameterWithDefaults() {
	c.ParseAndRunTypeInference("(/** @param {{x: (number|undefined)}} data */ function f({x = 3}) { X: x; });")
	c.AssertTypeOfExpression("X").ToStringIsEqualTo("number")
}

// Java: TypeInferenceTest.testObjectRest_inferredGivenObjectLiteralType()
func (c *TypeInferenceTest) TestObjectRest_inferredGivenObjectLiteralType() {
	c.InFunction("var obj = {a: 1, b: 2, c: 3}; const {a, ...rest} = obj;  A: a; REST: rest;")
	c.AssertTypeOfExpression("A").ToStringIsEqualTo("number")
	c.AssertTypeOfExpression("REST").IsEqualTo(c.registry.GetNativeType(rhino.JSTypeNative_OBJECT_TYPE))
}

// Java: TypeInferenceTest.testObjectLiteralNoSideEffect()
func (c *TypeInferenceTest) TestObjectLiteralNoSideEffect() {
	c.InFunction(`  for (let x = 0; x < 3; x++) {
    obj = {
      data: {tipsMetadata: ''},
      ...{}
    };
  }
`)
	AssertWithMessage("Type inference must not alter OBJECT_TYPE.").That(c.registry.GetNativeObjectType(rhino.JSTypeNative_OBJECT_TYPE).HasOwnProperty("data")).
		IsFalse()
	AssertWithMessage("Cannot resolve type of 'obj'").That(c.GetType("obj")).
		IsNotNull()
	AssertWithMessage("Expect property 'data' to be defined on 'obj'").That(c.registry.CanPropertyBeDefined(c.GetType("obj"), "data")).
		IsNotEqualTo(rhino.JSTypeRegistry_PropDefinitionKind_UNKNOWN)
}

// Java: TypeInferenceTest.testArrayDestructuringDeclaration()
func (c *TypeInferenceTest) TestArrayDestructuringDeclaration() {
	c.InFunction(`const /** !Iterable<number> */ numbers = [1, 2, 3];
let [x, y] = numbers;
X: x
Y: y
`)
	c.AssertTypeOfExpression("X").ToStringIsEqualTo("number")
	c.AssertTypeOfExpression("Y").ToStringIsEqualTo("number")
}

// Java: TypeInferenceTest.testArrayDestructuringDeclarationWithDefaultValue()
func (c *TypeInferenceTest) TestArrayDestructuringDeclarationWithDefaultValue() {
	c.InFunction(`const /** !Iterable<(number|undefined)> */ numbers = [1, 2, 3];
let [x = 'x', y = 'y'] = numbers;
X: x
Y: y
`)
	c.AssertTypeOfExpression("X").ToStringIsEqualTo("(number|string)")
	c.AssertTypeOfExpression("Y").ToStringIsEqualTo("(number|string)")
}

// Java: TypeInferenceTest.testArrayDestructuringDeclarationWithDefaultValueForNestedPattern()
func (c *TypeInferenceTest) TestArrayDestructuringDeclarationWithDefaultValueForNestedPattern() {
	c.InFunction(`const /** !Iterable<({x: number}|undefined)> */ xNumberObjs = [];
let [{x = 'foo'} = {}] = xNumberObjs;
X: x
Y: y
`)
	c.AssertTypeOfExpression("X").ToStringIsEqualTo("(number|string)")
}

// Java: TypeInferenceTest.testArrayDestructuringDeclarationWithRest()
func (c *TypeInferenceTest) TestArrayDestructuringDeclarationWithRest() {
	c.InFunction(`const /** !Iterable<number> */ numbers = [1, 2, 3];
let [x, ...y] = numbers;
X: x
Y: y
`)
	c.AssertTypeOfExpression("X").ToStringIsEqualTo("number")
	c.AssertTypeOfExpression("Y").ToStringIsEqualTo("Array<number>")
}

// Java: TypeInferenceTest.testArrayDestructuringDeclarationWithNestedArrayPattern()
func (c *TypeInferenceTest) TestArrayDestructuringDeclarationWithNestedArrayPattern() {
	c.InFunction(`const /** !Iterable<!Iterable<number>> */ numbers = [[1, 2, 3]];
let [[x], y] = numbers;
X: x
Y: y
`)
	c.AssertTypeOfExpression("X").ToStringIsEqualTo("number")
	c.AssertTypeOfExpression("Y").ToStringIsEqualTo("Iterable<number,?,?>")
}

// Java: TypeInferenceTest.testArrayDestructuringDeclarationWithNestedObjectPattern()
func (c *TypeInferenceTest) TestArrayDestructuringDeclarationWithNestedObjectPattern() {
	c.InFunction(`const /** !Iterable<{x: number}> */ numbers = [{x: 3}, {x: 4}];
let [{x}, {x: y}] = numbers;
X: x
Y: y
`)
	c.AssertTypeOfExpression("X").ToStringIsEqualTo("number")
	c.AssertTypeOfExpression("Y").ToStringIsEqualTo("number")
}

// Java: TypeInferenceTest.testArrayDestructuringDeclarationWithNonIterableRhs()
func (c *TypeInferenceTest) TestArrayDestructuringDeclarationWithNonIterableRhs() {
	c.InFunction("let [x] = 3; X: x;")
	c.AssertTypeOfExpression("X").ToStringIsEqualTo("?")
}

// Java: TypeInferenceTest.testArrayDestructuringAssignWithGetProp()
func (c *TypeInferenceTest) TestArrayDestructuringAssignWithGetProp() {
	c.InFunction(`const ns = {};
const /** !Iterable<number> */ numbers = [1, 2, 3];
[ns.x] = numbers;
NSX: ns.x;
`)
	c.AssertTypeOfExpression("NSX").ToStringIsEqualTo("number")
}

// Java: TypeInferenceTest.testArrayDestructuringAssignWithGetElem()
func (c *TypeInferenceTest) TestArrayDestructuringAssignWithGetElem() {
	c.InFunction(`const arr = [];
const /** !Iterable<number> */ numbers = [1, 2, 3];
[arr[1]] = numbers;
ARR1: arr[1];
`)
	c.AssertTypeOfExpression("ARR1").ToStringIsEqualTo("?")
}

// Java: TypeInferenceTest.testDeclarationDoesntOverrideInferredTypeInDestructuringPattern()
func (c *TypeInferenceTest) TestDeclarationDoesntOverrideInferredTypeInDestructuringPattern() {
	c.InFunction("var [/** number */ x] = /** @type {?} */ ([null]); X: x")
	c.AssertTypeOfExpression("X").ToStringIsEqualTo("number")
}

// Java: TypeInferenceTest.testDeclarationDoesntOverrideInferredTypeInForOfLoop()
func (c *TypeInferenceTest) TestDeclarationDoesntOverrideInferredTypeInForOfLoop() {
	c.InFunction("for (var /** number */ x of /** @type {?} */ [null]) { X: x; }")
	c.AssertTypeOfExpression("X").ToStringIsEqualTo("number")
}

// Java: TypeInferenceTest.testTypeInferenceOccursInDestructuringCatch()
func (c *TypeInferenceTest) TestTypeInferenceOccursInDestructuringCatch() {
	c.AssumingNative("x", rhino.JSTypeNative_NUMBER_TYPE)
	c.InFunction(`try {
  throw {err: 3};
} catch ({[x = 'err']: /** number */ err}) {
  ERR: err;
  X: x;
}
`)
	c.AssertTypeOfExpression("ERR").ToStringIsEqualTo("number")
	c.AssertTypeOfExpression("X").ToStringIsEqualTo("string")
}

// Java: TypeInferenceTest.testTypeInferenceOccursInDestructuringForIn()
func (c *TypeInferenceTest) TestTypeInferenceOccursInDestructuringForIn() {
	c.AssumingNative("x", rhino.JSTypeNative_NUMBER_TYPE)
	c.InFunction(`/** @type {number} */
String.prototype.length;

var obj = {};
for ({length: obj.length} in {'1': 1, '22': 22}) {
  LENGTH: obj.length; // set to '1'.length and '22'.length
}
`)
	c.AssertTypeOfExpression("LENGTH").ToStringIsEqualTo("number")
}

// Java: TypeInferenceTest.testInferringTypeInObjectPattern_fromTemplatizedProperty()
func (c *TypeInferenceTest) TestInferringTypeInObjectPattern_fromTemplatizedProperty() {
	var templateKey *rhino.TemplateType = c.registry.CreateTemplateType("T")
	var fooCtor rhino.FunctionType = c.registry.CreateConstructorType(jutil.StringPtr("Foo"), rhino.NilNode, c.registry.CreateParameters(), nil, []*rhino.TemplateType{templateKey}, false)
	var fooInstanceType rhino.ObjectType = fooCtor.GetInstanceType()
	fooInstanceType.DefineDeclaredProperty("data", templateKey, rhino.NilNode)
	var fooOfNumber rhino.JSType = c.Templatize(fooInstanceType, []rhino.JSType{c.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)})
	c.Assuming("obj", fooOfNumber)
	c.InFunction(`const {data} = obj;
OBJ: obj;
DATA: data
`)
	c.AssertTypeOfExpression("OBJ").ToStringIsEqualTo("Foo<number>")
	c.AssertTypeOfExpression("DATA").ToStringIsEqualTo("number")
}

// Java: TypeInferenceTest.testTypeInferenceOccursInsideVoidOperator()
func (c *TypeInferenceTest) TestTypeInferenceOccursInsideVoidOperator() {
	c.InFunction("var x; var y = void (x = 3); X: x; Y: y")
	c.AssertTypeOfExpression("X").ToStringIsEqualTo("number")
	c.AssertTypeOfExpression("Y").ToStringIsEqualTo("undefined")
}

// Java: TypeInferenceTest.constDeclarationWithReturnJSDoc_ignoresUnknownRhsType()
func (c *TypeInferenceTest) Test_constDeclarationWithReturnJSDoc_ignoresUnknownRhsType() {
	c.AssumingNative("foo", rhino.JSTypeNative_UNKNOWN_TYPE)
	c.InFunction("/** @return {number} */ const fn = foo;")
	var fooWithInterfaceType rhino.JSType = c.GetType("fn")
	TypeSubject_AssertType(fooWithInterfaceType).IsFunctionTypeThat().
		HasReturnTypeThat().
		IsNumber()
}

// Java: TypeInferenceTest.constDeclarationWithCtorJSDoc_ignoresKnownMixinReturnType()
func (c *TypeInferenceTest) Test_constDeclarationWithCtorJSDoc_ignoresKnownMixinReturnType() {
	var fooType rhino.JSType = rhino.FunctionType_Builder_(c.registry).ForConstructor().
		WithName(jutil.StringPtr("Foo")).
		Build()
	c.Assuming("Foo", fooType)
	var mixinType rhino.FunctionType = rhino.FunctionType_Builder_(c.registry).WithReturnType(fooType).
		Build()
	c.Assuming("mixin", mixinType)
	c.InFunction("/** @constructor @extends {Foo} */ const FooExtended = mixin();")
	var fooWithInterfaceType rhino.JSType = c.GetType("FooExtended")
	TypeSubject_AssertType(fooWithInterfaceType).IsNotEqualTo(fooType)
	TypeSubject_AssertType(fooWithInterfaceType).ToStringIsEqualTo("function(new:FooExtended): ?")
}

// Java: TypeInferenceTest.testSideEffectsInEsExportDefaultInferred()
func (c *TypeInferenceTest) TestSideEffectsInEsExportDefaultInferred() {
	c.AssumingNative("foo", rhino.JSTypeNative_NUMBER_TYPE)
	c.AssumingNative("bar", rhino.JSTypeNative_UNKNOWN_TYPE)
	c.InScript("export default (bar = foo, foo = 'not a number');")
	TypeSubject_AssertType(c.GetType("bar")).IsNumber()
	TypeSubject_AssertType(c.GetType("foo")).IsString()
}

// Java: TypeInferenceTest.testShneTightensUnknownOperandOnLhs()
func (c *TypeInferenceTest) TestShneTightensUnknownOperandOnLhs() {
	c.AssumingNative("foo", rhino.JSTypeNative_NUMBER_TYPE)
	c.AssumingNative("bar", rhino.JSTypeNative_UNKNOWN_TYPE)
	c.InFunction("if (bar === foo) { FOO: foo; BAR: bar; }")
	c.AssertTypeOfExpression("FOO").IsNumber()
	c.AssertTypeOfExpression("BAR").IsNumber()
}

// Java: TypeInferenceTest.testShneTightensUnknownOperandOnRhs()
func (c *TypeInferenceTest) TestShneTightensUnknownOperandOnRhs() {
	c.AssumingNative("foo", rhino.JSTypeNative_NUMBER_TYPE)
	c.AssumingNative("bar", rhino.JSTypeNative_UNKNOWN_TYPE)
	c.InFunction("if (foo === bar) { FOO: foo; BAR: bar; }")
	c.AssertTypeOfExpression("FOO").IsNumber()
	c.AssertTypeOfExpression("BAR").IsNumber()
}

// Java: TypeInferenceTest.testDynamicImport()
func (c *TypeInferenceTest) TestDynamicImport() {
	c.InScript("const foo = import('foo');")
	var promiseOfUnknownType rhino.JSType = c.registry.CreateTemplatizedType(c.registry.GetNativeObjectType(rhino.JSTypeNative_PROMISE_TYPE), c.registry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE))
	TypeSubject_AssertType(c.GetType("foo")).IsSubtypeOf(promiseOfUnknownType)
}

// Java: TypeInferenceTest.testDynamicImport2()
func (c *TypeInferenceTest) TestDynamicImport2() {
	c.WithModules([]string{"export default 1; export /** @return {string} */ function Bar() { return 'bar'; };", "const foo = import('./testcode0');"})
	TypeSubject_AssertType(c.GetType("foo")).ToStringIsEqualTo(`Promise<{
  Bar: function(): string,
  default: number
}>`)
}

// Java: TypeInferenceTest.testRequireDynamic()
func (c *TypeInferenceTest) TestRequireDynamic() {
	c.WithModules([]string{"goog.module('foo'); exports.barFunc = function Bar() { return 'bar'; };", "const f = goog.requireDynamic('foo');", "const e = goog.require('foo');"})
	TypeSubject_AssertType(c.GetType("f")).ToStringIsEqualTo("IThenable<{barFunc: function(): ?}>")
	TypeSubject_AssertType(c.GetType("e")).ToStringIsEqualTo("{barFunc: function(): ?}")
}

// Java: TypeInferenceTest.testDynamicImportAfterModuleRewriting()
func (c *TypeInferenceTest) TestDynamicImportAfterModuleRewriting() {
	c.WithModules([]string{`const module$testcode0 = {};
/** @const */ module$testcode0.default = 1;
/** @return {string} */ function Bar() { return 'bar'; };
/** @const */ module$testcode0.Bar = Bar;
`, "const foo = import('./testcode0');"})
	TypeSubject_AssertType(c.GetType("foo")).ToStringIsEqualTo(`Promise<{
  Bar: function(): string,
  default: number
}>`)
}

// Java: TypeInferenceTest.testForAllBigInt(JSType)
func TypeInferenceTest_TestForAllBigInt(type_ rhino.JSType) {
	AssertThat(TypeInference_GetBigIntPresence(type_)).IsEqualTo(TypeInference_BigIntPresence_ALL_BIGINT)
}

// Java: TypeInferenceTest.testForNoBigInt(JSType)
func TypeInferenceTest_TestForNoBigInt(type_ rhino.JSType) {
	AssertThat(TypeInference_GetBigIntPresence(type_)).IsEqualTo(TypeInference_BigIntPresence_NO_BIGINT)
}

// Java: TypeInferenceTest.testForBigIntOrNumber(JSType)
func TypeInferenceTest_TestForBigIntOrNumber(type_ rhino.JSType) {
	AssertThat(TypeInference_GetBigIntPresence(type_)).IsEqualTo(TypeInference_BigIntPresence_BIGINT_OR_NUMBER)
}

// Java: TypeInferenceTest.testForBigIntOrOther(JSType)
func TypeInferenceTest_TestForBigIntOrOther(type_ rhino.JSType) {
	AssertThat(TypeInference_GetBigIntPresence(type_)).IsEqualTo(TypeInference_BigIntPresence_BIGINT_OR_OTHER)
}

// Java: TypeInferenceTest.getNativeObjectType(JSTypeNative)
func (c *TypeInferenceTest) GetNativeObjectType(t rhino.JSTypeNative) rhino.ObjectType {
	return c.registry.GetNativeObjectType(t)
}

// Java: TypeInferenceTest.getNativeType(JSTypeNative)
func (c *TypeInferenceTest) GetNativeType(t rhino.JSTypeNative) rhino.JSType {
	return c.registry.GetNativeType(t)
}

// Java: TypeInferenceTest.templatize(ObjectType, ImmutableList<JSType>)
func (c *TypeInferenceTest) Templatize(objType rhino.ObjectType, t []rhino.JSType) rhino.JSType {
	return c.registry.CreateTemplatizedTypeOfList(objType, t)
}

// Java: TypeInferenceTest.includeGoogAssertionFn(String, JSType)
func (c *TypeInferenceTest) IncludeGoogAssertionFn(fnName string, returnType rhino.JSType) {
	fullName := "goog.asserts." + fnName
	var fnType rhino.FunctionType = rhino.FunctionType_Builder_(c.registry).WithReturnType(returnType).
		WithParametersOfList(c.registry.CreateParameters(c.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE))).
		WithName(jutil.StringPtr(fullName)).
		Build()
	c.Assuming(fullName, fnType)
}

// Java: TypeInferenceTest.includePrimitiveTruthyAssertionFunction(String)
func (c *TypeInferenceTest) IncludePrimitiveTruthyAssertionFunction(fnName string) {
	var t *rhino.TemplateType = c.registry.CreateTemplateType("T")
	// Go-only: language-adapter — the nullable ClosurePrimitive is a pointer.
	assertsTruthy := rhino.ClosurePrimitive_ASSERTS_TRUTHY
	var assertType rhino.FunctionType = rhino.FunctionType_Builder_(c.registry).WithName(jutil.StringPtr(fnName)).
		WithClosurePrimitiveId(&assertsTruthy).
		WithReturnType(t).
		WithParametersOfList(c.registry.CreateParameters(t)).
		WithTemplateKeys(t).
		Build()
	c.Assuming(fnName, assertType)
}

// Java: TypeInferenceTest.includePrimitiveAssertionFn(String, JSType)
func (c *TypeInferenceTest) IncludePrimitiveAssertionFn(fullName string, returnType rhino.JSType) {
	// Go-only: language-adapter — the nullable ClosurePrimitive is a pointer.
	assertsMatchesReturn := rhino.ClosurePrimitive_ASSERTS_MATCHES_RETURN
	var fnType rhino.FunctionType = rhino.FunctionType_Builder_(c.registry).WithReturnType(returnType).
		WithParametersOfList(c.registry.CreateParameters(c.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE))).
		WithName(jutil.StringPtr(fullName)).
		WithClosurePrimitiveId(&assertsMatchesReturn).
		Build()
	c.Assuming(fullName, fnType)
}

// Java: TypeInferenceTest.includeAssertInstanceof()
func (c *TypeInferenceTest) IncludeAssertInstanceof() {
	fullName := "goog.asserts.assertInstanceof"
	var templateType *rhino.TemplateType = c.registry.CreateTemplateType("T")
	var templateTypeCtor rhino.FunctionType = rhino.FunctionType_Builder_(c.registry).ForConstructor().
		WithTypeOfThis(templateType).
		Build()
	var fnType rhino.FunctionType = rhino.FunctionType_Builder_(c.registry).WithParametersOfList(c.registry.CreateParameters(c.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE), templateTypeCtor)).
		WithTemplateKeys(templateType).
		WithReturnType(templateType).
		WithName(jutil.StringPtr(fullName)).
		Build()
	c.Assuming(fullName, fnType)
}
