// Port of com/google/javascript/jscomp/ClosureReverseAbstractInterpreterTest.java @ closure-compiler v20260415

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

// Java: public final class ClosureReverseAbstractInterpreterTest extends CompilerTypeTestCase
type ClosureReverseAbstractInterpreterTest struct {
	CompilerTypeTestCaseBase
}

// Java: ClosureReverseAbstractInterpreterTest()
func NewClosureReverseAbstractInterpreterTest() *ClosureReverseAbstractInterpreterTest {
	c := &ClosureReverseAbstractInterpreterTest{}
	c.InitCompilerTypeTestCaseBase(c)
	return c
}

func TestClosureReverseAbstractInterpreterTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewClosureReverseAbstractInterpreterTest() })
}

// Java: ClosureReverseAbstractInterpreterTest.testGoogIsObjectOnNull()
func (c *ClosureReverseAbstractInterpreterTest) TestGoogIsObjectOnNull() {
	c.testClosureFunction("goog.isObject", nil, c.GetNativeObjectType(), nil)
}

// Java: ClosureReverseAbstractInterpreterTest.testGoogIsObject1()
func (c *ClosureReverseAbstractInterpreterTest) TestGoogIsObject1() {
	c.testClosureFunction(
		"goog.isObject", c.GetNativeAllType(), c.GetNativeNoObjectType(), c.GetNativeAllType())
}

// Java: ClosureReverseAbstractInterpreterTest.testGoogIsObject2a()
func (c *ClosureReverseAbstractInterpreterTest) TestGoogIsObject2a() {
	c.testClosureFunction(
		"goog.isObject",
		c.CreateUnionType(c.GetNativeObjectType(), c.GetNativeNumberStringBooleanType()),
		c.GetNativeObjectType(),
		c.GetNativeNumberStringBooleanType())
}

// Java: ClosureReverseAbstractInterpreterTest.testGoogIsObject2b()
func (c *ClosureReverseAbstractInterpreterTest) TestGoogIsObject2b() {
	c.testClosureFunction(
		"goog.isObject",
		c.CreateUnionType(c.GetNativeObjectType(), c.GetNativeValueTypes()),
		c.GetNativeObjectType(),
		c.GetNativeValueTypes())
}

// Java: ClosureReverseAbstractInterpreterTest.testGoogIsObject3a()
func (c *ClosureReverseAbstractInterpreterTest) TestGoogIsObject3a() {
	c.testClosureFunction(
		"goog.isObject",
		c.CreateUnionType(
			c.GetNativeObjectType(),
			c.GetNativeNumberStringBooleanType(),
			c.GetNativeNullType(),
			c.GetNativeVoidType()),
		c.GetNativeObjectType(),
		c.CreateUnionType(
			c.GetNativeNumberStringBooleanType(), c.GetNativeNullType(), c.GetNativeVoidType()))
}

// Java: ClosureReverseAbstractInterpreterTest.testGoogIsObject3b()
func (c *ClosureReverseAbstractInterpreterTest) TestGoogIsObject3b() {
	c.testClosureFunction(
		"goog.isObject",
		c.CreateUnionType(
			c.GetNativeObjectType(), c.GetNativeValueTypes(), c.GetNativeNullType(), c.GetNativeVoidType()),
		c.GetNativeObjectType(),
		c.CreateUnionType(c.GetNativeValueTypes(), c.GetNativeNullType(), c.GetNativeVoidType()))
}

// Java: ClosureReverseAbstractInterpreterTest.testGoogIsObject4()
func (c *ClosureReverseAbstractInterpreterTest) TestGoogIsObject4() {
	c.testClosureFunction(
		"goog.isObject",
		c.GetNativeUnknownType(),
		c.GetNativeNoObjectType(), // ? Should this be CHECKED_UNKNOWN?
		c.GetNativeCheckedUnknownType())
}

// Java: ClosureReverseAbstractInterpreterTest.testClosureFunction(String, JSType, JSType, JSType)
func (c *ClosureReverseAbstractInterpreterTest) testClosureFunction(
	function string, type_ rhino.JSType, trueType rhino.JSType, falseType rhino.JSType) {
	// function(a) where a : type
	n := c.compiler.ParseTestCode("var a; " + function + "(a)")
	call := n.GetLastChild().GetLastChild()
	name := call.GetLastChild()

	root := rhino.IR_Root(rhino.IR_Root(), rhino.IR_Root(n))
	scope := NewTypedScopeCreator(c.compiler).CreateScope(root, nil).(*TypedScope)
	var flowScope FlowScope = LinkedFlowScope_CreateEntryLattice(c.compiler, scope)

	AssertThat(call.GetToken()).IsEqualTo(rhino.Token_CALL)
	AssertThat(name.GetToken()).IsEqualTo(rhino.Token_NAME)

	flowScope = flowScope.InferSlotType("a", type_)
	rai := NewClosureReverseAbstractInterpreter(c.registry)

	// trueScope
	TypeSubject_AssertType(
		rai.GetPreciserScopeKnowingConditionOutcome(call, flowScope, rhino.Outcome_TRUE).
			GetSlot("a").(rhino.StaticTypedSlot).
			GetType()).
		IsEqualTo(trueType)

	// falseScope
	aType :=
		rai.GetPreciserScopeKnowingConditionOutcome(call, flowScope, rhino.Outcome_FALSE).
			GetSlot("a").(rhino.StaticTypedSlot).
			GetType()
	if falseType == nil {
		AssertThat(aType).IsNull()
	} else {
		TypeSubject_AssertType(aType).IsEqualTo(falseType)
	}
}
