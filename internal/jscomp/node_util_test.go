// Port of com/google/javascript/jscomp/NodeUtilTest.java @ closure-compiler v20261021

/*
 * Copyright 2004 The Closure Compiler Authors.
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

// Tests for NodeUtil.
//
// The JUnit test classes nested in NodeUtilTest become name prefixes
// (TestNodeUtil_<Class>_<Method>); a @Parameterized class becomes one Go test
// with a subtest per case, named like the Java @Parameters name.

import (
	"fmt"
	"math"
	"math/big"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/rhinotesting"
)

// nuAssertThat is Truth.assertThat(actual).
func nuAssertThat(t *testing.T, actual any) *rhinotesting.Subject {
	return rhinotesting.NewTruth(t).AssertThat(actual)
}

// nuAssertWithMessage is Truth.assertWithMessage(format, args...).
func nuAssertWithMessage(t *testing.T, format string, args ...any) *rhinotesting.StandardSubjectBuilder {
	return rhinotesting.NewTruth(t).AssertWithMessage(format, args...)
}

// nuNodeSubject is the part of com.google.javascript.rhino.testing.NodeSubject
// that NodeUtilTest uses.
type nuNodeSubject struct {
	t      *testing.T
	actual rhino.Node
}

// Java: NodeSubject.assertNode(Node)
func nuAssertNode(t *testing.T, actual rhino.Node) *nuNodeSubject {
	return &nuNodeSubject{t: t, actual: actual}
}

// IsEqualTo compares the node-trees of actual and expected, ignoring types,
// JSDoc and side-effects.
//
// Java: NodeSubject.isEqualTo(Node)
func (s *nuNodeSubject) IsEqualTo(expected rhino.Node) {
	s.t.Helper()
	if s.actual.IsNil() || expected.IsNil() {
		s.t.Fatalf("expected not to be null: actual %v, expected %v", s.actual, expected)
	}
	if m := nuNodeSubject_findFirstMismatch(s.actual, expected); m != nil {
		s.t.Fatalf("Node tree inequality\nAST diff\n  expected:\n%s\n  actual:\n%s",
			m[1].ToStringTree(), m[0].ToStringTree())
	}
}

// Java: NodeSubject.findFirstMismatch(Node, Node, boolean)
func nuNodeSubject_findFirstMismatch(actual, expected rhino.Node) []rhino.Node {
	if !actual.IsEquivalentToWithFlagsAndSideEffect(
		expected, false /* compareType */, false /* recurse */, false /* jsDoc */, false /* sideEffect */) {
		return []rhino.Node{actual, expected}
	}
	actualChild := actual.GetFirstChild()
	expectedChild := expected.GetFirstChild()
	for !actualChild.IsNil() {
		if mismatch := nuNodeSubject_findFirstMismatch(actualChild, expectedChild); mismatch != nil {
			return mismatch
		}
		actualChild = actualChild.GetNext()
		expectedChild = expectedChild.GetNext()
	}
	actualShadowAst := actual.GetClosureUnawareShadow()
	expectedShadowAst := expected.GetClosureUnawareShadow()
	if !actualShadowAst.IsNil() && !expectedShadowAst.IsNil() {
		if mismatch := nuNodeSubject_findFirstMismatch(actualShadowAst, expectedShadowAst); mismatch != nil {
			return mismatch
		}
	}
	return nil
}

func (s *nuNodeSubject) check(what string, ok bool) *nuNodeSubject {
	s.t.Helper()
	if !ok {
		s.t.Fatalf("value of: %s\nfor node: %v", what, s.actual)
	}
	return s
}

// Java: NodeSubject.hasToken(Token)
func (s *nuNodeSubject) HasToken(token rhino.Token) *nuNodeSubject {
	s.t.Helper()
	return s.check(fmt.Sprintf("getToken() = %v (was %v)", token, s.actual.GetToken()), s.actual.GetToken() == token)
}

// Java: NodeSubject.hasType(Token)
func (s *nuNodeSubject) HasType(token rhino.Token) *nuNodeSubject {
	s.t.Helper()
	return s.HasToken(token)
}

// Java: NodeSubject.isName(String)
func (s *nuNodeSubject) IsName(name string) *nuNodeSubject {
	s.t.Helper()
	s.check("isName()", s.actual.IsName())
	return s.check(fmt.Sprintf("getString() = %q (was %q)", name, s.actual.GetString()), s.actual.GetString() == name)
}

// Java: NodeSubject.isMemberFunctionDef(String)
func (s *nuNodeSubject) IsMemberFunctionDef(name string) *nuNodeSubject {
	s.t.Helper()
	s.check("isMemberFunctionDef()", s.actual.IsMemberFunctionDef())
	return s.check(fmt.Sprintf("getString() = %q (was %q)", name, s.actual.GetString()), s.actual.GetString() == name)
}

// Java: NodeSubject.matchesQualifiedName(String)
func (s *nuNodeSubject) MatchesQualifiedName(qname string) *nuNodeSubject {
	s.t.Helper()
	return s.check(fmt.Sprintf("matchesQualifiedName(%q)", qname), s.actual.MatchesQualifiedName(qname))
}

// nuParse is NodeUtilTest.parse (ParseHelper.parse): a Compiler with
// setStrictModeInput(false), the ES5_STRICT group off and languageIn
// UNSUPPORTED parses js with parseTestCode; it must report no errors.
func nuParse(t *testing.T, js string) rhino.Node {
	t.Helper()
	n, _ := nuParseWithCompiler(t, js)
	return n
}

// nuParseWithCompiler is ParseHelper.parse, also returning ParseHelper.compiler.
func nuParseWithCompiler(t *testing.T, js string) (rhino.Node, *Compiler) {
	t.Helper()
	options := NewCompilerOptions()

	// To allow octal literals such as 0123 to be parsed.
	options.SetStrictModeInput(false)
	options.SetWarningLevel(DiagnosticGroups_ES5_STRICT, CheckLevel_OFF)

	options.SetLanguageIn(CompilerOptions_LanguageMode_UNSUPPORTED)

	compiler := NewCompiler()
	compiler.InitOptions(options)
	n := compiler.ParseTestCode(js)
	nuAssertThat(t, compiler.GetErrors()).IsEmpty()
	return n, compiler
}

// nuParseFirstWithCompiler is ParseHelper.parseFirst, also returning
// ParseHelper.compiler (for ParseHelper.getAstAnalyzer).
func nuParseFirstWithCompiler(t *testing.T, token rhino.Token, js string) (rhino.Node, *Compiler) {
	t.Helper()
	rootNode, compiler := nuParseWithCompiler(t, js)
	jutil.CheckStateMsg1(rootNode.IsScript(), "%s", rootNode)
	if token == rhino.Token_SCRIPT {
		return rootNode, compiler
	}
	return nuGetNode(t, rootNode, token), compiler
}

// Java: NodeUtilTest.ParseHelper.parseFirst(Token, String)
func nuParseFirst(t *testing.T, token rhino.Token, js string) rhino.Node {
	t.Helper()
	rootNode := nuParse(t, js)
	jutil.CheckStateMsg1(rootNode.IsScript(), "%s", rootNode)
	if token == rhino.Token_SCRIPT {
		return rootNode
	}
	return nuGetNode(t, rootNode, token)
}

// Java: NodeUtilTest.ParseHelper.parseSecond(Token, String)
func nuParseSecond(t *testing.T, token rhino.Token, js string) rhino.Node {
	t.Helper()
	return nuGetNode(t, nuParseFirst(t, token, js), token)
}

// nuParseExpr returns the parsed expression (e.g. returns a NAME given 'a').
//
// Java: NodeUtilTest.parseExpr(String)
func nuParseExpr(t *testing.T, js string) rhino.Node {
	t.Helper()
	script := nuParse(t, "("+js+");") // Parens force interpretation as an expression.
	return script.
		GetFirstChild(). // EXPR_RESULT
		GetFirstChild()  // expr
}

// nuGetNode performs a DFS over root searching for a descendant matching
// token. The search doesn't include root. It fails if no matching node was
// found.
//
// Java: NodeUtilTest.getNode(Node, Token)
func nuGetNode(t *testing.T, root rhino.Node, token rhino.Token) rhino.Node {
	t.Helper()
	result := nuGetNodeOrNull(root, token)
	if result.IsNil() {
		t.Fatalf("No %v node found in:\n %s", token, root.ToStringTree())
	}
	return result
}

// nuGetNodeOrNull performs a DFS below root searching for a descendant
// matching token. The search doesn't include root.
//
// Java: NodeUtilTest.getNodeOrNull(Node, Token)
func nuGetNodeOrNull(root rhino.Node, token rhino.Token) rhino.Node {
	for n := root.GetFirstChild(); !n.IsNil(); n = n.GetNext() {
		if n.GetToken() == token {
			return n
		}
		potentialMatch := nuGetNodeOrNull(n, token)
		if !potentialMatch.IsNil() {
			return potentialMatch
		}
	}
	return rhino.NilNode
}

// nuAssertThatDouble is assertThat(Double): a nullable double compares like
// Double.equals (NaN equals NaN, 0.0 differs from -0.0).
func nuAssertThatDouble(t *testing.T, actual *float64) *rhinotesting.Subject {
	if actual == nil {
		return nuAssertThat(t, nil)
	}
	return nuAssertThat(t, *actual)
}

// nuAssertThatBigInt is assertThat(BigInteger): the value compares by its
// decimal string (BigInteger.equals compares values).
func nuAssertThatBigInt(t *testing.T, actual *big.Int) *rhinotesting.Subject {
	if actual == nil {
		return nuAssertThat(t, nil)
	}
	return nuAssertThat(t, actual.String())
}

// nuNodeIterator is the Iterator<Node> of a list.
type nuNodeIterator struct {
	nodes []rhino.Node
	i     int
}

func nuIteratorOf(nodes []rhino.Node) *nuNodeIterator { return &nuNodeIterator{nodes: nodes} }

// Java: Iterator.hasNext()
func (it *nuNodeIterator) HasNext() bool { return it.i < len(it.nodes) }

// Java: Iterator.next()
func (it *nuNodeIterator) Next() rhino.Node {
	if !it.HasNext() {
		panic(jutil.NewNoSuchElementException())
	}
	it.i++
	return it.nodes[it.i-1]
}

// nuFeatureSetSubject is com.google.javascript.jscomp.parsing.parser.testing.FeatureSetSubject.
type nuFeatureSetSubject struct {
	t      *testing.T
	actual *parser.FeatureSet
}

// Java: FeatureSetSubject.assertFS(FeatureSet)
func nuAssertFS(t *testing.T, fs *parser.FeatureSet) *nuFeatureSetSubject {
	return &nuFeatureSetSubject{t: t, actual: fs}
}

// Java: FeatureSetSubject.equals(FeatureSet)
func (s *nuFeatureSetSubject) Equals(other *parser.FeatureSet) {
	s.t.Helper()
	if !s.actual.Equals(other) {
		s.t.Fatal(jutil.LenientFormat("Expected a FeatureSet equal to: %s\nBut got: %s", other, s.actual))
	}
}

var (
	_ = math.NaN
	_ = big.NewInt
	_ = strings.Contains
	_ = base.Tri_TRUE
	_ = parser.FeatureSet_ES3
)

// ==== IsDefinedValueTests public static final class IsDefinedValueTests {
func TestNodeUtil_IsDefinedValueTests_TestIsDefinedValue(t *testing.T) {
	// while "null" is "defined" for the purposes of this method, it triggers the RHS.
	nuAssertThat(t, NodeUtil_IsDefinedValue(nuParseExpr(t, "null ?? undefined"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsDefinedValue(nuParseExpr(t, "null ?? null"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsDefinedValue(nuParseExpr(t, "undefined ?? undefined"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsDefinedValue(nuParseExpr(t, "undefined ?? null"))).IsTrue()

	// could be true but the logic is not refined enough
	nuAssertThat(t, NodeUtil_IsDefinedValue(nuParseExpr(t, "0 ?? undefined"))).IsFalse()
}

func TestNodeUtil_IsDefinedValueTests_IsDefinedValueOptionalChain(t *testing.T) {
	nuAssertThat(t, NodeUtil_IsDefinedValue(nuParseExpr(t, "x?.y"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsDefinedValue(nuParseExpr(t, "x?.[y]"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsDefinedValue(nuParseExpr(t, "x?.()"))).IsFalse()
}

// ==== IsFromTypeSummaryTests public static final class IsFromTypeSummaryTests {

func TestNodeUtil_IsFromTypeSummaryTests_IdentifiesTypeSummaryJSDoc(t *testing.T) {
	script := nuParse(t, "/** @fileoverview @typeSummary */ var x;")
	nuAssertThat(t, NodeUtil_IsFromTypeSummary(script)).IsTrue()

	scriptWithExterns :=
		nuParse(t, "// IJS\n/** @fileoverview @typeSummary */\n/** @externs */\n;")
	nuAssertThat(t, NodeUtil_IsFromTypeSummary(scriptWithExterns)).IsTrue()

	regularScript := nuParse(t, "var x;")
	nuAssertThat(t, NodeUtil_IsFromTypeSummary(regularScript)).IsFalse()
}

/**
 * Test the forms of getting a boolean value for a node.
 *
 * <p>The 3 forms differ based on when they decide the value is unknown. See more description on
 * the method definitions.
 *
 * <dl>
 *   <dt>pure
 *   <dd>known only for side-effect free literal(ish) values
 *   <dt>literal
 *   <dd>known for literal(ish) values regardless of side-effects
 *   <dt>impure
 *   <dd>known for literal(ish) values and some simple expressions regardless of side-effects
 * </dl>
 */

// ==== BooleanValueTests
//
// Test the forms of getting a boolean value for a node.
//
// Java: NodeUtilTest.BooleanValueTests (@Parameterized: jsExpression, expectedResult)
func TestNodeUtil_BooleanValueTests_GetBooleanValue(t *testing.T) {
	cases := []struct {
		// Snippet of JS that will be parsed as an expression.
		jsExpression string
		// Expected result of NodeUtil.getBooleanValue()
		expectedResult base.Tri
	}{
		// truly literal, side-effect free values are always known
		{"true", base.Tri_TRUE},
		{"10", base.Tri_TRUE},
		{"1n", base.Tri_TRUE},
		{"'0'", base.Tri_TRUE},
		{"/a/", base.Tri_TRUE},
		{"{}", base.Tri_TRUE},
		{"[]", base.Tri_TRUE},
		{"false", base.Tri_FALSE},
		{"null", base.Tri_FALSE},
		{"0", base.Tri_FALSE},
		{"0n", base.Tri_FALSE},
		{"''", base.Tri_FALSE},
		{"undefined", base.Tri_FALSE},

		// literals that have side-effects aren't pure
		{"{a:foo()}", base.Tri_TRUE},
		{"[foo()]", base.Tri_TRUE},

		// not really literals, but we pretend they are for our purposes
		{"void 0", base.Tri_FALSE},
		// side-effect keeps this one from being pure
		{"void foo()", base.Tri_FALSE},
		{"!true", base.Tri_FALSE},
		{"!false", base.Tri_TRUE},
		{"!''", base.Tri_TRUE},
		{"class Klass {}", base.Tri_TRUE},
		{"new Date()", base.Tri_TRUE},
		{"b", base.Tri_UNKNOWN},

		// Numbers
		{"+'0.0'", base.Tri_FALSE},
		{"+'1.0'", base.Tri_TRUE},
		{"+'3.1'", base.Tri_TRUE},
		{"+3.1", base.Tri_TRUE},
		{"-'0.0'", base.Tri_FALSE},
		{"-'1.0'", base.Tri_TRUE},
		{"-'1.0'", base.Tri_TRUE},
		{"-'3.1'", base.Tri_TRUE},
		{"-0.0", base.Tri_FALSE},
		{"-1.0", base.Tri_TRUE},
		{"-3.1", base.Tri_TRUE},
		{"-Infinity", base.Tri_TRUE},
		{"-NaN", base.Tri_FALSE},
		{"0.0", base.Tri_FALSE},
		{"1", base.Tri_TRUE},
		{"Infinity", base.Tri_TRUE},
		{"NaN", base.Tri_FALSE},
		{"~'-1.0'", base.Tri_FALSE},
		{"~-1.0", base.Tri_FALSE},
		{"~0.0", base.Tri_TRUE},
		{"~3.1", base.Tri_TRUE},
		{"~NaN", base.Tri_TRUE},

		// BigInts
		{"0n", base.Tri_FALSE},
		{"-0n", base.Tri_FALSE},
		{"~0n", base.Tri_TRUE},
		{"~-0n", base.Tri_TRUE},
		{"-~0n", base.Tri_TRUE},
		{"1n", base.Tri_TRUE},
		{"-1n", base.Tri_TRUE},
		{"~1n", base.Tri_TRUE},
		{"~-1n", base.Tri_FALSE},
		{"-~1n", base.Tri_TRUE},

		// template literals
		{"``", base.Tri_FALSE},
		{"`definiteLength`", base.Tri_TRUE},
		{"`${some}str`", base.Tri_UNKNOWN},

		// non-literal expressions
		{"a=true", base.Tri_TRUE},
		{"a=false", base.Tri_FALSE},
		{"a=(false,true)", base.Tri_TRUE},
		{"a=(true,false)", base.Tri_FALSE},
		{"a=(false || true)", base.Tri_TRUE},
		{"a=(true && false)", base.Tri_FALSE},
		{"a=!(true && false)", base.Tri_TRUE},
		{"a,true", base.Tri_TRUE},
		{"a,false", base.Tri_FALSE},
		{"true||false", base.Tri_TRUE},
		{"false||false", base.Tri_FALSE},
		{"true&&true", base.Tri_TRUE},
		{"true&&false", base.Tri_FALSE},

		// Logical assignments
		{"a||(a=true)", base.Tri_TRUE},
		{"a||(a=false)", base.Tri_UNKNOWN},
		{"a||=true", base.Tri_TRUE},
		{"a||=false", base.Tri_UNKNOWN},
		{"a&&(a=true)", base.Tri_UNKNOWN},
		{"a&&(a=false)", base.Tri_FALSE},
		{"a&&=true", base.Tri_UNKNOWN},
		{"a&&=false", base.Tri_FALSE},
		{"a??(a=true)", base.Tri_UNKNOWN},
		{"a??=true", base.Tri_UNKNOWN},
		{"a??(a=false)", base.Tri_UNKNOWN},
		{"a??=false", base.Tri_UNKNOWN},

		// Assignment ops other than ASSIGN are unknown.
		{"a *= 2", base.Tri_UNKNOWN},

		// Complex expressions that contain anything other then "=", ",", or "!" are
		// unknown.
		{"2 + 2", base.Tri_UNKNOWN},

		// assignment values are the RHS
		{"a=1", base.Tri_TRUE},
		{"a=/a/", base.Tri_TRUE},
		{"a={}", base.Tri_TRUE},

		// hooks have impure boolean value if both cases have same impure boolean value
		{"a?true:true", base.Tri_TRUE},
		{"a?false:false", base.Tri_FALSE},
		{"a?true:false", base.Tri_UNKNOWN},
		{"a?true:foo()", base.Tri_UNKNOWN},

		// coalesce returns LHS if LHS is truthy or if LHS and RHS have same boolean value
		{"null??false", base.Tri_FALSE}, // both false
		{"2??[]", base.Tri_TRUE},        // both true
		{"{}??false", base.Tri_TRUE},    // LHS is true
		{"undefined??[]", base.Tri_UNKNOWN},
		{"foo()??true", base.Tri_UNKNOWN},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("getBooleanValue(\"%s\") => %v", c.jsExpression, c.expectedResult), func(t *testing.T) {
			nuAssertThat(t, NodeUtil_GetBooleanValue(nuParseExpr(t, c.jsExpression))).IsEqualTo(c.expectedResult)
		})
	}
}

// ==== IsPropertyTestTests public static final class IsPropertyTestTests {

func TestNodeUtil_IsPropertyTestTests_OptChainGetPropIsPropertyTest(t *testing.T) {
	compiler := NewCompiler() // Java: new Compiler()
	getProp := nuParseExpr(t, "x.y?.z")
	nuAssertThat(t, NodeUtil_IsPropertyTest(compiler, getProp.GetFirstChild())).IsTrue()
}

func TestNodeUtil_IsPropertyTestTests_OptChainGetElemIsPropertyTest(t *testing.T) {
	compiler := NewCompiler() // Java: new Compiler()
	getElem := nuParseExpr(t, "x.y?.[z]")
	nuAssertThat(t, NodeUtil_IsPropertyTest(compiler, getElem.GetFirstChild())).IsTrue()
}

func TestNodeUtil_IsPropertyTestTests_OptChainCallIsPropertyTest(t *testing.T) {
	compiler := NewCompiler() // Java: new Compiler()
	call := nuParseExpr(t, "x.y?.(z)")
	nuAssertThat(t, NodeUtil_IsPropertyTest(compiler, call.GetFirstChild())).IsTrue()
}

func TestNodeUtil_IsPropertyTestTests_OptChainNonStartOfChainIsPropertyTest(t *testing.T) {
	compiler := NewCompiler() // Java: new Compiler()
	getProp := nuParseExpr(t, "x.y?.z.foo.bar")
	nuAssertThat(t, NodeUtil_IsPropertyTest(compiler, getProp.GetFirstChild())).IsTrue()
}

/**
 * A nested class to allow the `Enclosed` runner to run these tests since it doesn't run tests in
 * the outer class.
 *
 * <p>There's no particular organization to these tests since they were written in JUnit3. Many of
 * them could benefit from being broken out into their own classes.
 */

// ==== AssortedTests public static final class AssortedTests {

func TestNodeUtil_AssortedTests_TestIsLiteralOrConstValue(t *testing.T) {
	nuh_assertLiteralAndImmutable(t, nuParseExpr(t, "10"))
	nuh_assertLiteralAndImmutable(t, nuParseExpr(t, "-10"))
	nuh_assertLiteralButNotImmutable(t, nuParseExpr(t, "[10, 20]"))
	nuh_assertLiteralButNotImmutable(t, nuParseExpr(t, "{'a': 20}"))
	nuh_assertLiteralButNotImmutable(t, nuParseExpr(t, "[10, , 1.0, [undefined], 'a']"))
	nuh_assertLiteralButNotImmutable(t, nuParseExpr(t, "/abc/"))
	nuh_assertLiteralAndImmutable(t, nuParseExpr(t, "\"string\""))
	nuh_assertLiteralAndImmutable(t, nuParseExpr(t, "'aaa'"))
	nuh_assertLiteralAndImmutable(t, nuParseExpr(t, "null"))
	nuh_assertLiteralAndImmutable(t, nuParseExpr(t, "undefined"))
	nuh_assertLiteralAndImmutable(t, nuParseExpr(t, "void 0"))
	nuh_assertNotLiteral(t, nuParseExpr(t, "abc"))
	nuh_assertNotLiteral(t, nuParseExpr(t, "[10, foo(), 20]"))
	nuh_assertNotLiteral(t, nuParseExpr(t, "foo()"))
	nuh_assertNotLiteral(t, nuParseExpr(t, "c + d"))
	nuh_assertNotLiteral(t, nuParseExpr(t, "{'a': foo()}"))
	nuh_assertNotLiteral(t, nuParseExpr(t, "void foo()"))
}

func TestNodeUtil_AssortedTests_TestObjectLiteralIsLiteralValue(t *testing.T) {
	nuAssertThat(t, nuh_isLiteralValue(t, "{a: 20}")).IsTrue()
	nuAssertThat(t, nuh_isLiteralValue(t, "{'a': 20}")).IsTrue()
	nuAssertThat(t, nuh_isLiteralValue(t, "{a: function() {}}")).IsTrue()
	nuAssertThat(t, nuh_isLiteralValueExcludingFunctions(t, "{a: function() {}}")).IsFalse()
	nuAssertThat(t, nuh_isLiteralValue(t, "{a() {}}")).IsTrue()
	nuAssertThat(t, nuh_isLiteralValueExcludingFunctions(t, "{a() {}}")).IsFalse()
	nuAssertThat(t, nuh_isLiteralValue(t, "{'a'() {}}")).IsTrue()
	nuAssertThat(t, nuh_isLiteralValueExcludingFunctions(t, "{'a'() {}}")).IsFalse()

	nuAssertThat(t, nuh_isLiteralValue(t, "{['a']: 20}")).IsTrue()
	nuAssertThat(t, nuh_isLiteralValue(t, "{[b]: 20}")).IsFalse()
	nuAssertThat(t, nuh_isLiteralValue(t, "{['a']() {}}")).IsTrue()
	nuAssertThat(t, nuh_isLiteralValue(t, "{[b]() {}}")).IsFalse()
	nuAssertThat(t, nuh_isLiteralValueExcludingFunctions(t, "{['a']() {}}")).IsFalse()

	nuAssertThat(t, nuh_isLiteralValue(t, "{ get a() { return 0; } }")).IsTrue()
	nuAssertThat(t, nuh_isLiteralValue(t, "{ get 'a'() { return 0; } }")).IsTrue()
	nuAssertThat(t, nuh_isLiteralValue(t, "{ get ['a']() { return 0; } }")).IsTrue()
	nuAssertThat(t, nuh_isLiteralValue(t, "{ get 123() { return 0; } }")).IsTrue()
	nuAssertThat(t, nuh_isLiteralValue(t, "{ get [123]() { return 0; } }")).IsTrue()
	nuAssertThat(t, nuh_isLiteralValue(t, "{ get [b]() { return 0; } }")).IsFalse()
	nuAssertThat(t, nuh_isLiteralValue(t, "{ set a(x) { } }")).IsTrue()
	nuAssertThat(t, nuh_isLiteralValue(t, "{ set 'a'(x) { } }")).IsTrue()
	nuAssertThat(t, nuh_isLiteralValue(t, "{ set ['a'](x) { } }")).IsTrue()
	nuAssertThat(t, nuh_isLiteralValue(t, "{ set 123(x) { } }")).IsTrue()
	nuAssertThat(t, nuh_isLiteralValue(t, "{ set [123](x) { } }")).IsTrue()
	nuAssertThat(t, nuh_isLiteralValue(t, "{ set [b](x) { } }")).IsFalse()
}

func nuh_isLiteralValueExcludingFunctions(t *testing.T, code string) bool {
	return NodeUtil_IsLiteralValue(nuParseExpr(t, code) /* includeFunctions= */, false)
}

func nuh_isLiteralValue(t *testing.T, code string) bool {
	return NodeUtil_IsLiteralValue(nuParseExpr(t, code) /* includeFunctions= */, true)
}

func nuh_assertLiteralAndImmutable(t *testing.T, n rhino.Node) {
	nuAssertThat(t, NodeUtil_IsLiteralValue(n, true)).IsTrue()
	nuAssertThat(t, NodeUtil_IsLiteralValue(n, false)).IsTrue()
	nuAssertThat(t, NodeUtil_IsImmutableValue(n)).IsTrue()
}

func nuh_assertLiteralButNotImmutable(t *testing.T, n rhino.Node) {
	nuAssertThat(t, NodeUtil_IsLiteralValue(n, true)).IsTrue()
	nuAssertThat(t, NodeUtil_IsLiteralValue(n, false)).IsTrue()
	nuAssertThat(t, NodeUtil_IsImmutableValue(n)).IsFalse()
}

func nuh_assertNotLiteral(t *testing.T, n rhino.Node) {
	nuAssertThat(t, NodeUtil_IsLiteralValue(n, true)).IsFalse()
	nuAssertThat(t, NodeUtil_IsLiteralValue(n, false)).IsFalse()
	nuAssertThat(t, NodeUtil_IsImmutableValue(n)).IsFalse()
}

func TestNodeUtil_AssortedTests_TestIsStringLiteralStringLiteral(t *testing.T) {
	nuAssertThat(t, NodeUtil_IsSomeCompileTimeConstStringValue(nuParseExpr(t, "'b'"))).IsTrue() // String
}

func TestNodeUtil_AssortedTests_TestIsStringLiteralTemplateNoSubst(t *testing.T) {
	nuAssertThat(t, NodeUtil_IsSomeCompileTimeConstStringValue(nuParseExpr(t, "`b`"))).IsTrue() // Template
}

func TestNodeUtil_AssortedTests_TestIsStringLiteralTemplateWithSubstitution(t *testing.T) {
	// Template with substitution
	nuAssertThat(t, NodeUtil_IsSomeCompileTimeConstStringValue(nuParseExpr(t, "`foo${bar}`"))).IsFalse()
}

func TestNodeUtil_AssortedTests_TestIsStringLiteralVariable(t *testing.T) {
	nuAssertThat(t, NodeUtil_IsSomeCompileTimeConstStringValue(nuParseExpr(t, "b"))).IsFalse() // Variable
}

func TestNodeUtil_AssortedTests_TestIsStringLiteralConcatLiterals(t *testing.T) {
	nuAssertThat(t, NodeUtil_IsSomeCompileTimeConstStringValue(nuParseExpr(t, "'b' + 'c'"))).IsTrue() // String concatenation
}

func TestNodeUtil_AssortedTests_TestIsStringLiteralConcatLiteralVariable(t *testing.T) {
	nuAssertThat(t, NodeUtil_IsSomeCompileTimeConstStringValue(nuParseExpr(t, "b + 'c'"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsSomeCompileTimeConstStringValue(nuParseExpr(t, "'b' + c"))).IsFalse()
}

func TestNodeUtil_AssortedTests_TestIsStringLiteralTernary(t *testing.T) {
	nuAssertThat(t, NodeUtil_IsSomeCompileTimeConstStringValue(nuParseExpr(t, "a ? 'b' : 'c'"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsSomeCompileTimeConstStringValue(nuParseExpr(t, "a ? b : 'c'"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsSomeCompileTimeConstStringValue(nuParseExpr(t, "a ? 'b' : c"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsSomeCompileTimeConstStringValue(nuParseExpr(t, "a ? b : c"))).IsFalse()
}

func TestNodeUtil_AssortedTests_TestGetStringValue(t *testing.T) {
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "true"))).IsEqualTo("true")
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "10"))).IsEqualTo("10")
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "1.0"))).IsEqualTo("1")
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "-1.0"))).IsEqualTo("-1")

	/* See https://github.com/google/closure-compiler/issues/1262 */
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "1.2323919403474454e+21"))).IsEqualTo("1.2323919403474454e+21")

	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "'0'"))).IsEqualTo("0")
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "/a/"))).IsNull()
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "{}"))).IsEqualTo("[object Object]")
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "{get foo() { return 1; }}"))).IsNull()
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "{set foo(x) {}}"))).IsNull()
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "{toString: 1}"))).IsNull()
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "{'toString': 1}"))).IsNull()
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "{toString() { return 'foo'; }}"))).IsNull()
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "{get toString() { return 'foo'; }}"))).IsNull()
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "{set toString(x) {}}"))).IsNull()
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "{valueOf: 1}"))).IsNull()
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "{'valueOf': 1}"))).IsNull()
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "{valueOf() { return 1; }}"))).IsNull()
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "{get valueOf() { return 1; }}"))).IsNull()
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "{set valueOf(x) {}}"))).IsNull()
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "{__proto__: null}"))).IsNull()
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "{'__proto__': null}"))).IsNull()
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "{__proto__: p}"))).IsNull()
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "{[key]: 'val'}"))).IsNull()
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "{...spread}"))).IsNull()
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "[]"))).IsEmpty()
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "false"))).IsEqualTo("false")
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "null"))).IsEqualTo("null")
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "0"))).IsEqualTo("0")
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "1n"))).IsEqualTo("1")
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "''"))).IsEmpty()
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "undefined"))).IsEqualTo("undefined")
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "void 0"))).IsEqualTo("undefined")
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "void foo()"))).IsEqualTo("undefined")

	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "NaN"))).IsEqualTo("NaN")
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "Infinity"))).IsEqualTo("Infinity")
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "x"))).IsNull()

	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "`Hello`"))).IsEqualTo("Hello")
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "`Hello ${'foo'}`"))).IsEqualTo("Hello foo")
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "`Hello ${name}`"))).IsNull()
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "`${4} bananas`"))).IsEqualTo("4 bananas")
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "`This is ${true}.`"))).IsEqualTo("This is true.")
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "`${'hello'} ${name}`"))).IsNull()
}

func TestNodeUtil_AssortedTests_TestGetArrayStringValue(t *testing.T) {
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "[]"))).IsEmpty()
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "['']"))).IsEmpty()
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "[null]"))).IsEmpty()
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "[undefined]"))).IsEmpty()
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "[void 0]"))).IsEmpty()
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "[NaN]"))).IsEqualTo("NaN")
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "[,'']"))).IsEqualTo(",")
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "[[''],[''],['']]"))).IsEqualTo(",,")
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "[[1.0],[2.0]]"))).IsEqualTo("1,2")
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "[a]"))).IsNull()
	nuAssertThat(t, NodeUtil_GetStringValue(nuParseExpr(t, "[1,a]"))).IsNull()
}

func TestNodeUtil_AssortedTests_TestMayBeObjectLitKey(t *testing.T) {
	nuh_assertMayBeObjectLitKey(t, nuParseExpr(t, "({})"), false)
	nuh_assertMayBeObjectLitKey(t, nuParseExpr(t, "a"), false)
	nuh_assertMayBeObjectLitKey(t, nuParseExpr(t, "'a'"), false)
	nuh_assertMayBeObjectLitKey(t, nuParseExpr(t, "1"), false)
	nuh_assertMayBeObjectLitKey(t, nuParseExpr(t, "1n"), false)
	nuh_assertMayBeObjectLitKey(t, nuParseExpr(t, "({a: 1})").GetFirstChild(), true)
	nuh_assertMayBeObjectLitKey(t, nuParseExpr(t, "({1: 1})").GetFirstChild(), true)
	nuh_assertMayBeObjectLitKey(t, nuParseExpr(t, "({get a(){}})").GetFirstChild(), true)
	nuh_assertMayBeObjectLitKey(t, nuParseExpr(t, "({set a(b){}})").GetFirstChild(), true)
	// returns false for computed properties (b/111621528)
	nuh_assertMayBeObjectLitKey(t, nuParseExpr(t, "({['a']: 1})").GetFirstChild(), false)
	// returns true for non-object-literal keys
	nuh_assertMayBeObjectLitKey(t, nuParseExpr(t, "({a} = {})").GetFirstFirstChild(), true)
	nuh_assertMayBeObjectLitKey(t, nuParseExpr(t, "(class { a() {} })").GetLastChild().GetFirstChild(), true)
}

func TestNodeUtil_AssortedTests_TestIsObjectLitKey(t *testing.T) {
	nuh_assertIsObjectLitKey(t, nuParseExpr(t, "({a: 1})").GetFirstChild(), true)
	// returns false for computed properties (b/111621528)
	nuh_assertMayBeObjectLitKey(t, nuParseExpr(t, "({['a']: 1})").GetFirstChild(), false)
	// returns false for object patterns
	nuh_assertIsObjectLitKey(t, nuParseExpr(t, "({a} = {})").GetFirstFirstChild(), false)
	nuh_assertIsObjectLitKey(t, nuParseExpr(t, "(class { a() {} })").GetLastChild().GetFirstChild(), false)
}

func nuh_assertMayBeObjectLitKey(t *testing.T, node rhino.Node, expected bool) {
	nuAssertThat(t, NodeUtil_MayBeObjectLitKey(node)).IsEqualTo(expected)
}

func nuh_assertIsObjectLitKey(t *testing.T, node rhino.Node, expected bool) {
	nuAssertThat(t, NodeUtil_IsObjectLitKey(node)).IsEqualTo(expected)
}

func TestNodeUtil_AssortedTests_TestGetFunctionName1(t *testing.T) {
	parent := nuParse(t, "function name(){}")
	nuh_assertGetNameResult(t, parent.GetFirstChild(), "name")
}

func TestNodeUtil_AssortedTests_TestGetFunctionName2(t *testing.T) {
	parent := nuParse(t, "var name = function(){}").GetFirstFirstChild()

	nuh_assertGetNameResult(t, parent.GetFirstChild(), "name")
}

func TestNodeUtil_AssortedTests_TestGetFunctionName3(t *testing.T) {
	parent := nuParse(t, "qualified.name = function(){}").GetFirstFirstChild()

	nuh_assertGetNameResult(t, parent.GetLastChild(), "qualified.name")
}

func TestNodeUtil_AssortedTests_TestGetFunctionName4(t *testing.T) {
	parent := nuParse(t, "var name2 = function name1(){}").GetFirstFirstChild()

	nuh_assertGetNameResult(t, parent.GetFirstChild(), "name2")
}

func TestNodeUtil_AssortedTests_TestGetFunctionName5(t *testing.T) {
	n := nuParse(t, "qualified.name2 = function name1(){}")
	parent := n.GetFirstFirstChild()

	nuh_assertGetNameResult(t, parent.GetLastChild(), "qualified.name2")
}

func TestNodeUtil_AssortedTests_TestGetBestFunctionName1(t *testing.T) {
	parent := nuParse(t, "function func(){}")

	nuAssertThat(t, NodeUtil_GetNearestFunctionName(parent.GetFirstChild())).IsEqualTo("func")
}

func TestNodeUtil_AssortedTests_TestGetBestFunctionName2(t *testing.T) {
	parent := nuParse(t, "var obj = {memFunc(){}}").GetFirstFirstChild().GetFirstFirstChild()

	nuAssertThat(t, NodeUtil_GetNearestFunctionName(parent.GetLastChild())).IsEqualTo("memFunc")
}

func TestNodeUtil_AssortedTests_TestConstKeywordNamespace(t *testing.T) {
	decl := nuParse(t, "const ns = {};").GetFirstChild()
	jutil.CheckStateMsg1(decl.IsConst(), "%s", decl)

	nameNode := decl.GetFirstChild()
	jutil.CheckStateMsg1(nameNode.IsName(), "%s", nameNode)

	nuAssertThat(t, NodeUtil_IsNamespaceDecl(nameNode)).IsTrue()
}

func TestNodeUtil_AssortedTests_TestTypedefNamespace(t *testing.T) {
	statements :=
		nuParse(t,
			"/** @typedef */ const obj = {};\n/** @const */ obj.a = {};\n/** @typedef @const */ obj.b = {};\n/** @typedef */ obj.c = {};\n/** @typedef @const */ obj.d = function() {};\n")
	namespaceNames := jutil.NewLinkedHashSet[string]()
	nonNamespaceNames := jutil.NewLinkedHashSet[string]()
	for statement := range statements.Children() {
		var name rhino.Node
		if statement.IsExprResult() {
			name = statement.GetFirstFirstChild()
		} else {
			name = statement.GetFirstChild()
		}
		jutil.CheckStateMsg1(name.IsName() || name.IsGetProp(), "%s", name)
		if NodeUtil_IsNamespaceDecl(name) {
			namespaceNames.Add(*name.GetQualifiedName())
		} else {
			nonNamespaceNames.Add(*name.GetQualifiedName())
		}
	}
	nuAssertThat(t, namespaceNames).ContainsExactly("obj", "obj.a", "obj.b")
	nuAssertThat(t, nonNamespaceNames).ContainsExactly("obj.c", "obj.d")
}

func nuh_assertGetNameResult(t *testing.T, function rhino.Node, name string) {
	nuAssertNode(t, function).HasToken(rhino.Token_FUNCTION)
	nuAssertThat(t, NodeUtil_GetName(function)).IsEqualTo(name)
}

func TestNodeUtil_AssortedTests_TestIsFunctionDeclaration(t *testing.T) {
	nuAssertThat(t, NodeUtil_IsFunctionDeclaration(nuParseFirst(t, rhino.Token_FUNCTION, "function foo(){}"))).IsTrue()
	nuAssertThat(t,
		NodeUtil_IsFunctionDeclaration(nuParseFirst(t, rhino.Token_FUNCTION, "class C { constructor() {} }"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsFunctionDeclaration(nuParseFirst(t, rhino.Token_FUNCTION, "({ foo() {} })"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsFunctionDeclaration(nuParseFirst(t, rhino.Token_FUNCTION, "var x = function(){}"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsFunctionDeclaration(nuParseFirst(t, rhino.Token_FUNCTION, "export function f() {}"))).IsTrue()
	nuAssertThat(t,
		NodeUtil_IsFunctionDeclaration(nuParseFirst(t, rhino.Token_FUNCTION, "export default function() {}"))).IsFalse()
	nuAssertThat(t,
		NodeUtil_IsFunctionDeclaration(
			nuParseFirst(t, rhino.Token_FUNCTION, "export default function foo() {}"))).IsTrue()
	nuAssertThat(t,
		NodeUtil_IsFunctionDeclaration(
			nuParseFirst(t, rhino.Token_FUNCTION, "export default (foo) => { alert(foo); }"))).IsFalse()
}

// Java: NodeUtilTest.AssortedTests.checkParamSideEffects(String)
func nuh_checkParamSideEffects(t *testing.T, js string) bool {
	t.Helper()
	fn, compiler := nuParseFirstWithCompiler(t, rhino.Token_FUNCTION, js)
	return NodeUtil_FunctionParametersMayHaveSideEffects(fn, compiler.GetAstAnalyzer())
}

func TestNodeUtil_AssortedTests_TestFunctionParametersMayHaveSideEffects(t *testing.T) {
	nuAssertThat(t, nuh_checkParamSideEffects(t, "function f() {}")).IsFalse()
	nuAssertThat(t, nuh_checkParamSideEffects(t, "function f(a, b) {}")).IsFalse()
	nuAssertThat(t, nuh_checkParamSideEffects(t, "function f(a = 1) {}")).IsFalse()
	nuAssertThat(t, nuh_checkParamSideEffects(t, "function f(a = 'default') {}")).IsFalse()
	nuAssertThat(t, nuh_checkParamSideEffects(t, "function f(...a) {}")).IsFalse()

	nuAssertThat(t, nuh_checkParamSideEffects(t, "function f(a = audit()) {}")).IsTrue()
	nuAssertThat(t, nuh_checkParamSideEffects(t, "function f({a}) {}")).IsTrue()
	nuAssertThat(t, nuh_checkParamSideEffects(t, "function f([a]) {}")).IsTrue()
	nuAssertThat(t, nuh_checkParamSideEffects(t, "function f({a} = {}) {}")).IsTrue()
	nuAssertThat(t, nuh_checkParamSideEffects(t, "function f(...[a]) {}")).IsTrue()
	nuAssertThat(t, nuh_checkParamSideEffects(t, "function f(a, {b}, ...c) {}")).IsTrue()
}

func TestNodeUtil_AssortedTests_TestisBlockScopedFunctionDeclaration(t *testing.T) {
	nuAssertThat(t,
		NodeUtil_IsBlockScopedFunctionDeclaration(nuParseFirst(t, rhino.Token_FUNCTION, "function foo(){}"))).IsFalse()
	nuAssertThat(t,
		NodeUtil_IsBlockScopedFunctionDeclaration(
			nuParseFirst(t, rhino.Token_FUNCTION, "var x = function(){}"))).IsFalse()
	nuAssertThat(t,
		NodeUtil_IsBlockScopedFunctionDeclaration(
			nuParseFirst(t, rhino.Token_FUNCTION, "{ function foo(){} }"))).IsTrue()
	nuAssertThat(t,
		NodeUtil_IsBlockScopedFunctionDeclaration(
			nuParseFirst(t, rhino.Token_FUNCTION, "class C { static { function foo(){} } }"))).IsTrue()
	nuAssertThat(t,
		NodeUtil_IsBlockScopedFunctionDeclaration(
			nuParseSecond(t, rhino.Token_FUNCTION, "function foo(){ class C { static { function a(){} } } }"))).IsTrue()
}

func TestNodeUtil_AssortedTests_TestIsMethodDeclaration(t *testing.T) {
	nuAssertThat(t, NodeUtil_IsMethodDeclaration(nuParseFirst(t, rhino.Token_FUNCTION, "class C { constructor() {} }"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsMethodDeclaration(nuParseFirst(t, rhino.Token_FUNCTION, "class C { a() {} }"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsMethodDeclaration(nuParseFirst(t, rhino.Token_FUNCTION, "class C { static a() {} }"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsMethodDeclaration(nuParseFirst(t, rhino.Token_FUNCTION, "({ set foo(v) {} })"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsMethodDeclaration(nuParseFirst(t, rhino.Token_FUNCTION, "({ get foo() {} })"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsMethodDeclaration(nuParseFirst(t, rhino.Token_FUNCTION, "({ [foo]() {} })"))).IsTrue()
}

func TestNodeUtil_AssortedTests_TestIsClassDeclaration(t *testing.T) {
	nuAssertThat(t, NodeUtil_IsClassDeclaration(nuParseFirst(t, rhino.Token_CLASS, "class Foo {}"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsClassDeclaration(nuParseFirst(t, rhino.Token_CLASS, "var Foo = class {}"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsClassDeclaration(nuParseFirst(t, rhino.Token_CLASS, "var Foo = class Foo{}"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsClassDeclaration(nuParseFirst(t, rhino.Token_CLASS, "export default class Foo {}"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsClassDeclaration(nuParseFirst(t, rhino.Token_CLASS, "export class Foo {}"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsClassDeclaration(nuParseFirst(t, rhino.Token_CLASS, "export default class {}"))).IsFalse()
}

func TestNodeUtil_AssortedTests_TestIsFunctionExpression(t *testing.T) {
	nuh_assertContainsAnonFunc(t, true, "(function(){})")
	nuh_assertContainsAnonFunc(t, true, "[function a(){}]")
	nuh_assertContainsAnonFunc(t, true, "({x: function a(){}})")
	nuh_assertContainsAnonFunc(t, true, "({[0]: function a() {}})")
	nuh_assertContainsAnonFunc(t, true, "(function a(){})()")
	nuh_assertContainsAnonFunc(t, true, "x = function a(){};")
	nuh_assertContainsAnonFunc(t, true, "var x = function a(){};")
	nuh_assertContainsAnonFunc(t, true, "if (function a(){});")
	nuh_assertContainsAnonFunc(t, true, "while (function a(){});")
	nuh_assertContainsAnonFunc(t, true, "do; while (function a(){});")
	nuh_assertContainsAnonFunc(t, true, "for (function a(){};;);")
	nuh_assertContainsAnonFunc(t, true, "for (;function a(){};);")
	nuh_assertContainsAnonFunc(t, true, "for (;;function a(){});")
	nuh_assertContainsAnonFunc(t, true, "for (p in function a(){});")
	nuh_assertContainsAnonFunc(t, true, "with (function a(){}) {}")
	nuh_assertContainsAnonFunc(t, false, "function a(){}")
	nuh_assertContainsAnonFunc(t, false, "if (x) function a(){};")
	nuh_assertContainsAnonFunc(t, false, "if (x) { function a(){} }")
	nuh_assertContainsAnonFunc(t, false, "if (x); else function a(){};")
	nuh_assertContainsAnonFunc(t, false, "while (x) function a(){};")
	nuh_assertContainsAnonFunc(t, false, "do function a(){} while (0);")
	nuh_assertContainsAnonFunc(t, false, "for (;;) function a(){}")
	nuh_assertContainsAnonFunc(t, false, "for (p in o) function a(){};")
	nuh_assertContainsAnonFunc(t, false, "with (x) function a(){}")
	nuh_assertContainsAnonFunc(t, true, "export default function() {};")
	nuh_assertContainsAnonFunc(t, false, "export default function a() {};")
	nuh_assertContainsAnonFunc(t, false, "export function a() {};")
	nuh_assertContainsAnonFunc(t, false, "class C { a() {} }")
	nuh_assertContainsAnonFunc(t, false, "class C { static a() {} }")
	nuh_assertContainsAnonFunc(t, false, "x = { a() {} }")
	nuh_assertContainsAnonFunc(t, false, "({[0]() {}})")
	nuh_assertContainsAnonFunc(t, false, "({get [0]() {}})")
	nuh_assertContainsAnonFunc(t, false, "({set [0](x) {}})")
}

func nuh_assertContainsAnonFunc(t *testing.T, expected bool, js string) {
	funcParent := nuh_findParentOfFuncOrClassDescendant(t, nuParse(t, js), rhino.Token_FUNCTION)
	nuAssertWithMessage(t, "Expected function node in parse tree of: %s", js).That(funcParent).IsNotNull()

	funcNode := nuh_getFuncOrClassChild(t, funcParent, rhino.Token_FUNCTION)
	nuAssertThat(t, NodeUtil_IsFunctionExpression(funcNode)).IsEqualTo(expected)
}

func TestNodeUtil_AssortedTests_TestIsClassExpression(t *testing.T) {
	nuh_assertContainsAnonClass(t, true, "(class {})")
	nuh_assertContainsAnonClass(t, true, "[class Clazz {}]")
	nuh_assertContainsAnonClass(t, true, "({x: class Clazz {}})")
	nuh_assertContainsAnonClass(t, true, "x = class Clazz {};")
	nuh_assertContainsAnonClass(t, true, "var x = class Clazz {};")
	nuh_assertContainsAnonClass(t, true, "if (class Clazz {});")
	nuh_assertContainsAnonClass(t, true, "while (class Clazz {});")
	nuh_assertContainsAnonClass(t, true, "do; while (class Clazz {});")
	nuh_assertContainsAnonClass(t, true, "for (class Clazz {};;);")
	nuh_assertContainsAnonClass(t, true, "for (;class Clazz {};);")
	nuh_assertContainsAnonClass(t, true, "for (;;class Clazz {});")
	nuh_assertContainsAnonClass(t, true, "for (p in class Clazz {});")
	nuh_assertContainsAnonClass(t, true, "with (class Clazz {}) {}")
	nuh_assertContainsAnonClass(t, false, "class Clazz {}")
	nuh_assertContainsAnonClass(t, false, "if (x) class Clazz {};")
	nuh_assertContainsAnonClass(t, false, "if (x) { class Clazz {} }")
	nuh_assertContainsAnonClass(t, false, "if (x); else class Clazz {};")
	nuh_assertContainsAnonClass(t, false, "while (x) class Clazz {};")
	nuh_assertContainsAnonClass(t, false, "do class Clazz {} while (0);")
	nuh_assertContainsAnonClass(t, false, "for (;;) class Clazz {}")
	nuh_assertContainsAnonClass(t, false, "for (p in o) class Clazz {};")
	nuh_assertContainsAnonClass(t, false, "with (x) class Clazz {}")
	nuh_assertContainsAnonClass(t, true, "export default class {};")
	nuh_assertContainsAnonClass(t, false, "export default class Clazz {};")
	nuh_assertContainsAnonClass(t, false, "export class Clazz {};")
}

func nuh_assertContainsAnonClass(t *testing.T, expected bool, js string) {
	classParent := nuh_findParentOfFuncOrClassDescendant(t, nuParse(t, js), rhino.Token_CLASS)
	nuAssertWithMessage(t, "Expected class node in parse tree of: %s", js).That(classParent).IsNotNull()
	classNode := nuh_getFuncOrClassChild(t, classParent, rhino.Token_CLASS)
	nuAssertThat(t, NodeUtil_IsClassExpression(classNode)).IsEqualTo(expected)
}

func nuh_findParentOfFuncOrClassDescendant(t *testing.T, n rhino.Node, token rhino.Token) rhino.Node {
	jutil.CheckArgument(token == rhino.Token_CLASS || token == rhino.Token_FUNCTION)
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if c.GetToken() == token {
			return n
		}
		result := nuh_findParentOfFuncOrClassDescendant(t, c, token)
		if !result.IsNil() {
			return result
		}
	}
	return rhino.NilNode
}

func nuh_getFuncOrClassChild(t *testing.T, n rhino.Node, token rhino.Token) rhino.Node {
	jutil.CheckArgument(token == rhino.Token_CLASS || token == rhino.Token_FUNCTION)
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if c.GetToken() == token {
			return c
		}
	}
	return rhino.NilNode
}

func TestNodeUtil_AssortedTests_TestGetNodeTypeReferenceCount(t *testing.T) {
	nuAssertThat(t,
		NodeUtil_GetNodeTypeReferenceCount(
			nuParse(t, "function foo(){}"), rhino.Token_THIS, func(rhino.Node) bool { return true })).IsEqualTo(0)
	nuAssertThat(t,
		NodeUtil_GetNodeTypeReferenceCount(
			nuParse(t, "this"), rhino.Token_THIS, func(rhino.Node) bool { return true })).IsEqualTo(1)
	nuAssertThat(t,
		NodeUtil_GetNodeTypeReferenceCount(
			nuParse(t, "this;function foo(){}(this)"), rhino.Token_THIS, func(rhino.Node) bool { return true })).IsEqualTo(2)
}

func TestNodeUtil_AssortedTests_TestIsNameReferenceCount(t *testing.T) {
	nuAssertThat(t, NodeUtil_IsNameReferenced(nuParse(t, "function foo(){}"), "foo")).IsTrue()
	nuAssertThat(t, NodeUtil_IsNameReferenced(nuParse(t, "var foo = function(){}"), "foo")).IsTrue()
	nuAssertThat(t, NodeUtil_IsNameReferenced(nuParse(t, "function foo(){}"), "undefined")).IsFalse()
	nuAssertThat(t, NodeUtil_IsNameReferenced(nuParse(t, "undefined"), "undefined")).IsTrue()
	nuAssertThat(t,
		NodeUtil_IsNameReferenced(
			nuParse(t, "undefined;function foo(){}(undefined)"), "undefined")).IsTrue()

	nuAssertThat(t, NodeUtil_IsNameReferenced(nuParse(t, "goo.foo"), "goo")).IsTrue()
	nuAssertThat(t, NodeUtil_IsNameReferenced(nuParse(t, "goo.foo"), "foo")).IsFalse()
}

func TestNodeUtil_AssortedTests_TestGetNameReferenceCount(t *testing.T) {
	nuAssertThat(t, NodeUtil_GetNameReferenceCount(nuParse(t, "function foo(){}"), "undefined")).IsEqualTo(0)
	nuAssertThat(t, NodeUtil_GetNameReferenceCount(nuParse(t, "undefined"), "undefined")).IsEqualTo(1)
	nuAssertThat(t,
		NodeUtil_GetNameReferenceCount(
			nuParse(t, "undefined;function foo(){}(undefined)"), "undefined")).IsEqualTo(2)

	nuAssertThat(t, NodeUtil_GetNameReferenceCount(nuParse(t, "goo.foo"), "goo")).IsEqualTo(1)
	nuAssertThat(t, NodeUtil_GetNameReferenceCount(nuParse(t, "goo.foo"), "foo")).IsEqualTo(0)
	nuAssertThat(t, NodeUtil_GetNameReferenceCount(nuParse(t, "function foo(){}"), "foo")).IsEqualTo(1)
	nuAssertThat(t, NodeUtil_GetNameReferenceCount(nuParse(t, "var foo = function(){}"), "foo")).IsEqualTo(1)
}

func TestNodeUtil_AssortedTests_TestGetVarsDeclaredInBranch(t *testing.T) {
	nuh_assertNodeNames(t, []string{"foo"}, NodeUtil_GetVarsDeclaredInBranch(nuParse(t, "var foo;")))
	nuh_assertNodeNames(t,
		[]string{"foo", "goo"}, NodeUtil_GetVarsDeclaredInBranch(nuParse(t, "var foo,goo;")))
	nuh_assertNodeNames(t, []string{}, NodeUtil_GetVarsDeclaredInBranch(nuParse(t, "foo();")))
	nuh_assertNodeNames(t,
		[]string{},
		NodeUtil_GetVarsDeclaredInBranch(nuParse(t, "function f(){var foo;}")))
	nuh_assertNodeNames(t,
		[]string{"goo"},
		NodeUtil_GetVarsDeclaredInBranch(nuParse(t, "var goo;function f(){var foo;}")))
}

func nuh_assertNodeNames(t *testing.T, nodeNames []string, nodes []rhino.Node) {
	// Go-only: test-helper — Set.equals compares the names in any order.
	var actualNames []string
	for _, node := range nodes {
		actualNames = append(actualNames, node.GetString())
	}
	assertSameElements(t, "names", actualNames, nodeNames)
}

func TestNodeUtil_AssortedTests_TestIsNonlocalModuleExportNameOnExports1(t *testing.T) {
	root := nuParse(t, "export {localName as exportName};")
	moduleBody := root.GetFirstChild()
	exportNode := moduleBody.GetFirstChild()
	exportSpecs := exportNode.GetFirstChild()
	exportSpec := exportSpecs.GetFirstChild()

	localName := exportSpec.GetFirstChild()
	exportName := exportSpec.GetSecondChild()

	nuAssertThat(t, NodeUtil_IsNonlocalModuleExportName(localName)).IsFalse()
	nuAssertThat(t, NodeUtil_IsNonlocalModuleExportName(exportName)).IsTrue()
}

func TestNodeUtil_AssortedTests_TestIsNonlocalModuleExportNameOnExports2(t *testing.T) {
	root := nuParse(t, "let bar; export {bar};")
	moduleBody := root.GetFirstChild()
	exportNode := moduleBody.GetSecondChild()
	exportSpecs := exportNode.GetFirstChild()
	exportSpec := exportSpecs.GetFirstChild()

	name := exportSpec.GetFirstChild()

	// bar is defined locally, so isNonlocalModuleExportName is false.
	nuAssertThat(t, NodeUtil_IsNonlocalModuleExportName(name)).IsFalse()
}

func TestNodeUtil_AssortedTests_TestIsNonlocalModuleExportNameOnExportFrom(t *testing.T) {
	root := nuParse(t, "export {bar as baz} from './foo.js';")
	moduleBody := root.GetFirstChild()
	exportNode := moduleBody.GetFirstChild()
	exportSpecs := exportNode.GetFirstChild()
	exportSpec := exportSpecs.GetFirstChild()

	bar := exportSpec.GetFirstChild()
	baz := exportSpec.GetSecondChild()

	// Neither identifier is defined locally.
	nuAssertThat(t, NodeUtil_IsNonlocalModuleExportName(bar)).IsTrue()
	nuAssertThat(t, NodeUtil_IsNonlocalModuleExportName(baz)).IsTrue()
}

func TestNodeUtil_AssortedTests_TestIsNonlocalModuleExportNameOnImports1(t *testing.T) {
	root := nuParse(t, "import {exportName as localName} from './foo.js';")
	moduleBody := root.GetFirstChild()
	importNode := moduleBody.GetFirstChild()
	importSpecs := importNode.GetSecondChild()
	importSpec := importSpecs.GetFirstChild()

	exportName := importSpec.GetFirstChild()
	localName := importSpec.GetSecondChild()

	nuAssertThat(t, NodeUtil_IsNonlocalModuleExportName(exportName)).IsTrue()
	nuAssertThat(t, NodeUtil_IsNonlocalModuleExportName(localName)).IsFalse()
}

func TestNodeUtil_AssortedTests_TestIsNonlocalModuleExportNameOnImports2(t *testing.T) {
	root := nuParse(t, "import {bar} from './foo.js';")
	moduleBody := root.GetFirstChild()
	importNode := moduleBody.GetFirstChild()
	importSpecs := importNode.GetSecondChild()
	importSpec := importSpecs.GetFirstChild()

	name := importSpec.GetSecondChild()

	// bar is defined locally so isNonlocalModuleExportName is false
	nuAssertThat(t, NodeUtil_IsNonlocalModuleExportName(name)).IsFalse()
}

func TestNodeUtil_AssortedTests_TestIsControlStructureCodeBlock(t *testing.T) {
	root := nuParse(t, "if (x) foo(); else boo();")
	ifNode := root.GetFirstChild()

	ifCondition := ifNode.GetFirstChild()
	ifCase := ifNode.GetSecondChild()
	elseCase := ifNode.GetLastChild()

	nuAssertThat(t, NodeUtil_IsControlStructureCodeBlock(ifNode, ifCondition)).IsFalse()
	nuAssertThat(t, NodeUtil_IsControlStructureCodeBlock(ifNode, ifCase)).IsTrue()
	nuAssertThat(t, NodeUtil_IsControlStructureCodeBlock(ifNode, elseCase)).IsTrue()
}

func TestNodeUtil_AssortedTests_TestIsFunctionExpression1(t *testing.T) {
	root := nuParse(t, "(function foo() {})")
	statementNode := root.GetFirstChild()
	nuAssertThat(t, statementNode.IsExprResult()).IsTrue()
	functionNode := statementNode.GetFirstChild()
	nuAssertThat(t, functionNode.IsFunction()).IsTrue()
	nuAssertThat(t, NodeUtil_IsFunctionExpression(functionNode)).IsTrue()
}

func TestNodeUtil_AssortedTests_TestIsFunctionExpression2(t *testing.T) {
	root := nuParse(t, "function foo() {}")
	functionNode := root.GetFirstChild()
	nuAssertThat(t, functionNode.IsFunction()).IsTrue()
	nuAssertThat(t, NodeUtil_IsFunctionExpression(functionNode)).IsFalse()
}

func TestNodeUtil_AssortedTests_TestRemoveChildBlock(t *testing.T) {
	// Test removing the inner block.
	actual := nuParse(t, "{{x()}}")

	outerBlockNode := actual.GetFirstChild()
	innerBlockNode := outerBlockNode.GetFirstChild()

	NodeUtil_RemoveChild(outerBlockNode, innerBlockNode)
	expected := "{{}}"
	nuAssertNode(t, nuParse(t, expected)).IsEqualTo(actual)
}

func TestNodeUtil_AssortedTests_TestRemoveTryChild1(t *testing.T) {
	// Test removing the finally clause.
	actual := nuParse(t, "try {foo()} catch(e) {} finally {}")

	tryNode := actual.GetFirstChild()
	finallyBlock := tryNode.GetLastChild()

	NodeUtil_RemoveChild(tryNode, finallyBlock)
	expected := "try {foo()} catch(e) {}"
	nuAssertNode(t, nuParse(t, expected)).IsEqualTo(actual)
}

func TestNodeUtil_AssortedTests_TestRemoveTryChild2(t *testing.T) {
	// Test removing the try clause.
	actual := nuParse(t, "try {foo()} catch(e) {} finally {}")

	tryNode := actual.GetFirstChild()
	tryBlock := tryNode.GetFirstChild()

	NodeUtil_RemoveChild(tryNode, tryBlock)
	expected := "try {} catch(e) {} finally {}"
	nuAssertNode(t, nuParse(t, expected)).IsEqualTo(actual)
}

func TestNodeUtil_AssortedTests_TestRemoveTryChild3(t *testing.T) {
	// Test removing the catch clause.
	actual := nuParse(t, "try {foo()} catch(e) {} finally {}")

	tryNode := actual.GetFirstChild()
	catchBlocks := tryNode.GetSecondChild()
	catchBlock := catchBlocks.GetFirstChild()

	NodeUtil_RemoveChild(catchBlocks, catchBlock)
	expected := "try {foo()} finally {}"
	nuAssertNode(t, nuParse(t, expected)).IsEqualTo(actual)
}

func TestNodeUtil_AssortedTests_TestRemoveTryChild4(t *testing.T) {
	// Test removing the block that contains the catch clause.
	actual := nuParse(t, "try {foo()} catch(e) {} finally {}")

	tryNode := actual.GetFirstChild()
	catchBlocks := tryNode.GetSecondChild()

	NodeUtil_RemoveChild(tryNode, catchBlocks)
	expected := "try {foo()} finally {}"
	nuAssertNode(t, nuParse(t, expected)).IsEqualTo(actual)
}

func TestNodeUtil_AssortedTests_TestRemoveFromImport(t *testing.T) {
	// Remove imported function
	actual := nuParse(t, "import foo from './foo';")
	moduleBody := actual.GetFirstChild()
	importNode := moduleBody.GetFirstChild()
	functionFoo := importNode.GetFirstChild()

	NodeUtil_RemoveChild(importNode, functionFoo)
	expected := "import './foo';"
	nuAssertNode(t, nuParse(t, expected)).IsEqualTo(actual)
}

func TestNodeUtil_AssortedTests_TestRemoveParamChild1(t *testing.T) {
	// Remove traditional parameter
	actual := nuParse(t, "function f(p1) {}")
	functionNode := actual.GetFirstChild()
	paramList := functionNode.GetFirstChild().GetNext()
	p1 := paramList.GetFirstChild()

	NodeUtil_RemoveChild(paramList, p1)
	expected := "function f() {}"
	nuAssertNode(t, nuParse(t, expected)).IsEqualTo(actual)
}

func TestNodeUtil_AssortedTests_TestRemoveParamChild2(t *testing.T) {
	// Remove default parameter
	actual := nuParse(t, "function f(p1 = 0, p2) {}")
	functionNode := actual.GetFirstChild()
	paramList := functionNode.GetFirstChild().GetNext()
	p1 := paramList.GetFirstChild()

	NodeUtil_RemoveChild(paramList, p1)
	expected := "function f(p2) {}"
	nuAssertNode(t, nuParse(t, expected)).IsEqualTo(actual)
}

func TestNodeUtil_AssortedTests_TestRemoveVarChild(t *testing.T) {
	// Test removing the first child.
	actual := nuParse(t, "var foo, goo, hoo")

	varNode := actual.GetFirstChild()
	nameNode := varNode.GetFirstChild()

	NodeUtil_RemoveChild(varNode, nameNode)
	expected := "var goo, hoo"
	nuAssertNode(t, nuParse(t, expected)).IsEqualTo(actual)

	// Test removing the second child.
	actual = nuParse(t, "var foo, goo, hoo")

	varNode = actual.GetFirstChild()
	nameNode = varNode.GetSecondChild()

	NodeUtil_RemoveChild(varNode, nameNode)
	expected = "var foo, hoo"
	nuAssertNode(t, nuParse(t, expected)).IsEqualTo(actual)

	// Test removing the last child of several children.
	actual = nuParse(t, "var foo, hoo")

	varNode = actual.GetFirstChild()
	nameNode = varNode.GetSecondChild()

	NodeUtil_RemoveChild(varNode, nameNode)
	expected = "var foo"
	nuAssertNode(t, nuParse(t, expected)).IsEqualTo(actual)

	// Test removing the last.
	actual = nuParse(t, "var hoo")

	varNode = actual.GetFirstChild()
	nameNode = varNode.GetFirstChild()

	NodeUtil_RemoveChild(varNode, nameNode)
	expected = ""
	nuAssertNode(t, nuParse(t, expected)).IsEqualTo(actual)
}

func TestNodeUtil_AssortedTests_TestRemoveLetChild(t *testing.T) {
	// Test removing the first child.
	actual := nuParse(t, "let foo, goo, hoo")

	letNode := actual.GetFirstChild()
	nameNode := letNode.GetFirstChild()

	NodeUtil_RemoveChild(letNode, nameNode)
	expected := "let goo, hoo"
	nuAssertNode(t, nuParse(t, expected)).IsEqualTo(actual)

	// Test removing the second child.
	actual = nuParse(t, "let foo, goo, hoo")

	letNode = actual.GetFirstChild()
	nameNode = letNode.GetSecondChild()

	NodeUtil_RemoveChild(letNode, nameNode)
	expected = "let foo, hoo"
	nuAssertNode(t, nuParse(t, expected)).IsEqualTo(actual)

	// Test removing the last child of several children.
	actual = nuParse(t, "let foo, hoo")

	letNode = actual.GetFirstChild()
	nameNode = letNode.GetSecondChild()

	NodeUtil_RemoveChild(letNode, nameNode)
	expected = "let foo"
	nuAssertNode(t, nuParse(t, expected)).IsEqualTo(actual)

	// Test removing the last.
	actual = nuParse(t, "let hoo")

	letNode = actual.GetFirstChild()
	nameNode = letNode.GetFirstChild()

	NodeUtil_RemoveChild(letNode, nameNode)
	expected = ""
	nuAssertNode(t, nuParse(t, expected)).IsEqualTo(actual)
}

func TestNodeUtil_AssortedTests_TestRemoveLabelChild1(t *testing.T) {
	// Test removing the first child.
	actual := nuParse(t, "foo: goo()")

	labelNode := actual.GetFirstChild()
	callExpressNode := labelNode.GetLastChild()

	NodeUtil_RemoveChild(labelNode, callExpressNode)
	expected := ""
	nuAssertNode(t, nuParse(t, expected)).IsEqualTo(actual)
}

func TestNodeUtil_AssortedTests_TestRemoveLabelChild2(t *testing.T) {
	// Test removing the first child.
	actual := nuParse(t, "achoo: foo: goo()")

	labelNode := actual.GetFirstChild()
	callExpressNode := labelNode.GetLastChild()

	NodeUtil_RemoveChild(labelNode, callExpressNode)
	expected := ""
	nuAssertNode(t, nuParse(t, expected)).IsEqualTo(actual)
}

func TestNodeUtil_AssortedTests_TestRemovePatternChild(t *testing.T) {
	// remove variable declaration from object pattern
	actual := nuParse(t, "var {a, b, c} = {a:1, b:2, c:3}")
	varNode := actual.GetFirstChild()
	destructure := varNode.GetFirstChild()
	pattern := destructure.GetFirstChild()
	a := pattern.GetFirstChild()
	b := a.GetNext()
	c := b.GetNext()

	NodeUtil_RemoveChild(pattern, a)
	expected := "var {b, c} = {a:1, b:2, c:3};"
	nuAssertNode(t, nuParse(t, expected)).IsEqualTo(actual)

	// Remove all entries in object pattern
	NodeUtil_RemoveChild(pattern, b)
	NodeUtil_RemoveChild(pattern, c)
	expected = "var { } = {a:1, b:2, c:3};"
	nuAssertNode(t, nuParse(t, expected)).IsEqualTo(actual)

	// remove variable declaration from array pattern
	actual = nuParse(t, "var [a, b] = [1, 2]")
	varNode = actual.GetFirstChild()
	destructure = varNode.GetFirstChild()
	pattern = destructure.GetFirstChild()
	a = pattern.GetFirstChild()

	NodeUtil_RemoveChild(pattern, a)
	expected = "var [ , b] = [1, 2];"
	nuAssertNode(t, nuParse(t, expected)).IsEqualTo(actual)
}

func TestNodeUtil_AssortedTests_TestRemoveForChild(t *testing.T) {
	// Test removing the initializer.
	actual := nuParse(t, "for(var a=0;a<0;a++)foo()")

	forNode := actual.GetFirstChild()
	child := forNode.GetFirstChild()

	NodeUtil_RemoveChild(forNode, child)
	expected := "for(;a<0;a++)foo()"
	nuAssertNode(t, nuParse(t, expected)).IsEqualTo(actual)

	// Test removing the condition.
	actual = nuParse(t, "for(var a=0;a<0;a++)foo()")

	forNode = actual.GetFirstChild()
	child = forNode.GetSecondChild()

	NodeUtil_RemoveChild(forNode, child)
	expected = "for(var a=0;;a++)foo()"
	nuAssertNode(t, nuParse(t, expected)).IsEqualTo(actual)

	// Test removing the increment.
	actual = nuParse(t, "for(var a=0;a<0;a++)foo()")

	forNode = actual.GetFirstChild()
	child = forNode.GetSecondChild().GetNext()

	NodeUtil_RemoveChild(forNode, child)
	expected = "for(var a=0;a<0;)foo()"
	nuAssertNode(t, nuParse(t, expected)).IsEqualTo(actual)

	// Test removing the body.
	actual = nuParse(t, "for(var a=0;a<0;a++)foo()")

	forNode = actual.GetFirstChild()
	child = forNode.GetLastChild()

	NodeUtil_RemoveChild(forNode, child)
	expected = "for(var a=0;a<0;a++);"
	nuAssertNode(t, nuParse(t, expected)).IsEqualTo(actual)

	// Test removing the body.
	actual = nuParse(t, "for(a in ack)foo();")

	forNode = actual.GetFirstChild()
	child = forNode.GetLastChild()

	NodeUtil_RemoveChild(forNode, child)
	expected = "for(a in ack);"
	nuAssertNode(t, nuParse(t, expected)).IsEqualTo(actual)
}

func nuh_replaceDeclChild(t *testing.T, js string, declarationChild int, expected string) {
	actual := nuParse(t, js)
	declarationNode := actual.GetFirstChild()
	nameNode := declarationNode.GetChildAtIndex(declarationChild)

	NodeUtil_ReplaceDeclarationChild(nameNode, rhino.IR_Block())
	nuAssertNode(t, nuParse(t, expected)).IsEqualTo(actual)
}

func TestNodeUtil_AssortedTests_TestReplaceDeclarationName(t *testing.T) {
	nuh_replaceDeclChild(t, "var x;", 0, "{}")
	nuh_replaceDeclChild(t, "var x, y;", 0, "{} var y;")
	nuh_replaceDeclChild(t, "var x, y;", 1, "var x; {}")
	nuh_replaceDeclChild(t, "let x, y, z;", 0, "{} let y, z;")
	nuh_replaceDeclChild(t, "let x, y, z;", 1, "let x; {} let z;")
	nuh_replaceDeclChild(t, "let x, y, z;", 2, "let x, y; {}")
	nuh_replaceDeclChild(t, "const x = 1, y = 2, z = 3;", 1, "const x = 1; {} const z = 3;")
	nuh_replaceDeclChild(t,
		"const x =1, y = 2, z = 3, w = 4;", 1, "const x = 1; {} const z = 3, w = 4;")
}

func TestNodeUtil_AssortedTests_TestTryMergeBlock1(t *testing.T) {
	// Test removing the initializer.
	actual := nuParse(t, "{{a();b();}}")

	parentBlock := actual.GetFirstChild()
	childBlock := parentBlock.GetFirstChild()

	nuAssertThat(t, NodeUtil_TryMergeBlock(childBlock, false)).IsTrue()
	expected := "{a();b();}"
	nuAssertNode(t, nuParse(t, expected)).IsEqualTo(actual)
}

func TestNodeUtil_AssortedTests_TestTryMergeBlock2(t *testing.T) {
	// Test removing the initializer.
	actual := nuParse(t, "foo:{a();}")

	parentLabel := actual.GetFirstChild()
	childBlock := parentLabel.GetLastChild()

	nuAssertThat(t, NodeUtil_TryMergeBlock(childBlock, false)).IsFalse()
}

func TestNodeUtil_AssortedTests_TestTryMergeBlock3(t *testing.T) {
	// Test removing the initializer.
	code := "foo:{a();boo()}"
	actual := nuParse(t, "foo:{a();boo()}")

	parentLabel := actual.GetFirstChild()
	childBlock := parentLabel.GetLastChild()

	nuAssertThat(t, NodeUtil_TryMergeBlock(childBlock, false)).IsFalse()
	expected := code
	nuAssertNode(t, nuParse(t, expected)).IsEqualTo(actual)
}

func TestNodeUtil_AssortedTests_TestTryMergeBlock4(t *testing.T) {
	actual := nuParse(t, "{const module$exports$Foo=class{}}")
	expected := "const module$exports$Foo=class{}"

	block := actual.GetFirstChild()

	nuAssertThat(t, NodeUtil_TryMergeBlock(block, true)).IsTrue()
	nuAssertNode(t, nuParse(t, expected)).IsEqualTo(actual)
}

func TestNodeUtil_AssortedTests_TestCanMergeBlock1(t *testing.T) {
	actual := nuParse(t, "{a(); let x;}")

	block := actual.GetFirstChild()

	nuAssertThat(t, NodeUtil_CanMergeBlock(block)).IsFalse()
}

func TestNodeUtil_AssortedTests_TestCanMergeBlock2(t *testing.T) {
	actual := nuParse(t, "{a(); f(); var x; {const y = 2;}}")

	parentBlock := actual.GetFirstChild()
	childBlock := parentBlock.GetLastChild()

	nuAssertThat(t, NodeUtil_CanMergeBlock(parentBlock)).IsTrue()
	nuAssertThat(t, NodeUtil_CanMergeBlock(childBlock)).IsFalse()
}

func TestNodeUtil_AssortedTests_TestGetSourceName(t *testing.T) {
	n := rhino.NewNode(rhino.Token_BLOCK)
	parent := rhino.NewNodeWithChild(rhino.Token_BLOCK, n)
	parent.SetSourceFileForTesting("foo")

	nuAssertThat(t, NodeUtil_GetSourceName(n)).IsEqualTo("foo")
}

func TestNodeUtil_AssortedTests_TestLocalValue1(t *testing.T) {
	// Names are not known to be local.
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "x"))).IsFalse()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "x()"))).IsFalse()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "this"))).IsFalse()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "arguments"))).IsFalse()

	// We can't know if new objects are local unless we know
	// that they don't alias themselves.
	// TODO(tdeegan): Revisit this.
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "new x()"))).IsFalse()

	// property references are assumed to be non-local
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "(new x()).y"))).IsFalse()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "(new x())['y']"))).IsFalse()

	// Primitive values are local
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "null"))).IsTrue()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "undefined"))).IsTrue()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "Infinity"))).IsTrue()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "NaN"))).IsTrue()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "1"))).IsTrue()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "'a'"))).IsTrue()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "true"))).IsTrue()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "false"))).IsTrue()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "[]"))).IsTrue()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "{}"))).IsTrue()

	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "[x]"))).IsTrue()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "{'a':x}"))).IsTrue()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "{'a': {'b': 2}}"))).IsTrue()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "{'a': {'b': global}}"))).IsTrue()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "{get someGetter() { return 1; }}"))).IsTrue()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "{get someGetter() { return global; }}"))).IsTrue()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "{set someSetter(value) {}}"))).IsTrue()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "{[someComputedProperty]: {}}"))).IsTrue()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "{[someComputedProperty]: global}"))).IsTrue()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "{[someComputedProperty]: {'a':x}}"))).IsTrue()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "{[someComputedProperty]: {'a':1}}"))).IsTrue()

	// increment/decrement results in primitive number, the previous value is
	// always coersed to a number (even in the post.
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "++x"))).IsTrue()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "--x"))).IsTrue()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "x++"))).IsTrue()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "x--"))).IsTrue()

	// The left side of an only assign matters if it is an alias or mutable.
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "x=1"))).IsTrue()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "x=[]"))).IsFalse()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "x=y"))).IsFalse()
	// The right hand side of assignment opts don't matter, as they force
	// a local result.
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "x+=y"))).IsTrue()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "x*=y"))).IsTrue()
	// Comparisons always result in locals, as they force a local boolean
	// result.
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "x==y"))).IsTrue()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "x!=y"))).IsTrue()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "x>y"))).IsTrue()
	// Only the right side of a comma matters
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "(1,2)"))).IsTrue()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "(x,1)"))).IsTrue()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "(x,y)"))).IsFalse()

	// Both the operands of OR matter
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "1||2"))).IsTrue()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "x||1"))).IsFalse()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "x||y"))).IsFalse()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "1||y"))).IsFalse()

	// Both the operands of AND matter
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "1&&2"))).IsTrue()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "x&&1"))).IsFalse()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "x&&y"))).IsFalse()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "1&&y"))).IsFalse()

	// Logical assignments are always treated as non-local
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "x||=1"))).IsFalse()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "x&&=1"))).IsFalse()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "x??=1"))).IsFalse()

	// Only the results of HOOK matter
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "x?1:2"))).IsTrue()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "x?x:2"))).IsFalse()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "x?1:x"))).IsFalse()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "x?x:y"))).IsFalse()

	// Results of ops are local values
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "!y"))).IsTrue()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "~y"))).IsTrue()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "y + 1"))).IsTrue()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "y + z"))).IsTrue()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "y * z"))).IsTrue()

	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "'a' in x"))).IsTrue()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "typeof x"))).IsTrue()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "x instanceof y"))).IsTrue()

	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "void x"))).IsTrue()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "void 0"))).IsTrue()

	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "{}.x"))).IsFalse()

	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "{}.toString()"))).IsTrue()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "o.toString()"))).IsTrue()

	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "o.valueOf()"))).IsFalse()

	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "delete a.b"))).IsTrue()

	nuAssertThat(t,
		NodeUtil_EvaluatesToLocalValue(
			nuParseFirst(t, rhino.Token_NEW_TARGET, "function f() { new.target; }"))).IsFalse()

	// Dynamic import always returns a new Promise, which acts as a local value, even though
	// the result of the Promise is definitely non-local.
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "import('./foo.js')"))).IsTrue()
}

func TestNodeUtil_AssortedTests_TestLocalValueTemplateLit(t *testing.T) {
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "`hello`"))).IsTrue()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "`hello ${name}`"))).IsTrue()
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "`${'name'}`"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsLiteralValue(nuParseExpr(t, "`${'name'}`"), false)).IsTrue()
	nuAssertThat(t, NodeUtil_IsImmutableValue(nuParseExpr(t, "`${'name'}`"))).IsTrue()
}

func TestNodeUtil_AssortedTests_TestLocalValueTaggedTemplateLit1(t *testing.T) {
	n := nuParseExpr(t, "tag`simple string`")
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(n)).IsFalse()

	// Set 'tag' function as producing a local result
	flags := rhino.NewNode_SideEffectFlags()
	flags.ClearAllFlags()
	n.SetSideEffectFlagsObject(flags)

	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(n)).IsFalse()
}

func TestNodeUtil_AssortedTests_TestLocalValueTaggedTemplateLit2(t *testing.T) {
	// Here, replacement() may have side effects.
	n := nuParseExpr(t, "tag`string with ${replacement()}`")
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(n)).IsFalse()

	// Set 'tag' function as producing a local result
	flags := rhino.NewNode_SideEffectFlags()
	flags.ClearAllFlags()
	n.SetSideEffectFlagsObject(flags)

	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(n)).IsFalse()
}

func TestNodeUtil_AssortedTests_TestLocalValueNewExpr(t *testing.T) {
	newExpr := nuParseExpr(t, "new x()")
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(newExpr)).IsFalse()

	jutil.CheckState(newExpr.IsNew())
	flags := rhino.NewNode_SideEffectFlags()

	flags.ClearAllFlags()
	newExpr.SetSideEffectFlagsObject(flags)

	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(newExpr)).IsTrue()

	flags.ClearAllFlags()
	flags.SetMutatesThis()
	newExpr.SetSideEffectFlagsObject(flags)

	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(newExpr)).IsTrue()

	flags.ClearAllFlags()
	newExpr.SetSideEffectFlagsObject(flags)

	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(newExpr)).IsTrue()

	flags.ClearAllFlags()
	flags.SetThrows()
	newExpr.SetSideEffectFlagsObject(flags)

	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(newExpr)).IsFalse()

	flags.ClearAllFlags()
	flags.SetMutatesArguments()
	newExpr.SetSideEffectFlagsObject(flags)

	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(newExpr)).IsFalse()

	flags.ClearAllFlags()
	flags.SetMutatesGlobalState()
	newExpr.SetSideEffectFlagsObject(flags)

	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(newExpr)).IsFalse()
}

func TestNodeUtil_AssortedTests_TestLocalValueSpread(t *testing.T) {
	// Array literals are always known local values.
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "[...x]"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsLiteralValue(nuParseExpr(t, "[...x]"), false)).IsFalse()
	nuAssertThat(t, NodeUtil_IsImmutableValue(nuParseExpr(t, "[...x]"))).IsFalse()

	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "{...x}"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsLiteralValue(nuParseExpr(t, "{...x}"), false)).IsFalse()
	nuAssertThat(t, NodeUtil_IsImmutableValue(nuParseExpr(t, "{...x}"))).IsFalse()
}

func TestNodeUtil_AssortedTests_TestLocalValueAwait(t *testing.T) {
	var expr rhino.Node
	expr = nuParseFirst(t, rhino.Token_AWAIT, "async function f() { await someAsyncAction(); }")
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(expr)).IsFalse()

	expr = nuParseFirst(t, rhino.Token_AWAIT, "async function f() { await {then:function() { return p }}; }")
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(expr)).IsFalse()

	// it isn't clear why someone would want to wait on a non-thenable value...
	expr = nuParseFirst(t, rhino.Token_AWAIT, "async function f() { await 5; }")
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(expr)).IsFalse()
}

func TestNodeUtil_AssortedTests_TestLocalValueYield(t *testing.T) {
	var expr rhino.Node
	expr = nuParseFirst(t, rhino.Token_YIELD, "function *f() { yield; }")
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(expr)).IsFalse()

	expr = nuParseFirst(t, rhino.Token_YIELD, "function *f() { yield 'something'; }")
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(expr)).IsFalse()
}

func TestNodeUtil_AssortedTests_LocalValueOptChainGetProp(t *testing.T) {
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "x?.y"))).IsFalse()
}

func TestNodeUtil_AssortedTests_LocalValueOptChainGetElem(t *testing.T) {
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "x?.[y]"))).IsFalse()
}

func TestNodeUtil_AssortedTests_LocalValueOptChainCall(t *testing.T) {
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "x?.()"))).IsFalse()
}

func TestNodeUtil_AssortedTests_LocalValueOptChainCall_toString(t *testing.T) {
	// the toString() call is an optional call because it is part of an optional chain
	nuAssertThat(t, NodeUtil_EvaluatesToLocalValue(nuParseExpr(t, "x?.toString()"))).IsTrue()
}

func TestNodeUtil_AssortedTests_TestGetOctalNumberValue(t *testing.T) {
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "022"))).IsEqualTo(18.0)
}

func TestNodeUtil_AssortedTests_TestGetStringNumberValue(t *testing.T) {
	// Whitespace and empty
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("")).IsEqualTo(0.0)
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("  ")).IsEqualTo(0.0)
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("\t\n\r ")).IsEqualTo(0.0)
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue(" \t 42 \n ")).IsEqualTo(42.0)

	// Standard decimals
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("0")).IsEqualTo(0.0)
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("-0")).IsEqualTo(math.Copysign(0, -1)) // Java: -0.0
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("1")).IsEqualTo(1.0)
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("-1")).IsEqualTo(-1.0)
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("+1")).IsEqualTo(1.0)
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("1.5")).IsEqualTo(1.5)
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("-1.5")).IsEqualTo(-1.5)
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue(".5")).IsEqualTo(0.5)
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("-.5")).IsEqualTo(-0.5)
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("1e5")).IsEqualTo(100000.0)
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("1e-2")).IsEqualTo(0.01)

	// Binary literals
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("0b10")).IsEqualTo(2.0)
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("0B11")).IsEqualTo(3.0)
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("0b0")).IsEqualTo(0.0)
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("  0b1010  ")).IsEqualTo(10.0)

	// Binary invalid
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("0b")).IsEqualTo(math.NaN())
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("0b2")).IsEqualTo(math.NaN())
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("0b1a")).IsEqualTo(math.NaN())
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("+0b1")).IsEqualTo(math.NaN())
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("-0b1")).IsEqualTo(math.NaN())

	// Octal literals
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("0o10")).IsEqualTo(8.0)
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("0O77")).IsEqualTo(63.0)
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("0o0")).IsEqualTo(0.0)
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("  0o12  ")).IsEqualTo(10.0)

	// Octal invalid
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("0o")).IsEqualTo(math.NaN())
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("0o8")).IsEqualTo(math.NaN())
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("0o7a")).IsEqualTo(math.NaN())
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("+0o1")).IsEqualTo(math.NaN())
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("-0o1")).IsEqualTo(math.NaN())

	// Hexadecimal
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("0x10")).IsEqualTo(16.0)
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("0X1f")).IsEqualTo(31.0)
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("0xffffffff")).IsEqualTo(4294967295.0)
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("0x0")).IsEqualTo(0.0)

	// Large non-decimal literals (> 14 digits exercising BigInteger path)
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("0b10000000000000000000")).IsEqualTo(524288.0)
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("0o100000000000000")).IsEqualTo(4398046511104.0)
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("0x100000000000000")).
		IsEqualTo(7.205759403792794e+16)
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("0x"+strings.Repeat("f", 300))).IsEqualTo(math.Inf(1))

	// Hex invalid
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("0x")).IsEqualTo(math.NaN())
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("0xg")).IsEqualTo(math.NaN())
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("0x1.0p0")).IsEqualTo(math.NaN())

	// Hex with signs
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("+0x10")).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("-0x10")).IsNull()

	// Java suffixes rejected
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("1d")).IsEqualTo(math.NaN())
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("1D")).IsEqualTo(math.NaN())
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("1f")).IsEqualTo(math.NaN())
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("1F")).IsEqualTo(math.NaN())
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("1.0d")).IsEqualTo(math.NaN())
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("1.0f")).IsEqualTo(math.NaN())
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("1e2d")).IsEqualTo(math.NaN())
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("1e2f")).IsEqualTo(math.NaN())

	// Infinity
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("Infinity")).IsEqualTo(math.Inf(1))
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("+Infinity")).IsEqualTo(math.Inf(1))
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("-Infinity")).IsEqualTo(math.Inf(-1))

	// Case-insensitive infinity
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("infinity")).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("-infinity")).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("+infinity")).IsNull()

	// Vertical tab
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("\u000b")).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("\u000b1")).IsNull()

	// General invalid
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("abc")).IsEqualTo(math.NaN())
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("1.2.3")).IsEqualTo(math.NaN())
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("--1")).IsEqualTo(math.NaN())
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("+")).IsEqualTo(math.NaN())
	nuAssertThatDouble(t, NodeUtil_GetStringNumberValue("-")).IsEqualTo(math.NaN())
}

func TestNodeUtil_AssortedTests_TestGetNumberValue(t *testing.T) {
	// Strings
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "'\\uFEFF1'"))).IsEqualTo(1.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "''"))).IsEqualTo(0.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "' '"))).IsEqualTo(0.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "' \\t'"))).IsEqualTo(0.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "'+0'"))).IsEqualTo(0.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "'-0'"))).IsEqualTo(math.Copysign(0, -1)) // Java: -0.0
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "'+2'"))).IsEqualTo(2.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "'-1.6'"))).IsEqualTo(-1.6)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "'16'"))).IsEqualTo(16.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "' 16 '"))).IsEqualTo(16.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "'123e2'"))).IsEqualTo(12300.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "'123E2'"))).IsEqualTo(12300.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "'123e-2'"))).IsEqualTo(1.23)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "'123E-2'"))).IsEqualTo(1.23)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "'-123e-2'"))).IsEqualTo(-1.23)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "'-123E-2'"))).IsEqualTo(-1.23)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "'+123e-2'"))).IsEqualTo(1.23)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "'+123E-2'"))).IsEqualTo(1.23)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "'+123e+2'"))).IsEqualTo(12300.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "'+123E+2'"))).IsEqualTo(12300.0)

	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "'0xf'"))).IsEqualTo(15.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "'0xF'"))).IsEqualTo(15.0)

	// Chrome and rhino behavior differently from FF and IE. FF and IE
	// consider a negative hex number to be invalid
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "'-0xf'"))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "'-0xF'"))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "'+0xf'"))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "'+0xF'"))).IsNull()

	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "'0X10'"))).IsEqualTo(16.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "'0X10.8'"))).IsEqualTo(math.NaN())
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "'077'"))).IsEqualTo(77.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "'-077'"))).IsEqualTo(-77.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "'-077.5'"))).IsEqualTo(-77.5)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "'-Infinity'"))).IsEqualTo(math.Inf(-1))
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "'Infinity'"))).IsEqualTo(math.Inf(1))
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "'+Infinity'"))).IsEqualTo(math.Inf(1))

	// Firefox treats "infinity" as "Infinity", IE treats it as NaN
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "'-infinity'"))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "'infinity'"))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "'+infinity'"))).IsNull()

	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "'NaN'"))).IsEqualTo(math.NaN())
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "'some unknown string'"))).IsEqualTo(math.NaN())
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "'123 blah'"))).IsEqualTo(math.NaN())
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "\"0o10\""))).IsEqualTo(8.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "\"1d\""))).IsEqualTo(math.NaN())
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "\"1f\""))).IsEqualTo(math.NaN())

	// Literals
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "1"))).IsEqualTo(1.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "1n"))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "-1"))).IsEqualTo(-1.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "+1"))).IsEqualTo(1.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "22"))).IsEqualTo(22.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "022"))).IsEqualTo(18.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "0x22"))).IsEqualTo(34.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "-0.1"))).IsEqualTo(-0.1)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "-0.0"))).IsEqualTo(math.Copysign(0, -1)) // Java: -0.0

	// BITNOT
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "~1"))).IsEqualTo(-2.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "~-1"))).IsEqualTo(0.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "~22"))).IsEqualTo(-23.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "~022"))).IsEqualTo(-19.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "~0.0"))).IsEqualTo(-1.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "~0.1"))).IsEqualTo(-1.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "~NaN"))).IsEqualTo(-1.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "~Infinity"))).IsEqualTo(-1.0)

	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "true"))).IsEqualTo(1.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "false"))).IsEqualTo(0.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "null"))).IsEqualTo(0.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "void 0"))).IsEqualTo(math.NaN())
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "void f"))).IsEqualTo(math.NaN())
	// we pay no attention to possible side effects.
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "void f()"))).IsEqualTo(math.NaN())
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "NaN"))).IsEqualTo(math.NaN())
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "Infinity"))).IsEqualTo(math.Inf(1))
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "-Infinity"))).IsEqualTo(math.Inf(-1))

	// "infinity" is not a known name.
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "infinity"))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "-infinity"))).IsNull()

	// getNumberValue only converts literals
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "x"))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "x.y"))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "1/2"))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "1-2"))).IsNull()

	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "[1]"))).IsEqualTo(1.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "{}"))).IsEqualTo(math.NaN())
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "{a: 1}"))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "{'a': 1, 2: 3}"))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "{foo() {}}"))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "{get foo() { return 1; }}"))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "{set foo(x) {}}"))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "{valueOf: 1}"))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "{toString: 1}"))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "{valueOf() { return 1; }}"))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "{toString() { return '1'; }}"))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "{get valueOf() { return 1; }}"))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "{set valueOf(x) {}}"))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "{__proto__: null}"))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "{'__proto__': null}"))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "{[key]: 1}"))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "{...spread}"))).IsNull()

	// Unary with conversions
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "-\"1\""))).IsEqualTo(-1.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "+\"1\""))).IsEqualTo(1.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "~\"1\""))).IsEqualTo(-2.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "-true"))).IsEqualTo(-1.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "+true"))).IsEqualTo(1.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "~true"))).IsEqualTo(-2.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "-false"))).IsEqualTo(math.Copysign(0, -1)) // Java: -0.0
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "+false"))).IsEqualTo(0.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "~false"))).IsEqualTo(-1.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "-[1]"))).IsEqualTo(-1.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "+[1]"))).IsEqualTo(1.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "~[1]"))).IsEqualTo(-2.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "-{}"))).IsEqualTo(math.NaN())
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "+{}"))).IsEqualTo(math.NaN())
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "~{}"))).IsEqualTo(-1.0)

	// Nested unaries with conversions
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "- -\"1\""))).IsEqualTo(1.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "+ +\"1\""))).IsEqualTo(1.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "+ -\"1\""))).IsEqualTo(-1.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "~ ~\"1\""))).IsEqualTo(1.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "- -1"))).IsEqualTo(1.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "+ +1"))).IsEqualTo(1.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "+ -1"))).IsEqualTo(-1.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValue(nuParseExpr(t, "~ ~1"))).IsEqualTo(1.0)
}

func TestNodeUtil_AssortedTests_TestGetNumberValueNoConversions(t *testing.T) {
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "''"))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "``"))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "true"))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "false"))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "null"))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "undefined"))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "[1]"))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "{}"))).IsNull()

	// Unary operators without conversions should return null for non-numeric operands
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "-\"1\""))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "-\"1d\""))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "-true"))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "-false"))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "-[1]"))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "-{}"))).IsNull()

	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "+\"1\""))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "+\"1d\""))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "+true"))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "+false"))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "+[1]"))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "+{}"))).IsNull()

	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "~\"1\""))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "~true"))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "~false"))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "~[1]"))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "~{}"))).IsNull()

	// Nested unaries without conversions
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "- -\"1\""))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "+ +\"1\""))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "+ -\"1\""))).IsNull()

	// Literals
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "1"))).IsEqualTo(1.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "1n"))).IsNull()
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "-1"))).IsEqualTo(-1.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "+1"))).IsEqualTo(1.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "22"))).IsEqualTo(22.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "022"))).IsEqualTo(18.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "0x22"))).IsEqualTo(34.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "-0.1"))).IsEqualTo(-0.1)
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "-0.0"))).IsEqualTo(math.Copysign(0, -1)) // Java: -0.0

	// BITNOT
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "~1"))).IsEqualTo(-2.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "~-1"))).IsEqualTo(0.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "~22"))).IsEqualTo(-23.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "~022"))).IsEqualTo(-19.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "~0.0"))).IsEqualTo(-1.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "~0.1"))).IsEqualTo(-1.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "~NaN"))).IsEqualTo(-1.0)
	nuAssertThatDouble(t, NodeUtil_GetNumberValueNoConversions(nuParseExpr(t, "~Infinity"))).IsEqualTo(-1.0)
}

func TestNodeUtil_GetNumberValueOpt(t *testing.T) {
	exprs := []string{
		"0", "-0", "42", "-1.5", "0x22", "NaN", "Infinity", "-Infinity",
		"undefined", "null", "true", "false", `"123"`, `" 0b1010 "`, `"abc"`, `"infinity"`,
		"1n", "{}", "[1]", "x",
	}
	for _, src := range exprs {
		n := nuParseExpr(t, src)
		gotV, gotOK := NodeUtil_GetNumberValueOpt(n)
		wantPtr := NodeUtil_GetNumberValue(n)
		if (wantPtr != nil) != gotOK {
			t.Fatalf("GetNumberValueOpt(%s) ok=%v, want %v", src, gotOK, wantPtr != nil)
		}
		if gotOK && math.Float64bits(gotV) != math.Float64bits(*wantPtr) && !(math.IsNaN(gotV) && math.IsNaN(*wantPtr)) {
			t.Fatalf("GetNumberValueOpt(%s) = %v, want %v", src, gotV, *wantPtr)
		}

		ncV, ncOK := NodeUtil_GetNumberValueNoConversionsOpt(n)
		wantNcPtr := NodeUtil_GetNumberValueNoConversions(n)
		if (wantNcPtr != nil) != ncOK {
			t.Fatalf("GetNumberValueNoConversionsOpt(%s) ok=%v, want %v", src, ncOK, wantNcPtr != nil)
		}
		if ncOK && math.Float64bits(ncV) != math.Float64bits(*wantNcPtr) && !(math.IsNaN(ncV) && math.IsNaN(*wantNcPtr)) {
			t.Fatalf("GetNumberValueNoConversionsOpt(%s) = %v, want %v", src, ncV, *wantNcPtr)
		}
	}
	for _, raw := range []string{"", " 42 ", "-0", "0x10", "0b", "+0x1", "infinity", "-Infinity", "\v1", "notanumber"} {
		gotV, gotOK := NodeUtil_GetStringNumberValueOpt(raw)
		wantPtr := NodeUtil_GetStringNumberValue(raw)
		if (wantPtr != nil) != gotOK {
			t.Fatalf("GetStringNumberValueOpt(%q) ok=%v, want %v", raw, gotOK, wantPtr != nil)
		}
		if gotOK && math.Float64bits(gotV) != math.Float64bits(*wantPtr) && !(math.IsNaN(gotV) && math.IsNaN(*wantPtr)) {
			t.Fatalf("GetStringNumberValueOpt(%q) = %v, want %v", raw, gotV, *wantPtr)
		}
	}
}

func TestNodeUtil_AssortedTests_TestIsNumericLiteral(t *testing.T) {
	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "1"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "0"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "0.0"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "1.5"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "1e5"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "0x10"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "0b10"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "0o10"))).IsTrue()

	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "-1"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "-0"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "-0.0"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "-1.5"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "-1e5"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "-0x10"))).IsTrue()

	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "+1"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "+0"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "+0.0"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "+1.5"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "+1e5"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "+0x10"))).IsTrue()

	// False cases
	nuAssertThat(t, NodeUtil_IsNumericLiteral(rhino.NilNode)).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "null"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "undefined"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "void 0"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "NaN"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "Infinity"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "-Infinity"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "+Infinity"))).IsFalse()

	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "1n"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "-1n"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(rhino.NewNodeWithChild(rhino.Token_POS, rhino.Node_NewBigInt(big.NewInt(1))))).
		IsFalse()

	nuAssertThat(t, NodeUtil_IsNumericLiteral(rhino.NewNode(rhino.Token_NEG))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(rhino.NewNode(rhino.Token_POS))).IsFalse()

	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "\"1\""))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "-\"1\""))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "+\"1\""))).IsFalse()

	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "true"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "-true"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "+true"))).IsFalse()

	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "false"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "-false"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "+false"))).IsFalse()

	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "{}"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "[]"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "[1]"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "-[1]"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "+[1]"))).IsFalse()

	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "1 + 2"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "-(1 + 2)"))).IsFalse()

	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "-(-1)"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "+(+1)"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "+(-1)"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "-(+1)"))).IsFalse()

	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "~1"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericLiteral(nuParseExpr(t, "!1"))).IsFalse()

	nuAssertThat(t, NodeUtil_GetNumericLiteralValue(nuParseExpr(t, "1"))).IsEqualTo(1.0)
}

func TestNodeUtil_AssortedTests_TestGetNumericLiteralValue(t *testing.T) {
	nuAssertThat(t, NodeUtil_GetNumericLiteralValue(nuParseExpr(t, "1"))).IsEqualTo(1.0)
	nuAssertThat(t, NodeUtil_GetNumericLiteralValue(nuParseExpr(t, "0.0"))).IsEqualTo(0.0)
	nuAssertThat(t, NodeUtil_GetNumericLiteralValue(nuParseExpr(t, "-5.0"))).IsEqualTo(-5.0)
	nuAssertThat(t, NodeUtil_GetNumericLiteralValue(nuParseExpr(t, "+5.0"))).IsEqualTo(5.0)

	negZero := NodeUtil_GetNumericLiteralValue(nuParseExpr(t, "-0.0"))
	nuAssertThat(t, negZero).IsEqualTo(math.Copysign(0, -1)) // Java: -0.0
	nuAssertThat(t, 1.0/negZero).IsEqualTo(math.Inf(-1))

	posZero := NodeUtil_GetNumericLiteralValue(nuParseExpr(t, "+0.0"))
	nuAssertThat(t, posZero).IsEqualTo(0.0)
	nuAssertThat(t, 1.0/posZero).IsEqualTo(math.Inf(1))

	negIntZero := NodeUtil_GetNumericLiteralValue(nuParseExpr(t, "-0"))
	nuAssertThat(t, negIntZero).IsEqualTo(math.Copysign(0, -1)) // Java: -0.0
	nuAssertThat(t, 1.0/negIntZero).IsEqualTo(math.Inf(-1))

	posIntZero := NodeUtil_GetNumericLiteralValue(nuParseExpr(t, "+0"))
	nuAssertThat(t, posIntZero).IsEqualTo(0.0)
	nuAssertThat(t, 1.0/posIntZero).IsEqualTo(math.Inf(1))

	nuAssertThat(t, NodeUtil_GetNumericLiteralValue(nuParseExpr(t, "0x10"))).IsEqualTo(16.0)
	nuAssertThat(t, NodeUtil_GetNumericLiteralValue(nuParseExpr(t, "-0x10"))).IsEqualTo(-16.0)
	nuAssertThat(t, NodeUtil_GetNumericLiteralValue(nuParseExpr(t, "+0x10"))).IsEqualTo(16.0)
	nuAssertThat(t, NodeUtil_GetNumericLiteralValue(nuParseExpr(t, "0b10"))).IsEqualTo(2.0)
	nuAssertThat(t, NodeUtil_GetNumericLiteralValue(nuParseExpr(t, "0o10"))).IsEqualTo(8.0)
	nuAssertThat(t, NodeUtil_GetNumericLiteralValue(nuParseExpr(t, "1e5"))).IsEqualTo(100000.0)
	nuAssertThat(t, NodeUtil_GetNumericLiteralValue(nuParseExpr(t, "-1e5"))).IsEqualTo(-100000.0)
	nuAssertThat(t, NodeUtil_GetNumericLiteralValue(nuParseExpr(t, "+1e5"))).IsEqualTo(100000.0)

	assertThrows[jutil.IIllegalArgumentException](t, func() { NodeUtil_GetNumericLiteralValue(rhino.NilNode) })
	assertThrows[jutil.IIllegalArgumentException](t,
		func() { NodeUtil_GetNumericLiteralValue(nuParseExpr(t, "\"1\"")) })
	assertThrows[jutil.IIllegalArgumentException](t,
		func() { NodeUtil_GetNumericLiteralValue(nuParseExpr(t, "null")) })
	assertThrows[jutil.IIllegalArgumentException](t,
		func() { NodeUtil_GetNumericLiteralValue(nuParseExpr(t, "undefined")) })
	assertThrows[jutil.IIllegalArgumentException](t,
		func() { NodeUtil_GetNumericLiteralValue(nuParseExpr(t, "NaN")) })
	assertThrows[jutil.IIllegalArgumentException](t,
		func() { NodeUtil_GetNumericLiteralValue(nuParseExpr(t, "Infinity")) })
	assertThrows[jutil.IIllegalArgumentException](t,
		func() { NodeUtil_GetNumericLiteralValue(nuParseExpr(t, "-(1 + 2)")) })
	assertThrows[jutil.IIllegalArgumentException](t,
		func() { NodeUtil_GetNumericLiteralValue(nuParseExpr(t, "-(-1)")) })
	assertThrows[jutil.IIllegalArgumentException](t,
		func() { NodeUtil_GetNumericLiteralValue(nuParseExpr(t, "1n")) })
	assertThrows[jutil.IIllegalArgumentException](t,
		func() { NodeUtil_GetNumericLiteralValue(nuParseExpr(t, "-1n")) })
	assertThrows[jutil.IIllegalArgumentException](t,
		func() { NodeUtil_GetNumericLiteralValue(rhino.NewNode(rhino.Token_NEG)) })
	assertThrows[jutil.IIllegalArgumentException](t,
		func() { NodeUtil_GetNumericLiteralValue(rhino.NewNode(rhino.Token_POS)) })
}

func TestNodeUtil_AssortedTests_TestGetBigIntValue(t *testing.T) {
	// Strings
	// NOTE: Strings with the 'n' literal syntax can't be converted to BigInt
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "'1'"))).IsEqualTo("1")
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "'1n'"))).IsNull()

	// whitespace
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "'\\uFEFF1'"))).IsEqualTo("1")
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "''"))).IsEqualTo("0")
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "' '"))).IsEqualTo("0")
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "' \\t'"))).IsEqualTo("0")
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "' 1 '"))).IsEqualTo("1")

	// Even though unary plus cannot be applied to a bigint, numeric strings that have one can
	// still be converted, e.g. BigInt("+0") => 0n
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "'+0'"))).IsEqualTo("0")
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "'-0'"))).IsEqualTo("0")

	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "'10'"))).IsEqualTo("10")
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "'+10'"))).IsEqualTo("10")
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "'-10'"))).IsEqualTo("-10")
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "'   +10'"))).IsEqualTo("10")
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "'   -10'"))).IsEqualTo("-10")
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "'0xa'"))).IsEqualTo("10")
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "'0XA'"))).IsEqualTo("10")
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "'0o12'"))).IsEqualTo("10")
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "'0O12'"))).IsEqualTo("10")
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "'0b1010'"))).IsEqualTo("10")
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "'0B1010'"))).IsEqualTo("10")

	// Strings representing hex, octal, or binary literal values cannot include a unary plus or
	// minus when being passed to BigInt().
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "'+0xF'"))).IsNull()
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "'-0xF'"))).IsNull()
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "'+0oF'"))).IsNull()
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "'-0oF'"))).IsNull()
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "'+0bF'"))).IsNull()
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "'-0bF'"))).IsNull()

	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "'1.5'"))).IsNull()
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "'-Infinity'"))).IsNull()
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "'Infinity'"))).IsNull()
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "'+Infinity'"))).IsNull()

	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "'NaN'"))).IsNull()
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "'some unknown string'"))).IsNull()
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "'123 blah'"))).IsNull()

	// Literals
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "1"))).IsEqualTo("1")
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "1n"))).IsEqualTo("1")
	// "-1" is parsed as a literal
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "-1"))).IsEqualTo("-1")
	// "+1" is parse as an op + literal
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "+1"))).IsNull()
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "22"))).IsEqualTo("22")
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "022"))).IsEqualTo("18")
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "0x22"))).IsEqualTo("34")

	maxSafeInt := int64(9007199254740991) // Number.MAX_SAFE_INTEGER
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, jutil.LongToString(maxSafeInt)))).IsEqualTo(jutil.LongToString(maxSafeInt))
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, jutil.LongToString(-maxSafeInt)))).IsEqualTo(jutil.LongToString(-maxSafeInt))

	maxSafeIntPlusOne := maxSafeInt + 1
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, jutil.LongToString(maxSafeIntPlusOne)))).IsNull()
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, jutil.LongToString(-maxSafeIntPlusOne)))).IsNull()

	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "true"))).IsEqualTo("1")
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "false"))).IsEqualTo("0")
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "null"))).IsNull()
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "void 0"))).IsNull()
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "void f"))).IsNull()

	// getBigIntValue only converts literals
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "x"))).IsNull()
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "x.y"))).IsNull()
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "1/2"))).IsNull()
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "1-2"))).IsNull()
	nuAssertThatBigInt(t, NodeUtil_GetBigIntValue(nuParseExpr(t, "+1"))).IsNull()
}

func TestNodeUtil_AssortedTests_TestIsNumericResult(t *testing.T) {
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "1"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "true"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "+true"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "+1"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "-1"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "-Infinity"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "Infinity"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "NaN"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "undefined"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "void 0"))).IsFalse()

	// Technically the result when both operands are variables could be a `bigint`.
	// We would need to use type information to be certain, but `number` is much more likely.
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "a << b"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "a >> b"))).IsTrue()
	// if one side is a number then we can safely assume the result is a number
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "a << 2"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "a >> 2"))).IsTrue()
	// unsigned right shift is only allowed with numbers
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "a >>> b"))).IsTrue()

	// logical comparisons result in boolean, not number
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "a == b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "a != b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "a === b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "a !== b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "a < b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "a > b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "a <= b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "a >= b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "a in b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "a instanceof b"))).IsFalse()

	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "'a'"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "'a'+b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "a+'b'"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "a+b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "a()"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "''.a"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "a.b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "a.b()"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "a().b()"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "new a()"))).IsFalse()

	// Definitely not numeric
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "([1,2])"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "({a:1})"))).IsFalse()

	// Recurse into the expression when necessary.
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "1 && 2"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "1 || 2"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "a ? 2 : 3"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "a,1"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "a=1"))).IsTrue()

	// variable might be a string
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "a += 1"))).IsFalse()

	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "a -= 1"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "a *= 1"))).IsTrue()

	// Technically when the operand is a variable the result could be a `bigint`.
	// We would need to use type information to be certain, but `number` is much more likely.
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "--a"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "++a"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "a++"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsNumericResult(nuParseExpr(t, "a--"))).IsTrue()
}

func TestNodeUtil_AssortedTests_TestIsBigIntResult(t *testing.T) {
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "1n"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "-1n"))).IsTrue()

	// if one side is a bigint then we will assume the result is a bigint
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "a *= 1n"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "a * 1n"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "1n * a"))).IsTrue()

	// if neither side is explicitly bigint, then the the result could be a number or type error
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "a * b"))).IsFalse()
	// addition is a special case because of string concatenation
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "'a'+1n"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "1n+'b'"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "1n + 2n"))).IsTrue()
	// variable might be a string
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "a += 1n"))).IsFalse()

	// unsigned right shift and unary + are only allowed with numbers
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "+a"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "a >>> b"))).IsFalse()

	// recurse if necessary
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "a * b + 1n"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "a * b * 1n"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "a * b + c * 1n"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "a * b + c * d"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "a * ~b * -c * d * 1n"))).IsTrue()
	// even though the other side is definitely NUMBER, isBigIntResult returns true because it is
	// not responsible for type checking
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "a * 1 * b * c * d * 1n"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "a * b * +c * d * 1n"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "a * b >> c * d * 1n"))).IsTrue()
	// unsigned right shift is the overall parent and thus evaluated as NUMBER
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "a * b >>> c * d * 1n"))).IsFalse()

	// Infinity and NaN are numbers, not bigints
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "-Infinity"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "Infinity"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "NaN"))).IsFalse()

	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "undefined"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "void 0n"))).IsFalse()

	// logical comparisons result in boolean, not bigint
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "a == 1n"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "a != 1n"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "a === 1n"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "a !== 1n"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "a < 1n"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "a > 1n"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "a <= 1n"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "a >= 1n"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "a in 1n"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "a instanceof 1n"))).IsFalse()

	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "a()"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "''.a"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "a.b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "a.b()"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "a().b()"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "new a()"))).IsFalse()

	// Definitely not bigint
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "([1n,2n])"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "({a:1n})"))).IsFalse()

	// Recurse into the expression when necessary.
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "1n && 2n"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "1n || 2n"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "a ? 2n : 3n"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "a,1n"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "a=1n"))).IsTrue()

	// variable might be a number
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "--a"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "++a"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "a++"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBigIntResult(nuParseExpr(t, "a--"))).IsFalse()
}

func TestNodeUtil_AssortedTests_TestIsBooleanResult(t *testing.T) {
	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "1"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "true"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "+true"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "+1"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "-1"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "-Infinity"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "Infinity"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "NaN"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "undefined"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "void 0"))).IsFalse()

	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "a << b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "a >> b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "a >>> b"))).IsFalse()

	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "a == b"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "a != b"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "a === b"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "a !== b"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "a < b"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "a > b"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "a <= b"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "a >= b"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "a in b"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "a instanceof b"))).IsTrue()

	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "'a'"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "'a'+b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "a+'b'"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "a+b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "a()"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "''.a"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "a.b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "a.b()"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "a().b()"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "new a()"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "delete a"))).IsTrue()

	// Definitely not boolean
	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "([true,false])"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "({a:true})"))).IsFalse()

	// These are boolean
	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "true && false"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "true || false"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "a ? true : false"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "a,true"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "a=true"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsBooleanResult(nuParseExpr(t, "a=1"))).IsFalse()
}

func TestNodeUtil_AssortedTests_TestMayBeString(t *testing.T) {
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "1"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "1n"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "true"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "+true"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "+1"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "-1"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "-Infinity"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "Infinity"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "NaN"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "undefined"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "void 0"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "null"))).IsFalse()

	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "a << b"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "a >> b"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "a >>> b"))).IsFalse()

	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "a == b"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "a != b"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "a === b"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "a !== b"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "a < b"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "a > b"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "a <= b"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "a >= b"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "a in b"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "a instanceof b"))).IsFalse()

	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "'a'"))).IsTrue()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "'a'+b"))).IsTrue()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "a+'b'"))).IsTrue()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "a+b"))).IsTrue()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "a()"))).IsTrue()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "''.a"))).IsTrue()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "a.b"))).IsTrue()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "a.b()"))).IsTrue()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "a().b()"))).IsTrue()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "new a()"))).IsTrue()

	// These can't be strings but they aren't handled yet.
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "1 && 2"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "1 || 2"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "1 ? 2 : 3"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "1,2"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "a=1"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "1+1"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "true+true"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "null+null"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "NaN+NaN"))).IsFalse()

	// These are not strings but they aren't primitives either
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "([1,2])"))).IsTrue()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "({a:1})"))).IsTrue()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "({}+1)"))).IsTrue()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "(1+{})"))).IsTrue()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "([]+1)"))).IsTrue()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "(1+[])"))).IsTrue()

	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "a += 'x'"))).IsTrue()
	nuAssertThat(t, NodeUtil_MayBeString(nuParseExpr(t, "a += 1"))).IsTrue()
}

func TestNodeUtil_AssortedTests_TestMayBeStringWithTypes(t *testing.T) {
	registry := rhino.NewJSTypeRegistry(nil)
	nuAssertThat(t,
		NodeUtil_MayBeStringWithUseType(
			rhino.IR_Name("a").SetJSType(registry.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)), true)).IsFalse()
	nuAssertThat(t,
		NodeUtil_MayBeStringWithUseType(
			rhino.IR_Name("a").SetJSType(registry.GetNativeType(rhino.JSTypeNative_BIGINT_TYPE)), true)).IsFalse()
	nuAssertThat(t,
		NodeUtil_MayBeStringWithUseType(
			rhino.IR_Name("a").SetJSType(registry.GetNativeType(rhino.JSTypeNative_BOOLEAN_TYPE)), true)).IsFalse()
	nuAssertThat(t,
		NodeUtil_MayBeStringWithUseType(
			rhino.IR_Name("a").SetJSType(registry.GetNativeType(rhino.JSTypeNative_STRING_TYPE)), true)).IsTrue()
}

func TestNodeUtil_AssortedTests_TestIsStringResult(t *testing.T) {
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "1"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "true"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "+true"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "+1"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "-1"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "-Infinity"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "Infinity"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "NaN"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "undefined"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "void 0"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "null"))).IsFalse()

	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "a << b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "a >> b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "a >>> b"))).IsFalse()

	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "a == b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "a != b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "a === b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "a !== b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "a < b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "a > b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "a <= b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "a >= b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "a in b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "a instanceof b"))).IsFalse()

	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "'a'"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "'a'+b"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "a+'b'"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "a+b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "a()"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "''.a"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "a.b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "a.b()"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "a().b()"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "new a()"))).IsFalse()

	// These can't be strings but they aren't handled yet.
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "1 && 2"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "1 || 2"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "1 ? 2 : 3"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "1,2"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "a=1"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "1+1"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "true+true"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "null+null"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "NaN+NaN"))).IsFalse()

	// These are not strings but they aren't primitives either
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "([1,2])"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "({a:1})"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "({}+1)"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "(1+{})"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "([]+1)"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "(1+[])"))).IsFalse()

	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "a += 'x'"))).IsTrue()

	// Template literals
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "`x`"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsStringResult(nuParseExpr(t, "`a${b}c`"))).IsTrue()
}

func TestNodeUtil_AssortedTests_TestMayBeObject(t *testing.T) {
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "1"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "1n"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "a * b"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "true"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "+true"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "+1"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "-1"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "-Infinity"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "Infinity"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "NaN"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "undefined"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "void 0"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "null"))).IsFalse()

	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "a << b"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "a >> b"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "a >>> b"))).IsFalse()

	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "a == b"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "a != b"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "a === b"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "a !== b"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "a < b"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "a > b"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "a <= b"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "a >= b"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "a in b"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "a instanceof b"))).IsFalse()

	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "'a'"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "'a'+b"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "a+'b'"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "a+b"))).IsTrue()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "a()"))).IsTrue()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "''.a"))).IsTrue()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "a.b"))).IsTrue()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "a.b()"))).IsTrue()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "a().b()"))).IsTrue()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "new a()"))).IsTrue()

	// These can't be objects but they aren't handled yet.
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "1 && 2"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "1 || 2"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "1 ? 2 : 3"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "1,2"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "a=1"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "1+1"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "true+true"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "null+null"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "NaN+NaN"))).IsFalse()

	// These are not objects but they aren't primitives either
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "([1,2])"))).IsTrue()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "({a:1})"))).IsTrue()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "({}+1)"))).IsTrue()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "(1+{})"))).IsTrue()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "([]+1)"))).IsTrue()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "(1+[])"))).IsTrue()

	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "a += 'x'"))).IsFalse()
	nuAssertThat(t, NodeUtil_MayBeObject(nuParseExpr(t, "a += 1"))).IsTrue()
}

func TestNodeUtil_AssortedTests_TestIsObjectResult(t *testing.T) {
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "1"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "true"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "+true"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "+1"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "-1"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "-Infinity"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "Infinity"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "NaN"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "undefined"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "void 0"))).IsFalse()

	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "a << b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "a >> b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "a >>> b"))).IsFalse()

	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "a == b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "a != b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "a === b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "a !== b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "a < b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "a > b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "a <= b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "a >= b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "a in b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "a instanceof b"))).IsFalse()

	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "delete a"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "'a'"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "'a'+b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "a+'b'"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "a+b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "{},true"))).IsFalse()

	// "false" here means "unknown"
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "a()"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "''.a"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "a.b"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "a.b()"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "a().b()"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "a ? true : {}"))).IsFalse()

	// These are objects but aren't handled yet.
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "true && {}"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "true || {}"))).IsFalse()

	// Definitely objects
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "new a.b()"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "([true,false])"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "({a:true})"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "a={}"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "[] && {}"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "[] || {}"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "a ? [] : {}"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "{},[]"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsObjectResult(nuParseExpr(t, "/a/g"))).IsTrue()
}

func TestNodeUtil_AssortedTests_TestValidNames(t *testing.T) {
	nuAssertThat(t, nuh_isValidPropertyName(t, "a")).IsTrue()
	nuAssertThat(t, nuh_isValidPropertyName(t, "a3")).IsTrue()
	nuAssertThat(t, nuh_isValidPropertyName(t, "3a")).IsFalse()
	nuAssertThat(t, nuh_isValidPropertyName(t, "a.")).IsFalse()
	nuAssertThat(t, nuh_isValidPropertyName(t, ".a")).IsFalse()
	nuAssertThat(t, nuh_isValidPropertyName(t, "a.b")).IsFalse()
	nuAssertThat(t, nuh_isValidPropertyName(t, "true")).IsFalse()
	nuAssertThat(t, nuh_isValidPropertyName(t, "a.true")).IsFalse()
	nuAssertThat(t, nuh_isValidPropertyName(t, "a..b")).IsFalse()

	nuAssertThat(t, NodeUtil_IsValidSimpleName("a")).IsTrue()
	nuAssertThat(t, NodeUtil_IsValidSimpleName("a3")).IsTrue()
	nuAssertThat(t, NodeUtil_IsValidSimpleName("3a")).IsFalse()
	nuAssertThat(t, NodeUtil_IsValidSimpleName("a.")).IsFalse()
	nuAssertThat(t, NodeUtil_IsValidSimpleName(".a")).IsFalse()
	nuAssertThat(t, NodeUtil_IsValidSimpleName("a.b")).IsFalse()
	nuAssertThat(t, NodeUtil_IsValidSimpleName("true")).IsFalse()
	nuAssertThat(t, NodeUtil_IsValidSimpleName("a.true")).IsFalse()
	nuAssertThat(t, NodeUtil_IsValidSimpleName("a..b")).IsFalse()

	nuAssertThat(t, nuh_isValidQualifiedName(t, "a")).IsTrue()
	nuAssertThat(t, nuh_isValidQualifiedName(t, "a3")).IsTrue()
	nuAssertThat(t, nuh_isValidQualifiedName(t, "3a")).IsFalse()
	nuAssertThat(t, nuh_isValidQualifiedName(t, "a.")).IsFalse()
	nuAssertThat(t, nuh_isValidQualifiedName(t, ".a")).IsFalse()
	nuAssertThat(t, nuh_isValidQualifiedName(t, "a.b")).IsTrue()
	nuAssertThat(t, nuh_isValidQualifiedName(t, "true")).IsFalse()
	nuAssertThat(t, nuh_isValidQualifiedName(t, "a.true")).IsFalse()
	nuAssertThat(t, nuh_isValidQualifiedName(t, "a..b")).IsFalse()
}

func TestNodeUtil_AssortedTests_TestGetNearestFunctionName(t *testing.T) {
	nuh_testFunctionName(t, "(function() {})()", nil)
	nuh_testFunctionName(t, "function a() {}", "a")
	nuh_testFunctionName(t, "(function a() {})", "a")
	nuh_testFunctionName(t, "({a:function () {}})", "a")
	nuh_testFunctionName(t, "({get a() {}})", "a")
	nuh_testFunctionName(t, "({set a(b) {}})", "a")
	nuh_testFunctionName(t, "({set a(b) {}})", "a")
	nuh_testFunctionName(t, "({1:function () {}})", "1")
	nuh_testFunctionName(t, "var a = function a() {}", "a")
	nuh_testFunctionName(t, "var a;a = function a() {}", "a")
	nuh_testFunctionName(t, "var o;o.a = function a() {}", "o.a")
	nuh_testFunctionName(t, "this.a = function a() {}", "this.a")
}

func TestNodeUtil_AssortedTests_TestGetBestLValue(t *testing.T) {
	nuAssertThat(t, nuh_getFunctionLValue(t, "var x = function() {};")).IsEqualTo("x")
	nuAssertThat(t, nuh_getFunctionLValue(t, "x = function() {};")).IsEqualTo("x")
	nuAssertThat(t, nuh_getFunctionLValue(t, "function x() {};")).IsEqualTo("x")
	nuAssertThat(t, nuh_getFunctionLValue(t, "var x = y ? z : function() {};")).IsEqualTo("x")
	nuAssertThat(t, nuh_getFunctionLValue(t, "var x = y ? function() {} : z;")).IsEqualTo("x")
	nuAssertThat(t, nuh_getFunctionLValue(t, "var x = y && function() {};")).IsEqualTo("x")
	nuAssertThat(t, nuh_getFunctionLValue(t, "var x = y || function() {};")).IsEqualTo("x")
	nuAssertThat(t, nuh_getFunctionLValue(t, "var x = (y, function() {});")).IsEqualTo("x")
}

func TestNodeUtil_AssortedTests_TestGetBestLValueName(t *testing.T) {
	getBestName := func(js string) *string {
		return NodeUtil_GetBestLValueName(NodeUtil_GetBestLValue(nuParseFirst(t, rhino.Token_FUNCTION, js)))
	}
	nuAssertThat(t, getBestName("var x = function() {};")).IsEqualTo("x")
	nuAssertThat(t, getBestName("x = function() {};")).IsEqualTo("x")
	nuAssertThat(t, getBestName("function x() {};")).IsEqualTo("x")
	nuAssertThat(t, getBestName("var x = y ? z : function() {};")).IsEqualTo("x")
	nuAssertThat(t, getBestName("var x = y ? function() {} : z;")).IsEqualTo("x")
	nuAssertThat(t, getBestName("var x = y && function() {};")).IsEqualTo("x")
	nuAssertThat(t, getBestName("var x = y || function() {};")).IsEqualTo("x")
	nuAssertThat(t, getBestName("var x = (y, function() {});")).IsEqualTo("x")
	nuAssertThat(t, getBestName("C.prototype.d = function() {};")).IsEqualTo("C.prototype.d")
	nuAssertThat(t, getBestName("class C { d() {} };")).IsEqualTo("C.prototype.d")
	nuAssertThat(t, getBestName("C.d = function() {};")).IsEqualTo("C.d")
	nuAssertThat(t, getBestName("class C { static d() {} };")).IsEqualTo("C.d")
}

func TestNodeUtil_AssortedTests_TestGetRValueOfLValue(t *testing.T) {
	nuAssertThat(t, nuh_functionIsRValueOfAssign(t, "x = function() {};")).IsTrue()
	nuAssertThat(t, nuh_functionIsRValueOfAssign(t, "x += function() {};")).IsTrue()
	nuAssertThat(t, nuh_functionIsRValueOfAssign(t, "x -= function() {};")).IsTrue()
	nuAssertThat(t, nuh_functionIsRValueOfAssign(t, "x *= function() {};")).IsTrue()
	nuAssertThat(t, nuh_functionIsRValueOfAssign(t, "x /= function() {};")).IsTrue()
	nuAssertThat(t, nuh_functionIsRValueOfAssign(t, "x <<= function() {};")).IsTrue()
	nuAssertThat(t, nuh_functionIsRValueOfAssign(t, "x >>= function() {};")).IsTrue()
	nuAssertThat(t, nuh_functionIsRValueOfAssign(t, "x >>= function() {};")).IsTrue()
	nuAssertThat(t, nuh_functionIsRValueOfAssign(t, "x >>>= function() {};")).IsTrue()
	nuAssertThat(t, nuh_functionIsRValueOfAssign(t, "x = y ? x : function() {};")).IsFalse()
}

func TestNodeUtil_AssortedTests_TestGetRValueOfLValueInObjectLiteral(t *testing.T) {
	scriptNode := nuParse(t, "var x = {[computedProp] : 5};")
	computedProp :=
		scriptNode.GetFirstChild().GetFirstChild().GetFirstChild().GetFirstChild() // [computedProp] : 5 // {...} // x = {...} // VAR
	jutil.CheckStateMsg1(computedProp.IsComputedProp(), "%s", computedProp)
	nuAssertNode(t, NodeUtil_GetRValueOfLValue(computedProp)).HasType(rhino.Token_NUMBER)
}

func TestNodeUtil_AssortedTests_TestGetRValueOfLValueFromClass(t *testing.T) {
	scriptNode := nuParse(t, "class Foo { ['computedField'] = 5}")
	computedField :=
		scriptNode.GetFirstChild().GetLastChild().GetFirstChild() // COMPUTED_FIELD_DEF  ['computedField'] = 5 // CLASS_MEMBERS  { ['computedField'] = 5} // CLASS
	jutil.CheckStateMsg1(computedField.IsComputedFieldDef(), "%s", computedField)
	nuAssertNode(t, NodeUtil_GetRValueOfLValue(computedField)).HasType(rhino.Token_NUMBER)
}

/**
 * When the left side is a destructuring pattern, generally it's not possible to identify the
 * RHS for a specific name on the LHS.
 */
func TestNodeUtil_AssortedTests_TestGetRValueOfLValueInDestructuringPattern(t *testing.T) {
	nuAssertThat(t, NodeUtil_GetRValueOfLValue(nuh_getNameNode(t, nuParse(t, "var [x] = rhs;"), "x"))).IsNull()
	nuAssertThat(t, NodeUtil_GetRValueOfLValue(nuh_getNameNode(t, nuParse(t, "var [x, y] = rhs;"), "x"))).IsNull()
	nuAssertThat(t, NodeUtil_GetRValueOfLValue(nuh_getNameNode(t, nuParse(t, "var [y, x] = rhs;"), "x"))).IsNull()
	nuAssertThat(t, NodeUtil_GetRValueOfLValue(nuh_getNameNode(t, nuParse(t, "var {x: x} = rhs;"), "x"))).IsNull()
	nuAssertThat(t, NodeUtil_GetRValueOfLValue(nuh_getNameNode(t, nuParse(t, "var {y: x} = rhs;"), "x"))).IsNull()

	ast := nuParse(t, "var {x} = rhs;")
	x :=
		ast.GetFirstChild().GetFirstChild().GetFirstChild().GetFirstChild() // STRING_KEY // OBJECT_PATTERN // DESTRUCTURING_LHS // VAR
	jutil.CheckStateMsg1(x.IsStringKey(), "%s", x)
	nuAssertThat(t, NodeUtil_GetRValueOfLValue(x)).IsNull()
}

func TestNodeUtil_AssortedTests_TestGetRValueOfLValueDestructuringPattern(t *testing.T) {
	nuAssertNode(t, NodeUtil_GetRValueOfLValue(nuh_getPattern(t, nuParse(t, "var [x] = 'rhs';")))).HasType(rhino.Token_STRINGLIT)
	nuAssertNode(t, NodeUtil_GetRValueOfLValue(nuh_getPattern(t, nuParse(t, "var [x, y] = 'rhs';")))).HasType(rhino.Token_STRINGLIT)
	nuAssertNode(t, NodeUtil_GetRValueOfLValue(nuh_getPattern(t, nuParse(t, "var [y, x] = 'rhs';")))).HasType(rhino.Token_STRINGLIT)
	nuAssertNode(t, NodeUtil_GetRValueOfLValue(nuh_getPattern(t, nuParse(t, "var {x: x} = 'rhs';")))).HasType(rhino.Token_STRINGLIT)
	nuAssertNode(t, NodeUtil_GetRValueOfLValue(nuh_getPattern(t, nuParse(t, "var {y: x} = 'rhs';")))).HasType(rhino.Token_STRINGLIT)
	nuAssertNode(t, NodeUtil_GetRValueOfLValue(nuh_getPattern(t, nuParse(t, "var {x} = 'rhs';")))).HasType(rhino.Token_STRINGLIT)
}

func TestNodeUtil_AssortedTests_TestGetRValueOfLValueDestructuringLhs(t *testing.T) {
	nuAssertNode(t, NodeUtil_GetRValueOfLValue(nuParseFirst(t, rhino.Token_DESTRUCTURING_LHS, "var [x] = 'rhs';"))).HasType(rhino.Token_STRINGLIT)
	nuAssertNode(t, NodeUtil_GetRValueOfLValue(nuParseFirst(t, rhino.Token_DESTRUCTURING_LHS, "var [x, y] = 'rhs';"))).HasType(rhino.Token_STRINGLIT)
	nuAssertNode(t, NodeUtil_GetRValueOfLValue(nuParseFirst(t, rhino.Token_DESTRUCTURING_LHS, "var [y, x] = 'rhs';"))).HasType(rhino.Token_STRINGLIT)
	nuAssertNode(t, NodeUtil_GetRValueOfLValue(nuParseFirst(t, rhino.Token_DESTRUCTURING_LHS, "var {x: x} = 'rhs';"))).HasType(rhino.Token_STRINGLIT)
	nuAssertNode(t, NodeUtil_GetRValueOfLValue(nuParseFirst(t, rhino.Token_DESTRUCTURING_LHS, "var {y: x} = 'rhs';"))).HasType(rhino.Token_STRINGLIT)
	nuAssertNode(t, NodeUtil_GetRValueOfLValue(nuParseFirst(t, rhino.Token_DESTRUCTURING_LHS, "var {x} = 'rhs';"))).HasType(rhino.Token_STRINGLIT)
}

func TestNodeUtil_AssortedTests_TestIsNaN(t *testing.T) {
	nuAssertThat(t, NodeUtil_IsNaN(nuParseExpr(t, "NaN"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsNaN(nuParseExpr(t, "Infinity"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNaN(nuParseExpr(t, "x"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNaN(nuParseExpr(t, "0/0"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsNaN(nuParseExpr(t, "1/0"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNaN(nuParseExpr(t, "0/1"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsNaN(rhino.IR_Number(0.0))).IsFalse()
}

func TestNodeUtil_AssortedTests_TestIsExecutedExactlyOnce(t *testing.T) {
	nuAssertThat(t, nuh_executedOnceTestCase(t, "x;")).IsTrue()

	nuAssertThat(t, nuh_executedOnceTestCase(t, "x && 1;")).IsTrue()
	nuAssertThat(t, nuh_executedOnceTestCase(t, "1 && x;")).IsFalse()

	nuAssertThat(t, nuh_executedOnceTestCase(t, "1 && (x && 1);")).IsFalse()

	nuAssertThat(t, nuh_executedOnceTestCase(t, "x || 1;")).IsTrue()
	nuAssertThat(t, nuh_executedOnceTestCase(t, "1 || x;")).IsFalse()

	nuAssertThat(t, nuh_executedOnceTestCase(t, "1 && (x || 1);")).IsFalse()

	nuAssertThat(t, nuh_executedOnceTestCase(t, "x ||= 1;")).IsTrue()
	nuAssertThat(t, nuh_executedOnceTestCase(t, "y ||= x;")).IsFalse()
	nuAssertThat(t, nuh_executedOnceTestCase(t, "x &&= 1;")).IsTrue()
	nuAssertThat(t, nuh_executedOnceTestCase(t, "y &&= x;")).IsFalse()
	nuAssertThat(t, nuh_executedOnceTestCase(t, "x ??= 1;")).IsTrue()
	nuAssertThat(t, nuh_executedOnceTestCase(t, "y ??= x;")).IsFalse()
	nuAssertThat(t, nuh_executedOnceTestCase(t, "x ? 1 : 2;")).IsTrue()
	nuAssertThat(t, nuh_executedOnceTestCase(t, "1 ? 1 : x;")).IsFalse()
	nuAssertThat(t, nuh_executedOnceTestCase(t, "1 ? x : 2;")).IsFalse()

	nuAssertThat(t, nuh_executedOnceTestCase(t, "1 && (x ? 1 : 2);")).IsFalse()

	nuAssertThat(t, nuh_executedOnceTestCase(t, "if (x) {}")).IsTrue()
	nuAssertThat(t, nuh_executedOnceTestCase(t, "if (true) {x;}")).IsFalse()
	nuAssertThat(t, nuh_executedOnceTestCase(t, "if (true) {} else {x;}")).IsFalse()

	nuAssertThat(t, nuh_executedOnceTestCase(t, "if (1) { if (x) {} }")).IsFalse()

	nuAssertThat(t, nuh_executedOnceTestCase(t, "for(x;;){}")).IsTrue()
	nuAssertThat(t, nuh_executedOnceTestCase(t, "for(;x;){}")).IsFalse()
	nuAssertThat(t, nuh_executedOnceTestCase(t, "for(;;x){}")).IsFalse()
	nuAssertThat(t, nuh_executedOnceTestCase(t, "for(;;){x;}")).IsFalse()

	nuAssertThat(t, nuh_executedOnceTestCase(t, "if (1) { for(x;;){} }")).IsFalse()

	nuAssertThat(t, nuh_executedOnceTestCase(t, "for(x in {}){}")).IsFalse()
	nuAssertThat(t, nuh_executedOnceTestCase(t, "for({}.a in x){}")).IsTrue()
	nuAssertThat(t, nuh_executedOnceTestCase(t, "for({}.a in {}){x}")).IsFalse()

	nuAssertThat(t, nuh_executedOnceTestCase(t, "if (1) { for(x in {}){} }")).IsFalse()

	nuAssertThat(t, nuh_executedOnceTestCase(t, "for(x of {}){}")).IsFalse()
	nuAssertThat(t, nuh_executedOnceTestCase(t, "for({}.a of x){}")).IsTrue()
	nuAssertThat(t, nuh_executedOnceTestCase(t, "for({}.a of {}){x}")).IsFalse()

	nuAssertThat(t, nuh_executedOnceTestCase(t, "if (1) { for(x of {}){} }")).IsFalse()

	nuAssertThat(t, nuh_executedOnceTestCase(t, "async function f() { for await(x of {}){} }")).IsFalse()
	nuAssertThat(t, nuh_executedOnceTestCase(t, "async function f() { for await({}.a of x){} }")).IsTrue()
	nuAssertThat(t, nuh_executedOnceTestCase(t, "async function f() { for await({}.a of {}){x} }")).IsFalse()

	nuAssertThat(t, nuh_executedOnceTestCase(t, "async function f() { if (1) { for await(x of {}){} } }")).
		IsFalse()
	nuAssertThat(t, nuh_executedOnceTestCase(t, "switch (x) {}")).IsTrue()
	nuAssertThat(t, nuh_executedOnceTestCase(t, "switch (1) {case x:}")).IsFalse()
	nuAssertThat(t, nuh_executedOnceTestCase(t, "switch (1) {case 1: x}")).IsFalse()
	nuAssertThat(t, nuh_executedOnceTestCase(t, "switch (1) {default: x}")).IsFalse()

	nuAssertThat(t, nuh_executedOnceTestCase(t, "if (1) { switch (x) {} }")).IsFalse()

	nuAssertThat(t, nuh_executedOnceTestCase(t, "while (x) {}")).IsFalse()
	nuAssertThat(t, nuh_executedOnceTestCase(t, "while (1) {x}")).IsFalse()

	nuAssertThat(t, nuh_executedOnceTestCase(t, "do {} while (x)")).IsFalse()
	nuAssertThat(t, nuh_executedOnceTestCase(t, "do {x} while (1)")).IsFalse()

	nuAssertThat(t, nuh_executedOnceTestCase(t, "try {x} catch (e) {}")).IsFalse()
	nuAssertThat(t, nuh_executedOnceTestCase(t, "try {} catch (e) {x}")).IsFalse()
	nuAssertThat(t, nuh_executedOnceTestCase(t, "try {} finally {x}")).IsTrue()

	nuAssertThat(t, nuh_executedOnceTestCase(t, "if (1) { try {} finally {x} }")).IsFalse()
}

func nuh_assertLValueNamedX(t *testing.T, n rhino.Node) {
	nuAssertThat(t, n.GetString()).IsEqualTo("x")
	nuAssertThat(t, NodeUtil_IsLValue(n)).IsTrue()
}

func TestNodeUtil_AssortedTests_TestIsLValue(t *testing.T) {
	nuh_assertLValueNamedX(t, nuParse(t, "var x;").GetFirstFirstChild())
	nuh_assertLValueNamedX(t, nuParse(t, "var w, x;").GetFirstChild().GetLastChild())
	nuh_assertLValueNamedX(t,
		nuParse(t, "var [...x] = y;").GetFirstFirstChild().GetFirstFirstChild().GetFirstChild())
	nuh_assertLValueNamedX(t, nuParse(t, "var x = y;").GetFirstFirstChild())
	nuh_assertLValueNamedX(t, nuParse(t, "x++;").GetFirstFirstChild().GetFirstChild())
	nuh_assertLValueNamedX(t,
		NodeUtil_GetFunctionParameters(nuParse(t, "function f(x) {}").GetFirstChild()).GetFirstChild())

	x :=
		NodeUtil_GetFunctionParameters(nuParse(t, "function f(x = 3) {}").GetFirstChild()).GetFirstChild().GetFirstChild() // x // x = 3
	nuh_assertLValueNamedX(t, x)

	nuh_assertLValueNamedX(t,
		nuParse(t, "({x} = obj)").GetFirstFirstChild().GetFirstFirstChild().GetFirstChild())
	nuh_assertLValueNamedX(t, nuParse(t, "([x] = obj)").GetFirstFirstChild().GetFirstFirstChild())
	nuh_assertLValueNamedX(t,
		nuParse(t, "function foo (...x) {}").GetFirstChild().GetSecondChild().GetFirstFirstChild())
	nuh_assertLValueNamedX(t,
		nuParse(t, "({[0]: x} = obj)").GetFirstFirstChild().GetFirstFirstChild().GetSecondChild())
}

func nuh_assertNotLValueNamedX(t *testing.T, n rhino.Node) {
	nuAssertThat(t, n.GetString()).IsEqualTo("x")
	nuAssertThat(t, NodeUtil_IsLValue(n)).IsFalse()
}

func TestNodeUtil_AssortedTests_TestIsNotLValue(t *testing.T) {
	nuh_assertNotLValueNamedX(t, nuParse(t, "var a = x;").GetFirstFirstChild().GetFirstChild())

	x :=
		nuParse(t, "f(...x);").GetFirstChild().GetFirstChild().GetLastChild().GetFirstChild() // x // spread // call // expr result // script
	nuh_assertNotLValueNamedX(t, x)

	x =
		nuParse(t, "var a = [...x];").GetFirstChild().GetFirstChild().GetFirstChild().GetFirstChild().GetFirstChild() // x // spread // array // a // var // script
	nuh_assertNotLValueNamedX(t, x)
}

func nuh_assertDeclarationLValueNamedX(t *testing.T, js string) {
	x := nuh_getNameNode(t, nuParse(t, js), "x")
	nuAssertThat(t, NodeUtil_IsDeclarationLValue(x)).IsTrue()
}

func nuh_assertNotDeclarationLValueNamedX(t *testing.T, js string) {
	x := nuh_getNameNode(t, nuParse(t, js), "x")
	nuAssertThat(t, NodeUtil_IsDeclarationLValue(x)).IsFalse()
}

func TestNodeUtil_AssortedTests_TestIsDeclarationLValue(t *testing.T) {
	nuh_assertDeclarationLValueNamedX(t, "var x;")
	nuh_assertDeclarationLValueNamedX(t, "let x = 1;")
	nuh_assertDeclarationLValueNamedX(t, "const x = 1;")
	nuh_assertDeclarationLValueNamedX(t, "const [x] = arr;")
	nuh_assertDeclarationLValueNamedX(t, "const [[x]] = arr;")
	nuh_assertDeclarationLValueNamedX(t, "const {y: x} = obj;")
	nuh_assertDeclarationLValueNamedX(t, "const {z: {y: x}} = obj;")
	nuh_assertDeclarationLValueNamedX(t, "const [x = 1] = arr;")
	nuh_assertDeclarationLValueNamedX(t, "const [...x] = arr;")
	nuh_assertDeclarationLValueNamedX(t, "const {...x} = obj;")
	nuh_assertDeclarationLValueNamedX(t, "const {[k]: x} = obj;")
	nuh_assertDeclarationLValueNamedX(t, "const {[k]: [x = 1]} = obj;")
	nuh_assertDeclarationLValueNamedX(t, "for (const [x] of arr) {}")
	nuh_assertDeclarationLValueNamedX(t, "for (const [[x]] of arr) {}")
	nuh_assertDeclarationLValueNamedX(t, "function f(x) {}")
	nuh_assertDeclarationLValueNamedX(t, "function f(x = 1) {}")
	nuh_assertDeclarationLValueNamedX(t, "function f(...x) {}")
	nuh_assertDeclarationLValueNamedX(t, "function f({z: [x = 1]} = {}) {}")
	nuh_assertDeclarationLValueNamedX(t, "try {} catch ({y: [x]}) {}")
}

func TestNodeUtil_AssortedTests_TestIsNotDeclarationLValue(t *testing.T) {
	nuh_assertNotDeclarationLValueNamedX(t, "x = 1;")
	nuh_assertNotDeclarationLValueNamedX(t, "x++;")
	nuh_assertNotDeclarationLValueNamedX(t, "([x] = arr);")
	nuh_assertNotDeclarationLValueNamedX(t, "([[x]] = arr);")
	nuh_assertNotDeclarationLValueNamedX(t, "({y: x} = obj);")
	nuh_assertNotDeclarationLValueNamedX(t, "({z: {y: x}} = obj);")
	nuh_assertNotDeclarationLValueNamedX(t, "([x = 1] = arr);")
	nuh_assertNotDeclarationLValueNamedX(t, "([...x] = arr);")
	nuh_assertNotDeclarationLValueNamedX(t, "({...x} = obj);")
	nuh_assertNotDeclarationLValueNamedX(t, "({[k]: x} = obj);")
	nuh_assertNotDeclarationLValueNamedX(t, "for ([x] of arr) {}")
	nuh_assertNotDeclarationLValueNamedX(t, "for ([[x]] of arr) {}")
	nuh_assertNotDeclarationLValueNamedX(t, "var a = x;")
	nuh_assertNotDeclarationLValueNamedX(t, "const [a = x] = arr;")
	nuh_assertNotDeclarationLValueNamedX(t, "const {[x]: a} = obj;")
}

func TestNodeUtil_AssortedTests_TestGetBestJsDocInfoNodeStrict_declaredName_doesNotThrow(t *testing.T) {
	aName := nuParse(t, "/** some */ let A;").GetFirstFirstChild()
	nuAssertThat(t, aName.IsName()).IsTrue()

	bestJsDocInfoNode := NodeUtil_GetBestJSDocInfoNode(aName)
	nuAssertThat(t, bestJsDocInfoNode).IsNotNull()
	nuAssertThat(t, bestJsDocInfoNode == aName.GetParent()).IsTrue()

	// check that the strict version returns the same node
	bestJsDocInfoNodeStrict := NodeUtil_GetBestJsDocInfoNodeStrict(aName)
	nuAssertThat(t, bestJsDocInfoNodeStrict).IsNotNull()
	nuAssertThat(t, bestJsDocInfoNodeStrict == aName.GetParent()).IsTrue()
}

func TestNodeUtil_AssortedTests_TestGetBestJsDocInfoNodeStrict_classDeclaration_doesNotThrow(t *testing.T) {
	classNode := nuParse(t, "/** some */ class A {}").GetFirstChild()
	nuAssertThat(t, classNode.IsClass()).IsTrue()

	bestJsDocInfoNode := NodeUtil_GetBestJSDocInfoNode(classNode)
	nuAssertThat(t, bestJsDocInfoNode).IsNotNull()
	nuAssertThat(t, bestJsDocInfoNode == classNode).IsTrue()

	// check that the strict version returns the same node
	bestJsDocInfoNodeStrict := NodeUtil_GetBestJsDocInfoNodeStrict(classNode)
	nuAssertThat(t, bestJsDocInfoNodeStrict).IsNotNull()
	nuAssertThat(t, bestJsDocInfoNodeStrict == classNode).IsTrue()
}

func TestNodeUtil_AssortedTests_TestGetBestJsDocInfoNodeStrict_rhsNamedClassExpression_doesNotThrow(t *testing.T) {
	constNode := nuParse(t, "/** some */ const x = class A {}").GetFirstChild()
	xName := constNode.GetFirstChild()
	classNode := xName.GetFirstChild()
	nuAssertThat(t, classNode.IsClass()).IsTrue()

	bestJsDocInfoNode := NodeUtil_GetBestJSDocInfoNode(classNode)
	nuAssertThat(t, bestJsDocInfoNode).IsNotNull()
	nuAssertThat(t, bestJsDocInfoNode == constNode).IsTrue()

	// check that the strict version returns the same node
	bestJsDocInfoNodeStrict := NodeUtil_GetBestJsDocInfoNodeStrict(classNode)
	nuAssertThat(t, bestJsDocInfoNodeStrict).IsNotNull()
	nuAssertThat(t, bestJsDocInfoNodeStrict == constNode).IsTrue()
}

func TestNodeUtil_AssortedTests_TestGetBestJsDocInfoNodeStrict_rhsUnnamedClassExpression_throws(t *testing.T) {
	constNode := nuParse(t, "/** some */ const x = class {}").GetFirstChild()
	xName := constNode.GetFirstChild()
	classNode := xName.GetFirstChild()
	nuAssertThat(t, classNode.IsClass()).IsTrue()

	bestJsDocInfoNode := NodeUtil_GetBestJSDocInfoNode(classNode)
	nuAssertThat(t, bestJsDocInfoNode).IsNotNull()
	nuAssertThat(t, bestJsDocInfoNode == constNode).IsTrue()

	// check that the strict version rejects an unnamed class expression
	e := assertThrows[jutil.IIllegalStateException](t, func() { NodeUtil_GetBestJsDocInfoNodeStrict(classNode) })
	nuAssertThat(t, strings.Contains(e.GetMessage(), "Not allowed to get JSDocInfo node for node")).IsTrue()
}

func TestNodeUtil_AssortedTests_TestGetBestJsDocInfoNodeStrict_rhsClassName_doesNotThrow(t *testing.T) {
	constNode := nuParse(t, "/** some */ const x = class A {}").GetFirstChild()
	xName := constNode.GetFirstChild()
	classNode := xName.GetFirstChild()
	nuAssertThat(t, classNode.IsClass()).IsTrue()
	className := classNode.GetFirstChild()
	nuAssertThat(t, className.IsName()).IsTrue()

	bestJsDocInfoNode := NodeUtil_GetBestJSDocInfoNode(className)
	nuAssertThat(t, bestJsDocInfoNode).IsNotNull()
	nuAssertThat(t, bestJsDocInfoNode == constNode).IsTrue()

	// check that the strict version returns the same node
	bestJsDocInfoNodeStrict := NodeUtil_GetBestJsDocInfoNodeStrict(className)
	nuAssertThat(t, bestJsDocInfoNodeStrict).IsNotNull()
	nuAssertThat(t, bestJsDocInfoNodeStrict == constNode).IsTrue()
}

func TestNodeUtil_AssortedTests_TestGetBestJsDocInfoNodeStrict_rhsNamedFunctionExpression_doesNotThrow(t *testing.T) {
	constNode := nuParse(t, "/** some */ const x = function A() {}").GetFirstChild()
	xName := constNode.GetFirstChild()
	functionNode := xName.GetFirstChild()
	nuAssertThat(t, functionNode.IsFunction()).IsTrue()

	bestJsDocInfoNode := NodeUtil_GetBestJSDocInfoNode(functionNode)
	nuAssertThat(t, bestJsDocInfoNode).IsNotNull()
	nuAssertThat(t, bestJsDocInfoNode == constNode).IsTrue()

	// check that the strict version returns the same node
	bestJsDocInfoNodeStrict := NodeUtil_GetBestJsDocInfoNodeStrict(functionNode)
	nuAssertThat(t, bestJsDocInfoNodeStrict).IsNotNull()
	nuAssertThat(t, bestJsDocInfoNodeStrict == constNode).IsTrue()
}

func TestNodeUtil_AssortedTests_TestGetBestJsDocInfoNodeStrict_exportedFunction_doesNotThrow(t *testing.T) {
	exportNode := nuParse(t, "/** some */ export function A() {}").GetFirstFirstChild()
	nuAssertThat(t, exportNode.IsExport()).IsTrue()
	functionNode := exportNode.GetFirstChild()
	nuAssertThat(t, functionNode.IsFunction()).IsTrue()

	bestJsDocInfoNode := NodeUtil_GetBestJSDocInfoNode(functionNode)
	nuAssertThat(t, bestJsDocInfoNode).IsNotNull()
	// JSDoc gets attached to the children of export nodes, and there are no warnings.
	// See https://github.com/google/closure-compiler/issues/781
	// Hence, the best JSDocInfo node is the function node itself when the search is starting from
	// it.
	nuAssertThat(t, bestJsDocInfoNode == functionNode).IsTrue()

	// check that the strict version returns the same node
	bestJsDocInfoNodeStrict := NodeUtil_GetBestJsDocInfoNodeStrict(functionNode)
	nuAssertThat(t, bestJsDocInfoNodeStrict).IsNotNull()
	nuAssertThat(t, bestJsDocInfoNodeStrict == functionNode).IsTrue()
}

func TestNodeUtil_AssortedTests_TestGetBestJsDocInfoNodeStrict_exportNode_doesNotThrow(t *testing.T) {
	exportNode := nuParse(t, "/** some */ export function A() {}").GetFirstFirstChild()
	nuAssertThat(t, exportNode.IsExport()).IsTrue()
	functionNode := exportNode.GetFirstChild()
	nuAssertThat(t, functionNode.IsFunction()).IsTrue()

	bestJsDocInfoNode := NodeUtil_GetBestJSDocInfoNode(exportNode)
	nuAssertThat(t, bestJsDocInfoNode).IsNotNull()
	nuAssertThat(t, bestJsDocInfoNode == exportNode).IsTrue()

	// check that the strict version returns the same node
	bestJsDocInfoNodeStrict := NodeUtil_GetBestJsDocInfoNodeStrict(exportNode)
	nuAssertThat(t, bestJsDocInfoNodeStrict).IsNotNull()
	nuAssertThat(t, bestJsDocInfoNodeStrict == exportNode).IsTrue()
}

func TestNodeUtil_AssortedTests_TestGetBestJsDocInfoNodeStrict_rhsFunctionName_doesNotThrow(t *testing.T) {
	constNode := nuParse(t, "/** some */ const x = function A() {}").GetFirstChild()
	xName := constNode.GetFirstChild()
	functionNode := xName.GetFirstChild()
	nuAssertThat(t, functionNode.IsFunction()).IsTrue()
	functionName := functionNode.GetFirstChild()
	nuAssertThat(t, functionName.IsName()).IsTrue()

	bestJsDocInfoNode := NodeUtil_GetBestJSDocInfoNode(functionName)
	nuAssertThat(t, bestJsDocInfoNode).IsNotNull()
	nuAssertThat(t, bestJsDocInfoNode == constNode).IsTrue()

	// check that the strict version returns the same node
	bestJsDocInfoNodeStrict := NodeUtil_GetBestJsDocInfoNodeStrict(functionName)
	nuAssertThat(t, bestJsDocInfoNodeStrict).IsNotNull()
	nuAssertThat(t, bestJsDocInfoNodeStrict == constNode).IsTrue()
}

func TestNodeUtil_AssortedTests_TestIsConstantDeclaration(t *testing.T) {
	nuh_assertIsConstantDeclaration(t, false, nuParse(t, "var x = 1;").GetFirstFirstChild())
	nuh_assertIsConstantDeclaration(t, false, nuParse(t, "let x = 1;").GetFirstFirstChild())
	nuh_assertIsConstantDeclaration(t, true, nuParse(t, "const x = 1;").GetFirstFirstChild())

	nuh_assertIsConstantDeclaration(t, true, nuParse(t, "/** @const */ var x = 1;").GetFirstFirstChild())
	nuh_assertIsConstantDeclaration(t, true, nuParse(t, "var /** @const */ x = 1;").GetFirstFirstChild())
	nuh_assertIsConstantDeclaration(t, false, nuParse(t, "var x, /** @const */ y = 1;").GetFirstFirstChild())

	nuh_assertIsConstantDeclaration(t, true, nuh_getNameNodeFrom(t, "const [a] = [];", "a"))
	nuh_assertIsConstantDeclaration(t, true, nuh_getNameNodeFrom(t, "const [[[a]]] = [];", "a"))
	nuh_assertIsConstantDeclaration(t, true, nuh_getNameNodeFrom(t, "const [a = 1] = [];", "a"))
	nuh_assertIsConstantDeclaration(t, true, nuh_getNameNodeFrom(t, "const [...a] = [];", "a"))

	nuh_assertIsConstantDeclaration(t, true, nuh_getNameNodeFrom(t, "const {a} = {};", "a"))
	nuh_assertIsConstantDeclaration(t, true, nuh_getNameNodeFrom(t, "const {a = 1} = {};", "a"))
	nuh_assertIsConstantDeclaration(t, true, nuh_getNameNodeFrom(t, "const {b: a} = {};", "a"))
	nuh_assertIsConstantDeclaration(t, true, nuh_getNameNodeFrom(t, "const {[3]: a} = {};", "a"))
	nuh_assertIsConstantDeclaration(t, true, nuh_getNameNodeFrom(t, "const {a: [a]} = {};", "a"))
	nuh_assertIsConstantDeclaration(t,
		true, nuh_getNameNodeFrom(t, "/** @const */ var x = function a() {};", "a"))

	nuh_assertIsConstantDeclaration(t, false, nuh_getNameNodeFrom(t, "var FOO = 1;", "FOO"))

	nuh_assertIsConstantDeclaration(t, true, nuh_constructInferredConstantDeclaration(t))

	nuh_assertIsConstantDeclaration(t,
		false, nuParse(t, "class C {x = 2;}").GetFirstChild().GetLastChild().GetFirstChild())
	nuh_assertIsConstantDeclaration(t,
		true,
		nuParse(t, "class C {/** @const */ x = 2;}").GetFirstChild().GetLastChild().GetFirstChild())

	nuh_assertIsConstantDeclaration(t,
		false, nuParse(t, "class C { [x] = 2;}").GetFirstChild().GetLastChild().GetFirstChild())
	nuh_assertIsConstantDeclaration(t,
		true,
		nuParse(t, "class C { /** @const */ [x] = 2;}").GetFirstChild().GetLastChild().GetFirstChild())
}

func TestNodeUtil_AssortedTests_TestIsConstantDeclarations_FunctionClassLiterals(t *testing.T) {
	nuh_assertIsConstantDeclaration(t, false, nuh_getNameNodeFrom(t, "function Foo() {}", "Foo"))
	nuh_assertIsConstantDeclaration(t, false, nuh_getNameNodeFrom(t, "class Foo {}", "Foo"))
}

func TestNodeUtil_AssortedTests_TestIsConstantDeclaration_throwsOnNonDeclarationReferences(t *testing.T) {
	assertThrows[jutil.IIllegalStateException](t, func() { NodeUtil_IsConstantDeclaration(nil, nuh_getNameNodeFrom(t, "const x = y;", "y")) })

	assertThrows[jutil.IIllegalStateException](t, func() { NodeUtil_IsConstantDeclaration(nil, nuh_getNameNodeFrom(t, "x;", "x")) })

	assertThrows[jutil.IIllegalStateException](t, func() { NodeUtil_IsConstantDeclaration(nil, nuh_getNameNodeFrom(t, "const ns = {y: x};", "x")) })

	constAssignment := nuParse(t, "/** @const */ x.y = a.b;")
	rhs := constAssignment.GetFirstFirstChild().GetSecondChild()

	assertThrows[jutil.IIllegalArgumentException](t, func() { NodeUtil_IsConstantDeclaration(NodeUtil_GetBestJSDocInfo(rhs), rhs) })
}

func TestNodeUtil_AssortedTests_TestIsConstantDeclaration_qnames(t *testing.T) {
	constAssignment := nuParse(t, "/** @const */ x.y = a.b;")
	assign := constAssignment.GetFirstFirstChild()
	nuh_assertIsConstantDeclaration(t, true, assign.GetFirstChild())

	constByConventionAssignment := nuParse(t, "x.Y = a.b;")
	assign = constByConventionAssignment.GetFirstFirstChild()
	nuh_assertIsConstantDeclaration(t, false, assign.GetFirstChild())

	nonConstAssignment := nuParse(t, "x.y = a.b;")
	assign = nonConstAssignment.GetFirstFirstChild()
	nuh_assertIsConstantDeclaration(t, false, assign.GetFirstChild())

	expression := nuParse(t, "/** @const */ x.y;") // valid in externs
	nuh_assertIsConstantDeclaration(t, true, expression.GetFirstFirstChild())
}

func TestNodeUtil_AssortedTests_TestIsConstantDeclaration_keys(t *testing.T) {
	nuh_assertIsConstantDeclaration(t,
		true, nuh_getStringKeyNodeFrom(t, "const ns = {/** @const */ x: y};", "x"))
	nuh_assertIsConstantDeclaration(t,
		false, nuh_getStringKeyNodeFrom(t, "/** @const */ const ns = {x: y};", "x"))
	nuh_assertIsConstantDeclaration(t, false, nuh_getStringKeyNodeFrom(t, "const ns = {x: y};", "x"))
}

func nuh_assertIsConstantDeclaration(t *testing.T, isConstantDeclaration bool, node rhino.Node) {

	jsDocInfo := NodeUtil_GetBestJSDocInfo(node)
	nuAssertWithMessage(t, "Is %s a constant declaration?", node).That(NodeUtil_IsConstantDeclaration(jsDocInfo, node)).IsEqualTo(isConstantDeclaration)
}

/** Returns a NAME node from 'let foo = 1;` */
func nuh_constructInferredConstantDeclaration(t *testing.T) rhino.Node {
	nameNode := rhino.IR_Name("foo")
	nameNode.SetInferredConstantVar(true)
	rhino.IR_ScriptOfNodes(rhino.IR_LetWithValue(nameNode, rhino.IR_Number(1))) // attach the expected context to the name.
	return nameNode
}

func TestNodeUtil_AssortedTests_TestDestructuring1(t *testing.T) {
	root := nuParse(t, "var [a, b] = obj;")
	varNode := root.GetFirstChild()
	destructLhs := varNode.GetFirstChild()
	jutil.CheckStateMsg1(destructLhs.IsDestructuringLhs(), "%s", destructLhs)
	destructPat := destructLhs.GetFirstChild()
	jutil.CheckStateMsg1(destructPat.IsArrayPattern(), "%s", destructPat)

	nameNodeA := destructPat.GetFirstChild()
	nameNodeB := nameNodeA.GetNext()
	jutil.CheckStateMsg1(nameNodeA.GetString() == "a", "%s", nameNodeA)
	jutil.CheckStateMsg1(nameNodeB.GetString() == "b", "%s", nameNodeB)

	nameNodeObj := destructPat.GetNext()
	jutil.CheckStateMsg1(nameNodeObj.GetString() == "obj", "%s", nameNodeObj)

	nuh_assertLhsByDestructuring(t, nameNodeA)
	nuh_assertLhsByDestructuring(t, nameNodeB)
	nuh_assertNotLhsByDestructuring(t, nameNodeObj)

	nuAssertThat(t, NodeUtil_GetRootTarget(destructPat)).IsSameInstanceAs(destructPat)
	nuAssertThat(t, NodeUtil_GetRootTarget(nameNodeA)).IsSameInstanceAs(destructPat)
	nuAssertThat(t, NodeUtil_GetRootTarget(nameNodeB)).IsSameInstanceAs(destructPat)

	nuAssertThat(t, NodeUtil_GetDeclaringParent(nameNodeA)).IsSameInstanceAs(varNode)
	nuAssertThat(t, NodeUtil_GetDeclaringParent(nameNodeB)).IsSameInstanceAs(varNode)
}

func TestNodeUtil_AssortedTests_TestDestructuring1b(t *testing.T) {
	root := nuParse(t, "var {a: c, b: d} = obj;")
	varNode := root.GetFirstChild()
	destructLhs := varNode.GetFirstChild()
	jutil.CheckStateMsg1(destructLhs.IsDestructuringLhs(), "%s", destructLhs)
	destructPat := destructLhs.GetFirstChild()
	jutil.CheckStateMsg1(destructPat.IsObjectPattern(), "%s", destructPat)

	strKeyNodeA := destructPat.GetFirstChild()
	strKeyNodeB := strKeyNodeA.GetNext()
	nameNodeC := strKeyNodeA.GetFirstChild()
	nameNodeD := strKeyNodeB.GetFirstChild()
	jutil.CheckStateMsg1(strKeyNodeA.GetString() == "a", "%s", strKeyNodeA)
	jutil.CheckStateMsg1(strKeyNodeB.GetString() == "b", "%s", strKeyNodeB)
	jutil.CheckStateMsg1(nameNodeC.GetString() == "c", "%s", nameNodeC)
	jutil.CheckStateMsg1(nameNodeD.GetString() == "d", "%s", nameNodeD)

	nameNodeObj := destructPat.GetNext()
	jutil.CheckStateMsg1(nameNodeObj.GetString() == "obj", "%s", nameNodeObj)

	nuh_assertNotLhsByDestructuring(t, strKeyNodeA)
	nuh_assertNotLhsByDestructuring(t, strKeyNodeB)
	nuh_assertLhsByDestructuring(t, nameNodeC)
	nuh_assertLhsByDestructuring(t, nameNodeD)
	nuh_assertNotLhsByDestructuring(t, nameNodeObj)

	nuAssertThat(t, NodeUtil_GetRootTarget(destructPat)).IsSameInstanceAs(destructPat)
	nuAssertThat(t, NodeUtil_GetRootTarget(nameNodeC)).IsSameInstanceAs(destructPat)
	nuAssertThat(t, NodeUtil_GetRootTarget(nameNodeD)).IsSameInstanceAs(destructPat)

	nuAssertThat(t, NodeUtil_GetDeclaringParent(nameNodeC)).IsSameInstanceAs(varNode)
	nuAssertThat(t, NodeUtil_GetDeclaringParent(nameNodeD)).IsSameInstanceAs(varNode)
}

func TestNodeUtil_AssortedTests_TestDestructuring1c(t *testing.T) {
	root := nuParse(t, "var {a, b} = obj;")
	varNode := root.GetFirstChild()
	destructLhs := varNode.GetFirstChild()
	jutil.CheckStateMsg1(destructLhs.IsDestructuringLhs(), "%s", destructLhs)
	destructPat := destructLhs.GetFirstChild()
	jutil.CheckStateMsg1(destructPat.IsObjectPattern(), "%s", destructPat)

	strKeyNodeA := destructPat.GetFirstChild()
	strKeyNodeB := strKeyNodeA.GetNext()
	nameNodeA := strKeyNodeA.GetFirstChild()
	nameNodeB := strKeyNodeB.GetFirstChild()
	jutil.CheckStateMsg1(strKeyNodeA.GetString() == "a", "%s", strKeyNodeA)
	jutil.CheckStateMsg1(nameNodeA.GetString() == "a", "%s", nameNodeA)
	jutil.CheckStateMsg1(strKeyNodeB.GetString() == "b", "%s", strKeyNodeB)
	jutil.CheckStateMsg1(nameNodeB.GetString() == "b", "%s", nameNodeB)

	nameNodeObj := destructPat.GetNext()
	jutil.CheckStateMsg1(nameNodeObj.GetString() == "obj", "%s", nameNodeObj)

	nuh_assertNotLhsByDestructuring(t, strKeyNodeA)
	nuh_assertNotLhsByDestructuring(t, strKeyNodeB)
	nuh_assertLhsByDestructuring(t, nameNodeA)
	nuh_assertLhsByDestructuring(t, nameNodeB)
	nuh_assertNotLhsByDestructuring(t, nameNodeObj)

	nuAssertThat(t, NodeUtil_GetRootTarget(destructPat)).IsSameInstanceAs(destructPat)
	nuAssertThat(t, NodeUtil_GetRootTarget(nameNodeA)).IsSameInstanceAs(destructPat)
	nuAssertThat(t, NodeUtil_GetRootTarget(nameNodeB)).IsSameInstanceAs(destructPat)

	nuAssertThat(t, NodeUtil_GetDeclaringParent(nameNodeA)).IsSameInstanceAs(varNode)
	nuAssertThat(t, NodeUtil_GetDeclaringParent(nameNodeB)).IsSameInstanceAs(varNode)
}

func TestNodeUtil_AssortedTests_TestDestructuring1d(t *testing.T) {
	root := nuParse(t, "var {a = defaultValue} = obj;")
	varNode := root.GetFirstChild()
	destructLhs := varNode.GetFirstChild()
	jutil.CheckStateMsg1(destructLhs.IsDestructuringLhs(), "%s", destructLhs)
	destructPat := destructLhs.GetFirstChild()
	jutil.CheckStateMsg1(destructPat.IsObjectPattern(), "%s", destructPat)

	strKeyNodeA := destructPat.GetFirstChild()
	jutil.CheckState(strKeyNodeA.IsStringKey())
	jutil.CheckStateMsg1(strKeyNodeA.GetString() == "a", "%s", strKeyNodeA)

	defaultNodeA := strKeyNodeA.GetFirstChild()
	jutil.CheckStateMsg1(defaultNodeA.IsDefaultValue(), "%s", defaultNodeA)
	nameNodeA := defaultNodeA.GetFirstChild()
	jutil.CheckStateMsg1(nameNodeA.GetString() == "a", "%s", nameNodeA)

	nameNodeDefault := defaultNodeA.GetSecondChild()
	jutil.CheckStateMsg1(nameNodeDefault.GetString() == "defaultValue", "%s", nameNodeDefault)
	nameNodeObj := destructPat.GetNext()
	jutil.CheckStateMsg1(nameNodeObj.GetString() == "obj", "%s", nameNodeObj)

	nuh_assertNotLhsByDestructuring(t, strKeyNodeA)
	nuh_assertNotLhsByDestructuring(t, defaultNodeA)
	nuh_assertLhsByDestructuring(t, nameNodeA)
	nuh_assertNotLhsByDestructuring(t, nameNodeDefault)
	nuh_assertNotLhsByDestructuring(t, nameNodeObj)

	nuAssertThat(t, NodeUtil_GetRootTarget(destructPat)).IsSameInstanceAs(destructPat)
	nuAssertThat(t, NodeUtil_GetRootTarget(nameNodeA)).IsSameInstanceAs(destructPat)

	nuAssertThat(t, NodeUtil_GetDeclaringParent(nameNodeA)).IsSameInstanceAs(varNode)
}

func TestNodeUtil_AssortedTests_TestDestructuring1e(t *testing.T) {
	root := nuParse(t, "var {a: b = defaultValue} = obj;")
	varNode := root.GetFirstChild()
	destructLhs := varNode.GetFirstChild()
	jutil.CheckStateMsg1(destructLhs.IsDestructuringLhs(), "%s", destructLhs)
	destructPat := destructLhs.GetFirstChild()
	jutil.CheckStateMsg1(destructPat.IsObjectPattern(), "%s", destructPat)

	strKeyNodeA := destructPat.GetFirstChild()
	jutil.CheckStateMsg1(strKeyNodeA.GetString() == "a", "%s", strKeyNodeA)

	defaultNodeA := strKeyNodeA.GetOnlyChild()
	jutil.CheckStateMsg1(defaultNodeA.IsDefaultValue(), "%s", defaultNodeA)

	nameNodeB := defaultNodeA.GetFirstChild()
	jutil.CheckStateMsg1(nameNodeB.GetString() == "b", "%s", nameNodeB)

	nameNodeDefaultValue := defaultNodeA.GetSecondChild()
	jutil.CheckStateMsg1(nameNodeDefaultValue.GetString() == "defaultValue", "%s", nameNodeDefaultValue)
	nameNodeObj := destructPat.GetNext()
	jutil.CheckStateMsg1(nameNodeObj.GetString() == "obj", "%s", nameNodeObj)

	nuh_assertNotLhsByDestructuring(t, strKeyNodeA)
	nuh_assertLhsByDestructuring(t, nameNodeB)
	nuh_assertNotLhsByDestructuring(t, nameNodeDefaultValue)
	nuh_assertNotLhsByDestructuring(t, nameNodeObj)

	nuAssertThat(t, NodeUtil_GetRootTarget(destructPat)).IsSameInstanceAs(destructPat)
	nuAssertThat(t, NodeUtil_GetRootTarget(nameNodeB)).IsSameInstanceAs(destructPat)

	nuAssertThat(t, NodeUtil_GetDeclaringParent(nameNodeB)).IsSameInstanceAs(varNode)
}

func TestNodeUtil_AssortedTests_TestDestructuring1f(t *testing.T) {
	root := nuParse(t, "var [a  = defaultValue] = arr;")
	varNode := root.GetFirstChild()
	destructLhs := varNode.GetFirstChild()
	jutil.CheckStateMsg1(destructLhs.IsDestructuringLhs(), "%s", destructLhs)
	destructPat := destructLhs.GetFirstChild()
	jutil.CheckStateMsg1(destructPat.IsArrayPattern(), "%s", destructPat)

	defaultNodeA := destructPat.GetFirstChild()

	jutil.CheckStateMsg1(defaultNodeA.IsDefaultValue(), "%s", defaultNodeA)

	nameNodeA := defaultNodeA.GetFirstChild()
	jutil.CheckStateMsg1(nameNodeA.GetString() == "a", "%s", nameNodeA)

	nameNodeDefault := defaultNodeA.GetSecondChild()
	jutil.CheckStateMsg1(nameNodeDefault.GetString() == "defaultValue", "%s", nameNodeDefault)
	nameNodeObj := destructPat.GetNext()
	jutil.CheckStateMsg1(nameNodeObj.GetString() == "arr", "%s", nameNodeObj)

	nuh_assertLhsByDestructuring(t, nameNodeA)
	nuh_assertNotLhsByDestructuring(t, nameNodeDefault)
	nuh_assertNotLhsByDestructuring(t, nameNodeObj)

	nuAssertThat(t, NodeUtil_GetRootTarget(destructPat)).IsSameInstanceAs(destructPat)
	nuAssertThat(t, NodeUtil_GetRootTarget(nameNodeA)).IsSameInstanceAs(destructPat)

	nuAssertThat(t, NodeUtil_GetDeclaringParent(nameNodeA)).IsSameInstanceAs(varNode)
}

func TestNodeUtil_AssortedTests_TestDestructuring2(t *testing.T) {
	root := nuParse(t, "var [a, [b, c]] = obj;")
	varNode := root.GetFirstChild()
	destructLhs := varNode.GetFirstChild()
	jutil.CheckStateMsg1(destructLhs.IsDestructuringLhs(), "%s", destructLhs)
	destructPat := destructLhs.GetFirstChild()
	jutil.CheckStateMsg1(destructPat.IsArrayPattern(), "%s", destructPat)

	nameNodeA := destructPat.GetFirstChild()
	innerPat := nameNodeA.GetNext()
	nameNodeB := innerPat.GetFirstChild()
	nameNodeC := nameNodeB.GetNext()
	jutil.CheckStateMsg1(nameNodeA.GetString() == "a", "%s", nameNodeA)
	jutil.CheckStateMsg1(nameNodeB.GetString() == "b", "%s", nameNodeB)
	jutil.CheckStateMsg1(nameNodeC.GetString() == "c", "%s", nameNodeC)

	nameNodeObj := destructPat.GetNext()
	jutil.CheckStateMsg1(nameNodeObj.GetString() == "obj", "%s", nameNodeObj)

	nuh_assertLhsByDestructuring(t, nameNodeA)
	nuh_assertLhsByDestructuring(t, nameNodeB)
	nuh_assertLhsByDestructuring(t, nameNodeC)
	nuh_assertNotLhsByDestructuring(t, nameNodeObj)

	nuAssertThat(t, NodeUtil_GetRootTarget(destructPat)).IsSameInstanceAs(destructPat)
	nuAssertThat(t, NodeUtil_GetRootTarget(nameNodeA)).IsSameInstanceAs(destructPat)
	nuAssertThat(t, NodeUtil_GetRootTarget(innerPat)).IsSameInstanceAs(destructPat)
	nuAssertThat(t, NodeUtil_GetRootTarget(nameNodeB)).IsSameInstanceAs(destructPat)
	nuAssertThat(t, NodeUtil_GetRootTarget(nameNodeC)).IsSameInstanceAs(destructPat)

	nuAssertThat(t, NodeUtil_GetDeclaringParent(nameNodeA)).IsSameInstanceAs(varNode)
	nuAssertThat(t, NodeUtil_GetDeclaringParent(nameNodeB)).IsSameInstanceAs(varNode)
	nuAssertThat(t, NodeUtil_GetDeclaringParent(nameNodeC)).IsSameInstanceAs(varNode)
}

func TestNodeUtil_AssortedTests_TestDestructuring2b(t *testing.T) {
	root := nuParse(t, "var {a: e, b: {c: f, d: g}} = obj;")
	varNode := root.GetFirstChild()
	destructLhs := varNode.GetFirstChild()
	jutil.CheckStateMsg1(destructLhs.IsDestructuringLhs(), "%s", destructLhs)
	destructPat := destructLhs.GetFirstChild()
	jutil.CheckStateMsg1(destructPat.IsObjectPattern(), "%s", destructPat)

	strKeyNodeA := destructPat.GetFirstChild()
	strKeyNodeB := strKeyNodeA.GetNext()
	innerPat := strKeyNodeB.GetOnlyChild()
	strKeyNodeC := innerPat.GetFirstChild()
	strKeyNodeD := strKeyNodeC.GetNext()
	nameNodeE := strKeyNodeA.GetFirstChild()
	nameNodeF := strKeyNodeC.GetFirstChild()
	nameNodeG := strKeyNodeD.GetFirstChild()
	jutil.CheckStateMsg1(strKeyNodeA.GetString() == "a", "%s", strKeyNodeA)
	jutil.CheckStateMsg1(strKeyNodeB.GetString() == "b", "%s", strKeyNodeB)
	jutil.CheckStateMsg1(strKeyNodeC.GetString() == "c", "%s", strKeyNodeC)
	jutil.CheckStateMsg1(strKeyNodeD.GetString() == "d", "%s", strKeyNodeD)
	jutil.CheckStateMsg1(nameNodeE.GetString() == "e", "%s", nameNodeE)
	jutil.CheckStateMsg1(nameNodeF.GetString() == "f", "%s", nameNodeF)
	jutil.CheckStateMsg1(nameNodeG.GetString() == "g", "%s", nameNodeG)

	nameNodeObj := destructPat.GetNext()
	jutil.CheckStateMsg1(nameNodeObj.GetString() == "obj", "%s", nameNodeObj)

	nuh_assertNotLhsByDestructuring(t, strKeyNodeA)
	nuh_assertNotLhsByDestructuring(t, strKeyNodeB)
	nuh_assertNotLhsByDestructuring(t, strKeyNodeC)
	nuh_assertNotLhsByDestructuring(t, strKeyNodeD)
	nuh_assertLhsByDestructuring(t, nameNodeE)
	nuh_assertLhsByDestructuring(t, nameNodeF)
	nuh_assertLhsByDestructuring(t, nameNodeG)
	nuh_assertNotLhsByDestructuring(t, nameNodeObj)

	nuAssertThat(t, NodeUtil_GetRootTarget(destructPat)).IsSameInstanceAs(destructPat)
	nuAssertThat(t, NodeUtil_GetRootTarget(nameNodeE)).IsSameInstanceAs(destructPat)
	nuAssertThat(t, NodeUtil_GetRootTarget(innerPat)).IsSameInstanceAs(destructPat)
	nuAssertThat(t, NodeUtil_GetRootTarget(nameNodeF)).IsSameInstanceAs(destructPat)
	nuAssertThat(t, NodeUtil_GetRootTarget(nameNodeG)).IsSameInstanceAs(destructPat)

	nuAssertThat(t, NodeUtil_GetDeclaringParent(nameNodeE)).IsSameInstanceAs(varNode)
	nuAssertThat(t, NodeUtil_GetDeclaringParent(nameNodeF)).IsSameInstanceAs(varNode)
	nuAssertThat(t, NodeUtil_GetDeclaringParent(nameNodeG)).IsSameInstanceAs(varNode)
}

func TestNodeUtil_AssortedTests_TestDestructuring3(t *testing.T) {
	root := nuParse(t, "var [a, b] = [c, d];")
	varNode := root.GetFirstChild()
	destructLhs := varNode.GetFirstChild()
	jutil.CheckStateMsg1(destructLhs.IsDestructuringLhs(), "%s", destructLhs)
	destructPat := destructLhs.GetFirstChild()
	jutil.CheckStateMsg1(destructPat.IsArrayPattern(), "%s", destructPat)

	nameNodeA := destructPat.GetFirstChild()
	nameNodeB := nameNodeA.GetNext()
	jutil.CheckStateMsg1(nameNodeA.GetString() == "a", "%s", nameNodeA)
	jutil.CheckStateMsg1(nameNodeB.GetString() == "b", "%s", nameNodeB)

	nameNodeC := destructLhs.GetLastChild().GetFirstChild()
	nameNodeD := nameNodeC.GetNext()
	jutil.CheckStateMsg1(nameNodeC.GetString() == "c", "%s", nameNodeC)
	jutil.CheckStateMsg1(nameNodeD.GetString() == "d", "%s", nameNodeD)

	nuh_assertLhsByDestructuring(t, nameNodeA)
	nuh_assertLhsByDestructuring(t, nameNodeB)
	nuh_assertNotLhsByDestructuring(t, nameNodeC)
	nuh_assertNotLhsByDestructuring(t, nameNodeD)

	nuAssertThat(t, NodeUtil_GetRootTarget(destructPat)).IsSameInstanceAs(destructPat)
	nuAssertThat(t, NodeUtil_GetRootTarget(nameNodeA)).IsSameInstanceAs(destructPat)
	nuAssertThat(t, NodeUtil_GetRootTarget(nameNodeB)).IsSameInstanceAs(destructPat)

	nuAssertThat(t, NodeUtil_GetDeclaringParent(nameNodeA)).IsSameInstanceAs(varNode)
	nuAssertThat(t, NodeUtil_GetDeclaringParent(nameNodeB)).IsSameInstanceAs(varNode)
}

func TestNodeUtil_AssortedTests_TestDestructuring3b(t *testing.T) {
	root := nuParse(t, "var {a: c, b: d} = {a: 1, b: 2};")
	varNode := root.GetFirstChild()
	destructLhs := varNode.GetFirstChild()
	jutil.CheckStateMsg1(destructLhs.IsDestructuringLhs(), "%s", destructLhs)
	destructPat := destructLhs.GetFirstChild()
	jutil.CheckStateMsg1(destructPat.IsObjectPattern(), "%s", destructPat)

	strKeyNodeA := destructPat.GetFirstChild()
	strKeyNodeB := strKeyNodeA.GetNext()
	nameNodeC := strKeyNodeA.GetFirstChild()
	nameNodeD := strKeyNodeB.GetFirstChild()
	jutil.CheckStateMsg1(strKeyNodeA.GetString() == "a", "%s", strKeyNodeA)
	jutil.CheckStateMsg1(strKeyNodeB.GetString() == "b", "%s", strKeyNodeB)
	jutil.CheckStateMsg1(nameNodeC.GetString() == "c", "%s", nameNodeC)
	jutil.CheckStateMsg1(nameNodeD.GetString() == "d", "%s", nameNodeD)

	nuh_assertNotLhsByDestructuring(t, strKeyNodeA)
	nuh_assertNotLhsByDestructuring(t, strKeyNodeB)
	nuh_assertLhsByDestructuring(t, nameNodeC)
	nuh_assertLhsByDestructuring(t, nameNodeD)

	nuAssertThat(t, NodeUtil_GetRootTarget(destructPat)).IsSameInstanceAs(destructPat)
	nuAssertThat(t, NodeUtil_GetRootTarget(nameNodeC)).IsSameInstanceAs(destructPat)
	nuAssertThat(t, NodeUtil_GetRootTarget(nameNodeD)).IsSameInstanceAs(destructPat)

	nuAssertThat(t, NodeUtil_GetDeclaringParent(nameNodeC)).IsSameInstanceAs(varNode)
	nuAssertThat(t, NodeUtil_GetDeclaringParent(nameNodeD)).IsSameInstanceAs(varNode)
}

func TestNodeUtil_AssortedTests_TestDestructuring4(t *testing.T) {
	root := nuParse(t, "for ([a, b] of X){}")
	destructPat := root.GetFirstFirstChild()
	jutil.CheckStateMsg1(destructPat.IsArrayPattern(), "%s", destructPat)

	nameNodeA := destructPat.GetFirstChild()
	nameNodeB := destructPat.GetLastChild()
	jutil.CheckStateMsg1(nameNodeA.GetString() == "a", "%s", nameNodeA)
	jutil.CheckStateMsg1(nameNodeB.GetString() == "b", "%s", nameNodeB)

	nuh_assertLhsByDestructuring(t, nameNodeA)
	nuh_assertLhsByDestructuring(t, nameNodeB)

	nuAssertThat(t, NodeUtil_GetRootTarget(destructPat)).IsSameInstanceAs(destructPat)
	nuAssertThat(t, NodeUtil_GetRootTarget(nameNodeA)).IsSameInstanceAs(destructPat)
	nuAssertThat(t, NodeUtil_GetRootTarget(nameNodeB)).IsSameInstanceAs(destructPat)
}

func TestNodeUtil_AssortedTests_TestDestructuring5(t *testing.T) {
	root := nuParse(t, "function fn([a, b] = [c, d]){}")
	destructPat := root.GetFirstChild().GetSecondChild().GetFirstFirstChild()
	jutil.CheckStateMsg1(destructPat.IsArrayPattern(), "%s", destructPat)

	nameNodeA := destructPat.GetFirstChild()
	nameNodeB := destructPat.GetLastChild()
	nameNodeC := destructPat.GetNext().GetFirstChild()
	nameNodeD := destructPat.GetNext().GetLastChild()
	jutil.CheckStateMsg1(nameNodeA.GetString() == "a", "%s", nameNodeA)
	jutil.CheckStateMsg1(nameNodeB.GetString() == "b", "%s", nameNodeB)
	jutil.CheckStateMsg1(nameNodeC.GetString() == "c", "%s", nameNodeC)
	jutil.CheckStateMsg1(nameNodeD.GetString() == "d", "%s", nameNodeD)

	nuh_assertLhsByDestructuring(t, nameNodeA)
	nuh_assertLhsByDestructuring(t, nameNodeB)
	nuh_assertNotLhsByDestructuring(t, nameNodeC)
	nuh_assertNotLhsByDestructuring(t, nameNodeD)

	nuAssertThat(t, NodeUtil_GetRootTarget(destructPat)).IsSameInstanceAs(destructPat)
	nuAssertThat(t, NodeUtil_GetRootTarget(nameNodeA)).IsSameInstanceAs(destructPat)
	nuAssertThat(t, NodeUtil_GetRootTarget(nameNodeB)).IsSameInstanceAs(destructPat)
}

func TestNodeUtil_AssortedTests_TestRestParameter(t *testing.T) {
	root := nuParse(t, "function fn(a, b, ...rest){}")
	parameterList := root.GetFirstChild().GetSecondChild()

	nameNodeA := parameterList.GetFirstChild()
	nameNodeB := nameNodeA.GetNext()
	restNode := nameNodeB.GetNext()
	jutil.CheckStateMsg1(restNode.IsRest(), "%s", restNode)

	nameNodeRest := restNode.GetFirstChild()

	jutil.CheckStateMsg1(nameNodeA.GetString() == "a", "%s", nameNodeA)

	nuAssertThat(t, NodeUtil_GetRootTarget(nameNodeA)).IsSameInstanceAs(nameNodeA)
	nuAssertThat(t, NodeUtil_GetRootTarget(nameNodeB)).IsSameInstanceAs(nameNodeB)
	nuAssertThat(t, NodeUtil_GetRootTarget(nameNodeRest)).IsSameInstanceAs(nameNodeRest)

	nuAssertThat(t, NodeUtil_GetDeclaringParent(nameNodeA)).IsSameInstanceAs(parameterList)
	nuAssertThat(t, NodeUtil_GetDeclaringParent(nameNodeB)).IsSameInstanceAs(parameterList)
	nuAssertThat(t, NodeUtil_GetDeclaringParent(nameNodeRest)).IsSameInstanceAs(parameterList)
}

func TestNodeUtil_AssortedTests_TestDestructuring6(t *testing.T) {
	root := nuParse(t, "for ([{a: b}] of c) {}")
	forOfNode := root.GetFirstChild()
	forOfTargetPat := forOfNode.GetFirstChild()
	innerPat := forOfTargetPat.GetFirstChild()
	jutil.CheckStateMsg1(innerPat.IsObjectPattern(), "%s", innerPat)
	jutil.CheckStateMsg1(forOfTargetPat.IsArrayPattern(), "%s", forOfTargetPat)

	strKeyNodeA := innerPat.GetFirstChild()
	nameNodeB := strKeyNodeA.GetFirstChild()
	nameNodeC := innerPat.GetParent().GetNext()
	jutil.CheckStateMsg1(strKeyNodeA.GetString() == "a", "%s", strKeyNodeA)
	jutil.CheckStateMsg1(nameNodeB.GetString() == "b", "%s", nameNodeB)
	jutil.CheckStateMsg1(nameNodeC.GetString() == "c", "%s", nameNodeC)

	nuh_assertNotLhsByDestructuring(t, strKeyNodeA)
	nuh_assertLhsByDestructuring(t, nameNodeB)
	nuh_assertNotLhsByDestructuring(t, nameNodeC)

	nuAssertThat(t, NodeUtil_GetRootTarget(forOfTargetPat)).IsSameInstanceAs(forOfTargetPat)
	nuAssertThat(t, NodeUtil_GetRootTarget(innerPat)).IsSameInstanceAs(forOfTargetPat)
	nuAssertThat(t, NodeUtil_GetRootTarget(nameNodeB)).IsSameInstanceAs(forOfTargetPat)
}

func TestNodeUtil_AssortedTests_TestDestructuring6b(t *testing.T) {
	root := nuParse(t, "for ([{a: b}] in c) {}")
	forInNode := root.GetFirstChild()
	forInTargetPat := forInNode.GetFirstChild()
	innerPat := forInTargetPat.GetFirstChild()
	jutil.CheckStateMsg1(innerPat.IsObjectPattern(), "%s", innerPat)
	jutil.CheckStateMsg1(forInTargetPat.IsArrayPattern(), "%s", forInTargetPat)

	strKeyNodeA := innerPat.GetFirstChild()
	nameNodeB := strKeyNodeA.GetFirstChild()
	nameNodeC := innerPat.GetParent().GetNext()
	jutil.CheckStateMsg1(strKeyNodeA.GetString() == "a", "%s", strKeyNodeA)
	jutil.CheckStateMsg1(nameNodeB.GetString() == "b", "%s", nameNodeB)
	jutil.CheckStateMsg1(nameNodeC.GetString() == "c", "%s", nameNodeC)

	nuh_assertNotLhsByDestructuring(t, strKeyNodeA)
	nuh_assertLhsByDestructuring(t, nameNodeB)
	nuh_assertNotLhsByDestructuring(t, nameNodeC)

	nuAssertThat(t, NodeUtil_GetRootTarget(forInTargetPat)).IsSameInstanceAs(forInTargetPat)
	nuAssertThat(t, NodeUtil_GetRootTarget(innerPat)).IsSameInstanceAs(forInTargetPat)
	nuAssertThat(t, NodeUtil_GetRootTarget(nameNodeB)).IsSameInstanceAs(forInTargetPat)
}

func TestNodeUtil_AssortedTests_TestDestructuring6c(t *testing.T) {
	root := nuParse(t, "for (var [{a: b}] = [{a: 1}];;) {}")
	varNode := root.GetFirstFirstChild()
	destructArr := varNode.GetFirstFirstChild()
	jutil.CheckStateMsg1(destructArr.IsArrayPattern(), "%s", destructArr)
	destructPat := destructArr.GetFirstChild()
	jutil.CheckStateMsg1(destructPat.IsObjectPattern(), "%s", destructPat)

	strKeyNodeA := destructPat.GetFirstChild()
	nameNodeB := strKeyNodeA.GetFirstChild()
	jutil.CheckStateMsg1(strKeyNodeA.GetString() == "a", "%s", strKeyNodeA)
	jutil.CheckStateMsg1(nameNodeB.GetString() == "b", "%s", nameNodeB)

	nuh_assertNotLhsByDestructuring(t, strKeyNodeA)
	nuh_assertLhsByDestructuring(t, nameNodeB)

	nuAssertThat(t, NodeUtil_GetRootTarget(destructArr)).IsSameInstanceAs(destructArr)
	nuAssertThat(t, NodeUtil_GetRootTarget(destructPat)).IsSameInstanceAs(destructArr)
	nuAssertThat(t, NodeUtil_GetRootTarget(nameNodeB)).IsSameInstanceAs(destructArr)

	nuAssertThat(t, NodeUtil_GetDeclaringParent(nameNodeB)).IsSameInstanceAs(varNode)
}

func TestNodeUtil_AssortedTests_TestDestructuring7(t *testing.T) {
	root := nuParse(t, "for ([a] of c) {}")
	destructPat := root.GetFirstFirstChild()
	jutil.CheckStateMsg1(destructPat.IsArrayPattern(), "%s", destructPat)

	nameNodeA := destructPat.GetFirstChild()
	nameNodeC := destructPat.GetNext()
	jutil.CheckStateMsg1(nameNodeA.GetString() == "a", "%s", nameNodeA)
	jutil.CheckStateMsg1(nameNodeC.GetString() == "c", "%s", nameNodeC)

	nuh_assertLhsByDestructuring(t, nameNodeA)
	nuh_assertNotLhsByDestructuring(t, nameNodeC)

	nuAssertThat(t, NodeUtil_GetRootTarget(destructPat)).IsSameInstanceAs(destructPat)
	nuAssertThat(t, NodeUtil_GetRootTarget(nameNodeA)).IsSameInstanceAs(destructPat)
}

func TestNodeUtil_AssortedTests_TestDestructuring7b(t *testing.T) {
	root := nuParse(t, "for ([a] in c) {}")
	destructPat := root.GetFirstFirstChild()
	jutil.CheckStateMsg1(destructPat.IsArrayPattern(), "%s", destructPat)

	nameNodeA := destructPat.GetFirstChild()
	nameNodeC := destructPat.GetNext()
	jutil.CheckStateMsg1(nameNodeA.GetString() == "a", "%s", nameNodeA)
	jutil.CheckStateMsg1(nameNodeC.GetString() == "c", "%s", nameNodeC)

	nuh_assertLhsByDestructuring(t, nameNodeA)
	nuh_assertNotLhsByDestructuring(t, nameNodeC)

	nuAssertThat(t, NodeUtil_GetRootTarget(destructPat)).IsSameInstanceAs(destructPat)
	nuAssertThat(t, NodeUtil_GetRootTarget(nameNodeA)).IsSameInstanceAs(destructPat)
}

func TestNodeUtil_AssortedTests_TestDestructuring7c(t *testing.T) {
	root := nuParse(t, "for (var [a] = [1];;) {}")
	varNode := root.GetFirstFirstChild()
	destructLhs := varNode.GetFirstChild()
	jutil.CheckStateMsg1(destructLhs.IsDestructuringLhs(), "%s", destructLhs)
	destructPat := destructLhs.GetFirstChild()
	jutil.CheckStateMsg1(destructPat.IsArrayPattern(), "%s", destructPat)

	nameNodeA := destructPat.GetFirstChild()
	jutil.CheckStateMsg1(nameNodeA.GetString() == "a", "%s", nameNodeA)

	nuh_assertLhsByDestructuring(t, nameNodeA)

	nuAssertThat(t, NodeUtil_GetRootTarget(destructPat)).IsSameInstanceAs(destructPat)
	nuAssertThat(t, NodeUtil_GetRootTarget(nameNodeA)).IsSameInstanceAs(destructPat)

	nuAssertThat(t, NodeUtil_GetDeclaringParent(nameNodeA)).IsSameInstanceAs(varNode)
}

func TestNodeUtil_AssortedTests_TestDestructuring7d(t *testing.T) {
	root := nuParse(t, "for (let [a] = [1];;) {}")
	letNode := root.GetFirstFirstChild()
	destructLhs := letNode.GetFirstChild()
	jutil.CheckStateMsg1(destructLhs.IsDestructuringLhs(), "%s", destructLhs)
	destructPat := destructLhs.GetFirstChild()
	jutil.CheckStateMsg1(destructPat.IsArrayPattern(), "%s", destructPat)

	nameNodeA := destructPat.GetFirstChild()
	jutil.CheckStateMsg1(nameNodeA.GetString() == "a", "%s", nameNodeA)

	nuh_assertLhsByDestructuring(t, nameNodeA)

	nuAssertThat(t, NodeUtil_GetRootTarget(destructPat)).IsSameInstanceAs(destructPat)
	nuAssertThat(t, NodeUtil_GetRootTarget(nameNodeA)).IsSameInstanceAs(destructPat)

	nuAssertThat(t, NodeUtil_GetDeclaringParent(nameNodeA)).IsSameInstanceAs(letNode)
}

func TestNodeUtil_AssortedTests_TestDestructuring7e(t *testing.T) {
	root := nuParse(t, "for (const [a] = [1];;) {}")
	constNode := root.GetFirstFirstChild()
	destructLhs := constNode.GetFirstChild()
	jutil.CheckStateMsg1(destructLhs.IsDestructuringLhs(), "%s", destructLhs)
	destructPat := destructLhs.GetFirstChild()
	jutil.CheckStateMsg1(destructPat.IsArrayPattern(), "%s", destructPat)

	nameNodeA := destructPat.GetFirstChild()
	jutil.CheckStateMsg1(nameNodeA.GetString() == "a", "%s", nameNodeA)

	nuh_assertLhsByDestructuring(t, nameNodeA)

	nuAssertThat(t, NodeUtil_GetRootTarget(destructPat)).IsSameInstanceAs(destructPat)
	nuAssertThat(t, NodeUtil_GetRootTarget(nameNodeA)).IsSameInstanceAs(destructPat)

	nuAssertThat(t, NodeUtil_GetDeclaringParent(nameNodeA)).IsSameInstanceAs(constNode)
}

func TestNodeUtil_AssortedTests_TestDestructuring8(t *testing.T) {
	root := nuParse(t, "var [...a] = obj;")
	varNode := root.GetFirstChild()
	destructLhs := varNode.GetFirstChild()
	jutil.CheckStateMsg1(destructLhs.IsDestructuringLhs(), "%s", destructLhs)
	destructPat := destructLhs.GetFirstChild()
	jutil.CheckStateMsg1(destructPat.IsArrayPattern(), "%s", destructPat)
	restNode := destructPat.GetFirstChild()
	jutil.CheckStateMsg1(restNode.IsRest(), "%s", restNode)

	nameNodeA := restNode.GetFirstChild()
	jutil.CheckStateMsg1(nameNodeA.GetString() == "a", "%s", nameNodeA)

	nuh_assertLhsByDestructuring(t, nameNodeA)

	nuAssertThat(t, NodeUtil_GetRootTarget(destructPat)).IsSameInstanceAs(destructPat)
	nuAssertThat(t, NodeUtil_GetRootTarget(nameNodeA)).IsSameInstanceAs(destructPat)

	nuAssertThat(t, NodeUtil_GetDeclaringParent(nameNodeA)).IsSameInstanceAs(varNode)
}

func TestNodeUtil_AssortedTests_TestDestructuring8b(t *testing.T) {
	root := nuParse(t, "([...this.x] = obj);")
	assign := root.GetFirstFirstChild()
	jutil.CheckStateMsg1(assign.IsAssign(), "%s", assign)
	destructPat := assign.GetFirstChild()
	jutil.CheckStateMsg1(destructPat.IsArrayPattern(), "%s", destructPat)
	restNode := destructPat.GetFirstChild()
	jutil.CheckStateMsg1(restNode.IsRest(), "%s", restNode)
	getProp := restNode.GetFirstChild()
	jutil.CheckStateMsg1(getProp.IsGetProp(), "%s", getProp)

	nuh_assertLhsByDestructuring(t, getProp)

	nuAssertThat(t, NodeUtil_GetRootTarget(destructPat)).IsSameInstanceAs(destructPat)
	nuAssertThat(t, NodeUtil_GetRootTarget(getProp)).IsSameInstanceAs(destructPat)
}

func TestNodeUtil_AssortedTests_TestDestructuring9(t *testing.T) {
	root := nuParse(t, "var {['a']:a} = obj;")
	varNode := root.GetFirstChild()
	destructLhs := varNode.GetFirstChild()
	jutil.CheckStateMsg1(destructLhs.IsDestructuringLhs(), "%s", destructLhs)
	destructPat := destructLhs.GetFirstChild()
	jutil.CheckStateMsg1(destructPat.IsObjectPattern(), "%s", destructPat)
	computedPropNode := destructPat.GetFirstChild()
	jutil.CheckStateMsg1(computedPropNode.IsComputedProp(), "%s", computedPropNode)

	nameNodeA := computedPropNode.GetLastChild()
	jutil.CheckStateMsg1(nameNodeA.GetString() == "a", "%s", nameNodeA)

	nuh_assertLhsByDestructuring(t, nameNodeA)

	nuAssertThat(t, NodeUtil_GetRootTarget(destructPat)).IsSameInstanceAs(destructPat)
	nuAssertThat(t, NodeUtil_GetRootTarget(nameNodeA)).IsSameInstanceAs(destructPat)

	nuAssertThat(t, NodeUtil_GetDeclaringParent(nameNodeA)).IsSameInstanceAs(varNode)
}

func TestNodeUtil_AssortedTests_TestDestructuringComputedPropertyWithDefault(t *testing.T) {
	root := nuParse(t, "var {['a']: a = 1} = obj;")
	varNode := root.GetFirstChild()
	destructLhs := varNode.GetFirstChild()
	jutil.CheckStateMsg1(destructLhs.IsDestructuringLhs(), "%s", destructLhs)
	destructPat := destructLhs.GetFirstChild()
	jutil.CheckStateMsg1(destructPat.IsObjectPattern(), "%s", destructPat)
	computedPropNode := destructPat.GetFirstChild()
	jutil.CheckStateMsg1(computedPropNode.IsComputedProp(), "%s", computedPropNode)

	defaultValueNode := computedPropNode.GetLastChild()
	nameNodeA := defaultValueNode.GetFirstChild()
	jutil.CheckStateMsg1(nameNodeA.GetString() == "a", "%s", nameNodeA)

	nuh_assertLhsByDestructuring(t, nameNodeA)

	nuAssertThat(t, NodeUtil_GetRootTarget(destructPat)).IsSameInstanceAs(destructPat)
	nuAssertThat(t, NodeUtil_GetRootTarget(nameNodeA)).IsSameInstanceAs(destructPat)

	nuAssertThat(t, NodeUtil_GetDeclaringParent(nameNodeA)).IsSameInstanceAs(varNode)
}

func TestNodeUtil_AssortedTests_TestDestructuringWithArrayPatternInCatch(t *testing.T) {
	actual := nuParse(t, "try {} catch([a]) {}")

	tryNode := actual.GetFirstChild()
	catchBlock := tryNode.GetSecondChild()
	catchNode := catchBlock.GetFirstChild()
	nuAssertNode(t, catchNode).HasType(rhino.Token_CATCH)

	arrayPattern := catchNode.GetFirstChild()
	nameNodeA := arrayPattern.GetOnlyChild()
	nuAssertNode(t, nameNodeA).IsName("a")

	nuh_assertLhsByDestructuring(t, nameNodeA)

	nuAssertThat(t, NodeUtil_GetRootTarget(arrayPattern)).IsSameInstanceAs(arrayPattern)
	nuAssertThat(t, NodeUtil_GetRootTarget(nameNodeA)).IsSameInstanceAs(arrayPattern)

	nuAssertThat(t, NodeUtil_GetDeclaringParent(nameNodeA)).IsSameInstanceAs(catchNode)
}

func TestNodeUtil_AssortedTests_TestDestructuringObjectRest(t *testing.T) {
	root := nuParse(t, "var {...rest} = obj;")
	varNode := root.GetFirstChild()
	destructLhs := varNode.GetFirstChild()
	jutil.CheckStateMsg1(destructLhs.IsDestructuringLhs(), "%s", destructLhs)
	destructPat := destructLhs.GetFirstChild()
	jutil.CheckStateMsg1(destructPat.IsObjectPattern(), "%s", destructPat)

	restNode := destructPat.GetFirstChild()
	nameNodeRest := restNode.GetFirstChild()
	jutil.CheckStateMsg1(nameNodeRest.GetString() == "rest", "%s", nameNodeRest)

	nameNodeObj := destructPat.GetNext()
	jutil.CheckStateMsg1(nameNodeObj.GetString() == "obj", "%s", nameNodeObj)

	nuh_assertNotLhsByDestructuring(t, restNode)
	nuh_assertLhsByDestructuring(t, nameNodeRest)
	nuh_assertNotLhsByDestructuring(t, nameNodeObj)

	nuAssertThat(t, NodeUtil_GetRootTarget(destructPat)).IsSameInstanceAs(destructPat)
	nuAssertThat(t, NodeUtil_GetRootTarget(nameNodeRest)).IsSameInstanceAs(destructPat)

	nuAssertThat(t, NodeUtil_GetDeclaringParent(nameNodeRest)).IsSameInstanceAs(varNode)
}

func nuh_assertLhsByDestructuring(t *testing.T, n rhino.Node) {
	nuAssertThat(t, NodeUtil_IsLhsByDestructuring(n)).IsTrue()
}

func nuh_assertNotLhsByDestructuring(t *testing.T, n rhino.Node) {
	nuAssertThat(t, NodeUtil_IsLhsByDestructuring(n)).IsFalse()
}

func TestNodeUtil_AssortedTests_TestNewQName1(t *testing.T) {
	compiler := NewCompiler()
	options := NewCompilerOptions()
	options.SetCodingConvention(NewGoogleCodingConvention())
	compiler.Init([]*SourceFile{}, []*SourceFile{}, options)
	actual := NodeUtil_NewQName(compiler, "ns.prop")
	expected := rhino.IR_Getprop(rhino.IR_Name("ns"), "prop")
	nuAssertNode(t, actual).IsEqualTo(expected)
}

func TestNodeUtil_AssortedTests_TestNewQualifiedNameNode2(t *testing.T) {
	compiler := NewCompiler()
	options := NewCompilerOptions()
	options.SetCodingConvention(NewGoogleCodingConvention())
	compiler.Init([]*SourceFile{}, []*SourceFile{}, options)
	actual := NodeUtil_NewQName(compiler, "this.prop")
	expected := rhino.IR_Getprop(rhino.IR_ThisNode(), "prop")
	nuAssertNode(t, actual).IsEqualTo(expected)
}

func TestNodeUtil_AssortedTests_TestNewQNameDeclarationWithQualifiedName(t *testing.T) {
	nuAssertNode(t, nuh_createNewQNameDeclaration(t, "ns.prop", rhino.IR_Number(0), rhino.Token_VAR)).IsEqualTo(rhino.IR_ExprResult(rhino.IR_Assign(rhino.IR_Getprop(rhino.IR_Name("ns"), "prop"), rhino.IR_Number(0))))
}

func TestNodeUtil_AssortedTests_TestNewQNameDeclarationWithVar(t *testing.T) {
	nuAssertNode(t, nuh_createNewQNameDeclaration(t, "x", rhino.IR_Number(0), rhino.Token_VAR)).IsEqualTo(rhino.IR_VarWithValue(rhino.IR_Name("x"), rhino.IR_Number(0)))
}

func TestNodeUtil_AssortedTests_TestNewQNameDeclarationWithLet(t *testing.T) {
	nuAssertNode(t, nuh_createNewQNameDeclaration(t, "x", rhino.IR_Number(0), rhino.Token_LET)).IsEqualTo(rhino.IR_LetWithValue(rhino.IR_Name("x"), rhino.IR_Number(0)))
}

func TestNodeUtil_AssortedTests_TestNewQNameDeclarationWithConst(t *testing.T) {
	nuAssertNode(t, nuh_createNewQNameDeclaration(t, "x", rhino.IR_Number(0), rhino.Token_CONST)).IsEqualTo(rhino.IR_ConstNode(rhino.IR_Name("x"), rhino.IR_Number(0)))
}

func nuh_createNewQNameDeclaration(t *testing.T, name string, value rhino.Node, type_ rhino.Token) rhino.Node {
	compiler := NewCompiler()
	options := NewCompilerOptions()
	options.SetCodingConvention(NewGoogleCodingConvention())
	compiler.Init([]*SourceFile{}, []*SourceFile{}, options)
	return NodeUtil_NewQNameDeclarationWithType(compiler, name, value, nil, type_)
}

func TestNodeUtil_AssortedTests_TestGetBestJsDocInfoForClasses(t *testing.T) {
	classNode := nuParseFirst(t, rhino.Token_CLASS, "/** @export */ class Foo {}")
	nuAssertThat(t, NodeUtil_GetBestJSDocInfo(classNode).IsExport()).IsTrue()

	classNode = nuParseFirst(t, rhino.Token_CLASS, "/** @export */ var Foo = class {}")
	nuAssertThat(t, NodeUtil_GetBestJSDocInfo(classNode).IsExport()).IsTrue()

	classNode = nuParseFirst(t, rhino.Token_CLASS, "/** @export */ var Foo = class Bar {}")
	nuAssertThat(t, NodeUtil_GetBestJSDocInfo(classNode).IsExport()).IsTrue()

	classNode = nuParseFirst(t, rhino.Token_CLASS, "var Foo_1; /** @export */ let Foo = Foo_1 = class Foo {}")
	nuAssertThat(t, NodeUtil_GetBestJSDocInfo(classNode).IsExport()).IsTrue()
}

func TestNodeUtil_AssortedTests_TestGetBestJsDocInfoForMethods(t *testing.T) {
	function := nuParseFirst(t, rhino.Token_FUNCTION, "class C { /** @export */ foo() {} }")
	nuAssertThat(t, NodeUtil_GetBestJSDocInfo(function).IsExport()).IsTrue()

	function = nuParseFirst(t, rhino.Token_FUNCTION, "class C { /** @export */ [computedMethod]() {} }")
	nuAssertThat(t, NodeUtil_GetBestJSDocInfo(function).IsExport()).IsTrue()
}

func TestNodeUtil_AssortedTests_TestGetBestJsDocInfoExport(t *testing.T) {
	classNode := nuParseFirst(t, rhino.Token_CLASS, "/** @constructor */ export class Foo {}")
	nuAssertThat(t, NodeUtil_GetBestJSDocInfo(classNode).IsConstructor()).IsTrue()

	function := nuParseFirst(t, rhino.Token_FUNCTION, "/** @constructor */ export function Foo() {}")
	nuAssertThat(t, NodeUtil_GetBestJSDocInfo(function).IsConstructor()).IsTrue()

	function = nuParseFirst(t, rhino.Token_FUNCTION, "/** @constructor */ export var Foo = function() {}")
	nuAssertThat(t, NodeUtil_GetBestJSDocInfo(function).IsConstructor()).IsTrue()

	function = nuParseFirst(t, rhino.Token_FUNCTION, "/** @constructor */ export let Foo = function() {}")
	nuAssertThat(t, NodeUtil_GetBestJSDocInfo(function).IsConstructor()).IsTrue()
}

func TestNodeUtil_AssortedTests_TestGetDeclaredTypeExpression1(t *testing.T) {
	ast := nuParse(t, "function f(/** string */ x) {}")
	x := nuh_getNameNode(t, ast, "x")
	typeExpr := NodeUtil_GetDeclaredTypeExpression(x)
	nuAssertThat(t, typeExpr.GetRoot().GetString()).IsEqualTo("string")
}

func TestNodeUtil_AssortedTests_TestGetDeclaredTypeExpression2(t *testing.T) {
	ast := nuParse(t, "/** @param {string} x */ function f(x) {}")
	x := nuh_getNameNode(t, ast, "x")
	typeExpr := NodeUtil_GetDeclaredTypeExpression(x)
	nuAssertThat(t, typeExpr.GetRoot().GetString()).IsEqualTo("string")
}

func TestNodeUtil_AssortedTests_TestGetDeclaredTypeExpression3(t *testing.T) {
	ast := nuParse(t, "/** @param {...number} x */ function f(...x) {}")
	x := nuh_getNameNode(t, ast, "x")
	typeExpr := NodeUtil_GetDeclaredTypeExpression(x)
	nuAssertNode(t, typeExpr.GetRoot()).HasType(rhino.Token_ITER_REST)
	nuAssertThat(t, typeExpr.GetRoot().GetFirstChild().GetString()).IsEqualTo("number")
}

func TestNodeUtil_AssortedTests_TestGetDeclaredTypeExpression4(t *testing.T) {
	ast := nuParse(t, "/** @param {number=} x */ function f(x = -1) {}")
	x := nuh_getNameNode(t, ast, "x")
	typeExpr := NodeUtil_GetDeclaredTypeExpression(x)
	nuAssertNode(t, typeExpr.GetRoot()).HasType(rhino.Token_EQUALS)
	nuAssertThat(t, typeExpr.GetRoot().GetFirstChild().GetString()).IsEqualTo("number")
}

func TestNodeUtil_AssortedTests_TestFindLhsNodesInNodeWithNameDeclaration(t *testing.T) {
	nuAssertThat(t, nuh_visitLhsNodesInNode(t, "var x;")).HasSize(1)
	nuAssertThat(t, nuh_visitLhsNodesInNode(t, "var x, y;")).HasSize(2)
	nuAssertThat(t, nuh_visitLhsNodesInNode(t, "var f = function(x, y, z) {};")).HasSize(1)
}

func TestNodeUtil_AssortedTests_TestFindLhsNodesInNodeWithArrayPatternDeclaration(t *testing.T) {
	nuAssertThat(t, nuh_visitLhsNodesInNode(t, "var [x=a => a, y = b=>b+1] = arr;")).HasSize(2)
	nuAssertThat(t, nuh_visitLhsNodesInNode(t, "var [x=a => a, y = b=>b+1, ...z] = arr;")).HasSize(3)
	nuAssertThat(t, nuh_visitLhsNodesInNode(t, "var [ , , , y = b=>b+1, ...z] = arr;")).HasSize(2)
}

func TestNodeUtil_AssortedTests_TestFindLhsNodesInNodeWithObjectPatternDeclaration(t *testing.T) {
	nuAssertThat(t, nuh_visitLhsNodesInNode(t, "var {x = a=>a, y = b=>b+1} = obj;")).HasSize(2)
	nuAssertThat(t, nuh_visitLhsNodesInNode(t, "var {p1: x = a=>a, p2: y = b=>b+1} = obj;")).HasSize(2)
	nuAssertThat(t, nuh_visitLhsNodesInNode(t, "var {[pname]: x = a=>a, [p2name]: y} = obj;")).HasSize(2)
	nuAssertThat(t, nuh_visitLhsNodesInNode(t, "var {lhs1 = a, p2: [lhs2, lhs3 = b] = [notlhs]} = obj;")).HasSize(3)
}

func TestNodeUtil_AssortedTests_TestFindLhsNodesInNodeWithCastOnLhs(t *testing.T) {
	lhsNodes := nuh_visitLhsNodesInNode(t, "/** @type {*} */ (a.b) = 3;")
	nuAssertThat(t, lhsNodes).HasSize(1)
	nodeIterator := nuIteratorOf(lhsNodes)
	nuAssertNode(t, nodeIterator.Next()).MatchesQualifiedName("a.b")
}

func TestNodeUtil_AssortedTests_TestFindLhsNodesInNodeWithArrayPatternAssign(t *testing.T) {
	nuAssertThat(t, nuh_visitLhsNodesInNode(t, "[this.x] = rhs;")).HasSize(1)
	nuAssertThat(t, nuh_visitLhsNodesInNode(t, "[this.x, y] = rhs;")).HasSize(2)
	nuAssertThat(t, nuh_visitLhsNodesInNode(t, "[this.x, y, this.z] = rhs;")).HasSize(3)
	nuAssertThat(t, nuh_visitLhsNodesInNode(t, "[y, this.z] = rhs;")).HasSize(2)
	nuAssertThat(t, nuh_visitLhsNodesInNode(t, "[x[y]] = rhs;")).HasSize(1)
	nuAssertThat(t, nuh_visitLhsNodesInNode(t, "[x.y.z] = rhs;")).HasSize(1)
	nuAssertThat(t, nuh_visitLhsNodesInNode(t, "[ /** @type {*} */ (x.y.z) ] = rhs;")).HasSize(1)
}

func TestNodeUtil_AssortedTests_TestFindLhsNodesInNodeWithComplexAssign(t *testing.T) {
	nuAssertThat(t, nuh_visitLhsNodesInNode(t, "x += 1;")).HasSize(1)
	nuAssertThat(t, nuh_visitLhsNodesInNode(t, "x.y += 1;")).HasSize(1)
	nuAssertThat(t, nuh_visitLhsNodesInNode(t, "x -= 1;")).HasSize(1)
	nuAssertThat(t, nuh_visitLhsNodesInNode(t, "x.y -= 1;")).HasSize(1)
	nuAssertThat(t, nuh_visitLhsNodesInNode(t, "x *= 2;")).HasSize(1)
	nuAssertThat(t, nuh_visitLhsNodesInNode(t, "x.y *= 2;")).HasSize(1)
}

func TestNodeUtil_AssortedTests_TestFindLhsNodesInForOfWithDeclaration(t *testing.T) {
	lhsNodes := nuh_visitLhsNodesInNode(t, "for (const {x, y} of iterable) {}")
	nuAssertThat(t, lhsNodes).HasSize(2)
	nodeIterator := nuIteratorOf(lhsNodes)
	nuAssertNode(t, nodeIterator.Next()).IsName("x")
	nuAssertNode(t, nodeIterator.Next()).IsName("y")
}

func TestNodeUtil_AssortedTests_TestFindLhsNodesInForOfWithoutDeclaration(t *testing.T) {
	lhsNodes := nuh_visitLhsNodesInNode(t, "for ({x, y: a.b} of iterable) {}")
	nuAssertThat(t, lhsNodes).HasSize(2)
	nodeIterator := nuIteratorOf(lhsNodes)
	nuAssertNode(t, nodeIterator.Next()).IsName("x")
	nuAssertNode(t, nodeIterator.Next()).MatchesQualifiedName("a.b")
}

func TestNodeUtil_AssortedTests_TestFindLhsNodesInForInWithDeclaration(t *testing.T) {
	lhsNodes := nuh_visitLhsNodesInNode(t, "for (const x in obj) {}")
	nuAssertThat(t, lhsNodes).HasSize(1)
	nodeIterator := nuIteratorOf(lhsNodes)
	nuAssertNode(t, nodeIterator.Next()).IsName("x")
}

func TestNodeUtil_AssortedTests_TestFindLhsNodesInForInWithoutDeclaration(t *testing.T) {
	lhsNodes := nuh_visitLhsNodesInNode(t, "for (a.b in iterable) {}")
	nuAssertThat(t, lhsNodes).HasSize(1)
	nodeIterator := nuIteratorOf(lhsNodes)
	nuAssertNode(t, nodeIterator.Next()).MatchesQualifiedName("a.b")
}

func TestNodeUtil_AssortedTests_TestGetGoogRequireInfo_returnsNullForShadowedRequires(t *testing.T) {
	src := "goog.module('a.b.c'); const Foo = goog.require('d.Foo'); { const Foo = 0; }"

	first, compiler := nuParseWithCompiler(t, src)
	// Build the three layers of scopes: global scope, module scope, and block scope.
	scopeCreator := NewSyntacticScopeCreator(compiler)
	globalScope := scopeCreator.CreateScopeOfScope(first, nil)
	module := nuGetNode(t, first, rhino.Token_MODULE_BODY)
	moduleScope := scopeCreator.CreateScopeOfScope(module, globalScope)
	block := nuGetNode(t, first, rhino.Token_BLOCK)
	blockScope := scopeCreator.CreateScopeOfScope(block, moduleScope)

	nuAssertThat(t, NodeUtil_GetGoogRequireInfoWithScope("Foo", moduleScope)).IsEqualTo(NodeUtil_GoogRequire_FromNamespace("d.Foo", true))
	nuAssertThat(t, NodeUtil_GetGoogRequireInfoWithScope("Foo", blockScope)).IsNull()
}

func TestNodeUtil_AssortedTests_TestIsConstructor(t *testing.T) {
	nuAssertThat(t,
		NodeUtil_IsConstructor(nuParseFirst(t, rhino.Token_FUNCTION, "/** @constructor */ function Foo() {}"))).IsTrue()
	nuAssertThat(t,
		NodeUtil_IsConstructor(
			nuParseFirst(t, rhino.Token_FUNCTION, "/** @constructor */ var Foo = function() {}"))).IsTrue()
	nuAssertThat(t,
		NodeUtil_IsConstructor(
			nuParseFirst(t, rhino.Token_FUNCTION, "var x = {}; /** @constructor */ x.Foo = function() {}"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsConstructor(nuParseFirst(t, rhino.Token_FUNCTION, "class Foo { constructor() {} }"))).IsTrue()
	nuAssertThat(t,
		NodeUtil_IsConstructor(nuParseFirst(t, rhino.Token_FUNCTION, "class Foo { static constructor() {} }"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsConstructor(nuParseFirst(t, rhino.Token_FUNCTION, "let Foo = { constructor() {} }"))).IsFalse()

	nuAssertThat(t, NodeUtil_IsConstructor(nuParseFirst(t, rhino.Token_FUNCTION, "function Foo() {}"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsConstructor(nuParseFirst(t, rhino.Token_FUNCTION, "var Foo = function() {}"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsConstructor(nuParseFirst(t, rhino.Token_FUNCTION, "var x = {}; x.Foo = function() {};"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsConstructor(nuParseFirst(t, rhino.Token_FUNCTION, "function constructor() {}"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsConstructor(nuParseFirst(t, rhino.Token_FUNCTION, "class Foo { bar() {} }"))).IsFalse()
}

func TestNodeUtil_AssortedTests_TestGetEs6ClassConstructorMemberFunctionDef(t *testing.T) {
	// distinguish between static and non-static method named "constructor"
	constructorMemberFunctionDef :=
		NodeUtil_GetEs6ClassConstructorMemberFunctionDef(
			nuParseFirst(t,
				rhino.Token_CLASS, "class Foo { method() {} constructor() {} static constructor() {} };"))
	nuAssertNode(t, constructorMemberFunctionDef).IsMemberFunctionDef("constructor")
	nuAssertThat(t, constructorMemberFunctionDef.IsStaticMember()).IsFalse()

	// static method named "constructor" will not be returned as the constructor
	constructorMemberFunctionDef =
		NodeUtil_GetEs6ClassConstructorMemberFunctionDef(
			nuParseFirst(t, rhino.Token_CLASS, "class Foo { method() {} static constructor() {} }"))
	nuAssertThat(t, constructorMemberFunctionDef).IsNull()
}

func TestNodeUtil_AssortedTests_TestIsGetterOrSetter(t *testing.T) {
	fnNode :=
		nuParseFirst(t, rhino.Token_FUNCTION, "Object.defineProperty(this, 'bar', {get: function() {}});")
	nuAssertThat(t, NodeUtil_IsGetterOrSetter(fnNode.GetParent())).IsTrue()

	fnNode = nuParseFirst(t, rhino.Token_FUNCTION, "Object.defineProperty(this, 'bar', {set: function() {}});")
	nuAssertThat(t, NodeUtil_IsGetterOrSetter(fnNode.GetParent())).IsTrue()

	fnNode = nuParseFirst(t, rhino.Token_FUNCTION, "Object.defineProperties(this, {bar: {get: function() {}}});")
	nuAssertThat(t, NodeUtil_IsGetterOrSetter(fnNode.GetParent())).IsTrue()

	fnNode = nuParseFirst(t, rhino.Token_FUNCTION, "Object.defineProperties(this, {bar: {set: function() {}}});")
	nuAssertThat(t, NodeUtil_IsGetterOrSetter(fnNode.GetParent())).IsTrue()

	fnNode = nuParseFirst(t, rhino.Token_FUNCTION, "var x = {get bar() {}};")
	nuAssertThat(t, NodeUtil_IsGetterOrSetter(fnNode.GetParent())).IsTrue()

	fnNode = nuParseFirst(t, rhino.Token_FUNCTION, "var x = {set bar(z) {}};")
	nuAssertThat(t, NodeUtil_IsGetterOrSetter(fnNode.GetParent())).IsTrue()
}

func TestNodeUtil_AssortedTests_TestIsObjectDefinePropertiesDefinition(t *testing.T) {
	nuAssertThat(t,
		NodeUtil_IsObjectDefinePropertiesDefinition(
			nuParseFirst(t, rhino.Token_CALL, "Object.defineProperties(this, {});"))).IsTrue()
	nuAssertThat(t,
		NodeUtil_IsObjectDefinePropertiesDefinition(
			nuParseFirst(t, rhino.Token_CALL, "Object.defineProperties(this, foo);"))).IsTrue()
	nuAssertThat(t,
		NodeUtil_IsObjectDefinePropertiesDefinition(
			nuParseFirst(t, rhino.Token_CALL, "$jscomp.global.Object.defineProperties(this, foo);"))).IsTrue()
	nuAssertThat(t,
		NodeUtil_IsObjectDefinePropertiesDefinition(
			nuParseFirst(t, rhino.Token_CALL, "$jscomp$global.Object.defineProperties(this, foo);"))).IsTrue()

	nuAssertThat(t,
		NodeUtil_IsObjectDefinePropertiesDefinition(
			nuParseFirst(t, rhino.Token_CALL, "Object.defineProperties(this, {}, foo);"))).IsFalse()
	nuAssertThat(t,
		NodeUtil_IsObjectDefinePropertiesDefinition(
			nuParseFirst(t, rhino.Token_CALL, "Object.defineProperties(this);"))).IsFalse()
	nuAssertThat(t,
		NodeUtil_IsObjectDefinePropertiesDefinition(
			nuParseFirst(t, rhino.Token_CALL, "Object.defineProperties();"))).IsFalse()
}

func TestNodeUtil_AssortedTests_TestIsObjectCreateWithPropertiesDefinition(t *testing.T) {
	nuAssertThat(t,
		NodeUtil_IsObjectCreateWithPropertiesDefinition(
			nuParseFirst(t, rhino.Token_CALL, "Object.create(null, {});"))).IsTrue()
	nuAssertThat(t,
		NodeUtil_IsObjectCreateWithPropertiesDefinition(
			nuParseFirst(t, rhino.Token_CALL, "$jscomp.global.Object.create(proto, foo);"))).IsTrue()

	nuAssertThat(t,
		NodeUtil_IsObjectCreateWithPropertiesDefinition(
			nuParseFirst(t, rhino.Token_CALL, "Object.create(proto);"))).IsFalse()
	nuAssertThat(t,
		NodeUtil_IsObjectCreateWithPropertiesDefinition(
			nuParseFirst(t, rhino.Token_CALL, "other.Object.create(proto, {});"))).IsFalse()
}

func TestNodeUtil_AssortedTests_TestIsObjectDefinePropertyDefinition(t *testing.T) {
	nuAssertThat(t,
		NodeUtil_IsObjectDefinePropertyDefinition(
			nuParseFirst(t, rhino.Token_CALL, "Object.defineProperty(this, 'foo', {});"))).IsTrue()
	nuAssertThat(t,
		NodeUtil_IsObjectDefinePropertyDefinition(
			nuParseFirst(t, rhino.Token_CALL, "Object.defineProperty(this, 'foo', foo);"))).IsTrue()
	nuAssertThat(t,
		NodeUtil_IsObjectDefinePropertyDefinition(
			nuParseFirst(t, rhino.Token_CALL, "$jscomp.global.Object.defineProperty(this, 'foo', foo);"))).IsTrue()
	nuAssertThat(t,
		NodeUtil_IsObjectDefinePropertyDefinition(
			nuParseFirst(t, rhino.Token_CALL, "$jscomp$global.Object.defineProperty(this, 'foo', foo);"))).IsTrue()

	nuAssertThat(t,
		NodeUtil_IsObjectDefinePropertyDefinition(
			nuParseFirst(t, rhino.Token_CALL, "Object.defineProperty(this, {});"))).IsFalse()
	nuAssertThat(t,
		NodeUtil_IsObjectDefinePropertyDefinition(
			nuParseFirst(t, rhino.Token_CALL, "Object.defineProperty(this);"))).IsFalse()
	nuAssertThat(t,
		NodeUtil_IsObjectDefinePropertyDefinition(
			nuParseFirst(t, rhino.Token_CALL, "Object.defineProperty();"))).IsFalse()
}

func TestNodeUtil_AssortedTests_TestGetAllModuleVars(t *testing.T) {
	js :=
		"goog.module('m'); var h =2; function g(x, y) {var z; {let a; const b = 1} let c}"
	compiler := NewCompiler() // Java: new Compiler()

	compiler.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_NORMALIZED)
	scopeCreator := NewSyntacticScopeCreator(compiler)
	ast := nuParse(t, js)
	moduleNode := nuParseFirst(t, rhino.Token_MODULE_BODY, js)
	globalScope := Scope_CreateGlobalScope(ast)
	allVariables :=
		NodeUtil_GetAllVarNamesDeclaredInModule(moduleNode, compiler, scopeCreator, globalScope)
	nuAssertThat(t, allVariables).ContainsExactly("g", "h")
}

func TestNodeUtil_AssortedTests_TestGetAllModuleVars2(t *testing.T) {
	js := "var glob = 3;"
	compiler := NewCompiler() // Java: new Compiler()

	compiler.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_NORMALIZED)
	scopeCreator := NewSyntacticScopeCreator(compiler)
	ast := nuParse(t, js)
	globalScope := Scope_CreateGlobalScope(ast)
	ex :=
		assertThrows[jutil.IIllegalStateException](t, func() {
			NodeUtil_GetAllVarNamesDeclaredInModule(
				ast, compiler, scopeCreator, globalScope)
		})

	nuAssertThat(t, ex.GetMessage()).IsEqualTo("getAllVarsDeclaredInModule expects a module body node")
}

func TestNodeUtil_AssortedTests_TestGetAllVars1(t *testing.T) {
	fnString := "var h; function g(x, y) {var z; h = 2; {let a; const b = 1} let c}"
	compiler := NewCompiler() // Java: new Compiler()
	compiler.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_NORMALIZED)
	scopeCreator := NewSyntacticScopeCreator(compiler)

	ast := nuParse(t, fnString)
	functionNode := nuParseFirst(t, rhino.Token_FUNCTION, fnString)

	globalScope := Scope_CreateGlobalScope(ast)
	functionScope := scopeCreator.CreateScopeOfScope(functionNode, globalScope)

	allVarsDeclaredInFunction :=
		NodeUtil_GetAllVarsDeclaredInFunction(compiler, scopeCreator, functionScope)
	allVariables := allVarsDeclaredInFunction.GetAllVariables()

	nuAssertThat(t, allVariables.KeySlice()).ContainsExactly("a", "b", "c", "z", "x", "y")
}

func TestNodeUtil_AssortedTests_TestGetAllVars2(t *testing.T) {
	fnString :=
		"function g(x, y)\n{var z;\n{let a = (no1, no2) => { let no6, no7; };\nconst b = 1}\nlet c}\nfunction u(h) {let e}\n"

	compiler := NewCompiler() // Java: new Compiler()
	compiler.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_NORMALIZED)
	scopeCreator := NewSyntacticScopeCreator(compiler)

	ast := nuParse(t, fnString)
	functionNode := nuParseFirst(t, rhino.Token_FUNCTION, fnString)

	globalScope := Scope_CreateGlobalScope(ast)
	functionScope := scopeCreator.CreateScopeOfScope(functionNode, globalScope)

	allVarsDeclaredInFunction :=
		NodeUtil_GetAllVarsDeclaredInFunction(compiler, scopeCreator, functionScope)
	allVariables := allVarsDeclaredInFunction.GetAllVariables()

	nuAssertThat(t, allVariables.KeySlice()).ContainsExactly("x", "y", "z", "a", "b", "c")
}

func TestNodeUtil_AssortedTests_TestDoesFunctionReferenceOwnArgumentsObject(t *testing.T) {
	nuAssertThat(t,
		NodeUtil_DoesFunctionReferenceOwnArgumentsObject(
			nuParseExpr(t, "function() {return () => arguments}"))).IsTrue()

	nuAssertThat(t,
		NodeUtil_DoesFunctionReferenceOwnArgumentsObject(
			nuParseExpr(t, "function(a = arguments.length) {}"))).IsTrue()

	nuAssertThat(t, NodeUtil_DoesFunctionReferenceOwnArgumentsObject(nuParseExpr(t, "() => arguments"))).IsFalse()
}

/**
 * When the left side is a destructuring pattern, generally it's not possible to identify the
 * RHS for a specific name on the LHS.
 */
func TestNodeUtil_AssortedTests_TestIsExpressionResultUsed(t *testing.T) {
	nuAssertThat(t, NodeUtil_IsExpressionResultUsed(nuh_getNameNodeFrom(t, "for (x in y) z", "x"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsExpressionResultUsed(nuh_getNameNodeFrom(t, "for (x in y) z", "y"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsExpressionResultUsed(nuh_getNameNodeFrom(t, "for (x in y) z", "z"))).IsFalse()

	nuAssertThat(t, NodeUtil_IsExpressionResultUsed(nuh_getNameNodeFrom(t, "for (x of y) z", "x"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsExpressionResultUsed(nuh_getNameNodeFrom(t, "for (x of y) z", "y"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsExpressionResultUsed(nuh_getNameNodeFrom(t, "for (x of y) z", "z"))).IsFalse()

	nuAssertThat(t, NodeUtil_IsExpressionResultUsed(nuh_getNameNodeFrom(t, "for (x; y; z) a", "x"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsExpressionResultUsed(nuh_getNameNodeFrom(t, "for (x; y; z) a", "y"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsExpressionResultUsed(nuh_getNameNodeFrom(t, "for (x; y; z) a", "z"))).IsFalse()

	nuAssertThat(t, NodeUtil_IsExpressionResultUsed(nuh_getNameNodeFrom(t, "function f() { return y }", "y"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsExpressionResultUsed(nuh_getNameNodeFrom(t, "function f() { y }", "y"))).IsFalse()

	nuAssertThat(t, NodeUtil_IsExpressionResultUsed(nuh_getNameNodeFrom(t, "var x = ()=> y", "y"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsExpressionResultUsed(nuh_getNameNodeFrom(t, "var x = ()=>{ y }", "y"))).IsFalse()

	nuAssertThat(t, NodeUtil_IsExpressionResultUsed(nuh_getNameNodeFrom(t, "({a: x = y} = z)", "y"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsExpressionResultUsed(nuh_getNameNodeFrom(t, "[x = y] = z", "y"))).IsTrue()

	nuAssertThat(t, NodeUtil_IsExpressionResultUsed(nuh_getNameNodeFrom(t, "({x: y})", "y"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsExpressionResultUsed(nuh_getNameNodeFrom(t, "[y]", "y"))).IsTrue()

	nuAssertThat(t, NodeUtil_IsExpressionResultUsed(nuh_getNameNodeFrom(t, "y()", "y"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsExpressionResultUsed(nuh_getNameNodeFrom(t, "y``", "y"))).IsTrue()

	// We're using string literals instead of numbers here, because IRFactory will automatically
	// remove a number in an indirect call like `(0, some.callee)()` after applying FREE_CALL
	// to the `CALL` node.
	nuAssertThat(t, NodeUtil_IsExpressionResultUsed(nuh_getStringLitNodeFrom(t, "('',eval)()", ""))).IsTrue()
	nuAssertThat(t, NodeUtil_IsExpressionResultUsed(nuh_getStringLitNodeFrom(t, "('',x.y)()", ""))).IsTrue()
	nuAssertThat(t, NodeUtil_IsExpressionResultUsed(nuh_getStringLitNodeFrom(t, "('',x.y)``", ""))).IsTrue()
}

func TestNodeUtil_AssortedTests_TestIsSimpleOperator(t *testing.T) {
	nuAssertThat(t, NodeUtil_IsSimpleOperator(nuParseExpr(t, "!x"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsSimpleOperator(nuParseExpr(t, "5 + x"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsSimpleOperator(nuParseExpr(t, "typeof x"))).IsTrue()
	nuAssertThat(t, NodeUtil_IsSimpleOperator(nuParseExpr(t, "x instanceof y"))).IsTrue()
	// short curcuits aren't simple
	nuAssertThat(t, NodeUtil_IsSimpleOperator(nuParseExpr(t, "5 && x"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsSimpleOperator(nuParseExpr(t, "5 || x"))).IsFalse()
	// side-effects aren't simple
	nuAssertThat(t, NodeUtil_IsSimpleOperator(nuParseExpr(t, "x = 5"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsSimpleOperator(nuParseExpr(t, "x++"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsSimpleOperator(nuParseExpr(t, "--y"))).IsFalse()
	// prop access with IN are simple
	nuAssertThat(t, NodeUtil_IsSimpleOperator(nuParseExpr(t, "y in x"))).IsTrue()
	// prop access have side effects because of getters/setters
	nuAssertThat(t, NodeUtil_IsSimpleOperator(nuParseExpr(t, "x.y"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsSimpleOperator(nuParseExpr(t, "x[y]"))).IsFalse()
}

func TestNodeUtil_AssortedTests_TestIsCallToString(t *testing.T) {
	nuAssertThat(t, NodeUtil_IsCallTo(nuParseFirst(t, rhino.Token_CALL, "foo()"), "foo")).IsTrue()
	nuAssertThat(t, NodeUtil_IsCallTo(nuParseFirst(t, rhino.Token_CALL, "foo.bar()"), "foo.bar")).IsTrue()

	nuAssertThat(t, NodeUtil_IsCallTo(rhino.IR_Name("foo"), "foo")).IsFalse()
	nuAssertThat(t, NodeUtil_IsCallTo(rhino.IR_Getprop(rhino.IR_Name("foo"), "bar"), "foo.bar")).IsFalse()
	nuAssertThat(t, NodeUtil_IsCallTo(nuParseFirst(t, rhino.Token_CALL, "foo.bar()"), "foo")).IsFalse()
	nuAssertThat(t, NodeUtil_IsCallTo(nuParseFirst(t, rhino.Token_CALL, "foo[0]()"), "foo")).IsFalse()
}

func TestNodeUtil_AssortedTests_TestIsCallToNode(t *testing.T) {
	nuAssertThat(t, NodeUtil_IsCallToNode(nuParseFirst(t, rhino.Token_CALL, "foo()"), rhino.IR_Name("foo"))).IsTrue()
	nuAssertThat(t,
		NodeUtil_IsCallToNode(nuParseFirst(t, rhino.Token_CALL, "foo.bar()"), rhino.IR_Getprop(rhino.IR_Name("foo"), "bar"))).IsTrue()

	nuAssertThat(t, NodeUtil_IsCallToNode(rhino.IR_Name("foo"), rhino.IR_Name("foo"))).IsFalse()
	nuAssertThat(t,
		NodeUtil_IsCallToNode(
			rhino.IR_Getprop(rhino.IR_Name("foo"), "bar"), rhino.IR_Getprop(rhino.IR_Name("foo"), "bar"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsCallToNode(nuParseFirst(t, rhino.Token_CALL, "foo.bar()"), rhino.IR_Name("foo"))).IsFalse()
	nuAssertThat(t, NodeUtil_IsCallToNode(nuParseFirst(t, rhino.Token_CALL, "foo[0]()"), rhino.IR_Name("foo"))).IsFalse()
}

func TestNodeUtil_AssortedTests_TestCallTargetNodeResolution_qualifiedName(t *testing.T) {
	nuAssertThat(t,
		NodeUtil_GetCallTargetResolvingIndirectCalls(nuParse(t, "foo()").GetFirstFirstChild()).MatchesQualifiedName("foo")).IsTrue() // matches `foo`

	nuAssertThat(t,
		NodeUtil_GetCallTargetResolvingIndirectCalls(
			nuParse(t, "a.b.c.foo()").GetFirstFirstChild()).MatchesQualifiedName("a.b.c.foo")).IsTrue()
}

func TestNodeUtil_AssortedTests_TestCallTargetNodeResolution_notAQualifiedName(t *testing.T) {
	script := nuParse(t, "({valueOf: () => {}}.valueOf());")
	nuAssertThat(t, script.IsScript()).IsTrue()
	callNode := script.GetFirstFirstChild()
	nuAssertThat(t, callNode.IsCall()).IsTrue()
	nuAssertThat(t, NodeUtil_GetCallTargetResolvingIndirectCalls(callNode).IsQualifiedName()).IsFalse()
}

func TestNodeUtil_AssortedTests_TestIsGoogModuleGetCall(t *testing.T) {
	root := nuParse(t, "const Foo = goog.module.get('a.b.c.Foo');")
	call := root.GetFirstChild().GetFirstChild().GetFirstChild()

	nuAssertThat(t, NodeUtil_IsGoogModuleGetCall(call)).IsTrue()

	root = nuParse(t, "const Foo = goog.require('a.b.c.Foo');")
	call = root.GetFirstChild().GetFirstChild().GetFirstChild()

	nuAssertThat(t, NodeUtil_IsGoogModuleGetCall(call)).IsFalse()
}

func TestNodeUtil_AssortedTests_TestIsBundledGoogModule(t *testing.T) {
	googLoadModule :=
		nuParse(t, "goog.loadModule(function(exports) { goog.module('a.b'); return exports; });")

	googName := nuh_getNameNode(t, googLoadModule, "goog")
	callNode := googName.GetGrandparent()

	nuAssertNode(t, callNode).HasToken(rhino.Token_CALL)
	nuAssertThat(t, NodeUtil_IsBundledGoogModuleCall(callNode)).IsTrue()
}

func TestNodeUtil_AssortedTests_TestIsBundledGoogModule_onlyIfInScript(t *testing.T) {
	callNode :=
		rhino.IR_Call(
			rhino.IR_Getprop(rhino.IR_Name("goog"), "loadModule"), rhino.IR_String("imaginary module text here"))

	nuAssertThat(t, NodeUtil_IsBundledGoogModuleCall(callNode)).IsFalse()

	exprResult := rhino.IR_ExprResult(callNode)
	nuAssertThat(t, NodeUtil_IsBundledGoogModuleCall(callNode)).IsFalse()

	rhino.IR_ScriptOfNodes(exprResult) // Call this for its side effect of modifying `exprResult`.
	nuAssertThat(t, NodeUtil_IsBundledGoogModuleCall(callNode)).IsTrue()
}

func TestNodeUtil_AssortedTests_TestGetOriginalName_string(t *testing.T) {
	// No rename artifacts
	nuAssertThat(t, NodeUtil_GetOriginalName("foo")).IsEqualTo("foo")

	// Unique name suffix
	nuAssertThat(t, NodeUtil_GetOriginalName("foo$$1")).IsEqualTo("foo")
	nuAssertThat(t, NodeUtil_GetOriginalName("foo$$123")).IsEqualTo("foo")
	nuAssertThat(t, NodeUtil_GetOriginalName("foo_bar$$1")).IsEqualTo("foo_bar")

	// Contextual rename suffix
	nuAssertThat(t, NodeUtil_GetOriginalName("foo$jscomp$1")).IsEqualTo("foo")
	nuAssertThat(t, NodeUtil_GetOriginalName("foo_bar$jscomp$123")).IsEqualTo("foo_bar")

	// Module exports prefix
	nuAssertThat(t, NodeUtil_GetOriginalName("module$exports$foo$Bar")).IsEqualTo("Bar")
	nuAssertThat(t, NodeUtil_GetOriginalName("module$exports$foo$bar$Baz")).IsEqualTo("Baz")
	nuAssertThat(t, NodeUtil_GetOriginalName("module$exports$foo$bar$PREVENT_CSN_PROPERTY")).
		IsEqualTo("PREVENT_CSN_PROPERTY")

	// Module contents prefix (without underscores in module name)
	nuAssertThat(t, NodeUtil_GetOriginalName("module$contents$foo$Bar_baz")).IsEqualTo("baz")
	nuAssertThat(t, NodeUtil_GetOriginalName("module$contents$foo$Bar_baz_qux")).IsEqualTo("qux")

	// Suffix and prefix combined
	nuAssertThat(t, NodeUtil_GetOriginalName("module$exports$foo$Bar$$1")).IsEqualTo("Bar")
	nuAssertThat(t, NodeUtil_GetOriginalName("module$contents$foo$Bar_baz$$123")).IsEqualTo("baz")
	nuAssertThat(t, NodeUtil_GetOriginalName("module$exports$foo$Bar$jscomp$1$$1")).IsEqualTo("Bar")
}

// ==== ConvertToNonOptChainTests public static class ConvertToNonOptChainTests {
func TestNodeUtil_ConvertToNonOptChainTests_SimpleChain(t *testing.T) {
	// `expr?.prop`
	origChain := rhino.IR_StartOptChainGetprop(rhino.IR_Name("expr"), "prop")
	NodeUtil_ConvertToNonOptionalChainSegment(origChain)
	nuAssertThat(t, nuh_isChainConverted(t, origChain)).IsTrue()
}

func TestNodeUtil_ConvertToNonOptChainTests_ContinuedChain(t *testing.T) {
	// `expr?.prop1.prop2`
	innerGetProp := rhino.IR_StartOptChainGetprop(rhino.IR_Name("expr"), "pro1")
	outterGetProp := rhino.IR_ContinueOptChainGetprop(innerGetProp, "prop2")
	NodeUtil_ConvertToNonOptionalChainSegment(outterGetProp)
	nuAssertThat(t, nuh_isChainConverted(t, outterGetProp)).IsTrue()
}

func TestNodeUtil_ConvertToNonOptChainTests_NestedChain(t *testing.T) {
	// `expr2[expr1 ?.prop1]?.prop2`
	innerGetProp := rhino.IR_StartOptChainGetprop(rhino.IR_Name("expr1"), "pro1")
	getElem := rhino.IR_Getelem(rhino.IR_Name("expr2"), innerGetProp)
	outterGetProp := rhino.IR_StartOptChainGetprop(getElem, "prop2")

	NodeUtil_ConvertToNonOptionalChainSegment(outterGetProp)
	nuAssertThat(t, nuh_isChainConverted(t, outterGetProp)).IsTrue()

	// nested (inner) optional chain is unchanged
	nuAssertThat(t, outterGetProp.GetFirstChild().GetSecondChild().IsOptChainGetProp()).IsTrue()
}

func TestNodeUtil_ConvertToNonOptChainTests_ConvertSegmentDownToStopNode(t *testing.T) {
	// `expr?.prop1.prop2.prop3`
	innerGetProp := rhino.IR_StartOptChainGetprop(rhino.IR_Name("expr"), "prop1")
	middleGetProp := rhino.IR_ContinueOptChainGetprop(innerGetProp, "prop2")
	outerGetProp := rhino.IR_ContinueOptChainGetprop(middleGetProp, "prop3")

	NodeUtil_ConvertToNonOptionalChainSegmentDownTo(outerGetProp, innerGetProp)
	// outer and middle are converted
	nuAssertThat(t, outerGetProp.IsGetProp()).IsTrue()
	nuAssertThat(t, middleGetProp.IsGetProp()).IsTrue()
	// stopNode (innerGetProp) remains untouched
	nuAssertThat(t, innerGetProp.IsOptChainGetProp()).IsTrue()
}

// All nodes till the current chain's start are converted to non-optional.
func nuh_isChainConverted(t *testing.T, node rhino.Node) bool {
	if node.IsNil() {
		return true
	}
	return !NodeUtil_IsOptChainNode(node) && nuh_isChainConverted(t, node.GetFirstChild())
}

// ==== GetStartOfOptChainTests public static class GetStartOfOptChainTests {

func TestNodeUtil_GetStartOfOptChainTests_IsStartOfChain(t *testing.T) {
	// `expr?.prop`
	optChainGet := rhino.IR_StartOptChainGetprop(rhino.IR_Name("expr"), "prop")

	nuAssertThat(t, NodeUtil_GetStartOfOptChainSegment(optChainGet)).IsEqualTo(optChainGet)
}

func TestNodeUtil_GetStartOfOptChainTests_ShortChain(t *testing.T) {
	// `expr?.prop1.prop2`
	innerGetProp := rhino.IR_StartOptChainGetprop(rhino.IR_Name("expr"), "pro1")
	outterGetProp := rhino.IR_ContinueOptChainGetprop(innerGetProp, "prop2")

	nuAssertThat(t, NodeUtil_GetStartOfOptChainSegment(outterGetProp)).IsEqualTo(innerGetProp)
}

func TestNodeUtil_GetStartOfOptChainTests_MixedChain(t *testing.T) {
	// `expr().prop1?.prop2()[prop3]`
	call := rhino.IR_Call(rhino.IR_Name("expr"))
	getProp := rhino.IR_Getprop(call, "prop1")
	optChainGetProp := rhino.IR_StartOptChainGetprop(getProp, "prop2")
	optChainCall := rhino.IR_ContinueOptChainCall(optChainGetProp)
	optChainGetElem := rhino.IR_ContinueOptChainGetelem(optChainCall, rhino.IR_Name("prop3"))

	nuAssertThat(t, NodeUtil_GetStartOfOptChainSegment(optChainGetElem)).IsEqualTo(optChainGetProp)
	nuAssertThat(t, NodeUtil_GetStartOfOptChainSegment(optChainCall)).IsEqualTo(optChainGetProp)
}

// ==== GetEndOfOptChainTests public static class GetEndOfOptChainTests {

func TestNodeUtil_GetEndOfOptChainTests_IsEndOfChain(t *testing.T) {
	// `expr?.prop`
	optChainGet := rhino.IR_StartOptChainGetprop(rhino.IR_Name("expr"), "prop")

	nuAssertThat(t, NodeUtil_GetEndOfOptChainSegment(optChainGet)).IsEqualTo(optChainGet)
}

func TestNodeUtil_GetEndOfOptChainTests_IsEndOfChain_call_innerChain(t *testing.T) {
	optChainCall := nuParseExpr(t, "a?.b?.(x?.y)")
	nuAssertThat(t, optChainCall.IsOptChainCall()).IsTrue()
	nuAssertThat(t, optChainCall.IsOptionalChainStart()).IsTrue()

	innerOptChain := optChainCall.GetLastChild() // `x?.y`
	nuAssertThat(t, innerOptChain.IsOptChainGetProp()).IsTrue()
	nuAssertThat(t, innerOptChain.IsOptionalChainStart()).IsTrue()
	nuAssertThat(t, NodeUtil_IsEndOfOptChainSegment(innerOptChain)).IsTrue()
	nuAssertThat(t, NodeUtil_IsEndOfFullOptChain(innerOptChain)).IsTrue()
}

func TestNodeUtil_GetEndOfOptChainTests_IsEndOfChain_getProp_innerChain(t *testing.T) {
	optChainGetProp := nuParseExpr(t, "a?.b.x?.y")
	nuAssertThat(t, optChainGetProp.IsOptChainGetProp()).IsTrue()
	nuAssertThat(t, optChainGetProp.IsOptionalChainStart()).IsTrue()
	nuAssertThat(t, NodeUtil_IsEndOfOptChainSegment(optChainGetProp)).IsTrue()
	nuAssertThat(t, NodeUtil_IsEndOfFullOptChain(optChainGetProp)).IsTrue()

	// Check that `a?.b.x` is not the start of the optional chain, but it ends the segment that
	// starts with `a?.b`.
	innerOptChain := optChainGetProp.GetFirstChild() // `a?.b.x`
	nuAssertThat(t, innerOptChain.IsOptChainGetProp()).IsTrue()
	nuAssertThat(t, innerOptChain.IsOptionalChainStart()).IsFalse()
	nuAssertThat(t, NodeUtil_IsEndOfOptChainSegment(innerOptChain)).IsTrue()
	nuAssertThat(t, NodeUtil_IsEndOfFullOptChain(innerOptChain)).IsFalse()

	// Check that `a?.b` is the start of optChain but not the end
	innerMostOptChain := innerOptChain.GetFirstChild() // `a?.b`
	nuAssertThat(t, innerMostOptChain.IsOptChainGetProp()).IsTrue()
	nuAssertThat(t, innerMostOptChain.IsOptionalChainStart()).IsTrue()
	nuAssertThat(t, NodeUtil_IsEndOfOptChainSegment(innerMostOptChain)).IsFalse()
	nuAssertThat(t, NodeUtil_IsEndOfFullOptChain(innerMostOptChain)).IsFalse()
}

func TestNodeUtil_GetEndOfOptChainTests_IsEndOfChain_call_innerChain2(t *testing.T) {
	optChainCall := nuParseExpr(t, "a?.b?.(c.x?.y)")
	nuAssertThat(t, optChainCall.IsOptChainCall()).IsTrue()
	nuAssertThat(t, optChainCall.IsOptionalChainStart()).IsTrue()

	// Check that `c.x?.y` is the start of a new chain and also its end
	innerOptChain := optChainCall.GetLastChild()
	nuAssertThat(t, innerOptChain.IsOptChainGetProp()).IsTrue()
	nuAssertThat(t, innerOptChain.IsOptionalChainStart()).IsTrue()
	nuAssertThat(t, NodeUtil_IsEndOfOptChainSegment(innerOptChain)).IsTrue()
	nuAssertThat(t, NodeUtil_IsEndOfFullOptChain(innerOptChain)).IsTrue()
}

func TestNodeUtil_GetEndOfOptChainTests_IsEndOfChain_callAndGetProp_innerChain(t *testing.T) {
	optChainGetProp := nuParseExpr(t, "a?.b(x?.y).c")
	nuAssertThat(t, optChainGetProp.IsOptChainGetProp()).IsTrue()
	nuAssertThat(t, optChainGetProp.IsOptionalChainStart()).IsFalse()

	// Check that `a?.b(x?.y)` is not the start or end of the chain
	optChainCall := optChainGetProp.GetFirstChild()
	nuAssertThat(t, optChainCall.IsOptChainCall()).IsTrue()
	nuAssertThat(t, optChainCall.IsOptionalChainStart()).IsFalse()
	nuAssertThat(t, NodeUtil_IsEndOfOptChainSegment(optChainCall)).IsFalse()
	nuAssertThat(t, NodeUtil_IsEndOfFullOptChain(optChainCall)).IsFalse()

	// Check that `x?.y` is the start and end of the chain
	innerOptChain := optChainCall.GetLastChild()
	nuAssertThat(t, innerOptChain.IsOptChainGetProp()).IsTrue()
	nuAssertThat(t, innerOptChain.IsOptionalChainStart()).IsTrue()
	nuAssertThat(t, NodeUtil_IsEndOfOptChainSegment(innerOptChain)).IsTrue()
	nuAssertThat(t, NodeUtil_IsEndOfFullOptChain(innerOptChain)).IsTrue()
}

func TestNodeUtil_GetEndOfOptChainTests_IsEndOfChain_callAndGetProp_innerChain_multipleInnerArgs(t *testing.T) {
	optChainGetProp := nuParseExpr(t, "a?.b(x?.y, c, d).e")
	nuAssertThat(t, optChainGetProp.IsOptChainGetProp()).IsTrue()
	nuAssertThat(t, optChainGetProp.IsOptionalChainStart()).IsFalse()

	// Check that `a?.b(x?.y, c, d)` is not the start or end of its chain
	optChainCall := optChainGetProp.GetFirstChild()
	nuAssertThat(t, optChainCall.IsOptChainCall()).IsTrue()
	nuAssertThat(t, optChainCall.IsOptionalChainStart()).IsFalse()
	nuAssertThat(t, NodeUtil_IsEndOfOptChainSegment(optChainCall)).IsFalse()
	nuAssertThat(t, NodeUtil_IsEndOfFullOptChain(optChainCall)).IsFalse()

	// Check that `x?.y` is the start and end of its chain
	innerOptChain := optChainCall.GetSecondChild()
	nuAssertThat(t, innerOptChain.IsOptChainGetProp()).IsTrue()
	nuAssertThat(t, innerOptChain.IsOptionalChainStart()).IsTrue()
	nuAssertThat(t, NodeUtil_IsEndOfOptChainSegment(innerOptChain)).IsTrue()
	nuAssertThat(t, NodeUtil_IsEndOfFullOptChain(innerOptChain)).IsTrue()
}

func TestNodeUtil_GetEndOfOptChainTests_ShortChain(t *testing.T) {
	// `expr?.prop1.prop2`
	innerGetProp := rhino.IR_StartOptChainGetprop(rhino.IR_Name("expr"), "pro1")
	outerGetProp := rhino.IR_ContinueOptChainGetprop(innerGetProp, "prop2")

	nuAssertThat(t, NodeUtil_GetEndOfOptChainSegment(innerGetProp)).IsEqualTo(outerGetProp)
}

func TestNodeUtil_GetEndOfOptChainTests_TwoChains(t *testing.T) {
	// `expr()?.prop1.prop2()?.[prop3]`
	call := rhino.IR_Call(rhino.IR_Name("expr"))
	startOptGetProp := rhino.IR_StartOptChainGetprop(call, "prop1")
	optChainGetProp := rhino.IR_ContinueOptChainGetprop(startOptGetProp, "prop2")
	optChainCall := rhino.IR_ContinueOptChainCall(optChainGetProp)
	rhino.IR_StartOptChainGetelem(optChainCall, rhino.IR_Name("prop3"))

	nuAssertThat(t, NodeUtil_GetEndOfOptChainSegment(startOptGetProp)).IsEqualTo(optChainCall)
}

func TestNodeUtil_GetEndOfOptChainTests_BreakingOutOfOptChain(t *testing.T) {
	// `(expr?.prop1.prop2).prop3`
	startOptGetProp := rhino.IR_StartOptChainGetprop(rhino.IR_Name("expr"), "prop1")
	optChainGetProp := rhino.IR_ContinueOptChainGetprop(startOptGetProp, "prop2")
	rhino.IR_Getprop(optChainGetProp, "prop3")

	nuAssertThat(t, NodeUtil_GetEndOfOptChainSegment(startOptGetProp)).IsEqualTo(optChainGetProp)
}

// ==== CreateSynthesizedExternsSymbolTests (extends CompilerTestCase)

// Java: NodeUtilTest.CreateSynthesizedExternsSymbolTests
type NodeUtilTest_CreateSynthesizedExternsSymbolTests struct {
	CompilerTestCaseBase
}

// Java: NodeUtilTest.CreateSynthesizedExternsSymbolTests()
func NewNodeUtilTest_CreateSynthesizedExternsSymbolTests() *NodeUtilTest_CreateSynthesizedExternsSymbolTests {
	c := &NodeUtilTest_CreateSynthesizedExternsSymbolTests{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestNodeUtil_CreateSynthesizedExternsSymbolTests(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewNodeUtilTest_CreateSynthesizedExternsSymbolTests() })
}

// Java: NodeUtilTest.CreateSynthesizedExternsSymbolTests.getProcessor(Compiler)
func (c *NodeUtilTest_CreateSynthesizedExternsSymbolTests) GetProcessor(compiler *Compiler) CompilerPass {
	return CompilerPassFunc(func(externs rhino.Node, srcs rhino.Node) {
		expectedExtern := "TEST_NAME"
		NodeUtil_CreateSynthesizedExternsSymbol(compiler, expectedExtern)
	})
}

// Go-only: test-helper — renamed: the Go name Test is CompilerTestCase's
// test(...), which Java overloads.
//
// Java: NodeUtilTest.CreateSynthesizedExternsSymbolTests.test()
func (c *NodeUtilTest_CreateSynthesizedExternsSymbolTests) TestTest() {
	// setup a simple compilation job
	c.TestExternChanges(c.Externs(""), c.Srcs(""), c.Expected("var TEST_NAME"))
}

// ==== NodeTraversalTests public static class NodeTraversalTests {

func TestNodeUtil_NodeTraversalTests_TestVisitPreOrder(t *testing.T) {
	visitor := &nuRecordingVisitor{}
	NodeUtil_VisitPreOrder(nuh_buildTestTree(t), visitor.Visit)
	nuAssertThat(t, nuGetNodeStringsFromMockVisitor(t, visitor)).ContainsExactly("A", "B", "C", "D", "E", "F", "G").InOrder()
}

func TestNodeUtil_NodeTraversalTests_TestVisitPostOrder(t *testing.T) {
	visitor := &nuRecordingVisitor{}
	NodeUtil_VisitPostOrder(nuh_buildTestTree(t), visitor.Visit)
	nuAssertThat(t, nuGetNodeStringsFromMockVisitor(t, visitor)).ContainsExactly("C", "D", "B", "F", "G", "E", "A").InOrder()
}

func TestNodeUtil_NodeTraversalTests_TestIteratePreOrder(t *testing.T) {
	var nodeNames []string
	for n := range NodeUtil_PreOrderIterable(nuh_buildTestTree(t)) {
		nodeNames = append(nodeNames, n.GetString())
	}

	nuAssertThat(t, nodeNames).ContainsExactly("A", "B", "C", "D", "E", "F", "G").InOrder()
}

func TestNodeUtil_NodeTraversalTests_TestIteratePreOrderWithPredicate(t *testing.T) {
	isNotE := func(n rhino.Node) bool { return n.GetString() != "E" }
	var nodeNames []string
	for n := range NodeUtil_PreOrderIterableWithPred(nuh_buildTestTree(t), isNotE) {
		nodeNames = append(nodeNames, n.GetString())
	}

	nuAssertThat(t, nodeNames).ContainsExactly("A", "B", "C", "D", "E").InOrder()
}

func TestNodeUtil_NodeTraversalTests_ProhibitAddingNonAllowableFeatureToScriptExceptModules(t *testing.T) {
	scriptNode := nuParse(t, "")
	compiler := NewCompiler() // Java: new Compiler()
	compiler.SetAllowableFeatures(parser.FeatureSet_BARE_MINIMUM)
	NodeUtil_AddFeatureToScript(scriptNode, parser.FeatureSet_Feature_MODULES, compiler)

	// Adding MODULES to the script featureSet is okay, any other feature isn't
	nuAssertThat(t, NodeUtil_GetFeatureSetOfScript(scriptNode)).IsEqualTo(parser.FeatureSet_BARE_MINIMUM.With(parser.FeatureSet_Feature_MODULES))
	nuAssertFS(t, compiler.GetAllowableFeatures()).Equals(parser.FeatureSet_BARE_MINIMUM.With(parser.FeatureSet_Feature_MODULES))

	compiler.SetAllowableFeatures(parser.FeatureSet_BARE_MINIMUM) // reset
	assertThrows[jutil.IIllegalStateException](t, func() { NodeUtil_AddFeatureToScript(scriptNode, parser.FeatureSet_Feature_LET_DECLARATIONS, compiler) })
}

func TestNodeUtil_NodeTraversalTests_RemoveFeatureFromAllScriptUpdatesCompilerFeatureSet(t *testing.T) {
	rootNode := rhino.IR_Root()
	rootNode.AddChildToFront(nuParse(t, ""))
	compiler := NewCompiler() // Java: new Compiler()
	compiler.SetAllowableFeatures(parser.FeatureSet_ES2020)
	nuAssertFS(t, compiler.GetAllowableFeatures()).Equals(parser.FeatureSet_ES2020)

	NodeUtil_RemoveFeatureFromAllScripts(rootNode, parser.FeatureSet_Feature_LET_DECLARATIONS, compiler)
	nuAssertFS(t, compiler.GetAllowableFeatures()).Equals(parser.FeatureSet_ES2020.WithoutFeatures(parser.FeatureSet_Feature_LET_DECLARATIONS))
}

func TestNodeUtil_NodeTraversalTests_RemoveFeatureSetFromAllScriptUpdatesCompilerFeatureSet(t *testing.T) {
	rootNode := rhino.IR_Root()
	rootNode.AddChildToFront(nuParse(t, ""))
	compiler := NewCompiler() // Java: new Compiler()
	compiler.SetAllowableFeatures(parser.FeatureSet_ES2020)
	nuAssertFS(t, compiler.GetAllowableFeatures()).Equals(parser.FeatureSet_ES2020)

	NodeUtil_RemoveFeaturesFromAllScripts(
		rootNode,
		parser.FeatureSet_BARE_MINIMUM.WithFeatures(parser.FeatureSet_Feature_LET_DECLARATIONS, parser.FeatureSet_Feature_CONST_DECLARATIONS),
		compiler)
	nuAssertFS(t, compiler.GetAllowableFeatures()).Equals(parser.FeatureSet_ES2020.WithoutFeatures(parser.FeatureSet_Feature_LET_DECLARATIONS, parser.FeatureSet_Feature_CONST_DECLARATIONS))
}

/**
 * Builds a node tree of string {@link com.google.javascript.rhino.Node} with string labels that
 * can be used to verify order in tests.
 *
 * <p>The resultant tree looks like this.
 *
 * <pre>
 *      A
 *    ↙   ↘
 *   B     E
 *  ↙ ↘    ↙ ↘
 * C   D  F   G
 * </pre>
 */
func nuh_buildTestTree(t *testing.T) rhino.Node {
	a := rhino.Node_NewString("A")

	b := rhino.Node_NewString("B")
	b.AddChildToBack(rhino.Node_NewString("C"))
	b.AddChildToBack(rhino.Node_NewString("D"))
	a.AddChildToBack(b)

	e := rhino.Node_NewString("E")
	e.AddChildToBack(rhino.Node_NewString("F"))
	e.AddChildToBack(rhino.Node_NewString("G"))
	a.AddChildToBack(e)

	return a
}

// nuRecordingVisitor stands in for the Mockito mock of NodeUtil.Visitor
// and its ArgumentCaptor: it records every visited node.
//
// Go-only: test-helper — Mockito.
type nuRecordingVisitor struct {
	visited []rhino.Node
}

func (v *nuRecordingVisitor) Visit(node rhino.Node) { v.visited = append(v.visited, node) }

// Java: NodeUtilTest.NodeTraversalTests.getNodeStringsFromMockVisitor(NodeUtil.Visitor)
func nuGetNodeStringsFromMockVisitor(t *testing.T, mockVisitor *nuRecordingVisitor) []string {
	// verify(mockVisitor, Mockito.atLeastOnce()).visit(captor.capture())
	nuAssertThat(t, len(mockVisitor.visited) > 0).IsTrue()
	var strs []string
	for _, n := range mockVisitor.visited {
		strs = append(strs, n.GetString())
	}
	return strs
}

// ==== GoogRequireInfoTest (@Parameterized: src, name, require)

func TestNodeUtil_GoogRequireInfoTest_Test(t *testing.T) {
	cases := []struct {
		src     string
		name    string
		require *NodeUtil_GoogRequire
	}{
		{
			"goog.module('a.b.c'); const {Bar} = goog.require('d.Foo');",
			"Bar",
			NodeUtil_GoogRequire_FromNamespaceAndProperty("d.Foo", "Bar", true),
		},
		{"goog.module('a.b.c'); const {Bar} = goog.require('d.Foo');", "Foo", nil},
		{
			"goog.module('a.b.c'); const {Bar: BarLocal} = goog.require('d.Foo');",
			"BarLocal",
			NodeUtil_GoogRequire_FromNamespaceAndProperty("d.Foo", "Bar", true),
		},
		{
			"goog.module('a.b.c'); const {Bar: BarLocal} =\n goog.require('d.Foo');\n",
			"Bar",
			nil,
		},
		{
			"goog.module('a.b.c'); const Foo = goog.require('d.Foo');",
			"Foo",
			NodeUtil_GoogRequire_FromNamespace("d.Foo", true),
		},
		{
			"goog.module('a.b.c'); const dFoo = goog.require('d.Foo');",
			"dFoo",
			NodeUtil_GoogRequire_FromNamespace("d.Foo", true),
		},
		{
			"goog.module('a.b.c'); const Foo = goog.requireType('d.Foo');",
			"Foo",
			NodeUtil_GoogRequire_FromNamespace("d.Foo", false),
		},
		{
			"goog.module('a.b.c'); const {Bar} = goog.requireType('d.Foo');",
			"Bar",
			NodeUtil_GoogRequire_FromNamespaceAndProperty("d.Foo", "Bar", false),
		},
		// Test that non-requires just return null.
		{"goog.module('a.b.c'); let Foo;", "Foo", nil},
		{"goog.module('a.b.c'); let [Foo] = arr;", "Foo", nil},
		{"goog.module('a.b.c'); let {Bar: {Foo}} = obj;", "Foo", nil},
		{"goog.module('a.b.c'); const Foo = 0;", "Foo", nil},
		{"const Foo = goog.require('d.Foo');", "Foo", nil}, // NB: Foo will be null here.
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("src=%s, name=%s, GoogRequire=%s", c.src, c.name, jutil.ToJavaString(c.require)), func(t *testing.T) {
			first, compiler := nuParseWithCompiler(t, c.src)
			scopeCreator := NewSyntacticScopeCreator(compiler)
			globalScope := scopeCreator.CreateScopeOfScope(first, nil)
			module := nuGetNodeOrNull(first, rhino.Token_MODULE_BODY)
			localScope := globalScope
			if !module.IsNil() {
				localScope = scopeCreator.CreateScopeOfScope(module, globalScope)
			}
			got := NodeUtil_GetGoogRequireInfoWithScope(c.name, localScope)
			if c.require == nil {
				nuAssertThat(t, got).IsNull()
			} else {
				nuAssertThat(t, got).IsEqualTo(c.require)
			}
		})
	}
}

// ==== ReferencesReceiverTest (@Parameterized: js, refToEnclosing, refToOuterFnOwn)

type nuReferencesReceiverCase struct {
	js              string
	refToEnclosing  bool
	refToOuterFnOwn bool
}

// Java: NodeUtilTest.ReferencesReceiverTest.cases()
func nuReferencesReceiverCases() []nuReferencesReceiverCase {
	type entry struct {
		key   string
		value bool
	}
	//
	//
	templateToDefinesOwnReceiver := []entry{
		{"      (x = (%s)) => {}", false},
		{"      (      ) => (%s)", false},
		{"async (x = (%s)) => {}", false},
		{"async (      ) => (%s)", false},
		{"      function   f(        ) { (%s); }", true},
		{"async function   f(        ) { (%s); }", true},
		{"async function  *f(        ) { (%s); }", true},
		{"      function  *f(        ) { (%s); }", true},
		{"      function   f(x = (%s)) {       }", true},
		{"      function  *f(x = (%s)) {       }", true},
		{"async function   f(x = (%s)) {       }", true},
		{"async function  *f(x = (%s)) {       }", true},
		{"class F {        f(        ) { (%s); } }", true},
		{"class F {       *f(        ) { (%s); } }", true},
		{"class F { async  f(        ) { (%s); } }", true},
		{"class F { async *f(        ) { (%s); } }", true},
		{"class F {        f(x = (%s)) {       } }", true},
		{"class F {       *f(x = (%s)) {       } }", true},
		{"class F { async  f(x = (%s)) {       } }", true},
		{"class F { async *f(x = (%s)) {       } }", true},
		{"({               f(        ) { (%s); } })", true},
		{"({              *f(        ) { (%s); } })", true},
		{"({        async  f(        ) { (%s); } })", true},
		{"({        async *f(        ) { (%s); } })", true},
		{"({               f(x = (%s)) {       } })", true},
		{"({              *f(x = (%s)) {       } })", true},
		{"({        async  f(x = (%s)) {       } })", true},
		{"({        async *f(x = (%s)) {       } })", true},
	}

	//
	//
	exprToUsesReceiver := []entry{
		{"this", true},
		{"1 || this", true},
		{"{ x: this, }", true},
		{"1", false},
	}

	var cases []nuReferencesReceiverCase
	for _, outer := range templateToDefinesOwnReceiver {
		for _, inner := range templateToDefinesOwnReceiver {
			for _, e := range exprToUsesReceiver {
				caseSrc := fmt.Sprintf(outer.key, fmt.Sprintf(inner.key, e.key))
				cases = append(cases, nuReferencesReceiverCase{
					caseSrc,
					// refToEnclosing
					!outer.value && !inner.value && e.value,
					// refToOuterFnOwn
					outer.value && !inner.value && e.value,
				})
			}
		}
	}

	// Add a few cases using `super` to check it behaves the same.
	//
	// These aren't exhaustive becuase it's hard to construct valid strings that use `super`.
	cases = append(cases,
		nuReferencesReceiverCase{"class F { f() { super.a() } }", false, true},
		nuReferencesReceiverCase{"class F { f() { () => super.a(); } }", false, true},
		nuReferencesReceiverCase{"() => class F { f() { super.a(); } }", false, false},
		nuReferencesReceiverCase{"({ f() { super.a() } })", false, true},
		nuReferencesReceiverCase{"({ f() { () => super.a(); } })", false, true},
		nuReferencesReceiverCase{"() => ({ f() { super.a(); } })", false, false},
		// TODO(b/538170231): NodeUtil.referencesOwnReceiver accounts for new.target
		nuReferencesReceiverCase{"function f() { new.target; }", false, false},
		nuReferencesReceiverCase{"function f() { () => new.target; }", false, false},
		nuReferencesReceiverCase{"() => function f() { new.target; }", false, false},
	)

	return cases
}

func TestNodeUtil_ReferencesReceiverTest_TestReferencesEnclosingReceiver_ofScript(t *testing.T) {
	for _, c := range nuReferencesReceiverCases() {
		t.Run(fmt.Sprintf("\"%s\"", c.js), func(t *testing.T) {
			node := nuParseFirst(t, rhino.Token_SCRIPT, c.js)
			nuAssertThat(t, NodeUtil_ReferencesEnclosingReceiver(node)).IsEqualTo(c.refToEnclosing)
		})
	}
}

func TestNodeUtil_ReferencesReceiverTest_TestReferencesEnclosingReceiver_ofFn(t *testing.T) {
	for _, c := range nuReferencesReceiverCases() {
		t.Run(fmt.Sprintf("\"%s\"", c.js), func(t *testing.T) {
			node := nuParseFirst(t, rhino.Token_FUNCTION, c.js)
			nuAssertThat(t, NodeUtil_ReferencesEnclosingReceiver(node)).IsEqualTo(c.refToEnclosing)
		})
	}
}

func TestNodeUtil_ReferencesReceiverTest_TestReferencesOwnReceiver_ofFn(t *testing.T) {
	for _, c := range nuReferencesReceiverCases() {
		t.Run(fmt.Sprintf("\"%s\"", c.js), func(t *testing.T) {
			node := nuParseFirst(t, rhino.Token_FUNCTION, c.js)
			nuAssertThat(t, NodeUtil_ReferencesOwnReceiver(node)).IsEqualTo(c.refToOuterFnOwn)
		})
	}
}

// ==== GetRValueOfLValueTest (@Parameterized: token, member)

func TestNodeUtil_GetRValueOfLValueTest_Test(t *testing.T) {
	cases := []struct {
		token  rhino.Token
		member string
	}{
		// CLASS_MEMBERS
		{rhino.Token_MEMBER_FUNCTION_DEF, "constructor() { }"},
		{rhino.Token_MEMBER_FUNCTION_DEF, "foo() { }"},
		{rhino.Token_MEMBER_FUNCTION_DEF, "static foo() { }"},
		{rhino.Token_GETTER_DEF, "get foo() { }"},
		{rhino.Token_GETTER_DEF, "static get foo() { }"},
		{rhino.Token_SETTER_DEF, "set foo(x) { }"},
		{rhino.Token_SETTER_DEF, "static set foo(x) { }"},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("%v in \"%s\"", c.token, c.member), func(t *testing.T) {
			node := nuParseFirst(t, c.token, "class F { "+c.member+" }")
			nuAssertNode(t, NodeUtil_GetRValueOfLValue(node)).HasType(rhino.Token_FUNCTION)
		})
	}
}

// ==== IteratesImpureIterableTest (@Parameterized: token, source, expectation)

func TestNodeUtil_IteratesImpureIterableTest_Test(t *testing.T) {
	cases := []struct {
		token       rhino.Token
		source      string
		expectation bool
	}{
		{rhino.Token_ITER_SPREAD, "[...[]]", false},
		{rhino.Token_ITER_SPREAD, "[...[danger()]]", false},
		{rhino.Token_ITER_SPREAD, "[...'lit']", false},
		{rhino.Token_ITER_SPREAD, "[...`template`]", false},
		{rhino.Token_ITER_SPREAD, "[...`template ${sub}`]", false},
		{rhino.Token_ITER_SPREAD, "[...`template ${danger()}`]", false},
		{rhino.Token_ITER_SPREAD, "[...danger]", true},
		{rhino.Token_ITER_SPREAD, "[...danger()]", true},
		{rhino.Token_ITER_SPREAD, "[...5]", true},
		{rhino.Token_ITER_SPREAD, "foo(...[])", false},
		{rhino.Token_ITER_SPREAD, "foo(...[danger()])", false},
		{rhino.Token_ITER_SPREAD, "foo(...'lit')", false},
		{rhino.Token_ITER_SPREAD, "foo(...`template`)", false},
		{rhino.Token_ITER_SPREAD, "foo(...`template ${safe}`)", false},
		{rhino.Token_ITER_SPREAD, "foo(...`template ${danger()}`)", false},
		{rhino.Token_ITER_SPREAD, "foo(...danger)", true},
		{rhino.Token_ITER_SPREAD, "foo(...danger())", true},
		{rhino.Token_ITER_SPREAD, "foo(...5)", true},
		{rhino.Token_ITER_SPREAD, "new foo(...[])", false},
		{rhino.Token_ITER_SPREAD, "new foo(...[danger()])", false},
		{rhino.Token_ITER_SPREAD, "new foo(...'lit')", false},
		{rhino.Token_ITER_SPREAD, "new foo(...`template`)", false},
		{rhino.Token_ITER_SPREAD, "new foo(...`template ${safe}`)", false},
		{rhino.Token_ITER_SPREAD, "new foo(...`template ${danger()}`)", false},
		{rhino.Token_ITER_SPREAD, "new foo(...danger)", true},
		{rhino.Token_ITER_SPREAD, "new foo(...danger())", true},
		{rhino.Token_ITER_SPREAD, "new foo(...5)", true},
		{rhino.Token_ITER_REST, "const [...rest] = []", true},
		{rhino.Token_ITER_REST, "const [...rest] = 'lit'", true},
		{rhino.Token_ITER_REST, "const [...rest] = `template`", true},
		{rhino.Token_ITER_REST, "const [...rest] = safe", true},
		{rhino.Token_ITER_REST, "function f([...rest]) { }", true},
		{rhino.Token_ITER_REST, "const [[...rest]] = safe", true},
		{rhino.Token_ITER_REST, "const {key: [...rest]} = safe", true},
		{rhino.Token_ITER_REST, "function f(...x) { }", false},
		{rhino.Token_ITER_REST, "function f(a, ...x) { }", false},
		{rhino.Token_ITER_REST, "async function f(...x) { }", false},
		{rhino.Token_ITER_REST, "function* f(...x) { }", false},
		{rhino.Token_ITER_REST, "async function* f(...x) { }", false},
		{rhino.Token_ITER_REST, "((...x) => { })", false},
		{rhino.Token_FOR_OF, "for (let x of []) {}", false},
		{rhino.Token_FOR_OF, "for (let x of [danger()]) {}", false},
		{rhino.Token_FOR_OF, "for (let x of 'lit') {}", false},
		{rhino.Token_FOR_OF, "for (let x of `template`) {}", false},
		{rhino.Token_FOR_OF, "for (let x of `template ${safe}`) {}", false},
		{rhino.Token_FOR_OF, "for (let x of `template ${danger()}`) {}", false},
		{rhino.Token_FOR_OF, "for (let x of danger) {}", true},
		{rhino.Token_FOR_OF, "for (let x of danger()) {}", true},
		{rhino.Token_FOR_OF, "for (let x of 5) {}", true},
		{rhino.Token_FOR_AWAIT_OF, "(async()=>{ for await (let x of []) {} })", false},
		{rhino.Token_FOR_AWAIT_OF, "(async()=>{ for await (let x of [danger()]) {} })", false},
		{rhino.Token_FOR_AWAIT_OF, "(async()=>{ for await (let x of 'literal') {} })", false},
		{rhino.Token_FOR_AWAIT_OF, "(async()=>{ for await (let x of `template`) {} })", false},
		{rhino.Token_FOR_AWAIT_OF, "(async()=>{ for await (let x of `t ${safe}`) {} })", false},
		{rhino.Token_FOR_AWAIT_OF, "(async()=>{ for await (let x of `t ${dn()}`) {} })", false},
		{rhino.Token_FOR_AWAIT_OF, "(async()=>{ for await (let x of danger) {} })", true},
		{rhino.Token_FOR_AWAIT_OF, "(async()=>{ for await (let x of danger()) {} })", true},
		{rhino.Token_FOR_AWAIT_OF, "(async()=>{ for await (let x of 5) {} })", true},
		{rhino.Token_YIELD, "function* f() { yield* []; }", false},
		{rhino.Token_YIELD, "function* f() { yield* [danger()]; }", false},
		{rhino.Token_YIELD, "function* f() { yield* 'lit'; }", false},
		{rhino.Token_YIELD, "function* f() { yield* `template`; }", false},
		{rhino.Token_YIELD, "function* f() { yield* `template ${sub}`; }", false},
		{rhino.Token_YIELD, "function* f() { yield* `template ${danger()}`; }", false},
		{rhino.Token_YIELD, "function* f() { yield* danger; }", true},
		{rhino.Token_YIELD, "function* f() { yield* danger(); }", true},
		{rhino.Token_YIELD, "function* f() { yield* 5; }", true},
		{rhino.Token_YIELD, "function* f() { yield danger(); }", false},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("%v in \"%s\"", c.token, c.source), func(t *testing.T) {
			node := nuParseFirst(t, c.token, c.source)
			nuAssertThat(t, NodeUtil_IteratesImpureIterable(node)).IsEqualTo(c.expectation)
		})
	}
}

// ==== CanBeSideEffectedTest (@Parameterized: token, source, expectation)

func TestNodeUtil_CanBeSideEffectedTest_Test(t *testing.T) {
	cases := []struct {
		token       rhino.Token
		source      string
		expectation bool
	}{
		// TODO: Expand test cases for more node types.
		{rhino.Token_SUPER, "super()", false},
		{rhino.Token_SUPER, "super.foo()", false},
		{rhino.Token_CALL, "super()", true},
		{rhino.Token_CALL, "super.foo()", true},
		{rhino.Token_OPTCHAIN_CALL, "x?.()", true},
		{rhino.Token_OPTCHAIN_GETPROP, "x?.y", true},
		{rhino.Token_OPTCHAIN_GETELEM, "x?.[y]", true},
		{rhino.Token_GETPROP, "Symbol.iterator", false},
		{rhino.Token_GETPROP, "Fake.Symbol.iterator", true},

		// TAGGED_TEMPLATELIT
		{rhino.Token_TAGGED_TEMPLATELIT, "tag`literal`", true},

		// AWAIT
		{rhino.Token_AWAIT, "(async () => { await p; })", true},

		// DYNAMIC_IMPORT
		{rhino.Token_DYNAMIC_IMPORT, "import('m')", true},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("%v in \"%s\"", c.token, c.source), func(t *testing.T) {
			node := nuParseFirst(t, c.token, c.source)
			nuAssertThat(t, NodeUtil_CanBeSideEffected(node)).IsEqualTo(c.expectation)
		})
	}
}

// ==== NodeMetricsTest public static final class NodeMetricsTest {
func TestNodeUtil_NodeMetricsTest_TestEstimateNumLines(t *testing.T) {
	nuAssertThat(t, NodeUtil_EstimateNumLines(nuParse(t, ""))).IsEqualTo(2)
	nuAssertThat(t, NodeUtil_EstimateNumLines(nuParse(t, "const x = 1;\nconst y = 2;\n"))).IsEqualTo(3)
	nuAssertThat(t,
		NodeUtil_EstimateNumLines(
			nuParse(t,
				"/* some\nlong\nmulti\nline\ncomment\n*/\nconst x = 1;\nconst y = 2;\n"))).IsEqualTo(9)
}

// Java: NodeUtilTest.getNameNodeFrom(String, String)
func nuh_getNameNodeFrom(t *testing.T, code string, name string) rhino.Node {
	ast := nuParse(t, code)
	return nuh_getNameNode(t, ast, name)
}

// Java: NodeUtilTest.getStringKeyNodeFrom(String, String)
func nuh_getStringKeyNodeFrom(t *testing.T, code string, name string) rhino.Node {
	ast := nuParse(t, code)
	return nuh_getStringNode(t, ast, name, rhino.Token_STRING_KEY)
}

// Java: NodeUtilTest.executedOnceTestCase(String)
func nuh_executedOnceTestCase(t *testing.T, code string) bool {
	nameNode := nuh_getNameNodeFrom(t, code, "x")
	return NodeUtil_IsExecutedExactlyOnce(nameNode)
}

// Java: NodeUtilTest.getFunctionLValue(String)
func nuh_getFunctionLValue(t *testing.T, js string) *string {
	lVal := NodeUtil_GetBestLValue(nuParseFirst(t, rhino.Token_FUNCTION, js))
	if lVal.IsNil() {
		return nil
	}
	return jutil.StringPtr(lVal.GetString())
}

// Java: NodeUtilTest.functionIsRValueOfAssign(String)
func nuh_functionIsRValueOfAssign(t *testing.T, js string) bool {
	ast := nuParse(t, js)
	nameNode := nuh_getNameNode(t, ast, "x")
	funcNode := nuGetNode(t, ast, rhino.Token_FUNCTION)
	nuAssertWithMessage(t, "No function node to test").That(funcNode).IsNotNull()
	return funcNode == NodeUtil_GetRValueOfLValue(nameNode)
}

// Java: NodeUtilTest.testFunctionName(String, String)
func nuh_testFunctionName(t *testing.T, js string, expected any) {
	nuAssertThat(t, NodeUtil_GetNearestFunctionName(nuParseFirst(t, rhino.Token_FUNCTION, js))).IsEqualTo(expected)
}

// nuh_visitLhsNodesInNode: js is the JavaScript node to be passed to
// NodeUtil.visitLhsNodesInNode. Must be either an EXPR_RESULT containing an
// assignment operation (e.g. =, +=, /=, etc) in which case the assignment
// node will be passed to NodeUtil.visitLhsNodesInNode, or a VAR, LET, or
// CONST statement, in which case the declaration statement will be passed.
//
// Java: NodeUtilTest.visitLhsNodesInNode(String)
func nuh_visitLhsNodesInNode(t *testing.T, js string) []rhino.Node {
	root := nuParse(t, js)
	jutil.CheckStateMsg1(root.IsScript(), "%s", root)
	root = root.GetOnlyChild()
	if root.IsExprResult() {
		root = root.GetOnlyChild()
		jutil.CheckStateMsg1(NodeUtil_IsAssignmentOp(root), "%s", root)
	}
	var nodes []rhino.Node
	NodeUtil_VisitLhsNodesInNode(root, func(n rhino.Node) { nodes = append(nodes, n) })
	return nodes
}

// Java: NodeUtilTest.getNameNode(Node, String)
func nuh_getNameNode(t *testing.T, n rhino.Node, name string) rhino.Node {
	return nuh_getStringNode(t, n, name, rhino.Token_NAME)
}

// Java: NodeUtilTest.getStringNode(Node, String, Token)
func nuh_getStringNode(t *testing.T, n rhino.Node, name string, nodeType rhino.Token) rhino.Node {
	if nodeType == n.GetToken() && n.GetString() == name {
		return n
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		result := nuh_getStringNode(t, c, name, nodeType)
		if !result.IsNil() {
			return result
		}
	}
	return rhino.NilNode
}

// Java: NodeUtilTest.getStringLitNode(Node, String)
func nuh_getStringLitNode(t *testing.T, n rhino.Node, str string) rhino.Node {
	return nuh_getStringNode(t, n, str, rhino.Token_STRINGLIT)
}

// Java: NodeUtilTest.getStringLitNodeFrom(String, String)
func nuh_getStringLitNodeFrom(t *testing.T, code string, str string) rhino.Node {
	ast := nuParse(t, code)
	return nuh_getStringLitNode(t, ast, str)
}

// Java: NodeUtilTest.getNumberNode(Node, double)
func nuh_getNumberNode(t *testing.T, n rhino.Node, number float64) rhino.Node {
	if n.IsNumber() &&
		jutil.ObjectsEquals(n.GetDouble(), number) { // equals allow checks for -0 and NaN
		return n
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		result := nuh_getNumberNode(t, c, number)
		if !result.IsNil() {
			return result
		}
	}
	return rhino.NilNode
}

// nuh_getPattern returns the first node in tree that is an array pattern or
// object pattern.
//
// Java: NodeUtilTest.getPattern(Node)
func nuh_getPattern(t *testing.T, tree rhino.Node) rhino.Node {
	if tree.IsDestructuringPattern() {
		return tree
	}
	for c := tree.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		result := nuh_getPattern(t, c)
		if !result.IsNil() {
			return result
		}
	}
	return rhino.NilNode
}

// Java: NodeUtilTest.isValidPropertyName(String)
func nuh_isValidPropertyName(t *testing.T, s string) bool {
	return NodeUtil_IsValidPropertyName(parser.FeatureSet_ES3, s)
}

// Java: NodeUtilTest.isValidQualifiedName(String)
func nuh_isValidQualifiedName(t *testing.T, s string) bool {
	return NodeUtil_IsValidQualifiedNameFeatureSet(parser.FeatureSet_ES3, s)
}

// Go-only: the record methods of NodeUtil.GoogRequire (equals, hashCode,
// toString), against values printed by the JDK 21 for the same records.
func TestNodeUtil_GoogRequireRecordMethods(t *testing.T) {
	cases := []struct {
		r        *NodeUtil_GoogRequire
		hashCode int32
		str      string
	}{
		{NodeUtil_GoogRequire_FromNamespace("d.Foo", true), -58942997,
			"GoogRequire[namespace=d.Foo, property=null, isStrongRequire=true]"},
		{NodeUtil_GoogRequire_FromNamespace("d.Foo", false), -58942991,
			"GoogRequire[namespace=d.Foo, property=null, isStrongRequire=false]"},
		{NodeUtil_GoogRequire_FromNamespaceAndProperty("d.Foo", "Bar", true), -56880040,
			"GoogRequire[namespace=d.Foo, property=Bar, isStrongRequire=true]"},
		{NodeUtil_GoogRequire_FromNamespaceAndProperty("", "é", false), 8460,
			"GoogRequire[namespace=, property=é, isStrongRequire=false]"},
	}
	for i, c := range cases {
		nuAssertThat(t, c.r.HashCode()).IsEqualTo(c.hashCode)
		nuAssertThat(t, c.r.String()).IsEqualTo(c.str)
		for j, d := range cases {
			nuAssertThat(t, c.r.Equals(d.r)).IsEqualTo(i == j)
		}
		nuAssertThat(t, c.r.Equals(NewNodeUtil_GoogRequire(c.r.Namespace(), c.r.Property(), c.r.IsStrongRequire()))).IsTrue()
	}
}

// Go-only: test-helper — not in NodeUtilTest. new ArrayList<>(set) is an empty
// list, never null, and a nil slice is Java null to the callers of
// getParentChangeScopeNodes and removeNestedChangeScopeNodes.
func TestNodeUtil_GoOnly_ChangeScopeNodesEmptyIsNotNull(t *testing.T) {
	if NodeUtil_GetParentChangeScopeNodes([]rhino.Node{}) == nil {
		t.Error("getParentChangeScopeNodes(empty) returned null")
	}
	if NodeUtil_RemoveNestedChangeScopeNodes([]rhino.Node{}) == nil {
		t.Error("removeNestedChangeScopeNodes(empty) returned null")
	}
	script := rhino.IR_Script()
	fn := rhino.IR_Function(rhino.IR_Name("f"), rhino.IR_ParamList(), rhino.IR_Block())
	script.AddChildToBack(fn)
	if got := NodeUtil_RemoveNestedChangeScopeNodes([]rhino.Node{fn, script}); len(got) != 1 || got[0] != script {
		t.Errorf("removeNestedChangeScopeNodes([f, SCRIPT]) = %v, want [SCRIPT]", got)
	}
}
