// Port of com/google/javascript/rhino/testing/NodeSubject.java @ closure-compiler v20261021

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Bob Jervis
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

// (with rhino/testing/TextDiffFactsBuilder.java)

/*
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Bob Jervis
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package jscomp

// Go-only: placement — rhino.testing.NodeSubject uses
// jscomp.testing.ColorSubject, so in Go it lives with the jscomp test
// helpers (package rhinotesting cannot import jscomp).

import (
	"fmt"
	"math/big"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/rhinotesting"
	"github.com/paullepoulpe/jsgo/internal/testjunit"
)

// NodeSubject is a Truth Subject for the Node class. Usage:
//
//	NodeSubject_AssertNode(node1).IsEqualTo(node2)
//	NodeSubject_AssertNode(node1).HasType(rhino.Token_FUNCTION)
type NodeSubject struct {
	*Subject
	actual                  rhino.Node
	serializer              func(rhino.Node) string
	genericNameReplacements *jutil.LinkedHashMap[string, string]
}

// NodeSubject_AssertNode is NodeSubject.assertNode(Node).
//
// Java: NodeSubject.assertNode(Node)
func NodeSubject_AssertNode(node rhino.Node) *NodeSubject {
	return newNodeSubject(testjunit.NewSubject(node), node)
}

// AssertNode is NodeSubject_AssertNode for Java's
// `import static ...NodeSubject.assertNode` (Go-only: language-adapter).
func AssertNode(node rhino.Node) *NodeSubject {
	return NodeSubject_AssertNode(node)
}

// Java: NodeSubject(FailureMetadata, Node)
func newNodeSubject(s *Subject, node rhino.Node) *NodeSubject {
	return &NodeSubject{Subject: s, actual: node, genericNameReplacements: jutil.NewLinkedHashMap[string, string]()}
}

// check is Subject.check(format, args): a subject for a property of the
// node.
func (s *NodeSubject) check(format string, args ...any) *StandardSubjectBuilder {
	msg := "value of: node." + jutil.LenientFormat(format, args...) + "\nnode was: " + s.serializeNode(s.actual)
	if s.Message() != "" {
		msg = s.Message() + "\n" + msg
	}
	return testjunit.NewStandardSubjectBuilder(msg)
}

// UsingSerializer specifies a function to use when rendering Nodes into
// assertion failure messages.
//
// A common choice of serializer is Compiler::toSource, to as render
// JavaScript code.
//
// Java: NodeSubject.usingSerializer(Function<Node, String>)
func (s *NodeSubject) UsingSerializer(serializer func(rhino.Node) string) *NodeSubject {
	s.serializer = serializer
	return s
}

// Java: NodeSubject.withGenericNameReplacements(ImmutableMap<String, String>)
func (s *NodeSubject) WithGenericNameReplacements(genericNameReplacements *jutil.LinkedHashMap[string, string]) *NodeSubject {
	s.genericNameReplacements = genericNameReplacements
	return s
}

// IsEqualTo compares the node-trees of actual and expected.
//
// The per-node comparison ignores:
//   - Types
//   - JSDoc
//   - Side-effects
//
// Java's isEqualTo(Object) overload throws UnsupportedOperationException;
// Go's static typing makes it unreachable.
//
// TODO(nickreid): This isn't really equality based. Use a different name.
//
// Java: NodeSubject.isEqualTo(Node)
func (s *NodeSubject) IsEqualTo(expected rhino.Node) {
	s.IsEqualToInternal(expected /* checkJsdoc= */, false)
}

// IsEqualIncludingJsDocTo compares the node-trees of actual and expected.
//
// The per-node comparison ignores:
//   - Types
//   - Side-effects
//
// TODO(nickreid): This isn't really equality based. Use a different name.
//
// Java: NodeSubject.isEqualIncludingJsDocTo(Node)
func (s *NodeSubject) IsEqualIncludingJsDocTo(expected rhino.Node) {
	s.IsEqualToInternal(expected /* checkJsdoc= */, true)
}

// TODO(nickreid): This isn't really equality based. Use a different name.
//
// Java: NodeSubject.isEqualToInternal(Node, boolean)
func (s *NodeSubject) IsEqualToInternal(expected rhino.Node, checkJsdoc bool) {
	s.IsNotNull()
	NodeSubject_AssertNode(expected).IsNotNull()

	// Fast exit for the typical TAP case where all tests pass.
	if nodeSubject_findFirstMismatch(s.actual, expected, checkJsdoc) == nil {
		return
	}

	// If there's a text diff, always show the full actual JS. Intended to be machine-readable.
	_, parseableActual := os.LookupEnv("NodeSubject.PARSEABLE_ACTUAL")

	facts := []string{}

	// If applicable, diff each Script separately.
	nodeDiffs := s.maybeExtractScriptsSeparately(expected, s.actual)

	for i := 0; i < len(nodeDiffs); i++ {
		nodeDiff := nodeDiffs[i]
		expectedNode := nodeDiff.expected
		actualNode := nodeDiff.actual

		mismatch := nodeSubject_findFirstMismatch(actualNode, expectedNode, checkJsdoc)
		if mismatch == nil {
			continue
		}

		errorPrefix := nodeSubject_maybeGetErrorPrefix(nodeDiffs, i, actualNode)

		expectedOutputJs := s.serializeNode(expectedNode)
		actualOutputJs := s.serializeNode(actualNode)
		if expectedOutputJs == actualOutputJs {
			// The output code looks identical, so diff the AST instead to show what properties
			// are different.
			facts = append(facts,
				NewTextDiffFactsBuilder(errorPrefix+"AST diff").
					ExpectedText(mismatch.expected.ToStringTree()).
					ActualText(mismatch.actual.ToStringTree()).
					Build()...)
		} else if parseableActual {
			facts = append(facts,
				errorPrefix+
					"Actual JS (parseable): <<<"+
					s.applyGenericNames(actualOutputJs)+
					">>>")
		} else if jutil.StringTrim(expectedOutputJs) == "" {
			// Don't use the text diff fact because everything is new. Avoid prefixing the output with
			// +.
			facts = append(facts,
				errorPrefix+
					"Empty expected JS. Actual JS is:\n    "+
					// Indent by 4 spaces.
					strings.ReplaceAll(s.applyGenericNames(actualOutputJs), "\n", "\n    "))
		} else {
			expectedJsNormalized := s.applyGenericNames(expectedOutputJs)
			actualJsNormalized := s.applyGenericNames(actualOutputJs)
			if expectedJsNormalized == actualJsNormalized {
				// Somehow, the generic name replacement eliminated any diff. Use the original input.
				expectedJsNormalized = expectedOutputJs
				actualJsNormalized = actualOutputJs
			}
			facts = append(facts,
				NewTextDiffFactsBuilder(errorPrefix+"JS diff").
					ExpectedText(expectedJsNormalized).
					ActualText(actualJsNormalized).
					Build()...)
		}
		if checkJsdoc {
			expectedJsDoc := nodeSubject_jsdocToStringNullsafe(mismatch.expected.GetJSDocInfo())
			actualJsDoc := nodeSubject_jsdocToStringNullsafe(mismatch.actual.GetJSDocInfo())
			if expectedJsDoc != actualJsDoc {
				facts = append(facts,
					NewTextDiffFactsBuilder(errorPrefix+"JSDoc diff").
						ExpectedText(expectedJsDoc).
						ActualText(actualJsDoc).
						Build()...)
			}
		}
	}

	if len(facts) != 0 {
		s.Failf("Node tree inequality\n%s", strings.Join(facts, "\n"))
	}
}

// nodeSubject_NodeDiff is the record NodeDiff(Node expected, Node actual).
type nodeSubject_NodeDiff struct {
	expected rhino.Node
	actual   rhino.Node
}

// Java: NodeSubject.maybeExtractScriptsSeparately(Node, Node)
func (s *NodeSubject) maybeExtractScriptsSeparately(expected rhino.Node, actual rhino.Node) []nodeSubject_NodeDiff {
	nodeDiffs := []nodeSubject_NodeDiff{}
	if nodeSubject_isScriptSourceRoot(expected) &&
		nodeSubject_isScriptSourceRoot(actual) &&
		expected.GetChildCount() == actual.GetChildCount() {
		// The length is the same so we just check one. N^2 but N should be small.
		for i := 0; i < expected.GetChildCount(); i++ {
			nodeDiffs = append(nodeDiffs, nodeSubject_NodeDiff{expected.GetChildAtIndex(i), actual.GetChildAtIndex(i)})
		}
	} else {
		nodeDiffs = append(nodeDiffs, nodeSubject_NodeDiff{expected, actual})
	}
	return nodeDiffs
}

// Java: NodeSubject.isScriptSourceRoot(Node)
func nodeSubject_isScriptSourceRoot(node rhino.Node) bool {
	return node.IsRoot() &&
		node.HasParent() &&
		node.GetParent().IsRoot() &&
		node.GetParent().GetLastChild() == node
}

// Java: NodeSubject.maybeGetErrorPrefix(List<NodeDiff>, int, Node)
func nodeSubject_maybeGetErrorPrefix(nodeDiffs []nodeSubject_NodeDiff, i int, actualNode rhino.Node) string {
	if len(nodeDiffs) <= 1 {
		return ""
	}

	sourceFileName := actualNode.GetSourceFileName()
	if sourceFileName != nil && *sourceFileName != "" {
		return fmt.Sprintf("For expected %d, script \"%s\": ", i, *sourceFileName)
	}

	return fmt.Sprintf("For expected %d: ", i)
}

// Java: NodeSubject.applyGenericNames(String)
func (s *NodeSubject) applyGenericNames(input string) string {
	for k, v := range s.genericNameReplacements.All() {
		// Replace the actual name with the generic name.
		input = strings.ReplaceAll(input, k, v)
	}
	return input
}

// Java: NodeSubject.isEquivalentTo(Node)
func (s *NodeSubject) IsEquivalentTo(other rhino.Node) *NodeSubject {
	s.check("isEquivalentTo(%s)", other).That(s.actual.IsEquivalentTo(other)).IsTrue()
	return s
}

// Java: NodeSubject.isNotEquivalentTo(Node)
func (s *NodeSubject) IsNotEquivalentTo(other rhino.Node) *NodeSubject {
	s.check("isEquivalentTo(%s)", other).That(s.actual.IsEquivalentTo(other)).IsFalse()
	return s
}

// Java: NodeSubject.hasJSTypeThat()
func (s *NodeSubject) HasJSTypeThat() *rhinotesting.TypeSubject {
	b := s.check("getJSType()")
	return rhinotesting.NewTruth(assertionTB{}).AssertWithMessage("%s", b.Message()).AboutTypes().That(s.actual.GetJSTypeRequired())
}

// Java: NodeSubject.hasColorThat()
func (s *NodeSubject) HasColorThat() *ColorSubject {
	return newColorSubject(s.check("getColor()").That(s.actual.GetColor()), s.actual.GetColor())
}

// Java: NodeSubject.hasType(Token)
func (s *NodeSubject) HasType(type_ rhino.Token) {
	s.HasToken(type_)
}

// Java: NodeSubject.hasToken(Token)
func (s *NodeSubject) HasToken(token rhino.Token) *NodeSubject {
	s.check("getToken()").That(s.actual.GetToken()).IsEqualTo(token)
	return s
}

// Java: NodeSubject.isString(String)
func (s *NodeSubject) IsString(value string) *NodeSubject {
	s.check("getToken()").That(s.actual.GetToken()).IsEqualTo(rhino.Token_STRINGLIT)
	s.check("getString()").That(s.actual.GetString()).IsEqualTo(value)
	return s
}

// Java: NodeSubject.isName(String)
func (s *NodeSubject) IsName(name string) *NodeSubject {
	s.check("isName()").That(s.actual.IsName()).IsTrue()
	s.check("getString()").That(s.actual.GetString()).IsEqualTo(name)
	return s
}

// IsNumber takes Java's int and double overloads alike, as a float64.
//
// Java: NodeSubject.isNumber(int)
// Java: NodeSubject.isNumber(double)
func (s *NodeSubject) IsNumber(value float64) *NodeSubject {
	s.check("isNumber()").That(s.actual.IsNumber()).IsTrue()
	s.check("getNumber()").That(s.actual.GetDouble()).IsEqualTo(value)
	return s
}

// Java: NodeSubject.isBigInt(BigInteger)
func (s *NodeSubject) IsBigInt(value *big.Int) *NodeSubject {
	s.check("isBigInt()").That(s.actual.IsBigInt()).IsTrue()
	s.check("getBigInt()").That(s.actual.GetBigInt().Cmp(value) == 0).IsTrue()
	return s
}

// Java: NodeSubject.isAssign()
func (s *NodeSubject) IsAssign() *NodeSubject {
	s.check("isAssign()").That(s.actual.IsAssign()).IsTrue()
	return s
}

// Java: NodeSubject.isNullNode()
func (s *NodeSubject) IsNullNode() *NodeSubject {
	s.check("isNull()").That(s.actual.IsNull()).IsTrue()
	return s
}

// Java: NodeSubject.isThis()
func (s *NodeSubject) IsThis() *NodeSubject {
	s.check("isThis()").That(s.actual.IsThis()).IsTrue()
	return s
}

// Java: NodeSubject.isSuper()
func (s *NodeSubject) IsSuper() *NodeSubject {
	s.check("isSuper()").That(s.actual.IsSuper()).IsTrue()
	return s
}

// Java: NodeSubject.isFunction()
func (s *NodeSubject) IsFunction() *NodeSubject {
	s.HasToken(rhino.Token_FUNCTION)
	return s
}

// Java: NodeSubject.isCast()
func (s *NodeSubject) IsCast() *NodeSubject {
	s.HasToken(rhino.Token_CAST)
	return s
}

// Java: NodeSubject.isArrowFunction()
func (s *NodeSubject) IsArrowFunction() *NodeSubject {
	s.check("isArrowFunction()").That(s.actual.IsArrowFunction()).IsTrue()
	return s
}

// Java: NodeSubject.isStatic()
func (s *NodeSubject) IsStatic() *NodeSubject {
	s.check("isStatic()").That(s.actual.IsStaticMember()).IsTrue()
	return s
}

// Java: NodeSubject.hasTrailingComma()
func (s *NodeSubject) HasTrailingComma() *NodeSubject {
	s.check("hasTrailingComma()").That(s.actual.HasTrailingComma()).IsTrue()
	return s
}

// IsOptionalChainStart indicates whether the node we are asserting is the
// start of an optional chain e.g. `a?.b` of `a?.b.c`
//
// Java: NodeSubject.isOptionalChainStart()
func (s *NodeSubject) IsOptionalChainStart() *NodeSubject {
	s.check("isOptionalChainStart()").That(s.actual.IsOptionalChainStart()).IsTrue()
	return s
}

// IsNotOptionalChainStart indicates whether the node we are asserting is
// the start of an optional chain e.g. `b.c` of `a?.b.c`
//
// Java: NodeSubject.isNotOptionalChainStart()
func (s *NodeSubject) IsNotOptionalChainStart() *NodeSubject {
	s.check("isOptionalChainStart()").That(s.actual.IsOptionalChainStart()).IsFalse()
	return s
}

// Java: NodeSubject.isParamList()
func (s *NodeSubject) IsParamList() *NodeSubject {
	s.HasToken(rhino.Token_PARAM_LIST)
	return s
}

// Java: NodeSubject.isCall()
func (s *NodeSubject) IsCall() *NodeSubject {
	s.check("isCall()").That(s.actual.IsCall()).IsTrue()
	return s
}

// Java: NodeSubject.isFreeCall()
func (s *NodeSubject) IsFreeCall() *NodeSubject {
	s.check("callable").
		That(s.actual.IsCall() || s.actual.IsOptChainCall() || s.actual.IsTaggedTemplateLit()).
		IsTrue()
	s.check("getBooleanProp(Node.FREE_CALL)").That(s.actual.GetBooleanProp(rhino.Node_FREE_CALL)).IsTrue()
	return s
}

// Java: NodeSubject.isNotFreeCall()
func (s *NodeSubject) IsNotFreeCall() *NodeSubject {
	s.check("getBooleanProp(Node.FREE_CALL)").That(s.actual.GetBooleanProp(rhino.Node_FREE_CALL)).IsFalse()
	return s
}

// Java: NodeSubject.isConst()
func (s *NodeSubject) IsConst() *NodeSubject {
	s.check("isConst()").That(s.actual.IsConst()).IsTrue()
	return s
}

// Java: NodeSubject.isVar()
func (s *NodeSubject) IsVar() *NodeSubject {
	s.check("isVar()").That(s.actual.IsVar()).IsTrue()
	return s
}

// Java: NodeSubject.isGetProp()
func (s *NodeSubject) IsGetProp() *NodeSubject {
	s.check("isGetProp()").That(s.actual.IsGetProp()).IsTrue()
	return s
}

// Java: NodeSubject.isBlock()
func (s *NodeSubject) IsBlock() *NodeSubject {
	s.check("isBlock()").That(s.actual.IsBlock()).IsTrue()
	return s
}

// Java: NodeSubject.isObjectLit()
func (s *NodeSubject) IsObjectLit() *NodeSubject {
	s.check("isObjectLit()").That(s.actual.IsObjectLit()).IsTrue()
	return s
}

// Java: NodeSubject.isScript()
func (s *NodeSubject) IsScript() *NodeSubject {
	s.check("isScript()").That(s.actual.IsScript()).IsTrue()
	return s
}

// Java: NodeSubject.isReturn()
func (s *NodeSubject) IsReturn() *NodeSubject {
	s.check("isReturn()").That(s.actual.IsReturn()).IsTrue()
	return s
}

// Java: NodeSubject.isExprResult()
func (s *NodeSubject) IsExprResult() *NodeSubject {
	s.check("isExprResult()").That(s.actual.IsExprResult()).IsTrue()
	return s
}

// Java: NodeSubject.isMemberFunctionDef(String)
func (s *NodeSubject) IsMemberFunctionDef(name string) *NodeSubject {
	s.check("isMemberFunction()").That(s.actual.IsMemberFunctionDef()).IsTrue()
	s.check("getString()").That(s.actual.GetString()).IsEqualTo(name)
	return s
}

// Java: NodeSubject.matchesName(String)
func (s *NodeSubject) MatchesName(qname string) *NodeSubject {
	s.check("matchesName(%s)", qname).That(s.actual.MatchesName(qname)).IsTrue()
	return s
}

// Java: NodeSubject.matchesQualifiedName(String)
func (s *NodeSubject) MatchesQualifiedName(qname string) *NodeSubject {
	s.check("matchesQualifiedName(%s)", qname).That(s.actual.MatchesQualifiedName(qname)).IsTrue()
	return s
}

// HasSourceFileName takes the Java @Nullable String as *string.
//
// Java: NodeSubject.hasSourceFileName(String)
func (s *NodeSubject) HasSourceFileName(sourceFileName *string) *NodeSubject {
	s.check("getSourceFileName()").That(s.actual.GetSourceFileName()).IsEqualTo(sourceFileName)
	return s
}

// Java: NodeSubject.hasCharno(int)
func (s *NodeSubject) HasCharno(charno int) *NodeSubject {
	s.check("getCharno()").That(s.actual.GetCharno()).IsEqualTo(charno)
	return s
}

// Java: NodeSubject.hasLineno(int)
func (s *NodeSubject) HasLineno(lineno int) *NodeSubject {
	s.check("getLineno()").That(s.actual.GetLineno()).IsEqualTo(lineno)
	return s
}

// Java: NodeSubject.hasLength(int)
func (s *NodeSubject) HasLength(length int) *NodeSubject {
	s.check("getLength()").That(s.actual.GetLength()).IsEqualTo(length)
	return s
}

// Java: NodeSubject.hasEqualSourceInfoTo(Node)
func (s *NodeSubject) HasEqualSourceInfoTo(other rhino.Node) *NodeSubject {
	return s.HasSourceFileName(other.GetSourceFileName()).
		HasLineno(other.GetLineno()).
		HasCharno(other.GetCharno()).
		HasLength(other.GetLength())
}

// Java: NodeSubject.isIndexable(boolean)
func (s *NodeSubject) IsIndexable(isIndexable bool) *NodeSubject {
	s.check("isIndexable()").That(s.actual.IsIndexable()).IsEqualTo(isIndexable)
	return s
}

// Java: NodeSubject.hasOriginalName(String)
func (s *NodeSubject) HasOriginalName(originalName string) *NodeSubject {
	s.check("getOriginalName()").That(s.actual.GetOriginalName()).IsEqualTo(originalName)
	return s
}

// Java: NodeSubject.hasChildren(boolean)
func (s *NodeSubject) HasChildren(hasChildren bool) *NodeSubject {
	s.check("hasChildren()").That(s.actual.HasChildren()).IsEqualTo(hasChildren)
	return s
}

// Java: NodeSubject.hasXChildren(int)
func (s *NodeSubject) HasXChildren(numChildren int) *NodeSubject {
	s.check("getChildCount()").That(s.actual.GetChildCount()).IsEqualTo(numChildren)
	return s
}

// Java: NodeSubject.hasNoChildren()
func (s *NodeSubject) HasNoChildren() *NodeSubject {
	return s.HasXChildren(0)
}

// Java: NodeSubject.hasOneChild()
func (s *NodeSubject) HasOneChild() *NodeSubject {
	return s.HasXChildren(1)
}

// Java: NodeSubject.hasOneChildThat()
func (s *NodeSubject) HasOneChildThat() *NodeSubject {
	s.HasOneChild()
	return NodeSubject_AssertNode(s.actual.GetOnlyChild())
}

// Java: NodeSubject.hasFirstChildThat()
func (s *NodeSubject) HasFirstChildThat() *NodeSubject {
	s.HasChildren(true)
	return NodeSubject_AssertNode(s.actual.GetFirstChild())
}

// Java: NodeSubject.hasSecondChildThat()
func (s *NodeSubject) HasSecondChildThat() *NodeSubject {
	s.check("getChildCount()").That(s.actual.GetChildCount()).IsAtLeast(2)
	return NodeSubject_AssertNode(s.actual.GetSecondChild())
}

// Java: NodeSubject.hasLastChildThat()
func (s *NodeSubject) HasLastChildThat() *NodeSubject {
	s.HasChildren(true)
	return NodeSubject_AssertNode(s.actual.GetLastChild())
}

// Java: NodeSubject.isFromExterns()
func (s *NodeSubject) IsFromExterns() *NodeSubject {
	s.check("isFromExterns()").That(s.actual.IsFromExterns()).IsTrue()
	return s
}

// Java: NodeSubject.hasStringThat()
func (s *NodeSubject) HasStringThat() *Subject {
	return s.check("getString()").That(s.actual.GetString())
}

// Java: NodeSubject.hasOriginalNameThat()
func (s *NodeSubject) HasOriginalNameThat() *Subject {
	return s.check("getOriginalName()").That(s.actual.GetOriginalName())
}

// showActual is Subject.actualCustomStringRepresentation().
//
// Java: NodeSubject.actualCustomStringRepresentation()
func (s *NodeSubject) showActual() string {
	return s.serializeNode(s.actual)
}

// IsNotNull is Subject.isNotNull() with the node rendered by the serializer.
func (s *NodeSubject) IsNotNull() {
	if s.actual.IsNil() {
		s.Failf("expected not to be null")
	}
}

// nodeSubject_findFirstMismatch compares the given node-trees recursively
// and return the first pair of nodes that differs doing a pre-order
// traversal.
//
// jsDoc: Whether to check for differences in JSDoc.
//
// Java: NodeSubject.findFirstMismatch(Node, Node, boolean)
func nodeSubject_findFirstMismatch(actual rhino.Node, expected rhino.Node, jsDoc bool) *nodeSubject_NodeMismatch {
	if !actual.IsEquivalentToWithFlagsAndSideEffect(
		expected /* compareType= */, false /* recurse= */, false, jsDoc /* sideEffect= */, false) {
		return newNodeSubject_NodeMismatch(actual, expected)
	}

	// `isEquivalentTo` confirms that the number of children is the same.
	actualChild := actual.GetFirstChild()
	expectedChild := expected.GetFirstChild()
	for !actualChild.IsNil() {
		mismatch := nodeSubject_findFirstMismatch(actualChild, expectedChild, jsDoc)
		if mismatch != nil {
			return mismatch
		}

		actualChild = actualChild.GetNext()
		expectedChild = expectedChild.GetNext()
	}

	// `isEquivalentTo` confirms that either both nodes have a shadow root or none do.
	actualShadowAst := actual.GetClosureUnawareShadow()
	expectedShadowAst := expected.GetClosureUnawareShadow()
	if !actualShadowAst.IsNil() && !expectedShadowAst.IsNil() {
		mismatch := nodeSubject_findFirstMismatch(actualShadowAst, expectedShadowAst, jsDoc)
		if mismatch != nil {
			return mismatch
		}
	}

	return nil
}

// nodeSubject_NodeMismatch is a pair of nodes that were expected to match in
// some way but didn't.
//
// Java: NodeSubject.NodeMismatch
type nodeSubject_NodeMismatch struct {
	actual   rhino.Node
	expected rhino.Node
}

// Java: NodeSubject.NodeMismatch(Node, Node)
func newNodeSubject_NodeMismatch(actual rhino.Node, expected rhino.Node) *nodeSubject_NodeMismatch {
	return &nodeSubject_NodeMismatch{actual: actual, expected: expected}
}

var (
	// Java regex "\\{\\n\\s*\\}" (\s is Java's [ \t\n\x0B\f\r]).
	nodeSubject_emptyBlock = regexp.MustCompile(`\{\n[ \t\n\x0B\f\r]*\}`)
	// Java regex " +(\\n|$)".
	nodeSubject_trailingSpaces = regexp.MustCompile(` +(\n|$)`)
)

// Java: NodeSubject.serializeNode(Node)
func (s *NodeSubject) serializeNode(node rhino.Node) string {
	if s.serializer != nil {
		out := s.serializer(node)
		// Collapse empty blocks that span a newline.
		out = nodeSubject_emptyBlock.ReplaceAllString(out, "{}")
		// Trim trailing spaces from each line.
		out = nodeSubject_trailingSpaces.ReplaceAllString(out, "$1")
		return out
	}
	if node.IsNil() {
		return "<Java null>"
	}
	return node.ToStringTree()
}

// Java: NodeSubject.jsdocToStringNullsafe(JSDocInfo)
func nodeSubject_jsdocToStringNullsafe(jsdoc *rhino.JSDocInfo) string {
	if jsdoc == nil {
		return "(null)"
	}
	return jsdoc.ToStringVerbose()
}

// TextDiffFactsBuilder builds Fact objects for use with the Truth library
// that represent the difference between two text strings.
//
// Java: rhino.testing.TextDiffFactsBuilder
type TextDiffFactsBuilder struct {
	title        string
	expectedText *string
	actualText   *string
}

// Java: TextDiffFactsBuilder(String)
func NewTextDiffFactsBuilder(title string) *TextDiffFactsBuilder {
	return &TextDiffFactsBuilder{title: title}
}

// Java: TextDiffFactsBuilder.expectedText(String)
func (b *TextDiffFactsBuilder) ExpectedText(expectedText string) *TextDiffFactsBuilder {
	b.expectedText = &expectedText
	return b
}

// Java: TextDiffFactsBuilder.actualText(String)
func (b *TextDiffFactsBuilder) ActualText(actualText string) *TextDiffFactsBuilder {
	b.actualText = &actualText
	return b
}

// Build returns one or more Fact objects representing the difference
// between the expected and actual text strings, which must be specified by
// calling their methods before calling this one.
//
// Go-only: test-helper — java-diff-utils' unified diff is replaced by an
// LCS line diff in the same format (context of 10 lines); it only renders
// failure messages.
//
// Java: TextDiffFactsBuilder.build()
func (b *TextDiffFactsBuilder) Build() []string {
	// The diff algorithm expects to work on a list, so we need to split the text into
	// lines.
	expectedLines := strings.Split(*jutil.CheckNotNull(b.expectedText), "\n")
	actualLines := strings.Split(*jutil.CheckNotNull(b.actualText), "\n")
	return []string{b.title + ":\n" + textDiffFactsBuilder_unifiedDiff(expectedLines, actualLines, 10)}
}

// textDiffFactsBuilder_unifiedDiff renders a unified diff of two line lists.
func textDiffFactsBuilder_unifiedDiff(a, b []string, context int) string {
	n, m := len(a), len(b)
	lcs := make([][]int, n+1)
	for i := range lcs {
		lcs[i] = make([]int, m+1)
	}
	for i := n - 1; i >= 0; i-- {
		for j := m - 1; j >= 0; j-- {
			if a[i] == b[j] {
				lcs[i][j] = lcs[i+1][j+1] + 1
			} else {
				lcs[i][j] = max(lcs[i+1][j], lcs[i][j+1])
			}
		}
	}
	type op struct {
		kind byte // ' ', '-', '+'
		text string
	}
	var ops []op
	i, j := 0, 0
	for i < n || j < m {
		switch {
		case i < n && j < m && a[i] == b[j]:
			ops = append(ops, op{' ', a[i]})
			i++
			j++
		case j < m && (i == n || lcs[i][j+1] >= lcs[i+1][j]):
			ops = append(ops, op{'+', b[j]})
			j++
		default:
			ops = append(ops, op{'-', a[i]})
			i++
		}
	}
	var sb strings.Builder
	sb.WriteString("--- expected\n+++ actual\n")
	for k := 0; k < len(ops); {
		if ops[k].kind == ' ' {
			k++
			continue
		}
		start := max(0, k-context)
		end := k
		for end < len(ops) {
			if ops[end].kind != ' ' {
				end++
				continue
			}
			run := end
			for run < len(ops) && ops[run].kind == ' ' {
				run++
			}
			if run-end > 2*context || run == len(ops) {
				end = min(len(ops), end+context)
				break
			}
			end = run
		}
		sb.WriteString("@@ @@\n")
		for _, o := range ops[start:end] {
			sb.WriteByte(o.kind)
			sb.WriteString(o.text)
			sb.WriteByte('\n')
		}
		k = end
	}
	return strings.TrimSuffix(sb.String(), "\n")
}

// assertionTB is a testing.TB whose failures panic with *jutil.AssertionError, so
// that the rhinotesting subjects (written against testing.TB) fail like Java
// Truth subjects (Go-only: test-helper). Only the methods those subjects
// call are implemented.
type assertionTB struct {
	testing.TB
}

func (assertionTB) Helper() {}

func (assertionTB) Fatal(args ...any) { panic(jutil.NewAssertionError(fmt.Sprint(args...))) }

func (assertionTB) Fatalf(format string, args ...any) {
	panic(jutil.NewAssertionError(fmt.Sprintf(format, args...)))
}

func (assertionTB) Error(args ...any) { panic(jutil.NewAssertionError(fmt.Sprint(args...))) }

func (assertionTB) Errorf(format string, args ...any) {
	panic(jutil.NewAssertionError(fmt.Sprintf(format, args...)))
}

// TypeSubject_AssertType is rhino.testing.TypeSubject.assertType(JSType),
// failing with *jutil.AssertionError.
func TypeSubject_AssertType(type_ rhino.JSType) *rhinotesting.TypeSubject {
	return rhinotesting.NewTruth(assertionTB{}).AssertType(type_)
}

// ColorSubject is the Subject for Color.
//
// Java: jscomp.testing.ColorSubject
type ColorSubject struct {
	*Subject
	actual *colors.Color
}

// Java: ColorSubject(FailureMetadata, Color)
func newColorSubject(s *Subject, actual *colors.Color) *ColorSubject {
	return &ColorSubject{Subject: s, actual: actual}
}

// ColorSubject_AssertThat is ColorSubject.assertThat(Color).
//
// Java: ColorSubject.assertThat(Color)
func ColorSubject_AssertThat(actual *colors.Color) *ColorSubject {
	return newColorSubject(testjunit.NewSubject(actual), actual)
}

// Java: ColorSubject.actualNonNull()
func (s *ColorSubject) actualNonNull() *colors.Color {
	s.IsNotNull()
	return s.actual
}

func (s *ColorSubject) check(format string, args ...any) *StandardSubjectBuilder {
	msg := "value of: color." + jutil.LenientFormat(format, args...) + "\ncolor was: " + show(s.actual)
	if s.Message() != "" {
		msg = s.Message() + "\n" + msg
	}
	return testjunit.NewStandardSubjectBuilder(msg)
}

// Custom assertions

// Java: ColorSubject.hasId(ColorId)
func (s *ColorSubject) HasId(id colors.ColorId) {
	s.check("getId").That(s.actualNonNull().GetId()).IsEqualTo(id)
}

// Java: ColorSubject.isUnion()
func (s *ColorSubject) IsUnion() {
	s.check("isUnion()").That(s.actualNonNull().IsUnion()).IsTrue()
}

// Java: ColorSubject.isPrimitive()
func (s *ColorSubject) IsPrimitive() {
	s.check("isPrimitive()").That(s.actualNonNull().IsPrimitive()).IsTrue()
}

// Java: ColorSubject.isInvalidating()
func (s *ColorSubject) IsInvalidating() {
	s.check("isInvalidating").That(s.actualNonNull().IsInvalidating()).IsTrue()
}

// Java: ColorSubject.isNotInvalidating()
func (s *ColorSubject) IsNotInvalidating() {
	s.check("isInvalidating").That(s.actualNonNull().IsInvalidating()).IsFalse()
}

// Java: ColorSubject.isClosureAssert()
func (s *ColorSubject) IsClosureAssert() {
	s.check("isClosureAssert()").That(s.actualNonNull().IsClosureAssert()).IsTrue()
}

// Java: ColorSubject.isConstructor()
func (s *ColorSubject) IsConstructor() {
	s.check("isConstructor()").That(s.actualNonNull().IsConstructor()).IsTrue()
}

// IsEqualToSubject is isEqualTo(ColorSubject); the Subject's
// isEqualTo(Object) is promoted unchanged.
//
// Java: ColorSubject.isEqualTo(ColorSubject)
func (s *ColorSubject) IsEqualToSubject(otherSubject *ColorSubject) *ColorSubject {
	s.IsEqualTo(otherSubject.actual)
	return s
}

// Java: ColorSubject.propertiesKeepOriginalName()
func (s *ColorSubject) PropertiesKeepOriginalName() {
	s.check("propertiesKeepOriginalName").
		That(s.actualNonNull().GetPropertiesKeepOriginalName()).
		IsTrue()
}

// Java: ColorSubject.hasOwnPropertiesSetThat()
func (s *ColorSubject) HasOwnPropertiesSetThat() *Subject {
	return s.check("getOwnProperties()").That(s.actualNonNull().GetOwnProperties())
}

// Java: ColorSubject.hasAlternates(Color...)
func (s *ColorSubject) HasAlternates(alternates ...*colors.Color) {
	s.IsUnion()
	elems := make([]any, len(alternates))
	for i, a := range alternates {
		elems[i] = a
	}
	s.check("getAlternates().containsExactly()").
		That(s.actualNonNull().GetUnionElements()).
		// cast to Object[] to suppress warning about varargs vs. non-varargs call confusion
		ContainsExactly(elems...)
}

// Java: ColorSubject.hasDisambiguationSupertypesThat(ColorRegistry)
func (s *ColorSubject) HasDisambiguationSupertypesThat(registry *colors.ColorRegistry) *Subject {
	return s.check("registry.getDisambiguationSupertypes(%s)", s.actual).
		That(registry.GetDisambiguationSupertypes(s.actual))
}
