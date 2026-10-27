// Port of com/google/javascript/jscomp/CrossChunkReferenceCollectorTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/CrossChunkReferenceCollectorTest.java)

/*
 * Copyright 2017 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"slices"
	"testing"
)

type CrossChunkReferenceCollectorTest struct {
	CompilerTestCaseBase
	testedCollector *CrossChunkReferenceCollector
}

// Java: CrossChunkReferenceCollectorTest()
func NewCrossChunkReferenceCollectorTest() *CrossChunkReferenceCollectorTest {
	c := &CrossChunkReferenceCollectorTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestCrossChunkReferenceCollectorTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCrossChunkReferenceCollectorTest() })
}

// JUnitBefore runs the @Before methods: CompilerTestCase.setUp() (the
// superclass one) first, then customSetUp().
func (c *CrossChunkReferenceCollectorTest) JUnitBefore() {
	c.CompilerTestCaseBase.JUnitBefore()
	c.CustomSetUp()
}

// Java: CrossChunkReferenceCollectorTest.customSetUp()
func (c *CrossChunkReferenceCollectorTest) CustomSetUp() {
	c.EnableNormalize()
}

// Java: CrossChunkReferenceCollectorTest.getProcessor(Compiler)
func (c *CrossChunkReferenceCollectorTest) GetProcessor(compiler *Compiler) CompilerPass {
	scopeCreator := NewSyntacticScopeCreator(compiler)
	c.testedCollector = NewCrossChunkReferenceCollector(compiler, scopeCreator)
	return c.testedCollector
}

// Java: CrossChunkReferenceCollectorTest.testVarInBlock()
func (c *CrossChunkReferenceCollectorTest) TestVarInBlock() {
	c.TestSame(
		"  if (true) {\n" +
			"    var y = x;\n" +
			"    y;\n" +
			"    y;\n" +
			"  }\n")
	globalVariableNamesMap := c.testedCollector.GetGlobalVariableNamesMap()
	AssertThat(globalVariableNamesMap).ContainsKey("y")
	yVar := globalVariableNamesMap.Get("y")
	yRefs := c.testedCollector.GetReferences(yVar)
	AssertThat(yRefs.IsAssignedOnceInLifetime()).IsTrue()
	AssertThat(yRefs.IsWellDefined()).IsTrue()
}

// Java: CrossChunkReferenceCollectorTest.testVarInLoopNotAssignedOnlyOnceInLifetime()
func (c *CrossChunkReferenceCollectorTest) TestVarInLoopNotAssignedOnlyOnceInLifetime() {
	c.TestSame("var x; for (; true;) { x = 0; }")
	globalVariableNamesMap := c.testedCollector.GetGlobalVariableNamesMap()
	xVar := globalVariableNamesMap.Get("x")
	AssertThat(globalVariableNamesMap).ContainsKey("x")
	xRefs := c.testedCollector.GetReferences(xVar)
	AssertThat(xRefs.IsAssignedOnceInLifetime()).IsFalse()

	c.TestSame("let x; for (; true;) { x = 0; }")
	globalVariableNamesMap = c.testedCollector.GetGlobalVariableNamesMap()
	xVar = globalVariableNamesMap.Get("x")
	AssertThat(globalVariableNamesMap).ContainsKey("x")
	xRefs = c.testedCollector.GetReferences(xVar)
	AssertThat(xRefs.IsAssignedOnceInLifetime()).IsFalse()
}

/**
 * Although there is only one assignment to x in the code, it's in a function which could be
 * called multiple times, so {@code isAssignedOnceInLifetime()} returns false.
 */
// Java: CrossChunkReferenceCollectorTest.testVarInFunctionNotAssignedOnlyOnceInLifetime()
func (c *CrossChunkReferenceCollectorTest) TestVarInFunctionNotAssignedOnlyOnceInLifetime() {
	c.TestSame("var x; function f() { x = 0; }")
	globalVariableNamesMap := c.testedCollector.GetGlobalVariableNamesMap()
	xVar := globalVariableNamesMap.Get("x")
	AssertThat(globalVariableNamesMap).ContainsKey("x")
	xRefs := c.testedCollector.GetReferences(xVar)
	AssertThat(xRefs.IsAssignedOnceInLifetime()).IsFalse()

	c.TestSame("let x; function f() { x = 0; }")
	globalVariableNamesMap = c.testedCollector.GetGlobalVariableNamesMap()
	xVar = globalVariableNamesMap.Get("x")
	AssertThat(globalVariableNamesMap).ContainsKey("x")
	xRefs = c.testedCollector.GetReferences(xVar)
	AssertThat(xRefs.IsAssignedOnceInLifetime()).IsFalse()
}

// Java: CrossChunkReferenceCollectorTest.testVarAssignedOnceInLifetime1()
func (c *CrossChunkReferenceCollectorTest) TestVarAssignedOnceInLifetime1() {
	c.TestSame("var x = 0;")
	globalVariableNamesMap := c.testedCollector.GetGlobalVariableNamesMap()
	xVar := globalVariableNamesMap.Get("x")
	AssertThat(globalVariableNamesMap).ContainsKey("x")
	xRefs := c.testedCollector.GetReferences(xVar)
	AssertThat(xRefs.IsAssignedOnceInLifetime()).IsTrue()

	c.TestSame("let x = 0;")
	globalVariableNamesMap = c.testedCollector.GetGlobalVariableNamesMap()
	xVar = globalVariableNamesMap.Get("x")
	AssertThat(globalVariableNamesMap).ContainsKey("x")
	xRefs = c.testedCollector.GetReferences(xVar)
	AssertThat(xRefs.IsAssignedOnceInLifetime()).IsTrue()
}

// Java: CrossChunkReferenceCollectorTest.testVarAssignedOnceInLifetime2()
func (c *CrossChunkReferenceCollectorTest) TestVarAssignedOnceInLifetime2() {
	c.TestSame("{ var x = 0; }")
	globalVariableNamesMap := c.testedCollector.GetGlobalVariableNamesMap()
	xVar := globalVariableNamesMap.Get("x")
	AssertThat(globalVariableNamesMap).ContainsKey("x")
	xRefs := c.testedCollector.GetReferences(xVar)
	AssertThat(xRefs.IsAssignedOnceInLifetime()).IsTrue()
}

// Java: CrossChunkReferenceCollectorTest.testBasicBlocks()
func (c *CrossChunkReferenceCollectorTest) TestBasicBlocks() {
	c.TestSame(
		"var x = 0;\n" +
			"switch (x) {\n" +
			"  case 0:\n" +
			"    x;\n" +
			"}\n")
	globalVariableNamesMap := c.testedCollector.GetGlobalVariableNamesMap()
	xVar := globalVariableNamesMap.Get("x")
	AssertThat(globalVariableNamesMap).ContainsKey("x")
	xRefs := c.testedCollector.GetReferences(xVar)
	AssertThat(xRefs.references).HasSize(3)
	AssertNode(xRefs.references[0].GetBasicBlock().GetRoot()).HasType(rhino.Token_ROOT)
	AssertNode(xRefs.references[1].GetBasicBlock().GetRoot()).HasType(rhino.Token_ROOT)
	AssertNode(xRefs.references[2].GetBasicBlock().GetRoot()).HasType(rhino.Token_CASE)
}

// Java: CrossChunkReferenceCollectorTest.testClassStaticBlock()
func (c *CrossChunkReferenceCollectorTest) TestClassStaticBlock() {
	c.TestSame(
		"var x=1;\n" +
			"class C {\n" +
			"  static {\n" +
			"    x;\n" +
			"    x=2;\n" +
			"  }\n" +
			"}\n")

	globalVariableNamesMap := c.testedCollector.GetGlobalVariableNamesMap()
	AssertThat(globalVariableNamesMap).ContainsKey("x")
	xVar := globalVariableNamesMap.Get("x")
	xRefs := c.testedCollector.GetReferences(xVar)
	AssertThat(xRefs.IsAssignedOnceInLifetime()).IsFalse()
	AssertThat(xRefs.references).HasSize(3)
	AssertNode(xRefs.references[0].GetBasicBlock().GetRoot()).HasType(rhino.Token_ROOT)
	AssertNode(xRefs.references[1].GetBasicBlock().GetRoot()).HasType(rhino.Token_CLASS_MEMBERS)
	AssertNode(xRefs.references[2].GetBasicBlock().GetRoot()).HasType(rhino.Token_CLASS_MEMBERS)
}

// Java: CrossChunkReferenceCollectorTest.nullishCoalesce()
func (c *CrossChunkReferenceCollectorTest) TestNullishCoalesce() {
	c.TestSame("var x = 0; var y = x ?? (x = 1)")
	globalVariableNamesMap := c.testedCollector.GetGlobalVariableNamesMap()
	xVar := globalVariableNamesMap.Get("x")
	AssertThat(globalVariableNamesMap).ContainsKey("x")
	xRefs := c.testedCollector.GetReferences(xVar)
	AssertThat(xRefs.references).HasSize(3)
	AssertNode(xRefs.references[0].GetBasicBlock().GetRoot()).HasType(rhino.Token_ROOT)
	AssertNode(xRefs.references[1].GetBasicBlock().GetRoot()).HasType(rhino.Token_ROOT)
	AssertNode(xRefs.references[2].GetBasicBlock().GetRoot()).HasType(rhino.Token_ASSIGN)
}

// Java: CrossChunkReferenceCollectorTest.testTopLevelStatements()
func (c *CrossChunkReferenceCollectorTest) TestTopLevelStatements() {
	c.TestSame(
		"var x = 1;\n" +
			"const y = x;\n" +
			"let z = x - y;\n" +
			"function f(x1, y1) { // only f and z globals referenced\n" +
			"// NOTE: If we try to name the parameters the same as the global variables, then\n" +
			"// Normalization will rename them. Normalization always runs before\n" +
			"// CrossChunkReferenceCollector does.\n" +
			"  return x1 + y1 + z;\n" +
			"}\n")

	// Pull out all the references for comparison.
	globalVariableNamesMap := c.testedCollector.GetGlobalVariableNamesMap()
	xReferences :=
		slices.Collect(c.testedCollector.GetReferences(globalVariableNamesMap.Get("x")).All())
	yReferences :=
		slices.Collect(c.testedCollector.GetReferences(globalVariableNamesMap.Get("y")).All())
	zReferences :=
		slices.Collect(c.testedCollector.GetReferences(globalVariableNamesMap.Get("z")).All())
	fReferences :=
		slices.Collect(c.testedCollector.GetReferences(globalVariableNamesMap.Get("f")).All())

	// Make sure the statements have the references we expect.
	topLevelStatements := c.testedCollector.GetTopLevelStatements()
	AssertThat(topLevelStatements).HasSize(4)
	// var x = 1;
	xEquals1 := topLevelStatements[0]
	AssertThat(xEquals1.GetOriginalOrder()).IsEqualTo(0)
	AssertThat(xEquals1.GetDeclaredNameReference()).IsEqualTo(xReferences[0])
	AssertThat(xEquals1.GetNonDeclarationReferences()).IsEmpty()
	// const y = x;
	yEqualsX := topLevelStatements[1]
	AssertThat(yEqualsX.GetOriginalOrder()).IsEqualTo(1)
	AssertThat(yEqualsX.GetNonDeclarationReferences()).ContainsExactly(xReferences[1])
	// let z = x - y;
	zEqualsXMinusY := topLevelStatements[2]
	AssertThat(zEqualsXMinusY.GetOriginalOrder()).IsEqualTo(2)
	AssertThat(zEqualsXMinusY.GetNonDeclarationReferences()).
		ContainsExactly(xReferences[2], yReferences[1])
	// function f(x, y) { return x + y + z; }
	functionDeclaration := topLevelStatements[3]
	AssertThat(functionDeclaration.GetOriginalOrder()).IsEqualTo(3)
	AssertThat(functionDeclaration.GetDeclaredNameReference()).IsEqualTo(fReferences[0])
	AssertThat(functionDeclaration.GetNonDeclarationReferences()).
		ContainsExactly(zReferences[1])
}

// Java: CrossChunkReferenceCollectorTest.testVarDeclarationStatement()
func (c *CrossChunkReferenceCollectorTest) TestVarDeclarationStatement() {
	c.TestSame("var x = 1;")

	statements := c.testedCollector.GetTopLevelStatements()
	xRefs := c.GetReferencesForName("x", c.testedCollector)
	AssertThat(statements).HasSize(1)
	varStatement := statements[0]
	declaredNameReference := varStatement.GetDeclaredNameReference()
	AssertThat(declaredNameReference).IsEqualTo(xRefs.references[0])
	AssertThat(varStatement.GetNonDeclarationReferences()).IsEmpty()
	valueNode := varStatement.GetDeclaredValueNode()
	AssertThat(valueNode).IsNotNull()
	AssertThat(valueNode.GetDouble()).IsEqualTo(1.0)
}

// Java: CrossChunkReferenceCollectorTest.testFunctionDeclarationStatement()
func (c *CrossChunkReferenceCollectorTest) TestFunctionDeclarationStatement() {
	c.TestSame("function x() {}")

	statements := c.testedCollector.GetTopLevelStatements()
	AssertThat(statements).HasSize(1)
	xRefs := c.GetReferencesForName("x", c.testedCollector)

	functionDeclaration := statements[0]
	declaredNameReference := functionDeclaration.GetDeclaredNameReference()
	AssertThat(declaredNameReference).IsEqualTo(xRefs.references[0])
	AssertThat(functionDeclaration.GetNonDeclarationReferences()).IsEmpty()
}

// Java: CrossChunkReferenceCollectorTest.testVariableAssignmentStatement()
func (c *CrossChunkReferenceCollectorTest) TestVariableAssignmentStatement() {
	c.TestSame("var x; x = 1;")

	statements := c.testedCollector.GetTopLevelStatements()
	AssertThat(statements).HasSize(2)
	xRefs := c.GetReferencesForName("x", c.testedCollector)

	assignmentStatement := statements[1]
	declaredNameReference := assignmentStatement.GetDeclaredNameReference()
	AssertThat(declaredNameReference).IsEqualTo(xRefs.references[1])
	AssertThat(assignmentStatement.GetNonDeclarationReferences()).IsEmpty()
	valueNode := assignmentStatement.GetDeclaredValueNode()
	AssertThat(valueNode).IsNotNull()
	AssertThat(valueNode.GetDouble()).IsEqualTo(1.0)
}

// Java: CrossChunkReferenceCollectorTest.testPropertyAssignmentStatement()
func (c *CrossChunkReferenceCollectorTest) TestPropertyAssignmentStatement() {
	c.TestSame("var x = {}; x.prop = 1;")

	statements := c.testedCollector.GetTopLevelStatements()
	AssertThat(statements).HasSize(2)
	xRefs := c.GetReferencesForName("x", c.testedCollector)

	assignmentStatement := statements[1]
	declaredNameReference := assignmentStatement.GetDeclaredNameReference()
	AssertThat(declaredNameReference).IsEqualTo(xRefs.references[1])
	AssertThat(assignmentStatement.GetNonDeclarationReferences()).IsEmpty()
	valueNode := assignmentStatement.GetDeclaredValueNode()
	AssertThat(valueNode).IsNotNull()
	AssertThat(valueNode.GetDouble()).IsEqualTo(1.0)
}

// Java: CrossChunkReferenceCollectorTest.testGoogInheritsIsMovableDeclaration()
func (c *CrossChunkReferenceCollectorTest) TestGoogInheritsIsMovableDeclaration() {
	c.TestSame("function A() {} function B() {} goog.inherits(B, A);")

	statements := c.testedCollector.GetTopLevelStatements()
	AssertThat(statements).HasSize(3)

	refsToA := c.GetReferencesForName("A", c.testedCollector)
	refsToB := c.GetReferencesForName("B", c.testedCollector)

	inheritsStatement := statements[2]
	declaredNameReference := inheritsStatement.GetDeclaredNameReference()
	AssertThat(declaredNameReference).IsEqualTo(refsToB.references[1])
	AssertThat(inheritsStatement.GetNonDeclarationReferences()).
		ContainsExactly(refsToA.references[1])
	// inherits statements are always movable
	AssertThat(inheritsStatement.IsMovableDeclaration()).IsTrue()
}

// Java: CrossChunkReferenceCollectorTest.testDefinePropertiesIsMovableDeclaration()
func (c *CrossChunkReferenceCollectorTest) TestDefinePropertiesIsMovableDeclaration() {
	c.TestSame("function A() {} Object.defineProperties(A, {});")

	statements := c.testedCollector.GetTopLevelStatements()
	AssertThat(statements).HasSize(2)

	definePropertiesStatement := statements[1]
	declaredNameReference := definePropertiesStatement.GetDeclaredNameReference()
	AssertThat(declaredNameReference).IsNotNull()
	AssertThat(definePropertiesStatement.GetNonDeclarationReferences()).IsEmpty()
	// defineProperties statements are always movable
	AssertThat(definePropertiesStatement.IsMovableDeclaration()).IsTrue()
}

// Java: CrossChunkReferenceCollectorTest.testDefinePropertiesWithPrototypeIsMovableDeclaration()
func (c *CrossChunkReferenceCollectorTest) TestDefinePropertiesWithPrototypeIsMovableDeclaration() {
	c.TestSame("function A() {} Object.defineProperties(A.prototype, {});")

	statements := c.testedCollector.GetTopLevelStatements()
	AssertThat(statements).HasSize(2)

	definePropertiesStatement := statements[1]
	declaredNameReference := definePropertiesStatement.GetDeclaredNameReference()
	AssertThat(declaredNameReference).IsNotNull()
	AssertThat(definePropertiesStatement.GetNonDeclarationReferences()).IsEmpty()
	// defineProperties statements are always movable
	AssertThat(definePropertiesStatement.IsMovableDeclaration()).IsTrue()
}

// Java: CrossChunkReferenceCollectorTest.testFunctionDeclarationOrAssignmentIsMovable()
func (c *CrossChunkReferenceCollectorTest) TestFunctionDeclarationOrAssignmentIsMovable() {
	c.TestSame("function f() {}")
	AssertThat(c.testedCollector.GetTopLevelStatements()[0].IsMovableDeclaration()).IsTrue()
	c.TestSame("var f = function() {};")
	AssertThat(c.testedCollector.GetTopLevelStatements()[0].IsMovableDeclaration()).IsTrue()
}

// Java: CrossChunkReferenceCollectorTest.testLiteralValueIsMovable()
func (c *CrossChunkReferenceCollectorTest) TestLiteralValueIsMovable() {
	c.TestSame("var f = 1;")
	AssertThat(c.testedCollector.GetTopLevelStatements()[0].IsMovableDeclaration()).IsTrue()
}

// Java: CrossChunkReferenceCollectorTest.testFunctionCallsAreNotMovableExceptForMethodStubs()
func (c *CrossChunkReferenceCollectorTest) TestFunctionCallsAreNotMovableExceptForMethodStubs() {
	c.TestSame(
		"function Foo() {}\n" +
			"Foo.prototype.stub = JSCompiler_stubMethod(x);\n" +
			"Foo.prototype.unstub = JSCompiler_unstubMethod(x);\n" +
			"Foo.prototype.other = other();\n")
	statements := c.testedCollector.GetTopLevelStatements()
	AssertThat(statements[1].IsMovableDeclaration()).IsTrue()
	AssertThat(statements[2].IsMovableDeclaration()).IsFalse()
	AssertThat(statements[3].IsMovableDeclaration()).IsFalse()
}

// Java: CrossChunkReferenceCollectorTest.testUnknownNameValueIsImmovable()
func (c *CrossChunkReferenceCollectorTest) TestUnknownNameValueIsImmovable() {
	c.AssertStatementIsImmovable("var a = unknownName;")
}

// Java: CrossChunkReferenceCollectorTest.testWellDefinedNameValueIsMovable()
func (c *CrossChunkReferenceCollectorTest) TestWellDefinedNameValueIsMovable() {
	c.TestSame("var wellDefined = 1; var other = wellDefined;")
	AssertThat(c.testedCollector.GetTopLevelStatements()[1].IsMovableDeclaration()).IsTrue()
}

// Java: CrossChunkReferenceCollectorTest.testUninitializedNameValueIsNotMovable()
func (c *CrossChunkReferenceCollectorTest) TestUninitializedNameValueIsNotMovable() {
	c.TestSame("var value; var other = value;")
	AssertThat(c.testedCollector.GetTopLevelStatements()[1].IsMovableDeclaration()).IsFalse()
}

// Java: CrossChunkReferenceCollectorTest.testReDefinedNameValueIsNotMovable()
func (c *CrossChunkReferenceCollectorTest) TestReDefinedNameValueIsNotMovable() {
	c.TestSame("var redefined = 1; redefined = 2; var other = redefined;")
	AssertThat(c.testedCollector.GetTopLevelStatements()[2].IsMovableDeclaration()).IsFalse()
}

// Java: CrossChunkReferenceCollectorTest.testEmptyArrayLiteralIsMovable()
func (c *CrossChunkReferenceCollectorTest) TestEmptyArrayLiteralIsMovable() {
	c.TestSame("var a = [];")
	AssertThat(c.testedCollector.GetTopLevelStatements()[0].IsMovableDeclaration()).IsTrue()
}

// Java: CrossChunkReferenceCollectorTest.testArrayLiteralOfMovablesIsMovable()
func (c *CrossChunkReferenceCollectorTest) TestArrayLiteralOfMovablesIsMovable() {
	c.TestSame("var wellDefinedName = 1; var a = [function(){}, 1, wellDefinedName, []];")
	AssertThat(c.testedCollector.GetTopLevelStatements()[1].IsMovableDeclaration()).IsTrue()
}

// Java: CrossChunkReferenceCollectorTest.testArrayLiteralWithImmovableIsImmovable()
func (c *CrossChunkReferenceCollectorTest) TestArrayLiteralWithImmovableIsImmovable() {
	c.AssertStatementIsImmovable("var a = [unknownValue];")
}

// Java: CrossChunkReferenceCollectorTest.testEmptyObjectLiteralIsMovable()
func (c *CrossChunkReferenceCollectorTest) TestEmptyObjectLiteralIsMovable() {
	c.TestSame("var o = {};")
	AssertThat(c.testedCollector.GetTopLevelStatements()[0].IsMovableDeclaration()).IsTrue()
}

// Java: CrossChunkReferenceCollectorTest.testObjectLiteralOfMovablesIsMovable()
func (c *CrossChunkReferenceCollectorTest) TestObjectLiteralOfMovablesIsMovable() {
	c.TestSame(
		"var wellDefinedName = 1;\n" +
			"var o = {\n" +
			"  f: function(){},\n" +
			"  one: 1,\n" +
			"  n: wellDefinedName,\n" +
			"  o: {},\n" +
			"  'quoted': 1,\n" +
			"  123: 2,\n" +
			"// computed\n" +
			"  ['computed string']: 1,\n" +
			"  [234]: 1,\n" +
			"// method shorthand\n" +
			"  method() {},\n" +
			"  'quoted method'() {},\n" +
			"  ['computed method']() {},\n" +
			"  [345]() {},\n" +
			"// variable shorthand\n" +
			"  wellDefinedName,\n" +
			"// getters\n" +
			"  get x() {},\n" +
			"  get 'a'() {},\n" +
			"  get ['a']() {},\n" +
			"  get 456() {},\n" +
			"  get [567]() {},\n" +
			"// setters\n" +
			"  set x(x) {},\n" +
			"  set 'a'(v) {},\n" +
			"  set ['a'](v1) {},\n" +
			"  set 678(v2) {},\n" +
			"  set [678](v3) {},\n" +
			"// spread\n" +
			"  ...wellDefinedName,\n" +
			"};\n")
	AssertThat(c.testedCollector.GetTopLevelStatements()[1].IsMovableDeclaration()).IsTrue()
}

// Java: CrossChunkReferenceCollectorTest.testObjectLiteralWithImmovableIsImmovable()
func (c *CrossChunkReferenceCollectorTest) TestObjectLiteralWithImmovableIsImmovable() {
	c.AssertStatementIsImmovable("var o = { v: unknownValue };")
	c.AssertStatementIsImmovable("var o = { ...unknownValue };")
	c.AssertStatementIsImmovable("var o = { [unknownValue]: 1 };")
	c.AssertStatementIsImmovable("var o = { [unknownValue]() {} };")
	c.AssertStatementIsImmovable("var o = { get [unknownValue]() {} };")
	c.AssertStatementIsImmovable("var o = { set [unknownValue](x) {} };")
}

// Java: CrossChunkReferenceCollectorTest.assertStatementIsImmovable(String)
func (c *CrossChunkReferenceCollectorTest) AssertStatementIsImmovable(statement string) {
	c.TestSame(statement)
	AssertThat(c.testedCollector.GetTopLevelStatements()[0].IsMovableDeclaration()).IsFalse()
}

// Java: CrossChunkReferenceCollectorTest.testTemplateLiteralIsMovableIfSubstitutionsAreMovable()
func (c *CrossChunkReferenceCollectorTest) TestTemplateLiteralIsMovableIfSubstitutionsAreMovable() {
	c.TestSame(
		"var wellDefinedName = 1;\n" +
			"var t = `${wellDefinedName}`;\n")
	AssertThat(c.testedCollector.GetTopLevelStatements()[1].IsMovableDeclaration()).IsTrue()
}

// Java: CrossChunkReferenceCollectorTest.testTemplateLiteralIsImmovableIfSubstitutionsAreImmovable()
func (c *CrossChunkReferenceCollectorTest) TestTemplateLiteralIsImmovableIfSubstitutionsAreImmovable() {
	c.AssertStatementIsImmovable("var t = `${unknownValue}`")
}

// Java: CrossChunkReferenceCollectorTest.testPureOrBreakMyCodeAnnotatedStringConcatIsMovable()
func (c *CrossChunkReferenceCollectorTest) TestPureOrBreakMyCodeAnnotatedStringConcatIsMovable() {
	c.TestSame(
		"/** @pureOrBreakMyCode */\n" +
			"var t = 'a' + 'b';\n")
	AssertThat(c.testedCollector.GetTopLevelStatements()[0].IsMovableDeclaration()).IsTrue()
}

// Java: CrossChunkReferenceCollectorTest.testPureOrBreakMyCodeAnnotatedFunctionCallIsMovable()
func (c *CrossChunkReferenceCollectorTest) TestPureOrBreakMyCodeAnnotatedFunctionCallIsMovable() {
	c.TestSame(
		"/** @pureOrBreakMyCode */\n" +
			"var t = someFn();\n")
	AssertThat(c.testedCollector.GetTopLevelStatements()[0].IsMovableDeclaration()).IsTrue()
}

// Java: CrossChunkReferenceCollectorTest.testPureOrBreakMyCodeAnnotatedStaticClassPropertyInitializerIsMovable()
func (c *CrossChunkReferenceCollectorTest) TestPureOrBreakMyCodeAnnotatedStaticClassPropertyInitializerIsMovable() {
	c.TestSame(
		"class SomeClass {}\n" +
			"SomeClass.staticProp = /** @pureOrBreakMyCode */ someFn();\n")
	AssertThat(c.testedCollector.GetTopLevelStatements()[0].IsMovableDeclaration()).IsTrue()
	AssertThat(c.testedCollector.GetTopLevelStatements()[1].IsMovableDeclaration()).IsTrue()
}

// Java: CrossChunkReferenceCollectorTest.testPureOrBreakMyCodeAnnotatedStaticClassPropertyInitializerIsMovable2()
func (c *CrossChunkReferenceCollectorTest) TestPureOrBreakMyCodeAnnotatedStaticClassPropertyInitializerIsMovable2() {
	c.TestSame(
		"class SomeClass {}\n" +
			"SomeClass.staticProp = /** @pureOrBreakMyCode */ { prop: someFn() ? 'a' : 'b' };\n")
	AssertThat(c.testedCollector.GetTopLevelStatements()[0].IsMovableDeclaration()).IsTrue()
	AssertThat(c.testedCollector.GetTopLevelStatements()[1].IsMovableDeclaration()).IsTrue()
}

// Java: CrossChunkReferenceCollectorTest.testPureOrBreakMyCodeAnnotatedNonstaticClassFieldIsMovable()
func (c *CrossChunkReferenceCollectorTest) TestPureOrBreakMyCodeAnnotatedNonstaticClassFieldIsMovable() {
	c.TestSame(
		"class SomeClass {\n" +
			"  a;\n" +
			"  static b = 2;\n" +
			"  ['c'] = 3;\n" +
			"  static 'd' = 'hi';\n" +
			"  1 = 2;\n" +
			"}\n")
	AssertThat(c.testedCollector.GetTopLevelStatements()[0].IsMovableDeclaration()).IsTrue()
}

// Java: CrossChunkReferenceCollectorTest.testPureOrBreakMyCodeAnnotatedStaticClassFieldIsMovable()
func (c *CrossChunkReferenceCollectorTest) TestPureOrBreakMyCodeAnnotatedStaticClassFieldIsMovable() {
	c.TestSame(
		"class SomeClass {\n" +
			"  static staticProp = /** @pureOrBreakMyCode */ someFn();\n" +
			"}\n")
	AssertThat(c.testedCollector.GetTopLevelStatements()[0].IsMovableDeclaration()).IsTrue()
}

//	try to find cases to copy from CrossChunkCodeMotion
//
// Java: CrossChunkReferenceCollectorTest.getReferencesForName(String, CrossChunkReferenceCollector)
func (c *CrossChunkReferenceCollectorTest) GetReferencesForName(name string, collector *CrossChunkReferenceCollector) *ReferenceCollection {
	v := collector.GetGlobalVariableNamesMap().Get(name)
	AssertThat(v).IsNotNull()
	return collector.GetReferences(v)
}
