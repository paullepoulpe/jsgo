// Part of com/google/javascript/rhino/testing/NodeSubject.java @ closure-compiler v20261021:
// the methods ParserTest uses beyond those js_doc_info_parser_test.go defines
// on nodeSubject.

package parsing

import (
	"fmt"
	"math/big"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// IsEqualTo compares the node-trees of actual and expected, ignoring types,
// JSDoc and side-effects.
//
// Java: NodeSubject.isEqualTo(Node)
func (s *nodeSubject) IsEqualTo(expected rhino.Node) {
	s.t.Helper()
	if s.actual.IsNil() || expected.IsNil() {
		s.t.Fatalf("expected not to be null: actual %v, expected %v", s.actual, expected)
	}
	if m := nodeSubject_findFirstMismatch(s.actual, expected, false); m != nil {
		s.t.Fatalf("Node tree inequality\nAST diff\n  expected:\n%s\n  actual:\n%s",
			m[1].ToStringTree(), m[0].ToStringTree())
	}
}

// nodeSubject_findFirstMismatch returns the first (actual, expected) pair
// that is not equivalent, or nil.
//
// Java: NodeSubject.findFirstMismatch(Node, Node, boolean)
func nodeSubject_findFirstMismatch(actual, expected rhino.Node, jsDoc bool) []rhino.Node {
	if !actual.IsEquivalentToWithFlagsAndSideEffect(
		expected /* compareType= */, false /* recurse= */, false, jsDoc /* sideEffect= */, false) {
		return []rhino.Node{actual, expected}
	}

	// `isEquivalentTo` confirms that the number of children is the same.
	actualChild := actual.GetFirstChild()
	expectedChild := expected.GetFirstChild()
	for !actualChild.IsNil() {
		if mismatch := nodeSubject_findFirstMismatch(actualChild, expectedChild, jsDoc); mismatch != nil {
			return mismatch
		}
		actualChild = actualChild.GetNext()
		expectedChild = expectedChild.GetNext()
	}

	// `isEquivalentTo` confirms that either both nodes have a shadow root or
	// none do.
	actualShadowAst := actual.GetClosureUnawareShadow()
	expectedShadowAst := expected.GetClosureUnawareShadow()
	if !actualShadowAst.IsNil() && !expectedShadowAst.IsNil() {
		if mismatch := nodeSubject_findFirstMismatch(actualShadowAst, expectedShadowAst, jsDoc); mismatch != nil {
			return mismatch
		}
	}
	return nil
}

func (s *nodeSubject) check(what string, ok bool) *nodeSubject {
	s.t.Helper()
	if !ok {
		s.t.Fatalf("value of: %s\nfor node: %v", what, s.actual)
	}
	return s
}

// Java: NodeSubject.isName(String)
func (s *nodeSubject) IsName(name string) *nodeSubject {
	s.t.Helper()
	s.check("isName()", s.actual.IsName())
	return s.check(fmt.Sprintf("getString() = %q (was %q)", name, s.actual.GetString()), s.actual.GetString() == name)
}

// Java: NodeSubject.isNumber(double)
func (s *nodeSubject) IsNumber(value float64) *nodeSubject {
	s.t.Helper()
	s.check("isNumber()", s.actual.IsNumber())
	return s.check(fmt.Sprintf("getNumber() = %v (was %v)", value, s.actual.GetDouble()),
		jutil.ObjectsEquals(s.actual.GetDouble(), value))
}

// Java: NodeSubject.isBigInt(BigInteger)
func (s *nodeSubject) IsBigInt(value *big.Int) *nodeSubject {
	s.t.Helper()
	s.check("isBigInt()", s.actual.IsBigInt())
	return s.check(fmt.Sprintf("getBigInt() = %v (was %v)", value, s.actual.GetBigInt()),
		jutil.BigIntegerEquals(s.actual.GetBigInt(), value))
}

// Java: NodeSubject.isStatic()
func (s *nodeSubject) IsStatic() *nodeSubject {
	s.t.Helper()
	return s.check("isStatic()", s.actual.IsStaticMember())
}

// Java: NodeSubject.hasTrailingComma()
func (s *nodeSubject) HasTrailingComma() *nodeSubject {
	s.t.Helper()
	return s.check("hasTrailingComma()", s.actual.HasTrailingComma())
}

// Java: NodeSubject.isOptionalChainStart()
func (s *nodeSubject) IsOptionalChainStart() *nodeSubject {
	s.t.Helper()
	return s.check("isOptionalChainStart()", s.actual.IsOptionalChainStart())
}

// Java: NodeSubject.isNotOptionalChainStart()
func (s *nodeSubject) IsNotOptionalChainStart() *nodeSubject {
	s.t.Helper()
	return s.check("isOptionalChainStart() is false", !s.actual.IsOptionalChainStart())
}

// Java: NodeSubject.isFreeCall()
func (s *nodeSubject) IsFreeCall() *nodeSubject {
	s.t.Helper()
	s.check("callable", s.actual.IsCall() || s.actual.IsOptChainCall() || s.actual.IsTaggedTemplateLit())
	return s.check("getBooleanProp(Node.FREE_CALL)", s.actual.GetBooleanProp(rhino.Node_FREE_CALL))
}

// Java: NodeSubject.isNotFreeCall()
func (s *nodeSubject) IsNotFreeCall() *nodeSubject {
	s.t.Helper()
	return s.check("getBooleanProp(Node.FREE_CALL) is false", !s.actual.GetBooleanProp(rhino.Node_FREE_CALL))
}

// Java: NodeSubject.matchesQualifiedName(String)
func (s *nodeSubject) MatchesQualifiedName(qname string) *nodeSubject {
	s.t.Helper()
	return s.check("matchesQualifiedName("+qname+")", s.actual.MatchesQualifiedName(qname))
}

// Java: NodeSubject.hasXChildren(int)
func (s *nodeSubject) HasXChildren(numChildren int) *nodeSubject {
	s.t.Helper()
	return s.check(fmt.Sprintf("getChildCount() = %d (was %d)", numChildren, s.actual.GetChildCount()),
		s.actual.GetChildCount() == numChildren)
}

// Java: NodeSubject.hasChildren(boolean)
func (s *nodeSubject) HasChildren(hasChildren bool) *nodeSubject {
	s.t.Helper()
	return s.check(fmt.Sprintf("hasChildren() = %v", hasChildren), s.actual.HasChildren() == hasChildren)
}

// Java: NodeSubject.hasFirstChildThat()
func (s *nodeSubject) HasFirstChildThat() *nodeSubject {
	s.t.Helper()
	s.HasChildren(true)
	return assertNode(s.t, s.actual.GetFirstChild())
}

// Java: NodeSubject.isFunction()
func (s *nodeSubject) IsFunction() *nodeSubject {
	s.t.Helper()
	return s.HasToken(rhino.Token_FUNCTION)
}
