// Port of com/google/javascript/jscomp/SemanticReverseAbstractInterpreterTest.java @ closure-compiler v20260415

package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: test adapter — the Java test builds its blind scope with
// TypedScope.createGlobalScope / new TypedScope(...) and
// LinkedFlowScope.createEntryLattice(compiler, functionScope), which are not
// ported yet; newScope uses the copy-on-write test flow scope instead (the
// one the reverse-interpreter oracle checks against Java), and the registry
// is CompilerTypeTestCase's stand-in. When TypedScope and LinkedFlowScope
// land, switch newScope/createVar back to them.
type semanticRAITest struct {
	*compilerTypeTestCase
	interpreter ReverseAbstractInterpreter
}

func newSemanticRAITest(t *testing.T) *semanticRAITest {
	s := &semanticRAITest{compilerTypeTestCase: newCompilerTypeTestCase(t)}
	s.interpreter = NewSemanticReverseAbstractInterpreter(s.registry)
	return s
}

func (s *semanticRAITest) newScope() []FlowScope {
	return []FlowScope{newTestFlowScope()}
}

type semanticRAITest_TypedName struct {
	name  string
	type_ rhino.JSType
}

func typedName(name string, type_ rhino.JSType) semanticRAITest_TypedName {
	return semanticRAITest_TypedName{name: name, type_: type_}
}

func (s *semanticRAITest) testBinop(
	blind []FlowScope,
	binop rhino.Token,
	left *rhino.Node,
	right *rhino.Node,
	trueOutcome []semanticRAITest_TypedName,
	falseOutcome []semanticRAITest_TypedName) {
	s.t.Helper()
	condition := rhino.NewNode(binop)
	condition.AddChildToBack(left)
	condition.AddChildToBack(right)

	// true outcome.
	informedTrue := s.interpreter.GetPreciserScopeKnowingConditionOutcome(condition, blind[0], rhino.Outcome_TRUE)
	for _, p := range trueOutcome {
		s.assertTypeEqualsMsg(p.name, p.type_, s.getVarType(informedTrue, p.name))
	}

	// false outcome.
	informedFalse := s.interpreter.GetPreciserScopeKnowingConditionOutcome(condition, blind[0], rhino.Outcome_FALSE)
	for _, p := range falseOutcome {
		s.assertTypeEquals(p.type_, s.getVarType(informedFalse, p.name))
	}
}

func (s *semanticRAITest) createNull() *rhino.Node {
	n := rhino.NewNode(rhino.Token_NULL)
	n.SetJSType(s.getNativeNullType())
	return n
}

func (s *semanticRAITest) createNumber(n int) *rhino.Node {
	number := s.createUntypedNumber(n)
	number.SetJSType(s.getNativeNumberType())
	return number
}

func (s *semanticRAITest) createUntypedNumber(n int) *rhino.Node {
	return rhino.Node_NewNumber(float64(n))
}

func (s *semanticRAITest) getVarType(scope FlowScope, name string) rhino.JSType {
	return scope.GetSlot(name).(rhino.StaticTypedSlot).GetType()
}

func (s *semanticRAITest) createVar(scope []FlowScope, name string, type_ rhino.JSType) *rhino.Node {
	n := rhino.Node_NewStringWithToken(rhino.Token_NAME, name)
	// Java: functionScope.declare(name, n, null, null, true)
	scope[0] = scope[0].InferSlotType(name, type_)
	n.SetJSType(type_)
	return n
}

// Tests reverse interpretation of a NAME expression.
func TestSemanticReverseAbstractInterpreter_NameCondition(t *testing.T) {
	s := newSemanticRAITest(t)
	blind := s.newScope()
	condition := s.createVar(blind, "a", s.createNullableType(s.getNativeStringType()))

	// true outcome.
	informedTrue := s.interpreter.GetPreciserScopeKnowingConditionOutcome(condition, blind[0], rhino.Outcome_TRUE)
	s.assertTypeEquals(s.getNativeStringType(), s.getVarType(informedTrue, "a"))

	// false outcome.
	informedFalse := s.interpreter.GetPreciserScopeKnowingConditionOutcome(condition, blind[0], rhino.Outcome_FALSE)
	s.assertTypeEquals(s.createNullableType(s.getNativeStringType()), s.getVarType(informedFalse, "a"))
}

// Tests reverse interpretation of a NOT(NAME) expression.
func TestSemanticReverseAbstractInterpreter_NegatedNameCondition(t *testing.T) {
	s := newSemanticRAITest(t)
	blind := s.newScope()
	a := s.createVar(blind, "a", s.createNullableType(s.getNativeStringType()))
	condition := rhino.NewNode(rhino.Token_NOT)
	condition.AddChildToBack(a)

	// true outcome.
	informedTrue := s.interpreter.GetPreciserScopeKnowingConditionOutcome(condition, blind[0], rhino.Outcome_TRUE)
	s.assertTypeEquals(s.createNullableType(s.getNativeStringType()), s.getVarType(informedTrue, "a"))

	// false outcome.
	informedFalse := s.interpreter.GetPreciserScopeKnowingConditionOutcome(condition, blind[0], rhino.Outcome_FALSE)
	s.assertTypeEquals(s.getNativeStringType(), s.getVarType(informedFalse, "a"))
}

// Tests reverse interpretation of a ASSIGN expression.
func TestSemanticReverseAbstractInterpreter_AssignCondition1(t *testing.T) {
	s := newSemanticRAITest(t)
	blind := s.newScope()
	s.testBinop(
		blind,
		rhino.Token_ASSIGN,
		s.createVar(blind, "a", s.createNullableType(s.getNativeObjectType())),
		s.createVar(blind, "b", s.createNullableType(s.getNativeObjectType())),
		[]semanticRAITest_TypedName{
			typedName("a", s.getNativeObjectType()), typedName("b", s.getNativeObjectType())},
		[]semanticRAITest_TypedName{
			typedName("a", s.getNativeNullType()), typedName("b", s.getNativeNullType())})
}

// Tests reverse interpretation of a SHEQ(NAME, NUMBER) expression.
func TestSemanticReverseAbstractInterpreter_SheqCondition1(t *testing.T) {
	s := newSemanticRAITest(t)
	blind := s.newScope()
	s.testBinop(
		blind,
		rhino.Token_SHEQ,
		s.createVar(blind, "a", s.createUnionType(s.getNativeStringType(), s.getNativeNumberType())),
		s.createNumber(56),
		[]semanticRAITest_TypedName{typedName("a", s.getNativeNumberType())},
		[]semanticRAITest_TypedName{
			typedName("a", s.createUnionType(s.getNativeStringType(), s.getNativeNumberType()))})
}

// Tests reverse interpretation of a SHEQ(NUMBER, NAME) expression.
func TestSemanticReverseAbstractInterpreter_SheqCondition2(t *testing.T) {
	s := newSemanticRAITest(t)
	blind := s.newScope()
	s.testBinop(
		blind,
		rhino.Token_SHEQ,
		s.createNumber(56),
		s.createVar(blind, "a", s.createUnionType(s.getNativeStringType(), s.getNativeNumberType())),
		[]semanticRAITest_TypedName{typedName("a", s.getNativeNumberType())},
		[]semanticRAITest_TypedName{
			typedName("a", s.createUnionType(s.getNativeStringType(), s.getNativeNumberType()))})
}

// Tests reverse interpretation of a SHEQ(NAME, NAME) expression.
func TestSemanticReverseAbstractInterpreter_SheqCondition3(t *testing.T) {
	s := newSemanticRAITest(t)
	blind := s.newScope()
	s.testBinop(
		blind,
		rhino.Token_SHEQ,
		s.createVar(blind, "b", s.createUnionType(s.getNativeStringType(), s.getNativeBooleanType())),
		s.createVar(blind, "a", s.createUnionType(s.getNativeStringType(), s.getNativeNumberType())),
		[]semanticRAITest_TypedName{
			typedName("a", s.getNativeStringType()), typedName("b", s.getNativeStringType())},
		[]semanticRAITest_TypedName{
			typedName("a", s.createUnionType(s.getNativeStringType(), s.getNativeNumberType())),
			typedName("b", s.createUnionType(s.getNativeStringType(), s.getNativeBooleanType()))})
}

func TestSemanticReverseAbstractInterpreter_SheqCondition4(t *testing.T) {
	s := newSemanticRAITest(t)
	blind := s.newScope()
	s.testBinop(
		blind,
		rhino.Token_SHEQ,
		s.createVar(blind, "a", s.createUnionType(s.getNativeStringType(), s.getNativeVoidType())),
		s.createVar(blind, "b", s.createUnionType(s.getNativeVoidType())),
		[]semanticRAITest_TypedName{
			typedName("a", s.getNativeVoidType()), typedName("b", s.getNativeVoidType())},
		[]semanticRAITest_TypedName{
			typedName("a", s.getNativeStringType()), typedName("b", s.getNativeVoidType())})
}

func TestSemanticReverseAbstractInterpreter_SheqCondition5(t *testing.T) {
	s := newSemanticRAITest(t)
	blind := s.newScope()
	s.testBinop(
		blind,
		rhino.Token_SHEQ,
		s.createVar(blind, "a", s.createUnionType(s.getNativeNullType(), s.getNativeVoidType())),
		s.createVar(blind, "b", s.createUnionType(s.getNativeVoidType())),
		[]semanticRAITest_TypedName{
			typedName("a", s.getNativeVoidType()), typedName("b", s.getNativeVoidType())},
		[]semanticRAITest_TypedName{
			typedName("a", s.getNativeNullType()), typedName("b", s.getNativeVoidType())})
}

func TestSemanticReverseAbstractInterpreter_SheqCondition6(t *testing.T) {
	s := newSemanticRAITest(t)
	blind := s.newScope()
	s.testBinop(
		blind,
		rhino.Token_SHEQ,
		s.createVar(blind, "a", s.createUnionType(s.getNativeStringType(), s.getNativeVoidType())),
		s.createVar(blind, "b", s.createUnionType(s.getNativeNumberType(), s.getNativeVoidType())),
		[]semanticRAITest_TypedName{
			typedName("a", s.getNativeVoidType()), typedName("b", s.getNativeVoidType())},
		[]semanticRAITest_TypedName{
			typedName("a", s.createUnionType(s.getNativeStringType(), s.getNativeVoidType())),
			typedName("b", s.createUnionType(s.getNativeNumberType(), s.getNativeVoidType()))})
}

// Tests reverse interpretation of a SHNE(NAME, NUMBER) expression.
func TestSemanticReverseAbstractInterpreter_ShneCondition1(t *testing.T) {
	s := newSemanticRAITest(t)
	blind := s.newScope()
	s.testBinop(
		blind,
		rhino.Token_SHNE,
		s.createVar(blind, "a", s.createUnionType(s.getNativeStringType(), s.getNativeNumberType())),
		s.createNumber(56),
		[]semanticRAITest_TypedName{
			typedName("a", s.createUnionType(s.getNativeStringType(), s.getNativeNumberType()))},
		[]semanticRAITest_TypedName{typedName("a", s.getNativeNumberType())})
}

// Tests reverse interpretation of a SHNE(NUMBER, NAME) expression.
func TestSemanticReverseAbstractInterpreter_ShneCondition2(t *testing.T) {
	s := newSemanticRAITest(t)
	blind := s.newScope()
	s.testBinop(
		blind,
		rhino.Token_SHNE,
		s.createNumber(56),
		s.createVar(blind, "a", s.createUnionType(s.getNativeStringType(), s.getNativeNumberType())),
		[]semanticRAITest_TypedName{
			typedName("a", s.createUnionType(s.getNativeStringType(), s.getNativeNumberType()))},
		[]semanticRAITest_TypedName{typedName("a", s.getNativeNumberType())})
}

// Tests reverse interpretation of a SHNE(NAME, NAME) expression.
func TestSemanticReverseAbstractInterpreter_ShneCondition3(t *testing.T) {
	s := newSemanticRAITest(t)
	blind := s.newScope()
	s.testBinop(
		blind,
		rhino.Token_SHNE,
		s.createVar(blind, "b", s.createUnionType(s.getNativeStringType(), s.getNativeBooleanType())),
		s.createVar(blind, "a", s.createUnionType(s.getNativeStringType(), s.getNativeNumberType())),
		[]semanticRAITest_TypedName{
			typedName("a", s.createUnionType(s.getNativeStringType(), s.getNativeNumberType())),
			typedName("b", s.createUnionType(s.getNativeStringType(), s.getNativeBooleanType()))},
		[]semanticRAITest_TypedName{
			typedName("a", s.getNativeStringType()), typedName("b", s.getNativeStringType())})
}

func TestSemanticReverseAbstractInterpreter_ShneCondition4(t *testing.T) {
	s := newSemanticRAITest(t)
	blind := s.newScope()
	s.testBinop(
		blind,
		rhino.Token_SHNE,
		s.createVar(blind, "a", s.createUnionType(s.getNativeStringType(), s.getNativeVoidType())),
		s.createVar(blind, "b", s.createUnionType(s.getNativeVoidType())),
		[]semanticRAITest_TypedName{
			typedName("a", s.getNativeStringType()), typedName("b", s.getNativeVoidType())},
		[]semanticRAITest_TypedName{
			typedName("a", s.getNativeVoidType()), typedName("b", s.getNativeVoidType())})
}

func TestSemanticReverseAbstractInterpreter_ShneCondition5(t *testing.T) {
	s := newSemanticRAITest(t)
	blind := s.newScope()
	s.testBinop(
		blind,
		rhino.Token_SHNE,
		s.createVar(blind, "a", s.createUnionType(s.getNativeNullType(), s.getNativeVoidType())),
		s.createVar(blind, "b", s.createUnionType(s.getNativeNullType())),
		[]semanticRAITest_TypedName{
			typedName("a", s.getNativeVoidType()), typedName("b", s.getNativeNullType())},
		[]semanticRAITest_TypedName{
			typedName("a", s.getNativeNullType()), typedName("b", s.getNativeNullType())})
}

func TestSemanticReverseAbstractInterpreter_ShneCondition6(t *testing.T) {
	s := newSemanticRAITest(t)
	blind := s.newScope()
	s.testBinop(
		blind,
		rhino.Token_SHNE,
		s.createVar(blind, "a", s.createUnionType(s.getNativeStringType(), s.getNativeVoidType())),
		s.createVar(blind, "b", s.createUnionType(s.getNativeNumberType(), s.getNativeVoidType())),
		[]semanticRAITest_TypedName{
			typedName("a", s.createUnionType(s.getNativeStringType(), s.getNativeVoidType())),
			typedName("b", s.createUnionType(s.getNativeNumberType(), s.getNativeVoidType()))},
		[]semanticRAITest_TypedName{
			typedName("a", s.getNativeVoidType()), typedName("b", s.getNativeVoidType())})
}

// Tests reverse interpretation of a EQ(NAME, NULL) expression.
func TestSemanticReverseAbstractInterpreter_EqCondition1(t *testing.T) {
	s := newSemanticRAITest(t)
	blind := s.newScope()
	s.testBinop(
		blind,
		rhino.Token_EQ,
		s.createVar(blind, "a", s.createUnionType(s.getNativeBooleanType(), s.getNativeVoidType())),
		s.createNull(),
		[]semanticRAITest_TypedName{typedName("a", s.getNativeVoidType())},
		[]semanticRAITest_TypedName{typedName("a", s.getNativeBooleanType())})
}

// Tests reverse interpretation of a NE(NULL, NAME) expression.
func TestSemanticReverseAbstractInterpreter_EqCondition2(t *testing.T) {
	s := newSemanticRAITest(t)
	blind := s.newScope()
	s.testBinop(
		blind,
		rhino.Token_NE,
		s.createNull(),
		s.createVar(blind, "a", s.createUnionType(s.getNativeBooleanType(), s.getNativeVoidType())),
		[]semanticRAITest_TypedName{typedName("a", s.getNativeBooleanType())},
		[]semanticRAITest_TypedName{typedName("a", s.getNativeVoidType())})
}

// Tests reverse interpretation of a EQ(NAME, NULL) expression.
func TestSemanticReverseAbstractInterpreter_EqCondition3(t *testing.T) {
	s := newSemanticRAITest(t)
	blind := s.newScope()
	// (number,undefined,null)
	nullableOptionalNumber := s.createUnionType(s.getNativeNullType(), s.getNativeVoidType(), s.getNativeNumberType())
	// (null,undefined)
	nullUndefined := s.createUnionType(s.getNativeVoidType(), s.getNativeNullType())
	s.testBinop(
		blind,
		rhino.Token_EQ,
		s.createVar(blind, "a", nullableOptionalNumber),
		s.createNull(),
		[]semanticRAITest_TypedName{typedName("a", nullUndefined)},
		[]semanticRAITest_TypedName{typedName("a", s.getNativeNumberType())})
}

// Tests reverse interpretation of two undefineds.
func TestSemanticReverseAbstractInterpreter_EqCondition4(t *testing.T) {
	s := newSemanticRAITest(t)
	blind := s.newScope()
	s.testBinop(
		blind,
		rhino.Token_EQ,
		s.createVar(blind, "a", s.getNativeVoidType()),
		s.createVar(blind, "b", s.getNativeVoidType()),
		[]semanticRAITest_TypedName{
			typedName("a", s.getNativeVoidType()), typedName("b", s.getNativeVoidType())},
		[]semanticRAITest_TypedName{
			typedName("a", s.getNativeNoType()), typedName("b", s.getNativeNoType())})
}

// Tests reverse interpretation of a COMPARE(NAME, NUMBER) expression, where
// COMPARE can be LE, LT, GE or GT.
func TestSemanticReverseAbstractInterpreter_InequalitiesCondition1(t *testing.T) {
	s := newSemanticRAITest(t)
	for _, op := range []rhino.Token{rhino.Token_LT, rhino.Token_GT, rhino.Token_LE, rhino.Token_GE} {
		blind := s.newScope()
		s.testBinop(
			blind,
			op,
			s.createVar(blind, "a", s.createUnionType(s.getNativeStringType(), s.getNativeVoidType())),
			s.createNumber(8),
			[]semanticRAITest_TypedName{typedName("a", s.getNativeStringType())},
			[]semanticRAITest_TypedName{
				typedName("a", s.createUnionType(s.getNativeStringType(), s.getNativeVoidType()))})
	}
}

// Tests reverse interpretation of a COMPARE(NAME, NAME) expression, where
// COMPARE can be LE, LT, GE or GT.
func TestSemanticReverseAbstractInterpreter_InequalitiesCondition2(t *testing.T) {
	s := newSemanticRAITest(t)
	for _, op := range []rhino.Token{rhino.Token_LT, rhino.Token_GT, rhino.Token_LE, rhino.Token_GE} {
		blind := s.newScope()
		s.testBinop(
			blind,
			op,
			s.createVar(
				blind,
				"a",
				s.createUnionType(s.getNativeStringType(), s.getNativeNumberType(), s.getNativeVoidType())),
			s.createVar(blind, "b", s.createUnionType(s.getNativeNumberType(), s.getNativeNullType())),
			[]semanticRAITest_TypedName{
				typedName("a", s.createUnionType(s.getNativeStringType(), s.getNativeNumberType())),
				typedName("b", s.createUnionType(s.getNativeNumberType(), s.getNativeNullType()))},
			[]semanticRAITest_TypedName{
				typedName(
					"a",
					s.createUnionType(
						s.getNativeStringType(), s.getNativeNumberType(), s.getNativeVoidType())),
				typedName("b", s.createUnionType(s.getNativeNumberType(), s.getNativeNullType()))})
	}
}

// Tests reverse interpretation of a COMPARE(NUMBER-untyped, NAME) expression,
// where COMPARE can be LE, LT, GE or GT.
func TestSemanticReverseAbstractInterpreter_InequalitiesCondition3(t *testing.T) {
	s := newSemanticRAITest(t)
	for _, op := range []rhino.Token{rhino.Token_LT, rhino.Token_GT, rhino.Token_LE, rhino.Token_GE} {
		blind := s.newScope()
		s.testBinop(
			blind,
			op,
			s.createUntypedNumber(8),
			s.createVar(blind, "a", s.createUnionType(s.getNativeStringType(), s.getNativeVoidType())),
			[]semanticRAITest_TypedName{typedName("a", s.getNativeStringType())},
			[]semanticRAITest_TypedName{
				typedName("a", s.createUnionType(s.getNativeStringType(), s.getNativeVoidType()))})
	}
}

func TestSemanticReverseAbstractInterpreter_And(t *testing.T) {
	s := newSemanticRAITest(t)
	blind := s.newScope()
	s.testBinop(
		blind,
		rhino.Token_AND,
		s.createVar(blind, "b", s.createUnionType(s.getNativeStringType(), s.getNativeNullType())),
		s.createVar(blind, "a", s.createUnionType(s.getNativeNumberType(), s.getNativeVoidType())),
		[]semanticRAITest_TypedName{
			typedName("a", s.getNativeNumberType()), typedName("b", s.getNativeStringType())},
		[]semanticRAITest_TypedName{
			typedName("a", s.createUnionType(s.getNativeNumberType(), s.getNativeVoidType())),
			typedName("b", s.createUnionType(s.getNativeStringType(), s.getNativeNullType()))})
}

func TestSemanticReverseAbstractInterpreter_Typeof1(t *testing.T) {
	s := newSemanticRAITest(t)
	blind := s.newScope()
	s.testBinop(
		blind,
		rhino.Token_EQ,
		rhino.NewNodeWithChild(rhino.Token_TYPEOF, s.createVar(blind, "a", s.getNativeObjectType())),
		rhino.Node_NewString("function"),
		[]semanticRAITest_TypedName{typedName("a", s.getNativeFunctionType())},
		[]semanticRAITest_TypedName{typedName("a", s.getNativeObjectType())})
}

func TestSemanticReverseAbstractInterpreter_Typeof2(t *testing.T) {
	s := newSemanticRAITest(t)
	blind := s.newScope()
	s.testBinop(
		blind,
		rhino.Token_EQ,
		rhino.NewNodeWithChild(rhino.Token_TYPEOF, s.createVar(blind, "a", s.getNativeAllType())),
		rhino.Node_NewString("function"),
		[]semanticRAITest_TypedName{typedName("a", s.getNativeFunctionType())},
		[]semanticRAITest_TypedName{typedName("a", s.getNativeAllType())})
}

func TestSemanticReverseAbstractInterpreter_Typeof3(t *testing.T) {
	s := newSemanticRAITest(t)
	blind := s.newScope()
	s.testBinop(
		blind,
		rhino.Token_EQ,
		rhino.NewNodeWithChild(rhino.Token_TYPEOF, s.createVar(blind, "a", s.getNativeObjectNumberStringBooleanType())),
		rhino.Node_NewString("function"),
		[]semanticRAITest_TypedName{typedName("a", s.getNativeFunctionType())},
		[]semanticRAITest_TypedName{typedName("a", s.getNativeObjectNumberStringBooleanType())})
}

func TestSemanticReverseAbstractInterpreter_Typeof4(t *testing.T) {
	s := newSemanticRAITest(t)
	blind := s.newScope()
	s.testBinop(
		blind,
		rhino.Token_EQ,
		rhino.NewNodeWithChild(
			rhino.Token_TYPEOF,
			s.createVar(
				blind,
				"a",
				s.createUnionType(s.getNativeFunctionType(), s.getNativeNumberStringBooleanType()))),
		rhino.Node_NewString("function"),
		[]semanticRAITest_TypedName{typedName("a", s.getNativeFunctionType())},
		[]semanticRAITest_TypedName{typedName("a", s.getNativeNumberStringBooleanType())})
}

func TestSemanticReverseAbstractInterpreter_InstanceOf(t *testing.T) {
	s := newSemanticRAITest(t)
	blind := s.newScope()
	s.testBinop(
		blind,
		rhino.Token_INSTANCEOF,
		s.createVar(blind, "x", s.getNativeUnknownType()),
		s.createVar(blind, "s", s.getNativeStringObjectConstructorType()),
		[]semanticRAITest_TypedName{
			typedName("x", s.getNativeStringObjectType()),
			typedName("s", s.getNativeStringObjectConstructorType())},
		[]semanticRAITest_TypedName{typedName("s", s.getNativeStringObjectConstructorType())})
}

func TestSemanticReverseAbstractInterpreter_InstanceOf2(t *testing.T) {
	s := newSemanticRAITest(t)
	blind := s.newScope()
	s.testBinop(
		blind,
		rhino.Token_INSTANCEOF,
		s.createVar(
			blind, "x", s.createUnionType(s.getNativeStringObjectType(), s.getNativeNumberObjectType())),
		s.createVar(blind, "s", s.getNativeStringObjectConstructorType()),
		[]semanticRAITest_TypedName{
			typedName("x", s.getNativeStringObjectType()),
			typedName("s", s.getNativeStringObjectConstructorType())},
		[]semanticRAITest_TypedName{
			typedName("x", s.getNativeNumberObjectType()),
			typedName("s", s.getNativeStringObjectConstructorType())})
}

func TestSemanticReverseAbstractInterpreter_InstanceOf3(t *testing.T) {
	s := newSemanticRAITest(t)
	blind := s.newScope()
	s.testBinop(
		blind,
		rhino.Token_INSTANCEOF,
		s.createVar(blind, "x", s.getNativeObjectType()),
		s.createVar(blind, "s", s.getNativeStringObjectConstructorType()),
		[]semanticRAITest_TypedName{
			typedName("x", s.getNativeStringObjectType()),
			typedName("s", s.getNativeStringObjectConstructorType())},
		[]semanticRAITest_TypedName{
			typedName("x", s.getNativeObjectType()),
			typedName("s", s.getNativeStringObjectConstructorType())})
}

func TestSemanticReverseAbstractInterpreter_InstanceOf4(t *testing.T) {
	s := newSemanticRAITest(t)
	blind := s.newScope()
	s.testBinop(
		blind,
		rhino.Token_INSTANCEOF,
		s.createVar(blind, "x", s.getNativeAllType()),
		s.createVar(blind, "s", s.getNativeStringObjectConstructorType()),
		[]semanticRAITest_TypedName{
			typedName("x", s.getNativeStringObjectType()),
			typedName("s", s.getNativeStringObjectConstructorType())},
		[]semanticRAITest_TypedName{typedName("s", s.getNativeStringObjectConstructorType())})
}
