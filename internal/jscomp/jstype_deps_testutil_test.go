package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/rhinotesting"
)

// compilerTypeTestCase stands in for CompilerTypeTestCase (a
// test/com/google/javascript/jscomp class) until Compiler is ported: it
// keeps the registry and the native-type helpers the ported tests use, on a
// registry created directly instead of compiler.getTypeRegistry().
//
// Go-only: test adapter — the type-system state the tests below rely on
// (native types) is the same in both registries; the Compiler-only parts of
// CompilerTypeTestCase (parsing, TypedScopeCreator, warnings) are not
// reproduced, and the tests that need them are adapted or skipped.
type compilerTypeTestCase struct {
	t        *testing.T
	truth    rhinotesting.Truth
	registry *rhino.JSTypeRegistry
}

func newCompilerTypeTestCase(t *testing.T) *compilerTypeTestCase {
	return &compilerTypeTestCase{
		t:        t,
		truth:    rhinotesting.NewTruth(t),
		registry: rhino.NewJSTypeRegistry(rhinotesting.NewTestErrorReporter()),
	}
}

func (c *compilerTypeTestCase) createUnionType(variants ...rhino.JSType) rhino.JSType {
	return c.registry.CreateUnionType(variants...)
}

func (c *compilerTypeTestCase) createNullableType(type_ rhino.JSType) rhino.JSType {
	return c.registry.CreateNullableType(type_)
}

func (c *compilerTypeTestCase) assertTypeEquals(a rhino.JSType, b rhino.JSType) {
	c.t.Helper()
	c.truth.AssertType(b).IsEqualTo(a)
}

func (c *compilerTypeTestCase) assertTypeEqualsMsg(msg string, a rhino.JSType, b rhino.JSType) {
	c.t.Helper()
	c.truth.AssertWithMessage("%s", msg).AboutTypes().That(b).IsEqualTo(a)
}

func (c *compilerTypeTestCase) getNativeType(n rhino.JSTypeNative) rhino.JSType {
	return c.registry.GetNativeType(n)
}

func (c *compilerTypeTestCase) getNativeObjectTypeOf(n rhino.JSTypeNative) rhino.ObjectType {
	return c.registry.GetNativeObjectType(n)
}

func (c *compilerTypeTestCase) getNativeFunctionTypeOf(n rhino.JSTypeNative) rhino.FunctionType {
	return c.registry.GetNativeFunctionType(n)
}

func (c *compilerTypeTestCase) getNativeNoObjectType() rhino.ObjectType {
	return c.getNativeObjectTypeOf(rhino.JSTypeNative_NO_OBJECT_TYPE)
}

func (c *compilerTypeTestCase) getNativeStringObjectType() rhino.ObjectType {
	return c.getNativeObjectTypeOf(rhino.JSTypeNative_STRING_OBJECT_TYPE)
}

func (c *compilerTypeTestCase) getNativeNumberObjectType() rhino.ObjectType {
	return c.getNativeObjectTypeOf(rhino.JSTypeNative_NUMBER_OBJECT_TYPE)
}

func (c *compilerTypeTestCase) getNativeNoType() rhino.ObjectType {
	return c.getNativeObjectTypeOf(rhino.JSTypeNative_NO_TYPE)
}

func (c *compilerTypeTestCase) getNativeUnknownType() rhino.ObjectType {
	return c.getNativeObjectTypeOf(rhino.JSTypeNative_UNKNOWN_TYPE)
}

func (c *compilerTypeTestCase) getNativeCheckedUnknownType() rhino.ObjectType {
	return c.getNativeObjectTypeOf(rhino.JSTypeNative_CHECKED_UNKNOWN_TYPE)
}

func (c *compilerTypeTestCase) getNativeObjectType() rhino.ObjectType {
	return c.getNativeObjectTypeOf(rhino.JSTypeNative_OBJECT_TYPE)
}

func (c *compilerTypeTestCase) getNativeStringObjectConstructorType() rhino.FunctionType {
	return c.getNativeFunctionTypeOf(rhino.JSTypeNative_STRING_OBJECT_FUNCTION_TYPE)
}

func (c *compilerTypeTestCase) getNativeFunctionType() rhino.FunctionType {
	return c.getNativeFunctionTypeOf(rhino.JSTypeNative_FUNCTION_TYPE)
}

func (c *compilerTypeTestCase) getNativeVoidType() rhino.JSType {
	return c.getNativeType(rhino.JSTypeNative_VOID_TYPE)
}

func (c *compilerTypeTestCase) getNativeNullType() rhino.JSType {
	return c.getNativeType(rhino.JSTypeNative_NULL_TYPE)
}

func (c *compilerTypeTestCase) getNativeNumberType() rhino.JSType {
	return c.getNativeType(rhino.JSTypeNative_NUMBER_TYPE)
}

func (c *compilerTypeTestCase) getNativeBooleanType() rhino.JSType {
	return c.getNativeType(rhino.JSTypeNative_BOOLEAN_TYPE)
}

func (c *compilerTypeTestCase) getNativeStringType() rhino.JSType {
	return c.getNativeType(rhino.JSTypeNative_STRING_TYPE)
}

func (c *compilerTypeTestCase) getNativeObjectNumberStringBooleanType() rhino.JSType {
	return c.registry.CreateUnionTypeOfNatives(
		rhino.JSTypeNative_OBJECT_TYPE, rhino.JSTypeNative_NUMBER_TYPE,
		rhino.JSTypeNative_STRING_TYPE, rhino.JSTypeNative_BOOLEAN_TYPE)
}

func (c *compilerTypeTestCase) getNativeNumberStringBooleanType() rhino.JSType {
	return c.getNativeType(rhino.JSTypeNative_NUMBER_STRING_BOOLEAN)
}

func (c *compilerTypeTestCase) getNativeValueTypes() rhino.JSType {
	return c.getNativeType(rhino.JSTypeNative_VALUE_TYPES)
}

func (c *compilerTypeTestCase) getNativeAllType() rhino.JSType {
	return c.getNativeType(rhino.JSTypeNative_ALL_TYPE)
}

// newTestFlowScope stands in for LinkedFlowScope.createEntryLattice over a
// fresh TypedScope: an empty copy-on-write flow scope (the oracle's
// raiFlowScope).
//
// Go-only: test adapter — TypedScope and LinkedFlowScope are not ported yet.
func newTestFlowScope() FlowScope {
	return &raiFlowScope{slots: raiCopySlots(newRaiCase().slots)}
}
