// Port of com/google/javascript/rhino/NodeTest.java @ closure-compiler v20261015

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
 *   Nick Santos
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhino

// NodeSubject (com.google.javascript.rhino.testing) assertions are inlined:
// assertNode(a).isEquivalentTo(b) is a.isEquivalentTo(b), hasLineno(l) is
// getLineno() == l and hasTrailingComma() is hasTrailingComma().

import (
	"math"
	"slices"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func assertTrue(t *testing.T, what string, got bool) {
	t.Helper()
	if !got {
		t.Errorf("%s: got false, want true", what)
	}
}

func assertFalse(t *testing.T, what string, got bool) {
	t.Helper()
	if got {
		t.Errorf("%s: got true, want false", what)
	}
}

func assertEquivalent(t *testing.T, a, b Node) {
	t.Helper()
	if !a.IsEquivalentTo(b) {
		t.Errorf("not equivalent:\n%s\n%s", a.ToStringTree(), b.ToStringTree())
	}
}

func assertNotEquivalent(t *testing.T, a, b Node) {
	t.Helper()
	if a.IsEquivalentTo(b) {
		t.Errorf("unexpectedly equivalent:\n%s\n%s", a.ToStringTree(), b.ToStringTree())
	}
}

func TestNode_ValidatePropertiesForRoot(t *testing.T) {
	n := IR_Root()
	assertContainsExactly(t, getMessagesFromValidateProperties(n))

	// ROOT nodes shouldn't have properties, not even a source file
	n.SetSourceFileForTesting("file.js")
	assertContainsExactly(t, getMessagesFromValidateProperties(n), "ROOT has properties")
}

func getMessagesFromValidateProperties(n Node) []string {
	var listBuilder []string
	violationMessageConsumer := func(s string) { listBuilder = append(listBuilder, s) }
	n.ValidateProperties(violationMessageConsumer)

	return listBuilder
}

func TestNode_SideEffectFlagsSerialization(t *testing.T) {
	// Test each individual flag
	checkSideEffectFlagsRoundTrip(t, Node_SideEffectFlags_MUTATES_GLOBAL_STATE)
	checkSideEffectFlagsRoundTrip(t, Node_SideEffectFlags_MUTATES_THIS)
	checkSideEffectFlagsRoundTrip(t, Node_SideEffectFlags_MUTATES_ARGUMENTS)
	checkSideEffectFlagsRoundTrip(t, Node_SideEffectFlags_THROWS)
	// Test an arbitrary combination of 2 flags
	checkSideEffectFlagsRoundTrip(t, Node_SideEffectFlags_THROWS|Node_SideEffectFlags_MUTATES_THIS)
}

func checkSideEffectFlagsRoundTrip(t *testing.T, testFlags int32) {
	t.Helper()
	original := IR_Call(IR_Name("f"))
	// serialization only works for nodes that actually have a source file,
	// because in actual execution they always must have that.
	original.SetSourceFileForTesting("sourcefile")

	// Simulate the situation where we've deserialized the node itself, but
	// not its non-source-file properties yet, by cloning the orignal node
	// before adding the side effect flags to it.
	restored := original.CloneNode()

	original.SetSideEffectFlags(testFlags)
	serializedProperties := original.SerializeProperties()

	restored.DeserializeProperties(serializedProperties, false)
	if got := restored.GetSideEffectFlags(); got != testFlags {
		t.Errorf("restored side effect flags = %d, want %d", got, testFlags)
	}
}

func TestNode_ValidatePropertiesForIsParenthesized(t *testing.T) {
	n := IR_String("")
	n.SetSourceFileForTesting("file.js") // avoid error about missing source file

	n.SetIsParenthesized(true)
	n.SetToken(Token_STRING_KEY)
	assertContainsExactly(t, getMessagesFromValidateProperties(n), "non-expression is parenthesized")
}

func TestNode_ValidatePropertiesForFunctionProperties(t *testing.T) {
	n := IR_Empty()
	n.SetSourceFileForTesting("file.js") // avoid error about missing source file

	// Change the Node type so we won't get thrown errors when we try to set
	// the function-only properties.
	n.SetToken(Token_FUNCTION)
	n.SetIsArrowFunction(true)
	n.SetIsAsyncFunction(true)

	assertContainsExactly(t, getMessagesFromValidateProperties(n))
	n.SetToken(Token_EMPTY)
	assertContainsExactly(t, getMessagesFromValidateProperties(n),
		"invalid ARROW_FN prop", "invalid ASYNC_FN prop")
}

func TestNode_ValidatePropertiesForSyntheticProperty(t *testing.T) {
	n := IR_Block()
	n.SetSourceFileForTesting("file.js") // avoid error about missing source file
	n.SetIsSyntheticBlock(true)

	assertContainsExactly(t, getMessagesFromValidateProperties(n))
	n.SetToken(Token_EMPTY)
	assertContainsExactly(t, getMessagesFromValidateProperties(n), "invalid SYNTHETIC prop")
}

func TestNode_ValidatePropertiesForColorFromCast(t *testing.T) {
	n := IR_Name("a")
	n.SetSourceFileForTesting("file.js") // avoid error about missing source file
	n.SetColor(colors.StandardColors_NUMBER)
	n.SetColorFromTypeCast()

	assertContainsExactly(t, getMessagesFromValidateProperties(n))
	n.SetColor(nil)
	assertContainsExactly(t, getMessagesFromValidateProperties(n), "COLOR_FROM_CAST with no Color")
}

func TestNode_ValidatePropertiesForOptChain(t *testing.T) {
	n := IR_Empty()
	n.SetSourceFileForTesting("file.js") // avoid error about missing source file

	n.SetToken(Token_OPTCHAIN_CALL)
	n.SetIsOptionalChainStart(true)

	assertContainsExactly(t, getMessagesFromValidateProperties(n))
	n.SetToken(Token_OPTCHAIN_GETELEM)
	assertContainsExactly(t, getMessagesFromValidateProperties(n))
	n.SetToken(Token_OPTCHAIN_GETELEM)
	assertContainsExactly(t, getMessagesFromValidateProperties(n))
	n.SetToken(Token_EMPTY)
	assertContainsExactly(t, getMessagesFromValidateProperties(n),
		"START_OF_OPT_CHAIN on non-optional Node")
}

func TestNode_ValidatePropertiesForConstVarFlags(t *testing.T) {
	n := IR_Name("a")
	n.SetSourceFileForTesting("file.js") // avoid error about missing source file
	n.SetDeclaredConstantVar(true)
	n.SetInferredConstantVar(true)

	assertContainsExactly(t, getMessagesFromValidateProperties(n))
	n.SetToken(Token_IMPORT_STAR)
	assertContainsExactly(t, getMessagesFromValidateProperties(n))
	n.SetToken(Token_STRINGLIT)
	assertContainsExactly(t, getMessagesFromValidateProperties(n), "invalid CONST_VAR_FLAGS")
}

func TestNode_LinenoCharnoNormal(t *testing.T) {
	assertLinenoCharno(t, 5, 6, 5, 6)
	assertLinenoCharno(t, 456, 3423, 456, 3423)
	assertLinenoCharno(t, 0, 0, 0, 0)
}

func TestNode_LinenoCharnoErroneous(t *testing.T) {
	assertLinenoCharno(t, -5, 90, -1, -1)
	assertLinenoCharno(t, 90, -1, -1, -1)
}

func TestNode_MergeOverflowGraciously(t *testing.T) {
	assertLinenoCharno(t, 89, 4096, 89, 4095)
}

func assertLinenoCharno(t *testing.T, linenoIn, charnoIn, linenoOut, charnoOut int) {
	t.Helper()
	test := IR_Block()
	test.SetLinenoCharno(linenoIn, charnoIn)
	if got := test.GetLineno(); got != linenoOut {
		t.Errorf("setLinenoCharno(%d, %d): lineno = %d, want %d", linenoIn, charnoIn, got, linenoOut)
	}
	if got := test.GetCharno(); got != charnoOut {
		t.Errorf("setLinenoCharno(%d, %d): charno = %d, want %d", linenoIn, charnoIn, got, charnoOut)
	}
}

func TestNode_IsEquivalentToConsidersStartOfOptionalChainProperty(t *testing.T) {
	// `a?.b.c`
	singleSegmentOptChain :=
		IR_ContinueOptChainGetprop(IR_StartOptChainGetprop(IR_Name("a"), "b"), "c")
	assertEquivalent(t, singleSegmentOptChain, singleSegmentOptChain.CloneTree())

	twoSegmentOptChain := singleSegmentOptChain.CloneTree()
	twoSegmentOptChain.SetIsOptionalChainStart(true)
	assertNotEquivalent(t, singleSegmentOptChain, twoSegmentOptChain)
}

func TestNode_IsEquivalentToConsidersDirectEval(t *testing.T) {
	directEvalCall := IR_Call(IR_Name("eval"))
	directEvalCall.GetFirstChild().PutBooleanProp(Node_DIRECT_EVAL, true)

	indirectEvalCall := IR_Call(IR_Name("eval"))
	indirectEvalCall.GetFirstChild().PutBooleanProp(Node_DIRECT_EVAL, false)

	assertNotEquivalent(t, directEvalCall, indirectEvalCall)
}

func TestNode_IsEquivalentToForFunctionsConsidersKindOfFunction(t *testing.T) {
	normalFunction := IR_Function(IR_Name(""), IR_ParamList(), IR_Block())
	assertEquivalent(t, normalFunction, normalFunction.CloneTree())

	// normal vs async function
	asyncFunction := normalFunction.CloneTree()
	asyncFunction.SetIsAsyncFunction(true)
	assertEquivalent(t, asyncFunction, asyncFunction.CloneTree())
	assertNotEquivalent(t, asyncFunction, normalFunction)

	// normal vs arrow function
	arrowFunction := normalFunction.CloneTree()
	arrowFunction.SetIsArrowFunction(true)
	assertEquivalent(t, arrowFunction, arrowFunction.CloneTree())
	assertNotEquivalent(t, arrowFunction, normalFunction)

	// async arrow function vs async only vs arrow only
	asyncArrowFunction := arrowFunction.CloneTree()
	asyncArrowFunction.SetIsAsyncFunction(true)
	assertEquivalent(t, asyncArrowFunction, asyncArrowFunction.CloneTree())
	assertNotEquivalent(t, asyncArrowFunction, arrowFunction)
	assertNotEquivalent(t, asyncArrowFunction, asyncFunction)

	// normal vs generator function
	generatorFunction := normalFunction.CloneTree()
	generatorFunction.SetIsGeneratorFunction(true)
	assertEquivalent(t, generatorFunction, generatorFunction.CloneTree())
	assertNotEquivalent(t, generatorFunction, normalFunction)

	// async generator function vs only async vs only generator
	asyncGeneratorFunction := generatorFunction.CloneTree()
	asyncGeneratorFunction.SetIsAsyncFunction(true)
	assertEquivalent(t, asyncGeneratorFunction, asyncGeneratorFunction.CloneTree())
	assertNotEquivalent(t, asyncGeneratorFunction, asyncFunction)
	assertNotEquivalent(t, asyncGeneratorFunction, generatorFunction)
}

func TestNode_IsEquivalentTo_WithBoolean_IsSame(t *testing.T) {
	node1 := NewNode(Token_LET)
	assertTrue(t, "node1.isEquivalentTo(node1)", node1.IsEquivalentTo(node1))
}

func TestNode_IsEquivalentTo_WithBoolean_IsDifferent(t *testing.T) {
	node1 := NewNode(Token_LET)
	node2 := NewNode(Token_VAR)
	assertFalse(t, "node1.isEquivalentTo(node2)", node1.IsEquivalentTo(node2))
}

func TestNode_IsEquivalentTo_ConsidersDifferentEsModuleExports(t *testing.T) {
	exportAllFrom := NewNode(Token_EXPORT) // export * from './other/module'
	exportAllFrom.PutBooleanProp(Node_EXPORT_ALL_FROM, true)
	exportDefault := NewNode(Token_EXPORT) // export default function foo() {
	exportDefault.PutBooleanProp(Node_EXPORT_DEFAULT, true)
	simpleExport := NewNode(Token_EXPORT) // export {x} or export const x = 0;

	assertFalse(t, "exportAllFrom ~ exportDefault", exportAllFrom.IsEquivalentTo(exportDefault))
	assertFalse(t, "exportAllFrom ~ simpleExport", exportAllFrom.IsEquivalentTo(simpleExport))
	assertFalse(t, "exportDefault ~ simpleExport", exportDefault.IsEquivalentTo(simpleExport))
}

func TestNode_IsEquivalentToNumber(t *testing.T) {
	assertTrue(t, "1 ~ 1", Node_NewNumber(1).IsEquivalentTo(Node_NewNumber(1)))
	assertFalse(t, "1 ~ 2", Node_NewNumber(1).IsEquivalentTo(Node_NewNumber(2)))
}

func TestNode_IsEquivalentToConsidersPresenceOfClosureUnawareShadow(t *testing.T) {
	nameShadow1 := IR_Name("x")
	nameShadow2 := IR_Name("x")
	nameNoShadow := IR_Name("x")

	nameShadow1.SetClosureUnawareShadow(IR_Root())
	nameShadow2.SetClosureUnawareShadow(IR_Root())

	assertTrue(t, "nameShadow1 ~ nameShadow2", nameShadow1.IsEquivalentTo(nameShadow2))
	assertFalse(t, "nameShadow1 ~ nameNoShadow", nameShadow1.IsEquivalentTo(nameNoShadow))
	assertFalse(t, "nameNoShadow ~ nameShadow1", nameNoShadow.IsEquivalentTo(nameShadow1))
}

func TestNode_IsEquivalentToConsidersPresenceOfClosureUnawareShadow_Recursively(t *testing.T) {
	nameShadow1a := IR_Name("x")
	nameShadow1b := IR_Name("x")
	nameShadow2 := IR_Name("x")

	root1Shared := IR_Root(IR_ScriptOfNodes(IR_Block()))
	root2 := IR_Root()

	nameShadow1a.SetClosureUnawareShadow(root1Shared)
	nameShadow1b.SetClosureUnawareShadow(root1Shared.CloneTree())
	nameShadow2.SetClosureUnawareShadow(root2)

	assertTrue(t, "nameShadow1a ~ nameShadow1b", nameShadow1a.IsEquivalentTo(nameShadow1b))
	assertFalse(t, "nameShadow1a ~ nameShadow2", nameShadow1a.IsEquivalentTo(nameShadow2))
}

func TestNode_IsEquivalentTo_DoesNotConsiderPresenceOfClosureUnawareShadow_InDeepNoShadowMode(t *testing.T) {
	nameShadow1 := IR_Name("x")
	nameShadow2 := IR_Name("x")

	root1 := IR_Root(IR_ScriptOfNodes(IR_Block()))
	root2 := IR_Root()

	nameShadow1.SetClosureUnawareShadow(root1)
	nameShadow2.SetClosureUnawareShadow(root2)

	assertTrue(t, "DEEP_NO_SHADOW equivalence",
		nameShadow1.IsEquivalentToWithModes(
			nameShadow2,
			Node_RecursionMode_DEEP_NO_SHADOW,
			Node_TypeComparison_IGNORE,
			Node_JsDocComparison_IGNORE,
			Node_SideEffectComparison_COMPARE))
}

func TestNode_NumberRejects_IsNaN(t *testing.T) {
	assertNumberNodeRejects(t, math.NaN())
	assertNumberNodeRejects(t, -math.NaN())
}

func TestNode_NumberRejects_NegativeValues(t *testing.T) {
	assertNumberNodeRejects(t, -1394793.114)
	assertNumberNodeRejects(t, -1.0)
	assertNumberNodeRejects(t, math.Copysign(0, -1))
	assertNumberNodeRejects(t, math.Inf(-1))
}

func assertNumberNodeRejects(t *testing.T, d float64) {
	t.Helper()
	assertThrows[jutil.Throwable](t, func() { Node_NewNumber(d) })
	assertThrows[jutil.Throwable](t, func() { Node_NewNumber(0.0).SetDouble(d) })
}

func TestNode_BigintRejects_NegativeValues(t *testing.T) {
	assertBigIntNodeRejects(t, "-1394793")
	assertBigIntNodeRejects(t, "-1")
	if Node_NewBigInt(jutil.NewBigIntegerFromString("-0")).IsNil() {
		t.Error("newBigInt(-0) is null")
	}
}

func assertBigIntNodeRejects(t *testing.T, x string) {
	t.Helper()
	assertThrows[jutil.Throwable](t, func() { Node_NewBigInt(jutil.NewBigIntegerFromString(x)) })
	assertThrows[jutil.Throwable](t, func() {
		Node_NewBigInt(jutil.BigInteger_ZERO).SetBigInt(jutil.NewBigIntegerFromString(x))
	})
}

func TestNode_IsEquivalentToBigInt(t *testing.T) {
	assertTrue(t, "1n ~ 1n",
		Node_NewBigInt(jutil.BigInteger_ONE).IsEquivalentTo(Node_NewBigInt(jutil.BigInteger_ONE)))
	assertFalse(t, "1n ~ 10n",
		Node_NewBigInt(jutil.BigInteger_ONE).IsEquivalentTo(Node_NewBigInt(jutil.BigInteger_TEN)))
}

func TestNode_IsEquivalentToString(t *testing.T) {
	assertTrue(t, "'1' ~ '1'", Node_NewString("1").IsEquivalentTo(Node_NewString("1")))
	assertFalse(t, "'1' ~ '2'", Node_NewString("1").IsEquivalentTo(Node_NewString("2")))
}

func TestNode_CheckTreeTypeAwareEqualsSameNull(t *testing.T) {
	node1 := Node_NewStringWithToken(Token_NAME, "f")
	node2 := Node_NewStringWithToken(Token_NAME, "f")
	assertTrue(t, "node1.isEquivalentToTyped(node2)", node1.IsEquivalentToTyped(node2))
}

func TestNode_CheckTreeTypeAwareEqualsColorsSame(t *testing.T) {
	node1 := Node_NewStringWithToken(Token_NAME, "f")
	node1.SetColor(colors.StandardColors_NUMBER)
	node2 := Node_NewStringWithToken(Token_NAME, "f")
	node2.SetColor(colors.StandardColors_NUMBER)
	assertTrue(t, "node1.isEquivalentToTyped(node2)", node1.IsEquivalentToTyped(node2))
}

func TestNode_CheckTreeTypeAwareEqualsColorsSameNull(t *testing.T) {
	node1 := Node_NewStringWithToken(Token_NAME, "f")
	node2 := Node_NewStringWithToken(Token_NAME, "f")
	assertTrue(t, "node1.isEquivalentToTyped(node2)", node1.IsEquivalentToTyped(node2))
}

func TestNode_CheckTreeTypeAwareEqualsColorsDifferent(t *testing.T) {
	node1 := Node_NewStringWithToken(Token_NAME, "f")
	node1.SetColor(colors.StandardColors_NUMBER)
	node2 := Node_NewStringWithToken(Token_NAME, "f")
	node2.SetColor(colors.StandardColors_STRING)
	assertFalse(t, "node1.isEquivalentToTyped(node2)", node1.IsEquivalentToTyped(node2))
}

func TestNode_CheckTreeTypeAwareEqualsColorsDifferentNull(t *testing.T) {
	node1 := Node_NewStringWithToken(Token_NAME, "f")
	node1.SetColor(colors.StandardColors_NUMBER)
	node2 := Node_NewStringWithToken(Token_NAME, "f")
	assertFalse(t, "node1.isEquivalentToTyped(node2)", node1.IsEquivalentToTyped(node2))
}

func TestNode_IsQualifiedName(t *testing.T) {
	assertTrue(t, "a", IR_Name("a").IsQualifiedName())
	assertTrue(t, "$", IR_Name("$").IsQualifiedName())
	assertTrue(t, "_", IR_Name("_").IsQualifiedName())
	assertTrue(t, "a.b", IR_Getprop(IR_Name("a"), "b").IsQualifiedName())
	assertTrue(t, "this.b", IR_Getprop(IR_ThisNode(), "b").IsQualifiedName())
	assertFalse(t, "0", IR_Number(0).IsQualifiedName())
	assertFalse(t, "[]", IR_Arraylit().IsQualifiedName())
	assertFalse(t, "{}", IR_Objectlit().IsQualifiedName())
	assertFalse(t, "''", IR_String("").IsQualifiedName())
	assertFalse(t, "a['b']", IR_Getelem(IR_Name("a"), IR_String("b")).IsQualifiedName())
	assertFalse(t, "a[b].c", // a[b].c
		IR_Getprop(IR_Getelem(IR_Name("a"), IR_String("b")), "c").IsQualifiedName())
	assertFalse(t, "a.b[c]", // a.b[c]
		IR_Getelem(IR_Getprop(IR_Name("a"), "b"), IR_String("c")).IsQualifiedName())
	assertFalse(t, "a()", IR_Call(IR_Name("a")).IsQualifiedName())
	assertFalse(t, "a().b", // a().b
		IR_Getprop(IR_Call(IR_Name("a")), "b").IsQualifiedName())
	assertFalse(t, "(a.b)()", // (a.b)()
		IR_Call(IR_Getprop(IR_Name("a"), "b")).IsQualifiedName())
	assertFalse(t, "'a'", IR_String("a").IsQualifiedName())
	assertFalse(t, "/x/", IR_Regexp(IR_String("x")).IsQualifiedName())
	assertFalse(t, "x++", NewNodeWithChild(Token_INC, IR_Name("x")).IsQualifiedName())
}

func TestNode_MatchesQualifiedName1(t *testing.T) {
	type c struct {
		n    Node
		name string
		want bool
	}
	for i, tc := range []c{
		{IR_Name("a"), "a", true},
		{IR_Name("a"), "ab", false},
		{IR_Name("a"), "a.b", false},
		{IR_Name("a"), ".b", false},
		{IR_Name("a"), "a.", false},

		{nodeQname("a.b"), "a", false},
		{nodeQname("a.b"), "a.b", true},
		{nodeQname("a.b"), "a.bc", false},
		{nodeQname("a.b"), ".b", false},
		{nodeQname("a.b"), "this.b", false},

		{nodeQname("this"), "this", true},
		{nodeQname("this"), "thisx", false},

		{nodeQname("this.b"), "a", false},
		{nodeQname("this.b"), "a.b", false},
		{nodeQname("this.b"), ".b", false},
		{nodeQname("this.b"), "a.", false},
		{nodeQname("this.b"), "super.b", false},
		{nodeQname("this.b"), "this.b", true},

		{nodeQname("super"), "super", true},
		{nodeQname("super"), "superx", false},

		{nodeQname("super.b"), "a", false},
		{nodeQname("super.b"), "a.b", false},
		{nodeQname("super.b"), ".b", false},
		{nodeQname("super.b"), "a.", false},
		{nodeQname("super.b"), "this.b", false},
		{nodeQname("super.b"), "super.b", true},

		{nodeQname("a.b.c"), "a.b.c", true},
		{nodeQname("a.b.c"), "a.b.c", true},

		{IR_ImportStar("a"), "a", true},
		{IR_ImportStar("a"), "b", false},

		{IR_Number(0), "a.b", false},
		{IR_Arraylit(), "a.b", false},
		{IR_Objectlit(), "a.b", false},
		{IR_String(""), "a.b", false},
		{IR_Getelem(IR_Name("a"), IR_String("b")), "a.b", false},
		// a[b].c
		{IR_Getprop(IR_Getelem(IR_Name("a"), IR_String("b")), "c"), "a.b.c", false},
		// a.b[c]
		{IR_Getelem(IR_Getprop(IR_Name("a"), "b"), IR_String("c")), "a.b.c", false},
		{IR_Call(IR_Name("a")), "a", false},
		// a().b
		{IR_Getprop(IR_Call(IR_Name("a")), "b"), "a.b", false},
		// (a.b)()
		{IR_Call(IR_Getprop(IR_Name("a"), "b")), "a.b", false},
		{IR_String("a"), "a", false},
		{IR_Regexp(IR_String("x")), "x", false},
		{NewNodeWithChild(Token_INC, IR_Name("x")), "x", false},
	} {
		if got := tc.n.MatchesQualifiedName(tc.name); got != tc.want {
			t.Errorf("case %d: %s.matchesQualifiedName(%q) = %v, want %v",
				i, strings.TrimSpace(tc.n.ToStringTree()), tc.name, got, tc.want)
		}
	}
}

func TestNode_MatchesQualifiedName2(t *testing.T) {
	type c struct {
		n     Node
		other Node
		want  bool
	}
	for i, tc := range []c{
		{IR_Name("a"), nodeQname("a"), true},
		{IR_Name("a"), nodeQname("a.b"), false},

		{nodeQname("a.b"), nodeQname("a"), false},
		{nodeQname("a.b"), nodeQname("a.b"), true},
		{nodeQname("a.b"), nodeQname(".b"), false},
		{nodeQname("a.b"), nodeQname("this.b"), false},

		{nodeQname("this.b"), nodeQname("a"), false},
		{nodeQname("this.b"), nodeQname("a.b"), false},
		{nodeQname("this.b"), nodeQname("super.b"), false},
		{nodeQname("this.b"), nodeQname("this.b"), true},

		{nodeQname("super.b"), nodeQname("a"), false},
		{nodeQname("super.b"), nodeQname("a.b"), false},
		{nodeQname("super.b"), nodeQname("this.b"), false},
		{nodeQname("super.b"), nodeQname("super.b"), true},

		{nodeQname("a.b.c"), nodeQname("a.b.c"), true},
		{nodeQname("a.b.c"), nodeQname("a.b.c"), true},

		{IR_Number(0), nodeQname("a.b"), false},
		{IR_Arraylit(), nodeQname("a.b"), false},
		{IR_Objectlit(), nodeQname("a.b"), false},
		{IR_String(""), nodeQname("a.b"), false},
		{IR_Getelem(IR_Name("a"), IR_String("b")), nodeQname("a.b"), false},
		// a[b].c
		{IR_Getprop(IR_Getelem(IR_Name("a"), IR_String("b")), "c"), nodeQname("a.b.c"), false},
		{IR_Call(IR_Name("a")), nodeQname("a"), false},
		// a().b
		{IR_Getprop(IR_Call(IR_Name("a")), "b"), nodeQname("a.b"), false},
		// (a.b)()
		{IR_Call(IR_Getprop(IR_Name("a"), "b")), nodeQname("a.b"), false},
		{IR_String("a"), nodeQname("a"), false},
		{IR_Regexp(IR_String("x")), nodeQname("x"), false},
		{NewNodeWithChild(Token_INC, IR_Name("x")), nodeQname("x"), false},
	} {
		if got := tc.n.MatchesQualifiedNameNode(tc.other); got != tc.want {
			t.Errorf("case %d: %s.matchesQualifiedName(%s) = %v, want %v",
				i, strings.TrimSpace(tc.n.ToStringTree()), strings.TrimSpace(tc.other.ToStringTree()), got, tc.want)
		}
	}
	// a.b[c]
	assertFalse(t, "a.b[c] matches a.b.c",
		IR_Getelem(IR_Getprop(IR_Name("a"), "b"), IR_String("c")).MatchesQualifiedName("a.b.c"))
}

func TestNode_MatchesName(t *testing.T) {
	// Empty string are treat as unique.
	assertFalse(t, "'' matches ''", IR_Name("").MatchesName(""))

	assertTrue(t, "a matches a", IR_Name("a").MatchesName("a"))
	assertFalse(t, "a matches a.b", IR_Name("a").MatchesName("a.b"))
	assertFalse(t, "a matches ''", IR_Name("a").MatchesName(""))

	assertFalse(t, "this matches this", IR_ThisNode().MatchesName("this"))
	assertFalse(t, "super matches super", IR_SuperNode().MatchesName("super"))
}

func TestNode_MatchesNameNodes(t *testing.T) {
	assertTrue(t, "a matches a", IR_Name("a").MatchesNameNode(nodeQname("a")))
	assertFalse(t, "a matches a.b", IR_Name("a").MatchesNameNode(nodeQname("a.b")))

	assertFalse(t, "this matches this", IR_ThisNode().MatchesNameNode(nodeQname("this")))
	assertFalse(t, "super matches super", IR_SuperNode().MatchesNameNode(nodeQname("super")))
}

// nodeQname is NodeTest.qname(String) (renamed: QualifiedNameTest's
// qname(Node, String...) shares the package).
func nodeQname(name string) Node {
	endPos := strings.IndexByte(name, '.')
	if endPos == -1 {
		return IR_Name(name)
	}
	var node Node
	nodeName := name[:endPos]
	if nodeName == "this" {
		node = IR_ThisNode()
	} else if nodeName == "super" {
		node = IR_SuperNode()
	} else {
		node = IR_Name(nodeName)
	}
	var startPos int
	for {
		startPos = endPos + 1
		endPos = strings.IndexByte(name[startPos:], '.')
		if endPos != -1 {
			endPos += startPos
		}
		var part string
		if endPos == -1 {
			part = name[startPos:]
		} else {
			part = name[startPos:endPos]
		}
		node = IR_Getprop(node, part)
		if endPos == -1 {
			break
		}
	}

	return node
}

func TestNode_CloneAnnontations(t *testing.T) {
	n := getVarRef("a")
	n.SetLength(1)
	assertFalse(t, "IS_CONSTANT_NAME", n.GetBooleanProp(Node_IS_CONSTANT_NAME))
	n.PutBooleanProp(Node_IS_CONSTANT_NAME, true)
	assertTrue(t, "IS_CONSTANT_NAME", n.GetBooleanProp(Node_IS_CONSTANT_NAME))

	nodeClone := n.CloneNode()
	assertTrue(t, "clone IS_CONSTANT_NAME", nodeClone.GetBooleanProp(Node_IS_CONSTANT_NAME))
	if got := nodeClone.GetLength(); got != 1 {
		t.Errorf("clone length = %d, want 1", got)
	}
}

func TestNode_CloneValues(t *testing.T) {
	number := Node_NewNumber(100.0)
	if got := number.CloneNode().GetDouble(); got != 100.0 {
		t.Errorf("clone double = %v, want 100.0", got)
	}

	string_ := Node_NewString(newString("a"))
	assertTrue(t, "clone getString same instance",
		sameInstance(string_.CloneNode().GetString(), string_.GetString()))

	cooked := newString("a")
	template := Node_NewTemplateLitString(&cooked, newString("b"))
	assertTrue(t, "clone getCookedString same instance",
		sameInstance(*template.CloneNode().GetCookedString(), *template.GetCookedString()))
	assertTrue(t, "clone getRawString same instance",
		sameInstance(template.CloneNode().GetRawString(), template.GetRawString()))

	bigint := Node_NewBigInt(jutil.NewBigIntegerFromString("100"))
	assertTrue(t, "clone getBigInt same instance", bigint.CloneNode().GetBigInt() == bigint.GetBigInt())
}

func TestNode_SharedProps1(t *testing.T) {
	n := getCall("A")
	n.SetSideEffectFlags(5)
	m := NewNode(Token_TRUE)
	m.ClonePropsFrom(n)
	if n.GetPropListHeadForTesting() != m.GetPropListHeadForTesting() {
		t.Error("prop list heads differ")
	}
	assertSideEffectFlags(t, n, 5)
	assertSideEffectFlags(t, m, 5)
}

func assertSideEffectFlags(t *testing.T, n Node, want int32) {
	t.Helper()
	if got := n.GetSideEffectFlags(); got != want {
		t.Errorf("getSideEffectFlags() = %d, want %d", got, want)
	}
}

func TestNode_SharedProps2(t *testing.T) {
	n := getCall("A")
	n.SetSideEffectFlags(5)
	m := getCall("B")
	m.ClonePropsFrom(n)

	n.SetSideEffectFlags(6)
	assertSideEffectFlags(t, n, 6)
	assertSideEffectFlags(t, m, 5)
	assertFalse(t, "prop list heads identical", m.GetPropListHeadForTesting() == n.GetPropListHeadForTesting())

	m.SetSideEffectFlags(7)
	assertSideEffectFlags(t, n, 6)
	assertSideEffectFlags(t, m, 7)
}

func TestNode_SharedProps3(t *testing.T) {
	n := getCall("A")
	n.SetSideEffectFlags(2)
	n.PutBooleanProp(Node_INCRDECR_PROP, true)
	m := NewNode(Token_TRUE)
	m.ClonePropsFrom(n)
	n.SetSideEffectFlags(4)

	assertSideEffectFlags(t, n, 4)
	assertSideEffectFlags(t, m, 2)
}

func TestNode_BooleanProp(t *testing.T) {
	n := getVarRef("a")

	n.PutBooleanProp(Node_IS_CONSTANT_NAME, false)

	assertTrue(t, "lookupProperty is null", n.LookupProperty(Node_IS_CONSTANT_NAME) == nil)
	assertFalse(t, "IS_CONSTANT_NAME", n.GetBooleanProp(Node_IS_CONSTANT_NAME))

	n.PutBooleanProp(Node_IS_CONSTANT_NAME, true)

	assertTrue(t, "lookupProperty is not null", n.LookupProperty(Node_IS_CONSTANT_NAME) != nil)
	assertTrue(t, "IS_CONSTANT_NAME", n.GetBooleanProp(Node_IS_CONSTANT_NAME))

	n.PutBooleanProp(Node_IS_CONSTANT_NAME, false)

	assertTrue(t, "lookupProperty is null", n.LookupProperty(Node_IS_CONSTANT_NAME) == nil)
	assertFalse(t, "IS_CONSTANT_NAME", n.GetBooleanProp(Node_IS_CONSTANT_NAME))
}

// Verify that annotations on cloned nodes are properly handled.
func TestNode_CloneAnnontations2(t *testing.T) {
	n := getVarRef("a")
	n.PutBooleanProp(Node_IS_CONSTANT_NAME, true)
	assertTrue(t, "IS_CONSTANT_NAME", n.GetBooleanProp(Node_IS_CONSTANT_NAME))

	nodeClone := n.CloneNode()
	assertTrue(t, "clone IS_CONSTANT_NAME", nodeClone.GetBooleanProp(Node_IS_CONSTANT_NAME))

	assertTrue(t, "IS_CONSTANT_NAME", n.GetBooleanProp(Node_IS_CONSTANT_NAME))

	assertTrue(t, "clone IS_CONSTANT_NAME", nodeClone.GetBooleanProp(Node_IS_CONSTANT_NAME))
}

func bitsetFromNodeProperties(props ...NodeProperty) int64 {
	var bitset int64
	for _, prop := range props {
		bitset = Node_SetNodePropertyBit(bitset, prop)
	}
	return bitset
}

func assertSerialized(t *testing.T, got, want int64) {
	t.Helper()
	if got != want {
		t.Errorf("serializeProperties() = %#x, want %#x", got, want)
	}
}

func TestNode_SerializeProperties(t *testing.T) {
	node := IR_Function(IR_Name(""), IR_ParamList(), IR_Block())
	node.SetIsAsyncFunction(true)
	node.SetIsGeneratorFunction(true)
	result := node.SerializeProperties()

	assertSerialized(t, result,
		bitsetFromNodeProperties(NodeProperty_GENERATOR_FN, NodeProperty_ASYNC_FN))
}

func TestNode_SerializeProperties_IsDeclaredConstant(t *testing.T) {
	node := NewNode(Token_NAME)
	node.SetDeclaredConstantVar(true)
	result := node.SerializeProperties()
	assertSerialized(t, result, bitsetFromNodeProperties(NodeProperty_IS_DECLARED_CONSTANT))
}

func TestNode_SerializeProperties_IsInferredConstant(t *testing.T) {
	node := NewNode(Token_NAME)
	node.SetInferredConstantVar(true)
	result := node.SerializeProperties()
	assertSerialized(t, result, bitsetFromNodeProperties(NodeProperty_IS_INFERRED_CONSTANT))
}

func TestNode_SerializeProperties_UntranslatableRhinoProp(t *testing.T) {
	node := getCall("A")
	node.SetUseStrict(true)
	result := node.SerializeProperties()
	// Rhino node prop USE_STRICT does not have a corresponding NodeProperty
	assertTrue(t, "isUseStrict", node.IsUseStrict())
	assertSerialized(t, result, 0)
}

func TestNode_SerializeProperties_PrivateIdentifier(t *testing.T) {
	node := Node_NewStringWithToken(Token_NAME, "#field")
	node.SetPrivateIdentifier()
	result := node.SerializeProperties()
	assertSerialized(t, result, bitsetFromNodeProperties(NodeProperty_PRIVATE_IDENTIFIER))
}

func TestNode_SerializeProperties_PrivateIdentifierRoundTrip(t *testing.T) {
	original := Node_NewStringWithToken(Token_NAME, "#field")
	original.SetSourceFileForTesting("sourcefile")
	restored := original.CloneNode() // clone BEFORE setPrivateIdentifier

	assertFalse(t, "restored.isPrivateIdentifier()", restored.IsPrivateIdentifier())

	original.SetPrivateIdentifier()
	restored.DeserializeProperties(original.SerializeProperties(), false)

	assertTrue(t, "restored.isPrivateIdentifier()", restored.IsPrivateIdentifier())
}

func TestNode_GetIndexOfChild(t *testing.T) {
	assign := getAssignExpr("b", "c")
	if got := assign.GetChildCount(); got != 2 {
		t.Fatalf("child count = %d, want 2", got)
	}

	firstChild := assign.GetFirstChild()
	secondChild := firstChild.GetNext()
	if secondChild.IsNil() {
		t.Fatal("second child is null")
	}

	for _, tc := range []struct {
		child Node
		want  int
	}{{firstChild, 0}, {secondChild, 1}, {assign, -1}} {
		if got := assign.GetIndexOfChild(tc.child); got != tc.want {
			t.Errorf("getIndexOfChild = %d, want %d", got, tc.want)
		}
	}
}

func assertLinenoAndSourceFile(t *testing.T, n Node, lineno int, file string) {
	t.Helper()
	if got := n.GetLineno(); got != lineno {
		t.Errorf("lineno = %d, want %d", got, lineno)
	}
	if got := n.GetSourceFileName(); got == nil || *got != file {
		t.Errorf("source file = %v, want %q", jutil.ToJavaString(got), file)
	}
}

func TestNode_SrcrefIfMissing(t *testing.T) {
	assign := getAssignExpr("b", "c")
	assign.SetLinenoCharno(99, 0)
	assign.SetSourceFileForTesting("foo.js")

	lhs := assign.GetFirstChild()
	lhs.SrcrefIfMissing(assign)
	assertLinenoAndSourceFile(t, lhs, 99, "foo.js")

	assign.SetLinenoCharno(101, 0)
	assign.SetSourceFileForTesting("bar.js")
	lhs.SrcrefIfMissing(assign)
	assertLinenoAndSourceFile(t, lhs, 99, "foo.js")
}

func TestNode_Srcref(t *testing.T) {
	assign := getAssignExpr("b", "c")
	assign.SetLinenoCharno(99, 0)
	assign.SetSourceFileForTesting("foo.js")

	lhs := assign.GetFirstChild()
	lhs.Srcref(assign)
	assertLinenoAndSourceFile(t, lhs, 99, "foo.js")

	assign.SetLinenoCharno(101, 0)
	assign.SetSourceFileForTesting("bar.js")
	lhs.Srcref(assign)
	assertLinenoAndSourceFile(t, lhs, 101, "bar.js")
}

func TestNode_InvalidSourceOffset(t *testing.T) {
	string_ := Node_NewString("a")

	string_.SetLinenoCharno(-1, -1)
	if got := string_.GetSourceOffset(); got >= 0 {
		t.Errorf("source offset = %d, want < 0", got)
	}

	string_.SetSourceFileForTesting("foo.js")
	if got := string_.GetSourceOffset(); got >= 0 {
		t.Errorf("source offset = %d, want < 0", got)
	}
}

func TestNode_QualifiedName(t *testing.T) {
	if got := IR_Name("").GetQualifiedName(); got != nil {
		t.Errorf("qualified name of '' = %q, want null", *got)
	}
	for _, tc := range []struct {
		n    Node
		want string
	}{
		{IR_Name("a"), "a"},
		{IR_ThisNode(), "this"},
		{IR_SuperNode(), "super"},
		{IR_Getprop(IR_Name("a"), "b"), "a.b"},
		{IR_Getprop(IR_ThisNode(), "b"), "this.b"},
		{IR_Getprop(IR_SuperNode(), "b"), "super.b"},
	} {
		if got := tc.n.GetQualifiedName(); got == nil || *got != tc.want {
			t.Errorf("getQualifiedName() = %v, want %q", jutil.ToJavaString(got), tc.want)
		}
	}
	if got := IR_Getprop(IR_Call(IR_Name("a")), "b").GetQualifiedName(); got != nil {
		t.Errorf("qualified name of a().b = %q, want null", *got)
	}
}

func TestNode_JSDocInfoClone(t *testing.T) {
	original := IR_Var(IR_Name("varName"))
	builder := JSDocInfo_Builder_()
	builder.RecordType(NewJSTypeExpression(IR_Name("TypeName"), "blah"))
	info := builder.Build()
	original.GetFirstChild().SetJSDocInfo(info)

	// By default the JSDocInfo and JSTypeExpression objects are not cloned
	clone := original.CloneTree()
	assertTrue(t, "same JSDocInfo",
		clone.GetFirstChild().GetJSDocInfo() == original.GetFirstChild().GetJSDocInfo())
	assertTrue(t, "same type",
		clone.GetFirstChild().GetJSDocInfo().GetType() == original.GetFirstChild().GetJSDocInfo().GetType())
	assertTrue(t, "same type root",
		clone.GetFirstChild().GetJSDocInfo().GetType().GetRoot() ==
			original.GetFirstChild().GetJSDocInfo().GetType().GetRoot())

	// If requested the JSDocInfo and JSTypeExpression objects are cloned.
	// This is required because compiler classes are modifying the type
	// expressions in place
	clone = original.CloneTreeWithTypeExprs(true)
	assertFalse(t, "same JSDocInfo",
		clone.GetFirstChild().GetJSDocInfo() == original.GetFirstChild().GetJSDocInfo())
	assertFalse(t, "same type",
		clone.GetFirstChild().GetJSDocInfo().GetType() == original.GetFirstChild().GetJSDocInfo().GetType())
	assertFalse(t, "same type root",
		clone.GetFirstChild().GetJSDocInfo().GetType().GetRoot() ==
			original.GetFirstChild().GetJSDocInfo().GetType().GetRoot())
}

func TestNode_AddChildToFrontWithSingleNode(t *testing.T) {
	root := NewNode(Token_SCRIPT)
	nodeToAdd := NewNode(Token_SCRIPT)

	root.AddChildToFront(nodeToAdd)

	assertTrue(t, "parent", nodeToAdd.GetParent() == root)
	assertTrue(t, "first child", nodeToAdd == root.GetFirstChild())
	assertTrue(t, "last child", nodeToAdd == root.GetLastChild())
	assertTrue(t, "next is null", nodeToAdd.GetNext().IsNil())
}

func TestNode_AddChildToFrontWithLargerTree(t *testing.T) {
	left := Node_NewString("left")
	mid := Node_NewString("mid")
	right := Node_NewString("right")
	root := NewNodeWith3Children(Token_SCRIPT, left, mid, right)
	nodeToAdd := NewNode(Token_SCRIPT)

	root.AddChildToFront(nodeToAdd)

	assertTrue(t, "parent", nodeToAdd.GetParent() == root)
	assertTrue(t, "first child", nodeToAdd == root.GetFirstChild())
	assertTrue(t, "previous is null", nodeToAdd.GetPrevious().IsNil())
	assertTrue(t, "next is left", nodeToAdd.GetNext() == left)
	assertTrue(t, "left.previous", left.GetPrevious() == nodeToAdd)
}

func TestNode_Detach1(t *testing.T) {
	left := Node_NewString("left")
	mid := Node_NewString("mid")
	right := Node_NewString("right")
	root := NewNodeWith3Children(Token_SCRIPT, left, mid, right)

	assertTrue(t, "mid.parent", mid.GetParent() == root)
	assertTrue(t, "mid.previous", mid.GetPrevious() == left)
	assertTrue(t, "mid.next", mid.GetNext() == right)

	mid.Detach()

	assertTrue(t, "mid.parent is null", mid.GetParent().IsNil())
	assertTrue(t, "mid.next is null", mid.GetNext().IsNil())

	assertTrue(t, "right.previous", right.GetPrevious() == left)
	assertTrue(t, "left.next", left.GetNext() == right)
}

func TestNode_GetAncestors(t *testing.T) {
	grandparent := NewNode(Token_ROOT)
	parent := NewNode(Token_PLACEHOLDER1)
	node := NewNode(Token_PLACEHOLDER2)

	grandparent.AddChildToFront(parent)
	parent.AddChildToFront(node)

	got := slices.Collect(node.GetAncestors().All())
	if len(got) != 2 || !slices.Contains(got, parent) || !slices.Contains(got, grandparent) {
		t.Errorf("ancestors = %v, want exactly [parent, grandparent]", got)
	}
}

func TestNode_GetAncestors_Empty(t *testing.T) {
	node := NewNode(Token_ROOT)
	if got := slices.Collect(node.GetAncestors().All()); len(got) != 0 {
		t.Errorf("ancestors = %v, want empty", got)
	}
}

func TestNode_TrailingComma(t *testing.T) {
	list := NewNode(Token_ARRAYLIT)
	list.SetTrailingComma(true)
	assertTrue(t, "hasTrailingComma", list.HasTrailingComma())
}

func getVarRef(name string) Node {
	return Node_NewStringWithToken(Token_NAME, name)
}

func getAssignExpr(name1 string, name2 string) Node {
	return NewNodeWith2Children(Token_ASSIGN, getVarRef(name1), getVarRef(name2))
}

func getCall(name1 string) Node {
	return NewNodeWithChild(Token_CALL, getVarRef(name1))
}

func TestNode_ValidateMemoryGuarantees_ChecksForDifferentSourceFileInstances(t *testing.T) {
	b1 := IR_Block()
	b2 := IR_Block()
	other := IR_Block()

	b1.SetSourceFileForTesting("file.js")
	b2.SetStaticSourceFileFrom(b1)
	other.SetStaticSourceFile(b1.GetStaticSourceFile())

	assertTrue(t, "b1 vs b2", Node_ValidateMemorySensitivePropertyGuarantees(b1, false, b2, false))
	assertFalse(t, "b1 vs other", Node_ValidateMemorySensitivePropertyGuarantees(b1, false, other, false))
}

func TestNode_IsPos(t *testing.T) {
	assertTrue(t, "POS.isPos()", NewNode(Token_POS).IsPos())
	assertFalse(t, "NEG.isPos()", NewNode(Token_NEG).IsPos())
	assertFalse(t, "1.0.isPos()", Node_NewNumber(1.0).IsPos())
}
