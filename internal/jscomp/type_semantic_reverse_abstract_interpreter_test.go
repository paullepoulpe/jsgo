// Port of com/google/javascript/jscomp/SemanticReverseAbstractInterpreterTest.java @ closure-compiler v20261021

/*
 * Copyright 2007 The Closure Compiler Authors.
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

// Java: public final class SemanticReverseAbstractInterpreterTest extends CompilerTypeTestCase
type SemanticReverseAbstractInterpreterTest struct {
	CompilerTypeTestCaseBase

	interpreter   ReverseAbstractInterpreter
	functionScope *TypedScope
}

// Java: SemanticReverseAbstractInterpreterTest()
func NewSemanticReverseAbstractInterpreterTest() *SemanticReverseAbstractInterpreterTest {
	c := &SemanticReverseAbstractInterpreterTest{}
	c.InitCompilerTypeTestCaseBase(c)
	return c
}

func TestSemanticReverseAbstractInterpreterTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewSemanticReverseAbstractInterpreterTest() })
}

// Java: SemanticReverseAbstractInterpreterTest.setUp()
func (c *SemanticReverseAbstractInterpreterTest) SetUp() {
	c.CompilerTypeTestCaseBase.SetUp()

	c.interpreter = NewSemanticReverseAbstractInterpreter(c.registry)
}

// Java: SemanticReverseAbstractInterpreterTest.newScope()
func (c *SemanticReverseAbstractInterpreterTest) newScope() []FlowScope {
	globalScope := TypedScope_CreateGlobalScope(rhino.NewNode(rhino.Token_ROOT))
	c.functionScope = NewTypedScope(globalScope, rhino.NewNode(rhino.Token_FUNCTION))
	return []FlowScope{LinkedFlowScope_CreateEntryLattice(c.compiler, c.functionScope)}
}

// Tests reverse interpretation of a NAME expression.
// Java: SemanticReverseAbstractInterpreterTest.testNameCondition()
func (c *SemanticReverseAbstractInterpreterTest) TestNameCondition() {
	blind := c.newScope()
	condition := c.createVar(blind, "a", c.CreateNullableType(c.GetNativeStringType()))

	// true outcome.
	informedTrue := c.interpreter.GetPreciserScopeKnowingConditionOutcome(condition, blind[0], rhino.Outcome_TRUE)
	c.AssertTypeEquals(c.GetNativeStringType(), c.getVarType(informedTrue, "a"))

	// false outcome.
	informedFalse := c.interpreter.GetPreciserScopeKnowingConditionOutcome(condition, blind[0], rhino.Outcome_FALSE)
	c.AssertTypeEquals(c.CreateNullableType(c.GetNativeStringType()), c.getVarType(informedFalse, "a"))
}

// Tests reverse interpretation of a NOT(NAME) expression.
// Java: SemanticReverseAbstractInterpreterTest.testNegatedNameCondition()
func (c *SemanticReverseAbstractInterpreterTest) TestNegatedNameCondition() {
	blind := c.newScope()
	a := c.createVar(blind, "a", c.CreateNullableType(c.GetNativeStringType()))
	condition := rhino.NewNode(rhino.Token_NOT)
	condition.AddChildToBack(a)

	// true outcome.
	informedTrue := c.interpreter.GetPreciserScopeKnowingConditionOutcome(condition, blind[0], rhino.Outcome_TRUE)
	c.AssertTypeEquals(c.CreateNullableType(c.GetNativeStringType()), c.getVarType(informedTrue, "a"))

	// false outcome.
	informedFalse := c.interpreter.GetPreciserScopeKnowingConditionOutcome(condition, blind[0], rhino.Outcome_FALSE)
	c.AssertTypeEquals(c.GetNativeStringType(), c.getVarType(informedFalse, "a"))
}

// Tests reverse interpretation of a ASSIGN expression.
// Java: SemanticReverseAbstractInterpreterTest.testAssignCondition1()
func (c *SemanticReverseAbstractInterpreterTest) TestAssignCondition1() {
	blind := c.newScope()
	c.testBinop(
		blind,
		rhino.Token_ASSIGN,
		c.createVar(blind, "a", c.CreateNullableType(c.GetNativeObjectType())),
		c.createVar(blind, "b", c.CreateNullableType(c.GetNativeObjectType())),
		[]*semanticReverseAbstractInterpreterTest_TypedName{
			newSemanticReverseAbstractInterpreterTest_TypedName("a", c.GetNativeObjectType()), newSemanticReverseAbstractInterpreterTest_TypedName("b", c.GetNativeObjectType())},
		[]*semanticReverseAbstractInterpreterTest_TypedName{
			newSemanticReverseAbstractInterpreterTest_TypedName("a", c.GetNativeNullType()), newSemanticReverseAbstractInterpreterTest_TypedName("b", c.GetNativeNullType())})
}

// Tests reverse interpretation of a SHEQ(NAME, NUMBER) expression.
// Java: SemanticReverseAbstractInterpreterTest.testSheqCondition1()
func (c *SemanticReverseAbstractInterpreterTest) TestSheqCondition1() {
	blind := c.newScope()
	c.testBinop(
		blind,
		rhino.Token_SHEQ,
		c.createVar(blind, "a", c.CreateUnionType(c.GetNativeStringType(), c.GetNativeNumberType())),
		c.createNumber(56),
		[]*semanticReverseAbstractInterpreterTest_TypedName{newSemanticReverseAbstractInterpreterTest_TypedName("a", c.GetNativeNumberType())},
		[]*semanticReverseAbstractInterpreterTest_TypedName{
			newSemanticReverseAbstractInterpreterTest_TypedName("a", c.CreateUnionType(c.GetNativeStringType(), c.GetNativeNumberType()))})
}

// Tests reverse interpretation of a SHEQ(NUMBER, NAME) expression.
// Java: SemanticReverseAbstractInterpreterTest.testSheqCondition2()
func (c *SemanticReverseAbstractInterpreterTest) TestSheqCondition2() {
	blind := c.newScope()
	c.testBinop(
		blind,
		rhino.Token_SHEQ,
		c.createNumber(56),
		c.createVar(blind, "a", c.CreateUnionType(c.GetNativeStringType(), c.GetNativeNumberType())),
		[]*semanticReverseAbstractInterpreterTest_TypedName{newSemanticReverseAbstractInterpreterTest_TypedName("a", c.GetNativeNumberType())},
		[]*semanticReverseAbstractInterpreterTest_TypedName{
			newSemanticReverseAbstractInterpreterTest_TypedName("a", c.CreateUnionType(c.GetNativeStringType(), c.GetNativeNumberType()))})
}

// Tests reverse interpretation of a SHEQ(NAME, NAME) expression.
// Java: SemanticReverseAbstractInterpreterTest.testSheqCondition3()
func (c *SemanticReverseAbstractInterpreterTest) TestSheqCondition3() {
	blind := c.newScope()
	c.testBinop(
		blind,
		rhino.Token_SHEQ,
		c.createVar(blind, "b", c.CreateUnionType(c.GetNativeStringType(), c.GetNativeBooleanType())),
		c.createVar(blind, "a", c.CreateUnionType(c.GetNativeStringType(), c.GetNativeNumberType())),
		[]*semanticReverseAbstractInterpreterTest_TypedName{
			newSemanticReverseAbstractInterpreterTest_TypedName("a", c.GetNativeStringType()), newSemanticReverseAbstractInterpreterTest_TypedName("b", c.GetNativeStringType())},
		[]*semanticReverseAbstractInterpreterTest_TypedName{
			newSemanticReverseAbstractInterpreterTest_TypedName("a", c.CreateUnionType(c.GetNativeStringType(), c.GetNativeNumberType())),
			newSemanticReverseAbstractInterpreterTest_TypedName("b", c.CreateUnionType(c.GetNativeStringType(), c.GetNativeBooleanType()))})
}

// Java: SemanticReverseAbstractInterpreterTest.testSheqCondition4()
func (c *SemanticReverseAbstractInterpreterTest) TestSheqCondition4() {
	blind := c.newScope()
	c.testBinop(
		blind,
		rhino.Token_SHEQ,
		c.createVar(blind, "a", c.CreateUnionType(c.GetNativeStringType(), c.GetNativeVoidType())),
		c.createVar(blind, "b", c.CreateUnionType(c.GetNativeVoidType())),
		[]*semanticReverseAbstractInterpreterTest_TypedName{
			newSemanticReverseAbstractInterpreterTest_TypedName("a", c.GetNativeVoidType()), newSemanticReverseAbstractInterpreterTest_TypedName("b", c.GetNativeVoidType())},
		[]*semanticReverseAbstractInterpreterTest_TypedName{
			newSemanticReverseAbstractInterpreterTest_TypedName("a", c.GetNativeStringType()), newSemanticReverseAbstractInterpreterTest_TypedName("b", c.GetNativeVoidType())})
}

// Java: SemanticReverseAbstractInterpreterTest.testSheqCondition5()
func (c *SemanticReverseAbstractInterpreterTest) TestSheqCondition5() {
	blind := c.newScope()
	c.testBinop(
		blind,
		rhino.Token_SHEQ,
		c.createVar(blind, "a", c.CreateUnionType(c.GetNativeNullType(), c.GetNativeVoidType())),
		c.createVar(blind, "b", c.CreateUnionType(c.GetNativeVoidType())),
		[]*semanticReverseAbstractInterpreterTest_TypedName{
			newSemanticReverseAbstractInterpreterTest_TypedName("a", c.GetNativeVoidType()), newSemanticReverseAbstractInterpreterTest_TypedName("b", c.GetNativeVoidType())},
		[]*semanticReverseAbstractInterpreterTest_TypedName{
			newSemanticReverseAbstractInterpreterTest_TypedName("a", c.GetNativeNullType()), newSemanticReverseAbstractInterpreterTest_TypedName("b", c.GetNativeVoidType())})
}

// Java: SemanticReverseAbstractInterpreterTest.testSheqCondition6()
func (c *SemanticReverseAbstractInterpreterTest) TestSheqCondition6() {
	blind := c.newScope()
	c.testBinop(
		blind,
		rhino.Token_SHEQ,
		c.createVar(blind, "a", c.CreateUnionType(c.GetNativeStringType(), c.GetNativeVoidType())),
		c.createVar(blind, "b", c.CreateUnionType(c.GetNativeNumberType(), c.GetNativeVoidType())),
		[]*semanticReverseAbstractInterpreterTest_TypedName{
			newSemanticReverseAbstractInterpreterTest_TypedName("a", c.GetNativeVoidType()), newSemanticReverseAbstractInterpreterTest_TypedName("b", c.GetNativeVoidType())},
		[]*semanticReverseAbstractInterpreterTest_TypedName{
			newSemanticReverseAbstractInterpreterTest_TypedName("a", c.CreateUnionType(c.GetNativeStringType(), c.GetNativeVoidType())),
			newSemanticReverseAbstractInterpreterTest_TypedName("b", c.CreateUnionType(c.GetNativeNumberType(), c.GetNativeVoidType()))})
}

// Tests reverse interpretation of a SHNE(NAME, NUMBER) expression.
// Java: SemanticReverseAbstractInterpreterTest.testShneCondition1()
func (c *SemanticReverseAbstractInterpreterTest) TestShneCondition1() {
	blind := c.newScope()
	c.testBinop(
		blind,
		rhino.Token_SHNE,
		c.createVar(blind, "a", c.CreateUnionType(c.GetNativeStringType(), c.GetNativeNumberType())),
		c.createNumber(56),
		[]*semanticReverseAbstractInterpreterTest_TypedName{
			newSemanticReverseAbstractInterpreterTest_TypedName("a", c.CreateUnionType(c.GetNativeStringType(), c.GetNativeNumberType()))},
		[]*semanticReverseAbstractInterpreterTest_TypedName{newSemanticReverseAbstractInterpreterTest_TypedName("a", c.GetNativeNumberType())})
}

// Tests reverse interpretation of a SHNE(NUMBER, NAME) expression.
// Java: SemanticReverseAbstractInterpreterTest.testShneCondition2()
func (c *SemanticReverseAbstractInterpreterTest) TestShneCondition2() {
	blind := c.newScope()
	c.testBinop(
		blind,
		rhino.Token_SHNE,
		c.createNumber(56),
		c.createVar(blind, "a", c.CreateUnionType(c.GetNativeStringType(), c.GetNativeNumberType())),
		[]*semanticReverseAbstractInterpreterTest_TypedName{
			newSemanticReverseAbstractInterpreterTest_TypedName("a", c.CreateUnionType(c.GetNativeStringType(), c.GetNativeNumberType()))},
		[]*semanticReverseAbstractInterpreterTest_TypedName{newSemanticReverseAbstractInterpreterTest_TypedName("a", c.GetNativeNumberType())})
}

// Tests reverse interpretation of a SHNE(NAME, NAME) expression.
// Java: SemanticReverseAbstractInterpreterTest.testShneCondition3()
func (c *SemanticReverseAbstractInterpreterTest) TestShneCondition3() {
	blind := c.newScope()
	c.testBinop(
		blind,
		rhino.Token_SHNE,
		c.createVar(blind, "b", c.CreateUnionType(c.GetNativeStringType(), c.GetNativeBooleanType())),
		c.createVar(blind, "a", c.CreateUnionType(c.GetNativeStringType(), c.GetNativeNumberType())),
		[]*semanticReverseAbstractInterpreterTest_TypedName{
			newSemanticReverseAbstractInterpreterTest_TypedName("a", c.CreateUnionType(c.GetNativeStringType(), c.GetNativeNumberType())),
			newSemanticReverseAbstractInterpreterTest_TypedName("b", c.CreateUnionType(c.GetNativeStringType(), c.GetNativeBooleanType()))},
		[]*semanticReverseAbstractInterpreterTest_TypedName{
			newSemanticReverseAbstractInterpreterTest_TypedName("a", c.GetNativeStringType()), newSemanticReverseAbstractInterpreterTest_TypedName("b", c.GetNativeStringType())})
}

// Java: SemanticReverseAbstractInterpreterTest.testShneCondition4()
func (c *SemanticReverseAbstractInterpreterTest) TestShneCondition4() {
	blind := c.newScope()
	c.testBinop(
		blind,
		rhino.Token_SHNE,
		c.createVar(blind, "a", c.CreateUnionType(c.GetNativeStringType(), c.GetNativeVoidType())),
		c.createVar(blind, "b", c.CreateUnionType(c.GetNativeVoidType())),
		[]*semanticReverseAbstractInterpreterTest_TypedName{
			newSemanticReverseAbstractInterpreterTest_TypedName("a", c.GetNativeStringType()), newSemanticReverseAbstractInterpreterTest_TypedName("b", c.GetNativeVoidType())},
		[]*semanticReverseAbstractInterpreterTest_TypedName{
			newSemanticReverseAbstractInterpreterTest_TypedName("a", c.GetNativeVoidType()), newSemanticReverseAbstractInterpreterTest_TypedName("b", c.GetNativeVoidType())})
}

// Java: SemanticReverseAbstractInterpreterTest.testShneCondition5()
func (c *SemanticReverseAbstractInterpreterTest) TestShneCondition5() {
	blind := c.newScope()
	c.testBinop(
		blind,
		rhino.Token_SHNE,
		c.createVar(blind, "a", c.CreateUnionType(c.GetNativeNullType(), c.GetNativeVoidType())),
		c.createVar(blind, "b", c.CreateUnionType(c.GetNativeNullType())),
		[]*semanticReverseAbstractInterpreterTest_TypedName{
			newSemanticReverseAbstractInterpreterTest_TypedName("a", c.GetNativeVoidType()), newSemanticReverseAbstractInterpreterTest_TypedName("b", c.GetNativeNullType())},
		[]*semanticReverseAbstractInterpreterTest_TypedName{
			newSemanticReverseAbstractInterpreterTest_TypedName("a", c.GetNativeNullType()), newSemanticReverseAbstractInterpreterTest_TypedName("b", c.GetNativeNullType())})
}

// Java: SemanticReverseAbstractInterpreterTest.testShneCondition6()
func (c *SemanticReverseAbstractInterpreterTest) TestShneCondition6() {
	blind := c.newScope()
	c.testBinop(
		blind,
		rhino.Token_SHNE,
		c.createVar(blind, "a", c.CreateUnionType(c.GetNativeStringType(), c.GetNativeVoidType())),
		c.createVar(blind, "b", c.CreateUnionType(c.GetNativeNumberType(), c.GetNativeVoidType())),
		[]*semanticReverseAbstractInterpreterTest_TypedName{
			newSemanticReverseAbstractInterpreterTest_TypedName("a", c.CreateUnionType(c.GetNativeStringType(), c.GetNativeVoidType())),
			newSemanticReverseAbstractInterpreterTest_TypedName("b", c.CreateUnionType(c.GetNativeNumberType(), c.GetNativeVoidType()))},
		[]*semanticReverseAbstractInterpreterTest_TypedName{
			newSemanticReverseAbstractInterpreterTest_TypedName("a", c.GetNativeVoidType()), newSemanticReverseAbstractInterpreterTest_TypedName("b", c.GetNativeVoidType())})
}

// Tests reverse interpretation of a EQ(NAME, NULL) expression.
// Java: SemanticReverseAbstractInterpreterTest.testEqCondition1()
func (c *SemanticReverseAbstractInterpreterTest) TestEqCondition1() {
	blind := c.newScope()
	c.testBinop(
		blind,
		rhino.Token_EQ,
		c.createVar(blind, "a", c.CreateUnionType(c.GetNativeBooleanType(), c.GetNativeVoidType())),
		c.createNull(),
		[]*semanticReverseAbstractInterpreterTest_TypedName{newSemanticReverseAbstractInterpreterTest_TypedName("a", c.GetNativeVoidType())},
		[]*semanticReverseAbstractInterpreterTest_TypedName{newSemanticReverseAbstractInterpreterTest_TypedName("a", c.GetNativeBooleanType())})
}

// Tests reverse interpretation of a NE(NULL, NAME) expression.
// Java: SemanticReverseAbstractInterpreterTest.testEqCondition2()
func (c *SemanticReverseAbstractInterpreterTest) TestEqCondition2() {
	blind := c.newScope()
	c.testBinop(
		blind,
		rhino.Token_NE,
		c.createNull(),
		c.createVar(blind, "a", c.CreateUnionType(c.GetNativeBooleanType(), c.GetNativeVoidType())),
		[]*semanticReverseAbstractInterpreterTest_TypedName{newSemanticReverseAbstractInterpreterTest_TypedName("a", c.GetNativeBooleanType())},
		[]*semanticReverseAbstractInterpreterTest_TypedName{newSemanticReverseAbstractInterpreterTest_TypedName("a", c.GetNativeVoidType())})
}

// Tests reverse interpretation of a EQ(NAME, NULL) expression.
// Java: SemanticReverseAbstractInterpreterTest.testEqCondition3()
func (c *SemanticReverseAbstractInterpreterTest) TestEqCondition3() {
	blind := c.newScope()
	// (number,undefined,null)
	nullableOptionalNumber := c.CreateUnionType(c.GetNativeNullType(), c.GetNativeVoidType(), c.GetNativeNumberType())
	// (null,undefined)
	nullUndefined := c.CreateUnionType(c.GetNativeVoidType(), c.GetNativeNullType())
	c.testBinop(
		blind,
		rhino.Token_EQ,
		c.createVar(blind, "a", nullableOptionalNumber),
		c.createNull(),
		[]*semanticReverseAbstractInterpreterTest_TypedName{newSemanticReverseAbstractInterpreterTest_TypedName("a", nullUndefined)},
		[]*semanticReverseAbstractInterpreterTest_TypedName{newSemanticReverseAbstractInterpreterTest_TypedName("a", c.GetNativeNumberType())})
}

// Tests reverse interpretation of two undefineds.
// Java: SemanticReverseAbstractInterpreterTest.testEqCondition4()
func (c *SemanticReverseAbstractInterpreterTest) TestEqCondition4() {
	blind := c.newScope()
	c.testBinop(
		blind,
		rhino.Token_EQ,
		c.createVar(blind, "a", c.GetNativeVoidType()),
		c.createVar(blind, "b", c.GetNativeVoidType()),
		[]*semanticReverseAbstractInterpreterTest_TypedName{
			newSemanticReverseAbstractInterpreterTest_TypedName("a", c.GetNativeVoidType()), newSemanticReverseAbstractInterpreterTest_TypedName("b", c.GetNativeVoidType())},
		[]*semanticReverseAbstractInterpreterTest_TypedName{
			newSemanticReverseAbstractInterpreterTest_TypedName("a", c.GetNativeNoType()), newSemanticReverseAbstractInterpreterTest_TypedName("b", c.GetNativeNoType())})
}

// Tests reverse interpretation of a COMPARE(NAME, NUMBER) expression, where
// COMPARE can be LE, LT, GE or GT.
// Java: SemanticReverseAbstractInterpreterTest.testInequalitiesCondition1()
func (c *SemanticReverseAbstractInterpreterTest) TestInequalitiesCondition1() {
	for _, op := range []rhino.Token{rhino.Token_LT, rhino.Token_GT, rhino.Token_LE, rhino.Token_GE} {
		blind := c.newScope()
		c.testBinop(
			blind,
			op,
			c.createVar(blind, "a", c.CreateUnionType(c.GetNativeStringType(), c.GetNativeVoidType())),
			c.createNumber(8),
			[]*semanticReverseAbstractInterpreterTest_TypedName{newSemanticReverseAbstractInterpreterTest_TypedName("a", c.GetNativeStringType())},
			[]*semanticReverseAbstractInterpreterTest_TypedName{
				newSemanticReverseAbstractInterpreterTest_TypedName("a", c.CreateUnionType(c.GetNativeStringType(), c.GetNativeVoidType()))})
	}
}

// Tests reverse interpretation of a COMPARE(NAME, NAME) expression, where
// COMPARE can be LE, LT, GE or GT.
// Java: SemanticReverseAbstractInterpreterTest.testInequalitiesCondition2()
func (c *SemanticReverseAbstractInterpreterTest) TestInequalitiesCondition2() {
	for _, op := range []rhino.Token{rhino.Token_LT, rhino.Token_GT, rhino.Token_LE, rhino.Token_GE} {
		blind := c.newScope()
		c.testBinop(
			blind,
			op,
			c.createVar(
				blind,
				"a",
				c.CreateUnionType(c.GetNativeStringType(), c.GetNativeNumberType(), c.GetNativeVoidType())),
			c.createVar(blind, "b", c.CreateUnionType(c.GetNativeNumberType(), c.GetNativeNullType())),
			[]*semanticReverseAbstractInterpreterTest_TypedName{
				newSemanticReverseAbstractInterpreterTest_TypedName("a", c.CreateUnionType(c.GetNativeStringType(), c.GetNativeNumberType())),
				newSemanticReverseAbstractInterpreterTest_TypedName("b", c.CreateUnionType(c.GetNativeNumberType(), c.GetNativeNullType()))},
			[]*semanticReverseAbstractInterpreterTest_TypedName{
				newSemanticReverseAbstractInterpreterTest_TypedName(
					"a",
					c.CreateUnionType(
						c.GetNativeStringType(), c.GetNativeNumberType(), c.GetNativeVoidType())),
				newSemanticReverseAbstractInterpreterTest_TypedName("b", c.CreateUnionType(c.GetNativeNumberType(), c.GetNativeNullType()))})
	}
}

// Tests reverse interpretation of a COMPARE(NUMBER-untyped, NAME) expression,
// where COMPARE can be LE, LT, GE or GT.
// Java: SemanticReverseAbstractInterpreterTest.testInequalitiesCondition3()
func (c *SemanticReverseAbstractInterpreterTest) TestInequalitiesCondition3() {
	for _, op := range []rhino.Token{rhino.Token_LT, rhino.Token_GT, rhino.Token_LE, rhino.Token_GE} {
		blind := c.newScope()
		c.testBinop(
			blind,
			op,
			c.createUntypedNumber(8),
			c.createVar(blind, "a", c.CreateUnionType(c.GetNativeStringType(), c.GetNativeVoidType())),
			[]*semanticReverseAbstractInterpreterTest_TypedName{newSemanticReverseAbstractInterpreterTest_TypedName("a", c.GetNativeStringType())},
			[]*semanticReverseAbstractInterpreterTest_TypedName{
				newSemanticReverseAbstractInterpreterTest_TypedName("a", c.CreateUnionType(c.GetNativeStringType(), c.GetNativeVoidType()))})
	}
}

// Java: SemanticReverseAbstractInterpreterTest.testAnd()
func (c *SemanticReverseAbstractInterpreterTest) TestAnd() {
	blind := c.newScope()
	c.testBinop(
		blind,
		rhino.Token_AND,
		c.createVar(blind, "b", c.CreateUnionType(c.GetNativeStringType(), c.GetNativeNullType())),
		c.createVar(blind, "a", c.CreateUnionType(c.GetNativeNumberType(), c.GetNativeVoidType())),
		[]*semanticReverseAbstractInterpreterTest_TypedName{
			newSemanticReverseAbstractInterpreterTest_TypedName("a", c.GetNativeNumberType()), newSemanticReverseAbstractInterpreterTest_TypedName("b", c.GetNativeStringType())},
		[]*semanticReverseAbstractInterpreterTest_TypedName{
			newSemanticReverseAbstractInterpreterTest_TypedName("a", c.CreateUnionType(c.GetNativeNumberType(), c.GetNativeVoidType())),
			newSemanticReverseAbstractInterpreterTest_TypedName("b", c.CreateUnionType(c.GetNativeStringType(), c.GetNativeNullType()))})
}

// Java: SemanticReverseAbstractInterpreterTest.testTypeof1()
func (c *SemanticReverseAbstractInterpreterTest) TestTypeof1() {
	blind := c.newScope()
	c.testBinop(
		blind,
		rhino.Token_EQ,
		rhino.NewNodeWithChild(rhino.Token_TYPEOF, c.createVar(blind, "a", c.GetNativeObjectType())),
		rhino.Node_NewString("function"),
		[]*semanticReverseAbstractInterpreterTest_TypedName{newSemanticReverseAbstractInterpreterTest_TypedName("a", c.GetNativeFunctionType())},
		[]*semanticReverseAbstractInterpreterTest_TypedName{newSemanticReverseAbstractInterpreterTest_TypedName("a", c.GetNativeObjectType())})
}

// Java: SemanticReverseAbstractInterpreterTest.testTypeof2()
func (c *SemanticReverseAbstractInterpreterTest) TestTypeof2() {
	blind := c.newScope()
	c.testBinop(
		blind,
		rhino.Token_EQ,
		rhino.NewNodeWithChild(rhino.Token_TYPEOF, c.createVar(blind, "a", c.GetNativeAllType())),
		rhino.Node_NewString("function"),
		[]*semanticReverseAbstractInterpreterTest_TypedName{newSemanticReverseAbstractInterpreterTest_TypedName("a", c.GetNativeFunctionType())},
		[]*semanticReverseAbstractInterpreterTest_TypedName{newSemanticReverseAbstractInterpreterTest_TypedName("a", c.GetNativeAllType())})
}

// Java: SemanticReverseAbstractInterpreterTest.testTypeof3()
func (c *SemanticReverseAbstractInterpreterTest) TestTypeof3() {
	blind := c.newScope()
	c.testBinop(
		blind,
		rhino.Token_EQ,
		rhino.NewNodeWithChild(rhino.Token_TYPEOF, c.createVar(blind, "a", c.GetNativeObjectNumberStringBooleanType())),
		rhino.Node_NewString("function"),
		[]*semanticReverseAbstractInterpreterTest_TypedName{newSemanticReverseAbstractInterpreterTest_TypedName("a", c.GetNativeFunctionType())},
		[]*semanticReverseAbstractInterpreterTest_TypedName{newSemanticReverseAbstractInterpreterTest_TypedName("a", c.GetNativeObjectNumberStringBooleanType())})
}

// Java: SemanticReverseAbstractInterpreterTest.testTypeof4()
func (c *SemanticReverseAbstractInterpreterTest) TestTypeof4() {
	blind := c.newScope()
	c.testBinop(
		blind,
		rhino.Token_EQ,
		rhino.NewNodeWithChild(
			rhino.Token_TYPEOF,
			c.createVar(
				blind,
				"a",
				c.CreateUnionType(c.GetNativeFunctionType(), c.GetNativeNumberStringBooleanType()))),
		rhino.Node_NewString("function"),
		[]*semanticReverseAbstractInterpreterTest_TypedName{newSemanticReverseAbstractInterpreterTest_TypedName("a", c.GetNativeFunctionType())},
		[]*semanticReverseAbstractInterpreterTest_TypedName{newSemanticReverseAbstractInterpreterTest_TypedName("a", c.GetNativeNumberStringBooleanType())})
}

// Java: SemanticReverseAbstractInterpreterTest.testInstanceOf()
func (c *SemanticReverseAbstractInterpreterTest) TestInstanceOf() {
	blind := c.newScope()
	c.testBinop(
		blind,
		rhino.Token_INSTANCEOF,
		c.createVar(blind, "x", c.GetNativeUnknownType()),
		c.createVar(blind, "s", c.GetNativeStringObjectConstructorType()),
		[]*semanticReverseAbstractInterpreterTest_TypedName{
			newSemanticReverseAbstractInterpreterTest_TypedName("x", c.GetNativeStringObjectType()),
			newSemanticReverseAbstractInterpreterTest_TypedName("s", c.GetNativeStringObjectConstructorType())},
		[]*semanticReverseAbstractInterpreterTest_TypedName{newSemanticReverseAbstractInterpreterTest_TypedName("s", c.GetNativeStringObjectConstructorType())})
}

// Java: SemanticReverseAbstractInterpreterTest.testInstanceOf2()
func (c *SemanticReverseAbstractInterpreterTest) TestInstanceOf2() {
	blind := c.newScope()
	c.testBinop(
		blind,
		rhino.Token_INSTANCEOF,
		c.createVar(
			blind, "x", c.CreateUnionType(c.GetNativeStringObjectType(), c.GetNativeNumberObjectType())),
		c.createVar(blind, "s", c.GetNativeStringObjectConstructorType()),
		[]*semanticReverseAbstractInterpreterTest_TypedName{
			newSemanticReverseAbstractInterpreterTest_TypedName("x", c.GetNativeStringObjectType()),
			newSemanticReverseAbstractInterpreterTest_TypedName("s", c.GetNativeStringObjectConstructorType())},
		[]*semanticReverseAbstractInterpreterTest_TypedName{
			newSemanticReverseAbstractInterpreterTest_TypedName("x", c.GetNativeNumberObjectType()),
			newSemanticReverseAbstractInterpreterTest_TypedName("s", c.GetNativeStringObjectConstructorType())})
}

// Java: SemanticReverseAbstractInterpreterTest.testInstanceOf3()
func (c *SemanticReverseAbstractInterpreterTest) TestInstanceOf3() {
	blind := c.newScope()
	c.testBinop(
		blind,
		rhino.Token_INSTANCEOF,
		c.createVar(blind, "x", c.GetNativeObjectType()),
		c.createVar(blind, "s", c.GetNativeStringObjectConstructorType()),
		[]*semanticReverseAbstractInterpreterTest_TypedName{
			newSemanticReverseAbstractInterpreterTest_TypedName("x", c.GetNativeStringObjectType()),
			newSemanticReverseAbstractInterpreterTest_TypedName("s", c.GetNativeStringObjectConstructorType())},
		[]*semanticReverseAbstractInterpreterTest_TypedName{
			newSemanticReverseAbstractInterpreterTest_TypedName("x", c.GetNativeObjectType()),
			newSemanticReverseAbstractInterpreterTest_TypedName("s", c.GetNativeStringObjectConstructorType())})
}

// Java: SemanticReverseAbstractInterpreterTest.testInstanceOf4()
func (c *SemanticReverseAbstractInterpreterTest) TestInstanceOf4() {
	blind := c.newScope()
	c.testBinop(
		blind,
		rhino.Token_INSTANCEOF,
		c.createVar(blind, "x", c.GetNativeAllType()),
		c.createVar(blind, "s", c.GetNativeStringObjectConstructorType()),
		[]*semanticReverseAbstractInterpreterTest_TypedName{
			newSemanticReverseAbstractInterpreterTest_TypedName("x", c.GetNativeStringObjectType()),
			newSemanticReverseAbstractInterpreterTest_TypedName("s", c.GetNativeStringObjectConstructorType())},
		[]*semanticReverseAbstractInterpreterTest_TypedName{newSemanticReverseAbstractInterpreterTest_TypedName("s", c.GetNativeStringObjectConstructorType())})
}

// Java: SemanticReverseAbstractInterpreterTest.testBinop(FlowScope[], Token, Node, Node, Collection<TypedName>, Collection<TypedName>)
func (c *SemanticReverseAbstractInterpreterTest) testBinop(
	blind []FlowScope,
	binop rhino.Token,
	left *rhino.Node,
	right *rhino.Node,
	trueOutcome []*semanticReverseAbstractInterpreterTest_TypedName,
	falseOutcome []*semanticReverseAbstractInterpreterTest_TypedName) {
	condition := rhino.NewNode(binop)
	condition.AddChildToBack(left)
	condition.AddChildToBack(right)

	// true outcome.
	informedTrue :=
		c.interpreter.GetPreciserScopeKnowingConditionOutcome(condition, blind[0], rhino.Outcome_TRUE)
	for _, p := range trueOutcome {
		c.AssertTypeEqualsMsg(p.name, p.type_, c.getVarType(informedTrue, p.name))
	}

	// false outcome.
	informedFalse :=
		c.interpreter.GetPreciserScopeKnowingConditionOutcome(condition, blind[0], rhino.Outcome_FALSE)
	for _, p := range falseOutcome {
		c.AssertTypeEquals(p.type_, c.getVarType(informedFalse, p.name))
	}
}

// Java: SemanticReverseAbstractInterpreterTest.createNull()
func (c *SemanticReverseAbstractInterpreterTest) createNull() *rhino.Node {
	n := rhino.NewNode(rhino.Token_NULL)
	n.SetJSType(c.GetNativeNullType())
	return n
}

// Java: SemanticReverseAbstractInterpreterTest.createNumber(int)
func (c *SemanticReverseAbstractInterpreterTest) createNumber(n int) *rhino.Node {
	number := c.createUntypedNumber(n)
	number.SetJSType(c.GetNativeNumberType())
	return number
}

// Java: SemanticReverseAbstractInterpreterTest.createUntypedNumber(int)
func (c *SemanticReverseAbstractInterpreterTest) createUntypedNumber(n int) *rhino.Node {
	return rhino.Node_NewNumber(float64(n))
}

// Java: SemanticReverseAbstractInterpreterTest.getVarType(FlowScope, String)
func (c *SemanticReverseAbstractInterpreterTest) getVarType(scope FlowScope, name string) rhino.JSType {
	return scope.GetSlot(name).(rhino.StaticTypedSlot).GetType()
}

// Java: SemanticReverseAbstractInterpreterTest.createVar(FlowScope[], String, JSType)
func (c *SemanticReverseAbstractInterpreterTest) createVar(scope []FlowScope, name string, type_ rhino.JSType) *rhino.Node {
	n := rhino.Node_NewStringWithToken(rhino.Token_NAME, name)
	c.functionScope.Declare(name, n, nil, nil, true)
	scope[0] = scope[0].InferSlotType(name, type_)
	n.SetJSType(type_)
	return n
}

// Java: SemanticReverseAbstractInterpreterTest.TypedName
type semanticReverseAbstractInterpreterTest_TypedName struct {
	name  string
	type_ rhino.JSType
}

// Java: SemanticReverseAbstractInterpreterTest.TypedName(String, JSType)
func newSemanticReverseAbstractInterpreterTest_TypedName(name string, type_ rhino.JSType) *semanticReverseAbstractInterpreterTest_TypedName {
	return &semanticReverseAbstractInterpreterTest_TypedName{name: name, type_: type_}
}
