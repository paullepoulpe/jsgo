// Port of com/google/javascript/jscomp/ClosureReverseAbstractInterpreterTest.java @ closure-compiler v20260415

package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func TestClosureReverseAbstractInterpreter_GoogIsObjectOnNull(t *testing.T) {
	c := newCompilerTypeTestCase(t)
	testClosureFunction(c, "goog.isObject", nil, c.getNativeObjectType(), nil)
}

func TestClosureReverseAbstractInterpreter_GoogIsObject1(t *testing.T) {
	c := newCompilerTypeTestCase(t)
	testClosureFunction(c,
		"goog.isObject", c.getNativeAllType(), c.getNativeNoObjectType(), c.getNativeAllType())
}

func TestClosureReverseAbstractInterpreter_GoogIsObject2a(t *testing.T) {
	c := newCompilerTypeTestCase(t)
	testClosureFunction(c,
		"goog.isObject",
		c.createUnionType(c.getNativeObjectType(), c.getNativeNumberStringBooleanType()),
		c.getNativeObjectType(),
		c.getNativeNumberStringBooleanType())
}

func TestClosureReverseAbstractInterpreter_GoogIsObject2b(t *testing.T) {
	c := newCompilerTypeTestCase(t)
	testClosureFunction(c,
		"goog.isObject",
		c.createUnionType(c.getNativeObjectType(), c.getNativeValueTypes()),
		c.getNativeObjectType(),
		c.getNativeValueTypes())
}

func TestClosureReverseAbstractInterpreter_GoogIsObject3a(t *testing.T) {
	c := newCompilerTypeTestCase(t)
	testClosureFunction(c,
		"goog.isObject",
		c.createUnionType(
			c.getNativeObjectType(),
			c.getNativeNumberStringBooleanType(),
			c.getNativeNullType(),
			c.getNativeVoidType()),
		c.getNativeObjectType(),
		c.createUnionType(
			c.getNativeNumberStringBooleanType(), c.getNativeNullType(), c.getNativeVoidType()))
}

func TestClosureReverseAbstractInterpreter_GoogIsObject3b(t *testing.T) {
	c := newCompilerTypeTestCase(t)
	testClosureFunction(c,
		"goog.isObject",
		c.createUnionType(
			c.getNativeObjectType(), c.getNativeValueTypes(), c.getNativeNullType(), c.getNativeVoidType()),
		c.getNativeObjectType(),
		c.createUnionType(c.getNativeValueTypes(), c.getNativeNullType(), c.getNativeVoidType()))
}

func TestClosureReverseAbstractInterpreter_GoogIsObject4(t *testing.T) {
	c := newCompilerTypeTestCase(t)
	testClosureFunction(c,
		"goog.isObject",
		c.getNativeUnknownType(),
		c.getNativeNoObjectType(), // ? Should this be CHECKED_UNKNOWN?
		c.getNativeCheckedUnknownType())
}

// testClosureFunction checks `function(a)` where a : type.
//
// Go-only: test adapter — the Java test parses "var a; goog.isObject(a)"
// with compiler.parseTestCode and builds the flow scope with
// TypedScopeCreator and LinkedFlowScope.createEntryLattice, none of which is
// ported yet. Here the CALL node is built with IR (the shape the parser
// produces: a SCRIPT whose last statement is the call) and the flow scope is
// the copy-on-write test scope; `a` is declared without a type, as `var a;`
// declares it.
func testClosureFunction(
	c *compilerTypeTestCase, function string, type_ rhino.JSType, trueType rhino.JSType, falseType rhino.JSType) {
	c.t.Helper()
	// function(a) where a : type
	parts := rhino.QualifiedName_Of(function)
	n := rhino.IR_ScriptOfNodes(
		rhino.IR_Var(rhino.IR_Name("a")),
		rhino.IR_ExprResult(rhino.IR_Call(
			rhino.IR_Getprop(rhino.IR_Name(parts.GetOwner().Join()), parts.GetComponent()),
			rhino.IR_Name("a"))))
	call := n.GetLastChild().GetLastChild()
	name := call.GetLastChild()

	var flowScope FlowScope = newTestFlowScope()

	c.truth.AssertThat(call.GetToken()).IsEqualTo(rhino.Token_CALL)
	c.truth.AssertThat(name.GetToken()).IsEqualTo(rhino.Token_NAME)

	flowScope = flowScope.InferSlotType("a", type_)
	rai := NewClosureReverseAbstractInterpreter(c.registry)

	// trueScope
	c.truth.AssertType(
		rai.GetPreciserScopeKnowingConditionOutcome(call, flowScope, rhino.Outcome_TRUE).
			GetSlot("a").(rhino.StaticTypedSlot).
			GetType()).
		IsEqualTo(trueType)

	// falseScope
	aType := rai.GetPreciserScopeKnowingConditionOutcome(call, flowScope, rhino.Outcome_FALSE).
		GetSlot("a").(rhino.StaticTypedSlot).
		GetType()
	if falseType == nil {
		c.truth.AssertThat(aType).IsNull()
	} else {
		c.truth.AssertType(aType).IsEqualTo(falseType)
	}
}
