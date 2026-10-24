// Port of com/google/javascript/jscomp/LinkedFlowScopeTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/LinkedFlowScopeTest.java)

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

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// LinkedFlowScopeTest holds the tests for LinkedFlowScope.
type LinkedFlowScopeTest struct {
	CompilerTypeTestCaseBase

	functionNode rhino.Node
	rootNode     rhino.Node

	localScope *TypedScope

	globalEntry FlowScope

	localEntry FlowScope
}

const LinkedFlowScopeTest_LONG_CHAIN_LENGTH = 1050

// Java: LinkedFlowScopeTest()
func NewLinkedFlowScopeTest() *LinkedFlowScopeTest {
	c := &LinkedFlowScopeTest{}
	c.InitCompilerTypeTestCaseBase(c)
	c.functionNode = rhino.NewNode(rhino.Token_FUNCTION)
	c.rootNode = rhino.NewNodeWithChild(rhino.Token_ROOT, c.functionNode)
	return c
}

func TestLinkedFlowScopeTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewLinkedFlowScopeTest() })
}

// Java: LinkedFlowScopeTest.setUp()
func (c *LinkedFlowScopeTest) SetUp() {
	c.CompilerTypeTestCaseBase.SetUp()

	globalScope := TypedScope_CreateGlobalScope(c.rootNode)
	globalScope.Declare("globalA", rhino.Node{}, nil, nil, true)
	globalScope.Declare("globalB", rhino.Node{}, nil, nil, true)

	c.localScope = NewTypedScope(globalScope, c.functionNode)
	c.localScope.Declare("localA", rhino.Node{}, nil, nil, true)
	c.localScope.Declare("localB", rhino.Node{}, nil, nil, true)

	c.globalEntry = LinkedFlowScope_CreateEntryLattice(c.compiler, globalScope)
	c.localEntry = LinkedFlowScope_CreateEntryLattice(c.compiler, c.localScope)
}

// linkedFlowScopeTest_slotType is scope.getSlot(name).getType().
//
// Go-only: language-adapter — FlowScope.GetSlot keeps StaticScope's Go
// signature (see type_flow_scope.go), so the StaticTypedSlot it returns is
// asserted back.
func linkedFlowScopeTest_slotType(scope FlowScope, name string) rhino.JSType {
	return scope.GetSlot(name).(rhino.StaticTypedSlot).GetType()
}

// Java: LinkedFlowScopeTest.testJoin1()
func (c *LinkedFlowScopeTest) TestJoin1() {
	childA := c.localEntry.InferSlotType("localB", c.GetNativeNumberType())
	childAB := childA.InferSlotType("localB", c.GetNativeStringType())
	childB := c.localEntry.InferSlotType("localB", c.GetNativeBooleanType())

	c.AssertTypeEquals(c.GetNativeStringType(), linkedFlowScopeTest_slotType(childAB, "localB"))
	c.AssertTypeEquals(c.GetNativeBooleanType(), linkedFlowScopeTest_slotType(childB, "localB"))
	AssertThat(linkedFlowScopeTest_slotType(childB, "localA")).IsNull()

	joined := c.Join(childB, childAB)
	c.AssertTypeEquals(
		c.CreateUnionType(c.GetNativeStringType(), c.GetNativeBooleanType()),
		linkedFlowScopeTest_slotType(joined, "localB"))
	AssertThat(linkedFlowScopeTest_slotType(joined, "localA")).IsNull()

	joined = c.Join(childAB, childB)
	c.AssertTypeEquals(
		c.CreateUnionType(c.GetNativeStringType(), c.GetNativeBooleanType()),
		linkedFlowScopeTest_slotType(joined, "localB"))
	AssertThat(linkedFlowScopeTest_slotType(joined, "localA")).IsNull()

	AssertWithMessage("Join should be symmetric").
		That(c.Join(childAB, childB)).
		IsEqualTo(c.Join(childB, childAB))
}

// Java: LinkedFlowScopeTest.testJoin2()
func (c *LinkedFlowScopeTest) TestJoin2() {
	childA := c.localEntry.InferSlotType("localA", c.GetNativeStringType())
	childB := c.localEntry.InferSlotType("globalB", c.GetNativeBooleanType())

	c.AssertTypeEquals(c.GetNativeStringType(), linkedFlowScopeTest_slotType(childA, "localA"))
	c.AssertTypeEquals(c.GetNativeBooleanType(), linkedFlowScopeTest_slotType(childB, "globalB"))
	AssertThat(linkedFlowScopeTest_slotType(childB, "localB")).IsNull()

	joined := c.Join(childB, childA)
	c.AssertTypeEquals(c.GetNativeStringType(), linkedFlowScopeTest_slotType(joined, "localA"))
	c.AssertTypeEquals(c.GetNativeBooleanType(), linkedFlowScopeTest_slotType(joined, "globalB"))

	joined = c.Join(childA, childB)
	c.AssertTypeEquals(c.GetNativeStringType(), linkedFlowScopeTest_slotType(joined, "localA"))
	c.AssertTypeEquals(c.GetNativeBooleanType(), linkedFlowScopeTest_slotType(joined, "globalB"))

	AssertWithMessage("Join should be symmetric").
		That(c.Join(childA, childB)).
		IsEqualTo(c.Join(childB, childA))
}

// Java: LinkedFlowScopeTest.testJoin3()
func (c *LinkedFlowScopeTest) TestJoin3() {
	c.localScope.Declare("localC", rhino.Node{}, c.GetNativeStringType(), nil, true)
	c.localScope.Declare("localD", rhino.Node{}, c.GetNativeStringType(), nil, true)

	childA := c.localEntry.InferSlotType("localC", c.GetNativeNumberType())
	childB := c.localEntry.InferSlotType("localD", c.GetNativeBooleanType())

	joined := c.Join(childB, childA)
	c.AssertTypeEquals(
		c.CreateUnionType(c.GetNativeStringType(), c.GetNativeNumberType()),
		linkedFlowScopeTest_slotType(joined, "localC"))
	c.AssertTypeEquals(
		c.CreateUnionType(c.GetNativeStringType(), c.GetNativeBooleanType()),
		linkedFlowScopeTest_slotType(joined, "localD"))

	joined = c.Join(childA, childB)
	c.AssertTypeEquals(
		c.CreateUnionType(c.GetNativeStringType(), c.GetNativeNumberType()),
		linkedFlowScopeTest_slotType(joined, "localC"))
	c.AssertTypeEquals(
		c.CreateUnionType(c.GetNativeStringType(), c.GetNativeBooleanType()),
		linkedFlowScopeTest_slotType(joined, "localD"))

	AssertWithMessage("Join should be symmetric").
		That(c.Join(childA, childB)).
		IsEqualTo(c.Join(childB, childA))
}

// TestLongChain creates a long chain of flow scopes.
//
// Java: LinkedFlowScopeTest.testLongChain()
func (c *LinkedFlowScopeTest) TestLongChain() {
	chainA := c.localEntry
	chainB := c.localEntry
	for i := 0; i < LinkedFlowScopeTest_LONG_CHAIN_LENGTH; i++ {
		c.localScope.Declare("local"+strconv.Itoa(i), rhino.Node{}, nil, nil, true)
		typeA := c.GetNativeBooleanType()
		if i%2 == 0 {
			typeA = c.GetNativeNumberType()
		}
		chainA =
			chainA.InferSlotType("local"+strconv.Itoa(i), typeA)
		typeB := c.GetNativeBooleanType()
		if i%3 == 0 {
			typeB = c.GetNativeStringType()
		}
		chainB =
			chainB.InferSlotType("local"+strconv.Itoa(i), typeB)
	}

	joined := c.Join(chainA, chainB)
	for i := 0; i < LinkedFlowScopeTest_LONG_CHAIN_LENGTH; i++ {
		typeA := c.GetNativeBooleanType()
		if i%2 == 0 {
			typeA = c.GetNativeNumberType()
		}
		c.AssertTypeEquals(typeA, linkedFlowScopeTest_slotType(chainA, "local"+strconv.Itoa(i)))
		typeB := c.GetNativeBooleanType()
		if i%3 == 0 {
			typeB = c.GetNativeStringType()
		}
		c.AssertTypeEquals(typeB, linkedFlowScopeTest_slotType(chainB, "local"+strconv.Itoa(i)))

		joinedSlotType := linkedFlowScopeTest_slotType(joined, "local"+strconv.Itoa(i))
		if i%6 == 0 {
			c.AssertTypeEquals(
				c.CreateUnionType(c.GetNativeStringType(), c.GetNativeNumberType()), joinedSlotType)
		} else if i%2 == 0 {
			c.AssertTypeEquals(
				c.CreateUnionType(c.GetNativeNumberType(), c.GetNativeBooleanType()), joinedSlotType)
		} else if i%3 == 0 {
			c.AssertTypeEquals(
				c.CreateUnionType(c.GetNativeStringType(), c.GetNativeBooleanType()), joinedSlotType)
		} else {
			c.AssertTypeEquals(c.GetNativeBooleanType(), joinedSlotType)
		}
	}

	c.AssertScopesDiffer(chainA, chainB)
	c.AssertScopesDiffer(chainA, joined)
	c.AssertScopesDiffer(chainB, joined)
}

// Java: LinkedFlowScopeTest.testDiffer1()
func (c *LinkedFlowScopeTest) TestDiffer1() {
	childA := c.localEntry.InferSlotType("localB", c.GetNativeNumberType())
	childAB := childA.InferSlotType("localB", c.GetNativeStringType())
	childABC := childAB.InferSlotType("localA", c.GetNativeBooleanType())
	childB := childAB.InferSlotType("localB", c.GetNativeStringType())
	childBC := childB.InferSlotType("localA", c.GetNativeNoType())

	c.AssertScopesSame(childAB, childB)
	c.AssertScopesDiffer(childABC, childBC)

	c.AssertScopesDiffer(childABC, childB)
	c.AssertScopesDiffer(childAB, childBC)

	c.AssertScopesDiffer(childA, childAB)
	c.AssertScopesDiffer(childA, childABC)
	c.AssertScopesDiffer(childA, childB)
	c.AssertScopesDiffer(childA, childBC)
}

// Java: LinkedFlowScopeTest.testDiffer2()
func (c *LinkedFlowScopeTest) TestDiffer2() {
	childA := c.localEntry.InferSlotType("localA", c.GetNativeNumberType())
	childB := c.localEntry.InferSlotType("localA", c.GetNativeNoType())

	c.AssertScopesDiffer(childA, childB)
}

// Java: LinkedFlowScopeTest.assertScopesDiffer(FlowScope, FlowScope)
func (c *LinkedFlowScopeTest) AssertScopesDiffer(a FlowScope, b FlowScope) {
	AssertThat(a).IsNotEqualTo(b)
	AssertThat(b).IsNotEqualTo(a)
}

// Java: LinkedFlowScopeTest.assertScopesSame(FlowScope, FlowScope)
func (c *LinkedFlowScopeTest) AssertScopesSame(a FlowScope, b FlowScope) {
	AssertThat(b).IsEqualTo(a)
	AssertThat(a).IsEqualTo(b)
}

// Java: LinkedFlowScopeTest.join(FlowScope, FlowScope)
func (c *LinkedFlowScopeTest) Join(a FlowScope, b FlowScope) FlowScope {
	joiner := NewLinkedFlowScope_FlowScopeJoinOp(c.compiler)
	joiner.JoinFlow(a)
	joiner.JoinFlow(b)
	return joiner.Finish()
}
