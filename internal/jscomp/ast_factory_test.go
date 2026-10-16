// Port of com/google/javascript/jscomp/AstFactoryTest.java @ closure-compiler v20260415

package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/rhinotesting"
)

// The tests that parse source (parseWithoutTypes needs Compiler.init and
// parseInputs; parseAndAddTypes also TypeCheck; parseAndAddColors also
// ConvertTypesToColors) are skipped: those classes are above the core.

// astFactoryTest is AstFactoryTest's fixture (its fields and setUp()).
type astFactoryTest struct {
	t                   *testing.T
	truth               rhinotesting.Truth
	compiler            AbstractCompiler
	runtimeJsLibManager *RuntimeJsLibManager
}

// Java: AstFactoryTest.setUp()
func newAstFactoryTest(t *testing.T) *astFactoryTest {
	x := &astFactoryTest{t: t, truth: rhinotesting.NewTruth(t)}
	x.compiler = newTestCompiler(NewCompilerOptions())
	x.runtimeJsLibManager = RuntimeJsLibManager_Create(
		RuntimeJsLibManager_RuntimeLibraryMode_RECORD_AND_VALIDATE_FIELDS,
		astFactoryTest_TestResourceProvider,
		x.compiler.GetChangeTracker(),
		astFactoryTest_alwaysThrowNodeSupplier)
	return x
}

// Java: AstFactoryTest.alwaysThrowNodeSupplier()
func astFactoryTest_alwaysThrowNodeSupplier() *rhino.Node {
	panic(jutil.NewUnsupportedOperationException())
}

// Java: AstFactoryTest.TestResourceProvider.parse(String, String)
func astFactoryTest_TestResourceProvider(unused1 string, unused2 string) *rhino.Node {
	panic(jutil.NewUnsupportedOperationException())
}

func (x *astFactoryTest) getRegistry() *rhino.JSTypeRegistry {
	return x.compiler.GetTypeRegistry()
}

func (x *astFactoryTest) getNativeType(nativeType rhino.JSTypeNative) rhino.JSType {
	return x.getRegistry().GetNativeType(nativeType)
}

func (x *astFactoryTest) createTestAstFactory() *AstFactory {
	return AstFactory_CreateFactoryWithTypes(
		x.compiler.GetLifeCycleStage(), x.getRegistry(), x.runtimeJsLibManager)
}

func (x *astFactoryTest) createTestAstFactoryWithColors() *AstFactory {
	// the built-in color registry is available only if we've run parseAndAddColors()
	var registry *colors.ColorRegistry
	if x.compiler.HasOptimizationColors() {
		registry = x.compiler.GetColorRegistry()
	} else {
		registry = colors.ColorRegistry_Builder_().SetDefaultNativeColorsForTesting().Build()
	}
	return AstFactory_CreateFactoryWithColors(x.compiler.GetLifeCycleStage(), registry, x.runtimeJsLibManager)
}

func (x *astFactoryTest) createTestAstFactoryWithoutTypes() *AstFactory {
	return AstFactory_CreateFactoryWithoutTypes(x.compiler.GetLifeCycleStage(), x.runtimeJsLibManager)
}

// ---------------------------------------------------------------------------
// The parts of rhino.testing.NodeSubject / ColorSubject the test uses.

type afNodeSubject struct {
	x *astFactoryTest
	n *rhino.Node
}

func (x *astFactoryTest) assertNode(n *rhino.Node) *afNodeSubject {
	x.t.Helper()
	return &afNodeSubject{x: x, n: n}
}

func (s *afNodeSubject) fail(format string, args ...any) {
	s.x.t.Helper()
	s.x.t.Fatalf("node %v: "+format, append([]any{s.n}, args...)...)
}

func (s *afNodeSubject) HasType(token rhino.Token) *afNodeSubject {
	s.x.t.Helper()
	if s.n.GetToken() != token {
		s.fail("expected token %s", token)
	}
	return s
}

func (s *afNodeSubject) HasToken(token rhino.Token) *afNodeSubject {
	s.x.t.Helper()
	return s.HasType(token)
}

func (s *afNodeSubject) IsName(name string) *afNodeSubject {
	s.x.t.Helper()
	if !s.n.IsName() || s.n.GetString() != name {
		s.fail("expected NAME %s", name)
	}
	return s
}

func (s *afNodeSubject) MatchesName(name string) *afNodeSubject {
	s.x.t.Helper()
	if !s.n.MatchesName(name) {
		s.fail("expected to match name %s", name)
	}
	return s
}

func (s *afNodeSubject) MatchesQualifiedName(qname string) *afNodeSubject {
	s.x.t.Helper()
	if !s.n.MatchesQualifiedName(qname) {
		s.fail("expected to match qualified name %s", qname)
	}
	return s
}

func (s *afNodeSubject) IsVar() *afNodeSubject   { s.x.t.Helper(); return s.HasType(rhino.Token_VAR) }
func (s *afNodeSubject) IsConst() *afNodeSubject { s.x.t.Helper(); return s.HasType(rhino.Token_CONST) }
func (s *afNodeSubject) IsFunction() *afNodeSubject {
	s.x.t.Helper()
	return s.HasType(rhino.Token_FUNCTION)
}
func (s *afNodeSubject) IsParamList() *afNodeSubject {
	s.x.t.Helper()
	return s.HasType(rhino.Token_PARAM_LIST)
}
func (s *afNodeSubject) IsCall() *afNodeSubject { s.x.t.Helper(); return s.HasType(rhino.Token_CALL) }

func (s *afNodeSubject) HasXChildren(count int) *afNodeSubject {
	s.x.t.Helper()
	if s.n.GetChildCount() != count {
		s.fail("expected %d children", count)
	}
	return s
}

func (s *afNodeSubject) HasNoChildren() *afNodeSubject { s.x.t.Helper(); return s.HasChildren(false) }

func (s *afNodeSubject) HasChildren(hasChildren bool) *afNodeSubject {
	s.x.t.Helper()
	if s.n.HasChildren() != hasChildren {
		s.fail("expected hasChildren %t", hasChildren)
	}
	return s
}

func (s *afNodeSubject) child(what string, c *rhino.Node) *afNodeSubject {
	s.x.t.Helper()
	if c == nil {
		s.fail("expected a %s child", what)
	}
	return &afNodeSubject{x: s.x, n: c}
}

func (s *afNodeSubject) HasFirstChildThat() *afNodeSubject {
	s.x.t.Helper()
	return s.child("first", s.n.GetFirstChild())
}

func (s *afNodeSubject) HasSecondChildThat() *afNodeSubject {
	s.x.t.Helper()
	return s.child("second", s.n.GetSecondChild())
}

func (s *afNodeSubject) HasLastChildThat() *afNodeSubject {
	s.x.t.Helper()
	return s.child("last", s.n.GetLastChild())
}

func (s *afNodeSubject) HasOneChildThat() *afNodeSubject {
	s.x.t.Helper()
	if !s.n.HasOneChild() {
		s.fail("expected one child")
	}
	return s.child("only", s.n.GetOnlyChild())
}

// IsEqualTo is NodeSubject.isEqualTo(Node): tree equivalence ignoring types.
func (s *afNodeSubject) IsEqualTo(expected *rhino.Node) *afNodeSubject {
	s.x.t.Helper()
	if !s.n.IsEquivalentToWithFlags(expected, false, true, false) {
		s.fail("expected to equal %v", expected)
	}
	return s
}

// IsEquivalentTo is NodeSubject.isEquivalentTo(Node).
func (s *afNodeSubject) IsEquivalentTo(expected *rhino.Node) *afNodeSubject {
	s.x.t.Helper()
	if !s.n.IsEquivalentTo(expected) {
		s.fail("expected to be equivalent to %v", expected)
	}
	return s
}

func (s *afNodeSubject) IsSameInstanceAs(expected *rhino.Node) *afNodeSubject {
	s.x.t.Helper()
	if s.n != expected {
		s.fail("expected the same instance as %v", expected)
	}
	return s
}

func (s *afNodeSubject) IsOptionalChainStart() *afNodeSubject {
	s.x.t.Helper()
	if !s.n.IsOptionalChainStart() {
		s.fail("expected an optional chain start")
	}
	return s
}

func (s *afNodeSubject) IsNotOptionalChainStart() *afNodeSubject {
	s.x.t.Helper()
	if s.n.IsOptionalChainStart() {
		s.fail("expected not an optional chain start")
	}
	return s
}

type afStringSubject struct {
	s      *afNodeSubject
	actual string
}

func (s *afNodeSubject) HasStringThat() *afStringSubject {
	return &afStringSubject{s: s, actual: s.n.GetString()}
}

func (s *afStringSubject) IsEqualTo(expected string) {
	s.s.x.t.Helper()
	if s.actual != expected {
		s.s.fail("string %q, expected %q", s.actual, expected)
	}
}

func (s *afNodeSubject) HasJSTypeThat() *rhinotesting.TypeSubject {
	s.x.t.Helper()
	return s.x.truth.AssertType(s.n.GetJSType())
}

type afColorSubject struct {
	s      *afNodeSubject
	actual *colors.Color
}

func (s *afNodeSubject) HasColorThat() *afColorSubject {
	return &afColorSubject{s: s, actual: s.n.GetColor()}
}

func (c *afColorSubject) IsEqualTo(expected *colors.Color) {
	c.s.x.t.Helper()
	if !jutil.ObjectsEquals(c.actual, expected) {
		c.s.fail("color %v, expected %v", c.actual, expected)
	}
}

// HasAlternates is ColorSubject.hasAlternates(Color...): a union of exactly
// these colors.
func (c *afColorSubject) HasAlternates(alternates ...*colors.Color) {
	c.s.x.t.Helper()
	if c.actual == nil || !c.actual.IsUnion() {
		c.s.fail("expected a union color, got %v", c.actual)
	}
	elements := c.actual.GetUnionElements()
	if elements.Size() != len(alternates) {
		c.s.fail("union %v, expected alternates %v", c.actual, alternates)
	}
	for _, a := range alternates {
		if !elements.Contains(a) {
			c.s.fail("union %v lacks %v", c.actual, a)
		}
	}
}

// astFactoryTest_childList is AstFactoryTest.childList(Node).
func astFactoryTest_childList(parent *rhino.Node) []*rhino.Node {
	var list []*rhino.Node
	for child := parent.GetFirstChild(); child != nil; child = child.GetNext() {
		list = append(list, child)
	}
	return list
}

// assertChildren is assertThat(childList(parent)).containsExactly(nodes)
// (.inOrder() when inOrder); nodes compare by identity, as Node.equals
// does.
func (x *astFactoryTest) assertChildren(parent *rhino.Node, inOrder bool, nodes ...*rhino.Node) {
	x.t.Helper()
	children := astFactoryTest_childList(parent)
	if len(children) != len(nodes) {
		x.t.Fatalf("children of %v: got %d, want %d", parent, len(children), len(nodes))
	}
	for i, n := range nodes {
		if inOrder {
			if children[i] != n {
				x.t.Fatalf("child %d of %v: got %v, want %v", i, parent, children[i], n)
			}
			continue
		}
		found := false
		for _, c := range children {
			found = found || c == n
		}
		if !found {
			x.t.Fatalf("children of %v lack %v", parent, n)
		}
	}
}

func (x *astFactoryTest) assertType(type_ rhino.JSType) *rhinotesting.TypeSubject {
	x.t.Helper()
	return x.truth.AssertType(type_)
}

func (x *astFactoryTest) assertBool(what string, got bool, want bool) {
	x.t.Helper()
	if got != want {
		x.t.Fatalf("%s: got %t, want %t", what, got, want)
	}
}

// ---------------------------------------------------------------------------

func TestAstFactory_testStringLiteral_jstypes(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactory()

	stringLiteral := astFactory.CreateString("hello")
	x.assertNode(stringLiteral).HasType(rhino.Token_STRINGLIT)
	x.truth.AssertThat(stringLiteral.GetString()).IsEqualTo("hello")
	x.assertType(stringLiteral.GetJSType()).IsString()
}

func TestAstFactory_testStringLiteral_colors(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactoryWithColors()

	stringLiteral := astFactory.CreateString("hello")
	x.assertNode(stringLiteral).HasType(rhino.Token_STRINGLIT)
	x.truth.AssertThat(stringLiteral.GetString()).IsEqualTo("hello")
	x.assertNode(stringLiteral).HasColorThat().IsEqualTo(colors.StandardColors_STRING)
}

func TestAstFactory_testNumberLiteral_jstypes(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactory()

	numberLiteral := astFactory.CreateNumber(2112)
	x.assertNode(numberLiteral).HasType(rhino.Token_NUMBER)
	x.truth.AssertThat(numberLiteral.GetDouble()).IsEqualTo(2112.0)
	x.assertType(numberLiteral.GetJSType()).IsNumber()
}

func TestAstFactory_testNumberLiteral_colors(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactoryWithColors()

	numberLiteral := astFactory.CreateNumber(2112)
	x.assertNode(numberLiteral).HasType(rhino.Token_NUMBER)
	x.truth.AssertThat(numberLiteral.GetDouble()).IsEqualTo(2112.0)
	x.assertNode(numberLiteral).HasColorThat().IsEqualTo(colors.StandardColors_NUMBER)
}

func TestAstFactory_testBooleanLiteral_jstypes(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactory()

	trueNode := astFactory.CreateBoolean(true)
	x.assertNode(trueNode).HasType(rhino.Token_TRUE)
	x.assertType(trueNode.GetJSType()).IsBoolean()

	falseNode := astFactory.CreateBoolean(false)
	x.assertNode(falseNode).HasType(rhino.Token_FALSE)
	x.assertType(falseNode.GetJSType()).IsBoolean()
}

func TestAstFactory_testBooleanLiteral_colors(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactoryWithColors()

	trueNode := astFactory.CreateBoolean(true)
	x.assertNode(trueNode).HasType(rhino.Token_TRUE)
	x.assertNode(trueNode).HasColorThat().IsEqualTo(colors.StandardColors_BOOLEAN)
}

func TestAstFactory_testVoidExpression_jstypes(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactory()

	voidNode := astFactory.CreateVoid(astFactory.CreateNumber(0))
	x.assertNode(voidNode).HasType(rhino.Token_VOID)
	x.assertNode(voidNode).HasJSTypeThat().IsVoid()
}

func TestAstFactory_testVoidExpression_colors(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactoryWithColors()

	voidNode := astFactory.CreateVoid(astFactory.CreateNumber(0))
	x.assertNode(voidNode).HasType(rhino.Token_VOID)
	x.assertNode(voidNode).HasColorThat().IsEqualTo(colors.StandardColors_NULL_OR_VOID)
}

func TestAstFactory_testNotExpression_jstypes(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactory()

	notNode := astFactory.CreateNot(astFactory.CreateNumber(0))
	x.assertNode(notNode).HasType(rhino.Token_NOT)
	x.assertNode(notNode).HasJSTypeThat().IsBoolean()
}

func TestAstFactory_testNotExpression_colors(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactoryWithColors()

	notNode := astFactory.CreateNot(astFactory.CreateNumber(0))
	x.assertNode(notNode).HasType(rhino.Token_NOT)
	x.assertNode(notNode).HasColorThat().IsEqualTo(colors.StandardColors_BOOLEAN)
}

func TestAstFactory_testCreateArgumentsReference_jstypes(t *testing.T) {
	t.Skip("needs Compiler, TypeCheck, TestExternsBuilder")
}

func TestAstFactory_testCreateArgumentsReference_colors(t *testing.T) {
	t.Skip("needs Compiler, TypeCheck, ConvertTypesToColors, TestExternsBuilder")
}

func TestAstFactory_testCreateSingleVarNameDeclaration_jstypes(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactory()

	valueNode := astFactory.CreateBoolean(true)
	constNode := astFactory.CreateSingleVarNameDeclarationWithValue("myTrue", valueNode)
	x.assertNode(constNode).IsVar().HasOneChildThat().IsName("myTrue")
	nameNode := constNode.GetOnlyChild()
	x.assertNode(nameNode).HasOneChildThat().IsEqualTo(valueNode)
	x.assertNode(nameNode).HasJSTypeThat().IsEqualTo(valueNode.GetJSType())
}

func TestAstFactory_testCreateSingleVarNameDeclaration_colors(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactoryWithColors()

	valueNode := astFactory.CreateBoolean(true)
	constNode := astFactory.CreateSingleVarNameDeclarationWithValue("myTrue", valueNode)
	x.assertNode(constNode).IsVar().HasOneChildThat().IsName("myTrue")
	nameNode := constNode.GetOnlyChild()
	x.assertNode(nameNode).HasOneChildThat().IsEqualTo(valueNode)
	x.assertNode(nameNode).HasColorThat().IsEqualTo(valueNode.GetColor())
}

func TestAstFactory_testCreateSingleConstNameDeclaration_jstypes(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactory()

	valueNode := astFactory.CreateBoolean(true)
	constNode := astFactory.CreateSingleConstNameDeclaration("myTrue", valueNode)
	x.assertNode(constNode).IsConst().HasOneChildThat().IsName("myTrue")
	nameNode := constNode.GetOnlyChild()
	x.assertNode(nameNode).HasOneChildThat().IsEqualTo(valueNode)
	x.assertNode(nameNode).HasJSTypeThat().IsEqualTo(valueNode.GetJSType())
}

func TestAstFactory_testCreateSingleConstNameDeclaration_colors(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactoryWithColors()

	valueNode := astFactory.CreateBoolean(true)
	constNode := astFactory.CreateSingleConstNameDeclaration("myTrue", valueNode)
	x.assertNode(constNode).IsConst().HasOneChildThat().IsName("myTrue")
	nameNode := constNode.GetOnlyChild()
	x.assertNode(nameNode).HasOneChildThat().IsEqualTo(valueNode)
	x.assertNode(nameNode).HasColorThat().IsEqualTo(valueNode.GetColor())
}

func TestAstFactory_testCreateNameWithJSType(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactory()

	n := astFactory.CreateName("x", AstFactory_Type_OfJSTypeNative(rhino.JSTypeNative_STRING_TYPE))
	x.assertNode(n).HasType(rhino.Token_NAME)
	x.truth.AssertThat(n.GetString()).IsEqualTo("x")
	x.assertType(n.GetJSType()).IsString()
}

func TestAstFactory_testCreateNameWithNativeType(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactory()

	n := astFactory.CreateName("x", AstFactory_Type_OfJSTypeNative(rhino.JSTypeNative_STRING_TYPE))
	x.assertNode(n).HasType(rhino.Token_NAME)
	x.truth.AssertThat(n.GetString()).IsEqualTo("x")
	x.assertType(n.GetJSType()).IsString()
}

func TestAstFactory_testCreateNameWithColor(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactoryWithColors()

	n := astFactory.CreateName("x",
		AstFactory_Type_OfJSTypeNativeAndColor(rhino.JSTypeNative_STRING_TYPE, colors.StandardColors_STRING))
	x.assertNode(n).HasType(rhino.Token_NAME)
	x.truth.AssertThat(n.GetString()).IsEqualTo("x")
	x.assertNode(n).HasColorThat().IsEqualTo(colors.StandardColors_STRING)
}

func TestAstFactory_testCreateNameFromScope_jstypes(t *testing.T) {
	t.Skip("needs Compiler, TypeCheck")
}

func TestAstFactory_testCreateNameFromScope_colors(t *testing.T) {
	t.Skip("needs Compiler, TypeCheck, ConvertTypesToColors")
}

func TestAstFactory_testCreateNameFromScope_crashesIfMissingVariable(t *testing.T) {
	t.Skip("needs Compiler, TypeCheck")
}

func TestAstFactory_testCreateGetPropFromScope_defaultsToUnknownJSTypeWhenNull(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactory()
	receiver := astFactory.CreateNameWithUnknownType("JQ")
	typeTemplate := rhino.IR_Name("JQ")

	jutil.CheckStateMsg(typeTemplate.GetJSType() == nil, "getJSType does not return null ")

	n := astFactory.CreateGetProp(receiver, "$", AstFactory_Type_(typeTemplate))

	x.assertNode(n).HasType(rhino.Token_GETPROP)
	x.assertNode(n).MatchesQualifiedName("JQ.$")
	x.assertNode(n).HasJSTypeThat().IsUnknown()
}

func TestAstFactory_testCreateGetPropFromScope_defaultsToUnknownColorWhenNull(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactoryWithColors()
	receiver := astFactory.CreateNameWithUnknownType("JQ")
	typeTemplate := rhino.IR_Name("JQ")

	jutil.CheckStateMsg(typeTemplate.GetColor() == nil, "getColor does not return null ")

	n := astFactory.CreateGetProp(receiver, "$", AstFactory_Type_(typeTemplate))

	x.assertNode(n).HasType(rhino.Token_GETPROP)
	x.assertNode(n).MatchesQualifiedName("JQ.$")
	x.assertNode(n).HasColorThat().IsEqualTo(colors.StandardColors_UNKNOWN)
}

func TestAstFactory_testCreateThisReference(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactory()

	n := astFactory.CreateThis(AstFactory_Type_OfJSType(x.getNativeType(rhino.JSTypeNative_STRING_TYPE)))
	x.assertNode(n).HasType(rhino.Token_THIS)
	x.assertType(n.GetJSType()).IsString()
}

func TestAstFactory_testCreateSuperReference(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactory()

	n := astFactory.CreateSuper(AstFactory_Type_OfJSType(x.getNativeType(rhino.JSTypeNative_STRING_TYPE)))
	x.assertNode(n).HasType(rhino.Token_SUPER)
	x.assertType(n.GetJSType()).IsString()
}

func TestAstFactory_createThisForEs6ClassMember_jstypes(t *testing.T) {
	t.Skip("needs Compiler, TypeCheck")
}

func TestAstFactory_createThisForEs6ClassMember_colors(t *testing.T) {
	t.Skip("needs Compiler, TypeCheck, ConvertTypesToColors")
}

func TestAstFactory_createThisForEs6ClassStaticMember_jstypes(t *testing.T) {
	t.Skip("needs Compiler, TypeCheck")
}

func TestAstFactory_createThisForEs6ClassStaticMember_colors(t *testing.T) {
	t.Skip("needs Compiler, TypeCheck, ConvertTypesToColors")
}

func TestAstFactory_createThisAliasReferenceForFunction(t *testing.T) {
	t.Skip("needs Compiler, TypeCheck")
}

func TestAstFactory_createThisAliasReferenceForFunctionWithoutTypes(t *testing.T) {
	t.Skip("needs Compiler (init, parseInputs)")
}

func TestAstFactory_testCreateGetpropWithColorFromNode(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactoryWithColors()

	receiver := astFactory.CreateNameWithUnknownType("x")
	typeTemplate := astFactory.CreateNumber(0)

	getProp := astFactory.CreateGetProp(receiver, "y", AstFactory_Type_(typeTemplate))

	x.assertNode(getProp).HasToken(rhino.Token_GETPROP)
	x.assertNode(getProp).HasStringThat().IsEqualTo("y")
	x.assertNode(getProp).HasFirstChildThat().IsEqualTo(receiver)
	x.assertNode(getProp).HasColorThat().IsEqualTo(colors.StandardColors_NUMBER)
}

func TestAstFactory_testCreateStartOptChainGetprop(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactoryWithColors()

	receiver := astFactory.CreateNameWithUnknownType("x")
	typeTemplate := astFactory.CreateNumber(0)

	getProp := astFactory.CreateStartOptChainGetprop(receiver, "y", AstFactory_Type_(typeTemplate))

	x.assertNode(getProp).HasToken(rhino.Token_OPTCHAIN_GETPROP)
	x.assertNode(getProp).HasStringThat().IsEqualTo("y")
	x.assertNode(getProp).HasFirstChildThat().IsEqualTo(receiver)
	x.assertNode(getProp).HasColorThat().IsEqualTo(colors.StandardColors_NUMBER)
	x.assertNode(getProp).IsOptionalChainStart()
}

func TestAstFactory_testCreateContinueOptChainGetprop(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactoryWithColors()

	receiver := astFactory.CreateNameWithUnknownType("x")
	typeTemplate := astFactory.CreateNumber(0)

	getProp := astFactory.CreateContinueOptChainGetprop(receiver, "y", AstFactory_Type_(typeTemplate))

	x.assertNode(getProp).HasToken(rhino.Token_OPTCHAIN_GETPROP)
	x.assertNode(getProp).HasStringThat().IsEqualTo("y")
	x.assertNode(getProp).HasFirstChildThat().IsEqualTo(receiver)
	x.assertNode(getProp).HasColorThat().IsEqualTo(colors.StandardColors_NUMBER)
	x.assertNode(getProp).IsNotOptionalChainStart()
}

func TestAstFactory_testCreateStartOptChainGetelem(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactoryWithColors()

	receiver := astFactory.CreateNameWithUnknownType("x")
	elem := astFactory.CreateNumber(0)
	typeTemplate := astFactory.CreateNumber(0)

	getElem := astFactory.CreateStartOptChainGetelem(receiver, elem, AstFactory_Type_(typeTemplate))

	x.assertNode(getElem).HasToken(rhino.Token_OPTCHAIN_GETELEM)
	x.assertNode(getElem).HasFirstChildThat().IsEqualTo(receiver)
	x.assertNode(getElem).HasSecondChildThat().IsEqualTo(elem)
	x.assertNode(getElem).HasColorThat().IsEqualTo(colors.StandardColors_NUMBER)
	x.assertNode(getElem).IsOptionalChainStart()
}

func TestAstFactory_testCreateContinueOptChainGetelem(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactoryWithColors()

	receiver := astFactory.CreateNameWithUnknownType("x")
	elem := astFactory.CreateNumber(0)
	typeTemplate := astFactory.CreateNumber(0)

	getElem := astFactory.CreateContinueOptChainGetelem(receiver, elem, AstFactory_Type_(typeTemplate))

	x.assertNode(getElem).HasToken(rhino.Token_OPTCHAIN_GETELEM)
	x.assertNode(getElem).HasFirstChildThat().IsEqualTo(receiver)
	x.assertNode(getElem).HasSecondChildThat().IsEqualTo(elem)
	x.assertNode(getElem).HasColorThat().IsEqualTo(colors.StandardColors_NUMBER)
	x.assertNode(getElem).IsNotOptionalChainStart()
}

func TestAstFactory_testCreateStartOptChainCall(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactoryWithColors()

	receiver := astFactory.CreateNameWithUnknownType("x")
	typeTemplate := astFactory.CreateNumber(0)
	arg1 := astFactory.CreateNumber(1)
	arg2 := astFactory.CreateNumber(2)

	call := astFactory.CreateStartOptChainCall(receiver, AstFactory_Type_(typeTemplate), arg1, arg2)

	x.assertNode(call).HasToken(rhino.Token_OPTCHAIN_CALL)
	x.assertNode(call).HasFirstChildThat().IsEqualTo(receiver)
	x.assertNode(call).HasSecondChildThat().IsEqualTo(arg1)
	x.assertNode(call).HasLastChildThat().IsEqualTo(arg2)
	x.assertNode(call).HasColorThat().IsEqualTo(colors.StandardColors_NUMBER)
	x.assertNode(call).IsOptionalChainStart()
}

func TestAstFactory_testCreateContinueOptChainCall(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactoryWithColors()

	receiver := astFactory.CreateNameWithUnknownType("x")
	typeTemplate := astFactory.CreateNumber(0)
	arg1 := astFactory.CreateNumber(1)
	arg2 := astFactory.CreateNumber(2)

	call := astFactory.CreateContinueOptChainCall(receiver, AstFactory_Type_(typeTemplate), arg1, arg2)

	x.assertNode(call).HasToken(rhino.Token_OPTCHAIN_CALL)
	x.assertNode(call).HasFirstChildThat().IsEqualTo(receiver)
	x.assertNode(call).HasSecondChildThat().IsEqualTo(arg1)
	x.assertNode(call).HasLastChildThat().IsEqualTo(arg2)
	x.assertNode(call).HasColorThat().IsEqualTo(colors.StandardColors_NUMBER)
	x.assertNode(call).IsNotOptionalChainStart()
}

func TestAstFactory_testCreateStringKey_jstypes(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactory()

	numberNode := astFactory.CreateNumber(2112)
	stringKeyNode := astFactory.CreateStringKey("key", numberNode)

	x.assertNode(stringKeyNode).HasType(rhino.Token_STRING_KEY)
	x.truth.AssertThat(stringKeyNode.GetString()).IsEqualTo("key")
	x.assertChildren(stringKeyNode, false, numberNode)
	x.assertType(stringKeyNode.GetJSType()).IsNumber()
}

func TestAstFactory_testCreateStringKey_colors(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactoryWithColors()

	numberNode := astFactory.CreateNumber(2112)
	stringKeyNode := astFactory.CreateStringKey("key", numberNode)

	x.assertNode(stringKeyNode).HasType(rhino.Token_STRING_KEY)
	x.truth.AssertThat(stringKeyNode.GetString()).IsEqualTo("key")
	x.assertChildren(stringKeyNode, false, numberNode)
	x.assertNode(stringKeyNode).HasColorThat().IsEqualTo(colors.StandardColors_NUMBER)
}

func TestAstFactory_testCreateComputedProperty_jstypes(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactory()

	numberNode := astFactory.CreateNumber(2112)
	stringLiteral := astFactory.CreateString("string literal key")
	computedPropertyNode := astFactory.CreateComputedProperty(stringLiteral, numberNode)

	x.assertNode(computedPropertyNode).HasType(rhino.Token_COMPUTED_PROP)
	x.assertChildren(computedPropertyNode, true, stringLiteral, numberNode)
	x.assertType(computedPropertyNode.GetJSType()).IsNumber()
}

func TestAstFactory_testCreateComputedProperty_colors(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactoryWithColors()

	numberNode := astFactory.CreateNumber(2112)
	stringLiteral := astFactory.CreateString("string literal key")
	computedPropertyNode := astFactory.CreateComputedProperty(stringLiteral, numberNode)

	x.assertNode(computedPropertyNode).HasType(rhino.Token_COMPUTED_PROP)
	x.assertChildren(computedPropertyNode, true, stringLiteral, numberNode)
	x.assertNode(computedPropertyNode).HasColorThat().IsEqualTo(colors.StandardColors_NUMBER)
}

func TestAstFactory_testCreateGetterDef(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactory()

	valueNode := astFactory.CreateString("value")
	getterNode := astFactory.CreateGetterDef("name", valueNode)

	x.assertNode(getterNode).
		HasToken(rhino.Token_GETTER_DEF).
		HasStringThat().
		IsEqualTo("name")
	x.assertNode(getterNode).
		HasOneChildThat().
		IsFunction().
		HasXChildren(3)
	getterFunctionNode := getterNode.GetOnlyChild()
	// The first child is an empty string representing that the function has no name
	x.assertNode(getterFunctionNode).HasFirstChildThat().IsName("")
	// the second child is the parameter list, which should be empty.
	x.assertNode(getterFunctionNode).
		HasSecondChildThat().
		IsParamList().
		HasNoChildren()
	getterFunctionBody := getterFunctionNode.GetLastChild()
	x.assertNode(getterFunctionBody).
		HasOneChildThat().
		HasToken(rhino.Token_RETURN).
		HasOneChildThat().
		IsEqualTo(valueNode)
}

func TestAstFactory_testCreateGetElem_jstypes(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactory()

	objectName := astFactory.CreateName("obj", AstFactory_Type_OfJSTypeNative(rhino.JSTypeNative_OBJECT_TYPE))
	stringLiteral := astFactory.CreateString("string literal key")
	getElemNode := astFactory.CreateGetElem(objectName, stringLiteral)

	x.assertNode(getElemNode).HasType(rhino.Token_GETELEM)
	x.assertChildren(getElemNode, true, objectName, stringLiteral)
	// TODO(bradfordcsmith): When receiver is an Array<T> or an Object<K, V>, use the template type
	// here.
	x.assertType(getElemNode.GetJSType()).IsUnknown()
}

func TestAstFactory_testCreateGetElem_colors(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactoryWithColors()

	objectName := astFactory.CreateName("obj",
		AstFactory_Type_OfJSTypeNativeAndColor(rhino.JSTypeNative_STRING_TYPE, colors.StandardColors_STRING))
	stringLiteral := astFactory.CreateString("string literal key")
	getElemNode := astFactory.CreateGetElem(objectName, stringLiteral)

	x.assertNode(getElemNode).HasType(rhino.Token_GETELEM)
	x.assertChildren(getElemNode, true, objectName, stringLiteral)
	x.assertNode(getElemNode).HasColorThat().IsEqualTo(colors.StandardColors_UNKNOWN)
}

func TestAstFactory_testCreateComma_jstypes(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactory()

	stringNode := astFactory.CreateString("hi")
	numberNode := astFactory.CreateNumber(2112)
	commaNode := astFactory.CreateComma(stringNode, numberNode)

	x.assertNode(commaNode).HasType(rhino.Token_COMMA)
	x.assertChildren(commaNode, true, stringNode, numberNode)
	x.assertType(commaNode.GetJSType()).IsNumber()
}

func TestAstFactory_testCreateComma_colors(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactoryWithColors()

	stringNode := astFactory.CreateString("hi")
	numberNode := astFactory.CreateNumber(2112)
	commaNode := astFactory.CreateComma(stringNode, numberNode)

	x.assertNode(commaNode).HasType(rhino.Token_COMMA)
	x.assertChildren(commaNode, true, stringNode, numberNode)
	x.assertNode(commaNode).HasColorThat().IsEqualTo(colors.StandardColors_NUMBER)
}

func TestAstFactory_testCreateCommas_jstypes(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactory()

	stringNode := astFactory.CreateString("hi")
	numberNode := astFactory.CreateNumber(2112)
	trueNode := astFactory.CreateBoolean(true)
	falseNode := astFactory.CreateBoolean(false)
	// "hi", 2112, true, false
	stringNumberTrueFalse := astFactory.CreateCommas(stringNode, numberNode, trueNode, falseNode)

	// ("hi", 2112, true), false
	x.assertNode(stringNumberTrueFalse).HasType(rhino.Token_COMMA)
	stringNumberTrue := stringNumberTrueFalse.GetFirstChild()
	x.assertChildren(stringNumberTrueFalse, true, stringNumberTrue, falseNode)
	x.assertType(stringNumberTrueFalse.GetJSType()).IsBoolean()

	// ("hi", 2112), true
	stringNumber := stringNumberTrue.GetFirstChild()
	x.assertChildren(stringNumberTrue, true, stringNumber, trueNode)
	x.assertType(stringNumberTrue.GetJSType()).IsBoolean()

	// "hi", 2112
	x.assertChildren(stringNumber, false, stringNode, numberNode)
	x.assertType(stringNumber.GetJSType()).IsNumber()
}

func TestAstFactory_testCreateIn_jstypes(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactory()

	prop := astFactory.CreateString("prop")
	obj := rhino.IR_Name("obj") // TODO(bradfordcsmith): This should have a type on it.
	n := astFactory.CreateIn(prop, obj)
	x.assertNode(n).HasType(rhino.Token_IN)
	x.assertType(n.GetJSType()).IsBoolean()
	x.assertChildren(n, true, prop, obj)
}

func TestAstFactory_testCreateIn_colors(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactoryWithColors()

	prop := astFactory.CreateString("prop")
	obj := rhino.IR_Name("obj") // TODO(bradfordcsmith): This should have a type on it.
	n := astFactory.CreateIn(prop, obj)
	x.assertNode(n).HasType(rhino.Token_IN)
	x.assertNode(n).HasColorThat().IsEqualTo(colors.StandardColors_BOOLEAN)
	x.assertChildren(n, true, prop, obj)
}

func TestAstFactory_testCreateAnd_jstypes(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactory()

	numberLiteral := astFactory.CreateNumber(2112)
	stringLiteral := astFactory.CreateString("hello")
	andNode := astFactory.CreateAnd(numberLiteral, stringLiteral)

	x.assertNode(andNode).HasType(rhino.Token_AND)
	x.assertChildren(andNode, true, numberLiteral, stringLiteral)
	x.assertType(andNode.GetJSType()).ToStringIsEqualTo("(number|string)")
}

func TestAstFactory_testCreateAnd_colors(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactoryWithColors()

	numberLiteral := astFactory.CreateNumber(2112)
	stringLiteral := astFactory.CreateString("hello")
	andNode := astFactory.CreateAnd(numberLiteral, stringLiteral)

	x.assertNode(andNode).HasType(rhino.Token_AND)
	x.assertChildren(andNode, true, numberLiteral, stringLiteral)
	x.assertNode(andNode).HasColorThat().HasAlternates(colors.StandardColors_STRING, colors.StandardColors_NUMBER)
}

func TestAstFactory_testCreateAndWithAlwaysFalsyLhs(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactory()

	nullNode := astFactory.CreateNull()
	stringLiteral := astFactory.CreateString("hello")
	andNode := astFactory.CreateAnd(nullNode, stringLiteral)

	x.assertNode(andNode).HasType(rhino.Token_AND)
	x.assertChildren(andNode, true, nullNode, stringLiteral)
	// NULL_TYPE doesn't contain any truthy values, so its type is the only possibility
	// but AstFactory is simpler than the type inferencer and does not realize this.
	x.assertType(andNode.GetJSType()).ToStringIsEqualTo("(null|string)")
}

func TestAstFactory_testCreateAndWithAlwaysTruthyLhs(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactory()

	nonNullObject := astFactory.CreateName("nonNullObject", AstFactory_Type_OfJSTypeNative(rhino.JSTypeNative_OBJECT_TYPE))
	stringLiteral := astFactory.CreateString("hello")
	andNode := astFactory.CreateAnd(nonNullObject, stringLiteral)

	x.assertNode(andNode).HasType(rhino.Token_AND)
	x.assertChildren(andNode, true, nonNullObject, stringLiteral)
	// OBJECT_TYPE doesn't contain any falsy values, so the RHS type is the only possibility
	// but AstFactory is simpler than the type inferencer and does not realize this.
	x.assertType(andNode.GetJSType()).ToStringIsEqualTo("(Object|string)")
}

func TestAstFactory_testCreateOr(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactory()

	numberLiteral := astFactory.CreateNumber(2112)
	stringLiteral := astFactory.CreateString("hello")
	andNode := astFactory.CreateOr(numberLiteral, stringLiteral)

	x.assertNode(andNode).HasType(rhino.Token_OR)
	x.assertChildren(andNode, true, numberLiteral, stringLiteral)
	x.assertType(andNode.GetJSType()).ToStringIsEqualTo("(number|string)")
}

func TestAstFactory_testCreateOrWithAlwaysFalsyLhs(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactory()

	nullLiteral := astFactory.CreateNull()
	stringLiteral := astFactory.CreateString("hello")
	andNode := astFactory.CreateOr(nullLiteral, stringLiteral)

	x.assertNode(andNode).HasType(rhino.Token_OR)
	x.assertChildren(andNode, true, nullLiteral, stringLiteral)
	// NULL_TYPE doesn't contain any truthy values, so the RHS type is the only possibility
	// but AstFactory is simpler than the type inferencer and does not realize this.
	x.assertType(andNode.GetJSType()).ToStringIsEqualTo("(null|string)")
}

func TestAstFactory_testCreateOrWithAlwaysTruthyLhs(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactory()

	nonNullObject := astFactory.CreateName("nonNullObject", AstFactory_Type_OfJSTypeNative(rhino.JSTypeNative_OBJECT_TYPE))
	stringLiteral := astFactory.CreateString("hello")
	andNode := astFactory.CreateOr(nonNullObject, stringLiteral)

	x.assertNode(andNode).HasType(rhino.Token_OR)
	x.assertChildren(andNode, true, nonNullObject, stringLiteral)
	// OBJECT_TYPE doesn't contain any falsy values, so the RHS won't be evaluated
	// but AstFactory is simpler than the type inferencer and does not realize this.
	x.assertType(andNode.GetJSType()).ToStringIsEqualTo("(Object|string)")
}

func TestAstFactory_testCreateAdd_stringAndNumber_jstypes(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactory()

	zero := astFactory.CreateNumber(0)
	str := astFactory.CreateString("x")

	add := astFactory.CreateAdd(zero, str)
	x.assertNode(add).HasToken(rhino.Token_ADD)
	x.assertChildren(add, false, zero, str)
	x.assertNode(add).HasJSTypeThat().IsEqualTo(x.getNativeType(rhino.JSTypeNative_BIGINT_NUMBER_STRING))
}

func TestAstFactory_testCreateAdd_stringAndNumber_colors(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactoryWithColors()

	zero := astFactory.CreateNumber(0)
	str := astFactory.CreateString("x")

	add := astFactory.CreateAdd(zero, str)
	x.assertNode(add).HasToken(rhino.Token_ADD)
	x.assertChildren(add, false, zero, str)
	x.assertNode(add).
		HasColorThat().
		HasAlternates(colors.StandardColors_BIGINT, colors.StandardColors_NUMBER, colors.StandardColors_STRING)
}

// astFactoryTest_binary covers the createSub/createLessThan/createBitwiseAnd/
// createRightShift tests, which differ only in method, token and type.
func astFactoryTest_binary(
	t *testing.T, withColors bool, create func(*AstFactory, *rhino.Node, *rhino.Node) *rhino.Node,
	token rhino.Token, nativeType rhino.JSTypeNative, color *colors.Color) {
	x := newAstFactoryTest(t)
	var astFactory *AstFactory
	if withColors {
		astFactory = x.createTestAstFactoryWithColors()
	} else {
		astFactory = x.createTestAstFactory()
	}

	zero := astFactory.CreateNumber(0)
	one := astFactory.CreateNumber(1)

	n := create(astFactory, zero, one)

	x.assertNode(n).HasToken(token)
	x.assertChildren(n, false, zero, one)
	if withColors {
		x.assertNode(n).HasColorThat().IsEqualTo(color)
	} else {
		x.assertNode(n).HasJSTypeThat().IsEqualTo(x.getNativeType(nativeType))
	}
}

func TestAstFactory_testCreateSub_jstypes(t *testing.T) {
	astFactoryTest_binary(t, false, (*AstFactory).CreateSub, rhino.Token_SUB, rhino.JSTypeNative_NUMBER_TYPE, nil)
}

func TestAstFactory_testCreateSub_colors(t *testing.T) {
	astFactoryTest_binary(t, true, (*AstFactory).CreateSub, rhino.Token_SUB, 0, colors.StandardColors_NUMBER)
}

func TestAstFactory_testCreateLessThan_jstypes(t *testing.T) {
	astFactoryTest_binary(t, false, (*AstFactory).CreateLessThan, rhino.Token_LT, rhino.JSTypeNative_BOOLEAN_TYPE, nil)
}

func TestAstFactory_testCreateLessThan_colors(t *testing.T) {
	astFactoryTest_binary(t, true, (*AstFactory).CreateLessThan, rhino.Token_LT, 0, colors.StandardColors_BOOLEAN)
}

func TestAstFactory_testCreateBitwiseAnd_jstypes(t *testing.T) {
	astFactoryTest_binary(t, false, (*AstFactory).CreateBitwiseAnd, rhino.Token_BITAND, rhino.JSTypeNative_NUMBER_TYPE, nil)
}

func TestAstFactory_testCreateBitwiseAnd_colors(t *testing.T) {
	astFactoryTest_binary(t, true, (*AstFactory).CreateBitwiseAnd, rhino.Token_BITAND, 0, colors.StandardColors_NUMBER)
}

func TestAstFactory_testCreateRightShift_jstypes(t *testing.T) {
	astFactoryTest_binary(t, false, (*AstFactory).CreateRightShift, rhino.Token_RSH, rhino.JSTypeNative_NUMBER_TYPE, nil)
}

func TestAstFactory_testCreateRightShift_colors(t *testing.T) {
	astFactoryTest_binary(t, true, (*AstFactory).CreateRightShift, rhino.Token_RSH, 0, colors.StandardColors_NUMBER)
}

func TestAstFactory_testCreateInc_prefix_jstypes(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactory()

	n := astFactory.CreateNameWithUnknownType("x")

	inc := astFactory.CreateInc(n /* isPost= */, false)

	x.assertNode(inc).HasToken(rhino.Token_INC)
	x.assertBool("INCRDECR_PROP", inc.GetBooleanProp(rhino.Node_INCRDECR_PROP), false)
	x.assertChildren(inc, false, n)
	x.assertNode(inc).HasJSTypeThat().IsEqualTo(x.getNativeType(rhino.JSTypeNative_NUMBER_TYPE))
}

func TestAstFactory_testCreateInc_postfix_jstypes(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactory()

	n := astFactory.CreateNameWithUnknownType("x")

	inc := astFactory.CreateInc(n /* isPost= */, true)

	x.assertNode(inc).HasToken(rhino.Token_INC)
	x.assertBool("INCRDECR_PROP", inc.GetBooleanProp(rhino.Node_INCRDECR_PROP), true)
	x.assertChildren(inc, false, n)
	x.assertNode(inc).HasJSTypeThat().IsEqualTo(x.getNativeType(rhino.JSTypeNative_NUMBER_TYPE))
}

func TestAstFactory_testCreateInc_prefix_colors(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactoryWithColors()

	n := astFactory.CreateNameWithUnknownType("x")

	inc := astFactory.CreateInc(n /* isPost= */, false)

	x.assertNode(inc).HasToken(rhino.Token_INC)
	x.assertBool("INCRDECR_PROP", inc.GetBooleanProp(rhino.Node_INCRDECR_PROP), false)
	x.assertChildren(inc, false, n)
	x.assertNode(inc).HasColorThat().IsEqualTo(colors.StandardColors_NUMBER)
}

func TestAstFactory_testCreateCallWithTypeFromNode(t *testing.T) {
	t.Skip("needs Compiler, TypeCheck, SyntacticScopeCreator over a type-checked AST")
}

func TestAstFactory_testCreateCallWithColorFromNode(t *testing.T) {
	t.Skip("needs Compiler, TypeCheck, ConvertTypesToColors")
}

func TestAstFactory_testCreateFreeCall(t *testing.T) {
	t.Skip("needs Compiler, TypeCheck")
}

func TestAstFactory_testCreateMethodCall_throws(t *testing.T) {
	t.Skip("needs Compiler, TypeCheck, Compiler.getTranspilationNamespace")
}

func TestAstFactory_testCreateStaticMethodCallDotCallThrows(t *testing.T) {
	t.Skip("needs Compiler, TypeCheck, Compiler.getTranspilationNamespace")
}

func TestAstFactory_testCreateQNameWithUnknownTypeFromString_jstype(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactory()

	objDotInnerDotStr := astFactory.CreateQNameWithUnknownType("obj.inner.str")

	x.assertNode(objDotInnerDotStr).MatchesQualifiedName("obj.inner.str")
	objDotInner := objDotInnerDotStr.GetFirstChild()
	obj := objDotInner.GetFirstChild()

	x.assertNode(obj).HasJSTypeThat().IsUnknown()
	x.assertNode(objDotInner).HasJSTypeThat().IsUnknown()
	x.assertNode(objDotInnerDotStr).HasJSTypeThat().IsUnknown()
}

func TestAstFactory_testCreateQNameWithUnknownTypeFromString_colors(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactoryWithColors()

	objDotInnerDotStr := astFactory.CreateQNameWithUnknownType("obj.inner.str")

	x.assertNode(objDotInnerDotStr).MatchesQualifiedName("obj.inner.str")
	objDotInner := objDotInnerDotStr.GetFirstChild()
	obj := objDotInner.GetFirstChild()

	x.assertNode(obj).HasColorThat().IsEqualTo(colors.StandardColors_UNKNOWN)
	x.assertNode(objDotInner).HasColorThat().IsEqualTo(colors.StandardColors_UNKNOWN)
	x.assertNode(objDotInnerDotStr).HasColorThat().IsEqualTo(colors.StandardColors_UNKNOWN)
}

func TestAstFactory_testCreateQNameFromString(t *testing.T) {
	t.Skip("needs Compiler, TypeCheck, Compiler.getTranspilationNamespace")
}

func TestAstFactory_testCreateQNameFromBaseNamePlusStringIterable(t *testing.T) {
	t.Skip("needs Compiler, TypeCheck, Compiler.getTranspilationNamespace")
}

func TestAstFactory_testCreateNameMaintainsNormalization(t *testing.T) {
	t.Skip("needs Compiler, TypeCheck, Compiler.getTranspilationNamespace")
}

func TestAstFactory_testCreateQNameFromStringVarArgs(t *testing.T) {
	t.Skip("needs Compiler, TypeCheck, Compiler.getTranspilationNamespace")
}

func TestAstFactory_testCreateQNameFromSimpleStringAndTypedScope(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactory()

	scope := TypedScope_CreateGlobalScope(rhino.IR_Root())
	scope.Declare("x", rhino.IR_Name("x"), x.getNativeType(rhino.JSTypeNative_NUMBER_TYPE), nil, true)

	name := astFactory.CreateQNameUsingJSTypeInfo(scope, "x")

	x.assertNode(name).HasStringThat().IsEqualTo("x")
	x.assertNode(name).HasJSTypeThat().IsNumber()
}

func TestAstFactory_testCreateQNameFromDottedStringAndTypedScope(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactory()

	scope := TypedScope_CreateGlobalScope(rhino.IR_Root())
	// Declare a global "x" with the type "{y: number}".
	objectWithYProp := x.getRegistry().CreateAnonymousObjectType(nil)
	objectWithYProp.DefineDeclaredProperty("y", x.getNativeType(rhino.JSTypeNative_NUMBER_TYPE), nil)
	scope.Declare("x", rhino.IR_Name("x"), objectWithYProp, nil, true)

	name := astFactory.CreateQNameUsingJSTypeInfo(scope, "x.y")

	x.assertNode(name).MatchesQualifiedName("x.y")
	x.assertNode(name).HasJSTypeThat().IsNumber()
	x.assertNode(name.GetFirstChild()).HasJSTypeThat().IsObjectTypeWithProperty("y")
}

func TestAstFactory_testCreateQNameFromStringAndTypedScope_crashesGivenMissingName(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactory()

	scope := TypedScope_CreateGlobalScope(rhino.IR_Root())

	jsaThrows[jutil.Throwable](t, func() { astFactory.CreateQNameUsingJSTypeInfo(scope, "x") })
}

func TestAstFactory_testCreateQNameFromStringAndTypedScope_crashesGivenLocalScope(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactory()

	root := rhino.IR_Root()
	block := rhino.IR_Block()
	root.AddChildToFront(rhino.IR_ScriptOfNodes(block))

	globalScope := TypedScope_CreateGlobalScope(root)
	globalScope.Declare("x", rhino.IR_Name("x"), x.getNativeType(rhino.JSTypeNative_NUMBER_TYPE), nil, true)
	localScope := NewTypedScope(globalScope, block)

	_ = astFactory.CreateQNameUsingJSTypeInfo(globalScope, "x")
	jsaThrows[*jutil.IllegalArgumentException](t, func() {
		astFactory.CreateQNameUsingJSTypeInfo(localScope, "x")
	})
}

func TestAstFactory_testCreateConstructorCall(t *testing.T) {
	t.Skip("needs Compiler, TypeCheck")
}

func TestAstFactory_testCreateEmptyFunction(t *testing.T) {
	t.Skip("needs Compiler, TypeCheck")
}

func TestAstFactory_testCreateFunction(t *testing.T) {
	t.Skip("needs Compiler, TypeCheck")
}

func TestAstFactory_testCreateMemberFunctionDef_jstypes(t *testing.T) {
	t.Skip("needs Compiler, TypeCheck")
}

func TestAstFactory_testCreateMemberFunctionDef_colors(t *testing.T) {
	t.Skip("needs Compiler, TypeCheck, ConvertTypesToColors")
}

func TestAstFactory_testCreateZeroArgFunction(t *testing.T) {
	t.Skip("needs Compiler, TypeCheck")
}

func TestAstFactory_createZeroArgGeneratorFunction(t *testing.T) {
	t.Skip("needs Compiler, TypeCheck")
}

func TestAstFactory_testCreateZeroArgFunctionForExpression(t *testing.T) {
	t.Skip("needs Compiler, TypeCheck")
}

func TestAstFactory_testCreateAssignFromNodes_jstypes(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactory()

	lhs := astFactory.CreateName("x", AstFactory_Type_OfJSTypeNative(rhino.JSTypeNative_STRING_TYPE))
	rhs := astFactory.CreateNumber(0)

	assign := astFactory.CreateAssign(lhs, rhs)
	x.assertNode(assign).HasToken(rhino.Token_ASSIGN)
	x.assertNode(assign).HasFirstChildThat().IsEqualTo(lhs)
	x.assertNode(assign).HasSecondChildThat().IsEqualTo(rhs)
	x.assertNode(assign).HasJSTypeThat().IsNumber() // take the rhs type, not lhs type
}

func TestAstFactory_testCreateAssignFromNodes_colors(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactoryWithColors()

	lhs := astFactory.CreateName("x", AstFactory_Type_OfColor(colors.StandardColors_STRING))
	rhs := astFactory.CreateNumber(0)

	assign := astFactory.CreateAssign(lhs, rhs)
	x.assertNode(assign).HasToken(rhino.Token_ASSIGN)
	x.assertNode(assign).HasFirstChildThat().IsEqualTo(lhs)
	x.assertNode(assign).HasSecondChildThat().IsEqualTo(rhs)
	x.assertNode(assign).HasColorThat().IsEqualTo(colors.StandardColors_NUMBER) // rhs type, not lhs type
}

func TestAstFactory_testCreateObjectLit_empty(t *testing.T) {
	t.Skip("needs Compiler, TypeCheck")
}

func TestAstFactory_testCreateObjectLit_empty_withColor(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactoryWithColors()

	objColor := colors.Color_SingleBuilder().SetId(colors.ColorId_FromAscii("1")).Build()

	objectLit := astFactory.CreateObjectLitWithType(
		AstFactory_Type_OfJSTypeNativeAndColor(rhino.JSTypeNative_UNKNOWN_TYPE, objColor))

	x.assertNode(objectLit).HasColorThat().IsEqualTo(objColor)
	x.assertNode(objectLit).HasToken(rhino.Token_OBJECTLIT)
	x.assertNode(objectLit).HasChildren(false)
}

func TestAstFactory_testCreateObjectLit_withElements(t *testing.T) {
	t.Skip("needs Compiler, TypeCheck")
}

func TestAstFactory_testCreateObjectLit_withElementsAndType(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactory()

	// Try creating an object literal that implements Thenable and type it as one.
	thenableType := x.getRegistry().GetNativeType(rhino.JSTypeNative_THENABLE_TYPE)
	thenStringKey :=
		rhino.IR_StringKeyWithValue(
			"then",
			astFactory.CreateZeroArgArrowFunctionForExpression(astFactory.CreateString("hi")))

	objectLit := astFactory.CreateObjectLitWithType(AstFactory_Type_OfJSType(thenableType), thenStringKey)

	x.truth.AssertThat(objectLit.GetJSType()).IsEqualTo(thenableType)
	x.assertNode(objectLit).HasToken(rhino.Token_OBJECTLIT)
	x.assertNode(objectLit).HasOneChildThat().IsSameInstanceAs(thenStringKey)
}

func TestAstFactory_testCreateObjectLit_empty_colors(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactoryWithColors()

	objectLit := astFactory.CreateObjectLit()

	x.assertNode(objectLit).HasColorThat().IsEqualTo(colors.StandardColors_TOP_OBJECT)
	x.assertNode(objectLit).HasToken(rhino.Token_OBJECTLIT)
	x.assertNode(objectLit).HasChildren(false)
}

func TestAstFactory_testCreateDelProp_jstypes(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactory()

	getprop := rhino.IR_Getprop(rhino.IR_Name("obj"), "prop")

	delprop := astFactory.CreateDelProp(getprop)
	x.assertNode(delprop).HasToken(rhino.Token_DELPROP)
	x.assertType(delprop.GetJSType()).IsEqualTo(x.getNativeType(rhino.JSTypeNative_BOOLEAN_TYPE))
	x.assertNode(delprop).HasChildren(true)
}

func TestAstFactory_testCreateDelProp_colors(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactoryWithColors()

	getprop := rhino.IR_Getprop(rhino.IR_Name("obj"), "prop")

	delprop := astFactory.CreateDelProp(getprop)
	x.assertNode(delprop).HasToken(rhino.Token_DELPROP)
	x.assertNode(delprop).HasColorThat().IsEqualTo(colors.StandardColors_BOOLEAN)
	x.assertNode(delprop).HasChildren(true)
}

// astFactoryTest_comparison covers the createSheq/createEq/createNe tests.
func astFactoryTest_comparison(
	t *testing.T, withColors bool, create func(*AstFactory, *rhino.Node, *rhino.Node) *rhino.Node, token rhino.Token) {
	x := newAstFactoryTest(t)
	var astFactory *AstFactory
	if withColors {
		astFactory = x.createTestAstFactoryWithColors()
	} else {
		astFactory = x.createTestAstFactory()
	}

	left := rhino.IR_String("left")
	right := rhino.IR_Number(0)

	n := create(astFactory, left, right)
	x.assertNode(n).HasToken(token)
	if withColors {
		x.assertNode(n).HasColorThat().IsEqualTo(colors.StandardColors_BOOLEAN)
	} else {
		x.assertType(n.GetJSType()).IsEqualTo(x.getNativeType(rhino.JSTypeNative_BOOLEAN_TYPE))
	}
}

func TestAstFactory_testCreateSheq_jstypes(t *testing.T) {
	astFactoryTest_comparison(t, false, (*AstFactory).CreateSheq, rhino.Token_SHEQ)
}

func TestAstFactory_testCreateSheq_colors(t *testing.T) {
	astFactoryTest_comparison(t, true, (*AstFactory).CreateSheq, rhino.Token_SHEQ)
}

func TestAstFactory_testCreateEq_jstypes(t *testing.T) {
	astFactoryTest_comparison(t, false, (*AstFactory).CreateEq, rhino.Token_EQ)
}

func TestAstFactory_testCreateEq_colors(t *testing.T) {
	astFactoryTest_comparison(t, true, (*AstFactory).CreateEq, rhino.Token_EQ)
}

func TestAstFactory_testCreateNe_jstypes(t *testing.T) {
	astFactoryTest_comparison(t, false, (*AstFactory).CreateNe, rhino.Token_NE)
}

func TestAstFactory_testCreateNe_colors(t *testing.T) {
	astFactoryTest_comparison(t, true, (*AstFactory).CreateNe, rhino.Token_NE)
}

func TestAstFactory_testCreateHook_jstypes(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactory()
	stringType := x.getNativeType(rhino.JSTypeNative_STRING_TYPE)
	numberType := x.getNativeType(rhino.JSTypeNative_NUMBER_TYPE)

	condition := rhino.IR_FalseNode()
	left := rhino.IR_Name("left").SetJSType(stringType)
	right := rhino.IR_Number(0).SetJSType(numberType)

	hook := astFactory.CreateHook(condition, left, right)
	x.assertNode(hook).HasToken(rhino.Token_HOOK)
	x.assertType(hook.GetJSType()).IsEqualTo(x.getRegistry().CreateUnionType(stringType, numberType))
}

func TestAstFactory_testCreateHook_colors(t *testing.T) {
	x := newAstFactoryTest(t)
	astFactory := x.createTestAstFactoryWithColors()

	condition := rhino.IR_FalseNode()
	left := astFactory.CreateString("left")
	right := astFactory.CreateNumber(0)

	hook := astFactory.CreateHook(condition, left, right)
	x.assertNode(hook).HasToken(rhino.Token_HOOK)
	x.assertNode(hook).HasColorThat().HasAlternates(colors.StandardColors_STRING, colors.StandardColors_NUMBER)
}

func TestAstFactory_testCreateArraylit_jstypes(t *testing.T) {
	t.Skip("needs Compiler, TypeCheck")
}

func TestAstFactory_testCreateArraylit_colors(t *testing.T) {
	t.Skip("needs Compiler, TypeCheck, ConvertTypesToColors")
}

func TestAstFactory_testCreateNewNode_jstypes(t *testing.T) {
	t.Skip("needs Compiler, TypeCheck")
}

func TestAstFactory_testCreateNewNode_colors(t *testing.T) {
	t.Skip("needs Compiler, TypeCheck, ConvertTypesToColors")
}

func TestAstFactory_testCreateRuntimeField_throwsIfFieldNotInjected(t *testing.T) {
	x := newAstFactoryTest(t)
	// Given
	astFactory := x.createTestAstFactoryWithoutTypes()
	scope := rhinotesting.NewMapBasedScope(jutil.NewLinkedHashMap[string, rhino.JSType]())
	field := x.runtimeJsLibManager.GetJsLibField("$jscomp.global")

	// When
	jsaThrows[*jutil.IllegalStateException](t, func() { astFactory.CreateQNameOfField(scope, field) })
}

func TestAstFactory_testCreateRuntimeField_succeeds(t *testing.T) {
	x := newAstFactoryTest(t)
	// Given
	astFactory := x.createTestAstFactoryWithoutTypes()
	scope := rhinotesting.NewMapBasedScope(jutil.NewLinkedHashMap[string, rhino.JSType]())
	field := x.runtimeJsLibManager.GetJsLibField("$jscomp.global")
	x.runtimeJsLibManager.InjectLibForField("$jscomp.global")

	// When
	result := astFactory.CreateQNameOfField(scope, field)

	// Then
	x.assertNode(result).MatchesQualifiedName("$jscomp.global")
}

func TestAstFactory_testCreateRuntimeField_inExternMode_succeeds_andUsesUnderscoredName(t *testing.T) {
	x := newAstFactoryTest(t)
	// Given
	externs := rhino.IR_Script()
	externs.SetStaticSourceFile(SourceFile_FromCodeWithKind("externs.js", "", rhino.StaticSourceFile_SourceKind_EXTERN))
	x.runtimeJsLibManager = RuntimeJsLibManager_Create(
		RuntimeJsLibManager_RuntimeLibraryMode_EXTERN_FIELD_NAMES,
		astFactoryTest_TestResourceProvider,
		x.compiler.GetChangeTracker(),
		func() *rhino.Node { return externs })
	astFactory := x.createTestAstFactoryWithoutTypes()
	scope := rhinotesting.NewMapBasedScope(jutil.NewLinkedHashMap[string, rhino.JSType]())
	field := x.runtimeJsLibManager.GetJsLibField("$jscomp.global")
	x.runtimeJsLibManager.InjectLibForField("$jscomp.global")

	// When
	result := astFactory.CreateQNameOfField(scope, field)

	// Then
	x.assertNode(result).MatchesName("$jscomp_global")
}

func TestAstFactory_testCreateJscompMakeIteratorCall_throwsIfJscompMakeIteratorNotInjected(t *testing.T) {
	x := newAstFactoryTest(t)
	// Given
	astFactory := x.createTestAstFactoryWithoutTypes()
	iterable := rhino.IR_Name("arr")
	scope := rhinotesting.NewMapBasedScope(jutil.NewLinkedHashMap[string, rhino.JSType]())

	// When
	jsaThrows[*jutil.IllegalStateException](t, func() { astFactory.CreateJSCompMakeIteratorCall(iterable, scope) })
}

func TestAstFactory_testCreateJscompMakeIteratorCall_succeeds(t *testing.T) {
	x := newAstFactoryTest(t)
	// Given
	astFactory := x.createTestAstFactoryWithoutTypes()
	iterable := rhino.IR_Name("arr")
	x.runtimeJsLibManager.InjectLibForField("$jscomp.makeIterator")

	// When
	result := astFactory.CreateJSCompMakeIteratorCall(iterable, rhinotesting.MapBasedScope_EmptyScope())

	// Then
	x.assertNode(result).IsCall()
	x.assertNode(result).HasFirstChildThat().MatchesQualifiedName("$jscomp.makeIterator")
	x.assertNode(result).HasLastChildThat().IsEqualTo(iterable)
}
